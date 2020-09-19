package sharedguilddata

import (
	"bytes"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/config/xiongnu"
	"github.com/lightpaw/male7/entity/daily_amount"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/pb/guild"
	xiongnu2 "github.com/lightpaw/male7/gen/pb/xiongnu"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/collection"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"sort"
	"time"
)

func NewGuild(id int64, name, flagName string, datas *config.ConfigDatas, createTime time.Time) *Guild {
	g := &Guild{}
	g.id = id
	g.name = name
	g.flagName = flagName
	g.levelData = datas.GuildLevelData().MinKeyData
	g.createTime = createTime
	g.updateBuildingAmountTime = createTime
	g.country = datas.GuildConfig().DefaultGuildCountry
	g.prestigeHeroIdMap = make(map[int64]struct{})

	g.technologyMap = make(map[uint64]*guild_data.GuildTechnologyData)
	g.memberMap = make(map[int64]*GuildMember)

	g.memberLeaveTimeMap = make(map[int64]int64)

	g.donateRecords = make([]*shared_proto.GuildDonateRecordProto, 0, 2)

	g.donateHeroIdMap = make(map[int64]struct{})
	g.seekHelpMap = make(map[string]*shared_proto.GuildSeekHelpProto)

	g.bigBoxData = datas.GuildBigBoxData().MinKeyData
	g.fullBigBoxMemberIdMap = make(map[int64]struct{})

	g.unlockResistXiongNuData = datas.ResistXiongNuData().MinKeyData
	g.prestigeDaily = daily_amount.NewDailyAmount(datas.GuildConfig().KeepDailyPrestigeCount)
	g.prestigeCoreHourly = daily_amount.NewDailyAmount(datas.GuildConfig().KeepHourlyPrestigeCount)

	g.mcWarRecords = &shared_proto.McWarAllRecordProto{}
	g.mcWarRecordMsg = guild.NewS2cViewMcWarRecordMsg(nil, nil)

	g.marks = make([]*shared_proto.GuildMarkProto, datas.GuildGenConfig().GuildMarkCount)
	g.markMsgCache = make([]pbutil.Buffer, datas.GuildGenConfig().GuildMarkCount)

	g.SendYinliangToMe = make(map[int64]*shared_proto.GuildYinliangSendProto)
	g.yinliangSendToGuildRecords = collection.NewRingList(10)

	g.hostMingcIds = make(map[uint64]int64)
	g.weeklyTasks = make(map[server_proto.GuildTaskType]uint64)
	g.weeklyTaskStageIndexs = make(map[server_proto.GuildTaskType]int)
	g.weeklyTasksVersion = 1
	g.dailyMcBuildCounts = make(map[uint64]uint64)

	return g
}

type Guild struct {
	id int64

	name string

	flagName string

	country *country.CountryData // 联盟所属国家

	createTime time.Time

	lastPrestigeRank   uint64 // 前一次声望排名（用于每日联盟排行奖励）
	prestige           uint64 // 声望
	historyMaxPretige  uint64 // 历史最大声望
	prestigeHeroIdMap  map[int64]struct{}
	prestigeDaily      *daily_amount.DailyAmount // 声望每日数据
	prestigeCoreHourly *daily_amount.DailyAmount // 核心声望每小时数据,核心声望（管理层获得的声望统计，永久隐藏属性，前端不展示）

	hufu     uint64 // 虎符
	yinliang uint64 // 银两

	SendYinliangToMe map[int64]*shared_proto.GuildYinliangSendProto // 赠送银两给我的记录

	nextUpdatePrestigeTargetTime time.Time // 下次可以更新声望国家的时间

	nextChangeNameTime time.Time
	freeChangeName     bool

	template *guild_data.NpcGuildTemplate

	flagType uint64

	// 联盟科技，key是group
	technologyMap      map[uint64]*guild_data.GuildTechnologyData
	techUpgradeData    *guild_data.GuildTechnologyData
	techUpgradeEndTime time.Time
	techCdrTimes       uint64

	technologyArrayCache []*guild_data.GuildTechnologyData

	// 联盟等级数据
	levelData *guild_data.GuildLevelData

	// 联盟建设值
	buildingAmount           uint64
	updateBuildingAmountTime time.Time // 变更联盟建设值的时间
	donateHeroIdMap          map[int64]struct{}

	// 升级时间
	upgradeEndTime time.Time
	cdrTimes       uint64

	// leader

	text         string   // 对外公告
	internalText string   // 对内公告
	labels       []string // 联盟标签

	friendGuildText string // 友盟公告
	enemyGuildText  string // 友盟公告

	classNames      []string
	classTitleProto *shared_proto.GuildClassTitleProto

	leaderId          int64
	leaderOfflineTime time.Time

	// 禅让
	changeLeaderId   int64
	changeLeaderTime time.Time

	memberMap       map[int64]*GuildMember
	kickMemberCount uint64 // 每日踢人数，限制每日最多踢多少人

	// 入盟条件
	rejectAutoJoin        bool   // false表示达到条件直接入盟，true表示需要申请才能加入
	requiredHeroLevel     uint64 // 君主等级
	requiredJunXianLevel  uint64 // 百战军衔
	requiredTowerMaxFloor uint64 // 需要的最大千重楼层数

	// 弹劾盟主
	impeachLeader *impeach_leader

	// 邀请列表
	invateHeroIds     []int64
	invateExpiredTime []time.Time

	// 申请列表
	requestJoinHeroIds     []int64
	requestJoinExpiredTime []time.Time

	// 成员离开联盟时间（离开联盟后一段时间之内不能再加入联盟）
	memberLeaveTimeMap map[int64]int64

	// 捐献记录
	donateRecords []*shared_proto.GuildDonateRecordProto

	//// 大事记
	//bigEvents []*shared_proto.GuildBigEventProto
	//
	//// 动态
	//dynamics []*shared_proto.GuildDynamicProto

	// 完成了的联盟目标数量/当前联盟目标(可能为空)
	doingTargets []*targetWithEndTime

	statueRealmId  int64         // 联盟雕像放置的场景id，0表示没有放置联盟雕像
	statueCacheMsg pbutil.Buffer // 联盟雕像位置缓存的消息，如果没有放置联盟雕像，此处为空

	seekHelpMap map[string]*shared_proto.GuildSeekHelpProto

	// 当前的大宝箱
	bigBoxData   *guild_data.GuildBigBoxData
	bigBoxEnergy uint64

	// 满了的大宝箱（可领取）
	fullBigBoxData        *guild_data.GuildBigBoxData
	fullBigBoxMemberIdMap map[int64]struct{}

	// 联盟标记
	marks        []*shared_proto.GuildMarkProto
	markMsgCache []pbutil.Buffer

	changed bool

	dailyResetTime  time.Time
	weeklyResetTime time.Time

	unlockResistXiongNuData        *xiongnu.ResistXiongNuData            // 解锁了的抗击匈奴数据
	isStartResistXiongNuToday      bool                                  // 是否今天已经开启了
	resistXiongNuDefenders         []int64                               // 匈奴防守者
	lastResistXiongNuProto         *shared_proto.LastResistXiongNuProto  // 上次挑战匈奴数据
	lastResistXiongNuFightProto    *shared_proto.ResistXiongNuFightProto // 上次匈奴战斗排行榜数据
	lastResistXiongNuFightProtoMsg pbutil.Buffer

	mcWarRecords   *shared_proto.McWarAllRecordProto // 名城战记录
	mcWarRecordMsg pbutil.Buffer                     // 名城战记录消息缓存

	yinliangRecords   []*shared_proto.GuildYinliangRecordProto // 银两记录
	yinliangRecordMsg pbutil.Buffer

	//yinliangSendToGuildRecords   []*shared_proto.GuildYinliangSendToGuildProto // 最近赠送
	yinliangSendToGuildRecords *collection.RingList

	hostMingcIds map[uint64]int64 // 占领的名城

	// 联盟工坊
	workshop                 *Workshop // 生产次数
	workshopTodayCompleted   bool      // 今日竣工奖励
	workshopOutput           uint64    // 联盟工坊今日产出次数
	workshopOutputPrizeCount uint64    // 产出奖励个数
	workshopBeenHurtTimes    uint64    // 被破坏总次数

	recommendMcBuilds []uint64 // 推荐营建的名城

	dailyMcBuildCounts map[uint64]uint64 // 每日名城营建次数，和名城一起更新
	// 周任务
	weeklyTasks           map[server_proto.GuildTaskType]uint64
	weeklyTaskStageIndexs map[server_proto.GuildTaskType]int // 当前进度下标
	weeklyTasksVersion    int32                              // 任务版本号
	weeklyTasksMsg        pbutil.Buffer                      // 缓存消息，只生成一次

	// 联盟转国
	changeCountryWaitEndTime int64                // 转国等待结束时间
	changeCountryTarget      *country.CountryData // 转国目标（转到哪个国家）
	changeCountryNextTime    int64                // 下次可以转国时间
}

func (g *Guild) GetChangeCountryWaitEndTime() int64 {
	return g.changeCountryWaitEndTime
}

func (g *Guild) GetChangeCountryTarget() *country.CountryData {
	return g.changeCountryTarget
}

func (g *Guild) GetChangeCountryNextTime() int64 {
	return g.changeCountryNextTime
}

func (g *Guild) SetChangeCountry(country *country.CountryData, waitEndTime, nextTime int64) {
	g.changeCountryTarget = country
	g.changeCountryWaitEndTime = waitEndTime
	g.changeCountryNextTime = nextTime
}

func (g *Guild) CancelChangeCountry() {
	g.changeCountryTarget = nil
	g.changeCountryWaitEndTime = 0
}

type GuildOperateType uint64

const (
	// todo 以后统一到 OperateType 中
	JoinGuild       = GuildOperateType(1)
	ReplyJoinGuild  = GuildOperateType(2)
	InviteJoinGuild = GuildOperateType(5)
	LeaveGuild      = GuildOperateType(3)
	KickLeaveGuild  = GuildOperateType(4)
)

type GuildContext struct {
	OperType     GuildOperateType
	OperatorId   int64
	OperatorName string
}

func (g *Guild) RangeMarkMsg(f func(markMsg pbutil.Buffer) bool) {
	for i, markMsg := range g.markMsgCache {
		if markMsg == nil {
			mark := g.marks[i]
			if mark == nil {
				continue
			}
			markMsg = guild.NewS2cUpdateGuildMarkMsg(mark).Static()
			g.markMsgCache[i] = markMsg
		}

		if !f(markMsg) {
			break
		}
	}
}

func (g *Guild) AddMark(mark *shared_proto.GuildMarkProto, markMsg pbutil.Buffer) {
	idx := int(mark.Index - 1)
	if idx >= 0 && idx < len(g.marks) {
		g.marks[idx] = mark
		if markMsg != nil {
			g.markMsgCache[idx] = markMsg.Static()
		}
	}
}

