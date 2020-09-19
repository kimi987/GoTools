package grmon

import (
	"net/http"
	_ "net/http/pprof"
	"github.com/lightpaw/male7/util/netlis"
	"github.com/lightpaw/logrus"
)

func Start(addr string) {

	l, port, err := netlis.ListenTcp(addr)
	if err != nil {
		logrus.WithError(err).Panic("监听grmon端口失败")
	}
	logrus.WithField("port", port).Info("监听grmon端口")

	go http.Serve(l, nil)
}
