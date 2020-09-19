package singleton

import (
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/timeutil"
	"time"
	"github.com/lightpaw/male7/util/sortkeys"
	"sort"
	"github.com/lightpaw/male7/config/basedata"
)

//gogen:config
type GuildGenConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"联盟/联盟杂项.txt"`
	_ struct{} `protogen:"true"`

	LeaveAfterJoinDuration time.Duration `default:"4h"` // 成员离开联盟不能重新加入间隔

	GuildMarkCount        uint64 `default:"4"`  // 联盟标记个数
	GuildMarkMsgCharLimit uint64 `default:"20"` // 联盟标记文字个数限制，1个汉字算2个字符

	SendMinYinliangToMember uint64 `default:"10"`    // 单发银两给玩家的最小值
	SendMaxYinliangToMember uint64 `default:"10000"` // 单发银两给玩家的最大值

	SendMinYinliangToGuild uint64 `default:"10"`    // 发银两给联盟最小值
	SendMaxYinliangToGuild uint64 `default:"10000"` // 发银两给联盟最大值

	SendMinSalary uint64 `default:"10"`    // 设置工资最小值
	SendMaxSalary uint64 `default:"10000"` // 设置工资最大值

	ConveneCooldown time.Duration `default:"10s"` // 联盟官员召集CD

	// 修建工坊所需时间
	WorkshopBuildDuration time.Duration `default"72h"`

	// 盟友修建工坊加速时间
	WorkshopHeroBuildDuration time.Duration `default"5m"`

	// 联盟工坊每日加速次数上限
	WorkshopGuildBuildMaxTimes uint64 `default:"100"`

	// 英雄每日修建工坊次数上限
	WorkshopHeroBuildMaxTimes uint64 `default:"5"`

	// 每日工坊生产次数上限（这个不单独处理，根据联盟人数计算）
	// 英雄每日工坊生产次数上限
	WorkshopOutputInitTimes        uint64        `default"2" protofield:"-"`
	WorkshopOutputMaxTimes         uint64        `default"2"`
	WorkshopOutputRecoveryDuration time.Duration `default"1m"`

	// 每次生成的产出值
	WorkshopAddOutput uint64 `default"1" protofield:"-"`

	// 建设每次加的繁荣度
	WorkshopAddProsperity uint64 `default"1" protofield:"-"`

	// 敌人破坏减速时间
	WorkshopHurtDuration time.Duration `default"5m"`

	// 联盟工坊每日破坏次数上限
	WorkshopHurtTotalTimesLimit uint64 `default"20"`

	// 英雄每日破坏敌方工坊次数上限
	WorkshopHurtHeroTimesLimit uint64 `default"3"`

	// 英雄每次破坏CD
	WorkshopHurtCooldown time.Duration `default"1m"`

	// 破坏联盟工坊减少繁荣度
	WorkshopHurtProsperity uint64 `default"1"`

	// 联盟工坊产出所需生产次数
	WorkshopMaxOutput []uint64 `default:"1,2,3" validator:",duplicate" protofield:"-"`

	// 联盟工坊繁荣度上限
	WorkshopProsperityCapcity uint64 `head:"-"`

	// 联盟工坊荒芜繁荣度
	WorkshopBarrenProsperity uint64

	// 联盟工坊奖励初始个数
	WorkshopPrizeInitCount uint64 `validator:"uint" default:"1" protofield:"-"`

	// 联盟工坊奖励个数上限
	WorkshopPrizeMaxCount uint64 `validator:"uint" default:"99"`

	// 联盟工坊多久扣一次繁荣度
	WorkshopReduceProsperityDuration time.Duration `default"1m"`

	// 联盟工坊每次扣多少点繁荣度
	WorkshopReduceProsperity uint64 `default:"1"`

	// 联盟工坊距离自己主城的最大距离限制
	WorkshopDistanceLimit uint64 `default:"100"`

	// 联盟工坊对应的城池数据
	WorkshopBase *basedata.NpcBaseData `protofield:"-"`

	// 联盟转国消耗
	GuildChangeCountryCost *resdata.Cost

	// 联盟转国等待时间
	GuildChangeCountryWaitDuration time.Duration

	// 联盟转国CD
	GuildChangeCountryCooldown time.Duration

	// 联盟周任务开放等级
	TaskOpenLevel uint64
}

func (c *GuildGenConfig) Init(filename string) {

	check.PanicNotTrue(len(c.WorkshopMaxOutput) > 0, "%s 必须配置联盟工坊产出所需生产次数WorkshopMaxOutput", filename)

	c.WorkshopProsperityCapcity = c.WorkshopBase.ProsperityCapcity

}

//gogen:config
type GuildConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"联盟/联盟杂项.txt"`
	_ struct{} `proto:"shared_proto.GuildConfigProto"`
	_ struct{} `protoconfig:"GuildConfig"`

	HebiRobbedSuccessTaskProgress uint64 `default:"1" protofield:"-"` // 合璧干扰成功加的进度
	HebiCompleteTaskProgress      uint64 `default:"3" protofield:"-"` // 合璧完成加的进度

	CreateGuildCost         *resdata.Cost `default:"1"`                               // 创建帮派消耗
	ChangeGuildNameCost     *resdata.Cost `default:"1"`                               // 联盟改名消耗
	ChangeGuildNameDuration time.Duration `default:"72h" parser:"time.ParseDuration"` // 联盟改名CD
	GuildLabelLimitChar     uint64        `default:"5"`                               // 联盟标签字数限制
	GuildLabelLimitCount    uint64        `default:"4"`                               // 联盟标签个数限制
	GuildNumPerPage         uint64        `default:"10"`                              // 帮派列表个数

	TextLimitChar            uint64 `default:"128"` // 对外公告字符限制
	InternalTextLimitChar    uint64 `default:"128"` // 对内公告字符限制
	FriendGuildTextLimitChar uint64 `default:"128"` // 友盟公告字符限制
	EnemyGuildTextLimitChar  uint64 `default:"128"` // 敌盟公告字符限制

	ChangeLeaderCountdownMemberCount uint64        `default:"20"`                              // 联盟人数触发禅让倒计时
	ChangeLeaderCountdownDuration    time.Duration `default:"72h" parser:"time.ParseDuration"` // 禅让盟主倒计时

	ImpeachNpcLeaderHour         uint64        `default:"23"`
	ImpeachNpcLeaderMinute       uint64        `default:"40"`
	ImpeachUserLeaderMemberCount int           `default:"10" protofield:"ImpeachRequiredMemberCount"`    // 弹劾C类盟主最少需要有10个玩家
	ImpeachUserLeaderOffline     time.Duration `default:"48h" protofield:"ImpeachLeaderOfflineDuration"` // 弹劾C类盟主，需要盟主离线时间达到
	ImpeachUserLeaderDuration    time.Duration `default:"12h" protofield:"-"`                            // 弹劾持续时间
	ImpeachExtraCandidateCount   uint64        `default:"2" protofield:"-"`                              // 额外候选人个数

	UserMaxJoinRequestCount  uint64        `default:"5"`
	GuildMaxJoinRequestCount uint64        `default:"20" protofield:"-"`
	JoinRequestDuration      time.Duration `default:"48h" protofield:"-"`

	UserMaxBeenInvateCount uint64        `default:"10" protofield:"-"`
	GuildMaxInvateCount    uint64        `default:"50"`
	InvateDuration         time.Duration `default:"24h" protofield:"-"`

	ContributionDay        int           `default:"7" protofield:"-"`
	NpcKickOfflineDuration time.Duration `default:"24h" protofield:"-"`

	npcSetClassLevelArray []*guild_data.GuildClassLevelData
	leaderClassLevel      *guild_data.GuildClassLevelData
	lowestClassLevel      *guild_data.GuildClassLevelData

	NpcSetClassLevelDuration  time.Duration `default:"14h10m" protofield:"-"` // Npc设置成员职位时间
	npcSetClassLevelDailyTime *timeutil.CycleTime

	FreeNpcGuildKeepCount  uint64 `default:"4" protofield:"-"` // 最少保持多少个空闲npc帮派数量
	FreeNpcGuildEmptyCount uint64 `default:"5" protofield:"-"` // 空虚npc联盟成员空位数量

	DailyMaxKickCount uint64 `default:"3"`

	GuildClassTitleMaxCount     uint64 `default:"10"`
	GuildClassTitleMaxCharCount uint64 `default:"10"` // 支撑最大字符数

	GuildMaxDonateRecordCount uint64 `validator:"int>0" default:"10" protofield:"-"` // 帮派最大捐献记录
	GuildDonateNeedHeroLevel  uint64 `validator:"int>0" default:"4"`                 // 联盟捐献需要的英雄等级

	GuildMaxBigEventCount uint64 `validator:"int>0" default:"10" protofield:"-"` // 帮派最多记录的大事件的记录

	GuildMaxDynamicCount uint64 `validator:"int>0" default:"10" protofield:"-"` // 帮派最多记录的动态的记录

	DefaultGuildCountry *country.CountryData `default:"1" protofield:"-"` // 联盟默认声望目标

	KeepDailyPrestigeCount  uint64 `default:"2" protofield:"-"` // 保存的每日声望的数量
	KeepHourlyPrestigeCount uint64 `default:"24" protofield:"-"` // 保存的每小时核心声望的数量

	//npcGuildTargets    []*guild_data.GuildTarget `protofield:"-"` // npc联盟的联盟目标
	//notNpcGuildTargets []*guild_data.GuildTarget `protofield:"-"` // 非npc联盟的联盟目标
	guildTargetGroups [][]*guild_data.GuildTarget

	UpdateCountryYuanbao          uint64        `default:"2000"` // 更换联盟国家消耗元宝(策划配置字段作废，用 UpdateCountryCost)
	UpdateCountryCost             *resdata.Cost                  // 更换联盟国家消耗
	UpdateCountryDuration         time.Duration `default:"1m"`   // 更换联盟国家CD
	UpdateCountryLostPrestigeCoef float64       `default:"0.2"`  // 更换联盟损失声望系数，0.2表示20%

	template *guild_data.NpcGuildTemplate

	guildNameMap     map[string]struct{}
	guildFlagNameMap map[string]struct{}

	FirstJoinGuildPrize      *resdata.Prize           `default:"1"` // 首次加入联盟奖励
	FirstJoinGuildPrizeProto *shared_proto.PrizeProto `head:"-" protofield:"-"`

	ContributionPerHelp        uint64 `default:"200"` // 每次帮助盟友获得多少联盟贡献
	ContributionMaxCountPerDay uint64 `default:"5"`   // 每日帮助最大可获得联盟贡献次数

	// 联盟大宝箱，必须在联盟时间
	BigBoxCollectableDuration time.Duration `default:"4h"`
	EventPrizeMaxCount        uint64        `default:"300"`

	NotifyJoinGuildDuration                time.Duration `default:"12h"`               // 提醒加入联盟的时间间隔
	NotifyJoinGuildDurationOnOnlineOrLeave time.Duration `default:"10m"`               // 在上线或者离开联盟提醒加入联盟的时间间隔
	NotifyJoinGuildMaxPrestigeRank         uint64        `default:"30" protofield:"-"` // 最高的声望排名

	RecommendInviteHeroCount uint64 `default:"10" protofield:"-"` // 推荐入盟玩家列表长度

	SearchNoGuildHerosPerPageSize uint64 `default:"10"` // 联盟模糊搜索，每页条数

	SearchNoGuildHerosDuration time.Duration `default:"1s"` // 联盟模糊搜索，请求最小间隔

	// 联盟转国等待时间
	GuildChangeCountryWaitDuration time.Duration `protofield:"-"`

	prizeDataMap map[shared_proto.HeroEvent][]*guild_data.GuildEventPrizeData

	prestigeEventMap map[shared_proto.HeroEvent][]*guild_data.GuildPrestigeEventData
}

