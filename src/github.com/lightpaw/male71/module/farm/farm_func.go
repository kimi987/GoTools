package farm

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/farm"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/gen/iface"
	farmMsg "github.com/lightpaw/male7/gen/pb/farm"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
)

func NewFarmFunc(dbService iface.DbService, worldService iface.WorldService, datas iface.ConfigDatas, timeService iface.TimeService) *FarmFunc {
	m := FarmFunc{}
	m.datas = datas
	m.worldService = worldService
	m.timeService = timeService
	m.miscConf = datas.FarmMiscConfig()
	m.db = dbService

	return &m
}

type FarmFunc struct {
	datas        iface.ConfigDatas
	worldService iface.WorldService
	timeService  iface.TimeService
	miscConf     *farm.FarmMiscConfig
	db           iface.DbService
}

func (m *FarmFunc) GetFarm(heroId int64, allCubes []cb.Cube) (errMsg msg.ErrMsg, farmCubes []*entity.FarmCube) {
	dbSucc, farmCubes := m.loadFarmCubes(heroId)
	if !dbSucc {
		errMsg = farmMsg.ErrViewFarmFailServerErr
		return
	}

	for _, cube := range allCubes {
		if fc := getCube(farmCubes, cube); fc == nil {
			fc = entity.NewIdleFarmCube(heroId, cube)
			farmCubes = append(farmCubes, fc)
		}
	}

	return
}

func (m *FarmFunc) FirstGetFarm(hc iface.HeroController, okConf *farm.FarmOneKeyConfig, firstLevelCubes []cb.Cube, currFarmCubes []*entity.FarmCube) (farmCubes []*entity.FarmCube) {
	ctime := m.timeService.CurrentTime()
	goldCount := okConf.GoldHopeCubeCount
	stoneCount := okConf.StoneHopeCubeCount

	farmCubes = append(farmCubes, currFarmCubes...)

	var dbAllSucc bool
	ForEachCubes(firstLevelCubes, func(cube cb.Cube) {
		var resConf *farm.FarmResConfig
		if goldCount > 0 {
			resConf = m.miscConf.OneKeyResConfig[shared_proto.ResType_GOLD]
			goldCount--
		} else {
			resConf = m.miscConf.OneKeyResConfig[shared_proto.ResType_STONE]
			stoneCount--
		}

		startTime := ctime.Add(-resConf.RipeDuration)
		farmCube := entity.NewFarmCube(hc.Id(), cube, startTime, ctime, resConf, m.datas.FarmMiscConfig())
		farmCubes = append(farmCubes, farmCube)
		dbSucc := m.createFarmCube(farmCube)
		if !dbSucc {
			logrus.Errorf("FarmFunc.GetFirstFarm CreateFarmCube 失败 heroId:%v cube:%v", hc.Id(), farmCube.Cube)
		}
		// 有一个成功就可以
		dbAllSucc = dbAllSucc || dbSucc
	})

	if dbAllSucc {
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_VIEW_FARM) {
				result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_VIEW_FARM)))
			}
		})
	}

	return farmCubes
}

func (m *FarmFunc) Plant(hc iface.HeroController, cube cb.Cube, resConf *farm.FarmResConfig) (errMsg msg.ErrMsg) {
	heroId := hc.Id()

	dbSucc, fc := m.loadFarmCube(heroId, cube)
	if !dbSucc {
		return farmMsg.ErrPlantFailServerErr
	}

	var buffs []*entity.BuffInfo
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		buffs = hero.Buff().Buffs(shared_proto.BuffEffectType_Buff_ET_farm_harvest)
		return
	})

	ctime := m.timeService.CurrentTime()
	orgRipeTime := ctime.Add(resConf.RipeDuration)
	_, startTime, ripeTime := calcTimeWithBuffs(buffs, true, ctime, orgRipeTime, ctime)

	if fc == nil {
		fc = entity.NewFarmCube(heroId, cube, startTime, ripeTime, resConf, m.datas.FarmMiscConfig())
		dbSucc = m.createFarmCube(fc)
		if !dbSucc {
			return farmMsg.ErrPlantFailServerErr
		}
	} else if entity.CanPlant(fc) {
		dbSucc = m.plantFarmCube(heroId, cube, startTime, ripeTime, resConf.Id)
		if !dbSucc {
			return farmMsg.ErrPlantFailServerErr
		}
	} else {
		logrus.Debugf("农场种植，发来的地块不是可种状态 heroId:%v cube:%v", heroId, cube)
		return farmMsg.ErrPlantFailNotIdleCube
	}

	// 删掉这块地的偷菜记录
	m.removeFarmSteal(heroId, []cb.Cube{fc.Cube})

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_PLANT_FARM) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_PLANT_FARM)))
		}
	})

	return
}

