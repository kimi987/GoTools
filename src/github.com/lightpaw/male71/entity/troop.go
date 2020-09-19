package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func NewHeroTroopId(heroId int64, index uint64) int64 {
	return int64(uint64(heroId)<<4 | index)
}

func GetTroopHeroId(troopsId int64) int64 {
	if npcid.IsNpcId(troopsId) {
		return npcid.GetTroopNpcId(troopsId)
	}

	return int64(uint64(troopsId) >> 4)
}

const indexMask = 1<<4 - 1

func GetTroopIndex(troopsId int64) uint64 {
	if npcid.IsNpcId(troopsId) {
		return npcid.GetTroopSequence(troopsId)
	}
	return uint64(troopsId) & indexMask
}

func newTroop(sequence uint64, id int64, troopCaptainCount uint64) *Troop {
	t := &Troop{}
	t.sequence = sequence
	t.id = NewHeroTroopId(id, sequence)
	t.captains = newTroopPosArray(troopCaptainCount)

	// 默认是要自动补兵的
	//t.needRecoverSoldier = true

	return t
}

type Troop struct {
	sequence uint64
	id       int64

	// 这里是部队数据
	captains []*TroopPos // 武将（新的版本）

	// 这里是野外数据，nil表示在家
	invaseInfo *TroopInvaseInfo

	changed bool

	//needRecoverSoldier bool

	fightAmount        uint64 // 队伍战斗力
	fightAmountChanged bool   // 战斗力有更新
}

//SetId 设置id
func (t *Troop) SetId(id int64) {
	t.id = id
}

func (t *Troop) setFightAmountChanged() {
	t.fightAmountChanged = true
}

func (t *Troop) UpdateFightAmountIfChanged() (oldAmount uint64, newAmount uint64) {
	oldAmount = t.fightAmount
	newAmount = t.fightAmount

	if t.fightAmountChanged {
		t.fightAmountChanged = false
		newAmount = t.updateFightAmount()
	}
	return
}

func (t *Troop) updateFightAmount() uint64 {
	tfa := data.NewTroopFightAmount()
	for _, pos := range t.captains {
		v := pos.captain
		if v != nil {
			tfa.Add(v.FullSoldierFightAmount())
		}
	}
	t.fightAmount = tfa.ToU64()
	return t.fightAmount
}

func (t *Troop) ClearChanged() bool {
	if t.changed {
		t.changed = false
		return true
	}
	return false
}

func (t *Troop) setChanged() {
	t.changed = true
	t.setFightAmountChanged()

	//if !t.IsOutside() {
	//	for _, captain := range t.captains {
	//		if captain != nil && captain.Soldier() < captain.SoldierCapcity() {
	//			t.needRecoverSoldier = true
	//			break
	//		}
	//	}
	//}
}

//func (t *Troop) NeedRecoverSoldier() (result bool) {
//	result = t.needRecoverSoldier
//	t.needRecoverSoldier = false
//	return
//}

func (t *Troop) Id() int64 {
	return t.id
}

func (t *Troop) Sequence() uint64 {
	return t.sequence
}

func (t *Troop) GetInvateInfo() *TroopInvaseInfo {
	return t.invaseInfo
}

func (t *Troop) RemoveInvateInfo() {
	t.invaseInfo = nil
}

func (t *Troop) HasSoldier() bool {
	for _, pos := range t.captains {
		v := pos.captain
		if v != nil && v.Soldier() > 0 {
			return true
		}
	}

	return false
}

func (t *Troop) isNotEmpty() bool {
	if t.invaseInfo != nil {
		return true
	}

	for _, pos := range t.captains {
		v := pos.captain
		if v != nil {
			return true
		}
	}

	return false
}

func (t *Troop) CalFightAmount() uint64 {
	tfa := data.NewTroopFightAmount()
	for _, pos := range t.captains {
		c := pos.captain
		if c != nil {
			tfa.Add(c.FightAmount())
		}
	}
	return tfa.ToU64()
}

func (t *Troop) CalDefenseFightAmount(hero *Hero) uint64 {
	tfa := data.NewTroopFightAmount()
	for _, pos := range t.captains {
		c := pos.captain
		if c != nil {
			tfa.Add(c.CalDefenseFightAmountWithSoldier(hero, c.Soldier()))
		}
	}
	return tfa.ToU64()
}

func (t *Troop) FullFightAmount() uint64 {
	return t.fightAmount
}

func (t *Troop) CalFullFightAmount() uint64 {
	tfa := data.NewTroopFightAmount()
	for _, pos := range t.captains {
		c := pos.captain
		if c != nil {
			tfa.Add(c.FullSoldierFightAmount())
		}
	}
	return tfa.ToU64()
}

