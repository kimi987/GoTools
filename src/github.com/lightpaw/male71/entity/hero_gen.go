package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

var (
	_ = i64.Int32
	_ = u64.Int32
	_ = shared_proto.ErrIntOverflowHeroGen
	_ = server_proto.ErrIntOverflowHeroGen
	_ = time.Second
	_ = timeutil.GameZone
)

func newHeroGenXuanyuan() *HeroGenXuanyuan {
	return &HeroGenXuanyuan{}
}

// 轩辕会武
type HeroGenXuanyuan struct {
	challengeTarget []int64 // 挑战目标列表

	score uint64 // 挑战积分

	lostScore uint64 // 今日已损失积分

	win uint64 // 获胜场次

	lose uint64 // 失败场次

	rankPrizeCollected bool // 今日已领取排名奖励

	lastResetTime int64 // 上一次重置时间

}

func (hero *Hero) Xuanyuan() *HeroGenXuanyuan {
	return hero.heroGen.xuanyuan
}

func (d *HeroGenXuanyuan) AddChallengeTarget(toAdd int64) {
	d.challengeTarget = append(d.challengeTarget, toAdd)
}

func (d *HeroGenXuanyuan) GetChallengeTargetLen() int {
	return len(d.challengeTarget)
}

func (d *HeroGenXuanyuan) ContainsChallengeTarget(c int64) bool {
	return i64.Contains(d.challengeTarget, c)
}

func (d *HeroGenXuanyuan) GetScore() uint64 {
	return d.score
}

func (d *HeroGenXuanyuan) AddScore(toAdd uint64) uint64 {
	d.score += toAdd
	return d.score
}

func (d *HeroGenXuanyuan) SubScore(toSub uint64) uint64 {
	d.score = u64.Sub(d.score, toSub)
	return d.score
}

func (d *HeroGenXuanyuan) GetLostScore() uint64 {
	return d.lostScore
}

func (d *HeroGenXuanyuan) AddLostScore(toAdd uint64) uint64 {
	d.lostScore += toAdd
	return d.lostScore
}

func (d *HeroGenXuanyuan) GetWin() uint64 {
	return d.win
}

func (d *HeroGenXuanyuan) IncWin() uint64 {
	d.win++
	return d.win
}

func (d *HeroGenXuanyuan) GetLose() uint64 {
	return d.lose
}

func (d *HeroGenXuanyuan) IncLose() uint64 {
	d.lose++
	return d.lose
}

func (d *HeroGenXuanyuan) GetRankPrizeCollected() bool {
	return d.rankPrizeCollected
}

func (d *HeroGenXuanyuan) SetRankPrizeCollected(toSet bool) bool {
	o := d.rankPrizeCollected
	d.rankPrizeCollected = toSet
	return o
}

func (d *HeroGenXuanyuan) Reset() {
	d.challengeTarget = nil
	d.lostScore = 0
	d.rankPrizeCollected = false
}

func (d *HeroGenXuanyuan) EncodeClient() *shared_proto.HeroGenXuanyuanProto {
	proto := &shared_proto.HeroGenXuanyuanProto{}

	for _, v := range d.challengeTarget {
		proto.ChallengeTarget = append(proto.ChallengeTarget, i64.ToBytes(v))
	}
	proto.RankPrizeCollected = d.rankPrizeCollected

	return proto
}

func (d *HeroGenXuanyuan) EncodeServer() *server_proto.HeroGenXuanyuanServerProto {
	proto := &server_proto.HeroGenXuanyuanServerProto{}
	proto.ChallengeTarget = d.challengeTarget
	proto.Score = d.score
	proto.LostScore = d.lostScore
	proto.Win = d.win
	proto.Lose = d.lose
	proto.RankPrizeCollected = d.rankPrizeCollected
	proto.LastResetTime = d.lastResetTime

	return proto
}

