package idbytes

import (
	. "github.com/onsi/gomega"
	"math"
	"testing"
)

func TestIdBytes(t *testing.T) {
	RegisterTestingT(t)

	for i := 0; i > 0 && i < math.MaxInt64; {
		v := int64(i)
		bytes := ToBytes(v)
		toV, ok := ToId(bytes)
		Ω(ok).Should(BeTrue())
		Ω(v).Should(BeEquivalentTo(toV))
		Ω(IsPositive(bytes)).Should(BeTrue())
		//Ω(Len(v)).Should(BeEquivalentTo(len(bytes)))
		//Ω(Len(v)).Should(BeEquivalentTo(cap(bytes)))

		if i <= 65536 {
			i++
		} else {
			i *= 2
		}
	}

	for i := math.MinInt64; i < 0; {
		v := int64(i)
		bytes := ToBytes(v)
		toV, ok := ToId(bytes)
		Ω(ok).Should(BeTrue())
		Ω(v).Should(BeEquivalentTo(toV))
		Ω(IsPositive(bytes)).Should(BeFalse())
		//Ω(Len(v)).Should(BeEquivalentTo(len(bytes)))
		//Ω(Len(v)).Should(BeEquivalentTo(cap(bytes)))

		if i >= -65536 {
			i++
		} else {
			i /= 2
		}
	}
}
