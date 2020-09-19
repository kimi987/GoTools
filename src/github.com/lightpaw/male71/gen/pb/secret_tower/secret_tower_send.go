package secret_tower

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
	MODULE_ID = 22

	C2S_REQUEST_TEAM_COUNT = 2

	C2S_REQUEST_TEAM_LIST = 5

	C2S_CREATE_TEAM = 8

	C2S_JOIN_TEAM = 11

	C2S_LEAVE_TEAM = 15

	C2S_KICK_MEMBER = 19

	C2S_MOVE_MEMBER = 24

	C2S_UPDATE_MEMBER_POS = 67

	C2S_CHANGE_MODE = 27

	C2S_INVITE = 33

	C2S_INVITE_ALL = 71

	C2S_REQUEST_INVITE_LIST = 37

	C2S_REQUEST_TEAM_DETAIL = 39

	C2S_START_CHALLENGE = 42

	C2S_QUICK_QUERY_TEAM_BASIC = 58

	C2S_CHANGE_GUILD_MODE = 61

	C2S_LIST_RECORD = 74

	C2S_TEAM_TALK = 79
)

func NewS2cUnlockSecretTowerMsg(secret_tower_id int32) pbutil.Buffer {
	msg := &S2CUnlockSecretTowerProto{
		SecretTowerId: secret_tower_id,
	}
	return NewS2cUnlockSecretTowerProtoMsg(msg)
}

var s2c_unlock_secret_tower = [...]byte{22, 1} // 1
func NewS2cUnlockSecretTowerProtoMsg(object *S2CUnlockSecretTowerProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_unlock_secret_tower[:], "s2c_unlock_secret_tower")

}

func NewS2cRequestTeamCountMsg(secret_tower_id []int32, team_count []int32) pbutil.Buffer {
	msg := &S2CRequestTeamCountProto{
		SecretTowerId: secret_tower_id,
		TeamCount:     team_count,
	}
	return NewS2cRequestTeamCountProtoMsg(msg)
}

var s2c_request_team_count = [...]byte{22, 3} // 3
func NewS2cRequestTeamCountProtoMsg(object *S2CRequestTeamCountProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_team_count[:], "s2c_request_team_count")

}

// 密室没有开启
var ERR_REQUEST_TEAM_COUNT_FAIL_NOT_OPEN = pbutil.StaticBuffer{3, 22, 4, 1} // 4-1

func NewS2cRequestTeamListMsg(secret_tower_id int32, team_list [][]byte) pbutil.Buffer {
	msg := &S2CRequestTeamListProto{
		SecretTowerId: secret_tower_id,
		TeamList:      team_list,
	}
	return NewS2cRequestTeamListProtoMsg(msg)
}

var s2c_request_team_list = [...]byte{22, 6} // 6
func NewS2cRequestTeamListProtoMsg(object *S2CRequestTeamListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_team_list[:], "s2c_request_team_list")

}

// 密室没有开启
var ERR_REQUEST_TEAM_LIST_FAIL_NOT_OPEN = pbutil.StaticBuffer{3, 22, 7, 1} // 7-1

// 未知密室
var ERR_REQUEST_TEAM_LIST_FAIL_UNKNOWN_SECRET_TOWER = pbutil.StaticBuffer{3, 22, 7, 2} // 7-2

func NewS2cCreateTeamMsg(team_detail []byte) pbutil.Buffer {
	msg := &S2CCreateTeamProto{
		TeamDetail: team_detail,
	}
	return NewS2cCreateTeamProtoMsg(msg)
}

var s2c_create_team = [...]byte{22, 9} // 9
func NewS2cCreateTeamProtoMsg(object *S2CCreateTeamProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_create_team[:], "s2c_create_team")

}

// 未知密室
var ERR_CREATE_TEAM_FAIL_UNKNOWN_TOWER_ID = pbutil.StaticBuffer{3, 22, 10, 1} // 10-1

// 密室未开启
var ERR_CREATE_TEAM_FAIL_UNOPEN = pbutil.StaticBuffer{3, 22, 10, 2} // 10-2

// 当前有队伍
var ERR_CREATE_TEAM_FAIL_HAVE_TEAM_NOW = pbutil.StaticBuffer{3, 22, 10, 3} // 10-3

// 没有联盟
var ERR_CREATE_TEAM_FAIL_NO_GUILD = pbutil.StaticBuffer{3, 22, 10, 4} // 10-4

