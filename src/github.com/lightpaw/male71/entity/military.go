package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
	"github.com/lightpaw/male7/util/recovtimes"
	"github.com/lightpaw/male7/config/captain"
)

func newMilitary(id int64, initData *heroinit.HeroInitData, buildingEffect *building_effect, ctime time.Time) *Military {
	c := &Military{}
	c.buildingEffect = buildingEffect
	c.captains = make(map[uint64]*Captain)
	c.captainFriendship = newHeroCaptainFriendship()
	c.officialCounter = newCaptainOfficialCounter()
	c.soldierLevelData = initData.FirstLevelSoldierData
	c.soldierRaceData = make(map[shared_proto.Race]*soldier_race_data)
	c.lastRecruitTime = ctime
	c.jiuGuanTimes = newNextTimes()
	c.recruitTimes = NewRecoverableTimes(ctime, initData.JunYingRecoveryDuration, initData.JunYingMaxTimes)
	c.recruitTimes.SetTimes(initData.JunYingDefaultTimes, ctime) // 先给了默认次数了
	c.captainInitData = initData.CaptainInitData

	c.freeSoldier = NewRecoverableTimes(ctime, time.Hour/time.Duration(u64.Max(1, buildingEffect.soldierOutput)), buildingEffect.soldierCapcity)

	c.investigateTroop = newTroop(0, 0, 0)
	c.pveTroops = make([]*PveTroop, len(initData.PveTroopDatas))

	for i, troopData := range initData.PveTroopDatas {
		c.pveTroops[i] = newPveTroop(troopData)
	}

	c.troops = make([]*Troop, initData.FirstLevelHeroData.Sub.TroopsCount)
	for i := uint64(0); i < initData.FirstLevelHeroData.Sub.TroopsCount; i++ {
		c.troops[i] = newTroop(i, id, c.TroopCaptainCount())
	}

	c.candidateCaptainHeads = make([]string, 0)

	return c
}

// 军事
type Military struct {
	buildingEffect *building_effect

	// 武将初始化配置
	captainInitData *heroinit.CaptainInitData

	// 武将相关

	// 新版武将<武将id,Captain>
	captains map[uint64]*Captain

	// 武将羁绊
	captainFriendship *HeroCaptainFriendship

	// 武将官职
	officialCounter   *CaptainOfficialCounter

	troops    []*Troop    // 出征队伍
	investigateTroop *Troop      //侦察队伍
	pveTroops []*PveTroop // pve队伍

	defenserCacheFightAmount uint64 // 缓存的防守战斗力

	soldierLevelData *domestic_data.SoldierLevelData
	soldierRaceData  map[shared_proto.Race]*soldier_race_data

	// 士兵
	woundedSoldider uint64 // 伤兵
	//freeSoldier     uint64            // 空闲士兵
	recruitTimes *recovtimes.RecoverTimes // 招募次数

	// 空闲士兵，自动增加
	freeSoldier         *recovtimes.RecoverTimes
	overflowFreeSoldier uint64 // 溢出的空闲士兵（强征弄出来的），优先消耗这部分士兵

	forceAddSoldierTimes uint64 // 强征次数

	lastRecruitTime time.Time

	nextExpelTime time.Time

	// 修炼馆开始修炼时间
	globalTrainStartTime  time.Time
	captainTrainStartTime time.Time

	jiuGuanTimes *nexttimes // 酒馆次数
	refreshTimes uint64     // 刷新次数
	tutorIndex   uint64     // 导师Index

	candidateCaptainHeads []string // 头像候选

	recervedExp int64 // 修炼馆到期的 buff 给加的、但还没领的武将经验
}

func (c *Military) TrySetOfficialViewed(officialId uint64, officialIdx int32) (exist, success bool) {
	counter := c.officialCounter.Get(officialId)
	if counter != nil {
		exist = true
		_, ok := counter.noViewed[officialIdx]
		if ok {
			delete(counter.noViewed, officialIdx)
			success = true
		}
	}
	return
}

