package sub

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"time"
)

//gogen:config
type BuildingEffectData struct {
	_  struct{} `file:"内政/建筑效果.txt"`
	_  struct{} `proto:"shared_proto.DomesticEffectProto"`
	Id int      `protofield:"-"`

	// 仓库相关
	Capcity []*data.Amount `validator:"string,count=4,allNilOrNot,duplicate" parser:"data.ParseAmount"`
	//ProtectedCapcity []*data.ResAmount `parser:"data.ParseResAmount"`
	ProtectedCapcity *data.Amount `parser:"data.ParseAmount"`

	// 资源点相关
	OutputType    shared_proto.ResType `validator:"string" type:"enum"`
	Output        *data.Amount         `parser:"data.ParseAmount"`
	OutputCapcity *data.Amount         `parser:"data.ParseAmount"`

	// 军营相关
	SoldierCapcity        uint64 `validator:"uint"`
	SoldierOutput         uint64 `validator:"uint"`
	ForceSoldier          uint64 `validator:"uint"` // 强征士兵数
	NewSoldierOutput      uint64 `validator:"uint"`
	NewSoldierCapcity     uint64 `validator:"uint"`
	WoundedSoldierCapcity uint64 `validator:"uint"`
	RecruitSoldierCount   uint64 `validator:"uint"` // 招募数量

	// 士兵属性
	FarStat        *data.SpriteStat    `default:"nullable"` // 远程属性
	CloseStat      *data.SpriteStat    `default:"nullable"` // 近战属性
	SoldierRace    []shared_proto.Race `validator:"string" type:"enum"`
	SoldierStat    *data.SpriteStat    `default:"nullable"`
	AllSoldierStat *data.SpriteStat    `default:"nullable"` // 全部职业士兵属性加成

	// 单兵负载
	SoldierLoad uint64 `validator:"uint"`

	// 修炼馆基数
	TrainOutput  float64 `validator:"float64"`
	TrainCapcity float64 `validator:"float64"`
	TrainCoef    float64 `validator:"float64"` // 修炼馆加成系数

	TrainExpPerHour uint64 `validator:"uint"`

	// 建筑CD系数
	BuildingWorkerCdr             float64       `validator:"float64" protofield:"BuildingWorkerCoef"`
	BuildingWorkerFatigueDuration time.Duration `default:"0s"`

	// 科技CD系数
	TechWorkerCdr             float64       `validator:"float64" protofield:"TechWorkerCoef"`
	TechWorkerFatigueDuration time.Duration `default:"0s"`

	// 外使院
	SeekHelpCdr      time.Duration `default:"0s"`                 // 求助减少CD
	SeekHelpMaxTimes uint64        `validator:"uint" default:"0"` // 求助最大被帮助次数
	GuildDonateTimes uint64        `validator:"uint" default:"0"` // 联盟捐献次数

	// 城墙属性(攻击/防御/体力)
	HomeWallStat *data.SpriteStat `default:"nullable"` // 主城

	// 城墙固定伤害
	HomeWallFixDamage uint64 `validator:"uint"` // 主城

	FarmOutputType shared_proto.ResType `validator:"string" type:"enum"` // 农场产出类型
	FarmOutput     *data.Amount         `parser:"data.ParseAmount"`      // 农场产出加成

	// 税收
	Tax []*data.Amount `validator:"string,count=4,allNilOrNot,duplicate" parser:"data.ParseAmount"`

	// 守城属性
	AddedDefenseStat *data.SpriteStat `default:"nullable"`

	// 协防属性
	AddedAssistStat *data.SpriteStat `default:"nullable"`

	// 分身属性
	AddedCopyDefenseStat *data.SpriteStat `default:"nullable"`

	// 建筑消耗折扣
	BuildingCostReduceCoef float64 `validator:"float64" default:"0"`

	// 科技消耗折扣
	TechCostReduceCoef float64 `validator:"float64" default:"0"`
}

func (e *BuildingEffectData) Init(filename string) {
	if e.SoldierStat != nil {
		check.PanicNotTrue(len(e.SoldierRace) > 0, "%s 配置了给士兵加属性，但是士兵的职业没配置!%#v", filename, e)
	}
}

func (e *BuildingEffectData) GetFarmOutput(resType shared_proto.ResType) *data.Amount {
	if e.FarmOutputType == resType {
		return e.FarmOutput
	}
	return nil
}

func (e *BuildingEffectData) GetOutput(resType shared_proto.ResType) *data.Amount {
	if e.OutputType == resType {
		return e.Output
	}
	return nil
}

func (e *BuildingEffectData) GetOutputCapcity(resType shared_proto.ResType) *data.Amount {
	if e.OutputType == resType {
		return e.OutputCapcity
	}
	return nil
}

func (e *BuildingEffectData) GetCapcity(resType shared_proto.ResType) *data.Amount {
	index := int(resType) - 1
	if index >= 0 && index < len(e.Capcity) {
		return e.Capcity[index]
	}
	return nil
}

func (e *BuildingEffectData) GetTax(resType shared_proto.ResType) *data.Amount {
	index := int(resType) - 1
	if index >= 0 && index < len(e.Tax) {
		return e.Tax[index]
	}
	return nil
}

func (e *BuildingEffectData) IsCapcityEffect() bool {
	return len(e.Capcity) > 0 || !e.ProtectedCapcity.IsZero()
}

func (e *BuildingEffectData) IsTaxEffect() bool {
	return len(e.Tax) > 0
}

func (e *BuildingEffectData) IsResPointEffect() bool {
	return !e.Output.IsZero() || !e.OutputCapcity.IsZero()
}

func (e *BuildingEffectData) IsTrainEffect() bool {
	return e.TrainOutput > 0 || e.TrainCapcity > 0
}

func (e *BuildingEffectData) IsTrainCoefEffect() bool {
	return e.TrainCoef > 0
}

func (e *BuildingEffectData) IsSoldierEffect() bool {
	return e.SoldierCapcity > 0 || e.SoldierOutput > 0 || e.WoundedSoldierCapcity > 0 ||
		e.NewSoldierCapcity > 0 || e.NewSoldierOutput > 0 || e.RecruitSoldierCount > 0 ||
		e.SoldierLoad > 0
}

// 是否是远程士兵属性
func (e *BuildingEffectData) IsFarSoldierEffect() bool {
	return e.FarStat != nil
}

// 是否是近战士兵属性
func (e *BuildingEffectData) IsCloseSoldierEffect() bool {
	return e.CloseStat != nil
}

func (e *BuildingEffectData) GetSoldierFightTypeStat(isFar bool) *data.SpriteStat {
	if isFar {
		return e.FarStat
	} else {
		return e.CloseStat
	}
}

func (e *BuildingEffectData) IsSoldierStatEffect() bool {
	return e.SoldierStat != nil
}

func (e *BuildingEffectData) IsAllSoldierStatEffect() bool {
	return e.AllSoldierStat != nil
}

// 是城墙属性变更
func (e *BuildingEffectData) IsWallStatEffect() bool {
	return e.HomeWallStat != nil
}

// 是城墙固定伤害变更
func (e *BuildingEffectData) IsWallFixDamageEffect() bool {
	return e.HomeWallFixDamage != 0
}

func (e *BuildingEffectData) GetSoldierStat(race shared_proto.Race) *data.SpriteStat {
	for _, r := range e.SoldierRace {
		if r == race {
			return e.SoldierStat
		}
	}
	return nil
}
