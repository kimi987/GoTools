package shop

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/weight"
	"sort"
	"time"
)

const (
	YuanbaoShopType = 2
	VIPShopType     = 6
)

//gogen:config
type ShopMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"商店/商店杂项.txt"`
	_ struct{} `protogen:"true"`

	RefreshBlackMarketCost []*resdata.Cost `validator:",duplicate"` // 刷新云游商人消耗

	AutoRefreshBlackMarketDuration []time.Duration                                 // 政务自动刷新间隔
	autoRefreshBlackMarketCycle    []*timeutil.CycleTime `head:"-" protofield:"-"` // 下次自动刷新的cycle
}

func (data *ShopMiscData) Init(filename string) {

	var autoRefreshCycle []*timeutil.CycleTime
	for _, d := range data.AutoRefreshBlackMarketDuration {
		autoRefreshCycle = append(autoRefreshCycle,
			timeutil.NewOffsetDailyTime(int64(d/time.Second)))
	}

	data.autoRefreshBlackMarketCycle = autoRefreshCycle
}

func (d *ShopMiscData) GetRefreshBlackMarketCost(times uint64, extraTimes uint64) *resdata.Cost {
	if times < uint64(len(d.RefreshBlackMarketCost)) {
		return d.RefreshBlackMarketCost[times]
	}

	if u64.Sub(times, extraTimes) < uint64(len(d.RefreshBlackMarketCost)) {
		return d.GetMaxRefreshBlackMarketCost()
	}

	return nil
}

func (d *ShopMiscData) GetMaxRefreshBlackMarketCost() *resdata.Cost {
	if len(d.RefreshBlackMarketCost) <= 0 {
		return nil
	}
	return d.RefreshBlackMarketCost[len(d.RefreshBlackMarketCost)-1]
}

func (data *ShopMiscData) NextAutoRefreshBlackMarketTime(ctime time.Time) time.Time {

	nextRefreshTime := ctime.Add(24 * time.Hour)
	for _, c := range data.autoRefreshBlackMarketCycle {
		t := c.NextTime(ctime)
		nextRefreshTime = timeutil.Min(nextRefreshTime, t)
	}

	return nextRefreshTime
}

// 商店
//gogen:config
type Shop struct {
	_ struct{} `file:"商店/商店.txt"`
	_ struct{} `proto:"shared_proto.ShopProto"`
	_ struct{} `protoconfig:"shop"`

	// 商店类型
	//1表示帮派商店
	//2表示元宝商店
	//3表示银两商店
	//4表示玉璧商店
	//5表示玲珑阁
	//6表示VIP商店
	Type uint64 `key:"true"`

	NormalGoods    []*NormalShopGoods
	ZhenBaoGeGoods []*ZhenBaoGeShopGoods
}

func (shop *Shop) HasGoods(id uint64) bool {
	for _, g := range shop.NormalGoods {
		if id == g.Id {
			return true
		}
	}
	for _, g := range shop.ZhenBaoGeGoods {
		if id == g.Id {
			return true
		}
	}

	return false
}

func (*Shop) InitAll(filename string, datas interface {
	GetShopArray() []*Shop
	GetNormalShopGoods(uint64) *NormalShopGoods
	GetZhenBaoGeShopGoodsArray() []*ZhenBaoGeShopGoods
}) {
	for _, shop := range datas.GetShopArray() {
		for _, goods := range shop.NormalGoods {
			check.PanicNotTrue(goods.ShopType == 0, "%s 中多个 shop_goods.txt 中相同id的物品 [%d] 被配置在了多个商店中!%d, %d", goods.Id, shop.Type, goods.ShopType)
			check.PanicNotTrue(goods.ShowSale >= 0 && goods.ShowSale <= 100, "%v 商品:%v 折扣show_sale:%v 必须 >= 0 且 <= 100", filename, goods.Id, goods.ShowSale)
			goods.ShopType = shop.Type
		}
		for _, goods := range shop.ZhenBaoGeGoods {
			check.PanicNotTrue(goods.ShopType == 0, "%s 中多个 shop_goods.txt 中相同id的物品 [%d] 被配置在了多个商店中!%d, %d", goods.Id, shop.Type, goods.ShopType)
			goods.ShopType = shop.Type
		}
	}

	for _, goods := range datas.GetZhenBaoGeShopGoodsArray() {
		check.PanicNotTrue(datas.GetNormalShopGoods(goods.Id) == nil, "存在相同的id[%d]配置在了 normal_shop_goods.txt 跟 zhen_bao_ge_shop_goods.txt 里面", goods.Id)
	}
}

type ShopGoodsSub struct {
	// 限购个数，0表示不限购
	CountLimit uint64 `validator:"uint"`

	// 解锁条件
	UnlockCondition *data.UnlockCondition `type:"sub"`

	// 免费次数
	FreeTimes uint64 `validator:"uint"`

	//UseImmediately bool `default:"false" protofield:"-"` // 购买后立即使用

	UseImmediatelyGoods []*goods.GoodsData `protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"` // 购买后立即使用

	ShopType uint64 `head:"-" protofield:"-"` // 商店类型
}

