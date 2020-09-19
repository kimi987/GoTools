package logp

import (
	"runtime"
	"path/filepath"
	"github.com/lightpaw/logrus"
)

func init() {
	//logrus.SetLevel(logrus.DebugLevel)
}

func Debug(str string) {
	filename, funcname, line, _ := logPlus()
	logrus.Debugf("%v %v:%v\n  %v", filename, funcname, line, str)
}

func Debugf(format string, args ...interface{}) {
	filename, funcname, line, _ := logPlus()
	logrus.Debugf("%v %v:%v\n  " + format, filename, funcname, line, args)
}

func Info(str string) {
	filename, funcname, line, _ := logPlus()
	logrus.Infof("%v %v:%v\n  %v", filename, funcname, line, str)
}

func Infof(format string, args ...interface{}) {
	filename, funcname, line, _ := logPlus()
	logrus.Infof("%v %v:%v\n  " + format, filename, funcname, line, args)
}

func Warn(str string) {
	filename, funcname, line, _ := logPlus()
	logrus.Warnf("%v %v:%v\n  %v", filename, funcname, line, str)
}

func Warnf(format string, args ...interface{}) {
	filename, funcname, line, _ := logPlus()
	logrus.Warnf("%v %v:%v\n  " + format, filename, funcname, line, args)
}

func Error(format string, args ...interface{}) {
	filename, funcname, line, _ := logPlus()
	logrus.Errorf("%v %v:%v\n  " + format, filename, funcname, line, args)
}

func Errorf(format string, args ...interface{}) {
	filename, funcname, line, _ := logPlus()
	logrus.Errorf("%v %v:%v\n  " + format, filename, funcname, line, args)
}

func logPlus() (filename, funcname string, line int, ok bool) {
	pc, filename, line, ok := runtime.Caller(2)
	funcname = runtime.FuncForPC(pc).Name()
	filename = filepath.Base(filename)

	return
}