package imath

import (
	"math"
)

func Int32(x int) int32 {
	if x >= 0 {
		return positiveInt32(x)
	} else {
		return negativeInt32(x)
	}
}

func positiveInt32(x int) int32 {
	return int32(Min(Max(x, 0), math.MaxInt32))
}

func negativeInt32(x int) int32 {
	return int32(Max(Min(x, 0), math.MinInt32))
}
