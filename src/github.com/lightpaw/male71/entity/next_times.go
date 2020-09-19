package entity

import (
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"math"
	"github.com/lightpaw/male7/util/u64"
)

func newNextTimes() *nexttimes {
	return &nexttimes{}
}

type nexttimes struct {
	times uint64

	nextTime time.Time
}

func (t *nexttimes) Times() uint64 {
	return t.times
}

func (t *nexttimes) IncreseTimes() uint64 {
	t.times ++
	return t.times
}

func (t *nexttimes) SetTimes(toSet uint64) {
	t.times = toSet
}

func (t *nexttimes) NextTime() time.Time {
	return t.nextTime
}

func (t *nexttimes) SetNextTime(toSet time.Time) {
	t.nextTime = toSet
}

func (t *nexttimes) TimesAndNextTime() (int32, int32) {
	return t.encodeClient()
}

func (t *nexttimes) encodeClient() (int32, int32) {
	return u64.Int32(t.times), timeutil.Marshal32(t.nextTime)
}

func (t *nexttimes) encodeServer() uint64 {
	return uint64(timeutil.Marshal64(t.nextTime)) | (t.times << 32)
}

func (t *nexttimes) unmarshal(amt uint64) {
	t.times = amt >> 32
	t.nextTime = timeutil.Unix64(int64(amt & math.MaxInt32))
}
