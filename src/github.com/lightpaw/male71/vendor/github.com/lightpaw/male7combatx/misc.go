package combatx

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"math/rand"
)

func newMisc(p *server_proto.CombatXRequestServerProto) *misc {

	m := &misc{}
	m.request = p

	m.random = rand.New(rand.NewSource(p.GetSeed()))
	m.mapRes = p.GetMapRes()
	m.mapXLen = int(p.GetMapXLen())
	m.mapYLen = int(p.GetMapYLen())

	m.isDebug = p.Debug

	return m
}

type misc struct {
	request *server_proto.CombatXRequestServerProto

	// 战斗系数
	random *rand.Rand

	mapRes  string
	mapXLen int
	mapYLen int

	isDebug bool
}

//func (c *misc) getTargetPriority(sourceRace, targetRace shared_proto.Race) int {
//	return c.priorityMap[priorityKey(sourceRace, targetRace)]
//}
//
//func priorityKey(sourceRace, targetRace shared_proto.Race) int {
//	return (int(sourceRace) << 16) | int(targetRace)
//}
//
//func (c *misc) getScore(percent uint64) (score int32) {
//	for i := 0; i < len(c.scorePercent); i++ {
//		if percent > c.scorePercent[i] {
//			score = int32(i + 1)
//		} else {
//			break
//		}
//	}
//
//	return
//}
//
//func (c *misc) getTroopsCoef(attacker, beenAttacker *Troops) (coef float64) {
//	coef, ok := c.raceCoefMap[priorityKey(attacker.getRace(), beenAttacker.getRace())]
//	if ok {
//		return coef
//	}
//
//	return 1 // 没找到，系数为1
//}
