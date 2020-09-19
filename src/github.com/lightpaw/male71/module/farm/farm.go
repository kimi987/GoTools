package farm

import (
	"github.com/lightpaw/male7/gen/iface"
	farmMsg "github.com/lightpaw/male7/gen/pb/farm"
	"github.com/lightpaw/male7/config/farm"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/entity/npcid"
	"time"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/service/operate_type"
	"strconv"
)

func NewFarmModule(dep iface.ServiceDep, farmService iface.FarmService, dbService iface.DbService,
	worldService iface.WorldService, datas iface.ConfigDatas, timeService iface.TimeService,
	heroDataService iface.HeroDataService, heroSnapshot iface.HeroSnapshotService,
	guildSnapShotService iface.GuildSnapshotService, ticker iface.TickerService,
	seasonService iface.SeasonService) *FarmModule {

	m := FarmModule{}
	m.dep = dep
	m.datas = datas
	m.timeService = timeService
	m.worldService = worldService
	m.heroDataService = heroDataService
	m.heroSnapshot = heroSnapshot
	m.guildSnapShotService = guildSnapShotService
	m.miscConf = datas.FarmMiscConfig()
	m.farmService = farmService
	m.seasonService = seasonService
	m.farmFunc = NewFarmFunc(dbService, worldService, datas, timeService)
	m.farmCache = NewFarmCache(timeService, ticker)

	return &m
}

//gogen:iface
type FarmModule struct {
	dep                  iface.ServiceDep
	farmService          iface.FarmService
	worldService         iface.WorldService
	datas                iface.ConfigDatas
	timeService          iface.TimeService
	guildSnapShotService iface.GuildSnapshotService
	heroDataService      iface.HeroDataService
	heroSnapshot         iface.HeroSnapshotService
	seasonService        iface.SeasonService
	miscConf             *farm.FarmMiscConfig
	farmFunc             *FarmFunc
	farmCache            *FarmCache
}

func (m *FarmModule) copyAndAddSeasonEffect(effect *data.Amount) *data.Amount {
	if effect != nil {
		b := data.NewAmountBuilder()
		b.AddAmount(effect)

		if season := m.seasonService.Season(); season != nil && season.FarmBaseInc > 0 {
			b.AddPercent(season.FarmBaseInc)
		}
		return b.Amount()
	}

	return effect
}

//gogen:iface
func (m *FarmModule) ProcessPlant(proto *farmMsg.C2SPlantProto, hc iface.HeroController) {
	cube := cb.XYCube(int(proto.CubeX), int(proto.CubeY))

	resConf := m.datas.GetFarmResConfig(u64.FromInt32(proto.ResId))
	if resConf == nil {
		hc.Send(farmMsg.ERR_PLANT_FAIL_INVALID_RES_ID)
		return
	}

	var allCubes []cb.Cube
	var npcConflictCubes []cb.Cube
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		allCubes = m.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(hero.BaseLevel())
		npcConflictCubes = hero.GetNpcConflictResourcePointOffset()

		m.dep.Tlog().TlogFarmFlow(hero, uint64(cube), strconv.FormatInt(hc.Id(), 10), uint64(operate_type.FarmTypePlant), 0, uint64(resConf.ResType))
	})
	var cubeLegal bool
	for _, c := range allCubes {
		if c == cube {
			cubeLegal = true
			break
		}
	}

	if !cubeLegal {
		hc.Send(farmMsg.ERR_PLANT_FAIL_NOT_IDLE_CUBE)
		return
	}

	for _, b := range npcConflictCubes {
		if b == cube {
			logrus.Debugf("农场种植，发来的地块是冲突状态。heroId:%v cube:%v", hc.Id(), cube)
			hc.Send(farmMsg.ERR_PLANT_FAIL_NOT_IDLE_CUBE)
			return
		}
	}

	var errMsg msg.ErrMsg
	m.farmService.FuncWait("ProcessPlant", farmMsg.ERR_PLANT_FAIL_SERVER_ERR, hc, func() {
		errMsg = m.farmFunc.Plant(hc, cube, resConf)
	})
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}

	x, y := cube.XY()
	hc.Send(farmMsg.NewS2cPlantMsg(int32(x), int32(y), u64.Int32(resConf.Id)))
}

