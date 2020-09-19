package xiongnu

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
	MODULE_ID = 36

	C2S_SET_DEFENDER = 1

	C2S_START = 5

	C2S_TROOP_INFO = 10

	C2S_GET_XIONG_NU_NPC_BASE_INFO = 14

	C2S_GET_DEFENSER_FIGHT_AMOUNT = 17

	C2S_GET_XIONG_NU_FIGHT_INFO = 19
)

func NewS2cSetDefenderMsg(id []byte, to_set bool) pbutil.Buffer {
	msg := &S2CSetDefenderProto{
		Id:    id,
		ToSet: to_set,
	}
	return NewS2cSetDefenderProtoMsg(msg)
}

var s2c_set_defender = [...]byte{36, 2} // 2
func NewS2cSetDefenderProtoMsg(object *S2CSetDefenderProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_defender[:], "s2c_set_defender")

}

func NewS2cBroacastSetDefenderMsg(id []byte, to_set bool) pbutil.Buffer {
	msg := &S2CBroacastSetDefenderProto{
		Id:    id,
		ToSet: to_set,
	}
	return NewS2cBroacastSetDefenderProtoMsg(msg)
}

var s2c_broacast_set_defender = [...]byte{36, 3} // 3
func NewS2cBroacastSetDefenderProtoMsg(object *S2CBroacastSetDefenderProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_broacast_set_defender[:], "s2c_broacast_set_defender")

}

// 没有联盟
var ERR_SET_DEFENDER_FAIL_NO_GUILD = pbutil.StaticBuffer{3, 36, 4, 1} // 4-1

// 已经开启了
var ERR_SET_DEFENDER_FAIL_STARTED = pbutil.StaticBuffer{3, 36, 4, 2} // 4-2

// 目标没有找到
var ERR_SET_DEFENDER_FAIL_TARGET_NOT_FOUND = pbutil.StaticBuffer{3, 36, 4, 3} // 4-3

// 没有开启抗击匈奴的权限
var ERR_SET_DEFENDER_FAIL_NO_PERMISISON = pbutil.StaticBuffer{3, 36, 4, 4} // 4-4

// 目标不是联盟成员
var ERR_SET_DEFENDER_FAIL_TARGET_NOT_MEMBER = pbutil.StaticBuffer{3, 36, 4, 5} // 4-5

// 防守成员已满
var ERR_SET_DEFENDER_FAIL_FULL = pbutil.StaticBuffer{3, 36, 4, 6} // 4-6

// 今天已经参加了
var ERR_SET_DEFENDER_FAIL_TODAY_START = pbutil.StaticBuffer{3, 36, 4, 7} // 4-7

// 目标主城流亡了，无法设置防守
var ERR_SET_DEFENDER_FAIL_HOME_NOT_ALIVE = pbutil.StaticBuffer{3, 36, 4, 8} // 4-8

// 服务器忙，请稍后再试
var ERR_SET_DEFENDER_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 36, 4, 9} // 4-9

func NewS2cStartMsg(base_id int32, base_x int32, base_y int32) pbutil.Buffer {
	msg := &S2CStartProto{
		BaseId: base_id,
		BaseX:  base_x,
		BaseY:  base_y,
	}
	return NewS2cStartProtoMsg(msg)
}

var s2c_start = [...]byte{36, 6} // 6
func NewS2cStartProtoMsg(object *S2CStartProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_start[:], "s2c_start")

}

func NewS2cBroadcastStartMsg(name string) pbutil.Buffer {
	msg := &S2CBroadcastStartProto{
		Name: name,
	}
	return NewS2cBroadcastStartProtoMsg(msg)
}

var s2c_broadcast_start = [...]byte{36, 8} // 8
func NewS2cBroadcastStartProtoMsg(object *S2CBroadcastStartProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_broadcast_start[:], "s2c_broadcast_start")

}

// 没有联盟
var ERR_START_FAIL_NO_GUILD = pbutil.StaticBuffer{3, 36, 7, 1} // 7-1

// 已经开启了
var ERR_START_FAIL_STARTED = pbutil.StaticBuffer{3, 36, 7, 2} // 7-2

// 非法的等级数据
var ERR_START_FAIL_INVALID_LEVEL = pbutil.StaticBuffer{3, 36, 7, 3} // 7-3

// 你主城当前已经流亡
var ERR_START_FAIL_BASE_DEAD = pbutil.StaticBuffer{3, 36, 7, 9} // 7-9

// 没有开启抗击匈奴的权限
var ERR_START_FAIL_NO_PERMISISON = pbutil.StaticBuffer{3, 36, 7, 4} // 7-4

// 防守成员不足
var ERR_START_FAIL_NOT_ENOUGH = pbutil.StaticBuffer{3, 36, 7, 5} // 7-5

// 挑战难度未解锁
var ERR_START_FAIL_LOCK_LEVEL = pbutil.StaticBuffer{3, 36, 7, 6} // 7-6

// 联盟等级不足
var ERR_START_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 36, 7, 7} // 7-7

// 服务器忙，请稍后再试
var ERR_START_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 36, 7, 8} // 7-8

