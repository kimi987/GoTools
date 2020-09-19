package concurrent

import "testing"
import (
	"github.com/lightpaw/pbutil"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"strconv"
)

func TestI64BufferCacheMap(t *testing.T) {
	RegisterTestingT(t)

	buffs := []pbutil.Buffer{
		nil,
		pbutil.StaticBuffer([]byte("1")),
		pbutil.StaticBuffer([]byte("2")),
		pbutil.StaticBuffer([]byte("3")),
	}
	errs := []error{
		errors.New("error"),
	}

	f := func(key int64) (msg pbutil.Buffer, err error) {
		if len(buffs) > 0 {
			msg = buffs[0]
			buffs = buffs[1:]
		} else {
			msg = pbutil.StaticBuffer([]byte(strconv.FormatInt(key, 10)))
		}

		if len(errs) > 0 {
			err = errs[0]
			errs = errs[1:]
		}
		return
	}

	cacheMap := NewI64BufferCacheMap(f)

	// 第一次有错误
	version, msg, err := cacheMap.GetVersionBuffer(1)
	Ω(err).Should(HaveOccurred())
	Ω(msg).Should(BeNil())
	Ω(version).Should(Equal(uint64(1)))

	// 第二次正确
	version, msg, err = cacheMap.GetVersionBuffer(1)
	Ω(err).Should(Succeed())
	Ω(msg.Buffer()).Should(Equal([]byte("1")))
	Ω(version).Should(Equal(uint64(2)))

	// 再调一次，缓存中
	version, msg, err = cacheMap.GetVersionBuffer(1)
	Ω(err).Should(Succeed())
	Ω(msg.Buffer()).Should(Equal([]byte("1")))
	Ω(version).Should(Equal(uint64(2)))

	cacheMap.Clear(1) // version 3
	cacheMap.Clear(1) // version 3

	// 缓存清掉了
	version, msg, err = cacheMap.GetVersionBuffer(1)
	Ω(err).Should(Succeed())
	Ω(msg.Buffer()).Should(Equal([]byte("2")))
	Ω(version).Should(Equal(uint64(3)))

	// 再次调用
	version, msg, err = cacheMap.GetVersionBuffer(1)
	Ω(err).Should(Succeed())
	Ω(msg.Buffer()).Should(Equal([]byte("2")))
	Ω(version).Should(Equal(uint64(3)))

	// 新的key
	version, msg, err = cacheMap.GetVersionBuffer(2)
	Ω(err).Should(Succeed())
	Ω(msg.Buffer()).Should(Equal([]byte("3")))
	Ω(version).Should(Equal(uint64(4)))

	// 再次新的key
	version, msg, err = cacheMap.GetVersionBuffer(2)
	Ω(err).Should(Succeed())
	Ω(msg.Buffer()).Should(Equal([]byte("3")))
	Ω(version).Should(Equal(uint64(4)))

	// 清掉新的key
	cacheMap.Clear(2) // version 5

	// 不影响原来的key
	version, msg, err = cacheMap.GetVersionBuffer(1)
	Ω(err).Should(Succeed())
	Ω(msg.Buffer()).Should(Equal([]byte("2")))
	Ω(version).Should(Equal(uint64(3)))

	// 再次新的key
	version, msg, err = cacheMap.GetVersionBuffer(2)
	Ω(err).Should(Succeed())
	Ω(msg.Buffer()).Should(Equal([]byte("2")))
	Ω(version).Should(Equal(uint64(5)))
}
