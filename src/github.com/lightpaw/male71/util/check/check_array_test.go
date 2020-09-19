package check

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestName(t *testing.T) {
	RegisterTestingT(t)

	Ω(Int32Duplicate([]int32{})).Should(Equal(false))
	Ω(Int32Duplicate([]int32{1})).Should(Equal(false))
	Ω(Int32Duplicate([]int32{1, 2, 3})).Should(Equal(false))
	Ω(Int32Duplicate([]int32{1, 1})).Should(Equal(true))
	Ω(Int32Duplicate([]int32{2, 2})).Should(Equal(true))
	Ω(Int32Duplicate([]int32{1, 2, 1})).Should(Equal(true))
	Ω(Int32Duplicate([]int32{0, 0, 1})).Should(Equal(true))
	Ω(Int32Duplicate([]int32{1, 0, 1})).Should(Equal(true))

	Ω(Int32DuplicateIgnoreZero([]int32{})).Should(Equal(false))
	Ω(Int32DuplicateIgnoreZero([]int32{1})).Should(Equal(false))
	Ω(Int32DuplicateIgnoreZero([]int32{1, 2, 3})).Should(Equal(false))
	Ω(Int32DuplicateIgnoreZero([]int32{1, 1})).Should(Equal(true))
	Ω(Int32DuplicateIgnoreZero([]int32{2, 2})).Should(Equal(true))
	Ω(Int32DuplicateIgnoreZero([]int32{1, 2, 1})).Should(Equal(true))
	Ω(Int32DuplicateIgnoreZero([]int32{0, 0, 1})).Should(Equal(false))
	Ω(Int32DuplicateIgnoreZero([]int32{1, 0, 1})).Should(Equal(true))

	Ω(Int32AnyZero([]int32{})).Should(Equal(false))
	Ω(Int32AnyZero([]int32{0})).Should(Equal(true))
	Ω(Int32AnyZero([]int32{0, 0, 0, 0})).Should(Equal(true))
	Ω(Int32AnyZero([]int32{1})).Should(Equal(false))
	Ω(Int32AnyZero([]int32{1, 0})).Should(Equal(true))
	Ω(Int32AnyZero([]int32{0, 1, 1})).Should(Equal(true))

	Ω(IntDuplicate([]int{})).Should(Equal(false))
	Ω(IntDuplicate([]int{1})).Should(Equal(false))
	Ω(IntDuplicate([]int{1, 2, 3})).Should(Equal(false))
	Ω(IntDuplicate([]int{1, 1})).Should(Equal(true))
	Ω(IntDuplicate([]int{2, 2})).Should(Equal(true))
	Ω(IntDuplicate([]int{1, 2, 1})).Should(Equal(true))

	Ω(Int64Duplicate([]int64{})).Should(Equal(false))
	Ω(Int64Duplicate([]int64{1})).Should(Equal(false))
	Ω(Int64Duplicate([]int64{1, 2, 3})).Should(Equal(false))
	Ω(Int64Duplicate([]int64{1, 1})).Should(Equal(true))
	Ω(Int64Duplicate([]int64{2, 2})).Should(Equal(true))
	Ω(Int64Duplicate([]int64{1, 2, 1})).Should(Equal(true))

	Ω(Uint64Duplicate([]uint64{})).Should(Equal(false))
	Ω(Uint64Duplicate([]uint64{1})).Should(Equal(false))
	Ω(Uint64Duplicate([]uint64{1, 2, 3})).Should(Equal(false))
	Ω(Uint64Duplicate([]uint64{1, 1})).Should(Equal(true))
	Ω(Uint64Duplicate([]uint64{2, 2})).Should(Equal(true))
	Ω(Uint64Duplicate([]uint64{1, 2, 1})).Should(Equal(true))

	Ω(StringDuplicate([]string{})).Should(Equal(false))
	Ω(StringDuplicate([]string{"a"})).Should(Equal(false))
	Ω(StringDuplicate([]string{"", "a", "b", "c"})).Should(Equal(false))
	Ω(StringDuplicate([]string{"", ""})).Should(Equal(true))
	Ω(StringDuplicate([]string{"a", "b", "a"})).Should(Equal(true))
	Ω(StringDuplicate([]string{"i", "c", "b", "c"})).Should(Equal(true))
}
