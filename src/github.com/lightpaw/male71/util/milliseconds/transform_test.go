package milliseconds

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	RegisterTestingT(t)

	tm := time.Now()
	Ω(Time(tm)).Should(Equal(tm.UnixNano() / 1e6))

	Ω(Duration(time.Millisecond)).Should(Equal(int64(1)))
	Ω(Duration(time.Second)).Should(Equal(PerSecond))
	Ω(Duration(time.Minute)).Should(Equal(PerMinute))
	Ω(Duration(time.Hour)).Should(Equal(PerHour))
	Ω(Duration(24 * time.Hour)).Should(Equal(PerDay))

	Ω(FromSecond(1.5)).Should(Equal(int64(1500)))
	Ω(FromMinute(1.5)).Should(Equal(PerSecond * 90))
	Ω(FromHour(1.5)).Should(Equal(PerMinute * 90))
	Ω(FromDay(0.5)).Should(Equal(PerHour * 12))
}