// 次数不足
var ERR_CREATE_TEAM_FAIL_TIMES_NOT_ENOUGH = pbutil.StaticBuffer{3, 22, 10, 5} // 10-5

// 上阵武将未满
var ERR_CREATE_TEAM_FAIL_CAPTAIN_NOT_FULL = pbutil.StaticBuffer{3, 22, 10, 6} // 10-6

// 上阵武将超出上限
var ERR_CREATE_TEAM_FAIL_CAPTAIN_TOO_MUCH = pbutil.StaticBuffer{3, 22, 10, 7} // 10-7

// 上阵武将不存在
var ERR_CREATE_TEAM_FAIL_CAPTAIN_NOT_EXIST = pbutil.StaticBuffer{3, 22, 10, 8} // 10-8

// 上阵武将id重复
var ERR_CREATE_TEAM_FAIL_CAPTAIN_ID_DUPLICATE = pbutil.StaticBuffer{3, 22, 10, 9} // 10-9

// 服务器错误
var ERR_CREATE_TEAM_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 22, 10, 10} // 10-10

func NewS2cJoinTeamMsg(team_detail []byte) pbutil.Buffer {
	msg := &S2CJoinTeamProto{
		TeamDetail: team_detail,
	}
	return NewS2cJoinTeamProtoMsg(msg)
}

var s2c_join_team = [...]byte{22, 12} // 12
func NewS2cJoinTeamProtoMsg(object *S2CJoinTeamProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_join_team[:], "s2c_join_team")

}

func NewS2cOtherJoinJoinTeamMsg(member []byte, protect_end_time int32) pbutil.Buffer {
	msg := &S2COtherJoinJoinTeamProto{
		Member:         member,
		ProtectEndTime: protect_end_time,
	}
	return NewS2cOtherJoinJoinTeamProtoMsg(msg)
}

var s2c_other_join_join_team = [...]byte{22, 13} // 13
func NewS2cOtherJoinJoinTeamProtoMsg(object *S2COtherJoinJoinTeamProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_other_join_join_team[:], "s2c_other_join_join_team")

}

// 密室未开启
var ERR_JOIN_TEAM_FAIL_UNOPEN = pbutil.StaticBuffer{3, 22, 14, 1} // 14-1

// 当前有队伍
var ERR_JOIN_TEAM_FAIL_HAVE_TEAM_NOW = pbutil.StaticBuffer{3, 22, 14, 2} // 14-2

// 不是目标队伍联盟的
var ERR_JOIN_TEAM_FAIL_NOT_TARGET_GUILD = pbutil.StaticBuffer{3, 22, 14, 3} // 14-3

// 队伍没找到
var ERR_JOIN_TEAM_FAIL_TEAM_NOT_FOUND = pbutil.StaticBuffer{3, 22, 14, 4} // 14-4

// 队伍已满
var ERR_JOIN_TEAM_FAIL_TEAM_FULL = pbutil.StaticBuffer{3, 22, 14, 5} // 14-5

// 次数不足
var ERR_JOIN_TEAM_FAIL_TIMES_NOT_ENOUGH = pbutil.StaticBuffer{3, 22, 14, 6} // 14-6

// 上阵武将未满
var ERR_JOIN_TEAM_FAIL_CAPTAIN_NOT_FULL = pbutil.StaticBuffer{3, 22, 14, 7} // 14-7

// 上阵武将超出上限
var ERR_JOIN_TEAM_FAIL_CAPTAIN_TOO_MUCH = pbutil.StaticBuffer{3, 22, 14, 8} // 14-8

// 上阵武将不存在
var ERR_JOIN_TEAM_FAIL_CAPTAIN_NOT_EXIST = pbutil.StaticBuffer{3, 22, 14, 9} // 14-9

// 上阵武将id重复
var ERR_JOIN_TEAM_FAIL_CAPTAIN_ID_DUPLICATE = pbutil.StaticBuffer{3, 22, 14, 10} // 14-10

// 不可以协助自己能够参与的最高层的密室
var ERR_JOIN_TEAM_FAIL_CAN_NOT_HELP_MAX_TOWER = pbutil.StaticBuffer{3, 22, 14, 12} // 14-12

// 队伍中没有盟友，不能以协助模式加入队伍
var ERR_JOIN_TEAM_FAIL_CAN_NOT_HELP_NO_GUILD_MEMBER_TEAM = pbutil.StaticBuffer{3, 22, 14, 13} // 14-13

// 没找到合适的队伍可以加入，请尝试创建队伍
var ERR_JOIN_TEAM_FAIL_NOT_VALID_TEAM = pbutil.StaticBuffer{3, 22, 14, 14} // 14-14

