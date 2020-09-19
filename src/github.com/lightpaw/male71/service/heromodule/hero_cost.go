package heromodule

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/pb/depot"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/config/strategydata"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gamelogs"
)

func reduceSingleResource(heroResource *entity.ResourceStorage, result herolock.LockResult, resType shared_proto.ResType, toReduce uint64) bool {
	if toReduce <= 0 {
		return true
	}

	newRes := heroResource.ReduceRes(resType, toReduce)
	result.Add(domestic.NewS2cResourceUpdateSingleMsg(int32(resType), u64.Int32(newRes), heroResource.IsSafeResource()))

	result.Changed()
	return true
}

func ReduceResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gold, food, wood, stone uint64) {
	if gold|food|wood|stone <= 0 {
		return
	}

	// 先扣掉不安全的资源
	if unsafe := hero.GetUnsafeResource(); unsafe != nil {

		toReduceGold := u64.Min(unsafe.Gold(), gold)
		toReduceFood := u64.Min(unsafe.Food(), food)
		toReduceWood := u64.Min(unsafe.Wood(), wood)
		toReduceStone := u64.Min(unsafe.Stone(), stone)

		newGold, newFood, newWood, newStone := unsafe.ReduceResource(toReduceGold, toReduceFood, toReduceWood, toReduceStone)
		result.AddFunc(func() pbutil.Buffer {
			return domestic.NewS2cResourceUpdateMsg(
				u64.Int32(newGold),
				u64.Int32(newFood),
				u64.Int32(newWood),
				u64.Int32(newStone),
				false,
			)
		})

		gold -= toReduceGold
		food -= toReduceFood
		wood -= toReduceWood
		stone -= toReduceStone
	}

	if gold|food|wood|stone <= 0 {
		return
	}

	// 不够再扣安全资源
	newGold, newFood, newWood, newStone := hero.GetSafeResource().ReduceResource(gold, food, wood, stone)
	result.AddFunc(func() pbutil.Buffer {
		return domestic.NewS2cResourceUpdateMsg(
			u64.Int32(newGold),
			u64.Int32(newFood),
			u64.Int32(newWood),
			u64.Int32(newStone),
			true,
		)
	})

	result.Changed()

	tlogResource(hctx, hero, result, gold, food, wood, stone, false)
}

func ReduceUnsafeResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gold, food, wood, stone uint64) (changed bool, msg func() pbutil.Buffer) {
	if hero.GetUnsafeResource() == nil {
		return
	}

	changed, msg = ReduceStorageResource(hero.GetUnsafeResource(), gold, food, wood, stone)

	tlogResource(hctx, hero, result, gold, food, wood, stone, false)
	return
}

func ReduceSafeResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gold, food, wood, stone uint64) (changed bool, msg func() pbutil.Buffer) {
	changed, msg = ReduceStorageResource(hero.GetSafeResource(), gold, food, wood, stone)

	tlogResource(hctx, hero, result, gold, food, wood, stone, false)

	return
}

func tlogResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gold, food, wood, stone uint64, isAdd bool) {
	tlogSingleResource(hctx, hero, result, gold, shared_proto.ResType_GOLD, isAdd)
	tlogSingleResource(hctx, hero, result, stone, shared_proto.ResType_STONE, isAdd)
	tlogSingleResource(hctx, hero, result, wood, shared_proto.ResType_WOOD, isAdd)
	tlogSingleResource(hctx, hero, result, food, shared_proto.ResType_FOOD, isAdd)
}

func tlogSingleResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64, resType shared_proto.ResType, isAdd bool) {
	if amount <= 0 {
		return
	}

	var oper uint64
	if isAdd {
		oper = operate_type.Add
	} else {
		oper = operate_type.Reduce
	}

	if amount > 0 {
		resAmount := hero.GetAllRes(resType)
		hctx.Tlog().TlogMoneyFlow(hero, uint64(ChangeResTypeToAmountType(resType)), resAmount, amount, hctx.operateType.Id(), oper)
	}
}

