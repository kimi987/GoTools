package timeutil

import "testing"
import (
	. "github.com/onsi/gomega"
	"time"
)

func TestArray(t *testing.T) {
	RegisterTestingT(t)

	array := []time.Time{Unix64(0), Unix64(1), Unix64(2)}
	array = RemoveByIndex(array, 1)
	Ω(array).Should(Equal([]time.Time{Unix64(0), Unix64(2)}))

	array = RemoveByIndex(array, 1)
	Ω(array).Should(Equal([]time.Time{Unix64(0)}))

	array = RemoveByIndex(array, 1)
	Ω(array).Should(Equal([]time.Time{Unix64(0)}))

	array = RemoveByIndex(array, 0)
	Ω(array).Should(Equal([]time.Time{}))
}

func TestLeftShift2(t *testing.T) {
	RegisterTestingT(t)

	array := []time.Time{Unix64(0), Unix64(1), Unix64(2), Unix64(3),
			     Unix64(4), Unix64(5), Unix64(6), Unix64(7), Unix64(8), Unix64(9)}

	array = LeftShift(array, 0, 1)
	Ω(array).Should(Equal([]time.Time{Unix64(1), Unix64(2), Unix64(3), Unix64(4),
					  Unix64(5), Unix64(6), Unix64(7), Unix64(8), Unix64(9)}))

	array = LeftShift(array, 9, 1)
	Ω(array).Should(Equal([]time.Time{Unix64(1), Unix64(2), Unix64(3), Unix64(4),
					  Unix64(5), Unix64(6), Unix64(7), Unix64(8), Unix64(9)}))

	array = LeftShift(array, 100, 1)
	Ω(array).Should(Equal([]time.Time{Unix64(1), Unix64(2), Unix64(3), Unix64(4),
					  Unix64(5), Unix64(6), Unix64(7), Unix64(8), Unix64(9)}))

	array = LeftShift(array, 1, 3)
	Ω(array).Should(Equal([]time.Time{Unix64(1), Unix64(5), Unix64(6),
					  Unix64(7), Unix64(8), Unix64(9)}))

	array = LeftShift(array, 0, 100)
	Ω(array).Should(Equal([]time.Time{}))

	array = LeftShift(array, 0, 1)
	Ω(array).Should(Equal([]time.Time{}))
}
