package shop

import (
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/pbutil"
)

var (
	pool           = pbutil.Pool
	newProtoMsg    = util.NewProtoMsg
	newCompressMsg = util.NewCompressMsg
	safeMarshal    = util.SafeMarshal
	_              = shared_proto.ErrIntOverflowConfig
)

type marshaler util.Marshaler

const (
	MODULE_ID = 20

	C2S_BUY_GOODS = 2

	C2S_BUY_BLACK_MARKET_GOODS = 9

	C2S_REFRESH_BLACK_MARKET_GOODS = 12
)

func NewS2cUpdateDailyShopGoodsMsg(id int32, count int32) pbutil.Buffer {
	msg := &S2CUpdateDailyShopGoodsProto{
		Id:    id,
		Count: count,
	}
	return NewS2cUpdateDailyShopGoodsProtoMsg(msg)
}

var s2c_update_daily_shop_goods = [...]byte{20, 1} // 1
func NewS2cUpdateDailyShopGoodsProtoMsg(object *S2CUpdateDailyShopGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_daily_shop_goods[:], "s2c_update_daily_shop_goods")

}

func NewS2cBuyGoodsMsg(id int32, count int32, multi int32, prize []byte) pbutil.Buffer {
	msg := &S2CBuyGoodsProto{
		Id:    id,
		Count: count,
		Multi: multi,
		Prize: prize,
	}
	return NewS2cBuyGoodsProtoMsg(msg)
}

var s2c_buy_goods = [...]byte{20, 3} // 3
func NewS2cBuyGoodsProtoMsg(object *S2CBuyGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_buy_goods[:], "s2c_buy_goods")

}

// 无效的商品id
var ERR_BUY_GOODS_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 20, 4, 1} // 4-1

// 无效的购买个数
var ERR_BUY_GOODS_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 20, 4, 6} // 4-6

// 购买数量超过限购数量
var ERR_BUY_GOODS_FAIL_LIMIT = pbutil.StaticBuffer{3, 20, 4, 2} // 4-2

// 购买消耗不足
var ERR_BUY_GOODS_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 20, 4, 3} // 4-3

// 不满足解锁条件
var ERR_BUY_GOODS_FAIL_LOCKED = pbutil.StaticBuffer{3, 20, 4, 4} // 4-4

// 服务器忙，请稍后再试
var ERR_BUY_GOODS_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 20, 4, 5} // 4-5

// vip等级不够
var ERR_BUY_GOODS_FAIL_VIP_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 20, 4, 7} // 4-7

// 今天不能买了
var ERR_BUY_GOODS_FAIL_DAILY_BOUGHT_LIMIT = pbutil.StaticBuffer{3, 20, 4, 8} // 4-8

// 当前不能使用该物品
var ERR_BUY_GOODS_FAIL_CANT_USE_GOODS = pbutil.StaticBuffer{3, 20, 4, 9} // 4-9

func NewS2cMultiCritBroadcastMsg(shop_type int32, multi int32, name string, prize []byte) pbutil.Buffer {
	msg := &S2CMultiCritBroadcastProto{
		ShopType: shop_type,
		Multi:    multi,
		Name:     name,
		Prize:    prize,
	}
	return NewS2cMultiCritBroadcastProtoMsg(msg)
}

var s2c_multi_crit_broadcast = [...]byte{20, 5} // 5
func NewS2cMultiCritBroadcastProtoMsg(object *S2CMultiCritBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_multi_crit_broadcast[:], "s2c_multi_crit_broadcast")

}

func NewS2cPushBlackMarketGoodsMsg(refrash bool, next_refresh_time int32, goods_id []int32, discount []int32, buy []bool) pbutil.Buffer {
	msg := &S2CPushBlackMarketGoodsProto{
		Refrash:         refrash,
		NextRefreshTime: next_refresh_time,
		GoodsId:         goods_id,
		Discount:        discount,
		Buy:             buy,
	}
	return NewS2cPushBlackMarketGoodsProtoMsg(msg)
}

var s2c_push_black_market_goods = [...]byte{20, 8} // 8
func NewS2cPushBlackMarketGoodsProtoMsg(object *S2CPushBlackMarketGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_push_black_market_goods[:], "s2c_push_black_market_goods")

}

func NewS2cBuyBlackMarketGoodsMsg(index int32) pbutil.Buffer {
	msg := &S2CBuyBlackMarketGoodsProto{
		Index: index,
	}
	return NewS2cBuyBlackMarketGoodsProtoMsg(msg)
}

var s2c_buy_black_market_goods = [...]byte{20, 10} // 10
func NewS2cBuyBlackMarketGoodsProtoMsg(object *S2CBuyBlackMarketGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_buy_black_market_goods[:], "s2c_buy_black_market_goods")

}

// 无效的商品下标
var ERR_BUY_BLACK_MARKET_GOODS_FAIL_INVALID_INDEX = pbutil.StaticBuffer{3, 20, 11, 1} // 11-1

// 购买物品，消耗不足
var ERR_BUY_BLACK_MARKET_GOODS_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 20, 11, 2} // 11-2

// 这个物品已经购买过了
var ERR_BUY_BLACK_MARKET_GOODS_FAIL_BUYED = pbutil.StaticBuffer{3, 20, 11, 3} // 11-3

var REFRESH_BLACK_MARKET_GOODS_S2C = pbutil.StaticBuffer{2, 20, 13} // 13

// 刷新消耗不足
var ERR_REFRESH_BLACK_MARKET_GOODS_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 20, 14, 1} // 14-1

// 刷新次数已达上限
var ERR_REFRESH_BLACK_MARKET_GOODS_FAIL_TIMES_LIMIT = pbutil.StaticBuffer{3, 20, 14, 2} // 14-2

func NewS2cUpdateVipShopGoodsMsg(id int32, count int32) pbutil.Buffer {
	msg := &S2CUpdateVipShopGoodsProto{
		Id:    id,
		Count: count,
	}
	return NewS2cUpdateVipShopGoodsProtoMsg(msg)
}

var s2c_update_vip_shop_goods = [...]byte{20, 15} // 15
func NewS2cUpdateVipShopGoodsProtoMsg(object *S2CUpdateVipShopGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_vip_shop_goods[:], "s2c_update_vip_shop_goods")

}
