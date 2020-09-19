package shop

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/resdata"
	shop2 "github.com/lightpaw/male7/config/shop"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/shop"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
)

func NewShopModule(dep iface.ServiceDep, guildModule iface.GuildModule) *ShopModule {
	m := &ShopModule{}
	m.dep = dep
	m.guildModule = guildModule
	m.datas = dep.Datas()
	m.time = dep.Time()
	m.worldService = dep.World()
	m.sharedGuildService = dep.Guild()
	m.broadcast = dep.Broadcast()

	return m
}

//gogen:iface
type ShopModule struct {
	dep         iface.ServiceDep
	datas       iface.ConfigDatas
	time        iface.TimeService
	guildModule iface.GuildModule

	worldService       iface.WorldService
	sharedGuildService iface.GuildService
	broadcast          iface.BroadcastService

	effectResourceGoodsIdListMsg pbutil.Buffer
}

//gogen:iface
func (m *ShopModule) ProcessBuyGoods(proto *shop.C2SBuyGoodsProto, hc iface.HeroController) {

	count := u64.FromInt32(proto.Count)
	if count <= 0 {
		logrus.Debugf("购买物品，购买个数无效, %v", proto.Count)
		hc.Send(shop.ERR_BUY_GOODS_FAIL_INVALID_ID)
		return
	}

	if normalGoods := m.datas.GetNormalShopGoods(u64.FromInt32(proto.Id)); normalGoods != nil {
		m.buyNormalGoods(hc, normalGoods, count)
	} else if zhenBaoGeGoods := m.datas.GetZhenBaoGeShopGoods(u64.FromInt32(proto.Id)); zhenBaoGeGoods != nil {
		m.buyZhenBaoGeGoods(hc, zhenBaoGeGoods, count)
	} else {
		logrus.Debugf("购买物品，物品id不存在, %v", proto.Id)
		hc.Send(shop.ERR_BUY_GOODS_FAIL_INVALID_ID)
	}
}

func (m *ShopModule) buyZhenBaoGeGoods(hc iface.HeroController, shopGoods *shop2.ZhenBaoGeShopGoods, count uint64) {
	if heromodule.IsLocked(hc.Func, m.sharedGuildService.GetSnapshot, shopGoods.UnlockCondition) {
		logrus.Debugf("购买物品，商品未解锁")
		hc.Send(shop.ERR_BUY_GOODS_FAIL_LOCKED)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.ShopBuyZhenBaoGeGoods)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hasBuyCount := hero.GetShopGoodsCount(shopGoods.Id)

		if count != 1 {
			if len(shopGoods.BuyCosts) != 1 {
				logrus.Debugf("购买物品，一次只能够买一个")
				hc.Send(shop.ERR_BUY_GOODS_FAIL_INVALID_COUNT)
				return
			}

			if hasBuyCount < shopGoods.FreeTimes {
				logrus.Debugln("有免费次数的时候只能够一次买一个")
				hc.Send(shop.ERR_BUY_GOODS_FAIL_INVALID_COUNT)
				return
			}
		}

		if !shopGoods.CanBuy(hasBuyCount, count) {
			logrus.Debugf("购买物品，无法购买")
			hc.Send(shop.ERR_BUY_GOODS_FAIL_INVALID_COUNT)
			return
		}

		guanFuLevel := uint64(1)
		if guanFu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU); guanFu != nil {
			guanFuLevel = guanFu.Level
		}

		prize := shopGoods.GetPrize(guanFuLevel, count)
		if prize == nil {
			logrus.Errorf("购买物品，没找到奖励")
			hc.Send(shop.ERR_BUY_GOODS_FAIL_SERVER_ERROR)
			return
		}

		ctime := m.time.CurrentTime()

		var useGoods []*goods.GoodsData
		var useFuncs []heromodule.UseGoodsFunc
		if len(shopGoods.UseImmediatelyGoods) > 0 {

			for _, g := range shopGoods.UseImmediatelyGoods {
				maxCanUseCount, useFunc := heromodule.GetUseEffectGoodsFunc(hctx, hero, g, count, ctime)
				if useFunc == nil {
					logrus.Debugf("购买物品并使用，发现物品并不能直接使用")
					hc.Send(shop.ERR_BUY_GOODS_FAIL_CANT_USE_GOODS)
					return
				}

				if maxCanUseCount < count {
					logrus.Debugf("购买物品并使用，购买数量超出最大可使用个数")
					hc.Send(shop.ERR_BUY_GOODS_FAIL_INVALID_COUNT)
					return
				}

				useGoods = append(useGoods, g)
				useFuncs = append(useFuncs, useFunc)
			}
		}

		cost := shopGoods.GetBuyCost(hasBuyCount)
		if cost != nil {
			cost = cost.Multiple(count)
			if !heromodule.TryReduceCost(hctx, hero, result, cost) {
				logrus.Debugf("购买物品，消耗不够")
				hc.Send(shop.ERR_BUY_GOODS_FAIL_COST_NOT_ENOUGH)
				return
			}
		}

		multi, needBroadcast := shopGoods.RandomCrit()
		if multi != 1 {
			prize = prize.Multiple(multi)
		}

		prizeBytes := m.doAfterBuyGoods(hctx, hero, result, prize, multi, shopGoods.Id, count, shopGoods.ShopType, nil, useGoods, useFuncs, ctime)

		if needBroadcast {
			m.worldService.Broadcast(shop.NewS2cMultiCritBroadcastMsg(u64.Int32(shopGoods.ShopType), u64.Int32(multi), hero.Name(), prizeBytes))
		}

		if d := hctx.BroadcastHelp().ZhenBaoGeBaoJi; d != nil {
			hctx.AddBroadcast(d, hero, result, 0, multi, func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, multi)
				return text
			})
		}

	})
}

