package realm

import (
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/entity/npcid"
	"math/rand"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/gen/pb/misc"
)

func (t *troop) GetAssembly() *assembly {
	return t.assembly
}

func (t *troop) getFacadeTroop() *troop {
	if t.assembly != nil {
		return t.assembly.self
	}
	return t
}

func (t *troop) walkAll(f func(st *troop) (toContinue bool)) {

	if !f(t) {
		return
	}

	t.walkMember(f)
}

func (t *troop) walkMember(f func(st *troop) (toContinue bool)) {

	if t.assembly != nil && t.assembly.self == t {
		for _, st := range t.assembly.member {
			if st != nil {
				if !f(st) {
					return
				}
			}
		}
	}
}

func (t *troop) isAlive() (bool) {
	if t.HasAliveSoldier() {
		return true
	}

	if t.assembly != nil {
		for _, st := range t.assembly.member {
			if st != nil && st.HasAliveSoldier() {
				return true
			}
		}
	}

	return false
}

func (t *troop) TroopCount() uint64 {
	if t.assembly != nil {
		return t.assembly.Count()
	}
	return 1
}

func (t *troop) getAllTroops() []*troop {
	array := []*troop{t}
	if t.assembly != nil {
		for _, st := range t.assembly.member {
			if st != nil {
				array = append(array, st)
			}
		}
	}
	return array
}

func (t *troop) NextRobBaowuTroop() *troop {
	if t.assembly == nil {
		return t
	}

	// 从掠夺宝物个数最少的几个人中随机一个出来
	var minBaowuCount uint64
	var minBaowuTroops []*troop
	t.walkAll(func(st *troop) (toContinue bool) {
		if len(minBaowuTroops) <= 0 {
			minBaowuCount = st.accumRobBaowuCount
			minBaowuTroops = append(minBaowuTroops, st)
		} else {
			switch {
			case st.accumRobBaowuCount < minBaowuCount:
				minBaowuTroops = minBaowuTroops[0:]

				minBaowuCount = st.accumRobBaowuCount
				minBaowuTroops = append(minBaowuTroops, st)

			case st.accumRobBaowuCount == minBaowuCount:
				minBaowuTroops = append(minBaowuTroops, st)
			}
		}
		return true
	})

	n := len(minBaowuTroops)
	if n > 0 {
		return minBaowuTroops[rand.Intn(n)]
	}

	return t
}

func newAssembly(self *troop, totalCount uint64) *assembly {
	assemblyIdBytes := self.IdBytes()
	return &assembly{
		self:       self,
		totalCount: totalCount,
		changedMsg: region.NewS2cShowAssemblyChangedMsg(assemblyIdBytes).Static(),
	}
}

type assembly struct {
	// 集结者
	self *troop

	// 其他集结队伍（包含集结者）
	member []*troop

	// 集结属性
	addedStat *data.SpriteStat

	// 集结总共可以集结几个队伍
	totalCount uint64

	changedMsg pbutil.Buffer
}

//func (assembly *assembly) walkAll(f func(st *troop) (toContinue bool)) {
//
//	if !f(assembly.self) {
//		return
//	}
//
//	assembly.walkMember(f)
//}
//
//func (assembly *assembly) walkMember(f func(st *troop) (toContinue bool)) {
//
//	for _, st := range assembly.member {
//		if st != nil {
//			if !f(st) {
//				return
//			}
//		}
//	}
//}

func (assembly *assembly) broadcast(r *Realm, toSend pbutil.Buffer) {
	toSend = toSend.Static()

	assembly.self.walkAll(func(st *troop) (toContinue bool) {
		r.services.world.Send(st.startingBase.Id(), toSend)
		return true
	})
}

func (assembly *assembly) Count() (count uint64) {
	count = 1
	for _, v := range assembly.member {
		if v != nil {
			count++
		}
	}
	return
}

func (assembly *assembly) TotalCount() uint64 {
	return assembly.totalCount
}

