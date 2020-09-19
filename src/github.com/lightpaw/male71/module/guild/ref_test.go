package guild

import (
	"fmt"
	"sync"
	"testing"
)

type sub struct {
	i int
}

type t1 struct {
	sub
}

type t2 struct {
	*sub
}

func TestRef(t *testing.T) {

	t1 := &t1{}

	fmt.Println(t1.sub)
	fmt.Println(t1.i)

	t2 := &t2{}
	t2.sub = &t1.sub

	fmt.Println(t2.sub)
	fmt.Println(t2.i)

	t1.i = 10
	fmt.Println(t2.i)
}

var ints = []int{101, 2934, 291, 302, 123, 123, 123, 5, 43, 2, 123, 45, 6, 1235}
var rw = sync.RWMutex{}

func BenchmarkNotLock(b *testing.B) {
	n := len(ints)
	sum := 0
	for i := 0; i < b.N; i++ {
		sum += ints[i%n]
	}
}

func BenchmarkRLock(b *testing.B) {
	n := len(ints)
	sum := 0
	for i := 0; i < b.N; i++ {
		rw.RLock()
		sum += ints[i%n]
		rw.RUnlock()
	}
}

func BenchmarkLock(b *testing.B) {
	n := len(ints)
	sum := 0
	for i := 0; i < b.N; i++ {
		rw.Lock()
		sum += ints[i%n]
		rw.Unlock()
	}
}
