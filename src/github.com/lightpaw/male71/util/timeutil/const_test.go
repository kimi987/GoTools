package timeutil

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestDailyTime(t *testing.T) {
	RegisterTestingT(t)

	ctime := time.Now()

	prevTime := DailyTime.PrevTime(ctime)
	Ω(prevTime.Nanosecond()).Should(Equal(0))
	Ω(prevTime.Second()).Should(Equal(0))
	Ω(prevTime.Minute()).Should(Equal(0))

	Ω(DailyTime.NextTime(ctime)).Should(Equal(prevTime.Add(24 * time.Hour)))

	ctime = time.Date(2001, 1, 1, 2, 0, 0, 0, GameZone)
	Ω(DailyTime.PrevTime(ctime).UnixNano()).Should(Equal(time.Date(2001, 1, 1, 0, 0, 0, 0, GameZone).UnixNano()))
	Ω(DailyTime.NextTime(ctime).UnixNano()).Should(Equal(time.Date(2001, 1, 2, 0, 0, 0, 0, GameZone).UnixNano()))

	Ω(DailyTime.Duration(ctime)).Should(Equal(2 * time.Hour))
}

func TestNextWeekTime(t *testing.T) {
	RegisterTestingT(t)

	monday := time.Date(2017, 8, 14, 0, 0, 0, 0, GameZone)
	Ω(monday.Weekday()).Should(Equal(time.Monday))

	nextWeekExcepts := [...]time.Time{
		time.Date(2017, 8, 20, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 14, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 15, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 16, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 17, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 18, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 19, 0, 0, 0, 0, GameZone),
	}

	for i := 0; i < 7; i++ {
		// nextWeekTime
		Ω(NextWeekTime(monday, time.Weekday(i)).UnixNano()).Should(Equal(nextWeekExcepts[i].UnixNano()))
	}

}

func TestWeekCycleTime(t *testing.T) {
	RegisterTestingT(t)

	monday := time.Date(2017, 8, 14, 0, 0, 0, 0, GameZone)
	Ω(monday.Weekday()).Should(Equal(time.Monday))

	prevExcepts := [...]time.Time{
		time.Date(2017, 8, 13, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 14, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 8, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 9, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 10, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 11, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 12, 0, 0, 0, 0, GameZone),
	}

	nextExcepts := [...]time.Time{
		time.Date(2017, 8, 20, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 21, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 15, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 16, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 17, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 18, 0, 0, 0, 0, GameZone),
		time.Date(2017, 8, 19, 0, 0, 0, 0, GameZone),
	}

	for i := 0; i < 7; i++ {
		// prev time
		Ω(WeekCycleTime(time.Weekday(i)).PrevTime(monday).UnixNano()).Should(Equal(prevExcepts[i].UnixNano()))

		Ω(WeekCycleTime(time.Weekday(i)).NextTime(monday).UnixNano()).Should(Equal(nextExcepts[i].UnixNano()))
	}

}

func TestParseHMS(t *testing.T) {
	RegisterTestingT(t)

	taa := []struct {
		value string
		hh    int
		mm    int
		ss    int
		err   bool
	}{
		{"-1", 0, 0, 0, true},
		{"24", 0, 0, 0, true},
		{"0:-1", 0, 0, 0, true},
		{"0:60", 0, 0, 0, true},
		{"0:0:-1", 0, 0, 0, true},
		{"0:0:60", 0, 0, 0, true},

		{"", 0, 0, 0, false},
		{"0", 0, 0, 0, false},
		{"0:0", 0, 0, 0, false},
		{"0:0:0", 0, 0, 0, false},
		{"23", 23, 0, 0, false},
		{"23:59", 23, 59, 0, false},
		{"23:59:59", 23, 59, 59, false},

		{"05:08:09", 5, 8, 9, false},
	}

	for _, ta := range taa {
		hh, mm, ss, err := ParseHMS(ta.value)
		if ta.err {
			Ω(err).Should(HaveOccurred())
		} else {
			Ω(err).Should(Succeed())
		}

		Ω(hh).Should(Equal(hh))
		Ω(mm).Should(Equal(mm))
		Ω(ss).Should(Equal(ss))
	}

}
