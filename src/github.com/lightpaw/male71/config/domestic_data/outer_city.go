package domestic_data

import (
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/constants"
)

// 外城

// 外城的配置
//gogen:config
type OuterCityData struct {
	_ struct{} `file:"内政/外城.txt"`
	_ struct{} `proto:"shared_proto.OuterCityDataProto"`
	_ struct{} `protoconfig:"OuterCityDatas"`

	Id                    uint64                    `head:"-,uint64(%s.Diraction)"`                                               // 外城id
	Name                  string                                                                                                  // 外城的名字
	Desc                  string                                                                                                  // 外城的描述
	LockIcon              *icon.Icon                `protofield:"LockIconId,%s.Id"`                                               // 外城被锁的图标
	UnlockIcon            *icon.Icon                `protofield:"UnlockIconId,%s.Id"`                                             // 外城解锁的图标
	FirstLevelLayoutDatas []*OuterCityLayoutData    `head:"-" protofield:",config.U64a2I32a(GetOuterCityLayoutDataKeyArray(%s))"` // 外城布局
	Descs                 []*OuterCityBuildingDescData                                                                            // 外城建筑描述
	RegionModelRes        string                                                                                                  // 野外模型
	Diraction             shared_proto.PosDiraction `head:"id" protofield:"-"`                                                    // 方位
	funcType              uint64
	UnlockBeforeImage     string
	UnlockAfterImage      string
	UnlockDesc            []string                  `validator:"string,duplicate,notAllNil"`

	RecommandType uint64 `default:"0" validator:"int", protofield:"-"`
}

func (data *OuterCityData) Init(filename string, configs interface {
	GetOuterCityLayoutDataArray() []*OuterCityLayoutData
}) {
	for _, layout := range configs.GetOuterCityLayoutDataArray() {
		if layout.Level == 1 && layout.OuterCity == data {
			data.FirstLevelLayoutDatas = append(data.FirstLevelLayoutDatas, layout)
		}
	}

	switch data.Diraction {
	case shared_proto.PosDiraction_EAST:
		data.funcType = constants.FunctionType_TYPE_FEN_CHENG
	case shared_proto.PosDiraction_SOUTH:
		data.funcType = constants.FunctionType_TYPE_FEN_CHENG_2
	case shared_proto.PosDiraction_WEST:
		data.funcType = constants.FunctionType_TYPE_FEN_CHENG_3
	case shared_proto.PosDiraction_NORTH:
		data.funcType = constants.FunctionType_TYPE_FEN_CHENG_4
	default:
		logrus.Panicf("%s 分城配置里面居然有不是东南西北(1-4)的方位配置", filename)
	}
}

func (data *OuterCityData) GetFuncType() uint64 {
	return data.funcType
}

// 外城建筑描述
//gogen:config
type OuterCityBuildingDescData struct {
	_ struct{} `file:"内政/外城描述.txt"`
	_ struct{} `proto:"shared_proto.OuterCityBuildingDescDataProto"`

	Id   uint64     `protofield:"-"`
	Name string                                 // 建筑名字
	Icon *icon.Icon `protofield:"IconId,%s.Id"` // 建筑图标
	Desc string                                 // 建筑描述
}

const (
	MilitaryOuterCityType = 0
	EconomicOuterCityType = 1
)

// 外城布局
//gogen:config
type OuterCityLayoutData struct {
	_ struct{} `file:"内政/外城布局.txt"`
	_ struct{} `proto:"shared_proto.OuterCityLayoutDataProto"`
	_ struct{} `protoconfig:"OuterCityLayoutDatas"`

	Id     uint64                    // id
	Layout uint64 `validator:"uint"` // 布局Id

	OuterCity            *OuterCityData       `protofield:"-"`                                                                // 属于哪个外城
	ChangeTypeCost       *resdata.Cost        `default:"nullable"`
	UpgradeRequireLayout *OuterCityLayoutData `default:"nullable" protofield:"UpgradeRequireLayoutId,config.U64ToI32(%s.Id)"` // 升级到该级需要的布局，可为空，不为空的话该布局必须跟我在一个外城，且不等于自己
	UpgradeRequireIds    []*BuildingData      `protofield:",config.U64a2I32a(GetBuildingDataKeyArray(%s))"`                   // 升级到该级需要的主城建筑类型跟等级
	PrevLevel            *OuterCityLayoutData `default:"nullable" protofield:"-"`                                             // 上一级
	NextLevel            *OuterCityLayoutData `head:"-" protofield:",config.U64ToI32(%s.Id)"`                                 // 下一级
	DefaultUnlocked      bool                 `default:"true" protofield:"-"`                                                 // 只要外城解锁了，是否默认解锁，只有第一级需要配置
	Level                uint64               `head:"-" protofield:"-"`
	ChangeTypeLevel      *OuterCityLayoutData `head:"change_type_id" default:"nullable" protofield:"-"` // 改建后的指向等级
	MilitaryBuilding *OuterCityBuildingData `protofield:",config.U64ToI32(%s.Id)"`
	EconomicBuilding *OuterCityBuildingData `protofield:",config.U64ToI32(%s.Id)"`

	IsMain               bool               `head:"-" protofield:"-"` // 是否主建筑（布局）
}

