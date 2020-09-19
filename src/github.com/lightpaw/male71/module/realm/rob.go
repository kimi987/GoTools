package realm

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/maildata"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"math"
	"time"
	"github.com/lightpaw/male7/gamelogs"
	"github.com/lightpaw/male7/entity/heroid"
	"github.com/lightpaw/male7/constants"
)

// 具体抢劫. 计算抢劫的资源量和损失的繁荣度
var ()

// 抢劫结果
type robResult struct {
	gold, food, wood, stone       uint64 // 被抢走的资源. 资源已经给robber的troop加上了
	baowu                         *resdata.BaowuData
	lostProsperity                uint64 // 被打掉的繁荣度
	targetDestroyed               bool   // 被抢劫的base已经爆了
	robberIsFull                  bool   // 抢劫的人的负载已经满了
	targetWontLoseAnyProsperity   bool   // 被攻击方不会损失繁荣度
	targetWontLoseAnyMoreResource bool   // 被攻击方之后不会再损失资源了, 已经全部到了保护值了
	hasError                      bool
}

// 刚到对方的base, 抢了第一下, 如果有抢到资源, 则已经加到了troop中
// 并没有保存robber的troop到它的英雄中, 等调用方决定这个troop下一步行动时, 一起保存
func (r *Realm) initialRob(target *baseWithData, robber *troop, ctime time.Time) (result robResult) {

	// 只有主城存在第一次被抢的情况
	switch target.BaseType() {
	case realmface.BaseTypeHome:
		result = r.doRobHeroHome(target, robber, true, false, false, true, ctime)
	case realmface.BaseTypeNpc:
		if robber.startingBase.isNpcBase() {
			result.hasError = true
			logrus.Error("抢第一下的时候，Npc抢Npc?")
			return
		}

		info := target.internalBase.getBaseInfoByLevel(robber.targetBaseLevel)
		data := info.getNpcBaseData()
		if data != nil {

			var hateTypeData *regdata.RegionMultiLevelNpcTypeData
			toAddHate := uint64(0)
			if hateData := info.getHateData(); hateData != nil {
				hateTypeData = hateData.TypeData()
				toAddHate = hateData.FirstHate
			}

			toAddPrize := data.FirstPrize
			if data.FirstConditionPlunder != nil {

				var heroLevel uint64
				if hero := r.services.heroSnapshotService.Get(robber.startingBase.Id()); hero != nil {
					heroLevel = hero.Level
				}

				toAddPrize = data.FirstConditionPlunder.GetPrize(heroLevel)
			} else if data.FirstPlunder != nil {
				toAddPrize = resdata.AppendPrize(toAddPrize, data.FirstPlunder.Try())
			}

			result = r.doRobNpc(target, data, robber, data.FirstLoseProsperity, true, toAddPrize, hateTypeData, toAddHate, ctime)
		} else {
			result.hasError = true
			logrus.Error("抢第一下的时候，转换NpcBase类型出错")
		}
	default:
		result.hasError = true
		logrus.Error("抢第一下的时候，无效的目标类型")
	}

	return
}

func (r *Realm) recurringRob(target *baseWithData, robber *troop, addPrize, reduceProsperity, robBaowu, addHate bool, ctime time.Time) (result robResult) {
	switch target.BaseType() {
	case realmface.BaseTypeHome:
		result = r.recurringRobHome(target, robber, addPrize, reduceProsperity, robBaowu, ctime)
	case realmface.BaseTypeNpc:
		if robber.startingBase.isNpcBase() {
			result.hasError = true
			logrus.Error("持续掠夺，Npc抢Npc?")
			return
		}
		result = r.recurringRobNpc(target, robber, addPrize, reduceProsperity, addHate, ctime)
	default:
		result.hasError = true
		logrus.WithField("base_type", target.BaseType()).Error("Realm.recurringRob Unkown BaseType")
	}

	return
}

