package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"math"
	"time"
	"github.com/lightpaw/male7/config/domestic_data/sub"
)

func newDomestic(initData *heroinit.HeroInitData, buildingEffect *building_effect, ctime time.Time) *HeroDomestic {
	d := &HeroDomestic{
		buildingEffect:             buildingEffect,
		buildings:                  make(map[shared_proto.BuildingType]*domestic_data.BuildingData),
		resourcePoint:              make(map[uint64]*ResourcePoint),
		technologys:                make(map[uint64]*domestic_data.TechnologyData),
		workerRestEndTime:          make([]time.Time, initData.BuildingWorkerMaxCount),
		workerAlwaysUnlocked:       make([]bool, initData.BuildingWorkerMaxCount),
		workerLockStartTime:        make([]time.Time, initData.BuildingWorkerMaxCount),
		workerSeekHelp:             make([]bool, initData.BuildingWorkerMaxCount),
		checkRemoveWorkerSeekHelp:  make([]bool, initData.BuildingWorkerMaxCount),
		techWorkerRestEndTime:      make([]time.Time, initData.TechWorkerMaxCount),
		techSeekHelp:               make([]bool, initData.TechWorkerMaxCount),
		checkRemoveTechSeekHelp:    make([]bool, initData.TechWorkerMaxCount),
		forgingTimes:               newNextTimes(),
		cityEvent:                  &CityEvent{},
		resourceNextCollectTimeMap: make(map[shared_proto.ResType]time.Time),
		outerCities:                newOuterCities(),
		countdownPrize:             newCountdownPrize(initData.FirstLevelCountdownPrize, ctime),
	}

	d.workerAlwaysUnlocked[0] = true
	d.buildingInitEffect = initData.BuildingInitEffect

	for _, b := range initData.Building {
		d.buildings[b.Type] = b
	}

	return d
}

func (d *HeroDomestic) unmarshal(heroId int64, heroName string, p *server_proto.HeroDomesticServerProto, datas *config.ConfigDatas, ctime time.Time) {

	if p == nil {
		return
	}

	timeutil.CopyUnix32Array(d.workerRestEndTime, p.GetWorkerRestEndTime())
	timeutil.CopyUnix32Array(d.techWorkerRestEndTime, p.GetTechnologyRestEndTime())

	timeutil.CopyUnix32Array(d.workerLockStartTime, p.WorkerLockStartTime)
	copy(d.workerAlwaysUnlocked, p.WorkerAlwaysUnlocked)

	copy(d.workerSeekHelp, p.WorkerSeekHelp)
	copy(d.techSeekHelp, p.TechnologySeekHelp)
	d.dailyHelpMemberTimes = p.DailyHelpMemberTimes
	copy(d.checkRemoveWorkerSeekHelp, p.CheckRemoveWorkerSeekHelp)
	copy(d.checkRemoveTechSeekHelp, p.CheckRemoveTechnologySeekHelp)

	for _, rpp := range p.ResourcePoint {
		layoutData := datas.GetBuildingLayoutData(rpp.LayoutId)
		building := datas.GetBuildingData(rpp.BuildingId)
		if layoutData == nil || building == nil || !building.IsResPoint() {
			continue
		}

		d.resourcePoint[rpp.LayoutId] = NewResourcePoint(layoutData, building,
			timeutil.Unix64(rpp.OutputStartTime))
	}

	d.conflictResourcePoints = p.ConflictResourcePoints

	for _, bid := range p.Building {
		building := datas.GetBuildingData(bid)
		if building == nil {
			continue
		}

		d.buildings[building.Type] = building
	}

	for _, tid := range p.Technology {
		tech := datas.GetTechnologyData(tid)
		if tech == nil {
			continue
		}

		d.technologys[tech.Group] = tech
	}

	for _, tid := range p.GuildTechnology {
		tech := datas.GetGuildTechnologyData(tid)
		if tech == nil {
			continue
		}

		d.guildTechnologys = append(d.guildTechnologys, tech)
	}

	d.forgingTimes.unmarshal(p.ForgingCombineTimes)
	d.newForgingPos = p.NewForgingPos

	d.workshopIndex = p.WorkshopIndex
	d.workshopCollectTime = timeutil.Unix64(p.WorkshopCollectTime)
	d.nextRefreshWorkshopTime = timeutil.Unix64(p.NextRefreshWorkshopTime)
	d.workshopRefreshTimes = p.WorkshopRefreshTimes
	for _, id := range p.WorkshopEquipmentIds {
		data := datas.GetEquipmentData(id)
		if data != nil {
			d.workshopEquipment = append(d.workshopEquipment, data)
		} else {
			d.workshopIndex = 0
			d.workshopCollectTime = time.Time{}
		}
	}

	// 检查有效性
	if d.workshopIndex > uint64(len(d.workshopEquipment)) {
		d.workshopIndex = 0
		d.workshopCollectTime = time.Time{}
	}

	d.dailyResourceCollectTimes = p.DailyResourceCollectTimes
	if n := imath.Min(len(p.NextCollectTimeType), len(p.NextCollectTime)); n > 0 {
		for i := 0; i < n; i++ {
			t := timeutil.Unix64(p.NextCollectTime[i])
			if ctime.Before(t) {
				d.resourceNextCollectTimeMap[p.NextCollectTimeType[i]] = t
			}
		}
	}

	d.cityEvent.unmarshal(p.CityEvent, datas)

	d.sign = p.Sign

	d.voice = p.Voice

	d.outerCities.unmarshal(p.Cities, datas)

	d.countdownPrize.unmarshal(p.CountdownPrize, datas)

	d.isCollectSeasonPrize = p.IsCollectSeasonPrize
	d.collectProsperityDownCountDownPrizeTimes = p.CollectProsperityDownCountDownPrizeTimes
	d.nextCanCollectProsperityDownCountDownPrizeTime = timeutil.Unix64(p.NextCanCollectProsperityDownCountDownPrizeTime)
}

