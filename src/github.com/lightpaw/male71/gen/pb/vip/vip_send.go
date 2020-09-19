package vip

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
	MODULE_ID = 48

	C2S_VIP_COLLECT_DAILY_PRIZE = 4

	C2S_VIP_COLLECT_LEVEL_PRIZE = 7

	C2S_VIP_BUY_DUNGEON_TIMES = 12
)

func NewS2cVipLevelUpgradeNoticeMsg(vip_level int32, vip_exp int32) pbutil.Buffer {
	msg := &S2CVipLevelUpgradeNoticeProto{
		VipLevel: vip_level,
		VipExp:   vip_exp,
	}
	return NewS2cVipLevelUpgradeNoticeProtoMsg(msg)
}

var s2c_vip_level_upgrade_notice = [...]byte{48, 1} // 1
func NewS2cVipLevelUpgradeNoticeProtoMsg(object *S2CVipLevelUpgradeNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_vip_level_upgrade_notice[:], "s2c_vip_level_upgrade_notice")

}

func NewS2cVipAddExpNoticeMsg(vip_exp int32) pbutil.Buffer {
	msg := &S2CVipAddExpNoticeProto{
		VipExp: vip_exp,
	}
	return NewS2cVipAddExpNoticeProtoMsg(msg)
}

var s2c_vip_add_exp_notice = [...]byte{48, 2} // 2
func NewS2cVipAddExpNoticeProtoMsg(object *S2CVipAddExpNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_vip_add_exp_notice[:], "s2c_vip_add_exp_notice")

}

func NewS2cVipDailyLoginNoticeMsg(vip_level int32, vip_exp int32, continue_days int32, tomorrow_exp int32) pbutil.Buffer {
	msg := &S2CVipDailyLoginNoticeProto{
		VipLevel:     vip_level,
		VipExp:       vip_exp,
		ContinueDays: continue_days,
		TomorrowExp:  tomorrow_exp,
	}
	return NewS2cVipDailyLoginNoticeProtoMsg(msg)
}

var s2c_vip_daily_login_notice = [...]byte{48, 3} // 3
func NewS2cVipDailyLoginNoticeProtoMsg(object *S2CVipDailyLoginNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_vip_daily_login_notice[:], "s2c_vip_daily_login_notice")

}

func NewS2cVipCollectDailyPrizeMsg(vip_level int32) pbutil.Buffer {
	msg := &S2CVipCollectDailyPrizeProto{
		VipLevel: vip_level,
	}
	return NewS2cVipCollectDailyPrizeProtoMsg(msg)
}

var s2c_vip_collect_daily_prize = [...]byte{48, 5} // 5
func NewS2cVipCollectDailyPrizeProtoMsg(object *S2CVipCollectDailyPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_vip_collect_daily_prize[:], "s2c_vip_collect_daily_prize")

}

// vip_level不存在
var ERR_VIP_COLLECT_DAILY_PRIZE_FAIL_INVALID_LEVEL = pbutil.StaticBuffer{3, 48, 6, 2} // 6-2

// 等级不够
var ERR_VIP_COLLECT_DAILY_PRIZE_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 48, 6, 3} // 6-3

// 今天已经领过
var ERR_VIP_COLLECT_DAILY_PRIZE_FAIL_ALREADY_COLLECTED = pbutil.StaticBuffer{3, 48, 6, 1} // 6-1

// 没有每日奖励
var ERR_VIP_COLLECT_DAILY_PRIZE_FAIL_NO_DAILY_PRIZE = pbutil.StaticBuffer{3, 48, 6, 4} // 6-4

func NewS2cVipCollectLevelPrizeMsg(vip_level int32) pbutil.Buffer {
	msg := &S2CVipCollectLevelPrizeProto{
		VipLevel: vip_level,
	}
	return NewS2cVipCollectLevelPrizeProtoMsg(msg)
}

var s2c_vip_collect_level_prize = [...]byte{48, 8} // 8
func NewS2cVipCollectLevelPrizeProtoMsg(object *S2CVipCollectLevelPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_vip_collect_level_prize[:], "s2c_vip_collect_level_prize")

}

// vip_level不存在
var ERR_VIP_COLLECT_LEVEL_PRIZE_FAIL_INVALID_LEVEL = pbutil.StaticBuffer{3, 48, 9, 1} // 9-1

// 等级不够
var ERR_VIP_COLLECT_LEVEL_PRIZE_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 48, 9, 2} // 9-2

// 已经买过
var ERR_VIP_COLLECT_LEVEL_PRIZE_FAIL_ALREADY_COLLECTED = pbutil.StaticBuffer{3, 48, 9, 3} // 9-3

// 钱不够
var ERR_VIP_COLLECT_LEVEL_PRIZE_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 48, 9, 5} // 9-5

// 没有奖励
var ERR_VIP_COLLECT_LEVEL_PRIZE_FAIL_NO_LEVEL_PRIZE = pbutil.StaticBuffer{3, 48, 9, 6} // 9-6

func NewS2cVipBuyDungeonTimesMsg(dungeon_id int32, new_times int32) pbutil.Buffer {
	msg := &S2CVipBuyDungeonTimesProto{
		DungeonId: dungeon_id,
		NewTimes:  new_times,
	}
	return NewS2cVipBuyDungeonTimesProtoMsg(msg)
}

var s2c_vip_buy_dungeon_times = [...]byte{48, 13} // 13
func NewS2cVipBuyDungeonTimesProtoMsg(object *S2CVipBuyDungeonTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_vip_buy_dungeon_times[:], "s2c_vip_buy_dungeon_times")

}

// 钱不够
var ERR_VIP_BUY_DUNGEON_TIMES_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 48, 14, 1} // 14-1

// vip等级不够，买的次数已到最大
var ERR_VIP_BUY_DUNGEON_TIMES_FAIL_VIP_LEVEL_LIMIT = pbutil.StaticBuffer{3, 48, 14, 2} // 14-2

// 推图关卡 id 错误
var ERR_VIP_BUY_DUNGEON_TIMES_FAIL_DUNGEON_ID_INVALID = pbutil.StaticBuffer{3, 48, 14, 3} // 14-3
