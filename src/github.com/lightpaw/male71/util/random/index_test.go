package random

import (
	. "github.com/onsi/gomega"
	"math/rand"
	"testing"
	"time"
)

func TestMixIntArray(t *testing.T) {
	RegisterTestingT(t)
	rand.Seed(time.Now().UnixNano())

	list := []int{0, 1, 2, 3, 4, 5}
	mix := make([]int, len(list))
	copy(mix, list)

	Ω(mix).Should(Equal(list))

	MixIntArray(mix)
	Ω(mix).Should(ConsistOf(list[0], list[1], list[2], list[3], list[4], list[5]))
}

func TestMixU64Array(t *testing.T) {
	RegisterTestingT(t)
	rand.Seed(time.Now().UnixNano())

	list := []uint64{0, 1, 2, 3, 4, 5}
	mix := make([]uint64, len(list))
	copy(mix, list)

	Ω(mix).Should(Equal(list))

	MixU64Array(mix)
	Ω(mix).Should(ConsistOf(list[0], list[1], list[2], list[3], list[4], list[5]))
}

func TestNewIntIndexArray(t *testing.T) {
	RegisterTestingT(t)
	rand.Seed(time.Now().UnixNano())

	a := NewIntIndexArray(0)
	Ω(a).Should(BeEmpty())

	for i := 1; i < 10; i++ {
		a = NewIntIndexArray(i)
		Ω(a).Should(HaveLen(i))

		for n := 0; n < i; n++ {
			Ω(a).Should(ContainElement(n))
		}
	}
}

func TestNewU64IndexArray(t *testing.T) {
	RegisterTestingT(t)
	rand.Seed(time.Now().UnixNano())

	a := NewU64IndexArray(0)
	Ω(a).Should(BeEmpty())

	for i := uint64(1); i < 10; i++ {
		a = NewU64IndexArray(i)
		Ω(a).Should(HaveLen(int(i)))

		for n := uint64(0); n < i; n++ {
			Ω(a).Should(ContainElement(n))
		}
	}
}

func TestNewMNIntArray(t *testing.T) {
	RegisterTestingT(t)
	rand.Seed(time.Now().UnixNano())

	a := NewMNIntIndexArray(0, 0)
	Ω(a).Should(BeEmpty())

	a = NewMNIntIndexArray(2, 2)
	Ω(a).Should(HaveLen(2))
	Ω(a).Should(ContainElement(0))
	Ω(a).Should(ContainElement(1))

	for m := 1; m < 100; m++ {
		a = NewMNIntIndexArray(m, m)
		Ω(a).Should(HaveLen(m))
		for n := 0; n <= m; n++ {
			if n < m {
				Ω(a).Should(ContainElement(n))
			}

			mn := NewMNIntIndexArray(m, n)
			Ω(mn).Should(HaveLen(n))

			Ω(intDuplicate(a)).Should(BeFalse())
			for _, v := range mn {
				Ω(v < m).Should(BeTrue())
			}
		}
	}
}

func intDuplicate(array []int) bool {
	n := len(array)
	for i := 0; i < n; i++ {
		x := array[i]
		for j := i + 1; j < n; j++ {
			if x == array[j] {
				return true
			}
		}
	}

	return false
}
