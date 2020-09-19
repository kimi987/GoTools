package mingc_war

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/chat"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/logp"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"math"
	"time"
)

const (
	mc_chat_record_max_num = 100 // 名城聊天的最大缓存
	mc_rank_max_num        = 100 // 缓存的最大排名数量
)

var sameVersionRankMsg = mingc_war.NewS2cApplyRefreshRankMsg(0, &shared_proto.McWarTroopsRankProto{}).Static()

type McWarScene struct {
	dep           iface.ServiceDep
	mcId          uint64
	mcAllYinliang uint64

	startTime time.Time

	atkReliveBuilding *McWarBuilding
	defReliveBuilding *McWarBuilding
	atkHomeBuilding   *McWarBuilding
	defHomeBuilding   *McWarBuilding

	buildings map[cb.Cube]*McWarBuilding

	watchHeros map[int64]struct{}
	atkTroops  map[int64]*McWarTroop
	defTroops  map[int64]*McWarTroop

	atkDurmTimes uint64
	defDurmTimes uint64
	durmStopTime time.Time
	atkDurmStat  *shared_proto.SpriteStatProto
	defDurmStat  *shared_proto.SpriteStatProto

	troopsRank *McWarTroopsRank // 战绩排行

	ended        bool
	atkWin       bool
	endedNoticed bool // 已经通知上层

	viewMsg pbutil.Buffer

	record *McWarFightRecord

	npcs []*McWarAiTroop

	fightState shared_proto.MingcWarFightState

	// 攻方聊天记录
	atkChatRecord *chat.ChatRecord
	// 守方聊天记录
	defChatRecord *chat.ChatRecord
}

func newDefaultMcWarScene(datas iface.ConfigDatas) *McWarScene {
	c := &McWarScene{}

	c.atkTroops = make(map[int64]*McWarTroop)
	c.defTroops = make(map[int64]*McWarTroop)
	c.buildings = make(map[cb.Cube]*McWarBuilding)
	c.record = &McWarFightRecord{guilds: make(map[int64]*McWarFightGuildRecord)}
	c.fightState = shared_proto.MingcWarFightState_MC_F_FIGHT_PREPARE

	c.atkDurmStat = data.EmptyStatProto()
	c.defDurmStat = data.EmptyStatProto()

	firstSendNum := datas.MiscGenConfig().FirstHistoryChatSend
	c.atkChatRecord = chat.NewChatRecord(mc_chat_record_max_num, firstSendNum)
	c.defChatRecord = chat.NewChatRecord(mc_chat_record_max_num, firstSendNum)

	return c
}

func newMcWarScene(mc *MingcObj, dep iface.ServiceDep, startTime time.Time) *McWarScene {
	c := newDefaultMcWarScene(dep.Datas())
	c.dep = dep
	c.mcId = mc.id
	mingc := dep.Mingc().Mingc(c.mcId)
	c.mcAllYinliang = mingc.Yinliang() + mingc.ExtraYinliang()
	c.startTime = startTime
	c.durmStopTime = c.startTime.Add(c.dep.Datas().MingcMiscData().FightPrepareDuration).Add(-c.dep.Datas().MingcMiscData().DrumStopDuration)

	c.initBuilding(dep.Datas().GetMingcWarSceneData(mc.id), dep.Datas().GetMingcWarBuildingData)

	c.record.guilds[mc.atkId] = NewMcWarFightGuildRecord(mc.atkId)
	c.record.guilds[mc.defId] = NewMcWarFightGuildRecord(mc.defId)
	for _, gid := range mc.astAtkList {
		c.record.guilds[gid] = NewMcWarFightGuildRecord(gid)
	}
	for _, gid := range mc.astDefList {
		c.record.guilds[gid] = NewMcWarFightGuildRecord(gid)
	}

	if mc.defId <= 0 {
		c.npcJoinFight(mc, dep)
	}

	c.updateViewMsg()

	c.troopsRank = newMcWarTroopsRank()
	return c
}

