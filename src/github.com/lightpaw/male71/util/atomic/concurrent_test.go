package atomic

import (
	"math/rand"
	"sync"
	"testing"
)

var success = []uint64{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
}

var ai = NewUint64(0)

func BenchmarkAtomic(b *testing.B) {
	ai.Store(0)
	for i := 0; i < b.N; i++ {
		for _, y := range success {
			for {
				x := ai.Load()
				if x < y {
					if ai.CAS(x, y) {
						break
					}
					continue
				}
				break

			}

		}
	}
}

func BenchmarkAtomic2(b *testing.B) {

	ai.Store(0)

	b.SetParallelism(16)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := ai.Load()
			y := x
			if rand.Intn(2) > 0 {
				y++
			}

			for {
				if ai.CAS(x, y) {
					break
				}

				x = ai.Load()
				if x <= y {
					break
				}
			}

		}
	})
}

var locker sync.Mutex

func BenchmarkLock(b *testing.B) {
	x := uint64(0)
	for i := 0; i < b.N; i++ {
		for _, y := range success {
			locker.Lock()

			if x < y {
				x = y
			}

			locker.Unlock()
		}
	}
}

func BenchmarkLock2(b *testing.B) {

	x := uint64(0)

	b.SetParallelism(16)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			locker.Lock()

			if rand.Intn(2) > 0 {
				x++
			}

			locker.Unlock()
		}
	})
}
