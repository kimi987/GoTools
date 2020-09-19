package timer

//
// This code written by https://github.com/siddontang
// Copy from https://github.com/siddontang/go/blob/master/timingwheel/timingwheel.go
//

import (
	"sync"
	"time"
	"github.com/lightpaw/male7/util/call"
)

type TimingWheel struct {
	sync.Mutex
	interval   time.Duration
	ticker     *time.Ticker
	quit       chan struct{}
	maxTimeout time.Duration
	cs         []chan struct{}
	pos        int
}

func NewTimingWheel(interval time.Duration, buckets int) *TimingWheel {
	w := &TimingWheel{
		interval:   interval,
		quit:       make(chan struct{}),
		pos:        0,
		maxTimeout: time.Duration(interval * (time.Duration(buckets))),
		cs:         make([]chan struct{}, buckets),
		ticker:     time.NewTicker(interval),
	}
	for i := range w.cs {
		w.cs[i] = make(chan struct{})
	}
	go call.CatchLoopPanic(w.run, "TimingWheel")
	return w
}

func (w *TimingWheel) Stop() {
	close(w.quit)
}

func (w *TimingWheel) After(timeout time.Duration) <-chan struct{} {
	if timeout >= w.maxTimeout {
		panic("over maxtimeout")
	}
	w.Lock()
	index := (w.pos + int(timeout/w.interval)) % len(w.cs)
	b := w.cs[index]
	w.Unlock()
	return b
}

func (w *TimingWheel) run() {
	for {
		select {
		case <-w.ticker.C:
			w.onTicker()
		case <-w.quit:
			w.ticker.Stop()
			return
		}
	}
}

func (w *TimingWheel) onTicker() {
	w.Lock()
	lastC := w.cs[w.pos]
	w.cs[w.pos] = make(chan struct{})
	w.pos = (w.pos + 1) % len(w.cs)
	w.Unlock()
	close(lastC)
}
