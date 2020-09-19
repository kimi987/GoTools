package resdata

import "testing"
import (
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/must"
	. "github.com/onsi/gomega"
)

var goods1 *goods.GoodsData = &goods.GoodsData{Id: 1}
var goods2 *goods.GoodsData = &goods.GoodsData{Id: 2}

var equip1 *goods.EquipmentData = &goods.EquipmentData{Id: 1}
var equip2 *goods.EquipmentData = &goods.EquipmentData{Id: 2}

var gem1 *goods.GemData = &goods.GemData{Id: 1}
var gem2 *goods.GemData = &goods.GemData{Id: 2}

var captain1 *ResCaptainData = &ResCaptainData{Id: 1}
var captain2 *ResCaptainData = &ResCaptainData{Id: 2}

func TestNewPrizeBuilder(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	Expect(builder.safeGold).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(builder.safeFood).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(builder.safeWood).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(builder.safeStone).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")

	Expect(builder.unsafeGold).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(builder.unsafeFood).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(builder.unsafeWood).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(builder.unsafeStone).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")

	Expect(builder.captainExp).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(builder.heroExp).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(builder.yuanbao).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(builder.jade).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(builder.jadeOre).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(len(builder.goods)).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(len(builder.equipment)).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(len(builder.gem)).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
	Expect(len(builder.captain)).To(BeEquivalentTo(0), "初始化PrizeBuilder 的有问题？")
}

func TestPrizeBuilder_AddResource(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddSafeResource(1, 2, 3, 4)

	Expect(builder.safeGold).To(BeEquivalentTo(1), "加资源加 的有问题？")
	Expect(builder.safeFood).To(BeEquivalentTo(2), "加资源加 的有问题？")
	Expect(builder.safeWood).To(BeEquivalentTo(3), "加资源加 的有问题？")
	Expect(builder.safeStone).To(BeEquivalentTo(4), "加资源加 的有问题？")

	builder.AddUnsafeResource(11, 12, 13, 14)

	Expect(builder.unsafeGold).To(BeEquivalentTo(11), "加资源加 的有问题？")
	Expect(builder.unsafeFood).To(BeEquivalentTo(12), "加资源加 的有问题？")
	Expect(builder.unsafeWood).To(BeEquivalentTo(13), "加资源加 的有问题？")
	Expect(builder.unsafeStone).To(BeEquivalentTo(14), "加资源加 的有问题？")

	p := builder.Build()
	Ω(p.Gold).Should(BeEquivalentTo(12))
	Ω(p.Food).Should(BeEquivalentTo(14))
	Ω(p.Wood).Should(BeEquivalentTo(16))
	Ω(p.Stone).Should(BeEquivalentTo(18))

	g, f, w, s := p.GetSafeResource()
	Expect(g).To(BeEquivalentTo(1), "PrizeBuilder.Build() 的有问题？")
	Expect(f).To(BeEquivalentTo(2), "PrizeBuilder.Build() 的有问题？")
	Expect(w).To(BeEquivalentTo(3), "PrizeBuilder.Build() 的有问题？")
	Expect(s).To(BeEquivalentTo(4), "PrizeBuilder.Build() 的有问题？")

	g, f, w, s = p.GetUnsafeResource()
	Expect(g).To(BeEquivalentTo(11), "PrizeBuilder.Build() 的有问题？")
	Expect(f).To(BeEquivalentTo(12), "PrizeBuilder.Build() 的有问题？")
	Expect(w).To(BeEquivalentTo(13), "PrizeBuilder.Build() 的有问题？")
	Expect(s).To(BeEquivalentTo(14), "PrizeBuilder.Build() 的有问题？")
}

func TestPrizeBuilder_AddCaptainExp(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddCaptainExp(5)

	Expect(builder.captainExp).To(BeEquivalentTo(5), "加武将经验 的有问题？")

	builder.AddCaptainExp(15)

	Expect(builder.captainExp).To(BeEquivalentTo(20), "加武将经验 的有问题？")

	Expect(builder.Build().CaptainExp).To(BeEquivalentTo(20), "PrizeBuilder.Build() 的有问题？")
}

func TestPrizeBuilder_AddHeroExp(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddHeroExp(6)

	Expect(builder.heroExp).To(BeEquivalentTo(6), "加玩家经验 的有问题？")

	builder.AddHeroExp(15)

	Expect(builder.heroExp).To(BeEquivalentTo(21), "加玩家经验 的有问题？")

	Expect(builder.Build().HeroExp).To(BeEquivalentTo(21), "PrizeBuilder.Build() 的有问题？")
}

