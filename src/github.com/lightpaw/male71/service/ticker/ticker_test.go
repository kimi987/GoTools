package ticker

import (
	"github.com/lightpaw/male7/config/singleton"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestTick(t *testing.T) {
	RegisterTestingT(t)

	ctime := time.Now()
	nextDailyResetTime := singleton.GetNextResetDailyTime(ctime, 0)
	dailyTicker := NewTicker(ctime, nextDailyResetTime.Sub(ctime), 24*time.Hour)

	Ω(dailyTicker.GetTickTime().GetTickTime().Unix()).Should(Equal(nextDailyResetTime.Unix()))
	Ω(dailyTicker.GetTickTime().GetPrevTickTime().Unix()).Should(Equal(nextDailyResetTime.Add(-24 * time.Hour).Unix()))

	tt := nextDailyResetTime.Add(5 * time.Hour)
	nextCtime := tt.Add(24 * time.Hour)
	for i := 0; i < 23; i++ {
		d := time.Duration(i) * time.Hour

		if i == 5 {
			Ω(singleton.GetNextResetDailyTime(tt, d).Equal(nextCtime)).Should(BeTrue())
		} else {
			Ω(singleton.GetNextResetDailyTime(tt, d).Before(nextCtime)).Should(BeTrue())
		}
	}

}
