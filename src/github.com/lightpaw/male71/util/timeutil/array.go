package timeutil

import "time"

func RemoveByIndex(array []time.Time, idx int) []time.Time {
	if idx >= 0 && idx < len(array) {
		return removeIndex(array, idx)
	}

	return array
}

func removeIndex(array []time.Time, idx int) []time.Time {
	lastIndex := len(array) - 1
	if idx != lastIndex {
		// swap
		array[idx], array[lastIndex] = array[lastIndex], array[idx]
	}

	return array[0:lastIndex]
}

func LeftShift(array []time.Time, startIndex, count int) []time.Time {
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

func internalLeftShift(array []time.Time, startIndex, count int) []time.Time {

	n := len(array)

	copy(array[startIndex:], array[startIndex+count:])
	//array[n-count] = 0
	return array[:n-count]
}
