package resdata

import (
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/u64"
	"sync"
)

//gogen:config
type Prize struct {
	_ struct{} `file:"杂项/奖励.txt"`
	_ struct{} `proto:"shared_proto.PrizeProto"`

	Id    int    `protofield:"-"`
	Gold  uint64 `head:"-,%s.SafeGold + %s.UnsafeGold"`
	Food  uint64 `head:"-,%s.SafeFood + %s.UnsafeFood"`
	Wood  uint64 `head:"-,%s.SafeWood + %s.UnsafeWood"`
	Stone uint64 `head:"-,%s.SafeStone + %s.UnsafeStone"`

	// proto里面只包含总的资源，和安全资源，这样做为了兼容客户端，客户端已经大量的使用了总资源作为资源数据

	// 安全资源
	SafeGold  uint64 `validator:"uint"`
	SafeFood  uint64 `validator:"uint"`
	SafeWood  uint64 `validator:"uint"`
	SafeStone uint64 `validator:"uint"`

	// 不安全资源
	UnsafeGold  uint64 `validator:"uint" protofield:"-"`
	UnsafeFood  uint64 `validator:"uint" protofield:"-"`
	UnsafeWood  uint64 `validator:"uint" protofield:"-"`
	UnsafeStone uint64 `validator:"uint" protofield:"-"`

	Goods      []*goods.GoodsData `head:"goods_id" protofield:"GoodsId,config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`
	GoodsCount []uint64           `validator:"int,duplicate,"`

	Equipment      []*goods.EquipmentData `validator:"int,duplicate" head:"equipment_id" protofield:"EquipmentId,config.U64a2I32a(goods.GetEquipmentDataKeyArray(%s))"`
	EquipmentCount []uint64               `validator:"int,duplicate"`

	Gem      []*goods.GemData `validator:"int,duplicate" head:"gem_id" protofield:"GemId,config.U64a2I32a(goods.GetGemDataKeyArray(%s))"`
	GemCount []uint64         `validator:"int,duplicate"`

	Baowu      []*BaowuData `validator:"int,duplicate" head:"baowu_id" protofield:"BaowuId,config.U64a2I32a(GetBaowuDataKeyArray(%s))"`
	BaowuCount []uint64     `validator:"int,duplicate"`

	Captain      []*ResCaptainData `validator:"int,duplicate" head:"captain_id" protofield:"CaptainId,config.U64a2I32a(GetResCaptainDataKeyArray(%s))"`
	CaptainCount []uint64          `validator:"int,duplicate"`

	CaptainExp uint64 `validator:"uint"`

	HeroExp uint64 `validator:"uint"`

	// 繁荣度
	Prosperity uint64 `validator:"uint" default:"0"`

	// 元宝
	Yuanbao uint64 `validator:"uint" default:"0"`

	// 点券
	Dianquan uint64 `validator:"uint" default:"0"`

	// 银两
	Yinliang uint64 `validator:"uint" default:"0"`

	// 联盟贡献
	GuildContributionCoin uint64 `validator:"uint" default:"0"`

	// 玉石矿
	Jade    uint64 `validator:"uint"`
	JadeOre uint64 `validator:"uint"`

	// VIP经验
	VipExp uint64 `validator:"uint" default:"0"`

	// 体力值
	Sp uint64 `validator:"uint" default:"0"`

	IsNotEmpty bool `head:"-"`

	Sort *AmountShowSortData `default:"nil" protofield:"AmountShowSortId,config.U64ToI32(%s.Id)"` // 展示数值排序

	prizeProto     *shared_proto.PrizeProto
	prizeProtoOnce sync.Once
}

func (c *Prize) Encode4Init() *shared_proto.PrizeProto {

	// 可能还没有执行Init方法，所以，这里手动做一下
	c.IsNotEmpty = c.TypeCount() > 0

	var i interface{} = c
	m, ok := i.(interface {
		Encode() *shared_proto.PrizeProto
	})
	if !ok {
		logrus.Errorf("Prize.Encode4Init() cast type fail")
	}

	return m.Encode()
}

func (c *Prize) String() string {
	return fmt.Sprintf("Prize: %v", *c)
}

func (c *Prize) PrizeProto() *shared_proto.PrizeProto {
	c.prizeProtoOnce.Do(func() {
		c.prizeProto = c.Encode4Init()
	})
	return c.prizeProto
}

