package check

import "github.com/lightpaw/logrus"

func PanicNotTrue(b bool, format string, args ...interface{}) {
	if !b {
		logrus.Panicf(format, args...)
	}
}
