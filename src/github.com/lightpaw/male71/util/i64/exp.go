package i64

import "math"

func Exp(x int64) int64 {
	return int64(math.Exp(float64(x)))
}