//gogen:iface
func (m *FarmModule) ProcessHarvest(proto *farmMsg.C2SHarvestProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.FarmHarvest)

	cube := cb.XYCube(int(proto.CubeX), int(proto.CubeY))

	var errMsg msg.ErrMsg
	var fc *entity.FarmCube
	m.farmService.FuncWait("ProcessHarvest", farmMsg.ERR_HARVEST_FAIL_SERVER_ERR, hc, func() {
		errMsg, fc = m.farmFunc.Harvest(hc, cube)
	})
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}
	if fc == nil {
		hc.Send(farmMsg.ERR_HARVEST_FAIL_SERVER_ERR)
		return
	}

	ctime := m.timeService.CurrentTime()
	plantDuration := ctime.Sub(fc.StartTime)

	var heroName string
	var guildId int64
	var output uint64
	var resType shared_proto.ResType
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroName = hero.Name()
		guildId = hero.GuildId()

		resConf := fc.ResConf
		var effect *data.Amount
		if resConf.ResType == shared_proto.ResType_GOLD {
			effect = hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_GOLD)
		} else if resConf.ResType == shared_proto.ResType_STONE {
			effect = hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_STONE)
		} else {
			logrus.Debugf("农场收获，没有对应的产出加成 res:", resConf.ResType)
			effect = &data.Amount{}
		}
		effect = m.copyAndAddSeasonEffect(effect)

		resType = fc.ResType()
		output = entity.CalcHarvestOutput(plantDuration, fc.StealTimes, effect, resConf, m.miscConf)
		if output > 0 {
			heromodule.AddUnsafeSingleResource(hctx, hero, result, resConf.ResType, output)
			x, y := cube.XY()
			result.Add(farmMsg.NewS2cHarvestMsg(int32(x), int32(y), u64.Int32(output)))

			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AccumFarmHarvestRes, output)
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumFarmHarvestRes, uint64(resConf.ResType), output)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST)
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_FarmHarvestTimes)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST_TIMES)
		} else {
			result.Add(farmMsg.ERR_HARVEST_FAIL_NO_OUTPUT)
			return
		}

		x, y := cube.XY()
		result.Add(farmMsg.NewS2cHarvestMsg(int32(x), int32(y), u64.Int32(output)))

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_COLLECT_FARM) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_COLLECT_FARM)))
		}

		m.dep.Tlog().TlogFarmFlow(hero, uint64(cube), strconv.FormatInt(hc.Id(), 10), uint64(operate_type.FarmTypeHarvest), uint64(resType), uint64(resType))

		result.Changed()
		result.Ok()
	}) {
		return
	}

	var guildFlag string
	var countryId uint64
	if guildId > 0 {
		g := m.guildSnapShotService.GetSnapshot(guildId)
		if g != nil {
			guildFlag = g.FlagName
			countryId = g.Country.Id
		}
	}

	m.farmService.FuncNoWait(func() {
		var goldOutput, stoneOutput uint64
		if output > 0 {
			if resType == shared_proto.ResType_GOLD {
				goldOutput = output
			} else if resType == shared_proto.ResType_STONE {
				stoneOutput = output
			}

			m.farmFunc.AddLog(hc.Id(), heroName, guildFlag, countryId, hc.Id(), goldOutput, stoneOutput, ctime, shared_proto.FarmLogType_Harvest)
		}
	})

}

//gogen:iface
func (m *FarmModule) ProcessChange(proto *farmMsg.C2SChangeProto, hc iface.HeroController) {
	cube := cb.XYCube(int(proto.CubeX), int(proto.CubeY))
	newResConf := m.datas.GetFarmResConfig(u64.FromInt32(proto.ResId))
	if newResConf == nil {
		hc.Send(farmMsg.ERR_CHANGE_FAIL_INVALID_CUBE)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.FarmChange)

	var resConf *farm.FarmResConfig
	var stealTimes uint64
	var plantDuration time.Duration
	var errMsg msg.ErrMsg
	m.farmService.FuncWait("ProcessChange", farmMsg.ERR_CHANGE_FAIL_SERVER_ERR, hc, func() {
		errMsg, resConf, plantDuration, stealTimes = m.farmFunc.Change(hc, cube, newResConf)
	})
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}

	var heroName string
	var guildId int64
	var output uint64
	var resType shared_proto.ResType
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroName = hero.Name()
		guildId = hero.GuildId()
		var effect *data.Amount
		if resConf.ResType == shared_proto.ResType_GOLD {
			effect = hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_GOLD)
		} else if resConf.ResType == shared_proto.ResType_STONE {
			effect = hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_STONE)
		} else {
			logrus.Debugf("农场重建，没有对应的产出加成 res:", resConf.ResType)
			effect = &data.Amount{}
		}
		effect = m.copyAndAddSeasonEffect(effect)

		output = entity.CalcHarvestOutput(plantDuration, stealTimes, effect, resConf, m.miscConf)
		if output > 0 {
			heromodule.AddUnsafeSingleResource(hctx, hero, result, resConf.ResType, output)

			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AccumFarmHarvestRes, output)
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumFarmHarvestRes, uint64(resConf.ResType), output)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST)
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_FarmHarvestTimes)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST_TIMES)

			if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_COLLECT_FARM) {
				result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_COLLECT_FARM)))
			}
		}

		x, y := cube.XY()
		result.Add(farmMsg.NewS2cChangeMsg(
			int32(x),
			int32(y),
			u64.Int32(newResConf.Id),
			u64.Int32(resConf.Id),
			u64.Int32(output),
		))

		m.dep.Tlog().TlogFarmFlow(hero, uint64(cube), strconv.FormatInt(hc.Id(), 10), uint64(operate_type.FarmTypeChange), uint64(resConf.ResType), uint64(newResConf.ResType))

		result.Changed()
		result.Ok()
	})

	var guildFlag string
	var countryId uint64
	if guildId > 0 {
		g := m.guildSnapShotService.GetSnapshot(guildId)
		if g != nil {
			guildFlag = g.FlagName
			countryId = g.Country.Id
		}
	}

	ctime := m.timeService.CurrentTime()
	m.farmService.FuncNoWait(func() {
		var goldOutput, stoneOutput uint64
		if output > 0 {
			if resType == shared_proto.ResType_GOLD {
				goldOutput = output
			} else if resType == shared_proto.ResType_STONE {
				stoneOutput = output
			}

			m.farmFunc.AddLog(hc.Id(), heroName, guildFlag, countryId, hc.Id(), goldOutput, stoneOutput, ctime, shared_proto.FarmLogType_Harvest)
		}
	})

}

