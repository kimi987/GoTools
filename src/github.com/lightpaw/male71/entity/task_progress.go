package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"runtime/debug"
	"github.com/lightpaw/male7/config/fishing_data"
	"github.com/lightpaw/male7/config/shop"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/entity/npcid"
)

// progress
type TaskProgress struct {
	target *taskdata.TaskTargetData

	progress uint64
}

func (t *TaskProgress) encodeClient() uint64 {
	if t.target.DontUpdateProgress {
		return 0
	}
	return t.GetProgress()
}

func (t *TaskProgress) TargetType() shared_proto.TaskTargetType {
	return t.target.Type
}

func (t *TaskProgress) Target() *taskdata.TaskTargetData {
	return t.target
}

func (t *TaskProgress) GetProgress() uint64 {
	return u64.Min(t.progress, t.target.TotalProgress)
}

func (t *TaskProgress) IsCompleted() bool {
	return t.progress >= t.target.TotalProgress
}

func (t *TaskProgress) GmComplete() {
	t.setProgress(t.target.TotalProgress)
}

func (t *TaskProgress) setProgress(toSet uint64) bool {
	if t.progress != toSet {
		t.progress = toSet
		return true
	}

	return false
}

func (t *TaskProgress) IncreTaskTypeProgress(targetType shared_proto.TaskTargetType, hero *Hero, increAmount uint64) bool {
	if t.target.Type != targetType {
		return false
	}

	if _, exist := increTaskTargetTypes[targetType]; !exist {
		logrus.WithField("type", t.TargetType()).Debugf("未处理的任务incre类型")
		return false
	}

	return t.IncreTaskTypeProgressWithFunc(targetType, hero, func(hero *Hero, t *TaskProgress) uint64 {
		return increAmount
	})
}

func (t *TaskProgress) IncreTaskTypeProgressWithFunc(targetType shared_proto.TaskTargetType, hero *Hero, increFunc func(hero *Hero, t *TaskProgress) (increAmount uint64)) bool {
	if t.IsCompleted() {
		return false
	}

	if t.target.Type != targetType {
		return false
	}

	if increFunc == nil {
		logrus.WithField("type", t.TargetType()).WithField("stack", string(debug.Stack())).Debugf("未处理的任务incre类型")
		return false
	}

	updated := t.setProgress(t.progress + increFunc(hero, t))

	if t.target.DontUpdateProgress {
		return t.IsCompleted()
	}

	return updated
}

var increTaskTargetTypes = map[shared_proto.TaskTargetType]struct{}{
	shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL:      {},
	shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_TECH_LEVEL:          {},
	shared_proto.TaskTargetType_TASK_TARGET_RECRUIT_SOLDIER:             {},
	shared_proto.TaskTargetType_TASK_TARGET_JIU_GUAN_CONSULT:            {},
	shared_proto.TaskTargetType_TASK_TARGET_COLLECT_RESOURCE_TIMES:      {},
	shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_DUNGEON:           {},
	shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_TOWER:             {},
	shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_BAI_ZHAN:          {},
	shared_proto.TaskTargetType_TASK_TARGET_FIGHT_IN_JADE_REALM:         {},
	shared_proto.TaskTargetType_TASK_TARGET_CHAT_TIMES:                  {},
	shared_proto.TaskTargetType_TASK_TARGET_BUILD_EQUIP_DAILY:           {},
	shared_proto.TaskTargetType_TASK_TARGET_SMELT_EQUIP:                 {},
	shared_proto.TaskTargetType_TASK_TARGET_COLLECT_XIU_LIAN_EXP:        {},
	shared_proto.TaskTargetType_TASK_TARGET_GUILD_DONATE:                {},
	shared_proto.TaskTargetType_TASK_TARGET_ASSIST_GUILD_MEMBER:         {},
	shared_proto.TaskTargetType_TASK_TARGET_WATCH_VIDEO:                 {},
	shared_proto.TaskTargetType_TASK_TARGET_WORKER_SPEED_UP:             {},
	shared_proto.TaskTargetType_TASK_TARGET_FISHING:                     {},
	shared_proto.TaskTargetType_TASK_TARGET_HELP_GUILD_MEMBER:           {},
	shared_proto.TaskTargetType_TASK_TARGET_BUY_GOODS_COUNT:             {},
	shared_proto.TaskTargetType_TASK_TARGET_BUY_GOODS:                   {},
	shared_proto.TaskTargetType_TASK_TARGET_JADE_NPC:                    {},
	shared_proto.TaskTargetType_TASK_TARGET_ACTIVE_START_QUESTION_COUNT: {},
}

