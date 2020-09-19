package farm

import (
	"github.com/lightpaw/male7/gen/iface"
	farmMsg "github.com/lightpaw/male7/gen/pb/farm"
	"github.com/lightpaw/male7/config/farm"
	"github.com/lightpaw/logrus"
	"time"
	"github.com/lightpaw/male7/util/event"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
	"context"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
)

func NewFarmService(dbService iface.DbService, worldService iface.WorldService, datas iface.ConfigDatas, timeService iface.TimeService, tickerService iface.TickerService) *FarmService {
	m := FarmService{}
	m.datas = datas
	m.worldService = worldService
	m.timeService = timeService
	m.tickerService = tickerService
	m.miscConf = datas.FarmMiscConfig()
	m.db = dbService

	m.queue = event.NewEventQueue(2048, 5*time.Second, "FarmEvent")

	// db 中删除过期记录
	expiredTime := m.timeService.CurrentTime().Add(m.miscConf.StealLogExpiredDuration)
	ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		return m.db.RemoveFarmLog(ctx, timeutil.Marshal32(expiredTime))
	})

	return &m
}

//gogen:iface
type FarmService struct {
	datas               iface.ConfigDatas
	worldService        iface.WorldService
	guildService        iface.GuildService
	timeService         iface.TimeService
	tickerService       iface.TickerService
	heroDataService     iface.HeroDataService
	heroSnapshotService iface.HeroSnapshotService
	miscConf            *farm.FarmMiscConfig
	db                  iface.DbService
	tick                iface.TickerService

	queue *event.EventQueue
}

func (m *FarmService) FuncNoWait(f entity.FarmFuncType) {
	m.doFunc(false, f)
}

func (m *FarmService) FuncWait(handlerName string, serverErrMsg pbutil.Buffer, hc iface.HeroController, f entity.FarmFuncType) {
	if !m.doFunc(true, f) {
		logrus.Debugf("%s，超时", handlerName)
		hc.Send(serverErrMsg)
		return
	}
}

func (m *FarmService) doFunc(waitResult bool, f entity.FarmFuncType) bool {
	startTime := time.Now()
	ok := m.queue.TimeoutFunc(waitResult, f)
	logrus.Debugf("farm queue exec time: %vms", u64.Division2Float64(u64.FromInt64(time.Now().Sub(startTime).Nanoseconds()), 1000000))
	return ok
}

func (m *FarmService) Close() {
	m.queue.Stop()
}

/*
 更新农场地块
 absCubes 所有本次变化的地块
 allConflictedBlocks 本次变化的地块中，有冲突的地块
 npcConflictOffsets 所有 npc 造成的冲突地块
*/
func (m *FarmService) UpdateFarmCubes(heroId int64, baseLevel uint64, absCubes, allConflictedBlocks []cb.Cube, npcConflictOffsets []cb.Cube, baseX int, baseY int, ctime time.Time) {
	farmCubeMap := m.loadFarmCubeMap(heroId)

	// 空闲的受影响地块
	newFarmCubes := make([]*entity.FarmCube, 0)

	// 所有变化的，和自己有关的地块
	ForEachCubes(absCubes, func(absCube cb.Cube) {
		// 转成偏移
		x, y := absCube.XY()
		offset := hexagon.EvenOffsetBetween(baseX, baseY, x, y)
		allCubes := m.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(baseLevel)
		// 删除的
		isRemove := !cb.Contains(allCubes, offset)
		// 有冲突的
		conflicted := isConflicted(allConflictedBlocks, absCube) || isConflicted(npcConflictOffsets, offset)

		farmCube := farmCubeMap[offset]
		if farmCube == nil {
			farmCube = entity.NewIdleFarmCube(heroId, offset)
		}

		pauseTime := farmCube.GetPauseTime()

		// 更新状态
		if conflicted {
			if timeutil.IsZero(farmCube.ConflictTime) {
				// 原来不冲突
				farmCube.ConflictTime = ctime
			}
			if entity.IsIdle(farmCube) {
				newFarmCubes = append(newFarmCubes, farmCube)
			}
		} else {
			farmCube.ConflictTime = time.Time{}
		}
		if isRemove {
			if timeutil.IsZero(farmCube.RemoveTime) {
				farmCube.RemoveTime = ctime
			}
		} else {
			farmCube.RemoveTime = time.Time{}
		}

		if !conflicted && !isRemove {
			// 结束暂停状态
			if !entity.IsIdle(farmCube) && !timeutil.IsZero(pauseTime) {
				needDuration := farmCube.RipeTime.Sub(pauseTime)
				farmCube.StartTime = farmCube.StartTime.Add(needDuration)
				farmCube.RipeTime = farmCube.RipeTime.Add(needDuration)
			}
		}
	})

	// 更新 db
	for _, farmCube := range farmCubeMap {
		succ := m.updateFarmCubeState(heroId, farmCube)
		if !succ {
			logrus.Debugf("UpdateFarmCubes SaveFarmCube err. heroId:%v cube:%v", heroId, farmCube.Cube)
		}
	}

	for _, farmCube := range newFarmCubes {
		ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			err = m.db.CreateFarmCube(ctx, farmCube)
			if err != nil {
				logrus.WithError(err).Debugf("UpdateFarmCubes CreateFarmCube err. heroId:%v cube:%v", heroId, farmCube.Cube)
			}
			return
		})

	}

	// 同步更新农场通知
	m.worldService.Send(heroId, farmMsg.FARM_IS_UPDATE_S2C)
}

