package mingc_war

import (
	"time"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/pb/shared_proto"
	"math/rand"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/util/timeutil"
)

type McWarAiTroop struct {
	ai             Ai
	scene          *McWarScene
	data           *mingcdata.MingcWarSceneData
	troop          *McWarTroop
	speed          float64
	nextActionTime time.Time
}

var (
	idGener = atomic.NewUint64(1)
)

func nextActionTime(ctime time.Time) time.Time {
	return ctime.Add(time.Duration((5 + rand.Int63n(10)) * int64(time.Second))).Add(time.Duration(rand.Int63n(500) * int64(time.Millisecond)))
}

func firstActionTime(baseDur time.Duration, ctime time.Time) time.Time {
	return ctime.Add(baseDur).Add(time.Duration(rand.Int63n(30) * int64(time.Second)))
}

func (c *McWarAiTroop) unmershal(p *server_proto.McWarNpcTroop, troop *McWarTroop, scene *McWarScene, dep iface.ServiceDep) {
	c.scene = scene
	c.data = dep.Datas().GetMingcWarSceneData(scene.mcId)
	c.speed = dep.Datas().MingcMiscData().Speed
	c.troop = troop
	c.nextActionTime = firstActionTime(dep.Datas().MingcMiscData().FightPrepareDuration, dep.Time().CurrentTime())

	switch p.AiType {
	case server_proto.McWarAiType_MC_AI_ATK:
		c.ai = &AtkAi{}
	case server_proto.McWarAiType_MC_AI_DEF_HOME:
		c.ai = &DefHomeAi{}
	case server_proto.McWarAiType_MC_AI_DEF_CASTLE:
		c.ai = &DefCastleAi{}
	case server_proto.McWarAiType_MC_AI_DEF_GATE:
		c.ai = &DefGateAi{}
	default:
		c.ai = &StationAi{}
	}
}

func (c *McWarAiTroop) encodeServer() (p *server_proto.McWarNpcTroop) {
	p = &server_proto.McWarNpcTroop{}
	p.AiType = c.ai.AiType()
	p.TroopId = c.troop.heroId
	return p
}

func newMcWarAiTroop(npc *mingcdata.MingcWarNpcData, scene *McWarScene, dep iface.ServiceDep, ctime time.Time) *McWarAiTroop {
	ait := &McWarAiTroop{}
	ait.scene = scene
	ait.data = dep.Datas().GetMingcWarSceneData(scene.mcId)
	ait.speed = dep.Datas().MingcMiscData().Speed
	ait.nextActionTime = firstActionTime(dep.Datas().MingcMiscData().FightPrepareDuration, ctime)

	var captains []*shared_proto.CaptainInfoProto
	for _, capt := range npc.Npc.Captains {
		c := capt.EncodeCaptainInfo()
		c.FullFightAmount = u64.Int32(capt.FightAmount)
		captains = append(captains, c)
	}
	npcActionTime := ctime.Add(dep.Datas().MingcMiscData().FightPrepareDuration)
	npcId := npcid.GetNpcId(idGener.Inc(), scene.mcId, npcid.NpcType_MingCheng)
	ait.troop = newMcWarNpcTroop(npcId, npc, false, captains, scene.defReliveBuilding.pos, ctime, npcActionTime, dep.Datas())

	switch npc.AiType {
	case server_proto.McWarAiType_MC_AI_ATK:
		ait.ai = &AtkAi{}
	case server_proto.McWarAiType_MC_AI_DEF_HOME:
		ait.ai = &DefHomeAi{}
	case server_proto.McWarAiType_MC_AI_DEF_CASTLE:
		ait.ai = &DefCastleAi{}
	case server_proto.McWarAiType_MC_AI_DEF_GATE:
		ait.ai = &DefGateAi{}
	default:
		ait.ai = &StationAi{}
	}

	return ait
}

type Ai interface {
	onAction(t *McWarAiTroop, ctime time.Time) (updated bool)
	AiType() (t server_proto.McWarAiType)
}

type AtkAi struct{}

func (ai *AtkAi) AiType() (t server_proto.McWarAiType) { return server_proto.McWarAiType_MC_AI_ATK }
func (ai *AtkAi) onAction(t *McWarAiTroop, ctime time.Time) (updated bool) {
	return aiMove(0, t, ctime)
}

func (t *McWarAiTroop) station(pos cb.Cube, ctime time.Time) {
	t.troop.action = newStationAction(ctime, pos)
	x, y := pos.XYI32()
	t.scene.updateViewMsg()
	t.scene.broadcast(mingc_war.NewS2cSceneMoveStationMsg(idbytes.ToBytes(t.troop.heroId), x, y))
}

func (t *McWarAiTroop) move(start, dest cb.Cube, ctime time.Time) {
	endTime := calcEndTime(start, dest, ctime, t.speed)
	t.troop.action = newMoveAction(ctime, endTime, start, dest)
	t.nextActionTime = nextActionTime(endTime)
	sx, sy := start.XYI32()
	dx, dy := dest.XYI32()
	t.scene.updateViewMsg()
	t.scene.broadcast(mingc_war.NewS2cSceneOtherMoveMsg(idbytes.ToBytes(t.troop.heroId), sx, sy, dx, dy, timeutil.Marshal32(endTime)))
}