// 服务器错误
var ERR_JOIN_TEAM_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 22, 14, 11} // 14-11

var LEAVE_TEAM_S2C = pbutil.StaticBuffer{2, 22, 16} // 16

func NewS2cOtherLeaveLeaveTeamMsg(id []byte, new_team_leader_id []byte) pbutil.Buffer {
	msg := &S2COtherLeaveLeaveTeamProto{
		Id:              id,
		NewTeamLeaderId: new_team_leader_id,
	}
	return NewS2cOtherLeaveLeaveTeamProtoMsg(msg)
}

var s2c_other_leave_leave_team = [...]byte{22, 17} // 17
func NewS2cOtherLeaveLeaveTeamProtoMsg(object *S2COtherLeaveLeaveTeamProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_other_leave_leave_team[:], "s2c_other_leave_leave_team")

}

// 没有队伍
var ERR_LEAVE_TEAM_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 18, 1} // 18-1

func NewS2cKickMemberMsg(id []byte) pbutil.Buffer {
	msg := &S2CKickMemberProto{
		Id: id,
	}
	return NewS2cKickMemberProtoMsg(msg)
}

var s2c_kick_member = [...]byte{22, 20} // 20
func NewS2cKickMemberProtoMsg(object *S2CKickMemberProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_kick_member[:], "s2c_kick_member")

}

var KICK_MEMBER_S2C_YOU_BEEN_KICKED = pbutil.StaticBuffer{2, 22, 21} // 21

func NewS2cOtherBeenKickKickMemberMsg(id []byte) pbutil.Buffer {
	msg := &S2COtherBeenKickKickMemberProto{
		Id: id,
	}
	return NewS2cOtherBeenKickKickMemberProtoMsg(msg)
}

var s2c_other_been_kick_kick_member = [...]byte{22, 22} // 22
func NewS2cOtherBeenKickKickMemberProtoMsg(object *S2COtherBeenKickKickMemberProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_other_been_kick_kick_member[:], "s2c_other_been_kick_kick_member")

}

// 没有队伍
var ERR_KICK_MEMBER_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 23, 1} // 23-1

// 目标没找到
var ERR_KICK_MEMBER_FAIL_TARGET_NOT_FOUND = pbutil.StaticBuffer{3, 22, 23, 2} // 23-2

// 不是队长
var ERR_KICK_MEMBER_FAIL_NOT_LEADER = pbutil.StaticBuffer{3, 22, 23, 3} // 23-3

// 不能踢出自己
var ERR_KICK_MEMBER_FAIL_CANT_KICK_SELF = pbutil.StaticBuffer{3, 22, 23, 4} // 23-4

func NewS2cBroadcsatMoveMemberMsg(id []byte, up bool) pbutil.Buffer {
	msg := &S2CBroadcsatMoveMemberProto{
		Id: id,
		Up: up,
	}
	return NewS2cBroadcsatMoveMemberProtoMsg(msg)
}

var s2c_broadcsat_move_member = [...]byte{22, 25} // 25
func NewS2cBroadcsatMoveMemberProtoMsg(object *S2CBroadcsatMoveMemberProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_broadcsat_move_member[:], "s2c_broadcsat_move_member")

}

// 没有队伍
var ERR_MOVE_MEMBER_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 26, 1} // 26-1

// 不是队长
var ERR_MOVE_MEMBER_FAIL_NOT_LEADER = pbutil.StaticBuffer{3, 22, 26, 2} // 26-2

// 要移动的成员没找到
var ERR_MOVE_MEMBER_FAIL_TARGET_NOT_FOUND = pbutil.StaticBuffer{3, 22, 26, 3} // 26-3

// 要移动的成员已经是第一个了
var ERR_MOVE_MEMBER_FAIL_TARGET_IS_FIRST = pbutil.StaticBuffer{3, 22, 26, 4} // 26-4

// 要移动的成员已经是最后一个了
var ERR_MOVE_MEMBER_FAIL_TARGET_IS_LAST = pbutil.StaticBuffer{3, 22, 26, 5} // 26-5

func NewS2cUpdateMemberPosMsg(id [][]byte) pbutil.Buffer {
	msg := &S2CUpdateMemberPosProto{
		Id: id,
	}
	return NewS2cUpdateMemberPosProtoMsg(msg)
}

var s2c_update_member_pos = [...]byte{22, 68} // 68
func NewS2cUpdateMemberPosProtoMsg(object *S2CUpdateMemberPosProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_member_pos[:], "s2c_update_member_pos")

}

