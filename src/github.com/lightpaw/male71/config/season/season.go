package season

import (
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"time"
)

// 问卷调查
//gogen:config
type SeasonData struct {
	_ struct{} `file:"季节/季节.txt"`
	_ struct{} `proto:"shared_proto.SeasonDataProto"`
	_ struct{} `protoconfig:"SeasonDatas"`

	Id                    uint64         `head:"-,uint64(%s.Season)" protofield:"-"` // 问卷id，ID越大，越是没答过的越是优先答题
	Season                shared_proto.Season                                        // 季节
	Name                  string         `protofield:"-"`                            // 名字
	BgImg                 string                                                     // 背景图片
	ShowPrize             *resdata.Prize `protofield:"ShowPrize,%s.PrizeProto()"`    // 奖励
	Prize                 *resdata.Prize `protofield:"-"`                            // 奖励
	WorkerCdr             float64        `validator:"float64>=0"`                    //  增加建筑和研究效率, 这个值需要除以 1000，得到小数系数
	SecretTowerTimes      uint64         `validator:"uint"`                          // 增加重楼密室次数
	FarmBaseInc           float64        `validator:"float64>=0"`                    // 农场基础效率提升, 这个值需要除以 1000，得到小数系数
	AddMultiMonsterTimes  uint64         `validator:"uint"`                          // 增加讨伐野怪次数
	DecTroopSpeedRate     float64        `validator:"float64>=0"`                    // 降低行军速度, 这个值需要除以 1000，得到小数系数
	IncProsperityMultiple float64        `validator:"float64>=0"`                    // 增加领取繁荣度倍率, 这个值需要除以 1000，得到小数系数
	PrevSeason            *SeasonData    `head:"-" protofield:"-"`                   // 上一个季节
}

//gogen:config
type SeasonMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"季节/杂项.txt"`
	_ struct{} `proto:"shared_proto.SeasonMiscProto"`
	_ struct{} `protoconfig:"SeasonMisc"`

	SeasonDuration       time.Duration `head:"-" default:"24h"`
	SeasonSwitchDuration time.Duration `default:"10s"`
	seasons              []*SeasonData
}

func (data *SeasonMiscData) Init(filename string, configs interface {
	GetSeasonDataArray() []*SeasonData
}) {
	data.seasons = configs.GetSeasonDataArray()
	check.PanicNotTrue(len(data.seasons) == len(shared_proto.Season_name)-1, "%s 季节数据中，季节数据必须春夏秋冬都配置")

	prevSeason := data.seasons[len(data.seasons)-1]
	for idx, seasonData := range data.seasons {
		check.PanicNotTrue(shared_proto.Season(idx+1) == seasonData.Season, "%s 季节数据中，季节必须按照春夏秋冬配置")
		seasonData.PrevSeason = prevSeason
		prevSeason = seasonData
	}

	d, err := time.ParseDuration("24h")
	check.PanicNotTrue(err == nil, "%s 季节变化间隔错误", filename)
	data.SeasonDuration = d

	check.PanicNotTrue(data.SeasonDuration >= 10*time.Second, "%s 配置的季节变化间隔太短，最小需要10秒, 实际上是: %v秒", filename, data.SeasonDuration.Seconds())
}

func (c *SeasonMiscData) GetCurrentSeasonType(d time.Duration) shared_proto.Season {
	return c.GetCurrentSeason(d).Season
}

func (c *SeasonMiscData) GetCurrentSeason(d time.Duration) *SeasonData {
	if d < 0 {
		d = -d

		d = c.SeasonDuration - (d % c.SeasonDuration)
	}

	index := int(d/c.SeasonDuration) % len(c.seasons)
	return c.seasons[index]
}
