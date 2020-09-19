package transform

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestIa2I32a(t *testing.T) {
	RegisterTestingT(t)

	var in = []int{1, 2, 3}
	out := Ia2I32a(in)
	Ω(out).Should(Equal([]int32{1, 2, 3}))
}

func TestIa2I64a(t *testing.T) {
	RegisterTestingT(t)

	var in = []int{1, 2, 3}
	out := Ia2I64a(in)
	Ω(out).Should(Equal([]int64{1, 2, 3}))
}

func TestI32a2Ia(t *testing.T) {
	RegisterTestingT(t)

	var in = []int32{1, 2, 3}
	out := I32a2Ia(in)
	Ω(out).Should(Equal([]int{1, 2, 3}))
}

func TestI64a2Ia(t *testing.T) {
	RegisterTestingT(t)

	var in = []int64{1, 2, 3}
	out := I64a2Ia(in)
	Ω(out).Should(Equal([]int{1, 2, 3}))
}
