package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/util/recovtimes"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/config/heroinit"
)

func newHeroRegion(initData *heroinit.HeroInitData, ctime time.Time) *hero_region {

	hero := &hero_region{
		home: &home{
			base: &base{},
		},
		homeNpcBaseMap:     make(map[int64]*HomeNpcBase),
		favoritePoses:      newFavoritePoses(initData.MaxFavoritePosCount),
		investigationMap:   make(map[int64]*investigation),
		multiLevelNpcMap:   make(map[shared_proto.MultiLevelNpcType]*npc_type_Info),
		multiLevelNpcTimes: recovtimes.NewExtraRecoverTimes(ctime, initData.MultiLevelNpcRecoveryDuration, initData.MultiLevelNpcMaxTimes),
		invaseHeroTimes:    recovtimes.NewExtraRecoverTimes(ctime, initData.InvaseHeroRecoveryDuration, initData.InvaseHeroMaxTimes),
		junTuanNpcTimes:    recovtimes.NewExtraRecoverTimes(ctime, initData.JunTuanNpcRecoveryDuration, initData.JunTuanNpcMaxTimes),
	}

	hero.multiLevelNpcTimes.SetTimes(initData.MultiLevelNpcInitTimes, ctime, 0) // 设置默认次数
	hero.invaseHeroTimes.SetTimes(initData.InvaseHeroInitTimes, ctime, 0)
	hero.junTuanNpcTimes.SetTimes(initData.JunTuanNpcInitTimes, ctime, 0)

	return hero
}

type hero_region struct {
	home *home

	homeDefenseTroopIndex      uint64
	homeDefenseVersion         uint64
	homeTroopDefeatedMailProto *shared_proto.MailProto // 主城防守部队被打败的邮件

	homeDefenseTroopFightAmount uint64

	// 镜像驻防
	copyDefenser *copy_defenser

	// pve野怪
	createHomeNpcBaseIds []uint64               // 已经创建的pve野怪，因为只创建一次
	homeNpcBaseMap       map[int64]*HomeNpcBase // 当前入侵的野怪列表

	favoritePoses *FavoritePoses // 收藏点

	nextInvestigateTime time.Time
	investigationMap    map[int64]*investigation

	multiLevelNpcMap       map[shared_proto.MultiLevelNpcType]*npc_type_Info
	multiLevelNpcPassLevel uint64 // 怪物通关等级
	multiLevelNpcTimes     *recovtimes.ExtraRecoverTimes

	// 出征玩家次数
	invaseHeroTimes *recovtimes.ExtraRecoverTimes

	// 出征军团怪次数
	junTuanNpcTimes *recovtimes.ExtraRecoverTimes
}

func (hero *Hero) TryExpireCopyDefenser(ctime time.Time) bool {
	if hero.copyDefenser != nil {
		if hero.copyDefenser.endTime.Before(ctime) {
			hero.copyDefenser = nil
			return true
		}
	}
	return false
}

func (hero *Hero) NewUpdateCopyDefenserMsg() pbutil.Buffer {
	if hero.copyDefenser != nil {

		var soldier, totalSoldier uint64
		var captainFightAmount []uint64
		for _, c := range hero.copyDefenser.captains {
			if c == nil {
				continue
			}

			soldier += c.soldier
			totalSoldier += c.totalSoldier

			if c.soldier > 0 {
				captainFightAmount = append(captainFightAmount, c.GetFightAmount())
			}
		}
		totalFightAmount := data.TroopFightAmount(captainFightAmount...)
		return military.NewS2cUpdateCopyDefenserMsg(u64.Int32(soldier), u64.Int32(totalSoldier), u64.Int32(totalFightAmount))
	} else {
		return military.REMOVE_COPY_DEFENSER_S2C
	}
}

func (hero *Hero) getCaptainCopyDefenseStat(id uint64) *shared_proto.SpriteStatProto {
	if captain := hero.Military().Captain(id); captain != nil {
		return captain.GetCopyDefenseStat(hero).Encode()
	}
	return nil
}

func (hero *Hero) UpdateCopyDefenser(troop realmface.Troop) {
	if hero.copyDefenser != nil {
		for _, captain := range troop.Captains() {
			if c := hero.copyDefenser.captainMap[captain.Id()]; c != nil {
				c.soldier = u64.FromInt32(captain.Proto().Soldier)
			} else {
				logrus.WithField("captain", captain.Id()).WithField("troop", troop.Id()).Error("Hero.updateCopyDefenser, 竟然没找到troop中的captain")
			}
		}
	}
}

func (hero *Hero) RemoveCopyDefenser() {
	hero.copyDefenser = nil
}

func (hero *Hero) NewCopyDefCaptainInfo(ctime time.Time) []*shared_proto.CaptainInfoProto {

	if hero.copyDefenser == nil {
		return nil
	}

	if hero.copyDefenser.endTime.Before(ctime) {
		// 已过期，清掉
		hero.copyDefenser = nil
		return nil
	}

	if len(hero.copyDefenser.captainMap) <= 0 {
		hero.copyDefenser = nil
		return nil
	}

	captains := make([]*shared_proto.CaptainInfoProto, len(hero.copyDefenser.captains))
	captainCount := 0
	for i, v := range hero.copyDefenser.captains {
		if v == nil {
			continue
		}

		// 没兵的也不出战
		if v.soldier <= 0 || v.totalSoldier <= 0 {
			continue
		}

		c := hero.military.Captain(v.id)
		if c == nil {
			continue
		}

		captains[i] = c.EncodeCaptainInfoWithSoldier(v.totalStat, v.soldier, v.totalSoldier, v.xIndex)
		captainCount++
	}

	if captainCount <= 0 {
		hero.copyDefenser = nil
		return nil
	}

	return captains
}

