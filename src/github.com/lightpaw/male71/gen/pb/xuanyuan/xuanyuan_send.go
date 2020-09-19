package xuanyuan

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
	MODULE_ID = 40

	C2S_SELF_INFO = 1

	C2S_LIST_TARGET = 11

	C2S_QUERY_TARGET_TROOP = 5

	C2S_CHALLENGE = 15

	C2S_LIST_RECORD = 20

	C2S_COLLECT_RANK_PRIZE = 22
)

func NewS2cRankIsEmptyMsg(is_empty bool) pbutil.Buffer {
	msg := &S2CRankIsEmptyProto{
		IsEmpty: is_empty,
	}
	return NewS2cRankIsEmptyProtoMsg(msg)
}

var s2c_rank_is_empty = [...]byte{40, 26} // 26
func NewS2cRankIsEmptyProtoMsg(object *S2CRankIsEmptyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_rank_is_empty[:], "s2c_rank_is_empty")

}

func NewS2cSelfInfoMsg(rank int32, score int32, win int32, lose int32, range_id int32, first_target_rank int32, targets [][]byte) pbutil.Buffer {
	msg := &S2CSelfInfoProto{
		Rank:            rank,
		Score:           score,
		Win:             win,
		Lose:            lose,
		RangeId:         range_id,
		FirstTargetRank: first_target_rank,
		Targets:         targets,
	}
	return NewS2cSelfInfoProtoMsg(msg)
}

func NewS2cSelfInfoMarshalMsg(rank int32, score int32, win int32, lose int32, range_id int32, first_target_rank int32, targets [][]byte) pbutil.Buffer {
	msg := &S2CSelfInfoProto{
		Rank:            rank,
		Score:           score,
		Win:             win,
		Lose:            lose,
		RangeId:         range_id,
		FirstTargetRank: first_target_rank,
		Targets:         targets,
	}
	return NewS2cSelfInfoProtoMsg(msg)
}

var s2c_self_info = [...]byte{40, 2} // 2
func NewS2cSelfInfoProtoMsg(object *S2CSelfInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_self_info[:], "s2c_self_info")

}

func NewS2cListTargetMsg(range_id int32, first_target_rank int32, targets [][]byte) pbutil.Buffer {
	msg := &S2CListTargetProto{
		RangeId:         range_id,
		FirstTargetRank: first_target_rank,
		Targets:         targets,
	}
	return NewS2cListTargetProtoMsg(msg)
}

func NewS2cListTargetMarshalMsg(range_id int32, first_target_rank int32, targets [][]byte) pbutil.Buffer {
	msg := &S2CListTargetProto{
		RangeId:         range_id,
		FirstTargetRank: first_target_rank,
		Targets:         targets,
	}
	return NewS2cListTargetProtoMsg(msg)
}

var s2c_list_target = [...]byte{40, 12} // 12
func NewS2cListTargetProtoMsg(object *S2CListTargetProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_target[:], "s2c_list_target")

}

// 无效的分组id
var ERR_LIST_TARGET_FAIL_INVALID_RANGE = pbutil.StaticBuffer{3, 40, 13, 1} // 13-1

func NewS2cQueryTargetTroopMsg(id []byte, version int32, player []byte) pbutil.Buffer {
	msg := &S2CQueryTargetTroopProto{
		Id:      id,
		Version: version,
		Player:  player,
	}
	return NewS2cQueryTargetTroopProtoMsg(msg)
}

func NewS2cQueryTargetTroopMarshalMsg(id []byte, version int32, player marshaler) pbutil.Buffer {
	msg := &S2CQueryTargetTroopProto{
		Id:      id,
		Version: version,
		Player:  safeMarshal(player),
	}
	return NewS2cQueryTargetTroopProtoMsg(msg)
}

var s2c_query_target_troop = [...]byte{40, 6} // 6
func NewS2cQueryTargetTroopProtoMsg(object *S2CQueryTargetTroopProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_query_target_troop[:], "s2c_query_target_troop")

}

// 无效的id
var ERR_QUERY_TARGET_TROOP_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 40, 14, 1} // 14-1

