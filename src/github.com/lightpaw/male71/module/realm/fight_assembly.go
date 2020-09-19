package realm

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/random"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/idbytes"
	"time"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/gen/pb/relation"
	"github.com/lightpaw/male7/config/resdata"
	"runtime/debug"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/entity/npcid"
)

func newFightxContext() *fightxContext {
	return &fightxContext{
		subTroopMap: make(map[int64]*fightxTroop),
	}
}

type fightxContext struct {
	battles []*battle

	subTroopMap map[int64]*fightxTroop
}

func newFightSubtroop(sub *troop) *fightxTroop {

	beforeSoldier := make([]int32, len(sub.captains))
	hasSoldier := false
	for i, c := range sub.captains {
		if c.Proto().Soldier > 0 {
			beforeSoldier[i] = c.Proto().Soldier
			hasSoldier = true
		}
	}
	if !hasSoldier && sub.wallStat == nil {
		return nil
	}

	fst := &fightxTroop{}
	fst.sub = sub
	fst.facadeTroop = sub.getFacadeTroop()
	fst.beforeSoldier = beforeSoldier

	return fst
}

type fightxTroop struct {
	facadeTroop *troop

	sub *troop

	beforeSoldier []int32

	battles []*battle

	winTimes uint64
}

func newBattle(attacker, target *troop, attackerFightAmount, targetFightAmount, attackerAliveSoldier, targetAliveSoldier, attackerTotalSoldier, targetTotalSoldier int32, response *server_proto.CombatXResponseServerProto, winTimes int32) *battle {
	f := &battle{}

	f.attacker = attacker.getFacadeTroop()
	f.target = target.getFacadeTroop()

	f.attackerSubTroop = attacker
	f.targetSubTroop = target

	f.attackerFightAmount = attackerFightAmount
	f.targetFightAmount = targetFightAmount

	f.attackerAliveSoldier = attackerAliveSoldier
	f.targetAliveSoldier = targetAliveSoldier

	f.attackerTotalSoldier = attackerTotalSoldier
	f.targetTotalSoldier = targetTotalSoldier

	f.response = response

	f.winTimes = winTimes

	return f
}

type battle struct {
	attacker *troop
	target   *troop

	// 对象
	attackerSubTroop *troop
	targetSubTroop   *troop

	// 战力
	attackerFightAmount int32
	targetFightAmount   int32

	// 存活兵力
	attackerAliveSoldier int32
	targetAliveSoldier   int32

	// 总兵力
	attackerTotalSoldier int32
	targetTotalSoldier   int32

	response *server_proto.CombatXResponseServerProto

	// 胜利次数，3胜
	winTimes int32
}

func (f *battle) getOtherSubtroop(sub *troop) (other *troop, otherIsAttacker bool) {
	if f.attackerSubTroop == sub {
		return f.targetSubTroop, false
	} else {
		return f.targetSubTroop, true
	}
}

func (f *battle) toFightProto() *shared_proto.AssemblyFightProto {
	proto := &shared_proto.AssemblyFightProto{}
	proto.AttackerId = f.attackerSubTroop.startingBase.IdBytes()
	proto.DefenserId = f.targetSubTroop.startingBase.IdBytes()

	proto.AttackerFightAmount = f.attackerFightAmount
	proto.DefenserFightAmount = f.targetFightAmount

	proto.AttackerAliveSoldier = f.attackerAliveSoldier
	proto.DefenserAliveSoldier = f.targetAliveSoldier

	proto.AttackerTotalSoldier = f.attackerTotalSoldier
	proto.DefenserTotalSoldier = f.targetTotalSoldier

	proto.AttackerWin = f.response.AttackerWin
	// proto.Attacker 不设值，客户端自己设
	// proto.Defenser 不设值，客户端自己设

	proto.WinTimes = f.winTimes
	proto.Share = copyShare(f.response.DefenserShare)

	return proto
}

func (ctx *fightxContext) rollback() {
	for _, fst := range ctx.subTroopMap {
		for i, v := range fst.sub.captains {
			if i >= len(fst.beforeSoldier) {
				break
			}
			v.GetAndSetSoldierCount(u64.FromInt32(fst.beforeSoldier[i]))
		}
	}
}

func (ctx *fightxContext) addFightTroop(sub *troop) bool {
	fst := newFightSubtroop(sub)
	if fst != nil {
		ctx.subTroopMap[sub.Id()] = fst
	}
	return true
}

