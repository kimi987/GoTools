package event

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/timer"
	"time"
	"github.com/lightpaw/male7/util/call"
)

func NewEventQueue(queneLength uint64, timeout time.Duration, name string) *EventQueue {
	if timeout <= 0 {
		logrus.Panicf("NewEventQueue timeout <= 0")
	}

	m := &EventQueue{}

	m.funcChan = make(chan *event_action, queneLength)
	m.closeNotifier = make(chan struct{})
	m.loopExitNotifier = make(chan struct{})

	interval := 500 * time.Millisecond
	buckets := int(i64.DivideTimes(timeout.Nanoseconds(), interval.Nanoseconds()) + 1)
	m.queueTimeoutWheel = timer.NewTimingWheel(interval, buckets)
	m.queueTimeout = timeout

	go call.CatchLoopPanic(m.loop, name)

	return m
}

type EventQueue struct {
	// 事件处理队列
	funcChan         chan *event_action
	closeNotifier    chan struct{}
	loopExitNotifier chan struct{}

	queueTimeoutWheel *timer.TimingWheel
	queueTimeout      time.Duration

	name string
}

// 放进去直到超时
func (r *EventQueue) TimeoutFunc(waitResult bool, f func()) (funcCalled bool) {
	e := &event_action{f: f, called: make(chan struct{})}

	select {
	case r.funcChan <- e:
		if waitResult {
			select {
			case <-r.loopExitNotifier:
				return false // main loop exit

			case <-e.called:
				return true
			}
		} else {
			return true // put success
		}

	case <-r.queueTimeoutWheel.After(r.queueTimeout):
		return false

	case <-r.closeNotifier:
		return false
	}
}

// 放进去位置
func (r *EventQueue) Func(waitResult bool, f func()) (funcCalled bool) {
	e := &event_action{f: f, called: make(chan struct{})}

	select {
	case r.funcChan <- e:
		if waitResult {
			select {
			case <-r.loopExitNotifier:
				return false // main loop exit

			case <-e.called:
				return true
			}
		} else {
			return true // put success
		}

	case <-r.closeNotifier:
		return false
	}
}

func (m *EventQueue) loop() {

	defer close(m.loopExitNotifier)

	for {
		select {
		case f := <-m.funcChan:
			call.CatchPanic(f.f, m.name)
			close(f.called)
		case <-m.closeNotifier:
			return
		}
	}

}

func (r *EventQueue) Stop() {
	close(r.closeNotifier)

	<-r.loopExitNotifier
}

type event_action struct {
	f      func()
	called chan struct{}
}
