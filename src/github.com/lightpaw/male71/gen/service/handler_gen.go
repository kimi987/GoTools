package service

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lightpaw/male7/gen/iface"
	pb51 "github.com/lightpaw/male7/gen/pb/activity"
	pb24 "github.com/lightpaw/male7/gen/pb/bai_zhan"
	pb13 "github.com/lightpaw/male7/gen/pb/chat"
	pb21 "github.com/lightpaw/male7/gen/pb/client_config"
	pb27 "github.com/lightpaw/male7/gen/pb/country"
	pb11 "github.com/lightpaw/male7/gen/pb/depot"
	pb39 "github.com/lightpaw/male7/gen/pb/dianquan"
	pb2 "github.com/lightpaw/male7/gen/pb/domestic"
	pb26 "github.com/lightpaw/male7/gen/pb/dungeon"
	pb12 "github.com/lightpaw/male7/gen/pb/equipment"
	pb38 "github.com/lightpaw/male7/gen/pb/farm"
	pb16 "github.com/lightpaw/male7/gen/pb/fishing"
	pb31 "github.com/lightpaw/male7/gen/pb/garden"
	pb19 "github.com/lightpaw/male7/gen/pb/gem"
	pb3 "github.com/lightpaw/male7/gen/pb/gm"
	pb9 "github.com/lightpaw/male7/gen/pb/guild"
	pb41 "github.com/lightpaw/male7/gen/pb/hebi"
	pb1 "github.com/lightpaw/male7/gen/pb/login"
	pb8 "github.com/lightpaw/male7/gen/pb/mail"
	pb4 "github.com/lightpaw/male7/gen/pb/military"
	pb42 "github.com/lightpaw/male7/gen/pb/mingc"
	pb44 "github.com/lightpaw/male7/gen/pb/mingc_war"
	pb5 "github.com/lightpaw/male7/gen/pb/misc"
	pb43 "github.com/lightpaw/male7/gen/pb/promotion"
	pb34 "github.com/lightpaw/male7/gen/pb/question"
	pb45 "github.com/lightpaw/male7/gen/pb/random_event"
	pb23 "github.com/lightpaw/male7/gen/pb/rank"
	pb49 "github.com/lightpaw/male7/gen/pb/red_packet"
	pb7 "github.com/lightpaw/male7/gen/pb/region"
	pb35 "github.com/lightpaw/male7/gen/pb/relation"
	pb22 "github.com/lightpaw/male7/gen/pb/secret_tower"
	pb20 "github.com/lightpaw/male7/gen/pb/shop"
	pb46 "github.com/lightpaw/male7/gen/pb/strategy"
	pb10 "github.com/lightpaw/male7/gen/pb/stress"
	pb37 "github.com/lightpaw/male7/gen/pb/survey"
	pb29 "github.com/lightpaw/male7/gen/pb/tag"
	pb15 "github.com/lightpaw/male7/gen/pb/task"
	pb50 "github.com/lightpaw/male7/gen/pb/teach"
	pb14 "github.com/lightpaw/male7/gen/pb/tower"
	pb48 "github.com/lightpaw/male7/gen/pb/vip"
	pb36 "github.com/lightpaw/male7/gen/pb/xiongnu"
	pb40 "github.com/lightpaw/male7/gen/pb/xuanyuan"
	pb33 "github.com/lightpaw/male7/gen/pb/zhanjiang"
	pb32 "github.com/lightpaw/male7/gen/pb/zhengwu"
	"github.com/pkg/errors"
)

// 包含消息的信息
type MsgData struct {
	ModuleID   int
	SequenceID int

	// 处理这条消息需要的真实proto. 例如 *MountUpgradeProto
	Proto interface{}
}

// 在网络线程解析收到的消息
func Unmarshal(moduleID, sequenceID int, data []byte) (*MsgData, error) {
	switch moduleID {
	case 1: // login
		return unmarshal_login(sequenceID, data)

	case 2: // domestic
		return unmarshal_domestic(sequenceID, data)

	case 3: // gm
		return unmarshal_gm(sequenceID, data)

	case 4: // military
		return unmarshal_military(sequenceID, data)

	case 5: // misc
		return unmarshal_misc(sequenceID, data)

	case 7: // region
		return unmarshal_region(sequenceID, data)

	case 8: // mail
		return unmarshal_mail(sequenceID, data)

	case 9: // guild
		return unmarshal_guild(sequenceID, data)

	case 10: // stress
		return unmarshal_stress(sequenceID, data)

	case 11: // depot
		return unmarshal_depot(sequenceID, data)

	case 12: // equipment
		return unmarshal_equipment(sequenceID, data)

	case 13: // chat
		return unmarshal_chat(sequenceID, data)

	case 14: // tower
		return unmarshal_tower(sequenceID, data)

	case 15: // task
		return unmarshal_task(sequenceID, data)

	case 16: // fishing
		return unmarshal_fishing(sequenceID, data)

	case 19: // gem
		return unmarshal_gem(sequenceID, data)

	case 20: // shop
		return unmarshal_shop(sequenceID, data)

	case 21: // client_config
		return unmarshal_client_config(sequenceID, data)

	case 22: // secret_tower
		return unmarshal_secret_tower(sequenceID, data)

	case 23: // rank
		return unmarshal_rank(sequenceID, data)

	case 24: // bai_zhan
		return unmarshal_bai_zhan(sequenceID, data)

	case 26: // dungeon
		return unmarshal_dungeon(sequenceID, data)

	case 27: // country
		return unmarshal_country(sequenceID, data)

	case 29: // tag
		return unmarshal_tag(sequenceID, data)

	case 31: // garden
		return unmarshal_garden(sequenceID, data)

	case 32: // zhengwu
		return unmarshal_zhengwu(sequenceID, data)

	case 33: // zhanjiang
		return unmarshal_zhanjiang(sequenceID, data)

	case 34: // question
		return unmarshal_question(sequenceID, data)

	case 35: // relation
		return unmarshal_relation(sequenceID, data)

	case 36: // xiongnu
		return unmarshal_xiongnu(sequenceID, data)

	case 37: // survey
		return unmarshal_survey(sequenceID, data)

	case 38: // farm
		return unmarshal_farm(sequenceID, data)

	case 39: // dianquan
		return unmarshal_dianquan(sequenceID, data)

	case 40: // xuanyuan
		return unmarshal_xuanyuan(sequenceID, data)

	case 41: // hebi
		return unmarshal_hebi(sequenceID, data)

	case 42: // mingc
		return unmarshal_mingc(sequenceID, data)

	case 43: // promotion
		return unmarshal_promotion(sequenceID, data)

	case 44: // mingc_war
		return unmarshal_mingc_war(sequenceID, data)

	case 45: // random_event
		return unmarshal_random_event(sequenceID, data)

	case 46: // strategy
		return unmarshal_strategy(sequenceID, data)

	case 48: // vip
		return unmarshal_vip(sequenceID, data)

	case 49: // red_packet
		return unmarshal_red_packet(sequenceID, data)

	case 50: // teach
		return unmarshal_teach(sequenceID, data)

	case 51: // activity
		return unmarshal_activity(sequenceID, data)

	default:
		return nil, errors.New(fmt.Sprintf("收到未知消息: %d.%d", moduleID, sequenceID))
	}
}

func unmarshal_login(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_internal_login
		p := &pb1.C2SInternalLoginProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "login.c2s_internal_login UnmarshalMsgProto &C2SInternalLoginProto fail")
		}

		return &MsgData{ModuleID: 1, SequenceID: 1, Proto: p}, nil

	case 7: //c2s_login
		p := &pb1.C2SLoginProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "login.c2s_login UnmarshalMsgProto &C2SLoginProto fail")
		}

		return &MsgData{ModuleID: 1, SequenceID: 7, Proto: p}, nil

	case 18: //c2s_set_tutorial_progress
		p := &pb1.C2SSetTutorialProgressProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "login.c2s_set_tutorial_progress UnmarshalMsgProto &C2SSetTutorialProgressProto fail")
		}

		return &MsgData{ModuleID: 1, SequenceID: 18, Proto: p}, nil

	case 3: //c2s_create_hero
		p := &pb1.C2SCreateHeroProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "login.c2s_create_hero UnmarshalMsgProto &C2SCreateHeroProto fail")
		}

		return &MsgData{ModuleID: 1, SequenceID: 3, Proto: p}, nil

	case 10: //c2s_loaded
		return &MsgData{ModuleID: 1, SequenceID: 10}, nil

	case 13: //c2s_robot_login
		p := &pb1.C2SRobotLoginProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "login.c2s_robot_login UnmarshalMsgProto &C2SRobotLoginProto fail")
		}

		return &MsgData{ModuleID: 1, SequenceID: 13, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("login收到未知消息: %d", sequenceID))
	}
}

func unmarshal_domestic(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_create_building
		p := &pb2.C2SCreateBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_create_building UnmarshalMsgProto &C2SCreateBuildingProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_upgrade_building
		p := &pb2.C2SUpgradeBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_upgrade_building UnmarshalMsgProto &C2SUpgradeBuildingProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 4, Proto: p}, nil

	case 7: //c2s_rebuild_resource_building
		p := &pb2.C2SRebuildResourceBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_rebuild_resource_building UnmarshalMsgProto &C2SRebuildResourceBuildingProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 7, Proto: p}, nil

	case 108: //c2s_unlock_outer_city
		p := &pb2.C2SUnlockOuterCityProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_unlock_outer_city UnmarshalMsgProto &C2SUnlockOuterCityProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 108, Proto: p}, nil

	case 142: //c2s_update_outer_city_type
		p := &pb2.C2SUpdateOuterCityTypeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_update_outer_city_type UnmarshalMsgProto &C2SUpdateOuterCityTypeProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 142, Proto: p}, nil

	case 111: //c2s_upgrade_outer_city_building
		p := &pb2.C2SUpgradeOuterCityBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_upgrade_outer_city_building UnmarshalMsgProto &C2SUpgradeOuterCityBuildingProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 111, Proto: p}, nil

	case 15: //c2s_collect_resource
		p := &pb2.C2SCollectResourceProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_collect_resource UnmarshalMsgProto &C2SCollectResourceProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 15, Proto: p}, nil

	case 76: //c2s_collect_resource_v2
		p := &pb2.C2SCollectResourceV2Proto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_collect_resource_v2 UnmarshalMsgProto &C2SCollectResourceV2Proto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 76, Proto: p}, nil

	case 81: //c2s_request_resource_conflict
		return &MsgData{ModuleID: 2, SequenceID: 81}, nil

	case 18: //c2s_learn_technology
		p := &pb2.C2SLearnTechnologyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_learn_technology UnmarshalMsgProto &C2SLearnTechnologyProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 18, Proto: p}, nil

	case 87: //c2s_unlock_stable_building
		p := &pb2.C2SUnlockStableBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_unlock_stable_building UnmarshalMsgProto &C2SUnlockStableBuildingProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 87, Proto: p}, nil

	case 24: //c2s_upgrade_stable_building
		p := &pb2.C2SUpgradeStableBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_upgrade_stable_building UnmarshalMsgProto &C2SUpgradeStableBuildingProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 24, Proto: p}, nil

	case 90: //c2s_is_hero_name_exist
		p := &pb2.C2SIsHeroNameExistProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_is_hero_name_exist UnmarshalMsgProto &C2SIsHeroNameExistProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 90, Proto: p}, nil

	case 30: //c2s_change_hero_name
		p := &pb2.C2SChangeHeroNameProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_change_hero_name UnmarshalMsgProto &C2SChangeHeroNameProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 30, Proto: p}, nil

	case 33: //c2s_list_old_name
		p := &pb2.C2SListOldNameProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_list_old_name UnmarshalMsgProto &C2SListOldNameProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 33, Proto: p}, nil

	case 35: //c2s_view_other_hero
		p := &pb2.C2SViewOtherHeroProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_view_other_hero UnmarshalMsgProto &C2SViewOtherHeroProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 35, Proto: p}, nil

	case 125: //c2s_view_fight_info
		p := &pb2.C2SViewFightInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_view_fight_info UnmarshalMsgProto &C2SViewFightInfoProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 125, Proto: p}, nil

	case 41: //c2s_miao_building_worker_cd
		p := &pb2.C2SMiaoBuildingWorkerCdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_miao_building_worker_cd UnmarshalMsgProto &C2SMiaoBuildingWorkerCdProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 41, Proto: p}, nil

	case 44: //c2s_miao_tech_worker_cd
		p := &pb2.C2SMiaoTechWorkerCdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_miao_tech_worker_cd UnmarshalMsgProto &C2SMiaoTechWorkerCdProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 44, Proto: p}, nil

	case 51: //c2s_forging_equip
		p := &pb2.C2SForgingEquipProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_forging_equip UnmarshalMsgProto &C2SForgingEquipProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 51, Proto: p}, nil

	case 66: //c2s_sign
		p := &pb2.C2SSignProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_sign UnmarshalMsgProto &C2SSignProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 66, Proto: p}, nil

	case 69: //c2s_voice
		p := &pb2.C2SVoiceProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_voice UnmarshalMsgProto &C2SVoiceProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 69, Proto: p}, nil

	case 60: //c2s_request_city_exchange_event
		return &MsgData{ModuleID: 2, SequenceID: 60}, nil

	case 63: //c2s_city_event_exchange
		p := &pb2.C2SCityEventExchangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_city_event_exchange UnmarshalMsgProto &C2SCityEventExchangeProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 63, Proto: p}, nil

	case 94: //c2s_change_head
		p := &pb2.C2SChangeHeadProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_change_head UnmarshalMsgProto &C2SChangeHeadProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 94, Proto: p}, nil

	case 130: //c2s_change_body
		p := &pb2.C2SChangeBodyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_change_body UnmarshalMsgProto &C2SChangeBodyProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 130, Proto: p}, nil

	case 114: //c2s_collect_countdown_prize
		return &MsgData{ModuleID: 2, SequenceID: 114}, nil

	case 119: //c2s_start_workshop
		p := &pb2.C2SStartWorkshopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_start_workshop UnmarshalMsgProto &C2SStartWorkshopProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 119, Proto: p}, nil

	case 122: //c2s_collect_workshop
		p := &pb2.C2SCollectWorkshopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_collect_workshop UnmarshalMsgProto &C2SCollectWorkshopProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 122, Proto: p}, nil

	case 127: //c2s_workshop_miao_cd
		p := &pb2.C2SWorkshopMiaoCdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_workshop_miao_cd UnmarshalMsgProto &C2SWorkshopMiaoCdProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 127, Proto: p}, nil

	case 133: //c2s_refresh_workshop
		return &MsgData{ModuleID: 2, SequenceID: 133}, nil

	case 136: //c2s_collect_season_prize
		return &MsgData{ModuleID: 2, SequenceID: 136}, nil

	case 147: //c2s_buy_sp
		p := &pb2.C2SBuySpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_buy_sp UnmarshalMsgProto &C2SBuySpProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 147, Proto: p}, nil

	case 150: //c2s_use_buf_effect
		p := &pb2.C2SUseBufEffectProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_use_buf_effect UnmarshalMsgProto &C2SUseBufEffectProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 150, Proto: p}, nil

	case 154: //c2s_open_buf_effect_ui
		return &MsgData{ModuleID: 2, SequenceID: 154}, nil

	case 158: //c2s_use_advantage
		p := &pb2.C2SUseAdvantageProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_use_advantage UnmarshalMsgProto &C2SUseAdvantageProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 158, Proto: p}, nil

	case 166: //c2s_worker_unlock
		p := &pb2.C2SWorkerUnlockProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.c2s_worker_unlock UnmarshalMsgProto &C2SWorkerUnlockProto fail")
		}

		return &MsgData{ModuleID: 2, SequenceID: 166, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("domestic收到未知消息: %d", sequenceID))
	}
}

func unmarshal_gm(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 5: //c2s_list_cmd
		return &MsgData{ModuleID: 3, SequenceID: 5}, nil

	case 1: //c2s_gm
		p := &pb3.C2SGmProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gm.c2s_gm UnmarshalMsgProto &C2SGmProto fail")
		}

		return &MsgData{ModuleID: 3, SequenceID: 1, Proto: p}, nil

	case 7: //c2s_invase_target_id
		p := &pb3.C2SInvaseTargetIdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gm.c2s_invase_target_id UnmarshalMsgProto &C2SInvaseTargetIdProto fail")
		}

		return &MsgData{ModuleID: 3, SequenceID: 7, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("gm收到未知消息: %d", sequenceID))
	}
}

