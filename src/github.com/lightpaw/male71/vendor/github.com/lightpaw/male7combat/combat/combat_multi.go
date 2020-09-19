package combat

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/imath"
	"github.com/pkg/errors"
	"math"
)

func NewMultiCombat(p *server_proto.MultiCombatRequestServerProto) (*MultiCombat, error) {

	if p == nil {
		return nil, errors.Errorf("请求参数无效，传入的CombatRequestProto为nil")
	}

	if len(p.Attacker) == 0 {
		return nil, errors.Errorf("请求参数无效，传入的request.Attacker == 0")
	}

	if len(p.AttackerId) != len(p.Attacker) {
		return nil, errors.Errorf("请求参数无效，传入的 request.AttackerId(%d) 的长度 != request.Attacker(%d) 的长度", len(p.AttackerId), len(p.Attacker))
	}

	for _, attacker := range p.Attacker {
		if attacker == nil {
			return nil, errors.Errorf("请求参数无效，传入的request.Attacker中存在nil")
		}

		if len(attacker.Troops) == 0 {
			return nil, errors.Errorf("请求参数无效，传入的request.Attacker.Troops == 0")
		}
	}

	if len(p.Defenser) == 0 {
		return nil, errors.Errorf("请求参数无效，传入的request.Defenser == 0")
	}

	if len(p.DefenserId) != len(p.Defenser) {
		return nil, errors.Errorf("请求参数无效，传入的 request.DefenserId(%d) 的长度 != request.Defenser(%d) 的长度", len(p.DefenserId), len(p.Defenser))
	}

	for _, defenser := range p.Defenser {
		if defenser == nil {
			return nil, errors.Errorf("请求参数无效，传入的request.Defenser中存在nil")
		}

		if len(defenser.Troops) == 0 && defenser.WallStat == nil {
			return nil, errors.Errorf("请求参数无效，传入的request.Defenser.Troops == 0")
		}
	}

	if p.MaxRound <= 0 {
		return nil, errors.Errorf("请求参数无效，传入的request.MaxRound <= 0")
	}

	if len(p.Races) <= 0 {
		return nil, errors.Errorf("请求参数无效，传入的request.Races 长度 <= 0")
	}

	if p.GetMapXLen() <= 0 || p.GetMapXLen() > math.MaxInt16 {
		return nil, errors.Errorf("请求参数无效，传入的request.MapXLen should in [0, %v], v: %v", math.MaxInt16, p.GetMapXLen())
	}
	if p.GetMapYLen() <= 0 || p.GetMapYLen() > math.MaxInt16 {
		return nil, errors.Errorf("请求参数无效，传入的request.MapYLen should in [0, %v], v: %v", math.MaxInt16, p.GetMapYLen())
	}

	if p.GetConcurrentFightCount() <= 0 {
		return nil, errors.Errorf("请求参数无效，传入的request.ConcurrentFightCount should > 0, v: %v", p.GetConcurrentFightCount())
	}

	if p.GetConcurrentFightCount() > int32(imath.Max(len(p.Attacker), len(p.Defenser))) {
		return nil, errors.Errorf("请求参数无效，传入的request.ConcurrentFightCount should <=, v: %v", imath.Max(len(p.Attacker), len(p.Defenser)))
	}

	combat := &MultiCombat{
		misc:                 newMultiMisc(p),
		attackerPlayers:      make([]*CombatPlayer, 0, len(p.Attacker)),
		defenserPlayers:      make([]*CombatPlayer, 0, len(p.Defenser)),
		proto:                p,
		concurrentFightCount: int(p.GetConcurrentFightCount()),
		combatQueueWinner:    make([]*CombatPlayer, p.GetConcurrentFightCount()),
	}

	for idx, player := range p.Attacker {
		err := combat.addPlayer(p.AttackerId[idx], player, true)
		if err != nil {
			return nil, errors.Wrap(err, "添加进攻方失败")
		}
	}

	for idx, player := range p.Defenser {
		err := combat.addPlayer(p.DefenserId[idx], player, false)
		if err != nil {
			return nil, errors.Wrap(err, "添加防守方失败")
		}
	}

	return combat, nil
}

type MultiCombat struct {
	*misc

	proto *server_proto.MultiCombatRequestServerProto

	// 并发战斗数量
	concurrentFightCount int

	index int32

	// 进攻方/防守方
	attackerPlayers []*CombatPlayer
	defenserPlayers []*CombatPlayer

	// 下一次选择出战的人
	nextSelectAttackerIndex int
	nextSelectDefenserIndex int

	// 进攻方/防守方死亡数量
	attackerDeadCount int32
	defenserDeadCount int32

	// 进攻方/防守方离场数量
	attackerLeaveCount int32
	defenserLeaveCount int32

	// 进攻方/防守方胜利最大步数
	attackerWinMaxStep int
	defenserWinMaxStep int

	isAttackerWin bool

	// 战斗队列胜利方
	combatQueueWinner []*CombatPlayer
}

