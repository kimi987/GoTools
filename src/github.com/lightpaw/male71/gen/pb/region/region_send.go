package region

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
	MODULE_ID = 7

	C2S_UPDATE_SELF_VIEW = 148

	C2S_CLOSE_VIEW = 150

	C2S_PRE_INVASION_TARGET = 94

	C2S_WATCH_BASE_UNIT = 161

	C2S_REQUEST_TROOP_UNIT = 156

	C2S_REQUEST_RUINS_BASE = 131

	C2S_USE_MIAN_GOODS = 91

	C2S_UPGRADE_BASE = 40

	C2S_WHITE_FLAG_DETAIL = 86

	C2S_GET_BUY_PROSPERITY_COST = 211

	C2S_BUY_PROSPERITY = 106

	C2S_SWITCH_ACTION = 46

	C2S_REQUEST_MILITARY_PUSH = 165

	C2S_CREATE_BASE = 1

	C2S_FAST_MOVE_BASE = 14

	C2S_INVASION = 24

	C2S_CANCEL_INVASION = 26

	C2S_REPATRIATE = 71

	C2S_BAOZ_REPATRIATE = 186

	C2S_SPEED_UP = 139

	C2S_EXPEL = 30

	C2S_FAVORITE_POS = 99

	C2S_FAVORITE_POS_LIST = 102

	C2S_GET_PREV_INVESTIGATE = 175

	C2S_INVESTIGATE = 142

	C2S_INVESTIGATE_INVADE = 234

	C2S_USE_MULTI_LEVEL_NPC_TIMES_GOODS = 183

	C2S_USE_INVASE_HERO_TIMES_GOODS = 190

	C2S_CALC_MOVE_SPEED = 172

	C2S_LIST_ENEMY_POS = 178

	C2S_SEARCH_BAOZ_NPC = 180

	C2S_HOME_AST_DEFENDING_INFO = 193

	C2S_GUILD_PLEASE_HELP_ME = 196

	C2S_CREATE_ASSEMBLY = 199

	C2S_SHOW_ASSEMBLY = 202

	C2S_JOIN_ASSEMBLY = 206

	C2S_CREATE_GUILD_WORKSHOP = 214

	C2S_SHOW_GUILD_WORKSHOP = 217

	C2S_HURT_GUILD_WORKSHOP = 219

	C2S_REMOVE_GUILD_WORKSHOP = 223

	C2S_CATCH_GUILD_WORKSHOP_LOGS = 228

	C2S_GET_SELF_BAOZ = 232
)

func NewS2cUpdateMapRadiusMsg(center_x int32, center_y int32, radius int32) pbutil.Buffer {
	msg := &S2CUpdateMapRadiusProto{
		CenterX: center_x,
		CenterY: center_y,
		Radius:  radius,
	}
	return NewS2cUpdateMapRadiusProtoMsg(msg)
}

var s2c_update_map_radius = [...]byte{7, 159, 1} // 159
func NewS2cUpdateMapRadiusProtoMsg(object *S2CUpdateMapRadiusProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_map_radius[:], "s2c_update_map_radius")

}

func NewS2cUpdateSelfViewMsg(min_x int32, min_y int32, max_x int32, max_y int32) pbutil.Buffer {
	msg := &S2CUpdateSelfViewProto{
		MinX: min_x,
		MinY: min_y,
		MaxX: max_x,
		MaxY: max_y,
	}
	return NewS2cUpdateSelfViewProtoMsg(msg)
}

var s2c_update_self_view = [...]byte{7, 149, 1} // 149
func NewS2cUpdateSelfViewProtoMsg(object *S2CUpdateSelfViewProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_self_view[:], "s2c_update_self_view")

}

var CLOSE_VIEW_S2C = pbutil.StaticBuffer{3, 7, 151, 1} // 151

func NewS2cAddBaseUnitMsg(add_type int32, data []byte) pbutil.Buffer {
	msg := &S2CAddBaseUnitProto{
		AddType: add_type,
		Data:    data,
	}
	return NewS2cAddBaseUnitProtoMsg(msg)
}

func NewS2cAddBaseUnitMarshalMsg(add_type int32, data marshaler) pbutil.Buffer {
	msg := &S2CAddBaseUnitProto{
		AddType: add_type,
		Data:    safeMarshal(data),
	}
	return NewS2cAddBaseUnitProtoMsg(msg)
}

var s2c_add_base_unit = [...]byte{7, 152, 1} // 152
func NewS2cAddBaseUnitProtoMsg(object *S2CAddBaseUnitProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_base_unit[:], "s2c_add_base_unit")

}

func NewS2cUpdateNpcBaseInfoMsg(map_id int32, npc_id []byte, base_x int32, base_y int32, data_id int32, data_type int32, guild_id int32, guild_name string, guild_flag_name string, country int32, mian_disappear_time int32, prosperity int32, has_defenser bool, hero []byte, hero_end_time int32, hero_type int32, progress int32, total_progress int32, progress_type int32) pbutil.Buffer {
	msg := &S2CUpdateNpcBaseInfoProto{
		MapId:             map_id,
		NpcId:             npc_id,
		BaseX:             base_x,
		BaseY:             base_y,
		DataId:            data_id,
		DataType:          data_type,
		GuildId:           guild_id,
		GuildName:         guild_name,
		GuildFlagName:     guild_flag_name,
		Country:           country,
		MianDisappearTime: mian_disappear_time,
		Prosperity:        prosperity,
		HasDefenser:       has_defenser,
		Hero:              hero,
		HeroEndTime:       hero_end_time,
		HeroType:          hero_type,
		Progress:          progress,
		TotalProgress:     total_progress,
		ProgressType:      progress_type,
	}
	return NewS2cUpdateNpcBaseInfoProtoMsg(msg)
}

func NewS2cUpdateNpcBaseInfoMarshalMsg(map_id int32, npc_id []byte, base_x int32, base_y int32, data_id int32, data_type int32, guild_id int32, guild_name string, guild_flag_name string, country int32, mian_disappear_time int32, prosperity int32, has_defenser bool, hero marshaler, hero_end_time int32, hero_type int32, progress int32, total_progress int32, progress_type int32) pbutil.Buffer {
	msg := &S2CUpdateNpcBaseInfoProto{
		MapId:             map_id,
		NpcId:             npc_id,
		BaseX:             base_x,
		BaseY:             base_y,
		DataId:            data_id,
		DataType:          data_type,
		GuildId:           guild_id,
		GuildName:         guild_name,
		GuildFlagName:     guild_flag_name,
		Country:           country,
		MianDisappearTime: mian_disappear_time,
		Prosperity:        prosperity,
		HasDefenser:       has_defenser,
		Hero:              safeMarshal(hero),
		HeroEndTime:       hero_end_time,
		HeroType:          hero_type,
		Progress:          progress,
		TotalProgress:     total_progress,
		ProgressType:      progress_type,
	}
	return NewS2cUpdateNpcBaseInfoProtoMsg(msg)
}

var s2c_update_npc_base_info = [...]byte{7, 110} // 110
func NewS2cUpdateNpcBaseInfoProtoMsg(object *S2CUpdateNpcBaseInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_npc_base_info[:], "s2c_update_npc_base_info")

}

func NewS2cUpdateBaseProgressMsg(id []byte, progress int32, total_progress int32, progress_type int32) pbutil.Buffer {
	msg := &S2CUpdateBaseProgressProto{
		Id:            id,
		Progress:      progress,
		TotalProgress: total_progress,
		ProgressType:  progress_type,
	}
	return NewS2cUpdateBaseProgressProtoMsg(msg)
}

var s2c_update_base_progress = [...]byte{7, 213, 1} // 213
func NewS2cUpdateBaseProgressProtoMsg(object *S2CUpdateBaseProgressProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_base_progress[:], "s2c_update_base_progress")

}

func NewS2cRemoveBaseUnitMsg(remove_type int32, hero_id []byte) pbutil.Buffer {
	msg := &S2CRemoveBaseUnitProto{
		RemoveType: remove_type,
		HeroId:     hero_id,
	}
	return NewS2cRemoveBaseUnitProtoMsg(msg)
}

var s2c_remove_base_unit = [...]byte{7, 153, 1} // 153
func NewS2cRemoveBaseUnitProtoMsg(object *S2CRemoveBaseUnitProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_base_unit[:], "s2c_remove_base_unit")

}

func NewS2cPreInvasionTargetMsg(head string, level int32, max_tower_floor int32, jun_xian_level int32) pbutil.Buffer {
	msg := &S2CPreInvasionTargetProto{
		Head:          head,
		Level:         level,
		MaxTowerFloor: max_tower_floor,
		JunXianLevel:  jun_xian_level,
	}
	return NewS2cPreInvasionTargetProtoMsg(msg)
}

var s2c_pre_invasion_target = [...]byte{7, 95} // 95
func NewS2cPreInvasionTargetProtoMsg(object *S2CPreInvasionTargetProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_pre_invasion_target[:], "s2c_pre_invasion_target")

}

// 无效的目标id
var ERR_PRE_INVASION_TARGET_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 7, 96, 1} // 96-1

func NewS2cWatchBaseUnitMsg(target []byte, guild_name string, fight_amount int32, prosprity int32, prosprity_capcity int32, head string, level int32, max_tower_floor int32, jun_xian_level int32, soldier int32, captain_soldier []int32, hero []byte, hero_end_time int32, hero_type int32) pbutil.Buffer {
	msg := &S2CWatchBaseUnitProto{
		Target:           target,
		GuildName:        guild_name,
		FightAmount:      fight_amount,
		Prosprity:        prosprity,
		ProsprityCapcity: prosprity_capcity,
		Head:             head,
		Level:            level,
		MaxTowerFloor:    max_tower_floor,
		JunXianLevel:     jun_xian_level,
		Soldier:          soldier,
		CaptainSoldier:   captain_soldier,
		Hero:             hero,
		HeroEndTime:      hero_end_time,
		HeroType:         hero_type,
	}
	return NewS2cWatchBaseUnitProtoMsg(msg)
}

func NewS2cWatchBaseUnitMarshalMsg(target []byte, guild_name string, fight_amount int32, prosprity int32, prosprity_capcity int32, head string, level int32, max_tower_floor int32, jun_xian_level int32, soldier int32, captain_soldier []int32, hero marshaler, hero_end_time int32, hero_type int32) pbutil.Buffer {
	msg := &S2CWatchBaseUnitProto{
		Target:           target,
		GuildName:        guild_name,
		FightAmount:      fight_amount,
		Prosprity:        prosprity,
		ProsprityCapcity: prosprity_capcity,
		Head:             head,
		Level:            level,
		MaxTowerFloor:    max_tower_floor,
		JunXianLevel:     jun_xian_level,
		Soldier:          soldier,
		CaptainSoldier:   captain_soldier,
		Hero:             safeMarshal(hero),
		HeroEndTime:      hero_end_time,
		HeroType:         hero_type,
	}
	return NewS2cWatchBaseUnitProtoMsg(msg)
}

var s2c_watch_base_unit = [...]byte{7, 162, 1} // 162
func NewS2cWatchBaseUnitProtoMsg(object *S2CWatchBaseUnitProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_watch_base_unit[:], "s2c_watch_base_unit")

}