func (d *HeroGenXuanyuan) Unmarshal(proto *server_proto.HeroGenXuanyuanServerProto) {
	if proto == nil {
		return
	}
	d.challengeTarget = proto.ChallengeTarget
	d.score = proto.Score
	d.lostScore = proto.LostScore
	d.win = proto.Win
	d.lose = proto.Lose
	d.rankPrizeCollected = proto.RankPrizeCollected
	d.lastResetTime = proto.LastResetTime
}

func newHeroGenMiscData() *HeroGenMiscData {
	return &HeroGenMiscData{}
}

// 英雄杂项
type HeroGenMiscData struct {
	weeklyResetTime time.Time // 每周重置时间（上一次）

	dailyResetTime time.Time // 每日重置时间（上一次）

	seasonResetTime time.Time // 四季重置时间（上一次）

	defenserNextFullSoldierTime time.Time // 主城防守部队，下次自动补兵时间

	defenserDontAutoFullSoldier bool // 主城防守部队，是否不要自动补兵（false表示自动补兵）

	isRestoreProsperity bool // 是否恢复繁荣度

	nextCollectTaxTime time.Time // 下次收税时间

	fishPoint uint64 // 钓鱼积分

	fishCombo uint64 // 钓鱼连击（高级）

	fishPointExchangeIndex int32 // 钓鱼兑换映射的消耗数组下标

	banChatEndTime time.Time // 禁言结束时间

	banLoginEndTime time.Time // 封号结束时间

	dailyZeroResetTime time.Time // 每日零点重置时间（上一次）

	dailyMcResetTime time.Time // 每日22点重置时间（上一次）

}

func (hero *Hero) MiscData() *HeroGenMiscData {
	return hero.heroGen.misc_data
}

func (d *HeroGenMiscData) GetDefenserNextFullSoldierTime() time.Time {
	return d.defenserNextFullSoldierTime
}

func (d *HeroGenMiscData) SetDefenserNextFullSoldierTime(toSet time.Time) time.Time {
	o := d.defenserNextFullSoldierTime
	d.defenserNextFullSoldierTime = toSet
	return o
}

func (d *HeroGenMiscData) GetDefenserDontAutoFullSoldier() bool {
	return d.defenserDontAutoFullSoldier
}

func (d *HeroGenMiscData) SetDefenserDontAutoFullSoldier(toSet bool) bool {
	o := d.defenserDontAutoFullSoldier
	d.defenserDontAutoFullSoldier = toSet
	return o
}

func (d *HeroGenMiscData) GetIsRestoreProsperity() bool {
	return d.isRestoreProsperity
}

func (d *HeroGenMiscData) SetIsRestoreProsperity(toSet bool) bool {
	o := d.isRestoreProsperity
	d.isRestoreProsperity = toSet
	return o
}

func (d *HeroGenMiscData) GetNextCollectTaxTime() time.Time {
	return d.nextCollectTaxTime
}

func (d *HeroGenMiscData) SetNextCollectTaxTime(toSet time.Time) time.Time {
	o := d.nextCollectTaxTime
	d.nextCollectTaxTime = toSet
	return o
}

func (d *HeroGenMiscData) GetFishPoint() uint64 {
	return d.fishPoint
}

func (d *HeroGenMiscData) SetFishPoint(toSet uint64) uint64 {
	o := d.fishPoint
	d.fishPoint = toSet
	return o
}

func (d *HeroGenMiscData) GetFishCombo() uint64 {
	return d.fishCombo
}

func (d *HeroGenMiscData) SetFishCombo(toSet uint64) uint64 {
	o := d.fishCombo
	d.fishCombo = toSet
	return o
}

func (d *HeroGenMiscData) IncFishCombo() uint64 {
	d.fishCombo++
	return d.fishCombo
}

func (d *HeroGenMiscData) GetFishPointExchangeIndex() int32 {
	return d.fishPointExchangeIndex
}

func (d *HeroGenMiscData) IncFishPointExchangeIndex() int32 {
	d.fishPointExchangeIndex++
	return d.fishPointExchangeIndex
}

func (d *HeroGenMiscData) GetBanChatEndTime() time.Time {
	return d.banChatEndTime
}