func (r *Realm) fightAssembly(attacker *troop, defenserBase *baseWithData, assister []*troop, target *troop, fightType fightType) (fightErr, isAttackerAlive, isDefenserAlive bool, report *shared_proto.FightReportProto) {

	ctx := newFightxContext()

	defenserTroopCount := len(assister)
	if target != nil {
		defenserTroopCount++
	}

	// 先记录下对战双方的兵力，如果后续发生错误，将士兵数回滚
	attacker.walkAll(ctx.addFightTroop)

	if target != nil {
		target.walkAll(ctx.addFightTroop)
	}

	for _, t := range assister {
		t.walkAll(ctx.addFightTroop)
	}

	fightErr, isAttackerAlive, isDefenserAlive = r.doAssemblyFight(ctx, attacker, defenserBase, assister, target, fightType)
	if fightErr {
		// 把兵加回去
		ctx.rollback()
		return
	}

	// 如果防守方没死光，进攻方也还有兵，则进攻方回家（走回去）

	// 防守方死光了，算赢
	isAttackerWin := !isDefenserAlive

	ctime := r.services.timeService.CurrentTime()

	// 飘字，需要处理去重
	troopLastFight := make(map[*troop]*battle)
	for _, f := range ctx.battles {
		troopLastFight[f.target] = f
	}

	for _, f := range troopLastFight {
		r.sendShowText(f.attacker, f.target, fightType, f.response.AttackerWin)
	}

	// 更新双方数据
	r.updateFightTroop(ctx, attacker, defenserBase, !isAttackerAlive, true, fightType, ctime)

	// 防守方没有集结部队，否则里面需要额外的处理
	var targetFightTroops []*troop
	for t := range troopLastFight {
		targetFightTroops = append(targetFightTroops, t)

		isRemove := !t.isAlive()
		r.updateFightTroop(ctx, t, defenserBase, isRemove, false, fightType, ctime)
	}

	report = r.newAssemblyReport(attacker, defenserBase, defenserTroopCount, fightType, isAttackerWin, ctx.battles)

	if len(ctx.battles) > 0 {
		last := ctx.battles[len(ctx.battles)-1]

		// 发战报邮件
		r.sendAssemblyReportMail(attacker, targetFightTroops, defenserBase, last, fightType, isAttackerWin, report, ctime)
	}

	for _, f := range ctx.battles {
		killer := f.targetSubTroop
		beenKilled := f.attackerSubTroop
		if f.response.AttackerWin {
			killer, beenKilled = beenKilled, killer
		}

		// 击杀匈奴入侵怪物，减少士气
		if base := GetXiongNuBase(beenKilled.startingBase); base != nil {

			isReduceMorale := beenKilled.State() != realmface.Defending && beenKilled.State() != realmface.Temp

			var toReduceMorale uint64
			if isReduceMorale {
				toReduceMorale = r.services.datas.ResistXiongNuMisc().WipeOutReduceMorale
			}

			// 减少怪物
			base.Info().WipeOutMonster(ctime, toReduceMorale)

			if isReduceMorale {
				// 降低士气
				r.onWipeOutXiongNuNpcTroop(attacker.startingBase, base, toReduceMorale)

				if data := r.services.datas.GuildLogHelp().WipeOutXiongNuTroop; data != nil {
					hero := r.services.heroSnapshotService.Get(f.target.Id())
					if hero != nil {
						proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
						proto.Text = data.Text.New().
							WithHeroName(hero.Name).
							WithResistXiongNuName(attacker.getStartingBaseInfo().getName()).
							WithReduceMorale(toReduceMorale).JsonString()
						r.services.guildService.AddLog(base.Info().GuildId(), proto)
					}
				}
			}
		}
	}

	if npcid.IsXiongNuNpcId(defenserBase.Id()) || npcid.IsJunTuanNpcId(defenserBase.Id()) {
		defenserBase.updateRoBase()
		defenserBase.ClearUpdateBaseInfoMsg()
		r.broadcastBaseInfoToCared(defenserBase, addBaseTypeUpdate, 0)
	}
	return
}

func (r *Realm) doAssemblyFight(ctx *fightxContext, attacker *troop, defenserBase *baseWithData, assister []*troop, target *troop, fightType fightType) (fightErr, isAttackerAlive, isDefenserAlive bool) {
	// 进攻的队伍
	astRandom := random.NewArray()
	attacker.walkAll(func(st *troop) (toContinue bool) {
		fst := ctx.subTroopMap[st.Id()]
		if fst != nil {
			astRandom.Add(fst)
		}
		return true
	})
	isAttackerAlive = true

	// 开始打架

	var tlogBattleType uint64
	if fightType == fightTypeInvadeArrive {
		tlogBattleType = operate_type.BattleInvade
	} else if fightType == fightTypeAssistArrive {
		tlogBattleType = operate_type.BattleAssist
	} else if fightType == fightTypeExpel {
		tlogBattleType = operate_type.BattleExpel
	}

	tfctx := entity.NewTlogFightContext(tlogBattleType, 0, 0, 0)

	var attackerWinTimesLimit int32
	if npcid.IsJunTuanNpcId(defenserBase.Id()) {
		attackerWinTimesLimit = u64.Int32(r.genConfig().JunTuanWinTimeLimit)
	}

	var hasWinTimesLimitTroop bool
	dstRandom := random.NewArray()
	if len(assister) > 0 {
		// 防守的队伍
		for _, t := range assister {
			t.walkAll(func(st *troop) (toContinue bool) {
				fst := ctx.subTroopMap[st.Id()]
				if fst != nil {
					dstRandom.Add(fst)
				}
				return true
			})
		}

		if fightErr, hasWinTimesLimitTroop = r.doAssemblyRandomFight(ctx, tfctx, astRandom, dstRandom, fightType, attackerWinTimesLimit); fightErr {
			return
		}

		isAttackerAlive = hasWinTimesLimitTroop || astRandom.Len() > 0
		isDefenserAlive = dstRandom.Len() > 0 || target != nil
		if astRandom.Len() <= 0 {
			// 进攻方死光了
			return
		}

		if dstRandom.Len() > 0 {
			// 可能有连胜限制，导致没有继续打下去，算进攻方输
			return
		}
	}

	// 跟援助队伍打完，如果有防守队伍，继续干防守队伍
	if target != nil {
		target.walkAll(func(st *troop) (toContinue bool) {
			fst := ctx.subTroopMap[st.Id()]
			if fst != nil {
				dstRandom.Add(fst)
			}
			return true
		})

		isDefenserAlive = false
		if dstRandom.Len() > 0 {
			prevHasWinTimesLimitTroop := hasWinTimesLimitTroop
			if fightErr, hasWinTimesLimitTroop = r.doAssemblyRandomFight(ctx, tfctx, astRandom, dstRandom, fightType, attackerWinTimesLimit); fightErr {
				return
			}

			isAttackerAlive = prevHasWinTimesLimitTroop || hasWinTimesLimitTroop || astRandom.Len() > 0
			isDefenserAlive = dstRandom.Len() > 0
		}
	}

	return
}

