package realm

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/config/xiongnu"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"runtime/debug"
	"sync"
	"time"
)

var (
	_ realmface.Base = (*heroHome)(nil)
	_ realmface.Base = (*homeNpcBase)(nil)
	_ realmface.Base = (*basicNpcBase)(nil)
	_ realmface.Base = (*multiLevelNpcBase)(nil)
	_ realmface.Base = (*xiongNuNpcBase)(nil)
	_ realmface.Base = (*baozNpcBase)(nil)
	_ realmface.Base = (*junTuanNpcBase)(nil)
	_ realmface.Base = (*guildWorkshopBase)(nil)

	_ realmface.Home = (*heroHome)(nil)

	_ internalBase = (*heroHome)(nil)
	_ internalBase = (*homeNpcBase)(nil)
	_ internalBase = (*basicNpcBase)(nil)
	_ internalBase = (*multiLevelNpcBase)(nil)
	_ internalBase = (*xiongNuNpcBase)(nil)
	_ internalBase = (*baozNpcBase)(nil)
	_ internalBase = (*junTuanNpcBase)(nil)
	_ internalBase = (*guildWorkshopBase)(nil)
)

// 内部使用
type internalBase interface {
	realmface.Base
	SetBaseLevel(newLevel uint64)
	SetBaseXY(x, y int)
	NewUpdateBaseInfoMsg(r *Realm, addType int32, base *baseWithData) pbutil.Buffer // 不缓存
	GetCantSeeMeMsg() pbutil.Buffer
	ReduceProsperity(uint64) uint64
	ReduceProsperityDontKill(uint64) uint64
	AddProsperity(uint64) uint64
	SetProsperityCapcity(uint64)
	TrySetStopLostProsperity(bool) bool

	//ChangeGuild(id int64, name, flagName string)
	UpdateHeroBasicInfo(hero *entity.Hero)

	MianDisappearTime() int32
	SetMianDisappearTime(disappearTime int32)

	getBaseInfoByLevel(level uint64) baseInfo
	getNpcConfig() (dataId uint64, npcType npcid.NpcType)
	encodeNpcProto(self *baseWithData, r *Realm) *server_proto.NpcBaseProto
	ClearUpdateBaseInfoMsg()
}

//type npcBase interface {
//	getData() *basedata.NpcBaseData
//	newUpdateNpcBaseMsg(rid int64) pbutil.Buffer
//}

type baseWithData struct {
	internalBase internalBase

	// 如果部队正在返回状态, 则只在selfTroops里.
	// 即在这里的troop, troop里的targetBase也是这个base
	targetingTroops map[int64]*troop // 以这个基地为目标的部队, 或别人驻扎或掠夺或前往. id为troop的id. 如果是军营带着部队, 则默认带着的部队是驻扎的状态.

	// 在这里的troop, troop的startingBase也是这个
	selfTroops              map[int64]*troop // 出发地是这个基地的在外的部队. 包括自己驻扎在自己里面的部队(targetBase为nil, 不然是在帮别人守). 如果被打爆了, 里面的部队都消失
	selfNpcTroopSequenceGen uint64

	/*
		如果要迁移, 则需要看selfTroops, 里面只能有Defending状态的部队

		如果被攻击, 则2个map都需要看Defending状态的部队

		如果要驻扎, 则2个map都需要看Defending状态的部队是否超过上限

		如果要掠夺, 看targetingTroop里状态为robbing的是否超过上限

		如果来帮忙, 看targetingTroop里状态为robbing的, 赶走. 然后作为第一个defending的人
	*/

	//kimi 观察者列表
	watchList *sync.Map

	roBase atomic.Value
}

func (b *baseWithData) updateXiongNuTarget() {
	if base := GetXiongNuBase(b); base != nil {

		invateTargetCountMap := make(map[int64]uint64)
		for _, t := range b.selfTroops {
			if t.State().IsInvateState() && t.targetBase != nil && t.targetBase.isHeroHomeBase() {
				invateTargetCountMap[t.targetBase.Id()]++
			}
		}

		base.setInvadeTargetCount(invateTargetCountMap)
	}
}

func (b *baseWithData) updateRoBase() {

	proto := &server_proto.RoBaseProto{}
	proto.Id = b.Id()
	proto.Level = b.BaseLevel()
	proto.X = int32(b.BaseX())
	proto.Y = int32(b.BaseY())
	proto.Prosperity = b.Prosperity()
	proto.ProsperityCapcity = b.ProsperityCapcity()
	proto.BaseType = int32(b.BaseType())

	proto.Name = b.internalBase.HeroName()
	proto.GuildId = b.internalBase.GuildId()

	// Npc城池
	proto.NpcDataId, proto.NpcType = b.internalBase.getNpcConfig()

	// 宝藏Npc专用
	if bz := GetBaoZangBase(b); bz != nil {
		bz.updateRoBase(proto)
	}

	// 军团Npc专用
	if bz := GetJunTuanBase(b); bz != nil {
		bz.updateRoBase(b, proto)
	}

	b.roBase.Store(proto)
}

func (b *baseWithData) getRoBase() *server_proto.RoBaseProto {
	return b.roBase.Load().(*server_proto.RoBaseProto)
}

func (b *baseWithData) nextNpcTroopId() int64 {

	for i := 0; i < 10; i++ {
		// 尝试10次，找不到就算了
		b.selfNpcTroopSequenceGen++
		newTroopId := npcid.NewNpcTroopId(b.Id(), b.selfNpcTroopSequenceGen%npcid.TroopMaxSequence)

		if _, exist := b.selfTroops[newTroopId]; !exist {
			return newTroopId
		}
	}

	return 0
}

func (b *baseWithData) Id() int64 {
	return b.internalBase.Id()
}

func (b *baseWithData) IdBytes() []byte {
	return b.internalBase.IdBytes()
}

func (b *baseWithData) TroopType() shared_proto.TroopType {
	return b.internalBase.TroopType()
}

func (b *baseWithData) Base() realmface.Base {
	return b.internalBase
}

func (b *baseWithData) BaseX() int {
	return b.internalBase.BaseX()
}

func (b *baseWithData) BaseY() int {
	return b.internalBase.BaseY()
}

func (b *baseWithData) CanSee(area *realmface.ViewArea) bool {
	if area == nil {
		return false
	}
	return area.CanSeePos(b.BaseX(), b.BaseY())
}

//GetPos 获取坐标
func (b *baseWithData) GetPos() [][2]int {
	return [][2]int{[2]int{b.BaseX(), b.BaseY()}}
}

//AddWatcher 添加观察者
func (b *baseWithData) AddWatcher(hc iface.HeroController) {
	b.watchList.Store(hc.Id(), hc)
}

//RemoveWatcher 移除观察者
func (b *baseWithData) RemoveWatcher(hc iface.HeroController) {
	b.watchList.Delete(hc.Id())
}

//RangeWatchList 循环观察列表
func (b *baseWithData) RangeWatchList(f func(hc iface.HeroController)) {
	b.watchList.Range(func(k, v interface{}) bool {
		hc, ok := v.(iface.HeroController)
		if ok {
			f(hc)
		}
		return true
	})
}

func (b *baseWithData) SetBaseXY(x, y int) {
	b.internalBase.SetBaseXY(x, y)
}

func (b *baseWithData) BaseType() realmface.BaseType {
	return b.internalBase.BaseType()
}

func (b *baseWithData) BaseLevel() uint64 {
	return b.internalBase.GetBaseLevel()
}

func (b *baseWithData) SetBaseLevel(newLevel uint64) {
	b.internalBase.SetBaseLevel(newLevel)
}

func (b *baseWithData) RegionID() int64 {
	return b.internalBase.RegionID()
}

func (b *baseWithData) GuildId() int64 {
	return b.internalBase.GuildId()
}

func (b *baseWithData) Prosperity() uint64 {
	return b.internalBase.Prosperity()
}

func (b *baseWithData) MianDisappearTime() int32 {
	return b.internalBase.MianDisappearTime()
}

func (b *baseWithData) SetMianDisappearTime(disappearTime int32) {
	b.internalBase.SetMianDisappearTime(disappearTime)
}

func (b *baseWithData) ProsperityCapcity() uint64 {
	return b.internalBase.ProsperityCapcity()
}

func (b *baseWithData) AddProsperity(toAdd uint64) uint64 {
	return b.internalBase.AddProsperity(toAdd)
}

func (b *baseWithData) ReduceProsperity(toReduce uint64) uint64 {
	return b.internalBase.ReduceProsperity(toReduce)
}

func (b *baseWithData) ReduceProsperityDontKill(toReduce uint64) uint64 {
	return b.internalBase.ReduceProsperityDontKill(toReduce)
}

func (b *baseWithData) SetProsperityCapcity(toSet uint64) {
	b.internalBase.SetProsperityCapcity(toSet)
}

func (b *baseWithData) TrySetStopLostProsperity(toSet bool) bool {
	return b.internalBase.TrySetStopLostProsperity(toSet)
}

