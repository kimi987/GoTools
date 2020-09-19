package i64

func DivideTimes(x, y int64) int64 {
	if x <= 0 || y <= 0 {
		return 0
	}

	return (x + y - 1) / y
}

func Division2Float64(x, y int64) float64 {
	if y == 0 {
		return 0
	}

	return float64(x) / float64(y)
}