func (r *Realm) updateFightTroop(ctx *fightxContext, t *troop, defenserBase *baseWithData, isRemove, isInvadeTroop bool, fightType fightType, ctime time.Time) {

	t.walkAll(func(st *troop) (toContinue bool) {
		fst := ctx.subTroopMap[st.Id()]
		var fights []*battle
		if fst != nil {
			fights = fst.battles
		}

		if isRemove || len(fights) > 0 {
			r.updateHeroFightTroop(st, fights, defenserBase, isRemove, isInvadeTroop, fightType, ctime)
		}
		return true
	})

	// 删除失败的部队
	if isRemove {
		if t.State() != realmface.Temp {
			r.trySendTroopDoneMail(t, r.getTextHelp().MRDRExpelled4a.Text, r.getTextHelp().MRDRExpelled4d.Text, ctime)
		}

		t.walkAll(func(st *troop) (toContinue bool) {

			st.removeWithoutReturnCaptain(r)
			if st.State() != realmface.Temp {
				r.broadcastRemoveMaToCared(st)
			}
			st.clearMsgs()

			return true
		})
	} else {
		if t.startingBase != defenserBase {
			// 更新给关注的人
			r.broadcastMaToCared(t, addTroopTypeUpdate, 0)
		}
	}
}

func (r *Realm) updateHeroFightTroop(t *troop, fights []*battle, defenserBase *baseWithData, isRemove, isInvadeTroop bool, fightType fightType, ctime time.Time) {

	var defenserPrize []*resdata.PlunderPrize
	r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {

		// 更新士兵数
		if t.State() != realmface.Temp {
			if !isRemove {
				hero.UpdateTroop(t, true)
			} else {
				t.leaveTarget(hero, result, ctime)
				hero.RemoveTroop(t, true)
			}
		} else {
			if isInvadeTroop {
				logrus.WithField("fightType", fightType).WithField("stack", string(debug.Stack())).WithField("target", t.Id()).Error("fight时, attacker.state == temp")

				// 更新部队士兵
				hero.UpdateTroopSoldier(t)
			} else {
				//if fightType != fightTypeInvadeArrive && fightType != fightTypeExpel {
				//	logrus.WithField("fightType", fightType).WithField("stack", string(debug.Stack())).WithField("target", st.Id()).Error("fight时, target.state == temp, 但是fightType不是InvadeArrive || Expel")
				//}
				//
				//if attacker.targetBase != target.startingBase {
				//	logrus.WithField("attacker.targetBase", attacker.targetBase).WithField("target.startingBase", target.startingBase).WithField("stack", string(debug.Stack())).WithField("target",  st.Id()).Error("fight时, target.state == temp, 但是attacker.targetBase != target.startingBase")
				//}

				if t.Id() == 0 {
					// 镜像防守
					if isRemove {
						hero.RemoveCopyDefenser()
						result.Add(military.REMOVE_COPY_DEFENSER_S2C)
					} else {
						hero.UpdateCopyDefenser(t)
						result.Add(hero.NewUpdateCopyDefenserMsg())
					}
				} else {
					// 更新部队士兵
					hero.UpdateTroopSoldier(t)
				}

				// 任务相关
				if fightType == fightTypeExpel {
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_Expel)
					heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_EXPEL)
				}
			}
		}

		// 更新武将士兵数
		if t.Id() != 0 && len(t.captains) > 0 {
			result.AddFunc(func() pbutil.Buffer {
				proto := &region.S2CUpdateSelfTroopsProto{}
				for _, c := range t.captains {
					proto.Id = append(proto.Id, u64.Int32(c.Id()))
					proto.Soldier = append(proto.Soldier, c.Proto().Soldier)
					proto.FightAmount = append(proto.FightAmount, data.ProtoFightAmount(c.Proto().TotalStat, c.Proto().Soldier, c.Proto().SpellFightAmountCoef))
				}

				proto.RemoveOutside = isRemove
				if t.State() == realmface.Temp {
					proto.RemoveOutside = true
				}

				proto.TroopIndex = u64.Int32(entity.GetTroopIndex(t.Id()) + 1)

				return region.NewS2cUpdateSelfTroopsProtoMsg(proto)
			})
		}

		var updateTaskType []shared_proto.TaskTargetType
		for _, f := range fights {
			response := f.response

			otherSubTroop, otherIsAttacker := f.getOtherSubtroop(t)
			isAttacker := !otherIsAttacker

			var killSoldier int32
			if isAttacker {
				for _, v := range response.AttackerKillSoldier {
					killSoldier += v
				}
			} else {
				for _, v := range response.DefenserKillSoldier {
					killSoldier += v
				}
			}

			// 打的是玩家
			if !otherSubTroop.startingBase.isNpcBase() {
				if !isInvadeTroop {
					// 加仇人
					hero.Relation().AddEnemy(otherSubTroop.startingBase.Id(), ctime)
					result.Add(relation.NewS2cAddEnemyMsg(idbytes.ToBytes(otherSubTroop.startingBase.Id())))

					// TODO tlog 加仇人
				}

				if isAttacker == response.AttackerWin {
					// 战斗胜利次数
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RealmPvpSuccess)
				} else {
					// 战斗失败次数
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RealmPvpFail)
				}
			}

			if !isAttacker {
				if t.startingBase == defenserBase {
					// 防守消灭士兵数
					if killSoldier > 0 {
						hero.HistoryAmount().Increase(server_proto.HistoryAmountType_DefenseKillSoldier, u64.FromInt32(killSoldier))
						updateTaskType = append(updateTaskType, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_DEFENSE_KILL_SOLDIER)
					}
				} else {
					// 加援助次数
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_RealmPvpAssist)
					updateTaskType = append(updateTaskType, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_REALM_PVP_ASSIST)

					// 援助消灭士兵数
					if killSoldier > 0 {
						hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AssistKillSoldier, u64.FromInt32(killSoldier))
						updateTaskType = append(updateTaskType, shared_proto.TaskTargetType_TASK_TARGET_ASSIST_KILL_SOLDIER)
						updateTaskType = append(updateTaskType, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_ASSIST_KILL_SOLDIER)
					}
				}
			}

			// 击杀奖励
			if isAttacker == response.AttackerWin {
				if otherSubTroop.startingBase.isNpcBase() && otherSubTroop.mmd != nil {
					// 被杀的是个npc，给npc部队击杀奖励
					if otherSubTroop.mmd.BeenKillPrize != nil {
						prize := otherSubTroop.mmd.BeenKillPrize.GetPrize()
						//prizeProto := prize.Encode()

						// 击杀入侵奖励
						hctx := heromodule.NewContext(r.dep, operate_type.RealmKillDefNpc)
						heromodule.AddPrize(hctx, hero, result, prize, ctime)

						t.tryAddPrize(prize)
					}

					// 给城主加的
					if otherSubTroop.mmd.InvadePrize != nil {
						// 如果有击杀奖励，那么自己击杀的，不给入侵奖励
						if otherSubTroop.mmd.BeenKillPrize == nil || t.startingBase != defenserBase {
							defenserPrize = append(defenserPrize, otherSubTroop.mmd.InvadePrize)
						}
					}
				}
			}

			// 出征打Npc输了
			if fightType == fightTypeInvadeArrive && isAttacker && !response.AttackerWin {
				if defenseBaseInfo := otherSubTroop.getStartingBaseInfo(); defenseBaseInfo != nil {
					// 失败仇恨变化
					if hateData := defenseBaseInfo.getHateData(); hateData != nil {
						if info := hero.GetNpcTypeInfo(hateData.TypeData().Type); info != nil && hateData.FailHate != 0 {
							newHate := u64.AddInt(info.GetHate(), hateData.FailHate)
							newHate = u64.Min(newHate, hateData.TypeData().MaxHate)

							if info.SetHate(newHate) {
								result.Add(region.NewS2cUpdateMultiLevelNpcHateMsg(int32(hateData.TypeData().Type), u64.Int32(newHate)))
							}
						}
					}
				}
			}
		}

		for _, taskType := range updateTaskType {
			heromodule.UpdateTaskProgress(hero, result, taskType)
		}
	})

	if len(defenserPrize) > 0 {
		builder := resdata.NewPrizeBuilder()
		for _, p := range defenserPrize {
			builder.Add(p.GetPrize())
		}
		killInvadeTroopPrize := builder.Build()

		r.heroBaseFuncWithSend(defenserBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
			hctx := heromodule.NewContext(r.dep, operate_type.RealmKillInvadeNpc)
			heromodule.AddPrize(hctx, hero, result, killInvadeTroopPrize, ctime)
			//fctx.killInvadeTroopPrize = prizeProto
		})
	}
}