func (t *Troop) EncodeClient() *shared_proto.HeroTroopProto {
	proto := &shared_proto.HeroTroopProto{}
	proto.Sequence = u64.Int32(GetTroopIndex(t.id))
	proto.IsOutside = t.IsOutside()

	for _, pos := range t.captains {
		captain := pos.captain
		var cid int32
		if captain != nil {
			cid = u64.Int32(captain.Id())
		}
		proto.Captains = append(proto.Captains, cid)
		proto.XIndex = append(proto.XIndex, pos.xIndex)
	}

	return proto
}

func (t *Troop) encodeServer() *server_proto.HeroTroopsServerProto {
	result := &server_proto.HeroTroopsServerProto{}
	result.Sequence = u64.Int32(GetTroopIndex(t.id))

	// 新版武将
	result.Captains, result.CaptainsXIndex = doEncodeCaptainPos(t.captains)

	if info := t.invaseInfo; info != nil {
		result.InvateInfo = info.encodeServer()
	}

	return result
}

func (t *Troop) unmarshal(proto *server_proto.HeroTroopsServerProto, military *Military, hero *Hero, ctime time.Time) {

	doUnmarshalCaptainPos(t.captains, proto.Captains, proto.CaptainsXIndex, hero.military.Captain)
	for _, v := range t.captains {
		if v.captain != nil {
			v.captain.troop = t
		}
	}

	t.updateFightAmount()
	if proto.InvateInfo != nil {
		if t.isNotEmpty() {
			captainIds, xIndex := t.CaptainsIdPos()
			t.invaseInfo = unmarshalInvateInfo(t.id, captainIds, xIndex, proto.InvateInfo)
		} else {
			logrus.Error("玩家队伍包含野外出征信息，但是这个队伍没有武将")
		}
	}
}

func (t *Troop) IsOutside() bool {
	return t.invaseInfo != nil
}

func (t *Troop) Pos() []*TroopPos {
	return t.captains
}

func (t *Troop) CaptainIds() []uint64 {
	ids := make([]uint64, len(t.captains))
	for i, pos := range t.captains {
		v := pos.captain
		if v != nil {
			ids[i] = v.data.Id
		}
	}
	return ids
}

// 新版
func (t *Troop) CaptainsIdPos() (ids []uint64, pos []int32) {
	ids = make([]uint64, len(t.captains))
	pos = make([]int32, len(t.captains))
	for i, v := range t.captains {
		if v.captain != nil {
			ids[i] = v.captain.Id()
			pos[i] = v.xIndex
		}
	}
	return
}

func (troop *Troop) Set(captains []*Captain, xIndex []int32) {
	// 清掉老的数据
	for _, pos := range troop.captains {
		v := pos.captain
		if v != nil && v.troop == troop {
			v.troop = nil
		}
	}

	doSetTroopPos(troop.captains, captains, xIndex)

	// 更新队伍数据
	for _, pos := range troop.captains {
		v := pos.captain
		if v != nil {
			v.troop = troop
		}
	}

	troop.setChanged()
}

func (t *Troop) GetCaptain(index uint64) *Captain {
	if index < uint64(len(t.captains)) {
		return t.captains[index].captain
	} else {
		logrus.Errorf("entity.Troop.GetCaptain index >= len(t.captains), index:%v len:%v", index, len(t.captains))
	}

	return nil
}

func (t *Troop) GetPos(index uint64) *TroopPos {
	if index < uint64(len(t.captains)) {
		return t.captains[index]
	} else {
		logrus.Errorf("entity.Troop.GetCaptain index >= len(t.captains), index:%v len:%v", index, len(t.captains))
	}
	return nil
}

func (t *Troop) GetCaptainPos(index uint64) (*Captain, int32) {
	if index < uint64(len(t.captains)) {
		pos := t.captains[index]
		return pos.captain, pos.xIndex
	} else {
		logrus.Errorf("entity.Troop.GetCaptain index >= len(t.captains), index:%v len:%v", index, len(t.captains))
	}

	return nil, 0
}

func SwapTroopsCaptain(t1, t2 *Troop, index1, index2 uint64) bool {

	if t1.isValidIndex(index1) && t2.isValidIndex(index2) {
		captain1, xIndex1 := t1.GetCaptainPos(index1)
		captain2, xIndex2 := t2.SetCaptain(index2, captain1, xIndex1)
		t1.SetCaptain(index1, captain2, xIndex2)
		return true
	}
	return false
}

func (t *Troop) SetCaptain(index uint64, captain *Captain, xIndex int32) (*Captain, int32) {
	return t.setCaptain(index, captain, xIndex, false)
}

func (t *Troop) SetCaptainIfAbsent(index uint64, captain *Captain, xIndex int32) (*Captain, int32) {
	return t.setCaptain(index, captain, xIndex, true)
}

