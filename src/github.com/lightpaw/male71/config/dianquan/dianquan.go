package dianquan

import "github.com/lightpaw/male7/util/check"

//gogen:config
type ExchangeMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"点券/点券杂项.txt"`
	_ struct{} `proto:"shared_proto.DianquanMiscProto"`
	_ struct{} `protoconfig:"dianquan_misc"`

	ExchangeBaseYuanbao  uint64 `validator:"int>0"` // 兑换点券，单次元宝花费
	ExchangeBaseDianquan uint64 `validator:"int>0"` // 兑换点券，单次点券获得
}

func (d *ExchangeMiscData) Init(filename string) {
	check.PanicNotTrue(d.ExchangeBaseYuanbao <= d.ExchangeBaseDianquan, "%v，花费的元宝数不能大于获得的点券数 yuanbao:%v  dianquan:%v", filename, d.ExchangeBaseYuanbao, d.ExchangeBaseDianquan)
}
