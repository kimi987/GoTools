package u64

func Map2Int32Array(m map[uint64]uint64) (keys, values []int32) {

	for k, v := range m {
		keys = append(keys, Int32(k))
		values = append(values, Int32(v))
	}

	return
}

func MapKey2Int32Arrary(m map[uint64]struct{}) (keys []int32) {
	for k := range m {
		keys = append(keys, Int32(k))
	}
	return
}

func MapKey2Uint64Array(m map[uint64]struct{}) (keys []uint64) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func Uint64ArrayToMapKey(keys []uint64) (m map[uint64]struct{}) {
	m = make(map[uint64]struct{})
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return
}

func CopyMap(a map[uint64]uint64) map[uint64]uint64 {
	out := make(map[uint64]uint64, len(a))
	CopyMapTo(out, a)
	return out
}

func CopyMapTo(dest, src map[uint64]uint64) {
	for k, v := range src {
		dest[k] = v
	}
}

func CopyUi64Map(a map[uint64]int64) map[uint64]int64 {
	out := make(map[uint64]int64, len(a))
	CopyUi64MapTo(out, a)
	return out
}

func CopyUi64MapTo(dest, src map[uint64]int64) {
	for k, v := range src {
		dest[k] = v
	}
}