func (m *ShopModule) buyNormalGoods(hc iface.HeroController, shopGoods *shop2.NormalShopGoods, count uint64) {
	if heromodule.IsLocked(hc.Func, m.sharedGuildService.GetSnapshot, shopGoods.UnlockCondition) {
		logrus.Debugf("购买物品，商品未解锁")
		hc.Send(shop.ERR_BUY_GOODS_FAIL_LOCKED)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.ShopBuyNormalGoods)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hasBuyCount := hero.GetShopGoodsCount(shopGoods.Id)

		if count != 1 {
			if hasBuyCount < shopGoods.FreeTimes {
				logrus.Debugln("有免费次数的时候只能够一次买一个")
				hc.Send(shop.ERR_BUY_GOODS_FAIL_INVALID_COUNT)
				return
			}
		}

		if !shopGoods.CanBuy(hasBuyCount, count) {
			logrus.Debugf("购买物品，无法购买")
			hc.Send(shop.ERR_BUY_GOODS_FAIL_INVALID_COUNT)
			return
		}

		gid := m.dep.GuildSnapshot().GetGuildLevel(hero.GuildId())
		prize := shopGoods.PlunderPrize.GetGuildPrize(gid)

		if prize == nil {
			logrus.Errorf("购买物品，没找到奖励")
			hc.Send(shop.ERR_BUY_GOODS_FAIL_SERVER_ERROR)
			return
		}

		// vip商店单独处理
		if shopGoods.ShopType == shop2.VIPShopType {
			if hero.Vip().VipShopGoodsBoughtCount(shopGoods.Id)+count > shopGoods.VipDailyMaxCount {
				result.Add(shop.ERR_BUY_GOODS_FAIL_DAILY_BOUGHT_LIMIT)
				return
			}
		}

		ctime := m.time.CurrentTime()

		var useGoods []*goods.GoodsData
		var useFuncs []heromodule.UseGoodsFunc
		if len(shopGoods.UseImmediatelyGoods) > 0 {
			for _, g := range shopGoods.UseImmediatelyGoods {
				maxCanUseCount, useFunc := heromodule.GetUseEffectGoodsFunc(hctx, hero, g, count, ctime)
				if useFunc == nil {
					logrus.Debugf("购买物品并使用，发现物品并不能直接使用")
					hc.Send(shop.ERR_BUY_GOODS_FAIL_CANT_USE_GOODS)
					return
				}

				if maxCanUseCount < count {
					logrus.Debugf("购买物品并使用，购买数量超出最大可使用个数")
					hc.Send(shop.ERR_BUY_GOODS_FAIL_INVALID_COUNT)
					return
				}

				useGoods = append(useGoods, g)
				useFuncs = append(useFuncs, useFunc)
			}
		}

		cost := shopGoods.Cost
		if cost != nil {
			cost = cost.Multiple(count)
			if !heromodule.TryReduceCost(hctx, hero, result, cost) {
				logrus.Debugf("购买物品，消耗不够")
				hc.Send(shop.ERR_BUY_GOODS_FAIL_COST_NOT_ENOUGH)
				return
			}
		}

		prize = prize.Multiple(count)
		m.doAfterBuyGoods(hctx, hero, result, prize, 1, shopGoods.Id, count, shopGoods.ShopType, shopGoods.GuildEventPrize, useGoods, useFuncs, ctime)
	})
}

