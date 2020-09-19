package singleton

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

const (
	day  = 24 * time.Hour
	week = 7 * 24 * time.Hour
)

//gogen:config
type MiscGenConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"杂项/杂项.txt"`
	_ struct{} `protogen:"true"`

	DianquanToGold  uint64 `default:"10000"` // 1点券兑换的钱币
	DianquanToStone uint64 `default:"10000"` // 1点券兑换的石头

	Stronger4Coef []float64 `desc:"变强系数4" validator:"float64>=0" default:"0,1,1.5,2"`

	MiaoBaowuDuration   time.Duration      `desc:"秒宝物开启CD间隔"`
	MiaoBaowuCost       *resdata.Cost      `desc:"秒宝物开启CD消耗"`
	DailyMiaoBaowuLimit uint64             `validator:"uint" default:"3" desc:"每日秒宝物次数上限"`
	BaowuLogLimit       int                `default:"40" desc:"宝物日志条数上限"`
	FirstNpcBaowu       *resdata.BaowuData `default:"nullable" protofield:",config.U64ToI32(%s.Id),int32"` // 第一次从npc处获得的宝物

	FriendMaxCount uint64 `desc:"好友最大数量" default:"150"`
	BlackMaxCount  uint64 `desc:"黑名单最大数量" default:"150"`

	// 每次购买体力值增加的体力值
	BuySpValue uint64
	// 购买体力值消耗（元宝）
	BuySpCost uint64
	// 每日体力值购买次数上限
	BuySpLimit uint64

	// 体力值恢复间隔
	SpDuration time.Duration `default:"5m"`

	// 驻防自动补兵时间间隔
	AutoFullSoldoerDuration time.Duration `default:"5m"`

	// 加税收的间隔
	TaxDuration time.Duration `default:"30s"`

	// 钓鱼积分最大值
	FishMaxPoint []uint64 `validator:"uint"`

	// 钓鱼积分兑换将魂列表（现在仅前端作展示用）
	FishPointCaptain []*resdata.ResCaptainData `default:"nullable" protofield:"FishPointCaptainSoul,config.U64a2I32a(resdata.GetResCaptainDataKeyArray(%s)),int32"`

	FishPointPlunder *resdata.Plunder `protofield:"-"`

	DaZhaoSwitchLevelLimit uint64 `default:"4"` // 大招切换等级限制

	// 外城改建消耗
	UpdateOuterCityTypeCost *resdata.Cost

	// 聊天初始推送数量
	FirstHistoryChatSend int `default:"2" desc:"聊天初始推送数量" protofield:"-"`

	// 随机事件生成数量
	RandomEventNum int `default:"729" desc:"随机事件生成数量" protofield:"-"`

	// 随机事件自单元场景补足数量
	RandomEventOwnMinNum int `default:"3" desc:"随机事件自单元场景补足数量" protofield:"-"`

	// 随机事件大刷新间隔
	RandomEventBigRefreshDuration time.Duration `default:"6h"`

	// 随机事件小刷新间隔
	RandomEventSmallRefreshDuration time.Duration `default:"30m"`

	// 每日对相同玩家施计次限
	TargetUseStratagemLimit uint64 `default:"2"  protofield:"-"`

	// 中计次限
	TrappedStratagemLimit uint64 `default:"5"  protofield:"-"`

	// 武将重生消耗
	CaptainResetCost *resdata.Cost

	SkipFightingHeroLevel    uint64        `default:"10"`  // 跳过战斗的君主等级
	SkipFightingVipLevel     uint64        `default:"10"`  // 跳过战斗的Vip等级
	SkipFightingWaitDuration time.Duration `default:"10s"` // 跳过战斗等待时间，秒

	SecretTowerCd time.Duration `default:"30s"`
	XuanyuanCd    time.Duration `default:"30s"`
	BaizhanCd     time.Duration `default:"30s"`
	HebiCd        time.Duration `default:"30s"`
	XiongnuCd     time.Duration `default:"30s"`
	MailCd        time.Duration `default:"30s"`

	// 每个兵阵多少兵
	SoldierPerGroup uint64 `default:"2000"`

	// 召唤殷墟持续时间
	HeroBaozDuration    time.Duration `default:"2h"`
	HeroBaozMaxDistance uint64        `default:"200"`

	YuanbaoGiftPercent *data.Amount `default:"500%" desc:"可赠送元宝占充值元宝的最大比例。千分比"`
}

func (c *MiscGenConfig) Init(filename string) {

	check.PanicNotTrue(c.FishPointPlunder != nil, "%s 钓鱼积分掉落未配置！", filename)

	if c.FirstNpcBaowu != nil {
		check.PanicNotTrue(c.FirstNpcBaowu.CantRob, "%s 配置的首次掠夺Npc宝物[%d %s]，必须是一个不能被抢夺的宝物", filename, c.FirstNpcBaowu.Id, c.FirstNpcBaowu.Name)
		check.PanicNotTrue(c.FirstNpcBaowu.GetNextLevel() == nil, "%s 配置的首次掠夺Npc宝物[%d %s]，必须是一个不能合成下一级的宝物", filename, c.FirstNpcBaowu.Id, c.FirstNpcBaowu.Name)
	}

	check.PanicNotTrue(c.YuanbaoGiftPercent.Percent <= 1000, "%v YuanbaoGiftPercent：%v% 必须 < 1000%", filename, c.YuanbaoGiftPercent.Percent)
}

//gogen:config
type MiscConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"杂项/杂项.txt"`
	_ struct{} `proto:"shared_proto.MiscConfigProto"`
	_ struct{} `protoconfig:"MiscConfig"`

	// 日常重置
	DailyResetHour     uint64        `validator:"uint" default:"0"`
	DailyResetMinute   uint64        `validator:"uint" default:"0"`
	DailyResetDuration time.Duration `head:"-" protofield:"-"`
	DailyResetTime     *timeutil.CycleTime

	// 周常重置
	WeeklyResetHour     uint64        `validator:"uint" default:"0"`
	WeeklyResetMinute   uint64        `validator:"uint" default:"0"`
	WeeklyResetDuration time.Duration `head:"-" protofield:"-"`
	WeeklyResetTime     *timeutil.CycleTime

	MinNameCharLen uint64 `default:"2" protofield:"-"`
	MaxNameCharLen uint64 `default:"14" protofield:"-"`

	WorkshopRefreshHourMinute []uint64        `validator:"uint" default:"1200,1600,2000"`
	WorkshopRefreshDuration   []time.Duration `head:"-" protofield:"-"`

	SecondWorkerCost           *resdata.Cost // 解锁建筑二队的花费
	SecondWorkerUnlockDuration time.Duration // 解锁建筑二队的有效期

	ChangeHeroNameYuanbaoCost []uint64        `validator:"uint" default:"0,100"` // 改名消耗，策划配置字段作废，用 ChangeHeroNameCost
	ChangeHeroNameCost        []*resdata.Cost `validator:"uint"`                 // 改名消耗
	ChangeHeroNameDuration    time.Duration
	FirstChangeHeroNamePrize  *resdata.Prize `default:"null"` // 首次改名的奖励

	ChangeCaptainNameCost     *resdata.Cost `default:"1"` // 武将改名消耗
	ChangeCaptainRaceDuration time.Duration               // 武将转职cd
	ChangeCaptainRaceLevel    uint64 `default:"50"`       // 武将转职等级

	MailMinBatchCount           uint64        `default:"20"`                                // 邮件批量个数
	MailMaxBatchCount           uint64        `default:"50"`                                // 邮件批量个数
	TowerChallengeMaxTimes      uint64        `default:"3"`                                 // 千重楼最大挑战次数
	TowerResetChallengeDuration time.Duration `default:"30m" protofield:"-"`                // 千重楼挑战次数重置
	TowerReplayCount            uint64        `default:"4" protofield:"-"`                  // 千重楼回放个数
	TowerAutoKeepFloor          uint64        `validator:"uint" default:"5" protofield:"-"` // 千重楼扫荡减少层数
	EquipmentUpgradeMultiTimes  uint64        `default:"10" protofield:"-"`                 // 批量升级次数

	MiaoBuildingWorkerDuration time.Duration
	MiaoBuildingWorkerCost     *resdata.Cost // 秒建筑队消耗
	MiaoTechWorkerDuration     time.Duration
	MiaoTechWorkerCost         *resdata.Cost // 秒科研队消耗

	MiaoCaptainRebirthDuration time.Duration
	MiaoCaptainRebirthCost     *resdata.Cost // 秒武将转生消耗

	MiaoWorkshopDuration time.Duration
	MiaoWorkshopCost     *resdata.Cost // 秒锻造CD消耗

	// 初始锻造的次数
	DefaultForgingTimes uint64 `validator:"uint" default:"3" protofield:"-"`

	MaxDepotEquipCapacity uint64 `default:"100"` // 背包装备的最大容量

	TempDepotExpireDuration time.Duration `default:"24h" protofield:"-" parser:"time.ParseDuration"` // 临时背包的过期时间

	MaxSignLen  uint64 `default:"20"`                  // 最大的签名长度
	MaxVoiceLen uint64 `default:"2000" protofield:"-"` // 最大的语音长度

	StrategyRestoreDuration time.Duration `default:"1h"`

	MaxResourceCollectTimes     uint64        `validator:"uint" default:"10"`              // 最大资源采集次数
	DefaultResourceCollectTimes uint64        `validator:"int" default:"5" protofield:"-"` // 默认资源采集次数
	ResourceRecoveryDuration    time.Duration `default:"1h"`                               // 采集资源恢复间隔

	MondayZeroOClock int64 `head:"-"` // 周一的0点

	MaxFavoritePosCount uint64 `validator:"int>0" default:"10"`

	FlagHeroName *data.TextFormatter `default:"[%s];%s" protofield:"-"`

	CountryFlagHeroName *data.TextFormatter `default:"[%s];[%s];%s" protofield:"-"`

	// 聊天
	WorldChatLevel      uint64           `default:"1"`   // 世界聊天等级限制, >= 这个等级才能发世界聊天
	WorldChatDuration   time.Duration    `default:"10s"` // 世界聊天CD，每条聊天最少相隔X秒
	ChatTextLength      uint64           `default:"100"` // 聊天文字长度限制（英文=1字节，中文=2字节），必须<=这个值
	ChatJsonLength      uint64           `default:"600" protofield:"-"`
	ChatWindowLimit     uint64           `default:"100" protofield:"-"`
	ChatBatchCount      uint64           `default:"100" protofield:"-"`
	BroadcastGoods      *goods.GoodsData `protofield:",config.U64ToI32(%s.Id)"` // 广播小喇叭
	ChatPrivateMinLevel uint64           `default:"3"`                          // 私聊等级限制
	ChatDuration        time.Duration    `default:"5s"`                         // 通用聊天 CD
	ChatShareDuration   time.Duration    `default:"5s"`                         // 聊天分享 CD

	GiveAddCountDownPrizeDuration time.Duration `protofield:"-"` // 给插入主城繁荣度受损的奖励的时间间隔
	GiveAddCountDownPrizeTimes    uint64        `protofield:"-"` // 给插入主城繁荣度受损的奖励的次数

	StrongerCoef []float64 `validator:"float64>=0" default:"0,0.6,1"`

	BuildingInitEffect uint64 `default:"1"` // 建筑初始加成效果

	DbWorldChatExpireDuration   time.Duration `default:"72h" protofield:"-"`
	DbGuildChatExpireDuration   time.Duration `default:"168h" protofield:"-"`
	DbPrivateChatExpireDuration time.Duration `default:"2160h" protofield:"-"`
	DbGuildLogCountLimit        uint64        `default:"100" protofield:"-"`
	DbMailCountLimit            uint64        `default:"300" protofield:"-"`

	ExtraResDecayCoef     *data.Amount  `default:"30%" protected:"-"` // 仓库超出上限时的损耗系数
	ExtraResDecayDuration time.Duration `default:"6m"`                // 仓库超出上限时的损耗间隔

	CloseFightGuideDungeonId uint64 `default:"107"` // 关闭战前部署克制界面需要的幻境 id

	RefreshRecommendHeroDuration       time.Duration `default:"3s"`
	RefreshRecommendHeroPageSize       uint64        `default:"9"`
	RefreshRecommendHeroMinLevel       uint64        `default:"3"`
	RefreshRecommendHeroPageCount      uint64        `default:"10"`
	SearchHeroDuration                 time.Duration `default:"3s"`
	RecommendHeroOfflineExpireDuration time.Duration `default:"24h"`

	RedPacketServerDelDuration   time.Duration `protofield:"-" default:"168h"`
	RedPacketGuildMemberMinCount uint64        `default:"1"`
}