func (c *GuildConfig) Init(filename string, configs interface {
	GetGuildClassLevelDataArray() []*guild_data.GuildClassLevelData
	GetNpcGuildTemplateArray() []*guild_data.NpcGuildTemplate
	GetGuildTargetArray() []*guild_data.GuildTarget
	GetGuildEventPrizeDataArray() []*guild_data.GuildEventPrizeData
	GetGuildPrestigeEventDataArray() []*guild_data.GuildPrestigeEventData
}) {
	c.npcSetClassLevelDailyTime = timeutil.NewOffsetDailyTime(int64(c.NpcSetClassLevelDuration / time.Second))

	if n := len(configs.GetGuildClassLevelDataArray()); n > 0 {
		// 去掉最高最低，倒叙
		c.npcSetClassLevelArray = nil
		for i := n - 2; i > 0; i-- {
			c.npcSetClassLevelArray = append(c.npcSetClassLevelArray, configs.GetGuildClassLevelDataArray()[i])
		}

		c.leaderClassLevel = configs.GetGuildClassLevelDataArray()[n-1]
		c.lowestClassLevel = configs.GetGuildClassLevelDataArray()[0]
		check.PanicNotTrue(c.leaderClassLevel != c.lowestClassLevel, "联盟职位只有1个？")
	}

	c.guildNameMap = make(map[string]struct{})
	c.guildFlagNameMap = make(map[string]struct{})
	for _, t := range configs.GetNpcGuildTemplateArray() {

		if t.RejectUserJoin {
			// A类
			_, ok := c.guildNameMap[t.Name]
			check.PanicNotTrue(!ok, "Npc联盟中名字存在重复，%s", t.Name)
			c.guildNameMap[t.Name] = struct{}{}

			_, ok = c.guildFlagNameMap[t.FlagName]
			check.PanicNotTrue(!ok, "Npc联盟中旗号存在重复，%s", t.FlagName)
			c.guildFlagNameMap[t.FlagName] = struct{}{}
		} else {
			// B类
			if c.template == nil {
				c.template = t
				break
			}

			for _, name := range t.GetCombineNames() {
				_, ok := c.guildNameMap[name]
				check.PanicNotTrue(!ok, "Npc联盟中名字存在重复，%s", name)
				c.guildNameMap[name] = struct{}{}

				_, ok = c.guildFlagNameMap[t.FlagName]
				check.PanicNotTrue(!ok, "Npc联盟中旗号存在重复，%s", t.FlagName)
				c.guildFlagNameMap[t.FlagName] = struct{}{}
			}

		}
	}

	check.PanicNotTrue(c.template != nil, "至少要配置一个B类联盟模板")

	// 设置联盟目标
	targetGroupMap := make(map[uint64][]*guild_data.GuildTarget)
	for _, target := range configs.GetGuildTargetArray() {
		targetGroupMap[target.Group] = append(targetGroupMap[target.Group], target)
	}

	for _, ts := range targetGroupMap {

		var kvs []*sortkeys.U64KV
		for _, t := range ts {
			kvs = append(kvs, sortkeys.NewU64KV(t.Order, t))
		}
		sort.Sort(sortkeys.U64KVSlice(kvs))

		for i, kv := range kvs {
			ts[i] = kv.V.(*guild_data.GuildTarget)
		}

		c.guildTargetGroups = append(c.guildTargetGroups, ts)
	}

	check.PanicNotTrue(len(c.guildTargetGroups) > 0, "guild/guild_target.txt 没有配置一个非npc联盟的联盟目标!")

	prizeDataMap := make(map[shared_proto.HeroEvent][]*guild_data.GuildEventPrizeData)
	for _, data := range configs.GetGuildEventPrizeDataArray() {
		array := prizeDataMap[data.TriggerEvent]
		array = append(array, data)
		prizeDataMap[data.TriggerEvent] = array
	}
	c.prizeDataMap = prizeDataMap

	prestigeEventMap := make(map[shared_proto.HeroEvent][]*guild_data.GuildPrestigeEventData)
	for _, data := range configs.GetGuildPrestigeEventDataArray() {
		if data.TriggerEventCondition.Greater || data.TriggerEventCondition.Less {
			array := prestigeEventMap[data.TriggerEvent]
			array = append(array, data)
			prestigeEventMap[data.TriggerEvent] = array
		}
	}
	c.prestigeEventMap = prestigeEventMap

	if c.FirstJoinGuildPrize != nil {
		c.FirstJoinGuildPrizeProto = c.FirstJoinGuildPrize.Encode4Init()
	}
}

