package tower

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
	MODULE_ID = 14

	C2S_CHALLENGE = 1

	C2S_AUTO_CHALLENGE = 5

	C2S_COLLECT_BOX = 8

	C2S_LIST_PASS_REPLAY = 11
)

func NewS2cChallengeMsg(link string, share []byte, first_pass_prize []byte, prize []byte, auto_max_floor int32) pbutil.Buffer {
	msg := &S2CChallengeProto{
		Link:           link,
		Share:          share,
		FirstPassPrize: first_pass_prize,
		Prize:          prize,
		AutoMaxFloor:   auto_max_floor,
	}
	return NewS2cChallengeProtoMsg(msg)
}

func NewS2cChallengeMarshalMsg(link string, share []byte, first_pass_prize marshaler, prize marshaler, auto_max_floor int32) pbutil.Buffer {
	msg := &S2CChallengeProto{
		Link:           link,
		Share:          share,
		FirstPassPrize: safeMarshal(first_pass_prize),
		Prize:          safeMarshal(prize),
		AutoMaxFloor:   auto_max_floor,
	}
	return NewS2cChallengeProtoMsg(msg)
}

var s2c_challenge = [...]byte{14, 2} // 2
func NewS2cChallengeProtoMsg(object *S2CChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_challenge[:], "s2c_challenge")

}

func NewS2cFailureChallengeMsg(challenge_times int32, next_reset_challenge_time int32, link string, share []byte) pbutil.Buffer {
	msg := &S2CFailureChallengeProto{
		ChallengeTimes:         challenge_times,
		NextResetChallengeTime: next_reset_challenge_time,
		Link:  link,
		Share: share,
	}
	return NewS2cFailureChallengeProtoMsg(msg)
}

var s2c_failure_challenge = [...]byte{14, 3} // 3
func NewS2cFailureChallengeProtoMsg(object *S2CFailureChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_failure_challenge[:], "s2c_failure_challenge")

}

// 当前楼层无效
var ERR_CHALLENGE_FAIL_INVALID_FLOOR = pbutil.StaticBuffer{3, 14, 4, 1} // 4-1

// 没有挑战次数了
var ERR_CHALLENGE_FAIL_MAX_CHALLENGE_TIMES = pbutil.StaticBuffer{3, 14, 4, 2} // 4-2

// 已经到达最高层
var ERR_CHALLENGE_FAIL_MAX_FLOOR = pbutil.StaticBuffer{3, 14, 4, 9} // 4-9

// 领取重楼宝箱
var ERR_CHALLENGE_FAIL_BOX = pbutil.StaticBuffer{3, 14, 4, 3} // 4-3

// 上阵武将未满
var ERR_CHALLENGE_FAIL_CAPTAIN_NOT_FULL = pbutil.StaticBuffer{3, 14, 4, 5} // 4-5

// 上阵武将超出上限
var ERR_CHALLENGE_FAIL_CAPTAIN_TOO_MUCH = pbutil.StaticBuffer{3, 14, 4, 8} // 4-8

// 上阵武将不存在
var ERR_CHALLENGE_FAIL_CAPTAIN_NOT_EXIST = pbutil.StaticBuffer{3, 14, 4, 6} // 4-6

// 上阵武将id重复
var ERR_CHALLENGE_FAIL_CAPTAIN_ID_DUPLICATE = pbutil.StaticBuffer{3, 14, 4, 7} // 4-7

// 服务器忙，请稍后再试
var ERR_CHALLENGE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 14, 4, 4} // 4-4

func NewS2cAutoChallengeMsg(floor int32, prize [][]byte) pbutil.Buffer {
	msg := &S2CAutoChallengeProto{
		Floor: floor,
		Prize: prize,
	}
	return NewS2cAutoChallengeProtoMsg(msg)
}

func NewS2cAutoChallengeMarshalMsg(floor int32, prize [][]byte) pbutil.Buffer {
	msg := &S2CAutoChallengeProto{
		Floor: floor,
		Prize: prize,
	}
	return NewS2cAutoChallengeProtoMsg(msg)
}

