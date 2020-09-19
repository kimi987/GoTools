package herodata

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/military_data"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/config/captain"
)

//gogen:config
type HeroLevelData struct {
	_     struct{}               `file:"内政/君主等级.txt"`
	Level uint64
	Sub   *data.HeroLevelSubData `head:"-"`

	CaptainTrainingLevel []*military_data.TrainingLevelData `validator:"int>0,duplicate,notAllNil"`

	captainOfficialCountData *captain.CaptainOfficialCountData

	nextLevel *HeroLevelData `head:"-"`
}

// 获得君主官职最大数量
func (data *HeroLevelData) HeroOfficialCount(officialId uint64) uint64 {
	if officialId == 0 {
		return 0
	}
	return data.Sub.CaptainOfficialIdCount[officialId]
}

func (data *HeroLevelData) GetCaptainOfficialCountData() *captain.CaptainOfficialCountData {
	return data.captainOfficialCountData
}

func (data *HeroLevelData) NextLevel() *HeroLevelData {
	return data.nextLevel
}

func (data *HeroLevelData) Init(filename string, dataMap map[uint64]*HeroLevelData, configs interface {
	GetHeroLevelSubData(uint64) *data.HeroLevelSubData
	GetCaptainRebirthLevelDataArray() []*captain.CaptainRebirthLevelData
	GetCaptainLevelData(uint64) *captain.CaptainLevelData
	GetCaptainOfficialCountDataArray() []*captain.CaptainOfficialCountData
}) {
	data.Sub = configs.GetHeroLevelSubData(data.Level)
	check.PanicNotTrue(data.Sub != nil, "%s 没有找到等级[%v]的君主数据", filename, data.Level)

	data.nextLevel = dataMap[data.Level+1]

	if data.nextLevel == nil && data.Level < uint64(len(dataMap)) {
		logrus.Panicf("%s 没有找到等级[%v]的君主数据, 君主等级必须从1开始连续配置", filename, data.Level+1)
	}

	for _, v := range configs.GetCaptainRebirthLevelDataArray() {
		d := configs.GetCaptainLevelData(captain.CaptainLevelId(v.Level, data.Sub.CaptainLevelLimit))
		check.PanicNotTrue(d != nil, "%s 等级[%v]的君主数据配置的武将等级上限是[%v]级，但是没找到 %v转对应[%v]级的武将升级经验，请检查武将升级经验配置表",
			filename, data.Level, data.Sub.CaptainLevelLimit, v.Level, data.Sub.CaptainLevelLimit)
	}

	// 武将官职数量
	var countData *captain.CaptainOfficialCountData
	for _, ocd := range configs.GetCaptainOfficialCountDataArray() {
		if ocd.HeroLevel <= data.Level {
			countData = ocd
		} else {
			break
		}
	}

	data.captainOfficialCountData = countData

	check.PanicNotTrue(countData != nil, "君主获得的武将封官数量为 nil")
	check.PanicNotTrue(countData.HeroLevel <= data.Level, "封官数量等级竟然大于君主等级")
	data.Sub.CaptainOfficialId = countData.OfficialId
	data.Sub.CaptainOfficialCount = countData.MaxCount

	for i, id := range data.Sub.CaptainOfficialId {
		data.Sub.CaptainOfficialIdCount[id] = data.Sub.CaptainOfficialCount[i]
	}
}
