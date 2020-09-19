package event

import (
	"fmt"
	"time"
)

func calQueueCount(n uint64) uint64 {
	// 2的倍数
	x := uint64(1)
	for i := 0; i < 12; i++ {
		if x >= n {
			return x
		}
		x *= 2
	}
	return x
}

func NewHashFuncQueue(n, sizePerQuene uint64, name string) *HashFuncQueue {

	count := calQueueCount(n)
	s := &HashFuncQueue{
		h:        count - 1,
		queueMap: make(map[uint64]*FuncQueue, count),
	}

	for i := uint64(0); i < count; i++ {
		s.queueMap[i] = NewFuncQueue(sizePerQuene, fmt.Sprintf("%s-%d", name, i))
	}

	return s
}

type HashFuncQueue struct {
	h        uint64
	queueMap map[uint64]*FuncQueue
}

func (s *HashFuncQueue) Close() {
	for _, queue := range s.queueMap {
		queue.Close()
	}
}

func (s *HashFuncQueue) hash(id int64) uint64 {
	return uint64(id) & s.h
}

func (s *HashFuncQueue) getQueue(id int64) *FuncQueue {
	return s.queueMap[s.hash(id)]
}

func (s *HashFuncQueue) TryFunc(id int64, f func()) bool {
	return s.getQueue(id).TryFunc(f)
}

func (s *HashFuncQueue) TimeoutFunc(id int64, f func(), timeout time.Duration) bool {
	return s.getQueue(id).TimeoutFunc(f, timeout)
}
