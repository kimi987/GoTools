package realm

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/gen/pb/relation"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"runtime/debug"
	"time"
)

type fightType uint8

const (
	fightTypeInvadeArrive = 0
	fightTypeAssistArrive = 1
	fightTypeExpel        = 2
)

type fightContext struct {
	defeatTargetTroops bool // 打败了对方的队伍并且导致对方队伍取消防守了

	killInvadeTroopPrize  *shared_proto.PrizeProto
	killDefenseTroopPrize *shared_proto.PrizeProto
}

func copyShare(p *shared_proto.CombatShareProto) *shared_proto.CombatShareProto {
	proto := &shared_proto.CombatShareProto{}
	if p != nil {
		*proto = *p
	}
	return proto
}

// 战斗, 发战报, 发邮件. 修改队伍状态, 删除死掉的部队, 把资源还给被抢的一方
// 唯一没干的, 是发消息给双方, 修改双方的武将状态/士兵数/战斗力, 修改伤病数, 城防战斗力. 外面把所有战斗全算完后一次性一起发更新消息, 把信息合并为1条消息
func (r *Realm) fight(ctx *fightContext, attacker *troop, target *troop, fightType fightType, hasMoreTarget bool) (fightErr bool, attackerWin bool, response *server_proto.CombatXResponseServerProto, report *shared_proto.FightReportProto) {

	attackerHeroID := attacker.startingBase.Id()
	targetHeroID := target.startingBase.Id()

	var tlogBattleType uint64
	if fightType == fightTypeInvadeArrive {
		tlogBattleType = operate_type.BattleInvade
	} else if fightType == fightTypeAssistArrive {
		tlogBattleType = operate_type.BattleAssist
	} else if fightType == fightTypeExpel {
		tlogBattleType = operate_type.BattleExpel
	}
	tfctx := entity.NewTlogFightContext(tlogBattleType, 0, 0, 0)
	response = r.services.fightModule.SendFightRequest(tfctx, r.config().CombatScene, attackerHeroID, targetHeroID, attacker.toCombatPlayerProto(r), target.toCombatPlayerProto(r))

	if response.ReturnCode != 0 {
		logrus.WithField("troopid", attacker.Id()).WithField("target_troopid", target.Id()).WithField("fight_type", fightType).WithField("return_msg", response.ReturnMsg).WithField("return_code", response.ReturnCode).Error("战斗计算发生错误")
		fightErr = true
		return
	}

	// 已掠夺的资源
	//var gold, food, wood, stone uint64
	//var hasResource bool
	//if response.AttackerWin && target.state == realmface.Robbing {
	//	gold, food, wood, stone = target.ClearCarrying()
	//	hasResource = gold+food+wood+stone > 0
	//}

	// 战报, 要在修改troop的士兵数之前算出来开打前的总士兵数
	//report = r.newReport(attacker, target, response, fightType, gold, food, wood, stone)
	report = r.newReport(attacker, target, response, fightType)

	var aks, dks int32
	for _, v := range response.AttackerKillSoldier {
		aks += v
	}
	for _, v := range response.DefenserKillSoldier {
		dks += v
	}

	// 伤病率
	var attackerWoundedRate, targetWoundedRate float64
	switch fightType {
	case fightTypeInvadeArrive:
		attackerWoundedRate = r.config().AttackerWoundedRate
		if target.State() == realmface.Temp {
			targetWoundedRate = r.config().DefenserWoundedRate
		} else {
			targetWoundedRate = r.config().AssisterWoundedRate
		}
	case fightTypeExpel:
		attackerWoundedRate, targetWoundedRate = r.config().DefenserWoundedRate, r.config().AttackerWoundedRate
	default:
		attackerWoundedRate, targetWoundedRate = r.config().AssisterWoundedRate, r.config().AttackerWoundedRate
	}

	// 计算和设置队伍死伤
	var attackerWoundedSoldier, targetWoundedSoldier uint64
	if response.AttackerWin {
		// 进攻成功
		attackerWoundedSoldier = reduceSoldierToAlive(attacker, response.AttackerAliveSoldier, attackerWoundedRate)
		targetWoundedSoldier = reduceSoldierToZero(target, targetWoundedRate)
	} else {
		// 进攻失败
		attackerWoundedSoldier = reduceSoldierToZero(attacker, attackerWoundedRate)
		targetWoundedSoldier = reduceSoldierToAlive(target, response.DefenserAliveSoldier, targetWoundedRate)
	}

	ctime := r.services.timeService.CurrentTime()

	// 更新打人的一方的hero
	var attackerNewWoundedSoldier uint64
	r.heroBaseFuncWithSend(attacker.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {

		if attacker.State() != realmface.Temp {
			if response.AttackerWin {
				hero.UpdateTroop(attacker, true)
			} else {
				attacker.leaveTarget(hero, result, ctime)
				hero.RemoveTroop(attacker, true)
			}
		} else {
			logrus.WithField("fightType", fightType).WithField("stack", string(debug.Stack())).WithField("target", target).Error("fight时, attacker.state == temp")
			hero.UpdateTroopSoldier(attacker)
		}

		hero.Military().AddWoundedSoldier(attackerWoundedSoldier)
		attackerNewWoundedSoldier = hero.Military().WoundedSoldier()

		//// 掠夺方抢的就是attacker的资源, 把资源还回来（最新版本，掠夺的资源是拿不回来的）
		//if hasResource && target.targetBase != nil && target.targetBase.Id() == attackerHeroID {
		//	// fightType 一定是expel
		//	if fightType != fightTypeExpel {
		//		logrus.WithField("fightType", fightType).WithField("stack", string(debug.Stack())).WithField("target", target).Error("fight时, target在rob, 且target的targetid是attacker, 但是fightType不是Expel")
		//	} else {
		//		heromodule.AddUnsafeResource(hero, result, gold, food, wood, stone)
		//	}
		//}

		result.Changed()
		result.Ok()
	})

	// 更新被打的一方的hero
	var targetNewWoundedSoldier uint64
	if target.startingBase.isNpcBase() {
		if targetHomeNpc := GetHomeNpcBase(target.startingBase); targetHomeNpc != nil {
			if response.AttackerWin {
				r.services.heroDataService.FuncNotError(targetHomeNpc.ownerHeroId, func(hero *entity.Hero) (heroChanged bool) {
					// 被打挂了，删除掉防守部队
					npcBase := hero.GetHomeNpcBase(target.startingBase.Id())
					if npcBase != nil {
						npcBase.RemoveDefenseTroop()
						return true
					}

					return
				})
			}
		}
	} else {
		r.heroBaseFuncWithSend(target.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {

			if target.State() != realmface.Temp {
				if response.AttackerWin {
					target.leaveTarget(hero, result, ctime)
					hero.RemoveTroop(target, true)
				} else {
					hero.UpdateTroop(target, true)
				}

				//hero.Military().AddWoundedSoldier(targetWoundedSoldier)
				//targetNewWoundedSoldier = hero.Military().WoundedSoldier()
			} else {
				if fightType != fightTypeInvadeArrive && fightType != fightTypeExpel {
					logrus.WithField("fightType", fightType).WithField("stack", string(debug.Stack())).WithField("target", target).Error("fight时, target.state == temp, 但是fightType不是InvadeArrive || Expel")
				}

				if attacker.targetBase != target.startingBase {
					logrus.WithField("attacker.targetBase", attacker.targetBase).WithField("target.startingBase", target.startingBase).WithField("stack", string(debug.Stack())).WithField("target", target).Error("fight时, target.state == temp, 但是attacker.targetBase != target.startingBase")
				}

				if target.Id() == 0 {
					// 镜像防守
					if response.AttackerWin {
						hero.RemoveCopyDefenser()
						result.Add(military.REMOVE_COPY_DEFENSER_S2C)
					} else {
						hero.UpdateCopyDefenser(target)
						result.Add(hero.NewUpdateCopyDefenserMsg())
					}
				} else {
					// 出征，打的是基地，驻守部队
					hero.UpdateTroopSoldier(target)
				}

				if fightType == fightTypeExpel {
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_Expel)
					heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_EXPEL)
				}

				// 进攻方胜利了，就下掉防守部队
				//if response.AttackerWin && len(target.captains) > 0 {
				//	// 可能打败的是target不含部队，只有城墙的troop
				//
				//	if hero.GetHomeDefenseTroopIndex() != 0 {
				//		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_DEFENSER_FIGHTING)
				//
				//		hero.SetHomeDefenseTroopIndex(0)
				//		if ctx.defeatTargetTroops {
				//			result.AddFunc(func() pbutil.Buffer {
				//				return military.NewS2cSetDefenseTroopMsg(false, 0)
				//			})
				//		}
				//	}
				//}
			}

			if r.levelData.RegionType == shared_proto.RegionType_MONSTER {
				heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FIGHT_IN_JADE_REALM)
			}

			hero.Military().AddWoundedSoldier(targetWoundedSoldier)
			targetNewWoundedSoldier = hero.Military().WoundedSoldier()

			// pvp统计次数
			if !attacker.startingBase.isNpcBase() {
				// 加仇人
				hero.Relation().AddEnemy(attackerHeroID, ctime)
				result.Add(relation.NewS2cAddEnemyMsg(idbytes.ToBytes(attackerHeroID)))

				if !response.AttackerWin {
					// 战斗胜利次数
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RealmPvpSuccess)
				} else {
					// 战斗失败次数
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RealmPvpFail)
				}

				if fightType == fightTypeInvadeArrive {
					if attacker.targetBase != target.startingBase {
						// 加援助次数
						hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RealmPvpAssist)
						heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_REALM_PVP_ASSIST)

						// 援助消灭士兵数
						if dks > 0 {
							hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AssistKillSoldier, u64.FromInt32(dks))
							heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ASSIST_KILL_SOLDIER)
							heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_ASSIST_KILL_SOLDIER)
						}
					} else {
						// 防守消灭士兵数
						if dks > 0 {
							hero.HistoryAmount().Increase(server_proto.HistoryAmountType_DefenseKillSoldier, u64.FromInt32(dks))
							heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_DEFENSE_KILL_SOLDIER)
						}
					}
				}

				// tlog
				r.dep.Tlog().TlogSnsFlow(hero, operate_type.SNSAddEnemy, u64.FromInt64(attackerHeroID))
			} else {
				// 怪物攻城也加援助士兵数
				if npcid.IsMultiLevelMonsterNpcId(attacker.startingBase.Id()) {
					if fightType == fightTypeInvadeArrive {
						if attacker.targetBase != target.startingBase {

							// 援助消灭士兵数
							if dks > 0 {
								hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AssistKillSoldier, u64.FromInt32(dks))
								heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ASSIST_KILL_SOLDIER)
								heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_ASSIST_KILL_SOLDIER)
							}
						}
					}
				}
			}
			result.Changed()
			result.Ok()
		})
	}

	//// 如果驱逐成功, 且对方的target不是attacker, 一定是assist （最新版本，掠夺的资源是拿不回来的）
	//if hasResource && target.targetBase != nil && target.targetBase.Id() != attackerHeroID {
	//	if fightType != fightTypeAssistArrive {
	//		logrus.WithField("fightType", fightType).WithField("stack", string(debug.Stack())).WithField("target", target).Error("fight时, target在rob, 且target的targetid是attacker, 但是fightType不是Assist")
	//	} else {
	//		r.heroBaseFuncWithSend(target.targetBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
	//			heromodule.AddUnsafeResource(hero, result, gold, food, wood, stone)
	//			return
	//		})
	//	}
	//}

	// 处理部队被搞死逻辑
	r.handleTroopFight(ctx, attacker, target, fightType, response.AttackerWin, ctime)

	// 发战报邮件
	r.sendReportMail(response, attacker, attackerHeroID, target, targetHeroID, fightType, response.AttackerWin, report,
		ctx.killInvadeTroopPrize, ctx.killDefenseTroopPrize, ctime)
	ctx.killInvadeTroopPrize = nil
	ctx.killDefenseTroopPrize = nil

	// 击退列表
	if response.AttackerWin {
		if fightType == fightTypeAssistArrive {
			// 援助的人胜利了，更新击退列表
			attacker.AddKillEnemy(target.getStartingBaseFlagHeroName(r))
		}
	} else {
		if fightType == fightTypeInvadeArrive {
			if attacker.targetBase != target.startingBase {
				// 2个人的目标城池一致，属于援助
				// 掠夺的人失败了，更新击退列表
				target.AddKillEnemy(attacker.getStartingBaseFlagHeroName(r))
			}
		}
	}

	// 删除失败的部队
	if response.AttackerWin {
		if fightType != fightTypeExpel {
			r.trySendTroopDoneMail(target, r.getTextHelp().MRDRExpelled4a.Text, r.getTextHelp().MRDRExpelled4d.Text, ctime)
		} else {
			// 更新attacker的信息
			r.broadcastMaToCared(attacker, addTroopTypeUpdate, 0)
		}

		target.removeWithoutReturnCaptain(r)
		if fightType != fightTypeExpel {
			r.broadcastRemoveMaToCared(target)
		}
		target.clearMsgs()
	} else {
		if fightType == fightTypeExpel {
			r.trySendTroopDoneMail(attacker, r.getTextHelp().MRDRExpelled4a.Text, r.getTextHelp().MRDRExpelled4d.Text, ctime)
		}

		attacker.removeWithoutReturnCaptain(r)
		r.broadcastRemoveMaToCared(attacker)
		attacker.clearMsgs()

		// 更新防守部队
		if target.startingBase != target.targetBase {
			r.broadcastMaToCared(target, addTroopTypeUpdate, 0)
		}
	}

	// 发送更新士兵，武将消息
	if !hasMoreTarget || !response.AttackerWin {
		// 你输了，或者当前是最后一个对手，发消息
		r.services.world.SendFunc(attackerHeroID, func() pbutil.Buffer {
			proto := &region.S2CUpdateSelfTroopsProto{}
			for _, c := range attacker.captains {
				proto.Id = append(proto.Id, u64.Int32(c.Id()))
				proto.Soldier = append(proto.Soldier, c.Proto().Soldier)
				proto.FightAmount = append(proto.FightAmount, data.ProtoFightAmount(c.Proto().TotalStat, c.Proto().Soldier, c.Proto().SpellFightAmountCoef))
			}
			proto.WoundedSoldier = u64.Int32(attackerNewWoundedSoldier)

			proto.RemoveOutside = !response.AttackerWin
			if attacker.State() == realmface.Temp {
				proto.RemoveOutside = true
			}

			proto.TroopIndex = u64.Int32(entity.GetTroopIndex(attacker.Id()) + 1)

			return region.NewS2cUpdateSelfTroopsProtoMsg(proto)
		})
	}

	// 被打的人，更新部队士兵数
	if len(target.captains) > 0 {
		if target.Id() != 0 {
			r.services.world.SendFunc(targetHeroID, func() pbutil.Buffer {
				proto := &region.S2CUpdateSelfTroopsProto{}
				for _, c := range target.captains {
					proto.Id = append(proto.Id, u64.Int32(c.Id()))
					proto.Soldier = append(proto.Soldier, c.Proto().Soldier)
					proto.FightAmount = append(proto.FightAmount, data.ProtoFightAmount(c.Proto().TotalStat, c.Proto().Soldier, c.Proto().SpellFightAmountCoef))
				}
				proto.WoundedSoldier = u64.Int32(targetNewWoundedSoldier)
				proto.RemoveOutside = response.AttackerWin
				if target.State() == realmface.Temp {
					proto.RemoveOutside = true
				}
				proto.TroopIndex = u64.Int32(entity.GetTroopIndex(target.Id()) + 1)

				return region.NewS2cUpdateSelfTroopsProtoMsg(proto)
			})
		}
	}

	r.heroBaseFuncWithSend(attacker.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
		//heromodule.AutoRecoverCaptainSoldier(hero, result)

		if r.levelData.RegionType == shared_proto.RegionType_MONSTER {
			heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FIGHT_IN_JADE_REALM)

			if response.AttackerWin && npcid.IsMonsterNpcId(target.Id()) {
				// 打败的是npc怪物
				heromodule.IncreTaskProgressWithFunc(hero, result, shared_proto.TaskTargetType_TASK_TARGET_JADE_NPC, func(hero *entity.Hero, t *entity.TaskProgress) (increAmount uint64) {
					if t.Target().RegionData == r.levelData && t.Target().MonsterLevel == target.StartingBase().GetBaseLevel() {
						return 1
					}

					return 0
				})
			}
		}

		//if response.AttackerWin && npcid.IsHomeNpcId(target.Id()) {
		//	hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_KillHomeNpc)
		//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_KILL_HOME_NPC)
		//}

		// pvp统计次数
		if !target.startingBase.isNpcBase() {
			if response.AttackerWin {
				// 战斗胜利次数
				hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RealmPvpSuccess)
			} else {
				// 战斗失败次数
				hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RealmPvpFail)
			}

			if fightType == fightTypeAssistArrive {
				// 加援助次数
				hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RealmPvpAssist)
				heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_REALM_PVP_ASSIST)
			}

			if aks > 0 {
				toAdd := u64.FromInt32(aks)
				switch fightType {
				case fightTypeExpel:
					fallthrough
				case fightTypeInvadeArrive:
					// 进攻消灭士兵数
					hero.HistoryAmount().Increase(server_proto.HistoryAmountType_InvaseKillSoldier, toAdd)
					heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_INVASE_KILL_SOLDIER)
					heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_INVASE_KILL_SOLDIER)

				case fightTypeAssistArrive:
					// 援助消灭士兵数
					hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AssistKillSoldier, toAdd)
					heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ASSIST_KILL_SOLDIER)
					heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_ASSIST_KILL_SOLDIER)
				}
			}
		} else {
			// 怪物攻城也加援助士兵数
			if npcid.IsMultiLevelMonsterNpcId(attacker.startingBase.Id()) {
				if fightType == fightTypeInvadeArrive {
					if attacker.targetBase != target.startingBase {

						// 援助消灭士兵数
						if dks > 0 {
							hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AssistKillSoldier, u64.FromInt32(dks))
							heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ASSIST_KILL_SOLDIER)
							heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_ASSIST_KILL_SOLDIER)
						}
					}
				}
			}
		}

		result.Ok()
	})

	// 加被援助次数
	if attacker.targetBase != target.startingBase {
		// 进攻的人打的不是防守的人，进攻的目标就是被援助目标
		r.heroBaseFuncWithSend(attacker.targetBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
			// 加被援助次数
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RealmPvpBeenAssist)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_REALM_PVP_BEEN_ASSIST)
		})
	}

	// 更新宝藏Npc的防守士兵数和防守战力
	if npcid.IsBaoZangNpcId(target.startingBase.Id()) {
		target.startingBase.updateRoBase()
	}

	// 更新匈奴目标数量（匈奴部队如果输了，更新）
	if response.AttackerWin {
		if npcid.IsXiongNuNpcId(target.startingBase.Id()) {
			target.startingBase.updateXiongNuTarget()
		}
	} else {
		if npcid.IsXiongNuNpcId(attacker.startingBase.Id()) {
			attacker.startingBase.updateXiongNuTarget()
		}
	}

	// 匈奴数据战斗数据
	if xiongNu := GetXiongNuBase(attacker.startingBase); xiongNu != nil {
		xiongNu.info.AddHeroFightSoldier(target.startingBase.Id(), u64.FromInt32(dks), u64.FromInt32(aks))
	}
	if xiongNu := GetXiongNuBase(target.startingBase); xiongNu != nil {
		xiongNu.info.AddHeroFightSoldier(attacker.startingBase.Id(), u64.FromInt32(aks), u64.FromInt32(dks))
	}

	// 更新镜像队伍
	if response.AttackerWin && target.Id() == 0 {
		if data := r.getMailHelp().CopyDefenserBeenKilled; data != nil {
			proto := data.NewTextMail(shared_proto.MailType_MailNormal)
			r.services.mail.SendProtoMail(target.startingBase.Id(), proto, ctime)
		}
	}

	// 任务怪，更新英雄任务
	if response.AttackerWin {
		r.updateHeroSpecMonsterTask(target)
	} else {
		r.updateHeroSpecMonsterTask(attacker)
	}

	if defenserBase := attacker.targetBase; attackerWin && defenserBase != nil {
		if npcid.IsXiongNuNpcId(defenserBase.Id()) || npcid.IsJunTuanNpcId(defenserBase.Id()) {
			defenserBase.updateRoBase()
			defenserBase.ClearUpdateBaseInfoMsg()
			r.broadcastBaseInfoToCared(defenserBase, addBaseTypeUpdate, 0)
		}
	}

	return false, response.AttackerWin, response, report
}

