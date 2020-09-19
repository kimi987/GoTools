package weight

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestWeightRandomer_Index(t *testing.T) {
	RegisterTestingT(t)

	weight := []uint64{}
	r, err := NewWeightRandomer(weight)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{0}
	r, err = NewWeightRandomer(weight)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{0, 1}
	r, err = NewWeightRandomer(weight)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{1, 0, 1}
	r, err = NewWeightRandomer(weight)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{1, 2, 3}
	r, err = NewWeightRandomer(weight)
	Ω(err).Should(Succeed())

	Ω(r.Index(0)).Should(Equal(0))
	Ω(r.Index(1)).Should(Equal(1))
	Ω(r.Index(2)).Should(Equal(1))
	Ω(r.Index(3)).Should(Equal(2))
	Ω(r.Index(4)).Should(Equal(2))
	Ω(r.Index(5)).Should(Equal(2))

	Ω(r.Index(6)).Should(Equal(2))
	Ω(r.Index(7)).Should(Equal(2))
}

func TestU64WeightRandomer_Random(t *testing.T) {
	RegisterTestingT(t)

	weight := []uint64{}
	values := []uint64{}
	r, err := NewU64WeightRandomer(weight, values)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{0}
	values = []uint64{}
	r, err = NewU64WeightRandomer(weight, values)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{}
	values = []uint64{0}
	r, err = NewU64WeightRandomer(weight, values)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{0}
	values = []uint64{0}
	r, err = NewU64WeightRandomer(weight, values)
	Ω(err).Should(HaveOccurred())

	weight = []uint64{1, 2, 3}
	values = []uint64{111, 222, 333}
	r, err = NewU64WeightRandomer(weight, values)
	Ω(err).Should(Succeed())

	r.Random()
}

var r, _ = NewU64WeightRandomer([]uint64{1, 2, 3}, []uint64{111, 222, 333})

func BenchmarkU64WeightRandomer_Random(b *testing.B) {
	var x uint64
	for i := 0; i < b.N; i++ {
		x += r.Random()
	}
}
