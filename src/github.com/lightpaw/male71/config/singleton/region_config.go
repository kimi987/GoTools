package singleton

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/logrus"
)

//gogen:config
type RegionGenConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"地图/地区杂项.txt"`
	_ struct{} `protogen:"true"`

	UseGoodsMianMaxDuraion time.Duration `default:"30m"` // 破除免战后，最大的免战物品CD

	InvaseHeroInitTimes        uint64        `default:"3" protofield:"-"`
	InvaseHeroMaxTimes         uint64        `default:"5"`  // 讨伐玩家次数上限
	InvaseHeroRecoveryDuration time.Duration `default:"3h"` // 讨伐玩家次数恢复间隔

	JunTuanNpcInitTimes        uint64        `default:"3" protofield:"-"`
	JunTuanNpcMaxTimes         uint64        `default:"5"`  // 讨伐军团怪次数上限
	JunTuanNpcRecoveryDuration time.Duration `default:"3h"` // 讨伐军团怪次数恢复间隔
	JunTuanWinTimeLimit        uint64        `default:"3"`  // 讨伐军团怪进攻方连胜离场场次
}

//gogen:config
type RegionConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"地图/地区杂项.txt"`
	_ struct{} `proto:"shared_proto.RegionConfigProto"`
	_ struct{} `protoconfig:"RegionConfig"`

	SlowMoveDuration time.Duration `default:"3h" parser:"time.ParseDuration" protofield:"-"`   // 缓慢迁移城市cd
	ExpelDuration    time.Duration `default:"180s" parser:"time.ParseDuration" protofield:"-"` // 驱逐cd

	BasicTroopMoveVelocityPerSecond      float64       `default:"0.25"`                                                   // 基础移动速度
	BasicTroopMoveToNpcVelocityPerSecond float64       `default:"0.25" protofield:"-"`                                    // 基础攻击npc的移动速度
	MinTroopMoveVelocityPerSecond        float64       `head:"-,%s.BasicTroopMoveVelocityPerSecond * 0.5" protofield:"-"` // 基础移动速度
	Edge                                 float64       `default:"1"`                                                      // 六边形边长
	TroopMoveOffsetDuration              time.Duration `default:"0s" protofield:"-"`                                      // 部队移动偏移时间

	EdgeNotHomeLen uint64 `validator:"uint" default:"6"`

	AttackerWoundedRate float64 `default:"0.5" protofield:"-"` // 进攻方伤兵系数
	DefenserWoundedRate float64 `default:"0.7" protofield:"-"` // 防守方伤兵系数
	AssisterWoundedRate float64 `default:"0.5" protofield:"-"` // 援助方伤兵系数

	RobberCoef float64 `default:"1" protofield:"-"` // 抢的系数

	MaxLostProsperity  uint64  `default:"8000" protofield:"-"` // 每日最大损失繁荣度
	LostProsperityCoef float64 `default:"0.3" protofield:"-"`  // 损失系数

	RobTickDuration          time.Duration `default:"6s" parser:"time.ParseDuration" protofield:"-"`  // 抢多久计算一次结算一次
	RobMaxDuration           time.Duration `default:"30m" parser:"time.ParseDuration" protofield:"-"` // 最多抢劫多久
	ReduceProsperityDuration time.Duration `default:"10m" parser:"time.ParseDuration" protofield:"-"` // 抢多久计算扣一次繁荣度
	AssisterTickDuration     time.Duration `default:"5m" parser:"time.ParseDuration" protofield:"-"`  // 援助tick时间间隔
	AssisterMaxDuration      time.Duration `default:"24h" parser:"time.ParseDuration" protofield:"-"` // 援助最长时间间隔
	RobBaowuTickDuration     time.Duration `default:"10m" parser:"time.ParseDuration" protofield:"-"` // 抢多久计算抢一次宝物

	MaxRobbers uint64 `default:"4" protofield:"-"` // 最大抢劫部队数
	MaxAssist  uint64 `default:"4" protofield:"-"` // 最大驻扎部队数

	MaxInvationTroops  uint64 `default:"3" protofield:"-"` // 出征部队最大上限
	MaxInvationCaptain uint64 `default:"5" protofield:"-"` // 出征部队武将最大上限

	TentBuildingDuration time.Duration `default:"0m"`                  // 行营建造所需时间
	TentRestoreDuration  time.Duration `default:"100m" protofield:"-"` // 行营恢复满血所需时间
	TentFreeMoveCooldown time.Duration `default:"5m" protofield:"-"`   // 行营免费迁移CD

	TentHomeRegionEnterDuration    time.Duration `default:"6h" protofield:"-"` // 行营被打出来，再次进入主城地区所需时间
	TentMonsterRegionEnterDuration time.Duration `default:"6h" protofield:"-"` // 行营被打出来，再次进入荣誉地区所需时间

	MultiLevelNpcInitTimes        uint64        `default:"3" protofield:"-"`
	MultiLevelNpcMaxTimes         uint64        `default:"5"`  // 讨伐野怪次数上限
	MultiLevelNpcRecoveryDuration time.Duration `default:"3h"` // 讨伐野怪次数恢复间隔

	// 繁荣度损失速度
	NpcRobLostProsperityPerDuration uint64 `default:"1" protofield:"-"` // Npc持续掠夺玩家，每个duration减少的繁荣度

	MiaoTentBuildingDuration time.Duration // 秒行营建造间隔
	MiaoTentBuildingCost     *resdata.Cost // 秒行营建造消耗

	// 战斗场景
	CombatScene *scene.CombatScene `protofield:"-"`

	AstDefendRestoreHomeProsperityAmount *data.Amount  `default:"8"` // 每个间隔盟友驻扎给回复的繁荣度
	AstDefendRestoreHomeProsperity       uint64        `default:"8"` // 废弃，每个间隔盟友驻扎给回复的繁荣度
	RestoreHomeProsperity                uint64        `default:"10"`
	RestoreHomeProsperityDuration        time.Duration `default:"1m"`

	AstDefendLogLimit uint64 `default:"100" protofield:"-"`

	RestoreTentProsperity         uint64        `default:"10" protofield:"-"`
	RestoreTentProsperityDuration time.Duration `default:"1m" protofield:"-"`

	LossTentProsperity         uint64        `default:"10" protofield:"-"`
	LossTentProsperityDuration time.Duration `default:"1m" protofield:"-"`

	WhiteFlagDuration time.Duration `default:"5m"` // 插白旗时间

	RuinsBaseExpireDuration time.Duration `default:"24h" protofield:"-"` // 废墟过期时间

	// 联盟主城地区
	GuildRegionCenterX uint64   `default:"60"`
	GuildRegionCenterY uint64   `default:"52"`
	GuildRegionRadius  []uint64 `default:"10,30,50" protofield:"-"`

	// 侦查
	InvestigateCd               time.Duration `default:"5m"`
	MiaoInvestigateCdCost       *resdata.Cost
	InvestigateCost             *resdata.Cost
	InvestigateSpeedup          float64       `default:"7"`
	InvestigateMailTimeout      time.Duration `default:"30m"`
	InvestigateMaxDistance      uint64        `default:"300"`
	InvestigationLimit          uint64        `default:"10" protofield:"-"`
	InvestigationExpireDuration time.Duration `default:"30m" protofield:"-"`
	InvestigationBaowuCount     uint64        `default:"3" protofield:"-"`

	// 免战
	NewHeroMianDuration        time.Duration `default:"72h" protofield:"-"` // 新手免战时间
	NewHeroRemoveMianBaseLevel uint64        `default:"6"`                  // 移除新手免战的等级
	RebornMianDuration         time.Duration `default:"24h"`                // 流亡重建免战时间

	MinViewXLen int `default:"40" protofield:"-"`
	MinViewYLen int `default:"30" protofield:"-"`
	MaxViewXLen int `default:"80" protofield:"-"`
	MaxViewYLen int `default:"60" protofield:"-"`

	ListEnemyPosCount int `default:"30" protofield:"-"`

	SearchBaozNpcCount int `default:"30" protofield:"-"`

	KeepBaozMaxDistance int           `default:"300" protofield:"-"`
	KeepBaozDuration    time.Duration `default:"12h" protofield:"-"`

	//侦察部队的ID
	InvestigateTroopId uint64 `head:"-" protofield:"-"`

	// 根据位置偏移获取到layout
	evenOffsetLayoutMap map[cb.Cube]*domestic_data.BuildingLayoutData

	// 根据等级找到layout（只包含当前等级）
	evenOffsetCubeMapOnlyCurrentLevel map[uint64][]cb.Cube

	// 根据等级找到layout（包含下面的等级）
	evenOffsetCubeMapIncludeLowLevel [][]cb.Cube
}