func (t *Troop) isValidIndex(index uint64) bool {
	return index < uint64(len(t.captains))
}

func (t *Troop) setCaptain(index uint64, captain *Captain, xIndex int32, ifAbsent bool) (*Captain, int32) {
	if t.isValidIndex(index) {
		pos := t.captains[index]

		if ifAbsent && pos.captain != nil {
			return pos.captain, pos.xIndex
		}

		oldCaptain := pos.captain
		oldXIndex := pos.xIndex

		pos.captain = captain
		pos.xIndex = xIndex
		if captain != nil {
			captain.troop = t
		}

		t.setChanged()

		return oldCaptain, oldXIndex
	} else {
		logrus.Errorf("entity.Troop.SetCaptain index >= len(t.captains), index:%v len:%v", index, len(t.captains))
	}

	return nil, 0
}

//func (t *Troop) SetCaptains(toSet []*Captain) {
//
//	n := len(t.captains)
//	if n != len(toSet) {
//		logrus.WithField("len(t.captains)", n).WithField("len(toSet)", len(toSet)).Error("entity.Troop.SetCaptains len(toSet) != len(t.captain)")
//	}
//
//	count := imath.Min(n, len(toSet))
//	for i := 0; i < count; i++ {
//		captain := toSet[i]
//		t.captains[i] = captain
//
//		if captain != nil {
//			captain.troop = t
//		}
//	}
//
//	for i := count; i < n; i++ {
//		t.captains[i] = nil
//	}
//
//	t.setChanged()
//}

//type TroopInvaseInfo struct {
//}

// 存在英雄里, 只供存db用途
type TroopInvaseInfo struct {
	// 部队信息
	troopId  int64
	captains []uint64
	xIndex   []int32

	// 出征信息
	realmId int64

	originTargetId  int64
	targetBaseId    int64
	targetBaseLevel uint64
	state           realmface.TroopState

	createTime               time.Time
	moveStartTime            time.Time
	moveArriveTime           time.Time
	robbingEndTime           time.Time
	nextReduceProsperityTime time.Time
	nextAddHateTime          time.Time
	nextRobBaowuTime         time.Time

	backHomeTargetX, backHomeTargetY int // 正在回城的时候, 之前目标的坐标

	ownerCanSeeTarget bool // 英雄自己可以看到

	accumRobPrize         *shared_proto.PrizeProto
	accumReduceProsperity uint64

	accumAddHate uint64

	// 发起集结的id，这个是自己的话，说明这个集结是我发起的，0表示非集结队伍
	assemblyId       int64
	assemblyTargetId int64

	dialogue uint64

	npcTimes uint64
}

func (t *TroopInvaseInfo) NpcTimes() uint64 {
	return t.npcTimes
}

func (t *TroopInvaseInfo) AssemblyId() int64 {
	return t.assemblyId
}

func (t *TroopInvaseInfo) AssemblyTargetId() int64 {
	return t.assemblyTargetId
}

func (t *TroopInvaseInfo) Dialogue() uint64 {
	return t.dialogue
}

func (t *TroopInvaseInfo) Id() int64 {
	return t.troopId
}

func (t *TroopInvaseInfo) Captains() []uint64 {
	return t.captains
}

func (t *TroopInvaseInfo) CaptainXIndex() []int32 {
	return t.xIndex
}

func (t *TroopInvaseInfo) RegionID() int64 {
	return t.realmId
}

func (t *TroopInvaseInfo) OriginTargetID() int64 {
	return t.originTargetId
}

func (t *TroopInvaseInfo) TargetBaseID() int64 {
	return t.targetBaseId
}

func (t *TroopInvaseInfo) TargetBaseLevel() uint64 {
	return t.targetBaseLevel
}

func (t *TroopInvaseInfo) State() realmface.TroopState {
	return t.state
}

func (t *TroopInvaseInfo) BackHomeTargetX() int {
	return t.backHomeTargetX
}

func (t *TroopInvaseInfo) BackHomeTargetY() int {
	return t.backHomeTargetY
}

func (t *TroopInvaseInfo) TargetIsOwnerCanSee() bool {
	return t.ownerCanSeeTarget
}

func (t *TroopInvaseInfo) CreateTime() time.Time {
	return t.createTime
}

func (t *TroopInvaseInfo) MoveStartTime() time.Time {
	return t.moveStartTime
}

func (t *TroopInvaseInfo) MoveArriveTime() time.Time {
	return t.moveArriveTime
}

func (t *TroopInvaseInfo) RobbingEndTime() time.Time {
	return t.robbingEndTime
}

func (t *TroopInvaseInfo) NextReduceProsperityTime() time.Time {
	return t.nextReduceProsperityTime
}

