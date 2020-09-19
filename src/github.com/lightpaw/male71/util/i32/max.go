package i32

func Max(x, y int32) int32 {
	if x < y {
		return y
	}

	return x
}

func Maxx(a int32, b ...int32) (max int32) {
	max = a
	for _, x := range b {
		max = Max(max, x)
	}
	return
}

func Maxa(a []int32) (max int32) {
	if n := len(a); n > 0 {
		max = a[0]
		for i := 1; i < n; i++ {
			max = Max(max, a[i])
		}
	}
	return
}