func (m *ShopModule) doAfterBuyGoods(hctx *heromodule.HeroContext, hero *entity.Hero, result herolock.LockResult,
	prize *resdata.Prize, multi, shopGoodsId, buyCount, shopType uint64, guildEventPrize *guild_data.GuildEventPrizeData,
	useGoods []*goods.GoodsData, useFuncs []heromodule.UseGoodsFunc, ctime time.Time) (prizeBytes []byte) {

	if n := imath.Minx(len(useGoods), len(useFuncs)); n > 0 {
		for i := 0; i < n; i++ {
			g := useGoods[i]
			f := useFuncs[i]

			realUseCount := f(hctx, hero, result, g, buyCount, ctime)
			if realUseCount != buyCount {
				logrus.WithField("id", g.Id).WithField("name", g.Name).WithField("buyCount", buyCount).WithField("realUseCount", realUseCount).Error("商店购买并使用，发现使用的个数跟真正使用的个数不一样")
			}
		}
	}

	heromodule.AddPrize(hctx, hero, result, prize, m.time.CurrentTime())
	if guildEventPrize != nil {
		m.guildModule.HandleGiveGuildEventPrize(hero, result, hero.GuildId(), []*guild_data.GuildEventPrizeData{guildEventPrize}, 0)
	}

	// 更新购买数量
	if shopType == shop2.VIPShopType {
		hero.Vip().AddVipShopGoodsBoughtCount(shopGoodsId, buyCount)
		result.Add(shop.NewS2cUpdateVipShopGoodsMsg(u64.Int32(shopGoodsId), u64.Int32(buyCount)))
	} else {
		newCount := hero.AddShopGoodsCount(shopGoodsId, buyCount)
		result.Add(shop.NewS2cUpdateDailyShopGoodsMsg(u64.Int32(shopGoodsId), u64.Int32(newCount)))
	}

	prizeBytes = must.Marshal(prize.Encode())
	result.Add(shop.NewS2cBuyGoodsMsg(u64.Int32(shopGoodsId), u64.Int32(buyCount), u64.Int32(multi), prizeBytes))

	result.Changed()
	result.Ok()

	heromodule.IncreTaskProgressWithFunc(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUY_GOODS, func(hero *entity.Hero, t *entity.TaskProgress) uint64 {
		if shopType == t.Target().SubType {
			return buyCount
		}
		return 0
	})
	heromodule.IncreTaskProgressWithFunc(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUY_GOODS_COUNT, func(hero *entity.Hero, t *entity.TaskProgress) uint64 {
		if shopType == t.Target().SubType && shopGoodsId == t.Target().SubTypeId {
			return buyCount
		}
		return 0
	})

	hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_ShopBuy, shopType, buyCount)
	hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_ShopGoodsBuy, shop2.ShopTypeGoodsId(shopType, shopGoodsId), buyCount)

	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BUY_GOODS)
	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BUY_GOODS_COUNT)

	return
}

