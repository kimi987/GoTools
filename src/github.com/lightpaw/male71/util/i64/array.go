package i64

func Contains(array []int64, e int64) bool {
	return GetIndex(array, e) >= 0
}

func GetIndex(array []int64, e int64) int {
	for i, v := range array {
		if v == e {
			return i
		}
	}
	return -1
}

func AddIfAbsent(array []int64, e int64) []int64 {
	idx := GetIndex(array, e)
	if idx < 0 {
		array = append(array, e)
	}

	return array
}

func RemoveIfPresent(array []int64, e int64) []int64 {
	array, _ = RemoveIfPresentReturnIndex(array, e)
	return array
}

func RemoveIfPresentReturnIndex(array []int64, e int64) ([]int64, int) {
	idx := GetIndex(array, e)
	if idx >= 0 {
		return removeIndex(array, idx), idx
	}

	return array, idx
}

func RemoveByIndex(array []int64, idx int) []int64 {
	if idx >= 0 && idx < len(array) {
		return removeIndex(array, idx)
	}

	return array
}

func removeIndex(array []int64, idx int) []int64 {
	lastIndex := len(array) - 1
	if idx != lastIndex {
		// swap
		array[idx], array[lastIndex] = array[lastIndex], array[idx]
	}

	return array[0:lastIndex]
}

func LeftShiftRemoveIfPresent(array []int64, remove int64) []int64 {
	array, _ = LeftShiftRemoveIfPresentReturnIndex(array, remove)
	return array
}

func LeftShiftRemoveIfPresentReturnIndex(array []int64, remove int64) ([]int64, int) {
	if len(array) > 0 {
		for i, v := range array {
			if v == remove {
				return internalLeftShift(array, i, 1), i
			}
		}
	}

	return array, -1
}

func LeftShiftReturnRemovedValues(array []int64, startIndex, count int) ([]int64, []int64) {
	n := len(array)
	if n <= 0 || startIndex < 0 || count <= 0 {
		return array, nil
	}

	if startIndex >= n {
		return array, nil
	} else {
		c := n - startIndex
		if count > c {
			count = c
		}
	}

	removed := make([]int64, count)
	copy(removed, array[startIndex:])

	return internalLeftShift(array, startIndex, count), removed
}

func LeftShift(array []int64, startIndex, count int) []int64 {
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

func internalLeftShift(array []int64, startIndex, count int) []int64 {

	n := len(array)

	copy(array[startIndex:], array[startIndex+count:])
	//array[n-count] = 0
	return array[:n-count]
}

//func internalLeftShift(array []int64, startIndex,count int) []int64 {
//
//	n := len(array)
//	if c := n - startIndex;c < count {
//		count = c
//	}
//
//	if n <= startIndex {
//		return array
//	}
//
//	, count)
//
//	for i := startIndex + 1; i < n; i++ {
//		array[i-1] = array[i]
//	}
//
//	array[n-1] = 0
//	return array[:n-count]
//}

func Copy(a []int64) []int64 {
	out := make([]int64, len(a))
	copy(out, a)
	return out
}
