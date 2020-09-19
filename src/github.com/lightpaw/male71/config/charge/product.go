package charge

////gogen:config
//type ProductConfig struct {
//	_ struct{} `singleton:"true"`
//	_ struct{} `file:"充值/充值项.txt"`
//
//	product map[uint64]*ProductData
//}
//
//func (c *ProductConfig) Init(filename string, configs interface {
//	GetChargeObjDataArray() []*ChargeObjData
//}) {
//	productMap := make(map[uint64]*ProductData)
//	for _, v := range configs.GetChargeObjDataArray() {
//		exist := productMap[v.Id]
//		check.PanicNotTrue(exist == nil, "", )
//
//		pd := &ProductData{}
//		pd.id = v.Id
//
//		pd.data = v
//		pd.price = v.ChargeAmount
//		productMap[pd.id] = pd
//	}
//
//	c.product = productMap
//}
//
//func (c *ProductConfig) GetProductData(id uint64) *ProductData {
//	if len(c.product) > 0 {
//		return c.product[id]
//	}
//	return nil
//}

//gogen:config
type ProductData struct {
	_ struct{} `file:"充值/收费.txt"`

	Id uint64

	ProductId string

	ProductName string

	// 购买价格（人民币分）
	Price uint64

	// 具体的数据对象
	data interface{}
}

func (c *ProductData) GetData() interface{} {
	return c.data
}

func (c *ProductData) SetData(toSet interface{}) {
	c.data = toSet
}
