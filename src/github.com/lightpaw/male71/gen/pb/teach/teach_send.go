package teach

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
	MODULE_ID = 50

	C2S_FIGHT = 1

	C2S_COLLECT_PRIZE = 4
)

func NewS2cFightMsg(id int32) pbutil.Buffer {
	msg := &S2CFightProto{
		Id: id,
	}
	return NewS2cFightProtoMsg(msg)
}

var s2c_fight = [...]byte{50, 2} // 2
func NewS2cFightProtoMsg(object *S2CFightProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fight[:], "s2c_fight")

}

// id 不存在
var ERR_FIGHT_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 50, 3, 5} // 3-5

// 君主等级不够
var ERR_FIGHT_FAIL_HERO_LEVEL_LIMIT = pbutil.StaticBuffer{3, 50, 3, 1} // 3-1

// 霸业任务没通过
var ERR_FIGHT_FAIL_BA_YE_TASK_LIMIT = pbutil.StaticBuffer{3, 50, 3, 2} // 3-2

// 前置关卡没通过
var ERR_FIGHT_FAIL_PREV_CHAPTER_NOT_PASS = pbutil.StaticBuffer{3, 50, 3, 3} // 3-3

// 服务器错误
var ERR_FIGHT_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 50, 3, 4} // 3-4

// 推图没通过
var ERR_FIGHT_FAIL_DUNGEON_LIMIT = pbutil.StaticBuffer{3, 50, 3, 6} // 3-6

func NewS2cCollectPrizeMsg(id int32) pbutil.Buffer {
	msg := &S2CCollectPrizeProto{
		Id: id,
	}
	return NewS2cCollectPrizeProtoMsg(msg)
}

var s2c_collect_prize = [...]byte{50, 5} // 5
func NewS2cCollectPrizeProtoMsg(object *S2CCollectPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_prize[:], "s2c_collect_prize")

}

// id 不存在
var ERR_COLLECT_PRIZE_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 50, 6, 4} // 6-4

// 关卡没通过
var ERR_COLLECT_PRIZE_FAIL_NOT_PASS = pbutil.StaticBuffer{3, 50, 6, 1} // 6-1

// 奖励已领
var ERR_COLLECT_PRIZE_FAIL_ALREADY_COLLECTED = pbutil.StaticBuffer{3, 50, 6, 3} // 6-3