//gogen:iface c2s_one_key_harvest
func (m *FarmModule) ProcessOneKeyHarvest(proto *farmMsg.C2SOneKeyHarvestProto, hc iface.HeroController) {
	clientResType := proto.ResType
	var errMsg msg.ErrMsg
	var harvestCubes []*entity.FarmCube
	m.farmService.FuncWait("ProcessOneKeyHarvest", farmMsg.ERR_ONE_KEY_HARVEST_FAIL_SERVER_ERR, hc, func() {
		errMsg, harvestCubes = m.farmFunc.OneKeyHarvest(hc, clientResType)
	})
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.FarmOneKeyHarvest)

	var heroName string
	var guildId int64
	var goldOutput, stoneOutput uint64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroName = hero.Name()
		guildId = hero.GuildId()

		goldEffect := hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_GOLD)
		stoneEffect := hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_STONE)

		goldEffect = m.copyAndAddSeasonEffect(goldEffect)
		stoneEffect = m.copyAndAddSeasonEffect(stoneEffect)

		msgProto := &farmMsg.S2COneKeyHarvestProto{}

		var effect *data.Amount
		for _, harvestCube := range harvestCubes {
			if harvestCube == nil {
				continue
			}

			x, y := harvestCube.Cube.XY()
			msgProto.CubeX = append(msgProto.CubeX, int32(x))
			msgProto.CubeY = append(msgProto.CubeY, int32(y))

			if harvestCube.ResType() == shared_proto.ResType_GOLD {
				effect = goldEffect
			} else if harvestCube.ResType() == shared_proto.ResType_STONE {
				effect = stoneEffect
			} else {
				logrus.Debugf("农场收获，没有对应的产出加成 res:", harvestCube.ResType())
				effect = &data.Amount{}
			}

			output := entity.CalcHarvestOutput(harvestCube.ResConf.RipeDuration, harvestCube.StealTimes, effect, harvestCube.ResConf, m.miscConf)
			if harvestCube.ResType() == shared_proto.ResType_GOLD {
				goldOutput += output
				msgProto.GoldOutput = append(msgProto.GoldOutput, u64.Int32(output))
				msgProto.StoneOutput = append(msgProto.StoneOutput, 0)
			} else if harvestCube.ResType() == shared_proto.ResType_STONE {
				stoneOutput += output
				msgProto.GoldOutput = append(msgProto.GoldOutput, 0)
				msgProto.StoneOutput = append(msgProto.StoneOutput, u64.Int32(output))
			}
		}
		if goldOutput > 0 || stoneOutput > 0 {
			heromodule.AddUnsafeResource(hctx, hero, result, goldOutput, 0, 0, stoneOutput)

			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AccumFarmHarvestRes, goldOutput+stoneOutput)
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumFarmHarvestRes, uint64(shared_proto.ResType_GOLD), goldOutput)
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumFarmHarvestRes, uint64(shared_proto.ResType_STONE), stoneOutput)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST)
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_FarmHarvestTimes)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST_TIMES)
		}

		result.Add(farmMsg.NewS2cOneKeyHarvestProtoMsg(msgProto))

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_COLLECT_FARM) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_COLLECT_FARM)))
		}

		m.dep.Tlog().TlogFarmFlow(hero, 0, strconv.FormatInt(hero.Id(), 10), operate_type.FarmTypeReset, 0, 0)

		result.Changed()
		result.Ok()
	})

	var guildFlag string
	var countryId uint64
	if guildId > 0 {
		g := m.guildSnapShotService.GetSnapshot(guildId)
		if g != nil {
			guildFlag = g.FlagName
			countryId = g.Country.Id
		}
	}

	ctime := m.timeService.CurrentTime()
	m.farmService.FuncNoWait(func() {
		if goldOutput > 0 || stoneOutput > 0 {
			m.farmFunc.AddLog(hc.Id(), heroName, guildFlag, countryId, hc.Id(), goldOutput, stoneOutput, ctime, shared_proto.FarmLogType_Harvest)
		}
	})
}