// 更新农场地块
func (m *FarmService) UpdateFarmCubeWithOffset(heroId int64, baseLevel uint64, baseX, baseY int, npcOffsets []cb.Cube, isAdd bool, ctime time.Time) {
	farmCubeMap := m.loadFarmCubeMap(heroId)
	allCubes := m.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(baseLevel)

	newFarmCubes := make([]*entity.FarmCube, 0)
	// 处理 npc 冲突
	ForEachCubes(npcOffsets, func(offset cb.Cube) {
		if !cb.Contains(allCubes, offset) {
			return
		}

		farmCube := farmCubeMap[offset]
		if farmCube == nil {
			farmCube = entity.NewIdleFarmCube(heroId, offset)
			newFarmCubes = append(newFarmCubes, farmCube)
		}
		if isAdd {
			if entity.IsConflict(farmCube) {
				return
			}
			// 原来不冲突
			farmCube.ConflictTime = ctime
		} else {
			if !entity.IsConflict(farmCube) {
				return
			}
			farmCube.ConflictTime = time.Time{}
		}
	})

	// 更新 db
	for _, farmCube := range farmCubeMap {
		succ := m.updateFarmCubeState(heroId, farmCube)
		if !succ {
			logrus.Debugf("UpdateFarmCubeWithOffset SaveFarmCube err. heroId:%v cube:%v", heroId, farmCube.Cube)
		}
	}

	for _, farmCube := range newFarmCubes {
		ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			err = m.db.CreateFarmCube(ctx, farmCube)
			if err != nil {
				logrus.WithError(err).Debugf("UpdateFarmCubeWithOffset CreateFarmCube err. heroId:%v cube:%v", heroId, farmCube.Cube)
			}
			return
		})
	}

	// 同步更新农场通知
	m.worldService.Send(heroId, farmMsg.FARM_IS_UPDATE_S2C)
}

func (m *FarmService) updateFarmCubeState(heroId int64, fc *entity.FarmCube) bool {
	conflictedTimeInt := timeutil.Marshal64(fc.ConflictTime)
	removeTimeInt := timeutil.Marshal64(fc.RemoveTime)
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = m.db.UpdateFarmCubeState(ctx, heroId, fc.Cube, conflictedTimeInt, removeTimeInt)
		return err
	})
	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc updateFarmCubeConflict err. heroId:%v", heroId)
		return false
	}
	return true
}

func (m *FarmService) loadFarmCubeMap(heroId int64) map[cb.Cube]*entity.FarmCube {
	farmCubeMap := make(map[cb.Cube]*entity.FarmCube)
	farmCubes := make([]*entity.FarmCube, 0)
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		farmCubes, err = m.db.LoadFarmCubes(ctx, heroId)
		return err
	})
	if err != nil {
		logrus.WithError(err).Debugf("FarmFunc.LoadFarmCubeMap err. heroId:%v", heroId)
		return farmCubeMap
	}

	for _, farmCube := range farmCubes {
		if farmCube != nil {
			farmCubeMap[farmCube.Cube] = farmCube
		}
	}
	return farmCubeMap
}

// 代码生成工具不支持参数为struct{}类型
func isConflicted(allConflictedBlocks []cb.Cube, block cb.Cube) bool {
	for _, b := range allConflictedBlocks {
		if b == block {
			return true
		}
	}
	return false
}