func (r *Realm) doRobHeroHome(target *baseWithData, robber *troop, addPrize, reduceProsperity, robBaowu, isFirstRob bool, ctime time.Time) (result robResult) {
	if robber.startingBase.isNpcBase() {
		return r.doNpcRobHeroHome(target, robber, addPrize, reduceProsperity, isFirstRob, ctime)
	}

	// 抢的是老家
	var actProsperityCapcity uint64
	var actSoldierLevel uint64
	r.heroBaseFunc(robber.startingBase.Base(), func(hero *entity.Hero, err error) (heroChanged bool) {
		if err != nil {
			logrus.WithError(err).Error("realm.doRob时, lock英雄出错")
			result.hasError = true
			return
		}

		actProsperityCapcity = hero.ProsperityCapcity()
		actSoldierLevel = hero.Military().SoldierLevel()
		return
	})
	if result.hasError {
		return
	}

	originLevel := target.BaseLevel()

	var toReduceProsperity uint64
	var stopLostProsperityMsg pbutil.Buffer
	var toSendBeenRobMailFunc func()
	var robBaowuTroopId int64
	r.heroBaseFuncWithSend(target.Base(), func(hero *entity.Hero, heroResult herolock.LockResult) {

		// 计算能抢劫的资源量
		defProsperityCapcity := hero.ProsperityCapcity()

		// 衰减系数= arctan（（M守－M攻）/（M守＋M攻））/ π ＋1    π为圆周率
		weakCoef := math.Atan2(
			u64.Sub2Float64(defProsperityCapcity, actProsperityCapcity),
			float64(defProsperityCapcity+actProsperityCapcity),
		) / math.Pi + 1

		robberSoldier := robber.TotalAliveSoldier()

		if addPrize {
			// 收税收
			heromodule.TryUpdateTax(hero, heroResult, ctime, r.services.datas.MiscGenConfig().TaxDuration, r.dep.Datas().GetBuffEffectData)

			heroRes := hero.GetUnsafeResource()
			result.targetWontLoseAnyMoreResource = allProtected(heroRes.Gold(), heroRes.Food(), heroRes.Wood(), heroRes.Stone(), hero.ProtectedCapcity())

			if !result.targetWontLoseAnyMoreResource {

				transformResource := CalculateRecurringRobTransformResource
				if isFirstRob {
					transformResource = CalculateFirstRobTransformResource
				}

				// 还可以抢资源
				lostGold := u64.Min(u64.Sub(heroRes.Gold(), hero.ProtectedCapcity()),
					transformResource(robberSoldier, actSoldierLevel, heroRes.Gold(), hero.ProtectedCapcity(), weakCoef))
				lostFood := u64.Min(u64.Sub(heroRes.Food(), hero.ProtectedCapcity()),
					transformResource(robberSoldier, actSoldierLevel, heroRes.Food(), hero.ProtectedCapcity(), weakCoef))
				lostWood := u64.Min(u64.Sub(heroRes.Wood(), hero.ProtectedCapcity()),
					transformResource(robberSoldier, actSoldierLevel, heroRes.Wood(), hero.ProtectedCapcity(), weakCoef))
				lostStone := u64.Min(u64.Sub(heroRes.Stone(), hero.ProtectedCapcity()),
					transformResource(robberSoldier, actSoldierLevel, heroRes.Stone(), hero.ProtectedCapcity(), weakCoef))
				lostLoad := lostGold + lostFood + lostWood + lostStone

				if lostLoad > 0 {
					coef := math.Min(r.config().RobberCoef, 1)

					toAddGold := u64.MultiF64(lostGold, coef)
					toAddFood := u64.MultiF64(lostFood, coef)
					toAddWood := u64.MultiF64(lostWood, coef)
					toAddStone := u64.MultiF64(lostStone, coef)
					result.robberIsFull, result.gold, result.food, result.wood, result.stone = robber.tryAddResource(toAddGold, toAddFood, toAddWood, toAddStone)

					_, targetMsgFunc := heromodule.ReduceStorageResource(heroRes, lostGold, lostFood, lostWood, lostStone)
					heroResult.AddFunc(targetMsgFunc)
				}
			}
		}

		// 检查是否还可以损失繁荣度, 设置result
		if reduceProsperity {

			// 繁荣度不再设置上限
			//lostProsperity, stopLostProsperity := hero.LostProsperity(ctime, r.services.datas.MiscConfig().DailyResetTime)
			//if !stopLostProsperity {
			//	heroMaxLostProsperity := hero.MaxLostProsperity(r.config().GetMaxLostProsperity)
			//	canReduceProsperity := u64.Sub(heroMaxLostProsperity, lostProsperity)
			//	result.targetWontLoseAnyProsperity = canReduceProsperity <= 0
			//	if canReduceProsperity <= 0 {
			//		// 今日已达上限
			//		stopLostProsperity = true
			//		hero.SetStopLostProsperity()
			//	} else {
			//		//if reduceProsperity {

			// 繁荣度损失=INT（攻方剩余士兵数量/100*衰减系数）
			toReduce := u64.Max(1, u64.Multi(robberSoldier/100, weakCoef))
			//toReduce = u64.Min(toReduce, canReduceProsperity)
			toReduce = u64.Min(toReduce, hero.Prosperity())
			toReduceProsperity = toReduce

			logrus.WithField("weakcoef", weakCoef).
				WithField("attackerSoldier", robberSoldier).
				WithField("toReduce", toReduce).Debugf("扣繁荣度")

			if toReduce > 0 {
				robber.accumReduceProsperity += toReduce

				newProsperity := target.ReduceProsperity(toReduce)
				newBaseLevel := calculateBaseLevelByProsperity(newProsperity, r.services.datas.BaseLevelData().Must(originLevel))

				// 根据繁荣度计算等级，如果有变化，更新
				if originLevel != newBaseLevel {
					target.SetBaseLevel(newBaseLevel)

					if !npcid.IsNpcId(target.Id()) {
						gamelogs.UpdateBaseLevelLog(constants.PID, heroid.GetSid(target.Id()), target.Id(), target.BaseLevel())
					}
				}

				hero.UpdateBase(target.Base())

				//stopLostProsperity = hero.AddLostProsperity(toReduce, heroMaxLostProsperity)
				heroResult.AddFunc(func() pbutil.Buffer {
					return domestic.NewS2cHeroUpdateProsperityMsg(u64.Int32(newProsperity), u64.Int32(defProsperityCapcity))
				})

				heroResult.Changed()

				result.lostProsperity = toReduce
				result.targetDestroyed = newBaseLevel <= 0
			}

			//	}
			//}
			//
			//if stopLostProsperity {
			//	if target.TrySetStopLostProsperity(true) {
			//		stopLostProsperityMsg = region.NewS2cUpdateStopLostProsperityMsg(target.IdBytes()).Static()
			//	}
			//}

			if robBaowu {
				// 抢夺宝物
				var lostBaowu *resdata.BaowuData
				hero.Depot().RangeBaowu(func(id, count uint64) (toContinue bool) {
					data := r.services.datas.GetBaowuData(id)
					if data != nil && !data.CantRob {
						if lostBaowu == nil || resdata.BaowuLevelIsLarge(data, lostBaowu) {
							lostBaowu = data
						}
					}
					return true
				})

				if lostBaowu != nil {
					heromodule.ReduceBaowuAnyway(hero, heroResult, lostBaowu, 1)

					// 被抢日志
					hctx := heromodule.NewContext(r.dep, operate_type.BaoWuBeenRobbed)
					hctx.SetBaowuInfo(shared_proto.BaowuOpType_BOTBeenRob, robber.StartingBase().HeroName(), int32(robber.StartingBase().BaseX()), int32(robber.StartingBase().BaseY()), nil)
					heromodule.AddBaowuLog(hctx, hero, heroResult, lostBaowu, 1, ctime)

					robBaowuTroop := robber.NextRobBaowuTroop()
					robBaowuTroop.tryAddBaowu(lostBaowu, 1)
					robBaowuTroopId = robBaowuTroop.Id()
				}
				result.baowu = lostBaowu
			}
		}

		toSendBeenRobMailFunc = heromodule.TrySendMailFunc(r.services.mail, hero, shared_proto.HeroBoolType_BOOL_BEEN_ROB,
			r.services.datas.MailHelp().FirstBeenRob, ctime)

		heroResult.Ok()
	})

	if toSendBeenRobMailFunc != nil {
		toSendBeenRobMailFunc()
	}

	// 抢到的资源直接加给掠夺者
	totalRobResource := result.gold + result.food + result.wood + result.stone
	if totalRobResource > 0 || result.lostProsperity > 0 || result.baowu != nil {
		troopCount := robber.TroopCount()

		toAddGold, toAddFood, toAddWood, toAddStone := result.gold, result.food, result.wood, result.stone
		if troopCount > 1 {
			toAddGold = result.gold / troopCount
			toAddFood = result.food / troopCount
			toAddWood = result.wood / troopCount
			toAddStone = result.stone / troopCount
			totalRobResource = toAddGold + toAddFood + toAddWood + toAddStone
		}

		robber.walkAll(func(st *troop) (toContinue bool) {
			if totalRobResource > 0 || result.lostProsperity > 0 || (result.baowu != nil && st.Id() == robBaowuTroopId) {

				r.heroBaseFuncWithSend(st.startingBase.Base(), func(hero *entity.Hero, heroResult herolock.LockResult) {
					hctx := heromodule.NewContext(r.dep, operate_type.BaoWuRobHero)
					if totalRobResource > 0 {
						heromodule.AddUnsafeResource(hctx, hero, heroResult, toAddGold, toAddFood, toAddWood, toAddStone)

						hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumRobbingRes, uint64(shared_proto.ResType_GOLD), toAddGold)
						hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumRobbingRes, uint64(shared_proto.ResType_WOOD), toAddFood)
						hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumRobbingRes, uint64(shared_proto.ResType_FOOD), toAddWood)
						hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_AccumRobbingRes, uint64(shared_proto.ResType_STONE), toAddStone)
						hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AccumRobbingRes, totalRobResource)
						heromodule.UpdateTaskProgress(hero, heroResult, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_ROBBING_RES)
					}

					if result.targetDestroyed {
						hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumDestroyBase)
						heromodule.UpdateTaskProgress(hero, heroResult, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_DESTROY_BASE)
					}

					if result.lostProsperity > 0 {
						hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AccumDestroyProsperity, result.lostProsperity)
						heromodule.UpdateTaskProgress(hero, heroResult, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_DESTROY_PROSPERITY)
					}

					if result.baowu != nil && st.Id() == robBaowuTroopId {
						hctx := heromodule.NewContext(r.dep, operate_type.BaoWuRobHero)
						hctx.SetBaowuInfo(shared_proto.BaowuOpType_BOTRobHero, target.Base().HeroName(), int32(target.BaseX()), int32(target.BaseY()), nil)

						heromodule.AddBaowu(hctx, hero, heroResult, result.baowu, 1, ctime)

						hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_RobBaowu, result.baowu.Level)
						hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RobBaowu)
						heromodule.UpdateTaskProgress(hero, heroResult, shared_proto.TaskTargetType_TASK_TARGET_ROB_BAOWU)
					}

					heroResult.Ok()
				})
			}

			return true
		})
	}

	// 更新base繁荣度
	diffLevel := originLevel != target.BaseLevel()
	if toReduceProsperity > 0 {
		r.updateWatchBaseProsperityDiffLevel(target, diffLevel)

		target.updateRoBase()
	}

	if !diffLevel && stopLostProsperityMsg != nil {
		r.broadcastBaseToCared(target, stopLostProsperityMsg, 0)
	}

	// 广播飘字
	r.broadcastShowWords(target, robber, toReduceProsperity, result.gold, result.food, result.wood, result.stone)
	return
}

