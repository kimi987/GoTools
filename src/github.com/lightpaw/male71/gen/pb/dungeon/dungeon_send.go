package dungeon

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
	MODULE_ID = 26

	C2S_CHALLENGE = 1

	C2S_COLLECT_CHAPTER_PRIZE = 4

	C2S_COLLECT_PASS_DUNGEON_PRIZE = 13

	C2S_AUTO_CHALLENGE = 7

	C2S_COLLECT_CHAPTER_STAR_PRIZE = 17
)

func NewS2cChallengeMsg(id int32, link string, share []byte, prize []byte, pass bool, is_first_pass bool, enabled_star []bool, pass_seconds int32, chapter_star int32, pass_times int32, is_refresh bool) pbutil.Buffer {
	msg := &S2CChallengeProto{
		Id:          id,
		Link:        link,
		Share:       share,
		Prize:       prize,
		Pass:        pass,
		IsFirstPass: is_first_pass,
		EnabledStar: enabled_star,
		PassSeconds: pass_seconds,
		ChapterStar: chapter_star,
		PassTimes:   pass_times,
		IsRefresh:   is_refresh,
	}
	return NewS2cChallengeProtoMsg(msg)
}

var s2c_challenge = [...]byte{26, 2} // 2
func NewS2cChallengeProtoMsg(object *S2CChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_challenge[:], "s2c_challenge")

}

// 副本没找到
var ERR_CHALLENGE_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 26, 3, 1} // 3-1

// 已经通关了
var ERR_CHALLENGE_FAIL_HAS_PASS = pbutil.StaticBuffer{3, 26, 3, 2} // 3-2

// 前置副本没有通关
var ERR_CHALLENGE_FAIL_PREV_NOT_PASS = pbutil.StaticBuffer{3, 26, 3, 3} // 3-3

// 前置霸业目标没有完成，请继续进行霸业任务
var ERR_CHALLENGE_FAIL_BAYE_NOT_PASS = pbutil.StaticBuffer{3, 26, 3, 11} // 3-11

// 君主等级不足
var ERR_CHALLENGE_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 26, 3, 4} // 3-4

// 上阵武将未满
var ERR_CHALLENGE_FAIL_CAPTAIN_NOT_FULL = pbutil.StaticBuffer{3, 26, 3, 5} // 3-5

// 上阵起码要有一个
var ERR_CHALLENGE_FAIL_NEED_GT_ONE = pbutil.StaticBuffer{3, 26, 3, 6} // 3-6

// 上阵武将超出上限
var ERR_CHALLENGE_FAIL_CAPTAIN_TOO_MUCH = pbutil.StaticBuffer{3, 26, 3, 7} // 3-7

// 上阵武将不存在
var ERR_CHALLENGE_FAIL_CAPTAIN_NOT_EXIST = pbutil.StaticBuffer{3, 26, 3, 8} // 3-8

// 上阵武将id重复
var ERR_CHALLENGE_FAIL_CAPTAIN_ID_DUPLICATE = pbutil.StaticBuffer{3, 26, 3, 9} // 3-9

// 服务器忙，请稍后再试
var ERR_CHALLENGE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 26, 3, 10} // 3-10

// 当日挑战次数上限
var ERR_CHALLENGE_FAIL_PASS_LIMIT = pbutil.StaticBuffer{3, 26, 3, 13} // 3-13

// 体力值不够
var ERR_CHALLENGE_FAIL_SP_NOT_ENOUGH = pbutil.StaticBuffer{3, 26, 3, 14} // 3-14

func NewS2cUpdateChallengeTimesMsg(start_time int32) pbutil.Buffer {
	msg := &S2CUpdateChallengeTimesProto{
		StartTime: start_time,
	}
	return NewS2cUpdateChallengeTimesProtoMsg(msg)
}

var s2c_update_challenge_times = [...]byte{26, 16} // 16
func NewS2cUpdateChallengeTimesProtoMsg(object *S2CUpdateChallengeTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_challenge_times[:], "s2c_update_challenge_times")

}

func NewS2cCollectChapterPrizeMsg(id int32) pbutil.Buffer {
	msg := &S2CCollectChapterPrizeProto{
		Id: id,
	}
	return NewS2cCollectChapterPrizeProtoMsg(msg)
}

var s2c_collect_chapter_prize = [...]byte{26, 5} // 5
func NewS2cCollectChapterPrizeProtoMsg(object *S2CCollectChapterPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_chapter_prize[:], "s2c_collect_chapter_prize")

}

