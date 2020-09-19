package military

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
	MODULE_ID = 4

	C2S_RECRUIT_SOLDIER = 2

	C2S_RECRUIT_SOLDIER_V2 = 120

	C2S_HEAL_WOUNDED_SOLDIER = 6

	C2S_CAPTAIN_CHANGE_SOLDIER = 9

	C2S_CAPTAIN_FULL_SOLDIER = 66

	C2S_FORCE_ADD_SOLDIER = 149

	C2S_FIGHT = 12

	C2S_MULTI_FIGHT = 101

	C2S_FIGHTX = 198

	C2S_UPGRADE_SOLDIER_LEVEL = 15

	C2S_RECRUIT_CAPTAIN_V2 = 109

	C2S_RANDOM_CAPTAIN_HEAD = 176

	C2S_RECRUIT_CAPTAIN_SEEKER = 146

	C2S_SET_DEFENSE_TROOP = 106

	C2S_CLEAR_DEFENSE_TROOP_DEFEATED_MAIL = 129

	C2S_SET_DEFENSER_AUTO_FULL_SOLDIER = 188

	C2S_USE_COPY_DEFENSER_GOODS = 193

	C2S_SELL_SEEK_CAPTAIN = 34

	C2S_SET_MULTI_CAPTAIN_INDEX = 45

	C2S_SET_PVE_CAPTAIN = 143

	C2S_FIRE_CAPTAIN = 38

	C2S_CAPTAIN_REFINED = 48

	C2S_CAPTAIN_ENHANCE = 206

	C2S_CHANGE_CAPTAIN_NAME = 82

	C2S_CHANGE_CAPTAIN_RACE = 85

	C2S_CAPTAIN_REBIRTH_PREVIEW = 89

	C2S_CAPTAIN_REBIRTH = 92

	C2S_CAPTAIN_PROGRESS = 210

	C2S_CAPTAIN_REBIRTH_MIAO_CD = 166

	C2S_COLLECT_CAPTAIN_TRAINING_EXP = 136

	C2S_CAPTAIN_TRAIN_EXP = 213

	C2S_CAPTAIN_CAN_COLLECT_EXP = 261

	C2S_USE_TRAINING_EXP_GOODS = 140

	C2S_USE_LEVEL_EXP_GOODS = 216

	C2S_USE_LEVEL_EXP_GOODS2 = 243

	C2S_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP = 255

	C2S_GET_MAX_RECRUIT_SOLDIER = 74

	C2S_GET_MAX_HEAL_SOLDIER = 76

	C2S_JIU_GUAN_CONSULT = 112

	C2S_JIU_GUAN_REFRESH = 117

	C2S_UNLOCK_CAPTAIN_RESTRAINT_SPELL = 125

	C2S_GET_CAPTAIN_STAT_DETAILS = 132

	C2S_CAPTAIN_STAT_DETAILS = 219

	C2S_UPDATE_CAPTAIN_OFFICIAL = 169

	C2S_SET_CAPTAIN_OFFICIAL = 222

	C2S_LEAVE_CAPTAIN_OFFICIAL = 172

	C2S_USE_GONG_XUN_GOODS = 185

	C2S_USE_GONGXUN_GOODS = 228

	C2S_CLOSE_FIGHT_GUIDE = 181

	C2S_VIEW_OTHER_HERO_CAPTAIN = 190

	C2S_CAPTAIN_BORN = 231

	C2S_CAPTAIN_UPSTAR = 234

	C2S_CAPTAIN_EXCHANGE = 268

	C2S_NOTICE_CAPTAIN_HAS_VIEWED = 252

	C2S_ACTIVATE_CAPTAIN_FRIENDSHIP = 265

	C2S_NOTICE_OFFICIAL_HAS_VIEWED = 272
)

func NewS2cUpdateSoldierCapcityMsg(soldier_capcity int32, wounded_soldier_capcity int32, new_soldier_capcity int32, new_soldier_output int32, new_recruit_soldier_count int32) pbutil.Buffer {
	msg := &S2CUpdateSoldierCapcityProto{
		SoldierCapcity:         soldier_capcity,
		WoundedSoldierCapcity:  wounded_soldier_capcity,
		NewSoldierCapcity:      new_soldier_capcity,
		NewSoldierOutput:       new_soldier_output,
		NewRecruitSoldierCount: new_recruit_soldier_count,
	}
	return NewS2cUpdateSoldierCapcityProtoMsg(msg)
}

var s2c_update_soldier_capcity = [...]byte{4, 1} // 1
func NewS2cUpdateSoldierCapcityProtoMsg(object *S2CUpdateSoldierCapcityProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_soldier_capcity[:], "s2c_update_soldier_capcity")

}

func NewS2cRecruitSoldierMsg(new_soldier int32, free_soldier int32) pbutil.Buffer {
	msg := &S2CRecruitSoldierProto{
		NewSoldier:  new_soldier,
		FreeSoldier: free_soldier,
	}
	return NewS2cRecruitSoldierProtoMsg(msg)
}

var s2c_recruit_soldier = [...]byte{4, 3} // 3
func NewS2cRecruitSoldierProtoMsg(object *S2CRecruitSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_recruit_soldier[:], "s2c_recruit_soldier")

}

// 新兵数量没这么多
var ERR_RECRUIT_SOLDIER_FAIL_SOLDIER_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 4, 1} // 4-1

// 超出最大士兵上限
var ERR_RECRUIT_SOLDIER_FAIL_SOLDIER_CAPCITY_OVERFLOW = pbutil.StaticBuffer{3, 4, 4, 2} // 4-2

// 消耗资源不足
var ERR_RECRUIT_SOLDIER_FAIL_RESOURCE_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 4, 3} // 4-3

// 招募数量无效
var ERR_RECRUIT_SOLDIER_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 4, 4, 4} // 4-4

// 服务器忙，请稍后再试
var ERR_RECRUIT_SOLDIER_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 4, 4, 5} // 4-5

func NewS2cRecruitSoldierV2Msg(add_soldier int32) pbutil.Buffer {
	msg := &S2CRecruitSoldierV2Proto{
		AddSoldier: add_soldier,
	}
	return NewS2cRecruitSoldierV2ProtoMsg(msg)
}

var s2c_recruit_soldier_v2 = [...]byte{4, 121} // 121
func NewS2cRecruitSoldierV2ProtoMsg(object *S2CRecruitSoldierV2Proto) pbutil.Buffer {

	return newProtoMsg(object, s2c_recruit_soldier_v2[:], "s2c_recruit_soldier_v2")

}

// 超出最大士兵上限
var ERR_RECRUIT_SOLDIER_V2_FAIL_SOLDIER_CAPCITY_OVERFLOW = pbutil.StaticBuffer{3, 4, 122, 1} // 122-1

// 没有军营
var ERR_RECRUIT_SOLDIER_V2_FAIL_NO_JUN_YING = pbutil.StaticBuffer{3, 4, 122, 2} // 122-2

// 没有次数了
var ERR_RECRUIT_SOLDIER_V2_FAIL_NO_TIMES = pbutil.StaticBuffer{3, 4, 122, 3} // 122-3

// 服务器忙，请稍后再试
var ERR_RECRUIT_SOLDIER_V2_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 4, 122, 4} // 122-4

func NewS2cAutoRecoverSoldierMsg(free_soldier int32, captain_id []int32, captain_soldier_count []int32) pbutil.Buffer {
	msg := &S2CAutoRecoverSoldierProto{
		FreeSoldier:         free_soldier,
		CaptainId:           captain_id,
		CaptainSoldierCount: captain_soldier_count,
	}
	return NewS2cAutoRecoverSoldierProtoMsg(msg)
}

var s2c_auto_recover_soldier = [...]byte{4, 124} // 124
func NewS2cAutoRecoverSoldierProtoMsg(object *S2CAutoRecoverSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_auto_recover_soldier[:], "s2c_auto_recover_soldier")

}

func NewS2cRecruitSoldierTimesChangedMsg(start_recovery_time int32) pbutil.Buffer {
	msg := &S2CRecruitSoldierTimesChangedProto{
		StartRecoveryTime: start_recovery_time,
	}
	return NewS2cRecruitSoldierTimesChangedProtoMsg(msg)
}

var s2c_recruit_soldier_times_changed = [...]byte{4, 123} // 123
func NewS2cRecruitSoldierTimesChangedProtoMsg(object *S2CRecruitSoldierTimesChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_recruit_soldier_times_changed[:], "s2c_recruit_soldier_times_changed")

}

func NewS2cAddWoundedSoldierMsg(toAdd int32, total int32) pbutil.Buffer {
	msg := &S2CAddWoundedSoldierProto{
		ToAdd: toAdd,
		Total: total,
	}
	return NewS2cAddWoundedSoldierProtoMsg(msg)
}

var s2c_add_wounded_soldier = [...]byte{4, 5} // 5
func NewS2cAddWoundedSoldierProtoMsg(object *S2CAddWoundedSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_wounded_soldier[:], "s2c_add_wounded_soldier")

}

func NewS2cHealWoundedSoldierMsg(count int32) pbutil.Buffer {
	msg := &S2CHealWoundedSoldierProto{
		Count: count,
	}
	return NewS2cHealWoundedSoldierProtoMsg(msg)
}

var s2c_heal_wounded_soldier = [...]byte{4, 7} // 7
func NewS2cHealWoundedSoldierProtoMsg(object *S2CHealWoundedSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_heal_wounded_soldier[:], "s2c_heal_wounded_soldier")

}

// 伤兵数量没这么多
var ERR_HEAL_WOUNDED_SOLDIER_FAIL_SOLDIER_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 8, 1} // 8-1

// 超出最大士兵上限
var ERR_HEAL_WOUNDED_SOLDIER_FAIL_SOLDIER_CAPCITY_OVERFLOW = pbutil.StaticBuffer{3, 4, 8, 2} // 8-2

// 消耗资源不足
var ERR_HEAL_WOUNDED_SOLDIER_FAIL_RESOURCE_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 8, 3} // 8-3

// 治疗数量无效
var ERR_HEAL_WOUNDED_SOLDIER_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 4, 8, 4} // 8-4

// 服务器忙，请稍后再试
var ERR_HEAL_WOUNDED_SOLDIER_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 4, 8, 5} // 8-5

func NewS2cCaptainChangeSoldierMsg(id int32, soldier int32, fight_amount int32, free_soldier int32) pbutil.Buffer {
	msg := &S2CCaptainChangeSoldierProto{
		Id:          id,
		Soldier:     soldier,
		FightAmount: fight_amount,
		FreeSoldier: free_soldier,
	}
	return NewS2cCaptainChangeSoldierProtoMsg(msg)
}

var s2c_captain_change_soldier = [...]byte{4, 10} // 10
func NewS2cCaptainChangeSoldierProtoMsg(object *S2CCaptainChangeSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_change_soldier[:], "s2c_captain_change_soldier")

}

// 你还没有这个武将
var ERR_CAPTAIN_CHANGE_SOLDIER_FAIL_NOT_OWNER = pbutil.StaticBuffer{3, 4, 14, 1} // 14-1

// 武将出征中
var ERR_CAPTAIN_CHANGE_SOLDIER_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 4, 14, 3} // 14-3

// 服务器忙，请稍后再试
var ERR_CAPTAIN_CHANGE_SOLDIER_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 4, 14, 2} // 14-2

