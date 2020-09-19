package mingc_war

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
	MODULE_ID = 44

	C2S_VIEW_MC_WAR_SELF_GUILD = 31

	C2S_VIEW_MC_WAR = 29

	C2S_APPLY_ATK = 16

	C2S_APPLY_AST = 21

	C2S_CANCEL_APPLY_AST = 80

	C2S_REPLY_APPLY_AST = 25

	C2S_VIEW_MINGC_WAR_MC = 75

	C2S_JOIN_FIGHT = 35

	C2S_QUIT_FIGHT = 38

	C2S_SCENE_MOVE = 49

	C2S_SCENE_BACK = 85

	C2S_SCENE_SPEED_UP = 88

	C2S_SCENE_TROOP_RELIVE = 72

	C2S_VIEW_MC_WAR_SCENE = 46

	C2S_WATCH = 139

	C2S_QUIT_WATCH = 136

	C2S_VIEW_MC_WAR_RECORD = 91

	C2S_VIEW_MC_WAR_TROOP_RECORD = 94

	C2S_VIEW_SCENE_TROOP_RECORD = 99

	C2S_APPLY_REFRESH_RANK = 107

	C2S_VIEW_MY_GUILD_MEMBER_RANK = 111

	C2S_SCENE_CHANGE_MODE = 115

	C2S_SCENE_TOU_SHI_BUILDING_TURN_TO = 119

	C2S_SCENE_TOU_SHI_BUILDING_FIRE = 123

	C2S_SCENE_DRUM = 128
)

func NewS2cViewMcWarSelfGuildMsg(self_guild *shared_proto.McWarGuildProto) pbutil.Buffer {
	msg := &S2CViewMcWarSelfGuildProto{
		SelfGuild: self_guild,
	}
	return NewS2cViewMcWarSelfGuildProtoMsg(msg)
}

var s2c_view_mc_war_self_guild = [...]byte{44, 32} // 32
func NewS2cViewMcWarSelfGuildProtoMsg(object *S2CViewMcWarSelfGuildProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_mc_war_self_guild[:], "s2c_view_mc_war_self_guild")

}

// 没有联盟
var ERR_VIEW_MC_WAR_SELF_GUILD_FAIL_NO_GUILD = pbutil.StaticBuffer{3, 44, 33, 3} // 33-3

// 服务器错误
var ERR_VIEW_MC_WAR_SELF_GUILD_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 33, 2} // 33-2

func NewS2cViewMcWarMsg(ver int32, war *shared_proto.McWarProto) pbutil.Buffer {
	msg := &S2CViewMcWarProto{
		Ver: ver,
		War: war,
	}
	return NewS2cViewMcWarProtoMsg(msg)
}

var s2c_view_mc_war = [...]byte{44, 30} // 30
func NewS2cViewMcWarProtoMsg(object *S2CViewMcWarProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_mc_war[:], "s2c_view_mc_war")

}

func NewS2cApplyAtkMsg(mcid int32, cost int32) pbutil.Buffer {
	msg := &S2CApplyAtkProto{
		Mcid: mcid,
		Cost: cost,
	}
	return NewS2cApplyAtkProtoMsg(msg)
}

var s2c_apply_atk = [...]byte{44, 17} // 17
func NewS2cApplyAtkProtoMsg(object *S2CApplyAtkProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_apply_atk[:], "s2c_apply_atk")

}

// 当前不是名城战报名时间
var ERR_APPLY_ATK_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 18, 1} // 18-1

// 无效的名城
var ERR_APPLY_ATK_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 18, 2} // 18-2

// 已经报名攻打别的名城
var ERR_APPLY_ATK_FAIL_APPLIED = pbutil.StaticBuffer{3, 44, 18, 3} // 18-3

// 城主不能攻打自己的名城
var ERR_APPLY_ATK_FAIL_IS_HOST = pbutil.StaticBuffer{3, 44, 18, 4} // 18-4

// 虎符不足
var ERR_APPLY_ATK_FAIL_HUFU_NOT_ENOUGH = pbutil.StaticBuffer{3, 44, 18, 5} // 18-5

// 虎符未达最低要求
var ERR_APPLY_ATK_FAIL_HUFU_LIMIT = pbutil.StaticBuffer{3, 44, 18, 6} // 18-6

// 盟主才能申请
var ERR_APPLY_ATK_FAIL_NOT_LEADER = pbutil.StaticBuffer{3, 44, 18, 7} // 18-7

// 联盟等级不够
var ERR_APPLY_ATK_FAIL_GUILD_LEVEL_LIMIT = pbutil.StaticBuffer{3, 44, 18, 8} // 18-8

// 本国都城没开启
var ERR_APPLY_ATK_FAIL_DU_CHENG_NOT_OPEN = pbutil.StaticBuffer{3, 44, 18, 10} // 18-10

// 服务器错误
var ERR_APPLY_ATK_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 18, 9} // 18-9

// 他国都城没开启
var ERR_APPLY_ATK_FAIL_OTHER_DU_CHENG_NOT_OPEN = pbutil.StaticBuffer{3, 44, 18, 11} // 18-11

// 他国还在占领其他名城
var ERR_APPLY_ATK_FAIL_OTHER_COUNTRY_HAS_OTHER_MC = pbutil.StaticBuffer{3, 44, 18, 12} // 18-12

// 没有占领本国都城
var ERR_APPLY_ATK_FAIL_NOT_HOLD_CAPITAL = pbutil.StaticBuffer{3, 44, 18, 14} // 18-14

func NewS2cApplyAtkSuccMsg(mcid int32) pbutil.Buffer {
	msg := &S2CApplyAtkSuccProto{
		Mcid: mcid,
	}
	return NewS2cApplyAtkSuccProtoMsg(msg)
}

var s2c_apply_atk_succ = [...]byte{44, 19} // 19
func NewS2cApplyAtkSuccProtoMsg(object *S2CApplyAtkSuccProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_apply_atk_succ[:], "s2c_apply_atk_succ")

}

func NewS2cApplyAtkFailMsg(mcid int32, cost int32) pbutil.Buffer {
	msg := &S2CApplyAtkFailProto{
		Mcid: mcid,
		Cost: cost,
	}
	return NewS2cApplyAtkFailProtoMsg(msg)
}

var s2c_apply_atk_fail = [...]byte{44, 20} // 20
func NewS2cApplyAtkFailProtoMsg(object *S2CApplyAtkFailProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_apply_atk_fail[:], "s2c_apply_atk_fail")

}

func NewS2cApplyAstMsg(mcid int32, atk bool) pbutil.Buffer {
	msg := &S2CApplyAstProto{
		Mcid: mcid,
		Atk:  atk,
	}
	return NewS2cApplyAstProtoMsg(msg)
}

