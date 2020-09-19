package util

func StringArrayContains(array []string, s string) bool {

	for _, a := range array {
		if a == s {
			return true
		}
	}

	return false
}

func StringArrayContainsAny(array0, array1 []string) bool {

	if len(array0) > 0 && len(array1) > 0 {
		for _, a := range array0 {
			for _, b := range array1 {
				if a == b {
					return true
				}
			}
		}
	}

	return false
}

func StringArrayRemoveHead(array []string) []string {
	return StringArrayLeftShift(array, 0, 1)
}

func StringArrayLeftShift(array []string, startIndex, count int) []string {
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

	return internalStringArrayLeftShift(array, startIndex, count)
}

func internalStringArrayLeftShift(array []string, startIndex, count int) []string {

	n := len(array)

	copy(array[startIndex:], array[startIndex+count:])
	//array[n-count] = 0
	return array[:n-count]
}
