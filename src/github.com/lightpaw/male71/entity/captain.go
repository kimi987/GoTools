package entity

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/captain"
	"time"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/u64"
	"sort"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/config/heroinit"
)

type soldier_data interface {
	GetSoldierLevel() uint64
	GetSoldierStat() *data.SpriteStat
	GetSoldierLoad() uint64
	GetSoldierModel() uint64
}

//func NewCaptain(data *captain.CaptainData, datas *config.ConfigDatas, ctime time.Time,
//	getSoldierData func(raceData *race.RaceData) soldier_data,
//	getHeroLevelData func() *herodata.HeroLevelData,
//	getTitleData func() *taskdata.TitleData) *Captain {
//	return newCaptain(data, datas, ctime, getSoldierData, getHeroLevelData, getTitleData)
//}

func newCaptain(data *captain.CaptainData, datas *heroinit.CaptainInitData,
	friendship *HeroCaptainFriendship, ctime time.Time,
	getSoldierData func(raceData *race.RaceData) soldier_data,
	getHeroLevelData func() *herodata.HeroLevelData,
	getTitleData func() *taskdata.TitleData,
	getBuff func() *data.SpriteStat) *Captain {
	p := &Captain{
		data:        data,
		starData:    data.GetFirstStar(),
		abilityData: datas.AbilityData,
		rebirthData: datas.RebirthData,
		equipment:   make(map[shared_proto.EquipmentType]*Equipment),
		official:    captain.EmptyOfficialData,
		friendship:  friendship,
		gainTime:    ctime,
	}
	p.levelData = p.rebirthData.FirstCaptainLevel
	p.gems = make([]*goods.GemData, p.rebirthData.FirstCaptainLevel.GemSlotCount)
	p.soldierData = getSoldierData(p.data.Race)
	p.getHeroLevelData = getHeroLevelData
	p.getTitleData = getTitleData
	p.getBuff = getBuff
	return p
}

func (hero *Hero) NewCaptain(data *captain.CaptainData, ctime time.Time) *Captain {
	return newCaptain(data, hero.military.captainInitData,
		hero.military.captainFriendship, ctime,
		hero.military.GetOrCreateSoldierData,
		hero.LevelData, hero.TaskList().GetTitleData,
		hero.Buff().getCaptainBuffStat)
}

type GainTimeCaptainSlice []*Captain

func (a GainTimeCaptainSlice) Len() int      { return len(a) }
func (a GainTimeCaptainSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a GainTimeCaptainSlice) Less(i, j int) bool {
	if a[i].gainTime.Equal(a[j].gainTime) {
		return a[i].Id() < a[j].Id()
	}
	return a[i].gainTime.Before(a[j].gainTime)
}

// 武将
type Captain struct {
	// 星级相关
	data     *captain.CaptainData
	starData *captain.CaptainStarData

	// 成长相关
	abilityData    *captain.CaptainAbilityData
	abilityExp     uint64
	rebirthData    *captain.CaptainRebirthLevelData
	rebirthEndTime time.Time // 转生结束时间。zero 为没在转生中

	// 等级相关
	levelData *captain.CaptainLevelData
	exp       uint64

	// 装备相关
	equipment map[shared_proto.EquipmentType]*Equipment
	taoz      *goods.EquipmentTaozData // 套装

	// 镶嵌相关
	gems []*goods.GemData // 穿戴的宝石

	// 官职相关
	gongxun       uint64                       // 功勋
	official      *captain.CaptainOfficialData // 官职
	officialIndex uint64                       // 官职所在的位置下标

	// 兵力相关
	soldier     uint64 // 当前兵力
	soldierData soldier_data

	// 武将羁绊
	friendship *HeroCaptainFriendship

	// 缓存着的属性
	totalStat      *data.SpriteStat
	totalStatProto *shared_proto.SpriteStatProto

	// 其它相关
	gainTime    time.Time // 获得时间
	trainAccExp uint64    // 训练累积经验

	// 布阵相关
	troop *Troop // 隶属于哪支布阵，没有部队则为nil

	// 属性相关
	detailStatProto  *shared_proto.SpriteStatArrayProto
	getHeroLevelData func() *herodata.HeroLevelData
	getTitleData     func() *taskdata.TitleData
	getBuff          func() *data.SpriteStat

	// 预览相关
	viewed bool
}

