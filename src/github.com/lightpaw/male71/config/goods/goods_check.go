package goods

import (
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type GoodsCheck struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"物品/物品检查.txt"`
}

func (*GoodsCheck) InitAll(configDatas interface {
	GetGoodsDataArray() []*GoodsData
	GetGemDataArray() []*GemData
	GetEquipmentData(uint64) *EquipmentData
	GetGoodsData(uint64) *GoodsData
	GetGemData(uint64) *GemData
}) {
	allGoodsIdMap := map[uint64]uint64{}
	f := func(id uint64) {
		if _, ok := allGoodsIdMap[id]; ok {
			check.PanicNotTrue(false, "存在相同的物品id! [%d]", id)
		}

		allGoodsIdMap[id] = 1
	}

	for _, goodsData := range configDatas.GetGoodsDataArray() {
		f(goodsData.Id)
	}
	for _, gemData := range configDatas.GetGemDataArray() {
		f(gemData.Id)
	}

	//for _, goodsData := range configDatas.GetGoodsDataArray() {
	//	if goodsData.EffectType == shared_proto.GoodsEffectType_EFFECT_PARTS {
	//		check.PanicNotTrue(goodsData.GoodsEffect.PartsCombineCount > 0, "物品[%d-%s]碎片配置的合成所需个数 <= 0", goodsData.Id, goodsData.Name)
	//
	//		if data := configDatas.GetEquipmentData(goodsData.GoodsEffect.PartsCombineId); data != nil {
	//			continue
	//		}
	//		if data := configDatas.GetGoodsData(goodsData.GoodsEffect.PartsCombineId); data != nil {
	//			continue
	//		}
	//		if data := configDatas.GetGemData(goodsData.GoodsEffect.PartsCombineId); data != nil {
	//			continue
	//		}
	//
	//		logrus.Panicf("物品[%d-%s]碎片配置的合成id，必须是装备，物品或者宝石的id之一，PartsCombineId: %s", goodsData.Id, goodsData.Name, goodsData.GoodsEffect.PartsCombineId)
	//	}
	//}
}
