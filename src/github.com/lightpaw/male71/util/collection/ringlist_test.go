package collection

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestRingList(t *testing.T) {
	RegisterTestingT(t)

	q := NewRingList(10)

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

	for i := 0; i < 10; i++ {
		q.Add(i)
	}
	// 超出capcity，自动移除最小的值
	for i := 0; i < 10; i++ {
		q.Add(10 - i)
	}

	i = 10
	q.Range(func(v interface{}) (toContinue bool) {
		Ω(v).Should(Equal(i))
		i--
		return true
	})

	for i := 0; i < 10; i++ {
		Ω(q.Remove()).Should(Equal(10 - i))
	}
}