func (t *TroopInvaseInfo) NextAddHateTime() time.Time {
	return t.nextAddHateTime
}

func (t *TroopInvaseInfo) NextRobBaowuTime() time.Time {
	return t.nextRobBaowuTime
}

func (t *TroopInvaseInfo) AccumRobPrize() *shared_proto.PrizeProto {
	return t.accumRobPrize
}

func (t *TroopInvaseInfo) AccumReduceProsperity() uint64 {
	return t.accumReduceProsperity
}

func (t *TroopInvaseInfo) AccumAddHate() uint64 {
	return t.accumAddHate
}

func (t *TroopInvaseInfo) ClearAccumAddHate() uint64 {
	a := t.accumAddHate
	t.accumAddHate = 0
	return a
}

func (t *TroopInvaseInfo) AddHate(toAdd uint64) {
	t.accumAddHate += toAdd
}

func (t *TroopInvaseInfo) encodeServer() *server_proto.TroopsInvateProto {
	result := &server_proto.TroopsInvateProto{
		RealmId:         t.realmId,
		OriginTargetId:  t.originTargetId,
		TargetId:        t.targetBaseId,
		TargetBaseLevel: t.targetBaseLevel,

		State:                    int32(t.state),
		CreateTime:               timeutil.Marshal64(t.createTime),
		MoveStartTime:            timeutil.Marshal64(t.moveStartTime),
		MoveArriveTime:           timeutil.Marshal64(t.moveArriveTime),
		RobbingEndTime:           timeutil.Marshal64(t.robbingEndTime),
		NextReduceProsperityTime: timeutil.Marshal64(t.nextReduceProsperityTime),
		NextAddHateTime:          timeutil.Marshal64(t.nextAddHateTime),

		BackHomeTargetX: int32(t.backHomeTargetX),
		BackHomeTargetY: int32(t.backHomeTargetY),

		OwnerCanSeeTarget: t.ownerCanSeeTarget,

		AccumRobPrize:         t.accumRobPrize,
		AccumReduceProsperity: t.accumReduceProsperity,
		AccumAddHate:          t.accumAddHate,
		AssemblyId:            t.assemblyId,
		AssemblyTargetId:      t.assemblyTargetId,
		Dialogue:              t.dialogue,
		NpcTimes:              t.npcTimes,
	}

	return result
}

func unmarshalInvateInfo(troopId int64, captains []uint64, xIndex []int32, proto *server_proto.TroopsInvateProto) *TroopInvaseInfo {
	result := &TroopInvaseInfo{
		troopId:  troopId,
		captains: captains,
		xIndex:   xIndex,

		realmId:         proto.RealmId,
		originTargetId:  proto.OriginTargetId,
		targetBaseId:    proto.TargetId,
		targetBaseLevel: proto.TargetBaseLevel,
		state:           realmface.TroopState(proto.State),

		backHomeTargetX: int(proto.BackHomeTargetX),
		backHomeTargetY: int(proto.BackHomeTargetY),

		ownerCanSeeTarget: proto.OwnerCanSeeTarget,

		createTime:               time.Unix(proto.CreateTime, 0),
		moveStartTime:            time.Unix(proto.MoveStartTime, 0),
		moveArriveTime:           time.Unix(proto.MoveArriveTime, 0),
		robbingEndTime:           time.Unix(proto.RobbingEndTime, 0),
		nextReduceProsperityTime: time.Unix(proto.NextReduceProsperityTime, 0),
		nextAddHateTime:          timeutil.Unix64(proto.NextAddHateTime),
		nextRobBaowuTime:         timeutil.Unix64(proto.NextRobBaowuTime),

		accumRobPrize:         proto.AccumRobPrize,
		accumReduceProsperity: proto.AccumReduceProsperity,
		accumAddHate:          proto.AccumAddHate,
		assemblyId:            proto.AssemblyId,
		assemblyTargetId:      proto.AssemblyTargetId,
		dialogue:              proto.Dialogue,
		npcTimes:              proto.NpcTimes,
	}

	return result
}

// 只更新可能会更新的部分
func (t *TroopInvaseInfo) updateFrom(troop realmface.Troop) {
	if base := troop.TargetBase(); base == nil {
		t.targetBaseId = 0
		t.backHomeTargetX = troop.BackHomeTargetX()
		t.backHomeTargetY = troop.BackHomeTargetY()
	} else {
		t.targetBaseId = base.Id()
	}

	t.state = troop.State()
	t.moveStartTime = troop.MoveStartTime()
	t.moveArriveTime = troop.MoveArriveTime()
	t.robbingEndTime = troop.RobbingEndTime()
	t.nextReduceProsperityTime = troop.NextReduceProsperityTime()
	t.nextAddHateTime = troop.NextAddHateTime()
	t.nextRobBaowuTime = troop.NextRobBaowuTime()

	t.accumRobPrize = troop.AccumRobPrize()
	t.accumReduceProsperity = troop.AccumReduceProsperity()

	t.assemblyId = troop.AssemblyId()
	t.dialogue = troop.Dialogue()
}