//gogen:iface c2s_one_key_reset
func (m *FarmModule) ProcessOneKeyReset(hc iface.HeroController) {
	var errMsg msg.ErrMsg
	var harvestCubes []*entity.FarmCube
	m.farmService.FuncWait("ProcessOneKeyReset", farmMsg.ERR_ONE_KEY_RESET_FAIL_SERVER_ERR, hc, func() {
		errMsg, harvestCubes = m.farmFunc.OneKeyReset(hc)
	})
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.FarmOneKeyReset)

	ctime := m.timeService.CurrentTime()
	var heroName string
	var guildId int64
	var goldOutput, stoneOutput uint64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroName = hero.Name()
		guildId = hero.GuildId()

		goldEffect := hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_GOLD)
		stoneEffect := hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_STONE)

		goldEffect = m.copyAndAddSeasonEffect(goldEffect)
		stoneEffect = m.copyAndAddSeasonEffect(stoneEffect)

		msgProto := &farmMsg.S2COneKeyResetProto{}

		for _, harvestCube := range harvestCubes {
			if harvestCube == nil {
				continue
			}

			x, y := harvestCube.Cube.XY()
			msgProto.CubeX = append(msgProto.CubeX, int32(x))
			msgProto.CubeY = append(msgProto.CubeY, int32(y))

			var effect *data.Amount
			if harvestCube.ResType() == shared_proto.ResType_GOLD {
				effect = goldEffect
			} else if harvestCube.ResType() == shared_proto.ResType_STONE {
				effect = stoneEffect
			} else {
				logrus.Debugf("农场一键犁地，没有对应的产出加成 res:", harvestCube.ResType())
				effect = &data.Amount{}
			}

			var plantDuration time.Duration
			if timeutil.IsZero(harvestCube.StartTime) {
				plantDuration, _ = time.ParseDuration("0s")
			} else {
				plantDuration = ctime.Sub(harvestCube.StartTime)
			}

			output := entity.CalcHarvestOutput(plantDuration, harvestCube.StealTimes, effect, harvestCube.ResConf, m.miscConf)
			if harvestCube.ResType() == shared_proto.ResType_GOLD {
				goldOutput += output
				msgProto.GoldOutput = append(msgProto.GoldOutput, u64.Int32(output))
				msgProto.StoneOutput = append(msgProto.StoneOutput, 0)
			} else if harvestCube.ResType() == shared_proto.ResType_STONE {
				stoneOutput += output
				msgProto.GoldOutput = append(msgProto.GoldOutput, 0)
				msgProto.StoneOutput = append(msgProto.StoneOutput, u64.Int32(output))
			}
		}
		if goldOutput > 0 || stoneOutput > 0 {
			heromodule.AddUnsafeResource(hctx, hero, result, goldOutput, 0, 0, stoneOutput)

			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AccumFarmHarvestRes, goldOutput+stoneOutput)
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumFarmHarvestRes, uint64(shared_proto.ResType_GOLD), goldOutput)
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumFarmHarvestRes, uint64(shared_proto.ResType_STONE), stoneOutput)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST)
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_FarmHarvestTimes)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST_TIMES)
		}

		result.Add(farmMsg.NewS2cOneKeyResetProtoMsg(msgProto))

		m.dep.Tlog().TlogFarmFlow(hero, 0, strconv.FormatInt(hero.Id(), 10), operate_type.FarmTypeReset, 0, 0)

		result.Changed()
		result.Ok()
	})

	var guildFlag string
	var countryId uint64
	if guildId > 0 {
		g := m.guildSnapShotService.GetSnapshot(guildId)
		if g != nil {
			guildFlag = g.FlagName
			countryId = g.Country.Id
		}
	}

	m.farmService.FuncNoWait(func() {
		if goldOutput > 0 || stoneOutput > 0 {
			m.farmFunc.AddLog(hc.Id(), heroName, guildFlag, countryId, hc.Id(), goldOutput, stoneOutput, ctime, shared_proto.FarmLogType_Harvest)
		}
	})

}

