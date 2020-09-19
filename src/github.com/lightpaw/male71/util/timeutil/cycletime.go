package timeutil

import "time"

func NewDailyCycleTime(startTime int64) *CycleTime {
	return NewCycleTime(startTime, SecondsPerDay)
}

func NewWeeklyCycleTime(startTime int64) *CycleTime {
	return NewCycleTime(startTime, 7*SecondsPerDay)
}

func NewCycleTime(startTime, period int64) *CycleTime {
	return &CycleTime{
		startTime: startTime,
		period:    period,
	}
}

func NewOffsetDailyTime(offsetSeconds int64) *CycleTime {
	return NewDailyCycleTime(DailyTime.startTime + offsetSeconds)
}

func NewOffsetWeeklyTime(offsetSeconds int64) *CycleTime {
	return NewWeeklyCycleTime(Sunday.startTime + offsetSeconds)
}

type CycleTime struct {
	startTime int64 // unix时间戳，单位秒

	period int64 // 循环时间，单位秒
}

func (c *CycleTime) PrevTime(t time.Time) time.Time {

	ut := t.Unix()

	d := (ut - c.startTime) % c.period
	if d < 0 {
		d += c.period
	}

	return time.Unix(ut-d, 0)
}

func (c *CycleTime) NextTime(t time.Time) time.Time {

	ut := t.Unix()

	d := (ut - c.startTime) % c.period
	if d < 0 {
		d += c.period
	}

	return time.Unix(ut-d+c.period, 0)
}

func (c *CycleTime) Duration(t time.Time) time.Duration {
	return t.Sub(c.PrevTime(t))
}