func ReduceStorageResource(heroResource *entity.ResourceStorage, gold, food, wood, stone uint64) (changed bool, msg func() pbutil.Buffer) {
	if gold+food+wood+stone <= 0 {
		return false, nil
	}

	newGold, newFood, newWood, newStone := heroResource.ReduceResource(gold, food, wood, stone)
	isSafeResource := heroResource.IsSafeResource()

	return true, func() pbutil.Buffer {
		return domestic.NewS2cResourceUpdateMsg(
			u64.Int32(newGold),
			u64.Int32(newFood),
			u64.Int32(newWood),
			u64.Int32(newStone),
			isSafeResource,
		)
	}
}

func TryReduceSingleResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, resType shared_proto.ResType, toReduce uint64) (success bool) {

	unsafe := hero.GetUnsafeResource()
	res := unsafe.GetRes(resType)
	if res >= toReduce {
		newRes := unsafe.ReduceRes(resType, toReduce)
		result.Add(domestic.NewS2cResourceUpdateSingleMsg(int32(resType), u64.Int32(newRes), false))
		return true
	} else {
		toReduceSafe := toReduce - res

		safe := hero.GetSafeResource()
		if !safe.HasEnoughRes(resType, toReduceSafe) {
			logrus.Debug("heromodule.TryReduceSingleResource not enough resource")
			return false
		}

		reduceSingleResource(unsafe, result, resType, res)
		reduceSingleResource(safe, result, resType, toReduceSafe)

		tlogSingleResource(hctx, hero, result, toReduce, resType, false)

		return true
	}
}

func TryReduceResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gold, food, wood, stone uint64) (success bool) {

	if !hero.HasEnoughResource(gold, food, wood, stone) {
		logrus.Debugf("heromodule.TryReduceResource not enough resource")
		return false
	}

	ReduceResource(hctx, hero, result, gold, food, wood, stone)
	return true
}

func HasEnoughCombineCost(hero *entity.Hero, cost *domestic_data.CombineCost, ctime time.Time) bool {
	if !HasEnoughCost(hero, cost.Cost) {
		return false
	}

	if cost.BuildingWorkerTime > 0 && hero.Domestic().GetFreeWorker(ctime) < 0 {
		return false
	}

	if cost.TechWorkerTime > 0 && hero.Domestic().GetFreeTechWorker(ctime) < 0 {
		return false
	}

	if hero.Military().FreeSoldier(ctime) < cost.Soldier {
		return false
	}

	if !hero.GetInvaseHeroTimes().HasEnoughTimes(cost.InvadeTimes, ctime, 0) {
		return false
	}

	return true
}

func TryReduceCombineCost(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, cost *domestic_data.CombineCost, ctime time.Time) bool {
	if !HasEnoughCombineCost(hero, cost, ctime) {
		logrus.Debugf("hc.TryReduceCombineCost cost not enough")
		return false
	}

	ReduceCombineCostAnyway(hctx, hero, result, cost, ctime)
	return true
}

func ReduceCombineCostAnyway(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, cost *domestic_data.CombineCost, ctime time.Time) {
	ReduceCostAnyway(hctx, hero, result, cost.Cost)
	ReduceFreeSoldier(hero, result, cost.Soldier, ctime)
	UseBuildingWorkerTime(hero, result, cost.BuildingWorkerTime, ctime)
	UseTechWorkerTime(hero, result, cost.TechWorkerTime, ctime)

	oldTime := hero.GetInvaseHeroTimes().StartTime()
	hero.GetInvaseHeroTimes().ReduceTimes(cost.InvadeTimes, ctime, 0)
	newTime := hero.GetInvaseHeroTimes().StartTime()
	if !newTime.Equal(oldTime) {
		result.Add(region.NewS2cUpdateInvaseHeroTimesMsg(hero.GetInvaseHeroTimes().StartTimeUnix32(), nil))
	}
}

