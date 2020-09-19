package depot

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/depot"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/module/rank/ranklist"
)

func NewDepotModule(dep iface.ServiceDep, rankModule iface.RankModule) *DepotModule {
	m := &DepotModule{time: dep.Time(), datas: dep.Datas()}
	m.dep = dep
	m.rankModule = rankModule
	return m
}

//gogen:iface
type DepotModule struct {
	dep        iface.ServiceDep
	time       iface.TimeService
	datas      iface.ConfigDatas
	rankModule iface.RankModule
}

//gogen:iface
func (m *DepotModule) ProcessUseGoods(proto *depot.C2SUseGoodsProto, hc iface.HeroController) {
	id := u64.FromInt32(proto.Id)
	count := u64.FromInt32(proto.Count)
	if count <= 0 {
		logrus.Debugf("使用物品，发送的Count无效")
		hc.Send(depot.ERR_USE_GOODS_FAIL_INVALID_COUNT)
		return
	}

	goodsData := m.datas.GetGoodsData(id)
	if goodsData == nil {
		logrus.Debugf("使用物品，物品不存在")
		hc.Send(depot.ERR_USE_GOODS_FAIL_INVALID_ID)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.DepotUseGoods)
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroDepot := hero.Depot()
		if !heroDepot.HasEnoughGoods(id, count) {
			logrus.Debugf("使用物品，发送的Count不足")
			result.Add(depot.ERR_USE_GOODS_FAIL_COUNT_NOT_ENOUGH)
			return
		}

		if goodsData.GoodsEffect == nil {
			logrus.Debugf("使用物品，这个物品不能直接使用, %d-%s", goodsData.Id, goodsData.Name)
			result.Add(depot.ERR_USE_GOODS_FAIL_CANT_USE)
			return
		}

		ctime := m.time.CurrentTime()
		// 执行物品效果（核心函数），如果成功useCount必定大于0
		maxCanUseCount, useFunc := heromodule.GetUseEffectGoodsFunc(hctx, hero, goodsData, count, ctime)
		if useFunc == nil {
			logrus.Debugf("使用物品，这个物品不符合当前使用条件, %d-%s", goodsData.Id, goodsData.Name)
			result.Add(depot.ERR_USE_GOODS_FAIL_CANT_USE)
			return
		}

		count = u64.Min(count, maxCanUseCount)
		if count <= 0 {
			logrus.Debugf("使用物品，useCount <= 0, %d-%s", goodsData.Id, goodsData.Name)
			result.Add(depot.ERR_USE_GOODS_FAIL_CANT_USE)
			return
		}

		useCount := useFunc(hctx, hero, result, goodsData, count, ctime)

		// 扣物品
		heromodule.ReduceGoodsAnyway(hctx, hero, result, goodsData, useCount)

		result.Add(depot.NewS2cUseGoodsMsg(u64.Int32(id), u64.Int32(useCount)))

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DepotModule) ProcessUseCdrGoods(proto *depot.C2SUseCdrGoodsProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DepotUseCdrGoods)

	id := u64.FromInt32(proto.Id)
	count := u64.FromInt32(proto.Count)
	if count <= 0 {
		logrus.Debugf("使用CDR物品，发送的Count无效")
		hc.Send(depot.ERR_USE_CDR_GOODS_FAIL_INVALID_COUNT)
		return
	}

	goodsData := m.datas.GetGoodsData(id)
	if goodsData == nil {
		logrus.Debugf("使用CDR物品，物品不存在")
		hc.Send(depot.ERR_USE_CDR_GOODS_FAIL_INVALID_ID)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroDepot := hero.Depot()
		if !heroDepot.HasEnoughGoods(id, count) {
			logrus.Debugf("使用CDR物品，发送的Count不足")
			result.Add(depot.ERR_USE_CDR_GOODS_FAIL_COUNT_NOT_ENOUGH)
			return
		}

		if goodsData.GoodsEffect == nil {
			logrus.Debugf("使用CDR物品，这个物品不是减cd物品")
			result.Add(depot.ERR_USE_CDR_GOODS_FAIL_CANT_USE)
			return
		}

		index := int(proto.Index)
		ctime := m.time.CurrentTime()

		//var errMsg []byte

		useCount := count
		switch proto.CdrType {
		default:
			logrus.Debugf("使用CDR物品，CdrType无效, %d-%s", goodsData.Id, goodsData.Name)
			result.Add(depot.ERR_USE_CDR_GOODS_FAIL_INVALID_CDR_TYPE)
			return
		case 0:
			if !goodsData.GoodsEffect.BuildingCdr {
				logrus.Debugf("使用CDR物品，这个物品不是建筑减cd物品")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_INVALID_CDR_TYPE)
				return
			}

			ok, t := hero.Domestic().GetWorkerRestEndTime(index)
			if !ok {
				logrus.Debugf("使用CDR物品，index无效")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_INVALID_INDEX)
				return
			}

			if t.Before(ctime) {
				logrus.Debugf("使用CDR物品，不是cd中")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_NOT_IN_CD)
				return
			}

			d := t.Sub(ctime)

			useCount := u64.Min(useCount, uint64(i64.DivideTimes(d.Nanoseconds(), goodsData.GoodsEffect.Cdr.Nanoseconds())))
			if useCount <= 0 {
				logrus.Debugf("使用CDR物品，useCount <= 0")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_NOT_IN_CD)
				return
			}

			// 扣物品
			heromodule.ReduceGoodsAnyway(hctx, hero, result, goodsData, useCount)

			newTime, _ := hero.Domestic().AddWorkerRestEndTime(index, ctime, -(goodsData.GoodsEffect.Cdr * time.Duration(useCount)))

			// 发消息通知
			result.Add(depot.NewS2cUseCdrGoodsMsg(u64.Int32(id), u64.Int32(useCount), int32(proto.CdrType), int32(index), timeutil.Marshal32(newTime)))

			heromodule.IncreTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_WORKER_SPEED_UP, useCount)

			if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_WORKER_CDR) {
				result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_WORKER_CDR)))
			}

		case 1:
			if !goodsData.GoodsEffect.TechCdr {
				logrus.Debugf("使用CDR物品，这个物品不是科技减cd物品")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_INVALID_CDR_TYPE)
				return
			}

			ok, t := hero.Domestic().GetTechWorkerRestEndTime(index)
			if !ok {
				logrus.Debugf("使用CDR物品，index无效")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_INVALID_INDEX)
				return
			}

			if t.Before(ctime) {
				logrus.Debugf("使用CDR物品，不是cd中")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_NOT_IN_CD)
				return
			}

			d := t.Sub(ctime)

			useCount := u64.Min(useCount, uint64(i64.DivideTimes(d.Nanoseconds(), goodsData.GoodsEffect.Cdr.Nanoseconds())))
			if useCount <= 0 {
				logrus.Debugf("使用CDR物品，useCount <= 0")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_NOT_IN_CD)
				return
			}

			// 扣物品
			heromodule.ReduceGoodsAnyway(hctx, hero, result, goodsData, useCount)

			newTime, _ := hero.Domestic().AddTechWorkerRestEndTime(index, ctime, -(goodsData.GoodsEffect.Cdr * time.Duration(useCount)))

			// 发消息通知
			result.Add(depot.NewS2cUseCdrGoodsMsg(u64.Int32(id), u64.Int32(useCount), int32(proto.CdrType), int32(index), timeutil.Marshal32(newTime)))

			heromodule.IncreTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_WORKER_SPEED_UP, useCount)

			if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_TECH_CDR) {
				result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_TECH_CDR)))
			}

		case 2:
			if !goodsData.GoodsEffect.WorkshopCdr {
				logrus.Debugf("使用CDR物品，这个物品不是装备锻造减cd物品")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_INVALID_CDR_TYPE)
				return
			}

			workshopIndex := hero.Domestic().GetWorkshopIndex()
			if workshopIndex <= 0 {
				logrus.Debugf("使用CDR物品，没有锻造中的装备")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_NOT_IN_CD)
				return
			}

			t := hero.Domestic().GetWorkshopCollectTime()
			if t.Before(ctime) {
				logrus.Debugf("使用CDR物品，不是cd中")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_NOT_IN_CD)
				return
			}

			d := t.Sub(ctime)

			useCount := u64.Min(useCount, uint64(i64.DivideTimes(d.Nanoseconds(), goodsData.GoodsEffect.Cdr.Nanoseconds())))
			if useCount <= 0 {
				logrus.Debugf("使用CDR物品，useCount <= 0")
				result.Add(depot.ERR_USE_CDR_GOODS_FAIL_NOT_IN_CD)
				return
			}

			// 扣物品
			heromodule.ReduceGoodsAnyway(hctx, hero, result, goodsData, useCount)

			newTime := t.Add(-(goodsData.GoodsEffect.Cdr * time.Duration(useCount)))

			hero.Domestic().SetCurrentWorkshop(workshopIndex, newTime)

			result.Add(domestic.NewS2cStartWorkshopMsg(u64.Int32(workshopIndex), timeutil.Marshal32(newTime)))

			// 发消息通知
			result.Add(depot.NewS2cUseCdrGoodsMsg(u64.Int32(id), u64.Int32(useCount), int32(proto.CdrType), int32(index), timeutil.Marshal32(newTime)))

			heromodule.IncreTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_WORKER_SPEED_UP, useCount)

		}

		result.Changed()
		result.Ok()
	})
}

