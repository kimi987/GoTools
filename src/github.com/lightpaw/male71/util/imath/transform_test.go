package imath

import (
	. "github.com/onsi/gomega"
	"math"
	"testing"
)

func TestInt32(t *testing.T) {
	RegisterTestingT(t)

	Ω(Int32(0)).Should(Equal(int32(0)))
	Ω(Int32(math.MaxInt32)).Should(Equal(int32(math.MaxInt32)))
	Ω(Int32(math.MaxInt64)).Should(Equal(int32(math.MaxInt32)))
	Ω(Int32(math.MinInt32)).Should(Equal(int32(math.MinInt32)))
	Ω(Int32(math.MinInt64)).Should(Equal(int32(math.MinInt32)))
}