func UseTechWorkerTime(hero *entity.Hero, result herolock.LockResult, duration time.Duration, ctime time.Time) {
	if duration <= 0 {
		return
	}

	heroDomestic := hero.Domestic()
	workerPos := heroDomestic.GetFreeTechWorker(ctime)
	if workerPos < 0 {
		logrus.Errorln("存在玩家增加科研队时间，加不上去，问题是此前检查了")
	} else {
		workerRestEndTime, seekHelp := heroDomestic.AddTechWorkerRestEndTime(workerPos, ctime, duration)
		result.Add(domestic.NewS2cTechWorkerTimeChangedMsg(int32(workerPos), timeutil.Marshal32(workerRestEndTime)))

		if seekHelp {
			result.Add(guild.NewS2cUpdateSeekHelpMsg(constants.SeekTypeTech, int32(workerPos), seekHelp))
		}
	}
}

func UseBuildingWorkerTime(hero *entity.Hero, result herolock.LockResult, duration time.Duration, ctime time.Time) {
	if duration <= 0 {
		return
	}

	heroDomestic := hero.Domestic()
	workerPos := heroDomestic.GetFreeWorker(ctime)
	if workerPos < 0 {
		logrus.Errorln("存在玩家增加建筑队时间，加不上去，问题是此前检查了")
	} else {
		workerRestEndTime, seekHelp := heroDomestic.AddWorkerRestEndTime(workerPos, ctime, duration)
		result.Add(domestic.NewS2cBuildingWorkerTimeChangedMsg(int32(workerPos), timeutil.Marshal32(workerRestEndTime)))

		if seekHelp {
			// 发消息通知，可以求助
			result.Add(guild.NewS2cUpdateSeekHelpMsg(constants.SeekTypeWorker, imath.Int32(workerPos), seekHelp))
		}
	}
}

func ReduceFreeSoldier(hero *entity.Hero, result herolock.LockResult, toReduce uint64, ctime time.Time) {
	if toReduce <= 0 {
		return
	}

	heroMilitary := hero.Military()
	heroMilitary.ReduceFreeSoldier(toReduce, ctime)

	result.Add(hero.Military().NewUpdateFreeSoldierMsg())
	result.Changed()
}

func HasEnoughGoodsOrBuy(hero *entity.Hero, goodsData *goods.GoodsData, goodsCount uint64, autoBuy bool) bool {
	ok, toReduceGoodsCount, _, _, _ := planReduceOrBuyGoods(hero, goodsData, goodsCount)
	if autoBuy {
		return ok
	} else {
		return ok && toReduceGoodsCount >= goodsCount
	}
}

func planReduceOrBuyGoods(hero *entity.Hero, goodsData *goods.GoodsData, goodsCount uint64) (ok bool, toReduceGoodsCount, toReduceYinliang, toReduceDianquan, toReduceYuanbao uint64) {

	totalCount := hero.Depot().GetGoodsCount(goodsData.Id)
	if goodsCount <= totalCount {
		return true, goodsCount, 0, 0, 0
	}

	buyCount := u64.Sub(goodsCount, totalCount)
	switch {
	case goodsData.YinliangPrice > 0:
		toReduceYinliang = goodsData.YinliangPrice * buyCount
		if !hero.HasEnoughYinliang(toReduceYinliang) {
			return false, 0, 0, 0, 0
		}
		return true, totalCount, toReduceYinliang, 0, 0
	case goodsData.DianquanPrice > 0:
		toReduceDianquan = goodsData.DianquanPrice * buyCount
		if !hero.HasEnoughDianquan(toReduceDianquan) {
			return false, 0, 0, 0, 0
		}
		return true, totalCount, 0, toReduceDianquan, 0
	case goodsData.YuanbaoPrice > 0:
		toReduceYuanbao = goodsData.YuanbaoPrice * buyCount
		if !hero.HasEnoughYuanbao(toReduceYuanbao) {
			return false, 0, 0, 0, 0
		}
		return true, totalCount, 0, 0, toReduceYuanbao
	}

	return false, 0, 0, 0, 0
}

func TryReduceOrBuyGoods(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *goods.GoodsData, goodsCount uint64, autoBuy bool) bool {
	ok, _, _, _, _ := TryReduceOrBuyGoodsReturnPlan(hctx, hero, result, goodsData, goodsCount, autoBuy)
	return ok
}

