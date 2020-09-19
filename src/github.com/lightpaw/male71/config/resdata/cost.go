package resdata

import (
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/u64"
)

//gogen:config
type Cost struct {
	_     struct{} `file:"杂项/消耗.txt"`
	_     struct{} `proto:"shared_proto.CostProto"`
	Id    int      `protofield:"-"`
	Gold  uint64   `validator:"uint"`
	Food  uint64   `validator:"uint"`
	Wood  uint64   `validator:"uint"`
	Stone uint64   `validator:"uint"`

	Yuanbao  uint64 `validator:"uint"`
	Dianquan uint64 `validator:"uint"`
	Yinliang uint64 `validator:"uint"`

	GuildContributionCoin uint64 `validator:"uint"` // 帮派贡献币

	Jade    uint64 `validator:"uint"` // 玉璧
	JadeOre uint64 `validator:"uint"` // 玉石矿

	Goods      []*goods.GoodsData `head:"goods_id" validator:"int,duplicate" protofield:"GoodsId,config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`
	GoodsCount []uint64           `validator:"int,duplicate"`

	Gem      []*goods.GemData `head:"gem_id" validator:"int,duplicate" protofield:"GemId,config.U64a2I32a(goods.GetGemDataKeyArray(%s))"`
	GemCount []uint64         `validator:"int,duplicate"`

	IsNotEmpty bool `head:"-"`
	IsOnlyResource bool `head:"-" protofield:"-"`
}

func (c *Cost) Init(filename string) {
	check.PanicNotTrue(len(c.Goods) == len(c.GoodsCount), "%s 消耗配置%v 物品id跟物品个数必须一致", filename, c.Id)

	for i, g1 := range c.Goods {
		check.PanicNotTrue(c.GoodsCount[i] > 0, "%s 消耗配置%v 配置的物品个数必须 > 0", filename, c.Id)

		for j, g2 := range c.Goods {
			if i != j {
				check.PanicNotTrue(g1 != g2, "%s 消耗配置%v 物品id重复", filename, c.Id)
				check.PanicNotTrue(g1.Id != g2.Id, "%s 消耗配置%v 物品id重复", filename, c.Id)
			}
		}
	}

	for i, g1 := range c.Gem {
		check.PanicNotTrue(c.GemCount[i] > 0, "%s 奖励配置%v 配置的宝石个数必须 > 0", filename, c.Id)

		for j, g2 := range c.Gem {
			if i != j {
				check.PanicNotTrue(g1 != g2, "%s 消耗配置%v 宝石id重复", filename, c.Id)
				check.PanicNotTrue(g1.Id != g2.Id, "%s 消耗配置%v 宝石id重复", filename, c.Id)
			}
		}
	}

	c.IsNotEmpty = !c.IsEmpty()
	c.IsOnlyResource = c.isOnlyResource()
}

func (c *Cost) IsEmpty() bool {
	if c.Gold > 0 {
		return false
	}
	if c.Food > 0 {
		return false
	}
	if c.Wood > 0 {
		return false
	}
	if c.Stone > 0 {
		return false
	}
	if c.Yuanbao > 0 {
		return false
	}
	if c.Dianquan > 0 {
		return false
	}
	if c.Yinliang > 0 {
		return false
	}
	if c.GuildContributionCoin > 0 {
		return false
	}
	if c.Jade > 0 {
		return false
	}
	if c.JadeOre > 0 {
		return false
	}
	if len(c.Goods) > 0 {
		return false
	}
	if len(c.Gem) > 0 {
		return false
	}
	return true
}

func (c *Cost) isOnlyResource() bool {
	return c.HasResource() && c.Yuanbao <= 0 && c.Dianquan <= 0 && c.Yinliang <= 0 && c.GuildContributionCoin <= 0 && c.Jade <= 0 && c.JadeOre <= 0 && len(c.Goods) <= 0 && len(c.Gem) <= 0
}

func (c *Cost) TypeCount() (count int) {
	if c.Gold > 0 {
		count++
	}
	if c.Food > 0 {
		count++
	}
	if c.Wood > 0 {
		count++
	}
	if c.Stone > 0 {
		count++
	}
	if c.Yuanbao > 0 {
		count++
	}
	if c.Dianquan > 0 {
		count++
	}
	if c.Yinliang > 0 {
		count++
	}
	if c.GuildContributionCoin > 0 {
		count++
	}
	if c.Jade > 0 {
		count++
	}
	if c.JadeOre > 0 {
		count++
	}
	if len(c.Goods) > 0 {
		count++
	}
	if len(c.Gem) > 0 {
		count++
	}
	return
}

func (c *Cost) HasResource() bool {
	return c.Gold > 0 || c.Food > 0 || c.Wood > 0 || c.Stone > 0
}

func (c *Cost) GetResource() (uint64, uint64, uint64, uint64) {
	return c.Gold, c.Food, c.Wood, c.Stone
}

