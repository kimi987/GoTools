package entity

import (
	"github.com/lightpaw/male7/entity/cb"
	"time"
	"github.com/lightpaw/male7/config/farm"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/config"
)

func NewIdleFarmCube(heroId int64, cube cb.Cube) *FarmCube {
	fc := &FarmCube{}
	fc.HeroId = heroId
	fc.Cube = cube
	fc.StartTime = time.Time{}
	fc.RipeTime = time.Time{}
	fc.RemoveTime = time.Time{}
	fc.ConflictTime = time.Time{}

	return fc
}

func NewFarmCube(heroId int64, cube cb.Cube, startTime, ripeTime time.Time, resConf *farm.FarmResConfig, miscConf *farm.FarmMiscConfig) *FarmCube {
	fc := &FarmCube{}
	fc.HeroId = heroId
	fc.Cube = cube
	fc.StartTime = startTime
	fc.RipeTime = ripeTime
	fc.RemoveTime = time.Time{}
	fc.ConflictTime = time.Time{}
	fc.ResConf = resConf
	fc.ResId = resConf.Id
	fc.MiscConf = miscConf

	return fc
}

func (fc *FarmCube) PrepareSave() {
	fc.StartTimeInt = timeutil.Marshal64(fc.StartTime)
	fc.RemoveTimeInt = timeutil.Marshal64(fc.RemoveTime)
	fc.ConflictTimeInt = timeutil.Marshal64(fc.ConflictTime)
	fc.RipeTimeInt = timeutil.Marshal64(fc.RipeTime)
	if fc.ResConf != nil {
		fc.ResId = fc.ResConf.Id
	}
}

func (fc *FarmCube) OnLoad(datas *config.ConfigDatas) {
	fc.StartTime = timeutil.Unix64(fc.StartTimeInt)
	fc.RemoveTime = timeutil.Unix64(fc.RemoveTimeInt)
	fc.ConflictTime = timeutil.Unix64(fc.ConflictTimeInt)
	fc.RipeTime = timeutil.Unix64(fc.RipeTimeInt)
	fc.ResConf = datas.GetFarmResConfig(fc.ResId)
	fc.MiscConf = datas.FarmMiscConfig()
}

type FarmFuncType func()

/*
 空闲、种植、成熟状态互斥
 上面叠加冲突和删除状态
 */
type FarmCube struct {
	HeroId          int64
	Cube            cb.Cube
	StartTimeInt    int64
	StartTime       time.Time
	RipeTimeInt     int64
	RipeTime        time.Time
	RemoveTimeInt   int64
	RemoveTime      time.Time
	ConflictTimeInt int64
	ConflictTime    time.Time
	ResId           uint64
	ResConf         *farm.FarmResConfig
	MiscConf        *farm.FarmMiscConfig
	StealTimes      uint64
	ConflictHeroId  int64

	ShowUnlockAnimation bool
}

func (fc *FarmCube) ResType() shared_proto.ResType {
	if fc.ResConf == nil {
		return shared_proto.ResType_InvalidResType
	}
	return fc.ResConf.ResType
}

/*
 空闲、种植、成熟状态互斥
 上面叠加冲突和删除状态
 */
func IsIdle(fc *FarmCube) bool {
	if fc.ResConf == nil {
		return true
	}
	return timeutil.IsZero(fc.StartTime) && timeutil.IsZero(fc.RipeTime)
}

func (fc *FarmCube) GetPauseTime() (pauseTime time.Time) {
	if timeutil.IsZero(fc.RemoveTime) {
		if timeutil.IsZero(fc.ConflictTime) {
			pauseTime = time.Time{}
		} else {
			pauseTime = fc.ConflictTime
		}
	} else {
		if timeutil.IsZero(fc.ConflictTime) {
			pauseTime = fc.RemoveTime
		} else {
			pauseTime = timeutil.Min(fc.RemoveTime, fc.ConflictTime)
		}
	}
	return
}

/*
 空闲、种植、成熟状态互斥
 上面叠加冲突和删除状态
 */
