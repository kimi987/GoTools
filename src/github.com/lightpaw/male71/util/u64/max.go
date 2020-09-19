package u64

func Max(x, y uint64) uint64 {
	if x < y {
		return y
	}

	return x
}

func Maxa(a []uint64) (max uint64) {
	if n := len(a); n > 0 {
		max = a[0]
		for i := 1; i < n; i++ {
			max = Max(max, a[i])
		}
	}
	return
}
