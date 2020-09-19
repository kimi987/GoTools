package guild

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
	MODULE_ID = 9

	C2S_LIST_GUILD = 1

	C2S_SEARCH_GUILD = 4

	C2S_CREATE_GUILD = 7

	C2S_SELF_GUILD = 10

	C2S_LEAVE_GUILD = 13

	C2S_KICK_OTHER = 17

	C2S_UPDATE_TEXT = 20

	C2S_UPDATE_INTERNAL_TEXT = 65

	C2S_UPDATE_CLASS_NAMES = 23

	C2S_UPDATE_CLASS_TITLE = 122

	C2S_UPDATE_FLAG_TYPE = 26

	C2S_UPDATE_MEMBER_CLASS_LEVEL = 29

	C2S_CANCEL_CHANGE_LEADER = 80

	C2S_UPDATE_JOIN_CONDITION = 68

	C2S_UPDATE_GUILD_NAME = 71

	C2S_UPDATE_GUILD_LABEL = 75

	C2S_DONATE = 83

	C2S_UPGRADE_LEVEL = 90

	C2S_REDUCE_UPGRADE_LEVEL_CD = 93

	C2S_IMPEACH_LEADER = 96

	C2S_IMPEACH_LEADER_VOTE = 99

	C2S_LIST_GUILD_BY_IDS = 102

	C2S_USER_REQUEST_JOIN = 40

	C2S_USER_CANCEL_JOIN_REQUEST = 43

	C2S_GUILD_REPLY_JOIN_REQUEST = 55

	C2S_GUILD_INVATE_OTHER = 109

	C2S_GUILD_CANCEL_INVATE_OTHER = 112

	C2S_USER_REPLY_INVATE_REQUEST = 48

	C2S_LIST_INVITE_ME_GUILD = 193

	C2S_UPDATE_FRIEND_GUILD = 125

	C2S_UPDATE_ENEMY_GUILD = 128

	C2S_UPDATE_GUILD_PRESTIGE = 131

	C2S_PLACE_GUILD_STATUE = 134

	C2S_TAKE_BACK_GUILD_STATUE = 138

	C2S_COLLECT_FIRST_JOIN_GUILD_PRIZE = 143

	C2S_SEEK_HELP = 147

	C2S_HELP_GUILD_MEMBER = 151

	C2S_HELP_ALL_GUILD_MEMBER = 158

	C2S_COLLECT_GUILD_EVENT_PRIZE = 163

	C2S_COLLECT_FULL_BIG_BOX = 167

	C2S_UPGRADE_TECHNOLOGY = 172

	C2S_REDUCE_TECHNOLOGY_CD = 175

	C2S_LIST_GUILD_LOGS = 178

	C2S_REQUEST_RECOMMEND_GUILD = 181

	C2S_HELP_TECH = 184

	C2S_RECOMMEND_INVITE_HEROS = 187

	C2S_SEARCH_NO_GUILD_HEROS = 190

	C2S_VIEW_MC_WAR_RECORD = 199

	C2S_UPDATE_GUILD_MARK = 196

	C2S_VIEW_YINLIANG_RECORD = 202

	C2S_SEND_YINLIANG_TO_OTHER_GUILD = 205

	C2S_SEND_YINLIANG_TO_MEMBER = 208

	C2S_PAY_SALARY = 211

	C2S_SET_SALARY = 214

	C2S_VIEW_SEND_YINLIANG_TO_GUILD = 218

	C2S_CONVENE = 228

	C2S_COLLECT_DAILY_GUILD_RANK_PRIZE = 231

	C2S_VIEW_DAILY_GUILD_RANK = 234

	C2S_ADD_RECOMMEND_MC_BUILD = 240

	C2S_VIEW_TASK_PROGRESS = 243

	C2S_COLLECT_TASK_PRIZE = 247

	C2S_GUILD_CHANGE_COUNTRY = 250

	C2S_CANCEL_GUILD_CHANGE_COUNTRY = 253

	C2S_SHOW_WORKSHOP_NOT_EXIST = 256
)

func NewS2cListGuildMsg(guild_list []*shared_proto.GuildSnapshotProto) pbutil.Buffer {
	msg := &S2CListGuildProto{
		GuildList: guild_list,
	}
	return NewS2cListGuildProtoMsg(msg)
}

var s2c_list_guild = [...]byte{9, 2} // 2
func NewS2cListGuildProtoMsg(object *S2CListGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_guild[:], "s2c_list_guild")

}

// 服务器忙，请稍后再试
var ERR_LIST_GUILD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 3, 1} // 3-1

func NewS2cSearchGuildMsg(proto [][]byte, guild_list []*shared_proto.GuildSnapshotProto, yinliang_list []*shared_proto.GuildYinliangSendProto) pbutil.Buffer {
	msg := &S2CSearchGuildProto{
		Proto:        proto,
		GuildList:    guild_list,
		YinliangList: yinliang_list,
	}
	return NewS2cSearchGuildProtoMsg(msg)
}

func NewS2cSearchGuildMarshalMsg(proto [][]byte, guild_list []*shared_proto.GuildSnapshotProto, yinliang_list []*shared_proto.GuildYinliangSendProto) pbutil.Buffer {
	msg := &S2CSearchGuildProto{
		Proto:        proto,
		GuildList:    guild_list,
		YinliangList: yinliang_list,
	}
	return NewS2cSearchGuildProtoMsg(msg)
}

var s2c_search_guild = [...]byte{9, 5} // 5
func NewS2cSearchGuildProtoMsg(object *S2CSearchGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_search_guild[:], "s2c_search_guild")

}

// 无效的搜索名字
var ERR_SEARCH_GUILD_FAIL_INVALID_NAME = pbutil.StaticBuffer{3, 9, 6, 2} // 6-2

// 无效的页数
var ERR_SEARCH_GUILD_FAIL_INVALID_NUM = pbutil.StaticBuffer{3, 9, 6, 3} // 6-3

// 服务器忙，请稍后再试
var ERR_SEARCH_GUILD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 6, 1} // 6-1

func NewS2cCreateGuildMsg(proto []byte) pbutil.Buffer {
	msg := &S2CCreateGuildProto{
		Proto: proto,
	}
	return NewS2cCreateGuildProtoMsg(msg)
}

func NewS2cCreateGuildMarshalMsg(proto marshaler) pbutil.Buffer {
	msg := &S2CCreateGuildProto{
		Proto: safeMarshal(proto),
	}
	return NewS2cCreateGuildProtoMsg(msg)
}

var s2c_create_guild = [...]byte{9, 8} // 8
func NewS2cCreateGuildProtoMsg(object *S2CCreateGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_create_guild[:], "s2c_create_guild")

}

// 已经在联盟中，不能创建联盟
var ERR_CREATE_GUILD_FAIL_IN_THE_GUILD = pbutil.StaticBuffer{3, 9, 9, 1} // 9-1

// 无效的名字长度
var ERR_CREATE_GUILD_FAIL_INVALID_NAME_LEN = pbutil.StaticBuffer{3, 9, 9, 5} // 9-5

// 联盟名字已经存在
var ERR_CREATE_GUILD_FAIL_NAME_DUPLICATE = pbutil.StaticBuffer{3, 9, 9, 2} // 9-2

// 联盟旗号已经存在
var ERR_CREATE_GUILD_FAIL_FLAG_NAME_DUPLICATE = pbutil.StaticBuffer{3, 9, 9, 3} // 9-3

// 无效的旗号长度
var ERR_CREATE_GUILD_FAIL_INVALID_FLAG_NAME_LEN = pbutil.StaticBuffer{3, 9, 9, 6} // 9-6

// 创建联盟消耗不足
var ERR_CREATE_GUILD_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 9, 9, 7} // 9-7

// 玩家没有国家
var ERR_CREATE_GUILD_FAIL_HERO_NO_COUNTRY = pbutil.StaticBuffer{3, 9, 9, 9} // 9-9

// 输入包含敏感词
var ERR_CREATE_GUILD_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 9, 9, 8} // 9-8

// 服务器忙，请稍后再试
var ERR_CREATE_GUILD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 9, 4} // 9-4

func NewS2cSelfGuildMsg(varsion int32, proto []byte) pbutil.Buffer {
	msg := &S2CSelfGuildProto{
		Varsion: varsion,
		Proto:   proto,
	}
	return NewS2cSelfGuildProtoMsg(msg)
}

func NewS2cSelfGuildMarshalMsg(varsion int32, proto marshaler) pbutil.Buffer {
	msg := &S2CSelfGuildProto{
		Varsion: varsion,
		Proto:   safeMarshal(proto),
	}
	return NewS2cSelfGuildProtoMsg(msg)
}

var s2c_self_guild = [...]byte{9, 11} // 11
func NewS2cSelfGuildProtoMsg(object *S2CSelfGuildProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_self_guild[:], "s2c_self_guild")

}

// 你没有联盟
var ERR_SELF_GUILD_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 12, 2} // 12-2

// 服务器忙，请稍后再试
var ERR_SELF_GUILD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 12, 1} // 12-1

var SELF_GUILD_SAME_VERSION_S2C = pbutil.StaticBuffer{2, 9, 87} // 87

var SELF_GUILD_CHANGED_S2C = pbutil.StaticBuffer{2, 9, 88} // 88

var LEAVE_GUILD_S2C = pbutil.StaticBuffer{2, 9, 14} // 14

// 你没有联盟
var ERR_LEAVE_GUILD_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 15, 1} // 15-1

// 你是盟主，不能退出
var ERR_LEAVE_GUILD_FAIL_LEADER = pbutil.StaticBuffer{3, 9, 15, 2} // 15-2

// 匈奴入侵防守队员，且当前活动已开启
var ERR_LEAVE_GUILD_FAIL_XIONG_NU_DEFENDER = pbutil.StaticBuffer{3, 9, 15, 4} // 15-4

// 名城战进攻盟不能解散
var ERR_LEAVE_GUILD_FAIL_IS_MC_WAR_ATK = pbutil.StaticBuffer{3, 9, 15, 5} // 15-5

// 名城战防守盟不能解散
var ERR_LEAVE_GUILD_FAIL_IS_MC_WAR_DEF = pbutil.StaticBuffer{3, 9, 15, 6} // 15-6

// 名城战战斗阶段不能退出
var ERR_LEAVE_GUILD_FAIL_IN_MC_WAR_FIGHT = pbutil.StaticBuffer{3, 9, 15, 7} // 15-7

// 参与集结不能退出联盟
var ERR_LEAVE_GUILD_FAIL_ASSEMBLY = pbutil.StaticBuffer{3, 9, 15, 8} // 15-8

// 服务器忙，请稍后再试
var ERR_LEAVE_GUILD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 15, 3} // 15-3

func NewS2cLeaveGuildForOtherMsg(id []byte, name string, head string) pbutil.Buffer {
	msg := &S2CLeaveGuildForOtherProto{
		Id:   id,
		Name: name,
		Head: head,
	}
	return NewS2cLeaveGuildForOtherProtoMsg(msg)
}

var s2c_leave_guild_for_other = [...]byte{9, 16} // 16
func NewS2cLeaveGuildForOtherProtoMsg(object *S2CLeaveGuildForOtherProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_leave_guild_for_other[:], "s2c_leave_guild_for_other")

}

func NewS2cKickOtherMsg(id []byte, name string) pbutil.Buffer {
	msg := &S2CKickOtherProto{
		Id:   id,
		Name: name,
	}
	return NewS2cKickOtherProtoMsg(msg)
}

var s2c_kick_other = [...]byte{9, 18} // 18
func NewS2cKickOtherProtoMsg(object *S2CKickOtherProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_kick_other[:], "s2c_kick_other")

}

// 你没有联盟
var ERR_KICK_OTHER_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 19, 1} // 19-1

// Npc联盟不允许操作
var ERR_KICK_OTHER_FAIL_NPC = pbutil.StaticBuffer{3, 9, 19, 5} // 19-5

// 你没有权限操作
var ERR_KICK_OTHER_FAIL_DENY = pbutil.StaticBuffer{3, 9, 19, 2} // 19-2

// 盟主弹劾期间不能踢人
var ERR_KICK_OTHER_FAIL_IMPEACH_LEADER = pbutil.StaticBuffer{3, 9, 19, 7} // 19-7