var s2c_apply_ast = [...]byte{44, 22} // 22
func NewS2cApplyAstProtoMsg(object *S2CApplyAstProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_apply_ast[:], "s2c_apply_ast")

}

// 当前不是名城战报名时间
var ERR_APPLY_AST_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 23, 1} // 23-1

// 无效的名城
var ERR_APPLY_AST_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 23, 2} // 23-2

// 城主不能攻打自己的名城
var ERR_APPLY_AST_FAIL_IS_HOST = pbutil.StaticBuffer{3, 44, 23, 3} // 23-3

// 盟主才能申请
var ERR_APPLY_AST_FAIL_NOT_LEADER = pbutil.StaticBuffer{3, 44, 23, 4} // 23-4

// 这座城不能协助
var ERR_APPLY_AST_FAIL_MCID_CANNOT_AST = pbutil.StaticBuffer{3, 44, 23, 5} // 23-5

// 这座城协助已达上限
var ERR_APPLY_AST_FAIL_MC_AST_LIMIT = pbutil.StaticBuffer{3, 44, 23, 6} // 23-6

// 达到申请上限
var ERR_APPLY_AST_FAIL_APPLIED_LIMIT = pbutil.StaticBuffer{3, 44, 23, 7} // 23-7

// 已经申请
var ERR_APPLY_AST_FAIL_ALREADY_AST = pbutil.StaticBuffer{3, 44, 23, 9} // 23-9

// 服务器错误
var ERR_APPLY_AST_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 23, 8} // 23-8

func NewS2cReceiveApplyAstMsg(mcid int32) pbutil.Buffer {
	msg := &S2CReceiveApplyAstProto{
		Mcid: mcid,
	}
	return NewS2cReceiveApplyAstProtoMsg(msg)
}

var s2c_receive_apply_ast = [...]byte{44, 24} // 24
func NewS2cReceiveApplyAstProtoMsg(object *S2CReceiveApplyAstProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_receive_apply_ast[:], "s2c_receive_apply_ast")

}

func NewS2cCancelApplyAstMsg(mcid int32) pbutil.Buffer {
	msg := &S2CCancelApplyAstProto{
		Mcid: mcid,
	}
	return NewS2cCancelApplyAstProtoMsg(msg)
}

var s2c_cancel_apply_ast = [...]byte{44, 81} // 81
func NewS2cCancelApplyAstProtoMsg(object *S2CCancelApplyAstProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_cancel_apply_ast[:], "s2c_cancel_apply_ast")

}

// 当前不是名城战报名时间
var ERR_CANCEL_APPLY_AST_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 82, 1} // 82-1

// 无效的名城
var ERR_CANCEL_APPLY_AST_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 82, 2} // 82-2

// 没有申请
var ERR_CANCEL_APPLY_AST_FAIL_NOT_APPLY = pbutil.StaticBuffer{3, 44, 82, 3} // 82-3

// 盟主才能操作
var ERR_CANCEL_APPLY_AST_FAIL_NOT_LEADER = pbutil.StaticBuffer{3, 44, 82, 5} // 82-5

// 服务器错误
var ERR_CANCEL_APPLY_AST_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 82, 4} // 82-4

func NewS2cReceiveCancelApplyAstMsg(mcid int32) pbutil.Buffer {
	msg := &S2CReceiveCancelApplyAstProto{
		Mcid: mcid,
	}
	return NewS2cReceiveCancelApplyAstProtoMsg(msg)
}

var s2c_receive_cancel_apply_ast = [...]byte{44, 83} // 83
func NewS2cReceiveCancelApplyAstProtoMsg(object *S2CReceiveCancelApplyAstProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_receive_cancel_apply_ast[:], "s2c_receive_cancel_apply_ast")

}

func NewS2cReplyApplyAstMsg(mcid int32, gid int32, agree bool) pbutil.Buffer {
	msg := &S2CReplyApplyAstProto{
		Mcid:  mcid,
		Gid:   gid,
		Agree: agree,
	}
	return NewS2cReplyApplyAstProtoMsg(msg)
}

var s2c_reply_apply_ast = [...]byte{44, 26} // 26
func NewS2cReplyApplyAstProtoMsg(object *S2CReplyApplyAstProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_reply_apply_ast[:], "s2c_reply_apply_ast")

}

// 当前不是名城战报名时间
var ERR_REPLY_APPLY_AST_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 27, 1} // 27-1

// 无效的名城
var ERR_REPLY_APPLY_AST_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 27, 2} // 27-2

// 盟主才能审批
var ERR_REPLY_APPLY_AST_FAIL_NOT_LEADER = pbutil.StaticBuffer{3, 44, 27, 8} // 27-8

// 无权审批
var ERR_REPLY_APPLY_AST_FAIL_REPLY_PERMISSION_DENIED = pbutil.StaticBuffer{3, 44, 27, 6} // 27-6

// 联盟id 错误
var ERR_REPLY_APPLY_AST_FAIL_NOT_APPLY = pbutil.StaticBuffer{3, 44, 27, 7} // 27-7

// 协助的联盟数已经到最大了
var ERR_REPLY_APPLY_AST_FAIL_AST_LIMIT = pbutil.StaticBuffer{3, 44, 27, 3} // 27-3

// 对方协助的名城数已经最大了
var ERR_REPLY_APPLY_AST_FAIL_TARGET_AST_LIMIT = pbutil.StaticBuffer{3, 44, 27, 4} // 27-4

// 对方已经协助这座名城了
var ERR_REPLY_APPLY_AST_FAIL_TARGET_ALREADY_AST = pbutil.StaticBuffer{3, 44, 27, 9} // 27-9

// 服务器错误
var ERR_REPLY_APPLY_AST_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 27, 5} // 27-5

func NewS2cApplyAstPassMsg(mcid int32, agree bool) pbutil.Buffer {
	msg := &S2CApplyAstPassProto{
		Mcid:  mcid,
		Agree: agree,
	}
	return NewS2cApplyAstPassProtoMsg(msg)
}

var s2c_apply_ast_pass = [...]byte{44, 28} // 28
func NewS2cApplyAstPassProtoMsg(object *S2CApplyAstPassProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_apply_ast_pass[:], "s2c_apply_ast_pass")

}

func NewS2cMingcWarFightPrepareStartMsg(start_time int32, end_time int32) pbutil.Buffer {
	msg := &S2CMingcWarFightPrepareStartProto{
		StartTime: start_time,
		EndTime:   end_time,
	}
	return NewS2cMingcWarFightPrepareStartProtoMsg(msg)
}

var s2c_mingc_war_fight_prepare_start = [...]byte{44, 104} // 104
func NewS2cMingcWarFightPrepareStartProtoMsg(object *S2CMingcWarFightPrepareStartProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_mingc_war_fight_prepare_start[:], "s2c_mingc_war_fight_prepare_start")

}