func (b *baseWithData) NewUpdateBaseInfoMsg(r *Realm, addType int32) pbutil.Buffer {
	return b.internalBase.NewUpdateBaseInfoMsg(r, addType, b)
}

func (b *baseWithData) GetCantSeeMeMsg() pbutil.Buffer {
	return b.internalBase.GetCantSeeMeMsg()
}

func (b *baseWithData) ClearUpdateBaseInfoMsg() {
	b.internalBase.ClearUpdateBaseInfoMsg()
}

func (base *baseWithData) toReportHeroProto(r *Realm, level uint64) *shared_proto.ReportHeroProto {
	hero := base.internalBase.getBaseInfoByLevel(level).GetHeroBasicProto(r.services.heroSnapshotService.Get)
	return base.newReportHeroProto(r, hero)
}

func (base *baseWithData) newReportHeroProto(r *Realm, hero *shared_proto.HeroBasicProto) *shared_proto.ReportHeroProto {

	proto := &shared_proto.ReportHeroProto{
		Id:                base.IdBytes(),
		Name:              hero.Name,
		Level:             int32(hero.Level),
		Head:              hero.GetHead(),
		BaseRegion:        i64.Int32(base.RegionID()),
		BaseX:             imath.Int32(base.BaseX()),
		BaseY:             imath.Int32(base.BaseY()),
		Prosperity:        u64.Int32(base.Prosperity()),
		ProsperityCapcity: u64.Int32(base.ProsperityCapcity()),
		Country:           hero.CountryId,
	}

	if base.GuildId() != 0 {
		g := r.services.guildService.GetSnapshot(base.GuildId())
		if g != nil {
			proto.GuildFlagName = g.FlagName
		}
	}

	return proto
}

func (b *baseWithData) GetBeenAttackOrRobCount() (beenAttackCount, beenRobCount int64) {
	if len(b.targetingTroops) > 0 {
		for _, troop := range b.targetingTroops {
			switch troop.State() {
			case realmface.MovingToInvade, realmface.Assembly:
				beenAttackCount++
			case realmface.Robbing:
				beenRobCount++
			}
		}
	}

	return
}

func (b *baseWithData) remindAttackOrRobCountChanged(r *Realm) {
	if b.isNpcBase() {
		return
	}

	beenAttackCount, beenRobCount := b.GetBeenAttackOrRobCount()
	r.services.reminderService.ChangeAttackOrRobCount(b.Id(), beenAttackCount, beenRobCount, b.GuildId(), b.isHeroHomeBase())
}

func GetHomeBase(base *baseWithData) *heroHome {
	if heroHome, ok := base.internalBase.(*heroHome); ok {
		return heroHome
	}
	return nil
}

func newBaseWithData(basic internalBase) *baseWithData {
	b := &baseWithData{
		internalBase:    basic,
		targetingTroops: make(map[int64]*troop),
		selfTroops:      make(map[int64]*troop),
		watchList:       &sync.Map{},
	}

	b.updateRoBase()
	return b
}

func (r *Realm) newBase(hero *entity.Hero) *baseWithData {
	basic := &heroHome{
		basicBase: basicBase{
			realmId:  r.id,
			id:       hero.Id(),
			idBytes:  hero.IdBytes(),
			heroName: hero.Name(),
			vipLevel: hero.VipLevel(),

			baseX:              hero.BaseX(),
			baseY:              hero.BaseY(),
			baseLevel:          hero.BaseLevel(),
			prosperity:         hero.Prosperity(),
			prosperityCapcity:  hero.ProsperityCapcity(),
			stopLostProsperity: hero.GetStopLostProsperity(),
			guildId:            hero.GuildId(),
		},

		homeNpcBase: make(map[int64]*baseWithData),

		isRestore: hero.MiscData().GetIsRestoreProsperity(),

		whiteFlagGuildId:       hero.GetWhiteFlagGuildId(),
		whiteFlagDisappearTime: timeutil.Marshal32(hero.GetWhiteFlagDisappearTime()),

		mianDisappearTime:  timeutil.Marshal32(hero.GetMianDisappearTime()),
		outerCityUnlockBit: hero.Domestic().OuterCities().UnlockBit(),

		keepBaozMap: make(map[int64]int32),
	}

	basic.UpdateMoveBaseBuf(hero)

	return newBaseWithData(basic)
}

func (r *Realm) newAllHomeNpcBase(hero *entity.Hero, baseX, baseY int) []*baseWithData {

	var npcBases []*baseWithData
	hero.WalkHomeNpcBase(func(base *entity.HomeNpcBase) bool {
		npcBases = append(npcBases, r.newHomeNpcBase(hero.Id(), base, baseX, baseY))
		return false
	})

	return npcBases
}

func (r *Realm) newHomeNpcBase(ownerHeroId int64, base *entity.HomeNpcBase, baseX, baseY int) *baseWithData {

	homeNpcBaseData := base.GetData()

	basePos := hexagon.ShiftEvenOffset(baseX, baseY, homeNpcBaseData.EvenOffsetX, homeNpcBaseData.EvenOffsetY)
	baseX, baseY = basePos.XY()

	data := homeNpcBaseData.Data
	prosprity := u64.Min(base.GetProsprity(), data.ProsperityCapcity)

	basicNpc := r.newBasicNpcBase(base.Id(), data, baseX, baseY, realmface.BaseTypeNpc,
		shared_proto.TroopType_TT_NORMAL_NPC, data.Id, npcid.NpcType_HomeNpc)
	basicNpc.prosperity = prosprity

	basic := &homeNpcBase{
		basicNpcBase: basicNpc,
		ownerHeroId:  ownerHeroId,
		data:         homeNpcBaseData,
	}

	return newBaseWithData(basic)
}

func (r *Realm) newBasicNpcBaseWithData(id int64, data *basedata.NpcBaseData, baseX, baseY int) *baseWithData {

	basic := r.newBasicNpcBase(id, data, baseX, baseY, realmface.BaseTypeNpc,
		shared_proto.TroopType_TT_NORMAL_NPC, data.Id, npcid.NpcType_Monster)

	return newBaseWithData(basic)
}

func (r *Realm) newGuildWorkshopBase(id int64, data *basedata.NpcBaseData, guildId int64,
	baseX, baseY int, startTime, endTime int32, isComplete bool) *baseWithData {

	basic := r.newBasicNpcBase(id, data, baseX, baseY, realmface.BaseTypeNpc,
		shared_proto.TroopType_TT_NORMAL_NPC, data.Id, npcid.NpcType_Guild)
	basic.guildId = guildId

	workshop := &guildWorkshopBase{
		basicNpcBase: basic,
		startTime:    startTime,
		endTime:      endTime,
		isComplete:   isComplete,
	}

	return newBaseWithData(workshop)
}

func (r *Realm) newXiongNuNpcBase(id int64, info xiongnuface.RResistXiongNuInfo, baseX, baseY int) *baseWithData {

	basic := &xiongNuNpcBase{
		basicBase: basicBase{
			realmId:  r.id,
			id:       id,
			idBytes:  idbytes.ToBytes(id),
			heroName: info.Data().NpcBaseData.Npc.Name,

			baseX:             baseX,
			baseY:             baseY,
			baseLevel:         info.Data().NpcBaseData.BaseLevel,
			prosperity:        info.Data().NpcBaseData.ProsperityCapcity,
			prosperityCapcity: info.Data().NpcBaseData.ProsperityCapcity,
		},
		baseLevelInfo: newNpcDataBaseLevelInfo(info.Data().NpcBaseData, info.Data().NpcBaseData.EncodeSnapshot(id, r.id, baseX, baseY)),
		info:          info,
	}

	data := newBaseWithData(basic)

	basic.defenser = r.newNpcTroop(id, data, data.BaseLevel(), data, data.BaseLevel(), info.Data().NpcBaseData.Npc, realmface.Temp)
	basic.defenser.onChanged()

	return data
}

type basicBase struct {
	realmId  int64
	id       int64
	idBytes  []byte
	heroName string

	baseX, baseY       int
	baseLevel          uint64
	prosperity         uint64
	prosperityCapcity  uint64
	stopLostProsperity bool

	vipLevel uint64

	guildId int64

	cantSeeMeMsg pbutil.Buffer
}

type heroHome struct {
	basicBase

	homeNpcBase map[int64]*baseWithData

	whiteFlagGuildId       int64
	whiteFlagDisappearTime int32

	mianDisappearTime int32

	isRestore bool

	nextRestoreProsperityTime time.Time

	moveBaseRestoreProsperityBufEndTime time.Time // 迁城恢复繁荣度的buf的结束时间

	guanFuLevel        uint64 // 官府等级
	guanFuLevelData    *domestic_data.GuanFuLevelData
	outerCityUnlockBit int32 // 外城解锁数据

	sign  string // 个人签名
	title uint64 // 称号

	// key是宝藏id value是宝藏控制结束时间
	keepBaozMap map[int64]int32
}