func (t *Troop) InitInvateInfo(realmId int64, troop realmface.Troop) {
	// 先找troop

	captainIds, xIndex := t.CaptainsIdPos()

	t.invaseInfo = &TroopInvaseInfo{
		troopId:  t.id,
		captains: captainIds,
		xIndex:   xIndex,

		realmId:           realmId,
		originTargetId:    troop.TargetBase().Id(),
		targetBaseId:      troop.TargetBase().Id(),
		backHomeTargetX:   troop.TargetBase().BaseX(),
		backHomeTargetY:   troop.TargetBase().BaseY(),
		ownerCanSeeTarget: troop.TargetIsOwnerCanSee(),

		state:                 troop.State(),
		createTime:            troop.CreateTime(),
		moveStartTime:         troop.MoveStartTime(),
		moveArriveTime:        troop.MoveArriveTime(),
		robbingEndTime:        troop.RobbingEndTime(),
		accumRobPrize:         troop.AccumRobPrize(),
		accumReduceProsperity: troop.AccumReduceProsperity(),

		assemblyId:       troop.AssemblyId(),
		assemblyTargetId: troop.AssemblyTargetId(),
		npcTimes:         troop.NpcTimes(),
	}

	t.setChanged()
}

// 更新队伍状态, 可选是否更新每个武将带的兵数.
// must hold hero lock
func (hero *Hero) UpdateTroop(troop realmface.Troop, updateSoldier bool) {
	if t := hero.Troop(troop.Id()); t != nil {
		if t.invaseInfo != nil {
			t.invaseInfo.updateFrom(troop)

			if updateSoldier {
				hero.UpdateTroopSoldier(troop)
			}
		} else {
			logrus.WithField("troopid", troop.Id()).Error("Hero.UpdateTroop, 竟然找不到troop的invateInfo")
		}
	} else {
		logrus.WithField("troopid", troop.Id()).Error("Hero.UpdateTroop, 竟然找不到troop")
	}
}

func (hero *Hero) RemoveTroop(troop realmface.Troop, updateSoldier bool) {
	if t := hero.Troop(troop.Id()); t != nil {
		if t.invaseInfo != nil {
			t.RemoveInvateInfo() // 清空出征状态
			troop.RemoveRealmPoints() // 清空地图上的点
			if updateSoldier {
				hero.UpdateTroopSoldier(troop)
			}

			t.setChanged()
		} else {
			logrus.WithField("troopid", troop.Id()).Error("Hero.RemoveTroop, 竟然找不到troop的invateInfo")
		}
	} else {
		logrus.WithField("troopid", troop.Id()).Error("Hero.RemoveTroop, 竟然找不到troop")
	}
}

func (hero *Hero) UpdateTroopSoldier(troop realmface.Troop) {
	for _, captain := range troop.Captains() {
		if c := hero.military.Captain(captain.Id()); c != nil {

			c.SetSoldier(u64.FromInt32(captain.Proto().Soldier))
			c.updateTroopChanged()
		} else {
			logrus.WithField("captain", captain.Id()).WithField("troop", troop.Id()).Error("Hero.UpdateTroopSoldier, 竟然没找到troop中的captain")
		}
	}
}

func (hero *Hero) GetRecruitCaptainTroop() (*Troop, uint64) {

	for _, t := range hero.military.troops {
		if t.IsOutside() {
			continue
		}

		for i, v := range t.captains {
			if v.captain == nil {
				return t, uint64(i)
			}
		}
	}

	return nil, 0
}

//GetInvestigateTroop 获取侦察部队
func (hero *Hero) GetInvestigateTroop(troop_id uint64, hero_id int64) *Troop {
	if hero.military.investigateTroop == nil {
		hero.military.investigateTroop = newTroop(troop_id, hero.Id(), 0)
	}
	return hero.military.investigateTroop
}

//GetTroopOrInvestigateTroop 获取常规军队或者侦察军队
func (hero *Hero) GetTroopOrInvestigateTroop(index uint64) *Troop {
	c := hero.military
	if index < uint64(len(c.troops)) {
		return c.troops[index]
	}

	if hero.military.investigateTroop != nil && GetTroopIndex(hero.military.investigateTroop.Id()) == index {
		//返回侦察部队
		return hero.military.investigateTroop
	}
	return nil
}

func (hero *Hero) GetTroopByIndex(index uint64) *Troop {
	c := hero.military
	if index < uint64(len(c.troops)) {
		return c.troops[index]
	}

	return nil
}

