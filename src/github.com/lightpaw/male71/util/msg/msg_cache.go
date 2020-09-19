package msg

import (
	"github.com/lightpaw/pbutil"
	"sync"
	"time"
)

func NewMsgCache(expired time.Duration, timeSrv interface {
	CurrentTime() time.Time
}) *MsgCache {
	c := &MsgCache{}
	c.time = timeSrv
	c.msgInfos = make(map[uint64]*msgInfo)
	c.expired = expired

	return c
}

type MsgCache struct {
	time interface {
		CurrentTime() time.Time
	}

	lock sync.RWMutex

	msgInfos map[uint64]*msgInfo

	expired time.Duration
}

func NewMsg(msg pbutil.Buffer, expiredTime time.Time) *msgInfo {
	m := &msgInfo{}
	m.msg = msg
	m.expiredTime = expiredTime

	return m
}

type msgInfo struct {
	msg pbutil.Buffer

	expiredTime time.Time
}

func (c *MsgCache) Get(k uint64) (msg pbutil.Buffer) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	m := c.msgInfos[k]
	if m == nil {
		return
	}

	if c.time.CurrentTime().After(m.expiredTime) {
		return nil
	}
	return m.msg
}

func (c *MsgCache) Update(k uint64, f func() (result pbutil.Buffer)) (msg pbutil.Buffer) {
	c.lock.Lock()
	defer c.lock.Unlock()

	msg = f().Static()
	c.msgInfos[k] = NewMsg(msg, c.time.CurrentTime().Add(c.expired))
	return
}

func (c *MsgCache) GetOrUpdate(k uint64, f func() (result pbutil.Buffer)) (msg pbutil.Buffer) {
	if msg = c.Get(k); msg != nil {
		return
	}
	msg = c.Update(k, f)
	return
}

func (c *MsgCache) Disable(k uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.msgInfos, k)
}