// 没有队伍
var ERR_UPDATE_MEMBER_POS_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 69, 1} // 69-1

// 不是队长
var ERR_UPDATE_MEMBER_POS_FAIL_NOT_LEADER = pbutil.StaticBuffer{3, 22, 69, 2} // 69-2

// 要移动的成员没找到
var ERR_UPDATE_MEMBER_POS_FAIL_TARGET_NOT_FOUND = pbutil.StaticBuffer{3, 22, 69, 3} // 69-3

// 发送的成员列表重复
var ERR_UPDATE_MEMBER_POS_FAIL_TARGET_DUPLICATE = pbutil.StaticBuffer{3, 22, 69, 4} // 69-4

func NewS2cChangeModeMsg(mode int32) pbutil.Buffer {
	msg := &S2CChangeModeProto{
		Mode: mode,
	}
	return NewS2cChangeModeProtoMsg(msg)
}

var s2c_change_mode = [...]byte{22, 28} // 28
func NewS2cChangeModeProtoMsg(object *S2CChangeModeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_mode[:], "s2c_change_mode")

}

func NewS2cOtherChangedChangeModeMsg(id []byte, mode int32) pbutil.Buffer {
	msg := &S2COtherChangedChangeModeProto{
		Id:   id,
		Mode: mode,
	}
	return NewS2cOtherChangedChangeModeProtoMsg(msg)
}

var s2c_other_changed_change_mode = [...]byte{22, 29} // 29
func NewS2cOtherChangedChangeModeProtoMsg(object *S2COtherChangedChangeModeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_other_changed_change_mode[:], "s2c_other_changed_change_mode")

}

// 没有队伍
var ERR_CHANGE_MODE_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 30, 1} // 30-1

// 未知模式
var ERR_CHANGE_MODE_FAIL_UNKNOWN_MODE = pbutil.StaticBuffer{3, 22, 30, 7} // 30-7

// 没次数了，无法调整到目标模式
var ERR_CHANGE_MODE_FAIL_NO_TIMES = pbutil.StaticBuffer{3, 22, 30, 3} // 30-3

// 模式没变
var ERR_CHANGE_MODE_FAIL_MODE_NOT_CHANGE = pbutil.StaticBuffer{3, 22, 30, 4} // 30-4

// 队长不可以变更模式
var ERR_CHANGE_MODE_FAIL_IS_LEADER = pbutil.StaticBuffer{3, 22, 30, 5} // 30-5

// 不可以协助自己能够参与的最高层的密室
var ERR_CHANGE_MODE_FAIL_CAN_NOT_HELP_MAX_TOWER = pbutil.StaticBuffer{3, 22, 30, 8} // 30-8

// 队伍中没有盟友，不能协助
var ERR_CHANGE_MODE_FAIL_CAN_NOT_HELP_NO_GUILD_MEMBER = pbutil.StaticBuffer{3, 22, 30, 9} // 30-9

// 没有协助次数了
var ERR_CHANGE_MODE_FAIL_NO_HELP_TIMES = pbutil.StaticBuffer{3, 22, 30, 10} // 30-10

// 服务器错误
var ERR_CHANGE_MODE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 22, 30, 6} // 30-6

func NewS2cHelpTimesChangeMsg(new_times int32) pbutil.Buffer {
	msg := &S2CHelpTimesChangeProto{
		NewTimes: new_times,
	}
	return NewS2cHelpTimesChangeProtoMsg(msg)
}

var s2c_help_times_change = [...]byte{22, 31} // 31
func NewS2cHelpTimesChangeProtoMsg(object *S2CHelpTimesChangeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_help_times_change[:], "s2c_help_times_change")

}

func NewS2cTimesChangeMsg(challenge_times int32, history_challenge_times int32) pbutil.Buffer {
	msg := &S2CTimesChangeProto{
		ChallengeTimes:        challenge_times,
		HistoryChallengeTimes: history_challenge_times,
	}
	return NewS2cTimesChangeProtoMsg(msg)
}

var s2c_times_change = [...]byte{22, 32} // 32
func NewS2cTimesChangeProtoMsg(object *S2CTimesChangeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_times_change[:], "s2c_times_change")

}

func NewS2cInviteMsg(id []byte) pbutil.Buffer {
	msg := &S2CInviteProto{
		Id: id,
	}
	return NewS2cInviteProtoMsg(msg)
}

