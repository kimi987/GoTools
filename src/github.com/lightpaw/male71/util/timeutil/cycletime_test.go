package timeutil

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestCycleTime(t *testing.T) {
	RegisterTestingT(t)

	ctime := time.Now().Truncate(time.Second)
	c := NewCycleTime(ctime.Unix(), 60)

	prev := ctime.Add(-time.Second)

	Ω(c.PrevTime(prev)).Should(Equal(ctime.Add(-60 * time.Second)))
	Ω(c.NextTime(prev)).Should(Equal(ctime))

	Ω(c.PrevTime(ctime)).Should(Equal(ctime))
	Ω(c.NextTime(ctime)).Should(Equal(ctime.Add(60 * time.Second)))

	Ω(c.Duration(ctime)).Should(Equal(0*time.Second))
	Ω(c.Duration(prev)).Should(Equal(59*time.Second))

	Ω(c.Duration(ctime.Add(30*time.Second))).Should(Equal(30*time.Second))
}