func (r *Realm) doAssemblyRandomFight(ctx *fightxContext, tfctx *entity.TlogFightContext, astRandom, dstRandom *random.Array, fightType fightType, attackerWinTimesLimit int32) (fightErr, hasWinTimeLimitTroop bool) {

	fightMaxTimes := astRandom.Len() + dstRandom.Len() - 1
	for i := 0; i < fightMaxTimes; i++ {
		// 进攻方随机一个部队出来
		aidx, afst := randomFightSubtroop(astRandom.Random())
		if afst == nil {
			// wtf
			logrus.Error("野外战斗，随机出来的 ast == nil")
			fightErr = true
			return
		}

		// 防守方随机一个部队出来
		didx, dfst := randomFightSubtroop(dstRandom.Random())
		if dfst == nil {
			// wtf
			logrus.Error("野外战斗，随机出来的 dst == nil")
			fightErr = true
			return
		}

		ast := afst.sub
		dst := dfst.sub
		acp := ast.toCombatPlayerProto(r)
		dcp := dst.toCombatPlayerProto(r)

		afa := acp.TotalFightAmount
		dfa := dcp.TotalFightAmount

		var attackerTotalSoldier, defenserTotalSoldier int32
		for _, t := range acp.Troops {
			attackerTotalSoldier += t.Captain.Soldier
		}

		for _, t := range dcp.Troops {
			defenserTotalSoldier += t.Captain.Soldier
		}

		// 2个subtroop进行战斗
		response := r.services.fightModule.SendFightRequest(tfctx, r.config().CombatScene,
			ast.startingBase.Id(), dst.startingBase.Id(), acp, dcp)

		if response.ReturnCode != 0 {
			logrus.WithField("troopid", ast.Id()).WithField("target_troopid", dst.Id()).WithField("fight_type", fightType).WithField("return_msg", response.ReturnMsg).WithField("return_code", response.ReturnCode).Error("战斗计算发生错误")
			fightErr = true
			return
		}

		var attackerAliveSoldier, defenserAliveSoldier int32
		for k, v := range response.AttackerAliveSoldier {
			if k != 0 {
				attackerAliveSoldier += v
			}
		}
		for k, v := range response.DefenserAliveSoldier {
			if k != 0 {
				defenserAliveSoldier += v
			}
		}

		logrus.WithField("attacker", ast.startingBase.Id()).
			WithField("defenser", dst.startingBase.Id()).
			WithField("win", response.AttackerWin).
			WithField("a_total_soldier", attackerTotalSoldier).
			WithField("d_total_soldier", defenserTotalSoldier).
			WithField("a_alive_soldier", attackerAliveSoldier).
			WithField("d_alive_soldier", defenserAliveSoldier).
			Debug("集结战斗打架")

		var winTimes uint64
		if response.AttackerWin {
			afst.winTimes++
			winTimes = afst.winTimes
		} else {
			dfst.winTimes++
			winTimes = dfst.winTimes
		}

		fight := newBattle(ast, dst, afa, dfa, attackerAliveSoldier, defenserAliveSoldier,
			attackerTotalSoldier, defenserTotalSoldier, response, int32(winTimes))
		ctx.battles = append(ctx.battles, fight)

		afst.battles = append(afst.battles, fight)
		dfst.battles = append(dfst.battles, fight)

		if response.AttackerWin {
			// 进攻方胜利，删除防守方部队
			dstRandom.Remove(didx)

			reduceSoldierToAlive(ast, response.AttackerAliveSoldier, 0)
			reduceSoldierToZero(dst, 0)

			// 如果达到连胜上限，则不再参与战斗
			if attackerWinTimesLimit > 0 && fight.winTimes >= attackerWinTimesLimit {
				astRandom.Remove(aidx)

				hasWinTimeLimitTroop = true
				if astRandom.Len() <= 0 {
					// 进攻方没有战斗队伍了
					return
				}
			}

			if dstRandom.Len() <= 0 {
				// 防守方的援助队伍都死光了
				return
			}
		} else {
			// 进攻失败
			astRandom.Remove(aidx)

			reduceSoldierToZero(ast, 0)
			reduceSoldierToAlive(dst, response.DefenserAliveSoldier, 0)

			if astRandom.Len() <= 0 {
				// 进攻方全部死光了
				return
			}
		}
	}

	return
}

