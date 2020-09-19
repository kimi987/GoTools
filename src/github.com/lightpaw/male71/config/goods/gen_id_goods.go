package goods

// 物品
type GenIdGoods interface {
	Id() uint64
	GoodsData() GenIdGoodsData
}

// 自增id物品配置
type GenIdGoodsData interface {
	DataId() uint64
	GoodsType() GoodsType
}

// 生成id的物品类型
type GoodsType int32

const (
	EQUIPMENT GoodsType = 0
)