func getRefStat(ref *atomic.Value) *data.SpriteStat {
	if i := ref.Load(); i != nil {
		if stat, ok := i.(*data.SpriteStat); ok {
			if stat != data.EmptyStat() {
				return stat
			}
		}
	}
	return nil
}

func (b *heroHome) AddKeepBaozMap(baozId int64, keepEndTime int32) {
	b.keepBaozMap[baozId] = keepEndTime
}

func (b *heroHome) TryRemoveAllKeepBaoz(r *Realm, ctime time.Time) {

	ctimeUnix := timeutil.Marshal32(ctime)

	for k, v := range b.keepBaozMap {
		delete(b.keepBaozMap, k)

		if ctimeUnix < v {
			// 没过期
			if keepBase := r.getBase(k); keepBase != nil {
				if baoz := GetBaoZangBase(keepBase); baoz != nil {
					if baoz.heroType == HeroTypeKiller && baoz.heroId == b.Id() {
						baoz.setKillerKeepEndTime(ctimeUnix)
						r.broadcastBaseInfoToCared(keepBase, addBaseTypeUpdate, 0)
					}
				}
			}
		}
	}
}

func (b *heroHome) TryRemoveKeepBaozFarAway(r *Realm, ctime time.Time, newX, newY, distance int) {

	ctimeUnix := timeutil.Marshal32(ctime)

	for k, v := range b.keepBaozMap {

		if ctimeUnix < v {
			// 没过期
			if keepBase := r.getBase(k); keepBase != nil {
				if baoz := GetBaoZangBase(keepBase); baoz != nil {
					if baoz.heroType == HeroTypeKiller && baoz.heroId == b.Id() {
						if util.IsInRange(keepBase.BaseX(), keepBase.BaseY(), newX, newY, distance) {
							// 保留下来，不要删除
							continue
						}

						baoz.setKillerKeepEndTime(ctimeUnix)
						r.broadcastBaseInfoToCared(keepBase, addBaseTypeUpdate, 0)
					}
				}
			}
		}

		delete(b.keepBaozMap, k)
	}
}

func (b *heroHome) TryRemoveKeepBaoz(baozId int64) {
	delete(b.keepBaozMap, baozId)
}

func (b *heroHome) tryRemoveExpireKeepBaoz(ctime time.Time) {
	ctimeUnix := timeutil.Marshal32(ctime)

	for k, v := range b.keepBaozMap {
		if v <= ctimeUnix {
			delete(b.keepBaozMap, k)
		}
	}
}

func (b *heroHome) tryRestoreProsperity(ctime time.Time, duration time.Duration) bool {
	t := b.nextRestoreProsperityTime
	if ctime.Before(t) {
		return false
	}

	b.nextRestoreProsperityTime = t.Add(duration)
	if ctime.After(b.nextRestoreProsperityTime) {
		// 加完之后，还小于当前时间，以当前时间未准
		b.nextRestoreProsperityTime = ctime.Add(duration)
		return !timeutil.IsZero(t)
	}

	return true
}

func (base *heroHome) getLayoutCube(layoutData *domestic_data.BuildingLayoutData) cb.Cube {
	return hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), layoutData.RegionOffsetX, layoutData.RegionOffsetY)
}

func (h *heroHome) WhiteFlagGuildId() int64 {
	return h.whiteFlagGuildId
}

func (h *heroHome) WhiteFlagDisappearTime() int32 {
	return h.whiteFlagDisappearTime
}

func (h *heroHome) SetWhiteFlag(guildId int64, disappearTime int32) {
	h.whiteFlagGuildId = guildId
	h.whiteFlagDisappearTime = disappearTime
}

func (h *heroHome) MianDisappearTime() int32 {
	return h.mianDisappearTime
}

func (h *heroHome) SetMianDisappearTime(disappearTime int32) {
	h.mianDisappearTime = disappearTime
}

func (h *heroHome) SetSign(toSet string) {
	h.sign = toSet
}

func (h *heroHome) SetTitle(toSet uint64) {
	h.title = toSet
}

func (h *heroHome) UpdateMoveBaseBuf(hero *entity.Hero) {
	h.moveBaseRestoreProsperityBufEndTime = hero.GetMoveBaseRestoreProsperityBufEndTime()
	guanFu := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU)
	if guanFu != nil {
		h.guanFuLevel = guanFu.Level
	}
}

func (h *heroHome) GetGuanFuLevelData(r *Realm) *domestic_data.GuanFuLevelData {
	if h.guanFuLevelData != nil && h.guanFuLevelData.Level == h.guanFuLevel {
		return h.guanFuLevelData
	}

	data := r.services.datas.GuanFuLevelData().Must(h.guanFuLevel)
	h.guanFuLevelData = data
	return data
}

func (h *heroHome) UpdateOutCityUnlockBit(hero *entity.Hero) {
	h.outerCityUnlockBit = hero.Domestic().OuterCities().UnlockBit()
}

// BaseType
func (base *heroHome) BaseType() realmface.BaseType {
	return realmface.BaseTypeHome
}

func (base *heroHome) TroopType() shared_proto.TroopType {
	return shared_proto.TroopType_TT_HERO
}

func (base *multiLevelNpcBase) BaseType() realmface.BaseType {
	return realmface.BaseTypeNpc
}

func (base *multiLevelNpcBase) TroopType() shared_proto.TroopType {
	return shared_proto.TroopType_TT_MULTI_LEVEL
}

// 持续掠夺时间

func (b *heroHome) BeenRobTickDuration(r *Realm) time.Duration {
	return r.services.datas.RegionConfig().RobTickDuration
}

func (b *heroHome) BeenRobMaxDuration(r *Realm) time.Duration {
	return r.services.datas.RegionConfig().RobMaxDuration
}

func (b *heroHome) BeenRobLostProsperityDuration(r *Realm) time.Duration {
	return r.services.datas.RegionConfig().ReduceProsperityDuration
}

func GetGuildWorkshopBase(base *baseWithData) *guildWorkshopBase {
	if b, ok := base.internalBase.(*guildWorkshopBase); ok {
		return b
	}
	return nil
}

type guildWorkshopBase struct {
	*basicNpcBase

	// 显示修建进度
	startTime int32
	endTime   int32

	isComplete bool

	nextReduceProsperityTime int32
}

func (b *guildWorkshopBase) getProgressInfo() (progress, totalProgress, progressType int32) {
	if !b.isComplete {
		// 修建倒计时阶段
		progressType = TimeProgress
		progress = b.startTime
		totalProgress = b.endTime
	} else {
		// 建成阶段
		progressType = DefaultProgress
		progress = u64.Int32(b.Prosperity())
		totalProgress = u64.Int32(b.ProsperityCapcity())
	}
	return
}

func (b *guildWorkshopBase) newUpdateProgressBarMsg() pbutil.Buffer {
	progress, totalProgress, progressType := b.getProgressInfo()
	return region.NewS2cUpdateBaseProgressMsg(b.IdBytes(), progress, totalProgress, progressType).Static()
}

func (r *Realm) newBasicNpcBase(id int64, data *basedata.NpcBaseData, baseX, baseY int,
	baseType realmface.BaseType, troopType shared_proto.TroopType, npcDataId uint64, npcType npcid.NpcType) *basicNpcBase {

	basic := &basicNpcBase{
		basicBase: basicBase{
			realmId:  r.id,
			id:       id,
			idBytes:  idbytes.ToBytes(id),
			heroName: data.Npc.Name,

			baseX:             baseX,
			baseY:             baseY,
			baseLevel:         data.BaseLevel,
			prosperity:        data.ProsperityCapcity,
			prosperityCapcity: data.ProsperityCapcity,
			guildId:           0,
		},
		baseType:  baseType,
		troopType: troopType,
		npcDataId: npcDataId,
		npcType:   npcType,
		info:      newNpcDataBaseLevelInfo(data, data.EncodeSnapshot(id, r.id, baseX, baseY)),
	}

	return basic
}

// Npc
type basicNpcBase struct {
	basicBase

	baseType realmface.BaseType

	troopType shared_proto.TroopType

	npcDataId uint64

	npcType npcid.NpcType

	info *npcDataBaseLevelInfo

	updateBaseInfoMsg pbutil.Buffer
}

func (b *basicNpcBase) BaseType() realmface.BaseType {
	return b.baseType
}

func (base *basicNpcBase) TroopType() shared_proto.TroopType {
	return base.troopType
}

func (base *basicNpcBase) getBaseInfoByLevel(level uint64) baseInfo {
	return base.info
}

func (b *basicNpcBase) getNpcConfig() (dataId uint64, npcType npcid.NpcType) {
	return b.npcDataId, b.npcType
}

func (b *basicNpcBase) ClearUpdateBaseInfoMsg() {
	b.updateBaseInfoMsg = nil
}

func GetHomeNpcBase(base *baseWithData) *homeNpcBase {
	if homeNpcBase, ok := base.internalBase.(*homeNpcBase); ok {
		return homeNpcBase
	}
	return nil
}

// 英雄自己的Npc
type homeNpcBase struct {
	*basicNpcBase

	ownerHeroId int64                     // 属于哪个英雄的
	data        *basedata.HomeNpcBaseData // 野怪模板
}

