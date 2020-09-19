package concurrent

import (
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/pbutil"
)

func NewU64BufferCacheMap(f func(uint64) (pbutil.Buffer, error)) *U64BufferCacheMap {
	return NewU64VersionBufferCacheMap(func(key uint64, version uint64) (pbutil.Buffer, error) {
		return f(key)
	})
}

func NewU64VersionBufferCacheMap(f func(uint64, uint64) (pbutil.Buffer, error)) *U64BufferCacheMap {
	m := &U64BufferCacheMap{
		version:  atomic.NewUint64(0),
		cacheMap: Newu64_cache_map(),
		f:        f,
	}

	return m
}

type U64BufferCacheMap struct {
	version *atomic.Uint64

	cacheMap *u64_cache_map

	f func(uint64, uint64) (pbutil.Buffer, error)
}

func (m *U64BufferCacheMap) GetBuffer(key uint64) (pbutil.Buffer, error) {
	_, msg, err := m.GetVersionBuffer(key)
	return msg, err
}

func (m *U64BufferCacheMap) GetVersionBuffer(key uint64) (uint64, pbutil.Buffer, error) {

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

func (m *U64BufferCacheMap) GetI32VersionBuffer(key uint64) (int32, pbutil.Buffer, error) {
	version, msg, err := m.GetVersionBuffer(key)
	return I32Version(version), msg, err
}

func (m *U64BufferCacheMap) Clear(key uint64) {
	m.createAndSet(key)
}

func (m *U64BufferCacheMap) createAndSet(key uint64) *OnceBufferCache {
	c := NewOnceBufferCache(m.newFunc(key))
	m.cacheMap.Set(key, c)
	return c
}

func (m *U64BufferCacheMap) newFunc(key uint64) func() (uint64, pbutil.Buffer, error) {
	return func() (uint64, pbutil.Buffer, error) {
		version := m.version.Inc()
		msg, err := m.f(key, version)
		return version, msg, err
	}
}
