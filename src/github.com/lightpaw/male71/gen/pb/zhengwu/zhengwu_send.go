package zhengwu

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
	MODULE_ID = 32

	C2S_START = 1

	C2S_COLLECT = 4

	C2S_YUANBAO_COMPLETE = 7

	C2S_YUANBAO_REFRESH = 10

	C2S_VIP_COLLECT = 14
)

func NewS2cStartMsg(zheng_wu []byte) pbutil.Buffer {
	msg := &S2CStartProto{
		ZhengWu: zheng_wu,
	}
	return NewS2cStartProtoMsg(msg)
}

var s2c_start = [...]byte{32, 2} // 2
func NewS2cStartProtoMsg(object *S2CStartProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_start[:], "s2c_start")

}

// 政务没找到
var ERR_START_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 32, 3, 1} // 3-1

// 当前有政务已经完成，请先领取奖励
var ERR_START_FAIL_NEED_COLLECT = pbutil.StaticBuffer{3, 32, 3, 2} // 3-2

// 当前有政务正在做
var ERR_START_FAIL_HAVE_DOING = pbutil.StaticBuffer{3, 32, 3, 3} // 3-3

// 服务器繁忙，请稍后再试
var ERR_START_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 32, 3, 4} // 3-4

var COLLECT_S2C = pbutil.StaticBuffer{2, 32, 5} // 5

// 没有正在做的政务
var ERR_COLLECT_FAIL_NO_DOING_ZHENG_WU = pbutil.StaticBuffer{3, 32, 6, 1} // 6-1

// 政务还没完成
var ERR_COLLECT_FAIL_NOT_COMPLETE = pbutil.StaticBuffer{3, 32, 6, 2} // 6-2

// 服务器繁忙，请稍后再试
var ERR_COLLECT_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 32, 6, 3} // 6-3

var YUANBAO_COMPLETE_S2C = pbutil.StaticBuffer{2, 32, 8} // 8

// 当前没有正在做的政务
var ERR_YUANBAO_COMPLETE_FAIL_NOT_DOING = pbutil.StaticBuffer{3, 32, 9, 1} // 9-1

// 已经完成了
var ERR_YUANBAO_COMPLETE_FAIL_COMPLETE = pbutil.StaticBuffer{3, 32, 9, 2} // 9-2

// 元宝不足
var ERR_YUANBAO_COMPLETE_FAIL_NOT_ENOUGH_YUANBAO = pbutil.StaticBuffer{3, 32, 9, 3} // 9-3

// 服务器繁忙，请稍后再试
var ERR_YUANBAO_COMPLETE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 32, 9, 4} // 9-4

var YUANBAO_REFRESH_S2C = pbutil.StaticBuffer{2, 32, 11} // 11

// 当前正在做的政务已经完成，请先领取奖励再刷新
var ERR_YUANBAO_REFRESH_FAIL_COMPLETE = pbutil.StaticBuffer{3, 32, 12, 1} // 12-1

// 点券不足
var ERR_YUANBAO_REFRESH_FAIL_NOT_ENOUGH_COST = pbutil.StaticBuffer{3, 32, 12, 5} // 12-5

// 次数不足
var ERR_YUANBAO_REFRESH_FAIL_NOT_ENOUGH_TIMES = pbutil.StaticBuffer{3, 32, 12, 4} // 12-4

// 服务器繁忙，请稍后再试
var ERR_YUANBAO_REFRESH_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 32, 12, 3} // 12-3

func NewS2cRefreshMsg(new_zheng_wu []byte) pbutil.Buffer {
	msg := &S2CRefreshProto{
		NewZhengWu: new_zheng_wu,
	}
	return NewS2cRefreshProtoMsg(msg)
}

var s2c_refresh = [...]byte{32, 13} // 13
func NewS2cRefreshProtoMsg(object *S2CRefreshProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_refresh[:], "s2c_refresh")

}

func NewS2cVipCollectMsg(id int32) pbutil.Buffer {
	msg := &S2CVipCollectProto{
		Id: id,
	}
	return NewS2cVipCollectProtoMsg(msg)
}

var s2c_vip_collect = [...]byte{32, 15} // 15
func NewS2cVipCollectProtoMsg(object *S2CVipCollectProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_vip_collect[:], "s2c_vip_collect")

}

// 没有vip特权
var ERR_VIP_COLLECT_FAIL_VIP_LIMIT = pbutil.StaticBuffer{3, 32, 16, 1} // 16-1

// 没在政务列表中
var ERR_VIP_COLLECT_FAIL_NOT_IN_LIST = pbutil.StaticBuffer{3, 32, 16, 2} // 16-2

// 服务器错误
var ERR_VIP_COLLECT_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 32, 16, 3} // 16-3