func (c *Military) TryGetOfficialCaptainId(officialId, officialIdx uint64) (captainId uint64, ok bool) {
	if counter := c.officialCounter.Get(officialId); counter != nil && officialIdx < counter.capacity() {
		captainId = counter.get(officialIdx)
		ok = true
	}
	return
}

func (c *Military) TroopCaptainCount() uint64 {
	return c.captainInitData.TroopCaptainCount
}

func (c *Military) SetCandidateCaptainHeads(heads []string) {
	c.candidateCaptainHeads = heads
}

func (c *Military) GetCandidateCaptainHeads() []string {
	return c.candidateCaptainHeads
}

func (c *Military) ResetCandidateCaptainHeads() {
	c.candidateCaptainHeads = make([]string, 0)
}

func (c *Military) resetDaily() {
	c.jiuGuanTimes.SetTimes(0)
	c.forceAddSoldierTimes = 0
}

func (c *Military) GetGlobalTrainStartTime() time.Time {
	return c.globalTrainStartTime
}

func (c *Military) GetCaptainTrainStartTime() time.Time {
	return c.captainTrainStartTime
}

func (c *Military) SetCaptainTrainStartTime(toSet time.Time) {
	c.captainTrainStartTime = toSet
}

func (c *Military) SetGlobalTrainStartTime(toSet time.Time) {
	c.globalTrainStartTime = toSet
	c.SetCaptainTrainStartTime(toSet)
}

func (c *Military) Troops() []*Troop {
	return c.troops
}

func (c *Military) GetOrCreateSoldierData(newRaceData *race.RaceData) soldier_data {
	data := c.soldierRaceData[newRaceData.Race]
	if data == nil {
		data = newSoldierData(newRaceData, c)
		c.soldierRaceData[newRaceData.Race] = data
	}

	return data
}

func newSoldierData(raceData *race.RaceData, m *Military) *soldier_race_data {
	s := &soldier_race_data{}
	s.raceData = raceData
	s.m = m

	return s
}

type soldier_race_data struct {
	raceData *race.RaceData

	m *Military

	levelData          *domestic_data.SoldierLevelData
	extraStat          *data.SpriteStat
	allExtraStat       *data.SpriteStat
	fightTypeExtraStat *data.SpriteStat

	totalStat *data.SpriteStat
}

func (s *soldier_race_data) GetSoldierLevel() uint64 {
	return s.m.SoldierLevel()
}

func (s *soldier_race_data) GetSoldierStat() *data.SpriteStat {

	newLevelData := s.m.soldierLevelData
	extraStat := s.m.buildingEffect.GetSoldierExtraStat(s.raceData.Race)
	allExtraStat := s.m.buildingEffect.GetAllSoldierExtraStat()
	fightTypeExtraStat := s.m.buildingEffect.GetSoldierFightTypeExtraStat(s.raceData.IsFar)
	if s.extraStat == extraStat && s.levelData == newLevelData && s.allExtraStat == allExtraStat && s.fightTypeExtraStat == fightTypeExtraStat {
		return s.totalStat
	}

	s.levelData = newLevelData
	s.extraStat = extraStat
	s.allExtraStat = allExtraStat
	s.fightTypeExtraStat = fightTypeExtraStat
	s.totalStat = data.AppendSpriteStat(newLevelData.GetBaseStat(s.raceData.Race), extraStat, allExtraStat, fightTypeExtraStat)
	return s.totalStat
}

func (s *soldier_race_data) GetSoldierLoad() uint64 {
	return s.m.SoldierLoad()
}

func (s *soldier_race_data) GetSoldierModel() uint64 {
	return s.m.soldierLevelData.GetModel(s.raceData.Race)
}

func (c *Military) GetOfficialCount(officialId uint64) uint64 {
	return c.officialCounter.OfficialPositionCount(officialId)
}

func (c *Military) SetOfficialAnyway(officialId, officialIdx, captainId uint64) {
	c.officialCounter.Set(officialId, officialIdx, captainId)
}

