package regdata

import (
	"github.com/lightpaw/male7/config/basedata"
	"math/rand"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/util/sortkeys"
	"sort"
	"github.com/lightpaw/logrus"
)

//gogen:config
type JunTuanNpcData struct {
	_ struct{} `file:"地图/军团怪物.txt"`
	_ struct{} `protogen:"true"`

	// 军团怪id
	Id uint64

	// 野怪npc shared_proto.NpcBaseDataProto
	Npc *basedata.NpcBaseData `protofield:",config.U64ToI32(%s.Id),int32"`

	// 队伍数量
	TroopCount uint64

	// 分组
	Group uint64

	// 军团怪等级
	Level uint64

	// 出征所需君主等级
	RequiredHeroLevel uint64

	// 防守兵力
	totalSoldier uint64
}

func (b *JunTuanNpcData) Init(filename string) {
	check.PanicNotTrue(b.Id <= npcid.NpcDataMask, "%s npc城池的配置数据的id最大不能超过 %d, id: %d", filename, npcid.NpcDataMask, b.Id)

	b.Npc.DestroyWhenLose = true // TODO
	check.PanicNotTrue(b.Npc.DestroyWhenLose, "%s npc城池[%v]的怪物必须是击破流亡的，DestroyWhenLose必须设置为1", filename, b.Id)

	var totalSoldier uint64
	for _, captain := range b.Npc.Npc.Captains {
		totalSoldier += captain.Soldier
	}

	b.totalSoldier = totalSoldier * b.TroopCount

}

func (b *JunTuanNpcData) GetTotalSoldier() uint64 {
	return b.totalSoldier
}

//gogen:config
type JunTuanNpcPlaceConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"地图/军团怪物.txt"`

	groups []*JunTuanNpcPlaceGroup

	firstGroup *JunTuanNpcPlaceGroup

	lastGroup *JunTuanNpcPlaceGroup
}

func (c *JunTuanNpcPlaceConfig) Init(filename string, config interface {
	GetJunTuanNpcPlaceDataArray() []*JunTuanNpcPlaceData
}) {

	dayMap := make(map[uint64][]*JunTuanNpcPlaceData)
	for _, d := range config.GetJunTuanNpcPlaceDataArray() {
		dayMap[d.Day] = append(dayMap[d.Day], d)
	}

	var dayArray []uint64
	for k := range dayMap {
		dayArray = append(dayArray, k)
	}
	sortkeys.Uint64s(dayArray)

	c.groups = nil
	for _, day := range dayArray {
		array := dayMap[day]

		group := &JunTuanNpcPlaceGroup{}
		group.day = day
		group.blockPlaceMap = make(map[cb.Cube]map[uint64]*JunTuanNpcPlaceData)
		group.defaultPlace = make(map[uint64]*JunTuanNpcPlaceData)

		//dataMap := make(map[uint64]*JunTuanNpcPlaceData)
		for _, data := range array {
			//dataMap[data.Group] =
			if data.Area == nil {
				if old, exist := group.defaultPlace[data.Group]; exist {
					logrus.Panicf("军团刷新野怪数据初始化失败，第[%d]天的数据，剩余类别配置了重复的group[%d]刷新, id1: %d id2: %d", day, data.Group, old.Id, data.Id)
				}

				group.defaultPlace[data.Group] = data
			} else {
				for c := range data.Area.GetValidCubeMap() {
					bpm := group.blockPlaceMap[c]
					if bpm == nil {
						bpm = make(map[uint64]*JunTuanNpcPlaceData)
						group.blockPlaceMap[c] = bpm
					}

					if old, exist := bpm[data.Group]; exist {
						x, y := c.XY()
						logrus.Panicf("军团刷新野怪数据初始化失败，第[%d]天的数据，位置[%d, %d]配置了重复的group[%d]刷新, id1: %d id2: %d", day, x, y, data.Group, old.Id, data.Id)
					}

					bpm[data.Group] = data
				}
			}
		}

		check.PanicNotTrue(len(group.defaultPlace) > 0, "军团刷新野怪数据初始化失败，第[%d]天的数据，没有配置默认刷新列表", day)

		c.groups = append(c.groups, group)
	}

	check.PanicNotTrue(len(c.groups) > 0, "军团刷新野怪数据初始化失败，c.groups is empty")

	var kvs []*sortkeys.U64KV
	for _, t := range c.groups {
		kvs = append(kvs, sortkeys.NewU64KV(t.day, t))
	}
	sort.Sort(sortkeys.U64KVSlice(kvs))

	for i, kv := range kvs {
		c.groups[i] = kv.V.(*JunTuanNpcPlaceGroup)
	}

	c.firstGroup = c.groups[0]
	c.lastGroup = c.groups[len(c.groups)-1]
}

func (c *JunTuanNpcPlaceConfig) Must(day uint64) *JunTuanNpcPlaceGroup {
	if day >= c.lastGroup.day {
		return c.lastGroup
	}

	prev := c.firstGroup
	for _, g := range c.groups {
		if day < g.day {
			// 返回前一个
			return prev
		}
		prev = g
	}
	return prev
}

// 军团怪刷新时间数据
type JunTuanNpcPlaceGroup struct {
	// 开服后N天开始刷新
	day uint64

	// 安装block的xy反向存储
	blockPlaceMap map[cb.Cube]map[uint64]*JunTuanNpcPlaceData

	// 不再blockPlaceMap中的，统一使用这个数据刷新
	defaultPlace map[uint64]*JunTuanNpcPlaceData
}

func (g *JunTuanNpcPlaceGroup) GetPlaceData(c cb.Cube) map[uint64]*JunTuanNpcPlaceData {
	if array := g.blockPlaceMap[c]; len(array) > 0 {
		return array
	}
	return g.defaultPlace
}

// 军团怪刷新数据
//gogen:config
type JunTuanNpcPlaceData struct {
	_ struct{} `file:"地图/军团怪物刷新.txt"`

	Id uint64

	// 开服第几天
	Day uint64

	// 刷新的怪物分组
	Group uint64

	// 摆放的地区
	Area *AreaData `default:"nullable"`

	// 多个Npc中选一个
	npc []*JunTuanNpcData

	// 保持地区中的个数
	KeepCount uint64
}

func (d *JunTuanNpcPlaceData) Init(filename string, config interface {
	GetJunTuanNpcDataArray() []*JunTuanNpcData
}) {

	d.npc = nil
	for _, t := range config.GetJunTuanNpcDataArray() {
		if t.Group == d.Group {
			d.npc = append(d.npc, t)
		}
	}

	check.PanicNotTrue(len(d.npc) > 0, "%s %d 没有这个分组[%d]配置Npc野怪", filename, d.Id, d.Group)
}

func (d *JunTuanNpcPlaceData) Must(idx int) *JunTuanNpcData {
	if idx >= 0 && idx < len(d.npc) {
		return d.npc[idx]
	}

	if idx < 0 {
		return d.npc[0]
	}

	return d.npc[len(d.npc)-1]
}

func (d *JunTuanNpcPlaceData) Random() *JunTuanNpcData {
	return d.Must(rand.Intn(len(d.npc)))
}
