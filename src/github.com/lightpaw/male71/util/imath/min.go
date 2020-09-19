package imath

func Min(x, y int) int {
	if x > y {
		return y
	}

	return x
}

func Minx(a int, b ...int) (min int) {
	min = a
	for _, x := range b {
		min = Min(min, x)
	}
	return
}