// 超出每日踢人上限
var ERR_KICK_OTHER_FAIL_LIMIT = pbutil.StaticBuffer{3, 9, 19, 6} // 19-6

// 目标不在你的联盟
var ERR_KICK_OTHER_FAIL_TARGET_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 19, 3} // 19-3

// 匈奴入侵防守队员，且当前活动已开启
var ERR_KICK_OTHER_FAIL_XIONG_NU_DEFENDER = pbutil.StaticBuffer{3, 9, 19, 8} // 19-8

// 名城战战斗阶段不能踢人
var ERR_KICK_OTHER_FAIL_IN_MC_WAR_FIGHT = pbutil.StaticBuffer{3, 9, 19, 9} // 19-9

// 参与集结成员不能踢出联盟
var ERR_KICK_OTHER_FAIL_ASSEMBLY = pbutil.StaticBuffer{3, 9, 19, 10} // 19-10

// 服务器忙，请稍后再试
var ERR_KICK_OTHER_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 19, 4} // 19-4

var SELF_BEEN_KICKED_S2C = pbutil.StaticBuffer{2, 9, 89} // 89

func NewS2cUpdateTextMsg(text string) pbutil.Buffer {
	msg := &S2CUpdateTextProto{
		Text: text,
	}
	return NewS2cUpdateTextProtoMsg(msg)
}

var s2c_update_text = [...]byte{9, 21} // 21
func NewS2cUpdateTextProtoMsg(object *S2CUpdateTextProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_text[:], "s2c_update_text")

}

// 内容太长
var ERR_UPDATE_TEXT_FAIL_TEXT_TOO_LONG = pbutil.StaticBuffer{3, 9, 22, 4} // 22-4

// 你没有联盟
var ERR_UPDATE_TEXT_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 22, 1} // 22-1

// 你没有权限操作
var ERR_UPDATE_TEXT_FAIL_DENY = pbutil.StaticBuffer{3, 9, 22, 2} // 22-2

// Npc联盟不允许操作
var ERR_UPDATE_TEXT_FAIL_NPC = pbutil.StaticBuffer{3, 9, 22, 5} // 22-5

// 输入包含敏感词
var ERR_UPDATE_TEXT_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 9, 22, 6} // 22-6

// 服务器忙，请稍后再试
var ERR_UPDATE_TEXT_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 22, 3} // 22-3

func NewS2cUpdateInternalTextMsg(text string) pbutil.Buffer {
	msg := &S2CUpdateInternalTextProto{
		Text: text,
	}
	return NewS2cUpdateInternalTextProtoMsg(msg)
}

var s2c_update_internal_text = [...]byte{9, 66} // 66
func NewS2cUpdateInternalTextProtoMsg(object *S2CUpdateInternalTextProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_internal_text[:], "s2c_update_internal_text")

}

// 内容太长
var ERR_UPDATE_INTERNAL_TEXT_FAIL_TEXT_TOO_LONG = pbutil.StaticBuffer{3, 9, 67, 1} // 67-1

// 你没有联盟
var ERR_UPDATE_INTERNAL_TEXT_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 67, 2} // 67-2

// 你没有权限操作
var ERR_UPDATE_INTERNAL_TEXT_FAIL_DENY = pbutil.StaticBuffer{3, 9, 67, 3} // 67-3

// Npc联盟不允许操作
var ERR_UPDATE_INTERNAL_TEXT_FAIL_NPC = pbutil.StaticBuffer{3, 9, 67, 5} // 67-5

// 输入包含敏感词
var ERR_UPDATE_INTERNAL_TEXT_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 9, 67, 6} // 67-6

// 服务器忙，请稍后再试
var ERR_UPDATE_INTERNAL_TEXT_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 67, 4} // 67-4

func NewS2cUpdateClassNamesMsg(name []string) pbutil.Buffer {
	msg := &S2CUpdateClassNamesProto{
		Name: name,
	}
	return NewS2cUpdateClassNamesProtoMsg(msg)
}

var s2c_update_class_names = [...]byte{9, 24} // 24
func NewS2cUpdateClassNamesProtoMsg(object *S2CUpdateClassNamesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_class_names[:], "s2c_update_class_names")

}

// 你没有联盟
var ERR_UPDATE_CLASS_NAMES_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 25, 1} // 25-1

// 你没有权限操作
var ERR_UPDATE_CLASS_NAMES_FAIL_DENY = pbutil.StaticBuffer{3, 9, 25, 2} // 25-2

// 阶级个数无效
var ERR_UPDATE_CLASS_NAMES_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 9, 25, 3} // 25-3

// 阶级名称无效（空或重名）
var ERR_UPDATE_CLASS_NAMES_FAIL_INVALID_DUPLICATE = pbutil.StaticBuffer{3, 9, 25, 5} // 25-5

// Npc联盟不允许操作
var ERR_UPDATE_CLASS_NAMES_FAIL_NPC = pbutil.StaticBuffer{3, 9, 25, 6} // 25-6

// 输入包含敏感词
var ERR_UPDATE_CLASS_NAMES_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 9, 25, 7} // 25-7

// 服务器忙，请稍后再试
var ERR_UPDATE_CLASS_NAMES_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 25, 4} // 25-4

var UPDATE_CLASS_TITLE_S2C = pbutil.StaticBuffer{2, 9, 123} // 123

// 无效的proto
var ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_PROTO = pbutil.StaticBuffer{3, 9, 124, 5} // 124-5

// 无效的系统职称id
var ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_TITLE_ID = pbutil.StaticBuffer{3, 9, 124, 6} // 124-6

// 无效的联盟成员id
var ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_MEMBER_ID = pbutil.StaticBuffer{3, 9, 124, 7} // 124-7

// 不在联盟中
var ERR_UPDATE_CLASS_TITLE_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 124, 1} // 124-1

// 没有权限
var ERR_UPDATE_CLASS_TITLE_FAIL_DENY = pbutil.StaticBuffer{3, 9, 124, 2} // 124-2

// 自定义职称个数无效
var ERR_UPDATE_CLASS_TITLE_FAIL_COUNT_LIMIT = pbutil.StaticBuffer{3, 9, 124, 8} // 124-8

// 职称名字已经被使用了
var ERR_UPDATE_CLASS_TITLE_FAIL_NAME_EXIST = pbutil.StaticBuffer{3, 9, 124, 3} // 124-3

// Npc联盟不允许操作
var ERR_UPDATE_CLASS_TITLE_FAIL_NPC = pbutil.StaticBuffer{3, 9, 124, 9} // 124-9

// 输入包含敏感词
var ERR_UPDATE_CLASS_TITLE_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 9, 124, 10} // 124-10

// 服务器忙，请稍后再试
var ERR_UPDATE_CLASS_TITLE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 124, 4} // 124-4

func NewS2cUpdateFlagTypeMsg(flag_type int32) pbutil.Buffer {
	msg := &S2CUpdateFlagTypeProto{
		FlagType: flag_type,
	}
	return NewS2cUpdateFlagTypeProtoMsg(msg)
}

var s2c_update_flag_type = [...]byte{9, 27} // 27
func NewS2cUpdateFlagTypeProtoMsg(object *S2CUpdateFlagTypeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_flag_type[:], "s2c_update_flag_type")

}

// 你没有联盟
var ERR_UPDATE_FLAG_TYPE_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 28, 1} // 28-1

// 你没有权限操作
var ERR_UPDATE_FLAG_TYPE_FAIL_DENY = pbutil.StaticBuffer{3, 9, 28, 2} // 28-2

// 旗帜类型无效
var ERR_UPDATE_FLAG_TYPE_FAIL_INVALID_TYPE = pbutil.StaticBuffer{3, 9, 28, 3} // 28-3

// Npc联盟不允许操作
var ERR_UPDATE_FLAG_TYPE_FAIL_NPC = pbutil.StaticBuffer{3, 9, 28, 5} // 28-5

// 服务器忙，请稍后再试
var ERR_UPDATE_FLAG_TYPE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 28, 4} // 28-4

func NewS2cUpdateMemberClassLevelMsg(id []byte, class_level int32) pbutil.Buffer {
	msg := &S2CUpdateMemberClassLevelProto{
		Id:         id,
		ClassLevel: class_level,
	}
	return NewS2cUpdateMemberClassLevelProtoMsg(msg)
}

var s2c_update_member_class_level = [...]byte{9, 30} // 30
func NewS2cUpdateMemberClassLevelProtoMsg(object *S2CUpdateMemberClassLevelProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_member_class_level[:], "s2c_update_member_class_level")

}

// 无效的阶级
var ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_INVALID_CLASS_LEVEL = pbutil.StaticBuffer{3, 9, 31, 7} // 31-7

// 你没有联盟
var ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 31, 1} // 31-1

// 你没有权限操作
var ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_DENY = pbutil.StaticBuffer{3, 9, 31, 2} // 31-2

// 目标不在你的联盟
var ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_TARGET_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 31, 3} // 31-3

// 目标当前就是这个阶级
var ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_TARGET_SAME_CLASS_LEVEL = pbutil.StaticBuffer{3, 9, 31, 8} // 31-8

// 权限不足（目标权限不比你低）
var ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_DENY_TARGET = pbutil.StaticBuffer{3, 9, 31, 4} // 31-4

// 目标阶级已经满员
var ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_CLASS_FULL = pbutil.StaticBuffer{3, 9, 31, 5} // 31-5

// Npc联盟不允许操作
var ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_NPC = pbutil.StaticBuffer{3, 9, 31, 9} // 31-9

// 服务器忙，请稍后再试
var ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 31, 6} // 31-6

// 禅让国王，在CD中
var ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_CHANGE_KING_IN_CD = pbutil.StaticBuffer{3, 9, 31, 10} // 31-10

func NewS2cUpdateSelfClassLevelMsg(class_level int32) pbutil.Buffer {
	msg := &S2CUpdateSelfClassLevelProto{
		ClassLevel: class_level,
	}
	return NewS2cUpdateSelfClassLevelProtoMsg(msg)
}

var s2c_update_self_class_level = [...]byte{9, 130, 2} // 258
func NewS2cUpdateSelfClassLevelProtoMsg(object *S2CUpdateSelfClassLevelProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_self_class_level[:], "s2c_update_self_class_level")

}

var CANCEL_CHANGE_LEADER_S2C = pbutil.StaticBuffer{2, 9, 81} // 81

// 你没有联盟
var ERR_CANCEL_CHANGE_LEADER_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 82, 1} // 82-1

// 你不是盟主，不能取消
var ERR_CANCEL_CHANGE_LEADER_FAIL_DENY = pbutil.StaticBuffer{3, 9, 82, 2} // 82-2

// 当前没有禅让倒计时
var ERR_CANCEL_CHANGE_LEADER_FAIL_NOT_CHANGE_LEADER = pbutil.StaticBuffer{3, 9, 82, 3} // 82-3

// 服务器忙，请稍后再试
var ERR_CANCEL_CHANGE_LEADER_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 82, 4} // 82-4

func NewS2cUpdateJoinConditionMsg(reject_auto_join bool, required_hero_level int32, required_jun_xian_level int32, required_tower_max_floor int32) pbutil.Buffer {
	msg := &S2CUpdateJoinConditionProto{
		RejectAutoJoin:        reject_auto_join,
		RequiredHeroLevel:     required_hero_level,
		RequiredJunXianLevel:  required_jun_xian_level,
		RequiredTowerMaxFloor: required_tower_max_floor,
	}
	return NewS2cUpdateJoinConditionProtoMsg(msg)
}

var s2c_update_join_condition = [...]byte{9, 69} // 69
func NewS2cUpdateJoinConditionProtoMsg(object *S2CUpdateJoinConditionProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_join_condition[:], "s2c_update_join_condition")

}

// 你没有联盟
var ERR_UPDATE_JOIN_CONDITION_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 70, 1} // 70-1

// 你没有权限操作
var ERR_UPDATE_JOIN_CONDITION_FAIL_DENY = pbutil.StaticBuffer{3, 9, 70, 2} // 70-2

// 无效的君主等级
var ERR_UPDATE_JOIN_CONDITION_FAIL_INVALID_HERO_LEVEL = pbutil.StaticBuffer{3, 9, 70, 4} // 70-4

// 无效的百战军衔等级
var ERR_UPDATE_JOIN_CONDITION_FAIL_INVALID_JUN_XIAN_LEVEL = pbutil.StaticBuffer{3, 9, 70, 5} // 70-5

