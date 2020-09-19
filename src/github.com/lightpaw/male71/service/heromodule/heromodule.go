package heromodule

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/captain"
	country2 "github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/config/function"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/maildata"
	"github.com/lightpaw/male7/config/promdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/season"
	shopdata "github.com/lightpaw/male7/config/shop"
	"github.com/lightpaw/male7/config/vip"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/country"
	"github.com/lightpaw/male7/gen/pb/depot"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/promotion"
	"github.com/lightpaw/male7/gen/pb/random_event"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/gen/pb/shop"
	vipMsg "github.com/lightpaw/male7/gen/pb/vip"
	"github.com/lightpaw/male7/gen/pb/xuanyuan"
	"github.com/lightpaw/male7/gen/pb/zhengwu"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"math"
	"sync/atomic"
	"time"
	"github.com/lightpaw/male7/gamelogs"
)

func TryExchange(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult,
	exchangeData *resdata.ExchangeData, count uint64, ctime time.Time) (success bool) {

	if count <= 0 {
		return true
	}

	totalCost := exchangeData.Cost.Multiple(count)

	if !TryReduceCost(hctx, hero, result, totalCost) {
		logrus.Debugf("物品兑换，消耗不足")
		return false
	}

	prize := exchangeData.Prize.Multiple(count)

	AddPrize(hctx, hero, result, prize, ctime)

	return true
}

func TryResetXuanyuan(hc iface.HeroController, resetTime time.Time) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Xuanyuan().TryResetDaily(resetTime) {
			result.Add(xuanyuan.RESET_S2C)
		}
	})
}

func TryResetWeekly(hc iface.HeroController, ctime, resetTime time.Time) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !hero.GetWeeklyResetTime().Before(resetTime) {
			return
		}
		t := timeutil.Max(ctime, resetTime)
		hero.ResetWeekly(t)
		result.Add(misc.RESET_WEEKLY_S2C)
		result.Ok()
	})
}

func TryResetDailyZero(hc iface.HeroController, ctime, resetTime time.Time, services *Services) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !hero.GetDailyZeroResetTime().Before(resetTime) {
			return
		}

		t := timeutil.Max(ctime, resetTime)

		hero.ResetDailyZero(t, services.datas)
		result.Add(misc.RESET_DAILY_ZERO_S2C)

		// 凌晨刷新 vip
		hctx := NewContext(services.dep, operate_type.HeroLoginLoad)
		VipResetDailyZero(hctx, hero, result, services.dep.Datas(), t)

		hero.TaskList().ActiveDegreeTaskList().Walk(func(taskId uint64, task *entity.ActiveDegreeTask) bool {
			result.Add(task.NewUpdateTaskProgressMsg())
			return false
		})

		result.Changed()
	})
}

func TryResetDailyMc(hc iface.HeroController, ctime, resetTime time.Time, services *Services) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !hero.GetDailyMcResetTime().Before(resetTime) {
			return
		}

		t := timeutil.Max(ctime, resetTime)

		hero.ResetDailyMc(t, services.datas)
		result.Add(misc.RESET_DAILY_MC_S2C)

		hero.TaskList().ActiveDegreeTaskList().Walk(func(taskId uint64, task *entity.ActiveDegreeTask) bool {
			result.Add(task.NewUpdateTaskProgressMsg())
			return false
		})

		result.Changed()
	})
}

func TryResetDaily(hc iface.HeroController, ctime, resetTime time.Time, services *Services) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroResetDaily(hero, result, ctime, resetTime, services.datas)
		// tlog
		services.dep.Tlog().TlogResourceStockFlow(hero, hero.GetYuanbao(), hero.GetDianquan(), hero.Military().FreeSoldier(ctime), hero.GetAllRes(shared_proto.ResType_GOLD), hero.GetAllRes(shared_proto.ResType_STONE), hero.GetGuildContributionCoin())
	})
}

func heroResetDaily(hero *entity.Hero, result herolock.LockResult, ctime, resetTime time.Time, datas iface.ConfigDatas) bool {
	if hero.GetDailyResetTime().Before(resetTime) {
		if ctime.After(resetTime) {
			hero.ResetDaily(ctime, datas)
		} else {
			hero.ResetDaily(resetTime, datas)
		}
		result.Add(misc.RESET_DAILY_S2C)

		UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_LOGIN_DAY)

		hero.TaskList().ActiveDegreeTaskList().Walk(func(taskId uint64, task *entity.ActiveDegreeTask) bool {
			result.Add(task.NewUpdateTaskProgressMsg())
			return false
		})

		result.Changed()
		return true
	}
	return false
}

func TryResetSeason(hc iface.HeroController, ctime, resetTime time.Time, seasonData *season.SeasonData) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		HeroResetSeason(hero, result, ctime, resetTime, seasonData)
	})
}

func HeroResetSeason(hero *entity.Hero, result herolock.LockResult, ctime, resetTime time.Time, seasonData *season.SeasonData) bool {
	if !hero.GetSeasonResetTime().Before(resetTime) {
		return false
	}

	if ctime.After(resetTime) {
		hero.ResetSeason(ctime, seasonData)
	} else {
		hero.ResetSeason(resetTime, seasonData)
	}

	if seasonData.AddMultiMonsterTimes > 0 || seasonData.PrevSeason.AddMultiMonsterTimes > 0 {
		multiLevelNpcTimes := hero.GetMultiLevelNpcTimes()
		result.Add(region.NewS2cUpdateMultiLevelNpcTimesMsg(multiLevelNpcTimes.StartTimeUnix32(), nil))
	}

	result.Add(domestic.NewS2cSeasonStartBroadcastMsg(seasonData.Season, timeutil.Marshal32(resetTime), true))

	result.Changed()
	result.Ok()

	return true
}

func ResetRandomEvent(hero *entity.Hero, ctime time.Time, datas iface.ConfigDatas) (arrX []int32, arrY []int32) {
	e := hero.RandomEvent()

	isNeedRefresh := false
	isSmallRefresh := false
	if e.CheckBigRefreshTime(ctime) {
		e.SetBigRefreshTime(ctime.Add(datas.MiscGenConfig().RandomEventBigRefreshDuration))
		e.SetSmallRefreshTime(ctime.Add(datas.MiscGenConfig().RandomEventSmallRefreshDuration))
		isNeedRefresh = true
	} else if e.CheckSmallRefreshTime(ctime) {
		e.SetSmallRefreshTime(ctime.Add(datas.MiscGenConfig().RandomEventSmallRefreshDuration))
		isNeedRefresh = true
		isSmallRefresh = true
	}
	if isNeedRefresh {
		cubes := e.GetAllEventsList()
		// 清理掉老数据
		e.ClearAllEvents()
		if isSmallRefresh {
			oldLen := len(cubes)
			if oldLen <= 0 {
				return
			}
			// 筛选出区块内和非区块内的
			satisfyList := datas.RandomEventPositionDictionary().SelectPositions4Block(cubes, hero.BaseX(), hero.BaseY())
			e.PutEvents(satisfyList)
			for _, c := range satisfyList {
				x, y := c.XYI32()
				arrX = append(arrX, x)
				arrY = append(arrY, y)
			}
			// 重新生成新的剩余数据位置保存
			cubes = datas.RandomEventPositionDictionary().GetRandomPositionsWithout(satisfyList, oldLen-len(satisfyList))
			e.PutEvents(cubes)
			for _, c := range cubes {
				x, y := c.XYI32()
				arrX = append(arrX, x)
				arrY = append(arrY, y)
			}
		} else {
			// 随机一匹数据存入
			cubes := datas.RandomEventPositionDictionary().GetRandomPositions(datas.MiscGenConfig().RandomEventNum)
			e.PutEvents(cubes)
			for _, c := range cubes {
				x, y := c.XYI32()
				arrX = append(arrX, x)
				arrY = append(arrY, y)
			}
			// 老家区块补足检测
			if cubes = datas.RandomEventPositionDictionary().CheckAndCatch4Block(cubes, datas.MiscGenConfig().RandomEventOwnMinNum, hero.BaseX(), hero.BaseY()); cubes != nil {
				e.PutEvents(cubes)
				for _, c := range cubes {
					x, y := c.XYI32()
					arrX = append(arrX, x)
					arrY = append(arrY, y)
				}
			}
		}
	}
	return
}

// 刷新随机事件
func heroResetRandomEvent(hero *entity.Hero, result herolock.LockResult, ctime time.Time, datas iface.ConfigDatas) {
	arrX, arrY := ResetRandomEvent(hero, ctime, datas)
	if len(arrX) > 0 {
		result.Add(random_event.NewS2cNewEventMsg(arrX, arrY))
	}
}

// 体力值刷新
func heroRecoverSp(hero *entity.Hero, result herolock.LockResult, ctime time.Time, spDuration time.Duration) {
	if ctime.Sub(hero.LastRecoverSpTime()) >= spDuration {
		hero.SetLastRecoverSpTime(ctime)
		if hero.IncreaseOneSp() {
			result.Add(domestic.NewS2cUpdateSpMsg(u64.Int32(hero.GetSp())))
		}
	}
}

var servicesRef = &atomic.Value{}

func GetService(dep iface.ServiceDep, datas iface.ConfigDatas, realmService iface.RealmService,
	guildService iface.GuildService, world iface.WorldService, mail iface.MailModule, buff iface.BuffService) *Services {

	r := servicesRef.Load()
	if r != nil {
		return r.(*Services)
	}

	s := &Services{
		datas:        datas,
		realmService: realmService,
		guildService: guildService,
		world:        world,
		mail:         mail,
		dep:          dep,
		buff:         buff,
	}
	servicesRef.Store(s)

	return s
}

