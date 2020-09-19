package u64

import "testing"
import . "github.com/onsi/gomega"

func TestArray(t *testing.T) {
	RegisterTestingT(t)

	array := []uint64{}

	Ω(GetIndex(array, 0)).Should(Equal(-1))
	Ω(GetIndex(array, 1)).Should(Equal(-1))
	Ω(GetIndex(array, 2)).Should(Equal(-1))

	array = AddIfAbsent(array, 0)
	Ω(array).Should(Equal([]uint64{0}))

	array = AddIfAbsent(array, 0)
	Ω(array).Should(Equal([]uint64{0}))

	array = AddIfAbsent(array, 1)
	Ω(array).Should(Equal([]uint64{0, 1}))

	Ω(GetIndex(array, 0)).Should(Equal(0))
	Ω(GetIndex(array, 1)).Should(Equal(1))
	Ω(GetIndex(array, 2)).Should(Equal(-1))

	array = AddIfAbsent(array, 2)
	Ω(array).Should(Equal([]uint64{0, 1, 2}))

	array = RemoveIfPresent(array, 1)
	Ω(array).Should(Equal([]uint64{0, 2}))

	array = RemoveIfPresent(array, 1)
	Ω(array).Should(Equal([]uint64{0, 2}))

	array = RemoveIfPresent(array, 2)
	Ω(array).Should(Equal([]uint64{0}))

	array = RemoveIfPresent(array, 0)
	Ω(array).Should(Equal([]uint64{}))
}

func TestRemoveHead(t *testing.T) {
	RegisterTestingT(t)

	array := []uint64{}
	array = RemoveHead(array)
	Ω(len(array)).Should(Equal(0))

	array = append(array, 0)
	array = append(array, 0)
	array = append(array, 1)
	array = append(array, 2)

	Ω(array).Should(Equal([]uint64{0, 0, 1, 2}))

	array = RemoveHead(array)
	Ω(array).Should(Equal([]uint64{0, 1, 2}))

	array = RemoveHead(array)
	Ω(array).Should(Equal([]uint64{1, 2}))

	array = RemoveHead(array)
	Ω(array).Should(Equal([]uint64{2}))

	array = append(array, 0)
	Ω(array).Should(Equal([]uint64{2, 0}))

	array = RemoveHead(array)
	Ω(array).Should(Equal([]uint64{0}))

	array = RemoveHead(array)
	Ω(array).Should(Equal([]uint64{}))

	array = RemoveHead(array)
	Ω(array).Should(Equal([]uint64{}))
}