func NewS2cCaptainFullSoldierMsg(id []int32, soldier []int32, fight_amount []int32, free_soldier int32) pbutil.Buffer {
	msg := &S2CCaptainFullSoldierProto{
		Id:          id,
		Soldier:     soldier,
		FightAmount: fight_amount,
		FreeSoldier: free_soldier,
	}
	return NewS2cCaptainFullSoldierProtoMsg(msg)
}

var s2c_captain_full_soldier = [...]byte{4, 67} // 67
func NewS2cCaptainFullSoldierProtoMsg(object *S2CCaptainFullSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_full_soldier[:], "s2c_captain_full_soldier")

}

// 武将Id不存在
var ERR_CAPTAIN_FULL_SOLDIER_FAIL_CAPTAIN_NOT_EXIST = pbutil.StaticBuffer{3, 4, 68, 3} // 68-3

// 武将id重复
var ERR_CAPTAIN_FULL_SOLDIER_FAIL_DUPLICATE = pbutil.StaticBuffer{3, 4, 68, 4} // 68-4

// 没有空闲士兵
var ERR_CAPTAIN_FULL_SOLDIER_FAIL_EMPTY_SOLDIER = pbutil.StaticBuffer{3, 4, 68, 5} // 68-5

// 武将出征中
var ERR_CAPTAIN_FULL_SOLDIER_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 4, 68, 6} // 68-6

// 服务器忙，请稍后再试
var ERR_CAPTAIN_FULL_SOLDIER_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 4, 68, 2} // 68-2

func NewS2cUpdateFreeSoldierMsg(start_time int32, capcity int32, output int32, overflow int32) pbutil.Buffer {
	msg := &S2CUpdateFreeSoldierProto{
		StartTime: start_time,
		Capcity:   capcity,
		Output:    output,
		Overflow:  overflow,
	}
	return NewS2cUpdateFreeSoldierProtoMsg(msg)
}

var s2c_update_free_soldier = [...]byte{4, 80} // 80
func NewS2cUpdateFreeSoldierProtoMsg(object *S2CUpdateFreeSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_free_soldier[:], "s2c_update_free_soldier")

}

func NewS2cForceAddSoldierMsg(times int32) pbutil.Buffer {
	msg := &S2CForceAddSoldierProto{
		Times: times,
	}
	return NewS2cForceAddSoldierProtoMsg(msg)
}

var s2c_force_add_soldier = [...]byte{4, 150, 1} // 150
func NewS2cForceAddSoldierProtoMsg(object *S2CForceAddSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_force_add_soldier[:], "s2c_force_add_soldier")

}

// 军营建筑未解锁
var ERR_FORCE_ADD_SOLDIER_FAIL_LOCKED = pbutil.StaticBuffer{4, 4, 151, 1, 4} // 151-4

// 强征次数上限
var ERR_FORCE_ADD_SOLDIER_FAIL_TIMES_LIMIT = pbutil.StaticBuffer{4, 4, 151, 1, 1} // 151-1

// 强征消耗不足
var ERR_FORCE_ADD_SOLDIER_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 151, 1, 2} // 151-2

// 士兵数已满
var ERR_FORCE_ADD_SOLDIER_FAIL_FULL_SOLDIER = pbutil.StaticBuffer{4, 4, 151, 1, 3} // 151-3

func NewS2cCaptainChangeDataMsg(id int32, soldier int32, max_soldier int32, fight_amount int32) pbutil.Buffer {
	msg := &S2CCaptainChangeDataProto{
		Id:          id,
		Soldier:     soldier,
		MaxSoldier:  max_soldier,
		FightAmount: fight_amount,
	}
	return NewS2cCaptainChangeDataProtoMsg(msg)
}

var s2c_captain_change_data = [...]byte{4, 11} // 11
func NewS2cCaptainChangeDataProtoMsg(object *S2CCaptainChangeDataProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_change_data[:], "s2c_captain_change_data")

}

func NewS2cFightMsg(replay []byte) pbutil.Buffer {
	msg := &S2CFightProto{
		Replay: replay,
	}
	return NewS2cFightProtoMsg(msg)
}

func NewS2cFightMarshalMsg(replay marshaler) pbutil.Buffer {
	msg := &S2CFightProto{
		Replay: safeMarshal(replay),
	}
	return NewS2cFightProtoMsg(msg)
}

var s2c_fight = [...]byte{4, 13} // 13
func NewS2cFightProtoMsg(object *S2CFightProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_fight[:], "s2c_fight")

}

func NewS2cMultiFightMsg(replay []byte) pbutil.Buffer {
	msg := &S2CMultiFightProto{
		Replay: replay,
	}
	return NewS2cMultiFightProtoMsg(msg)
}

func NewS2cMultiFightMarshalMsg(replay marshaler) pbutil.Buffer {
	msg := &S2CMultiFightProto{
		Replay: safeMarshal(replay),
	}
	return NewS2cMultiFightProtoMsg(msg)
}

var s2c_multi_fight = [...]byte{4, 102} // 102
func NewS2cMultiFightProtoMsg(object *S2CMultiFightProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_multi_fight[:], "s2c_multi_fight")

}

func NewS2cFightxMsg(replay []byte) pbutil.Buffer {
	msg := &S2CFightxProto{
		Replay: replay,
	}
	return NewS2cFightxProtoMsg(msg)
}

func NewS2cFightxMarshalMsg(replay marshaler) pbutil.Buffer {
	msg := &S2CFightxProto{
		Replay: safeMarshal(replay),
	}
	return NewS2cFightxProtoMsg(msg)
}

var s2c_fightx = [...]byte{4, 199, 1} // 199
func NewS2cFightxProtoMsg(object *S2CFightxProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_fightx[:], "s2c_fightx")

}

func NewS2cUpgradeSoldierLevelMsg(level int32) pbutil.Buffer {
	msg := &S2CUpgradeSoldierLevelProto{
		Level: level,
	}
	return NewS2cUpgradeSoldierLevelProtoMsg(msg)
}

var s2c_upgrade_soldier_level = [...]byte{4, 16} // 16
func NewS2cUpgradeSoldierLevelProtoMsg(object *S2CUpgradeSoldierLevelProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_upgrade_soldier_level[:], "s2c_upgrade_soldier_level")

}

// 已经到最高级了
var ERR_UPGRADE_SOLDIER_LEVEL_FAIL_MAX_LEVEL = pbutil.StaticBuffer{3, 4, 17, 1} // 17-1

// 钱不够
var ERR_UPGRADE_SOLDIER_LEVEL_FAIL_RES_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 17, 2} // 17-2

// 军营等级不够
var ERR_UPGRADE_SOLDIER_LEVEL_FAIL_JUN_YING_LEVEL_TOO_LOW = pbutil.StaticBuffer{3, 4, 17, 4} // 17-4

// 服务器忙，请稍后再试
var ERR_UPGRADE_SOLDIER_LEVEL_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 4, 17, 3} // 17-3

func NewS2cRecruitCaptainV2Msg(captain []byte, captain_index int32) pbutil.Buffer {
	msg := &S2CRecruitCaptainV2Proto{
		Captain:      captain,
		CaptainIndex: captain_index,
	}
	return NewS2cRecruitCaptainV2ProtoMsg(msg)
}

func NewS2cRecruitCaptainV2MarshalMsg(captain marshaler, captain_index int32) pbutil.Buffer {
	msg := &S2CRecruitCaptainV2Proto{
		Captain:      safeMarshal(captain),
		CaptainIndex: captain_index,
	}
	return NewS2cRecruitCaptainV2ProtoMsg(msg)
}

var s2c_recruit_captain_v2 = [...]byte{4, 110} // 110
func NewS2cRecruitCaptainV2ProtoMsg(object *S2CRecruitCaptainV2Proto) pbutil.Buffer {

	return newProtoMsg(object, s2c_recruit_captain_v2[:], "s2c_recruit_captain_v2")

}

// 可招募武将已达上限
var ERR_RECRUIT_CAPTAIN_V2_FAIL_FULL = pbutil.StaticBuffer{3, 4, 111, 1} // 111-1

// 君主等级太低
var ERR_RECRUIT_CAPTAIN_V2_FAIL_HERO_LEVEL_TOO_LOW = pbutil.StaticBuffer{3, 4, 111, 2} // 111-2

// 当前没有武将可以招募
var ERR_RECRUIT_CAPTAIN_V2_FAIL_NO_CAPTAIN_CAN_RECRUIT = pbutil.StaticBuffer{3, 4, 111, 3} // 111-3

// 没有空闲的队伍(出征中)
var ERR_RECRUIT_CAPTAIN_V2_FAIL_NO_FREE_TROOP = pbutil.StaticBuffer{3, 4, 111, 4} // 111-4

// 服务器忙，请稍后再试
var ERR_RECRUIT_CAPTAIN_V2_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 4, 111, 5} // 111-5

func NewS2cRandomCaptainHeadMsg(head []string) pbutil.Buffer {
	msg := &S2CRandomCaptainHeadProto{
		Head: head,
	}
	return NewS2cRandomCaptainHeadProtoMsg(msg)
}

var s2c_random_captain_head = [...]byte{4, 177, 1} // 177
func NewS2cRandomCaptainHeadProtoMsg(object *S2CRandomCaptainHeadProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_random_captain_head[:], "s2c_random_captain_head")

}

func NewS2cRecruitCaptainSeekerMsg(index int32, captain []byte, captain_index int32) pbutil.Buffer {
	msg := &S2CRecruitCaptainSeekerProto{
		Index:        index,
		Captain:      captain,
		CaptainIndex: captain_index,
	}
	return NewS2cRecruitCaptainSeekerProtoMsg(msg)
}

func NewS2cRecruitCaptainSeekerMarshalMsg(index int32, captain marshaler, captain_index int32) pbutil.Buffer {
	msg := &S2CRecruitCaptainSeekerProto{
		Index:        index,
		Captain:      safeMarshal(captain),
		CaptainIndex: captain_index,
	}
	return NewS2cRecruitCaptainSeekerProtoMsg(msg)
}

var s2c_recruit_captain_seeker = [...]byte{4, 147, 1} // 147
func NewS2cRecruitCaptainSeekerProtoMsg(object *S2CRecruitCaptainSeekerProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_recruit_captain_seeker[:], "s2c_recruit_captain_seeker")

}

// 无效的编号
var ERR_RECRUIT_CAPTAIN_SEEKER_FAIL_INVALID_INDEX = pbutil.StaticBuffer{4, 4, 148, 1, 1} // 148-1

// 可招募武将已达上限
var ERR_RECRUIT_CAPTAIN_SEEKER_FAIL_FULL = pbutil.StaticBuffer{4, 4, 148, 1, 2} // 148-2

// 君主等级太低
var ERR_RECRUIT_CAPTAIN_SEEKER_FAIL_HERO_LEVEL_TOO_LOW = pbutil.StaticBuffer{4, 4, 148, 1, 3} // 148-3

// 当前没有武将可以招募
var ERR_RECRUIT_CAPTAIN_SEEKER_FAIL_NO_CAPTAIN_CAN_RECRUIT = pbutil.StaticBuffer{4, 4, 148, 1, 4} // 148-4