//gogen:iface c2s_one_key_plant
func (m *FarmModule) ProcessOneKeyPlant(proto *farmMsg.C2SOneKeyPlantProto, hc iface.HeroController) {
	if proto.GoldCount <= 0 && proto.StoneCount <= 0 {
		hc.Send(farmMsg.ERR_ONE_KEY_PLANT_FAIL_INVALID_COUNT)
		return
	}

	var baseLevel uint64
	var allCubes []cb.Cube
	var npcConflictCubes []cb.Cube

	hc.Func(func(hero *entity.Hero, err error) (heroChanged bool) {
		baseLevel = hero.BaseLevel()
		allCubes = m.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(baseLevel)
		npcConflictCubes = hero.GetNpcConflictResourcePointOffset()

		m.dep.Tlog().TlogFarmFlow(hero, 0, strconv.FormatInt(hero.Id(), 10), operate_type.FarmTypePlant, 0, 0)
		return false
	})

	okConf := m.datas.GetFarmOneKeyConfig(baseLevel)
	if okConf == nil {
		logrus.Debugf("农场一键种植，找不到%级主城的配置", baseLevel)
		hc.Send(farmMsg.ERR_ONE_KEY_PLANT_FAIL_SERVER_ERR)
		return
	}

	goldResConf := m.datas.GetFarmResConfig(u64.FromInt32(proto.GoldConfId))
	if goldResConf == nil {
		goldResConf = m.miscConf.OneKeyResConfig[shared_proto.ResType_GOLD]
	}
	stoneResConf := m.datas.GetFarmResConfig(u64.FromInt32(proto.StoneConfId))
	if stoneResConf == nil {
		stoneResConf = m.miscConf.OneKeyResConfig[shared_proto.ResType_STONE]
	}

	var errMsg msg.ErrMsg
	m.farmService.FuncWait("ProcessOneKeyPlant", farmMsg.ERR_ONE_KEY_PLANT_FAIL_SERVER_ERR, hc, func() {
		errMsg = m.farmFunc.OneKeyPlant(hc, okConf, allCubes, npcConflictCubes, goldResConf, stoneResConf, u64.FromInt32(proto.GoldCount), u64.FromInt32(proto.StoneCount))
	})
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}
}

//gogen:iface
func (m *FarmModule) ProcessSteal(proto *farmMsg.C2SStealProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.Target)
	if !ok {
		logrus.Debug("偷菜，无效的目标")
		hc.Send(farmMsg.ERR_STEAL_FAIL_INVALID_TARGET)
		return
	}

	if npcid.IsNpcId(targetId) {
		logrus.Debug("偷菜，目标是个npc?")
		hc.Send(farmMsg.ERR_STEAL_FAIL_INVALID_TARGET)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.FarmSteal)
	cube := cb.XYCube(int(proto.CubeX), int(proto.CubeY))

	// 今天还能偷多少
	var canStealGold, canStealStone uint64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		dailyStealGold := u64.FromInt32(hero.FarmExtra().GetDailyStealGold())
		dailyStealStone := u64.FromInt32(hero.FarmExtra().GetDailyStealStone())

		guanFuLevel := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU).Level
		maxDailyStealGoldAmount := m.farmFunc.MaxDailyStealGoldAmount(guanFuLevel)
		maxDailyStealStoneAmount := m.farmFunc.MaxDailyStealStoneAmount(guanFuLevel)

		if dailyStealGold >= maxDailyStealGoldAmount && dailyStealStone >= maxDailyStealStoneAmount {
			result.Add(farmMsg.ERR_STEAL_FAIL_DAILY_STEAL_FULL)
			return
		}
		canStealGold = u64.Sub(maxDailyStealGoldAmount, dailyStealGold)
		canStealStone = u64.Sub(maxDailyStealStoneAmount, dailyStealStone)

		m.dep.Tlog().TlogFarmFlow(hero, uint64(cube), strconv.FormatInt(targetId, 10), operate_type.FarmTypeSteal, 0, 0)
		result.Ok()
	}) {
		return
	}

	// 农场收获加成
	var goldEffect *data.Amount
	var stoneEffect *data.Amount
	m.heroDataService.Func(targetId, func(hero *entity.Hero, err error) (heroChanged bool) {
		if err != nil {
			logrus.Debug("偷菜，lock hero 失败")
			hc.Send(farmMsg.ERR_STEAL_FAIL_INVALID_TARGET)
			return
		}
		goldEffect = hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_GOLD)
		stoneEffect = hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_STONE)

		goldEffect = m.copyAndAddSeasonEffect(goldEffect)
		stoneEffect = m.copyAndAddSeasonEffect(stoneEffect)
		return false
	})

	// 偷菜
	ctime := m.timeService.CurrentTime()
	var errMsg msg.ErrMsg
	var goldOutput, stoneOutput uint64
	m.farmService.FuncWait("ProcessSteal", farmMsg.ERR_STEAL_FAIL_SERVER_ERR, hc, func() {
		errMsg, goldOutput, stoneOutput = m.farmFunc.Steal(hc, targetId, cube, ctime, goldEffect, stoneEffect, canStealGold, canStealStone)
	})
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}

	// 把偷的菜加到身上
	var heroName string
	var guildId int64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if goldOutput > 0 || stoneOutput > 0 {
			heromodule.AddUnsafeResource(hctx, hero, result, goldOutput, 0, 0, stoneOutput)
			hero.FarmExtra().SetDailyStealGold(hero.FarmExtra().GetDailyStealGold() + u64.Int32(goldOutput))
			hero.FarmExtra().SetDailyStealStone(hero.FarmExtra().GetDailyStealStone() + u64.Int32(stoneOutput))

			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AccumFarmStealRes, goldOutput+stoneOutput)
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumFarmStealRes, uint64(shared_proto.ResType_GOLD), goldOutput)
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumFarmStealRes, uint64(shared_proto.ResType_STONE), stoneOutput)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_FARM_STEAL)
		}

		heroName = hero.Name()
		guildId = hero.GuildId()

		result.Changed()
		result.Ok()
	})

	var guildFlag string
	var countryId uint64
	if guildId > 0 {
		g := m.guildSnapShotService.GetSnapshot(guildId)
		if g != nil {
			guildFlag = g.FlagName
			countryId = g.Country.Id
		}
	}

	m.farmService.FuncNoWait(func() {
		if goldOutput > 0 || stoneOutput > 0 {
			m.farmFunc.AddLog(hc.Id(), heroName, guildFlag, countryId, targetId, goldOutput, stoneOutput, ctime, shared_proto.FarmLogType_Steal)
		}
	})

}

