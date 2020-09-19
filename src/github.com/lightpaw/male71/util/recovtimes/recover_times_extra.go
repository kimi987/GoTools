package recovtimes

import (
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"math"
	"time"
)

func NewExtraRecoverTimes(startTime time.Time, duration time.Duration, maxTimes uint64) *ExtraRecoverTimes {
	return &ExtraRecoverTimes{
		startTime: startTime,
		duration:  duration,
		maxTimes:  maxTimes,
	}
}

// 额外次数恢复时间
type ExtraRecoverTimes struct {
	startTime time.Time     // 开始恢复时间
	duration  time.Duration // 恢复间隔
	maxTimes  uint64        // 最大恢复次数
}

// 开始恢复时间
func (rt *ExtraRecoverTimes) StartTime() time.Time {
	return rt.startTime
}

// 开始恢复时间
func (rt *ExtraRecoverTimes) StartTimeUnix64() int64 {
	return timeutil.Marshal64(rt.startTime)
}

// 开始恢复时间
func (rt *ExtraRecoverTimes) StartTimeUnix32() int32 {
	return timeutil.Marshal32(rt.startTime)
}

// 设置开始恢复时间
func (rt *ExtraRecoverTimes) SetStartTime(toSet time.Time) {
	rt.startTime = toSet
}

// 整数次
func (rt *ExtraRecoverTimes) Times(ctime time.Time, extraMaxTimes uint64) uint64 {
	times := ctime.Sub(rt.startTime) / rt.duration
	if times <= 0 {
		return 0
	}
	return u64.Min(rt.MaxTimes(extraMaxTimes), uint64(times))
}

func (rt *ExtraRecoverTimes) MaxTimes(extraMaxTimes uint64) uint64 {
	return rt.maxTimes + extraMaxTimes
}

func (rt *ExtraRecoverTimes) Duration() time.Duration {
	return rt.duration
}

// 浮点次数
func (rt *ExtraRecoverTimes) floatTimes(ctime time.Time, extraMaxTimes uint64) float64 {
	times := float64(ctime.Sub(rt.startTime)) / float64(rt.duration)
	if times <= 0 {
		return 0
	}
	return math.Min(float64(rt.MaxTimes(extraMaxTimes)), times)
}

// 是否有足够的次数
func (rt *ExtraRecoverTimes) HasEnoughTimes(times uint64, ctime time.Time, extraMaxTimes uint64) bool {
	return times <= rt.Times(ctime, extraMaxTimes)
}

// 减少一次
func (rt *ExtraRecoverTimes) ReduceOneTimes(ctime time.Time, extraMaxTimes uint64) {
	rt.ReduceTimes(1, ctime, extraMaxTimes)
}

// 减少次数
func (rt *ExtraRecoverTimes) ReduceTimes(times uint64, ctime time.Time, extraMaxTimes uint64) {
	curTimes := rt.Times(ctime, extraMaxTimes)
	leftTimes := u64.Sub(curTimes, times)
	if curTimes == rt.MaxTimes(extraMaxTimes) {
		// 已经是最大次数了
		rt.startTime = ctime.Add(-time.Duration(leftTimes) * rt.duration)
	} else {
		// 次数没满
		rt.startTime = rt.startTime.Add(time.Duration(times) * rt.duration)
	}
}

// 奖励次数
func (rt *ExtraRecoverTimes) AddTimes(toAdd uint64, ctime time.Time, extraMaxTimes uint64) {
	if timeutil.IsZero(rt.startTime) {
		// 以前没给
		toAdd = u64.Min(toAdd, rt.MaxTimes(extraMaxTimes))
		rt.startTime = ctime.Add(-time.Duration(toAdd) * rt.duration)
	} else if rt.Times(ctime, extraMaxTimes)+toAdd >= rt.MaxTimes(extraMaxTimes) {
		// 次数满了
		rt.startTime = ctime.Add(-time.Duration(rt.MaxTimes(extraMaxTimes)) * rt.duration)
	} else {
		// 次数没满
		rt.startTime = rt.startTime.Add(-time.Duration(toAdd) * rt.duration)
	}
}

// 奖励默认次数
func (rt *ExtraRecoverTimes) SetTimes(giveTimes uint64, ctime time.Time, extraMaxTimes uint64) *ExtraRecoverTimes {
	giveTimes = u64.Min(giveTimes, rt.MaxTimes(extraMaxTimes))
	rt.startTime = ctime.Add(-time.Duration(giveTimes) * rt.duration)
	return rt
}

// 变更恢复间隔
func (rt *ExtraRecoverTimes) ChangeDuration(newDuration time.Duration, ctime time.Time, extraMaxTimes uint64) {
	floatTimes := rt.floatTimes(ctime, extraMaxTimes)

	rt.duration = newDuration
	rt.startTime = ctime.Add(-time.Duration(floatTimes * float64(newDuration)))
}

// 变更恢复间隔跟最大次数
func (rt *ExtraRecoverTimes) ChangeMaxTimes(toSetMaxTimes uint64, ctime time.Time, extraMaxTimes uint64) {
	floatTimes := rt.floatTimes(ctime, extraMaxTimes)
	rt.maxTimes = toSetMaxTimes

	newMaxTimes := rt.MaxTimes(extraMaxTimes)
	if floatTimes > float64(newMaxTimes) {
		// 要减少
		rt.startTime = ctime.Add(-time.Duration(newMaxTimes) * rt.duration)
	}
}
