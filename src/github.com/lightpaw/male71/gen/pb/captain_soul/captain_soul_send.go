package captain_soul

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
	MODULE_ID = 17

	C2S_COLLECT_FETTERS_PRIZE = 3

	C2S_ACTIVATE_FETTERS = 6

	C2S_FU_SHEN = 9

	C2S_UPGRADE = 12

	C2S_UPGRADE_V2 = 15

	C2S_UNLOCK_SPELL = 18

	C2S_MARK_ALL = 21

	C2S_MARK = 29

	C2S_REBORN_PREVIEW = 27

	C2S_REBORN = 24
)

func NewS2cSoulActivatedMsg(data []byte) pbutil.Buffer {
	msg := &S2CSoulActivatedProto{
		Data: data,
	}
	return NewS2cSoulActivatedProtoMsg(msg)
}

var s2c_soul_activated = [...]byte{17, 1} // 1
func NewS2cSoulActivatedProtoMsg(object *S2CSoulActivatedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_soul_activated[:], "s2c_soul_activated")

}

func NewS2cCollectFettersPrizeMsg(id int32) pbutil.Buffer {
	msg := &S2CCollectFettersPrizeProto{
		Id: id,
	}
	return NewS2cCollectFettersPrizeProtoMsg(msg)
}

var s2c_collect_fetters_prize = [...]byte{17, 4} // 4
func NewS2cCollectFettersPrizeProtoMsg(object *S2CCollectFettersPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_fetters_prize[:], "s2c_collect_fetters_prize")

}

// 已经领取了
var ERR_COLLECT_FETTERS_PRIZE_FAIL_HAVE_COLLECTED = pbutil.StaticBuffer{3, 17, 5, 3} // 5-3

// 没有激活
var ERR_COLLECT_FETTERS_PRIZE_FAIL_HAVE_NOT_ACTIVATED = pbutil.StaticBuffer{3, 17, 5, 2} // 5-2

// 羁绊没有找到
var ERR_COLLECT_FETTERS_PRIZE_FAIL_FETTERS_NOT_FOUND = pbutil.StaticBuffer{3, 17, 5, 4} // 5-4

// 羁绊没奖励
var ERR_COLLECT_FETTERS_PRIZE_FAIL_FETTERS_NO_PRIZE = pbutil.StaticBuffer{3, 17, 5, 6} // 5-6

// 服务器错误
var ERR_COLLECT_FETTERS_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 17, 5, 5} // 5-5

func NewS2cActivateFettersMsg(id int32) pbutil.Buffer {
	msg := &S2CActivateFettersProto{
		Id: id,
	}
	return NewS2cActivateFettersProtoMsg(msg)
}

var s2c_activate_fetters = [...]byte{17, 7} // 7
func NewS2cActivateFettersProtoMsg(object *S2CActivateFettersProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_activate_fetters[:], "s2c_activate_fetters")

}

// 该羁绊中有将魂没有获得
var ERR_ACTIVATE_FETTERS_FAIL_CAN_NOT_ACTIVATE = pbutil.StaticBuffer{3, 17, 8, 1} // 8-1

// 已经激活了
var ERR_ACTIVATE_FETTERS_FAIL_HAVE_ACTIVATED = pbutil.StaticBuffer{3, 17, 8, 2} // 8-2

// 羁绊没有找到
var ERR_ACTIVATE_FETTERS_FAIL_FETTERS_NOT_FOUND = pbutil.StaticBuffer{3, 17, 8, 3} // 8-3

func NewS2cFuShenMsg(captain int32, up_soul_id int32, replace_soul_captain int32) pbutil.Buffer {
	msg := &S2CFuShenProto{
		Captain:            captain,
		UpSoulId:           up_soul_id,
		ReplaceSoulCaptain: replace_soul_captain,
	}
	return NewS2cFuShenProtoMsg(msg)
}

var s2c_fu_shen = [...]byte{17, 10} // 10
func NewS2cFuShenProtoMsg(object *S2CFuShenProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_fu_shen[:], "s2c_fu_shen")

}

// 该将魂已经被使用了
var ERR_FU_SHEN_FAIL_HAS_BEEN_USE = pbutil.StaticBuffer{3, 17, 11, 1} // 11-1