//gogen:iface
func (m *ShopModule) ProcessBuyBlackMarketGoods(proto *shop.C2SBuyBlackMarketGoodsProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.ShopBuyBlackMarketGoods)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		goods := hero.GetBlackMarketGoods(u64.FromInt32(proto.Index))
		if goods == nil {
			logrus.Debug("购买黑市商品，无效的索引")
			result.Add(shop.ERR_BUY_BLACK_MARKET_GOODS_FAIL_INVALID_INDEX)
			return
		}

		if goods.IsBuy() {
			logrus.Debug("购买黑市商品，物品已经购买过了")
			result.Add(shop.ERR_BUY_BLACK_MARKET_GOODS_FAIL_BUYED)
			return
		}

		// 扣钱
		cost := goods.Goods().Cost
		if goods.Discount() > 0 {
			cost = cost.MultipleF64(u64.Division2Float64(goods.Discount(), 1000))
		}
		if !heromodule.TryReduceCost(hctx, hero, result, cost) {
			logrus.Debug("购买黑市商品，消耗不足")
			result.Add(shop.ERR_BUY_BLACK_MARKET_GOODS_FAIL_COST_NOT_ENOUGH)
			return
		}

		goods.SetBuy()

		// 给奖励
		ctime := m.time.CurrentTime()
		heromodule.AddPrize(hctx, hero, result, goods.Goods().Prize, ctime)

		result.Add(shop.NewS2cBuyBlackMarketGoodsMsg(proto.Index))

		result.Ok()
	})

}

//gogen:iface c2s_refresh_black_market_goods
func (m *ShopModule) ProcessRefreshBlackMarketGoods(hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.ShopRefreshBlackMarketGoods)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		bmd := m.datas.GetBlackMarketData(shop2.BLACK_YUNYOU)
		if bmd == nil {
			logrus.Debug("刷新黑市商品，bmd == nil")
			result.Add(shop.ERR_REFRESH_BLACK_MARKET_GOODS_FAIL_TIMES_LIMIT)
			return
		}

		var vipExtraTimes uint64
		if d := m.dep.Datas().GetVipLevelData(hero.Vip().Level()); d != nil {
			vipExtraTimes = d.AddBlackMarketRefreshTimes
		}
		cost := m.datas.ShopMiscData().GetRefreshBlackMarketCost(hero.Shop().GetBlackMarketDailyRefreshTimes(), vipExtraTimes)
		if cost == nil {
			logrus.Debug("刷新黑市商品，刷新次数已达上限")
			result.Add(shop.ERR_REFRESH_BLACK_MARKET_GOODS_FAIL_TIMES_LIMIT)
			return
		}

		if !heromodule.TryReduceCost(hctx, hero, result, cost) {
			logrus.Debug("刷新黑市商品，消耗不足")
			result.Add(shop.ERR_REFRESH_BLACK_MARKET_GOODS_FAIL_COST_NOT_ENOUGH)
			return
		}

		hero.Shop().IncBlackMarketDailyRefreshTimes()

		// 更新商店数据
		goods, discount := bmd.Random(hero.Level())
		hero.UpdateBlackMarketGoods(goods, discount)

		heromodule.SendBlackMarketGoodsMsg(hero, result, true)

		result.Add(shop.REFRESH_BLACK_MARKET_GOODS_S2C)
		result.Ok()

		m.dep.Tlog().TlogRefreshFlow(hero, uint64(shared_proto.BuildingType_SHI_CHANG), operate_type.RefreshMoney, u64.FromInt(cost.Id))
	})

}
