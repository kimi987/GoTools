// AUTO_GEN, DONT MODIFY!!!
package singleton

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

var _ = strings.ToUpper("")      // import strings
var _ = strconv.IntSize          // import strconv
var _ = shared_proto.Int32Pair{} // import shared_proto
var _ = errors.Errorf("")        // import errors
var _ = time.Second              // import time

// start with GoodsConfig ----------------------------------

func LoadGoodsConfig(gos *config.GameObjects) (*GoodsConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.GoodsConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	dAtA, err := NewGoodsConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedGoodsConfig(gos *config.GameObjects, dAtA *GoodsConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GoodsConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	return dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS)
}

func NewGoodsConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*GoodsConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGoodsConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GoodsConfig{}

	// releated field: EquipmentUpgradeGoods
	// releated field: EquipmentRefinedGoods
	// releated field: CaptainRefinedGoods
	// releated field: ChangeCaptainRaceGoods
	// skip field: CaptainRebirthGoods
	// skip field: MoveBaseGoods
	// skip field: RandomMoveBaseGoods
	// skip field: GuildMoveBaseGoods
	// skip field: MianGoods
	// releated field: JiuGuanExpCaptainRefinedGoods
	// releated field: ExpCaptainSoulUpgradeGoods
	// skip field: SpeedUpGoods
	// skip field: TrainExpGoods
	// skip field: TrainExpGoods4CaptainUpgrade
	// skip field: BuildingCdrGoods
	// skip field: TechCdrGoods
	// skip field: WorkshopCdrGoods
	// skip field: SpecGoods
	// skip field: GoldGoods
	// skip field: GoldNormalShopGoods
	// skip field: StoneGoods
	// skip field: StoneNormalShopGoods
	// skip field: GongXunGoods
	// skip field: MultiLevelNpcTimesGoods
	// skip field: InvaseHeroTimesGoods
	// skip field: JunTuanNpcTimesGoods
	// releated field: FishGoods
	// skip field: CopyDefenserGoods
	// skip field: TigerGoods
	// skip field: ZhanjiangGoods
	// skip field: TowerGoods
	// skip field: BuffGoods
	dAtA.IndecomposableBaowuCount = pArSeR.Uint64("indecomposable_baowu_count")

	return dAtA, nil
}

var vAlIdAtOrGoodsConfig = map[string]*config.Validator{

	"equipment_upgrade_goods":            config.ParseValidator("string", "", false, nil, nil),
	"equipment_refined_goods":            config.ParseValidator("string", "", false, nil, nil),
	"captain_refined_goods":              config.ParseValidator("string", "", true, nil, nil),
	"change_captain_race_goods":          config.ParseValidator("string", "", false, nil, nil),
	"jiu_guan_exp_captain_refined_goods": config.ParseValidator("int>0,count=4", "", true, nil, nil),
	"exp_captain_soul_upgrade_goods":     config.ParseValidator("int>0,count=4", "", true, nil, nil),
	"fish_goods":                         config.ParseValidator("string", "", false, nil, nil),
	"indecomposable_baowu_count":         config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *GoodsConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GoodsConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GoodsConfig) Encode() *shared_proto.GoodsConfigProto {
	out := &shared_proto.GoodsConfigProto{}
	if dAtA.EquipmentUpgradeGoods != nil {
		out.EquipmentUpgradeGoods = config.U64ToI32(dAtA.EquipmentUpgradeGoods.Id)
	}
	if dAtA.EquipmentRefinedGoods != nil {
		out.EquipmentRefinedGoods = config.U64ToI32(dAtA.EquipmentRefinedGoods.Id)
	}
	if dAtA.CaptainRefinedGoods != nil {
		out.CaptainRefinedGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.CaptainRefinedGoods))
	}
	if dAtA.ChangeCaptainRaceGoods != nil {
		out.ChangeCaptainRaceGoods = config.U64ToI32(dAtA.ChangeCaptainRaceGoods.Id)
	}
	if dAtA.CaptainRebirthGoods != nil {
		out.CaptainRebirthGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.CaptainRebirthGoods))
	}
	if dAtA.MoveBaseGoods != nil {
		out.MoveBaseGoods = config.U64ToI32(dAtA.MoveBaseGoods.Id)
	}
	if dAtA.RandomMoveBaseGoods != nil {
		out.RandomMoveBaseGoods = config.U64ToI32(dAtA.RandomMoveBaseGoods.Id)
	}
	if dAtA.GuildMoveBaseGoods != nil {
		out.GuildMoveBaseGoods = config.U64ToI32(dAtA.GuildMoveBaseGoods.Id)
	}
	if dAtA.MianGoods != nil {
		out.MianGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.MianGoods))
	}
	if dAtA.JiuGuanExpCaptainRefinedGoods != nil {
		out.JiuGuanExpCaptainRefinedGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.JiuGuanExpCaptainRefinedGoods))
	}
	if dAtA.ExpCaptainSoulUpgradeGoods != nil {
		out.ExpCaptainSoulUpgradeGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.ExpCaptainSoulUpgradeGoods))
	}
	if dAtA.SpeedUpGoods != nil {
		out.SpeedUpGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.SpeedUpGoods))
	}
	if dAtA.TrainExpGoods != nil {
		out.TrainExpGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.TrainExpGoods))
	}
	if dAtA.BuildingCdrGoods != nil {
		out.BuildingCdrGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.BuildingCdrGoods))
	}
	if dAtA.TechCdrGoods != nil {
		out.TechCdrGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.TechCdrGoods))
	}
	if dAtA.WorkshopCdrGoods != nil {
		out.WorkshopCdrGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.WorkshopCdrGoods))
	}
	if dAtA.SpecGoods != nil {
		out.SpecGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.SpecGoods))
	}
	out.GoldGoods = config.U64a2I32a(dAtA.GoldGoods)
	out.GoldNormalShopGoods = config.U64a2I32a(dAtA.GoldNormalShopGoods)
	out.StoneGoods = config.U64a2I32a(dAtA.StoneGoods)
	out.StoneNormalShopGoods = config.U64a2I32a(dAtA.StoneNormalShopGoods)
	if dAtA.GongXunGoods != nil {
		out.GongXunGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.GongXunGoods))
	}
	if dAtA.MultiLevelNpcTimesGoods != nil {
		out.MultiLevelNpcTimesGoods = config.U64ToI32(dAtA.MultiLevelNpcTimesGoods.Id)
	}
	if dAtA.InvaseHeroTimesGoods != nil {
		out.InvaseHeroTimesGoods = config.U64ToI32(dAtA.InvaseHeroTimesGoods.Id)
	}
	if dAtA.JunTuanNpcTimesGoods != nil {
		out.JunTuanNpcTimesGoods = config.U64ToI32(dAtA.JunTuanNpcTimesGoods.Id)
	}
	if dAtA.FishGoods != nil {
		out.FishGoods = config.U64ToI32(dAtA.FishGoods.Id)
	}
	if dAtA.CopyDefenserGoods != nil {
		out.CopyDefenserGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.CopyDefenserGoods))
	}
	if dAtA.TigerGoods != nil {
		out.TigerGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.TigerGoods))
	}
	if dAtA.ZhanjiangGoods != nil {
		out.ZhanjiangGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.ZhanjiangGoods))
	}
	if dAtA.TowerGoods != nil {
		out.TowerGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.TowerGoods))
	}
	if dAtA.BuffGoods != nil {
		out.BuffGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.BuffGoods))
	}
	out.IndecomposableBaowuCount = config.U64ToI32(dAtA.IndecomposableBaowuCount)

	return out
}