func unmarshal_military(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 2: //c2s_recruit_soldier
		p := &pb4.C2SRecruitSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_recruit_soldier UnmarshalMsgProto &C2SRecruitSoldierProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 2, Proto: p}, nil

	case 120: //c2s_recruit_soldier_v2
		p := &pb4.C2SRecruitSoldierV2Proto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_recruit_soldier_v2 UnmarshalMsgProto &C2SRecruitSoldierV2Proto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 120, Proto: p}, nil

	case 6: //c2s_heal_wounded_soldier
		p := &pb4.C2SHealWoundedSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_heal_wounded_soldier UnmarshalMsgProto &C2SHealWoundedSoldierProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 6, Proto: p}, nil

	case 9: //c2s_captain_change_soldier
		p := &pb4.C2SCaptainChangeSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_change_soldier UnmarshalMsgProto &C2SCaptainChangeSoldierProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 9, Proto: p}, nil

	case 66: //c2s_captain_full_soldier
		p := &pb4.C2SCaptainFullSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_full_soldier UnmarshalMsgProto &C2SCaptainFullSoldierProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 66, Proto: p}, nil

	case 149: //c2s_force_add_soldier
		return &MsgData{ModuleID: 4, SequenceID: 149}, nil

	case 12: //c2s_fight
		p := &pb4.C2SFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_fight UnmarshalMsgProto &C2SFightProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 12, Proto: p}, nil

	case 101: //c2s_multi_fight
		return &MsgData{ModuleID: 4, SequenceID: 101}, nil

	case 198: //c2s_fightx
		p := &pb4.C2SFightxProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_fightx UnmarshalMsgProto &C2SFightxProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 198, Proto: p}, nil

	case 15: //c2s_upgrade_soldier_level
		return &MsgData{ModuleID: 4, SequenceID: 15}, nil

	case 109: //c2s_recruit_captain_v2
		return &MsgData{ModuleID: 4, SequenceID: 109}, nil

	case 176: //c2s_random_captain_head
		return &MsgData{ModuleID: 4, SequenceID: 176}, nil

	case 146: //c2s_recruit_captain_seeker
		p := &pb4.C2SRecruitCaptainSeekerProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_recruit_captain_seeker UnmarshalMsgProto &C2SRecruitCaptainSeekerProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 146, Proto: p}, nil

	case 106: //c2s_set_defense_troop
		p := &pb4.C2SSetDefenseTroopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_set_defense_troop UnmarshalMsgProto &C2SSetDefenseTroopProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 106, Proto: p}, nil

	case 129: //c2s_clear_defense_troop_defeated_mail
		p := &pb4.C2SClearDefenseTroopDefeatedMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_clear_defense_troop_defeated_mail UnmarshalMsgProto &C2SClearDefenseTroopDefeatedMailProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 129, Proto: p}, nil

	case 188: //c2s_set_defenser_auto_full_soldier
		p := &pb4.C2SSetDefenserAutoFullSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_set_defenser_auto_full_soldier UnmarshalMsgProto &C2SSetDefenserAutoFullSoldierProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 188, Proto: p}, nil

	case 193: //c2s_use_copy_defenser_goods
		p := &pb4.C2SUseCopyDefenserGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_use_copy_defenser_goods UnmarshalMsgProto &C2SUseCopyDefenserGoodsProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 193, Proto: p}, nil

	case 34: //c2s_sell_seek_captain
		p := &pb4.C2SSellSeekCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_sell_seek_captain UnmarshalMsgProto &C2SSellSeekCaptainProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 34, Proto: p}, nil

	case 45: //c2s_set_multi_captain_index
		p := &pb4.C2SSetMultiCaptainIndexProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_set_multi_captain_index UnmarshalMsgProto &C2SSetMultiCaptainIndexProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 45, Proto: p}, nil

	case 143: //c2s_set_pve_captain
		p := &pb4.C2SSetPveCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_set_pve_captain UnmarshalMsgProto &C2SSetPveCaptainProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 143, Proto: p}, nil

	case 38: //c2s_fire_captain
		p := &pb4.C2SFireCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_fire_captain UnmarshalMsgProto &C2SFireCaptainProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 38, Proto: p}, nil

	case 48: //c2s_captain_refined
		p := &pb4.C2SCaptainRefinedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_refined UnmarshalMsgProto &C2SCaptainRefinedProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 48, Proto: p}, nil

	case 206: //c2s_captain_enhance
		p := &pb4.C2SCaptainEnhanceProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_enhance UnmarshalMsgProto &C2SCaptainEnhanceProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 206, Proto: p}, nil

	case 82: //c2s_change_captain_name
		p := &pb4.C2SChangeCaptainNameProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_change_captain_name UnmarshalMsgProto &C2SChangeCaptainNameProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 82, Proto: p}, nil

	case 85: //c2s_change_captain_race
		p := &pb4.C2SChangeCaptainRaceProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_change_captain_race UnmarshalMsgProto &C2SChangeCaptainRaceProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 85, Proto: p}, nil

	case 89: //c2s_captain_rebirth_preview
		p := &pb4.C2SCaptainRebirthPreviewProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_rebirth_preview UnmarshalMsgProto &C2SCaptainRebirthPreviewProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 89, Proto: p}, nil

	case 92: //c2s_captain_rebirth
		p := &pb4.C2SCaptainRebirthProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_rebirth UnmarshalMsgProto &C2SCaptainRebirthProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 92, Proto: p}, nil

	case 210: //c2s_captain_progress
		p := &pb4.C2SCaptainProgressProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_progress UnmarshalMsgProto &C2SCaptainProgressProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 210, Proto: p}, nil

	case 166: //c2s_captain_rebirth_miao_cd
		p := &pb4.C2SCaptainRebirthMiaoCdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_rebirth_miao_cd UnmarshalMsgProto &C2SCaptainRebirthMiaoCdProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 166, Proto: p}, nil

	case 136: //c2s_collect_captain_training_exp
		return &MsgData{ModuleID: 4, SequenceID: 136}, nil

	case 213: //c2s_captain_train_exp
		return &MsgData{ModuleID: 4, SequenceID: 213}, nil

	case 261: //c2s_captain_can_collect_exp
		return &MsgData{ModuleID: 4, SequenceID: 261}, nil

	case 140: //c2s_use_training_exp_goods
		p := &pb4.C2SUseTrainingExpGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_use_training_exp_goods UnmarshalMsgProto &C2SUseTrainingExpGoodsProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 140, Proto: p}, nil

	case 216: //c2s_use_level_exp_goods
		p := &pb4.C2SUseLevelExpGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_use_level_exp_goods UnmarshalMsgProto &C2SUseLevelExpGoodsProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 216, Proto: p}, nil

	case 243: //c2s_use_level_exp_goods2
		p := &pb4.C2SUseLevelExpGoods2Proto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_use_level_exp_goods2 UnmarshalMsgProto &C2SUseLevelExpGoods2Proto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 243, Proto: p}, nil

	case 255: //c2s_auto_use_goods_until_captain_levelup
		p := &pb4.C2SAutoUseGoodsUntilCaptainLevelupProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_auto_use_goods_until_captain_levelup UnmarshalMsgProto &C2SAutoUseGoodsUntilCaptainLevelupProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 255, Proto: p}, nil

	case 74: //c2s_get_max_recruit_soldier
		return &MsgData{ModuleID: 4, SequenceID: 74}, nil

	case 76: //c2s_get_max_heal_soldier
		return &MsgData{ModuleID: 4, SequenceID: 76}, nil

	case 112: //c2s_jiu_guan_consult
		return &MsgData{ModuleID: 4, SequenceID: 112}, nil

	case 117: //c2s_jiu_guan_refresh
		p := &pb4.C2SJiuGuanRefreshProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_jiu_guan_refresh UnmarshalMsgProto &C2SJiuGuanRefreshProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 117, Proto: p}, nil

	case 125: //c2s_unlock_captain_restraint_spell
		p := &pb4.C2SUnlockCaptainRestraintSpellProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_unlock_captain_restraint_spell UnmarshalMsgProto &C2SUnlockCaptainRestraintSpellProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 125, Proto: p}, nil

	case 132: //c2s_get_captain_stat_details
		p := &pb4.C2SGetCaptainStatDetailsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_get_captain_stat_details UnmarshalMsgProto &C2SGetCaptainStatDetailsProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 132, Proto: p}, nil

	case 219: //c2s_captain_stat_details
		p := &pb4.C2SCaptainStatDetailsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_stat_details UnmarshalMsgProto &C2SCaptainStatDetailsProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 219, Proto: p}, nil

	case 169: //c2s_update_captain_official
		p := &pb4.C2SUpdateCaptainOfficialProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_update_captain_official UnmarshalMsgProto &C2SUpdateCaptainOfficialProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 169, Proto: p}, nil

	case 222: //c2s_set_captain_official
		p := &pb4.C2SSetCaptainOfficialProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_set_captain_official UnmarshalMsgProto &C2SSetCaptainOfficialProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 222, Proto: p}, nil

	case 172: //c2s_leave_captain_official
		p := &pb4.C2SLeaveCaptainOfficialProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_leave_captain_official UnmarshalMsgProto &C2SLeaveCaptainOfficialProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 172, Proto: p}, nil

	case 185: //c2s_use_gong_xun_goods
		p := &pb4.C2SUseGongXunGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_use_gong_xun_goods UnmarshalMsgProto &C2SUseGongXunGoodsProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 185, Proto: p}, nil

	case 228: //c2s_use_gongxun_goods
		p := &pb4.C2SUseGongxunGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_use_gongxun_goods UnmarshalMsgProto &C2SUseGongxunGoodsProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 228, Proto: p}, nil

	case 181: //c2s_close_fight_guide
		p := &pb4.C2SCloseFightGuideProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_close_fight_guide UnmarshalMsgProto &C2SCloseFightGuideProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 181, Proto: p}, nil

	case 190: //c2s_view_other_hero_captain
		p := &pb4.C2SViewOtherHeroCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_view_other_hero_captain UnmarshalMsgProto &C2SViewOtherHeroCaptainProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 190, Proto: p}, nil

	case 231: //c2s_captain_born
		p := &pb4.C2SCaptainBornProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_born UnmarshalMsgProto &C2SCaptainBornProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 231, Proto: p}, nil

	case 234: //c2s_captain_upstar
		p := &pb4.C2SCaptainUpstarProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_upstar UnmarshalMsgProto &C2SCaptainUpstarProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 234, Proto: p}, nil

	case 268: //c2s_captain_exchange
		p := &pb4.C2SCaptainExchangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_captain_exchange UnmarshalMsgProto &C2SCaptainExchangeProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 268, Proto: p}, nil

	case 252: //c2s_notice_captain_has_viewed
		return &MsgData{ModuleID: 4, SequenceID: 252}, nil

	case 265: //c2s_activate_captain_friendship
		p := &pb4.C2SActivateCaptainFriendshipProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_activate_captain_friendship UnmarshalMsgProto &C2SActivateCaptainFriendshipProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 265, Proto: p}, nil

	case 272: //c2s_notice_official_has_viewed
		p := &pb4.C2SNoticeOfficialHasViewedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.c2s_notice_official_has_viewed UnmarshalMsgProto &C2SNoticeOfficialHasViewedProto fail")
		}

		return &MsgData{ModuleID: 4, SequenceID: 272, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("military收到未知消息: %d", sequenceID))
	}
}

func unmarshal_misc(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_heart_beat
		return &MsgData{ModuleID: 5, SequenceID: 1}, nil

	case 35: //c2s_background_heart_beat
		return &MsgData{ModuleID: 5, SequenceID: 35}, nil

	case 36: //c2s_background_weakup
		return &MsgData{ModuleID: 5, SequenceID: 36}, nil

	case 3: //c2s_config
		p := &pb5.C2SConfigProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_config UnmarshalMsgProto &C2SConfigProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 3, Proto: p}, nil

	case 76: //c2s_configlua
		p := &pb5.C2SConfigluaProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_configlua UnmarshalMsgProto &C2SConfigluaProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 76, Proto: p}, nil

	case 7: //c2s_client_log
		p := &pb5.C2SClientLogProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_client_log UnmarshalMsgProto &C2SClientLogProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 7, Proto: p}, nil

	case 8: //c2s_sync_time
		p := &pb5.C2SSyncTimeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_sync_time UnmarshalMsgProto &C2SSyncTimeProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 8, Proto: p}, nil

	case 10: //c2s_block
		return &MsgData{ModuleID: 5, SequenceID: 10}, nil

	case 15: //c2s_ping
		return &MsgData{ModuleID: 5, SequenceID: 15}, nil

	case 17: //c2s_client_version
		p := &pb5.C2SClientVersionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_client_version UnmarshalMsgProto &C2SClientVersionProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 17, Proto: p}, nil

	case 19: //c2s_update_pf_token
		p := &pb5.C2SUpdatePfTokenProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_update_pf_token UnmarshalMsgProto &C2SUpdatePfTokenProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 19, Proto: p}, nil

	case 21: //c2s_settings
		p := &pb5.C2SSettingsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_settings UnmarshalMsgProto &C2SSettingsProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 21, Proto: p}, nil

	case 24: //c2s_settings_to_default
		return &MsgData{ModuleID: 5, SequenceID: 24}, nil

	case 31: //c2s_update_location
		p := &pb5.C2SUpdateLocationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_update_location UnmarshalMsgProto &C2SUpdateLocationProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 31, Proto: p}, nil

	case 42: //c2s_collect_charge_prize
		p := &pb5.C2SCollectChargePrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_collect_charge_prize UnmarshalMsgProto &C2SCollectChargePrizeProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 42, Proto: p}, nil

	case 46: //c2s_collect_daily_bargain
		p := &pb5.C2SCollectDailyBargainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_collect_daily_bargain UnmarshalMsgProto &C2SCollectDailyBargainProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 46, Proto: p}, nil

	case 51: //c2s_activate_duration_card
		p := &pb5.C2SActivateDurationCardProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_activate_duration_card UnmarshalMsgProto &C2SActivateDurationCardProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 51, Proto: p}, nil

	case 54: //c2s_collect_duration_card_daily_prize
		p := &pb5.C2SCollectDurationCardDailyPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_collect_duration_card_daily_prize UnmarshalMsgProto &C2SCollectDurationCardDailyPrizeProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 54, Proto: p}, nil

	case 64: //c2s_set_privacy_setting
		p := &pb5.C2SSetPrivacySettingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_set_privacy_setting UnmarshalMsgProto &C2SSetPrivacySettingProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 64, Proto: p}, nil

	case 67: //c2s_set_default_privacy_settings
		return &MsgData{ModuleID: 5, SequenceID: 67}, nil

	case 69: //c2s_get_product_info
		p := &pb5.C2SGetProductInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.c2s_get_product_info UnmarshalMsgProto &C2SGetProductInfoProto fail")
		}

		return &MsgData{ModuleID: 5, SequenceID: 69, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("misc收到未知消息: %d", sequenceID))
	}
}

func unmarshal_region(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 148: //c2s_update_self_view
		p := &pb7.C2SUpdateSelfViewProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_update_self_view UnmarshalMsgProto &C2SUpdateSelfViewProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 148, Proto: p}, nil

	case 150: //c2s_close_view
		return &MsgData{ModuleID: 7, SequenceID: 150}, nil

	case 94: //c2s_pre_invasion_target
		p := &pb7.C2SPreInvasionTargetProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_pre_invasion_target UnmarshalMsgProto &C2SPreInvasionTargetProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 94, Proto: p}, nil

	case 161: //c2s_watch_base_unit
		p := &pb7.C2SWatchBaseUnitProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_watch_base_unit UnmarshalMsgProto &C2SWatchBaseUnitProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 161, Proto: p}, nil

	case 156: //c2s_request_troop_unit
		p := &pb7.C2SRequestTroopUnitProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_request_troop_unit UnmarshalMsgProto &C2SRequestTroopUnitProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 156, Proto: p}, nil

	case 131: //c2s_request_ruins_base
		p := &pb7.C2SRequestRuinsBaseProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_request_ruins_base UnmarshalMsgProto &C2SRequestRuinsBaseProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 131, Proto: p}, nil

	case 91: //c2s_use_mian_goods
		p := &pb7.C2SUseMianGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_use_mian_goods UnmarshalMsgProto &C2SUseMianGoodsProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 91, Proto: p}, nil

	case 40: //c2s_upgrade_base
		return &MsgData{ModuleID: 7, SequenceID: 40}, nil

	case 86: //c2s_white_flag_detail
		p := &pb7.C2SWhiteFlagDetailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_white_flag_detail UnmarshalMsgProto &C2SWhiteFlagDetailProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 86, Proto: p}, nil

	case 211: //c2s_get_buy_prosperity_cost
		return &MsgData{ModuleID: 7, SequenceID: 211}, nil

	case 106: //c2s_buy_prosperity
		return &MsgData{ModuleID: 7, SequenceID: 106}, nil

	case 46: //c2s_switch_action
		p := &pb7.C2SSwitchActionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_switch_action UnmarshalMsgProto &C2SSwitchActionProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 46, Proto: p}, nil

	case 165: //c2s_request_military_push
		p := &pb7.C2SRequestMilitaryPushProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_request_military_push UnmarshalMsgProto &C2SRequestMilitaryPushProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 165, Proto: p}, nil

	case 1: //c2s_create_base
		p := &pb7.C2SCreateBaseProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_create_base UnmarshalMsgProto &C2SCreateBaseProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 1, Proto: p}, nil

	case 14: //c2s_fast_move_base
		p := &pb7.C2SFastMoveBaseProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_fast_move_base UnmarshalMsgProto &C2SFastMoveBaseProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 14, Proto: p}, nil

	case 24: //c2s_invasion
		p := &pb7.C2SInvasionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_invasion UnmarshalMsgProto &C2SInvasionProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 24, Proto: p}, nil

	case 26: //c2s_cancel_invasion
		p := &pb7.C2SCancelInvasionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_cancel_invasion UnmarshalMsgProto &C2SCancelInvasionProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 26, Proto: p}, nil

	case 71: //c2s_repatriate
		p := &pb7.C2SRepatriateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_repatriate UnmarshalMsgProto &C2SRepatriateProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 71, Proto: p}, nil

	case 186: //c2s_baoz_repatriate
		p := &pb7.C2SBaozRepatriateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_baoz_repatriate UnmarshalMsgProto &C2SBaozRepatriateProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 186, Proto: p}, nil

	case 139: //c2s_speed_up
		p := &pb7.C2SSpeedUpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_speed_up UnmarshalMsgProto &C2SSpeedUpProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 139, Proto: p}, nil

	case 30: //c2s_expel
		p := &pb7.C2SExpelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_expel UnmarshalMsgProto &C2SExpelProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 30, Proto: p}, nil

	case 99: //c2s_favorite_pos
		p := &pb7.C2SFavoritePosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_favorite_pos UnmarshalMsgProto &C2SFavoritePosProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 99, Proto: p}, nil

	case 102: //c2s_favorite_pos_list
		return &MsgData{ModuleID: 7, SequenceID: 102}, nil

	case 175: //c2s_get_prev_investigate
		p := &pb7.C2SGetPrevInvestigateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_get_prev_investigate UnmarshalMsgProto &C2SGetPrevInvestigateProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 175, Proto: p}, nil

	case 142: //c2s_investigate
		p := &pb7.C2SInvestigateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_investigate UnmarshalMsgProto &C2SInvestigateProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 142, Proto: p}, nil

	case 234: //c2s_investigate_invade
		p := &pb7.C2SInvestigateInvadeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_investigate_invade UnmarshalMsgProto &C2SInvestigateInvadeProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 234, Proto: p}, nil

	case 183: //c2s_use_multi_level_npc_times_goods
		p := &pb7.C2SUseMultiLevelNpcTimesGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_use_multi_level_npc_times_goods UnmarshalMsgProto &C2SUseMultiLevelNpcTimesGoodsProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 183, Proto: p}, nil

	case 190: //c2s_use_invase_hero_times_goods
		p := &pb7.C2SUseInvaseHeroTimesGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_use_invase_hero_times_goods UnmarshalMsgProto &C2SUseInvaseHeroTimesGoodsProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 190, Proto: p}, nil

	case 172: //c2s_calc_move_speed
		p := &pb7.C2SCalcMoveSpeedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_calc_move_speed UnmarshalMsgProto &C2SCalcMoveSpeedProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 172, Proto: p}, nil

	case 178: //c2s_list_enemy_pos
		return &MsgData{ModuleID: 7, SequenceID: 178}, nil

	case 180: //c2s_search_baoz_npc
		p := &pb7.C2SSearchBaozNpcProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_search_baoz_npc UnmarshalMsgProto &C2SSearchBaozNpcProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 180, Proto: p}, nil

	case 193: //c2s_home_ast_defending_info
		return &MsgData{ModuleID: 7, SequenceID: 193}, nil

	case 196: //c2s_guild_please_help_me
		return &MsgData{ModuleID: 7, SequenceID: 196}, nil

	case 199: //c2s_create_assembly
		p := &pb7.C2SCreateAssemblyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_create_assembly UnmarshalMsgProto &C2SCreateAssemblyProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 199, Proto: p}, nil

	case 202: //c2s_show_assembly
		p := &pb7.C2SShowAssemblyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_show_assembly UnmarshalMsgProto &C2SShowAssemblyProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 202, Proto: p}, nil

	case 206: //c2s_join_assembly
		p := &pb7.C2SJoinAssemblyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_join_assembly UnmarshalMsgProto &C2SJoinAssemblyProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 206, Proto: p}, nil

	case 214: //c2s_create_guild_workshop
		p := &pb7.C2SCreateGuildWorkshopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_create_guild_workshop UnmarshalMsgProto &C2SCreateGuildWorkshopProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 214, Proto: p}, nil

	case 217: //c2s_show_guild_workshop
		p := &pb7.C2SShowGuildWorkshopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_show_guild_workshop UnmarshalMsgProto &C2SShowGuildWorkshopProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 217, Proto: p}, nil

	case 219: //c2s_hurt_guild_workshop
		p := &pb7.C2SHurtGuildWorkshopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_hurt_guild_workshop UnmarshalMsgProto &C2SHurtGuildWorkshopProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 219, Proto: p}, nil

	case 223: //c2s_remove_guild_workshop
		return &MsgData{ModuleID: 7, SequenceID: 223}, nil

	case 228: //c2s_catch_guild_workshop_logs
		p := &pb7.C2SCatchGuildWorkshopLogsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.c2s_catch_guild_workshop_logs UnmarshalMsgProto &C2SCatchGuildWorkshopLogsProto fail")
		}

		return &MsgData{ModuleID: 7, SequenceID: 228, Proto: p}, nil

	case 232: //c2s_get_self_baoz
		return &MsgData{ModuleID: 7, SequenceID: 232}, nil

	default:
		return nil, errors.New(fmt.Sprintf("region收到未知消息: %d", sequenceID))
	}
}

func unmarshal_mail(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_list_mail
		p := &pb8.C2SListMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.c2s_list_mail UnmarshalMsgProto &C2SListMailProto fail")
		}

		return &MsgData{ModuleID: 8, SequenceID: 1, Proto: p}, nil

	case 8: //c2s_delete_mail
		p := &pb8.C2SDeleteMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.c2s_delete_mail UnmarshalMsgProto &C2SDeleteMailProto fail")
		}

		return &MsgData{ModuleID: 8, SequenceID: 8, Proto: p}, nil

	case 11: //c2s_keep_mail
		p := &pb8.C2SKeepMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.c2s_keep_mail UnmarshalMsgProto &C2SKeepMailProto fail")
		}

		return &MsgData{ModuleID: 8, SequenceID: 11, Proto: p}, nil

	case 14: //c2s_collect_mail_prize
		p := &pb8.C2SCollectMailPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.c2s_collect_mail_prize UnmarshalMsgProto &C2SCollectMailPrizeProto fail")
		}

		return &MsgData{ModuleID: 8, SequenceID: 14, Proto: p}, nil

	case 20: //c2s_read_mail
		p := &pb8.C2SReadMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.c2s_read_mail UnmarshalMsgProto &C2SReadMailProto fail")
		}

		return &MsgData{ModuleID: 8, SequenceID: 20, Proto: p}, nil

	case 24: //c2s_read_multi
		p := &pb8.C2SReadMultiProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.c2s_read_multi UnmarshalMsgProto &C2SReadMultiProto fail")
		}

		return &MsgData{ModuleID: 8, SequenceID: 24, Proto: p}, nil

	case 26: //c2s_delete_multi
		p := &pb8.C2SDeleteMultiProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.c2s_delete_multi UnmarshalMsgProto &C2SDeleteMultiProto fail")
		}

		return &MsgData{ModuleID: 8, SequenceID: 26, Proto: p}, nil

	case 28: //c2s_get_mail
		p := &pb8.C2SGetMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.c2s_get_mail UnmarshalMsgProto &C2SGetMailProto fail")
		}

		return &MsgData{ModuleID: 8, SequenceID: 28, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("mail收到未知消息: %d", sequenceID))
	}
}