// 内政
type HeroDomestic struct {
	buildingEffect *building_effect

	buildings   map[shared_proto.BuildingType]*domestic_data.BuildingData // 城内建筑
	technologys map[uint64]*domestic_data.TechnologyData                  // key 是group
	outerCities *OuterCities                                              // 外城

	guildTechnologys []*guild_data.GuildTechnologyData

	resourcePoint          map[uint64]*ResourcePoint
	conflictResourcePoints []uint64 // 冲突的layoutid

	workerRestEndTime         []time.Time // 建筑队
	workerSeekHelp            []bool      // 建筑队求助
	checkRemoveWorkerSeekHelp []bool      // 移除建筑队求助
	workerAlwaysUnlocked      []bool      // 建筑队永久解锁
	workerLockStartTime       []time.Time // 建筑队解锁结束重新锁住的时间

	techWorkerRestEndTime   []time.Time // 科研队
	techSeekHelp            []bool      // 科研队求助
	checkRemoveTechSeekHelp []bool      // 移除科研队求助

	// 今日帮助盟友次数
	dailyHelpMemberTimes uint64

	forgingTimes  *nexttimes
	newForgingPos []uint64

	// 装备作坊（新版锻造）
	workshopEquipment    []*goods.EquipmentData
	workshopIndex        uint64 // 从1开始
	workshopCollectTime  time.Time
	workshopRefreshTimes uint64 // 装备作坊刷新次数

	nextRefreshWorkshopTime time.Time

	sign  string // 签名
	voice []byte // 语音

	cityEvent *CityEvent

	dailyResourceCollectTimes  uint64 // 资源征收次数
	resourceNextCollectTimeMap map[shared_proto.ResType]time.Time

	// 倒计时礼包

	isCollectSeasonPrize bool // 是否领取了季节奖励

	countdownPrize                                 *countdown_prize
	isResetCountdownPrize                          bool
	collectProsperityDownCountDownPrizeTimes       uint64    // 领取了繁荣度降低的奖励的次数
	nextCanCollectProsperityDownCountDownPrizeTime time.Time // 下次可以领取繁荣度降低的奖励的时间

	buildingInitEffect *sub.BuildingEffectData // 建筑初始加成效果
}

func (m *HeroDomestic) GetGuildTechnology() []*guild_data.GuildTechnologyData {
	return m.guildTechnologys
}

func (m *HeroDomestic) SetGuildTechnology(toSet []*guild_data.GuildTechnologyData) []*guild_data.GuildTechnologyData {
	old := m.guildTechnologys
	m.guildTechnologys = toSet
	return old
}

func (m *HeroDomestic) GetNextRefreshWorkshopTime() time.Time {
	return m.nextRefreshWorkshopTime
}

func (m *HeroDomestic) RefreshWorkshop(nextTime time.Time, newEquipmets []*goods.EquipmentData) {
	m.nextRefreshWorkshopTime = nextTime

	// 如果当前有正在进行的装备，则将这部分装备移动到队伍头部
	var keep *goods.EquipmentData
	if m.workshopIndex > 0 {
		idx := u64.Sub(m.workshopIndex, 1)
		if idx < uint64(len(m.workshopEquipment)) {
			keep = m.workshopEquipment[idx]
		}
	}

	m.workshopEquipment = newEquipmets

	if keep != nil {
		m.workshopEquipment = append(m.workshopEquipment, keep)
		m.workshopIndex = uint64(len(m.workshopEquipment))
	} else {
		m.SetCurrentWorkshop(0, time.Time{})
	}
}

func (m *HeroDomestic) SetCurrentWorkshop(workshopIndex uint64, collectTime time.Time) {
	m.workshopIndex = workshopIndex
	m.workshopCollectTime = collectTime
}

func (m *HeroDomestic) GetWorkshopIndex() uint64 {
	return m.workshopIndex
}

func (m *HeroDomestic) SetWorkshopCollectTime(time time.Time) {
	m.workshopCollectTime = time
}

func (m *HeroDomestic) GetWorkshopCollectTime() time.Time {
	return m.workshopCollectTime
}

func (m *HeroDomestic) GetWorkshopEquipment() []*goods.EquipmentData {
	return m.workshopEquipment
}

// 能否刷新装备作坊
func (m *HeroDomestic) CanWorkshopRefresh() bool {
	for idx := range m.workshopEquipment {
		if uint64(idx+1) == m.workshopIndex {
			continue
		}
		return false
	}
	return true
}

func (m *HeroDomestic) GetWorkshopRefreshTimes() uint64 {
	return m.workshopRefreshTimes
}

func (m *HeroDomestic) IncWorkshopRefreshTimes() {
	m.workshopRefreshTimes++
}

func (m *HeroDomestic) CollectWorkshop() (e *goods.EquipmentData) {
	// 先将这个index的equipment移除掉
	if m.workshopIndex > 0 {
		idx := u64.Sub(m.workshopIndex, 1)
		n := uint64(len(m.workshopEquipment))
		if idx < n {
			e = m.workshopEquipment[idx]
			for i := idx + 1; i < n; i++ {
				m.workshopEquipment[i-1] = m.workshopEquipment[i]
			}
			m.workshopEquipment = m.workshopEquipment[:n-1]
		}

		m.SetCurrentWorkshop(0, time.Time{})
	}
	return
}

