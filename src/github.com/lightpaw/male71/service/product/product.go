package product

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/rpc7"
	"github.com/lightpaw/male7/pb/rpcpb/login2game"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/config/charge"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/ctxfunc"
	"context"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/gamelogs"
	"github.com/lightpaw/male7/constants"
)

func NewProductService(dep iface.ServiceDep, datas iface.ConfigDatas, dbService iface.DbService, time iface.TimeService, heroService iface.HeroDataService) *ProductService {
	s := &ProductService{
		dep:         dep,
		datas:       datas,
		dbService:   dbService,
		time:        time,
		heroService: heroService,
	}

	rpc7.Handle(login2game.NewBuyProductHandler(s.handleBuyProduct))

	return s
}

//gogen:iface
type ProductService struct {
	dep   iface.ServiceDep
	datas iface.ConfigDatas

	dbService   iface.DbService
	time        iface.TimeService
	heroService iface.HeroDataService
}

var (
	BuyProductSuccess       = &login2game.S2CBuyProductProto{Success: true}
	BuyProductUnkownProduct = &login2game.S2CBuyProductProto{Msg: "unkown product"}
	BuyProductHeroNotExist  = &login2game.S2CBuyProductProto{Msg: "hero not exist"}
	BuyProductInternalError = &login2game.S2CBuyProductProto{Msg: "internal error"}
)

func (s *ProductService) handleBuyProduct(r *login2game.C2SBuyProductProto) (*login2game.S2CBuyProductProto, error) {

	logrus.WithField("orderid", r.OrderId).
		WithField("orderamount", r.OrderAmount).
		WithField("ordertime", r.OrderTime).
		WithField("productid", r.ProductId).
		WithField("pid", r.Pid).
		WithField("sid", r.Sid).
		WithField("heroid", r.HeroId).Info("收到订单发货")

	// 几分钱
	orderAmountFen := r.OrderAmount

	// 1毛钱1元宝
	orderAmountYuanbao := orderAmountFen / 10

	productData := s.datas.GetProductData(r.ProductId)
	if productData == nil {
		logrus.WithField("id", r.ProductId).Error("处理商品发货，data == nil")
		return BuyProductUnkownProduct, nil
	}

	// 检查一下订单是否存在
	var orderExist bool
	if err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		orderExist, err = s.dbService.OrderExist(ctx, r.OrderId)
		return
	}); err != nil {
		logrus.WithError(err).Debug("处理商品发货，查询订单是否存在失败")
		return BuyProductInternalError, nil
	}

	if orderExist {
		// 说明这个订单处理过了，直接返回成功
		logrus.WithField("order", r.OrderId).Debug("处理商品发货，订单已经处理过了")
		return BuyProductSuccess, nil
	}

	exist, err := s.heroService.Exist(r.HeroId)
	if err != nil {
		logrus.WithError(err).Debug("处理商品发货，获取英雄是否存在失败")
		return BuyProductHeroNotExist, nil
	}

	if !exist {
		logrus.Debug("处理商品发货，英雄不存在")
		return BuyProductHeroNotExist, nil
	}

	ctime := s.time.CurrentTime()

	var payOrderFunc func() (*login2game.S2CBuyProductProto, error)
	switch data := productData.GetData().(type) {
	case *charge.ChargeObjData:

		payOrderFunc = func() (*login2game.S2CBuyProductProto, error) {
			// TODO r.sid 不同的server

			hctx := heromodule.NewContext(s.dep, operate_type.MiscActivateChargeObj)
			if s.heroService.FuncWithSend(r.HeroId, func(hero *entity.Hero, result herolock.LockResult) {
				isFirstCharge := hero.Misc().IsFirstChargedObj(data.Id)
				if isFirstCharge {
					hero.Misc().SetFirstChargedObj(data.Id)
					result.Add(misc.NewS2cUpdateFirstRechargeMsg(u64.Int32(data.Id)))
				}

				// 获取元宝
				rewardYuanbao := data.GetRewardYuanbao(isFirstCharge)
				heromodule.AddYuanbao(hctx, hero, result, rewardYuanbao)

				// 元宝赠送额度
				if toAddYuanbaoGift := s.dep.Datas().MiscGenConfig().YuanbaoGiftPercent.CalculateByPercent(rewardYuanbao); toAddYuanbaoGift > 0 {
					heromodule.AddYuanbaoGiftLimit(hctx, hero, result, toAddYuanbaoGift)
				}

				// 记录累计充值
				heromodule.AddRechargeAmount(hctx, hero, result, orderAmountYuanbao, ctime)

				result.Ok()
			}) {
				logrus.WithField("heroid", r.HeroId).Error("商品发货（充值），对英雄加锁失败")
				return BuyProductHeroNotExist, nil
			}

			// 记录充值日志
			gamelogs.RechargeLog(constants.PID, r.Sid, r.OrderId, r.ProductId, r.OrderAmount, r.HeroId)

			logrus.WithField("order", r.OrderId).Debug("处理商品发货(充值)，订单处理成功")
			return BuyProductSuccess, nil
		}
	default:
		logrus.WithField("id", productData.Id).Error("获取购买商品信息，该类型没有switch")
		return BuyProductUnkownProduct, nil
	}

	if payOrderFunc == nil {
		logrus.WithField("id", productData.Id).Error("获取购买商品信息，payOrderFunc == nil ")
		return BuyProductUnkownProduct, nil
	}

	// 尝试插入订单
	if err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		err = s.dbService.CreateOrder(ctx, r.OrderId, r.OrderAmount, r.OrderTime, r.Pid, r.Sid, r.HeroId, r.ProductId, ctime.Unix())
		return
	}); err != nil {
		logrus.WithField("order", r.OrderId).Error("处理商品发货，插入订单失败")
		return BuyProductInternalError, nil
	}

	return payOrderFunc()
}