func (m *FarmFunc) OneKeyPlant(hc iface.HeroController, oneKeyConf *farm.FarmOneKeyConfig, allCubes, npcConflictCubes []cb.Cube, goldResConf, stoneResConf *farm.FarmResConfig, goldCount, stoneCount uint64) (errMsg msg.ErrMsg) {
	dbSucc, farmCubes := m.loadFarmCubes(hc.Id())
	if !dbSucc {
		return farmMsg.ErrOneKeyPlantFailServerErr
	}

	ctime := m.timeService.CurrentTime()
	//goldCount, stoneCount := m.getOneKeyPlantResConf(farmCubes, oneKeyConf)
	msgProto := &farmMsg.S2COneKeyPlantProto{}

	var buffs []*entity.BuffInfo
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		buffs = hero.Buff().Buffs(shared_proto.BuffEffectType_Buff_ET_farm_harvest)
		return
	})

	// 种植的地块
	plantedCubes := make([]cb.Cube, 0)
	ForEachCubes(allCubes, func(cube cb.Cube) {
		for _, b := range npcConflictCubes {
			if b == cube {
				return
			}
		}

		if farmCube := getCube(farmCubes, cube); farmCube != nil {
			if !entity.CanPlant(farmCube) {
				return
			}
		}

		var resConf *farm.FarmResConfig
		if goldCount > 0 {
			resConf = goldResConf
			goldCount--
		} else {
			resConf = stoneResConf
			stoneCount--
		}

		orgRipeTime := ctime.Add(resConf.RipeDuration)
		_, startTime, ripeTime := calcTimeWithBuffs(buffs, true, ctime, orgRipeTime, ctime)

		var farmCube *entity.FarmCube
		if farmCube = getCube(farmCubes, cube); farmCube == nil {
			farmCube = entity.NewFarmCube(hc.Id(), cube, startTime, ripeTime, resConf, m.datas.FarmMiscConfig())
			if !m.createFarmCube(farmCube) {
				logrus.Errorf("FarmFunc.OneKeyPlant CreateFarmCube 失败 heroId:%v cube:%v", hc.Id, farmCube.Cube)
				return
			}
		} else {
			if !entity.CanPlant(farmCube) {
				return
			}
			if !m.plantFarmCube(hc.Id(), cube, ctime, ripeTime, resConf.Id) {
				logrus.Errorf("FarmFunc.OneKeyPlant plantFarmCube 失败 heroId:%v cube:%v", hc.Id, farmCube.Cube)
				return
			}
		}

		plantedCubes = append(plantedCubes, farmCube.Cube)

		x, y := farmCube.Cube.XY()
		msgProto.CubeX = append(msgProto.CubeX, int32(x))
		msgProto.CubeY = append(msgProto.CubeY, int32(y))
		msgProto.ResId = append(msgProto.ResId, int32(resConf.Id))
	})

	if len(plantedCubes) <= 0 {
		hc.Send(farmMsg.ERR_ONE_KEY_PLANT_FAIL_NONE_IDLE_CUBE)
		return
	}

	hc.Send(farmMsg.NewS2cOneKeyPlantProtoMsg(msgProto))
	// 删掉这块地的偷菜记录
	m.removeFarmSteal(hc.Id(), plantedCubes)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_PLANT_FARM) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_PLANT_FARM)))
		}
	})

	return nil
}