func (c *Military) encodeServer() *server_proto.HeroMilitaryServerProto {
	out := &server_proto.HeroMilitaryServerProto{}

	out.OfficialView = c.officialCounter.encodeServer()
	// 武将羁绊
	for id := range c.captainFriendship.friendship {
		out.CaptainFriendship = append(out.CaptainFriendship, id)
	}

	// 武将相关(新版)
	out.Captains = make([]*server_proto.HeroCaptainServerProto, 0, len(c.captains))
	for _, cp := range c.captains {
		out.Captains = append(out.Captains, cp.encodeServer())
	}

	// 士兵
	out.SoldierLevel = c.soldierLevelData.Level
	//out.FreeSoldier = c.freeSoldier
	out.LastRecruitTime = timeutil.Marshal32(c.lastRecruitTime)
	out.WoundedSoldier = c.woundedSoldider

	out.FreeSoldierStartRecoveryTime = c.freeSoldier.StartTimeUnix64()
	out.OverflowFreeSoldier = c.overflowFreeSoldier
	out.ForceAddSoldierTimes = c.forceAddSoldierTimes

	for _, t := range c.troops {
		if t.isNotEmpty() {
			out.Troops = append(out.Troops, t.encodeServer())
		}
	}

	//增加侦察Troop
	if c.investigateTroop != nil {
		out.InvestigateTroop = c.investigateTroop.encodeServer()
	}

	for _, v := range c.pveTroops {
		out.PveTroops = append(out.PveTroops, v.Encode())
	}

	out.JiuGuanRefreshTimes = c.refreshTimes
	out.JiuGuanTutorIndex = c.tutorIndex
	out.JiuGuanCombineTimes = c.jiuGuanTimes.encodeServer()

	out.JunYingRecruitStartRecoveyTime = c.recruitTimes.StartTimeUnix64()

	out.NextExpelTime = timeutil.Marshal64(c.nextExpelTime)

	out.GlobalTrainStartTime = timeutil.Marshal64(c.globalTrainStartTime)
	out.CaptainTrainStartTime = timeutil.Marshal64(c.captainTrainStartTime)

	out.CandidateCaptainHead = c.candidateCaptainHeads

	out.ReservedExp = c.recervedExp

	return out
}

