package transform

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestIntArray(t *testing.T) {
	RegisterTestingT(t)

	var in = []int{1, 2, 3}
	out := IntArray(in)
	Ω(out).Should(Equal(in))
	Ω(&out == &in).Should(BeFalse())
}

func TestInt32Array(t *testing.T) {
	RegisterTestingT(t)

	var in = []int32{1, 2, 3}
	out := Int32Array(in)
	Ω(out).Should(Equal(in))
	Ω(&out == &in).Should(BeFalse())
}

func TestInt64Array(t *testing.T) {
	RegisterTestingT(t)

	var in = []int64{1, 2, 3}
	out := Int64Array(in)
	Ω(out).Should(Equal(in))
	Ω(&out == &in).Should(BeFalse())
}

func TestStringArray(t *testing.T) {
	RegisterTestingT(t)

	var in = []string{"1", "b", "3"}
	out := StringArray(in)
	Ω(out).Should(Equal(in))
	Ω(&out == &in).Should(BeFalse())
}