func (t *TaskProgress) UpdateTaskProgress(hero *Hero) bool {
	updateFunc := updateTaskProgressFuncs[t.TargetType()]
	if updateFunc == nil {
		return false
	}
	return t.UpdateTaskTypeProgressWithFunc(t.TargetType(), hero, updateFunc)
}

func (t *TaskProgress) UpdateTaskTypeProgress(targetType shared_proto.TaskTargetType, hero *Hero) bool {
	if t.target.Type != targetType {
		return false
	}

	return t.UpdateTaskTypeProgressWithFunc(targetType, hero, updateTaskProgressFuncs[targetType])
}

func (t *TaskProgress) UpdateTaskTypeProgressWithFunc(targetType shared_proto.TaskTargetType, hero *Hero, updateFunc func(hero *Hero, t *TaskProgress) (newProgress uint64)) bool {
	if t.IsCompleted() {
		return false
	}

	if t.target.Type != targetType {
		return false
	}

	if updateFunc == nil {
		logrus.WithField("type", t.TargetType()).WithField("stack", string(debug.Stack())).Debugf("未处理的任务update类型")
		return false
	}

	updated := t.setProgress(updateFunc(hero, t))

	if t.target.DontUpdateProgress {
		return t.IsCompleted()
	}

	return updated
}