// 无效的目标id
var ERR_WATCH_BASE_UNIT_FAIL_INVALID_TARGET = pbutil.StaticBuffer{4, 7, 163, 1, 1} // 163-1

func NewS2cUpdateWatchBaseProsperityMsg(target []byte, prosprity int32, prosprity_capcity int32) pbutil.Buffer {
	msg := &S2CUpdateWatchBaseProsperityProto{
		Target:           target,
		Prosprity:        prosprity,
		ProsprityCapcity: prosprity_capcity,
	}
	return NewS2cUpdateWatchBaseProsperityProtoMsg(msg)
}

var s2c_update_watch_base_prosperity = [...]byte{7, 164, 1} // 164
func NewS2cUpdateWatchBaseProsperityProtoMsg(object *S2CUpdateWatchBaseProsperityProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_watch_base_prosperity[:], "s2c_update_watch_base_prosperity")

}

func NewS2cUpdateStopLostProsperityMsg(target []byte) pbutil.Buffer {
	msg := &S2CUpdateStopLostProsperityProto{
		Target: target,
	}
	return NewS2cUpdateStopLostProsperityProtoMsg(msg)
}

var s2c_update_stop_lost_prosperity = [...]byte{7, 174, 1} // 174
func NewS2cUpdateStopLostProsperityProtoMsg(object *S2CUpdateStopLostProsperityProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_stop_lost_prosperity[:], "s2c_update_stop_lost_prosperity")

}

func NewS2cAddTroopUnitMsg(add_type int32, data []byte) pbutil.Buffer {
	msg := &S2CAddTroopUnitProto{
		AddType: add_type,
		Data:    data,
	}
	return NewS2cAddTroopUnitProtoMsg(msg)
}

func NewS2cAddTroopUnitMarshalMsg(add_type int32, data marshaler) pbutil.Buffer {
	msg := &S2CAddTroopUnitProto{
		AddType: add_type,
		Data:    safeMarshal(data),
	}
	return NewS2cAddTroopUnitProtoMsg(msg)
}

var s2c_add_troop_unit = [...]byte{7, 154, 1} // 154
func NewS2cAddTroopUnitProtoMsg(object *S2CAddTroopUnitProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_troop_unit[:], "s2c_add_troop_unit")

}

func NewS2cRemoveTroopUnitMsg(remove_type int32, troop_id []byte) pbutil.Buffer {
	msg := &S2CRemoveTroopUnitProto{
		RemoveType: remove_type,
		TroopId:    troop_id,
	}
	return NewS2cRemoveTroopUnitProtoMsg(msg)
}

var s2c_remove_troop_unit = [...]byte{7, 155, 1} // 155
func NewS2cRemoveTroopUnitProtoMsg(object *S2CRemoveTroopUnitProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_troop_unit[:], "s2c_remove_troop_unit")

}

func NewS2cRequestTroopUnitMsg(data []byte) pbutil.Buffer {
	msg := &S2CRequestTroopUnitProto{
		Data: data,
	}
	return NewS2cRequestTroopUnitProtoMsg(msg)
}

func NewS2cRequestTroopUnitMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CRequestTroopUnitProto{
		Data: safeMarshal(data),
	}
	return NewS2cRequestTroopUnitProtoMsg(msg)
}

var s2c_request_troop_unit = [...]byte{7, 157, 1} // 157
func NewS2cRequestTroopUnitProtoMsg(object *S2CRequestTroopUnitProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_troop_unit[:], "s2c_request_troop_unit")

}

// 无效的部队id
var ERR_REQUEST_TROOP_UNIT_FAIL_INVALID_ID = pbutil.StaticBuffer{4, 7, 158, 1, 1} // 158-1

// 服务器忙，请稍后再试
var ERR_REQUEST_TROOP_UNIT_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 158, 1, 2} // 158-2

func NewS2cAddRuinsBaseMsg(pos_x int32, pos_y int32) pbutil.Buffer {
	msg := &S2CAddRuinsBaseProto{
		PosX: pos_x,
		PosY: pos_y,
	}
	return NewS2cAddRuinsBaseProtoMsg(msg)
}

var s2c_add_ruins_base = [...]byte{7, 134, 1} // 134
func NewS2cAddRuinsBaseProtoMsg(object *S2CAddRuinsBaseProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_ruins_base[:], "s2c_add_ruins_base")

}

func NewS2cRemoveRuinsBaseMsg(pos_x int32, pos_y int32) pbutil.Buffer {
	msg := &S2CRemoveRuinsBaseProto{
		PosX: pos_x,
		PosY: pos_y,
	}
	return NewS2cRemoveRuinsBaseProtoMsg(msg)
}

var s2c_remove_ruins_base = [...]byte{7, 130, 1} // 130
func NewS2cRemoveRuinsBaseProtoMsg(object *S2CRemoveRuinsBaseProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_ruins_base[:], "s2c_remove_ruins_base")

}

func NewS2cRequestRuinsBaseMsg(realm_id int32, pos_x int32, pos_y int32, hero_basic []byte) pbutil.Buffer {
	msg := &S2CRequestRuinsBaseProto{
		RealmId:   realm_id,
		PosX:      pos_x,
		PosY:      pos_y,
		HeroBasic: hero_basic,
	}
	return NewS2cRequestRuinsBaseProtoMsg(msg)
}

var s2c_request_ruins_base = [...]byte{7, 132, 1} // 132
func NewS2cRequestRuinsBaseProtoMsg(object *S2CRequestRuinsBaseProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_ruins_base[:], "s2c_request_ruins_base")

}

// 场景没找到
var ERR_REQUEST_RUINS_BASE_FAIL_REALM_NOT_FOUND = pbutil.StaticBuffer{4, 7, 133, 1, 1} // 133-1

// 错误的坐标
var ERR_REQUEST_RUINS_BASE_FAIL_INVALID_X_OR_Y = pbutil.StaticBuffer{4, 7, 133, 1, 2} // 133-2

// 废墟没找到
var ERR_REQUEST_RUINS_BASE_FAIL_NO_RUINS = pbutil.StaticBuffer{4, 7, 133, 1, 3} // 133-3

// 服务器忙，请稍后再试
var ERR_REQUEST_RUINS_BASE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 133, 1, 4} // 133-4

func NewS2cUpdateSelfMianDisappearTimeMsg(mian_disappear_time int32, mian_start_time int32) pbutil.Buffer {
	msg := &S2CUpdateSelfMianDisappearTimeProto{
		MianDisappearTime: mian_disappear_time,
		MianStartTime:     mian_start_time,
	}
	return NewS2cUpdateSelfMianDisappearTimeProtoMsg(msg)
}

var s2c_update_self_mian_disappear_time = [...]byte{7, 90} // 90
func NewS2cUpdateSelfMianDisappearTimeProtoMsg(object *S2CUpdateSelfMianDisappearTimeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_self_mian_disappear_time[:], "s2c_update_self_mian_disappear_time")

}

func NewS2cUseMianGoodsMsg(id int32, cooldown int32) pbutil.Buffer {
	msg := &S2CUseMianGoodsProto{
		Id:       id,
		Cooldown: cooldown,
	}
	return NewS2cUseMianGoodsProtoMsg(msg)
}

var s2c_use_mian_goods = [...]byte{7, 92} // 92
func NewS2cUseMianGoodsProtoMsg(object *S2CUseMianGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_mian_goods[:], "s2c_use_mian_goods")

}

// 无效的物品id
var ERR_USE_MIAN_GOODS_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 7, 93, 1} // 93-1

// 不是免战物品id
var ERR_USE_MIAN_GOODS_FAIL_INVALID_GOODS = pbutil.StaticBuffer{3, 7, 93, 2} // 93-2

// 物品个数不足
var ERR_USE_MIAN_GOODS_FAIL_COUNT_NOT_ENOUGH = pbutil.StaticBuffer{3, 7, 93, 3} // 93-3

// 部队出征中，请先召回部队
var ERR_USE_MIAN_GOODS_FAIL_TROOP_OUTSIDE = pbutil.StaticBuffer{3, 7, 93, 4} // 93-4

// 行营出征中，请先收回行营
var ERR_USE_MIAN_GOODS_FAIL_TENT_OUTSIDE = pbutil.StaticBuffer{3, 7, 93, 5} // 93-5

// 主城流亡了，请先重建主城
var ERR_USE_MIAN_GOODS_FAIL_HOME_NOT_ALIVE = pbutil.StaticBuffer{3, 7, 93, 6} // 93-6

// 免战中
var ERR_USE_MIAN_GOODS_FAIL_MIAN = pbutil.StaticBuffer{3, 7, 93, 7} // 93-7

// 免战物品冷却中
var ERR_USE_MIAN_GOODS_FAIL_COOLDOWN = pbutil.StaticBuffer{3, 7, 93, 9} // 93-9

// 服务器忙，请稍后再试
var ERR_USE_MIAN_GOODS_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 7, 93, 8} // 93-8

func NewS2cUpdateNewHeroMianDisappearTimeMsg(disappear_time int32) pbutil.Buffer {
	msg := &S2CUpdateNewHeroMianDisappearTimeProto{
		DisappearTime: disappear_time,
	}
	return NewS2cUpdateNewHeroMianDisappearTimeProtoMsg(msg)
}

var s2c_update_new_hero_mian_disappear_time = [...]byte{7, 160, 1} // 160
func NewS2cUpdateNewHeroMianDisappearTimeProtoMsg(object *S2CUpdateNewHeroMianDisappearTimeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_new_hero_mian_disappear_time[:], "s2c_update_new_hero_mian_disappear_time")

}

var UPGRADE_BASE_S2C = pbutil.StaticBuffer{2, 7, 41} // 41

// 流亡状态
var ERR_UPGRADE_BASE_FAIL_BASE_NOT_EXIST = pbutil.StaticBuffer{3, 7, 42, 1} // 42-1

// 已经达到最大等级
var ERR_UPGRADE_BASE_FAIL_BASE_MAX_LEVEL = pbutil.StaticBuffer{3, 7, 42, 3} // 42-3

// 繁荣度不足
var ERR_UPGRADE_BASE_FAIL_PROSPRITY_NOT_ENOUGH = pbutil.StaticBuffer{3, 7, 42, 2} // 42-2

// 服务器忙，请稍后再试
var ERR_UPGRADE_BASE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 7, 42, 4} // 42-4

func NewS2cSelfUpdateBaseLevelMsg(level int32) pbutil.Buffer {
	msg := &S2CSelfUpdateBaseLevelProto{
		Level: level,
	}
	return NewS2cSelfUpdateBaseLevelProtoMsg(msg)
}

var s2c_self_update_base_level = [...]byte{7, 52} // 52
func NewS2cSelfUpdateBaseLevelProtoMsg(object *S2CSelfUpdateBaseLevelProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_self_update_base_level[:], "s2c_self_update_base_level")

}

func NewS2cUpdateWhiteFlagMsg(hero_id []byte, white_flag_guild_id int32, white_flag_flag_name string, white_flag_disappear_time int32) pbutil.Buffer {
	msg := &S2CUpdateWhiteFlagProto{
		HeroId:                 hero_id,
		WhiteFlagGuildId:       white_flag_guild_id,
		WhiteFlagFlagName:      white_flag_flag_name,
		WhiteFlagDisappearTime: white_flag_disappear_time,
	}
	return NewS2cUpdateWhiteFlagProtoMsg(msg)
}