func (c *MultiCombat) genIndex() int32 {
	c.index++
	return c.index
}

func (c *MultiCombat) getNextPlayer(isAttacker bool) (nextPlayer *CombatPlayer) {
	if isAttacker {
		nextPlayer = c.attackerPlayers[c.nextSelectAttackerIndex]
		c.nextSelectAttackerIndex++
	} else {
		nextPlayer = c.defenserPlayers[c.nextSelectDefenserIndex]
		c.nextSelectDefenserIndex++
	}

	return
}

// 处理一次战斗结束后的计算
func (c *MultiCombat) onSingleFightEnd(isAttackerWin bool, isWinnerContinueWinLeave bool, winnerMaxStep int) {
	if isAttackerWin {
		c.defenserDeadCount++
		c.attackerWinMaxStep = imath.Max(c.attackerWinMaxStep, winnerMaxStep)
		if isWinnerContinueWinLeave {
			c.attackerLeaveCount++
		}
	} else {
		c.attackerDeadCount++
		c.defenserWinMaxStep = imath.Max(c.defenserWinMaxStep, winnerMaxStep)
		if isWinnerContinueWinLeave {
			c.defenserLeaveCount++
		}
	}
}

func (c *MultiCombat) addPlayer(id int64, p *shared_proto.CombatPlayerProto, isAttacker bool) error {
	player, err := newCombatPlayer(c.genIndex, id, isAttacker, p, c.misc)
	if err != nil {
		return errors.Wrap(err, "MultiCombat.addPlayer失败")
	}

	if isAttacker {
		c.attackerPlayers = append(c.attackerPlayers, player)
	} else {
		c.defenserPlayers = append(c.defenserPlayers, player)
	}

	return nil
}

func (c *MultiCombat) Calculate() (result *shared_proto.MultiCombatProto, err error) {
	// 最大轮次=双方人数的和
	maxRound := len(c.attackerPlayers) + len(c.defenserPlayers)

	combats := make([]*shared_proto.CombatProto, 0, maxRound)

	isEnd := false

	for r := 0; r < maxRound; r++ {
		combatResult, err := c.fight()
		if err != nil {
			return nil, errors.Wrap(err, "多人战斗失败")
		}

		if combatResult == nil {
			return nil, errors.New("combatResult 为空")
		}

		combats = append(combats, combatResult)

		if c.checkIsEnd() {
			isEnd = true
			logrus.Debugf("战斗结束了，进攻方死亡人数: %d, 防守方死亡人数: %d，进攻方连胜离场人数: %d, 防守方连胜离场人数: %d，进攻方最大胜利步数: %d, 防守方最大胜利步数: %d",
				c.attackerDeadCount, c.defenserDeadCount, c.attackerLeaveCount, c.defenserLeaveCount, c.attackerWinMaxStep, c.defenserWinMaxStep)
			break
		}
	}

	if !isEnd {
		return nil, errors.New("服务器计算是否结束所有轮次都跑完了，检查checkIsEnd 都没有结束，直接是跳出循环的")
	}

	logrus.Debugf("战斗结束，最终胜利方是: %v", c.isAttackerWin)

	// 处理
	return &shared_proto.MultiCombatProto{
		Attacker:             c.proto.Attacker,
		Defenser:             c.proto.Defenser,
		MapRes:               c.proto.MapRes,
		MapXLen:              c.proto.MapXLen,
		MapYLen:              c.proto.MapYLen,
		ConcurrentFightCount: int32(c.concurrentFightCount),
		AttackerWin:          c.isAttackerWin,
		Combats:              combats,
		Score:                c.calcScore(),
		AttackerDeadCount:    c.attackerDeadCount,
		DefenserDeadCount:    c.defenserDeadCount,
		AttackerLeaveCount:   c.attackerLeaveCount,
		DefenserLeaveCount:   c.defenserLeaveCount,
		WinTimes:             c.getWinTimesList(),
	}, nil
}