func (r *Realm) updateHeroSpecMonsterTask(t *troop) {
	if t.mmd == nil {
		return
	}

	if t.mmd.GetSpec() != monsterdata.InvasionTask {
		return
	}

	if npcid.IsNpcId(t.originTargetId) {
		return
	}

	monsterId := t.mmd.Id
	r.heroFuncWithSend(t.originTargetId, func(hero *entity.Hero, result herolock.LockResult) {
		// 移除任务怪
		hero.TaskList().RemoveTaskMonster(monsterId)

		// 添加任务怪
		hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_ExpelFightMonster, monsterId)

		// 更新任务
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_EXPEL_FIGHT_MONSTER)

		result.Ok()
	})
}

func (r *Realm) handleTroopFight(fctx *fightContext, attacker, target *troop, fightType fightType, attackerWin bool, ctime time.Time) {

	// 处理飘字
	switch fightType {
	case fightTypeAssistArrive:

		assister := attacker
		robber := target

		if attackerWin {
			// 援助成功飘字
			r.services.world.SendFunc(assister.startingBase.Id(), func() pbutil.Buffer {
				return misc.NewS2cScreenShowWordsMsg(
					r.getTextHelp().RealmAssistSuccess.New().
						WithTroopIndex(entity.GetTroopIndex(assister.Id()) + 1).
						JsonString())
			})

			// 掠夺者失败飘字
			r.services.world.SendFunc(robber.startingBase.Id(), func() pbutil.Buffer {
				return misc.NewS2cScreenShowWordsMsg(
					r.getTextHelp().BannerTroopFail.New().
						WithAttacker(assister.getStartingBaseFlagHeroName(r)).
						WithTroopIndex(entity.GetTroopIndex(robber.Id()) + 1).
						JsonString())
			})

		} else {
			// 援助失败飘字
			r.services.world.SendFunc(assister.startingBase.Id(), func() pbutil.Buffer {
				return misc.NewS2cScreenShowWordsMsg(
					r.getTextHelp().RealmAssistFail.New().
						WithTroopIndex(entity.GetTroopIndex(assister.Id()) + 1).
						JsonString())
			})

			// 掠夺者胜利飘字
			r.services.world.SendFunc(robber.startingBase.Id(), func() pbutil.Buffer {
				return misc.NewS2cScreenShowWordsMsg(
					r.getTextHelp().BannerTroopSuccess.New().
						WithAttacker(assister.getStartingBaseFlagHeroName(r)).
						WithTroopIndex(entity.GetTroopIndex(robber.Id()) + 1).
						JsonString())
			})
		}

	case fightTypeInvadeArrive:
		robber := attacker

		if !attackerWin {
			// 出征失败飘字
			r.services.world.SendFunc(robber.startingBase.Id(), func() pbutil.Buffer {
				return misc.NewS2cScreenShowWordsMsg(
					r.getTextHelp().RealmInvadeFail.New().
						WithTroopIndex(entity.GetTroopIndex(robber.Id()) + 1).
						JsonString())
			})

			if robber.targetBase != target.startingBase {
				// 防守的援助部队胜利飘字
				r.services.world.SendFunc(target.startingBase.Id(), func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().BannerTroopSuccess.New().
							WithAttacker(attacker.getStartingBaseFlagHeroName(r)).
							WithTroopIndex(entity.GetTroopIndex(target.Id()) + 1).
							JsonString())
				})
			}
		} else {
			if robber.targetBase != target.startingBase {
				// 防守的援助部队失败飘字
				r.services.world.SendFunc(target.startingBase.Id(), func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().BannerTroopFail.New().
							WithAttacker(attacker.getStartingBaseFlagHeroName(r)).
							WithTroopIndex(entity.GetTroopIndex(target.Id()) + 1).
							JsonString())
				})
			}
		}
	case fightTypeExpel:

		robbing := attacker
		if !attackerWin {
			// 遭受驱逐，战斗失败
			r.services.world.SendFunc(robbing.startingBase.Id(), func() pbutil.Buffer {
				return misc.NewS2cScreenShowWordsMsg(
					r.getTextHelp().BannerTroopFail.New().
						WithAttacker(target.getStartingBaseFlagHeroName(r)).
						WithTroopIndex(entity.GetTroopIndex(robbing.Id()) + 1).
						JsonString())
			})

		} else {
			// 遭受驱逐，战斗胜利
			r.services.world.SendFunc(robbing.startingBase.Id(), func() pbutil.Buffer {
				return misc.NewS2cScreenShowWordsMsg(
					r.getTextHelp().BannerTroopSuccess.New().
						WithAttacker(target.getStartingBaseFlagHeroName(r)).
						WithTroopIndex(entity.GetTroopIndex(robbing.Id()) + 1).
						JsonString())
			})
		}
	}

	// 处理击杀入侵部队
	switch fightType {
	case fightTypeAssistArrive:
		if attackerWin {
			r.tryKillInvadeTroop(fctx, attacker, target, attacker.targetBase, ctime)
		}

	case fightTypeInvadeArrive:
		if !attackerWin {
			r.tryKillInvadeTroop(fctx, target, attacker, attacker.targetBase, ctime)
		}
	case fightTypeExpel:
		if !attackerWin {
			r.tryKillInvadeTroop(fctx, target, attacker, target.startingBase, ctime)
		}
	}

	if attackerWin {
		r.tryKillTroop(fctx, attacker, target, ctime)
	} else {
		r.tryKillTroop(fctx, target, attacker, ctime)
	}

	// 出征打Npc输了
	if fightType == fightTypeInvadeArrive && !attackerWin {
		if defenseBaseInfo := attacker.getTargetBaseInfo(); defenseBaseInfo != nil {
			// 失败仇恨变化
			if hateData := defenseBaseInfo.getHateData(); hateData != nil && hateData.FailHate != 0 {
				r.heroBaseFuncWithSend(attacker.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
					if info := hero.GetNpcTypeInfo(hateData.TypeData().Type); info != nil {
						newHate := u64.AddInt(info.GetHate(), hateData.FailHate)
						newHate = u64.Min(newHate, hateData.TypeData().MaxHate)

						if info.SetHate(newHate) {
							result.Add(region.NewS2cUpdateMultiLevelNpcHateMsg(int32(hateData.TypeData().Type), u64.Int32(newHate)))
						}
					}
				})
			}
		}
	}
}

