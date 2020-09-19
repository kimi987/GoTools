package ticker

import (
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"github.com/lightpaw/male7/util/call"
	"sync/atomic"
	"time"
)

func NewTicker(ctime time.Time, delay, period time.Duration) *Ticker {
	t := &Ticker{}
	t.interval = period
	t.nextTickTimeRef = &atomic.Value{}
	t.nextTickTimeRef.Store(tickdata.New(ctime.Add(delay), period))
	t.quit = make(chan struct{})
	t.closed = make(chan struct{})

	go call.CatchLoopPanic(func() {
		t.loop(delay, period)
	}, "ticker")

	return t
}

type Ticker struct {
	//nextTickTime *timer.TickTime
	nextTickTimeRef *atomic.Value

	interval time.Duration

	quit   chan struct{}
	closed chan struct{}
}

func (t *Ticker) loop(delay, period time.Duration) {

	select {
	case <-time.After(delay):
		t.update()
	case <-t.quit:
		close(t.closed)
		return
	}

	periodTicker := time.NewTicker(period)
	for {
		select {
		case <-periodTicker.C:
			t.update()
		case <-t.quit:
			periodTicker.Stop()
			close(t.closed)
			return
		}
	}
}

func (t *Ticker) update() {
	tickTime := t.getData()
	nextTickTime := tickTime.GetTickTime().Add(t.interval)
	t.nextTickTimeRef.Store(tickdata.New(nextTickTime, t.interval))

	tickTime.Close()
}

func (t *Ticker) Stop() {
	close(t.quit)
	<-t.closed
}

func (t *Ticker) GetTickTime() tickdata.TickTime {
	return t.getData()
}

func (t *Ticker) getData() *tickdata.TickData {
	return t.nextTickTimeRef.Load().(*tickdata.TickData)
}
