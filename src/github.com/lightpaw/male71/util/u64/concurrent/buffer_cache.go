package concurrent

import (
	"github.com/lightpaw/male7/util/concurrent"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/pbutil"
	"time"
)

func newBufferCache(f func() (uint64, pbutil.Buffer, error), expireTime func() time.Time) *buffer_cache_with_expire_time {
	c := &buffer_cache_with_expire_time{}

	c.cache = concurrent.NewOnceBufferCache(func() (uint64, pbutil.Buffer, error) {
		c.expireTime = expireTime()
		return f()
	})

	return c
}

type buffer_cache_with_expire_time struct {
	cache      *concurrent.OnceBufferCache
	expireTime time.Time
}

func (c *buffer_cache_with_expire_time) Expired(ctime time.Time) bool {
	if timeutil.IsZero(c.expireTime) {
		return false
	}
	return ctime.After(c.expireTime)
}
