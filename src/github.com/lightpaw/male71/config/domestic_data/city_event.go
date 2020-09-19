package domestic_data

import (
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/weight"
	"time"
)

// 城内事件城池等级数据
//gogen:config
type CityEventLevelData struct {
	_ struct{} `file:"内政/城内事件等级.txt"`

	BaseLevel     uint64           `key:"true"`              // 主城等级
	BaseLevelData *BaseLevelData   `head:"-" protofield:"-"` // 主城等级数据
	EventDatas    []*CityEventData // 事件
	eventRandomer *weight.WeightRandomer
}

func (d *CityEventLevelData) Init(filename string) {
	check.PanicNotTrue(len(d.EventDatas) > 0, "%s 中没配置事件!%d", filename, len(d.EventDatas))

	weights := make([]uint64, 0, len(d.EventDatas))
	for _, data := range d.EventDatas {
		weights = append(weights, data.Weight)
	}
	randomer, err := weight.NewWeightRandomer(weights)
	check.PanicNotTrue(err == nil, "%s 创建事件随机报错了， err: %v", filename, err)
	d.eventRandomer = randomer
}

func (d *CityEventLevelData) Random() *CityEventData {
	return d.EventDatas[d.eventRandomer.RandomIndex()]
}

// 城内事件
//gogen:config
type CityEventData struct {
	_ struct{} `file:"内政/城内事件.txt"`
	_ struct{} `proto:"shared_proto.CityEventDataProto"`
	_ struct{} `protoconfig:"city_event_data"`

	Id     uint64 `validator:"int>0"`                            // 事件id
	Desc   string `validator:"string>0"`                         // 事件描述
	Weight uint64 `validator:"int>0" default:"1" protofield:"-"` // 权重

	Cost  *CombineCost   // 消耗
	Prize *resdata.Prize // 奖励
}

// 城内事件其他数据
//gogen:config
type CityEventMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"内政/城内事件杂项.txt"`
	_ struct{} `proto:"shared_proto.CityEventMiscProto"`
	_ struct{} `protoconfig:"city_event_misc"`

	MaxTimes        uint64        `validator:"int>0"` // 最大的次数
	RecoverDuration time.Duration // 恢复间隔
	UnlockBaseLevel uint64        `head:"-"` // 解锁需要的主城的等级
}

func (d *CityEventMiscData) Init(filename string, datas interface {
	GetCityEventLevelDataArray() []*CityEventLevelData
	GetBaseLevelDataArray() []*BaseLevelData
}) {
	check.PanicNotTrue(d.RecoverDuration > 0, "%s 中配置的恢复间隔[%v]必须>0!", filename, d.RecoverDuration)

	check.PanicNotTrue(len(datas.GetCityEventLevelDataArray()) > 0, "%s 中配置城市等级事件起码要配置一条!%d", filename, len(datas.GetCityEventLevelDataArray()))

	minBaseLevel := uint64(0)
	for _, data := range datas.GetCityEventLevelDataArray() {
		if minBaseLevel == 0 {
			minBaseLevel = data.BaseLevel
		} else {
			minBaseLevel = u64.Min(minBaseLevel, data.BaseLevel)
		}
	}

	d.UnlockBaseLevel = minBaseLevel
	check.PanicNotTrue(len(datas.GetBaseLevelDataArray()) == int(minBaseLevel)+len(datas.GetCityEventLevelDataArray())-1,
		"%s 主城等级总共有 total[%d] 级, 城市事件等级总共有 city[%d] 级, 开放城市事件等级为 unlock[%d] 级, total 必须 = city + unlock - 1",
		len(datas.GetBaseLevelDataArray()), len(datas.GetCityEventLevelDataArray()), int(minBaseLevel))
}
