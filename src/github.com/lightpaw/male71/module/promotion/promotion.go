package promotion

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/promotion"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
)

func NewPromotionModule(datas iface.ConfigDatas, timeService iface.TimeService, dep iface.ServiceDep, tlgService  iface.TimeLimitGiftService, guildModule iface.GuildModule) *PromotionModule {

	m := &PromotionModule{}
	m.datas = datas
	m.timeService = timeService
	m.dep = dep
	m.tlgService = tlgService
	m.guildModule = guildModule

	return m
}

//gogen:iface
type PromotionModule struct {
	datas       iface.ConfigDatas
	timeService iface.TimeService
	dep         iface.ServiceDep
	tlgService  iface.TimeLimitGiftService
	guildModule iface.GuildModule
}

//gogen:iface
func (m *PromotionModule) ProcessCollectLogin7DayPrize(proto *promotion.C2SCollectLoginDayPrizeProto, hc iface.HeroController) {

	loginData := m.datas.GetLoginDayData(u64.FromInt32(proto.Day))
	if loginData == nil {
		logrus.Debug("领取累积登陆奖励，无效的登陆天数")
		hc.Send(promotion.ERR_COLLECT_LOGIN_DAY_PRIZE_FAIL_INVALID_DAY)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.GetAccumLoginDay() < loginData.Day {
			logrus.Debug("领取累积登陆奖励，登陆天数不足")
			result.Add(promotion.ERR_COLLECT_LOGIN_DAY_PRIZE_FAIL_DAY_NOT_ENOUGH)
			return
		}

		if hero.Promotion().IsLoginDayCollected(loginData.Day) {
			logrus.Debug("领取累积登陆奖励，这天的奖励已经领取过了")
			result.Add(promotion.ERR_COLLECT_LOGIN_DAY_PRIZE_FAIL_COLLECTED)
			return
		}

		ctime := m.timeService.CurrentTime()
		hctx := heromodule.NewContext(m.dep, operate_type.PromotionCollectLoginPrize)
		heromodule.AddPrize(hctx, hero, result, loginData.Prize, ctime)

		hero.Promotion().CollectLoginDay(loginData.Day)
		result.Add(promotion.NewS2cCollectLoginDayPrizeMsg(u64.Int32(loginData.Day)))

		result.Ok()
	})

}

//gogen:iface c2s_buy_level_fund
func (m *PromotionModule) ProcessBuyHeroLevelFund(hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.PromotionBuyHeroLevelFund)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.Bools().Get(shared_proto.HeroBoolType_BOOL_BUY_LEVEL_FUND) {
			logrus.Debug("购买君主等级基金，已经购买过了")
			result.Add(promotion.ERR_BUY_LEVEL_FUND_FAIL_BUYED)
			return
		}

		// 扣钱
		if !heromodule.TryReduceCost(hctx, hero, result, m.datas.PromotionMiscData().HeroLevelFundCost) {
			logrus.Debug("购买君主等级基金，消耗不足")
			result.Add(promotion.ERR_BUY_LEVEL_FUND_FAIL_COST_NOT_ENOUGH)
			return
		}

		// 设置购买过标识
		hero.Bools().SetTrue(shared_proto.HeroBoolType_BOOL_BUY_LEVEL_FUND)
		result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_BUY_LEVEL_FUND)))

		result.Add(promotion.BUY_LEVEL_FUND_S2C)

		result.Ok()
	})

}

//gogen:iface
func (m *PromotionModule) ProcessCollectHeroLevelFund(proto *promotion.C2SCollectLevelFundProto, hc iface.HeroController) {

	fundData := m.datas.GetHeroLevelFundData(u64.FromInt32(proto.Level))
	if fundData == nil {
		logrus.Debug("领取累积登陆奖励，无效的登陆天数")
		hc.Send(promotion.ERR_COLLECT_LOGIN_DAY_PRIZE_FAIL_INVALID_DAY)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if !hero.Bools().Get(shared_proto.HeroBoolType_BOOL_BUY_LEVEL_FUND) {
			logrus.Debug("领取君主等级基金，还没有购买")
			result.Add(promotion.ERR_COLLECT_LEVEL_FUND_FAIL_NOT_BUY)
			return
		}

		if hero.Level() < fundData.Level {
			logrus.Debug("领取君主等级基金，等级不足")
			result.Add(promotion.ERR_COLLECT_LEVEL_FUND_FAIL_LEVEL_NOT_ENOUGH)
			return
		}

		if hero.Promotion().IsHeroLevelFundCollected(fundData.Level) {
			logrus.Debug("领取君主等级基金，这个等级的奖励已经领取过了")
			result.Add(promotion.ERR_COLLECT_LEVEL_FUND_FAIL_COLLECTED)
			return
		}

		ctime := m.timeService.CurrentTime()
		hctx := heromodule.NewContext(m.dep, operate_type.PromotionCollectLevelFund)
		heromodule.AddPrize(hctx, hero, result, fundData.Prize, ctime)

		hero.Promotion().CollectHeroLevelFund(fundData.Level)
		result.Add(promotion.NewS2cCollectLevelFundMsg(u64.Int32(fundData.Level)))

		result.Ok()
	})

}