func (d *HeroGenMiscData) SetBanChatEndTime(toSet time.Time) time.Time {
	o := d.banChatEndTime
	d.banChatEndTime = toSet
	return o
}

func (d *HeroGenMiscData) GetBanLoginEndTime() time.Time {
	return d.banLoginEndTime
}

func (d *HeroGenMiscData) SetBanLoginEndTime(toSet time.Time) time.Time {
	o := d.banLoginEndTime
	d.banLoginEndTime = toSet
	return o
}

func (d *HeroGenMiscData) Reset() {
}

func (d *HeroGenMiscData) EncodeClient() *shared_proto.HeroGenMiscDataProto {
	proto := &shared_proto.HeroGenMiscDataProto{}
	proto.WeeklyResetTime = timeutil.Marshal32(d.weeklyResetTime)
	proto.DailyResetTime = timeutil.Marshal32(d.dailyResetTime)
	proto.SeasonResetTime = timeutil.Marshal32(d.seasonResetTime)
	proto.DefenserNextFullSoldierTime = timeutil.Marshal32(d.defenserNextFullSoldierTime)
	proto.DefenserDontAutoFullSoldier = d.defenserDontAutoFullSoldier
	proto.FishPoint = u64.Int32(d.fishPoint)
	proto.FishCombo = u64.Int32(d.fishCombo)
	proto.FishPointExchangeIndex = d.fishPointExchangeIndex
	proto.BanChatEndTime = timeutil.Marshal32(d.banChatEndTime)
	proto.DailyZeroResetTime = timeutil.Marshal32(d.dailyZeroResetTime)
	proto.DailyMcResetTime = timeutil.Marshal32(d.dailyMcResetTime)

	return proto
}

func (d *HeroGenMiscData) EncodeServer() *server_proto.HeroGenMiscDataServerProto {
	proto := &server_proto.HeroGenMiscDataServerProto{}
	proto.WeeklyResetTime = timeutil.Marshal64(d.weeklyResetTime)
	proto.DailyResetTime = timeutil.Marshal64(d.dailyResetTime)
	proto.SeasonResetTime = timeutil.Marshal64(d.seasonResetTime)
	proto.DefenserNextFullSoldierTime = timeutil.Marshal64(d.defenserNextFullSoldierTime)
	proto.DefenserDontAutoFullSoldier = d.defenserDontAutoFullSoldier
	proto.IsRestoreProsperity = d.isRestoreProsperity
	proto.NextCollectTaxTime = timeutil.Marshal64(d.nextCollectTaxTime)
	proto.FishPoint = d.fishPoint
	proto.FishCombo = d.fishCombo
	proto.FishPointExchangeIndex = d.fishPointExchangeIndex
	proto.BanChatEndTime = timeutil.Marshal64(d.banChatEndTime)
	proto.BanLoginEndTime = timeutil.Marshal64(d.banLoginEndTime)
	proto.DailyZeroResetTime = timeutil.Marshal64(d.dailyZeroResetTime)
	proto.DailyMcResetTime = timeutil.Marshal64(d.dailyMcResetTime)

	return proto
}

func (d *HeroGenMiscData) Unmarshal(proto *server_proto.HeroGenMiscDataServerProto) {
	if proto == nil {
		return
	}
	d.weeklyResetTime = timeutil.Unix64(proto.WeeklyResetTime)
	d.dailyResetTime = timeutil.Unix64(proto.DailyResetTime)
	d.seasonResetTime = timeutil.Unix64(proto.SeasonResetTime)
	d.defenserNextFullSoldierTime = timeutil.Unix64(proto.DefenserNextFullSoldierTime)
	d.defenserDontAutoFullSoldier = proto.DefenserDontAutoFullSoldier
	d.isRestoreProsperity = proto.IsRestoreProsperity
	d.nextCollectTaxTime = timeutil.Unix64(proto.NextCollectTaxTime)
	d.fishPoint = proto.FishPoint
	d.fishCombo = proto.FishCombo
	d.fishPointExchangeIndex = proto.FishPointExchangeIndex
	d.banChatEndTime = timeutil.Unix64(proto.BanChatEndTime)
	d.banLoginEndTime = timeutil.Unix64(proto.BanLoginEndTime)
	d.dailyZeroResetTime = timeutil.Unix64(proto.DailyZeroResetTime)
	d.dailyMcResetTime = timeutil.Unix64(proto.DailyMcResetTime)
}