func (c *Captain) unmarshal(p *server_proto.HeroCaptainServerProto, datas *config.ConfigDatas, depot *Depot, buff *HeroBuff, ctime time.Time) {
	c.viewed = p.Viewed
	c.gainTime = timeutil.Unix64(p.GainTime)
	c.starData = c.data.GetStar(p.Star)
	c.abilityExp = p.AbilityExp
	c.abilityData = datas.CaptainAbilityData().Must(p.Ability)
	c.exp = p.LevelExp
	c.rebirthData = datas.CaptainRebirthLevelData().Must(p.Rebirth)
	c.rebirthEndTime = timeutil.Unix64(p.RebirthCdEndTime)
	c.levelData = c.rebirthData.MustCaptainLevelData(p.Level)
	c.gems = make([]*goods.GemData, c.levelData.GemSlotCount)
	c.gongxun = p.Gongxun

	if p.Official > 0 {
		officialCount := c.getHeroLevelData().HeroOfficialCount(p.Official)
		if p.OfficialIdx < officialCount {
			official := datas.CaptainOfficialData().Get(p.Official)
			if official != nil {
				c.official = official
				c.officialIndex = p.OfficialIdx
			}
		}
	}
	c.officialIndex = p.OfficialIdx

	c.trainAccExp = p.TrainAccExp

	for _, p := range p.Equipment {
		data := datas.GetEquipmentData(p.GetDataId())
		if data == nil {
			continue
		}
		id := depot.NewId()
		e := NewEquipment(id, data)
		e.unmarshal(p, datas)
		if old := c.equipment[data.Type]; old != nil {
			depot.AddGenIdGoods(old, ctime)
		}
		c.equipment[data.Type] = e
	}

	for _, p := range p.Gems {
		data := datas.GetGemData(p.GetGem())
		if data == nil {
			continue
		}
		slotIdx := int(p.GetSlotIdx())
		if slotIdx < 0 || slotIdx >= len(c.gems) {
			depot.AddGoods(data.Id, 1)
			continue
		}
		if old := c.gems[slotIdx]; old != nil {
			depot.AddGoods(old.Id, 1)
		}
		c.gems[slotIdx] = data
	}

	if timeutil.IsZero(c.rebirthEndTime) && c.NextLevelIsReBirth() {
		c.TriggerRebirthCD(ctime)
	}

	c.UpdateMorale(datas.EquipmentTaozConfig())

	c.CalculateProperties()

	// 设置士兵，需要获取属性种的士兵总数，所以需要先计算属性之后，才能设置士兵
	c.SetSoldier(p.Soldier)
}

func (c *Captain) Data() *captain.CaptainData {
	return c.data
}

func (c *Captain) Id() uint64 {
	return c.data.Id
}

func (c *Captain) NewTroopPos(xIndex int32) *TroopPos {
	return &TroopPos{
		captain: c,
		xIndex:  xIndex,
	}
}

func (c *Captain) Name() string {
	return c.data.Name
}

func (c *Captain) RaceId() int {
	return c.data.Race.Id
}

func (c *Captain) Race() *race.RaceData {
	return c.data.Race
}

func (c *Captain) ExchangeOfficial(captain *Captain) {
	official, index, gx := captain.official, captain.officialIndex, captain.gongxun
	captain.official, captain.officialIndex, captain.gongxun = c.official, c.officialIndex, c.gongxun
	c.official, c.officialIndex, c.gongxun = official, index, gx
}

func (c *Captain) OfficialId() uint64 {
	if c.official != nil {
		return c.official.Id
	}
	return 0
}

func (c *Captain) Official() *captain.CaptainOfficialData {
	return c.official
}

func (c *Captain) SetOfficial(official *captain.CaptainOfficialData, idx uint64) {
	c.official = official
	c.officialIndex = idx
}

func (c *Captain) OfficialIndex() uint64 {
	return c.officialIndex
}

func (c *Captain) encodeServer() *server_proto.HeroCaptainServerProto {
	p := &server_proto.HeroCaptainServerProto{}
	p.Id = c.Id()
	p.GainTime = timeutil.Marshal64(c.gainTime)
	p.LevelExp = c.exp
	p.Level = c.Level()
	p.Soldier = c.soldier
	p.Star = c.starData.Star
	p.AbilityExp = c.abilityExp
	p.Ability = c.Ability()
	p.Rebirth = c.RebirthLevel()
	p.RebirthCdEndTime = timeutil.Marshal64(c.rebirthEndTime)
	for _, e := range c.equipment {
		if e != nil {
			p.Equipment = append(p.Equipment, e.encodeServer())
		}
	}
	for idx, g := range c.gems {
		if g != nil {
			p.Gems = append(p.Gems, &server_proto.HeroCaptainGemServerProto{Gem: g.Id, SlotIdx: int32(idx)})
		}
	}
	p.TrainAccExp = c.trainAccExp
	p.Gongxun = c.gongxun
	p.Official = c.official.Id
	p.OfficialIdx = c.officialIndex
	p.Viewed = c.viewed
	return p
}

func (c *Captain) EncodeClient() *shared_proto.CaptainProto {
	p := &shared_proto.CaptainProto{}
	p.Id = u64.Int32(c.Id())
	p.IconId = c.data.Icon.Id
	p.Star = u64.Int32(c.starData.Star)
	p.Ability = u64.Int32(c.Ability())
	p.AbilityExp = u64.Int32(c.abilityExp)
	p.Gongxun = u64.Int32(c.gongxun)
	p.Rebirth = u64.Int32(c.levelData.Rebirth)
	p.RebirthCdEndTime = timeutil.Marshal32(c.rebirthEndTime)
	p.Exp = u64.Int32(c.exp)
	p.Level = u64.Int32(c.Level())
	p.Soldier = u64.Int32(c.soldier)
	p.FightAmount = u64.Int32(c.FightAmount())
	p.FullSoldierFightAmount = u64.Int32(c.FullSoldierFightAmount())
	p.TotalStat = c.totalStatProto
	p.Taoz = u64.Int32(c.TaozLevel())
	//p.RaceCdEndTime = timeutil.Marshal32(c.GetRaceCdEndTime())
	p.Equipment = c.EncodeEquipment()
	p.Gems = c.EncodeGems()
	p.Official = u64.Int32(c.official.Id)
	p.OfficialIdx = u64.Int32(c.officialIndex)
	p.Viewed = c.viewed
	return p
}

func (c *Captain) GetTrainAccExp() uint64 {
	return c.trainAccExp
}

func (c *Captain) SetTrainAccExp(toSet uint64) {
	c.trainAccExp = toSet
}

func (c *Captain) SoldierCapcity() uint64 {
	return c.totalStat.SoldierCapcity
}

