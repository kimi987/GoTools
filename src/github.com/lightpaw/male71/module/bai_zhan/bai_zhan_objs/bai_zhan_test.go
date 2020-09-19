package bai_zhan_objs

import (
	"testing"
	. "github.com/onsi/gomega"
	"github.com/lightpaw/male7/config/bai_zhan_data"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
)

func TestBaiZhanChallengeTimes(t *testing.T) {
	RegisterTestingT(t)

	dailyResetDuration := 4 * time.Hour
	dailyResetTime := timeutil.NewOffsetDailyTime(timeutil.DurationMarshal64(dailyResetDuration))

	d := &bai_zhan_data.BaiZhanMiscData{
		RecoverTimesTime: []time.Duration{4 * time.Hour, 12 * time.Hour, 18 * time.Hour, 22 * time.Hour},
		RecoverTimes:     []uint64{2, 2, 2, 2},
	}

	z := timeutil.DailyTime.PrevTime(time.Now())
	//fmt.Println(z)
	//fmt.Println(dailyResetTime.Duration(z))
	//fmt.Println(dailyResetTime.Duration(z.Add(4*time.Hour)))
	Ω(d.GetCanChallengeTimes(dailyResetTime.Duration(z.Add(4*time.Hour)), dailyResetDuration)).Should(Equal(uint64(2)))
	Ω(d.GetCanChallengeTimes(dailyResetTime.Duration(z.Add(11*time.Hour)), dailyResetDuration)).Should(Equal(uint64(2)))
	Ω(d.GetCanChallengeTimes(dailyResetTime.Duration(z.Add(12*time.Hour)), dailyResetDuration)).Should(Equal(uint64(4)))
	Ω(d.GetCanChallengeTimes(dailyResetTime.Duration(z.Add(17*time.Hour)), dailyResetDuration)).Should(Equal(uint64(4)))
	Ω(d.GetCanChallengeTimes(dailyResetTime.Duration(z.Add(18*time.Hour)), dailyResetDuration)).Should(Equal(uint64(6)))
	Ω(d.GetCanChallengeTimes(dailyResetTime.Duration(z.Add(21*time.Hour)), dailyResetDuration)).Should(Equal(uint64(6)))
	Ω(d.GetCanChallengeTimes(dailyResetTime.Duration(z.Add(22*time.Hour)), dailyResetDuration)).Should(Equal(uint64(8)))
	Ω(d.GetCanChallengeTimes(dailyResetTime.Duration(z.Add(24*time.Hour)), dailyResetDuration)).Should(Equal(uint64(8)))
	Ω(d.GetCanChallengeTimes(dailyResetTime.Duration(z.Add(3*time.Hour)), dailyResetDuration)).Should(Equal(uint64(8)))

}
