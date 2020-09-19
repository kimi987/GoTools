package fishing_data

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/weight"
	"time"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/captain"
)

const (
	FishTypeYuanbao = 0
	FishTypeFree    = 1

	FishCombo = 10
)

func FishId(fishType, times uint64) uint64 {
	return fishType*10000 + times
}

func FishTypeTimes(fishId uint64) (fishType, times uint64) {
	fishType = fishId / 10000
	times = fishId % 10000
	return
}

//钓鱼消耗
//gogen:config
type FishingCostData struct {
	_  struct{} `file:"钓鱼/钓鱼消耗.txt"`
	_  struct{} `proto:"shared_proto.FishingCostProto"`
	_  struct{} `protoconfig:"fishing_cost"`
	Id uint64   `head:"-,FishId(%s.FishType%c %s.Times)" protofield:"-"`

	Times uint64 // 钓鱼次数

	FreeTimes     uint64 `validator:"int" default:"0"` // 免费次数
	DiscountTimes uint64 `validator:"int" default:"0"` // 折扣次数
	DiscountCost  *resdata.Cost                        // 折扣消耗

	Cost *resdata.Cost // 非免费跟五折次数消耗

	FishType      uint64        `validator:"int"` // 钓鱼类型 0-高级 1-普通钓鱼
	DailyTimes    uint64        `validator:"int"` // 每日钓鱼次数，0表示不限次数
	FreeCountdown time.Duration `protofield:"-"`
}

func (d *FishingCostData) Init(filename string) {
	check.PanicNotTrue(d.Times <= 10, "%s 钓鱼次数配置太大，不能超过10次", filename)
	check.PanicNotTrue(d.FishType <= 1, "%s 钓鱼类型必须配置0或者1", filename)
}

func (data *FishingCostData) GetCost(fishingTimes uint64) *resdata.Cost {
	if fishingTimes < data.FreeTimes {
		return nil
	}

	if fishingTimes < data.FreeTimes+data.DiscountTimes {
		return data.DiscountCost
	}

	return data.Cost
}

// 钓鱼展示数据
//gogen:config
type FishingShowData struct {
	_ struct{} `file:"钓鱼/钓鱼展示.txt"`
	_ struct{} `proto:"shared_proto.FishingShowProto"`
	_ struct{} `protoconfig:"fishing_show"`

	Id            uint64               `protofield:"-"` // id
	GoodsData     *goods.GoodsData     `default:"nullable" protofield:",config.U64ToI32(%s.Id)"`
	GemData       *goods.GemData       `default:"nullable" protofield:",config.U64ToI32(%s.Id)"`
	EquipmentData *goods.EquipmentData `default:"nullable" protofield:",config.U64ToI32(%s.Id)"`
	CaptainData   *captain.CaptainData `default:"nullable" protofield:"CaptainSoulData,config.U64ToI32(%s.Id)"`

	Desc string // 描述

	FishType uint64 `validator:"int"` // 钓鱼类型 0-高级 1-普通钓鱼
	Out      bool                     // true表示展示在外面的奖励

	ShowType uint64 `default:"1"`
}

func (data *FishingShowData) Init(filename string) {
	notNilCount := 0
	if data.GoodsData != nil {
		notNilCount++
	}
	if data.GemData != nil {
		notNilCount++
	}
	if data.EquipmentData != nil {
		notNilCount++
	}
	if data.CaptainData != nil {
		notNilCount++
	}
	check.PanicNotTrue(notNilCount == 1, "%s 中配置的钓鱼展示只能够配置四种奖励类型中的一种!%s", filename, data.Id)
}