func (c *Captain) GetTroop() *Troop {
	return c.troop
}

func (c *Captain) GetTroopSequence() uint64 {
	if c.troop != nil {
		return c.troop.sequence + 1
	}
	return 0
}

func (c *Captain) IsOutSide() bool {
	if c.troop != nil {
		return c.troop.invaseInfo != nil
	}
	return false
}

func (c *Captain) updateTroopChanged() {
	if c.troop != nil {
		c.troop.setChanged()
	}
}

func (c *Captain) Star() uint64 {
	return c.starData.Star
}

func (c *Captain) StarData() *captain.CaptainStarData {
	return c.starData
}

func (c *Captain) Upstar() {
	if nextStar := c.starData.GetNextStar(); nextStar != nil {
		c.starData = nextStar
	}
}

func (c *Captain) GetSpellFightAmountCoef() uint64 {
	return c.starData.GetSpellFightAmountCoef(c.abilityData.UnlockSpellCount)
}

func (c *Captain) FightAmount() uint64 {
	return c.CalFightAmountWithSoldier(c.soldier)
}

func (c *Captain) FullSoldierFightAmount() uint64 {
	return c.CalFightAmountWithSoldier(c.SoldierCapcity())
}

func (c *Captain) CalDefenseFightAmountWithSoldier(hero *Hero, soldier uint64) uint64 {
	totalStat := c.getDefenseStat(hero)
	return totalStat.FightAmount(soldier, c.GetSpellFightAmountCoef())
}

func (c *Captain) CalFightAmountWithSoldier(soldier uint64) uint64 {
	return c.totalStat.FightAmount(soldier, c.GetSpellFightAmountCoef())
}

func (c *Captain) CalFightAmountWithSoldierAndAddedStat(soldier uint64, addedStat *data.SpriteStat) uint64 {
	if addedStat != nil {
		return data.AppendSpriteStat(c.totalStat, addedStat).FightAmount(soldier, c.GetSpellFightAmountCoef())
	}
	return c.totalStat.FightAmount(soldier, c.GetSpellFightAmountCoef())
}

func (c *Captain) Soldier() uint64 {
	return c.soldier
}

func (c *Captain) AddSoldier(toAdd uint64) uint64 {
	c.SetSoldier(c.soldier + toAdd)
	return c.soldier
}

func (c *Captain) ReduceSoldier(toReduce uint64) {
	c.SetSoldier(u64.Sub(c.soldier, toReduce))
}

func (c *Captain) GetSoldierData() soldier_data {
	return c.soldierData
}

func (c *Captain) GetSoldierLevel() uint64 {
	return c.soldierData.GetSoldierLevel()
}

func (c *Captain) GetSoldierStat() *data.SpriteStat {
	return c.soldierData.GetSoldierStat()
}

func (c *Captain) SetSoldier(newSoldier uint64) bool {
	if newSoldier > c.SoldierCapcity() {
		newSoldier = c.SoldierCapcity()
	}
	isChange := c.soldier != newSoldier
	c.soldier = newSoldier
	if isChange {
		c.updateTroopChanged()
	}
	return isChange
}

func (c *Captain) GongXun() uint64 {
	return c.gongxun
}

func (c *Captain) AddGongXun(toAdd uint64) uint64 {
	c.gongxun += toAdd
	return c.gongxun
}

func (c *Captain) ExchangeAbility(captain *Captain) {
	ability, exp := captain.abilityData, captain.abilityExp
	captain.abilityData, captain.abilityExp = c.abilityData, c.abilityExp
	c.abilityData, c.abilityExp = ability, exp
}

func (c *Captain) AbilityExp() uint64 {
	return c.abilityExp
}

func (c *Captain) AddAbilityExp(toAdd uint64) {
	c.abilityExp += toAdd
}

func (c *Captain) Ability() uint64 {
	return c.abilityData.Ability
}

func (c *Captain) AbilityData() *captain.CaptainAbilityData {
	return c.abilityData
}

func (c *Captain) IsMaxAbility() bool {
	return c.abilityData.Ability >= c.levelData.AbilityLimit || c.abilityData.NextLevel() == nil
}

func (c *Captain) TryUpgradeAbility() bool {
	if c.abilityExp < c.abilityData.UpgradeExp || c.IsMaxAbility() {
		return false
	}
	c.abilityData, c.abilityExp = UpgradeAbility(c.abilityData, c.abilityExp)
	c.updateTroopChanged()
	return true
}

func UpgradeAbility(currentLevel *captain.CaptainAbilityData, originExp uint64) (newLevel *captain.CaptainAbilityData, newExp uint64) {

	newLevel = currentLevel
	newExp = originExp

	n := int(newLevel.MaxLevel)
	for i := 0; i < n; i++ {
		if newExp < newLevel.UpgradeExp {
			break
		}

		if newLevel.NextLevel() == nil {
			newExp = 0 // 最高级，删除
			break
		}

		// 扣掉当前等级的升级经验，1升2，读1级经验
		newExp = u64.Sub(newExp, newLevel.UpgradeExp)
		newLevel = newLevel.NextLevel()
	}

	return
}

func (c *Captain) Quality() shared_proto.Quality {
	return c.abilityData.Quality
}

func (c *Captain) ExchangeLevel(captain *Captain) {
	rebirth, rebirthEndTime, level, exp := captain.rebirthData, captain.rebirthEndTime, captain.levelData, captain.exp
	captain.rebirthData, captain.rebirthEndTime, captain.levelData, captain.exp = c.rebirthData, c.rebirthEndTime, c.levelData, c.exp
	c.rebirthData, c.rebirthEndTime, c.levelData, c.exp = rebirth, rebirthEndTime, level, exp
}