func (g *Guild) RemoveMark(idx int) {
	if idx >= 0 && idx < len(g.marks) {
		g.marks[idx] = nil
		g.markMsgCache[idx] = nil
	}
}

func (g *Guild) RecommendMcBuilds() []uint64 {
	return g.recommendMcBuilds
}

func (g *Guild) RecommendMcBuildExist(mcId uint64) bool {
	for _, id := range g.recommendMcBuilds {
		if id == mcId {
			return true
		}
	}

	return false
}

func (g *Guild) AddRecommendMcBuild(newMcId uint64, maxLen int) (succ bool, result []uint64) {
	if g.RecommendMcBuildExist(newMcId) {
		return
	}
	if len(g.recommendMcBuilds) >= maxLen {
		g.recommendMcBuilds = g.recommendMcBuilds[1:]
	}
	g.recommendMcBuilds = append(g.recommendMcBuilds, newMcId)

	succ = true
	result = g.recommendMcBuilds
	return
}

func (g *Guild) GetTechnology(group uint64) *guild_data.GuildTechnologyData {
	return g.technologyMap[group]
}

func (g *Guild) SetTechnology(t *guild_data.GuildTechnologyData) {
	g.technologyMap[t.Group] = t
	g.technologyArrayCache = nil

	if t.BigBox != nil {
		g.setBigBoxData(t.BigBox)
	}
}

func (g *Guild) setBigBoxData(bigBox *guild_data.GuildBigBoxData) {
	g.bigBoxData = bigBox
	if g.bigBoxEnergy >= bigBox.UnlockEnergy {
		g.bigBoxEnergy = u64.Sub(bigBox.UnlockEnergy, 1)
	}
}

func (g *Guild) GetEffectTechnology() (technologys []*guild_data.GuildTechnologyData) {
	if g.technologyArrayCache == nil {
		g.technologyArrayCache = g.newEffectTechnology()
	}
	return g.technologyArrayCache
}

func (g *Guild) newEffectTechnology() (technologys []*guild_data.GuildTechnologyData) {
	for _, t := range g.technologyMap {
		if t.Effect != nil {
			technologys = append(technologys, t)
		}
	}
	return
}

func (g *Guild) GetTechUpgradeData() *guild_data.GuildTechnologyData {
	return g.techUpgradeData
}

func (g *Guild) GetTechUpgradeEndTime() time.Time {
	return g.techUpgradeEndTime
}

func (g *Guild) UpgradeTechnology(data *guild_data.GuildTechnologyData, endTime time.Time) {
	g.techUpgradeData = data
	g.techUpgradeEndTime = endTime
}

func (g *Guild) SetTechUpgradeEndTime(toSet time.Time) {
	g.techUpgradeEndTime = toSet
}

func (g *Guild) GetTechCdrTimes() uint64 {
	return g.techCdrTimes
}

func (g *Guild) IncTechCdrTimes() {
	g.techCdrTimes++
}

func (g *Guild) AddBigBoxEnergy(toAdd uint64) (newEnergy uint64, full bool) {
	g.bigBoxEnergy += toAdd
	return g.bigBoxEnergy, g.bigBoxEnergy >= g.bigBoxData.UnlockEnergy
}

func (g *Guild) ClearFullBigBox() (data *guild_data.GuildBigBoxData, memberIds []int64) {
	data = g.fullBigBoxData
	g.fullBigBoxData = nil

	for k := range g.fullBigBoxMemberIdMap {
		memberIds = append(memberIds, k)
		delete(g.fullBigBoxMemberIdMap, k)
	}

	return
}

func (g *Guild) GetBigBoxData() *guild_data.GuildBigBoxData {
	return g.bigBoxData
}

func (g *Guild) GetBigBoxEnergy() uint64 {
	return g.bigBoxEnergy
}

func (g *Guild) GetFullBigBoxData() *guild_data.GuildBigBoxData {
	return g.fullBigBoxData
}

func (g *Guild) IsFullBigBoxMember(memberId int64) bool {
	_, exist := g.fullBigBoxMemberIdMap[memberId]
	return exist
}

func (g *Guild) RemoveFullBigBoxMemberId(memberId int64) bool {
	if _, exist := g.fullBigBoxMemberIdMap[memberId]; exist {
		delete(g.fullBigBoxMemberIdMap, memberId)
		return true
	}
	return false
}

func (g *Guild) SetNextBigBox(toSet *guild_data.GuildBigBoxData, fullBigBoxMemberIds []int64) {

	g.fullBigBoxData = g.bigBoxData
	g.bigBoxData = toSet
	g.bigBoxEnergy = u64.Sub(g.bigBoxEnergy, g.fullBigBoxData.UnlockEnergy)
	if g.bigBoxEnergy >= toSet.UnlockEnergy {
		g.bigBoxEnergy = u64.Sub(toSet.UnlockEnergy, 1)
	}

	if len(g.fullBigBoxMemberIdMap) > 0 {
		for k := range g.fullBigBoxMemberIdMap {
			delete(g.fullBigBoxMemberIdMap, k)
		}
	}

	// 当前在联盟中的人都可以领取
	for _, memberId := range fullBigBoxMemberIds {
		g.fullBigBoxMemberIdMap[memberId] = struct{}{}
	}

}

type guild_seek_help struct {
	heroId int64

	proto *shared_proto.GuildSeekHelpProto
}

func (sh *guild_seek_help) TrySeekHelp(heroIdBytes []byte) bool {

	if len(sh.proto.HelpHeroIds) >= int(sh.proto.HelpMaxHeroCount) {
		return false
	}

	for _, v := range sh.proto.HelpHeroIds {
		if bytes.Equal(heroIdBytes, v) {
			return false
		}
	}

	sh.proto.HelpHeroIds = append(sh.proto.HelpHeroIds, heroIdBytes)

	return true
}

func NewHeroWorkerSeekHelpKey(heroId []byte, helpType, workerPos int32) string {

	n := len(heroId)
	data := make([]byte, n+2)
	copy(data, heroId)

	data[n] = byte(helpType)
	data[n+1] = byte(workerPos)

	return util.Byte2String(data)
}

func (g *Guild) AddSeekHelp(proto *shared_proto.GuildSeekHelpProto) {
	proto.Id = NewHeroWorkerSeekHelpKey(proto.HeroId, proto.HelpType, proto.WorkerPos)
	g.seekHelpMap[proto.Id] = proto
}

func (g *Guild) RangeSeekHelp(f func(proto *shared_proto.GuildSeekHelpProto) (isContinue bool)) {
	for _, v := range g.seekHelpMap {
		if !f(v) {
			return
		}
	}
}

func (g *Guild) GetSeekHelp(key string) *shared_proto.GuildSeekHelpProto {
	return g.seekHelpMap[key]
}

func (g *Guild) RemoveSeekHelp(key string) {
	delete(g.seekHelpMap, key)
}

func (g *Guild) SetChanged() {
	g.changed = true
}

func (g *Guild) SetFalseIfChanged() bool {
	if g.changed {
		g.changed = false
		return true
	}

	return false
}

func (g *Guild) GetNpcTemplate() *guild_data.NpcGuildTemplate {
	return g.template
}

func (g *Guild) SetNpcTemplate(toSet *guild_data.NpcGuildTemplate) {
	g.template = toSet
	if toSet != nil {
		g.text = toSet.Text
		g.internalText = toSet.InternalText
		g.labels = toSet.Labels
	}

}

func (g *Guild) NewBasicProto() *shared_proto.GuildBasicProto {
	proto := &shared_proto.GuildBasicProto{}
	proto.Id = i64.Int32(g.Id())
	proto.Name = g.Name()
	proto.FlagName = g.FlagName()
	proto.Level = u64.Int32(g.levelData.Level)
	proto.Country = u64.Int32(g.country.Id)

	return proto
}

func (g *Guild) NewHeroGuildProto() *shared_proto.HeroGuildProto {
	proto := &shared_proto.HeroGuildProto{}
	proto.Id = i64.Int32(g.Id())
	proto.Name = g.Name()
	proto.FlagName = g.FlagName()
	proto.Level = u64.Int32(g.levelData.Level)
	proto.Country = u64.Int32(g.country.Id)
	proto.Leader = idbytes.ToBytes(g.leaderId)

	return proto
}

func (g *Guild) NewSnapshot() *guildsnapshotdata.GuildSnapshot {
	s := &guildsnapshotdata.GuildSnapshot{}
	s.Id = g.id
	s.Name = g.name
	s.FlagName = g.flagName
	s.FlagType = g.flagType
	s.GuildLevel = g.LevelData()
	s.Country = g.country
	s.Prestige = g.prestige
	s.LeaderId = g.LeaderId()

	if npcid.IsNpcId(s.LeaderId) {
		s.IsNpcGuild = true
		member := g.GetMember(s.LeaderId)
		if member != nil {
			s.LeaderSnapshotIfIsNpc = member.npcProto
		}
	}

	s.MemberCount = uint64(len(g.memberMap))
	s.UserMemberIds = g.AllUserMemberIds()
	s.ResistXiongNuDefenders = g.ResistXiongNuDefenders()
	s.TotalPrestigeDaily = g.prestigeDaily.Total()

	// 入盟条件
	s.RejectAutoJoin = g.rejectAutoJoin
	s.RequiredHeroLevel = g.requiredHeroLevel
	s.RequiredJunXianLevel = g.requiredJunXianLevel
	s.RequiredTowerMaxFloor = g.requiredTowerMaxFloor

	s.Text = g.text

	s.Technologys = g.GetEffectTechnology()

	return s
}

func (g *Guild) GetInvateHeroIds() []int64 {
	return g.invateHeroIds
}

func (g *Guild) AddInvateHero(heroId int64, expiredTime time.Time) {
	if i64.GetIndex(g.invateHeroIds, heroId) < 0 {
		g.invateHeroIds = append(g.invateHeroIds, heroId)
		g.invateExpiredTime = append(g.invateExpiredTime, expiredTime)
	}
}

func (g *Guild) RemoveFirstInvateHeroId() int64 {
	if len(g.invateHeroIds) > 0 {
		first := g.invateHeroIds[0]
		g.invateHeroIds = i64.LeftShift(g.invateHeroIds, 0, 1)
		g.invateExpiredTime = timeutil.LeftShift(g.invateExpiredTime, 0, 1)
		return first
	}
	return 0
}

func (g *Guild) RemoveInvateHero(removeHeroId int64) bool {
	var idx int
	g.invateHeroIds, idx = i64.LeftShiftRemoveIfPresentReturnIndex(g.invateHeroIds, removeHeroId)
	if idx >= 0 {
		g.invateExpiredTime = timeutil.LeftShift(g.invateExpiredTime, idx, 1)
		return true
	}

	return false
}

