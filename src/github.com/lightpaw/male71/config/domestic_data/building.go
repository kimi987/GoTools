package domestic_data

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/config/domestic_data/sub"
)

const denominator uint64 = 10000

func BuildingId(Type shared_proto.BuildingType, level uint64) uint64 {
	return uint64(Type)*denominator + level
}

func BuildingType(id uint64) shared_proto.BuildingType {
	return shared_proto.BuildingType(id / denominator)
}

func BuildingLevel(id uint64) uint64 {
	return id % denominator
}

//gogen:config
type BuildingData struct {
	_          struct{}                  `file:"内政/建筑.txt"`
	_          struct{}                  `proto:"shared_proto.BuildingDataProto"`
	_          struct{}                  `protoconfig:"BuildingData"`
	Id         uint64
	Type       shared_proto.BuildingType `type:"enum"`
	Level      uint64                    `validator:"int>0"`
	WorkTime   time.Duration
	Prosperity uint64                    `validator:"uint"`
	HeroExp    uint64                    `validator:"uint"`
	Desc       string
	Tips       string                    `default:" "`
	Icon       *icon.Icon                `protofield:"IconId,%s.Id" default:"小卒"` // 图标
	Model      string                    `default:" "`                            // 模型

	EffectDesc string              `default:" "` // 建筑效果描述
	Effect     *sub.BuildingEffectData `default:"null"`
	Notice     string              `default:" "` // 预告通知

	Cost       *resdata.Cost
	RequireIds []*BuildingData `protofield:",config.U64a2I32a(GetBuildingDataKeyArray(%s))"`

	// 要求主城等级
	BaseLevel *BaseLevelData `validator:"uint" default:"nullable" protofield:",config.U64ToI32(%s.Level)"`
}

func (d *BuildingData) CalculateWorkTime(coef float64) time.Duration {

	if d.WorkTime == 0 || coef <= 0 || coef >= 1 {
		return d.WorkTime
	}

	return time.Duration(float64(d.WorkTime) * coef)
}