func (r *Realm) newMultiLevelNpcBase(id int64, data *regdata.RegionMultiLevelNpcData, baseX, baseY int) *baseWithData {

	firstLevel := data.GetFirstLevel().Npc
	basic := &multiLevelNpcBase{
		basicBase: basicBase{
			realmId:  r.id,
			id:       id,
			idBytes:  idbytes.ToBytes(id),
			heroName: firstLevel.Npc.Name,

			baseX:             baseX,
			baseY:             baseY,
			baseLevel:         firstLevel.BaseLevel,
			prosperity:        firstLevel.ProsperityCapcity,
			prosperityCapcity: firstLevel.ProsperityCapcity,
			guildId:           0,
		},
		data: data,
	}

	for _, v := range data.LevelBases {
		data := v.Npc
		basic.infos = append(basic.infos, newMultiLevelBaseInfo(data, data.EncodeSnapshot(id, r.id, baseX, baseY), v))
	}

	return newBaseWithData(basic)
}

func GetMultiLevelNpcBase(base *baseWithData) *multiLevelNpcBase {
	if b, ok := base.internalBase.(*multiLevelNpcBase); ok {
		return b
	}
	return nil
}

// 多等级的Npc
type multiLevelNpcBase struct {
	basicBase

	data *regdata.RegionMultiLevelNpcData

	infos []*npcDataBaseLevelInfo

	updateBaseInfoMsg pbutil.Buffer
}

func GetXiongNuBase(base *baseWithData) *xiongNuNpcBase {
	if b, ok := base.internalBase.(*xiongNuNpcBase); ok {
		return b
	}
	return nil
}

// 匈奴npc
type xiongNuNpcBase struct {
	basicBase

	baseLevelInfo *npcDataBaseLevelInfo

	info xiongnuface.RResistXiongNuInfo

	defenser *troop

	updateBaseInfoMsg pbutil.Buffer

	invadeTargetCountRef atomic.Value

	//removeTime time.Time
}

//func (b *xiongNuNpcBase) tryRemove(ctime time.Time) bool {
//	return b.removeTime.Before(ctime)
//}

func (b *xiongNuNpcBase) setInvadeTargetCount(itc map[int64]uint64) {
	b.invadeTargetCountRef.Store(i64.NewGetU64(itc))
}

func (b *xiongNuNpcBase) getInvadeTargetCount() i64.GetU64 {
	if itc := b.invadeTargetCountRef.Load(); itc != nil {
		return itc.(i64.GetU64)
	}
	return i64.EmptyGetU64()
}

func (r *Realm) GetXiongNuInvateTargetCount(id int64) i64.GetU64 {
	if b := r.getBase(id); b != nil {
		if base := GetXiongNuBase(b); base != nil {
			return base.getInvadeTargetCount()
		}
	}
	return nil
}

func (b *xiongNuNpcBase) BaseType() realmface.BaseType {
	return realmface.BaseTypeNpc
}

func (base *xiongNuNpcBase) TroopType() shared_proto.TroopType {
	return shared_proto.TroopType_TT_XIONG_NU
}

func (b *xiongNuNpcBase) Info() xiongnuface.RResistXiongNuInfo {
	return b.info
}

func (b *xiongNuNpcBase) Defenser() *troop {
	return b.defenser
}

func (b *xiongNuNpcBase) getBaseInfoByLevel(level uint64) baseInfo {
	return b.baseLevelInfo
}

func (b *xiongNuNpcBase) Data() *xiongnu.ResistXiongNuData {
	return b.info.Data()
}

func (b *xiongNuNpcBase) MianDisappearTime() int32 {
	return timeutil.Marshal32(b.info.InvadeEndTime())
}

func (b *xiongNuNpcBase) encodeNpcProto(self *baseWithData, r *Realm) *server_proto.NpcBaseProto {
	result := encodeNpcBaseProto(b, b.defenser, r)
	result.GuildId = b.Info().GuildId()

	for _, t := range self.selfTroops {
		if t.State() == realmface.Defending {
			troopProto := &server_proto.NpcDefendingTroopProto{}
			troopProto.Id = t.Id()
			if t.mmd != nil {
				troopProto.DataId = t.mmd.Id
			}

			for _, v := range t.Captains() {
				troopProto.CaptainIndex = append(troopProto.CaptainIndex, int32(v.Index()))
				troopProto.CaptainSoldier = append(troopProto.CaptainSoldier, v.Proto().Soldier)
			}

			result.DefendingTroop = append(result.DefendingTroop, troopProto)
		}
	}

	return result
}

func (r *Realm) newJunTuanNpcBase(id int64, data *regdata.JunTuanNpcData, baseX, baseY int) *baseWithData {

	basicNpc := r.newBasicNpcBase(id, data.Npc, baseX, baseY, realmface.BaseTypeNpc,
		shared_proto.TroopType_TT_NORMAL_NPC, data.Id, npcid.NpcType_JunTuan)

	basic := &junTuanNpcBase{
		basicNpcBase: basicNpc,
		data:         data,
	}

	base := newBaseWithData(basic)

	return base
}

func GetJunTuanBase(base *baseWithData) *junTuanNpcBase {
	if b, ok := base.internalBase.(*junTuanNpcBase); ok {
		return b
	}
	return nil
}

// 军团怪Npc
type junTuanNpcBase struct {
	*basicNpcBase

	data *regdata.JunTuanNpcData
}

func (b *junTuanNpcBase) Data() *regdata.JunTuanNpcData {
	return b.data
}

func (b *junTuanNpcBase) encodeNpcProto(self *baseWithData, r *Realm) *server_proto.NpcBaseProto {
	result := encodeNpcBaseProto(b, nil, r)

	for _, t := range self.selfTroops {
		if t.State() == realmface.Defending {
			troopProto := &server_proto.NpcDefendingTroopProto{}
			troopProto.Id = t.Id()
			if t.mmd != nil {
				troopProto.DataId = t.mmd.Id
			}

			for _, v := range t.Captains() {
				troopProto.CaptainIndex = append(troopProto.CaptainIndex, int32(v.Index()))
				troopProto.CaptainSoldier = append(troopProto.CaptainSoldier, v.Proto().Soldier)
			}

			result.DefendingTroop = append(result.DefendingTroop, troopProto)
		}
	}

	return result
}

func (b *junTuanNpcBase) updateRoBase(self *baseWithData, proto *server_proto.RoBaseProto) {
	proto.ProsperityCapcity = b.data.TroopCount

	var defendingTroopCount uint64
	var soldier int32
	for _, t := range self.targetingTroops {
		if t.State() == realmface.Defending {
			defendingTroopCount++
			for _, captain := range t.Captains() {
				soldier += captain.Proto().Soldier
			}
		}
	}
	proto.Prosperity = defendingTroopCount
	proto.Soldier = u64.FromInt32(soldier)
}

func (r *Realm) newBaozNpcBase(id int64, data *regdata.BaozNpcData, baseX, baseY int, heroId int64, expireTime int32) *baseWithData {

	basicNpc := r.newBasicNpcBase(id, data.Npc, baseX, baseY, realmface.BaseTypeNpc,
		shared_proto.TroopType_TT_BAO_ZANG, data.Id, npcid.NpcType_BaoZang)

	var heroType int32 = 0
	if heroId != 0 {
		heroType = 1
	}

	basic := &baozNpcBase{
		basicNpcBase: basicNpc,
		data:         data,
		heroType:     heroType,
		heroId:       heroId,
		heroEndTime:  expireTime,
	}

	base := newBaseWithData(basic)

	basic.defenser = r.newNpcTroop(id, base, base.BaseLevel(), base, base.BaseLevel(), data.Npc.Npc, realmface.Temp)
	basic.defenser.onChanged()

	// 更新一下士兵数
	basic.updateRoBase(base.getRoBase())

	return base
}

func GetBaoZangBase(base *baseWithData) *baozNpcBase {
	if b, ok := base.internalBase.(*baozNpcBase); ok {
		return b
	}
	return nil
}

const (
	HeroTypeKiller  = 0
	HeroTypeCreater = 1
)

// 宝藏Npc
type baozNpcBase struct {
	*basicNpcBase

	data *regdata.BaozNpcData

	defenser *troop

	// 0-击杀者 1-创建者
	heroType    int32
	heroId      int64
	heroEndTime int32

	heroBytes []byte
}

func (b *baozNpcBase) setKiller(killer int64, killerBytes []byte, killerKeepEndTime int32) {
	b.heroId = killer
	b.heroBytes = killerBytes
	b.setKillerKeepEndTime(killerKeepEndTime)
}

func (b *baozNpcBase) setKillerKeepEndTime(killerKeepEndTime int32) {
	b.heroEndTime = killerKeepEndTime
}

func (b *baozNpcBase) Defenser() *troop {
	return b.defenser
}

func (b *baozNpcBase) HasDefenser() bool {
	if b.defenser != nil {
		return b.defenser.AliveSoldier() > 0
	}
	return false
}