// 返回奖励类型数量
func (c *Prize) TypeCount() (count int) {
	if c.UnsafeGold > 0 || c.SafeGold > 0 {
		count++
	}
	if c.UnsafeFood > 0 || c.SafeFood > 0 {
		count++
	}
	if c.UnsafeWood > 0 || c.SafeWood > 0 {
		count++
	}
	if c.UnsafeStone > 0 || c.SafeStone > 0 {
		count++
	}
	if c.HeroExp > 0 {
		count++
	}
	if c.CaptainExp > 0 {
		count++
	}
	if len(c.Goods) > 0 {
		count++
	}
	if len(c.Equipment) > 0 {
		count++
	}
	if len(c.Gem) > 0 {
		count++
	}
	if len(c.Baowu) > 0 {
		count++
	}
	if len(c.Captain) > 0 {
		count++
	}
	if c.Prosperity > 0 {
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
	if c.Sp > 0 {
		count++
	}
	return
}

func SetPrizeProtoIsNotEmpty(c *shared_proto.PrizeProto) {
	c.IsNotEmpty = PrizeProtoTypeCount(c) > 0
}

func PrizeProtoTypeCount(c *shared_proto.PrizeProto) (count int) {

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
	if c.HeroExp > 0 {
		count++
	}
	if c.CaptainExp > 0 {
		count++
	}
	if len(c.GoodsId) > 0 {
		count++
	}
	if len(c.EquipmentId) > 0 {
		count++
	}
	if len(c.GemId) > 0 {
		count++
	}
	if len(c.BaowuId) > 0 {
		count++
	}
	if len(c.CaptainId) > 0 {
		count++
	}
	if c.Prosperity > 0 {
		count++
	}
	if c.Yuanbao > 0 {
		count++
	}
	if c.Dianquan > 0 {
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
	if c.Sp > 0 {
		count++
	}
	return
}

func (c *Prize) GetSafeResource() (uint64, uint64, uint64, uint64) {
	return c.SafeGold, c.SafeFood, c.SafeWood, c.SafeStone
}

func (c *Prize) GetUnsafeResource() (uint64, uint64, uint64, uint64) {
	return c.UnsafeGold, c.UnsafeFood, c.UnsafeWood, c.UnsafeStone
}

func (c *Prize) Init(filename string) {

	check.PanicNotTrue(len(c.Goods) == len(c.GoodsCount), "%s 奖励配置%v 物品id跟物品个数必须一致", filename, c.Id)
	check.PanicNotTrue(len(c.Equipment) == len(c.EquipmentCount), "%s 奖励配置%v 物品id跟物品个数必须一致", filename, c.Id)
	check.PanicNotTrue(len(c.Gem) == len(c.GemCount), "%s 奖励配置%v 宝石id跟宝石个数必须一致", filename, c.Id)
	check.PanicNotTrue(len(c.Baowu) == len(c.BaowuCount), "%s 奖励配置%v 宝物id跟宝物个数必须一致", filename, c.Id)
	check.PanicNotTrue(len(c.Captain) == len(c.CaptainCount), "%s 奖励配置%v 将魂id跟武将个数必须一致", filename, c.Id)

	for i, v := range c.GoodsCount {
		check.PanicNotTrue(c.Goods[i] != nil, "%s 奖励配置%v 配置的物品不存在", filename, c.Id)
		check.PanicNotTrue(v > 0, "%s 奖励配置%v 配置的物品个数必须 > 0", filename, c.Id)
	}

	for i, v := range c.EquipmentCount {
		check.PanicNotTrue(c.Equipment[i] != nil, "%s 奖励配置%v 配置的装备不存在", filename, c.Id)
		check.PanicNotTrue(v > 0, "%s 奖励配置%v 配置的装备个数必须 > 0", filename, c.Id)
	}

	for i, v := range c.GemCount {
		check.PanicNotTrue(c.Gem[i] != nil, "%s 奖励配置%v 配置的宝石不存在", filename, c.Id)
		check.PanicNotTrue(v > 0, "%s 奖励配置%v 配置的宝石个数必须 > 0", filename, c.Id)
	}

	for i, v := range c.BaowuCount {
		check.PanicNotTrue(c.Baowu[i] != nil, "%s 奖励配置%v 配置的宝物不存在", filename, c.Id)
		check.PanicNotTrue(v > 0, "%s 奖励配置%v 配置的宝物个数必须 > 0", filename, c.Id)
	}

	for i, v := range c.CaptainCount {
		check.PanicNotTrue(c.Captain[i] != nil, "%s 奖励配置%v 配置的武将不存在", filename, c.Id)
		check.PanicNotTrue(v > 0, "%s 奖励配置%v 配置的武将个数必须 > 0", filename, c.Id)
	}

	c.IsNotEmpty = c.TypeCount() > 0
}

func (c *Prize) initTotalResource() {
	c.Gold = c.SafeGold + c.UnsafeGold
	c.Food = c.SafeFood + c.UnsafeFood
	c.Wood = c.SafeWood + c.UnsafeWood
	c.Stone = c.SafeStone + c.UnsafeStone
}

func (c *Prize) Multiple(m uint64) *Prize {
	if m == 1 {
		return c
	}

	out := &Prize{}

	out.SafeGold = c.SafeGold * m
	out.SafeFood = c.SafeFood * m
	out.SafeWood = c.SafeWood * m
	out.SafeStone = c.SafeStone * m

	out.UnsafeGold = c.UnsafeGold * m
	out.UnsafeFood = c.UnsafeFood * m
	out.UnsafeWood = c.UnsafeWood * m
	out.UnsafeStone = c.UnsafeStone * m

	out.initTotalResource()

	out.CaptainExp = c.CaptainExp * m
	out.HeroExp = c.HeroExp * m

	out.Prosperity = c.Prosperity * m

	out.Yuanbao = c.Yuanbao * m
	out.Dianquan = c.Dianquan * m
	out.Yinliang = c.Yinliang * m

	out.GuildContributionCoin = c.GuildContributionCoin * m

	out.Jade = c.Jade * m
	out.JadeOre = c.JadeOre * m
	out.Sp = c.Sp * m

	out.Goods = c.Goods
	out.GoodsCount = u64.Copy(c.GoodsCount)
	for i, v := range out.GoodsCount {
		out.GoodsCount[i] = v * m
	}

	out.Equipment = c.Equipment
	out.EquipmentCount = u64.Copy(c.EquipmentCount)
	for i, v := range out.EquipmentCount {
		out.EquipmentCount[i] = v * m
	}

	out.Gem = c.Gem
	out.GemCount = u64.Copy(c.GemCount)
	for i, v := range out.GemCount {
		out.GemCount[i] = v * m
	}

	out.Baowu = c.Baowu
	out.BaowuCount = u64.Copy(c.BaowuCount)
	for i, v := range out.BaowuCount {
		out.BaowuCount[i] = v * m
	}

	out.Captain = c.Captain
	out.CaptainCount = u64.Copy(c.CaptainCount)
	for i, v := range out.CaptainCount {
		out.CaptainCount[i] = v * m
	}

	out.IsNotEmpty = out.TypeCount() > 0

	return out
}

func (c *Prize) MultiCoef(m float64) *Prize {
	if m == 1 {
		return c
	}

	out := &Prize{}

	out.SafeGold = u64.MultiCoef(c.SafeGold, m)
	out.SafeFood = u64.MultiCoef(c.SafeFood, m)
	out.SafeWood = u64.MultiCoef(c.SafeWood, m)
	out.SafeStone = u64.MultiCoef(c.SafeStone, m)

	out.UnsafeGold = u64.MultiCoef(c.UnsafeGold, m)
	out.UnsafeFood = u64.MultiCoef(c.UnsafeFood, m)
	out.UnsafeWood = u64.MultiCoef(c.UnsafeWood, m)
	out.UnsafeStone = u64.MultiCoef(c.UnsafeStone, m)

	out.initTotalResource()

	out.CaptainExp = u64.MultiCoef(c.CaptainExp, m)
	out.HeroExp = u64.MultiCoef(c.HeroExp, m)

	out.Prosperity = u64.MultiCoef(c.Prosperity, m)

	out.Yuanbao = u64.MultiCoef(c.Yuanbao, m)
	out.Dianquan = u64.MultiCoef(c.Dianquan, m)
	out.Yinliang = u64.MultiCoef(c.Yinliang, m)

	out.GuildContributionCoin = u64.MultiCoef(c.GuildContributionCoin, m)

	out.Jade = u64.MultiCoef(c.Jade, m)
	out.JadeOre = u64.MultiCoef(c.JadeOre, m)
	out.Sp = u64.MultiCoef(c.Sp, m)

	out.Goods = c.Goods
	out.GoodsCount = u64.Copy(c.GoodsCount)
	for i, v := range out.GoodsCount {
		out.GoodsCount[i] = u64.MultiCoef(v, m)
	}

	out.Equipment = c.Equipment
	out.EquipmentCount = u64.Copy(c.EquipmentCount)
	for i, v := range out.EquipmentCount {
		out.EquipmentCount[i] = u64.MultiCoef(v, m)
	}

	out.Gem = c.Gem
	out.GemCount = u64.Copy(c.GemCount)
	for i, v := range out.GemCount {
		out.GemCount[i] = u64.MultiCoef(v, m)
	}

	out.Baowu = c.Baowu
	out.BaowuCount = u64.Copy(c.BaowuCount)
	for i, v := range out.BaowuCount {
		out.BaowuCount[i] = u64.MultiCoef(v, m)
	}

	out.Captain = c.Captain
	out.CaptainCount = u64.Copy(c.CaptainCount)
	for i, v := range out.CaptainCount {
		out.CaptainCount[i] = u64.MultiCoef(v, m)
	}

	out.IsNotEmpty = out.TypeCount() > 0

	return out
}

func UnmarshalPrize(proto *shared_proto.PrizeProto, config interface {
	GetGoodsData(uint64) *goods.GoodsData
	GetEquipmentData(uint64) *goods.EquipmentData
	GetResCaptainData(uint64) *ResCaptainData
	GetGemData(uint64) *goods.GemData
	GetBaowuData(uint64) *BaowuData
}) *Prize {
	c := &Prize{}

	c.SafeGold = u64.FromInt32(proto.SafeGold)
	c.SafeFood = u64.FromInt32(proto.SafeFood)
	c.SafeWood = u64.FromInt32(proto.SafeWood)
	c.SafeStone = u64.FromInt32(proto.SafeStone)

	c.UnsafeGold = u64.FromInt32(proto.Gold - proto.SafeGold)
	c.UnsafeFood = u64.FromInt32(proto.Food - proto.SafeFood)
	c.UnsafeWood = u64.FromInt32(proto.Wood - proto.SafeWood)
	c.UnsafeStone = u64.FromInt32(proto.Stone - proto.SafeStone)

	c.initTotalResource()

	c.CaptainExp = u64.FromInt32(proto.CaptainExp)
	c.HeroExp = u64.FromInt32(proto.HeroExp)

	c.Prosperity = u64.FromInt32(proto.Prosperity)

	c.Yuanbao = u64.FromInt32(proto.Yuanbao)
	c.Dianquan = u64.FromInt32(proto.Dianquan)
	c.Yinliang = u64.FromInt32(proto.Yinliang)

	c.GuildContributionCoin = u64.FromInt32(proto.GuildContributionCoin)

	c.Jade = u64.FromInt32(proto.Jade)
	c.JadeOre = u64.FromInt32(proto.JadeOre)
	c.Sp = u64.FromInt32(proto.Sp)

	n := imath.Min(len(proto.GoodsId), len(proto.GoodsCount))
	c.Goods = make([]*goods.GoodsData, 0, n)
	c.GoodsCount = make([]uint64, 0, n)
	for i := 0; i < n; i++ {
		g := config.GetGoodsData(u64.FromInt32(proto.GoodsId[i]))
		if g != nil {
			if ct := proto.GoodsCount[i]; ct <= 0 {
				continue
			}
			c.Goods = append(c.Goods, g)
			c.GoodsCount = append(c.GoodsCount, u64.FromInt32(proto.GoodsCount[i]))
		}
	}

	n = imath.Min(len(proto.EquipmentId), len(proto.EquipmentCount))
	c.Equipment = make([]*goods.EquipmentData, 0, n)
	c.EquipmentCount = make([]uint64, 0, n)
	for i := 0; i < n; i++ {
		e := config.GetEquipmentData(u64.FromInt32(proto.EquipmentId[i]))
		if e != nil {
			if ct := proto.EquipmentCount[i]; ct <= 0 {
				continue
			}
			c.Equipment = append(c.Equipment, e)
			c.EquipmentCount = append(c.EquipmentCount, u64.FromInt32(proto.EquipmentCount[i]))
		}
	}

	n = imath.Min(len(proto.GemId), len(proto.GemCount))
	c.Gem = make([]*goods.GemData, 0, n)
	c.GemCount = make([]uint64, 0, n)
	for i := 0; i < n; i++ {
		e := config.GetGemData(u64.FromInt32(proto.GemId[i]))
		if e != nil {
			if ct := proto.GemCount[i]; ct <= 0 {
				continue
			}
			c.Gem = append(c.Gem, e)
			c.GemCount = append(c.GemCount, u64.FromInt32(proto.GemCount[i]))
		}
	}

	n = imath.Min(len(proto.BaowuId), len(proto.BaowuCount))
	c.Baowu = make([]*BaowuData, 0, n)
	c.BaowuCount = make([]uint64, 0, n)
	for i := 0; i < n; i++ {
		e := config.GetBaowuData(u64.FromInt32(proto.BaowuId[i]))
		if e != nil {
			if ct := proto.BaowuCount[i]; ct <= 0 {
				continue
			}
			c.Baowu = append(c.Baowu, e)
			c.BaowuCount = append(c.BaowuCount, u64.FromInt32(proto.BaowuCount[i]))
		}
	}

	n = imath.Min(len(proto.CaptainId), len(proto.CaptainCount))
	c.Captain = make([]*ResCaptainData, 0, n)
	c.CaptainCount = make([]uint64, 0, n)
	for i := 0; i < n; i++ {
		d := config.GetResCaptainData(u64.FromInt32(proto.CaptainId[i]))
		if d != nil {
			if ct := proto.CaptainCount[i]; ct <= 0 {
				continue
			}
			c.Captain = append(c.Captain, d)
			c.CaptainCount = append(c.CaptainCount, u64.FromInt32(proto.CaptainCount[i]))
		}
	}

	c.IsNotEmpty = c.TypeCount() > 0

	return c
}

func AppendPrize(a ...*Prize) *Prize {

	switch len(a) {
	case 0:
		return nil
	case 1:
		return a[0]
	default:
		b := NewPrizeBuilder()
		for _, p := range a {
			if p != nil {
				b.Add(p)
			}
		}
		return b.Build()
	}

}

func NewPrizeBuilder() *PrizeBuilder {
	return &PrizeBuilder{}
}

type PrizeBuilder struct {
	safeGold  uint64
	safeFood  uint64
	safeWood  uint64
	safeStone uint64

	unsafeGold  uint64
	unsafeFood  uint64
	unsafeWood  uint64
	unsafeStone uint64

	captainExp uint64

	heroExp uint64

	prosperity uint64

	yuanbao  uint64
	dianquan uint64
	yinliang uint64

	guildContributionCoin uint64

	jade    uint64
	jadeOre uint64
	sp      uint64

	goods     map[*goods.GoodsData]uint64
	equipment map[*goods.EquipmentData]uint64
	gem       map[*goods.GemData]uint64
	baowu     map[*BaowuData]uint64
	captain   map[*ResCaptainData]uint64
}

func (c *PrizeBuilder) GetSafeResource() (uint64, uint64, uint64, uint64) {
	return c.safeGold, c.safeFood, c.safeWood, c.safeStone
}

func (c *PrizeBuilder) GetUnsafeResource() (uint64, uint64, uint64, uint64) {
	return c.unsafeGold, c.unsafeFood, c.unsafeWood, c.unsafeStone
}

func (b *PrizeBuilder) AddSafeResource(gold, food, wood, stone uint64) *PrizeBuilder {
	b.safeGold += gold
	b.safeFood += food
	b.safeWood += wood
	b.safeStone += stone

	return b
}

func (b *PrizeBuilder) AddUnsafeResource(gold, food, wood, stone uint64) *PrizeBuilder {
	b.unsafeGold += gold
	b.unsafeFood += food
	b.unsafeWood += wood
	b.unsafeStone += stone

	return b
}

func (b *PrizeBuilder) AddCaptainExp(toAdd uint64) *PrizeBuilder {
	b.captainExp += toAdd
	return b
}

func (b *PrizeBuilder) AddHeroExp(toAdd uint64) *PrizeBuilder {
	b.heroExp += toAdd
	return b
}

func (b *PrizeBuilder) AddProsperity(toAdd uint64) *PrizeBuilder {
	b.prosperity += toAdd
	return b
}

func (b *PrizeBuilder) AddYuanbao(toAdd uint64) *PrizeBuilder {
	b.yuanbao += toAdd
	return b
}

func (b *PrizeBuilder) AddDianquan(toAdd uint64) *PrizeBuilder {
	b.dianquan += toAdd
	return b
}

func (b *PrizeBuilder) AddYinliang(toAdd uint64) *PrizeBuilder {
	b.yinliang += toAdd
	return b
}

func (b *PrizeBuilder) AddJade(toAdd uint64) *PrizeBuilder {
	b.jade += toAdd
	return b
}

func (b *PrizeBuilder) AddJadeOre(toAdd uint64) *PrizeBuilder {
	b.jadeOre += toAdd
	return b
}

func (b *PrizeBuilder) AddGoods(data *goods.GoodsData, count uint64) *PrizeBuilder {
	if b.goods == nil {
		b.goods = make(map[*goods.GoodsData]uint64, 1)
	}

	b.goods[data] = b.goods[data] + count
	return b
}

func (b *PrizeBuilder) ReduceGoods(data *goods.GoodsData, toReduce uint64) *PrizeBuilder {
	if b.goods != nil {
		if count, exist := b.goods[data]; exist {
			newCount := u64.Sub(count, toReduce)
			if newCount > 0 {
				b.goods[data] = newCount
			}else{
				delete(b.goods, data)
			}
		}
	}
	return b
}

func (b *PrizeBuilder) AddEquipment(data *goods.EquipmentData, count uint64) *PrizeBuilder {
	if b.equipment == nil {
		b.equipment = make(map[*goods.EquipmentData]uint64, 1)
	}

	b.equipment[data] = b.equipment[data] + count
	return b
}

func (b *PrizeBuilder) AddGem(data *goods.GemData, count uint64) *PrizeBuilder {
	if b.gem == nil {
		b.gem = make(map[*goods.GemData]uint64, 1)
	}

	b.gem[data] = b.gem[data] + count
	return b
}

func (b *PrizeBuilder) AddBaowu(data *BaowuData, count uint64) *PrizeBuilder {
	if b.baowu == nil {
		b.baowu = make(map[*BaowuData]uint64, 1)
	}

	b.baowu[data] = b.baowu[data] + count
	return b
}

func (b *PrizeBuilder) AddCaptain(data *ResCaptainData, count uint64) *PrizeBuilder {
	if b.captain == nil {
		b.captain = make(map[*ResCaptainData]uint64, 1)
	}

	b.captain[data] = b.captain[data] + count
	return b
}

func (b *PrizeBuilder) Add(p *Prize) *PrizeBuilder {
	if p == nil {
		return b
	}
	return b.AddMultiple(p, 1)
}

func (b *PrizeBuilder) AddMultiple(p *Prize, m uint64) *PrizeBuilder {

	b.safeGold += p.SafeGold * m
	b.safeFood += p.SafeFood * m
	b.safeWood += p.SafeWood * m
	b.safeStone += p.SafeStone * m

	b.unsafeGold += p.UnsafeGold * m
	b.unsafeFood += p.UnsafeFood * m
	b.unsafeWood += p.UnsafeWood * m
	b.unsafeStone += p.UnsafeStone * m

	b.captainExp += p.CaptainExp * m
	b.heroExp += p.HeroExp * m

	b.prosperity += p.Prosperity * m

	b.yuanbao += p.Yuanbao * m
	b.dianquan += p.Dianquan * m
	b.yinliang += p.Yinliang * m

	b.guildContributionCoin += p.GuildContributionCoin * m

	b.jade += p.Jade * m
	b.jadeOre += p.JadeOre * m
	b.sp += p.Sp * m

	n := imath.Min(len(p.Goods), len(p.GoodsCount))
	for i := 0; i < n; i++ {
		b.AddGoods(p.Goods[i], p.GoodsCount[i]*m)
	}

	n = imath.Min(len(p.Equipment), len(p.EquipmentCount))
	for i := 0; i < n; i++ {
		b.AddEquipment(p.Equipment[i], p.EquipmentCount[i]*m)
	}

	n = imath.Min(len(p.Gem), len(p.GemCount))
	for i := 0; i < n; i++ {
		b.AddGem(p.Gem[i], p.GemCount[i]*m)
	}

	n = imath.Min(len(p.Baowu), len(p.BaowuCount))
	for i := 0; i < n; i++ {
		b.AddBaowu(p.Baowu[i], p.BaowuCount[i]*m)
	}

	n = imath.Min(len(p.Captain), len(p.CaptainCount))
	for i := 0; i < n; i++ {
		b.AddCaptain(p.Captain[i], p.CaptainCount[i]*m)
	}

	return b
}

func (b *PrizeBuilder) Build() *Prize {
	p := &Prize{}

	p.SafeGold = b.safeGold
	p.SafeFood = b.safeFood
	p.SafeWood = b.safeWood
	p.SafeStone = b.safeStone

	p.UnsafeGold = b.unsafeGold
	p.UnsafeFood = b.unsafeFood
	p.UnsafeWood = b.unsafeWood
	p.UnsafeStone = b.unsafeStone

	p.initTotalResource()

	p.CaptainExp = b.captainExp
	p.HeroExp = b.heroExp

	p.Prosperity = b.prosperity

	p.Yuanbao = b.yuanbao
	p.Dianquan = b.dianquan
	p.Yinliang = b.yinliang

	p.GuildContributionCoin = b.guildContributionCoin

	p.Jade = b.jade
	p.JadeOre = b.jadeOre
	p.Sp = b.sp

	if len(b.goods) > 0 {
		p.Goods = make([]*goods.GoodsData, 0, len(b.goods))
		p.GoodsCount = make([]uint64, 0, len(b.goods))
		for k, v := range b.goods {
			if v <= 0 {
				continue
			}
			p.Goods = append(p.Goods, k)
			p.GoodsCount = append(p.GoodsCount, v)
		}
	}

	if len(b.equipment) > 0 {
		p.Equipment = make([]*goods.EquipmentData, 0, len(b.equipment))
		p.EquipmentCount = make([]uint64, 0, len(b.equipment))
		for k, v := range b.equipment {
			if v <= 0 {
				continue
			}
			p.Equipment = append(p.Equipment, k)
			p.EquipmentCount = append(p.EquipmentCount, v)
		}
	}

	if len(b.gem) > 0 {
		p.Gem = make([]*goods.GemData, 0, len(b.gem))
		p.GemCount = make([]uint64, 0, len(b.gem))
		for k, v := range b.gem {
			if v <= 0 {
				continue
			}
			p.Gem = append(p.Gem, k)
			p.GemCount = append(p.GemCount, v)
		}
	}

	if len(b.baowu) > 0 {
		p.Baowu = make([]*BaowuData, 0, len(b.baowu))
		p.BaowuCount = make([]uint64, 0, len(b.baowu))
		for k, v := range b.baowu {
			if v <= 0 {
				continue
			}
			p.Baowu = append(p.Baowu, k)
			p.BaowuCount = append(p.BaowuCount, v)
		}
	}

	if len(b.captain) > 0 {
		p.Captain = make([]*ResCaptainData, 0, len(b.captain))
		p.CaptainCount = make([]uint64, 0, len(b.captain))
		for k, v := range b.captain {
			if v <= 0 {
				continue
			}
			p.Captain = append(p.Captain, k)
			p.CaptainCount = append(p.CaptainCount, v)
		}
	}

	p.IsNotEmpty = p.TypeCount() > 0

	return p
}

func GenGuildLevelPrizeId(groupId, guildLevel uint64) uint64 {
	return uint64(groupId)<<8 | guildLevel
}

//gogen:config
type GuildLevelPrize struct {
	_ struct{} `file:"杂项/联盟等级奖励.txt"`
	_ struct{} `proto:"shared_proto.GuildLevelPrizeProto"`

	Id         uint64 `head:"-,GenGuildLevelPrizeId(%s.GroupId%c %s.GuildLevel%c)" protofield:"-"`
	GroupId    uint64 `validator:"int>0"`
	GuildLevel uint64
	Prize      *Prize
}

func GuildLevelPrizeGroup(groupId uint64, configs interface {
	GetGuildLevelPrizeArray() []*GuildLevelPrize
}) (prizes []*GuildLevelPrize) {
	prizes = make([]*GuildLevelPrize, 0)

	for _, glp := range configs.GetGuildLevelPrizeArray() {
		if glp.GroupId == groupId {
			prizes = append(prizes, glp)
		}
	}
	return
}

func GuildLevelPrizeGroupMap(groupId uint64, configs interface {
	GetGuildLevelPrizeArray() []*GuildLevelPrize
}) (prizes map[uint64]*GuildLevelPrize) {
	prizes = make(map[uint64]*GuildLevelPrize)

	for _, glp := range configs.GetGuildLevelPrizeArray() {
		if glp.GroupId == groupId {
			prizes[glp.GuildLevel] = glp
		}
	}
	return
}
