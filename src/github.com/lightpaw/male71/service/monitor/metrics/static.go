package metrics

import "github.com/lightpaw/male7/util/atomic"

var panicCounter = atomic.NewInt64(0)

func IncPanic() {
	panicCounter.Inc()
}

func GetIncPanicCount() int64 {
	c := panicCounter.Load()
	if c > 0 {
		panicCounter.Add(-c)
	}
	return c
}