func (r *Realm) doNpcRobHeroHome(target *baseWithData, robber *troop, addPrize, reduceProsperity, isFirstRob bool, ctime time.Time) (result robResult) {

	originLevel := target.BaseLevel()

	var targetProsperityMsgFunc func() pbutil.Buffer
	var toReduceProsperity uint64
	var stopLostProsperityMsg pbutil.Buffer
	if reduceProsperity {
		// 扣繁荣度
		toReduce := r.config().NpcRobLostProsperityPerDuration
		if toReduce > 0 {
			r.heroBaseFunc(target.Base(), func(hero *entity.Hero, err error) (heroChanged bool) {
				if err != nil {
					logrus.WithError(err).Error("realm.doRob时, lock英雄出错")
					result.hasError = true
					return
				}

				// 检查是否还可以损失繁荣度, 设置result
				//lostProsperity, stopLostProsperity := hero.LostProsperity(ctime, r.services.datas.MiscConfig().DailyResetTime)
				//if !stopLostProsperity {
				//	heroMaxLostProsperity := hero.MaxLostProsperity(r.config().GetMaxLostProsperity)
				//	canReduceProsperity := u64.Sub(heroMaxLostProsperity, lostProsperity)
				//	if canReduceProsperity <= 0 {
				//		// 今日已达上限
				//		stopLostProsperity = true
				//		hero.SetStopLostProsperity()
				//	} else {
				//		toReduce = u64.Min(toReduce, canReduceProsperity)
				toReduce = u64.Min(toReduce, hero.Prosperity())
				toReduceProsperity = toReduce

				if toReduce > 0 {
					robber.accumReduceProsperity += toReduce

					newProsperity := target.ReduceProsperity(toReduce)
					newBaseLevel := calculateBaseLevelByProsperity(newProsperity, r.services.datas.BaseLevelData().Must(originLevel))

					// 根据繁荣度计算等级，如果有变化，更新
					if originLevel != newBaseLevel {
						target.SetBaseLevel(newBaseLevel)

						if !npcid.IsNpcId(target.Id()) {
							gamelogs.UpdateBaseLevelLog(constants.PID, heroid.GetSid(target.Id()), target.Id(), target.BaseLevel())
						}
					}

					hero.UpdateBase(target.Base())

					//stopLostProsperity = hero.AddLostProsperity(toReduce, heroMaxLostProsperity)
					defProsperityCapcity := hero.ProsperityCapcity()
					targetProsperityMsgFunc = func() pbutil.Buffer {
						return domestic.NewS2cHeroUpdateProsperityMsg(u64.Int32(newProsperity), u64.Int32(defProsperityCapcity))
					}

					heroChanged = true

					result.lostProsperity = toReduce
					result.targetDestroyed = newBaseLevel <= 0
				}

				//	}
				//}
				//
				//if stopLostProsperity {
				//	if target.TrySetStopLostProsperity(true) {
				//		stopLostProsperityMsg = region.NewS2cUpdateStopLostProsperityMsg(target.IdBytes()).Static()
				//	}
				//}

				return
			})

			if toReduce > 0 {
				target.updateRoBase()
			}
		}
	}

	if targetProsperityMsgFunc != nil {
		r.services.world.SendFunc(target.Id(), targetProsperityMsgFunc)
	}

	// 更新base繁荣度
	diffLevel := originLevel != target.BaseLevel()
	if toReduceProsperity > 0 {
		r.updateWatchBaseProsperityDiffLevel(target, diffLevel)

		target.updateRoBase()
	}

	if !diffLevel && stopLostProsperityMsg != nil {
		r.broadcastBaseToCared(target, stopLostProsperityMsg, 0)
	}

	// 广播飘字
	r.broadcastShowWords(target, robber, toReduceProsperity, result.gold, result.food, result.wood, result.stone)
	return
}

func (r *Realm) updateWatchBaseProsperity(base *baseWithData) {
	r.updateWatchBaseProsperityDiffLevel(base, false)
}

func (r *Realm) updateWatchBaseProsperityDiffLevel(base *baseWithData, diffLevel bool) {
	baseId := base.Id()
	toSend := region.NewS2cUpdateWatchBaseProsperityMsg(base.IdBytes(), u64.Int32(base.Prosperity()), u64.Int32(base.ProsperityCapcity())).Static()
	r.services.world.WalkHero(func(id int64, hc iface.HeroController) {
		if cond := hc.GetCareCondition(); cond != nil && cond.BaseId == baseId {
			hc.Send(toSend)
		} else if !diffLevel && base.CanSee(hc.GetViewArea()) {
			hc.Send(toSend)
		}
	})
}

//func (r *Realm) doRobHeroTent(target *baseWithData, robber *troop, ctime time.Time) (result robResult) {
//
//	// 抢的是行营
//	originProsperity := target.Prosperity()
//
//	// 每6秒钟，繁荣度损失=INT（攻方剩余士兵数量/200）
//	toReduce := u64.Max(1, robber.AliveSoldier()/200)
//	toReduce = u64.Min(toReduce, target.Prosperity())
//
//	// 目标扣这么多
//	target.ReduceProsperity(toReduce)
//	result.targetDestroyed = target.Prosperity() <= 0 // 被打爆了
//
//	var updateTentProsperityMsgFunc func() pbutil.Buffer
//	r.heroBaseFuncNotError(target.Base(), func(hero *entity.Hero) (heroChanged bool) {
//		hero.UpdateBase(target.Base())
//
//		if result.targetDestroyed {
//			hero.RemoveTent()
//		}
//
//		updateTentProsperityMsgFunc = func() pbutil.Buffer {
//			return region.NewS2cUpdateTentProsperityMsg(u64.Int32(target.Prosperity()))
//		}
//		return true
//	})
//	r.services.world.SendFunc(target.Id(), updateTentProsperityMsgFunc)
//
//	// 没挂，更新基地繁荣度变化
//	r.updateBaseProsperityIfAlive(target, originProsperity)
//	return
//}
//
//func (r *Realm) doRobHeroHome(target *baseWithData, robber *troop, reduceProsperity bool,
//	ctime time.Time) (result robResult) {
//
//	// 扣老家繁荣度
//	// 抢玉石矿
//
//	// 抢的是老家
//	var actProsperityCapcity uint64
//	r.heroBaseFunc(robber.startingBase.Base(), func(hero *entity.Hero, err error) (heroChanged bool) {
//		if err != nil {
//			logrus.WithError(err).Error("realm.doRob时, lock英雄出错")
//			result.hasError = true
//			return
//		}
//
//		actProsperityCapcity = hero.ProsperityCapcity()
//		return
//	})
//	if result.hasError {
//		return
//	}
//
//	originLevel := target.BaseLevel()
//	originPrisperity := target.Prosperity()
//
//	var toAddJadeOre uint64
//	var toReduceProsperity uint64
//	var updateTargetBaseJadeOreFunc func()
//	result.hasError = r.heroBaseFuncWithSend(target.Base(), func(hero *entity.Hero, heroResult herolock.LockResult) {
//
//		// 计算能抢劫的资源量
//		defProsperityCapcity := hero.ProsperityCapcity()
//
//		// 衰减系数= arctan（（M守－M攻）/（M守＋M攻））/ π ＋1    π为圆周率
//		weakCoef := math.Atan2(
//			u64.Sub2Float64(defProsperityCapcity, actProsperityCapcity),
//			float64(defProsperityCapcity+actProsperityCapcity),
//		)/math.Pi + 1
//
//		// TODO 临时，每次抢千分之一，最少抢1点
//		if hero.JadeOre() > 0 {
//			toReduce := hero.JadeOre() / 1000
//			toReduce = u64.Max(1, toReduce)
//
//			newJadeOre := hero.ReduceJadeOre(toReduce)
//			heroResult.AddFunc(func() pbutil.Buffer {
//				return domestic.NewS2cUpdateJadeOreMsg(u64.Int32(newJadeOre))
//			})
//
//			toAddJadeOre = toReduce
//
//			heroResult.Changed()
//
//			updateTargetBaseJadeOreFunc = r.newUpdateBaseJadeFunc(target, hero.OtherBaseRegion(r.id), newJadeOre)
//		} else {
//			result.targetWontLoseAnyMoreResource = true
//		}
//
//		// 检查是否还可以损失繁荣度, 设置result
//		canReduceProsperity := u64.Sub(hero.MaxLostProsperity(r.config.GetMaxLostProsperity), hero.LostProsperity(ctime))
//		result.targetWontLoseAnyProsperity = canReduceProsperity <= 0
//		logrus.Debugf("canReduceProsperity: %d", canReduceProsperity)
//		if reduceProsperity && canReduceProsperity > 0 {
//
//			attackerTroopsSoldier := robber.AliveSoldier()
//
//			// 繁荣度损失=INT（攻方剩余士兵数量/3*衰减系数）
//			toReduce := u64.Max(1, u64.Multi(attackerTroopsSoldier/3, weakCoef))
//			toReduce = u64.Min(toReduce, canReduceProsperity)
//			toReduce = u64.Min(toReduce, hero.Prosperity())
//			toReduceProsperity = toReduce
//
//			logrus.WithField("weakcoef", weakCoef).
//				WithField("attackerSoldier", attackerTroopsSoldier).
//				WithField("toReduce", toReduce).Debugf("扣繁荣度")
//
//			if toReduce > 0 {
//				newProsperity := target.ReduceProsperity(toReduce)
//				newBaseLevel := calculateBaseLevelByProsperity(newProsperity, r.services.datas.BaseLevelData().Must(originLevel))
//
//				// 根据繁荣度计算等级，如果有变化，更新
//				if originLevel != newBaseLevel {
//					target.SetBaseLevel(newBaseLevel)
//				}
//
//				hero.UpdateBase(target.Base())
//
//				hero.AddLostProsperity(toReduce)
//				heroResult.AddFunc(func() pbutil.Buffer {
//					return domestic.NewS2cHeroUpdateProsperityMsg(u64.Int32(newProsperity), u64.Int32(defProsperityCapcity))
//				})
//
//				heroResult.Changed()
//
//				result.lostProsperity = toReduce
//				result.targetDestroyed = newBaseLevel <= 0
//
//				if result.targetDestroyed {
//					hero.ClearBase()
//				}
//			}
//
//		}
//
//		heroResult.Ok()
//		return
//	})
//	if result.hasError {
//		return
//	}
//
//	if updateTargetBaseJadeOreFunc != nil {
//		updateTargetBaseJadeOreFunc()
//	}
//
//	if originLevel == target.BaseLevel() {
//		// 没挂，更新基地繁荣度变化
//		r.updateBaseProsperityIfAlive(target, originPrisperity)
//	}
//
//	// 添加抢到的玉石矿
//	if toAddJadeOre > 0 {
//		var updateRobberBaseJadeOreFunc func()
//		r.heroBaseFuncWithSend(robber.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
//			newJadeOre := hero.AddJadeOre(toAddJadeOre)
//			result.Add(domestic.NewS2cUpdateJadeOreMsg(u64.Int32(newJadeOre)))
//
//			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_RobJadeOre, toAddJadeOre)
//			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_JADE_ORE)
//
//			updateRobberBaseJadeOreFunc = r.newUpdateBaseJadeFunc(robber.startingBase, hero.OtherBaseRegion(r.id), newJadeOre)
//
//			result.Changed()
//			result.Ok()
//			return
//		})
//
//		if updateRobberBaseJadeOreFunc != nil {
//			updateRobberBaseJadeOreFunc()
//		}
//
//		result.jadeOre += toAddJadeOre
//		robber.tryAddJadeOre(toAddJadeOre)
//
//		r.broadcastToCared(region.NewS2cRobJadeOreChangedMsg(robber.IdBytes(), u64.Int32(robber.JadeOre())), 0)
//	}
//
//	// 广播飘字
//	r.broadcastShowWords(target, robber, toReduceProsperity, 0, 0, 0, 0, 0, toAddJadeOre)
//	return
//
//}