func newHeroGenFarmExtra() *HeroGenFarmExtra {
	return &HeroGenFarmExtra{}
}

// 英雄农场相关数据
type HeroGenFarmExtra struct {
	dailyStealGold int32 // 每日偷铜币总量

	dailyStealStone int32 // 每日偷石料总量

}

func (hero *Hero) FarmExtra() *HeroGenFarmExtra {
	return hero.heroGen.farm_extra
}

func (d *HeroGenFarmExtra) GetDailyStealGold() int32 {
	return d.dailyStealGold
}

func (d *HeroGenFarmExtra) SetDailyStealGold(toSet int32) int32 {
	o := d.dailyStealGold
	d.dailyStealGold = toSet
	return o
}

func (d *HeroGenFarmExtra) GetDailyStealStone() int32 {
	return d.dailyStealStone
}

func (d *HeroGenFarmExtra) SetDailyStealStone(toSet int32) int32 {
	o := d.dailyStealStone
	d.dailyStealStone = toSet
	return o
}

func (d *HeroGenFarmExtra) Reset() {
	d.dailyStealGold = 0
	d.dailyStealStone = 0
}

func (d *HeroGenFarmExtra) EncodeClient() *shared_proto.HeroGenFarmExtraProto {
	proto := &shared_proto.HeroGenFarmExtraProto{}
	proto.DailyStealGold = d.dailyStealGold
	proto.DailyStealStone = d.dailyStealStone

	return proto
}

func (d *HeroGenFarmExtra) EncodeServer() *server_proto.HeroGenFarmExtraServerProto {
	proto := &server_proto.HeroGenFarmExtraServerProto{}
	proto.DailyStealGold = d.dailyStealGold
	proto.DailyStealStone = d.dailyStealStone

	return proto
}

func (d *HeroGenFarmExtra) Unmarshal(proto *server_proto.HeroGenFarmExtraServerProto) {
	if proto == nil {
		return
	}
	d.dailyStealGold = proto.DailyStealGold
	d.dailyStealStone = proto.DailyStealStone
}

func newHeroGenCountry() *HeroGenCountry {
	return &HeroGenCountry{}
}

// 英雄国家
type HeroGenCountry struct {
	countryId uint64 // Id

	newUserExpiredTime time.Time // 新手 cd 结束时间

	normalExpiredTime time.Time // 非新手 cd 结束时间

	appointOnSameDay bool // 当天刚被任命

	dailySalaryCollected bool // 今天的俸禄已领

}

func (hero *Hero) Country() *HeroGenCountry {
	return hero.heroGen.country
}

func (d *HeroGenCountry) GetCountryId() uint64 {
	return d.countryId
}

func (d *HeroGenCountry) SetCountryId(toSet uint64) uint64 {
	o := d.countryId
	d.countryId = toSet
	return o
}

func (d *HeroGenCountry) GetNewUserExpiredTime() time.Time {
	return d.newUserExpiredTime
}

func (d *HeroGenCountry) SetNewUserExpiredTime(toSet time.Time) time.Time {
	o := d.newUserExpiredTime
	d.newUserExpiredTime = toSet
	return o
}

func (d *HeroGenCountry) GetNormalExpiredTime() time.Time {
	return d.normalExpiredTime
}

func (d *HeroGenCountry) SetNormalExpiredTime(toSet time.Time) time.Time {
	o := d.normalExpiredTime
	d.normalExpiredTime = toSet
	return o
}

func (d *HeroGenCountry) GetAppointOnSameDay() bool {
	return d.appointOnSameDay
}

func (d *HeroGenCountry) SetAppointOnSameDay(toSet bool) bool {
	o := d.appointOnSameDay
	d.appointOnSameDay = toSet
	return o
}

