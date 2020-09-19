package pbutil

import (
	"github.com/lightpaw/logrus"
	"runtime/debug"
)

var (
	Pool = NewSyncPool(8, 1<<22, 4)

	_ Buffer = (*RecycleBuffer)(nil)

	_ Buffer = (StaticBuffer)([]byte{0})

	Empty = StaticBuffer(nil)
)

type Buffer interface {
	Buffer() []byte
	Free()
	DoFreeEvenItsStaticAndFuckMeIfItExplodes()
	Static() Buffer
	IsFreed() bool
}

type RecycleBuffer struct {
	buf []byte

	freed bool

	static bool
}

func (r *RecycleBuffer) Static() Buffer {
	r.static = true
	return r
}

func (r *RecycleBuffer) Buffer() []byte {
	return r.buf
}

func (r *RecycleBuffer) IsFreed() bool {
	return r.freed
}

func (r *RecycleBuffer) DoFreeEvenItsStaticAndFuckMeIfItExplodes() {
	if r.freed {
		logrus.WithField("stack", string(debug.Stack())).Error("重复free已经free的buf")
		return
	}
	r.freed = true
	r.static = false
	Pool.Free(r)
}

func (r *RecycleBuffer) Free() {
	if r.static {
		return
	}

	if r.freed {
		logrus.WithField("stack", string(debug.Stack())).Error("重复free已经free的buf")
		return
	}
	r.freed = true
	Pool.Free(r)
}

type StaticBuffer []byte

func (s StaticBuffer) Static() Buffer {
	return s
}

func (StaticBuffer) IsFreed() bool {
	return false
}

func (StaticBuffer) Free() {
}

func (StaticBuffer) DoFreeEvenItsStaticAndFuckMeIfItExplodes() {
}

func (s StaticBuffer) Buffer() []byte {
	return s
}
