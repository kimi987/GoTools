package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/body"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/head"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/config/season"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"strings"
	"time"
	"github.com/lightpaw/male7/entity/heroid"
	"github.com/lightpaw/male7/constants"
)

func NewHero(id int64, name string, initData *heroinit.HeroInitData, ctime time.Time) *Hero {
	return newHero(id, name, initData, nil, nil, ctime)
}

func UnmarshalHero(id int64, name string, initData *heroinit.HeroInitData, hsp *server_proto.HeroServerProto, datas *config.ConfigDatas, ctime time.Time) *Hero {
	return newHero(id, name, initData, hsp, datas, ctime)
}

func newHero(id int64, name string, initData *heroinit.HeroInitData, hsp *server_proto.HeroServerProto, datas *config.ConfigDatas, ctime time.Time) *Hero {

	buildingEffect := newBuildingEffect()

	maps := newHeroMaps()

	hero := &Hero{
		IdHolder:       idbytes.NewIdHolder(id),
		name:           name,
		level:          initData.FirstLevelHeroData,
		createTime:     ctime,
		buildingEffect: buildingEffect,
		misc:           NewHeroMisc(maps),
		countryMisc:    NewHeroCountryMisc(),
		domestic:       newDomestic(initData, buildingEffect, ctime),
		military:       newMilitary(id, initData, buildingEffect, ctime),
		depot:          newDepot(initData),
		tower:          newTower(),
		tasklist:       newTaskList(initData),
		//atomicData:       newAtomicData(initData),
		historyAmount:    NewHeroHistoryAmount(maps),
		openCombineEquip: newHeroOpenCombineEquip(),
		secretTower:      newHeroSecretTower(),
		randomEvent:      newHeroRandomEvent(ctime),
		hero_region:      newHeroRegion(initData, ctime),
		hero_resource: &hero_resource{
			unsafeResource:   &ResourceStorage{},
			safeResource:     &ResourceStorage{},
			nextResDecayTime: ctime,
		},
		hero_guild: newHeroGuild(initData, ctime),
		hero_shop:  &hero_shop{},
		dungeon:    newHeroDungeon(initData, ctime),
		tag:        newHeroTag(initData),
		reservation: &hero_reservation{
			goodsMap: make(map[uint64]uint64),
		},
		relation:          newRelation(),
		heroGen:           newHeroGen(),
		strategy:          newHeroStrategy(),
		buff:              NewHeroBuff(),
		clientDatas:       newClientDatas(maps),
		fishing:           NewHeroFishing(),
		function:          newHeroFunction(initData.FunctionOpenDataArray),
		zhengWu:           NewHeroZhengWu(initData.ZhengWuMiscData.NextAutoRefreshTime(ctime), initData.RandomZhengWu()),
		zhanJiang:         newHeroZhanJiang(),
		treasuryTree:      &hero_treasury_tree{},
		eventMap:          newHeroEventMap(maps),
		promotion:         newHeroPromotion(maps),
		activity:          newHeroActivity(),
		maps:              maps,
		bools:             newBools(),
		survey:            NewHeroSurvey(),
		settings:          NewHeroSettings(initData.DefaultSettings),
		lastChatTime:      ctime,
		lastRecoverSpTime: ctime,
		sp:                initData.FirstLevelHeroData.Sub.SpLimit,
	}

	buildingEffect.CalculateDomestic(hero)
	buildingEffect.CalculateSoldierStat(hero)

	hero.vip = NewHeroVip()

	hero.SetHead(initData.DefaultHead)
	hero.body = initData.DefaultBody

	hero.guildDonateTimes = make([]uint64, initData.GuildDonateTypeCount)

	hero.dailyShopGoodsMap = maps.getOrCreateMap(server_proto.HeroMapCategory_daily_shop_goods, true)

	hero.question = NewHeroQuestion()

	hero.hebi = NewHeroHebi()

	hero.mcBuild = NewHeroMcBuild()

	hero.redPacket = NewHeroRedPacket()

	hero.teach = NewHeroTeach()

	if hsp != nil && datas != nil {
		hero.unmarshal(hsp, datas, ctime)
	}

	if hero.Prosperity() > hero.ProsperityCapcity() {
		hero.SetProsperity(hero.ProsperityCapcity())
	}

	return hero
}

