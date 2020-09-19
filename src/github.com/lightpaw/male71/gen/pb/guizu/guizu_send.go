package guizu

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
	MODULE_ID = 30

	C2S_COLLECT_LEVEL_PRIZE = 1
)

func NewS2cCollectLevelPrizeMsg(level int32) pbutil.Buffer {
	msg := &S2CCollectLevelPrizeProto{
		Level: level,
	}
	return NewS2cCollectLevelPrizeProtoMsg(msg)
}

var s2c_collect_level_prize = [...]byte{30, 2} // 2
func NewS2cCollectLevelPrizeProtoMsg(object *S2CCollectLevelPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_level_prize[:], "s2c_collect_level_prize")

}

// 已经领取了
var ERR_COLLECT_LEVEL_PRIZE_FAIL_ALREADY_COLLECT = pbutil.StaticBuffer{3, 30, 3, 1} // 3-1

// 等级未达成
var ERR_COLLECT_LEVEL_PRIZE_FAIL_LEVEL_UNREACH = pbutil.StaticBuffer{3, 30, 3, 2} // 3-2

// 未知的等级
var ERR_COLLECT_LEVEL_PRIZE_FAIL_INVALID_LEVEL = pbutil.StaticBuffer{3, 30, 3, 3} // 3-3

// 服务器忙，请稍后再试
var ERR_COLLECT_LEVEL_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 30, 3, 4} // 3-4

func NewS2cLevelChangedMsg(level int32) pbutil.Buffer {
	msg := &S2CLevelChangedProto{
		Level: level,
	}
	return NewS2cLevelChangedProtoMsg(msg)
}

var s2c_level_changed = [...]byte{30, 4} // 4
func NewS2cLevelChangedProtoMsg(object *S2CLevelChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_level_changed[:], "s2c_level_changed")

}