func (c *Captain) Level() uint64 {
	return c.levelData.Level
}

func (c *Captain) LevelData() *captain.CaptainLevelData {
	return c.levelData
}

func (c *Captain) Exp() uint64 {
	return c.exp
}

// 返回是否升级等级了
func (c *Captain) AddExp(toAdd, levelLimit uint64) bool {
	if c.IsLevelLimit(levelLimit, false) {
		return false
	}
	c.exp += toAdd
	return c.exp >= c.levelData.UpgradeExp
}

func (c *Captain) IsMaxLevel() bool {
	return c.NextRebirthData() == nil && c.levelData.NextLevel() == nil
}

func (c *Captain) IsLevelLimit(levelLimit uint64, hasRebirthLimit bool) bool {
	// 优先判断一下当前是否可以加经验
	if c.IsMaxLevel() {
		return true
	}
	if c.Level() >= levelLimit {
		return true
	}
	nextLevel := c.levelData.NextLevel()
	if nextLevel == nil || nextLevel.Rebirth > c.levelData.Rebirth {
		return true
	}
	return false
}

func (c *Captain) GetMaxCanAddExp(levelLimit uint64) (canAddExp uint64) {
	cur := c.levelData
	for i := cur.Level; i < levelLimit; i++ {
		if cur.NextLevel() == nil {
			// 防御性
			break
		}
		if c.isNextLevelIsReBirth(cur.Level) {
			break
		}
		canAddExp += cur.UpgradeExp
		cur = cur.NextLevel()
	}
	return u64.Sub(canAddExp, c.exp)
}

func (c *Captain) UpgradeLevel(levelLimit uint64, ctime time.Time) (originLevel uint64, newLevel uint64, startRebrithCD bool) {
	originLevel = c.Level()
	// 转生CD中，不能升级
	if c.IsInRebrithing(ctime) {
		return
	}
	if c.NextLevelIsReBirth() {
		return
	}
	var upgraded bool
	for i := c.levelData.Level; i < levelLimit; i++ {
		if c.levelData.NextLevel() == nil {
			// 防御性
			break
		}
		upgradeExp := c.levelData.UpgradeExp
		if c.exp < upgradeExp {
			newLevel = c.Level()
			break
		}
		c.exp = u64.Sub(c.exp, upgradeExp)
		c.levelData = c.levelData.NextLevel()
		upgraded = true
		// 判定是否有宝石槽扩张
		if c.levelData.HasNewGemSlot {
			gems := make([]*goods.GemData, c.levelData.GemSlotCount)
			copy(gems, c.gems)
			c.gems = gems
		}
		if c.NextLevelIsReBirth() {
			break
		}
	}
	if upgraded {
		// 升级成功
		newLevel = c.Level()
		c.updateTroopChanged()
		if c.NextLevelIsReBirth() {
			// 触发转生 CD
			c.TriggerRebirthCD(ctime)
			startRebrithCD = true
		}
	}
	return
}

func (c *Captain) NextRebirthData() *captain.CaptainRebirthLevelData {
	return c.rebirthData.GetNextLevel()
}

func (c *Captain) RebirthLevel() uint64 {
	return c.rebirthData.Level
}

func (c *Captain) RebirthLevelData() *captain.CaptainRebirthLevelData {
	return c.rebirthData
}

func (c *Captain) RebirthEndTime() time.Time {
	return c.rebirthEndTime
}

func (c *Captain) SetRebirthEndTime(endTime time.Time) {
	c.rebirthEndTime = endTime
}

func (c *Captain) ExchangeEquipments(captain *Captain) {
	equip, taoz := captain.equipment, captain.taoz
	captain.equipment, captain.taoz = c.equipment, c.taoz
	c.equipment, c.taoz = equip, taoz
}

func (c *Captain) SetEquipment(e *Equipment) *Equipment {

	old := c.equipment[e.data.Type]
	c.equipment[e.data.Type] = e

	return old
}

func (c *Captain) GetEquipment(t shared_proto.EquipmentType) *Equipment {
	return c.equipment[t]
}

func (c *Captain) GetEquipmentById(id uint64) *Equipment {

	for _, v := range c.equipment {
		if v != nil && v.id == id {
			return v
		}
	}
	return nil
}

func (c *Captain) RemoveEquipment(id uint64) *Equipment {

	for k, v := range c.equipment {
		if v != nil && v.id == id {
			delete(c.equipment, k)
			return v
		}
	}
	return nil
}

func (c *Captain) RemoveAllEquipment() (removedEquipments []*Equipment) {
	if len(c.equipment) <= 0 {
		return
	}
	for _, equipment := range c.equipment {
		if equipment != nil {
			removedEquipments = append(removedEquipments, equipment)
		}
	}
	c.equipment = make(map[shared_proto.EquipmentType]*Equipment)
	return
}

func (c *Captain) GetEquipmentCount() uint64 {
	return u64.FromInt(len(c.equipment))
}

func (c *Captain) WalkEquipment(walkFunc func(e *Equipment) (walkEnd bool)) {
	for _, e := range c.equipment {
		if walkFunc(e) {
			return
		}
	}
}

func (c *Captain) RebirthAbilityLimit() uint64 {
	return c.rebirthData.AbilityLimit
}

