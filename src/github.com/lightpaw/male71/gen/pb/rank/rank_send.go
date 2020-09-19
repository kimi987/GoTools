package rank

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
	MODULE_ID = 23

	C2S_REQUEST_RANK = 1
)

func NewS2cRequestRankMsg(rank []byte) pbutil.Buffer {
	msg := &S2CRequestRankProto{
		Rank: rank,
	}
	return NewS2cRequestRankProtoMsg(msg)
}

var s2c_request_rank = [...]byte{23, 2} // 2
func NewS2cRequestRankProtoMsg(object *S2CRequestRankProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_request_rank[:], "s2c_request_rank")

}

// 未知的排行榜类型
var ERR_REQUEST_RANK_FAIL_UNKNOWN_RANK_TYPE = pbutil.StaticBuffer{3, 23, 3, 1} // 3-1

// 目标不存在
var ERR_REQUEST_RANK_FAIL_TARGET_NOT_FOUND = pbutil.StaticBuffer{3, 23, 3, 2} // 3-2

// 目标不在榜单上
var ERR_REQUEST_RANK_FAIL_TARGET_NOT_IN_RANK_LIST = pbutil.StaticBuffer{3, 23, 3, 3} // 3-3

// 服务器繁忙，请稍后再试
var ERR_REQUEST_RANK_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 23, 3, 4} // 3-4
