package strategy

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/strategy"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/config/strategydata"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/config/domestic_data"
	"time"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/gen/pb/farm"
	"github.com/lightpaw/male7/entity/hexagon"
)

func NewStrategyModule(dep iface.ServiceDep, datas *config.ConfigDatas, realmService iface.RealmService, farmService iface.FarmService) *StrategyModule {
	m := &StrategyModule{
		dep:          dep,
		datas:        datas,
		realmService: realmService,
		farmService:  farmService,
	}

	return m
}

//gogen:iface
type StrategyModule struct {
	dep          iface.ServiceDep
	datas        *config.ConfigDatas
	realmService iface.RealmService
	farmService  iface.FarmService
}

//gogen:iface
func (m *StrategyModule) ProcessUseStratagem(proto *strategy.C2SUseStratagemProto, hc iface.HeroController) {
	stratagemId := u64.FromInt32(proto.GetId())
	data := m.datas.GetStrategyData(stratagemId)
	if data == nil {
		logrus.Debug("施计，没有该计策")
		hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_STRATAGEM_ID)
		return
	}
	targetId, ok := idbytes.ToId(proto.GetTarget())
	if !ok || targetId == 0 {
		logrus.Debug("施计，无效的目标id")
		hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
		return
	}
	switch data.Target {
	case strategydata.StrategyTargetSelf:
		if targetId != hc.Id() {
			logrus.Debug("施计，目标必须是自己")
			hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
			return
		}
	case strategydata.StrategyTargetEnemy:
		if targetId == hc.Id() {
			logrus.Debug("施计，目标不允许是自己")
			hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
			return
		}
		if guildId, _ := hc.LockGetGuildId(); guildId != 0 {
			if target := m.dep.HeroSnapshot().Get(targetId); target == nil || target.GuildId == guildId {
				logrus.Debug("施计，目标不允许是盟友")
				hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
				return
			}
		}
	case strategydata.StrategyTargetGuild:
		if targetId != hc.Id() {
			guildId, _ := hc.LockGetGuildId();
			if guildId == 0 {
				logrus.Debug("施计，目标不是盟友")
				hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
				return
			}
			if target := m.dep.HeroSnapshot().Get(targetId); target == nil || target.GuildId != guildId {
				logrus.Debug("施计，目标不是盟友")
				hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
				return
			}
		}
	case strategydata.StrategyTargetGuildExceptMe:
		if targetId == hc.Id() {
			logrus.Debug("施计，目标不允许是自己")
			hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
			return
		}
		guildId, _ := hc.LockGetGuildId();
		if guildId == 0 {
			logrus.Debug("施计，目标不是盟友")
			hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
			return
		}
		if target := m.dep.HeroSnapshot().Get(targetId); target == nil || target.GuildId != guildId {
			logrus.Debug("施计，目标不是盟友")
			hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
			return
		}
	}

	posX := int(proto.PosX)
	posY := int(proto.PosY)

	hctx := heromodule.NewContext(m.dep, operate_type.StrategyUse)
	ctime := m.dep.Time().CurrentTime()
	var effectData *strategydata.StrategyEffectData
	var realCost *domestic_data.CombineCost
	var baozData *regdata.BaozNpcData

	var afterFunc func()
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		effectData = data.EffectMap[hero.Level()]
		if effectData == nil {
			logrus.Debug("施计，找不到对应君主等级的效果配置")
			result.Add(strategy.ERR_USE_STRATAGEM_FAIL_SERVER_ERR)
			return
		}

		if hero.Level() < data.UnlockHeroLevel {
			logrus.Debug("施计，计策未解锁")
			result.Add(strategy.ERR_USE_STRATAGEM_FAIL_LOCKED_STRATAGEM_ID)
			return
		}
		if !hero.HasEnoughSp(data.Sp) {
			logrus.Debug("施计，体力值不足")
			result.Add(strategy.ERR_USE_STRATAGEM_FAIL_SP_NOT_ENOUGH)
			return
		}
		if hero.Strategy().IsStratagemCd(stratagemId, ctime) {
			logrus.Debug("施计，计策CD")
			result.Add(strategy.ERR_USE_STRATAGEM_FAIL_STRATAGEM_CD)
			return
		}
		if hero.Strategy().GetTodayUsedTimes(stratagemId) >= data.TodayLimit {
			logrus.Debug("施计，该计策今日使用上限")
			result.Add(strategy.ERR_USE_STRATAGEM_FAIL_TIMES_LIMIT)
			return
		}
		realCost = heromodule.StrategyUsedRealCost(hero, effectData)
		if !heromodule.HasEnoughCombineCost(hero, realCost, ctime) {
			logrus.Debug("施计，消耗不够")
			result.Add(strategy.ERR_USE_STRATAGEM_FAIL_COST_NOT_ENOUGH)
			return
		}

		if effectData.EffectType == shared_proto.StrategyEffectType_Strategy_ET_baoz {
			if targetId != 0 && targetId != hero.Id() {
				logrus.Debug("施计（召唤殷墟），发送的目标不是自己")
				result.Add(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
				return
			}

			baozData = m.datas.GetBaozNpcData(u64.FromInt32(proto.DataId))
			if baozData == nil {
				logrus.Debug("施计（召唤殷墟），找不到殷墟的id")
				result.Add(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_DATA_ID)
				return
			}

			if hero.Level() < baozData.RequiredHeroLevel {
				logrus.Debugf("施计（召唤殷墟），君主等级不足")
				result.Add(strategy.ERR_USE_STRATAGEM_FAIL_LEVEL_NOT_ENOUGH)
				return
			}

			if distance := u64.FromInt(hexagon.Distance(posX, posY, hero.BaseX(), hero.BaseY()));
				distance > m.datas.MiscGenConfig().HeroBaozMaxDistance {
				logrus.Debugf("施计（召唤殷墟），距离自己的主城太远")
				result.Add(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_POS)
				return
			}
		}

		if targetId != hero.Id() {
			// 对别人释放
			if hero.Strategy().GetTodayTargetTimes(targetId) >= m.datas.MiscGenConfig().TargetUseStratagemLimit {
				logrus.Debug("施计，该目标施计次限")
				result.Add(strategy.ERR_USE_STRATAGEM_FAIL_TARGET_LIMIT)
				return
			}
		} else {
			ok, afterFunc = m.executeCommonEffect(hctx, hero, result, effectData, ctime)
			if ok {
				// 释放成功，更新数据发消息等
				heroUseStrategySuccess(hctx, hero, result, data, targetId, realCost, ctime)
			} else {
				// 没有处理函数
				switch effectData.EffectType {
				case shared_proto.StrategyEffectType_Strategy_ET_random_fast_move:
					// 释放成功，更新数据发消息等
					ok = true
					heroUseStrategySuccess(hctx, hero, result, data, targetId, realCost, ctime)

					originX, originY := hero.BaseX(), hero.BaseY()
					afterFunc = func() {
						newRealm, newX, newY := m.realmService.ReserveRandomHomePos(realmface.RPTRandom)

						m.realmService.DoMoveBase(shared_proto.GoodsMoveBaseType_MOVE_BASE_RANDOM, newRealm, hc, originX, originY, newX, newY, false)
					}

				case shared_proto.StrategyEffectType_Strategy_ET_baoz:
					// 什么都不干
				default:
					logrus.Error("施计，不存在对应的处理函数")
					result.Add(strategy.ERR_USE_STRATAGEM_FAIL_TARGET_CANNOT_EFFECT)
					return
				}
			}

		}

		result.Ok()
	}) {
		return
	}

	if afterFunc != nil {
		afterFunc()
	}

	if targetId == 0 || targetId == hc.Id() {
		if ok {
			return
		}

		if effectData.EffectType == shared_proto.StrategyEffectType_Strategy_ET_baoz {
			if baozData == nil {
				logrus.Error("施计（召唤殷墟），baozData == nil")
				hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_DATA_ID)
				return
			}

			r := m.realmService.GetBigMap()

			// 判断玩家是否已经召唤了殷墟
			if base := r.GetHeroBaozRoBase(hc.Id()); base != nil {
				logrus.Debugf("施计（召唤殷墟），已经存在召唤的殷墟")
				hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_SERVER_ERR)
				return
			}

			// 尝试申请坐标等等
			if !r.IsPosOpened(posX, posY) {
				logrus.Debugf("施计（召唤殷墟），坐标还未开放, %v,%v", posX, posY)
				hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_POS)
				return
			}

			if r.IsEdgeNotHomePos(posX, posY) {
				logrus.Debugf("施计（召唤殷墟），边界坐标, %v,%v", posX, posY)
				hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_POS)
				return
			}

			// 这里申请了野外的坐标，下面直接
			if ok := r.ReservePos(posX, posY); !ok {
				logrus.Debug("施计（召唤殷墟），申请坐标失败，位置被占用太近了")
				hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_POS)
				return
			}

			// 扣玩家次数
			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
				heroUseStrategySuccess(hctx, hero, result, data, targetId, realCost, ctime)
				result.Ok()
			})

			expireTime := ctime.Add(m.datas.MiscGenConfig().HeroBaozDuration)

			// 添加召唤的殷墟
			processed, success := r.AddHeroBaoZangMonster(baozData, posX, posY, hc.Id(), timeutil.Marshal32(expireTime))
			if !processed || !success {
				logrus.WithField("processed", processed).WithField("success", success).Error("召唤殷墟失败")
			}
		}
	} else {
		// 对别人释放

		var errMsg pbutil.Buffer
		if m.dep.HeroData().FuncWithSend(targetId, func(target *entity.Hero, r herolock.LockResult) {
			if target.Strategy().TodayTrappedTimes() >= m.datas.MiscGenConfig().TrappedStratagemLimit {
				logrus.Debug("施计（目标），该目标中计次限")
				errMsg = strategy.ERR_USE_STRATAGEM_FAIL_TARGET_TRAPPED_LIMIT
				return
			}
			if !target.Strategy().IsTrappedStratagemEnd(stratagemId, ctime) {
				logrus.Debug("施计（目标），该目标正中该计")
				errMsg = strategy.ERR_USE_STRATAGEM_FAIL_SAME_TRAPPED
				return
			}

			ok, afterFunc = m.executeCommonEffect(hctx, target, r, effectData, ctime)
			if !ok {
				logrus.Error("施计（目标），没找到处理函数")
				errMsg = strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET
				return
			}

			// 目标中计
			stratagem := target.Strategy().TrappedStratagem(stratagemId, ctime.Add(data.Cd)) // TODO: 暂时用CD代替
			r.Add(strategy.NewS2cTrappedStratagemMsg(stratagem.EncodeClient()))

			// TODO: 以后可能还需要加上前端飘字消息

			r.Ok()
		}) {
			logrus.Debug("施计，没有目标")
			hc.Send(strategy.ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)
			return
		}

		if errMsg != nil {
			hc.Send(errMsg)
			return
		}

		// 释放成功
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			heroUseStrategySuccess(hctx, hero, result, data, targetId, realCost, ctime)
			result.Ok()
		})

		if afterFunc != nil {
			afterFunc()
		}
	}

}