type Services struct {
	datas        iface.ConfigDatas
	realmService iface.RealmService
	guildService iface.GuildService
	world        iface.WorldService
	mail         iface.MailModule
	dep          iface.ServiceDep
	buff         iface.BuffService
}

func TryUpdatePerSeconds(hc iface.HeroController, ctime time.Time, service *Services) {

	// 处理一下tick func
	hc.TickFunc()

	var toDelBuffs []*entity.BuffInfo

	var removeWorkerSeekHelpFunc func()
	var addInvasionMonsters []func()
	var addProsperityFunc func()
	var sendCopyDefenserExpiredMailFunc func()
	var checkRealmInfoFunc func()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hctx := NewContext(service.dep, operate_type.HeroUpdatePerSecond)

		updateHeroPerSeconds(hctx, hero, result, ctime, service.datas)

		toDelBuffs = updateBuffPerSecond(hero, result, ctime, service)
		UpdateAdvantageCount(hero, result, service.datas, ctime)

		heroId := hero.Id()

		// 如果驻守部队变化了，更新到场景中
		homeDefenser := hero.GetHomeDefenser()
		if homeDefenser != nil {
			if homeDefenser.ClearChanged() {
				if update, newAmount := hero.TryUpdateHomeDefenserFightAmount(); update {
					result.Add(domestic.NewS2cUpdateHeroFightAmountMsg(u64.Int32(newAmount)))

					UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_DEFENSER_FIGHTING)
				}
			}

			// 驻防自动补兵
			TryAutoFullSoldier(hero, result, service.datas, homeDefenser, ctime, service.datas.MiscGenConfig().AutoFullSoldoerDuration)
		}

		//if hero.Military().FreeSoldier() > 0 && hero.NeedRecoverTroopsSoldier() {
		//	// 有士兵且需要恢复士兵
		//	AutoRecoverCaptainSoldier(hero, result)
		//}

		if hero.TryExpireCopyDefenser(ctime) {
			result.Add(military.REMOVE_COPY_DEFENSER_S2C)

			sendCopyDefenserExpiredMailFunc = func() {
				if mailData := service.datas.MailHelp().CopyDefenserExpired; mailData != nil {
					proto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
					service.mail.SendProtoMail(heroId, proto, ctime)
				}
			}
		}

		TryRefreshWorkshop(hctx, hero, result, service.datas, ctime)

		UpdateAllTroopFightAmount(hero, result)

		removeWorkerPos, removeTechPos := hero.Domestic().TryRemoveSeekHelp(ctime)

		removeWorkerSeekHelpFunc = newRemoveHeroSeekHelpFunc(hero.IdBytes(), hero.GuildId(), removeWorkerPos, removeTechPos, service)

		if hero.BaseRegion() != 0 && hero.BaseLevel() > 0 && hero.Prosperity() > 0 {
			// 怪物攻城
			for _, t := range service.datas.RegionMultiLevelNpcTypeData().Array {
				if t.FightHate > 0 {
					npcType := t.Type
					if info := hero.GetNpcTypeInfo(t.Type); info != nil && info.GetHate() >= t.FightHate {
						if ctime.Before(info.GetRevengeTime()) {
							continue
						}

						if timeutil.IsZero(info.GetRevengeTime()) {
							info.SetRevengeTime(ctime.Add(t.FightDelay))
							result.Add(misc.NewS2cScreenShowWordsMsg(
								service.datas.TextHelp().RealmMonsterInvasion.New().
									WithMonster(t.Name).
									JsonString()))
							continue
						}

						newHate := u64.Sub(info.GetHate(), t.FightReduceHate)
						if newHate >= t.FightHate {
							newHate = u64.Sub(t.FightHate, 1)
						}

						info.SetHate(newHate)
						info.SetRevengeTime(time.Time{})
						result.Add(region.NewS2cUpdateMultiLevelNpcHateMsg(int32(t.Type), u64.Int32(newHate)))

						revengeLevel := info.GetRevengeLevel() + 1
						addInvasionMonsters = append(addInvasionMonsters, func() {
							service.realmService.GetBigMap().AddInvasionMonster(heroId, npcType, revengeLevel)
						})
					}
				}
			}

			if checkTime := hero.TaskList().GetNextCheckRealmInfoTime(); checkTime.Before(ctime) {
				hero.TaskList().SetNextCheckRealmInfoTime(ctime.Add(time.Minute))

				if !timeutil.IsZero(checkTime) {
					hasHomeNpcTask := false
					hero.TaskList().WalkAllTask(func(task entity.Task) (endedWalk bool) {
						if !task.Progress().IsCompleted() {
							if task.Progress().Target().KillHomeNpcData != nil {
								hasHomeNpcTask = true
								return true
							}
						}
						return false
					})

					checkMonsterTask := hero.TaskList().HasTaskMonster()

					checkTroopOutside := false
					for _, t := range hero.Troops() {
						if t.IsOutside() {
							checkTroopOutside = true
							break
						}
					}

					if hasHomeNpcTask || checkMonsterTask || checkTroopOutside {
						checkRealmInfoFunc = func() {
							service.realmService.GetBigMap().UpdateHeroRealmInfo(heroId, hasHomeNpcTask, checkMonsterTask, checkTroopOutside)
						}
					}
				}
			}
		}

		// 刷新政务
		heroZhengWu := hero.ZhengWu()
		nextRefreshTime := heroZhengWu.NextRefreshTime()
		if ctime.After(nextRefreshTime) {
			// 刷新
			miscData := service.datas.ZhengWuMiscData()
			nextRefreshTime := miscData.NextAutoRefreshTime(ctime)

			heroZhengWu.SetNextRefreshTime(nextRefreshTime)

			// 刷新
			randomCount := miscData.RandomCount
			//if heroZhengWu.Doing() != nil {
			//	randomCount -= 1
			//}

			completedFirstZhengWu := hero.Bools().Get(shared_proto.HeroBoolType_BOOL_COMPELTED_FIRST_ZHENGWU)
			heroZhengWu.SetToDoList(service.datas.ZhengWuRandomData().Random(randomCount, completedFirstZhengWu))

			result.Add(zhengwu.NewS2cRefreshMsg(must.Marshal(heroZhengWu.EncodeClient())))
		}

		// 刷新黑市商品
		TryRefreshBlackMarketGoods(hero, result, service.datas, ctime, true)

		//// 解锁外城
		//if hero.Domestic().OuterCities().TryAutoUnlock() {
		//	toAddProsperity := TryUnlockOutsideCity(hero, result, service.datas, ctime)
		//	addProsperityFunc = service.realmService.AddProsperityFunc(hero.Id(), hero.BaseRegion(), toAddProsperity, "定时解锁分城")
		//}

	})

	service.buff.Cancel(hc.Id(), toDelBuffs) // 删掉过期buff效果

	if removeWorkerSeekHelpFunc != nil {
		removeWorkerSeekHelpFunc()
	}

	if len(addInvasionMonsters) > 0 {
		for _, f := range addInvasionMonsters {
			f()
		}
	}

	if addProsperityFunc != nil {
		addProsperityFunc()
	}

	if sendCopyDefenserExpiredMailFunc != nil {
		sendCopyDefenserExpiredMailFunc()
	}

	if checkRealmInfoFunc != nil {
		checkRealmInfoFunc()
	}
}

func TryRefreshBlackMarketGoods(hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas, ctime time.Time, isRefresh bool) bool {
	nextRefreshTime := hero.Shop().GetBlackMarketNextRefreshTime()
	if ctime.After(nextRefreshTime) {
		nextRefreshTime := datas.ShopMiscData().NextAutoRefreshBlackMarketTime(ctime)
		hero.Shop().SetBlackMarketNextRefreshTime(nextRefreshTime)

		bmd := datas.GetBlackMarketData(shopdata.BLACK_YUNYOU)
		if bmd != nil {
			goods, discount := bmd.Random(hero.Level())
			hero.UpdateBlackMarketGoods(goods, discount)

			SendBlackMarketGoodsMsg(hero, result, isRefresh)
		}
		return true
	}

	return false
}

func SendBlackMarketGoodsMsg(hero *entity.Hero, result herolock.LockResult, isRefresh bool) {

	object := &shop.S2CPushBlackMarketGoodsProto{}
	object.Refrash = isRefresh
	object.NextRefreshTime = timeutil.Marshal32(hero.Shop().GetBlackMarketNextRefreshTime())
	hero.RangeBlackMarketGoods(func(item *entity.BlackMarketGoodsItem) bool {
		object.GoodsId = append(object.GoodsId, u64.Int32(item.Goods().Id))
		object.Discount = append(object.Discount, u64.Int32(item.Discount()))
		object.Buy = append(object.Buy, item.IsBuy())

		return true
	})

	result.Add(shop.NewS2cPushBlackMarketGoodsProtoMsg(object))
}

func TryUnlockOutsideCity(hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas,
	ctime time.Time) uint64 {
	var unlockDatas []*domestic_data.OuterCityData
	var unlockTypes []uint64
	for _, data := range datas.GetOuterCityDataArray() {
		outerCity := hero.Domestic().OuterCities().OuterCity(data)
		if outerCity != nil {
			continue
		}

		var funcType uint64
		switch data.Diraction {
		case shared_proto.PosDiraction_EAST:
			funcType = constants.FunctionType_TYPE_FEN_CHENG
		case shared_proto.PosDiraction_SOUTH:
			funcType = constants.FunctionType_TYPE_FEN_CHENG_2
		case shared_proto.PosDiraction_WEST:
			funcType = constants.FunctionType_TYPE_FEN_CHENG_3
		case shared_proto.PosDiraction_NORTH:
			funcType = constants.FunctionType_TYPE_FEN_CHENG_4
		default:
			continue
		}

		if !hero.Function().IsFunctionOpened(funcType) {
			continue
		}

		unlockDatas = append(unlockDatas, data)
		unlockTypes = append(unlockTypes, data.RecommandType)
	}

	return UnlockOutsideCity(hero, result, unlockDatas, unlockTypes, datas, ctime)
}

