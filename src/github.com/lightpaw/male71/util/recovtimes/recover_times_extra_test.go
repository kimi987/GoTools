package recovtimes

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

//var duration time.Duration = time.Second * 10
//var maxTimes uint64 = 3

func NewExtra() *ExtraRecoverTimes {
	return NewExtraRecoverTimes(time.Time{}, duration, maxTimes)
}

// 测试
func TestNewExtraRecoverTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := NewExtra()
	Ω(rt.StartTime()).Should(BeEquivalentTo(time.Time{}))
	Ω(rt.Duration()).Should(BeEquivalentTo(duration))
	Ω(rt.MaxTimes(1)).Should(BeEquivalentTo(maxTimes + 1))
}

func TestExtraRecoverTimes_SetStartRecoverTime(t *testing.T) {
	RegisterTestingT(t)

	rt := NewExtra()

	ctime := time.Now()
	rt.SetStartTime(ctime)
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime))
}

func TestExtraRecoverTimes_Times(t *testing.T) {
	RegisterTestingT(t)

	rt := NewExtra()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	Ω(rt.Times(ctime.Add(time.Second), 0)).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second*1), 1)).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second*8), 2)).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second*10), 0)).Should(BeEquivalentTo(1))
	Ω(rt.Times(ctime.Add(time.Second*17), 1)).Should(BeEquivalentTo(1))
	Ω(rt.Times(ctime.Add(time.Second*20), 0)).Should(BeEquivalentTo(2))
	Ω(rt.Times(ctime.Add(time.Second*33), 0)).Should(BeEquivalentTo(3))
	Ω(rt.Times(ctime.Add(time.Second*39), 1)).Should(BeEquivalentTo(3))
	Ω(rt.Times(ctime.Add(time.Second*40), 1)).Should(BeEquivalentTo(4))
	Ω(rt.Times(ctime.Add(time.Second*330), 2)).Should(BeEquivalentTo(5))
}

func TestExtraRecoverTimes_floatTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := NewExtra()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	Ω(rt.floatTimes(ctime, 0)).Should(BeEquivalentTo(0.0))
	Ω(rt.floatTimes(ctime.Add(time.Second), 1)).Should(BeEquivalentTo(0.1))
	Ω(rt.floatTimes(ctime.Add(time.Second*8), 2)).Should(BeEquivalentTo(0.8))
	Ω(rt.floatTimes(ctime.Add(time.Second*10), 0)).Should(BeEquivalentTo(1))
	Ω(rt.floatTimes(ctime.Add(time.Second*17), 1)).Should(BeEquivalentTo(1.7))
	Ω(rt.floatTimes(ctime.Add(time.Second*20), 0)).Should(BeEquivalentTo(2.0))
	Ω(rt.floatTimes(ctime.Add(time.Second*33), 0)).Should(BeEquivalentTo(3))
	Ω(rt.floatTimes(ctime.Add(time.Second*33), 2)).Should(BeEquivalentTo(3.3))
	Ω(rt.floatTimes(ctime.Add(time.Second*330), 2)).Should(BeEquivalentTo(5))
}

func TestExtraRecoverTimes_HasEnoughTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := NewExtra()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	Ω(rt.HasEnoughTimes(1, ctime, 0)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(1, ctime.Add(time.Second), 0)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(1, ctime.Add(time.Second*11), 0)).Should(BeEquivalentTo(true))
	Ω(rt.HasEnoughTimes(2, ctime.Add(time.Second*11), 0)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(2, ctime.Add(time.Second*20), 0)).Should(BeEquivalentTo(true))
	Ω(rt.HasEnoughTimes(4, ctime.Add(time.Second*40), 0)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(4, ctime.Add(time.Second*40), 1)).Should(BeEquivalentTo(true))
}

func TestExtraRecoverTimes_ReduceOneTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := NewExtra()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	Ω(rt.HasEnoughTimes(1, ctime, 0)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(1, ctime.Add(time.Second), 0)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(1, ctime.Add(time.Second*11), 0)).Should(BeEquivalentTo(true))

	// 有1次
	rt.ReduceOneTimes(ctime.Add(time.Second*11), 0) // 减少掉一次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 10)))

	Ω(rt.Times(ctime.Add(time.Second*11), 0)).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second*20), 0)).Should(BeEquivalentTo(1))

	// 有2次
	rt.ReduceOneTimes(ctime.Add(time.Second*33), 0) // 减少掉一次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 20)))

	Ω(rt.Times(ctime.Add(time.Second*11), 0)).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second*20), 0)).Should(BeEquivalentTo(0))

	// 有3次
	rt.ReduceOneTimes(ctime.Add(time.Second*51), 0) // 减少掉一次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 31)))
	Ω(rt.Times(ctime.Add(time.Second*51), 0)).Should(BeEquivalentTo(2))
	Ω(rt.floatTimes(ctime.Add(time.Second*51), 0)).Should(BeEquivalentTo(2))

	// 有4次
	rt.ReduceOneTimes(ctime.Add(time.Second*81), 1) // 减少掉一次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 51)))
	Ω(rt.Times(ctime.Add(time.Second*81), 0)).Should(BeEquivalentTo(3))
	Ω(rt.floatTimes(ctime.Add(time.Second*81), 1)).Should(BeEquivalentTo(3))
}