// Npc联盟不允许操作
var ERR_UPDATE_JOIN_CONDITION_FAIL_NPC = pbutil.StaticBuffer{3, 9, 70, 6} // 70-6

// 服务器忙，请稍后再试
var ERR_UPDATE_JOIN_CONDITION_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 70, 3} // 70-3

func NewS2cUpdateGuildNameMsg(next_can_change_time int32) pbutil.Buffer {
	msg := &S2CUpdateGuildNameProto{
		NextCanChangeTime: next_can_change_time,
	}
	return NewS2cUpdateGuildNameProtoMsg(msg)
}

var s2c_update_guild_name = [...]byte{9, 72} // 72
func NewS2cUpdateGuildNameProtoMsg(object *S2CUpdateGuildNameProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_guild_name[:], "s2c_update_guild_name")

}

// 无效的联盟名字
var ERR_UPDATE_GUILD_NAME_FAIL_INVALID_NAME = pbutil.StaticBuffer{3, 9, 73, 1} // 73-1

// 无效的联盟旗号
var ERR_UPDATE_GUILD_NAME_FAIL_INVALID_FLAG_NAME = pbutil.StaticBuffer{3, 9, 73, 2} // 73-2

// 你没有联盟
var ERR_UPDATE_GUILD_NAME_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 73, 3} // 73-3

// 你没有权限操作
var ERR_UPDATE_GUILD_NAME_FAIL_DENY = pbutil.StaticBuffer{3, 9, 73, 4} // 73-4

// 联盟名字已存在
var ERR_UPDATE_GUILD_NAME_FAIL_EXIST_NAME = pbutil.StaticBuffer{3, 9, 73, 5} // 73-5

// 联盟旗号已存在
var ERR_UPDATE_GUILD_NAME_FAIL_EXIST_FLAG_NAME = pbutil.StaticBuffer{3, 9, 73, 6} // 73-6

// 改名消耗不足
var ERR_UPDATE_GUILD_NAME_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 9, 73, 7} // 73-7

// 改名CD中
var ERR_UPDATE_GUILD_NAME_FAIL_COOLDOWN = pbutil.StaticBuffer{3, 9, 73, 8} // 73-8

// Npc联盟不允许操作
var ERR_UPDATE_GUILD_NAME_FAIL_NPC = pbutil.StaticBuffer{3, 9, 73, 10} // 73-10

// 输入包含敏感词
var ERR_UPDATE_GUILD_NAME_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 9, 73, 11} // 73-11

// 服务器忙，请稍后再试
var ERR_UPDATE_GUILD_NAME_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 73, 9} // 73-9

func NewS2cUpdateGuildNameBroadcastMsg(id int32, name string, flag_name string) pbutil.Buffer {
	msg := &S2CUpdateGuildNameBroadcastProto{
		Id:       id,
		Name:     name,
		FlagName: flag_name,
	}
	return NewS2cUpdateGuildNameBroadcastProtoMsg(msg)
}

var s2c_update_guild_name_broadcast = [...]byte{9, 74} // 74
func NewS2cUpdateGuildNameBroadcastProtoMsg(object *S2CUpdateGuildNameBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_guild_name_broadcast[:], "s2c_update_guild_name_broadcast")

}

func NewS2cUpdateGuildLabelMsg(label []string) pbutil.Buffer {
	msg := &S2CUpdateGuildLabelProto{
		Label: label,
	}
	return NewS2cUpdateGuildLabelProtoMsg(msg)
}

var s2c_update_guild_label = [...]byte{9, 76} // 76
func NewS2cUpdateGuildLabelProtoMsg(object *S2CUpdateGuildLabelProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_guild_label[:], "s2c_update_guild_label")

}

// 你没有联盟
var ERR_UPDATE_GUILD_LABEL_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 77, 1} // 77-1

// 你没有权限操作
var ERR_UPDATE_GUILD_LABEL_FAIL_DENY = pbutil.StaticBuffer{3, 9, 77, 2} // 77-2

// 标签个数超出上限
var ERR_UPDATE_GUILD_LABEL_FAIL_COUNT_LIMIT = pbutil.StaticBuffer{3, 9, 77, 3} // 77-3

// 标签字数超出上限
var ERR_UPDATE_GUILD_LABEL_FAIL_CHAR_LIMIT = pbutil.StaticBuffer{3, 9, 77, 7} // 77-7

// 标签重名
var ERR_UPDATE_GUILD_LABEL_FAIL_DUPLICATE = pbutil.StaticBuffer{3, 9, 77, 5} // 77-5

// Npc联盟不允许操作
var ERR_UPDATE_GUILD_LABEL_FAIL_NPC = pbutil.StaticBuffer{3, 9, 77, 6} // 77-6

// 输入包含敏感词
var ERR_UPDATE_GUILD_LABEL_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 9, 77, 8} // 77-8

// 服务器忙，请稍后再试
var ERR_UPDATE_GUILD_LABEL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 77, 4} // 77-4

func NewS2cUpdateContributionCoinMsg(coin int32) pbutil.Buffer {
	msg := &S2CUpdateContributionCoinProto{
		Coin: coin,
	}
	return NewS2cUpdateContributionCoinProtoMsg(msg)
}

var s2c_update_contribution_coin = [...]byte{9, 86} // 86
func NewS2cUpdateContributionCoinProtoMsg(object *S2CUpdateContributionCoinProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_contribution_coin[:], "s2c_update_contribution_coin")

}

func NewS2cDonateMsg(sequence int32, times int32, donate_id int32, building_amount int32, contribution_amount int32, contribution_total_amount int32, contribution_amount7 int32, donation_amount int32, donation_total_amount int32, donation_amount7 int32, donation_total_yuanbao int32) pbutil.Buffer {
	msg := &S2CDonateProto{
		Sequence:                sequence,
		Times:                   times,
		DonateId:                donate_id,
		BuildingAmount:          building_amount,
		ContributionAmount:      contribution_amount,
		ContributionTotalAmount: contribution_total_amount,
		ContributionAmount7:     contribution_amount7,
		DonationAmount:          donation_amount,
		DonationTotalAmount:     donation_total_amount,
		DonationAmount7:         donation_amount7,
		DonationTotalYuanbao:    donation_total_yuanbao,
	}
	return NewS2cDonateProtoMsg(msg)
}

var s2c_donate = [...]byte{9, 84} // 84
func NewS2cDonateProtoMsg(object *S2CDonateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_donate[:], "s2c_donate")

}

// 你没有联盟
var ERR_DONATE_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 85, 4} // 85-4

// 无效的序号
var ERR_DONATE_FAIL_INVALID_SEQUENCE = pbutil.StaticBuffer{3, 9, 85, 1} // 85-1

// 已经达到最大捐献次数
var ERR_DONATE_FAIL_MAX_TIMES = pbutil.StaticBuffer{3, 9, 85, 2} // 85-2

// 消耗不足
var ERR_DONATE_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 9, 85, 3} // 85-3

// 捐献次数已经达到外使院捐献次数上限
var ERR_DONATE_FAIL_DONATE_TIMES_LIMIT = pbutil.StaticBuffer{3, 9, 85, 6} // 85-6

// 君主等级不够，无法捐献
var ERR_DONATE_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 9, 85, 7} // 85-7

// 服务器忙，请稍后再试
var ERR_DONATE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 85, 5} // 85-5

var UPGRADE_LEVEL_S2C = pbutil.StaticBuffer{2, 9, 91} // 91

// 你没有联盟
var ERR_UPGRADE_LEVEL_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 92, 1} // 92-1

// 你没有权限操作
var ERR_UPGRADE_LEVEL_FAIL_DENY = pbutil.StaticBuffer{3, 9, 92, 2} // 92-2

// 正在升级中
var ERR_UPGRADE_LEVEL_FAIL_UPGRADING = pbutil.StaticBuffer{3, 9, 92, 3} // 92-3

// 帮派已经达到最高级
var ERR_UPGRADE_LEVEL_FAIL_MAX_LEVEL = pbutil.StaticBuffer{3, 9, 92, 4} // 92-4

// 建设值不足
var ERR_UPGRADE_LEVEL_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 9, 92, 6} // 92-6

// Npc联盟不允许操作
var ERR_UPGRADE_LEVEL_FAIL_NPC = pbutil.StaticBuffer{3, 9, 92, 7} // 92-7

// 服务器忙，请稍后再试
var ERR_UPGRADE_LEVEL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 92, 5} // 92-5

var REDUCE_UPGRADE_LEVEL_CD_S2C = pbutil.StaticBuffer{2, 9, 94} // 94

// 你没有联盟
var ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 95, 1} // 95-1

// 你没有权限操作
var ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_DENY = pbutil.StaticBuffer{3, 9, 95, 2} // 95-2

// 联盟没有在升级，不能加速
var ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_NO_UPGRADING = pbutil.StaticBuffer{3, 9, 95, 6} // 95-6

// 帮派已经达到最大加速次数
var ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_MAX_TIMES = pbutil.StaticBuffer{3, 9, 95, 7} // 95-7

// 建设值不足
var ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 9, 95, 8} // 95-8

// Npc联盟不允许操作
var ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_NPC = pbutil.StaticBuffer{3, 9, 95, 9} // 95-9

// 服务器忙，请稍后再试
var ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 95, 5} // 95-5

var IMPEACH_LEADER_S2C = pbutil.StaticBuffer{2, 9, 97} // 97

// 你没有联盟
var ERR_IMPEACH_LEADER_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 98, 1} // 98-1

// 弹劾条件未满足
var ERR_IMPEACH_LEADER_FAIL_CONDITION_NOT_REACH = pbutil.StaticBuffer{3, 9, 98, 2} // 98-2

// 你没有权限操作
var ERR_IMPEACH_LEADER_FAIL_DENY = pbutil.StaticBuffer{3, 9, 98, 3} // 98-3

// 已经存在弹劾盟主
var ERR_IMPEACH_LEADER_FAIL_IMPEACH_EXIST = pbutil.StaticBuffer{3, 9, 98, 5} // 98-5

// 正在禅让盟主，不允许弹劾
var ERR_IMPEACH_LEADER_FAIL_CHANGING_LEADER = pbutil.StaticBuffer{3, 9, 98, 7} // 98-7

// 今日弹劾盟主时间已过，请明日再试
var ERR_IMPEACH_LEADER_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 9, 98, 6} // 98-6

// 服务器忙，请稍后再试
var ERR_IMPEACH_LEADER_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 98, 4} // 98-4

func NewS2cImpeachLeaderVoteMsg(impeach_end bool, impeach []byte) pbutil.Buffer {
	msg := &S2CImpeachLeaderVoteProto{
		ImpeachEnd: impeach_end,
		Impeach:    impeach,
	}
	return NewS2cImpeachLeaderVoteProtoMsg(msg)
}

var s2c_impeach_leader_vote = [...]byte{9, 100} // 100
func NewS2cImpeachLeaderVoteProtoMsg(object *S2CImpeachLeaderVoteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_impeach_leader_vote[:], "s2c_impeach_leader_vote")

}

// 你没有联盟
var ERR_IMPEACH_LEADER_VOTE_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 101, 1} // 101-1

// 无效的投票目标
var ERR_IMPEACH_LEADER_VOTE_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 9, 101, 4} // 101-4

// 当前没有弹劾盟主
var ERR_IMPEACH_LEADER_VOTE_FAIL_IMPEACH_NOT_EXIST = pbutil.StaticBuffer{3, 9, 101, 2} // 101-2

// 无法投票给原盟主
var ERR_IMPEACH_LEADER_VOTE_FAIL_OLD_LEADER = pbutil.StaticBuffer{3, 9, 101, 5} // 101-5

// 服务器忙，请稍后再试
var ERR_IMPEACH_LEADER_VOTE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 101, 3} // 101-3

func NewS2cListGuildByIdsMsg(guilds [][]byte) pbutil.Buffer {
	msg := &S2CListGuildByIdsProto{
		Guilds: guilds,
	}
	return NewS2cListGuildByIdsProtoMsg(msg)
}

func NewS2cListGuildByIdsMarshalMsg(guilds [][]byte) pbutil.Buffer {
	msg := &S2CListGuildByIdsProto{
		Guilds: guilds,
	}
	return NewS2cListGuildByIdsProtoMsg(msg)
}