func UnlockOutsideCity(hero *entity.Hero, result herolock.LockResult, unlockDatas []*domestic_data.OuterCityData, unlockTypes []uint64,
	datas iface.ConfigDatas, ctime time.Time) uint64 {

	n := imath.Min(len(unlockTypes), len(unlockDatas))
	if n <= 0 {
		return 0
	}

	heroDomestic := hero.Domestic()
	outerCities := heroDomestic.OuterCities()

	var toAddProsperity uint64
	var unlockCount uint64
	for i := 0; i < n; i++ {
		data := unlockDatas[i]
		t := unlockTypes[i]

		outerCity := outerCities.OuterCity(data)
		if outerCity != nil {
			continue
		}

		outerCity = entity.NewOuterCity(data, t)
		outerCities.Unlock(outerCity)

		// 发送解锁成功
		result.Add(domestic.NewS2cUnlockOuterCityMsg(must.Marshal(outerCity.Encode())))

		unlockCount += outerCity.Count()

		outerCity.WalkLayouts(func(layout *domestic_data.OuterCityLayoutData, building *domestic_data.OuterCityBuildingData) {
			toAddProsperity += building.BuildingData.Prosperity

			UpdateBuildingEffect(hero, result, datas, ctime, building.BuildingData.Effect)
			//m.updateBuildingLevel(hero, result, nil, layout.BuildingData, ctime)
		})

		if toAddProsperity > 0 {
			hero.AddProsperityCapcity(toAddProsperity)
		}

		result.Changed()
	}
	outerCities.UpdateUnlockBit()

	// 更新任务进度
	UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL)
	if unlockCount > 0 {
		IncreTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL, unlockCount)
	}

	return toAddProsperity
}

func TryRefreshWorkshop(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas, ctime time.Time) bool {
	if nextTime := hero.Domestic().GetNextRefreshWorkshopTime(); nextTime.Before(ctime) {
		newNextTime := datas.MiscConfig().GetWorkshopNextRefreshTime(ctime)
		RefreshWorkshopAnyway(hctx, hero, result, datas, newNextTime, operate_type.RefreshAuto, 0)
		return true
	}

	return false
}

func RefreshWorkshopAnyway(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas, newNextTime time.Time, refreshType, costId uint64) {
	// 到了要刷新的时间
	building := hero.Domestic().GetBuilding(shared_proto.BuildingType_TIE_JIANG_PU)

	var newEquipment []*goods.EquipmentData
	if building != nil {
		workshopData := datas.WorkshopLevelData().Must(building.Level)
		if workshopData != nil {
			newEquipment = workshopData.RandomEquipment()
		}
	}

	hero.Domestic().RefreshWorkshop(newNextTime, newEquipment)

	SendWorkshopMsg(hero, result, datas)

	// tlog
	hctx.Tlog().TlogRefreshFlow(hero, uint64(shared_proto.BuildingType_TIE_JIANG_PU), refreshType, costId)

}

func SendWorkshopMsg(hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas) {
	// 发更新消息
	var equipmentIds []int32
	var durations []int32
	for _, e := range hero.Domestic().GetWorkshopEquipment() {
		equipmentIds = append(equipmentIds, u64.Int32(e.Id))

		d := time.Hour
		if durationData := datas.GetWorkshopDuration(e.Id); durationData != nil {
			d = durationData.Duration
		}
		durations = append(durations, timeutil.DurationMarshal32(d))
	}

	result.Add(domestic.NewS2cListWorkshopEquipmentMsg(
		timeutil.Marshal32(hero.Domestic().GetNextRefreshWorkshopTime()),
		equipmentIds, durations,
		u64.Int32(hero.Domestic().GetWorkshopIndex()),
		timeutil.Marshal32(hero.Domestic().GetWorkshopCollectTime()),
		u64.Int32(hero.Domestic().GetWorkshopRefreshTimes()),
	))
}

func updateHeroPerSeconds(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, ctime time.Time, datas iface.ConfigDatas) {

	hero.UpdateOnlineTime(ctime)

	// 删除过期物品
	removeIds := hero.Depot().RemoveExpiredGoods(ctime)
	if len(removeIds) > 0 {
		result.Add(depot.NewS2cGoodsExpireTimeRemoveMsg(u64.Int32Array(removeIds)))
		result.Changed()
	}

	// 检查预约消耗时间，返回预约消耗
	tickAddBackReservation(hctx, hero, result, ctime)

	// 仓库衰减
	resDecay(hctx, hero, result, datas.MiscConfig().ExtraResDecayCoef, datas.MiscConfig().ExtraResDecayDuration, ctime)

	// 加税收
	TryUpdateTax(hero, result, ctime, datas.MiscGenConfig().TaxDuration, datas.GetBuffEffectData)
}

func resDecay(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, coef *data.Amount, decayDuration time.Duration, ctime time.Time) {
	if ctime.Before(hero.NextResDecayTime()) {
		return
	}
	if timeutil.IsZero(hero.NextResDecayTime()) {
		hero.SetNextResDecayTime(ctime)
	}
	defer hero.SetNextResDecayTime(ctime.Add(decayDuration))

	multi := u64.Max(1, u64.DivideTimes(uint64(ctime.Sub(hero.NextResDecayTime()).Seconds()), uint64(decayDuration.Seconds())))

	goldProtected, foodProtected, woodProtected, stoneProtected := hero.Domestic().StorageProtected(hero.GetUnsafeResource())
	var toReduceGold, toReduceStone, toReduceFood, toReduceWood uint64
	if !goldProtected {
		amount := hero.GetUnsafeResource().Gold()
		capacity := hero.Domestic().GoldCapcity()
		toReduceGold = u64.Min(coef.CalculateByPercent(amount)*multi, u64.Sub(amount, capacity))
	}
	if !stoneProtected {
		amount := hero.GetUnsafeResource().Stone()
		capacity := hero.Domestic().StoneCapcity()
		toReduceStone = u64.Min(coef.CalculateByPercent(amount)*multi, u64.Sub(amount, capacity))
	}
	if !foodProtected {
		amount := hero.GetUnsafeResource().Food()
		capacity := hero.Domestic().FoodCapcity()
		toReduceFood = u64.Min(coef.CalculateByPercent(amount)*multi, u64.Sub(amount, capacity))
	}
	if !woodProtected {
		amount := hero.GetUnsafeResource().Wood()
		capacity := hero.Domestic().WoodCapcity()
		toReduceWood = u64.Min(coef.CalculateByPercent(amount)*multi, u64.Sub(amount, capacity))
	}

	changed, msgFunc := ReduceUnsafeResource(hctx, hero, result, toReduceGold, toReduceFood, toReduceWood, toReduceStone)

	if changed && msgFunc != nil {
		result.Add(msgFunc())
	}
}

func TryUpdatePerMinute(hc iface.HeroController, ctime time.Time, service *Services) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		// 请求过期瞭望记录
		hero.ClearExpiredInvestigation(ctime, false)
		// 随机事件刷新
		heroResetRandomEvent(hero, result, ctime, service.datas)
		// 体力值恢复
		heroRecoverSp(hero, result, ctime, service.datas.MiscGenConfig().SpDuration)
	})
}

func AddGongXun(hero *entity.Hero, result herolock.LockResult, captain *entity.Captain, toAdd uint64) {
	newGongXun := captain.AddGongXun(toAdd)
	result.Add(military.NewS2cAddGongxunMsg(u64.Int32(captain.Id()), u64.Int32(newGongXun)))
	result.Changed()
}

func AddCaptainAbilityExp(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, captain *entity.Captain, toAddExp uint64, ctime time.Time) (upgrade bool) {
	if toAddExp <= 0 {
		return false
	}

	oldLevel := captain.Level()
	oldAbility := captain.Ability()
	oldRace := captain.Race().Id
	oldOfficial := captain.Official().Id
	oldUnlockSpellCount := captain.AbilityData().UnlockSpellCount

	captain.AddAbilityExp(toAddExp)

	// 判断是否升级
	if !captain.TryUpgradeAbility() {
		result.Add(military.NewS2cUpdateAbilityExpMsg(u64.Int32(captain.Id()), u64.Int32(captain.AbilityExp())))
		return
	}

	result.Add(military.NewS2cCaptainRefinedUpgradeMsg(u64.Int32(captain.Id()), u64.Int32(captain.AbilityExp()),
		u64.Int32(captain.Ability()), nil, int32(captain.Quality())))

	captain.CalculateProperties()
	result.Add(captain.NewUpdateCaptainStatMsg())
	UpdateTroopFightAmount(hero, captain.GetTroop(), result)

	// 更新任务进度
	UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_QUALITY_COUNT)
	UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_ABILITY_EXP)

	// 系统广播
	if d := hctx.BroadcastHelp().CaptainAbilityExp; d != nil {
		// 恭喜[color=#FFD633]{{self}}[/color]首次将麾下武将强化到了突破+{{quality}}
		hctx.AddBroadcast(d, hero, result, 0, uint64(captain.Quality()), func() *i18n.Fields {
			text := d.NewTextFields()
			text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
			text.WithFields(data.KeyQuality, uint64(captain.Quality()))
			return text
		})
	}

	// tlog
	hctx.Tlog().TlogPlayerCultivateFlow(hero, captain.Id(), operate_type.CaptainOperTypeRefine, oldLevel, captain.Level(), oldAbility, captain.Ability(), u64.FromInt(oldRace), u64.FromInt(captain.Race().Id), oldOfficial, captain.Official().Id, hctx.OperId())

	if oldUnlockSpellCount != captain.AbilityData().UnlockSpellCount {
		// 武将内政技能
		if captain.StarData().HasBuildingEffectSpell() {
			UpdateBuildingEffect(hero, result, hctx.datas, ctime,
				captain.StarData().GetBuildingEffectSpell(captain.AbilityData().UnlockSpellCount)...)
		}
	}

	return true
}