func NewS2cMingcWarFightStartMsg(start_time int32, end_time int32) pbutil.Buffer {
	msg := &S2CMingcWarFightStartProto{
		StartTime: start_time,
		EndTime:   end_time,
	}
	return NewS2cMingcWarFightStartProtoMsg(msg)
}

var s2c_mingc_war_fight_start = [...]byte{44, 41} // 41
func NewS2cMingcWarFightStartProtoMsg(object *S2CMingcWarFightStartProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_mingc_war_fight_start[:], "s2c_mingc_war_fight_start")

}

func NewS2cIsJoiningFightOnLoginMsg(mcid int32) pbutil.Buffer {
	msg := &S2CIsJoiningFightOnLoginProto{
		Mcid: mcid,
	}
	return NewS2cIsJoiningFightOnLoginProtoMsg(msg)
}

var s2c_is_joining_fight_on_login = [...]byte{44, 66} // 66
func NewS2cIsJoiningFightOnLoginProtoMsg(object *S2CIsJoiningFightOnLoginProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_is_joining_fight_on_login[:], "s2c_is_joining_fight_on_login")

}

var RED_POINT_NOTICE_S2C = pbutil.StaticBuffer{2, 44, 103} // 103

func NewS2cViewMingcWarMcMsg(mc *shared_proto.McWarMcProto) pbutil.Buffer {
	msg := &S2CViewMingcWarMcProto{
		Mc: mc,
	}
	return NewS2cViewMingcWarMcProtoMsg(msg)
}

var s2c_view_mingc_war_mc = [...]byte{44, 76} // 76
func NewS2cViewMingcWarMcProtoMsg(object *S2CViewMingcWarMcProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_mingc_war_mc[:], "s2c_view_mingc_war_mc")

}

// 无效的名城
var ERR_VIEW_MINGC_WAR_MC_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 77, 1} // 77-1

// 服务器错误
var ERR_VIEW_MINGC_WAR_MC_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 77, 2} // 77-2

func NewS2cJoinFightMsg(mcid int32, atk bool) pbutil.Buffer {
	msg := &S2CJoinFightProto{
		Mcid: mcid,
		Atk:  atk,
	}
	return NewS2cJoinFightProtoMsg(msg)
}

var s2c_join_fight = [...]byte{44, 36} // 36
func NewS2cJoinFightProtoMsg(object *S2CJoinFightProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_join_fight[:], "s2c_join_fight")

}

// 当前不是名城战战斗时间
var ERR_JOIN_FIGHT_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 37, 1} // 37-1

// 无效的名城
var ERR_JOIN_FIGHT_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 37, 2} // 37-2

// 这座城的名城战已经结束了
var ERR_JOIN_FIGHT_FAIL_MC_WAR_END = pbutil.StaticBuffer{3, 44, 37, 7} // 37-7

// 没有申请攻打/协助名城或不是占领盟
var ERR_JOIN_FIGHT_FAIL_NOT_APPLY = pbutil.StaticBuffer{3, 44, 37, 3} // 37-3

// 已经在其他名城参战了
var ERR_JOIN_FIGHT_FAIL_ALREADY_IN_WAR = pbutil.StaticBuffer{3, 44, 37, 4} // 37-4

// 武将不存在
var ERR_JOIN_FIGHT_FAIL_INVALID_CAPTAIN_ID = pbutil.StaticBuffer{3, 44, 37, 6} // 37-6

// 战斗开始后加入不让进
var ERR_JOIN_FIGHT_FAIL_JOIN_GUILD_TOO_LATE = pbutil.StaticBuffer{3, 44, 37, 8} // 37-8

// 君主等级不够
var ERR_JOIN_FIGHT_FAIL_HERO_LEVEL_LIMIT = pbutil.StaticBuffer{3, 44, 37, 9} // 37-9

// 服务器错误
var ERR_JOIN_FIGHT_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 37, 5} // 37-5

func NewS2cOtherJoinFightMsg(mcid int32, hero_id []byte, data *shared_proto.McWarTroopProto) pbutil.Buffer {
	msg := &S2COtherJoinFightProto{
		Mcid:   mcid,
		HeroId: hero_id,
		Data:   data,
	}
	return NewS2cOtherJoinFightProtoMsg(msg)
}

var s2c_other_join_fight = [...]byte{44, 58} // 58
func NewS2cOtherJoinFightProtoMsg(object *S2COtherJoinFightProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_other_join_fight[:], "s2c_other_join_fight")

}

func NewS2cQuitFightMsg(mcid int32) pbutil.Buffer {
	msg := &S2CQuitFightProto{
		Mcid: mcid,
	}
	return NewS2cQuitFightProtoMsg(msg)
}

var s2c_quit_fight = [...]byte{44, 39} // 39
func NewS2cQuitFightProtoMsg(object *S2CQuitFightProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_quit_fight[:], "s2c_quit_fight")

}

// 当前不是名城战战斗时间
var ERR_QUIT_FIGHT_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 40, 1} // 40-1

// 没有参战
var ERR_QUIT_FIGHT_FAIL_NOT_JOIN_FIGHT = pbutil.StaticBuffer{3, 44, 40, 2} // 40-2

// 服务器错误
var ERR_QUIT_FIGHT_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 40, 3} // 40-3

func NewS2cOtherQuitFightMsg(mcid int32, hero_id []byte) pbutil.Buffer {
	msg := &S2COtherQuitFightProto{
		Mcid:   mcid,
		HeroId: hero_id,
	}
	return NewS2cOtherQuitFightProtoMsg(msg)
}

var s2c_other_quit_fight = [...]byte{44, 57} // 57
func NewS2cOtherQuitFightProtoMsg(object *S2COtherQuitFightProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_other_quit_fight[:], "s2c_other_quit_fight")

}

func NewS2cSceneBuildingDestroyProsperityMsg(pos_x int32, pos_y int32, new_prosperity int32) pbutil.Buffer {
	msg := &S2CSceneBuildingDestroyProsperityProto{
		PosX:          pos_x,
		PosY:          pos_y,
		NewProsperity: new_prosperity,
	}
	return NewS2cSceneBuildingDestroyProsperityProtoMsg(msg)
}

var s2c_scene_building_destroy_prosperity = [...]byte{44, 79} // 79
func NewS2cSceneBuildingDestroyProsperityProtoMsg(object *S2CSceneBuildingDestroyProsperityProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_building_destroy_prosperity[:], "s2c_scene_building_destroy_prosperity")

}

func NewS2cSceneFightPrepareEndMsg(start_time int32) pbutil.Buffer {
	msg := &S2CSceneFightPrepareEndProto{
		StartTime: start_time,
	}
	return NewS2cSceneFightPrepareEndProtoMsg(msg)
}

