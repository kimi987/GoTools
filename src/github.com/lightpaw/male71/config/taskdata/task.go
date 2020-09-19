package taskdata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/dungeon"
	"github.com/lightpaw/male7/config/military_data"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/shop"
	"github.com/lightpaw/male7/config/towerdata"
	"github.com/lightpaw/male7/gen/pb/task"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"strconv"
	"strings"
	"github.com/lightpaw/male7/config/zhanjiang"
	"github.com/lightpaw/male7/config/xiongnu"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/monsterdata"
	"time"
	"github.com/lightpaw/male7/config/captain"
)

//gogen:config
type TaskMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"任务/杂项.txt"`
	_ struct{} `proto:"shared_proto.TaskMiscDataProto"`
	_ struct{} `protoconfig:"TaskMiscData"`

	MaxShowAchieveCount uint64 `default:"3"` // 最大展示成就数量
	BwzlBgImg           string               // 霸王之路背景图

	TaskMonsterArriveOffset time.Duration `default:"3s" protofield:"-"`
}

func (*TaskMiscData) InitAll(configs interface {
	GetMainTaskDataArray() []*MainTaskData
	GetBranchTaskDataArray() []*BranchTaskData
	GetAchieveTaskDataArray() []*AchieveTaskData
	GetBaYeTaskDataArray() []*BaYeTaskData
	GetBwzlTaskDataArray() []*BwzlTaskData
	GetActiveDegreeTaskDataArray() []*ActiveDegreeTaskData
	GetActiveDegreePrizeDataArray() []*ActiveDegreePrizeData
	GetTitleTaskDataArray() []*TitleTaskData
	GetActivityTaskDataArray() []*ActivityTaskData
}) {

	idMap := make(map[uint64]struct{})
	for _, t := range configs.GetMainTaskDataArray() {
		_, ok := idMap[t.Sequence]
		check.PanicNotTrue(!ok, "任务Id重复，主线任务、支线任务、成就任务、活跃度任务、霸王之路任务和霸业任务的ID不允许重复，重复id: %v", t.Sequence)

		idMap[t.Sequence] = struct{}{}
	}

	for _, t := range configs.GetBranchTaskDataArray() {
		_, ok := idMap[t.Id]
		check.PanicNotTrue(!ok, "任务Id重复，主线任务、支线任务、成就任务、活跃度任务、霸王之路任务和霸业任务的ID不允许重复，重复id: %v", t.Id)

		idMap[t.Id] = struct{}{}
	}

	for _, t := range configs.GetAchieveTaskDataArray() {
		_, ok := idMap[t.Id]
		check.PanicNotTrue(!ok, "任务Id重复，主线任务、支线任务、成就任务、活跃度任务、霸王之路任务和霸业任务的ID不允许重复，重复id: %v", t.Id)

		idMap[t.Id] = struct{}{}
	}

	for _, t := range configs.GetActiveDegreeTaskDataArray() {
		_, ok := idMap[t.Id]
		check.PanicNotTrue(!ok, "任务Id重复，主线任务、支线任务、成就任务、活跃度任务、霸王之路任务和霸业任务的ID不允许重复，重复id: %v", t.Id)

		idMap[t.Id] = struct{}{}
	}

	for _, t := range configs.GetBwzlTaskDataArray() {
		_, ok := idMap[t.Id]
		check.PanicNotTrue(!ok, "任务Id重复，主线任务、支线任务、成就任务、活跃度任务、霸王之路任务和霸业任务的ID不允许重复，重复id: %v", t.Id)

		idMap[t.Id] = struct{}{}
	}

	for _, t := range configs.GetBaYeTaskDataArray() {
		_, ok := idMap[t.Id]
		check.PanicNotTrue(!ok, "任务Id重复，主线任务、支线任务、成就任务、活跃度任务、霸王之路任务和霸业任务的ID不允许重复，重复id: %v", t.Id)

		idMap[t.Id] = struct{}{}
	}

	for _, t := range configs.GetTitleTaskDataArray() {
		_, ok := idMap[t.Id]
		check.PanicNotTrue(!ok, "任务Id重复，主线任务、支线任务、成就任务、活跃度任务、霸王之路任务和霸业任务的ID不允许重复，重复id: %v", t.Id)

		idMap[t.Id] = struct{}{}
	}

	for _, t := range configs.GetActivityTaskDataArray() {
		_, ok := idMap[t.Id]
		check.PanicNotTrue(!ok, "活动任务ID不允许和主线任务、支线任务、成就任务、活跃度任务、霸王之路任务和霸业任务的ID重复，重复id: %v", t.Id)
		idMap[t.Id] = struct{}{}
	}

	// 检查主线任务里面配置的成就任务不存在相同的成就类型配置多处

	totalDegree := uint64(0)
	for _, d := range configs.GetActiveDegreeTaskDataArray() {
		totalDegree += d.AddDegree * d.Target.TotalProgress
	}

	for _, prizeData := range configs.GetActiveDegreePrizeDataArray() {
		check.PanicNotTrue(prizeData.Degree <= totalDegree, "配置的所有活跃度任务奖励中任务总的活跃度[%d]不够领取 %d 活跃度奖励", totalDegree, prizeData.Degree)
	}
}

//gogen:config
type MainTaskData struct {
	_ struct{} `file:"任务/主线任务.txt"`
	_ struct{} `proto:"shared_proto.TaskDataProto"`
	_ struct{} `protoconfig:"MainTask"`

	// 主线
	Sequence uint64 `key:"true" protofield:"Id"`

	// 任务名字
	Name string

	// 任务内容
	Text string

	// 任务目标
	Target *TaskTargetData

	// 任务奖励
	Prize *resdata.Prize

	BranchTask []*BranchTaskData `default:"nullable" protofield:"-"`

	RemoveTaskMsg pbutil.Buffer `head:"-" protofield:"-"`

	nextTask *MainTaskData
}

func (d *MainTaskData) NextTask() *MainTaskData {
	return d.nextTask
}

func (d *MainTaskData) Init(filename string, dataMap map[uint64]*MainTaskData) {
	d.RemoveTaskMsg = task.NewS2cRemoveTaskMsg(u64.Int32(d.Sequence)).Static()

	if d.Sequence > 1 {
		prev := dataMap[d.Sequence-1]
		check.PanicNotTrue(prev != nil, "%s 主线任务的sequence[%v]没找到，必须从1开始连续配置", filename, d.Sequence-1)
		prev.nextTask = d
	}
}

//gogen:config
type BranchTaskData struct {
	_ struct{} `file:"任务/支线任务.txt"`
	_ struct{} `proto:"shared_proto.TaskDataProto"`
	_ struct{} `protoconfig:"BranchTask"`

	// 任务id
	Id uint64

	Next uint64 `validator:"uint" protofield:"-"`

	// 任务名字
	Name string

	// 任务内容
	Text string

	// 任务目标
	Target *TaskTargetData

	// 任务奖励
	Prize *resdata.Prize

	nextTask *BranchTaskData

	RemoveTaskMsg pbutil.Buffer `head:"-" protofield:"-"`
}

func (d *BranchTaskData) NextTask() *BranchTaskData {
	return d.nextTask
}

func (d *BranchTaskData) Init(filename string, dataMap map[uint64]*BranchTaskData) {
	d.RemoveTaskMsg = task.NewS2cRemoveTaskMsg(u64.Int32(d.Id)).Static()

	if d.Next > 0 {
		check.PanicNotTrue(d.Id < d.Next, "%s 支线任务[%v]配置的后续任务Next[%v]，必须配置比自己的ID大的ID", filename, d.Id, d.Next)

		d.nextTask = dataMap[d.Next]
		check.PanicNotTrue(d.nextTask != nil, "%s 支线任务[%v]配置的后续任务Next[%v]，不存在", filename, d.Id, d.Next)
	}
}

//gogen:config
type AchieveTaskData struct {
	_ struct{} `file:"任务/成就任务.txt"`
	_ struct{} `proto:"shared_proto.TaskDataProto"`
	_ struct{} `protoconfig:"AchieveTask"`

	// 任务id
	Id uint64

	// 任务名字
	Name string

	// 任务内容
	Text string

	// 图标
	Icon string

	// 任务目标
	Target *TaskTargetData

	// 任务奖励
	Prize *resdata.Prize

	AchieveType uint64 `validator:"uint"` // 成就类型

	Next uint64 `validator:"uint" protofield:"-"`

	Star      uint64 // 该任务的星数
	TotalStar uint64 // 完成该任务时，用来计算星数的值

	Quality shared_proto.Quality

	Order uint64 `validator:"uint" default:"0"`

	PrevTask      *AchieveTaskData `head:"-" protofield:",config.U64ToI32(%s.Id)"`
	nextTask      *AchieveTaskData
	RemoveTaskMsg pbutil.Buffer    `head:"-" protofield:"-"`
}

func (d *AchieveTaskData) Init(filename string, dataMap map[uint64]*AchieveTaskData) {
	d.RemoveTaskMsg = task.NewS2cRemoveTaskMsg(u64.Int32(d.Id)).Static()

	if d.Next > 0 {
		check.PanicNotTrue(d.Id < d.Next, "%s 成就任务[%v]配置的后续任务Next[%v]，必须配置比自己的ID大的ID", filename, d.Id, d.Next)

		d.nextTask = dataMap[d.Next]
		check.PanicNotTrue(d.nextTask != nil, "%s 成就任务[%v]配置的后续任务Next[%v]，不存在", filename, d.Id, d.Next)
		d.nextTask.PrevTask = d

		check.PanicNotTrue(d.nextTask.AchieveType == d.AchieveType, "%s 成就任务[%v]配置的后续任务Next[%v]配置的achieve_type [%d] 不相同", filename, d.Id, d.Next, d.AchieveType)
	}
}

func (d *AchieveTaskData) NextTask() *AchieveTaskData {
	return d.nextTask
}

// 成就星数奖励
//gogen:config
type AchieveTaskStarPrizeData struct {
	_ struct{} `file:"任务/成就星数奖励.txt"`
	_ struct{} `proto:"shared_proto.AchieveTaskStarPrizeProto"`
	_ struct{} `protoconfig:"AchieveTaskStarPrize"`

	Star    uint64           `validator:"int>0" key:"true"`      // 星数
	Icon    string                                               // 奖励图标
	Desc    string                                               // 奖励描述
	Prize   *resdata.Prize                                       // 奖励
	Plunder *resdata.Plunder `default:"nullable" protofield:"-"` // 真正给的掉落
}

// 霸业阶段
//gogen:config
type BaYeStageData struct {
	_ struct{} `file:"任务/霸业任务阶段.txt"`
	_ struct{} `proto:"shared_proto.BaYeStageDataProto"`
	_ struct{} `protoconfig:"BaYeStage"`

	Stage     uint64          `key:"true"` // 阶段
	StageName string          `default:" "`
	Name      string          `validator:"string>0"`                                                        // 阶段名字
	Tasks     []*BaYeTaskData `validator:"int" protofield:",config.U64a2I32a(GetBaYeTaskDataKeyArray(%s))"` // 任务目标
	Prize     *resdata.Prize                                                                                // 奖励
	Prev      *BaYeStageData  `head:"-" protofield:"-"`                                                     // 下一阶段，0表示没有下一阶段
	Next      *BaYeStageData  `head:"-" protofield:",config.U64ToI32(%s.Stage)"`                            // 下一阶段，0表示没有下一阶段
}

func (d *BaYeStageData) Init(filename string, dataMap map[uint64]*BaYeStageData) {
	if d.Stage > 1 {
		prev := dataMap[d.Stage-1]
		check.PanicNotTrue(prev != nil, "%s 霸业阶段的stage[%v]没找到，必须从1开始连续配置", filename, d.Stage-1)
		prev.Next = d
		d.Prev = prev
	}
}

// 霸业任务
//gogen:config
type BaYeTaskData struct {
	_ struct{} `file:"任务/霸业任务.txt"`
	_ struct{} `proto:"shared_proto.TaskDataProto"`
	_ struct{} `protoconfig:"BaYeTask"`

	// 任务id
	Id uint64

	// 任务名字
	Name string

	// 任务内容
	Text string

	// 图标
	Icon string

	// 任务目标
	Target *TaskTargetData

	// 任务奖励
	Prize *resdata.Prize
}

// 活跃度任务
//gogen:config
type ActiveDegreeTaskData struct {
	_ struct{} `file:"任务/活跃度任务.txt"`
	_ struct{} `proto:"shared_proto.TaskDataProto"`
	_ struct{} `protoconfig:"ActiveDegreeTask"`

	Id uint64 // 任务id

	// 任务名字
	Name string

	// 任务内容
	Text string

	// 图标
	Icon string

	// 任务目标
	Target *TaskTargetData

	AddDegree uint64 `validator:"int>0"` // 每一点进度增加的活跃度
}

func (d *ActiveDegreeTaskData) Init(filename string) {
	check.PanicNotTrue(d.Target.Daily, "%s 活跃度任务必须配置 daily: %d-%s", filename, d.Id, d.Name)
}

// 活跃度奖励
//gogen:config
type ActiveDegreePrizeData struct {
	_ struct{} `file:"任务/活跃度任务奖励.txt"`
	_ struct{} `proto:"shared_proto.ActiveDegreePrizeProto"`
	_ struct{} `protoconfig:"ActiveDegreePrize"`

	Degree  uint64           `validator:"int>0" key:"true"`      // 活跃度
	Prize   *resdata.Prize                                       // 奖励
	Plunder *resdata.Plunder `default:"nullable" protofield:"-"` // 真正给的掉落
}

//gogen:config
type TaskTargetData struct {
	_ struct{} `file:"任务/任务目标.txt"`
	_ struct{} `proto:"shared_proto.TaskTargetProto"`

	Id uint64 `protofield:"-"`

	// 任务目标类型
	Type shared_proto.TaskTargetType

	X uint64 `protofield:"-"`
	Y string `protofield:"-"`
	Z string `protofield:"-"`

	Daily bool `default:"false" protofield:"-"`

	SubType   uint64 `head:"-"`
	SubTypeId uint64 `head:"-"`

	// 不显示进度
	DontUpdateProgress bool `protofield:"-"`

	TotalProgress uint64 `head:"-" protofield:"-"`

	ShowProgress uint64 `head:"-" protofield:"TotalProgress"`

	//// 目标
	//
	//// 主城升到N级
	//BaseLevel uint64 `head:"-" protofield:"-"`
	//
	//// 君主到达X级
	//HeroLevel uint64 `head:"-" protofield:"-"`
	//
	//// 科技到达X级
	//TechId    uint64 `head:"-" protofield:"-"`
	TechGroup uint64 `head:"-" protofield:"-"`
	//TechLecel uint64 `head:"-" protofield:"-"`
	//
	//// 建筑到达X级
	BuildingType  shared_proto.BuildingType `head:"-" protofield:"-"`
	BuildingType2 shared_proto.BuildingType `head:"-" protofield:"-"`
	//BuildingLevel uint64                    `head:"-" protofield:"-"`
	//
	//UpgradeEquipmentTimes uint64 `head:"-" protofield:"-"` // 升级装备次数
	//
	//// X个武将到达X级，穿X件X级的装备
	//CaptainCount       uint64               `head:"-" protofield:"-"` // 有X个X级的武将，
	CaptainRebirth           uint64               `head:"-" protofield:"-"`
	CaptainLevel             uint64               `head:"-" protofield:"-"`
	CaptainQuality           shared_proto.Quality `head:"-" protofield:"-"`
	WearEquipmentCount       uint64               `head:"-" protofield:"-"` // 穿X件X级的装备
	WearEquipmentLevel       uint64               `head:"-" protofield:"-"`
	WearEquipmentRefineLevel uint64               `head:"-" protofield:"-"`
	WearEquipmentQuality     uint64               `head:"-" protofield:"-"`
	CaptainSoulQuality       shared_proto.Quality `head:"-" protofield:"-"` // x个y品质的将魂
	CaptainSoulLevel         uint64               `head:"-" protofield:"-"` // x个y等级的将魂
	// x个y星武将
	CaptainStar uint64 `head:"-" protofield:"-"`
	//// 武将强化X次
	//CaptainRefinedTimes uint64 `head:"-" protofield:"-"`
	//
	//// 卖武将个数
	//CaptainSellCount uint64 `head:"-" protofield:"-"`
	//
	//// X个修炼位升级到X级，修炼位在使用
	//TrainingCount uint64 `head:"-" protofield:"-"`
	TrainingLevel uint64 `head:"-" protofield:"-"`
	//
	//RecruitSoldierCount uint64 `head:"-" protofield:"-"` // 招募士兵数量
	//HealSoldierCount    uint64 `head:"-" protofield:"-"` // 治疗士兵数量
	//
	//// 士兵进阶到X级
	//SoldierLevel uint64 `head:"-" protofield:"-"`
	//
	//TowerFloor uint64 `head:"-" protofield:"-"` // 通关千重楼层数
	//
	//// X个X类型资源点
	//ResType            shared_proto.ResType      `head:"-" protofield:"-"`
	ResBuildingType shared_proto.BuildingType `head:"-" protofield:"-"`
	//ResourcePointCount uint64                    `head:"-" protofield:"-"`

	// 杀x个y级地区的z档怪物
	RegionData   *regdata.RegionData `head:"-" protofield:"-"`
	MonsterLevel uint64              `head:"-" protofield:"-"`

	// 野外任务怪物
	InvasionMonster *monsterdata.MonsterMasterData `head:"-" protofield:"-"`

	// 请教x次y导师
	Tutor *military_data.TutorData `head:"-" protofield:"-"`

	// 通关x副本
	PassDungeon *dungeon.DungeonData `head:"-" protofield:"-"`

	// 战胜X次Y层重楼密室
	PassSecretTower *towerdata.SecretTowerData `head:"-" protofield:"-"`

	// 将行营迁徙到X级的玉石地区
	TentRegion *regdata.RegionData `head:"-" protofield:"-"`

	// 通关斩将x关卡
	PassZhanjiangGuanqia *zhanjiang.ZhanJiangGuanQiaData `head:"-" protofield:"-"`

	// 官职，为nil表示不限官职
	CaptainOfficial *captain.CaptainOfficialData `head:"-" protofield:"-"`

	// 武将成长值
	CaptainAbilityExp uint64 `head:"-" protofield:"-"`

	// 拥有X个Y级Z类型的宝石
	GemType  uint64 `head:"-" protofield:"-"`
	GemLevel uint64 `head:"-" protofield:"-"`

	KillHomeNpcData *basedata.HomeNpcBaseData `head:"-" protofield:"-"`

	// 讨伐野怪类型
	MultiLevelNpcType shared_proto.MultiLevelNpcType `head:"-" protofield:"-"`
}

func (t *TaskTargetData) Init(filename string, configs interface {
	GetTechnologyData(uint64) *domestic_data.TechnologyData
	GetShop(uint64) *shop.Shop
	GetRegionData(uint64) *regdata.RegionData
	GetTutorData(uint64) *military_data.TutorData
	GetDungeonData(uint64) *dungeon.DungeonData
	GetSecretTowerData(uint64) *towerdata.SecretTowerData
	GetZhanJiangGuanQiaData(uint64) *zhanjiang.ZhanJiangGuanQiaData
	GetCaptainOfficialData(uint64) *captain.CaptainOfficialData
	GetResistXiongNuData(uint64) *xiongnu.ResistXiongNuData
	GetHomeNpcBaseData(uint64) *basedata.HomeNpcBaseData
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
}) {
	t.TotalProgress = t.X

	var err error
	switch t.Type {
	case shared_proto.TaskTargetType_TASK_TARGET_BASE_LEVEL:
	case shared_proto.TaskTargetType_TASK_TARGET_HERO_LEVEL:

	case shared_proto.TaskTargetType_TASK_TARGET_TECH_LEVEL:
		data := configs.GetTechnologyData(t.X)
		check.PanicNotTrue(data != nil, "%s 表中Id[%v]的数据配置的科技id没找到，%v", filename, t.Id, t.X)
		t.TotalProgress = data.Level
		t.TechGroup = data.Group

		t.SubType = t.X

	case shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL:

		intValue := t.parseEnumValue(filename, "Y", "升级建筑",
			shared_proto.BuildingType_name, shared_proto.BuildingType_value, 0)
		t.BuildingType = shared_proto.BuildingType(intValue)

		check.PanicNotTrue(!domestic_data.IsResourceBuilding(t.BuildingType), "%s 表中Id[%v]的Y配置建筑类型不能是资源点，实际配置的是[%v]，建筑类型格式%v", filename, t.Id, t.Y, config.EnumMapKeys(shared_proto.BuildingType_value, 0, 21, 22, 23, 24))

		t.SubType = uint64(t.BuildingType)

		if len(t.Z) > 0 {
			intValue := t.parseEnumValue(filename, "Z", "升级建筑",
				shared_proto.BuildingType_name, shared_proto.BuildingType_value, 0)
			t.BuildingType2 = shared_proto.BuildingType(intValue)

			check.PanicNotTrue(!domestic_data.IsResourceBuilding(t.BuildingType2), "%s 表中Id[%v]的Y配置建筑类型不能是资源点，实际配置的是[%v]，建筑类型格式%v", filename, t.Id, t.Z, config.EnumMapKeys(shared_proto.BuildingType_value, 0, 21, 22, 23, 24))
		}

	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_COUNT:

	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_LEVEL_COUNT:
		// 拥有X个Y级的武将
		if i, err := strconv.ParseInt(t.Y, 10, 32); err != nil {
			logrus.WithError(err).Panicf("%s 表中Id[%v]的Y应该武将等级(必须>0)，实际配置的是[%v]", filename, t.Id, t.Y)
		} else {
			t.CaptainLevel = u64.FromInt64(i)
		}

		if len(t.Z) > 0 {
			if t.CaptainRebirth, err = strconv.ParseUint(t.Z, 10, 32); err != nil {
				logrus.WithError(err).Panicf("%s 表中Id[%v]的Z应该武将转生等级(必须>=0)，实际配置的是[%v]", filename, t.Id, t.Z)
			}
		}

		check.PanicNotTrue(t.CaptainLevel > 0, "%s 表中Id[%v]的Y应该武将等级(必须>0)，实际配置的是[%v]", filename, t.Id, t.Y)

	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_SOUL_LEVEL_COUNT:
		// 拥有X个Y级的武将
		if i, err := strconv.ParseInt(t.Y, 10, 32); err != nil {
			logrus.WithError(err).Panicf("%s 表中Id[%v]的Y应该将魂等级(必须>0)，实际配置的是[%v]", filename, t.Id, t.Y)
		} else {
			t.CaptainSoulLevel = u64.FromInt64(i)
		}

		check.PanicNotTrue(t.CaptainSoulLevel > 0, "%s 表中Id[%v]的Y应该将魂等级(必须>0)，实际配置的是[%v]", filename, t.Id, t.Y)

	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_QUALITY_COUNT:
		// 拥有X个Y品质的武将

		intValue := t.parseEnumValue(filename, "Y", "武将品质",
			shared_proto.Quality_name, shared_proto.Quality_value, 0)

		t.CaptainQuality = shared_proto.Quality(intValue)

	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_UPSTAR:
		// X个Y星武将
		if i, err := strconv.ParseInt(t.Y, 10, 32); err == nil {
			t.CaptainStar = u64.FromInt64(i)
		}
		check.PanicNotTrue(t.CaptainStar > 1, "%s 表中Id[%v]的Y应该配置武将星级(必须>1)，实际配置的是[%v]", filename, t.Id, t.Y)

	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_EQUIPMENT:
		// 拥有X个武将身上穿Y件Z级装备
		if i, err := strconv.ParseInt(t.Y, 10, 32); err == nil {
			t.WearEquipmentCount = u64.FromInt64(i)
		}

		check.PanicNotTrue(t.WearEquipmentCount > 0, "%s 表中Id[%v]的Y应该配置武将装备件数(必须>0)，实际配置的是[%v]", filename, t.Id, t.Y)

		if i, err := strconv.ParseInt(t.Z, 10, 32); err == nil {
			t.WearEquipmentLevel = u64.FromInt64(i)
		}

		check.PanicNotTrue(t.WearEquipmentLevel > 0, "%s 表中Id[%v]的Z应该配置武将装备等级(必须>0)，实际配置的是[%v]", filename, t.Id, t.Z)

	case shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_LEVEL_Y:
		// 任意X个穿上的装备达到Y级
		if i, err := strconv.ParseInt(t.Y, 10, 32); err == nil {
			t.WearEquipmentLevel = u64.FromInt64(i)
		}

		check.PanicNotTrue(t.WearEquipmentLevel > 0, "%s 表中Id[%v]的Y应该配置武将装备等级(必须>0)，实际配置的是[%v]", filename, t.Id, t.Y)

	case shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_REFINE_LEVEL_Y:
		// 任意X个穿上的装备强化到Y级
		if i, err := strconv.ParseInt(t.Y, 10, 32); err == nil {
			t.WearEquipmentRefineLevel = u64.FromInt64(i)
		}

		check.PanicNotTrue(t.WearEquipmentRefineLevel > 0, "%s 表中Id[%v]的Y应该配置武将装备强化等级(必须>0)，实际配置的是[%v]", filename, t.Id, t.Y)

	case shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_QUALITY_Y:
		// 任意穿上X色装备Y件
		if i, err := strconv.ParseInt(t.Y, 10, 32); err == nil {
			t.WearEquipmentQuality = u64.FromInt64(i)
		}

		check.PanicNotTrue(t.WearEquipmentQuality > 0, "%s 表中Id[%v]的Y应该配置武将装备强化等级(必须>0)，实际配置的是[%v]", filename, t.Id, t.Y)

	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_REFINED_TIMES:

		//case shared_proto.TaskTargetType_TASK_TARGET_TRAINING_USE_COUNT:
		//
		//case shared_proto.TaskTargetType_TASK_TARGET_TRAINING_LEVEL_COUNT:
		//	// 拥有X个Y级的修炼位
		//	if i, err := strconv.ParseInt(t.Y, 10, 32); err == nil {
		//		t.TrainingLevel = u64.FromInt64(i)
		//	}
		//
		//	check.PanicNotTrue(t.TrainingLevel > 0, "%s 表中Id[%v]的Y应该配置修炼位等级(必须>0)，实际配置的是[%v]", filename, t.Id, t.Y)

	case shared_proto.TaskTargetType_TASK_TARGET_RECRUIT_SOLDIER_COUNT:

	case shared_proto.TaskTargetType_TASK_TARGET_HEAL_SOLDIER_COUNT:

	case shared_proto.TaskTargetType_TASK_TARGET_SOLDIER_LEVEL:

	case shared_proto.TaskTargetType_TASK_TARGET_TOWER_FLOOR:

	case shared_proto.TaskTargetType_TASK_TARGET_RESOURCE_POINT_COUNT:
		fallthrough
	case shared_proto.TaskTargetType_TASK_TARGET_COLLECT_RESOURCE:
		// 获取资源点个数
		if len(t.Y) > 0 {

			intValue := t.parseEnumValue(filename, "Y", "资源点资源",
				shared_proto.ResType_name, shared_proto.ResType_value)

			if bt, ok := domestic_data.GetResBuildingType(shared_proto.ResType(intValue)); ok {
				t.ResBuildingType = bt
			}
		}
	case shared_proto.TaskTargetType_TASK_TARGET_JOIN_GUILD:
		t.TotalProgress = 1

	case shared_proto.TaskTargetType_TASK_TARGET_BUY_GOODS:
		fallthrough
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BUY_GOODS:
		if shopType, err := strconv.ParseUint(t.Y, 10, 64); err != nil {
			logrus.Panicf("%s 表中Id[%v]的y[%s]应该配置商店类型", filename, t.Id, t.Y)
		} else {
			t.SubType = shopType
			check.PanicNotTrue(configs.GetShop(shopType) != nil, "%s 表中Id[%v]配置的商店类型没找到!%s", filename, t.Id, t.Y)
		}

	case shared_proto.TaskTargetType_TASK_TARGET_BUY_GOODS_COUNT:
		fallthrough
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BUY_GOODS_COUNT:
		if shopType, err := strconv.ParseUint(t.Y, 10, 64); err != nil {
			logrus.Panicf("%s 表中Id[%v]的y[%s]应该配置商店类型", filename, t.Id, t.Y)
		} else {
			t.SubType = shopType
			shop := configs.GetShop(shopType)
			check.PanicNotTrue(shop != nil, "%s 表中Id[%v]配置的商店类型没找到!%s", filename, t.Id, t.Y)

			if shopGoodsId, err := strconv.ParseUint(t.Z, 10, 64); err != nil {
				logrus.Panicf("%s 表中Id[%v]的z[%s]应该配置商店物品id", filename, t.Id, t.Z)
			} else {
				t.SubTypeId = shopGoodsId

				check.PanicNotTrue(shop.HasGoods(shopGoodsId), "%s 表中Id[%v]配置的商店物品id没找到!%s-%s", filename, t.Id, t.Y, t.Z)
			}
		}

	case shared_proto.TaskTargetType_TASK_TARGET_BAOWU_SELL:
		if baowuLevel, err := strconv.ParseUint(t.Y, 10, 64); err == nil {
			t.SubType = baowuLevel
		}
		check.PanicNotTrue(t.SubType > 0, "%s 表中Id[%v]的y(宝物的等级)应该大于0", filename, t.Id)

	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_SOUL_QUALITY_COUNT:
		// 拥有X个Y品质的武将
		intValue := t.parseEnumValue(filename, "Y", "将魂品质",
			shared_proto.Quality_name, shared_proto.Quality_value, 0)

		t.CaptainSoulQuality = shared_proto.Quality(intValue)

	case shared_proto.TaskTargetType_TASK_TARGET_JADE_NPC:
		// 杀x个y级地区的z档怪物
		if regionLevel, err := strconv.ParseUint(t.Y, 10, 64); err != nil {
			logrus.Panicf("%s 表中Id[%v]的y[%s]应该配置荣誉地区等级", filename, t.Id, t.Y)
		} else {
			t.RegionData = configs.GetRegionData(regdata.RegionDataID(shared_proto.RegionType_MONSTER, regionLevel))
			check.PanicNotTrue(t.RegionData != nil, "%s 表中Id[%v]配置的地图没找到!%s", filename, t.Id, t.Y)
			check.PanicNotTrue(t.RegionData.RegionType != shared_proto.RegionType_MONSTER, "%s 表中Id[%v]配置的地图不是荣誉地区!%s", filename, t.Id, t.Y)
		}

		if monLevel, err := strconv.ParseUint(t.Z, 10, 64); err != nil {
			logrus.Panicf("%s 表中Id[%v]的z[%s]应该配置怪物等级", filename, t.Id, t.Z)
		} else {
			t.MonsterLevel = monLevel
			check.PanicNotTrue(t.MonsterLevel > 0, "%s 表中Id[%v]配置的怪物等级必须>0!%s", filename, t.Id, t.Z)
		}

	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_CONSULT:

	case shared_proto.TaskTargetType_TASK_TARGET_HAS_CHALLENGE_DUNGEON:
		t.TotalProgress = 1
		t.PassDungeon = configs.GetDungeonData(t.X)
		check.PanicNotTrue(t.PassDungeon != nil, "%s 表中Id[%v]配置的X副本没找到!%d", filename, t.Id, t.X)

	case shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_SECRET_TOWER:
		fallthrough
	case shared_proto.TaskTargetType_TASK_TARGET_HELP_SECRET_TOWER:
		fallthrough
	case shared_proto.TaskTargetType_TASK_TARGET_HISTORY_CHALLENGE_SECRET_TOWER:
		if secretTowerId, err := strconv.ParseUint(t.Y, 10, 64); err == nil && secretTowerId != 0 {
			t.PassSecretTower = configs.GetSecretTowerData(secretTowerId)
			check.PanicNotTrue(t.PassSecretTower != nil, "%s 表中Id[%v]配置的Y重楼密室没找到!%s", filename, t.Id, t.Y)
		}

	case shared_proto.TaskTargetType_TASK_TARGET_HOME_IN_GUILD_REGION:
		t.TotalProgress = 1

	case shared_proto.TaskTargetType_TASK_TARGET_WIN_MULTI_LEVEL_MONSTER:
		if monsterLevel, err := strconv.ParseUint(t.Y, 10, 64); err == nil {
			t.MonsterLevel = monsterLevel
		}

		intValue := t.parseEnumValue(filename, "Z", "讨伐类型",
			shared_proto.MultiLevelNpcType_name, shared_proto.MultiLevelNpcType_value)
		t.MultiLevelNpcType = shared_proto.MultiLevelNpcType(intValue)

		t.SubType = SubTypeMultiLevelNpc(t.MultiLevelNpcType, t.MonsterLevel)

	case shared_proto.TaskTargetType_TASK_TARGET_EXPEL_FIGHT_MONSTER:

		t.TotalProgress = 1
		t.InvasionMonster = configs.GetMonsterMasterData(t.X)
		check.PanicNotTrue(t.InvasionMonster != nil, "%s 表中Id[%v]配置的怪物君主X没找到!%d", filename, t.Id, t.X)

		t.InvasionMonster.InitSpec(monsterdata.InvasionTask)
		t.SubType = t.InvasionMonster.Id

	case shared_proto.TaskTargetType_TASK_TARGET_ALL_RIGHT_QUESTION_AMOUNT:
		check.PanicNotTrue(t.X > 0, "%s 答题任务，答对题数必须>0. x:%v", filename, t.X)

	case shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL:
	case shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_TECH_LEVEL:
	case shared_proto.TaskTargetType_TASK_TARGET_RECRUIT_SOLDIER:
	case shared_proto.TaskTargetType_TASK_TARGET_JIU_GUAN_CONSULT:
	case shared_proto.TaskTargetType_TASK_TARGET_COLLECT_RESOURCE_TIMES:
	case shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_DUNGEON:
	case shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_TOWER:
	case shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_BAI_ZHAN:
	case shared_proto.TaskTargetType_TASK_TARGET_FIGHT_IN_JADE_REALM:
	case shared_proto.TaskTargetType_TASK_TARGET_CHAT_TIMES:
	case shared_proto.TaskTargetType_TASK_TARGET_BUILD_EQUIP_DAILY:
	case shared_proto.TaskTargetType_TASK_TARGET_SMELT_EQUIP:
	case shared_proto.TaskTargetType_TASK_TARGET_COLLECT_XIU_LIAN_EXP:
	case shared_proto.TaskTargetType_TASK_TARGET_GUILD_DONATE:
	case shared_proto.TaskTargetType_TASK_TARGET_ASSIST_GUILD_MEMBER:
	case shared_proto.TaskTargetType_TASK_TARGET_WATCH_VIDEO:
	case shared_proto.TaskTargetType_TASK_TARGET_BAI_ZHAN_JUN_XIAN:
	case shared_proto.TaskTargetType_TASK_TARGET_GUI_ZU_LEVEL:
	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_SOUL_COUNT:
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_JADE_ORE:
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_JADE:
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_GUILD_CONTRIBUTION:
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_GUILD_DONATE:
	case shared_proto.TaskTargetType_TASK_TARGET_EXPEL:
	case shared_proto.TaskTargetType_TASK_TARGET_DEFENSER_FIGHTING:
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_RECOVER_PROSPERITY:
	case shared_proto.TaskTargetType_TASK_TARGET_FU_SHEN:
	case shared_proto.TaskTargetType_TASK_TARGET_BUILD_EQUIP:
	case shared_proto.TaskTargetType_TASK_TARGET_KILL_HOME_NPC:

		t.TotalProgress = 1
		t.KillHomeNpcData = configs.GetHomeNpcBaseData(t.X)
		check.PanicNotTrue(t.KillHomeNpcData != nil, "%s 表中Id[%v]配置的玩家主城野怪没找到!%d", filename, t.Id, t.X)

	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_FISHING:
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BAI_ZHAN:
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_SMELT_EQUIP:
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_COLLECT_ACTIVE_BOX:
	case shared_proto.TaskTargetType_TASK_TARGET_ACTIVE_START_QUESTION_COUNT:
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_START_QUESTION:
	case shared_proto.TaskTargetType_TASK_TARGET_ZHANJIANG_GUANQIA_COMPLETE:
		//通关过关斩将X关卡，X填关卡编号
		t.TotalProgress = 1
		t.PassZhanjiangGuanqia = configs.GetZhanJiangGuanQiaData(t.X)
		check.PanicNotTrue(t.PassZhanjiangGuanqia != nil, "%s 表中Id[%v]配置的斩将X关卡没找到!%d", filename, t.Id, t.X)
	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_OFFICIAL_UPDATE:
		//册封X个武将Y官职，Y不填表示皆可
		if len(t.Y) > 0 {
			if oid, err := strconv.ParseUint(t.Y, 10, 64); err == nil {
				t.CaptainOfficial = configs.GetCaptainOfficialData(oid)
			}

			check.PanicNotTrue(t.CaptainOfficial != nil, "%s 表中Id[%v]配置的武将Y官职id没找到!%v", filename, t.Id, t.Y)
		}
	case shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_ABILITY_EXP:
		//X个武将成长值达到Y
		y, err := strconv.ParseUint(t.Y, 10, 64)
		check.PanicNotTrue(err == nil && y > 0, "%s 表中Id[%v]配置的Y必须>0!%d", filename, t.Id, t.Y)
		t.CaptainAbilityExp = y
	case shared_proto.TaskTargetType_TASK_TARGET_START_XIONGNU:
		//参与Y难度的抗击匈奴X次，Y不填表示皆可
		if len(t.Y) <= 0 {
			t.SubType = 0
		} else if level, err := strconv.ParseUint(t.Y, 10, 64); err != nil {
			logrus.Panicf("%s 表中Id[%v]配置的匈奴难度Y没找到!%d", filename, t.Id, t.Y)
		} else {
			check.PanicNotTrue(configs.GetResistXiongNuData(level) != nil, "%s 表中Id[%v]配置的匈奴难度Y没找到!%d", filename, t.Id, t.Y)
			t.SubType = level
		}
	case shared_proto.TaskTargetType_TASK_TARGET_FARM_HARVEST:
		//农场收获X数量的Y资源，若Y不填则两种皆可
		fallthrough
	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_FARM_STEAL:
		//累积从他人农场偷取X数量的Y资源，Y不填表示皆可
		intValue := t.parseEnumValue(filename, "Y", "偷菜资源类型",
			shared_proto.ResType_name, shared_proto.ResType_value)

		t.SubType = u64.FromInt32(intValue)

	case shared_proto.TaskTargetType_TASK_TARGET_ACCUM_ROBBING_RES:

		intValue := t.parseEnumValue(filename, "Y", "野区掠夺资源类型",
			shared_proto.ResType_name, shared_proto.ResType_value)

		t.SubType = u64.FromInt32(intValue)
	case shared_proto.TaskTargetType_TASK_TARGET_UNLOCK_BAOWU:
		fallthrough
	case shared_proto.TaskTargetType_TASK_TARGET_ROB_BAOWU:

		if len(t.Y) > 0 {
			level, err := strconv.ParseUint(t.Y, 10, 32)
			if err != nil {
				logrus.WithError(err).Panicf("%s 表中Id[%v]的%s应该配置宝物等级无效，实际配置的是[%v]", filename, t.Id, "Y", t.Y)
			}
			t.SubType = level
		}

	case shared_proto.TaskTargetType_TASK_TARGET_XUANYUAN_SCORE:
	case shared_proto.TaskTargetType_TASK_TARGET_HEBI:
		fallthrough
	case shared_proto.TaskTargetType_TASK_TARGET_HEBI_ROB:

		intValue := t.parseEnumValue(filename, "Y", "玉璧类型",
			shared_proto.Quality_name, shared_proto.Quality_value)

		t.SubType = u64.FromInt32(intValue)

	case shared_proto.TaskTargetType_TASK_TARGET_GEM:

		// 宝石等级 拥有X个Y级Z类型宝石
		t.GemLevel = t.parseU64Value(filename, "Y", "宝石等级", true)
		t.GemType = t.parseU64Value(filename, "Z", "宝石类型", true)

	case shared_proto.TaskTargetType_TASK_TARGET_INVASE_BAOZ:
		// 探索殷墟，到了就算
	case shared_proto.TaskTargetType_TASK_TARGET_KILL_JUN_TUAN:
		// 战胜Y级军团怪X次
		if len(t.Y) > 0 {
			level, err := strconv.ParseUint(t.Y, 10, 32)
			if err != nil {
				logrus.WithError(err).Panicf("%s 表中Id[%v]的%s应该配置军团怪等级无效，实际配置的是[%v]", filename, t.Id, "Y", t.Y)
			}
			t.SubType = level
		}
	case shared_proto.TaskTargetType_TASK_TARGET_BOOL:
		intValue := t.parseEnumValue(filename, "Y", "HeroBool类型",
			shared_proto.HeroBoolType_name, shared_proto.HeroBoolType_value)

		t.SubType = u64.FromInt32(intValue)
	}

	t.ShowProgress = t.TotalProgress
	if t.DontUpdateProgress {
		t.ShowProgress = 1
	}

}

func (t *TaskTargetData) parseEnumValue(filename, fieldName, typeKeyName string, enumName map[int32]string, enumValue map[string]int32, ignore ...int32) int32 {

	value := t.Z
	if fieldName == "Y" {
		value = t.Y
	}

	intValue := int32(0)
	if len(value) > 0 {
		if i, err := strconv.ParseInt(value, 10, 32); err != nil {
			intValue = -1
			for k, v := range enumName {
				if strings.ToUpper(v) == strings.ToUpper(value) {
					intValue = k
					break
				}
			}
		} else {
			intValue = int32(i)
		}
	}

	for _, v := range ignore {
		check.PanicNotTrue(intValue != v, "%s 表中Id[%v]的%s应该配置%s类型，实际配置的是[%v]，类型格式%v", filename, t.Id, fieldName, typeKeyName, value, config.EnumMapKeys(enumValue, ignore...))
	}

	check.PanicNotTrue(len(enumName[intValue]) > 0, "%s 表中Id[%v]的%s应该配置%s类型，实际配置的是[%v]，类型格式%v", filename, t.Id, fieldName, typeKeyName, value, config.EnumMapKeys(enumValue, ignore...))

	return intValue
}

func (t *TaskTargetData) parseU64Value(filename, fieldName, typeKeyName string, allowEmpty bool) uint64 {

	value := t.Z
	if fieldName == "Y" {
		value = t.Y
	}

	if len(value) <= 0 {
		if allowEmpty {
			return 0
		} else {
			logrus.Panicf("%s 表中Id[%v]的%s应该配置%s(整数类型)，实际配置的是[%v]", filename, t.Id, fieldName, typeKeyName, value)
		}
	}
	u64Value, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		logrus.WithError(err).Panicf("%s 表中Id[%v]的%s应该配置%s(整数类型)，实际配置的是[%v]", filename, t.Id, fieldName, typeKeyName, value)
	}

	return u64Value
}

//gogen:config
type TaskBoxData struct {
	_ struct{} `file:"任务/任务宝箱.txt"`
	_ struct{} `proto:"shared_proto.TaskBoxProto"`
	_ struct{} `protoconfig:"TaskBox"`

	// 宝箱id
	Id uint64

	Count uint64

	// 宝箱奖励
	Prize *resdata.Prize
}

func (d *TaskBoxData) Init(filename string, dataMap map[uint64]*TaskBoxData) {

	if d.Id > 1 {
		prev := dataMap[d.Id-1]
		check.PanicNotTrue(prev != nil, "%s 配置的任务宝箱奖励[%v]没找到，id必须从1开始连续配置", filename, d.Id-1)

		check.PanicNotTrue(prev.Count < d.Count, "%s 配置的任务宝箱奖励[%v]，完成的任务个数比上一级的要小", filename, d.Id)
	}
}

// 霸王之路任务
//gogen:config
type BwzlTaskData struct {
	_ struct{} `file:"任务/霸王之路.txt"`
	_ struct{} `proto:"shared_proto.TaskDataProto"`
	_ struct{} `protoconfig:"BwzlTask"`

	// 任务id
	Id uint64

	// 任务名字
	Name string

	// 任务内容
	Text string

	// 图标
	Icon string

	// 霸王之路第几天的任务
	Day uint64

	// 任务目标
	Target *TaskTargetData

	// 任务奖励
	Prize *resdata.Prize
}

// 霸王之路奖励
//gogen:config
type BwzlPrizeData struct {
	_ struct{} `file:"任务/霸王之路奖励.txt"`
	_ struct{} `proto:"shared_proto.BwzlPrizeDataProto"`
	_ struct{} `protoconfig:"BwzlPrize"`

	CollectPrizeTaskCount uint64 `key:"true"` // 领取了奖励的任务数量
	Icon                  string              // 奖励图标
	Prize                 *resdata.Prize      // 任务奖励

	CollectPrizeMsg pbutil.Buffer `head:"-" protofield:"-"` // 领取奖励缓存
}

func (data *BwzlPrizeData) Init() {
	data.CollectPrizeMsg = task.NewS2cCollectBwzlPrizeMsg(u64.Int32(data.CollectPrizeTaskCount)).Static()
}

// 称号数据
//gogen:config
type TitleData struct {
	_ struct{} `file:"任务/称号.txt"`
	_ struct{} `protogen:"true"`

	Id uint64

	// 称号名称
	Name string

	// 称号描述
	Desc string

	// 插图
	Image string `default:" "`

	// 称号品质
	Quality shared_proto.Quality

	// 称号属性
	SpriteStat *data.SpriteStat

	// 称号任务
	Task []*TitleTaskData `protofield:",config.U64a2I32a(GetTitleTaskDataKeyArray(%s)),int32"`

	// 消耗（贡品）
	TitleCost *resdata.Cost `head:"cost" default:"nullable"`

	TotalStat *data.SpriteStat `head:"-"`

	// 剧情id（客户端专用）
	PoltId uint64

	CountryChangeNameVoteCount uint64 `validator:"uint" desc:"国家改名票数"`

	nextData *TitleData

	completeMsg pbutil.Buffer
}

func (*TitleData) InitAll(filename string, dataMap map[uint64]*TitleData) {

	for i := 1; i <= len(dataMap); i++ {
		prev := dataMap[uint64(i-1)]
		next := dataMap[uint64(i)]
		if prev != nil {
			next.TotalStat = data.AppendSpriteStat(prev.TotalStat, next.SpriteStat)
		} else {
			next.TotalStat = next.SpriteStat
		}
	}

}

func (d *TitleData) Init(filename string, dataMap map[uint64]*TitleData) {

	if d.Id > 1 {
		prev := dataMap[d.Id-1]
		check.PanicNotTrue(prev != nil, "%s 的Id配置必须从1开始连续配置", filename)
		prev.nextData = d
	}

	d.completeMsg = task.NewS2cUpgradeTitleMsg(u64.Int32(d.Id)).Static()
}

func (d *TitleData) GetTotalStat() *data.SpriteStat {
	return d.TotalStat
}

func (d *TitleData) GetNextData() *TitleData {
	return d.nextData
}

func (d *TitleData) GetCompleteMsg() pbutil.Buffer {
	return d.completeMsg
}

//gogen:config
type TitleTaskData struct {
	_ struct{} `file:"任务/称号任务.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"task.proto"`

	Id uint64

	Name string

	Target *TaskTargetData
}

func SubTypeMultiLevelNpc(npcType shared_proto.MultiLevelNpcType, level uint64) uint64 {
	return uint64(npcType) | (level << 5)
}

// 活动任务
//gogen:config
type ActivityTaskData struct {
	_ struct{} `file:"活动/活动任务.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	Id     uint64
	Name   string
	Text   string
	Target *TaskTargetData
	Link   string
	Prize  *resdata.Prize
}

func (d *ActivityTaskData) Equal(data *ActivityTaskData) bool {
	return d == data
}