func (b *baozNpcBase) updateRoBase(proto *server_proto.RoBaseProto) {
	proto.Prosperity = b.prosperity

	if b.defenser != nil {
		proto.FightAmount = b.defenser.FightAmount()

		proto.CaptainSoldier = make([]uint64, constants.CaptainCountPerTroop)
		var totalSoldier uint64
		for _, c := range b.defenser.Captains() {
			if c != nil && c.Index() > 0 && c.Index() <= len(proto.CaptainSoldier) {
				soldier := u64.FromInt32(c.Proto().Soldier)
				proto.CaptainSoldier[c.Index()-1] = soldier
				totalSoldier += soldier
			}
		}
		proto.Soldier = totalSoldier
	}
	proto.HeroType = b.heroType
	proto.HeroId = b.heroId
	proto.HeroEndTime = b.heroEndTime
}

func (b *baozNpcBase) Data() *regdata.BaozNpcData {
	return b.data
}

func (b *baozNpcBase) encodeNpcProto(self *baseWithData, r *Realm) *server_proto.NpcBaseProto {

	proto := encodeNpcBaseProto(b, b.defenser, r)
	proto.HeroType = b.heroType
	proto.HeroId = b.heroId
	proto.HeroEndTime = int64(b.heroEndTime)

	return proto
}

func encodeNpcBaseProto(b realmface.Base, defenser *troop, r *Realm) *server_proto.NpcBaseProto {

	proto := &server_proto.NpcBaseProto{
		BaseRegion: r.Id(),
		BaseId:     b.Id(),
		BaseX:      int32(b.BaseX()),
		BaseY:      int32(b.BaseY()),
		Prosperity: b.Prosperity(),
	}

	if defenser != nil {
		for _, v := range defenser.Captains() {
			proto.CaptainIndex = append(proto.CaptainIndex, int32(v.Index()))
			proto.CaptainSoldier = append(proto.CaptainSoldier, v.Proto().Soldier)
		}
	}

	return proto
}

// Base更新消息

func (b *heroHome) NewUpdateBaseInfoMsg(r *Realm, addType int32, base *baseWithData) pbutil.Buffer { // 不缓存
	return b.newUpdateBaseInfoMsg(r, addType, b.mianDisappearTime,
		b.outerCityUnlockBit, b.sign, b.title, b.whiteFlagGuildId, b.whiteFlagDisappearTime)
}

func (b *basicNpcBase) NewUpdateBaseInfoMsg(r *Realm, addType int32, base *baseWithData) pbutil.Buffer {
	// 缓存
	if b.updateBaseInfoMsg == nil {
		b.updateBaseInfoMsg = r.newUpdateNpcBaseMsg(b, 0, false, 0, false, nil, 0, 0, 0, 0, DefaultProgress)
	}

	return b.updateBaseInfoMsg
}

func (b *multiLevelNpcBase) NewUpdateBaseInfoMsg(r *Realm, addType int32, base *baseWithData) pbutil.Buffer {
	// 缓存
	if b.updateBaseInfoMsg == nil {
		b.updateBaseInfoMsg = r.newUpdateNpcBaseMsg(b, 0, false, 0, false, nil, 0, 0, 0, 0, DefaultProgress)
	}

	return b.updateBaseInfoMsg
}

func (b *xiongNuNpcBase) NewUpdateBaseInfoMsg(r *Realm, addType int32, base *baseWithData) pbutil.Buffer {
	// 缓存
	if b.updateBaseInfoMsg == nil {
		var progress int32
		if b.defenser != nil {
			for _, captain := range b.defenser.Captains() {
				progress += captain.Proto().Soldier
			}
		}
		for _, t := range base.targetingTroops {
			if t.State() == realmface.Defending {
				for _, captain := range t.Captains() {
					progress += captain.Proto().Soldier
				}
			}
		}

		totalProgress := u64.Int32(b.info.Data().GetTotalSoldier())
		b.updateBaseInfoMsg = r.newUpdateNpcBaseMsg(b, b.Info().GuildId(), false, 0, false, nil, 0, 0, progress, totalProgress, DefaultProgress)
	}

	return b.updateBaseInfoMsg
}

func (b *junTuanNpcBase) NewUpdateBaseInfoMsg(r *Realm, addType int32, base *baseWithData) pbutil.Buffer {
	// 缓存
	if b.updateBaseInfoMsg == nil {
		var progress int32
		for _, t := range base.targetingTroops {
			if t.State() == realmface.Defending {
				for _, captain := range t.Captains() {
					progress += captain.Proto().Soldier
				}
			}
		}
		totalProgress := u64.Int32(b.data.GetTotalSoldier())
		b.updateBaseInfoMsg = r.newUpdateNpcBaseMsg(b, 0, false, b.Prosperity(), false, nil, 0, 0, progress, totalProgress, DefaultProgress)
	}

	return b.updateBaseInfoMsg
}

func (b *baozNpcBase) NewUpdateBaseInfoMsg(r *Realm, addType int32, base *baseWithData) pbutil.Buffer {
	// 缓存
	if b.updateBaseInfoMsg == nil {
		heroBytes := b.heroBytes
		if b.heroId != 0 {
			if len(heroBytes) > 0 {
				if hero := r.services.heroSnapshotService.GetFromCache(b.heroId); hero != nil {
					heroBytes = hero.EncodeBasic4ClientBytes()
					b.heroBytes = heroBytes
				}
			} else {
				if hero := r.services.heroSnapshotService.Get(b.heroId); hero != nil {
					heroBytes = hero.EncodeBasic4ClientBytes()
					b.heroBytes = heroBytes
				}
			}
		}

		b.updateBaseInfoMsg = r.newUpdateNpcBaseMsg(b, 0, false, b.prosperity, b.HasDefenser(),
			heroBytes, b.heroEndTime, b.heroType, 0, 0, DefaultProgress)
	}

	return b.updateBaseInfoMsg
}

func (b *guildWorkshopBase) NewUpdateBaseInfoMsg(r *Realm, addType int32, base *baseWithData) pbutil.Buffer {
	// 缓存
	if b.updateBaseInfoMsg == nil {
		progress, totalProgress, progressType := b.getProgressInfo()
		b.updateBaseInfoMsg = r.newUpdateNpcBaseMsg(b, b.guildId, true, b.prosperity, false,
			nil, 0, 0, progress, totalProgress, progressType)
	}

	return b.updateBaseInfoMsg
}

const (
	DefaultProgress = 0 // 分子分母类型
	TimeProgress    = 1 // 倒计时类型
)

func (r *Realm) newUpdateNpcBaseMsg(b internalBase, guildId int64, hasGuildDetails bool,
	prosperity uint64, hasDefenser bool, heroBytes []byte, heroEndTime, heroType, progress, totalProgress, progressType int32) pbutil.Buffer {
	dataId, npcType := b.getNpcConfig()

	var guildName, guildFlagName string
	var country int32
	if guildId != 0 && hasGuildDetails {
		if g := r.services.guildService.GetSnapshot(guildId); g != nil {
			guildName = g.Name
			guildFlagName = g.FlagName
			country = u64.Int32(g.Country.Id)
		}
	}

	return region.NewS2cUpdateNpcBaseInfoMsg(int32(r.id), b.IdBytes(), int32(b.BaseX()),
		int32(b.BaseY()), u64.Int32(dataId), npcType,
		i64.Int32(guildId), guildName, guildFlagName, country, b.MianDisappearTime(),
		u64.Int32(prosperity), hasDefenser, heroBytes, heroEndTime, heroType,
		progress, totalProgress, progressType).Static()
}

const (
	// 0-进入视野 1-新建角色 2-重建 3-迁移 4-更新
	addBaseTypeCanSee int32 = iota
	addBaseTypeNewHero
	addBaseTypeReborn
	addBaseTypeTransfer
	addBaseTypeUpdate
)

const (
	// 0-离开视野 1-迁移 2-流亡
	removeBaseTypeCantSee int32 = iota
	removeBaseTypeTransfer
	removeBaseTypeBroken
)

func (b *basicBase) newUpdateBaseInfoMsg(r *Realm, addType int32,
	mianDisappearTime, outerCityUnlockBit int32, sign string, title uint64,
	whiteFlagGuildId int64, whiteFlagDisappearTime int32) pbutil.Buffer { // 不缓存

	var guildFlagName, whiteGuildFlagName string
	var countryId uint64
	if b.guildId != 0 {
		if g := r.services.guildService.GetSnapshot(b.guildId); g != nil {
			guildFlagName = g.FlagName
			countryId = g.GetCountryId()
		}
	}

	if whiteFlagGuildId != 0 {
		if g := r.services.guildService.GetSnapshot(whiteFlagGuildId); g != nil {
			whiteGuildFlagName = g.FlagName
		}
	}

	return region.NewS2cAddBaseUnitMarshalMsg(addType, &shared_proto.BaseUnitProto{
		HeroId:                 b.idBytes,
		HeroName:               b.heroName,
		GuildId:                i64.Int32(b.guildId),
		GuildFlagName:          guildFlagName,
		CountryId:              u64.Int32(countryId),
		Level:                  u64.Int32(b.baseLevel),
		VipLevel:               u64.Int32(b.vipLevel),
		BaseX:                  int32(b.baseX),
		BaseY:                  int32(b.baseY),
		MianDisappearTime:      mianDisappearTime,
		WhiteFlagGuildId:       i64.Int32(whiteFlagGuildId),
		WhiteFlagGuildFlagName: whiteGuildFlagName,
		WhiteFlagDisappearTime: whiteFlagDisappearTime,
		OuterCityUnlockBit:     outerCityUnlockBit,
		Sign:                   sign,
		Title:                  u64.Int32(title),
		Prosperty:              u64.Int32(b.prosperity),
		StopLostProsperity:     b.stopLostProsperity,
	})
}

