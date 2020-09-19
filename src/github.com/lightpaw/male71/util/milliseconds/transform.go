package milliseconds

import (
	"time"
)

const (
	PerSecond int64 = 1000
	PerMinute       = 60 * PerSecond
	PerHour         = 60 * PerMinute
	PerDay          = 24 * PerHour
)

const (
	FloatPerSecond float64 = float64(PerSecond)
	FloatPerMinute float64 = float64(PerMinute)
	FloatPerHour   float64 = float64(PerHour)
	FloatPerDay    float64 = float64(PerDay)
)

func Time(time time.Time) int64 {
	return time.UnixNano() / 1e6
}

func Duration(d time.Duration) int64 {
	return d.Nanoseconds() / 1e6
}

func FromSecond(d float64) int64 {
	return int64(d * FloatPerSecond)
}

func FromMinute(d float64) int64 {
	return int64(d * FloatPerMinute)
}

func FromHour(d float64) int64 {
	return int64(d * FloatPerHour)
}

func FromDay(d float64) int64 {
	return int64(d * FloatPerDay)
}
