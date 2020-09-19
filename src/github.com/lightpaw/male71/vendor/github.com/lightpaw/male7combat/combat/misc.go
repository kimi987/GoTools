package combat

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"math"
	"math/rand"
)

func newMisc(p miscInterface) *misc {
	m := &misc{}

	m.maxRound = int(p.GetMaxRound())

	m.raceProtoMap = make(map[shared_proto.Race]*shared_proto.RaceDataProto, len(p.GetRaces()))
	m.priorityMap = make(map[int]int)
	m.raceCoefMap = make(map[int]float64)
	for _, race := range p.GetRaces() {
		m.raceProtoMap[race.Race] = race

		// 左边的最大
		for i, priorityRace := range race.Priority {
			m.priorityMap[priorityKey(int(race.Race), int(priorityRace))] = 1 + len(race.Priority) - i
		}

		for i, raceCoef := range race.RaceCoef {
			targetRace := i + 1
			m.raceCoefMap[priorityKey(int(race.Race), targetRace)] = float64(raceCoef) / Denominator
		}
	}

	m.random = rand.New(rand.NewSource(p.GetSeed()))
	m.coef = math.Max(0, float64(p.GetCoef())/Denominator)
	m.critRate = math.Max(0, float64(p.GetCritRate())/Denominator)
	m.restraintRate = math.Max(0, float64(p.GetRestraintRate())/Denominator)

	m.mapRes = p.GetMapRes()
	m.mapXLen = p.GetMapXLen()
	m.mapYLen = p.GetMapYLen()

	m.minWallAttackRound = p.GetMinWallAttackRound()
	m.maxWallAttachFixDamageRound = p.GetMaxWallAttachFixDamageRound()
	m.maxWallBeenHurtPercent = math.Max(0, float64(p.GetMaxWallBeenHurtPercent())/Denominator)

	m.scorePercent = p.GetScorePercent()

	return m
}

func newMultiMisc(p interface {
	miscInterface
	GetAttackerContinueWinCount() int32
	GetDefenserContinueWinCount() int32
}) *misc {
	m := newMisc(p)

	m.attackerContinueWinCount = int(p.GetAttackerContinueWinCount())
	m.defenserContinueWinCount = int(p.GetDefenserContinueWinCount())

	return m
}

type miscInterface interface {
	GetMaxRound() int32
	GetRaces() []*shared_proto.RaceDataProto
	GetSeed() int64
	GetCritRate() int32
	GetRestraintRate() int32
	GetCoef() int32
	GetMapRes() string
	GetMapXLen() int32
	GetMapYLen() int32
	GetMinWallAttackRound() int32
	GetMaxWallAttachFixDamageRound() int32
	GetMaxWallBeenHurtPercent() int32
	GetScorePercent() []uint64
}

type misc struct {
	raceProtoMap map[shared_proto.Race]*shared_proto.RaceDataProto
	priorityMap  map[int]int
	raceCoefMap  map[int]float64

	maxRound int

	minWallAttackRound          int32
	maxWallAttachFixDamageRound int32
	maxWallBeenHurtPercent      float64

	// 战斗系数
	random *rand.Rand

	coef          float64
	critRate      float64
	restraintRate float64 // 克制技系数

	mapRes  string
	mapXLen int32
	mapYLen int32

	// 进攻方的连胜数量
	attackerContinueWinCount int
	// 防守方的连胜数量
	defenserContinueWinCount int

	scorePercent []uint64
}

func (c *misc) getTroopsCoef(attacker, beenAttacker *Troops) (coef float64) {
	coef, ok := c.raceCoefMap[priorityKey(attacker.race, beenAttacker.race)]
	if ok {
		return coef
	}

	return 1 // 没找到，系数为1
}

func (c *misc) getTargetPriority(sourceRace, targetRace int) int {
	return c.priorityMap[priorityKey(sourceRace, targetRace)]
}

func priorityKey(sourceRace, targetRace int) int {
	return (sourceRace << 16) | targetRace
}

func (c *misc) getContinueWinTimes(isAttacker bool) int {
	if isAttacker {
		return c.attackerContinueWinCount
	} else {
		return c.defenserContinueWinCount
	}
}

func (c *misc) getScore(percent uint64) (score int32) {
	for i := 0; i < len(c.scorePercent); i++ {
		if percent > c.scorePercent[i] {
			score = int32(i + 1)
		} else {
			break
		}
	}

	return
}
