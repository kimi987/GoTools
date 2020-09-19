package combine

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/gen/pb/equipment"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
)

// 物品合成
//gogen:config
type GoodsCombineData struct {
	_ struct{} `file:"合成/物品合成.txt"`
	_ struct{} `proto:"shared_proto.GoodsCombineDataProto"`

	Id    uint64 `validator:"int>0"` // id
	Cost  *resdata.Cost              // 消耗
	Prize *resdata.Prize             // 奖励
}

// 装备合成
//gogen:config
type EquipCombineDatas struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"合成/装备合成.txt"`

	// 获得物品开启的铁匠铺装备合成数据
	obtainGoodsOpenCombineDataMap map[uint64]*EquipCombineData
}

func (d *EquipCombineDatas) Init(filename string, configData interface {
	GetEquipCombineDataArray() []*EquipCombineData
}) {
	d.obtainGoodsOpenCombineDataMap = make(map[uint64]*EquipCombineData, len(configData.GetEquipCombineDataArray()))

	for _, combineData := range configData.GetEquipCombineDataArray() {
		_, ok := d.obtainGoodsOpenCombineDataMap[combineData.CostGoods.Id]
		if ok {
			logrus.WithField("filename", filename).
				WithField("combineData", combineData.Id).
				WithField("combineName", combineData.GroupName).
				WithField("碎片id", combineData.CostGoods.Id).
				WithField("碎片名", combineData.CostGoods.Name).
				Panicln("装备合成中同一个碎片配置了多个合成分组")
		}
		d.obtainGoodsOpenCombineDataMap[combineData.CostGoods.Id] = combineData
	}
}

func (d *EquipCombineDatas) Len() int {
	return len(d.obtainGoodsOpenCombineDataMap)
}

func (d *EquipCombineDatas) GetOpenData(obtainGoodsId uint64) *EquipCombineData {
	return d.obtainGoodsOpenCombineDataMap[obtainGoodsId]
}

// 装备合成
//gogen:config
type EquipCombineData struct {
	_ struct{}             `file:"合成/装备合成.txt"`
	_ struct{}             `proto:"shared_proto.EquipCombineDataProto"`
	_ *goods.EquipmentData `protoconfig:"equip_combine"`

	Id           uint64                 `validator:"int>0"`                                                                         // 装备合成分组
	GroupName    string                 `validator:"string>0"`                                                                      // 分组的名字
	CostGoods    *goods.GoodsData       `head:"-" protofield:"CostGoodsId,config.U64ToI32(%s.Id)"`                                  // 消耗的物品
	CombineEquip []*goods.EquipmentData `head:"-" protofield:"CombineEquipId,config.U64a2I32a(goods.GetEquipmentDataKeyArray(%s))"` // 对应的合成的装备
	CombineData  []*GoodsCombineData                                                                                                // 合成的数据
	OpenMsg      pbutil.Buffer          `head:"-"`                                                                                  // 开启消息
}

func (d *EquipCombineData) Init(filename string) {
	check.PanicNotTrue(len(d.CombineData) > 0, "%s 中配置的合成数据起码要配置一条吧![%d]", filename, len(d.CombineData))

	d.CombineEquip = make([]*goods.EquipmentData, 0, len(d.CombineData))

	for idx, data := range d.CombineData {
		check.PanicNotTrue(data.Cost.TypeCount() == 1, "%s 中配置的合成数据，必须要消耗物品，且只能够消耗一种物品!%d", filename, len(data.Cost.Goods))
		check.PanicNotTrue(len(data.Cost.Goods) == 1, "%s 中配置的合成数据，必须要消耗物品，且必须消耗一种物品!%d", filename, len(data.Cost.Goods))

		if idx == 0 {
			d.CostGoods = data.Cost.Goods[0]
		} else {
			check.PanicNotTrue(d.CostGoods == data.Cost.Goods[0], "%s 中配置的合成数据，同一个分组消耗的物品必须相同!%v, %v", filename, d.CostGoods, data.Cost.Goods[0])
		}

		check.PanicNotTrue(data.Prize.TypeCount() == 1, "%s 中配置的合成数据，合成的东西不是装备，且只奖励一件装备!%d", filename, len(data.Prize.Equipment))
		check.PanicNotTrue(len(data.Prize.Equipment) == 1, "%s 中配置的合成数据，合成的东西不是装备，且只能够有一件装备!%d", filename, len(data.Prize.Equipment))

		d.CombineEquip = append(d.CombineEquip, data.Prize.Equipment[0])
	}

	d.OpenMsg = equipment.NewS2cOpenEquipCombineMsg(u64.Int32(d.Id)).Static()
}
