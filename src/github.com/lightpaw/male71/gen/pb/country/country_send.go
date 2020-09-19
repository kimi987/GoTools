package country

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
	MODULE_ID = 27

	C2S_REQUEST_COUNTRY_PRESTIGE = 16

	C2S_REQUEST_COUNTRIES = 19

	C2S_HERO_CHANGE_COUNTRY = 22

	C2S_COUNTRY_DETAIL = 31

	C2S_OFFICIAL_APPOINT = 40

	C2S_OFFICIAL_DEPOSE = 43

	C2S_OFFICIAL_LEAVE = 54

	C2S_COLLECT_OFFICIAL_SALARY = 46

	C2S_CHANGE_NAME_START = 72

	C2S_CHANGE_NAME_VOTE = 61

	C2S_SEARCH_TO_APPOINT_HERO_LIST = 66

	C2S_DEFAULT_TO_APPOINT_HERO_LIST = 69
)

func NewS2cRequestCountryPrestigeMsg(vsn int32, ids []int32, prestige []int32) pbutil.Buffer {
	msg := &S2CRequestCountryPrestigeProto{
		Vsn:      vsn,
		Ids:      ids,
		Prestige: prestige,
	}
	return NewS2cRequestCountryPrestigeProtoMsg(msg)
}

var s2c_request_country_prestige = [...]byte{27, 17} // 17
func NewS2cRequestCountryPrestigeProtoMsg(object *S2CRequestCountryPrestigeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_country_prestige[:], "s2c_request_country_prestige")

}

// 服务器繁忙，请稍后再试
var ERR_REQUEST_COUNTRY_PRESTIGE_FAIL_SEVER_ERROR = pbutil.StaticBuffer{3, 27, 18, 1} // 18-1

func NewS2cRequestCountriesMsg(vsn int32, countries *shared_proto.CountriesProto, mc_war *shared_proto.McWarProto, mc *shared_proto.MingcsProto) pbutil.Buffer {
	msg := &S2CRequestCountriesProto{
		Vsn:       vsn,
		Countries: countries,
		McWar:     mc_war,
		Mc:        mc,
	}
	return NewS2cRequestCountriesProtoMsg(msg)
}

var s2c_request_countries = [...]byte{27, 20} // 20
func NewS2cRequestCountriesProtoMsg(object *S2CRequestCountriesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_countries[:], "s2c_request_countries")

}

// 服务器繁忙，请稍后再试
var ERR_REQUEST_COUNTRIES_FAIL_SEVER_ERROR = pbutil.StaticBuffer{3, 27, 21, 1} // 21-1

func NewS2cCountriesUpdateNoticeMsg(countries *shared_proto.CountriesProto) pbutil.Buffer {
	msg := &S2CCountriesUpdateNoticeProto{
		Countries: countries,
	}
	return NewS2cCountriesUpdateNoticeProtoMsg(msg)
}

var s2c_countries_update_notice = [...]byte{27, 75} // 75
func NewS2cCountriesUpdateNoticeProtoMsg(object *S2CCountriesUpdateNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_countries_update_notice[:], "s2c_countries_update_notice")

}

func NewS2cHeroChangeCountryMsg(new_country int32, new_cd int32, normal_cd int32) pbutil.Buffer {
	msg := &S2CHeroChangeCountryProto{
		NewCountry: new_country,
		NewCd:      new_cd,
		NormalCd:   normal_cd,
	}
	return NewS2cHeroChangeCountryProtoMsg(msg)
}

var s2c_hero_change_country = [...]byte{27, 23} // 23
func NewS2cHeroChangeCountryProtoMsg(object *S2CHeroChangeCountryProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_hero_change_country[:], "s2c_hero_change_country")

}

// 国家错误
var ERR_HERO_CHANGE_COUNTRY_FAIL_INVALID_COUNTRY = pbutil.StaticBuffer{3, 27, 24, 1} // 24-1

// 在联盟中不能转国
var ERR_HERO_CHANGE_COUNTRY_FAIL_IN_GUILD = pbutil.StaticBuffer{3, 27, 24, 2} // 24-2

// 在 cd 中
var ERR_HERO_CHANGE_COUNTRY_FAIL_IN_CD = pbutil.StaticBuffer{3, 27, 24, 3} // 24-3

// 花费不够
var ERR_HERO_CHANGE_COUNTRY_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 27, 24, 4} // 24-4

// 国君不能转国
var ERR_HERO_CHANGE_COUNTRY_FAIL_IS_KING = pbutil.StaticBuffer{3, 27, 24, 5} // 24-5

