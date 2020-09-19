package farm

import (
	"time"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/config/data"
	"fmt"
)

// 农场资源
//gogen:config
type FarmResConfig struct {
	_ struct{} `file:"农场/农场资源.txt"`
	_ struct{} `proto:"shared_proto.FarmResConfigProto"`
	_ struct{} `protoconfig:"farm_res_config"`

	Id           uint64
	RipeDuration time.Duration        `parser:"time.ParseDuration"`    // 资源成熟时间间隔
	ResType      shared_proto.ResType `validator:"string" type:"enum"` // 资源类型
	BaseOutput   *data.Amount         `parser:"data.ParseAmount"`      // 基础产量
	Icon         string                                                // 资源图标
}

func GetResType(intType uint64) shared_proto.ResType {
	switch intType {
	case 1:
		return shared_proto.ResType_GOLD
	case 4:
		return shared_proto.ResType_STONE
	default:
		return 0
	}
}

//gogen:config
type FarmMaxStealConfig struct {
	_ struct{} `file:"农场/农场偷菜上限.txt"`

	GuanFuLevel              uint64 `validator:"uint" key:"true" desc:"官府等级"`
	MaxDailyStealGoldAmount  uint64 `validator:"uint" desc:"每天偷铜币最大上限"`
	MaxDailyStealStoneAmount uint64 `validator:"uint" desc:"每天偷石料最大上限"`
}

func (c *FarmMaxStealConfig) InitAll(filename string, array []*FarmMaxStealConfig, config interface{
	GetBuildingDataArray() []*domestic_data.BuildingData
}) {
	for _, b := range config.GetBuildingDataArray() {
		if b.Type != shared_proto.BuildingType_GUAN_FU {
			continue
		}

		var exist bool
		for _, a := range array {
			if b.Level == a.GuanFuLevel {
				exist = true
				break
			}
		}
		check.PanicNotTrue(exist, "%v 找不到官府%v级的农场推荐种植数据", filename, b.Level)
	}
}

// 农场一键种植
//gogen:config
type FarmOneKeyConfig struct {
	_ struct{} `file:"农场/农场一键种植.txt"`

	BaseLevel          uint64 `validator:"uint" key:"true"` // 主城等级
	StoneHopeCubeCount uint64 `validator:"uint"`            // 石料期望数量
	GoldHopeCubeCount  uint64 `validator:"uint"`            // 铜币期望数量

	MaxDailyStealGoldAmount uint64 `validator:"uint" head:"-" protofield:"-"`            // 石料期望数量
	MaxDailyStealStoneAmount  uint64 `validator:"uint" head:"-" protofield:"-"`            // 铜币期望数量
}

func (c *FarmOneKeyConfig) InitAll(filename string, array []*FarmOneKeyConfig, config interface {
	GetBaseLevelDataArray() []*domestic_data.BaseLevelData
}) {
	for _, baseLevel := range config.GetBaseLevelDataArray() {
		var exist bool
		for _, conf := range array {
			if conf.BaseLevel == baseLevel.Level {
				exist = true
				break
			}
		}
		check.PanicNotTrue(exist, "%v 找不到主城%v级的农场推荐种植数据", filename, baseLevel.Level)
	}
}

// 农场杂项
//gogen:config
type FarmMiscConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"农场/农场杂项.txt"`
	_ struct{} `proto:"shared_proto.FarmMiscConfigProto"`
	_ struct{} `protoconfig:"farm_misc_config"`

	EarlyHarvestPercent *data.Amount `parser:"data.ParseAmount"` // 提前收获千分比

	RipeProtectDuration     time.Duration // 成熟保护时间
	NegRipeProtectDuration  time.Duration `head:"-" protofield:"-"`
	StealGainPercent        *data.Amount  `parser:"data.ParseAmount"` // 偷菜收获千分比
	CubeStealMaxTime        uint64        `validator:"uint"`          // 每块地每次最多被偷次数
	StealLogMaxCount        uint64        `validator:"uint"`
	StealLogExpiredHours    string        `protofield:"-"`
	StealLogExpiredDuration time.Duration `head:"-" protofield:"-"` // 偷菜日志过期时间

	OneKeyResConfig map[shared_proto.ResType]*FarmResConfig `protofield:"-" head:"-"`
}

func (c *FarmMiscConfig) Init(filename string, configs interface {
	GetFarmResConfigArray() []*FarmResConfig
}) {
	logExpiredDuration, err := time.ParseDuration(fmt.Sprint("-", c.StealLogExpiredHours))
	check.PanicNotTrue(err == nil && logExpiredDuration < 0, "农场日志过期时间配置错误 %v", c.StealLogExpiredHours)
	c.StealLogExpiredDuration = logExpiredDuration
	check.PanicNotTrue(c.StealLogMaxCount <= 100, "农场日志最大条数不能大于100. %v", c.StealLogMaxCount)
	c.NegRipeProtectDuration = -c.RipeProtectDuration

	c.OneKeyResConfig = make(map[shared_proto.ResType]*FarmResConfig)
	for _, resConf := range configs.GetFarmResConfigArray() {
		curConf := c.OneKeyResConfig[resConf.ResType]
		if curConf == nil ||
			curConf.RipeDuration > resConf.RipeDuration {
			c.OneKeyResConfig[resConf.ResType] = resConf
		}
	}
	check.PanicNotTrue(
		c.OneKeyResConfig[shared_proto.ResType_GOLD] != nil &&
			c.OneKeyResConfig[shared_proto.ResType_STONE] != nil,
		"农场的铜币和石料资源一条都没配置")
}
