package combatx

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/logrus"
	"math"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/pkg/errors"
	"sort"
)

const Denominator float64 = 10000

func NewCombat(p *server_proto.CombatXRequestServerProto, config *Config) (*Combat, error) {

	if p == nil {
		return nil, errors.Errorf("请求参数无效，传入的CombatRequestProto为nil")
	}

	if p.Attacker == nil {
		return nil, errors.Errorf("请求参数无效，传入的request.Attacker为nil")
	}

	if len(p.Attacker.Troops) == 0 {
		return nil, errors.Errorf("请求参数无效，传入的request.Attacker.Troops == 0")
	}

	if p.Defenser == nil {
		return nil, errors.Errorf("请求参数无效，传入的request.Defenser为nil")
	}

	if len(p.Defenser.Troops) == 0 && p.Defenser.WallStat == nil {
		return nil, errors.Errorf("请求参数无效，传入的request.Defenser.Troops == 0")
	}

	if p.GetMapXLen() <= 0 || p.GetMapXLen() > math.MaxInt16 {
		return nil, errors.Errorf("请求参数无效，传入的request.MapXLen should in [0, %v], v: %v", math.MaxInt16, p.GetMapXLen())
	}
	if p.GetMapYLen() <= 0 || p.GetMapYLen() > math.MaxInt16 {
		return nil, errors.Errorf("请求参数无效，传入的request.MapYLen should in [0, %v], v: %v", math.MaxInt16, p.GetMapYLen())
	}

	misc := newMisc(p)

	var idCounter int32
	genId := func() int32 {
		idCounter++
		return idCounter
	}

	attacker, err := newCombatPlayer(genId, p.GetAttackerId(), true, p.Attacker, config, misc)
	if err != nil {
		return nil, errors.Wrap(err, "Combat.newCombatPlayer(attacker)失败")
	}

	defender, err := newCombatPlayer(genId, p.GetDefenserId(), false, p.Defenser, config, misc)
	if err != nil {
		return nil, errors.Wrap(err, "Combat.newCombatPlayer(attacker)失败")
	}

	return newCombat(config, misc, attacker, defender), nil
}

func newCombat(config *Config, misc *misc, attacker, defender *CombatPlayer) *Combat {

	if misc.isDebug {
		for _, t := range attacker.troops {
			logrus.Debugf("进攻方武将 %v %v %v %v %v %v %v", t.index, t.initProto.X, t.initProto.Y, t.soldier, t.proto.FightIndex, t.raceData.proto.Race, t.proto.Captain.CaptainId)
		}

		for _, t := range defender.troops {
			logrus.Debugf("防守方武将 %v %v %v %v %v %v %v", t.index, t.initProto.X, t.initProto.Y, t.soldier, t.proto.FightIndex, t.raceData.proto.Race, t.proto.Captain.CaptainId)
		}
	}

	c := &Combat{}
	c.Config = config
	c.misc = misc
	c.attacker = attacker
	c.defender = defender

	// 初始化对线目标
	for _, t := range attacker.troops {
		t.rushTarget, t.rushUpDownTarget = getRushTarget(t.proto.FightIndex, defender)
	}

	for _, t := range defender.troops {
		t.rushTarget, t.rushUpDownTarget = getRushTarget(t.proto.FightIndex, attacker)
	}

	actionTroops := make([]*Troops, 0, len(attacker.troops)+len(defender.troops))
	actionTroops = append(actionTroops, attacker.troops...)
	actionTroops = append(actionTroops, defender.troops...)
	// 按照行动顺序排序
	sort.Sort(TroopsSpeedSlice(actionTroops))
	c.actionTroops = actionTroops

	for _, t := range c.actionTroops {
		if !t.isAlive() {
			continue
		}

		t.resetPosition()

		ts := c.getCombatPlayer(t.isAttacker).troopRaceMap[t.getRace()]
		ts = append(ts, t)
		c.getCombatPlayer(t.isAttacker).troopRaceMap[t.getRace()] = ts
	}

	return c
}

func getRushTarget(fightIndex int32, otherHero *CombatPlayer) (rushTarget *Troops, rushUpDownTarget []*Troops) {
	if fightIndex <= 0 {
		return
	}

	// 初始化对线目标
	for _, dt := range otherHero.troops {
		if fightIndex == dt.proto.FightIndex {
			rushTarget = dt
			break
		}
	}

	if rushTarget != nil {
		for _, dt := range otherHero.troops {
			if diff := fightIndex - dt.proto.FightIndex; diff == 1 || diff == -1 {
				rushUpDownTarget = append(rushUpDownTarget, dt)
			}
		}
	}

	return
}

type Combat struct {
	*Config
	*misc

	// 战斗的双方
	attacker *CombatPlayer
	defender *CombatPlayer

	// 行动顺序，就按顺序吧
	actionTroops []*Troops

	triggerSpell []*TriggerSpell

	attackerWin bool
}

func (c *Combat) Calculate() *shared_proto.CombatXProto {
	aliveSoliderWhenStart := c.getAliveSolider()

	proto := &shared_proto.CombatXProto{}

	frameProto := &shared_proto.CombatFrameProto{}

	// 处理战前技能
	isCheckEnd := false

	for _, t := range c.actionTroops {
		if !t.isAlive() {
			continue
		}

		for _, s := range t.spellList.data.beginReleaseSpell {
			if s.spell != nil {
				// 根据技能找施法对象
				targets := t.getSpellTarget(c, s.spell, nil)
				if len(targets) <= 0 {
					continue
				}

				shouldCheckEnd := c.processReleaseSpell(t, targets, s.spell, 0,
					ReleaseTypeRelease, false, frameProto)
				isCheckEnd = isCheckEnd || shouldCheckEnd
			}
		}
	}

	if len(frameProto.Action) > 0 {
		proto.Frame = append(proto.Frame, frameProto)
		frameProto = nil
	}

	if !isCheckEnd || (!c.attacker.isAllDead() && !c.defender.isAllDead()) {
		// 开始行动

	out:
		for i := 0; i < c.MaxFrame; i++ {
			currentFrame := i + 1
			proto.MaxFrame = int32(currentFrame)

			if frameProto == nil {
				frameProto = &shared_proto.CombatFrameProto{}
			}
			frameProto.Frame = int32(currentFrame)

			// 部队行动
			hasTroopAlive := false
			for _, t := range c.actionTroops {
				if !t.isAlive() {
					continue
				}
				hasTroopAlive = true

				isCheckEnd, isCheckEnemy := t.Update(c, currentFrame, frameProto)
				if isCheckEnd {
					beenCheckerIsAttacker := t.isAttacker
					if isCheckEnemy {
						beenCheckerIsAttacker = !beenCheckerIsAttacker
					}

					beenChecker := c.getCombatPlayer(beenCheckerIsAttacker)
					if beenChecker.isAllDead() {
						// 敌人全死了
						c.attackerWin = !beenChecker.isAttacker

						// 添加frame
						proto.Frame = append(proto.Frame, frameProto)
						break out
					}
				}
			}

			if !hasTroopAlive {
				// 无部队存活，当防守方胜利
				logrus.Errorf("hasTroopAlive = false")
				c.attackerWin = false

				if len(frameProto.Action) > 0 {
					proto.Frame = append(proto.Frame, frameProto)
				}
				break out
			}

			// 城墙行动
			if checkEnd, isCheckAttacker := c.tryWallAction(currentFrame, frameProto); checkEnd {
				if isCheckAttacker {
					if c.attacker.isAllDead() {
						// 敌人全死了
						c.attackerWin = false
						// 添加frame
						proto.Frame = append(proto.Frame, frameProto)
						break out
					}
				} else {
					if c.defender.isAllDead() {
						// 敌人全死了
						c.attackerWin = true
						// 添加frame
						proto.Frame = append(proto.Frame, frameProto)
						break out
					}
				}

			}

			if len(frameProto.Action) > 0 {
				proto.Frame = append(proto.Frame, frameProto)
				frameProto = nil
			}
		}
	}

	// 防守方挂了，说明进攻胜利
	proto.Attacker = c.attacker.proto
	proto.AttackerTroopData = c.attacker.TroopInitData()
	proto.Defenser = c.defender.proto
	proto.DefenserTroopData = c.defender.TroopInitData()
	proto.MapRes = c.mapRes
	proto.MapXLen = c.request.MapXLen
	proto.MapYLen = c.request.MapYLen
	proto.CombatSolider = aliveSoliderWhenStart
	proto.Score = c.calcScore()
	proto.AttackerWin = c.isAttackerWin()
	proto.AliveSolider = c.getAliveSolider()
	proto.KillSolider = c.getKillSolider()

	return proto
}

func (c *Combat) isAttackerWin() (attackerWin bool) {
	return c.attackerWin
}

func (c *Combat) getWinner() (*CombatPlayer) {
	if c.attackerWin {
		return c.attacker
	} else {
		return c.defender
	}
}

