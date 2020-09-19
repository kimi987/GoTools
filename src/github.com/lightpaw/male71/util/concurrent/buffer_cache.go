package concurrent

import (
	atomic2 "github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/pbutil"
	"math"
	"sync"
	"sync/atomic"
)

type BufferCache interface {
	Get() (oc pbutil.Buffer, err error)
	GetWithVersion() (vsn uint64, oc pbutil.Buffer, err error)
	Clear()
}

func NewBufferCache(f func() (pbutil.Buffer, error)) BufferCache {
	return NewVersionBufferCache(func(uint64) (pbutil.Buffer, error) {
		return f()
	})
}

func NewVersionBufferCache(f func(uint64) (pbutil.Buffer, error)) BufferCache {
	version := atomic2.NewUint64(0)
	c := &buffer_cache{cache: &atomic.Value{}, f: func() (uint64, pbutil.Buffer, error) {
		vsn := version.Inc()
		msg, err := f(vsn)
		return vsn, msg, err
	}}
	c.Clear()

	return c
}

type buffer_cache struct {
	cache *atomic.Value
	f     func() (uint64, pbutil.Buffer, error)
}

func (bc *buffer_cache) Get() (oc pbutil.Buffer, err error) {
	_, oc, err = bc.GetWithVersion()
	return
}

func (bc *buffer_cache) GetWithVersion() (vsn uint64, oc pbutil.Buffer, err error) {
	c := bc.cache.Load().(*OnceBufferCache)
	oc, err = c.Get()
	if err != nil {
		bc.Clear()
	}

	return c.Version(), oc, err
}

func (bc *buffer_cache) Clear() {
	bc.cache.Store(NewOnceBufferCache(bc.f))
}

func NewOnceBufferCache(f func() (uint64, pbutil.Buffer, error)) *OnceBufferCache {
	return &OnceBufferCache{f: f}
}

type OnceBufferCache struct {
	version uint64
	v       pbutil.Buffer
	err     error
	once    sync.Once

	f func() (uint64, pbutil.Buffer, error)
}

func (c *OnceBufferCache) Version() uint64 {
	return c.version
}

func (c *OnceBufferCache) I32Version() int32 {
	return I32Version(c.version)
}

func I32Version(version uint64) int32 {
	return int32(version & math.MaxInt32)
}

func (c *OnceBufferCache) build() {
	c.version, c.v, c.err = c.f()
	if c.v != nil {
		c.v = c.v.Static()
	}
}

func (c *OnceBufferCache) Get() (pbutil.Buffer, error) {
	c.once.Do(c.build)
	return c.v, c.err
}
