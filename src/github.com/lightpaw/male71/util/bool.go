package util

func TrueCount(arr []bool) uint64 {
	if len(arr) <= 0 {
		return 0
	}
	count := uint64(0)
	for _, b := range arr {
		if b {
			count++
		}
	}
	return count
}

func FalseCount(arr []bool) uint64 {
	if len(arr) <= 0 {
		return 0
	}
	count := uint64(0)
	for _, b := range arr {
		if !b {
			count++
		}
	}
	return count
}