func (c *RegionConfig) MoveDurationCalc(distance float64) time.Duration {
	return MoveDuration(distance, c.BasicTroopMoveVelocityPerSecond)
}

func (c *RegionConfig) MoveToNpcDurationCalc(distance float64) time.Duration {
	return MoveDuration(distance, c.BasicTroopMoveToNpcVelocityPerSecond)
}

func MoveDuration(distance, moveSpeed float64) time.Duration {
	if moveSpeed <= 0 {
		return 0
	}

	return time.Duration(distance/moveSpeed) * time.Second
}

func (c *RegionConfig) Init(filename string, configs interface {
	GetHeroLevelSubDataArray() []*data.HeroLevelSubData
	GetBuildingLayoutDataArray() []*domestic_data.BuildingLayoutData
	GetGuanFuLevelDataArray() []*domestic_data.GuanFuLevelData
}) {

	heroAry := configs.GetHeroLevelSubDataArray()

	if len(heroAry) > 0 {
		c.InvestigateTroopId = heroAry[len(heroAry)-1].TroopsCount
		logrus.Debugf("初始化侦察部队Id 成功 InvestigateTroopId[%d]", c.InvestigateTroopId)
	}

	var originBaseOffset = hexagon.SpiralRing(0, 0, 1)

	c.evenOffsetLayoutMap = make(map[cb.Cube]*domestic_data.BuildingLayoutData)
	for _, v := range configs.GetBuildingLayoutDataArray() {
		cube := cb.XYCube(v.RegionOffsetX, v.RegionOffsetY)
		_, exist := c.evenOffsetLayoutMap[cube]
		check.PanicNotTrue(!exist, "野外资源点偏移配置中，配置了重复的偏移位置，重复位置：%d, %d", v.RegionOffsetX, v.RegionOffsetY)

		//check.PanicNotTrue(!cb.Contains(originBaseOffset, cube), "野外资源点偏移配置中，配置的偏移位置不能是原点和原点周围一圈，无效位置：%d, %d", v.RegionOffsetX, v.RegionOffsetY)

		c.evenOffsetLayoutMap[cube] = v
	}

	c.evenOffsetCubeMapOnlyCurrentLevel = make(map[uint64][]cb.Cube)
	maxBaseLevel := uint64(0)
	for _, v := range configs.GetBuildingLayoutDataArray() {
		maxBaseLevel = u64.Max(maxBaseLevel, v.RequireBaseLevel)

		check.PanicNotTrue(v.RequireBaseLevel >= 0, "野外资源点配置了0级主城的势力范围，layoutId: %d", v.Id)

		cubes := c.evenOffsetCubeMapOnlyCurrentLevel[v.RequireBaseLevel]
		cubes = append(cubes, cb.XYCube(v.RegionOffsetX, v.RegionOffsetY))

		c.evenOffsetCubeMapOnlyCurrentLevel[v.RequireBaseLevel] = cubes
	}

	// 1级包含主城坐标以及周围一圈的坐标
	layouts := c.evenOffsetCubeMapOnlyCurrentLevel[1]
	for _, v := range originBaseOffset {
		if !cb.Contains(layouts, v) {
			layouts = append(layouts, v)
		}
	}
	c.evenOffsetCubeMapOnlyCurrentLevel[1] = layouts

	// 每次将前一级的数据全加上
	var prevLevelCubes []cb.Cube
	for level := uint64(1); level <= maxBaseLevel; level++ {

		var array []cb.Cube
		array = append(array, prevLevelCubes...)

		currentLevel := c.evenOffsetCubeMapOnlyCurrentLevel[level]
		array = append(array, currentLevel...)

		c.evenOffsetCubeMapIncludeLowLevel = append(c.evenOffsetCubeMapIncludeLowLevel, array)

		prevLevelCubes = array
	}

	for _, data := range configs.GetGuanFuLevelDataArray() {
		check.PanicNotTrue(data.MoveBaseRestoreHomeProsperity >= c.RestoreHomeProsperity, "%s 官府配置的迁城buf期间每次恢复的繁荣度[%d]必须比默认恢复数量[%d]要多", filename, data.MoveBaseRestoreHomeProsperity, c.RestoreHomeProsperity)
	}

}

func (c *RegionConfig) GetLayoutDataByEvenOffset(offset cb.Cube) *domestic_data.BuildingLayoutData {
	return c.evenOffsetLayoutMap[offset]
}

func (c *RegionConfig) GetEvenOffsetCubesOnlyCurrentLevel(level uint64) []cb.Cube {
	return c.evenOffsetCubeMapOnlyCurrentLevel[level]
}

func (c *RegionConfig) GetEvenOffsetCubesIncludeLowLevel(level uint64) []cb.Cube {
	if level <= 0 {
		return nil
	}

	index := level - 1
	if n := uint64(len(c.evenOffsetCubeMapIncludeLowLevel)); index >= n {
		index = n - 1
	}

	return c.evenOffsetCubeMapIncludeLowLevel[index]
}

func (c *RegionConfig) GetMaxLostProsperity(prosperityCapcity uint64) uint64 {
	return u64.Max(c.MaxLostProsperity, u64.MultiCoef(prosperityCapcity, c.LostProsperityCoef))
}

func (c *RegionConfig) GetTentValidTime(ctime time.Time) time.Time {
	if c.TentBuildingDuration > 0 {
		return ctime.Add(c.TentBuildingDuration)
	}

	return time.Time{}
}
