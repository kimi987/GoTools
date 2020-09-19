package nsds

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestRetryDuration(t *testing.T) {
	RegisterTestingT(t)

	Ω(getRetryDuration(0)).Should(Equal(200 * time.Millisecond))
	Ω(getRetryDuration(1)).Should(Equal(300 * time.Millisecond))
	Ω(getRetryDuration(2)).Should(Equal(400 * time.Millisecond))
	Ω(getRetryDuration(3)).Should(Equal(500 * time.Millisecond))
	Ω(getRetryDuration(4)).Should(Equal(500 * time.Millisecond))
	Ω(getRetryDuration(5)).Should(Equal(500 * time.Millisecond))
	Ω(getRetryDuration(100)).Should(Equal(500 * time.Millisecond))
}

func TestRetry(t *testing.T) {
	RegisterTestingT(t)

	counter := 0

	DoBackoffRetry(func() (isRetry bool) {
		counter++
		return counter < 2
	})

	Ω(counter).Should(Equal(2))
}
