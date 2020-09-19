package atomic

import (
	"testing"
	. "github.com/onsi/gomega"
)

type testvalue struct {
	Value
}

func TestValue(t *testing.T) {
	RegisterTestingT(t)

	v := &testvalue{}
	v.Store("haha")

	Ω(v.Load()).Should(Equal("haha"))
}
