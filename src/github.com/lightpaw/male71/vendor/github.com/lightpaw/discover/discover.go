package discover

import (
	"bytes"
	"github.com/lightpaw/go-zookeeper/zk"
	"github.com/lightpaw/logrus"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	Added   EventType = 1
	Removed EventType = 2
	Updated EventType = 3
)

// 节点更新时间的类型
type EventType int

// 路径下有节点数据更新
type Event struct {
	Type     EventType
	FullPath string
	Path     string
	Data     []byte
	OldData  []byte
	ID       int64
}

type PathChildrenCache struct {
	path       string
	conn       *Conn
	eventChan  chan Event
	shouldQuit chan struct{}
	version    uint64
	quitOnce   sync.Once
	pathPrefix string // 取到节点路径时, 扣掉这个就是节点的id

	children   map[string]*pathChildrenNode
	dataEvents chan dataUpdateEvent
}

// 获得zk下一个路径下的所有节点和节点的数据. 节点的增加/删除/更新 都通过Event chan来提醒
func (c *Conn) Discover(path string) (*PathChildrenCache, <-chan Event) {
	// get the last /
	idx := strings.LastIndex(path, "/")

	realPath := ""
	pathPrefix := ""
	if idx > 0 {
		realPath = path[:idx]
		pathPrefix = path[idx+1:]
	} else if idx == 0 {
		realPath = ""
		pathPrefix = ""
	} else {
		panic("path must starts with /")
	}

	result := &PathChildrenCache{
		path:       realPath,
		pathPrefix: pathPrefix,
		conn:       c,
		eventChan:  make(chan Event, 1000),
		shouldQuit: make(chan struct{}),
		children:   make(map[string]*pathChildrenNode),
		dataEvents: make(chan dataUpdateEvent, 1000),
	}

	go result.loop()
	return result, result.eventChan
}

// 不再监听这个节点数据. 会close Event chan
func (p *PathChildrenCache) Close() {
	p.quitOnce.Do(func() {
		close(p.shouldQuit)
	})
}

func (p *PathChildrenCache) loop() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).Error("Discover.loop recovered from panic. SEVERE!!!")
			p.Close()
		}

		logrus.WithField("listen path", p.path).Debug("Discover.loop exited")

		// quit all nodes

		for _, node := range p.children {
			node.Close()
		}
		p.children = make(map[string]*pathChildrenNode) // gc old map
		close(p.eventChan)
	}()

loop:
	for {
		paths, _, events, err := p.conn.conn.ChildrenW(p.path)
		if err != nil {
			switch err {
			case zk.ErrClosing:
			// zk closing, using select to quit. do nothing
			case zk.ErrNoNode:
				exists, _, events, err := p.conn.conn.ExistsW(p.path)
				if err != nil {
					logrus.WithError(err).WithField("path", p.path).Error("fail to create exists watch")
					continue
				}

				if exists {
					continue
				}

				select {
				case <-events:
					continue

				case <-p.shouldQuit:
					return

				case <-p.conn.shouldQuit:
					return
				}
			default:
				logrus.WithError(err).WithField("path", p.path).Error("Zk error getting path chileren")
			}

			select {
			case <-time.After(100 * time.Millisecond):
				continue loop

			case <-p.shouldQuit:
				return

			case <-p.conn.shouldQuit:
				return
			}
		}

		// got all children
		logrus.WithField("node count", len(paths)).WithField("path", p.path).Debug("Zk read path children successful")
		version := p.version + 1
		p.version = version

		for _, nodePath := range paths {
			node, has := p.children[nodePath]
			if has {
				node.version = version
			} else {
				// 之前不存在, 创建一个
				logrus.WithField("path", nodePath).WithField("parent path", p.path).Debug("Zk found new path child")
				if !strings.HasPrefix(nodePath, p.pathPrefix) {
					logrus.WithField("path", nodePath).WithField("parent prefix", p.pathPrefix).Error("Zk found path child, but with different path prefix. ignore")
					continue
				}

				nodeId := nodePath[len(p.pathPrefix):]
				id, err := strconv.ParseInt(nodeId, 10, 0)
				if err != nil {
					logrus.WithError(err).WithField("path", nodePath).WithField("perent prefix", p.pathPrefix).Error("Zk found path child, but cannot parse to int64. ignore")
					continue
				}

				node = p.newNode(nodePath, id, version)
				p.children[nodePath] = node
				go node.loop()
			}
		}

		// 检查是否有删掉的节点
		for nodePath, node := range p.children {
			if node.version != version {
				// 已删掉
				logrus.WithField("path", nodePath).WithField("parent path", p.path).Debug("Zk found removed path child")

				delete(p.children, nodePath)
				node.Close()

				if node.addNotified {
					p.sendEvent(Removed, node.id, nodePath, node.fullPath, node.data, node.data)
				}
			}
		}

		// 监听路径变化事件
	outer:
		for {
			select {
			case event := <-events:
				logrus.WithField("event", event).WithField("path", p.path).Debug("Zk path children got new event")
				continue loop

			case dataEvent := <-p.dataEvents:
				node, has := p.children[dataEvent.path]
				if !has {
					// 节点已经删掉了
					continue outer
				}

				if node.addNotified {
					if !bytes.Equal(node.data, dataEvent.newData) {
						// node update
						logrus.WithField("path", node.path).WithField("parent path", p.path).Debug("Zk found node data update")
						p.sendEvent(Updated, node.id, node.path, node.fullPath, node.data, dataEvent.newData)
						node.data = dataEvent.newData
					}
				} else {
					// new node
					logrus.WithField("path", node.path).WithField("parent path", p.path).Debug("Zk found new node")
					p.sendEvent(Added, node.id, node.path, node.fullPath, dataEvent.newData, dataEvent.newData)
					node.data, node.addNotified = dataEvent.newData, true
				}

			case <-p.shouldQuit:
				return

			case <-p.conn.shouldQuit:
				return
			}
		}

	}
}

