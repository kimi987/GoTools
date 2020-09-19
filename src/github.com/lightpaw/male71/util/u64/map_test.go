package u64

import "testing"
import . "github.com/onsi/gomega"

func TestMap2Int32Array(t *testing.T) {
	RegisterTestingT(t)

	m := map[uint64]uint64{}
	for i := uint64(0); i < 10; i++ {
		m[i] = i * 1000
	}

	keys, values := Map2Int32Array(m)

	Ω(len(keys)).Should(Equal(len(values)))
	Ω(len(keys)).Should(Equal(len(m)))

	for i, k := range keys {
		Ω(uint64(values[i])).Should(Equal(m[uint64(k)]))
	}
}
