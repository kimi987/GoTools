package i64

func Map2Int32Array(m map[int64]int64) (keys, values []int32) {

	for k, v := range m {
		keys = append(keys, Int32(k))
		values = append(values, Int32(v))
	}

	return
}

func CopyMap(a map[int64]int64) map[int64]int64 {
	out := make(map[int64]int64, len(a))
	CopyMapTo(out, a)
	return out
}

func CopyMapTo(dest, src map[int64]int64) {
	for k, v := range src {
		dest[k] = v
	}
}