func (m *HeroDomestic) IsCollectSeasonPrize() bool {
	return m.isCollectSeasonPrize
}

func (m *HeroDomestic) SetCollectSeasonPrize() {
	m.isCollectSeasonPrize = true
}

func (m *HeroDomestic) NextCanCollectProsperityDownCountDownPrizeTime() time.Time {
	return m.nextCanCollectProsperityDownCountDownPrizeTime
}

func (m *HeroDomestic) SetNextCanCollectProsperityDownCountDownPrizeTime(toSet time.Time) {
	m.nextCanCollectProsperityDownCountDownPrizeTime = toSet
}

func (m *HeroDomestic) CollectProsperityDownCountDownPrizeTimes() uint64 {
	return m.collectProsperityDownCountDownPrizeTimes
}

func (m *HeroDomestic) IncCollectProsperityDownCountDownPrizeTimes() {
	m.collectProsperityDownCountDownPrizeTimes++
}

func (m *HeroDomestic) GetCountdownPrize() *countdown_prize {
	return m.countdownPrize
}

func (m *HeroDomestic) CollectCountdownPrize(ctime time.Time) *countdown_prize {

	data := m.countdownPrize.data.NextData()
	if m.isResetCountdownPrize {
		m.isResetCountdownPrize = false
		if m.countdownPrize.data != m.countdownPrize.data.FirstData() {
			data = m.countdownPrize.data.FirstData()
		}
	}

	if data != nil {
		m.countdownPrize = newCountdownPrize(data, ctime)
	} else {
		// 倒计时下一天
		m.countdownPrize = newCountdownPrize(
			m.countdownPrize.data.FirstData(),
			timeutil.DailyTime.NextTime(ctime),
		)
	}

	return m.countdownPrize
}

func newCountdownPrize(data *domestic_data.CountdownPrizeData, ctime time.Time) *countdown_prize {
	prize := data.Plunder.Try()
	collectTime := ctime.Add(data.WaitDuration)

	return &countdown_prize{
		data:        data,
		prize:       prize,
		collectTime: collectTime,
		desc:        data.RandomDesc(),
	}
}

type countdown_prize struct {
	data *domestic_data.CountdownPrizeData

	prize *resdata.Prize

	collectTime time.Time

	desc *domestic_data.CountdownPrizeDescData
}

func (cp *countdown_prize) NextCollectTime(ctime time.Time) time.Time {
	cp.collectTime = ctime.Add(cp.data.WaitDuration)
	return cp.collectTime
}

func (cp *countdown_prize) unmarshal(proto *server_proto.HeroCountdownPrizeServerProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	data := datas.GetCountdownPrizeData(proto.GetId())
	if data == nil {
		return
	}

	desc := datas.GetCountdownPrizeDescData(proto.GetDescId())
	if desc == nil {
		return
	}

	prize := resdata.UnmarshalPrize(proto.Prize, datas)
	if prize == nil || !prize.IsNotEmpty {
		return
	}

	cp.data = data
	cp.prize = prize
	cp.collectTime = timeutil.Unix64(proto.GetCollectTime())
	cp.desc = desc
}

func (cp *countdown_prize) encode() *server_proto.HeroCountdownPrizeServerProto {
	return &server_proto.HeroCountdownPrizeServerProto{
		Id:          cp.data.Id,
		Prize:       cp.prize.Encode(),
		CollectTime: timeutil.Marshal64(cp.collectTime),
		DescId:      cp.desc.Id,
	}
}

func (cp *countdown_prize) Prize() *resdata.Prize {
	return cp.prize
}

func (cp *countdown_prize) CollectTime() time.Time {
	return cp.collectTime
}

func (cp *countdown_prize) Desc() *domestic_data.CountdownPrizeDescData {
	return cp.desc
}

func (m *HeroDomestic) IncDailyHelpMemberTimes() uint64 {
	return m.AddDailyHelpMemberTimes(1)
}

func (m *HeroDomestic) AddDailyHelpMemberTimes(toAdd uint64) uint64 {
	m.dailyHelpMemberTimes += toAdd
	return m.dailyHelpMemberTimes
}

func (m *HeroDomestic) GetDailyHelpMemberTimes() uint64 {
	return m.dailyHelpMemberTimes
}

func (d *HeroDomestic) GetBuildingWorkerCdr() float64 {
	return d.buildingEffect.buildingWorkerCdr
}

func (d *HeroDomestic) GetBuildingWorkerCoef(extraWorkerCdr float64) float64 {
	return getWorkerCoef(d.buildingEffect.buildingWorkerCdr + extraWorkerCdr)
}

func (d *HeroDomestic) GetTechWorkerCdr() float64 {
	return d.buildingEffect.techWorkerCdr
}

func (d *HeroDomestic) GetTechWorkerCoef(extraWorkerCdr float64) float64 {
	return getWorkerCoef(d.buildingEffect.techWorkerCdr + extraWorkerCdr)
}

func getWorkerCoef(cdr float64) float64 {
	return 1 / (1 + math.Max(0, cdr))
}

func (h *HeroDomestic) StorageProtected(storage *ResourceStorage) (goldProtected, foodProtected, woodProtected, stoneProtected bool) {
	goldProtected = storage.Gold() <= h.GoldCapcity()
	foodProtected = storage.Food() <= h.FoodCapcity()
	woodProtected = storage.Wood() <= h.WoodCapcity()
	stoneProtected = storage.Stone() <= h.StoneCapcity()
	return
}