func (c *McWarScene) unmershal(p *server_proto.McWarSceneServerProto, dep iface.ServiceDep) {
	c.dep = dep
	c.mcId = p.McId
	c.ended = p.Ended
	mingc := dep.Mingc().Mingc(c.mcId)
	c.mcAllYinliang = mingc.Yinliang() + mingc.ExtraYinliang()
	c.startTime = timeutil.Unix64(p.StartTime)

	ctime := c.dep.Time().CurrentTime()
	fightRunStartTime := c.startTime.Add(c.dep.Datas().MingcMiscData().FightPrepareDuration)
	if ctime.After(fightRunStartTime) {
		c.fightState = shared_proto.MingcWarFightState_MC_F_FIGHT_RUNNING
	}

	c.buildings = make(map[cb.Cube]*McWarBuilding)
	for _, bp := range p.Building {
		b := &McWarBuilding{}
		b.unmarshal(bp, c.mcId, dep.Datas())
		c.buildings[b.pos] = b
		switch b.data.Type {
		case shared_proto.MingcWarBuildingType_MC_B_RELIVE:
			if b.atk {
				c.atkReliveBuilding = b
			} else {
				c.defReliveBuilding = b
			}
		case shared_proto.MingcWarBuildingType_MC_B_HOME:
			if b.atk {
				c.atkHomeBuilding = b
			} else {
				c.defHomeBuilding = b
			}
		}

		for _, t := range b.troops {
			if t.atk {
				c.atkTroops[t.heroId] = t
			} else {
				c.defTroops[t.heroId] = t
			}
		}
	}

	c.durmStopTime = fightRunStartTime.Add(-c.dep.Datas().MingcMiscData().DrumStopDuration)
	c.atkDurmTimes = p.AtkDrumTimes
	c.defDurmTimes = p.DefDrumTimes
	c.atkDurmStat = p.AtkDrumStat
	c.defDurmStat = p.DefDrumStat

	c.record.unmarshal(p.Record)

	for _, npcp := range p.Npc {
		if t := c.defTroops[npcp.TroopId]; t != nil {
			npc := &McWarAiTroop{}
			npc.unmershal(npcp, t, c, dep)
			c.npcs = append(c.npcs, npc)
		}
	}

	// todo 看需不需要保存数据
	c.troopsRank = newMcWarTroopsRank()

	c.updateViewMsg()
}

func (c *McWarScene) encodeServer() *server_proto.McWarSceneServerProto {
	p := &server_proto.McWarSceneServerProto{}
	p.McId = c.mcId
	p.Ended = c.ended
	p.AtkWin = c.atkWin
	p.StartTime = timeutil.Marshal64(c.startTime)
	p.EndedNoticed = c.endedNoticed
	for _, b := range c.buildings {
		p.Building = append(p.Building, b.encodeServer())
	}

	p.AtkDrumTimes = c.atkDurmTimes
	p.DefDrumTimes = c.defDurmTimes
	p.AtkDrumStat = c.atkDurmStat
	p.DefDrumStat = c.defDurmStat

	p.Record = c.record.encodeServer()
	for _, npc := range c.npcs {
		p.Npc = append(p.Npc, npc.encodeServer())
	}

	return p
}

func (c *McWarScene) encode() *shared_proto.McWarSceneProto {
	p := &shared_proto.McWarSceneProto{}
	p.McId = u64.Int32(c.mcId)
	p.McAllYinliang = u64.Int32(c.mcAllYinliang)
	p.StartTime = timeutil.Marshal32(c.startTime)
	p.FightState = c.fightState

	for _, b := range c.buildings {
		p.Building = append(p.Building, b.encode())
	}
	for _, t := range c.atkTroops {
		p.Troop = append(p.Troop, t.encode())
	}
	for _, t := range c.defTroops {
		p.Troop = append(p.Troop, t.encode())
	}

	p.AtkDrumTimes = u64.Int32(c.atkDurmTimes)
	p.DefDrumTimes = u64.Int32(c.defDurmTimes)
	p.AtkDrumStat = c.atkDurmStat
	p.DefDrumStat = c.defDurmStat
	p.DrumStopTime = timeutil.Marshal32(c.durmStopTime)

	return p
}

func (c *McWarScene) npcJoinFight(mc *MingcObj, dep iface.ServiceDep) {
	ctime := dep.Time().CurrentTime()
	for _, npcData := range dep.Datas().GetMingcWarNpcDataArray() {
		if c.mcId != npcData.Mingc.Id {
			continue
		}
		aiTroop := newMcWarAiTroop(npcData, c, dep, ctime)
		c.npcs = append(c.npcs, aiTroop)
		c.defTroops[aiTroop.troop.heroId] = aiTroop.troop
		if npcData.Def {
			mc.defId = npcData.GuildId()
		} else {
			mc.astDefList = append(mc.astDefList, int64(npcData.GuildId()))
		}

		troop := aiTroop.troop
		if g, ok := c.record.guilds[troop.gid]; !ok || g == nil {
			g = NewMcWarFightGuildRecord(troop.gid)
			c.record.guilds[troop.gid] = g
			g.addTroop(troop.heroId)
		} else {
			g.addTroop(troop.heroId)
		}

		joinFightMsg := mingc_war.NewS2cOtherJoinFightMsg(u64.Int32(c.mcId), idbytes.ToBytes(troop.heroId), nil).Static()
		c.broadcastExclude(joinFightMsg, troop.heroId)
		troopUpdateMsg := mingc_war.NewS2cSceneTroopUpdateMsg(idbytes.ToBytes(troop.heroId), troop.encode())
		c.broadcastExclude(troopUpdateMsg, troop.heroId)
	}
}