var s2c_list_guild_by_ids = [...]byte{9, 103} // 103
func NewS2cListGuildByIdsProtoMsg(object *S2CListGuildByIdsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_guild_by_ids[:], "s2c_list_guild_by_ids")

}

// 无效的id
var ERR_LIST_GUILD_BY_IDS_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 9, 104, 2} // 104-2

// 无效的id个数
var ERR_LIST_GUILD_BY_IDS_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 9, 104, 3} // 104-3

// 服务器忙，请稍后再试
var ERR_LIST_GUILD_BY_IDS_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 104, 1} // 104-1

func NewS2cUserRequestJoinMsg(id int32) pbutil.Buffer {
	msg := &S2CUserRequestJoinProto{
		Id: id,
	}
	return NewS2cUserRequestJoinProtoMsg(msg)
}

var s2c_user_request_join = [...]byte{9, 41} // 41
func NewS2cUserRequestJoinProtoMsg(object *S2CUserRequestJoinProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_user_request_join[:], "s2c_user_request_join")

}

// 无效的联盟id
var ERR_USER_REQUEST_JOIN_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 9, 42, 2} // 42-2

// 已达申请上限，请取消其他申请
var ERR_USER_REQUEST_JOIN_FAIL_SELF_FULL = pbutil.StaticBuffer{3, 9, 42, 4} // 42-4

// 不能申请加入自己的联盟
var ERR_USER_REQUEST_JOIN_FAIL_SELF_GUILD = pbutil.StaticBuffer{3, 9, 42, 7} // 42-7

// 要求未达到
var ERR_USER_REQUEST_JOIN_FAIL_CONDITION = pbutil.StaticBuffer{3, 9, 42, 5} // 42-5

// 申请的帮派已经满员
var ERR_USER_REQUEST_JOIN_FAIL_FULL = pbutil.StaticBuffer{3, 9, 42, 3} // 42-3

// 这个联盟已经申请了，不要重复申请
var ERR_USER_REQUEST_JOIN_FAIL_DUPLICATE = pbutil.StaticBuffer{3, 9, 42, 8} // 42-8

// 盟主不能申请加入其它帮派
var ERR_USER_REQUEST_JOIN_FAIL_LEADER = pbutil.StaticBuffer{3, 9, 42, 9} // 42-9

// 这个联盟是纯Npc联盟，不允许操作
var ERR_USER_REQUEST_JOIN_FAIL_NPC = pbutil.StaticBuffer{3, 9, 42, 10} // 42-10

// 匈奴入侵防守队员，且当前活动已开启
var ERR_USER_REQUEST_JOIN_FAIL_XIONG_NU_DEFENDER = pbutil.StaticBuffer{3, 9, 42, 11} // 42-11

// 离开此联盟不足4小时，不能加入
var ERR_USER_REQUEST_JOIN_FAIL_LEAVE_CD = pbutil.StaticBuffer{3, 9, 42, 12} // 42-12

// 不能加入其他国家的联盟
var ERR_USER_REQUEST_JOIN_FAIL_COUNTRY = pbutil.StaticBuffer{3, 9, 42, 13} // 42-13

// 服务器忙，请稍后再试
var ERR_USER_REQUEST_JOIN_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 42, 6} // 42-6

func NewS2cUserRemoveJoinRequestMsg(id int32) pbutil.Buffer {
	msg := &S2CUserRemoveJoinRequestProto{
		Id: id,
	}
	return NewS2cUserRemoveJoinRequestProtoMsg(msg)
}

var s2c_user_remove_join_request = [...]byte{9, 118} // 118
func NewS2cUserRemoveJoinRequestProtoMsg(object *S2CUserRemoveJoinRequestProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_user_remove_join_request[:], "s2c_user_remove_join_request")

}

var USER_CLEAR_JOIN_REQUEST_S2C = pbutil.StaticBuffer{2, 9, 119} // 119

func NewS2cUserCancelJoinRequestMsg(id int32) pbutil.Buffer {
	msg := &S2CUserCancelJoinRequestProto{
		Id: id,
	}
	return NewS2cUserCancelJoinRequestProtoMsg(msg)
}

var s2c_user_cancel_join_request = [...]byte{9, 44} // 44
func NewS2cUserCancelJoinRequestProtoMsg(object *S2CUserCancelJoinRequestProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_user_cancel_join_request[:], "s2c_user_cancel_join_request")

}

// 在联盟中
var ERR_USER_CANCEL_JOIN_REQUEST_FAIL_IN_THE_GUILD = pbutil.StaticBuffer{3, 9, 45, 1} // 45-1

// 无效的id
var ERR_USER_CANCEL_JOIN_REQUEST_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 9, 45, 2} // 45-2

// 服务器忙，请稍后再试
var ERR_USER_CANCEL_JOIN_REQUEST_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 45, 3} // 45-3

func NewS2cAddGuildMemberMsg(id []byte, name string, head string) pbutil.Buffer {
	msg := &S2CAddGuildMemberProto{
		Id:   id,
		Name: name,
		Head: head,
	}
	return NewS2cAddGuildMemberProtoMsg(msg)
}

var s2c_add_guild_member = [...]byte{9, 35} // 35
func NewS2cAddGuildMemberProtoMsg(object *S2CAddGuildMemberProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_guild_member[:], "s2c_add_guild_member")

}

func NewS2cUserJoinedMsg(id int32, name string, flag_name string, country int32) pbutil.Buffer {
	msg := &S2CUserJoinedProto{
		Id:       id,
		Name:     name,
		FlagName: flag_name,
		Country:  country,
	}
	return NewS2cUserJoinedProtoMsg(msg)
}

var s2c_user_joined = [...]byte{9, 36} // 36
func NewS2cUserJoinedProtoMsg(object *S2CUserJoinedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_user_joined[:], "s2c_user_joined")

}

func NewS2cGuildReplyJoinRequestMsg(id []byte, agree bool) pbutil.Buffer {
	msg := &S2CGuildReplyJoinRequestProto{
		Id:    id,
		Agree: agree,
	}
	return NewS2cGuildReplyJoinRequestProtoMsg(msg)
}

var s2c_guild_reply_join_request = [...]byte{9, 56} // 56
func NewS2cGuildReplyJoinRequestProtoMsg(object *S2CGuildReplyJoinRequestProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_guild_reply_join_request[:], "s2c_guild_reply_join_request")

}

// 不在联盟中
var ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_NO_GUILD = pbutil.StaticBuffer{3, 9, 57, 1} // 57-1

// 没有权限
var ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_DENY = pbutil.StaticBuffer{3, 9, 57, 2} // 57-2

// 玩家已经取消了申请
var ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_INVALID_REQUEST = pbutil.StaticBuffer{3, 9, 57, 4} // 57-4

// 匈奴入侵防守队员，且当前活动已开启
var ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_XIONG_NU_DEFENDER = pbutil.StaticBuffer{3, 9, 57, 5} // 57-5

// 玩家离开本联盟不足4小时，不能加入
var ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_LEAVE_CD = pbutil.StaticBuffer{3, 9, 57, 6} // 57-6

// 服务器忙，请稍后再试
var ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 57, 3} // 57-3

// 联盟已经满员
var ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_FULL_MEMBER = pbutil.StaticBuffer{3, 9, 57, 7} // 57-7

func NewS2cGuildInvateOtherMsg(id []byte) pbutil.Buffer {
	msg := &S2CGuildInvateOtherProto{
		Id: id,
	}
	return NewS2cGuildInvateOtherProtoMsg(msg)
}

var s2c_guild_invate_other = [...]byte{9, 110} // 110
func NewS2cGuildInvateOtherProtoMsg(object *S2CGuildInvateOtherProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_guild_invate_other[:], "s2c_guild_invate_other")

}

// 无效的玩家id
var ERR_GUILD_INVATE_OTHER_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 9, 111, 8} // 111-8

// 不在联盟中
var ERR_GUILD_INVATE_OTHER_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 111, 7} // 111-7

// 没有权限
var ERR_GUILD_INVATE_OTHER_FAIL_DENY = pbutil.StaticBuffer{3, 9, 111, 2} // 111-2

// 目标已经在邀请队列
var ERR_GUILD_INVATE_OTHER_FAIL_INVATED = pbutil.StaticBuffer{3, 9, 111, 5} // 111-5

// 邀请列表已满，取消之前申请
var ERR_GUILD_INVATE_OTHER_FAIL_FULL = pbutil.StaticBuffer{3, 9, 111, 6} // 111-6

// 邀请的玩家已经在自己的联盟中
var ERR_GUILD_INVATE_OTHER_FAIL_GUILD_MEMBER = pbutil.StaticBuffer{3, 9, 111, 3} // 111-3

// 服务器忙，请稍后再试
var ERR_GUILD_INVATE_OTHER_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 111, 4} // 111-4

// 联盟已经满员
var ERR_GUILD_INVATE_OTHER_FAIL_FULL_MEMBER = pbutil.StaticBuffer{3, 9, 111, 10} // 111-10

func NewS2cGuildCancelInvateOtherMsg(id []byte) pbutil.Buffer {
	msg := &S2CGuildCancelInvateOtherProto{
		Id: id,
	}
	return NewS2cGuildCancelInvateOtherProtoMsg(msg)
}

var s2c_guild_cancel_invate_other = [...]byte{9, 113} // 113
func NewS2cGuildCancelInvateOtherProtoMsg(object *S2CGuildCancelInvateOtherProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_guild_cancel_invate_other[:], "s2c_guild_cancel_invate_other")

}

// 无效的玩家id
var ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 9, 114, 7} // 114-7

// 不在联盟中
var ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 114, 6} // 114-6

// 没有权限
var ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_DENY = pbutil.StaticBuffer{3, 9, 114, 2} // 114-2

// 邀请的玩家已经在自己的联盟中
var ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_GUILD_MEMBER = pbutil.StaticBuffer{3, 9, 114, 3} // 114-3

// 玩家不在邀请列表中
var ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_ID_NOT_EXIST = pbutil.StaticBuffer{3, 9, 114, 4} // 114-4

// 服务器忙，请稍后再试
var ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 114, 5} // 114-5

func NewS2cUserAddBeenInvateGuildMsg(id int32) pbutil.Buffer {
	msg := &S2CUserAddBeenInvateGuildProto{
		Id: id,
	}
	return NewS2cUserAddBeenInvateGuildProtoMsg(msg)
}

var s2c_user_add_been_invate_guild = [...]byte{9, 120} // 120
func NewS2cUserAddBeenInvateGuildProtoMsg(object *S2CUserAddBeenInvateGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_user_add_been_invate_guild[:], "s2c_user_add_been_invate_guild")

}

func NewS2cUserRemoveBeenInvateGuildMsg(id int32) pbutil.Buffer {
	msg := &S2CUserRemoveBeenInvateGuildProto{
		Id: id,
	}
	return NewS2cUserRemoveBeenInvateGuildProtoMsg(msg)
}

var s2c_user_remove_been_invate_guild = [...]byte{9, 121} // 121
func NewS2cUserRemoveBeenInvateGuildProtoMsg(object *S2CUserRemoveBeenInvateGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_user_remove_been_invate_guild[:], "s2c_user_remove_been_invate_guild")

}

func NewS2cUserReplyInvateRequestMsg(id int32, agree bool) pbutil.Buffer {
	msg := &S2CUserReplyInvateRequestProto{
		Id:    id,
		Agree: agree,
	}
	return NewS2cUserReplyInvateRequestProtoMsg(msg)
}

var s2c_user_reply_invate_request = [...]byte{9, 49} // 49
func NewS2cUserReplyInvateRequestProtoMsg(object *S2CUserReplyInvateRequestProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_user_reply_invate_request[:], "s2c_user_reply_invate_request")

}

// 无效的id
var ERR_USER_REPLY_INVATE_REQUEST_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 9, 50, 2} // 50-2

// 你是盟主，请先卸任盟主再接受邀请
var ERR_USER_REPLY_INVATE_REQUEST_FAIL_LEADER = pbutil.StaticBuffer{3, 9, 50, 4} // 50-4

// 匈奴入侵防守队员，且当前活动已开启
var ERR_USER_REPLY_INVATE_REQUEST_FAIL_XIONG_NU_DEFENDER = pbutil.StaticBuffer{3, 9, 50, 5} // 50-5