func (hero *Hero) GetCopyDefenser() *copy_defenser {
	return hero.copyDefenser
}

func (hero *Hero) SetCopyDefenser(troopIndex uint64, captains []*CaptainIdSoldier, endTime time.Time) {
	hero.copyDefenser = newCopyDefenser(troopIndex, captains, endTime)
}

func newCopyDefenser(troopIndex uint64, captains []*CaptainIdSoldier, endTime time.Time) *copy_defenser {
	m := make(map[uint64]*CaptainIdSoldier, len(captains))
	for _, c := range captains {
		if c != nil {
			m[c.id] = c
		}
	}

	return &copy_defenser{
		troopIndex: troopIndex,
		captains:   captains,
		captainMap: m,
		endTime:    endTime,
	}
}

type copy_defenser struct {
	troopIndex uint64
	captains   []*CaptainIdSoldier
	captainMap map[uint64]*CaptainIdSoldier

	// 结束时间
	endTime time.Time
}

func (d *copy_defenser) GetCaptains() []*CaptainIdSoldier {
	return d.captains
}

func NewCaptainIdSoldier(id uint64, xIndex int32, totalStat *shared_proto.SpriteStatProto, soldier, totalSoldier, spellFightAmountCoef uint64) *CaptainIdSoldier {
	return &CaptainIdSoldier{
		id:                   id,
		xIndex:               xIndex,
		totalStat:            totalStat,
		soldier:              soldier,
		totalSoldier:         totalSoldier,
		spellFightAmountCoef: spellFightAmountCoef,
	}
}

type CaptainIdSoldier struct {
	id uint64

	xIndex int32

	totalStat *shared_proto.SpriteStatProto

	soldier uint64

	totalSoldier uint64

	spellFightAmountCoef uint64
}

func (cis *CaptainIdSoldier) GetId() uint64 {
	return cis.id
}

func (cis *CaptainIdSoldier) GetSoldier() uint64 {
	return cis.soldier
}

func (cis *CaptainIdSoldier) GetTotalSoldier() uint64 {
	return cis.totalSoldier
}

func (cis *CaptainIdSoldier) GetFightAmount() uint64 {
	return u64.FromInt32(data.ProtoFightAmount(cis.totalStat, u64.Int32(cis.soldier), u64.Int32(cis.spellFightAmountCoef)))
}

func (d *hero_region) GetInvestigation(targetId int64) *investigation {
	return d.investigationMap[targetId]
}

func (d *hero_region) RemoveInvestigationMail(mailId uint64) {
	if len(d.investigationMap) > 0 {
		for k, v := range d.investigationMap {
			if v.mailId == mailId {
				delete(d.investigationMap, k)
				break
			}
		}
	}
}

func (d *hero_region) AddInvestigation(targetId int64, expireTime time.Time, mailId uint64) {
	inv := &investigation{}
	inv.targetId = targetId
	inv.expireTime = expireTime
	inv.mailId = mailId

	d.investigationMap[targetId] = inv
}

func (d *hero_region) GetInvestigationCount() uint64 {
	return uint64(len(d.investigationMap))
}

func (d *hero_region) ClearExpiredInvestigation(ctime time.Time, removeMin bool) {

	if n := len(d.investigationMap); n > 0 {

		var minKey int64
		var minExpireTime time.Time
		for k, v := range d.investigationMap {
			if ctime.After(v.expireTime) {
				delete(d.investigationMap, k)
			}

			if removeMin && (minKey == 0 || v.expireTime.Before(minExpireTime)) {
				minKey = k
				minExpireTime = v.expireTime
			}
		}

		if removeMin && n == len(d.investigationMap) {
			// 一个都没减少，移除最小的id
			delete(d.investigationMap, minKey)
		}
	}
}

//gogen:hero
type investigation struct {
	_ struct{} `herofield:"-"`

	targetId int64 `desc:"瞭望目标id" shared:",bytes"`

	expireTime time.Time `desc:"瞭望过期时间"`

	mailId uint64 `desc:"瞭望战报id" shared:",bytes"`
}

func (d *investigation) GetExpireTime() time.Time {
	return d.expireTime
}

func (d *investigation) GetMailId() uint64 {
	return d.mailId
}

func (d *investigation) EncodeServer() *server_proto.HeroInvestigationServerProto {
	proto := &server_proto.HeroInvestigationServerProto{}
	proto.TargetId = d.targetId
	proto.ExpireTime = timeutil.Marshal64(d.expireTime)
	proto.MailId = d.mailId

	return proto
}

func (d *investigation) Unmarshal(proto *server_proto.HeroInvestigationServerProto) {
	if proto == nil {
		return
	}
	d.targetId = proto.TargetId
	d.expireTime = timeutil.Unix64(proto.ExpireTime)
	d.mailId = proto.MailId
}

func (d *investigation) GetTargetId() int64 {
	return d.targetId
}

