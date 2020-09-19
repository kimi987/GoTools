package charge

import "github.com/lightpaw/male7/config/resdata"

//gogen:config
type ChargePrizeData struct {
	_ struct{} `file:"充值/充值奖励.txt"`
	_ struct{} `protogen:"true"`

	Id     uint64
	Name   string
	Image  string `desc:"立绘"`
	Amount uint64 `desc:"充值额度，达到才有奖励"`
	Value  uint64 `desc:"价值(元宝)，前端展示用"`
	Prize  *resdata.Prize
	Desc   string
}

//gogen:config
type ChargeObjData struct {
	_ struct{} `file:"充值/充值项.txt"`
	_ struct{} `protogen:"true"`

	Id                 uint64
	Name               string
	Icon               string
	Image              string
	Product            *ProductData `default:"nullable" protofield:"ProductId,int32(%s.Id),int32"`
	ChargeAmount       uint64       `desc:"充值金额"`
	Yuanbao            uint64       `desc:"充值获得元宝"`
	YuanbaoAddition    uint64       `desc:"充值附赠元宝" validator:"uint" default:"0"`
	FirstChargeYuanbao uint64       `desc:"首充附赠元宝" validator:"uint" default:"0"`
	VipExp             uint64       `desc:"vip经验增加值"`
}

func (d *ChargeObjData) Init(filename string) {
	if d.Product != nil {
		d.Product.SetData(d)
	}
}

func (d *ChargeObjData) GetRewardYuanbao(firstCharge bool) uint64 {
	yuanbao := d.Yuanbao
	if firstCharge {
		yuanbao += d.FirstChargeYuanbao
	} else {
		yuanbao += d.YuanbaoAddition
	}
	return yuanbao
}