// 没有空闲的队伍(出征中)
var ERR_RECRUIT_CAPTAIN_SEEKER_FAIL_NO_FREE_TROOP = pbutil.StaticBuffer{4, 4, 148, 1, 5} // 148-5

// 服务器忙，请稍后再试
var ERR_RECRUIT_CAPTAIN_SEEKER_FAIL_SERVER_BUSY = pbutil.StaticBuffer{4, 4, 148, 1, 6} // 148-6

func NewS2cSetDefenseTroopMsg(is_tent bool, troop_index int32) pbutil.Buffer {
	msg := &S2CSetDefenseTroopProto{
		IsTent:     is_tent,
		TroopIndex: troop_index,
	}
	return NewS2cSetDefenseTroopProtoMsg(msg)
}

var s2c_set_defense_troop = [...]byte{4, 107} // 107
func NewS2cSetDefenseTroopProtoMsg(object *S2CSetDefenseTroopProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_defense_troop[:], "s2c_set_defense_troop")

}

// 无效的队伍index
var ERR_SET_DEFENSE_TROOP_FAIL_INVALID_TROOP_INDEX = pbutil.StaticBuffer{3, 4, 108, 1} // 108-1

func NewS2cSetDenfeseTroopDefeatedMailMsg(mail []byte, is_tent bool) pbutil.Buffer {
	msg := &S2CSetDenfeseTroopDefeatedMailProto{
		Mail:   mail,
		IsTent: is_tent,
	}
	return NewS2cSetDenfeseTroopDefeatedMailProtoMsg(msg)
}

var s2c_set_denfese_troop_defeated_mail = [...]byte{4, 128, 1} // 128
func NewS2cSetDenfeseTroopDefeatedMailProtoMsg(object *S2CSetDenfeseTroopDefeatedMailProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_set_denfese_troop_defeated_mail[:], "s2c_set_denfese_troop_defeated_mail")

}

func NewS2cClearDefenseTroopDefeatedMailMsg(is_tent bool) pbutil.Buffer {
	msg := &S2CClearDefenseTroopDefeatedMailProto{
		IsTent: is_tent,
	}
	return NewS2cClearDefenseTroopDefeatedMailProtoMsg(msg)
}

var s2c_clear_defense_troop_defeated_mail = [...]byte{4, 130, 1} // 130
func NewS2cClearDefenseTroopDefeatedMailProtoMsg(object *S2CClearDefenseTroopDefeatedMailProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_clear_defense_troop_defeated_mail[:], "s2c_clear_defense_troop_defeated_mail")

}

func NewS2cSetDefenserAutoFullSoldierMsg(dont bool) pbutil.Buffer {
	msg := &S2CSetDefenserAutoFullSoldierProto{
		Dont: dont,
	}
	return NewS2cSetDefenserAutoFullSoldierProtoMsg(msg)
}

var s2c_set_defenser_auto_full_soldier = [...]byte{4, 189, 1} // 189
func NewS2cSetDefenserAutoFullSoldierProtoMsg(object *S2CSetDefenserAutoFullSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_defenser_auto_full_soldier[:], "s2c_set_defenser_auto_full_soldier")

}

// vip等级不够
var ERR_SET_DEFENSER_AUTO_FULL_SOLDIER_FAIL_VIP_LEVEL_LIMIT = pbutil.StaticBuffer{4, 4, 136, 2, 1} // 264-1

func NewS2cUseCopyDefenserGoodsMsg(troop_index int32, end_time int32) pbutil.Buffer {
	msg := &S2CUseCopyDefenserGoodsProto{
		TroopIndex: troop_index,
		EndTime:    end_time,
	}
	return NewS2cUseCopyDefenserGoodsProtoMsg(msg)
}

var s2c_use_copy_defenser_goods = [...]byte{4, 194, 1} // 194
func NewS2cUseCopyDefenserGoodsProtoMsg(object *S2CUseCopyDefenserGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_copy_defenser_goods[:], "s2c_use_copy_defenser_goods")

}

// 无效的物品id
var ERR_USE_COPY_DEFENSER_GOODS_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 4, 195, 1, 1} // 195-1

// 无效的队伍编号
var ERR_USE_COPY_DEFENSER_GOODS_FAIL_INVALID_TROOP = pbutil.StaticBuffer{4, 4, 195, 1, 2} // 195-2

// 物品个数不足
var ERR_USE_COPY_DEFENSER_GOODS_FAIL_COUNT_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 195, 1, 3} // 195-3

// 物品自动购买，消耗不足
var ERR_USE_COPY_DEFENSER_GOODS_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 195, 1, 4} // 195-4

func NewS2cUpdateCopyDefenserMsg(soldier int32, total_soldier int32, fight_amount int32) pbutil.Buffer {
	msg := &S2CUpdateCopyDefenserProto{
		Soldier:      soldier,
		TotalSoldier: total_soldier,
		FightAmount:  fight_amount,
	}
	return NewS2cUpdateCopyDefenserProtoMsg(msg)
}

var s2c_update_copy_defenser = [...]byte{4, 196, 1} // 196
func NewS2cUpdateCopyDefenserProtoMsg(object *S2CUpdateCopyDefenserProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_copy_defenser[:], "s2c_update_copy_defenser")

}

var REMOVE_COPY_DEFENSER_S2C = pbutil.StaticBuffer{3, 4, 197, 1} // 197

func NewS2cSellSeekCaptainMsg(index int32) pbutil.Buffer {
	msg := &S2CSellSeekCaptainProto{
		Index: index,
	}
	return NewS2cSellSeekCaptainProtoMsg(msg)
}

var s2c_sell_seek_captain = [...]byte{4, 35} // 35
func NewS2cSellSeekCaptainProtoMsg(object *S2CSellSeekCaptainProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_sell_seek_captain[:], "s2c_sell_seek_captain")

}

// 武将id无效
var ERR_SELL_SEEK_CAPTAIN_FAIL_INVALID_INDEX = pbutil.StaticBuffer{3, 4, 37, 1} // 37-1

// 寻访武将列表为空
var ERR_SELL_SEEK_CAPTAIN_FAIL_EMPTY = pbutil.StaticBuffer{3, 4, 37, 2} // 37-2

// 服务器忙，请稍后再试
var ERR_SELL_SEEK_CAPTAIN_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 4, 37, 3} // 37-3

func NewS2cSetMultiCaptainIndexMsg(index int32, id []int32, x_index []int32) pbutil.Buffer {
	msg := &S2CSetMultiCaptainIndexProto{
		Index:  index,
		Id:     id,
		XIndex: x_index,
	}
	return NewS2cSetMultiCaptainIndexProtoMsg(msg)
}

var s2c_set_multi_captain_index = [...]byte{4, 46} // 46
func NewS2cSetMultiCaptainIndexProtoMsg(object *S2CSetMultiCaptainIndexProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_multi_captain_index[:], "s2c_set_multi_captain_index")

}

// 无效的序号
var ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_INVALID_INDEX = pbutil.StaticBuffer{3, 4, 47, 1} // 47-1

// 无效的武将id
var ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 4, 47, 2} // 47-2

// 出征武将不能修改
var ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 4, 47, 4} // 47-4

// 武将已经在别的编队中
var ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_OTHER_INDEX = pbutil.StaticBuffer{3, 4, 47, 3} // 47-3

// 必须包含所有武将
var ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_MUST_ALL = pbutil.StaticBuffer{3, 4, 47, 5} // 47-5

func NewS2cSetPveCaptainMsg(troop []byte) pbutil.Buffer {
	msg := &S2CSetPveCaptainProto{
		Troop: troop,
	}
	return NewS2cSetPveCaptainProtoMsg(msg)
}

var s2c_set_pve_captain = [...]byte{4, 144, 1} // 144
func NewS2cSetPveCaptainProtoMsg(object *S2CSetPveCaptainProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_pve_captain[:], "s2c_set_pve_captain")

}

// 无效的队伍类型
var ERR_SET_PVE_CAPTAIN_FAIL_INVALID_PVE_TYPE = pbutil.StaticBuffer{4, 4, 145, 1, 1} // 145-1

// 武将数量非法
var ERR_SET_PVE_CAPTAIN_FAIL_INVALID_CAPTAIN_COUNT = pbutil.StaticBuffer{4, 4, 145, 1, 2} // 145-2

// 武将没找到
var ERR_SET_PVE_CAPTAIN_FAIL_INVALID_ID = pbutil.StaticBuffer{4, 4, 145, 1, 3} // 145-3

// 存在重复的武将id
var ERR_SET_PVE_CAPTAIN_FAIL_DUP_CAPTAIN_ID = pbutil.StaticBuffer{4, 4, 145, 1, 4} // 145-4

// 不可以将队伍设置为没有武将出战
var ERR_SET_PVE_CAPTAIN_FAIL_NO_CAPTAIN = pbutil.StaticBuffer{4, 4, 145, 1, 5} // 145-5

// 无效的武将横向位置
var ERR_SET_PVE_CAPTAIN_FAIL_INVALID_X_INDEX = pbutil.StaticBuffer{4, 4, 145, 1, 7} // 145-7

func NewS2cFireCaptainMsg(id int32) pbutil.Buffer {
	msg := &S2CFireCaptainProto{
		Id: id,
	}
	return NewS2cFireCaptainProtoMsg(msg)
}

var s2c_fire_captain = [...]byte{4, 39} // 39
func NewS2cFireCaptainProtoMsg(object *S2CFireCaptainProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fire_captain[:], "s2c_fire_captain")

}

// 武将id没找到
var ERR_FIRE_CAPTAIN_FAIL_ID_NOT_FOUND = pbutil.StaticBuffer{3, 4, 40, 1} // 40-1

// 武将出征中
var ERR_FIRE_CAPTAIN_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 4, 40, 2} // 40-2

// 等级超过30级，不能解雇
var ERR_FIRE_CAPTAIN_FAIL_LEVEL_LIMIT = pbutil.StaticBuffer{3, 4, 40, 3} // 40-3

// 服务器忙，请稍后再试
var ERR_FIRE_CAPTAIN_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 4, 40, 4} // 40-4

// 装备背包满了
var ERR_FIRE_CAPTAIN_FAIL_DEPOT_EQUIPMENT_FULL = pbutil.StaticBuffer{3, 4, 40, 5} // 40-5

// 宝石背包满了
var ERR_FIRE_CAPTAIN_FAIL_DEPOT_GEM_FULL = pbutil.StaticBuffer{3, 4, 40, 6} // 40-6

// 在重楼密室队伍中
var ERR_FIRE_CAPTAIN_FAIL_SECRET_TOWER = pbutil.StaticBuffer{3, 4, 40, 7} // 40-7

func NewS2cCaptainRefinedMsg(captain int32, exp int32) pbutil.Buffer {
	msg := &S2CCaptainRefinedProto{
		Captain: captain,
		Exp:     exp,
	}
	return NewS2cCaptainRefinedProtoMsg(msg)
}

var s2c_captain_refined = [...]byte{4, 49} // 49
func NewS2cCaptainRefinedProtoMsg(object *S2CCaptainRefinedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_refined[:], "s2c_captain_refined")

}

// 无效的武将id
var ERR_CAPTAIN_REFINED_FAIL_INVALID_CAPTAIN = pbutil.StaticBuffer{3, 4, 50, 1} // 50-1

// 无效的强化符id
var ERR_CAPTAIN_REFINED_FAIL_INVALID_GOODS = pbutil.StaticBuffer{3, 4, 50, 2} // 50-2

