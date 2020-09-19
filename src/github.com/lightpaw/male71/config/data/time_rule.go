package data

import (
	"time"
	"github.com/gorhill/cronexpr"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/timeutil"
	"fmt"
	"github.com/lightpaw/male7/util/atomic"
	"strconv"
)

const (
	TimeRuleTypeDaily       = 0
	TimeRuleTypeDate        = 1
	TimeRuleTypeWeek        = 2
	TimeRuleTypeServerStart = 3
)

//gogen:config
type TimeRuleData struct {
	_ struct{} `file:"杂项/时间规则.txt"`

	Id           uint64
	RuleType     uint64 `validator:"uint" default:"0"` // 0-daily 1-date 2-week
	Rule         string
	Time         string
	TimeDuration time.Duration

	serverStartDuration time.Duration

	expressionRef atomic.Value
}

func (t *TimeRuleData) isCycle() bool {
	switch t.RuleType {
	case TimeRuleTypeDaily:
		return true
	case TimeRuleTypeWeek:
		return true
	}
	return false
}

func (t *TimeRuleData) getExpression() *cronexpr.Expression {
	if ref := t.expressionRef.Load(); ref != nil {
		return ref.(*cronexpr.Expression)
	}
	return nil
}

func (t *TimeRuleData) setExpresionIfNil(toSet *cronexpr.Expression) {
	t.expressionRef.Store(toSet)
}

func (t *TimeRuleData) Next(serverStartTime, ctime time.Time) time.Time {
	expression := t.getExpression()
	if expression == nil {
		// 里面已经set进去了
		expression = t.createServerStartTime(serverStartTime)
	}
	return expression.Next(ctime.Add(-t.TimeDuration))
}

func (d *TimeRuleData) Init(filename string) {

	hh, mm, ss, err := timeutil.ParseHMS(d.Time)
	if err != nil {
		logrus.WithError(err).Panicf("%s, time数据必须配置时分秒的格式[hh:mm:ss]:%s", filename, d.Time)
	}

	var cronLine string
	switch d.RuleType {
	case TimeRuleTypeDaily:
		cronLine = fmt.Sprintf("%d %d %d * * * *", ss, mm, hh)

	case TimeRuleTypeDate:

		rule := timeutil.CompletionMMDD(d.Rule, "/")
		date, err := timeutil.ParseDaySlashLayout(rule)
		if err != nil {
			logrus.WithError(err).Panicf("%s, rule数据是Date类型，必须配置年月日的格式[yyyy/mm/dd]:%s", filename, d.Rule)
		}

		year, month, day := date.Date()
		cronLine = fmt.Sprintf("%d %d %d %d %d * %d", ss, mm, hh, day, month, year)

	case TimeRuleTypeWeek:
		var dayOfWeek int
		switch d.Rule {
		case "w0", "W0", "w7", "W7":
			dayOfWeek = 0
		case "w1", "W1":
			dayOfWeek = 1
		case "w2", "W2":
			dayOfWeek = 2
		case "w3", "W3":
			dayOfWeek = 3
		case "w4", "W4":
			dayOfWeek = 4
		case "w5", "W5":
			dayOfWeek = 5
		case "w6", "W6":
			dayOfWeek = 6
		default:
			check.PanicNotTrue(false, "%s week类型的rule的格式必须是W1~W7或w1~w7 rule:%s", filename, d.Rule)
		}

		cronLine = fmt.Sprintf("%d %d %d * * %d *", ss, mm, hh, dayOfWeek)

	case TimeRuleTypeServerStart:
		serverStartDay, err := strconv.ParseUint(d.Rule, 10, 64)
		if err != nil {
			logrus.WithError(err).Panicf("%s 配置了无效的开服天数 %v", filename, d.Rule)
		}

		if serverStartDay <= 0 {
			logrus.WithError(err).Panicf("%s 配置了无效的开服天数[>=1] %v", filename, d.Rule)
		}

		d.serverStartDuration = time.Duration(serverStartDay-1)*timeutil.Day +
			time.Duration(hh)*time.Hour +
			time.Duration(mm)*time.Minute +
			time.Duration(ss)*time.Second

	default:
		// unkown type
		logrus.Panicf("%s 配置了无效的时间规则类型 %v", filename, d.Rule)
	}
	//fmt.Println(cronLine) // 经常需要打印看下

	if d.RuleType != TimeRuleTypeServerStart {
		d.setExpresionIfNil(cronexpr.MustParse(cronLine))
	}
}

func (t *TimeRuleData) createServerStartTime(serverStartTime time.Time) *cronexpr.Expression {
	startTime := serverStartTime.Add(t.serverStartDuration)

	year, month, day := startTime.Date()
	hh := startTime.Hour()
	mm := startTime.Minute()
	ss := startTime.Second()
	cronLine := fmt.Sprintf("%d %d %d %d %d * %d", ss, mm, hh, day, month, year)

	newCronExp := cronexpr.MustParse(cronLine)
	t.setExpresionIfNil(newCronExp)
	return newCronExp
}
