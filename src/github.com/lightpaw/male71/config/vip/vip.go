package vip

import (
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/check"
	"time"
)

//gogen:config
type VipMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"vip/vip杂项.txt"`
	_ struct{} `protogen:"true"`

	CollectVipDailyExpMinHeroLevel uint64          `desc:"可领 vip 每日经验的最小君主等级"`
	DungeonTimesCost               []*resdata.Cost `desc:"一次购买推图次数花费"`
	DungeonTimesEachBuy            []uint64        `desc:"一次购买获得的推图次数" validator:"uint,duplicate"`
}

func (d *VipMiscData) VipDungeonTimesCost(times uint64) (cost *resdata.Cost, toAdd uint64) {
	if times >= uint64(len(d.DungeonTimesCost)) {
		cost = d.DungeonTimesCost[len(d.DungeonTimesCost)-1]
	} else {
		cost = d.DungeonTimesCost[times]
	}

	if times >= uint64(len(d.DungeonTimesEachBuy)) {
		toAdd = d.DungeonTimesEachBuy[len(d.DungeonTimesEachBuy)-1]
	} else {
		toAdd = d.DungeonTimesEachBuy[times]
	}

	return
}

//gogen:config
type VipLevelData struct {
	_ struct{} `file:"vip/vip等级.txt"`
	_ struct{} `protogen:"true"`

	Level              uint64         `desc:"等级。0级表示还不是VIP" key:"true" validator:"uint"`
	Name               string         `desc:"名称"`
	Icon               *icon.Icon     `protofield:"IconId,%s.Id,string" desc:"图标"`
	UpgradeExp         uint64         `desc:"升到下一级需要的经验"`
	DailyExp           uint64         `desc:"每日赠送 vip 经验" validator:"uint"`
	DailyPrize         *resdata.Prize `desc:"vip 每日免费礼包" default:"nullable"`
	ShowDailyPrizeCost *resdata.Cost  `desc:"vip 每日免费礼包展示价格" default:"nullable"`
	LevelPrize         *resdata.Prize `desc:"vip 专属礼包" default:"nullable"`
	LevelPrizeCost     *resdata.Cost  `desc:"vip 专属礼包价格" default:"nullable"`
	ShowLevelPrizeCost *resdata.Cost  `desc:"vip 专属礼包展示价格" default:"nullable"`

	ContinuteDays *VipContinueDaysData `protofield:"-" head:"-"`

	NextLevelData *VipLevelData `protofield:"-" head:"-"`

	VipPrivilegeData `type:"sub"`
}

type VipPrivilegeData struct {
	BuyProsperity                    bool          `desc:"立即恢复繁荣度"`
	JiuGuanAutoMax                   bool          `desc:"酒馆一键白起"`
	JiuGuanCostRefreshCount          uint64        `desc:"新增酒馆花钱刷新次数" validator:"uint"`
	JiuGuanCostRefreshInfinite       bool          `desc:"酒馆无限次花钱刷新"`
	JiuGuanQuickConsult              bool          `desc:"酒馆立即请教"`
	CaptainTrainCoef                 float64       `desc:"提升修炼馆武将修炼经验" validator:"float64"`
	CaptainTrainCapacity             time.Duration `desc:"延长修炼馆武将修炼时间"`
	WallAutoFullSoldier              bool          `desc:"自动补兵"`
	BuySpMaxTimes                    uint64        `desc:"增加购买体力最大次数" validator:"uint"`
	DungeonMaxCostTimesLimit         uint64        `desc:"推图最大购买次数" validator:"uint"`
	WorkerUnlockPos                  uint64        `desc:"建筑队解锁pos" validator:"uint"`
	InvadeMultiLevelMonsterOnceCount uint64        `desc:"增加一次出征讨伐野怪次数" validator:"uint"`
	GuildPrizeOneKeyCollect          bool          `desc:"联盟礼包一键开启"`
	AddBlackMarketRefreshTimes       uint64        `desc:"增加游商最大刷新次数" validator:"uint"`
	FishingCaptainProbability        bool          `desc:"金杆钓"`
	ShowRegionHome                   bool          `desc:"城堡外显" default:"false"`
	ShowRegionSign                   bool          `desc:"城堡签名" default:"false"`
	ShowRegionTitle                  bool          `desc:"城堡铭牌" default:"false"`
	ShowHeadFrame                    bool          `desc:"君主头像框" default:"false"`
	ZhengWuAutoCompleted             bool          `desc:"政务自动完成" default:"false"`
	WorkshopAutoCompleted            bool          `desc:"铁匠铺自动完成" default:"false"`
	CollectDailySp                   bool          `desc:"可以领取日常体力" default:"false"`
}

func (data *VipLevelData) InitAll(filename string, conf interface {
	GetVipLevelDataArray() []*VipLevelData
}) {
	var preData *VipLevelData
	check.PanicNotTrue(len(conf.GetVipLevelDataArray()) > 0, "%v 必须配置数据。", filename)

	for _, d := range conf.GetVipLevelDataArray() {
		if preData == nil {
			check.PanicNotTrue(d.Level == 0, "%v, level 必须从0开始依次递增。%v", filename, d.Level)
		} else {
			check.PanicNotTrue(d.Level == preData.Level+1, "%v, level 必须从0开始依次递增。%v", filename, d.Level)
			check.PanicNotTrue(d.LevelPrize != nil, "%v, level:%v 必须配置专属礼包。%v", filename, d.Level)
			check.PanicNotTrue(d.DailyPrize != nil, "%v, level:%v 必须配置每日礼包。%v", filename, d.Level)
			check.PanicNotTrue(d.LevelPrizeCost != nil, "%v, level:%v 必须配置专属礼包消耗。%v", filename, d.Level)
			preData.NextLevelData = d
		}
		preData = d
	}
}

func (d *VipLevelData) Init(filename string, conf interface {
	HeroInitData() *heroinit.HeroInitData
}) {
	check.PanicNotTrue(d.UpgradeExp > 0, "%v, level:%v, upgrade_exp:%v 必须 > 0.", filename, d.Level, d.UpgradeExp)
	check.PanicNotTrue(d.WorkerUnlockPos < conf.HeroInitData().BuildingWorkerMaxCount, "%v, level:%v, worker_unlock_pos:%v 必须 < hero_init_data.txt building_worker_max_count:%v", filename, d.Level, d.WorkerUnlockPos, conf.HeroInitData().BuildingWorkerMaxCount)
	if d.JiuGuanCostRefreshInfinite {
		check.PanicNotTrue(d.JiuGuanCostRefreshCount == 0, "%v, level:%v 配置了酒馆无限刷新，就不能配置酒馆刷新次数：%v", filename, d.Level, d.JiuGuanCostRefreshCount)
	}
}

func (d *VipLevelData) ContinueDaysExp(days uint64) (exp uint64) {
	if dd := d.ContinuteDays; dd != nil {
		exp = dd.exps[days]
		return
	}
	return
}

//gogen:config
type VipContinueDaysData struct {
	_ struct{} `file:"vip/连续登录奖励.txt"`
	_ struct{} `protogen:"true"`

	Level uint64 `desc:"等级" key:"true" validator:"uint"`

	Days []uint64 `desc:"连续登录天数" validator:",duplicate"`
	Exp  []uint64 `desc:"连续登录VIP经验奖励，与 Days 一一对应" validator:"uint,duplicate"`
	exps map[uint64]uint64
}

func (d *VipContinueDaysData) Init(filename string, conf interface {
	GetVipLevelData(uint64) *VipLevelData
}) {
	levelData := conf.GetVipLevelData(d.Level)
	check.PanicNotTrue(levelData != nil, "%v, level:%v 不存在", filename, d.Level)
	check.PanicNotTrue(len(d.Days) == len(d.Exp), "%v, days 必须和 extra_exp 一一对应", filename)

	d.exps = make(map[uint64]uint64, len(d.Days))
	var preDay uint64
	for i, day := range d.Days {
		check.PanicNotTrue(day == preDay+1, "%v, days 必须依次递增。%v", filename, day)
		d.exps[day] = d.Exp[i]
		preDay = day
	}

	levelData.ContinuteDays = d
}
