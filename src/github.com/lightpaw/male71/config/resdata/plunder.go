package resdata

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/weight"
)

//gogen:config
type Plunder struct {
	_  struct{} `file:"杂项/掉落.txt"`
	Id uint64

	UnsafeGold  *data.RandAmount `default:" "`
	UnsafeFood  *data.RandAmount `default:" "`
	UnsafeWood  *data.RandAmount `default:" "`
	UnsafeStone *data.RandAmount `default:" "`

	SafeGold  *data.RandAmount `default:" "`
	SafeFood  *data.RandAmount `default:" "`
	SafeWood  *data.RandAmount `default:" "`
	SafeStone *data.RandAmount `default:" "`

	HeroExp    *data.RandAmount
	CaptainExp *data.RandAmount

	Item     []*PlunderItem
	ItemRate []*data.Rate `validator:"string,duplicate"`

	Group     []*PlunderGroup
	GroupRate []*data.Rate `validator:"string,duplicate"`
}

func (p *Plunder) Init(filename string, configs interface {
	GetPlunderItemArray() []*PlunderItem
}) {
	check.PanicNotTrue(len(p.Item) == len(p.ItemRate), "%s 掉落配置%v的item和item_rate 必须一致", filename, p.Id)

	check.PanicNotTrue(len(p.Group) == len(p.GroupRate), "%s 掉落配置的group和group_rate 必须一致")
}

func (p *Plunder) Try() *Prize {
	b := NewPrizeBuilder()
	b.AddUnsafeResource(
		p.UnsafeGold.Random(),
		p.UnsafeFood.Random(),
		p.UnsafeWood.Random(),
		p.UnsafeStone.Random(),
	)

	b.AddSafeResource(
		p.SafeGold.Random(),
		p.SafeFood.Random(),
		p.SafeWood.Random(),
		p.SafeStone.Random(),
	)

	b.AddHeroExp(p.HeroExp.Random())
	b.AddCaptainExp(p.CaptainExp.Random())

	n := imath.Min(len(p.Item), len(p.ItemRate))
	for i := 0; i < n; i++ {
		if p.ItemRate[i].Try() {
			p.Item[i].put(b)
		}
	}

	n = imath.Min(len(p.Group), len(p.GroupRate))
	for i := 0; i < n; i++ {
		if p.GroupRate[i].Try() {
			p.Group[i].put(b)
		}
	}

	return b.Build()
}

//gogen:config
type PlunderItem struct {
	_ struct{} `file:"杂项/掉落项.txt"`

	Id uint64

	// 掉什么，物品，装备，装备，掉落包，包中也分物品装备
	Goods     *goods.GoodsData     `default:"nullable"`
	Equipment *goods.EquipmentData `default:"nullable"`
	Gem       *goods.GemData       `default:"nullable"`
	Baowu     *BaowuData           `default:"nullable"`
	Captain   *ResCaptainData      `default:"nullable"`

	// 掉多少件
	Count uint64
}

func (item *PlunderItem) Init(filename string) {
	hasCount := 0
	if item.Goods != nil {
		hasCount++
	}
	if item.Equipment != nil {
		hasCount++
	}
	if item.Gem != nil {
		hasCount++
	}
	if item.Baowu != nil {
		hasCount++
	}
	if item.Captain != nil {
		hasCount ++
	}

	check.PanicNotTrue(hasCount == 1, "%s 掉落项配置%v 物品、装备、宝石不能同时配置", filename, item.Id)

}

func (item *PlunderItem) put(b *PrizeBuilder) {

	if item.Goods != nil {
		b.AddGoods(item.Goods, item.Count)
	}

	if item.Equipment != nil {
		b.AddEquipment(item.Equipment, item.Count)
	}

	if item.Gem != nil {
		b.AddGem(item.Gem, item.Count)
	}

	if item.Baowu != nil {
		b.AddBaowu(item.Baowu, item.Count)
	}

	if item.Captain != nil {
		b.AddCaptain(item.Captain, item.Count)
	}
}

//gogen:config
type PlunderGroup struct {
	_ struct{} `file:"杂项/掉落组.txt"`

	Id uint64

	Item     []*PlunderItem
	Weight   []uint64 `validator:",duplicate"`
	randomer *weight.WeightRandomer
}

func (g *PlunderGroup) Init(filename string) {

	check.PanicNotTrue(len(g.Item) == len(g.Weight), "%s 掉落配置%v的item和item_rate 必须一致", filename, g.Id)

	r, err := weight.NewWeightRandomer(g.Weight)
	if err != nil {
		logrus.WithError(err).Panicf("%s 掉落配置%v的weight无效，%v", filename, g.Id, g.Weight)
	}

	g.randomer = r
}

func (g *PlunderGroup) Try() *PlunderItem {
	return g.Item[g.randomer.RandomIndex()]
}

func (g *PlunderGroup) put(b *PrizeBuilder) {
	g.Try().put(b)
}
