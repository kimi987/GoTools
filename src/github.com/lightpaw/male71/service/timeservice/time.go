package timeservice

import (
	"time"
	"io/ioutil"
	"strings"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/logrus"
)

const file = "ctime.txt"

//gogen:iface
type TimeService struct {
	ctimeFunc func() time.Time
}

func NewDefaultTimeService() *TimeService {
	return newTimeService(false)
}

func NewTimeService(config iface.IndividualServerConfig) *TimeService {
	return newTimeService(config.GetIsDebug())
}

func newTimeService(isDebug bool) *TimeService {
	ctimeFunc := currentTime

	if isDebug {
		if data, _ := ioutil.ReadFile(file); len(data) > 0 {
			s := strings.TrimSpace(string(data))
			if startTime, err := timeutil.ParseSecondsLayout(s); err == nil {
				diff := startTime.Sub(currentTime())
				ctimeFunc = func() time.Time {
					return currentTime().Add(diff)
				}
				logrus.WithField("time", ctimeFunc()).Info("服务器开启时间设置")
			}
		}
	}

	return &TimeService{
		ctimeFunc: ctimeFunc,
	}
}

func (ts *TimeService) CurrentTime() time.Time {
	return ts.ctimeFunc()
}

func currentTime() time.Time {
	return time.Now()
}