func (d *HeroDomestic) GoldCapcity() uint64 {
	return d.buildingEffect.GoldCapcity()
}

func (d *HeroDomestic) FoodCapcity() uint64 {
	return d.buildingEffect.FoodCapcity()
}

func (d *HeroDomestic) WoodCapcity() uint64 {
	return d.buildingEffect.WoodCapcity()
}

func (d *HeroDomestic) StoneCapcity() uint64 {
	return d.buildingEffect.StoneCapcity()
}

//func (d *HeroDomestic) GetCapcity(resType shared_proto.ResType) uint64 {
//	return d.buildingEffect.GetCapcity(resType)
//}

func (d *HeroDomestic) GetFarmExtraOutput(resType shared_proto.ResType) *data.Amount {
	return d.buildingEffect.GetFarmExtraOutput(resType)
}

func (d *HeroDomestic) GetExtraOutput(resType shared_proto.ResType) *data.Amount {
	return d.buildingEffect.GetExtraOutput(resType)
}

func (d *HeroDomestic) GetExtraOutputCapcity(resType shared_proto.ResType) *data.Amount {
	return d.buildingEffect.GetExtraOutputCapcity(resType)
}

func (d *HeroDomestic) ProtectedCapcity() uint64 {
	return d.buildingEffect.protectedCapcity
}

func (d *HeroDomestic) ProsperityCapcity() uint64 {
	return d.buildingEffect.prosperityCapcity
}

func (d *HeroDomestic) OuterCities() *OuterCities {
	return d.outerCities
}

func (d *HeroDomestic) GetTechnology(group uint64) *domestic_data.TechnologyData {
	return d.technologys[group]
}

func (d *HeroDomestic) SetTechnology(group uint64, toSet *domestic_data.TechnologyData) *domestic_data.TechnologyData {
	old := d.technologys[group]
	d.technologys[group] = toSet
	return old
}

func (d *HeroDomestic) GetBuilding(buildingType shared_proto.BuildingType) *domestic_data.BuildingData {
	return d.buildings[buildingType]
}

func (d *HeroDomestic) SetBuilding(building *domestic_data.BuildingData) {
	d.buildings[building.Type] = building
}

func (d *HeroDomestic) ResourcePointCount() int {
	return len(d.resourcePoint)
}

func (d *HeroDomestic) WalkResourcePoint(f func(pos uint64, data *ResourcePoint) (endWalk bool)) {
	for k, v := range d.resourcePoint {
		if f(k, v) {
			break
		}
	}
}

func (d *HeroDomestic) WalkTechnology(f func(pos uint64, data *domestic_data.TechnologyData) bool) {
	for k, v := range d.technologys {
		if f(k, v) {
			break
		}
	}
}

func (d *HeroDomestic) GetBuildingWorkerFatigueDuration() time.Duration {
	b := d.GetBuilding(shared_proto.BuildingType_GUAN_FU)
	if b == nil || b.Effect == nil || b.Effect.BuildingWorkerFatigueDuration < 0 {
		return time.Duration(0)
	}

	return b.Effect.BuildingWorkerFatigueDuration
}

func (d *HeroDomestic) GetFreeWorker(ctime time.Time) int {

	fatigueTime := ctime.Add(d.GetBuildingWorkerFatigueDuration())
	for i, v := range d.workerRestEndTime {
		if !d.WorkerIsUnlock(i, ctime) {
			continue
		}
		if v.Before(fatigueTime) {
			return i
		}
	}

	return -1
}

func (d *HeroDomestic) GetWorkerRestEndTime(pos int) (bool, time.Time) {

	if 0 <= pos && pos < len(d.workerRestEndTime) {
		return true, d.workerRestEndTime[pos]
	}

	return false, time.Time{}
}

func (d *HeroDomestic) GetWorkerSeekHelp(pos int) bool {
	if pos >= 0 && pos < len(d.workerSeekHelp) {
		return d.workerSeekHelp[pos]
	}
	return false
}

func (d *HeroDomestic) UnsetWorkerSeekHelpIfTrue(pos int) bool {
	if pos >= 0 && pos < len(d.workerSeekHelp) {
		if d.workerSeekHelp[pos] {
			d.workerSeekHelp[pos] = false
			d.checkRemoveWorkerSeekHelp[pos] = true
			return true
		}
	}
	return false
}

func (d *HeroDomestic) WorkerPosLegal(pos int) bool {
	return pos >= 0 && pos < len(d.workerLockStartTime)
}

func (d *HeroDomestic) WorkerNeverUnlocked(pos int) bool {
	if pos < 0 || pos >= len(d.workerLockStartTime) {
		return true
	}

	if d.workerAlwaysUnlocked[pos] {
		return false
	}

	return timeutil.IsZero(d.workerLockStartTime[pos])
}

func (d *HeroDomestic) WorkerIsUnlock(pos int, ctime time.Time) bool {
	if pos < 0 || pos >= len(d.workerAlwaysUnlocked) {
		return false
	}
	if d.workerAlwaysUnlocked[pos] {
		return true
	}

	return ctime.Before(d.workerLockStartTime[pos])
}

func (d *HeroDomestic) UnlockWorkerForever(pos int) bool {
	if pos < 0 || pos >= len(d.workerAlwaysUnlocked) {
		return false
	}
	d.workerAlwaysUnlocked[pos] = true
	return true
}

func (d *HeroDomestic) UnlockWorker(pos int, t time.Time) time.Time {
	if pos < 0 || pos >= len(d.workerLockStartTime) {
		return time.Time{}
	}
	d.workerLockStartTime[pos] = t
	return t
}

