package entity

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/extratimesservice/extratimesface"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/util/recovtimes"
)

type RecoverableTimesWithExtraTimes interface {
	RecoverableTimes() *recovtimes.RecoverTimes
	UsedExtraTimes() UsedExtraTimes
	ReduceOneTimes(maxTimes extratimesface.ExtraMaxTimes, ctime time.Time)
	Times(maxTimes extratimesface.ExtraMaxTimes, ctime time.Time) uint64
	Encode() *shared_proto.RecoverableTimesWithExtraTimesProto
	Unmarshal(proto *shared_proto.RecoverableTimesWithExtraTimesProto)
}

func NewRecoverableTimesWithExtraTimes(startRecoveryTime time.Time, recoveryDuration time.Duration, maxTimes uint64) RecoverableTimesWithExtraTimes {
	return &recoverableTimesWithExtraTimes{
		recoverableTimes: recovtimes.NewRecoverTimes(startRecoveryTime, recoveryDuration, maxTimes),
		usedExtraTimes:   NewUsedExtraTimes(),
	}
}

// 可恢复的次数跟额外的次数
type recoverableTimesWithExtraTimes struct {
	recoverableTimes *recovtimes.RecoverTimes
	usedExtraTimes   UsedExtraTimes
}

func (r *recoverableTimesWithExtraTimes) Encode() *shared_proto.RecoverableTimesWithExtraTimesProto {
	proto := &shared_proto.RecoverableTimesWithExtraTimesProto{}

	proto.StartTime = r.recoverableTimes.StartTimeUnix32()
	proto.List = r.usedExtraTimes.Encode()

	return proto
}

func (r *recoverableTimesWithExtraTimes) Unmarshal(proto *shared_proto.RecoverableTimesWithExtraTimesProto) {
	if proto == nil {
		return
	}

	if proto.StartTime != 0 {
		r.recoverableTimes.SetStartTime(timeutil.Unix32(proto.StartTime))
	}

	for _, value := range proto.List.TypeList {
		r.usedExtraTimes.SetUsedExtraTimes(value.Type, u64.FromInt32(value.UsedTimes))
	}
}

func (r *recoverableTimesWithExtraTimes) RecoverableTimes() *recovtimes.RecoverTimes {
	return r.recoverableTimes
}

func (r *recoverableTimesWithExtraTimes) UsedExtraTimes() UsedExtraTimes {
	return r.usedExtraTimes
}

// 减少一次
func (r *recoverableTimesWithExtraTimes) ReduceOneTimes(maxTimes extratimesface.ExtraMaxTimes, ctime time.Time) {
	suc, _, _ := r.usedExtraTimes.ReduceOneTimes(maxTimes)
	if suc {
		return
	}

	r.recoverableTimes.ReduceOneTimes(ctime)
}

func (r *recoverableTimesWithExtraTimes) Times(maxTimes extratimesface.ExtraMaxTimes, ctime time.Time) uint64 {
	return r.usedExtraTimes.Times(maxTimes) + r.recoverableTimes.Times(ctime)
}

//type RecoverableTimes interface {
//	StartTime() time.Time
//	StartTimeUnix64() int64
//	StartTimeUnix32() int32
//	SetStartTime(toSet time.Time)
//	Times(ctime time.Time) uint64
//	MaxTimes() uint64
//	RecoveryDuration() time.Duration
//	HasEnoughTimes(times uint64, ctime time.Time) bool
//	ReduceOneTimes(ctime time.Time)
//	ReduceTimes(times uint64, ctime time.Time)
//	AddTimes(toAdd uint64, ctime time.Time)
//	SetTimes(giveTimes uint64, ctime time.Time) RecoverableTimes
//	ChangeRecoveryDuration(newRecoveryDuration time.Duration, ctime time.Time)
//	ChangeMaxTimes(newMaxTimes uint64, ctime time.Time)
//	floatTimes(ctime time.Time) float64
//}

func NewRecoverableTimes(startRecoveryTime time.Time, recoveryDuration time.Duration, maxTimes uint64) *recovtimes.RecoverTimes {
	return recovtimes.NewRecoverTimes(startRecoveryTime, recoveryDuration, maxTimes)
}

