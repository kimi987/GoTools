package u64

func Equals(a, b []uint64) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func Contains(array []uint64, e uint64) bool {
	return GetIndex(array, e) >= 0
}

func GetIndex(array []uint64, e uint64) int {
	for i, v := range array {
		if v == e {
			return i
		}
	}
	return -1
}

func AddIfAbsent(array []uint64, e uint64) []uint64 {
	idx := GetIndex(array, e)
	if idx < 0 {
		array = append(array, e)
	}

	return array
}

func TryAdd(array []uint64, e uint64) (bool, []uint64) {
	idx := GetIndex(array, e)
	if idx < 0 {
		return true, append(array, e)
	}

	return false, array
}

func RemoveIfPresent(array []uint64, e uint64) []uint64 {
	idx := GetIndex(array, e)
	if idx >= 0 {
		lastIndex := len(array) - 1
		if idx != lastIndex {
			// swap
			array[idx], array[lastIndex] = array[lastIndex], array[idx]
		}

		return array[0:lastIndex]
	}

	return array
}

func RemoveHead(array []uint64) []uint64 {
	return LeftShift(array, 0, 1)
}

func LeftShift(array []uint64, startIndex, count int) []uint64 {
	n := len(array)
	if n <= 0 || startIndex < 0 || count <= 0 {
		return array
	}

	if startIndex >= n {
		return array
	} else {
		c := n - startIndex
		if count > c {
			count = c
		}
	}

	return internalLeftShift(array, startIndex, count)
}

func internalLeftShift(array []uint64, startIndex, count int) []uint64 {

	n := len(array)

	copy(array[startIndex:], array[startIndex+count:])
	//array[n-count] = 0
	return array[:n-count]
}

func Copy(a []uint64) []uint64 {
	out := make([]uint64, len(a))
	copy(out, a)
	return out
}