func (hero *Hero) IterTroop(f func(inHero *TroopInvaseInfo) bool) {
	for _, troop := range hero.military.troops {
		if troop.invaseInfo == nil {
			continue
		}

		if ok := f(troop.invaseInfo); !ok {
			troop.RemoveInvateInfo()

			//// 把武将状态设为0
			//for _, cid := range troop.captains {
			//	if c := hero.military.captains[cid]; c != nil {
			//		if c.action == troop.id {
			//			c.SetAction(0)
			//		} else {
			//			logrus.WithField("action", c.action).WithField("troopid", troop.id).WithField("heroid", hero.id).WithField("captainid", cid).Error("troop在开服IterTroop时, 返回不ok, 要重置武将action时, action不等于队伍id. 同时很多个队伍都有同一个武将?")
			//		}
			//	} else {
			//		logrus.WithField("troopid", troop.id).WithField("heroid", hero.id).WithField("captainid", cid).Error("troop在开服IterTroop时, 返回不ok, 但是武将竟然不存在")
			//	}
			//}
		}
	}
}

func (hero *Hero) Troop(id int64) *Troop {
	index := GetTroopIndex(id)
	return hero.GetTroopOrInvestigateTroop(index)
}

func (hero *Hero) Troops() []*Troop {
	return hero.military.troops
}

// 是否存在参与集结的部队
func (hero *Hero) HasAssemblyTroop() bool {
	for _, t := range hero.military.troops {
		if t.invaseInfo != nil && t.invaseInfo.assemblyId != 0 {
			return true
		}
	}
	return false
}

func (hero *Hero) WalkPveTroop(walkFunc func(troop *PveTroop) (endWalk bool)) {
	for _, troop := range hero.military.pveTroops {
		if walkFunc(troop) {
			return
		}
	}
}

func (hero *Hero) PveTroop(troopType shared_proto.PveTroopType) *PveTroop {
	if troopType <= 0 {
		return nil
	}

	if int(troopType) > len(hero.military.pveTroops) {
		return nil
	}

	return hero.military.pveTroops[troopType-1]
}

//func (hero *Hero) NeedRecoverTroopsSoldier() (result bool) {
//	for _, t := range hero.military.troops {
//		if t.NeedRecoverSoldier() {
//			result = true
//		}
//	}
//
//	return
//}

func (hero *Hero) InvateTroopsCount() int {

	count := 0
	for _, t := range hero.military.troops {
		if t.invaseInfo != nil {
			count++
		}
	}

	return count
}

//func (t *Troop) encodeBaseDefenserProto() *shared_proto.BaseDefenserProto {
//	proto := &shared_proto.BaseDefenserProto{}
//
//	for i, pos := range t.captains {
//		c := pos.captain
//		if c == nil {
//
//			continue
//		}
//
//		proto.CaptainIndex = append(proto.CaptainIndex, imath.Int32(i+1))
//		proto.Captains = append(proto.Captains, pos.EncodeCaptainInfo(false))
//	}
//
//	return proto
//}

type GenCombatPlayerFailType int32

const (
	SUCCESS                  = iota // 成功
	SERVER_ERROR                    // 服务器错误
	CAPTAIN_COUNT_NOT_ENOUGH        // 武将人数不够
)

func (hero *Hero) CheckCanGenCombatPlayer(pveTroopType shared_proto.PveTroopType) (failType GenCombatPlayerFailType) {
	return hero.checkCanGenCombatPlayer(pveTroopType, 0)
}

func (hero *Hero) checkCanGenCombatPlayer(pveTroopType shared_proto.PveTroopType, yuanJunCount uint64) (failType GenCombatPlayerFailType) {
	pveTroop := hero.PveTroop(pveTroopType)
	if pveTroop == nil {
		logrus.Errorf("竟然没找到pve队伍, 类型: %v", pveTroopType)
		failType = SERVER_ERROR
		return
	}

	//if heroCaptainCount := uint64(hero.Military().CaptainCount()); heroCaptainCount+yuanJunCount >= pveTroop.TroopData().MinCaptainCount {
	//	if troopCaptainCount := pveTroop.CaptainCount(); troopCaptainCount+yuanJunCount < pveTroop.TroopData().MinCaptainCount {
	//		logrus.Debugf("队伍武将数量不够, 类型: %v, %d, %d", pveTroopType, troopCaptainCount, pveTroop.TroopData().MinCaptainCount)
	//		failType = CAPTAIN_COUNT_NOT_ENOUGH
	//		return
	//	}
	//} else {
	//	if yuanJunCount <= 0 {
	//		// 要全部出战
	//		if troopCaptainCount := pveTroop.CaptainCount(); troopCaptainCount != heroCaptainCount {
	//			logrus.Debugf("在武将数量没满的情况下，队伍人数还是没满, 类型: %v, %d, %d", pveTroopType, troopCaptainCount, heroCaptainCount)
	//			failType = CAPTAIN_COUNT_NOT_ENOUGH
	//			return
	//		}
	//	}
	//}

	failType = SUCCESS
	return
}