func (r *Realm) tryKillTroop(fctx *fightContext, killer, beenKilled *troop, ctime time.Time) {

	if beenKilled.startingBase.isNpcBase() {
		// 被杀的是个npc，给npc部队击杀奖励
		if beenKilled.mmd != nil && beenKilled.mmd.BeenKillPrize != nil {
			prize := beenKilled.mmd.BeenKillPrize.GetPrize()
			prizeProto := prize.Encode()

			// 击杀入侵奖励
			r.heroBaseFuncWithSend(killer.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
				hctx := heromodule.NewContext(r.dep, operate_type.RealmKillDefNpc)

				heromodule.AddPrize(hctx, hero, result, prize, ctime)
				fctx.killDefenseTroopPrize = prizeProto

				result.Ok()
			})
		}

		// 被杀的是个匈奴怪
		if base := GetXiongNuBase(beenKilled.startingBase); base != nil {

			// 防守怪物，不降低士气
			isReduceMorale := beenKilled.State() != realmface.Defending && beenKilled.State() != realmface.Temp

			var toReduceMorale uint64
			if isReduceMorale {
				toReduceMorale = r.services.datas.ResistXiongNuMisc().WipeOutReduceMorale
			}

			// 减少怪物
			base.Info().WipeOutMonster(ctime, toReduceMorale)
			if isReduceMorale {
				// 降低士气
				r.onWipeOutXiongNuNpcTroop(beenKilled.startingBase, base, toReduceMorale)

				if data := r.services.datas.GuildLogHelp().WipeOutXiongNuTroop; data != nil {
					hero := r.services.heroSnapshotService.Get(killer.Id())
					if hero != nil {
						proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
						proto.Text = data.Text.New().
							WithHeroName(hero.Name).
							WithResistXiongNuName(beenKilled.getStartingBaseInfo().getName()).
							WithReduceMorale(toReduceMorale).JsonString()
						r.services.guildService.AddLog(base.Info().GuildId(), proto)
					}
				}
			}

		}
	}

}