//gogen:iface goods_combine
func (m *DepotModule) ProcessGoodsCombine(proto *depot.C2SGoodsCombineProto, hc iface.HeroController) {
	goodsCombine := m.datas.GetGoodsCombineData(u64.FromInt32(proto.GetId()))

	if goodsCombine == nil {
		logrus.WithField("proto", proto).Debugln("物品合成没找到!")
		hc.Send(depot.ERR_GOODS_COMBINE_FAIL_COMBINE_NOT_FOUND)
		return
	}

	count := u64.FromInt32(proto.GetCount())
	if count <= 0 || count > 1000 {
		logrus.WithField("proto", proto).Debugln("物品合成数量非法，起码要合成一次，最多合成1000个!")
		hc.Send(depot.ERR_GOODS_COMBINE_FAIL_COMBINE_COUNT_INVALID)
		return
	}

	cost := goodsCombine.Cost.Multiple(count)
	hctx := heromodule.NewContext(m.dep, operate_type.DepotGoodsCombine)
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !heromodule.TryReduceCost(hctx, hero, result, cost) {
			logrus.WithField("proto", proto).Debugln("物品合成，物品不够，合成失败!")
			result.Add(depot.ERR_GOODS_COMBINE_FAIL_GOODS_NOT_ENOUGH)
			return
		}

		prize := goodsCombine.Prize.Multiple(count)

		hctx := heromodule.NewContext(m.dep, operate_type.DepotGoodsCombine)
		heromodule.AddPrize(hctx, hero, result, prize, m.time.CurrentTime())

		result.Add(depot.GOODS_COMBINE_S2C)

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DepotModule) ProcessGoodsPartCombine(proto *depot.C2SGoodsPartsCombineProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DepotGoodsPartCombine)

	goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.GetId()))

	if goodsData == nil {
		logrus.WithField("proto", proto).Debugln("物品碎片合成，物品没找到")
		hc.Send(depot.ERR_GOODS_PARTS_COMBINE_FAIL_GOODS_NOT_FOUND)
		return
	}

	if goodsData.EffectType != shared_proto.GoodsEffectType_EFFECT_PARTS {
		logrus.WithField("proto", proto).Debugln("物品碎片合成，发送的不是物品碎片")
		hc.Send(depot.ERR_GOODS_PARTS_COMBINE_FAIL_GOODS_NOT_FOUND)
		return
	}

	if goodsData.GoodsEffect.PartsCombineCount <= 0 {
		logrus.WithField("proto", proto).Error("物品碎片合成，配置的合成所需碎片个数 <= 0")
		hc.Send(depot.ERR_GOODS_PARTS_COMBINE_FAIL_GOODS_NOT_ENOUGH)
		return
	}

	selectIndex := int(proto.SelectIndex)
	if selectIndex < 0 || selectIndex >= len(goodsData.GoodsEffect.PartsPlunderPrizeId) {
		logrus.Debug("物品碎片合成，无效的SelectIndex")
		hc.Send(depot.ERR_GOODS_PARTS_COMBINE_FAIL_INVALID_SELECT_INDEX)
		return
	}

	plunderPrize := m.datas.GetPlunderPrize(goodsData.GoodsEffect.PartsPlunderPrizeId[selectIndex])
	if plunderPrize == nil {
		logrus.Error("物品碎片合成，plunderPrize == nil")
		hc.Send(depot.ERR_GOODS_PARTS_COMBINE_FAIL_INVALID_SELECT_INDEX)
		return
	}

	count := u64.FromInt32(proto.GetCount())
	if count <= 0 || count > 1000 {
		logrus.WithField("proto", proto).Debugln("物品碎片合成数量非法，起码要合成一次，最多合成1000个!")
		hc.Send(depot.ERR_GOODS_PARTS_COMBINE_FAIL_GOODS_NOT_ENOUGH)
		return
	}

	partCount := count * goodsData.GoodsEffect.PartsCombineCount

	//b := resdata.NewPrizeBuilder()
	//
	//if combineGoodsData := m.datas.GetGoodsData(goodsData.GoodsEffect.PartsCombineId); combineGoodsData != nil {
	//	b.AddGoods(combineGoodsData, count)
	//} else if combineGoodsData := m.datas.GetEquipmentData(goodsData.GoodsEffect.PartsCombineId); combineGoodsData != nil {
	//	b.AddEquipment(combineGoodsData, count)
	//} else if combineGoodsData := m.datas.GetGemData(goodsData.GoodsEffect.PartsCombineId); combineGoodsData != nil {
	//	b.AddGem(combineGoodsData, count)
	//} else {
	//	logrus.WithField("proto", proto).Error("物品碎片合成，碎片物品没找到合成物品")
	//	hc.Send(depot.ERR_GOODS_PARTS_COMBINE_FAIL_GOODS_NOT_FOUND)
	//	return
	//}

	b := resdata.NewPrizeBuilder()

	for i := uint64(0); i < count; i++ {
		b.Add(plunderPrize.GetPrize())
	}
	combinePrize := b.Build()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		//if !hero.Depot().HasEnoughGoods(goodsData.Id, partCount) {
		if !heromodule.TryReduceGoods(hctx, hero, result, goodsData, partCount) {
			logrus.Debug("物品碎片合成，英雄的碎片不足")
			result.Add(depot.ERR_GOODS_PARTS_COMBINE_FAIL_GOODS_NOT_ENOUGH)
			return
		}

		ctime := m.time.CurrentTime()
		hctx := heromodule.NewContext(m.dep, operate_type.DepotGoodsPartCombine)
		heromodule.AddPrize(hctx, hero, result, combinePrize, ctime)

		result.Add(depot.NewS2cGoodsPartsCombineMsg(proto.Id, proto.Count, proto.SelectIndex, must.Marshal(combinePrize)))

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DepotModule) ProcessUnlockBaowu(proto *depot.C2SUnlockBaowuProto, hc iface.HeroController) {

	data := m.datas.GetBaowuData(u64.FromInt32(proto.Id))
	if data == nil {
		logrus.WithField("id", proto.Id).Debug("解锁宝物，无效的宝物id")
		hc.Send(depot.ERR_UNLOCK_BAOWU_FAIL_NOT_FOUND)
		return
	}

	ctime := m.time.CurrentTime()
	guildId := int64(0)
	var heroName string
	var heroIdbytes []byte
	var heroHead string
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.Depot().GetUnlockBaowuData() != nil {
			logrus.WithField("id", proto.Id).Debug("解锁宝物，正在解锁其他宝物")
			hc.Send(depot.ERR_UNLOCK_BAOWU_FAIL_UNLOCKING)
			return
		}

		if !heromodule.TryReduceBaowu(hero, result, data, 1) {
			logrus.WithField("id", proto.Id).Debug("解锁宝物，英雄没有这个宝物")
			hc.Send(depot.ERR_UNLOCK_BAOWU_FAIL_NOT_FOUND)
			return
		}

		hero.Depot().UnlockBaowu(data, ctime)

		result.Add(depot.NewS2cUnlockBaowuMsg(u64.Int32(data.Id),
			timeutil.Marshal32(hero.Depot().GetUnlockBaowuEndTime())))

		hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_UnlockBaowu, data.Level)
		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_UnlockBaowu)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UNLOCK_BAOWU)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_UNLOCK_BAOWU) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_UNLOCK_BAOWU)))
		}

		guildId = hero.GuildId()
		heroName = hero.Name()
		heroIdbytes = hero.IdBytes()
		heroHead = hero.Head()
		result.Ok()
	})
	if guildId == 0 || data.Prestige == 0 {
		return
	}
	// 增加联盟声望
	m.dep.Guild().FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}
		if member := g.GetMember(hc.Id()); member != nil && member.ClassLevelData().CorePrestige {
			g.AddPrestigeCore(data.Prestige)
		}
		g.AddPrestige(data.Prestige)
		m.dep.Country().AddPrestige(g.Country().Id, data.Prestige)
		m.rankModule.AddOrUpdateRankObj(ranklist.NewGuildRankObj(m.dep.Guild().GetSnapshot, m.dep.HeroSnapshot().Get, g))

		if logData := m.datas.GuildLogHelp().BaowuAddPrestige; logData != nil {
			proto := logData.NewHeroLogProto(ctime, heroIdbytes, heroHead)
			proto.Text = logData.Text.New().WithHeroName(heroName).WithItem(data.Name).WithAmount(data.Prestige).JsonString()
			m.dep.Guild().AddLog(guildId, proto)
		}
	})
}

