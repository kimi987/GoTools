package u64

func DivideTimes(x, y uint64) uint64 {
	if x <= 0 || y <= 0 {
		return 0
	}

	return (x + y - 1) / y
}
