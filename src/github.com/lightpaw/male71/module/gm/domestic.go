package gm

import (
	"fmt"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/guild"
	domestic2 "github.com/lightpaw/male7/module/domestic"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"math/rand"
	"strings"
	"time"
	"github.com/lightpaw/male7/gen/pb/military"
)

// 内政
func (m *GmModule) newDomesticGmGroup() *gm_group {
	return &gm_group{
		tab: "内政",
		handler: []*gm_handler{
			newHeroIntHandler("清空资源", "", m.clearResources),
			newHeroIntHandler("资源全加", "100000", m.addResources),
			newHeroIntHandler("加钱(负数表示减钱)", "100000", m.addGold),
			//newHeroIntHandler("加粮食(负数表示减)", "100000", m.addFood),
			//newHeroIntHandler("加木材(负数表示减)", "100000", m.addWood),
			newHeroIntHandler("加石材(负数表示减)", "100000", m.addStone),
			newHeroIntHandler("加玉石矿(负数表示减)", "100000", m.addJadeOre),
			newHeroIntHandler("加玉璧(负数表示减)", "100000", m.addJade),
			newHeroIntHandler("加联盟贡献币(负数表示减)", "100000", m.addGuildContributionCoin),
			newHeroIntHandler("所有建筑升1级", "", m.upgradeBuild1Level),
			newHeroIntHandler("建筑升到X级", "19", m.upgradeBuildToLevelX),
			newStringHandler("收获资源", "gold", m.collectResource),
			newHeroIntHandler("满兵", "", m.fullSoldier),
			newHeroIntHandler("清兵", "", m.clearSoldier),
			newHeroStringHandler("加建筑队CD", "", m.addWorkerCD),
			newHeroStringHandler("加锻造次数", "", m.addForgeTimes),
			newIntHandler("使用增益", "1", m.advantageUse),
			newIntHandler("使用buff给自己", "1", m.buffToSelf),
			newIntHandler("清除buff", "1", m.clearBuff),
			newStringHandler("修炼馆可领经验", "", m.canCollectCapExp),
		},
	}
}

// 清空资源
func (m *GmModule) clearResources(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	unsafe := hero.GetUnsafeResource()
	heromodule.TryReduceResource(m.hctx, hero, result, unsafe.Gold(), unsafe.Food(), unsafe.Wood(), unsafe.Stone())
	safe := hero.GetSafeResource()
	heromodule.TryReduceResource(m.hctx, hero, result, safe.Gold(), safe.Food(), safe.Wood(), safe.Stone())
}

// 资源全加
func (m *GmModule) addResources(amt int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	amount := u64.FromInt64(amt)
	heromodule.AddUnsafeResource(m.hctx, hero, result, amount, amount, amount, amount)
}

// 加钱(负数表示减钱)
func (m *GmModule) addGold(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	m.changeSingleResource(shared_proto.ResType_GOLD, amount, hero, result)
}

// 加粮食(负数表示减粮食)
func (m *GmModule) addFood(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	m.changeSingleResource(shared_proto.ResType_FOOD, amount, hero, result)
}

// 加木材(负数表示减木材)
func (m *GmModule) addWood(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	m.changeSingleResource(shared_proto.ResType_WOOD, amount, hero, result)
}

// 加石材(负数表示减石材)
func (m *GmModule) addStone(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	m.changeSingleResource(shared_proto.ResType_STONE, amount, hero, result)
}

func (m *GmModule) changeSingleResource(resType shared_proto.ResType, amount int64, hero *entity.Hero, result herolock.LockResult) {
	if amount >= 0 {
		heromodule.AddUnsafeSingleResource(m.hctx, hero, result, resType, uint64(amount))
	} else {
		heromodule.TryReduceSingleResource(m.hctx, hero, result, resType, uint64(-amount))
	}
}

// 加玉石矿(负数表示减)
func (m *GmModule) addJadeOre(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	if amount >= 0 {
		heromodule.AddJadeOre(m.hctx, hero, result, uint64(amount))
	} else {
		heromodule.ReduceJadeOreAnyway(m.hctx, hero, result, uint64(-amount))
	}
}