func (g *ShopGoodsSub) CanBuy(hasBuyCount uint64, tryBuyCount uint64) (errBuyTimes bool) {
	if g.CountLimit == 0 {
		return true
	}

	if u64.Plus(hasBuyCount, tryBuyCount) > g.FreeTimes+g.CountLimit {
		// 没次数了或者一次买太多了
		return false
	}

	return true
}

//gogen:config
type NormalShopGoods struct {
	_ struct{} `file:"商店/普通商品.txt"`
	_ struct{} `proto:"shared_proto.NormalShopGoodsProto"`

	Id uint64

	ShopGoodsSub `type:"sub"`

	Cost *resdata.Cost

	ShowPrize *resdata.Prize

	PlunderPrize *resdata.PlunderPrize `protofield:"Prize,%s.Prize.PrizeProto()"` // 掉落奖励

	Tag string `protofield:",shared_proto.ShopGoodsTag_value[%s]"`

	ShowSale uint64 `default:"0" validator:"uint"` // 展示用:打n/10折。0表示不打折，85表示打8.5折

	ShowOriginCost *resdata.Cost `default:"nullable"` // 展示原价

	VipDailyMaxCount uint64 `default:"0" validator:"uint"` // Vip商品每个人每日能买几个。0点重置

	GuildEventPrize *guild_data.GuildEventPrizeData `default:"nullable" protofield:",config.U64ToI32(%s.Id)"` // 联盟礼包
}

//gogen:config
type ZhenBaoGeShopGoods struct {
	_ struct{} `file:"商店/珍宝阁商品.txt"`
	_ struct{} `proto:"shared_proto.ZhenBaoGeShopGoodsProto"`

	Id uint64

	ShopGoodsSub `type:"sub"`

	// 付费次数消耗
	BuyCosts []*resdata.Cost

	MaxBuyCountPerTimes uint64 `validator:"int" default:"0" protofield:"-"`

	// 获取奖励数据，传入等级 level
	// 奖励
	ShowPrizes []*resdata.Prize   `protofield:"Prizes"`
	Prizes     []*resdata.Prize   `protofield:"-"`
	Plunders   []*resdata.Plunder `protofield:"-"`
	Levels     []uint64           `validator:"int"` // 数组长度跟奖励长度匹配>=该等级获得该奖励

	// 暴击权重跟倍率, 广播的最小倍率, >= 该倍率的都会被广播
	CritWeight        []uint64 `validator:"int" protofield:"-"`
	CritMulti         []uint64 `validator:"int" protofield:"-"`
	critRandomer      *weight.U64WeightRandomer // 可能为空，为空表示没有暴击
	BroadcastMinMulti uint64 `validator:"uint" protofield:"-"`
	Tag               string `protofield:",shared_proto.ShopGoodsTag_value[%s]"`
}