func (d *HeroGenCountry) GetDailySalaryCollected() bool {
	return d.dailySalaryCollected
}

func (d *HeroGenCountry) SetDailySalaryCollected(toSet bool) bool {
	o := d.dailySalaryCollected
	d.dailySalaryCollected = toSet
	return o
}

func (d *HeroGenCountry) Reset() {
	d.appointOnSameDay = false
	d.dailySalaryCollected = false
}

func (d *HeroGenCountry) EncodeClient() *shared_proto.HeroGenCountryProto {
	proto := &shared_proto.HeroGenCountryProto{}
	proto.CountryId = u64.Int32(d.countryId)
	proto.NewUserExpiredTime = timeutil.Marshal32(d.newUserExpiredTime)
	proto.NormalExpiredTime = timeutil.Marshal32(d.normalExpiredTime)
	proto.AppointOnSameDay = d.appointOnSameDay
	proto.DailySalaryCollected = d.dailySalaryCollected

	return proto
}

func (d *HeroGenCountry) EncodeServer() *server_proto.HeroGenCountryServerProto {
	proto := &server_proto.HeroGenCountryServerProto{}
	proto.CountryId = d.countryId
	proto.NewUserExpiredTime = timeutil.Marshal64(d.newUserExpiredTime)
	proto.NormalExpiredTime = timeutil.Marshal64(d.normalExpiredTime)
	proto.AppointOnSameDay = d.appointOnSameDay
	proto.DailySalaryCollected = d.dailySalaryCollected

	return proto
}

func (d *HeroGenCountry) Unmarshal(proto *server_proto.HeroGenCountryServerProto) {
	if proto == nil {
		return
	}
	d.countryId = proto.CountryId
	d.newUserExpiredTime = timeutil.Unix64(proto.NewUserExpiredTime)
	d.normalExpiredTime = timeutil.Unix64(proto.NormalExpiredTime)
	d.appointOnSameDay = proto.AppointOnSameDay
	d.dailySalaryCollected = proto.DailySalaryCollected
}

func newHeroGenShop() *HeroGenShop {
	return &HeroGenShop{}
}

// 商店
type HeroGenShop struct {
	blackMarketDailyRefreshTimes uint64 // 黑市每日刷新次数

	blackMarketNextRefreshTime time.Time // 黑市下次刷新时间

}

func (hero *Hero) Shop() *HeroGenShop {
	return hero.heroGen.shop
}

func (d *HeroGenShop) GetBlackMarketDailyRefreshTimes() uint64 {
	return d.blackMarketDailyRefreshTimes
}

func (d *HeroGenShop) IncBlackMarketDailyRefreshTimes() uint64 {
	d.blackMarketDailyRefreshTimes++
	return d.blackMarketDailyRefreshTimes
}

func (d *HeroGenShop) GetBlackMarketNextRefreshTime() time.Time {
	return d.blackMarketNextRefreshTime
}

func (d *HeroGenShop) SetBlackMarketNextRefreshTime(toSet time.Time) time.Time {
	o := d.blackMarketNextRefreshTime
	d.blackMarketNextRefreshTime = toSet
	return o
}

func (d *HeroGenShop) Reset() {
	d.blackMarketDailyRefreshTimes = 0
}

func (d *HeroGenShop) EncodeClient() *shared_proto.HeroGenShopProto {
	proto := &shared_proto.HeroGenShopProto{}
	proto.BlackMarketDailyRefreshTimes = u64.Int32(d.blackMarketDailyRefreshTimes)
	proto.BlackMarketNextRefreshTime = timeutil.Marshal32(d.blackMarketNextRefreshTime)

	return proto
}

func (d *HeroGenShop) EncodeServer() *server_proto.HeroGenShopServerProto {
	proto := &server_proto.HeroGenShopServerProto{}
	proto.BlackMarketDailyRefreshTimes = d.blackMarketDailyRefreshTimes
	proto.BlackMarketNextRefreshTime = timeutil.Marshal64(d.blackMarketNextRefreshTime)

	return proto
}