func (b *baseWithData) getBaseDefenser(r *Realm, level uint64) *troop {
	switch b.BaseType() {
	case realmface.BaseTypeHome:
		return getHeroBaseDefenser(r, b)
	case realmface.BaseTypeNpc:
		return getNpcBaseDefenser(r, b, level)
	}

	return nil
}

func getNpcBaseDefenser(r *Realm, b *baseWithData, level uint64) *troop {
	if npcid.IsXiongNuNpcId(b.Id()) {
		if base, ok := b.internalBase.(*xiongNuNpcBase); ok {
			return base.Defenser()
		}
	}

	if npcid.IsBaoZangNpcId(b.Id()) {
		if base, ok := b.internalBase.(*baozNpcBase); ok {
			t := base.Defenser()
			if t != nil && t.AliveSoldier() > 0 {
				return t
			}
			return nil
		}
	}

	var targetDefenser *troop

	data := b.internalBase.getBaseInfoByLevel(level).getNpcBaseData()
	if data != nil {
		targetDefenser = r.newNpcTroop(b.Id(), b, level, b, level, data.Npc, realmface.Temp)
		targetDefenser.setWall(data.Npc.WallLevel, data.Npc.WallStat, data.Npc.WallFixDamage)
	}

	return targetDefenser
}

func newNpcTroopCaptains(data *monsterdata.MonsterMasterData) []realmface.Captain {
	troopCaptains := make([]realmface.Captain, 0, len(data.Captains))
	for i, c := range data.Captains {
		if c == nil {
			continue
		}

		captain := &captain{
			idAtHero: c.Id,
			index:    i + 1,
			proto:    c.EncodeCaptainInfo(),
		}

		troopCaptains = append(troopCaptains, captain)
	}

	return troopCaptains
}

func getHeroBaseDefenser(r *Realm, b *baseWithData) *troop {
	var targetDefenser *troop
	r.heroBaseFuncWithSend(b.Base(), func(hero *entity.Hero, result herolock.LockResult) {
		// 获取防守队伍
		var wallLevel uint64
		var wallStat *data.SpriteStat
		var wallFixDamage uint64
		// 主城防守队伍
		defenser := hero.GetHomeDefenser()
		if b := hero.Domestic().GetBuilding(shared_proto.BuildingType_CHENG_QIANG); b != nil {
			wallLevel = b.Level
			wallStat = hero.BuildingEffect().HomeWallStat()
			wallFixDamage = hero.BuildingEffect().HomeWallFixDamage()
		}

		ctime := r.services.timeService.CurrentTime()
		heromodule.TryAutoFullSoldier(hero, result, r.services.datas, defenser, ctime,
			r.services.datas.MiscGenConfig().AutoFullSoldoerDuration)

		if defenser != nil && !defenser.IsOutside() && defenser.HasSoldier() {
			targetDefenser = r.newTempTroop(defenser.Id(), b, hero.BaseLevel(), hero.EncodeDefenseCaptainInfo, defenser.Pos())
			targetDefenser.setWall(wallLevel, wallStat, wallFixDamage)
		} else if copyDefCaptain := hero.NewCopyDefCaptainInfo(ctime); len(copyDefCaptain) > 0 {
			// 特殊的id
			targetDefenser = r.newTempTroopWithInfo(0, b, hero.BaseLevel(), copyDefCaptain)
			targetDefenser.setWall(wallLevel, wallStat, wallFixDamage)
		} else if wallStat != nil {
			targetDefenser = r.newTempTroop(b.Id(), b, hero.BaseLevel(), hero.EncodeDefenseCaptainInfo, []*entity.TroopPos{})
			targetDefenser.setWall(wallLevel, wallStat, wallFixDamage)
		}

		return
	})
	return targetDefenser
}

// 返回当前正在掠夺的队伍中, 战斗力最低的那个
func (b *baseWithData) getARobbingTroopWithLowestFightAmount() (result *troop, hasMore bool) {
	//var resultFightAmount uint64

	for _, t := range b.targetingTroops {
		if t.state == realmface.Robbing {
			//fightAmount := t.FightAmount()
			//
			//hasMore = result != nil
			//if result == nil || resultFightAmount < fightAmount {
			//	result, resultFightAmount = t, fightAmount
			//}

			// 改成最晚到达的那个
			hasMore = result != nil
			if result == nil || result.moveArriveTime.Before(t.moveArriveTime) {
				result = t
			}
		}
	}

	return
}

// 返回当前正在这个基地防守的队伍中, 战斗力最低的那个
func (b *baseWithData) getADefendingTroopWithLowestFightAmount(selfDefender *troop) (result *troop, hasMore bool) {
	//var resultFightAmount uint64

	for _, t := range b.targetingTroops {
		if t.state == realmface.Defending {
			//fightAmount := t.FightAmount()
			//
			//hasMore = result != nil
			//if result == nil || resultFightAmount < fightAmount {
			//	result, resultFightAmount = t, fightAmount
			//}

			// 改成最晚到达的那个
			hasMore = result != nil
			if result == nil || result.moveArriveTime.Before(t.moveArriveTime) {
				result = t
			}
		}
	}

	// 先搞别人的防守部队, 最后才搞自己驻守着的
	if result != nil {
		if !hasMore && selfDefender != nil {
			hasMore = true
		}

		return
	}

	return selfDefender, false
}

func (b *baseWithData) getAssistDefendingTroopCount() (result uint64) {
	for _, t := range b.targetingTroops {
		if t.state == realmface.Defending {
			result++
		}
	}

	return result
}

func (b *baseWithData) getAssisterTroops() (assister []*troop) {
	for _, t := range b.targetingTroops {
		if t.state == realmface.Defending {
			assister = append(assister, t)
		}
	}
	return
}

func (b *baseWithData) getDefendingTroopCount() uint64 {
	var result uint64

	for _, t := range b.selfTroops {
		if t.state == realmface.Defending && t.targetBase == nil { // target != nil 表示在给别人守城
			result++
		}
	}

	for _, t := range b.targetingTroops {
		if t.state == realmface.Defending {
			result++
		}
	}

	return result
}

func (b *baseWithData) getRobbingTroopCount() uint64 {
	var result uint64

	for _, t := range b.targetingTroops {
		if t.state == realmface.Robbing {
			result++
		}
	}

	return result
}

func (b *baseWithData) isHeroHomeBase() bool {
	return b.BaseType() == realmface.BaseTypeHome
}

func (b *baseWithData) isNpcBase() bool {
	return b.BaseType() == realmface.BaseTypeNpc
}

func (b *baseWithData) isHeroBaseType() bool {
	return isHeroBaseType(b.Base())
}

func isHeroBaseType(base realmface.Base) bool {
	switch base.BaseType() {
	case realmface.BaseTypeHome:
		return true
	}

	return false
}

func (r *Realm) heroBaseFunc(base realmface.Base, f herolock.Func) {
	if isHeroBaseType(base) {
		r.services.heroDataService.Func(base.Id(), f)
	}
}

func (r *Realm) heroBaseFuncNotError(base realmface.Base, f herolock.FuncNotError) bool {
	if isHeroBaseType(base) {
		return r.services.heroDataService.FuncNotError(base.Id(), f)
	}
	return true
}

func (r *Realm) heroBaseFuncWithSend(base realmface.Base, f herolock.SendFunc) (hasError bool) {
	if isHeroBaseType(base) {
		return r.services.heroDataService.FuncWithSend(base.Id(), f)
	}
	return true
}

func (r *Realm) heroFuncWithSend(heroId int64, f herolock.SendFunc) (hasError bool) {
	if !npcid.IsNpcId(heroId) {
		return r.services.heroDataService.FuncWithSend(heroId, f)
	}
	return true
}

func (b *basicBase) GetHeroBasicProto(heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot) *shared_proto.HeroBasicProto {
	snapshot := heroSnapshotGetter(b.Id())
	if snapshot != nil {
		return snapshot.EncodeBasic4Client()
	}

	// 找不到？搞个默认的耍耍
	return &shared_proto.HeroBasicProto{
		Id:    b.idBytes,
		Name:  b.heroName,
		Level: 1,
	}
}