func (d *BuildingData) Init(filename string, dataMap map[uint64]*BuildingData, configDatas interface {
	GetTieJiangPuLevelDataArray() []*TieJiangPuLevelData
	GetGuanFuLevelData(level uint64) *GuanFuLevelData
	MainCityMiscData() *MainCityMiscData
}) {
	check.PanicNotTrue(d.Id == BuildingId(d.Type, d.Level), "%s 建筑数据 %v-%v id不符合规则，id = type(%d) * 10000 + level(%v)", filename, d.Id, d.Type, d.Type, d.Level)

	prevLevel := dataMap[BuildingId(d.Type, d.Level-1)]
	check.PanicNotTrue(d.Level == 1 || prevLevel != nil, "建筑数据 %v-%v 没找到%v级的数据，等级必须从1开始连续配置", d.Id, d.Type, d.Level-1)

	_, ok := shared_proto.BuildingType_name[int32(d.Type)]
	check.PanicNotTrue(ok && d.Type != shared_proto.BuildingType_InvalidBuildingType, "%s 建筑数据 %v-%v 配置了无效的类型，有效类型是 %v", filename, d.Id, d.Type, config.EnumMapKeys(shared_proto.BuildingType_value, 0))

	if d.Type == shared_proto.BuildingType_TIE_JIANG_PU {
		check.PanicNotTrue(d.Level <= uint64(len(configDatas.GetTieJiangPuLevelDataArray())), "建筑数据 %v-%v 是铁匠铺，铁匠铺等级数据没找到!%d", d.Id, d.Type, d.Level)
	}

	if d.Type == shared_proto.BuildingType_GUAN_FU {
		check.PanicNotTrue(d.BaseLevel != nil, "建筑数据 %v-%v 是官府，官府必须配置需要的主城等级!%d", d.Id, d.Type, d.Level)
		d.BaseLevel.UnlockGuanFuLevel = u64.Max(d.Level, d.BaseLevel.UnlockGuanFuLevel)
	}

	for _, requireBuilding := range d.RequireIds {
		check.PanicNotTrue(configDatas.MainCityMiscData().IsMainCityBuildingType(requireBuilding.Type), "建筑数据 %v-%v 解锁需要的建筑类型必须是主城建筑!%d", d.Id, d.Type, d.Level)
	}

	switch d.Type {
	default:
		check.PanicNotTrue(d.Effect == nil, "%s 建筑数据 %v-%v 不能配置建筑效果", filename, d.Id, d.Type)
	case shared_proto.BuildingType_GUAN_FU:
		fallthrough
	case shared_proto.BuildingType_CHENG_QIANG:
		fallthrough
	case shared_proto.BuildingType_SHU_YUAN:
		fallthrough
	case shared_proto.BuildingType_JUN_YING:
		fallthrough
	case shared_proto.BuildingType_CANG_KU:
		fallthrough
	case shared_proto.BuildingType_XIU_LIAN_GUAN:
		fallthrough
	case shared_proto.BuildingType_XING_YING:
		fallthrough
	case shared_proto.BuildingType_WAI_SHI_YUAN:
		fallthrough
	case shared_proto.BuildingType_GOLD_PRODUCER:
		fallthrough
	case shared_proto.BuildingType_FOOD_PRODUCER:
		fallthrough
	case shared_proto.BuildingType_WOOD_PRODUCER:
		fallthrough
	case shared_proto.BuildingType_STONE_PRODUCER:
		fallthrough
	case shared_proto.BuildingType_SI_TU_FU:
		fallthrough
	case shared_proto.BuildingType_ZHU_BI_CHANG:
		fallthrough
	case shared_proto.BuildingType_CAI_SHI_CHANGE:
		fallthrough
	case shared_proto.BuildingType_JI_XIA_XUE_GONG:
		fallthrough
	case shared_proto.BuildingType_LU_BAN_GONG_FANG:
		fallthrough
	case shared_proto.BuildingType_SI_MA_FU:
		fallthrough
	case shared_proto.BuildingType_SHEN_SHE_YING:
		fallthrough
	case shared_proto.BuildingType_HU_BEN_YING:
		fallthrough
	case shared_proto.BuildingType_CHENG_FANG_JI_GUAN:
		fallthrough
	case shared_proto.BuildingType_YU_BEI_BING_SUO:
		check.PanicNotTrue(d.Effect != nil, "%s建筑数据 %v-%v 没有配置建筑效果", filename, d.Id, d.Type)
	}

	if d.IsResPoint() {
		check.PanicNotTrue(d.Effect != nil, "%s建筑数据 %v-%v 是资源点，但是没有配置 effect对象", filename, d.Id, d.Type)

		resType, ok := GetBuildingResType(d.Type)
		check.PanicNotTrue(ok && resType == d.Effect.OutputType, "%s 建筑数据 %v-%v 是资源点，但是配置的产出资源类型不一致", filename, d.Id, d.Type)

		check.PanicNotTrue(!d.Effect.Output.IsZero(), "%s 建筑数据 %v-%v 是资源点，但是没有配置 产出", filename, d.Id, d.Type)
		check.PanicNotTrue(!d.Effect.OutputCapcity.IsZero(), "建筑数据 %v-%v 是资源点，但是没有配置 产出上限", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_JUN_YING {
		check.PanicNotTrue(d.Effect.SoldierCapcity > 0, "%s 建筑数据 %v-%v 是军营，但是没有配置 士兵上限", filename, d.Id, d.Type)
		check.PanicNotTrue(d.Effect.SoldierOutput > 0, "%s 建筑数据 %v-%v 是军营，但是没有配置 士兵产出", filename, d.Id, d.Type)
		check.PanicNotTrue(d.Effect.ForceSoldier > 0, "%s 建筑数据 %v-%v 是军营，但是没有配置 强征士兵数", filename, d.Id, d.Type)
		check.PanicNotTrue(d.Effect.NewSoldierOutput > 0, "%s 建筑数据 %v-%v 是军营，但是没有配置 新兵产出", filename, d.Id, d.Type)
		check.PanicNotTrue(d.Effect.NewSoldierCapcity > 0, "建筑数据 %v-%v 是军营，但是没有配置 新兵上限", d.Id, d.Type)
		check.PanicNotTrue(d.Effect.WoundedSoldierCapcity > 0, "建筑数据 %v-%v 是军营，但是没有配置 伤兵上限", d.Id, d.Type)
		check.PanicNotTrue(d.Effect.RecruitSoldierCount > 0, "建筑数据 %v-%v 是军营，但是没有配置 招募士兵数量", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_CANG_KU {
		check.PanicNotTrue(!d.Effect.ProtectedCapcity.IsZero(), "建筑数据 %v-%v 是仓库，但是没有配置 资源保护上限", d.Id, d.Type)

		for _, c := range d.Effect.Capcity {
			check.PanicNotTrue(!c.IsZero(), "建筑数据 %v-%v 是仓库，但是没有配置 资源上限", d.Id, d.Type)
		}
	}

	if d.Type == shared_proto.BuildingType_XIU_LIAN_GUAN {
		//check.PanicNotTrue(d.Effect.TrainOutput > 0, "建筑数据 %v-%v 是修炼馆，产出基数必须 > 0", d.Id, d.Type)
		//check.PanicNotTrue(d.Effect.TrainCapcity > 0, "建筑数据 %v-%v 是修炼馆，容量基数必须 > 0", d.Id, d.Type)
		check.PanicNotTrue(d.Effect.TrainExpPerHour > 0, "建筑数据 %v-%v 是修炼馆，每小时产出经验必须大于0", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_SI_TU_FU {
		//check.PanicNotTrue(d.Effect.TrainCoef > 0, "建筑数据 %v-%v 是司徒府，产出系数必须 > 0", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_SHEN_SHE_YING {
		check.PanicNotTrue(d.Effect.FarStat != nil, "建筑数据 %v-%v 是神射营，远程职业增加属性(far_stat)必须配置!", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_HU_BEN_YING {
		check.PanicNotTrue(d.Effect.CloseStat != nil, "建筑数据 %v-%v 是虎贲营，近战职业增加属性(close_stat)必须配置!", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_ZHU_BI_CHANG || d.Type == shared_proto.BuildingType_CAI_SHI_CHANGE {
		//check.PanicNotTrue(d.Effect.OutputType != 0, "建筑数据 %v-%v 是铸币厂/采石场，产出类型(output_type)必须配置!", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_JI_XIA_XUE_GONG {
		//check.PanicNotTrue(d.Effect.TechWorkerCdr > 0, "建筑数据 %v-%v 是稷下学宫，科技研究速度必须配置!", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_LU_BAN_GONG_FANG {
		//check.PanicNotTrue(d.Effect.BuildingWorkerCdr > 0, "建筑数据 %v-%v 是稷下学宫，建筑建造研究速度必须配置!", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_CHENG_FANG_JI_GUAN {
		//check.PanicNotTrue(d.Effect.HomeWallStat != nil, "建筑数据 %v-%v 是城防机关，主城属性必须配置!", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_YU_BEI_BING_SUO {
		//check.PanicNotTrue(d.Effect.RecruitSoldierCount > 0, "建筑数据 %v-%v 是预备兵所，recruit_soldier_count 必须配置!", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_GUAN_FU {
		check.PanicNotTrue(d.Effect.BuildingWorkerCdr > 0, "建筑数据 %v-%v 是官府，建筑cd系数必须 > 0", d.Id, d.Type)
		check.PanicNotTrue(configDatas.GetGuanFuLevelData(d.Level) != nil, "建筑数据 %v-%v 是官府，没找到官府等级数据>0", d.Id, d.Type)
		check.PanicNotTrue(d.Prosperity > 0, "建筑数据 %v-%v 是官府，升级获得的繁荣度必须>0!", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_CHENG_QIANG {
		check.PanicNotTrue(d.Effect.HomeWallStat != nil, "建筑数据 %v-%v 是城墙，建筑城墙属性必须配置", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_SHU_YUAN {
		check.PanicNotTrue(d.Effect.TechWorkerCdr > 0, "建筑数据 %v-%v 是书院，科技cd系数必须 > 0", d.Id, d.Type)
	}

	if d.Type == shared_proto.BuildingType_WAI_SHI_YUAN {
		check.PanicNotTrue(d.Effect.SeekHelpCdr > 0, "建筑数据 %v-%v 是外使院，每次被帮助减少CD必须 > 0", d.Id, d.Type)
		check.PanicNotTrue(d.Effect.SeekHelpMaxTimes > 0, "建筑数据 %v-%v 是外使院，每个求助最大被帮助次数必须 > 0", d.Id, d.Type)
	}

	if d.Type == DonateBuilding {
		check.PanicNotTrue(d.Effect.GuildDonateTimes > 0, "建筑数据 %v-%v 是联盟捐献次数建筑，联盟捐献次数必须 > 0", d.Id, d.Type)
	}
}

const DonateBuilding = shared_proto.BuildingType_GUAN_FU

func IsResourceBuilding(t shared_proto.BuildingType) bool {
	return t == shared_proto.BuildingType_GOLD_PRODUCER ||
		t == shared_proto.BuildingType_FOOD_PRODUCER ||
		t == shared_proto.BuildingType_WOOD_PRODUCER ||
		t == shared_proto.BuildingType_STONE_PRODUCER
}

func (e *BuildingData) IsResPoint() bool {
	return IsResourceBuilding(e.Type)
}

func (e *BuildingData) GetResPointEffect() *sub.BuildingEffectData {
	if IsResourceBuilding(e.Type) {
		return e.Effect
	}
	return nil
}

func GetBuildingResType(t shared_proto.BuildingType) (shared_proto.ResType, bool) {
	switch t {
	case shared_proto.BuildingType_GOLD_PRODUCER:
		return shared_proto.ResType_GOLD, true

	case shared_proto.BuildingType_FOOD_PRODUCER:
		return shared_proto.ResType_FOOD, true

	case shared_proto.BuildingType_WOOD_PRODUCER:
		return shared_proto.ResType_WOOD, true

	case shared_proto.BuildingType_STONE_PRODUCER:
		return shared_proto.ResType_STONE, true
	}

	return -1, false
}

func GetResBuildingType(t shared_proto.ResType) (shared_proto.BuildingType, bool) {
	switch t {
	case shared_proto.ResType_GOLD:
		return shared_proto.BuildingType_GOLD_PRODUCER, true

	case shared_proto.ResType_FOOD:
		return shared_proto.BuildingType_FOOD_PRODUCER, true

	case shared_proto.ResType_WOOD:
		return shared_proto.BuildingType_WOOD_PRODUCER, true

	case shared_proto.ResType_STONE:
		return shared_proto.BuildingType_STONE_PRODUCER, true
	}

	return -1, false
}

// 建筑解锁

//gogen:config
type BuildingUnlockData struct {
	_ struct{} `file:"内政/建筑解锁.txt"`
	_ struct{} `proto:"shared_proto.BuildingUnlockDataProto"`
	_ struct{} `protoconfig:"BuildingUnlockData"`

	IntBuildingType    uint64        `key:"true" head:"-,uint64(%s.BuildingType)" protofield:"-"` // 建筑类型
	BuildingType       shared_proto.BuildingType                                                  // 建筑类型
	Desc               string        `default:" "`                                                // 描述
	Icon               string        `default:"icon"`                                             // 图标
	NotifyOrder        uint64        `default:"1"`                                                // 提示排序
	GuanFuLevel        *BuildingData `default:"nullable" protofield:",config.U64ToI32(%s.Level)"` // 解锁需要的官府等级
	HeroLevel          uint64        `invalidator:"uint" default:"nullable"`                      // 解锁需要的君主等级
	MainTaskSequence   uint64        `invalidator:"uint" default:"nullable"`                      // 解锁需要的主线任务
	BaYeStage          uint64        `invalidator:"uint" default:"nullable"`                      // 解锁需要的霸业阶段
	UnlockBuildingData *BuildingData `head:"-" protofield:"-"`                                    // 解锁的建筑数据
}

func (d *BuildingUnlockData) Init(filepath string, configs interface {
	GetBuildingDataArray() []*BuildingData
	MainCityMiscData() *MainCityMiscData
}) {
	check.PanicNotTrue(d.BuildingType != shared_proto.BuildingType_GUAN_FU, "%s 官府必须默认就解锁!%+v", filepath, d.BuildingType)

	for _, data := range configs.GetBuildingDataArray() {
		if data.Type == d.BuildingType && data.Level == 1 {
			d.UnlockBuildingData = data
			break
		}
	}

	check.PanicNotTrue(configs.MainCityMiscData().IsMainCityBuildingType(d.BuildingType), "%s 建筑解锁中的建筑应该是主城建筑类型!%s, %v", filepath, d.Desc, d.BuildingType)

	check.PanicNotTrue(d.UnlockBuildingData != nil, "%s 没有找到解锁数据要解锁的建筑的一级数据, %+v", filepath, d.BuildingType)

	if d.GuanFuLevel != nil {
		check.PanicNotTrue(d.GuanFuLevel.Type == shared_proto.BuildingType_GUAN_FU, "%s 中配置的建筑解锁只能够配置官府!%+v, %+v", filepath, d.BuildingType, d.GuanFuLevel.Type)
	}

	condCount := 0
	if d.GuanFuLevel != nil {
		condCount++
	}
	if d.HeroLevel > 0 {
		condCount++
	}
	if d.MainTaskSequence > 0 {
		condCount++
	}
	if d.BaYeStage > 0 {
		condCount++
	}

	check.PanicNotTrue(condCount > 0, "%s 中配置的解锁条件起码要配置一个!%+v", filepath, d.BuildingType)
}
