package config

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

func ParseIntArray(in, sep string) ([]int, error) {
	return Sa2Ia(strings.Split(in, sep))
}

func ParseInt64Array(in, sep string) ([]int64, error) {
	return Sa2I64a(strings.Split(in, sep))
}

func ParseFloat32Array(in, sep string) ([]float32, error) {
	return Sa2F32a(strings.Split(in, sep))
}

func ParseFloat64Array(in, sep string) ([]float64, error) {
	return Sa2F64a(strings.Split(in, sep))
}

func Sa2Ia(in []string) ([]int, error) {

	out := make([]int, len(in))
	for i, s := range in {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.Wrapf(err, "string array convert to int array fail, %s", in)
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
			return nil, errors.Wrapf(err, "string array convert to int64 array fail, %s", in)
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
			return nil, errors.Wrapf(err, "string array convert to float32 array fail, %s", in)
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
			return nil, errors.Wrapf(err, "string array convert to float64 array fail, %s", in)
		}

		out[i] = v
	}

	return out, nil
}

func ParseDuration(s string) (time.Duration, error) {
	if len(s) > 0 {
		return time.ParseDuration(s)
	}

	return time.Duration(0), nil
}
