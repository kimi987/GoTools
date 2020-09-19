package zhengwu

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/random"
	"github.com/lightpaw/male7/util/timeutil"
	"time"
	"github.com/lightpaw/male7/config/singleton"
)

// 政务

//gogen:config
type ZhengWuData struct {
	_ struct{} `file:"政务/政务.txt"`
	_ struct{} `proto:"shared_proto.ZhengWuDataProto"`

	Id         uint64                                                                           // 政务id
	Name       string                                                                           // 政务名字
	Icon       string                                                                           // 政务图标
	Quality    shared_proto.Quality                                                             // 品质
	Cost       *ZhengWuCompleteData           `head:"-" protofield:",config.U64ToI32(%s.Cost)"` // 政务直接完成的消耗
	Prize      *resdata.Prize                                                                   // 奖励
	Duration   time.Duration                                                                    // 完成耗时
	Proto      *shared_proto.ZhengWuDataProto `head:"-" protofield:"-"`                         // 政务数据
	ProtoBytes []byte                         `head:"-" protofield:"-"`                         // 政务数据缓存
}

func (data *ZhengWuData) Init(filename string, configs interface {
	GetZhengWuCompleteData(key uint64) *ZhengWuCompleteData
}) {
	data.Cost = configs.GetZhengWuCompleteData(uint64(data.Quality))
	check.PanicNotTrue(data.Cost != nil, "%s, %d-%s 政务根据品质的消耗没找到!%v", filename, data.Id, data.Name, data.Quality)

	data.Proto = data.Encode4Init()
	data.ProtoBytes = must.Marshal(data.Proto)
}

func (data *ZhengWuData) Encode4Init() *shared_proto.ZhengWuDataProto {
	var i interface{}
	i = data

	m, ok := i.(interface {
		Encode() *shared_proto.ZhengWuDataProto
	})
	if !ok {
		logrus.Errorf("ZhengWuData.Encode4Init() cast type fail")
	}

	return m.Encode()
}

//gogen:config
type ZhengWuRefreshData struct {
	_ struct{} `file:"政务/刷新消耗.txt"`
	_ struct{} `proto:"shared_proto.ZhengWuRefreshDataProto"`
	_ struct{} `protoconfig:"ZhengWuRefresh"`

	Times   uint64 `key:"true"`                      // 刷新次数，<=该次数，使用该消耗，如果找不到，就找最高的那个
	Cost    uint64 `validator:"uint" default:"1000"` // （作废，用NewCost）政务刷新消耗
	NewCost *resdata.Cost                            // 政务刷新消耗
}

//gogen:config
type ZhengWuCompleteData struct {
	_ struct{} `file:"政务/完成消耗.txt"`

	Id      uint64 `head:"-,uint64(%s.Quality)"` // 完成id
	Quality shared_proto.Quality                 // 品质
	Cost    uint64 `validator:"uint"`            // 政务刷新消耗点券，可以为0
}

//gogen:config
type ZhengWuMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"政务/其他.txt"`
	_ struct{} `proto:"shared_proto.ZhengWuMiscProto"`
	_ struct{} `protoconfig:"ZhengWuMisc"`

	AutoRefreshDuration []time.Duration                                           // 政务自动刷新间隔
	autoRefreshCycle    []*timeutil.CycleTime `head:"-" protofield:"-"`           // 下次自动刷新的cycle
	RandomCount         uint64                `protofield:"-"`                    // 政务随机数量
	FirstZhengWu        *ZhengWuData          `default:"nullable" protofield:"-"` // 第一个政务
}

func (data *ZhengWuMiscData) Init(filename string, configs interface {
	GetZhengWuData(key uint64) *ZhengWuData
	MiscConfig() *singleton.MiscConfig
}) {

	var autoRefreshCycle []*timeutil.CycleTime
	for _, d := range data.AutoRefreshDuration {
		autoRefreshCycle = append(autoRefreshCycle,
			timeutil.NewOffsetDailyTime(int64(d/time.Second)))
	}

	data.autoRefreshCycle = autoRefreshCycle
}

func (data *ZhengWuMiscData) NextAutoRefreshTime(ctime time.Time) time.Time {

	nextRefreshTime := ctime.Add(24 * time.Hour)
	for _, c := range data.autoRefreshCycle {
		t := c.NextTime(ctime)
		nextRefreshTime = timeutil.Min(nextRefreshTime, t)
	}

	return nextRefreshTime
}

//gogen:config
type ZhengWuRandomData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"政务/其他.txt"`

	zhengWuDatas     []*ZhengWuData `head:"-"` // 数据
	firstZhengWuData *ZhengWuData   `head:"-"` // 第一个政务
}

func (data *ZhengWuRandomData) Init(configs interface {
	GetZhengWuDataArray() []*ZhengWuData
	ZhengWuMiscData() *ZhengWuMiscData
}) {
	data.firstZhengWuData = configs.ZhengWuMiscData().FirstZhengWu

	data.zhengWuDatas = configs.GetZhengWuDataArray()
	configs.GetZhengWuDataArray()

	array := make([]*ZhengWuData, 0, len(configs.GetZhengWuDataArray()))
	for _, d := range configs.GetZhengWuDataArray() {
		if d != data.firstZhengWuData {
			array = append(array, d)
		}
	}

	// 将新手引导的政务移除掉
	data.zhengWuDatas = array

	check.PanicNotTrue(uint64(len(data.zhengWuDatas)) >= configs.ZhengWuMiscData().RandomCount, "政务随机，随机的政务的数量比当前配置的政务数量还要多，做不到不重复")
}

func (data *ZhengWuRandomData) Random(count uint64, completedFirst bool) []*ZhengWuData {
	result := data.random(count)
	if data.firstZhengWuData == nil || completedFirst {
		return result
	}

	for i, d := range result {
		if d == data.firstZhengWuData {
			result[0], result[i] = result[i], result[0]
			break
		}
	}
	if result[0] != data.firstZhengWuData {
		result[0] = data.firstZhengWuData
	}

	return result
}

func (data *ZhengWuRandomData) random(count uint64) []*ZhengWuData {
	length := len(data.zhengWuDatas)
	intCount := int(count)
	if intCount >= length {
		return data.zhengWuDatas
	}

	indexArray := random.NewMNIntIndexArray(length, intCount)

	result := make([]*ZhengWuData, 0, count)
	for _, i := range indexArray {
		result = append(result, data.zhengWuDatas[i])
	}

	return result
}