func (g *Guild) RemoveExpiredInvateHero(ctime time.Time) (removeHeroIds []int64) {

	n := len(g.invateExpiredTime)
	expiredCount := 0
	for i := 0; i < n; i++ {
		if ctime.Before(g.invateExpiredTime[i]) {
			break
		}

		// 过期
		expiredCount++
	}
	if expiredCount <= 0 {
		return
	}

	// 需要移除这么多个
	g.invateExpiredTime = timeutil.LeftShift(g.invateExpiredTime, 0, expiredCount)
	g.invateHeroIds, removeHeroIds = i64.LeftShiftReturnRemovedValues(g.invateHeroIds, 0, expiredCount)
	return
}

func (g *Guild) GetRequestJoinHeroIds() []int64 {
	return g.requestJoinHeroIds
}

func (g *Guild) AddRequestJoinHeroId(heroId int64, expiredTime time.Time) {
	if i64.GetIndex(g.requestJoinHeroIds, heroId) < 0 {
		g.requestJoinHeroIds = append(g.requestJoinHeroIds, heroId)
		g.requestJoinExpiredTime = append(g.requestJoinExpiredTime, expiredTime)
	}
}

func (g *Guild) RemoveRequestJoinHeroId(removeHeroId int64) bool {
	var idx int
	g.requestJoinHeroIds, idx = i64.LeftShiftRemoveIfPresentReturnIndex(g.requestJoinHeroIds, removeHeroId)
	if idx >= 0 {
		g.requestJoinExpiredTime = timeutil.LeftShift(g.requestJoinExpiredTime, idx, 1)
		return true
	}

	return false
}

func (g *Guild) RemoveFirstRequestJoinHero() (removeHeroId int64) {

	if len(g.requestJoinHeroIds) > 0 {
		removeHeroId = g.requestJoinHeroIds[0]
	}

	g.requestJoinExpiredTime = timeutil.LeftShift(g.requestJoinExpiredTime, 0, 1)
	g.requestJoinHeroIds = i64.LeftShift(g.requestJoinHeroIds, 0, 1)
	return
}

func (g *Guild) RemoveExpiredRequestJoinHero(ctime time.Time) (removeHeroIds []int64) {

	n := len(g.requestJoinExpiredTime)
	expiredCount := 0
	for i := 0; i < n; i++ {
		if ctime.Before(g.requestJoinExpiredTime[i]) {
			break
		}

		// 过期
		expiredCount++
	}
	if expiredCount <= 0 {
		return
	}

	// 需要移除这么多个
	g.requestJoinExpiredTime = timeutil.LeftShift(g.requestJoinExpiredTime, 0, expiredCount)
	g.requestJoinHeroIds, removeHeroIds = i64.LeftShiftReturnRemovedValues(g.requestJoinHeroIds, 0, expiredCount)
	return
}

func UnmarshalGuild(id int64, proto *server_proto.GuildServerProto, datas *config.ConfigDatas, ctime time.Time) (*Guild, error) {

	g := NewGuild(id, proto.Name, proto.FlagName, datas, timeutil.Unix64(proto.CreateTime))

	g.lastPrestigeRank = proto.LastPrestigeRank
	g.prestige = proto.Prestige
	g.historyMaxPretige = u64.Max(proto.HistoryMaxPrestige, g.prestige)
	c := datas.GetCountryData(proto.Country)
	if c != nil {
		g.SetCountry(c)
	}
	g.prestigeDaily.Unmarshal(proto.PrestigeDaily)
	g.prestigeCoreHourly.Unmarshal(proto.PrestigeCoreHourly)

	g.yinliang = proto.Yinliang
	g.hufu = proto.Hufu

	g.template = datas.GetNpcGuildTemplate(proto.Template)
	g.flagType = proto.FlagType
	g.levelData = datas.GetGuildLevelData(proto.Level)
	g.buildingAmount = proto.BuildingAmount
	g.updateBuildingAmountTime = timeutil.Unix64(proto.UpdateBuildingAmountTime)
	g.upgradeEndTime = timeutil.Unix64(proto.UpgradeEndTime)
	g.cdrTimes = proto.CdrTimes
	g.text = proto.Text
	g.internalText = proto.InternalText
	g.labels = proto.Labels

	g.friendGuildText = proto.FriendGuildText
	g.enemyGuildText = proto.EnemyGuildText

	g.rejectAutoJoin = proto.RejectAutoJoin
	g.requiredHeroLevel = proto.RequiredHeroLevel
	g.requiredJunXianLevel = proto.RequiredJunXianLevel
	g.requiredTowerMaxFloor = proto.RequiredTowerMaxFloor

	g.classNames = proto.ClassNames
	g.classTitleProto = proto.ClassTitle

	var leaderMember *GuildMember
	for _, m := range proto.Members {

		var npcProto *shared_proto.HeroBasicSnapshotProto
		if npcid.IsNpcId(m.Id) {
			if g.template == nil {
				// 如果是Npc帮派，一定要有个 npc模板
				g.template = datas.GuildConfig().GetTemplate()
			}

			if d := g.template.GetNpc(npcid.GetNpcIdSequence(m.Id)); d != nil {
				npcProto = d.EncodeSnapshot(m.Id)
			} else {
				continue
			}
		}

		member := newMember(m.Id, datas.GuildClassLevelData().Must(m.ClassLevel), timeutil.Unix64(m.CreateTime), m, datas, npcProto)
		g.memberMap[m.Id] = member

		if leaderMember == nil || leaderMember.classLevelData.Level < member.classLevelData.Level {
			leaderMember = member
		}
	}
	g.kickMemberCount = proto.KickMemberCount

	if leaderMember != nil {
		g.leaderId = leaderMember.Id()
		leaderMember.classLevelData = datas.GuildClassLevelData().MaxKeyData
	}

	g.changeLeaderId = proto.ChangeLeaderId
	g.changeLeaderTime = timeutil.Unix64(proto.ChangeLeaderTime)
	g.nextChangeNameTime = timeutil.Unix64(proto.NextChangeNameTime)
	g.freeChangeName = proto.FreeChangeName

	g.donateRecords = proto.GetDonateRecords()
	//g.bigEvents = proto.GetBigEvents()
	//g.dynamics = proto.GetDynamics()

	if proto.ImpeachLeader != nil {
		var npcLeaderVote uint64
		npcLeader := g.leaderId
		if !npcid.IsNpcId(npcLeader) {
			npcLeader = 0
		} else {
			if g.template == nil {
				g.template = datas.GuildConfig().GetTemplate()
			}
			npcLeaderVote = g.template.NpcLeaderVote
		}

		candidates := make([]int64, 0)
		for _, cid := range proto.ImpeachLeader.Candidates {
			if _, ok := g.memberMap[cid]; ok {
				candidates = append(candidates, cid)
			}
		}
		if len(candidates) > 1 {
			proto.ImpeachLeader.Candidates = candidates
			if _, ok := g.memberMap[proto.ImpeachLeader.ImpeachMemberId]; !ok {
				proto.ImpeachLeader.ImpeachMemberId = 0
			}

			st := ctime
			if proto.ImpeachLeader.ImpeachStartTime > 0 {
				st = timeutil.Unix64(proto.ImpeachLeader.ImpeachStartTime)
			}

			g.impeachLeader = newImpeachLeader(npcLeader, st, timeutil.Unix64(proto.ImpeachLeader.ImpeachEndTime), proto.ImpeachLeader.Candidates, proto.ImpeachLeader.ImpeachMemberId, npcLeaderVote)

			n := imath.Min(len(proto.ImpeachLeader.VoteHeros), len(proto.ImpeachLeader.VoteTarget))
			for i := 0; i < n; i++ {
				g.impeachLeader.vote(proto.ImpeachLeader.VoteHeros[i], proto.ImpeachLeader.VoteTarget[i])
			}
		}
	}

	n := imath.Min(len(proto.InvateHeroIds), len(proto.InvateExpiredTime))
	for i := 0; i < n; i++ {
		g.invateHeroIds = append(g.invateHeroIds, proto.InvateHeroIds[i])
		g.invateExpiredTime = append(g.invateExpiredTime, timeutil.Unix64(proto.InvateExpiredTime[i]))
	}

	n = imath.Min(len(proto.RequestJoinHeroIds), len(proto.RequestJoinExpiredTime))
	for i := 0; i < n; i++ {
		g.requestJoinHeroIds = append(g.requestJoinHeroIds, proto.RequestJoinHeroIds[i])
		g.requestJoinExpiredTime = append(g.requestJoinExpiredTime, timeutil.Unix64(proto.RequestJoinExpiredTime[i]))
	}

	if proto.GetStatueRealmId() > 0 {
		g.PlaceStatue(proto.GetStatueRealmId())
	}

	for _, v := range proto.SeekHelp {
		g.AddSeekHelp(v)
	}

	if boxData := datas.GetGuildBigBoxData(proto.BigBoxId); boxData != nil {
		g.bigBoxData = boxData
		g.bigBoxEnergy = proto.BigBoxEnergy
	}

	if len(proto.FullBigBoxMemberIds) > 0 {
		if boxData := datas.GetGuildBigBoxData(proto.FullBigBoxId); boxData != nil {
			g.fullBigBoxData = boxData

			for _, id := range proto.FullBigBoxMemberIds {
				g.fullBigBoxMemberIdMap[id] = struct{}{}
			}
		}
	}

	if len(proto.Technologys) > 0 {
		for _, id := range proto.Technologys {
			t := datas.GetGuildTechnologyData(id)
			if t != nil {
				if ex := g.technologyMap[t.Group]; ex != nil {
					if ex.Level >= t.Level {
						continue
					}
				}

				g.technologyMap[t.Group] = t

				if t.BigBox != nil {
					g.setBigBoxData(t.BigBox)
				}
			}
		}
	}

	if proto.UpgradeTechnology != 0 {
		if t := datas.GetGuildTechnologyData(proto.UpgradeTechnology); t != nil {
			g.techUpgradeData = t
		}
	}
	g.techUpgradeEndTime = timeutil.Unix64(proto.TechUpgradeEndTime)
	g.techCdrTimes = proto.TechCdrTimes

	g.isStartResistXiongNuToday = proto.IsStartResistXiongNuToday
	for _, defenderId := range proto.ResistXiongNuDefenders {
		if g.GetMember(defenderId) != nil {
			g.resistXiongNuDefenders = i64.AddIfAbsent(g.resistXiongNuDefenders, defenderId)
			if uint64(len(g.resistXiongNuDefenders)) >= datas.ResistXiongNuMisc().DefenseMemberCount {
				break
			}
		}
	}
	g.unlockResistXiongNuData = datas.ResistXiongNuData().Must(proto.UnlockResistXiongNuLevel)
	g.lastResistXiongNuProto = proto.LastResistXiongNu
	g.lastResistXiongNuFightProto = proto.LastResistXiongNuFightProto

	g.nextUpdatePrestigeTargetTime = timeutil.Unix64(proto.NextUpdatePrestigeTargetTime)

	g.dailyResetTime = timeutil.Unix64(proto.DailyResetTime)
	g.weeklyResetTime = timeutil.Unix64(proto.WeeklyResetTime)

	g.TryUpdateTarget(datas.GuildConfig(), ctime, 0)

	i64.CopyMapTo(g.memberLeaveTimeMap, proto.MemberLeaveTimeMap)

	g.mcWarRecords = proto.Record
	if g.mcWarRecords == nil {
		g.mcWarRecords = &shared_proto.McWarAllRecordProto{}
	}

	g.yinliangRecords = proto.YinliangRecord

	for _, p := range proto.YinliangSendToGuild {
		if p.Guild != nil && p.Send != nil {
			g.yinliangSendToGuildRecords.Add(p)
		}
	}

	g.SendYinliangToMe = proto.YinliangSendToMe
	if g.SendYinliangToMe == nil {
		g.SendYinliangToMe = make(map[int64]*shared_proto.GuildYinliangSendProto)
	}

	g.hostMingcIds = proto.HostMingcIds
	if g.hostMingcIds == nil {
		g.hostMingcIds = make(map[uint64]int64)
	}

	for _, v := range proto.Mark {
		g.AddMark(v, nil)
	}

	if proto.Workshop != nil {
		g.workshop = NewWorkshop(proto.Workshop)
	}

	g.workshopTodayCompleted = proto.WorkshopTodayCompleted
	g.workshopOutput = proto.WorkshopOutput
	g.workshopOutputPrizeCount = proto.WorkshopPrizeCount
	g.workshopBeenHurtTimes = proto.WorkshopBeenHurtTimes

	g.recommendMcBuilds = proto.RecommendMcBuilds

	if len(proto.WeeklyTasks) > 0 {
		for id, progress := range proto.WeeklyTasks {
			data := datas.GetGuildTaskData(u64.FromInt32(id))
			if data == nil {
				continue
			}
			g.weeklyTasks[data.TaskType] = progress
			g.weeklyTaskStageIndexs[data.TaskType] = data.GetStageIndex(progress, 0)
		}
	}

	if proto.ChangeCountryTarget != 0 {
		if target := datas.GetCountryData(proto.ChangeCountryTarget); target != nil {
			g.changeCountryWaitEndTime = proto.ChangeCountryNextTime
			g.changeCountryTarget = target
		}
	}
	g.changeCountryNextTime = proto.ChangeCountryNextTime

	return g, nil
}

