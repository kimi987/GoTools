package domestic_data

import (
	"github.com/lightpaw/male7/pb/shared_proto"
)

// 主城其他数据
//gogen:config
type MainCityMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"内政/主城杂项.txt"`

	MainCityBuildingTypes    []shared_proto.BuildingType // 主城建筑类型
	mainCityBuildingTypesMap map[shared_proto.BuildingType]struct{}
}

func (data *MainCityMiscData) Init(filename string) {
	data.initTypeMap()
}

func (data *MainCityMiscData) IsMainCityBuildingType(buildingType shared_proto.BuildingType) bool {
	data.initTypeMap()

	_, exist := data.mainCityBuildingTypesMap[buildingType]
	return exist
}

func (data *MainCityMiscData) initTypeMap() {
	if data.mainCityBuildingTypesMap == nil {
		data.mainCityBuildingTypesMap = make(map[shared_proto.BuildingType]struct{})
		for _, tp := range data.MainCityBuildingTypes {
			data.mainCityBuildingTypesMap[tp] = struct{}{}
		}
	}
}
