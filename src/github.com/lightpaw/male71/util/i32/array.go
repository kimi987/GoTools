package i32

func Equals(a, b []int32) bool {
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

func Contains(array []int32, e int32) bool {
	return GetIndex(array, e) >= 0
}

func GetIndex(array []int32, e int32) int {
	for i, v := range array {
		if v == e {
			return i
		}
	}
	return -1
}

func AddIfAbsent(array []int32, e int32) []int32 {
	idx := GetIndex(array, e)
	if idx < 0 {
		array = append(array, e)
	}

	return array
}

func TryAdd(array []int32, e int32) (bool, []int32) {
	idx := GetIndex(array, e)
	if idx < 0 {
		return true, append(array, e)
	}

	return false, array
}

func RemoveIfPresent(array []int32, e int32) []int32 {
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

func RemoveHead(array []int32) []int32 {
	return LeftShift(array, 0, 1)
}

func LeftShift(array []int32, startIndex, count int) []int32 {
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

func internalLeftShift(array []int32, startIndex, count int) []int32 {

	n := len(array)

	copy(array[startIndex:], array[startIndex+count:])
	//array[n-count] = 0
	return array[:n-count]
}

func Copy(a []int32) []int32 {
	out := make([]int32, len(a))
	copy(out, a)
	return out
}