func (hero *Hero) EncodeBasicProto(guildGet guildsnapshotdata.Getter) *shared_proto.HeroBasicProto {
	proto := &shared_proto.HeroBasicProto{}

	proto.Id = hero.IdBytes()
	proto.Name = hero.Name()
	proto.Head = hero.Head()
	proto.Level = u64.Int32(hero.Level())
	proto.Male = hero.Male()
	proto.Location = u64.Int32(hero.location)
	proto.CountryId = u64.Int32(hero.CountryId())
	proto.Official = hero.countryMisc.showOfficialType
	proto.VipLevel = u64.Int32(hero.VipLevel())

	if hero.guildId != 0 {
		if g := guildGet(hero.guildId); g != nil {
			proto.GuildId = i64.Int32(g.Id)
			proto.GuildName = g.Name
			proto.GuildFlagName = g.FlagName
		}
	}

	return proto
}

//func (hero *Hero) GenCombatPlayerProto(fullSoldier bool, pveTroopType shared_proto.PveTroopType, getGuildSnapshot guildsnapshotdata.Getter) (player *shared_proto.CombatPlayerProto, failType GenCombatPlayerFailType) {
//	failType = hero.CheckCanGenCombatPlayer(pveTroopType)
//	if failType != SUCCESS {
//		return
//	}
//
//	pveTroop := hero.PveTroop(pveTroopType)
//
//	proto := hero.GenCombatPlayerProtoWithCaptains(fullSoldier, pveTroop.Captains(), getGuildSnapshot)
//	if proto == nil {
//		return nil, SERVER_ERROR
//	}
//
//	return proto, SUCCESS
//}
//
//func (hero *Hero) GenCombatPlayerProtoWithCaptains(fullSoldier bool, captains []*Captain, getGuildSnapshot guildsnapshotdata.Getter) (player *shared_proto.CombatPlayerProto) {
//	protos := make([]*shared_proto.CaptainInfoProto, len(captains))
//	for i, c := range captains {
//		if c == nil {
//			continue
//		}
//
//		protos[i] = c.EncodeCaptainInfo(fullSoldier, 0)
//	}
//
//	player = hero.GenCombatPlayerProtoWithCaptainProtos(protos, getGuildSnapshot)
//
//	return
//}

func (hero *Hero) GenCombatPlayerProtoWithCaptainProtos(captains []*shared_proto.CaptainInfoProto, getGuildSnapshot guildsnapshotdata.Getter) (player *shared_proto.CombatPlayerProto) {
	player = &shared_proto.CombatPlayerProto{}

	player.Hero = hero.EncodeBasicProto(getGuildSnapshot)

	player.Troops = make([]*shared_proto.CombatTroopsProto, 0, len(captains))
	tfa := data.NewTroopFightAmount()
	for i, c := range captains {
		if c == nil {
			continue
		}

		tps := &shared_proto.CombatTroopsProto{}
		tps.FightIndex = imath.Int32(i + 1)
		tps.Captain = c

		player.Troops = append(player.Troops, tps)
		tfa.AddInt32(tps.Captain.FightAmount)
	}
	player.TotalFightAmount = tfa.ToI32()

	if len(player.Troops) <= 0 {
		logrus.Errorf("竟然生成的proto中，武将数量长度竟然为0，前面怎么判断的")
		return nil
	}

	return
}