var updateTaskProgressFuncs = map[shared_proto.TaskTargetType]func(hero *Hero, t *TaskProgress) (newProgress uint64){
	shared_proto.TaskTargetType_TASK_TARGET_BASE_LEVEL: func(hero *Hero, t *TaskProgress) (newProgress uint64) { return hero.BaseLevel() },
	shared_proto.TaskTargetType_TASK_TARGET_HERO_LEVEL: func(hero *Hero, t *TaskProgress) (newProgress uint64) { return hero.Level() },
	shared_proto.TaskTargetType_TASK_TARGET_TECH_LEVEL: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		tech := hero.domestic.GetTechnology(t.target.TechGroup)
		if tech != nil {
			return tech.Level
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		b := hero.domestic.GetBuilding(t.target.BuildingType)
		if b != nil {
			return b.Level
		}

		hero.domestic.OuterCities().WalkBuildings(func(layout *domestic_data.OuterCityLayoutData, building *domestic_data.OuterCityBuildingData) {
			if building.BuildingType == t.target.BuildingType || building.BuildingType == t.target.BuildingType2 {
				newProgress = u64.Max(newProgress, layout.Level)
				return
			}
		})

		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return u64.FromInt(hero.military.CaptainCount())
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_LEVEL_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		// 拥有X个Y级的武将
		for _, c := range hero.military.captains {
			if c.RebirthLevel() >= t.target.CaptainRebirth && c.Level() >= t.target.CaptainLevel {
				newProgress++
			}
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_QUALITY_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		// 拥有X个Y品质的武将
		for _, c := range hero.military.captains {
			if c.Quality() >= t.target.CaptainQuality {
				newProgress++
			}
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_EQUIPMENT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		// 拥有X个武将身上穿Y件Z级装备
		for _, c := range hero.military.captains {
			n := uint64(0)
			for _, e := range c.equipment {
				if e.Level() >= t.target.WearEquipmentLevel {
					n++
				}
			}

			if n >= t.target.WearEquipmentCount {
				newProgress++
			}
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_LEVEL_Y: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		// 任意X个穿上的装备达到Y级
		for _, c := range hero.military.captains {
			for _, e := range c.equipment {
				if e.Level() >= t.target.WearEquipmentLevel {
					newProgress++

					if newProgress >= t.target.TotalProgress {
						return
					}
				}
			}
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_REFINE_LEVEL_Y: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		// 任意X个穿上的装备强化到Y级
		for _, c := range hero.military.captains {
			for _, e := range c.equipment {
				if e.RefinedLevel() >= t.target.WearEquipmentRefineLevel {
					newProgress++

					if newProgress >= t.target.TotalProgress {
						return
					}
				}
			}
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_QUALITY_Y: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		// 任意穿上X色装备Y件
		for _, c := range hero.military.captains {
			for _, e := range c.equipment {
				if e.Data().Quality.GoodsQuality.Level >= t.target.WearEquipmentQuality {
					newProgress++

					if newProgress >= t.target.TotalProgress {
						return
					}
				}
			}
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_REFINED_TIMES: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_CaptainRefinedTimes)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_JADE_ORE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_RobJadeOre)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_JADE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_Jade)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_GUILD_CONTRIBUTION: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_GuildContributionCoin)
	},
	shared_proto.TaskTargetType_TASK_TARGET_RECRUIT_SOLDIER_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_RecruitSoldierCount)
	},
	shared_proto.TaskTargetType_TASK_TARGET_HEAL_SOLDIER_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_HealSoldierCount)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_GUILD_DONATE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_GuildDonate)
	},
	shared_proto.TaskTargetType_TASK_TARGET_EXPEL: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_Expel)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_RECOVER_PROSPERITY: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_RestoreProsperity)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_CONSULT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_Consult)
	},
	shared_proto.TaskTargetType_TASK_TARGET_DEFENSER_FIGHTING: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.GetHomeDefenserFightAmount()
	},
	shared_proto.TaskTargetType_TASK_TARGET_SET_DEFENSER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		if hero.Bools().Get(shared_proto.HeroBoolType_BOOL_SET_DEFENSER) {
			newProgress = t.target.TotalProgress
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_SOLDIER_LEVEL: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.military.SoldierLevel()
	},
	shared_proto.TaskTargetType_TASK_TARGET_TOWER_FLOOR: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.tower.historyMaxFloor
	},
	//shared_proto.TaskTargetType_TASK_TARGET_TRAINING_USE_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
	//	for _, v := range hero.military.training {
	//		if v.captain != 0 {
	//			newProgress++
	//		}
	//	}
	//
	//	return
	//},
	//shared_proto.TaskTargetType_TASK_TARGET_TRAINING_LEVEL_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
	//	for _, v := range hero.military.training {
	//		if v.data.Level >= t.target.TrainingLevel {
	//			newProgress++
	//		}
	//	}
	//
	//	return
	//},
	shared_proto.TaskTargetType_TASK_TARGET_RESOURCE_POINT_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		// 拥有X个Y类型资源点
		if t.target.ResBuildingType == shared_proto.BuildingType_InvalidBuildingType {
			return u64.FromInt(len(hero.domestic.resourcePoint))
		}

		for _, b := range hero.domestic.resourcePoint {
			if b.building.Type == t.target.ResBuildingType {
				newProgress++
			}
		}

		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_COLLECT_RESOURCE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		switch t.target.ResBuildingType {
		default:
			return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_CollectResource)
		case shared_proto.BuildingType_GOLD_PRODUCER:
			return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_CollectResourceGold)
		case shared_proto.BuildingType_FOOD_PRODUCER:
			return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_CollectResourceFood)
		case shared_proto.BuildingType_WOOD_PRODUCER:
			return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_CollectResourceWood)
		case shared_proto.BuildingType_STONE_PRODUCER:
			return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_CollectResourceStone)
		}
	},
	shared_proto.TaskTargetType_TASK_TARGET_JOIN_GUILD: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		if hero.GuildId() != 0 {
			return t.target.TotalProgress // 直接完成
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_GUI_ZU_LEVEL: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.VipLevel()
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_SOUL_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return t.target.TotalProgress
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_SOUL_QUALITY_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return t.target.TotalProgress
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_SOUL_LEVEL_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return t.target.TotalProgress
	},
	shared_proto.TaskTargetType_TASK_TARGET_HAS_CHALLENGE_DUNGEON: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		if hero.Dungeon().IsPass(t.target.PassDungeon) {
			return t.target.TotalProgress
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_FU_SHEN: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return t.target.TotalProgress
	},
	shared_proto.TaskTargetType_TASK_TARGET_BUILD_EQUIP: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_CombineEquip)
	},
	shared_proto.TaskTargetType_TASK_TARGET_KILL_HOME_NPC: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		dataId := t.target.KillHomeNpcData.Id
		if hero.HasCreateHomeNpcBase(dataId) {
			if hero.GetHomeNpcBase(npcid.NewHomeNpcId(dataId)) == nil {
				return t.target.TotalProgress
			}
		}
		return 0
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_FISHING: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_Fishing, fishing_data.FishTypeYuanbao) +
			hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_Fishing, fishing_data.FishTypeFree)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BAI_ZHAN: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumBaiZhan)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_SMELT_EQUIP: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_SmeltEquip)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_COLLECT_ACTIVE_BOX: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_CollectActiveBox)
	},
	shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_SECRET_TOWER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		var id uint64
		if t.target.PassSecretTower != nil {
			id = t.target.PassSecretTower.Id
		}
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_ChallengeSecretTower, id)
	},
	shared_proto.TaskTargetType_TASK_TARGET_HELP_SECRET_TOWER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		var id uint64
		if t.target.PassSecretTower != nil {
			id = t.target.PassSecretTower.Id
		}
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_HelpSecretTower, id)
	},
	shared_proto.TaskTargetType_TASK_TARGET_HISTORY_CHALLENGE_SECRET_TOWER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		var id uint64
		if t.target.PassSecretTower != nil {
			id = t.target.PassSecretTower.Id
		}
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_ChallengeSecretTower, id)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BUY_GOODS: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_ShopBuy, t.target.SubType)
	},
	shared_proto.TaskTargetType_TASK_TARGET_HOME_IN_GUILD_REGION: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		regionType := realmface.ParseRegionType(hero.home.realmId)
		if regionType == shared_proto.RegionType_GUILD {
			return t.target.TotalProgress
		}

		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BUY_GOODS_COUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_ShopGoodsBuy, shop.ShopTypeGoodsId(t.target.SubType, t.target.SubTypeId))
	},
	shared_proto.TaskTargetType_TASK_TARGET_BAOWU_SELL: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_SellBaowu, t.target.SubType)
	},
	shared_proto.TaskTargetType_TASK_TARGET_BAI_ZHAN_JUN_XIAN: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_MaxJunXianLevel)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACTIVE_DEGREE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_MaxActiveDegree)
	},
	shared_proto.TaskTargetType_TASK_TARGET_INVADE_MULTI_LEVEL_MONSTER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_InvadeMultiLevelMonster)
	},
	shared_proto.TaskTargetType_TASK_TARGET_WIN_MULTI_LEVEL_MONSTER: historyAmountFunc(server_proto.HistoryAmountType_WinMultiLevelMonster),

	shared_proto.TaskTargetType_TASK_TARGET_COMPLETE_ZHENG_WU: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_CompleteZhengWu)
	},
	shared_proto.TaskTargetType_TASK_TARGET_TREASURY_TREE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_TreasuryTree)
	},
	shared_proto.TaskTargetType_TASK_TARGET_INVASE_KILL_SOLDIER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_InvaseKillSoldier)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ASSIST_KILL_SOLDIER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AssistKillSoldier)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_LOGIN_DAY: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumLoginDays)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ALL_RIGHT_QUESTION_AMOUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_QuestionRightAmount)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_START_QUESTION: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumStartQuestion)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_HELP_WATER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumHelpWater)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BASE_DEAD: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumBaseDead)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_MOVE_BASE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumMoveBase)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_DESTROY_BASE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumDestroyBase)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_DESTROY_PROSPERITY: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumDestroyProsperity)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_ROBBING_RES: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_AccumRobbingRes, t.target.SubType)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_INVESTIGATION: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_Inverstigation)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BEEN_INVESTIGATION: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_BeenInverstigation)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_REALM_PVP_ASSIST: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_RealmPvpAssist)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_REALM_PVP_BEEN_ASSIST: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_RealmPvpBeenAssist)
	},
	shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_AccumFarmHarvestRes, t.target.SubType)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ZHANJIANG_START_TIMES: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumZhanJiangGuanqia)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_MC_WAR_KILL_SOLDIER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_McWarKillSolider)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_MC_WAR_WIN: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_McWarWin)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_MC_WAR_DESTROY_BUILDING: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_McWarDestroyBuilding)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ZHANJIANG_GUANQIA_COMPLETE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		if hero.ZhanJiang().IsPass(t.target.PassZhanjiangGuanqia) {
			return t.target.TotalProgress
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_START_XIONGNU: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_AccumXiongNuStart, t.target.SubType)
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_OFFICIAL_UPDATE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		if t.target.CaptainOfficial != nil {
			newProgress = hero.Military().GetOfficialCount(t.target.CaptainOfficial.Id)
		} else {
			for _, off := range hero.military.officialCounter.officials {
				newProgress += off.positionCount()
			}
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_FARM_STEAL: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, server_proto.HistoryAmountType_AccumFarmStealRes, t.target.SubType)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_AUTO_DUNGEON: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumDungeonAutoComplete)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_COUNT_DOWN_PRIZE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AccumCountDownPrizeTimes)
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_ABILITY_EXP: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		for _, c := range hero.Military().Captains() {
			if c != nil && c.abilityData != nil && c.abilityData.Ability >= t.target.CaptainAbilityExp {
				newProgress++
			}
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BWZL_COMPLETE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		newProgress = hero.TaskList().BwzlTaskList().CollectPrizeCount()
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_INVASE_KILL_SOLDIER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_InvaseKillSoldier)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_ASSIST_KILL_SOLDIER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_AssistKillSoldier)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_DEFENSE_KILL_SOLDIER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_DefenseKillSoldier)
	},
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_USE_TU_FEI_GOODS: boolFunc(
		shared_proto.HeroBoolType_BOOL_USE_TU_FEI_GOODS),
	shared_proto.TaskTargetType_TASK_TARGET_ACCUM_HELP_GUILD_MEMBER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_HelpGuildMember)
	},

	shared_proto.TaskTargetType_TASK_TARGET_ROB_MULTI_LEVEL_MONSTER: historyAmountFunc(
		server_proto.HistoryAmountType_RobMultiLevelMonsterRes),

	shared_proto.TaskTargetType_TASK_TARGET_EXPEL_FIGHT_MONSTER: historyAmountFunc(
		server_proto.HistoryAmountType_ExpelFightMonster),

	shared_proto.TaskTargetType_TASK_TARGET_UNLOCK_BAOWU: historyAmountFunc(
		server_proto.HistoryAmountType_UnlockBaowu),

	shared_proto.TaskTargetType_TASK_TARGET_ROB_BAOWU: historyAmountFunc(
		server_proto.HistoryAmountType_RobBaowu),

	shared_proto.TaskTargetType_TASK_TARGET_ROB_NPC_BAOWU: historyAmountFunc(
		server_proto.HistoryAmountType_RobNpcBaowu),

	shared_proto.TaskTargetType_TASK_TARGET_XUANYUAN_SCORE: historyAmountFunc(
		server_proto.HistoryAmountType_XuanyuanHisMaxScore),

	shared_proto.TaskTargetType_TASK_TARGET_HEBI: historyAmountFunc(
		server_proto.HistoryAmountType_StartHebi),

	shared_proto.TaskTargetType_TASK_TARGET_HEBI_ROB: historyAmountFunc(
		server_proto.HistoryAmountType_RobHebi),

	shared_proto.TaskTargetType_TASK_TARGET_GEM: func(hero *Hero, t *TaskProgress) (newProgress uint64) {

		// 增加拥有X个Y级Z类型宝石的任务条件
		for _, c := range hero.Military().Captains() {
			for _, v := range c.gems {
				if v != nil {
					if t.target.GemType != 0 && t.target.GemType != v.GemType {
						continue
					}

					if t.target.GemLevel > v.Level {
						continue
					}

					newProgress++
					if newProgress >= t.target.TotalProgress {
						return
					}
				}
			}
		}

		if newProgress < t.target.TotalProgress {
			// 遍历背包宝石
			for goodsId, count := range hero.depot.goodsMap {
				gemType, gemLevel := goods.GetGemTypeLevelByGoodsId(goodsId)
				if gemLevel <= 0 {
					continue
				}

				if t.target.GemType != 0 && t.target.GemType != gemType {
					continue
				}

				if t.target.GemLevel > gemLevel {
					continue
				}

				newProgress += count
				if newProgress >= t.target.TotalProgress {
					return
				}
			}
		}

		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_TITLE: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.TaskList().TitleId()
	},
	shared_proto.TaskTargetType_TASK_TARGET_FRIEND_AMOUNT: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.Relation().friendCount
	},
	shared_proto.TaskTargetType_TASK_TARGET_TEAM_POWER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().Amount(server_proto.HistoryAmountType_MaxTroopFightAmount)
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_UPSTAR: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		// 将X个武将升至Y星
		for _, c := range hero.military.captains {
			if c.Star() >= t.target.CaptainStar {
				newProgress++
			}
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_RNDEVENT_HANDBOOKS: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return u64.FromInt(len(hero.randomEvent.handbooks))
	},
	shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_EXP_GOODS_USE: historyAmountFunc(
		server_proto.HistoryAmountType_CaptainExpGoodsUsed),
	shared_proto.TaskTargetType_TASK_TARGET_STRATEGY_USE: historyAmountFunc(
		server_proto.HistoryAmountType_StrategyUsed),
	shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST_TIMES: historyAmountFunc(
		server_proto.HistoryAmountType_FarmHarvestTimes),
	shared_proto.TaskTargetType_TASK_TARGET_INVASE_BAOZ: boolFunc(
		shared_proto.HeroBoolType_BOOL_INVASE_BAOZ),
	shared_proto.TaskTargetType_TASK_TARGET_KILL_JUN_TUAN: historyAmountFunc(
		server_proto.HistoryAmountType_KillJunTuanTimes),
	shared_proto.TaskTargetType_TASK_TARGET_ADD_CAPTAIN_SOLDIER: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		if hero.Bools().Get(shared_proto.HeroBoolType_BOOL_RECOVER_SOLDIER) {
			newProgress = t.target.TotalProgress
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_BOOL: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		if hero.Bools().Get(shared_proto.HeroBoolType(t.target.SubType)) {
			newProgress = t.target.TotalProgress
		}
		return
	},
	shared_proto.TaskTargetType_TASK_TARGET_MCWAR_JOIN: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_McWarJoin)
	},
	shared_proto.TaskTargetType_TASK_TARGET_MCWAR_OCCUPY: func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmount(t.target.Daily, server_proto.HistoryAmountType_McOccupy)
	},
}

func boolFunc(boolType shared_proto.HeroBoolType) func(hero *Hero, t *TaskProgress) (newProgress uint64) {
	return func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		if hero.Bools().Get(boolType) {
			return t.target.TotalProgress
		}
		return
	}
}

func historyAmountFunc(hisAmtType server_proto.HistoryAmountType) func(hero *Hero, t *TaskProgress) (newProgress uint64) {
	return func(hero *Hero, t *TaskProgress) (newProgress uint64) {
		return hero.HistoryAmount().GetAmountWithSubType(t.target.Daily, hisAmtType, t.target.SubType)
	}
}

func RegisterUpdateTaskProgressFunc(targetType shared_proto.TaskTargetType, updateFunc func(hero *Hero, t *TaskProgress) (newProgress uint64)) {
	if updateTaskProgressFuncs[targetType] != nil {
		logrus.WithField("type", targetType).WithField("stack", string(debug.Stack())).Errorln("同一个任务类型被注册了多次")
		return
	}

	updateTaskProgressFuncs[targetType] = updateFunc
}