func (g *Guild) Marshal() ([]byte, error) {
	return g.EncodeServer().Marshal()
}

func (g *Guild) EncodeServer() *server_proto.GuildServerProto {
	// 创建帮派

	proto := &server_proto.GuildServerProto{}
	proto.Name = g.name
	proto.FlagName = g.flagName
	proto.CreateTime = timeutil.Marshal64(g.createTime)

	if g.template != nil {
		proto.Template = g.template.Id
	}

	proto.FlagType = g.flagType
	proto.Level = g.levelData.Level
	proto.BuildingAmount = g.buildingAmount
	proto.UpdateBuildingAmountTime = timeutil.Marshal64(g.updateBuildingAmountTime)
	proto.UpgradeEndTime = timeutil.Marshal64(g.upgradeEndTime)
	proto.CdrTimes = g.cdrTimes
	proto.Text = g.text
	proto.InternalText = g.internalText
	proto.Labels = g.labels

	proto.LastPrestigeRank = g.lastPrestigeRank
	proto.Prestige = g.prestige
	proto.HistoryMaxPrestige = g.historyMaxPretige
	proto.PrestigeDaily = g.prestigeDaily.Amounts()
	proto.PrestigeCoreHourly = g.prestigeCoreHourly.Amounts()
	if g.country != nil {
		proto.Country = g.country.Id
	}

	proto.Hufu = g.hufu
	proto.Yinliang = g.yinliang

	proto.FriendGuildText = g.friendGuildText
	proto.EnemyGuildText = g.enemyGuildText

	proto.RejectAutoJoin = g.rejectAutoJoin
	proto.RequiredHeroLevel = g.requiredHeroLevel
	proto.RequiredJunXianLevel = g.requiredJunXianLevel
	proto.RequiredTowerMaxFloor = g.requiredTowerMaxFloor

	proto.ClassNames = g.classNames
	proto.ClassTitle = g.classTitleProto

	for _, m := range g.memberMap {
		proto.Members = append(proto.Members, m.encodeServer())
	}
	proto.KickMemberCount = g.kickMemberCount

	proto.ChangeLeaderId = g.changeLeaderId
	proto.ChangeLeaderTime = timeutil.Marshal64(g.changeLeaderTime)
	proto.NextChangeNameTime = timeutil.Marshal64(g.nextChangeNameTime)
	proto.FreeChangeName = g.freeChangeName

	if g.impeachLeader != nil {
		proto.ImpeachLeader = g.impeachLeader.encodeServer()
	}

	proto.InvateHeroIds = g.invateHeroIds
	for _, t := range g.invateExpiredTime {
		proto.InvateExpiredTime = append(proto.InvateExpiredTime, timeutil.Marshal64(t))
	}

	proto.RequestJoinHeroIds = g.requestJoinHeroIds
	for _, t := range g.requestJoinExpiredTime {
		proto.RequestJoinExpiredTime = append(proto.RequestJoinExpiredTime, timeutil.Marshal64(t))
	}

	proto.DonateRecords = g.donateRecords
	//proto.BigEvents = g.bigEvents
	//proto.Dynamics = g.dynamics

	if g.statueRealmId > 0 {
		proto.StatueRealmId = g.statueRealmId
	}

	for _, v := range g.seekHelpMap {
		proto.SeekHelp = append(proto.SeekHelp, v)
	}

	proto.BigBoxId = g.bigBoxData.Id
	proto.BigBoxEnergy = g.bigBoxEnergy

	if g.fullBigBoxData != nil {
		proto.FullBigBoxId = g.fullBigBoxData.Id

		for id := range g.fullBigBoxMemberIdMap {
			proto.FullBigBoxMemberIds = append(proto.FullBigBoxMemberIds, id)
		}
	}

	for _, t := range g.technologyMap {
		proto.Technologys = append(proto.Technologys, t.Id)
	}
	if g.techUpgradeData != nil {
		proto.UpgradeTechnology = g.techUpgradeData.Id
	}
	proto.TechUpgradeEndTime = timeutil.Marshal64(g.techUpgradeEndTime)
	proto.TechCdrTimes = g.techCdrTimes

	proto.NextUpdatePrestigeTargetTime = timeutil.Marshal64(g.nextUpdatePrestigeTargetTime)

	proto.DailyResetTime = timeutil.Marshal64(g.dailyResetTime)
	proto.WeeklyResetTime = timeutil.Marshal64(g.weeklyResetTime)

	proto.IsStartResistXiongNuToday = g.isStartResistXiongNuToday
	proto.ResistXiongNuDefenders = g.resistXiongNuDefenders
	proto.UnlockResistXiongNuLevel = g.unlockResistXiongNuData.Level
	proto.LastResistXiongNu = g.lastResistXiongNuProto
	proto.LastResistXiongNuFightProto = g.lastResistXiongNuFightProto

	proto.MemberLeaveTimeMap = i64.CopyMap(g.memberLeaveTimeMap)

	proto.Record = g.mcWarRecords

	proto.YinliangRecord = g.yinliangRecords

	g.RangeYinliangSendToGuildRecord(func(r *shared_proto.GuildYinliangSendToGuildProto) (toContinue bool) {
		proto.YinliangSendToGuild = append(proto.YinliangSendToGuild, r)
		return true
	})
	proto.YinliangSendToMe = g.SendYinliangToMe

	proto.HostMingcIds = g.hostMingcIds

	for _, v := range g.marks {
		if v != nil {
			proto.Mark = append(proto.Mark, v)
		}
	}

	if g.workshop != nil {
		proto.Workshop = &server_proto.GuildWorkshopServerProto{}
		proto.Workshop.StartTime = g.workshop.startTime
		proto.Workshop.EndTime = g.workshop.endTime
		proto.Workshop.X = int32(g.workshop.x)
		proto.Workshop.Y = int32(g.workshop.y)
		proto.Workshop.IsComplete = g.workshop.isComplete
		proto.Workshop.Prosperity = g.workshop.prosperity
		if g.workshop.log.Length() > 0 {
			g.workshop.log.Range(func(v interface{}) (toContinue bool) {
				log := v.(*shared_proto.GuildWorkshopLogProto)
				proto.Workshop.Log = append(proto.Workshop.Log, log)
				return true
			})
		}
	}

	proto.WorkshopTodayCompleted = g.workshopTodayCompleted
	proto.WorkshopOutput = g.workshopOutput
	proto.WorkshopPrizeCount = g.workshopOutputPrizeCount
	proto.WorkshopBeenHurtTimes = g.workshopBeenHurtTimes
	proto.RecommendMcBuilds = g.recommendMcBuilds

	proto.WeeklyTasks = make(map[int32]uint64)
	for id, progress := range g.weeklyTasks {
		if progress > 0 {
			proto.WeeklyTasks[int32(id)] = progress
		}
	}

	if g.changeCountryTarget != nil {
		proto.ChangeCountryWaitEndTime = g.changeCountryWaitEndTime
		proto.ChangeCountryTarget = g.changeCountryTarget.Id
	}
	proto.ChangeCountryNextTime = g.changeCountryNextTime

	return proto
}

type snapshotService interface {
	Get(int64) *snapshotdata.HeroSnapshot
}

func (g *Guild) encodeMember(member *GuildMember, heroSnapshotService snapshotService, isTodayJoinResistXiongNu IsTodayJoinXiongNuFunc) *shared_proto.GuildMemberProto {

	var proto *shared_proto.HeroBasicSnapshotProto

	if npcid.IsNpcId(member.Id()) {
		proto = member.npcProto
		if proto == nil {
			logrus.WithField("id", member.Id()).WithField("guild_id", g.id).WithField("guild_name", g.name).Error("帮派有Npc成员，Npc成员的npcProto没找到")
			return nil
		}
	} else {
		snapshot := heroSnapshotService.Get(member.Id())
		if snapshot == nil {
			logrus.WithField("id", member.Id()).WithField("guild_id", g.id).WithField("guild_name", g.name).Error("帮派从snapshot里加载失败")
			return nil
		}

		proto = snapshot.EncodeClient()
	}

	return member.encodeClient(proto, isTodayJoinResistXiongNu)
}

