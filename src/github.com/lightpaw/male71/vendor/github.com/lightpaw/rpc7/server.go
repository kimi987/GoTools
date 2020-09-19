package rpc7

import (
	"net"
	"strconv"
	"github.com/pkg/errors"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	"github.com/lightpaw/logrus"
)

func (mux *HandlerMux) Serve(name, addr string) (port uint32, listener net.Listener, loopNotifier chan struct{}, err error) {
	return DoServe(name, addr, mux.NewHandler())
}

func Serve(name, addr string) (port uint32, listener net.Listener, loopNotifier chan struct{}, err error) {
	return DefaultHandlerMux.Serve(name, addr)
}

func (mux *HandlerMux) ServeListener(name string, listener net.Listener) (port uint32, loopNotifier chan struct{}, err error) {
	return DoServeListener(name, listener, mux.NewHandler())
}

func ServeListener(name string, listener net.Listener) (port uint32, loopNotifier chan struct{}, err error) {
	return DefaultHandlerMux.ServeListener(name, listener)
}

func DoServe(name, addr string, handler RpcServiceServer) (port uint32, listener net.Listener, loopNotifier chan struct{}, err error) {
	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return 0, nil, nil, errors.Wrapf(err, "监听Tcp失败（rpc）")
	}

	port, loopNotifier, err = DoServeListener(name, listener, handler)
	return
}

func DoServeListener(name string, listener net.Listener, handler RpcServiceServer) (port uint32, loopNotifier chan struct{}, err error) {

	if handler == nil {
		return 0, nil, errors.Wrapf(err, "handler is nil")
	}

	_, rpcPort, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return 0, nil, errors.Wrapf(err, "截取Tcp端口号失败（rpc）")
	}
	rpcPortUint, err := strconv.ParseUint(rpcPort, 10, 32)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "截取Tcp端口号不是个数字（rpc） %s", rpcPort)
	}

	// 开启RPC
	rpcLoopNotifier := make(chan struct{})
	go CatchLoopPanic(func() {
		defer close(rpcLoopNotifier)

		s := grpc.NewServer()
		RegisterRpcServiceServer(s, handler)
		// Register reflection service on gRPC server.
		reflection.Register(s)
		//defer s.GracefulStop() // may be block
		defer s.Stop()

		if err := s.Serve(listener); err != nil {
			logrus.WithField("name", name).WithError(err).Error("rpc.Serve stoped")
		}
	}, name)

	return uint32(rpcPortUint), rpcLoopNotifier, nil
}