func AddVipExp(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, toAdd uint64, ctime time.Time) (upgrade bool) {
	if toAdd <= 0 {
		return
	}

	v := hero.Vip()
	oldData := hctx.Datas().GetVipLevelData(v.Level())
	if oldData == nil || oldData.NextLevelData == nil {
		return
	}

	v.AddExp(toAdd)

	newData := oldData
	for d := oldData; v.Exp() >= d.UpgradeExp; d = hctx.Datas().GetVipLevelData(v.Level()) {
		nextLevel := d.NextLevelData
		if nextLevel == nil || d.Level >= nextLevel.Level {
			break
		}
		newData = nextLevel

		v.SetExp(u64.Sub(v.Exp(), d.UpgradeExp))
		v.SetLevel(nextLevel.Level)
		v.AddDailyPrizeCanCollectLevel(v.Level())
		upgrade = true
	}

	if newData.NextLevelData == nil {
		v.SetExp(0)
	}

	if upgrade {
		result.Changed()
		result.Add(vipMsg.NewS2cVipLevelUpgradeNoticeMsg(u64.Int32(v.Level()), u64.Int32(v.Exp())))
		onVipLevelUp(hero, result, oldData, newData, hctx.Datas(), ctime)
	} else {
		result.Add(vipMsg.NewS2cVipAddExpNoticeMsg(u64.Int32(v.Exp())))
	}

	return
}

// 添加充值金额，比如购买元宝，购买每日礼包等等
// 对应的人民币需要增加活动相关值，比如累积充值金额，vip经验等等
func AddRechargeAmount(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, toAdd uint64, ctime time.Time) {

	// 累积充值金额
	hero.Misc().AddChargeAmount(toAdd, ctime)
	result.Add(misc.NewS2cUpdateChargeAmountMsg(u64.Int32(hero.Misc().ChargeAmount())))

	// 增加Vip经验（1元宝 = 1点Vip经验）
	// （如果某些花钱项目不需要加Vip经验，则这个方法新增一个bool值表示是否需要添加Vip经验）
	AddVipExp(hctx, hero, result, toAdd, ctime)
}

func AddExp(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, toAdd uint64, ctime time.Time) (upgrade bool) {
	if toAdd <= 0 {
		return
	}

	if hero.IsMaxLevel() {
		return
	}

	oldLevelData := hero.LevelData()
	upgrade = hero.AddExp(toAdd)
	result.Changed()
	if upgrade {
		newLevelData := hero.LevelData()

		// 升级了
		result.Add(domestic.NewS2cHeroUpgradeLevelMsg(u64.Int32(hero.Exp()), u64.Int32(hero.Level())))

		onHeroLevelUp(hero, result, oldLevelData, hero.LevelData(), ctime)

		// 系统广播
		if d := hctx.BroadcastHelp().HeroLevel; d != nil {
			hctx.AddBroadcast(d, hero, result, 0, hero.Level(), func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, hero.Level())
				return text
			})
		}

		// tlog
		hctx.Tlog().TlogKingExpFlow(hero, oldLevelData.Level, hero.Level(), hctx.OperId())

		if giftData := hctx.datas.EventLimitGiftConfig().GetHeroLevelGift(hero.Level()); giftData != nil {
			ActivateEventLimitGift(hero, result, giftData, ctime)
		}

		// eventlog
		for i := oldLevelData.Level + 1; i <= newLevelData.Level; i++ {
			gamelogs.UpgradeHeroLevelLog(constants.PID, hero.Sid(), hero.Id(), i)
		}

	} else {
		result.Add(domestic.NewS2cHeroUpdateExpMsg(u64.Int32(hero.Exp())))
	}

	return
}

func onHeroLevelUp(hero *entity.Hero, result herolock.LockResult, oldLevelData, newLevelData *herodata.HeroLevelData, ctime time.Time) {
	UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HERO_LEVEL)

	OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_HERO_LEVEL, newLevelData.Level)

	//originTraining := hero.Military().OriginTraining()
	//newTraining := hero.Military().GetTraining(newLevelData)
	//if len(originTraining) < len(newTraining) {
	//	index := make([]int32, 0, len(newTraining)-len(originTraining))
	//	output := make([]int32, 0, len(newTraining)-len(originTraining))
	//	capcity := make([]int32, 0, len(newTraining)-len(originTraining))
	//	exps := make([]int32, 0, len(newTraining)-len(originTraining))
	//	for i := len(originTraining); i < len(newTraining); i++ {
	//		t := newTraining[i]
	//
	//		index = append(index, int32(i))
	//		output = append(output, u64.Int32(t.Output()))
	//		capcity = append(capcity, u64.Int32(t.Capcity()))
	//		exps = append(exps, 0) // 新开启的，经验肯定是0
	//	}
	//
	//	// 更新修炼位产出
	//	result.Add(military.NewS2cUpdateTrainingOutputMsg(index, output, capcity, exps))
	//
	//	UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_TRAINING_LEVEL_COUNT)
	//
	//	result.Changed()
	//}

	//if oldLevelData.Sub.StrategyLimit != newLevelData.Sub.StrategyLimit {
	//	// 更新君主策略开始恢复时间
	//	if hero.GetStrategyCount(ctime) >= oldLevelData.Sub.StrategyLimit {
	//		hero.SetStrategyCount(oldLevelData.Sub.StrategyLimit, ctime)
	//		result.Add(domestic.NewS2cUpdateStrategyRestoreStartTimeMsg(timeutil.Marshal32(hero.GetStrategyRestoreStartTime())))
	//	}
	//}

	if oldLevelData.Sub.AddSoldierCapacity != newLevelData.Sub.AddSoldierCapacity {
		for _, captain := range hero.Military().Captains() {
			captain.CalculateProperties()
			result.Add(captain.NewUpdateCaptainStatMsg())
		}
		UpdateAllTroopFightAmount(hero, result)
	}

	oldTroopCount := len(hero.Troops())
	hero.Military().OnHeroLevelChanged(hero.Id(), newLevelData)
	newTroopCount := len(hero.Troops())

	for index := oldTroopCount; index < newTroopCount; index++ {
		troop := hero.GetTroopByIndex(u64.FromInt(index))
		if troop == nil {
			logrus.Errorf("升级以后，新增的部队竟然没有发给客户端: %d", index)
			continue
		}

		result.Add(military.NewS2cNewTroopsMsg(must.Marshal(troop.EncodeClient())))
	}

	CheckFuncsOpened(hero, result)
}

func newRemoveHeroSeekHelpFunc(heroIdBytes []byte, guildId int64, removeWorkerPos, removeTechPos []int, service *Services) func() {

	if guildId == 0 {
		return nil
	}

	if len(removeWorkerPos) <= 0 && len(removeTechPos) <= 0 {
		return nil
	}

	return func() {
		var allMemberIds []int64
		var broadcastMsgs []pbutil.Buffer
		service.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
			if g == nil {
				return
			}

			for _, pos := range removeWorkerPos {
				key := sharedguilddata.NewHeroWorkerSeekHelpKey(heroIdBytes, constants.SeekTypeWorker, int32(pos))
				if g.GetSeekHelp(key) != nil {
					g.RemoveSeekHelp(key)
					broadcastMsgs = append(broadcastMsgs, guild.NewS2cRemoveGuildSeekHelpMsg(key))
				}
			}

			for _, pos := range removeTechPos {
				key := sharedguilddata.NewHeroWorkerSeekHelpKey(heroIdBytes, constants.SeekTypeTech, int32(pos))
				if g.GetSeekHelp(key) != nil {
					g.RemoveSeekHelp(key)
					broadcastMsgs = append(broadcastMsgs, guild.NewS2cRemoveGuildSeekHelpMsg(key))
				}
			}

			if len(broadcastMsgs) > 0 {
				allMemberIds = g.AllUserMemberIds()
			}
		})

		service.world.MultiSendMsgs(allMemberIds, broadcastMsgs)
	}
}

//func newHeroDefenseData(hero *entity.Hero, realmService iface.RealmService) (r iface.Realm, version uint64, protoBytes []byte) {
//
//	baseRegion := hero.BaseRegion()
//	if baseRegion == 0 {
//		return
//	}
//
//	r = realmService.GetBigMap()
//	if r == nil {
//		logrus.WithField("region", baseRegion).Error("更新玩家防守队伍，找不到地图")
//		return
//	}
//
//	version, protoBytes = hero.NewHomeDefenserBytes()
//
//	return
//}

func SendResourcePointUpdateMsg(hero *entity.Hero, result herolock.LockResult, point *entity.ResourcePoint) {

	effect := point.Building().GetResPointEffect()
	if effect == nil {
		logrus.Error("heromodule.UpdateResourcePoint effect == nil")
		return
	}

	outputPerHour, outputCapcity := hero.CalculateResourcePointOutput(effect)
	result.Add(domestic.NewS2cUpdateResourceBuildingMsg(u64.Int32(point.LayoutData().Id), u64.Int32(point.OutputAmount()),
		u64.Int32(outputCapcity), u64.Int32(outputPerHour), hero.IsConflictResourcePoint(point.LayoutData()),
		hero.BaseLevel() < point.LayoutData().RequireBaseLevel))
}

