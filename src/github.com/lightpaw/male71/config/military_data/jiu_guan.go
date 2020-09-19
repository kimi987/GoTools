package military_data

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/weight"
	"time"
)

// 酒馆等级
//gogen:config
type JiuGuanData struct {
	_ struct{} `file:"军事/酒馆.txt"`
	_ struct{} `proto:"shared_proto.JiuGuanDataProto"`
	_ struct{} `protoconfig:"JiuGuanData"`

	Level uint64 `validator:"uint" key:"true"` // 酒馆等级

	MaxTimes         uint64
	RecoveryDuration time.Duration `default:"3m"` // 酒馆次数恢复时间

	TutorDatas        []*TutorData                                                    // 导师数据
	RandomTutorWeight []uint64 `validator:"int>0,duplicate,notAllNil" protofield:"-"` // 导师的权重
	tutorRandomer     []*weight.U64WeightRandomer                                     // 导师随机
	InitTutorWeight   []uint64 `validator:"int,duplicate,notAllNil" protofield:"-"`   // 初始导师的权重
	initTutorRandomer *weight.U64WeightRandomer

	BroadcastContent      string `validator:"string>0"`           // 暴击广播内容
	BroadcastMinCritMulti uint64 `validator:"int" protofield:"-"` // 暴击广播的最小倍率
}

func (d *JiuGuanData) Init(filepath string) {
	check.PanicNotTrue(len(d.TutorDatas) == len(d.RandomTutorWeight), "%s 请教 %d 中配置的导师数量跟权重没有一一对应!", filepath, d.Level)
	check.PanicNotTrue(len(d.TutorDatas) == len(d.InitTutorWeight), "%s 请教 %d 中配置的导师数量跟初始权重没有一一对应!", filepath, d.Level)

	d.tutorRandomer = make([]*weight.U64WeightRandomer, len(d.TutorDatas)-1)

	indexArray := make([]uint64, len(d.RandomTutorWeight))
	for i := 0; i < len(d.RandomTutorWeight); i++ {
		indexArray[i] = uint64(i)
	}

	for i := 0; i < len(d.TutorDatas)-1; i++ {
		r, err := weight.NewU64WeightRandomer(d.RandomTutorWeight[i:], indexArray[i:])
		if err != nil {
			logrus.WithError(err).Panicln("%s 请教 %d 中配置的导师权重非法", filepath, d.Level)
		}
		d.tutorRandomer[i] = r
	}

	// 初始权重
	var initIndexWeight []uint64
	var initIndexArray []uint64
	for i, w := range d.InitTutorWeight {
		if w > 0 {
			initIndexWeight = append(initIndexWeight, w)
			initIndexArray = append(initIndexArray, uint64(i))
		}
	}

	if len(initIndexWeight) <= 0 {
		logrus.Panicln("%s 请教 %d 中配置的导师的初始权重全部为0", filepath, d.Level)
	}

	r, err := weight.NewU64WeightRandomer(initIndexWeight, initIndexArray)
	if err != nil {
		logrus.WithError(err).Panicln("%s 请教 %d 中配置的导师的初始权重非法", filepath, d.Level)
	}
	d.initTutorRandomer = r

}

func (d *JiuGuanData) MustTutorData(tutorIndex uint64) *TutorData {
	if tutorIndex >= uint64(len(d.TutorDatas)) {
		return d.TutorDatas[len(d.TutorDatas)-1]
	}

	return d.TutorDatas[tutorIndex]
}

func (d *JiuGuanData) InitRefresh() uint64 {
	return d.initTutorRandomer.Random()
}

// 刷新
func (d *JiuGuanData) Refresh(oldTutorIndex uint64) (newTutorIndex uint64, consultSuc bool) {
	if oldTutorIndex >= uint64(len(d.tutorRandomer)) {
		return
	}

	// 还是可以随机的
	consultSuc = true
	newTutorIndex = uint64(d.tutorRandomer[oldTutorIndex].Random())

	return
}

// 导师数据
//gogen:config
type TutorData struct {
	_ struct{} `file:"军事/酒馆导师.txt"`
	_ struct{} `proto:"shared_proto.TutorDataProto"`

	Id             uint64        `validator:"uint" protofield:"-"`                      // 导师id
	Name           string        `validator:"string>0"`                                 // 导师名字
	Image          string        `validator:"string>0"`                                 // 导师图片
	Weight         []uint64      `validator:"int>0,duplicate,notAllNil" protofield:"-"` // 权重
	Crit           []uint64      `validator:"int>0,notAllNil" protofield:"-"`           // 暴击倍率
	CritImgIndex   []int         `validator:"int,notAllNil" protofield:"-"`             // 暴击倍率图片下标
	Prize          *resdata.Prize                                                       // 奖励
	critRandomer   *weight.WeightRandomer                                               // 暴击随机
	ChatContent    string                                                               // 聊天内容
	RefreshMaxCost *resdata.Cost `default:"nullable"`
}

func (d *TutorData) Init(filepath string) {
	check.PanicNotTrue(len(d.CritImgIndex) == len(d.Weight), "%s TutorData.Init CritImgIndex配置的数量必须等于Weights配置的数量", filepath)
	check.PanicNotTrue(len(d.Crit) == len(d.Weight), "%s TutorData.Init Crit配置的数量必须等于Weights配置的数量", filepath)
	critRandomer, err := weight.NewWeightRandomer(d.Weight)
	if err != nil {
		logrus.WithError(err).Panicln("TutorData.Init() 权重报错")
	}

	d.critRandomer = critRandomer
}

func (d *TutorData) RandomCritMulti() (critMulti uint64, critImgIndex int) {
	index := d.critRandomer.RandomIndex()
	return d.Crit[index], d.CritImgIndex[index]
}

// 酒馆其他数据
//gogen:config
type JiuGuanMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"军事/酒馆杂项.txt"`
	_ struct{} `proto:"shared_proto.JiuGuanMiscDataProto"`
	_ struct{} `protoconfig:"JiuGuanMisc"`

	//DefaultTimes uint64 `validator:"uint" protofield:"-"` // 默认次数
	//MaxTimes     uint64 `validator:"uint"`                // 最大次数
	//RefreshFreeTimes   uint64 `validator:"int"`                 // 免费刷新的次数
	RefreshCostYuanBao []uint64 `validator:"uint" default:"0,10,20,30,40,50"` // 每次刷新消耗的点券, (作废，用 RefreshCost)
	RefreshCost        []*resdata.Cost                                        // 每次刷新消耗
	RecoveryDuration   time.Duration                                          // 恢复间隔

	FirstRefreshIndex uint64 `default:"2" protofield:"-"`
}

func (d *JiuGuanMiscData) Init(filename string) {
	check.PanicNotTrue(len(d.RefreshCost) > 0, "%v, 刷新消耗必须至少配一个", filename)
	if l := len(d.RefreshCost); l > 1 {
		check.PanicNotTrue(d.RefreshCost[l-1] != nil, "%v, 刷新消耗如果配置多个，至少最后一个不能是空的", filename)
	}
}

func (d *JiuGuanMiscData) GetRefreshCostDianquan(times uint64) *resdata.Cost {
	if times >= uint64(len(d.RefreshCost)) {
		return d.RefreshCost[len(d.RefreshCost)-1]
	}

	return d.RefreshCost[times]
}