func NewS2cCountryDetailMsg(country *shared_proto.CountryDetailProto) pbutil.Buffer {
	msg := &S2CCountryDetailProto{
		Country: country,
	}
	return NewS2cCountryDetailProtoMsg(msg)
}

var s2c_country_detail = [...]byte{27, 32} // 32
func NewS2cCountryDetailProtoMsg(object *S2CCountryDetailProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_country_detail[:], "s2c_country_detail")

}

// 国家id不存在
var ERR_COUNTRY_DETAIL_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 27, 33, 1} // 33-1

var OFFICIAL_APPOINT_S2C = pbutil.StaticBuffer{2, 27, 41} // 41

// 国家不存在
var ERR_OFFICIAL_APPOINT_FAIL_INVALID_COUNTRY = pbutil.StaticBuffer{3, 27, 42, 1} // 42-1

// 没有权限
var ERR_OFFICIAL_APPOINT_FAIL_DENY = pbutil.StaticBuffer{3, 27, 42, 2} // 42-2

// 职位错误
var ERR_OFFICIAL_APPOINT_FAIL_INVALID_OFFICIAL = pbutil.StaticBuffer{3, 27, 42, 3} // 42-3

// 职位满了
var ERR_OFFICIAL_APPOINT_FAIL_OFFICIAL_FULL = pbutil.StaticBuffer{3, 27, 42, 4} // 42-4

// 对方找不到人
var ERR_OFFICIAL_APPOINT_FAIL_INVALID_HERO = pbutil.StaticBuffer{3, 27, 42, 5} // 42-5

// 对方不是一国人
var ERR_OFFICIAL_APPOINT_FAIL_HERO_NOT_SAME_COUNTRY = pbutil.StaticBuffer{3, 27, 42, 6} // 42-6

// 对方已有职位
var ERR_OFFICIAL_APPOINT_FAIL_HERO_IS_OFFICIAL = pbutil.StaticBuffer{3, 27, 42, 7} // 42-7

// 服务器错误
var ERR_OFFICIAL_APPOINT_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 27, 42, 8} // 42-8

// 对方在联盟转国中
var ERR_OFFICIAL_APPOINT_FAIL_HERO_IN_GUILD_CHANGE_COUNTRY = pbutil.StaticBuffer{3, 27, 42, 9} // 42-9

func NewS2cOfficialAppointNoticeMsg(official_type int32) pbutil.Buffer {
	msg := &S2COfficialAppointNoticeProto{
		OfficialType: official_type,
	}
	return NewS2cOfficialAppointNoticeProtoMsg(msg)
}

var s2c_official_appoint_notice = [...]byte{27, 49} // 49
func NewS2cOfficialAppointNoticeProtoMsg(object *S2COfficialAppointNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_official_appoint_notice[:], "s2c_official_appoint_notice")

}

var OFFICIAL_DEPOSE_NOTICE_S2C = pbutil.StaticBuffer{2, 27, 50} // 50

var OFFICIAL_DEPOSE_S2C = pbutil.StaticBuffer{2, 27, 44} // 44

// 国家不存在
var ERR_OFFICIAL_DEPOSE_FAIL_INVALID_COUNTRY = pbutil.StaticBuffer{3, 27, 45, 1} // 45-1

// 没有权限
var ERR_OFFICIAL_DEPOSE_FAIL_DENY = pbutil.StaticBuffer{3, 27, 45, 2} // 45-2

// 对方找不到人
var ERR_OFFICIAL_DEPOSE_FAIL_INVALID_HERO = pbutil.StaticBuffer{3, 27, 45, 4} // 45-4

// 对方不是一国人
var ERR_OFFICIAL_DEPOSE_FAIL_HERO_NOT_SAME_COUNTRY = pbutil.StaticBuffer{3, 27, 45, 7} // 45-7

// 对方没有职位或不在职位上
var ERR_OFFICIAL_DEPOSE_FAIL_HERO_NOT_IN_POST = pbutil.StaticBuffer{3, 27, 45, 6} // 45-6

// 是国王
var ERR_OFFICIAL_DEPOSE_FAIL_IS_KING = pbutil.StaticBuffer{3, 27, 45, 8} // 45-8

// 服务器错误
var ERR_OFFICIAL_DEPOSE_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 27, 45, 5} // 45-5

// cd 中
var ERR_OFFICIAL_DEPOSE_FAIL_IN_CD = pbutil.StaticBuffer{3, 27, 45, 9} // 45-9

var OFFICIAL_LEAVE_S2C = pbutil.StaticBuffer{2, 27, 55} // 55

