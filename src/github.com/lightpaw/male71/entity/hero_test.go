package entity

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestMapArray(t *testing.T) {
	RegisterTestingT(t)

	array := []int{0, 1}
	array[0] += 10
	array[1] += 10
	array[0] += 10

	Ω(array[0]).Should(Equal(20))
	Ω(array[1]).Should(Equal(11))

	historyAmountMap := map[int]int{}

	historyAmountMap[0] += 10
	historyAmountMap[1] += 10
	historyAmountMap[0] += 10

	Ω(historyAmountMap[0]).Should(Equal(20))
	Ω(historyAmountMap[1]).Should(Equal(10))
}

func TestHero_EncodeCliente(t *testing.T) {
	RegisterTestingT(t)

	//id := int64(10928)
	//hero := NewHero(id, fmt.Sprintf("player_%d", id), &config.HeroInitData{})
	//hero.Domestic().ChangeRes(shared_proto.ResType_GOLD, 1000)
	//hero.Domestic().ChangeRes(shared_proto.ResType_FOOD, 1000)
	//hero.Domestic().ChangeRes(shared_proto.ResType_WOOD, 1000)
	//hero.Domestic().ChangeRes(shared_proto.ResType_STONE, 1000)
	//
	//Ω(hero.Domestic().GetGold()).Should(Equal(1000))
	//Ω(hero.Domestic().GetFood()).Should(Equal(1000))
	//Ω(hero.Domestic().GetWood()).Should(Equal(1000))
	//Ω(hero.Domestic().GetStone()).Should(Equal(1000))
	//
	//proto := hero.EncodeClient(time.Now().Unix64() * 1000)
	//Ω(proto.Id).Should(Equal(int64(10928)))
	//Ω(proto.Name).Should(Equal("player_10928"))
	//Ω(proto.GetDomestic().Gold).Should(Equal(int32(1000)))
	//Ω(proto.GetDomestic().Food).Should(Equal(int32(1000)))
	//Ω(proto.GetDomestic().Wood).Should(Equal(int32(1000)))
	//Ω(proto.GetDomestic().Stone).Should(Equal(int32(1000)))

}