func (c *MiscConfig) Init(filename string) {
	check.PanicNotTrue(c.DailyResetHour < 24, "%s 每日重置小时数必须[0 <= hour < 24], hour: %v", filename, c.DailyResetHour)
	check.PanicNotTrue(c.DailyResetMinute < 60, "%s 每日重置分钟数必须[0 <= minute < 60], minute: %v", filename, c.DailyResetMinute)

	c.DailyResetDuration = time.Duration(c.DailyResetHour)*time.Hour + time.Duration(c.DailyResetMinute)*time.Minute
	c.DailyResetTime = timeutil.NewOffsetDailyTime(int64(c.DailyResetDuration / time.Second))

	oneDay, _ := time.ParseDuration("24h")
	check.PanicNotTrue(oneDay > c.DailyResetDuration, "%s 每日重置时间配置不能大于24小时，hour:%v minute: %v", filename, c.DailyResetHour, c.DailyResetMinute)

	check.PanicNotTrue(c.WeeklyResetHour < 168, "%s 每周重置小时数必须[0 <= hour < 168], hour: %v", filename, c.WeeklyResetHour)
	check.PanicNotTrue(c.WeeklyResetMinute < 60, "%s 每周重置分钟数必须[0 <= minute < 60], minute: %v", filename, c.WeeklyResetMinute)

	c.WeeklyResetDuration = time.Duration(c.WeeklyResetHour)*time.Hour + time.Duration(c.WeeklyResetMinute)*time.Minute
	c.WeeklyResetTime = timeutil.NewOffsetWeeklyTime(int64(c.WeeklyResetDuration / time.Second))

	check.PanicNotTrue(!check.Uint64Duplicate(c.WorkshopRefreshHourMinute), "%s 装备作坊配置了重复的刷新时间, hour minute: %v", filename, c.WorkshopRefreshHourMinute)

	u64.Sort(c.WorkshopRefreshHourMinute)
	for _, hhmm := range c.WorkshopRefreshHourMinute {
		hh := hhmm / 100
		mm := hhmm % 100

		check.PanicNotTrue(hh < 24, "%s 装备作坊刷新小时数必须[0 <= hour < 24], hour: %v", filename, hh)
		check.PanicNotTrue(mm < 60, "%s 装备作坊刷新分钟数必须[0 <= minute < 60], minute: %v", filename, mm)

		dd := time.Duration(hh)*time.Hour + time.Duration(mm)*time.Minute
		c.WorkshopRefreshDuration = append(c.WorkshopRefreshDuration, dd)
	}

	check.PanicNotTrue(len(c.ChangeHeroNameCost) > 0, "%s 配置的改名消耗起码也要配置一条吧! 条数: %v", filename, len(c.ChangeHeroNameCost))

	for idx, cost := range c.ChangeHeroNameCost {
		if idx > 0 {
			check.PanicNotTrue(cost != nil, "%v 改名消耗，除了第一次，其他不能为空。第%v次", filename, idx+1)
		}
	}

	check.PanicNotTrue(c.MaxDepotEquipCapacity > 0, "%s 配置的 背包中装备的最大容量必须>0!%d", filename, c.MaxDepotEquipCapacity)
	check.PanicNotTrue(c.TempDepotExpireDuration > 0, "%s 配置的 临时背包的过期时间必须>0!%d", filename, c.TempDepotExpireDuration)

	// 计算一个周一的0点
	c.MondayZeroOClock = timeutil.WeekCycleTime(time.Monday).PrevTime(time.Now()).Unix()

	check.PanicNotTrue(c.BroadcastGoods.DianquanPrice > 0 || c.BroadcastGoods.YuanbaoPrice > 0, "%v 广播小喇叭的元宝或点券价格必须>0. goods:%v price:%v", filename, c.BroadcastGoods.Id, c.BroadcastGoods.DianquanPrice)
}