func (data *OuterCityLayoutData) GetBuilding(t uint64) *OuterCityBuildingData {
	switch t {
	case MilitaryOuterCityType:
		return data.MilitaryBuilding
	default:
		return data.EconomicBuilding
	}
}

func (data *OuterCityLayoutData) Init(filename string) {

	data.Level = data.MilitaryBuilding.Level

	f := func(building, prevLevel *OuterCityBuildingData) {
		if prevLevel == nil {
			check.PanicNotTrue(building.Level == 1, "%s 配置的分城建筑没有配置前置建筑，但是等级不为1!%d-%v", filename, building.Level, building.BuildingType)
		} else {
			check.PanicNotTrue(building.BuildingType == prevLevel.BuildingType, "%s 配置的分城建筑配置的前置建筑 %d-%v，但是类型不相同!%d-%v", filename, prevLevel.Level, prevLevel.BuildingType, building.Level, building.BuildingType)
			check.PanicNotTrue(building.Level == prevLevel.Level+1, "%s 配置的分城建筑配置的前置建筑 %d-%v，但是等级差不为1!%d-%v", filename, prevLevel.Level, prevLevel.BuildingType, building.Level, building.BuildingType)
		}
	}
	if data.PrevLevel == nil {
		f(data.MilitaryBuilding, nil)
		f(data.EconomicBuilding, nil)
	} else {
		f(data.MilitaryBuilding, data.PrevLevel.MilitaryBuilding)
		f(data.EconomicBuilding, data.PrevLevel.EconomicBuilding)

		data.PrevLevel.NextLevel = data
	}

	if data.UpgradeRequireLayout != nil {
		check.PanicNotTrue(data.OuterCity == data.UpgradeRequireLayout.OuterCity, "%s 配置的分城建筑[%d]配置所属的外城跟前置升级需要的布局所在的外城不同! %s %s", filename, data.Id, data.OuterCity.Name, data.UpgradeRequireLayout.OuterCity.Name)
	}

	if data.MilitaryBuilding.BuildingType == shared_proto.BuildingType_SI_MA_FU {
		check.PanicNotTrue(data.ChangeTypeCost != nil, "%s 主建筑布局[%d]没有配置改建消耗", filename, data.Id)
		data.IsMain = true
	} else {
		check.PanicNotTrue(data.ChangeTypeCost == nil, "%s 不是主建筑布局[%d]不允许配置改建消耗", filename, data.Id)
	}
}

// 外城布局
//gogen:config
type OuterCityBuildingData struct {
	_ struct{} `file:"内政/外城建筑.txt"`
	_ struct{} `protogen:"true"`

	Id           uint64
	LockIcon     *icon.Icon                `protofield:"LockIconId,%s.Id,string"`                          // 建筑被锁的图标
	UnlockIcon   *icon.Icon                `protofield:"UnlockIconId,%s.Id,string"`                        // 建筑解锁的图标
	BuildingType shared_proto.BuildingType `protofield:"-"`                                                // 建筑类型
	Level        uint64                    `protofield:"-"`                                                // 等级
	BuildingData *BuildingData             `head:"-" protofield:"BuildingId,config.U64ToI32(%s.Id),int32"` // 建筑，服务器自己根据类型读取的一级数据
	Desc         string                                                                                    // 描述
}

func (data *OuterCityBuildingData) Init(filename string, configs interface {
	GetBuildingData(uint64) *BuildingData
	MainCityMiscData() *MainCityMiscData
}) {

	buildingId := BuildingId(data.BuildingType, data.Level)
	data.BuildingData = configs.GetBuildingData(buildingId)
	check.PanicNotTrue(data.BuildingData != nil, "%s 配置的分城建筑解锁的建筑没找到!%d-%v", filename, data.Level, data.BuildingType)

	check.PanicNotTrue(!configs.MainCityMiscData().IsMainCityBuildingType(data.BuildingType), "%s 外城中的建筑不应该是主城建筑类型!%d, %v", filename, data.Id, data.BuildingType)
}
