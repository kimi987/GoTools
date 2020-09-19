package timeutil

import (
	"testing"
	. "github.com/onsi/gomega"
	"time"
	"fmt"
)

func TestBuildWeekDurTime(t *testing.T) {
	RegisterTestingT(t)

	dur, _ := time.ParseDuration("20m")

	t1, err := BuildWeekDurTime("w1,20h30m", dur)
	Ω(err).Should(BeNil())
	Ω(t1.Week).Should(Equal(Mon))
	Ω(int(t1.Time.Minutes())).Should(Equal(20*60 + 30))

	t2, err := BuildWeekDurTime("w7,2h0m", dur)
	Ω(err).Should(BeNil())
	Ω(t2.Week).Should(Equal(Sun))
	Ω(int(t2.Time.Minutes())).Should(Equal(2 * 60))

	_, err = BuildWeekDurTime("w0,2h0m", dur)
	Ω(err).ShouldNot(BeNil())

	dur60m, _ := time.ParseDuration("60m")
	t3, err := BuildWeekDurTime("w3,23h0m", dur60m)
	Ω(err).Should(BeNil())
	Ω(t3.Week).Should(Equal(Wed))
	fmt.Println(t3.Time)
	Ω(int(t3.Time.Minutes())).Should(Equal(23 * 60))
	Ω(int(t3.Dur.Minutes())).Should(Equal(60))
}

func TestWeekDurTime_CurrWeekTime(t *testing.T) {
	RegisterTestingT(t)

	// 2018-01-01 10:00:00 Wednesday
	ctime := time.Date(2018, 1, 3, 10, 0, 0, 0, East8)
	dur, _ := time.ParseDuration("20m")

	t1, _ := BuildWeekDurTime("w7,20h30m", dur)
	ct := t1.CurrWeekTime(ctime)
	Ω(ct.Day()).Should(Equal(7))

	t2, _ := BuildWeekDurTime("w1,1h30m", dur)
	ct = t2.CurrWeekTime(ctime)
	Ω(ct.Day()).Should(Equal(1))

	t3, _ := BuildWeekDurTime("w3,1h30m", dur)
	ct = t3.CurrWeekTime(ctime)
	Ω(ct.Day()).Should(Equal(3))
}

func TestWeekDurTime_NextTime(t *testing.T) {
	RegisterTestingT(t)

	// 2018-01-01 10:00:00 Wednesday
	ctime := time.Date(2018, 1, 3, 10, 0, 0, 0, East8)
	dur, _ := time.ParseDuration("24h20m")

	t1, _ := BuildWeekDurTime("w1,2h30m", dur)
	start, end := t1.NextTime(ctime)
	Ω(start.Day()).Should(Equal(8))
	Ω(end.Day()).Should(Equal(9))

	t2, _ := BuildWeekDurTime("w2,2h30m", dur)
	start, end = t2.NextTime(ctime)
	Ω(start.Day()).Should(Equal(9))
	Ω(end.Day()).Should(Equal(10))

	t3, _ := BuildWeekDurTime("w3,12h30m", dur)
	start, end = t3.NextTime(ctime)
	Ω(start.Day()).Should(Equal(3))
	Ω(end.Day()).Should(Equal(4))

	dur60m, _ := time.ParseDuration("60m")
	t4, _ := BuildWeekDurTime("w3,23h0m", dur60m)
	start, end = t4.NextTime(ctime)
	fmt.Println(start)
	fmt.Println(end)
	Ω(start.Day()).Should(Equal(3))
	Ω(end.Day()).Should(Equal(4))
}

func TestWeekTime_Add(t *testing.T) {
	RegisterTestingT(t)

	dur, _ := time.ParseDuration("20m")

	t1 := &WeekTime{Week: Mon, Time: dur}
	dur1, _ := time.ParseDuration("24h20m")
	t1 = t1.Add(dur1)

	Ω(t1.Week).Should(Equal(Tue))

	dur2, _ := time.ParseDuration("144h0m")
	t1 = t1.Add(dur2)
	Ω(t1.Week).Should(Equal(Mon))

	dur3, _ := time.ParseDuration("72h0m")
	t1 = t1.Add(dur3)
	Ω(t1.Week).Should(Equal(Thu))

	dur4, _ := time.ParseDuration("72h0m")
	t1 = t1.Add(dur4)
	Ω(t1.Week).Should(Equal(Sun))

	dur5, _ := time.ParseDuration("2h0m")
	t1 = t1.Add(dur5)
	Ω(t1.Week).Should(Equal(Sun))
}

func TestWeekTime_DurationWith(t *testing.T) {
	RegisterTestingT(t)

	dur, _ := time.ParseDuration("24h20m")

	t1, _ := BuildWeekDurTime("w1,20h30m", dur)
	t2, _ := BuildWeekDurTime("w7,20h30m", dur)

	d := t1.DurationWithInSameWeek(t2.WeekTime)
	Ω(d.Hours()).Should(Equal(float64(24 * 6)))

	d = t2.DurationWithInSameWeek(t1.WeekTime)
	Ω(d.Hours()).Should(Equal(float64(24 * 6)))
}