func (c *Military) unmarshal(hero *Hero, p *server_proto.HeroMilitaryServerProto, datas *config.ConfigDatas, ctime time.Time) {
	if p == nil {
		logrus.Errorf("military unmarshal proto == nil, %d-%s", hero.Id(), hero.Name())
		return
	}

	if hero.LevelData().Sub.TroopsCount > u64.FromInt(len(c.troops)) {
		for i := u64.FromInt(len(c.troops)); i < hero.LevelData().Sub.TroopsCount; i++ {
			c.troops = append(c.troops, newTroop(i, hero.Id(), c.TroopCaptainCount()))
		}
	}

	c.officialCounter.unmarshal(hero.LevelData().Sub.CaptainOfficialIdCount, p.OfficialView)

	c.soldierLevelData = datas.SoldierLevelData().Must(p.SoldierLevel)

	// 兵种属性
	c.buildingEffect.CalculateSoldierStat(hero)

	// 先处理羁绊（下面武将属性计算需要用到）
	for _, id := range p.CaptainFriendship {
		if data := datas.GetCaptainFriendshipData(id); data != nil {
			c.captainFriendship.friendship[id] = data
		}
	}
	c.captainFriendship.calculateProsperity()

	// 武将（新版）
	for _, cp := range p.GetCaptains() {
		if cp == nil {
			continue
		}
		captainData := datas.GetCaptainData(cp.Id)
		if captainData == nil {
			continue
		}

		captain := hero.NewCaptain(captainData, ctime)
		captain.unmarshal(cp, datas, hero.Depot(), hero.Buff(), ctime)
		c.captains[captain.Id()] = captain
	}

	// 内政数据（武将有内政技能，所以放在这里处理）
	c.buildingEffect.CalculateDomestic(hero)

	// 检查羁绊中的武将，是否全部拥有，如果没拥有，则将羁绊移除
	isRemoveFriendship := false
	for id, cf := range c.captainFriendship.friendship {
		for _, data := range cf.Captains {
			if captain := c.captains[data.Id]; captain == nil {
				delete(c.captainFriendship.friendship, id)
				isRemoveFriendship = true
				break
			}
		}
	}

	if isRemoveFriendship {
		// 存在移除的羁绊，重新计算羁绊属性
		c.captainFriendship.calculateProsperity()

		// 重新计算一下所有的武将属性
		for _, captain := range c.captains {
			captain.CalculateProperties()
		}
	}

	//// 武将
	//for _, cp := range p.GetCaptains() {
	//	if cp == nil {
	//		continue
	//	}
	//
	//	raceData := datas.GetRaceData(int(cp.Race))
	//	if raceData == nil {
	//		continue
	//	}
	//
	//	//captain := newCaptain(cp.Id, datas, c.GetOrCreateSoldierData)
	//	//captain.unmarshal(heroId, heroName, depot, heroCaptainSoul, cp, raceData, datas, ctime)
	//	//c.captains[captain.id] = captain
	//}

	//for _, count := range p.GetOfficialCount() {
	//	if count == nil {
	//		continue
	//	}
	//	oid := u64.FromInt32(count.OfficialId)
	//	if c.officialCounter[oid] > 0 {
	//		logrus.Errorf("Military.unmarshal proto.official_count 竟然有重复的 official_id:%v", oid)
	//		// todo 不会出现，暂时不处理
	//	}
	//	c.officialCounter[oid] = u64.FromInt32(count.Count)
	//}

	// 士兵

	//c.freeSoldier = p.GetFreeSoldier()
	c.lastRecruitTime = timeutil.Unix32(p.GetLastRecruitTime())
	c.woundedSoldider = p.GetWoundedSoldier()

	st := ctime
	if p.FreeSoldierStartRecoveryTime > 0 {
		st = timeutil.Unix64(p.FreeSoldierStartRecoveryTime)
	}
	c.freeSoldier = NewRecoverableTimes(st, time.Hour/time.Duration(u64.Max(1, c.buildingEffect.soldierOutput)), c.buildingEffect.soldierCapcity)
	c.overflowFreeSoldier = p.OverflowFreeSoldier
	c.forceAddSoldierTimes = p.ForceAddSoldierTimes

	for _, tp := range p.Troops {
		if int(tp.Sequence) < len(c.troops) {
			c.troops[tp.Sequence].unmarshal(tp, c, hero, ctime)
		}
	}
	if p.InvestigateTroop != nil {
		//反序列化对象
		c.investigateTroop.unmarshal(p.InvestigateTroop, c, hero, ctime)
		c.investigateTroop.SetId(NewHeroTroopId(hero.Id(), datas.RegionConfig().InvestigateTroopId))
	}

	for _, v := range p.PveTroops {
		pveTroop := hero.PveTroop(v.Type)
		if pveTroop == nil {
			continue
		}

		pveTroop.unmarshal(c.Captain, v)

		// 处理下队伍没满的，兼容旧的版本
		if pveTroop.CaptainCount() > 0 {
			continue
		}

		// 没满
		if c.CaptainCount() <= 0 {
			// 没有武将，不管
			continue
		}

		for _, captain := range c.captains {
			if !pveTroop.AddCaptain(captain) {
				// 满了
				break
			}
		}
	}

	c.unmarshalOfficialCounter()

	c.jiuGuanTimes.unmarshal(p.JiuGuanCombineTimes)
	c.refreshTimes = p.JiuGuanRefreshTimes
	c.tutorIndex = p.JiuGuanTutorIndex

	c.recruitTimes.SetStartTime(timeutil.Unix64(p.JunYingRecruitStartRecoveyTime)) // 老号默认给满次数

	c.nextExpelTime = timeutil.Unix64(p.NextExpelTime)

	c.globalTrainStartTime = timeutil.Unix64(p.GlobalTrainStartTime)
	c.captainTrainStartTime = timeutil.Unix64(p.CaptainTrainStartTime)

	c.candidateCaptainHeads = p.CandidateCaptainHead

	c.recervedExp = p.ReservedExp

}