func unmarshal_guild(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_list_guild
		return &MsgData{ModuleID: 9, SequenceID: 1}, nil

	case 4: //c2s_search_guild
		p := &pb9.C2SSearchGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_search_guild UnmarshalMsgProto &C2SSearchGuildProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 4, Proto: p}, nil

	case 7: //c2s_create_guild
		p := &pb9.C2SCreateGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_create_guild UnmarshalMsgProto &C2SCreateGuildProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 7, Proto: p}, nil

	case 10: //c2s_self_guild
		p := &pb9.C2SSelfGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_self_guild UnmarshalMsgProto &C2SSelfGuildProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 10, Proto: p}, nil

	case 13: //c2s_leave_guild
		return &MsgData{ModuleID: 9, SequenceID: 13}, nil

	case 17: //c2s_kick_other
		p := &pb9.C2SKickOtherProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_kick_other UnmarshalMsgProto &C2SKickOtherProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 17, Proto: p}, nil

	case 20: //c2s_update_text
		p := &pb9.C2SUpdateTextProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_text UnmarshalMsgProto &C2SUpdateTextProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 20, Proto: p}, nil

	case 65: //c2s_update_internal_text
		p := &pb9.C2SUpdateInternalTextProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_internal_text UnmarshalMsgProto &C2SUpdateInternalTextProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 65, Proto: p}, nil

	case 23: //c2s_update_class_names
		p := &pb9.C2SUpdateClassNamesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_class_names UnmarshalMsgProto &C2SUpdateClassNamesProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 23, Proto: p}, nil

	case 122: //c2s_update_class_title
		p := &pb9.C2SUpdateClassTitleProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_class_title UnmarshalMsgProto &C2SUpdateClassTitleProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 122, Proto: p}, nil

	case 26: //c2s_update_flag_type
		p := &pb9.C2SUpdateFlagTypeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_flag_type UnmarshalMsgProto &C2SUpdateFlagTypeProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 26, Proto: p}, nil

	case 29: //c2s_update_member_class_level
		p := &pb9.C2SUpdateMemberClassLevelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_member_class_level UnmarshalMsgProto &C2SUpdateMemberClassLevelProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 29, Proto: p}, nil

	case 80: //c2s_cancel_change_leader
		return &MsgData{ModuleID: 9, SequenceID: 80}, nil

	case 68: //c2s_update_join_condition
		p := &pb9.C2SUpdateJoinConditionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_join_condition UnmarshalMsgProto &C2SUpdateJoinConditionProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 68, Proto: p}, nil

	case 71: //c2s_update_guild_name
		p := &pb9.C2SUpdateGuildNameProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_guild_name UnmarshalMsgProto &C2SUpdateGuildNameProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 71, Proto: p}, nil

	case 75: //c2s_update_guild_label
		p := &pb9.C2SUpdateGuildLabelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_guild_label UnmarshalMsgProto &C2SUpdateGuildLabelProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 75, Proto: p}, nil

	case 83: //c2s_donate
		p := &pb9.C2SDonateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_donate UnmarshalMsgProto &C2SDonateProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 83, Proto: p}, nil

	case 90: //c2s_upgrade_level
		return &MsgData{ModuleID: 9, SequenceID: 90}, nil

	case 93: //c2s_reduce_upgrade_level_cd
		p := &pb9.C2SReduceUpgradeLevelCdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_reduce_upgrade_level_cd UnmarshalMsgProto &C2SReduceUpgradeLevelCdProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 93, Proto: p}, nil

	case 96: //c2s_impeach_leader
		return &MsgData{ModuleID: 9, SequenceID: 96}, nil

	case 99: //c2s_impeach_leader_vote
		p := &pb9.C2SImpeachLeaderVoteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_impeach_leader_vote UnmarshalMsgProto &C2SImpeachLeaderVoteProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 99, Proto: p}, nil

	case 102: //c2s_list_guild_by_ids
		p := &pb9.C2SListGuildByIdsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_list_guild_by_ids UnmarshalMsgProto &C2SListGuildByIdsProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 102, Proto: p}, nil

	case 40: //c2s_user_request_join
		p := &pb9.C2SUserRequestJoinProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_user_request_join UnmarshalMsgProto &C2SUserRequestJoinProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 40, Proto: p}, nil

	case 43: //c2s_user_cancel_join_request
		p := &pb9.C2SUserCancelJoinRequestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_user_cancel_join_request UnmarshalMsgProto &C2SUserCancelJoinRequestProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 43, Proto: p}, nil

	case 55: //c2s_guild_reply_join_request
		p := &pb9.C2SGuildReplyJoinRequestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_guild_reply_join_request UnmarshalMsgProto &C2SGuildReplyJoinRequestProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 55, Proto: p}, nil

	case 109: //c2s_guild_invate_other
		p := &pb9.C2SGuildInvateOtherProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_guild_invate_other UnmarshalMsgProto &C2SGuildInvateOtherProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 109, Proto: p}, nil

	case 112: //c2s_guild_cancel_invate_other
		p := &pb9.C2SGuildCancelInvateOtherProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_guild_cancel_invate_other UnmarshalMsgProto &C2SGuildCancelInvateOtherProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 112, Proto: p}, nil

	case 48: //c2s_user_reply_invate_request
		p := &pb9.C2SUserReplyInvateRequestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_user_reply_invate_request UnmarshalMsgProto &C2SUserReplyInvateRequestProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 48, Proto: p}, nil

	case 193: //c2s_list_invite_me_guild
		return &MsgData{ModuleID: 9, SequenceID: 193}, nil

	case 125: //c2s_update_friend_guild
		p := &pb9.C2SUpdateFriendGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_friend_guild UnmarshalMsgProto &C2SUpdateFriendGuildProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 125, Proto: p}, nil

	case 128: //c2s_update_enemy_guild
		p := &pb9.C2SUpdateEnemyGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_enemy_guild UnmarshalMsgProto &C2SUpdateEnemyGuildProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 128, Proto: p}, nil

	case 131: //c2s_update_guild_prestige
		p := &pb9.C2SUpdateGuildPrestigeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_guild_prestige UnmarshalMsgProto &C2SUpdateGuildPrestigeProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 131, Proto: p}, nil

	case 134: //c2s_place_guild_statue
		p := &pb9.C2SPlaceGuildStatueProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_place_guild_statue UnmarshalMsgProto &C2SPlaceGuildStatueProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 134, Proto: p}, nil

	case 138: //c2s_take_back_guild_statue
		return &MsgData{ModuleID: 9, SequenceID: 138}, nil

	case 143: //c2s_collect_first_join_guild_prize
		return &MsgData{ModuleID: 9, SequenceID: 143}, nil

	case 147: //c2s_seek_help
		p := &pb9.C2SSeekHelpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_seek_help UnmarshalMsgProto &C2SSeekHelpProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 147, Proto: p}, nil

	case 151: //c2s_help_guild_member
		p := &pb9.C2SHelpGuildMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_help_guild_member UnmarshalMsgProto &C2SHelpGuildMemberProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 151, Proto: p}, nil

	case 158: //c2s_help_all_guild_member
		return &MsgData{ModuleID: 9, SequenceID: 158}, nil

	case 163: //c2s_collect_guild_event_prize
		p := &pb9.C2SCollectGuildEventPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_collect_guild_event_prize UnmarshalMsgProto &C2SCollectGuildEventPrizeProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 163, Proto: p}, nil

	case 167: //c2s_collect_full_big_box
		return &MsgData{ModuleID: 9, SequenceID: 167}, nil

	case 172: //c2s_upgrade_technology
		p := &pb9.C2SUpgradeTechnologyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_upgrade_technology UnmarshalMsgProto &C2SUpgradeTechnologyProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 172, Proto: p}, nil

	case 175: //c2s_reduce_technology_cd
		p := &pb9.C2SReduceTechnologyCdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_reduce_technology_cd UnmarshalMsgProto &C2SReduceTechnologyCdProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 175, Proto: p}, nil

	case 178: //c2s_list_guild_logs
		p := &pb9.C2SListGuildLogsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_list_guild_logs UnmarshalMsgProto &C2SListGuildLogsProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 178, Proto: p}, nil

	case 181: //c2s_request_recommend_guild
		return &MsgData{ModuleID: 9, SequenceID: 181}, nil

	case 184: //c2s_help_tech
		return &MsgData{ModuleID: 9, SequenceID: 184}, nil

	case 187: //c2s_recommend_invite_heros
		return &MsgData{ModuleID: 9, SequenceID: 187}, nil

	case 190: //c2s_search_no_guild_heros
		p := &pb9.C2SSearchNoGuildHerosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_search_no_guild_heros UnmarshalMsgProto &C2SSearchNoGuildHerosProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 190, Proto: p}, nil

	case 199: //c2s_view_mc_war_record
		return &MsgData{ModuleID: 9, SequenceID: 199}, nil

	case 196: //c2s_update_guild_mark
		p := &pb9.C2SUpdateGuildMarkProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_update_guild_mark UnmarshalMsgProto &C2SUpdateGuildMarkProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 196, Proto: p}, nil

	case 202: //c2s_view_yinliang_record
		return &MsgData{ModuleID: 9, SequenceID: 202}, nil

	case 205: //c2s_send_yinliang_to_other_guild
		p := &pb9.C2SSendYinliangToOtherGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_send_yinliang_to_other_guild UnmarshalMsgProto &C2SSendYinliangToOtherGuildProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 205, Proto: p}, nil

	case 208: //c2s_send_yinliang_to_member
		p := &pb9.C2SSendYinliangToMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_send_yinliang_to_member UnmarshalMsgProto &C2SSendYinliangToMemberProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 208, Proto: p}, nil

	case 211: //c2s_pay_salary
		return &MsgData{ModuleID: 9, SequenceID: 211}, nil

	case 214: //c2s_set_salary
		p := &pb9.C2SSetSalaryProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_set_salary UnmarshalMsgProto &C2SSetSalaryProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 214, Proto: p}, nil

	case 218: //c2s_view_send_yinliang_to_guild
		return &MsgData{ModuleID: 9, SequenceID: 218}, nil

	case 228: //c2s_convene
		p := &pb9.C2SConveneProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_convene UnmarshalMsgProto &C2SConveneProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 228, Proto: p}, nil

	case 231: //c2s_collect_daily_guild_rank_prize
		return &MsgData{ModuleID: 9, SequenceID: 231}, nil

	case 234: //c2s_view_daily_guild_rank
		return &MsgData{ModuleID: 9, SequenceID: 234}, nil

	case 240: //c2s_add_recommend_mc_build
		p := &pb9.C2SAddRecommendMcBuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_add_recommend_mc_build UnmarshalMsgProto &C2SAddRecommendMcBuildProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 240, Proto: p}, nil

	case 243: //c2s_view_task_progress
		p := &pb9.C2SViewTaskProgressProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_view_task_progress UnmarshalMsgProto &C2SViewTaskProgressProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 243, Proto: p}, nil

	case 247: //c2s_collect_task_prize
		p := &pb9.C2SCollectTaskPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_collect_task_prize UnmarshalMsgProto &C2SCollectTaskPrizeProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 247, Proto: p}, nil

	case 250: //c2s_guild_change_country
		p := &pb9.C2SGuildChangeCountryProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.c2s_guild_change_country UnmarshalMsgProto &C2SGuildChangeCountryProto fail")
		}

		return &MsgData{ModuleID: 9, SequenceID: 250, Proto: p}, nil

	case 253: //c2s_cancel_guild_change_country
		return &MsgData{ModuleID: 9, SequenceID: 253}, nil

	case 256: //c2s_show_workshop_not_exist
		return &MsgData{ModuleID: 9, SequenceID: 256}, nil

	default:
		return nil, errors.New(fmt.Sprintf("guild收到未知消息: %d", sequenceID))
	}
}

func unmarshal_stress(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_robot_ping
		p := &pb10.C2SRobotPingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "stress.c2s_robot_ping UnmarshalMsgProto &C2SRobotPingProto fail")
		}

		return &MsgData{ModuleID: 10, SequenceID: 1, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("stress收到未知消息: %d", sequenceID))
	}
}

func unmarshal_depot(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 2: //c2s_use_goods
		p := &pb11.C2SUseGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.c2s_use_goods UnmarshalMsgProto &C2SUseGoodsProto fail")
		}

		return &MsgData{ModuleID: 11, SequenceID: 2, Proto: p}, nil

	case 6: //c2s_use_cdr_goods
		p := &pb11.C2SUseCdrGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.c2s_use_cdr_goods UnmarshalMsgProto &C2SUseCdrGoodsProto fail")
		}

		return &MsgData{ModuleID: 11, SequenceID: 6, Proto: p}, nil

	case 9: //c2s_goods_combine
		p := &pb11.C2SGoodsCombineProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.c2s_goods_combine UnmarshalMsgProto &C2SGoodsCombineProto fail")
		}

		return &MsgData{ModuleID: 11, SequenceID: 9, Proto: p}, nil

	case 18: //c2s_goods_parts_combine
		p := &pb11.C2SGoodsPartsCombineProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.c2s_goods_parts_combine UnmarshalMsgProto &C2SGoodsPartsCombineProto fail")
		}

		return &MsgData{ModuleID: 11, SequenceID: 18, Proto: p}, nil

	case 23: //c2s_unlock_baowu
		p := &pb11.C2SUnlockBaowuProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.c2s_unlock_baowu UnmarshalMsgProto &C2SUnlockBaowuProto fail")
		}

		return &MsgData{ModuleID: 11, SequenceID: 23, Proto: p}, nil

	case 26: //c2s_collect_baowu
		p := &pb11.C2SCollectBaowuProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.c2s_collect_baowu UnmarshalMsgProto &C2SCollectBaowuProto fail")
		}

		return &MsgData{ModuleID: 11, SequenceID: 26, Proto: p}, nil

	case 30: //c2s_list_baowu_log
		return &MsgData{ModuleID: 11, SequenceID: 30}, nil

	case 35: //c2s_decompose_baowu
		p := &pb11.C2SDecomposeBaowuProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.c2s_decompose_baowu UnmarshalMsgProto &C2SDecomposeBaowuProto fail")
		}

		return &MsgData{ModuleID: 11, SequenceID: 35, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("depot收到未知消息: %d", sequenceID))
	}
}

func unmarshal_equipment(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 40: //c2s_view_chat_equip
		p := &pb12.C2SViewChatEquipProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.c2s_view_chat_equip UnmarshalMsgProto &C2SViewChatEquipProto fail")
		}

		return &MsgData{ModuleID: 12, SequenceID: 40, Proto: p}, nil

	case 1: //c2s_wear_equipment
		p := &pb12.C2SWearEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.c2s_wear_equipment UnmarshalMsgProto &C2SWearEquipmentProto fail")
		}

		return &MsgData{ModuleID: 12, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_upgrade_equipment
		p := &pb12.C2SUpgradeEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.c2s_upgrade_equipment UnmarshalMsgProto &C2SUpgradeEquipmentProto fail")
		}

		return &MsgData{ModuleID: 12, SequenceID: 4, Proto: p}, nil

	case 19: //c2s_upgrade_equipment_all
		p := &pb12.C2SUpgradeEquipmentAllProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.c2s_upgrade_equipment_all UnmarshalMsgProto &C2SUpgradeEquipmentAllProto fail")
		}

		return &MsgData{ModuleID: 12, SequenceID: 19, Proto: p}, nil

	case 7: //c2s_refined_equipment
		p := &pb12.C2SRefinedEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.c2s_refined_equipment UnmarshalMsgProto &C2SRefinedEquipmentProto fail")
		}

		return &MsgData{ModuleID: 12, SequenceID: 7, Proto: p}, nil

	case 10: //c2s_smelt_equipment
		p := &pb12.C2SSmeltEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.c2s_smelt_equipment UnmarshalMsgProto &C2SSmeltEquipmentProto fail")
		}

		return &MsgData{ModuleID: 12, SequenceID: 10, Proto: p}, nil

	case 13: //c2s_rebuild_equipment
		p := &pb12.C2SRebuildEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.c2s_rebuild_equipment UnmarshalMsgProto &C2SRebuildEquipmentProto fail")
		}

		return &MsgData{ModuleID: 12, SequenceID: 13, Proto: p}, nil

	case 43: //c2s_one_key_take_off
		p := &pb12.C2SOneKeyTakeOffProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.c2s_one_key_take_off UnmarshalMsgProto &C2SOneKeyTakeOffProto fail")
		}

		return &MsgData{ModuleID: 12, SequenceID: 43, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("equipment收到未知消息: %d", sequenceID))
	}
}

func unmarshal_chat(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_world_chat
		p := &pb13.C2SWorldChatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.c2s_world_chat UnmarshalMsgProto &C2SWorldChatProto fail")
		}

		return &MsgData{ModuleID: 13, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_guild_chat
		p := &pb13.C2SGuildChatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.c2s_guild_chat UnmarshalMsgProto &C2SGuildChatProto fail")
		}

		return &MsgData{ModuleID: 13, SequenceID: 4, Proto: p}, nil

	case 8: //c2s_self_chat_window
		return &MsgData{ModuleID: 13, SequenceID: 8}, nil

	case 21: //c2s_create_self_chat_window
		p := &pb13.C2SCreateSelfChatWindowProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.c2s_create_self_chat_window UnmarshalMsgProto &C2SCreateSelfChatWindowProto fail")
		}

		return &MsgData{ModuleID: 13, SequenceID: 21, Proto: p}, nil

	case 10: //c2s_remove_chat_window
		p := &pb13.C2SRemoveChatWindowProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.c2s_remove_chat_window UnmarshalMsgProto &C2SRemoveChatWindowProto fail")
		}

		return &MsgData{ModuleID: 13, SequenceID: 10, Proto: p}, nil

	case 12: //c2s_list_history_chat
		p := &pb13.C2SListHistoryChatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.c2s_list_history_chat UnmarshalMsgProto &C2SListHistoryChatProto fail")
		}

		return &MsgData{ModuleID: 13, SequenceID: 12, Proto: p}, nil

	case 14: //c2s_send_chat
		p := &pb13.C2SSendChatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.c2s_send_chat UnmarshalMsgProto &C2SSendChatProto fail")
		}

		return &MsgData{ModuleID: 13, SequenceID: 14, Proto: p}, nil

	case 18: //c2s_read_chat_msg
		p := &pb13.C2SReadChatMsgProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.c2s_read_chat_msg UnmarshalMsgProto &C2SReadChatMsgProto fail")
		}

		return &MsgData{ModuleID: 13, SequenceID: 18, Proto: p}, nil

	case 25: //c2s_get_hero_chat_info
		p := &pb13.C2SGetHeroChatInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.c2s_get_hero_chat_info UnmarshalMsgProto &C2SGetHeroChatInfoProto fail")
		}

		return &MsgData{ModuleID: 13, SequenceID: 25, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("chat收到未知消息: %d", sequenceID))
	}
}

func unmarshal_tower(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_challenge
		p := &pb14.C2SChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tower.c2s_challenge UnmarshalMsgProto &C2SChallengeProto fail")
		}

		return &MsgData{ModuleID: 14, SequenceID: 1, Proto: p}, nil

	case 5: //c2s_auto_challenge
		return &MsgData{ModuleID: 14, SequenceID: 5}, nil

	case 8: //c2s_collect_box
		p := &pb14.C2SCollectBoxProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tower.c2s_collect_box UnmarshalMsgProto &C2SCollectBoxProto fail")
		}

		return &MsgData{ModuleID: 14, SequenceID: 8, Proto: p}, nil

	case 11: //c2s_list_pass_replay
		p := &pb14.C2SListPassReplayProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tower.c2s_list_pass_replay UnmarshalMsgProto &C2SListPassReplayProto fail")
		}

		return &MsgData{ModuleID: 14, SequenceID: 11, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("tower收到未知消息: %d", sequenceID))
	}
}

