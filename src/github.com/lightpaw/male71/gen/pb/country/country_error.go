package country

import (
	"github.com/lightpaw/pbutil"
)

// request_country_prestige
var (
	ErrRequestCountryPrestigeFailSeverError = newMsgError("request_country_prestige 服务器繁忙，请稍后再试", ERR_REQUEST_COUNTRY_PRESTIGE_FAIL_SEVER_ERROR) // 18-1
)

// request_countries
var (
	ErrRequestCountriesFailSeverError = newMsgError("request_countries 服务器繁忙，请稍后再试", ERR_REQUEST_COUNTRIES_FAIL_SEVER_ERROR) // 21-1
)

// hero_change_country
var (
	ErrHeroChangeCountryFailInvalidCountry = newMsgError("hero_change_country 国家错误", ERR_HERO_CHANGE_COUNTRY_FAIL_INVALID_COUNTRY) // 24-1
	ErrHeroChangeCountryFailInGuild        = newMsgError("hero_change_country 在联盟中不能转国", ERR_HERO_CHANGE_COUNTRY_FAIL_IN_GUILD)    // 24-2
	ErrHeroChangeCountryFailInCd           = newMsgError("hero_change_country 在 cd 中", ERR_HERO_CHANGE_COUNTRY_FAIL_IN_CD)         // 24-3
	ErrHeroChangeCountryFailCostNotEnough  = newMsgError("hero_change_country 花费不够", ERR_HERO_CHANGE_COUNTRY_FAIL_COST_NOT_ENOUGH) // 24-4
	ErrHeroChangeCountryFailIsKing         = newMsgError("hero_change_country 国君不能转国", ERR_HERO_CHANGE_COUNTRY_FAIL_IS_KING)       // 24-5
)

// country_detail
var (
	ErrCountryDetailFailInvalidId = newMsgError("country_detail 国家id不存在", ERR_COUNTRY_DETAIL_FAIL_INVALID_ID) // 33-1
)

// official_appoint
var (
	ErrOfficialAppointFailInvalidCountry           = newMsgError("official_appoint 国家不存在", ERR_OFFICIAL_APPOINT_FAIL_INVALID_COUNTRY)                 // 42-1
	ErrOfficialAppointFailDeny                     = newMsgError("official_appoint 没有权限", ERR_OFFICIAL_APPOINT_FAIL_DENY)                             // 42-2
	ErrOfficialAppointFailInvalidOfficial          = newMsgError("official_appoint 职位错误", ERR_OFFICIAL_APPOINT_FAIL_INVALID_OFFICIAL)                 // 42-3
	ErrOfficialAppointFailOfficialFull             = newMsgError("official_appoint 职位满了", ERR_OFFICIAL_APPOINT_FAIL_OFFICIAL_FULL)                    // 42-4
	ErrOfficialAppointFailInvalidHero              = newMsgError("official_appoint 对方找不到人", ERR_OFFICIAL_APPOINT_FAIL_INVALID_HERO)                   // 42-5
	ErrOfficialAppointFailHeroNotSameCountry       = newMsgError("official_appoint 对方不是一国人", ERR_OFFICIAL_APPOINT_FAIL_HERO_NOT_SAME_COUNTRY)         // 42-6
	ErrOfficialAppointFailHeroIsOfficial           = newMsgError("official_appoint 对方已有职位", ERR_OFFICIAL_APPOINT_FAIL_HERO_IS_OFFICIAL)               // 42-7
	ErrOfficialAppointFailServerErr                = newMsgError("official_appoint 服务器错误", ERR_OFFICIAL_APPOINT_FAIL_SERVER_ERR)                      // 42-8
	ErrOfficialAppointFailHeroInGuildChangeCountry = newMsgError("official_appoint 对方在联盟转国中", ERR_OFFICIAL_APPOINT_FAIL_HERO_IN_GUILD_CHANGE_COUNTRY) // 42-9
)

// official_depose
var (
	ErrOfficialDeposeFailInvalidCountry     = newMsgError("official_depose 国家不存在", ERR_OFFICIAL_DEPOSE_FAIL_INVALID_COUNTRY)         // 45-1
	ErrOfficialDeposeFailDeny               = newMsgError("official_depose 没有权限", ERR_OFFICIAL_DEPOSE_FAIL_DENY)                     // 45-2
	ErrOfficialDeposeFailInvalidHero        = newMsgError("official_depose 对方找不到人", ERR_OFFICIAL_DEPOSE_FAIL_INVALID_HERO)           // 45-4
	ErrOfficialDeposeFailHeroNotSameCountry = newMsgError("official_depose 对方不是一国人", ERR_OFFICIAL_DEPOSE_FAIL_HERO_NOT_SAME_COUNTRY) // 45-7
	ErrOfficialDeposeFailHeroNotInPost      = newMsgError("official_depose 对方没有职位或不在职位上", ERR_OFFICIAL_DEPOSE_FAIL_HERO_NOT_IN_POST) // 45-6
	ErrOfficialDeposeFailIsKing             = newMsgError("official_depose 是国王", ERR_OFFICIAL_DEPOSE_FAIL_IS_KING)                   // 45-8
	ErrOfficialDeposeFailServerErr          = newMsgError("official_depose 服务器错误", ERR_OFFICIAL_DEPOSE_FAIL_SERVER_ERR)              // 45-5
	ErrOfficialDeposeFailInCd               = newMsgError("official_depose cd 中", ERR_OFFICIAL_DEPOSE_FAIL_IN_CD)                    // 45-9
)