func (hero *Hero) unmarshal(proto *server_proto.HeroServerProto, datas *config.ConfigDatas, ctime time.Time) {
	if strings.HasPrefix(proto.Head, "http") {
		hero.SetHeadUrl(proto.Head)
	} else if data := datas.GetHeadData(proto.Head); data != nil {
		hero.SetHead(data)
	} else {
		hero.SetHead(datas.HeroInitData().DefaultHead)
	}

	if data := datas.GetBodyData(proto.Body); data != nil {
		hero.SetBody(data)
	}

	hero.setLevel(datas.HeroLevelData().Must(proto.Level))

	hero.exp = proto.Exp
	hero.male = proto.Male
	hero.location = proto.Location

	hero.oldName = proto.OldName
	hero.nextChangeNameTime = timeutil.Unix64(proto.NextChangeNameTime)
	hero.changeHeroNameTimes = proto.GetChangeHeroNameTimes()
	hero.giveFirstChangeHeroNamePrize = proto.GiveFirstChangeHeroNamePrize

	hero.collectedFirstJoinGuildPrize = proto.CollectedFirstJoinGuildPrize
	hero.collectedDailyGuildRankPrize = proto.CollectedDailyGuildRankPrize
	hero.workshopOutputTimes.SetStartTime(timeutil.Unix64(proto.WorkshopOutputStartTime))
	if len(proto.CollectedGuildTaskStages) > 0 {
		for taskId, stages := range proto.CollectedGuildTaskStages {
			m := make(map[int32]struct{})
			for _, stage := range stages.V {
				m[stage] = struct{}{}
			}
			hero.collectedTaskStages[taskId] = m
		}
	}
	hero.guildId = proto.GuildId
	hero.joinGuildTime = timeutil.Unix64(proto.JoinGuildTime)
	hero.guildContributionCoin = proto.ContributionCoin
	copy(hero.guildDonateTimes, proto.GuildDonateTimes)
	hero.joinGuildIds = i64.Copy(proto.JoinGuildIds)
	hero.beenInvateGuildIds = i64.Copy(proto.BeenInvateGuildIds)
	hero.nextNotifyGuildTime = timeutil.Unix64(proto.NextNotifyGuildTime)

	// 元宝
	hero.yuanbao = proto.Yuanbao
	// 点券
	hero.dianquan = proto.Dianquan
	hero.yinliang = proto.Yinliang
	// 体力值
	hero.sp = proto.Sp
	hero.buySpTimes = proto.BuySpTimes
	lastSpRecoverTime := timeutil.Unix64(proto.LastRecoverSpTime)
	if dlt := ctime.Sub(lastSpRecoverTime); dlt >= datas.MiscGenConfig().SpDuration {
		if hero.sp < hero.level.Sub.SpLimit {
			hero.sp += uint64(timeutil.DurationMarshal64(dlt) / timeutil.DurationMarshal64(datas.MiscGenConfig().SpDuration))
			if hero.sp > hero.level.Sub.SpLimit {
				hero.sp = hero.level.Sub.SpLimit
			}
		}
		dlt %= datas.MiscGenConfig().SpDuration
		hero.lastRecoverSpTime = timeutil.Unix64(timeutil.Marshal64(ctime) - timeutil.DurationMarshal64(datas.MiscGenConfig().SpDuration) + timeutil.DurationMarshal64(dlt))
	} else {
		hero.lastRecoverSpTime = lastSpRecoverTime
	}

	hero.createTime = timeutil.Unix64(proto.CreateTime)
	hero.totalOnlineTime = timeutil.Duration64(proto.TotalOnlineTime)
	hero.loginTime = timeutil.Unix64(proto.LoginTime)

	// 资源
	hero.unsafeResource.gold = proto.Gold
	hero.unsafeResource.food = proto.Food
	hero.unsafeResource.wood = proto.Wood
	hero.unsafeResource.stone = proto.Stone

	hero.safeResource.gold = proto.SafeGold
	hero.safeResource.food = proto.SafeFood
	hero.safeResource.wood = proto.SafeWood
	hero.safeResource.stone = proto.SafeStone

	hero.jade = proto.Jade
	hero.jadeOre = proto.JadeOre
	hero.historyJade = proto.HistoryJade
	hero.todayObtainJade = proto.TodayObtainJade
	hero.nextResDecayTime = timeutil.Unix64(proto.NextResDecayTime)

	hero.misc.Unmarshal(proto.Misc)
	hero.promotion.Unmarshal(proto.Promotion, datas, ctime)
	hero.activity.unmarshal(proto.Activity, ctime)

	hero.vip.unmarshal(proto.Vip, datas)

	hero.depot.unmarshal(hero.Id(), hero.name, proto.GetDepot(), datas, ctime)

	hero.domestic.unmarshal(hero.Id(), hero.name, proto.GetDomestic(), datas, ctime)

	// 武将unmarshal之前先解码数据
	hero.countryMisc.unmarshal(proto.CountryMisc, datas)

	// 武将unmarshal之前先解码数据，武将要用到称号数据
	taskProto := proto.GetTask()
	hero.tasklist.unmarshal(taskProto, datas)
	if taskProto != nil {
		hero.tasklist.unmarshalActivityListMode(taskProto.ActivityTaskListModeMap, datas, ctime)
	}

	// 武将unmarshal之前先解码数据
	hero.buff.unmarshal(proto.Buff, datas.BuffEffectData(), ctime)

	hero.military.unmarshal(hero, proto.GetMilitary(), datas, ctime)

	// 武将unmarshal之后再操作
	hero.hero_region.unmarshal(hero, proto.Region, datas, ctime)

	hero.tower.unmarshal(proto.GetTower())

	hero.secretTower.unmarshal(proto.GetSecretTower(), hero.Tower(), datas, ctime)

	hero.tag.unmarshal(proto.GetTag())

	hero.maps.unmarshal(proto)
	hero.bools.unmarshal(proto)

	hero.openCombineEquip.unmarshal(proto.GetOpenCombineEquip(), datas)

	hero.dungeon.unmarshal(proto.GetDungeon(), datas)

	if proto.GetLastOfflineTime() != 0 {
		hero.lastOfflineTime = timeutil.Unix64(proto.GetLastOfflineTime())
	}

	hero.reservation.unmarshal(proto.Reservation)

	hero.relation.unmarshal(proto.Relation, ctime)

	hero.randomEvent.unmarshal(proto.RandomEvent)

	hero.heroGen.Unmarshal(proto.HeroGen)
	if hero.CountryId() == 0 {
		// 老号初始化处理
		hero.SetCountryId(datas.CountryData().MinKeyData.Id)
	}

	hero.strategy.unmarshal(proto.Strategy)

	hero.clientDatas.unmarshal(proto.ClientDatas)

	hero.fishing.unmarshal(proto.Fishing)

	hero.function.unmarshal(proto.Function)

	hero.zhengWu.unmarshal(proto.ZhengWu, datas, ctime)

	hero.survey.unmarshal(proto.Survey, datas)

	hero.settings.unmarshal(proto.Settings)

	hero.zhanJiang.unmarshal(proto.ZhanJiang, datas, hero.Military().Captain)

	hero.treasuryTree.unmarshal(proto.TreasuryTree)

	hero.question.unmarshal(proto.Question)

	hero.hebi.unmarshal(proto.Hebi)

	hero.mcBuild.unmarshal(proto.McBuild)

	hero.redPacket.unmarshal(proto.RedPacket)

	hero.teach.unmarshal(proto.Teach)

	hero.hero_shop.unmarshal(proto.Shop, datas)

	for _, v := range proto.GuildEventPrizes {
		data := datas.GetGuildEventPrizeData(v.DataId)
		if data != nil {
			hero.AddGuildEventPrize(data, v.SendHeroId, timeutil.Unix64(v.ExpireTime), v.HideGiver)
		}
	}

}

