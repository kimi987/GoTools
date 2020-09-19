package promdata

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/timeutil"
	"time"
)

//gogen:config
type PromotionMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"福利/福利杂项.txt"`
	_ struct{} `protogen:"true"`

	// 君主等级基金购买消耗
	HeroLevelFundCost *resdata.Cost
}

//gogen:config
type HeroLevelFundData struct {
	_ struct{} `file:"福利/君主等级基金.txt"`
	_ struct{} `protogen:"true"`

	// 君主等级
	Level uint64
	// 返利
	Rebate uint64
	// 奖励
	Prize *resdata.Prize
}

//gogen:config
type SpCollectionData struct {
	_ struct{} `file:"福利/体力领取.txt"`
	_ struct{} `protogen:"true"`

	// 宴席id
	Id uint64
	// 名称
	Name string
	// 图标
	Icon string
	// 时间段
	TimeShow string
	// 从0点为基准的开始相对时间
	StartDuration time.Duration
	// 从0点为基准的结束相对时间
	EndDuration time.Duration
	// 获得体力（作废，可作展示用）
	Sp uint64
	// 补领贵族等级
	RepairVip uint64
	// 奖励
	SpPrize *resdata.Prize `head:"prize"`
	// 补领消耗
	RepairCost *resdata.Cost
}

//gogen:config
type DailyBargainData struct {
	_ struct{} `file:"福利/每日特惠.txt"`
	_ struct{} `protogen:"true"`

	// 特惠id
	Id uint64
	// 名称
	Name string
	// 购买即可获得N元宝
	GiveYuanbao uint64 `validator:"uint" default:"0"`
	// 还可以获得价值N元宝大礼包
	ShowYuanbao uint64 `validator:"uint" default:"0"`
	// 每日限购(领)
	Limit uint64
	// 充值金额
	ChargeAmount uint64 `validator:"int>0"`
	// 奖励包，prize里面要配上give_yuanbao的值
	Prize *resdata.Prize
}

//gogen:config
type DurationCardData struct {
	_ struct{} `file:"福利/尊享卡.txt"`
	_ struct{} `protogen:"true"`

	// 尊享卡id
	Id uint64
	// 卡名
	Name string
	// 图标
	Icon string
	// 描述
	Desc string
	// 持续时间（为0则永久生效）
	Duration time.Duration `desc:"为0则永久生效（永久卡）"`
	// 充值金额
	ChargeAmount uint64
	// 购买后立即获得的奖励
	Prize *resdata.Prize
	// 每日获赠（领取）的奖励
	DailyPrize *resdata.Prize
	// 离到期前间隔多久提示续费
	BeforePromptDuration time.Duration
}

//gogen:config
type FreeGiftData struct {
	_ struct{} `file:"福利/免费礼包.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"promotion.proto"`

	Id       uint64
	Name     string
	GiftType shared_proto.GiftType
	Daily    bool `desc:"false则为一次性奖励礼包，否则为每日刷新的礼包" default:"false"`
	Prize    *resdata.Prize
}

// 时限礼包组
//gogen:config
type TimeLimitGiftGroupData struct {
	_ struct{} `file:"福利/定时时限礼包组.txt"`
	_ struct{} `proto:"shared_proto.TimeLimitGiftGroupDataProto"`
	_ struct{} `protoconfig:"TimeLimitGiftGroups"`

	Id       uint64 // 礼包组id
	TimeRule *data.TimeRuleData   `protofield:"-"`
	Gifts    []*TimeLimitGiftData `validator:"uint,notAllNil,duplicate" protofield:"GiftIds,config.U64a2I32a(GetTimeLimitGiftDataKeyArray(%s))"`

	MinHeroLevel   uint64 `validator:"uint" default:"0"` // 推送所需英雄等级下限
	MaxHeroLevel   uint64 `validator:"uint" default:"0"` // 推送所需英雄等级上限
	MinGuanfuLevel uint64 `validator:"uint" default:"0"` // 推送所需官府等级上限
	MaxGuanfuLevel uint64 `validator:"uint" default:"0"` // 推送所需官府等级上限
}

