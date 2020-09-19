package i32

func Pow(a, n int32) int32 {
	if n < 0 {
		return 0
	}
	var x int32 = 1
	for n != 0 {
		if (n & 1) == 1 {
			x *= a
		}
		n >>= 1
		a *= a
	}

	return x
}