func (c *Combat) getAliveSolider() (aliveSolider []*shared_proto.Int32Pair) {

	if w := c.defender.Wall(); w != nil {
		if w.life > 0 {
			aliveSolider = append(aliveSolider, &shared_proto.Int32Pair{
				Key:   0,
				Value: int32(w.life),
			})
		}
	}

	// 胜利者产生，结算
	for _, t := range c.actionTroops {
		if t.soldier <= 0 {
			continue
		}

		aliveSolider = append(aliveSolider, &shared_proto.Int32Pair{
			Key:   t.index,
			Value: int32(t.soldier),
		})
	}

	return
}

func (c *Combat) getKillSolider() (killSoldier []*shared_proto.Int32Pair) {
	if w := c.defender.Wall(); w != nil {
		if w.killSoldier > 0 {
			killSoldier = append(killSoldier, &shared_proto.Int32Pair{
				Key:   0,
				Value: int32(w.killSoldier),
			})
		}
	}

	// 胜利者产生，结算
	for _, t := range c.actionTroops {
		if t.killSoldier <= 0 {
			continue
		}

		killSoldier = append(killSoldier, &shared_proto.Int32Pair{
			Key:   t.index,
			Value: int32(t.killSoldier),
		})
	}

	return
}

func (c *Combat) calcScore() (score int32) {
	// 计算评分
	winner := c.getWinner()

	aliveSoldier := 0
	totalSoldier := 0

	for _, v := range winner.troops {
		aliveSoldier += v.soldier
		totalSoldier += int(v.proto.Captain.Soldier)
	}

	// 积分
	if totalSoldier <= 0 {
		return
	}

	percent := int32(aliveSoldier * 100 / totalSoldier)
	score = c.getScore(percent)

	logrus.Debugf("计算评分，alive: %d total: %d percent: %d score: %d", aliveSoldier, totalSoldier, percent, score)
	return
}

func (c *Combat) getOtherCombatPlayer(isAttacker bool) *CombatPlayer {
	return c.getCombatPlayer(!isAttacker)
}

func (c *Combat) getCombatPlayer(isAttacker bool) *CombatPlayer {
	if isAttacker {
		return c.attacker
	} else {
		return c.defender
	}
}

func (c *Combat) isEnemyAllDead(isAttacker bool) (allDead bool) {
	return c.getEnemy(isAttacker).isAllDead()
}

func (c *Combat) getEnemyTroops(isAttacker bool) []*Troops {
	return c.getEnemy(isAttacker).troops
}

func (c *Combat) getEnemy(isAttacker bool) *CombatPlayer {
	if isAttacker {
		return c.defender
	} else {
		return c.attacker
	}
}

func (c *Troops) updateState(combat *Combat, frame int, proto *shared_proto.CombatFrameProto) (moveSpeedChanged, rageRecoverChanged bool) {

	// 尝试移除过期的State
	for _, s := range c.stateMap {
		if frame >= s.getNextTickFrame() {
			s.tick()

			if s.tickDamage > 0 || s.data.Rage > 0 {

				effectProto := &shared_proto.TroopTickEffectProto{}
				effectProto.Id = s.data.Id()

				if s.data.Rage > 0 {
					c.addRage(s.data.Rage)
				}

				if s.tickDamage > 0 {
					c.MarkRage()

					changeShield, changeSoldier, showChangeSoldier, aliveSoldier, removeStateIds := combat.tryReduceDamage(c, s.tickDamage, frame, proto)
					s.caster.addKillSoldier(changeSoldier)

					effectProto.ChangeShield = int32(changeShield)
					effectProto.RemoveState = removeStateIds
					effectProto.ChangeSoldier = int32(IMax(showChangeSoldier, 1))
					effectProto.Soldier = int32(aliveSoldier)
					effectProto.Rage = c.newRageUpdateProtoIfChanged()
				}

				if effectProto.Rage == nil && s.data.Rage > 0 {
					effectProto.Rage = c.newRageUpdateProto()
				}

				addStateEffectAction(proto, c.index, effectProto)

				if !c.isAlive() {
					// 挂了...
					return
				}
			}

		}

		if frame >= s.endFrame {
			// 状态结束了，移除掉
			c.removeState(s.data)
		}
	}

	if c.isAttackSpeedChanged {
		c.isAttackSpeedChanged = false
		c.updateAttackSpeed(combat.FramePerSecond, combat.MinAttackFrame, combat.MaxAttackFrame)
	}

	if c.isMoveSpeedChanged {
		c.isMoveSpeedChanged = false
		c.updateMoveSpeed(combat.FramePerSecond, combat.MinMoveSpeedPerFrame, combat.MaxMoveSpeedPerFrame)

		// 更新当前移动路径
		moveSpeedChanged = true
	}

	if c.isStatChanged {
		c.isStatChanged = false
		c.updateTotalStat(combat.MinStat)
	}

	if c.isRageRecoverChanged {
		c.isRageRecoverChanged = false
		c.updateRageRecover(combat.RageRecoverSpeed, combat.FramePerSecond)

		// 更新怒气
		rageRecoverChanged = true
	}

	if c.hasStepToMove() && (c.isUnmovable() || c.isStun()) {
		// 不可行走
		c.stopMove(frame)
	}

	if c.isNotAttackable() || c.isStun() {
		c.spellList.baseSpell.nextReleaseFrame = 0
	}

	return
}

func (combat *Combat) tryReduceDamage(target *Troops, damage, frame int, frameProto *shared_proto.CombatFrameProto) (changeShield, changeSoldier, showChangeSoldier, aliveSoldier int, removeStateIds []int32) {

	// 掉血，先扣护盾，然后再扣血
	changeShield, removeStateIds = combat.tryReduceShield(target, damage)
	damage -= changeShield

	if damage > 0 {
		lifePercent := target.getLifePercent()

		changeSoldier = target.reduceLife(damage)
		showChangeSoldier = changeSoldier
		aliveSoldier = target.getSoldier()

		var reliveSpell *PassiveSpellData
		if !target.isAlive() {
			// 挂了...
			showChangeSoldier = damage / target.lifePerSoldier

			if target.reliveIndex >= len(target.spellList.data.reliveSpell) {
				target.stopMove(frame)
				return
			}

			reliveSpell = target.spellList.data.reliveSpell[target.reliveIndex]
			target.reliveIndex++

			// 有复活技能，先别死
			target.setLife(1)
			aliveSoldier = target.getSoldier()
		}

		// 被打加怒气
		if target.spellList.HasRageSpell() && !target.isFullRage() && changeSoldier > 0 {
			lostLifePercent := lifePercent - target.getLifePercent()
			if lostLifePercent > 0 {
				target.addRage(combat.AddRageLost1Percent * lostLifePercent)
			}
		}

		// 触发复活
		if reliveSpell != nil {
			// 加血
			reliveLife := IMax(int(float64(target.totalLife)*reliveSpell.relivePercent), 1)
			target.setLife(reliveLife)

			// 触发被动
			triggerProto := &shared_proto.TroopTriggerPassiveSpellActionProto{}
			triggerProto.PassiveSpellId = reliveSpell.proto.Id
			triggerProto.SelfIndex = target.index
			triggerProto.TargetIndex = target.index
			addTriggerPassiveSpellAction(frameProto, triggerProto)

			// 加血
			effectProto := &shared_proto.TroopTickEffectProto{}
			effectProto.Id = reliveSpell.proto.Id
			effectProto.ChangeSoldier = -int32(target.getSoldier())
			effectProto.Soldier = int32(target.getSoldier())
			addSpellEffectAction(frameProto, target.index, effectProto)
		}
	}

	return
}