func TestPrizeBuilder_AddProsperity(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddProsperity(6)

	Expect(builder.prosperity).To(BeEquivalentTo(6), "加繁荣度 的有问题？")

	builder.AddProsperity(15)

	Expect(builder.prosperity).To(BeEquivalentTo(21), "加繁荣度 的有问题？")

	Expect(builder.Build().Prosperity).To(BeEquivalentTo(21), "PrizeBuilder.Build() 的有问题？")
}

func TestPrizeBuilder_AddYuanbao(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddYuanbao(6)

	Expect(builder.yuanbao).To(BeEquivalentTo(6), "加元宝 的有问题？")

	builder.AddYuanbao(15)

	Expect(builder.yuanbao).To(BeEquivalentTo(21), "加元宝 的有问题？")

	Expect(builder.Build().Yuanbao).To(BeEquivalentTo(21), "PrizeBuilder.Build() 的有问题？")
}

func TestPrizeBuilder_AddDianquan(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddDianquan(6)

	Expect(builder.dianquan).To(BeEquivalentTo(6), "加点券 的有问题？")

	builder.AddDianquan(15)

	Expect(builder.dianquan).To(BeEquivalentTo(21), "加点券 的有问题？")

	Expect(builder.Build().Dianquan).To(BeEquivalentTo(21), "PrizeBuilder.Build() 的有问题？")
}

func TestPrizeBuilder_AddYinliang(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddYinliang(6)

	Expect(builder.yinliang).To(BeEquivalentTo(6), "加yinliang 的有问题？")

	builder.AddYinliang(15)

	Expect(builder.yinliang).To(BeEquivalentTo(21), "加yinliang 的有问题？")

	Expect(builder.Build().Yinliang).To(BeEquivalentTo(21), "PrizeBuilder.Build() 的有问题？")
}

func TestPrizeBuilder_AddJade(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddJade(6)

	Expect(builder.jade).To(BeEquivalentTo(6), "加玉璧 的有问题？")

	builder.AddJade(15)

	Expect(builder.jade).To(BeEquivalentTo(21), "加玉璧 的有问题？")

	Expect(builder.Build().Jade).To(BeEquivalentTo(21), "PrizeBuilder.Build() 的有问题？")
}

func TestPrizeBuilder_AddJadeOre(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddJadeOre(6)

	Expect(builder.jadeOre).To(BeEquivalentTo(6), "加玉石矿 的有问题？")

	builder.AddJadeOre(15)

	Expect(builder.jadeOre).To(BeEquivalentTo(21), "加玉石矿 的有问题？")

	Expect(builder.Build().JadeOre).To(BeEquivalentTo(21), "PrizeBuilder.Build() 的有问题？")
}

func TestPrizeBuilder_AddGoods(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddGoods(goods1, 3)
	builder.AddGoods(goods2, 10)

	Expect(builder.goods[goods1]).To(BeEquivalentTo(3))
	Expect(builder.goods[goods2]).To(BeEquivalentTo(10))

	builder.AddGoods(goods1, 13)
	builder.AddGoods(goods2, 110)

	Expect(builder.goods[goods1]).To(BeEquivalentTo(16))
	Expect(builder.goods[goods2]).To(BeEquivalentTo(120))

	prize := builder.Build()

	Expect(len(prize.Goods)).To(BeEquivalentTo(2), "PrizeBuilder.Build() 物品在 build之后就数量就有问题了？？")

	for idx, goods := range prize.Goods {
		Expect(prize.GoodsCount[idx]).To(BeEquivalentTo(builder.goods[goods]), "PrizeBuilder.Build() 物品在 build之后就数量就有问题了？？")
	}
}

func TestPrizeBuilder_AddEquipment(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddEquipment(equip1, 3)
	builder.AddEquipment(equip2, 10)

	Expect(builder.equipment[equip1]).To(BeEquivalentTo(3))
	Expect(builder.equipment[equip2]).To(BeEquivalentTo(10))

	builder.AddEquipment(equip1, 13)
	builder.AddEquipment(equip2, 110)

	Expect(builder.equipment[equip1]).To(BeEquivalentTo(16))
	Expect(builder.equipment[equip2]).To(BeEquivalentTo(120))

	prize := builder.Build()

	Expect(len(prize.Equipment)).To(BeEquivalentTo(2), "PrizeBuilder.Build() 物品在 build之后就数量就有问题了？？")

	for idx, equip := range prize.Equipment {
		Expect(prize.EquipmentCount[idx]).To(BeEquivalentTo(builder.equipment[equip]), "PrizeBuilder.Build() 物品在 build之后就数量就有问题了？？")
	}
}

