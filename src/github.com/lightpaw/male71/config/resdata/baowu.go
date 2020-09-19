package resdata

import (
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/pb/shared_proto"
	"time"
	"github.com/lightpaw/male7/util/check"
)

func BaoDataId(group, level uint64) uint64 {
	return group*10000 + level
}

//gogen:config
type BaowuData struct {
	_ struct{} `file:"物品/宝物.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"goods.proto"`

	// 宝物id
	Id uint64 `desc:"宝物id" head:"-,BaoDataId(%s.Group%c %s.Level)"`

	Group uint64 `desc:"宝物类型" validator:"int"`

	Level uint64 `desc:"宝物等级"`

	Name string `desc:"宝物名字"`

	Desc string `desc:"宝物描述"`

	Icon *icon.Icon `desc:"宝物图标" protofield:",%s.Id,string"` // 图标

	Quality shared_proto.Quality `desc:"宝物品质" head:"-"` // GoodsQuality

	GoodsQuality *goods.GoodsQuality `desc:"物品品质" protofield:",config.U64ToI32(%s.Level),int32"`

	UnlockDuration time.Duration `desc:"解锁时间"`

	PlunderPrize *PlunderPrize `desc:"解锁奖励" protofield:",%s.Prize.PrizeProto(),PrizeProto"`

	UpgradeNeedCount uint64 `desc:"升下一级所需个数"`

	DecomposeGold  uint64 `desc:"分解获得的铜币" validator:"uint"`
	DecomposeStone uint64 `desc:"分解获得的石料" validator:"uint"`

	MiaoDuration time.Duration `desc:"秒 CD 间隔"`

	// 不能被抢夺
	CantRob bool `default:"false" protofield:"-"`

	Prestige uint64 // 联盟增加声望值

	nextLevel *BaowuData
}

func (d *BaowuData) GetNextLevel() *BaowuData {
	return d.nextLevel
}

func (d *BaowuData) Init(filename string, dataMap map[uint64]*BaowuData) {

	d.Quality = d.GoodsQuality.Quality

	check.PanicNotTrue(d.Level < 10000, "%s 配置的等级必须 < 10000", filename)

	check.PanicNotTrue(d.MiaoDuration >= time.Second, "%s miao_duration 至少要1秒", filename)

	if d.Level > 1 {
		prevId := BaoDataId(d.Group, d.Level-1)
		prev := dataMap[prevId]
		check.PanicNotTrue(prev != nil, "%s 配置的同一个组[%d]的宝物，等级必须从1开始连续配置", filename, d.Group)

		prev.nextLevel = d
	}

}

func BaowuLevelIsLarge(pi, pj *BaowuData) bool {
	if pi.Level != pj.Level {
		// 等级大的在前面
		return pi.Level > pj.Level
	}

	// id大的在前面
	return pi.Id > pj.Id
}

type BaowuDataLevelSlice []*BaowuData

func (p BaowuDataLevelSlice) Len() int { return len(p) }
func (p BaowuDataLevelSlice) Less(i, j int) bool {
	// 大的在前面
	return BaowuLevelIsLarge(p[i], p[j])
}
func (p BaowuDataLevelSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