func (c *McWarScene) onUpdate(ctime time.Time) {
	if c.ended {
		return
	}

	var updated bool
	if c.fightState == shared_proto.MingcWarFightState_MC_F_FIGHT_PREPARE {
		fightRunStartTime := c.startTime.Add(c.dep.Datas().MingcMiscData().FightPrepareDuration)
		if ctime.After(fightRunStartTime) {
			c.fightState = shared_proto.MingcWarFightState_MC_F_FIGHT_RUNNING
			c.addDrumStat()
			c.broadcast(mingc_war.NewS2cSceneFightPrepareEndMsg(timeutil.Marshal32(fightRunStartTime)))
			updated = true
		}
	}

	for _, npc := range c.npcs {
		updated = npc.ai.onAction(npc, ctime) || updated
	}
	for _, t := range c.atkTroops {
		updated = t.onUpdate(c, ctime) || updated
	}
	for _, t := range c.defTroops {
		updated = t.onUpdate(c, ctime) || updated
	}
	for _, b := range c.buildings {
		updated = b.onUpdate(c, ctime) || updated
		if b.data.Type == shared_proto.MingcWarBuildingType_MC_B_HOME && b.prosperity <= 0 {
			c.onEnd()
			break
		}
	}

	if updated {
		c.updateViewMsg()
	}
}

func (c *McWarScene) updateViewMsg() {
	c.viewMsg = mingc_war.NewS2cViewMcWarSceneMsg(c.encode()).Static()
}

func (c *McWarScene) broadcast(msg pbutil.Buffer) {
	c.broadcastExclude(msg, 0)
}

func (c *McWarScene) broadcastExclude(msg pbutil.Buffer, exclude int64) {
	c.broadcastCampExclude(msg, true, false, exclude)
	c.broadcastCampExclude(msg, false, true, exclude)
}

func (c *McWarScene) broadcastCampExclude(msg pbutil.Buffer, atk bool, excludeWatch bool, exclude int64) {
	var troops map[int64]*McWarTroop
	if atk {
		troops = c.atkTroops
	} else {
		troops = c.defTroops
	}

	var ids []int64
	for _, t := range troops {
		if t.heroId != exclude {
			ids = append(ids, t.heroId)
		}
	}

	if !excludeWatch {
		for wid := range c.watchHeros {
			ids = append(ids, wid)
		}
	}

	c.dep.World().MultiSend(ids, msg)
}

func (c *McWarScene) onEnd() {
	c.atkWin = c.defHomeBuilding.prosperity <= 0
	c.fightState = shared_proto.MingcWarFightState_MC_F_INVALID
	c.ended = true
}

func (c *McWarScene) watch(heroId int64) {
	c.watchHeros[heroId] = struct{}{}
}

func (c *McWarScene) quitWatch(heroId int64) {
	delete(c.watchHeros, heroId)
}

func (c *McWarScene) joinFight(t *McWarTroop) {
	if t.atk {
		c.atkTroops[t.heroId] = t
		t.addStat(c.atkDurmStat)
	} else {
		c.defTroops[t.heroId] = t
		t.addStat(c.defDurmStat)
	}

	if g, ok := c.record.guilds[t.gid]; !ok {
		g = NewMcWarFightGuildRecord(t.gid)
		c.record.guilds[t.gid] = g
		g.addTroop(t.heroId)
	} else {
		g.addTroop(t.heroId)
	}

	otherJoinMsg := mingc_war.NewS2cOtherJoinFightMsg(u64.Int32(c.mcId), idbytes.ToBytes(t.heroId), nil).Static()
	c.broadcastExclude(otherJoinMsg, t.heroId)
	joinTroopUpdateMsg := mingc_war.NewS2cSceneTroopUpdateMsg(idbytes.ToBytes(t.heroId), t.encode())
	c.broadcastExclude(joinTroopUpdateMsg, t.heroId)

	c.updateViewMsg()
}

func (c *McWarScene) quitFight(heroId int64) (troop *McWarTroop) {
	if t, ok := c.atkTroops[heroId]; ok {
		delete(c.atkTroops, heroId)
		troop = t
	} else if t, ok := c.defTroops[heroId]; ok {
		delete(c.defTroops, heroId)
		troop = t
	} else {
		logrus.Warnf("名城战退出战斗，找不到队伍.mcid:%v heroid:%v", c.mcId, heroId)
	}
	if b, ok := c.buildings[troop.action.getPos()]; ok {
		delete(b.troops, troop.heroId)
	} else {
		logrus.Warnf("名城战退出战斗，找不到队伍所属的建筑.mcid:%v heroid:%v b:%v", c.mcId, heroId, troop.action.getPos())
	}

	if troop.rankObj != nil {
		troop.rankObj.resetMultiKill()
	}

	msg := mingc_war.NewS2cOtherQuitFightMsg(u64.Int32(c.mcId), idbytes.ToBytes(heroId)).Static()
	c.broadcastExclude(msg, heroId)

	c.updateViewMsg()

	return
}

func (c *McWarScene) pos(heroId int64) (p cb.Cube, ok bool) {
	t, _ := c.troop(heroId)
	if t == nil {
		return
	}
	p = t.action.getPos()
	ok = true
	return
}

