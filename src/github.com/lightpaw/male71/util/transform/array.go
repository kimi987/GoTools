package transform

func Int64Array(a []int64) []int64 {
	out := make([]int64, len(a))
	copy(out, a)
	return out
}

func Int32Array(a []int32) []int32 {
	out := make([]int32, len(a))
	copy(out, a)
	return out
}

func IntArray(a []int) []int {
	out := make([]int, len(a))
	copy(out, a)
	return out
}

func StringArray(a []string) []string {
	out := make([]string, len(a))
	copy(out, a)
	return out
}
