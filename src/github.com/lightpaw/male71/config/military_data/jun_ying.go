package military_data

import (
	"time"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/config/resdata"
)

// 军营等级
//gogen:config
type JunYingLevelData struct {
	_ struct{} `file:"军事/军营.txt"`
	_ struct{} `proto:"shared_proto.JunYingLevelDataProto"`
	_ struct{} `protoconfig:"JunYingLevelData"`

	Level    uint64 `validator:"uint" key:"true"` // 酒馆等级
	MaxTimes uint64 `validator:"uint"`            // 最大次数

	RecoveryDuration time.Duration // 恢复间隔
}

// 军营其他数据
//gogen:config
type JunYingMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"军事/军营杂项.txt"`
	_ struct{} `proto:"shared_proto.JunYingMiscProto"`
	_ struct{} `protoconfig:"JunYingMisc"`

	DefaultTimes uint64 `validator:"uint" protofield:"-"` // 默认次数

	ForceAddSoldierMaxTimes uint64   `validator:"uint"`                                           // 强征最大次数，0表示不限次数
	ForceAddSoldierCost     []uint64 `validator:"uint,duplicate" default:"10,50,100,200,300,500"` // (作废，用ForceAddSoldierNewCost) 强征消耗
	ForceAddSoldierNewCost  []*resdata.Cost                                                       // 强征消耗
}

func (d *JunYingMiscData) Init(filename string) {

	check.PanicNotTrue(len(d.ForceAddSoldierNewCost) > 0, "%s 强征消耗没有配置", filename)
	check.PanicNotTrue(d.ForceAddSoldierNewCost[len(d.ForceAddSoldierNewCost)-1] != nil, "%s 最后一个强征消耗必须 > 0", filename)
}

func (d *JunYingMiscData) GetForceAddSoldierCost(times uint64) *resdata.Cost {
	n := uint64(len(d.ForceAddSoldierNewCost))
	if times < n {
		return d.ForceAddSoldierNewCost[times]
	}

	return d.ForceAddSoldierNewCost[n-1]
}
