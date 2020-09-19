package cluster

import (
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/rpc7"
	"net"
	"github.com/lightpaw/logrus"
	"os"
	"os/signal"
	"syscall"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/gen/iface"
	"fmt"
	"github.com/lightpaw/muxface"
	"github.com/lightpaw/male7/util/netlis"
)

func NewGameServer(serverConfig *kv.IndividualServerConfig) *GameServer {

	// 先把需要监听的端口，先监听了，先不启动处理
	tcpAddr := fmt.Sprintf(":%d", serverConfig.GetPort())
	tcpListener, tcpPort, err := netlis.ListenTcp(tcpAddr)
	if err != nil {
		logrus.WithError(err).Panic("监听Game端口失败")
	}

	rpcAddr := fmt.Sprintf(":%d", serverConfig.RpcPort)
	rpcListener, rpcPort, err := netlis.ListenTcp(rpcAddr)
	if err != nil {
		logrus.WithError(err).Panic("监听Rpc端口失败")
	}

	return &GameServer{
		serverConfig: serverConfig,
		sid:          uint32(serverConfig.GetServerID()),
		tcpListener:  tcpListener,
		tcpPort:      tcpPort,
		rpcListener:  rpcListener,
		rpcPort:      rpcPort,
	}
}

//gogen:iface
type GameServer struct {
	serverConfig *kv.IndividualServerConfig
	sid          uint32

	tcpListener net.Listener
	tcpPort     uint32

	rpcListener net.Listener
	rpcPort     uint32
}

func (s *GameServer) GetTcpPort() uint32 {
	return s.tcpPort
}

func (s *GameServer) GetRpcPort() uint32 {
	return s.rpcPort
}

func (s *GameServer) Serve(tcpServe iface.ServeListener, handler iface.ConnHandler) {

	// 启动rpc
	rpcPort, rpcLoopNotify, err := rpc7.ServeListener("game_rpc", s.rpcListener)
	if err != nil {
		logrus.WithError(err).Panic("启动游戏服Rpc失败")
	}
	logrus.WithField("port", rpcPort).Debug("启动Rpc监听端口")

	// 启动tcp
	lsn, err := tcpServe(s.tcpListener)
	if err != nil {
		logrus.WithError(err).Panic("启动游戏服MuxListener失败")
	}

	tcpLoopNotify, err := s.serveMuxListener(lsn, handler)
	if err != nil {
		logrus.WithError(err).Panic("启动游戏服Tcp失败")
	}

	go call.CatchLoopPanic(func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-rpcLoopNotify:
		case <-tcpLoopNotify:
		case <-sigs:
		}
		s.tcpListener.Close()
		s.rpcListener.Close()

	}, "main.sigs")

	<-tcpLoopNotify
	<-rpcLoopNotify
}

func (s *GameServer) serveMuxListener(listener muxface.Listener, handle iface.ConnHandler) (chan struct{}, error) {

	logrus.WithField("addr", listener.Addr().String()).Info("开始监听连接")

	loopNotifier := make(chan struct{})
	go call.CatchLoopPanic(func() {
		defer close(loopNotifier)

		for {
			conn, err := listener.Accept()
			if err != nil {
				logrus.WithError(err).Error("muxface.Listener Accept() fail")
				return
			}

			if IsRobotConn(conn) && !s.serverConfig.GetIsAllowRobot() {
				logrus.Debug("收到机器人登陆，服务器不允许机器人登陆，踢掉")
				conn.Close()
				continue
			}

			if s.sid != conn.GetLoginToken().GameServerID {
				logrus.WithField("sid", s.sid).WithField("sidto", conn.GetLoginToken().GameServerID).Debug("收到不是这个服的用户登陆，踢掉")
				conn.Close()
				continue
			}

			go handle(conn)
		}
	}, "TpcLoop")

	return loopNotifier, nil
}

func IsRobotConn(conn muxface.Conn) bool {
	return conn.GetLoginToken().Reserved&1 == 1
}

func (s *GameServer) serveHttp() {

}
