package timeutil

import (
	"time"
	"strings"
	"strconv"
	"github.com/pkg/errors"
)

const (
	SecondsPerDay int64 = 24 * 60 * 60
	offsetHour          = 8 // 东八区
)

var (
	East8 = time.FixedZone("East-8", offsetHour*60*60)
	//East8Offset = 8 * time.Hour // (duration east of UTC).

	GameZone = East8
	//GameZoneOffset = East8Offset

	StartTime = time.Date(2000, 1, 1, 0, 0, 0, 0, GameZone)

	DailyTime = NewDailyCycleTime(StartTime.Unix())

	Sunday    = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Sunday).Unix())
	Monday    = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Monday).Unix())
	Tuesday   = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Tuesday).Unix())
	Wednesday = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Wednesday).Unix())
	Thursday  = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Thursday).Unix())
	Friday    = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Friday).Unix())
	Saturday  = NewWeeklyCycleTime(NextWeekTime(StartTime, time.Saturday).Unix())

	weeks = [...]*CycleTime{
		Sunday,
		Monday,
		Tuesday,
		Wednesday,
		Thursday,
		Friday,
		Saturday,
	}

	DayLayout         = "2006-01-02"
	SecondsLayout     = "2006-01-02_15:04:05"
	DaySlashLayout    = "2006/01/02"
	DaySlashLayoutLen = len(DaySlashLayout)
)

func WeekCycleTime(d time.Weekday) *CycleTime {
	return weeks[d]
}

func NextWeekTime(t time.Time, weekday time.Weekday) time.Time {

	wd := t.Weekday()
	diff := weekday - wd
	if diff < 0 {
		diff += 7
	}

	return t.Add(time.Duration(int64(diff)*SecondsPerDay) * time.Second)
}

func ParseDayLayout(value string) (time.Time, error) {
	return time.ParseInLocation(DayLayout, value, GameZone)
}

func ParseDaySlashLayout(value string) (time.Time, error) {
	return time.ParseInLocation(DaySlashLayout, value, GameZone)
}

func ParseSecondsLayout(value string) (time.Time, error) {
	return time.ParseInLocation(SecondsLayout, value, GameZone)
}

// 自动补全日期
func CompletionMMDD(value, sep string) string {
	// 2018-1-2 -> 2018-01-02

	if len(value) < DaySlashLayoutLen {
		array := strings.SplitN(value, sep, 3)
		if len(array) == 3 {
			if len(array[1]) == 1 {
				array[1] = "0" + array[1]
			}
			if len(array[2]) == 1 {
				array[2] = "0" + array[2]
			}

			return array[0] + sep + array[1] + sep + array[2]
		}
	}
	return value
}

func ParseHMS(value string) (hour, minute, second int, err error) {
	if len(value) <= 0 {
		return
	}

	hhmmss := strings.Split(value, ":")
	n := len(hhmmss)
	if n <= 0 {
		return 0, 0, 0, errors.Errorf("parse hms fail, [%v]", value)
	}

	// 小时分钟秒
	if hour, err = strconv.Atoi(hhmmss[0]); err != nil {
		return 0, 0, 0, errors.Wrapf(err, "parse hms fail, %v", value)
	}

	if hour < 0 || hour >= 24 {
		return 0, 0, 0, errors.Errorf("parse hms fail, invalid hour, %v", value)
	}

	if n <= 1 {
		return
	}

	// 分钟
	if minute, err = strconv.Atoi(hhmmss[1]); err != nil {
		return 0, 0, 0, errors.Wrapf(err, "parse hms fail, %v", value)
	}

	if minute < 0 || minute >= 60 {
		return 0, 0, 0, errors.Errorf("parse hms fail, invalid minute, %v", value)
	}

	if n <= 2 {
		return
	}

	// 秒
	if second, err = strconv.Atoi(hhmmss[2]); err != nil {
		return 0, 0, 0, errors.Wrapf(err, "parse hms fail, %v", value)
	}

	if second < 0 || second >= 60 {
		return 0, 0, 0, errors.Errorf("parse hms fail, invalid second, %v", value)
	}

	return
}