func (d *HeroDomestic) UnsetTechSeekHelpIfTrue(pos int) bool {
	if pos >= 0 && pos < len(d.techSeekHelp) {
		if d.techSeekHelp[pos] {
			d.techSeekHelp[pos] = false
			d.checkRemoveTechSeekHelp[pos] = true
			return true
		}
	}
	return false
}

func (d *HeroDomestic) TryRemoveSeekHelp(ctime time.Time) (workerPos, techPos []int) {

	for i, b := range d.checkRemoveWorkerSeekHelp {
		if b && ctime.After(d.workerRestEndTime[i]) {
			d.checkRemoveWorkerSeekHelp[i] = false
			workerPos = append(workerPos, i)
		}
	}

	for i, b := range d.checkRemoveTechSeekHelp {
		if b && ctime.After(d.techWorkerRestEndTime[i]) {
			d.checkRemoveTechSeekHelp[i] = false
			techPos = append(techPos, i)
		}
	}

	return
}

func (d *HeroDomestic) AddWorkerRestEndTime(pos int, ctime time.Time, toAdd time.Duration) (time.Time, bool) {

	if 0 <= pos && pos < len(d.workerRestEndTime) {
		ret := d.workerRestEndTime[pos]
		if ret.Before(ctime) {
			ret = ctime
		}
		newRestEndTime := ret.Add(toAdd)
		d.workerRestEndTime[pos] = newRestEndTime

		// 如果增加时间，则看下是不是刷新联盟求助按钮
		seekHelp := d.workerSeekHelp[pos]
		if toAdd > 0 {
			fatigueTime := ctime.Add(d.GetBuildingWorkerFatigueDuration())

			if newRestEndTime.After(fatigueTime) && ret.Before(fatigueTime) {
				// 刷新联盟求助按钮
				seekHelp = true
			}

			d.workerSeekHelp[pos] = seekHelp
		}

		return newRestEndTime, seekHelp
	}

	return ctime, false
}

func (d *HeroDomestic) GetTechWorkerFatigueDuration() time.Duration {
	b := d.GetBuilding(shared_proto.BuildingType_SHU_YUAN)
	if b == nil || b.Effect == nil || b.Effect.TechWorkerFatigueDuration < 0 {
		return time.Duration(0)
	}

	return b.Effect.TechWorkerFatigueDuration
}

func (d *HeroDomestic) GetFreeTechWorker(ctime time.Time) int {

	fatigueTime := ctime.Add(d.GetTechWorkerFatigueDuration())
	for i, v := range d.techWorkerRestEndTime {
		if v.Before(fatigueTime) {
			return i
		}
	}

	return -1
}

func (d *HeroDomestic) GetTechWorkerRestEndTime(pos int) (bool, time.Time) {

	if 0 <= pos && pos < len(d.techWorkerRestEndTime) {
		return true, d.techWorkerRestEndTime[pos]
	}

	return false, time.Time{}
}

func (d *HeroDomestic) AddTechWorkerRestEndTime(pos int, ctime time.Time, toAdd time.Duration) (time.Time, bool) {

	if 0 <= pos && pos < len(d.techWorkerRestEndTime) {
		ret := d.techWorkerRestEndTime[pos]
		if ret.Before(ctime) {
			ret = ctime
		}
		newRestEndTime := ret.Add(toAdd)
		d.techWorkerRestEndTime[pos] = newRestEndTime

		// 如果增加时间，则看下是不是刷新联盟求助按钮
		seekHelp := d.techSeekHelp[pos]
		if toAdd > 0 {
			fatigueTime := ctime.Add(d.GetTechWorkerFatigueDuration())

			if newRestEndTime.After(fatigueTime) && ret.Before(fatigueTime) {
				// 刷新联盟求助按钮
				seekHelp = true
			}

			d.techSeekHelp[pos] = seekHelp
		}

		return newRestEndTime, seekHelp
	}

	return ctime, false
}

func (d *HeroDomestic) ResourceCollectTimes() uint64 {
	return d.dailyResourceCollectTimes
}

func (d *HeroDomestic) IncreseResourceCollectTimes() uint64 {
	d.dailyResourceCollectTimes++
	return d.dailyResourceCollectTimes
}

func (d *HeroDomestic) GetResourceNextCollectTime(resType shared_proto.ResType) time.Time {
	return d.resourceNextCollectTimeMap[resType]
}

func (d *HeroDomestic) SetResourceNextCollectTime(resType shared_proto.ResType, toSet time.Time) {
	d.resourceNextCollectTimeMap[resType] = toSet
}

func (d *HeroDomestic) GetForgingTimes() *nexttimes {
	return d.forgingTimes
}

func (d *HeroDomestic) AddNewForgingPos(toAdds []uint64) {
	if len(toAdds) > 0 {
		for _, toAdd := range toAdds {
			d.newForgingPos = u64.AddIfAbsent(d.newForgingPos, toAdd)
		}
	}
}

func (d *HeroDomestic) RemoveNewForgingPos(toRemove uint64) {
	if len(d.newForgingPos) > 0 {
		d.newForgingPos = u64.RemoveIfPresent(d.newForgingPos, toRemove)
	}
}

func (hero *Hero) EncodeOtherDomestic(d *HeroDomestic) *shared_proto.HeroDomesticOtherProto {
	proto := &shared_proto.HeroDomesticOtherProto{}

	proto.Prosperity = u64.Int32(hero.home.prosperity)
	proto.LostProsperity = u64.Int32(hero.home.lostProsperity)
	proto.MaxProsperity = u64.Int32(d.ProsperityCapcity())

	proto.Sign = d.sign
	proto.Voice = d.voice

	return proto
}