func (c *Captain) NextLevelIsReBirth() bool {
	return c.isNextLevelIsReBirth(c.Level())
}

func (c *Captain) isNextLevelIsReBirth(level uint64) bool {
	if c.rebirthData.GetNextLevel() == nil {
		return false
	}
	return level >= c.rebirthData.GetNextLevel().CaptainLevelLimit.Level
}

func (c *Captain) IsUpgrade(levelLimit uint64) bool {

	if c.exp < c.levelData.UpgradeExp || c.IsLevelLimit(levelLimit, false) {
		return false
	}

	return true
}

// 是否正在转生CD中
func (c *Captain) IsInRebrithing(ctime time.Time) bool {
	return c.rebirthEndTime.After(ctime)
}

// 触发转生CD
func (c *Captain) TriggerRebirthCD(ctime time.Time) bool {
	if !timeutil.IsZero(c.rebirthEndTime) {
		logrus.Debugf("达到转生等级，却没有触发转生cd %s %s", c.RebirthLevel(), c.Level())
		return false
	}
	if c.rebirthData.GetNextLevel() == nil {
		return false
	}
	c.rebirthEndTime = ctime.Add(c.rebirthData.GetNextLevel().Cd)
	return true
}

// 执行转生
// 改变的属性在外面处理
func (c *Captain) ProcessRebirth(ctime time.Time) bool {
	if c.IsInRebrithing(ctime) {
		return false
	}
	if timeutil.IsZero(c.rebirthEndTime) {
		return false
	}
	if c.rebirthData.GetNextLevel() == nil {
		return false
	}
	c.rebirthData = c.rebirthData.GetNextLevel()
	c.rebirthEndTime = time.Time{}
	c.levelData = c.rebirthData.FirstCaptainLevel
	c.updateTroopChanged()
	c.exp = 0
	return true
}

func (c *Captain) GemSlotCapacity() int {
	return len(c.gems)
}

func (c *Captain) GetGem(idx int) *goods.GemData {
	if idx >= 0 && idx < len(c.gems) {
		return c.gems[idx]
	}
	return nil
}

func (c *Captain) SetGem(slotIdx int, g *goods.GemData) (old *goods.GemData) {
	old = c.gems[slotIdx]
	c.gems[slotIdx] = g
	return
}

func (c *Captain) RemoveGem(slotIdx int) *goods.GemData {
	gem := c.gems[slotIdx]
	c.gems[slotIdx] = nil
	return gem
}

func (c *Captain) Gems() []*goods.GemData {
	return c.gems
}

func (c *Captain) CopyGems() []*goods.GemData {
	result := make([]*goods.GemData, len(c.gems))
	copy(result, c.gems)
	return result
}

func (c *Captain) RemoveAllGem() []*goods.GemData {

	array := make([]*goods.GemData, 0, len(c.gems))
	for _, v := range c.gems {
		if v != nil {
			array = append(array, v)
		}
	}
	c.gems = make([]*goods.GemData, c.levelData.GemSlotCount)
	return array
}

func (c *Captain) GetGemCount() (count uint64) {
	for _, v := range c.gems {
		if v != nil {
			count++
		}
	}

	return count
}

func (c *Captain) GemStat() *data.SpriteStat {
	b := data.NewSpriteStatBuilder()
	for _, g := range c.gems {
		if g != nil {
			b.Add(g.BaseStat)
		}
	}
	return b.Build()
}

// 装备属性
func (c *Captain) EquipmentStat() *data.SpriteStat {

	b := data.NewSpriteStatBuilder()

	for _, e := range c.equipment {
		b.AddProto(e.TotalStat())
	}

	if c.taoz != nil {
		b.Add(c.taoz.SpriteStat)
	}

	return b.Build()
}

func (c *Captain) EquipmentRankLevels() []uint64 {

	if len(c.equipment) <= 0 {
		return nil
	}

	levels := make([]uint64, 0, len(c.equipment))
	for _, e := range c.equipment {
		levels = append(levels, e.RefinedLevel())
	}

	// 等级从低到高排序
	sort.Sort(u64.Uint64Slice(levels))
	return levels
}

func (c *Captain) TaozData() *goods.EquipmentTaozData {
	return c.taoz
}

func (c *Captain) TaozStar() uint64 {
	if c.taoz != nil {
		return c.taoz.Star
	}
	return 0
}

func (c *Captain) TaozLevel() uint64 {
	if c.taoz != nil {
		return c.taoz.Level
	}
	return 0
}

func (c *Captain) Morale() uint64 {
	morale := uint64(0)
	if c.taoz != nil {
		morale += c.taoz.Morale
	}

	return morale
}

func (c *Captain) UpdateMorale(config *goods.EquipmentTaozConfig) {

	refinedLevels := c.EquipmentRankLevels()
	n := len(refinedLevels)

	// 看下最低件数要求达到没有
	c.taoz = nil
	for i := 0; i < n; i++ {
		level := refinedLevels[i]
		if level <= 0 {
			continue
		}

		count := u64.FromInt(n - i)
		taozData := config.GetTaoz(count, level)
		if taozData == nil {
			continue
		}

		c.taoz = taozData
		break
	}
}

func (c *Captain) GetRarityStarCoef() float64 {
	return c.starData.Coef
}

