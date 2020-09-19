package domestic_data

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type BuildingLayoutData struct {
	_            struct{}                    `file:"内政/建筑布局.txt"`
	_            struct{}                    `proto:"shared_proto.BuildingLayoutProto"`
	_            struct{}                    `protoconfig:"building_layout"`
	Id           uint64                      `validator:"int>0"`
	BuildingType []shared_proto.BuildingType `type:"enum" protofield:"Building"`

	RequireBaseLevel uint64 `validator:"int>0"`
	RegionOffsetX    int    `validator:"int"`
	RegionOffsetY    int    `validator:"int"`

	IgnoreConflict bool `default:"false"`
}

func (d *BuildingLayoutData) Init(filename string) {

	check.PanicNotTrue(len(d.BuildingType) > 0, "城外布局配置%v 建筑类型列不能全部为空", filename)

	for _, t := range d.BuildingType {
		_, ok := GetBuildingResType(t)
		check.PanicNotTrue(ok, "城外布局-%v 配置了非资源点建筑类型[%d-%v]，具体请看 %v", d.Id, t, t, filename)
	}
}

//gogen:config
type BuildingLayoutMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"内政/建筑布局杂项.txt"`

	Gold  *data.Amount `validator:"uint"`
	Food  *data.Amount `validator:"uint"`
	Wood  *data.Amount `validator:"uint"`
	Stone *data.Amount `validator:"uint"`
}

func (d *BuildingLayoutMiscData) SingleResOutPutAmount(resType shared_proto.ResType) *data.Amount {
	switch resType {
	case shared_proto.ResType_GOLD:
		return d.Gold
	case shared_proto.ResType_FOOD:
		return d.Food
	case shared_proto.ResType_WOOD:
		return d.Wood
	case shared_proto.ResType_STONE:
		return d.Stone
	}

	logrus.Errorf("未知的资源类型:%+v", resType)

	return nil
}
