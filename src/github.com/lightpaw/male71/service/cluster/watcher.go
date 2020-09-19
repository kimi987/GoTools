package cluster

import (
	"github.com/lightpaw/go-zookeeper/zk"
	"github.com/lightpaw/logrus"
	"time"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/discover"
)

func WatchEvent(zkCluster *discover.Conn, zkPath string, closeChan <-chan struct{}, f func(event discover.NodeEvent)) {

	_, clientVersionEvents := zkCluster.ListenChildrens(zkPath)

	for {
		select {
		case event, ok := <-clientVersionEvents:
			if !ok {
				return
			}

			call.CatchPanic(func() {
				f(event)
			}, "watch data changed")

		case <-closeChan:
			return
		}
	}

}

func Watch(conn *zk.Conn, zkPath string, closeChan <-chan struct{}, f func(data []byte)) {

loop:
	for {

		// 先看下node是不是存在
		data, _, events, err := conn.GetW(zkPath)
		if err != nil {
			if err != zk.ErrClosing {
				if err != zk.ErrNoNode {
					// if no node, do not complain
					logrus.WithError(err).WithField("path", zkPath).Error("Zk watch fail to get data")
				} else {
					// node不存在，监听节点创建
					exist, _, events, err := conn.ExistsW(zkPath)
					if err != nil {
						logrus.WithError(err).WithField("path", zkPath).Error("Zk watch fail to get data")
					} else {
						if !exist {
							select {
							case event := <-events:
								// 节点创建了
								logrus.WithField("path", zkPath).WithField("event", event.Type).Info("node create")
								continue loop

							case <-closeChan:
								return
							}
						}
					}

				}
			}

			select {
			case <-time.After(time.Second):
				continue loop

			case <-closeChan:
				return
			}
		}

		call.CatchPanic(func() {
			f(data)
		}, "watch data changed")

		select {
		case event := <-events:
			logrus.WithField("path", zkPath).WithField("event", event.Type).Info("node changed")

			if event.Type == zk.EventNodeDeleted {
				call.CatchPanic(func() {
					f(nil)
				}, "watch data deleted")
			}

			continue loop

		case <-closeChan:
			return
		}

	}

}
