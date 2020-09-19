package discover

import (
	"github.com/lightpaw/go-zookeeper/zk"
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Connected    Status = 1
	Disconnected Status = 2
	Expired      Status = 3
)

type Status int32

type Conn struct {
	conn              *zk.Conn
	eventChan         <-chan zk.Event
	shouldQuit        chan struct{}
	connStateCallback ConnectionStateCallback
	quitOnce          sync.Once
	serverStatus      int32
}

// 与zk的连接callback
type ConnectionStateCallback interface {
	// 与zk连接expire时调用, 只会调用一次
	OnSessionExpired()
}

// 连接到zk. 只有连接成功时才会返回.
func Connect(servers []string, sessionTimeout time.Duration, listener ConnectionStateCallback) (*Conn, error) {
	zkConn, eventChan, err := zk.Connect(servers, sessionTimeout, zk.WithSessionExpireAndQuit())

	if err != nil {
		return nil, err
	}

	logrus.Info("Waiting for zk connection establish")

	for {
		event, ok := <-eventChan
		if !ok {
			return nil, errors.New("zk already closed")
		}

		if event.Type == zk.EventSession && event.State == zk.StateHasSession {
			logrus.WithField("session id", zkConn.SessionID()).Info("Zk connection established")
			break
		}
	}

	result := &Conn{
		conn:              zkConn,
		connStateCallback: listener,
		eventChan:         eventChan,
		shouldQuit:        make(chan struct{}),
		serverStatus:      int32(Connected),
	}

	go result.eventListenerLoop()

	return result, nil
}

func (c *Conn) RawZkConn() *zk.Conn {
	return c.conn
}

func (c *Conn) ServerStatus() Status {
	return Status(atomic.LoadInt32(&c.serverStatus))
}

func (c *Conn) setServerStatus(status Status) {
	atomic.StoreInt32(&c.serverStatus, int32(status))
}

func (c *Conn) Close() {
	c.quitOnce.Do(func() {
		c.conn.Close()
		close(c.shouldQuit)
	})
}

// 监听zk的session expire事件
func (c *Conn) eventListenerLoop() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).Error("Conn.eventListenerLoop recovered from panic. SEVERE!!!")
		}

		c.Close()
		c.setServerStatus(Expired)
		c.connStateCallback.OnSessionExpired()
		logrus.WithField("session id", c.conn.SessionID()).Debug("Conn.eventListenerLoop exited")
	}()

	for {
		select {
		case event, ok := <-c.eventChan:
			if !ok {
				logrus.WithField("session id", c.conn.SessionID()).Info("Discover found zk event chan closed. session must be closed")

				return
			}
			if event.Type == zk.EventSession {
				switch event.State {
				case zk.StateExpired:
					logrus.WithField("session id", c.conn.SessionID()).Info("Discover found zk session expired")
					return // quit loop
				case zk.StateConnected:
					c.setServerStatus(Connected)
				case zk.StateDisconnected:
					c.setServerStatus(Disconnected)
				}
			}

		case <-c.conn.ClosedChan():
			logrus.WithField("session id", c.conn.SessionID()).Info("Discover found zk session closed")
			return
		}
	}
}
