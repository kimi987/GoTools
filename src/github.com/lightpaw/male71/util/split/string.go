package split

import (
	"github.com/pkg/errors"
	"strconv"
)

func Sa2Ia(in []string) ([]int, error) {

	out := make([]int, len(in))
	for i, s := range in {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.Wrapf(err, "split.Sa2Ia atoi error, %s", in)
		}

		out[i] = v
	}

	return out, nil
}

func Sa2I64a(in []string) ([]int64, error) {

	out := make([]int64, len(in))
	for i, s := range in {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "split.Sa2I64a ParseInt error, %s", in)
		}

		out[i] = v
	}

	return out, nil
}

func Sa2F32a(in []string) ([]float32, error) {

	out := make([]float32, len(in))
	for i, s := range in {
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, errors.Wrap(err, "split.Sa2F32a ParseFloat error")
		}

		out[i] = float32(v)
	}

	return out, nil
}

func Sa2F64a(in []string) ([]float64, error) {

	if in == nil {
		return nil, nil
	}

	out := make([]float64, len(in))
	for i, s := range in {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, errors.Wrap(err, "split.Sa2F64a ParseFloat error")
		}

		out[i] = v
	}

	return out, nil
}