type npc_type_Info struct {
	npcType shared_proto.MultiLevelNpcType

	passLevel uint64 // 通关等级

	hate uint64 // 仇恨值

	revengeLevel uint64 // 怪物攻城等级

	revengeTime time.Time // 怪物攻城时间
}

//func (info *npc_type_Info) GetPassLevel() uint64 {
//	return info.passLevel
//}
//
//func (info *npc_type_Info) SetPassLevel(toSet uint64) {
//	info.passLevel = toSet
//}

func (d *hero_region) GetMultiLevelNpcPassLevel() uint64 {
	return d.multiLevelNpcPassLevel
}

func (d *hero_region) SetMultiLevelNpcPassLevel(toSet uint64) {
	d.multiLevelNpcPassLevel = toSet
}

func (info *npc_type_Info) GetHate() uint64 {
	return info.hate
}

func (info *npc_type_Info) SetHate(toSet uint64) bool {
	if info.hate != toSet {
		info.hate = toSet
		return true
	}
	return false
}

func (info *npc_type_Info) GetRevengeTime() time.Time {
	return info.revengeTime
}

func (info *npc_type_Info) SetRevengeTime(toSet time.Time) {
	info.revengeTime = toSet
}

func (info *npc_type_Info) GetRevengeLevel() uint64 {
	return info.revengeLevel
}

func (info *npc_type_Info) SetRevengeLevel(toSet uint64) {
	info.revengeLevel = toSet
}

func (hero *hero_region) GetNpcTypeInfo(t shared_proto.MultiLevelNpcType) *npc_type_Info {
	return hero.multiLevelNpcMap[t]
}

func (hero *hero_region) GetOrCreateNpcTypeInfo(t shared_proto.MultiLevelNpcType) *npc_type_Info {
	info := hero.multiLevelNpcMap[t]
	if info == nil {
		info = &npc_type_Info{
			npcType: t,
		}
		hero.multiLevelNpcMap[t] = info
	}
	return info
}

func (hero *hero_region) GetMultiLevelNpcTimes() *recovtimes.ExtraRecoverTimes {
	return hero.multiLevelNpcTimes
}

func (hero *hero_region) GetInvaseHeroTimes() *recovtimes.ExtraRecoverTimes {
	return hero.invaseHeroTimes
}

func (hero *hero_region) GetJunTuanNpcTimes() *recovtimes.ExtraRecoverTimes {
	return hero.junTuanNpcTimes
}

func (hero *hero_region) GetNextInvestigateTime() time.Time {
	return hero.nextInvestigateTime
}

func (hero *hero_region) SetNextInvestigateTime(toSet time.Time) {
	hero.nextInvestigateTime = toSet
}

func (hero *hero_region) FavoritePoses() *FavoritePoses {
	return hero.favoritePoses
}

func (hero *hero_region) HasCreateHomeNpcBase(dataId uint64) bool {
	return u64.Contains(hero.createHomeNpcBaseIds, dataId)
}

func (hero *hero_region) CreateHomeNpcBase(data *basedata.HomeNpcBaseData) *HomeNpcBase {

	id := npcid.NewHomeNpcId(data.Id)
	hero.createHomeNpcBaseIds = u64.AddIfAbsent(hero.createHomeNpcBaseIds, data.Id)

	npcBase := newHomeNpcBase(id, data)
	hero.homeNpcBaseMap[npcBase.Id()] = npcBase

	return npcBase
}

func (hero *hero_region) GetHomeNpcBase(id int64) *HomeNpcBase {
	return hero.homeNpcBaseMap[id]
}

func (hero *hero_region) RemoveHomeNpcBase(id int64) {
	delete(hero.homeNpcBaseMap, id)
}

func (hero *hero_region) WalkHomeNpcBase(f func(base *HomeNpcBase) (toBreak bool)) {
	for _, v := range hero.homeNpcBaseMap {
		if f(v) {
			return
		}
	}
}

func newHomeNpcBase(id int64, data *basedata.HomeNpcBaseData) *HomeNpcBase {
	return &HomeNpcBase{
		IdHolder: idbytes.NewIdHolder(id),
		data:     data,
	}
}

// 主城周围的野怪
type HomeNpcBase struct {
	idbytes.IdHolder

	// 野怪模板
	data *basedata.HomeNpcBaseData

	// 剩余繁荣度，0表示野怪满血，否则表示野怪当前剩余的繁荣度
	prosprity uint64
}

func (b *HomeNpcBase) encode() *server_proto.HeroHomeNpcBaseServerProto {
	proto := &server_proto.HeroHomeNpcBaseServerProto{}

	proto.Sequence = npcid.GetNpcIdSequence(b.Id())

	proto.DataId = b.data.Id
	proto.Prosperity = b.prosprity

	return proto
}

func (b *HomeNpcBase) GetData() *basedata.HomeNpcBaseData {
	return b.data
}

func (b *HomeNpcBase) GetProsprity() uint64 {
	return b.prosprity
}

func (b *HomeNpcBase) SetProsprity(toSet uint64) {
	b.prosprity = toSet
}

func (b *HomeNpcBase) RemoveDefenseTroop() {
	b.prosprity = b.data.Data.ProsperityCapcity
}

