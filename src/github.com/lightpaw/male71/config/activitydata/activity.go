package activitydata

import (
	"github.com/lightpaw/male7/config/combine"
	"github.com/lightpaw/male7/config/data"
	"time"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/pb/shared_proto"
)

// 活动展示
//gogen:config
type ActivityShowData struct {
	_ struct{} `file:"活动/活动展示.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	Id        uint64             `validator:"int>0"` // 展示id
	SpineId   uint64             `validator:"uint" default"0"`
	Name      string                                 // 活动名
	NameIcon  string                                 // 活动名，美术字
	TimeShow  string                                 // 活动时间展示
	Desc      string                                 // 描述
	PrizeDesc string                                 // 奖励描述
	TimeRule  *data.TimeRuleData `protofield:"-"`

	Icon                string                 // 图标
	IconSelect          string                 // 选中图标
	ShowCountdown       bool `default:"false"` // 是否显示倒计时
	Sort                uint64                 // 排序
	Image               string                 // 底图
	ImagePos            uint64                 // 底图位置 0 左 1 中 2 上
	LinkName            string                 // 前往按钮名字
	LinkTaskTarget      *taskdata.TaskTargetData `default:"nullable" protofield:"-"` // 仅仅用作链接（前端要求）

	LinkTargetType      shared_proto.TaskTargetType `head:"-"` // 和[任务目标.txt]的type一样
	LinkTargetSubType   uint64 `head:"-"` // 每个类型的子类型
	LinkTargetSubTypeId uint64 `head:"-"` // 子类型id
}

func (d *ActivityShowData) Init(fileName string) {
	if d.LinkTaskTarget != nil {
		d.LinkTargetType = d.LinkTaskTarget.Type
		d.LinkTargetSubType = d.LinkTaskTarget.SubType
		d.LinkTargetSubTypeId = d.LinkTaskTarget.SubTypeId
	}
}

func (d *ActivityShowData) Equal(data *ActivityShowData) bool {
	return d == data
}

func (d *ActivityShowData) GetRecentShowTime(serverStartTime, ctime time.Time) (ok bool, startTime, endTime time.Time) {
	startTime = d.TimeRule.Next(serverStartTime, ctime)
	if startTime.IsZero() {
		return
	}

	endTime = startTime.Add(d.TimeRule.TimeDuration)
	ok = true
	return
}

// 收集活动
//gogen:config
type ActivityCollectionData struct {
	_ struct{} `file:"活动/收集活动.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	Id         uint64 `validator:"int>0"` // 活动ID
	Name       string                     // 活动名
	NameIcon   string                     // 活动名，美术字
	TimeShow   string                     // 活动时间展示
	Desc       string
	Icon       string // 图标
	IconSelect string // 选中图标
	Image      string // 图片
	Sort       uint64 // 排序

	Exchanges []*CollectionExchangeData
	TimeRule  *data.TimeRuleData `protofield:"-"`
}

func (d *ActivityCollectionData) GetExchange(id uint64) *CollectionExchangeData {
	for _, e := range d.Exchanges {
		if e.Id == id {
			return e
		}
	}
	return nil
}

func (d *ActivityCollectionData) Equal(data *ActivityCollectionData) bool {
	return d == data
}

func (d *ActivityCollectionData) GetRecentActivityTime(serverStartTime, ctime time.Time) (ok bool, startTime, endTime time.Time) {
	startTime = d.TimeRule.Next(serverStartTime, ctime)
	if startTime.IsZero() {
		return
	}

	endTime = startTime.Add(d.TimeRule.TimeDuration)
	ok = true
	return
}

// 收集兑换
//gogen:config
type CollectionExchangeData struct {
	_ struct{} `file:"活动/收集兑换.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	Id uint64 `validator:"int>0"` // 兑换id

	Combine *combine.GoodsCombineData
	Limit   uint64 `validator:"uint" default:"0"` // 兑换次限,0的话则没有兑换次数限制
}

func (d *CollectionExchangeData) equal(data *CollectionExchangeData) bool {
	return d == data
}

// 列表式任务活动
//gogen:config
type ActivityTaskListModeData struct {
	_ struct{} `file:"活动/列表式任务活动.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	Id         uint64 `protofield:"-"`
	Name       string // 活动名
	NameIcon   string // 活动名，美术字
	TimeShow   string // 活动时间展示
	Desc       string
	Icon       string // 图标
	IconSelect string // 选中图标
	Image      string // 图片
	Sort       uint64 // 排序

	TimeRule *data.TimeRuleData `protofield:"-"`
	Tasks    []*taskdata.ActivityTaskData
}

func (d *ActivityTaskListModeData) Equal(data *ActivityTaskListModeData) bool {
	return d == data
}

func (d *ActivityTaskListModeData) GetRecentActivityTime(serverStartTime, ctime time.Time) (ok bool, startTime, endTime time.Time) {
	startTime = d.TimeRule.Next(serverStartTime, ctime)
	if startTime.IsZero() {
		return
	}

	endTime = startTime.Add(d.TimeRule.TimeDuration)
	ok = true
	return
}

func (*ActivityTaskListModeData) InitAll(configs interface {
	GetActivityTaskListModeDataArray() []*ActivityTaskListModeData
}) {
	idMap := make(map[uint64]struct{})
	for _, d := range configs.GetActivityTaskListModeDataArray() {
		for _, task := range d.Tasks {
			_, ok := idMap[task.Id]
			check.PanicNotTrue(!ok, "列表式任务活动Id %v 有任务Id %v 与别的任务重复", d.Id, task.Id)
			idMap[task.Id] = struct{}{}
		}
	}
}