func (c *Captain) CalculateProperties() {
	// 武将四维属性之一 =（成长*20*该属性比例+该属性精研值）*武将等级+ 转生属性*该属性比例
	captainStat := CalculateCaptainStat(c.Ability(), c.Level(), c.rebirthData, c.data.Race, c.GetSoldierStat(), c.StarData())
	cfRaceStat := c.friendship.getRaceStat(c.Race().Race)
	equipmentStat := c.EquipmentStat()
	gemStat := c.GemStat()
	officialStat := c.official.SpriteStat
	buffStat := c.getBuff()
	c.totalStat = CalculateCaptainTotalStat(captainStat, cfRaceStat, equipmentStat, gemStat, officialStat, buffStat, c.rebirthData, c.levelData, c.getHeroLevelData(), c.getTitleData())
	c.detailStatProto = toDetailStatProto2(captainStat, equipmentStat, gemStat, officialStat)

	c.totalStatProto = c.totalStat.Encode()
	c.updateTroopChanged()
}

func (c *Captain) GetTotalStat() *shared_proto.SpriteStatProto {
	return c.totalStatProto
}

func (c *Captain) GetDetailStat() *shared_proto.SpriteStatArrayProto {
	return c.detailStatProto
}

func (c *Captain) GetPreviewStat(titleData *taskdata.TitleData) *data.SpriteStat {

	// 武将四维属性之一 =（成长*20*该属性比例+该属性精研值）*武将等级+ 转生属性*该属性比例
	captainStat := CalculateCaptainStat(c.Ability(), c.Level(), c.rebirthData, c.data.Race, c.GetSoldierStat(), c.StarData())
	cfRaceStat := c.friendship.getRaceStat(c.Race().Race)
	equipmentStat := c.EquipmentStat()
	gemStat := c.GemStat()
	officialStat := c.official.SpriteStat
	buffStat := c.getBuff()

	return CalculateCaptainTotalStat(captainStat, cfRaceStat, equipmentStat, gemStat,
		officialStat, buffStat, c.rebirthData, c.levelData, c.getHeroLevelData(), titleData)
}

func (c *Captain) NewUpdateCaptainStatMsg() pbutil.Buffer {
	return military.NewS2cUpdateCaptainStatMarshalMsg(u64.Int32(c.Id()), c.totalStatProto, u64.Int32(c.FightAmount()), u64.Int32(c.FullSoldierFightAmount()))
}

func toDetailStatProto2(captainStat, equipmentStat, gemStat, officialStat *data.SpriteStat) *shared_proto.SpriteStatArrayProto {
	// 属性列表，依次为 武将属性，装备属性，宝石属性, 官职属性

	captainStat = data.Nil2EmptyStat(captainStat)
	equipmentStat = data.Nil2EmptyStat(equipmentStat)
	gemStat = data.Nil2EmptyStat(gemStat)
	officialStat = data.Nil2EmptyStat(officialStat)

	dsp := &shared_proto.SpriteStatArrayProto{}
	dsp.SpriteStat = append(dsp.SpriteStat, captainStat.Encode())
	dsp.SpriteStat = append(dsp.SpriteStat, equipmentStat.Encode())
	dsp.SpriteStat = append(dsp.SpriteStat, gemStat.Encode())
	dsp.SpriteStat = append(dsp.SpriteStat, officialStat.Encode())

	return dsp
}

func (c *Captain) EncodeGems() []*shared_proto.HeroCaptainGemProto {
	p := make([]*shared_proto.HeroCaptainGemProto, 0, len(c.gems))
	for idx, g := range c.gems {
		if g != nil {
			p = append(p, &shared_proto.HeroCaptainGemProto{Gem: u64.Int32(g.Id), SlotIdx: int32(idx)})
		}
	}
	return p
}

func (c *Captain) EncodeEquipment() []*shared_proto.EquipmentProto {
	p := make([]*shared_proto.EquipmentProto, 0, len(c.equipment))
	for _, e := range c.equipment {
		if e != nil {
			p = append(p, e.EncodeClient())
		}
	}
	return p
}

func (captain *Captain) EncodeOtherView() *shared_proto.HeroCaptainOtherProto {

	out := &shared_proto.HeroCaptainOtherProto{}
	out.Id = u64.Int32(captain.Id())
	out.Quality = captain.Quality()
	out.Ability = u64.Int32(captain.Ability())

	out.Rebirth = u64.Int32(captain.levelData.Rebirth)

	out.Level = u64.Int32(captain.Level())
	out.SoldierLevel = u64.Int32(captain.GetSoldierLevel())

	out.FullSoldierFightAmount = u64.Int32(captain.FullSoldierFightAmount())
	out.TotalStat = captain.totalStatProto

	out.Equipment = captain.EncodeEquipment()

	out.Gems = captain.EncodeGems()

	out.Official = u64.Int32(captain.Official().Id)
	out.Star = u64.Int32(captain.Star())

	return out
}

// pve 直接使用 captain.GetTotalStat()

// 野外出击
func (c *Captain) getInvaseStat(hero *Hero) *data.SpriteStat {
	return c.getCaptainStat(hero, false, false, false, true)
}

// 野外防守
func (c *Captain) getDefenseStat(hero *Hero) *data.SpriteStat {
	return c.getCaptainStat(hero, true, false, false, true)
}

// 野外援助
func (c *Captain) getAssistStat(hero *Hero) *data.SpriteStat {
	return c.getCaptainStat(hero, false, true, false, true)
}

// 镜像防守
func (c *Captain) GetCopyDefenseStat(hero *Hero) *data.SpriteStat {
	return c.getCaptainStat(hero, true, false, true, true)
}