func unmarshal_task(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 2: //c2s_collect_task_prize
		p := &pb15.C2SCollectTaskPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.c2s_collect_task_prize UnmarshalMsgProto &C2SCollectTaskPrizeProto fail")
		}

		return &MsgData{ModuleID: 15, SequenceID: 2, Proto: p}, nil

	case 6: //c2s_collect_task_box_prize
		p := &pb15.C2SCollectTaskBoxPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.c2s_collect_task_box_prize UnmarshalMsgProto &C2SCollectTaskBoxPrizeProto fail")
		}

		return &MsgData{ModuleID: 15, SequenceID: 6, Proto: p}, nil

	case 9: //c2s_collect_ba_ye_stage_prize
		return &MsgData{ModuleID: 15, SequenceID: 9}, nil

	case 12: //c2s_collect_active_degree_prize
		p := &pb15.C2SCollectActiveDegreePrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.c2s_collect_active_degree_prize UnmarshalMsgProto &C2SCollectActiveDegreePrizeProto fail")
		}

		return &MsgData{ModuleID: 15, SequenceID: 12, Proto: p}, nil

	case 16: //c2s_collect_achieve_star_prize
		p := &pb15.C2SCollectAchieveStarPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.c2s_collect_achieve_star_prize UnmarshalMsgProto &C2SCollectAchieveStarPrizeProto fail")
		}

		return &MsgData{ModuleID: 15, SequenceID: 16, Proto: p}, nil

	case 20: //c2s_change_select_show_achieve
		p := &pb15.C2SChangeSelectShowAchieveProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.c2s_change_select_show_achieve UnmarshalMsgProto &C2SChangeSelectShowAchieveProto fail")
		}

		return &MsgData{ModuleID: 15, SequenceID: 20, Proto: p}, nil

	case 23: //c2s_collect_bwzl_prize
		p := &pb15.C2SCollectBwzlPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.c2s_collect_bwzl_prize UnmarshalMsgProto &C2SCollectBwzlPrizeProto fail")
		}

		return &MsgData{ModuleID: 15, SequenceID: 23, Proto: p}, nil

	case 26: //c2s_view_other_achieve_task_list
		p := &pb15.C2SViewOtherAchieveTaskListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.c2s_view_other_achieve_task_list UnmarshalMsgProto &C2SViewOtherAchieveTaskListProto fail")
		}

		return &MsgData{ModuleID: 15, SequenceID: 26, Proto: p}, nil

	case 29: //c2s_get_troop_title_fight_amount
		p := &pb15.C2SGetTroopTitleFightAmountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.c2s_get_troop_title_fight_amount UnmarshalMsgProto &C2SGetTroopTitleFightAmountProto fail")
		}

		return &MsgData{ModuleID: 15, SequenceID: 29, Proto: p}, nil

	case 31: //c2s_upgrade_title
		return &MsgData{ModuleID: 15, SequenceID: 31}, nil

	case 35: //c2s_complete_bool_task
		p := &pb15.C2SCompleteBoolTaskProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.c2s_complete_bool_task UnmarshalMsgProto &C2SCompleteBoolTaskProto fail")
		}

		return &MsgData{ModuleID: 15, SequenceID: 35, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("task收到未知消息: %d", sequenceID))
	}
}

func unmarshal_fishing(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_fishing
		p := &pb16.C2SFishingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "fishing.c2s_fishing UnmarshalMsgProto &C2SFishingProto fail")
		}

		return &MsgData{ModuleID: 16, SequenceID: 1, Proto: p}, nil

	case 8: //c2s_fish_point_exchange
		return &MsgData{ModuleID: 16, SequenceID: 8}, nil

	case 11: //c2s_set_fishing_captain
		p := &pb16.C2SSetFishingCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "fishing.c2s_set_fishing_captain UnmarshalMsgProto &C2SSetFishingCaptainProto fail")
		}

		return &MsgData{ModuleID: 16, SequenceID: 11, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("fishing收到未知消息: %d", sequenceID))
	}
}

func unmarshal_gem(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 3: //c2s_use_gem
		p := &pb19.C2SUseGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.c2s_use_gem UnmarshalMsgProto &C2SUseGemProto fail")
		}

		return &MsgData{ModuleID: 19, SequenceID: 3, Proto: p}, nil

	case 21: //c2s_inlay_gem
		p := &pb19.C2SInlayGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.c2s_inlay_gem UnmarshalMsgProto &C2SInlayGemProto fail")
		}

		return &MsgData{ModuleID: 19, SequenceID: 21, Proto: p}, nil

	case 6: //c2s_combine_gem
		p := &pb19.C2SCombineGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.c2s_combine_gem UnmarshalMsgProto &C2SCombineGemProto fail")
		}

		return &MsgData{ModuleID: 19, SequenceID: 6, Proto: p}, nil

	case 9: //c2s_one_key_use_gem
		p := &pb19.C2SOneKeyUseGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.c2s_one_key_use_gem UnmarshalMsgProto &C2SOneKeyUseGemProto fail")
		}

		return &MsgData{ModuleID: 19, SequenceID: 9, Proto: p}, nil

	case 11: //c2s_one_key_combine_gem
		p := &pb19.C2SOneKeyCombineGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.c2s_one_key_combine_gem UnmarshalMsgProto &C2SOneKeyCombineGemProto fail")
		}

		return &MsgData{ModuleID: 19, SequenceID: 11, Proto: p}, nil

	case 15: //c2s_request_one_key_combine_cost
		p := &pb19.C2SRequestOneKeyCombineCostProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.c2s_request_one_key_combine_cost UnmarshalMsgProto &C2SRequestOneKeyCombineCostProto fail")
		}

		return &MsgData{ModuleID: 19, SequenceID: 15, Proto: p}, nil

	case 18: //c2s_one_key_combine_depot_gem
		p := &pb19.C2SOneKeyCombineDepotGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.c2s_one_key_combine_depot_gem UnmarshalMsgProto &C2SOneKeyCombineDepotGemProto fail")
		}

		return &MsgData{ModuleID: 19, SequenceID: 18, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("gem收到未知消息: %d", sequenceID))
	}
}

func unmarshal_shop(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 2: //c2s_buy_goods
		p := &pb20.C2SBuyGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "shop.c2s_buy_goods UnmarshalMsgProto &C2SBuyGoodsProto fail")
		}

		return &MsgData{ModuleID: 20, SequenceID: 2, Proto: p}, nil

	case 9: //c2s_buy_black_market_goods
		p := &pb20.C2SBuyBlackMarketGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "shop.c2s_buy_black_market_goods UnmarshalMsgProto &C2SBuyBlackMarketGoodsProto fail")
		}

		return &MsgData{ModuleID: 20, SequenceID: 9, Proto: p}, nil

	case 12: //c2s_refresh_black_market_goods
		return &MsgData{ModuleID: 20, SequenceID: 12}, nil

	default:
		return nil, errors.New(fmt.Sprintf("shop收到未知消息: %d", sequenceID))
	}
}

func unmarshal_client_config(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_config
		p := &pb21.C2SConfigProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "client_config.c2s_config UnmarshalMsgProto &C2SConfigProto fail")
		}

		return &MsgData{ModuleID: 21, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_set_client_data
		p := &pb21.C2SSetClientDataProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "client_config.c2s_set_client_data UnmarshalMsgProto &C2SSetClientDataProto fail")
		}

		return &MsgData{ModuleID: 21, SequenceID: 4, Proto: p}, nil

	case 5: //c2s_set_client_key
		p := &pb21.C2SSetClientKeyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "client_config.c2s_set_client_key UnmarshalMsgProto &C2SSetClientKeyProto fail")
		}

		return &MsgData{ModuleID: 21, SequenceID: 5, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("client_config收到未知消息: %d", sequenceID))
	}
}

func unmarshal_secret_tower(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 2: //c2s_request_team_count
		return &MsgData{ModuleID: 22, SequenceID: 2}, nil

	case 5: //c2s_request_team_list
		p := &pb22.C2SRequestTeamListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_request_team_list UnmarshalMsgProto &C2SRequestTeamListProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 5, Proto: p}, nil

	case 8: //c2s_create_team
		p := &pb22.C2SCreateTeamProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_create_team UnmarshalMsgProto &C2SCreateTeamProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 8, Proto: p}, nil

	case 11: //c2s_join_team
		p := &pb22.C2SJoinTeamProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_join_team UnmarshalMsgProto &C2SJoinTeamProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 11, Proto: p}, nil

	case 15: //c2s_leave_team
		return &MsgData{ModuleID: 22, SequenceID: 15}, nil

	case 19: //c2s_kick_member
		p := &pb22.C2SKickMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_kick_member UnmarshalMsgProto &C2SKickMemberProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 19, Proto: p}, nil

	case 24: //c2s_move_member
		p := &pb22.C2SMoveMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_move_member UnmarshalMsgProto &C2SMoveMemberProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 24, Proto: p}, nil

	case 67: //c2s_update_member_pos
		p := &pb22.C2SUpdateMemberPosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_update_member_pos UnmarshalMsgProto &C2SUpdateMemberPosProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 67, Proto: p}, nil

	case 27: //c2s_change_mode
		p := &pb22.C2SChangeModeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_change_mode UnmarshalMsgProto &C2SChangeModeProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 27, Proto: p}, nil

	case 33: //c2s_invite
		p := &pb22.C2SInviteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_invite UnmarshalMsgProto &C2SInviteProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 33, Proto: p}, nil

	case 71: //c2s_invite_all
		p := &pb22.C2SInviteAllProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_invite_all UnmarshalMsgProto &C2SInviteAllProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 71, Proto: p}, nil

	case 37: //c2s_request_invite_list
		return &MsgData{ModuleID: 22, SequenceID: 37}, nil

	case 39: //c2s_request_team_detail
		return &MsgData{ModuleID: 22, SequenceID: 39}, nil

	case 42: //c2s_start_challenge
		return &MsgData{ModuleID: 22, SequenceID: 42}, nil

	case 58: //c2s_quick_query_team_basic
		p := &pb22.C2SQuickQueryTeamBasicProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_quick_query_team_basic UnmarshalMsgProto &C2SQuickQueryTeamBasicProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 58, Proto: p}, nil

	case 61: //c2s_change_guild_mode
		return &MsgData{ModuleID: 22, SequenceID: 61}, nil

	case 74: //c2s_list_record
		return &MsgData{ModuleID: 22, SequenceID: 74}, nil

	case 79: //c2s_team_talk
		p := &pb22.C2STeamTalkProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.c2s_team_talk UnmarshalMsgProto &C2STeamTalkProto fail")
		}

		return &MsgData{ModuleID: 22, SequenceID: 79, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("secret_tower收到未知消息: %d", sequenceID))
	}
}

func unmarshal_rank(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_request_rank
		p := &pb23.C2SRequestRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "rank.c2s_request_rank UnmarshalMsgProto &C2SRequestRankProto fail")
		}

		return &MsgData{ModuleID: 23, SequenceID: 1, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("rank收到未知消息: %d", sequenceID))
	}
}

func unmarshal_bai_zhan(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_query_bai_zhan_info
		return &MsgData{ModuleID: 24, SequenceID: 1}, nil

	case 34: //c2s_clear_last_jun_xian
		return &MsgData{ModuleID: 24, SequenceID: 34}, nil

	case 4: //c2s_bai_zhan_challenge
		return &MsgData{ModuleID: 24, SequenceID: 4}, nil

	case 7: //c2s_collect_salary
		return &MsgData{ModuleID: 24, SequenceID: 7}, nil

	case 10: //c2s_collect_jun_xian_prize
		p := &pb24.C2SCollectJunXianPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "bai_zhan.c2s_collect_jun_xian_prize UnmarshalMsgProto &C2SCollectJunXianPrizeProto fail")
		}

		return &MsgData{ModuleID: 24, SequenceID: 10, Proto: p}, nil

	case 29: //c2s_self_record
		p := &pb24.C2SSelfRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "bai_zhan.c2s_self_record UnmarshalMsgProto &C2SSelfRecordProto fail")
		}

		return &MsgData{ModuleID: 24, SequenceID: 29, Proto: p}, nil

	case 23: //c2s_request_rank
		p := &pb24.C2SRequestRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "bai_zhan.c2s_request_rank UnmarshalMsgProto &C2SRequestRankProto fail")
		}

		return &MsgData{ModuleID: 24, SequenceID: 23, Proto: p}, nil

	case 26: //c2s_request_self_rank
		return &MsgData{ModuleID: 24, SequenceID: 26}, nil

	default:
		return nil, errors.New(fmt.Sprintf("bai_zhan收到未知消息: %d", sequenceID))
	}
}

func unmarshal_dungeon(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_challenge
		p := &pb26.C2SChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.c2s_challenge UnmarshalMsgProto &C2SChallengeProto fail")
		}

		return &MsgData{ModuleID: 26, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_collect_chapter_prize
		p := &pb26.C2SCollectChapterPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.c2s_collect_chapter_prize UnmarshalMsgProto &C2SCollectChapterPrizeProto fail")
		}

		return &MsgData{ModuleID: 26, SequenceID: 4, Proto: p}, nil

	case 13: //c2s_collect_pass_dungeon_prize
		p := &pb26.C2SCollectPassDungeonPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.c2s_collect_pass_dungeon_prize UnmarshalMsgProto &C2SCollectPassDungeonPrizeProto fail")
		}

		return &MsgData{ModuleID: 26, SequenceID: 13, Proto: p}, nil

	case 7: //c2s_auto_challenge
		p := &pb26.C2SAutoChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.c2s_auto_challenge UnmarshalMsgProto &C2SAutoChallengeProto fail")
		}

		return &MsgData{ModuleID: 26, SequenceID: 7, Proto: p}, nil

	case 17: //c2s_collect_chapter_star_prize
		p := &pb26.C2SCollectChapterStarPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.c2s_collect_chapter_star_prize UnmarshalMsgProto &C2SCollectChapterStarPrizeProto fail")
		}

		return &MsgData{ModuleID: 26, SequenceID: 17, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("dungeon收到未知消息: %d", sequenceID))
	}
}

func unmarshal_country(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 16: //c2s_request_country_prestige
		p := &pb27.C2SRequestCountryPrestigeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.c2s_request_country_prestige UnmarshalMsgProto &C2SRequestCountryPrestigeProto fail")
		}

		return &MsgData{ModuleID: 27, SequenceID: 16, Proto: p}, nil

	case 19: //c2s_request_countries
		p := &pb27.C2SRequestCountriesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.c2s_request_countries UnmarshalMsgProto &C2SRequestCountriesProto fail")
		}

		return &MsgData{ModuleID: 27, SequenceID: 19, Proto: p}, nil

	case 22: //c2s_hero_change_country
		p := &pb27.C2SHeroChangeCountryProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.c2s_hero_change_country UnmarshalMsgProto &C2SHeroChangeCountryProto fail")
		}

		return &MsgData{ModuleID: 27, SequenceID: 22, Proto: p}, nil

	case 31: //c2s_country_detail
		p := &pb27.C2SCountryDetailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.c2s_country_detail UnmarshalMsgProto &C2SCountryDetailProto fail")
		}

		return &MsgData{ModuleID: 27, SequenceID: 31, Proto: p}, nil

	case 40: //c2s_official_appoint
		p := &pb27.C2SOfficialAppointProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.c2s_official_appoint UnmarshalMsgProto &C2SOfficialAppointProto fail")
		}

		return &MsgData{ModuleID: 27, SequenceID: 40, Proto: p}, nil

	case 43: //c2s_official_depose
		p := &pb27.C2SOfficialDeposeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.c2s_official_depose UnmarshalMsgProto &C2SOfficialDeposeProto fail")
		}

		return &MsgData{ModuleID: 27, SequenceID: 43, Proto: p}, nil

	case 54: //c2s_official_leave
		return &MsgData{ModuleID: 27, SequenceID: 54}, nil

	case 46: //c2s_collect_official_salary
		return &MsgData{ModuleID: 27, SequenceID: 46}, nil

	case 72: //c2s_change_name_start
		p := &pb27.C2SChangeNameStartProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.c2s_change_name_start UnmarshalMsgProto &C2SChangeNameStartProto fail")
		}

		return &MsgData{ModuleID: 27, SequenceID: 72, Proto: p}, nil

	case 61: //c2s_change_name_vote
		p := &pb27.C2SChangeNameVoteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.c2s_change_name_vote UnmarshalMsgProto &C2SChangeNameVoteProto fail")
		}

		return &MsgData{ModuleID: 27, SequenceID: 61, Proto: p}, nil

	case 66: //c2s_search_to_appoint_hero_list
		p := &pb27.C2SSearchToAppointHeroListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.c2s_search_to_appoint_hero_list UnmarshalMsgProto &C2SSearchToAppointHeroListProto fail")
		}

		return &MsgData{ModuleID: 27, SequenceID: 66, Proto: p}, nil

	case 69: //c2s_default_to_appoint_hero_list
		return &MsgData{ModuleID: 27, SequenceID: 69}, nil

	default:
		return nil, errors.New(fmt.Sprintf("country收到未知消息: %d", sequenceID))
	}
}

func unmarshal_tag(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_add_or_update_tag
		p := &pb29.C2SAddOrUpdateTagProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tag.c2s_add_or_update_tag UnmarshalMsgProto &C2SAddOrUpdateTagProto fail")
		}

		return &MsgData{ModuleID: 29, SequenceID: 1, Proto: p}, nil

	case 5: //c2s_delete_tag
		p := &pb29.C2SDeleteTagProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tag.c2s_delete_tag UnmarshalMsgProto &C2SDeleteTagProto fail")
		}

		return &MsgData{ModuleID: 29, SequenceID: 5, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("tag收到未知消息: %d", sequenceID))
	}
}

func unmarshal_garden(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_list_treasury_tree_hero
		return &MsgData{ModuleID: 31, SequenceID: 1}, nil

	case 13: //c2s_list_help_me
		p := &pb31.C2SListHelpMeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "garden.c2s_list_help_me UnmarshalMsgProto &C2SListHelpMeProto fail")
		}

		return &MsgData{ModuleID: 31, SequenceID: 13, Proto: p}, nil

	case 3: //c2s_list_treasury_tree_times
		return &MsgData{ModuleID: 31, SequenceID: 3}, nil

	case 5: //c2s_water_treasury_tree
		p := &pb31.C2SWaterTreasuryTreeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "garden.c2s_water_treasury_tree UnmarshalMsgProto &C2SWaterTreasuryTreeProto fail")
		}

		return &MsgData{ModuleID: 31, SequenceID: 5, Proto: p}, nil

	case 10: //c2s_collect_treasury_tree_prize
		return &MsgData{ModuleID: 31, SequenceID: 10}, nil

	default:
		return nil, errors.New(fmt.Sprintf("garden收到未知消息: %d", sequenceID))
	}
}

func unmarshal_zhengwu(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_start
		p := &pb32.C2SStartProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhengwu.c2s_start UnmarshalMsgProto &C2SStartProto fail")
		}

		return &MsgData{ModuleID: 32, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_collect
		return &MsgData{ModuleID: 32, SequenceID: 4}, nil

	case 7: //c2s_yuanbao_complete
		return &MsgData{ModuleID: 32, SequenceID: 7}, nil

	case 10: //c2s_yuanbao_refresh
		return &MsgData{ModuleID: 32, SequenceID: 10}, nil

	case 14: //c2s_vip_collect
		p := &pb32.C2SVipCollectProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhengwu.c2s_vip_collect UnmarshalMsgProto &C2SVipCollectProto fail")
		}

		return &MsgData{ModuleID: 32, SequenceID: 14, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("zhengwu收到未知消息: %d", sequenceID))
	}
}

func unmarshal_zhanjiang(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_open
		p := &pb33.C2SOpenProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhanjiang.c2s_open UnmarshalMsgProto &C2SOpenProto fail")
		}

		return &MsgData{ModuleID: 33, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_give_up
		return &MsgData{ModuleID: 33, SequenceID: 4}, nil

	case 7: //c2s_update_captain
		p := &pb33.C2SUpdateCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhanjiang.c2s_update_captain UnmarshalMsgProto &C2SUpdateCaptainProto fail")
		}

		return &MsgData{ModuleID: 33, SequenceID: 7, Proto: p}, nil

	case 10: //c2s_challenge
		p := &pb33.C2SChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhanjiang.c2s_challenge UnmarshalMsgProto &C2SChallengeProto fail")
		}

		return &MsgData{ModuleID: 33, SequenceID: 10, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("zhanjiang收到未知消息: %d", sequenceID))
	}
}

