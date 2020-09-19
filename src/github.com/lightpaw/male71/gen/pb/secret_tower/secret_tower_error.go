package secret_tower

import (
	"github.com/lightpaw/pbutil"
)

// request_team_count
var (
	ErrRequestTeamCountFailNotOpen = newMsgError("request_team_count 密室没有开启", ERR_REQUEST_TEAM_COUNT_FAIL_NOT_OPEN) // 4-1
)

// request_team_list
var (
	ErrRequestTeamListFailNotOpen            = newMsgError("request_team_list 密室没有开启", ERR_REQUEST_TEAM_LIST_FAIL_NOT_OPEN)           // 7-1
	ErrRequestTeamListFailUnknownSecretTower = newMsgError("request_team_list 未知密室", ERR_REQUEST_TEAM_LIST_FAIL_UNKNOWN_SECRET_TOWER) // 7-2
)

// create_team
var (
	ErrCreateTeamFailUnknownTowerId     = newMsgError("create_team 未知密室", ERR_CREATE_TEAM_FAIL_UNKNOWN_TOWER_ID)         // 10-1
	ErrCreateTeamFailUnopen             = newMsgError("create_team 密室未开启", ERR_CREATE_TEAM_FAIL_UNOPEN)                  // 10-2
	ErrCreateTeamFailHaveTeamNow        = newMsgError("create_team 当前有队伍", ERR_CREATE_TEAM_FAIL_HAVE_TEAM_NOW)           // 10-3
	ErrCreateTeamFailNoGuild            = newMsgError("create_team 没有联盟", ERR_CREATE_TEAM_FAIL_NO_GUILD)                 // 10-4
	ErrCreateTeamFailTimesNotEnough     = newMsgError("create_team 次数不足", ERR_CREATE_TEAM_FAIL_TIMES_NOT_ENOUGH)         // 10-5
	ErrCreateTeamFailCaptainNotFull     = newMsgError("create_team 上阵武将未满", ERR_CREATE_TEAM_FAIL_CAPTAIN_NOT_FULL)       // 10-6
	ErrCreateTeamFailCaptainTooMuch     = newMsgError("create_team 上阵武将超出上限", ERR_CREATE_TEAM_FAIL_CAPTAIN_TOO_MUCH)     // 10-7
	ErrCreateTeamFailCaptainNotExist    = newMsgError("create_team 上阵武将不存在", ERR_CREATE_TEAM_FAIL_CAPTAIN_NOT_EXIST)     // 10-8
	ErrCreateTeamFailCaptainIdDuplicate = newMsgError("create_team 上阵武将id重复", ERR_CREATE_TEAM_FAIL_CAPTAIN_ID_DUPLICATE) // 10-9
	ErrCreateTeamFailServerError        = newMsgError("create_team 服务器错误", ERR_CREATE_TEAM_FAIL_SERVER_ERROR)            // 10-10
)

// join_team
var (
	ErrJoinTeamFailUnopen                      = newMsgError("join_team 密室未开启", ERR_JOIN_TEAM_FAIL_UNOPEN)                                          // 14-1
	ErrJoinTeamFailHaveTeamNow                 = newMsgError("join_team 当前有队伍", ERR_JOIN_TEAM_FAIL_HAVE_TEAM_NOW)                                   // 14-2
	ErrJoinTeamFailNotTargetGuild              = newMsgError("join_team 不是目标队伍联盟的", ERR_JOIN_TEAM_FAIL_NOT_TARGET_GUILD)                            // 14-3
	ErrJoinTeamFailTeamNotFound                = newMsgError("join_team 队伍没找到", ERR_JOIN_TEAM_FAIL_TEAM_NOT_FOUND)                                  // 14-4
	ErrJoinTeamFailTeamFull                    = newMsgError("join_team 队伍已满", ERR_JOIN_TEAM_FAIL_TEAM_FULL)                                        // 14-5
	ErrJoinTeamFailTimesNotEnough              = newMsgError("join_team 次数不足", ERR_JOIN_TEAM_FAIL_TIMES_NOT_ENOUGH)                                 // 14-6
	ErrJoinTeamFailCaptainNotFull              = newMsgError("join_team 上阵武将未满", ERR_JOIN_TEAM_FAIL_CAPTAIN_NOT_FULL)                               // 14-7
	ErrJoinTeamFailCaptainTooMuch              = newMsgError("join_team 上阵武将超出上限", ERR_JOIN_TEAM_FAIL_CAPTAIN_TOO_MUCH)                             // 14-8
	ErrJoinTeamFailCaptainNotExist             = newMsgError("join_team 上阵武将不存在", ERR_JOIN_TEAM_FAIL_CAPTAIN_NOT_EXIST)                             // 14-9
	ErrJoinTeamFailCaptainIdDuplicate          = newMsgError("join_team 上阵武将id重复", ERR_JOIN_TEAM_FAIL_CAPTAIN_ID_DUPLICATE)                         // 14-10
	ErrJoinTeamFailCanNotHelpMaxTower          = newMsgError("join_team 不可以协助自己能够参与的最高层的密室", ERR_JOIN_TEAM_FAIL_CAN_NOT_HELP_MAX_TOWER)             // 14-12
	ErrJoinTeamFailCanNotHelpNoGuildMemberTeam = newMsgError("join_team 队伍中没有盟友，不能以协助模式加入队伍", ERR_JOIN_TEAM_FAIL_CAN_NOT_HELP_NO_GUILD_MEMBER_TEAM) // 14-13
	ErrJoinTeamFailNotValidTeam                = newMsgError("join_team 没找到合适的队伍可以加入，请尝试创建队伍", ERR_JOIN_TEAM_FAIL_NOT_VALID_TEAM)                   // 14-14
	ErrJoinTeamFailServerError                 = newMsgError("join_team 服务器错误", ERR_JOIN_TEAM_FAIL_SERVER_ERROR)                                    // 14-11
)