func (b *basicBase) EncodeAsHeroBasicSnapshot(heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot) *shared_proto.HeroBasicSnapshotProto {

	proto := &shared_proto.HeroBasicSnapshotProto{}

	proto.Basic = b.GetHeroBasicProto(heroSnapshotGetter)
	proto.BaseRegion = i64.Int32(b.realmId)
	proto.BaseLevel = u64.Int32(b.baseLevel)
	proto.BaseX = imath.Int32(b.baseX)
	proto.BaseY = imath.Int32(b.baseY)

	return proto
}

func (b *basicBase) RegionID() int64 {
	return b.realmId
}

func (b *basicBase) Id() int64 {
	return b.id
}

func (b *basicBase) IdBytes() []byte {
	return b.idBytes
}

func (b *basicBase) HeroName() string {
	return b.heroName
}

func (b *basicBase) SetBaseLevel(newBaseLevel uint64) {
	if newBaseLevel <= 0 && b.prosperity > 0 {
		logrus.WithField("prosprity", b.prosperity).WithField("stack", string(debug.Stack())).Debug("将base等级设置为0")
	}
	b.baseLevel = newBaseLevel
}

func (b *basicBase) SetBaseXY(x, y int) {
	b.baseX = x
	b.baseY = y
}

func (b *basicBase) GetBaseLevel() uint64 {
	return b.baseLevel
}

func (b *basicBase) BaseX() int {
	return b.baseX
}

func (b *basicBase) BaseY() int {
	return b.baseY
}

func (b *basicBase) Prosperity() uint64 {
	return b.prosperity
}

func (b *basicBase) AddProsperity(toAdd uint64) uint64 {
	b.prosperity += toAdd
	b.prosperity = u64.Max(1, u64.Min(b.prosperity, b.prosperityCapcity))

	return b.prosperity
}

func (b *basicBase) ReduceProsperity(toReduce uint64) uint64 {
	b.prosperity = u64.Sub(b.prosperity, toReduce)
	return b.prosperity
}

func (b *basicBase) ReduceProsperityDontKill(toReduce uint64) uint64 {
	b.prosperity = u64.Max(u64.Sub(b.prosperity, toReduce), 1)
	return b.prosperity
}

func (b *basicBase) ProsperityCapcity() uint64 {
	return b.prosperityCapcity
}

func (b *basicBase) SetProsperityCapcity(toSet uint64) {
	b.prosperityCapcity = toSet
}

func (b *basicBase) TrySetStopLostProsperity(toSet bool) bool {
	if b.stopLostProsperity != toSet {
		b.stopLostProsperity = toSet
		return true
	}
	return false
}

func (b *basicBase) GuildId() int64 {
	return b.guildId
}

func (b *basicBase) UpdateHeroBasicInfo(hero *entity.Hero) {

	b.heroName = hero.Name()
	b.guildId = hero.GuildId()
}

// default

func (h *basicBase) GetCantSeeMeMsg() pbutil.Buffer {
	if h.cantSeeMeMsg == nil {
		h.cantSeeMeMsg = region.NewS2cRemoveBaseUnitMsg(removeBaseTypeCantSee, h.idBytes)
	}
	return h.cantSeeMeMsg
}

func (h *basicBase) MianDisappearTime() int32                 { return 0 }
func (h *basicBase) SetMianDisappearTime(disappearTime int32) {}

// base level info

func (base *multiLevelNpcBase) getBaseInfoByLevel(level uint64) baseInfo {
	index := u64.Sub(level, 1)
	n := len(base.infos)
	if index < uint64(n) {
		return base.infos[index]
	}

	return base.infos[n-1]
}

func (data *heroHome) getName() string {
	return data.heroName
}

func (base *heroHome) getBaseInfoByLevel(level uint64) baseInfo {
	return base
}

func (base *heroHome) BeenRobMaxCount(r *Realm) uint64 {
	return r.config().MaxRobbers
}

func (base *heroHome) getNpcBaseData() *basedata.NpcBaseData {
	return nil
}

func (base *heroHome) getHateData() *regdata.RegionMultiLevelNpcLevelData {
	return nil
}

func (base *heroHome) IsDestroyWhenLose() bool {
	return false
}

type baseInfo interface {
	getName() string

	getNpcBaseData() *basedata.NpcBaseData

	// 如果以后有多个地方使用，把返回对象换掉
	getHateData() *regdata.RegionMultiLevelNpcLevelData

	GetBaseLevel() uint64

	BeenRobTickDuration(r *Realm) time.Duration

	BeenRobMaxDuration(r *Realm) time.Duration

	BeenRobLostProsperityDuration(r *Realm) time.Duration

	BeenRobMaxCount(r *Realm) uint64

	IsDestroyWhenLose() bool

	GetHeroBasicProto(heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot) *shared_proto.HeroBasicProto
	EncodeAsHeroBasicSnapshot(heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot) *shared_proto.HeroBasicSnapshotProto
}

var _ baseInfo = (*npcDataBaseLevelInfo)(nil)
var _ baseInfo = (*heroHome)(nil)

func newNpcDataBaseLevelInfo(data *basedata.NpcBaseData, proto *shared_proto.HeroBasicSnapshotProto) *npcDataBaseLevelInfo {
	return &npcDataBaseLevelInfo{
		NpcBaseData: data,
		proto:       proto,
	}
}

func newMultiLevelBaseInfo(data *basedata.NpcBaseData, proto *shared_proto.HeroBasicSnapshotProto, multiLevelBaseData *regdata.RegionMultiLevelNpcLevelData) *npcDataBaseLevelInfo {
	return &npcDataBaseLevelInfo{
		NpcBaseData:        data,
		proto:              proto,
		multiLevelBaseData: multiLevelBaseData,
	}
}

type npcDataBaseLevelInfo struct {
	*basedata.NpcBaseData

	proto *shared_proto.HeroBasicSnapshotProto

	multiLevelBaseData *regdata.RegionMultiLevelNpcLevelData
}

func (data *npcDataBaseLevelInfo) getName() string {
	return data.Npc.Name
}

func (data *npcDataBaseLevelInfo) getNpcBaseData() *basedata.NpcBaseData {
	return data.NpcBaseData
}

func (data *npcDataBaseLevelInfo) getHateData() *regdata.RegionMultiLevelNpcLevelData {
	return data.multiLevelBaseData
}

func (data *npcDataBaseLevelInfo) GetBaseLevel() uint64 {
	return data.BaseLevel
}

func (data *npcDataBaseLevelInfo) BeenRobTickDuration(r *Realm) time.Duration {
	return data.TickDuration
}

func (data *npcDataBaseLevelInfo) BeenRobMaxDuration(r *Realm) time.Duration {
	return data.RobMaxDuration
}

func (data *npcDataBaseLevelInfo) BeenRobLostProsperityDuration(r *Realm) time.Duration {
	return data.LostProsperityDuration
}

func (data *npcDataBaseLevelInfo) BeenRobMaxCount(r *Realm) uint64 {
	return data.MaxRobbers
}

func (data *npcDataBaseLevelInfo) IsDestroyWhenLose() bool {
	return data.DestroyWhenLose
}

func (b *npcDataBaseLevelInfo) GetHeroBasicProto(heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot) *shared_proto.HeroBasicProto {
	return b.proto.Basic
}

func (b *npcDataBaseLevelInfo) EncodeAsHeroBasicSnapshot(heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot) *shared_proto.HeroBasicSnapshotProto {
	return b.proto
}

// npc config
func (b *basicBase) getNpcConfig() (dataId uint64, npcType npcid.NpcType) {
	return
}

func (b *basicBase) encodeNpcProto(self *baseWithData, r *Realm) *server_proto.NpcBaseProto {
	return nil
}

func (b *basicBase) ClearUpdateBaseInfoMsg() {
	return
}

func (b *multiLevelNpcBase) getNpcConfig() (dataId uint64, npcType npcid.NpcType) {
	return b.data.Id, npcid.NpcType_MultiLevelMonster
}

func (b *xiongNuNpcBase) getNpcConfig() (dataId uint64, npcType npcid.NpcType) {
	return b.info.Data().Level, npcid.NpcType_XiongNu
}

func NewRuinsBasePosInfoMap(r *Realm) *ruinsBasePosInfoMap {
	return &ruinsBasePosInfoMap{
		realmId:             r.Id(),
		infoMap:             make(map[cb.Cube]*ruinsBase),
		expireDuration:      r.services.datas.RegionConfig().RuinsBaseExpireDuration,
		regionConfig:        r.services.datas.RegionConfig(),
		broadcastToCaredPos: r.broadcastToCaredPos,
	}
}

func newRuinsBase(proto *server_proto.RuinsBasePosInfoProto) *ruinsBase {
	return &ruinsBase{
		proto:     proto,
		addMsg:    region.NewS2cAddRuinsBaseMsg(proto.PosX, proto.PosY),
		removeMsg: region.NewS2cRemoveRuinsBaseMsg(proto.PosX, proto.PosY),
		watchList: &sync.Map{},
	}
}