func (c *McWarScene) back(heroId int64, ctime time.Time) (endTime time.Time, succ bool) {
	t, _ := c.troop(heroId)
	if t == nil {
		return
	}
	if t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_MOVING {
		return
	}
	if c.buildings[t.action.getPos()] == nil {
		return
	}

	endTime = ctime.Add(ctime.Sub(t.action.getStartTime()))
	if endTime.Before(ctime) {
		endTime = ctime
	}

	dest := t.action.getPos()
	t.action.(*MoveAction).pos = t.action.(*MoveAction).destPos
	t.move(c, ctime, endTime, dest, heroId)
	c.updateViewMsg()
	succ = true
	return
}

func (c *McWarScene) speedUp(heroId int64, speedUpRate float64, ctime time.Time) (endTime time.Time, succ bool) {
	t, _ := c.troop(heroId)
	if t == nil || speedUpRate <= 0 {
		return
	}
	if t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_MOVING {
		return
	}
	dur := t.action.getEndTime().Sub(ctime)
	newDur := i64.Max(0, int64(float64(dur)/(1+speedUpRate)))
	endTime = ctime.Add(time.Duration(newDur))
	t.action.(*MoveAction).endTime = endTime
	c.updateViewMsg()
	succ = true

	x, y := t.action.getPos().XYI32()
	dx, dy := t.action.(*MoveAction).destPos.XYI32()
	msg := mingc_war.NewS2cSceneOtherMoveMsg(idbytes.ToBytes(t.heroId), x, y, dx, dy, timeutil.Marshal32(endTime)).Static()
	c.broadcast(msg)

	return
}

func (c *McWarScene) canArrive(heroId int64, dest cb.Cube) (ok bool) {
	if c.buildings[dest] == nil {
		return
	}
	t, _ := c.troop(heroId)
	if t == nil {
		return
	}

	// 正在摧毁对方据点时，不能再移动到对方没被摧毁的据点
	if stationBuilding := c.buildings[t.action.getPos()]; stationBuilding.atk != t.atk && stationBuilding.prosperity > 0 {
		if c.buildings[dest].atk != t.atk && c.buildings[dest].prosperity > 0 {
			return
		}
	}
	ok = true
	return
}

func (c *McWarScene) move(heroId int64, dest cb.Cube, ctime time.Time) (endTime time.Time, succ bool) {
	if c.buildings[dest] == nil {
		return
	}
	t, _ := c.troop(heroId)
	if t == nil {
		return
	}
	if t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION {
		return
	}

	endTime = moveEndTime(t.action.getPos(), dest, t.getSpeed(), ctime)
	t.move(c, ctime, endTime, dest, heroId)
	c.updateViewMsg()
	succ = true
	return
}

func (c *McWarScene) relive(heroId int64, ctime time.Time) (endTime time.Time, succ bool) {
	t, _ := c.troop(heroId)
	if t == nil {
		return
	}
	if t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION {
		return
	}

	b := c.buildings[t.action.getPos()]
	if b.atk != t.atk || b.data.Type != shared_proto.MingcWarBuildingType_MC_B_RELIVE {
		return
	}
	if t.isFullSolider() {
		return
	}

	endTime = t.relive(c, ctime)
	c.updateViewMsg()
	succ = true
	return
}

func (c *McWarScene) changeMode(heroId int64, newMode shared_proto.MingcWarModeType) (succ bool) {
	t, _ := c.troop(heroId)
	if t == nil {
		return
	}
	if t.mode == newMode {
		return
	}
	if t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION &&
		t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_WAIT &&
		t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_RELIVE {
		return
	}

	b := c.buildings[t.action.getPos()]
	if b.atk != t.atk || b.data.Type != shared_proto.MingcWarBuildingType_MC_B_RELIVE {
		return
	}

	t.mode = newMode
	c.updateViewMsg()
	succ = true

	c.broadcast(mingc_war.NewS2cSceneChangeModeNoticeMsg(idbytes.ToBytes(heroId), int32(t.mode)))
	return
}

func (c *McWarScene) troop(heroId int64) (t *McWarTroop, isAtk bool) {
	if t, ok := c.atkTroops[heroId]; ok {
		return t, true
	}
	if t, ok := c.defTroops[heroId]; ok {
		return t, false
	}
	return
}

// ********** building **********