func TryReduceOrBuyGoodsReturnPlan(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *goods.GoodsData, goodsCount uint64, autoBuy bool) (ok bool, toReduceGoodsCount, toReduceYinliang, toReduceDianquan, toReduceYuanbao uint64) {

	ok, toReduceGoodsCount, toReduceYinliang, toReduceDianquan, toReduceYuanbao = planReduceOrBuyGoods(hero, goodsData, goodsCount)
	if !ok {
		return false, 0, 0, 0, 0
	}

	if !autoBuy && toReduceGoodsCount < goodsCount {
		return false, 0, 0, 0, 0
	}

	ReduceGoodsAnyway(hctx, hero, result, goodsData, toReduceGoodsCount)
	ReduceYinliangAnyway(hctx, hero, result, toReduceYinliang)
	ReduceDianquanAnyway(hctx, hero, result, toReduceDianquan)
	ReduceYuanbaoAnyway(hctx, hero, result, toReduceYuanbao)

	return
}

func calHasEnoughCostWithDianquanConvert(hero *entity.Hero, cost *resdata.Cost, datas iface.ConfigDatas) (dianquanReduce, goldAdd, stoneAdd uint64, ok bool) {
	if !cost.IsOnlyResource {
		// 配置加载判定过，所以获取到的calCost不可能为Empty
		calCost := resdata.NewCostBuilder().Add(cost).ReduceResource(cost.Gold, cost.Food, cost.Wood, cost.Stone).Build()
		if !HasEnoughCost(hero, calCost) {
			return
		}
	}
	// 走到这里必然是只有Resource
	if gold := u64.Sub(cost.Gold, hero.GetGold()); gold > 0 {
		// 资源1兑换资源2通过兑换比例得出资源1需要的数量
		dianquanReduce = u64.DivideTimes(gold, datas.MiscGenConfig().DianquanToGold)
		goldAdd = datas.MiscGenConfig().DianquanToGold * dianquanReduce
	}
	if stone := u64.Sub(cost.Stone, hero.GetStone()); stone > 0 {
		dq := u64.DivideTimes(stone, datas.MiscGenConfig().DianquanToStone)
		stoneAdd = datas.MiscGenConfig().DianquanToStone * dq
		dianquanReduce += dq
	}
	ok = hero.HasEnoughDianquan(dianquanReduce + cost.Dianquan)
	return
}

func HasEnoughCostWithDianquanConvert(hero *entity.Hero, cost *resdata.Cost, datas iface.ConfigDatas) bool {
	_, _, _, ok := calHasEnoughCostWithDianquanConvert(hero, cost, datas)
	return ok
}

func HasEnoughCost(hero *entity.Hero, cost *resdata.Cost) bool {
	return GetCostCount(hero, cost) > 0
}

func GetCostCount(hero *entity.Hero, cost *resdata.Cost) uint64 {

	n := hero.GetResourceMultiple(cost.GetResource())
	if n <= 0 {
		return n
	}

	if cost.Jade > 0 {
		n = u64.Min(n, hero.Jade()/cost.Jade)
		if n <= 0 {
			return n
		}
	}

	if cost.JadeOre > 0 {
		n = u64.Min(n, hero.JadeOre()/cost.JadeOre)
		if n <= 0 {
			return n
		}
	}

	if cost.GuildContributionCoin > 0 {
		n = u64.Min(n, hero.GetGuildContributionCoin()/cost.GuildContributionCoin)
		if n <= 0 {
			return n
		}
	}

	if cost.Yuanbao > 0 {
		n = u64.Min(n, hero.GetYuanbao()/cost.Yuanbao)
		if n <= 0 {
			return n
		}
	}

	if cost.Dianquan > 0 {
		n = u64.Min(n, hero.GetDianquan()/cost.Dianquan)
		if n <= 0 {
			return n
		}
	}

	if cost.Yinliang > 0 {
		n = u64.Min(n, hero.GetYinliang()/cost.Yinliang)
		if n <= 0 {
			return n
		}
	}

	heroDepot := hero.Depot()
	for i, g := range cost.Goods {
		if cost.GoodsCount[i] > 0 {
			count := heroDepot.GetGoodsCount(g.Id) / cost.GoodsCount[i]
			if count <= 0 {
				return 0
			}

			n = u64.Min(n, count)
		}
	}

	for i, g := range cost.Gem {
		if cost.GemCount[i] > 0 {
			count := heroDepot.GetGoodsCount(g.Id) / cost.GemCount[i]
			if count <= 0 {
				return 0
			}

			n = u64.Min(n, count)
		}
	}

	return n
}

