package tickdata

import "time"

type TickTime interface {
	GetPrevTickTime() time.Time
	GetTickTime() time.Time
	Tick() <-chan struct{}
}

func New(tickTime time.Time, period time.Duration) *TickData {
	return NewTickTime(tickTime, tickTime.Add(-period))
}

func Copy(tickTime TickTime) *TickData {
	return NewTickTime(tickTime.GetTickTime(), tickTime.GetPrevTickTime())
}

func NewTickTime(tickTime, prevTickTime time.Time) *TickData {
	t := &TickData{}
	t.prevTickTime = prevTickTime
	t.tickTime = tickTime
	t.c = make(chan struct{})
	return t
}

type TickData struct {
	prevTickTime time.Time
	tickTime     time.Time

	c chan struct{}
}

func (t *TickData) GetPrevTickTime() time.Time {
	return t.prevTickTime
}

func (t *TickData) GetTickTime() time.Time {
	return t.tickTime
}

func (t *TickData) Close() {
	close(t.c)
}

func (t *TickData) Tick() <-chan struct{} {
	return t.c
}
