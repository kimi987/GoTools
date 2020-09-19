package concurrent

import (
	"github.com/lightpaw/male7/util/concurrent"
	"github.com/lightpaw/pbutil"
	"time"
)

func NewBufferCacheMapBuilder(f func(key int64) (pbutil.Buffer, error)) *BufferCacheMapBuilder {
	return NewVersionBufferCacheMapBuilder(func(key int64, version uint64) (pbutil.Buffer, error) {
		return f(key)
	})
}

func NewVersionBufferCacheMapBuilder(f Func) *BufferCacheMapBuilder {
	return &BufferCacheMapBuilder{
		f: f,
	}
}

type Func func(key int64, version uint64) (pbutil.Buffer, error)

type BufferCacheMapBuilder struct {
	f Func

	// 带过期时间的
	timeProvider func() time.Time
	expireAfter  time.Duration
}

func (b *BufferCacheMapBuilder) TimeProvider(timeProvider func() time.Time) *BufferCacheMapBuilder {
	b.timeProvider = timeProvider
	return b
}

func (b *BufferCacheMapBuilder) ExpireAfter(duration time.Duration) *BufferCacheMapBuilder {
	b.expireAfter = duration
	if b.timeProvider == nil {
		b.timeProvider = time.Now
	}
	return b
}

func (b *BufferCacheMapBuilder) Build() I64BufferMap {
	if b.expireAfter > 0 {
		// 有过期时间
		return newVersionBufferCacheMap(b.f, b.timeProvider, b.expireAfter)
	}

	return concurrent.NewI64VersionBufferCacheMap(b.f)
}
