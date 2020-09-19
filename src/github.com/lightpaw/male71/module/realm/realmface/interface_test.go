package realmface

import (
	. "github.com/onsi/gomega"
	"testing"
)

type obj struct {
	parent
}

type i2 interface {
	getInt2() string
}

type s2 interface {
	getString2() int
}

type parentImpl struct {
}

func (s *parentImpl) getInt() int {
	return 0
}

func (s *parentImpl) getInt2() string {
	return "1"
}

func (s *parentImpl) getString() string {
	return ""
}

type parentImpl2 struct {
}

func (s *parentImpl2) getInt() int {
	return 0
}

func (s *parentImpl2) getString2() int {
	return 2
}

func (s *parentImpl2) getString() string {
	return ""
}

type parent interface {
	sub
	getInt() int
}

type sub interface {
	getString() string
}

func TestObj(t *testing.T) {
	RegisterTestingT(t)

	o1 := &obj{
		parent: &parentImpl{},
	}

	var o interface{} = o1

	_, ok := o.(sub)
	Ω(ok).Should(BeTrue())

	_, ok = o.(parent)
	Ω(ok).Should(BeTrue())

	// 外面那个不行
	_, ok = o.(i2)
	Ω(ok).Should(BeFalse())

	_, ok = o.(interface {
		getInt2() string
	})
	Ω(ok).Should(BeFalse())

	// 里面的可以
	_, ok = o1.parent.(i2)
	Ω(ok).Should(BeTrue())

	_, ok = o1.parent.(interface {
		getInt2() string
	})
	Ω(ok).Should(BeTrue())

	op2 := &obj{
		parent: &parentImpl2{},
	}

	var o2 interface{} = op2

	_, ok = o2.(sub)
	Ω(ok).Should(BeTrue())

	_, ok = o2.(parent)
	Ω(ok).Should(BeTrue())

	// 外面那个不行
	_, ok = o2.(s2)
	Ω(ok).Should(BeFalse())

	_, ok = o2.(interface {
		getString2() int
	})
	Ω(ok).Should(BeFalse())

	// 里面的可以
	_, ok = op2.parent.(s2)
	Ω(ok).Should(BeTrue())

	_, ok = op2.parent.(interface {
		getString2() int
	})
	Ω(ok).Should(BeTrue())

	// 接口来试试
	checkI2(o1, false)
	checkI2(o1.parent, true)

	checkS2(op2, false)
	checkS2(op2.parent, true)

	// 接口的结果一样
}

func checkI2(parent parent, assert bool) {
	_, ok := parent.(i2)
	Ω(ok).Should(Equal(assert))
}

func checkS2(parent parent, assert bool) {
	_, ok := parent.(s2)
	Ω(ok).Should(Equal(assert))
}

func TestNilInterface(t *testing.T) {
	RegisterTestingT(t)

	var err error

	Ω(err).Should(BeNil())
	Ω(funcError()).Should(BeNil())

	Ω(err == nil).Should(BeTrue())
	Ω(funcError() == nil).Should(BeFalse()) // 看这里看这里看这里
}

func funcError() error {
	return funcError0()
}

func funcError0() *error0 {
	return nil
}

type error0 struct {
}

func (e *error0) Error() string {
	return ""
}

func TestInterface(t *testing.T) {
	RegisterTestingT(t)

	b1 := &a1{}
	b2 := &a2{}
	b3 := &a3{}

	Ω(b1.getInt2()).Should(Equal("a3"))
	Ω(b2.getInt2()).Should(Equal("a3"))
	Ω(b3.getInt2()).Should(Equal("a3"))

	var i interface{}

	i = b1
	_, ok := i.(i2)
	Ω(ok).Should(BeTrue())

	i = b2
	_, ok = i.(i2)
	Ω(ok).Should(BeTrue())

	i = b3
	_, ok = i.(i2)
	Ω(ok).Should(BeTrue())
}

type a1 struct {
	a2
}

type a2 struct {
	a3
}

type a3 struct {
}

func (a *a3) getInt2() string {
	return "a3"
}

func BenchmarkMap1(b *testing.B) {
	map1 := make(map[int]struct{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		map1[i&255] = struct{}{}
	}
}

func BenchmarkStruct(b *testing.B) {

	var s struct{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = struct{}{}
	}

	func(s struct{}) {
	}(s)
}

func BenchmarkMap1Exist(b *testing.B) {

	map1 := make(map[int]struct{})
	for i := 0; i < 65535/2; i++ {
		map1[i] = struct{}{}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, exist := map1[i&65535]; !exist {
		}
	}
}

func BenchmarkMap2(b *testing.B) {
	map2 := make(map[int]struct{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, exist := map2[i&255]; !exist {
			map2[i&255] = struct{}{}
		}
	}
}