func (m *FarmFunc) getOneKeyPlantResConf(currFarmCubes []*entity.FarmCube, okConf *farm.FarmOneKeyConfig) (goldCount, stoneCount uint64) {
	goldCount = okConf.GoldHopeCubeCount
	stoneCount = okConf.StoneHopeCubeCount
	ForEachFarmCubes(currFarmCubes, func(fc *entity.FarmCube) {
		if fc == nil || fc.ResConf == nil {
			return
		}
		if entity.IsRemoved(fc) {
			return
		}
		if fc.ResConf.ResType == shared_proto.ResType_GOLD {
			goldCount--
		} else if fc.ResConf.ResType == shared_proto.ResType_STONE {
			stoneCount--
		}
	})
	return
}

func (m *FarmFunc) Harvest(hc iface.HeroController, cube cb.Cube) (errMsg msg.ErrMsg, fc *entity.FarmCube) {
	dbSucc, fc := m.loadFarmCube(hc.Id(), cube)
	if !dbSucc {
		errMsg = farmMsg.ErrHarvestFailServerErr
		return
	}

	if fc == nil {
		errMsg = farmMsg.ErrHarvestFailInvalidCube
		return
	}

	ctime := m.timeService.CurrentTime()
	if !entity.CanHarvest(fc, ctime) {
		errMsg = farmMsg.ErrHarvestFailInvalidCube
		return
	}

	var buffs []*entity.BuffInfo
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		buffs = hero.Buff().Buffs(shared_proto.BuffEffectType_Buff_ET_farm_harvest)
		return
	})

	orgRipeTime := ctime.Add(fc.ResConf.RipeDuration)
	_, startTime, ripeTime := calcTimeWithBuffs(buffs, true, ctime, orgRipeTime, ctime)

	dbSucc = m.plantFarmCube(hc.Id(), cube, startTime, ripeTime, fc.ResConf.Id)
	if !dbSucc {
		errMsg = farmMsg.ErrHarvestFailServerErr
		return
	}

	// 删掉这块地的偷菜记录
	m.removeFarmSteal(hc.Id(), []cb.Cube{fc.Cube})

	return
}

func (m *FarmFunc) OneKeyReset(hc iface.HeroController) (errMsg msg.ErrMsg, harvestCubes []*entity.FarmCube) {
	dbSucc, farmCubes := m.loadFarmCubes(hc.Id())
	if !dbSucc {
		errMsg = farmMsg.ErrOneKeyResetFailServerErr
		return
	}

	ForEachFarmCubes(farmCubes, func(fc *entity.FarmCube) {
		if entity.IsIdle(fc) {
			return
		}

		if !entity.IsConflict(fc) {
			harvestCubes = append(harvestCubes, fc)
		}

		dbSucc := m.dbExec(hc.Id(), func(ctx context.Context) error {
			return m.db.ResetFarmCubes(ctx, hc.Id())
		})
		if !dbSucc {
			logrus.Errorf("农场一键犁地，重置地块失败 heroId:%v ", hc.Id())
			return
		}

		dbSucc = m.dbExec(hc.Id(), func(ctx context.Context) error {
			return m.db.ResetConflictFarmCubes(ctx, hc.Id())
		})
		if !dbSucc {
			logrus.Errorf("农场一键犁地，重置冲突地块失败 heroId:%v ", hc.Id())
			return
		}
	})

	return
}