func (c *Troops) Update(combat *Combat, frame int, proto *shared_proto.CombatFrameProto) (isCheckEnd, isCheckEnemy bool) {

	// 处理延时伤害
	if c.nextDelayDamageFrame <= frame {

		c.nextDelayDamageFrame = math.MaxInt32
		if len(c.delayDamage) > 0 {
			for i, dd := range c.delayDamage {
				if dd == nil {
					continue
				}

				if dd.effectFrame <= frame {
					c.delayDamage[i] = nil

					c.MarkRage()

					// 触发
					changeShield, changeSoldier, showChangeSoldier, aliveSoldier, removeStateIds := combat.tryReduceDamage(c, dd.damage, frame, proto)
					dd.caster.addKillSoldier(changeSoldier)

					effectProto := &shared_proto.TroopTickEffectProto{}
					effectProto.Id = dd.spellId
					effectProto.ChangeShield = int32(changeShield)
					effectProto.RemoveState = removeStateIds
					effectProto.ChangeSoldier = int32(IMax(showChangeSoldier, 1))
					effectProto.Soldier = int32(aliveSoldier)
					effectProto.Rage = c.newRageUpdateProtoIfChanged()

					addSpellEffectAction(proto, c.index, effectProto)

					if !c.isAlive() {
						return true, false
					}
				} else {
					c.nextDelayDamageFrame = IMin(c.nextDelayDamageFrame, dd.effectFrame)
				}
			}
		}
	}

	// 前面处理掉各种状态
	moveSpeedChanged, rageRecoverChanged := c.updateState(combat, frame, proto)
	if !c.isAlive() {
		return true, false
	}

	// 恢复怒气
	c.RecoverRage()

	// 更新怒气
	if rageRecoverChanged {
		addStateEffectAction(proto, c.index, &shared_proto.TroopTickEffectProto{
			Rage: c.newRageUpdateProto(),
		})
	}

	if c.isStun() {
		// 晕眩中，退出
		return
	}

	if frame < c.strongeEndFrame {
		// 硬直中不处理打人和移动
		return
	}

	// 当前有路可以走，先走了这一下
	if c.hasStepToMove() {
		c.updatePosition(proto, frame)
	}

	// 看一下对方是否还有部队，如果没有，则找城墙
	enemy := combat.getEnemy(c.isAttacker)
	if enemy.isAllTroopDead() {
		// 尝试打城墙
		wall := enemy.wall
		if wall == nil || !wall.isAlive() {
			return true, true
		}

		// 打城墙只能用普攻
		baseSpell := c.spellList.baseSpell
		if IAbs(c.x-wall.x) <= baseSpell.data.ReleaseRange {
			// 攻击范围内
			c.stopMove(frame)

			// 检查一下上一个普攻目标是否挂了，如果挂了，则清空普攻CD
			if baseSpell.lastTarget != nil &&
				!baseSpell.lastTarget.isAlive() {
				baseSpell.lastTarget = nil
				baseSpell.nextReleaseFrame = 0
			}

			if frame < baseSpell.nextReleaseFrame {
				// 技能CD中
				return false, true
			}

			// 放技能
			checkEnd := combat.tryAttackWall(c, baseSpell.data, frame, proto)

			// 加技能cd
			baseSpell.nextReleaseFrame = frame + c.spellList.baseSpellCooldownFrame

			return checkEnd, true
		} else {
			// 攻击范围外，往城墙方向移动
			targetX := wall.x
			targetY := c.y
			if p := c.currentPath; p == nil || int(p.MoveEndX) != targetX || int(p.MoveEndY) != targetY {
				// 转向
				c.startMove(combat, frame, targetX, targetY, combat.CheckMoveFrame, proto)
			}
		}
		return
	}

	// 先把技能都遍历一遍，看下有没有可以释放的技能，如果没有可以释放的技能
	for _, spell := range c.spellList.spells {
		cd := spell.data.CooldownFrame

		// 普攻CD
		if c.isBaseSpell(spell.data) {
			if c.isNotAttackable() {
				continue
			}

			cd = c.spellList.baseSpellCooldownFrame

			// 检查一下上一个普攻目标是否挂了，如果挂了，则清空普攻CD
			if spell.lastTarget != nil &&
				!spell.lastTarget.isAlive() {
				spell.lastTarget = nil
				spell.nextReleaseFrame = 0
			}
		} else {
			if c.isSilence() {
				// 沉默
				continue
			}
		}

		if frame < spell.nextReleaseFrame {
			// 技能CD中
			continue
		}

		// 怒气技能
		if spell.data.proto.RageSpell && !c.isFullRage() {
			// 怒气未满
			continue
		}

		// 根据技能找施法对象
		targets := c.getSpellTarget(combat, spell.data, spell.lastTarget)
		if len(targets) <= 0 {
			continue
		}

		if !spell.data.proto.KeepMove {
			// 需要停下来放技能
			c.stopMove(frame)
		}

		shouldCheckEnd := combat.processReleaseSpell(c, targets, spell.data, frame,
			ReleaseTypeRelease, true, proto)
		isCheckEnd = isCheckEnd || shouldCheckEnd

		// 加硬直
		c.strongeEndFrame = IMax(c.strongeEndFrame, frame+spell.data.StrongeFrame)

		// 加技能cd
		spell.nextReleaseFrame = frame + cd

		if c.isBaseSpell(spell.data) {
			c.checkShortMove = true
		}

		spell.lastTarget = targets[0]

		c.clearRushTarget()
	}

	// 放完技能，如果有硬直，就不再走了
	if c.strongeEndFrame > frame {
		return
	}

	// 看下普攻范围内有没有可以打的人，有的话，不走，没有的话就走
	hasTargetInBaseSpellRange, hasAnyCaptainAlive := c.hasEnemyInBaseSpellRange(combat)
	if !hasAnyCaptainAlive {
		// 对面全死了，退出
		return !enemy.HasAliveWall(), true
	}

	if hasTargetInBaseSpellRange {
		// 如果普攻的攻击范围内有人，则停止移动，等攻击CD
		c.stopMove(frame)
		return false, true
	}

	if moveSpeedChanged && c.hasStepToMove() {
		// 移动速度改变了，停下来，重新找一条路来走
		c.stopMove(frame)
	}

	if !c.isUnmovable() {
		// 当前位置没有可以打的人，根据逻辑寻找移动位置
		if !c.hasStepToMove() || c.nextCheckMoveTargetFrame <= frame {
			// 没有移动目标，找一个目标，开始走
			targetX, targetY := c.selectMoveTarget(combat)

			if p := c.currentPath; p == nil || int(p.MoveEndX) != targetX || int(p.MoveEndY) != targetY {
				// 转向
				c.startMove(combat, frame, targetX, targetY, combat.CheckMoveFrame, proto)
			} else {
				// 目标位置不变，那么不做额外的处理，只修改下次检查目标位置的时间
				c.updateNextCheckMoveTargetFrame(frame, combat.CheckMoveFrame)
			}
		}
	}

	return false, true
}

func (c *Troops) hasEnemyInBaseSpellRange(combat *Combat) (hasTargetInBaseSpellRange, hasAnyCaptainAlive bool) {
	for _, captain := range combat.getEnemyTroops(c.isAttacker) {
		if !captain.isAlive() {
			continue
		}
		hasAnyCaptainAlive = true
		if c.isInRange(captain.x, captain.y, c.spellList.data.baseSpell.ReleaseRange) {
			hasTargetInBaseSpellRange = true
			break
		}
	}
	return
}

const (
	ReleaseTypeRelease = 0
	ReleaseTypeTrigger = 1
)

// 行动类型，1-移动 2-放技能 3-技能生效 4-buff生效 5-触发被动 6-短距离移动（不改阵型）
const (
	ActionTypeMove                = 1
	ActionTypeReleaseSpell        = 2
	ActionTypeSpellEffect         = 3
	ActionTypeStateEffect         = 4
	ActionTypeTriggerPassiveSpell = 5
	ActionTypeShortMove           = 6
)

func addMoveAction(proto *shared_proto.CombatFrameProto, index int32, moveProto *shared_proto.TroopMoveActionProto, isShortMove bool) {
	action := &shared_proto.TroopActionProto{
		Index:      index,
		ActionType: ActionTypeMove,
		Move:       moveProto,
	}

	if isShortMove {
		action.ActionType = ActionTypeShortMove
	}

	proto.Action = append(proto.Action, action)
}

func addReleaseSpellAction(proto *shared_proto.CombatFrameProto, index int32, releaseSpellProto *shared_proto.TroopReleaseSpellActionProto) {
	proto.Action = append(proto.Action, &shared_proto.TroopActionProto{
		Index:        index,
		ActionType:   ActionTypeReleaseSpell,
		ReleaseSpell: releaseSpellProto,
	})
}

func addSpellEffectAction(proto *shared_proto.CombatFrameProto, index int32, tickProto *shared_proto.TroopTickEffectProto) {
	proto.Action = append(proto.Action, &shared_proto.TroopActionProto{
		Index:       index,
		ActionType:  ActionTypeSpellEffect,
		SpellEffect: tickProto,
	})
}

func addStateEffectAction(proto *shared_proto.CombatFrameProto, index int32, tickProto *shared_proto.TroopTickEffectProto) {
	proto.Action = append(proto.Action, &shared_proto.TroopActionProto{
		Index:       index,
		ActionType:  ActionTypeStateEffect,
		StateEffect: tickProto,
	})
}

func addTriggerPassiveSpellAction(proto *shared_proto.CombatFrameProto, triggerProto *shared_proto.TroopTriggerPassiveSpellActionProto) {
	proto.Action = append(proto.Action, &shared_proto.TroopActionProto{
		Index:               triggerProto.SelfIndex,
		ActionType:          ActionTypeTriggerPassiveSpell,
		TriggerPassiveSpell: triggerProto,
	})
}