// 加玉璧(负数表示减)
func (m *GmModule) addJade(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	if amount >= 0 {
		heromodule.AddJade(m.hctx, hero, result, uint64(amount))
	} else {
		heromodule.ReduceJadeAnyway(m.hctx, hero, result, uint64(-amount))
	}
}

// 加元宝(负数表示减)
func (m *GmModule) addYuanbao(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	if amount >= 0 {
		heromodule.AddYuanbao(m.hctx, hero, result, uint64(amount))

		// 触发充值任务
		for _, data := range m.datas.GetGuildEventPrizeDataArray() {
			if data.TriggerEvent == shared_proto.HeroEvent_HERO_EVENT_RECHARGE {
				heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_RECHARGE, data.TriggerEventCondition.Amount)
			}
		}
		heromodule.AddYuanbaoGiftLimit(m.hctx, hero, result, uint64(amount))
	} else {
		heromodule.ReduceYuanbaoAnyway(m.hctx, hero, result, uint64(-amount))
		heromodule.ReduceYuanbaoGiftLimit(m.hctx, hero, result, uint64(-amount))
	}
}

// 加点券(负数表示减)
func (m *GmModule) addDianquan(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	if amount >= 0 {
		heromodule.AddDianquan(m.hctx, hero, result, uint64(amount))
	} else {
		heromodule.ReduceDianquanAnyway(m.hctx, hero, result, uint64(-amount))
	}
}

// 加yinliang(负数表示减)
func (m *GmModule) addYinliang(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	if amount >= 0 {
		heromodule.AddYinliang(m.hctx, hero, result, uint64(amount))
	} else {
		heromodule.ReduceYinliangAnyway(m.hctx, hero, result, uint64(-amount))
	}
}

// 加联盟贡献币(负数表示减)
func (m *GmModule) addGuildContributionCoin(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	if amount >= 0 {
		hero.AddGuildContributionCoin(uint64(amount))
	} else {
		hero.ReduceGuildContributionCoin(uint64(-amount))
	}
	result.Add(guild.NewS2cUpdateContributionCoinMsg(u64.Int32(hero.GetGuildContributionCoin())))
}