func newFavoritePoses(maxCount uint64) *FavoritePoses {
	return &FavoritePoses{
		MaxCount:   maxCount,
		posesProto: &shared_proto.FavoritePosesProto{},
	}
}

// 所有收藏的点
type FavoritePoses struct {
	MaxCount   uint64 // 最大收藏数量
	posesProto *shared_proto.FavoritePosesProto
}

// 是否有相同的点
func (f *FavoritePoses) findPos(realmId, x, y int32) int {
	for idx, pos := range f.posesProto.Poses {
		if pos.Id == realmId && pos.X == x && pos.Y == y {
			return idx
		}
	}

	return -1
}

// 新增
func (f *FavoritePoses) Add(realmId, x, y int32) (full, exist bool) {
	if f.PosCount() >= f.MaxCount {
		full = true
		return
	}

	if f.findPos(realmId, x, y) >= 0 {
		exist = true
		return
	}

	f.posesProto.Poses = append(f.posesProto.Poses, &shared_proto.FavoritePosProto{Id: realmId, X: x, Y: y})

	return
}

// 删除
func (f *FavoritePoses) Del(realmId, x, y int32) (suc bool) {
	index := f.findPos(realmId, x, y)
	if index < 0 {
		return
	}

	suc = true

	poses := f.posesProto.Poses

	if uint64(index+1) != f.PosCount() {
		// 不是最后一个
		copy(poses[index:], poses[index+1:])
	}

	f.posesProto.Poses = poses[:f.PosCount()-1]

	return
}

func (f *FavoritePoses) PosCount() uint64 {
	return uint64(len(f.posesProto.Poses))
}

// 遍历
func (f *FavoritePoses) Walk(walkFunc func(pos *shared_proto.FavoritePosProto)) {
	for _, pos := range f.posesProto.Poses {
		walkFunc(pos)
	}
}

func (f *FavoritePoses) encode() *shared_proto.FavoritePosesProto {
	return f.posesProto
}

func (f *FavoritePoses) unmarshal(proto *shared_proto.FavoritePosesProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	if len(proto.Poses) > 0 {
		f.posesProto.Poses = make([]*shared_proto.FavoritePosProto, 0, len(proto.Poses))

		for _, pos := range proto.Poses {
			level, _, regionType := realmface.ParseRealmId(int64(pos.Id))
			levelData := datas.GetRegionData(regdata.RegionDataID(regionType, level))
			if levelData == nil {
				logrus.WithField("level", level).Debugln("没找到玩家收藏的场景")
				continue
			}

			if pos.X < 0 || uint64(pos.X) >= levelData.Block.XLen {
				logrus.WithField("x", pos.X).Debugln("玩家收藏的场景的x越界了")
				continue
			}

			if pos.Y < 0 || uint64(pos.Y) >= levelData.Block.YLen {
				logrus.WithField("y", pos.Y).Debugln("玩家收藏的场景的y越界了")
				continue
			}

			f.Add(pos.Id, pos.X, pos.Y)
		}
	}
}

func (hero *Hero) GetHomeDefenser() *Troop {
	if hero.homeDefenseTroopIndex > 0 {
		return hero.GetTroopByIndex(u64.Sub(hero.homeDefenseTroopIndex, 1))
	}
	return nil
}

func (hero *Hero) calculateHomeDefenserFightAmount() uint64 {
	t := hero.GetHomeDefenser()
	if t != nil {
		tfa := data.NewTroopFightAmount()
		for _, pos := range t.captains {
			c := pos.captain
			if c != nil && c.soldier > 0 {
				totalStat := c.getDefenseStat(hero)
				tfa.Add(totalStat.FightAmount(c.soldier, c.GetSpellFightAmountCoef()))
			}
		}
		return tfa.ToU64()
	}
	return 0
}

func (hero *Hero) GetHomeDefenserFightAmount() uint64 {
	return hero.homeDefenseTroopFightAmount
}

func (hero *Hero) TryUpdateHomeDefenserFightAmount() (bool, uint64) {
	// 所有建筑附加的
	newAmount := hero.calculateHomeDefenserFightAmount()

	if hero.homeDefenseTroopFightAmount != newAmount {
		hero.homeDefenseTroopFightAmount = newAmount
		return true, newAmount
	}

	return false, newAmount
}

func (hero *Hero) IsHomeDefenseCaptain(captain *Captain) bool {

	if hero.homeDefenseTroopIndex <= 0 {
		return false
	}

	return captain.GetTroopSequence() == hero.homeDefenseTroopIndex
}

//func (hero *Hero) NewHomeDefenserBytes() (version uint64, protoBytes []byte) {
//	hero.homeDefenseVersion++
//	version = hero.homeDefenseVersion
//
//	t := hero.GetHomeDefenser()
//	if t != nil && !t.IsOutside() {
//		protoBytes = must.Marshal(t.encodeBaseDefenserProto())
//	}
//
//	return
//}

func (hero *hero_region) GetHomeDefenseTroopIndex() uint64 {
	return hero.homeDefenseTroopIndex
}

func (hero *hero_region) SetHomeDefenseTroopIndex(toSet uint64) {
	hero.homeDefenseTroopIndex = toSet
}

func (hero *hero_region) GetTroopDefeatedMailProto() *shared_proto.MailProto {
	return hero.homeTroopDefeatedMailProto
}

