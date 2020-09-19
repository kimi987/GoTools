package combat

import (
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/imath"
	"github.com/pkg/errors"
	"math"
)

func newTroopsArray(genIndex func() int32, isAttacker bool, tps []*shared_proto.CombatTroopsProto, misc *misc) ([]*Troops, error) {
	out := make([]*Troops, 0, len(tps))
	for i, v := range tps {
		troops, err := newTroops(genIndex(), isAttacker, v, misc)
		if err != nil {
			return nil, errors.Wrapf(err, "Troops[%v]", i)
		}

		out = append(out, troops)
	}

	return out, nil
}

func newTroops(index int32, isAttacker bool, p *shared_proto.CombatTroopsProto, misc *misc) (*Troops, error) {
	if index <= 0 {
		return nil, errors.Errorf("Troops.index should > 0, v: %v", index)
	}

	c := p.Captain
	if c == nil {
		return nil, errors.Errorf("Troops.Captain == nil")
	}

	if c.Soldier <= 0 {
		return nil, errors.Errorf("Troops.Soldier should > 0,  v: %v", c.Soldier)
	}

	if c.TotalSoldier < c.Soldier {
		return nil, errors.Errorf("Troops.TotalSoldier should > Troops.Soldier,  v: %v", c.TotalSoldier)
	}

	if c.LifePerSoldier <= 0 {
		return nil, errors.Errorf("Troops.LifePerSoldier should > 0,  v: %v", c.LifePerSoldier)
	}

	raceProto := misc.raceProtoMap[c.Race]
	if raceProto == nil {
		return nil, errors.Errorf("Troops.race data not found, v: %v", c.Race)
	}

	posX := p.X
	posY := p.Y
	if p.FightIndex > 0 {
		// 根据index设置位置
		if isAttacker {
			posX = 0
		} else {
			posX = misc.mapXLen - 1
		}

		posY = p.FightIndex - 1
	} else {
		// 自己设置位置
		if !(0 <= posX && posX < misc.mapXLen) {
			return nil, errors.Errorf("Troops.X not in [0 - %v), v: %v", misc.mapXLen, posX)
		}

		if !(0 <= posY && posY < misc.mapYLen) {
			return nil, errors.Errorf("Troops.Y not in [0 - %v), v: %v", misc.mapYLen, posY)
		}
	}

	//p.Index = index
	pos := &shared_proto.CombatTroopsPosProto{
		Index: index,
		X:     posX,
		Y:     posY,
	}

	out := &Troops{
		index:      index,
		isAttacker: isAttacker,
		proto:      p,
		pos:        pos,
	}

	out.race = int(c.Race)
	out.raceProto = raceProto
	out.attackRange = int(raceProto.AttackRange)
	out.moveTimesPerRound = int(raceProto.MoveTimesPerRound)
	out.moveSpeed = int(raceProto.MoveSpeed)

	out.lifePerSoldier = int(c.LifePerSoldier)

	// 攻防体敏
	out.attack = float64(c.TotalStat.Attack)
	out.defense = float64(c.TotalStat.Defense)
	out.strength = float64(c.TotalStat.Strength)
	out.dexterity = float64(c.TotalStat.Dexterity)
	out.Scoef = math.Max(1, float64(c.Morale)/Denominator)
	out.damageIncrePer = float64(c.TotalStat.DamageIncrePer) / Denominator
	out.damageDecrePer = float64(c.TotalStat.DamageDecrePer) / Denominator

	out.wallCoef = math.Max(float64(raceProto.WallCoef)/Denominator, 0.1)

	out.canTriggerRestraintSpell = c.GetCanTriggerRestraintSpell()

	out.setLife(int(c.LifePerSoldier) * int(c.Soldier))

	return out, nil
}