// 所有建筑升1级
func (m *GmModule) upgradeBuild1Level(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	heroDomestic := hero.Domestic()

	for t := range shared_proto.BuildingType_name {
		if t <= 0 {
			continue
		}

		buildingType := shared_proto.BuildingType(t)
		if domestic_data.IsResourceBuilding(buildingType) {
			continue
		}

		currentBuilding := heroDomestic.GetBuilding(buildingType)
		nextLevelAmount := uint64(1)
		if currentBuilding == nil {
			continue
		}

		nextLevelAmount = currentBuilding.Level + 1

		nextLevelId := domestic_data.BuildingId(buildingType, nextLevelAmount)
		nextLevel := m.datas.GetBuildingData(nextLevelId)
		if nextLevel == nil {
			continue
		}

		var prosperity uint64

		heroDomestic.SetBuilding(nextLevel)

		toAddProsperity := nextLevel.Prosperity
		if currentBuilding != nil {
			toAddProsperity = u64.Sub(toAddProsperity, currentBuilding.Prosperity)
		}

		if toAddProsperity > 0 {
			hero.AddProsperityCapcity(toAddProsperity)
			prosperity = hero.Prosperity()
		}

		heromodule.AddExp(m.hctx, hero, result, nextLevel.HeroExp, m.time.CurrentTime())

		_, wt := heroDomestic.GetWorkerRestEndTime(0)
		result.Add(domestic.NewS2cUpgradeStableBuildingMsg(u64.Int32(nextLevel.Id), int32(0), timeutil.Marshal32(wt)))

		if prosperity > 0 {
			result.Add(domestic.NewS2cHeroUpdateProsperityMsg(u64.Int32(prosperity), u64.Int32(heroDomestic.ProsperityCapcity())))
		}

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL)
		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL)
	}

	outerCities := heroDomestic.OuterCities()
	for _, outerCityData := range m.datas.GetOuterCityDataArray() {
		if !hero.Function().IsFunctionOpened(outerCityData.GetFuncType()) {
			continue
		}

		outerCity := outerCities.OuterCity(outerCityData)
		if outerCity == nil {
			outerCity = entity.NewOuterCity(outerCityData, outerCityData.Id%2)
			outerCities.Unlock(outerCity)
		}
	}
	outerCities.UpdateUnlockBit()

	for _, layoutData := range m.datas.GetOuterCityLayoutDataArray() {
		outerCity := outerCities.OuterCity(layoutData.OuterCity)
		if outerCity == nil {
			continue
		}

		var nextLevel *domestic_data.OuterCityLayoutData
		curLayout := outerCity.Layout(layoutData)
		if curLayout != nil {
			nextLevel = curLayout.NextLevel

			// 建筑已经最高级
			if nextLevel == nil {
				continue
			}
		} else {
			nextLevel = layoutData
			if layoutData.Level != 1 {
				continue
			}
		}

		// 主城等级要求
		if nextLevel.GetBuilding(outerCity.Type()).BuildingData.BaseLevel != nil && hero.BaseLevel() < nextLevel.GetBuilding(outerCity.Type()).BuildingData.BaseLevel.Level {
			continue
		}

		// 前提条件
		for _, requireBuilding := range nextLevel.GetBuilding(outerCity.Type()).BuildingData.RequireIds {
			heroBuilding := heroDomestic.GetBuilding(requireBuilding.Type)
			if heroBuilding == nil || heroBuilding.Level < requireBuilding.Level {
				continue
			}
		}

		for _, requireBuilding := range nextLevel.UpgradeRequireIds {
			heroBuilding := heroDomestic.GetBuilding(requireBuilding.Type)
			if heroBuilding == nil || heroBuilding.Level < requireBuilding.Level {
				continue
			}
		}

		if layoutData.UpgradeRequireLayout != nil {
			// 升级需要其他的布局
			if layout := outerCity.Layout(layoutData.UpgradeRequireLayout); layout == nil || layout.GetBuilding(outerCity.Type()).BuildingData.Level < layoutData.UpgradeRequireLayout.Level {
				continue
			}
		}

		outerCity.SetLayout(nextLevel)

		toAddProsperity := nextLevel.GetBuilding(outerCity.Type()).BuildingData.Prosperity
		if curLayout != nil {
			toAddProsperity = u64.Sub(toAddProsperity, curLayout.GetBuilding(outerCity.Type()).BuildingData.Prosperity)
		}

		if toAddProsperity > 0 {
			hero.AddProsperityCapcity(toAddProsperity)
		}

		heromodule.AddExp(m.hctx, hero, result, nextLevel.GetBuilding(outerCity.Type()).BuildingData.HeroExp, m.time.CurrentTime())

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL)
		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL)

		result.Add(domestic.NewS2cUpgradeOuterCityBuildingMsg(u64.Int32(layoutData.OuterCity.Id), u64.Int32(layoutData.Id), u64.Int32(nextLevel.Id)))

		result.Changed()
		result.Ok()
	}
}

