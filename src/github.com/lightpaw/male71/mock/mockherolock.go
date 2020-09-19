package mock

import (
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/pbutil"
)

var LockResult = newLockResult()

func newLockResult() *lock_result {
	return &lock_result{}
}

type lock_result struct {
	changed bool
	ok      bool
	list    []pbutil.Buffer
	bcList  []pbutil.Buffer
}

func (m *lock_result) IsChanged() bool {
	return m.changed
}

func (m *lock_result) IsOk() bool {
	return m.ok
}

func (m *lock_result) PopMsg() pbutil.Buffer {
	n := len(m.list)
	if n <= 0 {
		return nil
	}

	msg := m.list[0]
	copy(m.list, m.list[1:])
	m.list = m.list[:n-1]

	return msg
}

func (m *lock_result) Add(msg pbutil.Buffer) {
	m.list = append(m.list, msg)
}

func (m *lock_result) AddBroadcast(msg pbutil.Buffer) {
	m.bcList = append(m.bcList, msg)
}

func (m *lock_result) AddTlog(f func()) {
}

func (m *lock_result) AddFunc(f func() pbutil.Buffer) {
}

func (m *lock_result) Changed() {
	m.changed = true
}

func (m *lock_result) Ok() {
	m.ok = true
}

func (m *lock_result) send(sender sender.Sender) {
	sender.SendAll(m.list)
}

func (m *lock_result) Reset() {
	m.changed = false
	m.ok = false
	m.list = m.list[:0]
	m.bcList = m.bcList[:0]
}

func ReadMsgList(m *lock_result) []pbutil.Buffer {
	return m.list
}
