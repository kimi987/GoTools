package mingcdata

import (
	"github.com/lightpaw/male7/config/military_data"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/check"
	"time"
)

//gogen:config
type McBuildMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"名城营建/营建杂项.txt"`
	_ struct{} `protogen:"true"`

	MaxDailyAddSupport       uint64        `desc:"名城每日最高可加民心"`
	DailyReduceSupport       uint64        `desc:"名城每日下降的民心"`
	BuildCd                  time.Duration `desc:"营建CD"`
	DailyBuildMaxCount       uint64        `desc:"每日营建最大次数"`
	BuildMinHeroLevel        uint64        `desc:"可以营建的最小玩家等级"`
	MaxRecommendMcGuildCount int           `desc:"最多可推荐名城数量" validator:"int>0"`
	MaxMcBuildLogCount       int           `desc:"营建名城记录最大排名数量" validator:"int>0"`
}

//gogen:config
type McBuildGuildMemberPrizeData struct {
	_ struct{} `file:"名城营建/联盟成员奖励.txt"`
	_ struct{} `protogen:"true"`

	Id             uint64    `protofield:"-"`
	MinBuildCount  uint64    `desc:"最小营建次数"`
	MaxBuildCount  uint64    `desc:"最大营建次数"`
	BuildCountRage *U64Range `head:"-" protofield:"-"`

	Prize *resdata.Prize `desc:"营建奖励"`
}

func (mcbgmpdata *McBuildGuildMemberPrizeData) InitAll(filename string, conf interface {
	GetMcBuildGuildMemberPrizeDataArray() []*McBuildGuildMemberPrizeData
}) {
	var prev *McBuildGuildMemberPrizeData
	for _, d := range conf.GetMcBuildGuildMemberPrizeDataArray() {
		check.PanicNotTrue(d.MinBuildCount < d.MaxBuildCount, "%v, id:%v 最小营建次数必须小于最大营建次数", filename, d.Id)
		d.BuildCountRage = NewU64Range(d.MinBuildCount, d.MaxBuildCount)
		if prev != nil {
			check.PanicNotTrue(d.MinBuildCount == prev.MaxBuildCount+1, "%v, id:%v 最小营建次数必须==上一条最大营建次数+1", filename, d.Id)
		}
		prev = d
	}
}

//gogen:config
type McBuildMcSupportData struct {
	_ struct{} `file:"名城营建/名城民心.txt"`
	_ struct{} `protogen:"true"`

	Level          uint64 `key:"true" desc:"民心等级"`
	UpgradeSupport uint64 `desc:"升到下一级需要的民心"`

	AddDailyYinliang     uint64 `desc:"增加的名城每日新增银两"`
	AddMaxYinliang       uint64 `desc:"增加的名城仓库银两上限"`
	AddHostDailyYinliang uint64 `desc:"增加的占领盟每日收益银两"`

	PrevLevelData *McBuildMcSupportData `head:"-" protofield:"-"`
	NextLevelData *McBuildMcSupportData `head:"-" protofield:"-"`
}

func (mcbgmpdata *McBuildMcSupportData) InitAll(filename string, conf interface {
	GetMcBuildMcSupportDataArray() []*McBuildMcSupportData
}) {
	var prev *McBuildMcSupportData
	for _, d := range conf.GetMcBuildMcSupportDataArray() {
		if prev != nil {
			check.PanicNotTrue(d.Level == prev.Level+1, "%v, level:%v 从1开始必须依次递增", filename, d.Level)
			prev.NextLevelData = d
		} else {
			check.PanicNotTrue(d.Level == 1, "%v, level:%v 从1开始必须依次递增", filename, d.Level)
		}
		d.PrevLevelData = prev

		prev = d
	}
}

//gogen:config
type McBuildAddSupportData struct {
	_ struct{} `file:"名城营建/增加民心.txt"`
	_ struct{} `protogen:"true"`

	BaiZhanLevel uint64 `key:"true" desc:"百战军衔"`
	AddSupport   uint64 `desc:"增加的民心"`
}

func (d *McBuildAddSupportData) Init(filename string, conf interface {
	GetJunYingLevelData(uint64) *military_data.JunYingLevelData
}) {
	check.PanicNotTrue(conf.GetJunYingLevelData(d.BaiZhanLevel) != nil, "%v, 百战军衔%v 不存在", filename, d.BaiZhanLevel)
}