func (hero *hero_region) SetTroopDefeatedMailProto(proto *shared_proto.MailProto) bool {
	if hero.homeTroopDefeatedMailProto != proto {
		hero.homeTroopDefeatedMailProto = proto
		return true
	}

	return false
}

type base struct {
	realmId      int64 // 第几张地图
	baseLevel    uint64
	baseX, baseY int
	prosperity   uint64 // 繁荣度
}

func (b *base) updateBase(base realmface.Base) {
	b.realmId = base.RegionID()
	b.baseLevel = base.GetBaseLevel()
	b.baseX = base.BaseX()
	b.baseY = base.BaseY()
	b.prosperity = base.Prosperity()

	if b.baseLevel <= 0 {
		b.realmId = 0
	}
}

func (b *base) encode() *server_proto.HeroBase {
	return &server_proto.HeroBase{
		BaseRegion: b.realmId,
		BaseLevel:  b.baseLevel,
		BaseX:      imath.Int32(b.baseX),
		BaseY:      imath.Int32(b.baseY),
		Prosperity: b.prosperity,
	}
}

func (b *base) encodeClient() *shared_proto.HeroBaseProto {
	return &shared_proto.HeroBaseProto{
		BaseRegion: i64.Int32(b.realmId),
		BaseLevel:  u64.Int32(b.baseLevel),
		BaseX:      imath.Int32(b.baseX),
		BaseY:      imath.Int32(b.baseY),
		Prosperity: u64.Int32(b.prosperity),
	}
}

func unmarshalHeroBase(pb *server_proto.HeroBase) *base {
	return &base{
		realmId:    pb.BaseRegion,
		baseLevel:  pb.BaseLevel,
		baseX:      int(pb.BaseX),
		baseY:      int(pb.BaseY),
		prosperity: pb.Prosperity,
	}
}

type home struct {
	*base

	// 主城信息
	lostProsperity              uint64 // 今天还可以扣了多少繁荣度
	nextResetLostProsperityTime time.Time
	maxLostProsperity           uint64 // 今天最多可以扣多少繁荣度
	stopLostProsperity          bool

	moveBaseRestoreProsperityBufEndTime time.Time // 迁城恢复繁荣度的buf的结束时间

	// 白旗
	whiteFlagHeroId        int64
	whiteFlagGuildId       int64
	whiteFlagDisappearTime time.Time

	// 免战结束时间
	newHeroMianDisappearTime time.Time
	mianStartTime            time.Time
	mianDisappearTime        time.Time
	mianRebackTime           time.Time
	nextUseMianGoodsTime     time.Time

	historyMaxBaseLevel uint64 // 历史最高主城等级
}

func (hero *Hero) GetMianStartTime() time.Time {
	return hero.home.mianStartTime
}

func (hero *Hero) SetMianRebackTime(t time.Time) {
	hero.home.mianRebackTime = t
}

func (hero *Hero) GetMianRebackTime() time.Time {
	return hero.home.mianRebackTime
}

func (hero *Hero) GetWhiteFlagHeroId() int64 {
	return hero.home.whiteFlagHeroId
}

func (hero *Hero) GetWhiteFlagGuildId() int64 {
	return hero.home.whiteFlagGuildId
}

func (hero *Hero) GetWhiteFlagDisappearTime() time.Time {
	return hero.home.whiteFlagDisappearTime
}

func (hero *Hero) SetWhiteFlag(heroId, guildId int64, disappearTime time.Time) {
	hero.home.whiteFlagHeroId = heroId
	hero.home.whiteFlagGuildId = guildId
	hero.home.whiteFlagDisappearTime = disappearTime
}

func (hero *Hero) GetMianDisappearTime() time.Time {
	return hero.home.mianDisappearTime
}

func (hero *Hero) SetMianDisappearTime(startTime, disappearTime time.Time) {
	hero.home.mianStartTime = startTime
	hero.home.mianDisappearTime = disappearTime
}

func (hero *Hero) GetNextUseMianGoodsTime() time.Time {
	return hero.home.nextUseMianGoodsTime
}

func (hero *Hero) SetNextUseMianGoodsTime(disappearTime time.Time) {
	hero.home.nextUseMianGoodsTime = disappearTime
}

func (hero *Hero) GetNewHeroMianDisappearTime() time.Time {
	return hero.home.newHeroMianDisappearTime
}

func (hero *Hero) SetNewHeroMianDisappearTime(disappearTime time.Time) {
	hero.home.newHeroMianDisappearTime = disappearTime
}

func (hero *Hero) GetMoveBaseRestoreProsperityBufEndTime() time.Time {
	return hero.home.moveBaseRestoreProsperityBufEndTime
}

func (hero *Hero) SetMoveBaseRestoreProsperityBufEndTime(endTime time.Time) {
	hero.home.moveBaseRestoreProsperityBufEndTime = endTime
}

func (hero *Hero) RemoveWhiteFlag() {
	hero.home.whiteFlagHeroId = 0
	hero.home.whiteFlagGuildId = 0
	hero.home.whiteFlagDisappearTime = time.Time{}
}

func (hero *Hero) UpdateBase(base realmface.Base) {
	switch base.BaseType() {
	case realmface.BaseTypeHome:
		home, ok := base.(realmface.Home)
		if !ok {
			hero.home.updateBase(base)
		} else {
			hero.UpdateHome(home)
		}
	default:
		logrus.WithField("baseType", base.BaseType()).Error("hero.UpdateBase unkown BaseType")
	}
}