// 无效的使用个数
var ERR_CAPTAIN_REFINED_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 4, 50, 3} // 50-3

// 武将出征中
var ERR_CAPTAIN_REFINED_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 4, 50, 4} // 50-4

// 成长值到达转生上限
var ERR_CAPTAIN_REFINED_FAIL_REBIRTH_LIMIT = pbutil.StaticBuffer{3, 4, 50, 5} // 50-5

func NewS2cCaptainEnhanceMsg(captain int32, ability int32, ability_exp int32, quality int32) pbutil.Buffer {
	msg := &S2CCaptainEnhanceProto{
		Captain:    captain,
		Ability:    ability,
		AbilityExp: ability_exp,
		Quality:    quality,
	}
	return NewS2cCaptainEnhanceProtoMsg(msg)
}

var s2c_captain_enhance = [...]byte{4, 207, 1} // 207
func NewS2cCaptainEnhanceProtoMsg(object *S2CCaptainEnhanceProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_enhance[:], "s2c_captain_enhance")

}

// 无效的武将id
var ERR_CAPTAIN_ENHANCE_FAIL_INVALID_CAPTAIN = pbutil.StaticBuffer{4, 4, 208, 1, 1} // 208-1

// 无效的强化符id
var ERR_CAPTAIN_ENHANCE_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 4, 208, 1, 2} // 208-2

// 无效的使用个数
var ERR_CAPTAIN_ENHANCE_FAIL_INVALID_COUNT = pbutil.StaticBuffer{4, 4, 208, 1, 3} // 208-3

// 武将出征中
var ERR_CAPTAIN_ENHANCE_FAIL_OUTSIDE = pbutil.StaticBuffer{4, 4, 208, 1, 4} // 208-4

// 成长值到达转生上限
var ERR_CAPTAIN_ENHANCE_FAIL_REBIRTH_LIMIT = pbutil.StaticBuffer{4, 4, 208, 1, 5} // 208-5

// 成长值上限
var ERR_CAPTAIN_ENHANCE_FAIL_ABILITY_MAX = pbutil.StaticBuffer{4, 4, 208, 1, 6} // 208-6

func NewS2cCaptainRefinedUpgradeMsg(captain int32, exp int32, ability int32, name []byte, quality int32) pbutil.Buffer {
	msg := &S2CCaptainRefinedUpgradeProto{
		Captain: captain,
		Exp:     exp,
		Ability: ability,
		Name:    name,
		Quality: quality,
	}
	return NewS2cCaptainRefinedUpgradeProtoMsg(msg)
}

var s2c_captain_refined_upgrade = [...]byte{4, 51} // 51
func NewS2cCaptainRefinedUpgradeProtoMsg(object *S2CCaptainRefinedUpgradeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_refined_upgrade[:], "s2c_captain_refined_upgrade")

}

func NewS2cUpdateAbilityExpMsg(captain int32, exp int32) pbutil.Buffer {
	msg := &S2CUpdateAbilityExpProto{
		Captain: captain,
		Exp:     exp,
	}
	return NewS2cUpdateAbilityExpProtoMsg(msg)
}

var s2c_update_ability_exp = [...]byte{4, 184, 1} // 184
func NewS2cUpdateAbilityExpProtoMsg(object *S2CUpdateAbilityExpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_ability_exp[:], "s2c_update_ability_exp")

}

func NewS2cUpdateCaptainExpMsg(captain int32, exp int32) pbutil.Buffer {
	msg := &S2CUpdateCaptainExpProto{
		Captain: captain,
		Exp:     exp,
	}
	return NewS2cUpdateCaptainExpProtoMsg(msg)
}

var s2c_update_captain_exp = [...]byte{4, 52} // 52
func NewS2cUpdateCaptainExpProtoMsg(object *S2CUpdateCaptainExpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_captain_exp[:], "s2c_update_captain_exp")

}

func NewS2cUpdateCaptainLevelMsg(captain int32, exp int32, level int32, name []byte, soldier_capcity int32) pbutil.Buffer {
	msg := &S2CUpdateCaptainLevelProto{
		Captain:        captain,
		Exp:            exp,
		Level:          level,
		Name:           name,
		SoldierCapcity: soldier_capcity,
	}
	return NewS2cUpdateCaptainLevelProtoMsg(msg)
}

var s2c_update_captain_level = [...]byte{4, 53} // 53
func NewS2cUpdateCaptainLevelProtoMsg(object *S2CUpdateCaptainLevelProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_captain_level[:], "s2c_update_captain_level")

}

func NewS2cCaptainLevelupMsg(captain int32, exp int32, level int32, soldier_capcity int32) pbutil.Buffer {
	msg := &S2CCaptainLevelupProto{
		Captain:        captain,
		Exp:            exp,
		Level:          level,
		SoldierCapcity: soldier_capcity,
	}
	return NewS2cCaptainLevelupProtoMsg(msg)
}

var s2c_captain_levelup = [...]byte{4, 209, 1} // 209
func NewS2cCaptainLevelupProtoMsg(object *S2CCaptainLevelupProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_levelup[:], "s2c_captain_levelup")

}

func NewS2cUpdateCaptainStatMsg(captain int32, total_stat []byte, fight_amount int32, full_fight_amount int32) pbutil.Buffer {
	msg := &S2CUpdateCaptainStatProto{
		Captain:         captain,
		TotalStat:       total_stat,
		FightAmount:     fight_amount,
		FullFightAmount: full_fight_amount,
	}
	return NewS2cUpdateCaptainStatProtoMsg(msg)
}

func NewS2cUpdateCaptainStatMarshalMsg(captain int32, total_stat marshaler, fight_amount int32, full_fight_amount int32) pbutil.Buffer {
	msg := &S2CUpdateCaptainStatProto{
		Captain:         captain,
		TotalStat:       safeMarshal(total_stat),
		FightAmount:     fight_amount,
		FullFightAmount: full_fight_amount,
	}
	return NewS2cUpdateCaptainStatProtoMsg(msg)
}

var s2c_update_captain_stat = [...]byte{4, 69} // 69
func NewS2cUpdateCaptainStatProtoMsg(object *S2CUpdateCaptainStatProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_captain_stat[:], "s2c_update_captain_stat")

}

func NewS2cChangeCaptainNameMsg(id int32, name string) pbutil.Buffer {
	msg := &S2CChangeCaptainNameProto{
		Id:   id,
		Name: name,
	}
	return NewS2cChangeCaptainNameProtoMsg(msg)
}

var s2c_change_captain_name = [...]byte{4, 83} // 83
func NewS2cChangeCaptainNameProtoMsg(object *S2CChangeCaptainNameProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_captain_name[:], "s2c_change_captain_name")

}

// 武将id无效
var ERR_CHANGE_CAPTAIN_NAME_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 4, 84, 1} // 84-1

// 武将名字无效
var ERR_CHANGE_CAPTAIN_NAME_FAIL_INVALID_NAME = pbutil.StaticBuffer{3, 4, 84, 2} // 84-2

// 新的名字跟原来的名字相同
var ERR_CHANGE_CAPTAIN_NAME_FAIL_SAME_NAME = pbutil.StaticBuffer{3, 4, 84, 3} // 84-3

// 新的名字已经被其他武将使用了
var ERR_CHANGE_CAPTAIN_NAME_FAIL_DUPLICATE_NAME = pbutil.StaticBuffer{3, 4, 84, 5} // 84-5

// 消耗不足
var ERR_CHANGE_CAPTAIN_NAME_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 84, 4} // 84-4

// 名字包含敏感词
var ERR_CHANGE_CAPTAIN_NAME_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 4, 84, 6} // 84-6

// 服务忙，请稍后重试
var ERR_CHANGE_CAPTAIN_NAME_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 4, 84, 7} // 84-7

func NewS2cChangeCaptainRaceMsg(id int32, race int32, cooldown int32, name []byte) pbutil.Buffer {
	msg := &S2CChangeCaptainRaceProto{
		Id:       id,
		Race:     race,
		Cooldown: cooldown,
		Name:     name,
	}
	return NewS2cChangeCaptainRaceProtoMsg(msg)
}

var s2c_change_captain_race = [...]byte{4, 86} // 86
func NewS2cChangeCaptainRaceProtoMsg(object *S2CChangeCaptainRaceProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_captain_race[:], "s2c_change_captain_race")

}

// 无效的武将id
var ERR_CHANGE_CAPTAIN_RACE_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 4, 87, 1} // 87-1

// 等级未达到转职要求
var ERR_CHANGE_CAPTAIN_RACE_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 87, 9} // 87-9

// 转职消耗不足
var ERR_CHANGE_CAPTAIN_RACE_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 87, 2} // 87-2

// 不支持点券购买
var ERR_CHANGE_CAPTAIN_RACE_FAIL_NOT_SUPPORT_YUANBAO = pbutil.StaticBuffer{3, 4, 87, 4} // 87-4

// 武将转职cd中
var ERR_CHANGE_CAPTAIN_RACE_FAIL_COOLDOWN = pbutil.StaticBuffer{3, 4, 87, 3} // 87-3

// 无效的职业类型
var ERR_CHANGE_CAPTAIN_RACE_FAIL_INVALID_RACE = pbutil.StaticBuffer{3, 4, 87, 5} // 87-5

// 新职业跟当前职业一样
var ERR_CHANGE_CAPTAIN_RACE_FAIL_SAME_RACE = pbutil.StaticBuffer{3, 4, 87, 6} // 87-6

// 武将出征中，不能转职
var ERR_CHANGE_CAPTAIN_RACE_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 4, 87, 7} // 87-7

// 服务器忙，请稍后再试
var ERR_CHANGE_CAPTAIN_RACE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 4, 87, 8} // 87-8

func NewS2cCaptainRebirthPreviewMsg(id int32, name []byte, rebirth_level int32, quality int32, ability int32, ability_limit int32, total_stat int32, soldier_capcity int32, add_stat []byte) pbutil.Buffer {
	msg := &S2CCaptainRebirthPreviewProto{
		Id:             id,
		Name:           name,
		RebirthLevel:   rebirth_level,
		Quality:        quality,
		Ability:        ability,
		AbilityLimit:   ability_limit,
		TotalStat:      total_stat,
		SoldierCapcity: soldier_capcity,
		AddStat:        add_stat,
	}
	return NewS2cCaptainRebirthPreviewProtoMsg(msg)
}

func NewS2cCaptainRebirthPreviewMarshalMsg(id int32, name []byte, rebirth_level int32, quality int32, ability int32, ability_limit int32, total_stat int32, soldier_capcity int32, add_stat marshaler) pbutil.Buffer {
	msg := &S2CCaptainRebirthPreviewProto{
		Id:             id,
		Name:           name,
		RebirthLevel:   rebirth_level,
		Quality:        quality,
		Ability:        ability,
		AbilityLimit:   ability_limit,
		TotalStat:      total_stat,
		SoldierCapcity: soldier_capcity,
		AddStat:        safeMarshal(add_stat),
	}
	return NewS2cCaptainRebirthPreviewProtoMsg(msg)
}

var s2c_captain_rebirth_preview = [...]byte{4, 90} // 90
func NewS2cCaptainRebirthPreviewProtoMsg(object *S2CCaptainRebirthPreviewProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_rebirth_preview[:], "s2c_captain_rebirth_preview")

}

