package u64

import (
	. "github.com/onsi/gomega"
	"math"
	"testing"
)

func TestFromInt32(t *testing.T) {
	RegisterTestingT(t)
	Ω(FromInt32(-1)).Should(Equal(uint64(0)))
	Ω(FromInt32(0)).Should(Equal(uint64(0)))
	Ω(FromInt32(1)).Should(Equal(uint64(1)))
	Ω(FromInt32(100)).Should(Equal(uint64(100)))
}

func TestInt32(t *testing.T) {
	RegisterTestingT(t)
	Ω(Int32(math.MaxInt64)).Should(Equal(int32(math.MaxInt32)))
	Ω(Int32(0)).Should(Equal(int32(0)))
	Ω(Int32(1)).Should(Equal(int32(1)))
	Ω(Int32(100)).Should(Equal(int32(100)))
}
