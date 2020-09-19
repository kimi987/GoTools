package fishing

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
	MODULE_ID = 16

	C2S_FISHING = 1

	C2S_FISH_POINT_EXCHANGE = 8

	C2S_SET_FISHING_CAPTAIN = 11
)

func NewS2cFishingMsg(fishing_result [][]byte, have_soul_to_goods []bool, show_index []int32, times int32, fish_type int32, fishing_times int32, next_time int32) pbutil.Buffer {
	msg := &S2CFishingProto{
		FishingResult:   fishing_result,
		HaveSoulToGoods: have_soul_to_goods,
		ShowIndex:       show_index,
		Times:           times,
		FishType:        fish_type,
		FishingTimes:    fishing_times,
		NextTime:        next_time,
	}
	return NewS2cFishingProtoMsg(msg)
}

var s2c_fishing = [...]byte{16, 2} // 2
func NewS2cFishingProtoMsg(object *S2CFishingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fishing[:], "s2c_fishing")

}

// 钱不够
var ERR_FISHING_FAIL_RES_NOT_ENOUGH = pbutil.StaticBuffer{3, 16, 3, 1} // 3-1

// 消耗没找到
var ERR_FISHING_FAIL_COST_NOUT_FOUND = pbutil.StaticBuffer{3, 16, 3, 3} // 3-3

// 无效的钓鱼类型
var ERR_FISHING_FAIL_INVALID_FISH_TYPE = pbutil.StaticBuffer{3, 16, 3, 4} // 3-4

// 钓鱼倒计时未结束
var ERR_FISHING_FAIL_COUNTDOWN = pbutil.StaticBuffer{3, 16, 3, 5} // 3-5

// 钓鱼次数达到当日上限
var ERR_FISHING_FAIL_DAILY_TIMES_LIMIT = pbutil.StaticBuffer{3, 16, 3, 6} // 3-6

func NewS2cFishingBroadcastMsg(id []byte, name string, flagname string, prize []byte) pbutil.Buffer {
	msg := &S2CFishingBroadcastProto{
		Id:       id,
		Name:     name,
		Flagname: flagname,
		Prize:    prize,
	}
	return NewS2cFishingBroadcastProtoMsg(msg)
}

var s2c_fishing_broadcast = [...]byte{16, 5} // 5
func NewS2cFishingBroadcastProtoMsg(object *S2CFishingBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fishing_broadcast[:], "s2c_fishing_broadcast")

}

func NewS2cUpdateFishPointMsg(point int32) pbutil.Buffer {
	msg := &S2CUpdateFishPointProto{
		Point: point,
	}
	return NewS2cUpdateFishPointProtoMsg(msg)
}

var s2c_update_fish_point = [...]byte{16, 7} // 7
func NewS2cUpdateFishPointProtoMsg(object *S2CUpdateFishPointProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_fish_point[:], "s2c_update_fish_point")

}

func NewS2cFishPointExchangeMsg(prize []byte, exchange_index int32, exist bool) pbutil.Buffer {
	msg := &S2CFishPointExchangeProto{
		Prize:         prize,
		ExchangeIndex: exchange_index,
		Exist:         exist,
	}
	return NewS2cFishPointExchangeProtoMsg(msg)
}

func NewS2cFishPointExchangeMarshalMsg(prize marshaler, exchange_index int32, exist bool) pbutil.Buffer {
	msg := &S2CFishPointExchangeProto{
		Prize:         safeMarshal(prize),
		ExchangeIndex: exchange_index,
		Exist:         exist,
	}
	return NewS2cFishPointExchangeProtoMsg(msg)
}

var s2c_fish_point_exchange = [...]byte{16, 9} // 9
func NewS2cFishPointExchangeProtoMsg(object *S2CFishPointExchangeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fish_point_exchange[:], "s2c_fish_point_exchange")

}

// 钓鱼积分兑换，积分不足
var ERR_FISH_POINT_EXCHANGE_FAIL_POINT_NOT_ENOUGH = pbutil.StaticBuffer{3, 16, 10, 1} // 10-1

// 钓鱼积分兑换，已拥有所有的兑换物
var ERR_FISH_POINT_EXCHANGE_FAIL_OWNER_ALL = pbutil.StaticBuffer{3, 16, 10, 2} // 10-2

func NewS2cSetFishingCaptainMsg(captain_id int32) pbutil.Buffer {
	msg := &S2CSetFishingCaptainProto{
		CaptainId: captain_id,
	}
	return NewS2cSetFishingCaptainProtoMsg(msg)
}

var s2c_set_fishing_captain = [...]byte{16, 12} // 12
func NewS2cSetFishingCaptainProtoMsg(object *S2CSetFishingCaptainProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_fishing_captain[:], "s2c_set_fishing_captain")

}

// vip等级不足
var ERR_SET_FISHING_CAPTAIN_FAIL_VIP_LEVEL_LIMIT = pbutil.StaticBuffer{3, 16, 13, 1} // 13-1

// 无效的配置id
var ERR_SET_FISHING_CAPTAIN_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 16, 13, 2} // 13-2

// 重复设置
var ERR_SET_FISHING_CAPTAIN_FAIL_DUPLICATE = pbutil.StaticBuffer{3, 16, 13, 3} // 13-3
