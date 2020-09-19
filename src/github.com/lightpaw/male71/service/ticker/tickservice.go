package ticker

import (
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/timeutil"
	"time"
)

func NewTickerService(timeService iface.TimeService, datas iface.ConfigDatas) *TickerService {

	ctime := timeService.CurrentTime()
	nextDailyResetTime := datas.MiscConfig().GetNextResetDailyTime(ctime)
	nextWeeklyResetTime := datas.MiscConfig().GetNextResetWeeklyTime(ctime)
	nextDailyMcResetTime := singleton.GetNextResetDailyTime(ctime, datas.MingcMiscData().DailyUpdateMingcTime)

	t := &TickerService{}
	t.dailyMcTicker = NewTicker(ctime, nextDailyMcResetTime.Sub(ctime), 24*time.Hour)
	t.weeklyTicker = NewTicker(ctime, nextWeeklyResetTime.Sub(ctime), 24*7*time.Hour)
	t.dailyZeroTicker = NewTicker(ctime, singleton.GetNextResetDailyTime(ctime, time.Duration(0)).Sub(ctime), 24*time.Hour)
	t.dailyTicker = NewTicker(ctime, nextDailyResetTime.Sub(ctime), 24*time.Hour)

	prevDailyResetTime := nextDailyResetTime.Add(-24 * time.Hour)

	nextTickTime := timeutil.NextTickTime(prevDailyResetTime, ctime, time.Minute)
	t.tickPerMinute = NewTicker(ctime, nextTickTime.Sub(ctime), time.Minute)

	nextTickTime = timeutil.NextTickTime(prevDailyResetTime, ctime, 10*time.Minute)
	t.tickPer10Minute = NewTicker(ctime, nextTickTime.Sub(ctime), 10*time.Minute)

	nextTickTime = timeutil.NextTickTime(prevDailyResetTime, ctime, 30*time.Minute)
	t.tickPer30Minute = NewTicker(ctime, nextTickTime.Sub(ctime), 30*time.Minute)

	nextTickTime = timeutil.NextTickTime(prevDailyResetTime, ctime, 60*time.Minute)
	t.tickPerHour = NewTicker(ctime, nextTickTime.Sub(ctime), 60*time.Minute)

	return t
}

//gogen:iface
type TickerService struct {
	weeklyTicker    *Ticker

	dailyTicker     *Ticker
	dailyZeroTicker *Ticker
	dailyMcTicker   *Ticker

	tickPerMinute   *Ticker
	tickPer10Minute *Ticker
	tickPer30Minute *Ticker
	tickPerHour     *Ticker
}

func (ts *TickerService) Close() {
	ts.dailyMcTicker.Stop()
	ts.weeklyTicker.Stop()
	ts.dailyZeroTicker.Stop()
	ts.dailyTicker.Stop()

	ts.tickPerMinute.Stop()
	ts.tickPer10Minute.Stop()
	ts.tickPer30Minute.Stop()
	ts.tickPerHour.Stop()
}

func (ts *TickerService) GetDailyMcTickTime() tickdata.TickTime {
	return ts.dailyMcTicker.GetTickTime()
}

func (ts *TickerService) GetWeeklyTickTime() tickdata.TickTime {
	return ts.weeklyTicker.GetTickTime()
}

func (ts *TickerService) GetDailyZeroTickTime() tickdata.TickTime {
	return ts.dailyZeroTicker.GetTickTime()
}

func (ts *TickerService) GetDailyTickTime() tickdata.TickTime {
	return ts.dailyTicker.GetTickTime()
}

func (ts *TickerService) GetPerMinuteTickTime() tickdata.TickTime {
	return ts.tickPerMinute.GetTickTime()
}

func (ts *TickerService) GetPer10MinuteTickTime() tickdata.TickTime {
	return ts.tickPer10Minute.GetTickTime()
}

func (ts *TickerService) GetPer30MinuteTickTime() tickdata.TickTime {
	return ts.tickPer30Minute.GetTickTime()
}

func (ts *TickerService) GetPerHourTickTime() tickdata.TickTime {
	return ts.tickPerHour.GetTickTime()
}

func (ts *TickerService) TickTickPerWeek(name string, f iface.TickFunc) iface.Func {
	return tick(ts.weeklyTicker, name, f)
}

func (ts *TickerService) TickPerDayZero(name string, f iface.TickFunc) iface.Func {
	return tick(ts.dailyZeroTicker, name, f)
}

func (ts *TickerService) TickPerDay(name string, f iface.TickFunc) iface.Func {
	return tick(ts.dailyTicker, name, f)
}

func (ts *TickerService) TickPerMinute(name string, f iface.TickFunc) iface.Func {
	return tick(ts.tickPerMinute, name, f)
}

func (ts *TickerService) TickPer10Minute(name string, f iface.TickFunc) iface.Func {
	return tick(ts.tickPer10Minute, name, f)
}

func (ts *TickerService) TickPer30Minute(name string, f iface.TickFunc) iface.Func {
	return tick(ts.tickPer30Minute, name, f)
}

func (ts *TickerService) TickPerHour(name string, f iface.TickFunc) iface.Func {
	return tick(ts.tickPerHour, name, f)
}

func tick(ticker *Ticker, name string, f iface.TickFunc) (stop iface.Func) {

	closeNotify := make(chan struct{})
	loopNotify := make(chan struct{})

	go func() {
		defer close(loopNotify)

		t := ticker.GetTickTime()
		for {
			select {
			case <-t.Tick():
				t = ticker.GetTickTime()
				call.CatchLoopPanic(func() {
					f(t)
				}, name)
			case <-closeNotify:
				return
			}
		}
	}()

	return func() {
		close(closeNotify)
		<-loopNotify
	}
}