func (m *Military) unmarshalOfficialCounter() {
	for cid, c := range m.Captains() {
		if c.official.Id <= 0 {
			continue
		}
		if counter := m.officialCounter.Get(c.official.Id); counter != nil && counter.checkEmpty(c.officialIndex) {
			counter.set(c.officialIndex, cid)
		} else {
			c.SetOfficial(captain.EmptyOfficialData, 0)
		}
	}
}

func (c *Military) OnHeroLevelChanged(heroId int64, levelData *herodata.HeroLevelData) {

	c.officialCounter.change(levelData.Sub.CaptainOfficialIdCount)

	if levelData.Sub.TroopsCount <= u64.FromInt(len(c.troops)) {
		return
	}

	newTroops := make([]*Troop, levelData.Sub.TroopsCount)
	copy(newTroops, c.troops)

	for i := uint64(len(c.troops)); i < levelData.Sub.TroopsCount; i++ {
		newTroops[i] = newTroop(i, heroId, c.TroopCaptainCount())
	}

	c.troops = newTroops
}

//func (c *Military) EncodeClient(ctime time.Time) *shared_proto.HeroMilitaryProto {
//	out := &shared_proto.HeroMilitaryProto{}
//
//	// 武将
//	out.Captains = make([]*shared_proto.HeroCaptainProto, 0, len(c.captains))
//	for _, captain := range c.captains {
//		out.Captains = append(out.Captains, captain.EncodeClient())
//	}
//
//	out.CaptainGenTime = i64.Int32(c.captainGenTime.Unix64())
//	out.Seeker = c.seeker
//	out.CaptainIndex = u64.Int32Array(c.captainIndex)
//	out.Defenser = u64.Int32Array(c.defenser)
//
//	// 士兵
//	out.SoldierLevel = u64.Int32(c.soldierLevelData.Level)
//	out.SoldierCapcity = u64.Int32(c.soldierCapcity)
//	out.FreeSoldier = int32(c.freeSoldier)
//
//	out.WoundedSoldierCapcity = int32(c.woundedCapcity)
//
//	out.NewSoldier = int32(c.NewSoldierCount(ctime))
//	out.NewSoldierCapcity = int32(c.newSoldierCapcity)
//	out.NewSoldierOutput = int32(c.newSoldierOutput)
//
//	return out
//}

func (hero *Hero) GetMoveSpeedRate() float64 {
	return hero.military.captainFriendship.moveSpeedRate
}

func (c *Military) GetCaptainFriendship() *HeroCaptainFriendship {
	return c.captainFriendship
}

func (c *Military) NextExpelTime() time.Time {
	return c.nextExpelTime
}

func (c *Military) SetNextExpelTime(toSet time.Time) {
	c.nextExpelTime = toSet
}

func (c *Military) SoldierLevel() uint64 {
	return c.soldierLevelData.Level
}

func (c *Military) SoldierLevelData() *domestic_data.SoldierLevelData {
	return c.soldierLevelData
}

func (c *Military) SetSoldierLevelData(toSet *domestic_data.SoldierLevelData) {
	c.soldierLevelData = toSet
}

func (d *Military) GetSoldierExtraStat(race shared_proto.Race) *data.SpriteStat {
	return d.buildingEffect.GetSoldierExtraStat(race)
}

func (d *Military) GetSoldierTotalStat(race shared_proto.Race) *data.SpriteStat {
	return data.AppendSpriteStat(d.soldierLevelData.GetBaseStat(race), d.buildingEffect.GetSoldierExtraStat(race), d.buildingEffect.GetAllSoldierExtraStat())
}

func (c *Military) ExtraLoad() uint64 {
	return c.buildingEffect.extraLoad
}