type Hero struct {
	*hero_region
	*hero_resource
	*hero_guild
	*hero_shop
	reservation *hero_reservation
	relation    *hero_relation
	randomEvent *hero_random_event

	heroGen *hero_gen

	idbytes.IdHolder

	name    string
	headUrl string
	head    *head.HeadData
	body    *body.BodyData
	level   *herodata.HeroLevelData
	exp     uint64
	male    bool

	location uint64

	createTime time.Time

	// 玩家改名
	// 变更英雄名字的次数/时间
	oldName                      []string
	nextChangeNameTime           time.Time
	changeHeroNameTimes          uint64
	giveFirstChangeHeroNamePrize bool // 是否有给首次改名奖励

	lastChatTime      time.Time // 上次聊天的时间
	lastWorldChatTime time.Time // 上次聊天的时间

	// 元宝
	yuanbao uint64

	// 点券
	dianquan uint64

	// 银两
	yinliang uint64

	// 体力值
	sp                uint64
	buySpTimes        uint64 // 当天已经购买体力值的次数
	lastRecoverSpTime time.Time

	// 上次离线时间
	lastOnlineTime  time.Time
	lastOfflineTime time.Time
	loginTime       time.Time
	// 上次下线为止的总在线时长
	totalOnlineTime time.Duration

	vip *HeroVip

	misc *HeroMisc // 零碎的额外数据

	// 建筑效果数据
	buildingEffect *building_effect

	// 内政
	domestic *HeroDomestic

	// 战斗
	military *Military

	// 背包
	depot *Depot

	// 千重楼
	tower *Tower

	// 任务
	tasklist *TaskList

	//// atomic
	//atomicData *AtomicData

	tag *HeroTag

	historyAmount *HeroHistoryAmount // 历史记录值

	// 开启了的装备合成
	openCombineEquip *HeroOpenCombineEquip

	secretTower *HeroSecretTower

	dungeon *HeroDungeon

	strategy *HeroStrategy // 策略

	buff *HeroBuff // buff

	currentAdvantageCount int // 当前增益数量

	clientDatas *HeroClientDatas

	fishing *HeroFishing // 钓鱼

	function *HeroFunction // 功能开启

	zhengWu *HeroZhengWu // 政务

	question *HeroQuestion // 答题

	zhanJiang *HeroZhanJiang // 过关斩将

	survey *HeroSurvey // 问卷调查

	hebi *HeroHebi // 合璧

	mcBuild *HeroMcBuild // 名城营建

	redPacket *HeroRedPacket // 红包

	teach *HeroTeach // 教学

	settings *HeroSettings

	treasuryTree *hero_treasury_tree

	eventMap *guild_event_prize_map

	promotion *hero_promotion

	countryMisc *HeroCountryMisc

	activity *hero_activity // 活动

	maps  *hero_maps
	bools *hero_bools
}

func (hero *Hero) Pid() uint32 {
	return constants.PID
}

func (hero *Hero) Sid() uint32 {
	return heroid.GetSid(hero.Id())
}

func (hero *Hero) CountryMisc() *HeroCountryMisc {
	return hero.countryMisc
}

func (hero *Hero) Teach() *HeroTeach {
	return hero.teach
}

func (hero *Hero) HeroRedPacket() *HeroRedPacket {
	return hero.redPacket
}

func (hero *Hero) Vip() *HeroVip {
	return hero.vip
}

func (hero *Hero) VipLevel() uint64 {
	return hero.vip.level
}

func (hero *Hero) Name() string {
	return hero.name
}

func (hero *Hero) Head() string {
	if len(hero.headUrl) > 0 {
		return hero.headUrl
	}

	return hero.head.Id
}

func (hero *Hero) SetHeadUrl(headUrl string) {
	hero.headUrl = headUrl
	hero.head = nil
}

func (hero *Hero) SetHead(head *head.HeadData) {
	hero.headUrl = ""
	hero.head = head
}

func (hero *Hero) Body() uint64 {
	return hero.body.Id
}

func (hero *Hero) SetBody(body *body.BodyData) {
	hero.body = body
}

func (hero *Hero) Male() bool {
	return hero.male
}

func (hero *Hero) SetMale(male bool) {
	hero.male = male
}

func (hero *Hero) Location() uint64 {
	return hero.location
}

func (hero *Hero) SetLocation(location uint64) {
	hero.location = location
}

func (hero *Hero) SetCountryId(countryId uint64) {
	hero.Country().SetCountryId(countryId)
}

func (hero *Hero) CountryId() uint64 {
	return hero.Country().GetCountryId()
}

func (hero *Hero) Level() uint64 {
	return hero.level.Level
}

func (hero *Hero) setLevel(toSet *herodata.HeroLevelData) {
	if toSet != nil {
		hero.level = toSet
	}
}

