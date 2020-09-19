package bai_zhan

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
	MODULE_ID = 24

	C2S_QUERY_BAI_ZHAN_INFO = 1

	C2S_CLEAR_LAST_JUN_XIAN = 34

	C2S_BAI_ZHAN_CHALLENGE = 4

	C2S_COLLECT_SALARY = 7

	C2S_COLLECT_JUN_XIAN_PRIZE = 10

	C2S_SELF_RECORD = 29

	C2S_REQUEST_RANK = 23

	C2S_REQUEST_SELF_RANK = 26
)

func NewS2cQueryBaiZhanInfoMsg(data []byte) pbutil.Buffer {
	msg := &S2CQueryBaiZhanInfoProto{
		Data: data,
	}
	return NewS2cQueryBaiZhanInfoProtoMsg(msg)
}

func NewS2cQueryBaiZhanInfoMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CQueryBaiZhanInfoProto{
		Data: safeMarshal(data),
	}
	return NewS2cQueryBaiZhanInfoProtoMsg(msg)
}

var s2c_query_bai_zhan_info = [...]byte{24, 2} // 2
func NewS2cQueryBaiZhanInfoProtoMsg(object *S2CQueryBaiZhanInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_query_bai_zhan_info[:], "s2c_query_bai_zhan_info")

}

// 服务器忙，请稍后再试
var ERR_QUERY_BAI_ZHAN_INFO_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 24, 3, 1} // 3-1

var CLEAR_LAST_JUN_XIAN_S2C = pbutil.StaticBuffer{2, 24, 35} // 35

func NewS2cBaiZhanChallengeMsg(win bool, challenge_times int32, link string, share []byte, point int32, history_max_point int32) pbutil.Buffer {
	msg := &S2CBaiZhanChallengeProto{
		Win:             win,
		ChallengeTimes:  challenge_times,
		Link:            link,
		Share:           share,
		Point:           point,
		HistoryMaxPoint: history_max_point,
	}
	return NewS2cBaiZhanChallengeProtoMsg(msg)
}

var s2c_bai_zhan_challenge = [...]byte{24, 5} // 5
func NewS2cBaiZhanChallengeProtoMsg(object *S2CBaiZhanChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_bai_zhan_challenge[:], "s2c_bai_zhan_challenge")

}

// 上阵武将未满
var ERR_BAI_ZHAN_CHALLENGE_FAIL_CAPTAIN_NOT_FULL = pbutil.StaticBuffer{3, 24, 6, 1} // 6-1

// 上阵武将超出上限
var ERR_BAI_ZHAN_CHALLENGE_FAIL_CAPTAIN_TOO_MUCH = pbutil.StaticBuffer{3, 24, 6, 2} // 6-2

// 上阵武将不存在
var ERR_BAI_ZHAN_CHALLENGE_FAIL_CAPTAIN_NOT_EXIST = pbutil.StaticBuffer{3, 24, 6, 3} // 6-3

// 上阵武将id重复
var ERR_BAI_ZHAN_CHALLENGE_FAIL_CAPTAIN_ID_DUPLICATE = pbutil.StaticBuffer{3, 24, 6, 4} // 6-4

// 服务器忙，请稍后再试
var ERR_BAI_ZHAN_CHALLENGE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 24, 6, 5} // 6-5

// 没有挑战次数了
var ERR_BAI_ZHAN_CHALLENGE_FAIL_NO_CHALLENGE_TIEMS = pbutil.StaticBuffer{3, 24, 6, 6} // 6-6

var COLLECT_SALARY_S2C = pbutil.StaticBuffer{2, 24, 8} // 8

// 没有俸禄
var ERR_COLLECT_SALARY_FAIL_NO_SALARY = pbutil.StaticBuffer{3, 24, 9, 1} // 9-1

// 俸禄已经领取了
var ERR_COLLECT_SALARY_FAIL_SALARY_COLLECT = pbutil.StaticBuffer{3, 24, 9, 2} // 9-2

// 服务器忙，请稍后再试
var ERR_COLLECT_SALARY_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 24, 9, 3} // 9-3

func NewS2cCollectJunXianPrizeMsg(id int32) pbutil.Buffer {
	msg := &S2CCollectJunXianPrizeProto{
		Id: id,
	}
	return NewS2cCollectJunXianPrizeProtoMsg(msg)
}

var s2c_collect_jun_xian_prize = [...]byte{24, 11} // 11
func NewS2cCollectJunXianPrizeProtoMsg(object *S2CCollectJunXianPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_jun_xian_prize[:], "s2c_collect_jun_xian_prize")

}

// 奖励已经领取了
var ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_PRIZE_COLLECTED = pbutil.StaticBuffer{3, 24, 12, 1} // 12-1

