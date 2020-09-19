package heroservice

import (
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/pbutil"
	"sync"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/logrus"
	"runtime/debug"
	"github.com/lightpaw/male7/service/monitor/metrics"
)

var (
	cache = sync.Pool{New: func() interface{} {
		return &MsgSender{
			list:     make([]pbutil.Buffer, 0, 3),
			bcList:   make([]pbutil.Buffer, 0, 3),
			funcList: make([]func() pbutil.Buffer, 0, 3),
		}
	}}
)

type MsgSender struct {
	changed  bool
	ok       bool
	list     []pbutil.Buffer
	funcList []func() pbutil.Buffer

	bcList []pbutil.Buffer
}

func (m *MsgSender) Add(msg pbutil.Buffer) {
	m.list = append(m.list, msg)
}

func (m *MsgSender) AddBroadcast(msg pbutil.Buffer) {
	m.bcList = append(m.bcList, msg)
}

func (m *MsgSender) AddFunc(msgFunc func() pbutil.Buffer) {
	m.funcList = append(m.funcList, msgFunc)
}

func (m *MsgSender) Changed() {
	m.changed = true
}

func (m *MsgSender) Ok() {
	m.ok = true
}

func getMsgSender() *MsgSender {
	return cache.Get().(*MsgSender)
}

func (m *MsgSender) send(sender sender.Sender) {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", err).Errorf("recovered from MsgSender.%v panic. SEVERE!!!", "MsgSender.send")
			metrics.IncPanic()
			if entry, ok := err.(error); ok {
				err = entry
			}
		}
	}()

	sender.SendAll(m.list)

	if len(m.funcList) > 0 {
		for _, f := range m.funcList {
			msg := f()
			if msg != nil {
				sender.Send(f())
			}
		}
	}

}

func (m *MsgSender) sendBroadcast(world iface.WorldService) {
	for _, msg := range m.bcList {
		world.Broadcast(msg)
	}
}

func (m *MsgSender) clear() {
	m.changed = false
	m.ok = false
	m.list = m.list[:0]
	m.bcList = m.bcList[:0]
	m.funcList = m.funcList[:0]
	cache.Put(m)
}
