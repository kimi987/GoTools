package atomic

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestU64AddLimit(t *testing.T) {
	RegisterTestingT(t)

	amt := NewUint64(0)

	realAdd, total := U46AddLimit(amt, 50, 100)
	Ω(realAdd).Should(Equal(uint64(50)))
	Ω(total).Should(Equal(uint64(50)))

	realAdd, total = U46AddLimit(amt, 51, 100)
	Ω(realAdd).Should(Equal(uint64(50)))
	Ω(total).Should(Equal(uint64(100)))

	realAdd, total = U46AddLimit(amt, 11, 100)
	Ω(realAdd).Should(Equal(uint64(0)))
	Ω(total).Should(Equal(uint64(100)))
}
