package dianquan

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/dianquan"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
)

 func NewDianquanModule(dep iface.ServiceDep, datas iface.ConfigDatas) *DianquanModule {
	m := &DianquanModule{}
	m.dep = dep
	m.datas = datas
	return m
}

//gogen:iface
type DianquanModule struct {
	dep iface.ServiceDep
	datas iface.ConfigDatas
}

//gogen:iface
func (m *DianquanModule) ProcessExchange(proto *dianquan.C2SExchangeProto, hc iface.HeroController) {
	times := proto.Times
	if times <= 0 {
		logrus.Debugf("点券兑换，times <=0 ")
		hc.Send(dianquan.ERR_EXCHANGE_FAIL_INVALID_TIMES)
		return
	}

	cost := m.datas.ExchangeMiscData().ExchangeBaseYuanbao * u64.FromInt32(times)
	toAdd := m.datas.ExchangeMiscData().ExchangeBaseDianquan * u64.FromInt32(times)
	if cost <= 0 || toAdd <= 0 {
		logrus.Debugf("点券兑换，实际的花费或获得<=0，times 太大了？")
		hc.Send(dianquan.ERR_EXCHANGE_FAIL_INVALID_TIMES)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.DianquanExchange)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !heromodule.ReduceYuanbao(hctx, hero, result, cost) {
			logrus.Debugf("点券兑换，元宝不够")
			result.Add(dianquan.ERR_EXCHANGE_FAIL_COST_NOT_ENOUGH)
			return
		}

		heromodule.AddDianquan(hctx, hero, result, toAdd)
		result.Add(dianquan.NewS2cExchangeMsg(times, u64.Int32(toAdd)))

		result.Changed()
		result.Ok()
	})
}


