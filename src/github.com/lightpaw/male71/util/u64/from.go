package u64

import "math"

func Int32(x uint64) int32 {
	if x > math.MaxInt32 {
		return math.MaxInt32
	}

	return int32(x)
}

func Int32Array(x []uint64) []int32 {
	ia := make([]int32, len(x))
	for i, v := range x {
		ia[i] = Int32(v)
	}

	return ia
}

func FromInt32(x int32) uint64 {
	if x < 0 {
		return 0
	}

	return uint64(x)
}

func FromInt32Array(x []int32) []uint64 {
	ia := make([]uint64, len(x))
	for i, v := range x {
		ia[i] = FromInt32(v)
	}

	return ia
}

func Int(x uint64) int {
	if x > math.MaxInt32 {
		return math.MaxInt32
	}

	return int(x)
}

func FromInt(x int) uint64 {
	if x < 0 {
		return 0
	}

	return uint64(x)
}

func FromInt64(x int64) uint64 {
	if x < 0 {
		return 0
	}

	return uint64(x)
}

//func From2I32(a, b int32) uint64 {
//	ua := FromInt32(a)
//	ub := FromInt32(b)
//	return ua<<32 | ub
//}
//
//func From2U32(a, b uint32) uint64 {
//	ua := uint64(a)
//	ub := uint64(b)
//	return ua<<32 | ub
//}