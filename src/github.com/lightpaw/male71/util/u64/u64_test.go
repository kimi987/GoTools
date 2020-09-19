package u64

import "testing"
import (
	"fmt"
	. "github.com/onsi/gomega"
)

func TestMulti(t *testing.T) {
	RegisterTestingT(t)

	f20 := 2.0
	f18 := 1.8
	f02 := f20 - f18
	f005 := f02 / 4

	fmt.Println(f20, f18, f02, f005)

	Ω(f02).ShouldNot(Equal(0.2))
	Ω(f1000(f02)).Should(Equal(uint64(200)))

	Ω(f1000(f005)).Should(Equal(uint64(50)))
	Ω(uint64(20 * f005)).ShouldNot(Equal(uint64(1)))
	Ω(MultiCoef(20, f005)).Should(Equal(uint64(1)))
}

func TestMultiF64(t *testing.T) {
	RegisterTestingT(t)

	f20 := 2.0
	f18 := 1.8
	f02 := f20 - f18
	f005 := f02 / 4

	fmt.Println(f20, f18, f02, f005)

	Ω(f02).ShouldNot(Equal(0.2))
	Ω(MultiF64(1000, f02)).Should(Equal(uint64(200)))

	Ω(MultiF64(1000, f005)).Should(Equal(uint64(50)))
	Ω(uint64(20 * f005)).ShouldNot(Equal(uint64(1)))
	Ω(MultiF64(20, f005)).Should(Equal(uint64(1)))

	Ω(MultiF64(0, f005)).Should(Equal(uint64(0)))
	Ω(MultiF64(100, 0)).Should(Equal(uint64(0)))
}

func TestLeftShift(t *testing.T) {

	var i32 int32 = -1
	var i int = -100
	var i64 int64 = -100

	fmt.Println(uint32(i32))
	fmt.Println(uint64(i32))

	fmt.Println(uint32(i))
	fmt.Println(uint64(i))

	fmt.Println(uint32(i64))
	fmt.Println(uint64(i64))

}