// 离开此联盟不足4小时，不能加入
var ERR_USER_REPLY_INVATE_REQUEST_FAIL_LEAVE_CD = pbutil.StaticBuffer{3, 9, 50, 6} // 50-6

// 服务器忙，请稍后再试
var ERR_USER_REPLY_INVATE_REQUEST_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 50, 3} // 50-3

// 联盟已满员
var ERR_USER_REPLY_INVATE_REQUEST_FAIL_FULL_MEMBER = pbutil.StaticBuffer{3, 9, 50, 7} // 50-7

func NewS2cListInviteMeGuildMsg(guild_list []*shared_proto.GuildSnapshotProto) pbutil.Buffer {
	msg := &S2CListInviteMeGuildProto{
		GuildList: guild_list,
	}
	return NewS2cListInviteMeGuildProtoMsg(msg)
}

var s2c_list_invite_me_guild = [...]byte{9, 194, 1} // 194
func NewS2cListInviteMeGuildProtoMsg(object *S2CListInviteMeGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_invite_me_guild[:], "s2c_list_invite_me_guild")

}

// 服务器忙，请稍后再试
var ERR_LIST_INVITE_ME_GUILD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 195, 1, 2} // 195-2

func NewS2cUpdateFriendGuildMsg(text string) pbutil.Buffer {
	msg := &S2CUpdateFriendGuildProto{
		Text: text,
	}
	return NewS2cUpdateFriendGuildProtoMsg(msg)
}

var s2c_update_friend_guild = [...]byte{9, 126} // 126
func NewS2cUpdateFriendGuildProtoMsg(object *S2CUpdateFriendGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_friend_guild[:], "s2c_update_friend_guild")

}

// 内容太长
var ERR_UPDATE_FRIEND_GUILD_FAIL_TEXT_TOO_LONG = pbutil.StaticBuffer{3, 9, 127, 1} // 127-1

// 你没有联盟
var ERR_UPDATE_FRIEND_GUILD_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 9, 127, 2} // 127-2

// 你没有权限操作
var ERR_UPDATE_FRIEND_GUILD_FAIL_DENY = pbutil.StaticBuffer{3, 9, 127, 3} // 127-3

// Npc联盟不允许操作
var ERR_UPDATE_FRIEND_GUILD_FAIL_NPC = pbutil.StaticBuffer{3, 9, 127, 4} // 127-4

// 输入包含敏感词
var ERR_UPDATE_FRIEND_GUILD_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 9, 127, 6} // 127-6

// 服务器忙，请稍后再试
var ERR_UPDATE_FRIEND_GUILD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 9, 127, 5} // 127-5

func NewS2cUpdateEnemyGuildMsg(text string) pbutil.Buffer {
	msg := &S2CUpdateEnemyGuildProto{
		Text: text,
	}
	return NewS2cUpdateEnemyGuildProtoMsg(msg)
}

var s2c_update_enemy_guild = [...]byte{9, 129, 1} // 129
func NewS2cUpdateEnemyGuildProtoMsg(object *S2CUpdateEnemyGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_enemy_guild[:], "s2c_update_enemy_guild")

}

// 内容太长
var ERR_UPDATE_ENEMY_GUILD_FAIL_TEXT_TOO_LONG = pbutil.StaticBuffer{4, 9, 130, 1, 1} // 130-1

// 你没有联盟
var ERR_UPDATE_ENEMY_GUILD_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 130, 1, 2} // 130-2

// 你没有权限操作
var ERR_UPDATE_ENEMY_GUILD_FAIL_DENY = pbutil.StaticBuffer{4, 9, 130, 1, 3} // 130-3

// Npc联盟不允许操作
var ERR_UPDATE_ENEMY_GUILD_FAIL_NPC = pbutil.StaticBuffer{4, 9, 130, 1, 4} // 130-4

// 输入包含敏感词
var ERR_UPDATE_ENEMY_GUILD_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{4, 9, 130, 1, 6} // 130-6

// 服务器忙，请稍后再试
var ERR_UPDATE_ENEMY_GUILD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 130, 1, 5} // 130-5

func NewS2cUpdateGuildPrestigeMsg(old_target int32, target int32) pbutil.Buffer {
	msg := &S2CUpdateGuildPrestigeProto{
		OldTarget: old_target,
		Target:    target,
	}
	return NewS2cUpdateGuildPrestigeProtoMsg(msg)
}

var s2c_update_guild_prestige = [...]byte{9, 132, 1} // 132
func NewS2cUpdateGuildPrestigeProtoMsg(object *S2CUpdateGuildPrestigeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_guild_prestige[:], "s2c_update_guild_prestige")

}

// 声望目标没找到
var ERR_UPDATE_GUILD_PRESTIGE_FAIL_TARGET_NOT_FOUND = pbutil.StaticBuffer{4, 9, 133, 1, 1} // 133-1

// 你没有联盟
var ERR_UPDATE_GUILD_PRESTIGE_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 133, 1, 2} // 133-2

// 你没有权限操作
var ERR_UPDATE_GUILD_PRESTIGE_FAIL_DENY = pbutil.StaticBuffer{4, 9, 133, 1, 3} // 133-3

// Npc联盟不允许操作
var ERR_UPDATE_GUILD_PRESTIGE_FAIL_NPC = pbutil.StaticBuffer{4, 9, 133, 1, 4} // 133-4

// 消耗不足
var ERR_UPDATE_GUILD_PRESTIGE_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 9, 133, 1, 6} // 133-6

// 修改倒计时
var ERR_UPDATE_GUILD_PRESTIGE_FAIL_COUNTDOWN = pbutil.StaticBuffer{4, 9, 133, 1, 7} // 133-7

// 修改的目标跟现在的目标一样
var ERR_UPDATE_GUILD_PRESTIGE_FAIL_SAME_TARGET = pbutil.StaticBuffer{4, 9, 133, 1, 8} // 133-8

// 服务器忙，请稍后再试
var ERR_UPDATE_GUILD_PRESTIGE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 133, 1, 5} // 133-5

var PLACE_GUILD_STATUE_S2C = pbutil.StaticBuffer{3, 9, 135, 1} // 135

// 没有联盟
var ERR_PLACE_GUILD_STATUE_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 9, 136, 1, 1} // 136-1

// 不是盟主
var ERR_PLACE_GUILD_STATUE_FAIL_NOT_LEADER = pbutil.StaticBuffer{4, 9, 136, 1, 2} // 136-2

// 有放置了，请先取回
var ERR_PLACE_GUILD_STATUE_FAIL_HAS_PLACED = pbutil.StaticBuffer{4, 9, 136, 1, 3} // 136-3

// 地图没找到
var ERR_PLACE_GUILD_STATUE_FAIL_MAP_NOT_FOUND = pbutil.StaticBuffer{4, 9, 136, 1, 4} // 136-4

// x非法
var ERR_PLACE_GUILD_STATUE_FAIL_X_INVALID = pbutil.StaticBuffer{4, 9, 136, 1, 5} // 136-5

// y非法
var ERR_PLACE_GUILD_STATUE_FAIL_Y_INVALID = pbutil.StaticBuffer{4, 9, 136, 1, 6} // 136-6

// 服务器忙，请稍后再试
var ERR_PLACE_GUILD_STATUE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 136, 1, 7} // 136-7

func NewS2cGuildStatueMsg(realm_id int32) pbutil.Buffer {
	msg := &S2CGuildStatueProto{
		RealmId: realm_id,
	}
	return NewS2cGuildStatueProtoMsg(msg)
}

var s2c_guild_statue = [...]byte{9, 137, 1} // 137
func NewS2cGuildStatueProtoMsg(object *S2CGuildStatueProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_guild_statue[:], "s2c_guild_statue")

}

var TAKE_BACK_GUILD_STATUE_S2C = pbutil.StaticBuffer{3, 9, 139, 1} // 139

var TAKE_BACK_GUILD_STATUE_S2C_BROADCAST = pbutil.StaticBuffer{3, 9, 140, 1} // 140

// 没有联盟
var ERR_TAKE_BACK_GUILD_STATUE_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 9, 141, 1, 1} // 141-1

// 不是盟主
var ERR_TAKE_BACK_GUILD_STATUE_FAIL_NOT_LEADER = pbutil.StaticBuffer{4, 9, 141, 1, 2} // 141-2

// 没有放置
var ERR_TAKE_BACK_GUILD_STATUE_FAIL_NOT_PLACE = pbutil.StaticBuffer{4, 9, 141, 1, 3} // 141-3

// 服务器忙，请稍后再试
var ERR_TAKE_BACK_GUILD_STATUE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 141, 1, 4} // 141-4

var COLLECT_FIRST_JOIN_GUILD_PRIZE_S2C = pbutil.StaticBuffer{3, 9, 144, 1} // 144

// 没有加入联盟
var ERR_COLLECT_FIRST_JOIN_GUILD_PRIZE_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 9, 145, 1, 3} // 145-3

// 奖励已经被领取了
var ERR_COLLECT_FIRST_JOIN_GUILD_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{4, 9, 145, 1, 2} // 145-2

func NewS2cUpdateSeekHelpMsg(help_type int32, worker_pos int32, enable bool) pbutil.Buffer {
	msg := &S2CUpdateSeekHelpProto{
		HelpType:  help_type,
		WorkerPos: worker_pos,
		Enable:    enable,
	}
	return NewS2cUpdateSeekHelpProtoMsg(msg)
}

var s2c_update_seek_help = [...]byte{9, 146, 1} // 146
func NewS2cUpdateSeekHelpProtoMsg(object *S2CUpdateSeekHelpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_seek_help[:], "s2c_update_seek_help")

}

func NewS2cUpdateHelpMemberTimesMsg(daily_help_member_times int32) pbutil.Buffer {
	msg := &S2CUpdateHelpMemberTimesProto{
		DailyHelpMemberTimes: daily_help_member_times,
	}
	return NewS2cUpdateHelpMemberTimesProtoMsg(msg)
}

var s2c_update_help_member_times = [...]byte{9, 157, 1} // 157
func NewS2cUpdateHelpMemberTimesProtoMsg(object *S2CUpdateHelpMemberTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_help_member_times[:], "s2c_update_help_member_times")

}

func NewS2cSeekHelpMsg(help_type int32, worker_pos int32) pbutil.Buffer {
	msg := &S2CSeekHelpProto{
		HelpType:  help_type,
		WorkerPos: worker_pos,
	}
	return NewS2cSeekHelpProtoMsg(msg)
}

var s2c_seek_help = [...]byte{9, 148, 1} // 148
func NewS2cSeekHelpProtoMsg(object *S2CSeekHelpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_seek_help[:], "s2c_seek_help")

}

// 当前求助状态不可用
var ERR_SEEK_HELP_FAIL_DISABLE = pbutil.StaticBuffer{4, 9, 149, 1, 1} // 149-1

// 自己不在联盟中，不能求助
var ERR_SEEK_HELP_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 149, 1, 2} // 149-2

// 自己还没有外使院建筑
var ERR_SEEK_HELP_FAIL_WAI_SHI_YUAN = pbutil.StaticBuffer{4, 9, 149, 1, 3} // 149-3

func NewS2cAddGuildSeekHelpMsg(data []byte) pbutil.Buffer {
	msg := &S2CAddGuildSeekHelpProto{
		Data: data,
	}
	return NewS2cAddGuildSeekHelpProtoMsg(msg)
}

func NewS2cAddGuildSeekHelpMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CAddGuildSeekHelpProto{
		Data: safeMarshal(data),
	}
	return NewS2cAddGuildSeekHelpProtoMsg(msg)
}

var s2c_add_guild_seek_help = [...]byte{9, 150, 1} // 150
func NewS2cAddGuildSeekHelpProtoMsg(object *S2CAddGuildSeekHelpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_guild_seek_help[:], "s2c_add_guild_seek_help")

}

func NewS2cHelpGuildMemberMsg(id string) pbutil.Buffer {
	msg := &S2CHelpGuildMemberProto{
		Id: id,
	}
	return NewS2cHelpGuildMemberProtoMsg(msg)
}

var s2c_help_guild_member = [...]byte{9, 152, 1} // 152
func NewS2cHelpGuildMemberProtoMsg(object *S2CHelpGuildMemberProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_help_guild_member[:], "s2c_help_guild_member")

}