//gogen:iface
func (m *DepotModule) ProcessCollectBaowu(proto *depot.C2SCollectBaowuProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DepotCollectBaowu)

	var baowuLevel uint64
	var guildId int64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		unlockData := hero.Depot().GetUnlockBaowuData()
		if unlockData == nil {
			logrus.Debug("开启宝物，没有解锁中的宝物")
			hc.Send(depot.ERR_COLLECT_BAOWU_FAIL_NO_UNLOCKING)
			return
		}

		ctime := m.time.CurrentTime()
		diff := hero.Depot().GetUnlockBaowuEndTime().Sub(ctime)

		if diff > 0 {
			if proto.Miao {
				if limit := m.datas.MiscGenConfig().DailyMiaoBaowuLimit;
					limit > 0 && hero.Depot().GetMiaoBaowuTimes() >= limit {
					logrus.Debug("开启宝物，没有解锁中的宝物")
					hc.Send(depot.ERR_COLLECT_BAOWU_FAIL_NO_UNLOCKING)
					return
				}

				// 扣钱
				times := i64.DivideTimes(diff.Nanoseconds(), unlockData.MiaoDuration.Nanoseconds())
				if times <= 0 {
					logrus.Debug("开启宝物，开启时间未到，不能领取")
					hc.Send(depot.ERR_COLLECT_BAOWU_FAIL_UNLOCKING)
					return
				}

				cost := m.datas.MiscGenConfig().MiaoBaowuCost.Multiple(uint64(times))
				if !heromodule.TryReduceCost(hctx, hero, result, cost) {
					logrus.Debug("开启宝物，秒CD消耗不足")
					hc.Send(depot.ERR_COLLECT_BAOWU_FAIL_COST_NOT_ENOUGH)
					return
				}
			} else {
				logrus.Debug("开启宝物，开启时间未到，不能领取")
				hc.Send(depot.ERR_COLLECT_BAOWU_FAIL_UNLOCKING)
				return
			}

			hero.Depot().IncMiaoBaowuTimes()
		}

		// 清空奖励
		hero.Depot().ClearUnlockBaowu()

		// 给奖励
		toAdd := unlockData.PlunderPrize.GetPrize()

		hctx := heromodule.NewContext(m.dep, operate_type.DepotCollectBaowu)
		heromodule.AddPrize(hctx, hero, result, toAdd, ctime)

		// 发消息
		result.Add(depot.NewS2cCollectBaowuMsg(toAdd.Encode()))

		guildId = hero.GuildId()
		baowuLevel = unlockData.Level

		result.Changed()
		result.Ok()
	})

	if baowuLevel > 0 { // 开启成功，更新联盟任务进度
		data := m.datas.GetGuildTaskData(u64.FromInt32(int32(server_proto.GuildTaskType_BaoWu)))
		m.dep.Guild().AddGuildTaskProgress(guildId, data, baowuLevel)
	}
}