func (hero *Hero) LevelData() *herodata.HeroLevelData {
	return hero.level
}

func (hero *Hero) IsMaxLevel() bool {
	return hero.level.NextLevel() == nil
}

func (hero *Hero) Exp() uint64 {
	return hero.exp
}

func (hero *Hero) AddExp(toAdd uint64) bool {
	if hero.level.NextLevel() == nil {
		return false
	}

	hero.exp += toAdd
	if hero.exp >= hero.level.Sub.UpgradeExp {
		// 升级
		hero.exp = u64.Sub(hero.exp, hero.level.Sub.UpgradeExp)
		hero.level = hero.level.NextLevel()

		if hero.level.NextLevel() == nil {
			// 最高级，把经验设置成0
			hero.exp = 0
		}

		return true
	}

	return false
}

const oldname_count = 10

func (hero *Hero) ChangeName(newName string) {
	hero.oldName = append(hero.oldName, hero.name)
	if len(hero.oldName) > oldname_count {
		for i := 0; i < oldname_count; i++ {
			hero.oldName[i] = hero.oldName[i+1]
		}

		hero.oldName = hero.oldName[:oldname_count]
	}

	hero.name = newName
}

func (hero *Hero) OldName() []string {
	return hero.oldName
}

func (hero *Hero) GetNextChangeNameTime() time.Time {
	return hero.nextChangeNameTime
}

func (hero *Hero) SetNextChangeNameTime(toSet time.Time) {
	hero.nextChangeNameTime = toSet
}

func (hero *Hero) GetChangeHeroNameTimes() uint64 {
	return hero.changeHeroNameTimes
}

func (hero *Hero) SetChangeHeroNameTimes(toSet uint64) {
	hero.changeHeroNameTimes = toSet
}

func (hero *Hero) HasGiveFirstChangeHeroNamePrize() bool {
	return hero.giveFirstChangeHeroNamePrize
}

func (hero *Hero) GiveFirstChangeHeroNamePrize() {
	hero.giveFirstChangeHeroNamePrize = true
}

func (hero *Hero) CreateTime() time.Time {
	return hero.createTime
}

func (hero *Hero) GetDailyMcResetTime() time.Time {
	return hero.MiscData().dailyMcResetTime
}

func (hero *Hero) SetDailyMcResetTime(toSet time.Time) {
	hero.MiscData().dailyMcResetTime = toSet
}

func (hero *Hero) GetWeeklyResetTime() time.Time {
	return hero.MiscData().weeklyResetTime
}

func (hero *Hero) SetWeeklyResetTime(toSet time.Time) {
	hero.MiscData().weeklyResetTime = toSet
}

func (hero *Hero) GetDailyZeroResetTime() time.Time {
	return hero.MiscData().dailyZeroResetTime
}

func (hero *Hero) SetDailyZeroResetTime(toSet time.Time) {
	hero.MiscData().dailyZeroResetTime = toSet
}

func (hero *Hero) GetDailyResetTime() time.Time {
	return hero.MiscData().dailyResetTime
}

func (hero *Hero) SetDailyResetTime(toSet time.Time) {
	hero.MiscData().dailyResetTime = toSet
}

func (hero *Hero) GetSeasonResetTime() time.Time {
	return hero.MiscData().seasonResetTime
}

func (hero *Hero) SetSeasonResetTime(toSet time.Time) {
	hero.MiscData().seasonResetTime = toSet
}

func (hero *Hero) HasEnoughYuanbao(amount uint64) bool {
	return hero.yuanbao >= amount
}

func (hero *Hero) GetYuanbao() uint64 {
	return hero.yuanbao
}

func (hero *Hero) ReduceYuanbao(toReduce uint64) {
	hero.yuanbao = u64.Sub(hero.yuanbao, toReduce)
}

func (hero *Hero) AddYuanbao(toAdd uint64) {
	hero.yuanbao += toAdd
}

func (hero *Hero) HasEnoughYinliang(amount uint64) bool {
	return hero.yinliang >= amount
}

func (hero *Hero) GetYinliang() uint64 {
	return hero.yinliang
}

func (hero *Hero) ReduceYinliang(toReduce uint64) {
	hero.yinliang = u64.Sub(hero.yinliang, toReduce)
}

func (hero *Hero) AddYinliang(toAdd uint64) {
	hero.yinliang += toAdd
}

func (hero *Hero) HasEnoughSp(sp uint64) bool {
	return hero.sp >= sp
}

func (hero *Hero) GetSp() uint64 {
	return hero.sp
}

// 用于自然恢复
func (hero *Hero) IncreaseOneSp() (changed bool) {
	if hero.sp >= hero.level.Sub.SpLimit {
		return
	}
	hero.sp++
	changed = true
	return
}

// 可能用于使用体力值道具，允许超出体力值上限
func (hero *Hero) AddSp(toAdd uint64) {
	hero.sp += toAdd
}

func (hero *Hero) ReduceSp(toReduce uint64) {
	hero.sp = u64.Sub(hero.sp, toReduce)
}

func (hero *Hero) GetBuySpTimes() uint64 {
	return hero.buySpTimes
}

func (hero *Hero) AddBuySpTimes(toAdd uint64) {
	hero.buySpTimes += toAdd
}

func (hero *Hero) LastRecoverSpTime() time.Time {
	return hero.lastRecoverSpTime
}

func (hero *Hero) SetLastRecoverSpTime(toSet time.Time) {
	hero.lastRecoverSpTime = toSet
}

func (hero *Hero) HasEnoughDianquan(amount uint64) bool {
	return hero.dianquan >= amount
}

func (hero *Hero) GetDianquan() uint64 {
	return hero.dianquan
}