func (d *HeroDomestic) CityEvent() *CityEvent {
	return d.cityEvent
}

func (d *HeroDomestic) ResetSeason() {
	d.isCollectSeasonPrize = false
}

func (d *HeroDomestic) ResetDaily() {
	d.cityEvent.resetDaily()
	d.dailyHelpMemberTimes = 0
	d.isResetCountdownPrize = true
	d.dailyResourceCollectTimes = 0
	d.forgingTimes.SetTimes(0)
	d.workshopRefreshTimes = 0
	d.collectProsperityDownCountDownPrizeTimes = 0
}

func (d *HeroDomestic) encodeServer() *server_proto.HeroDomesticServerProto {

	out := &server_proto.HeroDomesticServerProto{}

	out.WorkerRestEndTime = timeutil.MarshalArray32(d.workerRestEndTime)
	out.WorkerSeekHelp = d.workerSeekHelp
	out.CheckRemoveWorkerSeekHelp = d.checkRemoveWorkerSeekHelp

	out.WorkerAlwaysUnlocked = d.workerAlwaysUnlocked
	out.WorkerLockStartTime = timeutil.MarshalArray32(d.workerLockStartTime)

	// building
	out.ResourcePoint = make([]*server_proto.HeroResourcePointServerProto, 0, len(d.resourcePoint))
	for _, b := range d.resourcePoint {
		out.ResourcePoint = append(out.ResourcePoint, b.encodeServer())
	}
	out.ConflictResourcePoints = d.conflictResourcePoints

	for _, b := range d.buildings {
		out.Building = append(out.Building, b.Id)
	}

	// technology
	for _, tech := range d.technologys {
		out.Technology = append(out.Technology, tech.Id)
	}

	// guildTechnologys
	for _, tech := range d.guildTechnologys {
		out.GuildTechnology = append(out.GuildTechnology, tech.Id)
	}

	out.TechnologyRestEndTime = timeutil.MarshalArray32(d.techWorkerRestEndTime)
	out.TechnologySeekHelp = d.techSeekHelp
	out.CheckRemoveTechnologySeekHelp = d.checkRemoveTechSeekHelp

	out.DailyHelpMemberTimes = d.dailyHelpMemberTimes
	out.ForgingCombineTimes = d.forgingTimes.encodeServer()
	out.NewForgingPos = d.newForgingPos

	out.WorkshopIndex = d.workshopIndex
	out.WorkshopCollectTime = timeutil.Marshal64(d.workshopCollectTime)
	out.NextRefreshWorkshopTime = timeutil.Marshal64(d.nextRefreshWorkshopTime)
	out.WorkshopEquipmentIds = goods.GetEquipmentDataKeyArray(d.workshopEquipment)
	out.WorkshopRefreshTimes = d.workshopRefreshTimes

	out.DailyResourceCollectTimes = d.dailyResourceCollectTimes
	for k, v := range d.resourceNextCollectTimeMap {
		out.NextCollectTimeType = append(out.NextCollectTimeType, k)
		out.NextCollectTime = append(out.NextCollectTime, timeutil.Marshal64(v))
	}

	out.Sign = d.sign

	out.Voice = d.voice

	out.CityEvent = d.cityEvent.encode()

	out.Cities = d.outerCities.Encode()

	out.CountdownPrize = d.countdownPrize.encode()

	out.IsCollectSeasonPrize = d.isCollectSeasonPrize
	out.CollectProsperityDownCountDownPrizeTimes = d.collectProsperityDownCountDownPrizeTimes
	out.NextCanCollectProsperityDownCountDownPrizeTime = timeutil.Marshal64(d.nextCanCollectProsperityDownCountDownPrizeTime)

	return out
}

func (d *HeroDomestic) MarshalWorkerRestEndTime() []int32 {
	return timeutil.MarshalArray32(d.workerRestEndTime)
}

func (d *HeroDomestic) MarshalTechWorkerRestEndTime() []int32 {
	return timeutil.MarshalArray32(d.techWorkerRestEndTime)
}

func (d *HeroDomestic) GetBuildingIds() []int32 {

	ids := make([]int32, 0, len(d.buildings))
	for _, b := range d.buildings {
		ids = append(ids, u64.Int32(b.Id))
	}

	return ids
}

func (d *HeroDomestic) GetTechnologyIds() []int32 {
	ids := make([]int32, 0, len(d.technologys))
	for _, b := range d.technologys {
		ids = append(ids, u64.Int32(b.Id))
	}

	return ids
}

func (d *HeroDomestic) WalkMainAndOuterCityBuildings(walkFunc func(buildingData *domestic_data.BuildingData)) {
	for _, building := range d.buildings {
		walkFunc(building)
	}

	d.outerCities.WalkBuildings(func(layout *domestic_data.OuterCityLayoutData, building *domestic_data.OuterCityBuildingData) {
		walkFunc(building.BuildingData)
	})
}

func (d *HeroDomestic) WalkMainAndOuterCityBuildingsAndTechnologysHasEffect(walkFunc func(effect *sub.BuildingEffectData)) {
	// 先计算初始加成
	if d.buildingInitEffect != nil {
		walkFunc(d.buildingInitEffect)
	}

	d.WalkMainAndOuterCityBuildings(func(buildingData *domestic_data.BuildingData) {
		if buildingData.Effect != nil {
			walkFunc(buildingData.Effect)
		}
	})

	for _, tech := range d.technologys {
		if tech.Effect != nil {
			walkFunc(tech.Effect)
		}
	}

	for _, tech := range d.guildTechnologys {
		if tech.Effect != nil {
			walkFunc(tech.Effect)
		}
	}
}