// 所有建筑升1级
func (m *GmModule) upgradeBuildToLevelX(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	if amount < 0 {
		return
	}

	heroDomestic := hero.Domestic()

	targetLevel := uint64(amount)

	for i := 0; i < int(targetLevel); i++ {
		for t := range shared_proto.BuildingType_name {
			if t <= 0 {
				continue
			}

			buildingType := shared_proto.BuildingType(t)
			if domestic_data.IsResourceBuilding(buildingType) {
				continue
			}

			currentBuilding := heroDomestic.GetBuilding(buildingType)
			nextLevelAmount := uint64(1)
			if currentBuilding == nil {
				continue
			}

			nextLevelAmount = currentBuilding.Level + 1

			nextLevelId := domestic_data.BuildingId(buildingType, nextLevelAmount)
			nextLevel := m.datas.GetBuildingData(nextLevelId)
			if nextLevel == nil {
				continue
			}

			if nextLevel.Level > targetLevel {
				continue
			}

			var prosperity uint64

			heroDomestic.SetBuilding(nextLevel)

			toAddProsperity := nextLevel.Prosperity
			if currentBuilding != nil {
				toAddProsperity = u64.Sub(toAddProsperity, currentBuilding.Prosperity)
			}

			if toAddProsperity > 0 {
				hero.AddProsperityCapcity(toAddProsperity)
				prosperity = hero.Prosperity()
			}

			heromodule.AddExp(m.hctx, hero, result, nextLevel.HeroExp, m.time.CurrentTime())

			_, wt := heroDomestic.GetWorkerRestEndTime(0)
			result.Add(domestic.NewS2cUpgradeStableBuildingMsg(u64.Int32(nextLevel.Id), int32(0), timeutil.Marshal32(wt)))

			if prosperity > 0 {
				result.Add(domestic.NewS2cHeroUpdateProsperityMsg(u64.Int32(prosperity), u64.Int32(heroDomestic.ProsperityCapcity())))
			}

			// 更新任务进度
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL)
			heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL)
		}

		outerCities := heroDomestic.OuterCities()
		for _, outerCityData := range m.datas.GetOuterCityDataArray() {
			if !hero.Function().IsFunctionOpened(outerCityData.GetFuncType()) {
				continue
			}

			outerCity := outerCities.OuterCity(outerCityData)
			if outerCity == nil {
				outerCity = entity.NewOuterCity(outerCityData, outerCityData.Id%2)
				outerCities.Unlock(outerCity)
			}
		}
		outerCities.UpdateUnlockBit()

		for _, layoutData := range m.datas.GetOuterCityLayoutDataArray() {
			outerCity := outerCities.OuterCity(layoutData.OuterCity)
			if outerCity == nil {
				continue
			}

			var nextLevel *domestic_data.OuterCityLayoutData
			curLayout := outerCity.Layout(layoutData)
			if curLayout != nil {
				nextLevel = curLayout.NextLevel

				// 建筑已经最高级
				if nextLevel == nil {
					continue
				}
			} else {
				nextLevel = layoutData
				if layoutData.Level != 1 {
					continue
				}
			}

			if nextLevel.Level > targetLevel {
				continue
			}

			// 主城等级要求
			if nextLevel.GetBuilding(outerCity.Type()).BuildingData.BaseLevel != nil && hero.BaseLevel() < nextLevel.GetBuilding(outerCity.Type()).BuildingData.BaseLevel.Level {
				continue
			}

			// 前提条件
			for _, requireBuilding := range nextLevel.GetBuilding(outerCity.Type()).BuildingData.RequireIds {
				heroBuilding := heroDomestic.GetBuilding(requireBuilding.Type)
				if heroBuilding == nil || heroBuilding.Level < requireBuilding.Level {
					continue
				}
			}

			for _, requireBuilding := range nextLevel.UpgradeRequireIds {
				heroBuilding := heroDomestic.GetBuilding(requireBuilding.Type)
				if heroBuilding == nil || heroBuilding.Level < requireBuilding.Level {
					continue
				}
			}

			if layoutData.UpgradeRequireLayout != nil {
				// 升级需要其他的布局
				if layout := outerCity.Layout(layoutData.UpgradeRequireLayout); layout == nil || layout.GetBuilding(outerCity.Type()).BuildingData.Level < layoutData.UpgradeRequireLayout.Level {
					continue
				}
			}

			outerCity.SetLayout(nextLevel)

			toAddProsperity := nextLevel.GetBuilding(outerCity.Type()).BuildingData.Prosperity
			if curLayout != nil {
				toAddProsperity = u64.Sub(toAddProsperity, curLayout.GetBuilding(outerCity.Type()).BuildingData.Prosperity)
			}

			if toAddProsperity > 0 {
				hero.AddProsperityCapcity(toAddProsperity)
			}

			heromodule.AddExp(m.hctx, hero, result, nextLevel.GetBuilding(outerCity.Type()).BuildingData.HeroExp, m.time.CurrentTime())

			// 更新任务进度
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL)
			heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL)

			result.Add(domestic.NewS2cUpgradeOuterCityBuildingMsg(u64.Int32(layoutData.OuterCity.Id), u64.Int32(layoutData.Id), u64.Int32(nextLevel.Id)))

			result.Changed()
			result.Ok()
		}
	}
}