type ruinsBase struct {
	proto *server_proto.RuinsBasePosInfoProto

	addMsg    pbutil.Buffer
	removeMsg pbutil.Buffer

	//Kimi 索引地块
	blockIndex *block_index
	//kimi 观察者列表
	watchList *sync.Map
}

func (t *ruinsBase) CanSee(area *realmface.ViewArea) bool {
	if area == nil {
		return false
	}
	return area.CanSeePos(int(t.proto.PosX), int(t.proto.PosY))
}

//GetPos 获取坐标
func (t *ruinsBase) GetPos() [][2]int {
	return [][2]int{[2]int{int(t.proto.PosX), int(t.proto.PosY)}}
}

//GetCantSeeMeMsg 获取移除消息
func (t *ruinsBase) GetCantSeeMeMsg() pbutil.Buffer {
	return t.removeMsg
}

//AddWatcher 添加观察者(这里是为了玩家在野外视野移动中,给观察对象加一个观察者索引)
func (t *ruinsBase) AddWatcher(hc iface.HeroController) {
	t.watchList.Store(hc.Id(), hc)
}

//RemoveWatcher 移除观察者
func (t *ruinsBase) RemoveWatcher(hc iface.HeroController) {
	t.watchList.Delete(hc.Id())
}

//RangeWatchList 遍历观察者
func (t *ruinsBase) RangeWatchList(f func(hc iface.HeroController)) {
	t.watchList.Range(func(k, v interface{}) bool {
		hc, ok := v.(iface.HeroController)
		if ok {
			f(hc)
		}
		return true
	})
}

type ruinsBasePosInfoMap struct {
	sync.RWMutex
	realmId             int64
	infoMap             map[cb.Cube]*ruinsBase
	expireDuration      time.Duration // 过期间隔
	minTime             time.Time     // 最小的时间
	regionConfig        *singleton.RegionConfig
	broadcastToCaredPos func(posX, posY int, msg pbutil.Buffer, except int64)
}

func (m *ruinsBasePosInfoMap) Update(ctime time.Time) {
	if timeutil.IsZero(m.minTime) || ctime.Before(m.minTime.Add(m.expireDuration)) {
		return
	}

	var removes []*ruinsBase
	func() {
		m.Lock()
		defer m.Unlock()

		minTime := time.Time{}
		expireRuinsTime := ctime.Add(-m.expireDuration)

		for cb, info := range m.infoMap {
			// 摧毁时间
			ruinsTime := timeutil.Unix64(info.proto.Time)
			if ruinsTime.Before(expireRuinsTime) {
				// 过期了
				removes = append(removes, info)
				delete(m.infoMap, cb)
				if info.blockIndex != nil {
					info.blockIndex.RemoveRuinIndex(info)
				}
				continue
			}

			// 没过期
			if timeutil.IsZero(minTime) || ruinsTime.Before(minTime) {
				minTime = ruinsTime
			}
		}

		m.minTime = minTime
	}()

	if len(removes) > 0 {
		for _, rb := range removes {
			m.broadcastToCaredPos(int(rb.proto.PosX), int(rb.proto.PosY), rb.removeMsg, 0)
		}
	}
}

// 势力范围变更了
func (m *ruinsBasePosInfoMap) OnBaseChanged(base realmface.Base) {
	if npcid.IsHomeNpcId(base.Id()) {
		// 玩家周围的npc，不添加
		return
	}

	var targetCubes = []cb.Cube{cb.XYCube(base.BaseX(), base.BaseY())}
	if base.BaseType() == realmface.BaseTypeHome {
		evenOffsetCubes := m.regionConfig.GetEvenOffsetCubesIncludeLowLevel(base.GetBaseLevel())
		if len(evenOffsetCubes) > 0 {
			targetCubes = targetCubes[:0]

			for _, offset := range evenOffsetCubes {
				evenOffsetX, evenOffsetY := offset.XY()
				targetCube := hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), evenOffsetX, evenOffsetY)
				targetCubes = append(targetCubes, targetCube)
			}
		}
	}

	var removes []*ruinsBase
	func() {
		m.Lock()
		defer m.Unlock()
		for _, cb := range targetCubes {
			if ruinsProto := m.infoMap[cb]; ruinsProto != nil {
				delete(m.infoMap, cb)

				removes = append(removes, ruinsProto)
			}
		}
	}()

	if len(removes) > 0 {
		for _, rb := range removes {
			m.broadcastToCaredPos(int(rb.proto.PosX), int(rb.proto.PosY), rb.removeMsg, 0)
		}
	}
}

func (m *ruinsBasePosInfoMap) Walk(walkFunc func(cb.Cube, *server_proto.RuinsBasePosInfoProto)) {
	m.RLock()
	defer m.RUnlock()

	if len(m.infoMap) <= 0 {
		return
	}

	for cb, proto := range m.infoMap {
		walkFunc(cb, proto.proto)
	}
}

func (m *ruinsBasePosInfoMap) AddRuinsBase(r *Realm, heroId int64, posX, posY int, ctime time.Time) {
	m.AddRuinsBaseProto(r, &server_proto.RuinsBasePosInfoProto{
		Id:   heroId,
		Time: timeutil.Marshal64(ctime),
		PosX: int32(posX),
		PosY: int32(posY),
	}, ctime)
}

func (m *ruinsBasePosInfoMap) AddRuinsBaseProto(r *Realm, proto *server_proto.RuinsBasePosInfoProto, ctime time.Time) {
	m.Lock()
	defer m.Unlock()

	rb := newRuinsBase(proto)
	m.infoMap[cb.XYCube(int(proto.PosX), int(proto.PosY))] = rb

	//Kimi 索引废墟
	r.blockManager.AddRuinIndex(int(proto.PosX), int(proto.PosY), rb)

	if timeutil.IsZero(m.minTime) || ctime.Before(m.minTime) {
		m.minTime = ctime
	}

	// 不发消息，客户端自己根据城池流亡的消息添加
	//m.broadcastToCaredPos(int(proto.PosX), int(proto.PosY), rb.addMsg, 0)
}

func (m *ruinsBasePosInfoMap) GetRuinsBase(x, y int) int64 {
	m.RLock()
	defer m.RUnlock()

	proto := m.infoMap[cb.XYCube(x, y)]
	if proto == nil {
		return 0
	}

	return proto.proto.Id
}

func (m *ruinsBasePosInfoMap) EncodeRuinsBasePoses() (ruinsBaseX []int32, ruinsBaseY []int32) {
	m.RLock()
	defer m.RUnlock()

	if len(m.infoMap) <= 0 {
		return
	}

	ruinsBaseX = make([]int32, 0, len(m.infoMap))
	ruinsBaseY = make([]int32, 0, len(m.infoMap))

	for cb, _ := range m.infoMap {
		x, y := cb.XY()
		ruinsBaseX = append(ruinsBaseX, int32(x))
		ruinsBaseY = append(ruinsBaseY, int32(y))
	}

	return
}

func (m *ruinsBasePosInfoMap) Encode() *server_proto.RealmRuinsBasePosInfosProto {
	m.RLock()
	defer m.RUnlock()

	if len(m.infoMap) <= 0 {
		return nil
	}

	proto := &server_proto.RealmRuinsBasePosInfosProto{
		RealmId: m.realmId,
		Infos:   make([]*server_proto.RuinsBasePosInfoProto, 0, len(m.infoMap)),
	}

	for _, info := range m.infoMap {
		proto.Infos = append(proto.Infos, info.proto)
	}

	return proto
}

type basePosInfoMap struct {
	posIdMap sync.Map
}

func (m *basePosInfoMap) GetBase(x, y int) int64 {

	if v, ok := m.posIdMap.Load(cb.XYCube(x, y)); ok {
		return v.(int64)
	}
	return 0
}

// 增加base
func (m *basePosInfoMap) AddBase(base *baseWithData) {
	if npcid.IsHomeNpcId(base.Id()) {
		// 玩家周围的npc，不添加
		return
	}

	m.posIdMap.Store(cb.XYCube(base.BaseX(), base.BaseY()), base.Id())
}

func (m *basePosInfoMap) RemoveBase(base realmface.Base) {
	if npcid.IsHomeNpcId(base.Id()) {
		// 玩家周围的npc，不清理
		return
	}

	m.posIdMap.Delete(cb.XYCube(base.BaseX(), base.BaseY()))
}

func (m *basePosInfoMap) ChangeBasePos(baseId int64, oldX, oldY, newX, newY int) {
	if npcid.IsHomeNpcId(baseId) {
		// 玩家周围的npc，不清理
		return
	}

	if id := m.GetBase(oldX, oldY); id == baseId {
		m.posIdMap.Delete(cb.XYCube(oldX, oldY))
	} else {
		logrus.Errorf("迁移城池，更新坐标时候，发现原来的坐标存储的baseId，不是自己的，heroID:%d old:(%d,%d) new:(%d,%d)",
			baseId, oldX, oldY, newX, newY)
	}

	m.posIdMap.Store(cb.XYCube(newX, newY), baseId)
}
