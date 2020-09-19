package pvetroop

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type PveTroopData struct {
	_ struct{} `file:"杂项/pve部队.txt"`
	_ struct{} `proto:"shared_proto.PveTroopDataProto"`
	_ struct{} `protoconfig:"PveTroopDatas"`

	Id              uint64                    `head:"-,uint64(%s.PveTroopType)" protofield:"-"`
	PveTroopType    shared_proto.PveTroopType // 类型
	Capacity        uint64                    `default:"5"` // 队伍容量(队伍最高人数)
	MinCaptainCount uint64                    `default:"5"` // pve 时最少需要的武将数量
}

func (*PveTroopData) InitAll(filename string, configs interface {
	GetPveTroopDataArray() []*PveTroopData
}) {
	check.PanicNotTrue(len(configs.GetPveTroopDataArray()) == len(shared_proto.PveTroopType_name)-1, "%s 配置的pve组队配置应该每个队伍类型都配置一条!%d, %d", filename, len(configs.GetPveTroopDataArray()), len(shared_proto.PveTroopType_name)-1)
}

func (d *PveTroopData) Init(filename string) {
	check.PanicNotTrue(d.PveTroopType >= 1 && int(d.PveTroopType) <= len(shared_proto.PveTroopType_name), "pveTroopType必须从1开始逐渐递增 %v", d.PveTroopType)
}