func ArrayEncodeGoodsConfig(datas []*GoodsConfig) []*shared_proto.GoodsConfigProto {

	out := make([]*shared_proto.GoodsConfigProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *GoodsConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.EquipmentUpgradeGoods = cOnFigS.GetGoodsData(pArSeR.Uint64("equipment_upgrade_goods"))
	if dAtA.EquipmentUpgradeGoods == nil {
		return errors.Errorf("%s 配置的关联字段[equipment_upgrade_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("equipment_upgrade_goods"), *pArSeR)
	}

	dAtA.EquipmentRefinedGoods = cOnFigS.GetGoodsData(pArSeR.Uint64("equipment_refined_goods"))
	if dAtA.EquipmentRefinedGoods == nil {
		return errors.Errorf("%s 配置的关联字段[equipment_refined_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("equipment_refined_goods"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("captain_refined_goods", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetGoodsData(v)
		if obj != nil {
			dAtA.CaptainRefinedGoods = append(dAtA.CaptainRefinedGoods, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[captain_refined_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain_refined_goods"), *pArSeR)
		}
	}

	dAtA.ChangeCaptainRaceGoods = cOnFigS.GetGoodsData(pArSeR.Uint64("change_captain_race_goods"))
	if dAtA.ChangeCaptainRaceGoods == nil {
		return errors.Errorf("%s 配置的关联字段[change_captain_race_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("change_captain_race_goods"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("jiu_guan_exp_captain_refined_goods", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetGoodsData(v)
		if obj != nil {
			dAtA.JiuGuanExpCaptainRefinedGoods = append(dAtA.JiuGuanExpCaptainRefinedGoods, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[jiu_guan_exp_captain_refined_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("jiu_guan_exp_captain_refined_goods"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("exp_captain_soul_upgrade_goods", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetGoodsData(v)
		if obj != nil {
			dAtA.ExpCaptainSoulUpgradeGoods = append(dAtA.ExpCaptainSoulUpgradeGoods, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[exp_captain_soul_upgrade_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("exp_captain_soul_upgrade_goods"), *pArSeR)
		}
	}

	dAtA.FishGoods = cOnFigS.GetGoodsData(pArSeR.Uint64("fish_goods"))
	if dAtA.FishGoods == nil {
		return errors.Errorf("%s 配置的关联字段[fish_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fish_goods"), *pArSeR)
	}

	return nil
}

// start with GuildConfig ----------------------------------

func LoadGuildConfig(gos *config.GameObjects) (*GuildConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	dAtA, err := NewGuildConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedGuildConfig(gos *config.GameObjects, dAtA *GuildConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	return dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS)
}

func NewGuildConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildConfig{}

	dAtA.HebiRobbedSuccessTaskProgress = 1
	if pArSeR.KeyExist("hebi_robbed_success_task_progress") {
		dAtA.HebiRobbedSuccessTaskProgress = pArSeR.Uint64("hebi_robbed_success_task_progress")
	}

	dAtA.HebiCompleteTaskProgress = 3
	if pArSeR.KeyExist("hebi_complete_task_progress") {
		dAtA.HebiCompleteTaskProgress = pArSeR.Uint64("hebi_complete_task_progress")
	}

	// releated field: CreateGuildCost
	// releated field: ChangeGuildNameCost
	if pArSeR.KeyExist("change_guild_name_duration") {
		dAtA.ChangeGuildNameDuration, err = time.ParseDuration(pArSeR.String("change_guild_name_duration"))
	} else {
		dAtA.ChangeGuildNameDuration, err = time.ParseDuration("72h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[change_guild_name_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("change_guild_name_duration"), dAtA)
	}

	dAtA.GuildLabelLimitChar = 5
	if pArSeR.KeyExist("guild_label_limit_char") {
		dAtA.GuildLabelLimitChar = pArSeR.Uint64("guild_label_limit_char")
	}

	dAtA.GuildLabelLimitCount = 4
	if pArSeR.KeyExist("guild_label_limit_count") {
		dAtA.GuildLabelLimitCount = pArSeR.Uint64("guild_label_limit_count")
	}

	dAtA.GuildNumPerPage = 10
	if pArSeR.KeyExist("guild_num_per_page") {
		dAtA.GuildNumPerPage = pArSeR.Uint64("guild_num_per_page")
	}

	dAtA.TextLimitChar = 128
	if pArSeR.KeyExist("text_limit_char") {
		dAtA.TextLimitChar = pArSeR.Uint64("text_limit_char")
	}

	dAtA.InternalTextLimitChar = 128
	if pArSeR.KeyExist("internal_text_limit_char") {
		dAtA.InternalTextLimitChar = pArSeR.Uint64("internal_text_limit_char")
	}

	dAtA.FriendGuildTextLimitChar = 128
	if pArSeR.KeyExist("friend_guild_text_limit_char") {
		dAtA.FriendGuildTextLimitChar = pArSeR.Uint64("friend_guild_text_limit_char")
	}

	dAtA.EnemyGuildTextLimitChar = 128
	if pArSeR.KeyExist("enemy_guild_text_limit_char") {
		dAtA.EnemyGuildTextLimitChar = pArSeR.Uint64("enemy_guild_text_limit_char")
	}

	dAtA.ChangeLeaderCountdownMemberCount = 20
	if pArSeR.KeyExist("change_leader_countdown_member_count") {
		dAtA.ChangeLeaderCountdownMemberCount = pArSeR.Uint64("change_leader_countdown_member_count")
	}

	if pArSeR.KeyExist("change_leader_countdown_duration") {
		dAtA.ChangeLeaderCountdownDuration, err = time.ParseDuration(pArSeR.String("change_leader_countdown_duration"))
	} else {
		dAtA.ChangeLeaderCountdownDuration, err = time.ParseDuration("72h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[change_leader_countdown_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("change_leader_countdown_duration"), dAtA)
	}

	dAtA.ImpeachNpcLeaderHour = 23
	if pArSeR.KeyExist("impeach_npc_leader_hour") {
		dAtA.ImpeachNpcLeaderHour = pArSeR.Uint64("impeach_npc_leader_hour")
	}

	dAtA.ImpeachNpcLeaderMinute = 40
	if pArSeR.KeyExist("impeach_npc_leader_minute") {
		dAtA.ImpeachNpcLeaderMinute = pArSeR.Uint64("impeach_npc_leader_minute")
	}

	dAtA.ImpeachUserLeaderMemberCount = 10
	if pArSeR.KeyExist("impeach_user_leader_member_count") {
		dAtA.ImpeachUserLeaderMemberCount = pArSeR.Int("impeach_user_leader_member_count")
	}

	if pArSeR.KeyExist("impeach_user_leader_offline") {
		dAtA.ImpeachUserLeaderOffline, err = config.ParseDuration(pArSeR.String("impeach_user_leader_offline"))
	} else {
		dAtA.ImpeachUserLeaderOffline, err = config.ParseDuration("48h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[impeach_user_leader_offline] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("impeach_user_leader_offline"), dAtA)
	}

	if pArSeR.KeyExist("impeach_user_leader_duration") {
		dAtA.ImpeachUserLeaderDuration, err = config.ParseDuration(pArSeR.String("impeach_user_leader_duration"))
	} else {
		dAtA.ImpeachUserLeaderDuration, err = config.ParseDuration("12h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[impeach_user_leader_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("impeach_user_leader_duration"), dAtA)
	}

	dAtA.ImpeachExtraCandidateCount = 2
	if pArSeR.KeyExist("impeach_extra_candidate_count") {
		dAtA.ImpeachExtraCandidateCount = pArSeR.Uint64("impeach_extra_candidate_count")
	}

	dAtA.UserMaxJoinRequestCount = 5
	if pArSeR.KeyExist("user_max_join_request_count") {
		dAtA.UserMaxJoinRequestCount = pArSeR.Uint64("user_max_join_request_count")
	}

	dAtA.GuildMaxJoinRequestCount = 20
	if pArSeR.KeyExist("guild_max_join_request_count") {
		dAtA.GuildMaxJoinRequestCount = pArSeR.Uint64("guild_max_join_request_count")
	}

	if pArSeR.KeyExist("join_request_duration") {
		dAtA.JoinRequestDuration, err = config.ParseDuration(pArSeR.String("join_request_duration"))
	} else {
		dAtA.JoinRequestDuration, err = config.ParseDuration("48h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[join_request_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("join_request_duration"), dAtA)
	}

	dAtA.UserMaxBeenInvateCount = 10
	if pArSeR.KeyExist("user_max_been_invate_count") {
		dAtA.UserMaxBeenInvateCount = pArSeR.Uint64("user_max_been_invate_count")
	}

	dAtA.GuildMaxInvateCount = 50
	if pArSeR.KeyExist("guild_max_invate_count") {
		dAtA.GuildMaxInvateCount = pArSeR.Uint64("guild_max_invate_count")
	}

	if pArSeR.KeyExist("invate_duration") {
		dAtA.InvateDuration, err = config.ParseDuration(pArSeR.String("invate_duration"))
	} else {
		dAtA.InvateDuration, err = config.ParseDuration("24h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[invate_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("invate_duration"), dAtA)
	}

	dAtA.ContributionDay = 7
	if pArSeR.KeyExist("contribution_day") {
		dAtA.ContributionDay = pArSeR.Int("contribution_day")
	}

	if pArSeR.KeyExist("npc_kick_offline_duration") {
		dAtA.NpcKickOfflineDuration, err = config.ParseDuration(pArSeR.String("npc_kick_offline_duration"))
	} else {
		dAtA.NpcKickOfflineDuration, err = config.ParseDuration("24h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[npc_kick_offline_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("npc_kick_offline_duration"), dAtA)
	}

	if pArSeR.KeyExist("npc_set_class_level_duration") {
		dAtA.NpcSetClassLevelDuration, err = config.ParseDuration(pArSeR.String("npc_set_class_level_duration"))
	} else {
		dAtA.NpcSetClassLevelDuration, err = config.ParseDuration("14h10m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[npc_set_class_level_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("npc_set_class_level_duration"), dAtA)
	}

	dAtA.FreeNpcGuildKeepCount = 4
	if pArSeR.KeyExist("free_npc_guild_keep_count") {
		dAtA.FreeNpcGuildKeepCount = pArSeR.Uint64("free_npc_guild_keep_count")
	}

	dAtA.FreeNpcGuildEmptyCount = 5
	if pArSeR.KeyExist("free_npc_guild_empty_count") {
		dAtA.FreeNpcGuildEmptyCount = pArSeR.Uint64("free_npc_guild_empty_count")
	}

	dAtA.DailyMaxKickCount = 3
	if pArSeR.KeyExist("daily_max_kick_count") {
		dAtA.DailyMaxKickCount = pArSeR.Uint64("daily_max_kick_count")
	}

	dAtA.GuildClassTitleMaxCount = 10
	if pArSeR.KeyExist("guild_class_title_max_count") {
		dAtA.GuildClassTitleMaxCount = pArSeR.Uint64("guild_class_title_max_count")
	}

	dAtA.GuildClassTitleMaxCharCount = 10
	if pArSeR.KeyExist("guild_class_title_max_char_count") {
		dAtA.GuildClassTitleMaxCharCount = pArSeR.Uint64("guild_class_title_max_char_count")
	}

	dAtA.GuildMaxDonateRecordCount = 10
	if pArSeR.KeyExist("guild_max_donate_record_count") {
		dAtA.GuildMaxDonateRecordCount = pArSeR.Uint64("guild_max_donate_record_count")
	}

	dAtA.GuildDonateNeedHeroLevel = 4
	if pArSeR.KeyExist("guild_donate_need_hero_level") {
		dAtA.GuildDonateNeedHeroLevel = pArSeR.Uint64("guild_donate_need_hero_level")
	}

	dAtA.GuildMaxBigEventCount = 10
	if pArSeR.KeyExist("guild_max_big_event_count") {
		dAtA.GuildMaxBigEventCount = pArSeR.Uint64("guild_max_big_event_count")
	}

	dAtA.GuildMaxDynamicCount = 10
	if pArSeR.KeyExist("guild_max_dynamic_count") {
		dAtA.GuildMaxDynamicCount = pArSeR.Uint64("guild_max_dynamic_count")
	}

	// releated field: DefaultGuildCountry
	dAtA.KeepDailyPrestigeCount = 2
	if pArSeR.KeyExist("keep_daily_prestige_count") {
		dAtA.KeepDailyPrestigeCount = pArSeR.Uint64("keep_daily_prestige_count")
	}

	dAtA.KeepHourlyPrestigeCount = 24
	if pArSeR.KeyExist("keep_hourly_prestige_count") {
		dAtA.KeepHourlyPrestigeCount = pArSeR.Uint64("keep_hourly_prestige_count")
	}

	dAtA.UpdateCountryYuanbao = 2000
	if pArSeR.KeyExist("update_country_yuanbao") {
		dAtA.UpdateCountryYuanbao = pArSeR.Uint64("update_country_yuanbao")
	}

	// releated field: UpdateCountryCost
	if pArSeR.KeyExist("update_country_duration") {
		dAtA.UpdateCountryDuration, err = config.ParseDuration(pArSeR.String("update_country_duration"))
	} else {
		dAtA.UpdateCountryDuration, err = config.ParseDuration("1m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[update_country_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("update_country_duration"), dAtA)
	}

	dAtA.UpdateCountryLostPrestigeCoef = 0.2
	if pArSeR.KeyExist("update_country_lost_prestige_coef") {
		dAtA.UpdateCountryLostPrestigeCoef = pArSeR.Float64("update_country_lost_prestige_coef")
	}

	// releated field: FirstJoinGuildPrize
	// skip field: FirstJoinGuildPrizeProto
	dAtA.ContributionPerHelp = 200
	if pArSeR.KeyExist("contribution_per_help") {
		dAtA.ContributionPerHelp = pArSeR.Uint64("contribution_per_help")
	}

	dAtA.ContributionMaxCountPerDay = 5
	if pArSeR.KeyExist("contribution_max_count_per_day") {
		dAtA.ContributionMaxCountPerDay = pArSeR.Uint64("contribution_max_count_per_day")
	}

	if pArSeR.KeyExist("big_box_collectable_duration") {
		dAtA.BigBoxCollectableDuration, err = config.ParseDuration(pArSeR.String("big_box_collectable_duration"))
	} else {
		dAtA.BigBoxCollectableDuration, err = config.ParseDuration("4h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[big_box_collectable_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("big_box_collectable_duration"), dAtA)
	}

	dAtA.EventPrizeMaxCount = 300
	if pArSeR.KeyExist("event_prize_max_count") {
		dAtA.EventPrizeMaxCount = pArSeR.Uint64("event_prize_max_count")
	}

	if pArSeR.KeyExist("notify_join_guild_duration") {
		dAtA.NotifyJoinGuildDuration, err = config.ParseDuration(pArSeR.String("notify_join_guild_duration"))
	} else {
		dAtA.NotifyJoinGuildDuration, err = config.ParseDuration("12h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[notify_join_guild_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("notify_join_guild_duration"), dAtA)
	}

	if pArSeR.KeyExist("notify_join_guild_duration_on_online_or_leave") {
		dAtA.NotifyJoinGuildDurationOnOnlineOrLeave, err = config.ParseDuration(pArSeR.String("notify_join_guild_duration_on_online_or_leave"))
	} else {
		dAtA.NotifyJoinGuildDurationOnOnlineOrLeave, err = config.ParseDuration("10m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[notify_join_guild_duration_on_online_or_leave] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("notify_join_guild_duration_on_online_or_leave"), dAtA)
	}

	dAtA.NotifyJoinGuildMaxPrestigeRank = 30
	if pArSeR.KeyExist("notify_join_guild_max_prestige_rank") {
		dAtA.NotifyJoinGuildMaxPrestigeRank = pArSeR.Uint64("notify_join_guild_max_prestige_rank")
	}

	dAtA.RecommendInviteHeroCount = 10
	if pArSeR.KeyExist("recommend_invite_hero_count") {
		dAtA.RecommendInviteHeroCount = pArSeR.Uint64("recommend_invite_hero_count")
	}

	dAtA.SearchNoGuildHerosPerPageSize = 10
	if pArSeR.KeyExist("search_no_guild_heros_per_page_size") {
		dAtA.SearchNoGuildHerosPerPageSize = pArSeR.Uint64("search_no_guild_heros_per_page_size")
	}

	if pArSeR.KeyExist("search_no_guild_heros_duration") {
		dAtA.SearchNoGuildHerosDuration, err = config.ParseDuration(pArSeR.String("search_no_guild_heros_duration"))
	} else {
		dAtA.SearchNoGuildHerosDuration, err = config.ParseDuration("1s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[search_no_guild_heros_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("search_no_guild_heros_duration"), dAtA)
	}

	dAtA.GuildChangeCountryWaitDuration, err = config.ParseDuration(pArSeR.String("guild_change_country_wait_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[guild_change_country_wait_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("guild_change_country_wait_duration"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrGuildConfig = map[string]*config.Validator{

	"hebi_robbed_success_task_progress":             config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"hebi_complete_task_progress":                   config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"create_guild_cost":                             config.ParseValidator("string", "", false, nil, []string{"1"}),
	"change_guild_name_cost":                        config.ParseValidator("string", "", false, nil, []string{"1"}),
	"change_guild_name_duration":                    config.ParseValidator("string", "", false, nil, []string{"72h"}),
	"guild_label_limit_char":                        config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"guild_label_limit_count":                       config.ParseValidator("int>0", "", false, nil, []string{"4"}),
	"guild_num_per_page":                            config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"text_limit_char":                               config.ParseValidator("int>0", "", false, nil, []string{"128"}),
	"internal_text_limit_char":                      config.ParseValidator("int>0", "", false, nil, []string{"128"}),
	"friend_guild_text_limit_char":                  config.ParseValidator("int>0", "", false, nil, []string{"128"}),
	"enemy_guild_text_limit_char":                   config.ParseValidator("int>0", "", false, nil, []string{"128"}),
	"change_leader_countdown_member_count":          config.ParseValidator("int>0", "", false, nil, []string{"20"}),
	"change_leader_countdown_duration":              config.ParseValidator("string", "", false, nil, []string{"72h"}),
	"impeach_npc_leader_hour":                       config.ParseValidator("int>0", "", false, nil, []string{"23"}),
	"impeach_npc_leader_minute":                     config.ParseValidator("int>0", "", false, nil, []string{"40"}),
	"impeach_user_leader_member_count":              config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"impeach_user_leader_offline":                   config.ParseValidator("string", "", false, nil, []string{"48h"}),
	"impeach_user_leader_duration":                  config.ParseValidator("string", "", false, nil, []string{"12h"}),
	"impeach_extra_candidate_count":                 config.ParseValidator("int>0", "", false, nil, []string{"2"}),
	"user_max_join_request_count":                   config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"guild_max_join_request_count":                  config.ParseValidator("int>0", "", false, nil, []string{"20"}),
	"join_request_duration":                         config.ParseValidator("string", "", false, nil, []string{"48h"}),
	"user_max_been_invate_count":                    config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"guild_max_invate_count":                        config.ParseValidator("int>0", "", false, nil, []string{"50"}),
	"invate_duration":                               config.ParseValidator("string", "", false, nil, []string{"24h"}),
	"contribution_day":                              config.ParseValidator("int>0", "", false, nil, []string{"7"}),
	"npc_kick_offline_duration":                     config.ParseValidator("string", "", false, nil, []string{"24h"}),
	"npc_set_class_level_duration":                  config.ParseValidator("string", "", false, nil, []string{"14h10m"}),
	"free_npc_guild_keep_count":                     config.ParseValidator("int>0", "", false, nil, []string{"4"}),
	"free_npc_guild_empty_count":                    config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"daily_max_kick_count":                          config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"guild_class_title_max_count":                   config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"guild_class_title_max_char_count":              config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"guild_max_donate_record_count":                 config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"guild_donate_need_hero_level":                  config.ParseValidator("int>0", "", false, nil, []string{"4"}),
	"guild_max_big_event_count":                     config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"guild_max_dynamic_count":                       config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"default_guild_country":                         config.ParseValidator("string", "", false, nil, []string{"1"}),
	"keep_daily_prestige_count":                     config.ParseValidator("int>0", "", false, nil, []string{"2"}),
	"keep_hourly_prestige_count":                    config.ParseValidator("int>0", "", false, nil, []string{"24"}),
	"update_country_yuanbao":                        config.ParseValidator("int>0", "", false, nil, []string{"2000"}),
	"update_country_cost":                           config.ParseValidator("string", "", false, nil, nil),
	"update_country_duration":                       config.ParseValidator("string", "", false, nil, []string{"1m"}),
	"update_country_lost_prestige_coef":             config.ParseValidator("float64>0", "", false, nil, []string{"0.2"}),
	"first_join_guild_prize":                        config.ParseValidator("string", "", false, nil, []string{"1"}),
	"contribution_per_help":                         config.ParseValidator("int>0", "", false, nil, []string{"200"}),
	"contribution_max_count_per_day":                config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"big_box_collectable_duration":                  config.ParseValidator("string", "", false, nil, []string{"4h"}),
	"event_prize_max_count":                         config.ParseValidator("int>0", "", false, nil, []string{"300"}),
	"notify_join_guild_duration":                    config.ParseValidator("string", "", false, nil, []string{"12h"}),
	"notify_join_guild_duration_on_online_or_leave": config.ParseValidator("string", "", false, nil, []string{"10m"}),
	"notify_join_guild_max_prestige_rank":           config.ParseValidator("int>0", "", false, nil, []string{"30"}),
	"recommend_invite_hero_count":                   config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"search_no_guild_heros_per_page_size":           config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"search_no_guild_heros_duration":                config.ParseValidator("string", "", false, nil, []string{"1s"}),
	"guild_change_country_wait_duration":            config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *GuildConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildConfig) Encode() *shared_proto.GuildConfigProto {
	out := &shared_proto.GuildConfigProto{}
	if dAtA.CreateGuildCost != nil {
		out.CreateGuildCost = dAtA.CreateGuildCost.Encode()
	}
	if dAtA.ChangeGuildNameCost != nil {
		out.ChangeGuildNameCost = dAtA.ChangeGuildNameCost.Encode()
	}
	out.ChangeGuildNameDuration = config.Duration2I32Seconds(dAtA.ChangeGuildNameDuration)
	out.GuildLabelLimitChar = config.U64ToI32(dAtA.GuildLabelLimitChar)
	out.GuildLabelLimitCount = config.U64ToI32(dAtA.GuildLabelLimitCount)
	out.GuildNumPerPage = config.U64ToI32(dAtA.GuildNumPerPage)
	out.TextLimitChar = config.U64ToI32(dAtA.TextLimitChar)
	out.InternalTextLimitChar = config.U64ToI32(dAtA.InternalTextLimitChar)
	out.FriendGuildTextLimitChar = config.U64ToI32(dAtA.FriendGuildTextLimitChar)
	out.EnemyGuildTextLimitChar = config.U64ToI32(dAtA.EnemyGuildTextLimitChar)
	out.ChangeLeaderCountdownMemberCount = config.U64ToI32(dAtA.ChangeLeaderCountdownMemberCount)
	out.ChangeLeaderCountdownDuration = config.Duration2I32Seconds(dAtA.ChangeLeaderCountdownDuration)
	out.ImpeachNpcLeaderHour = config.U64ToI32(dAtA.ImpeachNpcLeaderHour)
	out.ImpeachNpcLeaderMinute = config.U64ToI32(dAtA.ImpeachNpcLeaderMinute)
	out.ImpeachRequiredMemberCount = int32(dAtA.ImpeachUserLeaderMemberCount)
	out.ImpeachLeaderOfflineDuration = config.Duration2I32Seconds(dAtA.ImpeachUserLeaderOffline)
	out.UserMaxJoinRequestCount = config.U64ToI32(dAtA.UserMaxJoinRequestCount)
	out.GuildMaxInvateCount = config.U64ToI32(dAtA.GuildMaxInvateCount)
	out.DailyMaxKickCount = config.U64ToI32(dAtA.DailyMaxKickCount)
	out.GuildClassTitleMaxCount = config.U64ToI32(dAtA.GuildClassTitleMaxCount)
	out.GuildClassTitleMaxCharCount = config.U64ToI32(dAtA.GuildClassTitleMaxCharCount)
	out.GuildDonateNeedHeroLevel = config.U64ToI32(dAtA.GuildDonateNeedHeroLevel)
	out.UpdateCountryYuanbao = config.U64ToI32(dAtA.UpdateCountryYuanbao)
	if dAtA.UpdateCountryCost != nil {
		out.UpdateCountryCost = dAtA.UpdateCountryCost.Encode()
	}
	out.UpdateCountryDuration = config.Duration2I32Seconds(dAtA.UpdateCountryDuration)
	out.UpdateCountryLostPrestigeCoef = config.F64ToI32X1000(dAtA.UpdateCountryLostPrestigeCoef)
	if dAtA.FirstJoinGuildPrize != nil {
		out.FirstJoinGuildPrize = dAtA.FirstJoinGuildPrize.Encode()
	}
	out.ContributionPerHelp = config.U64ToI32(dAtA.ContributionPerHelp)
	out.ContributionMaxCountPerDay = config.U64ToI32(dAtA.ContributionMaxCountPerDay)
	out.BigBoxCollectableDuration = config.Duration2I32Seconds(dAtA.BigBoxCollectableDuration)
	out.EventPrizeMaxCount = config.U64ToI32(dAtA.EventPrizeMaxCount)
	out.NotifyJoinGuildDuration = config.Duration2I32Seconds(dAtA.NotifyJoinGuildDuration)
	out.NotifyJoinGuildDurationOnOnlineOrLeave = config.Duration2I32Seconds(dAtA.NotifyJoinGuildDurationOnOnlineOrLeave)
	out.SearchNoGuildHerosPerPageSize = config.U64ToI32(dAtA.SearchNoGuildHerosPerPageSize)
	out.SearchNoGuildHerosDuration = config.Duration2I32Seconds(dAtA.SearchNoGuildHerosDuration)

	return out
}

func ArrayEncodeGuildConfig(datas []*GuildConfig) []*shared_proto.GuildConfigProto {

	out := make([]*shared_proto.GuildConfigProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *GuildConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("create_guild_cost") {
		dAtA.CreateGuildCost = cOnFigS.GetCost(pArSeR.Int("create_guild_cost"))
	} else {
		dAtA.CreateGuildCost = cOnFigS.GetCost(1)
	}
	if dAtA.CreateGuildCost == nil {
		return errors.Errorf("%s 配置的关联字段[create_guild_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("create_guild_cost"), *pArSeR)
	}

	if pArSeR.KeyExist("change_guild_name_cost") {
		dAtA.ChangeGuildNameCost = cOnFigS.GetCost(pArSeR.Int("change_guild_name_cost"))
	} else {
		dAtA.ChangeGuildNameCost = cOnFigS.GetCost(1)
	}
	if dAtA.ChangeGuildNameCost == nil {
		return errors.Errorf("%s 配置的关联字段[change_guild_name_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("change_guild_name_cost"), *pArSeR)
	}

	if pArSeR.KeyExist("default_guild_country") {
		dAtA.DefaultGuildCountry = cOnFigS.GetCountryData(pArSeR.Uint64("default_guild_country"))
	} else {
		dAtA.DefaultGuildCountry = cOnFigS.GetCountryData(1)
	}
	if dAtA.DefaultGuildCountry == nil {
		return errors.Errorf("%s 配置的关联字段[default_guild_country] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("default_guild_country"), *pArSeR)
	}

	dAtA.UpdateCountryCost = cOnFigS.GetCost(pArSeR.Int("update_country_cost"))
	if dAtA.UpdateCountryCost == nil {
		return errors.Errorf("%s 配置的关联字段[update_country_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("update_country_cost"), *pArSeR)
	}

	if pArSeR.KeyExist("first_join_guild_prize") {
		dAtA.FirstJoinGuildPrize = cOnFigS.GetPrize(pArSeR.Int("first_join_guild_prize"))
	} else {
		dAtA.FirstJoinGuildPrize = cOnFigS.GetPrize(1)
	}
	if dAtA.FirstJoinGuildPrize == nil {
		return errors.Errorf("%s 配置的关联字段[first_join_guild_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_join_guild_prize"), *pArSeR)
	}

	return nil
}

// start with GuildGenConfig ----------------------------------

func LoadGuildGenConfig(gos *config.GameObjects) (*GuildGenConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildGenConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	dAtA, err := NewGuildGenConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedGuildGenConfig(gos *config.GameObjects, dAtA *GuildGenConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildGenConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	return dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS)
}

func NewGuildGenConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildGenConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildGenConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildGenConfig{}

	if pArSeR.KeyExist("leave_after_join_duration") {
		dAtA.LeaveAfterJoinDuration, err = config.ParseDuration(pArSeR.String("leave_after_join_duration"))
	} else {
		dAtA.LeaveAfterJoinDuration, err = config.ParseDuration("4h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[leave_after_join_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("leave_after_join_duration"), dAtA)
	}

	dAtA.GuildMarkCount = 4
	if pArSeR.KeyExist("guild_mark_count") {
		dAtA.GuildMarkCount = pArSeR.Uint64("guild_mark_count")
	}

	dAtA.GuildMarkMsgCharLimit = 20
	if pArSeR.KeyExist("guild_mark_msg_char_limit") {
		dAtA.GuildMarkMsgCharLimit = pArSeR.Uint64("guild_mark_msg_char_limit")
	}

	dAtA.SendMinYinliangToMember = 10
	if pArSeR.KeyExist("send_min_yinliang_to_member") {
		dAtA.SendMinYinliangToMember = pArSeR.Uint64("send_min_yinliang_to_member")
	}

	dAtA.SendMaxYinliangToMember = 10000
	if pArSeR.KeyExist("send_max_yinliang_to_member") {
		dAtA.SendMaxYinliangToMember = pArSeR.Uint64("send_max_yinliang_to_member")
	}

	dAtA.SendMinYinliangToGuild = 10
	if pArSeR.KeyExist("send_min_yinliang_to_guild") {
		dAtA.SendMinYinliangToGuild = pArSeR.Uint64("send_min_yinliang_to_guild")
	}

	dAtA.SendMaxYinliangToGuild = 10000
	if pArSeR.KeyExist("send_max_yinliang_to_guild") {
		dAtA.SendMaxYinliangToGuild = pArSeR.Uint64("send_max_yinliang_to_guild")
	}

	dAtA.SendMinSalary = 10
	if pArSeR.KeyExist("send_min_salary") {
		dAtA.SendMinSalary = pArSeR.Uint64("send_min_salary")
	}

	dAtA.SendMaxSalary = 10000
	if pArSeR.KeyExist("send_max_salary") {
		dAtA.SendMaxSalary = pArSeR.Uint64("send_max_salary")
	}

	if pArSeR.KeyExist("convene_cooldown") {
		dAtA.ConveneCooldown, err = config.ParseDuration(pArSeR.String("convene_cooldown"))
	} else {
		dAtA.ConveneCooldown, err = config.ParseDuration("10s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[convene_cooldown] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("convene_cooldown"), dAtA)
	}

	dAtA.WorkshopBuildDuration, err = config.ParseDuration(pArSeR.String("workshop_build_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[workshop_build_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("workshop_build_duration"), dAtA)
	}

	dAtA.WorkshopHeroBuildDuration, err = config.ParseDuration(pArSeR.String("workshop_hero_build_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[workshop_hero_build_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("workshop_hero_build_duration"), dAtA)
	}

	dAtA.WorkshopGuildBuildMaxTimes = 100
	if pArSeR.KeyExist("workshop_guild_build_max_times") {
		dAtA.WorkshopGuildBuildMaxTimes = pArSeR.Uint64("workshop_guild_build_max_times")
	}

	dAtA.WorkshopHeroBuildMaxTimes = 5
	if pArSeR.KeyExist("workshop_hero_build_max_times") {
		dAtA.WorkshopHeroBuildMaxTimes = pArSeR.Uint64("workshop_hero_build_max_times")
	}

	dAtA.WorkshopOutputInitTimes = pArSeR.Uint64("workshop_output_init_times")
	dAtA.WorkshopOutputMaxTimes = pArSeR.Uint64("workshop_output_max_times")
	dAtA.WorkshopOutputRecoveryDuration, err = config.ParseDuration(pArSeR.String("workshop_output_recovery_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[workshop_output_recovery_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("workshop_output_recovery_duration"), dAtA)
	}

	dAtA.WorkshopAddOutput = pArSeR.Uint64("workshop_add_output")
	dAtA.WorkshopAddProsperity = pArSeR.Uint64("workshop_add_prosperity")
	dAtA.WorkshopHurtDuration, err = config.ParseDuration(pArSeR.String("workshop_hurt_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[workshop_hurt_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("workshop_hurt_duration"), dAtA)
	}

	dAtA.WorkshopHurtTotalTimesLimit = pArSeR.Uint64("workshop_hurt_total_times_limit")
	dAtA.WorkshopHurtHeroTimesLimit = pArSeR.Uint64("workshop_hurt_hero_times_limit")
	dAtA.WorkshopHurtCooldown, err = config.ParseDuration(pArSeR.String("workshop_hurt_cooldown"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[workshop_hurt_cooldown] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("workshop_hurt_cooldown"), dAtA)
	}

	dAtA.WorkshopHurtProsperity = pArSeR.Uint64("workshop_hurt_prosperity")
	dAtA.WorkshopMaxOutput = []uint64{1, 2, 3}
	if pArSeR.KeyExist("workshop_max_output") {
		dAtA.WorkshopMaxOutput = pArSeR.Uint64Array("workshop_max_output", "", false)
	}

	// skip field: WorkshopProsperityCapcity
	dAtA.WorkshopBarrenProsperity = pArSeR.Uint64("workshop_barren_prosperity")
	dAtA.WorkshopPrizeInitCount = 1
	if pArSeR.KeyExist("workshop_prize_init_count") {
		dAtA.WorkshopPrizeInitCount = pArSeR.Uint64("workshop_prize_init_count")
	}

	dAtA.WorkshopPrizeMaxCount = 99
	if pArSeR.KeyExist("workshop_prize_max_count") {
		dAtA.WorkshopPrizeMaxCount = pArSeR.Uint64("workshop_prize_max_count")
	}

	dAtA.WorkshopReduceProsperityDuration, err = config.ParseDuration(pArSeR.String("workshop_reduce_prosperity_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[workshop_reduce_prosperity_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("workshop_reduce_prosperity_duration"), dAtA)
	}

	dAtA.WorkshopReduceProsperity = 1
	if pArSeR.KeyExist("workshop_reduce_prosperity") {
		dAtA.WorkshopReduceProsperity = pArSeR.Uint64("workshop_reduce_prosperity")
	}

	dAtA.WorkshopDistanceLimit = 100
	if pArSeR.KeyExist("workshop_distance_limit") {
		dAtA.WorkshopDistanceLimit = pArSeR.Uint64("workshop_distance_limit")
	}

	// releated field: WorkshopBase
	// releated field: GuildChangeCountryCost
	dAtA.GuildChangeCountryWaitDuration, err = config.ParseDuration(pArSeR.String("guild_change_country_wait_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[guild_change_country_wait_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("guild_change_country_wait_duration"), dAtA)
	}

	dAtA.GuildChangeCountryCooldown, err = config.ParseDuration(pArSeR.String("guild_change_country_cooldown"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[guild_change_country_cooldown] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("guild_change_country_cooldown"), dAtA)
	}

	dAtA.TaskOpenLevel = pArSeR.Uint64("task_open_level")

	return dAtA, nil
}

var vAlIdAtOrGuildGenConfig = map[string]*config.Validator{

	"leave_after_join_duration":           config.ParseValidator("string", "", false, nil, []string{"4h"}),
	"guild_mark_count":                    config.ParseValidator("int>0", "", false, nil, []string{"4"}),
	"guild_mark_msg_char_limit":           config.ParseValidator("int>0", "", false, nil, []string{"20"}),
	"send_min_yinliang_to_member":         config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"send_max_yinliang_to_member":         config.ParseValidator("int>0", "", false, nil, []string{"10000"}),
	"send_min_yinliang_to_guild":          config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"send_max_yinliang_to_guild":          config.ParseValidator("int>0", "", false, nil, []string{"10000"}),
	"send_min_salary":                     config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"send_max_salary":                     config.ParseValidator("int>0", "", false, nil, []string{"10000"}),
	"convene_cooldown":                    config.ParseValidator("string", "", false, nil, []string{"10s"}),
	"workshop_build_duration":             config.ParseValidator("string", "", false, nil, nil),
	"workshop_hero_build_duration":        config.ParseValidator("string", "", false, nil, nil),
	"workshop_guild_build_max_times":      config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"workshop_hero_build_max_times":       config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"workshop_output_init_times":          config.ParseValidator("int>0", "", false, nil, nil),
	"workshop_output_max_times":           config.ParseValidator("int>0", "", false, nil, nil),
	"workshop_output_recovery_duration":   config.ParseValidator("string", "", false, nil, nil),
	"workshop_add_output":                 config.ParseValidator("int>0", "", false, nil, nil),
	"workshop_add_prosperity":             config.ParseValidator("int>0", "", false, nil, nil),
	"workshop_hurt_duration":              config.ParseValidator("string", "", false, nil, nil),
	"workshop_hurt_total_times_limit":     config.ParseValidator("int>0", "", false, nil, nil),
	"workshop_hurt_hero_times_limit":      config.ParseValidator("int>0", "", false, nil, nil),
	"workshop_hurt_cooldown":              config.ParseValidator("string", "", false, nil, nil),
	"workshop_hurt_prosperity":            config.ParseValidator("int>0", "", false, nil, nil),
	"workshop_max_output":                 config.ParseValidator(",duplicate", "", true, nil, []string{"1", "2", "3"}),
	"workshop_barren_prosperity":          config.ParseValidator("int>0", "", false, nil, nil),
	"workshop_prize_init_count":           config.ParseValidator("uint", "", false, nil, []string{"1"}),
	"workshop_prize_max_count":            config.ParseValidator("uint", "", false, nil, []string{"99"}),
	"workshop_reduce_prosperity_duration": config.ParseValidator("string", "", false, nil, nil),
	"workshop_reduce_prosperity":          config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"workshop_distance_limit":             config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"workshop_base":                       config.ParseValidator("string", "", false, nil, nil),
	"guild_change_country_cost":           config.ParseValidator("string", "", false, nil, nil),
	"guild_change_country_wait_duration":  config.ParseValidator("string", "", false, nil, nil),
	"guild_change_country_cooldown":       config.ParseValidator("string", "", false, nil, nil),
	"task_open_level":                     config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *GuildGenConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildGenConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildGenConfig) Encode() *shared_proto.GuildGenConfigProto {
	out := &shared_proto.GuildGenConfigProto{}
	out.LeaveAfterJoinDuration = config.Duration2I32Seconds(dAtA.LeaveAfterJoinDuration)
	out.GuildMarkCount = config.U64ToI32(dAtA.GuildMarkCount)
	out.GuildMarkMsgCharLimit = config.U64ToI32(dAtA.GuildMarkMsgCharLimit)
	out.SendMinYinliangToMember = config.U64ToI32(dAtA.SendMinYinliangToMember)
	out.SendMaxYinliangToMember = config.U64ToI32(dAtA.SendMaxYinliangToMember)
	out.SendMinYinliangToGuild = config.U64ToI32(dAtA.SendMinYinliangToGuild)
	out.SendMaxYinliangToGuild = config.U64ToI32(dAtA.SendMaxYinliangToGuild)
	out.SendMinSalary = config.U64ToI32(dAtA.SendMinSalary)
	out.SendMaxSalary = config.U64ToI32(dAtA.SendMaxSalary)
	out.ConveneCooldown = config.Duration2I32Seconds(dAtA.ConveneCooldown)
	out.WorkshopBuildDuration = config.Duration2I32Seconds(dAtA.WorkshopBuildDuration)
	out.WorkshopHeroBuildDuration = config.Duration2I32Seconds(dAtA.WorkshopHeroBuildDuration)
	out.WorkshopGuildBuildMaxTimes = config.U64ToI32(dAtA.WorkshopGuildBuildMaxTimes)
	out.WorkshopHeroBuildMaxTimes = config.U64ToI32(dAtA.WorkshopHeroBuildMaxTimes)
	out.WorkshopOutputMaxTimes = config.U64ToI32(dAtA.WorkshopOutputMaxTimes)
	out.WorkshopOutputRecoveryDuration = config.Duration2I32Seconds(dAtA.WorkshopOutputRecoveryDuration)
	out.WorkshopHurtDuration = config.Duration2I32Seconds(dAtA.WorkshopHurtDuration)
	out.WorkshopHurtTotalTimesLimit = config.U64ToI32(dAtA.WorkshopHurtTotalTimesLimit)
	out.WorkshopHurtHeroTimesLimit = config.U64ToI32(dAtA.WorkshopHurtHeroTimesLimit)
	out.WorkshopHurtCooldown = config.Duration2I32Seconds(dAtA.WorkshopHurtCooldown)
	out.WorkshopHurtProsperity = config.U64ToI32(dAtA.WorkshopHurtProsperity)
	out.WorkshopProsperityCapcity = config.U64ToI32(dAtA.WorkshopProsperityCapcity)
	out.WorkshopBarrenProsperity = config.U64ToI32(dAtA.WorkshopBarrenProsperity)
	out.WorkshopPrizeMaxCount = config.U64ToI32(dAtA.WorkshopPrizeMaxCount)
	out.WorkshopReduceProsperityDuration = config.Duration2I32Seconds(dAtA.WorkshopReduceProsperityDuration)
	out.WorkshopReduceProsperity = config.U64ToI32(dAtA.WorkshopReduceProsperity)
	out.WorkshopDistanceLimit = config.U64ToI32(dAtA.WorkshopDistanceLimit)
	if dAtA.GuildChangeCountryCost != nil {
		out.GuildChangeCountryCost = dAtA.GuildChangeCountryCost.Encode()
	}
	out.GuildChangeCountryWaitDuration = config.Duration2I32Seconds(dAtA.GuildChangeCountryWaitDuration)
	out.GuildChangeCountryCooldown = config.Duration2I32Seconds(dAtA.GuildChangeCountryCooldown)
	out.TaskOpenLevel = config.U64ToI32(dAtA.TaskOpenLevel)

	return out
}

func ArrayEncodeGuildGenConfig(datas []*GuildGenConfig) []*shared_proto.GuildGenConfigProto {

	out := make([]*shared_proto.GuildGenConfigProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *GuildGenConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.WorkshopBase = cOnFigS.GetNpcBaseData(pArSeR.Uint64("workshop_base"))
	if dAtA.WorkshopBase == nil {
		return errors.Errorf("%s 配置的关联字段[workshop_base] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("workshop_base"), *pArSeR)
	}

	dAtA.GuildChangeCountryCost = cOnFigS.GetCost(pArSeR.Int("guild_change_country_cost"))
	if dAtA.GuildChangeCountryCost == nil {
		return errors.Errorf("%s 配置的关联字段[guild_change_country_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_change_country_cost"), *pArSeR)
	}

	return nil
}

// start with MilitaryConfig ----------------------------------

func LoadMilitaryConfig(gos *config.GameObjects) (*MilitaryConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.MilitaryConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	dAtA, err := NewMilitaryConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedMilitaryConfig(gos *config.GameObjects, dAtA *MilitaryConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MilitaryConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	return dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS)
}

func NewMilitaryConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*MilitaryConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMilitaryConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MilitaryConfig{}

	dAtA.GenSeekCount = 6
	if pArSeR.KeyExist("gen_seek_count") {
		dAtA.GenSeekCount = pArSeR.Uint64("gen_seek_count")
	}

	if pArSeR.KeyExist("seek_duration") {
		dAtA.SeekDuration, err = time.ParseDuration(pArSeR.String("seek_duration"))
	} else {
		dAtA.SeekDuration, err = time.ParseDuration("1h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[seek_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("seek_duration"), dAtA)
	}

	dAtA.SeekMaxTimes = 15
	if pArSeR.KeyExist("seek_max_times") {
		dAtA.SeekMaxTimes = pArSeR.Uint64("seek_max_times")
	}

	dAtA.DefenserCount = 5
	if pArSeR.KeyExist("defenser_count") {
		dAtA.DefenserCount = pArSeR.Uint64("defenser_count")
	}

	dAtA.FireLevelLimit = 15
	if pArSeR.KeyExist("fire_level_limit") {
		dAtA.FireLevelLimit = pArSeR.Uint64("fire_level_limit")
	}

	// releated field: CombatRes
	dAtA.CombatXLen = 10
	if pArSeR.KeyExist("combat_xlen") {
		dAtA.CombatXLen = pArSeR.Uint64("combat_xlen")
	}

	dAtA.CombatYLen = 5
	if pArSeR.KeyExist("combat_ylen") {
		dAtA.CombatYLen = pArSeR.Uint64("combat_ylen")
	}

	dAtA.CombatMaxRound = 100
	if pArSeR.KeyExist("combat_max_round") {
		dAtA.CombatMaxRound = pArSeR.Uint64("combat_max_round")
	}

	dAtA.CombatCoef = 10
	if pArSeR.KeyExist("combat_coef") {
		dAtA.CombatCoef = pArSeR.Float64("combat_coef")
	}

	dAtA.CombatCritRate = 0.3
	if pArSeR.KeyExist("combat_crit_rate") {
		dAtA.CombatCritRate = pArSeR.Float64("combat_crit_rate")
	}

	dAtA.CombatRestraintRate = 0.2
	if pArSeR.KeyExist("combat_restraint_rate") {
		dAtA.CombatRestraintRate = pArSeR.Float64("combat_restraint_rate")
	}

	dAtA.CombatScorePercent = []uint64{33, 66, 90}
	if pArSeR.KeyExist("combat_score_percent") {
		dAtA.CombatScorePercent = pArSeR.Uint64Array("combat_score_percent", "", false)
	}

	// skip field: TrainingHeroLevel
	// skip field: TrainingInitLevel
	// skip field: TrainingLevelLimit
	if pArSeR.KeyExist("training_max_duration") {
		dAtA.TrainingMaxDuration, err = config.ParseDuration(pArSeR.String("training_max_duration"))
	} else {
		dAtA.TrainingMaxDuration, err = config.ParseDuration("8h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[training_max_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("training_max_duration"), dAtA)
	}

	dAtA.CaptainInitTrainExp = 1000
	if pArSeR.KeyExist("captain_init_train_exp") {
		dAtA.CaptainInitTrainExp = pArSeR.Uint64("captain_init_train_exp")
	}

	dAtA.MinWallAttackRound = 2
	if pArSeR.KeyExist("min_wall_attack_round") {
		dAtA.MinWallAttackRound = pArSeR.Uint64("min_wall_attack_round")
	}

	dAtA.MaxWallAttachFixDamageRound = 4
	if pArSeR.KeyExist("max_wall_attach_fix_damage_round") {
		dAtA.MaxWallAttachFixDamageRound = pArSeR.Uint64("max_wall_attach_fix_damage_round")
	}

	dAtA.MaxWallBeenHurtPercent = 0.1
	if pArSeR.KeyExist("max_wall_been_hurt_percent") {
		dAtA.MaxWallBeenHurtPercent = pArSeR.Float64("max_wall_been_hurt_percent")
	}

	// skip field: TroopsUnlockLevel
	dAtA.CaptainSeekerCandidateCount = 5
	if pArSeR.KeyExist("captain_seeker_candidate_count") {
		dAtA.CaptainSeekerCandidateCount = pArSeR.Uint64("captain_seeker_candidate_count")
	}

	return dAtA, nil
}

var vAlIdAtOrMilitaryConfig = map[string]*config.Validator{

	"gen_seek_count":                   config.ParseValidator("int>0", "", false, nil, []string{"6"}),
	"seek_duration":                    config.ParseValidator("string", "", false, nil, []string{"1h"}),
	"seek_max_times":                   config.ParseValidator("int>0", "", false, nil, []string{"15"}),
	"defenser_count":                   config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"fire_level_limit":                 config.ParseValidator("int>0", "", false, nil, []string{"15"}),
	"combat_res":                       config.ParseValidator("string", "", false, nil, []string{"Battle_Field_1"}),
	"combat_xlen":                      config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"combat_ylen":                      config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"combat_max_round":                 config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"combat_coef":                      config.ParseValidator("float64>0", "", false, nil, []string{"10"}),
	"combat_crit_rate":                 config.ParseValidator("float64>0", "", false, nil, []string{"0.3"}),
	"combat_restraint_rate":            config.ParseValidator("float64>0", "", false, nil, []string{"0.2"}),
	"combat_score_percent":             config.ParseValidator("float64>0", "", true, nil, []string{"33", "66", "90"}),
	"training_max_duration":            config.ParseValidator("string", "", false, nil, []string{"8h"}),
	"captain_init_train_exp":           config.ParseValidator("int>0", "", false, nil, []string{"1000"}),
	"min_wall_attack_round":            config.ParseValidator("int>0", "", false, nil, []string{"2"}),
	"max_wall_attach_fix_damage_round": config.ParseValidator("int>0", "", false, nil, []string{"4"}),
	"max_wall_been_hurt_percent":       config.ParseValidator("float64>0", "", false, nil, []string{"0.1"}),
	"captain_seeker_candidate_count":   config.ParseValidator("int>0", "", false, nil, []string{"5"}),
}

func (dAtA *MilitaryConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MilitaryConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MilitaryConfig) Encode() *shared_proto.MilitaryConfigProto {
	out := &shared_proto.MilitaryConfigProto{}
	out.GenSeekCount = config.U64ToI32(dAtA.GenSeekCount)
	out.SeekDuration = config.Duration2I32Seconds(dAtA.SeekDuration)
	out.SeekMaxTimes = config.U64ToI32(dAtA.SeekMaxTimes)
	out.DefenserCount = config.U64ToI32(dAtA.DefenserCount)
	out.FireLevelLimit = config.U64ToI32(dAtA.FireLevelLimit)
	out.TrainingHeroLevel = config.U64a2I32a(dAtA.TrainingHeroLevel)
	out.TrainingInitLevel = config.U64a2I32a(dAtA.TrainingInitLevel)
	out.TrainingLevelLimit = config.U64a2I32a(dAtA.TrainingLevelLimit)
	out.TrainingMaxDuration = config.Duration2I32Seconds(dAtA.TrainingMaxDuration)
	out.CaptainInitTrainExp = config.U64ToI32(dAtA.CaptainInitTrainExp)
	out.TroopsUnlockLevel = config.U64a2I32a(dAtA.TroopsUnlockLevel)

	return out
}

func ArrayEncodeMilitaryConfig(datas []*MilitaryConfig) []*shared_proto.MilitaryConfigProto {

	out := make([]*shared_proto.MilitaryConfigProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *MilitaryConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("combat_res") {
		dAtA.CombatRes = cOnFigS.GetCombatScene(pArSeR.String("combat_res"))
	} else {
		dAtA.CombatRes = cOnFigS.GetCombatScene("Battle_Field_1")
	}
	if dAtA.CombatRes == nil {
		return errors.Errorf("%s 配置的关联字段[combat_res] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("combat_res"), *pArSeR)
	}

	return nil
}

// start with MiscConfig ----------------------------------

func LoadMiscConfig(gos *config.GameObjects) (*MiscConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.MiscConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	dAtA, err := NewMiscConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedMiscConfig(gos *config.GameObjects, dAtA *MiscConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MiscConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	return dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS)
}

func NewMiscConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*MiscConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMiscConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MiscConfig{}

	dAtA.DailyResetHour = 0
	if pArSeR.KeyExist("daily_reset_hour") {
		dAtA.DailyResetHour = pArSeR.Uint64("daily_reset_hour")
	}

	dAtA.DailyResetMinute = 0
	if pArSeR.KeyExist("daily_reset_minute") {
		dAtA.DailyResetMinute = pArSeR.Uint64("daily_reset_minute")
	}

	// skip field: DailyResetDuration
	dAtA.WeeklyResetHour = 0
	if pArSeR.KeyExist("weekly_reset_hour") {
		dAtA.WeeklyResetHour = pArSeR.Uint64("weekly_reset_hour")
	}

	dAtA.WeeklyResetMinute = 0
	if pArSeR.KeyExist("weekly_reset_minute") {
		dAtA.WeeklyResetMinute = pArSeR.Uint64("weekly_reset_minute")
	}

	// skip field: WeeklyResetDuration
	dAtA.MinNameCharLen = 2
	if pArSeR.KeyExist("min_name_char_len") {
		dAtA.MinNameCharLen = pArSeR.Uint64("min_name_char_len")
	}

	dAtA.MaxNameCharLen = 14
	if pArSeR.KeyExist("max_name_char_len") {
		dAtA.MaxNameCharLen = pArSeR.Uint64("max_name_char_len")
	}

	dAtA.WorkshopRefreshHourMinute = []uint64{1200, 1600, 2000}
	if pArSeR.KeyExist("workshop_refresh_hour_minute") {
		dAtA.WorkshopRefreshHourMinute = pArSeR.Uint64Array("workshop_refresh_hour_minute", "", false)
	}

	// skip field: WorkshopRefreshDuration
	// releated field: SecondWorkerCost
	dAtA.SecondWorkerUnlockDuration, err = config.ParseDuration(pArSeR.String("second_worker_unlock_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[second_worker_unlock_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("second_worker_unlock_duration"), dAtA)
	}

	dAtA.ChangeHeroNameYuanbaoCost = []uint64{0, 100}
	if pArSeR.KeyExist("change_hero_name_yuanbao_cost") {
		dAtA.ChangeHeroNameYuanbaoCost = pArSeR.Uint64Array("change_hero_name_yuanbao_cost", "", false)
	}

	// releated field: ChangeHeroNameCost
	dAtA.ChangeHeroNameDuration, err = config.ParseDuration(pArSeR.String("change_hero_name_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[change_hero_name_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("change_hero_name_duration"), dAtA)
	}

	// releated field: FirstChangeHeroNamePrize
	// releated field: ChangeCaptainNameCost
	dAtA.ChangeCaptainRaceDuration, err = config.ParseDuration(pArSeR.String("change_captain_race_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[change_captain_race_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("change_captain_race_duration"), dAtA)
	}

	dAtA.ChangeCaptainRaceLevel = 50
	if pArSeR.KeyExist("change_captain_race_level") {
		dAtA.ChangeCaptainRaceLevel = pArSeR.Uint64("change_captain_race_level")
	}

	dAtA.MailMinBatchCount = 20
	if pArSeR.KeyExist("mail_min_batch_count") {
		dAtA.MailMinBatchCount = pArSeR.Uint64("mail_min_batch_count")
	}

	dAtA.MailMaxBatchCount = 50
	if pArSeR.KeyExist("mail_max_batch_count") {
		dAtA.MailMaxBatchCount = pArSeR.Uint64("mail_max_batch_count")
	}

	dAtA.TowerChallengeMaxTimes = 3
	if pArSeR.KeyExist("tower_challenge_max_times") {
		dAtA.TowerChallengeMaxTimes = pArSeR.Uint64("tower_challenge_max_times")
	}

	if pArSeR.KeyExist("tower_reset_challenge_duration") {
		dAtA.TowerResetChallengeDuration, err = config.ParseDuration(pArSeR.String("tower_reset_challenge_duration"))
	} else {
		dAtA.TowerResetChallengeDuration, err = config.ParseDuration("30m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tower_reset_challenge_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tower_reset_challenge_duration"), dAtA)
	}

	dAtA.TowerReplayCount = 4
	if pArSeR.KeyExist("tower_replay_count") {
		dAtA.TowerReplayCount = pArSeR.Uint64("tower_replay_count")
	}

	dAtA.TowerAutoKeepFloor = 5
	if pArSeR.KeyExist("tower_auto_keep_floor") {
		dAtA.TowerAutoKeepFloor = pArSeR.Uint64("tower_auto_keep_floor")
	}

	dAtA.EquipmentUpgradeMultiTimes = 10
	if pArSeR.KeyExist("equipment_upgrade_multi_times") {
		dAtA.EquipmentUpgradeMultiTimes = pArSeR.Uint64("equipment_upgrade_multi_times")
	}

	dAtA.MiaoBuildingWorkerDuration, err = config.ParseDuration(pArSeR.String("miao_building_worker_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[miao_building_worker_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("miao_building_worker_duration"), dAtA)
	}

	// releated field: MiaoBuildingWorkerCost
	dAtA.MiaoTechWorkerDuration, err = config.ParseDuration(pArSeR.String("miao_tech_worker_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[miao_tech_worker_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("miao_tech_worker_duration"), dAtA)
	}

	// releated field: MiaoTechWorkerCost
	dAtA.MiaoCaptainRebirthDuration, err = config.ParseDuration(pArSeR.String("miao_captain_rebirth_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[miao_captain_rebirth_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("miao_captain_rebirth_duration"), dAtA)
	}

	// releated field: MiaoCaptainRebirthCost
	dAtA.MiaoWorkshopDuration, err = config.ParseDuration(pArSeR.String("miao_workshop_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[miao_workshop_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("miao_workshop_duration"), dAtA)
	}

	// releated field: MiaoWorkshopCost
	dAtA.DefaultForgingTimes = 3
	if pArSeR.KeyExist("default_forging_times") {
		dAtA.DefaultForgingTimes = pArSeR.Uint64("default_forging_times")
	}

	dAtA.MaxDepotEquipCapacity = 100
	if pArSeR.KeyExist("max_depot_equip_capacity") {
		dAtA.MaxDepotEquipCapacity = pArSeR.Uint64("max_depot_equip_capacity")
	}

	if pArSeR.KeyExist("temp_depot_expire_duration") {
		dAtA.TempDepotExpireDuration, err = time.ParseDuration(pArSeR.String("temp_depot_expire_duration"))
	} else {
		dAtA.TempDepotExpireDuration, err = time.ParseDuration("24h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[temp_depot_expire_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("temp_depot_expire_duration"), dAtA)
	}

	dAtA.MaxSignLen = 20
	if pArSeR.KeyExist("max_sign_len") {
		dAtA.MaxSignLen = pArSeR.Uint64("max_sign_len")
	}

	dAtA.MaxVoiceLen = 2000
	if pArSeR.KeyExist("max_voice_len") {
		dAtA.MaxVoiceLen = pArSeR.Uint64("max_voice_len")
	}

	if pArSeR.KeyExist("strategy_restore_duration") {
		dAtA.StrategyRestoreDuration, err = config.ParseDuration(pArSeR.String("strategy_restore_duration"))
	} else {
		dAtA.StrategyRestoreDuration, err = config.ParseDuration("1h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[strategy_restore_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("strategy_restore_duration"), dAtA)
	}

	dAtA.MaxResourceCollectTimes = 10
	if pArSeR.KeyExist("max_resource_collect_times") {
		dAtA.MaxResourceCollectTimes = pArSeR.Uint64("max_resource_collect_times")
	}

	dAtA.DefaultResourceCollectTimes = 5
	if pArSeR.KeyExist("default_resource_collect_times") {
		dAtA.DefaultResourceCollectTimes = pArSeR.Uint64("default_resource_collect_times")
	}

	if pArSeR.KeyExist("resource_recovery_duration") {
		dAtA.ResourceRecoveryDuration, err = config.ParseDuration(pArSeR.String("resource_recovery_duration"))
	} else {
		dAtA.ResourceRecoveryDuration, err = config.ParseDuration("1h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[resource_recovery_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("resource_recovery_duration"), dAtA)
	}

	// skip field: MondayZeroOClock
	dAtA.MaxFavoritePosCount = 10
	if pArSeR.KeyExist("max_favorite_pos_count") {
		dAtA.MaxFavoritePosCount = pArSeR.Uint64("max_favorite_pos_count")
	}

	if pArSeR.KeyExist("flag_hero_name") {
		dAtA.FlagHeroName, err = data.ParseTextFormatter(pArSeR.String("flag_hero_name"))
	} else {
		dAtA.FlagHeroName, err = data.ParseTextFormatter("[%s];%s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[flag_hero_name] 解析失败(data.ParseTextFormatter)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("flag_hero_name"), dAtA)
	}

	if pArSeR.KeyExist("country_flag_hero_name") {
		dAtA.CountryFlagHeroName, err = data.ParseTextFormatter(pArSeR.String("country_flag_hero_name"))
	} else {
		dAtA.CountryFlagHeroName, err = data.ParseTextFormatter("[%s];[%s];%s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[country_flag_hero_name] 解析失败(data.ParseTextFormatter)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("country_flag_hero_name"), dAtA)
	}

	dAtA.WorldChatLevel = 1
	if pArSeR.KeyExist("world_chat_level") {
		dAtA.WorldChatLevel = pArSeR.Uint64("world_chat_level")
	}

	if pArSeR.KeyExist("world_chat_duration") {
		dAtA.WorldChatDuration, err = config.ParseDuration(pArSeR.String("world_chat_duration"))
	} else {
		dAtA.WorldChatDuration, err = config.ParseDuration("10s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[world_chat_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("world_chat_duration"), dAtA)
	}

	dAtA.ChatTextLength = 100
	if pArSeR.KeyExist("chat_text_length") {
		dAtA.ChatTextLength = pArSeR.Uint64("chat_text_length")
	}

	dAtA.ChatJsonLength = 600
	if pArSeR.KeyExist("chat_json_length") {
		dAtA.ChatJsonLength = pArSeR.Uint64("chat_json_length")
	}

	dAtA.ChatWindowLimit = 100
	if pArSeR.KeyExist("chat_window_limit") {
		dAtA.ChatWindowLimit = pArSeR.Uint64("chat_window_limit")
	}

	dAtA.ChatBatchCount = 100
	if pArSeR.KeyExist("chat_batch_count") {
		dAtA.ChatBatchCount = pArSeR.Uint64("chat_batch_count")
	}

	// releated field: BroadcastGoods
	dAtA.ChatPrivateMinLevel = 3
	if pArSeR.KeyExist("chat_private_min_level") {
		dAtA.ChatPrivateMinLevel = pArSeR.Uint64("chat_private_min_level")
	}

	if pArSeR.KeyExist("chat_duration") {
		dAtA.ChatDuration, err = config.ParseDuration(pArSeR.String("chat_duration"))
	} else {
		dAtA.ChatDuration, err = config.ParseDuration("5s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[chat_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("chat_duration"), dAtA)
	}

	if pArSeR.KeyExist("chat_share_duration") {
		dAtA.ChatShareDuration, err = config.ParseDuration(pArSeR.String("chat_share_duration"))
	} else {
		dAtA.ChatShareDuration, err = config.ParseDuration("5s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[chat_share_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("chat_share_duration"), dAtA)
	}

	dAtA.GiveAddCountDownPrizeDuration, err = config.ParseDuration(pArSeR.String("give_add_count_down_prize_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[give_add_count_down_prize_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("give_add_count_down_prize_duration"), dAtA)
	}

	dAtA.GiveAddCountDownPrizeTimes = pArSeR.Uint64("give_add_count_down_prize_times")
	dAtA.StrongerCoef = []float64{0, 0.6, 1}
	if pArSeR.KeyExist("stronger_coef") {
		dAtA.StrongerCoef = pArSeR.Float64Array("stronger_coef", "", false)
	}

	dAtA.BuildingInitEffect = 1
	if pArSeR.KeyExist("building_init_effect") {
		dAtA.BuildingInitEffect = pArSeR.Uint64("building_init_effect")
	}

	if pArSeR.KeyExist("db_world_chat_expire_duration") {
		dAtA.DbWorldChatExpireDuration, err = config.ParseDuration(pArSeR.String("db_world_chat_expire_duration"))
	} else {
		dAtA.DbWorldChatExpireDuration, err = config.ParseDuration("72h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[db_world_chat_expire_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("db_world_chat_expire_duration"), dAtA)
	}

	if pArSeR.KeyExist("db_guild_chat_expire_duration") {
		dAtA.DbGuildChatExpireDuration, err = config.ParseDuration(pArSeR.String("db_guild_chat_expire_duration"))
	} else {
		dAtA.DbGuildChatExpireDuration, err = config.ParseDuration("168h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[db_guild_chat_expire_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("db_guild_chat_expire_duration"), dAtA)
	}

	if pArSeR.KeyExist("db_private_chat_expire_duration") {
		dAtA.DbPrivateChatExpireDuration, err = config.ParseDuration(pArSeR.String("db_private_chat_expire_duration"))
	} else {
		dAtA.DbPrivateChatExpireDuration, err = config.ParseDuration("2160h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[db_private_chat_expire_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("db_private_chat_expire_duration"), dAtA)
	}

	dAtA.DbGuildLogCountLimit = 100
	if pArSeR.KeyExist("db_guild_log_count_limit") {
		dAtA.DbGuildLogCountLimit = pArSeR.Uint64("db_guild_log_count_limit")
	}

	dAtA.DbMailCountLimit = 300
	if pArSeR.KeyExist("db_mail_count_limit") {
		dAtA.DbMailCountLimit = pArSeR.Uint64("db_mail_count_limit")
	}

	if pArSeR.KeyExist("extra_res_decay_coef") {
		dAtA.ExtraResDecayCoef, err = data.ParseAmount(pArSeR.String("extra_res_decay_coef"))
	} else {
		dAtA.ExtraResDecayCoef, err = data.ParseAmount("30%")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[extra_res_decay_coef] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("extra_res_decay_coef"), dAtA)
	}

	if pArSeR.KeyExist("extra_res_decay_duration") {
		dAtA.ExtraResDecayDuration, err = config.ParseDuration(pArSeR.String("extra_res_decay_duration"))
	} else {
		dAtA.ExtraResDecayDuration, err = config.ParseDuration("6m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[extra_res_decay_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("extra_res_decay_duration"), dAtA)
	}

	dAtA.CloseFightGuideDungeonId = 107
	if pArSeR.KeyExist("close_fight_guide_dungeon_id") {
		dAtA.CloseFightGuideDungeonId = pArSeR.Uint64("close_fight_guide_dungeon_id")
	}

	if pArSeR.KeyExist("refresh_recommend_hero_duration") {
		dAtA.RefreshRecommendHeroDuration, err = config.ParseDuration(pArSeR.String("refresh_recommend_hero_duration"))
	} else {
		dAtA.RefreshRecommendHeroDuration, err = config.ParseDuration("3s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[refresh_recommend_hero_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("refresh_recommend_hero_duration"), dAtA)
	}

	dAtA.RefreshRecommendHeroPageSize = 9
	if pArSeR.KeyExist("refresh_recommend_hero_page_size") {
		dAtA.RefreshRecommendHeroPageSize = pArSeR.Uint64("refresh_recommend_hero_page_size")
	}

	dAtA.RefreshRecommendHeroMinLevel = 3
	if pArSeR.KeyExist("refresh_recommend_hero_min_level") {
		dAtA.RefreshRecommendHeroMinLevel = pArSeR.Uint64("refresh_recommend_hero_min_level")
	}

	dAtA.RefreshRecommendHeroPageCount = 10
	if pArSeR.KeyExist("refresh_recommend_hero_page_count") {
		dAtA.RefreshRecommendHeroPageCount = pArSeR.Uint64("refresh_recommend_hero_page_count")
	}

	if pArSeR.KeyExist("search_hero_duration") {
		dAtA.SearchHeroDuration, err = config.ParseDuration(pArSeR.String("search_hero_duration"))
	} else {
		dAtA.SearchHeroDuration, err = config.ParseDuration("3s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[search_hero_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("search_hero_duration"), dAtA)
	}

	if pArSeR.KeyExist("recommend_hero_offline_expire_duration") {
		dAtA.RecommendHeroOfflineExpireDuration, err = config.ParseDuration(pArSeR.String("recommend_hero_offline_expire_duration"))
	} else {
		dAtA.RecommendHeroOfflineExpireDuration, err = config.ParseDuration("24h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[recommend_hero_offline_expire_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("recommend_hero_offline_expire_duration"), dAtA)
	}

	if pArSeR.KeyExist("red_packet_server_del_duration") {
		dAtA.RedPacketServerDelDuration, err = config.ParseDuration(pArSeR.String("red_packet_server_del_duration"))
	} else {
		dAtA.RedPacketServerDelDuration, err = config.ParseDuration("168h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[red_packet_server_del_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("red_packet_server_del_duration"), dAtA)
	}

	dAtA.RedPacketGuildMemberMinCount = 1
	if pArSeR.KeyExist("red_packet_guild_member_min_count") {
		dAtA.RedPacketGuildMemberMinCount = pArSeR.Uint64("red_packet_guild_member_min_count")
	}

	return dAtA, nil
}

var vAlIdAtOrMiscConfig = map[string]*config.Validator{

	"daily_reset_hour":                       config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"daily_reset_minute":                     config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"weekly_reset_hour":                      config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"weekly_reset_minute":                    config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"min_name_char_len":                      config.ParseValidator("int>0", "", false, nil, []string{"2"}),
	"max_name_char_len":                      config.ParseValidator("int>0", "", false, nil, []string{"14"}),
	"workshop_refresh_hour_minute":           config.ParseValidator("uint", "", true, nil, []string{"1200", "1600", "2000"}),
	"second_worker_cost":                     config.ParseValidator("string", "", false, nil, nil),
	"second_worker_unlock_duration":          config.ParseValidator("string", "", false, nil, nil),
	"change_hero_name_yuanbao_cost":          config.ParseValidator("uint", "", true, nil, []string{"0", "100"}),
	"change_hero_name_cost":                  config.ParseValidator("uint", "", true, nil, nil),
	"change_hero_name_duration":              config.ParseValidator("string", "", false, nil, nil),
	"first_change_hero_name_prize":           config.ParseValidator("string", "", false, nil, nil),
	"change_captain_name_cost":               config.ParseValidator("string", "", false, nil, []string{"1"}),
	"change_captain_race_duration":           config.ParseValidator("string", "", false, nil, nil),
	"change_captain_race_level":              config.ParseValidator("int>0", "", false, nil, []string{"50"}),
	"mail_min_batch_count":                   config.ParseValidator("int>0", "", false, nil, []string{"20"}),
	"mail_max_batch_count":                   config.ParseValidator("int>0", "", false, nil, []string{"50"}),
	"tower_challenge_max_times":              config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"tower_reset_challenge_duration":         config.ParseValidator("string", "", false, nil, []string{"30m"}),
	"tower_replay_count":                     config.ParseValidator("int>0", "", false, nil, []string{"4"}),
	"tower_auto_keep_floor":                  config.ParseValidator("uint", "", false, nil, []string{"5"}),
	"equipment_upgrade_multi_times":          config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"miao_building_worker_duration":          config.ParseValidator("string", "", false, nil, nil),
	"miao_building_worker_cost":              config.ParseValidator("string", "", false, nil, nil),
	"miao_tech_worker_duration":              config.ParseValidator("string", "", false, nil, nil),
	"miao_tech_worker_cost":                  config.ParseValidator("string", "", false, nil, nil),
	"miao_captain_rebirth_duration":          config.ParseValidator("string", "", false, nil, nil),
	"miao_captain_rebirth_cost":              config.ParseValidator("string", "", false, nil, nil),
	"miao_workshop_duration":                 config.ParseValidator("string", "", false, nil, nil),
	"miao_workshop_cost":                     config.ParseValidator("string", "", false, nil, nil),
	"default_forging_times":                  config.ParseValidator("uint", "", false, nil, []string{"3"}),
	"max_depot_equip_capacity":               config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"temp_depot_expire_duration":             config.ParseValidator("string", "", false, nil, []string{"24h"}),
	"max_sign_len":                           config.ParseValidator("int>0", "", false, nil, []string{"20"}),
	"max_voice_len":                          config.ParseValidator("int>0", "", false, nil, []string{"2000"}),
	"strategy_restore_duration":              config.ParseValidator("string", "", false, nil, []string{"1h"}),
	"max_resource_collect_times":             config.ParseValidator("uint", "", false, nil, []string{"10"}),
	"default_resource_collect_times":         config.ParseValidator("int", "", false, nil, []string{"5"}),
	"resource_recovery_duration":             config.ParseValidator("string", "", false, nil, []string{"1h"}),
	"max_favorite_pos_count":                 config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"flag_hero_name":                         config.ParseValidator("string", "", false, nil, []string{"[%s];%s"}),
	"country_flag_hero_name":                 config.ParseValidator("string", "", false, nil, []string{"[%s];[%s];%s"}),
	"world_chat_level":                       config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"world_chat_duration":                    config.ParseValidator("string", "", false, nil, []string{"10s"}),
	"chat_text_length":                       config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"chat_json_length":                       config.ParseValidator("int>0", "", false, nil, []string{"600"}),
	"chat_window_limit":                      config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"chat_batch_count":                       config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"broadcast_goods":                        config.ParseValidator("string", "", false, nil, nil),
	"chat_private_min_level":                 config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"chat_duration":                          config.ParseValidator("string", "", false, nil, []string{"5s"}),
	"chat_share_duration":                    config.ParseValidator("string", "", false, nil, []string{"5s"}),
	"give_add_count_down_prize_duration":     config.ParseValidator("string", "", false, nil, nil),
	"give_add_count_down_prize_times":        config.ParseValidator("int>0", "", false, nil, nil),
	"stronger_coef":                          config.ParseValidator("float64>=0", "", true, nil, []string{"0", "0.6", "1"}),
	"building_init_effect":                   config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"db_world_chat_expire_duration":          config.ParseValidator("string", "", false, nil, []string{"72h"}),
	"db_guild_chat_expire_duration":          config.ParseValidator("string", "", false, nil, []string{"168h"}),
	"db_private_chat_expire_duration":        config.ParseValidator("string", "", false, nil, []string{"2160h"}),
	"db_guild_log_count_limit":               config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"db_mail_count_limit":                    config.ParseValidator("int>0", "", false, nil, []string{"300"}),
	"extra_res_decay_coef":                   config.ParseValidator("string", "", false, nil, []string{"30%"}),
	"extra_res_decay_duration":               config.ParseValidator("string", "", false, nil, []string{"6m"}),
	"close_fight_guide_dungeon_id":           config.ParseValidator("int>0", "", false, nil, []string{"107"}),
	"refresh_recommend_hero_duration":        config.ParseValidator("string", "", false, nil, []string{"3s"}),
	"refresh_recommend_hero_page_size":       config.ParseValidator("int>0", "", false, nil, []string{"9"}),
	"refresh_recommend_hero_min_level":       config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"refresh_recommend_hero_page_count":      config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"search_hero_duration":                   config.ParseValidator("string", "", false, nil, []string{"3s"}),
	"recommend_hero_offline_expire_duration": config.ParseValidator("string", "", false, nil, []string{"24h"}),
	"red_packet_server_del_duration":         config.ParseValidator("string", "", false, nil, []string{"168h"}),
	"red_packet_guild_member_min_count":      config.ParseValidator("int>0", "", false, nil, []string{"1"}),
}

func (dAtA *MiscConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MiscConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MiscConfig) Encode() *shared_proto.MiscConfigProto {
	out := &shared_proto.MiscConfigProto{}
	out.DailyResetHour = config.U64ToI32(dAtA.DailyResetHour)
	out.DailyResetMinute = config.U64ToI32(dAtA.DailyResetMinute)
	out.WeeklyResetHour = config.U64ToI32(dAtA.WeeklyResetHour)
	out.WeeklyResetMinute = config.U64ToI32(dAtA.WeeklyResetMinute)
	out.WorkshopRefreshHourMinute = config.U64a2I32a(dAtA.WorkshopRefreshHourMinute)
	if dAtA.SecondWorkerCost != nil {
		out.SecondWorkerCost = dAtA.SecondWorkerCost.Encode()
	}
	out.SecondWorkerUnlockDuration = config.Duration2I32Seconds(dAtA.SecondWorkerUnlockDuration)
	out.ChangeHeroNameYuanbaoCost = config.U64a2I32a(dAtA.ChangeHeroNameYuanbaoCost)
	if dAtA.ChangeHeroNameCost != nil {
		out.ChangeHeroNameCost = resdata.ArrayEncodeCost(dAtA.ChangeHeroNameCost)
	}
	out.ChangeHeroNameDuration = config.Duration2I32Seconds(dAtA.ChangeHeroNameDuration)
	if dAtA.FirstChangeHeroNamePrize != nil {
		out.FirstChangeHeroNamePrize = dAtA.FirstChangeHeroNamePrize.Encode()
	}
	if dAtA.ChangeCaptainNameCost != nil {
		out.ChangeCaptainNameCost = dAtA.ChangeCaptainNameCost.Encode()
	}
	out.ChangeCaptainRaceDuration = config.Duration2I32Seconds(dAtA.ChangeCaptainRaceDuration)
	out.ChangeCaptainRaceLevel = config.U64ToI32(dAtA.ChangeCaptainRaceLevel)
	out.MailMinBatchCount = config.U64ToI32(dAtA.MailMinBatchCount)
	out.MailMaxBatchCount = config.U64ToI32(dAtA.MailMaxBatchCount)
	out.TowerChallengeMaxTimes = config.U64ToI32(dAtA.TowerChallengeMaxTimes)
	out.MiaoBuildingWorkerDuration = config.Duration2I32Seconds(dAtA.MiaoBuildingWorkerDuration)
	if dAtA.MiaoBuildingWorkerCost != nil {
		out.MiaoBuildingWorkerCost = dAtA.MiaoBuildingWorkerCost.Encode()
	}
	out.MiaoTechWorkerDuration = config.Duration2I32Seconds(dAtA.MiaoTechWorkerDuration)
	if dAtA.MiaoTechWorkerCost != nil {
		out.MiaoTechWorkerCost = dAtA.MiaoTechWorkerCost.Encode()
	}
	out.MiaoCaptainRebirthDuration = config.Duration2I32Seconds(dAtA.MiaoCaptainRebirthDuration)
	if dAtA.MiaoCaptainRebirthCost != nil {
		out.MiaoCaptainRebirthCost = dAtA.MiaoCaptainRebirthCost.Encode()
	}
	out.MiaoWorkshopDuration = config.Duration2I32Seconds(dAtA.MiaoWorkshopDuration)
	if dAtA.MiaoWorkshopCost != nil {
		out.MiaoWorkshopCost = dAtA.MiaoWorkshopCost.Encode()
	}
	out.MaxDepotEquipCapacity = config.U64ToI32(dAtA.MaxDepotEquipCapacity)
	out.MaxSignLen = config.U64ToI32(dAtA.MaxSignLen)
	out.StrategyRestoreDuration = config.Duration2I32Seconds(dAtA.StrategyRestoreDuration)
	out.MaxResourceCollectTimes = config.U64ToI32(dAtA.MaxResourceCollectTimes)
	out.ResourceRecoveryDuration = config.Duration2I32Seconds(dAtA.ResourceRecoveryDuration)
	out.MondayZeroOClock = config.I64ToI32(dAtA.MondayZeroOClock)
	out.MaxFavoritePosCount = config.U64ToI32(dAtA.MaxFavoritePosCount)
	out.WorldChatLevel = config.U64ToI32(dAtA.WorldChatLevel)
	out.WorldChatDuration = config.Duration2I32Seconds(dAtA.WorldChatDuration)
	out.ChatTextLength = config.U64ToI32(dAtA.ChatTextLength)
	if dAtA.BroadcastGoods != nil {
		out.BroadcastGoods = config.U64ToI32(dAtA.BroadcastGoods.Id)
	}
	out.ChatPrivateMinLevel = config.U64ToI32(dAtA.ChatPrivateMinLevel)
	out.ChatDuration = config.Duration2I32Seconds(dAtA.ChatDuration)
	out.ChatShareDuration = config.Duration2I32Seconds(dAtA.ChatShareDuration)
	out.StrongerCoef = config.F64a2I32aX1000(dAtA.StrongerCoef)
	out.BuildingInitEffect = config.U64ToI32(dAtA.BuildingInitEffect)
	if dAtA.ExtraResDecayCoef != nil {
		out.ExtraResDecayCoef = dAtA.ExtraResDecayCoef.Encode()
	}
	out.ExtraResDecayDuration = config.Duration2I32Seconds(dAtA.ExtraResDecayDuration)
	out.CloseFightGuideDungeonId = config.U64ToI32(dAtA.CloseFightGuideDungeonId)
	out.RefreshRecommendHeroDuration = config.Duration2I32Seconds(dAtA.RefreshRecommendHeroDuration)
	out.RefreshRecommendHeroPageSize = config.U64ToI32(dAtA.RefreshRecommendHeroPageSize)
	out.RefreshRecommendHeroMinLevel = config.U64ToI32(dAtA.RefreshRecommendHeroMinLevel)
	out.RefreshRecommendHeroPageCount = config.U64ToI32(dAtA.RefreshRecommendHeroPageCount)
	out.SearchHeroDuration = config.Duration2I32Seconds(dAtA.SearchHeroDuration)
	out.RecommendHeroOfflineExpireDuration = config.Duration2I32Seconds(dAtA.RecommendHeroOfflineExpireDuration)
	out.RedPacketGuildMemberMinCount = config.U64ToI32(dAtA.RedPacketGuildMemberMinCount)

	return out
}

func ArrayEncodeMiscConfig(datas []*MiscConfig) []*shared_proto.MiscConfigProto {

	out := make([]*shared_proto.MiscConfigProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *MiscConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.SecondWorkerCost = cOnFigS.GetCost(pArSeR.Int("second_worker_cost"))
	if dAtA.SecondWorkerCost == nil {
		return errors.Errorf("%s 配置的关联字段[second_worker_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("second_worker_cost"), *pArSeR)
	}

	intKeys = pArSeR.IntArray("change_hero_name_cost", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetCost(v)
		if obj != nil {
			dAtA.ChangeHeroNameCost = append(dAtA.ChangeHeroNameCost, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[change_hero_name_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("change_hero_name_cost"), *pArSeR)
		}
	}

	dAtA.FirstChangeHeroNamePrize = cOnFigS.GetPrize(pArSeR.Int("first_change_hero_name_prize"))
	if dAtA.FirstChangeHeroNamePrize == nil && pArSeR.Int("first_change_hero_name_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[first_change_hero_name_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_change_hero_name_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("change_captain_name_cost") {
		dAtA.ChangeCaptainNameCost = cOnFigS.GetCost(pArSeR.Int("change_captain_name_cost"))
	} else {
		dAtA.ChangeCaptainNameCost = cOnFigS.GetCost(1)
	}
	if dAtA.ChangeCaptainNameCost == nil {
		return errors.Errorf("%s 配置的关联字段[change_captain_name_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("change_captain_name_cost"), *pArSeR)
	}

	dAtA.MiaoBuildingWorkerCost = cOnFigS.GetCost(pArSeR.Int("miao_building_worker_cost"))
	if dAtA.MiaoBuildingWorkerCost == nil {
		return errors.Errorf("%s 配置的关联字段[miao_building_worker_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("miao_building_worker_cost"), *pArSeR)
	}

	dAtA.MiaoTechWorkerCost = cOnFigS.GetCost(pArSeR.Int("miao_tech_worker_cost"))
	if dAtA.MiaoTechWorkerCost == nil {
		return errors.Errorf("%s 配置的关联字段[miao_tech_worker_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("miao_tech_worker_cost"), *pArSeR)
	}

	dAtA.MiaoCaptainRebirthCost = cOnFigS.GetCost(pArSeR.Int("miao_captain_rebirth_cost"))
	if dAtA.MiaoCaptainRebirthCost == nil {
		return errors.Errorf("%s 配置的关联字段[miao_captain_rebirth_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("miao_captain_rebirth_cost"), *pArSeR)
	}

	dAtA.MiaoWorkshopCost = cOnFigS.GetCost(pArSeR.Int("miao_workshop_cost"))
	if dAtA.MiaoWorkshopCost == nil {
		return errors.Errorf("%s 配置的关联字段[miao_workshop_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("miao_workshop_cost"), *pArSeR)
	}

	dAtA.BroadcastGoods = cOnFigS.GetGoodsData(pArSeR.Uint64("broadcast_goods"))
	if dAtA.BroadcastGoods == nil {
		return errors.Errorf("%s 配置的关联字段[broadcast_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("broadcast_goods"), *pArSeR)
	}

	return nil
}

// start with MiscGenConfig ----------------------------------

func LoadMiscGenConfig(gos *config.GameObjects) (*MiscGenConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.MiscGenConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	dAtA, err := NewMiscGenConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedMiscGenConfig(gos *config.GameObjects, dAtA *MiscGenConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MiscGenConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	return dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS)
}

func NewMiscGenConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*MiscGenConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMiscGenConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MiscGenConfig{}

	dAtA.DianquanToGold = 10000
	if pArSeR.KeyExist("dianquan_to_gold") {
		dAtA.DianquanToGold = pArSeR.Uint64("dianquan_to_gold")
	}

	dAtA.DianquanToStone = 10000
	if pArSeR.KeyExist("dianquan_to_stone") {
		dAtA.DianquanToStone = pArSeR.Uint64("dianquan_to_stone")
	}

	dAtA.Stronger4Coef = []float64{0, 1, 1.5, 2}
	if pArSeR.KeyExist("stronger4_coef") {
		dAtA.Stronger4Coef = pArSeR.Float64Array("stronger4_coef", "", false)
	}

	dAtA.MiaoBaowuDuration, err = config.ParseDuration(pArSeR.String("miao_baowu_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[miao_baowu_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("miao_baowu_duration"), dAtA)
	}

	// releated field: MiaoBaowuCost
	dAtA.DailyMiaoBaowuLimit = 3
	if pArSeR.KeyExist("daily_miao_baowu_limit") {
		dAtA.DailyMiaoBaowuLimit = pArSeR.Uint64("daily_miao_baowu_limit")
	}

	dAtA.BaowuLogLimit = 40
	if pArSeR.KeyExist("baowu_log_limit") {
		dAtA.BaowuLogLimit = pArSeR.Int("baowu_log_limit")
	}

	// releated field: FirstNpcBaowu
	dAtA.FriendMaxCount = 150
	if pArSeR.KeyExist("friend_max_count") {
		dAtA.FriendMaxCount = pArSeR.Uint64("friend_max_count")
	}

	dAtA.BlackMaxCount = 150
	if pArSeR.KeyExist("black_max_count") {
		dAtA.BlackMaxCount = pArSeR.Uint64("black_max_count")
	}

	dAtA.BuySpValue = pArSeR.Uint64("buy_sp_value")
	dAtA.BuySpCost = pArSeR.Uint64("buy_sp_cost")
	dAtA.BuySpLimit = pArSeR.Uint64("buy_sp_limit")
	if pArSeR.KeyExist("sp_duration") {
		dAtA.SpDuration, err = config.ParseDuration(pArSeR.String("sp_duration"))
	} else {
		dAtA.SpDuration, err = config.ParseDuration("5m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[sp_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("sp_duration"), dAtA)
	}

	if pArSeR.KeyExist("auto_full_soldoer_duration") {
		dAtA.AutoFullSoldoerDuration, err = config.ParseDuration(pArSeR.String("auto_full_soldoer_duration"))
	} else {
		dAtA.AutoFullSoldoerDuration, err = config.ParseDuration("5m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[auto_full_soldoer_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("auto_full_soldoer_duration"), dAtA)
	}

	if pArSeR.KeyExist("tax_duration") {
		dAtA.TaxDuration, err = config.ParseDuration(pArSeR.String("tax_duration"))
	} else {
		dAtA.TaxDuration, err = config.ParseDuration("30s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tax_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tax_duration"), dAtA)
	}

	dAtA.FishMaxPoint = pArSeR.Uint64Array("fish_max_point", "", false)
	// releated field: FishPointCaptain
	// releated field: FishPointPlunder
	dAtA.DaZhaoSwitchLevelLimit = 4
	if pArSeR.KeyExist("da_zhao_switch_level_limit") {
		dAtA.DaZhaoSwitchLevelLimit = pArSeR.Uint64("da_zhao_switch_level_limit")
	}

	// releated field: UpdateOuterCityTypeCost
	dAtA.FirstHistoryChatSend = 2
	if pArSeR.KeyExist("first_history_chat_send") {
		dAtA.FirstHistoryChatSend = pArSeR.Int("first_history_chat_send")
	}

	dAtA.RandomEventNum = 729
	if pArSeR.KeyExist("random_event_num") {
		dAtA.RandomEventNum = pArSeR.Int("random_event_num")
	}

	dAtA.RandomEventOwnMinNum = 3
	if pArSeR.KeyExist("random_event_own_min_num") {
		dAtA.RandomEventOwnMinNum = pArSeR.Int("random_event_own_min_num")
	}

	if pArSeR.KeyExist("random_event_big_refresh_duration") {
		dAtA.RandomEventBigRefreshDuration, err = config.ParseDuration(pArSeR.String("random_event_big_refresh_duration"))
	} else {
		dAtA.RandomEventBigRefreshDuration, err = config.ParseDuration("6h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[random_event_big_refresh_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("random_event_big_refresh_duration"), dAtA)
	}

	if pArSeR.KeyExist("random_event_small_refresh_duration") {
		dAtA.RandomEventSmallRefreshDuration, err = config.ParseDuration(pArSeR.String("random_event_small_refresh_duration"))
	} else {
		dAtA.RandomEventSmallRefreshDuration, err = config.ParseDuration("30m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[random_event_small_refresh_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("random_event_small_refresh_duration"), dAtA)
	}

	dAtA.TargetUseStratagemLimit = 2
	if pArSeR.KeyExist("target_use_stratagem_limit") {
		dAtA.TargetUseStratagemLimit = pArSeR.Uint64("target_use_stratagem_limit")
	}

	dAtA.TrappedStratagemLimit = 5
	if pArSeR.KeyExist("trapped_stratagem_limit") {
		dAtA.TrappedStratagemLimit = pArSeR.Uint64("trapped_stratagem_limit")
	}

	// releated field: CaptainResetCost
	dAtA.SkipFightingHeroLevel = 10
	if pArSeR.KeyExist("skip_fighting_hero_level") {
		dAtA.SkipFightingHeroLevel = pArSeR.Uint64("skip_fighting_hero_level")
	}

	dAtA.SkipFightingVipLevel = 10
	if pArSeR.KeyExist("skip_fighting_vip_level") {
		dAtA.SkipFightingVipLevel = pArSeR.Uint64("skip_fighting_vip_level")
	}

	if pArSeR.KeyExist("skip_fighting_wait_duration") {
		dAtA.SkipFightingWaitDuration, err = config.ParseDuration(pArSeR.String("skip_fighting_wait_duration"))
	} else {
		dAtA.SkipFightingWaitDuration, err = config.ParseDuration("10s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[skip_fighting_wait_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("skip_fighting_wait_duration"), dAtA)
	}

	if pArSeR.KeyExist("secret_tower_cd") {
		dAtA.SecretTowerCd, err = config.ParseDuration(pArSeR.String("secret_tower_cd"))
	} else {
		dAtA.SecretTowerCd, err = config.ParseDuration("30s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[secret_tower_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("secret_tower_cd"), dAtA)
	}

	if pArSeR.KeyExist("xuanyuan_cd") {
		dAtA.XuanyuanCd, err = config.ParseDuration(pArSeR.String("xuanyuan_cd"))
	} else {
		dAtA.XuanyuanCd, err = config.ParseDuration("30s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[xuanyuan_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("xuanyuan_cd"), dAtA)
	}

	if pArSeR.KeyExist("baizhan_cd") {
		dAtA.BaizhanCd, err = config.ParseDuration(pArSeR.String("baizhan_cd"))
	} else {
		dAtA.BaizhanCd, err = config.ParseDuration("30s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[baizhan_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("baizhan_cd"), dAtA)
	}

	if pArSeR.KeyExist("hebi_cd") {
		dAtA.HebiCd, err = config.ParseDuration(pArSeR.String("hebi_cd"))
	} else {
		dAtA.HebiCd, err = config.ParseDuration("30s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[hebi_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("hebi_cd"), dAtA)
	}

	if pArSeR.KeyExist("xiongnu_cd") {
		dAtA.XiongnuCd, err = config.ParseDuration(pArSeR.String("xiongnu_cd"))
	} else {
		dAtA.XiongnuCd, err = config.ParseDuration("30s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[xiongnu_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("xiongnu_cd"), dAtA)
	}

	if pArSeR.KeyExist("mail_cd") {
		dAtA.MailCd, err = config.ParseDuration(pArSeR.String("mail_cd"))
	} else {
		dAtA.MailCd, err = config.ParseDuration("30s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[mail_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("mail_cd"), dAtA)
	}

	dAtA.SoldierPerGroup = 2000
	if pArSeR.KeyExist("soldier_per_group") {
		dAtA.SoldierPerGroup = pArSeR.Uint64("soldier_per_group")
	}

	if pArSeR.KeyExist("hero_baoz_duration") {
		dAtA.HeroBaozDuration, err = config.ParseDuration(pArSeR.String("hero_baoz_duration"))
	} else {
		dAtA.HeroBaozDuration, err = config.ParseDuration("2h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[hero_baoz_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("hero_baoz_duration"), dAtA)
	}

	dAtA.HeroBaozMaxDistance = 200
	if pArSeR.KeyExist("hero_baoz_max_distance") {
		dAtA.HeroBaozMaxDistance = pArSeR.Uint64("hero_baoz_max_distance")
	}

	if pArSeR.KeyExist("yuanbao_gift_percent") {
		dAtA.YuanbaoGiftPercent, err = data.ParseAmount(pArSeR.String("yuanbao_gift_percent"))
	} else {
		dAtA.YuanbaoGiftPercent, err = data.ParseAmount("500%")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[yuanbao_gift_percent] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("yuanbao_gift_percent"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrMiscGenConfig = map[string]*config.Validator{

	"dianquan_to_gold":                    config.ParseValidator("int>0", "", false, nil, []string{"10000"}),
	"dianquan_to_stone":                   config.ParseValidator("int>0", "", false, nil, []string{"10000"}),
	"stronger4_coef":                      config.ParseValidator("float64>=0", "", true, nil, []string{"0", "1", "1.5", "2"}),
	"miao_baowu_duration":                 config.ParseValidator("string", "", false, nil, nil),
	"miao_baowu_cost":                     config.ParseValidator("string", "", false, nil, nil),
	"daily_miao_baowu_limit":              config.ParseValidator("uint", "", false, nil, []string{"3"}),
	"baowu_log_limit":                     config.ParseValidator("int>0", "", false, nil, []string{"40"}),
	"first_npc_baowu":                     config.ParseValidator("string", "", false, nil, nil),
	"friend_max_count":                    config.ParseValidator("int>0", "", false, nil, []string{"150"}),
	"black_max_count":                     config.ParseValidator("int>0", "", false, nil, []string{"150"}),
	"buy_sp_value":                        config.ParseValidator("int>0", "", false, nil, nil),
	"buy_sp_cost":                         config.ParseValidator("int>0", "", false, nil, nil),
	"buy_sp_limit":                        config.ParseValidator("int>0", "", false, nil, nil),
	"sp_duration":                         config.ParseValidator("string", "", false, nil, []string{"5m"}),
	"auto_full_soldoer_duration":          config.ParseValidator("string", "", false, nil, []string{"5m"}),
	"tax_duration":                        config.ParseValidator("string", "", false, nil, []string{"30s"}),
	"fish_max_point":                      config.ParseValidator("uint", "", true, nil, nil),
	"fish_point_captain":                  config.ParseValidator("string", "", true, nil, nil),
	"fish_point_plunder":                  config.ParseValidator("string", "", false, nil, nil),
	"da_zhao_switch_level_limit":          config.ParseValidator("int>0", "", false, nil, []string{"4"}),
	"update_outer_city_type_cost":         config.ParseValidator("string", "", false, nil, nil),
	"first_history_chat_send":             config.ParseValidator("int>0", "", false, nil, []string{"2"}),
	"random_event_num":                    config.ParseValidator("int>0", "", false, nil, []string{"729"}),
	"random_event_own_min_num":            config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"random_event_big_refresh_duration":   config.ParseValidator("string", "", false, nil, []string{"6h"}),
	"random_event_small_refresh_duration": config.ParseValidator("string", "", false, nil, []string{"30m"}),
	"target_use_stratagem_limit":          config.ParseValidator("int>0", "", false, nil, []string{"2"}),
	"trapped_stratagem_limit":             config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"captain_reset_cost":                  config.ParseValidator("string", "", false, nil, nil),
	"skip_fighting_hero_level":            config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"skip_fighting_vip_level":             config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"skip_fighting_wait_duration":         config.ParseValidator("string", "", false, nil, []string{"10s"}),
	"secret_tower_cd":                     config.ParseValidator("string", "", false, nil, []string{"30s"}),
	"xuanyuan_cd":                         config.ParseValidator("string", "", false, nil, []string{"30s"}),
	"baizhan_cd":                          config.ParseValidator("string", "", false, nil, []string{"30s"}),
	"hebi_cd":                             config.ParseValidator("string", "", false, nil, []string{"30s"}),
	"xiongnu_cd":                          config.ParseValidator("string", "", false, nil, []string{"30s"}),
	"mail_cd":                             config.ParseValidator("string", "", false, nil, []string{"30s"}),
	"soldier_per_group":                   config.ParseValidator("int>0", "", false, nil, []string{"2000"}),
	"hero_baoz_duration":                  config.ParseValidator("string", "", false, nil, []string{"2h"}),
	"hero_baoz_max_distance":              config.ParseValidator("int>0", "", false, nil, []string{"200"}),
	"yuanbao_gift_percent":                config.ParseValidator("string", "", false, nil, []string{"500%"}),
}

func (dAtA *MiscGenConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MiscGenConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MiscGenConfig) Encode() *shared_proto.MiscGenConfigProto {
	out := &shared_proto.MiscGenConfigProto{}
	out.DianquanToGold = config.U64ToI32(dAtA.DianquanToGold)
	out.DianquanToStone = config.U64ToI32(dAtA.DianquanToStone)
	out.Stronger4Coef = config.F64a2I32aX1000(dAtA.Stronger4Coef)
	out.MiaoBaowuDuration = config.Duration2I32Seconds(dAtA.MiaoBaowuDuration)
	if dAtA.MiaoBaowuCost != nil {
		out.MiaoBaowuCost = dAtA.MiaoBaowuCost.Encode()
	}
	out.DailyMiaoBaowuLimit = config.U64ToI32(dAtA.DailyMiaoBaowuLimit)
	out.BaowuLogLimit = int32(dAtA.BaowuLogLimit)
	if dAtA.FirstNpcBaowu != nil {
		out.FirstNpcBaowu = config.U64ToI32(dAtA.FirstNpcBaowu.Id)
	}
	out.FriendMaxCount = config.U64ToI32(dAtA.FriendMaxCount)
	out.BlackMaxCount = config.U64ToI32(dAtA.BlackMaxCount)
	out.BuySpValue = config.U64ToI32(dAtA.BuySpValue)
	out.BuySpCost = config.U64ToI32(dAtA.BuySpCost)
	out.BuySpLimit = config.U64ToI32(dAtA.BuySpLimit)
	out.SpDuration = config.Duration2I32Seconds(dAtA.SpDuration)
	out.AutoFullSoldoerDuration = config.Duration2I32Seconds(dAtA.AutoFullSoldoerDuration)
	out.TaxDuration = config.Duration2I32Seconds(dAtA.TaxDuration)
	out.FishMaxPoint = config.U64a2I32a(dAtA.FishMaxPoint)
	if dAtA.FishPointCaptain != nil {
		out.FishPointCaptainSoul = config.U64a2I32a(resdata.GetResCaptainDataKeyArray(dAtA.FishPointCaptain))
	}
	out.DaZhaoSwitchLevelLimit = config.U64ToI32(dAtA.DaZhaoSwitchLevelLimit)
	if dAtA.UpdateOuterCityTypeCost != nil {
		out.UpdateOuterCityTypeCost = dAtA.UpdateOuterCityTypeCost.Encode()
	}
	out.RandomEventBigRefreshDuration = config.Duration2I32Seconds(dAtA.RandomEventBigRefreshDuration)
	out.RandomEventSmallRefreshDuration = config.Duration2I32Seconds(dAtA.RandomEventSmallRefreshDuration)
	if dAtA.CaptainResetCost != nil {
		out.CaptainResetCost = dAtA.CaptainResetCost.Encode()
	}
	out.SkipFightingHeroLevel = config.U64ToI32(dAtA.SkipFightingHeroLevel)
	out.SkipFightingVipLevel = config.U64ToI32(dAtA.SkipFightingVipLevel)
	out.SkipFightingWaitDuration = config.Duration2I32Seconds(dAtA.SkipFightingWaitDuration)
	out.SecretTowerCd = config.Duration2I32Seconds(dAtA.SecretTowerCd)
	out.XuanyuanCd = config.Duration2I32Seconds(dAtA.XuanyuanCd)
	out.BaizhanCd = config.Duration2I32Seconds(dAtA.BaizhanCd)
	out.HebiCd = config.Duration2I32Seconds(dAtA.HebiCd)
	out.XiongnuCd = config.Duration2I32Seconds(dAtA.XiongnuCd)
	out.MailCd = config.Duration2I32Seconds(dAtA.MailCd)
	out.SoldierPerGroup = config.U64ToI32(dAtA.SoldierPerGroup)
	out.HeroBaozDuration = config.Duration2I32Seconds(dAtA.HeroBaozDuration)
	out.HeroBaozMaxDistance = config.U64ToI32(dAtA.HeroBaozMaxDistance)
	if dAtA.YuanbaoGiftPercent != nil {
		out.YuanbaoGiftPercent = dAtA.YuanbaoGiftPercent.Encode()
	}

	return out
}

func ArrayEncodeMiscGenConfig(datas []*MiscGenConfig) []*shared_proto.MiscGenConfigProto {

	out := make([]*shared_proto.MiscGenConfigProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *MiscGenConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.MiaoBaowuCost = cOnFigS.GetCost(pArSeR.Int("miao_baowu_cost"))
	if dAtA.MiaoBaowuCost == nil {
		return errors.Errorf("%s 配置的关联字段[miao_baowu_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("miao_baowu_cost"), *pArSeR)
	}

	dAtA.FirstNpcBaowu = cOnFigS.GetBaowuData(pArSeR.Uint64("first_npc_baowu"))
	if dAtA.FirstNpcBaowu == nil && pArSeR.Uint64("first_npc_baowu") != 0 {
		return errors.Errorf("%s 配置的关联字段[first_npc_baowu] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_npc_baowu"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("fish_point_captain", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetResCaptainData(v)
		if obj != nil {
			dAtA.FishPointCaptain = append(dAtA.FishPointCaptain, obj)
		} else if v != 0 {
			return errors.Errorf("%s 配置的关联字段[fish_point_captain] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fish_point_captain"), *pArSeR)
		}
	}

	dAtA.FishPointPlunder = cOnFigS.GetPlunder(pArSeR.Uint64("fish_point_plunder"))
	if dAtA.FishPointPlunder == nil {
		return errors.Errorf("%s 配置的关联字段[fish_point_plunder] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fish_point_plunder"), *pArSeR)
	}

	dAtA.UpdateOuterCityTypeCost = cOnFigS.GetCost(pArSeR.Int("update_outer_city_type_cost"))
	if dAtA.UpdateOuterCityTypeCost == nil {
		return errors.Errorf("%s 配置的关联字段[update_outer_city_type_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("update_outer_city_type_cost"), *pArSeR)
	}

	dAtA.CaptainResetCost = cOnFigS.GetCost(pArSeR.Int("captain_reset_cost"))
	if dAtA.CaptainResetCost == nil {
		return errors.Errorf("%s 配置的关联字段[captain_reset_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain_reset_cost"), *pArSeR)
	}

	return nil
}

// start with RegionConfig ----------------------------------

func LoadRegionConfig(gos *config.GameObjects) (*RegionConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.RegionConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	dAtA, err := NewRegionConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedRegionConfig(gos *config.GameObjects, dAtA *RegionConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RegionConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	return dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS)
}

func NewRegionConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*RegionConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRegionConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RegionConfig{}

	if pArSeR.KeyExist("slow_move_duration") {
		dAtA.SlowMoveDuration, err = time.ParseDuration(pArSeR.String("slow_move_duration"))
	} else {
		dAtA.SlowMoveDuration, err = time.ParseDuration("3h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[slow_move_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("slow_move_duration"), dAtA)
	}

	if pArSeR.KeyExist("expel_duration") {
		dAtA.ExpelDuration, err = time.ParseDuration(pArSeR.String("expel_duration"))
	} else {
		dAtA.ExpelDuration, err = time.ParseDuration("180s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[expel_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("expel_duration"), dAtA)
	}

	dAtA.BasicTroopMoveVelocityPerSecond = 0.25
	if pArSeR.KeyExist("basic_troop_move_velocity_per_second") {
		dAtA.BasicTroopMoveVelocityPerSecond = pArSeR.Float64("basic_troop_move_velocity_per_second")
	}

	dAtA.BasicTroopMoveToNpcVelocityPerSecond = 0.25
	if pArSeR.KeyExist("basic_troop_move_to_npc_velocity_per_second") {
		dAtA.BasicTroopMoveToNpcVelocityPerSecond = pArSeR.Float64("basic_troop_move_to_npc_velocity_per_second")
	}

	dAtA.Edge = 1
	if pArSeR.KeyExist("edge") {
		dAtA.Edge = pArSeR.Float64("edge")
	}

	if pArSeR.KeyExist("troop_move_offset_duration") {
		dAtA.TroopMoveOffsetDuration, err = config.ParseDuration(pArSeR.String("troop_move_offset_duration"))
	} else {
		dAtA.TroopMoveOffsetDuration, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[troop_move_offset_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("troop_move_offset_duration"), dAtA)
	}

	dAtA.EdgeNotHomeLen = 6
	if pArSeR.KeyExist("edge_not_home_len") {
		dAtA.EdgeNotHomeLen = pArSeR.Uint64("edge_not_home_len")
	}

	dAtA.AttackerWoundedRate = 0.5
	if pArSeR.KeyExist("attacker_wounded_rate") {
		dAtA.AttackerWoundedRate = pArSeR.Float64("attacker_wounded_rate")
	}

	dAtA.DefenserWoundedRate = 0.7
	if pArSeR.KeyExist("defenser_wounded_rate") {
		dAtA.DefenserWoundedRate = pArSeR.Float64("defenser_wounded_rate")
	}

	dAtA.AssisterWoundedRate = 0.5
	if pArSeR.KeyExist("assister_wounded_rate") {
		dAtA.AssisterWoundedRate = pArSeR.Float64("assister_wounded_rate")
	}

	dAtA.RobberCoef = 1
	if pArSeR.KeyExist("robber_coef") {
		dAtA.RobberCoef = pArSeR.Float64("robber_coef")
	}

	dAtA.MaxLostProsperity = 8000
	if pArSeR.KeyExist("max_lost_prosperity") {
		dAtA.MaxLostProsperity = pArSeR.Uint64("max_lost_prosperity")
	}

	dAtA.LostProsperityCoef = 0.3
	if pArSeR.KeyExist("lost_prosperity_coef") {
		dAtA.LostProsperityCoef = pArSeR.Float64("lost_prosperity_coef")
	}

	if pArSeR.KeyExist("rob_tick_duration") {
		dAtA.RobTickDuration, err = time.ParseDuration(pArSeR.String("rob_tick_duration"))
	} else {
		dAtA.RobTickDuration, err = time.ParseDuration("6s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[rob_tick_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("rob_tick_duration"), dAtA)
	}

	if pArSeR.KeyExist("rob_max_duration") {
		dAtA.RobMaxDuration, err = time.ParseDuration(pArSeR.String("rob_max_duration"))
	} else {
		dAtA.RobMaxDuration, err = time.ParseDuration("30m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[rob_max_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("rob_max_duration"), dAtA)
	}

	if pArSeR.KeyExist("reduce_prosperity_duration") {
		dAtA.ReduceProsperityDuration, err = time.ParseDuration(pArSeR.String("reduce_prosperity_duration"))
	} else {
		dAtA.ReduceProsperityDuration, err = time.ParseDuration("10m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[reduce_prosperity_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("reduce_prosperity_duration"), dAtA)
	}

	if pArSeR.KeyExist("assister_tick_duration") {
		dAtA.AssisterTickDuration, err = time.ParseDuration(pArSeR.String("assister_tick_duration"))
	} else {
		dAtA.AssisterTickDuration, err = time.ParseDuration("5m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[assister_tick_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("assister_tick_duration"), dAtA)
	}

	if pArSeR.KeyExist("assister_max_duration") {
		dAtA.AssisterMaxDuration, err = time.ParseDuration(pArSeR.String("assister_max_duration"))
	} else {
		dAtA.AssisterMaxDuration, err = time.ParseDuration("24h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[assister_max_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("assister_max_duration"), dAtA)
	}

	if pArSeR.KeyExist("rob_baowu_tick_duration") {
		dAtA.RobBaowuTickDuration, err = time.ParseDuration(pArSeR.String("rob_baowu_tick_duration"))
	} else {
		dAtA.RobBaowuTickDuration, err = time.ParseDuration("10m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[rob_baowu_tick_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("rob_baowu_tick_duration"), dAtA)
	}

	dAtA.MaxRobbers = 4
	if pArSeR.KeyExist("max_robbers") {
		dAtA.MaxRobbers = pArSeR.Uint64("max_robbers")
	}

	dAtA.MaxAssist = 4
	if pArSeR.KeyExist("max_assist") {
		dAtA.MaxAssist = pArSeR.Uint64("max_assist")
	}

	dAtA.MaxInvationTroops = 3
	if pArSeR.KeyExist("max_invation_troops") {
		dAtA.MaxInvationTroops = pArSeR.Uint64("max_invation_troops")
	}

	dAtA.MaxInvationCaptain = 5
	if pArSeR.KeyExist("max_invation_captain") {
		dAtA.MaxInvationCaptain = pArSeR.Uint64("max_invation_captain")
	}

	if pArSeR.KeyExist("tent_building_duration") {
		dAtA.TentBuildingDuration, err = config.ParseDuration(pArSeR.String("tent_building_duration"))
	} else {
		dAtA.TentBuildingDuration, err = config.ParseDuration("0m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tent_building_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tent_building_duration"), dAtA)
	}

	if pArSeR.KeyExist("tent_restore_duration") {
		dAtA.TentRestoreDuration, err = config.ParseDuration(pArSeR.String("tent_restore_duration"))
	} else {
		dAtA.TentRestoreDuration, err = config.ParseDuration("100m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tent_restore_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tent_restore_duration"), dAtA)
	}

	if pArSeR.KeyExist("tent_free_move_cooldown") {
		dAtA.TentFreeMoveCooldown, err = config.ParseDuration(pArSeR.String("tent_free_move_cooldown"))
	} else {
		dAtA.TentFreeMoveCooldown, err = config.ParseDuration("5m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tent_free_move_cooldown] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tent_free_move_cooldown"), dAtA)
	}

	if pArSeR.KeyExist("tent_home_region_enter_duration") {
		dAtA.TentHomeRegionEnterDuration, err = config.ParseDuration(pArSeR.String("tent_home_region_enter_duration"))
	} else {
		dAtA.TentHomeRegionEnterDuration, err = config.ParseDuration("6h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tent_home_region_enter_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tent_home_region_enter_duration"), dAtA)
	}

	if pArSeR.KeyExist("tent_monster_region_enter_duration") {
		dAtA.TentMonsterRegionEnterDuration, err = config.ParseDuration(pArSeR.String("tent_monster_region_enter_duration"))
	} else {
		dAtA.TentMonsterRegionEnterDuration, err = config.ParseDuration("6h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tent_monster_region_enter_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tent_monster_region_enter_duration"), dAtA)
	}

	dAtA.MultiLevelNpcInitTimes = 3
	if pArSeR.KeyExist("multi_level_npc_init_times") {
		dAtA.MultiLevelNpcInitTimes = pArSeR.Uint64("multi_level_npc_init_times")
	}

	dAtA.MultiLevelNpcMaxTimes = 5
	if pArSeR.KeyExist("multi_level_npc_max_times") {
		dAtA.MultiLevelNpcMaxTimes = pArSeR.Uint64("multi_level_npc_max_times")
	}

	if pArSeR.KeyExist("multi_level_npc_recovery_duration") {
		dAtA.MultiLevelNpcRecoveryDuration, err = config.ParseDuration(pArSeR.String("multi_level_npc_recovery_duration"))
	} else {
		dAtA.MultiLevelNpcRecoveryDuration, err = config.ParseDuration("3h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[multi_level_npc_recovery_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("multi_level_npc_recovery_duration"), dAtA)
	}

	dAtA.NpcRobLostProsperityPerDuration = 1
	if pArSeR.KeyExist("npc_rob_lost_prosperity_per_duration") {
		dAtA.NpcRobLostProsperityPerDuration = pArSeR.Uint64("npc_rob_lost_prosperity_per_duration")
	}

	dAtA.MiaoTentBuildingDuration, err = config.ParseDuration(pArSeR.String("miao_tent_building_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[miao_tent_building_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("miao_tent_building_duration"), dAtA)
	}

	// releated field: MiaoTentBuildingCost
	// releated field: CombatScene
	if pArSeR.KeyExist("ast_defend_restore_home_prosperity_amount") {
		dAtA.AstDefendRestoreHomeProsperityAmount, err = data.ParseAmount(pArSeR.String("ast_defend_restore_home_prosperity_amount"))
	} else {
		dAtA.AstDefendRestoreHomeProsperityAmount, err = data.ParseAmount("8")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[ast_defend_restore_home_prosperity_amount] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("ast_defend_restore_home_prosperity_amount"), dAtA)
	}

	dAtA.AstDefendRestoreHomeProsperity = 8
	if pArSeR.KeyExist("ast_defend_restore_home_prosperity") {
		dAtA.AstDefendRestoreHomeProsperity = pArSeR.Uint64("ast_defend_restore_home_prosperity")
	}

	dAtA.RestoreHomeProsperity = 10
	if pArSeR.KeyExist("restore_home_prosperity") {
		dAtA.RestoreHomeProsperity = pArSeR.Uint64("restore_home_prosperity")
	}

	if pArSeR.KeyExist("restore_home_prosperity_duration") {
		dAtA.RestoreHomeProsperityDuration, err = config.ParseDuration(pArSeR.String("restore_home_prosperity_duration"))
	} else {
		dAtA.RestoreHomeProsperityDuration, err = config.ParseDuration("1m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[restore_home_prosperity_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("restore_home_prosperity_duration"), dAtA)
	}

	dAtA.AstDefendLogLimit = 100
	if pArSeR.KeyExist("ast_defend_log_limit") {
		dAtA.AstDefendLogLimit = pArSeR.Uint64("ast_defend_log_limit")
	}

	dAtA.RestoreTentProsperity = 10
	if pArSeR.KeyExist("restore_tent_prosperity") {
		dAtA.RestoreTentProsperity = pArSeR.Uint64("restore_tent_prosperity")
	}

	if pArSeR.KeyExist("restore_tent_prosperity_duration") {
		dAtA.RestoreTentProsperityDuration, err = config.ParseDuration(pArSeR.String("restore_tent_prosperity_duration"))
	} else {
		dAtA.RestoreTentProsperityDuration, err = config.ParseDuration("1m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[restore_tent_prosperity_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("restore_tent_prosperity_duration"), dAtA)
	}

	dAtA.LossTentProsperity = 10
	if pArSeR.KeyExist("loss_tent_prosperity") {
		dAtA.LossTentProsperity = pArSeR.Uint64("loss_tent_prosperity")
	}

	if pArSeR.KeyExist("loss_tent_prosperity_duration") {
		dAtA.LossTentProsperityDuration, err = config.ParseDuration(pArSeR.String("loss_tent_prosperity_duration"))
	} else {
		dAtA.LossTentProsperityDuration, err = config.ParseDuration("1m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[loss_tent_prosperity_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("loss_tent_prosperity_duration"), dAtA)
	}

	if pArSeR.KeyExist("white_flag_duration") {
		dAtA.WhiteFlagDuration, err = config.ParseDuration(pArSeR.String("white_flag_duration"))
	} else {
		dAtA.WhiteFlagDuration, err = config.ParseDuration("5m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[white_flag_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("white_flag_duration"), dAtA)
	}

	if pArSeR.KeyExist("ruins_base_expire_duration") {
		dAtA.RuinsBaseExpireDuration, err = config.ParseDuration(pArSeR.String("ruins_base_expire_duration"))
	} else {
		dAtA.RuinsBaseExpireDuration, err = config.ParseDuration("24h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[ruins_base_expire_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("ruins_base_expire_duration"), dAtA)
	}

	dAtA.GuildRegionCenterX = 60
	if pArSeR.KeyExist("guild_region_center_x") {
		dAtA.GuildRegionCenterX = pArSeR.Uint64("guild_region_center_x")
	}

	dAtA.GuildRegionCenterY = 52
	if pArSeR.KeyExist("guild_region_center_y") {
		dAtA.GuildRegionCenterY = pArSeR.Uint64("guild_region_center_y")
	}

	dAtA.GuildRegionRadius = []uint64{10, 30, 50}
	if pArSeR.KeyExist("guild_region_radius") {
		dAtA.GuildRegionRadius = pArSeR.Uint64Array("guild_region_radius", "", false)
	}

	if pArSeR.KeyExist("investigate_cd") {
		dAtA.InvestigateCd, err = config.ParseDuration(pArSeR.String("investigate_cd"))
	} else {
		dAtA.InvestigateCd, err = config.ParseDuration("5m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[investigate_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("investigate_cd"), dAtA)
	}

	// releated field: MiaoInvestigateCdCost
	// releated field: InvestigateCost
	dAtA.InvestigateSpeedup = 7
	if pArSeR.KeyExist("investigate_speedup") {
		dAtA.InvestigateSpeedup = pArSeR.Float64("investigate_speedup")
	}

	if pArSeR.KeyExist("investigate_mail_timeout") {
		dAtA.InvestigateMailTimeout, err = config.ParseDuration(pArSeR.String("investigate_mail_timeout"))
	} else {
		dAtA.InvestigateMailTimeout, err = config.ParseDuration("30m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[investigate_mail_timeout] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("investigate_mail_timeout"), dAtA)
	}

	dAtA.InvestigateMaxDistance = 300
	if pArSeR.KeyExist("investigate_max_distance") {
		dAtA.InvestigateMaxDistance = pArSeR.Uint64("investigate_max_distance")
	}

	dAtA.InvestigationLimit = 10
	if pArSeR.KeyExist("investigation_limit") {
		dAtA.InvestigationLimit = pArSeR.Uint64("investigation_limit")
	}

	if pArSeR.KeyExist("investigation_expire_duration") {
		dAtA.InvestigationExpireDuration, err = config.ParseDuration(pArSeR.String("investigation_expire_duration"))
	} else {
		dAtA.InvestigationExpireDuration, err = config.ParseDuration("30m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[investigation_expire_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("investigation_expire_duration"), dAtA)
	}

	dAtA.InvestigationBaowuCount = 3
	if pArSeR.KeyExist("investigation_baowu_count") {
		dAtA.InvestigationBaowuCount = pArSeR.Uint64("investigation_baowu_count")
	}

	if pArSeR.KeyExist("new_hero_mian_duration") {
		dAtA.NewHeroMianDuration, err = config.ParseDuration(pArSeR.String("new_hero_mian_duration"))
	} else {
		dAtA.NewHeroMianDuration, err = config.ParseDuration("72h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[new_hero_mian_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("new_hero_mian_duration"), dAtA)
	}

	dAtA.NewHeroRemoveMianBaseLevel = 6
	if pArSeR.KeyExist("new_hero_remove_mian_base_level") {
		dAtA.NewHeroRemoveMianBaseLevel = pArSeR.Uint64("new_hero_remove_mian_base_level")
	}

	if pArSeR.KeyExist("reborn_mian_duration") {
		dAtA.RebornMianDuration, err = config.ParseDuration(pArSeR.String("reborn_mian_duration"))
	} else {
		dAtA.RebornMianDuration, err = config.ParseDuration("24h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[reborn_mian_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("reborn_mian_duration"), dAtA)
	}

	dAtA.MinViewXLen = 40
	if pArSeR.KeyExist("min_view_xlen") {
		dAtA.MinViewXLen = pArSeR.Int("min_view_xlen")
	}

	dAtA.MinViewYLen = 30
	if pArSeR.KeyExist("min_view_ylen") {
		dAtA.MinViewYLen = pArSeR.Int("min_view_ylen")
	}

	dAtA.MaxViewXLen = 80
	if pArSeR.KeyExist("max_view_xlen") {
		dAtA.MaxViewXLen = pArSeR.Int("max_view_xlen")
	}

	dAtA.MaxViewYLen = 60
	if pArSeR.KeyExist("max_view_ylen") {
		dAtA.MaxViewYLen = pArSeR.Int("max_view_ylen")
	}

	dAtA.ListEnemyPosCount = 30
	if pArSeR.KeyExist("list_enemy_pos_count") {
		dAtA.ListEnemyPosCount = pArSeR.Int("list_enemy_pos_count")
	}

	dAtA.SearchBaozNpcCount = 30
	if pArSeR.KeyExist("search_baoz_npc_count") {
		dAtA.SearchBaozNpcCount = pArSeR.Int("search_baoz_npc_count")
	}

	dAtA.KeepBaozMaxDistance = 300
	if pArSeR.KeyExist("keep_baoz_max_distance") {
		dAtA.KeepBaozMaxDistance = pArSeR.Int("keep_baoz_max_distance")
	}

	if pArSeR.KeyExist("keep_baoz_duration") {
		dAtA.KeepBaozDuration, err = config.ParseDuration(pArSeR.String("keep_baoz_duration"))
	} else {
		dAtA.KeepBaozDuration, err = config.ParseDuration("12h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[keep_baoz_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("keep_baoz_duration"), dAtA)
	}

	// skip field: InvestigateTroopId

	// calculate fields
	dAtA.MinTroopMoveVelocityPerSecond = dAtA.BasicTroopMoveVelocityPerSecond * 0.5

	return dAtA, nil
}

var vAlIdAtOrRegionConfig = map[string]*config.Validator{

	"slow_move_duration":                          config.ParseValidator("string", "", false, nil, []string{"3h"}),
	"expel_duration":                              config.ParseValidator("string", "", false, nil, []string{"180s"}),
	"basic_troop_move_velocity_per_second":        config.ParseValidator("float64>0", "", false, nil, []string{"0.25"}),
	"basic_troop_move_to_npc_velocity_per_second": config.ParseValidator("float64>0", "", false, nil, []string{"0.25"}),
	"edge": config.ParseValidator("float64>0", "", false, nil, []string{"1"}),
	"troop_move_offset_duration":                config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"edge_not_home_len":                         config.ParseValidator("uint", "", false, nil, []string{"6"}),
	"attacker_wounded_rate":                     config.ParseValidator("float64>0", "", false, nil, []string{"0.5"}),
	"defenser_wounded_rate":                     config.ParseValidator("float64>0", "", false, nil, []string{"0.7"}),
	"assister_wounded_rate":                     config.ParseValidator("float64>0", "", false, nil, []string{"0.5"}),
	"robber_coef":                               config.ParseValidator("float64>0", "", false, nil, []string{"1"}),
	"max_lost_prosperity":                       config.ParseValidator("int>0", "", false, nil, []string{"8000"}),
	"lost_prosperity_coef":                      config.ParseValidator("float64>0", "", false, nil, []string{"0.3"}),
	"rob_tick_duration":                         config.ParseValidator("string", "", false, nil, []string{"6s"}),
	"rob_max_duration":                          config.ParseValidator("string", "", false, nil, []string{"30m"}),
	"reduce_prosperity_duration":                config.ParseValidator("string", "", false, nil, []string{"10m"}),
	"assister_tick_duration":                    config.ParseValidator("string", "", false, nil, []string{"5m"}),
	"assister_max_duration":                     config.ParseValidator("string", "", false, nil, []string{"24h"}),
	"rob_baowu_tick_duration":                   config.ParseValidator("string", "", false, nil, []string{"10m"}),
	"max_robbers":                               config.ParseValidator("int>0", "", false, nil, []string{"4"}),
	"max_assist":                                config.ParseValidator("int>0", "", false, nil, []string{"4"}),
	"max_invation_troops":                       config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"max_invation_captain":                      config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"tent_building_duration":                    config.ParseValidator("string", "", false, nil, []string{"0m"}),
	"tent_restore_duration":                     config.ParseValidator("string", "", false, nil, []string{"100m"}),
	"tent_free_move_cooldown":                   config.ParseValidator("string", "", false, nil, []string{"5m"}),
	"tent_home_region_enter_duration":           config.ParseValidator("string", "", false, nil, []string{"6h"}),
	"tent_monster_region_enter_duration":        config.ParseValidator("string", "", false, nil, []string{"6h"}),
	"multi_level_npc_init_times":                config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"multi_level_npc_max_times":                 config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"multi_level_npc_recovery_duration":         config.ParseValidator("string", "", false, nil, []string{"3h"}),
	"npc_rob_lost_prosperity_per_duration":      config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"miao_tent_building_duration":               config.ParseValidator("string", "", false, nil, nil),
	"miao_tent_building_cost":                   config.ParseValidator("string", "", false, nil, nil),
	"combat_scene":                              config.ParseValidator("string", "", false, nil, []string{"CombatScene"}),
	"ast_defend_restore_home_prosperity_amount": config.ParseValidator("string", "", false, nil, []string{"8"}),
	"ast_defend_restore_home_prosperity":        config.ParseValidator("int>0", "", false, nil, []string{"8"}),
	"restore_home_prosperity":                   config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"restore_home_prosperity_duration":          config.ParseValidator("string", "", false, nil, []string{"1m"}),
	"ast_defend_log_limit":                      config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"restore_tent_prosperity":                   config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"restore_tent_prosperity_duration":          config.ParseValidator("string", "", false, nil, []string{"1m"}),
	"loss_tent_prosperity":                      config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"loss_tent_prosperity_duration":             config.ParseValidator("string", "", false, nil, []string{"1m"}),
	"white_flag_duration":                       config.ParseValidator("string", "", false, nil, []string{"5m"}),
	"ruins_base_expire_duration":                config.ParseValidator("string", "", false, nil, []string{"24h"}),
	"guild_region_center_x":                     config.ParseValidator("int>0", "", false, nil, []string{"60"}),
	"guild_region_center_y":                     config.ParseValidator("int>0", "", false, nil, []string{"52"}),
	"guild_region_radius":                       config.ParseValidator("uint", "", true, nil, []string{"10", "30", "50"}),
	"investigate_cd":                            config.ParseValidator("string", "", false, nil, []string{"5m"}),
	"miao_investigate_cd_cost":                  config.ParseValidator("string", "", false, nil, nil),
	"investigate_cost":                          config.ParseValidator("string", "", false, nil, nil),
	"investigate_speedup":                       config.ParseValidator("float64>0", "", false, nil, []string{"7"}),
	"investigate_mail_timeout":                  config.ParseValidator("string", "", false, nil, []string{"30m"}),
	"investigate_max_distance":                  config.ParseValidator("int>0", "", false, nil, []string{"300"}),
	"investigation_limit":                       config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"investigation_expire_duration":             config.ParseValidator("string", "", false, nil, []string{"30m"}),
	"investigation_baowu_count":                 config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"new_hero_mian_duration":                    config.ParseValidator("string", "", false, nil, []string{"72h"}),
	"new_hero_remove_mian_base_level":           config.ParseValidator("int>0", "", false, nil, []string{"6"}),
	"reborn_mian_duration":                      config.ParseValidator("string", "", false, nil, []string{"24h"}),
	"min_view_xlen":                             config.ParseValidator("int>0", "", false, nil, []string{"40"}),
	"min_view_ylen":                             config.ParseValidator("int>0", "", false, nil, []string{"30"}),
	"max_view_xlen":                             config.ParseValidator("int>0", "", false, nil, []string{"80"}),
	"max_view_ylen":                             config.ParseValidator("int>0", "", false, nil, []string{"60"}),
	"list_enemy_pos_count":                      config.ParseValidator("int>0", "", false, nil, []string{"30"}),
	"search_baoz_npc_count":                     config.ParseValidator("int>0", "", false, nil, []string{"30"}),
	"keep_baoz_max_distance":                    config.ParseValidator("int>0", "", false, nil, []string{"300"}),
	"keep_baoz_duration":                        config.ParseValidator("string", "", false, nil, []string{"12h"}),
}

func (dAtA *RegionConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *RegionConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *RegionConfig) Encode() *shared_proto.RegionConfigProto {
	out := &shared_proto.RegionConfigProto{}
	out.BasicTroopMoveVelocityPerSecond = config.F64ToI32X1000(dAtA.BasicTroopMoveVelocityPerSecond)
	out.Edge = config.F64ToI32X1000(dAtA.Edge)
	out.EdgeNotHomeLen = config.U64ToI32(dAtA.EdgeNotHomeLen)
	out.TentBuildingDuration = config.Duration2I32Seconds(dAtA.TentBuildingDuration)
	out.MultiLevelNpcMaxTimes = config.U64ToI32(dAtA.MultiLevelNpcMaxTimes)
	out.MultiLevelNpcRecoveryDuration = config.Duration2I32Seconds(dAtA.MultiLevelNpcRecoveryDuration)
	out.MiaoTentBuildingDuration = config.Duration2I32Seconds(dAtA.MiaoTentBuildingDuration)
	if dAtA.MiaoTentBuildingCost != nil {
		out.MiaoTentBuildingCost = dAtA.MiaoTentBuildingCost.Encode()
	}
	if dAtA.AstDefendRestoreHomeProsperityAmount != nil {
		out.AstDefendRestoreHomeProsperityAmount = dAtA.AstDefendRestoreHomeProsperityAmount.Encode()
	}
	out.AstDefendRestoreHomeProsperity = config.U64ToI32(dAtA.AstDefendRestoreHomeProsperity)
	out.RestoreHomeProsperity = config.U64ToI32(dAtA.RestoreHomeProsperity)
	out.RestoreHomeProsperityDuration = config.Duration2I32Seconds(dAtA.RestoreHomeProsperityDuration)
	out.WhiteFlagDuration = config.Duration2I32Seconds(dAtA.WhiteFlagDuration)
	out.GuildRegionCenterX = config.U64ToI32(dAtA.GuildRegionCenterX)
	out.GuildRegionCenterY = config.U64ToI32(dAtA.GuildRegionCenterY)
	out.InvestigateCd = config.Duration2I32Seconds(dAtA.InvestigateCd)
	if dAtA.MiaoInvestigateCdCost != nil {
		out.MiaoInvestigateCdCost = dAtA.MiaoInvestigateCdCost.Encode()
	}
	if dAtA.InvestigateCost != nil {
		out.InvestigateCost = dAtA.InvestigateCost.Encode()
	}
	out.InvestigateSpeedup = config.F64ToI32X1000(dAtA.InvestigateSpeedup)
	out.InvestigateMailTimeout = config.Duration2I32Seconds(dAtA.InvestigateMailTimeout)
	out.InvestigateMaxDistance = config.U64ToI32(dAtA.InvestigateMaxDistance)
	out.NewHeroRemoveMianBaseLevel = config.U64ToI32(dAtA.NewHeroRemoveMianBaseLevel)
	out.RebornMianDuration = config.Duration2I32Seconds(dAtA.RebornMianDuration)

	return out
}

func ArrayEncodeRegionConfig(datas []*RegionConfig) []*shared_proto.RegionConfigProto {

	out := make([]*shared_proto.RegionConfigProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *RegionConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.MiaoTentBuildingCost = cOnFigS.GetCost(pArSeR.Int("miao_tent_building_cost"))
	if dAtA.MiaoTentBuildingCost == nil {
		return errors.Errorf("%s 配置的关联字段[miao_tent_building_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("miao_tent_building_cost"), *pArSeR)
	}

	if pArSeR.KeyExist("combat_scene") {
		dAtA.CombatScene = cOnFigS.GetCombatScene(pArSeR.String("combat_scene"))
	} else {
		dAtA.CombatScene = cOnFigS.GetCombatScene("CombatScene")
	}
	if dAtA.CombatScene == nil {
		return errors.Errorf("%s 配置的关联字段[combat_scene] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("combat_scene"), *pArSeR)
	}

	dAtA.MiaoInvestigateCdCost = cOnFigS.GetCost(pArSeR.Int("miao_investigate_cd_cost"))
	if dAtA.MiaoInvestigateCdCost == nil {
		return errors.Errorf("%s 配置的关联字段[miao_investigate_cd_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("miao_investigate_cd_cost"), *pArSeR)
	}

	dAtA.InvestigateCost = cOnFigS.GetCost(pArSeR.Int("investigate_cost"))
	if dAtA.InvestigateCost == nil {
		return errors.Errorf("%s 配置的关联字段[investigate_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("investigate_cost"), *pArSeR)
	}

	return nil
}

// start with RegionGenConfig ----------------------------------

func LoadRegionGenConfig(gos *config.GameObjects) (*RegionGenConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.RegionGenConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	dAtA, err := NewRegionGenConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedRegionGenConfig(gos *config.GameObjects, dAtA *RegionGenConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RegionGenConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	return dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS)
}

func NewRegionGenConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*RegionGenConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRegionGenConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RegionGenConfig{}

	if pArSeR.KeyExist("use_goods_mian_max_duraion") {
		dAtA.UseGoodsMianMaxDuraion, err = config.ParseDuration(pArSeR.String("use_goods_mian_max_duraion"))
	} else {
		dAtA.UseGoodsMianMaxDuraion, err = config.ParseDuration("30m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[use_goods_mian_max_duraion] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("use_goods_mian_max_duraion"), dAtA)
	}

	dAtA.InvaseHeroInitTimes = 3
	if pArSeR.KeyExist("invase_hero_init_times") {
		dAtA.InvaseHeroInitTimes = pArSeR.Uint64("invase_hero_init_times")
	}

	dAtA.InvaseHeroMaxTimes = 5
	if pArSeR.KeyExist("invase_hero_max_times") {
		dAtA.InvaseHeroMaxTimes = pArSeR.Uint64("invase_hero_max_times")
	}

	if pArSeR.KeyExist("invase_hero_recovery_duration") {
		dAtA.InvaseHeroRecoveryDuration, err = config.ParseDuration(pArSeR.String("invase_hero_recovery_duration"))
	} else {
		dAtA.InvaseHeroRecoveryDuration, err = config.ParseDuration("3h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[invase_hero_recovery_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("invase_hero_recovery_duration"), dAtA)
	}

	dAtA.JunTuanNpcInitTimes = 3
	if pArSeR.KeyExist("jun_tuan_npc_init_times") {
		dAtA.JunTuanNpcInitTimes = pArSeR.Uint64("jun_tuan_npc_init_times")
	}

	dAtA.JunTuanNpcMaxTimes = 5
	if pArSeR.KeyExist("jun_tuan_npc_max_times") {
		dAtA.JunTuanNpcMaxTimes = pArSeR.Uint64("jun_tuan_npc_max_times")
	}

	if pArSeR.KeyExist("jun_tuan_npc_recovery_duration") {
		dAtA.JunTuanNpcRecoveryDuration, err = config.ParseDuration(pArSeR.String("jun_tuan_npc_recovery_duration"))
	} else {
		dAtA.JunTuanNpcRecoveryDuration, err = config.ParseDuration("3h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[jun_tuan_npc_recovery_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("jun_tuan_npc_recovery_duration"), dAtA)
	}

	dAtA.JunTuanWinTimeLimit = 3
	if pArSeR.KeyExist("jun_tuan_win_time_limit") {
		dAtA.JunTuanWinTimeLimit = pArSeR.Uint64("jun_tuan_win_time_limit")
	}

	return dAtA, nil
}

var vAlIdAtOrRegionGenConfig = map[string]*config.Validator{

	"use_goods_mian_max_duraion":     config.ParseValidator("string", "", false, nil, []string{"30m"}),
	"invase_hero_init_times":         config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"invase_hero_max_times":          config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"invase_hero_recovery_duration":  config.ParseValidator("string", "", false, nil, []string{"3h"}),
	"jun_tuan_npc_init_times":        config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"jun_tuan_npc_max_times":         config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"jun_tuan_npc_recovery_duration": config.ParseValidator("string", "", false, nil, []string{"3h"}),
	"jun_tuan_win_time_limit":        config.ParseValidator("int>0", "", false, nil, []string{"3"}),
}

func (dAtA *RegionGenConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *RegionGenConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *RegionGenConfig) Encode() *shared_proto.RegionGenConfigProto {
	out := &shared_proto.RegionGenConfigProto{}
	out.UseGoodsMianMaxDuraion = config.Duration2I32Seconds(dAtA.UseGoodsMianMaxDuraion)
	out.InvaseHeroMaxTimes = config.U64ToI32(dAtA.InvaseHeroMaxTimes)
	out.InvaseHeroRecoveryDuration = config.Duration2I32Seconds(dAtA.InvaseHeroRecoveryDuration)
	out.JunTuanNpcMaxTimes = config.U64ToI32(dAtA.JunTuanNpcMaxTimes)
	out.JunTuanNpcRecoveryDuration = config.Duration2I32Seconds(dAtA.JunTuanNpcRecoveryDuration)
	out.JunTuanWinTimeLimit = config.U64ToI32(dAtA.JunTuanWinTimeLimit)

	return out
}

func ArrayEncodeRegionGenConfig(datas []*RegionGenConfig) []*shared_proto.RegionGenConfigProto {

	out := make([]*shared_proto.RegionGenConfigProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *RegionGenConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	return nil
}

type related_configs interface {
	GetBaowuData(uint64) *resdata.BaowuData
	GetCombatScene(string) *scene.CombatScene
	GetCost(int) *resdata.Cost
	GetCountryData(uint64) *country.CountryData
	GetGoodsData(uint64) *goods.GoodsData
	GetNpcBaseData(uint64) *basedata.NpcBaseData
	GetPlunder(uint64) *resdata.Plunder
	GetPrize(int) *resdata.Prize
	GetResCaptainData(uint64) *resdata.ResCaptainData
}