func TestPrizeBuilder_AddGem(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddGem(gem1, 3)
	builder.AddGem(gem2, 10)

	Expect(builder.gem[gem1]).To(BeEquivalentTo(3))
	Expect(builder.gem[gem2]).To(BeEquivalentTo(10))

	builder.AddGem(gem1, 13)
	builder.AddGem(gem2, 110)

	Expect(builder.gem[gem1]).To(BeEquivalentTo(16))
	Expect(builder.gem[gem2]).To(BeEquivalentTo(120))

	prize := builder.Build()

	Expect(len(prize.Gem)).To(BeEquivalentTo(2), "PrizeBuilder.Build() 宝石在 build之后就数量就有问题了？？")

	for idx, gem := range prize.Gem {
		Expect(prize.GemCount[idx]).To(BeEquivalentTo(builder.gem[gem]), "PrizeBuilder.Build() 宝石在 build之后就数量就有问题了？？")
	}
}

func TestPrizeBuilder_AddCaptain(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	builder.AddCaptain(captain1, 3)
	builder.AddCaptain(captain2, 10)

	Expect(builder.captain[captain1]).To(BeEquivalentTo(3))
	Expect(builder.captain[captain2]).To(BeEquivalentTo(10))

	builder.AddCaptain(captain1, 13)
	builder.AddCaptain(captain2, 110)

	Expect(builder.captain[captain1]).To(BeEquivalentTo(16))
	Expect(builder.captain[captain2]).To(BeEquivalentTo(120))

	prize := builder.Build()

	Expect(len(prize.Captain)).To(BeEquivalentTo(2), "PrizeBuilder.Build() 物品在 build之后就数量就有问题了？？")

	for idx, captain := range prize.Captain {
		Expect(prize.CaptainCount[idx]).To(BeEquivalentTo(builder.captain[captain]), "PrizeBuilder.Build() 物品在 build之后就数量就有问题了？？")
	}
}

var prize1 *Prize = &Prize{
	Gold:       1,
	Food:       2,
	SafeGold:   1,
	SafeFood:   1,
	UnsafeFood: 1,

	CaptainExp: 5,

	Goods:      []*goods.GoodsData{goods1, goods2},
	GoodsCount: []uint64{1, 2},

	Gem:      []*goods.GemData{gem1, gem2},
	GemCount: []uint64{1, 2},

	Captain:      []*ResCaptainData{captain1, captain2},
	CaptainCount: []uint64{1, 2},

	Jade: 3,
}

var prize2 *Prize = &Prize{
	Wood:        3,
	Stone:       4,
	SafeWood:    1,
	UnsafeWood:  2,
	UnsafeStone: 4,

	HeroExp: 6,

	Equipment:      []*goods.EquipmentData{equip1, equip2},
	EquipmentCount: []uint64{1, 2},

	Captain:      []*ResCaptainData{captain1, captain2},
	CaptainCount: []uint64{1, 2},
}

func TestCostBuilder_Add(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	prize := builder.Add(prize1).Build()
	testPrizeEqual(prize, prize1, 1)

	prize = builder.Add(prize2).Build()
	testPrizeEqual(prize, AppendPrize(prize1, prize2), 1)
}

func TestPrizeBuilder_AddMultiple(t *testing.T) {
	RegisterTestingT(t)

	builder := NewPrizeBuilder()

	prize := builder.AddMultiple(prize1, 33).Build()
	testPrizeEqual(prize, prize1, 33)
}

