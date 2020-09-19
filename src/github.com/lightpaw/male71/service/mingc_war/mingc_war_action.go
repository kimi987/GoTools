package mingc_war

import (
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/gen/pb/misc"
	"time"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/mingcdata"
)

type Action struct {
	state     shared_proto.MingcWarTroopState
	startTime time.Time
	endTime   time.Time
	pos       cb.Cube
}

func (c *Action) getStartTime() time.Time                                                   { return c.startTime }
func (c *Action) getEndTime() time.Time                                                     { return c.endTime }
func (c *Action) getState() shared_proto.MingcWarTroopState                                 { return c.state }
func (c *Action) getPos() cb.Cube                                                           { return c.pos }
func (c *Action) onAction(scene *McWarScene, t *McWarTroop, ctime time.Time) (updated bool) { return }

type Actionable interface {
	getStartTime() time.Time
	getEndTime() time.Time
	getState() shared_proto.MingcWarTroopState
	getPos() cb.Cube
	onAction(scene *McWarScene, t *McWarTroop, ctime time.Time) (updated bool)
}

// 移动行为
type MoveAction struct {
	*Action

	startPos cb.Cube
	destPos  cb.Cube
}

func newMoveAction(startTime, endTime time.Time, startPos, destPos cb.Cube) *MoveAction {
	c := &MoveAction{Action: &Action{state: shared_proto.MingcWarTroopState_MC_TP_MOVING}}
	c.startTime = startTime
	c.pos = startPos
	c.startPos, c.destPos = startPos, destPos
	c.endTime = endTime
	return c
}

func (c *MoveAction) onAction(scene *McWarScene, t *McWarTroop, ctime time.Time) (updated bool) {
	if ctime.Before(c.endTime.Add(-mingcdata.McWarLoopDuration)) {
		return
	}

	// 到达
	updated = true
	destBuilding := scene.buildings[c.destPos]
	if destBuilding.arriveOrFight(t, scene, ctime) {
		// 打赢了
		t.action = newStationAction(ctime, c.destPos)
		x, y := destBuilding.pos.XYI32()
		delete(scene.buildings[c.startPos].troops, t.heroId)
		scene.buildings[c.destPos].troops[t.heroId] = t

		msg := mingc_war.NewS2cSceneMoveStationMsg(idbytes.ToBytes(t.heroId), x, y)
		scene.broadcast(msg)
	} else {
		t.failToRelive(scene, destBuilding, ctime, 0)
	}

	return
}

func (scene *McWarScene) getReliveBuildingPos(t *McWarTroop) cb.Cube {
	if t.atk {
		return scene.atkReliveBuilding.pos
	} else {
		return scene.defReliveBuilding.pos
	}
}

// 补兵行为
type ReliveAction struct {
	*Action
}

func newReliveAction(startTime, endTime time.Time, pos cb.Cube) *ReliveAction {
	c := &ReliveAction{Action: &Action{state: shared_proto.MingcWarTroopState_MC_TP_RELIVE}}
	c.startTime, c.endTime = startTime, endTime
	c.pos = pos
	return c
}

func (c *ReliveAction) onAction(scene *McWarScene, t *McWarTroop, ctime time.Time) (updated bool) {
	if ctime.Before(c.endTime) {
		return
	}

	for _, pc := range t.captains {
		if pc == nil {
			continue
		}
		pc.Soldier = pc.TotalSoldier
		pc.FightAmount = pc.FullFightAmount
	}

	updated = true
	t.action = newStationAction(ctime, c.pos)

	troopMsg := mingc_war.NewS2cSceneTroopUpdateMsg(idbytes.ToBytes(t.heroId), t.encode()).Static()
	scene.broadcast(troopMsg)

	x, y := c.pos.XYI32()
	stationMsg := mingc_war.NewS2cSceneMoveStationMsg(idbytes.ToBytes(t.heroId), x, y)
	scene.broadcast(stationMsg)

	if d := scene.dep.Datas().TextHelp().McWarReliveSucc; d != nil {
		scene.dep.World().Send(t.heroId, misc.NewS2cScreenShowWordsMsg(d.Text.New().JsonString()))
	}

	return
}

// 驻扎行为
type StationAction struct {
	*Action
}

func newStationAction(startTime time.Time, pos cb.Cube) *StationAction {
	c := &StationAction{Action: &Action{state: shared_proto.MingcWarTroopState_MC_TP_STATION}}
	c.startTime = startTime
	c.pos = pos
	return c
}

// 入场准备行为
type JoinAction struct {
	*Action
}

func newJoinAction(startTime time.Time, endTime time.Time, pos cb.Cube) *JoinAction {
	c := &JoinAction{Action: &Action{state: shared_proto.MingcWarTroopState_MC_TP_WAIT}}
	c.startTime = startTime
	c.endTime = endTime
	c.pos = pos
	return c
}

func (c *JoinAction) onAction(scene *McWarScene, t *McWarTroop, ctime time.Time) (updated bool) {
	if ctime.Before(c.endTime) {
		return
	}

	updated = true
	t.action = newStationAction(ctime, c.pos)

	x, y := c.pos.XYI32()
	msg := mingc_war.NewS2cSceneMoveStationMsg(idbytes.ToBytes(t.heroId), x, y)
	scene.broadcast(msg)

	return
}