func (r *Realm) broadcastShowWords(target *baseWithData, robber *troop, prosperity, gold, food, wood, stone uint64) {

	var robberIdBytes []byte
	if robber != nil {
		robberIdBytes = robber.IdBytes()
	}

	r.broadcastBaseToCared(target, region.NewS2cShowWordsMsg(
		target.IdBytes(), robberIdBytes, u64.Int32(prosperity),
		u64.Int32(gold), u64.Int32(food), u64.Int32(wood),
		u64.Int32(stone), 0, 0,
	), 0)

}

func (r *Realm) doRobNpc(target *baseWithData, data *basedata.NpcBaseData, robber *troop,
	reduceProsperity uint64, isFirstRob bool, toAddPrize *resdata.Prize, hateTypeData *regdata.RegionMultiLevelNpcTypeData,
	toAddHate uint64, ctime time.Time) (robResult robResult) {
	// 抢的是npc基地

	var toReduceProsperity uint64
	if reduceProsperity > 0 {
		// 每次扣除繁荣度
		toReduceProsperity = u64.Min(reduceProsperity, target.Prosperity())

		robber.accumReduceProsperity += toReduceProsperity
		// 目标扣这么多
		target.ReduceProsperity(toReduceProsperity)
		robResult.targetDestroyed = target.Prosperity() <= 0 // 被打爆了

		// 更新base繁荣度
		r.updateWatchBaseProsperity(target)

		// 宝藏怪物清除缓存
		if b := GetBaoZangBase(target); b != nil {
			b.ClearUpdateBaseInfoMsg()
		}

		target.updateRoBase()
	}

	// 我加奖励
	if toAddPrize != nil {
		var rareBaowuBaseName string
		var rareBaowuIds []uint64
		if b := GetBaoZangBase(target); b != nil {
			rareBaowuBaseName = b.data.Npc.Name
			rareBaowuIds = b.data.RareBaowuIds
		}

		robber.walkAll(func(st *troop) (toContinue bool) {
			r.heroBaseFuncWithSend(st.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {

				// 首次掠夺Npc宝物
				if b := GetBaoZangBase(target); b != nil {
					if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_NPC_BAOWU) {
						result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_NPC_BAOWU)))

						// 奖励中添加一个npc宝物
						if toAdd := r.services.datas.MiscGenConfig().FirstNpcBaowu; toAdd != nil {
							toAddPrize = resdata.NewPrizeBuilder().
								Add(toAddPrize).
								AddBaowu(toAdd, 1).
								Build()
						}
					}
				}

				var captainIds []uint64
				for _, c := range st.captains {
					captainIds = append(captainIds, c.Id())
				}

				// 讨伐野怪可以多次
				toAddPrize = toAddPrize.Multiple(u64.Max(1, robber.npcTimes))

				hctx := heromodule.NewContext(r.dep, operate_type.RealmRobNpc)
				hctx.SetBaowuInfo(shared_proto.BaowuOpType_BOTRobNpc, rareBaowuBaseName, int32(target.BaseX()), int32(target.BaseY()), rareBaowuIds)
				heromodule.AddCaptainPrize(hctx, hero, result, toAddPrize, captainIds, ctime)

				// 加仇恨
				if toAddHate > 0 && hateTypeData != nil {
					heroTroop := hero.GetTroopByIndex(entity.GetTroopIndex(st.Id()))
					if heroTroop != nil {
						if ii := heroTroop.GetInvateInfo(); ii != nil {
							ii.AddHate(toAddHate)
						}
					}
				}

				if robResult.targetDestroyed {
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumDestroyBase)
					heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_DESTROY_BASE)
				}
				hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AccumDestroyProsperity, toReduceProsperity)
				heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_DESTROY_PROSPERITY)

				// 掠夺野怪资源
				if npcid.IsMultiLevelMonsterNpcId(target.Id()) {

					totalRes := toAddPrize.Gold + toAddPrize.Food + toAddPrize.Wood + toAddPrize.Stone
					if totalRes > 0 {
						hero.HistoryAmount().IncreaseWithSubType(
							server_proto.HistoryAmountType_RobMultiLevelMonsterRes,
							uint64(shared_proto.ResType_GOLD),
							toAddPrize.Gold)

						hero.HistoryAmount().IncreaseWithSubType(
							server_proto.HistoryAmountType_RobMultiLevelMonsterRes,
							uint64(shared_proto.ResType_FOOD),
							toAddPrize.Food)

						hero.HistoryAmount().IncreaseWithSubType(
							server_proto.HistoryAmountType_RobMultiLevelMonsterRes,
							uint64(shared_proto.ResType_WOOD),
							toAddPrize.Wood)

						hero.HistoryAmount().IncreaseWithSubType(
							server_proto.HistoryAmountType_RobMultiLevelMonsterRes,
							uint64(shared_proto.ResType_STONE),
							toAddPrize.Stone)

						hero.HistoryAmount().IncreaseWithSubType(
							server_proto.HistoryAmountType_RobMultiLevelMonsterRes, 0, totalRes)

						heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ROB_MULTI_LEVEL_MONSTER)
					}
				}

				if n := imath.Min(len(toAddPrize.Baowu), len(toAddPrize.BaowuCount)); n > 0 {
					var totalCount uint64
					for i := 0; i < n; i++ {
						baowu := toAddPrize.Baowu[i]
						count := toAddPrize.BaowuCount[i]
						if baowu != nil && count > 0 {
							hero.HistoryAmount().IncreaseWithSubType(server_proto.HistoryAmountType_RobNpcBaowu, baowu.Level, count)
							totalCount += count
						}
					}

					if totalCount > 0 {
						hero.HistoryAmount().Increase(server_proto.HistoryAmountType_RobNpcBaowu, totalCount)
						heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ROB_NPC_BAOWU)
					}
				}

				if robResult.targetDestroyed {
					if b := GetBaoZangBase(target); b != nil {
						heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_KILL_BAOZ, b.data.Npc.BaseLevel)
					}
				}

				result.Changed()
				result.Ok()
				return
			})

			st.tryAddPrize(toAddPrize)

			return true
		})
	}

	// TODO 繁荣度改变，宝藏Npc需要更新

	// 广播飘字
	if toAddPrize != nil {
		r.broadcastShowWords(target, robber, toReduceProsperity, toAddPrize.Gold, toAddPrize.Food,
			toAddPrize.Wood, toAddPrize.Stone)
	} else {
		r.broadcastShowWords(target, robber, toReduceProsperity, 0, 0, 0, 0)
	}
	return
}