func (c *Captain) getCaptainStat(hero *Hero, hasDefenseStat, hasAssistStat, hasCopyDefenseStat, hasPvpStat bool) *data.SpriteStat {

	if !hasDefenseStat && !hasAssistStat && !hasCopyDefenseStat && !hasPvpStat {
		return c.totalStat
	}

	b := data.NewSpriteStatBuilder()
	b.Add(c.totalStat)

	if hasDefenseStat || hasAssistStat || hasCopyDefenseStat {
		if bd := hero.Domestic().GetBuilding(shared_proto.BuildingType_CHENG_QIANG); bd != nil && bd.Effect != nil {
			if hasDefenseStat {
				b.Add(bd.Effect.AddedDefenseStat)
			}

			if hasAssistStat {
				b.Add(bd.Effect.AddedAssistStat)
			}

			if hasCopyDefenseStat {
				b.Add(bd.Effect.AddedCopyDefenseStat)
			}
		}
	}

	if hasPvpStat {
		b.Add(hero.Buff().getPvpBuffStat())
	}

	return b.Build()
}

//func (c *Captain) EncodeCaptainInfo(fullSoldier bool, xIndex int32) *shared_proto.CaptainInfoProto {
//	totalSoldier := c.SoldierCapcity()
//	soldier := totalSoldier
//	if !fullSoldier {
//		soldier = c.Soldier()
//	}
//
//	return c.EncodeCaptainInfoWithSoldier(soldier, totalSoldier, xIndex)
//}
//
//func (c *Captain) EncodeCaptainInfoWithBuff(buff *data.StatBuff, fullSoldier bool, xIndex int32) (cp *shared_proto.CaptainInfoProto) {
//	cp = c.EncodeCaptainInfo(fullSoldier, xIndex)
//	if buff != nil {
//		stat := data.CalcStatProtoBuff(c.GetTotalStat(), buff)
//		cp.TotalStat = stat
//		cp.FightAmount = u64.Int32(data.CalFightAmount(u64.FromInt32(stat.Attack), u64.FromInt32(stat.Defense), u64.FromInt32(stat.Strength), u64.FromInt32(stat.Dexterity), int64(stat.DamageIncrePer), int64(stat.DamageDecrePer), u64.FromInt32(cp.Soldier)))
//	}
//
//	return
//}

func (c *Captain) EncodeInvaseCaptainInfo(hero *Hero, fullSoldier bool, xIndex int32) *shared_proto.CaptainInfoProto {
	return c.encodeCaptainInfo(c.getInvaseStat(hero).Encode(), fullSoldier, xIndex)
}

func (c *Captain) EncodeDefenseCaptainInfo(hero *Hero, fullSoldier bool, xIndex int32) *shared_proto.CaptainInfoProto {
	return c.encodeCaptainInfo(c.getDefenseStat(hero).Encode(), fullSoldier, xIndex)
}

func (c *Captain) EncodeCopyDefenseCaptainInfo(hero *Hero, fullSoldier bool, xIndex int32) *shared_proto.CaptainInfoProto {
	return c.encodeCaptainInfo(c.GetCopyDefenseStat(hero).Encode(), fullSoldier, xIndex)
}

func (c *Captain) EncodeAssistCaptainInfo(hero *Hero, fullSoldier bool, xIndex int32) *shared_proto.CaptainInfoProto {
	return c.encodeCaptainInfo(c.getAssistStat(hero).Encode(), fullSoldier, xIndex)
}

func (c *Captain) EncodeCaptainInfo(fullSoldier bool, xIndex int32) *shared_proto.CaptainInfoProto {
	soldier, totalSoldier := c.getSoldier(fullSoldier)
	return c.EncodeCaptainInfoWithSoldier(c.totalStatProto, soldier, totalSoldier, xIndex)
}

func (c *Captain) encodeCaptainInfo(totalStat *shared_proto.SpriteStatProto, fullSoldier bool, xIndex int32) *shared_proto.CaptainInfoProto {
	soldier, totalSoldier := c.getSoldier(fullSoldier)
	return c.EncodeCaptainInfoWithSoldier(totalStat, soldier, totalSoldier, xIndex)
}

func (c *Captain) getSoldier(fullSoldier bool) (soldier, totalSoldier uint64) {
	totalSoldier = c.SoldierCapcity()
	soldier = totalSoldier
	if !fullSoldier {
		soldier = c.Soldier()
	}
	return
}

func (c *Captain) EncodeCaptainInfoWithSoldier(totalStat *shared_proto.SpriteStatProto, soldier, totalSoldier uint64, xIndex int32) *shared_proto.CaptainInfoProto {
	soldier = u64.Min(soldier, totalSoldier)

	out := &shared_proto.CaptainInfoProto{}

	out.Id = u64.Int32(c.Id())
	//out.Name = c.Name()
	//out.IconId = c.icon.Id
	//out.Male = c.male
	out.Race = c.data.Race.Race // 老版战斗需要
	out.Quality = c.abilityData.Quality
	out.Soldier = u64.Int32(soldier)
	out.SpellFightAmountCoef = u64.Int32(c.GetSpellFightAmountCoef())
	out.FightAmount = data.ProtoFightAmount(totalStat, out.Soldier, out.SpellFightAmountCoef)
	if soldier >= totalSoldier {
		out.TotalSoldier = out.Soldier
		out.FullFightAmount = out.FightAmount
	} else {
		out.TotalSoldier = u64.Int32(totalSoldier)
		out.FullFightAmount = data.ProtoFightAmount(totalStat, out.TotalSoldier, out.SpellFightAmountCoef)
	}

	out.Level = u64.Int32(c.Level())
	out.SoldierLevel = u64.Int32(c.GetSoldierLevel())

	out.Outside = c.IsOutSide()
	out.TotalStat = totalStat
	out.LifePerSoldier = u64.Int32(data.ProtoLife(totalStat))

	out.Model = u64.Int32(c.GetSoldierData().GetSoldierModel())

	out.RebirthLevel = u64.Int32(c.RebirthLevel())

	// 技能列表
	out.XIndex = xIndex

	out.CaptainId = out.Id

	out.Star = u64.Int32(c.Star())
	out.UnlockSpellCount = u64.Int32(c.abilityData.UnlockSpellCount)

	out.CanTriggerRestraintSpell = true

	return out
}

