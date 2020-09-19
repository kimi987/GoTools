package entity

import (
	"testing"
	. "github.com/onsi/gomega"
	"time"
)

func TestNextTimes(t *testing.T) {
	RegisterTestingT(t)

	nt := newNextTimes()

	times, nextTime := nt.encodeClient()
	Ω(times).Should(Equal(int32(0)))
	Ω(nextTime).Should(Equal(int32(0)))
	encoded := nt.encodeServer()
	Ω(encoded).Should(Equal(uint64(0)))

	nnt := newNextTimes()
	nnt.unmarshal(encoded)
	Ω(nnt.times).Should(Equal(uint64(0)))
	Ω(nnt.nextTime.Unix()).Should(Equal(int64(0)))

	// 次数+1
	nt.IncreseTimes()
	times, nextTime = nt.encodeClient()
	Ω(times).Should(Equal(int32(1)))
	Ω(nextTime).Should(Equal(int32(0)))
	encoded = nt.encodeServer()

	nnt = newNextTimes()
	nnt.unmarshal(encoded)
	Ω(nnt.times).Should(Equal(uint64(1)))
	Ω(nnt.nextTime.Unix()).Should(Equal(int64(0)))

	// 设置时间
	now := time.Now()
	nt.SetTimes(10)
	nt.SetNextTime(now)
	times, nextTime = nt.encodeClient()
	Ω(times).Should(Equal(int32(10)))
	Ω(nextTime).Should(Equal(int32(now.Unix())))
	encoded = nt.encodeServer()

	nnt = newNextTimes()
	nnt.unmarshal(encoded)
	Ω(nnt.times).Should(Equal(uint64(10)))
	Ω(nnt.nextTime.Unix()).Should(Equal(now.Unix()))
}