var s2c_invite = [...]byte{22, 34} // 34
func NewS2cInviteProtoMsg(object *S2CInviteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_invite[:], "s2c_invite")

}

func NewS2cFailTargetNotFoundInviteMsg(id []byte) pbutil.Buffer {
	msg := &S2CFailTargetNotFoundInviteProto{
		Id: id,
	}
	return NewS2cFailTargetNotFoundInviteProtoMsg(msg)
}

var s2c_fail_target_not_found_invite = [...]byte{22, 51} // 51
func NewS2cFailTargetNotFoundInviteProtoMsg(object *S2CFailTargetNotFoundInviteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_target_not_found_invite[:], "s2c_fail_target_not_found_invite")

}

func NewS2cFailTargetNotInMyGuildInviteMsg(id []byte) pbutil.Buffer {
	msg := &S2CFailTargetNotInMyGuildInviteProto{
		Id: id,
	}
	return NewS2cFailTargetNotInMyGuildInviteProtoMsg(msg)
}

var s2c_fail_target_not_in_my_guild_invite = [...]byte{22, 52} // 52
func NewS2cFailTargetNotInMyGuildInviteProtoMsg(object *S2CFailTargetNotInMyGuildInviteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_target_not_in_my_guild_invite[:], "s2c_fail_target_not_in_my_guild_invite")

}

func NewS2cFailTargetNotOpenInviteMsg(id []byte) pbutil.Buffer {
	msg := &S2CFailTargetNotOpenInviteProto{
		Id: id,
	}
	return NewS2cFailTargetNotOpenInviteProtoMsg(msg)
}

var s2c_fail_target_not_open_invite = [...]byte{22, 53} // 53
func NewS2cFailTargetNotOpenInviteProtoMsg(object *S2CFailTargetNotOpenInviteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_target_not_open_invite[:], "s2c_fail_target_not_open_invite")

}

func NewS2cFailTargetNotOnlineInviteMsg(id []byte) pbutil.Buffer {
	msg := &S2CFailTargetNotOnlineInviteProto{
		Id: id,
	}
	return NewS2cFailTargetNotOnlineInviteProtoMsg(msg)
}

var s2c_fail_target_not_online_invite = [...]byte{22, 54} // 54
func NewS2cFailTargetNotOnlineInviteProtoMsg(object *S2CFailTargetNotOnlineInviteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_target_not_online_invite[:], "s2c_fail_target_not_online_invite")

}

func NewS2cFailTargetInYourTeamInviteMsg(id []byte) pbutil.Buffer {
	msg := &S2CFailTargetInYourTeamInviteProto{
		Id: id,
	}
	return NewS2cFailTargetInYourTeamInviteProtoMsg(msg)
}

var s2c_fail_target_in_your_team_invite = [...]byte{22, 55} // 55
func NewS2cFailTargetInYourTeamInviteProtoMsg(object *S2CFailTargetInYourTeamInviteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_target_in_your_team_invite[:], "s2c_fail_target_in_your_team_invite")

}

func NewS2cFailTargetNoTimesInviteMsg(id []byte) pbutil.Buffer {
	msg := &S2CFailTargetNoTimesInviteProto{
		Id: id,
	}
	return NewS2cFailTargetNoTimesInviteProtoMsg(msg)
}

var s2c_fail_target_no_times_invite = [...]byte{22, 56} // 56
func NewS2cFailTargetNoTimesInviteProtoMsg(object *S2CFailTargetNoTimesInviteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_target_no_times_invite[:], "s2c_fail_target_no_times_invite")

}

// 没有队伍
var ERR_INVITE_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 35, 1} // 35-1

// 队伍已满
var ERR_INVITE_FAIL_TEAM_FULL = pbutil.StaticBuffer{3, 22, 35, 2} // 35-2

// 服务器错误
var ERR_INVITE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 22, 35, 9} // 35-9

func NewS2cInviteAllMsg(id [][]byte) pbutil.Buffer {
	msg := &S2CInviteAllProto{
		Id: id,
	}
	return NewS2cInviteAllProtoMsg(msg)
}

var s2c_invite_all = [...]byte{22, 72} // 72
func NewS2cInviteAllProtoMsg(object *S2CInviteAllProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_invite_all[:], "s2c_invite_all")

}

// 无效的id
var ERR_INVITE_ALL_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 22, 73, 4} // 73-4

// 没有队伍
var ERR_INVITE_ALL_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 73, 1} // 73-1

// 队伍已满
var ERR_INVITE_ALL_FAIL_TEAM_FULL = pbutil.StaticBuffer{3, 22, 73, 2} // 73-2