func (c *Military) SoldierLoad() uint64 {
	return c.soldierLevelData.Load + c.ExtraLoad()
}

func (c *Military) RecruitTimes() *recovtimes.RecoverTimes {
	return c.recruitTimes
}

func (c *Military) WoundedCapcity() uint64 {
	return c.buildingEffect.woundedCapcity
}

func (c *Military) NewSoldierCapcity() uint64 {
	return c.buildingEffect.newSoldierCapcity
}

func (c *Military) GetForceAddSoldierTimes() uint64 {
	return c.forceAddSoldierTimes
}

func (c *Military) IncreseForceAddSoldierTimes() uint64 {
	c.forceAddSoldierTimes++
	return c.forceAddSoldierTimes
}

func (c *Military) FreeSoldier(ctime time.Time) uint64 {
	return c.freeSoldier.Times(ctime) + c.overflowFreeSoldier
}

func (c *Military) UpdateFreeSoldierOutputCapcity(ctime time.Time) {

	if c.freeSoldier.MaxTimes() != c.buildingEffect.soldierCapcity {
		c.freeSoldier.ChangeMaxTimes(c.buildingEffect.soldierCapcity, ctime)
	}

	if rd := time.Hour / time.Duration(u64.Max(1, c.buildingEffect.soldierOutput)); rd != c.freeSoldier.Duration() {
		c.freeSoldier.ChangeDuration(rd, ctime)
	}
}

func (c *Military) SetFreeSoldier(toSet uint64, ctime time.Time) {

	if toSet > c.freeSoldier.MaxTimes() {
		c.freeSoldier.SetTimes(c.freeSoldier.MaxTimes(), ctime)
		c.overflowFreeSoldier = u64.Sub(toSet, c.freeSoldier.MaxTimes())
	} else {
		c.freeSoldier.SetTimes(toSet, ctime)
		c.overflowFreeSoldier = 0
	}
}

func (c *Military) AddFreeSoldier(toAdd uint64, ctime time.Time) {
	// 如果没满，则加时间，满了加溢出士兵
	t := c.freeSoldier.Times(ctime)
	if t+toAdd <= c.freeSoldier.MaxTimes() {
		c.freeSoldier.AddTimes(toAdd, ctime)
	} else {
		if ta := u64.Sub(c.freeSoldier.MaxTimes(), t); ta > 0 {
			c.freeSoldier.AddTimes(ta, ctime)
		}

		overflow := u64.Sub(t+toAdd, c.freeSoldier.MaxTimes())
		c.overflowFreeSoldier += overflow
	}
}

func (c *Military) ReduceFreeSoldier(toReduce uint64, ctime time.Time) {
	// 优先扣溢出的士兵
	if c.overflowFreeSoldier > 0 {
		if c.overflowFreeSoldier >= toReduce {
			c.overflowFreeSoldier = u64.Sub(c.overflowFreeSoldier, toReduce)
			return
		}
		toReduce = u64.Sub(toReduce, c.overflowFreeSoldier)
		c.overflowFreeSoldier = 0
	}

	c.freeSoldier.ReduceTimes(toReduce, ctime)
}

func (c *Military) NewUpdateFreeSoldierMsg() pbutil.Buffer {
	return military.NewS2cUpdateFreeSoldierMsg(
		c.freeSoldier.StartTimeUnix32(),
		u64.Int32(c.buildingEffect.soldierCapcity),
		u64.Int32(c.buildingEffect.soldierOutput),
		u64.Int32(c.overflowFreeSoldier),
	)
}

//func (c *Military) FightSoldier() uint64 {
//	var captainSoldier uint64 = 0
//	for _, ca := range c.captains {
//		captainSoldier += ca.soldier
//	}
//
//	return captainSoldier
//}
//
//// 计算需要补兵的数量
//func (c *Military) NeedRecoverSoldier() uint64 {
//	var needSoldier uint64 = 0
//	for _, ca := range c.captains {
//		if ca.IsOutSide() {
//			continue
//		}
//
//		needSoldier += u64.Sub(ca.SoldierCapcity(), ca.soldier)
//	}
//
//	return needSoldier
//}

