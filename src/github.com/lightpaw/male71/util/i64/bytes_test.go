package i64

import (
	. "github.com/onsi/gomega"
	"math"
	"testing"
)

func i64ToBytes(a int64) []byte {

	if a == 0 {
		return nil
	}

	if a == math.MinInt64 {
		return []byte{1}
	}

	var x uint64
	if a < 0 {
		x = (uint64(-a) * 2) + 1
	} else {
		x = uint64(a) * 2
	}

	bn := i64Len(x)
	buf := make([]byte, bn)
	for i := 0; i < bn; i++ {
		buf[i] = uint8(x % 256)
		x = x / 256
	}

	return buf
}

func bytesToi64(buf []byte) int64 {

	switch n := len(buf); n {
	case 0:
		return 0
	case 1:
		if buf[0] == 1 {
			return math.MinInt64
		}
	default:
		if n > 8 {
			return 0
		}
	}

	var a uint64
	p := uint64(1)
	for _, b := range buf {
		a += uint64(b) * p
		p = p * 256
	}

	x := int64(a / 2)
	if a%2 == 1 {
		x = -x
	}

	return x
}

func i64Len(a uint64) int {
	x := uint64(256)
	for i := 1; i <= 8; i++ {
		if a < x {
			return i
		}
		x = x * 256
	}
	return 8
}

func TestVarintBytes(t *testing.T) {
	RegisterTestingT(t)

	for i := 0; i > 0 && i < math.MaxInt64; {
		v := int64(i)
		bytes := ToBytes(v)
		toV, ok := FromBytes(bytes)
		Ω(ok).Should(BeTrue())
		Ω(v).Should(BeEquivalentTo(toV))
		Ω(IsPositiveVarint(bytes)).Should(BeTrue())

		cbytes := i64ToBytes(v)
		Ω(bytes).Should(BeEquivalentTo(cbytes))
		toV = bytesToi64(cbytes)
		Ω(v).Should(BeEquivalentTo(toV))

		if i <= 65536 {
			i++
		} else {
			i *= 2
		}
	}

	for i := math.MinInt64; i < 0; {
		v := int64(i)
		bytes := ToBytes(v)
		V, ok := FromBytes(bytes)
		Ω(ok).Should(BeTrue())
		Ω(v).Should(BeEquivalentTo(V))
		Ω(IsPositiveVarint(bytes)).Should(BeFalse())

		cbytes := i64ToBytes(v)
		Ω(bytes).Should(BeEquivalentTo(cbytes))
		V = bytesToi64(cbytes)
		Ω(v).Should(BeEquivalentTo(V))

		if i >= -65536 {
			i++
		} else {
			i /= 2
		}
	}
}

func TestVaruintBytes(t *testing.T) {
	RegisterTestingT(t)

	for i := 0; i > 0 && i < math.MaxInt64; {
		v := uint64(i)
		bytes := ToBytes(int64(v))
		toV, ok := FromBytes(bytes)
		Ω(ok).Should(BeTrue())
		Ω(v).Should(BeEquivalentTo(uint64(toV)))
		Ω(IsPositiveVarint(bytes)).Should(BeTrue())

		if i <= 65536 {
			i++
		} else {
			i *= 2
		}
	}

	for i := math.MinInt64; i < 0; {
		v := uint64(i)
		bytes := ToBytes(int64(v))
		V, ok := FromBytes(bytes)
		Ω(ok).Should(BeTrue())
		Ω(v).Should(BeEquivalentTo(uint64(V)))
		Ω(IsPositiveVarint(bytes)).Should(BeFalse())

		if i >= -65536 {
			i++
		} else {
			i /= 2
		}
	}

	for i := 0; i < 65536; i++ {
		v := uint64(math.MaxUint64) - uint64(i)
		bytes := ToBytes(int64(v))
		toV, ok := FromBytes(bytes)
		Ω(ok).Should(BeTrue())
		Ω(v).Should(BeEquivalentTo(uint64(toV)))
	}

	for i := 0; i < 65536; i++ {
		v := uint64(math.MaxInt64) - uint64(i)
		bytes := ToBytes(int64(v))
		toV, ok := FromBytes(bytes)
		Ω(ok).Should(BeTrue())
		Ω(v).Should(BeEquivalentTo(uint64(toV)))
	}

	for i := 0; i < 65536; i++ {
		v := uint64(math.MaxInt64) + uint64(i)
		bytes := ToBytes(int64(v))
		toV, ok := FromBytes(bytes)
		Ω(ok).Should(BeTrue())
		Ω(v).Should(BeEquivalentTo(uint64(toV)))
	}
}