// official_leave
var (
	ErrOfficialLeaveFailInvalidCountry = newMsgError("official_leave 国家不存在", ERR_OFFICIAL_LEAVE_FAIL_INVALID_COUNTRY) // 56-1
	ErrOfficialLeaveFailNoOfficial     = newMsgError("official_leave 没有官职", ERR_OFFICIAL_LEAVE_FAIL_NO_OFFICIAL)      // 56-2
	ErrOfficialLeaveFailIsKing         = newMsgError("official_leave 是国王", ERR_OFFICIAL_LEAVE_FAIL_IS_KING)           // 56-3
	ErrOfficialLeaveFailInCd           = newMsgError("official_leave cd 中", ERR_OFFICIAL_LEAVE_FAIL_IN_CD)            // 56-4
	ErrOfficialLeaveFailServerErr      = newMsgError("official_leave 服务器错误", ERR_OFFICIAL_LEAVE_FAIL_SERVER_ERR)      // 56-5
)

// collect_official_salary
var (
	ErrCollectOfficialSalaryFailNoCountry        = newMsgError("collect_official_salary 没有国家", ERR_COLLECT_OFFICIAL_SALARY_FAIL_NO_COUNTRY)             // 48-5
	ErrCollectOfficialSalaryFailNoOfficial       = newMsgError("collect_official_salary 没有官职", ERR_COLLECT_OFFICIAL_SALARY_FAIL_NO_OFFICIAL)            // 48-1
	ErrCollectOfficialSalaryFailAppointOnSameDay = newMsgError("collect_official_salary 任职当天不能领", ERR_COLLECT_OFFICIAL_SALARY_FAIL_APPOINT_ON_SAME_DAY) // 48-2
	ErrCollectOfficialSalaryFailAlreadyCollected = newMsgError("collect_official_salary 今天已领", ERR_COLLECT_OFFICIAL_SALARY_FAIL_ALREADY_COLLECTED)      // 48-3
	ErrCollectOfficialSalaryFailServerErr        = newMsgError("collect_official_salary 服务器错误", ERR_COLLECT_OFFICIAL_SALARY_FAIL_SERVER_ERR)            // 48-4
)

// change_name_start
var (
	ErrChangeNameStartFailNewNameLenLimit      = newMsgError("change_name_start // 长度错误", ERR_CHANGE_NAME_START_FAIL_NEW_NAME_LEN_LIMIT)           // 74-1
	ErrChangeNameStartFailNotKing              = newMsgError("change_name_start // 不是国君", ERR_CHANGE_NAME_START_FAIL_NOT_KING)                     // 74-2
	ErrChangeNameStartFailCostNotEnough        = newMsgError("change_name_start // 消耗不够", ERR_CHANGE_NAME_START_FAIL_COST_NOT_ENOUGH)              // 74-3
	ErrChangeNameStartFailInCd                 = newMsgError("change_name_start // 在改名CD中", ERR_CHANGE_NAME_START_FAIL_IN_CD)                      // 74-7
	ErrChangeNameStartFailNewNameInDefaultName = newMsgError("change_name_start 名字不能为默认其他国名", ERR_CHANGE_NAME_START_FAIL_NEW_NAME_IN_DEFAULT_NAME) // 74-5
	ErrChangeNameStartFailNewNameInCurrentName = newMsgError("change_name_start 名字不能为当前其他国名", ERR_CHANGE_NAME_START_FAIL_NEW_NAME_IN_CURRENT_NAME) // 74-6
)

// change_name_vote
var (
	ErrChangeNameVoteFailNotInVote   = newMsgError("change_name_vote // 不在投票中", ERR_CHANGE_NAME_VOTE_FAIL_NOT_IN_VOTE) // 63-1
	ErrChangeNameVoteFailVoted       = newMsgError("change_name_vote // 在同一方投过票了", ERR_CHANGE_NAME_VOTE_FAIL_VOTED)    // 63-2
	ErrChangeNameVoteFailNoVoteCount = newMsgError("change_name_vote // 没有票", ERR_CHANGE_NAME_VOTE_FAIL_NO_VOTE_COUNT) // 63-3
)

// search_to_appoint_hero_list
var (
	ErrSearchToAppointHeroListFailServerErr   = newMsgError("search_to_appoint_hero_list 服务器错误", ERR_SEARCH_TO_APPOINT_HERO_LIST_FAIL_SERVER_ERR)     // 68-1
	ErrSearchToAppointHeroListFailNameIsEmpty = newMsgError("search_to_appoint_hero_list 关键字是空的", ERR_SEARCH_TO_APPOINT_HERO_LIST_FAIL_NAME_IS_EMPTY) // 68-2
)

// default_to_appoint_hero_list
var (
	ErrDefaultToAppointHeroListFailServerErr = newMsgError("default_to_appoint_hero_list 服务器错误", ERR_DEFAULT_TO_APPOINT_HERO_LIST_FAIL_SERVER_ERR) // 71-1
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