func (m *FarmFunc) OneKeyHarvest(hc iface.HeroController, resType int32) (errMsg msg.ErrMsg, harvestCubes []*entity.FarmCube) {
	dbSucc, farmCubes := m.loadFarmHarvestCubes(hc.Id())
	if !dbSucc {
		errMsg = farmMsg.ErrOneKeyHarvestFailServerErr
		return
	}

	ctime := m.timeService.CurrentTime()

	var buffs []*entity.BuffInfo
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		buffs = hero.Buff().Buffs(shared_proto.BuffEffectType_Buff_ET_farm_harvest)
		return
	})

	// 重新种植的地块
	rePlantCubes := make([]cb.Cube, 0)

	ForEachFarmCubes(farmCubes, func(fc *entity.FarmCube) {
		if entity.IsIdle(fc) {
			return
		}
		if !entity.IsRipe(fc, ctime) && !entity.IsRemoved(fc) {
			return
		}
		if resType != 0 && int32(fc.ResType()) != resType {
			return
		}

		harvestCubes = append(harvestCubes, fc)

		if entity.IsRemoved(fc) {
			dbSucc = m.removeFarmCube(hc.Id(), fc.Cube)
			if !dbSucc {
				logrus.Errorf("农场一键收获，删除地块失败 heroId:%v cube:%v", hc.Id(), fc.Cube)
			}
			return
		}

		orgRipeTime := ctime.Add(fc.ResConf.RipeDuration)
		_, startTime, ripeTime := calcTimeWithBuffs(buffs, true, ctime, orgRipeTime, ctime)

		dbSucc := m.plantFarmCube(fc.HeroId, fc.Cube, startTime, ripeTime, fc.ResConf.Id)
		if !dbSucc {
			logrus.Errorf("农场一键收获，保存地块失败 heroId:%v cube:%v", hc.Id(), fc.Cube)
			return
		}

		rePlantCubes = append(rePlantCubes, fc.Cube)
	})

	if len(rePlantCubes) <= 0 {
		hc.Send(farmMsg.ERR_ONE_KEY_HARVEST_FAIL_NONE_IDLE_CUBE)
		return
	}

	// 删掉这块地的偷菜记录
	if len(rePlantCubes) > 0 {
		m.removeFarmSteal(hc.Id(), rePlantCubes)
	}

	return
}

func (m *FarmFunc) Change(hc iface.HeroController, cube cb.Cube, newResConf *farm.FarmResConfig) (errMsg msg.ErrMsg, resConf *farm.FarmResConfig, plantDuration time.Duration, stealTimes uint64) {
	dbSucc, fc := m.loadFarmCube(hc.Id(), cube)
	if !dbSucc {
		errMsg = farmMsg.ErrChangeFailServerErr
		return
	}

	if fc == nil {
		errMsg = farmMsg.ErrChangeFailInvalidCube
		return
	}

	ctime := m.timeService.CurrentTime()
	if !entity.CanHarvest(fc, ctime) {
		errMsg = farmMsg.ErrChangeFailInvalidCube
		return
	}

	var buffs []*entity.BuffInfo
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		buffs = hero.Buff().Buffs(shared_proto.BuffEffectType_Buff_ET_farm_harvest)
		return
	})

	resConf = fc.ResConf
	plantDuration = ctime.Sub(fc.StartTime)
	stealTimes = fc.StealTimes

	orgRipeTime := ctime.Add(newResConf.RipeDuration)
	_, startTime, ripeTime := calcTimeWithBuffs(buffs, true, ctime, orgRipeTime, ctime)

	dbSucc = m.plantFarmCube(hc.Id(), cube, startTime, ripeTime, newResConf.Id)
	if !dbSucc {
		errMsg = farmMsg.ErrChangeFailServerErr
		return
	}

	// 删掉这块地的偷菜记录
	m.removeFarmSteal(hc.Id(), []cb.Cube{fc.Cube})

	return
}

