package collection

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestQueue(t *testing.T) {
	RegisterTestingT(t)

	q := NewQueue()

	for i := 0; i < 10; i++ {
		q.Add(i)
	}

	i := 0
	q.Range(func(v interface{}) (toContinue bool) {
		Ω(v).Should(Equal(i))
		i++
		return true
	})

	q.ReverseRange(func(v interface{}) (toContinue bool) {
		i--
		Ω(v).Should(Equal(i))
		return true
	})

	for i := 0; i < 10; i++ {
		idx := i
		q.RangeWithStartIndex(idx, func(v interface{}) (toContinue bool) {
			Ω(v).Should(Equal(idx))
			idx++
			return true
		})
	}

	for i := 0; i < 10; i++ {
		idx := 10 - i
		q.ReverseRangeWithStartIndex(i, func(v interface{}) (toContinue bool) {
			idx--
			Ω(v).Should(Equal(idx))
			return true
		})
	}

	for i := 0; i < 10; i++ {
		Ω(q.Remove()).Should(Equal(i))
	}

}