func NewS2cChallengeMsg(id []byte, link string, add_score int32) pbutil.Buffer {
	msg := &S2CChallengeProto{
		Id:       id,
		Link:     link,
		AddScore: add_score,
	}
	return NewS2cChallengeProtoMsg(msg)
}

var s2c_challenge = [...]byte{40, 16} // 16
func NewS2cChallengeProtoMsg(object *S2CChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_challenge[:], "s2c_challenge")

}

// 无效的id
var ERR_CHALLENGE_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 40, 17, 1} // 17-1

// 目标已经挑战过了
var ERR_CHALLENGE_FAIL_CHALLENGED = pbutil.StaticBuffer{3, 40, 17, 2} // 17-2

// 挑战次数不足
var ERR_CHALLENGE_FAIL_TIMES_LIMIT = pbutil.StaticBuffer{3, 40, 17, 3} // 17-3

// 挑战目标阵容已改变
var ERR_CHALLENGE_FAIL_VERSION = pbutil.StaticBuffer{3, 40, 17, 4} // 17-4

// 上阵武将个数不足
var ERR_CHALLENGE_FAIL_CAPTAIN_NOT_ENOUGH = pbutil.StaticBuffer{3, 40, 17, 5} // 17-5

// 服务器忙，请稍后重试
var ERR_CHALLENGE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 40, 17, 6} // 17-6

func NewS2cUpdateXyInfoMsg(score int32, win int32, lose int32) pbutil.Buffer {
	msg := &S2CUpdateXyInfoProto{
		Score: score,
		Win:   win,
		Lose:  lose,
	}
	return NewS2cUpdateXyInfoProtoMsg(msg)
}

var s2c_update_xy_info = [...]byte{40, 18} // 18
func NewS2cUpdateXyInfoProtoMsg(object *S2CUpdateXyInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_xy_info[:], "s2c_update_xy_info")

}

func NewS2cAddRecordMsg(id int32, data []byte) pbutil.Buffer {
	msg := &S2CAddRecordProto{
		Id:   id,
		Data: data,
	}
	return NewS2cAddRecordProtoMsg(msg)
}

func NewS2cAddRecordMarshalMsg(id int32, data marshaler) pbutil.Buffer {
	msg := &S2CAddRecordProto{
		Id:   id,
		Data: safeMarshal(data),
	}
	return NewS2cAddRecordProtoMsg(msg)
}

var s2c_add_record = [...]byte{40, 19} // 19
func NewS2cAddRecordProtoMsg(object *S2CAddRecordProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_record[:], "s2c_add_record")

}

func NewS2cListRecordMsg(id int32, up bool, ids []int32, data [][]byte) pbutil.Buffer {
	msg := &S2CListRecordProto{
		Id:   id,
		Up:   up,
		Ids:  ids,
		Data: data,
	}
	return NewS2cListRecordProtoMsg(msg)
}

func NewS2cListRecordMarshalMsg(id int32, up bool, ids []int32, data [][]byte) pbutil.Buffer {
	msg := &S2CListRecordProto{
		Id:   id,
		Up:   up,
		Ids:  ids,
		Data: data,
	}
	return NewS2cListRecordProtoMsg(msg)
}

var s2c_list_record = [...]byte{40, 21} // 21
func NewS2cListRecordProtoMsg(object *S2CListRecordProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_list_record[:], "s2c_list_record")

}

func NewS2cCollectRankPrizeMsg(prize *shared_proto.PrizeProto) pbutil.Buffer {
	msg := &S2CCollectRankPrizeProto{
		Prize: prize,
	}
	return NewS2cCollectRankPrizeProtoMsg(msg)
}

var s2c_collect_rank_prize = [...]byte{40, 23} // 23
func NewS2cCollectRankPrizeProtoMsg(object *S2CCollectRankPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_rank_prize[:], "s2c_collect_rank_prize")

}

// 已领取过奖励
var ERR_COLLECT_RANK_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 40, 24, 1} // 24-1

var RESET_S2C = pbutil.StaticBuffer{2, 40, 25} // 25
