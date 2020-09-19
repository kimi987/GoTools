package discover

import (
	"bytes"
	"github.com/lightpaw/go-zookeeper/zk"
	"github.com/lightpaw/logrus"
	"sync"
	"time"
)

type Node struct {
	FullPath       string
	conn           *zk.Conn
	nodeShouldQuit chan struct{} // 自己节点是否该退出
	eventChan      chan NodeChangeEvent
	closeOnce      sync.Once

	lastNotifiedData []byte
}

// 节点数据变化事件. 节点如果被删除, 则event中的data为nil
type NodeChangeEvent struct {
	FullPath string
	Data     []byte
}

// 监听一个节点
func NewListenNode(conn *zk.Conn, fullPath string) *Node {
	node := &Node{FullPath: fullPath, conn: conn, nodeShouldQuit: make(chan struct{}), eventChan: make(chan NodeChangeEvent, 10)}
	go node.loop()
	return node
}

// 不再监听节点变化
func (n *Node) Close() {
	n.closeOnce.Do(func() {
		close(n.nodeShouldQuit)
	})
}

// 获取节点数据变化事件, 节点如果被删除, 则event中的data为nil
func (n *Node) EventChan() <-chan NodeChangeEvent {
	// 只读chan, 外面不能write/close
	return n.eventChan
}

func (n *Node) loop() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).Error("discover.listennode.loop recovered from panic. SEVERE!!!")
		}

		close(n.eventChan)
	}()

	for {
		data, _, events, err := n.conn.GetW(n.FullPath)
		if err != nil {
			switch err {
			case zk.ErrClosing:
				// node closing, using select to quit. do nothing
			case zk.ErrNoNode:
				if n.lastNotifiedData != nil {
					// node deleted
					n.eventChan <- NodeChangeEvent{FullPath: n.FullPath}
					n.lastNotifiedData = nil
				}

				// set exists watcher
				exists, _, events, err := n.conn.ExistsW(n.FullPath)
				if err != nil {
					logrus.WithError(err).WithField("path", n.FullPath).Error("fail to create exists watcher")
					continue
				}

				if exists {
					continue
				}

				select {
				case <-events:
					continue

				case <-n.nodeShouldQuit:
					return

				case <-n.conn.ClosedChan():
					return
				}
			default:
				logrus.WithError(err).WithField("full path", n.FullPath).Error("discover.listennode.loop fail to get data")
			}

			select {
			case <-time.After(100 * time.Millisecond):
				continue

			case <-n.nodeShouldQuit:
				return

			case <-n.conn.ClosedChan():
				return
			}
		}

		logrus.WithField("data", data).Debug("receive node update")
		if !bytes.Equal(data, n.lastNotifiedData) {
			n.eventChan <- NodeChangeEvent{FullPath: n.FullPath, Data: data}
			n.lastNotifiedData = data
		}

		select {
		case <-events:
			continue

		case <-n.nodeShouldQuit:
			return

		case <-n.conn.ClosedChan():
			return
		}
	}
}
