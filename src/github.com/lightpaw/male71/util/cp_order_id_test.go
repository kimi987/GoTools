package util

import (
	"testing"
	"time"
	"fmt"
	. "github.com/onsi/gomega"
)

func TestCpOrderId(t *testing.T) {
	RegisterTestingT(t)
	//name := "M7"

	sid := uint32(10)

	heroId := int64(11928313)

	productId := uint64(133312312)

	money := uint64(64800)

	timestamp := time.Now().Unix()

	cpOrderId := NewCpOrderId(sid, heroId, productId, money, timestamp)
	fmt.Println(cpOrderId)

	version, randNum1, sid1, heroId1, productId1, moneyFen1, ctime1, err := ParseCpOrderId(cpOrderId)
	Ω(err).Should(Succeed())

	Ω(version).Should(BeEquivalentTo(1))
	Ω(randNum1 >> 25).Should(BeEquivalentTo(0))
	Ω(sid).Should(Equal(sid1))
	Ω(heroId).Should(Equal(heroId1))
	Ω(productId).Should(Equal(productId1))
	Ω(money).Should(Equal(moneyFen1))
	Ω(timestamp).Should(Equal(ctime1))

}
