package service

import (
	"github.com/golang/protobuf/proto"
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
	"strconv"
)

func PrintObject(moduleID, sequenceID int, data []byte) (proto.Message, error) {
	switch moduleID {
	case 1: // login
		return print_login(sequenceID, data)

	case 2: // domestic
		return print_domestic(sequenceID, data)

	case 3: // gm
		return print_gm(sequenceID, data)

	case 4: // military
		return print_military(sequenceID, data)

	case 5: // misc
		return print_misc(sequenceID, data)

	case 7: // region
		return print_region(sequenceID, data)

	case 8: // mail
		return print_mail(sequenceID, data)

	case 9: // guild
		return print_guild(sequenceID, data)

	case 10: // stress
		return print_stress(sequenceID, data)

	case 11: // depot
		return print_depot(sequenceID, data)

	case 12: // equipment
		return print_equipment(sequenceID, data)

	case 13: // chat
		return print_chat(sequenceID, data)

	case 14: // tower
		return print_tower(sequenceID, data)

	case 15: // task
		return print_task(sequenceID, data)

	case 16: // fishing
		return print_fishing(sequenceID, data)

	case 19: // gem
		return print_gem(sequenceID, data)

	case 20: // shop
		return print_shop(sequenceID, data)

	case 21: // client_config
		return print_client_config(sequenceID, data)

	case 22: // secret_tower
		return print_secret_tower(sequenceID, data)

	case 23: // rank
		return print_rank(sequenceID, data)

	case 24: // bai_zhan
		return print_bai_zhan(sequenceID, data)

	case 26: // dungeon
		return print_dungeon(sequenceID, data)

	case 27: // country
		return print_country(sequenceID, data)

	case 29: // tag
		return print_tag(sequenceID, data)

	case 31: // garden
		return print_garden(sequenceID, data)

	case 32: // zhengwu
		return print_zhengwu(sequenceID, data)

	case 33: // zhanjiang
		return print_zhanjiang(sequenceID, data)

	case 34: // question
		return print_question(sequenceID, data)

	case 35: // relation
		return print_relation(sequenceID, data)

	case 36: // xiongnu
		return print_xiongnu(sequenceID, data)

	case 37: // survey
		return print_survey(sequenceID, data)

	case 38: // farm
		return print_farm(sequenceID, data)

	case 39: // dianquan
		return print_dianquan(sequenceID, data)

	case 40: // xuanyuan
		return print_xuanyuan(sequenceID, data)

	case 41: // hebi
		return print_hebi(sequenceID, data)

	case 42: // mingc
		return print_mingc(sequenceID, data)

	case 43: // promotion
		return print_promotion(sequenceID, data)

	case 44: // mingc_war
		return print_mingc_war(sequenceID, data)

	case 45: // random_event
		return print_random_event(sequenceID, data)

	case 46: // strategy
		return print_strategy(sequenceID, data)

	case 48: // vip
		return print_vip(sequenceID, data)

	case 49: // red_packet
		return print_red_packet(sequenceID, data)

	case 50: // teach
		return print_teach(sequenceID, data)

	case 51: // activity
		return print_activity(sequenceID, data)

	default:
		return nil, errors.Errorf("打印未知消息: %d.%d", moduleID, sequenceID)
	}
}