func SendResourcePointUpdateMsgByResType(hero *entity.Hero, result herolock.LockResult, resType shared_proto.ResType) {

	extraOutputPerHour := hero.Domestic().GetExtraOutput(resType)
	extraCapcity := hero.Domestic().GetExtraOutputCapcity(resType)

	proto := &domestic.S2CUpdateMultiResourceBuildingProto{}
	hero.Domestic().WalkResourcePoint(func(pos uint64, point *entity.ResourcePoint) bool {

		effect := point.Building().GetResPointEffect()
		if effect != nil && effect.OutputType == resType {

			// 前面已经计算过了，这里不再计算
			// data.CalculateCurrentOutput(hero, ctime)

			proto.Id = append(proto.Id, u64.Int32(point.LayoutData().Id))
			proto.Amount = append(proto.Amount, u64.Int32(point.OutputAmount()))

			outputPerHour, outputCapcity := entity.CalculateResourcePointOutputInfo(effect, extraOutputPerHour, extraCapcity)
			proto.Capcity = append(proto.Capcity, u64.Int32(outputCapcity))
			proto.Output = append(proto.Output, u64.Int32(outputPerHour))
			proto.Conflict = append(proto.Conflict, hero.IsConflictResourcePoint(point.LayoutData()))
			proto.BaseLevelLock = append(proto.BaseLevelLock, hero.BaseLevel() < point.LayoutData().RequireBaseLevel)
		}

		return false
	})

	if len(proto.Id) > 0 {
		result.Add(domestic.NewS2cUpdateMultiResourceBuildingProtoMsg(proto))
	}

}

// 资源点冲突更新
func SendResourcePointUpdateMsgByLayoutIds(hero *entity.Hero, result herolock.LockResult, layoutIds []uint64) {
	if len(layoutIds) <= 0 {
		return
	}

	proto := &domestic.S2CUpdateMultiResourceBuildingProto{}

	for _, layoutId := range layoutIds {
		point := hero.Domestic().GetLayoutRes(layoutId)
		if point == nil {
			continue
		}

		effect := point.Building().GetResPointEffect()
		if effect == nil {
			continue
		}

		// 前面已经计算过了，这里不再计算
		// data.CalculateCurrentOutput(hero, ctime)

		proto.Id = append(proto.Id, u64.Int32(point.LayoutData().Id))
		proto.Amount = append(proto.Amount, u64.Int32(point.OutputAmount()))

		extraOutputPerHour := hero.Domestic().GetExtraOutput(effect.OutputType)
		extraCapcity := hero.Domestic().GetExtraOutputCapcity(effect.OutputType)
		outputPerHour, outputCapcity := entity.CalculateResourcePointOutputInfo(effect, extraOutputPerHour, extraCapcity)
		proto.Capcity = append(proto.Capcity, u64.Int32(outputCapcity))
		proto.Output = append(proto.Output, u64.Int32(outputPerHour))
		proto.Conflict = append(proto.Conflict, hero.IsConflictResourcePoint(point.LayoutData()))
		proto.BaseLevelLock = append(proto.BaseLevelLock, hero.BaseLevel() < point.LayoutData().RequireBaseLevel)
	}

	if len(proto.Id) > 0 {
		result.Add(domestic.NewS2cUpdateMultiResourceBuildingProtoMsg(proto))
	}

}

//func AutoRecoverCaptainSoldier(hero *entity.Hero, result herolock.LockResult) {
//	heroMilitary := hero.Military()
//
//	if heroMilitary.FreeSoldier() <= 0 {
//		return
//	}
//
//	troops := make([]*entity.Troop, 0, len(hero.Troops()))
//	for _, troop := range hero.Troops() {
//		if troop.IsOutside() {
//			continue
//		}
//
//		troops = append(troops, troop)
//	}
//
//	if len(troops) <= 0 {
//		return
//	}
//
//	sort.Sort(entity.TroopFullFightAmountSlice(troops))
//
//	var captainIds []int32
//	var captainSoldiers []int32
//
//	for _, troop := range troops {
//		for _, captain := range troop.Captains() {
//			if captain == nil || captain.Soldier() >= captain.SoldierCapcity() {
//				continue
//			}
//
//			// 恢复
//			recover := u64.Min(captain.SoldierCapcity()-captain.Soldier(), heroMilitary.FreeSoldier())
//			captain.AddSoldier(recover)
//			result.Add(captain.NewUpdateCaptainStatMsg())
//
//			heroMilitary.ReduceFreeSoldier(recover)
//
//			captainIds = append(captainIds, u64.Int32(captain.Id()))
//			captainSoldiers = append(captainSoldiers, u64.Int32(captain.Soldier()))
//
//			if heroMilitary.FreeSoldier() <= 0 {
//				break
//			}
//		}
//
//		if heroMilitary.FreeSoldier() <= 0 {
//			break
//		}
//	}
//
//	// 更新君主战斗力
//	if update, newAmount := hero.TryUpdateHomeDefenserFightAmount(); update {
//		result.Add(domestic.NewS2cUpdateHeroFightAmountMsg(u64.Int32(newAmount)))
//	}
//
//	if len(captainIds) <= 0 {
//		return
//	}
//
//	result.Add(military.NewS2cAutoRecoverSoldierMsg(u64.Int32(heroMilitary.FreeSoldier()), captainIds, captainSoldiers))
//
//	result.Changed()
//	result.Ok()
//}

func CheckFuncsOpened(hero *entity.Hero, result herolock.LockResult) {
	heroFunction := hero.Function()
	var funcs []*function.FunctionOpenData
	for _, data := range heroFunction.FunctionOpenDataArray {
		if GetIsFuncOpened(hero, data) {
			heroFunction.OpenFunction(data.FunctionType)
			funcs = append(funcs, data)

			// 功能解锁
			OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_UNLOCK_FUNC, uint64(data.FunctionType))

			// 分城解锁
			switch data.FunctionType {
			case constants.FunctionType_TYPE_FEN_CHENG,
				constants.FunctionType_TYPE_FEN_CHENG_2,
				constants.FunctionType_TYPE_FEN_CHENG_3,
				constants.FunctionType_TYPE_FEN_CHENG_4:
				hero.Domestic().OuterCities().SetAutoUnlock()
			}

			result.Changed()
		}
	}

	// 发消息
	switch len(funcs) {
	case 0:
	case 1:
		data := funcs[0]
		result.Add(data.OpenMsg)
		result.Changed()
	default:
		// 一次性解锁多个
		toSendProto := &misc.S2COpenMultiFunctionProto{}
		for _, data := range funcs {
			toSendProto.FunctionType = append(toSendProto.FunctionType, int32(data.FunctionType))
		}
		result.Add(misc.NewS2cOpenMultiFunctionProtoMsg(toSendProto))
	}
}

// 检查功能是否有开启
func GetIsFuncOpened(hero *entity.Hero, data *function.FunctionOpenData) bool {
	if hero.Function().IsFunctionOpened(data.FunctionType) {
		return false
	}

	if data.GuanFuLevel != nil {
		if guanFu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU); guanFu == nil || data.GuanFuLevel.Level > guanFu.Level {
			return false
		}
	}

	if data.Building != nil {
		if building := hero.Domestic().GetBuilding(data.Building.Type); building == nil || data.Building.Level > building.Level {
			return false
		}
	}

	if data.HeroLevel != nil && data.HeroLevel.Level > hero.Level() {
		return false
	}

	if data.MainTask != nil && data.MainTask.Sequence > hero.TaskList().CompletedMainTaskSequence() {
		return false
	}

	if data.BaYeStage != nil && data.BaYeStage.Stage > hero.TaskList().GetCompletedBaYeStage() {
		return false
	}

	if data.TowerFloor != nil && data.TowerFloor.Floor > hero.Tower().HistoryMaxFloor() {
		return false
	}

	if data.Dungeon != nil && !hero.Dungeon().IsPass(data.Dungeon) {
		return false
	}

	return true
}

func UpdateAllTroopFightAmount(hero *entity.Hero, result herolock.LockResult) {
	var troopIndex, startAmount, endAmount []int32
	var fightAmountMax uint64
	for _, t := range hero.Troops() {
		oldAmount, newAmount := t.UpdateFightAmountIfChanged()
		if oldAmount < newAmount {
			troopIndex = append(troopIndex, u64.Int32(t.Sequence())+1)
			startAmount = append(startAmount, u64.Int32(oldAmount))
			endAmount = append(endAmount, u64.Int32(newAmount))
			if fightAmountMax < newAmount {
				fightAmountMax = newAmount
			}
		}
	}
	if len(troopIndex) > 0 {
		result.Add(military.NewS2cUpdateTroopFightAmountMsg(troopIndex, startAmount, endAmount))
	}
	if fightAmountMax > 0 && hero.HistoryAmount().Amount(server_proto.HistoryAmountType_MaxTroopFightAmount) < fightAmountMax {
		hero.HistoryAmount().Set(server_proto.HistoryAmountType_MaxTroopFightAmount, fightAmountMax)
		UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_TEAM_POWER)
	}
}

