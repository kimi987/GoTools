package realmface

import "testing"

func BenchmarkDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bf(i)
	}
}

func bf(i int) {
	i += 1
}

func BenchmarkFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			bf(i)
		}()
	}
}

func BenchmarkFunc2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func(i int) {
			bf(i)
		}(i)
	}
}

func TestChan(t *testing.T) {
	c := make(chan struct{})

	close(c)

	<-c
	<-c

}
