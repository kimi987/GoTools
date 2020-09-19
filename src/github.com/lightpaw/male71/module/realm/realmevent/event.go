package realmevent

import (
	"container/heap"
	"github.com/lightpaw/logrus"
	"runtime/debug"
	"time"
)

var (
	_ EventQueue = (*eventQueue)(nil)
	_ Event      = (*event)(nil)
)

type EventQueue interface {
	NewEvent(time.Time, interface{}) Event
	Peek() Event
}

type Event interface {
	Time() time.Time
	Data() interface{}

	UpdateTime(time.Time)
	RemoveFromQueue()
}

func NewEventQueue() EventQueue {
	return &eventQueue{}
}

type eventQueue struct {
	heap          []*event
	nextEventTime <-chan time.Time
}

type event struct {
	time  time.Time // 触发的时间
	queue *eventQueue
	index int

	data interface{}
}

func (e *event) Time() time.Time {
	return e.time
}

func (e *event) Data() interface{} {
	return e.data
}

func (e *event) UpdateTime(t time.Time) {
	if e.queue == nil {
		logrus.WithField("stack", string(debug.Stack())).Error("已经被删除的event, 又调用了UpdateTime")
		return
	}

	e.queue.update(e, t)
}

func (e *event) RemoveFromQueue() {
	if e.queue == nil {
		logrus.WithField("stack", string(debug.Stack())).Error("已经被删除的event, 又调用了RemoveFromQueue")
		return
	}

	e.queue.removeEvent(e)
}

// --- event heap ---

func (pq *eventQueue) NewEvent(t time.Time, data interface{}) Event {
	e := &event{
		time:  t,
		queue: pq,
		data:  data,
	}

	heap.Push(pq, e)
	return e
}

func (pq *eventQueue) Len() int { return len(pq.heap) }

func (pq *eventQueue) Less(i, j int) bool {
	return pq.heap[i].time.Before(pq.heap[j].time)
}

func (pq *eventQueue) Swap(i, j int) {
	pq.heap[i], pq.heap[j] = pq.heap[j], pq.heap[i]
	pq.heap[i].index = i
	pq.heap[j].index = j
}

func (pq *eventQueue) Push(x interface{}) {
	n := len(pq.heap)
	item := x.(*event)
	item.index = n
	pq.heap = append(pq.heap, item)
}

func (pq *eventQueue) Pop() interface{} {
	old := pq.heap
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	item.queue = nil
	pq.heap = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *eventQueue) update(item *event, newTime time.Time) {
	item.time = newTime
	heap.Fix(pq, item.index)
}

func (pq *eventQueue) removeEvent(event *event) {
	heap.Remove(pq, event.index)
	event.index = -1
	event.queue = nil
}

func (pq *eventQueue) Peek() Event {
	if len(pq.heap) > 0 {
		return pq.heap[0]
	}

	return nil
}