func (m *FarmFunc) Steal(hc iface.HeroController, targetId int64, cube cb.Cube, ctime time.Time, goldEffect, stoneEffect *data.Amount, canStealGold, canStealStone uint64) (errMsg msg.ErrMsg, goldOutput, stoneOutput uint64) {
	dbSucc, fc := m.loadFarmCube(targetId, cube)
	if !dbSucc {
		errMsg = farmMsg.ErrStealFailServerErr
		return
	}

	if fc == nil || !entity.CanSteal(fc, ctime) {
		errMsg = farmMsg.ErrStealFailNoStealTime
		return
	}

	dbSucc, stealCount := m.loadFarmStealCount(targetId, hc.Id(), cube)
	if !dbSucc {
		errMsg = farmMsg.ErrStealFailServerErr
		return
	}
	if stealCount > 0 {
		logrus.Debugf("农场偷菜，这块地已经偷过了")
		errMsg = farmMsg.ErrStealFailCubeAlreadyStealed
		return
	}

	// 计算
	resType, output := fc.CalcStealOutput(ctime, goldEffect, stoneEffect)
	if resType == shared_proto.ResType_GOLD {
		goldOutput = u64.Min(canStealGold, output)
	} else if resType == shared_proto.ResType_STONE {
		stoneOutput = u64.Min(canStealStone, output)
	}
	if goldOutput <= 0 && stoneOutput <= 0 {
		errMsg = farmMsg.ErrStealFailDailyStealFull
		return
	}

	// 添加偷菜数据
	dbSucc = m.addFarmSteal(targetId, hc.Id(), fc.Cube)
	if !dbSucc {
		errMsg = farmMsg.ErrStealFailServerErr
		return
	}

	// 更新地块数据
	dbSucc = m.updateFarmStealTimes(targetId, []cb.Cube{fc.Cube})
	if !dbSucc {
		errMsg = farmMsg.ErrStealFailServerErr
		return
	}

	x, y := cube.XY()
	hc.Send(farmMsg.NewS2cStealMsg(idbytes.ToBytes(targetId), int32(x), int32(y), u64.Int32(output)))

	// 被偷消息
	m.worldService.Send(targetId, farmMsg.NewS2cWhoStealFromMeMsg(idbytes.ToBytes(hc.Id()), u64.Int32(goldOutput), u64.Int32(stoneOutput)))

	return
}

func (m *FarmFunc) OneKeySteal(hc iface.HeroController, targetId int64, ctime time.Time, goldEffect, stoneEffect *data.Amount, canStealGold, canStealStone uint64) (errMsg msg.ErrMsg, goldOutput, stoneOutput uint64) {
	dbSucc, canStealCubes := m.loadCanStealCube(targetId, hc.Id(), ctime.Add(m.miscConf.NegRipeProtectDuration))
	if !dbSucc {
		errMsg = farmMsg.ErrHarvestFailServerErr
		return
	}

	whoStealMeMsgProto := &farmMsg.S2CWhoOneKeyStealFromMeProto{}
	stealMsgProto := &farmMsg.S2COneKeyStealProto{}

	stealedCubes := make([]cb.Cube, 0)

	// 真正被偷的地块
	ForEachFarmCubes(canStealCubes, func(fc *entity.FarmCube) {
		resType, output := fc.CalcStealOutput(ctime, goldEffect, stoneEffect)

		if resType == shared_proto.ResType_GOLD {
			if canStealGold <= 0 {
				return
			}

			output = u64.Min(canStealGold, output)
			goldOutput += output
			canStealGold = u64.Sub(canStealGold, output)

			stealMsgProto.CubeGoldOutput = append(stealMsgProto.CubeGoldOutput, int32(output))
			stealMsgProto.CubeStoneOutput = append(stealMsgProto.CubeStoneOutput, 0)
		} else if resType == shared_proto.ResType_STONE {
			if canStealStone <= 0 {
				return
			}

			output = u64.Min(canStealStone, output)
			stoneOutput += output
			canStealStone = u64.Sub(canStealStone, output)

			stealMsgProto.CubeGoldOutput = append(stealMsgProto.CubeGoldOutput, 0)
			stealMsgProto.CubeStoneOutput = append(stealMsgProto.CubeStoneOutput, int32(output))
		}

		// 添加偷菜数据
		dbSucc = m.addFarmSteal(targetId, hc.Id(), fc.Cube)
		if !dbSucc {
			return
		}

		x, y := fc.Cube.XY()
		whoStealMeMsgProto.CubeX = append(whoStealMeMsgProto.CubeX, int32(x))
		whoStealMeMsgProto.CubeY = append(whoStealMeMsgProto.CubeY, int32(y))
		whoStealMeMsgProto.StealTimes = append(whoStealMeMsgProto.StealTimes, u64.Int32(fc.StealTimes))

		stealMsgProto.CubeX = append(stealMsgProto.CubeX, int32(x))
		stealMsgProto.CubeY = append(stealMsgProto.CubeY, int32(y))

		stealedCubes = append(stealedCubes, fc.Cube)
	})

	if len(stealedCubes) > 0 {
		// 批量更新地块数据
		dbSucc = m.updateFarmStealTimes(targetId, stealedCubes)
		if !dbSucc {
			errMsg = farmMsg.ErrOneKeyStealFailServerErr
			return
		}
	} else {
		errMsg = farmMsg.ErrOneKeyStealFailNoCanStealCube
		return
	}

	stealMsgProto.GoldOutput = u64.Int32(goldOutput)
	stealMsgProto.StoneOutput = u64.Int32(stoneOutput)
	hc.Send(farmMsg.NewS2cOneKeyStealProtoMsg(stealMsgProto))

	// 被偷消息
	m.worldService.Send(targetId, farmMsg.NewS2cWhoOneKeyStealFromMeProtoMsg(whoStealMeMsgProto))
	return
}