// 国家不存在
var ERR_OFFICIAL_LEAVE_FAIL_INVALID_COUNTRY = pbutil.StaticBuffer{3, 27, 56, 1} // 56-1

// 没有官职
var ERR_OFFICIAL_LEAVE_FAIL_NO_OFFICIAL = pbutil.StaticBuffer{3, 27, 56, 2} // 56-2

// 是国王
var ERR_OFFICIAL_LEAVE_FAIL_IS_KING = pbutil.StaticBuffer{3, 27, 56, 3} // 56-3

// cd 中
var ERR_OFFICIAL_LEAVE_FAIL_IN_CD = pbutil.StaticBuffer{3, 27, 56, 4} // 56-4

// 服务器错误
var ERR_OFFICIAL_LEAVE_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 27, 56, 5} // 56-5

func NewS2cCountryHostChangedNoticeMsg(country_id int32) pbutil.Buffer {
	msg := &S2CCountryHostChangedNoticeProto{
		CountryId: country_id,
	}
	return NewS2cCountryHostChangedNoticeProtoMsg(msg)
}

var s2c_country_host_changed_notice = [...]byte{27, 51} // 51
func NewS2cCountryHostChangedNoticeProtoMsg(object *S2CCountryHostChangedNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_country_host_changed_notice[:], "s2c_country_host_changed_notice")

}

func NewS2cCountryDestroyNoticeMsg(country_id int32) pbutil.Buffer {
	msg := &S2CCountryDestroyNoticeProto{
		CountryId: country_id,
	}
	return NewS2cCountryDestroyNoticeProtoMsg(msg)
}

var s2c_country_destroy_notice = [...]byte{27, 53} // 53
func NewS2cCountryDestroyNoticeProtoMsg(object *S2CCountryDestroyNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_country_destroy_notice[:], "s2c_country_destroy_notice")

}

func NewS2cKingChangedNoticeMsg(country_id int32) pbutil.Buffer {
	msg := &S2CKingChangedNoticeProto{
		CountryId: country_id,
	}
	return NewS2cKingChangedNoticeProtoMsg(msg)
}

var s2c_king_changed_notice = [...]byte{27, 52} // 52
func NewS2cKingChangedNoticeProtoMsg(object *S2CKingChangedNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_king_changed_notice[:], "s2c_king_changed_notice")

}

func NewS2cCollectOfficialSalaryMsg(salary *shared_proto.PrizeProto) pbutil.Buffer {
	msg := &S2CCollectOfficialSalaryProto{
		Salary: salary,
	}
	return NewS2cCollectOfficialSalaryProtoMsg(msg)
}

var s2c_collect_official_salary = [...]byte{27, 47} // 47
func NewS2cCollectOfficialSalaryProtoMsg(object *S2CCollectOfficialSalaryProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_official_salary[:], "s2c_collect_official_salary")

}

// 没有国家
var ERR_COLLECT_OFFICIAL_SALARY_FAIL_NO_COUNTRY = pbutil.StaticBuffer{3, 27, 48, 5} // 48-5

// 没有官职
var ERR_COLLECT_OFFICIAL_SALARY_FAIL_NO_OFFICIAL = pbutil.StaticBuffer{3, 27, 48, 1} // 48-1

// 任职当天不能领
var ERR_COLLECT_OFFICIAL_SALARY_FAIL_APPOINT_ON_SAME_DAY = pbutil.StaticBuffer{3, 27, 48, 2} // 48-2

// 今天已领
var ERR_COLLECT_OFFICIAL_SALARY_FAIL_ALREADY_COLLECTED = pbutil.StaticBuffer{3, 27, 48, 3} // 48-3

// 服务器错误
var ERR_COLLECT_OFFICIAL_SALARY_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 27, 48, 4} // 48-4

var CHANGE_NAME_START_S2C = pbutil.StaticBuffer{2, 27, 73} // 73

// // 长度错误
var ERR_CHANGE_NAME_START_FAIL_NEW_NAME_LEN_LIMIT = pbutil.StaticBuffer{3, 27, 74, 1} // 74-1

// // 不是国君
var ERR_CHANGE_NAME_START_FAIL_NOT_KING = pbutil.StaticBuffer{3, 27, 74, 2} // 74-2

// // 消耗不够
var ERR_CHANGE_NAME_START_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 27, 74, 3} // 74-3

// // 在改名CD中
var ERR_CHANGE_NAME_START_FAIL_IN_CD = pbutil.StaticBuffer{3, 27, 74, 7} // 74-7

