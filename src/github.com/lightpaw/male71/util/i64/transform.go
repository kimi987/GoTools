package i64

import (
	"math"
)

func Int32(x int64) int32 {
	if x >= 0 {
		return positiveInt32(x)
	} else {
		return negativeInt32(x)
	}
}

func Int32Array(x []int64) []int32 {
	ia := make([]int32, len(x))
	for i, v := range x {
		ia[i] = Int32(v)
	}

	return ia
}

func positiveInt32(x int64) int32 {
	return int32(Min(Max(x, 0), math.MaxInt32))
}

func negativeInt32(x int64) int32 {
	return int32(Max(Min(x, 0), math.MinInt32))
}

type GetU64 func(k int64) uint64

func NewGetU64(m map[int64]uint64) GetU64 {
	return func(k int64) uint64 {
		return m[k]
	}
}

func EmptyGetU64() GetU64 {
	return func(k int64) uint64 {
		return 0
	}
}

func Rate(start, end, cur int64) float64 {

	// 特殊情况处理
	if end <= cur {
		return 1
	}

	if cur <= start {
		return 0
	}

	diff := cur - start
	total := end - start

	return float64(diff) / float64(total)
}
