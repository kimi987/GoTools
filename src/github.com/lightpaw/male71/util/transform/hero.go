package transform

func Int2Long(high, low int32) int64 {
	return int64(uint64(uint32(high))<<32 | uint64(uint32(low)))
}

func Long2Int(x int64) (int32, int32) {
	y := uint64(x)
	return int32(y >> 32), int32(y & 0xffffffff)
}