var s2c_update_white_flag = [...]byte{7, 85} // 85
func NewS2cUpdateWhiteFlagProtoMsg(object *S2CUpdateWhiteFlagProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_white_flag[:], "s2c_update_white_flag")

}

func NewS2cWhiteFlagDetailMsg(hero_id []byte, white_flag_hero_id []byte, white_flag_hero_name string, white_flag_guild_id int32, white_flag_guild_name string, white_flag_disappear_time int32, white_flag_country int32) pbutil.Buffer {
	msg := &S2CWhiteFlagDetailProto{
		HeroId:                 hero_id,
		WhiteFlagHeroId:        white_flag_hero_id,
		WhiteFlagHeroName:      white_flag_hero_name,
		WhiteFlagGuildId:       white_flag_guild_id,
		WhiteFlagGuildName:     white_flag_guild_name,
		WhiteFlagDisappearTime: white_flag_disappear_time,
		WhiteFlagCountry:       white_flag_country,
	}
	return NewS2cWhiteFlagDetailProtoMsg(msg)
}

var s2c_white_flag_detail = [...]byte{7, 87} // 87
func NewS2cWhiteFlagDetailProtoMsg(object *S2CWhiteFlagDetailProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_white_flag_detail[:], "s2c_white_flag_detail")

}

// 英雄当前没有插白旗
var ERR_WHITE_FLAG_DETAIL_FAIL_NO_FLAG = pbutil.StaticBuffer{3, 7, 88, 1} // 88-1

func NewS2cSelfBaseDestroyMsg(is_tent bool, destroy_type int32) pbutil.Buffer {
	msg := &S2CSelfBaseDestroyProto{
		IsTent:      is_tent,
		DestroyType: destroy_type,
	}
	return NewS2cSelfBaseDestroyProtoMsg(msg)
}

var s2c_self_base_destroy = [...]byte{7, 135, 1} // 135
func NewS2cSelfBaseDestroyProtoMsg(object *S2CSelfBaseDestroyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_self_base_destroy[:], "s2c_self_base_destroy")

}

func NewS2cProsperityBufMsg(end_time int32) pbutil.Buffer {
	msg := &S2CProsperityBufProto{
		EndTime: end_time,
	}
	return NewS2cProsperityBufProtoMsg(msg)
}

var s2c_prosperity_buf = [...]byte{7, 105} // 105
func NewS2cProsperityBufProtoMsg(object *S2CProsperityBufProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_prosperity_buf[:], "s2c_prosperity_buf")

}

func NewS2cShowWordsMsg(base_id []byte, troop_id []byte, prosperity int32, gold int32, food int32, wood int32, stone int32, jade int32, jade_ore int32) pbutil.Buffer {
	msg := &S2CShowWordsProto{
		BaseId:     base_id,
		TroopId:    troop_id,
		Prosperity: prosperity,
		Gold:       gold,
		Food:       food,
		Wood:       wood,
		Stone:      stone,
		Jade:       jade,
		JadeOre:    jade_ore,
	}
	return NewS2cShowWordsProtoMsg(msg)
}

var s2c_show_words = [...]byte{7, 122} // 122
func NewS2cShowWordsProtoMsg(object *S2CShowWordsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_show_words[:], "s2c_show_words")

}

func NewS2cGetBuyProsperityCostMsg(cost int32) pbutil.Buffer {
	msg := &S2CGetBuyProsperityCostProto{
		Cost: cost,
	}
	return NewS2cGetBuyProsperityCostProtoMsg(msg)
}

var s2c_get_buy_prosperity_cost = [...]byte{7, 212, 1} // 212
func NewS2cGetBuyProsperityCostProtoMsg(object *S2CGetBuyProsperityCostProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_buy_prosperity_cost[:], "s2c_get_buy_prosperity_cost")

}

func NewS2cBuyProsperityMsg(add_prosperity int32) pbutil.Buffer {
	msg := &S2CBuyProsperityProto{
		AddProsperity: add_prosperity,
	}
	return NewS2cBuyProsperityProtoMsg(msg)
}

var s2c_buy_prosperity = [...]byte{7, 107} // 107
func NewS2cBuyProsperityProtoMsg(object *S2CBuyProsperityProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_buy_prosperity[:], "s2c_buy_prosperity")

}

// 主城流亡了，请先重建主城
var ERR_BUY_PROSPERITY_FAIL_HOME_NOT_ALIVE = pbutil.StaticBuffer{3, 7, 108, 1} // 108-1

// 繁荣度已满
var ERR_BUY_PROSPERITY_FAIL_PROSPERITY_FULL = pbutil.StaticBuffer{3, 7, 108, 2} // 108-2

// 点券不足
var ERR_BUY_PROSPERITY_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 7, 108, 5} // 108-5

// 服务器忙，请稍后再试
var ERR_BUY_PROSPERITY_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 7, 108, 4} // 108-4

// vip等级不够
var ERR_BUY_PROSPERITY_FAIL_VIP_LEVEL_LIMIT = pbutil.StaticBuffer{3, 7, 108, 6} // 108-6

func NewS2cSelfBeenAttackRobChangedMsg(been_attack int32, been_rob int32) pbutil.Buffer {
	msg := &S2CSelfBeenAttackRobChangedProto{
		BeenAttack: been_attack,
		BeenRob:    been_rob,
	}
	return NewS2cSelfBeenAttackRobChangedProtoMsg(msg)
}

var s2c_self_been_attack_rob_changed = [...]byte{7, 128, 1} // 128
func NewS2cSelfBeenAttackRobChangedProtoMsg(object *S2CSelfBeenAttackRobChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_self_been_attack_rob_changed[:], "s2c_self_been_attack_rob_changed")

}

func NewS2cGuildBeenAttackRobChangedMsg(total int32) pbutil.Buffer {
	msg := &S2CGuildBeenAttackRobChangedProto{
		Total: total,
	}
	return NewS2cGuildBeenAttackRobChangedProtoMsg(msg)
}

var s2c_guild_been_attack_rob_changed = [...]byte{7, 129, 1} // 129
func NewS2cGuildBeenAttackRobChangedProtoMsg(object *S2CGuildBeenAttackRobChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_guild_been_attack_rob_changed[:], "s2c_guild_been_attack_rob_changed")

}

func NewS2cSwitchActionMsg(open bool) pbutil.Buffer {
	msg := &S2CSwitchActionProto{
		Open: open,
	}
	return NewS2cSwitchActionProtoMsg(msg)
}

var s2c_switch_action = [...]byte{7, 47} // 47
func NewS2cSwitchActionProtoMsg(object *S2CSwitchActionProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_switch_action[:], "s2c_switch_action")

}

func NewS2cRequestMilitaryPushMsg(main_military bool, guild_military bool, to_target []byte, to_target_base bool, from_target []byte) pbutil.Buffer {
	msg := &S2CRequestMilitaryPushProto{
		MainMilitary:  main_military,
		GuildMilitary: guild_military,
		ToTarget:      to_target,
		ToTargetBase:  to_target_base,
		FromTarget:    from_target,
	}
	return NewS2cRequestMilitaryPushProtoMsg(msg)
}

var s2c_request_military_push = [...]byte{7, 166, 1} // 166
func NewS2cRequestMilitaryPushProtoMsg(object *S2CRequestMilitaryPushProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_military_push[:], "s2c_request_military_push")

}

// 无效的玩家id
var ERR_REQUEST_MILITARY_PUSH_FAIL_INVALID_ID = pbutil.StaticBuffer{4, 7, 167, 1, 1} // 167-1

// 你没有联盟，不能请求联盟军情
var ERR_REQUEST_MILITARY_PUSH_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 7, 167, 1, 2} // 167-2

func NewS2cUpdateMilitaryInfoMsg(data []byte, region bool, ma bool) pbutil.Buffer {
	msg := &S2CUpdateMilitaryInfoProto{
		Data:   data,
		Region: region,
		Ma:     ma,
	}
	return NewS2cUpdateMilitaryInfoProtoMsg(msg)
}

var s2c_update_military_info = [...]byte{7, 22} // 22
func NewS2cUpdateMilitaryInfoProtoMsg(object *S2CUpdateMilitaryInfoProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_update_military_info[:], "s2c_update_military_info")

}

func NewS2cRemoveMilitaryInfoMsg(id []byte, region bool, ma bool) pbutil.Buffer {
	msg := &S2CRemoveMilitaryInfoProto{
		Id:     id,
		Region: region,
		Ma:     ma,
	}
	return NewS2cRemoveMilitaryInfoProtoMsg(msg)
}

var s2c_remove_military_info = [...]byte{7, 23} // 23
func NewS2cRemoveMilitaryInfoProtoMsg(object *S2CRemoveMilitaryInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_military_info[:], "s2c_remove_military_info")

}

func NewS2cUpdateSelfMilitaryInfoMsg(troop_index int32, troop_id []byte, data []byte) pbutil.Buffer {
	msg := &S2CUpdateSelfMilitaryInfoProto{
		TroopIndex: troop_index,
		TroopId:    troop_id,
		Data:       data,
	}
	return NewS2cUpdateSelfMilitaryInfoProtoMsg(msg)
}

var s2c_update_self_military_info = [...]byte{7, 50} // 50
func NewS2cUpdateSelfMilitaryInfoProtoMsg(object *S2CUpdateSelfMilitaryInfoProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_update_self_military_info[:], "s2c_update_self_military_info")

}

func NewS2cRemoveSelfMilitaryInfoMsg(id []byte) pbutil.Buffer {
	msg := &S2CRemoveSelfMilitaryInfoProto{
		Id: id,
	}
	return NewS2cRemoveSelfMilitaryInfoProtoMsg(msg)
}

var s2c_remove_self_military_info = [...]byte{7, 51} // 51
func NewS2cRemoveSelfMilitaryInfoProtoMsg(object *S2CRemoveSelfMilitaryInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_self_military_info[:], "s2c_remove_self_military_info")

}

func NewS2cNpcBaseInfoMsg(map_id int32, npc_id [][]byte, base_x []int32, base_y []int32, data_id []int32) pbutil.Buffer {
	msg := &S2CNpcBaseInfoProto{
		MapId:  map_id,
		NpcId:  npc_id,
		BaseX:  base_x,
		BaseY:  base_y,
		DataId: data_id,
	}
	return NewS2cNpcBaseInfoProtoMsg(msg)
}

var s2c_npc_base_info = [...]byte{7, 109} // 109
func NewS2cNpcBaseInfoProtoMsg(object *S2CNpcBaseInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_npc_base_info[:], "s2c_npc_base_info")

}

func NewS2cCreateBaseMsg(map_id int32, new_x int32, new_y int32, level int32, prosperity int32) pbutil.Buffer {
	msg := &S2CCreateBaseProto{
		MapId:      map_id,
		NewX:       new_x,
		NewY:       new_y,
		Level:      level,
		Prosperity: prosperity,
	}
	return NewS2cCreateBaseProtoMsg(msg)
}

var s2c_create_base = [...]byte{7, 2} // 2
func NewS2cCreateBaseProtoMsg(object *S2CCreateBaseProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_create_base[:], "s2c_create_base")

}

