package i64

func Pow(a, n int64) int64 {
	if n < 0 {
		return 0
	}
	var x int64 = 1
	for n != 0 {
		if (n & 1) == 1 {
			x *= a
		}
		n >>= 1
		a *= a
	}

	return x
}