func (c *MultiCombat) fight() (combatResult *shared_proto.CombatProto, err error) {
	combat, fightQueue, err := c.newCombat()
	if err != nil {
		return
	}

	combatResult = combat.Calculate(fightQueue)

	winner := combat.getWinner()

	c.combatQueueWinner[fightQueue-1] = winner

	totalQueueStep := winner.CurStep() + combat.step
	combat.attacker.SetCurStep(totalQueueStep)
	combat.defender.SetCurStep(totalQueueStep)

	c.onSingleFightEnd(combat.isAttackerWin(), winner.IsContinueWinLeave(c.getContinueWinTimes(winner.isAttacker)), totalQueueStep)

	logrus.Debugf("战斗队列 %d 胜利方id: %d, 是进攻方: %v，消耗: %d 步，该战斗队列总共使用 %d 步，连胜次数: %d，是否连胜离场: %t \n",
		fightQueue, winner.Id(), winner.isAttacker, combat.step, totalQueueStep, winner.winTimes, winner.IsContinueWinLeave(c.getContinueWinTimes(winner.isAttacker)))

	return
}

func (c *MultiCombat) newCombat() (combat *Combat, fightQueue int, err error) {
	// 找到当前有空位的位置
	attacker, defenser, fightQueue, err := c.getFightQueue()
	if err != nil {
		return
	}

	if attacker == nil || defenser == nil {
		err = errors.New("没找到进攻方或者防守方")
		return
	}

	if fightQueue <= 0 {
		err = errors.New("没找到战斗队列")
		return
	}

	combat, err = newCombat(c.misc, attacker, defenser)
	if err != nil {
		err = errors.Wrap(err, "newCombat出错")
		return
	}

	oldWinner := c.combatQueueWinner[fightQueue-1]
	if oldWinner != nil {
		attacker.SetCurStep(oldWinner.CurStep())
		defenser.SetCurStep(oldWinner.CurStep())
	}

	if attacker.lastFightQueue != 0 && attacker.lastFightQueue != fightQueue {
		c.combatQueueWinner[attacker.lastFightQueue-1] = nil
	}

	if defenser.lastFightQueue != 0 && defenser.lastFightQueue != fightQueue {
		c.combatQueueWinner[defenser.lastFightQueue-1] = nil
	}

	attacker.lastFightQueue = fightQueue
	defenser.lastFightQueue = fightQueue

	return
}

// 获得战斗队列
func (c *MultiCombat) getFightQueue() (attacker, defenser *CombatPlayer, fightQueue int, err error) {
	hasAttackerNextCanSelect := c.nextSelectAttackerIndex < len(c.attackerPlayers)
	hasDefenserNextCanSelect := c.nextSelectDefenserIndex < len(c.defenserPlayers)

	if !hasAttackerNextCanSelect && !hasDefenserNextCanSelect {
		// 没新人了, 找在场的
		for _, v := range c.combatQueueWinner {
			if v == nil {
				continue
			}

			if v.IsContinueWinLeave(c.getContinueWinTimes(v.isAttacker)) {
				// 是连胜的
				continue
			}

			if v.isAttacker {
				if attacker == nil {
					attacker = v
				} else if attacker.CurStep() > v.CurStep() {
					attacker = v
				}
			} else {
				if defenser == nil {
					defenser = v
				} else if defenser.CurStep() > v.CurStep() {
					defenser = v
				}
			}
		}

		if attacker != nil && defenser != nil {
			if attacker.CurStep() < defenser.curStep {
				fightQueue = attacker.lastFightQueue
			} else {
				fightQueue = defenser.lastFightQueue
			}
		}

		return
	}

	if hasAttackerNextCanSelect && hasDefenserNextCanSelect {
		// 都还有新人，找空的位置，以及连胜下场的空位

		// 找空位
		for idx, v := range c.combatQueueWinner {
			if v == nil {
				fightQueue = idx + 1
				attacker = c.getNextPlayer(true)
				defenser = c.getNextPlayer(false)
				return
			}
		}

		var player *CombatPlayer

		// 位置都被占掉了，找到进攻方防守方中步数最少的
		for _, v := range c.combatQueueWinner {
			if v == nil {
				continue
			}

			if player == nil {
				player = v
			} else if player.CurStep() > v.CurStep() {
				player = v
			}
		}

		if player == nil {
			// 我去，什么鬼
			err = errors.New("在没有空位，没有离场的情况下，竟然没找到步数最小的在场战斗的")
			return
		}

		fightQueue = player.lastFightQueue

		if player.IsContinueWinLeave(c.getContinueWinTimes(player.isAttacker)) {
			// 最少的是连胜的，重新找新人啊
			attacker = c.getNextPlayer(true)
			defenser = c.getNextPlayer(false)
			return

		}

		if player.isAttacker {
			attacker = player
			defenser = c.getNextPlayer(false)
			return
		} else {
			attacker = c.getNextPlayer(true)
			defenser = player
			return
		}
	}

	// 进攻方有新人
	if hasAttackerNextCanSelect {
		attacker = c.getNextPlayer(true)

		// 找在场的
		for _, v := range c.combatQueueWinner {
			if v == nil {
				continue
			}

			if v.isAttacker {
				continue
			}

			if v.IsContinueWinLeave(c.defenserContinueWinCount) {
				// 是连胜的
				continue
			}

			if defenser == nil {
				defenser = v
			} else if defenser.CurStep() > v.CurStep() {
				defenser = v
			}
		}

		if defenser != nil {
			fightQueue = defenser.lastFightQueue
		}
	} else {
		defenser = c.getNextPlayer(false)

		// 找在场的
		for _, v := range c.combatQueueWinner {
			if v == nil {
				continue
			}

			if !v.isAttacker {
				continue
			}

			if v.IsContinueWinLeave(c.attackerContinueWinCount) {
				// 是连胜的
				continue
			}

			if attacker == nil {
				attacker = v
			} else if attacker.CurStep() > v.CurStep() {
				attacker = v
			}
		}

		if attacker != nil {
			fightQueue = attacker.lastFightQueue
		}
	}

	return
}