func (c *MiscConfig) GetChangeHeroNameCost(changeTimes uint64) *resdata.Cost {
	if changeTimes < 0 {
		return c.ChangeHeroNameCost[0]
	}

	if changeTimes >= uint64(len(c.ChangeHeroNameCost)) {
		return c.ChangeHeroNameCost[len(c.ChangeHeroNameCost)-1]
	}

	return c.ChangeHeroNameCost[changeTimes]
}

func (c *MiscConfig) GetNextResetDailyTime(ctime time.Time) time.Time {
	return GetNextResetDailyTime(ctime, c.DailyResetDuration)
}

func (c *MiscConfig) GetNextResetWeeklyTime(ctime time.Time) time.Time {
	return GetNextResetWeeklyTime(ctime, c.WeeklyResetDuration)
}

func (c *MiscConfig) GetWorkshopNextRefreshTime(ctime time.Time) time.Time {
	// 多个里面找一个时间最小的
	t := ctime.Add(day)
	for _, d := range c.WorkshopRefreshDuration {
		nextTime := GetNextResetDailyTime(ctime, d)
		if nextTime.Before(t) {
			t = nextTime
		}
	}

	return t
}

// 每日0点重置，返回下一天的0点（传入今日0点，也返回下一天0点）
func GetNextResetDailyTime(ctime time.Time, resetDuration time.Duration) time.Time {

	todayZeroTime := timeutil.DailyTime.PrevTime(ctime)
	resetTime := todayZeroTime.Add(resetDuration)

	if ctime.Before(resetTime) {
		return resetTime
	}

	return resetTime.Add(day)
}

func GetNextResetWeeklyTime(ctime time.Time, resetDuration time.Duration) time.Time {
	zeroTime := timeutil.Sunday.PrevTime(ctime)
	resetTime := zeroTime.Add(resetDuration)
	if ctime.Before(resetTime) {
		return resetTime
	}
	return resetTime.Add(week)
}
