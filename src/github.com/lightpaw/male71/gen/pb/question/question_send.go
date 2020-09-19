package question

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
	MODULE_ID = 34

	C2S_START = 1

	C2S_ANSWER = 4

	C2S_NEXT = 6

	C2S_GET_PRIZE = 9
)

var START_S2C = pbutil.StaticBuffer{2, 34, 2} // 2

// 问题 ID 错误
var ERR_START_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 34, 3, 2} // 3-2

// 次数用完了
var ERR_START_FAIL_NO_TIMES = pbutil.StaticBuffer{3, 34, 3, 1} // 3-1

// 正在答题中
var ERR_START_FAIL_IN_QUESTION = pbutil.StaticBuffer{3, 34, 3, 3} // 3-3

func NewS2cAnswerMsg(id int32) pbutil.Buffer {
	msg := &S2CAnswerProto{
		Id: id,
	}
	return NewS2cAnswerProtoMsg(msg)
}

var s2c_answer = [...]byte{34, 5} // 5
func NewS2cAnswerProtoMsg(object *S2CAnswerProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_answer[:], "s2c_answer")

}

// 答题没有开始
var ERR_ANSWER_FAIL_NOT_START = pbutil.StaticBuffer{3, 34, 12, 1} // 12-1

// 问题 ID 错误
var ERR_ANSWER_FAIL_INVALID_QID = pbutil.StaticBuffer{3, 34, 12, 2} // 12-2

// 当前最后一题已经答过了
var ERR_ANSWER_FAIL_ALREADY_ANSWERED = pbutil.StaticBuffer{3, 34, 12, 3} // 12-3

// 答案必须>0
var ERR_ANSWER_FAIL_INVALID_AID = pbutil.StaticBuffer{3, 34, 12, 4} // 12-4

var NEXT_S2C = pbutil.StaticBuffer{2, 34, 7} // 7

// 本轮题数答够了
var ERR_NEXT_FAIL_ENOUGH = pbutil.StaticBuffer{3, 34, 8, 1} // 8-1

// 问题 ID 不存在
var ERR_NEXT_FAIL_INVALID_QID = pbutil.StaticBuffer{3, 34, 8, 3} // 8-3

// 当前最后一题还没回答
var ERR_NEXT_FAIL_LAST_NOT_ANSWER = pbutil.StaticBuffer{3, 34, 8, 5} // 8-5

var GET_PRIZE_S2C = pbutil.StaticBuffer{2, 34, 10} // 10

// 还没答完
var ERR_GET_PRIZE_FAIL_NOT_FINISH = pbutil.StaticBuffer{3, 34, 11, 3} // 11-3

// 本轮分数不对
var ERR_GET_PRIZE_FAIL_SCORE_ERR = pbutil.StaticBuffer{3, 34, 11, 2} // 11-2

// 没有要领的奖励
var ERR_GET_PRIZE_FAIL_NO_PRIZE = pbutil.StaticBuffer{3, 34, 11, 4} // 11-4