func print_login(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_internal_login
		p := &pb1.S2CInternalLoginProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "login.s2c_internal_login PrintMsgProto &S2CInternalLoginProto fail")
		}

		p.HeroProto = nil

		return p, nil

	case 5: //s2c_fail_internal_login
		return toErrCodeMessage(1, 5, data), nil

	case 8: //s2c_login
		p := &pb1.S2CLoginProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "login.s2c_login PrintMsgProto &S2CLoginProto fail")
		}

		return p, nil

	case 9: //s2c_fail_login
		return toErrCodeMessage(1, 9, data), nil

	case 17: //s2c_tutorial_progress
		p := &pb1.S2CTutorialProgressProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "login.s2c_tutorial_progress PrintMsgProto &S2CTutorialProgressProto fail")
		}

		return p, nil

	case 4: //s2c_create_hero
		return toStringMessage(1, 4), nil

	case 6: //s2c_fail_create_hero
		return toErrCodeMessage(1, 6, data), nil

	case 11: //s2c_loaded
		p := &pb1.S2CLoadedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "login.s2c_loaded PrintMsgProto &S2CLoadedProto fail")
		}

		p.HeroProto = nil

		return p, nil

	case 12: //s2c_fail_loaded
		return toErrCodeMessage(1, 12, data), nil

	case 14: //s2c_robot_login
		return toStringMessage(1, 14), nil

	case 19: //s2c_ban_login
		p := &pb1.S2CBanLoginProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "login.s2c_ban_login PrintMsgProto &S2CBanLoginProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_domestic(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 74: //s2c_update_resource_building
		p := &pb2.S2CUpdateResourceBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_resource_building PrintMsgProto &S2CUpdateResourceBuildingProto fail")
		}

		return p, nil

	case 75: //s2c_update_multi_resource_building
		p := &pb2.S2CUpdateMultiResourceBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_multi_resource_building PrintMsgProto &S2CUpdateMultiResourceBuildingProto fail")
		}

		return p, nil

	case 2: //s2c_create_building
		p := &pb2.S2CCreateBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_create_building PrintMsgProto &S2CCreateBuildingProto fail")
		}

		return p, nil

	case 3: //s2c_fail_create_building
		return toErrCodeMessage(2, 3, data), nil

	case 5: //s2c_upgrade_building
		p := &pb2.S2CUpgradeBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_upgrade_building PrintMsgProto &S2CUpgradeBuildingProto fail")
		}

		return p, nil

	case 6: //s2c_fail_upgrade_building
		return toErrCodeMessage(2, 6, data), nil

	case 8: //s2c_rebuild_resource_building
		p := &pb2.S2CRebuildResourceBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_rebuild_resource_building PrintMsgProto &S2CRebuildResourceBuildingProto fail")
		}

		return p, nil

	case 9: //s2c_fail_rebuild_resource_building
		return toErrCodeMessage(2, 9, data), nil

	case 109: //s2c_unlock_outer_city
		p := &pb2.S2CUnlockOuterCityProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_unlock_outer_city PrintMsgProto &S2CUnlockOuterCityProto fail")
		}

		p.OuterCity = nil

		return p, nil

	case 110: //s2c_fail_unlock_outer_city
		return toErrCodeMessage(2, 110, data), nil

	case 143: //s2c_update_outer_city_type
		p := &pb2.S2CUpdateOuterCityTypeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_outer_city_type PrintMsgProto &S2CUpdateOuterCityTypeProto fail")
		}

		return p, nil

	case 144: //s2c_fail_update_outer_city_type
		return toErrCodeMessage(2, 144, data), nil

	case 112: //s2c_upgrade_outer_city_building
		p := &pb2.S2CUpgradeOuterCityBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_upgrade_outer_city_building PrintMsgProto &S2CUpgradeOuterCityBuildingProto fail")
		}

		return p, nil

	case 113: //s2c_fail_upgrade_outer_city_building
		return toErrCodeMessage(2, 113, data), nil

	case 13: //s2c_resource_update
		p := &pb2.S2CResourceUpdateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_resource_update PrintMsgProto &S2CResourceUpdateProto fail")
		}

		return p, nil

	case 28: //s2c_resource_update_single
		p := &pb2.S2CResourceUpdateSingleProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_resource_update_single PrintMsgProto &S2CResourceUpdateSingleProto fail")
		}

		return p, nil

	case 14: //s2c_resource_capcity_update
		p := &pb2.S2CResourceCapcityUpdateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_resource_capcity_update PrintMsgProto &S2CResourceCapcityUpdateProto fail")
		}

		return p, nil

	case 16: //s2c_collect_resource
		p := &pb2.S2CCollectResourceProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_collect_resource PrintMsgProto &S2CCollectResourceProto fail")
		}

		return p, nil

	case 17: //s2c_fail_collect_resource
		return toErrCodeMessage(2, 17, data), nil

	case 77: //s2c_collect_resource_v2
		p := &pb2.S2CCollectResourceV2Proto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_collect_resource_v2 PrintMsgProto &S2CCollectResourceV2Proto fail")
		}

		return p, nil

	case 78: //s2c_fail_collect_resource_v2
		return toErrCodeMessage(2, 78, data), nil

	case 79: //s2c_collect_resource_times_changed
		p := &pb2.S2CCollectResourceTimesChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_collect_resource_times_changed PrintMsgProto &S2CCollectResourceTimesChangedProto fail")
		}

		return p, nil

	case 80: //s2c_resource_point_change_v2
		p := &pb2.S2CResourcePointChangeV2Proto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_resource_point_change_v2 PrintMsgProto &S2CResourcePointChangeV2Proto fail")
		}

		p.Data = nil

		return p, nil

	case 82: //s2c_request_resource_conflict
		p := &pb2.S2CRequestResourceConflictProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_request_resource_conflict PrintMsgProto &S2CRequestResourceConflictProto fail")
		}

		return p, nil

	case 83: //s2c_fail_request_resource_conflict
		return toErrCodeMessage(2, 83, data), nil

	case 84: //s2c_resource_conflict_changed
		return toStringMessage(2, 84), nil

	case 19: //s2c_learn_technology
		p := &pb2.S2CLearnTechnologyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_learn_technology PrintMsgProto &S2CLearnTechnologyProto fail")
		}

		return p, nil

	case 20: //s2c_fail_learn_technology
		return toErrCodeMessage(2, 20, data), nil

	case 88: //s2c_unlock_stable_building
		p := &pb2.S2CUnlockStableBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_unlock_stable_building PrintMsgProto &S2CUnlockStableBuildingProto fail")
		}

		return p, nil

	case 89: //s2c_fail_unlock_stable_building
		return toErrCodeMessage(2, 89, data), nil

	case 25: //s2c_upgrade_stable_building
		p := &pb2.S2CUpgradeStableBuildingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_upgrade_stable_building PrintMsgProto &S2CUpgradeStableBuildingProto fail")
		}

		return p, nil

	case 26: //s2c_fail_upgrade_stable_building
		return toErrCodeMessage(2, 26, data), nil

	case 22: //s2c_hero_update_exp
		p := &pb2.S2CHeroUpdateExpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_hero_update_exp PrintMsgProto &S2CHeroUpdateExpProto fail")
		}

		return p, nil

	case 21: //s2c_hero_upgrade_level
		p := &pb2.S2CHeroUpgradeLevelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_hero_upgrade_level PrintMsgProto &S2CHeroUpgradeLevelProto fail")
		}

		return p, nil

	case 27: //s2c_hero_update_prosperity
		p := &pb2.S2CHeroUpdateProsperityProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_hero_update_prosperity PrintMsgProto &S2CHeroUpdateProsperityProto fail")
		}

		return p, nil

	case 91: //s2c_is_hero_name_exist
		p := &pb2.S2CIsHeroNameExistProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_is_hero_name_exist PrintMsgProto &S2CIsHeroNameExistProto fail")
		}

		return p, nil

	case 92: //s2c_fail_is_hero_name_exist
		return toErrCodeMessage(2, 92, data), nil

	case 31: //s2c_change_hero_name
		p := &pb2.S2CChangeHeroNameProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_change_hero_name PrintMsgProto &S2CChangeHeroNameProto fail")
		}

		return p, nil

	case 32: //s2c_fail_change_hero_name
		return toErrCodeMessage(2, 32, data), nil

	case 93: //s2c_give_first_change_hero_name_prize
		return toStringMessage(2, 93), nil

	case 49: //s2c_hero_name_changed_broadcast
		p := &pb2.S2CHeroNameChangedBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_hero_name_changed_broadcast PrintMsgProto &S2CHeroNameChangedBroadcastProto fail")
		}

		p.Id = nil

		return p, nil

	case 34: //s2c_list_old_name
		p := &pb2.S2CListOldNameProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_list_old_name PrintMsgProto &S2CListOldNameProto fail")
		}

		return p, nil

	case 37: //s2c_fail_list_old_name
		return toErrCodeMessage(2, 37, data), nil

	case 36: //s2c_view_other_hero
		p := &pb2.S2CViewOtherHeroProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_view_other_hero PrintMsgProto &S2CViewOtherHeroProto fail")
		}

		p.Hero = nil

		return p, nil

	case 38: //s2c_fail_view_other_hero
		return toErrCodeMessage(2, 38, data), nil

	case 126: //s2c_view_fight_info
		p := &pb2.S2CViewFightInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_view_fight_info PrintMsgProto &S2CViewFightInfoProto fail")
		}

		p.Id = nil

		return p, nil

	case 39: //s2c_update_building_worker_coef
		p := &pb2.S2CUpdateBuildingWorkerCoefProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_building_worker_coef PrintMsgProto &S2CUpdateBuildingWorkerCoefProto fail")
		}

		return p, nil

	case 40: //s2c_update_tech_worker_coef
		p := &pb2.S2CUpdateTechWorkerCoefProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_tech_worker_coef PrintMsgProto &S2CUpdateTechWorkerCoefProto fail")
		}

		return p, nil

	case 55: //s2c_update_building_worker_fatigue_duration
		p := &pb2.S2CUpdateBuildingWorkerFatigueDurationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_building_worker_fatigue_duration PrintMsgProto &S2CUpdateBuildingWorkerFatigueDurationProto fail")
		}

		return p, nil

	case 56: //s2c_update_tech_worker_fatigue_duration
		p := &pb2.S2CUpdateTechWorkerFatigueDurationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_tech_worker_fatigue_duration PrintMsgProto &S2CUpdateTechWorkerFatigueDurationProto fail")
		}

		return p, nil

	case 42: //s2c_miao_building_worker_cd
		p := &pb2.S2CMiaoBuildingWorkerCdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_miao_building_worker_cd PrintMsgProto &S2CMiaoBuildingWorkerCdProto fail")
		}

		return p, nil

	case 43: //s2c_fail_miao_building_worker_cd
		return toErrCodeMessage(2, 43, data), nil

	case 45: //s2c_miao_tech_worker_cd
		p := &pb2.S2CMiaoTechWorkerCdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_miao_tech_worker_cd PrintMsgProto &S2CMiaoTechWorkerCdProto fail")
		}

		return p, nil

	case 46: //s2c_fail_miao_tech_worker_cd
		return toErrCodeMessage(2, 46, data), nil

	case 47: //s2c_update_yuanbao
		p := &pb2.S2CUpdateYuanbaoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_yuanbao PrintMsgProto &S2CUpdateYuanbaoProto fail")
		}

		return p, nil

	case 169: //s2c_update_yuanbao_gift_limit
		p := &pb2.S2CUpdateYuanbaoGiftLimitProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_yuanbao_gift_limit PrintMsgProto &S2CUpdateYuanbaoGiftLimitProto fail")
		}

		return p, nil

	case 140: //s2c_update_dianquan
		p := &pb2.S2CUpdateDianquanProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_dianquan PrintMsgProto &S2CUpdateDianquanProto fail")
		}

		return p, nil

	case 141: //s2c_update_yinliang
		p := &pb2.S2CUpdateYinliangProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_yinliang PrintMsgProto &S2CUpdateYinliangProto fail")
		}

		return p, nil

	case 48: //s2c_update_hero_fight_amount
		p := &pb2.S2CUpdateHeroFightAmountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_hero_fight_amount PrintMsgProto &S2CUpdateHeroFightAmountProto fail")
		}

		return p, nil

	case 54: //s2c_recovery_forging_time_change
		p := &pb2.S2CRecoveryForgingTimeChangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_recovery_forging_time_change PrintMsgProto &S2CRecoveryForgingTimeChangeProto fail")
		}

		return p, nil

	case 52: //s2c_forging_equip
		p := &pb2.S2CForgingEquipProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_forging_equip PrintMsgProto &S2CForgingEquipProto fail")
		}

		return p, nil

	case 53: //s2c_fail_forging_equip
		return toErrCodeMessage(2, 53, data), nil

	case 117: //s2c_update_new_forging_pos
		p := &pb2.S2CUpdateNewForgingPosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_new_forging_pos PrintMsgProto &S2CUpdateNewForgingPosProto fail")
		}

		return p, nil

	case 67: //s2c_sign
		return toStringMessage(2, 67), nil

	case 68: //s2c_fail_sign
		return toErrCodeMessage(2, 68, data), nil

	case 70: //s2c_voice
		return toStringMessage(2, 70), nil

	case 71: //s2c_fail_voice
		return toErrCodeMessage(2, 71, data), nil

	case 57: //s2c_building_worker_time_changed
		p := &pb2.S2CBuildingWorkerTimeChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_building_worker_time_changed PrintMsgProto &S2CBuildingWorkerTimeChangedProto fail")
		}

		return p, nil

	case 58: //s2c_tech_worker_time_changed
		p := &pb2.S2CTechWorkerTimeChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_tech_worker_time_changed PrintMsgProto &S2CTechWorkerTimeChangedProto fail")
		}

		return p, nil

	case 59: //s2c_city_event_time_changed
		p := &pb2.S2CCityEventTimeChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_city_event_time_changed PrintMsgProto &S2CCityEventTimeChangedProto fail")
		}

		return p, nil

	case 61: //s2c_request_city_exchange_event
		p := &pb2.S2CRequestCityExchangeEventProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_request_city_exchange_event PrintMsgProto &S2CRequestCityExchangeEventProto fail")
		}

		return p, nil

	case 62: //s2c_fail_request_city_exchange_event
		return toErrCodeMessage(2, 62, data), nil

	case 64: //s2c_city_event_exchange
		p := &pb2.S2CCityEventExchangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_city_event_exchange PrintMsgProto &S2CCityEventExchangeProto fail")
		}

		return p, nil

	case 65: //s2c_fail_city_event_exchange
		return toErrCodeMessage(2, 65, data), nil

	case 72: //s2c_update_strategy_restore_start_time
		p := &pb2.S2CUpdateStrategyRestoreStartTimeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_strategy_restore_start_time PrintMsgProto &S2CUpdateStrategyRestoreStartTimeProto fail")
		}

		return p, nil

	case 73: //s2c_update_strategy_next_use_time
		p := &pb2.S2CUpdateStrategyNextUseTimeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_strategy_next_use_time PrintMsgProto &S2CUpdateStrategyNextUseTimeProto fail")
		}

		return p, nil

	case 85: //s2c_update_jade_ore
		p := &pb2.S2CUpdateJadeOreProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_jade_ore PrintMsgProto &S2CUpdateJadeOreProto fail")
		}

		return p, nil

	case 86: //s2c_update_jade
		p := &pb2.S2CUpdateJadeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_jade PrintMsgProto &S2CUpdateJadeProto fail")
		}

		return p, nil

	case 95: //s2c_change_head
		p := &pb2.S2CChangeHeadProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_change_head PrintMsgProto &S2CChangeHeadProto fail")
		}

		return p, nil

	case 96: //s2c_fail_change_head
		return toErrCodeMessage(2, 96, data), nil

	case 131: //s2c_change_body
		p := &pb2.S2CChangeBodyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_change_body PrintMsgProto &S2CChangeBodyProto fail")
		}

		return p, nil

	case 132: //s2c_fail_change_body
		return toErrCodeMessage(2, 132, data), nil

	case 115: //s2c_collect_countdown_prize
		p := &pb2.S2CCollectCountdownPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_collect_countdown_prize PrintMsgProto &S2CCollectCountdownPrizeProto fail")
		}

		p.Prize = nil

		return p, nil

	case 116: //s2c_fail_collect_countdown_prize
		return toErrCodeMessage(2, 116, data), nil

	case 118: //s2c_list_workshop_equipment
		p := &pb2.S2CListWorkshopEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_list_workshop_equipment PrintMsgProto &S2CListWorkshopEquipmentProto fail")
		}

		return p, nil

	case 120: //s2c_start_workshop
		p := &pb2.S2CStartWorkshopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_start_workshop PrintMsgProto &S2CStartWorkshopProto fail")
		}

		return p, nil

	case 121: //s2c_fail_start_workshop
		return toErrCodeMessage(2, 121, data), nil

	case 123: //s2c_collect_workshop
		p := &pb2.S2CCollectWorkshopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_collect_workshop PrintMsgProto &S2CCollectWorkshopProto fail")
		}

		return p, nil

	case 124: //s2c_fail_collect_workshop
		return toErrCodeMessage(2, 124, data), nil

	case 128: //s2c_workshop_miao_cd
		p := &pb2.S2CWorkshopMiaoCdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_workshop_miao_cd PrintMsgProto &S2CWorkshopMiaoCdProto fail")
		}

		return p, nil

	case 129: //s2c_fail_workshop_miao_cd
		return toErrCodeMessage(2, 129, data), nil

	case 134: //s2c_refresh_workshop
		return toStringMessage(2, 134), nil

	case 135: //s2c_fail_refresh_workshop
		return toErrCodeMessage(2, 135, data), nil

	case 137: //s2c_collect_season_prize
		return toStringMessage(2, 137), nil

	case 138: //s2c_fail_collect_season_prize
		return toErrCodeMessage(2, 138, data), nil

	case 139: //s2c_season_start_broadcast
		p := &pb2.S2CSeasonStartBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_season_start_broadcast PrintMsgProto &S2CSeasonStartBroadcastProto fail")
		}

		return p, nil

	case 145: //s2c_update_cost_reduce_coef
		p := &pb2.S2CUpdateCostReduceCoefProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_cost_reduce_coef PrintMsgProto &S2CUpdateCostReduceCoefProto fail")
		}

		return p, nil

	case 146: //s2c_update_sp
		p := &pb2.S2CUpdateSpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_sp PrintMsgProto &S2CUpdateSpProto fail")
		}

		return p, nil

	case 148: //s2c_buy_sp
		p := &pb2.S2CBuySpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_buy_sp PrintMsgProto &S2CBuySpProto fail")
		}

		return p, nil

	case 149: //s2c_fail_buy_sp
		return toErrCodeMessage(2, 149, data), nil

	case 151: //s2c_use_buf_effect
		p := &pb2.S2CUseBufEffectProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_use_buf_effect PrintMsgProto &S2CUseBufEffectProto fail")
		}

		return p, nil

	case 152: //s2c_fail_use_buf_effect
		return toErrCodeMessage(2, 152, data), nil

	case 155: //s2c_open_buf_effect_ui
		p := &pb2.S2COpenBufEffectUiProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_open_buf_effect_ui PrintMsgProto &S2COpenBufEffectUiProto fail")
		}

		return p, nil

	case 156: //s2c_fail_open_buf_effect_ui
		return toErrCodeMessage(2, 156, data), nil

	case 159: //s2c_use_advantage
		p := &pb2.S2CUseAdvantageProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_use_advantage PrintMsgProto &S2CUseAdvantageProto fail")
		}

		return p, nil

	case 160: //s2c_fail_use_advantage
		return toErrCodeMessage(2, 160, data), nil

	case 161: //s2c_update_advantage_count
		p := &pb2.S2CUpdateAdvantageCountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_update_advantage_count PrintMsgProto &S2CUpdateAdvantageCountProto fail")
		}

		return p, nil

	case 167: //s2c_worker_unlock
		p := &pb2.S2CWorkerUnlockProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_worker_unlock PrintMsgProto &S2CWorkerUnlockProto fail")
		}

		return p, nil

	case 168: //s2c_fail_worker_unlock
		return toErrCodeMessage(2, 168, data), nil

	case 165: //s2c_worker_always_unlock
		p := &pb2.S2CWorkerAlwaysUnlockProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "domestic.s2c_worker_always_unlock PrintMsgProto &S2CWorkerAlwaysUnlockProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_gm(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 6: //s2c_list_cmd
		p := &pb3.S2CListCmdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gm.s2c_list_cmd PrintMsgProto &S2CListCmdProto fail")
		}

		p.Datas = nil

		return p, nil

	case 2: //s2c_gm
		p := &pb3.S2CGmProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gm.s2c_gm PrintMsgProto &S2CGmProto fail")
		}

		return p, nil

	case 8: //s2c_invase_target_id
		p := &pb3.S2CInvaseTargetIdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gm.s2c_invase_target_id PrintMsgProto &S2CInvaseTargetIdProto fail")
		}

		p.TargetId = nil

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_military(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 1: //s2c_update_soldier_capcity
		p := &pb4.S2CUpdateSoldierCapcityProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_update_soldier_capcity PrintMsgProto &S2CUpdateSoldierCapcityProto fail")
		}

		return p, nil

	case 3: //s2c_recruit_soldier
		p := &pb4.S2CRecruitSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_recruit_soldier PrintMsgProto &S2CRecruitSoldierProto fail")
		}

		return p, nil

	case 4: //s2c_fail_recruit_soldier
		return toErrCodeMessage(4, 4, data), nil

	case 121: //s2c_recruit_soldier_v2
		p := &pb4.S2CRecruitSoldierV2Proto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_recruit_soldier_v2 PrintMsgProto &S2CRecruitSoldierV2Proto fail")
		}

		return p, nil

	case 122: //s2c_fail_recruit_soldier_v2
		return toErrCodeMessage(4, 122, data), nil

	case 124: //s2c_auto_recover_soldier
		p := &pb4.S2CAutoRecoverSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_auto_recover_soldier PrintMsgProto &S2CAutoRecoverSoldierProto fail")
		}

		return p, nil

	case 123: //s2c_recruit_soldier_times_changed
		p := &pb4.S2CRecruitSoldierTimesChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_recruit_soldier_times_changed PrintMsgProto &S2CRecruitSoldierTimesChangedProto fail")
		}

		return p, nil

	case 5: //s2c_add_wounded_soldier
		p := &pb4.S2CAddWoundedSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_add_wounded_soldier PrintMsgProto &S2CAddWoundedSoldierProto fail")
		}

		return p, nil

	case 7: //s2c_heal_wounded_soldier
		p := &pb4.S2CHealWoundedSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_heal_wounded_soldier PrintMsgProto &S2CHealWoundedSoldierProto fail")
		}

		return p, nil

	case 8: //s2c_fail_heal_wounded_soldier
		return toErrCodeMessage(4, 8, data), nil

	case 10: //s2c_captain_change_soldier
		p := &pb4.S2CCaptainChangeSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_change_soldier PrintMsgProto &S2CCaptainChangeSoldierProto fail")
		}

		return p, nil

	case 14: //s2c_fail_captain_change_soldier
		return toErrCodeMessage(4, 14, data), nil

	case 67: //s2c_captain_full_soldier
		p := &pb4.S2CCaptainFullSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_full_soldier PrintMsgProto &S2CCaptainFullSoldierProto fail")
		}

		return p, nil

	case 68: //s2c_fail_captain_full_soldier
		return toErrCodeMessage(4, 68, data), nil

	case 80: //s2c_update_free_soldier
		p := &pb4.S2CUpdateFreeSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_update_free_soldier PrintMsgProto &S2CUpdateFreeSoldierProto fail")
		}

		return p, nil

	case 150: //s2c_force_add_soldier
		p := &pb4.S2CForceAddSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_force_add_soldier PrintMsgProto &S2CForceAddSoldierProto fail")
		}

		return p, nil

	case 151: //s2c_fail_force_add_soldier
		return toErrCodeMessage(4, 151, data), nil

	case 11: //s2c_captain_change_data
		p := &pb4.S2CCaptainChangeDataProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_change_data PrintMsgProto &S2CCaptainChangeDataProto fail")
		}

		return p, nil

	case 13: //s2c_fight
		p := &pb4.S2CFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_fight PrintMsgProto &S2CFightProto fail")
		}

		p.Replay = nil

		return p, nil

	case 102: //s2c_multi_fight
		p := &pb4.S2CMultiFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_multi_fight PrintMsgProto &S2CMultiFightProto fail")
		}

		p.Replay = nil

		return p, nil

	case 199: //s2c_fightx
		p := &pb4.S2CFightxProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_fightx PrintMsgProto &S2CFightxProto fail")
		}

		p.Replay = nil

		return p, nil

	case 16: //s2c_upgrade_soldier_level
		p := &pb4.S2CUpgradeSoldierLevelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_upgrade_soldier_level PrintMsgProto &S2CUpgradeSoldierLevelProto fail")
		}

		return p, nil

	case 17: //s2c_fail_upgrade_soldier_level
		return toErrCodeMessage(4, 17, data), nil

	case 110: //s2c_recruit_captain_v2
		p := &pb4.S2CRecruitCaptainV2Proto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_recruit_captain_v2 PrintMsgProto &S2CRecruitCaptainV2Proto fail")
		}

		p.Captain = nil

		return p, nil

	case 111: //s2c_fail_recruit_captain_v2
		return toErrCodeMessage(4, 111, data), nil

	case 177: //s2c_random_captain_head
		p := &pb4.S2CRandomCaptainHeadProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_random_captain_head PrintMsgProto &S2CRandomCaptainHeadProto fail")
		}

		return p, nil

	case 147: //s2c_recruit_captain_seeker
		p := &pb4.S2CRecruitCaptainSeekerProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_recruit_captain_seeker PrintMsgProto &S2CRecruitCaptainSeekerProto fail")
		}

		p.Captain = nil

		return p, nil

	case 148: //s2c_fail_recruit_captain_seeker
		return toErrCodeMessage(4, 148, data), nil

	case 107: //s2c_set_defense_troop
		p := &pb4.S2CSetDefenseTroopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_set_defense_troop PrintMsgProto &S2CSetDefenseTroopProto fail")
		}

		return p, nil

	case 108: //s2c_fail_set_defense_troop
		return toErrCodeMessage(4, 108, data), nil

	case 128: //s2c_set_denfese_troop_defeated_mail
		p := &pb4.S2CSetDenfeseTroopDefeatedMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_set_denfese_troop_defeated_mail PrintMsgProto &S2CSetDenfeseTroopDefeatedMailProto fail")
		}

		p.Mail = nil

		return p, nil

	case 130: //s2c_clear_defense_troop_defeated_mail
		p := &pb4.S2CClearDefenseTroopDefeatedMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_clear_defense_troop_defeated_mail PrintMsgProto &S2CClearDefenseTroopDefeatedMailProto fail")
		}

		return p, nil

	case 189: //s2c_set_defenser_auto_full_soldier
		p := &pb4.S2CSetDefenserAutoFullSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_set_defenser_auto_full_soldier PrintMsgProto &S2CSetDefenserAutoFullSoldierProto fail")
		}

		return p, nil

	case 264: //s2c_fail_set_defenser_auto_full_soldier
		return toErrCodeMessage(4, 264, data), nil

	case 194: //s2c_use_copy_defenser_goods
		p := &pb4.S2CUseCopyDefenserGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_use_copy_defenser_goods PrintMsgProto &S2CUseCopyDefenserGoodsProto fail")
		}

		return p, nil

	case 195: //s2c_fail_use_copy_defenser_goods
		return toErrCodeMessage(4, 195, data), nil

	case 196: //s2c_update_copy_defenser
		p := &pb4.S2CUpdateCopyDefenserProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_update_copy_defenser PrintMsgProto &S2CUpdateCopyDefenserProto fail")
		}

		return p, nil

	case 197: //s2c_remove_copy_defenser
		return toStringMessage(4, 197), nil

	case 35: //s2c_sell_seek_captain
		p := &pb4.S2CSellSeekCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_sell_seek_captain PrintMsgProto &S2CSellSeekCaptainProto fail")
		}

		return p, nil

	case 37: //s2c_fail_sell_seek_captain
		return toErrCodeMessage(4, 37, data), nil

	case 46: //s2c_set_multi_captain_index
		p := &pb4.S2CSetMultiCaptainIndexProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_set_multi_captain_index PrintMsgProto &S2CSetMultiCaptainIndexProto fail")
		}

		return p, nil

	case 47: //s2c_fail_set_multi_captain_index
		return toErrCodeMessage(4, 47, data), nil

	case 144: //s2c_set_pve_captain
		p := &pb4.S2CSetPveCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_set_pve_captain PrintMsgProto &S2CSetPveCaptainProto fail")
		}

		p.Troop = nil

		return p, nil

	case 145: //s2c_fail_set_pve_captain
		return toErrCodeMessage(4, 145, data), nil

	case 39: //s2c_fire_captain
		p := &pb4.S2CFireCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_fire_captain PrintMsgProto &S2CFireCaptainProto fail")
		}

		return p, nil

	case 40: //s2c_fail_fire_captain
		return toErrCodeMessage(4, 40, data), nil

	case 49: //s2c_captain_refined
		p := &pb4.S2CCaptainRefinedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_refined PrintMsgProto &S2CCaptainRefinedProto fail")
		}

		return p, nil

	case 50: //s2c_fail_captain_refined
		return toErrCodeMessage(4, 50, data), nil

	case 207: //s2c_captain_enhance
		p := &pb4.S2CCaptainEnhanceProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_enhance PrintMsgProto &S2CCaptainEnhanceProto fail")
		}

		return p, nil

	case 208: //s2c_fail_captain_enhance
		return toErrCodeMessage(4, 208, data), nil

	case 51: //s2c_captain_refined_upgrade
		p := &pb4.S2CCaptainRefinedUpgradeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_refined_upgrade PrintMsgProto &S2CCaptainRefinedUpgradeProto fail")
		}

		p.Name = nil

		return p, nil

	case 184: //s2c_update_ability_exp
		p := &pb4.S2CUpdateAbilityExpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_update_ability_exp PrintMsgProto &S2CUpdateAbilityExpProto fail")
		}

		return p, nil

	case 52: //s2c_update_captain_exp
		p := &pb4.S2CUpdateCaptainExpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_update_captain_exp PrintMsgProto &S2CUpdateCaptainExpProto fail")
		}

		return p, nil

	case 53: //s2c_update_captain_level
		p := &pb4.S2CUpdateCaptainLevelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_update_captain_level PrintMsgProto &S2CUpdateCaptainLevelProto fail")
		}

		p.Name = nil

		return p, nil

	case 209: //s2c_captain_levelup
		p := &pb4.S2CCaptainLevelupProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_levelup PrintMsgProto &S2CCaptainLevelupProto fail")
		}

		return p, nil

	case 69: //s2c_update_captain_stat
		p := &pb4.S2CUpdateCaptainStatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_update_captain_stat PrintMsgProto &S2CUpdateCaptainStatProto fail")
		}

		p.TotalStat = nil

		return p, nil

	case 83: //s2c_change_captain_name
		p := &pb4.S2CChangeCaptainNameProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_change_captain_name PrintMsgProto &S2CChangeCaptainNameProto fail")
		}

		return p, nil

	case 84: //s2c_fail_change_captain_name
		return toErrCodeMessage(4, 84, data), nil

	case 86: //s2c_change_captain_race
		p := &pb4.S2CChangeCaptainRaceProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_change_captain_race PrintMsgProto &S2CChangeCaptainRaceProto fail")
		}

		p.Name = nil

		return p, nil

	case 87: //s2c_fail_change_captain_race
		return toErrCodeMessage(4, 87, data), nil

	case 90: //s2c_captain_rebirth_preview
		p := &pb4.S2CCaptainRebirthPreviewProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_rebirth_preview PrintMsgProto &S2CCaptainRebirthPreviewProto fail")
		}

		p.Name = nil

		p.AddStat = nil

		return p, nil

	case 91: //s2c_fail_captain_rebirth_preview
		return toErrCodeMessage(4, 91, data), nil

	case 161: //s2c_captain_rebirth_cd_start
		p := &pb4.S2CCaptainRebirthCdStartProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_rebirth_cd_start PrintMsgProto &S2CCaptainRebirthCdStartProto fail")
		}

		return p, nil

	case 93: //s2c_captain_rebirth
		p := &pb4.S2CCaptainRebirthProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_rebirth PrintMsgProto &S2CCaptainRebirthProto fail")
		}

		p.Name = nil

		p.TotalStat = nil

		return p, nil

	case 94: //s2c_fail_captain_rebirth
		return toErrCodeMessage(4, 94, data), nil

	case 211: //s2c_captain_progress
		p := &pb4.S2CCaptainProgressProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_progress PrintMsgProto &S2CCaptainProgressProto fail")
		}

		p.TotalStat = nil

		return p, nil

	case 212: //s2c_fail_captain_progress
		return toErrCodeMessage(4, 212, data), nil

	case 167: //s2c_captain_rebirth_miao_cd
		p := &pb4.S2CCaptainRebirthMiaoCdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_rebirth_miao_cd PrintMsgProto &S2CCaptainRebirthMiaoCdProto fail")
		}

		return p, nil

	case 168: //s2c_fail_captain_rebirth_miao_cd
		return toErrCodeMessage(4, 168, data), nil

	case 137: //s2c_collect_captain_training_exp
		p := &pb4.S2CCollectCaptainTrainingExpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_collect_captain_training_exp PrintMsgProto &S2CCollectCaptainTrainingExpProto fail")
		}

		return p, nil

	case 138: //s2c_fail_collect_captain_training_exp
		return toErrCodeMessage(4, 138, data), nil

	case 214: //s2c_captain_train_exp
		p := &pb4.S2CCaptainTrainExpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_train_exp PrintMsgProto &S2CCaptainTrainExpProto fail")
		}

		return p, nil

	case 215: //s2c_fail_captain_train_exp
		return toErrCodeMessage(4, 215, data), nil

	case 139: //s2c_update_training
		p := &pb4.S2CUpdateTrainingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_update_training PrintMsgProto &S2CUpdateTrainingProto fail")
		}

		return p, nil

	case 262: //s2c_captain_can_collect_exp
		p := &pb4.S2CCaptainCanCollectExpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_can_collect_exp PrintMsgProto &S2CCaptainCanCollectExpProto fail")
		}

		return p, nil

	case 263: //s2c_fail_captain_can_collect_exp
		return toErrCodeMessage(4, 263, data), nil

	case 141: //s2c_use_training_exp_goods
		p := &pb4.S2CUseTrainingExpGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_use_training_exp_goods PrintMsgProto &S2CUseTrainingExpGoodsProto fail")
		}

		return p, nil

	case 142: //s2c_fail_use_training_exp_goods
		return toErrCodeMessage(4, 142, data), nil

	case 217: //s2c_use_level_exp_goods
		p := &pb4.S2CUseLevelExpGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_use_level_exp_goods PrintMsgProto &S2CUseLevelExpGoodsProto fail")
		}

		return p, nil

	case 218: //s2c_fail_use_level_exp_goods
		return toErrCodeMessage(4, 218, data), nil

	case 244: //s2c_use_level_exp_goods2
		p := &pb4.S2CUseLevelExpGoods2Proto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_use_level_exp_goods2 PrintMsgProto &S2CUseLevelExpGoods2Proto fail")
		}

		return p, nil

	case 245: //s2c_fail_use_level_exp_goods2
		return toErrCodeMessage(4, 245, data), nil

	case 256: //s2c_auto_use_goods_until_captain_levelup
		p := &pb4.S2CAutoUseGoodsUntilCaptainLevelupProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_auto_use_goods_until_captain_levelup PrintMsgProto &S2CAutoUseGoodsUntilCaptainLevelupProto fail")
		}

		return p, nil

	case 257: //s2c_fail_auto_use_goods_until_captain_levelup
		return toErrCodeMessage(4, 257, data), nil

	case 75: //s2c_get_max_recruit_soldier
		p := &pb4.S2CGetMaxRecruitSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_get_max_recruit_soldier PrintMsgProto &S2CGetMaxRecruitSoldierProto fail")
		}

		return p, nil

	case 77: //s2c_get_max_heal_soldier
		p := &pb4.S2CGetMaxHealSoldierProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_get_max_heal_soldier PrintMsgProto &S2CGetMaxHealSoldierProto fail")
		}

		return p, nil

	case 113: //s2c_jiu_guan_consult
		p := &pb4.S2CJiuGuanConsultProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_jiu_guan_consult PrintMsgProto &S2CJiuGuanConsultProto fail")
		}

		p.Prize = nil

		return p, nil

	case 114: //s2c_fail_jiu_guan_consult
		return toErrCodeMessage(4, 114, data), nil

	case 115: //s2c_jiu_guan_consult_broadcast
		p := &pb4.S2CJiuGuanConsultBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_jiu_guan_consult_broadcast PrintMsgProto &S2CJiuGuanConsultBroadcastProto fail")
		}

		return p, nil

	case 116: //s2c_jiu_guan_times_changed
		p := &pb4.S2CJiuGuanTimesChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_jiu_guan_times_changed PrintMsgProto &S2CJiuGuanTimesChangedProto fail")
		}

		return p, nil

	case 118: //s2c_jiu_guan_refresh
		p := &pb4.S2CJiuGuanRefreshProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_jiu_guan_refresh PrintMsgProto &S2CJiuGuanRefreshProto fail")
		}

		return p, nil

	case 119: //s2c_fail_jiu_guan_refresh
		return toErrCodeMessage(4, 119, data), nil

	case 126: //s2c_unlock_captain_restraint_spell
		p := &pb4.S2CUnlockCaptainRestraintSpellProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_unlock_captain_restraint_spell PrintMsgProto &S2CUnlockCaptainRestraintSpellProto fail")
		}

		return p, nil

	case 127: //s2c_fail_unlock_captain_restraint_spell
		return toErrCodeMessage(4, 127, data), nil

	case 131: //s2c_new_troops
		p := &pb4.S2CNewTroopsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_new_troops PrintMsgProto &S2CNewTroopsProto fail")
		}

		p.Troop = nil

		return p, nil

	case 133: //s2c_get_captain_stat_details
		p := &pb4.S2CGetCaptainStatDetailsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_get_captain_stat_details PrintMsgProto &S2CGetCaptainStatDetailsProto fail")
		}

		p.Stats = nil

		return p, nil

	case 134: //s2c_fail_get_captain_stat_details
		return toErrCodeMessage(4, 134, data), nil

	case 220: //s2c_captain_stat_details
		p := &pb4.S2CCaptainStatDetailsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_stat_details PrintMsgProto &S2CCaptainStatDetailsProto fail")
		}

		p.Stats = nil

		return p, nil

	case 221: //s2c_fail_captain_stat_details
		return toErrCodeMessage(4, 221, data), nil

	case 135: //s2c_update_troop_fight_amount
		p := &pb4.S2CUpdateTroopFightAmountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_update_troop_fight_amount PrintMsgProto &S2CUpdateTroopFightAmountProto fail")
		}

		return p, nil

	case 170: //s2c_update_captain_official
		p := &pb4.S2CUpdateCaptainOfficialProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_update_captain_official PrintMsgProto &S2CUpdateCaptainOfficialProto fail")
		}

		return p, nil

	case 171: //s2c_fail_update_captain_official
		return toErrCodeMessage(4, 171, data), nil

	case 223: //s2c_set_captain_official
		p := &pb4.S2CSetCaptainOfficialProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_set_captain_official PrintMsgProto &S2CSetCaptainOfficialProto fail")
		}

		return p, nil

	case 224: //s2c_fail_set_captain_official
		return toErrCodeMessage(4, 224, data), nil

	case 173: //s2c_leave_captain_official
		p := &pb4.S2CLeaveCaptainOfficialProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_leave_captain_official PrintMsgProto &S2CLeaveCaptainOfficialProto fail")
		}

		return p, nil

	case 174: //s2c_fail_leave_captain_official
		return toErrCodeMessage(4, 174, data), nil

	case 175: //s2c_add_gongxun
		p := &pb4.S2CAddGongxunProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_add_gongxun PrintMsgProto &S2CAddGongxunProto fail")
		}

		return p, nil

	case 186: //s2c_use_gong_xun_goods
		p := &pb4.S2CUseGongXunGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_use_gong_xun_goods PrintMsgProto &S2CUseGongXunGoodsProto fail")
		}

		return p, nil

	case 187: //s2c_fail_use_gong_xun_goods
		return toErrCodeMessage(4, 187, data), nil

	case 229: //s2c_use_gongxun_goods
		p := &pb4.S2CUseGongxunGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_use_gongxun_goods PrintMsgProto &S2CUseGongxunGoodsProto fail")
		}

		return p, nil

	case 230: //s2c_fail_use_gongxun_goods
		return toErrCodeMessage(4, 230, data), nil

	case 182: //s2c_close_fight_guide
		p := &pb4.S2CCloseFightGuideProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_close_fight_guide PrintMsgProto &S2CCloseFightGuideProto fail")
		}

		return p, nil

	case 183: //s2c_fail_close_fight_guide
		return toErrCodeMessage(4, 183, data), nil

	case 191: //s2c_view_other_hero_captain
		p := &pb4.S2CViewOtherHeroCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_view_other_hero_captain PrintMsgProto &S2CViewOtherHeroCaptainProto fail")
		}

		p.HeroId = nil

		p.Captain = nil

		return p, nil

	case 192: //s2c_fail_view_other_hero_captain
		return toErrCodeMessage(4, 192, data), nil

	case 232: //s2c_captain_born
		p := &pb4.S2CCaptainBornProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_born PrintMsgProto &S2CCaptainBornProto fail")
		}

		p.Captain = nil

		return p, nil

	case 233: //s2c_fail_captain_born
		return toErrCodeMessage(4, 233, data), nil

	case 235: //s2c_captain_upstar
		p := &pb4.S2CCaptainUpstarProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_upstar PrintMsgProto &S2CCaptainUpstarProto fail")
		}

		return p, nil

	case 236: //s2c_fail_captain_upstar
		return toErrCodeMessage(4, 236, data), nil

	case 269: //s2c_captain_exchange
		p := &pb4.S2CCaptainExchangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_captain_exchange PrintMsgProto &S2CCaptainExchangeProto fail")
		}

		p.Cap1 = nil

		p.Cap2 = nil

		return p, nil

	case 270: //s2c_fail_captain_exchange
		return toErrCodeMessage(4, 270, data), nil

	case 253: //s2c_notice_captain_has_viewed
		return toStringMessage(4, 253), nil

	case 254: //s2c_fail_notice_captain_has_viewed
		return toErrCodeMessage(4, 254, data), nil

	case 266: //s2c_activate_captain_friendship
		p := &pb4.S2CActivateCaptainFriendshipProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_activate_captain_friendship PrintMsgProto &S2CActivateCaptainFriendshipProto fail")
		}

		return p, nil

	case 267: //s2c_fail_activate_captain_friendship
		return toErrCodeMessage(4, 267, data), nil

	case 271: //s2c_show_prize_captain
		p := &pb4.S2CShowPrizeCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_show_prize_captain PrintMsgProto &S2CShowPrizeCaptainProto fail")
		}

		return p, nil

	case 273: //s2c_notice_official_has_viewed
		p := &pb4.S2CNoticeOfficialHasViewedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "military.s2c_notice_official_has_viewed PrintMsgProto &S2CNoticeOfficialHasViewedProto fail")
		}

		return p, nil

	case 274: //s2c_fail_notice_official_has_viewed
		return toErrCodeMessage(4, 274, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_misc(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 4: //s2c_config
		p := &pb5.S2CConfigProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_config PrintMsgProto &S2CConfigProto fail")
		}

		p.Config = nil

		return p, nil

	case 77: //s2c_configlua
		p := &pb5.S2CConfigluaProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_configlua PrintMsgProto &S2CConfigluaProto fail")
		}

		p.Config = nil

		return p, nil

	case 5: //s2c_fail_disconect_reason
		return toErrCodeMessage(5, 5, data), nil

	case 6: //s2c_reset_daily
		return toStringMessage(5, 6), nil

	case 9: //s2c_sync_time
		p := &pb5.S2CSyncTimeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_sync_time PrintMsgProto &S2CSyncTimeProto fail")
		}

		return p, nil

	case 11: //s2c_block
		p := &pb5.S2CBlockProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_block PrintMsgProto &S2CBlockProto fail")
		}

		p.Data = nil

		return p, nil

	case 12: //s2c_open_function
		p := &pb5.S2COpenFunctionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_open_function PrintMsgProto &S2COpenFunctionProto fail")
		}

		return p, nil

	case 34: //s2c_open_multi_function
		p := &pb5.S2COpenMultiFunctionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_open_multi_function PrintMsgProto &S2COpenMultiFunctionProto fail")
		}

		return p, nil

	case 14: //s2c_set_hero_bool
		p := &pb5.S2CSetHeroBoolProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_set_hero_bool PrintMsgProto &S2CSetHeroBoolProto fail")
		}

		return p, nil

	case 30: //s2c_reset_hero_bool
		p := &pb5.S2CResetHeroBoolProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_reset_hero_bool PrintMsgProto &S2CResetHeroBoolProto fail")
		}

		return p, nil

	case 13: //s2c_screen_show_words
		p := &pb5.S2CScreenShowWordsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_screen_show_words PrintMsgProto &S2CScreenShowWordsProto fail")
		}

		return p, nil

	case 16: //s2c_ping
		return toStringMessage(5, 16), nil

	case 18: //s2c_client_version
		p := &pb5.S2CClientVersionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_client_version PrintMsgProto &S2CClientVersionProto fail")
		}

		return p, nil

	case 20: //s2c_update_pf_token
		return toStringMessage(5, 20), nil

	case 22: //s2c_settings
		p := &pb5.S2CSettingsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_settings PrintMsgProto &S2CSettingsProto fail")
		}

		return p, nil

	case 23: //s2c_fail_settings
		return toErrCodeMessage(5, 23, data), nil

	case 25: //s2c_settings_to_default
		p := &pb5.S2CSettingsToDefaultProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_settings_to_default PrintMsgProto &S2CSettingsToDefaultProto fail")
		}

		return p, nil

	case 26: //s2c_hero_broadcast
		p := &pb5.S2CHeroBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_hero_broadcast PrintMsgProto &S2CHeroBroadcastProto fail")
		}

		return p, nil

	case 28: //s2c_sys_timing_broadcast
		p := &pb5.S2CSysTimingBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_sys_timing_broadcast PrintMsgProto &S2CSysTimingBroadcastProto fail")
		}

		return p, nil

	case 27: //s2c_sys_broadcast
		p := &pb5.S2CSysBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_sys_broadcast PrintMsgProto &S2CSysBroadcastProto fail")
		}

		return p, nil

	case 32: //s2c_update_location
		p := &pb5.S2CUpdateLocationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_update_location PrintMsgProto &S2CUpdateLocationProto fail")
		}

		return p, nil

	case 33: //s2c_fail_update_location
		return toErrCodeMessage(5, 33, data), nil

	case 43: //s2c_collect_charge_prize
		p := &pb5.S2CCollectChargePrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_collect_charge_prize PrintMsgProto &S2CCollectChargePrizeProto fail")
		}

		p.Prize = nil

		return p, nil

	case 44: //s2c_fail_collect_charge_prize
		return toErrCodeMessage(5, 44, data), nil

	case 45: //s2c_update_charge_amount
		p := &pb5.S2CUpdateChargeAmountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_update_charge_amount PrintMsgProto &S2CUpdateChargeAmountProto fail")
		}

		return p, nil

	case 47: //s2c_collect_daily_bargain
		p := &pb5.S2CCollectDailyBargainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_collect_daily_bargain PrintMsgProto &S2CCollectDailyBargainProto fail")
		}

		p.Prize = nil

		return p, nil

	case 48: //s2c_fail_collect_daily_bargain
		return toErrCodeMessage(5, 48, data), nil

	case 52: //s2c_activate_duration_card
		p := &pb5.S2CActivateDurationCardProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_activate_duration_card PrintMsgProto &S2CActivateDurationCardProto fail")
		}

		p.Prize = nil

		return p, nil

	case 53: //s2c_fail_activate_duration_card
		return toErrCodeMessage(5, 53, data), nil

	case 55: //s2c_collect_duration_card_daily_prize
		p := &pb5.S2CCollectDurationCardDailyPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_collect_duration_card_daily_prize PrintMsgProto &S2CCollectDurationCardDailyPrizeProto fail")
		}

		p.Prize = nil

		return p, nil

	case 56: //s2c_fail_collect_duration_card_daily_prize
		return toErrCodeMessage(5, 56, data), nil

	case 62: //s2c_reset_daily_zero
		return toStringMessage(5, 62), nil

	case 72: //s2c_reset_daily_mc
		return toStringMessage(5, 72), nil

	case 65: //s2c_set_privacy_setting
		p := &pb5.S2CSetPrivacySettingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_set_privacy_setting PrintMsgProto &S2CSetPrivacySettingProto fail")
		}

		return p, nil

	case 66: //s2c_fail_set_privacy_setting
		return toErrCodeMessage(5, 66, data), nil

	case 68: //s2c_fail_set_default_privacy_settings
		return toErrCodeMessage(5, 68, data), nil

	case 70: //s2c_get_product_info
		p := &pb5.S2CGetProductInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_get_product_info PrintMsgProto &S2CGetProductInfoProto fail")
		}

		return p, nil

	case 71: //s2c_fail_get_product_info
		return toErrCodeMessage(5, 71, data), nil

	case 73: //s2c_reset_weekly
		return toStringMessage(5, 73), nil

	case 74: //s2c_update_first_recharge
		p := &pb5.S2CUpdateFirstRechargeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_update_first_recharge PrintMsgProto &S2CUpdateFirstRechargeProto fail")
		}

		return p, nil

	case 75: //s2c_update_buff_notice
		p := &pb5.S2CUpdateBuffNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "misc.s2c_update_buff_notice PrintMsgProto &S2CUpdateBuffNoticeProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_region(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 159: //s2c_update_map_radius
		p := &pb7.S2CUpdateMapRadiusProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_map_radius PrintMsgProto &S2CUpdateMapRadiusProto fail")
		}

		return p, nil

	case 149: //s2c_update_self_view
		p := &pb7.S2CUpdateSelfViewProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_self_view PrintMsgProto &S2CUpdateSelfViewProto fail")
		}

		return p, nil

	case 151: //s2c_close_view
		return toStringMessage(7, 151), nil

	case 152: //s2c_add_base_unit
		p := &pb7.S2CAddBaseUnitProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_add_base_unit PrintMsgProto &S2CAddBaseUnitProto fail")
		}

		p.Data = nil

		return p, nil

	case 110: //s2c_update_npc_base_info
		p := &pb7.S2CUpdateNpcBaseInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_npc_base_info PrintMsgProto &S2CUpdateNpcBaseInfoProto fail")
		}

		p.NpcId = nil

		p.Hero = nil

		return p, nil

	case 213: //s2c_update_base_progress
		p := &pb7.S2CUpdateBaseProgressProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_base_progress PrintMsgProto &S2CUpdateBaseProgressProto fail")
		}

		p.Id = nil

		return p, nil

	case 153: //s2c_remove_base_unit
		p := &pb7.S2CRemoveBaseUnitProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_remove_base_unit PrintMsgProto &S2CRemoveBaseUnitProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 95: //s2c_pre_invasion_target
		p := &pb7.S2CPreInvasionTargetProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_pre_invasion_target PrintMsgProto &S2CPreInvasionTargetProto fail")
		}

		return p, nil

	case 96: //s2c_fail_pre_invasion_target
		return toErrCodeMessage(7, 96, data), nil

	case 162: //s2c_watch_base_unit
		p := &pb7.S2CWatchBaseUnitProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_watch_base_unit PrintMsgProto &S2CWatchBaseUnitProto fail")
		}

		p.Target = nil

		p.Hero = nil

		return p, nil

	case 163: //s2c_fail_watch_base_unit
		return toErrCodeMessage(7, 163, data), nil

	case 164: //s2c_update_watch_base_prosperity
		p := &pb7.S2CUpdateWatchBaseProsperityProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_watch_base_prosperity PrintMsgProto &S2CUpdateWatchBaseProsperityProto fail")
		}

		p.Target = nil

		return p, nil

	case 174: //s2c_update_stop_lost_prosperity
		p := &pb7.S2CUpdateStopLostProsperityProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_stop_lost_prosperity PrintMsgProto &S2CUpdateStopLostProsperityProto fail")
		}

		p.Target = nil

		return p, nil

	case 154: //s2c_add_troop_unit
		p := &pb7.S2CAddTroopUnitProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_add_troop_unit PrintMsgProto &S2CAddTroopUnitProto fail")
		}

		p.Data = nil

		return p, nil

	case 155: //s2c_remove_troop_unit
		p := &pb7.S2CRemoveTroopUnitProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_remove_troop_unit PrintMsgProto &S2CRemoveTroopUnitProto fail")
		}

		p.TroopId = nil

		return p, nil

	case 157: //s2c_request_troop_unit
		p := &pb7.S2CRequestTroopUnitProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_request_troop_unit PrintMsgProto &S2CRequestTroopUnitProto fail")
		}

		p.Data = nil

		return p, nil

	case 158: //s2c_fail_request_troop_unit
		return toErrCodeMessage(7, 158, data), nil

	case 134: //s2c_add_ruins_base
		p := &pb7.S2CAddRuinsBaseProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_add_ruins_base PrintMsgProto &S2CAddRuinsBaseProto fail")
		}

		return p, nil

	case 130: //s2c_remove_ruins_base
		p := &pb7.S2CRemoveRuinsBaseProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_remove_ruins_base PrintMsgProto &S2CRemoveRuinsBaseProto fail")
		}

		return p, nil

	case 132: //s2c_request_ruins_base
		p := &pb7.S2CRequestRuinsBaseProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_request_ruins_base PrintMsgProto &S2CRequestRuinsBaseProto fail")
		}

		p.HeroBasic = nil

		return p, nil

	case 133: //s2c_fail_request_ruins_base
		return toErrCodeMessage(7, 133, data), nil

	case 90: //s2c_update_self_mian_disappear_time
		p := &pb7.S2CUpdateSelfMianDisappearTimeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_self_mian_disappear_time PrintMsgProto &S2CUpdateSelfMianDisappearTimeProto fail")
		}

		return p, nil

	case 92: //s2c_use_mian_goods
		p := &pb7.S2CUseMianGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_use_mian_goods PrintMsgProto &S2CUseMianGoodsProto fail")
		}

		return p, nil

	case 93: //s2c_fail_use_mian_goods
		return toErrCodeMessage(7, 93, data), nil

	case 160: //s2c_update_new_hero_mian_disappear_time
		p := &pb7.S2CUpdateNewHeroMianDisappearTimeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_new_hero_mian_disappear_time PrintMsgProto &S2CUpdateNewHeroMianDisappearTimeProto fail")
		}

		return p, nil

	case 41: //s2c_upgrade_base
		return toStringMessage(7, 41), nil

	case 42: //s2c_fail_upgrade_base
		return toErrCodeMessage(7, 42, data), nil

	case 52: //s2c_self_update_base_level
		p := &pb7.S2CSelfUpdateBaseLevelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_self_update_base_level PrintMsgProto &S2CSelfUpdateBaseLevelProto fail")
		}

		return p, nil

	case 85: //s2c_update_white_flag
		p := &pb7.S2CUpdateWhiteFlagProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_white_flag PrintMsgProto &S2CUpdateWhiteFlagProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 87: //s2c_white_flag_detail
		p := &pb7.S2CWhiteFlagDetailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_white_flag_detail PrintMsgProto &S2CWhiteFlagDetailProto fail")
		}

		p.HeroId = nil

		p.WhiteFlagHeroId = nil

		return p, nil

	case 88: //s2c_fail_white_flag_detail
		return toErrCodeMessage(7, 88, data), nil

	case 135: //s2c_self_base_destroy
		p := &pb7.S2CSelfBaseDestroyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_self_base_destroy PrintMsgProto &S2CSelfBaseDestroyProto fail")
		}

		return p, nil

	case 105: //s2c_prosperity_buf
		p := &pb7.S2CProsperityBufProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_prosperity_buf PrintMsgProto &S2CProsperityBufProto fail")
		}

		return p, nil

	case 122: //s2c_show_words
		p := &pb7.S2CShowWordsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_show_words PrintMsgProto &S2CShowWordsProto fail")
		}

		p.BaseId = nil

		p.TroopId = nil

		return p, nil

	case 212: //s2c_get_buy_prosperity_cost
		p := &pb7.S2CGetBuyProsperityCostProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_get_buy_prosperity_cost PrintMsgProto &S2CGetBuyProsperityCostProto fail")
		}

		return p, nil

	case 107: //s2c_buy_prosperity
		p := &pb7.S2CBuyProsperityProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_buy_prosperity PrintMsgProto &S2CBuyProsperityProto fail")
		}

		return p, nil

	case 108: //s2c_fail_buy_prosperity
		return toErrCodeMessage(7, 108, data), nil

	case 128: //s2c_self_been_attack_rob_changed
		p := &pb7.S2CSelfBeenAttackRobChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_self_been_attack_rob_changed PrintMsgProto &S2CSelfBeenAttackRobChangedProto fail")
		}

		return p, nil

	case 129: //s2c_guild_been_attack_rob_changed
		p := &pb7.S2CGuildBeenAttackRobChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_guild_been_attack_rob_changed PrintMsgProto &S2CGuildBeenAttackRobChangedProto fail")
		}

		return p, nil

	case 47: //s2c_switch_action
		p := &pb7.S2CSwitchActionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_switch_action PrintMsgProto &S2CSwitchActionProto fail")
		}

		return p, nil

	case 166: //s2c_request_military_push
		p := &pb7.S2CRequestMilitaryPushProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_request_military_push PrintMsgProto &S2CRequestMilitaryPushProto fail")
		}

		p.ToTarget = nil

		p.FromTarget = nil

		return p, nil

	case 167: //s2c_fail_request_military_push
		return toErrCodeMessage(7, 167, data), nil

	case 22: //s2c_update_military_info
		p := &pb7.S2CUpdateMilitaryInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_military_info PrintMsgProto &S2CUpdateMilitaryInfoProto fail")
		}

		p.Data = nil

		return p, nil

	case 23: //s2c_remove_military_info
		p := &pb7.S2CRemoveMilitaryInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_remove_military_info PrintMsgProto &S2CRemoveMilitaryInfoProto fail")
		}

		p.Id = nil

		return p, nil

	case 50: //s2c_update_self_military_info
		p := &pb7.S2CUpdateSelfMilitaryInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_self_military_info PrintMsgProto &S2CUpdateSelfMilitaryInfoProto fail")
		}

		p.TroopId = nil

		p.Data = nil

		return p, nil

	case 51: //s2c_remove_self_military_info
		p := &pb7.S2CRemoveSelfMilitaryInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_remove_self_military_info PrintMsgProto &S2CRemoveSelfMilitaryInfoProto fail")
		}

		p.Id = nil

		return p, nil

	case 109: //s2c_npc_base_info
		p := &pb7.S2CNpcBaseInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_npc_base_info PrintMsgProto &S2CNpcBaseInfoProto fail")
		}

		p.NpcId = nil

		return p, nil

	case 2: //s2c_create_base
		p := &pb7.S2CCreateBaseProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_create_base PrintMsgProto &S2CCreateBaseProto fail")
		}

		return p, nil

	case 3: //s2c_fail_create_base
		return toErrCodeMessage(7, 3, data), nil

	case 15: //s2c_fast_move_base
		p := &pb7.S2CFastMoveBaseProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_fast_move_base PrintMsgProto &S2CFastMoveBaseProto fail")
		}

		return p, nil

	case 16: //s2c_fail_fast_move_base
		return toErrCodeMessage(7, 16, data), nil

	case 25: //s2c_invasion
		p := &pb7.S2CInvasionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_invasion PrintMsgProto &S2CInvasionProto fail")
		}

		p.Target = nil

		return p, nil

	case 34: //s2c_fail_invasion
		return toErrCodeMessage(7, 34, data), nil

	case 56: //s2c_update_self_troops
		p := &pb7.S2CUpdateSelfTroopsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_self_troops PrintMsgProto &S2CUpdateSelfTroopsProto fail")
		}

		return p, nil

	case 210: //s2c_update_self_troops_outside
		p := &pb7.S2CUpdateSelfTroopsOutsideProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_self_troops_outside PrintMsgProto &S2CUpdateSelfTroopsOutsideProto fail")
		}

		return p, nil

	case 28: //s2c_cancel_invasion
		p := &pb7.S2CCancelInvasionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_cancel_invasion PrintMsgProto &S2CCancelInvasionProto fail")
		}

		p.Id = nil

		return p, nil

	case 29: //s2c_fail_cancel_invasion
		return toErrCodeMessage(7, 29, data), nil

	case 72: //s2c_repatriate
		p := &pb7.S2CRepatriateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_repatriate PrintMsgProto &S2CRepatriateProto fail")
		}

		p.Id = nil

		return p, nil

	case 73: //s2c_fail_repatriate
		return toErrCodeMessage(7, 73, data), nil

	case 187: //s2c_baoz_repatriate
		p := &pb7.S2CBaozRepatriateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_baoz_repatriate PrintMsgProto &S2CBaozRepatriateProto fail")
		}

		p.BaseId = nil

		p.TroopId = nil

		return p, nil

	case 188: //s2c_fail_baoz_repatriate
		return toErrCodeMessage(7, 188, data), nil

	case 140: //s2c_speed_up
		p := &pb7.S2CSpeedUpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_speed_up PrintMsgProto &S2CSpeedUpProto fail")
		}

		p.Id = nil

		return p, nil

	case 141: //s2c_fail_speed_up
		return toErrCodeMessage(7, 141, data), nil

	case 31: //s2c_expel
		p := &pb7.S2CExpelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_expel PrintMsgProto &S2CExpelProto fail")
		}

		p.Id = nil

		return p, nil

	case 32: //s2c_fail_expel
		return toErrCodeMessage(7, 32, data), nil

	case 100: //s2c_favorite_pos
		p := &pb7.S2CFavoritePosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_favorite_pos PrintMsgProto &S2CFavoritePosProto fail")
		}

		return p, nil

	case 101: //s2c_fail_favorite_pos
		return toErrCodeMessage(7, 101, data), nil

	case 103: //s2c_favorite_pos_list
		p := &pb7.S2CFavoritePosListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_favorite_pos_list PrintMsgProto &S2CFavoritePosListProto fail")
		}

		p.Data = nil

		return p, nil

	case 104: //s2c_fail_favorite_pos_list
		return toErrCodeMessage(7, 104, data), nil

	case 176: //s2c_get_prev_investigate
		p := &pb7.S2CGetPrevInvestigateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_get_prev_investigate PrintMsgProto &S2CGetPrevInvestigateProto fail")
		}

		p.HeroId = nil

		p.MailId = nil

		return p, nil

	case 143: //s2c_investigate
		p := &pb7.S2CInvestigateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_investigate PrintMsgProto &S2CInvestigateProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 144: //s2c_fail_investigate
		return toErrCodeMessage(7, 144, data), nil

	case 235: //s2c_investigate_invade
		p := &pb7.S2CInvestigateInvadeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_investigate_invade PrintMsgProto &S2CInvestigateInvadeProto fail")
		}

		p.Target = nil

		return p, nil

	case 236: //s2c_fail_investigate_invade
		return toErrCodeMessage(7, 236, data), nil

	case 168: //s2c_update_multi_level_npc_pass_level
		p := &pb7.S2CUpdateMultiLevelNpcPassLevelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_multi_level_npc_pass_level PrintMsgProto &S2CUpdateMultiLevelNpcPassLevelProto fail")
		}

		return p, nil

	case 169: //s2c_update_multi_level_npc_hate
		p := &pb7.S2CUpdateMultiLevelNpcHateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_multi_level_npc_hate PrintMsgProto &S2CUpdateMultiLevelNpcHateProto fail")
		}

		return p, nil

	case 170: //s2c_update_multi_level_npc_times
		p := &pb7.S2CUpdateMultiLevelNpcTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_multi_level_npc_times PrintMsgProto &S2CUpdateMultiLevelNpcTimesProto fail")
		}

		return p, nil

	case 184: //s2c_use_multi_level_npc_times_goods
		p := &pb7.S2CUseMultiLevelNpcTimesGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_use_multi_level_npc_times_goods PrintMsgProto &S2CUseMultiLevelNpcTimesGoodsProto fail")
		}

		return p, nil

	case 185: //s2c_fail_use_multi_level_npc_times_goods
		return toErrCodeMessage(7, 185, data), nil

	case 189: //s2c_update_invase_hero_times
		p := &pb7.S2CUpdateInvaseHeroTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_invase_hero_times PrintMsgProto &S2CUpdateInvaseHeroTimesProto fail")
		}

		return p, nil

	case 209: //s2c_update_jun_tuan_npc_times
		p := &pb7.S2CUpdateJunTuanNpcTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_jun_tuan_npc_times PrintMsgProto &S2CUpdateJunTuanNpcTimesProto fail")
		}

		return p, nil

	case 191: //s2c_use_invase_hero_times_goods
		p := &pb7.S2CUseInvaseHeroTimesGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_use_invase_hero_times_goods PrintMsgProto &S2CUseInvaseHeroTimesGoodsProto fail")
		}

		return p, nil

	case 192: //s2c_fail_use_invase_hero_times_goods
		return toErrCodeMessage(7, 192, data), nil

	case 173: //s2c_calc_move_speed
		p := &pb7.S2CCalcMoveSpeedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_calc_move_speed PrintMsgProto &S2CCalcMoveSpeedProto fail")
		}

		p.Id = nil

		return p, nil

	case 179: //s2c_list_enemy_pos
		p := &pb7.S2CListEnemyPosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_list_enemy_pos PrintMsgProto &S2CListEnemyPosProto fail")
		}

		return p, nil

	case 181: //s2c_search_baoz_npc
		p := &pb7.S2CSearchBaozNpcProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_search_baoz_npc PrintMsgProto &S2CSearchBaozNpcProto fail")
		}

		p.BaseId = nil

		return p, nil

	case 182: //s2c_fail_search_baoz_npc
		return toErrCodeMessage(7, 182, data), nil

	case 194: //s2c_home_ast_defending_info
		p := &pb7.S2CHomeAstDefendingInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_home_ast_defending_info PrintMsgProto &S2CHomeAstDefendingInfoProto fail")
		}

		return p, nil

	case 195: //s2c_fail_home_ast_defending_info
		return toErrCodeMessage(7, 195, data), nil

	case 197: //s2c_guild_please_help_me
		return toStringMessage(7, 197), nil

	case 198: //s2c_fail_guild_please_help_me
		return toErrCodeMessage(7, 198, data), nil

	case 200: //s2c_create_assembly
		p := &pb7.S2CCreateAssemblyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_create_assembly PrintMsgProto &S2CCreateAssemblyProto fail")
		}

		p.Target = nil

		p.Id = nil

		return p, nil

	case 201: //s2c_fail_create_assembly
		return toErrCodeMessage(7, 201, data), nil

	case 203: //s2c_show_assembly
		p := &pb7.S2CShowAssemblyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_show_assembly PrintMsgProto &S2CShowAssemblyProto fail")
		}

		p.Id = nil

		p.Data = nil

		return p, nil

	case 205: //s2c_show_assembly_changed
		p := &pb7.S2CShowAssemblyChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_show_assembly_changed PrintMsgProto &S2CShowAssemblyChangedProto fail")
		}

		p.Id = nil

		return p, nil

	case 207: //s2c_join_assembly
		p := &pb7.S2CJoinAssemblyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_join_assembly PrintMsgProto &S2CJoinAssemblyProto fail")
		}

		p.Id = nil

		return p, nil

	case 208: //s2c_fail_join_assembly
		return toErrCodeMessage(7, 208, data), nil

	case 215: //s2c_create_guild_workshop
		p := &pb7.S2CCreateGuildWorkshopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_create_guild_workshop PrintMsgProto &S2CCreateGuildWorkshopProto fail")
		}

		return p, nil

	case 216: //s2c_fail_create_guild_workshop
		return toErrCodeMessage(7, 216, data), nil

	case 218: //s2c_show_guild_workshop
		p := &pb7.S2CShowGuildWorkshopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_show_guild_workshop PrintMsgProto &S2CShowGuildWorkshopProto fail")
		}

		p.BaseId = nil

		return p, nil

	case 220: //s2c_hurt_guild_workshop
		p := &pb7.S2CHurtGuildWorkshopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_hurt_guild_workshop PrintMsgProto &S2CHurtGuildWorkshopProto fail")
		}

		p.BaseId = nil

		return p, nil

	case 221: //s2c_fail_hurt_guild_workshop
		return toErrCodeMessage(7, 221, data), nil

	case 224: //s2c_remove_guild_workshop
		return toStringMessage(7, 224), nil

	case 225: //s2c_fail_remove_guild_workshop
		return toErrCodeMessage(7, 225, data), nil

	case 222: //s2c_update_guild_workshop_prize_count
		p := &pb7.S2CUpdateGuildWorkshopPrizeCountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_guild_workshop_prize_count PrintMsgProto &S2CUpdateGuildWorkshopPrizeCountProto fail")
		}

		return p, nil

	case 226: //s2c_update_hero_build_workshop_times
		p := &pb7.S2CUpdateHeroBuildWorkshopTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_hero_build_workshop_times PrintMsgProto &S2CUpdateHeroBuildWorkshopTimesProto fail")
		}

		return p, nil

	case 227: //s2c_update_hero_output_workshop_times
		p := &pb7.S2CUpdateHeroOutputWorkshopTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_update_hero_output_workshop_times PrintMsgProto &S2CUpdateHeroOutputWorkshopTimesProto fail")
		}

		return p, nil

	case 229: //s2c_catch_guild_workshop_logs
		p := &pb7.S2CCatchGuildWorkshopLogsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_catch_guild_workshop_logs PrintMsgProto &S2CCatchGuildWorkshopLogsProto fail")
		}

		return p, nil

	case 230: //s2c_fail_catch_guild_workshop_logs
		return toErrCodeMessage(7, 230, data), nil

	case 233: //s2c_get_self_baoz
		p := &pb7.S2CGetSelfBaozProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "region.s2c_get_self_baoz PrintMsgProto &S2CGetSelfBaozProto fail")
		}

		p.BaseId = nil

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_mail(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_list_mail
		p := &pb8.S2CListMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.s2c_list_mail PrintMsgProto &S2CListMailProto fail")
		}

		p.Mail = nil

		return p, nil

	case 3: //s2c_fail_list_mail
		return toErrCodeMessage(8, 3, data), nil

	case 4: //s2c_receive_mail
		p := &pb8.S2CReceiveMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.s2c_receive_mail PrintMsgProto &S2CReceiveMailProto fail")
		}

		p.Mail = nil

		return p, nil

	case 9: //s2c_delete_mail
		p := &pb8.S2CDeleteMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.s2c_delete_mail PrintMsgProto &S2CDeleteMailProto fail")
		}

		p.Id = nil

		return p, nil

	case 10: //s2c_fail_delete_mail
		return toErrCodeMessage(8, 10, data), nil

	case 12: //s2c_keep_mail
		p := &pb8.S2CKeepMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.s2c_keep_mail PrintMsgProto &S2CKeepMailProto fail")
		}

		p.Id = nil

		return p, nil

	case 13: //s2c_fail_keep_mail
		return toErrCodeMessage(8, 13, data), nil

	case 15: //s2c_collect_mail_prize
		p := &pb8.S2CCollectMailPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.s2c_collect_mail_prize PrintMsgProto &S2CCollectMailPrizeProto fail")
		}

		p.Id = nil

		return p, nil

	case 16: //s2c_fail_collect_mail_prize
		return toErrCodeMessage(8, 16, data), nil

	case 21: //s2c_read_mail
		p := &pb8.S2CReadMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.s2c_read_mail PrintMsgProto &S2CReadMailProto fail")
		}

		p.Id = nil

		return p, nil

	case 22: //s2c_fail_read_mail
		return toErrCodeMessage(8, 22, data), nil

	case 23: //s2c_notify_mail_count
		p := &pb8.S2CNotifyMailCountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.s2c_notify_mail_count PrintMsgProto &S2CNotifyMailCountProto fail")
		}

		return p, nil

	case 25: //s2c_read_multi
		p := &pb8.S2CReadMultiProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.s2c_read_multi PrintMsgProto &S2CReadMultiProto fail")
		}

		p.Ids = nil

		p.Prize = nil

		return p, nil

	case 27: //s2c_delete_multi
		p := &pb8.S2CDeleteMultiProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.s2c_delete_multi PrintMsgProto &S2CDeleteMultiProto fail")
		}

		p.Ids = nil

		return p, nil

	case 29: //s2c_get_mail
		p := &pb8.S2CGetMailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mail.s2c_get_mail PrintMsgProto &S2CGetMailProto fail")
		}

		p.Data = nil

		return p, nil

	case 30: //s2c_fail_get_mail
		return toErrCodeMessage(8, 30, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_guild(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_list_guild
		p := &pb9.S2CListGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_list_guild PrintMsgProto &S2CListGuildProto fail")
		}

		return p, nil

	case 3: //s2c_fail_list_guild
		return toErrCodeMessage(9, 3, data), nil

	case 5: //s2c_search_guild
		p := &pb9.S2CSearchGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_search_guild PrintMsgProto &S2CSearchGuildProto fail")
		}

		p.Proto = nil

		return p, nil

	case 6: //s2c_fail_search_guild
		return toErrCodeMessage(9, 6, data), nil

	case 8: //s2c_create_guild
		p := &pb9.S2CCreateGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_create_guild PrintMsgProto &S2CCreateGuildProto fail")
		}

		p.Proto = nil

		return p, nil

	case 9: //s2c_fail_create_guild
		return toErrCodeMessage(9, 9, data), nil

	case 11: //s2c_self_guild
		p := &pb9.S2CSelfGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_self_guild PrintMsgProto &S2CSelfGuildProto fail")
		}

		p.Proto = nil

		return p, nil

	case 12: //s2c_fail_self_guild
		return toErrCodeMessage(9, 12, data), nil

	case 87: //s2c_self_guild_same_version
		return toStringMessage(9, 87), nil

	case 88: //s2c_self_guild_changed
		return toStringMessage(9, 88), nil

	case 14: //s2c_leave_guild
		return toStringMessage(9, 14), nil

	case 15: //s2c_fail_leave_guild
		return toErrCodeMessage(9, 15, data), nil

	case 16: //s2c_leave_guild_for_other
		p := &pb9.S2CLeaveGuildForOtherProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_leave_guild_for_other PrintMsgProto &S2CLeaveGuildForOtherProto fail")
		}

		p.Id = nil

		return p, nil

	case 18: //s2c_kick_other
		p := &pb9.S2CKickOtherProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_kick_other PrintMsgProto &S2CKickOtherProto fail")
		}

		p.Id = nil

		return p, nil

	case 19: //s2c_fail_kick_other
		return toErrCodeMessage(9, 19, data), nil

	case 89: //s2c_self_been_kicked
		return toStringMessage(9, 89), nil

	case 21: //s2c_update_text
		p := &pb9.S2CUpdateTextProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_text PrintMsgProto &S2CUpdateTextProto fail")
		}

		return p, nil

	case 22: //s2c_fail_update_text
		return toErrCodeMessage(9, 22, data), nil

	case 66: //s2c_update_internal_text
		p := &pb9.S2CUpdateInternalTextProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_internal_text PrintMsgProto &S2CUpdateInternalTextProto fail")
		}

		return p, nil

	case 67: //s2c_fail_update_internal_text
		return toErrCodeMessage(9, 67, data), nil

	case 24: //s2c_update_class_names
		p := &pb9.S2CUpdateClassNamesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_class_names PrintMsgProto &S2CUpdateClassNamesProto fail")
		}

		return p, nil

	case 25: //s2c_fail_update_class_names
		return toErrCodeMessage(9, 25, data), nil

	case 123: //s2c_update_class_title
		return toStringMessage(9, 123), nil

	case 124: //s2c_fail_update_class_title
		return toErrCodeMessage(9, 124, data), nil

	case 27: //s2c_update_flag_type
		p := &pb9.S2CUpdateFlagTypeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_flag_type PrintMsgProto &S2CUpdateFlagTypeProto fail")
		}

		return p, nil

	case 28: //s2c_fail_update_flag_type
		return toErrCodeMessage(9, 28, data), nil

	case 30: //s2c_update_member_class_level
		p := &pb9.S2CUpdateMemberClassLevelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_member_class_level PrintMsgProto &S2CUpdateMemberClassLevelProto fail")
		}

		p.Id = nil

		return p, nil

	case 31: //s2c_fail_update_member_class_level
		return toErrCodeMessage(9, 31, data), nil

	case 258: //s2c_update_self_class_level
		p := &pb9.S2CUpdateSelfClassLevelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_self_class_level PrintMsgProto &S2CUpdateSelfClassLevelProto fail")
		}

		return p, nil

	case 81: //s2c_cancel_change_leader
		return toStringMessage(9, 81), nil

	case 82: //s2c_fail_cancel_change_leader
		return toErrCodeMessage(9, 82, data), nil

	case 69: //s2c_update_join_condition
		p := &pb9.S2CUpdateJoinConditionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_join_condition PrintMsgProto &S2CUpdateJoinConditionProto fail")
		}

		return p, nil

	case 70: //s2c_fail_update_join_condition
		return toErrCodeMessage(9, 70, data), nil

	case 72: //s2c_update_guild_name
		p := &pb9.S2CUpdateGuildNameProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_guild_name PrintMsgProto &S2CUpdateGuildNameProto fail")
		}

		return p, nil

	case 73: //s2c_fail_update_guild_name
		return toErrCodeMessage(9, 73, data), nil

	case 74: //s2c_update_guild_name_broadcast
		p := &pb9.S2CUpdateGuildNameBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_guild_name_broadcast PrintMsgProto &S2CUpdateGuildNameBroadcastProto fail")
		}

		return p, nil

	case 76: //s2c_update_guild_label
		p := &pb9.S2CUpdateGuildLabelProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_guild_label PrintMsgProto &S2CUpdateGuildLabelProto fail")
		}

		return p, nil

	case 77: //s2c_fail_update_guild_label
		return toErrCodeMessage(9, 77, data), nil

	case 86: //s2c_update_contribution_coin
		p := &pb9.S2CUpdateContributionCoinProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_contribution_coin PrintMsgProto &S2CUpdateContributionCoinProto fail")
		}

		return p, nil

	case 84: //s2c_donate
		p := &pb9.S2CDonateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_donate PrintMsgProto &S2CDonateProto fail")
		}

		return p, nil

	case 85: //s2c_fail_donate
		return toErrCodeMessage(9, 85, data), nil

	case 91: //s2c_upgrade_level
		return toStringMessage(9, 91), nil

	case 92: //s2c_fail_upgrade_level
		return toErrCodeMessage(9, 92, data), nil

	case 94: //s2c_reduce_upgrade_level_cd
		return toStringMessage(9, 94), nil

	case 95: //s2c_fail_reduce_upgrade_level_cd
		return toErrCodeMessage(9, 95, data), nil

	case 97: //s2c_impeach_leader
		return toStringMessage(9, 97), nil

	case 98: //s2c_fail_impeach_leader
		return toErrCodeMessage(9, 98, data), nil

	case 100: //s2c_impeach_leader_vote
		p := &pb9.S2CImpeachLeaderVoteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_impeach_leader_vote PrintMsgProto &S2CImpeachLeaderVoteProto fail")
		}

		p.Impeach = nil

		return p, nil

	case 101: //s2c_fail_impeach_leader_vote
		return toErrCodeMessage(9, 101, data), nil

	case 103: //s2c_list_guild_by_ids
		p := &pb9.S2CListGuildByIdsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_list_guild_by_ids PrintMsgProto &S2CListGuildByIdsProto fail")
		}

		p.Guilds = nil

		return p, nil

	case 104: //s2c_fail_list_guild_by_ids
		return toErrCodeMessage(9, 104, data), nil

	case 41: //s2c_user_request_join
		p := &pb9.S2CUserRequestJoinProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_user_request_join PrintMsgProto &S2CUserRequestJoinProto fail")
		}

		return p, nil

	case 42: //s2c_fail_user_request_join
		return toErrCodeMessage(9, 42, data), nil

	case 118: //s2c_user_remove_join_request
		p := &pb9.S2CUserRemoveJoinRequestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_user_remove_join_request PrintMsgProto &S2CUserRemoveJoinRequestProto fail")
		}

		return p, nil

	case 119: //s2c_user_clear_join_request
		return toStringMessage(9, 119), nil

	case 44: //s2c_user_cancel_join_request
		p := &pb9.S2CUserCancelJoinRequestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_user_cancel_join_request PrintMsgProto &S2CUserCancelJoinRequestProto fail")
		}

		return p, nil

	case 45: //s2c_fail_user_cancel_join_request
		return toErrCodeMessage(9, 45, data), nil

	case 35: //s2c_add_guild_member
		p := &pb9.S2CAddGuildMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_add_guild_member PrintMsgProto &S2CAddGuildMemberProto fail")
		}

		p.Id = nil

		return p, nil

	case 36: //s2c_user_joined
		p := &pb9.S2CUserJoinedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_user_joined PrintMsgProto &S2CUserJoinedProto fail")
		}

		return p, nil

	case 56: //s2c_guild_reply_join_request
		p := &pb9.S2CGuildReplyJoinRequestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_guild_reply_join_request PrintMsgProto &S2CGuildReplyJoinRequestProto fail")
		}

		p.Id = nil

		return p, nil

	case 57: //s2c_fail_guild_reply_join_request
		return toErrCodeMessage(9, 57, data), nil

	case 110: //s2c_guild_invate_other
		p := &pb9.S2CGuildInvateOtherProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_guild_invate_other PrintMsgProto &S2CGuildInvateOtherProto fail")
		}

		p.Id = nil

		return p, nil

	case 111: //s2c_fail_guild_invate_other
		return toErrCodeMessage(9, 111, data), nil

	case 113: //s2c_guild_cancel_invate_other
		p := &pb9.S2CGuildCancelInvateOtherProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_guild_cancel_invate_other PrintMsgProto &S2CGuildCancelInvateOtherProto fail")
		}

		p.Id = nil

		return p, nil

	case 114: //s2c_fail_guild_cancel_invate_other
		return toErrCodeMessage(9, 114, data), nil

	case 120: //s2c_user_add_been_invate_guild
		p := &pb9.S2CUserAddBeenInvateGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_user_add_been_invate_guild PrintMsgProto &S2CUserAddBeenInvateGuildProto fail")
		}

		return p, nil

	case 121: //s2c_user_remove_been_invate_guild
		p := &pb9.S2CUserRemoveBeenInvateGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_user_remove_been_invate_guild PrintMsgProto &S2CUserRemoveBeenInvateGuildProto fail")
		}

		return p, nil

	case 49: //s2c_user_reply_invate_request
		p := &pb9.S2CUserReplyInvateRequestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_user_reply_invate_request PrintMsgProto &S2CUserReplyInvateRequestProto fail")
		}

		return p, nil

	case 50: //s2c_fail_user_reply_invate_request
		return toErrCodeMessage(9, 50, data), nil

	case 194: //s2c_list_invite_me_guild
		p := &pb9.S2CListInviteMeGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_list_invite_me_guild PrintMsgProto &S2CListInviteMeGuildProto fail")
		}

		return p, nil

	case 195: //s2c_fail_list_invite_me_guild
		return toErrCodeMessage(9, 195, data), nil

	case 126: //s2c_update_friend_guild
		p := &pb9.S2CUpdateFriendGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_friend_guild PrintMsgProto &S2CUpdateFriendGuildProto fail")
		}

		return p, nil

	case 127: //s2c_fail_update_friend_guild
		return toErrCodeMessage(9, 127, data), nil

	case 129: //s2c_update_enemy_guild
		p := &pb9.S2CUpdateEnemyGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_enemy_guild PrintMsgProto &S2CUpdateEnemyGuildProto fail")
		}

		return p, nil

	case 130: //s2c_fail_update_enemy_guild
		return toErrCodeMessage(9, 130, data), nil

	case 132: //s2c_update_guild_prestige
		p := &pb9.S2CUpdateGuildPrestigeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_guild_prestige PrintMsgProto &S2CUpdateGuildPrestigeProto fail")
		}

		return p, nil

	case 133: //s2c_fail_update_guild_prestige
		return toErrCodeMessage(9, 133, data), nil

	case 135: //s2c_place_guild_statue
		return toStringMessage(9, 135), nil

	case 136: //s2c_fail_place_guild_statue
		return toErrCodeMessage(9, 136, data), nil

	case 137: //s2c_guild_statue
		p := &pb9.S2CGuildStatueProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_guild_statue PrintMsgProto &S2CGuildStatueProto fail")
		}

		return p, nil

	case 139: //s2c_take_back_guild_statue
		return toStringMessage(9, 139), nil

	case 140: //s2c_broadcast_take_back_guild_statue
		return toStringMessage(9, 140), nil

	case 141: //s2c_fail_take_back_guild_statue
		return toErrCodeMessage(9, 141, data), nil

	case 144: //s2c_collect_first_join_guild_prize
		return toStringMessage(9, 144), nil

	case 145: //s2c_fail_collect_first_join_guild_prize
		return toErrCodeMessage(9, 145, data), nil

	case 146: //s2c_update_seek_help
		p := &pb9.S2CUpdateSeekHelpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_seek_help PrintMsgProto &S2CUpdateSeekHelpProto fail")
		}

		return p, nil

	case 157: //s2c_update_help_member_times
		p := &pb9.S2CUpdateHelpMemberTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_help_member_times PrintMsgProto &S2CUpdateHelpMemberTimesProto fail")
		}

		return p, nil

	case 148: //s2c_seek_help
		p := &pb9.S2CSeekHelpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_seek_help PrintMsgProto &S2CSeekHelpProto fail")
		}

		return p, nil

	case 149: //s2c_fail_seek_help
		return toErrCodeMessage(9, 149, data), nil

	case 150: //s2c_add_guild_seek_help
		p := &pb9.S2CAddGuildSeekHelpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_add_guild_seek_help PrintMsgProto &S2CAddGuildSeekHelpProto fail")
		}

		p.Data = nil

		return p, nil

	case 152: //s2c_help_guild_member
		p := &pb9.S2CHelpGuildMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_help_guild_member PrintMsgProto &S2CHelpGuildMemberProto fail")
		}

		return p, nil

	case 153: //s2c_fail_help_guild_member
		return toErrCodeMessage(9, 153, data), nil

	case 159: //s2c_help_all_guild_member
		p := &pb9.S2CHelpAllGuildMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_help_all_guild_member PrintMsgProto &S2CHelpAllGuildMemberProto fail")
		}

		return p, nil

	case 160: //s2c_fail_help_all_guild_member
		return toErrCodeMessage(9, 160, data), nil

	case 154: //s2c_add_guild_seek_help_hero_ids
		p := &pb9.S2CAddGuildSeekHelpHeroIdsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_add_guild_seek_help_hero_ids PrintMsgProto &S2CAddGuildSeekHelpHeroIdsProto fail")
		}

		p.HelpHeroId = nil

		return p, nil

	case 155: //s2c_remove_guild_seek_help
		p := &pb9.S2CRemoveGuildSeekHelpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_remove_guild_seek_help PrintMsgProto &S2CRemoveGuildSeekHelpProto fail")
		}

		return p, nil

	case 156: //s2c_list_guild_seek_help
		p := &pb9.S2CListGuildSeekHelpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_list_guild_seek_help PrintMsgProto &S2CListGuildSeekHelpProto fail")
		}

		p.Data = nil

		return p, nil

	case 161: //s2c_list_guild_event_prize
		p := &pb9.S2CListGuildEventPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_list_guild_event_prize PrintMsgProto &S2CListGuildEventPrizeProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 162: //s2c_add_guild_event_prize
		p := &pb9.S2CAddGuildEventPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_add_guild_event_prize PrintMsgProto &S2CAddGuildEventPrizeProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 170: //s2c_remove_guild_event_prize
		p := &pb9.S2CRemoveGuildEventPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_remove_guild_event_prize PrintMsgProto &S2CRemoveGuildEventPrizeProto fail")
		}

		return p, nil

	case 164: //s2c_collect_guild_event_prize
		p := &pb9.S2CCollectGuildEventPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_collect_guild_event_prize PrintMsgProto &S2CCollectGuildEventPrizeProto fail")
		}

		p.Prize = nil

		return p, nil

	case 165: //s2c_fail_collect_guild_event_prize
		return toErrCodeMessage(9, 165, data), nil

	case 166: //s2c_update_full_big_box
		p := &pb9.S2CUpdateFullBigBoxProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_full_big_box PrintMsgProto &S2CUpdateFullBigBoxProto fail")
		}

		return p, nil

	case 168: //s2c_collect_full_big_box
		return toStringMessage(9, 168), nil

	case 169: //s2c_fail_collect_full_big_box
		return toErrCodeMessage(9, 169, data), nil

	case 171: //s2c_update_hero_join_guild_time
		p := &pb9.S2CUpdateHeroJoinGuildTimeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_hero_join_guild_time PrintMsgProto &S2CUpdateHeroJoinGuildTimeProto fail")
		}

		return p, nil

	case 173: //s2c_upgrade_technology
		p := &pb9.S2CUpgradeTechnologyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_upgrade_technology PrintMsgProto &S2CUpgradeTechnologyProto fail")
		}

		return p, nil

	case 174: //s2c_fail_upgrade_technology
		return toErrCodeMessage(9, 174, data), nil

	case 176: //s2c_reduce_technology_cd
		return toStringMessage(9, 176), nil

	case 177: //s2c_fail_reduce_technology_cd
		return toErrCodeMessage(9, 177, data), nil

	case 179: //s2c_list_guild_logs
		p := &pb9.S2CListGuildLogsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_list_guild_logs PrintMsgProto &S2CListGuildLogsProto fail")
		}

		p.Data = nil

		return p, nil

	case 180: //s2c_add_guild_log
		p := &pb9.S2CAddGuildLogProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_add_guild_log PrintMsgProto &S2CAddGuildLogProto fail")
		}

		p.Data = nil

		return p, nil

	case 182: //s2c_request_recommend_guild
		p := &pb9.S2CRequestRecommendGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_request_recommend_guild PrintMsgProto &S2CRequestRecommendGuildProto fail")
		}

		return p, nil

	case 183: //s2c_push_tech_helpable
		p := &pb9.S2CPushTechHelpableProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_push_tech_helpable PrintMsgProto &S2CPushTechHelpableProto fail")
		}

		return p, nil

	case 185: //s2c_help_tech
		p := &pb9.S2CHelpTechProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_help_tech PrintMsgProto &S2CHelpTechProto fail")
		}

		return p, nil

	case 186: //s2c_fail_help_tech
		return toErrCodeMessage(9, 186, data), nil

	case 188: //s2c_recommend_invite_heros
		p := &pb9.S2CRecommendInviteHerosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_recommend_invite_heros PrintMsgProto &S2CRecommendInviteHerosProto fail")
		}

		p.Heros = nil

		return p, nil

	case 189: //s2c_fail_recommend_invite_heros
		return toErrCodeMessage(9, 189, data), nil

	case 191: //s2c_search_no_guild_heros
		p := &pb9.S2CSearchNoGuildHerosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_search_no_guild_heros PrintMsgProto &S2CSearchNoGuildHerosProto fail")
		}

		p.Heros = nil

		return p, nil

	case 192: //s2c_fail_search_no_guild_heros
		return toErrCodeMessage(9, 192, data), nil

	case 200: //s2c_view_mc_war_record
		p := &pb9.S2CViewMcWarRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_view_mc_war_record PrintMsgProto &S2CViewMcWarRecordProto fail")
		}

		return p, nil

	case 201: //s2c_fail_view_mc_war_record
		return toErrCodeMessage(9, 201, data), nil

	case 197: //s2c_update_guild_mark
		p := &pb9.S2CUpdateGuildMarkProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_guild_mark PrintMsgProto &S2CUpdateGuildMarkProto fail")
		}

		return p, nil

	case 198: //s2c_fail_update_guild_mark
		return toErrCodeMessage(9, 198, data), nil

	case 203: //s2c_view_yinliang_record
		p := &pb9.S2CViewYinliangRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_view_yinliang_record PrintMsgProto &S2CViewYinliangRecordProto fail")
		}

		return p, nil

	case 204: //s2c_fail_view_yinliang_record
		return toErrCodeMessage(9, 204, data), nil

	case 206: //s2c_send_yinliang_to_other_guild
		p := &pb9.S2CSendYinliangToOtherGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_send_yinliang_to_other_guild PrintMsgProto &S2CSendYinliangToOtherGuildProto fail")
		}

		return p, nil

	case 207: //s2c_fail_send_yinliang_to_other_guild
		return toErrCodeMessage(9, 207, data), nil

	case 209: //s2c_send_yinliang_to_member
		p := &pb9.S2CSendYinliangToMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_send_yinliang_to_member PrintMsgProto &S2CSendYinliangToMemberProto fail")
		}

		p.MemId = nil

		return p, nil

	case 210: //s2c_fail_send_yinliang_to_member
		return toErrCodeMessage(9, 210, data), nil

	case 212: //s2c_pay_salary
		p := &pb9.S2CPaySalaryProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_pay_salary PrintMsgProto &S2CPaySalaryProto fail")
		}

		return p, nil

	case 213: //s2c_fail_pay_salary
		return toErrCodeMessage(9, 213, data), nil

	case 215: //s2c_set_salary
		p := &pb9.S2CSetSalaryProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_set_salary PrintMsgProto &S2CSetSalaryProto fail")
		}

		p.MemId = nil

		return p, nil

	case 216: //s2c_fail_set_salary
		return toErrCodeMessage(9, 216, data), nil

	case 217: //s2c_update_hero_guild
		p := &pb9.S2CUpdateHeroGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_hero_guild PrintMsgProto &S2CUpdateHeroGuildProto fail")
		}

		return p, nil

	case 219: //s2c_view_send_yinliang_to_guild
		p := &pb9.S2CViewSendYinliangToGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_view_send_yinliang_to_guild PrintMsgProto &S2CViewSendYinliangToGuildProto fail")
		}

		return p, nil

	case 220: //s2c_fail_view_send_yinliang_to_guild
		return toErrCodeMessage(9, 220, data), nil

	case 221: //s2c_update_hufu
		p := &pb9.S2CUpdateHufuProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_update_hufu PrintMsgProto &S2CUpdateHufuProto fail")
		}

		return p, nil

	case 229: //s2c_convene
		p := &pb9.S2CConveneProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_convene PrintMsgProto &S2CConveneProto fail")
		}

		p.Target = nil

		return p, nil

	case 230: //s2c_fail_convene
		return toErrCodeMessage(9, 230, data), nil

	case 232: //s2c_collect_daily_guild_rank_prize
		p := &pb9.S2CCollectDailyGuildRankPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_collect_daily_guild_rank_prize PrintMsgProto &S2CCollectDailyGuildRankPrizeProto fail")
		}

		p.Prize = nil

		return p, nil

	case 233: //s2c_fail_collect_daily_guild_rank_prize
		return toErrCodeMessage(9, 233, data), nil

	case 235: //s2c_view_daily_guild_rank
		p := &pb9.S2CViewDailyGuildRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_view_daily_guild_rank PrintMsgProto &S2CViewDailyGuildRankProto fail")
		}

		return p, nil

	case 236: //s2c_fail_view_daily_guild_rank
		return toErrCodeMessage(9, 236, data), nil

	case 238: //s2c_view_last_guild_rank
		p := &pb9.S2CViewLastGuildRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_view_last_guild_rank PrintMsgProto &S2CViewLastGuildRankProto fail")
		}

		return p, nil

	case 241: //s2c_add_recommend_mc_build
		p := &pb9.S2CAddRecommendMcBuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_add_recommend_mc_build PrintMsgProto &S2CAddRecommendMcBuildProto fail")
		}

		return p, nil

	case 242: //s2c_fail_add_recommend_mc_build
		return toErrCodeMessage(9, 242, data), nil

	case 244: //s2c_view_task_progress
		p := &pb9.S2CViewTaskProgressProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_view_task_progress PrintMsgProto &S2CViewTaskProgressProto fail")
		}

		return p, nil

	case 245: //s2c_fail_view_task_progress
		return toErrCodeMessage(9, 245, data), nil

	case 246: //s2c_notice_task_stage_update
		p := &pb9.S2CNoticeTaskStageUpdateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_notice_task_stage_update PrintMsgProto &S2CNoticeTaskStageUpdateProto fail")
		}

		return p, nil

	case 248: //s2c_collect_task_prize
		p := &pb9.S2CCollectTaskPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_collect_task_prize PrintMsgProto &S2CCollectTaskPrizeProto fail")
		}

		return p, nil

	case 249: //s2c_fail_collect_task_prize
		return toErrCodeMessage(9, 249, data), nil

	case 251: //s2c_guild_change_country
		p := &pb9.S2CGuildChangeCountryProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_guild_change_country PrintMsgProto &S2CGuildChangeCountryProto fail")
		}

		return p, nil

	case 252: //s2c_fail_guild_change_country
		return toErrCodeMessage(9, 252, data), nil

	case 254: //s2c_cancel_guild_change_country
		return toStringMessage(9, 254), nil

	case 255: //s2c_fail_cancel_guild_change_country
		return toErrCodeMessage(9, 255, data), nil

	case 257: //s2c_show_workshop_not_exist
		p := &pb9.S2CShowWorkshopNotExistProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "guild.s2c_show_workshop_not_exist PrintMsgProto &S2CShowWorkshopNotExistProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_stress(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_robot_ping
		p := &pb10.S2CRobotPingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "stress.s2c_robot_ping PrintMsgProto &S2CRobotPingProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_depot(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 1: //s2c_update_goods
		p := &pb11.S2CUpdateGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_update_goods PrintMsgProto &S2CUpdateGoodsProto fail")
		}

		return p, nil

	case 5: //s2c_update_multi_goods
		p := &pb11.S2CUpdateMultiGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_update_multi_goods PrintMsgProto &S2CUpdateMultiGoodsProto fail")
		}

		return p, nil

	case 3: //s2c_use_goods
		p := &pb11.S2CUseGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_use_goods PrintMsgProto &S2CUseGoodsProto fail")
		}

		return p, nil

	case 4: //s2c_fail_use_goods
		return toErrCodeMessage(11, 4, data), nil

	case 7: //s2c_use_cdr_goods
		p := &pb11.S2CUseCdrGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_use_cdr_goods PrintMsgProto &S2CUseCdrGoodsProto fail")
		}

		return p, nil

	case 8: //s2c_fail_use_cdr_goods
		return toErrCodeMessage(11, 8, data), nil

	case 10: //s2c_goods_combine
		return toStringMessage(11, 10), nil

	case 11: //s2c_fail_goods_combine
		return toErrCodeMessage(11, 11, data), nil

	case 19: //s2c_goods_parts_combine
		p := &pb11.S2CGoodsPartsCombineProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_goods_parts_combine PrintMsgProto &S2CGoodsPartsCombineProto fail")
		}

		p.Prize = nil

		return p, nil

	case 20: //s2c_fail_goods_parts_combine
		return toErrCodeMessage(11, 20, data), nil

	case 16: //s2c_goods_expired
		p := &pb11.S2CGoodsExpiredProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_goods_expired PrintMsgProto &S2CGoodsExpiredProto fail")
		}

		return p, nil

	case 17: //s2c_goods_expire_time_remove
		p := &pb11.S2CGoodsExpireTimeRemoveProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_goods_expire_time_remove PrintMsgProto &S2CGoodsExpireTimeRemoveProto fail")
		}

		return p, nil

	case 21: //s2c_update_baowu
		p := &pb11.S2CUpdateBaowuProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_update_baowu PrintMsgProto &S2CUpdateBaowuProto fail")
		}

		return p, nil

	case 22: //s2c_update_multi_baowu
		p := &pb11.S2CUpdateMultiBaowuProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_update_multi_baowu PrintMsgProto &S2CUpdateMultiBaowuProto fail")
		}

		return p, nil

	case 24: //s2c_unlock_baowu
		p := &pb11.S2CUnlockBaowuProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_unlock_baowu PrintMsgProto &S2CUnlockBaowuProto fail")
		}

		return p, nil

	case 25: //s2c_fail_unlock_baowu
		return toErrCodeMessage(11, 25, data), nil

	case 27: //s2c_collect_baowu
		p := &pb11.S2CCollectBaowuProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_collect_baowu PrintMsgProto &S2CCollectBaowuProto fail")
		}

		return p, nil

	case 28: //s2c_fail_collect_baowu
		return toErrCodeMessage(11, 28, data), nil

	case 29: //s2c_add_baowu_log
		p := &pb11.S2CAddBaowuLogProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_add_baowu_log PrintMsgProto &S2CAddBaowuLogProto fail")
		}

		p.Data = nil

		return p, nil

	case 31: //s2c_list_baowu_log
		p := &pb11.S2CListBaowuLogProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_list_baowu_log PrintMsgProto &S2CListBaowuLogProto fail")
		}

		p.Datas = nil

		return p, nil

	case 36: //s2c_decompose_baowu
		p := &pb11.S2CDecomposeBaowuProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "depot.s2c_decompose_baowu PrintMsgProto &S2CDecomposeBaowuProto fail")
		}

		return p, nil

	case 37: //s2c_fail_decompose_baowu
		return toErrCodeMessage(11, 37, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_equipment(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 41: //s2c_view_chat_equip
		p := &pb12.S2CViewChatEquipProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_view_chat_equip PrintMsgProto &S2CViewChatEquipProto fail")
		}

		p.Data = nil

		return p, nil

	case 42: //s2c_fail_view_chat_equip
		return toErrCodeMessage(12, 42, data), nil

	case 18: //s2c_add_equipment
		p := &pb12.S2CAddEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_add_equipment PrintMsgProto &S2CAddEquipmentProto fail")
		}

		p.Data = nil

		return p, nil

	case 34: //s2c_add_equipment_with_expire_time
		p := &pb12.S2CAddEquipmentWithExpireTimeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_add_equipment_with_expire_time PrintMsgProto &S2CAddEquipmentWithExpireTimeProto fail")
		}

		p.Data = nil

		return p, nil

	case 2: //s2c_wear_equipment
		p := &pb12.S2CWearEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_wear_equipment PrintMsgProto &S2CWearEquipmentProto fail")
		}

		return p, nil

	case 3: //s2c_fail_wear_equipment
		return toErrCodeMessage(12, 3, data), nil

	case 5: //s2c_upgrade_equipment
		p := &pb12.S2CUpgradeEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_upgrade_equipment PrintMsgProto &S2CUpgradeEquipmentProto fail")
		}

		return p, nil

	case 6: //s2c_fail_upgrade_equipment
		return toErrCodeMessage(12, 6, data), nil

	case 20: //s2c_upgrade_equipment_all
		p := &pb12.S2CUpgradeEquipmentAllProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_upgrade_equipment_all PrintMsgProto &S2CUpgradeEquipmentAllProto fail")
		}

		return p, nil

	case 21: //s2c_fail_upgrade_equipment_all
		return toErrCodeMessage(12, 21, data), nil

	case 8: //s2c_refined_equipment
		p := &pb12.S2CRefinedEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_refined_equipment PrintMsgProto &S2CRefinedEquipmentProto fail")
		}

		return p, nil

	case 9: //s2c_fail_refined_equipment
		return toErrCodeMessage(12, 9, data), nil

	case 25: //s2c_update_equipment
		p := &pb12.S2CUpdateEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_update_equipment PrintMsgProto &S2CUpdateEquipmentProto fail")
		}

		p.Data = nil

		return p, nil

	case 26: //s2c_update_multi_equipment
		p := &pb12.S2CUpdateMultiEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_update_multi_equipment PrintMsgProto &S2CUpdateMultiEquipmentProto fail")
		}

		p.Data = nil

		return p, nil

	case 11: //s2c_smelt_equipment
		p := &pb12.S2CSmeltEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_smelt_equipment PrintMsgProto &S2CSmeltEquipmentProto fail")
		}

		return p, nil

	case 12: //s2c_fail_smelt_equipment
		return toErrCodeMessage(12, 12, data), nil

	case 14: //s2c_rebuild_equipment
		p := &pb12.S2CRebuildEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_rebuild_equipment PrintMsgProto &S2CRebuildEquipmentProto fail")
		}

		return p, nil

	case 15: //s2c_fail_rebuild_equipment
		return toErrCodeMessage(12, 15, data), nil

	case 33: //s2c_open_equip_combine
		p := &pb12.S2COpenEquipCombineProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_open_equip_combine PrintMsgProto &S2COpenEquipCombineProto fail")
		}

		return p, nil

	case 35: //s2c_rebuild_upgrade_equipment
		p := &pb12.S2CRebuildUpgradeEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_rebuild_upgrade_equipment PrintMsgProto &S2CRebuildUpgradeEquipmentProto fail")
		}

		return p, nil

	case 36: //s2c_rebuild_refine_equipment
		p := &pb12.S2CRebuildRefineEquipmentProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_rebuild_refine_equipment PrintMsgProto &S2CRebuildRefineEquipmentProto fail")
		}

		return p, nil

	case 44: //s2c_one_key_take_off
		p := &pb12.S2COneKeyTakeOffProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "equipment.s2c_one_key_take_off PrintMsgProto &S2COneKeyTakeOffProto fail")
		}

		return p, nil

	case 45: //s2c_fail_one_key_take_off
		return toErrCodeMessage(12, 45, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_chat(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_world_chat
		return toStringMessage(13, 2), nil

	case 3: //s2c_world_other_chat
		p := &pb13.S2CWorldOtherChatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_world_other_chat PrintMsgProto &S2CWorldOtherChatProto fail")
		}

		p.Id = nil

		return p, nil

	case 5: //s2c_guild_chat
		return toStringMessage(13, 5), nil

	case 7: //s2c_fail_guild_chat
		return toErrCodeMessage(13, 7, data), nil

	case 6: //s2c_guild_other_chat
		p := &pb13.S2CGuildOtherChatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_guild_other_chat PrintMsgProto &S2CGuildOtherChatProto fail")
		}

		p.Id = nil

		return p, nil

	case 9: //s2c_self_chat_window
		p := &pb13.S2CSelfChatWindowProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_self_chat_window PrintMsgProto &S2CSelfChatWindowProto fail")
		}

		p.Sender = nil

		return p, nil

	case 22: //s2c_create_self_chat_window
		p := &pb13.S2CCreateSelfChatWindowProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_create_self_chat_window PrintMsgProto &S2CCreateSelfChatWindowProto fail")
		}

		p.Target = nil

		return p, nil

	case 11: //s2c_remove_chat_window
		p := &pb13.S2CRemoveChatWindowProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_remove_chat_window PrintMsgProto &S2CRemoveChatWindowProto fail")
		}

		p.ChatTarget = nil

		return p, nil

	case 13: //s2c_list_history_chat
		p := &pb13.S2CListHistoryChatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_list_history_chat PrintMsgProto &S2CListHistoryChatProto fail")
		}

		p.ChatMsg = nil

		return p, nil

	case 15: //s2c_send_chat
		p := &pb13.S2CSendChatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_send_chat PrintMsgProto &S2CSendChatProto fail")
		}

		p.ChatId = nil

		p.Receiver = nil

		return p, nil

	case 16: //s2c_fail_send_chat
		return toErrCodeMessage(13, 16, data), nil

	case 17: //s2c_other_send_chat
		p := &pb13.S2COtherSendChatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_other_send_chat PrintMsgProto &S2COtherSendChatProto fail")
		}

		p.ChatMsg = nil

		return p, nil

	case 19: //s2c_read_chat_msg
		p := &pb13.S2CReadChatMsgProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_read_chat_msg PrintMsgProto &S2CReadChatMsgProto fail")
		}

		p.ChatTarget = nil

		return p, nil

	case 24: //s2c_offline_chat
		return toStringMessage(13, 24), nil

	case 26: //s2c_get_hero_chat_info
		p := &pb13.S2CGetHeroChatInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_get_hero_chat_info PrintMsgProto &S2CGetHeroChatInfoProto fail")
		}

		p.Id = nil

		return p, nil

	case 27: //s2c_ban_chat
		p := &pb13.S2CBanChatProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "chat.s2c_ban_chat PrintMsgProto &S2CBanChatProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_tower(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_challenge
		p := &pb14.S2CChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tower.s2c_challenge PrintMsgProto &S2CChallengeProto fail")
		}

		p.Share = nil

		p.FirstPassPrize = nil

		p.Prize = nil

		return p, nil

	case 3: //s2c_failure_challenge
		p := &pb14.S2CFailureChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tower.s2c_failure_challenge PrintMsgProto &S2CFailureChallengeProto fail")
		}

		p.Share = nil

		return p, nil

	case 4: //s2c_fail_challenge
		return toErrCodeMessage(14, 4, data), nil

	case 6: //s2c_auto_challenge
		p := &pb14.S2CAutoChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tower.s2c_auto_challenge PrintMsgProto &S2CAutoChallengeProto fail")
		}

		p.Prize = nil

		return p, nil

	case 7: //s2c_fail_auto_challenge
		return toErrCodeMessage(14, 7, data), nil

	case 9: //s2c_collect_box
		p := &pb14.S2CCollectBoxProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tower.s2c_collect_box PrintMsgProto &S2CCollectBoxProto fail")
		}

		return p, nil

	case 10: //s2c_fail_collect_box
		return toErrCodeMessage(14, 10, data), nil

	case 12: //s2c_list_pass_replay
		p := &pb14.S2CListPassReplayProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tower.s2c_list_pass_replay PrintMsgProto &S2CListPassReplayProto fail")
		}

		p.Data = nil

		return p, nil

	case 13: //s2c_fail_list_pass_replay
		return toErrCodeMessage(14, 13, data), nil

	case 14: //s2c_update_current_floor
		p := &pb14.S2CUpdateCurrentFloorProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tower.s2c_update_current_floor PrintMsgProto &S2CUpdateCurrentFloorProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_task(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 1: //s2c_update_task_progress
		p := &pb15.S2CUpdateTaskProgressProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_update_task_progress PrintMsgProto &S2CUpdateTaskProgressProto fail")
		}

		return p, nil

	case 3: //s2c_collect_task_prize
		p := &pb15.S2CCollectTaskPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_collect_task_prize PrintMsgProto &S2CCollectTaskPrizeProto fail")
		}

		return p, nil

	case 4: //s2c_fail_collect_task_prize
		return toErrCodeMessage(15, 4, data), nil

	case 5: //s2c_new_task
		p := &pb15.S2CNewTaskProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_new_task PrintMsgProto &S2CNewTaskProto fail")
		}

		return p, nil

	case 7: //s2c_collect_task_box_prize
		p := &pb15.S2CCollectTaskBoxPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_collect_task_box_prize PrintMsgProto &S2CCollectTaskBoxPrizeProto fail")
		}

		return p, nil

	case 8: //s2c_fail_collect_task_box_prize
		return toErrCodeMessage(15, 8, data), nil

	case 10: //s2c_collect_ba_ye_stage_prize
		p := &pb15.S2CCollectBaYeStagePrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_collect_ba_ye_stage_prize PrintMsgProto &S2CCollectBaYeStagePrizeProto fail")
		}

		p.Stage = nil

		return p, nil

	case 11: //s2c_fail_collect_ba_ye_stage_prize
		return toErrCodeMessage(15, 11, data), nil

	case 13: //s2c_collect_active_degree_prize
		p := &pb15.S2CCollectActiveDegreePrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_collect_active_degree_prize PrintMsgProto &S2CCollectActiveDegreePrizeProto fail")
		}

		return p, nil

	case 14: //s2c_fail_collect_active_degree_prize
		return toErrCodeMessage(15, 14, data), nil

	case 17: //s2c_collect_achieve_star_prize
		p := &pb15.S2CCollectAchieveStarPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_collect_achieve_star_prize PrintMsgProto &S2CCollectAchieveStarPrizeProto fail")
		}

		return p, nil

	case 18: //s2c_fail_collect_achieve_star_prize
		return toErrCodeMessage(15, 18, data), nil

	case 19: //s2c_achieve_reach
		p := &pb15.S2CAchieveReachProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_achieve_reach PrintMsgProto &S2CAchieveReachProto fail")
		}

		return p, nil

	case 21: //s2c_change_select_show_achieve
		p := &pb15.S2CChangeSelectShowAchieveProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_change_select_show_achieve PrintMsgProto &S2CChangeSelectShowAchieveProto fail")
		}

		return p, nil

	case 22: //s2c_fail_change_select_show_achieve
		return toErrCodeMessage(15, 22, data), nil

	case 15: //s2c_remove_task
		p := &pb15.S2CRemoveTaskProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_remove_task PrintMsgProto &S2CRemoveTaskProto fail")
		}

		return p, nil

	case 24: //s2c_collect_bwzl_prize
		p := &pb15.S2CCollectBwzlPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_collect_bwzl_prize PrintMsgProto &S2CCollectBwzlPrizeProto fail")
		}

		return p, nil

	case 25: //s2c_fail_collect_bwzl_prize
		return toErrCodeMessage(15, 25, data), nil

	case 27: //s2c_view_other_achieve_task_list
		p := &pb15.S2CViewOtherAchieveTaskListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_view_other_achieve_task_list PrintMsgProto &S2CViewOtherAchieveTaskListProto fail")
		}

		p.Id = nil

		return p, nil

	case 28: //s2c_fail_view_other_achieve_task_list
		return toErrCodeMessage(15, 28, data), nil

	case 30: //s2c_get_troop_title_fight_amount
		p := &pb15.S2CGetTroopTitleFightAmountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_get_troop_title_fight_amount PrintMsgProto &S2CGetTroopTitleFightAmountProto fail")
		}

		return p, nil

	case 32: //s2c_upgrade_title
		p := &pb15.S2CUpgradeTitleProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_upgrade_title PrintMsgProto &S2CUpgradeTitleProto fail")
		}

		return p, nil

	case 33: //s2c_fail_upgrade_title
		return toErrCodeMessage(15, 33, data), nil

	case 36: //s2c_complete_bool_task
		p := &pb15.S2CCompleteBoolTaskProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "task.s2c_complete_bool_task PrintMsgProto &S2CCompleteBoolTaskProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_fishing(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_fishing
		p := &pb16.S2CFishingProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "fishing.s2c_fishing PrintMsgProto &S2CFishingProto fail")
		}

		p.FishingResult = nil

		return p, nil

	case 3: //s2c_fail_fishing
		return toErrCodeMessage(16, 3, data), nil

	case 5: //s2c_fishing_broadcast
		p := &pb16.S2CFishingBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "fishing.s2c_fishing_broadcast PrintMsgProto &S2CFishingBroadcastProto fail")
		}

		p.Id = nil

		p.Prize = nil

		return p, nil

	case 7: //s2c_update_fish_point
		p := &pb16.S2CUpdateFishPointProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "fishing.s2c_update_fish_point PrintMsgProto &S2CUpdateFishPointProto fail")
		}

		return p, nil

	case 9: //s2c_fish_point_exchange
		p := &pb16.S2CFishPointExchangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "fishing.s2c_fish_point_exchange PrintMsgProto &S2CFishPointExchangeProto fail")
		}

		p.Prize = nil

		return p, nil

	case 10: //s2c_fail_fish_point_exchange
		return toErrCodeMessage(16, 10, data), nil

	case 12: //s2c_set_fishing_captain
		p := &pb16.S2CSetFishingCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "fishing.s2c_set_fishing_captain PrintMsgProto &S2CSetFishingCaptainProto fail")
		}

		return p, nil

	case 13: //s2c_fail_set_fishing_captain
		return toErrCodeMessage(16, 13, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_gem(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 4: //s2c_use_gem
		p := &pb19.S2CUseGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.s2c_use_gem PrintMsgProto &S2CUseGemProto fail")
		}

		return p, nil

	case 5: //s2c_fail_use_gem
		return toErrCodeMessage(19, 5, data), nil

	case 22: //s2c_inlay_gem
		p := &pb19.S2CInlayGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.s2c_inlay_gem PrintMsgProto &S2CInlayGemProto fail")
		}

		return p, nil

	case 23: //s2c_fail_inlay_gem
		return toErrCodeMessage(19, 23, data), nil

	case 7: //s2c_combine_gem
		p := &pb19.S2CCombineGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.s2c_combine_gem PrintMsgProto &S2CCombineGemProto fail")
		}

		return p, nil

	case 8: //s2c_fail_combine_gem
		return toErrCodeMessage(19, 8, data), nil

	case 10: //s2c_one_key_use_gem
		p := &pb19.S2COneKeyUseGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.s2c_one_key_use_gem PrintMsgProto &S2COneKeyUseGemProto fail")
		}

		return p, nil

	case 14: //s2c_fail_one_key_use_gem
		return toErrCodeMessage(19, 14, data), nil

	case 12: //s2c_one_key_combine_gem
		p := &pb19.S2COneKeyCombineGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.s2c_one_key_combine_gem PrintMsgProto &S2COneKeyCombineGemProto fail")
		}

		return p, nil

	case 13: //s2c_fail_one_key_combine_gem
		return toErrCodeMessage(19, 13, data), nil

	case 16: //s2c_request_one_key_combine_cost
		p := &pb19.S2CRequestOneKeyCombineCostProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.s2c_request_one_key_combine_cost PrintMsgProto &S2CRequestOneKeyCombineCostProto fail")
		}

		return p, nil

	case 17: //s2c_fail_request_one_key_combine_cost
		return toErrCodeMessage(19, 17, data), nil

	case 19: //s2c_one_key_combine_depot_gem
		p := &pb19.S2COneKeyCombineDepotGemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "gem.s2c_one_key_combine_depot_gem PrintMsgProto &S2COneKeyCombineDepotGemProto fail")
		}

		return p, nil

	case 20: //s2c_fail_one_key_combine_depot_gem
		return toErrCodeMessage(19, 20, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_shop(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 1: //s2c_update_daily_shop_goods
		p := &pb20.S2CUpdateDailyShopGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "shop.s2c_update_daily_shop_goods PrintMsgProto &S2CUpdateDailyShopGoodsProto fail")
		}

		return p, nil

	case 3: //s2c_buy_goods
		p := &pb20.S2CBuyGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "shop.s2c_buy_goods PrintMsgProto &S2CBuyGoodsProto fail")
		}

		p.Prize = nil

		return p, nil

	case 4: //s2c_fail_buy_goods
		return toErrCodeMessage(20, 4, data), nil

	case 5: //s2c_multi_crit_broadcast
		p := &pb20.S2CMultiCritBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "shop.s2c_multi_crit_broadcast PrintMsgProto &S2CMultiCritBroadcastProto fail")
		}

		p.Prize = nil

		return p, nil

	case 8: //s2c_push_black_market_goods
		p := &pb20.S2CPushBlackMarketGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "shop.s2c_push_black_market_goods PrintMsgProto &S2CPushBlackMarketGoodsProto fail")
		}

		return p, nil

	case 10: //s2c_buy_black_market_goods
		p := &pb20.S2CBuyBlackMarketGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "shop.s2c_buy_black_market_goods PrintMsgProto &S2CBuyBlackMarketGoodsProto fail")
		}

		return p, nil

	case 11: //s2c_fail_buy_black_market_goods
		return toErrCodeMessage(20, 11, data), nil

	case 13: //s2c_refresh_black_market_goods
		return toStringMessage(20, 13), nil

	case 14: //s2c_fail_refresh_black_market_goods
		return toErrCodeMessage(20, 14, data), nil

	case 15: //s2c_update_vip_shop_goods
		p := &pb20.S2CUpdateVipShopGoodsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "shop.s2c_update_vip_shop_goods PrintMsgProto &S2CUpdateVipShopGoodsProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_client_config(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_config
		p := &pb21.S2CConfigProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "client_config.s2c_config PrintMsgProto &S2CConfigProto fail")
		}

		p.Data = nil

		return p, nil

	case 3: //s2c_fail_config
		return toErrCodeMessage(21, 3, data), nil

	case 6: //s2c_set_client_key
		p := &pb21.S2CSetClientKeyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "client_config.s2c_set_client_key PrintMsgProto &S2CSetClientKeyProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_secret_tower(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 1: //s2c_unlock_secret_tower
		p := &pb22.S2CUnlockSecretTowerProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_unlock_secret_tower PrintMsgProto &S2CUnlockSecretTowerProto fail")
		}

		return p, nil

	case 3: //s2c_request_team_count
		p := &pb22.S2CRequestTeamCountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_request_team_count PrintMsgProto &S2CRequestTeamCountProto fail")
		}

		return p, nil

	case 4: //s2c_fail_request_team_count
		return toErrCodeMessage(22, 4, data), nil

	case 6: //s2c_request_team_list
		p := &pb22.S2CRequestTeamListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_request_team_list PrintMsgProto &S2CRequestTeamListProto fail")
		}

		p.TeamList = nil

		return p, nil

	case 7: //s2c_fail_request_team_list
		return toErrCodeMessage(22, 7, data), nil

	case 9: //s2c_create_team
		p := &pb22.S2CCreateTeamProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_create_team PrintMsgProto &S2CCreateTeamProto fail")
		}

		p.TeamDetail = nil

		return p, nil

	case 10: //s2c_fail_create_team
		return toErrCodeMessage(22, 10, data), nil

	case 12: //s2c_join_team
		p := &pb22.S2CJoinTeamProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_join_team PrintMsgProto &S2CJoinTeamProto fail")
		}

		p.TeamDetail = nil

		return p, nil

	case 13: //s2c_other_join_join_team
		p := &pb22.S2COtherJoinJoinTeamProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_other_join_join_team PrintMsgProto &S2COtherJoinJoinTeamProto fail")
		}

		p.Member = nil

		return p, nil

	case 14: //s2c_fail_join_team
		return toErrCodeMessage(22, 14, data), nil

	case 16: //s2c_leave_team
		return toStringMessage(22, 16), nil

	case 17: //s2c_other_leave_leave_team
		p := &pb22.S2COtherLeaveLeaveTeamProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_other_leave_leave_team PrintMsgProto &S2COtherLeaveLeaveTeamProto fail")
		}

		p.Id = nil

		p.NewTeamLeaderId = nil

		return p, nil

	case 18: //s2c_fail_leave_team
		return toErrCodeMessage(22, 18, data), nil

	case 20: //s2c_kick_member
		p := &pb22.S2CKickMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_kick_member PrintMsgProto &S2CKickMemberProto fail")
		}

		p.Id = nil

		return p, nil

	case 21: //s2c_you_been_kicked_kick_member
		return toStringMessage(22, 21), nil

	case 22: //s2c_other_been_kick_kick_member
		p := &pb22.S2COtherBeenKickKickMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_other_been_kick_kick_member PrintMsgProto &S2COtherBeenKickKickMemberProto fail")
		}

		p.Id = nil

		return p, nil

	case 23: //s2c_fail_kick_member
		return toErrCodeMessage(22, 23, data), nil

	case 25: //s2c_broadcsat_move_member
		p := &pb22.S2CBroadcsatMoveMemberProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_broadcsat_move_member PrintMsgProto &S2CBroadcsatMoveMemberProto fail")
		}

		p.Id = nil

		return p, nil

	case 26: //s2c_fail_move_member
		return toErrCodeMessage(22, 26, data), nil

	case 68: //s2c_update_member_pos
		p := &pb22.S2CUpdateMemberPosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_update_member_pos PrintMsgProto &S2CUpdateMemberPosProto fail")
		}

		p.Id = nil

		return p, nil

	case 69: //s2c_fail_update_member_pos
		return toErrCodeMessage(22, 69, data), nil

	case 28: //s2c_change_mode
		p := &pb22.S2CChangeModeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_change_mode PrintMsgProto &S2CChangeModeProto fail")
		}

		return p, nil

	case 29: //s2c_other_changed_change_mode
		p := &pb22.S2COtherChangedChangeModeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_other_changed_change_mode PrintMsgProto &S2COtherChangedChangeModeProto fail")
		}

		p.Id = nil

		return p, nil

	case 30: //s2c_fail_change_mode
		return toErrCodeMessage(22, 30, data), nil

	case 31: //s2c_help_times_change
		p := &pb22.S2CHelpTimesChangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_help_times_change PrintMsgProto &S2CHelpTimesChangeProto fail")
		}

		return p, nil

	case 32: //s2c_times_change
		p := &pb22.S2CTimesChangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_times_change PrintMsgProto &S2CTimesChangeProto fail")
		}

		return p, nil

	case 34: //s2c_invite
		p := &pb22.S2CInviteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_invite PrintMsgProto &S2CInviteProto fail")
		}

		p.Id = nil

		return p, nil

	case 51: //s2c_fail_target_not_found_invite
		p := &pb22.S2CFailTargetNotFoundInviteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_target_not_found_invite PrintMsgProto &S2CFailTargetNotFoundInviteProto fail")
		}

		p.Id = nil

		return p, nil

	case 52: //s2c_fail_target_not_in_my_guild_invite
		p := &pb22.S2CFailTargetNotInMyGuildInviteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_target_not_in_my_guild_invite PrintMsgProto &S2CFailTargetNotInMyGuildInviteProto fail")
		}

		p.Id = nil

		return p, nil

	case 53: //s2c_fail_target_not_open_invite
		p := &pb22.S2CFailTargetNotOpenInviteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_target_not_open_invite PrintMsgProto &S2CFailTargetNotOpenInviteProto fail")
		}

		p.Id = nil

		return p, nil

	case 54: //s2c_fail_target_not_online_invite
		p := &pb22.S2CFailTargetNotOnlineInviteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_target_not_online_invite PrintMsgProto &S2CFailTargetNotOnlineInviteProto fail")
		}

		p.Id = nil

		return p, nil

	case 55: //s2c_fail_target_in_your_team_invite
		p := &pb22.S2CFailTargetInYourTeamInviteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_target_in_your_team_invite PrintMsgProto &S2CFailTargetInYourTeamInviteProto fail")
		}

		p.Id = nil

		return p, nil

	case 56: //s2c_fail_target_no_times_invite
		p := &pb22.S2CFailTargetNoTimesInviteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_target_no_times_invite PrintMsgProto &S2CFailTargetNoTimesInviteProto fail")
		}

		p.Id = nil

		return p, nil

	case 35: //s2c_fail_invite
		return toErrCodeMessage(22, 35, data), nil

	case 72: //s2c_invite_all
		p := &pb22.S2CInviteAllProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_invite_all PrintMsgProto &S2CInviteAllProto fail")
		}

		p.Id = nil

		return p, nil

	case 73: //s2c_fail_invite_all
		return toErrCodeMessage(22, 73, data), nil

	case 36: //s2c_receive_invite
		p := &pb22.S2CReceiveInviteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_receive_invite PrintMsgProto &S2CReceiveInviteProto fail")
		}

		return p, nil

	case 38: //s2c_request_invite_list
		p := &pb22.S2CRequestInviteListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_request_invite_list PrintMsgProto &S2CRequestInviteListProto fail")
		}

		p.InviteList = nil

		return p, nil

	case 40: //s2c_request_team_detail
		p := &pb22.S2CRequestTeamDetailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_request_team_detail PrintMsgProto &S2CRequestTeamDetailProto fail")
		}

		p.TeamDetail = nil

		return p, nil

	case 41: //s2c_fail_request_team_detail
		return toErrCodeMessage(22, 41, data), nil

	case 43: //s2c_broadcast_start_challenge
		p := &pb22.S2CBroadcastStartChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_broadcast_start_challenge PrintMsgProto &S2CBroadcastStartChallengeProto fail")
		}

		p.Result = nil

		return p, nil

	case 44: //s2c_fail_with_member_times_not_enough_start_challenge
		p := &pb22.S2CFailWithMemberTimesNotEnoughStartChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_with_member_times_not_enough_start_challenge PrintMsgProto &S2CFailWithMemberTimesNotEnoughStartChallengeProto fail")
		}

		p.Id = nil

		return p, nil

	case 45: //s2c_fail_with_member_help_times_not_enough_start_challenge
		p := &pb22.S2CFailWithMemberHelpTimesNotEnoughStartChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_with_member_help_times_not_enough_start_challenge PrintMsgProto &S2CFailWithMemberHelpTimesNotEnoughStartChallengeProto fail")
		}

		p.Id = nil

		return p, nil

	case 46: //s2c_fail_with_member_no_guild_start_challenge
		p := &pb22.S2CFailWithMemberNoGuildStartChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_with_member_no_guild_start_challenge PrintMsgProto &S2CFailWithMemberNoGuildStartChallengeProto fail")
		}

		p.Id = nil

		return p, nil

	case 47: //s2c_fail_with_member_not_my_guild_start_challenge
		p := &pb22.S2CFailWithMemberNotMyGuildStartChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_with_member_not_my_guild_start_challenge PrintMsgProto &S2CFailWithMemberNotMyGuildStartChallengeProto fail")
		}

		p.Id = nil

		return p, nil

	case 48: //s2c_fail_with_member_is_help_but_no_guild_start_challenge
		p := &pb22.S2CFailWithMemberIsHelpButNoGuildStartChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_with_member_is_help_but_no_guild_start_challenge PrintMsgProto &S2CFailWithMemberIsHelpButNoGuildStartChallengeProto fail")
		}

		p.Id = nil

		return p, nil

	case 49: //s2c_fail_with_member_is_help_but_no_guild_member_start_challenge
		p := &pb22.S2CFailWithMemberIsHelpButNoGuildMemberStartChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_fail_with_member_is_help_but_no_guild_member_start_challenge PrintMsgProto &S2CFailWithMemberIsHelpButNoGuildMemberStartChallengeProto fail")
		}

		p.Id = nil

		return p, nil

	case 50: //s2c_fail_start_challenge
		return toErrCodeMessage(22, 50, data), nil

	case 57: //s2c_team_expired
		return toStringMessage(22, 57), nil

	case 65: //s2c_team_destroyed_because_of_leader_leave
		return toStringMessage(22, 65), nil

	case 59: //s2c_quick_query_team_basic
		p := &pb22.S2CQuickQueryTeamBasicProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_quick_query_team_basic PrintMsgProto &S2CQuickQueryTeamBasicProto fail")
		}

		p.Basics = nil

		return p, nil

	case 60: //s2c_fail_quick_query_team_basic
		return toErrCodeMessage(22, 60, data), nil

	case 62: //s2c_change_guild_mode
		return toStringMessage(22, 62), nil

	case 63: //s2c_fail_change_guild_mode
		return toErrCodeMessage(22, 63, data), nil

	case 64: //s2c_change_guild_mode_broadcast
		p := &pb22.S2CChangeGuildModeBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_change_guild_mode_broadcast PrintMsgProto &S2CChangeGuildModeBroadcastProto fail")
		}

		return p, nil

	case 66: //s2c_member_troop_changed
		p := &pb22.S2CMemberTroopChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_member_troop_changed PrintMsgProto &S2CMemberTroopChangedProto fail")
		}

		p.Member = nil

		return p, nil

	case 75: //s2c_list_record
		p := &pb22.S2CListRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_list_record PrintMsgProto &S2CListRecordProto fail")
		}

		return p, nil

	case 80: //s2c_team_talk
		p := &pb22.S2CTeamTalkProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_team_talk PrintMsgProto &S2CTeamTalkProto fail")
		}

		return p, nil

	case 81: //s2c_fail_team_talk
		return toErrCodeMessage(22, 81, data), nil

	case 82: //s2c_team_who_talk
		p := &pb22.S2CTeamWhoTalkProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_team_who_talk PrintMsgProto &S2CTeamWhoTalkProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 83: //s2c_team_history_talk
		p := &pb22.S2CTeamHistoryTalkProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "secret_tower.s2c_team_history_talk PrintMsgProto &S2CTeamHistoryTalkProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_rank(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_request_rank
		p := &pb23.S2CRequestRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "rank.s2c_request_rank PrintMsgProto &S2CRequestRankProto fail")
		}

		p.Rank = nil

		return p, nil

	case 3: //s2c_fail_request_rank
		return toErrCodeMessage(23, 3, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_bai_zhan(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_query_bai_zhan_info
		p := &pb24.S2CQueryBaiZhanInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "bai_zhan.s2c_query_bai_zhan_info PrintMsgProto &S2CQueryBaiZhanInfoProto fail")
		}

		p.Data = nil

		return p, nil

	case 3: //s2c_fail_query_bai_zhan_info
		return toErrCodeMessage(24, 3, data), nil

	case 35: //s2c_clear_last_jun_xian
		return toStringMessage(24, 35), nil

	case 5: //s2c_bai_zhan_challenge
		p := &pb24.S2CBaiZhanChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "bai_zhan.s2c_bai_zhan_challenge PrintMsgProto &S2CBaiZhanChallengeProto fail")
		}

		p.Share = nil

		return p, nil

	case 6: //s2c_fail_bai_zhan_challenge
		return toErrCodeMessage(24, 6, data), nil

	case 8: //s2c_collect_salary
		return toStringMessage(24, 8), nil

	case 9: //s2c_fail_collect_salary
		return toErrCodeMessage(24, 9, data), nil

	case 11: //s2c_collect_jun_xian_prize
		p := &pb24.S2CCollectJunXianPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "bai_zhan.s2c_collect_jun_xian_prize PrintMsgProto &S2CCollectJunXianPrizeProto fail")
		}

		return p, nil

	case 12: //s2c_fail_collect_jun_xian_prize
		return toErrCodeMessage(24, 12, data), nil

	case 13: //s2c_reset
		return toStringMessage(24, 13), nil

	case 30: //s2c_self_record
		p := &pb24.S2CSelfRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "bai_zhan.s2c_self_record PrintMsgProto &S2CSelfRecordProto fail")
		}

		p.Data = nil

		return p, nil

	case 31: //s2c_no_change_self_record
		return toStringMessage(24, 31), nil

	case 32: //s2c_fail_self_record
		return toErrCodeMessage(24, 32, data), nil

	case 22: //s2c_self_defence_record_changed
		return toStringMessage(24, 22), nil

	case 24: //s2c_request_rank
		p := &pb24.S2CRequestRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "bai_zhan.s2c_request_rank PrintMsgProto &S2CRequestRankProto fail")
		}

		p.Data = nil

		return p, nil

	case 28: //s2c_fail_request_rank
		return toErrCodeMessage(24, 28, data), nil

	case 27: //s2c_request_self_rank
		p := &pb24.S2CRequestSelfRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "bai_zhan.s2c_request_self_rank PrintMsgProto &S2CRequestSelfRankProto fail")
		}

		return p, nil

	case 33: //s2c_max_jun_xian_level_changed
		p := &pb24.S2CMaxJunXianLevelChangedProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "bai_zhan.s2c_max_jun_xian_level_changed PrintMsgProto &S2CMaxJunXianLevelChangedProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_dungeon(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_challenge
		p := &pb26.S2CChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.s2c_challenge PrintMsgProto &S2CChallengeProto fail")
		}

		p.Share = nil

		p.Prize = nil

		return p, nil

	case 3: //s2c_fail_challenge
		return toErrCodeMessage(26, 3, data), nil

	case 16: //s2c_update_challenge_times
		p := &pb26.S2CUpdateChallengeTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.s2c_update_challenge_times PrintMsgProto &S2CUpdateChallengeTimesProto fail")
		}

		return p, nil

	case 5: //s2c_collect_chapter_prize
		p := &pb26.S2CCollectChapterPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.s2c_collect_chapter_prize PrintMsgProto &S2CCollectChapterPrizeProto fail")
		}

		return p, nil

	case 6: //s2c_fail_collect_chapter_prize
		return toErrCodeMessage(26, 6, data), nil

	case 14: //s2c_collect_pass_dungeon_prize
		p := &pb26.S2CCollectPassDungeonPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.s2c_collect_pass_dungeon_prize PrintMsgProto &S2CCollectPassDungeonPrizeProto fail")
		}

		return p, nil

	case 15: //s2c_fail_collect_pass_dungeon_prize
		return toErrCodeMessage(26, 15, data), nil

	case 8: //s2c_auto_challenge
		p := &pb26.S2CAutoChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.s2c_auto_challenge PrintMsgProto &S2CAutoChallengeProto fail")
		}

		p.Prizes = nil

		return p, nil

	case 9: //s2c_fail_auto_challenge
		return toErrCodeMessage(26, 9, data), nil

	case 18: //s2c_collect_chapter_star_prize
		p := &pb26.S2CCollectChapterStarPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dungeon.s2c_collect_chapter_star_prize PrintMsgProto &S2CCollectChapterStarPrizeProto fail")
		}

		p.Prize = nil

		return p, nil

	case 19: //s2c_fail_collect_chapter_star_prize
		return toErrCodeMessage(26, 19, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_country(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 17: //s2c_request_country_prestige
		p := &pb27.S2CRequestCountryPrestigeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_request_country_prestige PrintMsgProto &S2CRequestCountryPrestigeProto fail")
		}

		return p, nil

	case 18: //s2c_fail_request_country_prestige
		return toErrCodeMessage(27, 18, data), nil

	case 20: //s2c_request_countries
		p := &pb27.S2CRequestCountriesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_request_countries PrintMsgProto &S2CRequestCountriesProto fail")
		}

		return p, nil

	case 21: //s2c_fail_request_countries
		return toErrCodeMessage(27, 21, data), nil

	case 75: //s2c_countries_update_notice
		p := &pb27.S2CCountriesUpdateNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_countries_update_notice PrintMsgProto &S2CCountriesUpdateNoticeProto fail")
		}

		return p, nil

	case 23: //s2c_hero_change_country
		p := &pb27.S2CHeroChangeCountryProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_hero_change_country PrintMsgProto &S2CHeroChangeCountryProto fail")
		}

		return p, nil

	case 24: //s2c_fail_hero_change_country
		return toErrCodeMessage(27, 24, data), nil

	case 32: //s2c_country_detail
		p := &pb27.S2CCountryDetailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_country_detail PrintMsgProto &S2CCountryDetailProto fail")
		}

		return p, nil

	case 33: //s2c_fail_country_detail
		return toErrCodeMessage(27, 33, data), nil

	case 41: //s2c_official_appoint
		return toStringMessage(27, 41), nil

	case 42: //s2c_fail_official_appoint
		return toErrCodeMessage(27, 42, data), nil

	case 49: //s2c_official_appoint_notice
		p := &pb27.S2COfficialAppointNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_official_appoint_notice PrintMsgProto &S2COfficialAppointNoticeProto fail")
		}

		return p, nil

	case 50: //s2c_official_depose_notice
		return toStringMessage(27, 50), nil

	case 44: //s2c_official_depose
		return toStringMessage(27, 44), nil

	case 45: //s2c_fail_official_depose
		return toErrCodeMessage(27, 45, data), nil

	case 55: //s2c_official_leave
		return toStringMessage(27, 55), nil

	case 56: //s2c_fail_official_leave
		return toErrCodeMessage(27, 56, data), nil

	case 51: //s2c_country_host_changed_notice
		p := &pb27.S2CCountryHostChangedNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_country_host_changed_notice PrintMsgProto &S2CCountryHostChangedNoticeProto fail")
		}

		return p, nil

	case 53: //s2c_country_destroy_notice
		p := &pb27.S2CCountryDestroyNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_country_destroy_notice PrintMsgProto &S2CCountryDestroyNoticeProto fail")
		}

		return p, nil

	case 52: //s2c_king_changed_notice
		p := &pb27.S2CKingChangedNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_king_changed_notice PrintMsgProto &S2CKingChangedNoticeProto fail")
		}

		return p, nil

	case 47: //s2c_collect_official_salary
		p := &pb27.S2CCollectOfficialSalaryProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_collect_official_salary PrintMsgProto &S2CCollectOfficialSalaryProto fail")
		}

		return p, nil

	case 48: //s2c_fail_collect_official_salary
		return toErrCodeMessage(27, 48, data), nil

	case 73: //s2c_change_name_start
		return toStringMessage(27, 73), nil

	case 74: //s2c_fail_change_name_start
		return toErrCodeMessage(27, 74, data), nil

	case 60: //s2c_change_name_start_notice
		return toStringMessage(27, 60), nil

	case 62: //s2c_change_name_vote
		p := &pb27.S2CChangeNameVoteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_change_name_vote PrintMsgProto &S2CChangeNameVoteProto fail")
		}

		return p, nil

	case 63: //s2c_fail_change_name_vote
		return toErrCodeMessage(27, 63, data), nil

	case 64: //s2c_hero_change_name_vote_count_update_notice
		p := &pb27.S2CHeroChangeNameVoteCountUpdateNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_hero_change_name_vote_count_update_notice PrintMsgProto &S2CHeroChangeNameVoteCountUpdateNoticeProto fail")
		}

		return p, nil

	case 65: //s2c_change_name_succ_notice
		p := &pb27.S2CChangeNameSuccNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_change_name_succ_notice PrintMsgProto &S2CChangeNameSuccNoticeProto fail")
		}

		return p, nil

	case 67: //s2c_search_to_appoint_hero_list
		p := &pb27.S2CSearchToAppointHeroListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_search_to_appoint_hero_list PrintMsgProto &S2CSearchToAppointHeroListProto fail")
		}

		return p, nil

	case 68: //s2c_fail_search_to_appoint_hero_list
		return toErrCodeMessage(27, 68, data), nil

	case 70: //s2c_default_to_appoint_hero_list
		p := &pb27.S2CDefaultToAppointHeroListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "country.s2c_default_to_appoint_hero_list PrintMsgProto &S2CDefaultToAppointHeroListProto fail")
		}

		return p, nil

	case 71: //s2c_fail_default_to_appoint_hero_list
		return toErrCodeMessage(27, 71, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_tag(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_add_or_update_tag
		p := &pb29.S2CAddOrUpdateTagProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tag.s2c_add_or_update_tag PrintMsgProto &S2CAddOrUpdateTagProto fail")
		}

		p.Id = nil

		p.Record = nil

		p.Tag = nil

		return p, nil

	case 3: //s2c_fail_add_or_update_tag
		return toErrCodeMessage(29, 3, data), nil

	case 4: //s2c_other_tag_me
		p := &pb29.S2COtherTagMeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tag.s2c_other_tag_me PrintMsgProto &S2COtherTagMeProto fail")
		}

		p.Record = nil

		p.Tag = nil

		return p, nil

	case 6: //s2c_delete_tag
		p := &pb29.S2CDeleteTagProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "tag.s2c_delete_tag PrintMsgProto &S2CDeleteTagProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_garden(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_list_treasury_tree_hero
		p := &pb31.S2CListTreasuryTreeHeroProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "garden.s2c_list_treasury_tree_hero PrintMsgProto &S2CListTreasuryTreeHeroProto fail")
		}

		p.HeroId = nil

		p.HelpMeHeroId = nil

		return p, nil

	case 14: //s2c_list_help_me
		p := &pb31.S2CListHelpMeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "garden.s2c_list_help_me PrintMsgProto &S2CListHelpMeProto fail")
		}

		p.TargetId = nil

		p.HelpMeHeroId = nil

		return p, nil

	case 15: //s2c_fail_list_help_me
		return toErrCodeMessage(31, 15, data), nil

	case 4: //s2c_list_treasury_tree_times
		p := &pb31.S2CListTreasuryTreeTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "garden.s2c_list_treasury_tree_times PrintMsgProto &S2CListTreasuryTreeTimesProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 6: //s2c_water_treasury_tree
		p := &pb31.S2CWaterTreasuryTreeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "garden.s2c_water_treasury_tree PrintMsgProto &S2CWaterTreasuryTreeProto fail")
		}

		p.Target = nil

		return p, nil

	case 7: //s2c_fail_water_treasury_tree
		return toErrCodeMessage(31, 7, data), nil

	case 8: //s2c_update_self_treasury_tree_times
		p := &pb31.S2CUpdateSelfTreasuryTreeTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "garden.s2c_update_self_treasury_tree_times PrintMsgProto &S2CUpdateSelfTreasuryTreeTimesProto fail")
		}

		p.HelpMeHeroId = nil

		return p, nil

	case 9: //s2c_update_self_treasury_tree_full
		p := &pb31.S2CUpdateSelfTreasuryTreeFullProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "garden.s2c_update_self_treasury_tree_full PrintMsgProto &S2CUpdateSelfTreasuryTreeFullProto fail")
		}

		return p, nil

	case 11: //s2c_collect_treasury_tree_prize
		return toStringMessage(31, 11), nil

	case 12: //s2c_fail_collect_treasury_tree_prize
		return toErrCodeMessage(31, 12, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_zhengwu(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_start
		p := &pb32.S2CStartProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhengwu.s2c_start PrintMsgProto &S2CStartProto fail")
		}

		p.ZhengWu = nil

		return p, nil

	case 3: //s2c_fail_start
		return toErrCodeMessage(32, 3, data), nil

	case 5: //s2c_collect
		return toStringMessage(32, 5), nil

	case 6: //s2c_fail_collect
		return toErrCodeMessage(32, 6, data), nil

	case 8: //s2c_yuanbao_complete
		return toStringMessage(32, 8), nil

	case 9: //s2c_fail_yuanbao_complete
		return toErrCodeMessage(32, 9, data), nil

	case 11: //s2c_yuanbao_refresh
		return toStringMessage(32, 11), nil

	case 12: //s2c_fail_yuanbao_refresh
		return toErrCodeMessage(32, 12, data), nil

	case 13: //s2c_refresh
		p := &pb32.S2CRefreshProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhengwu.s2c_refresh PrintMsgProto &S2CRefreshProto fail")
		}

		p.NewZhengWu = nil

		return p, nil

	case 15: //s2c_vip_collect
		p := &pb32.S2CVipCollectProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhengwu.s2c_vip_collect PrintMsgProto &S2CVipCollectProto fail")
		}

		return p, nil

	case 16: //s2c_fail_vip_collect
		return toErrCodeMessage(32, 16, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_zhanjiang(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_open
		p := &pb33.S2COpenProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhanjiang.s2c_open PrintMsgProto &S2COpenProto fail")
		}

		return p, nil

	case 3: //s2c_fail_open
		return toErrCodeMessage(33, 3, data), nil

	case 5: //s2c_give_up
		return toStringMessage(33, 5), nil

	case 6: //s2c_fail_give_up
		return toErrCodeMessage(33, 6, data), nil

	case 8: //s2c_update_captain
		p := &pb33.S2CUpdateCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhanjiang.s2c_update_captain PrintMsgProto &S2CUpdateCaptainProto fail")
		}

		return p, nil

	case 9: //s2c_fail_update_captain
		return toErrCodeMessage(33, 9, data), nil

	case 11: //s2c_challenge
		p := &pb33.S2CChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhanjiang.s2c_challenge PrintMsgProto &S2CChallengeProto fail")
		}

		p.Share = nil

		return p, nil

	case 12: //s2c_fail_challenge
		return toErrCodeMessage(33, 12, data), nil

	case 13: //s2c_pass
		p := &pb33.S2CPassProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhanjiang.s2c_pass PrintMsgProto &S2CPassProto fail")
		}

		return p, nil

	case 14: //s2c_update_open_times
		p := &pb33.S2CUpdateOpenTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "zhanjiang.s2c_update_open_times PrintMsgProto &S2CUpdateOpenTimesProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_question(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_start
		return toStringMessage(34, 2), nil

	case 3: //s2c_fail_start
		return toErrCodeMessage(34, 3, data), nil

	case 5: //s2c_answer
		p := &pb34.S2CAnswerProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "question.s2c_answer PrintMsgProto &S2CAnswerProto fail")
		}

		return p, nil

	case 12: //s2c_fail_answer
		return toErrCodeMessage(34, 12, data), nil

	case 7: //s2c_next
		return toStringMessage(34, 7), nil

	case 8: //s2c_fail_next
		return toErrCodeMessage(34, 8, data), nil

	case 10: //s2c_get_prize
		return toStringMessage(34, 10), nil

	case 11: //s2c_fail_get_prize
		return toErrCodeMessage(34, 11, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_relation(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_add_relation
		p := &pb35.S2CAddRelationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_add_relation PrintMsgProto &S2CAddRelationProto fail")
		}

		p.Id = nil

		p.Proto = nil

		return p, nil

	case 3: //s2c_fail_add_relation
		return toErrCodeMessage(35, 3, data), nil

	case 9: //s2c_add_enemy
		p := &pb35.S2CAddEnemyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_add_enemy PrintMsgProto &S2CAddEnemyProto fail")
		}

		p.Id = nil

		return p, nil

	case 11: //s2c_remove_enemy
		p := &pb35.S2CRemoveEnemyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_remove_enemy PrintMsgProto &S2CRemoveEnemyProto fail")
		}

		p.Id = nil

		return p, nil

	case 12: //s2c_fail_remove_enemy
		return toErrCodeMessage(35, 12, data), nil

	case 5: //s2c_remove_relation
		p := &pb35.S2CRemoveRelationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_remove_relation PrintMsgProto &S2CRemoveRelationProto fail")
		}

		p.Id = nil

		return p, nil

	case 6: //s2c_fail_remove_relation
		return toErrCodeMessage(35, 6, data), nil

	case 8: //s2c_list_relation
		p := &pb35.S2CListRelationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_list_relation PrintMsgProto &S2CListRelationProto fail")
		}

		p.Proto = nil

		return p, nil

	case 29: //s2c_new_list_relation
		p := &pb35.S2CNewListRelationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_new_list_relation PrintMsgProto &S2CNewListRelationProto fail")
		}

		p.Proto = nil

		return p, nil

	case 17: //s2c_recommend_hero_list
		p := &pb35.S2CRecommendHeroListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_recommend_hero_list PrintMsgProto &S2CRecommendHeroListProto fail")
		}

		return p, nil

	case 18: //s2c_fail_recommend_hero_list
		return toErrCodeMessage(35, 18, data), nil

	case 23: //s2c_search_heros
		p := &pb35.S2CSearchHerosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_search_heros PrintMsgProto &S2CSearchHerosProto fail")
		}

		return p, nil

	case 24: //s2c_fail_search_heros
		return toErrCodeMessage(35, 24, data), nil

	case 26: //s2c_search_hero_by_id
		p := &pb35.S2CSearchHeroByIdProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_search_hero_by_id PrintMsgProto &S2CSearchHeroByIdProto fail")
		}

		return p, nil

	case 27: //s2c_fail_search_hero_by_id
		return toErrCodeMessage(35, 27, data), nil

	case 34: //s2c_set_important_friend
		p := &pb35.S2CSetImportantFriendProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_set_important_friend PrintMsgProto &S2CSetImportantFriendProto fail")
		}

		p.Id = nil

		return p, nil

	case 35: //s2c_fail_set_important_friend
		return toErrCodeMessage(35, 35, data), nil

	case 37: //s2c_cancel_important_friend
		p := &pb35.S2CCancelImportantFriendProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "relation.s2c_cancel_important_friend PrintMsgProto &S2CCancelImportantFriendProto fail")
		}

		p.Id = nil

		return p, nil

	case 38: //s2c_fail_cancel_important_friend
		return toErrCodeMessage(35, 38, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_xiongnu(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_set_defender
		p := &pb36.S2CSetDefenderProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.s2c_set_defender PrintMsgProto &S2CSetDefenderProto fail")
		}

		p.Id = nil

		return p, nil

	case 3: //s2c_broacast_set_defender
		p := &pb36.S2CBroacastSetDefenderProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.s2c_broacast_set_defender PrintMsgProto &S2CBroacastSetDefenderProto fail")
		}

		p.Id = nil

		return p, nil

	case 4: //s2c_fail_set_defender
		return toErrCodeMessage(36, 4, data), nil

	case 6: //s2c_start
		p := &pb36.S2CStartProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.s2c_start PrintMsgProto &S2CStartProto fail")
		}

		return p, nil

	case 8: //s2c_broadcast_start
		p := &pb36.S2CBroadcastStartProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.s2c_broadcast_start PrintMsgProto &S2CBroadcastStartProto fail")
		}

		return p, nil

	case 7: //s2c_fail_start
		return toErrCodeMessage(36, 7, data), nil

	case 9: //s2c_info_broadcast
		p := &pb36.S2CInfoBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.s2c_info_broadcast PrintMsgProto &S2CInfoBroadcastProto fail")
		}

		p.Info = nil

		return p, nil

	case 11: //s2c_troop_info
		p := &pb36.S2CTroopInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.s2c_troop_info PrintMsgProto &S2CTroopInfoProto fail")
		}

		p.BaseTroops = nil

		return p, nil

	case 12: //s2c_fail_troop_info
		return toErrCodeMessage(36, 12, data), nil

	case 13: //s2c_end_broadcast
		p := &pb36.S2CEndBroadcastProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.s2c_end_broadcast PrintMsgProto &S2CEndBroadcastProto fail")
		}

		p.ResistXiongNu = nil

		return p, nil

	case 15: //s2c_get_xiong_nu_npc_base_info
		p := &pb36.S2CGetXiongNuNpcBaseInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.s2c_get_xiong_nu_npc_base_info PrintMsgProto &S2CGetXiongNuNpcBaseInfoProto fail")
		}

		return p, nil

	case 16: //s2c_fail_get_xiong_nu_npc_base_info
		return toErrCodeMessage(36, 16, data), nil

	case 18: //s2c_get_defenser_fight_amount
		p := &pb36.S2CGetDefenserFightAmountProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.s2c_get_defenser_fight_amount PrintMsgProto &S2CGetDefenserFightAmountProto fail")
		}

		p.DefenserId = nil

		return p, nil

	case 20: //s2c_get_xiong_nu_fight_info
		p := &pb36.S2CGetXiongNuFightInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xiongnu.s2c_get_xiong_nu_fight_info PrintMsgProto &S2CGetXiongNuFightInfoProto fail")
		}

		p.Data = nil

		return p, nil

	case 21: //s2c_fail_get_xiong_nu_fight_info
		return toErrCodeMessage(36, 21, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_survey(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 1: //s2c_complete
		p := &pb37.S2CCompleteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "survey.s2c_complete PrintMsgProto &S2CCompleteProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_farm(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 50: //s2c_farm_is_update
		return toStringMessage(38, 50), nil

	case 3: //s2c_plant
		p := &pb38.S2CPlantProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_plant PrintMsgProto &S2CPlantProto fail")
		}

		return p, nil

	case 4: //s2c_fail_plant
		return toErrCodeMessage(38, 4, data), nil

	case 6: //s2c_harvest
		p := &pb38.S2CHarvestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_harvest PrintMsgProto &S2CHarvestProto fail")
		}

		return p, nil

	case 7: //s2c_fail_harvest
		return toErrCodeMessage(38, 7, data), nil

	case 9: //s2c_change
		p := &pb38.S2CChangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_change PrintMsgProto &S2CChangeProto fail")
		}

		return p, nil

	case 10: //s2c_fail_change
		return toErrCodeMessage(38, 10, data), nil

	case 13: //s2c_one_key_plant
		p := &pb38.S2COneKeyPlantProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_one_key_plant PrintMsgProto &S2COneKeyPlantProto fail")
		}

		return p, nil

	case 14: //s2c_fail_one_key_plant
		return toErrCodeMessage(38, 14, data), nil

	case 29: //s2c_one_key_harvest
		p := &pb38.S2COneKeyHarvestProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_one_key_harvest PrintMsgProto &S2COneKeyHarvestProto fail")
		}

		return p, nil

	case 30: //s2c_fail_one_key_harvest
		return toErrCodeMessage(38, 30, data), nil

	case 53: //s2c_one_key_reset
		p := &pb38.S2COneKeyResetProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_one_key_reset PrintMsgProto &S2COneKeyResetProto fail")
		}

		return p, nil

	case 54: //s2c_fail_one_key_reset
		return toErrCodeMessage(38, 54, data), nil

	case 44: //s2c_view_farm
		p := &pb38.S2CViewFarmProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_view_farm PrintMsgProto &S2CViewFarmProto fail")
		}

		p.Target = nil

		p.TargetBasic = nil

		p.HeroFarm = nil

		return p, nil

	case 45: //s2c_fail_view_farm
		return toErrCodeMessage(38, 45, data), nil

	case 19: //s2c_steal
		p := &pb38.S2CStealProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_steal PrintMsgProto &S2CStealProto fail")
		}

		p.Target = nil

		return p, nil

	case 20: //s2c_fail_steal
		return toErrCodeMessage(38, 20, data), nil

	case 21: //s2c_who_steal_from_me
		p := &pb38.S2CWhoStealFromMeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_who_steal_from_me PrintMsgProto &S2CWhoStealFromMeProto fail")
		}

		p.Target = nil

		return p, nil

	case 32: //s2c_one_key_steal
		p := &pb38.S2COneKeyStealProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_one_key_steal PrintMsgProto &S2COneKeyStealProto fail")
		}

		return p, nil

	case 37: //s2c_fail_one_key_steal
		return toErrCodeMessage(38, 37, data), nil

	case 33: //s2c_who_one_key_steal_from_me
		p := &pb38.S2CWhoOneKeyStealFromMeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_who_one_key_steal_from_me PrintMsgProto &S2CWhoOneKeyStealFromMeProto fail")
		}

		p.TargetId = nil

		return p, nil

	case 40: //s2c_steal_log_list
		p := &pb38.S2CStealLogListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_steal_log_list PrintMsgProto &S2CStealLogListProto fail")
		}

		p.Logs = nil

		return p, nil

	case 42: //s2c_fail_steal_log_list
		return toErrCodeMessage(38, 42, data), nil

	case 46: //s2c_can_steal_list
		p := &pb38.S2CCanStealListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "farm.s2c_can_steal_list PrintMsgProto &S2CCanStealListProto fail")
		}

		p.CanStealId = nil

		return p, nil

	case 49: //s2c_fail_can_steal_list
		return toErrCodeMessage(38, 49, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_dianquan(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_exchange
		p := &pb39.S2CExchangeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "dianquan.s2c_exchange PrintMsgProto &S2CExchangeProto fail")
		}

		return p, nil

	case 3: //s2c_fail_exchange
		return toErrCodeMessage(39, 3, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_xuanyuan(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 26: //s2c_rank_is_empty
		p := &pb40.S2CRankIsEmptyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.s2c_rank_is_empty PrintMsgProto &S2CRankIsEmptyProto fail")
		}

		return p, nil

	case 2: //s2c_self_info
		p := &pb40.S2CSelfInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.s2c_self_info PrintMsgProto &S2CSelfInfoProto fail")
		}

		p.Targets = nil

		return p, nil

	case 12: //s2c_list_target
		p := &pb40.S2CListTargetProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.s2c_list_target PrintMsgProto &S2CListTargetProto fail")
		}

		p.Targets = nil

		return p, nil

	case 13: //s2c_fail_list_target
		return toErrCodeMessage(40, 13, data), nil

	case 6: //s2c_query_target_troop
		p := &pb40.S2CQueryTargetTroopProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.s2c_query_target_troop PrintMsgProto &S2CQueryTargetTroopProto fail")
		}

		p.Id = nil

		p.Player = nil

		return p, nil

	case 14: //s2c_fail_query_target_troop
		return toErrCodeMessage(40, 14, data), nil

	case 16: //s2c_challenge
		p := &pb40.S2CChallengeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.s2c_challenge PrintMsgProto &S2CChallengeProto fail")
		}

		p.Id = nil

		return p, nil

	case 17: //s2c_fail_challenge
		return toErrCodeMessage(40, 17, data), nil

	case 18: //s2c_update_xy_info
		p := &pb40.S2CUpdateXyInfoProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.s2c_update_xy_info PrintMsgProto &S2CUpdateXyInfoProto fail")
		}

		return p, nil

	case 19: //s2c_add_record
		p := &pb40.S2CAddRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.s2c_add_record PrintMsgProto &S2CAddRecordProto fail")
		}

		p.Data = nil

		return p, nil

	case 21: //s2c_list_record
		p := &pb40.S2CListRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.s2c_list_record PrintMsgProto &S2CListRecordProto fail")
		}

		p.Data = nil

		return p, nil

	case 23: //s2c_collect_rank_prize
		p := &pb40.S2CCollectRankPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "xuanyuan.s2c_collect_rank_prize PrintMsgProto &S2CCollectRankPrizeProto fail")
		}

		return p, nil

	case 24: //s2c_fail_collect_rank_prize
		return toErrCodeMessage(40, 24, data), nil

	case 25: //s2c_reset
		return toStringMessage(40, 25), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_hebi(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_room_list
		p := &pb41.S2CRoomListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_room_list PrintMsgProto &S2CRoomListProto fail")
		}

		return p, nil

	case 36: //s2c_hero_record_list
		p := &pb41.S2CHeroRecordListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_hero_record_list PrintMsgProto &S2CHeroRecordListProto fail")
		}

		return p, nil

	case 4: //s2c_change_captain
		p := &pb41.S2CChangeCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_change_captain PrintMsgProto &S2CChangeCaptainProto fail")
		}

		return p, nil

	case 5: //s2c_fail_change_captain
		return toErrCodeMessage(41, 5, data), nil

	case 34: //s2c_change_room_captain
		p := &pb41.S2CChangeRoomCaptainProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_change_room_captain PrintMsgProto &S2CChangeRoomCaptainProto fail")
		}

		return p, nil

	case 29: //s2c_check_in_room
		p := &pb41.S2CCheckInRoomProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_check_in_room PrintMsgProto &S2CCheckInRoomProto fail")
		}

		return p, nil

	case 30: //s2c_fail_check_in_room
		return toErrCodeMessage(41, 30, data), nil

	case 32: //s2c_copy_self
		p := &pb41.S2CCopySelfProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_copy_self PrintMsgProto &S2CCopySelfProto fail")
		}

		return p, nil

	case 33: //s2c_fail_copy_self
		return toErrCodeMessage(41, 33, data), nil

	case 10: //s2c_join_room
		p := &pb41.S2CJoinRoomProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_join_room PrintMsgProto &S2CJoinRoomProto fail")
		}

		return p, nil

	case 11: //s2c_fail_join_room
		return toErrCodeMessage(41, 11, data), nil

	case 25: //s2c_someone_joined_room
		p := &pb41.S2CSomeoneJoinedRoomProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_someone_joined_room PrintMsgProto &S2CSomeoneJoinedRoomProto fail")
		}

		return p, nil

	case 13: //s2c_rob_pos
		p := &pb41.S2CRobPosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_rob_pos PrintMsgProto &S2CRobPosProto fail")
		}

		return p, nil

	case 14: //s2c_fail_rob_pos
		return toErrCodeMessage(41, 14, data), nil

	case 26: //s2c_someone_robbed_my_pos
		p := &pb41.S2CSomeoneRobbedMyPosProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_someone_robbed_my_pos PrintMsgProto &S2CSomeoneRobbedMyPosProto fail")
		}

		return p, nil

	case 19: //s2c_leave_room
		p := &pb41.S2CLeaveRoomProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_leave_room PrintMsgProto &S2CLeaveRoomProto fail")
		}

		return p, nil

	case 20: //s2c_fail_leave_room
		return toErrCodeMessage(41, 20, data), nil

	case 22: //s2c_rob
		p := &pb41.S2CRobProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_rob PrintMsgProto &S2CRobProto fail")
		}

		return p, nil

	case 23: //s2c_fail_rob
		return toErrCodeMessage(41, 23, data), nil

	case 27: //s2c_someone_robbed_my_prize
		p := &pb41.S2CSomeoneRobbedMyPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_someone_robbed_my_prize PrintMsgProto &S2CSomeoneRobbedMyPrizeProto fail")
		}

		return p, nil

	case 24: //s2c_complete
		p := &pb41.S2CCompleteProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_complete PrintMsgProto &S2CCompleteProto fail")
		}

		return p, nil

	case 38: //s2c_view_show_prize
		p := &pb41.S2CViewShowPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "hebi.s2c_view_show_prize PrintMsgProto &S2CViewShowPrizeProto fail")
		}

		return p, nil

	case 39: //s2c_fail_view_show_prize
		return toErrCodeMessage(41, 39, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_mingc(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 5: //s2c_mingc_list
		p := &pb42.S2CMingcListProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc.s2c_mingc_list PrintMsgProto &S2CMingcListProto fail")
		}

		return p, nil

	case 6: //s2c_fail_mingc_list
		return toErrCodeMessage(42, 6, data), nil

	case 8: //s2c_view_mingc
		p := &pb42.S2CViewMingcProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc.s2c_view_mingc PrintMsgProto &S2CViewMingcProto fail")
		}

		return p, nil

	case 9: //s2c_fail_view_mingc
		return toErrCodeMessage(42, 9, data), nil

	case 11: //s2c_mc_build
		p := &pb42.S2CMcBuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc.s2c_mc_build PrintMsgProto &S2CMcBuildProto fail")
		}

		return p, nil

	case 12: //s2c_fail_mc_build
		return toErrCodeMessage(42, 12, data), nil

	case 14: //s2c_mc_build_log
		p := &pb42.S2CMcBuildLogProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc.s2c_mc_build_log PrintMsgProto &S2CMcBuildLogProto fail")
		}

		return p, nil

	case 15: //s2c_fail_mc_build_log
		return toErrCodeMessage(42, 15, data), nil

	case 16: //s2c_reset_daily_mc
		return toStringMessage(42, 16), nil

	case 21: //s2c_mingc_host_guild
		p := &pb42.S2CMingcHostGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc.s2c_mingc_host_guild PrintMsgProto &S2CMingcHostGuildProto fail")
		}

		return p, nil

	case 22: //s2c_fail_mingc_host_guild
		return toErrCodeMessage(42, 22, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_promotion(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 5: //s2c_collect_login_day_prize
		p := &pb43.S2CCollectLoginDayPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.s2c_collect_login_day_prize PrintMsgProto &S2CCollectLoginDayPrizeProto fail")
		}

		return p, nil

	case 6: //s2c_fail_collect_login_day_prize
		return toErrCodeMessage(43, 6, data), nil

	case 8: //s2c_buy_level_fund
		return toStringMessage(43, 8), nil

	case 9: //s2c_fail_buy_level_fund
		return toErrCodeMessage(43, 9, data), nil

	case 11: //s2c_collect_level_fund
		p := &pb43.S2CCollectLevelFundProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.s2c_collect_level_fund PrintMsgProto &S2CCollectLevelFundProto fail")
		}

		return p, nil

	case 12: //s2c_fail_collect_level_fund
		return toErrCodeMessage(43, 12, data), nil

	case 14: //s2c_collect_daily_sp
		p := &pb43.S2CCollectDailySpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.s2c_collect_daily_sp PrintMsgProto &S2CCollectDailySpProto fail")
		}

		return p, nil

	case 15: //s2c_fail_collect_daily_sp
		return toErrCodeMessage(43, 15, data), nil

	case 17: //s2c_collect_free_gift
		p := &pb43.S2CCollectFreeGiftProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.s2c_collect_free_gift PrintMsgProto &S2CCollectFreeGiftProto fail")
		}

		p.Prize = nil

		return p, nil

	case 18: //s2c_fail_collect_free_gift
		return toErrCodeMessage(43, 18, data), nil

	case 24: //s2c_notice_time_limit_gifts
		p := &pb43.S2CNoticeTimeLimitGiftsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.s2c_notice_time_limit_gifts PrintMsgProto &S2CNoticeTimeLimitGiftsProto fail")
		}

		return p, nil

	case 22: //s2c_buy_time_limit_gift
		p := &pb43.S2CBuyTimeLimitGiftProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.s2c_buy_time_limit_gift PrintMsgProto &S2CBuyTimeLimitGiftProto fail")
		}

		return p, nil

	case 23: //s2c_fail_buy_time_limit_gift
		return toErrCodeMessage(43, 23, data), nil

	case 25: //s2c_notice_event_limit_gift
		p := &pb43.S2CNoticeEventLimitGiftProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.s2c_notice_event_limit_gift PrintMsgProto &S2CNoticeEventLimitGiftProto fail")
		}

		return p, nil

	case 27: //s2c_buy_event_limit_gift
		p := &pb43.S2CBuyEventLimitGiftProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "promotion.s2c_buy_event_limit_gift PrintMsgProto &S2CBuyEventLimitGiftProto fail")
		}

		return p, nil

	case 28: //s2c_fail_buy_event_limit_gift
		return toErrCodeMessage(43, 28, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_mingc_war(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 32: //s2c_view_mc_war_self_guild
		p := &pb44.S2CViewMcWarSelfGuildProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_view_mc_war_self_guild PrintMsgProto &S2CViewMcWarSelfGuildProto fail")
		}

		return p, nil

	case 33: //s2c_fail_view_mc_war_self_guild
		return toErrCodeMessage(44, 33, data), nil

	case 30: //s2c_view_mc_war
		p := &pb44.S2CViewMcWarProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_view_mc_war PrintMsgProto &S2CViewMcWarProto fail")
		}

		return p, nil

	case 17: //s2c_apply_atk
		p := &pb44.S2CApplyAtkProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_apply_atk PrintMsgProto &S2CApplyAtkProto fail")
		}

		return p, nil

	case 18: //s2c_fail_apply_atk
		return toErrCodeMessage(44, 18, data), nil

	case 19: //s2c_apply_atk_succ
		p := &pb44.S2CApplyAtkSuccProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_apply_atk_succ PrintMsgProto &S2CApplyAtkSuccProto fail")
		}

		return p, nil

	case 20: //s2c_apply_atk_fail
		p := &pb44.S2CApplyAtkFailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_apply_atk_fail PrintMsgProto &S2CApplyAtkFailProto fail")
		}

		return p, nil

	case 22: //s2c_apply_ast
		p := &pb44.S2CApplyAstProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_apply_ast PrintMsgProto &S2CApplyAstProto fail")
		}

		return p, nil

	case 23: //s2c_fail_apply_ast
		return toErrCodeMessage(44, 23, data), nil

	case 24: //s2c_receive_apply_ast
		p := &pb44.S2CReceiveApplyAstProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_receive_apply_ast PrintMsgProto &S2CReceiveApplyAstProto fail")
		}

		return p, nil

	case 81: //s2c_cancel_apply_ast
		p := &pb44.S2CCancelApplyAstProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_cancel_apply_ast PrintMsgProto &S2CCancelApplyAstProto fail")
		}

		return p, nil

	case 82: //s2c_fail_cancel_apply_ast
		return toErrCodeMessage(44, 82, data), nil

	case 83: //s2c_receive_cancel_apply_ast
		p := &pb44.S2CReceiveCancelApplyAstProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_receive_cancel_apply_ast PrintMsgProto &S2CReceiveCancelApplyAstProto fail")
		}

		return p, nil

	case 26: //s2c_reply_apply_ast
		p := &pb44.S2CReplyApplyAstProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_reply_apply_ast PrintMsgProto &S2CReplyApplyAstProto fail")
		}

		return p, nil

	case 27: //s2c_fail_reply_apply_ast
		return toErrCodeMessage(44, 27, data), nil

	case 28: //s2c_apply_ast_pass
		p := &pb44.S2CApplyAstPassProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_apply_ast_pass PrintMsgProto &S2CApplyAstPassProto fail")
		}

		return p, nil

	case 104: //s2c_mingc_war_fight_prepare_start
		p := &pb44.S2CMingcWarFightPrepareStartProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_mingc_war_fight_prepare_start PrintMsgProto &S2CMingcWarFightPrepareStartProto fail")
		}

		return p, nil

	case 41: //s2c_mingc_war_fight_start
		p := &pb44.S2CMingcWarFightStartProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_mingc_war_fight_start PrintMsgProto &S2CMingcWarFightStartProto fail")
		}

		return p, nil

	case 66: //s2c_is_joining_fight_on_login
		p := &pb44.S2CIsJoiningFightOnLoginProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_is_joining_fight_on_login PrintMsgProto &S2CIsJoiningFightOnLoginProto fail")
		}

		return p, nil

	case 103: //s2c_red_point_notice
		return toStringMessage(44, 103), nil

	case 76: //s2c_view_mingc_war_mc
		p := &pb44.S2CViewMingcWarMcProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_view_mingc_war_mc PrintMsgProto &S2CViewMingcWarMcProto fail")
		}

		return p, nil

	case 77: //s2c_fail_view_mingc_war_mc
		return toErrCodeMessage(44, 77, data), nil

	case 36: //s2c_join_fight
		p := &pb44.S2CJoinFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_join_fight PrintMsgProto &S2CJoinFightProto fail")
		}

		return p, nil

	case 37: //s2c_fail_join_fight
		return toErrCodeMessage(44, 37, data), nil

	case 58: //s2c_other_join_fight
		p := &pb44.S2COtherJoinFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_other_join_fight PrintMsgProto &S2COtherJoinFightProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 39: //s2c_quit_fight
		p := &pb44.S2CQuitFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_quit_fight PrintMsgProto &S2CQuitFightProto fail")
		}

		return p, nil

	case 40: //s2c_fail_quit_fight
		return toErrCodeMessage(44, 40, data), nil

	case 57: //s2c_other_quit_fight
		p := &pb44.S2COtherQuitFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_other_quit_fight PrintMsgProto &S2COtherQuitFightProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 79: //s2c_scene_building_destroy_prosperity
		p := &pb44.S2CSceneBuildingDestroyProsperityProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_building_destroy_prosperity PrintMsgProto &S2CSceneBuildingDestroyProsperityProto fail")
		}

		return p, nil

	case 78: //s2c_scene_fight_prepare_end
		p := &pb44.S2CSceneFightPrepareEndProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_fight_prepare_end PrintMsgProto &S2CSceneFightPrepareEndProto fail")
		}

		return p, nil

	case 64: //s2c_scene_war_end
		p := &pb44.S2CSceneWarEndProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_war_end PrintMsgProto &S2CSceneWarEndProto fail")
		}

		return p, nil

	case 50: //s2c_scene_move
		p := &pb44.S2CSceneMoveProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_move PrintMsgProto &S2CSceneMoveProto fail")
		}

		return p, nil

	case 51: //s2c_fail_scene_move
		return toErrCodeMessage(44, 51, data), nil

	case 86: //s2c_scene_back
		p := &pb44.S2CSceneBackProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_back PrintMsgProto &S2CSceneBackProto fail")
		}

		return p, nil

	case 87: //s2c_fail_scene_back
		return toErrCodeMessage(44, 87, data), nil

	case 89: //s2c_scene_speed_up
		p := &pb44.S2CSceneSpeedUpProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_speed_up PrintMsgProto &S2CSceneSpeedUpProto fail")
		}

		p.Id = nil

		return p, nil

	case 90: //s2c_fail_scene_speed_up
		return toErrCodeMessage(44, 90, data), nil

	case 52: //s2c_scene_other_move
		p := &pb44.S2CSceneOtherMoveProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_other_move PrintMsgProto &S2CSceneOtherMoveProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 62: //s2c_scene_move_station
		p := &pb44.S2CSceneMoveStationProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_move_station PrintMsgProto &S2CSceneMoveStationProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 63: //s2c_scene_building_fight
		p := &pb44.S2CSceneBuildingFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_building_fight PrintMsgProto &S2CSceneBuildingFightProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 54: //s2c_scene_troop_relive
		p := &pb44.S2CSceneTroopReliveProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_troop_relive PrintMsgProto &S2CSceneTroopReliveProto fail")
		}

		return p, nil

	case 73: //s2c_fail_scene_troop_relive
		return toErrCodeMessage(44, 73, data), nil

	case 74: //s2c_scene_other_troop_relive
		p := &pb44.S2CSceneOtherTroopReliveProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_other_troop_relive PrintMsgProto &S2CSceneOtherTroopReliveProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 56: //s2c_scene_troop_update
		p := &pb44.S2CSceneTroopUpdateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_troop_update PrintMsgProto &S2CSceneTroopUpdateProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 47: //s2c_view_mc_war_scene
		p := &pb44.S2CViewMcWarSceneProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_view_mc_war_scene PrintMsgProto &S2CViewMcWarSceneProto fail")
		}

		return p, nil

	case 48: //s2c_fail_view_mc_war_scene
		return toErrCodeMessage(44, 48, data), nil

	case 140: //s2c_watch
		return toStringMessage(44, 140), nil

	case 141: //s2c_fail_watch
		return toErrCodeMessage(44, 141, data), nil

	case 137: //s2c_quit_watch
		return toStringMessage(44, 137), nil

	case 138: //s2c_fail_quit_watch
		return toErrCodeMessage(44, 138, data), nil

	case 67: //s2c_mc_war_end_record
		p := &pb44.S2CMcWarEndRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_mc_war_end_record PrintMsgProto &S2CMcWarEndRecordProto fail")
		}

		return p, nil

	case 92: //s2c_view_mc_war_record
		p := &pb44.S2CViewMcWarRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_view_mc_war_record PrintMsgProto &S2CViewMcWarRecordProto fail")
		}

		return p, nil

	case 93: //s2c_fail_view_mc_war_record
		return toErrCodeMessage(44, 93, data), nil

	case 95: //s2c_view_mc_war_troop_record
		p := &pb44.S2CViewMcWarTroopRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_view_mc_war_troop_record PrintMsgProto &S2CViewMcWarTroopRecordProto fail")
		}

		return p, nil

	case 96: //s2c_fail_view_mc_war_troop_record
		return toErrCodeMessage(44, 96, data), nil

	case 100: //s2c_view_scene_troop_record
		p := &pb44.S2CViewSceneTroopRecordProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_view_scene_troop_record PrintMsgProto &S2CViewSceneTroopRecordProto fail")
		}

		p.Record = nil

		return p, nil

	case 101: //s2c_fail_view_scene_troop_record
		return toErrCodeMessage(44, 101, data), nil

	case 105: //s2c_scene_troop_record_add_notice
		return toStringMessage(44, 105), nil

	case 110: //s2c_my_rank
		p := &pb44.S2CMyRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_my_rank PrintMsgProto &S2CMyRankProto fail")
		}

		return p, nil

	case 108: //s2c_apply_refresh_rank
		p := &pb44.S2CApplyRefreshRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_apply_refresh_rank PrintMsgProto &S2CApplyRefreshRankProto fail")
		}

		return p, nil

	case 109: //s2c_fail_apply_refresh_rank
		return toErrCodeMessage(44, 109, data), nil

	case 112: //s2c_view_my_guild_member_rank
		p := &pb44.S2CViewMyGuildMemberRankProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_view_my_guild_member_rank PrintMsgProto &S2CViewMyGuildMemberRankProto fail")
		}

		return p, nil

	case 113: //s2c_fail_view_my_guild_member_rank
		return toErrCodeMessage(44, 113, data), nil

	case 114: //s2c_cur_multi_kill
		p := &pb44.S2CCurMultiKillProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_cur_multi_kill PrintMsgProto &S2CCurMultiKillProto fail")
		}

		return p, nil

	case 135: //s2c_special_multi_kill
		p := &pb44.S2CSpecialMultiKillProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_special_multi_kill PrintMsgProto &S2CSpecialMultiKillProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 116: //s2c_scene_change_mode
		p := &pb44.S2CSceneChangeModeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_change_mode PrintMsgProto &S2CSceneChangeModeProto fail")
		}

		return p, nil

	case 117: //s2c_fail_scene_change_mode
		return toErrCodeMessage(44, 117, data), nil

	case 133: //s2c_scene_change_mode_notice
		p := &pb44.S2CSceneChangeModeNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_change_mode_notice PrintMsgProto &S2CSceneChangeModeNoticeProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 120: //s2c_scene_tou_shi_building_turn_to
		p := &pb44.S2CSceneTouShiBuildingTurnToProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_tou_shi_building_turn_to PrintMsgProto &S2CSceneTouShiBuildingTurnToProto fail")
		}

		return p, nil

	case 121: //s2c_fail_scene_tou_shi_building_turn_to
		return toErrCodeMessage(44, 121, data), nil

	case 122: //s2c_scene_tou_shi_building_turn_to_notice
		p := &pb44.S2CSceneTouShiBuildingTurnToNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_tou_shi_building_turn_to_notice PrintMsgProto &S2CSceneTouShiBuildingTurnToNoticeProto fail")
		}

		return p, nil

	case 124: //s2c_scene_tou_shi_building_fire
		p := &pb44.S2CSceneTouShiBuildingFireProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_tou_shi_building_fire PrintMsgProto &S2CSceneTouShiBuildingFireProto fail")
		}

		return p, nil

	case 125: //s2c_fail_scene_tou_shi_building_fire
		return toErrCodeMessage(44, 125, data), nil

	case 126: //s2c_scene_tou_shi_building_fire_notice
		p := &pb44.S2CSceneTouShiBuildingFireNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_tou_shi_building_fire_notice PrintMsgProto &S2CSceneTouShiBuildingFireNoticeProto fail")
		}

		return p, nil

	case 127: //s2c_scene_tou_shi_bomb_explode_notice
		p := &pb44.S2CSceneTouShiBombExplodeNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_tou_shi_bomb_explode_notice PrintMsgProto &S2CSceneTouShiBombExplodeNoticeProto fail")
		}

		p.FireTroopId = nil

		return p, nil

	case 129: //s2c_scene_drum
		p := &pb44.S2CSceneDrumProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_drum PrintMsgProto &S2CSceneDrumProto fail")
		}

		return p, nil

	case 130: //s2c_fail_scene_drum
		return toErrCodeMessage(44, 130, data), nil

	case 131: //s2c_scene_drum_notice
		p := &pb44.S2CSceneDrumNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_drum_notice PrintMsgProto &S2CSceneDrumNoticeProto fail")
		}

		p.HeroId = nil

		return p, nil

	case 132: //s2c_scene_drum_add_stat_notice
		p := &pb44.S2CSceneDrumAddStatNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_scene_drum_add_stat_notice PrintMsgProto &S2CSceneDrumAddStatNoticeProto fail")
		}

		p.Troops = nil

		return p, nil

	case 134: //s2c_mingc_host_update_notice
		p := &pb44.S2CMingcHostUpdateNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "mingc_war.s2c_mingc_host_update_notice PrintMsgProto &S2CMingcHostUpdateNoticeProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_random_event(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_choose_option
		p := &pb45.S2CChooseOptionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "random_event.s2c_choose_option PrintMsgProto &S2CChooseOptionProto fail")
		}

		p.Prize = nil

		return p, nil

	case 3: //s2c_fail_choose_option
		return toErrCodeMessage(45, 3, data), nil

	case 5: //s2c_open_event
		p := &pb45.S2COpenEventProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "random_event.s2c_open_event PrintMsgProto &S2COpenEventProto fail")
		}

		return p, nil

	case 6: //s2c_fail_open_event
		return toErrCodeMessage(45, 6, data), nil

	case 8: //s2c_new_event
		p := &pb45.S2CNewEventProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "random_event.s2c_new_event PrintMsgProto &S2CNewEventProto fail")
		}

		return p, nil

	case 9: //s2c_add_event_handbook
		p := &pb45.S2CAddEventHandbookProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "random_event.s2c_add_event_handbook PrintMsgProto &S2CAddEventHandbookProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_strategy(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_use_stratagem
		p := &pb46.S2CUseStratagemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "strategy.s2c_use_stratagem PrintMsgProto &S2CUseStratagemProto fail")
		}

		return p, nil

	case 3: //s2c_fail_use_stratagem
		return toErrCodeMessage(46, 3, data), nil

	case 4: //s2c_trapped_stratagem
		p := &pb46.S2CTrappedStratagemProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "strategy.s2c_trapped_stratagem PrintMsgProto &S2CTrappedStratagemProto fail")
		}

		return p, nil

	case 5: //s2c_use_stratagem_fail
		p := &pb46.S2CUseStratagemFailProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "strategy.s2c_use_stratagem_fail PrintMsgProto &S2CUseStratagemFailProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_vip(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 1: //s2c_vip_level_upgrade_notice
		p := &pb48.S2CVipLevelUpgradeNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "vip.s2c_vip_level_upgrade_notice PrintMsgProto &S2CVipLevelUpgradeNoticeProto fail")
		}

		return p, nil

	case 2: //s2c_vip_add_exp_notice
		p := &pb48.S2CVipAddExpNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "vip.s2c_vip_add_exp_notice PrintMsgProto &S2CVipAddExpNoticeProto fail")
		}

		return p, nil

	case 3: //s2c_vip_daily_login_notice
		p := &pb48.S2CVipDailyLoginNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "vip.s2c_vip_daily_login_notice PrintMsgProto &S2CVipDailyLoginNoticeProto fail")
		}

		return p, nil

	case 5: //s2c_vip_collect_daily_prize
		p := &pb48.S2CVipCollectDailyPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "vip.s2c_vip_collect_daily_prize PrintMsgProto &S2CVipCollectDailyPrizeProto fail")
		}

		return p, nil

	case 6: //s2c_fail_vip_collect_daily_prize
		return toErrCodeMessage(48, 6, data), nil

	case 8: //s2c_vip_collect_level_prize
		p := &pb48.S2CVipCollectLevelPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "vip.s2c_vip_collect_level_prize PrintMsgProto &S2CVipCollectLevelPrizeProto fail")
		}

		return p, nil

	case 9: //s2c_fail_vip_collect_level_prize
		return toErrCodeMessage(48, 9, data), nil

	case 13: //s2c_vip_buy_dungeon_times
		p := &pb48.S2CVipBuyDungeonTimesProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "vip.s2c_vip_buy_dungeon_times PrintMsgProto &S2CVipBuyDungeonTimesProto fail")
		}

		return p, nil

	case 14: //s2c_fail_vip_buy_dungeon_times
		return toErrCodeMessage(48, 14, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_red_packet(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_buy
		p := &pb49.S2CBuyProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "red_packet.s2c_buy PrintMsgProto &S2CBuyProto fail")
		}

		return p, nil

	case 3: //s2c_fail_buy
		return toErrCodeMessage(49, 3, data), nil

	case 5: //s2c_create
		p := &pb49.S2CCreateProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "red_packet.s2c_create PrintMsgProto &S2CCreateProto fail")
		}

		return p, nil

	case 6: //s2c_fail_create
		return toErrCodeMessage(49, 6, data), nil

	case 8: //s2c_grab
		p := &pb49.S2CGrabProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "red_packet.s2c_grab PrintMsgProto &S2CGrabProto fail")
		}

		return p, nil

	case 9: //s2c_fail_grab
		return toErrCodeMessage(49, 9, data), nil

	case 10: //s2c_all_grabbed_notice
		p := &pb49.S2CAllGrabbedNoticeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "red_packet.s2c_all_grabbed_notice PrintMsgProto &S2CAllGrabbedNoticeProto fail")
		}

		p.Id = nil

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_teach(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 2: //s2c_fight
		p := &pb50.S2CFightProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "teach.s2c_fight PrintMsgProto &S2CFightProto fail")
		}

		return p, nil

	case 3: //s2c_fail_fight
		return toErrCodeMessage(50, 3, data), nil

	case 5: //s2c_collect_prize
		p := &pb50.S2CCollectPrizeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "teach.s2c_collect_prize PrintMsgProto &S2CCollectPrizeProto fail")
		}

		return p, nil

	case 6: //s2c_fail_collect_prize
		return toErrCodeMessage(50, 6, data), nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func print_activity(sequenceID int, data []byte) (proto.Message, error) {
	switch sequenceID {

	case 9: //s2c_notice_activity_show
		p := &pb51.S2CNoticeActivityShowProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "activity.s2c_notice_activity_show PrintMsgProto &S2CNoticeActivityShowProto fail")
		}

		return p, nil

	case 2: //s2c_collect_collection
		p := &pb51.S2CCollectCollectionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "activity.s2c_collect_collection PrintMsgProto &S2CCollectCollectionProto fail")
		}

		return p, nil

	case 3: //s2c_fail_collect_collection
		return toErrCodeMessage(51, 3, data), nil

	case 5: //s2c_notice_collection
		p := &pb51.S2CNoticeCollectionProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "activity.s2c_notice_collection PrintMsgProto &S2CNoticeCollectionProto fail")
		}

		return p, nil

	case 6: //s2c_notice_collection_counts
		p := &pb51.S2CNoticeCollectionCountsProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "activity.s2c_notice_collection_counts PrintMsgProto &S2CNoticeCollectionCountsProto fail")
		}

		return p, nil

	case 7: //s2c_notice_task_list_mode
		p := &pb51.S2CNoticeTaskListModeProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "activity.s2c_notice_task_list_mode PrintMsgProto &S2CNoticeTaskListModeProto fail")
		}

		return p, nil

	case 8: //s2c_notice_task_list_mode_progress
		p := &pb51.S2CNoticeTaskListModeProgressProto{}
		if err := proto.UnmarshalMerge(data, p); err != nil {
			return nil, errors.Wrap(err, "activity.s2c_notice_task_list_mode_progress PrintMsgProto &S2CNoticeTaskListModeProgressProto fail")
		}

		return p, nil

	default:
		return nil, errors.Errorf("achieve打印未知消息: %d", sequenceID)
	}
}

func toStringMessage(moduleID, sequenceID int) stringProto {
	s := strconv.Itoa(moduleID) + "-" + strconv.Itoa(sequenceID)
	return stringProto(s)
}

func toErrCodeMessage(moduleID, sequenceID int, data []byte) stringProto {

	s := strconv.Itoa(moduleID) + "-" + strconv.Itoa(sequenceID)

	for _, b := range data {
		s += "-"
		s += strconv.Itoa(int(b))
	}

	return stringProto(s)
}

type stringProto string

func (s stringProto) Reset()         {}
func (s stringProto) String() string { return string(s) }
func (s stringProto) ProtoMessage()  {}

/*
----------- 自动识别分割线 -----------
package service

import (
	"github.com/lightpaw/protobuf/proto"
)

func PrintObject(moduleID, sequenceID int, data []byte) (proto.Message, error) {
	return nil, nil
}

----------- 自动识别分割线 -----------
*/
