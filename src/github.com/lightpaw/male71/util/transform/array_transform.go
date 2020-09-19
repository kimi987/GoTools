package transform

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