// 无效的武将id
var ERR_CAPTAIN_REBIRTH_PREVIEW_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 4, 91, 1} // 91-1

// 已转生到最高等级
var ERR_CAPTAIN_REBIRTH_PREVIEW_FAIL_MAX_LEVEL = pbutil.StaticBuffer{3, 4, 91, 4} // 91-4

// 服务器忙，请稍后再试
var ERR_CAPTAIN_REBIRTH_PREVIEW_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 4, 91, 3} // 91-3

func NewS2cCaptainRebirthCdStartMsg(id int32, cd_endtime int32) pbutil.Buffer {
	msg := &S2CCaptainRebirthCdStartProto{
		Id:        id,
		CdEndtime: cd_endtime,
	}
	return NewS2cCaptainRebirthCdStartProtoMsg(msg)
}

var s2c_captain_rebirth_cd_start = [...]byte{4, 161, 1} // 161
func NewS2cCaptainRebirthCdStartProtoMsg(object *S2CCaptainRebirthCdStartProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_rebirth_cd_start[:], "s2c_captain_rebirth_cd_start")

}

func NewS2cCaptainRebirthMsg(id int32, name []byte, rebirth_level int32, rebirth_exp int32, quality int32, ability int32, ability_exp int32, ability_limit int32, soldier int32, soldier_capcity int32, total_stat []byte, fight_amount int32, full_fight_amount int32) pbutil.Buffer {
	msg := &S2CCaptainRebirthProto{
		Id:              id,
		Name:            name,
		RebirthLevel:    rebirth_level,
		RebirthExp:      rebirth_exp,
		Quality:         quality,
		Ability:         ability,
		AbilityExp:      ability_exp,
		AbilityLimit:    ability_limit,
		Soldier:         soldier,
		SoldierCapcity:  soldier_capcity,
		TotalStat:       total_stat,
		FightAmount:     fight_amount,
		FullFightAmount: full_fight_amount,
	}
	return NewS2cCaptainRebirthProtoMsg(msg)
}

func NewS2cCaptainRebirthMarshalMsg(id int32, name []byte, rebirth_level int32, rebirth_exp int32, quality int32, ability int32, ability_exp int32, ability_limit int32, soldier int32, soldier_capcity int32, total_stat marshaler, fight_amount int32, full_fight_amount int32) pbutil.Buffer {
	msg := &S2CCaptainRebirthProto{
		Id:              id,
		Name:            name,
		RebirthLevel:    rebirth_level,
		RebirthExp:      rebirth_exp,
		Quality:         quality,
		Ability:         ability,
		AbilityExp:      ability_exp,
		AbilityLimit:    ability_limit,
		Soldier:         soldier,
		SoldierCapcity:  soldier_capcity,
		TotalStat:       safeMarshal(total_stat),
		FightAmount:     fight_amount,
		FullFightAmount: full_fight_amount,
	}
	return NewS2cCaptainRebirthProtoMsg(msg)
}

var s2c_captain_rebirth = [...]byte{4, 93} // 93
func NewS2cCaptainRebirthProtoMsg(object *S2CCaptainRebirthProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_rebirth[:], "s2c_captain_rebirth")

}

// 无效的武将id
var ERR_CAPTAIN_REBIRTH_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 4, 94, 1} // 94-1

// 未达到转职等级
var ERR_CAPTAIN_REBIRTH_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 94, 2} // 94-2

// 已转生到最高等级
var ERR_CAPTAIN_REBIRTH_FAIL_MAX_LEVEL = pbutil.StaticBuffer{3, 4, 94, 4} // 94-4

// 服务器忙，请稍后再试
var ERR_CAPTAIN_REBIRTH_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 4, 94, 5} // 94-5

// 还在转生CD中
var ERR_CAPTAIN_REBIRTH_FAIL_IN_REBIRTHING = pbutil.StaticBuffer{3, 4, 94, 6} // 94-6

// 秒CD消耗不足
var ERR_CAPTAIN_REBIRTH_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 94, 7} // 94-7

func NewS2cCaptainProgressMsg(id int32, rebirth_level int32, rebirth_exp int32, quality int32, ability int32, ability_exp int32, ability_limit int32, soldier int32, soldier_capcity int32, total_stat []byte, fight_amount int32, full_fight_amount int32) pbutil.Buffer {
	msg := &S2CCaptainProgressProto{
		Id:              id,
		RebirthLevel:    rebirth_level,
		RebirthExp:      rebirth_exp,
		Quality:         quality,
		Ability:         ability,
		AbilityExp:      ability_exp,
		AbilityLimit:    ability_limit,
		Soldier:         soldier,
		SoldierCapcity:  soldier_capcity,
		TotalStat:       total_stat,
		FightAmount:     fight_amount,
		FullFightAmount: full_fight_amount,
	}
	return NewS2cCaptainProgressProtoMsg(msg)
}

func NewS2cCaptainProgressMarshalMsg(id int32, rebirth_level int32, rebirth_exp int32, quality int32, ability int32, ability_exp int32, ability_limit int32, soldier int32, soldier_capcity int32, total_stat marshaler, fight_amount int32, full_fight_amount int32) pbutil.Buffer {
	msg := &S2CCaptainProgressProto{
		Id:              id,
		RebirthLevel:    rebirth_level,
		RebirthExp:      rebirth_exp,
		Quality:         quality,
		Ability:         ability,
		AbilityExp:      ability_exp,
		AbilityLimit:    ability_limit,
		Soldier:         soldier,
		SoldierCapcity:  soldier_capcity,
		TotalStat:       safeMarshal(total_stat),
		FightAmount:     fight_amount,
		FullFightAmount: full_fight_amount,
	}
	return NewS2cCaptainProgressProtoMsg(msg)
}

var s2c_captain_progress = [...]byte{4, 211, 1} // 211
func NewS2cCaptainProgressProtoMsg(object *S2CCaptainProgressProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_progress[:], "s2c_captain_progress")

}

// 无效的武将id
var ERR_CAPTAIN_PROGRESS_FAIL_INVALID_ID = pbutil.StaticBuffer{4, 4, 212, 1, 1} // 212-1

// 未达到转职等级
var ERR_CAPTAIN_PROGRESS_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 212, 1, 2} // 212-2

// 已转生到最高等级
var ERR_CAPTAIN_PROGRESS_FAIL_MAX_LEVEL = pbutil.StaticBuffer{4, 4, 212, 1, 3} // 212-3

// 服务器忙，请稍后再试
var ERR_CAPTAIN_PROGRESS_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 4, 212, 1, 4} // 212-4

// 还在转生CD中
var ERR_CAPTAIN_PROGRESS_FAIL_IN_REBIRTHING = pbutil.StaticBuffer{4, 4, 212, 1, 5} // 212-5

// 秒CD消耗不足
var ERR_CAPTAIN_PROGRESS_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 212, 1, 6} // 212-6

func NewS2cCaptainRebirthMiaoCdMsg(id int32) pbutil.Buffer {
	msg := &S2CCaptainRebirthMiaoCdProto{
		Id: id,
	}
	return NewS2cCaptainRebirthMiaoCdProtoMsg(msg)
}

var s2c_captain_rebirth_miao_cd = [...]byte{4, 167, 1} // 167
func NewS2cCaptainRebirthMiaoCdProtoMsg(object *S2CCaptainRebirthMiaoCdProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_rebirth_miao_cd[:], "s2c_captain_rebirth_miao_cd")

}

// 无效的武将id
var ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_INVALID_ID = pbutil.StaticBuffer{4, 4, 168, 1, 1} // 168-1

// 没有在转生，或CD已结束
var ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_NOT_IN_REBIRTH = pbutil.StaticBuffer{4, 4, 168, 1, 2} // 168-2

// 点券不够
var ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 168, 1, 5} // 168-5

// 服务器忙，请稍后再试
var ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 4, 168, 1, 4} // 168-4

func NewS2cCollectCaptainTrainingExpMsg(time int32) pbutil.Buffer {
	msg := &S2CCollectCaptainTrainingExpProto{
		Time: time,
	}
	return NewS2cCollectCaptainTrainingExpProtoMsg(msg)
}

var s2c_collect_captain_training_exp = [...]byte{4, 137, 1} // 137
func NewS2cCollectCaptainTrainingExpProtoMsg(object *S2CCollectCaptainTrainingExpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_captain_training_exp[:], "s2c_collect_captain_training_exp")

}

// 修炼馆未解锁
var ERR_COLLECT_CAPTAIN_TRAINING_EXP_FAIL_BUILDING_LOCKED = pbutil.StaticBuffer{4, 4, 138, 1, 1} // 138-1

func NewS2cCaptainTrainExpMsg(time int32) pbutil.Buffer {
	msg := &S2CCaptainTrainExpProto{
		Time: time,
	}
	return NewS2cCaptainTrainExpProtoMsg(msg)
}

var s2c_captain_train_exp = [...]byte{4, 214, 1} // 214
func NewS2cCaptainTrainExpProtoMsg(object *S2CCaptainTrainExpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_train_exp[:], "s2c_captain_train_exp")

}

// 修炼馆未解锁
var ERR_CAPTAIN_TRAIN_EXP_FAIL_BUILDING_LOCKED = pbutil.StaticBuffer{4, 4, 215, 1, 1} // 215-1

func NewS2cUpdateTrainingMsg(gst int32, cst int32, exp_per_hour int32, coef int32) pbutil.Buffer {
	msg := &S2CUpdateTrainingProto{
		Gst:        gst,
		Cst:        cst,
		ExpPerHour: exp_per_hour,
		Coef:       coef,
	}
	return NewS2cUpdateTrainingProtoMsg(msg)
}

var s2c_update_training = [...]byte{4, 139, 1} // 139
func NewS2cUpdateTrainingProtoMsg(object *S2CUpdateTrainingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_training[:], "s2c_update_training")

}

func NewS2cCaptainCanCollectExpMsg(exp int32, buff_coef []int32, max_duration int32) pbutil.Buffer {
	msg := &S2CCaptainCanCollectExpProto{
		Exp:         exp,
		BuffCoef:    buff_coef,
		MaxDuration: max_duration,
	}
	return NewS2cCaptainCanCollectExpProtoMsg(msg)
}

var s2c_captain_can_collect_exp = [...]byte{4, 134, 2} // 262
func NewS2cCaptainCanCollectExpProtoMsg(object *S2CCaptainCanCollectExpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_can_collect_exp[:], "s2c_captain_can_collect_exp")

}

// 修炼馆未解锁
var ERR_CAPTAIN_CAN_COLLECT_EXP_FAIL_UNLOCK = pbutil.StaticBuffer{4, 4, 135, 2, 1} // 263-1

func NewS2cUseTrainingExpGoodsMsg(captain_id int32, goods_id int32, count int32, upgrade bool) pbutil.Buffer {
	msg := &S2CUseTrainingExpGoodsProto{
		CaptainId: captain_id,
		GoodsId:   goods_id,
		Count:     count,
		Upgrade:   upgrade,
	}
	return NewS2cUseTrainingExpGoodsProtoMsg(msg)
}

var s2c_use_training_exp_goods = [...]byte{4, 141, 1} // 141
func NewS2cUseTrainingExpGoodsProtoMsg(object *S2CUseTrainingExpGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_training_exp_goods[:], "s2c_use_training_exp_goods")

}

