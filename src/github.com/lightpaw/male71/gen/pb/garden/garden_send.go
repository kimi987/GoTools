package garden

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
	MODULE_ID = 31

	C2S_LIST_TREASURY_TREE_HERO = 1

	C2S_LIST_HELP_ME = 13

	C2S_LIST_TREASURY_TREE_TIMES = 3

	C2S_WATER_TREASURY_TREE = 5

	C2S_COLLECT_TREASURY_TREE_PRIZE = 10
)

func NewS2cListTreasuryTreeHeroMsg(hero_id [][]byte, hero_name []string, hero_head []string, hero_guild []int32, hero_friend []bool, hero_water_times []int32, guild_id []int32, flag_name []string, help_me_hero_id [][]byte, help_me_hero_name []string, help_me_hero_guild []int32, help_me_hero_season []int32) pbutil.Buffer {
	msg := &S2CListTreasuryTreeHeroProto{
		HeroId:           hero_id,
		HeroName:         hero_name,
		HeroHead:         hero_head,
		HeroGuild:        hero_guild,
		HeroFriend:       hero_friend,
		HeroWaterTimes:   hero_water_times,
		GuildId:          guild_id,
		FlagName:         flag_name,
		HelpMeHeroId:     help_me_hero_id,
		HelpMeHeroName:   help_me_hero_name,
		HelpMeHeroGuild:  help_me_hero_guild,
		HelpMeHeroSeason: help_me_hero_season,
	}
	return NewS2cListTreasuryTreeHeroProtoMsg(msg)
}

var s2c_list_treasury_tree_hero = [...]byte{31, 2} // 2
func NewS2cListTreasuryTreeHeroProtoMsg(object *S2CListTreasuryTreeHeroProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_treasury_tree_hero[:], "s2c_list_treasury_tree_hero")

}

func NewS2cListHelpMeMsg(target_id []byte, help_me_hero_id [][]byte, help_me_hero_name []string, help_me_hero_guild []int32, help_me_hero_flag_name []string, help_me_hero_season []int32) pbutil.Buffer {
	msg := &S2CListHelpMeProto{
		TargetId:           target_id,
		HelpMeHeroId:       help_me_hero_id,
		HelpMeHeroName:     help_me_hero_name,
		HelpMeHeroGuild:    help_me_hero_guild,
		HelpMeHeroFlagName: help_me_hero_flag_name,
		HelpMeHeroSeason:   help_me_hero_season,
	}
	return NewS2cListHelpMeProtoMsg(msg)
}

var s2c_list_help_me = [...]byte{31, 14} // 14
func NewS2cListHelpMeProtoMsg(object *S2CListHelpMeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_help_me[:], "s2c_list_help_me")

}

// 无效的目标
var ERR_LIST_HELP_ME_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 31, 15, 1} // 15-1

func NewS2cListTreasuryTreeTimesMsg(hero_id [][]byte, water_times []int32) pbutil.Buffer {
	msg := &S2CListTreasuryTreeTimesProto{
		HeroId:     hero_id,
		WaterTimes: water_times,
	}
	return NewS2cListTreasuryTreeTimesProtoMsg(msg)
}

var s2c_list_treasury_tree_times = [...]byte{31, 4} // 4
func NewS2cListTreasuryTreeTimesProtoMsg(object *S2CListTreasuryTreeTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_treasury_tree_times[:], "s2c_list_treasury_tree_times")

}

func NewS2cWaterTreasuryTreeMsg(target []byte, water_times int32) pbutil.Buffer {
	msg := &S2CWaterTreasuryTreeProto{
		Target:     target,
		WaterTimes: water_times,
	}
	return NewS2cWaterTreasuryTreeProtoMsg(msg)
}

var s2c_water_treasury_tree = [...]byte{31, 6} // 6
func NewS2cWaterTreasuryTreeProtoMsg(object *S2CWaterTreasuryTreeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_water_treasury_tree[:], "s2c_water_treasury_tree")

}

// 无效的目标
var ERR_WATER_TREASURY_TREE_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 31, 7, 1} // 7-1

// 你已经给目标浇过水了
var ERR_WATER_TREASURY_TREE_FAIL_WATERED = pbutil.StaticBuffer{3, 31, 7, 2} // 7-2

// 当季摇钱树的健康度已满，不能再照料了
var ERR_WATER_TREASURY_TREE_FAIL_FULL = pbutil.StaticBuffer{3, 31, 7, 3} // 7-3

func NewS2cUpdateSelfTreasuryTreeTimesMsg(water_times int32, help_me_hero_id []byte, help_me_hero_name string, help_me_hero_guild int32, help_me_hero_flag_name string) pbutil.Buffer {
	msg := &S2CUpdateSelfTreasuryTreeTimesProto{
		WaterTimes:         water_times,
		HelpMeHeroId:       help_me_hero_id,
		HelpMeHeroName:     help_me_hero_name,
		HelpMeHeroGuild:    help_me_hero_guild,
		HelpMeHeroFlagName: help_me_hero_flag_name,
	}
	return NewS2cUpdateSelfTreasuryTreeTimesProtoMsg(msg)
}

var s2c_update_self_treasury_tree_times = [...]byte{31, 8} // 8
func NewS2cUpdateSelfTreasuryTreeTimesProtoMsg(object *S2CUpdateSelfTreasuryTreeTimesProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_self_treasury_tree_times[:], "s2c_update_self_treasury_tree_times")

}

func NewS2cUpdateSelfTreasuryTreeFullMsg(season int32, collect_time int32) pbutil.Buffer {
	msg := &S2CUpdateSelfTreasuryTreeFullProto{
		Season:      season,
		CollectTime: collect_time,
	}
	return NewS2cUpdateSelfTreasuryTreeFullProtoMsg(msg)
}

var s2c_update_self_treasury_tree_full = [...]byte{31, 9} // 9
func NewS2cUpdateSelfTreasuryTreeFullProtoMsg(object *S2CUpdateSelfTreasuryTreeFullProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_self_treasury_tree_full[:], "s2c_update_self_treasury_tree_full")

}

var COLLECT_TREASURY_TREE_PRIZE_S2C = pbutil.StaticBuffer{2, 31, 11} // 11

// 摇钱树还未满，不能领取奖励
var ERR_COLLECT_TREASURY_TREE_PRIZE_FAIL_NO_FULL = pbutil.StaticBuffer{3, 31, 12, 1} // 12-1

// 领奖倒计时未结束
var ERR_COLLECT_TREASURY_TREE_PRIZE_FAIL_COUNTDOWN = pbutil.StaticBuffer{3, 31, 12, 2} // 12-2