func TryReduceCostWithDianquanConvert(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, cost *resdata.Cost) bool {
	dianquan, gold, stone, ok := calHasEnoughCostWithDianquanConvert(hero, cost, hctx.datas)
	if !ok {
		logrus.Debugf("hc.TryReduceCostWithDianquanConvertUnderLock not enough")
		return false
	}
	if dianquan > 0 {
		ReduceDianquanAnyway(hctx, hero, result, dianquan)
		AddSafeResource(hctx, hero, result, gold, 0, 0, stone)
	}
	ReduceCostAnyway(hctx, hero, result, cost)
	return true
}

func TryReduceCost(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, cost *resdata.Cost) bool {

	if !HasEnoughCost(hero, cost) {
		logrus.Debugf("hc.TryReduceCostUnderLock not enough")
		return false
	}

	ReduceCostAnyway(hctx, hero, result, cost)
	return true
}

func ReduceCostAnyway(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, cost *resdata.Cost) {

	// 元宝
	if cost.Yuanbao > 0 {
		ReduceYuanbaoAnyway(hctx, hero, result, cost.Yuanbao)
	}

	// 点券
	if cost.Dianquan > 0 {
		ReduceDianquanAnyway(hctx, hero, result, cost.Dianquan)
	}

	// 银两
	if cost.Yinliang > 0 {
		hero.ReduceYinliang(cost.Yinliang)
		result.Add(domestic.NewS2cUpdateYinliangMsg(u64.Int32(hero.GetYinliang())))

		result.Changed()

		yinliang := hero.GetYinliang()
		hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_YINLIANG), yinliang, cost.Yinliang, hctx.operateType.Id(), operate_type.Reduce)
	}

	// 帮派贡献币
	if cost.GuildContributionCoin > 0 {
		hero.ReduceGuildContributionCoin(cost.GuildContributionCoin)
		result.Add(guild.NewS2cUpdateContributionCoinMsg(u64.Int32(hero.GetGuildContributionCoin())))

		result.Changed()

		gcCoin := hero.GetGuildContributionCoin()
		hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_GUILD_CONTRIBUTION_COIN), gcCoin, cost.GuildContributionCoin, hctx.operateType.Id(), operate_type.Reduce)
	}

	if cost.HasResource() {
		gold, food, wood, stone := cost.GetResource()
		ReduceResource(hctx, hero, result, gold, food, wood, stone)
	}

	// 玉璧
	ReduceJadeAnyway(hctx, hero, result, cost.Jade)
	// 玉石矿
	ReduceJadeOreAnyway(hctx, hero, result, cost.JadeOre)

	ReduceGoodsArray(hctx, hero, result, cost.Goods, cost.GoodsCount)
	ReduceGemArray(hctx, hero, result, cost.Gem, cost.GemCount)

	result.Changed()
}

func ReduceGoodsArray(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsDatas []*goods.GoodsData, goodsCount []uint64) {
	n := imath.Min(len(goodsDatas), len(goodsCount))
	if n <= 0 {
		return
	}

	heroDepot := hero.Depot()

	newCount := make([]uint64, n)
	oldCount := make([]uint64, n)
	for i := 0; i < n; i++ {
		g, count := goodsDatas[i], goodsCount[i]
		oldCount[i] = heroDepot.GetGoodsCount(g.Id)
		newCount[i] = heroDepot.RemoveGoods(g.Id, count)
	}

	result.Add(depot.NewS2cUpdateMultiGoodsMsg(u64.Int32Array(goods.GetGoodsDataKeyArray(goodsDatas)), u64.Int32Array(newCount)))
	result.Changed()

	// tlog
	for i := 0; i < n; i++ {
		g, count, oldCnt, newCnt := goodsDatas[i], goodsCount[i], oldCount[i], newCount[i]
		if count <= 0 {
			continue
		}
		hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeConsumable, uint64(g.EffectType)), g.Id, oldCnt, newCnt, hctx.operateType.Id(), operate_type.Reduce, g.GoodsQuality.Level, g.Id, -int64(count))
	}
}