func (c *McWarScene) initBuilding(sceneData *mingcdata.MingcWarSceneData, buildingGetter func(uint64) *mingcdata.MingcWarBuildingData) {
	c.buildings = make(map[cb.Cube]*McWarBuilding)
	c.atkReliveBuilding = newMcWarBuilding(true, sceneData.AtkRelivePos, sceneData.AtkReliveName, buildingGetter(uint64(shared_proto.MingcWarBuildingType_MC_B_RELIVE)))
	c.buildings[sceneData.AtkRelivePos] = c.atkReliveBuilding
	c.atkHomeBuilding = newMcWarBuilding(true, sceneData.AtkHomePos, sceneData.AtkHomeName, buildingGetter(uint64(shared_proto.MingcWarBuildingType_MC_B_HOME)))
	c.buildings[sceneData.AtkHomePos] = c.atkHomeBuilding

	for i, pos := range sceneData.AtkCastlePos {
		c.buildings[pos] = newMcWarBuilding(true, pos, sceneData.AtkCastleName[i], buildingGetter(uint64(shared_proto.MingcWarBuildingType_MC_B_CASTLE)))
	}
	for i, pos := range sceneData.AtkGatePos {
		c.buildings[pos] = newMcWarBuilding(true, pos, sceneData.AtkGateName[i], buildingGetter(uint64(shared_proto.MingcWarBuildingType_MC_B_GATE)))
	}

	c.defReliveBuilding = newMcWarBuilding(false, sceneData.DefRelivePos, sceneData.DefReliveName, buildingGetter(uint64(shared_proto.MingcWarBuildingType_MC_B_RELIVE)))
	c.buildings[sceneData.DefRelivePos] = c.defReliveBuilding
	c.defHomeBuilding = newMcWarBuilding(false, sceneData.DefHomePos, sceneData.DefHomeName, buildingGetter(uint64(shared_proto.MingcWarBuildingType_MC_B_HOME)))
	c.buildings[sceneData.DefHomePos] = c.defHomeBuilding

	for i, pos := range sceneData.DefCastlePos {
		c.buildings[pos] = newMcWarBuilding(false, pos, sceneData.DefCastleName[i], buildingGetter(uint64(shared_proto.MingcWarBuildingType_MC_B_CASTLE)))
	}
	for i, pos := range sceneData.DefGatePos {
		c.buildings[pos] = newMcWarBuilding(false, pos, sceneData.DefGateName[i], buildingGetter(uint64(shared_proto.MingcWarBuildingType_MC_B_GATE)))
	}
	for i, pos := range sceneData.AtkTouShiPos {
		c.buildings[pos] = newMcWarBuilding(true, pos, sceneData.AtkTouShiName[i], buildingGetter(uint64(shared_proto.MingcWarBuildingType_MC_B_TOU_SHI)))
		c.buildings[pos].touShiTargets = sceneData.GetTouShiTarget(pos)
	}
	for i, pos := range sceneData.DefTouShiPos {
		c.buildings[pos] = newMcWarBuilding(false, pos, sceneData.DefTouShiName[i], buildingGetter(uint64(shared_proto.MingcWarBuildingType_MC_B_TOU_SHI)))
		c.buildings[pos].touShiTargets = sceneData.GetTouShiTarget(pos)
	}

}

func doFight(b *McWarBuilding, troop1, troop2 *McWarTroop, troop1IsInvade bool, scene *McWarScene, ctime time.Time) (troop1Win bool) {
	var atkWallAtk, defWallAtk bool
	if b.data.WallAtk {
		atkWallAtk, defWallAtk = !troop1IsInvade, troop1IsInvade
	}
	tfctx := entity.NewTlogFightContext(operate_type.BattleMcWar, 0, 0, 0)

	t1Proto := troop1.genCombatPlayerProto(atkWallAtk)
	t2Proto := troop2.genCombatPlayerProto(defWallAtk)

	var resp *server_proto.CombatXResponseServerProto
	if t1Proto != nil && t2Proto != nil {
		resp = scene.dep.FightX().SendFightRequest(tfctx, b.data.CombatScene, troop1.heroId, troop2.heroId, t1Proto, t2Proto)
	}

	var t1Killed, t1LastBeatKilled, t2Killed, t2LastBeatKilled uint64
	if resp == nil || resp.ReturnCode != 0 {
		logp.Errorf("名城战战斗异常 mcid:%v resp:%v ", scene.mcId, resp)
		// 服务器错误算守方赢
		troop1.reduceSoldierToZero()
	} else {
		if resp.AttackerWin {
			troop1Win = true
			t2OldSoldier := troop2.soldier()
			t2Killed = troop1.reduceSoldierToAlive(resp.AttackerAliveSoldier)
			t1Killed = troop2.reduceSoldierToZero()
			_, t2LastBeatKilled = scene.lastBeatWhenFail(troop2, t2OldSoldier, troop1)
		} else {
			t1OldSoldier := troop1.soldier()
			t2Killed = troop1.reduceSoldierToZero()
			t1Killed = troop2.reduceSoldierToAlive(resp.DefenserAliveSoldier)
			_, t1LastBeatKilled = scene.lastBeatWhenFail(troop1, t1OldSoldier, troop2)
		}
		scene.addFightRecord(troop1, troop2, resp, t1Killed, t2Killed, t1LastBeatKilled, t2LastBeatKilled, ctime, b.pos, false)
		scene.dep.World().Send(troop1.heroId, mingc_war.SCENE_TROOP_RECORD_ADD_NOTICE_S2C)
		scene.dep.World().Send(troop2.heroId, mingc_war.SCENE_TROOP_RECORD_ADD_NOTICE_S2C)
	}

	stationTroopMsg := mingc_war.NewS2cSceneTroopUpdateMsg(idbytes.ToBytes(troop2.heroId), troop2.encode()).Static()
	scene.broadcast(stationTroopMsg)
	moveTroopMsg := mingc_war.NewS2cSceneTroopUpdateMsg(idbytes.ToBytes(troop1.heroId), troop1.encode()).Static()
	scene.broadcast(moveTroopMsg)

	return
}

