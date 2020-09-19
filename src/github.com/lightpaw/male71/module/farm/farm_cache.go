package farm

import (
	"sync"
	"github.com/lightpaw/pbutil"
	"time"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/timer"
)

var (
	newestLogExpiredDuration, _ = time.ParseDuration("1s")
	logExpiredDuration, _       = time.ParseDuration("5s")
	relationExpiredDuration, _  = time.ParseDuration("5s")
)

type FarmCache struct {
	newestLogCache *sync.Map
	logCache       *sync.Map
	stealListCache *sync.Map

	timeService iface.TimeService
	ticker      iface.TickerService
}

func NewFarmCache(timeService iface.TimeService, ticker iface.TickerService) *FarmCache {
	m := &FarmCache{
		logCache:       &sync.Map{},
		stealListCache: &sync.Map{},
		newestLogCache: &sync.Map{},
		timeService:    timeService,
		ticker:         ticker,
	}

	go call.CatchLoopPanic(m.loop, "FarmCache.loop")

	return m
}

type CacheEntry struct {
	msg         pbutil.Buffer
	expiredTime time.Time
}

func (c *FarmCache) LoadNewestLog(heroId int64) pbutil.Buffer {
	obj, ok := c.newestLogCache.Load(heroId)
	if !ok || obj == nil {
		return nil
	}
	if entry, ok := obj.(*CacheEntry); ok {
		return entry.msg
	}
	return nil
}

func (c *FarmCache) LoadLog(heroId int64) pbutil.Buffer {
	obj, ok := c.logCache.Load(heroId)
	if !ok || obj == nil {
		return nil
	}
	if entry, ok := obj.(*CacheEntry); ok {
		return entry.msg
	}
	return nil
}

func (c *FarmCache) LoadStealList(heroId int64) pbutil.Buffer {
	obj, ok := c.stealListCache.Load(heroId)
	if !ok || obj == nil {
		return nil
	}
	if entry, ok := obj.(*CacheEntry); ok {
		return entry.msg
	}
	return nil
}

func (m *FarmCache) StoreNewestLog(heroId int64, buffer pbutil.Buffer) {
	ctime := m.timeService.CurrentTime()
	entry := &CacheEntry{msg: buffer, expiredTime: ctime.Add(newestLogExpiredDuration)}
	m.newestLogCache.Store(heroId, entry)
}

func (m *FarmCache) StoreLog(heroId int64, buffer pbutil.Buffer) {
	ctime := m.timeService.CurrentTime()
	entry := &CacheEntry{msg: buffer, expiredTime: ctime.Add(logExpiredDuration)}
	m.logCache.Store(heroId, entry)
}

func (m *FarmCache) StoreStealList(heroId int64, buffer pbutil.Buffer) {
	ctime := m.timeService.CurrentTime()
	entry := &CacheEntry{msg: buffer, expiredTime: ctime.Add(relationExpiredDuration)}
	m.stealListCache.Store(heroId, entry)
}

func (m *FarmCache) loop() {
	loopWheel := timer.NewTimingWheel(500 * time.Millisecond, 32)
	secondTick := loopWheel.After(time.Second)
	for {
		select {
		case <-secondTick:
			secondTick = loopWheel.After(time.Second)
			call.CatchPanic(m.removeExpiredCache, "农场，定时删除缓存")
		}
	}
}

func (m *FarmCache) removeExpiredCache() {
	ctime := m.timeService.CurrentTime()

	m.newestLogCache.Range(func(key, value interface{}) bool {
		if entry, ok := value.(*CacheEntry); ok {
			if entry.expiredTime.Before(ctime) {
				m.newestLogCache.Delete(key)
			}
		}
		return true
	})

	m.logCache.Range(func(key, value interface{}) bool {
		if entry, ok := value.(*CacheEntry); ok {
			if entry.expiredTime.Before(ctime) {
				m.logCache.Delete(key)
			}
		}
		return true
	})

	m.stealListCache.Range(func(key, value interface{}) bool {
		if entry, ok := value.(*CacheEntry); ok {
			if entry.expiredTime.Before(ctime) {
				m.stealListCache.Delete(key)
			}
		}
		return true
	})
}