func (m *FarmFunc) AddLog(heroId int64, heroName string, guildFlag string, countryId uint64, targetId int64, goldOutput, stoneOutput uint64, logTime time.Time, logType shared_proto.FarmLogType) {
	log := &shared_proto.FarmStealLogProto{}
	log.HeroId = idbytes.ToBytes(targetId)
	log.ThiefId = idbytes.ToBytes(heroId)
	log.ThiefName = heroName
	log.ThiefGuildFlag = guildFlag
	log.ThiefCountryId = u64.Int32(countryId)
	log.GoldOutput = u64.Int32(goldOutput)
	log.StoneOutput = u64.Int32(stoneOutput)
	log.LogTime = timeutil.Marshal32(logTime)
	log.Type = logType

	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = m.db.CreateFarmLog(ctx, log)
		return err
	})

	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc.AddLog err.")
		return
	}
}

func (m *FarmFunc) FarmLogList(heroId int64, size uint64, newest bool) (succMsg pbutil.Buffer, errMsg msg.ErrMsg) {
	var logs []*shared_proto.FarmStealLogProto
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		logs, err = m.db.LoadFarmLog(ctx, heroId, size)
		return err
	})
	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc.FarmLogList LoadFarmCubes err.")
		errMsg = farmMsg.ErrStealLogListFailServerErr
		return
	}

	proto := &shared_proto.FarmStealLogListProto{}
	proto.Log = logs
	protoBytes, err := proto.Marshal()
	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc.FarmLogList FarmStealLogListProto.Marshal() err.")
		errMsg = farmMsg.ErrStealLogListFailServerErr
		return
	}

	msgProto := &farmMsg.S2CStealLogListProto{}
	msgProto.Logs = protoBytes
	msgProto.Newest = newest

	succMsg = farmMsg.NewS2cStealLogListProtoMsg(msgProto).Static()
	return
}