var s2c_scene_fight_prepare_end = [...]byte{44, 78} // 78
func NewS2cSceneFightPrepareEndProtoMsg(object *S2CSceneFightPrepareEndProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_fight_prepare_end[:], "s2c_scene_fight_prepare_end")

}

func NewS2cSceneWarEndMsg(mc_id int32) pbutil.Buffer {
	msg := &S2CSceneWarEndProto{
		McId: mc_id,
	}
	return NewS2cSceneWarEndProtoMsg(msg)
}

var s2c_scene_war_end = [...]byte{44, 64} // 64
func NewS2cSceneWarEndProtoMsg(object *S2CSceneWarEndProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_war_end[:], "s2c_scene_war_end")

}

func NewS2cSceneMoveMsg(dest_pos_x int32, dest_pos_y int32, end_time int32) pbutil.Buffer {
	msg := &S2CSceneMoveProto{
		DestPosX: dest_pos_x,
		DestPosY: dest_pos_y,
		EndTime:  end_time,
	}
	return NewS2cSceneMoveProtoMsg(msg)
}

var s2c_scene_move = [...]byte{44, 50} // 50
func NewS2cSceneMoveProtoMsg(object *S2CSceneMoveProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_move[:], "s2c_scene_move")

}

// 当前不是名城战战斗时间
var ERR_SCENE_MOVE_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 51, 1} // 51-1

// 无效的名城
var ERR_SCENE_MOVE_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 51, 2} // 51-2

// 没有参战
var ERR_SCENE_MOVE_FAIL_NOT_IN_SCENE = pbutil.StaticBuffer{3, 44, 51, 6} // 51-6

// 这座城的名城战已经结束了
var ERR_SCENE_MOVE_FAIL_MC_WAR_END = pbutil.StaticBuffer{3, 44, 51, 11} // 51-11

// 在准备时间
var ERR_SCENE_MOVE_FAIL_IN_PREPARE_DURATION = pbutil.StaticBuffer{3, 44, 51, 9} // 51-9

// 在入场时间
var ERR_SCENE_MOVE_FAIL_IN_JOIN_DURATION = pbutil.StaticBuffer{3, 44, 51, 10} // 51-10

// 不在驻扎状态
var ERR_SCENE_MOVE_FAIL_NOT_STATION = pbutil.StaticBuffer{3, 44, 51, 8} // 51-8

// 已经在目的地了
var ERR_SCENE_MOVE_FAIL_ALREADY_ON_DEST_POS = pbutil.StaticBuffer{3, 44, 51, 4} // 51-4

// 服务器错误
var ERR_SCENE_MOVE_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 51, 5} // 51-5

// 目的地不能直接到达
var ERR_SCENE_MOVE_FAIL_DEST_CANNOT_ARRIVE = pbutil.StaticBuffer{3, 44, 51, 12} // 51-12

// 当前占领据点没有被摧毁
var ERR_SCENE_MOVE_FAIL_NOT_DESTROY = pbutil.StaticBuffer{3, 44, 51, 13} // 51-13

// 在行军中
var ERR_SCENE_MOVE_FAIL_IS_MOVING = pbutil.StaticBuffer{3, 44, 51, 14} // 51-14

// 在补兵中
var ERR_SCENE_MOVE_FAIL_IS_RELIVE = pbutil.StaticBuffer{3, 44, 51, 15} // 51-15

func NewS2cSceneBackMsg(dest_pos_x int32, dest_pos_y int32, end_time int32) pbutil.Buffer {
	msg := &S2CSceneBackProto{
		DestPosX: dest_pos_x,
		DestPosY: dest_pos_y,
		EndTime:  end_time,
	}
	return NewS2cSceneBackProtoMsg(msg)
}

var s2c_scene_back = [...]byte{44, 86} // 86
func NewS2cSceneBackProtoMsg(object *S2CSceneBackProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_back[:], "s2c_scene_back")

}

// 当前不是名城战战斗时间
var ERR_SCENE_BACK_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 87, 1} // 87-1

// 无效的名城
var ERR_SCENE_BACK_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 87, 2} // 87-2

// 没有参战
var ERR_SCENE_BACK_FAIL_NOT_IN_SCENE = pbutil.StaticBuffer{3, 44, 87, 3} // 87-3

// 这座城的名城战已经结束了
var ERR_SCENE_BACK_FAIL_MC_WAR_END = pbutil.StaticBuffer{3, 44, 87, 4} // 87-4

// 在准备时间
var ERR_SCENE_BACK_FAIL_IN_PREPARE_DURATION = pbutil.StaticBuffer{3, 44, 87, 5} // 87-5

// 在入场时间
var ERR_SCENE_BACK_FAIL_IN_JOIN_DURATION = pbutil.StaticBuffer{3, 44, 87, 6} // 87-6

// 不是移动状态
var ERR_SCENE_BACK_FAIL_NOT_MOVE = pbutil.StaticBuffer{3, 44, 87, 7} // 87-7

// 服务器错误
var ERR_SCENE_BACK_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 87, 8} // 87-8

func NewS2cSceneSpeedUpMsg(id []byte, end_time int32) pbutil.Buffer {
	msg := &S2CSceneSpeedUpProto{
		Id:      id,
		EndTime: end_time,
	}
	return NewS2cSceneSpeedUpProtoMsg(msg)
}

var s2c_scene_speed_up = [...]byte{44, 89} // 89
func NewS2cSceneSpeedUpProtoMsg(object *S2CSceneSpeedUpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_speed_up[:], "s2c_scene_speed_up")

}

// 发送的不是行军加速道具
var ERR_SCENE_SPEED_UP_FAIL_INVALID_GOODS = pbutil.StaticBuffer{3, 44, 90, 1} // 90-1

// 物品个数不足
var ERR_SCENE_SPEED_UP_FAIL_GOODS_NOT_ENOUGH = pbutil.StaticBuffer{3, 44, 90, 2} // 90-2

// 不支持点券购买
var ERR_SCENE_SPEED_UP_FAIL_COST_NOT_SUPPORT = pbutil.StaticBuffer{3, 44, 90, 3} // 90-3

// 点券购买，点券不足
var ERR_SCENE_SPEED_UP_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 44, 90, 4} // 90-4

// 当前不是名城战战斗时间
var ERR_SCENE_SPEED_UP_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 90, 5} // 90-5

// 无效的名城
var ERR_SCENE_SPEED_UP_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 90, 6} // 90-6

// 没有参战
var ERR_SCENE_SPEED_UP_FAIL_NOT_IN_SCENE = pbutil.StaticBuffer{3, 44, 90, 7} // 90-7

// 这座城的名城战已经结束了
var ERR_SCENE_SPEED_UP_FAIL_MC_WAR_END = pbutil.StaticBuffer{3, 44, 90, 8} // 90-8