func TryReduceGoods(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *goods.GoodsData, goodsCount uint64) bool {
	if hero.Depot().HasEnoughGoods(goodsData.Id, goodsCount) {
		ReduceGoodsAnyway(hctx, hero, result, goodsData, goodsCount)
		return true
	}
	return false
}

func ReduceGoodsAnyway(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, g *goods.GoodsData, count uint64) {
	if count <= 0 {
		return
	}

	heroDepot := hero.Depot()

	oldCount := heroDepot.GetGoodsCount(g.Id)
	newCount := heroDepot.RemoveGoods(g.Id, count)

	result.Add(depot.NewS2cUpdateGoodsMsg(u64.Int32(g.Id), u64.Int32(newCount)))

	// tlog
	hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeConsumable, uint64(g.EffectType)), g.Id, oldCount, newCount, hctx.operateType.Id(), operate_type.Reduce, g.GoodsQuality.Level, g.Id, -int64(count))
}

func ReduceGemArray(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gemDatas []*goods.GemData, gemCount []uint64) {
	n := imath.Min(len(gemDatas), len(gemCount))
	if n <= 0 {
		return
	}

	heroDepot := hero.Depot()

	oldCount := make([]uint64, n)
	newCount := make([]uint64, n)

	for i := 0; i < n; i++ {
		gemData, count := gemDatas[i], gemCount[i]
		oldCount[i] = heroDepot.GetGoodsCount(gemData.Id)
		newCount[i] = heroDepot.RemoveGoods(gemData.Id, count)
	}

	result.Add(depot.NewS2cUpdateMultiGoodsMsg(u64.Int32Array(goods.GetGemDataKeyArray(gemDatas)), u64.Int32Array(newCount)))
	result.Changed()

	// tlog
	for i := 0; i < n; i++ {
		gemData, count, oldCnt, newCnt := gemDatas[i], gemCount[i], oldCount[i], newCount[i]
		if count <= 0 {
			continue
		}
		hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeGem, gemData.GemType), gemData.Id, oldCnt, newCnt, hctx.operateType.Id(), operate_type.Reduce, uint64(gemData.Quality), gemData.Id, -int64(count))
	}
}

func ReduceGemAnyway(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gemData *goods.GemData, count uint64) {
	if count <= 0 {
		return
	}

	heroDepot := hero.Depot()

	oldCount := heroDepot.GetGoodsCount(gemData.Id)
	newCount := heroDepot.RemoveGoods(gemData.Id, count)

	result.Add(depot.NewS2cUpdateGoodsMsg(u64.Int32(gemData.Id), u64.Int32(newCount)))

	// tlog
	hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeGem, gemData.GemType), gemData.Id, oldCount, newCount, hctx.operateType.Id(), operate_type.Reduce, uint64(gemData.Quality), gemData.Id, -int64(count))
}

func ReduceYuanbao(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) bool {
	if hero.GetYuanbao() < amount {
		return false
	}

	ReduceYuanbaoAnyway(hctx, hero, result, amount)
	return true
}

func ReduceYuanbaoAnyway(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount == 0 {
		return
	}

	hero.ReduceYuanbao(amount)
	result.Add(domestic.NewS2cUpdateYuanbaoMsg(u64.Int32(hero.GetYuanbao())))

	result.Changed()

	yuanbao := hero.GetYuanbao()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_YUANBAO), yuanbao, amount, hctx.operateType.Id(), operate_type.Reduce)

	gamelogs.YuanbaoReduceLog(constants.PID, hero.Sid(), hero.Id(), hctx.OperType().Id(), amount) // TODO
}

func ReduceDianquan(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) bool {
	if hero.GetDianquan() < amount {
		return false
	}

	ReduceDianquanAnyway(hctx, hero, result, amount)
	return true
}