func (r *Realm) tryKillInvadeTroop(fctx *fightContext, killer, invadeTroop *troop, defenseBase *baseWithData, ctime time.Time) {
	if !invadeTroop.startingBase.isNpcBase() {
		// 非Npc无特殊处理
		return
	}

	if defenseBase != nil && invadeTroop.mmd != nil && invadeTroop.mmd.InvadePrize != nil {
		// 如果有击杀奖励，那么自己击杀的，不给入侵奖励
		if invadeTroop.mmd.BeenKillPrize == nil || killer.startingBase != defenseBase {

			prize := invadeTroop.mmd.InvadePrize.GetPrize()
			prizeProto := prize.Encode()

			// 击杀入侵奖励
			r.heroBaseFuncWithSend(defenseBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
				hctx := heromodule.NewContext(r.dep, operate_type.RealmKillInvadeNpc)

				heromodule.AddPrize(hctx, hero, result, prize, ctime)
				fctx.killInvadeTroopPrize = prizeProto

				result.Ok()
			})
		}

	}

	//baseInfo := invadeTroop.getStartingBaseInfo()
	//if hateData := baseInfo.getHateData(); hateData != nil && defenseBase != nil {
	//	r.heroBaseFuncWithSend(defenseBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
	//		hctx := heromodule.NewContext(r.dep, operate_type.RealmKillInvadeNpc)
	//
	//		// 攻城胜利奖励
	//		if hateData.FightNpcPrize != nil {
	//			heromodule.AddPrize(hctx, hero, result, hateData.FightNpcPrize, ctime)
	//			fctx.killInvadeTroopPrize = hateData.FightNpcPrizeProto
	//		}
	//
	//		result.Ok()
	//	})
	//}

}