// 舍命一击伤害
func (scene *McWarScene) lastBeatWhenFail(atk *McWarTroop, atkSoldier uint64, def *McWarTroop) (succ bool, hurt uint64) {
	var baiZhanLevel uint64

	if atk.heroId > 0 {
		if hero := scene.dep.HeroSnapshot().Get(atk.heroId); hero != nil {
			baiZhanLevel = hero.BaiZhanJunXianLevel
		}
	} else {
		if d := scene.dep.Datas().GetMingcWarNpcData(atk.npcDataId); d != nil {
			baiZhanLevel = d.BaiZhanLevel
		}
	}

	if baiZhanLevel <= 0 {
		return
	}

	var percent *data.Amount
	for i := baiZhanLevel; i > 0; i-- {
		d := scene.dep.Datas().GetMingcWarTroopLastBeatWhenFailData(i)
		if d == nil {
			continue
		}
		if atkSoldier < d.SoliderAmount {
			continue
		}
		if scene.isSelfCaptainAtkBack(atk.heroId) {
			percent = d.AtkBackHurtPercent
		} else {
			percent = d.HurtPercent
		}

		break
	}

	if percent == nil || percent.Percent <= 0 {
		return
	}

	succ = true
	hurt = def.hurtByPercentToAlive(percent)

	logrus.Debugf("名城战舍命一击 atk:%v def:%v hurt:%v defsoldier:%v", atk.heroId, def.heroId, hurt, def.soldier())

	return
}

func (scene *McWarScene) isSelfCaptainAtkBack(heroId int64) bool {
	if _, isAtk := scene.troop(heroId); !isAtk {
		return false
	}

	mcId := scene.dep.Country().LockHeroCapital(heroId)
	return scene.mcId == mcId
}

func (scene *McWarScene) touShiBuildingTurnTo(heroId int64, pos cb.Cube, left bool, ctime time.Time) (succ bool, newTargetIndex int, turnEndTime time.Time) {
	b, ok := scene.buildings[pos]
	if !ok || b.data.Type != shared_proto.MingcWarBuildingType_MC_B_TOU_SHI {
		return
	}

	if ctime.Before(b.touShiTurnEndTime) || ctime.Before(b.touShiPrepareEndTime) {
		return
	}

	t, ok := b.troops[heroId]
	if !ok || t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION {
		return
	}

	b.touShiTurnTo(left, ctime, scene.dep.Datas().MingcMiscData().TouShiBuildingTurnDuration)
	b.updateLastFireHero(t, scene.dep.Datas().MiscConfig())
	scene.updateViewMsg()

	x, y := b.pos.XYI32()
	scene.broadcast(mingc_war.NewS2cSceneTouShiBuildingTurnToNoticeMsg(x, y, left, int32(b.touShiTargetIndex), timeutil.Marshal32(b.touShiTurnEndTime)))

	newTargetIndex = b.touShiTargetIndex
	turnEndTime = b.touShiTurnEndTime
	succ = true

	return
}

func (scene *McWarScene) touShiBuildingFire(heroId int64, pos cb.Cube, ctime time.Time) (succ bool, bombExplodeTime time.Time) {
	b, ok := scene.buildings[pos]
	if !ok || b.data.Type != shared_proto.MingcWarBuildingType_MC_B_TOU_SHI {
		return
	}

	if b.touShiTargetIndex < 0 || b.touShiTargetIndex >= len(b.touShiTargets) {
		return
	}
	target, ok := scene.buildings[b.touShiTargets[b.touShiTargetIndex]]
	if !ok {
		return
	}

	if ctime.Before(b.touShiTurnEndTime) || ctime.Before(b.touShiPrepareEndTime) {
		return
	}

	t, ok := b.troops[heroId]
	if !ok || t.action.getState() != shared_proto.MingcWarTroopState_MC_TP_STATION {
		return
	}

	miscData := scene.dep.Datas().MingcMiscData()
	b.touShiPrepareEndTime = ctime.Add(miscData.TouShiBuildingPrepareDuration)
	bomb := newTouShiBomb(t, ctime, miscData)
	target.explodeBombs = append(target.explodeBombs, bomb)
	b.updateLastFireHero(t, scene.dep.Datas().MiscConfig())
	scene.updateViewMsg()

	x, y := b.pos.XYI32()
	scene.broadcast(mingc_war.NewS2cSceneTouShiBuildingFireNoticeMsg(x, y, int32(b.touShiTargetIndex), timeutil.Marshal32(b.touShiPrepareEndTime), timeutil.Marshal32(bomb.explodeTime), b.lastFireHeroName, b.lastFireHeroCountry))

	bombExplodeTime = bomb.explodeTime
	succ = true

	return
}

