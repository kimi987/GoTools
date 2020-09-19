package heroinit

import (
	"github.com/lightpaw/male7/config/body"
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/dungeon"
	"github.com/lightpaw/male7/config/function"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/head"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/military_data"
	"github.com/lightpaw/male7/config/pvetroop"
	"github.com/lightpaw/male7/config/settings"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/config/tag"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/male7/config/zhengwu"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/config/domestic_data/sub"
)

//gogen:config
type HeroCreateData struct {
	_     struct{} `singleton:"true"`
	_     struct{} `file:"杂项/创建英雄基础数据.txt"`
	Gold  uint64   `validator:"uint"`
	Food  uint64   `validator:"uint"`
	Wood  uint64   `validator:"uint"`
	Stone uint64   `validator:"uint"`

	NewSoldier uint64 `validator:"uint"`

	Prosperity uint64 `head:"-"`

	Captain []*captain.CaptainData
}

func (d *HeroCreateData) Init(filename string, configs interface {
	GetBuildingDataArray() []*domestic_data.BuildingData
	GetBuildingUnlockData(buildingType uint64) *domestic_data.BuildingUnlockData
	MainCityMiscData() *domestic_data.MainCityMiscData
}) {
	prosperity := uint64(0)
	for _, b := range configs.GetBuildingDataArray() {
		if b.Level == 1 && configs.GetBuildingUnlockData(uint64(b.Type)) == nil && configs.MainCityMiscData().IsMainCityBuildingType(b.Type) {
			prosperity += b.Prosperity
		}
	}

	d.Prosperity = prosperity
}

//gogen:config
type HeroInitData struct {
	_                      struct{}                        `singleton:"true"`
	_                      struct{}                        `file:"singleton/hero_init_data.txt"`
	BuildingWorkerMaxCount uint64                          `default:"2"`
	TechWorkerMaxCount     uint64                          `default:"1"`
	FirstLevelSoldierData  *domestic_data.SoldierLevelData `default:"1"`
	FirstLevelHeroData     *herodata.HeroLevelData         `default:"1"`
	FirstMainTask          *taskdata.MainTaskData          `default:"1"`
	FirstBaYeStageData     *taskdata.BaYeStageData         `default:"1"`
	FirstTitleData         *taskdata.TitleData             `default:"1"`

	FirstLevelCountdownPrize *domestic_data.CountdownPrizeData `default:"1"`

	CaptainIndexCount uint64 `default:"15"`

	Building []*domestic_data.BuildingData `head:"-"`

	// 最低等级的铁匠铺
	MinLevelTieJiangPuData *domestic_data.TieJiangPuLevelData `head:"-"`

	// 初始锻造的次数
	DefaultForgingTimes uint64 `head:"-"`

	MaxDepotEquipCapacity uint64 `head:"-"` // 背包装备的最大容量

	TempDepotExpireDuration time.Duration `head:"-"` // 临时背包的过期时间

	// 帮派捐献种类个数
	GuildDonateTypeCount uint64 `head:"-"`

	TroopCaptainCount uint64 `default:"5"`

	// 最大的标签数量
	MaxTagColorType uint64 `head:"-"`

	// 最大记录的标签日志数量
	MaxTagRecordCount uint64 `head:"-"`

	// 展示给查看的标签数量
	MaxShowForViewCount uint64 `head:"-"`

	StrategyRestoreDuration time.Duration `head:"-"`

	DungeonChallengeDefaultTimes    uint64        `head:"-"`
	DungeonChallengeMaxTimes        uint64        `head:"-"`
	DungeonChallengeRecoverDuration time.Duration `head:"-"`

	JunYingRecoveryDuration time.Duration `head:"-"` // 军营恢复间隔
	JunYingMaxTimes         uint64        `head:"-"` // 军营最大次数
	JunYingDefaultTimes     uint64        `head:"-"` // 军营默认次数

	MaxFavoritePosCount uint64 `head:"-"` // 最大收藏点的数量

	ActiveDegreeTaskDatas []*taskdata.ActiveDegreeTaskData `head:"-"`

	BwzlTaskDatas []*taskdata.BwzlTaskData `head:"-"`

	AchieveTaskDatas map[uint64]*taskdata.AchieveTaskData `head:"-"`

	FunctionOpenDataArray []*function.FunctionOpenData `head:"-"`

	DefaultHead *head.HeadData `head:"-"` // 默认头像

	DefaultBody *body.BodyData `head:"-"` // 默认形象

	PveTroopDatas []*pvetroop.PveTroopData `head:"-"`

	MultiLevelNpcInitTimes        uint64        `head:"-"` // 讨伐野怪初始次数
	MultiLevelNpcMaxTimes         uint64        `head:"-"` // 讨伐野怪次数上限
	MultiLevelNpcRecoveryDuration time.Duration `head:"-"` // 讨伐野怪次数恢复间隔

	InvaseHeroInitTimes        uint64        `head:"-"`
	InvaseHeroMaxTimes         uint64        `head:"-"`
	InvaseHeroRecoveryDuration time.Duration `head:"-"`

	JunTuanNpcInitTimes        uint64        `head:"-"`
	JunTuanNpcMaxTimes         uint64        `head:"-"`
	JunTuanNpcRecoveryDuration time.Duration `head:"-"`

	WorkshopOutputInitTimes        uint64        `head:"-"`
	WorkshopOutputMaxTimes         uint64        `head:"-"`
	WorkshopOutputRecoveryDuration time.Duration `head:"-"`

	ZhengWuMiscData *zhengwu.ZhengWuMiscData      `head:"-"` // 政务其他数据
	RandomZhengWu   func() []*zhengwu.ZhengWuData `head:"-"` // 政务随机

	DefaultSettings []shared_proto.SettingType `head:"-"` // 默认设置

	BuildingInitEffect *sub.BuildingEffectData `head:"-"` // 初始建筑加成效果

	BaowuLogLimit int `head:"-"`

	CaptainInitData *CaptainInitData `head:"-"`
}