// 服务器错误
var ERR_INVITE_ALL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 22, 73, 3} // 73-3

func NewS2cReceiveInviteMsg(count int32, have_new bool) pbutil.Buffer {
	msg := &S2CReceiveInviteProto{
		Count:   count,
		HaveNew: have_new,
	}
	return NewS2cReceiveInviteProtoMsg(msg)
}

var s2c_receive_invite = [...]byte{22, 36} // 36
func NewS2cReceiveInviteProtoMsg(object *S2CReceiveInviteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_receive_invite[:], "s2c_receive_invite")

}

func NewS2cRequestInviteListMsg(invite_list [][]byte) pbutil.Buffer {
	msg := &S2CRequestInviteListProto{
		InviteList: invite_list,
	}
	return NewS2cRequestInviteListProtoMsg(msg)
}

var s2c_request_invite_list = [...]byte{22, 38} // 38
func NewS2cRequestInviteListProtoMsg(object *S2CRequestInviteListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_invite_list[:], "s2c_request_invite_list")

}

func NewS2cRequestTeamDetailMsg(team_detail []byte) pbutil.Buffer {
	msg := &S2CRequestTeamDetailProto{
		TeamDetail: team_detail,
	}
	return NewS2cRequestTeamDetailProtoMsg(msg)
}

var s2c_request_team_detail = [...]byte{22, 40} // 40
func NewS2cRequestTeamDetailProtoMsg(object *S2CRequestTeamDetailProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_team_detail[:], "s2c_request_team_detail")

}

// 没有队伍
var ERR_REQUEST_TEAM_DETAIL_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 41, 1} // 41-1

func NewS2cBroadcastStartChallengeMsg(result []byte) pbutil.Buffer {
	msg := &S2CBroadcastStartChallengeProto{
		Result: result,
	}
	return NewS2cBroadcastStartChallengeProtoMsg(msg)
}

var s2c_broadcast_start_challenge = [...]byte{22, 43} // 43
func NewS2cBroadcastStartChallengeProtoMsg(object *S2CBroadcastStartChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_broadcast_start_challenge[:], "s2c_broadcast_start_challenge")

}

func NewS2cFailWithMemberTimesNotEnoughStartChallengeMsg(id []byte, name string, guild_flag string) pbutil.Buffer {
	msg := &S2CFailWithMemberTimesNotEnoughStartChallengeProto{
		Id:        id,
		Name:      name,
		GuildFlag: guild_flag,
	}
	return NewS2cFailWithMemberTimesNotEnoughStartChallengeProtoMsg(msg)
}

var s2c_fail_with_member_times_not_enough_start_challenge = [...]byte{22, 44} // 44
func NewS2cFailWithMemberTimesNotEnoughStartChallengeProtoMsg(object *S2CFailWithMemberTimesNotEnoughStartChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_with_member_times_not_enough_start_challenge[:], "s2c_fail_with_member_times_not_enough_start_challenge")

}

func NewS2cFailWithMemberHelpTimesNotEnoughStartChallengeMsg(id []byte, name string, guild_flag string) pbutil.Buffer {
	msg := &S2CFailWithMemberHelpTimesNotEnoughStartChallengeProto{
		Id:        id,
		Name:      name,
		GuildFlag: guild_flag,
	}
	return NewS2cFailWithMemberHelpTimesNotEnoughStartChallengeProtoMsg(msg)
}

var s2c_fail_with_member_help_times_not_enough_start_challenge = [...]byte{22, 45} // 45
func NewS2cFailWithMemberHelpTimesNotEnoughStartChallengeProtoMsg(object *S2CFailWithMemberHelpTimesNotEnoughStartChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_with_member_help_times_not_enough_start_challenge[:], "s2c_fail_with_member_help_times_not_enough_start_challenge")

}

func NewS2cFailWithMemberNoGuildStartChallengeMsg(id []byte, name string, guild_flag string) pbutil.Buffer {
	msg := &S2CFailWithMemberNoGuildStartChallengeProto{
		Id:        id,
		Name:      name,
		GuildFlag: guild_flag,
	}
	return NewS2cFailWithMemberNoGuildStartChallengeProtoMsg(msg)
}

var s2c_fail_with_member_no_guild_start_challenge = [...]byte{22, 46} // 46
func NewS2cFailWithMemberNoGuildStartChallengeProtoMsg(object *S2CFailWithMemberNoGuildStartChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_with_member_no_guild_start_challenge[:], "s2c_fail_with_member_no_guild_start_challenge")

}