// 发送的mapid无效
var ERR_CREATE_BASE_FAIL_INVALID_MAP_ID = pbutil.StaticBuffer{3, 7, 3, 1} // 3-1

// 发送的坐标无效（城市不能建在边缘，需要满足周围6格都有位置）
var ERR_CREATE_BASE_FAIL_INVALID_POS = pbutil.StaticBuffer{3, 7, 3, 2} // 3-2

// 不是流亡状态
var ERR_CREATE_BASE_FAIL_BASE_EXIST = pbutil.StaticBuffer{3, 7, 3, 6} // 3-6

// 距离其他玩家太近
var ERR_CREATE_BASE_FAIL_TOO_CLOSE_OTHER = pbutil.StaticBuffer{3, 7, 3, 3} // 3-3

// 地图已满，不能再新建主城
var ERR_CREATE_BASE_FAIL_FULL = pbutil.StaticBuffer{3, 7, 3, 4} // 3-4

// 服务器忙，请稍后再试
var ERR_CREATE_BASE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 7, 3, 5} // 3-5

func NewS2cFastMoveBaseMsg(map_id int32, new_x int32, new_y int32, is_tent bool) pbutil.Buffer {
	msg := &S2CFastMoveBaseProto{
		MapId:  map_id,
		NewX:   new_x,
		NewY:   new_y,
		IsTent: is_tent,
	}
	return NewS2cFastMoveBaseProtoMsg(msg)
}

var s2c_fast_move_base = [...]byte{7, 15} // 15
func NewS2cFastMoveBaseProtoMsg(object *S2CFastMoveBaseProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fast_move_base[:], "s2c_fast_move_base")

}

// 无效的地图id
var ERR_FAST_MOVE_BASE_FAIL_INVALID_MAP_ID = pbutil.StaticBuffer{3, 7, 16, 6} // 16-6

// 发送的坐标无效（城市不能建在边缘，需要满足周围6格都有位置）
var ERR_FAST_MOVE_BASE_FAIL_INVALID_POS = pbutil.StaticBuffer{3, 7, 16, 1} // 16-1

// 迁移的坐标跟现在的坐标一样
var ERR_FAST_MOVE_BASE_FAIL_SELF_POS = pbutil.StaticBuffer{3, 7, 16, 22} // 16-22

// 流亡状态
var ERR_FAST_MOVE_BASE_FAIL_BASE_NOT_EXIST = pbutil.StaticBuffer{3, 7, 16, 2} // 16-2

// 距离其他玩家太近
var ERR_FAST_MOVE_BASE_FAIL_TOO_CLOSE_OTHER = pbutil.StaticBuffer{3, 7, 16, 3} // 16-3

// 地图已满，没有空位迁移
var ERR_FAST_MOVE_BASE_FAIL_FULL = pbutil.StaticBuffer{3, 7, 16, 4} // 16-4

// 无效的物品id
var ERR_FAST_MOVE_BASE_FAIL_INVALID_GOODS = pbutil.StaticBuffer{3, 7, 16, 7} // 16-7

// 物品个数不足
var ERR_FAST_MOVE_BASE_FAIL_GOODS_NOT_ENOUGH = pbutil.StaticBuffer{3, 7, 16, 8} // 16-8

// 当前有武将出征，不能迁城
var ERR_FAST_MOVE_BASE_FAIL_CAPTAIN_OUT_SIDE = pbutil.StaticBuffer{3, 7, 16, 9} // 16-9

// 当前行营在外面，不能迁移主城
var ERR_FAST_MOVE_BASE_FAIL_TENT_OUT_SIDE = pbutil.StaticBuffer{3, 7, 16, 10} // 16-10

// 迁移行营，但是行营不在野外
var ERR_FAST_MOVE_BASE_FAIL_TENT_NOT_EXIST = pbutil.StaticBuffer{3, 7, 16, 11} // 16-11

// 迁移行营，主城也在这张地图中
var ERR_FAST_MOVE_BASE_FAIL_MAP_HAS_HOME = pbutil.StaticBuffer{3, 7, 16, 12} // 16-12

// 迁移主城，行营也在这张地图中
var ERR_FAST_MOVE_BASE_FAIL_MAP_HAS_TENT = pbutil.StaticBuffer{3, 7, 16, 13} // 16-13

// 这个地图不允许主城进入
var ERR_FAST_MOVE_BASE_FAIL_MAP_DENY_HOME = pbutil.StaticBuffer{3, 7, 16, 14} // 16-14

// 这个地图只允许所属联盟主城进入
var ERR_FAST_MOVE_BASE_FAIL_MAP_GUILD_MEMBER_ONLY = pbutil.StaticBuffer{3, 7, 16, 15} // 16-15

// 免费迁移CD中
var ERR_FAST_MOVE_BASE_FAIL_COOLDOWN = pbutil.StaticBuffer{3, 7, 16, 16} // 16-16

// 免费迁移不能使用（只能行营使用，只能迁移到荣誉地区）
var ERR_FAST_MOVE_BASE_FAIL_FREE_CANT_USE = pbutil.StaticBuffer{3, 7, 16, 17} // 16-17

// 荣誉地区CD中（被打出来CD时间内不能再进入）
var ERR_FAST_MOVE_BASE_FAIL_MONSTER_COOLDOWN = pbutil.StaticBuffer{3, 7, 16, 18} // 16-18

// 主城地区CD中（被打出来CD时间内不能再进入）
var ERR_FAST_MOVE_BASE_FAIL_HOME_COOLDOWN = pbutil.StaticBuffer{3, 7, 16, 19} // 16-19

// 这个等级的荣誉地区未解锁
var ERR_FAST_MOVE_BASE_FAIL_MONSTER_LEVEL_LOCKED = pbutil.StaticBuffer{3, 7, 16, 20} // 16-20

// 联盟地区内圈行营个数超出上限，请放在外圈
var ERR_FAST_MOVE_BASE_FAIL_HOME_AREA_TENT_COUNT_LIMIT = pbutil.StaticBuffer{3, 7, 16, 21} // 16-21

// 使用的是联盟迁城令，但没加入联盟
var ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_NOT_IN_GUILD = pbutil.StaticBuffer{3, 7, 16, 23} // 16-23

// 联盟迁城令，地图满了
var ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_FULL = pbutil.StaticBuffer{3, 7, 16, 24} // 16-24

// 联盟迁城令，自己是盟主
var ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_IS_LEADER = pbutil.StaticBuffer{3, 7, 16, 25} // 16-25

// 联盟迁城令，盟主流亡
var ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_LEADER_NO_BASE = pbutil.StaticBuffer{3, 7, 16, 26} // 16-26

// 联盟迁城令，已经在盟主周围
var ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_ALREADY_IN_LEADER_AROUND = pbutil.StaticBuffer{3, 7, 16, 27} // 16-27

// 联盟迁城令，盟主是 npc
var ERR_FAST_MOVE_BASE_FAIL_GUILD_MOVE_BASE_LEADER_IS_NPC = pbutil.StaticBuffer{3, 7, 16, 28} // 16-28

// 服务器忙，请稍后再试
var ERR_FAST_MOVE_BASE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 7, 16, 5} // 16-5

func NewS2cInvasionMsg(target []byte, troop_index int32) pbutil.Buffer {
	msg := &S2CInvasionProto{
		Target:     target,
		TroopIndex: troop_index,
	}
	return NewS2cInvasionProtoMsg(msg)
}

var s2c_invasion = [...]byte{7, 25} // 25
func NewS2cInvasionProtoMsg(object *S2CInvasionProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_invasion[:], "s2c_invasion")

}

// 无效的目标id
var ERR_INVASION_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 7, 34, 1} // 34-1

// 无效的目标出征类型，
var ERR_INVASION_FAIL_INVALID_TARGET_INVATION = pbutil.StaticBuffer{3, 7, 34, 13} // 34-13

// 无效的队伍序号
var ERR_INVASION_FAIL_INVALID_TROOP_INDEX = pbutil.StaticBuffer{3, 7, 34, 18} // 34-18

// 队伍出征中
var ERR_INVASION_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 7, 34, 5} // 34-5

// 武将士兵数为0
var ERR_INVASION_FAIL_NO_SOLDIER = pbutil.StaticBuffer{3, 7, 34, 8} // 34-8

// 目标处于流亡状态
var ERR_INVASION_FAIL_TARGET_NOT_EXIST = pbutil.StaticBuffer{3, 7, 34, 10} // 34-10

// 自己处于流亡状态
var ERR_INVASION_FAIL_SELF_NOT_EXIST = pbutil.StaticBuffer{3, 7, 34, 12} // 34-12

// 不在同一个地图
var ERR_INVASION_FAIL_NOT_SAME_MAP = pbutil.StaticBuffer{3, 7, 34, 11} // 34-11

// 出征部队已达上限
var ERR_INVASION_FAIL_MAX_INVATION_TROOPS = pbutil.StaticBuffer{3, 7, 34, 14} // 34-14

// 服务器忙，请稍后再试
var ERR_INVASION_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 7, 34, 6} // 34-6

// 那张地图上没有你的主城或可用的行营
var ERR_INVASION_FAIL_NO_BASE_IN_MAP = pbutil.StaticBuffer{3, 7, 34, 15} // 34-15

// 行营不能攻击行营
var ERR_INVASION_FAIL_NO_TENT_TO_TENT = pbutil.StaticBuffer{3, 7, 34, 16} // 34-16

// 出发的行营还未建好. 等valid time
var ERR_INVASION_FAIL_TENT_NOT_VALID = pbutil.StaticBuffer{3, 7, 34, 17} // 34-17

// 目标主城免战中，不能出征
var ERR_INVASION_FAIL_MIAN = pbutil.StaticBuffer{3, 7, 34, 19} // 34-19

// 选择的目标等级太高，无法出征
var ERR_INVASION_FAIL_TARGET_LEVEL_LOCKED = pbutil.StaticBuffer{3, 7, 34, 23} // 34-23

// 选择的目标已经有部队进行出征了
var ERR_INVASION_FAIL_DUPLICATE_TARGET = pbutil.StaticBuffer{3, 7, 34, 24} // 34-24

// 讨伐野怪功能还未开启
var ERR_INVASION_FAIL_MLN_FUNC_LOCKED = pbutil.StaticBuffer{3, 7, 34, 25} // 34-25

// 讨伐野怪次数不足
var ERR_INVASION_FAIL_MLN_TIMES_LIMIT = pbutil.StaticBuffer{3, 7, 34, 26} // 34-26

// 今日已经参与过反击匈奴了
var ERR_INVASION_FAIL_TODAY_JOIN_XIONG_NU = pbutil.StaticBuffer{3, 7, 34, 27} // 34-27

// 主城等级不足，不能出征
var ERR_INVASION_FAIL_REQUIRED_BASE_LEVEL = pbutil.StaticBuffer{3, 7, 34, 28} // 34-28

// 君主等级不足，不能出征
var ERR_INVASION_FAIL_REQUIRED_HERO_LEVEL = pbutil.StaticBuffer{3, 7, 34, 30} // 34-30

// 其他队伍正在打宝藏，不能出征
var ERR_INVASION_FAIL_BAOZ_TROOP_LIMIT = pbutil.StaticBuffer{3, 7, 34, 29} // 34-29