func (c *Combat) processReleaseSpell(attacker *Troops, targets []*Troops, spell *SpellData, frame, releaseType int, canTriggerSpell bool,
	proto *shared_proto.CombatFrameProto) (shouldCheckEnd bool) {

	// 技能替换逻辑，如果目标身上具有某种状态，则替换施法技能 TODO

	releaseSpellProto := &shared_proto.TroopReleaseSpellActionProto{
		ReleaseType: int32(releaseType),
		SpellId:     spell.Id(),
		EndFrame:    0, // 持续生效技能有值
		SpellX:      int32(attacker.x),
		SpellY:      int32(attacker.y),
	}
	addReleaseSpellAction(proto, attacker.index, releaseSpellProto)

	//if len(targets) == 1 {
	//	if target := targets[0]; target != nil {
	//		releaseSpellProto.SpellTarget = target.index
	//		releaseSpellProto.SpellX = int32(target.x)
	//		releaseSpellProto.SpellY = int32(target.y)
	//	}
	//}

	attacker.MarkRage()
	if spell.proto.RageSpell {
		// 扣怒气
		attacker.ClearRage()
		releaseSpellProto.ClearRage = true
	}

	//var self
	var selfProto *shared_proto.TroopReleaseSpellEffectProto

	// 给自己加状态
	for _, stateWithRate := range spell.SelfState {
		stateData := stateWithRate.Data
		if !c.tryRate(int(stateWithRate.TriggerRate)) {
			continue
		}

		newState := attacker.addState(stateData, frame, c.FramePerSecond, attacker, 0)
		if newState != nil {
			selfProto = attacker.newReleaseSpellEffectProtoIfNil(selfProto)
			selfProto.AddState = append(selfProto.AddState, newState.data.Id())
			selfProto.AddStateEndFrame = append(selfProto.AddStateEndFrame, int32(newState.endFrame))
		}
	}

	// 给自己加怒气
	if spell.SelfRage > 0 {
		attacker.addRage(spell.SelfRage)
	}

	for _, target := range targets {
		if c.isDebug {
			logrus.Debugf("Frame: %v 释放技能 %v %v %v", frame, attacker.index, target.index, spell.Id())
		}

		c.processReleaseSpellToTarget(attacker, target, spell, frame, canTriggerSpell,
			proto, releaseSpellProto)

		shouldCheckEnd = shouldCheckEnd || !target.isAlive()
	}

	if rageProto := attacker.newRageUpdateProtoIfChanged(); rageProto != nil {
		selfProto = attacker.newReleaseSpellEffectProtoIfNil(selfProto)
		selfProto.Rage = rageProto
	}

	if selfProto != nil {
		selfProto.HurtType = shared_proto.HurtType_NORMAL
		releaseSpellProto.Target = append(releaseSpellProto.Target, selfProto)
	}

	// 处理触发技能
	if canTriggerSpell && len(c.triggerSpell) > 0 {
		for _, ts := range c.triggerSpell {
			if ts.attacker.isAlive() {

				spell := ts.spell

				var targets []*Troops
				if !spell.proto.FriendSpell && spell.ReleaseRange == 0 && spell.HurtRange == 0 {
					// 对目标触发
					if !ts.target.isAlive() {
						// 目标挂了
						continue
					}
					targets = []*Troops{ts.target}
				} else {
					// 重新查找目标
					targets = ts.attacker.getSpellTarget(c, spell, nil)
					if len(targets) <= 0 {
						continue
					}
				}

				needCheckEnd := c.processReleaseSpell(ts.attacker, targets, ts.spell,
					frame, ReleaseTypeTrigger, false, proto)
				shouldCheckEnd = shouldCheckEnd || needCheckEnd
			}
		}

		c.triggerSpell = nil
	}

	return
}

func (c *Combat) processReleaseSpellToTarget(attacker *Troops, target *Troops, spell *SpellData, frame int, canTriggerSpell bool,
	frameProto *shared_proto.CombatFrameProto, proto *shared_proto.TroopReleaseSpellActionProto) {

	targetProto := target.newReleaseSpellEffectProtoIfNil(nil)
	targetProto.IsTarget = true
	proto.Target = append(proto.Target, targetProto)

	// 是否是普攻技能
	isBaseSpell := attacker.isBaseSpell(spell)

	// 给目标加减怒气
	if attacker != target {
		target.MarkRage()
	}
	if spell.TargetRage != 0 {
		target.addRage(spell.TargetRage)
	}

	targetProto.HurtType = shared_proto.HurtType_NORMAL
	if spell.Coef > 0 {

		var damage int
		if isBaseSpell {
			// 先看看是否闪避攻击
			if c.tryDodge(attacker, target) {
				targetProto.HurtType = shared_proto.HurtType_MISS
			} else {
				// 普攻
				damage, targetProto.HurtType = c.calculateDamage(attacker, target, spell)

				if attacker.spellList.HasRageSpell() && !attacker.isFullRage() {
					// 命中加怒气
					attacker.addRage(c.AddRagePerHint)
				}
			}
		} else {
			// 技能
			damage = c.calculateSpellDamage(attacker, target, spell.Coef, spell.proto.EffectType, spell.proto.HurtType)
		}

		if damage > 0 {
			delayFrame := 0
			if spell.FlySpeedPerFrame > 0 {
				distance := Distance(attacker.x, attacker.y, target.x, target.y)
				if distance > 0 {
					delayFrame = distance / spell.FlySpeedPerFrame
				}
			}

			if delayFrame > 0 {
				effectFrame := frame + delayFrame
				targetProto.EffectFrame = int32(effectFrame)
				// 延时生效技能
				target.addDelayDamage(attacker, damage, effectFrame, spell.Id())
			} else {

				changeShield, changeSoldier, showChangeSoldier, aliveSoldier, removeStateIds := c.tryReduceDamage(target, damage, frame, frameProto)
				attacker.addKillSoldier(changeSoldier)

				targetProto.ChangeShield = int32(changeShield)
				targetProto.RemoveState = removeStateIds
				targetProto.ChangeSoldier = int32(showChangeSoldier)
				targetProto.Soldier = int32(aliveSoldier)

				if !target.isAlive() {
					// 挂了...
					return
				}
			}
		}
	}

	// 给目标加状态
	for _, stateWithRate := range spell.TargetState {
		stateData := stateWithRate.Data
		if !c.tryRate(int(stateWithRate.TriggerRate)) {
			continue
		}

		var tickDamage int
		if stateData.DamageCoef > 0 {
			tickDamage = c.calculateSpellDamage(attacker, target, stateData.DamageCoef, stateData.proto.EffectType, stateData.proto.DamageHurtType)
		}

		newState := target.addState(stateData, frame, c.FramePerSecond, attacker, tickDamage)
		if newState != nil {
			targetProto.AddState = append(targetProto.AddState, newState.data.Id())
			targetProto.AddStateEndFrame = append(targetProto.AddStateEndFrame, int32(newState.endFrame))
		}
	}

	if attacker != target {
		targetProto.Rage = target.newRageUpdateProtoIfChanged()
	}

	// 打断移动
	if target.isStun() || target.isUnmovable() {
		target.stopMove(frame)
	}

	// 尝试触发其他技能
	if canTriggerSpell && attacker.isAttacker != target.isAttacker {

		if isBaseSpell {
			// 先战
			if !attacker.spellList.baseSpellUsed {
				attacker.spellList.baseSpellUsed = true

				for _, ps := range attacker.spellList.data.firstAttackSpell {
					if c.tryTriggerPassiveSpell(attacker, target, ps, frame) {
						triggerProto := c.doTriggerPassiveSpell(attacker, target, ps, frame, frameProto)
						addTriggerPassiveSpellAction(frameProto, triggerProto)
					}
				}
			}

			// 首次攻击目标
			if len(attacker.spellList.data.firstAttackTargetSpell) > 0 {
				if _, exist := attacker.spellList.spellTargetMap[target.index]; !exist {
					attacker.spellList.spellTargetMap[target.index] = struct{}{}

					for _, ps := range attacker.spellList.data.firstAttackTargetSpell {
						if c.tryTriggerPassiveSpell(attacker, target, ps, frame) {
							triggerProto := c.doTriggerPassiveSpell(attacker, target, ps, frame, frameProto)
							addTriggerPassiveSpellAction(frameProto, triggerProto)
						}
					}
				}
			}

			// 每次普攻
			for _, ps := range attacker.spellList.data.times1Spell {
				if c.tryTriggerPassiveSpell(attacker, target, ps, frame) {
					triggerProto := c.doTriggerPassiveSpell(attacker, target, ps, frame, frameProto)
					addTriggerPassiveSpellAction(frameProto, triggerProto)
				}
			}

			if attacker.spellList.baseSpellHurtTarget == target.index {
				attacker.spellList.baseSpellHurtTimes++

				// N次普攻
				for _, ps := range attacker.spellList.data.timesNSpell {
					if ps.triggerHit <= 0 {
						// 防御性
						continue
					}

					if attacker.spellList.baseSpellHurtTimes%ps.triggerHit != 0 {
						// 没触发
						continue
					}

					if c.tryTriggerPassiveSpell(attacker, target, ps, frame) {
						triggerProto := c.doTriggerPassiveSpell(attacker, target, ps, frame, frameProto)
						addTriggerPassiveSpellAction(frameProto, triggerProto)
					}
				}
			} else {
				attacker.spellList.baseSpellHurtTarget = target.index
				attacker.spellList.baseSpellHurtTimes = 1
			}

			// 被打触发
			if target.isAlive() {
				for _, ps := range target.spellList.data.beenHurtSpell {
					if c.tryTriggerPassiveSpell(target, attacker, ps, frame) {
						triggerProto := c.doTriggerPassiveSpell(target, attacker, ps, frame, frameProto)
						addTriggerPassiveSpellAction(frameProto, triggerProto)
					}
				}
			}
		}
	}
}

