package netlis

import (
	"net"
	"strconv"
	"github.com/pkg/errors"
)

func ListenTcp(addr string) (net.Listener, uint32, error) {
	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "监听Tcp端口失败，"+addr)
	}
	tcpPort, err := GetListenerPort(tcpListener)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "获取Tcp端口失败，"+addr)
	}

	return tcpListener, tcpPort, nil
}

func GetListenerPort(listener net.Listener) (uint32, error) {
	_, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return 0, errors.Wrapf(err, "截取Tcp端口号失败（rpc）")
	}
	portUint, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return 0, errors.Wrapf(err, "截取Tcp端口号不是个数字（rpc） %s", port)
	}
	return uint32(portUint), nil
}
