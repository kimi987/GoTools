package combat

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/imath"
	"github.com/pkg/errors"
	"math"
	"math/rand"
	"sort"
)

func newCombat(misc *misc, attacker, defender *CombatPlayer) (sc *Combat, err error) {
	sc = &Combat{
		misc:     misc,
		attacker: attacker,
		defender: defender,
		rounds:   make([]*shared_proto.CombatRoundProto, 0, imath.Max(len(attacker.troops), len(defender.troops))),
	}

	actionTroops := make([]*Troops, 0, len(attacker.troops)+len(defender.troops))
	actionTroops = append(actionTroops, attacker.troops...)
	actionTroops = append(actionTroops, defender.troops...)
	// 按照行动顺序排序
	sort.Sort(TroopsSpeedSlice(actionTroops))
	sc.actionTroops = actionTroops

	if misc.mapXLen*misc.mapYLen <= 1024 {
		sc.blockInfo = NewFastBlock(int(misc.mapXLen), int(misc.mapYLen))
	} else {
		sc.blockInfo = (*IterateBlockInfo)(sc)
	}

	for _, t := range sc.actionTroops {
		if !t.isAlive() {
			continue
		}

		t.resetPosition()

		// 设置单回合最大移动力
		sc.maxMovePerRound = imath.Max(sc.maxMovePerRound, t.moveTimesPerRound)

		// 设置站位阻挡
		if !sc.blockInfo.Walkable(t.x, t.y) {
			return nil, errors.Errorf("请求参数无效，部队位置重叠, pos: %v,%v", t.x, t.y)
		}

		sc.blockInfo.SetUnwalkable(t.x, t.y)
	}

	return
}

// 战斗
type Combat struct {
	*misc

	blockInfo BlockInfo

	// 按照行动顺序开始动(行动力高的前面，一样的话，攻方编号小的在前面)
	actionTroops []*Troops

	attacker *CombatPlayer
	defender *CombatPlayer

	maxMovePerRound int

	step int

	attackerWin bool

	rounds []*shared_proto.CombatRoundProto
}

