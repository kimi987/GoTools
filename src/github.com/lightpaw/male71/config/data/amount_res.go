package data

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/transform"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func ParseResAmount(data string) (*ResAmount, error) {
	i1 := strings.Index(data, ":")
	if i1 < 0 {
		return nil, errors.Errorf("ResType解析失败，正确格式: gold:100+10%% 类型[%v], %s", transform.EnumMapKeys(shared_proto.ResType_value), data)
	}
	typ := data[0:i1]
	it, ok := shared_proto.ResType_value[strings.ToUpper(typ)]
	if !ok {
		return nil, errors.Errorf("ResType解析失败，正确格式: gold:100+10%% 类型[%v], %s", transform.EnumMapKeys(shared_proto.ResType_value), data)
	}

	amt, pct, err := parseAmount(data[i1+1:])
	if err != nil {
		return nil, errors.Wrapf(err, "ResType解析失败，正确格式: gold:100+10%% 类型[%v], %s", transform.EnumMapKeys(shared_proto.ResType_value), data)
	}

	return &ResAmount{Type: shared_proto.ResType(it), Amount: amt, Percent: pct}, nil
}

func parseAmount(data string) (uint64, uint64, error) {
	if len(data) == 0 {
		return 0, 0, nil
	}

	i2 := strings.Index(data, "+")
	if i2 > 0 {
		if !strings.HasSuffix(data, "%") {
			return 0, 0, errors.Errorf("Amount解析失败，正确格式: 100+10%% (有+号后面必须以%%结尾) 类型[%v], %s", data)
		}

		amount, err := strconv.ParseUint(data[:i2], 10, 64)
		if err != nil {
			return 0, 0, errors.Errorf("Amount解析失败，正确格式: 100+10%% 类型[%v], %s", data)
		}

		percent, err := strconv.ParseUint(data[i2+1:len(data)-1], 10, 64)
		if err != nil {
			return 0, 0, errors.Errorf("Amount解析失败，正确格式: 100+10%% 类型[%v], %s", data)
		}

		return amount, percent, nil
	}

	if strings.HasSuffix(data, "%") {
		percent, err := strconv.ParseUint(data[:len(data)-1], 10, 64)
		if err != nil {
			return 0, 0, errors.Errorf("Amount解析失败，正确格式: 100+10%% 类型[%v], %s", data)
		}
		return 0, percent, nil
	} else {
		amount, err := strconv.ParseUint(data, 10, 64)
		if err != nil {
			return 0, 0, errors.Errorf("Amount解析失败，正确格式: 100+10%% 类型[%v], %s", data)
		}
		return amount, 0, nil
	}

}

//
//type Amount struct {
//	Amount  int // 数值
//	Percent int // 万分比
//}
//
//func (t *Amount) Reset() {
//	t.Amount = 0
//	t.Percent = 0
//}
//
//func (t *Amount) TotalAmount() int {
//	if t.Percent == 0 {
//		return t.Amount
//	}
//
//	if t.Percent <= -iPercentRate {
//		// 最多减少到0
//		return 0
//	}
//
//	percentAmount := float64(t.Amount) * (float64(t.Percent) / fPercentRate)
//	return t.Amount + int(percentAmount)
//}

//gogen:config
type ResAmount struct {
	_       struct{}             `proto:"shared_proto.ResProto"`
	Type    shared_proto.ResType `type:"enum"`
	Amount  uint64               `validator:"uint"` // 数值
	Percent uint64               `validator:"uint"` // 万分比
}

func (t *ResAmount) Reset() {
	t.Amount = 0
	t.Percent = 0
}

func NewResAmountBuilder(typ shared_proto.ResType) *ResAmountBuilder {
	return &ResAmountBuilder{Type: typ}
}

func NewResBuilderMap() map[shared_proto.ResType]*ResAmountBuilder {
	resBdMap := make(map[shared_proto.ResType]*ResAmountBuilder, len(shared_proto.ResType_value))
	for _, v := range shared_proto.ResType_value {
		t := shared_proto.ResType(v)
		resBdMap[t] = NewResAmountBuilder(t)
	}

	return resBdMap
}

type ResAmountBuilder struct {
	Type shared_proto.ResType
	AmountBuilder
}

func (t *ResAmountBuilder) Add(amount *ResAmount) {
	if t.Type == amount.Type {
		t.AmountBuilder.Add(amount.Amount, amount.Percent)
	}
}