func randomFightSubtroop(idx int, value interface{}) (int, *fightxTroop) {
	if value == nil {
		return idx, nil
	}
	return idx, value.(*fightxTroop)
}

func (r *Realm) sendAssemblyReportMail(attacker *troop, fightTroops []*troop, defenserBase *baseWithData, last *battle, fightType fightType, isAttackerWin bool, report *shared_proto.FightReportProto, ctime time.Time) {

	// 对所有进行战斗的人，发送战报消息

	// 构建对所有人有效的战报数据

	// 区分不同的角色，获取战报信息，

	attackerId := attacker.startingBase.Id()
	defenserId := defenserBase.Id()

	attackerFlagName := attacker.getStartingBaseFlagHeroName(r)
	defenserFlagName := r.toBaseFlagHeroName(defenserBase, attacker.targetBaseLevel)

	mailTag := getMailTagByTroop(attacker, last.target)

	switch fightType {
	case fightTypeAssistArrive:
		assister := last.target
		assisterFlagName := assister.getStartingBaseFlagHeroName(r)

		if isAttackerWin {
			// 进攻方胜利

			// 进攻方成功
			if data := r.getMailHelp().AssemblySaaSuccess; data != nil {
				//report.ShowPrize = killDefPrize
				attacker.walkAll(func(st *troop) (toContinue bool) {
					r.sendHeroReportMail(st.startingBase.Id(), data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, true, 0, ctime, mailTag, shared_proto.MailType_MailAssemblyReport)
					return true
				})
			}

			// 援助方失败
			if data := r.getMailHelp().AssemblySasFail; data != nil {
				assister.walkAll(func(st *troop) (toContinue bool) {
					r.sendHeroReportMail(st.startingBase.Id(), data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, false, 0, ctime, mailTag, shared_proto.MailType_MailAssemblyReport)
					return true
				})
			}

			// 城主失败
			if data := r.getMailHelp().AssemblySadFail; data != nil {
				r.sendHeroReportMail(defenserId, data,
					attackerFlagName, defenserFlagName, assisterFlagName,
					report, false, 0, ctime, mailTag, shared_proto.MailType_MailAssemblyReport)
			}

		} else {
			// 进攻方失败

			// 进攻方失败
			if data := r.getMailHelp().AssemblySaaFail; data != nil {
				attacker.walkAll(func(st *troop) (toContinue bool) {
					r.sendHeroReportMail(st.startingBase.Id(), data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, true, 0, ctime,
						mailTag, shared_proto.MailType_MailAssemblyReport)
					return true
				})
			}

			// 援助方成功
			if data := r.getMailHelp().AssemblySasSuccess; data != nil {
				assister.walkAll(func(st *troop) (toContinue bool) {
					r.sendHeroReportMail(st.startingBase.Id(), data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, false, 0, ctime,
						mailTag, shared_proto.MailType_MailAssemblyReport)
					return true
				})
				//report.ShowPrize = killDefPrize
			}

			// 城主成功
			if data := r.getMailHelp().AssemblySadSuccess; data != nil {
				//report.ShowPrize = fightNpcPrize

				r.sendHeroReportMail(defenserId, data,
					attackerFlagName, defenserFlagName, assisterFlagName,
					report, false, 0, ctime,
					mailTag, shared_proto.MailType_MailAssemblyReport)
			}
		}

	case fightTypeInvadeArrive:

		// 进攻到达
		if isAttackerWin {

			// 掠夺者干掉了盟友
			// 进攻方胜利 在外面再发消息
			//if data := r.getMailHelp().AssemblyAdaSuccess; data != nil {
			//	//report.ShowPrize = killDefPrize
			//	attacker.walkAll(func(st *troop) (toContinue bool) {
			//		r.sendHeroReportMail(st.startingBase.Id(), data,
			//			attackerFlagName, defenserFlagName, "",
			//			report, true, 0, ctime,
			//			mailTag, shared_proto.MailType_MailAssemblyReport)
			//		return true
			//	})
			//}

			// 援助方失败
			if data := r.getMailHelp().AssemblyAssFail; data != nil {
				for _, assister := range fightTroops {
					if assister.startingBase == defenserBase {
						continue
					}
					assisterFlagName := assister.getStartingBaseFlagHeroName(r)
					assister.walkAll(func(st *troop) (toContinue bool) {
						r.sendHeroReportMail(st.startingBase.Id(), data,
							attackerFlagName, defenserFlagName, assisterFlagName,
							report, false, 0, ctime,
							mailTag, shared_proto.MailType_MailAssemblyReport)
						return true
					})
				}
			}

			// 城主失败
			//if data := r.getMailHelp().AssemblyAddFail; data != nil {
			//	r.sendHeroReportMail(defenserId, data,
			//		attackerFlagName, defenserFlagName, "",
			//		report, false, 0, ctime,
			//		mailTag, shared_proto.MailType_MailAssemblyReport)
			//}

		} else {

			attacker := last.attacker
			target := last.target

			if target.startingBase == defenserBase {
				// 防守阵容干掉的

				// 进攻方失败
				if data := r.getMailHelp().AssemblyAdaFail; data != nil {
					attacker.walkAll(func(st *troop) (toContinue bool) {
						report.ShowPrize = st.clearAndReturnAccumRobPrize()
						r.sendHeroReportMail(st.startingBase.Id(), data,
							attackerFlagName, defenserFlagName, "",
							report, true, 0, ctime,
							mailTag, shared_proto.MailType_MailAssemblyReport)
						report.ShowPrize = nil
						return true
					})
				}

				// 援助方成功
				if data := r.getMailHelp().AssemblyAssSuccess; data != nil {
					//report.ShowPrize = killDefPrize

					for _, assister := range fightTroops {
						if assister.startingBase == defenserBase {
							continue
						}
						assisterFlagName := assister.getStartingBaseFlagHeroName(r)
						assister.walkAll(func(st *troop) (toContinue bool) {
							r.sendHeroReportMail(st.startingBase.Id(), data,
								attackerFlagName, defenserFlagName, assisterFlagName,
								report, false, 0, ctime,
								mailTag, shared_proto.MailType_MailAssemblyReport)
							return true
						})
					}
				}

				// 城主成功
				if data := r.getMailHelp().AssemblyAddSuccess; data != nil {
					//report.ShowPrize = fightNpcPrize
					//if report.ShowPrize == nil {
					//	report.ShowPrize = killDefPrize
					//}

					r.sendHeroReportMail(defenserId, data,
						attackerFlagName, defenserFlagName, "",
						report, false, 0, ctime,
						mailTag, shared_proto.MailType_MailAssemblyReport)
				}
			} else {
				// 盟友干掉的
				assisterFlagName := target.getStartingBaseFlagHeroName(r)

				// 进攻方失败
				if data := r.getMailHelp().AssemblyAsaFail; data != nil {
					attacker.walkAll(func(st *troop) (toContinue bool) {
						report.ShowPrize = st.clearAndReturnAccumRobPrize()
						r.sendHeroReportMail(attackerId, data,
							attackerFlagName, defenserFlagName, assisterFlagName,
							report, true, 0, ctime,
							mailTag, shared_proto.MailType_MailAssemblyReport)
						report.ShowPrize = nil
						return true
					})
				}

				// 援助方成功
				if data := r.getMailHelp().AssemblyAssSuccess; data != nil {
					//report.ShowPrize = killDefPrize

					for _, assister := range fightTroops {
						if assister.startingBase == defenserBase {
							continue
						}
						assisterFlagName := assister.getStartingBaseFlagHeroName(r)
						assister.walkAll(func(st *troop) (toContinue bool) {
							assisterId := st.startingBase.Id()
							r.sendHeroReportMail(assisterId, data,
								attackerFlagName, defenserFlagName, assisterFlagName,
								report, false, 0, ctime,
								mailTag, shared_proto.MailType_MailAssemblyReport)
							return true
						})
					}
				}

				// 城主成功
				if data := r.getMailHelp().AssemblyAsdSuccess; data != nil {
					//report.ShowPrize = fightNpcPrize

					r.sendHeroReportMail(defenserId, data,
						attackerFlagName, defenserFlagName, assisterFlagName,
						report, false, 0, ctime,
						mailTag, shared_proto.MailType_MailAssemblyReport)
				}
			}
		}

	case fightTypeExpel:

		if isAttackerWin {
			// 驱逐失败
			// 进攻方成功
			if data := r.getMailHelp().AssemblyExpelAttackerSuccess; data != nil {
				//report.ShowPrize = killDefPrize
				attacker.walkAll(func(st *troop) (toContinue bool) {
					r.sendHeroReportMail(st.startingBase.Id(), data,
						attackerFlagName, defenserFlagName, "",
						report, true, 0, ctime,
						mailTag, shared_proto.MailType_MailAssemblyReport)
					return true
				})
			}

			// 城主失败
			if data := r.getMailHelp().AssemblyExpelDefenserFail; data != nil {
				r.sendHeroReportMail(defenserId, data,
					attackerFlagName, defenserFlagName, "",
					report, false, 0, ctime,
					mailTag, shared_proto.MailType_MailAssemblyReport)
			}
		} else {
			// 驱逐成功
			// 进攻方失败
			if data := r.getMailHelp().AssemblyExpelAttackerFail; data != nil {
				attacker.walkAll(func(st *troop) (toContinue bool) {
					r.sendHeroReportMail(st.startingBase.Id(), data,
						attackerFlagName, defenserFlagName, "",
						report, true, 0, ctime,
						mailTag, shared_proto.MailType_MailAssemblyReport)
					return true
				})

			}

			// 城主成功
			if data := r.getMailHelp().AssemblyExpelDefenserSuccess; data != nil {
				//report.ShowPrize = fightNpcPrize
				//if report.ShowPrize == nil {
				//	report.ShowPrize = killDefPrize
				//}

				r.sendHeroReportMail(defenserId, data,
					attackerFlagName, defenserFlagName, "",
					report, false, 0, ctime,
					mailTag, shared_proto.MailType_MailAssemblyReport)
			}
		}
	}
}

