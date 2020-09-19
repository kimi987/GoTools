package i32

func Uint64(x int32) uint64 {
	if x < 0 {
		return 0
	}

	return uint64(x)
}

func MultiF64(d uint64, coef float64) int32 {
	if d == 0 || coef <= 0 {
		return 0
	}

	fd := float64(d)
	return int32((coef + (1 / (fd * 10))) * fd)
}

func MultiCoef(d int32, coef float64) int32 {
	if d == 0 || coef <= 0 {
		return 0
	}

	fd := float64(d)
	return int32((coef + (1 / (fd * 10))) * fd)
}
