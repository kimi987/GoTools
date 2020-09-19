package i64

import "testing"
import (
	. "github.com/onsi/gomega"
)

func TestArray(t *testing.T) {
	RegisterTestingT(t)

	array := []int64{}

	Ω(GetIndex(array, 0)).Should(Equal(-1))
	Ω(GetIndex(array, 1)).Should(Equal(-1))
	Ω(GetIndex(array, 2)).Should(Equal(-1))

	array = AddIfAbsent(array, 0)
	Ω(array).Should(Equal([]int64{0}))

	array = AddIfAbsent(array, 0)
	Ω(array).Should(Equal([]int64{0}))

	array = AddIfAbsent(array, 1)
	Ω(array).Should(Equal([]int64{0, 1}))

	Ω(GetIndex(array, 0)).Should(Equal(0))
	Ω(GetIndex(array, 1)).Should(Equal(1))
	Ω(GetIndex(array, 2)).Should(Equal(-1))

	array = AddIfAbsent(array, 2)
	Ω(array).Should(Equal([]int64{0, 1, 2}))

	array = RemoveIfPresent(array, 1)
	Ω(array).Should(Equal([]int64{0, 2}))

	array = RemoveIfPresent(array, 1)
	Ω(array).Should(Equal([]int64{0, 2}))

	array = RemoveIfPresent(array, 2)
	Ω(array).Should(Equal([]int64{0}))

	array = RemoveIfPresent(array, 0)
	Ω(array).Should(Equal([]int64{}))

	array = []int64{0, 1, 2}
	idx := -1
	array, idx = RemoveIfPresentReturnIndex(array, 1)
	Ω(array).Should(Equal([]int64{0, 2}))
	Ω(idx).Should(Equal(1))

	array, idx = RemoveIfPresentReturnIndex(array, 1)
	Ω(array).Should(Equal([]int64{0, 2}))
	Ω(idx).Should(Equal(-1))

	array, idx = RemoveIfPresentReturnIndex(array, 2)
	Ω(array).Should(Equal([]int64{0}))
	Ω(idx).Should(Equal(1))

	array, idx = RemoveIfPresentReturnIndex(array, 0)
	Ω(array).Should(Equal([]int64{}))
	Ω(idx).Should(Equal(0))

	array = []int64{0, 1, 2, 3}
	array = RemoveByIndex(array, 1)
	Ω(array).Should(Equal([]int64{0, 3, 2}))

	array = RemoveByIndex(array, 1)
	Ω(array).Should(Equal([]int64{0, 2}))

	array = RemoveByIndex(array, 1)
	Ω(array).Should(Equal([]int64{0}))

	array = RemoveByIndex(array, 1)
	Ω(array).Should(Equal([]int64{0}))

	array = RemoveByIndex(array, 0)
	Ω(array).Should(Equal([]int64{}))
}

func TestLeftShift(t *testing.T) {
	RegisterTestingT(t)

	array := []int64{0, 1, 2, 3}
	array = LeftShiftRemoveIfPresent(array, 1)
	Ω(array).Should(Equal([]int64{0, 2, 3}))

	array = LeftShiftRemoveIfPresent(array, 1)
	Ω(array).Should(Equal([]int64{0, 2, 3}))

	array = LeftShiftRemoveIfPresent(array, 0)
	Ω(array).Should(Equal([]int64{2, 3}))

	array = LeftShiftRemoveIfPresent(array, 3)
	Ω(array).Should(Equal([]int64{2}))

	array = LeftShiftRemoveIfPresent(array, 2)
	Ω(array).Should(Equal([]int64{}))

	array = []int64{0, 1, 2, 3}
	array = internalLeftShift(array, 1, 1)
	Ω(array).Should(Equal([]int64{0, 2, 3}))

	array = internalLeftShift(array, 1, 1)
	Ω(array).Should(Equal([]int64{0, 3}))

	array = internalLeftShift(array, 1, 1)
	Ω(array).Should(Equal([]int64{0}))

	array = internalLeftShift(array, 0, 1)
	Ω(array).Should(Equal([]int64{}))

	array = []int64{0, 1, 2, 3}
	array = leftShift0(array, 1)
	Ω(array).Should(Equal([]int64{0, 2, 3}))

	array = leftShift0(array, 1)
	Ω(array).Should(Equal([]int64{0, 3}))

	array = leftShift0(array, 1)
	Ω(array).Should(Equal([]int64{0}))

	array = leftShift0(array, 0)
	Ω(array).Should(Equal([]int64{}))
}

