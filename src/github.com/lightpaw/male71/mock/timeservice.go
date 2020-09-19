package mock

import (
	"github.com/lightpaw/male7/gen/ifacemock"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
)

var ctime = time.Now()

func init() {
	ifacemock.TimeService.Mock(ifacemock.TimeService.CurrentTime, func() time.Time {
		return ctime
	})
}

func SetTime(toSet time.Time) {
	ctime = toSet
}

func IncSecond() {
	SetTime(ctime.Add(time.Second))
}

func IncMinute() {
	SetTime(ctime.Add(time.Minute))
}

func IncHour() {
	SetTime(ctime.Add(time.Hour))
}

func IncDay() {
	SetTime(ctime.Add(24 * time.Hour))
}

func TruncateDay() {
	SetTime(timeutil.DailyTime.PrevTime(ctime))
}
