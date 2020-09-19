package domestic_data

import (
	"github.com/lightpaw/male7/config/resdata"
	"time"
)

//gogen:config
type CombineCost struct {
	_ struct{} `file:"杂项/组合消耗.txt"`
	_ struct{} `proto:"shared_proto.CombineCostProto"`

	Id int `protofield:"-"`

	Cost *resdata.Cost

	BuildingWorkerTime time.Duration `default:"0s"` // 建筑队时间

	TechWorkerTime time.Duration `default:"0s"` // 科研时间

	Soldier uint64 `validator:"uint"` // 消耗的军营中的士兵数量

	InvadeTimes uint64 `validator:"uint"` // 消耗的出征次数
}
