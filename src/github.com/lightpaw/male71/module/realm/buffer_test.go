package realm

import (
	"github.com/lightpaw/pbutil"
	. "github.com/onsi/gomega"
	"testing"
)

func TestFree_Msg(t *testing.T) {
	RegisterTestingT(t)
	msg := pbutil.Buffer(pbutil.StaticBuffer([]byte{1, 2, 3}))

	freeMsg(&msg)
	Ω(msg).Should(BeNil())
}

func TestName(t *testing.T) {
	RegisterTestingT(t)
	tp := &troop{}

	base := tp.TargetBase()
	Ω(base).Should(BeNil())
	Ω(base == nil).Should(BeTrue())

	base = tp.StartingBase()
	Ω(base).Should(BeNil())
	Ω(base == nil).Should(BeTrue())
}
