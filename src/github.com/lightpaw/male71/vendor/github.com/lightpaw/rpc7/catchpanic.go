package rpc7

import (
	"github.com/lightpaw/logrus"
	"runtime/debug"
)

func CatchPanic(f func(), name string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", r).Errorf("%s recovered from panic. SEVERE!!!", name)
		}
	}()

	f()
}

func CatchLoopPanic(f func(), name string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", r).Errorf("%s recovered from panic. SEVERE!!!", name)
		}

		logrus.Infof("%s exit", name)
	}()

	f()
}