func (assembly *assembly) isBaseJoined(baseId int64) bool {
	for _, t := range assembly.member {
		if t != nil && t.startingBase.Id() == baseId {
			return true
		}
	}
	return false
}

func (assembly *assembly) broadcastChanged(r *Realm) {
	assembly.broadcast(r, assembly.changedMsg)

	// 目标也发一下
	r.services.world.Send(assembly.self.originTargetId, assembly.changedMsg)
}

func (assembly *assembly) destroyAndTroopBackHome(r *Realm, fromBaseId int64, fromX, fromY int, ctime time.Time, updateTroop bool) {

	if assembly.self.startingBase.Id() == fromBaseId {
		// 就在自己家，执行到家逻辑
		assembly.self.doReturnedToBase(r, false)
		assembly.self.assembly = nil
	} else {
		// 从目标处来说回家
		assembly.self.updateCaptainStat(assembly.addedStat, nil)
		assembly.self.backHomeFrom(r, fromBaseId, fromX, fromY, 1, ctime)
		assembly.self.assembly = nil
		r.broadcastMaToCared(assembly.self, addTroopTypeUpdate, 0)

		if updateTroop {
			t := assembly.self
			r.heroBaseFuncNotError(t.startingBase.Base(), func(hero *entity.Hero) (heroChanged bool) {
				hero.UpdateTroop(t, false)
				return true
			})
		}
	}

	for i, st := range assembly.member {
		if st == nil {
			continue
		}

		assembly.removeTroopByIndex(i)
		st.updateCaptainStat(assembly.addedStat, nil)

		st.backHomeFrom(r, fromBaseId, fromX, fromY, 1, ctime)
		r.broadcastMaToCared(st, addTroopTypeUpdate, 0)

		t := st
		r.heroBaseFuncNotError(t.startingBase.Base(), func(hero *entity.Hero) (heroChanged bool) {
			hero.UpdateTroop(t, false)
			return true
		})
	}

}

func (assembly *assembly) removeTroopByIndex(index int) {
	if index < 0 || index >= len(assembly.member) {
		return
	}

	toRemove := assembly.member[index]
	assembly.member[index] = nil

	if toRemove != nil {
		toRemove.assembly = nil
	}
}

func (assembly *assembly) addTroop(toAdd *troop) {

	for i, m := range assembly.member {
		if m == nil {
			assembly.member[i] = toAdd
			return
		}
	}

	assembly.member = append(assembly.member, toAdd)
}

func (assembly *assembly) updateAddedStat(r *Realm, skipTroopId int64) {
	builder := data.NewSpriteStatBuilder()
	if o := r.services.baiZhanService.GetBaiZhanObj(assembly.self.startingBase.Id()); o != nil {
		builder.Add(o.LevelData().AssemblyStat)
	}

	for _, m := range assembly.member {
		if m != nil && m.State() == realmface.AssemblyArrived {
			if o := r.services.baiZhanService.GetBaiZhanObj(m.startingBase.Id()); o != nil {
				builder.Add(o.LevelData().AssemblyStat)
			}
		}
	}

	oldStat := assembly.addedStat
	newStat := builder.Build()
	assembly.addedStat = newStat

	if !data.IsEqualsStat(oldStat, newStat) {
		assembly.self.updateCaptainStat(oldStat, newStat)
		for _, m := range assembly.member {
			if m != nil && m.Id() != skipTroopId && m.State() == realmface.AssemblyArrived {
				m.updateCaptainStat(oldStat, newStat)
			}
		}
	}
}

func (t *troop) updateCaptainStat(toRemove, toAdd *data.SpriteStat) {

	for _, c := range t.captains {
		builder := data.NewSpriteStatBuilder()
		builder.AddProto(c.Proto().TotalStat)
		builder.Add(toAdd)
		builder.Sub(toRemove)

		c.Proto().TotalStat = builder.Build().Encode()

		if c.Proto().Soldier > 0 {
			c.Proto().FightAmount = data.ProtoFightAmount(c.Proto().TotalStat, c.Proto().Soldier, c.Proto().SpellFightAmountCoef)
		}
	}
}