var (
	homeBrokenMsg = region.NewS2cSelfBaseDestroyMsg(false, int32(realmface.BDTBroken)).Static()
)

func (r *Realm) recurringRobHome(target *baseWithData, robber *troop, addPrize, reduceProsperity, robBaowu bool, ctime time.Time) (result robResult) {

	originLevel := target.BaseLevel()
	result = r.doRobHeroHome(target, robber, addPrize, reduceProsperity, robBaowu, false, ctime)

	// 降级邮件
	if originLevel != target.BaseLevel() {
		//r.sendLevelDownMail(target, originLevel, ctime)

		// 移除占用的资源点
		r.updateHomeResourcePointBlockWhenLevelChanged(target, originLevel)

		if target.BaseLevel() > 0 {
			// 发给自己
			r.services.world.SendFunc(target.Id(), func() pbutil.Buffer {
				return region.NewS2cSelfUpdateBaseLevelMsg(u64.Int32(target.BaseLevel()))
			})

			// 降级
			r.broadcastBaseInfoToCared(target, addBaseTypeUpdate, 0)
		} else {
			// 流亡
			r.removeRealmBase(target, r.services.datas.TextHelp().MRDRBroken4a.Text, r.services.datas.TextHelp().MRDRBroken4d.Text, removeBaseTypeBroken, ctime)

			// 发给自己
			r.services.world.Send(target.Id(), homeBrokenMsg)
			r.services.world.SendFunc(target.Id(), func() pbutil.Buffer {
				return region.NewS2cSelfUpdateBaseLevelMsg(u64.Int32(0))
			})

			// 产生废墟
			r.ruinsBasePosInfoMap.AddRuinsBase(r, target.Id(), target.BaseX(), target.BaseY(), ctime)

			// 将玩家从联盟防守列表中移除
			guildSnapshot := r.services.guildService.GetSnapshot(target.GuildId())
			if guildSnapshot != nil && i64.Contains(guildSnapshot.ResistXiongNuDefenders, target.Id()) {
				// 怎么丢到联盟线程呢
				r.services.otherRealmEventQueue.MustFunc(func() {
					r.services.guildService.FuncGuild(target.GuildId(), func(g *sharedguilddata.Guild) {
						if g == nil {
							return
						}

						g.SetChanged()
						g.RemoveResistXiongNuDefender(target.Id())
					})
				})
			}

			attackerFlagName := robber.getStartingBaseFlagHeroName(r)
			targetFlagName := r.toBaseFlagHeroName(target, robber.targetBaseLevel)

			// 把别人打流亡
			if gid := robber.startingBase.GuildId(); gid != 0 {
				if data := r.services.datas.GuildLogHelp().FAttDestroy; data != nil {
					if hero := r.services.heroSnapshotService.Get(robber.startingBase.Id()); hero != nil {
						if hero.GuildId == gid {
							proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
							proto.Text = data.Text.New().WithAttacker(attackerFlagName).WithDefenser(targetFlagName).JsonString()
							proto.FightX = int32(target.BaseX())
							proto.FightY = int32(target.BaseY())

							r.services.guildService.AddLog(hero.GuildId, proto)
						}
					}
				}
			}

			// 被别人打流亡
			if gid := target.GuildId(); gid != 0 {
				if data := r.services.datas.GuildLogHelp().FDefDestroy; data != nil {
					if hero := r.services.heroSnapshotService.Get(target.Id()); hero != nil {
						if hero.GuildId == gid {
							proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
							proto.Text = data.Text.New().WithAttacker(attackerFlagName).WithDefenser(targetFlagName).JsonString()
							proto.FightX = int32(target.BaseX())
							proto.FightY = int32(target.BaseY())

							r.services.guildService.AddLog(hero.GuildId, proto)
						}
					}
				}
			}
		}
	}

	return result
}

func (r *Realm) recurringRobNpc(target *baseWithData, robber *troop, addPrize, reduceProsperity, addHate bool, ctime time.Time) (result robResult) {
	result.targetWontLoseAnyMoreResource = true // 外面不发更新消息

	info := target.internalBase.getBaseInfoByLevel(robber.targetBaseLevel)
	data := info.getNpcBaseData()
	if data == nil {
		logrus.WithField("id", target.Id()).WithField("base_name", target.internalBase.HeroName()).Error("Realm.recurringRobNpc 不是NpcBase")
		return
	}

	var hateTypeData *regdata.RegionMultiLevelNpcTypeData
	toAddHate := uint64(0)
	if addHate {
		if hateData := info.getHateData(); hateData != nil {
			hateTypeData = hateData.TypeData()
			toAddHate = hateData.TickHate
		}
	}

	var toAddPrize *resdata.Prize
	if addPrize {
		toAddPrize = data.TickPrize
		if data.TickConditionPlunder != nil {

			var heroLevel uint64
			if hero := r.services.heroSnapshotService.Get(robber.startingBase.Id()); hero != nil {
				heroLevel = hero.Level
			}

			toAddPrize = data.TickConditionPlunder.GetPrize(heroLevel)
		} else if data.TickPlunder != nil {
			toAddPrize = resdata.AppendPrize(toAddPrize, data.TickPlunder.Try())
		}
	}

	var toReduceProsperity uint64
	if reduceProsperity {
		toReduceProsperity = data.LostProsperityPerDuration
	}

	result = r.doRobNpc(target, data, robber, toReduceProsperity, false, toAddPrize, hateTypeData, toAddHate, ctime)
	if result.targetDestroyed {
		// 对方被打爆了
		r.removeRealmBase(target, r.services.datas.TextHelp().MRDRBroken4a.Text, r.services.datas.TextHelp().MRDRBroken4d.Text, removeBaseTypeBroken, ctime)
	}

	return
}

func allProtected(gold, food, wood, stone, protectedCapacity uint64) bool {
	return gold <= protectedCapacity && food <= protectedCapacity && wood <= protectedCapacity && stone <= protectedCapacity
}

func CalculateFirstRobTransformResource(soldier, soldierLevel, amount, protectedAmount uint64, weakCoef float64) uint64 {
	if soldier <= 0 || amount <= protectedAmount || weakCoef <= 0 {
		return 0
	}

	// 守方每种资源损失量 = INT(攻方剩余士兵数量*攻方士兵等级*15+12000) *衰减系数
	toAdd := u64.Multi(soldier*soldierLevel*15+12000, weakCoef)
	if toAdd >= 1000000 {
		logrus.WithField("soldier", soldier).
			WithField("amount", amount).
			WithField("protectedAmount", protectedAmount).
			WithField("weakCoef", weakCoef).
			Error("首次抢夺资源，超过百万，预警，可能有大bug")
	}

	return toAdd
}