func (m *FarmService) ReduceRipeTime(heroId int64, toReduce time.Duration) {
	m.doFunc(true, func() {
		var cubes []*entity.FarmCube
		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			cubes, err = m.db.LoadFarmCubes(ctx, heroId)
			return
		}); err != nil {
			logrus.WithError(err).Errorf("ReduceRipeTime err.heroId:%v", heroId)
			return
		}

		for _, c := range cubes {
			if c == nil || entity.IsIdle(c) || entity.IsRemoved(c) {
				continue
			}

			if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
				ripeTime := c.RipeTime.Add(-toReduce)
				startTime := c.StartTime.Add(-toReduce)
				err = m.db.UpdateFarmCubeRipeTime(ctx, heroId, c.Cube, timeutil.Marshal64(startTime), timeutil.Marshal64(ripeTime))
				return
			}); err != nil {
				logrus.WithError(err).Debugf("FarmService.GMRipe err. heroId:%v", heroId)
			}
		}
	})
}

func (m *FarmService) ReduceRipeTimePercent(heroId int64, toReduceBuff, toAddBuff *entity.BuffInfo) {
	m.doFunc(true, func() {
		m.reduceRipeTimePercent(heroId, toReduceBuff, toAddBuff)
	})
}

func (m *FarmService) reduceRipeTimePercent(heroId int64, toReduceBuff, toAddBuff *entity.BuffInfo) {
	var cubes []*entity.FarmCube
	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		cubes, err = m.db.LoadFarmCubes(ctx, heroId)
		return
	}); err != nil {
		logrus.WithError(err).Errorf("ReduceRipeTime err.heroId:%v", heroId)
		return
	}

	ctime := m.timeService.CurrentTime()
	for _, c := range cubes {
		if c == nil || entity.IsIdle(c) || entity.IsRemoved(c) {
			continue
		}

		reduceSucc, startTime, ripeTime := calcTimeWithBuff(toReduceBuff, false, c.StartTime, c.RipeTime, ctime)
		toAddSucc, startTime, ripeTime := calcTimeWithBuff(toAddBuff, true, startTime, ripeTime, ctime)
		if !reduceSucc && !toAddSucc {
			continue
		}

		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			err = m.db.UpdateFarmCubeRipeTime(ctx, heroId, c.Cube, timeutil.Marshal64(startTime), timeutil.Marshal64(ripeTime))
			return
		}); err != nil {
			logrus.WithError(err).Debugf("FarmService.ReduceRipeTimePercent err. heroId:%v", heroId)
		}
	}
}

func (m *FarmService) GMRipe(heroId int64) {
	m.doFunc(true, func() {
		startTime := m.timeService.CurrentTime().Add(-72 * time.Hour)
		ripeTime := m.timeService.CurrentTime()
		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			err = m.db.GMFarmRipe(ctx, heroId, timeutil.Marshal64(startTime), timeutil.Marshal64(ripeTime))
			return err
		}); err != nil {
			logrus.WithError(err).Debugf("FarmService.GMRipe err. heroId:%v", heroId)
			return
		}
	})
}

func (m *FarmService) GMCanSteal(heroId int64) {
	m.doFunc(true, func() {
		startTime := m.timeService.CurrentTime().Add(-72 * time.Hour)
		canStealTime := m.timeService.CurrentTime().Add(-m.datas.FarmMiscConfig().RipeProtectDuration)
		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			err = m.db.GMFarmRipe(ctx, heroId, timeutil.Marshal64(startTime), timeutil.Marshal64(canStealTime))
			return err
		}); err != nil {
			logrus.WithError(err).Debugf("FarmService.GMRipe err. heroId:%v", heroId)
			return
		}
	})
}

func calcTimeWithBuffs(buffs []*entity.BuffInfo, buffIsAdd bool, startTime, ripeTime, ctime time.Time) (succ bool, newStartTime, newRipeTime time.Time) {
	newStartTime, newRipeTime = startTime, ripeTime
	var success bool
	for _, b := range buffs {
		success, newStartTime, newRipeTime = calcTimeWithBuff(b, buffIsAdd, newStartTime, newRipeTime, ctime)
		succ = succ || success
	}
	return
}

func calcTimeWithBuff(buff *entity.BuffInfo, buffIsAdd bool, startTime, ripeTime, ctime time.Time) (succ bool, newStartTime, newRipeTime time.Time) {
	newStartTime, newRipeTime = startTime, ripeTime
	if buff == nil {
		return
	}

	buffEndTime := timeutil.Min(buff.EndTime, ripeTime)
	buffStartTime := timeutil.Max(buff.StartTime, ctime)
	dur := buffEndTime.Sub(buffStartTime)
	if dur <= 0 {
		return
	}

	var realDur time.Duration
	if buffIsAdd {
		realDur = time.Duration(buff.EffectData.FarmHarvest.CalculateByPercent(uint64(dur)))
	} else {
		realDur = time.Duration(buff.EffectData.FarmHarvest.Return(uint64(dur)))
	}

	newStartTime = startTime.Add(realDur - dur)
	newRipeTime = ripeTime.Add(realDur - dur)
	succ = true
	return
}