func (as *assembly) encodeInfoProto(r *Realm) *shared_proto.AssemblyInfoProto {
	proto := &shared_proto.AssemblyInfoProto{}

	t := as.self

	_, moveType := stateToActionMoveType(t.state)
	proto.MoveType = moveType.Int32()

	proto.MoveStartTime = timeutil.Marshal32(t.moveStartTime)
	proto.MoveArrivedTime = timeutil.Marshal32(t.moveArriveTime)
	proto.RobbingEndTime = timeutil.Marshal32(t.robbingEndTime)

	proto.Self = t.getStartingBaseSnapshot(r)
	proto.SelfFightAmount = u64.Int32(as.self.FightAmount())

	// 目标
	tb := t.targetBase
	if tb == nil {
		tb = r.getBase(t.originTargetId)
	}
	if tb != nil {
		proto.Target = tb.internalBase.getBaseInfoByLevel(t.targetBaseLevel).EncodeAsHeroBasicSnapshot(r.services.heroSnapshotService.Get)

		npcDataId, npcType := tb.internalBase.getNpcConfig()
		proto.TargetNpcDataId, proto.TargetNpcType = u64.Int32(npcDataId), npcType

		// 士气
		if b := GetXiongNuBase(tb); b != nil {
			proto.TargetMorale = u64.Int32(b.info.Morale())
		}
	} else {

		if npcid.IsNpcId(t.originTargetId) {
			npcDataId := npcid.GetNpcDataId(t.originTargetId)
			npcType := npcid.GetNpcIdType(t.originTargetId)
			proto.TargetNpcDataId, proto.TargetNpcType = u64.Int32(npcDataId), npcType

			switch npcType {
			case npcid.NpcType_MultiLevelMonster:
				data := r.services.datas.GetRegionMultiLevelNpcData(npcDataId)
				if data != nil {
					proto.Target = data.GetLevelBaseData(t.targetBaseLevel).Npc.EncodeSnapshot(t.originTargetId, r.id, t.originTargetX, t.originTargetY)
				}
			default:
				data := r.services.datas.GetNpcBaseData(npcDataId)
				if data != nil {
					proto.Target = data.EncodeSnapshot(t.originTargetId, r.id, t.originTargetX, t.originTargetY)
				}
			}

		} else {
			targetSnapshot := r.services.heroSnapshotService.Get(t.originTargetId)
			if targetSnapshot != nil {
				proto.Target = targetSnapshot.EncodeClient()
			}
		}
	}

	proto.TargetBaseX = imath.Int32(t.originTargetX)
	proto.TargetBaseY = imath.Int32(t.originTargetY)

	for _, m := range as.member {
		if m != nil {
			mp := &shared_proto.AssemblyMemberProto{}
			mp.Hero = m.getStartingBaseSnapshot(r)
			mp.FightAmount = u64.Int32(m.FightAmount())

			if m.State() != realmface.AssemblyArrived {
				mp.MoveStartTime = timeutil.Marshal32(m.moveStartTime)
				mp.MoveArrivedTime = timeutil.Marshal32(m.moveArriveTime)
			}

			proto.Member = append(proto.Member, mp)
		}
	}

	return proto
}