func (m *FarmFunc) GetRelationCanStealList(hc iface.HeroController, relationIds []int64) (succMsg pbutil.Buffer, errMsg msg.ErrMsg) {
	msgProto := &farmMsg.S2CCanStealListProto{}
	thiefId := hc.Id()
	ctime := m.timeService.CurrentTime()

	var stealEnough bool
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		guanFuLevel := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU).Level
		if u64.FromInt32(hero.FarmExtra().GetDailyStealGold()) >= m.MaxDailyStealGoldAmount(guanFuLevel) &&
			u64.FromInt32(hero.FarmExtra().GetDailyStealStone()) >= m.MaxDailyStealStoneAmount(guanFuLevel) {
			stealEnough = true
		}
		return false
	})

	if stealEnough {
		succMsg = farmMsg.NewS2cCanStealListProtoMsg(msgProto).Static()
		return
	}

	ForEachHeroIds(relationIds, func(relationId int64) {
		dbSucc, canStealCount := m.loadCanStealCount(relationId, thiefId, ctime)
		if !dbSucc {
			return
		}

		if canStealCount <= 0 {
			return
		}

		msgProto.CanStealId = append(msgProto.CanStealId, idbytes.ToBytes(relationId))
	})

	succMsg = farmMsg.NewS2cCanStealListProtoMsg(msgProto).Static()
	return
}

func ForEachHeroIds(relationIds []int64, f func(relationId int64)) {
	for _, id := range relationIds {
		if id > 0 {
			f(id)
		}
	}
}

func ForEachFarmCubes(farmCubes []*entity.FarmCube, f func(fc *entity.FarmCube)) {
	for _, fc := range farmCubes {
		if fc != nil {
			f(fc)
		}
	}
}

func ForEachCubes(cubes []cb.Cube, f func(cube cb.Cube)) {
	for _, cube := range cubes {
		f(cube)
	}
}

func (m *FarmFunc) addFarmSteal(heroId, thiefId int64, cube cb.Cube) bool {
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = m.db.AddFarmSteal(ctx, heroId, thiefId, cube)
		return err
	})
	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc AddFarmSteal err. heroId:%v cube:%v", heroId, cube)
		return false
	}
	return true
}

func (m *FarmFunc) removeFarmSteal(heroId int64, cubes []cb.Cube) bool {
	// 删掉这块地的偷菜记录
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		return m.db.RemoveFarmSteal(ctx, heroId, cubes)
	})
	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc RemoveFarmSteal err. heroId:%v ", heroId)
		return false
	}
	return true
}

func (m *FarmFunc) loadFarmCube(heroId int64, cube cb.Cube) (succ bool, fc *entity.FarmCube) {
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		fc, err = m.db.LoadFarmCube(ctx, heroId, cube)
		return err
	})

	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc LoadFarmCube err. heroId:%v cube:%v", heroId, cube)
		return
	}
	succ = true
	return
}

func (m *FarmFunc) loadCanStealCube(heroId, thiefId int64, minRipeTime time.Time) (succ bool, cubes []*entity.FarmCube) {
	minRipeTimeInt := timeutil.Marshal64(minRipeTime)
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		cubes, err = m.db.LoadCanStealCube(ctx, heroId, thiefId, minRipeTimeInt, m.miscConf.CubeStealMaxTime)
		return err
	})
	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc loadCanStealCube err. heroId:%v", heroId)
		return
	}
	succ = true
	return
}

func (m *FarmFunc) loadFarmStealCubes(heroId int64) (succ bool, farmCubes []*entity.FarmCube) {
	ctime := m.timeService.CurrentTime()
	ripeTime := ctime.Add(m.miscConf.NegRipeProtectDuration)
	ripeTimeInt := timeutil.Marshal64(ripeTime)
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		farmCubes, err = m.db.LoadFarmStealCubes(ctx, heroId, ripeTimeInt, m.miscConf.CubeStealMaxTime)
		return err
	})
	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc LoadFarmCubes err. heroId:%v", heroId)
		return
	}
	succ = true
	return
}

func (m *FarmFunc) loadFarmCubes(heroId int64) (succ bool, farmCubes []*entity.FarmCube) {
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		farmCubes, err = m.db.LoadFarmCubes(ctx, heroId)
		return err
	})

	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc.LoadFarmCubes err. heroId:%v", heroId)
		return
	}
	succ = true
	return
}