// 奖励没找到
var ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_PRIZE_NOT_FOUND = pbutil.StaticBuffer{3, 24, 12, 2} // 12-2

// 军衔等级太低
var ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_JUN_XIAN_LEVEL_TOO_LOW = pbutil.StaticBuffer{3, 24, 12, 3} // 12-3

// 积分太低
var ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_POINT_TOO_LOW = pbutil.StaticBuffer{3, 24, 12, 4} // 12-4

// 无效的奖励id
var ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_INVALID_PRIZE = pbutil.StaticBuffer{3, 24, 12, 5} // 12-5

// 服务器忙，请稍后再试
var ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 24, 12, 6} // 12-6

var RESET_S2C = pbutil.StaticBuffer{2, 24, 13} // 13

func NewS2cSelfRecordMsg(version int32, data [][]byte) pbutil.Buffer {
	msg := &S2CSelfRecordProto{
		Version: version,
		Data:    data,
	}
	return NewS2cSelfRecordProtoMsg(msg)
}

func NewS2cSelfRecordMarshalMsg(version int32, data [][]byte) pbutil.Buffer {
	msg := &S2CSelfRecordProto{
		Version: version,
		Data:    data,
	}
	return NewS2cSelfRecordProtoMsg(msg)
}

var s2c_self_record = [...]byte{24, 30} // 30
func NewS2cSelfRecordProtoMsg(object *S2CSelfRecordProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_self_record[:], "s2c_self_record")

}

var SELF_RECORD_S2C_NO_CHANGE = pbutil.StaticBuffer{2, 24, 31} // 31

// 服务器繁忙，请稍后再试
var ERR_SELF_RECORD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 24, 32, 1} // 32-1

var SELF_DEFENCE_RECORD_CHANGED_S2C = pbutil.StaticBuffer{2, 24, 22} // 22

func NewS2cRequestRankMsg(self bool, jun_xian_level int32, start_rank int32, total_rank_count int32, level_up_need_min_point int32, level_keep_need_min_point int32, data [][]byte) pbutil.Buffer {
	msg := &S2CRequestRankProto{
		Self:                  self,
		JunXianLevel:          jun_xian_level,
		StartRank:             start_rank,
		TotalRankCount:        total_rank_count,
		LevelUpNeedMinPoint:   level_up_need_min_point,
		LevelKeepNeedMinPoint: level_keep_need_min_point,
		Data: data,
	}
	return NewS2cRequestRankProtoMsg(msg)
}

func NewS2cRequestRankMarshalMsg(self bool, jun_xian_level int32, start_rank int32, total_rank_count int32, level_up_need_min_point int32, level_keep_need_min_point int32, data [][]byte) pbutil.Buffer {
	msg := &S2CRequestRankProto{
		Self:                  self,
		JunXianLevel:          jun_xian_level,
		StartRank:             start_rank,
		TotalRankCount:        total_rank_count,
		LevelUpNeedMinPoint:   level_up_need_min_point,
		LevelKeepNeedMinPoint: level_keep_need_min_point,
		Data: data,
	}
	return NewS2cRequestRankProtoMsg(msg)
}

var s2c_request_rank = [...]byte{24, 24} // 24
func NewS2cRequestRankProtoMsg(object *S2CRequestRankProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_rank[:], "s2c_request_rank")

}

// 服务器繁忙，请稍后再试
var ERR_REQUEST_RANK_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 24, 28, 1} // 28-1

func NewS2cRequestSelfRankMsg(rank int32, level_change_type int32, level_up_need_min_point int32, level_keep_need_min_point int32) pbutil.Buffer {
	msg := &S2CRequestSelfRankProto{
		Rank:                  rank,
		LevelChangeType:       level_change_type,
		LevelUpNeedMinPoint:   level_up_need_min_point,
		LevelKeepNeedMinPoint: level_keep_need_min_point,
	}
	return NewS2cRequestSelfRankProtoMsg(msg)
}

var s2c_request_self_rank = [...]byte{24, 27} // 27
func NewS2cRequestSelfRankProtoMsg(object *S2CRequestSelfRankProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_self_rank[:], "s2c_request_self_rank")

}

func NewS2cMaxJunXianLevelChangedMsg(level int32) pbutil.Buffer {
	msg := &S2CMaxJunXianLevelChangedProto{
		Level: level,
	}
	return NewS2cMaxJunXianLevelChangedProtoMsg(msg)
}

var s2c_max_jun_xian_level_changed = [...]byte{24, 33} // 33
func NewS2cMaxJunXianLevelChangedProtoMsg(object *S2CMaxJunXianLevelChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_max_jun_xian_level_changed[:], "s2c_max_jun_xian_level_changed")

}