func UpdateTroopFightAmount(hero *entity.Hero, troop *entity.Troop, result herolock.LockResult) {
	if troop == nil {
		return
	}

	oldAmount, newAmount := troop.UpdateFightAmountIfChanged()
	if oldAmount < newAmount {
		result.Add(military.NewS2cUpdateTroopFightAmountMsg(
			[]int32{u64.Int32(troop.Sequence()) + 1},
			[]int32{u64.Int32(oldAmount)},
			[]int32{u64.Int32(newAmount)}))
		if hero.HistoryAmount().Amount(server_proto.HistoryAmountType_MaxTroopFightAmount) < newAmount {
			hero.HistoryAmount().Set(server_proto.HistoryAmountType_MaxTroopFightAmount, newAmount)
			UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_TEAM_POWER)
		}
	}
}

func UpdateBuildingEffect(hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas, ctime time.Time, effects ...*sub.BuildingEffectData) {

	if len(effects) <= 0 {
		return
	}

	isEffect := func(effects []*sub.BuildingEffectData, f func(effect *sub.BuildingEffectData) bool) bool {
		for _, effect := range effects {
			if effect != nil && f(effect) {
				return true
			}
		}
		return false
	}

	// 仓库
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.IsCapcityEffect()
	}) {
		hero.Calculate(shared_proto.BuildingType_CANG_KU)

		d := hero.Domestic()
		result.Add(domestic.NewS2cResourceCapcityUpdateMsg(
			u64.Int32(d.GoldCapcity()), u64.Int32(d.FoodCapcity()), u64.Int32(d.WoodCapcity()), u64.Int32(d.StoneCapcity()),
			u64.Int32(d.ProtectedCapcity())))
	}

	// 兵营
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.IsSoldierEffect()
	}) {
		//m := hero.Military()
		//result.Add(military.NewS2cUpdateSoldierCapcityMsg(u64.Int32(m.SoldierCapcity()), u64.Int32(m.WoundedCapcity()),
		//	u64.Int32(m.NewSoldierCapcity()), u64.Int32(m.NewSoldierOutput()), u64.Int32(m.NewSoldierRecruitCount())))

		// 更新空闲士兵
		soldier := hero.Military().FreeSoldier(ctime)
		hero.Calculate(shared_proto.BuildingType_JUN_YING)
		hero.Military().UpdateFreeSoldierOutputCapcity(ctime)
		hero.Military().SetFreeSoldier(soldier, ctime)

		result.Add(hero.Military().NewUpdateFreeSoldierMsg())
	}

	// 官府
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.BuildingWorkerCdr > 0
	}) {
		hero.Calculate(shared_proto.BuildingType_GUAN_FU)

		result.Add(domestic.NewS2cUpdateBuildingWorkerCoefMsg(i32.MultiF64(1000, hero.Domestic().GetBuildingWorkerCdr())))
	}

	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.BuildingWorkerFatigueDuration > 0
	}) {
		result.Add(domestic.NewS2cUpdateBuildingWorkerFatigueDurationMsg(timeutil.DurationMarshal32(hero.Domestic().GetBuildingWorkerFatigueDuration())))
	}

	// 书院
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.TechWorkerCdr > 0
	}) {
		hero.Calculate(shared_proto.BuildingType_SHU_YUAN)

		result.Add(domestic.NewS2cUpdateTechWorkerCoefMsg(i32.MultiF64(1000, hero.Domestic().GetTechWorkerCdr())))
	}

	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.TechWorkerFatigueDuration > 0
	}) {
		result.Add(domestic.NewS2cUpdateTechWorkerFatigueDurationMsg(timeutil.DurationMarshal32(hero.Domestic().GetTechWorkerFatigueDuration())))
	}

	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.IsTrainCoefEffect()
	}) {
		oldCoef := hero.BuildingEffect().GetTrainCoef()
		hero.Calculate(shared_proto.BuildingType_SI_TU_FU)
		newCoef := hero.BuildingEffect().GetTrainCoef()

		updateCaptainTrain(hero, result, oldCoef, newCoef, datas, ctime)
	}

	// 资源点
	outputTypeMap := make(map[shared_proto.ResType]struct{})
	for _, effect := range effects {
		if effect != nil && effect.IsResPointEffect() {
			outputTypeMap[effect.OutputType] = struct{}{}
		}
	}
	for t := range outputTypeMap {
		if by, ok := domestic_data.GetResBuildingType(t); ok {
			hero.Calculate(by)

			// 更新消息
			SendResourcePointUpdateMsgByResType(hero, result, t)

			result.Add(domestic.NewS2cResourcePointChangeV2Msg(must.Marshal(hero.EncodeResourcePointV2(datas))))
		}
	}

	// 所有士兵属性
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.IsAllSoldierStatEffect()
	}) {
		hero.BuildingEffect().CalculateAllSoldierExtraStatEffect(hero)

		// 更新武将数据
		for _, sharedCaptain := range hero.Military().Captains() {
			sharedCaptain.CalculateProperties()
			result.Add(sharedCaptain.NewUpdateCaptainStatMsg())
		}
	}

	// 士兵属性
	raceMap := make(map[shared_proto.Race]struct{})
	for _, effect := range effects {
		if effect != nil && effect.IsSoldierStatEffect() {
			for _, race := range effect.SoldierRace {
				raceMap[race] = struct{}{}
			}
		}
	}

	if len(raceMap) > 0 {
		for race := range raceMap {
			hero.BuildingEffect().CalculateSoldierExtraStatEffect(hero, race)

			// 更新武将数据
			for _, sharedCaptain := range hero.Military().Captains() {
				if sharedCaptain.Race().Race == race {
					sharedCaptain.CalculateProperties()
					result.Add(sharedCaptain.NewUpdateCaptainStatMsg())
				}
			}
		}
	}

	// 远程士兵属性
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.IsFarSoldierEffect()
	}) {
		hero.BuildingEffect().CalculateSoldierFightTypeStatEffect(hero, true)
		// 更新武将数据
		for _, sharedCaptain := range hero.Military().Captains() {
			if sharedCaptain.Race().IsFar {
				sharedCaptain.CalculateProperties()
				result.Add(sharedCaptain.NewUpdateCaptainStatMsg())
			}
		}
	}

	// 近战士兵属性
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.IsCloseSoldierEffect()
	}) {
		hero.BuildingEffect().CalculateSoldierFightTypeStatEffect(hero, false)
		// 更新武将数据
		for _, sharedCaptain := range hero.Military().Captains() {
			if !sharedCaptain.Race().IsFar {
				sharedCaptain.CalculateProperties()
				result.Add(sharedCaptain.NewUpdateCaptainStatMsg())
			}
		}
	}

	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.IsAllSoldierStatEffect() || effect.IsSoldierStatEffect() || effect.IsFarSoldierEffect() || effect.IsCloseSoldierEffect()
	}) {
		UpdateAllTroopFightAmount(hero, result)
	}

	// 城墙属性
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.IsWallStatEffect()
	}) {
		hero.BuildingEffect().CalculateWallStat(hero)
	}

	// 城墙固定伤害
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.IsWallFixDamageEffect()
	}) {
		hero.BuildingEffect().CalculateWallFixDamage(hero)
	}

	// 税收
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.IsTaxEffect()
	}) {
		hero.BuildingEffect().CalculateTax(hero)
	}

	costReduceCoefChanged := false
	// 建筑升级消耗减少系数
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.BuildingCostReduceCoef > 0
	}) {
		hero.BuildingEffect().CalculateBuildingCostReduceCoef(hero)
		costReduceCoefChanged = true
	}

	// 科技升级消耗减少系数
	if isEffect(effects, func(effect *sub.BuildingEffectData) bool {
		return effect.TechCostReduceCoef > 0
	}) {
		hero.BuildingEffect().CalculateTechCostReduceCoef(hero)
		costReduceCoefChanged = true
	}

	if costReduceCoefChanged {
		result.Add(domestic.NewS2cUpdateCostReduceCoefMsg(
			i32.MultiF64(1000, hero.BuildingEffect().GetBuildingCostReduceCoef()),
			i32.MultiF64(1000, hero.BuildingEffect().GetTechCostReduceCoef()),
		))
	}
}

func TryMailProto(hero *entity.Hero, boolType shared_proto.HeroBoolType, data *maildata.MailData) *shared_proto.MailProto {
	if data != nil && hero.Bools().TrySet(boolType) {
		return data.NewTextMail(shared_proto.MailType_MailNormal)
	}
	return nil
}

func TrySendMailFunc(mail iface.MailModule, hero *entity.Hero, boolType shared_proto.HeroBoolType, data *maildata.MailData, ctime time.Time) func() {
	proto := TryMailProto(hero, boolType, data)
	if proto != nil {
		heroId := hero.Id()
		return func() {
			mail.SendProtoMail(heroId, proto, ctime)
		}
	}
	return nil
}