func (d *HeroGenShop) Unmarshal(proto *server_proto.HeroGenShopServerProto) {
	if proto == nil {
		return
	}
	d.blackMarketDailyRefreshTimes = proto.BlackMarketDailyRefreshTimes
	d.blackMarketNextRefreshTime = timeutil.Unix64(proto.BlackMarketNextRefreshTime)
}

func newHeroGenGuildWorkshop() *HeroGenGuildWorkshop {
	return &HeroGenGuildWorkshop{}
}

// 联盟工坊
type HeroGenGuildWorkshop struct {
	dailyBuildTimes uint64 // 今日建设次数

	dailyHurtTimes uint64 // 今日破坏次数

	nextHurtTime time.Time // 下次破坏时间

}

func (hero *Hero) GuildWorkshop() *HeroGenGuildWorkshop {
	return hero.heroGen.guild_workshop
}

func (d *HeroGenGuildWorkshop) GetDailyBuildTimes() uint64 {
	return d.dailyBuildTimes
}

func (d *HeroGenGuildWorkshop) IncDailyBuildTimes() uint64 {
	d.dailyBuildTimes++
	return d.dailyBuildTimes
}

func (d *HeroGenGuildWorkshop) GetDailyHurtTimes() uint64 {
	return d.dailyHurtTimes
}

func (d *HeroGenGuildWorkshop) IncDailyHurtTimes() uint64 {
	d.dailyHurtTimes++
	return d.dailyHurtTimes
}

func (d *HeroGenGuildWorkshop) GetNextHurtTime() time.Time {
	return d.nextHurtTime
}

func (d *HeroGenGuildWorkshop) SetNextHurtTime(toSet time.Time) time.Time {
	o := d.nextHurtTime
	d.nextHurtTime = toSet
	return o
}

func (d *HeroGenGuildWorkshop) Reset() {
	d.dailyBuildTimes = 0
	d.dailyHurtTimes = 0
}

func (d *HeroGenGuildWorkshop) EncodeClient() *shared_proto.HeroGenGuildWorkshopProto {
	proto := &shared_proto.HeroGenGuildWorkshopProto{}
	proto.DailyBuildTimes = u64.Int32(d.dailyBuildTimes)
	proto.DailyHurtTimes = u64.Int32(d.dailyHurtTimes)
	proto.NextHurtTime = timeutil.Marshal32(d.nextHurtTime)

	return proto
}

func (d *HeroGenGuildWorkshop) EncodeServer() *server_proto.HeroGenGuildWorkshopServerProto {
	proto := &server_proto.HeroGenGuildWorkshopServerProto{}
	proto.DailyBuildTimes = d.dailyBuildTimes
	proto.DailyHurtTimes = d.dailyHurtTimes
	proto.NextHurtTime = timeutil.Marshal64(d.nextHurtTime)

	return proto
}

func (d *HeroGenGuildWorkshop) Unmarshal(proto *server_proto.HeroGenGuildWorkshopServerProto) {
	if proto == nil {
		return
	}
	d.dailyBuildTimes = proto.DailyBuildTimes
	d.dailyHurtTimes = proto.DailyHurtTimes
	d.nextHurtTime = timeutil.Unix64(proto.NextHurtTime)
}

func newHeroGenProduct() *HeroGenProduct {
	return &HeroGenProduct{}
}

// 渠道相关
type HeroGenProduct struct {
	yuanbaoGiftLimit uint64 // 可赠送的元宝额度

}

func (hero *Hero) Product() *HeroGenProduct {
	return hero.heroGen.product
}

func (d *HeroGenProduct) GetYuanbaoGiftLimit() uint64 {
	return d.yuanbaoGiftLimit
}

func (d *HeroGenProduct) SetYuanbaoGiftLimit(toSet uint64) uint64 {
	o := d.yuanbaoGiftLimit
	d.yuanbaoGiftLimit = toSet
	return o
}

