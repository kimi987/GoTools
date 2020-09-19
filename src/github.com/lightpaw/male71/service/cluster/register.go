package cluster

import (
	"github.com/golang/protobuf/proto"
	"github.com/lightpaw/discover"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/kv"
	"github.com/pkg/errors"
	"strconv"
	"time"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/build"
	"runtime/debug"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/rpc7"
	"github.com/lightpaw/male7/gen/iface"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/pbutil"
	"sync"
	"strings"
	"github.com/lightpaw/male7/util/atomic"
)

const (
	sessionTimeout      = 10 * time.Second
	PathPrefix          = "/m7/game"
	registerPath        = PathPrefix + "/t-"
	clientVersionZkPath = "/m7/cv/"
	clientVersionSep    = "#"
	LoginGrpcTarget     = "/m7/grpc/login"
	osAndroid           = "android"
	osiOS               = "ios"
)

//gogen:iface
type ClusterService struct {
	conn     *zkConn
	quitChan chan struct{}

	etcdClient  *clientv3.Client
	loginClient *rpc7.Client

	rpcAddr string

	clientVersion *client_version
}

func (c *ClusterService) Close() {
	close(c.quitChan)
	c.loginClient.Close()
	if c.etcdClient != nil {
		c.etcdClient.Close()
	}
}

func NewClusterService(serverConfig *kv.IndividualServerConfig, gs iface.GameServer, world iface.WorldService) *ClusterService {
	if serverConfig.IsCheckConfig {
		return &ClusterService{}
	}

	rpcAddr := fmt.Sprintf("%s:%d", serverConfig.RpcAddr, gs.GetRpcPort())

	// --- 压成proto ---
	data := &GameServerInfoProto{}
	data.ConnAddr = serverConfig.ConnectionServerAddr
	data.Name = serverConfig.ServerName
	data.Id = uint32(serverConfig.ServerID)
	data.Address = serverConfig.LocalAddr
	data.Port = gs.GetTcpPort()
	data.MetricsAddr = serverConfig.MetricsAddr
	data.Version = build.GetVersion()
	data.ConfigVersion = build.GetConfigVersion()
	data.ClientVersion = build.GetClientVersion()
	data.RpcAddr = rpcAddr
	data.BuildTime = int32(build.GetBuildUnixTime())

	dataBytes, err := proto.Marshal(data)
	if err != nil {
		logrus.WithError(err).Panic("marshal 注册的proto出错")
	}

	// --- 正式启动 ---
	logrus.Info("正在连接到zk")

	clientVersion := newClientVersion(serverConfig.ServerID)

	path := registerPath + strconv.Itoa(serverConfig.ServerID)
	closedChan := make(chan struct{})
	conn, err := newZkConn(serverConfig.ZkAddr, closedChan, path, dataBytes, world, clientVersion)
	if err != nil {
		logrus.WithError(err).Panic("无法连接到zk")
	}

	result := &ClusterService{
		conn:          conn,
		quitChan:      make(chan struct{}),
		clientVersion: clientVersion,
	}
	result.rpcAddr = rpcAddr

	if len(serverConfig.EtcdAddr) > 0 {
		// 连接etcd
		result.etcdClient, err = clientv3.New(clientv3.Config{
			Endpoints:   serverConfig.EtcdAddr,
			DialTimeout: 2 * time.Second,
		})
		if err != nil {
			logrus.WithError(err).Panic("连接etcd失败")
		}

		result.loginClient, err = rpc7.NewEtcdClient(result.etcdClient, LoginGrpcTarget)
		if err != nil {
			logrus.WithError(err).Panic("初始化登录服RPC（etcd）失败")
		}
	} else {
		result.loginClient, err = rpc7.NewClient(serverConfig.LoginClusterAddr)
		if err != nil {
			logrus.WithError(err).Panic("初始化登录服RPC失败")
		}
	}

	go call.CatchLoopPanic(func() {
		for {
			select {
			case <-result.quitChan:
				result.conn.conn.Close()
				return

			case <-closedChan:
				closedChan = make(chan struct{})

				for {
					logrus.Info("与zk连接已断开, 正在重连")
					conn, err := newZkConn(serverConfig.ZkAddr, closedChan, path, dataBytes, world, clientVersion)
					if err != nil {
						logrus.WithError(err).Error("与zk重连失败")
						time.Sleep(1 * time.Second)
						continue
					}
					result.conn = conn
					logrus.Info("zk 重连成功")
					break
				}
			}

		}
	}, "cluster.loop")

	return result
}

type zkConn struct {
	conn *discover.Conn
}

func newZkConn(zkAddr []string, closedChan chan struct{}, path string, data []byte,
	world iface.WorldService, clientVersion *client_version) (*zkConn, error) {
	closeNotify := &closedNotifier{closedChan: closedChan}

	conn, err := discover.Connect(zkAddr, sessionTimeout, closeNotify)
	if err != nil {
		return nil, errors.Wrap(err, "连接zk失败")
	}

	if err := conn.CreateEphemeral(path, data); err != nil {
		return nil, errors.Wrap(err, "创建节点失败: "+path)
	}

	if err := watchClientVersion(conn, closedChan, world, clientVersion); err != nil {
		return nil, errors.Wrap(err, "监听客户端版本号失败")
	}

	return &zkConn{conn: conn}, nil
}

func watchClientVersion(conn *discover.Conn, closedChan chan struct{}, world iface.WorldService,
	clientVersion *client_version) error {

	go WatchEvent(conn, clientVersionZkPath, closedChan, func(event discover.NodeEvent) {

		if ok, msg := clientVersion.handleEvent(clientVersionZkPath, event); ok {
			if msg != nil {
				// 广播给客户端，版本号变化了
				world.Broadcast(msg)
			}
		}
	})

	return nil
}