// 收获资源
func (m *GmModule) collectResource(amount string, hc iface.HeroController) {
	var f func()

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		layoutData := m.datas.BuildingLayoutData().MinKeyData
		layoutId := layoutData.Id

		amount = strings.TrimSpace(amount)
		if len(amount) <= 0 {
			switch rand.Intn(4) {
			case 0:
				amount = "gold"
			case 1:
				amount = "food"
			case 2:
				amount = "wood"
			case 3:
				amount = "stone"
			}
		}

		buildingType := shared_proto.BuildingType_GOLD_PRODUCER
		switch strings.ToLower(amount) {
		case "food":
			buildingType = shared_proto.BuildingType_FOOD_PRODUCER
		case "wood":
			buildingType = shared_proto.BuildingType_WOOD_PRODUCER
		case "stone":
			buildingType = shared_proto.BuildingType_STONE_PRODUCER
		}

		buildingId := domestic_data.BuildingId(buildingType, 1)
		building := m.datas.GetBuildingData(buildingId)
		if building == nil {
			return
		}

		ctime := m.time.CurrentTime()

		resourcePoint := hero.Domestic().GetLayoutRes(layoutId)
		if resourcePoint == nil {
			resourcePoint = hero.Domestic().SetResourcePoint(layoutData, building, ctime.Add(-time.Hour))
		}

		resourcePoint.SetBuilding(building)
		resourcePoint.AddOutputAmount(100)

		f = func() {
			m.modules.DomesticModule().(*domestic2.DomesticModule).ProcessCollectResource(&domestic.C2SCollectResourceProto{
				Id: u64.Int32(layoutId),
			}, hc)
		}
	})

	if f != nil {
		f()
	}
}

