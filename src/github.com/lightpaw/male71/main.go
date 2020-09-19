package main

import (
	"flag"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/service"
	"github.com/lightpaw/male7/service/unmarshal"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/lfshook"
	"github.com/lightpaw/nonmux"
	"time"
	"net"
	"github.com/lightpaw/muxface"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/lightpaw/male7/bootstrap"
)

var (
	//conf = flag.String("config", "conf", "config base path")

	log       = flag.String("log", "", "log file path")
	logDay    = flag.Uint("logday", 7, "log day")
	logLevel  = flag.String("loglv", "debug", "log level")
	version   = flag.Bool("version", false, "version")
	pprofAddr = flag.String("pprof", ":0", "pprof port")
)

func main() {
	flag.Parse()

	if len(*log) > 0 {
		path := *log

		//pathMap := lfshook.PathMap{}
		//for _, lv := range logrus.AllLevels {
		//	pathMap[lv] = *log
		//}

		// rotatelog
		writer, _ := rotatelogs.New(
			path+".%Y%m%d",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithRotationTime(24*time.Hour),
			rotatelogs.WithRotationCount(*logDay),
		)

		writerMap := lfshook.WriterMap{}
		for _, lv := range logrus.AllLevels {
			writerMap[lv] = writer
		}

		logrus.AddHook(lfshook.NewHook(writerMap))
	}

	if toSet, err := logrus.ParseLevel(*logLevel); err != nil {
		logrus.WithField("lvl", *logLevel).WithError(err).Error("无效的日志级别")
	} else {
		logrus.SetLevel(toSet)
	}

	call.CatchPanic(func() {
		bootstrap.Start(*pprofAddr, *version, createNonMuxListener)
	}, "main")
}

func createNonMuxListener(listener net.Listener) (muxface.Listener, error) {
	config := &nonmux.Config{CloseWaitTime: 30 * time.Second, Unmarshaller: unmarshal.NewProtoUnmarshaller(), DontEncrypt: service.IndividualServerConfig.GetDontEncrypt()}
	return nonmux.Wrap(nonmux.Listen(listener, config))
}