// leave_team
var (
	ErrLeaveTeamFailNoTeam = newMsgError("leave_team 没有队伍", ERR_LEAVE_TEAM_FAIL_NO_TEAM) // 18-1
)

// kick_member
var (
	ErrKickMemberFailNoTeam         = newMsgError("kick_member 没有队伍", ERR_KICK_MEMBER_FAIL_NO_TEAM)           // 23-1
	ErrKickMemberFailTargetNotFound = newMsgError("kick_member 目标没找到", ERR_KICK_MEMBER_FAIL_TARGET_NOT_FOUND) // 23-2
	ErrKickMemberFailNotLeader      = newMsgError("kick_member 不是队长", ERR_KICK_MEMBER_FAIL_NOT_LEADER)        // 23-3
	ErrKickMemberFailCantKickSelf   = newMsgError("kick_member 不能踢出自己", ERR_KICK_MEMBER_FAIL_CANT_KICK_SELF)  // 23-4
)

// move_member
var (
	ErrMoveMemberFailNoTeam         = newMsgError("move_member 没有队伍", ERR_MOVE_MEMBER_FAIL_NO_TEAM)                  // 26-1
	ErrMoveMemberFailNotLeader      = newMsgError("move_member 不是队长", ERR_MOVE_MEMBER_FAIL_NOT_LEADER)               // 26-2
	ErrMoveMemberFailTargetNotFound = newMsgError("move_member 要移动的成员没找到", ERR_MOVE_MEMBER_FAIL_TARGET_NOT_FOUND)    // 26-3
	ErrMoveMemberFailTargetIsFirst  = newMsgError("move_member 要移动的成员已经是第一个了", ERR_MOVE_MEMBER_FAIL_TARGET_IS_FIRST) // 26-4
	ErrMoveMemberFailTargetIsLast   = newMsgError("move_member 要移动的成员已经是最后一个了", ERR_MOVE_MEMBER_FAIL_TARGET_IS_LAST) // 26-5
)

// update_member_pos
var (
	ErrUpdateMemberPosFailNoTeam          = newMsgError("update_member_pos 没有队伍", ERR_UPDATE_MEMBER_POS_FAIL_NO_TEAM)               // 69-1
	ErrUpdateMemberPosFailNotLeader       = newMsgError("update_member_pos 不是队长", ERR_UPDATE_MEMBER_POS_FAIL_NOT_LEADER)            // 69-2
	ErrUpdateMemberPosFailTargetNotFound  = newMsgError("update_member_pos 要移动的成员没找到", ERR_UPDATE_MEMBER_POS_FAIL_TARGET_NOT_FOUND) // 69-3
	ErrUpdateMemberPosFailTargetDuplicate = newMsgError("update_member_pos 发送的成员列表重复", ERR_UPDATE_MEMBER_POS_FAIL_TARGET_DUPLICATE) // 69-4
)

// change_mode
var (
	ErrChangeModeFailNoTeam                  = newMsgError("change_mode 没有队伍", ERR_CHANGE_MODE_FAIL_NO_TEAM)                              // 30-1
	ErrChangeModeFailUnknownMode             = newMsgError("change_mode 未知模式", ERR_CHANGE_MODE_FAIL_UNKNOWN_MODE)                         // 30-7
	ErrChangeModeFailNoTimes                 = newMsgError("change_mode 没次数了，无法调整到目标模式", ERR_CHANGE_MODE_FAIL_NO_TIMES)                   // 30-3
	ErrChangeModeFailModeNotChange           = newMsgError("change_mode 模式没变", ERR_CHANGE_MODE_FAIL_MODE_NOT_CHANGE)                      // 30-4
	ErrChangeModeFailIsLeader                = newMsgError("change_mode 队长不可以变更模式", ERR_CHANGE_MODE_FAIL_IS_LEADER)                       // 30-5
	ErrChangeModeFailCanNotHelpMaxTower      = newMsgError("change_mode 不可以协助自己能够参与的最高层的密室", ERR_CHANGE_MODE_FAIL_CAN_NOT_HELP_MAX_TOWER) // 30-8
	ErrChangeModeFailCanNotHelpNoGuildMember = newMsgError("change_mode 队伍中没有盟友，不能协助", ERR_CHANGE_MODE_FAIL_CAN_NOT_HELP_NO_GUILD_MEMBER) // 30-9
	ErrChangeModeFailNoHelpTimes             = newMsgError("change_mode 没有协助次数了", ERR_CHANGE_MODE_FAIL_NO_HELP_TIMES)                     // 30-10
	ErrChangeModeFailServerError             = newMsgError("change_mode 服务器错误", ERR_CHANGE_MODE_FAIL_SERVER_ERROR)                        // 30-6
)