func (hero *Hero) ReduceDianquan(toReduce uint64) {
	hero.dianquan = u64.Sub(hero.dianquan, toReduce)
}

func (hero *Hero) AddDianquan(toAdd uint64) {
	hero.dianquan += toAdd
}

func (hero *Hero) LastChatTime() time.Time {
	return hero.lastChatTime
}

func (hero *Hero) SetLastChatTime(toSet time.Time) {
	hero.lastChatTime = toSet
}

func (hero *Hero) LastWorldChatTime() time.Time {
	return hero.lastChatTime
}

func (hero *Hero) SetLastWorldChatTime(toSet time.Time) {
	hero.lastWorldChatTime = toSet
}

func (hero *Hero) LastOnlineTime() time.Time {
	return hero.lastOnlineTime
}

func (hero *Hero) LastOfflineTime() time.Time {
	return hero.lastOfflineTime
}

func (hero *Hero) Online(toSet time.Time) {
	hero.loginTime = toSet
	hero.lastOnlineTime = toSet
	hero.lastOfflineTime = time.Time{}
}

func (hero *Hero) UpdateOnlineTime(toSet time.Time) {
	hero.lastOnlineTime = toSet
}

func (hero *Hero) Offline(toSet time.Time) {
	hero.lastOfflineTime = toSet
	hero.totalOnlineTime += toSet.Sub(hero.loginTime)
}

func (hero *Hero) GetAccumLoginDay() uint64 {
	return hero.HistoryAmount().Amount(server_proto.HistoryAmountType_AccumLoginDays)
}

func (hero *Hero) BuildingEffect() *building_effect {
	return hero.buildingEffect
}

func (hero *Hero) Calculate(buildingType shared_proto.BuildingType) bool {
	return hero.buildingEffect.calculateBuildingType(hero, buildingType)
}

func (hero *Hero) Misc() *HeroMisc {
	return hero.misc
}

func (hero *Hero) Domestic() *HeroDomestic {
	return hero.domestic
}

func (hero *Hero) Military() *Military {
	return hero.military
}

func (hero *Hero) Depot() *Depot {
	return hero.depot
}

func (hero *Hero) Tower() *Tower {
	return hero.tower
}

func (hero *Hero) HistoryAmount() *HeroHistoryAmount {
	return hero.historyAmount
}

func (hero *Hero) TaskList() *TaskList {
	return hero.tasklist
}

//func (hero *Hero) AtomicData() *AtomicData {
//	return hero.atomicData
//}

func (hero *Hero) OpenCombineEquip() *HeroOpenCombineEquip {
	return hero.openCombineEquip
}

func (hero *Hero) SecretTower() *HeroSecretTower {
	return hero.secretTower
}

func (hero *Hero) Tag() *HeroTag {
	return hero.tag
}

func (hero *Hero) Strategy() *HeroStrategy {
	return hero.strategy
}

func (hero *Hero) Buff() *HeroBuff {
	return hero.buff
}

func (hero *Hero) CurrentAdvantageCount() int {
	return hero.currentAdvantageCount
}

func (hero *Hero) SetCurrentAdvantageCount(c int) {
	hero.currentAdvantageCount = c
}

func (hero *Hero) Dungeon() *HeroDungeon {
	return hero.dungeon
}

func (hero *Hero) Reservation() *hero_reservation {
	return hero.reservation
}

func (hero *Hero) Relation() *hero_relation {
	return hero.relation
}

func (hero *Hero) RandomEvent() *hero_random_event {
	return hero.randomEvent
}

func (hero *Hero) ClientDatas() *HeroClientDatas {
	return hero.clientDatas
}

func (hero *Hero) Fishing() *HeroFishing {
	return hero.fishing
}

func (hero *Hero) Function() *HeroFunction {
	return hero.function
}

func (hero *Hero) ZhengWu() *HeroZhengWu {
	return hero.zhengWu
}

func (hero *Hero) Survey() *HeroSurvey {
	return hero.survey
}

func (hero *Hero) Settings() *HeroSettings {
	return hero.settings
}

func (hero *Hero) ZhanJiang() *HeroZhanJiang {
	return hero.zhanJiang
}

func (hero *Hero) Question() *HeroQuestion {
	return hero.question
}

func (hero *Hero) Hebi() *HeroHebi {
	return hero.hebi
}

func (hero *Hero) McBuild() *HeroMcBuild {
	return hero.mcBuild
}

func (hero *Hero) FriendsCount() uint64 {
	return hero.relation.friendCount
}

func (hero *Hero) TotalOnlineTime() time.Duration {
	return hero.totalOnlineTime
}

//func (hero *Hero) Id() int64 {
//	return hero.IdHolder.Id()
//}

func (hero *Hero) GetAllRes(t shared_proto.ResType) (result uint64) {
	if hero.unsafeResource != nil {
		result += hero.unsafeResource.GetRes(t)
	}
	if hero.safeResource != nil {
		result += hero.safeResource.GetRes(t)
	}
	return
}