func ReduceDianquanAnyway(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount == 0 {
		return
	}

	hero.ReduceDianquan(amount)
	result.Add(domestic.NewS2cUpdateDianquanMsg(u64.Int32(hero.GetDianquan())))

	result.Changed()

	dianquan := hero.GetDianquan()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_DIANQUAN), dianquan, amount, hctx.operateType.Id(), operate_type.Reduce)

	gamelogs.DianquanReduceLog(constants.PID, hero.Sid(), hero.Id(), hctx.OperType().Id(), amount) // TODO
}

func ReduceYinliang(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) bool {
	if hero.GetYinliang() < amount {
		return false
	}

	ReduceYinliangAnyway(hctx, hero, result, amount)
	return true
}

func ReduceYinliangAnyway(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount == 0 {
		return
	}

	hero.ReduceYinliang(amount)
	result.Add(domestic.NewS2cUpdateYinliangMsg(u64.Int32(hero.GetYinliang())))

	result.Changed()

	yinliang := hero.GetYinliang()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_YINLIANG), yinliang, amount, hctx.operateType.Id(), operate_type.Reduce)
}

func ReduceJade(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) bool {
	if amount < 0 {
		return false
	}

	if hero.Jade() < amount {
		return false
	}

	ReduceJadeAnyway(hctx, hero, result, amount)
	return true
}

func ReduceJadeAnyway(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount < 0 {
		return
	}

	hero.ReduceJade(amount)
	result.Add(domestic.NewS2cUpdateJadeMsg(u64.Int32(hero.Jade()), u64.Int32(hero.HistoryJade()), u64.Int32(hero.TodayObtainJade())))

	result.Changed()

	jade := hero.Jade()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_JADE), jade, amount, hctx.operateType.Id(), operate_type.Reduce)
}

func HasEnoughJadeOre(hero *entity.Hero, amount uint64) bool {
	if amount < 0 {
		return false
	}

	return hero.JadeOre() >= amount
}

func ReduceJadeOre(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) bool {
	if amount < 0 {
		return false
	}

	if hero.JadeOre() < amount {
		return false
	}

	ReduceJadeOreAnyway(hctx, hero, result, amount)

	return true
}

func ReduceJadeOreAnyway(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount < 0 {
		return
	}

	hero.ReduceJadeOre(amount)
	result.Add(domestic.NewS2cUpdateJadeOreMsg(u64.Int32(hero.JadeOre())))

	result.Changed()

	jadeOre := hero.JadeOre()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_JADE_ORE), jadeOre, amount, hctx.operateType.Id(), operate_type.Reduce)

}

func TryReduceBaowu(hero *entity.Hero, result herolock.LockResult, goodsData *resdata.BaowuData, goodsCount uint64) bool {
	if hero.Depot().HasEnoughBaowu(goodsData.Id, goodsCount) {
		ReduceBaowuAnyway(hero, result, goodsData, goodsCount)
		return true
	}
	return false
}

func ReduceBaowuAnyway(hero *entity.Hero, result herolock.LockResult, goodsData *resdata.BaowuData, goodsCount uint64) {
	if goodsCount <= 0 {
		return
	}

	heroDepot := hero.Depot()

	newCount := heroDepot.RemoveBaowu(goodsData.Id, goodsCount)

	result.Add(depot.NewS2cUpdateBaowuMsg(u64.Int32(goodsData.Id), u64.Int32(newCount)))
}

// 计算使用君主策略的真实消耗
func StrategyUsedRealCost(hero *entity.Hero, effect *strategydata.StrategyEffectData) (result *domestic_data.CombineCost) {
	result = effect.Cost

	if effect.Prize.Prosperity <= 0 {
		return
	}

	toAdd := u64.Sub(hero.ProsperityCapcity(), hero.Prosperity())
	if effect.Prize.Prosperity <= toAdd {
		return
	}

	per := u64.Division2Float64(toAdd, effect.Prize.Prosperity)

	var realSolider uint64
	if effect.Cost.Soldier > 0 {
		realSolider = u64.Max(0, u64.MultiF64(effect.Cost.Soldier, per))
	}
	realCost := effect.Cost.Cost.MultipleF64(per)
	result = &domestic_data.CombineCost{Cost: realCost, Soldier: realSolider}
	return
}