// copy
func (hero *Hero) UpdateHome(home realmface.Home) {
	hero.home.updateBase(home)
}

func (h *home) encode() *server_proto.HeroHome {
	return &server_proto.HeroHome{
		Base:                        h.base.encode(),
		LostProsperity:              h.lostProsperity,
		StopLostProsperity:          h.stopLostProsperity,
		NextResetLostProsperityTime: timeutil.Marshal32(h.nextResetLostProsperityTime),

		MaxBaseLevel: h.historyMaxBaseLevel,

		WhiteFlagHeroId:        h.whiteFlagHeroId,
		WhiteFlagGuildId:       h.whiteFlagGuildId,
		WhiteFlagDisappearTime: timeutil.Marshal64(h.whiteFlagDisappearTime),

		MianStartTime:                       timeutil.Marshal64(h.mianStartTime),
		MianDisappearTime:                   timeutil.Marshal64(h.mianDisappearTime),
		NextUseMianGoodsTime:                timeutil.Marshal64(h.nextUseMianGoodsTime),
		NewHeroMianDisappearTime:            timeutil.Marshal64(h.newHeroMianDisappearTime),
		MoveBaseRestoreProsperityBufEndTime: timeutil.Marshal64(h.moveBaseRestoreProsperityBufEndTime),
		MianRebackTime:                      timeutil.Marshal64(h.mianRebackTime),
	}
}

func (result *hero_region) unmarshal(hero *Hero, pb *server_proto.HeroRegionServerProto, datas *config.ConfigDatas, ctime time.Time) {
	if pb == nil {
		return
	}

	if h := pb.Home; h != nil {
		if base := h.Base; base != nil {
			result.home.base = unmarshalHeroBase(base)
			result.home.lostProsperity = h.LostProsperity
			result.home.stopLostProsperity = h.StopLostProsperity
			result.home.nextResetLostProsperityTime = timeutil.Unix32(h.NextResetLostProsperityTime)
			result.home.whiteFlagHeroId = h.WhiteFlagHeroId
			result.home.whiteFlagGuildId = h.WhiteFlagGuildId
			result.home.whiteFlagDisappearTime = timeutil.Unix64(h.WhiteFlagDisappearTime)
			result.home.mianStartTime = timeutil.Unix64(h.MianStartTime)
			result.home.mianDisappearTime = timeutil.Unix64(h.MianDisappearTime)
			result.home.nextUseMianGoodsTime = timeutil.Unix64(h.NextUseMianGoodsTime)
			result.home.newHeroMianDisappearTime = timeutil.Unix64(h.NewHeroMianDisappearTime)
			result.home.historyMaxBaseLevel = h.MaxBaseLevel
			result.home.moveBaseRestoreProsperityBufEndTime = timeutil.Unix64(h.MoveBaseRestoreProsperityBufEndTime)
			result.home.mianRebackTime = timeutil.Unix64(h.MianRebackTime)
		}
	}

	result.homeDefenseTroopIndex = pb.HomeDefenseTroopIndex
	result.homeTroopDefeatedMailProto = pb.HomeTroopDefeatedMailProto

	// 野怪Npc
	result.createHomeNpcBaseIds = u64.Copy(pb.CreateHomeNpcBaseIds)
	for _, v := range pb.HomeNpcBase {
		data := datas.GetHomeNpcBaseData(v.DataId)
		if data == nil {
			continue
		}

		npcBase := newHomeNpcBase(npcid.NewHomeNpcId(data.Id), data)
		npcBase.prosprity = v.Prosperity

		result.homeNpcBaseMap[npcBase.Id()] = npcBase
	}

	result.favoritePoses.unmarshal(pb.FavoritePoses, datas)

	result.nextInvestigateTime = timeutil.Unix64(pb.NextInvestigateTime)
	for _, v := range pb.Investigation {
		result.AddInvestigation(v.TargetId, timeutil.Unix64(v.ExpireTime), v.MailId)
	}

	for _, v := range pb.MultiLevelNpc {
		info := result.GetOrCreateNpcTypeInfo(v.Type)
		info.passLevel = v.PassLevel
		info.hate = v.Hate
		info.revengeLevel = v.RevengeLevel
	}

	hero.multiLevelNpcPassLevel = pb.MultiLevelNpcPassLevel
	hero.multiLevelNpcTimes.SetStartTime(timeutil.Unix64(pb.MultiLevelNpcStartTime))

	hero.invaseHeroTimes.SetStartTime(timeutil.Unix64(pb.InvaseHeroStartTime))

	hero.junTuanNpcTimes.SetStartTime(timeutil.Unix64(pb.JunTuanNpcStartTime))

	if endTime := timeutil.Unix64(pb.CopyDefenserEndTime); ctime.Before(endTime) {
		n := imath.Minx(len(pb.CopyDefenserCaptainId),
			len(pb.CopyDefenserCaptainSoldier),
			len(pb.CopyDefenserCaptainTotalSoldier),
			len(pb.CopyDefenserCaptainIndex),
			len(pb.CopyDefenserCaptainXIndex),
			len(pb.CopyDefenserCaptainStat),
			len(pb.CopyDefenserCaptainSpellFightAmountCoef))
		if n > 0 {

			max := u64.Min(u64.Maxa(pb.CopyDefenserCaptainIndex)+1, constants.CaptainCountPerTroop)
			captains := make([]*CaptainIdSoldier, max)
			captainCount := 0
			for i := 0; i < n; i++ {
				idx := pb.CopyDefenserCaptainIndex[i]
				if idx < max {
					captains[idx] = NewCaptainIdSoldier(
						pb.CopyDefenserCaptainId[i],
						pb.CopyDefenserCaptainXIndex[i],
						pb.CopyDefenserCaptainStat[i],
						pb.CopyDefenserCaptainSoldier[i],
						pb.CopyDefenserCaptainTotalSoldier[i],
						pb.CopyDefenserCaptainSpellFightAmountCoef[i],
					)
					captainCount++
				}
			}

			if captainCount > 0 {
				troopIndex := u64.Max(pb.CopyDefenserTroopIndex, 1)
				hero.copyDefenser = newCopyDefenser(troopIndex, captains, endTime)
			}
		}
	}

	hero.TryUpdateHomeDefenserFightAmount()
}