func (c *Combat) doTriggerPassiveSpell(trigger, target *Troops, spell *PassiveSpellData, frame int, frameProto *shared_proto.CombatFrameProto) *shared_proto.TroopTriggerPassiveSpellActionProto {

	// 加触发CD
	if spell.targetCooldownFrame > 0 {
		target.spellList.SetTriggerCd(spell.proto.Id, frame+spell.targetCooldownFrame)
	}

	proto := &shared_proto.TroopTriggerPassiveSpellActionProto{}
	proto.PassiveSpellId = spell.proto.Id
	proto.SelfIndex = trigger.index
	if target != nil {
		proto.TargetIndex = target.index
	}

	if len(spell.selfState) > 0 {
		// 给自己加状态
		for _, stateData := range spell.selfState {
			newState := trigger.addState(stateData, frame, c.FramePerSecond, trigger, 0)
			if newState != nil {
				proto.SelfAddState = append(proto.SelfAddState, newState.data.Id())
				proto.SelfAddStateEndFrame = append(proto.SelfAddStateEndFrame, int32(newState.endFrame))
			}
		}
	}

	if target != nil && target.isAlive() && trigger.isAttacker != target.isAttacker {
		if len(spell.targetState) > 0 {
			// 给目标加状态
			for _, stateData := range spell.targetState {
				var tickDamage int
				if stateData.DamageCoef > 0 {
					tickDamage = c.calculateSpellDamage(trigger, target, stateData.DamageCoef, stateData.proto.EffectType, stateData.proto.DamageHurtType)
				}

				newState := target.addState(stateData, frame, c.FramePerSecond, trigger, tickDamage)
				if newState != nil {
					proto.TargetAddState = append(proto.TargetAddState, newState.data.Id())
					proto.TargetAddStateEndFrame = append(proto.TargetAddStateEndFrame, int32(newState.endFrame))
				}
			}
		}

		// 持续效果立即生效
		if spell.proto.ExciteEffectType != 0 {
			if target.hasEffectState(spell.proto.ExciteEffectType) {
				// 遍历状态，将持续掉血的状态弄出来
				for _, s := range target.stateMap {
					if s.tickDamage <= 0 {
						continue
					}

					if s.data.proto.EffectType != spell.proto.ExciteEffectType {
						continue
					}

					// 移除状态
					target.removeState(s.data)

					// 立即生效
					damage := s.getRemainTickTimes() * s.tickDamage

					target.MarkRage()
					changeShield, changeSoldier, showChangeSoldier, aliveSoldier, removeStateIds := c.tryReduceDamage(target, damage, frame, frameProto)
					s.caster.addKillSoldier(changeSoldier)

					effectProto := &shared_proto.TroopTickEffectProto{}
					effectProto.Id = s.data.Id()
					effectProto.ChangeShield = int32(changeShield)
					effectProto.RemoveState = removeStateIds
					effectProto.ChangeSoldier = int32(IMax(showChangeSoldier, 1))
					effectProto.Soldier = int32(aliveSoldier)
					effectProto.Rage = target.newRageUpdateProtoIfChanged()

					effectProto.RemoveState = append(effectProto.RemoveState, s.data.Id())

					proto.TargetStateEffect = append(proto.TargetStateEffect, effectProto)

					if !target.isAlive() {
						// 挂了...
						break
					}
				}
			}
		}
	}

	// 释放技能
	if spell.spell != nil {
		c.addTriggerSpell(trigger, target, spell.spell)
	}

	return proto
}

func (c *Combat) tryTriggerPassiveSpell(attacker, target *Troops, spell *PassiveSpellData, frame int) bool {

	switch spell.proto.TriggerType {
	case shared_proto.SpellTriggerType_STNone:
		return false
	case shared_proto.SpellTriggerType_STBeginRelease:
		return false
	case shared_proto.SpellTriggerType_STFirstHit:
		if !spell.triggerTarget.IsValidTroop(target) {
			return false
		}

		// 先战技能，在外面判断过了，其他情况都触发
		return true
	case shared_proto.SpellTriggerType_STHitN:

		// 打X次触发这个要看职业
		if !spell.triggerTarget.IsValidTroop(target) {
			return false
		}
	case shared_proto.SpellTriggerType_STBeenHit:

	case shared_proto.SpellTriggerType_STShieldBroken:
		// 破盾一般不看其他，都不会进来这里问
		return true
	}

	if spell.targetCooldownFrame > 0 && target.spellList.IsInTriggerCd(spell.proto.Id, frame) {
		return false
	}

	return c.tryRate(spell.triggerRate)
}

func (c *Combat) tryReduceShield(t *Troops, damage int) (totalReduceAmount int, removeStateIds []int32) {

	// 护盾抵扣伤害
	if !t.hasShield() {
		return
	}

	for _, s := range t.stateMap {
		if s.shield <= 0 {
			continue
		}

		maxReduceAmount := IMax(int(float64(damage)*s.data.ShieldEffectRate), 1)
		if maxReduceAmount <= totalReduceAmount {
			// 已经达到这个护盾可以抵消的数值了
			continue
		}

		toReduce := IMin(s.shield, maxReduceAmount)

		s.shield -= toReduce
		totalReduceAmount += toReduce

		if s.shield <= 0 {
			// 破盾
			t.removeState(s.data)
			removeStateIds = append(removeStateIds, s.data.Id())
		}
	}

	return
}

const RateDenominator = 10000

func (c *Combat) tryRate(rate int) bool {
	if rate < c.IConfigDenominator {
		return c.randomRate(c.IConfigDenominator) < rate
	}

	return true
}

func (c *Combat) randomRate(n int) int {
	return c.random.Intn(n)
}

func (c *Combat) addTriggerSpell(attacker, target *Troops, spell *SpellData) {
	c.triggerSpell = append(c.triggerSpell, &TriggerSpell{
		attacker: attacker,
		target:   target,
		spell:    spell,
	})
}

type TriggerSpell struct {
	attacker *Troops
	target   *Troops

	spell *SpellData
}

const closeDistance = 100

func (c *Troops) selectMoveTarget(combat *Combat) (targetX, targetY int) {

	// 如果有对线目标，则走冲刺流程
	if ok, targetX, targetY := c.selectRushMoveTarget(combat); ok {
		return targetX, targetY
	}

	// 没有对线目标，普通流程
	// 视野内有敌军，向着敌军移动

	var target *Troops
	targetPriority := 0
	targetDistance := 0
	for _, captain := range combat.getEnemyTroops(c.isAttacker) {
		if !captain.isAlive() {
			continue
		}

		distance := Distance(c.x, c.y, captain.x, captain.y)
		if c.getRaceData().viewRange < distance {
			// 看不见，跳过
			continue
		}

		priority := c.getRaceData().getTargetPriority(captain.getRace())
		if target != nil {
			if priority < targetPriority {
				// 优先职业（如果都是在视野内，如果优先级没超过target，直接跳过）
				continue
			}

			if priority == targetPriority && targetDistance <= distance {
				// 优先级一致，距离远的跳过
				continue
			}
		}

		target = captain
		targetPriority = priority
		targetDistance = distance
	}

	if target != nil {
		return target.x, target.y
	}

	// 视野内没有敌军
	if !c.isFar() {
		// 近战，直接根据优先级，找第一优先作为目标，移动
		target := combat.getEnemy(c.isAttacker)
		for _, targetRace := range c.getRaceData().proto.Priority {
			raceTroops := target.troopRaceMap[targetRace]
			if len(raceTroops) <= 0 {
				continue
			}

			var nearestCaptain *Troops
			var nearestDistance int
			for _, captain := range raceTroops {
				if !captain.isAlive() {
					continue
				}

				d := Distance(c.x, c.y, captain.x, captain.y)
				if nearestCaptain == nil || d < nearestDistance {
					nearestCaptain = captain
					nearestDistance = d
				}
			}

			if nearestCaptain != nil {
				return nearestCaptain.x, nearestCaptain.y
			}
		}

		// 近战实在没找到，也走一遍远程的逻辑
		logrus.Error("近战寻找移动目标，居然没找到人？")
	}

	// 远程，存在前方的敌人，则向前移动（附加一些提前量，防止冲过头），如果最后一个家伙，则Y方向上接近
	var lastEnemy *Troops
	for _, captain := range combat.getEnemyTroops(c.isAttacker) {
		if !captain.isAlive() {
			continue
		}

		if lastEnemy != nil {
			if c.isAttacker {
				// 从左到右
				if lastEnemy.x >= captain.x {
					continue
				}
			} else {
				// 从右到左
				if lastEnemy.x <= captain.x {
					continue
				}
			}
		}

		lastEnemy = captain
	}

	if lastEnemy != nil {
		if c.isAttacker {
			// 攻击方，从左到右（x从小到大）
			if c.x+closeDistance >= lastEnemy.x {
				return lastEnemy.x, lastEnemy.y
			}
		} else {
			if lastEnemy.x+closeDistance >= c.x {
				return lastEnemy.x, lastEnemy.y
			}
		}
	}

	// 如果这个也没有，那么朝着X方向向前冲
	if c.isAttacker {
		return combat.mapXLen, c.y
	} else {
		return 0, c.y
	}
}