// 无效的物品id
var ERR_INVASION_FAIL_INVALID_GOODS = pbutil.StaticBuffer{3, 7, 34, 31} // 34-31

// 消耗不足
var ERR_INVASION_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 7, 34, 32} // 34-32

// 在名城战战斗中
var ERR_INVASION_FAIL_IN_MC_WAR_FIGHT = pbutil.StaticBuffer{3, 7, 34, 33} // 34-33

// 野怪讨伐多次，vip等级不够
var ERR_INVASION_FAIL_MULTI_LEVEL_MONSTER_COUNT_VIP_LIMIT = pbutil.StaticBuffer{3, 7, 34, 35} // 34-35

func NewS2cUpdateSelfTroopsMsg(id []int32, soldier []int32, fight_amount []int32, wounded_soldier int32, remove_outside bool, troop_index int32) pbutil.Buffer {
	msg := &S2CUpdateSelfTroopsProto{
		Id:             id,
		Soldier:        soldier,
		FightAmount:    fight_amount,
		WoundedSoldier: wounded_soldier,
		RemoveOutside:  remove_outside,
		TroopIndex:     troop_index,
	}
	return NewS2cUpdateSelfTroopsProtoMsg(msg)
}

var s2c_update_self_troops = [...]byte{7, 56} // 56
func NewS2cUpdateSelfTroopsProtoMsg(object *S2CUpdateSelfTroopsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_self_troops[:], "s2c_update_self_troops")

}

func NewS2cUpdateSelfTroopsOutsideMsg(troop_index int32, outside bool) pbutil.Buffer {
	msg := &S2CUpdateSelfTroopsOutsideProto{
		TroopIndex: troop_index,
		Outside:    outside,
	}
	return NewS2cUpdateSelfTroopsOutsideProtoMsg(msg)
}

var s2c_update_self_troops_outside = [...]byte{7, 210, 1} // 210
func NewS2cUpdateSelfTroopsOutsideProtoMsg(object *S2CUpdateSelfTroopsOutsideProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_self_troops_outside[:], "s2c_update_self_troops_outside")

}

func NewS2cCancelInvasionMsg(id []byte) pbutil.Buffer {
	msg := &S2CCancelInvasionProto{
		Id: id,
	}
	return NewS2cCancelInvasionProtoMsg(msg)
}

var s2c_cancel_invasion = [...]byte{7, 28} // 28
func NewS2cCancelInvasionProtoMsg(object *S2CCancelInvasionProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_cancel_invasion[:], "s2c_cancel_invasion")

}

// 无效的军情id
var ERR_CANCEL_INVASION_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 7, 29, 5} // 29-5

// 这条军情的部队不是你的，不能操作
var ERR_CANCEL_INVASION_FAIL_NOT_SELF = pbutil.StaticBuffer{3, 7, 29, 1} // 29-1

// 不是驻扎在你城里的同盟部队，不能叫回家
var ERR_CANCEL_INVASION_FAIL_NO_ARRIVED = pbutil.StaticBuffer{3, 7, 29, 3} // 29-3

// 正在召回中，不要反复操作
var ERR_CANCEL_INVASION_FAIL_BACKING = pbutil.StaticBuffer{3, 7, 29, 2} // 29-2

// 集结已出发，不能召回
var ERR_CANCEL_INVASION_FAIL_ASSEMBLY_STARTED = pbutil.StaticBuffer{3, 7, 29, 6} // 29-6

// 服务器忙，请稍后再试
var ERR_CANCEL_INVASION_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 7, 29, 4} // 29-4

func NewS2cRepatriateMsg(id []byte, is_tent bool) pbutil.Buffer {
	msg := &S2CRepatriateProto{
		Id:     id,
		IsTent: is_tent,
	}
	return NewS2cRepatriateProtoMsg(msg)
}

var s2c_repatriate = [...]byte{7, 72} // 72
func NewS2cRepatriateProtoMsg(object *S2CRepatriateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_repatriate[:], "s2c_repatriate")

}

// 军情id没找到
var ERR_REPATRIATE_FAIL_ID_NOT_FOUND = pbutil.StaticBuffer{3, 7, 73, 1} // 73-1

// 不是驻扎在你城里的同盟部队，不能叫回家
var ERR_REPATRIATE_FAIL_NO_DEFENDING = pbutil.StaticBuffer{3, 7, 73, 2} // 73-2

// 遣返盟友集结部队，你不是集结创建者
var ERR_REPATRIATE_FAIL_NO_ASSEMBLY_CREATER = pbutil.StaticBuffer{3, 7, 73, 4} // 73-4

// 集结已出发，不能遣返
var ERR_REPATRIATE_FAIL_ASSEMBLY_STARTED = pbutil.StaticBuffer{3, 7, 73, 5} // 73-5

// 服务器忙，请稍后再试
var ERR_REPATRIATE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 7, 73, 3} // 73-3

func NewS2cBaozRepatriateMsg(base_id []byte, troop_id []byte) pbutil.Buffer {
	msg := &S2CBaozRepatriateProto{
		BaseId:  base_id,
		TroopId: troop_id,
	}
	return NewS2cBaozRepatriateProtoMsg(msg)
}

var s2c_baoz_repatriate = [...]byte{7, 187, 1} // 187
func NewS2cBaozRepatriateProtoMsg(object *S2CBaozRepatriateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_baoz_repatriate[:], "s2c_baoz_repatriate")

}

// 部队id没找到
var ERR_BAOZ_REPATRIATE_FAIL_ID_NOT_FOUND = pbutil.StaticBuffer{4, 7, 188, 1, 1} // 188-1

// 不是你控制的宝藏，不能遣返
var ERR_BAOZ_REPATRIATE_FAIL_NOT_KEEP = pbutil.StaticBuffer{4, 7, 188, 1, 2} // 188-2

// 服务器忙，请稍后再试
var ERR_BAOZ_REPATRIATE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 188, 1, 3} // 188-3

func NewS2cSpeedUpMsg(id []byte) pbutil.Buffer {
	msg := &S2CSpeedUpProto{
		Id: id,
	}
	return NewS2cSpeedUpProtoMsg(msg)
}

var s2c_speed_up = [...]byte{7, 140, 1} // 140
func NewS2cSpeedUpProtoMsg(object *S2CSpeedUpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_speed_up[:], "s2c_speed_up")

}

// 军情id没找到
var ERR_SPEED_UP_FAIL_ID_NOT_FOUND = pbutil.StaticBuffer{4, 7, 141, 1, 1} // 141-1

// 发送的不是行军加速道具
var ERR_SPEED_UP_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 7, 141, 1, 2} // 141-2

// 物品个数不足
var ERR_SPEED_UP_FAIL_GOODS_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 141, 1, 5} // 141-5

// 不支持点券购买
var ERR_SPEED_UP_FAIL_COST_NOT_SUPPORT = pbutil.StaticBuffer{4, 7, 141, 1, 8} // 141-8

// 点券购买，点券不足
var ERR_SPEED_UP_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 141, 1, 9} // 141-9

// 部队不是行军中，不能加速
var ERR_SPEED_UP_FAIL_NO_MOVING = pbutil.StaticBuffer{4, 7, 141, 1, 3} // 141-3

// 加入目标部队没找到
var ERR_SPEED_UP_FAIL_OTHER_ID_NOT_FOUND = pbutil.StaticBuffer{4, 7, 141, 1, 10} // 141-10

// 部队集结等待中，不能加速
var ERR_SPEED_UP_FAIL_ASSEMBLY_WAIT = pbutil.StaticBuffer{4, 7, 141, 1, 11} // 141-11

// 服务器忙，请稍后再试
var ERR_SPEED_UP_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 141, 1, 4} // 141-4

func NewS2cExpelMsg(id []byte, cooldown int32, link string) pbutil.Buffer {
	msg := &S2CExpelProto{
		Id:       id,
		Cooldown: cooldown,
		Link:     link,
	}
	return NewS2cExpelProtoMsg(msg)
}

var s2c_expel = [...]byte{7, 31} // 31
func NewS2cExpelProtoMsg(object *S2CExpelProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_expel[:], "s2c_expel")

}

// 无效的军情id
var ERR_EXPEL_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 7, 32, 11} // 32-11

// 无效的地图id
var ERR_EXPEL_FAIL_INVALID_MAP = pbutil.StaticBuffer{3, 7, 32, 12} // 32-12

// 这条军情的部队不是正在掠夺你的，不能操作
var ERR_EXPEL_FAIL_NOT_SELF = pbutil.StaticBuffer{3, 7, 32, 1} // 32-1

// 这个军队当前不处于掠夺状态（还没到，或者回去了）
var ERR_EXPEL_FAIL_NOT_ARRIVED = pbutil.StaticBuffer{3, 7, 32, 2} // 32-2

// 驱逐CD中
var ERR_EXPEL_FAIL_COOLDOWN = pbutil.StaticBuffer{3, 7, 32, 3} // 32-3

// 无效的队伍编号
var ERR_EXPEL_FAIL_INVALID_TROOP_INDEX = pbutil.StaticBuffer{3, 7, 32, 13} // 32-13

// 队伍出征中，不能驱逐
var ERR_EXPEL_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 7, 32, 9} // 32-9

// 武将士兵数为0
var ERR_EXPEL_FAIL_NO_SOLDIER = pbutil.StaticBuffer{3, 7, 32, 10} // 32-10

// 服务器忙，请稍后再试
var ERR_EXPEL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 7, 32, 4} // 32-4

func NewS2cFavoritePosMsg(add bool, id int32, pos_x int32, pos_y int32) pbutil.Buffer {
	msg := &S2CFavoritePosProto{
		Add:  add,
		Id:   id,
		PosX: pos_x,
		PosY: pos_y,
	}
	return NewS2cFavoritePosProtoMsg(msg)
}

var s2c_favorite_pos = [...]byte{7, 100} // 100
func NewS2cFavoritePosProtoMsg(object *S2CFavoritePosProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_favorite_pos[:], "s2c_favorite_pos")

}

// 没有在收藏列表里面找到，无法删除
var ERR_FAVORITE_POS_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 7, 101, 1} // 101-1

// 场景没有找到
var ERR_FAVORITE_POS_FAIL_SCENE_NOT_FOUND = pbutil.StaticBuffer{3, 7, 101, 2} // 101-2

// 坐标非法
var ERR_FAVORITE_POS_FAIL_POS_INVALID = pbutil.StaticBuffer{3, 7, 101, 3} // 101-3

// 收藏点数量已满，无法添加
var ERR_FAVORITE_POS_FAIL_FULL = pbutil.StaticBuffer{3, 7, 101, 4} // 101-4

// 收藏点已经存在，无法添加
var ERR_FAVORITE_POS_FAIL_EXIST = pbutil.StaticBuffer{3, 7, 101, 5} // 101-5

// 服务器繁忙，请稍后再试
var ERR_FAVORITE_POS_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 7, 101, 6} // 101-6

func NewS2cFavoritePosListMsg(data []byte) pbutil.Buffer {
	msg := &S2CFavoritePosListProto{
		Data: data,
	}
	return NewS2cFavoritePosListProtoMsg(msg)
}