func IsPlanting(fc *FarmCube, ctime time.Time) bool {
	if fc.ResConf == nil {
		return false
	}
	if fc.RipeTime.Before(ctime) {
		return false
	}
	return true
}

/*
 空闲、种植、成熟状态互斥
 上面叠加冲突和删除状态
 */
func IsRipe(fc *FarmCube, ctime time.Time) bool {
	if fc.ResConf == nil {
		return false
	}
	if fc.RipeTime.After(ctime) {
		return false
	}
	return true
}

/*
 空闲、种植、成熟状态互斥
 上面叠加冲突和删除状态
 */
func IsRemoved(fc *FarmCube) bool {
	return !timeutil.IsZero(fc.RemoveTime)
}

func IsConflict(fc *FarmCube) bool {
	return !timeutil.IsZero(fc.ConflictTime)
}

func CanPlant(fc *FarmCube) bool {
	if IsConflict(fc) || IsRemoved(fc) {
		return false
	}
	return IsIdle(fc)
}

func CanHarvest(fc *FarmCube, ctime time.Time) bool {
	if IsIdle(fc) {
		return false
	}
	if IsConflict(fc) {
		return false
	}
	return IsPlanting(fc, ctime) || IsRipe(fc, ctime)
}

func CanSteal(fc *FarmCube, ctime time.Time) bool {
	if IsIdle(fc) {
		return false
	}
	if IsConflict(fc) {
		return false
	}
	if ctime.Before(fc.RipeTime.Add(fc.MiscConf.RipeProtectDuration)) {
		return false
	}
	return fc.StealTimes < fc.MiscConf.CubeStealMaxTime
}

func CalcMaxOutput(farmEffect *data.Amount, resConf *farm.FarmResConfig) uint64 {
	return resConf.BaseOutput.Amount + farmEffect.CalculateByPercent(resConf.BaseOutput.Amount)
}

func CalcHarvestOutput(plantDuration time.Duration, stealTimes uint64, farmEffect *data.Amount, resConf *farm.FarmResConfig, miscConf *farm.FarmMiscConfig) (output uint64) {
	isRipe := plantDuration >= resConf.RipeDuration

	if isRipe {
		maxOutput := CalcMaxOutput(farmEffect, resConf)
		stealOutput := miscConf.StealGainPercent.CalculateByPercent(maxOutput) * stealTimes
		output = u64.Sub(maxOutput, stealOutput)
	} else {
		defaultMaxOutput := CalcMaxOutput(farmEffect, resConf)
		mulit := u64.Division2Float64(uint64(plantDuration), uint64(resConf.RipeDuration))
		maxOutput := u64.MultiCoef(defaultMaxOutput, mulit)
		output = miscConf.EarlyHarvestPercent.CalculateByPercent(maxOutput)
	}
	return
}

func (fc *FarmCube) CalcStealOutput(ctime time.Time, goldEffect, stoneEffect *data.Amount) (resType shared_proto.ResType, output uint64) {
	if !CanSteal(fc, ctime) {
		return
	}
	resType = fc.ResType()
	var effect *data.Amount
	switch resType {
	case shared_proto.ResType_GOLD:
		effect = goldEffect
	case shared_proto.ResType_STONE:
		effect = stoneEffect
	default:
		effect = &data.Amount{}
	}
	maxOutput := CalcMaxOutput(effect, fc.ResConf)
	output = fc.MiscConf.StealGainPercent.CalculateByPercent(maxOutput)
	return
}

func (fc *FarmCube) Encode() *shared_proto.FarmCubeProto {
	proto := &shared_proto.FarmCubeProto{}
	x, y := fc.Cube.XY()
	proto.CubeX = int32(x)
	proto.CubeY = int32(y)
	proto.StartTime = timeutil.Marshal32(fc.StartTime)
	proto.ConflictTime = timeutil.Marshal32(fc.ConflictTime)
	proto.RemoveTime = timeutil.Marshal32(fc.RemoveTime)
	proto.ResConfig = u64.Int32(fc.ResId)
	proto.StealTimes = u64.Int32(fc.StealTimes)
	proto.ShowAnimation = fc.ShowUnlockAnimation

	return proto
}
