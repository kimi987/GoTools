package data

import (
	"github.com/lightpaw/male7/util/u64"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"strings"
)

func ParseRandAmount(s string) (*RandAmount, error) {

	if len(s) <= 0 {
		return EmptyRandAmount, nil
	}

	array := strings.Split(s, "-")

	switch len(array) {
	case 1:
		min, err := strconv.ParseUint(array[0], 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "RandAmount parse fail, %s", s)
		}

		return &RandAmount{Min: min}, nil
	case 2:
		min, err := strconv.ParseUint(array[0], 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "RandAmount parse fail, %s", s)
		}

		max, err := strconv.ParseUint(array[1], 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "RandAmount parse fail, %s", s)
		}

		if min > max {
			return nil, errors.Wrapf(err, "RandAmount parse fail, min > max, %s", s)
		}

		rand := int64(u64.Sub(max, min))

		return &RandAmount{Min: min, Max: max, rand: rand}, nil
	default:
		return nil, errors.Errorf("RandAmount parse fail, %s", s)
	}

}

var EmptyRandAmount = &RandAmount{}

//gogen:config
type RandAmount struct {
	Min  uint64 `head:"-"` // 最小值
	Max  uint64 `head:"-"` // 最大值
	rand int64  `head:"-"`
}

func (r *RandAmount) IsZero() bool {
	return r.Min == 0 && r.Max == 0
}

func (r *RandAmount) Random() uint64 {
	if r.rand > 1 {
		return r.Min + u64.FromInt64(rand.Int63n(r.rand))
	}

	return r.Min
}