func (s *McWarScene) drum(heroId int64, ctime time.Time) (succ bool, toAdd *shared_proto.SpriteStatProto, desc string, nextDrumTime time.Time) {
	t, isAtk := s.troop(heroId)
	if t == nil {
		return
	}

	nextDrumTime = ctime.Add(s.dep.Datas().MingcMiscData().DurmDuration)
	t.nextDrumTime = nextDrumTime
	if t.rankObj != nil {
		t.rankObj.increaseDrumTimes()
		s.troopsRank.needSort = true
	}

	if isAtk {
		s.atkDurmTimes++
		succ, desc, toAdd = s.dep.Datas().GetMingcWarSceneData(s.mcId).DrumStat(s.atkDurmTimes)
		if s.isSelfCaptainAtkBack(heroId) {
			toAdd = data.AppendSpriteStatProto(toAdd, toAdd)
		}
		if succ {
			s.atkDurmStat = data.AppendSpriteStatProto(s.atkDurmStat, toAdd)
			s.broadcast(mingc_war.NewS2cSceneDrumNoticeMsg(idbytes.ToBytes(heroId), isAtk, u64.Int32(s.atkDurmTimes), toAdd, s.atkDurmStat))
		}
	} else {
		s.defDurmTimes++
		succ, desc, toAdd = s.dep.Datas().GetMingcWarSceneData(s.mcId).DrumStat(s.defDurmTimes)
		if s.isSelfCaptainAtkBack(heroId) {
			toAdd = data.AppendSpriteStatProto(toAdd, toAdd)
		}
		if succ {
			s.defDurmStat = data.AppendSpriteStatProto(s.defDurmStat, toAdd)
			s.broadcast(mingc_war.NewS2cSceneDrumNoticeMsg(idbytes.ToBytes(heroId), isAtk, u64.Int32(s.atkDurmTimes), toAdd, s.atkDurmStat))
		}
	}

	succ = true
	s.updateViewMsg()

	return
}

func (s *McWarScene) addDrumStat() {
	if s.atkDurmStat != nil {
		s.addCampDrumStat(s.atkTroops, s.atkDurmStat)
	}

	if s.defDurmStat != nil {
		s.addCampDrumStat(s.defTroops, s.defDurmStat)
	}
}

func (s *McWarScene) addCampDrumStat(troops map[int64]*McWarTroop, stat *shared_proto.SpriteStatProto) {
	var bytes [][]byte
	for _, t := range troops {
		t.addStat(stat)
		if b, err := t.encode().Marshal(); err == nil {
			bytes = append(bytes, b)
		} else {
			logrus.WithError(err).Errorf("addDrumStat troop.Marshal() 异常")
		}
	}
	s.broadcast(mingc_war.NewS2cSceneDrumAddStatNoticeMsg(bytes))
}