func TestLeftShift2(t *testing.T) {
	RegisterTestingT(t)

	array := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	array = LeftShift(array, 0, 1)
	Ω(array).Should(Equal([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9}))

	array = LeftShift(array, 9, 1)
	Ω(array).Should(Equal([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9}))

	array = LeftShift(array, 100, 1)
	Ω(array).Should(Equal([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9}))

	array = LeftShift(array, 1, 3)
	Ω(array).Should(Equal([]int64{1, 5, 6, 7, 8, 9}))

	array = LeftShift(array, 0, 100)
	Ω(array).Should(Equal([]int64{}))

	array = LeftShift(array, 0, 1)
	Ω(array).Should(Equal([]int64{}))
}

func TestLeftShift3(t *testing.T) {
	RegisterTestingT(t)

	array := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var removed []int64

	array, removed = LeftShiftReturnRemovedValues(array, 0, 1)
	Ω(array).Should(Equal([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	Ω(removed).Should(Equal([]int64{0}))

	array, removed = LeftShiftReturnRemovedValues(array, 9, 1)
	Ω(array).Should(Equal([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	Ω(removed).Should(BeNil())

	array, removed = LeftShiftReturnRemovedValues(array, 100, 1)
	Ω(array).Should(Equal([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	Ω(removed).Should(BeNil())

	array, removed = LeftShiftReturnRemovedValues(array, 1, 3)
	Ω(array).Should(Equal([]int64{1, 5, 6, 7, 8, 9}))
	Ω(removed).Should(Equal([]int64{2, 3, 4}))

	array, removed = LeftShiftReturnRemovedValues(array, 0, 100)
	Ω(array).Should(Equal([]int64{}))
	Ω(removed).Should(Equal([]int64{1, 5, 6, 7, 8, 9}))

	array, removed = LeftShiftReturnRemovedValues(array, 0, 1)
	Ω(array).Should(Equal([]int64{}))
	Ω(removed).Should(BeNil())
}

func leftShift0(array []int64, startIndex int) []int64 {

	n := len(array)
	for i := startIndex + 1; i < n; i++ {
		array[i-1] = array[i]
	}

	//array[n-1] = 0
	return array[:n-1]
}

func BenchmarkLeftShift(b *testing.B) {

	array := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for i := 0; i < b.N; i++ {
		array = internalLeftShift(array, 0, 1)
		array = append(array, int64(i))
	}
}

func BenchmarkLeftShift0(b *testing.B) {
	array := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for i := 0; i < b.N; i++ {
		array = leftShift0(array, 0)
		array = append(array, int64(i))
	}
}

func BenchmarkAddIfAbsent8(b *testing.B) {
	benchAddIfAbsent(b, 8)
}

func BenchmarkAddIfAbsent32(b *testing.B) {
	benchAddIfAbsent(b, 32)
}

func BenchmarkAddIfAbsent64(b *testing.B) {
	benchAddIfAbsent(b, 64)
}

func benchAddIfAbsent(b *testing.B, n int) {

	array := make([]int64, n)
	for i := 0; i < n; i++ {
		array[i] = int64(i)
	}

	var target []int64
	for i := 0; i < b.N; i++ {
		target = AddIfAbsent(target, array[i%n])
	}
}

func BenchmarkI64Map8(b *testing.B) {
	benchMap(b, 8)
}

func BenchmarkI64Map32(b *testing.B) {
	benchMap(b, 32)
}

func BenchmarkI64Map64(b *testing.B) {
	benchMap(b, 64)
}

func benchMap(b *testing.B, n int) {

	array := make([]int64, n)
	for i := 0; i < n; i++ {
		array[i] = int64(i)
	}

	target := make(map[int64]struct{})
	v := struct{}{}
	for i := 0; i < b.N; i++ {
		x := i % n
		if _, ok := target[array[x]]; !ok {
			target[array[x]] = v
		}

	}
}