func (d *HeroInitData) Init(filename string, configs interface {
	GetBuildingDataArray() []*domestic_data.BuildingData
	GetBuildingEffectData(int) *sub.BuildingEffectData
	GetBuildingUnlockData(buildingType uint64) *domestic_data.BuildingUnlockData
	GetTieJiangPuLevelDataArray() []*domestic_data.TieJiangPuLevelData
	MiscConfig() *singleton.MiscConfig
	JiuGuanMiscData() *military_data.JiuGuanMiscData
	JunYingMiscData() *military_data.JunYingMiscData
	GetJunYingLevelDataArray() []*military_data.JunYingLevelData
	TagMiscData() *tag.TagMiscData
	GetGuildDonateDataArray() []*guild_data.GuildDonateData
	DungeonMiscData() *dungeon.DungeonMiscData
	GetActiveDegreeTaskDataArray() []*taskdata.ActiveDegreeTaskData
	GetBwzlTaskDataArray() []*taskdata.BwzlTaskData
	GetAchieveTaskDataArray() []*taskdata.AchieveTaskData
	GetFunctionOpenDataArray() []*function.FunctionOpenData
	GetHeadDataArray() []*head.HeadData
	GetBodyDataArray() []*body.BodyData
	GetFunctionOpenData(intType uint64) *function.FunctionOpenData
	MainCityMiscData() *domestic_data.MainCityMiscData
	GetPveTroopDataArray() []*pvetroop.PveTroopData
	RegionConfig() *singleton.RegionConfig
	ZhengWuRandomData() *zhengwu.ZhengWuRandomData
	ZhengWuMiscData() *zhengwu.ZhengWuMiscData
	SettingMiscData() *settings.SettingMiscData
	MiscGenConfig() *singleton.MiscGenConfig
	RegionGenConfig() *singleton.RegionGenConfig
	GetCaptainAbilityDataArray() []*captain.CaptainAbilityData
	GetCaptainRebirthLevelDataArray() []*captain.CaptainRebirthLevelData
	GuildGenConfig() *singleton.GuildGenConfig
}) {
	for _, b := range configs.GetBuildingDataArray() {
		if b.Level == 1 && configs.GetBuildingUnlockData(uint64(b.Type)) == nil && configs.MainCityMiscData().IsMainCityBuildingType(b.Type) {
			d.Building = append(d.Building, b)
		}
	}

	d.MinLevelTieJiangPuData = configs.GetTieJiangPuLevelDataArray()[0]
	d.DefaultForgingTimes = configs.MiscConfig().DefaultForgingTimes

	d.MaxDepotEquipCapacity = configs.MiscConfig().MaxDepotEquipCapacity
	d.TempDepotExpireDuration = configs.MiscConfig().TempDepotExpireDuration

	typeCount := uint64(0)
	for _, b := range configs.GetGuildDonateDataArray() {
		if b.Times == 1 {
			typeCount++
		}
	}
	d.GuildDonateTypeCount = typeCount

	d.DungeonChallengeDefaultTimes = configs.DungeonMiscData().DefaultAutoTimes
	d.DungeonChallengeMaxTimes = configs.DungeonMiscData().MaxAutoTimes
	d.DungeonChallengeRecoverDuration = configs.DungeonMiscData().RecoverAutoDuration

	d.MaxTagColorType = configs.TagMiscData().MaxTagColorType
	d.MaxTagRecordCount = configs.TagMiscData().MaxRecordCount
	d.MaxShowForViewCount = configs.TagMiscData().MaxShowForViewCount

	d.StrategyRestoreDuration = configs.MiscConfig().StrategyRestoreDuration

	junYingMinLevelData := configs.GetJunYingLevelDataArray()[0]
	d.JunYingRecoveryDuration = junYingMinLevelData.RecoveryDuration
	d.JunYingMaxTimes = junYingMinLevelData.MaxTimes
	d.JunYingDefaultTimes = configs.JunYingMiscData().DefaultTimes

	d.MaxFavoritePosCount = configs.MiscConfig().MaxFavoritePosCount

	d.ActiveDegreeTaskDatas = configs.GetActiveDegreeTaskDataArray()
	d.BwzlTaskDatas = configs.GetBwzlTaskDataArray()

	d.AchieveTaskDatas = map[uint64]*taskdata.AchieveTaskData{}
	for _, data := range configs.GetAchieveTaskDataArray() {
		if data.PrevTask == nil {
			d.AchieveTaskDatas[data.AchieveType] = data
		}
	}

	d.FunctionOpenDataArray = configs.GetFunctionOpenDataArray()

	for _, head := range configs.GetHeadDataArray() {
		if head.DefaultHead {
			d.DefaultHead = head
			break
		}
	}
	check.PanicNotTrue(d.DefaultHead != nil, "头像.txt 没有配置默认头像")

	for _, body := range configs.GetBodyDataArray() {
		if body.DefaultBody {
			d.DefaultBody = body
			break
		}
	}
	check.PanicNotTrue(d.DefaultBody != nil, "形象.txt 没有配置默认形象")

	// 不再处理这个
	//for fType := range shared_proto.FunctionType_name {
	//	if shared_proto.FunctionType(fType) != shared_proto.FunctionType_Type_Invalid {
	//		if configs.GetFunctionOpenData(uint64(fType)) == nil {
	//			d.DefaultOpenFunctionTypes = append(d.DefaultOpenFunctionTypes, shared_proto.FunctionType(fType))
	//		}
	//	}
	//}

	d.PveTroopDatas = configs.GetPveTroopDataArray()

	d.MultiLevelNpcInitTimes = configs.RegionConfig().MultiLevelNpcInitTimes
	d.MultiLevelNpcMaxTimes = configs.RegionConfig().MultiLevelNpcMaxTimes
	d.MultiLevelNpcRecoveryDuration = configs.RegionConfig().MultiLevelNpcRecoveryDuration
	check.PanicNotTrue(d.MultiLevelNpcRecoveryDuration > 0, "讨伐野怪次数恢复间隔必须>0")

	d.InvaseHeroInitTimes = configs.RegionGenConfig().InvaseHeroInitTimes
	d.InvaseHeroMaxTimes = configs.RegionGenConfig().InvaseHeroMaxTimes
	d.InvaseHeroRecoveryDuration = configs.RegionGenConfig().InvaseHeroRecoveryDuration
	check.PanicNotTrue(d.InvaseHeroRecoveryDuration > 0, "讨伐玩家次数恢复间隔必须>0")

	d.JunTuanNpcInitTimes = configs.RegionGenConfig().JunTuanNpcInitTimes
	d.JunTuanNpcMaxTimes = configs.RegionGenConfig().JunTuanNpcMaxTimes
	d.JunTuanNpcRecoveryDuration = configs.RegionGenConfig().JunTuanNpcRecoveryDuration
	check.PanicNotTrue(d.JunTuanNpcRecoveryDuration > 0, "讨伐军团怪次数恢复间隔必须>0")

	d.WorkshopOutputInitTimes = configs.GuildGenConfig().WorkshopOutputInitTimes
	d.WorkshopOutputMaxTimes = configs.GuildGenConfig().WorkshopOutputMaxTimes
	d.WorkshopOutputRecoveryDuration = configs.GuildGenConfig().WorkshopOutputRecoveryDuration
	check.PanicNotTrue(d.WorkshopOutputRecoveryDuration > 0, "联盟工坊玩家生成次数恢复间隔必须>0")

	d.RandomZhengWu = func() []*zhengwu.ZhengWuData {
		return configs.ZhengWuRandomData().Random(configs.ZhengWuMiscData().RandomCount, false)
	}
	d.ZhengWuMiscData = configs.ZhengWuMiscData()

	d.DefaultSettings = configs.SettingMiscData().DefaultSettings

	d.BuildingInitEffect = configs.GetBuildingEffectData(u64.Int(configs.MiscConfig().BuildingInitEffect))

	check.PanicNotTrue(d.BuildingInitEffect != nil, "初始建筑加成效果没有配置。建筑效果.txt id:%v", configs.MiscConfig().BuildingInitEffect)

	d.BaowuLogLimit = configs.MiscGenConfig().BaowuLogLimit
	check.PanicNotTrue(d.BaowuLogLimit > 0, "宝物日志条数必须 > 0")

	cid := &CaptainInitData{}
	cid.AbilityData = configs.GetCaptainAbilityDataArray()[0]
	cid.RebirthData = configs.GetCaptainRebirthLevelDataArray()[0]
	cid.TroopCaptainCount = d.TroopCaptainCount

	d.CaptainInitData = cid
}

type CaptainInitData struct {
	AbilityData *captain.CaptainAbilityData

	RebirthData *captain.CaptainRebirthLevelData

	TroopCaptainCount uint64 // 部队中武将的数量
}