func (t *troop) assemblyTimesUpEvent(r *Realm) {
	t.event.RemoveFromQueue()
	t.event = nil

	logrus.WithField("troopid", t.Id()).Debug("处理集结到时间出发")

	if t.state != realmface.Assembly {
		logrus.WithField("troopid", t.Id()).WithField("state", t.state).Error("realm竟然触发了assemblyTimesUpEvent, 但是队伍状态不是Assembly")
		return
	}

	defenserBase := t.targetBase
	if defenserBase == nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了assemblyTimesUpEvent, 但是troop的targetBase为nil")
		return
	}

	assembly := t.GetAssembly()
	if assembly == nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了assemblyTimesUpEvent, 但是 assembly == nil")
		return
	}

	startingBase := t.startingBase

	ctime := r.services.timeService.CurrentTime()

	// 将还未到达的目标遣返
	arrivedMemberCount := 0
	for i, t := range assembly.member {
		if t == nil {
			continue
		}

		if t.State() == realmface.AssemblyArrived {
			arrivedMemberCount++
			continue
		}

		assembly.removeTroopByIndex(i)

		rate := timeutil.Rate(t.moveStartTime, t.moveArriveTime, ctime)
		t.backHomeFrom(r, startingBase.Id(), startingBase.BaseX(), startingBase.BaseY(), rate, ctime)

		r.broadcastMaToCared(t, addTroopTypeUpdate, 0)

		r.heroBaseFuncNotError(t.startingBase.Base(), func(hero *entity.Hero) (heroChanged bool) {
			hero.UpdateTroop(t, false)
			return true
		})

		r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().RealmAssemblyBackTimeout.New().
					WithTroopIndex(entity.GetTroopIndex(t.Id()) + 1).
					JsonString())
		})
	}

	assembly.broadcastChanged(r)

	// 如果集结中只有自己一个人到达，那么解散集结
	if arrivedMemberCount <= 0 {
		t.removeWithoutReturnCaptain(r)
		r.broadcastRemoveMaToCared(t)
		t.clearMsgs()

		r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
			r.doRemoveHeroTroop(hero, result, t, nil)
			result.Ok()
			return
		})

		r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().RealmAssemblyMemberCountNotEnough.New().
					JsonString())
		})
		return
	}

	moveDuration := t.getMoveDuration(r)

	// 出发
	t.state = realmface.MovingToInvade
	t.moveStartTime = ctime
	t.moveArriveTime = ctime.Add(moveDuration)

	t.event = r.newEvent(t.moveArriveTime, t.invasionArrivedEvent)
	t.onChanged()

	// 广播消息
	r.broadcastMaToCared(t, addTroopTypeInvate, 0)

	r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
		hero.UpdateTroop(t, false)
		result.Ok()
		return
	})

	// 对应飘字
	assembly.broadcast(r, misc.NewS2cScreenShowWordsMsg(r.getTextHelp().RealmAssemblyStart.New().JsonString()).Static())
}

func (t *troop) assemblyArrivedEvent(r *Realm) {
	t.event.RemoveFromQueue()
	t.event = nil

	logrus.WithField("troopid", t.Id()).Debug("处理到达集结地点")

	if t.state != realmface.MovingToAssembly {
		logrus.WithField("troopid", t.Id()).WithField("state", t.state).Error("realm竟然触发了assemblyArrivedEvent, 但是队伍状态不是MovingToAssembly")
		return
	}

	defenserBase := t.targetBase
	if defenserBase == nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了assemblyArrivedEvent, 但是troop的targetBase为nil")
		return
	}

	assembly := t.GetAssembly()
	if assembly == nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了assemblyArrivedEvent, 但是 assembly == nil")
		return
	}

	// 到了就到了，改一下状态，没你什么事了

	t.state = realmface.AssemblyArrived
	t.onChanged()

	// 野外的马要删掉
	r.broadcastRemoveSeeMeMsg(t)

	// 将出征时间设置一下
	t.moveStartTime = assembly.self.moveStartTime
	t.moveArriveTime = assembly.self.moveArriveTime

	// 将马的消息更新给自己
	r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
		return t.getUpdateMsgToSelf(r)
	})

	// 更新主马
	assembly.self.onChanged()
	r.broadcastMaToCared(assembly.self, addTroopTypeUpdate, 0)

	r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
		hero.UpdateTroop(t, false)
		result.Ok()
		return
	})

	assembly.updateAddedStat(r, t.Id())
	t.updateCaptainStat(nil, assembly.addedStat)

	assembly.broadcastChanged(r)

	// 对应飘字
	r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
		return misc.NewS2cScreenShowWordsMsg(
			r.getTextHelp().RealmAssemblyArrived.New().
				WithTroopIndex(entity.GetTroopIndex(t.Id()) + 1).
				JsonString())
	})
}