func (c *GuildConfig) GetEventPrizes(t shared_proto.HeroEvent) []*guild_data.GuildEventPrizeData {
	return c.prizeDataMap[t]
}

func (c *GuildConfig) GetPrestigeEvent(t shared_proto.HeroEvent) []*guild_data.GuildPrestigeEventData {
	return c.prestigeEventMap[t]
}

func (c *GuildConfig) ExistName(name string) bool {
	_, ok := c.guildNameMap[name]
	return ok
}

func (c *GuildConfig) ExistFlagName(name string) bool {
	_, ok := c.guildFlagNameMap[name]
	return ok
}

func (c *GuildConfig) GetTemplate() *guild_data.NpcGuildTemplate {
	return c.template
}

func (c *GuildConfig) GetPrevNpcSetClassLevelTime(ctime time.Time) time.Time {
	return c.npcSetClassLevelDailyTime.PrevTime(ctime)
}

func (c *GuildConfig) GetNextNpcSetClassLevelTime(ctime time.Time) time.Time {
	return c.npcSetClassLevelDailyTime.NextTime(ctime)
}

func (c *GuildConfig) GetNpcSetClassLevelArray() []*guild_data.GuildClassLevelData {
	return c.npcSetClassLevelArray
}

func (c *GuildConfig) GetLeaderClassLevel() *guild_data.GuildClassLevelData {
	return c.leaderClassLevel
}

func (c *GuildConfig) GetLowestClassLevel() *guild_data.GuildClassLevelData {
	return c.lowestClassLevel
}

func (c *GuildConfig) GetImpeachLeaderTime(ctime time.Time, isNpcLeader bool) time.Time {

	if isNpcLeader {
		// npc盟主，到一个固定时间点
		todayZeroTime := timeutil.DailyTime.PrevTime(ctime)
		resetTime := todayZeroTime.Add(time.Duration(c.ImpeachNpcLeaderHour)*time.Hour + time.Duration(c.ImpeachNpcLeaderMinute)*time.Minute)
		return resetTime
	} else {
		// 玩家盟主，当前时间往后推一段时间
		return ctime.Add(c.ImpeachUserLeaderDuration)
	}
}

func (c *GuildConfig) GetGuildTargetGroups() [][]*guild_data.GuildTarget {
	return c.guildTargetGroups
}