// 没到开启时间
var ERR_START_FAIL_START_TIME_LIMIT = pbutil.StaticBuffer{3, 36, 7, 10} // 7-10

func NewS2cInfoBroadcastMsg(info []byte) pbutil.Buffer {
	msg := &S2CInfoBroadcastProto{
		Info: info,
	}
	return NewS2cInfoBroadcastProtoMsg(msg)
}

var s2c_info_broadcast = [...]byte{36, 9} // 9
func NewS2cInfoBroadcastProtoMsg(object *S2CInfoBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_info_broadcast[:], "s2c_info_broadcast")

}

func NewS2cTroopInfoMsg(baseTroops []byte, morale int32) pbutil.Buffer {
	msg := &S2CTroopInfoProto{
		BaseTroops: baseTroops,
		Morale:     morale,
	}
	return NewS2cTroopInfoProtoMsg(msg)
}

var s2c_troop_info = [...]byte{36, 11} // 11
func NewS2cTroopInfoProtoMsg(object *S2CTroopInfoProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_troop_info[:], "s2c_troop_info")

}

// 没有联盟
var ERR_TROOP_INFO_FAIL_NO_GUILD = pbutil.StaticBuffer{3, 36, 12, 1} // 12-1

// 没有开启
var ERR_TROOP_INFO_FAIL_NOT_STARTED = pbutil.StaticBuffer{3, 36, 12, 2} // 12-2

// 服务器忙，请稍后再试
var ERR_TROOP_INFO_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 36, 12, 3} // 12-3

func NewS2cEndBroadcastMsg(guild_id int32, resist_xiong_nu []byte, unlock_next_level bool, add_prestige int32) pbutil.Buffer {
	msg := &S2CEndBroadcastProto{
		GuildId:         guild_id,
		ResistXiongNu:   resist_xiong_nu,
		UnlockNextLevel: unlock_next_level,
		AddPrestige:     add_prestige,
	}
	return NewS2cEndBroadcastProtoMsg(msg)
}

var s2c_end_broadcast = [...]byte{36, 13} // 13
func NewS2cEndBroadcastProtoMsg(object *S2CEndBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_end_broadcast[:], "s2c_end_broadcast")

}

func NewS2cGetXiongNuNpcBaseInfoMsg(guild_id int32, guild_name string, guild_flag string, morale int32, start_time int32, fighting_amount int32) pbutil.Buffer {
	msg := &S2CGetXiongNuNpcBaseInfoProto{
		GuildId:        guild_id,
		GuildName:      guild_name,
		GuildFlag:      guild_flag,
		Morale:         morale,
		StartTime:      start_time,
		FightingAmount: fighting_amount,
	}
	return NewS2cGetXiongNuNpcBaseInfoProtoMsg(msg)
}

var s2c_get_xiong_nu_npc_base_info = [...]byte{36, 15} // 15
func NewS2cGetXiongNuNpcBaseInfoProtoMsg(object *S2CGetXiongNuNpcBaseInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_xiong_nu_npc_base_info[:], "s2c_get_xiong_nu_npc_base_info")

}

// 活动没开启
var ERR_GET_XIONG_NU_NPC_BASE_INFO_FAIL_NOT_STARTED = pbutil.StaticBuffer{3, 36, 16, 1} // 16-1

func NewS2cGetDefenserFightAmountMsg(version int32, defenser_id [][]byte, defenser_fight_amount []int32, defenser_enemy_count []int32) pbutil.Buffer {
	msg := &S2CGetDefenserFightAmountProto{
		Version:             version,
		DefenserId:          defenser_id,
		DefenserFightAmount: defenser_fight_amount,
		DefenserEnemyCount:  defenser_enemy_count,
	}
	return NewS2cGetDefenserFightAmountProtoMsg(msg)
}

var s2c_get_defenser_fight_amount = [...]byte{36, 18} // 18
func NewS2cGetDefenserFightAmountProtoMsg(object *S2CGetDefenserFightAmountProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_defenser_fight_amount[:], "s2c_get_defenser_fight_amount")

}

func NewS2cGetXiongNuFightInfoMsg(data []byte) pbutil.Buffer {
	msg := &S2CGetXiongNuFightInfoProto{
		Data: data,
	}
	return NewS2cGetXiongNuFightInfoProtoMsg(msg)
}

func NewS2cGetXiongNuFightInfoMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CGetXiongNuFightInfoProto{
		Data: safeMarshal(data),
	}
	return NewS2cGetXiongNuFightInfoProtoMsg(msg)
}

var s2c_get_xiong_nu_fight_info = [...]byte{36, 20} // 20
func NewS2cGetXiongNuFightInfoProtoMsg(object *S2CGetXiongNuFightInfoProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_get_xiong_nu_fight_info[:], "s2c_get_xiong_nu_fight_info")

}

// 你没有联盟，不能获取匈奴战斗排行榜
var ERR_GET_XIONG_NU_FIGHT_INFO_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 36, 21, 1} // 21-1

// 服务器忙，请稍后再试
var ERR_GET_XIONG_NU_FIGHT_INFO_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 36, 21, 2} // 21-2