func (c *Military) SoldierCapcity() uint64 {
	return c.buildingEffect.soldierCapcity
}

func (c *Military) NewSoldierOutput() uint64 {
	return c.buildingEffect.newSoldierOutput
}

func (c *Military) NewSoldierRecruitCount() uint64 {
	return c.buildingEffect.newSoldierRecruitCount
}

func calNewSoldierCount(output uint64, duration time.Duration) uint64 {
	return u64.Multi(output, duration.Hours())
}

func calNewSoldierDuration(output, newSoldierCount uint64) time.Duration {
	return time.Duration(u64.Division2Float64(newSoldierCount, output)) * time.Hour
}

func (c *Military) SetNewSoldierCount(newSoldierCount uint64, ctime time.Time) {
	if newSoldierCount > 0 {
		c.lastRecruitTime = ctime.Add(-calNewSoldierDuration(c.NewSoldierOutput(), newSoldierCount))
	} else {
		c.lastRecruitTime = ctime
	}
}

func (c *Military) NewSoldierCount(ctime time.Time) uint64 {
	return u64.Min(calNewSoldierCount(c.NewSoldierOutput(), ctime.Sub(c.lastRecruitTime)), c.NewSoldierCapcity())
}

func (c *Military) JiuGuanConsult(nextTime time.Time, toSetTutorIndex uint64) {
	c.jiuGuanTimes.IncreseTimes()        // 加次数
	c.jiuGuanTimes.SetNextTime(nextTime) // 下次时间
	c.refreshTimes = 0
	c.tutorIndex = toSetTutorIndex
}

func (c *Military) JiuGuanTimes() *nexttimes {
	return c.jiuGuanTimes
}

func (c *Military) JiuGuanRefreshTimes() uint64 {
	return c.refreshTimes
}

func (c *Military) IncreJiuGuanRefreshTimes() {
	c.refreshTimes++
}

func (c *Military) JiuGuanTutorIndex() uint64 {
	return c.tutorIndex
}

func (c *Military) SetJiuGuanTutorIndex(toSet uint64) {
	c.tutorIndex = toSet
}

func (c *Military) TrySetCaptainFriendship(data *captain.CaptainFriendshipData) bool {
	return c.captainFriendship.tryAdd(data)
}

// 武将
func (c *Military) Captain(cid uint64) *Captain {
	return c.captains[cid]
}

func (c *Military) Captains() map[uint64]*Captain {
	return c.captains
}

func (c *Military) CaptainCount() int {
	return len(c.captains)
}

// 新版获取武将数量
func (c *Military) GetCaptainCount() int {
	return len(c.captains)
}

func (c *Military) AddCaptain(captain *Captain) {
	c.captains[captain.data.Id] = captain
}

func (c *Military) MaxFightingAmountCaptain() (captain *Captain) {
	for _, c := range c.captains {
		if captain == nil || c.FightAmount() > captain.FightAmount() {
			captain = c
		}
	}
	return
}

//func (c *Military) RemoveCaptain(id uint64) {
//
//	captain := c.captains[id]
//	delete(c.captains, id)
//
//	if captain != nil {
//		// 解绑
//		if t := captain.troop; t != nil {
//			for i, v := range t.captains {
//				if captain == v {
//					t.captains[i] = nil
//					break
//				}
//			}
//		} else {
//			logrus.Error("删除武将时，发现武将引用的队伍居然为nil")
//			// 从每个队伍都看一下，有就删除
//		out:
//			for _, t := range c.troops {
//				for i, v := range t.captains {
//					if captain == v {
//						t.captains[i] = nil
//						break out
//					}
//				}
//			}
//		}
//
//		captain.troop = nil
//	}
//}

func (c *Military) WoundedSoldier() uint64 {
	return c.woundedSoldider
}

func (c *Military) WoundedSoldierCapcity() uint64 {
	return c.buildingEffect.woundedCapcity
}

