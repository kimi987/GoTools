package i64

func Max(x, y int64) int64 {
	if x < y {
		return y
	}

	return x
}
