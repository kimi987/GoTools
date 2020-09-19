package domestic_data

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type SoldierLevelData struct {
	_     struct{} `file:"军事/士兵等级.txt"`
	_     struct{} `protoconfig:"soldier"`
	_     struct{} `proto:"shared_proto.SoldierLevelProto"`
	Level uint64   `key:"1"`

	Load uint64

	// 招募消耗
	RecruitCost *resdata.Cost

	// 伤兵治疗消耗
	WoundedCost *resdata.Cost

	// 升级消耗
	UpgradeCost *resdata.Cost

	// 升级到该等级需要的军营等级
	JunYingLevel uint64

	BaseStat []*data.SpriteStat `validator:"int>0,count=5,notNil,duplicate"`

	TotalStatSum uint64 `head:"-"`

	Desc string

	Models []uint64 `validator:"string,count=5" protofield:",config.U64a2I32a(%s)"` // 模型
}

func (d *SoldierLevelData) Init(fileName string) {
	for _, stat := range d.BaseStat {
		if d.TotalStatSum == 0 {
			d.TotalStatSum = stat.Sum4D()
		} else {
			check.PanicNotTrue(d.TotalStatSum == stat.Sum4D(), "%s 配置的士兵属性[base_stat]的4D属性必须都相同!%s", fileName, d.Level)
		}
	}
}

func (d *SoldierLevelData) GetBaseStat(race shared_proto.Race) *data.SpriteStat {

	index := int(race) - 1
	if index >= 0 && index < len(d.BaseStat) {
		return d.BaseStat[index]
	}

	return data.EmptyStat()
}

func (d *SoldierLevelData) GetModel(race shared_proto.Race) uint64 {
	index := int(race) - 1
	if index >= 0 && index < len(d.Models) {
		return d.Models[index]
	}
	return 0
}