//gogen:iface
func (m *PromotionModule) ProcessCollectDailySp(proto *promotion.C2SCollectDailySpProto, hc iface.HeroController) {
	id := u64.FromInt32(proto.GetId())
	data := m.datas.GetSpCollectionData(id)
	if data == nil {
		logrus.Debug("领取日常体力奖励，无效的id")
		hc.Send(promotion.ERR_COLLECT_DAILY_SP_FAIL_INVALID_ID)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Promotion().IsDailySpCollected(data.Id) {
			logrus.Debug("领取日常体力奖励，已经领取过这顿宴席了")
			result.Add(promotion.ERR_COLLECT_DAILY_SP_FAIL_COLLECTED)
			return
		}
		hctx := heromodule.NewContext(m.dep, operate_type.PromotionCollectDailySp)
		ctime := m.timeService.CurrentTime()
		midnightTime := timeutil.DailyTime.PrevTime(ctime)
		startTime := midnightTime.Add(data.StartDuration)
		endTime := midnightTime.Add(data.EndDuration)
		if ctime.Before(startTime) { // 还未开饭，vip大佬也没用
			logrus.Debug("领取日常体力奖励，还未到饭点")
			result.Add(promotion.ERR_COLLECT_DAILY_SP_FAIL_OVERDUE)
			return
		} else if ctime.After(endTime) { // 错过饭点，只有满足vip等级才可以补领
			if vipData := m.dep.Datas().GetVipLevelData(hero.VipLevel()); vipData == nil || !vipData.CollectDailySp { // vip等级不足，无法补领
				logrus.Debug("领取日常体力奖励，错过饭点")
				result.Add(promotion.ERR_COLLECT_DAILY_SP_FAIL_OVERDUE)
				return
			}
			if !heromodule.TryReduceCost(hctx, hero, result, data.RepairCost) {
				logrus.Debug("领取日常体力奖励，补领消耗不足")
				result.Add(promotion.ERR_COLLECT_DAILY_SP_FAIL_REPAIR_COST_NOT_ENOUGH)
				return
			}
		}
		hero.Promotion().CollectDailySp(data.Id)
		heromodule.AddPrize(hctx, hero, result, data.SpPrize, ctime)
		result.Add(promotion.NewS2cCollectDailySpMsg(proto.Id, u64.Int32(hero.GetSp())))
		result.Ok()
	})
}

//gogen:iface
func (m *PromotionModule) ProcessCollectFreeGift(proto *promotion.C2SCollectFreeGiftProto, hc iface.HeroController) {
	id := u64.FromInt32(proto.GetId())
	data := m.datas.GetFreeGiftData(id)
	if data == nil {
		logrus.Debug("领取免费礼包，无效的id")
		hc.Send(promotion.ERR_COLLECT_FREE_GIFT_FAIL_INVALID_ID)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Promotion().IsFreeGiftCollected(data) {
			logrus.Debug("领取免费礼包，已经领取过")
			result.Add(promotion.ERR_COLLECT_FREE_GIFT_FAIL_COLLECTED)
			return
		}
		hero.Promotion().CollectFreeGift(data)

		hctx := heromodule.NewContext(m.dep, operate_type.PromotionCollectFreeGift)
		heromodule.AddPrize(hctx, hero, result, data.Prize, m.timeService.CurrentTime())

		result.Add(promotion.NewS2cCollectFreeGiftMsg(proto.Id, must.Marshal(data.Prize.Encode())))

		result.Ok()
	})
}