type GetGuildRankFunc func(guildId int64, country *country.CountryData) (uint64, uint64)

type IsTodayJoinXiongNuFunc func(heroId int64) bool

func (g *Guild) EncodeClient(details bool, heroSnapshotService snapshotService, f GetGuildRankFunc, isTodayJoinResistXiongNu IsTodayJoinXiongNuFunc) *shared_proto.GuildProto {
	// 创建帮派

	proto := &shared_proto.GuildProto{}
	proto.Id = i64.Int32(g.id)
	proto.Name = g.name
	proto.FlagName = g.flagName
	proto.FlagType = u64.Int32(g.flagType)
	proto.Level = u64.Int32(g.levelData.Level)
	proto.MemberCount = imath.Int32(len(g.memberMap))

	if f != nil {
		rank, rankByCountry := f(g.id, g.country)
		proto.Rank = int32(rank)
		proto.RankByCountry = int32(rankByCountry)
	}

	leader := g.getLeader()
	if leader != nil {
		proto.Leader = g.encodeMember(leader, heroSnapshotService, isTodayJoinResistXiongNu)
	}

	proto.Text = g.text
	proto.InternalText = g.internalText
	proto.Labels = g.labels

	proto.RejectAutoJoin = g.rejectAutoJoin
	proto.RequiredHeroLevel = u64.Int32(g.requiredHeroLevel)
	proto.RequiredJunXianLevel = u64.Int32(g.requiredJunXianLevel)
	proto.RequiredTowerMaxFloor = u64.Int32(g.requiredTowerMaxFloor)

	if g.country != nil {
		proto.PrestigeTarget = u64.Int32(g.country.Id)
	}

	if details {

		proto.Prestige = u64.Int32(g.prestige)
		proto.HistoryMaxPrestige = u64.Int32(g.historyMaxPretige)

		proto.Hufu = u64.Int32(g.hufu)
		proto.Yinliang = u64.Int32(g.yinliang)

		proto.FriendGuildText = g.friendGuildText
		proto.EnemyGuildText = g.enemyGuildText

		// 详情数据
		proto.ClassNames = g.classNames
		proto.ClassTitle = g.classTitleProto
		for _, m := range g.memberMap {
			if m != leader {
				memberProto := g.encodeMember(m, heroSnapshotService, isTodayJoinResistXiongNu)
				if memberProto != nil {
					proto.Members = append(proto.Members, memberProto)
				}
			}
		}
		proto.KickMemberCount = u64.Int32(g.kickMemberCount)

		proto.BuildingAmount = u64.Int32(g.buildingAmount)
		proto.UpgradeEndTime = timeutil.Marshal32(g.upgradeEndTime)
		proto.CdrTimes = u64.Int32(g.cdrTimes)

		if g.changeLeaderId != 0 {
			proto.ChangeLeaderId = idbytes.ToBytes(g.changeLeaderId)
			proto.ChangeLeaderTime = timeutil.Marshal32(g.changeLeaderTime)
		}

		proto.NextChangeNameTime = timeutil.Marshal32(g.nextChangeNameTime)
		proto.FreeChangeName = g.freeChangeName

		proto.ImpeachLeader = g.EncodeImpeachLeader()

		for _, heroId := range g.invateHeroIds {
			if snapshot := heroSnapshotService.Get(heroId); snapshot != nil {
				proto.InvateHero = append(proto.InvateHero, snapshot.EncodeClient())
			}
		}

		for _, heroId := range g.requestJoinHeroIds {
			if snapshot := heroSnapshotService.Get(heroId); snapshot != nil {
				proto.RequestJoinHero = append(proto.RequestJoinHero, snapshot.EncodeClient())
			}
		}

		proto.DonateRecords = g.donateRecords
		//proto.BigEvents = g.bigEvents
		//proto.Dynamics = g.dynamics

		for _, t := range g.doingTargets {
			proto.GuildTargetId = append(proto.GuildTargetId, u64.Int32(t.target.Id))
			proto.GuildTargetStartTime = append(proto.GuildTargetStartTime, timeutil.Marshal32(t.startTime))
			proto.GuildTargetEndTime = append(proto.GuildTargetEndTime, timeutil.Marshal32(t.endTime))
		}

		proto.BigBoxId = u64.Int32(g.bigBoxData.Id)
		proto.BigBoxEnergy = u64.Int32(g.bigBoxEnergy)

		for _, t := range g.technologyMap {
			proto.Technologys = append(proto.Technologys, u64.Int32(t.Id))
		}

		if g.techUpgradeData != nil {
			if g.techUpgradeData.GetPrevLevel() != nil {
				proto.UpgradeTechnology = u64.Int32(g.techUpgradeData.GetPrevLevel().Id)
			} else {
				proto.UpgradeTechnology = u64.Int32(g.techUpgradeData.Id)
			}

			proto.TechUpgradeEndTime = timeutil.Marshal32(g.techUpgradeEndTime)
			proto.TechCdrTimes = u64.Int32(g.techCdrTimes)
		}

		proto.NextUpdatePrestigeTargetTime = timeutil.Marshal32(g.nextUpdatePrestigeTargetTime)

		proto.IsStartResistXiongNuToday = g.isStartResistXiongNuToday
		proto.ResistXiongNuDefenders = make([][]byte, 0, len(g.resistXiongNuDefenders))
		for _, defenderId := range g.resistXiongNuDefenders {
			bytes := idbytes.ToBytes(defenderId)
			proto.ResistXiongNuDefenders = append(proto.ResistXiongNuDefenders, bytes)
		}
		proto.UnlockResistXiongNuLevel = u64.Int32(g.unlockResistXiongNuData.Level)

		if g.lastResistXiongNuProto != nil {
			proto.HasLastResistXiongNu = true
			proto.LastResistXiongNu = g.lastResistXiongNuProto
		}

		proto.MingcHostCount = int32(len(g.hostMingcIds))

		proto.RecommendMcBuilds = u64.Int32Array(g.recommendMcBuilds)

		proto.DailyMcBuildId, proto.DailyMcBuildCount = u64.Map2Int32Array(g.dailyMcBuildCounts)

		// 转国
		if g.changeCountryTarget != nil {
			proto.ChangeCountryWaitEndTime = int32(g.changeCountryWaitEndTime)
			proto.ChangeCountryTarget = u64.Int32(g.changeCountryTarget.Id)
		}
		proto.ChangeCountryNextTime = int32(g.changeCountryNextTime)

		if g.workshop != nil {
			proto.WorkshopX = int32(g.workshop.x)
			proto.WorkshopY = int32(g.workshop.y)
		}
	}

	return proto
}

func (g *Guild) GetWeeklyTasksMsg() pbutil.Buffer {
	if g.weeklyTasksMsg == nil {
		g.weeklyTasksMsg = guild.NewS2cViewTaskProgressMsg(g.weeklyTasksVersion, g.encodeTask()).Static()
	}
	return g.weeklyTasksMsg
}

// 周联盟任务进度
func (g *Guild) encodeTask() []*shared_proto.Int32Pair {
	length := len(g.weeklyTasks)
	if length <= 0 {
		return []*shared_proto.Int32Pair{}
	}
	p := make([]*shared_proto.Int32Pair, 0, length)
	for id, progress := range g.weeklyTasks {
		p = append(p, &shared_proto.Int32Pair{
			Key:   int32(id),
			Value: u64.Int32(progress),
		})
	}
	return p
}

func (g *Guild) EncodeImpeachLeader() *shared_proto.GuildImpeachProto {
	if g.impeachLeader == nil {
		return nil
	}

	return g.impeachLeader.encodeClient(g.memberMap)
}

func (g *Guild) Id() int64 {
	return g.id
}

func (g *Guild) Name() string {
	return g.name
}

func (g *Guild) FlagName() string {
	return g.flagName
}

func (g *Guild) CreateTime() time.Time {
	return g.createTime
}

func (g *Guild) IsFreeChangeName() bool {
	return g.freeChangeName
}

func (g *Guild) NextChangeNameTime() time.Time {
	return g.nextChangeNameTime
}

func (g *Guild) ChangeName(name, flagName string, nextChangeNameTime time.Time) {
	g.name = name
	g.flagName = flagName
	g.nextChangeNameTime = nextChangeNameTime

	g.freeChangeName = false
}

func (g *Guild) FlagType() uint64 {
	return g.flagType
}

func (g *Guild) SetFlagType(toSet uint64) {
	g.flagType = toSet
}

func (g *Guild) LevelData() *guild_data.GuildLevelData {
	return g.levelData
}

func (g *Guild) SetLevelData(toSet *guild_data.GuildLevelData) {
	g.levelData = toSet
}

func (g *Guild) GetText() string {
	return g.text
}

func (g *Guild) SetText(toSet string) {
	g.text = toSet
}

func (g *Guild) GetInternalText() string {
	return g.internalText
}

func (g *Guild) SetInternalText(toSet string) {
	g.internalText = toSet
}

func (g *Guild) SetFriendGuildText(toSet string) {
	g.friendGuildText = toSet
}

func (g *Guild) SetEnemyGuildText(toSet string) {
	g.enemyGuildText = toSet
}

func (g *Guild) GetHufu() uint64 {
	return g.hufu
}

func (g *Guild) SetHufu(toSet uint64) {
	g.hufu = toSet
}

func (g *Guild) AddHufu(toAdd uint64) uint64 {
	g.hufu += toAdd
	return g.hufu
}

func (g *Guild) ReduceHufu(toReduce uint64) (new uint64, ok bool) {
	if g.hufu < toReduce {
		return g.hufu, false
	}

	g.hufu = u64.Sub(g.hufu, toReduce)
	return g.hufu, true
}

func (g *Guild) HasEnoughYinliang(amount uint64) bool {
	return g.yinliang >= amount
}

func (g *Guild) GetYinliang() uint64 {
	return g.yinliang
}

func (g *Guild) SetYinliang(toSet uint64) {
	g.yinliang = toSet
}

func (g *Guild) AddYinliang(toAdd uint64) uint64 {
	g.yinliang += toAdd
	return g.yinliang
}

func (g *Guild) ReduceYinliang(toReduce uint64) (new uint64, ok bool) {
	if g.yinliang < toReduce {
		return g.yinliang, false
	}

	g.yinliang = u64.Sub(g.yinliang, toReduce)
	return g.yinliang, true
}

var (
	yinliangRecordMaxLen = 500
)

func (g *Guild) AddYinliangRecord(text string, ctime time.Time) {
	p := &shared_proto.GuildYinliangRecordProto{}
	p.Time = timeutil.Marshal32(ctime)
	p.Text = text

	if len(g.yinliangRecords) >= yinliangRecordMaxLen {
		g.yinliangRecords = g.yinliangRecords[1:]
	}
	g.yinliangRecords = append(g.yinliangRecords, p)
	g.yinliangRecordMsg = g.BuildYinliangRecordMsg()
}

