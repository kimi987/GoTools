package data

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type HeroLevelSubData struct {
	_          struct{} `file:"内政/君主等级.txt"`
	_          struct{} `protoconfig:"hero"`
	_          struct{} `proto:"shared_proto.HeroLevelProto"`
	Level      uint64   `key:"1" validator:"int>0"`
	UpgradeExp uint64

	AddSoldierCapacity    uint64              `validator:"uint" default:"0"` // 增加士兵上限
	EquipmentLevelLimit   uint64                                             // 装备等级上限
	CaptainSoulLevelLimit uint64                                             // 将魂等级上限
	CaptainLevelLimit     uint64
	UnlockedRaces         []shared_proto.Race `type:"enum"`

	CaptainTrainingLevel      []uint64 `validator:"int>0,duplicate,notAllNil" protofield:"-"`
	CaptainTrainingLevelLimit []uint64 `validator:"int>0,duplicate,notAllNil" protofield:"-"`

	StrategyLimit uint64
	SpLimit uint64

	TroopsCount        uint64 `default:"3"` // 队伍数量
	TroopsCaptainCount uint64 `default:"2"` // 队伍武将数量

	nextLevel *HeroLevelSubData `head:"-" protofield:"-"`

	CaptainOfficialId    []uint64 `head:"-"` // 武将官职ID
	CaptainOfficialCount []uint64 `head:"-"` // 武将官职数量

	CaptainOfficialIdCount map[uint64]uint64 `head:"-" protofield:"-"`
}

func (data *HeroLevelSubData) NextLevel() *HeroLevelSubData {
	return data.nextLevel
}

func (data *HeroLevelSubData) Init(filename string, dataMap map[uint64]*HeroLevelSubData) {
	// 武将官职数量配置
	data.CaptainOfficialIdCount = make(map[uint64]uint64)

	data.nextLevel = dataMap[data.Level+1]

	if data.nextLevel == nil && data.Level < uint64(len(dataMap)) {
		logrus.Panicf("%s 没有找到等级[%v]的君主数据, 君主等级必须从1开始连续配置", filename, data.Level+1)
	}

	check.PanicNotTrue(len(data.CaptainTrainingLevel) == len(data.CaptainTrainingLevelLimit), "%s 君主[%v]级数据配置的武将修炼位初始等级和最大等级个数不一致，%d, %d", filename, data.Level, len(data.CaptainTrainingLevel), len(data.CaptainTrainingLevelLimit))
	for i := 0; i < len(data.CaptainTrainingLevel); i++ {
		check.PanicNotTrue(data.CaptainTrainingLevel[i] <= data.CaptainTrainingLevelLimit[i], "%s 君主[%v]级数据配置的武将修炼位初始等级比等级上限还要大", filename, data.Level)
	}

	if data.nextLevel != nil {
		// 1级数据，放的1升2的数据
		data.UpgradeExp = data.nextLevel.UpgradeExp

		check.PanicNotTrue(data.EquipmentLevelLimit <= data.nextLevel.EquipmentLevelLimit, "%s 君主[%v]级数据配置的装备等级限制必上一级的小", filename, data.Level+1)
		check.PanicNotTrue(data.CaptainLevelLimit <= data.nextLevel.CaptainLevelLimit, "%s 君主[%v]级数据配置的武将等级限制必上一级的小", filename, data.Level+1)

		for i := 0; i < len(data.UnlockedRaces); i++ {
			check.PanicNotTrue(data.UnlockedRaces[i] == data.nextLevel.UnlockedRaces[i], "%s 君主[%v]级数据配置的武将开放职业跟上一级的不一致", filename, data.Level+1)
		}

		check.PanicNotTrue(len(data.CaptainTrainingLevel) <= len(data.nextLevel.CaptainTrainingLevelLimit), "%s 君主[%v]级数据配置的武将修炼位必须必上一级的要多", filename, data.Level+1)
		for i := 0; i < len(data.CaptainTrainingLevel); i++ {
			check.PanicNotTrue(data.CaptainTrainingLevel[i] == data.nextLevel.CaptainTrainingLevel[i], "%s 君主[%v]级数据配置的武将修炼位等级跟上一级的不一致", filename, data.Level+1)
			check.PanicNotTrue(data.CaptainTrainingLevelLimit[i] == data.nextLevel.CaptainTrainingLevelLimit[i], "%s 君主[%v]级数据配置的武将修炼位等级限制跟上一级的不一致", filename, data.Level+1)
		}

		check.PanicNotTrue(data.nextLevel.AddSoldierCapacity >= data.AddSoldierCapacity, "%s 君主[%v]级数据配置的武将增加带兵量必须>=上一级的", filename, data.Level+1)

		check.PanicNotTrue(data.TroopsCount <= data.nextLevel.TroopsCount, "英雄队伍数量下一级[%d]必须>=上一级[%d]", data.nextLevel.Level, data.Level)

	}

	check.PanicNotTrue(data.TroopsCount <= 16, "英雄队伍数量最大不能超过16个!Level: %d, Count: %d", data.Level, data.TroopsCount)
}