var s2c_favorite_pos_list = [...]byte{7, 103} // 103
func NewS2cFavoritePosListProtoMsg(object *S2CFavoritePosListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_favorite_pos_list[:], "s2c_favorite_pos_list")

}

// 服务器繁忙，请稍后再试
var ERR_FAVORITE_POS_LIST_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 7, 104, 1} // 104-1

func NewS2cGetPrevInvestigateMsg(hero_id []byte, mail_id []byte, expire_time int32) pbutil.Buffer {
	msg := &S2CGetPrevInvestigateProto{
		HeroId:     hero_id,
		MailId:     mail_id,
		ExpireTime: expire_time,
	}
	return NewS2cGetPrevInvestigateProtoMsg(msg)
}

var s2c_get_prev_investigate = [...]byte{7, 176, 1} // 176
func NewS2cGetPrevInvestigateProtoMsg(object *S2CGetPrevInvestigateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_prev_investigate[:], "s2c_get_prev_investigate")

}

func NewS2cInvestigateMsg(hero_id []byte, next_investigate_time int32) pbutil.Buffer {
	msg := &S2CInvestigateProto{
		HeroId:              hero_id,
		NextInvestigateTime: next_investigate_time,
	}
	return NewS2cInvestigateProtoMsg(msg)
}

var s2c_investigate = [...]byte{7, 143, 1} // 143
func NewS2cInvestigateProtoMsg(object *S2CInvestigateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_investigate[:], "s2c_investigate")

}

// 不能侦查自己
var ERR_INVESTIGATE_FAIL_SELF_ID = pbutil.StaticBuffer{4, 7, 144, 1, 5} // 144-5

// 英雄id没找到
var ERR_INVESTIGATE_FAIL_HERO_NOT_FOUND = pbutil.StaticBuffer{4, 7, 144, 1, 1} // 144-1

// 流亡状态不能侦查
var ERR_INVESTIGATE_FAIL_BASE_DESTROY = pbutil.StaticBuffer{4, 7, 144, 1, 3} // 144-3

// 侦查CD未到
var ERR_INVESTIGATE_FAIL_COOLDOWN = pbutil.StaticBuffer{4, 7, 144, 1, 2} // 144-2

// 盟友城池不能侦查
var ERR_INVESTIGATE_FAIL_SAME_GUILD = pbutil.StaticBuffer{4, 7, 144, 1, 6} // 144-6

// 清除侦查CD所需消耗不足
var ERR_INVESTIGATE_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 144, 1, 7} // 144-7

// 距离太远，不能侦查
var ERR_INVESTIGATE_FAIL_DISTANCE = pbutil.StaticBuffer{4, 7, 144, 1, 8} // 144-8

// 自己免战中，不能侦查
var ERR_INVESTIGATE_FAIL_SELF_MIAN = pbutil.StaticBuffer{4, 7, 144, 1, 9} // 144-9

// 目标免战中，不能侦查
var ERR_INVESTIGATE_FAIL_TARGET_MIAN = pbutil.StaticBuffer{4, 7, 144, 1, 10} // 144-10

// 名城战期间不能侦查
var ERR_INVESTIGATE_FAIL_IN_MC_WAR_FIGHT = pbutil.StaticBuffer{4, 7, 144, 1, 11} // 144-11

// 服务器忙，请稍后再试
var ERR_INVESTIGATE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 144, 1, 4} // 144-4

func NewS2cInvestigateInvadeMsg(target []byte) pbutil.Buffer {
	msg := &S2CInvestigateInvadeProto{
		Target: target,
	}
	return NewS2cInvestigateInvadeProtoMsg(msg)
}

var s2c_investigate_invade = [...]byte{7, 235, 1} // 235
func NewS2cInvestigateInvadeProtoMsg(object *S2CInvestigateInvadeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_investigate_invade[:], "s2c_investigate_invade")

}

// 不能侦查自己
var ERR_INVESTIGATE_INVADE_FAIL_SELF_ID = pbutil.StaticBuffer{4, 7, 236, 1, 1} // 236-1

// 英雄id没找到
var ERR_INVESTIGATE_INVADE_FAIL_HERO_NOT_FOUND = pbutil.StaticBuffer{4, 7, 236, 1, 2} // 236-2

// 地图没找到
var ERR_INVESTIGATE_INVADE_FAIL_MAP_NOT_FOUND = pbutil.StaticBuffer{4, 7, 236, 1, 3} // 236-3

// 流亡状态不能侦查
var ERR_INVESTIGATE_INVADE_FAIL_BASE_DESTROY = pbutil.StaticBuffer{4, 7, 236, 1, 4} // 236-4

// 盟友城池不能侦查
var ERR_INVESTIGATE_INVADE_FAIL_SAME_GUILD = pbutil.StaticBuffer{4, 7, 236, 1, 5} // 236-5

// 侦察消耗不足
var ERR_INVESTIGATE_INVADE_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 236, 1, 6} // 236-6

// 自己免战中，不能侦查
var ERR_INVESTIGATE_INVADE_FAIL_SELF_MIAN = pbutil.StaticBuffer{4, 7, 236, 1, 7} // 236-7

// 目标免战中，不能侦查
var ERR_INVESTIGATE_INVADE_FAIL_TARGET_MIAN = pbutil.StaticBuffer{4, 7, 236, 1, 8} // 236-8

// 名城战期间不能侦查
var ERR_INVESTIGATE_INVADE_FAIL_IN_MC_WAR_FIGHT = pbutil.StaticBuffer{4, 7, 236, 1, 9} // 236-9

// 服务器忙，请稍后再试
var ERR_INVESTIGATE_INVADE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 236, 1, 10} // 236-10

func NewS2cUpdateMultiLevelNpcPassLevelMsg(npc_type int32, level int32) pbutil.Buffer {
	msg := &S2CUpdateMultiLevelNpcPassLevelProto{
		NpcType: npc_type,
		Level:   level,
	}
	return NewS2cUpdateMultiLevelNpcPassLevelProtoMsg(msg)
}

var s2c_update_multi_level_npc_pass_level = [...]byte{7, 168, 1} // 168
func NewS2cUpdateMultiLevelNpcPassLevelProtoMsg(object *S2CUpdateMultiLevelNpcPassLevelProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_multi_level_npc_pass_level[:], "s2c_update_multi_level_npc_pass_level")

}

func NewS2cUpdateMultiLevelNpcHateMsg(npc_type int32, hate int32) pbutil.Buffer {
	msg := &S2CUpdateMultiLevelNpcHateProto{
		NpcType: npc_type,
		Hate:    hate,
	}
	return NewS2cUpdateMultiLevelNpcHateProtoMsg(msg)
}

var s2c_update_multi_level_npc_hate = [...]byte{7, 169, 1} // 169
func NewS2cUpdateMultiLevelNpcHateProtoMsg(object *S2CUpdateMultiLevelNpcHateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_multi_level_npc_hate[:], "s2c_update_multi_level_npc_hate")

}

func NewS2cUpdateMultiLevelNpcTimesMsg(start_recovey_time int32, times *shared_proto.RecoverableTimesWithExtraTimesProto) pbutil.Buffer {
	msg := &S2CUpdateMultiLevelNpcTimesProto{
		StartRecoveyTime: start_recovey_time,
		Times:            times,
	}
	return NewS2cUpdateMultiLevelNpcTimesProtoMsg(msg)
}

var s2c_update_multi_level_npc_times = [...]byte{7, 170, 1} // 170
func NewS2cUpdateMultiLevelNpcTimesProtoMsg(object *S2CUpdateMultiLevelNpcTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_multi_level_npc_times[:], "s2c_update_multi_level_npc_times")

}

func NewS2cUseMultiLevelNpcTimesGoodsMsg(id int32, buy bool) pbutil.Buffer {
	msg := &S2CUseMultiLevelNpcTimesGoodsProto{
		Id:  id,
		Buy: buy,
	}
	return NewS2cUseMultiLevelNpcTimesGoodsProtoMsg(msg)
}

var s2c_use_multi_level_npc_times_goods = [...]byte{7, 184, 1} // 184
func NewS2cUseMultiLevelNpcTimesGoodsProtoMsg(object *S2CUseMultiLevelNpcTimesGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_multi_level_npc_times_goods[:], "s2c_use_multi_level_npc_times_goods")

}

// 无效的物品id
var ERR_USE_MULTI_LEVEL_NPC_TIMES_GOODS_FAIL_INVALID_ID = pbutil.StaticBuffer{4, 7, 185, 1, 1} // 185-1

// 不是讨伐令物品id
var ERR_USE_MULTI_LEVEL_NPC_TIMES_GOODS_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 7, 185, 1, 2} // 185-2

// 物品个数不足
var ERR_USE_MULTI_LEVEL_NPC_TIMES_GOODS_FAIL_COUNT_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 185, 1, 3} // 185-3

// 购买消耗不足
var ERR_USE_MULTI_LEVEL_NPC_TIMES_GOODS_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 185, 1, 4} // 185-4

// 讨伐次数已满
var ERR_USE_MULTI_LEVEL_NPC_TIMES_GOODS_FAIL_FULL_TIMES = pbutil.StaticBuffer{4, 7, 185, 1, 6} // 185-6

func NewS2cUpdateInvaseHeroTimesMsg(start_recovey_time int32, times *shared_proto.RecoverableTimesWithExtraTimesProto) pbutil.Buffer {
	msg := &S2CUpdateInvaseHeroTimesProto{
		StartRecoveyTime: start_recovey_time,
		Times:            times,
	}
	return NewS2cUpdateInvaseHeroTimesProtoMsg(msg)
}

var s2c_update_invase_hero_times = [...]byte{7, 189, 1} // 189
func NewS2cUpdateInvaseHeroTimesProtoMsg(object *S2CUpdateInvaseHeroTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_invase_hero_times[:], "s2c_update_invase_hero_times")

}

func NewS2cUpdateJunTuanNpcTimesMsg(start_recovey_time int32, times *shared_proto.RecoverableTimesWithExtraTimesProto) pbutil.Buffer {
	msg := &S2CUpdateJunTuanNpcTimesProto{
		StartRecoveyTime: start_recovey_time,
		Times:            times,
	}
	return NewS2cUpdateJunTuanNpcTimesProtoMsg(msg)
}

var s2c_update_jun_tuan_npc_times = [...]byte{7, 209, 1} // 209
func NewS2cUpdateJunTuanNpcTimesProtoMsg(object *S2CUpdateJunTuanNpcTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_jun_tuan_npc_times[:], "s2c_update_jun_tuan_npc_times")

}

func NewS2cUseInvaseHeroTimesGoodsMsg(id int32, buy bool) pbutil.Buffer {
	msg := &S2CUseInvaseHeroTimesGoodsProto{
		Id:  id,
		Buy: buy,
	}
	return NewS2cUseInvaseHeroTimesGoodsProtoMsg(msg)
}

var s2c_use_invase_hero_times_goods = [...]byte{7, 191, 1} // 191
func NewS2cUseInvaseHeroTimesGoodsProtoMsg(object *S2CUseInvaseHeroTimesGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_invase_hero_times_goods[:], "s2c_use_invase_hero_times_goods")

}

// 无效的物品id
var ERR_USE_INVASE_HERO_TIMES_GOODS_FAIL_INVALID_ID = pbutil.StaticBuffer{4, 7, 192, 1, 1} // 192-1

