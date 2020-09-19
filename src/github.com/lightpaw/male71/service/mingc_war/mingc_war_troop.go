package mingc_war

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"time"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/i32"
)

type McWarTroop struct {
	datas iface.ConfigDatas

	heroId   int64
	hero     *shared_proto.HeroBasicProto
	gid      int64
	atk      bool
	joinTime time.Time

	npcDataId uint64

	mode shared_proto.MingcWarModeType

	action Actionable // 部队行为

	woundedAmount   uint64 // 损兵
	killAmount      uint64 // 歼兵
	destroyBuilding uint64 // 破坏

	nextDrumTime time.Time // 下次可以击鼓时间

	rankObj *McWarTroopRankObject

	captains []*shared_proto.CaptainInfoProto // 空位为nil

	records []*shared_proto.McWarTroopRecordProto
}

func newMcWarTroop(rankObj *McWarTroopRankObject, gid int64, captains []*shared_proto.CaptainInfoProto, pos cb.Cube, ctime time.Time, endTime time.Time, datas iface.ConfigDatas) *McWarTroop {
	t := &McWarTroop{}
	t.mode = shared_proto.MingcWarModeType_MC_MT_NORMAL
	t.rankObj = rankObj
	t.heroId = rankObj.heroId
	t.hero = rankObj.heroProto
	t.gid = gid

	t.atk = rankObj.isAtk
	t.captains = captains
	t.joinTime = ctime
	t.datas = datas
	t.nextDrumTime = ctime.Add(datas.MingcMiscData().DurmDuration)

	t.action = newJoinAction(ctime, endTime, pos)
	return t
}

func newMcWarNpcTroop(npcId int64, npc *mingcdata.MingcWarNpcData, atk bool, captains []*shared_proto.CaptainInfoProto, pos cb.Cube, ctime time.Time, endTime time.Time, datas iface.ConfigDatas) *McWarTroop {
	t := &McWarTroop{}
	t.mode = shared_proto.MingcWarModeType_MC_MT_NORMAL
	t.heroId = npcId
	t.hero = npc.Npc.EncodeHeroBasicProto(idbytes.ToBytes(npcId))
	t.gid = npc.GuildId()
	t.atk = atk
	t.captains = captains
	t.joinTime = ctime
	t.datas = datas
	t.npcDataId = npc.Id

	t.action = newJoinAction(ctime, endTime, pos)

	return t
}

func (t *McWarTroop) unmarshal(p *server_proto.McWarTroopServerProto) (hasCaptain bool) {
	t.mode = p.Mode
	t.heroId = p.HeroId
	t.hero = p.Hero
	t.gid = p.Gid
	t.atk = p.Atk
	t.joinTime = timeutil.Unix64(p.JoinTime)
	t.killAmount = p.KillAmount
	t.woundedAmount = p.WoundedAmount
	t.destroyBuilding = p.DestroyBuilding
	t.nextDrumTime = timeutil.Unix64(p.NextDrumTime)
	t.records = p.Record
	t.npcDataId = p.NpcDataId

	var maxIndex int32
	for index, _ := range p.Captain {
		if index > maxIndex {
			maxIndex = index
		}
	}
	t.captains = make([]*shared_proto.CaptainInfoProto, maxIndex+1)
	for index, c := range p.Captain {
		if c != nil {
			t.captains[index] = c
			hasCaptain = true
		}
	}

	switch p.State {
	case shared_proto.MingcWarTroopState_MC_TP_STATION:
		t.action = newStationAction(timeutil.Unix64(p.StartTime), cb.Cube(p.Pos))
	case shared_proto.MingcWarTroopState_MC_TP_RELIVE:
		t.action = newReliveAction(timeutil.Unix64(p.StartTime), timeutil.Unix64(p.EndTime), cb.Cube(p.Pos))
	case shared_proto.MingcWarTroopState_MC_TP_MOVING:
		t.action = newMoveAction(timeutil.Unix64(p.StartTime), timeutil.Unix64(p.EndTime), cb.Cube(p.StartPos), cb.Cube(p.DestPos))
	case shared_proto.MingcWarTroopState_MC_TP_WAIT:
		t.action = newJoinAction(timeutil.Unix64(p.StartTime), timeutil.Unix64(p.EndTime), cb.Cube(p.StartPos))
	}

	return
}