// 修炼馆未解锁
var ERR_USE_TRAINING_EXP_GOODS_FAIL_BUILDING_LOCKED = pbutil.StaticBuffer{4, 4, 142, 1, 3} // 142-3

// 无效的武将id
var ERR_USE_TRAINING_EXP_GOODS_FAIL_INVALID_CAPTAIN = pbutil.StaticBuffer{4, 4, 142, 1, 1} // 142-1

// 无效的物品id
var ERR_USE_TRAINING_EXP_GOODS_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 4, 142, 1, 2} // 142-2

// 物品个数不足
var ERR_USE_TRAINING_EXP_GOODS_FAIL_GOODS_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 142, 1, 4} // 142-4

// 武将已达最大等级
var ERR_USE_TRAINING_EXP_GOODS_FAIL_CAPTAIN_MAX_LEVEL = pbutil.StaticBuffer{4, 4, 142, 1, 5} // 142-5

// 武将等级受限，请提升君主等级
var ERR_USE_TRAINING_EXP_GOODS_FAIL_CAPTAIN_LEVEL_LIMIT = pbutil.StaticBuffer{4, 4, 142, 1, 6} // 142-6

// 武将在转生 CD 中
var ERR_USE_TRAINING_EXP_GOODS_FAIL_IN_REBIRTHING = pbutil.StaticBuffer{4, 4, 142, 1, 7} // 142-7

func NewS2cUseLevelExpGoodsMsg(captain_id int32, goods_id int32, count int32, upgrade bool) pbutil.Buffer {
	msg := &S2CUseLevelExpGoodsProto{
		CaptainId: captain_id,
		GoodsId:   goods_id,
		Count:     count,
		Upgrade:   upgrade,
	}
	return NewS2cUseLevelExpGoodsProtoMsg(msg)
}

var s2c_use_level_exp_goods = [...]byte{4, 217, 1} // 217
func NewS2cUseLevelExpGoodsProtoMsg(object *S2CUseLevelExpGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_level_exp_goods[:], "s2c_use_level_exp_goods")

}

// 修炼馆未解锁
var ERR_USE_LEVEL_EXP_GOODS_FAIL_BUILDING_LOCKED = pbutil.StaticBuffer{4, 4, 218, 1, 1} // 218-1

// 无效的武将id
var ERR_USE_LEVEL_EXP_GOODS_FAIL_INVALID_CAPTAIN = pbutil.StaticBuffer{4, 4, 218, 1, 2} // 218-2

// 无效的物品id
var ERR_USE_LEVEL_EXP_GOODS_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 4, 218, 1, 3} // 218-3

// 物品个数不足
var ERR_USE_LEVEL_EXP_GOODS_FAIL_GOODS_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 218, 1, 4} // 218-4

// 武将已达最大等级
var ERR_USE_LEVEL_EXP_GOODS_FAIL_CAPTAIN_MAX_LEVEL = pbutil.StaticBuffer{4, 4, 218, 1, 5} // 218-5

// 武将等级受限，请提升君主等级
var ERR_USE_LEVEL_EXP_GOODS_FAIL_CAPTAIN_LEVEL_LIMIT = pbutil.StaticBuffer{4, 4, 218, 1, 6} // 218-6

// 武将在转生 CD 中
var ERR_USE_LEVEL_EXP_GOODS_FAIL_IN_REBIRTHING = pbutil.StaticBuffer{4, 4, 218, 1, 7} // 218-7

func NewS2cUseLevelExpGoods2Msg(captain int32, level int32, exp int32, upgrade bool) pbutil.Buffer {
	msg := &S2CUseLevelExpGoods2Proto{
		Captain: captain,
		Level:   level,
		Exp:     exp,
		Upgrade: upgrade,
	}
	return NewS2cUseLevelExpGoods2ProtoMsg(msg)
}

var s2c_use_level_exp_goods2 = [...]byte{4, 244, 1} // 244
func NewS2cUseLevelExpGoods2ProtoMsg(object *S2CUseLevelExpGoods2Proto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_level_exp_goods2[:], "s2c_use_level_exp_goods2")

}

// 修炼馆未解锁
var ERR_USE_LEVEL_EXP_GOODS2_FAIL_BUILDING_LOCKED = pbutil.StaticBuffer{4, 4, 245, 1, 1} // 245-1

// 没有该武将
var ERR_USE_LEVEL_EXP_GOODS2_FAIL_NO_CAPTAIN = pbutil.StaticBuffer{4, 4, 245, 1, 8} // 245-8

// 不是经验书
var ERR_USE_LEVEL_EXP_GOODS2_FAIL_NOT_EXP_GOODS = pbutil.StaticBuffer{4, 4, 245, 1, 9} // 245-9

// 经验书不足
var ERR_USE_LEVEL_EXP_GOODS2_FAIL_NOT_ENOUTH_GOODS = pbutil.StaticBuffer{4, 4, 245, 1, 10} // 245-10

// 武将已达最大等级
var ERR_USE_LEVEL_EXP_GOODS2_FAIL_CAPTAIN_MAX_LEVEL = pbutil.StaticBuffer{4, 4, 245, 1, 5} // 245-5

// 武将等级受限，请提升君主等级
var ERR_USE_LEVEL_EXP_GOODS2_FAIL_CAPTAIN_LEVEL_LIMIT = pbutil.StaticBuffer{4, 4, 245, 1, 6} // 245-6

// 武将在转生 CD 中
var ERR_USE_LEVEL_EXP_GOODS2_FAIL_IN_REBIRTHING = pbutil.StaticBuffer{4, 4, 245, 1, 7} // 245-7

func NewS2cAutoUseGoodsUntilCaptainLevelupMsg(captain int32, level int32, exp int32, upgrade bool) pbutil.Buffer {
	msg := &S2CAutoUseGoodsUntilCaptainLevelupProto{
		Captain: captain,
		Level:   level,
		Exp:     exp,
		Upgrade: upgrade,
	}
	return NewS2cAutoUseGoodsUntilCaptainLevelupProtoMsg(msg)
}

var s2c_auto_use_goods_until_captain_levelup = [...]byte{4, 128, 2} // 256
func NewS2cAutoUseGoodsUntilCaptainLevelupProtoMsg(object *S2CAutoUseGoodsUntilCaptainLevelupProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_auto_use_goods_until_captain_levelup[:], "s2c_auto_use_goods_until_captain_levelup")

}

// 修炼馆未解锁
var ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_BUILDING_LOCKED = pbutil.StaticBuffer{4, 4, 129, 2, 1} // 257-1

// 没有该武将
var ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_NO_CAPTAIN = pbutil.StaticBuffer{4, 4, 129, 2, 2} // 257-2

// 没有任何经验书
var ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_NO_EXP_GOODS = pbutil.StaticBuffer{4, 4, 129, 2, 3} // 257-3

// 武将已达最大等级
var ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_CAPTAIN_MAX_LEVEL = pbutil.StaticBuffer{4, 4, 129, 2, 4} // 257-4

// 武将等级受限，请提升君主等级
var ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_CAPTAIN_LEVEL_LIMIT = pbutil.StaticBuffer{4, 4, 129, 2, 5} // 257-5

// 武将在转生 CD 中
var ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_IN_REBIRTHING = pbutil.StaticBuffer{4, 4, 129, 2, 6} // 257-6

func NewS2cGetMaxRecruitSoldierMsg(count int32) pbutil.Buffer {
	msg := &S2CGetMaxRecruitSoldierProto{
		Count: count,
	}
	return NewS2cGetMaxRecruitSoldierProtoMsg(msg)
}

var s2c_get_max_recruit_soldier = [...]byte{4, 75} // 75
func NewS2cGetMaxRecruitSoldierProtoMsg(object *S2CGetMaxRecruitSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_max_recruit_soldier[:], "s2c_get_max_recruit_soldier")

}

func NewS2cGetMaxHealSoldierMsg(count int32) pbutil.Buffer {
	msg := &S2CGetMaxHealSoldierProto{
		Count: count,
	}
	return NewS2cGetMaxHealSoldierProtoMsg(msg)
}

var s2c_get_max_heal_soldier = [...]byte{4, 77} // 77
func NewS2cGetMaxHealSoldierProtoMsg(object *S2CGetMaxHealSoldierProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_max_heal_soldier[:], "s2c_get_max_heal_soldier")

}

func NewS2cJiuGuanConsultMsg(prize []byte, crit_multi int32, crit_multi_img_index int32, original_index int32, tutor_index int32) pbutil.Buffer {
	msg := &S2CJiuGuanConsultProto{
		Prize:             prize,
		CritMulti:         crit_multi,
		CritMultiImgIndex: crit_multi_img_index,
		OriginalIndex:     original_index,
		TutorIndex:        tutor_index,
	}
	return NewS2cJiuGuanConsultProtoMsg(msg)
}

var s2c_jiu_guan_consult = [...]byte{4, 113} // 113
func NewS2cJiuGuanConsultProtoMsg(object *S2CJiuGuanConsultProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_jiu_guan_consult[:], "s2c_jiu_guan_consult")

}

// 没有次数了
var ERR_JIU_GUAN_CONSULT_FAIL_NO_TIMES = pbutil.StaticBuffer{3, 4, 114, 1} // 114-1

// 下次请教倒计时未结束
var ERR_JIU_GUAN_CONSULT_FAIL_COUNTDOWN = pbutil.StaticBuffer{3, 4, 114, 4} // 114-4

// 没有酒馆
var ERR_JIU_GUAN_CONSULT_FAIL_NO_JIU_GUAN = pbutil.StaticBuffer{3, 4, 114, 2} // 114-2

// 服务器忙，请稍后再试
var ERR_JIU_GUAN_CONSULT_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 4, 114, 3} // 114-3

func NewS2cJiuGuanConsultBroadcastMsg(level int32, crit_multi int32, hero_name string) pbutil.Buffer {
	msg := &S2CJiuGuanConsultBroadcastProto{
		Level:     level,
		CritMulti: crit_multi,
		HeroName:  hero_name,
	}
	return NewS2cJiuGuanConsultBroadcastProtoMsg(msg)
}

var s2c_jiu_guan_consult_broadcast = [...]byte{4, 115} // 115
func NewS2cJiuGuanConsultBroadcastProtoMsg(object *S2CJiuGuanConsultBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_jiu_guan_consult_broadcast[:], "s2c_jiu_guan_consult_broadcast")

}

func NewS2cJiuGuanTimesChangedMsg(times int32, next_time int32) pbutil.Buffer {
	msg := &S2CJiuGuanTimesChangedProto{
		Times:    times,
		NextTime: next_time,
	}
	return NewS2cJiuGuanTimesChangedProtoMsg(msg)
}

var s2c_jiu_guan_times_changed = [...]byte{4, 116} // 116
func NewS2cJiuGuanTimesChangedProtoMsg(object *S2CJiuGuanTimesChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_jiu_guan_times_changed[:], "s2c_jiu_guan_times_changed")

}

func NewS2cJiuGuanRefreshMsg(tutor_index int32, refresh_times int32, auto_max bool) pbutil.Buffer {
	msg := &S2CJiuGuanRefreshProto{
		TutorIndex:   tutor_index,
		RefreshTimes: refresh_times,
		AutoMax:      auto_max,
	}
	return NewS2cJiuGuanRefreshProtoMsg(msg)
}

