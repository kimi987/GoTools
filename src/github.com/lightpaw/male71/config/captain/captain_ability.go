package captain

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
)

//gogen:config
type CaptainAbilityData struct {
	_          struct{}       `file:"武将/武将成长.txt"`
	_          struct{}       `protoconfig:"captain_ability"`
	_          struct{}       `proto:"shared_proto.CaptainAbilityProto"`

	Ability    uint64         `key:"1" validator:"int>0"`
	UpgradeExp uint64         `validator:"int>0"`
	SellPrice  *resdata.Prize // 卖掉获得的强化符
	FirePrice  *resdata.Prize // 解雇获得的强化符
	Desc       string

	Quality shared_proto.Quality `type:"enum"`
	Title   string               `protofield:"-"`

	nextLevel *CaptainAbilityData `head:"-" protofield:"-"`
	MaxLevel  uint64              `head:"-" protofield:"-"`

	// 解锁技能个数
	UnlockSpellCount uint64
}

func (data *CaptainAbilityData) NextLevel() *CaptainAbilityData {
	return data.nextLevel
}

//func (data *CaptainAbilityData) Model(race shared_proto.Race) int32 {
//	intRace := int(race)
//
//	var model uint64
//
//	if intRace <= 0 {
//		model = data.Models[0]
//	} else if intRace >= len(data.Models) {
//		model = data.Models[len(data.Models)-1]
//	} else {
//		model = data.Models[intRace-1]
//	}
//
//	return u64.Int32(model)
//}

func (data *CaptainAbilityData) Init(dataMap map[uint64]*CaptainAbilityData) {
	data.nextLevel = dataMap[data.Ability+1]

	data.MaxLevel = uint64(len(dataMap))
	if data.NextLevel() == nil && data.Ability < data.MaxLevel {
		logrus.Panicf("没有找到成长值[%v]的武将成长数据数据, 等级必须从1开始连续配置", data.Ability+1)
	}

	if data.NextLevel() != nil {
		// 1级数据，放的1升2的数据
		data.UpgradeExp = data.NextLevel().UpgradeExp
	}
}