// invite
var (
	ErrInviteFailNoTeam      = newMsgError("invite 没有队伍", ERR_INVITE_FAIL_NO_TEAM)       // 35-1
	ErrInviteFailTeamFull    = newMsgError("invite 队伍已满", ERR_INVITE_FAIL_TEAM_FULL)     // 35-2
	ErrInviteFailServerError = newMsgError("invite 服务器错误", ERR_INVITE_FAIL_SERVER_ERROR) // 35-9
)

// invite_all
var (
	ErrInviteAllFailInvalidId   = newMsgError("invite_all 无效的id", ERR_INVITE_ALL_FAIL_INVALID_ID)   // 73-4
	ErrInviteAllFailNoTeam      = newMsgError("invite_all 没有队伍", ERR_INVITE_ALL_FAIL_NO_TEAM)       // 73-1
	ErrInviteAllFailTeamFull    = newMsgError("invite_all 队伍已满", ERR_INVITE_ALL_FAIL_TEAM_FULL)     // 73-2
	ErrInviteAllFailServerError = newMsgError("invite_all 服务器错误", ERR_INVITE_ALL_FAIL_SERVER_ERROR) // 73-3
)

// request_team_detail
var (
	ErrRequestTeamDetailFailNoTeam = newMsgError("request_team_detail 没有队伍", ERR_REQUEST_TEAM_DETAIL_FAIL_NO_TEAM) // 41-1
)

// start_challenge
var (
	ErrStartChallengeFailNoTeam              = newMsgError("start_challenge 没有队伍", ERR_START_CHALLENGE_FAIL_NO_TEAM)                      // 50-1
	ErrStartChallengeFailNotTeamLeader       = newMsgError("start_challenge 不是队长", ERR_START_CHALLENGE_FAIL_NOT_TEAM_LEADER)              // 50-2
	ErrStartChallengeFailTeamMemberNotEnough = newMsgError("start_challenge 队伍人数不够", ERR_START_CHALLENGE_FAIL_TEAM_MEMBER_NOT_ENOUGH)     // 50-3
	ErrStartChallengeFailWaitProtectEnd      = newMsgError("start_challenge 请等待其他人都准备好", ERR_START_CHALLENGE_FAIL_WAIT_PROTECT_END)       // 50-4
	ErrStartChallengeFailNoChallengePeople   = newMsgError("start_challenge 没有人是挑战模式，无法开启", ERR_START_CHALLENGE_FAIL_NO_CHALLENGE_PEOPLE) // 50-5
	ErrStartChallengeFailServerError         = newMsgError("start_challenge 服务器繁忙，请稍后再试", ERR_START_CHALLENGE_FAIL_SERVER_ERROR)          // 50-6
)

// quick_query_team_basic
var (
	ErrQuickQueryTeamBasicFailServerError = newMsgError("quick_query_team_basic 服务器繁忙，请稍后再试", ERR_QUICK_QUERY_TEAM_BASIC_FAIL_SERVER_ERROR) // 60-1
)

// change_guild_mode
var (
	ErrChangeGuildModeFailNoTeam         = newMsgError("change_guild_mode 没有队伍", ERR_CHANGE_GUILD_MODE_FAIL_NO_TEAM)                   // 63-1
	ErrChangeGuildModeFailNotLeader      = newMsgError("change_guild_mode 不是队长", ERR_CHANGE_GUILD_MODE_FAIL_NOT_LEADER)                // 63-2
	ErrChangeGuildModeFailNoGuild        = newMsgError("change_guild_mode 没有加入联盟，无法变更为联盟模式", ERR_CHANGE_GUILD_MODE_FAIL_NO_GUILD)      // 63-3
	ErrChangeGuildModeFailSbNotInMyGuild = newMsgError("change_guild_mode 有玩家没有加入我们联盟", ERR_CHANGE_GUILD_MODE_FAIL_SB_NOT_IN_MY_GUILD) // 63-4
	ErrChangeGuildModeFailServerError    = newMsgError("change_guild_mode 服务器繁忙，请稍后再试", ERR_CHANGE_GUILD_MODE_FAIL_SERVER_ERROR)       // 63-5
)

// team_talk
var (
	ErrTeamTalkFailInvalidWords = newMsgError("team_talk id 错误", ERR_TEAM_TALK_FAIL_INVALID_WORDS) // 81-1
	ErrTeamTalkFailNoTeam       = newMsgError("team_talk 没有队伍", ERR_TEAM_TALK_FAIL_NO_TEAM)        // 81-2
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