func (c *Combat) isEnemyAllDead(isAttacker bool) (allDead bool) {
	return c.getEnemy(isAttacker).AllDead()
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

func (c *Combat) incStep() {
	c.step++
}

func (c *Combat) tryCrit(A, B *Troops) (bool, float64) {
	//A打B，若A敏捷≤B敏捷，不能触发暴击
	if A.dexterity <= B.dexterity {
		return false, 0
	}

	if c.random.Float64() < c.critRate {
		// 暴击伤害系数=（（2*A敏*A敏/（B敏*B敏+A敏*B敏））^0.5）*3-2
		coef := math.Sqrt(2*A.dexterity*A.dexterity/(B.dexterity*B.dexterity+A.dexterity*B.dexterity))*3 - 2
		return true, coef
	}

	return false, 0
}

func (c *Combat) tryDodge(A, B *Troops) bool {
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

// --- 战斗 ---
func (c *Combat) Calculate(fightQueue int) (combatResult *shared_proto.CombatProto) {
	logrus.Debugf("战斗队列：%d \n进攻方: %v\n 防守方: %v", fightQueue, c.attacker, c.defender)

	// 保存下获得士兵数量
	aliveSoliderWhenStart := c.getAliveSolider()

	rounds := make([]*shared_proto.CombatRoundProto, 0, imath.Max(len(c.attacker.troops), len(c.defender.troops)))

	if c.attacker.AllDead() {
		c.attackerWin = false
	} else if c.defender.AllDead() {
		c.attackerWin = true
	} else {
		// 最多进行这么多个回合数
		for round := 1; round <= c.maxRound; round++ {
			roundProto, attackerWin, isEnd := c.CalculateRound(int32(round))
			rounds = append(rounds, roundProto)
			if isEnd {
				c.attackerWin = attackerWin
				break
			}
		}
	}

	winner := c.getWinner()
	winner.IncreWinTimes()

	logrus.Debugf("战斗队列 %d 战斗结果: \n进攻方: %v\n 防守方: %v \n", fightQueue, c.attacker, c.defender)

	// 计算
	combatResult = &shared_proto.CombatProto{
		Attacker:                 c.attacker.proto,
		AttackerTroopPos:         c.attacker.TroopPos(),
		Defenser:                 c.defender.proto,
		DefenserTroopPos:         c.defender.TroopPos(),
		MapRes:                   c.mapRes,
		MapXLen:                  c.mapXLen,
		MapYLen:                  c.mapYLen,
		CombatSolider:            aliveSoliderWhenStart,
		FightQueue:               int32(fightQueue),
		Rounds:                   rounds,
		Score:                    c.calcScore(),
		AttackerWin:              c.isAttackerWin(),
		AliveSolider:             c.getAliveSolider(),
		KillSolider:              c.getKillSolider(),
		WinnerContinueWinTimes:   int32(winner.winTimes),
		IsWinnerContinueWinLeave: winner.IsContinueWinLeave(c.getContinueWinTimes(winner.isAttacker)),
	}

	return
}

func (c *Combat) isAttackerWin() (attackerWin bool) {
	return c.attackerWin
}

func (c *Combat) getWinner() (winner *CombatPlayer) {
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

func (c *Combat) getKillSolider() (aliveSolider []*shared_proto.Int32Pair) {
	if w := c.defender.Wall(); w != nil {
		if w.killSoldier > 0 {
			aliveSolider = append(aliveSolider, &shared_proto.Int32Pair{
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

		aliveSolider = append(aliveSolider, &shared_proto.Int32Pair{
			Key:   t.index,
			Value: int32(t.killSoldier),
		})
	}

	return
}

func (c *Combat) CalculateRound(round int32) (roundProto *shared_proto.CombatRoundProto, isAttackerWin, isFightEnd bool) {
	logrus.Debugf("进入第 %d 轮", round)

	roundProto = &shared_proto.CombatRoundProto{Round: round}

	// 新的回合开始，重置移动力
	for _, v := range c.actionTroops {
		if v.isAlive() {
			v.resetMoveTimes()
		}
	}

	for i := 0; i < c.maxMovePerRound; i++ {
		roundContinue := false

		// 根据机动力，依次行动
		for _, troops := range c.actionTroops {
			// 跳过 挂掉了 移动力为0
			if !troops.isAlive() || !troops.isMovable() {
				continue
			}

			actionProto := c.tryAction(troops, round)
			if isValidActionProto(actionProto) {
				actionProto.LeftMoveTimes = int32(troops.moveTimes)
				roundProto.Actions = append(roundProto.Actions, actionProto)
			}

			if c.isEnemyAllDead(troops.isAttacker) {
				// 这个家伙动完之后，没有敌人了，跳出循环，结算
				if troops.isAttacker {
					logrus.Debugf("防守方部队全部阵亡，进攻方胜利")
				} else {
					logrus.Debugf("进攻方部队全部阵亡，进攻方胜利")
				}

				isAttackerWin = troops.isAttacker
				isFightEnd = true
				return
			}

			if troops.isMovable() {
				// 还有可以动的单位，回合继续
				roundContinue = true
			}
		}

		if !roundContinue {
			break
		}
	}

	if c.tryWallAction(round, roundProto) {
		// 这个家伙动完之后，没有敌人了，跳出循环，结算
		logrus.Debugf("进攻方部队全部阵亡，防守方胜利")
		isAttackerWin = false
		isFightEnd = true
		return
	}

	// 防止错误，加上这个判断
	if c.attacker.AllDead() {
		logrus.Errorf("进攻方部队全部阵亡，防守方胜利，但是前面没判断结束")
		isAttackerWin = false
		isFightEnd = true
	} else if c.defender.AllDead() {
		logrus.Errorf("防守方部队全部阵亡，进攻方胜利，但是前面没判断结束")
		isAttackerWin = true
		isFightEnd = true
	}

	return
}

func (c *Combat) tryAction(attacker *Troops, round int32) (actionProto *shared_proto.CombatActionProto) {
	actionProto = &shared_proto.CombatActionProto{Index: attacker.index}

	// 找目标攻击
	if c.tryAttack(attacker, round, actionProto) {
		return
	}

	// 找目标移动
	if !c.tryMove(attacker, actionProto) {
		// 移动失败
		return
	}

	// 移动成功了，步骤也要算一次
	c.incStep()

	// 移动成功，再找人打一下
	c.tryAttack(attacker, round, actionProto)

	return
}

func (c *Combat) tryMove(attacker *Troops, proto *shared_proto.CombatActionProto) (suc bool) {
	logrus.Debugf("%d-%s 尝试移动", attacker.index, attacker.proto.Captain.Name)
	if c.tryMoveToEnemyTroop(attacker, proto) {
		return true
	} else if c.tryMoveToWall(attacker, proto) {
		return true
	}

	logrus.Errorf("%d-%s 寻找移动目标，但是目标没有找到（这里不应该进来，打人的地方就应该处理掉）", attacker.index, attacker.proto.Captain.Name)
	attacker.moveTimes = 0
	return
}

func (c *Combat) tryMoveToWall(attacker *Troops, proto *shared_proto.CombatActionProto) (suc bool) {
	enemy := c.getEnemy(attacker.isAttacker)
	if !enemy.HasAliveWall() {
		// 对方没城墙或者城墙已死
		logrus.Debugf("%d-%s 尝试移动到敌方城墙旁边，地方没有城墙或者城墙已经挂了", attacker.index, attacker.proto.Captain.Name)
		return
	}

	// 先将移动次数减掉
	attacker.moveTimes--

	oldX, oldY := attacker.x, attacker.y
	moveX := attacker.x
	moveY := attacker.y

	direction := shared_proto.Direction_ORIGIN

	// 走到城墙那里去
	if attacker.isAttacker {
		moveX++
		direction = shared_proto.Direction_RIGHT
	} else {
		moveX--
		direction = shared_proto.Direction_LEFT
	}

	// 可以走，那就走
	if !c.tryMoveTroops(attacker, moveX, moveY) {
		return
	}

	proto.MoveDirection = direction

	logrus.Debugf("%d-%s 往 %+v 方向移动去城墙，从 %d, %d 移动到了 %d, %d", attacker.index, attacker.proto.Captain.Name, direction, oldX, oldY, moveX, moveY)
	return true
}

// 移动到敌对部队去
func (c *Combat) tryMoveToEnemyTroop(attacker *Troops, proto *shared_proto.CombatActionProto) (suc bool) {
	moveToTroops := c.findCanMoveToTroop(attacker)
	if moveToTroops == nil {
		// 没找到敌对部队
		logrus.Debugf("%d-%s 尝试移动到敌方部队旁边，敌方没有活的部队了", attacker.index, attacker.proto.Captain.Name)
		return
	}

	// 先将移动次数减掉，这后面的都算成功
	attacker.moveTimes--

	oldX, oldY := attacker.x, attacker.y
	moveX := attacker.x
	moveY := attacker.y

	suc = true

	direction := shared_proto.Direction_ORIGIN

	// 往目标方向移动，优先横向移动，纵坐标一致时候，
	if attacker.x != moveToTroops.x {
		if attacker.x < moveToTroops.x {
			moveX++
			direction = shared_proto.Direction_RIGHT
		} else {
			moveX--
			direction = shared_proto.Direction_LEFT
		}

		// 可以走，那就走
		if c.tryMoveTroops(attacker, moveX, moveY) {
			proto.MoveDirection = direction
			logrus.Debugf("%d-%s 往 %+v 方向移动去 %d-%s，从 %d, %d 移动到了 %d, %d", attacker.index, attacker.proto.Captain.Name, direction, moveToTroops.index, moveToTroops.proto.Captain.Name, oldX, oldY, moveX, moveY)
			return
		}

		// 重置回来
		moveX = attacker.x
	}

	if attacker.y != moveToTroops.y {
		if attacker.y < moveToTroops.y {
			moveY++
			direction = shared_proto.Direction_DOWN // Y轴反向
		} else {
			moveY--
			direction = shared_proto.Direction_UP // Y轴反向
		}

		// 可以走，那就走
		if c.tryMoveTroops(attacker, moveX, moveY) {
			proto.MoveDirection = direction
			logrus.Debugf("%d-%s 往 %+v 方向移动去 %d-%s，从 %d, %d 移动到了 %d, %d", attacker.index, attacker.proto.Captain.Name, direction, moveToTroops.index, moveToTroops.proto.Captain.Name, oldX, oldY, moveX, moveY)
			return
		}
	}

	// 移动失败，可能我想走的方向站了人
	logrus.Debugf("%d-%s 往地方部队%d-%s移动失败，周围的点都无法移动", attacker.index, attacker.proto.Captain.Name, moveToTroops.index, moveToTroops.proto.Captain.Name)
	return
}

func (c *Combat) findCanMoveToTroop(attacker *Troops) (moveToTroops *Troops) {
	var currentTargetRange int
	var currentPriority int

	for _, enemyTroop := range c.getEnemy(attacker.isAttacker).troops {
		if !enemyTroop.isAlive() {
			continue
		}

		r := attacker.getRange(enemyTroop)

		moveToTroops, currentTargetRange, currentPriority = c.selectTarget(attacker, moveToTroops, currentTargetRange, currentPriority, enemyTroop, r)
	}

	return
}

func (c *Combat) tryAttack(attacker *Troops, round int32, proto *shared_proto.CombatActionProto) (suc bool) {
	logrus.Debugf("%d-%s 尝试攻击", attacker.index, attacker.proto.Captain.Name)
	if suc, hasAliveTroop := c.tryAttackTroop(attacker, round, proto); suc {
		// 攻击成功
		return true
	} else if hasAliveTroop {
		// 攻击敌方部队失败，但是还有活着的敌方部队
		logrus.Debugf("%d-%s 尝试攻击，敌方活着的不对没有在我的攻击范围的", attacker.index, attacker.proto.Captain.Name)
		return false
	} else {
		// 对方没得部队可以攻击了
		return c.tryAttackWall(attacker, round, proto)
	}
}

func (c *Combat) tryAttackWall(attacker *Troops, round int32, proto *shared_proto.CombatActionProto) (suc bool) {
	logrus.Debugf("%d-%s 尝试攻击城墙", attacker.index, attacker.proto.Captain.Name)
	enemy := c.getEnemy(attacker.isAttacker)
	wall := enemy.Wall()
	if wall == nil || !wall.isAlive() {
		logrus.Debugf("%d-%s 尝试攻击城墙，地方没有城墙或者城墙已死", attacker.index, attacker.proto.Captain.Name)
		return false
	}

	// 对方的兵都死了，那就看对方的城墙了
	if attacker.isAttacker {
		if attacker.attackRange < (int(c.mapXLen) - attacker.x) {
			logrus.Debugf("%d-%s 尝试攻击城墙，城墙超出攻击距离了, x: %d, y: %d, range: %d", attacker.index, attacker.proto.Captain.Name, attacker.x, attacker.y, attacker.attackRange)
			return false
		}
	} else {
		if attacker.attackRange < attacker.x+1 {
			logrus.Debugf("%d-%s 尝试攻击城墙，城墙超出攻击距离了, x: %d, y: %d, range: %d", attacker.index, attacker.proto.Captain.Name, attacker.x, attacker.y, attacker.attackRange)
			return false
		}
	}

	//普通伤害= A攻击*（A攻击+c）/（A攻击+城墙防御+c）*1.5*D*H*T*r*P
	//其中，防御系数D=1，体力系数H=1，加减伤系数P=1
	//T为A单位对城墙的克制系数，读配置表
	//打城墙不会触发闪避和暴击，百分比加减伤无效
	//对城墙的最终伤害=MAX{1，INT(普通伤害×攻击部队当前士兵数) }

	coef := c.coef
	A := attacker
	B := wall

	D := float64(1)
	H := float64(1)
	T := A.wallCoef
	S := float64(1)
	r := 0.9 + c.random.Float64()*0.2

	//普通伤害= A攻击*（A攻击+c）/（A攻击+城墙防御+c）*1.5*D*H*T*r*P

	damage := A.attack * (A.attack + coef) / math.Max(A.attack+B.defense+coef, 1) * 1.5 * D * H * T * r * S

	//对城墙的最终伤害=MAX{1，INT(普通伤害×攻击部队当前士兵数) }
	hurtWallLife := int(math.Max(1, damage*float64(A.soldier)))

	// 城墙最大受伤血量（防止高级玩家一下就打爆了）
	if max := B.GetMaxBeenHurt(); hurtWallLife > max && max > 0 {
		hurtWallLife = max
	}

	// 打城墙啊
	enemy.Wall().ReduceLife(hurtWallLife)

	logrus.Debugf("%d-%s 攻击城墙，城墙掉血: %d，剩余血量: %d", attacker.index, attacker.proto.Captain.Name, hurtWallLife, enemy.Wall().Life())

	proto.HurtWallLife = int32(hurtWallLife)
	proto.WallLeftLife = int32(enemy.Wall().Life())
	proto.HurtWallSpell = attacker.raceProto.GetNormalSpellId()

	attacker.moveTimes = 0 // 打到人/墙，不走了

	c.incStep()

	return true
}

func (c *Combat) tryAttackTroop(attacker *Troops, round int32, proto *shared_proto.CombatActionProto) (suc, hasAliveTroop bool) {
	enemyTroop, hasAliveTroop := c.findCanAttackTroop(attacker)
	if enemyTroop == nil {
		logrus.Debugf("%d-%s 尝试攻击，没有找到可以攻击的目标", attacker.index, attacker.proto.Captain.Name)
		return
	}

	logrus.Debugf("%d-%s 尝试攻击，找到攻击目标: %d-%s", attacker.index, attacker.proto.Captain.Name, enemyTroop.index, enemyTroop.proto.Captain.Name)

	c.triggerRestraintSpell(attacker, round, enemyTroop, proto)

	proto.TargetIndex = enemyTroop.index
	proto.NormalSpell = attacker.raceProto.GetNormalSpellId()

	// 闪避
	if c.tryDodge(attacker, enemyTroop) {
		proto.HurtType = shared_proto.HurtType_MISS

		logrus.Debugf("%d-%s 攻击(%d-%s)，目标闪避了 ", attacker.index, attacker.proto.Captain.Name, enemyTroop.index, enemyTroop.proto.Captain.Name)
	} else {
		// 计算伤害
		damage, hurtType := c.calculateDamage(attacker, enemyTroop)
		reduceSoldier := c.tryReduceLife(enemyTroop, damage)
		attacker.killSoldier += reduceSoldier

		proto.Damage += int32(reduceSoldier)
		proto.HurtType = hurtType

		logrus.Debugf("%d-%s 攻击(%d-%s)，释放技能: %d, 伤害: %v，杀死士兵数: %v，剩余士兵数: %v", attacker.index, attacker.proto.Captain.Name, enemyTroop.index, enemyTroop.proto.Captain.Name, proto.NormalSpell, damage, reduceSoldier, enemyTroop.Soldier())
	}

	proto.LeftSoldier = int32(enemyTroop.Soldier())

	attacker.moveTimes = 0 // 打到人/墙，不走了

	c.incStep()
	suc = true

	return
}

// 计算克制技伤害
func (c *Combat) calculateRestraintDamage(A, B *Troops) int {
	damage := c.calculateNormalDamage(A, B)

	//若A敏捷<B敏捷，
	//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数×（（A敏*（A敏+ B敏）/（2*B敏*B敏））^0.5)；
	//若A敏捷=B敏捷，
	//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数)；
	//若A敏捷>B敏捷，
	//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数×（2*A敏*A敏/（B敏*B敏+A敏*B敏））^0.5)
	//以上，技能伤害系数=0.2
	//克制技伤害不能被闪避，也不会暴击。

	switch {
	case A.dexterity < B.dexterity:
		//若A敏捷<B敏捷，
		//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数×（（A敏*（A敏+ B敏）/（2*B敏*B敏））^0.5)；
		coef := math.Sqrt(A.dexterity * (A.dexterity + B.dexterity) / (2 * B.dexterity * B.dexterity))
		return int(math.Max(1, damage*c.restraintRate*float64(A.soldier)*coef))
	case A.dexterity > B.dexterity:
		//若A敏捷>B敏捷，
		//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数×（2*A敏*A敏/（B敏*B敏+A敏*B敏））^0.5)
		coef := math.Sqrt(2 * A.dexterity * A.dexterity / (B.dexterity*B.dexterity + A.dexterity*B.dexterity))
		return int(math.Max(1, damage*c.restraintRate*float64(A.soldier)*coef))
	default:
		//若A敏捷=B敏捷，
		//克制技伤害= INT(当次普通伤害×技能伤害系数×A当前士兵数)；
		return int(math.Max(1, damage*c.restraintRate*float64(A.soldier)))
	}
}

// 触发克制技能
func (c *Combat) triggerRestraintSpell(attacker *Troops, round int32, enemy *Troops, proto *shared_proto.CombatActionProto) {
	if !attacker.canTriggerRestraintSpell {
		// 进攻方不可以释放克制技
		logrus.Debugf("%d-%s 攻击(%d-%s)，未触发克制技, 不可以触发克制技", attacker.index, attacker.proto.Captain.Name, enemy.index, enemy.proto.Captain.Name)
		return
	}

	// 克制技
	if !attacker.isRestraintSpellTriggerRound(round) {
		// 不是克制轮
		logrus.Debugf("%d-%s 攻击(%d-%s)，未触发克制技, 不是克制轮(%d)", attacker.index, attacker.proto.Captain.Name, enemy.index, enemy.proto.Captain.Name, round)
		return
	}

	if !attacker.isRestraintRace(enemy.raceProto.Race) {
		// 不是克制职业
		logrus.Debugf("%d-%s 攻击(%d-%s)，未触发克制技, 不是克制职业", attacker.index, attacker.proto.Captain.Name, enemy.index, enemy.proto.Captain.Name)
		return
	}

	spellDamage := c.calculateRestraintDamage(attacker, enemy)
	proto.RestraintSpell = attacker.raceProto.GetRestraintSpellId()

	if enemy.isAlive() {
		reduceSoldier := c.tryReduceLife(enemy, spellDamage)
		attacker.killSoldier += reduceSoldier

		proto.RestraintSpellDamage = int32(reduceSoldier)
		proto.Damage += int32(reduceSoldier)

		logrus.Debugf("%d-%s 攻击(%d-%s)，触发克制技, 伤害: %v，杀死士兵数: %v，剩余士兵数: %v", attacker.index, attacker.proto.Captain.Name, enemy.index, enemy.proto.Captain.Name, spellDamage, reduceSoldier, enemy.Soldier())
	} else {
		// 这里伤害都是士兵数，对方都已经死了的，那就都是0
		logrus.Debugf("%d-%s 攻击(%d-%s)，触发克制技, 但是未对目标造成伤害", attacker.index, attacker.proto.Captain.Name, enemy.index, enemy.proto.Captain.Name)
	}

}

func (c *Combat) findCanAttackTroop(attacker *Troops) (enemyTroop *Troops, hasAliveTroop bool) {
	var currentTargetRange int
	var currentPriority int

	for _, troop := range c.getEnemy(attacker.isAttacker).troops {
		if !troop.isAlive() {
			continue
		}

		hasAliveTroop = true

		r := attacker.getRange(troop)
		if attacker.attackRange < r {
			// 攻击范围不够
			continue
		}

		enemyTroop, currentTargetRange, currentPriority = c.selectTarget(attacker, enemyTroop, currentTargetRange, currentPriority, troop, r)
	}

	return
}

func (c *Combat) selectTarget(attacker *Troops, currentTarget *Troops, currentTargetRange int, currentPriority int, newTarget *Troops, newTargetRange int) (*Troops, int, int) {
	// 当前没有目标，则选为初始目标
	if currentTarget == nil {
		return newTarget, newTargetRange, c.getTargetPriority(attacker.race, newTarget.race)
	}

	// 当前有目标，则选择优先级高的为目标
	newPriority := currentPriority
	if currentTarget.race != newTarget.race {
		newPriority = c.getTargetPriority(attacker.race, newTarget.race)
		if currentPriority < newPriority {
			return newTarget, newTargetRange, newPriority
		} else {
			return currentTarget, currentTargetRange, currentPriority
		}
	}

	// 优先级相同，则选择距离小的为目标
	if currentTargetRange != newTargetRange {
		if currentTargetRange < newTargetRange {
			return currentTarget, currentTargetRange, currentPriority
		} else {
			return newTarget, newTargetRange, newPriority
		}
	}

	// 距离相同，则选择index小的为目标
	if currentTarget.index < newTarget.index {
		return currentTarget, currentTargetRange, currentPriority
	} else {
		return newTarget, newTargetRange, newPriority
	}
}

// 城墙攻击
func (c *Combat) tryWallAction(round int32, roundProto *shared_proto.CombatRoundProto) (attackerAllDead bool) {
	logrus.Debugf("进入防守方城墙进攻阶段")

	if round < c.minWallAttackRound {
		logrus.Debugf("当前轮次 %d 不在城墙进攻轮 %d-%d次", round, c.minWallAttackRound, c.maxWallAttachFixDamageRound)
		return false
	}

	// 判断防守方是不是有城墙还可以继续干
	wall := c.defender.Wall()
	if wall == nil || !wall.isAlive() {
		logrus.Debugf("防守方没有城墙或者城墙已经死亡")
		return
	}

	var nearlyTroop *Troops // 最近的部队

	randomIndex := rand.Intn(len(c.attacker.troops))
	for i := 0; i < len(c.attacker.troops); i++ {
		targetTroop := c.attacker.troops[(i+randomIndex)%len(c.attacker.troops)]
		if targetTroop == nil || !targetTroop.isAlive() {
			continue
		}

		if nearlyTroop == nil {
			nearlyTroop = targetTroop
		} else if targetTroop.x > nearlyTroop.x {
			nearlyTroop = targetTroop
		} else if targetTroop.x == nearlyTroop.x && targetTroop.y < nearlyTroop.y {
			// 优先打y小的
			nearlyTroop = targetTroop
		}
	}

	if nearlyTroop == nil {
		logrus.Errorf("城墙攻击阶段，进攻方没有可以被攻击的部队，什么鬼")
		return
	}

	//床弩打B，依然使用上述伤害计算公式
	//普通伤害= 床弩攻击*（床弩攻击+c）/（床弩攻击+B防御+c）*1.5*D*H*T*r*P
	//其中，防御系数D=1，体力系数H=1，克制系数T=1，加减伤系数P=1
	//弩箭打人不会触发闪避和暴击，百分比加减伤无效
	//弩箭打人的伤害=MAX{1，INT(普通伤害×床弩统率) }
	//城防术附加死兵数量由科技属性配置

	coef := c.coef

	D := float64(1)
	H := float64(1)
	T := float64(1)
	S := float64(1)
	r := 0.9 + c.random.Float64()*0.2

	//普通伤害= 床弩攻击*（床弩攻击+c）/（床弩攻击+B防御+c）*1.5*D*H*T*r*P
	A := wall
	B := nearlyTroop
	damage := A.attack * (A.attack + coef) / math.Max(A.attack+B.defense+coef, 1) * 1.5 * D * H * T * r * S

	//弩箭打人的伤害=MAX{1，INT(普通伤害×床弩统率) }
	totalDamage := int(math.Max(1, damage*A.soldierCapcity))

	reduceSoldier := c.tryReduceLife(nearlyTroop, totalDamage)

	if round <= c.maxWallAttachFixDamageRound {
		// 固定死兵
		d := c.tryReduceSoldier(nearlyTroop, c.defender.Wall().fixDamage)
		reduceSoldier += d
	}
	wall.killSoldier += reduceSoldier

	roundProto.WallAction = &shared_proto.WallActionProto{
		TargetIndex: nearlyTroop.index,
		Damage:      int32(reduceSoldier),
		LeftSoldier: int32(nearlyTroop.Soldier()),
	}

	logrus.Debugf("城墙攻击阶段，目标(%d-%s)死亡士兵数: %d，目标剩余士兵数量: %d", nearlyTroop.index, nearlyTroop.proto.Captain.Name, damage, nearlyTroop.Soldier())

	return c.attacker.AllDead()
}

func (c *Combat) calculateNormalDamage(A, B *Troops) float64 {

	coef := c.coef

	// A打B
	// 防御系数D=（（A体力+A）* （B体力+A））^0.5 /（B防御+A）
	D := math.Sqrt((A.strength+coef)*(B.strength+coef)) / math.Max(B.defense+coef, 1)

	// 体力系数H，如果A体力≤B体力，则H=1.00，否则H=arctan（A体力 / B体力）*1.083+0.15
	H := 1.0
	if A.strength > B.strength {
		H = math.Atan(A.strength/math.Max(B.strength, 1))*1.083 + 0.15
	}

	// T为兵种克制系数，读表获得
	T := math.Max(c.getTroopsCoef(A, B), 0.1)

	//// 士气系数S（武将属性）（已过期，不再使用）
	//S := math.Max(A.Scoef, 1)

	// 加减伤系数，S=（1+A的加伤1+A的加伤2+…）/（1+B的减伤1+B的减伤2+…）
	S := (1 + A.damageIncrePer) / (1 + B.damageDecrePer)

	r := 0.9 + c.random.Float64()*0.2

	// 普通伤害= A攻击*（A攻击+A）/（A攻击+B防御+A）*1.5*D*H*T*r*S*P
	damage := A.attack * (A.attack + coef) / math.Max(A.attack+B.defense+coef, 1) * 1.5 * D * H * T * r * S

	return damage
}

func (c *Combat) calculateDamage(A, B *Troops) (int, shared_proto.HurtType) {

	damage := c.calculateNormalDamage(A, B)
	hurtType := shared_proto.HurtType_NORMAL

	// 计算暴击
	isCrit, critCoef := c.tryCrit(A, B)
	if isCrit {
		damage = damage * math.Max(1, critCoef)
		hurtType = shared_proto.HurtType_CRIT
	}

	// 最终伤害=MAX{1，INT(单兵伤害×攻击部队当前士兵数) }
	return int(math.Max(1, damage*float64(A.soldier))), hurtType
}

func (c *Combat) tryMoveTroops(troops *Troops, x, y int) bool {
	b := c.blockInfo
	if !b.Walkable(x, y) {
		return false
	}

	b.SetUnwalkable(x, y)
	b.SetWalkable(troops.x, troops.y)

	troops.x = x
	troops.y = y
	return true
}

func (c *Combat) tryReduceLife(troops *Troops, damage int) (reduceSoldier int) {
	reduceSoldier = troops.reduceLife(damage)

	if !troops.isAlive() {
		c.blockInfo.SetWalkable(troops.x, troops.y)
	}

	return
}

func (c *Combat) tryReduceSoldier(troops *Troops, toReduce int) (reduceSoldier int) {
	reduceSoldier = troops.reduceSoldier(toReduce)

	if !troops.isAlive() {
		c.blockInfo.SetWalkable(troops.x, troops.y)
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

	percent := uint64(aliveSoldier * 100 / totalSoldier)
	score = c.getScore(percent)

	logrus.Debugf("计算评分，alive: %d total: %d percent: %d score: %d", aliveSoldier, totalSoldier, percent, score)
	return
}

func isValidActionProto(proto *shared_proto.CombatActionProto) bool {
	if proto == nil || proto.Index <= 0 {
		return false
	}

	// 现在要更新行动力
	return true

	//if proto.MoveDirection != shared_proto.Direction_ORIGIN {
	//	return true
	//}
	//
	//if proto.TargetIndex != 0 {
	//	return true
	//}
	//
	//if proto.HurtWallSpell != 0 {
	//	return true
	//}
	//
	//return false
}

type IterateBlockInfo Combat

func (b *IterateBlockInfo) Walkable(x, y int) bool {
	for _, t := range b.actionTroops {
		if t.isAlive() && t.x == x && t.y == y {
			return false
		}
	}

	return true
}

func (b *IterateBlockInfo) SetWalkable(x, y int)   {}
func (b *IterateBlockInfo) SetUnwalkable(x, y int) {}