func testPrizeEqual(result, origin *Prize, m uint64) {
	Expect(result.Gold).To(BeEquivalentTo(origin.Gold * m))
	Expect(result.Food).To(BeEquivalentTo(origin.Food * m))
	Expect(result.Wood).To(BeEquivalentTo(origin.Wood * m))
	Expect(result.Stone).To(BeEquivalentTo(origin.Stone * m))

	Expect(result.SafeGold).To(BeEquivalentTo(origin.SafeGold * m))
	Expect(result.SafeFood).To(BeEquivalentTo(origin.SafeFood * m))
	Expect(result.SafeWood).To(BeEquivalentTo(origin.SafeWood * m))
	Expect(result.SafeStone).To(BeEquivalentTo(origin.SafeStone * m))

	Expect(result.UnsafeGold).To(BeEquivalentTo(origin.UnsafeGold * m))
	Expect(result.UnsafeFood).To(BeEquivalentTo(origin.UnsafeFood * m))
	Expect(result.UnsafeWood).To(BeEquivalentTo(origin.UnsafeWood * m))
	Expect(result.UnsafeStone).To(BeEquivalentTo(origin.UnsafeStone * m))

	Expect(result.CaptainExp).To(BeEquivalentTo(origin.CaptainExp * m))
	Expect(result.HeroExp).To(BeEquivalentTo(origin.HeroExp * m))
	Expect(result.Yuanbao).To(BeEquivalentTo(origin.Yuanbao * m))
	Expect(result.Dianquan).To(BeEquivalentTo(origin.Dianquan * m))
	Expect(result.Jade).To(BeEquivalentTo(origin.Jade * m))
	Expect(result.JadeOre).To(BeEquivalentTo(origin.JadeOre * m))

	Expect(len(result.Goods)).To(BeEquivalentTo(len(origin.Goods)))
	Expect(len(result.Equipment)).To(BeEquivalentTo(len(origin.Equipment)))
	Expect(len(result.Gem)).To(BeEquivalentTo(len(origin.GemCount)))
	Expect(len(result.Captain)).To(BeEquivalentTo(len(origin.Captain)))

	for i := 0; i < len(result.Goods); i++ {
		for idx, goods := range result.Goods {
			havSame := false
			for originIdx, originGoods := range origin.Goods {
				if goods == originGoods {
					Expect(result.GoodsCount[idx]).To(BeEquivalentTo(origin.GoodsCount[originIdx] * m))
					havSame = true
					break
				}
			}

			Expect(havSame).To(BeTrue())
		}
	}

	for i := 0; i < len(result.Equipment); i++ {
		for idx, equip := range result.Equipment {
			havSame := false
			for originIdx, originEquip := range origin.Equipment {
				if equip == originEquip {
					Expect(result.EquipmentCount[idx]).To(BeEquivalentTo(origin.EquipmentCount[originIdx] * m))
					havSame = true
					break
				}
			}

			Expect(havSame).To(BeTrue())
		}
	}

	for i := 0; i < len(result.Gem); i++ {
		for idx, gem := range result.Gem {
			havSame := false
			for originIdx, originGem := range origin.Gem {
				if gem == originGem {
					Expect(result.GemCount[idx]).To(BeEquivalentTo(origin.GemCount[originIdx] * m))
					havSame = true
					break
				}
			}

			Expect(havSame).To(BeTrue())
		}
	}

	for i := 0; i < len(result.Captain); i++ {
		for idx, equip := range result.Captain {
			havSame := false
			for originIdx, originEquip := range origin.Captain {
				if equip == originEquip {
					Expect(result.CaptainCount[idx]).To(BeEquivalentTo(origin.CaptainCount[originIdx] * m))
					havSame = true
					break
				}
			}

			Expect(havSame).To(BeTrue())
		}
	}
}

type configDatas struct {
}

func (c *configDatas) GetGoodsData(id uint64) *goods.GoodsData {
	if id == goods1.Id {
		return goods1
	} else if id == goods2.Id {
		return goods2
	} else {
		return nil
	}
}

func (c *configDatas) GetEquipmentData(id uint64) *goods.EquipmentData {
	if id == equip1.Id {
		return equip1
	} else if id == equip2.Id {
		return equip2
	} else {
		return nil
	}
}

func (c *configDatas) GetGemData(id uint64) *goods.GemData {
	if id == gem1.Id {
		return gem1
	} else if id == gem2.Id {
		return gem2
	} else {
		return nil
	}
}

func (c *configDatas) GetResCaptainData(id uint64) *ResCaptainData {
	if id == captain1.Id {
		return captain1
	} else if id == captain2.Id {
		return captain2
	} else {
		return nil
	}
}

func (c *configDatas) GetBaowuData(id uint64) *BaowuData {
	return nil
}

func TestPrize_Encode(t *testing.T) {
	RegisterTestingT(t)

	proto := prize1.Encode()
	bytes := must.Marshal(proto)

	prizeProto := &shared_proto.PrizeProto{}
	prizeProto.Unmarshal(bytes)

	configdatas := &configDatas{}
	prize := UnmarshalPrize(prizeProto, configdatas)
	testPrizeEqual(prize1, prize, 1)
}
