package strategydata

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"time"
)

const (
	StrategyTypeInternal = 1
	StrategyTypeMilitary = 2

	StrategyTargetSelf          = 1
	StrategyTargetEnemy         = 2
	StrategyTargetGuild         = 3
	StrategyTargetGuildExceptMe = 4
)

//gogen:config
type StrategyData struct {
	_ struct{} `file:"策略/策略.txt"`
	_ struct{} `proto:"shared_proto.StrategyDataProto"`
	_ struct{} `protoconfig:"Strategy"`

	Id              uint64                                      // 策略ID
	Type            uint64                                      // 策略类型: 1.内政 2.军事
	Target          uint64                                      // 目标类型: 1.仅自己，2.仅敌人，3.友方，4.仅盟友
	Name            string                                      // 策略名称
	Sp              uint64 `validator:"uint" default:"0"`       // 体力值消耗
	UnlockHeroLevel uint64                                      // 解锁等级（君主）
	Cd              time.Duration                               // 冷却时间
	TodayLimit      uint64                                      // 今日使用次限
	Icon            *icon.Icon `protofield:"Icon,%s.Id,string"` // 图标
	Desc            string                                      // 描述

	EffectMap map[uint64]*StrategyEffectData `head:"-" protofield:"-"`
}

//gogen:config
type StrategyEffectData struct {
	_ struct{} `file:"策略/策略效果.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"domestic.proto"`

	Id uint64 `protofield:"-"`

	StrategyId uint64 `desc:"策略 id"`

	IntEffectType int `protofield:"-"`

	EffectType shared_proto.StrategyEffectType `desc:"策略效果类型" head:"-"`

	HeroLevel uint64 `desc:"君主等级"`

	Cost *domestic_data.CombineCost `desc:"消耗"`

	Prize *resdata.Prize `desc:"策略数值收益"`

	FarmFastHarvestDuration time.Duration `desc:"农场减少时间"`

	TargetReduceSolider uint64 `desc:"对方减兵" validator:"uint"`
}

func (d *StrategyEffectData) InitAll(filename string, conf interface {
	GetStrategyEffectDataArray() []*StrategyEffectData
	GetStrategyData(uint64) *StrategyData
	GetHeroLevelData(uint64) *herodata.HeroLevelData
	GetHeroLevelDataArray() []*herodata.HeroLevelData
}) {
	allMap := make(map[uint64]map[uint64]*StrategyEffectData)
	for _, d := range conf.GetStrategyEffectDataArray() {
		strategy := conf.GetStrategyData(d.StrategyId)
		check.PanicNotTrue(strategy != nil, "%v, strategy_id 不存在. strategy_id:%v id:%v", filename, d.StrategyId, d.Id)

		d.EffectType = shared_proto.StrategyEffectType(d.IntEffectType)
		switch d.EffectType {
		case shared_proto.StrategyEffectType_Strategy_ET_prize:

			check.PanicNotTrue(d.Prize != nil && d.Cost != nil, "%v, 数值奖励类型必须配 prize 和 cost。id:%v", filename, d.Id)

			check.PanicNotTrue(strategy.Target == StrategyTargetSelf ||
				strategy.Target == StrategyTargetGuild ||
				strategy.Target == StrategyTargetGuildExceptMe,
				"%v, 配置的是奖励效果，但是策略释放目标是敌人?（给敌人加奖励），策略: ", filename, d.Id, strategy.Id)
		case shared_proto.StrategyEffectType_Strategy_ET_farm:

			check.PanicNotTrue(d.FarmFastHarvestDuration > 0 && d.Cost != nil, "%v, 农场类型必须配 farm_fast_harvest_duration 和 cost。id:%v", filename, d.Id)

			check.PanicNotTrue(strategy.Target == StrategyTargetSelf, "%v, 配置的是农场效果，但是策略释放目标不是自己，策略: ", filename, d.Id, strategy.Id)
		case shared_proto.StrategyEffectType_Strategy_ET_random_fast_move:

			check.PanicNotTrue(strategy.Target == StrategyTargetSelf, "%v, 配置的是随机迁城效果，但是策略释放目标不是自己，策略: ", filename, d.Id, strategy.Id)
		case shared_proto.StrategyEffectType_Strategy_ET_reduce_solider:

			check.PanicNotTrue(d.TargetReduceSolider > 0 && d.Cost != nil, "%v, 减兵类型必须配  target_reduce_solider 和 cost。id:%v", filename, d.Id)

			check.PanicNotTrue(strategy.Target == StrategyTargetEnemy, "%v, 配置的是扣目标兵力效果，但是策略释放目标不是敌人，策略: ", filename, d.Id, strategy.Id)
		case shared_proto.StrategyEffectType_Strategy_ET_baoz:

			check.PanicNotTrue(strategy.Target == StrategyTargetSelf, "%v, 配置的是召唤殷墟效果，但是策略释放目标不是自己，策略: ", filename, d.Id, strategy.Id)
		default:
			logrus.Panicf("%v, 找不到 策略效果类型：%v", filename, d.IntEffectType)
		}

		if d.Prize != nil && d.Prize.Prosperity > 0 {
			check.PanicNotTrue(d.Prize.TypeCount() == 1, "%v, 有繁荣度的奖励会超出返还，所以奖励只能配繁荣度。id:%v", filename, d.Id)
			if d.Cost.Soldier > 0 {
				check.PanicNotTrue(d.Cost.Cost.TypeCount() == 0, "%v, 有繁荣度的奖励会超出返还，所以消耗只能配一种。id:%v", filename, d.Id)
			} else {
				check.PanicNotTrue(d.Cost.Cost.TypeCount() == 1, "%v, 有繁荣度的奖励会超出返还，所以消耗只能配一种。id:%v", filename, d.Id)
			}
		}

		check.PanicNotTrue(conf.GetHeroLevelData(d.HeroLevel) != nil, "%v, hero_level 不存在. hero_level:%v id:%v", filename, d.HeroLevel, d.Id)

		m, ok := allMap[d.StrategyId]
		if !ok {
			m = make(map[uint64]*StrategyEffectData)
			allMap[d.StrategyId] = m
		}

		if _, ok := m[d.HeroLevel]; ok {
			logrus.Panicf("%v, hero_level 重复。id:%v", filename, d.Id)
		}

		m[d.HeroLevel] = d
	}

	levelCount := len(conf.GetHeroLevelDataArray())
	for id, m := range allMap {
		check.PanicNotTrue(len(m) == levelCount, "%v, strategy_id:%v 的君主等级没有配全。", filename, id)
		conf.GetStrategyData(id).EffectMap = m
	}

}
