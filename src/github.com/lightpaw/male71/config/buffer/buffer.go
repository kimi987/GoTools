package buffer

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/util/check"
	"time"
)

//gogen:config
type BufferData struct {
	_ struct{} `file:"策略/增益.txt"`
	_ struct{} `protogen:"true"`

	// 增益id
	Id uint64 `validator:"int>0"`

	// 增益名称
	Name string

	// 增益名称描述
	NameDesc string

	// 图标
	Icon *icon.Icon `protofield:"IconId,%s.Id,string"` // 图标

	// 描述
	Desc string

	// 类型
	Type     uint64          `validator:"int>0"`
	TypeData *BufferTypeData `head:"-" protofield:"-"`

	ShowKeepDuration time.Duration
	ShowLevel        uint64

	BuffGoodsId   uint64           `desc:"buff物品ID"`
	BuffGoodsData *goods.GoodsData `head:"-" protofield:"-"`
}

func (d *BufferData) Init(filename string, conf interface {
	GetGoodsData(uint64) *goods.GoodsData
	GetBufferTypeData(uint64) *BufferTypeData
}) {
	d.BuffGoodsData = conf.GetGoodsData(d.BuffGoodsId)
	check.PanicNotTrue(d.BuffGoodsData != nil, "%v, buff_goods_id 不存在。id：%v buff_goods_id:%v", filename, d.Id, d.BuffGoodsId)
	check.PanicNotTrue(d.BuffGoodsData.GoodsEffect.BuffId > 0, "%v, buff_goods_id:%v 不是 buff 物品。id：%v", filename, d.BuffGoodsId, d.Id)
	d.TypeData = conf.GetBufferTypeData(d.Type)
	check.PanicNotTrue(d.TypeData != nil, "%v, type：%v 不存在。id：%v", filename, d.Type, d.Id)
}

func (d *BufferData) InitAll(filename string, conf interface {
	GetBufferData(uint64) *BufferData
	GetBuffEffectDataArray() []*data.BuffEffectData
}) {
	for _, d := range conf.GetBuffEffectDataArray() {
		if d.AdvantageId > 0 {
			check.PanicNotTrue(conf.GetBufferData(d.AdvantageId) != nil, "buff.txt, buff id:%v 的 advantage_id:%v 不存在", d.Id, d.AdvantageId)
		}
	}
}

//gogen:config
type BufferTypeData struct {
	_ struct{} `file:"策略/增益类型.txt"`
	_ struct{} `protogen:"true"`

	Id uint64

	// 类型名称
	TypeName string

	// 类型描述
	TypeDesc string

	// 是免战类型
	IsMian bool

	// 对应 buff 分组
	BuffGroup uint64 `validator:"uint"`

	// 图标
	Icon *icon.Icon `protofield:"IconId,%s.Id,string"` // 图标

	Sort uint64 `desc:"排序" default:"0"`
}

func (d *BufferTypeData) Init(filename string, conf interface {
	GetBuffEffectDataArray() []*data.BuffEffectData
}) {
	if !d.IsMian {
		var groupExisted bool
		for _, eff := range conf.GetBuffEffectDataArray() {
			if eff.Group == d.BuffGroup {
				groupExisted = true
				break
			}
		}
		check.PanicNotTrue(groupExisted, "%v 增益类型 id:%v 的 buff_group：%v 不存在", d.Id, d.BuffGroup)
	}
}