// 将魂没找到
var ERR_FU_SHEN_FAIL_UNKNOW_SOUL_ID = pbutil.StaticBuffer{3, 17, 11, 3} // 11-3

// 武将没找到
var ERR_FU_SHEN_FAIL_CAPTAIN_NOT_FOUND = pbutil.StaticBuffer{3, 17, 11, 4} // 11-4

// 武将外出了
var ERR_FU_SHEN_FAIL_CAPTAIN_OUTSIDE = pbutil.StaticBuffer{3, 17, 11, 5} // 11-5

// 附身将魂穿在其他出征的武将身上，无法附身
var ERR_FU_SHEN_FAIL_CAPTAIN_SOUL_OUTSIDE = pbutil.StaticBuffer{3, 17, 11, 8} // 11-8

// 没有附身的将魂
var ERR_FU_SHEN_FAIL_NO_FU_SHEN_SOUL = pbutil.StaticBuffer{3, 17, 11, 6} // 11-6

// 当前已经有附身的将魂了
var ERR_FU_SHEN_FAIL_HAS_SOUL = pbutil.StaticBuffer{3, 17, 11, 7} // 11-7

// 要附身的武将在外面，必须用更高等级的将魂附身
var ERR_FU_SHEN_FAIL_CAPTAIN_SOUL_OUTSIDE_QUALITY_ERR = pbutil.StaticBuffer{3, 17, 11, 9} // 11-9

// 替换外出武将的将魂，必须继承
var ERR_FU_SHEN_FAIL_CAPTAIN_SOUL_OUTSIDE_MUST_INHERIT = pbutil.StaticBuffer{3, 17, 11, 10} // 11-10

func NewS2cUpgradeMsg(captain_soul_id int32, captain_soul_level int32) pbutil.Buffer {
	msg := &S2CUpgradeProto{
		CaptainSoulId:    captain_soul_id,
		CaptainSoulLevel: captain_soul_level,
	}
	return NewS2cUpgradeProtoMsg(msg)
}

var s2c_upgrade = [...]byte{17, 13} // 13
func NewS2cUpgradeProtoMsg(object *S2CUpgradeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_upgrade[:], "s2c_upgrade")

}

// 将魂没找到
var ERR_UPGRADE_FAIL_UNKNOW_SOUL_ID = pbutil.StaticBuffer{3, 17, 14, 1} // 14-1

// 将魂等级已满
var ERR_UPGRADE_FAIL_LEVEL_MAX = pbutil.StaticBuffer{3, 17, 14, 2} // 14-2

// 物品不够
var ERR_UPGRADE_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 17, 14, 3} // 14-3

// 武将出征中
var ERR_UPGRADE_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 17, 14, 4} // 14-4

func NewS2cUpgradeV2Msg(captain_soul int32, exp int32) pbutil.Buffer {
	msg := &S2CUpgradeV2Proto{
		CaptainSoul: captain_soul,
		Exp:         exp,
	}
	return NewS2cUpgradeV2ProtoMsg(msg)
}

var s2c_upgrade_v2 = [...]byte{17, 16} // 16
func NewS2cUpgradeV2ProtoMsg(object *S2CUpgradeV2Proto) pbutil.Buffer {

	return newProtoMsg(object, s2c_upgrade_v2[:], "s2c_upgrade_v2")

}

// 无效的将魂id
var ERR_UPGRADE_V2_FAIL_INVALID_CAPTAIN_SOUL = pbutil.StaticBuffer{3, 17, 17, 1} // 17-1

// 将魂已经满级了
var ERR_UPGRADE_V2_FAIL_LEVEL_MAX = pbutil.StaticBuffer{3, 17, 17, 2} // 17-2

// 无效的强化符id
var ERR_UPGRADE_V2_FAIL_INVALID_GOODS = pbutil.StaticBuffer{3, 17, 17, 3} // 17-3

// 无效的使用个数
var ERR_UPGRADE_V2_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 17, 17, 4} // 17-4

// 武将出征中
var ERR_UPGRADE_V2_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 17, 17, 5} // 17-5

func NewS2cLevelUpgradeMsg(captain_soul int32, captain_soul_level int32, stat []byte) pbutil.Buffer {
	msg := &S2CLevelUpgradeProto{
		CaptainSoul:      captain_soul,
		CaptainSoulLevel: captain_soul_level,
		Stat:             stat,
	}
	return NewS2cLevelUpgradeProtoMsg(msg)
}