// 通关奖励没找到
var ERR_COLLECT_CHAPTER_PRIZE_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 26, 6, 1} // 6-1

// 没通关
var ERR_COLLECT_CHAPTER_PRIZE_FAIL_NOT_PASS = pbutil.StaticBuffer{3, 26, 6, 2} // 6-2

// 奖励已经领取了
var ERR_COLLECT_CHAPTER_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 26, 6, 3} // 6-3

func NewS2cCollectPassDungeonPrizeMsg(id int32) pbutil.Buffer {
	msg := &S2CCollectPassDungeonPrizeProto{
		Id: id,
	}
	return NewS2cCollectPassDungeonPrizeProtoMsg(msg)
}

var s2c_collect_pass_dungeon_prize = [...]byte{26, 14} // 14
func NewS2cCollectPassDungeonPrizeProtoMsg(object *S2CCollectPassDungeonPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_pass_dungeon_prize[:], "s2c_collect_pass_dungeon_prize")

}

// 副本id无效
var ERR_COLLECT_PASS_DUNGEON_PRIZE_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 26, 15, 4} // 15-4

// 副本没有配置通关奖励
var ERR_COLLECT_PASS_DUNGEON_PRIZE_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 26, 15, 1} // 15-1

// 没通关
var ERR_COLLECT_PASS_DUNGEON_PRIZE_FAIL_NOT_PASS = pbutil.StaticBuffer{3, 26, 15, 2} // 15-2

// 奖励已经领取了
var ERR_COLLECT_PASS_DUNGEON_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 26, 15, 3} // 15-3

func NewS2cAutoChallengeMsg(id int32, prizes [][]byte, pass_times int32) pbutil.Buffer {
	msg := &S2CAutoChallengeProto{
		Id:        id,
		Prizes:    prizes,
		PassTimes: pass_times,
	}
	return NewS2cAutoChallengeProtoMsg(msg)
}

var s2c_auto_challenge = [...]byte{26, 8} // 8
func NewS2cAutoChallengeProtoMsg(object *S2CAutoChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_auto_challenge[:], "s2c_auto_challenge")

}

// 副本没找到
var ERR_AUTO_CHALLENGE_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 26, 9, 1} // 9-1

// 未通关
var ERR_AUTO_CHALLENGE_FAIL_NOT_PASS = pbutil.StaticBuffer{3, 26, 9, 2} // 9-2

// 扫荡次数非法
var ERR_AUTO_CHALLENGE_FAIL_INVALID_AUTO_TIMES = pbutil.StaticBuffer{3, 26, 9, 3} // 9-3

// 体力值不足
var ERR_AUTO_CHALLENGE_FAIL_SP_NOT_ENOUGH = pbutil.StaticBuffer{3, 26, 9, 7} // 9-7

// 无法扫荡未满星通关的副本
var ERR_AUTO_CHALLENGE_FAIL_NOT_FULL_STAR = pbutil.StaticBuffer{3, 26, 9, 5} // 9-5

// 扫荡次数超出每日通关上限
var ERR_AUTO_CHALLENGE_FAIL_PASS_LIMIT = pbutil.StaticBuffer{3, 26, 9, 6} // 9-6

func NewS2cCollectChapterStarPrizeMsg(id int32, collect_n int32, prize []byte) pbutil.Buffer {
	msg := &S2CCollectChapterStarPrizeProto{
		Id:       id,
		CollectN: collect_n,
		Prize:    prize,
	}
	return NewS2cCollectChapterStarPrizeProtoMsg(msg)
}

var s2c_collect_chapter_star_prize = [...]byte{26, 18} // 18
func NewS2cCollectChapterStarPrizeProtoMsg(object *S2CCollectChapterStarPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_chapter_star_prize[:], "s2c_collect_chapter_star_prize")

}

// 章节id无效
var ERR_COLLECT_CHAPTER_STAR_PRIZE_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 26, 19, 1} // 19-1

// 奖励下标无效
var ERR_COLLECT_CHAPTER_STAR_PRIZE_FAIL_INVALID_PRIZE_INDEX = pbutil.StaticBuffer{3, 26, 19, 2} // 19-2

// 星数不足
var ERR_COLLECT_CHAPTER_STAR_PRIZE_FAIL_STAR_NOT_ENOUGH = pbutil.StaticBuffer{3, 26, 19, 3} // 19-3

// 奖励已经领取了
var ERR_COLLECT_CHAPTER_STAR_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 26, 19, 4} // 19-4