func (d *TimeLimitGiftGroupData) Init(filepath string, configs interface {
	GetGuanFuLevelDataArray() []*domestic_data.GuanFuLevelData
}) {
	heroMaxLevel := uint64(100)
	check.PanicNotTrue(d.MaxHeroLevel <= heroMaxLevel, "%s 礼包组ID：%v [max_hero_level]:%d 超出君主等级上限", filepath, d.Id, d.MaxHeroLevel)
	check.PanicNotTrue(d.MinHeroLevel <= heroMaxLevel, "%s 礼包组ID：%v [min_hero_level]:%d 超出君主等级上限", filepath, d.Id, d.MinHeroLevel)
	if d.MaxHeroLevel <= 0 {
		d.MaxHeroLevel = heroMaxLevel
	}
	check.PanicNotTrue(d.MinHeroLevel <= d.MaxHeroLevel, "%s 礼包组ID：%v [min_hero_level]:%d 不允许大于 [max_hero_level]:%d", filepath, d.Id, d.MinHeroLevel, d.MaxHeroLevel)

	guanfuLevelArr := configs.GetGuanFuLevelDataArray()
	guanfuMaxLevel := guanfuLevelArr[len(guanfuLevelArr)-1].Level
	check.PanicNotTrue(d.MaxGuanfuLevel <= guanfuMaxLevel, "%s 礼包组ID：%v [max_guanfu_level]:%d 超出官府等级上限", filepath, d.Id, d.MaxGuanfuLevel)
	check.PanicNotTrue(d.MinGuanfuLevel <= guanfuMaxLevel, "%s 礼包组ID：%v [min_guanfu_level]:%d 超出官府等级上限", filepath, d.Id, d.MinGuanfuLevel)
	if d.MaxGuanfuLevel <= 0 {
		d.MaxGuanfuLevel = guanfuMaxLevel
	}
	check.PanicNotTrue(d.MinGuanfuLevel <= d.MaxGuanfuLevel, "%s 礼包组ID：%v [min_guanfu_level]:%d 不允许大于 [max_guanfu_level]:%d", filepath, d.Id, d.MinGuanfuLevel, d.MaxGuanfuLevel)
}

// 礼包组是否可以启动，如果可以启动，返回可以启动的TimeRule的结束时间以及时间规则id作为refreshId
func (d *TimeLimitGiftGroupData) IsCanOpenUp(serverStartTime, ctime time.Time) (endTime time.Time, ok bool) {
	openUpTime := d.TimeRule.Next(serverStartTime, ctime)
	if openUpTime.IsZero() {
		return
	}

	endTime = openUpTime.Add(d.TimeRule.TimeDuration)
	if timeutil.Between(ctime, openUpTime, endTime) {
		ok = true
	}
	return
}

// 礼包是否可以推送(购买)
func (d *TimeLimitGiftGroupData) IsCanBuy(heroLevel uint64, guanfuLevel uint64) bool {
	return heroLevel >= d.MinHeroLevel && heroLevel <= d.MaxHeroLevel && guanfuLevel >= d.MinGuanfuLevel && guanfuLevel <= d.MaxGuanfuLevel
}

// 超值时限礼包
//gogen:config
type TimeLimitGiftData struct {
	_ struct{} `file:"福利/定时时限礼包.txt"`
	_ struct{} `protogen:"true"`

	Id           uint64                                // 礼包id
	Name         string                                // 礼包名
	Icon         string                                // 礼包icon
	Image        string                                // 礼包插画
	Desc         string                                // 礼包文本
	YuanbaoPrice uint64                                // 现价(元宝)，不走cost表了，确定是元宝，否则old_price也无法展示
	OldPrice     uint64                                // 原价（展示用）
	SignIcon     string                                // 显示的标签（文件路径）
	SignName     string                                // 标签内显示的文本
	DiscountIcon string                                // 有红热和银亮等折扣标签
	Discount     string                                // 标签内折扣值
	Dianquan     uint64                                // 购买后立即获得的点券
	Priority     uint64 `validator:"uint" default:"0"` // 贴脸优先级，取最高值贴脸，如果是0则不贴脸
	Prize        *resdata.Prize
	BuyLimit     uint64 `validator:"uint" default:"0"` // 该种礼包购买次限（不是一次性购买数量，而是购买次数，每次只能买1个），填0不作限制
	Sort         uint64 `validator:"uint" default:"0"` // 客户端用于排序

	ShowPrize         *resdata.Prize `desc:"展示奖励"`
	GuildEventPrizeId uint64         `desc:"联盟礼包Id" default:"0" validator:"uint"`
}