func (d *hero_region) encodeRegion() *server_proto.HeroRegionServerProto {
	result := &server_proto.HeroRegionServerProto{}

	result.Home = d.home.encode()

	result.HomeDefenseTroopIndex = d.homeDefenseTroopIndex
	result.HomeTroopDefeatedMailProto = d.homeTroopDefeatedMailProto

	result.CreateHomeNpcBaseIds = d.createHomeNpcBaseIds

	for _, v := range d.homeNpcBaseMap {
		result.HomeNpcBase = append(result.HomeNpcBase, v.encode())
	}

	result.FavoritePoses = d.favoritePoses.encode()

	result.NextInvestigateTime = timeutil.Marshal64(d.nextInvestigateTime)
	for _, v := range d.investigationMap {
		result.Investigation = append(result.Investigation, v.EncodeServer())
	}

	for _, v := range d.multiLevelNpcMap {
		result.MultiLevelNpc = append(result.MultiLevelNpc, &server_proto.HeroMultiLevelNpcServerProto{
			Type:         v.npcType,
			PassLevel:    v.passLevel,
			Hate:         v.hate,
			RevengeLevel: v.revengeLevel,
		})
	}
	result.MultiLevelNpcPassLevel = d.multiLevelNpcPassLevel
	result.MultiLevelNpcStartTime = d.multiLevelNpcTimes.StartTimeUnix64()

	result.InvaseHeroStartTime = d.invaseHeroTimes.StartTimeUnix64()

	result.JunTuanNpcStartTime = d.junTuanNpcTimes.StartTimeUnix64()

	if d.copyDefenser != nil {
		for i, c := range d.copyDefenser.captains {
			if c == nil {
				continue
			}
			result.CopyDefenserCaptainId = append(result.CopyDefenserCaptainId, c.id)
			result.CopyDefenserCaptainSoldier = append(result.CopyDefenserCaptainSoldier, c.soldier)
			result.CopyDefenserCaptainTotalSoldier = append(result.CopyDefenserCaptainTotalSoldier, c.totalSoldier)
			result.CopyDefenserCaptainIndex = append(result.CopyDefenserCaptainIndex, uint64(i))
			result.CopyDefenserCaptainXIndex = append(result.CopyDefenserCaptainXIndex, c.xIndex)
			result.CopyDefenserCaptainStat = append(result.CopyDefenserCaptainStat, c.totalStat)
			result.CopyDefenserCaptainSpellFightAmountCoef = append(result.CopyDefenserCaptainSpellFightAmountCoef, c.spellFightAmountCoef)
		}

		result.CopyDefenserEndTime = timeutil.Marshal64(d.copyDefenser.endTime)
		result.CopyDefenserTroopIndex = d.copyDefenser.troopIndex
	}

	return result
}