func (d *HeroGenProduct) AddYuanbaoGiftLimit(toAdd uint64) uint64 {
	d.yuanbaoGiftLimit += toAdd
	return d.yuanbaoGiftLimit
}

func (d *HeroGenProduct) SubYuanbaoGiftLimit(toSub uint64) uint64 {
	d.yuanbaoGiftLimit = u64.Sub(d.yuanbaoGiftLimit, toSub)
	return d.yuanbaoGiftLimit
}

func (d *HeroGenProduct) Reset() {
}

func (d *HeroGenProduct) EncodeClient() *shared_proto.HeroGenProductProto {
	proto := &shared_proto.HeroGenProductProto{}
	proto.YuanbaoGiftLimit = u64.Int32(d.yuanbaoGiftLimit)

	return proto
}

func (d *HeroGenProduct) EncodeServer() *server_proto.HeroGenProductServerProto {
	proto := &server_proto.HeroGenProductServerProto{}
	proto.YuanbaoGiftLimit = d.yuanbaoGiftLimit

	return proto
}

func (d *HeroGenProduct) Unmarshal(proto *server_proto.HeroGenProductServerProto) {
	if proto == nil {
		return
	}
	d.yuanbaoGiftLimit = proto.YuanbaoGiftLimit
}

// hero gen
func newHeroGen() *hero_gen {
	return &hero_gen{
		xuanyuan: newHeroGenXuanyuan(),

		misc_data: newHeroGenMiscData(),

		farm_extra: newHeroGenFarmExtra(),

		country: newHeroGenCountry(),

		shop: newHeroGenShop(),

		guild_workshop: newHeroGenGuildWorkshop(),

		product: newHeroGenProduct(),
	}
}

type hero_gen struct {
	xuanyuan *HeroGenXuanyuan // 轩辕会武

	misc_data *HeroGenMiscData // 英雄杂项

	farm_extra *HeroGenFarmExtra // 英雄农场相关数据

	country *HeroGenCountry // 英雄国家

	shop *HeroGenShop // 商店

	guild_workshop *HeroGenGuildWorkshop // 联盟工坊

	product *HeroGenProduct // 渠道相关

}

func (g *hero_gen) resetDaily() {

	g.misc_data.Reset()

	g.country.Reset()

	g.shop.Reset()

	g.guild_workshop.Reset()

	g.product.Reset()

}

func (g *hero_gen) EncodeClient() *shared_proto.HeroGenProto {
	proto := &shared_proto.HeroGenProto{}
	proto.Xuanyuan = g.xuanyuan.EncodeClient()

	proto.MiscData = g.misc_data.EncodeClient()

	proto.FarmExtra = g.farm_extra.EncodeClient()

	proto.Country = g.country.EncodeClient()

	proto.Shop = g.shop.EncodeClient()

	proto.GuildWorkshop = g.guild_workshop.EncodeClient()

	proto.Product = g.product.EncodeClient()

	return proto
}

func (g *hero_gen) EncodeServer() *server_proto.HeroGenServerProto {
	proto := &server_proto.HeroGenServerProto{}
	proto.Xuanyuan = g.xuanyuan.EncodeServer()

	proto.MiscData = g.misc_data.EncodeServer()

	proto.FarmExtra = g.farm_extra.EncodeServer()

	proto.Country = g.country.EncodeServer()

	proto.Shop = g.shop.EncodeServer()

	proto.GuildWorkshop = g.guild_workshop.EncodeServer()

	proto.Product = g.product.EncodeServer()

	return proto
}

func (g *hero_gen) Unmarshal(proto *server_proto.HeroGenServerProto) {
	if proto == nil {
		return
	}
	g.xuanyuan.Unmarshal(proto.Xuanyuan)

	g.misc_data.Unmarshal(proto.MiscData)

	g.farm_extra.Unmarshal(proto.FarmExtra)

	g.country.Unmarshal(proto.Country)

	g.shop.Unmarshal(proto.Shop)

	g.guild_workshop.Unmarshal(proto.GuildWorkshop)

	g.product.Unmarshal(proto.Product)

}