//gogen:iface
func (m *FarmModule) ProcessOneKeySteal(proto *farmMsg.C2SOneKeyStealProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.Target)
	if !ok {
		logrus.Debug("一键偷菜，无效的目标")
		hc.Send(farmMsg.ERR_ONE_KEY_STEAL_FAIL_INVALID_TARGET)
		return
	}

	if npcid.IsNpcId(targetId) {
		logrus.Debug("一键偷菜，目标是个npc?")
		hc.Send(farmMsg.ERR_ONE_KEY_STEAL_FAIL_INVALID_TARGET)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.FarmOneKeySteal)

	// 今天还能偷多少
	var canStealGold, canStealStone uint64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		dailyStealGold := u64.FromInt32(hero.FarmExtra().GetDailyStealGold())
		dailyStealStone := u64.FromInt32(hero.FarmExtra().GetDailyStealStone())

		guanFuLevel := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU).Level
		if dailyStealGold >= m.farmFunc.MaxDailyStealGoldAmount(guanFuLevel) && dailyStealStone >= m.farmFunc.MaxDailyStealStoneAmount(guanFuLevel) {
			result.Add(farmMsg.ERR_ONE_KEY_STEAL_FAIL_DAILY_STEAL_FULL)
			return
		}
		canStealGold = u64.Sub(m.farmFunc.MaxDailyStealGoldAmount(guanFuLevel), dailyStealGold)
		canStealStone = u64.Sub(m.farmFunc.MaxDailyStealStoneAmount(guanFuLevel), dailyStealStone)

		result.Ok()
	}) {
		return
	}

	// 农场加成
	var goldEffect *data.Amount
	var stoneEffect *data.Amount
	m.heroDataService.Func(targetId, func(hero *entity.Hero, err error) (heroChanged bool) {
		if err != nil {
			logrus.Debug("偷菜，lock hero 失败")
			hc.Send(farmMsg.ERR_ONE_KEY_STEAL_FAIL_INVALID_TARGET)
			return
		}
		goldEffect = hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_GOLD)
		stoneEffect = hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_STONE)

		goldEffect = m.copyAndAddSeasonEffect(goldEffect)
		stoneEffect = m.copyAndAddSeasonEffect(stoneEffect)
		return false
	})

	// 偷菜
	ctime := m.timeService.CurrentTime()
	var goldOutput uint64
	var stoneOutput uint64
	var errMsg msg.ErrMsg
	m.farmService.FuncWait("ProcessOneKeySteal", farmMsg.ERR_ONE_KEY_STEAL_FAIL_SERVER_ERR, hc, func() {
		errMsg, goldOutput, stoneOutput = m.farmFunc.OneKeySteal(hc, targetId, ctime, goldEffect, stoneEffect, canStealGold, canStealStone)
	})
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}

	// 把偷的菜加到身上
	var heroName string
	var guildId int64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if goldOutput > 0 || stoneOutput > 0 {
			heromodule.AddUnsafeResource(hctx, hero, result, goldOutput, 0, 0, stoneOutput)
			hero.FarmExtra().SetDailyStealGold(hero.FarmExtra().GetDailyStealGold() + u64.Int32(goldOutput))
			hero.FarmExtra().SetDailyStealStone(hero.FarmExtra().GetDailyStealStone() + u64.Int32(stoneOutput))

			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AccumFarmStealRes, goldOutput+stoneOutput)
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumFarmStealRes, uint64(shared_proto.ResType_GOLD), goldOutput)
			hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumFarmStealRes, uint64(shared_proto.ResType_STONE), stoneOutput)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_FARM_STEAL)
		}
		heroName = hero.Name()
		guildId = hero.GuildId()

		m.dep.Tlog().TlogFarmFlow(hero, 0, strconv.FormatInt(targetId, 10), operate_type.FarmTypeSteal, 0, 0)

		result.Changed()
		result.Ok()
	})

	var guildFlag string
	var countryId uint64
	if guildId > 0 {
		g := m.guildSnapShotService.GetSnapshot(guildId)
		if g != nil {
			guildFlag = g.FlagName
			countryId = g.Country.Id
		}
	}

	m.farmService.FuncNoWait(func() {
		if goldOutput > 0 || stoneOutput > 0 {
			m.farmFunc.AddLog(hc.Id(), heroName, guildFlag, countryId, targetId, goldOutput, stoneOutput, ctime, shared_proto.FarmLogType_Steal)
		}
	})

}

