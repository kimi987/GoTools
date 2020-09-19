package regdata

import (
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/entity/npcid"
)

//gogen:config
type RegionMultiLevelNpcData struct {
	_ struct{} `file:"地图/多等级野怪.txt"`
	_ struct{} `proto:"shared_proto.RegionMultiLevelNpcDataProto"`
	_ struct{} `protoconfig:"MultiLevelNpcData"`

	// 怪物配置id
	Id uint64

	// 怪物类型
	TypeData *RegionMultiLevelNpcTypeData `protofield:"Type,%s.Type"`

	// 可以有多个等级的怪物配表
	LevelBases []*RegionMultiLevelNpcLevelData `head:"-" protofield:"LevelBaseId,encodeMultiLevelNpcLevelData(%s)"`
	firstLevel *RegionMultiLevelNpcLevelData

	// 坐标
	OffsetBaseX uint64
	OffsetBaseY uint64

	// 所属的Region Data Id
	RegionType  shared_proto.RegionType `protofield:"-"`
	RegionLevel uint64                  `protofield:"-"`
	RegionId    uint64                  `head:"-,RegionDataID(%s.RegionType%c %s.RegionLevel)" protofield:"-"`
}

func (d *RegionMultiLevelNpcData) Init(filename string, configDatas interface {
	GetRegionMultiLevelNpcLevelDataArray() []*RegionMultiLevelNpcLevelData
}) {

	check.PanicNotTrue(d.Id <= npcid.NpcDataMask, "%s npc城池的配置数据的id最大不能超过 %d, id: %d", filename, npcid.NpcDataMask, d.Id)

	var levelBases []*RegionMultiLevelNpcLevelData
	for _, v := range configDatas.GetRegionMultiLevelNpcLevelDataArray() {
		if d.Id == v.MultiLevelNpcId {
			levelBases = append(levelBases, v)
			v.typeData = d.TypeData
		}
	}

	check.PanicNotTrue(len(levelBases) > 0, "多等级怪物[%d]的等级数据没有配置", d.Id)
	for i := 0; i < len(levelBases); i++ {
		// 这个必须从1开始，英雄判断怪物是否通过，使用这个来判断
		check.PanicNotTrue(levelBases[i].Level == uint64(i+1), "多等级怪物[%d]的等级必须从1开始配置，没找到[%d]级数据", d.Id, i+1)
	}

	d.firstLevel = levelBases[0]
	d.LevelBases = levelBases

}

func encodeMultiLevelNpcLevelData(levelBases []*RegionMultiLevelNpcLevelData) (ids []int32) {
	for _, v := range levelBases {
		ids = append(ids, u64.Int32(v.Npc.Id))
	}
	return
}

func (d *RegionMultiLevelNpcData) GetFirstLevel() *RegionMultiLevelNpcLevelData {
	return d.firstLevel
}

func (d *RegionMultiLevelNpcData) GetLevelBaseData(level uint64) *RegionMultiLevelNpcLevelData {

	index := u64.Sub(level, 1)
	n := len(d.LevelBases)
	if index < uint64(n) {
		return d.LevelBases[index]
	}

	return d.LevelBases[n-1]
}

//gogen:config
type RegionMultiLevelNpcLevelData struct {
	_  struct{} `file:"地图/多等级野怪等级.txt"`
	Id uint64   `head:"-,%s.MultiLevelNpcId*10000+%s.Level"`

	MultiLevelNpcId uint64

	Level uint64

	typeData *RegionMultiLevelNpcTypeData

	Npc *basedata.NpcBaseData

	// 持续掠夺加仇恨
	HateTickDuration time.Duration
	TickHate         uint64 `validator:"uint"` // 每个duration加多少仇恨
	FirstHate        uint64 `validator:"uint"` // 第一次加多少仇恨
	FailHate         int    `validator:"int"`  // 进攻失败，仇恨值变化，可能是负数

	// 出征怪物
	FightNpc *monsterdata.MonsterMasterData
	//FightNpcPrize      *resdata.Prize           `default:"nullable"`
	//FightNpcPrizeProto *shared_proto.PrizeProto `head:"-" protofield:"-"`
}

func (d *RegionMultiLevelNpcLevelData) Init(filename string) {

	//if d.FightNpcPrize != nil {
	//	d.FightNpcPrizeProto = d.FightNpcPrize.Encode4Init()
	//}
}

func (d *RegionMultiLevelNpcLevelData) TypeData() *RegionMultiLevelNpcTypeData {
	return d.typeData
}

//gogen:config
type RegionMultiLevelNpcTypeData struct {
	_ struct{} `file:"地图/多等级野怪类型.txt"`
	_ struct{} `proto:"shared_proto.RegionMultiLevelNpcTypeProto"`
	_ struct{} `protoconfig:"MultiLevelNpcType"`

	Id int `head:"-,int(%s.Type)" protofield:"-"`

	// 名称
	Name *i18n.I18nRef

	// 类型
	Type shared_proto.MultiLevelNpcType

	// 初始仇恨值
	InitHate uint64 `validator:"uint" protofield:"-"`

	// 仇恨最大值
	MaxHate uint64

	// 怪物攻城仇恨值，>=这个仇恨值，怪物会出兵
	FightHate uint64

	// 怪物攻城后减少多少仇恨值
	FightReduceHate uint64 `protofield:"-"`

	// 怪物攻城，入侵怪物最小距离，必须>=这个值才能选择这个作为开始位置
	FightMustDistance uint64 `protofield:"-"`

	// 怪物延迟多久，开始入侵
	FightDelay time.Duration `default:"30s" protofield:"-"`
}

func (d *RegionMultiLevelNpcTypeData) Init(filename string) {

	check.PanicNotTrue(d.InitHate <= d.MaxHate, "野怪类型 %s 配置的初始仇恨值（%d） > 最大仇恨值（%d）", d.Type, d.InitHate, d.MaxHate)
	check.PanicNotTrue(d.FightHate <= d.MaxHate, "野怪类型 %s 配置的怪物攻城仇恨值（%d） > 最大仇恨值（%d）", d.Type, d.InitHate, d.MaxHate)

}
