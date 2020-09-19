package concurrent

import (
	"github.com/lightpaw/pbutil"
	. "github.com/onsi/gomega"
	"math"
	"testing"
)

func TestBufferCache(t *testing.T) {
	RegisterTestingT(t)

	b1 := pbutil.StaticBuffer([]byte("b1"))
	b2 := pbutil.StaticBuffer([]byte("b2"))
	b3 := pbutil.StaticBuffer([]byte("b3"))

	arr := []pbutil.Buffer{
		b1, b2, b3,
	}

	f := func() (pbutil.Buffer, error) {
		if len(arr) > 0 {
			out := arr[0]
			arr = arr[1:]
			return out, nil
		}
		return nil, nil
	}

	c := NewBufferCache(f)

	b, err := c.Get()
	Ω(err).Should(Succeed())
	Ω(b.Buffer()).Should(Equal(b1.Buffer()))

	b, err = c.Get()
	Ω(err).Should(Succeed())
	Ω(b.Buffer()).Should(Equal(b1.Buffer()))

	b, err = c.Get()
	Ω(err).Should(Succeed())
	Ω(b.Buffer()).Should(Equal(b1.Buffer()))

	b, err = c.Get()
	Ω(err).Should(Succeed())
	Ω(b.Buffer()).Should(Equal(b1.Buffer()))

	b, err = c.Get()
	Ω(err).Should(Succeed())
	Ω(b.Buffer()).Should(Equal(b1.Buffer()))

	c.Clear()
	b, err = c.Get()
	Ω(err).Should(Succeed())
	Ω(b.Buffer()).Should(Equal(b2.Buffer()))

	b, err = c.Get()
	Ω(err).Should(Succeed())
	Ω(b.Buffer()).Should(Equal(b2.Buffer()))

	b, err = c.Get()
	Ω(err).Should(Succeed())
	Ω(b.Buffer()).Should(Equal(b2.Buffer()))

	c.Clear()
	b, err = c.Get()
	Ω(err).Should(Succeed())
	Ω(b.Buffer()).Should(Equal(b3.Buffer()))

	b, err = c.Get()
	Ω(err).Should(Succeed())
	Ω(b.Buffer()).Should(Equal(b3.Buffer()))

	c.Clear()
	b, err = c.Get()
	Ω(err).Should(Succeed())
	Ω(b).Should(BeNil())

}

func TestVersion(t *testing.T) {
	RegisterTestingT(t)

	Ω(I32Version(0)).Should(Equal(int32(0)))
	Ω(I32Version(1)).Should(Equal(int32(1)))
	Ω(I32Version(2)).Should(Equal(int32(2)))

	var version uint64 = math.MaxInt32
	Ω(I32Version(version)).Should(Equal(int32(math.MaxInt32)))
	Ω(I32Version(version + 1)).Should(Equal(int32(0)))
	Ω(I32Version(version + 2)).Should(Equal(int32(1)))
	Ω(I32Version(version + 3)).Should(Equal(int32(2)))
}
