package i32

import "math"

func Exp(x int32) int32 {
	return int32(math.Exp(float64(x)))
}