func (c *Troops) selectRushMoveTarget(combat *Combat) (ok bool, targetX, targetY int) {
	if c.rushTarget == nil {
		return
	}

	if !c.rushTarget.isAlive() {
		c.clearRushTarget()
		return
	}

	// 如果对位目标足够近，不看别的，直接上去干
	rushTargetXDiff := IAbs(c.x - c.rushTarget.x)
	if rushTargetXDiff <= closeDistance {
		return true, c.rushTarget.x, c.rushTarget.y
	}

	// 我是远程的，那么不管，往前冲就好了
	if c.getRaceData().proto.IsFar {
		if c.isAttacker {
			return true, combat.mapXLen, c.y
		} else {
			return true, 0, c.y
		}
	}

	frame3Distance := c.rushTarget.movePerFrame * 3

	// 如果相邻的部队在Y方向上有交叉，则转向干架
	var target *Troops
	targetPriority := 0
	for _, t := range c.rushUpDownTarget {
		if !t.isAlive() {
			continue
		}

		diff := IAbs(c.x - t.x)
		if diff > closeDistance {
			continue
		}

		if rushTargetXDiff-diff <= frame3Distance {
			// 距离太近，不考虑
			continue
		}

		priority := c.getRaceData().getTargetPriority(t.getRace())
		if target != nil {
			if priority <= targetPriority {
				// 优先职业（如果都是在视野内，如果优先级没超过target，直接跳过）
				continue
			}
		}

		target = t
		targetPriority = priority
	}

	if target != nil {
		return true, target.x, target.y
	}

	// 如果这个也没有，那么朝着X方向向前冲
	if c.isAttacker {
		return true, combat.mapXLen, c.y
	} else {
		return true, 0, c.y
	}
}

func newMovePath(x1, y1, x2, y2, endFrame int) *shared_proto.TroopMoveActionProto {

	return &shared_proto.TroopMoveActionProto{
		MoveStartX:   int32(x1),
		MoveStartY:   int32(y1),
		MoveEndX:     int32(x2),
		MoveEndY:     int32(y2),
		MoveEndFrame: int32(endFrame),
	}

}

func (c *Troops) tryTruncCurrentPath(frame int) {

	if c.currentPath == nil {
		return
	}

	c.currentPath.MoveEndX = int32(c.x)
	c.currentPath.MoveEndY = int32(c.y)
	c.currentPath.MoveEndFrame = int32(frame)
	c.currentPath = nil
}

func (c *Troops) startMove(combat *Combat, frame, targetX, targetY, nextCheckDiffFrame int, proto *shared_proto.CombatFrameProto) {

	c.tryTruncCurrentPath(frame)

	distance := Distance(c.x, c.y, targetX, targetY)
	if distance <= 0 {
		logrus.Errorf("Troops.startMove() distance <= 0")
		return
	}

	isShortMove := false
	if c.checkShortMove {
		if distance <= combat.ShortMoveDistance {
			isShortMove = true
		} else {
			c.checkShortMove = false
		}
	}

	moveFrame := (distance + c.movePerFrame - 1) / c.movePerFrame

	c.currentPath = newMovePath(c.x, c.y, targetX, targetY, frame+moveFrame)
	c.currentPathStartFrame = int32(frame)

	addMoveAction(proto, c.index, c.currentPath, isShortMove)

	c.updateNextCheckMoveTargetFrame(frame, nextCheckDiffFrame)
}

func (c *Troops) updateNextCheckMoveTargetFrame(frame, toAdd int) {
	c.nextCheckMoveTargetFrame = frame + toAdd
}

func (c *Troops) updatePosition(proto *shared_proto.CombatFrameProto, frame int) {

	currentPath := c.currentPath
	if currentPath == nil {
		// 防御性
		return
	}

	if frame >= int(currentPath.MoveEndFrame) {
		c.x = int(currentPath.MoveEndX)
		c.y = int(currentPath.MoveEndY)

		// 理论上很少发生，会因为打人，或者目标位置改变而更改移动位置，不会一条道走到黑
		c.currentPath = nil
		return
	}

	totalMoveFrame := currentPath.MoveEndFrame - c.currentPathStartFrame
	curMoveFrame := int32(frame) - c.currentPathStartFrame

	newX := currentPath.MoveStartX + (currentPath.MoveEndX-currentPath.MoveStartX)*curMoveFrame/totalMoveFrame
	newY := currentPath.MoveStartY + (currentPath.MoveEndY-currentPath.MoveStartY)*curMoveFrame/totalMoveFrame

	c.x = int(newX)
	c.y = int(newY)
}

func (c *Troops) stopMove(frame int) {
	c.tryTruncCurrentPath(frame)
}

func (c *Troops) getSpellTarget(combat *Combat, spell *SpellData, lastSpellTarget *Troops) []*Troops {

	isAtk := c.isAttacker
	if spell.proto.FriendSpell {
		isAtk = !isAtk
	}
	targetPlayer := combat.getEnemy(isAtk)

	if spell.ReleaseRange > 0 {
		if spell.HurtRange > 0 {
			// 根据施法距离，找多个目标
			targets := c.findReleaseSpellTarget(combat, targetPlayer, spell, spell.ReleaseRange, 1, lastSpellTarget)

			if len(targets) > 0 {
				target := targets[0]

				for _, t := range targetPlayer.troops {
					if t == target {
						continue
					}

					if !target.isInRange(t.x, t.y, spell.HurtRange) {
						continue
					}

					targets = append(targets, t)
					if len(targets) >= spell.HurtCount {
						return targets
					}
				}
			}

			return targets
		}

		return c.findReleaseSpellTarget(combat, targetPlayer, spell, spell.ReleaseRange, spell.HurtCount, lastSpellTarget)
	}

	// 以自己为目标
	if spell.HurtRange > 0 {
		return c.findReleaseSpellTarget(combat, targetPlayer, spell, spell.HurtRange, spell.HurtCount, lastSpellTarget)
	}

	return []*Troops{c}
}

func (c *Troops) findReleaseSpellTarget(combat *Combat, target *CombatPlayer, spell *SpellData, radius, count int, lastSpellTarget *Troops) []*Troops {

	if spell.proto.TargetSubType == shared_proto.SpellTargetSubType_SSTNone {

		var validTargets []*Troops
		if c.isAttacker == target.isAttacker {
			// 友方技能
			if lastSpellTarget != nil && c.isValidSpellTarget(lastSpellTarget, spell, radius) {
				validTargets = append(validTargets, lastSpellTarget)
			}

			if len(validTargets) < count {
				// 找到第一个符合条件的目标返回
				for _, t := range target.troops {
					if t == lastSpellTarget {
						continue
					}

					if c.isValidSpellTarget(t, spell, radius) {
						validTargets = append(validTargets, t)

						if len(validTargets) >= count {
							break
						}
					}
				}
			}
		} else {
			if lastSpellTarget != nil && !c.isValidSpellTarget(lastSpellTarget, spell, radius) {
				lastSpellTarget = nil
			}

			// 敌方技能
		out:
			for _, targetRace := range c.getRaceData().proto.Priority {
				raceTroops := target.troopRaceMap[targetRace]
				if len(raceTroops) <= 0 {
					continue
				}

				// 最后的目标，优先选择
				if lastSpellTarget != nil && lastSpellTarget.getRace() == targetRace {
					validTargets = append(validTargets, lastSpellTarget)
					if len(validTargets) >= count {
						break out
					}
				}

				for _, t := range raceTroops {
					if t == lastSpellTarget {
						continue
					}

					if c.isValidSpellTarget(t, spell, radius) {
						validTargets = append(validTargets, t)

						if len(validTargets) >= count {
							break out
						}
					}
				}
			}
		}

		return validTargets
	}

	// 找到所有符合条件的目标
	var validTargets []*Troops
	for _, t := range target.troops {
		if c.isValidSpellTarget(t, spell, radius) {
			validTargets = append(validTargets, t)
		}
	}

	if len(validTargets) <= count {
		return validTargets
	}

	switch spell.proto.TargetSubType {
	case shared_proto.SpellTargetSubType_SSTRandom:
		// 随机数组
		mixTroops(validTargets, count, combat.random)
	case shared_proto.SpellTargetSubType_STTMaxLifePercent:
		for _, t := range validTargets {
			t.MarkLifePercent()
		}

		// 根据血量百分比排序
		sortTroopsByMaxLifePercent(validTargets)

	case shared_proto.SpellTargetSubType_STTMinLifePercent:
		for _, t := range validTargets {
			t.MarkLifePercent()
		}

		// 根据血量百分比最低
		sortTroopsByMinLifePercent(validTargets)

	case shared_proto.SpellTargetSubType_SSTMaxRage:
		// 根据怒气排序
		sortTroopsByMaxRage(validTargets)

	case shared_proto.SpellTargetSubType_SSTMinRage:
		sortTroopsByMinRage(validTargets)

	}

	return validTargets[:count]
}

