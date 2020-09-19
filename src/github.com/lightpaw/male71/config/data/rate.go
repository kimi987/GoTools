package data

import (
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"strings"
)

func ParseRate(s string) (*Rate, error) {

	if len(s) <= 0 {
		return EmptyRate, nil
	}

	array := strings.Split(s, "/")

	switch len(array) {
	case 1:
		denominator := int64(10000)

		min, err := strconv.ParseInt(array[0], 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "Rate parse fail, %s", s)
		}

		if min < 0 {
			return nil, errors.Errorf("Rate parse fail, 分子 < 0, %s", s)
		}

		if min > denominator {
			return nil, errors.Errorf("Rate parse fail, 分子 > 分母, %s", denominator, s)
		}

		return &Rate{Numerator: min, Denominator: denominator}, nil
	case 2:
		min, err := strconv.ParseInt(array[0], 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "Rate parse fail, %s", s)
		}

		max, err := strconv.ParseInt(array[1], 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "Rate parse fail, %s", s)
		}

		if min < 0 {
			return nil, errors.Errorf("Rate parse fail, 分子 < 0, %s", s)
		}

		if max <= 0 {
			return nil, errors.Errorf("Rate parse fail, 分母 <= 0, %s", s)
		}

		if min > max {
			return nil, errors.Errorf("Rate parse fail, 分子 > 分母, %s", s)
		}

		return &Rate{Numerator: min, Denominator: max}, nil
	default:
		return nil, errors.Errorf("RandAmount parse fail, %s", s)
	}

}

var EmptyRate = &Rate{}

//gogen:config
type Rate struct {
	Numerator   int64 `head:"-"`
	Denominator int64 `head:"-"`
}

func (r *Rate) Try() bool {
	if r.Numerator <= 0 || r.Denominator <= 0 {
		return false
	}

	if r.Denominator <= r.Numerator {
		return true
	}

	return rand.Int63n(r.Denominator) < r.Numerator
}
