package captain

import (
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

//gogen:config
type CaptainRebirthLevelData struct {
	_ struct{} `file:"武将/武将转生.txt"`
	_ struct{} `protoconfig:"captain_rebirth_level"`
	_ struct{} `proto:"shared_proto.CaptainRebirthLevelProto"`

	Level uint64 `validator:"uint"`

	Cd time.Duration

	// 转生等级要求
	RequiredCaptainLevel uint64 `protofield:"-"`

	captainLevelDatas []*CaptainLevelData

	// 转生最小等级（升到本级需要的上一级转生的武将等级）
	CaptainLevelLimit *CaptainLevelData `head:"-" protofield:",config.U64ToI32(%s.Level)"`

	FirstCaptainLevel *CaptainLevelData `head:"-" protofield:"-"`

	// 转生属性
	SpriteStatPoint uint64 `validator:"uint"`

	// 转生统帅
	SoldierCapcity uint64 `validator:"uint"`

	// 成长值上限
	AbilityLimit uint64

	// 转生赠送的成长经验
	AbilityExp uint64

	// 转生前等级
	BeforeRebirthLevel uint64 `validator:"uint" default:"0"`

	HeroLevelLimit uint64 `validator:"uint" default:"0"`

	nextLevel *CaptainRebirthLevelData
}

func (d *CaptainRebirthLevelData) GetNextLevel() *CaptainRebirthLevelData {
	return d.nextLevel
}

func (d *CaptainRebirthLevelData) GetFirstLevelData() *CaptainLevelData {
	return d.captainLevelDatas[0]
}

func (d *CaptainRebirthLevelData) MustCaptainLevelData(level uint64) *CaptainLevelData {
	index := u64.Sub(level, 1)
	len := uint64(len(d.captainLevelDatas))
	if index < len {
		return d.captainLevelDatas[index]
	}

	return d.captainLevelDatas[len-1]
}

func (d *CaptainRebirthLevelData) Init(filename string, dataMap map[uint64]*CaptainRebirthLevelData, configs interface {
	GetCaptainLevelDataArray() []*CaptainLevelData
	GetCaptainLevelData(uint64) *CaptainLevelData
}) {

	for _, levelData := range configs.GetCaptainLevelDataArray() {
		if levelData.Rebirth == d.Level {
			d.captainLevelDatas = append(d.captainLevelDatas, levelData)
		}
	}
	check.PanicNotTrue(len(d.captainLevelDatas) > 0, "%s 武将转生配置，没有找到%v转的武将经验配置，看下武将升级经验配置是不是没有配这个", filename, d.Level)
	for i, v := range d.captainLevelDatas {
		lv := uint64(i + 1)
		check.PanicNotTrue(lv == v.Level, "%s 武将转生配置，没有找到%v转的%v级武将经验配置，看下武将升级经验配置是不是没有配这个", filename, d.Level, lv)
	}

	if d.Level > 0 {
		prev := dataMap[d.Level-1]
		check.PanicNotTrue(prev != nil, "%s 武将转生配置，没有找到%v转的配置，武将转生等级必须从0开始连续配置", filename, d.Level-1)
		prev.nextLevel = d
		d.CaptainLevelLimit = prev.MustCaptainLevelData(d.RequiredCaptainLevel)
		check.PanicNotTrue(d.CaptainLevelLimit != nil, "%s 武将转生配置，没有找到%v转的%v级的武将升级经验配置%v，看下武将升级经验配置是不是没有配这个", filename, prev.Level, d.RequiredCaptainLevel, prev.CaptainLevelLimit.Level)
	} else {
		d.CaptainLevelLimit = d.MustCaptainLevelData(d.RequiredCaptainLevel)
		check.PanicNotTrue(d.RequiredCaptainLevel == 1, "%s 武将转生配置，0转配置的转生要求等级必须是1", filename)
		check.PanicNotTrue(d.CaptainLevelLimit != nil, "%s 武将转生配置，没有找到0转的1级的武将升级经验配置，看下武将升级经验配置是不是没有配这个")
	}

	d.FirstCaptainLevel = d.MustCaptainLevelData(1)
}
