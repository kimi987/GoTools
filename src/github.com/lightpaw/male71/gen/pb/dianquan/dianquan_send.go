package dianquan

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
	MODULE_ID = 39

	C2S_EXCHANGE = 1
)

func NewS2cExchangeMsg(times int32, amount int32) pbutil.Buffer {
	msg := &S2CExchangeProto{
		Times:  times,
		Amount: amount,
	}
	return NewS2cExchangeProtoMsg(msg)
}

var s2c_exchange = [...]byte{39, 2} // 2
func NewS2cExchangeProtoMsg(object *S2CExchangeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_exchange[:], "s2c_exchange")

}

// 兑换次数错误
var ERR_EXCHANGE_FAIL_INVALID_TIMES = pbutil.StaticBuffer{3, 39, 3, 1} // 3-1

// 消耗不足
var ERR_EXCHANGE_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 39, 3, 2} // 3-2

// 服务器错误
var ERR_EXCHANGE_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 39, 3, 3} // 3-3