func (c *Cost) Multiple(m uint64) *Cost {
	if m == 1 {
		return c
	}
	return c.MultipleF64(float64(m))
}

func (c *Cost) MultipleF64(m float64) *Cost {
	if m == 1 {
		return c
	}

	out := &Cost{}
	out.Gold = multipleAmount(c.Gold, m)
	out.Food = multipleAmount(c.Food, m)
	out.Wood = multipleAmount(c.Wood, m)
	out.Stone = multipleAmount(c.Stone, m)

	out.Yuanbao = multipleAmount(c.Yuanbao, m)
	out.Dianquan = multipleAmount(c.Dianquan, m)
	out.Yinliang = multipleAmount(c.Yinliang, m)

	out.GuildContributionCoin = multipleAmount(c.GuildContributionCoin, m)

	out.Jade = multipleAmount(c.Jade, m)
	out.JadeOre = multipleAmount(c.JadeOre, m)

	out.Goods = c.Goods
	out.GoodsCount = u64.Copy(c.GoodsCount)

	for i, v := range out.GoodsCount {
		out.GoodsCount[i] = multipleAmount(v, m)
	}

	out.Gem = c.Gem
	out.GemCount = u64.Copy(c.GemCount)

	for i, v := range out.GemCount {
		out.GemCount[i] = multipleAmount(v, m)
	}

	out.IsNotEmpty = !out.IsEmpty()
	out.IsOnlyResource = !out.isOnlyResource()
	return out
}

func multipleAmount(amount uint64, multi float64) uint64 {
	if amount <= 0 {
		return 0
	}

	return u64.Max(1, u64.MultiF64(amount, multi))
}

func NewCostBuilder() *CostBuilder {
	return &CostBuilder{}
}

type CostBuilder struct {
	gold  uint64
	food  uint64
	wood  uint64
	stone uint64

	yuanbao  uint64
	dianquan uint64

	guildContributionCoin uint64

	jade    uint64
	jadeOre uint64

	goods map[*goods.GoodsData]uint64
	gem   map[*goods.GemData]uint64
}

func (b *CostBuilder) ReduceResource(gold, food, wood, stone uint64) *CostBuilder {
	b.gold = u64.Sub(b.gold, gold)
	b.food = u64.Sub(b.food, food)
	b.wood = u64.Sub(b.wood, wood)
	b.stone = u64.Sub(b.stone, stone)
	return b
}

func (b *CostBuilder) AddDianquan(toAdd uint64) *CostBuilder {
	b.dianquan += toAdd
	return b
}

func (b *CostBuilder) AddGoods(data *goods.GoodsData, count uint64) *CostBuilder {
	if b.goods == nil {
		b.goods = make(map[*goods.GoodsData]uint64, 1)
	}
	b.goods[data] = b.goods[data] + count
	return b
}

func (b *CostBuilder) Add(p *Cost) *CostBuilder {
	return b.AddMultiple(p, 1)
}

func (b *CostBuilder) AddMultiple(p *Cost, m uint64) *CostBuilder {

	b.gold += p.Gold * m
	b.food += p.Food * m
	b.wood += p.Wood * m
	b.stone += p.Stone * m

	b.yuanbao += p.Yuanbao * m
	b.dianquan += p.Dianquan * m

	b.guildContributionCoin += p.GuildContributionCoin * m

	b.jade += p.Jade * m
	b.jadeOre += p.JadeOre * m

	n := imath.Min(len(p.Goods), len(p.GoodsCount))
	for i := 0; i < n; i++ {
		b.AddGoods(p.Goods[i], p.GoodsCount[i]*m)
	}

	n = imath.Min(len(p.Gem), len(p.GemCount))
	for i := 0; i < n; i++ {
		b.gem[p.Gem[i]] = b.gem[p.Gem[i]] + p.GemCount[i]*m
	}

	return b
}

func (b *CostBuilder) Build() *Cost {
	p := &Cost{}

	p.Gold = b.gold
	p.Food = b.food
	p.Wood = b.wood
	p.Stone = b.stone

	p.Yuanbao = b.yuanbao
	p.Dianquan = b.dianquan

	p.GuildContributionCoin = b.guildContributionCoin

	p.Jade = b.jade
	p.JadeOre = b.jadeOre

	if len(b.goods) > 0 {
		p.Goods = make([]*goods.GoodsData, 0, len(b.goods))
		p.GoodsCount = make([]uint64, 0, len(b.goods))
		for k, v := range b.goods {
			p.Goods = append(p.Goods, k)
			p.GoodsCount = append(p.GoodsCount, v)
		}
	}

	if len(b.gem) > 0 {
		p.Gem = make([]*goods.GemData, 0, len(b.gem))
		p.GemCount = make([]uint64, 0, len(b.gem))
		for k, v := range b.gem {
			p.Gem = append(p.Gem, k)
			p.GemCount = append(p.GemCount, v)
		}
	}

	p.IsNotEmpty = !p.IsEmpty()
	p.IsOnlyResource = p.isOnlyResource()

	return p
}
