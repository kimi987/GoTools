package domestic_data

import (
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"time"
	"github.com/lightpaw/male7/config/domestic_data/sub"
)

//gogen:config
type TechnologyData struct {
	_ struct{} `file:"内政/科技.txt"`
	_ struct{} `proto:"shared_proto.TechnologyDataProto"`
	_ struct{} `protoconfig:"technology_data"`

	Id       uint64
	Name     string
	Desc     string
	Type     shared_proto.TechType `type:"enum"`
	Sequence uint64                `validator:"int>0"`
	Icon     string
	IntIcon  uint64 `validator:"uint" default:"0"` // int 图标，客户端控制的使用
	Group    uint64 `head:"-,uint64(%s.Type)*10000 + %s.Sequence"`
	Level    uint64 `validator:"int>0"`

	IsBigTech bool `default:"false"` // 是否是大科技

	Effect *sub.BuildingEffectData

	RequireBuildingIds []*BuildingData   `protofield:",config.U64a2I32a(GetBuildingDataKeyArray(%s))"`
	RequireTechIds     []*TechnologyData `protofield:",config.U64a2I32a(GetTechnologyDataKeyArray(%s))"`

	Cost     *resdata.Cost
	WorkTime time.Duration

	NextLevel *TechnologyData `head:"-" protofield:"NextLevelId,config.U64ToI32(%s.Id)"`

	MaxLevel uint64 `head:"-"` // 该科技最大等级
}

func (d *TechnologyData) GetNextLevel() *TechnologyData {
	return d.NextLevel
}

func (d *TechnologyData) CalculateWorkTime(coef float64) time.Duration {

	if d.WorkTime == 0 || coef <= 0 || coef >= 1 {
		return d.WorkTime
	}

	return time.Duration(float64(d.WorkTime) * coef)
}

func (*TechnologyData) InitAll(filename string, array []*TechnologyData) {

	groupMap := make(map[uint64]map[uint64]*TechnologyData)
	for _, v := range array {
		levelMap := groupMap[v.Group]
		if levelMap == nil {
			levelMap = make(map[uint64]*TechnologyData)
			groupMap[v.Group] = levelMap
		}

		check.PanicNotTrue(levelMap[v.Level] == nil, "科技配置%v 同一类型[Type:%v]-[Sequence:%v]存在重复的等级[level:%v] id: %v", filename, v.Type, v.Sequence, v.Level, v.Id)
		levelMap[v.Level] = v
	}

	for k, levelMap := range groupMap {
		for i := 0; i < len(levelMap); i++ {
			level := uint64(i + 1)
			data := levelMap[level]
			check.PanicNotTrue(data != nil, "科技配置%v 同一类型[Type:%v]-[Sequence:%v]等级[level:%v]的数据没找到，必须从1级开始连续配置", filename, shared_proto.TechType(k/10000), k%10000, level)

			data.MaxLevel = uint64(len(levelMap))

			if i > 0 {
				prevLevel := levelMap[uint64(i)]
				prevLevel.NextLevel = data
			}

			check.PanicNotTrue(data.WorkTime >= 0, "科技配置%v 同一类型[Type:%v]-[Sequence:%v]等级[level:%v]的升级所需时间必须>0", filename, shared_proto.TechType(k/10000), k%10000, level)
		}
	}

}