func unmarshal_question(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_start
		p := &pb34.C2SStartProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "question.c2s_start UnmarshalMsgProto &C2SStartProto fail")
		}

		return &MsgData{ModuleID: 34, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_answer
		p := &pb34.C2SAnswerProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "question.c2s_answer UnmarshalMsgProto &C2SAnswerProto fail")
		}

		return &MsgData{ModuleID: 34, SequenceID: 4, Proto: p}, nil

	case 6: //c2s_next
		p := &pb34.C2SNextProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "question.c2s_next UnmarshalMsgProto &C2SNextProto fail")
		}

		return &MsgData{ModuleID: 34, SequenceID: 6, Proto: p}, nil

	case 9: //c2s_get_prize
		p := &pb34.C2SGetPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "question.c2s_get_prize UnmarshalMsgProto &C2SGetPrizeProto fail")
		}

		return &MsgData{ModuleID: 34, SequenceID: 9, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("question收到未知消息: %d", sequenceID))
	}
}

func unmarshal_relation(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_add_relation
		p := &pb35.C2SAddRelationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.c2s_add_relation UnmarshalMsgProto &C2SAddRelationProto fail")
		}

		return &MsgData{ModuleID: 35, SequenceID: 1, Proto: p}, nil

	case 10: //c2s_remove_enemy
		p := &pb35.C2SRemoveEnemyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.c2s_remove_enemy UnmarshalMsgProto &C2SRemoveEnemyProto fail")
		}

		return &MsgData{ModuleID: 35, SequenceID: 10, Proto: p}, nil

	case 4: //c2s_remove_relation
		p := &pb35.C2SRemoveRelationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.c2s_remove_relation UnmarshalMsgProto &C2SRemoveRelationProto fail")
		}

		return &MsgData{ModuleID: 35, SequenceID: 4, Proto: p}, nil

	case 7: //c2s_list_relation
		p := &pb35.C2SListRelationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.c2s_list_relation UnmarshalMsgProto &C2SListRelationProto fail")
		}

		return &MsgData{ModuleID: 35, SequenceID: 7, Proto: p}, nil

	case 28: //c2s_new_list_relation
		p := &pb35.C2SNewListRelationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.c2s_new_list_relation UnmarshalMsgProto &C2SNewListRelationProto fail")
		}

		return &MsgData{ModuleID: 35, SequenceID: 28, Proto: p}, nil

	case 16: //c2s_recommend_hero_list
		p := &pb35.C2SRecommendHeroListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.c2s_recommend_hero_list UnmarshalMsgProto &C2SRecommendHeroListProto fail")
		}

		return &MsgData{ModuleID: 35, SequenceID: 16, Proto: p}, nil

	case 22: //c2s_search_heros
		p := &pb35.C2SSearchHerosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.c2s_search_heros UnmarshalMsgProto &C2SSearchHerosProto fail")
		}

		return &MsgData{ModuleID: 35, SequenceID: 22, Proto: p}, nil

	case 25: //c2s_search_hero_by_id
		p := &pb35.C2SSearchHeroByIdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.c2s_search_hero_by_id UnmarshalMsgProto &C2SSearchHeroByIdProto fail")
		}

		return &MsgData{ModuleID: 35, SequenceID: 25, Proto: p}, nil

	case 33: //c2s_set_important_friend
		p := &pb35.C2SSetImportantFriendProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.c2s_set_important_friend UnmarshalMsgProto &C2SSetImportantFriendProto fail")
		}

		return &MsgData{ModuleID: 35, SequenceID: 33, Proto: p}, nil

	case 36: //c2s_cancel_important_friend
		p := &pb35.C2SCancelImportantFriendProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.c2s_cancel_important_friend UnmarshalMsgProto &C2SCancelImportantFriendProto fail")
		}

		return &MsgData{ModuleID: 35, SequenceID: 36, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("relation收到未知消息: %d", sequenceID))
	}
}

func unmarshal_xiongnu(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_set_defender
		p := &pb36.C2SSetDefenderProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.c2s_set_defender UnmarshalMsgProto &C2SSetDefenderProto fail")
		}

		return &MsgData{ModuleID: 36, SequenceID: 1, Proto: p}, nil

	case 5: //c2s_start
		p := &pb36.C2SStartProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.c2s_start UnmarshalMsgProto &C2SStartProto fail")
		}

		return &MsgData{ModuleID: 36, SequenceID: 5, Proto: p}, nil

	case 10: //c2s_troop_info
		return &MsgData{ModuleID: 36, SequenceID: 10}, nil

	case 14: //c2s_get_xiong_nu_npc_base_info
		p := &pb36.C2SGetXiongNuNpcBaseInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.c2s_get_xiong_nu_npc_base_info UnmarshalMsgProto &C2SGetXiongNuNpcBaseInfoProto fail")
		}

		return &MsgData{ModuleID: 36, SequenceID: 14, Proto: p}, nil

	case 17: //c2s_get_defenser_fight_amount
		p := &pb36.C2SGetDefenserFightAmountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.c2s_get_defenser_fight_amount UnmarshalMsgProto &C2SGetDefenserFightAmountProto fail")
		}

		return &MsgData{ModuleID: 36, SequenceID: 17, Proto: p}, nil

	case 19: //c2s_get_xiong_nu_fight_info
		return &MsgData{ModuleID: 36, SequenceID: 19}, nil

	default:
		return nil, errors.New(fmt.Sprintf("xiongnu收到未知消息: %d", sequenceID))
	}
}

func unmarshal_survey(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 2: //c2s_complete
		p := &pb37.C2SCompleteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "survey.c2s_complete UnmarshalMsgProto &C2SCompleteProto fail")
		}

		return &MsgData{ModuleID: 37, SequenceID: 2, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("survey收到未知消息: %d", sequenceID))
	}
}

func unmarshal_farm(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 2: //c2s_plant
		p := &pb38.C2SPlantProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.c2s_plant UnmarshalMsgProto &C2SPlantProto fail")
		}

		return &MsgData{ModuleID: 38, SequenceID: 2, Proto: p}, nil

	case 5: //c2s_harvest
		p := &pb38.C2SHarvestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.c2s_harvest UnmarshalMsgProto &C2SHarvestProto fail")
		}

		return &MsgData{ModuleID: 38, SequenceID: 5, Proto: p}, nil

	case 8: //c2s_change
		p := &pb38.C2SChangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.c2s_change UnmarshalMsgProto &C2SChangeProto fail")
		}

		return &MsgData{ModuleID: 38, SequenceID: 8, Proto: p}, nil

	case 12: //c2s_one_key_plant
		p := &pb38.C2SOneKeyPlantProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.c2s_one_key_plant UnmarshalMsgProto &C2SOneKeyPlantProto fail")
		}

		return &MsgData{ModuleID: 38, SequenceID: 12, Proto: p}, nil

	case 28: //c2s_one_key_harvest
		p := &pb38.C2SOneKeyHarvestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.c2s_one_key_harvest UnmarshalMsgProto &C2SOneKeyHarvestProto fail")
		}

		return &MsgData{ModuleID: 38, SequenceID: 28, Proto: p}, nil

	case 52: //c2s_one_key_reset
		return &MsgData{ModuleID: 38, SequenceID: 52}, nil

	case 43: //c2s_view_farm
		p := &pb38.C2SViewFarmProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.c2s_view_farm UnmarshalMsgProto &C2SViewFarmProto fail")
		}

		return &MsgData{ModuleID: 38, SequenceID: 43, Proto: p}, nil

	case 18: //c2s_steal
		p := &pb38.C2SStealProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.c2s_steal UnmarshalMsgProto &C2SStealProto fail")
		}

		return &MsgData{ModuleID: 38, SequenceID: 18, Proto: p}, nil

	case 31: //c2s_one_key_steal
		p := &pb38.C2SOneKeyStealProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.c2s_one_key_steal UnmarshalMsgProto &C2SOneKeyStealProto fail")
		}

		return &MsgData{ModuleID: 38, SequenceID: 31, Proto: p}, nil

	case 39: //c2s_steal_log_list
		p := &pb38.C2SStealLogListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.c2s_steal_log_list UnmarshalMsgProto &C2SStealLogListProto fail")
		}

		return &MsgData{ModuleID: 38, SequenceID: 39, Proto: p}, nil

	case 48: //c2s_can_steal_list
		return &MsgData{ModuleID: 38, SequenceID: 48}, nil

	default:
		return nil, errors.New(fmt.Sprintf("farm收到未知消息: %d", sequenceID))
	}
}

func unmarshal_dianquan(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_exchange
		p := &pb39.C2SExchangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dianquan.c2s_exchange UnmarshalMsgProto &C2SExchangeProto fail")
		}

		return &MsgData{ModuleID: 39, SequenceID: 1, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("dianquan收到未知消息: %d", sequenceID))
	}
}

func unmarshal_xuanyuan(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_self_info
		return &MsgData{ModuleID: 40, SequenceID: 1}, nil

	case 11: //c2s_list_target
		p := &pb40.C2SListTargetProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.c2s_list_target UnmarshalMsgProto &C2SListTargetProto fail")
		}

		return &MsgData{ModuleID: 40, SequenceID: 11, Proto: p}, nil

	case 5: //c2s_query_target_troop
		p := &pb40.C2SQueryTargetTroopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.c2s_query_target_troop UnmarshalMsgProto &C2SQueryTargetTroopProto fail")
		}

		return &MsgData{ModuleID: 40, SequenceID: 5, Proto: p}, nil

	case 15: //c2s_challenge
		p := &pb40.C2SChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.c2s_challenge UnmarshalMsgProto &C2SChallengeProto fail")
		}

		return &MsgData{ModuleID: 40, SequenceID: 15, Proto: p}, nil

	case 20: //c2s_list_record
		p := &pb40.C2SListRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.c2s_list_record UnmarshalMsgProto &C2SListRecordProto fail")
		}

		return &MsgData{ModuleID: 40, SequenceID: 20, Proto: p}, nil

	case 22: //c2s_collect_rank_prize
		return &MsgData{ModuleID: 40, SequenceID: 22}, nil

	default:
		return nil, errors.New(fmt.Sprintf("xuanyuan收到未知消息: %d", sequenceID))
	}
}

func unmarshal_hebi(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_room_list
		p := &pb41.C2SRoomListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.c2s_room_list UnmarshalMsgProto &C2SRoomListProto fail")
		}

		return &MsgData{ModuleID: 41, SequenceID: 1, Proto: p}, nil

	case 35: //c2s_hero_record_list
		return &MsgData{ModuleID: 41, SequenceID: 35}, nil

	case 3: //c2s_change_captain
		p := &pb41.C2SChangeCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.c2s_change_captain UnmarshalMsgProto &C2SChangeCaptainProto fail")
		}

		return &MsgData{ModuleID: 41, SequenceID: 3, Proto: p}, nil

	case 28: //c2s_check_in_room
		p := &pb41.C2SCheckInRoomProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.c2s_check_in_room UnmarshalMsgProto &C2SCheckInRoomProto fail")
		}

		return &MsgData{ModuleID: 41, SequenceID: 28, Proto: p}, nil

	case 31: //c2s_copy_self
		p := &pb41.C2SCopySelfProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.c2s_copy_self UnmarshalMsgProto &C2SCopySelfProto fail")
		}

		return &MsgData{ModuleID: 41, SequenceID: 31, Proto: p}, nil

	case 9: //c2s_join_room
		p := &pb41.C2SJoinRoomProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.c2s_join_room UnmarshalMsgProto &C2SJoinRoomProto fail")
		}

		return &MsgData{ModuleID: 41, SequenceID: 9, Proto: p}, nil

	case 12: //c2s_rob_pos
		p := &pb41.C2SRobPosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.c2s_rob_pos UnmarshalMsgProto &C2SRobPosProto fail")
		}

		return &MsgData{ModuleID: 41, SequenceID: 12, Proto: p}, nil

	case 18: //c2s_leave_room
		p := &pb41.C2SLeaveRoomProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.c2s_leave_room UnmarshalMsgProto &C2SLeaveRoomProto fail")
		}

		return &MsgData{ModuleID: 41, SequenceID: 18, Proto: p}, nil

	case 21: //c2s_rob
		p := &pb41.C2SRobProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.c2s_rob UnmarshalMsgProto &C2SRobProto fail")
		}

		return &MsgData{ModuleID: 41, SequenceID: 21, Proto: p}, nil

	case 37: //c2s_view_show_prize
		p := &pb41.C2SViewShowPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.c2s_view_show_prize UnmarshalMsgProto &C2SViewShowPrizeProto fail")
		}

		return &MsgData{ModuleID: 41, SequenceID: 37, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("hebi收到未知消息: %d", sequenceID))
	}
}

func unmarshal_mingc(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 4: //c2s_mingc_list
		p := &pb42.C2SMingcListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc.c2s_mingc_list UnmarshalMsgProto &C2SMingcListProto fail")
		}

		return &MsgData{ModuleID: 42, SequenceID: 4, Proto: p}, nil

	case 7: //c2s_view_mingc
		p := &pb42.C2SViewMingcProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc.c2s_view_mingc UnmarshalMsgProto &C2SViewMingcProto fail")
		}

		return &MsgData{ModuleID: 42, SequenceID: 7, Proto: p}, nil

	case 10: //c2s_mc_build
		p := &pb42.C2SMcBuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc.c2s_mc_build UnmarshalMsgProto &C2SMcBuildProto fail")
		}

		return &MsgData{ModuleID: 42, SequenceID: 10, Proto: p}, nil

	case 13: //c2s_mc_build_log
		p := &pb42.C2SMcBuildLogProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc.c2s_mc_build_log UnmarshalMsgProto &C2SMcBuildLogProto fail")
		}

		return &MsgData{ModuleID: 42, SequenceID: 13, Proto: p}, nil

	case 20: //c2s_mingc_host_guild
		p := &pb42.C2SMingcHostGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc.c2s_mingc_host_guild UnmarshalMsgProto &C2SMingcHostGuildProto fail")
		}

		return &MsgData{ModuleID: 42, SequenceID: 20, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("mingc收到未知消息: %d", sequenceID))
	}
}

func unmarshal_promotion(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 4: //c2s_collect_login_day_prize
		p := &pb43.C2SCollectLoginDayPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.c2s_collect_login_day_prize UnmarshalMsgProto &C2SCollectLoginDayPrizeProto fail")
		}

		return &MsgData{ModuleID: 43, SequenceID: 4, Proto: p}, nil

	case 7: //c2s_buy_level_fund
		return &MsgData{ModuleID: 43, SequenceID: 7}, nil

	case 10: //c2s_collect_level_fund
		p := &pb43.C2SCollectLevelFundProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.c2s_collect_level_fund UnmarshalMsgProto &C2SCollectLevelFundProto fail")
		}

		return &MsgData{ModuleID: 43, SequenceID: 10, Proto: p}, nil

	case 13: //c2s_collect_daily_sp
		p := &pb43.C2SCollectDailySpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.c2s_collect_daily_sp UnmarshalMsgProto &C2SCollectDailySpProto fail")
		}

		return &MsgData{ModuleID: 43, SequenceID: 13, Proto: p}, nil

	case 16: //c2s_collect_free_gift
		p := &pb43.C2SCollectFreeGiftProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.c2s_collect_free_gift UnmarshalMsgProto &C2SCollectFreeGiftProto fail")
		}

		return &MsgData{ModuleID: 43, SequenceID: 16, Proto: p}, nil

	case 21: //c2s_buy_time_limit_gift
		p := &pb43.C2SBuyTimeLimitGiftProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.c2s_buy_time_limit_gift UnmarshalMsgProto &C2SBuyTimeLimitGiftProto fail")
		}

		return &MsgData{ModuleID: 43, SequenceID: 21, Proto: p}, nil

	case 26: //c2s_buy_event_limit_gift
		p := &pb43.C2SBuyEventLimitGiftProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.c2s_buy_event_limit_gift UnmarshalMsgProto &C2SBuyEventLimitGiftProto fail")
		}

		return &MsgData{ModuleID: 43, SequenceID: 26, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("promotion收到未知消息: %d", sequenceID))
	}
}

func unmarshal_mingc_war(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 31: //c2s_view_mc_war_self_guild
		return &MsgData{ModuleID: 44, SequenceID: 31}, nil

	case 29: //c2s_view_mc_war
		p := &pb44.C2SViewMcWarProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_view_mc_war UnmarshalMsgProto &C2SViewMcWarProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 29, Proto: p}, nil

	case 16: //c2s_apply_atk
		p := &pb44.C2SApplyAtkProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_apply_atk UnmarshalMsgProto &C2SApplyAtkProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 16, Proto: p}, nil

	case 21: //c2s_apply_ast
		p := &pb44.C2SApplyAstProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_apply_ast UnmarshalMsgProto &C2SApplyAstProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 21, Proto: p}, nil

	case 80: //c2s_cancel_apply_ast
		p := &pb44.C2SCancelApplyAstProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_cancel_apply_ast UnmarshalMsgProto &C2SCancelApplyAstProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 80, Proto: p}, nil

	case 25: //c2s_reply_apply_ast
		p := &pb44.C2SReplyApplyAstProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_reply_apply_ast UnmarshalMsgProto &C2SReplyApplyAstProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 25, Proto: p}, nil

	case 75: //c2s_view_mingc_war_mc
		p := &pb44.C2SViewMingcWarMcProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_view_mingc_war_mc UnmarshalMsgProto &C2SViewMingcWarMcProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 75, Proto: p}, nil

	case 35: //c2s_join_fight
		p := &pb44.C2SJoinFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_join_fight UnmarshalMsgProto &C2SJoinFightProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 35, Proto: p}, nil

	case 38: //c2s_quit_fight
		return &MsgData{ModuleID: 44, SequenceID: 38}, nil

	case 49: //c2s_scene_move
		p := &pb44.C2SSceneMoveProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_scene_move UnmarshalMsgProto &C2SSceneMoveProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 49, Proto: p}, nil

	case 85: //c2s_scene_back
		return &MsgData{ModuleID: 44, SequenceID: 85}, nil

	case 88: //c2s_scene_speed_up
		p := &pb44.C2SSceneSpeedUpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_scene_speed_up UnmarshalMsgProto &C2SSceneSpeedUpProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 88, Proto: p}, nil

	case 72: //c2s_scene_troop_relive
		return &MsgData{ModuleID: 44, SequenceID: 72}, nil

	case 46: //c2s_view_mc_war_scene
		p := &pb44.C2SViewMcWarSceneProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_view_mc_war_scene UnmarshalMsgProto &C2SViewMcWarSceneProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 46, Proto: p}, nil

	case 139: //c2s_watch
		p := &pb44.C2SWatchProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_watch UnmarshalMsgProto &C2SWatchProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 139, Proto: p}, nil

	case 136: //c2s_quit_watch
		p := &pb44.C2SQuitWatchProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_quit_watch UnmarshalMsgProto &C2SQuitWatchProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 136, Proto: p}, nil

	case 91: //c2s_view_mc_war_record
		p := &pb44.C2SViewMcWarRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_view_mc_war_record UnmarshalMsgProto &C2SViewMcWarRecordProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 91, Proto: p}, nil

	case 94: //c2s_view_mc_war_troop_record
		p := &pb44.C2SViewMcWarTroopRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_view_mc_war_troop_record UnmarshalMsgProto &C2SViewMcWarTroopRecordProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 94, Proto: p}, nil

	case 99: //c2s_view_scene_troop_record
		return &MsgData{ModuleID: 44, SequenceID: 99}, nil

	case 107: //c2s_apply_refresh_rank
		p := &pb44.C2SApplyRefreshRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_apply_refresh_rank UnmarshalMsgProto &C2SApplyRefreshRankProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 107, Proto: p}, nil

	case 111: //c2s_view_my_guild_member_rank
		p := &pb44.C2SViewMyGuildMemberRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_view_my_guild_member_rank UnmarshalMsgProto &C2SViewMyGuildMemberRankProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 111, Proto: p}, nil

	case 115: //c2s_scene_change_mode
		p := &pb44.C2SSceneChangeModeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_scene_change_mode UnmarshalMsgProto &C2SSceneChangeModeProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 115, Proto: p}, nil

	case 119: //c2s_scene_tou_shi_building_turn_to
		p := &pb44.C2SSceneTouShiBuildingTurnToProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_scene_tou_shi_building_turn_to UnmarshalMsgProto &C2SSceneTouShiBuildingTurnToProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 119, Proto: p}, nil

	case 123: //c2s_scene_tou_shi_building_fire
		p := &pb44.C2SSceneTouShiBuildingFireProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.c2s_scene_tou_shi_building_fire UnmarshalMsgProto &C2SSceneTouShiBuildingFireProto fail")
		}

		return &MsgData{ModuleID: 44, SequenceID: 123, Proto: p}, nil

	case 128: //c2s_scene_drum
		return &MsgData{ModuleID: 44, SequenceID: 128}, nil

	default:
		return nil, errors.New(fmt.Sprintf("mingc_war收到未知消息: %d", sequenceID))
	}
}