func (hero *Hero) EncodeServer() *server_proto.HeroServerProto {

	proto := &server_proto.HeroServerProto{}
	proto.Id = hero.Id()
	proto.Name = hero.name
	proto.Head = hero.Head()
	proto.Body = hero.Body()
	proto.Level = hero.Level()
	proto.Exp = hero.exp
	proto.Male = hero.male
	proto.Location = hero.location
	proto.CreateTime = timeutil.Marshal64(hero.createTime)
	proto.TotalOnlineTime = timeutil.DurationMarshal64(hero.totalOnlineTime)
	proto.LoginTime = timeutil.Marshal64(hero.loginTime)

	proto.OldName = hero.oldName
	proto.NextChangeNameTime = timeutil.Marshal64(hero.nextChangeNameTime)
	proto.ChangeHeroNameTimes = hero.changeHeroNameTimes
	proto.GiveFirstChangeHeroNamePrize = hero.HasGiveFirstChangeHeroNamePrize()

	proto.CollectedFirstJoinGuildPrize = hero.collectedFirstJoinGuildPrize
	proto.CollectedDailyGuildRankPrize = hero.collectedDailyGuildRankPrize
	proto.WorkshopOutputStartTime = hero.workshopOutputTimes.StartTimeUnix64()
	if len(hero.collectedTaskStages) > 0 {
		proto.CollectedGuildTaskStages = make(map[uint64]*shared_proto.Int32ArrayProto)
		for taskId, stages := range hero.collectedTaskStages {
			arr := make([]int32, 0, len(stages))
			for stage, _ := range stages {
				arr = append(arr, stage)
			}
			proto.CollectedGuildTaskStages[taskId] = &shared_proto.Int32ArrayProto{
				V: arr,
			}
		}
	}
	proto.GuildId = hero.guildId
	proto.JoinGuildTime = timeutil.Marshal64(hero.joinGuildTime)
	proto.ContributionCoin = hero.guildContributionCoin
	proto.GuildDonateTimes = hero.guildDonateTimes
	proto.JoinGuildIds = hero.joinGuildIds
	proto.BeenInvateGuildIds = hero.beenInvateGuildIds
	proto.NextNotifyGuildTime = timeutil.Marshal64(hero.nextNotifyGuildTime)

	proto.Region = hero.encodeRegion()

	// 元宝
	proto.Yuanbao = hero.yuanbao
	// 点券
	proto.Dianquan = hero.dianquan
	proto.Yinliang = hero.yinliang
	// 体力值
	proto.Sp = hero.sp
	proto.LastRecoverSpTime = timeutil.Marshal64(hero.lastRecoverSpTime)
	proto.BuySpTimes = hero.buySpTimes

	proto.Gold = hero.unsafeResource.gold
	proto.Food = hero.unsafeResource.food
	proto.Wood = hero.unsafeResource.wood
	proto.Stone = hero.unsafeResource.stone

	proto.SafeGold = hero.safeResource.gold
	proto.SafeFood = hero.safeResource.food
	proto.SafeWood = hero.safeResource.wood
	proto.SafeStone = hero.safeResource.stone

	proto.Jade = hero.jade
	proto.JadeOre = hero.jadeOre
	proto.HistoryJade = hero.historyJade
	proto.TodayObtainJade = hero.todayObtainJade
	proto.NextResDecayTime = timeutil.Marshal64(hero.nextResDecayTime)

	proto.Misc = hero.misc.EncodeServer()
	proto.Promotion = hero.promotion.EncodeServer()
	proto.Activity = hero.activity.encodeServer()
	proto.Domestic = hero.domestic.encodeServer()
	proto.Military = hero.military.encodeServer()
	proto.Depot = hero.depot.encodeServer()
	proto.Tower = hero.tower.encodeServer()
	proto.SecretTower = hero.SecretTower().encodeServer()
	proto.Task = hero.tasklist.encodeServer()

	proto.OpenCombineEquip = hero.OpenCombineEquip().encodeServer()

	proto.HeroMaps, proto.HeroKeys = hero.maps.encode()
	proto.HeroBools = hero.bools.serverBools

	proto.Tag = hero.Tag().Encode()

	proto.Dungeon = hero.Dungeon().EncodeServer()

	proto.LastOfflineTime = timeutil.Marshal64(hero.lastOfflineTime)

	proto.Reservation = hero.reservation.encode()

	proto.Relation = hero.relation.encode()

	proto.RandomEvent = hero.randomEvent.encode()

	proto.HeroGen = hero.heroGen.EncodeServer()

	proto.Strategy = hero.strategy.encode()

	proto.Buff = hero.buff.encodeServer()

	proto.ClientDatas = hero.clientDatas.encodeServer()

	proto.Fishing = hero.fishing.encode()

	proto.Function = hero.function.encodeServer()

	proto.ZhengWu = hero.zhengWu.EncodeServer()

	proto.Survey = hero.survey.Encode()

	proto.Settings = hero.settings.Encode()

	proto.ZhanJiang = hero.zhanJiang.EncodeServer()

	proto.TreasuryTree = hero.treasuryTree.encode()

	proto.Question = hero.question.encodeServer()

	proto.Hebi = hero.hebi.encode()

	proto.Shop = hero.hero_shop.encode()

	proto.Vip = hero.vip.encodeServer()

	proto.McBuild = hero.mcBuild.encodeServer()

	proto.RedPacket = hero.redPacket.encodeServer()

	proto.Teach = hero.teach.encodeServer()

	proto.CountryMisc = hero.countryMisc.encodeServer()

	for _, v := range hero.eventMap.guildEventPrizes {
		proto.GuildEventPrizes = append(proto.GuildEventPrizes, &server_proto.HeroGuildEventPrizeServerProto{
			DataId:     v.Data.Id,
			SendHeroId: v.SendHeroId,
			ExpireTime: timeutil.Marshal64(v.ExpireTime),
			HideGiver:  v.HideGiver,
		})
	}

	return proto
}

func (hero *Hero) Marshal() ([]byte, error) {
	return hero.EncodeServer().Marshal()
}

func (hero *Hero) ResetSeason(resetTime time.Time, seasonData *season.SeasonData) {
	hero.MiscData().seasonResetTime = resetTime

	// 添加额外次数
	if seasonData.AddMultiMonsterTimes > 0 {
		// 在重置的时间点，加额外次数
		hero.multiLevelNpcTimes.AddTimes(seasonData.AddMultiMonsterTimes, resetTime, seasonData.AddMultiMonsterTimes)
	}

	hero.domestic.ResetSeason()
}

