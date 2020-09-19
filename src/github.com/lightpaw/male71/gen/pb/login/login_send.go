package login

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
	MODULE_ID = 1

	C2S_INTERNAL_LOGIN = 1

	C2S_LOGIN = 7

	C2S_SET_TUTORIAL_PROGRESS = 18

	C2S_CREATE_HERO = 3

	C2S_LOADED = 10

	C2S_ROBOT_LOGIN = 13
)

func NewS2cInternalLoginMsg(heroProto []byte) pbutil.Buffer {
	msg := &S2CInternalLoginProto{
		HeroProto: heroProto,
	}
	return NewS2cInternalLoginProtoMsg(msg)
}

func NewS2cInternalLoginMarshalMsg(heroProto marshaler) pbutil.Buffer {
	msg := &S2CInternalLoginProto{
		HeroProto: safeMarshal(heroProto),
	}
	return NewS2cInternalLoginProtoMsg(msg)
}

var s2c_internal_login = [...]byte{1, 2} // 2
func NewS2cInternalLoginProtoMsg(object *S2CInternalLoginProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_internal_login[:], "s2c_internal_login")

}

// 已经登陆了，不要重复登陆
var ERR_INTERNAL_LOGIN_FAIL_ALREADY_LOGIN = pbutil.StaticBuffer{3, 1, 5, 1} // 5-1

// 发送上来的proto解析不了
var ERR_INTERNAL_LOGIN_FAIL_INVALID_PROTO = pbutil.StaticBuffer{3, 1, 5, 2} // 5-2

// 发送的id无效
var ERR_INTERNAL_LOGIN_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 1, 5, 3} // 5-3

// 被T下线
var ERR_INTERNAL_LOGIN_FAIL_KICK = pbutil.StaticBuffer{3, 1, 5, 5} // 5-5

// 服务器忙，请稍后再试
var ERR_INTERNAL_LOGIN_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 1, 5, 4} // 5-4

func NewS2cLoginMsg(created bool, male bool, head string, building []int32, is_debug bool, countries *shared_proto.CountriesProto) pbutil.Buffer {
	msg := &S2CLoginProto{
		Created:   created,
		Male:      male,
		Head:      head,
		Building:  building,
		IsDebug:   is_debug,
		Countries: countries,
	}
	return NewS2cLoginProtoMsg(msg)
}

var s2c_login = [...]byte{1, 8} // 8
func NewS2cLoginProtoMsg(object *S2CLoginProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_login[:], "s2c_login")

}

// 已经登陆了，不要重复登陆
var ERR_LOGIN_FAIL_ALREADY_LOGIN = pbutil.StaticBuffer{3, 1, 9, 1} // 9-1

// 发送上来的proto解析不了
var ERR_LOGIN_FAIL_INVALID_PROTO = pbutil.StaticBuffer{3, 1, 9, 2} // 9-2

// 发送的id无效
var ERR_LOGIN_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 1, 9, 3} // 9-3

// 被T下线
var ERR_LOGIN_FAIL_KICK = pbutil.StaticBuffer{3, 1, 9, 4} // 9-4

// 服务器忙，请稍后再试
var ERR_LOGIN_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 1, 9, 5} // 9-5

// 无效的token
var ERR_LOGIN_FAIL_INVALID_TOKEN = pbutil.StaticBuffer{3, 1, 9, 6} // 9-6

// 账号已被查封
var ERR_LOGIN_FAIL_BAN_LOGIN = pbutil.StaticBuffer{3, 1, 9, 2} // 9-2

// 无效的 tencent_info
var ERR_LOGIN_FAIL_INVALID_TENCENT_INFO = pbutil.StaticBuffer{3, 1, 9, 7} // 9-7

func NewS2cTutorialProgressMsg(progress int32) pbutil.Buffer {
	msg := &S2CTutorialProgressProto{
		Progress: progress,
	}
	return NewS2cTutorialProgressProtoMsg(msg)
}

var s2c_tutorial_progress = [...]byte{1, 17} // 17
func NewS2cTutorialProgressProtoMsg(object *S2CTutorialProgressProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_tutorial_progress[:], "s2c_tutorial_progress")

}

var CREATE_HERO_S2C = pbutil.StaticBuffer{2, 1, 4} // 4

// 还没登陆
var ERR_CREATE_HERO_FAIL_NO_LOGIN = pbutil.StaticBuffer{3, 1, 6, 1} // 6-1

// 发送上来的proto解析不了
var ERR_CREATE_HERO_FAIL_INVALID_PROTO = pbutil.StaticBuffer{3, 1, 6, 2} // 6-2

// 名字无效（长度不对）
var ERR_CREATE_HERO_FAIL_INVALID_NAME = pbutil.StaticBuffer{3, 1, 6, 3} // 6-3

// 名字已经存在
var ERR_CREATE_HERO_FAIL_NAME_EXIST = pbutil.StaticBuffer{3, 1, 6, 4} // 6-4

// 国家错误
var ERR_CREATE_HERO_FAIL_COUNTRY_ERR = pbutil.StaticBuffer{3, 1, 6, 9} // 6-9

// 英雄已经创建过了
var ERR_CREATE_HERO_FAIL_CREATED = pbutil.StaticBuffer{3, 1, 6, 6} // 6-6

// 名字包含敏感词
var ERR_CREATE_HERO_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 1, 6, 8} // 6-8

// 服务器忙，请稍后再试
var ERR_CREATE_HERO_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 1, 6, 7} // 6-7

func NewS2cLoadedMsg(heroProto []byte) pbutil.Buffer {
	msg := &S2CLoadedProto{
		HeroProto: heroProto,
	}
	return NewS2cLoadedProtoMsg(msg)
}

func NewS2cLoadedMarshalMsg(heroProto marshaler) pbutil.Buffer {
	msg := &S2CLoadedProto{
		HeroProto: safeMarshal(heroProto),
	}
	return NewS2cLoadedProtoMsg(msg)
}

var s2c_loaded = [...]byte{1, 11} // 11
func NewS2cLoadedProtoMsg(object *S2CLoadedProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_loaded[:], "s2c_loaded")

}

// 还没登陆
var ERR_LOADED_FAIL_NO_LOGIN = pbutil.StaticBuffer{3, 1, 12, 1} // 12-1

// 还没创建角色
var ERR_LOADED_FAIL_NO_CREATED = pbutil.StaticBuffer{3, 1, 12, 2} // 12-2

var ROBOT_LOGIN_S2C = pbutil.StaticBuffer{2, 1, 14} // 14

func NewS2cBanLoginMsg(result_time int32) pbutil.Buffer {
	msg := &S2CBanLoginProto{
		ResultTime: result_time,
	}
	return NewS2cBanLoginProtoMsg(msg)
}

var s2c_ban_login = [...]byte{1, 19} // 19
func NewS2cBanLoginProtoMsg(object *S2CBanLoginProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_ban_login[:], "s2c_ban_login")

}