// 在准备时间
var ERR_SCENE_SPEED_UP_FAIL_IN_PREPARE_DURATION = pbutil.StaticBuffer{3, 44, 90, 9} // 90-9

// 在入场时间
var ERR_SCENE_SPEED_UP_FAIL_IN_JOIN_DURATION = pbutil.StaticBuffer{3, 44, 90, 10} // 90-10

// 部队不是行军中，不能加速
var ERR_SCENE_SPEED_UP_FAIL_NO_MOVING = pbutil.StaticBuffer{3, 44, 90, 11} // 90-11

// 服务器忙，请稍后再试
var ERR_SCENE_SPEED_UP_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 44, 90, 12} // 90-12

func NewS2cSceneOtherMoveMsg(hero_id []byte, start_pos_x int32, start_pos_y int32, dest_pos_x int32, dest_pos_y int32, end_time int32) pbutil.Buffer {
	msg := &S2CSceneOtherMoveProto{
		HeroId:    hero_id,
		StartPosX: start_pos_x,
		StartPosY: start_pos_y,
		DestPosX:  dest_pos_x,
		DestPosY:  dest_pos_y,
		EndTime:   end_time,
	}
	return NewS2cSceneOtherMoveProtoMsg(msg)
}

var s2c_scene_other_move = [...]byte{44, 52} // 52
func NewS2cSceneOtherMoveProtoMsg(object *S2CSceneOtherMoveProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_other_move[:], "s2c_scene_other_move")

}

func NewS2cSceneMoveStationMsg(hero_id []byte, pos_x int32, pos_y int32) pbutil.Buffer {
	msg := &S2CSceneMoveStationProto{
		HeroId: hero_id,
		PosX:   pos_x,
		PosY:   pos_y,
	}
	return NewS2cSceneMoveStationProtoMsg(msg)
}

var s2c_scene_move_station = [...]byte{44, 62} // 62
func NewS2cSceneMoveStationProtoMsg(object *S2CSceneMoveStationProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_move_station[:], "s2c_scene_move_station")

}

func NewS2cSceneBuildingFightMsg(hero_id []byte, gid int32, win bool, pos_x int32, pos_y int32) pbutil.Buffer {
	msg := &S2CSceneBuildingFightProto{
		HeroId: hero_id,
		Gid:    gid,
		Win:    win,
		PosX:   pos_x,
		PosY:   pos_y,
	}
	return NewS2cSceneBuildingFightProtoMsg(msg)
}

var s2c_scene_building_fight = [...]byte{44, 63} // 63
func NewS2cSceneBuildingFightProtoMsg(object *S2CSceneBuildingFightProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_building_fight[:], "s2c_scene_building_fight")

}

func NewS2cSceneTroopReliveMsg(end_time int32) pbutil.Buffer {
	msg := &S2CSceneTroopReliveProto{
		EndTime: end_time,
	}
	return NewS2cSceneTroopReliveProtoMsg(msg)
}

var s2c_scene_troop_relive = [...]byte{44, 54} // 54
func NewS2cSceneTroopReliveProtoMsg(object *S2CSceneTroopReliveProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_troop_relive[:], "s2c_scene_troop_relive")

}

// 当前不是名城战战斗时间
var ERR_SCENE_TROOP_RELIVE_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 73, 1} // 73-1

// 无效的名城
var ERR_SCENE_TROOP_RELIVE_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 73, 2} // 73-2

// 没有参战
var ERR_SCENE_TROOP_RELIVE_FAIL_NOT_IN_SCENE = pbutil.StaticBuffer{3, 44, 73, 3} // 73-3

// 这座城的名城战已经结束了
var ERR_SCENE_TROOP_RELIVE_FAIL_MC_WAR_END = pbutil.StaticBuffer{3, 44, 73, 4} // 73-4

// 在准备时间
var ERR_SCENE_TROOP_RELIVE_FAIL_IN_PREPARE_DURATION = pbutil.StaticBuffer{3, 44, 73, 5} // 73-5

// 在入场时间
var ERR_SCENE_TROOP_RELIVE_FAIL_IN_JOIN_DURATION = pbutil.StaticBuffer{3, 44, 73, 6} // 73-6

// 不是驻扎状态
var ERR_SCENE_TROOP_RELIVE_FAIL_NOT_STATION = pbutil.StaticBuffer{3, 44, 73, 7} // 73-7

// 不在复活点
var ERR_SCENE_TROOP_RELIVE_FAIL_NOT_IN_RELIVE_POS = pbutil.StaticBuffer{3, 44, 73, 10} // 73-10

// 兵是满的
var ERR_SCENE_TROOP_RELIVE_FAIL_FULL_SOLIDER = pbutil.StaticBuffer{3, 44, 73, 8} // 73-8

// 服务器错误
var ERR_SCENE_TROOP_RELIVE_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 73, 9} // 73-9

func NewS2cSceneOtherTroopReliveMsg(hero_id []byte, end_time int32) pbutil.Buffer {
	msg := &S2CSceneOtherTroopReliveProto{
		HeroId:  hero_id,
		EndTime: end_time,
	}
	return NewS2cSceneOtherTroopReliveProtoMsg(msg)
}

var s2c_scene_other_troop_relive = [...]byte{44, 74} // 74
func NewS2cSceneOtherTroopReliveProtoMsg(object *S2CSceneOtherTroopReliveProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_other_troop_relive[:], "s2c_scene_other_troop_relive")

}

func NewS2cSceneTroopUpdateMsg(hero_id []byte, data *shared_proto.McWarTroopProto) pbutil.Buffer {
	msg := &S2CSceneTroopUpdateProto{
		HeroId: hero_id,
		Data:   data,
	}
	return NewS2cSceneTroopUpdateProtoMsg(msg)
}

var s2c_scene_troop_update = [...]byte{44, 56} // 56
func NewS2cSceneTroopUpdateProtoMsg(object *S2CSceneTroopUpdateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_troop_update[:], "s2c_scene_troop_update")

}

func NewS2cViewMcWarSceneMsg(scene *shared_proto.McWarSceneProto) pbutil.Buffer {
	msg := &S2CViewMcWarSceneProto{
		Scene: scene,
	}
	return NewS2cViewMcWarSceneProtoMsg(msg)
}

var s2c_view_mc_war_scene = [...]byte{44, 47} // 47
func NewS2cViewMcWarSceneProtoMsg(object *S2CViewMcWarSceneProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_view_mc_war_scene[:], "s2c_view_mc_war_scene")

}

// 当前不是名城战战斗时间
var ERR_VIEW_MC_WAR_SCENE_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 48, 1} // 48-1

// 无效的名城
var ERR_VIEW_MC_WAR_SCENE_FAIL_INVALID_MCID = pbutil.StaticBuffer{3, 44, 48, 2} // 48-2