// 这里开始资源点

func (hero *HeroDomestic) GetLayoutRes(layoutId uint64) *ResourcePoint {
	return hero.resourcePoint[layoutId]
}

func (hero *HeroDomestic) SetResourcePoint(layoutData *domestic_data.BuildingLayoutData, building *domestic_data.BuildingData, outputStartTime time.Time) *ResourcePoint {
	rp := NewResourcePoint(layoutData, building, outputStartTime)
	hero.resourcePoint[layoutData.Id] = rp
	return rp
}

func (hero *HeroDomestic) LayoutResMap() map[uint64]*ResourcePoint {
	return hero.resourcePoint
}

func (hero *HeroDomestic) SetSign(toSet string) {
	hero.sign = toSet
}

func (hero *HeroDomestic) SetVoice(toSet []byte) {
	hero.voice = toSet
}

func (hero *Hero) encodeDomestic(d *HeroDomestic, configDatas config.Configs, ctime time.Time) *shared_proto.HeroDomesticProto {

	out := &shared_proto.HeroDomesticProto{}

	out.Gold = u64.Int32(hero.unsafeResource.gold)
	out.Food = u64.Int32(hero.unsafeResource.food)
	out.Wood = u64.Int32(hero.unsafeResource.wood)
	out.Stone = u64.Int32(hero.unsafeResource.stone)

	out.SafeGold = u64.Int32(hero.safeResource.gold)
	out.SafeFood = u64.Int32(hero.safeResource.food)
	out.SafeWood = u64.Int32(hero.safeResource.wood)
	out.SafeStone = u64.Int32(hero.safeResource.stone)

	out.Prosperity = u64.Int32(hero.home.prosperity)
	out.LostProsperity = u64.Int32(hero.home.lostProsperity)

	// 黄金分割线 -----------

	out.GoldCapcity = u64.Int32(d.GoldCapcity())
	out.FoodCapcity = u64.Int32(d.FoodCapcity())
	out.WoodCapcity = u64.Int32(d.WoodCapcity())
	out.StoneCapcity = u64.Int32(d.StoneCapcity())

	out.ProtectedCapcity = u64.Int32(d.ProtectedCapcity())

	out.WorkerRestEndTime = d.MarshalWorkerRestEndTime()
	out.WorkerSeekHelp = d.workerSeekHelp

	out.WorkerAlwaysUnlocked = d.workerAlwaysUnlocked
	out.WorkerLockStartTime = timeutil.MarshalArray32(d.workerLockStartTime)

	for _, b := range d.resourcePoint {
		out.ResourcePoint = append(out.ResourcePoint, b.encodeClient(hero, ctime))
	}

	out.Building = d.GetBuildingIds()
	out.Technology = d.GetTechnologyIds()

	out.TechnologyRestEndTime = d.MarshalTechWorkerRestEndTime()
	out.TechnologySeekHelp = d.techSeekHelp
	out.DailyHelpMemberTimes = u64.Int32(d.dailyHelpMemberTimes)

	out.MaxProsperity = u64.Int32(d.ProsperityCapcity())

	out.BuildingWorkerCoef = i32.MultiF64(1000, d.GetBuildingWorkerCdr())
	out.TechWorkerCoef = i32.MultiF64(1000, d.GetTechWorkerCdr())
	out.BuildingWorkerFatigueDuration = timeutil.DurationMarshal32(d.GetBuildingWorkerFatigueDuration())
	out.TechWorkerFatigueDuration = timeutil.DurationMarshal32(d.GetTechWorkerFatigueDuration())

	out.ForgingTimes, out.ForgingNextTime = d.forgingTimes.encodeClient()
	out.NewForgingPos = u64.Int32Array(d.newForgingPos)

	out.DailyResourceCollectTimes = u64.Int32(d.dailyResourceCollectTimes)
	for k, v := range d.resourceNextCollectTimeMap {
		out.NextCollectTimeType = append(out.NextCollectTimeType, k)
		out.NextCollectTime = append(out.NextCollectTime, timeutil.Marshal32(v))
	}

	out.Sign = d.sign

	out.Voice = d.voice

	out.CityEvent = d.cityEvent.encode()

	out.ResourcePointV2 = hero.EncodeResourcePointV2(configDatas)

	out.Jade = u64.Int32(hero.jade)
	out.JadeOre = u64.Int32(hero.jadeOre)
	out.HistoryJade = u64.Int32(hero.historyJade)
	out.TodayObtainJade = u64.Int32(hero.todayObtainJade)

	out.Cities = d.outerCities.Encode()

	out.CountdownCollectTime = timeutil.Marshal32(d.countdownPrize.collectTime)

	out.IsCollectSeasonPrize = d.isCollectSeasonPrize

	if d.workshopIndex > 0 {
		out.WorkshopCollectTime = timeutil.Marshal32(d.workshopCollectTime)
	}

	out.BuildingCostReduceCoef = i32.MultiF64(1000, hero.BuildingEffect().buildingCostReduceCoef)
	out.TechCostReduceCoef = i32.MultiF64(1000, hero.BuildingEffect().techCostReduceCoef)

	return out
}

var resTypes = []shared_proto.ResType{shared_proto.ResType_GOLD, shared_proto.ResType_STONE}