func (c *Troops) isValidSpellTarget(target *Troops, spell *SpellData, r int) bool {
	if !target.isAlive() {
		return false
	}

	if !spell.proto.SelfAsTarget && target == c {
		return false
	}

	switch spell.proto.TargetSubType {
	case shared_proto.SpellTargetSubType_SSTMinRage,
		shared_proto.SpellTargetSubType_SSTMaxRage:
		if !target.spellList.HasRageSpell() {
			return false
		}
	}

	if !c.isInRange(target.x, target.y, r) {
		return false
	}

	if !spell.Target.IsValidTroop(target) {
		return false
	}

	return true
}

type FrameSlice []*shared_proto.CombatFrameProto

func (p FrameSlice) Len() int           { return len(p) }
func (p FrameSlice) Less(i, j int) bool { return p[i].Frame < p[j].Frame }
func (p FrameSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

const (
	HurtTypeKillSoldier = 1
	HurtTypeAvgStat     = 2
)

// 计算技能伤害
func (c *Combat) calculateSpellDamage(troopA, troopB *Troops, spellFixCoef float64, effectType, hurtType int32) int {

	A := troopA.totalStat
	B := troopB.totalStat
	switch hurtType {
	case HurtTypeKillSoldier:
		// 真实伤害
		killSoldier := float64(troopA.getSoldier()) * spellFixCoef
		return int(math.Max(1, killSoldier)) * troopB.lifePerSoldier
	case HurtTypeAvgStat:
		A = troopA.avgStat
		B = troopB.avgStat
	default:
		// 伤害公式
	}

	damage := c.calculateNormalDamage(troopA, troopB, A, B)

	spellCoef := spellFixCoef
	if effectType > 0 {
		if coef := troopB.beenHurtEffectIncCoef[effectType]; coef > 0 {
			spellCoef = spellCoef * (1 + coef)
		}

		if coef := troopB.beenHurtEffectDecCoef[effectType]; coef > 0 {
			spellCoef = spellCoef / (1 + coef)
		}
	}

	//若A敏捷<B敏捷，
	//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数×（（A敏*（A敏+ B敏）/（2*B敏*B敏））^0.5)；
	//若A敏捷=B敏捷，
	//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数)；
	//若A敏捷>B敏捷，
	//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数×（2*A敏*A敏/（B敏*B敏+A敏*B敏））^0.5)
	//以上，技能伤害系数=0.2
	//克制技伤害不能被闪避，也不会暴击。

	result := 0
	switch {
	case A.dexterity < B.dexterity:
		//若A敏捷<B敏捷，
		//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数×（（A敏*（A敏+ B敏）/（2*B敏*B敏））^0.5)；
		coef := math.Sqrt(A.dexterity * (A.dexterity + B.dexterity) / (2 * B.dexterity * B.dexterity))
		result = int(math.Max(1, damage*spellCoef*float64(troopA.soldier)*coef))
	case A.dexterity > B.dexterity:
		//若A敏捷>B敏捷，
		//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数×（2*A敏*A敏/（B敏*B敏+A敏*B敏））^0.5)
		coef := math.Sqrt(2 * A.dexterity * A.dexterity / (B.dexterity*B.dexterity + A.dexterity*B.dexterity))
		result = int(math.Max(1, damage*spellCoef*float64(troopA.soldier)*coef))
	default:
		//若A敏捷=B敏捷，
		//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数)；
		result = int(math.Max(1, damage*spellCoef*float64(troopA.soldier)))
	}

	if c.isDebug {
		logrus.Debugf("spell_result: %v %v %v %v %v", result, troopA.soldier, troopB.soldier, c.Coef, spellCoef)
	}

	return result
}

func (c *Combat) calculateDamage(A, B *Troops, spell *SpellData) (int, shared_proto.HurtType) {

	statA := A.totalStat
	statB := B.totalStat
	if spell.proto.HurtType == HurtTypeAvgStat {
		statA = A.avgStat
		statB = B.avgStat
	}

	damage := c.calculateNormalDamage(A, B, statA, statB)
	hurtType := shared_proto.HurtType_NORMAL

	// 计算暴击
	isCrit, critCoef := c.tryCrit(A, B, statA, statB)
	if isCrit {
		damage = damage * math.Max(1, critCoef)
		hurtType = shared_proto.HurtType_CRIT
	}

	result := int(math.Max(1, damage*float64(A.soldier)*spell.Coef))

	if c.isDebug {
		logrus.Debugf("result: %v %v %v %v %v %v %v", result, hurtType, A.soldier, B.soldier, c.Coef, spell.Coef, critCoef)
	}
	// 最终伤害=MAX{1，INT(单兵伤害×攻击部队当前士兵数) }
	return result, hurtType
}

func (c *Combat) calculateNormalDamage(troopA, troopB *Troops, A, B *F64Stat) float64 {

	coef := c.Coef

	// A打B
	// 防御系数D=（（A体力+A）* （B体力+A））^0.5 /（B防御+A）
	D := math.Sqrt((A.strength+coef)*(B.strength+coef)) / math.Max(B.defense+coef, 1)

	// 体力系数H，如果A体力≤B体力，则H=1.00，否则H=arctan（A体力 / B体力）*1.083+0.15
	H := 1.0
	if A.strength > B.strength {
		H = math.Atan(A.strength/math.Max(B.strength, 1))*1.083 + 0.15
	}

	// T为兵种克制系数，读表获得
	//T := math.Max(troopA.getRaceData().getTroopsCoef(troopB.getRace()), 0.1)
	T := 1.0

	//// 士气系数S（武将属性）（已过期，不再使用）
	//S := math.Max(A.Scoef, 1)

	// 加减伤系数S=（1+A的加伤1+A的加伤2+B的伤害加深1+B的伤害加深2…）/（1+B的减伤1+B的减伤2+A的伤害减弱1+A的伤害减弱2…）
	S := (1 + A.damageIncrePer + B.beenHurtIncrePer) / (1 + B.damageDecrePer + A.beenHurtDecrePer)

	r := 0.9 + c.random.Float64()*0.2

	// 普通伤害= A攻击*（A攻击+A）/（A攻击+B防御+A）*1.5*D*H*T*r*S*P
	damage := A.attack * (A.attack + coef) / math.Max(A.attack+B.defense+coef, 1) * 1.5 * D * H * T * r * S

	if c.isDebug {
		logrus.Debugf("TroopA index: %v stat:[%v %v %v %v %v %v]", troopA.index, A.attack, A.defense, A.strength, A.dexterity, A.damageIncrePer, A.damageDecrePer)
		logrus.Debugf("TroopB index: %v stat:[%v %v %v %v %v %v %v]", troopB.index, B.attack, B.defense, B.strength, B.dexterity, B.damageIncrePer, B.damageDecrePer, damage)
	}
	return damage
}

func (c *Combat) tryCrit(troopA, troopB *Troops, A, B *F64Stat) (bool, float64) {

	//A打B，若A敏捷≤B敏捷，不能触发暴击
	if A.dexterity <= B.dexterity {
		return false, 0
	}

	if c.tryRate(c.CritRate) {
		// 暴击伤害系数=（（2*A敏*A敏/（B敏*B敏+A敏*B敏））^0.5）*3-2
		coef := math.Sqrt(2*A.dexterity*A.dexterity/(B.dexterity*B.dexterity+A.dexterity*B.dexterity))*3 - 2
		return true, coef
	}

	return false, 0
}

func (c *Combat) tryDodge(troopA, troopB *Troops) bool {
	A := troopA.totalStat
	B := troopB.totalStat

	//A打B，若B敏捷≤A敏捷，不能触发闪避
	if B.dexterity <= A.dexterity {
		return false
	}

	// B的闪避几率=1-（（A敏*（A敏+ B敏）/（2*B敏*B敏））^0.5
	rate := 1 - math.Sqrt(A.dexterity*(A.dexterity+B.dexterity)/math.Max(2*B.dexterity*B.dexterity, 1))
	if rate <= 0 {
		return false
	}

	return c.random.Float64() < rate
}

// 攻击城墙
func (c *Combat) tryAttackWall(attacker *Troops, spell *SpellData, frame int, proto *shared_proto.CombatFrameProto) (suc bool) {
	enemy := c.getEnemy(attacker.isAttacker)
	wall := enemy.Wall()
	if wall == nil || !wall.isAlive() {
		logrus.Error("%d-%s 尝试攻击城墙，地方没有城墙或者城墙已死", attacker.index, attacker.proto.Captain.Name)
		return true
	}

	releaseSpellProto := &shared_proto.TroopReleaseSpellActionProto{
		ReleaseType: ReleaseTypeRelease,
		SpellId:     spell.Id(),
		EndFrame:    0, // 持续生效技能有值
		SpellX:      int32(attacker.x),
		SpellY:      int32(attacker.y),
	}
	addReleaseSpellAction(proto, attacker.index, releaseSpellProto)

	//普通伤害= A攻击*（A攻击+c）/（A攻击+城墙防御+c）*1.5*D*H*T*r*P
	//其中，防御系数D=1，体力系数H=1，加减伤系数P=1
	//T为A单位对城墙的克制系数，读配置表
	//打城墙不会触发闪避和暴击，百分比加减伤无效
	//对城墙的最终伤害=MAX{1，INT(普通伤害×攻击部队当前士兵数) }

	coef := c.Coef
	troopA := attacker
	A := troopA.totalStat
	B := wall

	D := float64(1)
	H := float64(1)
	T := troopA.raceData.wallCoef
	S := float64(1)
	r := 0.9 + c.random.Float64()*0.2

	//普通伤害= A攻击*（A攻击+c）/（A攻击+城墙防御+c）*1.5*D*H*T*r*P

	damage := A.attack * (A.attack + coef) / math.Max(A.attack+B.defense+coef, 1) * 1.5 * D * H * T * r * S

	//对城墙的最终伤害=MAX{1，INT(普通伤害×攻击部队当前士兵数) }
	hurtWallLife := int(math.Max(1, damage*float64(troopA.getSoldier())))

	// 城墙最大受伤血量（防止高级玩家一下就打爆了）
	if max := B.GetMaxBeenHurt(); hurtWallLife > max && max > 0 {
		hurtWallLife = max
	}

	targetProto := &shared_proto.TroopReleaseSpellEffectProto{
		TargetIndex: 0,
	}
	targetProto.IsTarget = true
	targetProto.HurtType = shared_proto.HurtType_NORMAL
	releaseSpellProto.Target = append(releaseSpellProto.Target, targetProto)

	if hurtWallLife > 0 {
		delayFrame := 0
		if spell.FlySpeedPerFrame > 0 {
			distance := IAbs(wall.x - attacker.x)
			if distance > 0 {
				delayFrame = distance / spell.FlySpeedPerFrame
			}
		}

		if delayFrame > 0 {
			effectFrame := frame + delayFrame
			targetProto.EffectFrame = int32(effectFrame)
			// 延时生效技能
			wall.addDelayDamage(attacker, hurtWallLife, effectFrame, spell.Id())
		} else {

			// 打城墙啊
			wall.ReduceLife(hurtWallLife)

			targetProto.ChangeSoldier = int32(hurtWallLife)
			targetProto.Soldier = int32(wall.Life())
		}
	}

	return !wall.isAlive()
}

// 城墙行动
func (c *Combat) tryWallAction(frame int, proto *shared_proto.CombatFrameProto) (checkEnd, isCheckAttacker bool) {

	// 判断防守方是不是有城墙还可以继续干
	wall := c.defender.Wall()
	if wall == nil {
		return false, false
	}

	if !wall.isAlive() {
		return true, false
	}

	// 处理延时伤害
	if wall.nextDelayDamageFrame <= frame {

		wall.nextDelayDamageFrame = math.MaxInt32
		if len(wall.delayDamage) > 0 {
			for i, dd := range wall.delayDamage {
				if dd == nil {
					continue
				}

				if dd.effectFrame <= frame {
					wall.delayDamage[i] = nil

					// 触发

					wall.ReduceLife(dd.damage)

					effectProto := &shared_proto.TroopTickEffectProto{}
					effectProto.Id = dd.spellId
					effectProto.ChangeSoldier = int32(IMax(dd.damage, 1))
					effectProto.Soldier = int32(wall.Life())

					addSpellEffectAction(proto, 0, effectProto)

					if !wall.isAlive() {
						return true, false
					}
				} else {
					wall.nextDelayDamageFrame = IMin(wall.nextDelayDamageFrame, dd.effectFrame)
				}
			}
		}
	}

	if frame < c.WallWaitFrame {
		return false, false
	}

	spell := c.WallSpell
	if frame < wall.nextReleaseFrame {
		// CD没到
		return false, false
	}

	var target *Troops // 最近的部队
	//randomIndex := rand.Intn(len(c.attacker.troops))
	for i := 0; i < len(c.attacker.troops); i++ {
		//targetTroop := c.attacker.troops[(i+randomIndex)%len(c.attacker.troops)]
		targetTroop := c.attacker.troops[i]
		if targetTroop == nil || !targetTroop.isAlive() {
			continue
		}

		if target == nil {
			target = targetTroop
		} else if targetTroop.x > target.x {
			target = targetTroop
		} else if targetTroop.x == target.x && targetTroop.y < target.y {
			// 优先打y小的
			target = targetTroop
		}
	}

	if target == nil {
		logrus.Errorf("城墙攻击阶段，进攻方没有可以被攻击的部队，什么鬼")
		return true, true
	}
	// 设置CD
	wall.nextReleaseFrame = frame + spell.CooldownFrame

	releaseSpellProto := &shared_proto.TroopReleaseSpellActionProto{
		ReleaseType: ReleaseTypeRelease,
		SpellId:     spell.Id(),
		EndFrame:    0, // 持续生效技能有值
		//SpellTarget: target.index,
		//SpellX:      int32(target.x),
		//SpellY:      int32(target.y),
	}
	addReleaseSpellAction(proto, 0, releaseSpellProto)

	reduceSoldier := 0

	wall.attackTimes++
	if wall.attackTimes <= c.WallAttackFixDamageTimes {
		// 固定死兵
		reduceSoldier = IMin(target.getSoldier(), wall.fixDamage)
	}

	// 如果固定死兵就被打死了，不计算其他技能伤害
	totalDamage := reduceSoldier * target.lifePerSoldier
	if reduceSoldier < target.getSoldier() {
		//床弩打B，依然使用上述伤害计算公式
		//普通伤害= 床弩攻击*（床弩攻击+c）/（床弩攻击+B防御+c）*1.5*D*H*T*r*P
		//其中，防御系数D=1，体力系数H=1，克制系数T=1，加减伤系数P=1
		//弩箭打人不会触发闪避和暴击，百分比加减伤无效
		//弩箭打人的伤害=MAX{1，INT(普通伤害×床弩统率) }
		//城防术附加死兵数量由科技属性配置

		coef := c.Coef

		D := float64(1)
		H := float64(1)
		T := float64(1)
		S := float64(1)
		r := 0.9 + c.random.Float64()*0.2

		//普通伤害= 床弩攻击*（床弩攻击+c）/（床弩攻击+B防御+c）*1.5*D*H*T*r*P
		A := wall
		troopB := target
		B := troopB.totalStat
		damage := A.attack * (A.attack + coef) / math.Max(A.attack+B.defense+coef, 1) * 1.5 * D * H * T * r * S

		//弩箭打人的伤害=MAX{1，INT(普通伤害×床弩统率) }
		wallDamage := int(math.Max(1, damage*A.soldierCapcity))

		if c.isDebug {
			logrus.Debugf("wall: %v %v %v %v %v %v %v", wallDamage, totalDamage, A.soldierCapcity, c.Coef, spell.Coef, A.attack, B.defense)
		}

		totalDamage += wallDamage
	}

	targetProto := target.newReleaseSpellEffectProtoIfNil(nil)
	targetProto.IsTarget = true
	releaseSpellProto.Target = append(releaseSpellProto.Target, targetProto)

	delayFrame := 0
	if spell.FlySpeedPerFrame > 0 {
		distance := IAbs(wall.x - target.x)
		if distance > 0 {
			delayFrame = distance / spell.FlySpeedPerFrame
		}
		// 最低
		delayFrame = IMax(delayFrame, c.WallDelayMinFrame)
	}

	if delayFrame > 0 {
		effectFrame := frame + delayFrame
		targetProto.EffectFrame = int32(effectFrame)
		// 延时生效技能
		target.addDelayDamage(wall, totalDamage, effectFrame, spell.Id())
	} else {

		changeShield, changeSoldier, showChangeSoldier, aliveSoldier, removeStateIds := c.tryReduceDamage(target, totalDamage, frame, proto)
		wall.addKillSoldier(changeSoldier)

		targetProto.HurtType = shared_proto.HurtType_NORMAL
		targetProto.ChangeShield = int32(changeShield)
		targetProto.RemoveState = removeStateIds
		targetProto.ChangeSoldier = int32(showChangeSoldier)
		targetProto.Soldier = int32(aliveSoldier)

		if !target.isAlive() {
			return true, true
		}
	}

	return false, false
}