// 服务器错误
var ERR_VIEW_MC_WAR_SCENE_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 48, 3} // 48-3

var WATCH_S2C = pbutil.StaticBuffer{3, 44, 140, 1} // 140

// 当前不是名城战战斗时间
var ERR_WATCH_FAIL_INVALID_TIME = pbutil.StaticBuffer{4, 44, 141, 1, 1} // 141-1

// 无效的名城
var ERR_WATCH_FAIL_INVALID_MCID = pbutil.StaticBuffer{4, 44, 141, 1, 2} // 141-2

// 服务器错误
var ERR_WATCH_FAIL_SERVER_ERR = pbutil.StaticBuffer{4, 44, 141, 1, 3} // 141-3

var QUIT_WATCH_S2C = pbutil.StaticBuffer{3, 44, 137, 1} // 137

// 当前不是名城战战斗时间
var ERR_QUIT_WATCH_FAIL_INVALID_TIME = pbutil.StaticBuffer{4, 44, 138, 1, 2} // 138-2

// 无效的名城
var ERR_QUIT_WATCH_FAIL_INVALID_MCID = pbutil.StaticBuffer{4, 44, 138, 1, 3} // 138-3

// 服务器错误
var ERR_QUIT_WATCH_FAIL_SERVER_ERR = pbutil.StaticBuffer{4, 44, 138, 1, 4} // 138-4

func NewS2cMcWarEndRecordMsg(war_id int32, record *shared_proto.McWarFightRecordProto) pbutil.Buffer {
	msg := &S2CMcWarEndRecordProto{
		WarId:  war_id,
		Record: record,
	}
	return NewS2cMcWarEndRecordProtoMsg(msg)
}

var s2c_mc_war_end_record = [...]byte{44, 67} // 67
func NewS2cMcWarEndRecordProtoMsg(object *S2CMcWarEndRecordProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_mc_war_end_record[:], "s2c_mc_war_end_record")

}

func NewS2cViewMcWarRecordMsg(war_id int32, record *shared_proto.McWarFightRecordProto) pbutil.Buffer {
	msg := &S2CViewMcWarRecordProto{
		WarId:  war_id,
		Record: record,
	}
	return NewS2cViewMcWarRecordProtoMsg(msg)
}

var s2c_view_mc_war_record = [...]byte{44, 92} // 92
func NewS2cViewMcWarRecordProtoMsg(object *S2CViewMcWarRecordProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_mc_war_record[:], "s2c_view_mc_war_record")

}

// 找不到记录
var ERR_VIEW_MC_WAR_RECORD_FAIL_NO_RECORD = pbutil.StaticBuffer{3, 44, 93, 1} // 93-1

func NewS2cViewMcWarTroopRecordMsg(record *shared_proto.McWarTroopAllRecordProto) pbutil.Buffer {
	msg := &S2CViewMcWarTroopRecordProto{
		Record: record,
	}
	return NewS2cViewMcWarTroopRecordProtoMsg(msg)
}

var s2c_view_mc_war_troop_record = [...]byte{44, 95} // 95
func NewS2cViewMcWarTroopRecordProtoMsg(object *S2CViewMcWarTroopRecordProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_mc_war_troop_record[:], "s2c_view_mc_war_troop_record")

}

// 找不到记录
var ERR_VIEW_MC_WAR_TROOP_RECORD_FAIL_NO_RECORD = pbutil.StaticBuffer{3, 44, 96, 1} // 96-1

func NewS2cViewSceneTroopRecordMsg(record [][]byte) pbutil.Buffer {
	msg := &S2CViewSceneTroopRecordProto{
		Record: record,
	}
	return NewS2cViewSceneTroopRecordProtoMsg(msg)
}

var s2c_view_scene_troop_record = [...]byte{44, 100} // 100
func NewS2cViewSceneTroopRecordProtoMsg(object *S2CViewSceneTroopRecordProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_scene_troop_record[:], "s2c_view_scene_troop_record")

}

// 当前不是名城战战斗时间
var ERR_VIEW_SCENE_TROOP_RECORD_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 101, 1} // 101-1

// 没有参战
var ERR_VIEW_SCENE_TROOP_RECORD_FAIL_NOT_IN_SCENE = pbutil.StaticBuffer{3, 44, 101, 3} // 101-3

// 这座城的名城战已经结束了
var ERR_VIEW_SCENE_TROOP_RECORD_FAIL_MC_WAR_END = pbutil.StaticBuffer{3, 44, 101, 4} // 101-4

// 服务器错误
var ERR_VIEW_SCENE_TROOP_RECORD_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 101, 5} // 101-5

var SCENE_TROOP_RECORD_ADD_NOTICE_S2C = pbutil.StaticBuffer{2, 44, 105} // 105

func NewS2cMyRankMsg(info *shared_proto.McWarTroopRankProto) pbutil.Buffer {
	msg := &S2CMyRankProto{
		Info: info,
	}
	return NewS2cMyRankProtoMsg(msg)
}

var s2c_my_rank = [...]byte{44, 110} // 110
func NewS2cMyRankProtoMsg(object *S2CMyRankProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_my_rank[:], "s2c_my_rank")

}

func NewS2cApplyRefreshRankMsg(version int32, rank *shared_proto.McWarTroopsRankProto) pbutil.Buffer {
	msg := &S2CApplyRefreshRankProto{
		Version: version,
		Rank:    rank,
	}
	return NewS2cApplyRefreshRankProtoMsg(msg)
}

var s2c_apply_refresh_rank = [...]byte{44, 108} // 108
func NewS2cApplyRefreshRankProtoMsg(object *S2CApplyRefreshRankProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_apply_refresh_rank[:], "s2c_apply_refresh_rank")

}

// 服务器错误
var ERR_APPLY_REFRESH_RANK_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 109, 1} // 109-1

// 当前不是名城战战斗时间
var ERR_APPLY_REFRESH_RANK_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 109, 2} // 109-2

// 没有参战
var ERR_APPLY_REFRESH_RANK_FAIL_NOT_IN_SCENE = pbutil.StaticBuffer{3, 44, 109, 3} // 109-3

func NewS2cViewMyGuildMemberRankMsg(rank *shared_proto.McWarTroopsInfoProto) pbutil.Buffer {
	msg := &S2CViewMyGuildMemberRankProto{
		Rank: rank,
	}
	return NewS2cViewMyGuildMemberRankProtoMsg(msg)
}

var s2c_view_my_guild_member_rank = [...]byte{44, 112} // 112
func NewS2cViewMyGuildMemberRankProtoMsg(object *S2CViewMyGuildMemberRankProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_my_guild_member_rank[:], "s2c_view_my_guild_member_rank")

}