// 名字不能为默认其他国名
var ERR_CHANGE_NAME_START_FAIL_NEW_NAME_IN_DEFAULT_NAME = pbutil.StaticBuffer{3, 27, 74, 5} // 74-5

// 名字不能为当前其他国名
var ERR_CHANGE_NAME_START_FAIL_NEW_NAME_IN_CURRENT_NAME = pbutil.StaticBuffer{3, 27, 74, 6} // 74-6

var CHANGE_NAME_START_NOTICE_S2C = pbutil.StaticBuffer{2, 27, 60} // 60

func NewS2cChangeNameVoteMsg(agree bool, count int32, reduce_count int32) pbutil.Buffer {
	msg := &S2CChangeNameVoteProto{
		Agree:       agree,
		Count:       count,
		ReduceCount: reduce_count,
	}
	return NewS2cChangeNameVoteProtoMsg(msg)
}

var s2c_change_name_vote = [...]byte{27, 62} // 62
func NewS2cChangeNameVoteProtoMsg(object *S2CChangeNameVoteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_name_vote[:], "s2c_change_name_vote")

}

// // 不在投票中
var ERR_CHANGE_NAME_VOTE_FAIL_NOT_IN_VOTE = pbutil.StaticBuffer{3, 27, 63, 1} // 63-1

// // 在同一方投过票了
var ERR_CHANGE_NAME_VOTE_FAIL_VOTED = pbutil.StaticBuffer{3, 27, 63, 2} // 63-2

// // 没有票
var ERR_CHANGE_NAME_VOTE_FAIL_NO_VOTE_COUNT = pbutil.StaticBuffer{3, 27, 63, 3} // 63-3

func NewS2cHeroChangeNameVoteCountUpdateNoticeMsg(new_count int32) pbutil.Buffer {
	msg := &S2CHeroChangeNameVoteCountUpdateNoticeProto{
		NewCount: new_count,
	}
	return NewS2cHeroChangeNameVoteCountUpdateNoticeProtoMsg(msg)
}

var s2c_hero_change_name_vote_count_update_notice = [...]byte{27, 64} // 64
func NewS2cHeroChangeNameVoteCountUpdateNoticeProtoMsg(object *S2CHeroChangeNameVoteCountUpdateNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_hero_change_name_vote_count_update_notice[:], "s2c_hero_change_name_vote_count_update_notice")

}

func NewS2cChangeNameSuccNoticeMsg(succ bool, country int32, new_name string) pbutil.Buffer {
	msg := &S2CChangeNameSuccNoticeProto{
		Succ:    succ,
		Country: country,
		NewName: new_name,
	}
	return NewS2cChangeNameSuccNoticeProtoMsg(msg)
}

var s2c_change_name_succ_notice = [...]byte{27, 65} // 65
func NewS2cChangeNameSuccNoticeProtoMsg(object *S2CChangeNameSuccNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_name_succ_notice[:], "s2c_change_name_succ_notice")

}

func NewS2cSearchToAppointHeroListMsg(heros []*shared_proto.HeroBasicSnapshotProto) pbutil.Buffer {
	msg := &S2CSearchToAppointHeroListProto{
		Heros: heros,
	}
	return NewS2cSearchToAppointHeroListProtoMsg(msg)
}

var s2c_search_to_appoint_hero_list = [...]byte{27, 67} // 67
func NewS2cSearchToAppointHeroListProtoMsg(object *S2CSearchToAppointHeroListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_search_to_appoint_hero_list[:], "s2c_search_to_appoint_hero_list")

}

// 服务器错误
var ERR_SEARCH_TO_APPOINT_HERO_LIST_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 27, 68, 1} // 68-1

// 关键字是空的
var ERR_SEARCH_TO_APPOINT_HERO_LIST_FAIL_NAME_IS_EMPTY = pbutil.StaticBuffer{3, 27, 68, 2} // 68-2

func NewS2cDefaultToAppointHeroListMsg(heros []*shared_proto.HeroBasicSnapshotProto) pbutil.Buffer {
	msg := &S2CDefaultToAppointHeroListProto{
		Heros: heros,
	}
	return NewS2cDefaultToAppointHeroListProtoMsg(msg)
}

var s2c_default_to_appoint_hero_list = [...]byte{27, 70} // 70
func NewS2cDefaultToAppointHeroListProtoMsg(object *S2CDefaultToAppointHeroListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_default_to_appoint_hero_list[:], "s2c_default_to_appoint_hero_list")

}

// 服务器错误
var ERR_DEFAULT_TO_APPOINT_HERO_LIST_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 27, 71, 1} // 71-1