func (g *ZhenBaoGeShopGoods) Init2(filename string) {
	for _, cost := range g.BuyCosts {
		check.PanicNotTrue(cost.TypeCount() == 1, "%s 兑换配置消耗必须配置一种类型的消耗!%v", filename, g)
	}
	check.PanicNotTrue(len(g.BuyCosts) > 0, "%s 配置的购买消耗起码要配置一个!%d", filename, g.Id)

	if len(g.CritWeight) > 0 {
		critRandomer, err := weight.NewU64WeightRandomer(g.CritWeight, g.CritMulti)
		if err != nil {
			logrus.WithError(err).Panicf("%s 配置的暴击非法!", filename)
		}

		g.critRandomer = critRandomer
	}

	check.PanicNotTrue(sort.IsSorted(u64.Uint64Slice(g.Levels)), "%s 配置的等级(levels)应该逐级递增!%d", filename, g.Id)
	check.PanicNotTrue(len(g.Levels) == len(g.ShowPrizes) && len(g.ShowPrizes) > 0, "%s 配置的等级(levels)跟奖励(show_prizes)应该一一对应，且至少要配置一个!%d", filename, g.Id)

	if g.Prizes != nil {
		for _, prize := range g.Prizes {
			if len(g.UseImmediatelyGoods) > 0 {
				check.PanicNotTrue(prize.TypeCount() == 0, "%s 兑换并使用功能，奖励必须配空奖励（奖励里面什么都没有）!%v", filename, g)
			} else {
				check.PanicNotTrue(prize.TypeCount() == 1, "%s 兑换配置奖励必须配置一种类型的奖励!%v", filename, g)
			}
		}
		check.PanicNotTrue(len(g.Levels) == len(g.Prizes) && len(g.Prizes) > 0, "%s 配置的等级(levels)跟奖励(prizes)应该一一对应，且至少要配置一个!%d", filename, g.Id)
	} else {
		check.PanicNotTrue(g.Plunders != nil, "%s 配置的等级(levels)跟奖励(plunders)应该一一对应，且至少要配置一个!%d", filename, g.Id)
		check.PanicNotTrue(len(g.Levels) == len(g.Plunders) && len(g.Plunders) > 0, "%s 配置的等级(levels)跟奖励(plunders)应该一一对应，且至少要配置一个!%d", filename, g.Id)
	}

	check.PanicNotTrue(g.GetPrize(1, 1) != nil, "%s 配置的等级(levels)配置的在1级的时候取不到数据!%d", filename, g.Id)
}

func (g *ZhenBaoGeShopGoods) GetBuyCost(hasBuyCount uint64) *resdata.Cost {
	if hasBuyCount < g.FreeTimes {
		// 免费
		return nil
	}

	times := hasBuyCount - g.FreeTimes
	if times >= uint64(len(g.BuyCosts)) {
		return g.BuyCosts[len(g.BuyCosts)-1]
	}

	return g.BuyCosts[times]
}

func (g *ZhenBaoGeShopGoods) GetPrize(level, count uint64) *resdata.Prize {

	for i := len(g.Levels) - 1; i >= 0; i-- {
		if level >= g.Levels[i] {
			if g.Plunders != nil {
				plunder := g.Plunders[i]
				if count > 1 {
					builder := resdata.NewPrizeBuilder()
					for i := uint64(0); i < count; i++ {
						builder.Add(plunder.Try())
					}
					return builder.Build()
				}
				return plunder.Try()
			} else {
				prize := g.Prizes[i]
				if count > 1 {
					return prize.Multiple(count)
				}
				return prize
			}
		}
	}

	return nil
}

func (g *ZhenBaoGeShopGoods) RandomCrit() (multi uint64, needBroadcast bool) {
	if g.critRandomer == nil {
		return 1, false
	}

	multi = g.critRandomer.Random()
	needBroadcast = multi >= g.BroadcastMinMulti
	return
}

func ShopTypeGoodsId(shopType, goodsId uint64) uint64 {
	return goodsId<<4 | shopType
}