// 购买定时时限礼包
//gogen:iface
func (m *PromotionModule) ProcessBuyTimeLimitGift(proto *promotion.C2SBuyTimeLimitGiftProto, hc iface.HeroController) {
	data := m.datas.GetTimeLimitGiftGroupData(u64.FromInt32(proto.GetGrpId()))
	if data == nil {
		logrus.Debug("购买时限礼包，无效的id")
		hc.Send(promotion.ERR_BUY_TIME_LIMIT_GIFT_FAIL_INVALID_ID)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		var guanfuLevel uint64
		if guanFu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU); guanFu != nil {
			guanfuLevel = guanFu.Level
		}
		if !data.IsCanBuy(hero.Level(), guanfuLevel) {
			logrus.Debug("购买时限礼包，等级限定，无法购买")
			result.Add(promotion.ERR_BUY_TIME_LIMIT_GIFT_FAIL_CANNOT_BUY)
			return
		}
		endTime, ok := m.tlgService.GetGiftEndTime(data.Id)
		if !ok {
			logrus.Debug("购买时限礼包，无法购买")
			result.Add(promotion.ERR_BUY_TIME_LIMIT_GIFT_FAIL_CANNOT_BUY)
			return
		}
		giftIndex := hero.Promotion().GetTimeLimitGiftIndex(data.Id, endTime)
		if giftIndex >= u64.FromInt(len(data.Gifts)) {
			logrus.Debug("购买时限礼包，礼包已购买")
			result.Add(promotion.ERR_BUY_TIME_LIMIT_GIFT_FAIL_BOUGHT)
			return
		}
		gift := data.Gifts[giftIndex]
		if gift.BuyLimit > 0 && hero.Promotion().GetTimeLimitGiftBuyTimes(gift.Id) >= gift.BuyLimit {
			logrus.Debug("购买时限礼包，购买次数超限，无法购买")
			result.Add(promotion.ERR_BUY_TIME_LIMIT_GIFT_FAIL_CANNOT_BUY)
			return
		}
		if !hero.HasEnoughYuanbao(gift.YuanbaoPrice) {
			logrus.Debug("购买时限礼包，元宝不足")
			result.Add(promotion.ERR_BUY_TIME_LIMIT_GIFT_FAIL_NOT_ENOUGH_YUANBAO)
			return
		}
		// 扣元宝
		hctx := heromodule.NewContext(m.dep, operate_type.PromotionBuyTimeLimitGift)
		heromodule.ReduceYuanbaoAnyway(hctx, hero, result, gift.YuanbaoPrice)
		// 给奖励
		heromodule.AddDianquan(hctx, hero, result, gift.Dianquan)
		if prize := gift.Prize; prize != nil {
			heromodule.AddPrize(hctx, hero, result, prize, m.timeService.CurrentTime())
		}
		if guildEventPrize := m.dep.Datas().GetGuildEventPrizeData(gift.GuildEventPrizeId); guildEventPrize != nil {
			m.guildModule.HandleGiveGuildEventPrize(hero, result, hero.GuildId(), []*guild_data.GuildEventPrizeData{guildEventPrize}, 0)
		}

		giftIndex++
		hero.Promotion().SetTimeLimitGiftIndex(data.Id, giftIndex, endTime)
		hero.Promotion().IncTimeLimitGiftBuyTimes(gift.Id)

		result.Add(promotion.NewS2cBuyTimeLimitGiftMsg(u64.Int32(data.Id), u64.Int32(giftIndex)))
		result.Ok()
	})
}

//gogen:iface
func (m *PromotionModule) ProcessBuyEventLimitGift(proto *promotion.C2SBuyEventLimitGiftProto, hc iface.HeroController) {
	data := m.datas.GetEventLimitGiftData(u64.FromInt32(proto.GetId()))
	if data == nil {
		logrus.Debug("购买事件礼包，无效的id")
		hc.Send(promotion.ERR_BUY_EVENT_LIMIT_GIFT_FAIL_INVALID_ID)
		return
	}
	ctime := m.timeService.CurrentTime()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		gift := hero.Promotion().GetEventLimitGiftWithCheck(data.Id, ctime)
		if gift == nil {
			logrus.Debug("购买事件礼包，无法购买")
			result.Add(promotion.ERR_BUY_EVENT_LIMIT_GIFT_FAIL_CANNOT_BUY)
			return
		}
		var guanfuLevel uint64
		if guanFu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU); guanFu != nil {
			guanfuLevel = guanFu.Level
		}
		if !data.IsCanBuy(hero.Level(), guanfuLevel) {
			logrus.Debug("购买事件礼包，无法购买")
			result.Add(promotion.ERR_BUY_EVENT_LIMIT_GIFT_FAIL_CANNOT_BUY)
			return
		}
		if !hero.HasEnoughYuanbao(data.YuanbaoPrice) {
			logrus.Debug("购买时限礼包，元宝不足")
			result.Add(promotion.ERR_BUY_EVENT_LIMIT_GIFT_FAIL_NOT_ENOUGH_YUANBAO)
			return
		}
		// 扣元宝
		hctx := heromodule.NewContext(m.dep, operate_type.PromotionBuyTimeLimitGift)
		heromodule.ReduceYuanbaoAnyway(hctx, hero, result, data.YuanbaoPrice)
		// 给奖励
		heromodule.AddDianquan(hctx, hero, result, data.Dianquan)
		heromodule.AddPrize(hctx, hero, result, gift.Data().Prize, ctime)
		if guildEventPrize := m.dep.Datas().GetGuildEventPrizeData(gift.Data().GuildEventPrizeId); guildEventPrize != nil {
			m.guildModule.HandleGiveGuildEventPrize(hero, result, hero.GuildId(), []*guild_data.GuildEventPrizeData{guildEventPrize}, 0)
		}

		hero.Promotion().RemoveEventLimitGift(data.Id)
		hero.Promotion().IncEventLimitGiftBuyTimes(data.Id)

		result.Add(promotion.NewS2cBuyEventLimitGiftMsg(proto.Id))
		result.Ok()
	})
}