// 求助id没找到
var ERR_HELP_GUILD_MEMBER_FAIL_ID_NOT_FOUND = pbutil.StaticBuffer{4, 9, 153, 1, 1} // 153-1

// 这条求助你已经帮助过了
var ERR_HELP_GUILD_MEMBER_FAIL_HELPED = pbutil.StaticBuffer{4, 9, 153, 1, 2} // 153-2

// 你不在联盟中，不能帮助盟友的求助
var ERR_HELP_GUILD_MEMBER_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 153, 1, 3} // 153-3

func NewS2cHelpAllGuildMemberMsg(ids []string) pbutil.Buffer {
	msg := &S2CHelpAllGuildMemberProto{
		Ids: ids,
	}
	return NewS2cHelpAllGuildMemberProtoMsg(msg)
}

var s2c_help_all_guild_member = [...]byte{9, 159, 1} // 159
func NewS2cHelpAllGuildMemberProtoMsg(object *S2CHelpAllGuildMemberProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_help_all_guild_member[:], "s2c_help_all_guild_member")

}

// 你不在联盟中，不能帮助盟友的求助
var ERR_HELP_ALL_GUILD_MEMBER_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 160, 1, 1} // 160-1

func NewS2cAddGuildSeekHelpHeroIdsMsg(id string, help_hero_id []byte) pbutil.Buffer {
	msg := &S2CAddGuildSeekHelpHeroIdsProto{
		Id:         id,
		HelpHeroId: help_hero_id,
	}
	return NewS2cAddGuildSeekHelpHeroIdsProtoMsg(msg)
}

var s2c_add_guild_seek_help_hero_ids = [...]byte{9, 154, 1} // 154
func NewS2cAddGuildSeekHelpHeroIdsProtoMsg(object *S2CAddGuildSeekHelpHeroIdsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_guild_seek_help_hero_ids[:], "s2c_add_guild_seek_help_hero_ids")

}

func NewS2cRemoveGuildSeekHelpMsg(id string) pbutil.Buffer {
	msg := &S2CRemoveGuildSeekHelpProto{
		Id: id,
	}
	return NewS2cRemoveGuildSeekHelpProtoMsg(msg)
}

var s2c_remove_guild_seek_help = [...]byte{9, 155, 1} // 155
func NewS2cRemoveGuildSeekHelpProtoMsg(object *S2CRemoveGuildSeekHelpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_guild_seek_help[:], "s2c_remove_guild_seek_help")

}

func NewS2cListGuildSeekHelpMsg(data [][]byte) pbutil.Buffer {
	msg := &S2CListGuildSeekHelpProto{
		Data: data,
	}
	return NewS2cListGuildSeekHelpProtoMsg(msg)
}

func NewS2cListGuildSeekHelpMarshalMsg(data [][]byte) pbutil.Buffer {
	msg := &S2CListGuildSeekHelpProto{
		Data: data,
	}
	return NewS2cListGuildSeekHelpProtoMsg(msg)
}

var s2c_list_guild_seek_help = [...]byte{9, 156, 1} // 156
func NewS2cListGuildSeekHelpProtoMsg(object *S2CListGuildSeekHelpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_guild_seek_help[:], "s2c_list_guild_seek_help")

}

func NewS2cListGuildEventPrizeMsg(id []int32, data_id []int32, expire_time []int32, hero_id [][]byte, hero_name []string) pbutil.Buffer {
	msg := &S2CListGuildEventPrizeProto{
		Id:         id,
		DataId:     data_id,
		ExpireTime: expire_time,
		HeroId:     hero_id,
		HeroName:   hero_name,
	}
	return NewS2cListGuildEventPrizeProtoMsg(msg)
}

var s2c_list_guild_event_prize = [...]byte{9, 161, 1} // 161
func NewS2cListGuildEventPrizeProtoMsg(object *S2CListGuildEventPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_guild_event_prize[:], "s2c_list_guild_event_prize")

}

func NewS2cAddGuildEventPrizeMsg(id int32, data_id int32, expire_time int32, hero_id []byte, hero_name string) pbutil.Buffer {
	msg := &S2CAddGuildEventPrizeProto{
		Id:         id,
		DataId:     data_id,
		ExpireTime: expire_time,
		HeroId:     hero_id,
		HeroName:   hero_name,
	}
	return NewS2cAddGuildEventPrizeProtoMsg(msg)
}

var s2c_add_guild_event_prize = [...]byte{9, 162, 1} // 162
func NewS2cAddGuildEventPrizeProtoMsg(object *S2CAddGuildEventPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_guild_event_prize[:], "s2c_add_guild_event_prize")

}

func NewS2cRemoveGuildEventPrizeMsg(id int32) pbutil.Buffer {
	msg := &S2CRemoveGuildEventPrizeProto{
		Id: id,
	}
	return NewS2cRemoveGuildEventPrizeProtoMsg(msg)
}

var s2c_remove_guild_event_prize = [...]byte{9, 170, 1} // 170
func NewS2cRemoveGuildEventPrizeProtoMsg(object *S2CRemoveGuildEventPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_guild_event_prize[:], "s2c_remove_guild_event_prize")

}

func NewS2cCollectGuildEventPrizeMsg(id int32, new_energy int32, add_energy int32, prize []byte) pbutil.Buffer {
	msg := &S2CCollectGuildEventPrizeProto{
		Id:        id,
		NewEnergy: new_energy,
		AddEnergy: add_energy,
		Prize:     prize,
	}
	return NewS2cCollectGuildEventPrizeProtoMsg(msg)
}

var s2c_collect_guild_event_prize = [...]byte{9, 164, 1} // 164
func NewS2cCollectGuildEventPrizeProtoMsg(object *S2CCollectGuildEventPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_guild_event_prize[:], "s2c_collect_guild_event_prize")

}

// id不存在
var ERR_COLLECT_GUILD_EVENT_PRIZE_FAIL_ID_NOT_FOUND = pbutil.StaticBuffer{4, 9, 165, 1, 1} // 165-1

// 你不在联盟中
var ERR_COLLECT_GUILD_EVENT_PRIZE_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 165, 1, 2} // 165-2

// vip等级不够，不能一键领
var ERR_COLLECT_GUILD_EVENT_PRIZE_FAIL_VIP_LIMIT = pbutil.StaticBuffer{4, 9, 165, 1, 3} // 165-3

func NewS2cUpdateFullBigBoxMsg(box_id int32, collectable bool, energy int32) pbutil.Buffer {
	msg := &S2CUpdateFullBigBoxProto{
		BoxId:       box_id,
		Collectable: collectable,
		Energy:      energy,
	}
	return NewS2cUpdateFullBigBoxProtoMsg(msg)
}

var s2c_update_full_big_box = [...]byte{9, 166, 1} // 166
func NewS2cUpdateFullBigBoxProtoMsg(object *S2CUpdateFullBigBoxProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_full_big_box[:], "s2c_update_full_big_box")

}

var COLLECT_FULL_BIG_BOX_S2C = pbutil.StaticBuffer{3, 9, 168, 1} // 168

// 你不在联盟中
var ERR_COLLECT_FULL_BIG_BOX_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 169, 1, 1} // 169-1

// 宝箱尚未解锁
var ERR_COLLECT_FULL_BIG_BOX_FAIL_LOCKED = pbutil.StaticBuffer{4, 9, 169, 1, 2} // 169-2

// 进入联盟时间不足，不能领取
var ERR_COLLECT_FULL_BIG_BOX_FAIL_TIME_NOT_ENOUGH = pbutil.StaticBuffer{4, 9, 169, 1, 4} // 169-4

// 服务器忙，请稍后再试
var ERR_COLLECT_FULL_BIG_BOX_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 169, 1, 3} // 169-3

func NewS2cUpdateHeroJoinGuildTimeMsg(join_guild_time int32) pbutil.Buffer {
	msg := &S2CUpdateHeroJoinGuildTimeProto{
		JoinGuildTime: join_guild_time,
	}
	return NewS2cUpdateHeroJoinGuildTimeProtoMsg(msg)
}

var s2c_update_hero_join_guild_time = [...]byte{9, 171, 1} // 171
func NewS2cUpdateHeroJoinGuildTimeProtoMsg(object *S2CUpdateHeroJoinGuildTimeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_hero_join_guild_time[:], "s2c_update_hero_join_guild_time")

}

func NewS2cUpgradeTechnologyMsg(group int32, end_time int32) pbutil.Buffer {
	msg := &S2CUpgradeTechnologyProto{
		Group:   group,
		EndTime: end_time,
	}
	return NewS2cUpgradeTechnologyProtoMsg(msg)
}

var s2c_upgrade_technology = [...]byte{9, 173, 1} // 173
func NewS2cUpgradeTechnologyProtoMsg(object *S2CUpgradeTechnologyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_upgrade_technology[:], "s2c_upgrade_technology")

}

// 你没有联盟
var ERR_UPGRADE_TECHNOLOGY_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 174, 1, 7} // 174-7

// 无效的科技组
var ERR_UPGRADE_TECHNOLOGY_FAIL_INVALID_GROUP = pbutil.StaticBuffer{4, 9, 174, 1, 1} // 174-1

// 科技已经达到最大等级
var ERR_UPGRADE_TECHNOLOGY_FAIL_MAX_LEVEL = pbutil.StaticBuffer{4, 9, 174, 1, 2} // 174-2

// 联盟等级不足
var ERR_UPGRADE_TECHNOLOGY_FAIL_REQUIRED = pbutil.StaticBuffer{4, 9, 174, 1, 3} // 174-3

// 联盟建设值不足
var ERR_UPGRADE_TECHNOLOGY_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 9, 174, 1, 4} // 174-4

// 存在正在升级的科技
var ERR_UPGRADE_TECHNOLOGY_FAIL_UPGRADING = pbutil.StaticBuffer{4, 9, 174, 1, 5} // 174-5

// 你没有权限操作
var ERR_UPGRADE_TECHNOLOGY_FAIL_DENY = pbutil.StaticBuffer{4, 9, 174, 1, 6} // 174-6

// 服务器忙，请稍后再试
var ERR_UPGRADE_TECHNOLOGY_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 174, 1, 8} // 174-8

var REDUCE_TECHNOLOGY_CD_S2C = pbutil.StaticBuffer{3, 9, 176, 1} // 176

// 你没有联盟
var ERR_REDUCE_TECHNOLOGY_CD_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 177, 1, 1} // 177-1

// 你没有权限操作
var ERR_REDUCE_TECHNOLOGY_CD_FAIL_DENY = pbutil.StaticBuffer{4, 9, 177, 1, 2} // 177-2

// 科技没有在升级，不能加速
var ERR_REDUCE_TECHNOLOGY_CD_FAIL_NO_UPGRADING = pbutil.StaticBuffer{4, 9, 177, 1, 3} // 177-3

// 帮派已经达到最大加速次数
var ERR_REDUCE_TECHNOLOGY_CD_FAIL_MAX_TIMES = pbutil.StaticBuffer{4, 9, 177, 1, 4} // 177-4

// 建设值不足
var ERR_REDUCE_TECHNOLOGY_CD_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 9, 177, 1, 5} // 177-5

// Npc联盟不允许操作
var ERR_REDUCE_TECHNOLOGY_CD_FAIL_NPC = pbutil.StaticBuffer{4, 9, 177, 1, 6} // 177-6

// 服务器忙，请稍后再试
var ERR_REDUCE_TECHNOLOGY_CD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 177, 1, 7} // 177-7

func NewS2cListGuildLogsMsg(data [][]byte) pbutil.Buffer {
	msg := &S2CListGuildLogsProto{
		Data: data,
	}
	return NewS2cListGuildLogsProtoMsg(msg)
}

func NewS2cListGuildLogsMarshalMsg(data [][]byte) pbutil.Buffer {
	msg := &S2CListGuildLogsProto{
		Data: data,
	}
	return NewS2cListGuildLogsProtoMsg(msg)
}

var s2c_list_guild_logs = [...]byte{9, 179, 1} // 179
func NewS2cListGuildLogsProtoMsg(object *S2CListGuildLogsProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_list_guild_logs[:], "s2c_list_guild_logs")

}

func NewS2cAddGuildLogMsg(data []byte) pbutil.Buffer {
	msg := &S2CAddGuildLogProto{
		Data: data,
	}
	return NewS2cAddGuildLogProtoMsg(msg)
}

