package call

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"runtime/debug"
)

func CatchPanic(f func(), name string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", r).Errorf("%s recovered from panic. SEVERE!!!", name)
			metrics.IncPanic()
		}
	}()

	f()
}

func CatchLoopPanic(f func(), name string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", r).Errorf("%s recovered from panic. SEVERE!!!", name)
			metrics.IncPanic()
		}

		logrus.Infof("%s exit", name)
	}()

	f()
}
