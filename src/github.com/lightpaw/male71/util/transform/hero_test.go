package transform

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestId(t *testing.T) {
	RegisterTestingT(t)

	test(0, 0)
	test(0, 1)
	test(0, -1)
	test(-1, 0)
	test(1, 0)
	test(1, 1)
	test(-1, -1)

	test(10, 10)
	test(10, -10)
	test(-10, 10)
	test(-10, -10)

}

func test(high, low int32) {

	id := Int2Long(high, low)
	x, y := Long2Int(id)
	Ω(x).Should(Equal(high))
	Ω(y).Should(Equal(low))
}
