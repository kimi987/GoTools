package u64

func Min(x, y uint64) uint64 {
	if x > y {
		return y
	}

	return x
}
