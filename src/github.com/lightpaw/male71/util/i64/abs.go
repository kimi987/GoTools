package i64

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}

	return x
}
