package concurrent

import (
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/pbutil"
)

func NewI64BufferCacheMap(f func(int64) (pbutil.Buffer, error)) *I64BufferCacheMap {
	return NewI64VersionBufferCacheMap(func(key int64, version uint64) (pbutil.Buffer, error) {
		return f(key)
	})
}

func NewI64VersionBufferCacheMap(f func(int64, uint64) (pbutil.Buffer, error)) *I64BufferCacheMap {
	m := &I64BufferCacheMap{
		version:  atomic.NewUint64(0),
		cacheMap: Newi64_cache_map(),
		f:        f,
	}

	return m
}

type I64BufferCacheMap struct {
	version *atomic.Uint64

	cacheMap *i64_cache_map

	f func(int64, uint64) (pbutil.Buffer, error)
}

func (m *I64BufferCacheMap) GetBuffer(key int64) (pbutil.Buffer, error) {
	_, msg, err := m.GetVersionBuffer(key)
	return msg, err
}

func (m *I64BufferCacheMap) GetVersionBuffer(key int64) (uint64, pbutil.Buffer, error) {

	cache, _ := m.cacheMap.Get(key)
	if cache == nil {
		cache = m.createAndSet(key)
	}

	oc, err := cache.Get()
	if err != nil {
		m.Clear(key)
	}

	return cache.Version(), oc, err
}

func (m *I64BufferCacheMap) GetI32VersionBuffer(key int64) (int32, pbutil.Buffer, error) {
	version, msg, err := m.GetVersionBuffer(key)
	return I32Version(version), msg, err
}

func (m *I64BufferCacheMap) Clear(key int64) {
	m.createAndSet(key)
}

func (m *I64BufferCacheMap) createAndSet(key int64) *OnceBufferCache {
	c := NewOnceBufferCache(m.newFunc(key))
	m.cacheMap.Set(key, c)
	return c
}

func (m *I64BufferCacheMap) newFunc(key int64) func() (uint64, pbutil.Buffer, error) {
	return func() (uint64, pbutil.Buffer, error) {
		version := m.version.Inc()
		msg, err := m.f(key, version)
		return version, msg, err
	}
}
