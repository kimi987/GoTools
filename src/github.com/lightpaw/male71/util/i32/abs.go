package i32

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}

	return x
}