//// 可恢复的次数
//type recoverableTimes struct {
//	startRecoveryTime time.Time     // 开始恢复时间
//	recoveryDuration  time.Duration // 恢复间隔
//	maxTimes          uint64        // 最大恢复次数
//}
//
//// 开始恢复时间
//func (rt *recoverableTimes) StartRecoveryTime() time.Time {
//	return rt.startRecoveryTime
//}
//
//// 开始恢复时间
//func (rt *recoverableTimes) StartRecoveryTimeUnix64() int64 {
//	return timeutil.Marshal64(rt.startRecoveryTime)
//}
//
//// 开始恢复时间
//func (rt *recoverableTimes) StartRecoveryTimeUnix32() int32 {
//	return timeutil.Marshal32(rt.startRecoveryTime)
//}
//
//// 设置开始恢复时间
//func (rt *recoverableTimes) SetStartRecoveryTime(toSet time.Time) {
//	rt.startRecoveryTime = toSet
//}
//
//// 整数次
//func (rt *recoverableTimes) Times(ctime time.Time) uint64 {
//	times := ctime.Sub(rt.startRecoveryTime) / rt.recoveryDuration
//	if times <= 0 {
//		return 0
//	}
//	return u64.Min(rt.MaxTimes(), uint64(times))
//}
//
//func (rt *recoverableTimes) MaxTimes() uint64 {
//	return rt.maxTimes
//}
//
//func (rt *recoverableTimes) RecoveryDuration() time.Duration {
//	return rt.recoveryDuration
//}
//
//// 浮点次数
//func (rt *recoverableTimes) floatTimes(ctime time.Time) float64 {
//	times := float64(ctime.Sub(rt.startRecoveryTime)) / float64(rt.recoveryDuration)
//	if times <= 0 {
//		return 0
//	}
//	return math.Min(float64(rt.MaxTimes()), times)
//}
//
//// 是否有足够的次数
//func (rt *recoverableTimes) HasEnoughTimes(times uint64, ctime time.Time) bool {
//	return times <= rt.Times(ctime)
//}
//
//// 减少一次
//func (rt *recoverableTimes) ReduceOneTimes(ctime time.Time) {
//	rt.ReduceTimes(1, ctime)
//}
//
//// 减少次数
//func (rt *recoverableTimes) ReduceTimes(times uint64, ctime time.Time) {
//	curTimes := rt.Times(ctime)
//	leftTimes := u64.Sub(curTimes, times)
//	if curTimes == rt.MaxTimes() {
//		// 已经是最大次数了
//		rt.startRecoveryTime = ctime.Add(-time.Duration(leftTimes) * rt.recoveryDuration)
//	} else {
//		// 次数没满
//		rt.startRecoveryTime = rt.startRecoveryTime.Add(time.Duration(times) * rt.recoveryDuration)
//	}
//}
//
//// 奖励次数
//func (rt *recoverableTimes) AddTimes(toAdd uint64, ctime time.Time) {
//	if timeutil.IsZero(rt.startRecoveryTime) {
//		// 以前没给
//		toAdd = u64.Min(toAdd, rt.MaxTimes())
//		rt.startRecoveryTime = ctime.Add(-time.Duration(toAdd) * rt.recoveryDuration)
//	} else if rt.Times(ctime)+toAdd >= rt.MaxTimes() {
//		// 次数满了
//		rt.startRecoveryTime = ctime.Add(-time.Duration(rt.MaxTimes()) * rt.recoveryDuration)
//	} else {
//		// 次数没满
//		rt.startRecoveryTime = rt.startRecoveryTime.Add(-time.Duration(toAdd) * rt.recoveryDuration)
//	}
//}
//
//// 奖励默认次数
//func (rt *recoverableTimes) SetTimes(giveTimes uint64, ctime time.Time) RecoverableTimes {
//	giveTimes = u64.Min(giveTimes, rt.MaxTimes())
//	rt.startRecoveryTime = ctime.Add(-time.Duration(giveTimes) * rt.recoveryDuration)
//	return rt
//}
//
//// 变更恢复间隔
//func (rt *recoverableTimes) ChangeRecoveryDuration(newRecoveryDuration time.Duration, ctime time.Time) {
//	floatTimes := rt.floatTimes(ctime)
//
//	rt.recoveryDuration = newRecoveryDuration
//	rt.startRecoveryTime = ctime.Add(-time.Duration(floatTimes * float64(newRecoveryDuration)))
//}
//
//// 变更恢复间隔跟最大次数
//func (rt *recoverableTimes) ChangeMaxTimes(toSetMaxTimes uint64, ctime time.Time) {
//	floatTimes := rt.floatTimes(ctime)
//	rt.maxTimes = toSetMaxTimes
//
//	newMaxTimes := rt.MaxTimes()
//	if floatTimes > float64(newMaxTimes) {
//		// 要减少
//		rt.startRecoveryTime = ctime.Add(-time.Duration(newMaxTimes) * rt.recoveryDuration)
//	}
//}
