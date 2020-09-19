package event

import (
	"testing"
	"time"
	"sync"
	"fmt"
	. "github.com/onsi/gomega"
)

func TestName(t *testing.T) {

	queue := NewEventQueue(100, time.Second, "test")

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		queue.TimeoutFunc(true, func() {
			panic("haha")
		})
		fmt.Println("done")
		wg.Done()
	}()

	wg.Wait()

	queue.Stop()

}

func BenchmarkName(b *testing.B) {

	queue := NewEventQueue(1, time.Second, "test")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		queue.Func(true, func() {})
	}
}

func TestCalQueueCount(t *testing.T) {

	for i := 0; i < 10; i++ {
		fmt.Println(calQueueCount(uint64(i)))
	}
	fmt.Println(calQueueCount(50))
	fmt.Println(calQueueCount(100))
	fmt.Println(calQueueCount(500))
	fmt.Println(calQueueCount(1000))
	fmt.Println(calQueueCount(5000))
	fmt.Println(calQueueCount(5000))
}

func TestHashFuncQueue(t *testing.T) {
	RegisterTestingT(t)

	hfq := NewHashFuncQueue(0, 0, "hero")
	Ω(hfq.h).Should(BeEquivalentTo(0))
	Ω(hfq.queueMap).Should(HaveLen(1))
	for i := 0; i < 10; i++ {
		Ω(hfq.getQueue(int64(i))).ShouldNot(BeNil())
		Ω(hfq.getQueue(int64(i))).Should(Equal(hfq.queueMap[uint64(i)%1]))
	}

	hfq = NewHashFuncQueue(1, 10, "hero")
	Ω(hfq.h).Should(BeEquivalentTo(0))
	Ω(hfq.queueMap).Should(HaveLen(1))
	for i := 0; i < 10; i++ {
		Ω(hfq.getQueue(int64(i))).ShouldNot(BeNil())
		Ω(hfq.getQueue(int64(i))).Should(Equal(hfq.queueMap[uint64(i)%1]))
	}

	hfq = NewHashFuncQueue(2, 10, "hero")
	Ω(hfq.h).Should(BeEquivalentTo(1))
	Ω(hfq.queueMap).Should(HaveLen(2))
	for i := 0; i < 10; i++ {
		Ω(hfq.getQueue(int64(i))).ShouldNot(BeNil())
		Ω(hfq.getQueue(int64(i))).Should(Equal(hfq.queueMap[uint64(i)%2]))
	}

	hfq = NewHashFuncQueue(3, 10, "hero")
	Ω(hfq.h).Should(BeEquivalentTo(3))
	Ω(hfq.queueMap).Should(HaveLen(4))
	for i := 0; i < 10; i++ {
		Ω(hfq.getQueue(int64(i))).ShouldNot(BeNil())
		Ω(hfq.getQueue(int64(i))).Should(Equal(hfq.queueMap[uint64(i)%4]))
	}

	hfq = NewHashFuncQueue(4, 10, "hero")
	Ω(hfq.h).Should(BeEquivalentTo(3))
	Ω(hfq.queueMap).Should(HaveLen(4))
	for i := 0; i < 10; i++ {
		Ω(hfq.getQueue(int64(i))).ShouldNot(BeNil())
		Ω(hfq.getQueue(int64(i))).Should(Equal(hfq.queueMap[uint64(i)%4]))
	}

	hfq.Close()
}

func BenchmarkHashQueue(b *testing.B) {

	queue := NewHashFuncQueue(4, 0, "test")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		queue.TimeoutFunc(int64(i), func() {
		}, time.Second)
	}
}
