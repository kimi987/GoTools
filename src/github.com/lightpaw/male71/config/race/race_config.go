package race

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
)

var Array = []shared_proto.Race{
	shared_proto.Race_BU,
	shared_proto.Race_QI,
	shared_proto.Race_GONG,
	shared_proto.Race_CHE,
	shared_proto.Race_XIE,
}

//gogen:config
type RaceConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"杂项/职业杂项.txt"`

	dataMap map[shared_proto.Race]*RaceData
}

func (c *RaceConfig) Init(filename string, configs interface {
	GetRaceDataArray() []*RaceData
}) {
	c.dataMap = make(map[shared_proto.Race]*RaceData)

	for _, data := range configs.GetRaceDataArray() {
		check.PanicNotTrue(c.dataMap[data.Race] == nil, "race.txt 配置了重复的兵种，%v", data.Race)
		c.dataMap[data.Race] = data
	}

	for _, race := range Array {
		check.PanicNotTrue(c.dataMap[race] != nil, "race.txt 没有配置%v兵种的数据", race)
	}

	check.PanicNotTrue(len(c.dataMap) == len(Array), "兵种配置个数不一致，dataMap: %v, array: %v", c.dataMap, Array)
}

func (c *RaceConfig) GetData(race shared_proto.Race) *RaceData {
	return c.dataMap[race]
}

func (c *RaceConfig) GetProto(race shared_proto.Race) *shared_proto.RaceDataProto {
	data := c.GetData(race)
	if data != nil {
		return data.GetProto()
	}

	return nil
}
