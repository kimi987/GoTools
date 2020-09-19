package data

import (
	"github.com/lightpaw/male7/util/u64"
)

func ParseAmount(data string) (*Amount, error) {
	amt, pct, err := parseAmount(data)
	if err != nil {
		return nil, err
	}

	if amt == 0 && pct == 0 {
		return nil, nil
	}

	return &Amount{Amount: amt, Percent: pct}, nil
}

//gogen:config
type Amount struct {
	_       struct{} `proto:"shared_proto.AmountProto"`
	Amount  uint64   `validator:"uint"` // 数值
	Percent uint64   `validator:"uint"` // 千分比
}

func (t *Amount) IsZero() bool {
	if t != nil {
		return t.Amount == 0 && t.Percent == 0
	}
	return true
}

func (t *Amount) Reset() {
	t.Amount = 0
	t.Percent = 0
}

func (t *Amount) GetAmount() uint64 {
	if t != nil {
		return t.Amount
	}
	return 0
}

func TotalAmount(toAdds ...*Amount) uint64 {
	var amount, percent uint64
	for _, toAdd := range toAdds {
		if toAdd != nil {
			amount += toAdd.Amount
			percent += toAdd.Percent
		}
	}

	return totalAmount(amount, percent)
}

func (a *Amount) CalculateByPercent(amount uint64) uint64 {
	if a.Percent <= 0 {
		return 0
	}

	return u64.MultiCoef(amount, u64.Division2Float64(a.Percent, iPercentRate))
}

func (a *Amount) Calc(amount uint64) uint64 {
	if a.Percent <= 0 {
		return amount + a.Amount
	}

	return a.CalculateByPercent(amount) + a.Amount
}

// 计算出原始值
func (a *Amount) Return(amount uint64) uint64 {
	if a.Percent <= 0 {
		return u64.Sub(amount, a.Amount)
	}

	return a.ReturnByPercent(u64.Sub(amount, a.Amount))
}

// 计算出原始值
func (a *Amount) ReturnByPercent(amount uint64) uint64 {
	if a.Percent <= 0 {
		return 0
	}

	return uint64(u64.Division2Float64(amount, a.Percent)) * iPercentRate
}

func NewAmountBuilder() *AmountBuilder {
	return &AmountBuilder{}
}

type AmountBuilder struct {
	amount  uint64
	percent uint64
}

func (t *AmountBuilder) Amount() *Amount {
	return &Amount{Amount: t.amount, Percent: t.percent}
}

func (t *AmountBuilder) Reset() {
	t.amount = 0
	t.percent = 0
}

func (t *AmountBuilder) Add(amt, pct uint64) {
	t.amount += amt
	t.percent += pct
}

func (t *AmountBuilder) AddAmount(amt *Amount) {
	if amt != nil {
		t.amount += amt.Amount
		t.percent += amt.Percent
	}
}

func (t *AmountBuilder) AddPercent(pct float64) {
	if pct > 0 {
		t.percent += u64.MultiCoef(iPercentRate, pct)
	}
}

func (t *AmountBuilder) TotalAmount() uint64 {
	return totalAmount(t.amount, t.percent)
}

const iPercentRate uint64 = 1000

func totalAmount(amount, percent uint64) uint64 {
	if percent == 0 {
		return amount
	}

	percentAmount := u64.MultiCoef(amount, u64.Division2Float64(percent, iPercentRate))
	return amount + percentAmount
}