func CalculateRecurringRobTransformResource(soldier, soldierLevel, amount, protectedAmount uint64, weakCoef float64) uint64 {
	if soldier <= 0 || soldierLevel <= 0 || amount <= protectedAmount || weakCoef <= 0 {
		return 0
	}

	// 守方每种资源损失量 = INT(攻方剩余士兵数量*攻方士兵等级*3/2+1200)*衰减系数
	toAdd := u64.Multi(soldier*soldierLevel*3/2+1200, weakCoef)
	if toAdd >= 1000000 {
		logrus.WithField("soldier", soldier).
			WithField("amount", amount).
			WithField("protectedAmount", protectedAmount).
			WithField("weakCoef", weakCoef).
			Error("持续抢夺资源，超过百万，预警，可能有大bug")
	}
	return toAdd
}

const (
	attackerKeyName = "attacker"
	defenserKeyName = "defenser"
	assisterKeyName = "assister"
)

func fieldsSetHeroFlagName(fields *i18n.Fields, attackerFlagName, defenserFlagName, assisterFlagName string) *i18n.Fields {

	if len(attackerFlagName) > 0 {
		fields.WithFields(attackerKeyName, attackerFlagName)
	}

	if len(defenserFlagName) > 0 {
		fields.WithFields(defenserKeyName, defenserFlagName)
	}

	if len(assisterFlagName) > 0 {
		fields.WithFields(assisterKeyName, assisterFlagName)
	}

	return fields
}

func (r *Realm) trySendTroopDoneMail(t *troop, reason4a, reason4d *i18n.I18nRef, ctime time.Time) {

	// 队伍状态为robbing的，才有掠夺结束邮件
	if t.state != realmface.Robbing {
		return
	}

	if t.targetBase == nil {
		logrus.Error("发送掠夺结束战报，但是 targetBase == nil")
		return
	}

	if reason4a == nil && reason4d == nil {
		logrus.Error("发送掠夺结束战报，但是 reason == nil")
		return
	}

	isAssembly := t.assembly != nil

	attackerId := t.startingBase.Id()
	attackerFlagName := t.getStartingBaseFlagHeroName(r)

	defenserId := t.targetBase.Id()
	defenserFlagName := t.getTargetBaseFlagHeroName(r)

	// 持续掠夺时间
	targetBaseInfo := t.targetBase.internalBase.getBaseInfoByLevel(t.targetBaseLevel)
	robDuration := t.getRobDuration(r, targetBaseInfo)
	if d := t.robbingEndTime.Sub(ctime); d > 0 && d < robDuration {
		robDuration -= d
	}

	report := &shared_proto.FightReportProto{
		Attacker: t.getStartingBaseReportProto(r),
		Defenser: t.getTargetBaseReportProto(r),

		FightX: int32(t.targetBase.BaseX()),
		FightY: int32(t.targetBase.BaseY()),

		RobDuration: timeutil.DurationMarshal32(robDuration),
	}

	if npcid.IsMultiLevelMonsterNpcId(defenserId) {
		// 加仇恨，根据持续掠夺时间，计算仇恨值
		report.DoneDesc = r.services.datas.TextHelp().MailAddHateDesc.Text.KeysOnlyJson()
		if hateData := targetBaseInfo.getHateData(); hateData != nil {
			toAddHate := hateData.FirstHate
			if hateData.HateTickDuration > 0 {
				toAddHate += uint64(robDuration/hateData.HateTickDuration) * hateData.TickHate
			}
			report.DoneDescAmount = u64.Int32(toAddHate)
		}
		if t.npcTimes > 1 {
			if d := r.dep.Datas().TextHelp().MultiMonsterPrizeCount; d != nil {
				report.Desc = d.Text.New().WithCount(t.npcTimes).JsonString()
			}

		}
	} else {
		// 扣繁荣度
		report.DoneDesc = r.services.datas.TextHelp().MailReduceProsperityDesc.Text.KeysOnlyJson()
		report.DoneDescAmount = u64.Int32(t.accumReduceProsperity)
	}

	mailType := shared_proto.MailType_MailRobFinished
	if isAssembly {
		mailType = shared_proto.MailType_MailAssemblyRobFinished

		report.AttackerTroopCount = int32(t.TroopCount())
		report.AttackerTroopTotalCount = int32(t.assembly.TotalCount())

		t.walkAll(func(st *troop) (toContinue bool) {
			hero := report.Attacker
			if st != t {
				hero = st.getStartingBaseReportProto(r)
			}

			fight := &shared_proto.AssemblyFightProto{
				Attacker:  hero,
				ShowPrize: st.AccumRobPrize(),
			}
			report.Fight = append(report.Fight, fight)
			return true
		})

		// 集结奖励，不在这里处理
	} else {
		report.ShowPrize = t.AccumRobPrize()
	}

	//gold, food, wood, stone := t.Carrying()

	if !npcid.IsNpcId(attackerId) && reason4a != nil {

		if bz := GetBaoZangBase(t.targetBase); bz != nil {

			if data := r.getMailHelp().ReportDoneBaoz; data != nil {
				baseName := bz.data.Npc.Name

				mailProto := data.NewTextMail(mailType)
				mailProto.SubTitle = fieldsSetHeroFlagName(data.NewSubTitleFields(), attackerFlagName, baseName, "").
					WithReason(reason4a).JsonString()
				mailProto.Text = fieldsSetHeroFlagName(data.NewTextFields(), attackerFlagName, baseName, "").
					WithReason(reason4a).JsonString()
				//mailProto.Text = fieldsSetHeroFlagName(data.NewTextFields(), attackerFlagName, defenserFlagName, "").
				//	WithFields("gold", gold).WithFields("food", food).
				//	WithFields("wood", wood).WithFields("stone", stone).JsonString()

				mailProto.Report = report
				report.AttackerSide = true
				report.Desc = data.Desc.KeysOnlyJson()
				report.DoneDescAmount = 0

				mailProto.ReportTag = getMailTag(attackerId, defenserId)

				t.walkAll(func(st *troop) (toContinue bool) {
					r.services.mail.SendReportMail(st.startingBase.Id(), mailProto, ctime)
					return true
				})
			}
		} else {
			if data := r.getMailHelp().ReportDoneAttacker; data != nil {
				mailProto := data.NewTextMail(mailType)
				mailProto.SubTitle = fieldsSetHeroFlagName(data.NewSubTitleFields(), attackerFlagName, defenserFlagName, "").
					WithReason(reason4a).JsonString()
				mailProto.Text = fieldsSetHeroFlagName(data.NewTextFields(), attackerFlagName, defenserFlagName, "").
					WithReason(reason4a).JsonString()
				//mailProto.Text = fieldsSetHeroFlagName(data.NewTextFields(), attackerFlagName, defenserFlagName, "").
				//	WithFields("gold", gold).WithFields("food", food).
				//	WithFields("wood", wood).WithFields("stone", stone).JsonString()

				mailProto.Report = report
				report.AttackerSide = true

				if t.npcTimes > 1 {
					if d := r.dep.Datas().TextHelp().MultiMonsterPrizeCount; d != nil {
						report.Desc = d.Text.New().WithCount(t.npcTimes).JsonString()
					}
				} else {
					report.Desc = data.Desc.KeysOnlyJson()
				}

				mailProto.ReportTag = getMailTag(attackerId, defenserId)

				t.walkAll(func(st *troop) (toContinue bool) {
					// 已经缓存，可以直接使用
					report.ShowPrize = st.AccumRobPrize()
					r.services.mail.SendReportMail(st.startingBase.Id(), mailProto, ctime)
					return true
				})
			}
		}
	}

	if !npcid.IsNpcId(defenserId) && reason4d != nil {
		if data := r.getMailHelp().ReportDoneDefenser; data != nil {
			mailProto := data.NewTextMail(mailType)
			mailProto.SubTitle = fieldsSetHeroFlagName(data.NewSubTitleFields(), attackerFlagName, defenserFlagName, "").
				WithReason(reason4a).JsonString()
			mailProto.Text = fieldsSetHeroFlagName(data.NewTextFields(), attackerFlagName, defenserFlagName, "").
				WithReason(reason4d).JsonString()
			//mailProto.Text = fieldsSetHeroFlagName(data.NewTextFields(), attackerFlagName, defenserFlagName, "").
			//	WithFields("gold", gold).WithFields("food", food).
			//	WithFields("wood", wood).WithFields("stone", stone).JsonString()

			mailProto.Report = report
			report.AttackerSide = false
			report.Desc = data.Desc.KeysOnlyJson()
			report.DoneDesc = r.services.datas.TextHelp().MailLostProsperityDesc.Text.KeysOnlyJson()

			if isAssembly {
				// 把这些人抢的都加起来才是损失的
				pb := resdata.NewPrizeBuilder()
				t.walkAll(func(st *troop) (toContinue bool) {
					pb.Add(st.BuildAccumRobPrize())
					return true
				})
				report.ShowPrize = pb.Build().Encode()
			}

			// 还原到之前的值
			if report.ShowPrize != nil {
				if coef := math.Min(r.config().RobberCoef, 1); coef > 0 && coef < 1 {
					coef = 1 / coef
					report.ShowPrize.Gold = i32.MultiCoef(report.ShowPrize.Gold, coef)
					report.ShowPrize.Food = i32.MultiCoef(report.ShowPrize.Food, coef)
					report.ShowPrize.Wood = i32.MultiCoef(report.ShowPrize.Wood, coef)
					report.ShowPrize.Stone = i32.MultiCoef(report.ShowPrize.Stone, coef)
				}
			}

			mailProto.ReportTag = getMailTag(attackerId, defenserId)
			r.services.mail.SendReportMail(defenserId, mailProto, ctime)
		}
	}
}