type closedNotifier struct {
	closedChan chan struct{}
}

func (c *closedNotifier) OnSessionExpired() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Error("Cluster.OnSessionExpired recovered from panic")
			metrics.IncPanic()
		}
	}()

	close(c.closedChan)
}

func (c *ClusterService) GetConfig(path string) ([]byte, error) {
	return loadAndUnencrypt(c.conn.conn.RawZkConn(), path)
}

func (c *ClusterService) EtcdClient() *clientv3.Client {
	return c.etcdClient
}

func (c *ClusterService) LoginClient() *rpc7.Client {
	return c.loginClient
}

func (c *ClusterService) RpcAddr() string {
	return c.rpcAddr
}

func (c *ClusterService) GMUpdateClientVersion(newVersion string) {
	//clientVersionPath := clientVersionZkPath + strconv.FormatInt(build.GetBuildUnixTime(), 10)
	//
	//if len(newVersion) > 0 {
	//	err := createNode(c.conn.conn.RawZkConn(), []byte(newVersion), clientVersionPath)
	//	if err != nil {
	//		logrus.WithError(err).Error("GMUpdateClientVersion createNode fail")
	//	}
	//} else {
	//	err := c.conn.conn.RawZkConn().Delete(clientVersionPath, -1)
	//	if err != nil {
	//		logrus.WithError(err).Error("GMUpdateClientVersion deleteNode fail")
	//	}
	//}

}

func (c *ClusterService) GetClientVersionMsg(os, tag string) pbutil.Buffer {
	return c.clientVersion.getMsg(os, tag)
}

func newClientVersion(sid int) *client_version {
	return &client_version{
		sid:       sid,
		sidString: strconv.Itoa(sid),
		android:   newOsClientVersion(osAndroid),
		ios:       newOsClientVersion(osiOS),
	}
}

type client_version struct {
	sid       int
	sidString string
	lock      sync.RWMutex

	android *os_client_version
	ios     *os_client_version
}

func (c *client_version) getMsg(os, tag string) pbutil.Buffer {
	oscv := c.android
	if os == osiOS {
		oscv = c.ios
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	if cvm := oscv.selfMap[tag]; cvm != nil {
		return cvm.getMsg()
	}
	return oscv.globalVersion.getMsg()
}

func (c *client_version) handleEvent(pathPrefix string, event discover.NodeEvent) (bool, pbutil.Buffer) {

	// sid_os_tag
	key := strings.Replace(event.Path, pathPrefix, "", 1)

	fields := strings.Split(key, clientVersionSep)

	if len(fields) != 2 && len(fields) != 3 {
		logrus.WithField("key", key).Error("热更新版本号，无效的key")
		return false, nil
	}

	sidStr := fields[0]
	os := fields[1]
	tag := ""
	if len(fields) > 2 {
		tag = fields[2]
	}

	isGlobalVersion := sidStr == "0"
	if isGlobalVersion {
		if tag != "" {
			logrus.WithField("key", key).Error("热更新版本号，全局版本不支持tag配置")
			return false, nil
		}
	} else {
		if sidStr != c.sidString {
			return false, nil
		}
	}

	var oscv *os_client_version
	switch strings.ToLower(os) {
	case osAndroid:
		oscv = c.android
	case osiOS:
		oscv = c.ios
	default:
		logrus.WithField("os", os).Error("热更新版本号，无效的os类型")
		return false, nil
	}

	if oscv == nil {
		logrus.Error("热更新版本号，oscv == nil")
		return false, nil
	}

	switch event.Type {
	case discover.Added, discover.Updated:

		version := string(event.Data)
		// 检查version是否有效的字符串

		cvm := newCvm(os, tag, version)

		c.lock.Lock()
		if isGlobalVersion {
			oscv.globalVersion = cvm
		} else {
			oscv.selfMap[tag] = cvm
		}
		c.lock.Unlock()

		return true, cvm.getMsg()

	case discover.Removed:

		if isGlobalVersion {
			logrus.WithField("key", key).Error("global version remove?")
			return false, nil
		}

		c.lock.Lock()
		delete(oscv.selfMap, tag)
		newClientVersion := oscv.globalVersion.cv
		c.lock.Unlock()

		return true, misc.NewS2cClientVersionMsg(newClientVersion, os, tag).Static()

	default:
		logrus.WithField("event", event).Error("未知的client_version更新事件")
		return false, nil
	}

}

func newOsClientVersion(os string) *os_client_version {
	return &os_client_version{
		os:            os,
		globalVersion: newCvm(os, "", build.GetConfigVersion()),
		selfMap:       make(map[string]*cvm),
	}
}

type os_client_version struct {
	os string

	// 公共版本号 key为
	globalVersion *cvm

	// 本服版本号
	selfMap map[string]*cvm
}

func newCvm(os, tag, cv string) *cvm {
	return &cvm{
		os:  os,
		tag: tag,
		cv:  cv,
	}
}

type cvm struct {
	os     string
	tag    string
	cv     string
	msgRef atomic.Value
}

func (m *cvm) getMsg() pbutil.Buffer {

	if ref := m.msgRef.Load(); ref != nil {
		return ref.(pbutil.Buffer)
	}

	msg := misc.NewS2cClientVersionMsg(m.cv, m.os, m.tag).Static()
	m.msgRef.Store(msg)
	return msg
}