func NewS2cAddGuildLogMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CAddGuildLogProto{
		Data: safeMarshal(data),
	}
	return NewS2cAddGuildLogProtoMsg(msg)
}

var s2c_add_guild_log = [...]byte{9, 180, 1} // 180
func NewS2cAddGuildLogProtoMsg(object *S2CAddGuildLogProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_guild_log[:], "s2c_add_guild_log")

}

func NewS2cRequestRecommendGuildMsg(has bool, next_notify_guild_time int32, id int32, name string, flag string, country int32) pbutil.Buffer {
	msg := &S2CRequestRecommendGuildProto{
		Has:                 has,
		NextNotifyGuildTime: next_notify_guild_time,
		Id:                  id,
		Name:                name,
		Flag:                flag,
		Country:             country,
	}
	return NewS2cRequestRecommendGuildProtoMsg(msg)
}

var s2c_request_recommend_guild = [...]byte{9, 182, 1} // 182
func NewS2cRequestRecommendGuildProtoMsg(object *S2CRequestRecommendGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_recommend_guild[:], "s2c_request_recommend_guild")

}

func NewS2cPushTechHelpableMsg(helpable bool) pbutil.Buffer {
	msg := &S2CPushTechHelpableProto{
		Helpable: helpable,
	}
	return NewS2cPushTechHelpableProtoMsg(msg)
}

var s2c_push_tech_helpable = [...]byte{9, 183, 1} // 183
func NewS2cPushTechHelpableProtoMsg(object *S2CPushTechHelpableProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_push_tech_helpable[:], "s2c_push_tech_helpable")

}

func NewS2cHelpTechMsg(tech int32) pbutil.Buffer {
	msg := &S2CHelpTechProto{
		Tech: tech,
	}
	return NewS2cHelpTechProtoMsg(msg)
}

var s2c_help_tech = [...]byte{9, 185, 1} // 185
func NewS2cHelpTechProtoMsg(object *S2CHelpTechProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_help_tech[:], "s2c_help_tech")

}

// 你没有联盟
var ERR_HELP_TECH_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 186, 1, 1} // 186-1

// 你不能协助
var ERR_HELP_TECH_FAIL_CANT_HELP = pbutil.StaticBuffer{4, 9, 186, 1, 2} // 186-2

// 没有升级中的科技
var ERR_HELP_TECH_FAIL_NO_TECH_UPGRADING = pbutil.StaticBuffer{4, 9, 186, 1, 3} // 186-3

// 服务器忙，请稍后再试
var ERR_HELP_TECH_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 186, 1, 4} // 186-4

func NewS2cRecommendInviteHerosMsg(heros []byte) pbutil.Buffer {
	msg := &S2CRecommendInviteHerosProto{
		Heros: heros,
	}
	return NewS2cRecommendInviteHerosProtoMsg(msg)
}

var s2c_recommend_invite_heros = [...]byte{9, 188, 1} // 188
func NewS2cRecommendInviteHerosProtoMsg(object *S2CRecommendInviteHerosProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_recommend_invite_heros[:], "s2c_recommend_invite_heros")

}

// 服务器忙，请稍后再试
var ERR_RECOMMEND_INVITE_HEROS_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 189, 1, 1} // 189-1

func NewS2cSearchNoGuildHerosMsg(heros []byte) pbutil.Buffer {
	msg := &S2CSearchNoGuildHerosProto{
		Heros: heros,
	}
	return NewS2cSearchNoGuildHerosProtoMsg(msg)
}

var s2c_search_no_guild_heros = [...]byte{9, 191, 1} // 191
func NewS2cSearchNoGuildHerosProtoMsg(object *S2CSearchNoGuildHerosProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_search_no_guild_heros[:], "s2c_search_no_guild_heros")

}

// 参数错误
var ERR_SEARCH_NO_GUILD_HEROS_FAIL_INVALID_ARG = pbutil.StaticBuffer{4, 9, 192, 1, 2} // 192-2

// 服务器忙，请稍后再试
var ERR_SEARCH_NO_GUILD_HEROS_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 192, 1, 1} // 192-1

// 请求太频繁
var ERR_SEARCH_NO_GUILD_HEROS_FAIL_TOO_FAST = pbutil.StaticBuffer{4, 9, 192, 1, 3} // 192-3

func NewS2cViewMcWarRecordMsg(record *shared_proto.McWarAllRecordProto, record2 *shared_proto.McWarAllRecordWithJoinedProto) pbutil.Buffer {
	msg := &S2CViewMcWarRecordProto{
		Record:  record,
		Record2: record2,
	}
	return NewS2cViewMcWarRecordProtoMsg(msg)
}

var s2c_view_mc_war_record = [...]byte{9, 200, 1} // 200
func NewS2cViewMcWarRecordProtoMsg(object *S2CViewMcWarRecordProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_view_mc_war_record[:], "s2c_view_mc_war_record")

}

// 没有联盟
var ERR_VIEW_MC_WAR_RECORD_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 9, 201, 1, 1} // 201-1

func NewS2cUpdateGuildMarkMsg(mark *shared_proto.GuildMarkProto) pbutil.Buffer {
	msg := &S2CUpdateGuildMarkProto{
		Mark: mark,
	}
	return NewS2cUpdateGuildMarkProtoMsg(msg)
}

var s2c_update_guild_mark = [...]byte{9, 197, 1} // 197
func NewS2cUpdateGuildMarkProtoMsg(object *S2CUpdateGuildMarkProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_guild_mark[:], "s2c_update_guild_mark")

}

// 无效的序号
var ERR_UPDATE_GUILD_MARK_FAIL_INVALID_INDEX = pbutil.StaticBuffer{4, 9, 198, 1, 1} // 198-1

// 无效的坐标
var ERR_UPDATE_GUILD_MARK_FAIL_INVALID_POS = pbutil.StaticBuffer{4, 9, 198, 1, 2} // 198-2

// 无效的标记内容
var ERR_UPDATE_GUILD_MARK_FAIL_INVALID_MSG = pbutil.StaticBuffer{4, 9, 198, 1, 3} // 198-3

// 标记内容包含敏感词
var ERR_UPDATE_GUILD_MARK_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{4, 9, 198, 1, 4} // 198-4

// 你没有联盟
var ERR_UPDATE_GUILD_MARK_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 198, 1, 5} // 198-5

// 你没有权限操作
var ERR_UPDATE_GUILD_MARK_FAIL_DENY = pbutil.StaticBuffer{4, 9, 198, 1, 6} // 198-6

// 服务器忙，请稍后再试
var ERR_UPDATE_GUILD_MARK_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 198, 1, 7} // 198-7

func NewS2cViewYinliangRecordMsg(record *shared_proto.GuildAllYinliangRecordProto) pbutil.Buffer {
	msg := &S2CViewYinliangRecordProto{
		Record: record,
	}
	return NewS2cViewYinliangRecordProtoMsg(msg)
}

var s2c_view_yinliang_record = [...]byte{9, 203, 1} // 203
func NewS2cViewYinliangRecordProtoMsg(object *S2CViewYinliangRecordProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_view_yinliang_record[:], "s2c_view_yinliang_record")

}

// 没有联盟
var ERR_VIEW_YINLIANG_RECORD_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 9, 204, 1, 1} // 204-1

func NewS2cSendYinliangToOtherGuildMsg(gid int32, amount int32) pbutil.Buffer {
	msg := &S2CSendYinliangToOtherGuildProto{
		Gid:    gid,
		Amount: amount,
	}
	return NewS2cSendYinliangToOtherGuildProtoMsg(msg)
}

var s2c_send_yinliang_to_other_guild = [...]byte{9, 206, 1} // 206
func NewS2cSendYinliangToOtherGuildProtoMsg(object *S2CSendYinliangToOtherGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_send_yinliang_to_other_guild[:], "s2c_send_yinliang_to_other_guild")

}

// 没有权限
var ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_DENY = pbutil.StaticBuffer{4, 9, 207, 1, 1} // 207-1

// 对方联盟不存在
var ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 9, 207, 1, 2} // 207-2

// 钱不够
var ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_NOT_ENOUGH = pbutil.StaticBuffer{4, 9, 207, 1, 3} // 207-3

// 参数非法
var ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_INVALID_AMOUNT = pbutil.StaticBuffer{4, 9, 207, 1, 4} // 207-4

// 服务器错误
var ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_SERVER_ERR = pbutil.StaticBuffer{4, 9, 207, 1, 5} // 207-5

func NewS2cSendYinliangToMemberMsg(mem_id []byte, amount int32) pbutil.Buffer {
	msg := &S2CSendYinliangToMemberProto{
		MemId:  mem_id,
		Amount: amount,
	}
	return NewS2cSendYinliangToMemberProtoMsg(msg)
}

var s2c_send_yinliang_to_member = [...]byte{9, 209, 1} // 209
func NewS2cSendYinliangToMemberProtoMsg(object *S2CSendYinliangToMemberProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_send_yinliang_to_member[:], "s2c_send_yinliang_to_member")

}

// 没有权限
var ERR_SEND_YINLIANG_TO_MEMBER_FAIL_DENY = pbutil.StaticBuffer{4, 9, 210, 1, 1} // 210-1

// 成员不存在
var ERR_SEND_YINLIANG_TO_MEMBER_FAIL_NO_MEMBER = pbutil.StaticBuffer{4, 9, 210, 1, 2} // 210-2

// 钱不够
var ERR_SEND_YINLIANG_TO_MEMBER_FAIL_NOT_ENOUGH = pbutil.StaticBuffer{4, 9, 210, 1, 3} // 210-3

// 参数非法
var ERR_SEND_YINLIANG_TO_MEMBER_FAIL_INVALID_AMOUNT = pbutil.StaticBuffer{4, 9, 210, 1, 4} // 210-4

// 服务器错误
var ERR_SEND_YINLIANG_TO_MEMBER_FAIL_SERVER_ERR = pbutil.StaticBuffer{4, 9, 210, 1, 5} // 210-5

func NewS2cPaySalaryMsg(amount int32) pbutil.Buffer {
	msg := &S2CPaySalaryProto{
		Amount: amount,
	}
	return NewS2cPaySalaryProtoMsg(msg)
}

var s2c_pay_salary = [...]byte{9, 212, 1} // 212
func NewS2cPaySalaryProtoMsg(object *S2CPaySalaryProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_pay_salary[:], "s2c_pay_salary")

}

// 没有权限
var ERR_PAY_SALARY_FAIL_DENY = pbutil.StaticBuffer{4, 9, 213, 1, 1} // 213-1

// 钱不够
var ERR_PAY_SALARY_FAIL_NOT_ENOUGH = pbutil.StaticBuffer{4, 9, 213, 1, 2} // 213-2

// 参数非法
var ERR_PAY_SALARY_FAIL_INVALID_AMOUNT = pbutil.StaticBuffer{4, 9, 213, 1, 3} // 213-3

// 服务器错误
var ERR_PAY_SALARY_FAIL_SERVER_ERR = pbutil.StaticBuffer{4, 9, 213, 1, 4} // 213-4

func NewS2cSetSalaryMsg(mem_id []byte, salary int32) pbutil.Buffer {
	msg := &S2CSetSalaryProto{
		MemId:  mem_id,
		Salary: salary,
	}
	return NewS2cSetSalaryProtoMsg(msg)
}

var s2c_set_salary = [...]byte{9, 215, 1} // 215
func NewS2cSetSalaryProtoMsg(object *S2CSetSalaryProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_salary[:], "s2c_set_salary")

}

// 没有权限
var ERR_SET_SALARY_FAIL_DENY = pbutil.StaticBuffer{4, 9, 216, 1, 1} // 216-1

// 成员不存在
var ERR_SET_SALARY_FAIL_NO_MEMBER = pbutil.StaticBuffer{4, 9, 216, 1, 2} // 216-2

// 工资非法
var ERR_SET_SALARY_FAIL_INVALID_SALARY = pbutil.StaticBuffer{4, 9, 216, 1, 3} // 216-3

// 服务器错误
var ERR_SET_SALARY_FAIL_SERVER_ERR = pbutil.StaticBuffer{4, 9, 216, 1, 4} // 216-4

