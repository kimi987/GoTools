package nsds

import (
	"time"
)

func DoBackoffRetry(f func() (isRetry bool)) {

	for i := 0; i < maxRetryTimes; i++ {
		if i > 0 {
			// 等
			time.Sleep(getRetryDuration(i - 1))
		}

		if !f() {
			return
		}
	}
}

const (
	maxRetryTimes      = 10
	initRetryDuration  = 200 * time.Millisecond
	increRetryDuration = 100 * time.Millisecond
	maxRetryDuration   = 500 * time.Millisecond
)

func getRetryDuration(times int) time.Duration {
	d := initRetryDuration

	if times > 0 {
		d += time.Duration(times) * increRetryDuration
	}

	if d > maxRetryDuration {
		d = maxRetryDuration
	}

	return d
}
