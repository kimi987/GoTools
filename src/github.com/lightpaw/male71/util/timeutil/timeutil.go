package timeutil

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/imath"
	"math"
	"runtime/debug"
	"time"
)

func MarshalArray64(array []time.Time) []int64 {
	out := make([]int64, len(array))
	for i, v := range array {
		out[i] = Marshal64(v)
	}

	return out
}

func MarshalArray32(array []time.Time) []int32 {
	out := make([]int32, len(array))
	for i, v := range array {
		out[i] = Marshal32(v)
	}

	return out
}

func CopyUnix32Array(array []time.Time, intArray []int32) {
	n := imath.Min(len(array), len(intArray))
	if n > 0 {
		for i := 0; i < n; i++ {
			array[i] = Unix32(intArray[i])
		}
	}
}

func Marshal64(t time.Time) int64 {
	if t.IsZero() {
		// 为了防止当前时间就是 unixZeroTime, 这里不用 IsZero(t)
		return 0
	}

	if unix := t.Unix(); unix < math.MinInt32 {
		logrus.WithField("stack", string(debug.Stack())).WithField("unix", unix).Error("timeutil.Marshal t.unix < math.MinInt32")
		return 0
	} else {
		return unix
	}
}

func Marshal32(t time.Time) int32 {
	return int32(Marshal64(t))
}

func Unix32(second int32) time.Time {
	return Unix64(int64(second))
}

func Unix64(second int64) time.Time {
	return time.Unix(second, 0)
}

const Day = 24 * time.Hour

func IsSameDay(t1, t2 time.Time) bool {
	ns := i64.Abs(t1.Sub(t2).Nanoseconds())
	if ns > Day.Nanoseconds() {
		return false
	}

	return DailyTime.PrevTime(t1).Equal(DailyTime.PrevTime(t2))
}

func DivideTimes(x, y time.Duration) uint64 {
	if x <= 0 || y <= 0 {
		return 0
	}

	return uint64((x + y - 1) / y)
}

func DurationMarshal32(duration time.Duration) int32 {
	return config.Duration2I32Seconds(duration)
}

func DurationMarshal64(duration time.Duration) int64 {
	return int64(duration / time.Second)
}

// 将duration数组转成seconds数组
func DurationArrayToSecondArray(array []time.Duration) (result []int32) {
	result = make([]int32, 0, len(array))

	for _, duration := range array {
		result = append(result, int32(duration/time.Second))
	}

	return result
}

func Duration32(seconds int32) time.Duration {
	return time.Duration(seconds) * time.Second
}

func Duration64(seconds int64) time.Duration {
	return time.Duration(seconds) * time.Second
}

var unixZeroTime = Unix64(0)

func IsZero(time time.Time) bool {
	return time.IsZero() || time.Equal(unixZeroTime)
}

func Min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func Max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func MinDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func MaxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

func MultiDuration(multi float64, d time.Duration) time.Duration {
	return time.Duration(multi * float64(d))
}

func NextTickTime(prevTime, ctime time.Time, d time.Duration) time.Time {
	if ctime.Before(prevTime) {
		return prevTime
	}
	return prevTime.Add((ctime.Sub(prevTime)/d)*d + d)
}

func Between(t, start, end time.Time) bool {
	return t.After(start) && t.Before(end)
}

func Rate(startTime, endTime, ctime time.Time) float64 {
	return i64.Rate(startTime.Unix(), endTime.Unix(), ctime.Unix())
}

func Midnight(t time.Time) time.Time {
	return DailyTime.PrevTime(t)
}