func unmarshal_random_event(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_choose_option
		p := &pb45.C2SChooseOptionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "random_event.c2s_choose_option UnmarshalMsgProto &C2SChooseOptionProto fail")
		}

		return &MsgData{ModuleID: 45, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_open_event
		p := &pb45.C2SOpenEventProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "random_event.c2s_open_event UnmarshalMsgProto &C2SOpenEventProto fail")
		}

		return &MsgData{ModuleID: 45, SequenceID: 4, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("random_event收到未知消息: %d", sequenceID))
	}
}

func unmarshal_strategy(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_use_stratagem
		p := &pb46.C2SUseStratagemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "strategy.c2s_use_stratagem UnmarshalMsgProto &C2SUseStratagemProto fail")
		}

		return &MsgData{ModuleID: 46, SequenceID: 1, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("strategy收到未知消息: %d", sequenceID))
	}
}

func unmarshal_vip(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 4: //c2s_vip_collect_daily_prize
		p := &pb48.C2SVipCollectDailyPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "vip.c2s_vip_collect_daily_prize UnmarshalMsgProto &C2SVipCollectDailyPrizeProto fail")
		}

		return &MsgData{ModuleID: 48, SequenceID: 4, Proto: p}, nil

	case 7: //c2s_vip_collect_level_prize
		p := &pb48.C2SVipCollectLevelPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "vip.c2s_vip_collect_level_prize UnmarshalMsgProto &C2SVipCollectLevelPrizeProto fail")
		}

		return &MsgData{ModuleID: 48, SequenceID: 7, Proto: p}, nil

	case 12: //c2s_vip_buy_dungeon_times
		p := &pb48.C2SVipBuyDungeonTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "vip.c2s_vip_buy_dungeon_times UnmarshalMsgProto &C2SVipBuyDungeonTimesProto fail")
		}

		return &MsgData{ModuleID: 48, SequenceID: 12, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("vip收到未知消息: %d", sequenceID))
	}
}

func unmarshal_red_packet(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_buy
		p := &pb49.C2SBuyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "red_packet.c2s_buy UnmarshalMsgProto &C2SBuyProto fail")
		}

		return &MsgData{ModuleID: 49, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_create
		p := &pb49.C2SCreateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "red_packet.c2s_create UnmarshalMsgProto &C2SCreateProto fail")
		}

		return &MsgData{ModuleID: 49, SequenceID: 4, Proto: p}, nil

	case 7: //c2s_grab
		p := &pb49.C2SGrabProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "red_packet.c2s_grab UnmarshalMsgProto &C2SGrabProto fail")
		}

		return &MsgData{ModuleID: 49, SequenceID: 7, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("red_packet收到未知消息: %d", sequenceID))
	}
}

func unmarshal_teach(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_fight
		p := &pb50.C2SFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "teach.c2s_fight UnmarshalMsgProto &C2SFightProto fail")
		}

		return &MsgData{ModuleID: 50, SequenceID: 1, Proto: p}, nil

	case 4: //c2s_collect_prize
		p := &pb50.C2SCollectPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "teach.c2s_collect_prize UnmarshalMsgProto &C2SCollectPrizeProto fail")
		}

		return &MsgData{ModuleID: 50, SequenceID: 4, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("teach收到未知消息: %d", sequenceID))
	}
}

func unmarshal_activity(sequenceID int, data []byte) (*MsgData, error) {
	switch sequenceID {

	case 1: //c2s_collect_collection
		p := &pb51.C2SCollectCollectionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "activity.c2s_collect_collection UnmarshalMsgProto &C2SCollectCollectionProto fail")
		}

		return &MsgData{ModuleID: 51, SequenceID: 1, Proto: p}, nil

	default:
		return nil, errors.New(fmt.Sprintf("activity收到未知消息: %d", sequenceID))
	}
}

// 在英雄线程处理收到的消息
func Handle(msg *MsgData, hc iface.HeroController) error {
	switch msg.ModuleID {
	case 1: // login
		return handle_login(msg.SequenceID, hc, msg.Proto)

	case 2: // domestic
		return handle_domestic(msg.SequenceID, hc, msg.Proto)

	case 3: // gm
		return handle_gm(msg.SequenceID, hc, msg.Proto)

	case 4: // military
		return handle_military(msg.SequenceID, hc, msg.Proto)

	case 5: // misc
		return handle_misc(msg.SequenceID, hc, msg.Proto)

	case 7: // region
		return handle_region(msg.SequenceID, hc, msg.Proto)

	case 8: // mail
		return handle_mail(msg.SequenceID, hc, msg.Proto)

	case 9: // guild
		return handle_guild(msg.SequenceID, hc, msg.Proto)

	case 10: // stress
		return handle_stress(msg.SequenceID, hc, msg.Proto)

	case 11: // depot
		return handle_depot(msg.SequenceID, hc, msg.Proto)

	case 12: // equipment
		return handle_equipment(msg.SequenceID, hc, msg.Proto)

	case 13: // chat
		return handle_chat(msg.SequenceID, hc, msg.Proto)

	case 14: // tower
		return handle_tower(msg.SequenceID, hc, msg.Proto)

	case 15: // task
		return handle_task(msg.SequenceID, hc, msg.Proto)

	case 16: // fishing
		return handle_fishing(msg.SequenceID, hc, msg.Proto)

	case 19: // gem
		return handle_gem(msg.SequenceID, hc, msg.Proto)

	case 20: // shop
		return handle_shop(msg.SequenceID, hc, msg.Proto)

	case 21: // client_config
		return handle_client_config(msg.SequenceID, hc, msg.Proto)

	case 22: // secret_tower
		return handle_secret_tower(msg.SequenceID, hc, msg.Proto)

	case 23: // rank
		return handle_rank(msg.SequenceID, hc, msg.Proto)

	case 24: // bai_zhan
		return handle_bai_zhan(msg.SequenceID, hc, msg.Proto)

	case 26: // dungeon
		return handle_dungeon(msg.SequenceID, hc, msg.Proto)

	case 27: // country
		return handle_country(msg.SequenceID, hc, msg.Proto)

	case 29: // tag
		return handle_tag(msg.SequenceID, hc, msg.Proto)

	case 31: // garden
		return handle_garden(msg.SequenceID, hc, msg.Proto)

	case 32: // zhengwu
		return handle_zhengwu(msg.SequenceID, hc, msg.Proto)

	case 33: // zhanjiang
		return handle_zhanjiang(msg.SequenceID, hc, msg.Proto)

	case 34: // question
		return handle_question(msg.SequenceID, hc, msg.Proto)

	case 35: // relation
		return handle_relation(msg.SequenceID, hc, msg.Proto)

	case 36: // xiongnu
		return handle_xiongnu(msg.SequenceID, hc, msg.Proto)

	case 37: // survey
		return handle_survey(msg.SequenceID, hc, msg.Proto)

	case 38: // farm
		return handle_farm(msg.SequenceID, hc, msg.Proto)

	case 39: // dianquan
		return handle_dianquan(msg.SequenceID, hc, msg.Proto)

	case 40: // xuanyuan
		return handle_xuanyuan(msg.SequenceID, hc, msg.Proto)

	case 41: // hebi
		return handle_hebi(msg.SequenceID, hc, msg.Proto)

	case 42: // mingc
		return handle_mingc(msg.SequenceID, hc, msg.Proto)

	case 43: // promotion
		return handle_promotion(msg.SequenceID, hc, msg.Proto)

	case 44: // mingc_war
		return handle_mingc_war(msg.SequenceID, hc, msg.Proto)

	case 45: // random_event
		return handle_random_event(msg.SequenceID, hc, msg.Proto)

	case 46: // strategy
		return handle_strategy(msg.SequenceID, hc, msg.Proto)

	case 48: // vip
		return handle_vip(msg.SequenceID, hc, msg.Proto)

	case 49: // red_packet
		return handle_red_packet(msg.SequenceID, hc, msg.Proto)

	case 50: // teach
		return handle_teach(msg.SequenceID, hc, msg.Proto)

	case 51: // activity
		return handle_activity(msg.SequenceID, hc, msg.Proto)

	default:
		return errors.New(fmt.Sprintf("消息没有模块处理. 模块号: %d", msg.ModuleID))
	}
}