func getMailTagByTroop(attacker, defenser *troop) int32 {
	if attacker != nil {
		if defenser != nil {
			return getMailTag(attacker.startingBase.Id(), attacker.originTargetId, defenser.startingBase.Id())
		} else {
			return getMailTag(attacker.startingBase.Id(), attacker.originTargetId)
		}
	}
	return maildata.TagPvp
}

func getMailTag(ids ...int64) int32 {
	for _, id := range ids {
		if npcid.IsBaoZangNpcId(id) || npcid.IsMultiLevelMonsterNpcId(id) {
			return maildata.TagPve
		}

		if npcid.IsXiongNuNpcId(id) {
			return maildata.TagAct
		}
	}

	return maildata.TagPvp
}

func (r *Realm) sendHeroReportMail(toSendHeroId int64, data *maildata.MailData,
	attackerFlagName, defenserFlagName, assisterFlagName string,
	report *shared_proto.FightReportProto, isAttackerSide bool, score int32,
	ctime time.Time, reportTag int32, mailType shared_proto.MailType) {
	if toSendHeroId == 0 || npcid.IsNpcId(toSendHeroId) {
		return
	}

	r.sendHeroReportMailReturnProto(toSendHeroId, data,
		attackerFlagName, defenserFlagName, assisterFlagName,
		report, isAttackerSide, score, ctime, reportTag, mailType)
}

func (r *Realm) sendHeroReportMailReturnProto(toSendHeroId int64, data *maildata.MailData,
	attackerFlagName, defenserFlagName, assisterFlagName string,
	report *shared_proto.FightReportProto, isAttackerSide bool, score int32,
	ctime time.Time, reportTag int32, mailType shared_proto.MailType) *shared_proto.MailProto {

	mailProto := data.NewTextMail(mailType)
	// title 不带参数，这里不做处理
	mailProto.SubTitle = fieldsSetHeroFlagName(data.NewSubTitleFields(), attackerFlagName, defenserFlagName, assisterFlagName).JsonString()
	mailProto.Text = fieldsSetHeroFlagName(data.NewTextFields(), attackerFlagName, defenserFlagName, assisterFlagName).JsonString()

	if report != nil {
		mailProto.Report = report

		if report.Desc == "" {
			report.Desc = data.Desc.KeysOnlyJson()
		}

		report.AttackerSide = isAttackerSide
		if report.Share != nil {
			report.Share.IsAttacker = isAttackerSide
		}

		if len(report.Fight) > 0 {
			for _, f := range report.Fight {
				if f.Share != nil {
					f.Share.IsAttacker = isAttackerSide
				}
			}
		}
	}

	// 小胜
	mailProto.ImageWord += score

	mailProto.ReportTag = reportTag
	r.services.mail.SendReportMail(toSendHeroId, mailProto, ctime)

	return mailProto
}

func getMailData(report, assembly *maildata.MailData, isAssembly bool) *maildata.MailData {
	if isAssembly {
		return assembly
	}
	return report
}

func (r *Realm) sendInvadeSuccessMail(attackerTroop *troop, attackerFightTroop []*troop, targetBase *baseWithData, targetDefenser *troop, score int32, report *shared_proto.FightReportProto, ctime time.Time, isBack bool) {

	if report == nil {
		return
	}
	//report = r.newNoFightReport(attackerTroop, targetBase, i18n.MailKey.MailReportWinnerDesc().KeysOnlyJson(), i18n.MailKey.MailReportLoserDesc().KeysOnlyJson())

	attackerId := attackerTroop.startingBase.Id()
	targetId := targetBase.Id()
	attackerFlagName := attackerTroop.getStartingBaseFlagHeroName(r)
	targetFlagName := r.toBaseFlagHeroName(targetBase, attackerTroop.targetBaseLevel)

	isAssemblyMail := len(attackerFightTroop) > 1

	// 战斗胜利邮件
	var showPrizeProtos []*shared_proto.PrizeProto
	if data := getMailData(
		r.getMailHelp().ReportAdaSuccess,
		r.getMailHelp().AssemblyAdaSuccess,
		isAssemblyMail,
	); data != nil {
		mailTag := getMailTagByTroop(attackerTroop, targetDefenser)

		for _, t := range attackerFightTroop {
			report.ShowPrize = t.clearAndReturnAccumRobPrize()
			showPrizeProtos = append(showPrizeProtos, report.ShowPrize)

			mailType := shared_proto.MailType_MailReport
			if isAssemblyMail {
				mailType = shared_proto.MailType_MailAssemblyReport
			}

			r.sendHeroReportMail(t.startingBase.Id(), data,
				attackerFlagName, targetFlagName, "",
				report, true, score, ctime,
				mailTag, mailType)
		}
	}

	if !npcid.IsNpcId(targetId) {
		var targetMailProto *shared_proto.MailProto
		if data := getMailData(
			r.getMailHelp().ReportAddFail,
			r.getMailHelp().AssemblyAddFail,
			isAssemblyMail,
		); data != nil {
			if !npcid.IsNpcId(attackerId) && report.ShowPrize != nil {
				if len(showPrizeProtos) > 1 {
					prizeBuilder := resdata.NewPrizeBuilder()
					for _, p := range showPrizeProtos {
						prize := resdata.UnmarshalPrize(p, r.services.datas)
						prizeBuilder.Add(prize)
					}
					report.ShowPrize = prizeBuilder.Build().Encode()
				}

				if coef := math.Min(r.config().RobberCoef, 1); coef > 0 && coef < 1 {
					coef = 1 / coef
					report.ShowPrize.Gold = i32.MultiCoef(report.ShowPrize.Gold, coef)
					report.ShowPrize.Food = i32.MultiCoef(report.ShowPrize.Food, coef)
					report.ShowPrize.Wood = i32.MultiCoef(report.ShowPrize.Wood, coef)
					report.ShowPrize.Stone = i32.MultiCoef(report.ShowPrize.Stone, coef)
				}
			}

			mailType := shared_proto.MailType_MailReport
			if isAssemblyMail {
				mailType = shared_proto.MailType_MailAssemblyReport
			}

			targetMailProto = r.sendHeroReportMailReturnProto(targetId, data,
				attackerFlagName, targetFlagName, "",
				report, false, score, ctime,
				getMailTagByTroop(attackerTroop, targetDefenser), mailType)
		}
		if d := r.getTextHelp().BannerBaseFail; d != nil {
			r.services.world.Send(targetId, misc.NewS2cScreenShowWordsMsg(d.Text.New().WithAttacker(attackerFlagName).JsonString()))
		}

		if targetMailProto != nil {
			// 进攻方赢了，并且守城的是有部队的，设置被打战报
			if targetDefenser != nil && len(targetDefenser.captains) > 0 {

				// 发送消息
				r.heroBaseFuncWithSend(targetBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
					hero.SetTroopDefeatedMailProto(targetMailProto)
					result.AddFunc(func() pbutil.Buffer {
						return military.NewS2cSetDenfeseTroopDefeatedMailMsg(must.Marshal(targetMailProto), false)
					})

					result.Changed()
					result.Ok()
				})
			}
		}
	}

	if !npcid.IsNpcId(attackerId) && !npcid.IsNpcId(targetId) {
		fightX := int32(targetBase.BaseX())
		fightY := int32(targetBase.BaseY())

		if gid := attackerTroop.startingBase.GuildId(); gid != 0 {
			if data := r.services.datas.GuildLogHelp().FAttSucc; data != nil {
				if hero := r.services.heroSnapshotService.Get(attackerId); hero != nil {
					if hero.GuildId == gid {
						proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
						proto.Text = data.Text.New().WithAttacker(attackerFlagName).WithDefenser(targetFlagName).JsonString()
						proto.FightX = fightX
						proto.FightY = fightY

						r.services.guildService.AddLog(hero.GuildId, proto)
					}
				}
			}
		}

		if gid := targetBase.GuildId(); gid != 0 {
			if data := r.services.datas.GuildLogHelp().FDefFail; data != nil {
				if hero := r.services.heroSnapshotService.Get(targetId); hero != nil {
					if hero.GuildId == gid {
						proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
						proto.Text = data.Text.New().WithAttacker(attackerFlagName).WithDefenser(targetFlagName).JsonString()
						proto.FightX = fightX
						proto.FightY = fightY

						r.services.guildService.AddLog(hero.GuildId, proto)
					}
				}
			}
		}
	}
}