func TestExtraRecoverTimes_ReduceTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := NewExtra()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	Ω(rt.HasEnoughTimes(2, ctime, 0)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(2, ctime.Add(time.Second), 0)).Should(BeEquivalentTo(false))
	Ω(rt.HasEnoughTimes(2, ctime.Add(time.Second*21), 0)).Should(BeEquivalentTo(true))

	// 有2次
	rt.ReduceTimes(2, ctime.Add(time.Second*21), 0) // 减少掉2次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 20)))

	Ω(rt.Times(ctime.Add(time.Second*21), 0)).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(time.Second*30), 0)).Should(BeEquivalentTo(1))

	// 有3次
	rt.ReduceTimes(2, ctime.Add(time.Second*51), 0) // 减少掉2次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 41)))
	Ω(rt.Times(ctime.Add(time.Second*51), 0)).Should(BeEquivalentTo(1))
	Ω(rt.floatTimes(ctime.Add(time.Second*51), 0)).Should(BeEquivalentTo(1))

	// 有3次
	rt.ReduceTimes(3, ctime.Add(time.Second*100), 0) // 减少掉3次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 100)))
	Ω(rt.floatTimes(ctime.Add(time.Second*51), 0)).Should(BeEquivalentTo(0))
	Ω(rt.floatTimes(ctime.Add(time.Second*100), 0)).Should(BeEquivalentTo(0))
	Ω(rt.floatTimes(ctime.Add(time.Second*110), 0)).Should(BeEquivalentTo(1))

	// 有4次
	rt.ReduceTimes(3, ctime.Add(time.Second*150), 1) // 减少掉3次
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(time.Second * 140)))
	Ω(rt.floatTimes(ctime.Add(time.Second*51), 0)).Should(BeEquivalentTo(0))
	Ω(rt.floatTimes(ctime.Add(time.Second*100), 1)).Should(BeEquivalentTo(0))
	Ω(rt.floatTimes(ctime.Add(time.Second*150), 0)).Should(BeEquivalentTo(1))
}

func TestExtraRecoverTimes_GiveTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := NewExtra()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	rt.AddTimes(2, ctime, 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(2))
	rt.AddTimes(4, ctime, 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(3))

	rt.ReduceTimes(2, ctime, 0)

	rt.AddTimes(1, ctime.Add(duration), 0)
	Ω(rt.Times(ctime.Add(duration), 0)).Should(BeEquivalentTo(3))

	rt.AddTimes(1, ctime.Add(duration), 0)
	Ω(rt.Times(ctime.Add(duration), 0)).Should(BeEquivalentTo(3))

	rt.AddTimes(1, ctime.Add(duration), 0)
	Ω(rt.Times(ctime.Add(duration), 0)).Should(BeEquivalentTo(3))

	rt.AddTimes(1, ctime.Add(duration), 1)
	Ω(rt.Times(ctime.Add(duration), 1)).Should(BeEquivalentTo(4))

	// 已满不会再涨
	rt.AddTimes(1, ctime.Add(duration), 1)
	Ω(rt.Times(ctime.Add(duration), 1)).Should(BeEquivalentTo(4))

	// 加大上限，再涨
	rt.AddTimes(1, ctime.Add(duration), 3)
	Ω(rt.Times(ctime.Add(duration), 3)).Should(BeEquivalentTo(5))
}

func TestExtraRecoverTimes_ChangeDuration(t *testing.T) {
	RegisterTestingT(t)

	rt := NewExtra()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	rt.AddTimes(4, ctime, 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(3))

	rt.ChangeDuration(duration/2, ctime, 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(3))
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(-duration * 3 / 2)))

	// 变更回来
	rt.ChangeDuration(duration, ctime, 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(3))
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(-duration * 3)))

	// 减掉2次
	rt.ReduceTimes(2, ctime, 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(1))

	rt.ChangeDuration(duration/2, ctime.Add(duration), 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(0))
	Ω(rt.Times(ctime.Add(duration), 0)).Should(BeEquivalentTo(2))
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime))

	// 现在有2次次数，加2次，再将duration变回来
	rt.AddTimes(2, ctime.Add(duration), 1)

	rt.ChangeDuration(duration, ctime.Add(duration), 1)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(3))
	Ω(rt.Times(ctime.Add(duration), 0)).Should(BeEquivalentTo(3))
	Ω(rt.Times(ctime.Add(duration), 1)).Should(BeEquivalentTo(4))
	Ω(rt.StartTime()).Should(BeEquivalentTo(ctime.Add(-3 * duration)))
}

func TestExtraRecoverTimes_ChangeMaxTimes(t *testing.T) {
	RegisterTestingT(t)

	rt := NewExtra()

	ctime := time.Now()
	rt.SetStartTime(ctime)

	rt.AddTimes(4, ctime, 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(3))
	Ω(rt.Times(ctime, 1)).Should(BeEquivalentTo(3))

	rt.ChangeMaxTimes(6, ctime, 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(3))

	rt.AddTimes(4, ctime, 1)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(6))
	Ω(rt.Times(ctime, 2)).Should(BeEquivalentTo(7))

	rt.ChangeMaxTimes(3, ctime, 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(3))
	Ω(rt.Times(ctime, 1)).Should(BeEquivalentTo(3))

	rt.ChangeMaxTimes(6, ctime, 0)
	Ω(rt.Times(ctime, 0)).Should(BeEquivalentTo(3))
}