func (t *McWarTroop) encodeServer() *server_proto.McWarTroopServerProto {
	p := &server_proto.McWarTroopServerProto{}
	p.Mode = t.mode
	p.HeroId = t.heroId
	p.Hero = t.hero
	p.Gid = t.gid
	p.Atk = t.atk
	p.JoinTime = timeutil.Marshal64(t.joinTime)
	p.KillAmount = t.killAmount
	p.WoundedAmount = t.woundedAmount
	p.DestroyBuilding = t.destroyBuilding
	p.NextDrumTime = timeutil.Marshal64(t.nextDrumTime)
	p.Record = t.records
	p.NpcDataId = t.npcDataId

	p.Captain = make(map[int32]*shared_proto.CaptainInfoProto)
	for i, c := range t.captains {
		if c == nil {
			continue
		}

		p.Captain[int32(i)] = c
	}

	p.State = t.action.getState()
	p.StartTime = timeutil.Marshal64(t.action.getStartTime())
	p.EndTime = timeutil.Marshal64(t.action.getEndTime())
	p.Pos = uint64(t.action.getPos())

	if o, ok := t.action.(*MoveAction); ok {
		p.StartPos = uint64(o.startPos)
		p.DestPos = uint64(o.destPos)
	}

	return p
}

func (t *McWarTroop) encode() *shared_proto.McWarTroopProto {
	p := &shared_proto.McWarTroopProto{}
	p.Mode = t.mode
	p.Hero = t.hero
	p.Atk = t.atk
	p.JoinTime = timeutil.Marshal32(t.joinTime)
	p.State = t.action.getState()
	p.StateStartTime = timeutil.Marshal32(t.action.getStartTime())
	p.StateEndTime = timeutil.Marshal32(t.action.getEndTime())
	p.StartPosX, p.StartPosY = t.action.getPos().XYI32()
	//p.Record = t.records
	p.NextDrumTime = timeutil.Marshal32(t.nextDrumTime)

	if t.action.getState() == shared_proto.MingcWarTroopState_MC_TP_MOVING {
		p.DestPosX, p.DestPosY = t.action.(*MoveAction).destPos.XYI32()
	}

	for i, c := range t.captains {
		if c == nil {
			continue
		}
		p.CaptainIndex = append(p.CaptainIndex, int32(i+1))
		p.Captains = append(p.Captains, c)
	}

	return p
}

func (t *McWarTroop) isFullSolider() bool {
	for _, c := range t.captains {
		if c == nil {
			continue
		}
		if c.Soldier < c.TotalSoldier {
			return false
		}
	}
	return true
}

func (t *McWarTroop) fullSoliderAmount() (amount uint64) {
	for _, c := range t.captains {
		if c == nil {
			continue
		}
		amount += u64.FromInt32(c.TotalSoldier)
	}
	return
}

func (t *McWarTroop) failToRelive(scene *McWarScene, fightBuilding *McWarBuilding, ctime time.Time, excludeId int64) {
	delete(fightBuilding.troops, t.heroId)
	delete(scene.buildings[t.action.getPos()].troops, t.heroId)

	destPos := scene.getReliveBuildingPos(t)
	scene.buildings[destPos].troops[t.heroId] = t
	t.action = newStationAction(ctime, destPos)

	x, y := destPos.XYI32()
	scene.broadcast(mingc_war.NewS2cSceneMoveStationMsg(idbytes.ToBytes(t.heroId), x, y))

	endTime := t.relive(scene, ctime)
	scene.dep.World().Send(t.heroId, mingc_war.NewS2cSceneTroopReliveMsg(timeutil.Marshal32(endTime)))
}