//gogen:iface c2s_list_baowu_log
func (m *DepotModule) ProcessListBaowuLog(hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		toSendProto := &depot.S2CListBaowuLogProto{}
		hero.Depot().RangeBaowuLogWithStartIndex(constants.BaowuLogPreviewCount, func(log []byte) (toContinue bool) {
			toSendProto.Datas = append(toSendProto.Datas, log)
			return true
		})

		result.Add(depot.NewS2cListBaowuLogProtoMsg(toSendProto))

		result.Ok()
	})

}

//gogen:iface
func (m *DepotModule) ProcessDecomposeBaowu(proto *depot.C2SDecomposeBaowuProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		baowuId := u64.FromInt32(proto.BaowuId)
		toRemoveCount := u64.FromInt32(proto.Count)

		baowuData := m.datas.GetBaowuData(baowuId)
		if baowuData == nil {
			result.Add(depot.ERR_DECOMPOSE_BAOWU_FAIL_INVALID_BAOWU)
			return
		}

		allCount := hero.Depot().GetBaowuCount(baowuId)
		if allCount <= 0 {
			result.Add(depot.ERR_DECOMPOSE_BAOWU_FAIL_INVALID_BAOWU)
			return
		}

		if toRemoveCount > allCount {
			result.Add(depot.ERR_DECOMPOSE_BAOWU_FAIL_INVALID_COUNT)
			return
		}

		if unlockBaowu := hero.Depot().GetUnlockBaowuData(); unlockBaowu != nil {
			if unlockBaowu.Id == baowuId {
				result.Add(depot.ERR_DECOMPOSE_BAOWU_FAIL_BAOWU_IS_UNLOCKING)
				return
			}
		}

		hctx := heromodule.NewContext(m.dep, operate_type.BaoWuDecompose)
		newCount := hero.Depot().RemoveBaowu(baowuId, toRemoveCount)
		heromodule.AddUnsafeResource(hctx, hero, result, baowuData.DecomposeGold, 0, 0, baowuData.DecomposeStone)

		result.Add(depot.NewS2cDecomposeBaowuMsg(proto.BaowuId, u64.Int32(newCount)))
		result.Add(depot.NewS2cUpdateBaowuMsg(proto.BaowuId, u64.Int32(newCount)))

		for level := baowuData.Level; level > 0; level-- {
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_SellBaowu, level, toRemoveCount)
		}
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BAOWU_SELL)

		result.Changed()
		result.Ok()
	})
}
