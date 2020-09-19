package red_packet

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
	MODULE_ID = 49

	C2S_BUY = 1

	C2S_CREATE = 4

	C2S_GRAB = 7
)

func NewS2cBuyMsg(data_id int32, new_count int32) pbutil.Buffer {
	msg := &S2CBuyProto{
		DataId:   data_id,
		NewCount: new_count,
	}
	return NewS2cBuyProtoMsg(msg)
}

var s2c_buy = [...]byte{49, 2} // 2
func NewS2cBuyProtoMsg(object *S2CBuyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_buy[:], "s2c_buy")

}

// 红包 data.id 错误
var ERR_BUY_FAIL_INVALID_DATA_ID = pbutil.StaticBuffer{3, 49, 3, 1} // 3-1

// 钱不够
var ERR_BUY_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 49, 3, 2} // 3-2

// 服务器错误
var ERR_BUY_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 49, 3, 3} // 3-3

// 元宝赠送额度不够
var ERR_BUY_FAIL_RECHANGE_YUANBAO_LIMIT = pbutil.StaticBuffer{3, 49, 3, 6} // 3-6

func NewS2cCreateMsg(data_id int32) pbutil.Buffer {
	msg := &S2CCreateProto{
		DataId: data_id,
	}
	return NewS2cCreateProtoMsg(msg)
}

var s2c_create = [...]byte{49, 5} // 5
func NewS2cCreateProtoMsg(object *S2CCreateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_create[:], "s2c_create")

}

// 红包 data.id 错误
var ERR_CREATE_FAIL_INVALID_DATA_ID = pbutil.StaticBuffer{3, 49, 6, 1} // 6-1

// 聊天类型错误
var ERR_CREATE_FAIL_INVALID_CHAT_TYPE = pbutil.StaticBuffer{3, 49, 6, 2} // 6-2

// 没买
var ERR_CREATE_FAIL_NOT_BOUGHT = pbutil.StaticBuffer{3, 49, 6, 4} // 6-4

// 数量错误
var ERR_CREATE_FAIL_COUNT_ERR = pbutil.StaticBuffer{3, 49, 6, 5} // 6-5

// 联盟人太少
var ERR_CREATE_FAIL_GUILD_LIMIT = pbutil.StaticBuffer{3, 49, 6, 6} // 6-6

// 发联盟红包，但没有联盟
var ERR_CREATE_FAIL_NO_GUILD = pbutil.StaticBuffer{3, 49, 6, 7} // 6-7

// 服务器错误
var ERR_CREATE_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 49, 6, 3} // 6-3

// 文字太长
var ERR_CREATE_FAIL_TEXT_TOO_LONG = pbutil.StaticBuffer{3, 49, 6, 8} // 6-8

// 充钱不够
var ERR_CREATE_FAIL_RECHARGE_NOT_ENOUGH = pbutil.StaticBuffer{3, 49, 6, 9} // 6-9

func NewS2cGrabMsg(red_packet *shared_proto.RedPacketProto, grab_money int32) pbutil.Buffer {
	msg := &S2CGrabProto{
		RedPacket: red_packet,
		GrabMoney: grab_money,
	}
	return NewS2cGrabProtoMsg(msg)
}

var s2c_grab = [...]byte{49, 8} // 8
func NewS2cGrabProtoMsg(object *S2CGrabProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_grab[:], "s2c_grab")

}

// 红包id不存在或已经删掉了
var ERR_GRAB_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 49, 9, 1} // 9-1

// 不是同联盟
var ERR_GRAB_FAIL_NOT_SAME_GUILD = pbutil.StaticBuffer{3, 49, 9, 3} // 9-3

// 过期了
var ERR_GRAB_FAIL_EXPIRED = pbutil.StaticBuffer{3, 49, 9, 4} // 9-4

// 服务器错误
var ERR_GRAB_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 49, 9, 2} // 9-2

func NewS2cAllGrabbedNoticeMsg(id []byte) pbutil.Buffer {
	msg := &S2CAllGrabbedNoticeProto{
		Id: id,
	}
	return NewS2cAllGrabbedNoticeProtoMsg(msg)
}

var s2c_all_grabbed_notice = [...]byte{49, 10} // 10
func NewS2cAllGrabbedNoticeProtoMsg(object *S2CAllGrabbedNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_all_grabbed_notice[:], "s2c_all_grabbed_notice")

}