func (r *Realm) sendShowText(attacker, target *troop, fightType fightType, attackerWin bool) {

	// 处理战斗过程飘字
	switch fightType {
	case fightTypeAssistArrive:

		assister := attacker
		robber := target

		if attackerWin {

			// 援助成功飘字
			assister.walkAll(func(st *troop) (toContinue bool) {
				r.services.world.SendFunc(st.startingBase.Id(), func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().RealmAssistSuccess.New().
							WithTroopIndex(entity.GetTroopIndex(st.Id()) + 1).
							JsonString())
				})
				return true
			})

			// 掠夺者失败飘字
			robber.walkAll(func(st *troop) (toContinue bool) {
				r.services.world.SendFunc(st.startingBase.Id(), func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().BannerTroopFail.New().
							WithAttacker(assister.getStartingBaseFlagHeroName(r)).
							WithTroopIndex(entity.GetTroopIndex(st.Id()) + 1).
							JsonString())
				})
				return true
			})

		} else {
			// 援助失败飘字
			assister.walkAll(func(st *troop) (toContinue bool) {
				r.services.world.SendFunc(st.startingBase.Id(), func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().RealmAssistFail.New().
							WithTroopIndex(entity.GetTroopIndex(st.Id()) + 1).
							JsonString())
				})
				return true
			})

			// 掠夺者胜利飘字
			robber.walkAll(func(st *troop) (toContinue bool) {
				r.services.world.SendFunc(st.startingBase.Id(), func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().BannerTroopSuccess.New().
							WithAttacker(assister.getStartingBaseFlagHeroName(r)).
							WithTroopIndex(entity.GetTroopIndex(st.Id()) + 1).
							JsonString())
				})
				return true
			})
		}

	case fightTypeInvadeArrive:
		robber := attacker

		if !attackerWin {
			// 出征失败飘字
			robber.walkAll(func(st *troop) (toContinue bool) {
				r.services.world.SendFunc(st.startingBase.Id(), func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().RealmInvadeFail.New().
							WithTroopIndex(entity.GetTroopIndex(st.Id()) + 1).
							JsonString())
				})
				return true
			})

			if robber.targetBase != target.startingBase {
				// 防守的援助部队胜利飘字
				target.walkAll(func(st *troop) (toContinue bool) {
					r.services.world.SendFunc(st.startingBase.Id(), func() pbutil.Buffer {
						return misc.NewS2cScreenShowWordsMsg(
							r.getTextHelp().BannerTroopSuccess.New().
								WithAttacker(attacker.getStartingBaseFlagHeroName(r)).
								WithTroopIndex(entity.GetTroopIndex(st.Id()) + 1).
								JsonString())
					})
					return true
				})
			} else {
				if d := r.getTextHelp().BannerBaseSuccess; d != nil {
					defenserId := attacker.originTargetId
					r.services.world.Send(defenserId, misc.NewS2cScreenShowWordsMsg(d.Text.New().WithAttacker(attacker.getStartingBaseFlagHeroName(r)).JsonString()))
				}
			}

		} else {
			if robber.targetBase != target.startingBase {
				// 防守的援助部队失败飘字
				target.walkAll(func(st *troop) (toContinue bool) {
					r.services.world.SendFunc(st.startingBase.Id(), func() pbutil.Buffer {
						return misc.NewS2cScreenShowWordsMsg(
							r.getTextHelp().BannerTroopFail.New().
								WithAttacker(attacker.getStartingBaseFlagHeroName(r)).
								WithTroopIndex(entity.GetTroopIndex(st.Id()) + 1).
								JsonString())
					})
					return true
				})
			}
		}
	case fightTypeExpel:

		robbing := attacker
		if !attackerWin {
			// 遭受驱逐，战斗失败
			robbing.walkAll(func(st *troop) (toContinue bool) {
				r.services.world.SendFunc(st.startingBase.Id(), func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().BannerTroopFail.New().
							WithAttacker(target.getStartingBaseFlagHeroName(r)).
							WithTroopIndex(entity.GetTroopIndex(st.Id()) + 1).
							JsonString())
				})
				return true
			})

		} else {
			// 遭受驱逐，战斗胜利
			robbing.walkAll(func(st *troop) (toContinue bool) {
				r.services.world.SendFunc(st.startingBase.Id(), func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().BannerTroopSuccess.New().
							WithAttacker(target.getStartingBaseFlagHeroName(r)).
							WithTroopIndex(entity.GetTroopIndex(st.Id()) + 1).
							JsonString())
				})
				return true
			})
		}
	}
}