func TryAutoFullSoldier(hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas, troop *entity.Troop, ctime time.Time, duration time.Duration) {

	if troop == nil || troop.IsOutside() {
		return
	}

	if vipData := datas.GetVipLevelData(hero.VipLevel()); vipData == nil || !vipData.WallAutoFullSoldier {
		return
	}

	if hero.MiscData().GetDefenserDontAutoFullSoldier() {
		return
	}

	if ctime.Before(hero.MiscData().GetDefenserNextFullSoldierTime()) {
		return
	}
	hero.MiscData().SetDefenserNextFullSoldierTime(ctime.Add(duration))

	// 先看下还有没有空闲士兵
	freeSoldier := hero.Military().FreeSoldier(ctime)
	if freeSoldier <= 0 {
		return
	}

	originFreeSoldier := freeSoldier

	var ids, newSoldiers, fightAmounts []int32
	for _, pos := range troop.Pos() {
		captain := pos.Captain()
		if captain == nil {
			continue
		}

		toAdd := u64.Sub(captain.SoldierCapcity(), captain.Soldier())
		if toAdd <= 0 {
			continue
		}

		toAdd = u64.Min(freeSoldier, toAdd)

		captain.AddSoldier(toAdd)
		ids = append(ids, u64.Int32(captain.Id()))
		newSoldiers = append(newSoldiers, u64.Int32(captain.Soldier()))
		fightAmounts = append(fightAmounts, u64.Int32(captain.FightAmount()))

		freeSoldier = u64.Sub(freeSoldier, toAdd)
		if freeSoldier <= 0 {
			break
		}
	}

	if len(ids) > 0 {
		result.Add(military.NewS2cCaptainFullSoldierMsg(ids, newSoldiers, fightAmounts, u64.Int32(freeSoldier)))
		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_RECOVER_SOLDIER) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_RECOVER_SOLDIER)))
			UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ADD_CAPTAIN_SOLDIER)
		}
	}

	if originFreeSoldier != freeSoldier {
		hero.Military().ReduceFreeSoldier(u64.Sub(originFreeSoldier, freeSoldier), ctime)

		result.Add(hero.Military().NewUpdateFreeSoldierMsg())
	}

}

func TryUpdateTax(hero *entity.Hero, result herolock.LockResult, ctime time.Time, duration time.Duration, getBuffEffectData func(uint64) *data.BuffEffectData) {

	if hero.UpdateTax(ctime, duration, getBuffEffectData) {
		// 有更新，发送消息
		storage := hero.GetUnsafeResource()
		newGold := storage.Gold()
		newFood := storage.Food()
		newWood := storage.Wood()
		newStone := storage.Stone()
		isSafe := storage.IsSafeResource()

		result.AddFunc(func() pbutil.Buffer {
			return domestic.NewS2cResourceUpdateMsg(
				u64.Int32(newGold),
				u64.Int32(newFood),
				u64.Int32(newWood),
				u64.Int32(newStone),
				isSafe,
			)
		})
	}
}

func TryAddCaptain(hero *entity.Hero, result herolock.LockResult, data *captain.CaptainData, ctime time.Time) bool {
	// 武将已经存在
	if hero.Military().Captain(data.Id) != nil {
		return false
	}

	AddCaptainAnyway(hero, result, data, ctime)
	return true
}

func AddCaptainAnyway(hero *entity.Hero, result herolock.LockResult, data *captain.CaptainData, ctime time.Time) {
	captain := hero.NewCaptain(data, ctime)
	captain.SetTrainAccExp(data.GetInitTrainExp())
	captain.CalculateProperties()
	hero.Military().AddCaptain(captain)

	result.Add(military.NewS2cCaptainBornMsg(must.Marshal(captain.EncodeClient())))

	hero.WalkPveTroop(func(troop *entity.PveTroop) (endWalk bool) {
		if troop.AddCaptain(captain) {
			result.Add(military.NewS2cSetPveCaptainMsg(must.Marshal(troop.Encode())))
		}
		return
	})

	UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_COUNT)

	if data.GiftData != nil {
		ActivateEventLimitGift(hero, result, data.GiftData, ctime)
	}

	gamelogs.UnlockCaptainLog(hero.Pid(), hero.Sid(), hero.Id(), data.Id, uint64(data.Rarity.Color))
}

// 激活玩家事件时限礼包
func ActivateEventLimitGift(hero *entity.Hero, result herolock.LockResult, giftData *promdata.EventLimitGiftData, ctime time.Time) {
	var guanfuLevel uint64
	if guanFu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU); guanFu != nil {
		guanfuLevel = guanFu.Level
	}
	if !giftData.IsCanBuy(hero.Level(), guanfuLevel) {
		return
	}
	// 礼包购买次数超限也不会出现
	if giftData.BuyLimit > 0 && hero.Promotion().GetEventLimitGiftBuyTimes(giftData.Id) >= giftData.BuyLimit {
		return
	}
	if gift, ok := hero.Promotion().TrySetEventLimitGift(giftData, ctime); ok {
		result.Add(promotion.NewS2cNoticeEventLimitGiftMsg(gift.Encode()))
	}
}

func updateBuffPerSecond(hero *entity.Hero, result herolock.LockResult, ctime time.Time, service *Services) (toDelBuffs []*entity.BuffInfo) {
	hero.Buff().Walk(func(buff *entity.BuffInfo) {
		if !buff.EffectData.NoDuration && ctime.After(buff.EndTime) {
			hero.Buff().Del(buff.EffectData.Group, ctime)
			toDelBuffs = append(toDelBuffs, buff)
			return
		}

		service.buff.UpdatePerSecond(hero, result, buff)
	})

	return
}

func UpdateOnBuffChanged(hero *entity.Hero, result herolock.LockResult) {
	for _, c := range hero.Military().Captains() {
		if c == nil {
			continue
		}
		c.CalculateProperties()
		result.Add(c.NewUpdateCaptainStatMsg())
		UpdateTroopFightAmount(hero, c.GetTroop(), result)
	}
}

func UpdateAdvantageCount(hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas, ctime time.Time) {
	var count int
	// todo 以后保存起来
	for _, d := range datas.GetBufferTypeDataArray() {
		if hero.Buff().Buff(d.BuffGroup) != nil {
			count++
		}
	}
	if ctime.After(hero.GetMianStartTime()) && ctime.Before(hero.GetMianDisappearTime()) {
		count++
	}
	if count != hero.CurrentAdvantageCount() {
		hero.SetCurrentAdvantageCount(count)
		result.Add(domestic.NewS2cUpdateAdvantageCountMsg(int32(count)))
	}
}

func CalcTrainCurrentBuffExp(hero *entity.Hero, datas iface.ConfigDatas, ctime time.Time) (toAdd int64) {
	hero.Buff().Walk(func(buff *entity.BuffInfo) {
		toAdd += CalcTranBuffExp(hero, buff, datas, ctime)
	})

	return
}

func CalcTranBuffExp(hero *entity.Hero, buff *entity.BuffInfo, datas iface.ConfigDatas, ctime time.Time) (toUpdate int64) {
	building := hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN)
	if building == nil {
		return
	}

	endTime := hero.Military().GetGlobalTrainStartTime().Add(GetTrainingMaxDuration(hero, datas))
	endTime = timeutil.Min(endTime, ctime)
	endTime = timeutil.Min(endTime, buff.EndTime)
	startTime := timeutil.Max(hero.Military().GetCaptainTrainStartTime(), hero.Military().GetGlobalTrainStartTime())
	startTime = timeutil.Max(startTime, buff.StartTime)

	dur := endTime.Sub(startTime)
	if dur <= 0 {
		return
	}

	hours := math.Min(dur.Hours(), GetTrainingMaxDuration(hero, datas).Hours())
	exp := u64.Multi(building.Effect.TrainExpPerHour, hours*(1+hero.BuildingEffect().GetTrainCoef()))
	toUpdate = int64(buff.EffectData.CaptainTrain.CalculateByPercent(exp)) - int64(exp)

	return
}

func CalcTrainBuildingExp(hero *entity.Hero, datas iface.ConfigDatas, ctime time.Time) (toAdd uint64) {
	building := hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN)
	if building == nil {
		return
	}

	// 本次修炼结束时间
	globalStartTrainTime := hero.Military().GetGlobalTrainStartTime()
	endTime := globalStartTrainTime.Add(GetTrainingMaxDuration(hero, datas))
	endTime = timeutil.Min(endTime, ctime)
	startTime := timeutil.Max(globalStartTrainTime, hero.Military().GetCaptainTrainStartTime())

	dur := endTime.Sub(startTime)
	if dur <= 0 {
		return
	}

	hours := math.Min(dur.Hours(), GetTrainingMaxDuration(hero, datas).Hours())
	toAdd = u64.Multi(building.Effect.TrainExpPerHour, hours*(1+hero.BuildingEffect().GetTrainCoef()))

	return
}

func CalcTrainAllExp(hero *entity.Hero, datas iface.ConfigDatas, ctime time.Time) (toAdd uint64) {
	buildingExp := CalcTrainBuildingExp(hero, datas, ctime)
	currBuffExp := CalcTrainCurrentBuffExp(hero, datas, ctime)
	resevdBuffExp := hero.Military().ReservedExp()
	allExp := u64.FromInt64(int64(buildingExp) + currBuffExp + resevdBuffExp)

	return allExp
}

func VipResetDailyZero(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas, ctime time.Time) {
	vip := hero.Vip()
	vip.ResetDaily(ctime)
	result.Changed()

	if timeutil.IsSameDay(vip.DailyFirstLoginTime(), ctime) {
		return
	}
	if vip.DailyFirstLoginTime().After(ctime) {
		return
	}

	vipDailyAddExp(hctx, hero, result, datas, ctime)
}

func GmVipResetDailyZero(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas, ctime time.Time) {
	hero.Vip().ResetDaily(ctime)
	hero.Vip().IncrContinueDays() // 正常同一天不会加连续天数，强制修改
	vipDailyAddExp(hctx, hero, result, datas, ctime)
	result.Changed()
}

func vipDailyAddExp(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, datas iface.ConfigDatas, ctime time.Time) {
	vip := hero.Vip()

	if hero.Level() < datas.VipMiscData().CollectVipDailyExpMinHeroLevel {
		return
	}

	if datas.GetVipLevelData(vip.Level()).NextLevelData == nil {
		return
	}

	vip.UpdateContinueDays(ctime)

	// 计算每日vip经验
	d := datas.GetVipLevelData(vip.Level())
	toAdd := d.DailyExp
	toAdd += d.ContinueDaysExp(vip.ContinueDays())

	nextToAdd := d.DailyExp + d.ContinueDaysExp(vip.ContinueDays()+1)

	result.Add(vipMsg.NewS2cVipDailyLoginNoticeMsg(u64.Int32(vip.Level()), u64.Int32(toAdd), u64.Int32(vip.ContinueDays()), u64.Int32(nextToAdd)))

	// 发过消息再加经验
	AddVipExp(hctx, hero, result, toAdd, ctime)

	result.Changed()
}

