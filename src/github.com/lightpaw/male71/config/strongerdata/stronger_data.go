package strongerdata

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/util/check"
)

func GetStrongerDataId(level, tp uint64) uint64 {
	return level*10000 + tp
}

//gogen:config
type StrongerData struct {
	_ struct{} `file:"杂项/变强.txt"`
	_ struct{} `proto:"shared_proto.StrongerDataProto"`
	_ struct{} `protoconfig:"StrongerData"`

	Id uint64 `head:"-,GetStrongerDataId(%s.Level%c %s.Type)" protofield:"-"`

	Level uint64

	// 类型
	// 1、达到推荐战力值
	// 2、X个武将强化点数达到Y
	// 3、X个武将穿上Y件Z级装备
	// 4、X个武将穿上Y件Z星装备
	// 5、X个武将穿上Y个Z级宝石
	// 6、X个武将穿上Y级将魂
	// 7、君主升到X级
	// 8、士兵升到X阶
	// 9、宝石属性之和达到X
	Type uint64

	X uint64

	Y uint64 `validator:"uint"`

	Z uint64 `validator:"uint"`
}

func (*StrongerData) InitAll(filename string, array []*StrongerData, configs interface {
	GetHeroLevelSubDataArray() []*data.HeroLevelSubData
}) {
	dataMap := make(map[uint64][]*StrongerData)
	for _, v := range array {
		dataMap[v.Level] = append(dataMap[v.Level], v)
	}

	for _, v := range configs.GetHeroLevelSubDataArray() {
		curLevel := dataMap[v.Level]
		check.PanicNotTrue(len(curLevel) > 0, "%s 变强数据没有配置%v 级的数据", filename, v.Level)
	}
}