func (r *Realm) newAssemblyReport(attacker *troop, defenserBase *baseWithData, defenserTroopCount int,
	fightType fightType, isAttackerWin bool, fight []*battle) *shared_proto.FightReportProto {
	proto := &shared_proto.FightReportProto{}

	proto.FightX = int32(defenserBase.BaseX())
	proto.FightY = int32(defenserBase.BaseY())

	proto.FightType = int32(fightType)

	proto.AttackerWin = isAttackerWin
	proto.Attacker = attacker.toReportHeroProto(r, nil, nil)
	proto.Defenser = defenserBase.toReportHeroProto(r, attacker.targetBaseLevel)

	baseIdMap := make(map[int64]struct{})
	baseIdMap[attacker.startingBase.Id()] = struct{}{}
	baseIdMap[defenserBase.Id()] = struct{}{}
	for _, f := range fight {
		proto.Fight = append(proto.Fight, f.toFightProto())

		if f.response.AttackerWin {
			proto.AttackerTroopWinTimes++
		} else {
			proto.DefenserTroopWinTimes++
		}

		base := f.attackerSubTroop.startingBase
		if _, exist := baseIdMap[base.Id()]; !exist {
			proto.FightHero = append(proto.FightHero, base.toReportHeroProto(r, base.BaseLevel()))
		}

		base = f.targetSubTroop.startingBase
		if _, exist := baseIdMap[base.Id()]; !exist {
			proto.FightHero = append(proto.FightHero, base.toReportHeroProto(r, base.BaseLevel()))
		}
	}

	if as := attacker.GetAssembly(); as != nil {
		proto.AttackerTroopCount = int32(as.Count())
		proto.AttackerTroopTotalCount = u64.Int32(as.TotalCount())
	}

	proto.DefenserTroopCount = int32(defenserTroopCount)
	proto.DefenserTroopTotalCount = proto.DefenserTroopCount

	return proto
}

//func reduceSubtroopSoldierToZero(st *troop) {
//
//	for _, captain := range st.captains {
//		captain.GetAndSetSoldierCount(0)
//	}
//
//	st.t.onChanged()
//	return
//}
//
//func reduceSubtroopSoldierToAlive(st *troop, aliveSoldierMap map[int32]int32) {
//
//	for _, captain := range st.captains {
//		if alive, has := aliveSoldierMap[int32(captain.Id())]; has && alive > 0 {
//			oldCount := captain.GetAndSetSoldierCount(uint64(alive))
//
//			if oldCount < uint64(alive) {
//				logrus.WithField("captain", captain).WithField("old_count", oldCount).WithField("new_count", alive).Error("打完一场后, 活着的士兵竟然超过了开打时的士兵数")
//				captain.GetAndSetSoldierCount(oldCount)
//			}
//		} else {
//			// 没找到, 当做全死光了
//			captain.GetAndSetSoldierCount(0)
//		}
//	}
//
//	st.t.onChanged()
//	return
//}
