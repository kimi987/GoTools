package check

func Int32Duplicate(array []int32) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func Int32DuplicateIgnoreZero(array []int32) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		if x == 0 {
			continue
		}

		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func Int32CountIgnoreZero(array []int32) int {
	c := 0
	for _, v := range array {
		if v != 0 {
			c++
		}
	}

	return c
}

func Int32AnyZero(array []int32) bool {
	for _, v := range array {
		if v == 0 {
			return true
		}
	}

	return false
}

func Int32AnyLe0(array []int32) bool {
	for _, v := range array {
		if v <= 0 {
			return true
		}
	}

	return false
}

func Int32AnyLt0(array []int32) bool {
	for _, v := range array {
		if v < 0 {
			return true
		}
	}

	return false
}

func Int32AnyGt0(array []int32) bool {
	for _, v := range array {
		if v > 0 {
			return true
		}
	}

	return false
}

func IntDuplicate(array []int) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func Int64Duplicate(array []int64) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func Uint64Duplicate(array []uint64) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func StringDuplicate(array []string) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}

func StringNilOrDuplicate(array []string) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		if len(x) <= 0 {
			return true
		}

		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}