var s2c_auto_challenge = [...]byte{14, 6} // 6
func NewS2cAutoChallengeProtoMsg(object *S2CAutoChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_auto_challenge[:], "s2c_auto_challenge")

}

// 已经到达扫荡最高层
var ERR_AUTO_CHALLENGE_FAIL_AUTO_MAX = pbutil.StaticBuffer{3, 14, 7, 1} // 7-1

// 上阵武将未满
var ERR_AUTO_CHALLENGE_FAIL_CAPTAIN_NOT_FULL = pbutil.StaticBuffer{3, 14, 7, 3} // 7-3

// 上阵武将超出上限
var ERR_AUTO_CHALLENGE_FAIL_CAPTAIN_TOO_MUCH = pbutil.StaticBuffer{3, 14, 7, 6} // 7-6

// 上阵武将不存在
var ERR_AUTO_CHALLENGE_FAIL_CAPTAIN_NOT_EXIST = pbutil.StaticBuffer{3, 14, 7, 4} // 7-4

// 上阵武将id重复
var ERR_AUTO_CHALLENGE_FAIL_CAPTAIN_ID_DUPLICATE = pbutil.StaticBuffer{3, 14, 7, 5} // 7-5

// 服务器忙，请稍后再试
var ERR_AUTO_CHALLENGE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 14, 7, 2} // 7-2

func NewS2cCollectBoxMsg(box_floor int32) pbutil.Buffer {
	msg := &S2CCollectBoxProto{
		BoxFloor: box_floor,
	}
	return NewS2cCollectBoxProtoMsg(msg)
}

var s2c_collect_box = [...]byte{14, 9} // 9
func NewS2cCollectBoxProtoMsg(object *S2CCollectBoxProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_box[:], "s2c_collect_box")

}

// 当前楼层没有宝箱
var ERR_COLLECT_BOX_FAIL_NOT_BOX = pbutil.StaticBuffer{3, 14, 10, 1} // 10-1

// 重楼宝箱已经领取过了
var ERR_COLLECT_BOX_FAIL_COLLECTED = pbutil.StaticBuffer{3, 14, 10, 2} // 10-2

// 未通关该层，无法领取
var ERR_COLLECT_BOX_FAIL_CAN_NOT_COLLECT = pbutil.StaticBuffer{3, 14, 10, 4} // 10-4

// 服务器忙，请稍后再试
var ERR_COLLECT_BOX_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 14, 10, 3} // 10-3

func NewS2cListPassReplayMsg(floor int32, data []byte) pbutil.Buffer {
	msg := &S2CListPassReplayProto{
		Floor: floor,
		Data:  data,
	}
	return NewS2cListPassReplayProtoMsg(msg)
}

func NewS2cListPassReplayMarshalMsg(floor int32, data marshaler) pbutil.Buffer {
	msg := &S2CListPassReplayProto{
		Floor: floor,
		Data:  safeMarshal(data),
	}
	return NewS2cListPassReplayProtoMsg(msg)
}

var s2c_list_pass_replay = [...]byte{14, 12} // 12
func NewS2cListPassReplayProtoMsg(object *S2CListPassReplayProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_pass_replay[:], "s2c_list_pass_replay")

}

// 无效的楼层
var ERR_LIST_PASS_REPLAY_FAIL_INVALID_FLOOR = pbutil.StaticBuffer{3, 14, 13, 1} // 13-1

func NewS2cUpdateCurrentFloorMsg(floor int32) pbutil.Buffer {
	msg := &S2CUpdateCurrentFloorProto{
		Floor: floor,
	}
	return NewS2cUpdateCurrentFloorProtoMsg(msg)
}

var s2c_update_current_floor = [...]byte{14, 14} // 14
func NewS2cUpdateCurrentFloorProtoMsg(object *S2CUpdateCurrentFloorProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_current_floor[:], "s2c_update_current_floor")

}