// 钓鱼随机数据
//gogen:config
type FishRandomer struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"钓鱼/钓鱼数据.txt"`

	fishGroupMap      map[uint64]*FishGroup
	fishGroup0        *FishGroup
	fishGroupPriority *FishGroup
	fishGroupMaxTimes uint64

	freeFishGroupMap      map[uint64]*FishGroup
	freeFishGroup0        *FishGroup
	freeFishGroupPriority *FishGroup
	freeFishGroupMaxTimes uint64
}

func (f *FishRandomer) Init(filename string, configData interface {
	GetFishDataArray() []*FishData
	GetFishingCaptainProbabilityDataArray() []*FishingCaptainProbabilityData
}) {

	f.fishGroupMap = make(map[uint64]*FishGroup)
	f.freeFishGroupMap = make(map[uint64]*FishGroup)

	f.fishGroupPriority = &FishGroup{}
	f.freeFishGroupPriority = &FishGroup{}

	for _, data := range configData.GetFishDataArray() {
		if data.IsPriority {
			group := f.freeFishGroupPriority
			if data.FishType == FishTypeYuanbao {
				group = f.fishGroupPriority
			}
			group.datas = append(group.datas, data)
			continue
		}

		groupMap := f.freeFishGroupMap
		if data.FishType == FishTypeYuanbao {
			groupMap = f.fishGroupMap
		}

		group := groupMap[data.FishTimes]
		if group == nil {
			group = &FishGroup{}
			groupMap[data.FishTimes] = group
		}

		group.datas = append(group.datas, data)
	}

	f.fishGroup0 = f.fishGroupMap[0]
	f.freeFishGroup0 = f.freeFishGroupMap[0]

	check.PanicNotTrue(f.fishGroup0 != nil, "%s 加载钓鱼随机数据失败，没有找到元宝钓鱼次数为0的配置", filename)
	check.PanicNotTrue(f.freeFishGroup0 != nil, "%s 加载钓鱼随机数据失败，没有找到免费钓鱼次数为0的配置", filename)

	arr := configData.GetFishingCaptainProbabilityDataArray()
	for times, g := range f.fishGroupMap {
		g.initRandomer(arr)
		f.fishGroupMaxTimes = u64.Max(f.fishGroupMaxTimes, times)
	}

	for times, g := range f.freeFishGroupMap {
		g.initRandomer(arr)
		f.freeFishGroupMaxTimes = u64.Max(f.freeFishGroupMaxTimes, times)
	}

	// 保底
	if len(f.fishGroupPriority.datas) > 0 {
		f.fishGroupPriority.initRandomer(arr)
	}

	if len(f.freeFishGroupPriority.datas) > 0 {
		f.freeFishGroupPriority.initRandomer(arr)
	}

}

// 钓鱼
func (f *FishRandomer) Fishing(fishType, fishTimes, setId uint64) *FishData {
	if fishType == FishTypeYuanbao {
		if fishTimes <= f.fishGroupMaxTimes {
			g := f.fishGroupMap[fishTimes]
			if g != nil {
				return g.Fishing(setId)
			}
		}
		return f.fishGroup0.Fishing(setId)
	} else {
		if fishTimes <= f.freeFishGroupMaxTimes {
			g := f.freeFishGroupMap[fishTimes]
			if g != nil {
				return g.Fishing(setId)
			}
		}
		return f.freeFishGroup0.Fishing(setId)
	}
}

// 保底钓鱼 可能为nil
func (f *FishRandomer) PriorityFishing(fishType, setId uint64) *FishData {
	if fishType == FishTypeYuanbao {
		if len(f.fishGroupPriority.datas) > 0 {
			return f.fishGroupPriority.Fishing(setId)
		}
	} else {
		if len(f.freeFishGroupPriority.datas) > 0 {
			return f.freeFishGroupPriority.Fishing(setId)
		}
	}

	return nil
}

type FishGroup struct {
	datas     []*FishData            // 所有的鱼
	// 钓鱼随机
	randomers map[uint64]*weight.WeightRandomer
	randomer  *weight.WeightRandomer
}

func (f *FishGroup) initRandomer(datas []*FishingCaptainProbabilityData) {
	// 获得权重
	wLen := len(f.datas)
	weights := make([]uint64, 0, wLen)
	for _, fishData := range f.datas {
		weights = append(weights, fishData.Weight)
	}

	if w, err := weight.NewWeightRandomer(weights); err != nil {
		logrus.WithError(err).Panicf("加载钓鱼随机数据失败")
	} else {
		f.randomer = w
	}

	f.randomers = make(map[uint64]*weight.WeightRandomer)
	for _, data := range datas {
		var found []int
		for i, fishData := range f.datas {
			if len(fishData.Prize.Captain) > 0 {
				for _, captain := range fishData.Prize.Captain {
					if captain.Id == data.CaptainId {
						found = append(found, i)
						break
					}
				}
			}
		}
		if len(found) > 0 {
			weights0 := make([]uint64, wLen)
			copy(weights0, weights)
			for _, idx := range found {
				d := f.datas[idx]
				weights0[idx] = d.Weight * data.Multiple
			}
			if w, err := weight.NewWeightRandomer(weights0); err != nil {
				logrus.WithError(err).Panicf("加载钓鱼随机数据失败")
			} else {
				f.randomers[data.CaptainId] = w
			}
		}
	}
}

// 钓鱼
func (f *FishGroup) Fishing(setId uint64) *FishData {
	if r := f.randomers[setId]; r != nil {
		return  f.datas[r.RandomIndex()]
	}
	return f.datas[f.randomer.RandomIndex()]
}

// 钓鱼数据
//gogen:config
type FishData struct {
	_ struct{} `file:"钓鱼/钓鱼数据.txt"`

	Id uint64 `validator:"int>0"` // 钓鱼id

	Prize      *resdata.Prize // 奖励
	prizeBytes []byte

	IsShow bool `protofield:"-"` // 是否钓鱼展示

	//作废字段
	//IsBroadcast bool // 是否广播

	Weight uint64 `validator:"int>0" protofield:"-"` // 权重

	FishType  uint64 `validator:"int"`             // 钓鱼类型 0-高级 1-普通钓鱼
	FishTimes uint64 `validator:"int" default:"0"` // 0表示除固定次数之外的次数

	IsPriority bool `default:"false" protofield:"-"`
}

func (f *FishData) GetPrizeBytes() []byte {
	return f.prizeBytes
}

func (f *FishData) Init(fileName string) {
	f.prizeBytes = must.Marshal(f.Prize.Encode4Init())

	check.PanicNotTrue(len(f.Prize.Captain) <= 1, "钓鱼[%s][%d]奖励里面配置的武将最多只能够配置一种!", fileName, f.Id)

	if len(f.Prize.Captain) == 1 {
		check.PanicNotTrue(f.Prize.CaptainCount[0] == 1, "钓鱼[%s][%d]奖励里面配置的武将最多只能够配置一个!", fileName, f.Id)
	}

	check.PanicNotTrue(f.FishType <= 1, "%s 钓鱼类型必须配置0或者1", fileName)

	check.PanicNotTrue(!(f.FishTimes != 0 && f.IsPriority), "%s 钓鱼次数[fish_times] 和 是否保底[is_priority]不能同时有值", fileName)
}

// 金杆钓
//gogen:config
type FishingCaptainProbabilityData struct {
	_ struct{} `file:"钓鱼/金杆钓.txt"`
	_ struct{} `protogen:"true"`

	CaptainId   uint64  `validator:"uint" key:"true"` // 武将id
	Multiple    uint64  // 倍率
}