//gogen:iface
func (m *FarmModule) ProcessViewFarm(proto *farmMsg.C2SViewFarmProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.Target)
	if !ok {
		hc.Send(farmMsg.ERR_VIEW_FARM_FAIL_INVALID_TARGET)
		return
	}

	var baseLevel uint64
	var lastBaseLevel uint64
	var heroBasic []byte
	var firstView bool
	var goldOutputEffect *data.Amount
	var stoneOutputEffect *data.Amount
	var canSteal bool

	ctime := m.timeService.CurrentTime()

	npc := make(map[cb.Cube]*entity.HomeNpcBase)

	if hc.Id() != targetId {
		heroSnapShot := m.heroSnapshot.Get(targetId)
		if heroSnapShot != nil {
			heroBasic = heroSnapShot.EncodeBasic4ClientBytes()
			baseLevel = heroSnapShot.BaseLevel
		}

		dbSucc, canStealCount := m.farmFunc.loadCanStealCount(targetId, hc.Id(), ctime)
		if dbSucc && canStealCount > 0 {
			hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
				if u64.FromInt32(hero.FarmExtra().GetDailyStealGold()) >= m.farmFunc.MaxDailyStealGoldAmount(hero.BaseLevel()) &&
					u64.FromInt32(hero.FarmExtra().GetDailyStealStone()) >= m.farmFunc.MaxDailyStealStoneAmount(hero.BaseLevel()) {
					canSteal = true
				}
				return
			})
		}
	} else {
		hc.Func(func(hero *entity.Hero, err error) (heroChanged bool) {
			baseLevel = hero.BaseLevel()
			lastBaseLevel = hero.Misc().FarmUnlockAnimationBaseLevel
			if proto.OpenWin {
				hero.Misc().FarmUnlockAnimationBaseLevel = baseLevel
			}
			firstView = !hero.Bools().Get(shared_proto.HeroBoolType_BOOL_COLLECT_FARM) && !hero.Bools().Get(shared_proto.HeroBoolType_BOOL_VIEW_FARM)
			goldOutputEffect = hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_GOLD)
			stoneOutputEffect = hero.Domestic().GetFarmExtraOutput(shared_proto.ResType_STONE)

			goldOutputEffect = m.copyAndAddSeasonEffect(goldOutputEffect)
			stoneOutputEffect = m.copyAndAddSeasonEffect(stoneOutputEffect)

			hero.WalkHomeNpcBase(func(base *entity.HomeNpcBase) (toBreak bool) {
				cube := cb.XYCube(base.GetData().EvenOffsetX, base.GetData().EvenOffsetY)
				npc[cube] = base
				return false
			})

			return true
		})
	}

	msgProto := &farmMsg.S2CViewFarmProto{}

	// 所有农场地块
	allCubes := m.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(baseLevel)

	errMsg, farmCubes := m.farmFunc.GetFarm(targetId, allCubes)
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}

	if proto.OpenWin && firstView && hc.Id() == targetId {
		// 农场新手引导
		firstLevel := baseLevel
		if firstLevel <= 0 {
			firstLevel = 1
		}
		firstLevelCubes := m.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(1)
		farmCubes = m.farmFunc.FirstGetFarm(hc, m.datas.GetFarmOneKeyConfig(firstLevel), firstLevelCubes, farmCubes)
	}

	if hc.Id() == targetId {
		if proto.OpenWin {
			// 下一级即将开启
			nextBaseLevel := baseLevel + 1
			nextCubes := m.datas.RegionConfig().GetEvenOffsetCubesOnlyCurrentLevel(nextBaseLevel)
			if nextCubes != nil {
				for _, cube := range nextCubes {
					x, y := cube.XY()
					msgProto.NextLevelCubeX = append(msgProto.NextLevelCubeX, int32(x))
					msgProto.NextLevelCubeY = append(msgProto.NextLevelCubeY, int32(y))
				}
			}

			if !firstView {
				// 处理解锁动画
				showAnimationCubes := make([]cb.Cube, 0)
				for level := baseLevel; level > lastBaseLevel; level-- {
					thisCubes := m.datas.RegionConfig().GetEvenOffsetCubesOnlyCurrentLevel(level)
					showAnimationCubes = append(showAnimationCubes, thisCubes...)
				}

				if len(showAnimationCubes) > 0 {
					for _, fc := range farmCubes {
						for _, scube := range showAnimationCubes {
							if fc.Cube == scube {
								fc.ShowUnlockAnimation = true
								break
							}
						}
					}
				}
			}
		}
	}

	farmProto := &shared_proto.HeroFarmProto{}

	if goldOutputEffect != nil {
		farmProto.GoldEffect = goldOutputEffect.Encode()
	}
	if stoneOutputEffect != nil {
		farmProto.StoneEffect = stoneOutputEffect.Encode()
	}

	targetIdBytes := idbytes.ToBytes(targetId)
	farmProto.HeroId = targetIdBytes
	for _, farmCube := range farmCubes {
		proto := farmCube.Encode()
		if _, ok := npc[farmCube.Cube]; ok {
			proto.NpcConflicted = true
		}
		farmProto.FarmCube = append(farmProto.FarmCube, proto)
	}
	farmProtoBytes, err := farmProto.Marshal()
	if err != nil {
		hc.Send(farmMsg.ERR_VIEW_FARM_FAIL_SERVER_ERR)
		return
	}

	msgProto.Target = targetIdBytes
	msgProto.HeroFarm = farmProtoBytes
	if len(heroBasic) > 0 {
		msgProto.TargetBasic = heroBasic
	}
	msgProto.CanSteal = canSteal

	hc.Send(farmMsg.NewS2cViewFarmProtoMsg(msgProto))
}

