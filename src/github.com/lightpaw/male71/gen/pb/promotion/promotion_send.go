package promotion

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
	MODULE_ID = 43

	C2S_COLLECT_LOGIN_DAY_PRIZE = 4

	C2S_BUY_LEVEL_FUND = 7

	C2S_COLLECT_LEVEL_FUND = 10

	C2S_COLLECT_DAILY_SP = 13

	C2S_COLLECT_FREE_GIFT = 16

	C2S_BUY_TIME_LIMIT_GIFT = 21

	C2S_BUY_EVENT_LIMIT_GIFT = 26
)

func NewS2cCollectLoginDayPrizeMsg(day int32) pbutil.Buffer {
	msg := &S2CCollectLoginDayPrizeProto{
		Day: day,
	}
	return NewS2cCollectLoginDayPrizeProtoMsg(msg)
}

var s2c_collect_login_day_prize = [...]byte{43, 5} // 5
func NewS2cCollectLoginDayPrizeProtoMsg(object *S2CCollectLoginDayPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_login_day_prize[:], "s2c_collect_login_day_prize")

}

// 无效的天数
var ERR_COLLECT_LOGIN_DAY_PRIZE_FAIL_INVALID_DAY = pbutil.StaticBuffer{3, 43, 6, 1} // 6-1

// 登陆天数不足
var ERR_COLLECT_LOGIN_DAY_PRIZE_FAIL_DAY_NOT_ENOUGH = pbutil.StaticBuffer{3, 43, 6, 2} // 6-2

// 已经领取过奖励
var ERR_COLLECT_LOGIN_DAY_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 43, 6, 3} // 6-3

var BUY_LEVEL_FUND_S2C = pbutil.StaticBuffer{2, 43, 8} // 8

// 你已购买过了
var ERR_BUY_LEVEL_FUND_FAIL_BUYED = pbutil.StaticBuffer{3, 43, 9, 1} // 9-1

// 消耗不足
var ERR_BUY_LEVEL_FUND_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 43, 9, 2} // 9-2

func NewS2cCollectLevelFundMsg(level int32) pbutil.Buffer {
	msg := &S2CCollectLevelFundProto{
		Level: level,
	}
	return NewS2cCollectLevelFundProtoMsg(msg)
}

var s2c_collect_level_fund = [...]byte{43, 11} // 11
func NewS2cCollectLevelFundProtoMsg(object *S2CCollectLevelFundProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_level_fund[:], "s2c_collect_level_fund")

}

// 无效的君主等级
var ERR_COLLECT_LEVEL_FUND_FAIL_INVALID_LEVEL = pbutil.StaticBuffer{3, 43, 12, 1} // 12-1

// 你没有购买君主等级基金
var ERR_COLLECT_LEVEL_FUND_FAIL_NOT_BUY = pbutil.StaticBuffer{3, 43, 12, 2} // 12-2

// 君主等级不足
var ERR_COLLECT_LEVEL_FUND_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 43, 12, 3} // 12-3

// 领取过这个等级的奖励了
var ERR_COLLECT_LEVEL_FUND_FAIL_COLLECTED = pbutil.StaticBuffer{3, 43, 12, 4} // 12-4

func NewS2cCollectDailySpMsg(id int32, sp int32) pbutil.Buffer {
	msg := &S2CCollectDailySpProto{
		Id: id,
		Sp: sp,
	}
	return NewS2cCollectDailySpProtoMsg(msg)
}

var s2c_collect_daily_sp = [...]byte{43, 14} // 14
func NewS2cCollectDailySpProtoMsg(object *S2CCollectDailySpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_daily_sp[:], "s2c_collect_daily_sp")

}

// 无效的id
var ERR_COLLECT_DAILY_SP_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 43, 15, 1} // 15-1

// 已经领取过
var ERR_COLLECT_DAILY_SP_FAIL_COLLECTED = pbutil.StaticBuffer{3, 43, 15, 2} // 15-2

// 不在领取时间段内
var ERR_COLLECT_DAILY_SP_FAIL_OVERDUE = pbutil.StaticBuffer{3, 43, 15, 3} // 15-3

// 补领消耗不足
var ERR_COLLECT_DAILY_SP_FAIL_REPAIR_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 43, 15, 4} // 15-4