func handle_login(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	default:
		return errors.New(fmt.Sprintf("模块login消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_domestic(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		domesticModule.ProcessCreateResourceBuilding(proto.(*pb2.C2SCreateBuildingProto), hc)

	case 4:
		domesticModule.ProcessUpgradeResourceBuilding(proto.(*pb2.C2SUpgradeBuildingProto), hc)

	case 7:
		domesticModule.ProcessRebuildBuilding(proto.(*pb2.C2SRebuildResourceBuildingProto), hc)

	case 15:
		domesticModule.ProcessCollectResource(proto.(*pb2.C2SCollectResourceProto), hc)

	case 18:
		domesticModule.ProcessLearnTechnology(proto.(*pb2.C2SLearnTechnologyProto), hc)

	case 24:
		domesticModule.ProcessUpgradeStableBuilding(proto.(*pb2.C2SUpgradeStableBuildingProto), hc)

	case 30:
		domesticModule.ProcessChangeHeroName(proto.(*pb2.C2SChangeHeroNameProto), hc)

	case 33:
		domesticModule.ProcessListOldName(proto.(*pb2.C2SListOldNameProto), hc)

	case 35:
		domesticModule.ProcessViewOtherHero(proto.(*pb2.C2SViewOtherHeroProto), hc)

	case 41:
		domesticModule.ProcessMiaoBuildingWorkerCd(proto.(*pb2.C2SMiaoBuildingWorkerCdProto), hc)

	case 44:
		domesticModule.ProcessMiaoTechWorkerCd(proto.(*pb2.C2SMiaoTechWorkerCdProto), hc)

	case 51:
		domesticModule.ProcessForgingEquip(proto.(*pb2.C2SForgingEquipProto), hc)

	case 60:
		domesticModule.ProcessRequestCityExchangeEvent(hc)

	case 63:
		domesticModule.ProcessCityEventExchange(proto.(*pb2.C2SCityEventExchangeProto), hc)

	case 66:
		domesticModule.ProcessSign(proto.(*pb2.C2SSignProto), hc)

	case 69:
		domesticModule.ProcessVoice(proto.(*pb2.C2SVoiceProto), hc)

	case 76:
		domesticModule.ProcessCollectResourceV2(proto.(*pb2.C2SCollectResourceV2Proto), hc)

	case 81:
		domesticModule.ProcessRequestResourceConflict(hc)

	case 87:
		domesticModule.ProcessUnlockStableBuilding(proto.(*pb2.C2SUnlockStableBuildingProto), hc)

	case 90:
		domesticModule.ProcessIsHeroNameExist(proto.(*pb2.C2SIsHeroNameExistProto), hc)

	case 94:
		domesticModule.ProcessChangeHead(proto.(*pb2.C2SChangeHeadProto), hc)

	case 108:
		domesticModule.ProcessUnlockOuterCity(proto.(*pb2.C2SUnlockOuterCityProto), hc)

	case 111:
		domesticModule.ProcessUpgradeOuterCityBuilding(proto.(*pb2.C2SUpgradeOuterCityBuildingProto), hc)

	case 114:
		domesticModule.ProcessCollectCountdownPrize(hc)

	case 119:
		domesticModule.ProcessWorkshopStartForge(proto.(*pb2.C2SStartWorkshopProto), hc)

	case 122:
		domesticModule.ProcessWorkshopCollect(proto.(*pb2.C2SCollectWorkshopProto), hc)

	case 125:
		domesticModule.ProcessViewFightInfo(proto.(*pb2.C2SViewFightInfoProto), hc)

	case 127:
		domesticModule.ProcessWorkshopMiaoCd(proto.(*pb2.C2SWorkshopMiaoCdProto), hc)

	case 130:
		domesticModule.ProcessChangeBody(proto.(*pb2.C2SChangeBodyProto), hc)

	case 133:
		domesticModule.ProcessRefreshWorkshop(hc)

	case 136:
		domesticModule.ProcessCollectSeasonPrize(hc)

	case 142:
		domesticModule.ProcessUpdateOuterCityType(proto.(*pb2.C2SUpdateOuterCityTypeProto), hc)

	case 147:
		domesticModule.ProcessBuySp(proto.(*pb2.C2SBuySpProto), hc)

	case 150:
		domesticModule.ProcessUseBufEffect(proto.(*pb2.C2SUseBufEffectProto), hc)

	case 154:
		domesticModule.ProcessOpenBufEffectUi(hc)

	case 158:
		domesticModule.ProcessUseAdvantage(proto.(*pb2.C2SUseAdvantageProto), hc)

	case 166:
		domesticModule.ProcessWorkerUnlock(proto.(*pb2.C2SWorkerUnlockProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块domestic消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_gm(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		gmModule.ProcessGmMsg(proto.(*pb3.C2SGmProto), hc)

	case 5:
		gmModule.ProcessListCmdMsg(hc)

	case 7:
		gmModule.ProcessInvaseTargetIdMsg(proto.(*pb3.C2SInvaseTargetIdProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块gm消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_military(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 2:
		militaryModule.ProcessC2SRecruitSoldierMsg_deprecated(proto.(*pb4.C2SRecruitSoldierProto), hc)

	case 6:
		militaryModule.ProcessC2SHealWoundedSoldierMsg_deprecated(proto.(*pb4.C2SHealWoundedSoldierProto), hc)

	case 9:
		militaryModule.ProcessC2SCaptainChangeSoldierMsg(proto.(*pb4.C2SCaptainChangeSoldierProto), hc)

	case 12:
		militaryModule.ProcessC2SFightMsg(proto.(*pb4.C2SFightProto), hc)

	case 15:
		militaryModule.ProcessC2SUpgradeSoldierMsg(hc)

	case 34:
		militaryModule.ProcessSellSeekCaptain_deprecated(proto.(*pb4.C2SSellSeekCaptainProto), hc)

	case 38:
		militaryModule.ProcessFireCaptain_deprecated(proto.(*pb4.C2SFireCaptainProto), hc)

	case 45:
		militaryModule.ProcessSetMultiCaptainIndex(proto.(*pb4.C2SSetMultiCaptainIndexProto), hc)

	case 48:
		militaryModule.ProcessCaptainRefined_deprecated(proto.(*pb4.C2SCaptainRefinedProto), hc)

	case 66:
		militaryModule.ProcessC2SCaptainFullSoldierMsg(proto.(*pb4.C2SCaptainFullSoldierProto), hc)

	case 74:
		militaryModule.ProcessGetMaxRecruitSoldier_deprecated(hc)

	case 76:
		militaryModule.ProcessGetMaxHealSoldier_deprecated(hc)

	case 82:
		militaryModule.ProcessChangeCaptainName_deprecated(proto.(*pb4.C2SChangeCaptainNameProto), hc)

	case 85:
		militaryModule.ProcessChangeCaptainRace_deprecated(proto.(*pb4.C2SChangeCaptainRaceProto), hc)

	case 89:
		militaryModule.ProcessCaptainRevirthPreview(proto.(*pb4.C2SCaptainRebirthPreviewProto), hc)

	case 92:
		militaryModule.ProcessCaptainRebirth_deprecated(proto.(*pb4.C2SCaptainRebirthProto), hc)

	case 101:
		militaryModule.ProcessC2SMultiFightMsg(hc)

	case 106:
		militaryModule.ProcessSetDefenseTroop(proto.(*pb4.C2SSetDefenseTroopProto), hc)

	case 109:
		militaryModule.ProcessRecruitCaptainV2_deprecated(hc)

	case 112:
		militaryModule.ProcessJiuGuanConsult(hc)

	case 117:
		militaryModule.ProcessJiuGuanRefresh(proto.(*pb4.C2SJiuGuanRefreshProto), hc)

	case 120:
		militaryModule.ProcessC2SRecruitSoldierV2Msg_deprecated(proto.(*pb4.C2SRecruitSoldierV2Proto), hc)

	case 125:
		militaryModule.ProcessUnlockCaptainRestraintSpell_deprecated(proto.(*pb4.C2SUnlockCaptainRestraintSpellProto), hc)

	case 129:
		militaryModule.ProcessClearDefenseTroopDefeatedMail(proto.(*pb4.C2SClearDefenseTroopDefeatedMailProto), hc)

	case 132:
		militaryModule.ProcessGetCaptainStatDetails_deprecated(proto.(*pb4.C2SGetCaptainStatDetailsProto), hc)

	case 136:
		militaryModule.ProcessCollectCaptainTrainingExp(hc)

	case 140:
		militaryModule.ProcessUseTrainingExpGoods_deprecated(proto.(*pb4.C2SUseTrainingExpGoodsProto), hc)

	case 143:
		militaryModule.ProcessSetPveCaptain(proto.(*pb4.C2SSetPveCaptainProto), hc)

	case 146:
		militaryModule.ProcessRecruitCaptainSeeker_deprecated(proto.(*pb4.C2SRecruitCaptainSeekerProto), hc)

	case 149:
		militaryModule.ProcessForceAddSoldierMsg(hc)

	case 166:
		militaryModule.ProcessCaptainRebirthMiaoCd(proto.(*pb4.C2SCaptainRebirthMiaoCdProto), hc)

	case 169:
		militaryModule.ProcessUpdateCaptainOfficial_deprecated(proto.(*pb4.C2SUpdateCaptainOfficialProto), hc)

	case 172:
		militaryModule.ProcessLeaveCaptainOfficial_deprecated(proto.(*pb4.C2SLeaveCaptainOfficialProto), hc)

	case 176:
		militaryModule.ProcessRandomCaptainHead_deprecated(hc)

	case 181:
		militaryModule.ProcessCloseFightGuide(proto.(*pb4.C2SCloseFightGuideProto), hc)

	case 185:
		militaryModule.ProcessUseGongXunGoods(proto.(*pb4.C2SUseGongXunGoodsProto), hc)

	case 188:
		militaryModule.ProcessSetDefenserAutoFullSoldier(proto.(*pb4.C2SSetDefenserAutoFullSoldierProto), hc)

	case 190:
		militaryModule.ProcessViewOtherHeroCaptain(proto.(*pb4.C2SViewOtherHeroCaptainProto), hc)

	case 193:
		militaryModule.ProcessCopyDefenserGoods(proto.(*pb4.C2SUseCopyDefenserGoodsProto), hc)

	case 198:
		militaryModule.ProcessC2SFightxMsg(proto.(*pb4.C2SFightxProto), hc)

	case 206:
		militaryModule.ProcessCaptainEnhance(proto.(*pb4.C2SCaptainEnhanceProto), hc)

	case 210:
		militaryModule.ProcessCaptainProgress(proto.(*pb4.C2SCaptainProgressProto), hc)

	case 213:
		militaryModule.ProcessCaptainTrainExp(hc)

	case 216:
		militaryModule.ProcessUseLevelExpGoods_deprecated(proto.(*pb4.C2SUseLevelExpGoodsProto), hc)

	case 219:
		militaryModule.ProcessCaptainStatDetails(proto.(*pb4.C2SCaptainStatDetailsProto), hc)

	case 222:
		militaryModule.ProcessSetCaptainOfficial(proto.(*pb4.C2SSetCaptainOfficialProto), hc)

	case 228:
		militaryModule.ProcessUseGongxunGoods_deprecated(proto.(*pb4.C2SUseGongxunGoodsProto), hc)

	case 231:
		militaryModule.ProcessCaptainBorn(proto.(*pb4.C2SCaptainBornProto), hc)

	case 234:
		militaryModule.ProcessCaptainUpstar(proto.(*pb4.C2SCaptainUpstarProto), hc)

	case 243:
		militaryModule.ProcessUseLevelExpGoods2(proto.(*pb4.C2SUseLevelExpGoods2Proto), hc)

	case 252:
		militaryModule.ProcessNoticeCaptainHasViewed(hc)

	case 255:
		militaryModule.ProcessAutoUseGoodsUntilCaptainLevelup(proto.(*pb4.C2SAutoUseGoodsUntilCaptainLevelupProto), hc)

	case 261:
		militaryModule.ProcessCaptainCanCollectExp(hc)

	case 265:
		militaryModule.ProcessActivateCaptainFetter(proto.(*pb4.C2SActivateCaptainFriendshipProto), hc)

	case 268:
		militaryModule.ProcessCaptainExchange(proto.(*pb4.C2SCaptainExchangeProto), hc)

	case 272:
		militaryModule.ProcessNoticeOfficialHasViewed(proto.(*pb4.C2SNoticeOfficialHasViewedProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块military消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_misc(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		miscModule.ProcessHeartBeat(hc)

	case 3:
		miscModule.ProcessRequestConfig(proto.(*pb5.C2SConfigProto), hc)

	case 7:
		miscModule.ProcessClientLog(proto.(*pb5.C2SClientLogProto), hc)

	case 8:
		miscModule.ProcessSyncTime(proto.(*pb5.C2SSyncTimeProto), hc)

	case 10:
		miscModule.ProcessGetBlock(hc)

	case 15:
		miscModule.ProcessPing(hc)

	case 17:
		miscModule.ProcessClientVersion(proto.(*pb5.C2SClientVersionProto), hc)

	case 19:
		miscModule.ProcessUpdatePfToken(proto.(*pb5.C2SUpdatePfTokenProto), hc)

	case 21:
		miscModule.ProcessSettings(proto.(*pb5.C2SSettingsProto), hc)

	case 24:
		miscModule.ProcessSettingsToDefault(hc)

	case 31:
		miscModule.ProcessUpdateLocation(proto.(*pb5.C2SUpdateLocationProto), hc)

	case 35:
		miscModule.ProcessBackgroudHeartBeat(hc)

	case 36:
		miscModule.ProcessBackgroudWeakup(hc)

	case 42:
		miscModule.ProcessCollectChargePrize(proto.(*pb5.C2SCollectChargePrizeProto), hc)

	case 46:
		miscModule.ProcessCollectDailyBargain(proto.(*pb5.C2SCollectDailyBargainProto), hc)

	case 51:
		miscModule.ProcessActivateDurationCard(proto.(*pb5.C2SActivateDurationCardProto), hc)

	case 54:
		miscModule.ProcessCollectDurationCardDailyPrize(proto.(*pb5.C2SCollectDurationCardDailyPrizeProto), hc)

	case 64:
		miscModule.ProcessSetPrivacySetting(proto.(*pb5.C2SSetPrivacySettingProto), hc)

	case 67:
		miscModule.ProcessSetDefaultPrivacySettings(hc)

	case 69:
		miscModule.ProcessGetProductInfo(proto.(*pb5.C2SGetProductInfoProto), hc)

	case 76:
		miscModule.ProcessRequestLuaConfig(proto.(*pb5.C2SConfigluaProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块misc消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_region(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		regionModule.ProcessCreateBase(proto.(*pb7.C2SCreateBaseProto), hc)

	case 14:
		regionModule.ProcessFastMoveBase(proto.(*pb7.C2SFastMoveBaseProto), hc)

	case 24:
		regionModule.ProcessInvasion(proto.(*pb7.C2SInvasionProto), hc)

	case 26:
		regionModule.ProcessCancelInvasion(proto.(*pb7.C2SCancelInvasionProto), hc)

	case 30:
		regionModule.ProcessExpel(proto.(*pb7.C2SExpelProto), hc)

	case 40:
		regionModule.ProcessUpgradeBase(hc)

	case 46:
		regionModule.ProcessSwitchAction(proto.(*pb7.C2SSwitchActionProto), hc)

	case 71:
		regionModule.ProcessRepatriate(proto.(*pb7.C2SRepatriateProto), hc)

	case 86:
		regionModule.ProcessGetWhiteFlagDetail(proto.(*pb7.C2SWhiteFlagDetailProto), hc)

	case 91:
		regionModule.ProcessUseMianGoods(proto.(*pb7.C2SUseMianGoodsProto), hc)

	case 94:
		regionModule.ProcessPreInvasionTarget(proto.(*pb7.C2SPreInvasionTargetProto), hc)

	case 99:
		regionModule.ProcessFavoritePos(proto.(*pb7.C2SFavoritePosProto), hc)

	case 102:
		regionModule.ProcessFavoritePosList(hc)

	case 106:
		regionModule.ProcessBuyProsperity(hc)

	case 131:
		regionModule.ProcessRequestRuinsBase(proto.(*pb7.C2SRequestRuinsBaseProto), hc)

	case 139:
		regionModule.ProcessSpeedUp(proto.(*pb7.C2SSpeedUpProto), hc)

	case 142:
		regionModule.ProcessInvestigate(proto.(*pb7.C2SInvestigateProto), hc)

	case 148:
		regionModule.ProcessUpdateSelfView(proto.(*pb7.C2SUpdateSelfViewProto), hc)

	case 150:
		regionModule.ProcessCloseView(hc)

	case 156:
		regionModule.ProcessRequestTroopUnit(proto.(*pb7.C2SRequestTroopUnitProto), hc)

	case 161:
		regionModule.ProcessWatchBaseUnit(proto.(*pb7.C2SWatchBaseUnitProto), hc)

	case 165:
		regionModule.ProcessRequestMilitaryPush(proto.(*pb7.C2SRequestMilitaryPushProto), hc)

	case 172:
		regionModule.ProcessCalcMoveSpeed(proto.(*pb7.C2SCalcMoveSpeedProto), hc)

	case 175:
		regionModule.ProcessGetPrevInvestigate(proto.(*pb7.C2SGetPrevInvestigateProto), hc)

	case 178:
		regionModule.ProcessListEnemyPos(hc)

	case 180:
		regionModule.ProcessSearchBaozNpc(proto.(*pb7.C2SSearchBaozNpcProto), hc)

	case 183:
		regionModule.ProcessUseMultiLevelNpcTimesGoods(proto.(*pb7.C2SUseMultiLevelNpcTimesGoodsProto), hc)

	case 186:
		regionModule.ProcessBaozRepatriate(proto.(*pb7.C2SBaozRepatriateProto), hc)

	case 190:
		regionModule.ProcessUseInvaseHeroTimesGoods(proto.(*pb7.C2SUseInvaseHeroTimesGoodsProto), hc)

	case 193:
		regionModule.ProcessHomeAstDefendingInfo(hc)

	case 196:
		regionModule.ProcessGuildPleaseHelpMe(hc)

	case 199:
		regionModule.ProcessCreateAssembly(proto.(*pb7.C2SCreateAssemblyProto), hc)

	case 202:
		regionModule.ProcessShowAssembly(proto.(*pb7.C2SShowAssemblyProto), hc)

	case 206:
		regionModule.ProcessJoinAssembly(proto.(*pb7.C2SJoinAssemblyProto), hc)

	case 211:
		regionModule.ProcessGetBuyProsperityCost(hc)

	case 214:
		regionModule.ProcessCreateGuildWorkshop(proto.(*pb7.C2SCreateGuildWorkshopProto), hc)

	case 217:
		regionModule.ProcessShowGuildWorkshop(proto.(*pb7.C2SShowGuildWorkshopProto), hc)

	case 219:
		regionModule.ProcessHurtGuildWorkshop(proto.(*pb7.C2SHurtGuildWorkshopProto), hc)

	case 223:
		regionModule.ProcessRemoveGuildWorkshop(hc)

	case 228:
		regionModule.ProcessCatchGuildWorkshopLogs(proto.(*pb7.C2SCatchGuildWorkshopLogsProto), hc)

	case 232:
		regionModule.ProcessGetSelfBaoz(hc)

	case 234:
		regionModule.ProcessInvestigateInvade(proto.(*pb7.C2SInvestigateInvadeProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块region消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_mail(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		mailModule.ListMail(proto.(*pb8.C2SListMailProto), hc)

	case 8:
		mailModule.DeleteMail(proto.(*pb8.C2SDeleteMailProto), hc)

	case 11:
		mailModule.KeepMail(proto.(*pb8.C2SKeepMailProto), hc)

	case 14:
		mailModule.ProcessCollectMailPrize(proto.(*pb8.C2SCollectMailPrizeProto), hc)

	case 20:
		mailModule.ReadMail(proto.(*pb8.C2SReadMailProto), hc)

	case 24:
		mailModule.ProcessReadMulti(proto.(*pb8.C2SReadMultiProto), hc)

	case 26:
		mailModule.ProcessDeleteMulti(proto.(*pb8.C2SDeleteMultiProto), hc)

	case 28:
		mailModule.ProcessGetMail(proto.(*pb8.C2SGetMailProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块mail消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_guild(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		guildModule.ProcessListGuild(hc)

	case 4:
		guildModule.ProcessSearchGuild(proto.(*pb9.C2SSearchGuildProto), hc)

	case 7:
		guildModule.ProcessCreateGuild(proto.(*pb9.C2SCreateGuildProto), hc)

	case 10:
		guildModule.ProcessSelfGuild(proto.(*pb9.C2SSelfGuildProto), hc)

	case 13:
		guildModule.ProcessLeaveGuild(hc)

	case 17:
		guildModule.ProcessKickOther(proto.(*pb9.C2SKickOtherProto), hc)

	case 20:
		guildModule.ProcessUpdateText(proto.(*pb9.C2SUpdateTextProto), hc)

	case 23:
		guildModule.ProcessUpdateClassNames(proto.(*pb9.C2SUpdateClassNamesProto), hc)

	case 26:
		guildModule.ProcessUpdateFlagType(proto.(*pb9.C2SUpdateFlagTypeProto), hc)

	case 29:
		guildModule.ProcessUpdateMemberClassLevel(proto.(*pb9.C2SUpdateMemberClassLevelProto), hc)

	case 40:
		guildModule.ProcessUserRequestJoin(proto.(*pb9.C2SUserRequestJoinProto), hc)

	case 43:
		guildModule.ProcessUserCancelJoinRequest(proto.(*pb9.C2SUserCancelJoinRequestProto), hc)

	case 48:
		guildModule.ProcessUserReplyInvateRequest(proto.(*pb9.C2SUserReplyInvateRequestProto), hc)

	case 55:
		guildModule.ProcessGuildReplyJoinRequest(proto.(*pb9.C2SGuildReplyJoinRequestProto), hc)

	case 65:
		guildModule.ProcessUpdateInternalText(proto.(*pb9.C2SUpdateInternalTextProto), hc)

	case 68:
		guildModule.ProcessUpdateJoinCondition(proto.(*pb9.C2SUpdateJoinConditionProto), hc)

	case 71:
		guildModule.ProcessUpdateGuildName(proto.(*pb9.C2SUpdateGuildNameProto), hc)

	case 75:
		guildModule.ProcessUpdateLabels(proto.(*pb9.C2SUpdateGuildLabelProto), hc)

	case 80:
		guildModule.ProcessCancelChangeLeader(hc)

	case 83:
		guildModule.ProcessDonation(proto.(*pb9.C2SDonateProto), hc)

	case 90:
		guildModule.ProcessUpgradeLevel(hc)

	case 93:
		guildModule.ProcessReduceUpgradeLevelCd(proto.(*pb9.C2SReduceUpgradeLevelCdProto), hc)

	case 96:
		guildModule.ProcessImpeachLeader(hc)

	case 99:
		guildModule.ProcessImpeachLeaderVote(proto.(*pb9.C2SImpeachLeaderVoteProto), hc)

	case 102:
		guildModule.ProcessListGuildByIds(proto.(*pb9.C2SListGuildByIdsProto), hc)

	case 109:
		guildModule.ProcessInvateOtherRequest(proto.(*pb9.C2SGuildInvateOtherProto), hc)

	case 112:
		guildModule.ProcessCancelInvateOtherRequest(proto.(*pb9.C2SGuildCancelInvateOtherProto), hc)

	case 122:
		guildModule.ProcessUpdateClassTitle(proto.(*pb9.C2SUpdateClassTitleProto), hc)

	case 125:
		guildModule.ProcessUpdateFriendGuild(proto.(*pb9.C2SUpdateFriendGuildProto), hc)

	case 128:
		guildModule.ProcessUpdateEnemyGuild(proto.(*pb9.C2SUpdateEnemyGuildProto), hc)

	case 131:
		guildModule.ProcessUpdateGuildPrestige(proto.(*pb9.C2SUpdateGuildPrestigeProto), hc)

	case 134:
		guildModule.ProcessPlaceGuildStatue(proto.(*pb9.C2SPlaceGuildStatueProto), hc)

	case 138:
		guildModule.ProcessTakeBackGuildStatue(hc)

	case 143:
		guildModule.ProcessCollectFirstJoinGuildPrize(hc)

	case 147:
		guildModule.ProcessSeekHelp(proto.(*pb9.C2SSeekHelpProto), hc)

	case 151:
		guildModule.ProcessHelpGuildMember(proto.(*pb9.C2SHelpGuildMemberProto), hc)

	case 158:
		guildModule.ProcessHelpAllGuildMember(hc)

	case 163:
		guildModule.ProcessCollectGuildEventPrize(proto.(*pb9.C2SCollectGuildEventPrizeProto), hc)

	case 167:
		guildModule.ProcessCollectFullBigBox(hc)

	case 172:
		guildModule.ProcessUpgradeTechnology(proto.(*pb9.C2SUpgradeTechnologyProto), hc)

	case 175:
		guildModule.ProcessReduceTechnologyCd(proto.(*pb9.C2SReduceTechnologyCdProto), hc)

	case 178:
		guildModule.ProcessListGuildLogs(proto.(*pb9.C2SListGuildLogsProto), hc)

	case 181:
		guildModule.ProcessRequestRecommendGuild(hc)

	case 184:
		guildModule.ProcessHelpTech(hc)

	case 187:
		guildModule.ProcessRecommendInviteHeros(hc)

	case 190:
		guildModule.ProcessSearchNoGuildHeros(proto.(*pb9.C2SSearchNoGuildHerosProto), hc)

	case 193:
		guildModule.ProcessListInviteMeGuild(hc)

	case 196:
		guildModule.ProcessUpdateGuildMark(proto.(*pb9.C2SUpdateGuildMarkProto), hc)

	case 199:
		guildModule.ProcessViewMcWarRecord(hc)

	case 202:
		guildModule.ProcessViewYinliangRecord(hc)

	case 205:
		guildModule.ProcessSendYinliangToOtherGuild(proto.(*pb9.C2SSendYinliangToOtherGuildProto), hc)

	case 208:
		guildModule.ProcessSendYinliangToMember(proto.(*pb9.C2SSendYinliangToMemberProto), hc)

	case 211:
		guildModule.ProcessPaySalary(hc)

	case 214:
		guildModule.ProcessSetSalary(proto.(*pb9.C2SSetSalaryProto), hc)

	case 218:
		guildModule.ProcessViewSendYinliangToGuild(hc)

	case 228:
		guildModule.ProcessConvene(proto.(*pb9.C2SConveneProto), hc)

	case 231:
		guildModule.ProcessCollectDailyGuildRankPrize(hc)

	case 234:
		guildModule.ProcessViewDailyGuildRank(hc)

	case 240:
		guildModule.ProcessAddRecommendMcBuild(proto.(*pb9.C2SAddRecommendMcBuildProto), hc)

	case 243:
		guildModule.ProcessViewTaskProgress(proto.(*pb9.C2SViewTaskProgressProto), hc)

	case 247:
		guildModule.ProcessCollectTaskPrizeProgress(proto.(*pb9.C2SCollectTaskPrizeProto), hc)

	case 250:
		guildModule.ProcessGuildChangeCountry(proto.(*pb9.C2SGuildChangeCountryProto), hc)

	case 253:
		guildModule.ProcessCancelGuildChangeCountry(hc)

	case 256:
		guildModule.ProcessShowWorkshopNotExist(hc)

	default:
		return errors.New(fmt.Sprintf("模块guild消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_stress(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		stressModule.Ping(proto.(*pb10.C2SRobotPingProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块stress消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_depot(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 2:
		depotModule.ProcessUseGoods(proto.(*pb11.C2SUseGoodsProto), hc)

	case 6:
		depotModule.ProcessUseCdrGoods(proto.(*pb11.C2SUseCdrGoodsProto), hc)

	case 9:
		depotModule.ProcessGoodsCombine(proto.(*pb11.C2SGoodsCombineProto), hc)

	case 18:
		depotModule.ProcessGoodsPartCombine(proto.(*pb11.C2SGoodsPartsCombineProto), hc)

	case 23:
		depotModule.ProcessUnlockBaowu(proto.(*pb11.C2SUnlockBaowuProto), hc)

	case 26:
		depotModule.ProcessCollectBaowu(proto.(*pb11.C2SCollectBaowuProto), hc)

	case 30:
		depotModule.ProcessListBaowuLog(hc)

	case 35:
		depotModule.ProcessDecomposeBaowu(proto.(*pb11.C2SDecomposeBaowuProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块depot消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_equipment(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		equipmentModule.ProcessWearEquipment(proto.(*pb12.C2SWearEquipmentProto), hc)

	case 4:
		equipmentModule.ProcessUpgradeEquipment(proto.(*pb12.C2SUpgradeEquipmentProto), hc)

	case 7:
		equipmentModule.ProcessRefinedEquipment(proto.(*pb12.C2SRefinedEquipmentProto), hc)

	case 10:
		equipmentModule.ProcessSmeltEquipment(proto.(*pb12.C2SSmeltEquipmentProto), hc)

	case 13:
		equipmentModule.ProcessRebuildEquipment(proto.(*pb12.C2SRebuildEquipmentProto), hc)

	case 19:
		equipmentModule.ProcessUpgradeEquipmentAll(proto.(*pb12.C2SUpgradeEquipmentAllProto), hc)

	case 40:
		equipmentModule.ProcessViewChatEquip(proto.(*pb12.C2SViewChatEquipProto), hc)

	case 43:
		equipmentModule.ProcessOneKeyTakeOff(proto.(*pb12.C2SOneKeyTakeOffProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块equipment消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_chat(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		chatModule.ProcessWorldChat(proto.(*pb13.C2SWorldChatProto), hc)

	case 4:
		chatModule.ProcessGuildChat(proto.(*pb13.C2SGuildChatProto), hc)

	case 8:
		chatModule.ProcessSelfChatWindow(hc)

	case 10:
		chatModule.ProcessRemoveChatWindow(proto.(*pb13.C2SRemoveChatWindowProto), hc)

	case 12:
		chatModule.ProcessListHistoryChat(proto.(*pb13.C2SListHistoryChatProto), hc)

	case 14:
		chatModule.ProcessSendChat(proto.(*pb13.C2SSendChatProto), hc)

	case 18:
		chatModule.ProcessReadChatMsg(proto.(*pb13.C2SReadChatMsgProto), hc)

	case 21:
		chatModule.ProcessCreateSelfChatWindow(proto.(*pb13.C2SCreateSelfChatWindowProto), hc)

	case 25:
		chatModule.ProcessGetHeroChatInfo(proto.(*pb13.C2SGetHeroChatInfoProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块chat消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_tower(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		towerModule.ProcessChallenge(proto.(*pb14.C2SChallengeProto), hc)

	case 5:
		towerModule.ProcessAutoChallenge(hc)

	case 8:
		towerModule.ProcessCollectBox(proto.(*pb14.C2SCollectBoxProto), hc)

	case 11:
		towerModule.ProcessListPassReplay(proto.(*pb14.C2SListPassReplayProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块tower消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_task(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 2:
		taskModule.ProcessCollectTaskPrize(proto.(*pb15.C2SCollectTaskPrizeProto), hc)

	case 6:
		taskModule.ProcessCollectTaskBoxPrize(proto.(*pb15.C2SCollectTaskBoxPrizeProto), hc)

	case 9:
		taskModule.ProcessCollectBaYeStagePrize(hc)

	case 12:
		taskModule.ProcessCollectActiveDegreePrize(proto.(*pb15.C2SCollectActiveDegreePrizeProto), hc)

	case 16:
		taskModule.ProcessCollectAchieveStarPrize(proto.(*pb15.C2SCollectAchieveStarPrizeProto), hc)

	case 20:
		taskModule.ProcessChangeSelectShowAchieve(proto.(*pb15.C2SChangeSelectShowAchieveProto), hc)

	case 23:
		taskModule.ProcessCollectBwzlPrize(proto.(*pb15.C2SCollectBwzlPrizeProto), hc)

	case 26:
		taskModule.ProcessViewOtherAchieveTaskList(proto.(*pb15.C2SViewOtherAchieveTaskListProto), hc)

	case 29:
		taskModule.ProcessGetUpgradeTitleFightAmount(proto.(*pb15.C2SGetTroopTitleFightAmountProto), hc)

	case 31:
		taskModule.ProcessUpgradeTitle(hc)

	case 35:
		taskModule.ProcessCompleteBoolTask(proto.(*pb15.C2SCompleteBoolTaskProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块task消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_fishing(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		fishingModule.ProcessFishing(proto.(*pb16.C2SFishingProto), hc)

	case 8:
		fishingModule.ProcessFishPointExchange(hc)

	case 11:
		fishingModule.ProcessSetFishingCaptain(proto.(*pb16.C2SSetFishingCaptainProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块fishing消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_gem(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 3:
		gemModule.ProcessUseGem(proto.(*pb19.C2SUseGemProto), hc)

	case 6:
		gemModule.ProcessCombineGem(proto.(*pb19.C2SCombineGemProto), hc)

	case 9:
		gemModule.ProcessOneKeyUseGem(proto.(*pb19.C2SOneKeyUseGemProto), hc)

	case 11:
		gemModule.ProcessOneKeyCombineGem(proto.(*pb19.C2SOneKeyCombineGemProto), hc)

	case 15:
		gemModule.ProcessRequestOneKeyCombineGemCost(proto.(*pb19.C2SRequestOneKeyCombineCostProto), hc)

	case 18:
		gemModule.ProcessOneKeyCombineDepotGem(proto.(*pb19.C2SOneKeyCombineDepotGemProto), hc)

	case 21:
		gemModule.ProcessInlayGem(proto.(*pb19.C2SInlayGemProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块gem消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_shop(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 2:
		shopModule.ProcessBuyGoods(proto.(*pb20.C2SBuyGoodsProto), hc)

	case 9:
		shopModule.ProcessBuyBlackMarketGoods(proto.(*pb20.C2SBuyBlackMarketGoodsProto), hc)

	case 12:
		shopModule.ProcessRefreshBlackMarketGoods(hc)

	default:
		return errors.New(fmt.Sprintf("模块shop消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_client_config(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		clientConfigModule.ProcessConfig(proto.(*pb21.C2SConfigProto), hc)

	case 4:
		clientConfigModule.ProcessSetClientData(proto.(*pb21.C2SSetClientDataProto), hc)

	case 5:
		clientConfigModule.ProcessSetClientKey(proto.(*pb21.C2SSetClientKeyProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块client_config消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_secret_tower(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 2:
		secretTowerModule.ProcessRequestTeamCount(hc)

	case 5:
		secretTowerModule.ProcessRequestTeamList(proto.(*pb22.C2SRequestTeamListProto), hc)

	case 8:
		secretTowerModule.ProcessCreateTeam(proto.(*pb22.C2SCreateTeamProto), hc)

	case 11:
		secretTowerModule.ProcessJoinTeam(proto.(*pb22.C2SJoinTeamProto), hc)

	case 15:
		secretTowerModule.ProcessLeaveTeam(hc)

	case 19:
		secretTowerModule.ProcessKickMember(proto.(*pb22.C2SKickMemberProto), hc)

	case 24:
		secretTowerModule.ProcessMoveMember(proto.(*pb22.C2SMoveMemberProto), hc)

	case 27:
		secretTowerModule.ProcessChangeMode(proto.(*pb22.C2SChangeModeProto), hc)

	case 33:
		secretTowerModule.ProcessInvite(proto.(*pb22.C2SInviteProto), hc)

	case 37:
		secretTowerModule.ProcessRequestInviteList(hc)

	case 39:
		secretTowerModule.ProcessRequestTeamDetail(hc)

	case 42:
		secretTowerModule.ProcessStartChallenge(hc)

	case 58:
		secretTowerModule.ProcessQuickQueryTeamBasic(proto.(*pb22.C2SQuickQueryTeamBasicProto), hc)

	case 61:
		secretTowerModule.ProcessChangeGuildMode(hc)

	case 67:
		secretTowerModule.ProcessUpdateMemberPos(proto.(*pb22.C2SUpdateMemberPosProto), hc)

	case 71:
		secretTowerModule.ProcessInviteAll(proto.(*pb22.C2SInviteAllProto), hc)

	case 74:
		secretTowerModule.ProcessListRecord(hc)

	case 79:
		secretTowerModule.ProcessTeamTalk(proto.(*pb22.C2STeamTalkProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块secret_tower消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_rank(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		rankModule.ProcessRequestRank(proto.(*pb23.C2SRequestRankProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块rank消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_bai_zhan(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		baiZhanModule.ProcessQueryBaiZhanInfo(hc)

	case 4:
		baiZhanModule.ProcessChallenge(hc)

	case 7:
		baiZhanModule.ProcessCollectSalary(hc)

	case 10:
		baiZhanModule.ProcessCollectJunXianPrize(proto.(*pb24.C2SCollectJunXianPrizeProto), hc)

	case 23:
		baiZhanModule.ProcessRequestRank(proto.(*pb24.C2SRequestRankProto), hc)

	case 26:
		baiZhanModule.ProcessRequestSelfRank(hc)

	case 29:
		baiZhanModule.ProcessSelfRecord(proto.(*pb24.C2SSelfRecordProto), hc)

	case 34:
		baiZhanModule.ProcesClearLastJunXian(hc)

	default:
		return errors.New(fmt.Sprintf("模块bai_zhan消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_dungeon(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		dungeonModule.ProcessChallenge(proto.(*pb26.C2SChallengeProto), hc)

	case 4:
		dungeonModule.ProcessCollectChapterPrize(proto.(*pb26.C2SCollectChapterPrizeProto), hc)

	case 7:
		dungeonModule.ProcessAutoChallenge(proto.(*pb26.C2SAutoChallengeProto), hc)

	case 13:
		dungeonModule.ProcessCollectPassDungeonPrize(proto.(*pb26.C2SCollectPassDungeonPrizeProto), hc)

	case 17:
		dungeonModule.ProcessCollectChapterStarPrize(proto.(*pb26.C2SCollectChapterStarPrizeProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块dungeon消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_country(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 16:
		countryModule.ProcessRequestCountryPrestige(proto.(*pb27.C2SRequestCountryPrestigeProto), hc)

	case 19:
		countryModule.ProcessRequestCountries(proto.(*pb27.C2SRequestCountriesProto), hc)

	case 22:
		countryModule.ProcessHeroChangeCountry(proto.(*pb27.C2SHeroChangeCountryProto), hc)

	case 31:
		countryModule.ProcessCountryDetail(proto.(*pb27.C2SCountryDetailProto), hc)

	case 40:
		countryModule.ProcessOfficialAppoint(proto.(*pb27.C2SOfficialAppointProto), hc)

	case 43:
		countryModule.ProcessOfficialDepose(proto.(*pb27.C2SOfficialDeposeProto), hc)

	case 46:
		countryModule.ProcessCollectOfficialSalary(hc)

	case 54:
		countryModule.ProcessOfficialLeave(hc)

	case 61:
		countryModule.ProcessChangeNameVote(proto.(*pb27.C2SChangeNameVoteProto), hc)

	case 66:
		countryModule.ProcessSearchToAppointHeroList(proto.(*pb27.C2SSearchToAppointHeroListProto), hc)

	case 69:
		countryModule.ProcessDefaultToAppointHeroList(hc)

	case 72:
		countryModule.ProcessChangeNameStart(proto.(*pb27.C2SChangeNameStartProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块country消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_tag(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		tagModule.ProcessAddOrUpdateTag(proto.(*pb29.C2SAddOrUpdateTagProto), hc)

	case 5:
		tagModule.ProcessDeleteTag(proto.(*pb29.C2SDeleteTagProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块tag消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_garden(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		gardenModule.ProcessListTreasuryTreeHero(hc)

	case 3:
		gardenModule.ProcessListTreasuryTreeTimes(hc)

	case 5:
		gardenModule.ProcessWaterTreasuryTree(proto.(*pb31.C2SWaterTreasuryTreeProto), hc)

	case 10:
		gardenModule.ProcessCollectTreasureTreePrize(hc)

	case 13:
		gardenModule.ProcessListHelpMe(proto.(*pb31.C2SListHelpMeProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块garden消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_zhengwu(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		zhengWuModule.ProcessStart(proto.(*pb32.C2SStartProto), hc)

	case 4:
		zhengWuModule.ProcessCollect(hc)

	case 7:
		zhengWuModule.ProcessYuanBaoComplete(hc)

	case 10:
		zhengWuModule.ProcessYuanBaoRefresh(hc)

	case 14:
		zhengWuModule.ProcessVipCollect(proto.(*pb32.C2SVipCollectProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块zhengwu消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_zhanjiang(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		zhanJiangModule.ProcessOpen(proto.(*pb33.C2SOpenProto), hc)

	case 4:
		zhanJiangModule.ProcessGiveUp(hc)

	case 7:
		zhanJiangModule.ProcessUpdateCaptain(proto.(*pb33.C2SUpdateCaptainProto), hc)

	case 10:
		zhanJiangModule.ProcessChallenge(proto.(*pb33.C2SChallengeProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块zhanjiang消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_question(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		questionModule.ProcessStart(proto.(*pb34.C2SStartProto), hc)

	case 4:
		questionModule.ProcessAnswer(proto.(*pb34.C2SAnswerProto), hc)

	case 6:
		questionModule.ProcessNext(proto.(*pb34.C2SNextProto), hc)

	case 9:
		questionModule.ProcessGetPrize(proto.(*pb34.C2SGetPrizeProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块question消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_relation(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		relationModule.ProcessAddRelation(proto.(*pb35.C2SAddRelationProto), hc)

	case 4:
		relationModule.ProcessRemoveRelation(proto.(*pb35.C2SRemoveRelationProto), hc)

	case 7:
		relationModule.ProcessListRelation(proto.(*pb35.C2SListRelationProto), hc)

	case 10:
		relationModule.ProcessRemoveEnemy(proto.(*pb35.C2SRemoveEnemyProto), hc)

	case 16:
		relationModule.ProcessRecommendHeroList(proto.(*pb35.C2SRecommendHeroListProto), hc)

	case 22:
		relationModule.ProcessSearchHeros(proto.(*pb35.C2SSearchHerosProto), hc)

	case 25:
		relationModule.ProcessSearchHeroById(proto.(*pb35.C2SSearchHeroByIdProto), hc)

	case 28:
		relationModule.ProcessNewListRelation(proto.(*pb35.C2SNewListRelationProto), hc)

	case 33:
		relationModule.ProcessSetImportantFriend(proto.(*pb35.C2SSetImportantFriendProto), hc)

	case 36:
		relationModule.ProcessCancelImportantFriend(proto.(*pb35.C2SCancelImportantFriendProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块relation消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_xiongnu(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		xiongNuModule.ProcessSetDefender(proto.(*pb36.C2SSetDefenderProto), hc)

	case 5:
		xiongNuModule.ProcessStart(proto.(*pb36.C2SStartProto), hc)

	case 10:
		xiongNuModule.ProcessTroopInfo(hc)

	case 14:
		xiongNuModule.ProcessGetXiongNuNpcBaseInfo(proto.(*pb36.C2SGetXiongNuNpcBaseInfoProto), hc)

	case 17:
		xiongNuModule.ProcessGetDefenserFightAmount(proto.(*pb36.C2SGetDefenserFightAmountProto), hc)

	case 19:
		xiongNuModule.ProcessGetXiongNuFightInfo(hc)

	default:
		return errors.New(fmt.Sprintf("模块xiongnu消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_survey(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 2:
		surveyModule.ProcessComplete(proto.(*pb37.C2SCompleteProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块survey消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_farm(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 2:
		farmModule.ProcessPlant(proto.(*pb38.C2SPlantProto), hc)

	case 5:
		farmModule.ProcessHarvest(proto.(*pb38.C2SHarvestProto), hc)

	case 8:
		farmModule.ProcessChange(proto.(*pb38.C2SChangeProto), hc)

	case 12:
		farmModule.ProcessOneKeyPlant(proto.(*pb38.C2SOneKeyPlantProto), hc)

	case 18:
		farmModule.ProcessSteal(proto.(*pb38.C2SStealProto), hc)

	case 28:
		farmModule.ProcessOneKeyHarvest(proto.(*pb38.C2SOneKeyHarvestProto), hc)

	case 31:
		farmModule.ProcessOneKeySteal(proto.(*pb38.C2SOneKeyStealProto), hc)

	case 39:
		farmModule.ProcessStealLogList(proto.(*pb38.C2SStealLogListProto), hc)

	case 43:
		farmModule.ProcessViewFarm(proto.(*pb38.C2SViewFarmProto), hc)

	case 48:
		farmModule.ProcessCanStealList(hc)

	case 52:
		farmModule.ProcessOneKeyReset(hc)

	default:
		return errors.New(fmt.Sprintf("模块farm消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_dianquan(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		dianquanModule.ProcessExchange(proto.(*pb39.C2SExchangeProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块dianquan消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_xuanyuan(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		xuanyuanModule.ProcessSelfInfo(hc)

	case 5:
		xuanyuanModule.ProcessQueryTargetTroop(proto.(*pb40.C2SQueryTargetTroopProto), hc)

	case 11:
		xuanyuanModule.ProcessListTarget(proto.(*pb40.C2SListTargetProto), hc)

	case 15:
		xuanyuanModule.ProcessChallenge(proto.(*pb40.C2SChallengeProto), hc)

	case 20:
		xuanyuanModule.ProcessListRecord(proto.(*pb40.C2SListRecordProto), hc)

	case 22:
		xuanyuanModule.ProcessCollectRankPrize(hc)

	default:
		return errors.New(fmt.Sprintf("模块xuanyuan消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_hebi(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		hebiModule.ProcessRoomList(proto.(*pb41.C2SRoomListProto), hc)

	case 3:
		hebiModule.ProcessChangeCaptain(proto.(*pb41.C2SChangeCaptainProto), hc)

	case 9:
		hebiModule.ProcessJoinRoom(proto.(*pb41.C2SJoinRoomProto), hc)

	case 12:
		hebiModule.ProcessRobPos(proto.(*pb41.C2SRobPosProto), hc)

	case 18:
		hebiModule.ProcessLeave(proto.(*pb41.C2SLeaveRoomProto), hc)

	case 21:
		hebiModule.ProcessRob(proto.(*pb41.C2SRobProto), hc)

	case 28:
		hebiModule.ProcessCheckInRoom(proto.(*pb41.C2SCheckInRoomProto), hc)

	case 31:
		hebiModule.ProcessCopySelf(proto.(*pb41.C2SCopySelfProto), hc)

	case 35:
		hebiModule.ProcessHebiHeroRecordList(hc)

	case 37:
		hebiModule.ProcessViewHebiShowPrize(proto.(*pb41.C2SViewShowPrizeProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块hebi消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_mingc(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 4:
		mingcModule.ProcessMingcList(proto.(*pb42.C2SMingcListProto), hc)

	case 7:
		mingcModule.ProcessViewMingc(proto.(*pb42.C2SViewMingcProto), hc)

	case 10:
		mingcModule.ProcessMcBuild(proto.(*pb42.C2SMcBuildProto), hc)

	case 13:
		mingcModule.ProcessMcBuildLog(proto.(*pb42.C2SMcBuildLogProto), hc)

	case 20:
		mingcModule.ProcessMingcHostGuild(proto.(*pb42.C2SMingcHostGuildProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块mingc消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_promotion(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 4:
		promotionModule.ProcessCollectLogin7DayPrize(proto.(*pb43.C2SCollectLoginDayPrizeProto), hc)

	case 7:
		promotionModule.ProcessBuyHeroLevelFund(hc)

	case 10:
		promotionModule.ProcessCollectHeroLevelFund(proto.(*pb43.C2SCollectLevelFundProto), hc)

	case 13:
		promotionModule.ProcessCollectDailySp(proto.(*pb43.C2SCollectDailySpProto), hc)

	case 16:
		promotionModule.ProcessCollectFreeGift(proto.(*pb43.C2SCollectFreeGiftProto), hc)

	case 21:
		promotionModule.ProcessBuyTimeLimitGift(proto.(*pb43.C2SBuyTimeLimitGiftProto), hc)

	case 26:
		promotionModule.ProcessBuyEventLimitGift(proto.(*pb43.C2SBuyEventLimitGiftProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块promotion消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_mingc_war(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 16:
		mingcWarModule.ProcessApplyAtk(proto.(*pb44.C2SApplyAtkProto), hc)

	case 21:
		mingcWarModule.ProcessApplyAst(proto.(*pb44.C2SApplyAstProto), hc)

	case 25:
		mingcWarModule.ProcessReplyApplyAst(proto.(*pb44.C2SReplyApplyAstProto), hc)

	case 29:
		mingcWarModule.ProcessViewMcWar(proto.(*pb44.C2SViewMcWarProto), hc)

	case 31:
		mingcWarModule.ProcessViewMcWarSelfGuild(hc)

	case 35:
		mingcWarModule.ProcessJoinFight(proto.(*pb44.C2SJoinFightProto), hc)

	case 38:
		mingcWarModule.ProcessQuitFight(hc)

	case 46:
		mingcWarModule.ProcessViewMcWarScene(proto.(*pb44.C2SViewMcWarSceneProto), hc)

	case 49:
		mingcWarModule.ProcessSceneMove(proto.(*pb44.C2SSceneMoveProto), hc)

	case 72:
		mingcWarModule.ProcessSceneTroopRelive(hc)

	case 75:
		mingcWarModule.ProcessViewMingcWarMc(proto.(*pb44.C2SViewMingcWarMcProto), hc)

	case 80:
		mingcWarModule.ProcessCancelApplyAst(proto.(*pb44.C2SCancelApplyAstProto), hc)

	case 85:
		mingcWarModule.ProcessSceneBack(hc)

	case 88:
		mingcWarModule.ProcessSceneSpeedUp(proto.(*pb44.C2SSceneSpeedUpProto), hc)

	case 91:
		mingcWarModule.ProcessViewMcWarRecord(proto.(*pb44.C2SViewMcWarRecordProto), hc)

	case 94:
		mingcWarModule.ProcessViewMcWarTroopRecord(proto.(*pb44.C2SViewMcWarTroopRecordProto), hc)

	case 99:
		mingcWarModule.ProcessViewSceneTroopRecord(hc)

	case 107:
		mingcWarModule.ProcessApplyRefreshRank(proto.(*pb44.C2SApplyRefreshRankProto), hc)

	case 111:
		mingcWarModule.ProcessViewMyGuildMemberRank(proto.(*pb44.C2SViewMyGuildMemberRankProto), hc)

	case 115:
		mingcWarModule.ProcessSceneChangeMode(proto.(*pb44.C2SSceneChangeModeProto), hc)

	case 119:
		mingcWarModule.ProcessSceneTouShiBuildingTurnTo(proto.(*pb44.C2SSceneTouShiBuildingTurnToProto), hc)

	case 123:
		mingcWarModule.ProcessSceneTouShiBuildingFire(proto.(*pb44.C2SSceneTouShiBuildingFireProto), hc)

	case 128:
		mingcWarModule.ProcessSceneDrum(hc)

	case 136:
		mingcWarModule.ProcessQuitWatch(proto.(*pb44.C2SQuitWatchProto), hc)

	case 139:
		mingcWarModule.ProcessWatch(proto.(*pb44.C2SWatchProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块mingc_war消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_random_event(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		randomEventModule.ProcessChooseOption(proto.(*pb45.C2SChooseOptionProto), hc)

	case 4:
		randomEventModule.ProcessOpenEvent(proto.(*pb45.C2SOpenEventProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块random_event消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_strategy(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		strategyModule.ProcessUseStratagem(proto.(*pb46.C2SUseStratagemProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块strategy消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_vip(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 4:
		vipModule.ProcessVipCollectDailyPrize(proto.(*pb48.C2SVipCollectDailyPrizeProto), hc)

	case 7:
		vipModule.ProcessVipCollectLevelPrize(proto.(*pb48.C2SVipCollectLevelPrizeProto), hc)

	case 12:
		vipModule.ProcessVipBuyDungeonTimes(proto.(*pb48.C2SVipBuyDungeonTimesProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块vip消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_red_packet(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		redPacketModule.ProcessBuy(proto.(*pb49.C2SBuyProto), hc)

	case 4:
		redPacketModule.ProcessCreate(proto.(*pb49.C2SCreateProto), hc)

	case 7:
		redPacketModule.ProcessGrab(proto.(*pb49.C2SGrabProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块red_packet消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_teach(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		teachModule.ProcessFight(proto.(*pb50.C2SFightProto), hc)

	case 4:
		teachModule.ProcessCollectPrize(proto.(*pb50.C2SCollectPrizeProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块teach消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

func handle_activity(sequenceID int, hc iface.HeroController, proto interface{}) error {
	switch sequenceID {
	case 1:
		activityModule.ProcessCollectCollection(proto.(*pb51.C2SCollectCollectionProto), hc)

	default:
		return errors.New(fmt.Sprintf("模块activity消息处理时没有handler: %d", sequenceID))
	}
	return nil
}

/*
----------- 自动识别分割线 -----------
package service

import(
	"github.com/lightpaw/male7/gen/iface"
)

// 包含消息的信息
type MsgData struct {
	ModuleID   int
	SequenceID int

	// 处理这条消息需要的真实proto. 例如 *MountUpgradeProto
	Proto interface{}
}

func Unmarshal(moduleID, sequenceID int, data []byte) (*MsgData, error) {
	return nil, nil
}

func Handle(msg *MsgData, hc iface.HeroController) error {
	return nil
}

----------- 自动识别分割线 -----------
*/