func (c *MultiCombat) calcScore() (score int32) {
	// 计算评分
	winPlayers := c.getWinPlayers()

	aliveSoldier := 0
	totalSoldier := 0
	for _, p := range winPlayers {
		for _, v := range p.troops {
			aliveSoldier += v.soldier
			totalSoldier += int(v.proto.Captain.Soldier)
		}
	}

	// 积分
	if totalSoldier <= 0 {
		return
	}

	percent := uint64(aliveSoldier * 100 / totalSoldier)

	return c.getScore(percent)
}

// 获得胜利方
func (c *MultiCombat) getWinPlayers() []*CombatPlayer {
	if c.isAttackerWin {
		return c.attackerPlayers
	} else {
		return c.defenserPlayers
	}
}

func (c *MultiCombat) checkIsEnd() (isEnd bool) {
	isAttackerAllDeadOrLeave := int(c.attackerDeadCount+c.attackerLeaveCount) >= len(c.attackerPlayers)
	isDefenserAllDeadOrLeave := int(c.defenserDeadCount+c.defenserLeaveCount) >= len(c.defenserPlayers)

	if isAttackerAllDeadOrLeave {
		if isDefenserAllDeadOrLeave {
			// 双方都死光了或者离场，谁连胜步数高，就算谁赢
			c.isAttackerWin = c.attackerWinMaxStep > c.defenserWinMaxStep
			return true
		} else {
			// 防守方有没死光
			c.isAttackerWin = false
			return true
		}
	} else {
		if isDefenserAllDeadOrLeave {
			// 进攻方有没死光
			c.isAttackerWin = true
			return true
		} else {
			// 都没死光
			return false
		}
	}
}

func (c *MultiCombat) getWinTimesList() (winTimesArray []*shared_proto.BytesInt32Pair) {
	winTimesArray = make([]*shared_proto.BytesInt32Pair, 0, imath.Min(len(c.attackerPlayers), len(c.defenserPlayers)))

	c.walkPlayers(func(p *CombatPlayer) {
		if p.WinTimes() > 0 {
			winTimesArray = append(winTimesArray, &shared_proto.BytesInt32Pair{Key: p.IdBytes(), Value: int32(p.WinTimes())})
		}
	})

	return
}

func (c *MultiCombat) getWinTimesMap() (winTimesMap map[int64]int64) {
	winTimesMap = make(map[int64]int64)

	c.walkPlayers(func(p *CombatPlayer) {
		if p.WinTimes() > 0 {
			winTimesMap[p.id] = int64(p.WinTimes())
		}
	})

	return
}

func (c *MultiCombat) getAliveSoldiers() (aliveSoldiers []*server_proto.AliveSoldierProto) {
	aliveSoldiers = make([]*server_proto.AliveSoldierProto, 0, imath.Min(len(c.attackerPlayers), len(c.defenserPlayers)))

	c.walkPlayers(func(p *CombatPlayer) {
		if !p.AllDead() {
			aliveSoldiers = append(aliveSoldiers, p.EncodeAliveSoldiers())
		}
	})

	return
}

func (c *MultiCombat) walkPlayers(f func(player *CombatPlayer)) {
	for _, p := range c.attackerPlayers {
		f(p)
	}
	for _, p := range c.defenserPlayers {
		f(p)
	}
}
