package daily_amount

import (
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/u64"
)

func NewDailyAmount(dayCount uint64) *DailyAmount {
	return &DailyAmount{amounts: make([]uint64, dayCount)}
}

// 每日数据
type DailyAmount struct {
	amounts []uint64
}

func (da *DailyAmount) GetAmount(index int) uint64 {
	if index < 0 || index >= len(da.amounts) {
		return 0
	}
	return da.amounts[index]
}

func (da *DailyAmount) Amounts() []uint64 {
	return da.amounts
}

func (da *DailyAmount) Add(toAdd uint64) {
	if len(da.amounts) <= 0 {
		return
	}

	da.amounts[0] = u64.Plus(da.amounts[0], toAdd)
}

// 总计
func (da *DailyAmount) Total() (result uint64) {
	for _, amount := range da.amounts {
		result = u64.Plus(result, amount)
	}
	return
}

// 按日统计
func (da *DailyAmount) TotalByDay(day int) (result uint64) {
	for i := imath.Min(day, len(da.amounts)) - 1; i >= 0; i-- {
		result = u64.Plus(result, da.amounts[i])
	}
	return
}

// 重置
func (da *DailyAmount) ResetDaily(day int) {
	if day <= 0 {
		return
	}

	if day < len(da.amounts) {
		copy(da.amounts[day:], da.amounts)
	}

	for i := 0; i < imath.Min(day, len(da.amounts)); i++ {
		da.amounts[i] = 0
	}
}

func (da *DailyAmount) Unmarshal(amounts []uint64) {
	copy(da.amounts, amounts)
}