// 无效的物品类型
var ERR_USE_INVASE_HERO_TIMES_GOODS_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 7, 192, 1, 2} // 192-2

// 物品个数不足
var ERR_USE_INVASE_HERO_TIMES_GOODS_FAIL_COUNT_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 192, 1, 3} // 192-3

// 购买消耗不足
var ERR_USE_INVASE_HERO_TIMES_GOODS_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 192, 1, 4} // 192-4

// 次数已满，无需购买
var ERR_USE_INVASE_HERO_TIMES_GOODS_FAIL_FULL_TIMES = pbutil.StaticBuffer{4, 7, 192, 1, 5} // 192-5

func NewS2cCalcMoveSpeedMsg(id []byte, speed int32) pbutil.Buffer {
	msg := &S2CCalcMoveSpeedProto{
		Id:    id,
		Speed: speed,
	}
	return NewS2cCalcMoveSpeedProtoMsg(msg)
}

var s2c_calc_move_speed = [...]byte{7, 173, 1} // 173
func NewS2cCalcMoveSpeedProtoMsg(object *S2CCalcMoveSpeedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_calc_move_speed[:], "s2c_calc_move_speed")

}

func NewS2cListEnemyPosMsg(pos_x []int32, pos_y []int32) pbutil.Buffer {
	msg := &S2CListEnemyPosProto{
		PosX: pos_x,
		PosY: pos_y,
	}
	return NewS2cListEnemyPosProtoMsg(msg)
}

var s2c_list_enemy_pos = [...]byte{7, 179, 1} // 179
func NewS2cListEnemyPosProtoMsg(object *S2CListEnemyPosProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_enemy_pos[:], "s2c_list_enemy_pos")

}

func NewS2cSearchBaozNpcMsg(data_id int32, base_id [][]byte, base_x []int32, base_y []int32) pbutil.Buffer {
	msg := &S2CSearchBaozNpcProto{
		DataId: data_id,
		BaseId: base_id,
		BaseX:  base_x,
		BaseY:  base_y,
	}
	return NewS2cSearchBaozNpcProtoMsg(msg)
}

var s2c_search_baoz_npc = [...]byte{7, 181, 1} // 181
func NewS2cSearchBaozNpcProtoMsg(object *S2CSearchBaozNpcProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_search_baoz_npc[:], "s2c_search_baoz_npc")

}

// 无效的配置id
var ERR_SEARCH_BAOZ_NPC_FAIL_INVALID_DATA_ID = pbutil.StaticBuffer{4, 7, 182, 1, 1} // 182-1

// 主城流亡了，请先重建主城
var ERR_SEARCH_BAOZ_NPC_FAIL_HOME_NOT_ALIVE = pbutil.StaticBuffer{4, 7, 182, 1, 2} // 182-2

func NewS2cHomeAstDefendingInfoMsg(heros []*shared_proto.HeroBasicProto, logs []*shared_proto.AstDefendLogProto) pbutil.Buffer {
	msg := &S2CHomeAstDefendingInfoProto{
		Heros: heros,
		Logs:  logs,
	}
	return NewS2cHomeAstDefendingInfoProtoMsg(msg)
}

var s2c_home_ast_defending_info = [...]byte{7, 194, 1} // 194
func NewS2cHomeAstDefendingInfoProtoMsg(object *S2CHomeAstDefendingInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_home_ast_defending_info[:], "s2c_home_ast_defending_info")

}

// 没有主城
var ERR_HOME_AST_DEFENDING_INFO_FAIL_NO_BASE = pbutil.StaticBuffer{4, 7, 195, 1, 1} // 195-1

// 繁荣度满
var ERR_HOME_AST_DEFENDING_INFO_FAIL_PROSPERITY_FULL = pbutil.StaticBuffer{4, 7, 195, 1, 2} // 195-2

var GUILD_PLEASE_HELP_ME_S2C = pbutil.StaticBuffer{3, 7, 197, 1} // 197

// 没有主城
var ERR_GUILD_PLEASE_HELP_ME_FAIL_NO_BASE = pbutil.StaticBuffer{4, 7, 198, 1, 1} // 198-1

// 繁荣度满
var ERR_GUILD_PLEASE_HELP_ME_FAIL_PROSPERITY_FULL = pbutil.StaticBuffer{4, 7, 198, 1, 2} // 198-2

// 援助部队满
var ERR_GUILD_PLEASE_HELP_ME_FAIL_TROOPS_LIMIT = pbutil.StaticBuffer{4, 7, 198, 1, 3} // 198-3

// 没有联盟
var ERR_GUILD_PLEASE_HELP_ME_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 7, 198, 1, 4} // 198-4

func NewS2cCreateAssemblyMsg(troop_index int32, target []byte, id []byte) pbutil.Buffer {
	msg := &S2CCreateAssemblyProto{
		TroopIndex: troop_index,
		Target:     target,
		Id:         id,
	}
	return NewS2cCreateAssemblyProtoMsg(msg)
}

var s2c_create_assembly = [...]byte{7, 200, 1} // 200
func NewS2cCreateAssemblyProtoMsg(object *S2CCreateAssemblyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_create_assembly[:], "s2c_create_assembly")

}

// 无效的目标id
var ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET = pbutil.StaticBuffer{4, 7, 201, 1, 1} // 201-1

// 无效的目标集结类型，
var ERR_CREATE_ASSEMBLY_FAIL_INVALID_TARGET_ASSEMBLY = pbutil.StaticBuffer{4, 7, 201, 1, 2} // 201-2

// 无效的集结等待时间
var ERR_CREATE_ASSEMBLY_FAIL_INVALID_WAIT_INDEX = pbutil.StaticBuffer{4, 7, 201, 1, 3} // 201-3

// 无效的队伍序号
var ERR_CREATE_ASSEMBLY_FAIL_INVALID_TROOP_INDEX = pbutil.StaticBuffer{4, 7, 201, 1, 4} // 201-4

// 队伍出征中
var ERR_CREATE_ASSEMBLY_FAIL_OUTSIDE = pbutil.StaticBuffer{4, 7, 201, 1, 5} // 201-5

// 武将士兵数为0
var ERR_CREATE_ASSEMBLY_FAIL_NO_SOLDIER = pbutil.StaticBuffer{4, 7, 201, 1, 6} // 201-6

// 目标处于流亡状态
var ERR_CREATE_ASSEMBLY_FAIL_TARGET_NOT_EXIST = pbutil.StaticBuffer{4, 7, 201, 1, 7} // 201-7

// 自己处于流亡状态
var ERR_CREATE_ASSEMBLY_FAIL_SELF_NOT_EXIST = pbutil.StaticBuffer{4, 7, 201, 1, 8} // 201-8

// 自己没有联盟
var ERR_CREATE_ASSEMBLY_FAIL_SELF_NOT_GUILD = pbutil.StaticBuffer{4, 7, 201, 1, 9} // 201-9

// 服务器忙，请稍后再试
var ERR_CREATE_ASSEMBLY_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 201, 1, 10} // 201-10

// 目标主城免战中，不能集结
var ERR_CREATE_ASSEMBLY_FAIL_MIAN = pbutil.StaticBuffer{4, 7, 201, 1, 11} // 201-11

// 集结次数不足
var ERR_CREATE_ASSEMBLY_FAIL_TIMES_LIMIT = pbutil.StaticBuffer{4, 7, 201, 1, 12} // 201-12

// 今日已经参与过反击匈奴了
var ERR_CREATE_ASSEMBLY_FAIL_TODAY_JOIN_XIONG_NU = pbutil.StaticBuffer{4, 7, 201, 1, 13} // 201-13

// 主城等级不足，不能集结
var ERR_CREATE_ASSEMBLY_FAIL_REQUIRED_BASE_LEVEL = pbutil.StaticBuffer{4, 7, 201, 1, 14} // 201-14

// 君主等级不足，不能集结
var ERR_CREATE_ASSEMBLY_FAIL_REQUIRED_HERO_LEVEL = pbutil.StaticBuffer{4, 7, 201, 1, 15} // 201-15

// 在名城战战斗中
var ERR_CREATE_ASSEMBLY_FAIL_IN_MC_WAR_FIGHT = pbutil.StaticBuffer{4, 7, 201, 1, 16} // 201-16

// 无效的物品id
var ERR_CREATE_ASSEMBLY_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 7, 201, 1, 17} // 201-17

// 消耗不足
var ERR_CREATE_ASSEMBLY_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 201, 1, 18} // 201-18

// 不能对同一个目标发起多次集结
var ERR_CREATE_ASSEMBLY_FAIL_SAME_TARGET = pbutil.StaticBuffer{4, 7, 201, 1, 19} // 201-19

func NewS2cShowAssemblyMsg(not_exist bool, id []byte, version int32, data []byte) pbutil.Buffer {
	msg := &S2CShowAssemblyProto{
		NotExist: not_exist,
		Id:       id,
		Version:  version,
		Data:     data,
	}
	return NewS2cShowAssemblyProtoMsg(msg)
}

func NewS2cShowAssemblyMarshalMsg(not_exist bool, id []byte, version int32, data marshaler) pbutil.Buffer {
	msg := &S2CShowAssemblyProto{
		NotExist: not_exist,
		Id:       id,
		Version:  version,
		Data:     safeMarshal(data),
	}
	return NewS2cShowAssemblyProtoMsg(msg)
}

var s2c_show_assembly = [...]byte{7, 203, 1} // 203
func NewS2cShowAssemblyProtoMsg(object *S2CShowAssemblyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_show_assembly[:], "s2c_show_assembly")

}

func NewS2cShowAssemblyChangedMsg(id []byte) pbutil.Buffer {
	msg := &S2CShowAssemblyChangedProto{
		Id: id,
	}
	return NewS2cShowAssemblyChangedProtoMsg(msg)
}

var s2c_show_assembly_changed = [...]byte{7, 205, 1} // 205
func NewS2cShowAssemblyChangedProtoMsg(object *S2CShowAssemblyChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_show_assembly_changed[:], "s2c_show_assembly_changed")

}

func NewS2cJoinAssemblyMsg(id []byte, troop_index int32) pbutil.Buffer {
	msg := &S2CJoinAssemblyProto{
		Id:         id,
		TroopIndex: troop_index,
	}
	return NewS2cJoinAssemblyProtoMsg(msg)
}

var s2c_join_assembly = [...]byte{7, 207, 1} // 207
func NewS2cJoinAssemblyProtoMsg(object *S2CJoinAssemblyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_join_assembly[:], "s2c_join_assembly")

}

// 无效的目标id
var ERR_JOIN_ASSEMBLY_FAIL_INVALID_TARGET = pbutil.StaticBuffer{4, 7, 208, 1, 1} // 208-1

// 无效的队伍序号
var ERR_JOIN_ASSEMBLY_FAIL_INVALID_TROOP_INDEX = pbutil.StaticBuffer{4, 7, 208, 1, 2} // 208-2

// 队伍出征中
var ERR_JOIN_ASSEMBLY_FAIL_OUTSIDE = pbutil.StaticBuffer{4, 7, 208, 1, 3} // 208-3

