package realmevent

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestEventQueue_NewEvent(t *testing.T) {
	RegisterTestingT(t)
	q := NewEventQueue()

	e1 := q.NewEvent(time.Unix(124, 0), nil)
	Ω(e1).ShouldNot(BeNil())
	Ω(q.Peek()).Should(BeIdenticalTo(e1))

	e := q.NewEvent(time.Unix(123, 0), nil)
	Ω(e).ShouldNot(BeNil())

	Ω(e.Time().Before(e1.Time())).Should(BeTrue())

	fmt.Println(q.Peek())
	Ω(q.Peek()).Should(BeIdenticalTo(e))

	e1.UpdateTime(time.Unix(120, 0))
	Ω(e1.Time()).Should(Equal(time.Unix(120, 0)))
	Ω(q.Peek()).Should(BeIdenticalTo(e1))

	e2 := q.NewEvent(time.Unix(130, 0), func() {})
	Ω(e2.Data()).ShouldNot(BeNil())
	Ω(q.Peek()).Should(BeIdenticalTo(e1))

	e1.RemoveFromQueue()
	Ω(q.Peek()).Should(BeIdenticalTo(e))
	e.RemoveFromQueue()
	Ω(q.Peek()).Should(BeIdenticalTo(e2))
	e2.RemoveFromQueue()
	Ω(q.Peek()).Should(BeNil())
}
