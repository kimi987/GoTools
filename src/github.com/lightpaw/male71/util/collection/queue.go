package collection

import "github.com/eapache/queue"

func NewQueue() *Queue {
	return &Queue{
		Queue: queue.New(),
	}
}

type Queue struct {
	*queue.Queue
}

func (q *Queue) Range(f func(v interface{}) (toContinue bool)) {
	for i := 0; i < q.Length(); i++ {
		if !f(q.Get(i)) {
			break
		}
	}
}

func (q *Queue) RangeWithStartIndex(startIndex int, f func(v interface{}) (toContinue bool)) {
	for i := startIndex; i < q.Length(); i++ {
		if !f(q.Get(i)) {
			break
		}
	}
}

func (q *Queue) ReverseRange(f func(v interface{}) (toContinue bool)) {
	for i := 1; i <= q.Length(); i++ {
		if !f(q.Get(-i)) {
			break
		}
	}
}

func (q *Queue) ReverseRangeWithStartIndex(startIndex int, f func(v interface{}) (toContinue bool)) {
	for i := startIndex + 1; i <= q.Length(); i++ {
		if !f(q.Get(-i)) {
			break
		}
	}
}