// 武将士兵数为0
var ERR_JOIN_ASSEMBLY_FAIL_NO_SOLDIER = pbutil.StaticBuffer{4, 7, 208, 1, 4} // 208-4

// 目标处于流亡状态
var ERR_JOIN_ASSEMBLY_FAIL_TARGET_NOT_EXIST = pbutil.StaticBuffer{4, 7, 208, 1, 5} // 208-5

// 自己处于流亡状态
var ERR_JOIN_ASSEMBLY_FAIL_SELF_NOT_EXIST = pbutil.StaticBuffer{4, 7, 208, 1, 6} // 208-6

// 自己没有联盟
var ERR_JOIN_ASSEMBLY_FAIL_SELF_NOT_GUILD = pbutil.StaticBuffer{4, 7, 208, 1, 7} // 208-7

// 服务器忙，请稍后再试
var ERR_JOIN_ASSEMBLY_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 208, 1, 8} // 208-8

// 目标主城免战中，不能集结
var ERR_JOIN_ASSEMBLY_FAIL_MIAN = pbutil.StaticBuffer{4, 7, 208, 1, 9} // 208-9

// 集结次数不足
var ERR_JOIN_ASSEMBLY_FAIL_TIMES_LIMIT = pbutil.StaticBuffer{4, 7, 208, 1, 10} // 208-10

// 今日已经参与过反击匈奴了
var ERR_JOIN_ASSEMBLY_FAIL_TODAY_JOIN_XIONG_NU = pbutil.StaticBuffer{4, 7, 208, 1, 11} // 208-11

// 主城等级不足，不能集结
var ERR_JOIN_ASSEMBLY_FAIL_REQUIRED_BASE_LEVEL = pbutil.StaticBuffer{4, 7, 208, 1, 12} // 208-12

// 君主等级不足，不能集结
var ERR_JOIN_ASSEMBLY_FAIL_REQUIRED_HERO_LEVEL = pbutil.StaticBuffer{4, 7, 208, 1, 13} // 208-13

// 在名城战战斗中
var ERR_JOIN_ASSEMBLY_FAIL_IN_MC_WAR_FIGHT = pbutil.StaticBuffer{4, 7, 208, 1, 14} // 208-14

// 集结已满
var ERR_JOIN_ASSEMBLY_FAIL_FULL = pbutil.StaticBuffer{4, 7, 208, 1, 15} // 208-15

// 不能多个队伍加入同一个集结
var ERR_JOIN_ASSEMBLY_FAIL_MULTI_JOIN = pbutil.StaticBuffer{4, 7, 208, 1, 16} // 208-16

// 集结已经出发
var ERR_JOIN_ASSEMBLY_FAIL_STARTED = pbutil.StaticBuffer{4, 7, 208, 1, 17} // 208-17

// 无效的物品id
var ERR_JOIN_ASSEMBLY_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 7, 208, 1, 18} // 208-18

// 消耗不足
var ERR_JOIN_ASSEMBLY_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 208, 1, 19} // 208-19

func NewS2cCreateGuildWorkshopMsg(pos_x int32, pos_y int32) pbutil.Buffer {
	msg := &S2CCreateGuildWorkshopProto{
		PosX: pos_x,
		PosY: pos_y,
	}
	return NewS2cCreateGuildWorkshopProtoMsg(msg)
}

var s2c_create_guild_workshop = [...]byte{7, 215, 1} // 215
func NewS2cCreateGuildWorkshopProtoMsg(object *S2CCreateGuildWorkshopProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_create_guild_workshop[:], "s2c_create_guild_workshop")

}

// 你不在联盟中
var ERR_CREATE_GUILD_WORKSHOP_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 7, 216, 1, 1} // 216-1

// 你没有权限操作
var ERR_CREATE_GUILD_WORKSHOP_FAIL_DENY = pbutil.StaticBuffer{4, 7, 216, 1, 2} // 216-2

// 无效的坐标位置
var ERR_CREATE_GUILD_WORKSHOP_FAIL_INVALID_POS = pbutil.StaticBuffer{4, 7, 216, 1, 3} // 216-3

// 已存在联盟工坊
var ERR_CREATE_GUILD_WORKSHOP_FAIL_BASE_EXIST = pbutil.StaticBuffer{4, 7, 216, 1, 4} // 216-4

// 超出最大距离限制
var ERR_CREATE_GUILD_WORKSHOP_FAIL_DISTANCE_LIMIT = pbutil.StaticBuffer{4, 7, 216, 1, 6} // 216-6

// 服务器忙，请稍后再试
var ERR_CREATE_GUILD_WORKSHOP_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 216, 1, 5} // 216-5

func NewS2cShowGuildWorkshopMsg(base_id []byte, guild_id int32, output int32, total_output int32, prize_count int32, been_hurt_times int32) pbutil.Buffer {
	msg := &S2CShowGuildWorkshopProto{
		BaseId:        base_id,
		GuildId:       guild_id,
		Output:        output,
		TotalOutput:   total_output,
		PrizeCount:    prize_count,
		BeenHurtTimes: been_hurt_times,
	}
	return NewS2cShowGuildWorkshopProtoMsg(msg)
}

var s2c_show_guild_workshop = [...]byte{7, 218, 1} // 218
func NewS2cShowGuildWorkshopProtoMsg(object *S2CShowGuildWorkshopProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_show_guild_workshop[:], "s2c_show_guild_workshop")

}

func NewS2cHurtGuildWorkshopMsg(base_id []byte, next_time int32, times int32) pbutil.Buffer {
	msg := &S2CHurtGuildWorkshopProto{
		BaseId:   base_id,
		NextTime: next_time,
		Times:    times,
	}
	return NewS2cHurtGuildWorkshopProtoMsg(msg)
}

var s2c_hurt_guild_workshop = [...]byte{7, 220, 1} // 220
func NewS2cHurtGuildWorkshopProtoMsg(object *S2CHurtGuildWorkshopProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_hurt_guild_workshop[:], "s2c_hurt_guild_workshop")

}

// 无效的id
var ERR_HURT_GUILD_WORKSHOP_FAIL_INVALID_BASE_ID = pbutil.StaticBuffer{4, 7, 221, 1, 1} // 221-1

// 破坏次数不足
var ERR_HURT_GUILD_WORKSHOP_FAIL_HURT_TIMES_NOT_ENOUGH = pbutil.StaticBuffer{4, 7, 221, 1, 2} // 221-2

// 目标被破坏次数已达上限
var ERR_HURT_GUILD_WORKSHOP_FAIL_BEEN_HURT_TIMES_LIMIT = pbutil.StaticBuffer{4, 7, 221, 1, 3} // 221-3

// 破坏CD中
var ERR_HURT_GUILD_WORKSHOP_FAIL_COOLDOWN = pbutil.StaticBuffer{4, 7, 221, 1, 5} // 221-5

// 服务器忙，请稍后再试
var ERR_HURT_GUILD_WORKSHOP_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 221, 1, 4} // 221-4

var REMOVE_GUILD_WORKSHOP_S2C = pbutil.StaticBuffer{3, 7, 224, 1} // 224

// 你不在联盟中
var ERR_REMOVE_GUILD_WORKSHOP_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{4, 7, 225, 1, 1} // 225-1

// 你没有权限操作
var ERR_REMOVE_GUILD_WORKSHOP_FAIL_DENY = pbutil.StaticBuffer{4, 7, 225, 1, 2} // 225-2

// 联盟工坊不存在
var ERR_REMOVE_GUILD_WORKSHOP_FAIL_BASE_NOT_EXIST = pbutil.StaticBuffer{4, 7, 225, 1, 3} // 225-3

// 服务器忙，请稍后再试
var ERR_REMOVE_GUILD_WORKSHOP_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 7, 225, 1, 4} // 225-4

func NewS2cUpdateGuildWorkshopPrizeCountMsg(count int32) pbutil.Buffer {
	msg := &S2CUpdateGuildWorkshopPrizeCountProto{
		Count: count,
	}
	return NewS2cUpdateGuildWorkshopPrizeCountProtoMsg(msg)
}

var s2c_update_guild_workshop_prize_count = [...]byte{7, 222, 1} // 222
func NewS2cUpdateGuildWorkshopPrizeCountProtoMsg(object *S2CUpdateGuildWorkshopPrizeCountProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_guild_workshop_prize_count[:], "s2c_update_guild_workshop_prize_count")

}

func NewS2cUpdateHeroBuildWorkshopTimesMsg(times int32) pbutil.Buffer {
	msg := &S2CUpdateHeroBuildWorkshopTimesProto{
		Times: times,
	}
	return NewS2cUpdateHeroBuildWorkshopTimesProtoMsg(msg)
}

var s2c_update_hero_build_workshop_times = [...]byte{7, 226, 1} // 226
func NewS2cUpdateHeroBuildWorkshopTimesProtoMsg(object *S2CUpdateHeroBuildWorkshopTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_hero_build_workshop_times[:], "s2c_update_hero_build_workshop_times")

}

func NewS2cUpdateHeroOutputWorkshopTimesMsg(start_recovey_time int32) pbutil.Buffer {
	msg := &S2CUpdateHeroOutputWorkshopTimesProto{
		StartRecoveyTime: start_recovey_time,
	}
	return NewS2cUpdateHeroOutputWorkshopTimesProtoMsg(msg)
}

var s2c_update_hero_output_workshop_times = [...]byte{7, 227, 1} // 227
func NewS2cUpdateHeroOutputWorkshopTimesProtoMsg(object *S2CUpdateHeroOutputWorkshopTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_hero_output_workshop_times[:], "s2c_update_hero_output_workshop_times")

}

func NewS2cCatchGuildWorkshopLogsMsg(version int32, logs []*shared_proto.GuildWorkshopLogProto) pbutil.Buffer {
	msg := &S2CCatchGuildWorkshopLogsProto{
		Version: version,
		Logs:    logs,
	}
	return NewS2cCatchGuildWorkshopLogsProtoMsg(msg)
}

var s2c_catch_guild_workshop_logs = [...]byte{7, 229, 1} // 229
func NewS2cCatchGuildWorkshopLogsProtoMsg(object *S2CCatchGuildWorkshopLogsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_catch_guild_workshop_logs[:], "s2c_catch_guild_workshop_logs")

}

// 没有联盟
var ERR_CATCH_GUILD_WORKSHOP_LOGS_FAIL_NO_GUILD = pbutil.StaticBuffer{4, 7, 230, 1, 1} // 230-1

// 没有联盟工坊
var ERR_CATCH_GUILD_WORKSHOP_LOGS_FAIL_NOT_EXIST = pbutil.StaticBuffer{4, 7, 230, 1, 2} // 230-2

func NewS2cGetSelfBaozMsg(exist bool, base_id []byte, base_x int32, base_y int32, expire_time int32) pbutil.Buffer {
	msg := &S2CGetSelfBaozProto{
		Exist:      exist,
		BaseId:     base_id,
		BaseX:      base_x,
		BaseY:      base_y,
		ExpireTime: expire_time,
	}
	return NewS2cGetSelfBaozProtoMsg(msg)
}

var s2c_get_self_baoz = [...]byte{7, 233, 1} // 233
func NewS2cGetSelfBaozProtoMsg(object *S2CGetSelfBaozProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_self_baoz[:], "s2c_get_self_baoz")

}
