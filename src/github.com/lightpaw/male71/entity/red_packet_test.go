package entity

import (
	"fmt"
	"github.com/lightpaw/male7/config/red_packet"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestCreateRedPacket(t *testing.T) {
	RegisterTestingT(t)

	testWithErrCount(10000, 1, 20, true)
	testWithErrCount(100, 21, 10, false)
}

func testWithErrCount(allMoney, partMin uint64, count int, enoughCount bool) {
	for i := 0; i < 10000; i++ {
		parts := createParts(allMoney, partMin, u64.FromInt(count))
		Ω(len(parts)).Should(Equal(count))

		var sumMoney uint64
		for _, p := range parts {
			if enoughCount {
				Ω(p.money >= partMin).Should(BeTrue())
			}
			sumMoney += p.money
		}
		Ω(sumMoney).Should(Equal(allMoney))

		//mix(parts)
		//fmt.Printf("===== ")
		//for _, p := range parts {
		//	fmt.Printf("%v ", p.money)
		//}
		//fmt.Println()
	}
}

func TestRedPacket_Grab(t *testing.T) {
	RegisterTestingT(t)

	var count uint64 = 10
	ctime := time.Now()

	p, succ := CreateRedPacket(1, &red_packet.RedPacketData{Id: 1, Money: 100, MinPartMoney: 1}, count, ctime, "xxx", &shared_proto.HeroBasicProto{}, shared_proto.ChatType_ChatWorld)
	Ω(succ).Should(BeTrue())
	j, _ := p.BuildChatJson()
	fmt.Printf("json:%v\n", j)

	for i := 1; i < (int(count) + 10); i++ {
		money, allGrabbed := p.Grab(int64(i), &shared_proto.HeroBasicProto{}, ctime)
		fmt.Printf("i:%v money:%v allGrabbed:%v\n", i, money, allGrabbed)
		if uint64(i) < count {
			Ω(money > 0).Should(BeTrue())
			Ω(allGrabbed).Should(BeFalse())
		} else if uint64(i) == count {
			Ω(money > 0).Should(BeTrue())
			Ω(allGrabbed).Should(BeTrue())
		} else {
			Ω(money == 0).Should(BeTrue())
			Ω(allGrabbed).Should(BeTrue())
		}
	}
}