func (g *Guild) BuildYinliangRecordMsg() pbutil.Buffer {
	return guild.NewS2cViewYinliangRecordMsg(&shared_proto.GuildAllYinliangRecordProto{Record: g.yinliangRecords}).Static()
}

func (g *Guild) GetLastPrestigeRank() uint64 {
	return g.lastPrestigeRank
}

func (g *Guild) SetLastPrestigeRank(rank uint64) {
	g.lastPrestigeRank = rank
}

func (g *Guild) GetPrestige() uint64 {
	return g.prestige
}

func (g *Guild) SetPrestige(target uint64) {
	g.prestige = target
	//if g.prestige > g.historyMaxPretige {
	//	g.historyMaxPretige = g.prestige
	//}
}

func (g *Guild) AddPrestige(toAdd uint64) uint64 {
	g.prestige += toAdd
	g.prestigeDaily.Add(toAdd)
	return g.prestige
}

func (g *Guild) RefreshHistoryMaxPrestige(prestige uint64) uint64 {
	if prestige > g.historyMaxPretige {
		g.historyMaxPretige = prestige
	}
	return g.historyMaxPretige
}

func (g *Guild) GetHistoryMaxPrestige() uint64 {
	return g.historyMaxPretige
}

func (g *Guild) GetPrestigeCore() uint64 {
	return g.prestigeCoreHourly.Total()
}

func (g *Guild) AddPrestigeCore(toAdd uint64) {
	g.prestigeCoreHourly.Add(toAdd)
}

func (g *Guild) Country() *country.CountryData {
	return g.country
}

func (g *Guild) CountryId() uint64 {
	return g.country.Id
}

func (g *Guild) SetCountry(target *country.CountryData) {
	g.country = target
}

func (g *Guild) GetNextUpdatePrestigeTargetTime() time.Time {
	return g.nextUpdatePrestigeTargetTime
}

func (g *Guild) SetNextUpdatePrestigeTargetTime(toSet time.Time) {
	g.nextUpdatePrestigeTargetTime = toSet
}

func (g *Guild) GetPrestigeHeroCount() int {
	return len(g.prestigeHeroIdMap)
}

func (g *Guild) IsPrestigeHero(heroId int64) bool {
	_, exist := g.prestigeHeroIdMap[heroId]
	return exist
}

func (g *Guild) PutPrestigeHero(heroId int64) {
	g.prestigeHeroIdMap[heroId] = struct{}{}
}

func (g *Guild) Statue() (realmId int64) {
	return g.statueRealmId
}

func (g *Guild) StatueCacheMsg() pbutil.Buffer {
	return g.statueCacheMsg
}

func (g *Guild) PlaceStatue(realmId int64) {
	g.statueRealmId = realmId
	g.statueCacheMsg = guild.NewS2cGuildStatueMsg(i64.Int32(realmId)).Static()
}

func (g *Guild) TakeBackStatue() {
	g.statueRealmId = 0
	g.statueCacheMsg = nil
}

func (g *Guild) GetClassNames() []string {
	return g.classNames
}

func (g *Guild) SetClassNames(toSet []string) {
	g.classNames = toSet
}

func (g *Guild) GetCustomClassTitle() []string {
	if g.classTitleProto != nil {
		return g.classTitleProto.CustomClassTitleName
	}

	return nil
}

func (g *Guild) SetClassTitle(proto *shared_proto.GuildClassTitleProto) {
	g.classTitleProto = proto
}

func (g *Guild) GetLabels() []string {
	return g.labels
}

func (g *Guild) SetLabels(toSet []string) {
	g.labels = toSet
}

func (g *Guild) SetJoinCondition(rejectAutoJoin bool, requiredHeroLevel, requiredJunXianLevel, requiredTowerMaxFloor uint64) {
	g.rejectAutoJoin = rejectAutoJoin
	g.requiredHeroLevel = requiredHeroLevel
	g.requiredJunXianLevel = requiredJunXianLevel
	g.requiredTowerMaxFloor = requiredTowerMaxFloor
}

func (g *Guild) IsRejectAutoJoin() bool {
	return g.rejectAutoJoin
}

func (g *Guild) GetRequiredHeroLevel() uint64 {
	return g.requiredHeroLevel
}

func (g *Guild) GetRequiredJunXianLevel() uint64 {
	return g.requiredJunXianLevel
}

func (g *Guild) GetRequiredTowerMaxFloor() uint64 {
	return g.requiredTowerMaxFloor
}

func (g *Guild) GetBuildingAmount() uint64 {
	return g.buildingAmount
}

func (g *Guild) GetDonateHeroCount() int {
	return len(g.donateHeroIdMap)
}

func (g *Guild) IsDonate(heroId int64) bool {
	_, exist := g.donateHeroIdMap[heroId]
	return exist
}

func (g *Guild) AddDonateBuildingAmount(toAdd uint64, toSetTime time.Time, donateHeroId int64) {
	g.AddBuildingAmount(toAdd, toSetTime)

	if donateHeroId != 0 {
		g.donateHeroIdMap[donateHeroId] = struct{}{}
	}
}

func (g *Guild) AddBuildingAmount(toAdd uint64, toSetTime time.Time) {
	g.buildingAmount += toAdd
	g.updateBuildingAmountTime = toSetTime
}

func (g *Guild) ReduceBuildingAmount(toReduce uint64, toSetTime time.Time) {
	g.buildingAmount = u64.Sub(g.buildingAmount, toReduce)
	g.updateBuildingAmountTime = toSetTime
}

func (g *Guild) UpdateBuildingAmountTime() time.Time {
	return g.updateBuildingAmountTime
}

func (g *Guild) GetUpgradeEndTime() time.Time {
	return g.upgradeEndTime
}

func (g *Guild) SetUpgradeEndTime(toSet time.Time) {
	g.upgradeEndTime = toSet
}

func (g *Guild) GetCdrTimes() uint64 {
	return g.cdrTimes
}

func (g *Guild) IncCdrTimes() {
	g.cdrTimes++
}

func (g *Guild) TryUpgradeLevel(ctime time.Time) bool {

	if timeutil.IsZero(g.upgradeEndTime) {
		return false
	}

	if ctime.Before(g.upgradeEndTime) {
		// 时间还没到
		return false
	}

	if nextLevel := g.levelData.NextLevel(); nextLevel != nil {
		g.SetLevelData(nextLevel)
	}

	g.upgradeEndTime = time.Time{}
	g.cdrTimes = 0

	return true
}

func (g *Guild) TryUpgradeTechnology(ctime time.Time) (ok bool, toUpgrade *guild_data.GuildTechnologyData) {

	if timeutil.IsZero(g.techUpgradeEndTime) {
		return
	}

	if ctime.Before(g.techUpgradeEndTime) {
		// 时间还没到
		return
	}

	toUpgrade = g.techUpgradeData
	if toUpgrade != nil {
		g.SetTechnology(toUpgrade)
	}

	g.techUpgradeData = nil
	g.techUpgradeEndTime = time.Time{}
	g.techCdrTimes = 0

	return true, toUpgrade
}

func (g *Guild) GetKickMemberCount() uint64 {
	return g.kickMemberCount
}

func (g *Guild) IncKickMemberCount() {
	g.kickMemberCount++
}

func (g *Guild) MemberCount() int {
	return len(g.memberMap)
}

func (g *Guild) EmptyMemberCount() uint64 {
	return u64.Sub(g.levelData.MemberCount, uint64(len(g.memberMap)))
}

func (g *Guild) GetMember(id int64) *GuildMember {
	return g.memberMap[id]
}

func (g *Guild) RemoveMember(id, leaveTime int64) {
	delete(g.memberMap, id)

	// 如果是帮主转让候选人，取消帮主转让
	if id == g.GetChangeLeaderId() {
		g.CancelChangeLeader()
	}

	// 如果当前有盟主弹劾
	if g.impeachLeader != nil {
		// 如果是候选人，删掉候选人，将投票给他的人，设置为未投票
		// 如果不是候选人，人变少了，看下是否弹劾成功

		// 删掉这个人的投票，如果是候选人，删掉候选人，以及删掉所有投给他的票
		g.impeachLeader.removeMember(id)
	}

	g.RemoveResistXiongNuDefender(id)

	g.AddLeaveMemver(id, leaveTime)
}

func (g *Guild) AddMember(toAdd *GuildMember) {
	g.memberMap[toAdd.Id()] = toAdd
}

func (g *Guild) WalkMember(f func(member *GuildMember)) {
	for _, member := range g.memberMap {
		f(member)
	}
}

func (g *Guild) SetLeader(id int64) {
	g.leaderId = id
	g.leaderOfflineTime = time.Time{}
}

func (g *Guild) GetLeaderOfflineTime() time.Time {
	return g.leaderOfflineTime
}

func (g *Guild) UpdateLeaderOfflineTime(toUpdate time.Time) {
	g.leaderOfflineTime = toUpdate
}

func (g *Guild) HasChangeLeaderCountDown() bool {
	if g.changeLeaderId == 0 {
		return false
	}

	newleader := g.memberMap[g.changeLeaderId]
	if newleader != nil && newleader.Id() != g.leaderId && !timeutil.IsZero(g.changeLeaderTime) {
		return true
	}
	g.CancelChangeLeader()

	return false
}

func (g *Guild) TryTickChangeLeader(ctime time.Time, lowestClassData, leaderClassData *guild_data.GuildClassLevelData) bool {
	if g.changeLeaderId == 0 {
		return false
	}

	// 时间没到
	if ctime.Before(g.changeLeaderTime) {
		return false
	}

	newleader := g.memberMap[g.changeLeaderId]
	g.CancelChangeLeader() // 清掉数据
	if newleader == nil {
		return false
	}

	return g.changeLeader(newleader, lowestClassData, leaderClassData)
}

func (g *Guild) ChangeLeaderCountDown(changeLeaderId int64, changeLeaderTime time.Time) {
	g.changeLeaderId = changeLeaderId
	g.changeLeaderTime = changeLeaderTime
}

func (g *Guild) GetChangeLeaderId() int64 {
	return g.changeLeaderId
}

func (g *Guild) CancelChangeLeader() {
	g.changeLeaderId = 0
	g.changeLeaderTime = time.Time{}
}

func (g *Guild) LeaderId() int64 {
	return g.leaderId
}

func (g *Guild) IsNpcLeader() bool {
	return npcid.IsNpcId(g.leaderId)
}

func (g *Guild) GetMaxClassLevelMember() *GuildMember {
	return g.getLeader()
}

