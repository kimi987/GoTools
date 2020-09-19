package weight

import (
	"testing"
	. "github.com/onsi/gomega"
	"math/rand"
	"fmt"
)

func TestRandomN(t *testing.T) {
	RegisterTestingT(t)

	for _, weight := range [][]uint64{
		{0, 0}, {1, 0}, {0, 1}, {1, 0, 1},
	} {
		index, err := RandomN(weight, 1)
		Ω(err).Should(Equal(ErrZeroWeight))
		Ω(index).Should(BeNil())
	}

	for _, weight := range [][]uint64{
		{0}, {1, 1}, {2, 2, 2}, {3, 0, 3, 4},
	} {
		count := len(weight)

		index, err := RandomN(weight, count)
		Ω(err).Should(Succeed())

		a := make([]int, count)
		for i := 0; i < count; i++ {
			a[i] = i
		}
		Ω(index).Should(BeEquivalentTo(a))
	}

	for _, weight := range [][]uint64{
		{}, {0}, {1, 1}, {2, 2, 2}, {3, 0, 3, 4},
	} {
		count := len(weight)

		index, err := RandomN(weight, count+1)
		Ω(err).Should(Equal(ErrChoiceNotEnough))
		Ω(index).Should(BeNil())
	}

	weight := []uint64{0}
	index, err := RandomN(weight, 1)
	Ω(err).Should(Succeed())
	Ω(index).Should(BeEquivalentTo([]int{0}))

	weight = []uint64{}
	index, err = RandomN(weight, 1)
	Ω(err).Should(Equal(ErrChoiceNotEnough))
	Ω(index).Should(BeNil())

	for i := 0; i < 1000; i++ {
		weight := []uint64{10, 4, 1, 5, 6}
		index, err := RandomN(weight, rand.Intn(len(weight)-1)+1)
		Ω(err).Should(Succeed())

		// 没有重复元素
		for i, v0 := range index {
			for j, v1 := range index {
				if i != j {
					Ω(v0).ShouldNot(BeEquivalentTo(v1))
				}
			}

			Ω(v0 >= 0 && v0 < len(weight)).Should(BeTrue())
		}
	}
}

func TestRandomN2(t *testing.T) {
	RegisterTestingT(t)

	weight := []uint64{1, 2, 3, 4, 5}

	for n := 0; n < 5; n++ {
		timesMap := make(map[int]int)
		for i := 0; i < 10000; i++ {
			index, err := RandomN(weight, n+1)
			Ω(err).Should(Succeed())

			for _, v := range index {
				timesMap[v]++
			}
		}
		fmt.Println(timesMap)
	}

}
