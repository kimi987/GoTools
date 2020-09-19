package shop

import (
	"github.com/lightpaw/male7/util/weight"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/random"
)

const (
	BLACK_YUNYOU = 1
)

//gogen:config
type BlackMarketData struct {
	_ struct{} `file:"商店/黑市.txt"`

	Id uint64

	Group []*BlackMarketGoodsGroupData

	// 优先折扣
	MustDiscount      []float64 `protofield:"-"`
	MustDiscountCount []uint64  `validator:",duplicate" protofield:"-"`

	// 折扣
	Discount         []float64 `protofield:"-"`
	DiscountWeight   []uint64  `validator:",duplicate"  protofield:"-"`
	discountRandomer *weight.WeightRandomer
}

func (b *BlackMarketData) Init(filename string) {

	check.PanicNotTrue(len(b.MustDiscount) == len(b.MustDiscountCount),
		"%s 商店[%d]配置的优先折扣个数必须跟权重个数一致", filename, b.Id)

	check.PanicNotTrue(len(b.Discount) == len(b.DiscountWeight),
		"%s 商店[%d]配置的折扣个数必须跟权重个数一致", filename, b.Id)

	randomer, err := weight.NewWeightRandomer(b.DiscountWeight)
	check.PanicNotTrue(err == nil, "%s 折扣权重无效， err: %v", filename, err)
	b.discountRandomer = randomer
}

func (b *BlackMarketData) Random(heroLevel uint64) ([]*BlackMarketGoodsData, []uint64) {

	// 先获取到所有的物品id列表
	var goods []*BlackMarketGoodsData
	for _, group := range b.Group {
		goods = append(goods, group.RandomGoods(heroLevel)...)
	}

	discount := make([]uint64, len(goods))

	idx := 0
	index := random.NewIntIndexArray(len(goods))

	if idx < len(index) {
		// 先处理优先折扣
	out:
		for i, count := range b.MustDiscountCount {
			d := u64.MultiCoef(1000, b.MustDiscount[i])

			for i := uint64(0); i < count; i++ {
				discount[index[idx]] = d

				idx++
				if idx >= len(index) {
					break out
				}
			}
		}
	}

	// 没有折扣的，再随机一次
	if idx < len(index) {
		for i := idx; i < len(index); i++ {
			d := u64.MultiCoef(1000, b.Discount[b.discountRandomer.RandomIndex()])

			discount[index[i]] = d
		}
	}

	return goods, discount
}

//gogen:config
type BlackMarketGoodsGroupData struct {
	_ struct{} `file:"商店/黑市商品分组.txt"`

	Id uint64

	Count uint64

	Goods  []*BlackMarketGoodsData
	Weight []uint64 `validator:",duplicate"`

	maxRequiredHeroLevel uint64
}

func (g *BlackMarketGoodsGroupData) Init(filename string) {

	check.PanicNotTrue(len(g.Goods) == len(g.Weight), "%s 分组[%d]配置的商品个数必须跟权重个数一致", filename, g.Id)

	check.PanicNotTrue(g.Count <= uint64(len(g.Goods)), "%s 分组[%d]配置的随机个数，必须 <= 商品个数", filename, g.Id)

	var lv0Count uint64 = 0
	for _, v := range g.Goods {
		if v.RequiredHeroLevel <= 0 {
			lv0Count++
		}

		g.maxRequiredHeroLevel = u64.Max(g.maxRequiredHeroLevel, v.RequiredHeroLevel)
	}
	check.PanicNotTrue(g.Count <= lv0Count, "%s 分组[%d]配置的随机个数，必须 <= 不含等级限制的物品个数", filename, g.Id)
}

func (g *BlackMarketGoodsGroupData) RandomGoods(heroLevel uint64) []*BlackMarketGoodsData {

	var newGoods []*BlackMarketGoodsData
	var newWeight []uint64
	if g.maxRequiredHeroLevel <= heroLevel {
		newGoods = g.Goods
		newWeight = g.Weight
	} else {
		for i, v := range g.Goods {
			if v.RequiredHeroLevel <= heroLevel {
				newGoods = append(newGoods, v)
				newWeight = append(newWeight, g.Weight[i])
			}
		}
	}

	index, err := weight.RandomN(newWeight, int(g.Count))
	if err != nil {
		logrus.WithError(err).Error("BlackMarketGoodsGroup.RandomGoods() fail")
		return nil
	}

	array := make([]*BlackMarketGoodsData, 0, len(index))
	for _, i := range index {
		array = append(array, newGoods[i])
	}

	return array
}

//gogen:config
type BlackMarketGoodsData struct {
	_ struct{} `file:"商店/黑市商品.txt"`
	_ struct{} `protogen:"true"`

	// 商品id
	Id uint64

	// 商品奖励
	Prize *resdata.Prize `protofield:"-"`

	// 展示奖励
	ShowPrize *resdata.Prize

	// 商品价格
	Cost *resdata.Cost

	// 品质
	Quality shared_proto.Quality

	// 需要君主等级
	RequiredHeroLevel uint64 `validator:"uint" default:"0" protofield:"-"`
}

//gogen:config
type DiscountColorData struct {
	_ struct{} `file:"商店/折扣颜色.txt"`
	_ struct{} `protogen:"true"`

	// 折扣
	Discount uint64 `key:"true"`
	// 颜色
	Color string
}