func (g *Guild) getLeader() *GuildMember {
	m := g.memberMap[g.leaderId]
	if m != nil {
		return m
	}

	// 找一个
	var leaderMember *GuildMember
	for _, member := range g.memberMap {
		if leaderMember == nil || leaderMember.classLevelData.Level < member.classLevelData.Level {
			leaderMember = member
		}
	}

	return leaderMember
}

func (g *Guild) ResetLeader() {
	leader := g.getLeader()
	if leader != nil {
		g.SetLeader(leader.Id())
	}
}

func (g *Guild) IsFull() bool {
	return uint64(len(g.memberMap)) >= g.levelData.MemberCount
}

func (g *Guild) IsClassFull(classLevelData *guild_data.GuildClassLevelData) bool {

	c := g.levelData.GetClassMemberCount(classLevelData.Level)
	if c > 0 {
		for _, v := range g.memberMap {
			if v.ClassLevelData() != classLevelData {
				continue
			}

			c--
			if c <= 0 {
				return true
			}
		}
	}

	return false
}

func (g *Guild) AllUserMemberIds() []int64 {

	ids := make([]int64, 0, len(g.memberMap))
	for k := range g.memberMap {
		if !npcid.IsNpcId(k) {
			ids = append(ids, k)
		}
	}

	return ids
}

func (g *Guild) GetContribution7Slice() []*GuildMember {

	var array []*GuildMember
	for _, m := range g.memberMap {
		array = append(array, m)
	}

	sort.Sort(contribution7Slice(array))

	return array
}

func (g *Guild) GetImpeachLeader() *impeach_leader {
	return g.impeachLeader
}

func (g *Guild) SetImpeachLeader(toSet *impeach_leader) {
	g.impeachLeader = toSet
}

func (g *Guild) StartImpeachLeader(impeachMemberId int64, impeachStartTime, impeachEndTime time.Time, candidateCount uint64) {
	var npcLeaderVote uint64
	if g.template != nil {
		npcLeaderVote = g.template.NpcLeaderVote
	}
	g.impeachLeader = createImpeachLeader(g.leaderId, impeachMemberId, impeachStartTime, impeachEndTime, g.memberMap, candidateCount, npcLeaderVote)
}

func (g *Guild) TryTickImpeachLeader(ctime time.Time, lowestClassData, leaderClassData *guild_data.GuildClassLevelData) (changed, success bool) {

	// 没有在弹劾
	if g.impeachLeader == nil {
		return
	}

	// 弹劾时间未到
	if ctime.Before(g.impeachLeader.impeachEndTime) {
		return
	}

	changed = true

	// 弹劾时间到了，找到还剩的那个家伙出来
	// 找到票数最高的那个家伙出来
	isNpcLeader := g.IsNpcLeader()
	newLeader := g.impeachLeader.getMaxScoreCandidate(g.memberMap)
	if newLeader != nil {
		if g.changeLeader(newLeader, lowestClassData, leaderClassData) {
			if isNpcLeader {
				g.kickAllNpc()
			}
			success = true
		}
	}

	g.impeachLeader = nil

	return
}

func (g *Guild) TryVote(selfId, voteTargetId int64) bool {
	if g.impeachLeader != nil {
		g.impeachLeader.vote(selfId, voteTargetId)
		return false
	}

	return false
}

func (g *Guild) TryImpeachLeader(lowestClassData, leaderClassData *guild_data.GuildClassLevelData) bool {
	if g.impeachLeader == nil {
		return false
	}

	// 如果只剩一个候选人，弹劾结束
	if len(g.impeachLeader.candidates) <= 1 {

		for _, newLeaderId := range g.impeachLeader.candidates {
			newLeader := g.memberMap[newLeaderId]
			if newLeader != nil {
				// 已经产生新的盟主，更新盟主
				if g.changeLeader(newLeader, lowestClassData, leaderClassData) {
					// 将所有的NPC 踢掉
					g.kickAllNpc()
					return true
				}
			}

			return false
		}

		return false
	}

	// 如果是非NPC弹劾，重新计算一下票数
	return g.tryImpeachLeader0(lowestClassData, leaderClassData)
}

func (g *Guild) tryImpeachLeader0(lowestClassData, leaderClassData *guild_data.GuildClassLevelData) bool {

	// 重新计算一遍投票情况，看下是否弹劾成功
	newLeader := g.impeachLeader.tryImpeach(g.memberMap)
	if newLeader != nil {
		// 已经产生新的盟主，更新盟主
		if g.changeLeader(newLeader, lowestClassData, leaderClassData) {
			// 将所有的NPC 踢掉
			g.kickAllNpc()
			return true
		}
	}

	return false
}

func (g *Guild) kickAllNpc() {
	for _, m := range g.memberMap {
		if npcid.IsNpcId(m.Id()) {
			delete(g.memberMap, m.Id())
		}
	}

	// 免费改名
	g.freeChangeName = true
}

func (g *Guild) changeLeader(newLeader *GuildMember, lowestClassData, leaderClassData *guild_data.GuildClassLevelData) (success bool) {
	if newLeader.Id() == g.leaderId {
		return
	}

	// 先取消原来帮助的职务
	originLeader := g.memberMap[g.leaderId]
	if originLeader != nil {
		originLeader.classLevelData = lowestClassData
	}

	g.SetLeader(newLeader.Id())
	newLeader.classLevelData = leaderClassData

	g.impeachLeader = nil

	return true
}

func (g *Guild) GmResetDaily(contributionDay int, resetTime time.Time) {
	g.resetDaily(resetTime, contributionDay)
	g.SetChanged()
}

func (g *Guild) GmSetResetTime(toSet time.Time) {
	g.dailyResetTime = toSet
}

func (g *Guild) ResetHourly() {
	g.prestigeCoreHourly.ResetDaily(1)
}

func (g *Guild) ResetWeekly(resetTime time.Time) bool {
	if !g.weeklyResetTime.Before(resetTime) {
		return false
	}
	logrus.Debug("联盟每周重置", g.name, g.weeklyResetTime, resetTime)
	g.weeklyResetTime = resetTime
	g.resetWeekly()
	return true
}

func (g *Guild) ResetDaily(contributionDay int, resetTime time.Time) bool {

	if !g.dailyResetTime.Before(resetTime) {
		return false
	}

	logrus.Debug("联盟每日重置", g.name, g.dailyResetTime, resetTime)
	g.resetDaily(resetTime, contributionDay)

	return true
}

func (g *Guild) resetWeekly() {
	g.weeklyTasksVersion++
	g.weeklyTasksMsg = nil
	if len(g.weeklyTasks) > 0 {
		g.weeklyTasks = make(map[server_proto.GuildTaskType]uint64)
	}
	if len(g.weeklyTaskStageIndexs) > 0 {
		g.weeklyTaskStageIndexs = make(map[server_proto.GuildTaskType]int)
	}
}

func (g *Guild) resetDaily(resetTime time.Time, contributionDay int) {

	g.dailyResetTime = resetTime
	// 每日踢人个数
	g.kickMemberCount = 0

	// 每日添加建设值
	g.donateHeroIdMap = make(map[int64]struct{})

	// 每日重置声望member
	g.prestigeHeroIdMap = make(map[int64]struct{})

	// 重置是否有开启抗击匈奴
	g.isStartResistXiongNuToday = false

	g.prestigeDaily.ResetDaily(1)

	for _, m := range g.memberMap {
		if !npcid.IsNpcId(m.Id()) {
			m.resetDaily(contributionDay)
		}
	}

	if resetTime.Weekday() == time.Monday {
		for _, sendProto := range g.SendYinliangToMe {
			sendProto.WeeklySend = 0
		}
	}

	g.workshopOutput = 0
	g.workshopTodayCompleted = false
	g.workshopOutputPrizeCount = 0
	g.workshopBeenHurtTimes = 0
}

// 添加捐献记录
func (g *Guild) AddDonateRecord(record *shared_proto.GuildDonateRecordProto, maxDonateRecordCount uint64) {
	if uint64(len(g.donateRecords)) >= maxDonateRecordCount {
		// 删除掉最左边的
		deleteCount := uint64(len(g.donateRecords)) - maxDonateRecordCount + 1
		copy(g.donateRecords, g.donateRecords[deleteCount:])
		g.donateRecords = g.donateRecords[:maxDonateRecordCount-1]
	}

	g.donateRecords = append(g.donateRecords, record)

	g.SetChanged()
}

// 今天是否有开启抗击匈奴
func (g *Guild) IsStartResistXiongNuToday() bool {
	return g.isStartResistXiongNuToday
}

func (g *Guild) SetIsStartResistXiongNuToday(toSet bool) {
	g.isStartResistXiongNuToday = toSet
}

// 抗击匈奴的防守者
func (g *Guild) ResistXiongNuDefenders() []int64 {
	return g.resistXiongNuDefenders
}

func (g *Guild) IsResistXiongNuDefenders(id int64) bool {
	return i64.Contains(g.resistXiongNuDefenders, id)
}

// 抗击匈奴的等级数据
func (g *Guild) UnlockResistXiongNuData() *xiongnu.ResistXiongNuData {
	return g.unlockResistXiongNuData
}

func (g *Guild) SetUnlockResistXiongNuData(toSet *xiongnu.ResistXiongNuData) {
	g.unlockResistXiongNuData = toSet
}

// 添加
func (g *Guild) AddResistXiongNuDefender(toAdd int64) {
	g.resistXiongNuDefenders = i64.AddIfAbsent(g.resistXiongNuDefenders, toAdd)
}

// 移除
func (g *Guild) RemoveResistXiongNuDefender(toRemove int64) {
	g.resistXiongNuDefenders = i64.RemoveIfPresent(g.resistXiongNuDefenders, toRemove)
}

// 设置最近的一场战斗
func (g *Guild) SetLastResistXiongNuProto(toSet *shared_proto.LastResistXiongNuProto) {
	g.lastResistXiongNuProto = toSet
}

func (g *Guild) SetLastResistXiongNuFightProto(toSet *shared_proto.ResistXiongNuFightProto) {
	g.lastResistXiongNuFightProto = toSet
	g.lastResistXiongNuFightProtoMsg = nil
}

func (g *Guild) GetLastResistXiongNuFightMsg() pbutil.Buffer {

	if g.lastResistXiongNuFightProtoMsg == nil {
		if g.lastResistXiongNuFightProto != nil {
			g.lastResistXiongNuFightProtoMsg = xiongnu2.NewS2cGetXiongNuFightInfoMsg(
				must.Marshal(g.lastResistXiongNuFightProto)).Static()
		}
	}

	return g.lastResistXiongNuFightProtoMsg
}