// 找不到记录
var ERR_VIEW_MY_GUILD_MEMBER_RANK_FAIL_NO_RECORD = pbutil.StaticBuffer{3, 44, 113, 1} // 113-1

func NewS2cCurMultiKillMsg(multi_kill int32) pbutil.Buffer {
	msg := &S2CCurMultiKillProto{
		MultiKill: multi_kill,
	}
	return NewS2cCurMultiKillProtoMsg(msg)
}

var s2c_cur_multi_kill = [...]byte{44, 114} // 114
func NewS2cCurMultiKillProtoMsg(object *S2CCurMultiKillProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_cur_multi_kill[:], "s2c_cur_multi_kill")

}

func NewS2cSpecialMultiKillMsg(hero_id []byte, multi_kill int32) pbutil.Buffer {
	msg := &S2CSpecialMultiKillProto{
		HeroId:    hero_id,
		MultiKill: multi_kill,
	}
	return NewS2cSpecialMultiKillProtoMsg(msg)
}

var s2c_special_multi_kill = [...]byte{44, 135, 1} // 135
func NewS2cSpecialMultiKillProtoMsg(object *S2CSpecialMultiKillProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_special_multi_kill[:], "s2c_special_multi_kill")

}

func NewS2cSceneChangeModeMsg(mode int32) pbutil.Buffer {
	msg := &S2CSceneChangeModeProto{
		Mode: mode,
	}
	return NewS2cSceneChangeModeProtoMsg(msg)
}

var s2c_scene_change_mode = [...]byte{44, 116} // 116
func NewS2cSceneChangeModeProtoMsg(object *S2CSceneChangeModeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_change_mode[:], "s2c_scene_change_mode")

}

// mode 不存在
var ERR_SCENE_CHANGE_MODE_FAIL_INVALID_MODE = pbutil.StaticBuffer{3, 44, 117, 9} // 117-9

// 当前不是名城战战斗时间
var ERR_SCENE_CHANGE_MODE_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 117, 5} // 117-5

// 没有参战
var ERR_SCENE_CHANGE_MODE_FAIL_NOT_IN_SCENE = pbutil.StaticBuffer{3, 44, 117, 7} // 117-7

// 这座城的名城战已经结束了
var ERR_SCENE_CHANGE_MODE_FAIL_MC_WAR_END = pbutil.StaticBuffer{3, 44, 117, 8} // 117-8

// 已经是新形态了
var ERR_SCENE_CHANGE_MODE_FAIL_SAME_MODE = pbutil.StaticBuffer{3, 44, 117, 1} // 117-1

// 不是驻扎状态
var ERR_SCENE_CHANGE_MODE_FAIL_NOT_STATION = pbutil.StaticBuffer{3, 44, 117, 2} // 117-2

// 不在复活点
var ERR_SCENE_CHANGE_MODE_FAIL_NOT_IN_RELIVE_POS = pbutil.StaticBuffer{3, 44, 117, 3} // 117-3

// 服务器错误
var ERR_SCENE_CHANGE_MODE_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 117, 4} // 117-4

func NewS2cSceneChangeModeNoticeMsg(hero_id []byte, mode int32) pbutil.Buffer {
	msg := &S2CSceneChangeModeNoticeProto{
		HeroId: hero_id,
		Mode:   mode,
	}
	return NewS2cSceneChangeModeNoticeProtoMsg(msg)
}

var s2c_scene_change_mode_notice = [...]byte{44, 133, 1} // 133
func NewS2cSceneChangeModeNoticeProtoMsg(object *S2CSceneChangeModeNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_change_mode_notice[:], "s2c_scene_change_mode_notice")

}

func NewS2cSceneTouShiBuildingTurnToMsg(pos_x int32, pos_y int32, left bool, new_target_index int32, turn_end_time int32) pbutil.Buffer {
	msg := &S2CSceneTouShiBuildingTurnToProto{
		PosX:           pos_x,
		PosY:           pos_y,
		Left:           left,
		NewTargetIndex: new_target_index,
		TurnEndTime:    turn_end_time,
	}
	return NewS2cSceneTouShiBuildingTurnToProtoMsg(msg)
}

var s2c_scene_tou_shi_building_turn_to = [...]byte{44, 120} // 120
func NewS2cSceneTouShiBuildingTurnToProtoMsg(object *S2CSceneTouShiBuildingTurnToProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_tou_shi_building_turn_to[:], "s2c_scene_tou_shi_building_turn_to")

}

// 当前不是名城战战斗时间
var ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 121, 6} // 121-6

// 没有参战
var ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_NOT_IN_SCENE = pbutil.StaticBuffer{3, 44, 121, 7} // 121-7

// 这座城的名城战已经结束了
var ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_MC_WAR_END = pbutil.StaticBuffer{3, 44, 121, 8} // 121-8

// 坐标错误或不是投石机
var ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_INVALID_POS = pbutil.StaticBuffer{3, 44, 121, 1} // 121-1

// 不是占领者
var ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_NOT_HOST = pbutil.StaticBuffer{3, 44, 121, 2} // 121-2

// 没有目标
var ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_NO_TARGET = pbutil.StaticBuffer{3, 44, 121, 3} // 121-3

// 正在转向中
var ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_IN_TURN_CD = pbutil.StaticBuffer{3, 44, 121, 4} // 121-4

// 正在装填中
var ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_IN_PREPARE_CD = pbutil.StaticBuffer{3, 44, 121, 5} // 121-5

// 服务器错误
var ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 121, 9} // 121-9

func NewS2cSceneTouShiBuildingTurnToNoticeMsg(pos_x int32, pos_y int32, left bool, new_target_index int32, turn_end_time int32) pbutil.Buffer {
	msg := &S2CSceneTouShiBuildingTurnToNoticeProto{
		PosX:           pos_x,
		PosY:           pos_y,
		Left:           left,
		NewTargetIndex: new_target_index,
		TurnEndTime:    turn_end_time,
	}
	return NewS2cSceneTouShiBuildingTurnToNoticeProtoMsg(msg)
}

var s2c_scene_tou_shi_building_turn_to_notice = [...]byte{44, 122} // 122
func NewS2cSceneTouShiBuildingTurnToNoticeProtoMsg(object *S2CSceneTouShiBuildingTurnToNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_tou_shi_building_turn_to_notice[:], "s2c_scene_tou_shi_building_turn_to_notice")

}

func NewS2cSceneTouShiBuildingFireMsg(pos_x int32, pos_y int32, target_index int32, prepare_end_time int32, bomb_explode_time int32) pbutil.Buffer {
	msg := &S2CSceneTouShiBuildingFireProto{
		PosX:            pos_x,
		PosY:            pos_y,
		TargetIndex:     target_index,
		PrepareEndTime:  prepare_end_time,
		BombExplodeTime: bomb_explode_time,
	}
	return NewS2cSceneTouShiBuildingFireProtoMsg(msg)
}