func onVipLevelUp(hero *entity.Hero, result herolock.LockResult, oldVipData, newVipData *vip.VipLevelData, datas iface.ConfigDatas, ctime time.Time) {
	captainTrainOnVipLevelUp(hero, result, oldVipData, newVipData, datas, ctime)
	unlockWorker(hero, result, newVipData)
}

func unlockWorker(hero *entity.Hero, result herolock.LockResult, newVipData *vip.VipLevelData) {
	if hero.Domestic().UnlockWorkerForever(u64.Int(newVipData.WorkerUnlockPos)) {
		result.Add(domestic.NewS2cWorkerAlwaysUnlockMsg(u64.Int32(newVipData.WorkerUnlockPos)))
	}
}

func captainTrainOnVipLevelUp(hero *entity.Hero, result herolock.LockResult, oldVipData, newVipData *vip.VipLevelData, datas iface.ConfigDatas, ctime time.Time) {
	updateCaptainTrain(hero, result, oldVipData.CaptainTrainCoef, newVipData.CaptainTrainCoef, datas, ctime)
}

func updateCaptainTrain(hero *entity.Hero, result herolock.LockResult, oldCoef, newCoef float64, datas iface.ConfigDatas, ctime time.Time) {
	building := hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN)
	if building == nil {
		return
	}

	if oldCoef == newCoef {
		return
	}

	// 修改武将的修炼进度，计算每个武将的累积时间
	globalStartTrainTime := hero.Military().GetGlobalTrainStartTime()

	endTime := globalStartTrainTime.Add(GetTrainingMaxDuration(hero, datas))
	endTime = timeutil.Min(endTime, ctime)

	captainStartTrainTime := timeutil.Max(globalStartTrainTime, hero.Military().GetCaptainTrainStartTime())
	if d := endTime.Sub(captainStartTrainTime); d > 0 {
		hour := d.Hours() * (1 + oldCoef) / (1 + newCoef)
		newStartTime := endTime.Add(-timeutil.MultiDuration(hour, time.Hour))
		hero.Military().SetCaptainTrainStartTime(newStartTime)
	}

	result.Add(military.NewS2cUpdateTrainingMsg(
		timeutil.Marshal32(hero.Military().GetGlobalTrainStartTime()),
		timeutil.Marshal32(hero.Military().GetCaptainTrainStartTime()),
		u64.Int32(building.Effect.TrainExpPerHour),
		i32.MultiF64(1000, newCoef)))
}

func GetTrainingMaxDuration(hero *entity.Hero, datas iface.ConfigDatas) (dur time.Duration) {
	dur = datas.MilitaryConfig().TrainingMaxDuration
	if vipData := datas.GetVipLevelData(hero.VipLevel()); vipData != nil {
		dur += vipData.CaptainTrainCapacity
	}
	return
}

func DoChangeCountry(dep iface.ServiceDep, hero *entity.Hero, result herolock.LockResult, newCountry uint64) {
	oldCountry := hero.Country().GetCountryId()
	// 清掉CD
	hero.SetNextNotifyGuildTime(time.Time{})
	// 改变国家
	hero.SetCountryId(newCountry)
	result.Add(country.NewS2cHeroChangeCountryMsg(u64.Int32(newCountry),
		timeutil.Marshal32(hero.Country().GetNewUserExpiredTime()),
		timeutil.Marshal32(hero.Country().GetNormalExpiredTime())))
	// tlog
	dep.Tlog().TlogNationalFlow(hero, operate_type.CountryChange, oldCountry, newCountry)
}

// 转国
func ChangeCountry(dep iface.ServiceDep, hero *entity.Hero, result herolock.LockResult, ctime time.Time, newCountry uint64, autoBuy bool) (errMsg msg.ErrMsg) {
	var isNewUser bool
	if hero.Level() <= dep.Datas().CountryMiscData().NewHeroMaxLevel {
		isNewUser = true
	}
	if isNewUser {
		if ctime.Before(hero.Country().GetNewUserExpiredTime()) {
			errMsg = country.ErrHeroChangeCountryFailInCd
			return
		}
	} else {
		if ctime.Before(hero.Country().GetNormalExpiredTime()) {
			errMsg = country.ErrHeroChangeCountryFailInCd
			return
		}
	}

	if hero.GuildId() != 0 {
		errMsg = country.ErrHeroChangeCountryFailInGuild
		return
	}

	if isNewUser {
		hero.Country().SetNewUserExpiredTime(ctime.Add(dep.Datas().CountryMiscData().NewHeroChangeCountryCd))
	} else {
		if !TryReduceOrBuyGoods(NewContext(dep, operate_type.CountryChangeCountry), hero, result, dep.Datas().CountryMiscData().NormalChangeCountryGoods, 1, autoBuy) {
			logrus.Debugf("转国物品不足")
			errMsg = country.ErrHeroChangeCountryFailCostNotEnough
			return
		}
		hero.Country().SetNormalExpiredTime(ctime.Add(dep.Datas().CountryMiscData().NormalChangeCountryCd))
	}

	DoChangeCountry(dep, hero, result, newCountry)

	return
}

func CheckChangeCountry(hero *entity.Hero, datas iface.ConfigDatas, ctime time.Time) (errMsg msg.ErrMsg) {
	if hero.GuildId() != 0 {
		errMsg = country.ErrHeroChangeCountryFailInGuild
		return
	}

	if hero.Level() <= datas.CountryMiscData().NewHeroMaxLevel {
		if ctime.Before(hero.Country().GetNewUserExpiredTime()) {
			errMsg = country.ErrHeroChangeCountryFailInCd
		}
	} else {
		if ctime.Before(hero.Country().GetNormalExpiredTime()) {
			errMsg = country.ErrHeroChangeCountryFailInCd
			return
		}
		if !HasEnoughGoodsOrBuy(hero, datas.CountryMiscData().NormalChangeCountryGoods, 1, true) {
			errMsg = country.ErrHeroChangeCountryFailCostNotEnough
		}
	}
	return
}

func ChangeCountryAnyway(dep iface.ServiceDep, hero *entity.Hero, result herolock.LockResult, ctime time.Time, newCountry uint64) {
	if hero.Level() <= dep.Datas().CountryMiscData().NewHeroMaxLevel {
		hero.Country().SetNewUserExpiredTime(ctime.Add(dep.Datas().CountryMiscData().NewHeroChangeCountryCd))
	} else {
		hero.Country().SetNormalExpiredTime(ctime.Add(dep.Datas().CountryMiscData().NormalChangeCountryCd))
		TryReduceOrBuyGoods(NewContext(dep, operate_type.CountryChangeCountry), hero, result, dep.Datas().CountryMiscData().NormalChangeCountryGoods, 1, true)
	}
	DoChangeCountry(dep, hero, result, newCountry)
}

func CountryOfficialAppoint(hero *entity.Hero, result herolock.LockResult, toAppointOfficialData *country2.CountryOfficialData, datas iface.ConfigDatas, ctime time.Time) {
	hero.Country().SetAppointOnSameDay(true)
	if toAppointOfficialData == nil {
		logrus.Errorf("heromodule, 国家任职，找不到职位data")
		return
	}

	hero.CountryMisc().Appoint(toAppointOfficialData, ctime)
	if effect := toAppointOfficialData.BuildingEffect; effect != nil {
		UpdateBuildingEffect(hero, result, datas, ctime, effect)
	}

	if toAppointOfficialData.Body != nil {
		hero.SetBody(toAppointOfficialData.Body)
		result.Add(toAppointOfficialData.Body.ChangeBodyMsg)
	}

	result.Add(country.NewS2cOfficialAppointNoticeMsg(int32(toAppointOfficialData.OfficialType)))
	result.Changed()
	result.Ok()
}

func CountryOfficialDepose(hero *entity.Hero, result herolock.LockResult, toDeposeOldOfficialData *country2.CountryOfficialData, datas iface.ConfigDatas, ctime time.Time) {
	if toDeposeOldOfficialData == nil {
		logrus.Errorf("heromodule, 国家免职，找不到职位data")
		return
	}

	hero.CountryMisc().Depose()
	if effect := toDeposeOldOfficialData.BuildingEffect; effect != nil {
		UpdateBuildingEffect(hero, result, datas, ctime, effect)
	}

	if toDeposeOldOfficialData.Head != nil && hero.Head() == toDeposeOldOfficialData.Head.Id {
		hero.SetHead(datas.HeroInitData().DefaultHead)
		result.Add(datas.HeroInitData().DefaultHead.ChangeHeadMsg)
	}
	if toDeposeOldOfficialData.Body != nil && hero.Body() == toDeposeOldOfficialData.Body.Id {
		hero.SetBody(datas.HeroInitData().DefaultBody)
		result.Add(datas.HeroInitData().DefaultBody.ChangeBodyMsg)
	}

	result.Add(country.OFFICIAL_DEPOSE_NOTICE_S2C)
	result.Changed()
}

func GetCountryChangeNameVoteCount(hero *entity.Hero) (count uint64) {
	currentTitle := hero.TaskList().GetTitleData()
	if currentTitle == nil {
		return
	}

	return currentTitle.CountryChangeNameVoteCount
}