//// 添加大事记记录
//func (g *Guild) AddBigEvent(event *shared_proto.GuildBigEventProto, maxBigEventCount uint64) {
//	if uint64(len(g.bigEvents)) >= maxBigEventCount {
//		// 删除掉最左边的
//		deleteCount := uint64(len(g.bigEvents)) - maxBigEventCount + 1
//		copy(g.bigEvents, g.bigEvents[deleteCount:])
//		g.bigEvents = g.bigEvents[:maxBigEventCount-1]
//	}
//
//	g.bigEvents = append(g.bigEvents, event)
//
//	g.SetChanged()
//}
//
//// 添加联盟动态
//func (g *Guild) AddDynamic(proto *shared_proto.GuildDynamicProto, maxDynamicCount uint64) {
//	if uint64(len(g.dynamics)) >= maxDynamicCount {
//		// 删除掉最左边的
//		deleteCount := uint64(len(g.dynamics)) - maxDynamicCount + 1
//		copy(g.dynamics, g.dynamics[deleteCount:])
//		g.dynamics = g.dynamics[:maxDynamicCount-1]
//	}
//
//	g.dynamics = append(g.dynamics, proto)
//
//	g.SetChanged()
//}

func (g *Guild) TryUpdateTarget(c *singleton.GuildConfig, ctime time.Time, targetType shared_proto.GuildTargetType) bool {

	if targetType != 0 {
		update := false
		isDoingTarget := false
		for _, t := range g.doingTargets {
			if t.target.TargetType == targetType {
				isDoingTarget = true
				if t.isFinished(g, ctime) {
					update = true
					break
				}
			}
		}

		if isDoingTarget && !update {
			return false
		}
	}

	// 找到第一个符合条件的情况
	var ts []*targetWithEndTime
	for _, ta := range c.GetGuildTargetGroups() {
		for _, t := range ta {
			if ok, startTime, endTime := g.IsDoingTarget(t, c, ctime); ok {
				ts = append(ts, &targetWithEndTime{
					target:    t,
					startTime: startTime,
					endTime:   endTime,
				})
				break
			}
		}
	}
	sort.Sort(targetWithEndTimeSlice(ts))

	if len(g.doingTargets) != len(ts) {
		g.doingTargets = ts
		return true
	}

	for i, t := range ts {
		doing := g.doingTargets[i]
		if t.target != doing.target || !t.endTime.Equal(doing.endTime) {
			g.doingTargets = ts
			return true
		}
	}

	return false
}

func (g *Guild) IsDoingTarget(t *guild_data.GuildTarget, c *singleton.GuildConfig, ctime time.Time) (bool, time.Time, time.Time) {

	switch t.TargetType {
	case shared_proto.GuildTargetType_GuildLevelUp:
		return g.levelData.Level < t.Target, time.Time{}, time.Time{}
	case shared_proto.GuildTargetType_PrestigeUp:
		return g.historyMaxPretige < t.Target, time.Time{}, time.Time{}
	case shared_proto.GuildTargetType_ImpeachNpcLeader:
		if g.IsNpcLeader() && g.impeachLeader != nil {
			return true, g.impeachLeader.impeachStartTime, g.impeachLeader.impeachEndTime
		}
	case shared_proto.GuildTargetType_ImpeachUserLeader:
		if !g.IsNpcLeader() && g.impeachLeader != nil {
			return true, g.impeachLeader.impeachStartTime, g.impeachLeader.impeachEndTime
		}
	case shared_proto.GuildTargetType_UpdateMemverClass:
		if g.IsNpcLeader() {
			// 固定2点刷新
			return true, c.GetPrevNpcSetClassLevelTime(ctime), c.GetNextNpcSetClassLevelTime(ctime)
		}
	case shared_proto.GuildTargetType_UserLeaderUseless:
		if !g.IsNpcLeader() && g.impeachLeader == nil && !timeutil.IsZero(g.leaderOfflineTime) && ctime.After(g.leaderOfflineTime.Add(c.ImpeachUserLeaderOffline)) {
			// 检查盟主是否离线超过弹劾时间
			return true, time.Time{}, time.Time{}
		}
	case shared_proto.GuildTargetType_GuildChangeCountry:
		if g.changeCountryTarget != nil {
			endTime := timeutil.Unix64(g.changeCountryWaitEndTime)
			return true, endTime.Add(-c.GuildChangeCountryWaitDuration), endTime
		}
	}

	return false, time.Time{}, time.Time{}
}

type targetWithEndTime struct {
	target    *guild_data.GuildTarget
	startTime time.Time // 用来给客户端做进度条用的
	endTime   time.Time
}

func (tw *targetWithEndTime) isFinished(g *Guild, ctime time.Time) bool {

	// 有CD的，时间到了，直接完成
	if !timeutil.IsZero(tw.endTime) && !ctime.Before(tw.endTime) {
		return true
	}

	t := tw.target
	switch t.TargetType {
	case shared_proto.GuildTargetType_GuildLevelUp:
		return g.levelData.Level >= t.Target
	case shared_proto.GuildTargetType_PrestigeUp:
		return g.historyMaxPretige >= t.Target
	case shared_proto.GuildTargetType_ImpeachNpcLeader:
		return !g.IsNpcLeader() || g.impeachLeader == nil
	case shared_proto.GuildTargetType_ImpeachUserLeader:
		return g.IsNpcLeader() || g.impeachLeader == nil
	case shared_proto.GuildTargetType_UpdateMemverClass:
		return !g.IsNpcLeader()
	case shared_proto.GuildTargetType_UserLeaderUseless:
		return g.IsNpcLeader() || g.impeachLeader != nil
	case shared_proto.GuildTargetType_GuildChangeCountry:
		return g.changeCountryTarget == nil
	}

	return false
}

type targetWithEndTimeSlice []*targetWithEndTime

func (p targetWithEndTimeSlice) Len() int { return len(p) }
func (p targetWithEndTimeSlice) Less(i, j int) bool {
	return p[i].target.OrderAmount < p[j].target.OrderAmount
}
func (p targetWithEndTimeSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (g *Guild) GetMemverLeaveMemver(memberId int64) int64 {
	return g.memberLeaveTimeMap[memberId]
}

func (g *Guild) AddLeaveMemver(memberId, leaveTime int64) {
	if leaveTime > 0 {
		g.memberLeaveTimeMap[memberId] = leaveTime
	}
}

func (g *Guild) RemoveExpireMemberLeaveTime(t int64) {
	for k, v := range g.memberLeaveTimeMap {
		if v < t {
			delete(g.memberLeaveTimeMap, k)
		}
	}
}

func (g *Guild) AddMcWarRecord(r *shared_proto.McWarRecordProto) {
	g.mcWarRecords.Record = append(g.mcWarRecords.Record, r)
}

func (g *Guild) McWarRecord() (p *shared_proto.McWarAllRecordProto) {
	return g.mcWarRecords
}

func (g *Guild) YinliangRecordMsg() pbutil.Buffer {
	return g.yinliangRecordMsg
}

func (g *Guild) UpdateYinliangSendToGuild(receiver *shared_proto.GuildBasicProto, sendProto *shared_proto.GuildYinliangSendProto) (exist *shared_proto.GuildYinliangSendToGuildProto) {
	var toReturn *shared_proto.GuildYinliangSendToGuildProto
	g.RangeYinliangSendToGuildRecord(func(p *shared_proto.GuildYinliangSendToGuildProto) (toContinue bool) {
		if p.Guild != nil && p.Guild.Id == receiver.Id {
			p.Guild = receiver
			p.Send = sendProto
			toReturn = p
			return false
		}
		return true
	})

	if toReturn != nil {
		return toReturn
	}

	toAdd := &shared_proto.GuildYinliangSendToGuildProto{
		Guild: receiver,
		Send:  sendProto,
	}

	g.addYinliangSendToGuild(toAdd)

	return nil
}

func (g *Guild) addYinliangSendToGuild(p *shared_proto.GuildYinliangSendToGuildProto) {
	g.yinliangSendToGuildRecords.Add(p)
}

func (g *Guild) RangeYinliangSendToGuildRecord(f func(r *shared_proto.GuildYinliangSendToGuildProto) (toContinue bool)) {
	g.yinliangSendToGuildRecords.Range(func(v interface{}) (toContinue bool) {
		if v != nil {
			if r, ok := v.(*shared_proto.GuildYinliangSendToGuildProto); ok {
				return f(r)
			}
		}
		return true
	})
}

// 增加坐拥名城id记录，返回当前坐拥名城数
func (g *Guild) AddHostMingc(mcId uint64) uint64 {
	g.hostMingcIds[mcId] = 0
	return u64.FromInt(len(g.hostMingcIds))
}

func (g *Guild) DelHostMingc(mcId uint64) {
	delete(g.hostMingcIds, mcId)
}

// 获取昨天的声望值
func (g *Guild) YesterdayPrestige() uint64 {
	return g.prestigeDaily.GetAmount(1)
}

func (g *Guild) SetMcBuildCount(mcId, count uint64) {
	g.dailyMcBuildCounts[mcId] = count
}

func (g *Guild) McBuildCount(mcId uint64) uint64 {
	return g.dailyMcBuildCounts[mcId]
}

// 返回是否激活新的stage
func (g *Guild) AddGuildTaskProgress(data *guild_data.GuildTaskData, progress uint64) bool {
	if data == nil {
		return false
	}
	// 刷新版本和重置公共消息
	g.weeklyTasksVersion++
	g.weeklyTasksMsg = nil
	// 增加进度值
	progress += g.weeklyTasks[data.TaskType]
	g.weeklyTasks[data.TaskType] = progress
	// 获得新的下标作比较，如果改变了，就是有起码大于等于1的激活
	oldIndex := g.weeklyTaskStageIndexs[data.TaskType]
	newIndex := data.GetStageIndex(progress, oldIndex)
	if newIndex != oldIndex {
		g.weeklyTaskStageIndexs[data.TaskType] = newIndex
		return true
	}
	return false
}

func (g *Guild) GetGuildTaskStageIndex(t server_proto.GuildTaskType) int {
	return g.weeklyTaskStageIndexs[t]
}

func (g *Guild) GetGuildTaskProgress(t server_proto.GuildTaskType) uint64 {
	return g.weeklyTasks[t]
}

func (g *Guild) GetGuildTaskVersion() int32 {
	return g.weeklyTasksVersion
}

func (g *Guild) getGuildTaskCompletedStageCount(data *guild_data.GuildTaskData) uint64 {
	progress := g.GetGuildTaskProgress(server_proto.GuildTaskType(u64.Int32(data.Id)))
	if progress <= 0 {
		return 0
	}
	return data.GetCompletedStageCount(progress)
}

func (g *Guild) GetAllGuildTasksCompletedStageCount(datas []*guild_data.GuildTaskData) uint64 {
	count := uint64(0)
	for _, data := range datas {
		count += g.getGuildTaskCompletedStageCount(data)
	}
	return count
}