func (c *Military) AddWoundedSoldier(toAdd uint64) bool {
	if c.woundedSoldider < c.WoundedSoldierCapcity() {
		c.woundedSoldider = u64.Min(c.woundedSoldider+toAdd, c.WoundedSoldierCapcity())
		return true
	}

	return false
}

func (c *Military) ReduceWoundedSoldier(toReduce uint64) {
	c.woundedSoldider = u64.Sub(c.woundedSoldider, toReduce)
}

func (c *Military) encodeMilitary(ctime time.Time) *shared_proto.HeroMilitaryProto {
	out := &shared_proto.HeroMilitaryProto{}

	for id := range c.captainFriendship.friendship {
		out.CaptainFriendship = append(out.CaptainFriendship, u64.Int32(id))
	}

	for _, v := range c.troops {
		out.Troops = append(out.Troops, v.EncodeClient())
	}
	if c.investigateTroop != nil {
		out.InvestigateTroop = c.investigateTroop.EncodeClient()
	}

	for _, v := range c.pveTroops {
		out.PveTroops = append(out.PveTroops, v.Encode())
	}

	// 士兵
	out.SoldierLevel = u64.Int32(c.SoldierLevelData().Level)
	out.SoldierCapcity = u64.Int32(c.SoldierCapcity())
	//out.FreeSoldier = int32(c.FreeSoldier())

	out.WoundedSoldierCapcity = int32(c.WoundedCapcity())

	out.FreeSoldierStartRecoveryTime = c.freeSoldier.StartTimeUnix32()
	out.FreeSoldierCapcity = u64.Int32(c.buildingEffect.soldierCapcity)
	out.FreeSoldierOutput = u64.Int32(c.buildingEffect.soldierOutput)
	out.OverflowFreeSoldier = u64.Int32(c.overflowFreeSoldier)
	out.ForceAddSoldierTimes = u64.Int32(c.forceAddSoldierTimes)

	out.NewSoldier = int32(c.NewSoldierCount(ctime))
	out.NewSoldierCapcity = int32(c.NewSoldierCapcity())
	out.NewSoldierOutput = int32(c.NewSoldierOutput())
	out.NewSoldierRecruitCount = int32(c.NewSoldierRecruitCount())

	// -----黄金分割线---------

	out.WoundedSoldier = u64.Int32(c.woundedSoldider)

	// 新版武将压入
	if captainLen := len(c.captains); captainLen > 0 {
		out.Captain = make([]*shared_proto.CaptainProto, 0, captainLen)
		for _, captain := range c.captains {
			out.Captain = append(out.Captain, captain.EncodeClient())
		}
	}

	for id, official := range c.officialCounter.officials {
		out.OfficialCount = append(out.OfficialCount, &shared_proto.HeroCaptainOfficialCount {
			OfficialId: u64.Int32(id),
			Count: u64.Int32(official.positionCount()),
		})
	}
	out.OfficialView = c.officialCounter.encode()

	out.JiuGuanRefreshTimes = u64.Int32(c.refreshTimes)
	out.JiuGuanTutorIndex = u64.Int32(c.tutorIndex)
	out.JiuGuanTimes, out.JiuGuanNextTime = c.jiuGuanTimes.encodeClient()

	out.JunYingRecruitStartRecoveyTime = c.recruitTimes.StartTimeUnix32()

	out.NextExpelTime = timeutil.Marshal32(c.nextExpelTime)

	out.GlobalTrainStartTime = timeutil.Marshal32(c.globalTrainStartTime)
	out.CaptainTrainStartTime = timeutil.Marshal32(c.captainTrainStartTime)
	out.TrainCoef = i32.MultiF64(1000, c.buildingEffect.trainCoef)
	return out
}

func (m *Military) ReservedExp() int64 {
	return m.recervedExp
}

func (m *Military) ClearReservedExp() {
	m.recervedExp = 0
}

func (m *Military) AddReservedExp(toAdd int64) {
	m.recervedExp += toAdd
}
