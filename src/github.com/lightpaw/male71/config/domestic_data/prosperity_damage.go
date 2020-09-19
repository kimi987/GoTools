package domestic_data

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/util/check"
)

const Percent = 100

//gogen:config
type ProsperityDamageBuffData struct {
	_ struct{} `file:"内政/繁荣度buff.txt"`
	_ struct{} `protogen:"true"`

	Id         uint64
	MinPercent uint64               `desc:"繁荣度受损最小范围。百分比 (min_percent, max_percent]" validator:"uint"`
	MaxPercent uint64               `desc:"繁荣度受损最大范围。百分比 (min_percent, max_percent]"`
	Desc       string               `desc:"描述"`
	BuffId     uint64               `desc:"触发的 buff id。BuffEffectDataProto.Id"`
	BuffData   *data.BuffEffectData `head:"-" protofield:"-"`
}

func (data *ProsperityDamageBuffData) InitAll(filename string, conf interface {
	GetProsperityDamageBuffDataArray() []*ProsperityDamageBuffData
	GetBuffEffectData(uint64) *data.BuffEffectData
}) {
	var preData *ProsperityDamageBuffData
	for _, d := range conf.GetProsperityDamageBuffDataArray() {
		if preData != nil {
			check.PanicNotTrue(d.MinPercent == preData.MaxPercent, "%v, min_damage 必须从0开始，且等于 上一条的 max_damage。id:%v min_damage:%v", filename, d.Id, d.MinPercent)
		}

		check.PanicNotTrue(d.MaxPercent <= 100, "%v, max_damage 必须 <= 100。id:%v max_damage:%v", filename, d.Id, d.MaxPercent)
		d.BuffData = conf.GetBuffEffectData(d.BuffId)
		check.PanicNotTrue(d.BuffData != nil, "%v, buff_id:%v 不存在。id:%v", filename, d.BuffId, d.Id)

		preData = d
	}
}

func (d *ProsperityDamageBuffData) Contains(per uint64) bool {
	return per > d.MinPercent && per <= d.MaxPercent
}
