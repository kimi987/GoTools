package concurrent

import (
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/concurrent"
	"github.com/lightpaw/pbutil"
	"time"
)

func newVersionBufferCacheMap(f func(uint64, uint64) (pbutil.Buffer, error), curTime func() time.Time, expireAfter time.Duration) *BufferCacheWithExpireTimeMap {
	m := &BufferCacheWithExpireTimeMap{
		version:     atomic.NewUint64(0),
		cacheMap:    Newcache_map(),
		f:           f,
		curTime:     curTime,
		expireAfter: expireAfter,
	}

	return m
}

type BufferCacheWithExpireTimeMap struct {
	version *atomic.Uint64

	cacheMap *cache_map

	f func(uint64, uint64) (pbutil.Buffer, error)

	curTime func() time.Time

	expireAfter time.Duration
}

func (m *BufferCacheWithExpireTimeMap) GetBuffer(key uint64) (pbutil.Buffer, error) {
	_, msg, err := m.GetVersionBuffer(key)
	return msg, err
}

func (m *BufferCacheWithExpireTimeMap) GetVersionBuffer(key uint64) (uint64, pbutil.Buffer, error) {

	cacheWithExpireTime, _ := m.cacheMap.Get(key)
	if cacheWithExpireTime == nil {
		cacheWithExpireTime = m.createAndSet(key, cacheWithExpireTime)
	} else if cacheWithExpireTime.Expired(m.curTime()) {
		cacheWithExpireTime = m.createAndSet(key, cacheWithExpireTime)
	}

	cache := cacheWithExpireTime.cache

	oc, err := cache.Get()
	if err != nil {
		m.Clear(key)
	}

	return cache.Version(), oc, err
}

func (m *BufferCacheWithExpireTimeMap) GetI32VersionBuffer(key uint64) (int32, pbutil.Buffer, error) {
	version, msg, err := m.GetVersionBuffer(key)
	return concurrent.I32Version(version), msg, err
}

func (m *BufferCacheWithExpireTimeMap) Clear(key uint64) {
	m.cacheMap.Remove(key)
}

func (m *BufferCacheWithExpireTimeMap) expireTime() time.Time {
	return m.curTime().Add(m.expireAfter)
}

func (m *BufferCacheWithExpireTimeMap) createAndSet(key uint64, oldValue *buffer_cache_with_expire_time) *buffer_cache_with_expire_time {
	return m.cacheMap.Upsert(key, newBufferCache(m.newFunc(key), m.expireTime), upsertFunc(oldValue))
}

func upsertFunc(oldValue *buffer_cache_with_expire_time) func(exist bool, valueInMap, newValue *buffer_cache_with_expire_time) *buffer_cache_with_expire_time {
	return func(exist bool, valueInMap, newValue *buffer_cache_with_expire_time) *buffer_cache_with_expire_time {
		if !exist {
			// 以前都不存在
			return newValue
		}

		if valueInMap != oldValue {
			return valueInMap
		}

		return newValue
	}
}

func (m *BufferCacheWithExpireTimeMap) newFunc(key uint64) func() (uint64, pbutil.Buffer, error) {
	return func() (uint64, pbutil.Buffer, error) {
		version := m.version.Inc()
		msg, err := m.f(key, version)
		return version, msg, err
	}
}
