package domestic_data

import (
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
)

//gogen:config
type BaseLevelData struct {
	_ struct{} `file:"内政/主城等级.txt"`
	_ struct{} `proto:"shared_proto.BaseLevelProto"`
	_ struct{} `protoconfig:"base_level"`
	// 主城等级
	Level uint64 `validator:"int>0"`

	// 更新主城等级协议
	UpdateBaseLevelMsg pbutil.Buffer `head:"-" protofield:"-"`

	// 本级繁荣度，1级时候，表示0升1所需繁荣度
	Prosperity uint64 `head:"upgrade_prosperity" protofield:"-"`

	// 升级所需繁荣度，1级放的是1升2的数据
	UpgradeProsperity uint64 `head:"-"`

	UnlockGuanFuLevel uint64 `head:"-"`    // 解锁的官府等级，0表示不解锁官府等级
	UnlockPowerRange  string `default:" "` // 解锁的势力范围图标，空表示不解锁势力范围
	AppearanceRes     string `default:" "` // 外观资源

	TriggerCountdownPrizeProsperity uint64                  `validator:"uint" protofield:"-"` // 等级小于历史最高等级或者当前繁荣度<该等级繁荣度，插入马车奖励
	AddCountdownPrizeProsperity     uint64                  `validator:"uint" protofield:"-"` // 加多少繁荣度
	AddCountdownPrizeDesc           *CountdownPrizeDescData `protofield:"-"`                  // 马车描述

	prevLevel *BaseLevelData
}

func (d *BaseLevelData) GetPrevLevel() *BaseLevelData {
	return d.prevLevel
}

func (d *BaseLevelData) Init(filename string, dataMap map[uint64]*BaseLevelData) {

	d.prevLevel = dataMap[d.Level-1]
	if d.Level > 1 {
		check.PanicNotTrue(d.prevLevel != nil, "%s 没有找到等级[%v]的主城数据, 等级必须从1开始连续配置", filename, d.Level-1)
		check.PanicNotTrue(d.prevLevel.Prosperity < d.Prosperity, "%s [%v]级的主城数据, 配置的升级所需繁荣度必比上一级要大", filename, d.Level)

		d.prevLevel.UpgradeProsperity = d.Prosperity
	} else {
		check.PanicNotTrue(d.Prosperity == 1, "%s [%v]级的主城数据, 升级所需繁荣度必须==1", filename, d.Level)
	}

	d.UpdateBaseLevelMsg = region.NewS2cSelfUpdateBaseLevelMsg(u64.Int32(d.Level)).Static()
}
