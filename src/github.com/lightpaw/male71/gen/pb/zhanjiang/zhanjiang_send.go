package zhanjiang

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
	MODULE_ID = 33

	C2S_OPEN = 1

	C2S_GIVE_UP = 4

	C2S_UPDATE_CAPTAIN = 7

	C2S_CHALLENGE = 10
)

func NewS2cOpenMsg(id int32, captain_id int32) pbutil.Buffer {
	msg := &S2COpenProto{
		Id:        id,
		CaptainId: captain_id,
	}
	return NewS2cOpenProtoMsg(msg)
}

var s2c_open = [...]byte{33, 2} // 2
func NewS2cOpenProtoMsg(object *S2COpenProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_open[:], "s2c_open")

}

// 关卡没找到
var ERR_OPEN_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 33, 3, 1} // 3-1

// 前置没有通关
var ERR_OPEN_FAIL_PREV_NOT_PASS = pbutil.StaticBuffer{3, 33, 3, 2} // 3-2

// 功能未开启
var ERR_OPEN_FAIL_FUNC_NOT_OPEN = pbutil.StaticBuffer{3, 33, 3, 3} // 3-3

// 关卡已经开启
var ERR_OPEN_FAIL_IS_OPEN = pbutil.StaticBuffer{3, 33, 3, 4} // 3-4

// 没有设置上场武将
var ERR_OPEN_FAIL_NO_CAPTAIN = pbutil.StaticBuffer{3, 33, 3, 5} // 3-5

// 开启次数不足
var ERR_OPEN_FAIL_NO_OPEN_TIMES = pbutil.StaticBuffer{3, 33, 3, 6} // 3-6

// 服务器忙，请稍后再试
var ERR_OPEN_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 33, 3, 7} // 3-7

var GIVE_UP_S2C = pbutil.StaticBuffer{2, 33, 5} // 5

// 未开启
var ERR_GIVE_UP_FAIL_NOT_OPEN = pbutil.StaticBuffer{3, 33, 6, 1} // 6-1

// 服务器忙，请稍后再试
var ERR_GIVE_UP_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 33, 6, 2} // 6-2

func NewS2cUpdateCaptainMsg(id int32) pbutil.Buffer {
	msg := &S2CUpdateCaptainProto{
		Id: id,
	}
	return NewS2cUpdateCaptainProtoMsg(msg)
}

var s2c_update_captain = [...]byte{33, 8} // 8
func NewS2cUpdateCaptainProtoMsg(object *S2CUpdateCaptainProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_captain[:], "s2c_update_captain")

}

// 武将没找到
var ERR_UPDATE_CAPTAIN_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 33, 9, 1} // 9-1

// 没有开启
var ERR_UPDATE_CAPTAIN_FAIL_NOT_OPEN = pbutil.StaticBuffer{3, 33, 9, 2} // 9-2

// 服务器忙，请稍后再试
var ERR_UPDATE_CAPTAIN_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 33, 9, 3} // 9-3

func NewS2cChallengeMsg(pass bool, link string, share []byte, zhan_jiang_data_id int32) pbutil.Buffer {
	msg := &S2CChallengeProto{
		Pass:            pass,
		Link:            link,
		Share:           share,
		ZhanJiangDataId: zhan_jiang_data_id,
	}
	return NewS2cChallengeProtoMsg(msg)
}

var s2c_challenge = [...]byte{33, 11} // 11
func NewS2cChallengeProtoMsg(object *S2CChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_challenge[:], "s2c_challenge")

}

// 挑战没开启
var ERR_CHALLENGE_FAIL_NOT_OPEN = pbutil.StaticBuffer{3, 33, 12, 1} // 12-1

// 没有设置武将
var ERR_CHALLENGE_FAIL_NO_CAPTAIN = pbutil.StaticBuffer{3, 33, 12, 2} // 12-2

// 请勿重复发送挑战请求
var ERR_CHALLENGE_FAIL_NO_DUPLICATE = pbutil.StaticBuffer{3, 33, 12, 3} // 12-3

// 未设置出战武将
var ERR_CHALLENGE_FAIL_NOT_SET_CAPTAIN = pbutil.StaticBuffer{3, 33, 12, 4} // 12-4

// 服务器忙，请稍后再试
var ERR_CHALLENGE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 33, 12, 5} // 12-5

func NewS2cPassMsg(id int32) pbutil.Buffer {
	msg := &S2CPassProto{
		Id: id,
	}
	return NewS2cPassProtoMsg(msg)
}

var s2c_pass = [...]byte{33, 13} // 13
func NewS2cPassProtoMsg(object *S2CPassProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_pass[:], "s2c_pass")

}

func NewS2cUpdateOpenTimesMsg(open_times int32) pbutil.Buffer {
	msg := &S2CUpdateOpenTimesProto{
		OpenTimes: open_times,
	}
	return NewS2cUpdateOpenTimesProtoMsg(msg)
}

var s2c_update_open_times = [...]byte{33, 14} // 14
func NewS2cUpdateOpenTimesProtoMsg(object *S2CUpdateOpenTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_open_times[:], "s2c_update_open_times")

}
