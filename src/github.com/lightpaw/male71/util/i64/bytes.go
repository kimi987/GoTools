package i64

import (
	"math/bits"
	"math"
)

func ToBytes(a int64) []byte {
	if a == 0 {
		return nil
	}
	if a == math.MinInt64 {
		return []byte{1}
	}

	var x uint64
	if a < 0 {
		x = (uint64(-a) << 1) | 1
	} else {
		x = uint64(a) << 1
	}

	n := bits.Len64(x)
	bn := (n + 7) / 8
	buf := make([]byte, bn)
	for i := 0; i < bn; i++ {
		buf[i] = uint8(x & 0xff)
		x >>= 8
	}
	return buf
}

func ToBytesU64(a uint64) []byte {
	return ToBytes(int64(a))
}

func BytesArray(as []int64) [][]byte {
	bufs := make([][]byte, len(as))
	for i, a := range as {
		bufs[i] = ToBytes(a)
	}
	return bufs
}

func BytesArrayU64(as []uint64) [][]byte {
	bufs := make([][]byte, len(as))
	for i, a := range as {
		bufs[i] = ToBytes(int64(a))
	}
	return bufs
}

func FromBytes(buf []byte) (int64, bool) {

	switch n := len(buf); n {
	case 0:
		return 0, true
	case 1:
		if buf[0] == 1 {
			return math.MinInt64, true
		}
	default:
		if n > 8 {
			return 0, false
		}
	}

	var a uint64
	for i, b := range buf {
		a |= uint64(b) << uint64(i*8)
	}

	x := int64(a >> 1)
	if a&1 == 1 {
		x = -x
	}

	return x, true
}

func FromBytesU64(buf []byte) (uint64, bool) {
	i64, ok := FromBytes(buf)
	return uint64(i64), ok
}

func FromBytesArray(bufs [][]byte) ([]int64, bool) {
	as := make([]int64, len(bufs))

	var ok bool
	for i, buf := range bufs {
		if as[i], ok = FromBytes(buf); !ok {
			return nil, false
		}
	}

	return as, true
}

func FromBytesArrayU64(bufs [][]byte) ([]uint64, bool) {
	as := make([]uint64, len(bufs))

	var ok bool
	for i, buf := range bufs {
		if as[i], ok = FromBytesU64(buf); !ok {
			return nil, false
		}
	}

	return as, true
}

func DefFromBytes(buf []byte, def int64) int64 {
	if a, ok := FromBytes(buf); ok {
		return a
	}

	return def
}

func IsPositiveVarint(buf []byte) bool {
	if len(buf) <= 0 {
		return true
	}

	return buf[0]&1 == 0
}