func toDetailStatProto(captainStat, equipmentStat, gemStat, fushenStat, officialStat *data.SpriteStat) *shared_proto.SpriteStatArrayProto {
	// 属性列表，依次为 武将属性，装备属性，宝石属性，将魂属性, 官职属性

	captainStat = data.Nil2EmptyStat(captainStat)
	equipmentStat = data.Nil2EmptyStat(equipmentStat)
	gemStat = data.Nil2EmptyStat(gemStat)
	fushenStat = data.Nil2EmptyStat(fushenStat)
	officialStat = data.Nil2EmptyStat(officialStat)

	dsp := &shared_proto.SpriteStatArrayProto{}
	dsp.SpriteStat = append(dsp.SpriteStat, captainStat.Encode())
	dsp.SpriteStat = append(dsp.SpriteStat, equipmentStat.Encode())
	dsp.SpriteStat = append(dsp.SpriteStat, gemStat.Encode())
	dsp.SpriteStat = append(dsp.SpriteStat, fushenStat.Encode())
	dsp.SpriteStat = append(dsp.SpriteStat, officialStat.Encode())

	return dsp
}

func CalculateCaptainStat(ability, level uint64, rebirthData *captain.CaptainRebirthLevelData, raceData *race.RaceData, soldierStat *data.SpriteStat, starData *captain.CaptainStarData) *data.SpriteStat {
	captainStatPoint := CalculateCaptainStatPoint(ability, level, rebirthData, starData.Coef)

	b := data.NewSpriteStatBuilder()
	if soldierStat != nil {
		b.Add(soldierStat)
	}

	if captainStatPoint > 0 {
		b.Add4D(CalculateCaptain4DStatPoint(captainStatPoint, raceData))
	}

	if starData.SpriteStat != nil {
		b.Add(starData.SpriteStat)
	}

	return b.Build()
}

func CalculateCaptain4DStatPoint(captainStatPoint uint64, raceData *race.RaceData) (attack, defense, strength, dexterity uint64) {

	attack = u64.MultiCoef(captainStatPoint, raceData.GetAbilityRate(shared_proto.StatType_ATTACK))
	defense = u64.MultiCoef(captainStatPoint, raceData.GetAbilityRate(shared_proto.StatType_DEFENSE))
	strength = u64.MultiCoef(captainStatPoint, raceData.GetAbilityRate(shared_proto.StatType_STRENGTH))
	dexterity = u64.MultiCoef(captainStatPoint, raceData.GetAbilityRate(shared_proto.StatType_DEXTERITY))
	return
}

func CalculateCaptainStatPoint(ability, level uint64, rebirthData *captain.CaptainRebirthLevelData, rarityStarCoef float64) uint64 {
	// 武将四维属性之一 =（成长)*(武将等级+30)*星级系数+ 转生属性）*该属性比例
	point := ability * (level + rebirthData.BeforeRebirthLevel + 30)
	if rarityStarCoef > 0 {
		point = u64.MultiF64(point, rarityStarCoef)
	}
	return point + rebirthData.SpriteStatPoint
}

func CalculateCaptainTotalStat(captainStat, cfRaceStat, equipmentStat, gemStat, officialStat, buffStat *data.SpriteStat,
	rebirthData *captain.CaptainRebirthLevelData, levelData *captain.CaptainLevelData,
	heroLevelData *herodata.HeroLevelData, titleData *taskdata.TitleData) *data.SpriteStat {

	b := data.NewSpriteStatBuilder()

	if equipmentStat != nil {
		b.Add(equipmentStat)
	}

	if gemStat != nil {
		b.Add(gemStat)
	}

	if captainStat != nil {
		b.Add(captainStat)
	}

	if cfRaceStat != nil {
		b.Add(cfRaceStat)
	}

	if officialStat != nil {
		b.Add(officialStat)
	}

	if buffStat != nil {
		b.Add(buffStat)
	}

	b.AddSoldierCapcity(levelData.SoldierCapcity)
	b.AddSoldierCapcity(rebirthData.SoldierCapcity)
	b.AddSoldierCapcity(heroLevelData.Sub.AddSoldierCapacity)

	if titleData != nil {
		b.Add(titleData.GetTotalStat())
	}

	return b.Build()
}

func (c *Captain) EncodeHebiCaptain() *shared_proto.HebiCaptainProto {
	out := &shared_proto.HebiCaptainProto{}
	out.Id = u64.Int32(c.Id())
	out.Quality = c.abilityData.Quality
	out.FightAmount = u64.Int32(c.FullSoldierFightAmount())
	out.Level = u64.Int32(c.Level())
	out.Race = c.Race().Race

	return out
}

func (c *Captain) TrySetViewedTrue() bool {
	if c.viewed {
		return false
	}
	c.viewed = true
	return true
}