func (r *Realm) sendReportMail(response *server_proto.CombatXResponseServerProto, attacker *troop, attackerHeroID int64,
	target *troop, targetHeroID int64, fightType fightType, attackerWin bool, report *shared_proto.FightReportProto,
	fightNpcPrize, killDefPrize *shared_proto.PrizeProto, ctime time.Time) {
	// 发送战报邮件

	switch fightType {
	case fightTypeAssistArrive:

		attackerId := target.startingBase.Id()
		defenserId := attacker.targetBase.Id()
		assisterId := attacker.startingBase.Id()

		attackerFlagName := target.getStartingBaseFlagHeroName(r)
		defenserFlagName := attacker.getTargetBaseFlagHeroName(r)
		assisterFlagName := attacker.getStartingBaseFlagHeroName(r)
		if attackerWin {
			// 进攻方失败
			if data := r.getMailHelp().ReportSaaFail; data != nil {
				r.sendHeroReportMail(attackerId, data,
					attackerFlagName, defenserFlagName, assisterFlagName,
					report, false, response.GetScore(), ctime,
					getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
			}

			// 援助方成功
			if data := r.getMailHelp().ReportSasSuccess; data != nil {
				report.ShowPrize = killDefPrize

				r.sendHeroReportMail(assisterId, data,
					attackerFlagName, defenserFlagName, assisterFlagName,
					report, true, response.GetScore(), ctime,
					getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
			}

			// 城主成功
			if data := r.getMailHelp().ReportSadSuccess; data != nil {
				report.ShowPrize = fightNpcPrize

				r.sendHeroReportMail(defenserId, data,
					attackerFlagName, defenserFlagName, assisterFlagName,
					report, true, response.GetScore(), ctime,
					getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
			}
		} else {
			// 进攻方成功
			if data := r.getMailHelp().ReportSaaSuccess; data != nil {
				report.ShowPrize = killDefPrize

				r.sendHeroReportMail(attackerId, data,
					attackerFlagName, defenserFlagName, assisterFlagName,
					report, false, response.GetScore(), ctime,
					getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
			}

			// 援助方失败
			if data := r.getMailHelp().ReportSasFail; data != nil {
				r.sendHeroReportMail(assisterId, data,
					attackerFlagName, defenserFlagName, assisterFlagName,
					report, true, response.GetScore(), ctime,
					getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
			}

			// 城主失败
			if data := r.getMailHelp().ReportSadFail; data != nil {
				r.sendHeroReportMail(defenserId, data,
					attackerFlagName, defenserFlagName, assisterFlagName,
					report, true, response.GetScore(), ctime,
					getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
			}
		}

	case fightTypeInvadeArrive:
		attackerId := attacker.startingBase.Id()
		defenserId := attacker.targetBase.Id()
		assisterId := target.startingBase.Id()

		attackerFlagName := attacker.getStartingBaseFlagHeroName(r)
		defenserFlagName := attacker.getTargetBaseFlagHeroName(r)
		assisterFlagName := defenserFlagName
		if assisterId != defenserId {
			assisterFlagName = target.getStartingBaseFlagHeroName(r)
		}

		if attackerWin {

			// 最后一波不发送（匈奴入侵，starting是匈奴，但是状态是defending）
			if target.startingBase != attacker.targetBase || target.State() == realmface.Defending {
				// 掠夺者干掉了盟友
				// 进攻方胜利
				if data := r.getMailHelp().ReportAsaSuccess; data != nil {
					report.ShowPrize = killDefPrize

					r.sendHeroReportMail(attackerId, data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, true, response.GetScore(), ctime,
						getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
				}

				// 援助方失败
				if data := r.getMailHelp().ReportAssFail; data != nil {
					r.sendHeroReportMail(assisterId, data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, false, response.GetScore(), ctime,
						getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
				}

				// 城主失败
				if data := r.getMailHelp().ReportAsdFail; data != nil {
					r.sendHeroReportMail(defenserId, data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, false, response.GetScore(), ctime,
						getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
				}
			}

		} else {
			// 来抢劫的, 自己挂了
			// 掠夺者被干掉
			if target.startingBase == attacker.targetBase {
				// 防守阵容干掉的

				// 进攻方失败
				if data := r.getMailHelp().ReportAdaFail; data != nil {
					r.sendHeroReportMail(attackerId, data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, true, response.GetScore(), ctime,
						getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
				}

				// 城主成功
				if data := r.getMailHelp().ReportAddSuccess; data != nil {
					report.ShowPrize = fightNpcPrize
					if report.ShowPrize == nil {
						report.ShowPrize = killDefPrize
					}

					r.sendHeroReportMail(defenserId, data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, false, response.GetScore(), ctime,
						getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
				}
				if d := r.getTextHelp().BannerBaseSuccess; d != nil {
					r.services.world.Send(defenserId, misc.NewS2cScreenShowWordsMsg(d.Text.New().WithAttacker(attackerFlagName).JsonString()))
				}
			} else {
				// 盟友干掉的
				// 进攻方失败
				if data := r.getMailHelp().ReportAsaFail; data != nil {
					r.sendHeroReportMail(attackerId, data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, true, response.GetScore(), ctime,
						getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
				}

				// 援助方成功
				if data := r.getMailHelp().ReportAssSuccess; data != nil {
					report.ShowPrize = killDefPrize

					r.sendHeroReportMail(assisterId, data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, false, response.GetScore(), ctime,
						getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
				}

				// 城主成功
				if data := r.getMailHelp().ReportAsdSuccess; data != nil {
					report.ShowPrize = fightNpcPrize

					r.sendHeroReportMail(defenserId, data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, false, response.GetScore(), ctime,
						getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
				}
			}

		}

	case fightTypeExpel:

		attackerId := attacker.startingBase.Id()
		defenserId := target.startingBase.Id()

		attackerFlagName := attacker.getStartingBaseFlagHeroName(r)
		defenserFlagName := target.getStartingBaseFlagHeroName(r)

		if attackerWin {
			// 驱逐失败
			// 进攻方成功
			if data := r.getMailHelp().ReportExpelAttackerSuccess; data != nil {
				report.ShowPrize = killDefPrize

				r.sendHeroReportMail(attackerId, data,
					attackerFlagName, defenserFlagName, "",
					report, true, response.GetScore(), ctime,
					getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
			}

			// 城主失败
			if data := r.getMailHelp().ReportExpelDefenserFail; data != nil {
				r.sendHeroReportMail(defenserId, data,
					attackerFlagName, defenserFlagName, "",
					report, false, response.GetScore(), ctime,
					getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
			}
		} else {
			// 驱逐成功
			// 进攻方失败
			if data := r.getMailHelp().ReportExpelAttackerFail; data != nil {
				r.sendHeroReportMail(attackerId, data,
					attackerFlagName, defenserFlagName, "",
					report, true, response.GetScore(), ctime,
					getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
			}

			// 城主成功
			if data := r.getMailHelp().ReportExpelDefenserSuccess; data != nil {
				report.ShowPrize = fightNpcPrize
				if report.ShowPrize == nil {
					report.ShowPrize = killDefPrize
				}

				r.sendHeroReportMail(defenserId, data,
					attackerFlagName, defenserFlagName, "",
					report, false, response.GetScore(), ctime,
					getMailTagByTroop(attacker, target), shared_proto.MailType_MailReport)
			}
		}
	}
}

func (r *Realm) newReport(attacker *troop, target *troop, result *server_proto.CombatXResponseServerProto, fightType fightType) *shared_proto.FightReportProto {
	proto := &shared_proto.FightReportProto{}

	proto.AttackerWin = result.AttackerWin
	proto.ReplayUrl = result.Link
	proto.Share = copyShare(result.DefenserShare)

	proto.Attacker = attacker.toReportHeroProto(r, result.AttackerAliveSoldier, result.AttackerKillSoldier)
	proto.Defenser = target.toReportHeroProto(r, result.DefenserAliveSoldier, result.DefenserKillSoldier)

	//proto.Gold = u64.Int32(gold)
	//proto.Food = u64.Int32(food)
	//proto.Wood = u64.Int32(wood)
	//proto.Stone = u64.Int32(stone)
	//proto.JadeOre = u64.Int32(attacker.jadeOre)

	proto.Score = result.Score
	proto.FightType = int32(fightType)
	// 战斗坐标
	if attacker.targetBase != nil && attacker.targetBase != target.startingBase {
		proto.FightX, proto.FightY = imath.Int32(attacker.targetBase.BaseX()), imath.Int32(attacker.targetBase.BaseY())

		proto.FightTargetId = attacker.targetBase.IdBytes()
		proto.FightTargetName = attacker.getTargetBaseInfo().getName()
		proto.FightTargetFlagName = r.getFlagName(attacker.targetBase.GuildId())

	} else {
		if attacker.targetBase == nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("fightType", fightType).Error("竟然attacker的target为nil")
		}
		proto.FightX, proto.FightY = proto.Defenser.BaseX, proto.Defenser.BaseY

		proto.FightTargetId = proto.Defenser.Id
		proto.FightTargetName = proto.Defenser.Name
		proto.FightTargetFlagName = proto.Defenser.GuildFlagName
	}

	// 提示语言
	proto.AttackerDesc, proto.DefenserDesc = r.getTextHelp().MailReportWinnerDesc.Text.KeysOnlyJson(), r.getTextHelp().MailReportLoserDesc.Text.KeysOnlyJson()
	//switch fightType {
	//case fightTypeAssistArrive:
	//	proto.AttackerDesc, proto.DefenserDesc = r.services.datas.MailConfig().RegionAssistActDesc.OneText, r.services.datas.MailConfig().RegionAssistDefDesc.OneText
	//
	//case fightTypeInvadeArrive:
	//	proto.AttackerDesc, proto.DefenserDesc = r.services.datas.MailConfig().RegionInvadeActDesc.OneText, r.services.datas.MailConfig().RegionInvadeDefDesc.OneText
	//
	//case fightTypeExpel:
	//	proto.AttackerDesc, proto.DefenserDesc = r.services.datas.MailConfig().RegionExpelActDesc.OneText, r.services.datas.MailConfig().RegionExpelDefDesc.OneText
	//}

	if attacker.npcTimes > 1 {
		if d := r.dep.Datas().TextHelp().MultiMonsterPrizeCount; d != nil {
			proto.Desc = d.Text.New().WithCount(attacker.npcTimes).JsonString()
		}
	}

	return proto
}

//func (r *Realm) newNoFightReport(attacker *troop, target *baseWithData, attackerDesc, defenserDesc string) *shared_proto.FightReportProto {
//	proto := &shared_proto.FightReportProto{}
//
//	proto.AttackerWin = true
//
//	proto.Attacker = attacker.toReportHeroProto(r, nil)
//	proto.Defenser = target.toReportHeroProto(r.services.guildService.GetSnapshot)
//
//	proto.Gold = u64.Int32(attacker.gold)
//	proto.Food = u64.Int32(attacker.food)
//	proto.Wood = u64.Int32(attacker.wood)
//	proto.Stone = u64.Int32(attacker.stone)
//	proto.JadeOre = u64.Int32(attacker.jadeOre)
//
//	// 战斗坐标
//	proto.FightX, proto.FightY = imath.Int32(target.BaseX()), imath.Int32(target.BaseY())
//
//	// 提示语言
//	proto.AttackerDesc, proto.DefenserDesc = attackerDesc, defenserDesc
//
//	return proto
//}

func (r *Realm) toBaseFlagHeroName(b *baseWithData, baseLevel uint64) string {
	flagName := r.getFlagName(b.GuildId())
	heroName := b.internalBase.getBaseInfoByLevel(baseLevel).getName()
	return r.toFlagHeroName(flagName, heroName)
}

func (r *Realm) toFlagHeroName(flagName, heroName string) string {
	return r.services.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(flagName, heroName)
}

func (r *Realm) toFlagHeroNameByGuildId(guildId int64, heroName string) string {
	return r.toFlagHeroName(r.getFlagName(guildId), heroName)
}

func (r *Realm) toFlagHeroNameByHeroId(heroId int64) string {
	if hero := r.services.heroSnapshotService.Get(heroId); hero != nil {
		return r.toFlagHeroNameByGuildId(hero.GuildId, hero.Name)
	}
	return idbytes.PlayerName(heroId)
}

func (r *Realm) getGuildNameAndFlagName(guildId int64) (guildName string, flagName string) {
	if guildId != 0 {
		g := r.services.guildService.GetSnapshot(guildId)
		if g != nil {
			return g.Name, g.FlagName
		}
	}
	return "", ""
}

func (r *Realm) getFlagName(guildId int64) string {
	if guildId != 0 {
		g := r.services.guildService.GetSnapshot(guildId)
		if g != nil {
			return g.FlagName
		}
	}
	return ""
}

func (r *Realm) getCountry(guildId int64) uint64 {
	if guildId != 0 {
		g := r.services.guildService.GetSnapshot(guildId)
		if g != nil && g.Country != nil {
			return g.Country.Id
		}
	}
	return 0
}

func reduceSoldierToZero(t *troop, woundRate float64) (woundedCount uint64) {
	var deadSoldierCounter uint64

	for _, captain := range t.captains {
		deadSoldierCounter += captain.GetAndSetSoldierCount(0)
	}

	woundedCount = u64.MultiCoef(deadSoldierCounter, woundRate)
	t.onChanged()
	return
}

func reduceSoldierToAlive(t *troop, aliveSoldierMap map[int32]int32, woundRate float64) (woundedCount uint64) {
	var deadSoldierCounter uint64

	for _, captain := range t.captains {
		if alive, has := aliveSoldierMap[int32(captain.Id())]; has && alive > 0 {
			oldCount := captain.GetAndSetSoldierCount(uint64(alive))

			if oldCount < uint64(alive) {
				logrus.WithField("captain", captain).WithField("old_count", oldCount).WithField("new_count", alive).Error("打完一场后, 活着的士兵竟然超过了开打时的士兵数")
				captain.GetAndSetSoldierCount(oldCount)
				continue
			}
			deadCount := oldCount - uint64(alive)
			deadSoldierCounter += deadCount
		} else {
			// 没找到, 当做全死光了
			oldCount := captain.GetAndSetSoldierCount(0)
			deadSoldierCounter += oldCount
		}
	}

	woundedCount = u64.MultiCoef(deadSoldierCounter, woundRate)
	t.onChanged()
	return
}