func NewS2cCollectFreeGiftMsg(id int32, prize []byte) pbutil.Buffer {
	msg := &S2CCollectFreeGiftProto{
		Id:    id,
		Prize: prize,
	}
	return NewS2cCollectFreeGiftProtoMsg(msg)
}

var s2c_collect_free_gift = [...]byte{43, 17} // 17
func NewS2cCollectFreeGiftProtoMsg(object *S2CCollectFreeGiftProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_free_gift[:], "s2c_collect_free_gift")

}

// 无效的id
var ERR_COLLECT_FREE_GIFT_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 43, 18, 1} // 18-1

// 已经领取过
var ERR_COLLECT_FREE_GIFT_FAIL_COLLECTED = pbutil.StaticBuffer{3, 43, 18, 2} // 18-2

func NewS2cNoticeTimeLimitGiftsMsg(gifts []*shared_proto.TimeLimitGiftProto) pbutil.Buffer {
	msg := &S2CNoticeTimeLimitGiftsProto{
		Gifts: gifts,
	}
	return NewS2cNoticeTimeLimitGiftsProtoMsg(msg)
}

var s2c_notice_time_limit_gifts = [...]byte{43, 24} // 24
func NewS2cNoticeTimeLimitGiftsProtoMsg(object *S2CNoticeTimeLimitGiftsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_notice_time_limit_gifts[:], "s2c_notice_time_limit_gifts")

}

func NewS2cBuyTimeLimitGiftMsg(grp_id int32, index int32) pbutil.Buffer {
	msg := &S2CBuyTimeLimitGiftProto{
		GrpId: grp_id,
		Index: index,
	}
	return NewS2cBuyTimeLimitGiftProtoMsg(msg)
}

var s2c_buy_time_limit_gift = [...]byte{43, 22} // 22
func NewS2cBuyTimeLimitGiftProtoMsg(object *S2CBuyTimeLimitGiftProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_buy_time_limit_gift[:], "s2c_buy_time_limit_gift")

}

// 无效的id
var ERR_BUY_TIME_LIMIT_GIFT_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 43, 23, 1} // 23-1

// 无法购买
var ERR_BUY_TIME_LIMIT_GIFT_FAIL_CANNOT_BUY = pbutil.StaticBuffer{3, 43, 23, 4} // 23-4

// 请等待下次刷新
var ERR_BUY_TIME_LIMIT_GIFT_FAIL_BOUGHT = pbutil.StaticBuffer{3, 43, 23, 2} // 23-2

// 礼包已刷新
var ERR_BUY_TIME_LIMIT_GIFT_FAIL_REFRESHED = pbutil.StaticBuffer{3, 43, 23, 3} // 23-3

// 元宝不足
var ERR_BUY_TIME_LIMIT_GIFT_FAIL_NOT_ENOUGH_YUANBAO = pbutil.StaticBuffer{3, 43, 23, 5} // 23-5

func NewS2cNoticeEventLimitGiftMsg(gift *shared_proto.EventLimitGiftProto) pbutil.Buffer {
	msg := &S2CNoticeEventLimitGiftProto{
		Gift: gift,
	}
	return NewS2cNoticeEventLimitGiftProtoMsg(msg)
}

var s2c_notice_event_limit_gift = [...]byte{43, 25} // 25
func NewS2cNoticeEventLimitGiftProtoMsg(object *S2CNoticeEventLimitGiftProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_notice_event_limit_gift[:], "s2c_notice_event_limit_gift")

}

func NewS2cBuyEventLimitGiftMsg(id int32) pbutil.Buffer {
	msg := &S2CBuyEventLimitGiftProto{
		Id: id,
	}
	return NewS2cBuyEventLimitGiftProtoMsg(msg)
}

var s2c_buy_event_limit_gift = [...]byte{43, 27} // 27
func NewS2cBuyEventLimitGiftProtoMsg(object *S2CBuyEventLimitGiftProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_buy_event_limit_gift[:], "s2c_buy_event_limit_gift")

}

// 无效的id
var ERR_BUY_EVENT_LIMIT_GIFT_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 43, 28, 1} // 28-1

// 无法购买
var ERR_BUY_EVENT_LIMIT_GIFT_FAIL_CANNOT_BUY = pbutil.StaticBuffer{3, 43, 28, 2} // 28-2

// 元宝不足
var ERR_BUY_EVENT_LIMIT_GIFT_FAIL_NOT_ENOUGH_YUANBAO = pbutil.StaticBuffer{3, 43, 28, 3} // 28-3
