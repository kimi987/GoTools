package recovtimes

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

var duration time.Duration = time.Second * 10
var maxTimes uint64 = 3

func New() *RecoverTimes {
	return NewRecoverTimes(time.Time{}, duration, maxTimes)
}

// 测试
func TestNewRecoverTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := New()
	Ω(rt.StartTime()).Should(BeEquivalentTo(time.Time{}))
	Ω(rt.Duration()).Should(BeEquivalentTo(duration))
	Ω(rt.MaxTimes()).Should(BeEquivalentTo(maxTimes))
}

func TestRecoverTimes_SetStartRecoverTime(t *testing.T) {
	RegisterTestingT(t)

	rt := New()

	ctime := time.Now()
	rt.SetStartTime(ctime)
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime))
}

func TestRecoverTimes_Times(t *testing.T) {
	RegisterTestingT(t)

	rt := New()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	Ω(rt.Times(ctime.Add(time.Second))).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second * 1))).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second * 8))).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second * 10))).Should(BeEquivalentTo(1))
	Ω(rt.Times(ctime.Add(time.Second * 17))).Should(BeEquivalentTo(1))
	Ω(rt.Times(ctime.Add(time.Second * 20))).Should(BeEquivalentTo(2))
	Ω(rt.Times(ctime.Add(time.Second * 33))).Should(BeEquivalentTo(3))
	Ω(rt.Times(ctime.Add(time.Second * 330))).Should(BeEquivalentTo(3))
}

func TestRecoverTimes_floatTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := New()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	Ω(rt.floatTimes(ctime)).Should(BeEquivalentTo(0.0))
	Ω(rt.floatTimes(ctime.Add(time.Second))).Should(BeEquivalentTo(0.1))
	Ω(rt.floatTimes(ctime.Add(time.Second * 8))).Should(BeEquivalentTo(0.8))
	Ω(rt.floatTimes(ctime.Add(time.Second * 10))).Should(BeEquivalentTo(1))
	Ω(rt.floatTimes(ctime.Add(time.Second * 17))).Should(BeEquivalentTo(1.7))
	Ω(rt.floatTimes(ctime.Add(time.Second * 20))).Should(BeEquivalentTo(2.0))
	Ω(rt.floatTimes(ctime.Add(time.Second * 33))).Should(BeEquivalentTo(3))
	Ω(rt.floatTimes(ctime.Add(time.Second * 330))).Should(BeEquivalentTo(3))
}

func TestRecoverTimes_HasEnoughTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := New()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	Ω(rt.HasEnoughTimes(1, ctime)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(1, ctime.Add(time.Second))).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(1, ctime.Add(time.Second*11))).Should(BeEquivalentTo(true))
	Ω(rt.HasEnoughTimes(2, ctime.Add(time.Second*11))).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(2, ctime.Add(time.Second*20))).Should(BeEquivalentTo(true))
}

func TestRecoverTimes_ReduceOneTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := New()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	Ω(rt.HasEnoughTimes(1, ctime)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(1, ctime.Add(time.Second))).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(1, ctime.Add(time.Second*11))).Should(BeEquivalentTo(true))

	// 有1次
	rt.ReduceOneTimes(ctime.Add(time.Second * 11)) // 减少掉一次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 10)))

	Ω(rt.Times(ctime.Add(time.Second * 11))).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second * 20))).Should(BeEquivalentTo(1))

	// 有2次
	rt.ReduceOneTimes(ctime.Add(time.Second * 33)) // 减少掉一次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 20)))

	Ω(rt.Times(ctime.Add(time.Second * 11))).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second * 20))).Should(BeEquivalentTo(0))

	// 有3次
	rt.ReduceOneTimes(ctime.Add(time.Second * 51)) // 减少掉一次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 31)))
	Ω(rt.Times(ctime.Add(time.Second * 51))).Should(BeEquivalentTo(2))
	Ω(rt.floatTimes(ctime.Add(time.Second * 51))).Should(BeEquivalentTo(2))
}

func TestRecoverTimes_ReduceTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := New()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	Ω(rt.HasEnoughTimes(2, ctime)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(2, ctime.Add(time.Second))).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(2, ctime.Add(time.Second*21))).Should(BeEquivalentTo(true))

	// 有2次
	rt.ReduceTimes(2, ctime.Add(time.Second*21)) // 减少掉2次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 20)))

	Ω(rt.Times(ctime.Add(time.Second * 21))).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second * 30))).Should(BeEquivalentTo(1))

	// 有3次
	rt.ReduceTimes(2, ctime.Add(time.Second*51)) // 减少掉2次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 41)))
	Ω(rt.Times(ctime.Add(time.Second * 51))).Should(BeEquivalentTo(1))
	Ω(rt.floatTimes(ctime.Add(time.Second * 51))).Should(BeEquivalentTo(1))

	// 有3次
	rt.ReduceTimes(3, ctime.Add(time.Second*100)) // 减少掉3次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 100)))
	Ω(rt.floatTimes(ctime.Add(time.Second * 51))).Should(BeEquivalentTo(0))
	Ω(rt.floatTimes(ctime.Add(time.Second * 100))).Should(BeEquivalentTo(0))
	Ω(rt.floatTimes(ctime.Add(time.Second * 110))).Should(BeEquivalentTo(1))
}

func TestRecoverTimes_GiveTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := New()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	rt.AddTimes(2, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(2))
	rt.AddTimes(4, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(3))

	rt.ReduceTimes(2, ctime)

	rt.AddTimes(1, ctime.Add(duration))
	Ω(rt.Times(ctime.Add(duration))).Should(BeEquivalentTo(3))

	rt.AddTimes(1, ctime.Add(duration))
	Ω(rt.Times(ctime.Add(duration))).Should(BeEquivalentTo(3))
}

func TestRecoverTimes_ChangeDuration(t *testing.T) {
	RegisterTestingT(t)

	rt := New()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	rt.AddTimes(4, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(3))

	rt.ChangeDuration(duration/2, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(3))
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(-duration * 3 / 2)))

	// 变更回来
	rt.ChangeDuration(duration, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(3))
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(-duration * 3)))

	// 减掉2次
	rt.ReduceTimes(2, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(1))

	rt.ChangeDuration(duration/2, ctime.Add(duration))
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(duration))).Should(BeEquivalentTo(2))
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime))
}

func TestRecoverTimes_ChangeMaxTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := New()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	rt.AddTimes(4, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(3))

	rt.ChangeMaxTimes(6, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(3))

	rt.AddTimes(4, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(6))

	rt.ChangeMaxTimes(3, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(3))

	rt.ChangeMaxTimes(6, ctime)
	Ω(rt.Times(ctime)).Should(BeEquivalentTo(3))
}
