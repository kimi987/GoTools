package timeutil

import (
	"time"
	"regexp"
	"github.com/pkg/errors"
	"fmt"
	"strings"
	"strconv"
)

type CnWeekday int

const (
	Mon CnWeekday = iota + 1
	Tue
	Wed
	Thu
	Fri
	Sat
	Sun
)

func BuildWeekDurTime(str string, dur time.Duration) (t *WeekDurTime, err error) {
	w, d, e := getWeekAndDuration(str)
	if e != nil {
		err = e
		return
	}

	t = &WeekDurTime{WeekTime: &WeekTime{Week: w, Time: d}, Dur: dur}
	return
}

func getWeekAndDuration(str string) (w CnWeekday, t time.Duration, err error) {
	if ok, e := regexp.MatchString(`w[1-7],[0-2]??[0-9]h[0-5]??[0-9]m`, str); !ok || e != nil {
		err = errors.Errorf(fmt.Sprintf("时间格式错误:%v。必须为：w1,20h30m", str))
		return
	}

	ss := strings.Split(str, ",")
	if cw, e := strconv.ParseInt(ss[0][1:], 10, 32); e != nil {
		err = errors.Wrapf(e, fmt.Sprintf("时间格式错误:%v。必须为：w1,20h30m", str))
		return
	} else {
		w = CnWeekday(cw)
	}

	if t, err = time.ParseDuration(ss[1]); err != nil {
		err = errors.Wrapf(err, fmt.Sprintf("时间格式错误:%v。必须为：w1,20h30m", str))
		return
	}

	return
}

type WeekDurTime struct {
	*WeekTime
	Dur time.Duration
}

type WeekTime struct {
	Week CnWeekday
	Time time.Duration
}

func (t *WeekDurTime) ToString() string {
	return fmt.Sprintf("w%v %vm %vm", t.Week, t.Time.Minutes(), t.Dur.Minutes())
}

func ConvCnWeekday(w time.Weekday) CnWeekday {
	if w == time.Sunday {
		return CnWeekday(7)
	}
	return CnWeekday(w)
}

func (t *WeekDurTime) NextTime(ctime time.Time) (startTime, endTime time.Time) {
	startTime = t.CurrWeekTime(ctime)
	if ctime.After(startTime) {
		startTime = startTime.AddDate(0, 0, 7)
	}
	endTime = startTime.Add(t.Dur)
	return
}

func (t *WeekDurTime) CurrWeekTime(ctime time.Time) time.Time {
	result := DailyTime.PrevTime(ctime)
	d := int64(t.Week-ConvCnWeekday(ctime.Weekday())) * 24 * int64(time.Hour)
	result = result.Add(time.Duration(d)).Add(t.Time)

	return result
}

func (t *WeekTime) DurationWithInSameWeek(t1 *WeekTime) time.Duration {
	if t.Week == t1.Week && t.Time == t1.Time {
		return 0
	}

	var start, end *WeekTime
	if t.Week > t1.Week {
		start, end = t1, t
	} else if t.Week == t1.Week {
		if t.Time > t1.Time {
			start, end = t1, t
		} else {
			start, end = t, t1
		}
	} else {
		start, end = t, t1
	}

	return time.Duration(int64(end.Week-start.Week)*int64(time.Hour)*24 + int64(end.Time-start.Time))
}

func (t *WeekTime) Add(dur time.Duration) *WeekTime {
	r := &WeekTime{Week:t.Week, Time:t.Time}

	day := (r.Time + dur) / (time.Hour * 24)
	if day <= 0 {
		r.Time = r.Time + dur
		return r
	}

	w := int(t.Week) + int(day)
	if w > int(Sun) {
		w = w % int(Sun)
	}

	r.Week = CnWeekday(w)
	r.Time = (r.Time + dur) % (time.Hour * 24)

	return r
}