func NewS2cFailWithMemberNotMyGuildStartChallengeMsg(id []byte, name string, guild_flag string) pbutil.Buffer {
	msg := &S2CFailWithMemberNotMyGuildStartChallengeProto{
		Id:        id,
		Name:      name,
		GuildFlag: guild_flag,
	}
	return NewS2cFailWithMemberNotMyGuildStartChallengeProtoMsg(msg)
}

var s2c_fail_with_member_not_my_guild_start_challenge = [...]byte{22, 47} // 47
func NewS2cFailWithMemberNotMyGuildStartChallengeProtoMsg(object *S2CFailWithMemberNotMyGuildStartChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_with_member_not_my_guild_start_challenge[:], "s2c_fail_with_member_not_my_guild_start_challenge")

}

func NewS2cFailWithMemberIsHelpButNoGuildStartChallengeMsg(id []byte, name string, guild_flag string) pbutil.Buffer {
	msg := &S2CFailWithMemberIsHelpButNoGuildStartChallengeProto{
		Id:        id,
		Name:      name,
		GuildFlag: guild_flag,
	}
	return NewS2cFailWithMemberIsHelpButNoGuildStartChallengeProtoMsg(msg)
}

var s2c_fail_with_member_is_help_but_no_guild_start_challenge = [...]byte{22, 48} // 48
func NewS2cFailWithMemberIsHelpButNoGuildStartChallengeProtoMsg(object *S2CFailWithMemberIsHelpButNoGuildStartChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_with_member_is_help_but_no_guild_start_challenge[:], "s2c_fail_with_member_is_help_but_no_guild_start_challenge")

}

func NewS2cFailWithMemberIsHelpButNoGuildMemberStartChallengeMsg(id []byte, name string, guild_flag string) pbutil.Buffer {
	msg := &S2CFailWithMemberIsHelpButNoGuildMemberStartChallengeProto{
		Id:        id,
		Name:      name,
		GuildFlag: guild_flag,
	}
	return NewS2cFailWithMemberIsHelpButNoGuildMemberStartChallengeProtoMsg(msg)
}

var s2c_fail_with_member_is_help_but_no_guild_member_start_challenge = [...]byte{22, 49} // 49
func NewS2cFailWithMemberIsHelpButNoGuildMemberStartChallengeProtoMsg(object *S2CFailWithMemberIsHelpButNoGuildMemberStartChallengeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fail_with_member_is_help_but_no_guild_member_start_challenge[:], "s2c_fail_with_member_is_help_but_no_guild_member_start_challenge")

}

// 没有队伍
var ERR_START_CHALLENGE_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 50, 1} // 50-1

// 不是队长
var ERR_START_CHALLENGE_FAIL_NOT_TEAM_LEADER = pbutil.StaticBuffer{3, 22, 50, 2} // 50-2

// 队伍人数不够
var ERR_START_CHALLENGE_FAIL_TEAM_MEMBER_NOT_ENOUGH = pbutil.StaticBuffer{3, 22, 50, 3} // 50-3

// 请等待其他人都准备好
var ERR_START_CHALLENGE_FAIL_WAIT_PROTECT_END = pbutil.StaticBuffer{3, 22, 50, 4} // 50-4

// 没有人是挑战模式，无法开启
var ERR_START_CHALLENGE_FAIL_NO_CHALLENGE_PEOPLE = pbutil.StaticBuffer{3, 22, 50, 5} // 50-5

// 服务器繁忙，请稍后再试
var ERR_START_CHALLENGE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 22, 50, 6} // 50-6

var TEAM_EXPIRED_S2C = pbutil.StaticBuffer{2, 22, 57} // 57

var TEAM_DESTROYED_BECAUSE_OF_LEADER_LEAVE_S2C = pbutil.StaticBuffer{2, 22, 65} // 65

func NewS2cQuickQueryTeamBasicMsg(basics [][]byte, not_exist_ids []int32) pbutil.Buffer {
	msg := &S2CQuickQueryTeamBasicProto{
		Basics:      basics,
		NotExistIds: not_exist_ids,
	}
	return NewS2cQuickQueryTeamBasicProtoMsg(msg)
}

var s2c_quick_query_team_basic = [...]byte{22, 59} // 59
func NewS2cQuickQueryTeamBasicProtoMsg(object *S2CQuickQueryTeamBasicProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_quick_query_team_basic[:], "s2c_quick_query_team_basic")

}

