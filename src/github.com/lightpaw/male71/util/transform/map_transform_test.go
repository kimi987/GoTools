package transform

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestI32m2Im(t *testing.T) {
	RegisterTestingT(t)

	in := map[int32]int32{
		1: 100,
		2: 200,
		3: 300,
	}

	out := I32m2Im(in)
	Ω(out).Should(Equal(map[int]int{
		1: 100,
		2: 200,
		3: 300,
	}))

}

func TestI3264m2I64m(t *testing.T) {
	RegisterTestingT(t)

	in := map[int32]int64{
		1: 100,
		2: 200,
		3: 300,
	}

	out := I3264m2I64m(in)
	Ω(out).Should(Equal(map[int]int64{
		1: 100,
		2: 200,
		3: 300,
	}))

}