func (m *FarmFunc) loadFarmHarvestCubes(heroId int64) (succ bool, farmCubes []*entity.FarmCube) {
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		farmCubes, err = m.db.LoadFarmHarvestCubes(ctx, heroId)
		return err
	})

	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc.OneKeyHarvest LoadFarmCubes err. heroId:%v", heroId)
		return
	}
	succ = true
	return
}

func (m *FarmFunc) removeFarmCube(heroId int64, cube cb.Cube) bool {
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = m.db.RemoveFarmCube(ctx, heroId, cube)
		return err
	})

	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc RemoveFarmCube err. heroId:%v cube:%v", heroId, cube)
		return false
	}
	return true
}

func (m *FarmFunc) dbExec(heroId int64, f func(ctx context.Context) error) bool {
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = f(ctx)
		return err
	})

	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc dbExec err. heroId:%v ", heroId)
		return false
	}
	return true
}

func (m *FarmFunc) plantFarmCube(heroId int64, cube cb.Cube, startTime, ripeTime time.Time, resId uint64) bool {
	startTimeInt := timeutil.Marshal64(startTime)
	ripeTimeInt := timeutil.Marshal64(ripeTime)
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = m.db.PlantFarmCube(ctx, heroId, cube, startTimeInt, ripeTimeInt, resId)
		return err
	})

	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc plantFarmCube err. heroId:%v ", heroId)
		return false
	}
	return true
}

func (m *FarmFunc) updateFarmStealTimes(heroId int64, cubes []cb.Cube) bool {
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = m.db.UpdateFarmStealTimes(ctx, heroId, cubes)
		return err
	})
	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc UpdateFarmStealTimes err. heroId:%v ", heroId)
		return false
	}
	return true
}

func (m *FarmFunc) createFarmCube(fc *entity.FarmCube) bool {
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = m.db.CreateFarmCube(ctx, fc)
		return err
	})
	if err != nil {
		logrus.WithError(err).Errorf("FarmFunc CreateFarmCube 失败 heroId:%v cube:%v", fc.HeroId, fc.Cube)
		return false
	}
	return true
}

func (m *FarmFunc) loadFarmStealCount(relationId, thiefId int64, cube cb.Cube) (succ bool, stealCount uint64) {
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		stealCount, err = m.db.LoadFarmStealCount(ctx, relationId, thiefId, cube)
		return
	})
	if err != nil {
		logrus.WithError(err).Errorf("FarmFunc loadFarmStealCount err")
		return
	}
	succ = true
	return
}

func (m *FarmFunc) loadCanStealCount(relationId, thiefId int64, ctime time.Time) (succ bool, count uint64) {
	minExpireProtcetTime := ctime.Add(m.miscConf.NegRipeProtectDuration)
	minExpreProtectTimeInt := timeutil.Marshal64(minExpireProtcetTime)
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		count, err = m.db.LoadCanStealCount(ctx, relationId, thiefId, minExpreProtectTimeInt, m.miscConf.CubeStealMaxTime)
		return
	})
	if err != nil {
		logrus.WithError(err).Errorf("FarmFunc loadFarmStealCount err")
		return
	}
	succ = true
	return
}

func getCube(fcs []*entity.FarmCube, cb cb.Cube) *entity.FarmCube {
	for _, fc := range fcs {
		if fc.Cube == cb {
			return fc
		}
	}
	return nil
}

func (m *FarmFunc) MaxDailyStealGoldAmount(guanfuLevel uint64) uint64 {
	conf := m.datas.GetFarmMaxStealConfig(guanfuLevel)
	if conf == nil {
		return 0
	}

	return conf.MaxDailyStealGoldAmount
}

func (m *FarmFunc) MaxDailyStealStoneAmount(guanfuLevel uint64) uint64 {
	conf := m.datas.GetFarmMaxStealConfig(guanfuLevel)
	if conf == nil {
		return 0
	}

	return conf.MaxDailyStealStoneAmount
}
