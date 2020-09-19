package bai_zhan

import (
	"github.com/lightpaw/pbutil"
)

// query_bai_zhan_info
var (
	ErrQueryBaiZhanInfoFailServerError = newMsgError("query_bai_zhan_info 服务器忙，请稍后再试", ERR_QUERY_BAI_ZHAN_INFO_FAIL_SERVER_ERROR) // 3-1
)

// bai_zhan_challenge
var (
	ErrBaiZhanChallengeFailCaptainNotFull     = newMsgError("bai_zhan_challenge 上阵武将未满", ERR_BAI_ZHAN_CHALLENGE_FAIL_CAPTAIN_NOT_FULL)       // 6-1
	ErrBaiZhanChallengeFailCaptainTooMuch     = newMsgError("bai_zhan_challenge 上阵武将超出上限", ERR_BAI_ZHAN_CHALLENGE_FAIL_CAPTAIN_TOO_MUCH)     // 6-2
	ErrBaiZhanChallengeFailCaptainNotExist    = newMsgError("bai_zhan_challenge 上阵武将不存在", ERR_BAI_ZHAN_CHALLENGE_FAIL_CAPTAIN_NOT_EXIST)     // 6-3
	ErrBaiZhanChallengeFailCaptainIdDuplicate = newMsgError("bai_zhan_challenge 上阵武将id重复", ERR_BAI_ZHAN_CHALLENGE_FAIL_CAPTAIN_ID_DUPLICATE) // 6-4
	ErrBaiZhanChallengeFailServerError        = newMsgError("bai_zhan_challenge 服务器忙，请稍后再试", ERR_BAI_ZHAN_CHALLENGE_FAIL_SERVER_ERROR)       // 6-5
	ErrBaiZhanChallengeFailNoChallengeTiems   = newMsgError("bai_zhan_challenge 没有挑战次数了", ERR_BAI_ZHAN_CHALLENGE_FAIL_NO_CHALLENGE_TIEMS)    // 6-6
)

// collect_salary
var (
	ErrCollectSalaryFailNoSalary      = newMsgError("collect_salary 没有俸禄", ERR_COLLECT_SALARY_FAIL_NO_SALARY)          // 9-1
	ErrCollectSalaryFailSalaryCollect = newMsgError("collect_salary 俸禄已经领取了", ERR_COLLECT_SALARY_FAIL_SALARY_COLLECT)  // 9-2
	ErrCollectSalaryFailServerError   = newMsgError("collect_salary 服务器忙，请稍后再试", ERR_COLLECT_SALARY_FAIL_SERVER_ERROR) // 9-3
)

// collect_jun_xian_prize
var (
	ErrCollectJunXianPrizeFailPrizeCollected     = newMsgError("collect_jun_xian_prize 奖励已经领取了", ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_PRIZE_COLLECTED)       // 12-1
	ErrCollectJunXianPrizeFailPrizeNotFound      = newMsgError("collect_jun_xian_prize 奖励没找到", ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_PRIZE_NOT_FOUND)         // 12-2
	ErrCollectJunXianPrizeFailJunXianLevelTooLow = newMsgError("collect_jun_xian_prize 军衔等级太低", ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_JUN_XIAN_LEVEL_TOO_LOW) // 12-3
	ErrCollectJunXianPrizeFailPointTooLow        = newMsgError("collect_jun_xian_prize 积分太低", ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_POINT_TOO_LOW)            // 12-4
	ErrCollectJunXianPrizeFailInvalidPrize       = newMsgError("collect_jun_xian_prize 无效的奖励id", ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_INVALID_PRIZE)         // 12-5
	ErrCollectJunXianPrizeFailServerError        = newMsgError("collect_jun_xian_prize 服务器忙，请稍后再试", ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_SERVER_ERROR)       // 12-6
)

// self_record
var (
	ErrSelfRecordFailServerError = newMsgError("self_record 服务器繁忙，请稍后再试", ERR_SELF_RECORD_FAIL_SERVER_ERROR) // 32-1
)

// request_rank
var (
	ErrRequestRankFailServerError = newMsgError("request_rank 服务器繁忙，请稍后再试", ERR_REQUEST_RANK_FAIL_SERVER_ERROR) // 28-1
)

func newMsgError(msg string, buffer pbutil.StaticBuffer) *error_msg {
	return &error_msg{
		msg:  msg,
		buff: buffer,
	}
}

type error_msg struct {
	msg  string
	buff pbutil.Buffer
}

func (f *error_msg) Error() string         { return f.msg }
func (f *error_msg) ErrMsg() pbutil.Buffer { return f.buff }