// 部队
type Troops struct {
	index int32

	isAttacker bool

	proto *shared_proto.CombatTroopsProto
	pos   *shared_proto.CombatTroopsPosProto

	race      int
	raceProto *shared_proto.RaceDataProto

	attackRange int

	moveSpeed int

	moveTimesPerRound int

	lifePerSoldier int

	// 属性 float64 攻防体敏
	attack         float64
	defense        float64
	strength       float64
	dexterity      float64
	damageIncrePer float64 // 伤害增加百分比
	damageDecrePer float64 // 伤害减少百分比

	wallCoef float64 // 城墙克制系数

	Scoef float64 // 士气系数

	// 下面是战斗过程计算使用的变量

	moveTimes int // 移动力

	x, y int

	life int

	soldier int

	killSoldier int

	// 能否触发克制技
	canTriggerRestraintSpell bool
}

func (c *Troops) String() string {
	return fmt.Sprintf("Troops{index: %d, isAttacker: %v, lifePerSoldier:%d, TotalSoldier: %d, soldier: %d, x: %d, y: %d}",
		c.index, c.isAttacker, c.lifePerSoldier, c.proto.Captain.TotalSoldier, c.soldier, c.x, c.y)
}

func (c *Troops) resetPosition() {
	c.x = int(c.pos.X)
	c.y = int(c.pos.Y)
}

func (c *Troops) reduceLife(toReduce int) (reduceSoldier int) {
	oldSoldier := c.soldier
	c.setLife(c.life - toReduce)
	return oldSoldier - c.soldier
}

func (c *Troops) setLife(toSet int) int {
	c.life = imath.Max(toSet, 0)
	c.soldier = (c.life + c.lifePerSoldier - 1) / c.lifePerSoldier
	return c.life
}

func (c *Troops) reduceSoldier(toReduce int) (reduceSoldier int) {
	return c.reduceLife(toReduce * c.lifePerSoldier)
}

// 士兵
func (c *Troops) Soldier() int {
	return c.soldier
}

func (c *Troops) resetMoveTimes() {
	c.moveTimes = c.moveTimesPerRound
}

func (c *Troops) isAlive() bool {
	return c.life > 0
}

func (c *Troops) isMovable() bool {
	return c.moveTimes > 0
}

// 是否克制的职业
func (c *Troops) isRestraintRace(enemyRace shared_proto.Race) (isRestraintRace bool) {
	for _, race := range c.raceProto.RestraintRace {
		if race == enemyRace {
			return true
		}
	}
	return false
}

// 是否克制技的触发轮
func (c *Troops) isRestraintSpellTriggerRound(round int32) (isTriggerRound bool) {
	switch c.raceProto.RestraintRoundType {
	case shared_proto.RestraintRoundType_ODD:
		//isTriggerRound = round%2 == 1
		// 近战首次发动回合为第3回合，之后每个3个回合发动1次
		if round >= 3 {
			isTriggerRound = round%3 == 0
		}
	case shared_proto.RestraintRoundType_EVEN:
		//isTriggerRound = round%2 == 0
		// 远程首次发动回合为第2回合，之后每个3个回合发动1次
		if round >= 2 {
			isTriggerRound = (round-2)%3 == 0
		}
	default:
		logrus.Debugln("未处理的克制技能的触发轮次类型!%v", c.raceProto.RestraintRoundType)
	}
	return
}

func (c *Troops) getRange(target *Troops) int {
	return imath.Abs(c.x-target.x) + imath.Abs(c.y-target.y)
}

type TroopsSpeedSlice []*Troops

func (p TroopsSpeedSlice) Len() int      { return len(p) }
func (p TroopsSpeedSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p TroopsSpeedSlice) Less(i, j int) bool {
	t1 := p[i]
	t2 := p[j]

	// 机动力从高到低排序
	if t1.moveSpeed != t2.moveSpeed {
		return t1.moveSpeed > t2.moveSpeed
	}

	// 机动力相同时，攻方武将按照编号从小到大排在前面，守方武将按照编号从小到大排在后面
	if t1.isAttacker != t2.isAttacker {
		return t1.isAttacker
	}

	return t1.index < t2.index
}