//gogen:iface
func (m *FarmModule) ProcessStealLogList(proto *farmMsg.C2SStealLogListProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.Target)
	if !ok {
		hc.Send(farmMsg.ERR_STEAL_LOG_LIST_FAIL_INVALID_TARGET)
		return
	}

	if proto.Newest {
		if succMsg := m.farmCache.LoadNewestLog(targetId); succMsg != nil {
			hc.Send(succMsg)
			return
		}

		succMsg, errMsg := m.farmFunc.FarmLogList(targetId, 1, true)
		if errMsg != nil {
			hc.Send(errMsg.ErrMsg())
			return
		}

		hc.Send(succMsg)
		m.farmCache.StoreNewestLog(targetId, succMsg)
	} else {
		size := m.miscConf.StealLogMaxCount
		if succMsg := m.farmCache.LoadLog(targetId); succMsg != nil {
			hc.Send(succMsg)
			return
		}

		succMsg, errMsg := m.farmFunc.FarmLogList(targetId, size, false)
		if errMsg != nil {
			hc.Send(errMsg.ErrMsg())
			return
		}

		hc.Send(succMsg)
		m.farmCache.StoreLog(targetId, succMsg)
	}
}

//gogen:iface c2s_can_steal_list
func (m *FarmModule) ProcessCanStealList(hc iface.HeroController) {
	if succMsg := m.farmCache.LoadStealList(hc.Id()); succMsg != nil {
		hc.Send(succMsg)
		return
	}

	var relationIds []int64
	hc.Func(func(hero *entity.Hero, err error) (heroChanged bool) {
		relationIds = hero.Relation().FriendAndEnemyIds()
		return false
	})

	succMsg, errMsg := m.farmFunc.GetRelationCanStealList(hc, relationIds)
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}
	hc.Send(succMsg)
	m.farmCache.StoreStealList(hc.Id(), succMsg)

}