//func (hero *Hero) BuildPveCaptainProtosWithHeroAndMonster(fullSoldier bool, pveType shared_proto.PveTroopType, toAdds []*monsterdata.MonsterCaptainData, yuanJun bool) (dungeonTroop []*shared_proto.CaptainInfoProto, failType GenCombatPlayerFailType) {
//	failType = hero.checkCanGenCombatPlayer(pveType, uint64(len(toAdds)))
//	if failType != SUCCESS {
//		return
//	}
//
//	pveTroop := hero.PveTroop(pveType)
//	// pveTroop.Captains() 数组的长度就是队伍长度，包含空位，不能加上 len(toAdds)
//	dungeonTroop = make([]*shared_proto.CaptainInfoProto, len(pveTroop.Captains()))
//
//	// 先放援军，援军放完，再放玩家的部队，玩家部队放完，如果跟援军位置冲突的部队，再找空位放进去
//	for _, c := range toAdds {
//		if c == nil {
//			continue
//		}
//		index := c.Index
//		if index <= 0 || int(index) > len(dungeonTroop) {
//			continue
//		}
//
//		dungeonTroop[index-1] = c.EncodeCaptainInfo()
//	}
//
//	var findEmptyPosCaptains []*Captain
//	for i, c := range pveTroop.Captains() {
//		if c == nil {
//			continue
//		}
//
//		if dungeonTroop[i] == nil {
//			dungeonTroop[i] = c.EncodeCaptainInfo(fullSoldier, 0)
//		} else {
//			findEmptyPosCaptains = append(findEmptyPosCaptains, c)
//		}
//	}
//
//	if len(findEmptyPosCaptains) > 0 {
//
//		startIndex := 0
//	out:
//		for _, c := range findEmptyPosCaptains {
//
//			for i := startIndex; i < len(dungeonTroop); i++ {
//				startIndex = i + 1
//
//				if dungeonTroop[i] == nil {
//					dungeonTroop[i] = c.EncodeCaptainInfo(fullSoldier, 0)
//					continue out
//				}
//			}
//
//			// 已经没有空位了
//			break
//		}
//	}
//
//	return
//}

func (hero *Hero) GenCombatPlayerProto(fullSoldier bool, pveTroopType shared_proto.PveTroopType, getGuildSnapshot guildsnapshotdata.Getter) (player *shared_proto.CombatPlayerProto, failType GenCombatPlayerFailType) {
	failType = hero.CheckCanGenCombatPlayer(pveTroopType)
	if failType != SUCCESS {
		return
	}

	pveTroop := hero.PveTroop(pveTroopType)

	proto := hero.GenCombatPlayerProtoWithCaptains(fullSoldier, pveTroop.Captains(), getGuildSnapshot)
	if proto == nil {
		return nil, SERVER_ERROR
	}

	return proto, SUCCESS
}

func (hero *Hero) GenCombatPlayerProtoWithCaptains(fullSoldier bool, captains []*TroopPos, getGuildSnapshot guildsnapshotdata.Getter) (player *shared_proto.CombatPlayerProto) {
	protos := make([]*shared_proto.CaptainInfoProto, len(captains))
	for i, c := range captains {
		if c == nil {
			continue
		}

		protos[i] = c.EncodeCaptainInfo(fullSoldier)
	}

	player = hero.GenCombatPlayerProtoWithCaptainProtos(protos, getGuildSnapshot)

	return
}

func (hero *Hero) BuildPveCaptainProtosWithHeroAndMonster(fullSoldier bool, pveType shared_proto.PveTroopType, toAdds []*monsterdata.MonsterCaptainData, yuanJun bool) (dungeonTroop []*shared_proto.CaptainInfoProto, failType GenCombatPlayerFailType) {
	failType = hero.checkCanGenCombatPlayer(pveType, uint64(len(toAdds)))
	if failType != SUCCESS {
		return
	}

	pveTroop := hero.PveTroop(pveType)
	// pveTroop.Captains() 数组的长度就是队伍长度，包含空位，不能加上 len(toAdds)
	dungeonTroop = make([]*shared_proto.CaptainInfoProto, len(pveTroop.Captains()))

	// 先放援军，援军放完，再放玩家的部队，玩家部队放完，如果跟援军位置冲突的部队，再找空位放进去
	for _, c := range toAdds {
		if c == nil {
			continue
		}
		index := c.Index
		if index <= 0 || int(index) > len(dungeonTroop) {
			continue
		}

		dungeonTroop[index-1] = c.EncodeCaptainInfo()
	}

	var findEmptyPosCaptains []*TroopPos
	for i, c := range pveTroop.Captains() {
		if c.captain == nil {
			continue
		}

		if dungeonTroop[i] == nil {
			dungeonTroop[i] = c.EncodeCaptainInfo(fullSoldier)
		} else {
			findEmptyPosCaptains = append(findEmptyPosCaptains, c)
		}
	}

	if len(findEmptyPosCaptains) > 0 {

		startIndex := 0
	out:
		for _, c := range findEmptyPosCaptains {

			for i := startIndex; i < len(dungeonTroop); i++ {
				startIndex = i + 1

				if dungeonTroop[i] == nil {
					dungeonTroop[i] = c.EncodeCaptainInfo(fullSoldier)
					continue out
				}
			}

			// 已经没有空位了
			break
		}
	}

	return
}

type TroopFullFightAmountSlice []*Troop

func (p TroopFullFightAmountSlice) Len() int      { return len(p) }
func (p TroopFullFightAmountSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p TroopFullFightAmountSlice) Less(i, j int) bool {
	t1 := p[i]
	t2 := p[j]

	// 战斗力从高到低排序
	return t1.FullFightAmount() > t2.FullFightAmount()
}
