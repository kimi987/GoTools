package entity

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/shop"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/util/imath"
)

type hero_shop struct {
	// 物品限购数量
	dailyShopGoodsMap *heromap

	// 黑市商店数据
	blackMarketGoodsItems []*BlackMarketGoodsItem
}

func (hero *hero_shop) GetShopGoodsCount(shopGoodsId uint64) uint64 {
	return hero.dailyShopGoodsMap.internalMap[shopGoodsId]
}

func (hero *hero_shop) AddShopGoodsCount(shopGoodsId, count uint64) uint64 {
	c := hero.dailyShopGoodsMap.internalMap[shopGoodsId]
	c += count
	hero.dailyShopGoodsMap.internalMap[shopGoodsId] = c

	return c
}

func (hero *hero_shop) GetBlackMarketGoods(idx uint64) *BlackMarketGoodsItem {
	if idx < uint64(len(hero.blackMarketGoodsItems)) {
		return hero.blackMarketGoodsItems[idx]
	}
	return nil
}

func (hero *hero_shop) UpdateBlackMarketGoods(goods []*shop.BlackMarketGoodsData, discount []uint64) {

	n := imath.Min(len(goods), len(discount))
	items := make([]*BlackMarketGoodsItem, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, &BlackMarketGoodsItem{
			goods:    goods[i],
			discount: discount[i],
		})
	}

	hero.blackMarketGoodsItems = items
}

func (hero *hero_shop) RangeBlackMarketGoods(f func(item *BlackMarketGoodsItem) bool) {
	for _, item := range hero.blackMarketGoodsItems {
		if !f(item) {
			break
		}
	}
}

func (hero *hero_shop) encodeShopClient() *shared_proto.HeroShopProto {
	proto := &shared_proto.HeroShopProto{}
	proto.DailyShopGoods, proto.DailyBuyTimes = u64.Map2Int32Array(hero.dailyShopGoodsMap.internalMap)

	return proto
}

func (hero *hero_shop) encode() *server_proto.HeroShopServerProto {
	proto := &server_proto.HeroShopServerProto{}

	for _, v := range hero.blackMarketGoodsItems {
		proto.Item = append(proto.Item, &server_proto.BlackMarketGoodsItemServerProto{
			Goods:    v.goods.Id,
			Discount: v.discount,
			Buy:      v.buy,
		})
	}

	return proto
}

func (hero *hero_shop) unmarshal(proto *server_proto.HeroShopServerProto, datas *config.ConfigDatas, ) {
	if proto == nil {
		return
	}

	for _, v := range proto.Item {
		goods := datas.GetBlackMarketGoodsData(v.Goods)
		if goods == nil {
			continue
		}

		hero.blackMarketGoodsItems = append(hero.blackMarketGoodsItems, &BlackMarketGoodsItem{
			goods:    goods,
			discount: v.Discount,
			buy:      v.Buy,
		})
	}
}

type BlackMarketGoodsItem struct {
	// 黑市商品
	goods *shop.BlackMarketGoodsData

	// 折扣
	discount uint64

	// true表示已经购买过了
	buy bool
}

func (item *BlackMarketGoodsItem) Goods() *shop.BlackMarketGoodsData {
	return item.goods
}

func (item *BlackMarketGoodsItem) Discount() uint64 {
	return item.discount
}

func (item *BlackMarketGoodsItem) IsBuy() bool {
	return item.buy
}

func (item *BlackMarketGoodsItem) SetBuy() {
	item.buy = true
}