func (hero *Hero) ResetWeekly(resetTime time.Time) {
	hero.SetWeeklyResetTime(resetTime)
	hero.guildResetWeekly()
}

func (hero *Hero) ResetDailyZero(resetTime time.Time, datas interface {
	GetGoodsData(uint64 uint64) *goods.GoodsData
}) {
	logrus.Debugf("hero:%v %v 0点重置", hero.Id(), hero.Level())
	hero.MiscData().dailyZeroResetTime = resetTime
}

func (hero *Hero) ResetDailyMc(resetTime time.Time, datas interface {
	GetGoodsData(uint64 uint64) *goods.GoodsData
}) {
	logrus.Debugf("hero:%v %v 名城重置", hero.Id(), hero.Level())
	hero.MiscData().dailyMcResetTime = resetTime
	hero.mcBuild.ResetDaily()
}

func (hero *Hero) ResetDaily(resetTime time.Time, datas interface {
	GetGoodsData(uint64 uint64) *goods.GoodsData
}) {
	hero.MiscData().dailyResetTime = resetTime

	hero.tower.resetDaily()

	hero.SecretTower().ResetDaily()

	hero.Domestic().ResetDaily()
	hero.Military().resetDaily()

	hero.guildResetDaily()

	hero.Fishing().ResetDaily()

	hero.todayObtainJade = 0

	hero.maps.resetDaily()

	hero.treasuryTree.resetDaily()

	hero.eventMap.resetDaily(resetTime)

	hero.zhengWu.ResetDaily()

	hero.question.ResetDaily()

	hero.zhanJiang.ResetDaily()

	hero.Hebi().ResetDaily()

	hero.depot.resetDaily(datas)

	hero.FarmExtra().Reset()

	hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumLoginDays)

	hero.heroGen.resetDaily()

	// 最后再重置活跃度任务，否则会因为别的数据没有重置，导致任务直接完成了
	hero.TaskList().ResetDaily(hero)

	hero.dungeon.ResetDaily()

	hero.buySpTimes = 0 // 重置购买体力值次数

	hero.strategy.ResetDaily()
}

func (hero *Hero) EncodeOther(ctime time.Time, getGuildSnapshot guildsnapshotdata.Getter, baiZhanGetter interface {
	GetHistoryMaxJunXianLevel(int64) uint64
	GetJunXianLevel(int64) uint64
}) *shared_proto.OtherHeroProto {
	heroProto := &shared_proto.OtherHeroProto{}

	heroProto.Id = hero.IdBytes()
	heroProto.Name = hero.Name()
	heroProto.Head = hero.Head()
	heroProto.Body = u64.Int32(hero.Body())
	heroProto.Level = u64.Int32(hero.Level())
	heroProto.Male = hero.Male()
	heroProto.Location = u64.Int32(hero.Location())
	heroProto.HasOldName = len(hero.OldName()) > 0
	heroProto.CountryId = u64.Int32(hero.CountryId())

	for _, t := range hero.Troops() {
		if t != nil {
			heroProto.FightAmount = i32.Max(heroProto.FightAmount, u64.Int32(t.FullFightAmount()))
		}
	}

	if hero.guildId != 0 {
		heroProto.GuildId = i64.Int32(hero.guildId)
		g := getGuildSnapshot(hero.guildId)
		if g != nil {
			heroProto.GuildName = g.Name
			heroProto.GuildFlagName = g.FlagName
		}
	}

	heroProto.MaxTowerFloor = u64.Int32(hero.Tower().HistoryMaxFloor())

	heroProto.JunXianLevel = u64.Int32(baiZhanGetter.GetJunXianLevel(hero.Id()))
	heroProto.HistoryMaxJunXianLevel = u64.Int32(baiZhanGetter.GetHistoryMaxJunXianLevel(hero.Id()))

	heroProto.Home = hero.home.encodeClient()

	heroProto.HasWhiteFlag = hero.home.whiteFlagGuildId != 0 && ctime.Before(hero.home.whiteFlagDisappearTime)
	if heroProto.HasWhiteFlag {
		g := getGuildSnapshot(hero.home.whiteFlagGuildId)
		if g != nil {
			heroProto.WhiteFlagName = g.FlagName
		} else {
			heroProto.HasWhiteFlag = false
		}
	}

	heroProto.MaxBaseLevel = u64.Int32(hero.home.historyMaxBaseLevel)

	heroProto.Domestic = hero.EncodeOtherDomestic(hero.Domestic())

	heroProto.Tag = hero.Tag().Encode()

	heroProto.SelectShowAchieves = hero.TaskList().AchieveTaskList().EncodeSelectShowAchieves()

	return heroProto
}

