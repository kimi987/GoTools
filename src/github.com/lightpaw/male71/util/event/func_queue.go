package event

import (
	"github.com/lightpaw/male7/util/call"
	"time"
)

func NewFuncQueue(n uint64, name string) *FuncQueue {
	q := &FuncQueue{
		name:           name,
		funcQueue:      make(chan func(), n),
		closeNotify:    make(chan struct{}),
		loopExitNotify: make(chan struct{}),
	}

	go call.CatchLoopPanic(q.loop, name)

	return q
}

type FuncQueue struct {
	name string

	funcQueue chan func()

	closeNotify    chan struct{}
	loopExitNotify chan struct{}
}

func (s *FuncQueue) Close() {
	close(s.closeNotify)
	<-s.loopExitNotify
}

func (s *FuncQueue) loop() {
	defer close(s.loopExitNotify)

	for {
		select {
		case f := <-s.funcQueue:
			call.CatchPanic(f, s.name)
		case <-s.closeNotify:
			return
		}
	}
}

func (s *FuncQueue) TryFunc(f func()) bool {
	select {
	case s.funcQueue <- f:
		return true
	default:
		return false
	}
}

func (s *FuncQueue) TimeoutFunc(f func(), timeout time.Duration) bool {
	if s.TryFunc(f) {
		return true
	}

	select {
	case s.funcQueue <- f:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (s *FuncQueue) MustFunc(f func()) {
	if !s.TryFunc(f) {
		go call.CatchPanic(f, s.name)
	}
}