func (t *McWarTroop) move(scene *McWarScene, startTime, endTime time.Time, destPos cb.Cube, excludeId int64) {
	t.action = newMoveAction(startTime, endTime, t.action.getPos(), destPos)
	startX, startY := t.action.getPos().XYI32()
	destX, destY := destPos.XYI32()
	msg := mingc_war.NewS2cSceneOtherMoveMsg(idbytes.ToBytes(t.heroId), startX, startY, destX, destY, timeutil.Marshal32(endTime)).Static()
	scene.broadcast(msg)
}

func (t *McWarTroop) relive(scene *McWarScene, ctime time.Time) (endTime time.Time) {
	endTime = ctime.Add(scene.dep.Datas().MingcMiscData().ReliveDuration)
	t.action = newReliveAction(ctime, endTime, t.action.getPos())
	msg := mingc_war.NewS2cSceneOtherTroopReliveMsg(idbytes.ToBytes(t.heroId), timeutil.Marshal32(endTime)).Static()
	scene.broadcastExclude(msg, t.heroId)

	if t.heroId > 0 {
		t.updateCaptain(scene)
	}

	return
}

func (t *McWarTroop) updateCaptain(scene *McWarScene) {
	var caps []*shared_proto.CaptainInfoProto
	var updated bool
	scene.dep.HeroData().FuncNotError(t.heroId, func(hero *entity.Hero) (heroChanged bool) {
		for _, oldProto := range t.captains {
			if oldProto == nil {
				caps = append(caps, nil)
				continue
			}

			if c := hero.Military().Captain(u64.FromInt32(oldProto.Id)); c != nil {
				p := c.EncodeInvaseCaptainInfo(hero, true, oldProto.XIndex)
				p.Soldier = oldProto.Soldier
				caps = append(caps, p)

				updated = true
			} else {
				logrus.Warnf("名城战更新武将，找不到武将, 不更新.id:%v", oldProto.Id)
				return
			}
		}
		return
	})

	if updated {
		var drumStat *shared_proto.SpriteStatProto
		if t.atk {
			drumStat = scene.atkDurmStat
		} else {
			drumStat = scene.defDurmStat
		}
		t.addStat(drumStat)

		t.captains = caps
		scene.broadcast(mingc_war.NewS2cSceneTroopUpdateMsg(idbytes.ToBytes(t.heroId), t.encode()).Static())
	}
}

func (t *McWarTroop) genCombatPlayerProto(wall bool) (p *shared_proto.CombatPlayerProto) {
	p = &shared_proto.CombatPlayerProto{}

	p.Hero = t.hero
	p.Troops = make([]*shared_proto.CombatTroopsProto, 0, len(t.captains))
	tfa := data.NewTroopFightAmount()
	for i, c := range t.captains {
		if c == nil || c.Soldier <= 0 {
			continue
		}

		tps := &shared_proto.CombatTroopsProto{}
		tps.FightIndex = imath.Int32(i + 1)
		tps.Captain = c

		p.Troops = append(p.Troops, tps)
		tfa.AddInt32(tps.Captain.FightAmount)
	}
	p.TotalFightAmount = tfa.ToI32()

	if len(p.Troops) <= 0 {
		logrus.Errorf("名城战，生成的proto中，武将数量长度为0 t:%v", t.heroId)
		return nil
	}

	if wall {
		p.WallStat = t.datas.MingcMiscData().WallStat.Encode()
		p.WallFixDamage = u64.Int32(t.datas.MingcMiscData().WallFixDamage)
		p.TotalWallLife = i32.Max(p.WallStat.Strength, 1)
		p.WallLevel = u64.Int32(t.datas.MingcMiscData().WallLevel)
	}

	return
}

func (t *McWarTroop) onUpdate(scene *McWarScene, ctime time.Time) (updated bool) {
	return t.action.onAction(scene, t, ctime)
}