func NewS2cUpdateHeroGuildMsg(update_type int32, data *shared_proto.HeroGuildProto) pbutil.Buffer {
	msg := &S2CUpdateHeroGuildProto{
		UpdateType: update_type,
		Data:       data,
	}
	return NewS2cUpdateHeroGuildProtoMsg(msg)
}

var s2c_update_hero_guild = [...]byte{9, 217, 1} // 217
func NewS2cUpdateHeroGuildProtoMsg(object *S2CUpdateHeroGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_hero_guild[:], "s2c_update_hero_guild")

}

func NewS2cViewSendYinliangToGuildMsg(guilds *shared_proto.GuildAllYinliangSendToGuildProto) pbutil.Buffer {
	msg := &S2CViewSendYinliangToGuildProto{
		Guilds: guilds,
	}
	return NewS2cViewSendYinliangToGuildProtoMsg(msg)
}

var s2c_view_send_yinliang_to_guild = [...]byte{9, 219, 1} // 219
func NewS2cViewSendYinliangToGuildProtoMsg(object *S2CViewSendYinliangToGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_send_yinliang_to_guild[:], "s2c_view_send_yinliang_to_guild")

}

// 服务器错误
var ERR_VIEW_SEND_YINLIANG_TO_GUILD_FAIL_SERVER_ERR = pbutil.StaticBuffer{4, 9, 220, 1, 1} // 220-1

func NewS2cUpdateHufuMsg(hufu int32) pbutil.Buffer {
	msg := &S2CUpdateHufuProto{
		Hufu: hufu,
	}
	return NewS2cUpdateHufuProtoMsg(msg)
}

var s2c_update_hufu = [...]byte{9, 221, 1} // 221
func NewS2cUpdateHufuProtoMsg(object *S2CUpdateHufuProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_hufu[:], "s2c_update_hufu")

}

func NewS2cConveneMsg(target []byte) pbutil.Buffer {
	msg := &S2CConveneProto{
		Target: target,
	}
	return NewS2cConveneProtoMsg(msg)
}

var s2c_convene = [...]byte{9, 229, 1} // 229
func NewS2cConveneProtoMsg(object *S2CConveneProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_convene[:], "s2c_convene")

}

// 无效的目标id
var ERR_CONVENE_FAIL_INVALID_TARGET = pbutil.StaticBuffer{4, 9, 230, 1, 1} // 230-1

// 你没有联盟
var ERR_CONVENE_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 230, 1, 2} // 230-2

// 权限不足
var ERR_CONVENE_FAIL_DENY = pbutil.StaticBuffer{4, 9, 230, 1, 3} // 230-3

// 目标不在你的联盟中
var ERR_CONVENE_FAIL_TARGET_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 230, 1, 4} // 230-4

// 盟友召集CD中，请稍后再试
var ERR_CONVENE_FAIL_COOLDOWN = pbutil.StaticBuffer{4, 9, 230, 1, 5} // 230-5

// 服务器忙，请稍后再试
var ERR_CONVENE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 230, 1, 6} // 230-6

func NewS2cCollectDailyGuildRankPrizeMsg(prize []byte) pbutil.Buffer {
	msg := &S2CCollectDailyGuildRankPrizeProto{
		Prize: prize,
	}
	return NewS2cCollectDailyGuildRankPrizeProtoMsg(msg)
}

var s2c_collect_daily_guild_rank_prize = [...]byte{9, 232, 1} // 232
func NewS2cCollectDailyGuildRankPrizeProtoMsg(object *S2CCollectDailyGuildRankPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_daily_guild_rank_prize[:], "s2c_collect_daily_guild_rank_prize")

}

// 没有联盟
var ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 9, 233, 1, 1} // 233-1

// 联盟未上榜
var ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_NO_GUILD_RANK = pbutil.StaticBuffer{4, 9, 233, 1, 2} // 233-2

// 无法领取
var ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{4, 9, 233, 1, 3} // 233-3

// 服务器忙，请稍后再试
var ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 233, 1, 4} // 233-4

func NewS2cViewDailyGuildRankMsg(rank *shared_proto.RankProto) pbutil.Buffer {
	msg := &S2CViewDailyGuildRankProto{
		Rank: rank,
	}
	return NewS2cViewDailyGuildRankProtoMsg(msg)
}

var s2c_view_daily_guild_rank = [...]byte{9, 235, 1} // 235
func NewS2cViewDailyGuildRankProtoMsg(object *S2CViewDailyGuildRankProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_daily_guild_rank[:], "s2c_view_daily_guild_rank")

}

// 还未加入国家
var ERR_VIEW_DAILY_GUILD_RANK_FAIL_NO_COUNTRY = pbutil.StaticBuffer{4, 9, 236, 1, 1} // 236-1

// 服务器忙，请稍后再试
var ERR_VIEW_DAILY_GUILD_RANK_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 9, 236, 1, 2} // 236-2

func NewS2cViewLastGuildRankMsg(rank int32) pbutil.Buffer {
	msg := &S2CViewLastGuildRankProto{
		Rank: rank,
	}
	return NewS2cViewLastGuildRankProtoMsg(msg)
}

var s2c_view_last_guild_rank = [...]byte{9, 238, 1} // 238
func NewS2cViewLastGuildRankProtoMsg(object *S2CViewLastGuildRankProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_last_guild_rank[:], "s2c_view_last_guild_rank")

}

func NewS2cAddRecommendMcBuildMsg(new_mc_ids []int32) pbutil.Buffer {
	msg := &S2CAddRecommendMcBuildProto{
		NewMcIds: new_mc_ids,
	}
	return NewS2cAddRecommendMcBuildProtoMsg(msg)
}

var s2c_add_recommend_mc_build = [...]byte{9, 241, 1} // 241
func NewS2cAddRecommendMcBuildProtoMsg(object *S2CAddRecommendMcBuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_recommend_mc_build[:], "s2c_add_recommend_mc_build")

}

// 没有权限
var ERR_ADD_RECOMMEND_MC_BUILD_FAIL_DENY = pbutil.StaticBuffer{4, 9, 242, 1, 1} // 242-1

// 名城不存在
var ERR_ADD_RECOMMEND_MC_BUILD_FAIL_INVALID_MC_ID = pbutil.StaticBuffer{4, 9, 242, 1, 2} // 242-2

// 名城已经被推荐
var ERR_ADD_RECOMMEND_MC_BUILD_FAIL_MC_IS_RECOMMENDED = pbutil.StaticBuffer{4, 9, 242, 1, 3} // 242-3

// 没有联盟
var ERR_ADD_RECOMMEND_MC_BUILD_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 9, 242, 1, 5} // 242-5

// 服务器错误
var ERR_ADD_RECOMMEND_MC_BUILD_FAIL_SERVER_ERR = pbutil.StaticBuffer{4, 9, 242, 1, 4} // 242-4

func NewS2cViewTaskProgressMsg(version int32, progress []*shared_proto.Int32Pair) pbutil.Buffer {
	msg := &S2CViewTaskProgressProto{
		Version:  version,
		Progress: progress,
	}
	return NewS2cViewTaskProgressProtoMsg(msg)
}

var s2c_view_task_progress = [...]byte{9, 244, 1} // 244
func NewS2cViewTaskProgressProtoMsg(object *S2CViewTaskProgressProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_task_progress[:], "s2c_view_task_progress")

}

// 没有联盟
var ERR_VIEW_TASK_PROGRESS_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 9, 245, 1, 1} // 245-1

// 联盟等级不足
var ERR_VIEW_TASK_PROGRESS_FAIL_GUILD_LEVEL_LIMIT = pbutil.StaticBuffer{4, 9, 245, 1, 2} // 245-2

func NewS2cNoticeTaskStageUpdateMsg(task_id int32, stage int32) pbutil.Buffer {
	msg := &S2CNoticeTaskStageUpdateProto{
		TaskId: task_id,
		Stage:  stage,
	}
	return NewS2cNoticeTaskStageUpdateProtoMsg(msg)
}

var s2c_notice_task_stage_update = [...]byte{9, 246, 1} // 246
func NewS2cNoticeTaskStageUpdateProtoMsg(object *S2CNoticeTaskStageUpdateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_notice_task_stage_update[:], "s2c_notice_task_stage_update")

}

func NewS2cCollectTaskPrizeMsg(task_id int32, stage int32) pbutil.Buffer {
	msg := &S2CCollectTaskPrizeProto{
		TaskId: task_id,
		Stage:  stage,
	}
	return NewS2cCollectTaskPrizeProtoMsg(msg)
}

var s2c_collect_task_prize = [...]byte{9, 248, 1} // 248
func NewS2cCollectTaskPrizeProtoMsg(object *S2CCollectTaskPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_task_prize[:], "s2c_collect_task_prize")

}

// 无效的数据
var ERR_COLLECT_TASK_PRIZE_FAIL_INVALID_VALUE = pbutil.StaticBuffer{4, 9, 249, 1, 1} // 249-1

// 没有联盟
var ERR_COLLECT_TASK_PRIZE_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 9, 249, 1, 2} // 249-2

// 阶段奖励未激活
var ERR_COLLECT_TASK_PRIZE_FAIL_NO_PRIZE = pbutil.StaticBuffer{4, 9, 249, 1, 3} // 249-3

// 已领取
var ERR_COLLECT_TASK_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{4, 9, 249, 1, 4} // 249-4

func NewS2cGuildChangeCountryMsg(country int32) pbutil.Buffer {
	msg := &S2CGuildChangeCountryProto{
		Country: country,
	}
	return NewS2cGuildChangeCountryProtoMsg(msg)
}

var s2c_guild_change_country = [...]byte{9, 251, 1} // 251
func NewS2cGuildChangeCountryProtoMsg(object *S2CGuildChangeCountryProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_guild_change_country[:], "s2c_guild_change_country")

}

// 无效的国家
var ERR_GUILD_CHANGE_COUNTRY_FAIL_INVALID_COUNTRY = pbutil.StaticBuffer{4, 9, 252, 1, 1} // 252-1

// 你没有联盟
var ERR_GUILD_CHANGE_COUNTRY_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 252, 1, 2} // 252-2

// 你不是盟主，不能联盟转国
var ERR_GUILD_CHANGE_COUNTRY_FAIL_NOT_LEARDER = pbutil.StaticBuffer{4, 9, 252, 1, 3} // 252-3

// 消耗不足
var ERR_GUILD_CHANGE_COUNTRY_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 9, 252, 1, 4} // 252-4

// 转国cd中
var ERR_GUILD_CHANGE_COUNTRY_FAIL_COOLDOWN = pbutil.StaticBuffer{4, 9, 252, 1, 5} // 252-5

// 当前正在转国中
var ERR_GUILD_CHANGE_COUNTRY_FAIL_EXIST = pbutil.StaticBuffer{4, 9, 252, 1, 6} // 252-6

// 你的联盟已经是这个国家
var ERR_GUILD_CHANGE_COUNTRY_FAIL_SAME_COUNTRY = pbutil.StaticBuffer{4, 9, 252, 1, 7} // 252-7

// 国王不能转国
var ERR_GUILD_CHANGE_COUNTRY_FAIL_IS_KING = pbutil.StaticBuffer{4, 9, 252, 1, 8} // 252-8

var CANCEL_GUILD_CHANGE_COUNTRY_S2C = pbutil.StaticBuffer{3, 9, 254, 1} // 254

// 你没有联盟
var ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 9, 255, 1, 1} // 255-1

// 你不是盟主，不能取消联盟转国
var ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_LEARDER = pbutil.StaticBuffer{4, 9, 255, 1, 2} // 255-2

// 转国cd中
var ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_COOLDOWN = pbutil.StaticBuffer{4, 9, 255, 1, 3} // 255-3

// 联盟没有转国中
var ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_EXIST = pbutil.StaticBuffer{4, 9, 255, 1, 4} // 255-4

func NewS2cShowWorkshopNotExistMsg(show bool) pbutil.Buffer {
	msg := &S2CShowWorkshopNotExistProto{
		Show: show,
	}
	return NewS2cShowWorkshopNotExistProtoMsg(msg)
}

var s2c_show_workshop_not_exist = [...]byte{9, 129, 2} // 257
func NewS2cShowWorkshopNotExistProtoMsg(object *S2CShowWorkshopNotExistProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_show_workshop_not_exist[:], "s2c_show_workshop_not_exist")

}
