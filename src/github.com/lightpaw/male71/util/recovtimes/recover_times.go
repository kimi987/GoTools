package recovtimes

import (
	"time"
	"github.com/lightpaw/male7/util/timeutil"
)

func NewRecoverTimes(startTime time.Time, duration time.Duration, maxTimes uint64) *RecoverTimes {
	return (*RecoverTimes)(NewExtraRecoverTimes(startTime, duration, maxTimes))
}

type RecoverTimes ExtraRecoverTimes

func (rt *RecoverTimes) GetExtraRecoverTimes() *ExtraRecoverTimes {
	return rt.extra()
}

func (rt *RecoverTimes) extra() *ExtraRecoverTimes {
	return (*ExtraRecoverTimes)(rt)
}

// 开始恢复时间
func (rt *RecoverTimes) StartTime() time.Time {
	return rt.startTime
}

// 开始恢复时间
func (rt *RecoverTimes) StartTimeUnix64() int64 {
	return timeutil.Marshal64(rt.startTime)
}

// 开始恢复时间
func (rt *RecoverTimes) StartTimeUnix32() int32 {
	return timeutil.Marshal32(rt.startTime)
}

// 设置开始恢复时间
func (rt *RecoverTimes) SetStartTime(toSet time.Time) {
	rt.startTime = toSet
}

// 整数次
func (rt *RecoverTimes) Times(ctime time.Time) uint64 {
	return rt.extra().Times(ctime, 0)
}

func (rt *RecoverTimes) MaxTimes() uint64 {
	return rt.maxTimes
}

func (rt *RecoverTimes) Duration() time.Duration {
	return rt.duration
}

// 浮点次数
func (rt *RecoverTimes) floatTimes(ctime time.Time) float64 {
	return rt.extra().floatTimes(ctime, 0)
}

// 是否有足够的次数
func (rt *RecoverTimes) HasEnoughTimes(times uint64, ctime time.Time) bool {
	return times <= rt.Times(ctime)
}

// 减少一次
func (rt *RecoverTimes) ReduceOneTimes(ctime time.Time) {
	rt.ReduceTimes(1, ctime)
}

// 减少次数
func (rt *RecoverTimes) ReduceTimes(times uint64, ctime time.Time) {
	rt.extra().ReduceTimes(times, ctime, 0)
}

// 奖励次数
func (rt *RecoverTimes) AddTimes(toAdd uint64, ctime time.Time) {
	rt.extra().AddTimes(toAdd, ctime, 0)
}

// 奖励默认次数
func (rt *RecoverTimes) SetTimes(giveTimes uint64, ctime time.Time) *RecoverTimes {
	rt.extra().SetTimes(giveTimes, ctime, 0)
	return rt
}

// 变更恢复间隔
func (rt *RecoverTimes) ChangeDuration(newDuration time.Duration, ctime time.Time) {
	rt.extra().ChangeDuration(newDuration, ctime, 0)
}

// 变更恢复间隔跟最大次数
func (rt *RecoverTimes) ChangeMaxTimes(toSetMaxTimes uint64, ctime time.Time) {
	rt.extra().ChangeMaxTimes(toSetMaxTimes, ctime, 0)
}
