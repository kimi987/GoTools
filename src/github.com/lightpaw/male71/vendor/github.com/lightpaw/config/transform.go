package config

import (
	"fmt"
	"math"
	"time"
)

func Ia2I32a(a []int) []int32 {
	out := make([]int32, len(a))
	for i, v := range a {
		out[i] = int32(v)
	}
	return out
}

func Ia2I64a(a []int) []int64 {
	out := make([]int64, len(a))
	for i, v := range a {
		out[i] = int64(v)
	}
	return out
}

func U64a2I32a(a []uint64) []int32 {
	out := make([]int32, len(a))
	for i, v := range a {
		out[i] = U64ToI32(v)
	}
	return out
}

func U64ToI32(x uint64) int32 {
	if x > math.MaxInt32 {
		return math.MaxInt32
	}

	return int32(x)
}

func I64a2I32a(a []int64) []int32 {
	out := make([]int32, len(a))
	for i, v := range a {
		out[i] = I64ToI32(v)
	}
	return out
}

func I64ToI32(x int64) int32 {
	if x > math.MaxInt32 {
		return math.MaxInt32
	} else if x < math.MinInt32 {
		return math.MinInt32
	}

	return int32(x)
}

func F64a2F32a(a []float64) []float32 {
	out := make([]float32, len(a))
	for i, v := range a {
		out[i] = F64ToF32(v)
	}
	return out
}

func F64ToF32(x float64) float32 {
	return float32(x)
}

func F64a2I32a(a []float64) []int32 {
	out := make([]int32, len(a))
	for i, v := range a {
		out[i] = F64ToI32(v)
	}
	return out
}

func F64ToI32(x float64) int32 {
	return int32(x)
}

func F32a2I32a(a []float32) []int32 {
	out := make([]int32, len(a))
	for i, v := range a {
		out[i] = F32ToI32(v)
	}
	return out
}

func F32ToI32(x float32) int32 {
	return int32(x)
}

func F64a2I32aX1000(a []float64) []int32 {
	out := make([]int32, len(a))
	for i, v := range a {
		out[i] = F64ToI32X1000(v)
	}
	return out
}

func F64ToI32X1000(x float64) int32 {
	return MultiF64(1000, x)
}

func F32a2I32aX1000(a []float32) []int32 {
	out := make([]int32, len(a))
	for i, v := range a {
		out[i] = F32ToI32X1000(v)
	}
	return out
}

func F32ToI32X1000(x float32) int32 {
	return MultiF64(1000, float64(x))
}

func MultiF64(d uint64, coef float64) int32 {
	if d == 0 || coef <= 0 {
		return 0
	}

	fd := float64(d)
	return int32((coef + (1 / (fd * 10))) * fd)
}

func I32a2Ia(a []int32) []int {
	out := make([]int, len(a))
	for i, v := range a {
		out[i] = int(v)
	}
	return out
}

func I64a2Ia(a []int64) []int {
	out := make([]int, len(a))
	for i, v := range a {
		out[i] = int(v)
	}
	return out
}

func I32m2Im(in map[int32]int32) map[int]int {
	out := make(map[int]int, len(in))

	for k, v := range in {
		out[int(k)] = int(v)
	}

	return out
}

func EnumMapKeys(enumMap map[string]int32, ignores ...int32) []string {
	out := make([]string, 0, len(enumMap))

enum:
	for k, v := range enumMap {
		for _, ign := range ignores {
			if v == ign {
				continue enum
			}
		}

		out = append(out, k)
		out = append(out, fmt.Sprintf("%d", v))
	}

	return out
}

func Duration2I32Seconds(duration time.Duration) int32 {
	return int32(duration / time.Second)
}

func DurationArr2I32Seconds(duration []time.Duration) []int32 {
	seconds := make([]int32, 0, len(duration))
	for _, d := range duration {
		seconds = append(seconds, Duration2I32Seconds(d))
	}
	return seconds
}