// 如果有很特殊的逻辑需要加进来，那么，可以在外面加个判断，单独处理掉，当然，最好还是可以放进来
func (m *StrategyModule) executeCommonEffect(hctx *heromodule.HeroContext, hero *entity.Hero, result herolock.LockResult,

	effect *strategydata.StrategyEffectData, ctime time.Time) (ok bool, afterFunc func()) {
	switch effect.EffectType {
	case shared_proto.StrategyEffectType_Strategy_ET_prize:

		heromodule.AddPrize(hctx, hero, result, effect.Prize, ctime)
		if effect.Prize.Prosperity > 0 {
			hid := hero.Id()
			region := hero.BaseRegion()
			afterFunc = func() {
				heromodule.AddProsperity(m.realmService, hid, region, effect.Prize.Prosperity)
			}
		}

		return true, afterFunc

	case shared_proto.StrategyEffectType_Strategy_ET_farm:

		if effect.FarmFastHarvestDuration > 0 {
			heroId := hero.Id()
			afterFunc = func() {
				m.farmService.ReduceRipeTime(heroId, effect.FarmFastHarvestDuration)
				m.dep.World().Send(heroId, farm.FARM_IS_UPDATE_S2C)
			}
		}

		return true, afterFunc

	case shared_proto.StrategyEffectType_Strategy_ET_reduce_solider:
		hero.Military().ReduceFreeSoldier(effect.TargetReduceSolider, ctime)
		result.Add(hero.Military().NewUpdateFreeSoldierMsg())

		return true, afterFunc
	}

	return
}

func heroUseStrategySuccess(hctx *heromodule.HeroContext, hero *entity.Hero, result herolock.LockResult,
	data *strategydata.StrategyData, targetId int64, realCost *domestic_data.CombineCost, ctime time.Time) {

	// 扣消耗
	heromodule.ReduceCombineCostAnyway(hctx, hero, result, realCost, ctime)
	hero.ReduceSp(data.Sp)
	result.Add(domestic.NewS2cUpdateSpMsg(u64.Int32(hero.GetSp())))

	// 加次数
	stratagem := hero.Strategy().UseStratagem(data.Id, targetId, ctime.Add(data.Cd))

	result.Add(strategy.NewS2cUseStratagemMsg(u64.Int32(data.Id), data.Name, u64.Int32(stratagem.DailyUsedTimes()), timeutil.Marshal32(stratagem.NextUseableTime()), hero.Name()))

	hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_StrategyUsed)
	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_STRATEGY_USE)
}

func (m *StrategyModule) GMStrategy(sid uint64, hc iface.HeroController) {
	m.ProcessUseStratagem(&strategy.C2SUseStratagemProto{Id: u64.Int32(sid), Target: idbytes.ToBytes(hc.Id())}, hc)
}