func (p *PathChildrenCache) sendEvent(eventType EventType, id int64, path, fullPath string, oldData, data []byte) {
	event := Event{
		Type:     eventType,
		ID:       id,
		Path:     path,
		FullPath: fullPath,
		OldData:  oldData,
		Data:     data,
	}

	p.eventChan <- event
}

// --- node ---

type dataUpdateEvent struct {
	path    string
	newData []byte
}

type pathChildrenNode struct {
	path     string // 在目录中的名字
	fullPath string // 包括完整的目录的名字
	id       int64  // 节点id, zk生成的int64

	parent  *PathChildrenCache
	version uint64

	addNotified    bool
	nodeShouldQuit chan struct{}

	data []byte
}

func (p *PathChildrenCache) newNode(path string, id int64, version uint64) *pathChildrenNode {
	return &pathChildrenNode{
		path:           path,
		fullPath:       p.path + "/" + path,
		id:             id,
		parent:         p,
		version:        version,
		nodeShouldQuit: make(chan struct{}),
	}
}

func (n *pathChildrenNode) Close() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).Error("pathChildrenNode.Close recovered from panic. SEVERE!!!")
		}
	}()

	close(n.nodeShouldQuit)
}

func (n *pathChildrenNode) loop() {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).Error("pathChildrenNode.loop recovered from panic. SEVERE!!!")
		}

		logrus.WithField("listen path", n.fullPath).Debug("pathChildrenCnode.loop exited")
	}()

loop:
	for {
		data, _, events, err := n.parent.conn.conn.GetW(n.fullPath)
		if err != nil {
			if err != zk.ErrClosing {
				if err != zk.ErrNoNode {
					// if no node, do not complain
					logrus.WithError(err).WithField("full path", n.fullPath).Error("Zk pathChildrenNode fail to get data")
				}
			}

			select {
			case <-time.After(100 * time.Millisecond):
				continue loop

			case <-n.nodeShouldQuit:
				return

			case <-n.parent.shouldQuit:
				return

			case <-n.parent.conn.shouldQuit:
				return
			}
		}

		// 正常获得
		logrus.WithField("full path", n.fullPath).WithField("data", data).Debug("Zk pathChildrenNode got data")
		n.parent.dataEvents <- dataUpdateEvent{path: n.path, newData: data}

		select {
		case event := <-events:
			logrus.WithField("event", event).WithField("path", n.path).WithField("parent path", n.parent.path).Debug("Zk path children node got new event")
			continue loop

		case <-n.nodeShouldQuit:
			return

		case <-n.parent.shouldQuit:
			return

		case <-n.parent.conn.shouldQuit:
			return
		}
	}
}