var s2c_level_upgrade = [...]byte{17, 23} // 23
func NewS2cLevelUpgradeProtoMsg(object *S2CLevelUpgradeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_level_upgrade[:], "s2c_level_upgrade")

}

func NewS2cUnlockSpellMsg(captain_soul int32, index int32, stat []byte) pbutil.Buffer {
	msg := &S2CUnlockSpellProto{
		CaptainSoul: captain_soul,
		Index:       index,
		Stat:        stat,
	}
	return NewS2cUnlockSpellProtoMsg(msg)
}

var s2c_unlock_spell = [...]byte{17, 19} // 19
func NewS2cUnlockSpellProtoMsg(object *S2CUnlockSpellProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_unlock_spell[:], "s2c_unlock_spell")

}

// 无效的将魂id
var ERR_UNLOCK_SPELL_FAIL_INVALID_CAPTAIN_SOUL = pbutil.StaticBuffer{3, 17, 20, 1} // 20-1

// 技能未找到
var ERR_UNLOCK_SPELL_FAIL_INVALID_SPELL = pbutil.StaticBuffer{3, 17, 20, 2} // 20-2

// 将魂等级不够
var ERR_UNLOCK_SPELL_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 17, 20, 3} // 20-3

// 已经解锁了
var ERR_UNLOCK_SPELL_FAIL_UNLOCKED = pbutil.StaticBuffer{3, 17, 20, 4} // 20-4

// 东西不够，无法解锁
var ERR_UNLOCK_SPELL_FAIL_NOT_ENOUGH = pbutil.StaticBuffer{3, 17, 20, 5} // 20-5

// 武将出征中
var ERR_UNLOCK_SPELL_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 17, 20, 6} // 20-6

// 服务器繁忙，请稍后再试
var ERR_UNLOCK_SPELL_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 17, 20, 7} // 20-7

var MARK_ALL_S2C = pbutil.StaticBuffer{2, 17, 22} // 22

func NewS2cMarkMsg(captain_soul int32) pbutil.Buffer {
	msg := &S2CMarkProto{
		CaptainSoul: captain_soul,
	}
	return NewS2cMarkProtoMsg(msg)
}

var s2c_mark = [...]byte{17, 30} // 30
func NewS2cMarkProtoMsg(object *S2CMarkProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_mark[:], "s2c_mark")

}

func NewS2cRebornPreviewMsg(captain_soul int32, prize []byte) pbutil.Buffer {
	msg := &S2CRebornPreviewProto{
		CaptainSoul: captain_soul,
		Prize:       prize,
	}
	return NewS2cRebornPreviewProtoMsg(msg)
}

func NewS2cRebornPreviewMarshalMsg(captain_soul int32, prize marshaler) pbutil.Buffer {
	msg := &S2CRebornPreviewProto{
		CaptainSoul: captain_soul,
		Prize:       safeMarshal(prize),
	}
	return NewS2cRebornPreviewProtoMsg(msg)
}

var s2c_reborn_preview = [...]byte{17, 28} // 28
func NewS2cRebornPreviewProtoMsg(object *S2CRebornPreviewProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_reborn_preview[:], "s2c_reborn_preview")

}

func NewS2cRebornMsg(captain_soul int32) pbutil.Buffer {
	msg := &S2CRebornProto{
		CaptainSoul: captain_soul,
	}
	return NewS2cRebornProtoMsg(msg)
}

var s2c_reborn = [...]byte{17, 25} // 25
func NewS2cRebornProtoMsg(object *S2CRebornProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_reborn[:], "s2c_reborn")

}

// 无效的将魂id
var ERR_REBORN_FAIL_INVALID_CAPTAIN_SOUL = pbutil.StaticBuffer{3, 17, 26, 1} // 26-1

// 将魂未升级，不能重生
var ERR_REBORN_FAIL_EMPTY_EXP = pbutil.StaticBuffer{3, 17, 26, 2} // 26-2

// 武将出征中
var ERR_REBORN_FAIL_OUTSIDE = pbutil.StaticBuffer{3, 17, 26, 3} // 26-3