// 把资源点序列化给客户端
func (hero *Hero) EncodeResourcePointV2(configDatas config.Configs) *shared_proto.ResourcePointV2Proto {
	cubes := configDatas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(hero.BaseLevel())

	resCount := int32(0)
	conflictResCount := int32(0)

	for _, cube := range cubes {
		layoutData := configDatas.RegionConfig().GetLayoutDataByEvenOffset(cube)
		if layoutData == nil {
			continue
		}

		resCount++
		if hero.IsConflictResourcePoint(layoutData) {
			conflictResCount++
		}
	}

	proto := &shared_proto.ResourcePointV2Proto{ResCount: resCount, ConflictResCount: conflictResCount}

	for _, resType := range resTypes {
		srp := &shared_proto.SingleResourcePointV2Proto{ResType: resType}
		proto.ResourcePoint = append(proto.ResourcePoint, srp)

		outputAmount := configDatas.BuildingLayoutMiscData().SingleResOutPutAmount(resType)
		if outputAmount == nil {
			continue
		}

		totalAmount := u64.Int32(data.TotalAmount(outputAmount))
		if extraOutput := hero.Domestic().GetExtraOutput(resType); extraOutput != nil {
			totalAmount = u64.Int32(data.TotalAmount(outputAmount, extraOutput))
		}

		srp.OriginOutput = resCount * totalAmount
		srp.RealOutput = (resCount - conflictResCount) * totalAmount
	}

	return proto
}

// 城内事件
type CityEvent struct {
	acceptTimes     uint64                       // 第几次的城内事件
	canExchangeTime time.Time                    // 能够兑换/放弃城内事件的时间
	eventData       *domestic_data.CityEventData // 当前事件
}

func (d *CityEvent) AcceptTimes() uint64 {
	return d.acceptTimes
}

func (d *CityEvent) EventData() *domestic_data.CityEventData {
	return d.eventData
}

func (d *CityEvent) CanExchangeTime() time.Time {
	return d.canExchangeTime
}

func (d *CityEvent) SetCanExchangeTime(toSet time.Time) time.Time {
	d.canExchangeTime = toSet
	return d.canExchangeTime
}

func (d *CityEvent) Accept(data *domestic_data.CityEventData) {
	d.acceptTimes++
	d.eventData = data
}

func (d *CityEvent) ExchangeOrGiveUp() {
	d.eventData = nil
}

func (d *CityEvent) resetDaily() {
	d.acceptTimes = 0
}

func (d *CityEvent) encode() *shared_proto.CityEventProto {
	proto := &shared_proto.CityEventProto{}

	proto.AcceptTimes = u64.Int32(d.acceptTimes)
	proto.CanExchangeTime = timeutil.Marshal32(d.canExchangeTime)
	if d.eventData != nil {
		proto.EventId = u64.Int32(d.eventData.Id)
	}

	return proto
}

func (d *CityEvent) unmarshal(p *shared_proto.CityEventProto, datas *config.ConfigDatas) {
	if p == nil {
		return
	}

	d.acceptTimes = u64.FromInt32(p.GetAcceptTimes())
	d.canExchangeTime = timeutil.Unix32(p.GetCanExchangeTime())
	if p.GetEventId() != 0 {
		d.eventData = datas.GetCityEventData(u64.FromInt32(p.GetEventId()))
		if d.eventData == nil {
			logrus.WithField("cityEventId", p.EventId).Errorln("玩家的城内事件丢失")
		}
	}
}

var taxResType = []shared_proto.ResType{
	shared_proto.ResType_GOLD,
	shared_proto.ResType_STONE,
}

// 税收
func (hero *Hero) UpdateTax(ctime time.Time, duration time.Duration, getBuffEffectData func(uint64) *data.BuffEffectData) (update bool) {

	if duration <= 0 {
		return
	}

	// 看下时间是否到了下次收税时间
	//b := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU)
	//if b == nil || b.Effect == nil {
	//	return
	//}

	startTime := hero.MiscData().GetNextCollectTaxTime()
	diff := ctime.Sub(startTime)
	if diff < 0 {
		// 时间未到
		return
	}

	// 保护一下
	if timeutil.IsZero(startTime) {
		hero.MiscData().SetNextCollectTaxTime(ctime.Add(duration))
		return
	}

	hero.MiscData().SetNextCollectTaxTime(ctime.Add(duration - (diff % duration)))

	// 计算buff

	// 根据时间计算，要加多少份
	toAddDuration := ((diff / duration) + 1) * duration
	hours := toAddDuration.Hours()

	for _, resType := range taxResType {
		current := hero.GetUnsafeResource().GetRes(resType)
		capcity := hero.buildingEffect.GetCapcity(resType)

		canAdd := u64.Sub(capcity, current)
		if canAdd > 0 {
			// 可以加
			if toAddPerHour := hero.BuildingEffect().GetTax(resType); toAddPerHour > 0 {
				totalToAdd := u64.Min(u64.MultiF64(toAddPerHour, hours), canAdd)
				totalToAdd = hero.calcTaxDiffWithBuff(totalToAdd, getBuffEffectData)
				if totalToAdd > 0 {
					hero.GetUnsafeResource().AddRes(resType, totalToAdd)
					update = true
				}
			}
		}
	}

	return
}

func (hero *Hero) calcTaxDiffWithBuff(toAdd uint64, getBuffEffectData func(uint64) *data.BuffEffectData) (result uint64) {
	result = toAdd
	for _, b := range hero.Buff().Buffs(shared_proto.BuffEffectType_Buff_ET_tax) {
		result = b.EffectData.Tax.CalculateByPercent(result)
	}

	return
}