type DefHomeAi struct{}

func (ai *DefHomeAi) AiType() (t server_proto.McWarAiType) { return server_proto.McWarAiType_MC_AI_DEF_HOME }
func (ai *DefHomeAi) onAction(t *McWarAiTroop, ctime time.Time) (updated bool) {
	return aiMove(shared_proto.MingcWarBuildingType_MC_B_HOME, t, ctime)
}

type DefCastleAi struct{}

func (ai *DefCastleAi) AiType() (t server_proto.McWarAiType) { return server_proto.McWarAiType_MC_AI_DEF_CASTLE }
func (ai *DefCastleAi) onAction(t *McWarAiTroop, ctime time.Time) (updated bool) {
	return aiMove(shared_proto.MingcWarBuildingType_MC_B_CASTLE, t, ctime)
}

type DefGateAi struct{}

func (ai *DefGateAi) AiType() (t server_proto.McWarAiType) { return server_proto.McWarAiType_MC_AI_DEF_GATE }
func (ai *DefGateAi) onAction(t *McWarAiTroop, ctime time.Time) (updated bool) {
	return aiMove(shared_proto.MingcWarBuildingType_MC_B_GATE, t, ctime)
}

func aiMove(bType shared_proto.MingcWarBuildingType, t *McWarAiTroop, ctime time.Time) (updated bool) {
	if t.troop.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION {
		return
	}

	if ctime.Before(t.nextActionTime) {
		return
	}

	t.nextActionTime = nextActionTime(ctime)

	pos := t.troop.action.getPos()
	currBuilding := t.scene.buildings[pos]
	if currBuilding != nil && currBuilding.atk && currBuilding.prosperity > 0 {
		return
	}

	var isDefAi bool
	if t.ai.AiType() != server_proto.McWarAiType_MC_AI_ATK {
		isDefAi = true
	}

	if isDefAi {
		if currBuilding.data.Type == bType {
			// 是否继续驻扎
			var key int
			if bType == shared_proto.MingcWarBuildingType_MC_B_GATE {
				key = 7
			} else {
				key = 5
			}
			if rand.Intn(10) >= key {
				return
			}
		}
	}

	// 走
	var targetDests, otherDests []cb.Cube
	for _, destPos := range t.data.McWarMap()[pos] {
		var moveType aiDestType

		if isDefAi {
			moveType = t.defAiMoveTarget(pos, destPos, bType)
		} else {
			moveType = t.atkAiMoveTarget(pos, destPos)
		}

		switch moveType {
		case aiDestTarget:
			targetDests = append(targetDests, destPos)
		case aiDestOther:
			otherDests = append(otherDests, destPos)
		}
		continue
	}

	var dest cb.Cube
	if len(targetDests) > 0 {
		dest = targetDests[rand.Intn(len(targetDests))]
	} else if len(otherDests) > 0 {
		dest = otherDests[rand.Intn(len(otherDests))]
	} else {
		return
	}

	t.move(pos, dest, ctime)
	return true
}

type aiDestType int

var (
	aiDestNoWay  = aiDestType(0)
	aiDestTarget = aiDestType(1)
	aiDestOther  = aiDestType(2)
)

func (t *McWarAiTroop) atkAiMoveTarget(pos, destPos cb.Cube) (destType aiDestType) {
	destType = aiDestNoWay

	// 攻方始终往左走
	if posX(destPos) >= posX(pos) {
		return
	}

	b, ok := t.scene.buildings[destPos]
	if !ok {
		return
	}

	if b.data.Type == shared_proto.MingcWarBuildingType_MC_B_RELIVE {
		return
	}

	if b.data.Type == shared_proto.MingcWarBuildingType_MC_B_TOU_SHI {
		return
	}

	if b.atk {
		return aiDestTarget
	} else {
		return aiDestOther
	}

	return
}

func (t *McWarAiTroop) defAiMoveTarget(pos, destPos cb.Cube, targetBuildingType shared_proto.MingcWarBuildingType) (destType aiDestType) {
	destType = aiDestNoWay

	b, ok := t.scene.buildings[destPos]
	if !ok {
		return
	}

	if b.data.Type == shared_proto.MingcWarBuildingType_MC_B_RELIVE {
		return
	}

	if b.data.Type == shared_proto.MingcWarBuildingType_MC_B_TOU_SHI {
		return
	}

	if b.atk {
		return
	}

	if b.data.Type == targetBuildingType {
		return aiDestTarget
	} else {
		return aiDestOther
	}

	return
}

type StationAi struct{}

func (ai *StationAi) AiType() (t server_proto.McWarAiType) { return server_proto.McWarAiType_MC_AI_STATION }
func (ai *StationAi) onAction(t *McWarAiTroop, ctime time.Time) (updated bool) {
	return
}

func calcEndTime(start, dest cb.Cube, ctime time.Time, speed float64) time.Time {
	x1, y1 := start.XY()
	x2, y2 := dest.XY()
	distance := util.Distance(x1, y1, x2, y2)

	seconds := uint64(float64(distance) / speed)
	return ctime.Add(time.Duration(seconds * uint64(time.Second)))
}

func posX(c cb.Cube) int {
	x, _ := c.XY()
	return x
}
