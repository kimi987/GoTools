package daily_amount

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestDaily(t *testing.T) {
	RegisterTestingT(t)

	dailyAmount := NewDailyAmount(2)
	dailyAmount.Add(10)
	Ω(dailyAmount.amounts).Should(BeEquivalentTo([]uint64{10, 0}))

	dailyAmount.Add(10)
	Ω(dailyAmount.amounts).Should(BeEquivalentTo([]uint64{20, 0}))

	dailyAmount.ResetDaily(1)

	Ω(dailyAmount.amounts).Should(BeEquivalentTo([]uint64{0, 20}))

	dailyAmount.Add(10)
	Ω(dailyAmount.amounts).Should(BeEquivalentTo([]uint64{10, 20}))

	Ω(dailyAmount.GetAmount(0)).Should(BeEquivalentTo(10))
	Ω(dailyAmount.GetAmount(1)).Should(BeEquivalentTo(20))

	dailyAmount.ResetDaily(2)
	Ω(dailyAmount.amounts).Should(BeEquivalentTo([]uint64{0, 0}))
}
