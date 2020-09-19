package data

import (
	"testing"
	. "github.com/onsi/gomega"
	"github.com/gorhill/cronexpr"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
)

func TestTimeRule(t *testing.T) {
	RegisterTestingT(t)

	cron := cronexpr.MustParse("0 0 0 9 1 * 2019")

	t20190109 := time.Date(2019, 1, 9, 0, 0, 0, 0, timeutil.East8)
	Ω(cron.Next(t20190109)).Should(Equal(time.Time{}))

	t20190108 := t20190109.Add(-time.Second)
	Ω(cron.Next(t20190108)).Should(Equal(t20190109))

	t201901091 := t20190109.Add(time.Second)
	Ω(cron.Next(t201901091)).Should(Equal(time.Time{}))

	serverStartTime, _ := timeutil.ParseSecondsLayout("2019-01-09_00:00:00")

	// daily
	d := TimeRuleData{}
	d.RuleType = TimeRuleTypeDaily
	d.Rule = ""
	d.Time = "15:30"
	d.TimeDuration = time.Hour
	d.Init("")

	var array = [][2]string{
		{"2019-01-09_00:00:00", "2019-01-09_15:30:00"},
		{"2019-01-09_15:30:00", "2019-01-09_15:30:00"},
		{"2019-01-09_16:30:00", "2019-01-10_15:30:00"},
	}

	for _, a := range array {
		ctime, _ := timeutil.ParseSecondsLayout(a[0])
		nextTime, _ := timeutil.ParseSecondsLayout(a[1])
		Ω(d.Next(serverStartTime, ctime)).Should(Equal(nextTime))
	}

	// date
	d = TimeRuleData{}
	d.RuleType = TimeRuleTypeDate
	d.Rule = "2019-01-10"
	d.Time = "15:30"
	d.TimeDuration = time.Hour
	d.Init("")

	array = [][2]string{
		{"2019-01-10_00:00:00", "2019-01-10_15:30:00"},
		{"2019-01-10_15:30:00", "2019-01-10_15:30:00"},
		{"2019-01-10_16:30:00", ""},
	}

	for _, a := range array {
		ctime, _ := timeutil.ParseSecondsLayout(a[0])

		var nextTime time.Time
		if len(a[1]) > 0 {
			nextTime, _ = timeutil.ParseSecondsLayout(a[1])
		}
		Ω(d.Next(serverStartTime, ctime)).Should(Equal(nextTime))
	}

	// TimeRuleTypeWeek
	d = TimeRuleData{}
	d.RuleType = TimeRuleTypeWeek
	d.Rule = "w4"
	d.Time = "15:30"
	d.TimeDuration = time.Hour
	d.Init("")

	array = [][2]string{
		{"2019-01-10_00:00:00", "2019-01-10_15:30:00"},
		{"2019-01-10_15:30:00", "2019-01-10_15:30:00"},
		{"2019-01-10_16:30:00", "2019-01-17_15:30:00"},
	}

	for _, a := range array {
		ctime, _ := timeutil.ParseSecondsLayout(a[0])

		var nextTime time.Time
		if len(a[1]) > 0 {
			nextTime, _ = timeutil.ParseSecondsLayout(a[1])
		}
		Ω(d.Next(serverStartTime, ctime)).Should(Equal(nextTime))
	}

	// TimeRuleTypeWeek
	d = TimeRuleData{}
	d.RuleType = TimeRuleTypeServerStart
	d.Rule = "3"
	d.Time = "15:30"
	d.TimeDuration = time.Hour
	d.Init("")

	array = [][2]string{
		{"2019-01-10_00:00:00", "2019-01-11_15:30:00"},
		{"2019-01-11_15:30:00", "2019-01-11_15:30:00"},
		{"2019-01-11_16:30:00", ""},
	}

	for _, a := range array {
		ctime, _ := timeutil.ParseSecondsLayout(a[0])

		var nextTime time.Time
		if len(a[1]) > 0 {
			nextTime, _ = timeutil.ParseSecondsLayout(a[1])
		}
		Ω(d.Next(serverStartTime, ctime)).Should(Equal(nextTime))
	}
}
