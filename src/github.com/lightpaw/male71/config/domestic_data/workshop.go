package domestic_data

import (
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/check"
	"time"
)

//gogen:config
type WorkshopLevelData struct {
	_     struct{} `file:"内政/装备作坊.txt"`
	Level uint64

	Group []*resdata.PlunderGroup `validator:"int>0,duplicate,notAllNil"`
}

func (d *WorkshopLevelData) Init(filename string, configs interface {
	GetWorkshopDuration(uint64) *WorkshopDuration
}) {

	check.PanicNotTrue(len(d.Group) > 0, "%s 的配置中%d级没有配置掉落组包", filename, d.Level)

	for _, g := range d.Group {
		for _, item := range g.Item {
			check.PanicNotTrue(item.Equipment != nil, "%s 的配置中%d级配置的掉落组包[%d]中的掉落项[%d]，配置的掉落不是装备", filename, d.Level, g.Id, item.Id)

			check.PanicNotTrue(configs.GetWorkshopDuration(item.Equipment.Id) != nil, "%s 的配置中%d级配置的掉落组包[%d]中的掉落项[%d]中配置的装备[%d-%v]，没有配置对应的锻造时间", filename, d.Level, g.Id, item.Id, item.Equipment.Id, item.Equipment.Name)
		}
	}

}

func (d *WorkshopLevelData) RandomEquipment() []*goods.EquipmentData {
	result := make([]*goods.EquipmentData, len(d.Group))
	for i, g := range d.Group {
		result[i] = g.Try().Equipment
	}

	return result
}

//gogen:config
type WorkshopDuration struct {
	_  struct{} `file:"内政/装备作坊时间.txt"`
	Id uint64

	Duration time.Duration // 每一件装备的锻造时间
}

func (d *WorkshopDuration) Init(filename string) {
	check.PanicNotTrue(d.Duration > 0, "%s 配置%d 的锻造时间必须 > 0", filename, d.Id)
}

//gogen:config
type WorkshopRefreshCost struct {
	_ struct{} `file:"内政/装备作坊刷新消耗.txt"`
	_ struct{} `protoconfig:"WorkshopRefreshCosts"`
	_ struct{} `proto:"shared_proto.WorkshopRefreshCostProto"`

	Id   uint64
	Cost *resdata.Cost // 刷新消耗
}

func (*WorkshopRefreshCost) InitAll(filename string, configs interface {
	GetWorkshopRefreshCostArray() []*WorkshopRefreshCost
}) {
	for idx, data := range configs.GetWorkshopRefreshCostArray() {
		check.PanicNotTrue(data.Id == uint64(idx+1), "%s 配置%d 的锻造次数id必须从1开始，逐行加1> 0", filename, data.Id)
	}
}