var s2c_jiu_guan_refresh = [...]byte{4, 118} // 118
func NewS2cJiuGuanRefreshProtoMsg(object *S2CJiuGuanRefreshProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_jiu_guan_refresh[:], "s2c_jiu_guan_refresh")

}

// 没有请教次数了, 无法刷新
var ERR_JIU_GUAN_REFRESH_FAIL_NO_TIMES = pbutil.StaticBuffer{3, 4, 119, 1} // 119-1

// 下次请教倒计时未结束
var ERR_JIU_GUAN_REFRESH_FAIL_COUNTDOWN = pbutil.StaticBuffer{3, 4, 119, 6} // 119-6

// 没有酒馆, 无法刷新
var ERR_JIU_GUAN_REFRESH_FAIL_NO_JIU_GUAN = pbutil.StaticBuffer{3, 4, 119, 2} // 119-2

// 已经是最好的导师了
var ERR_JIU_GUAN_REFRESH_FAIL_TUTOR_NO_NEED = pbutil.StaticBuffer{3, 4, 119, 3} // 119-3

// 点券不够
var ERR_JIU_GUAN_REFRESH_FAIL_YUANBAO_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 119, 4} // 119-4

// 未开放一键刷新
var ERR_JIU_GUAN_REFRESH_FAIL_AUTO_MAX_NOT_OPEN = pbutil.StaticBuffer{3, 4, 119, 7} // 119-7

// 服务器忙，请稍后再试
var ERR_JIU_GUAN_REFRESH_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 4, 119, 5} // 119-5

// 没有刷新次数了, 无法刷新
var ERR_JIU_GUAN_REFRESH_FAIL_NO_REFRESH_TIMES = pbutil.StaticBuffer{3, 4, 119, 9} // 119-9

func NewS2cUnlockCaptainRestraintSpellMsg(captain int32) pbutil.Buffer {
	msg := &S2CUnlockCaptainRestraintSpellProto{
		Captain: captain,
	}
	return NewS2cUnlockCaptainRestraintSpellProtoMsg(msg)
}

var s2c_unlock_captain_restraint_spell = [...]byte{4, 126} // 126
func NewS2cUnlockCaptainRestraintSpellProtoMsg(object *S2CUnlockCaptainRestraintSpellProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_unlock_captain_restraint_spell[:], "s2c_unlock_captain_restraint_spell")

}

// 武将没找到
var ERR_UNLOCK_CAPTAIN_RESTRAINT_SPELL_FAIL_CAPTAIN_NOT_FOUND = pbutil.StaticBuffer{3, 4, 127, 1} // 127-1

// 武将出征中
var ERR_UNLOCK_CAPTAIN_RESTRAINT_SPELL_FAIL_OUT_SIDE = pbutil.StaticBuffer{3, 4, 127, 2} // 127-2

// 成长不够，无法解锁
var ERR_UNLOCK_CAPTAIN_RESTRAINT_SPELL_FAIL_ABILITY_NOT_ENOUGH = pbutil.StaticBuffer{3, 4, 127, 3} // 127-3

func NewS2cNewTroopsMsg(troop []byte) pbutil.Buffer {
	msg := &S2CNewTroopsProto{
		Troop: troop,
	}
	return NewS2cNewTroopsProtoMsg(msg)
}

var s2c_new_troops = [...]byte{4, 131, 1} // 131
func NewS2cNewTroopsProtoMsg(object *S2CNewTroopsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_new_troops[:], "s2c_new_troops")

}

func NewS2cGetCaptainStatDetailsMsg(captain int32, stats []byte) pbutil.Buffer {
	msg := &S2CGetCaptainStatDetailsProto{
		Captain: captain,
		Stats:   stats,
	}
	return NewS2cGetCaptainStatDetailsProtoMsg(msg)
}

func NewS2cGetCaptainStatDetailsMarshalMsg(captain int32, stats marshaler) pbutil.Buffer {
	msg := &S2CGetCaptainStatDetailsProto{
		Captain: captain,
		Stats:   safeMarshal(stats),
	}
	return NewS2cGetCaptainStatDetailsProtoMsg(msg)
}

var s2c_get_captain_stat_details = [...]byte{4, 133, 1} // 133
func NewS2cGetCaptainStatDetailsProtoMsg(object *S2CGetCaptainStatDetailsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_captain_stat_details[:], "s2c_get_captain_stat_details")

}

// 武将没找到
var ERR_GET_CAPTAIN_STAT_DETAILS_FAIL_CAPTAIN_NOT_FOUND = pbutil.StaticBuffer{4, 4, 134, 1, 1} // 134-1

func NewS2cCaptainStatDetailsMsg(captain int32, stats []byte) pbutil.Buffer {
	msg := &S2CCaptainStatDetailsProto{
		Captain: captain,
		Stats:   stats,
	}
	return NewS2cCaptainStatDetailsProtoMsg(msg)
}

func NewS2cCaptainStatDetailsMarshalMsg(captain int32, stats marshaler) pbutil.Buffer {
	msg := &S2CCaptainStatDetailsProto{
		Captain: captain,
		Stats:   safeMarshal(stats),
	}
	return NewS2cCaptainStatDetailsProtoMsg(msg)
}

var s2c_captain_stat_details = [...]byte{4, 220, 1} // 220
func NewS2cCaptainStatDetailsProtoMsg(object *S2CCaptainStatDetailsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_stat_details[:], "s2c_captain_stat_details")

}

// 武将没找到
var ERR_CAPTAIN_STAT_DETAILS_FAIL_CAPTAIN_NOT_FOUND = pbutil.StaticBuffer{4, 4, 221, 1, 1} // 221-1

func NewS2cUpdateTroopFightAmountMsg(troop_index []int32, start_value []int32, end_value []int32) pbutil.Buffer {
	msg := &S2CUpdateTroopFightAmountProto{
		TroopIndex: troop_index,
		StartValue: start_value,
		EndValue:   end_value,
	}
	return NewS2cUpdateTroopFightAmountProtoMsg(msg)
}

var s2c_update_troop_fight_amount = [...]byte{4, 135, 1} // 135
func NewS2cUpdateTroopFightAmountProtoMsg(object *S2CUpdateTroopFightAmountProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_troop_fight_amount[:], "s2c_update_troop_fight_amount")

}

func NewS2cUpdateCaptainOfficialMsg(captain int32, official int32) pbutil.Buffer {
	msg := &S2CUpdateCaptainOfficialProto{
		Captain:  captain,
		Official: official,
	}
	return NewS2cUpdateCaptainOfficialProtoMsg(msg)
}

var s2c_update_captain_official = [...]byte{4, 170, 1} // 170
func NewS2cUpdateCaptainOfficialProtoMsg(object *S2CUpdateCaptainOfficialProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_captain_official[:], "s2c_update_captain_official")

}

// 武将不存在
var ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_INVALID_CAPTAIN = pbutil.StaticBuffer{4, 4, 171, 1, 6} // 171-6

// 官职不存在
var ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_INVALID_OFFICIAL = pbutil.StaticBuffer{4, 4, 171, 1, 7} // 171-7

// 已经在这个职位上了
var ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_ALREADY_ON_OFFICIAL = pbutil.StaticBuffer{4, 4, 171, 1, 9} // 171-9

// 武将功勋不够
var ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_GONGXUN_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 171, 1, 3} // 171-3

// 职位册封人数已达到最大数量
var ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_MAX_COUNT = pbutil.StaticBuffer{4, 4, 171, 1, 4} // 171-4

// 降职时，武将在外面
var ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_CAPTAIN_IS_OUTSIDE = pbutil.StaticBuffer{4, 4, 171, 1, 5} // 171-5

func NewS2cSetCaptainOfficialMsg(captain []int32, official []int32, official_idx []int32) pbutil.Buffer {
	msg := &S2CSetCaptainOfficialProto{
		Captain:     captain,
		Official:    official,
		OfficialIdx: official_idx,
	}
	return NewS2cSetCaptainOfficialProtoMsg(msg)
}

var s2c_set_captain_official = [...]byte{4, 223, 1} // 223
func NewS2cSetCaptainOfficialProtoMsg(object *S2CSetCaptainOfficialProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_captain_official[:], "s2c_set_captain_official")

}

// 数据错误
var ERR_SET_CAPTAIN_OFFICIAL_FAIL_ERR_DATA = pbutil.StaticBuffer{4, 4, 224, 1, 8} // 224-8

// 武将不存在
var ERR_SET_CAPTAIN_OFFICIAL_FAIL_INVALID_CAPTAIN = pbutil.StaticBuffer{4, 4, 224, 1, 1} // 224-1

// 官职不存在
var ERR_SET_CAPTAIN_OFFICIAL_FAIL_INVALID_OFFICIAL = pbutil.StaticBuffer{4, 4, 224, 1, 2} // 224-2

// 已经在这个职位上并且位置相同
var ERR_SET_CAPTAIN_OFFICIAL_FAIL_ALREADY_ON_OFFICIAL = pbutil.StaticBuffer{4, 4, 224, 1, 3} // 224-3

// 武将功勋不够
var ERR_SET_CAPTAIN_OFFICIAL_FAIL_GONGXUN_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 224, 1, 4} // 224-4

// 武将在外面
var ERR_SET_CAPTAIN_OFFICIAL_FAIL_CAPTAIN_IS_OUTSIDE = pbutil.StaticBuffer{4, 4, 224, 1, 6} // 224-6

// 已经没有职位了
var ERR_SET_CAPTAIN_OFFICIAL_FAIL_ALREADY_NO_OFFICIAL = pbutil.StaticBuffer{4, 4, 224, 1, 7} // 224-7

func NewS2cLeaveCaptainOfficialMsg(captain int32) pbutil.Buffer {
	msg := &S2CLeaveCaptainOfficialProto{
		Captain: captain,
	}
	return NewS2cLeaveCaptainOfficialProtoMsg(msg)
}

var s2c_leave_captain_official = [...]byte{4, 173, 1} // 173
func NewS2cLeaveCaptainOfficialProtoMsg(object *S2CLeaveCaptainOfficialProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_leave_captain_official[:], "s2c_leave_captain_official")

}

// 武将不存在
var ERR_LEAVE_CAPTAIN_OFFICIAL_FAIL_INVALID_CAPTAIN = pbutil.StaticBuffer{4, 4, 174, 1, 1} // 174-1

// 已经没有职位了
var ERR_LEAVE_CAPTAIN_OFFICIAL_FAIL_ALREADY_NO_OFFICIAL = pbutil.StaticBuffer{4, 4, 174, 1, 2} // 174-2

// 卸任时，武将在外面
var ERR_LEAVE_CAPTAIN_OFFICIAL_FAIL_CAPTAIN_IS_OUTSIDE = pbutil.StaticBuffer{4, 4, 174, 1, 3} // 174-3

func NewS2cAddGongxunMsg(captain int32, new_gongxun int32) pbutil.Buffer {
	msg := &S2CAddGongxunProto{
		Captain:    captain,
		NewGongxun: new_gongxun,
	}
	return NewS2cAddGongxunProtoMsg(msg)
}

var s2c_add_gongxun = [...]byte{4, 175, 1} // 175
func NewS2cAddGongxunProtoMsg(object *S2CAddGongxunProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_gongxun[:], "s2c_add_gongxun")

}