// 清兵
func (m *GmModule) fullSoldier(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {

	ctime := m.time.CurrentTime()
	hero.Military().AddFreeSoldier(u64.Sub(hero.BuildingEffect().SoldierCapcity(), hero.Military().FreeSoldier(ctime)), ctime)
	result.Add(hero.Military().NewUpdateFreeSoldierMsg())

	for _, troop := range hero.Troops() {
		if troop.IsOutside() {
			continue
		}

		for _, pos := range troop.Pos() {
			c := pos.Captain()
			if c != nil {
				c.AddSoldier(u64.Sub(c.SoldierCapcity(), c.Soldier()))
				result.Add(c.NewUpdateCaptainStatMsg())
				heromodule.UpdateTroopFightAmount(hero, c.GetTroop(), result)

				newSoldier := c.Soldier()
				fightAmount := c.FightAmount()
				result.Add(military.NewS2cCaptainChangeSoldierMsg(u64.Int32(c.Id()), u64.Int32(newSoldier), u64.Int32(fightAmount), u64.Int32(hero.Military().FreeSoldier(ctime))))
			}
		}
	}

	result.Changed()
}

// 清兵
func (m *GmModule) clearSoldier(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {

	ctime := m.time.CurrentTime()
	hero.Military().ReduceFreeSoldier(hero.Military().FreeSoldier(ctime), ctime)
	result.Add(hero.Military().NewUpdateFreeSoldierMsg())

	for _, troop := range hero.Troops() {
		if troop.IsOutside() {
			continue
		}

		for _, pos := range troop.Pos() {
			c := pos.Captain()
			if c != nil {
				c.ReduceSoldier(c.Soldier())
				result.Add(c.NewUpdateCaptainStatMsg())
				heromodule.UpdateTroopFightAmount(hero, c.GetTroop(), result)

				newSoldier := c.Soldier()
				fightAmount := c.FightAmount()
				result.Add(military.NewS2cCaptainChangeSoldierMsg(u64.Int32(c.Id()), u64.Int32(newSoldier), u64.Int32(fightAmount), u64.Int32(hero.Military().FreeSoldier(ctime))))
			}
		}
	}

	result.Changed()
}

// 加建筑队CD
func (m *GmModule) addWorkerCD(amount string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {

	ctime := m.time.CurrentTime()
	if hero.Domestic().GetBuildingWorkerFatigueDuration() > 0 {
		if pos := hero.Domestic().GetFreeWorker(ctime); pos >= 0 {
			heromodule.UseBuildingWorkerTime(hero, result, time.Hour, ctime)
		}
	}

	if hero.Domestic().GetTechWorkerFatigueDuration() > 0 {
		if pos := hero.Domestic().GetFreeTechWorker(ctime); pos >= 0 {
			heromodule.UseTechWorkerTime(hero, result, time.Hour, ctime)
		}
	}

	result.Changed()
}

// 加锻造次数
func (m *GmModule) addForgeTimes(amount string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {

	b := hero.Domestic().GetBuilding(shared_proto.BuildingType_TIE_JIANG_PU)
	if b != nil {
		hero.Domestic().GetForgingTimes().SetTimes(0)
		hero.Domestic().GetForgingTimes().SetNextTime(time.Time{})
		result.Add(domestic.NewS2cRecoveryForgingTimeChangeMsg(0, 0))
	}

	result.Changed()
}

// 使用主城增益
func (m *GmModule) advantageUse(id int64, hc iface.HeroController) {
	d := m.datas.GetBufferData(u64.FromInt64(id))
	if d == nil {
		hc.Send(domestic.ERR_USE_ADVANTAGE_FAIL_INVALID_ID)
		return
	}

	if d.TypeData.IsMian {
		m.modules.DomesticModule().UseMianGoods(d, u64.FromInt64(id), hc)
	} else {
		m.modules.DomesticModule().UseBuffGoods(d, u64.FromInt64(id), hc)
	}
}


func (m *GmModule) buffToSelf(bid int64, hc iface.HeroController) {
	buff := m.datas.GetBuffEffectData(u64.FromInt64(bid))
	succ := m.buffService.AddBuffToSelf(buff, hc.Id())
	fmt.Printf("=== buffToSelf %v %v \n", buff.Id, succ)
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.Buff().Walk(func(buff *entity.BuffInfo) {
			fmt.Printf("buff:%+v %+v \n", buff.EffectData, buff)
		})
	})
}

func (m *GmModule) clearBuff(gid int64, hc iface.HeroController) {
	fmt.Printf("clear buff:%v \n", gid)
	m.buffService.CancelGroup(hc.Id(), u64.FromInt64(gid))

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.Buff().Walk(func(buff *entity.BuffInfo) {
			fmt.Printf("buff:%+v %+v \n", buff.EffectData, buff)
		})
	})
}

func (m *GmModule) canCollectCapExp(str string, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		buildingExp := heromodule.CalcTrainBuildingExp(hero, m.datas, m.time.CurrentTime())
		currBuffExp := heromodule.CalcTrainCurrentBuffExp(hero, m.dep.Datas(), m.dep.Time().CurrentTime())
		resevdBuffExp := hero.Military().ReservedExp()

		fmt.Printf("====== can collect cap tran building exp:%v\n", buildingExp)
		fmt.Printf("====== can collect cap tran curr buff exp:%v\n", currBuffExp)
		fmt.Printf("====== can collect cap tran reserved buff exp:%v\n", resevdBuffExp)
		fmt.Printf("====== can collect cap tran all buff exp:%v\n",  buildingExp + u64.FromInt64(currBuffExp + resevdBuffExp))
	})
}
