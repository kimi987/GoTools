package i32

func Min(x, y int32) int32 {
	if x > y {
		return y
	}

	return x
}

func Minx(a int32, b ...int32) (min int32) {
	min = a
	for _, x := range b {
		min = Min(min, x)
	}
	return
}