// 服务器繁忙，请稍后再试
var ERR_QUICK_QUERY_TEAM_BASIC_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 22, 60, 1} // 60-1

var CHANGE_GUILD_MODE_S2C = pbutil.StaticBuffer{2, 22, 62} // 62

// 没有队伍
var ERR_CHANGE_GUILD_MODE_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 63, 1} // 63-1

// 不是队长
var ERR_CHANGE_GUILD_MODE_FAIL_NOT_LEADER = pbutil.StaticBuffer{3, 22, 63, 2} // 63-2

// 没有加入联盟，无法变更为联盟模式
var ERR_CHANGE_GUILD_MODE_FAIL_NO_GUILD = pbutil.StaticBuffer{3, 22, 63, 3} // 63-3

// 有玩家没有加入我们联盟
var ERR_CHANGE_GUILD_MODE_FAIL_SB_NOT_IN_MY_GUILD = pbutil.StaticBuffer{3, 22, 63, 4} // 63-4

// 服务器繁忙，请稍后再试
var ERR_CHANGE_GUILD_MODE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 22, 63, 5} // 63-5

func NewS2cChangeGuildModeBroadcastMsg(guild_id int32) pbutil.Buffer {
	msg := &S2CChangeGuildModeBroadcastProto{
		GuildId: guild_id,
	}
	return NewS2cChangeGuildModeBroadcastProtoMsg(msg)
}

var s2c_change_guild_mode_broadcast = [...]byte{22, 64} // 64
func NewS2cChangeGuildModeBroadcastProtoMsg(object *S2CChangeGuildModeBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_guild_mode_broadcast[:], "s2c_change_guild_mode_broadcast")

}

func NewS2cMemberTroopChangedMsg(member []byte, protect_end_time int32) pbutil.Buffer {
	msg := &S2CMemberTroopChangedProto{
		Member:         member,
		ProtectEndTime: protect_end_time,
	}
	return NewS2cMemberTroopChangedProtoMsg(msg)
}

var s2c_member_troop_changed = [...]byte{22, 66} // 66
func NewS2cMemberTroopChangedProtoMsg(object *S2CMemberTroopChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_member_troop_changed[:], "s2c_member_troop_changed")

}

func NewS2cListRecordMsg(record []*shared_proto.SecretRecordProto) pbutil.Buffer {
	msg := &S2CListRecordProto{
		Record: record,
	}
	return NewS2cListRecordProtoMsg(msg)
}

var s2c_list_record = [...]byte{22, 75} // 75
func NewS2cListRecordProtoMsg(object *S2CListRecordProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_record[:], "s2c_list_record")

}

func NewS2cTeamTalkMsg(words_id int32, text string) pbutil.Buffer {
	msg := &S2CTeamTalkProto{
		WordsId: words_id,
		Text:    text,
	}
	return NewS2cTeamTalkProtoMsg(msg)
}

var s2c_team_talk = [...]byte{22, 80} // 80
func NewS2cTeamTalkProtoMsg(object *S2CTeamTalkProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_team_talk[:], "s2c_team_talk")

}

// id 错误
var ERR_TEAM_TALK_FAIL_INVALID_WORDS = pbutil.StaticBuffer{3, 22, 81, 1} // 81-1

// 没有队伍
var ERR_TEAM_TALK_FAIL_NO_TEAM = pbutil.StaticBuffer{3, 22, 81, 2} // 81-2

func NewS2cTeamWhoTalkMsg(hero_id []byte, words_id int32, text string) pbutil.Buffer {
	msg := &S2CTeamWhoTalkProto{
		HeroId:  hero_id,
		WordsId: words_id,
		Text:    text,
	}
	return NewS2cTeamWhoTalkProtoMsg(msg)
}

var s2c_team_who_talk = [...]byte{22, 82} // 82
func NewS2cTeamWhoTalkProtoMsg(object *S2CTeamWhoTalkProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_team_who_talk[:], "s2c_team_who_talk")

}

func NewS2cTeamHistoryTalkMsg(records []*shared_proto.SecretTowerChatRecordProto) pbutil.Buffer {
	msg := &S2CTeamHistoryTalkProto{
		Records: records,
	}
	return NewS2cTeamHistoryTalkProtoMsg(msg)
}

var s2c_team_history_talk = [...]byte{22, 83} // 83
func NewS2cTeamHistoryTalkProtoMsg(object *S2CTeamHistoryTalkProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_team_history_talk[:], "s2c_team_history_talk")

}