func (hero *Hero) encodeRegionClient(ctime time.Time, getGuildSnapshot guildsnapshotdata.Getter) *shared_proto.HeroRegionProto {
	proto := &shared_proto.HeroRegionProto{}

	// 主城
	proto.Home = hero.home.encodeClient()
	proto.MaxBaseLevel = u64.Int32(hero.home.historyMaxBaseLevel)
	proto.MoveBaseRestoreProsperityBufEndTime = timeutil.Marshal32(hero.home.moveBaseRestoreProsperityBufEndTime)

	if hero.home.whiteFlagGuildId != 0 {
		if ctime.Before(hero.home.whiteFlagDisappearTime) {
			g := getGuildSnapshot(hero.home.whiteFlagGuildId)
			if g != nil {
				proto.WhiteFlagGuildId = i64.Int32(hero.home.whiteFlagGuildId)
				proto.WhiteFlagGuildFlagName = g.FlagName
				proto.WhiteFlagDisappearTime = timeutil.Marshal32(hero.home.whiteFlagDisappearTime)
			} else {
				// 帮派没了，拔旗
				hero.RemoveWhiteFlag()
			}
		} else {
			// 白旗过期
			hero.RemoveWhiteFlag()
		}
	}
	proto.MianStartTime = timeutil.Marshal32(hero.home.mianStartTime)
	proto.MianDisappearTime = timeutil.Marshal32(hero.home.mianDisappearTime)
	proto.NextUseMianGoodsTime = timeutil.Marshal32(hero.home.nextUseMianGoodsTime)

	// 行营
	proto.HomeDefenseTroopIndex = u64.Int32(hero.homeDefenseTroopIndex)
	proto.FavoritePoses = hero.favoritePoses.encode()

	proto.HomeTroopDefeatedMailProto = hero.homeTroopDefeatedMailProto

	// 驻守镜像
	if hero.copyDefenser != nil {
		proto.CopyDefenserTroopIndex = u64.Int32(hero.copyDefenser.troopIndex)
		proto.CopyDefenserCaptainId = make([]int32, len(hero.copyDefenser.captains))
		proto.CopyDefenserCaptainXindex = make([]int32, len(hero.copyDefenser.captains))
		for i, c := range hero.copyDefenser.captains {
			if c == nil {
				continue
			}
			proto.CopyDefenserCaptainId[i] = u64.Int32(c.id)
			proto.CopyDefenserCaptainXindex[i] = c.xIndex
			proto.CopyDefenserSoldier += u64.Int32(c.soldier)
			proto.CopyDefenserTotalSoldier += u64.Int32(c.totalSoldier)
		}

		proto.CopyDefenserEndTime = timeutil.Marshal32(hero.copyDefenser.endTime)
	}

	proto.NextInvestigateTime = timeutil.Marshal32(hero.nextInvestigateTime)

	for _, v := range hero.multiLevelNpcMap {
		proto.MultiLevelNpc = append(proto.MultiLevelNpc, &shared_proto.HeroMultiLevelNpcProto{
			Type:      v.npcType,
			PassLevel: u64.Int32(v.passLevel),
			Hate:      u64.Int32(v.hate),
		})
	}
	proto.MultiLevelNpcPassLevel = u64.Int32(hero.multiLevelNpcPassLevel)
	proto.MultiLevelNpcStartRecoveryTime = &shared_proto.RecoverableTimesWithExtraTimesProto{
		StartTime: hero.multiLevelNpcTimes.StartTimeUnix32(),
	}

	proto.InvaseHeroStartRecoveryTime = &shared_proto.RecoverableTimesWithExtraTimesProto{
		StartTime: hero.invaseHeroTimes.StartTimeUnix32(),
	}

	proto.JunTuanNpcStartTime = &shared_proto.RecoverableTimesWithExtraTimesProto{
		StartTime: hero.junTuanNpcTimes.StartTimeUnix32(),
	}

	return proto
}

func (hero *Hero) BaseRegion() int64 {
	return hero.home.realmId
}

func (hero *Hero) BaseLevel() uint64 {
	return hero.home.baseLevel
}

func (hero *Hero) HomeHistoryMaxLevel() uint64 {
	return hero.home.historyMaxBaseLevel
}

func (hero *Hero) BaseX() int {
	return hero.home.baseX
}

func (hero *Hero) BaseY() int {
	return hero.home.baseY
}

func (hero *Hero) SetBaseLevel(baseLevel uint64) {
	hero.home.baseLevel = baseLevel
	hero.home.historyMaxBaseLevel = u64.Max(hero.home.historyMaxBaseLevel, baseLevel)

	if baseLevel <= 0 {
		hero.home.realmId = 0
	}
}

func (hero *Hero) Prosperity() uint64 {
	return hero.home.prosperity
}

func (hero *Hero) SetProsperity(amount uint64) {
	hero.home.prosperity = u64.Min(amount, hero.ProsperityCapcity())
}

func (hero *Hero) AddProsperityCapcity(toAdd uint64) {
	hero.buildingEffect.prosperityCapcity += toAdd
	hero.home.maxLostProsperity = 0
}

func (hero *Hero) AddLostProsperity(toAdd, maxLostProsperity uint64) bool {
	hero.home.lostProsperity += toAdd
	if hero.home.lostProsperity >= maxLostProsperity {
		hero.SetStopLostProsperity()
	}
	return hero.home.stopLostProsperity
}

func (hero *Hero) SetStopLostProsperity() {
	hero.home.stopLostProsperity = true
}

func (hero *Hero) LostProsperity(ctime time.Time, dailyResetTime *timeutil.CycleTime) (uint64, bool) {
	if hero.home.nextResetLostProsperityTime.Before(ctime) {
		hero.home.lostProsperity = 0
		hero.home.stopLostProsperity = false
		hero.home.nextResetLostProsperityTime = dailyResetTime.NextTime(ctime)
	}

	return hero.home.lostProsperity, hero.home.stopLostProsperity
}

func (hero *Hero) MaxLostProsperity(f func(prosperityCapcity uint64) uint64) uint64 {
	if hero.home.maxLostProsperity == 0 {
		hero.home.maxLostProsperity = f(hero.ProsperityCapcity())
	}

	return hero.home.maxLostProsperity
}

func (hero *Hero) GetStopLostProsperity() bool {
	return hero.home.stopLostProsperity
}

func (hero *Hero) ProsperityCapcity() uint64 {
	return hero.buildingEffect.prosperityCapcity
}

func (hero *Hero) ClearBase() {
	hero.home.realmId = 0
	hero.home.baseX, hero.home.baseY = 0, 0
}

func (hero *Hero) SetBaseXY(region int64, x, y int) {
	hero.home.realmId = region
	hero.home.baseX, hero.home.baseY = x, y
}