func NewS2cUseGongXunGoodsMsg(captain int32, new_gongxun int32) pbutil.Buffer {
	msg := &S2CUseGongXunGoodsProto{
		Captain:    captain,
		NewGongxun: new_gongxun,
	}
	return NewS2cUseGongXunGoodsProtoMsg(msg)
}

var s2c_use_gong_xun_goods = [...]byte{4, 186, 1} // 186
func NewS2cUseGongXunGoodsProtoMsg(object *S2CUseGongXunGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_gong_xun_goods[:], "s2c_use_gong_xun_goods")

}

// 无效的武将id
var ERR_USE_GONG_XUN_GOODS_FAIL_INVALID_CAPTAIN = pbutil.StaticBuffer{4, 4, 187, 1, 1} // 187-1

// 功勋已达最高封官
var ERR_USE_GONG_XUN_GOODS_FAIL_OFFICIAL_LIMIT = pbutil.StaticBuffer{4, 4, 187, 1, 2} // 187-2

// 无效的功勋令牌id
var ERR_USE_GONG_XUN_GOODS_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 4, 187, 1, 3} // 187-3

// 无效的使用个数
var ERR_USE_GONG_XUN_GOODS_FAIL_INVALID_COUNT = pbutil.StaticBuffer{4, 4, 187, 1, 4} // 187-4

func NewS2cUseGongxunGoodsMsg(captain int32, new_gongxun int32) pbutil.Buffer {
	msg := &S2CUseGongxunGoodsProto{
		Captain:    captain,
		NewGongxun: new_gongxun,
	}
	return NewS2cUseGongxunGoodsProtoMsg(msg)
}

var s2c_use_gongxun_goods = [...]byte{4, 229, 1} // 229
func NewS2cUseGongxunGoodsProtoMsg(object *S2CUseGongxunGoodsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_gongxun_goods[:], "s2c_use_gongxun_goods")

}

// 无效的武将id
var ERR_USE_GONGXUN_GOODS_FAIL_INVALID_CAPTAIN = pbutil.StaticBuffer{4, 4, 230, 1, 1} // 230-1

// 功勋已达最高封官
var ERR_USE_GONGXUN_GOODS_FAIL_OFFICIAL_LIMIT = pbutil.StaticBuffer{4, 4, 230, 1, 2} // 230-2

// 无效的功勋令牌id
var ERR_USE_GONGXUN_GOODS_FAIL_INVALID_GOODS = pbutil.StaticBuffer{4, 4, 230, 1, 3} // 230-3

// 无效的使用个数
var ERR_USE_GONGXUN_GOODS_FAIL_INVALID_COUNT = pbutil.StaticBuffer{4, 4, 230, 1, 4} // 230-4

func NewS2cCloseFightGuideMsg(close bool) pbutil.Buffer {
	msg := &S2CCloseFightGuideProto{
		Close: close,
	}
	return NewS2cCloseFightGuideProtoMsg(msg)
}

var s2c_close_fight_guide = [...]byte{4, 182, 1} // 182
func NewS2cCloseFightGuideProtoMsg(object *S2CCloseFightGuideProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_close_fight_guide[:], "s2c_close_fight_guide")

}

// 没有权限，还没通关指定幻境
var ERR_CLOSE_FIGHT_GUIDE_FAIL_NO_AUTH = pbutil.StaticBuffer{4, 4, 183, 1, 1} // 183-1

func NewS2cViewOtherHeroCaptainMsg(hero_id []byte, hero_name string, captain []byte) pbutil.Buffer {
	msg := &S2CViewOtherHeroCaptainProto{
		HeroId:   hero_id,
		HeroName: hero_name,
		Captain:  captain,
	}
	return NewS2cViewOtherHeroCaptainProtoMsg(msg)
}

var s2c_view_other_hero_captain = [...]byte{4, 191, 1} // 191
func NewS2cViewOtherHeroCaptainProtoMsg(object *S2CViewOtherHeroCaptainProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_other_hero_captain[:], "s2c_view_other_hero_captain")

}

// 无效的英雄id
var ERR_VIEW_OTHER_HERO_CAPTAIN_FAIL_INVALID_HERO_ID = pbutil.StaticBuffer{4, 4, 192, 1, 1} // 192-1

// 武将不存在
var ERR_VIEW_OTHER_HERO_CAPTAIN_FAIL_NOT_FOUND = pbutil.StaticBuffer{4, 4, 192, 1, 2} // 192-2

func NewS2cCaptainBornMsg(captain []byte) pbutil.Buffer {
	msg := &S2CCaptainBornProto{
		Captain: captain,
	}
	return NewS2cCaptainBornProtoMsg(msg)
}

func NewS2cCaptainBornMarshalMsg(captain marshaler) pbutil.Buffer {
	msg := &S2CCaptainBornProto{
		Captain: safeMarshal(captain),
	}
	return NewS2cCaptainBornProtoMsg(msg)
}

var s2c_captain_born = [...]byte{4, 232, 1} // 232
func NewS2cCaptainBornProtoMsg(object *S2CCaptainBornProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_born[:], "s2c_captain_born")

}

// 没有该武将
var ERR_CAPTAIN_BORN_FAIL_NO_CAPTAIN = pbutil.StaticBuffer{4, 4, 233, 1, 1} // 233-1

// 已经获得该武将
var ERR_CAPTAIN_BORN_FAIL_EXISTED = pbutil.StaticBuffer{4, 4, 233, 1, 2} // 233-2

// 道具数量不足
var ERR_CAPTAIN_BORN_FAIL_ITEM_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 233, 1, 3} // 233-3

func NewS2cCaptainUpstarMsg(captain_id int32, star int32) pbutil.Buffer {
	msg := &S2CCaptainUpstarProto{
		CaptainId: captain_id,
		Star:      star,
	}
	return NewS2cCaptainUpstarProtoMsg(msg)
}

var s2c_captain_upstar = [...]byte{4, 235, 1} // 235
func NewS2cCaptainUpstarProtoMsg(object *S2CCaptainUpstarProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_upstar[:], "s2c_captain_upstar")

}

// 没有该武将
var ERR_CAPTAIN_UPSTAR_FAIL_NO_CAPTAIN = pbutil.StaticBuffer{4, 4, 236, 1, 1} // 236-1

// 未获得该武将
var ERR_CAPTAIN_UPSTAR_FAIL_NOT_GAIN = pbutil.StaticBuffer{4, 4, 236, 1, 2} // 236-2

// 已经升至最高星
var ERR_CAPTAIN_UPSTAR_FAIL_MAX_STAR = pbutil.StaticBuffer{4, 4, 236, 1, 3} // 236-3

// 道具数量不足
var ERR_CAPTAIN_UPSTAR_FAIL_ITEM_NOT_ENOUGH = pbutil.StaticBuffer{4, 4, 236, 1, 4} // 236-4

func NewS2cCaptainExchangeMsg(cap1 []byte, cap2 []byte) pbutil.Buffer {
	msg := &S2CCaptainExchangeProto{
		Cap1: cap1,
		Cap2: cap2,
	}
	return NewS2cCaptainExchangeProtoMsg(msg)
}

func NewS2cCaptainExchangeMarshalMsg(cap1 marshaler, cap2 marshaler) pbutil.Buffer {
	msg := &S2CCaptainExchangeProto{
		Cap1: safeMarshal(cap1),
		Cap2: safeMarshal(cap2),
	}
	return NewS2cCaptainExchangeProtoMsg(msg)
}

var s2c_captain_exchange = [...]byte{4, 141, 2} // 269
func NewS2cCaptainExchangeProtoMsg(object *S2CCaptainExchangeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_captain_exchange[:], "s2c_captain_exchange")

}

// 有武将不存在或相同
var ERR_CAPTAIN_EXCHANGE_FAIL_INVALID_CAPTAIN = pbutil.StaticBuffer{4, 4, 142, 2, 1} // 270-1

// 有武将出征中
var ERR_CAPTAIN_EXCHANGE_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{4, 4, 142, 2, 2} // 270-2

// 点券不足
var ERR_CAPTAIN_EXCHANGE_FAIL_NOT_ENOUGH_COST = pbutil.StaticBuffer{4, 4, 142, 2, 3} // 270-3

var NOTICE_CAPTAIN_HAS_VIEWED_S2C = pbutil.StaticBuffer{3, 4, 253, 1} // 253

// 没有任何武将未被浏览过
var ERR_NOTICE_CAPTAIN_HAS_VIEWED_FAIL_NO_CAPTAIN_NO_VIEWED = pbutil.StaticBuffer{4, 4, 254, 1, 3} // 254-3

func NewS2cActivateCaptainFriendshipMsg(id int32) pbutil.Buffer {
	msg := &S2CActivateCaptainFriendshipProto{
		Id: id,
	}
	return NewS2cActivateCaptainFriendshipProtoMsg(msg)
}

var s2c_activate_captain_friendship = [...]byte{4, 138, 2} // 266
func NewS2cActivateCaptainFriendshipProtoMsg(object *S2CActivateCaptainFriendshipProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_activate_captain_friendship[:], "s2c_activate_captain_friendship")

}

// 无效的id
var ERR_ACTIVATE_CAPTAIN_FRIENDSHIP_FAIL_INVALID_ID = pbutil.StaticBuffer{4, 4, 139, 2, 1} // 267-1

// 羁绊条件不足
var ERR_ACTIVATE_CAPTAIN_FRIENDSHIP_FAIL_NO_FETTER = pbutil.StaticBuffer{4, 4, 139, 2, 2} // 267-2

// 已激活
var ERR_ACTIVATE_CAPTAIN_FRIENDSHIP_FAIL_ACTIVATED = pbutil.StaticBuffer{4, 4, 139, 2, 3} // 267-3

func NewS2cShowPrizeCaptainMsg(captain int32, exist bool) pbutil.Buffer {
	msg := &S2CShowPrizeCaptainProto{
		Captain: captain,
		Exist:   exist,
	}
	return NewS2cShowPrizeCaptainProtoMsg(msg)
}

var s2c_show_prize_captain = [...]byte{4, 143, 2} // 271
func NewS2cShowPrizeCaptainProtoMsg(object *S2CShowPrizeCaptainProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_show_prize_captain[:], "s2c_show_prize_captain")

}

func NewS2cNoticeOfficialHasViewedMsg(official_id int32, official_idx int32) pbutil.Buffer {
	msg := &S2CNoticeOfficialHasViewedProto{
		OfficialId:  official_id,
		OfficialIdx: official_idx,
	}
	return NewS2cNoticeOfficialHasViewedProtoMsg(msg)
}

var s2c_notice_official_has_viewed = [...]byte{4, 145, 2} // 273
func NewS2cNoticeOfficialHasViewedProtoMsg(object *S2CNoticeOfficialHasViewedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_notice_official_has_viewed[:], "s2c_notice_official_has_viewed")

}

// 官职不存在
var ERR_NOTICE_OFFICIAL_HAS_VIEWED_FAIL_NO_OFFICIAL = pbutil.StaticBuffer{4, 4, 146, 2, 1} // 274-1

// 已经设置过
var ERR_NOTICE_OFFICIAL_HAS_VIEWED_FAIL_VIEWED = pbutil.StaticBuffer{4, 4, 146, 2, 2} // 274-2