var s2c_scene_tou_shi_building_fire = [...]byte{44, 124} // 124
func NewS2cSceneTouShiBuildingFireProtoMsg(object *S2CSceneTouShiBuildingFireProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_tou_shi_building_fire[:], "s2c_scene_tou_shi_building_fire")

}

// 当前不是名城战战斗时间
var ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_INVALID_TIME = pbutil.StaticBuffer{3, 44, 125, 6} // 125-6

// 没有参战
var ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_NOT_IN_SCENE = pbutil.StaticBuffer{3, 44, 125, 7} // 125-7

// 这座城的名城战已经结束了
var ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_MC_WAR_END = pbutil.StaticBuffer{3, 44, 125, 8} // 125-8

// 坐标错误或不是投石机
var ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_INVALID_POS = pbutil.StaticBuffer{3, 44, 125, 1} // 125-1

// 不是占领者
var ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_NOT_HOST = pbutil.StaticBuffer{3, 44, 125, 2} // 125-2

// 没有目标
var ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_NO_TARGET = pbutil.StaticBuffer{3, 44, 125, 3} // 125-3

// 正在转向中
var ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_IN_TURN_CD = pbutil.StaticBuffer{3, 44, 125, 4} // 125-4

// 正在装填中
var ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_IN_PREPARE_CD = pbutil.StaticBuffer{3, 44, 125, 5} // 125-5

// 服务器错误
var ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 44, 125, 9} // 125-9

func NewS2cSceneTouShiBuildingFireNoticeMsg(pos_x int32, pos_y int32, target_index int32, prepare_end_time int32, bomb_explode_time int32, fire_hero_name string, fire_hero_country int32) pbutil.Buffer {
	msg := &S2CSceneTouShiBuildingFireNoticeProto{
		PosX:            pos_x,
		PosY:            pos_y,
		TargetIndex:     target_index,
		PrepareEndTime:  prepare_end_time,
		BombExplodeTime: bomb_explode_time,
		FireHeroName:    fire_hero_name,
		FireHeroCountry: fire_hero_country,
	}
	return NewS2cSceneTouShiBuildingFireNoticeProtoMsg(msg)
}

var s2c_scene_tou_shi_building_fire_notice = [...]byte{44, 126} // 126
func NewS2cSceneTouShiBuildingFireNoticeProtoMsg(object *S2CSceneTouShiBuildingFireNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_tou_shi_building_fire_notice[:], "s2c_scene_tou_shi_building_fire_notice")

}

func NewS2cSceneTouShiBombExplodeNoticeMsg(fire_troop_id []byte) pbutil.Buffer {
	msg := &S2CSceneTouShiBombExplodeNoticeProto{
		FireTroopId: fire_troop_id,
	}
	return NewS2cSceneTouShiBombExplodeNoticeProtoMsg(msg)
}

var s2c_scene_tou_shi_bomb_explode_notice = [...]byte{44, 127} // 127
func NewS2cSceneTouShiBombExplodeNoticeProtoMsg(object *S2CSceneTouShiBombExplodeNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_tou_shi_bomb_explode_notice[:], "s2c_scene_tou_shi_bomb_explode_notice")

}

func NewS2cSceneDrumMsg(drum_desc string, next_drum_time int32, add_drum_stat *shared_proto.SpriteStatProto) pbutil.Buffer {
	msg := &S2CSceneDrumProto{
		DrumDesc:     drum_desc,
		NextDrumTime: next_drum_time,
		AddDrumStat:  add_drum_stat,
	}
	return NewS2cSceneDrumProtoMsg(msg)
}

var s2c_scene_drum = [...]byte{44, 129, 1} // 129
func NewS2cSceneDrumProtoMsg(object *S2CSceneDrumProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_drum[:], "s2c_scene_drum")

}

// 没有参战
var ERR_SCENE_DRUM_FAIL_NOT_IN_SCENE = pbutil.StaticBuffer{4, 44, 130, 1, 1} // 130-1

// 击鼓时间已经结束
var ERR_SCENE_DRUM_FAIL_INVALID_TIME = pbutil.StaticBuffer{4, 44, 130, 1, 2} // 130-2

// 还在击鼓 CD 中
var ERR_SCENE_DRUM_FAIL_IN_CD = pbutil.StaticBuffer{4, 44, 130, 1, 3} // 130-3

// 百战千军等级不够
var ERR_SCENE_DRUM_FAIL_BAI_ZHAN_LEVEL_LIMIT = pbutil.StaticBuffer{4, 44, 130, 1, 5} // 130-5

// 服务器错误
var ERR_SCENE_DRUM_FAIL_SERVER_ERR = pbutil.StaticBuffer{4, 44, 130, 1, 4} // 130-4

func NewS2cSceneDrumNoticeMsg(hero_id []byte, is_atk bool, new_drum_times int32, add_drum_stat *shared_proto.SpriteStatProto, new_drum_stat *shared_proto.SpriteStatProto) pbutil.Buffer {
	msg := &S2CSceneDrumNoticeProto{
		HeroId:       hero_id,
		IsAtk:        is_atk,
		NewDrumTimes: new_drum_times,
		AddDrumStat:  add_drum_stat,
		NewDrumStat:  new_drum_stat,
	}
	return NewS2cSceneDrumNoticeProtoMsg(msg)
}

var s2c_scene_drum_notice = [...]byte{44, 131, 1} // 131
func NewS2cSceneDrumNoticeProtoMsg(object *S2CSceneDrumNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_drum_notice[:], "s2c_scene_drum_notice")

}

func NewS2cSceneDrumAddStatNoticeMsg(troops [][]byte) pbutil.Buffer {
	msg := &S2CSceneDrumAddStatNoticeProto{
		Troops: troops,
	}
	return NewS2cSceneDrumAddStatNoticeProtoMsg(msg)
}

var s2c_scene_drum_add_stat_notice = [...]byte{44, 132, 1} // 132
func NewS2cSceneDrumAddStatNoticeProtoMsg(object *S2CSceneDrumAddStatNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_scene_drum_add_stat_notice[:], "s2c_scene_drum_add_stat_notice")

}

func NewS2cMingcHostUpdateNoticeMsg(mc_id int32, new_host *shared_proto.GuildBasicProto) pbutil.Buffer {
	msg := &S2CMingcHostUpdateNoticeProto{
		McId:    mc_id,
		NewHost: new_host,
	}
	return NewS2cMingcHostUpdateNoticeProtoMsg(msg)
}

var s2c_mingc_host_update_notice = [...]byte{44, 134, 1} // 134
func NewS2cMingcHostUpdateNoticeProtoMsg(object *S2CMingcHostUpdateNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_mingc_host_update_notice[:], "s2c_mingc_host_update_notice")

}