// 总的能承载的资源量
func (t *troop) totalCarryLoad() uint64 {
	var c uint64
	for _, captain := range t.captains {
		soldierCount := captain.Proto().Soldier
		loadPerSoldier := captain.Proto().LoadPerSoldier

		if soldierCount > 0 && loadPerSoldier > 0 {
			c += uint64(soldierCount) * uint64(loadPerSoldier)
		}
	}

	return c
}

func (t *troop) tryAddResource(gold, food, wood, stone uint64) (full bool, addedGold, addedFood, addedWood, addedStone uint64) {

	if t.assembly != nil {
		// 集结掠夺的东西，按人头平分
		if count := t.assembly.Count(); count > 0 {
			toAddGold := gold / count
			toAddFood := food / count
			toAddWood := wood / count
			toAddStone := stone / count

			if toAddGold|toAddFood|toAddWood|toAddStone > 0 {
				t.walkAll(func(st *troop) (toContinue bool) {
					st.tryAddResource0(toAddGold, toAddFood, toAddWood, toAddStone)
					return true
				})
			}

			// 除不尽的，给集结者
			toAddGold = gold % count
			toAddFood = food % count
			toAddWood = wood % count
			toAddStone = stone % count
			if toAddGold|toAddFood|toAddWood|toAddStone > 0 {
				t.tryAddResource0(toAddGold, toAddFood, toAddWood, toAddStone)
			}
		}
	} else {
		t.tryAddResource0(gold, food, wood, stone)
	}

	return false, gold, food, wood, stone
}

func (t *troop) tryAddResource0(gold, food, wood, stone uint64) {

	p := t.AccumRobPrize()
	p.Gold += int32(gold)
	p.Food += int32(food)
	p.Wood += int32(wood)
	p.Stone += int32(stone)
	p.IsNotEmpty = true

	t.accumRobPrize.AddUnsafeResource(gold, food, wood, stone)
	t.onChanged()
}

func (t *troop) tryAddPrize(toAdd *resdata.Prize) {
	t.accumRobPrize.Add(toAdd)
	t.clearAccumRobPrizeProto()
	t.onChanged()
}

func (t *troop) tryAddBaowu(data *resdata.BaowuData, count uint64) {
	dataId := u64.Int32(data.Id)

	p := t.AccumRobPrize()
	for i, v := range p.BaowuId {
		if v == dataId {
			p.BaowuCount[i] += u64.Int32(count)
			dataId = 0
			break
		}
	}
	if dataId != 0 {
		p.BaowuId = append(p.BaowuId, dataId)
		p.BaowuCount = append(p.BaowuCount, u64.Int32(count))

		t.accumRobBaowuCount += count
	}

	p.IsNotEmpty = true

	t.accumRobPrize.AddBaowu(data, count)
	t.onChanged()
}

func (t *troop) clearAndReturnAccumRobPrize() *shared_proto.PrizeProto {
	prize := t.AccumRobPrize()
	t.accumRobPrize = resdata.NewPrizeBuilder()
	t.clearAccumRobPrizeProto()
	t.onChanged()
	return prize
}

//func (t *troop) tryAddResource(gold, food, wood, stone uint64) (full bool, addedGold, addedFood, addedWood, addedStone uint64) {
//	totalCapacity := t.totalCarryLoad()               // 队伍最大承载上限
//	currentLoad := t.gold + t.food + t.wood + t.stone // 当前承载量
//	if currentLoad >= totalCapacity {
//		full = true
//		return
//	}
//
//	toAddLoad := gold + food + wood + stone // 要加上的承载量
//	newLoad := currentLoad + toAddLoad      // 如果全加上, 最新的承载量
//	if newLoad <= totalCapacity {
//		t.gold += gold
//		t.food += food
//		t.wood += wood
//		t.stone += stone
//
//		return newLoad == totalCapacity, gold, food, wood, stone
//	} else {
//		// 超出上限了
//		full = true
//
//		// 超出最大上限，调整
//		rate := u64.Division2Float64(totalCapacity-currentLoad, toAddLoad) // < 1
//
//		addedGold = u64.Multi(gold, rate)
//		addedFood = u64.Multi(food, rate)
//		addedWood = u64.Multi(wood, rate)
//		addedStone = u64.Multi(stone, rate)
//
//		diff := totalCapacity - currentLoad - addedGold - addedFood - addedWood - addedStone
//		if diff > 0 {
//			if more := gold - addedGold; more > 0 {
//				more = u64.Min(more, diff)
//				addedGold += more
//				diff = u64.Sub(diff, more)
//			}
//
//			if diff > 0 {
//				if more := food - addedFood; more > 0 {
//					more = u64.Min(more, diff)
//					addedFood += more
//					diff = u64.Sub(diff, more)
//				}
//
//				if diff > 0 {
//					if more := wood - addedWood; more > 0 {
//						more = u64.Min(more, diff)
//						addedWood += more
//						diff = u64.Sub(diff, more)
//					}
//
//					if diff > 0 {
//						if more := stone - addedStone; more > 0 {
//							more = u64.Min(more, diff)
//							addedStone += more
//							diff = u64.Sub(diff, more)
//						}
//
//						if diff > 0 {
//							logrus.WithField("totalCapacity", totalCapacity).WithField("currentLoad", currentLoad).WithField("gold", gold).WithField("food", food).WithField("wood", wood).WithField("stone", stone).Error("你他妈到底想怎么样, 怎么加都不对")
//						}
//					}
//				}
//			}
//		}
//
//		t.gold += addedGold
//		t.food += addedFood
//		t.wood += addedWood
//		t.stone += addedStone
//		return
//	}
//}

func (r *Realm) GetBaseLevel(prosperity uint64) *domestic_data.BaseLevelData {
	baseLevelDataArray := r.services.datas.GetBaseLevelDataArray()
	for i := len(baseLevelDataArray) - 1; i >= 0; i-- {
		data := baseLevelDataArray[i]
		if prosperity >= data.Prosperity {
			return data
		}
	}

	return nil
}

func calculateBaseLevelByProsperity(prosperity uint64, levelData *domestic_data.BaseLevelData) uint64 {
	if prosperity <= 0 {
		return 0
	}

	n := int(levelData.Level)
	for i := 0; i < n; i++ {

		if prosperity >= levelData.Prosperity {
			return levelData.Level
		}

		if levelData.GetPrevLevel() == nil {
			return levelData.Level
		}

		levelData = levelData.GetPrevLevel()
	}

	return 0
}