//gogen:config
type EventLimitGiftConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"福利/事件时限礼包.txt"`

	homeBaseGifts map[uint64]*EventLimitGiftData // 主城建设礼包

	captainGifts map[uint64]*EventLimitGiftData // 武将培养礼包

	towerGifts map[uint64]*EventLimitGiftData // 爬塔勇士礼包

	heroLevelGifts map[uint64]*EventLimitGiftData // 君主等级礼包

	rebornGift *EventLimitGiftData // 重建礼包

	supplyGift *EventLimitGiftData // 前线补给礼包

	towerHelpGift *EventLimitGiftData // 爬楼协助礼包
}

func (c *EventLimitGiftConfig) GetHomeBaseGift(baseLevel uint64) *EventLimitGiftData {
	return c.homeBaseGifts[baseLevel]
}

func (c *EventLimitGiftConfig) GetCaptainGift(rarityId uint64) *EventLimitGiftData {
	return c.captainGifts[rarityId]
}

func (c *EventLimitGiftConfig) GetTowerGift(floor uint64) *EventLimitGiftData {
	return c.towerGifts[floor]
}

func (c *EventLimitGiftConfig) GetHeroLevelGift(heroLevel uint64) *EventLimitGiftData {
	return c.heroLevelGifts[heroLevel]
}

func (c *EventLimitGiftConfig) GetRebornGift() *EventLimitGiftData {
	return c.rebornGift
}

func (c *EventLimitGiftConfig) GetSupplyGift() *EventLimitGiftData {
	return c.supplyGift
}

func (c *EventLimitGiftConfig) GetTowerHelpGift() *EventLimitGiftData {
	return c.towerHelpGift
}

func (c *EventLimitGiftConfig) Init(filename string, configs interface {
	GetEventLimitGiftDataArray() []*EventLimitGiftData
}) {
	c.homeBaseGifts = make(map[uint64]*EventLimitGiftData)
	c.captainGifts = make(map[uint64]*EventLimitGiftData)
	c.towerGifts = make(map[uint64]*EventLimitGiftData)
	c.heroLevelGifts = make(map[uint64]*EventLimitGiftData)

	for _, data := range configs.GetEventLimitGiftDataArray() {
		switch data.Condition {
		case 1: // 主城建设礼包
			check.PanicNotTrue(c.homeBaseGifts[data.ConditionValue] == nil, "%s %s [条件]:%d级 主城建设礼包的条件值已经重复", filename, data.Name, data.ConditionValue)
			c.homeBaseGifts[data.ConditionValue] = data
		case 2: // 武将培养礼包
			check.PanicNotTrue(c.captainGifts[data.ConditionValue] == nil, "%s %s [条件]:%d品质 武将培养礼包的条件值已经重复", filename, data.Name, data.ConditionValue)
			c.captainGifts[data.ConditionValue] = data
		case 3: // 爬塔勇士礼包
			check.PanicNotTrue(c.towerGifts[data.ConditionValue] == nil, "%s %s [条件]:%d楼 爬塔勇士礼包的条件值已经重复", filename, data.Name, data.ConditionValue)
			c.towerGifts[data.ConditionValue] = data
		case 4: // 君主等级礼包
			check.PanicNotTrue(c.heroLevelGifts[data.ConditionValue] == nil, "%s %s [条件]:%d级 君主等级礼包的条件值已经重复", filename, data.Name, data.ConditionValue)
			c.heroLevelGifts[data.ConditionValue] = data
		case 5: // 重建礼包
			check.PanicNotTrue(c.rebornGift == nil, "%s %s 重建礼包重复", filename, data.Name)
			c.rebornGift = data
		case 6: // 前线补给礼包
			check.PanicNotTrue(c.supplyGift == nil, "%s %s 前线补给礼包重复", filename, data.Name)
			c.supplyGift = data
		case 7: // 爬楼协助礼包
			check.PanicNotTrue(c.towerHelpGift == nil, "%s %s 爬楼协助礼包重复", filename, data.Name)
			c.towerHelpGift = data
		default:
			logrus.Panicf("%s 未知的事件时限礼包条件：%d", filename, data.Condition)
		}
	}
}

// 超值时限事件礼包
//gogen:config
type EventLimitGiftData struct {
	_ struct{} `file:"福利/事件时限礼包.txt"`
	_ struct{} `protogen:"true"`

	Id             uint64                                       // 礼包id
	Name           string                                       // 礼包名
	Icon           string                                       // 礼包icon
	Image          string                                       // 礼包插画
	Desc           string                                       // 礼包文本
	YuanbaoPrice   uint64                                       // 现价(元宝)，不走cost表了，确定是元宝，否则old_price也无法展示
	OldPrice       uint64                                       // 原价（展示用）
	SignIcon       string                                       // 显示的标签（文件路径）
	SignName       string                                       // 标签内显示的文本
	DiscountIcon   string                                       // 有红热和银亮等折扣标签
	Discount       string                                       // 标签内折扣值
	Dianquan       uint64                                       // 购买后立即获得的点券
	Priority       uint64        `validator:"uint" default:"0"` // 贴脸优先级，取最高值贴脸，如果是0则不贴脸
	TimeDuration   time.Duration `default:"1h" protofield:"-"`
	MinHeroLevel   uint64        `validator:"uint" default:"0" protofield:"-"` // 推送所需英雄等级下限
	MaxHeroLevel   uint64        `validator:"uint" default:"0" protofield:"-"` // 推送所需英雄等级上限
	MinGuanfuLevel uint64        `validator:"uint" default:"0" protofield:"-"` // 推送所需官府等级上限
	MaxGuanfuLevel uint64        `validator:"uint" default:"0" protofield:"-"` // 推送所需官府等级上限
	BuyLimit       uint64        `validator:"uint" default:"0"`                // 该种礼包购买次限（不是一次性购买数量，而是购买次数，每次只能买1个），填0不作限制

	// 以下字段仅服务器用
	Prize          *resdata.Prize   `protofield:"-"`
	Condition      uint64           `protofield:"-"`
	ConditionValue uint64           `validator:"uint" default:"0" protofield:"-"`

	ShowPrize         *resdata.Prize `desc:"展示奖励"`
	GuildEventPrizeId uint64         `desc:"联盟礼包Id" default:"0" validator:"uint"`
}

// 礼包是否可以推送(购买)
func (d *EventLimitGiftData) IsCanBuy(heroLevel uint64, guanfuLevel uint64) bool {
	return heroLevel >= d.MinHeroLevel && heroLevel <= d.MaxHeroLevel && guanfuLevel >= d.MinGuanfuLevel && guanfuLevel <= d.MaxGuanfuLevel
}

func (d *EventLimitGiftData) Init(filepath string, configs interface {
	GetGuanFuLevelDataArray() []*domestic_data.GuanFuLevelData
}) {
	heroMaxLevel := uint64(100)
	check.PanicNotTrue(d.MaxHeroLevel <= heroMaxLevel, "%s %s [max_hero_level]:%d 超出君主等级上限", filepath, d.Name, d.MaxHeroLevel)
	check.PanicNotTrue(d.MinHeroLevel <= heroMaxLevel, "%s %s [min_hero_level]:%d 超出君主等级上限", filepath, d.Name, d.MinHeroLevel)
	if d.MaxHeroLevel <= 0 {
		d.MaxHeroLevel = heroMaxLevel
	}
	check.PanicNotTrue(d.MinHeroLevel <= d.MaxHeroLevel, "%s %s [min_hero_level]:%d 不允许大于 [max_hero_level]:%d", filepath, d.Name, d.MinHeroLevel, d.MaxHeroLevel)

	guanfuLevelArr := configs.GetGuanFuLevelDataArray()
	guanfuMaxLevel := guanfuLevelArr[len(guanfuLevelArr)-1].Level
	check.PanicNotTrue(d.MaxGuanfuLevel <= guanfuMaxLevel, "%s %s [max_guanfu_level]:%d 超出官府等级上限", filepath, d.Name, d.MaxGuanfuLevel)
	check.PanicNotTrue(d.MinGuanfuLevel <= guanfuMaxLevel, "%s %s [min_guanfu_level]:%d 超出官府等级上限", filepath, d.Name, d.MinGuanfuLevel)
	if d.MaxGuanfuLevel <= 0 {
		d.MaxGuanfuLevel = guanfuMaxLevel
	}
	check.PanicNotTrue(d.MinGuanfuLevel <= d.MaxGuanfuLevel, "%s %s [min_guanfu_level]:%d 不允许大于 [max_guanfu_level]:%d", filepath, d.Name, d.MinGuanfuLevel, d.MaxGuanfuLevel)
}