func (t *McWarTroop) hurtByPercentToAlive(percent *data.Amount) (allHurt uint64) {
	defSoliderAliveMap := make(map[int32]int32)
	for _, c := range t.captains {
		if c == nil {
			continue
		}
		hurt := u64.Min(percent.CalculateByPercent(u64.FromInt32(c.Soldier)), u64.FromInt32(c.Soldier))
		defSoliderAliveMap[c.Id] = i32.Max(1, c.Soldier-u64.Int32(hurt))
		allHurt += hurt
	}
	t.reduceSoldierToAlive(defSoliderAliveMap)

	return
}

func (t *McWarTroop) hurtByAmount(amount uint64) (alive bool, allHurt uint64) {
	if t.soldier() <= amount {
		allHurt = t.reduceSoldierToZero()
		return
	}

	alive, allHurt = true, amount
	currAmount := amount
	defSoliderAliveMap := make(map[int32]int32)
	for _, c := range t.captains {
		if currAmount <= 0 {
			break
		}
		if c == nil {
			continue
		}
		var hurt uint64
		if u64.FromInt32(c.Soldier) >= currAmount {
			hurt = currAmount
		} else {
			hurt = u64.FromInt32(c.Soldier)
		}
		defSoliderAliveMap[c.Id] = i32.Max(0, c.Soldier-u64.Int32(hurt))
		currAmount -= hurt
	}
	t.reduceSoldierToAlive(defSoliderAliveMap)

	return
}

func (t *McWarTroop) reduceSoldierToZero() (woundedCount uint64) {
	var deadSoldierCounter uint64

	for _, captain := range t.captains {
		if captain != nil {
			deadSoldierCounter += getAndSetSoldierCount(captain, 0)
		}
	}

	woundedCount = deadSoldierCounter
	return
}

func (t *McWarTroop) reduceSoldierToAlive(aliveSoldierMap map[int32]int32) (woundedCount uint64) {
	var deadSoldierCounter uint64

	for _, captain := range t.captains {
		if captain == nil {
			continue
		}
		if alive, has := aliveSoldierMap[captain.Id]; has && alive > 0 {
			oldCount := getAndSetSoldierCount(captain, uint64(alive))

			if oldCount < uint64(alive) {
				logrus.WithField("captain", captain).WithField("old_count", oldCount).WithField("new_count", alive).Error("打完一场后, 活着的士兵竟然超过了开打时的士兵数")
				getAndSetSoldierCount(captain, oldCount)
				continue
			}
			deadCount := oldCount - uint64(alive)
			deadSoldierCounter += deadCount
		} else {
			// 没找到, 当做全死光了
			oldCount := getAndSetSoldierCount(captain, 0)
			deadSoldierCounter += oldCount
		}
	}

	woundedCount = deadSoldierCounter
	return
}

func (t *McWarTroop) soldier() uint64 {
	var soldier int32
	for _, c := range t.captains {
		if c == nil {
			continue
		}
		soldier += c.Soldier
	}
	return u64.FromInt32(soldier)
}

func (t *McWarTroop) getSpeed() float64 {
	if t.mode == shared_proto.MingcWarModeType_MC_MT_FREE_TANK {
		return t.datas.MingcMiscData().FreeTankSpeed
	}
	return t.datas.MingcMiscData().Speed
}

func (t *McWarTroop) getDestroyProsperity() uint64 {
	if t.mode == shared_proto.MingcWarModeType_MC_MT_FREE_TANK {
		return t.datas.MingcMiscData().FreeTankPerDestroyProsperity
	}
	return t.datas.MingcMiscData().PerDestroyProsperity
}

func (t *McWarTroop) addStat(toAdd *shared_proto.SpriteStatProto) {
	for _, c := range t.captains {
		if c == nil {
			continue
		}
		c.TotalStat = data.AppendSpriteStatProto(c.TotalStat, toAdd)
		c.FightAmount = data.ProtoFightAmount(c.TotalStat, c.Soldier, c.SpellFightAmountCoef)
		c.FullFightAmount = data.ProtoFightAmount(c.TotalStat, c.TotalSoldier, c.SpellFightAmountCoef)
	}
}