// 更新战斗相关的记录
func (scene *McWarScene) addFightRecord(atk, def *McWarTroop, resp *server_proto.CombatXResponseServerProto, atkKilled, defKilled, atkLastBeatKilled, defLastBeatKilled uint64, ctime time.Time, buildingPos cb.Cube, touShi bool) {
	// 结算
	atkAllKilled := atkKilled + atkLastBeatKilled
	defAllKilled := defKilled + defLastBeatKilled

	if atkGuild, ok := scene.record.guilds[atk.gid]; ok {
		atkGuild.killedAmount += atkAllKilled
		atkGuild.woundedAmount += defAllKilled
	}
	if defGuild, ok := scene.record.guilds[def.gid]; ok {
		defGuild.woundedAmount += atkAllKilled
		defGuild.killedAmount += defAllKilled
	}

	atk.killAmount += atkAllKilled
	atk.woundedAmount += defAllKilled

	def.killAmount += defAllKilled
	def.woundedAmount += atkAllKilled

	// 排行榜
	if atk.rankObj != nil {
		multiKillRefreshed := atk.rankObj.refresh4FightResult(atkAllKilled, defAllKilled, resp.AttackerWin)
		if multiKillRefreshed {
			scene.dep.World().Send(atk.heroId, mingc_war.NewS2cCurMultiKillMsg(u64.Int32(atk.rankObj.multiKill)))
			if scene.dep.Datas().GetMingcWarMultiKillData(atk.rankObj.multiKill) != nil {
				scene.broadcast(mingc_war.NewS2cSpecialMultiKillMsg(idbytes.ToBytes(atk.heroId), u64.Int32(atk.rankObj.multiKill)))
			}
		}
		scene.troopsRank.needSort = true
	}
	if def.rankObj != nil {
		multiKillRefreshed := def.rankObj.refresh4FightResult(defAllKilled, atkAllKilled, !resp.AttackerWin)
		if multiKillRefreshed {
			scene.dep.World().Send(def.heroId, mingc_war.NewS2cCurMultiKillMsg(u64.Int32(def.rankObj.multiKill)))
			if scene.dep.Datas().GetMingcWarMultiKillData(def.rankObj.multiKill) != nil {
				scene.broadcast(mingc_war.NewS2cSpecialMultiKillMsg(idbytes.ToBytes(def.heroId), u64.Int32(def.rankObj.multiKill)))
			}
		}
		scene.troopsRank.needSort = true
	}

	maxLen := u64.Int(scene.dep.Datas().MingcMiscData().SceneRecordMaxLen)
	// 攻方部队对战详细记录
	atkp := &shared_proto.McWarTroopRecordProto{}
	atkp.McId = u64.Int32(scene.mcId)
	atkp.TouShi = touShi
	atkp.Time = timeutil.Marshal32(ctime)
	atkp.Atk = true
	atkp.BuildingX, atkp.BuildingY = buildingPos.XYI32()
	atkp.Target = def.hero
	atkp.Win = resp.AttackerWin
	atkp.Killed = u64.Int32(atkKilled)
	atkp.LastBeatKilled = u64.Int32(atkLastBeatKilled)
	atkp.Wounded = u64.Int32(defAllKilled)
	atkp.LastBeatKilled = u64.Int32(atkLastBeatKilled)
	atkp.Combat = &shared_proto.CombatShareProto{}
	atkp.Combat.Link = resp.Link
	atkp.Combat.IsAttacker = true
	atkp.Combat.Type = shared_proto.CombatType_SINGLE
	recLen := len(atk.records)
	if recLen >= maxLen {
		idx := recLen - maxLen + 1
		atk.records = atk.records[idx:maxLen]
	}
	atk.records = append(atk.records, atkp)

	// 守方部队对战详细记录
	defp := &shared_proto.McWarTroopRecordProto{}
	defp.McId = u64.Int32(scene.mcId)
	defp.TouShi = touShi
	defp.Time = timeutil.Marshal32(ctime)
	defp.BuildingX, defp.BuildingY = buildingPos.XYI32()
	defp.Target = atk.hero
	defp.Win = !resp.AttackerWin
	defp.Killed = u64.Int32(defKilled)
	defp.LastBeatKilled = u64.Int32(defLastBeatKilled)
	defp.Wounded = u64.Int32(atkAllKilled)
	defp.Combat = &shared_proto.CombatShareProto{}
	defp.Combat.Link = resp.Link
	defp.Combat.IsAttacker = false
	defp.Combat.Type = shared_proto.CombatType_SINGLE
	recLen = len(def.records)
	if recLen >= maxLen {
		idx := recLen - maxLen + 1
		def.records = def.records[idx:maxLen]
	}
	def.records = append(def.records, defp)
}

func distance(start, dest cb.Cube) float64 {
	x1, y1 := start.XY()
	x2, y2 := dest.XY()
	x := imath.Abs(x1 - x2)
	y := imath.Abs(y1 - y2)
	return math.Sqrt(float64(x*x + y*y))
}

func moveDuration(dist float64, speed float64) (seconds int64) {
	return i64.Max(1, int64(dist/speed)) * int64(time.Second)
}

func moveEndTime(start, dest cb.Cube, speed float64, ctime time.Time) time.Time {
	return ctime.Add(time.Duration(moveDuration(distance(start, dest), speed)))
}

func getAndSetSoldierCount(c *shared_proto.CaptainInfoProto, newCount uint64) (oldCount uint64) {
	oldCount = u64.FromInt32(c.Soldier)

	c.Soldier = u64.Int32(newCount)
	c.FightAmount = data.ProtoFightAmount(c.TotalStat, c.Soldier, c.SpellFightAmountCoef)
	return
}

func (s *McWarScene) viewTroopRecord(heroId int64) (records []*shared_proto.McWarTroopRecordProto) {
	t, _ := s.troop(heroId)
	if t == nil {
		return
	}
	records = t.records
	return
}

func (s *McWarScene) recordAtkChat(proto *shared_proto.ChatMsgProto) {
	s.atkChatRecord.AddChat(proto)
}

func (s *McWarScene) recordDefChat(proto *shared_proto.ChatMsgProto) {
	s.defChatRecord.AddChat(proto)
}

func (s *McWarScene) catchAtkChatRecord(minChatId int64) (sendMsg pbutil.Buffer) {
	if minChatId == 0 {
		sendMsg = s.atkChatRecord.GetFirstChatRecord()
	} else {
		sendMsg = s.atkChatRecord.GetChatRecored(minChatId)
	}
	return
}

func (s *McWarScene) catchDefChatRecord(minChatId int64) (sendMsg pbutil.Buffer) {
	if minChatId == 0 {
		sendMsg = s.defChatRecord.GetFirstChatRecord()
	} else {
		sendMsg = s.defChatRecord.GetChatRecored(minChatId)
	}
	return
}

func (s *McWarScene) sort4TroopRank() {
	if !s.troopsRank.needSort {
		return
	}
	s.troopsRank.needSort = false
	s.troopsRank.sort()
}
