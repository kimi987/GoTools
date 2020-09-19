package heromodule

import (
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"time"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/operate_type"
)

// 预订，就是先扣掉，失败再还回去

func ReserveCost(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, cost *resdata.Cost, ctime time.Time) (bool, *entity.ReserveResult) {

	if TryReduceCost(hctx, hero, result, cost) {
		return true, hero.Reservation().ReserveCost(cost, ctime)
	}
	return false, nil
}

func ReserveGoods(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *goods.GoodsData, goodsCount uint64, ctime time.Time) (bool, *entity.ReserveResult) {

	if TryReduceGoods(hctx, hero, result, goodsData, goodsCount) {
		return true, hero.Reservation().ReserveGoods(goodsData.Id, goodsCount, ctime)
	}
	return false, nil
}

func ReserveGoodsOrBuy(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *goods.GoodsData, goodsCount uint64, autoBuy bool, ctime time.Time) (bool, *entity.ReserveResult) {

	ok, toReduceGoodsCount, toReduceYinliang, toReduceDianquan, toReduceYuanbao := TryReduceOrBuyGoodsReturnPlan(hctx, hero, result, goodsData, goodsCount, autoBuy)
	if ok {
		return true, hero.Reservation().ReserveGoodsOrBuy(goodsData.Id, toReduceGoodsCount, toReduceYinliang, toReduceDianquan, toReduceYuanbao, ctime)
	}
	return false, nil
}

func ConfirmReserveResult(hctx *HeroContext, hc iface.HeroController, reserveResult *entity.ReserveResult, f func() (success bool)) {
	success := f()

	// 根据配置，可能不扣东西，此时reserveResult == nil

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if reserveResult != nil {
			ConfirmHeroReserveResult(hctx, hero, result, reserveResult, success)
			result.Ok()
		}
		if success {
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumMoveBase)
			UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_MOVE_BASE)
		}
	})

}

func ConfirmReserveResultWithHeroFunc(hctx *HeroContext, hc iface.HeroController, reserveResult *entity.ReserveResult, f func() (success bool, heroFunc herolock.SendFunc)) {
	success, heroFunc := f()

	// 根据配置，可能不扣东西，此时reserveResult == nil
	if reserveResult != nil {
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			ConfirmHeroReserveResult(hctx, hero, result, reserveResult, success)

			if heroFunc != nil {
				heroFunc(hero, result)
			}

			result.Ok()
		})
	}
}

func ReserveYuanbao(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, yuanbao uint64, ctime time.Time) (bool, *entity.ReserveResult) {
	if ReduceYuanbao(hctx, hero, result, yuanbao) {
		return true, hero.Reservation().ReserveYuanbao(yuanbao, ctime)
	}
	return false, nil
}

func ReserveDianquan(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, dianquan uint64, ctime time.Time) (bool, *entity.ReserveResult) {
	if ReduceDianquan(hctx, hero, result, dianquan) {
		return true, hero.Reservation().ReserveDianquan(dianquan, ctime)
	}
	return false, nil
}

func ReserveYinliang(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, yinliang uint64, ctime time.Time) (bool, *entity.ReserveResult) {
	if ReduceYinliang(hctx, hero, result, yinliang) {
		return true, hero.Reservation().ReserveYinliang(yinliang, ctime)
	}
	return false, nil
}

func ConfirmHeroReserveResult(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, reserveResult *entity.ReserveResult, success bool) {

	hero.Reservation().ConfirmResult(reserveResult)
	if !success {
		// 退回扣掉的东西
		addBackReserveResult(hctx, hero, result, reserveResult)
	}
}

func addBackReserveResult(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, r *entity.ReserveResult) {
	hctx = hctx.Copy(operate_type.HeroAddBackReservation)

	// 还资源
	gold, food, wood, stone := r.GetResource()
	AddUnsafeResource(hctx, hero, result, gold, food, wood, stone)

	// 还元宝
	AddYuanbao(hctx, hero, result, r.GetYuanbao())

	// 还点券
	AddDianquan(hctx, hero, result, r.GetDianquan())

	// 还银两
	AddYinliang(hctx, hero, result, r.GetYinliang())

	// 帮派贡献币
	AddGuildContributionCoin(hctx, hero, result, r.GetGuildContributionCoin())

	// 物品
	goodsIds, goodsCounts := r.GetGoodsIdCounts()
	AddGoodsIdArray(hctx, hero, result, goodsIds, goodsCounts)
}

// 定时检查，是不是要还钱
func tickAddBackReservation(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, ctime time.Time) {
	reserveResult := hero.Reservation().TryClearAndGetBackReserveResult(ctime)
	if reserveResult == nil {
		return
	}

	reserveResult.Print("超时退还预约消耗")

	addBackReserveResult(hctx, hero, result, reserveResult)
}
