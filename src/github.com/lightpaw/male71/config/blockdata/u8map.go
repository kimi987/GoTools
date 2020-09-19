package blockdata

func newUint8Map(xLen, yLen uint64, overflowReturn uint8) *uint8_map {
	b := &uint8_map{}

	b.xLen = xLen
	b.yLen = yLen
	b.array = make([]uint8, xLen*yLen)
	b.overflowReturn = overflowReturn

	return b
}

type uint8_map struct {
	xLen, yLen uint64

	array []uint8

	overflowReturn uint8
}

func (b *uint8_map) SetInt(x, y int, toSet uint8) bool {
	if x < 0 || y < 0 {
		return false
	}

	return b.Set(uint64(x), uint64(y), toSet)
}

func (b *uint8_map) Set(x, y uint64, toSet uint8) bool {
	if x >= b.xLen || y >= b.yLen {
		return false
	}

	b.array[b.index(x, y)] = toSet
	return true
}

func (b *uint8_map) GetInt(x, y int) uint8 {
	if x < 0 || y < 0 {
		return b.overflowReturn
	}

	return b.Get(uint64(x), uint64(y))
}

func (b *uint8_map) Get(x, y uint64) uint8 {
	if x >= b.xLen || y >= b.yLen {
		return b.overflowReturn
	}

	return b.array[b.index(x, y)]
}

func (b *uint8_map) index(x, y uint64) uint64 {
	return x*b.yLen + y
}
