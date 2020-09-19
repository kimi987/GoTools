package captain

import (
	"github.com/lightpaw/male7/util/check"
)

func CaptainLevelId(rebirth, level uint64) uint64 {
	return rebirth*10000 + level
}

//gogen:config
type CaptainLevelData struct {
	_              struct{} `file:"武将/武将等级.txt"`
	_              struct{} `protoconfig:"captain_level"`
	_              struct{} `proto:"shared_proto.CaptainLevelProto"`
	Id             uint64   `head:"-,CaptainLevelId(%s.Rebirth%c %s.Level)" protofield:"-"`
	Rebirth        uint64   `validator:"uint"`
	Level          uint64   `validator:"int>0"`
	SoldierCapcity uint64   `validator:"int>0" protofield:"-"`
	AbilityLimit   uint64   `validator:"uint"`

	UpgradeExp uint64
	//TotalUpgradeExp []uint64 `head:"-" protofield:"-"`
	GemSlotCount uint64 // 宝石槽位数量
	HasNewGemSlot bool `head:"-"` // 是否有新的槽位增加

	nextLevel *CaptainLevelData `head:"-" protofield:"-"`
}

func (data *CaptainLevelData) NextLevel() *CaptainLevelData {
	return data.nextLevel
}

//func (*CaptainLevelData) InitAll(array []*CaptainLevelData) {
//	var prev *CaptainLevelData
//	for _, v := range array {
//		n := len(v.UpgradeExp)
//		v.TotalUpgradeExp = make([]uint64, n)
//		for i := 0; i < n; i++ {
//			v.TotalUpgradeExp[i] += v.UpgradeExp[i]
//
//			if prev == nil {
//				v.TotalUpgradeExp[i] += prev.TotalUpgradeExp[i]
//			}
//		}
//
//		prev = v
//	}
//}

func (data *CaptainLevelData) Init(filename string, dataMap map[uint64]*CaptainLevelData) {
	check.PanicNotTrue(data.Level < 10000, "%s 武将等级不能超过10000", filename)

	data.nextLevel = dataMap[CaptainLevelId(data.Rebirth, data.Level+1)]

	if data.NextLevel() != nil {
		if data.Rebirth == data.nextLevel.Rebirth {
			check.PanicNotTrue(data.Level + 1 == data.nextLevel.Level, "%s 武将等级必须依次递增 %s %s", filename, data.nextLevel.Rebirth, data.nextLevel.Level)
		} else if data.Rebirth + 1 == data.NextLevel().Rebirth {
			check.PanicNotTrue(data.nextLevel.Level == 1, "%s 武将转生后的第一个等级必须是1 %s %s", filename, data.nextLevel.Rebirth, data.nextLevel.Level)
		} else {
			check.PanicNotTrue(false, "%s 武将的转生等级必须 等于或+1 上一级的转生等级 rebirth:%s level:%s", filename, data.nextLevel.Rebirth, data.nextLevel.Level)
		}
		// 1级数据，放的1升2的数据
		data.UpgradeExp = data.NextLevel().UpgradeExp
	}

	nextLevel := data.NextLevel()
	if nextLevel == nil {
		nextLevel = dataMap[CaptainLevelId(data.Rebirth + 1, 1)]
	}
	if nextLevel != nil {
		if data.GemSlotCount != nextLevel.GemSlotCount {
			check.PanicNotTrue(data.GemSlotCount <= nextLevel.GemSlotCount, "%s 武将镶嵌数必须依次递增 %s %s %s", filename, nextLevel.Rebirth, nextLevel.Level, nextLevel.GemSlotCount)
			nextLevel.HasNewGemSlot = true
		}
	}
}