func (hero *Hero) EncodeClient(ctime time.Time, getGuildSnapshot guildsnapshotdata.Getter, configDatas config.Configs, serverInfo *shared_proto.HeroServerInfoProto) *shared_proto.HeroProto {

	heroProto := &shared_proto.HeroProto{}
	heroProto.Ctime = timeutil.Marshal32(ctime)
	heroProto.ServerInfo = serverInfo

	heroProto.Id = hero.IdBytes()
	heroProto.Name = hero.Name()
	heroProto.Head = 1 // 设置个默认值
	heroProto.Head2 = hero.Head()
	heroProto.Body = u64.Int32(hero.Body())
	heroProto.Level = u64.Int32(hero.Level())
	heroProto.Exp = u64.Int32(hero.Exp())
	heroProto.Male = hero.Male()
	heroProto.Location = u64.Int32(hero.location)
	heroProto.HasOldName = len(hero.OldName()) > 0
	heroProto.NextChangeNameTime = timeutil.Marshal32(hero.GetNextChangeNameTime())
	heroProto.ChangeHeroNameTimes = u64.Int32(hero.GetChangeHeroNameTimes())
	heroProto.GiveFirstChangeHeroNamePrize = hero.HasGiveFirstChangeHeroNamePrize()

	heroProto.FightAmount = u64.Int32(hero.GetHomeDefenserFightAmount())
	heroProto.CreateTime = timeutil.Marshal32(hero.CreateTime())

	heroProto.CollectedFirstJoinGuildPrize = hero.collectedFirstJoinGuildPrize
	heroProto.CollectedDailyGuildRankPrize = hero.collectedDailyGuildRankPrize
	heroProto.WorkshopOutputStartTime = hero.workshopOutputTimes.StartTimeUnix32()
	if l := len(hero.collectedTaskStages); l > 0 {
		heroProto.CollectedGuildTaskStages = make([]*shared_proto.Int32PairInt32ArrayProto, 0, l)
		for taskId, stages := range hero.collectedTaskStages {
			arr := make([]int32, 0, len(stages))
			for stage, _ := range stages {
				arr = append(arr, stage)
			}
			heroProto.CollectedGuildTaskStages = append(heroProto.CollectedGuildTaskStages, &shared_proto.Int32PairInt32ArrayProto{
				K: u64.Int32(taskId),
				V: arr,
			})
		}
	}
	if hero.guildId != 0 {
		heroProto.GuildId = i64.Int32(hero.guildId)
		g := getGuildSnapshot(hero.guildId)
		if g != nil {
			// 待删除
			heroProto.GuildName = g.Name
			heroProto.GuildFlagName = g.FlagName

			heroProto.Guild = g.HeroGuildProto()
		}
		heroProto.JoinGuildTime = timeutil.Marshal32(hero.joinGuildTime)
	}
	heroProto.ContributionCoin = u64.Int32(hero.guildContributionCoin)
	heroProto.GuildDonateTimes = u64.Int32Array(hero.guildDonateTimes)
	heroProto.JoinGuildIds = i64.Int32Array(hero.joinGuildIds)
	heroProto.BeenInvateGuildIds = i64.Int32Array(hero.beenInvateGuildIds)
	heroProto.NextNotifyGuildTime = timeutil.Marshal32(hero.nextNotifyGuildTime)

	heroProto.Region = hero.encodeRegionClient(ctime, getGuildSnapshot)

	//// 自己主界面的军情
	//for _, t := range hero.troopsMap {
	//	heroProto.SelfMilitaryInfo = append(heroProto.SelfMilitaryInfo, t.encodeMilitaryInfo())
	//}

	heroProto.Buff = hero.Buff().encode()

	// 元宝
	heroProto.Yuanbao = u64.Int32(hero.GetYuanbao())
	// 点券
	heroProto.Dianquan = u64.Int32(hero.GetDianquan())

	heroProto.Yinliang = u64.Int32(hero.yinliang)

	heroProto.Sp = u64.Int32(hero.sp)
	heroProto.BuySpTimes = u64.Int32(hero.buySpTimes)

	heroProto.Domestic = hero.encodeDomestic(hero.Domestic(), configDatas, ctime)

	heroProto.Military = hero.Military().encodeMilitary(ctime)

	heroProto.Charge = hero.misc.EncodeClient()

	heroProto.Depot = hero.Depot().EncodeClient()

	heroProto.Tower = hero.Tower().EncodeClient()

	heroProto.SecretTower = hero.SecretTower().EncodeClient()

	heroProto.Task = hero.TaskList().EncodeClient()

	heroProto.OpenCombineEquip = hero.OpenCombineEquip().EncodeClient()

	heroProto.Shop = hero.encodeShopClient()

	heroProto.Tag = hero.Tag().Encode()

	heroProto.Dungeon = hero.Dungeon().EncodeClient(configDatas)

	heroProto.Strategy = hero.strategy.encodeClient()

	heroProto.ClientDatas = hero.clientDatas.encodeClient()

	heroProto.Fishing = hero.fishing.encode()

	heroProto.FuncOpen = hero.function.encodeClient()

	heroProto.ZhengWu = hero.zhengWu.EncodeClient()

	heroProto.Survey = hero.survey.Encode()

	heroProto.Settings = hero.settings.Encode()

	heroProto.ZhanJiang = hero.zhanJiang.EncodeClient()

	heroProto.TreasuryTree = hero.treasuryTree.encodeClient()

	heroProto.Question = hero.question.encodeClient(hero.HistoryAmount().Amount(server_proto.HistoryAmountType_QuestionRightAmount))

	heroProto.Hebi = hero.hebi.encode()

	heroProto.Relation = hero.relation.encodeClient()

	heroProto.Vip = hero.vip.Encode()

	heroProto.McBuild = hero.mcBuild.Encode()

	heroProto.Teach = hero.teach.encode()

	heroProto.RedPacket = hero.redPacket.encode()

	heroProto.CountryOfficial = hero.countryMisc.encode()

	heroProto.RandomEvent = hero.randomEvent.encodeClient()

	heroProto.HeroGen = hero.heroGen.EncodeClient()

	heroProto.Bools = hero.bools.encodeClient()
	heroProto.BoolsArray = hero.bools.encodeArrayClient()

	heroProto.AccumLoginDay = u64.Int32(hero.GetAccumLoginDay())
	heroProto.Promotion = hero.promotion.encodeClient()

	return heroProto
}
