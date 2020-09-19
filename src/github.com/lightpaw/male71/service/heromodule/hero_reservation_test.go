package heromodule

import (
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/pbutil"
	. "github.com/onsi/gomega"
	"testing"
	"time"
	"github.com/lightpaw/male7/service/operate_type"
)

func TestReservation(t *testing.T) {
	RegisterTestingT(t)

	ctime := time.Now()
	ifacemock.TimeService.Mock(ifacemock.TimeService.CurrentTime, func() time.Time {
		return ctime
	})

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	Ω(err).Should(Succeed())

	testGoodsReservation(datas, datas.GoodsData().MinKeyData, false)
	testGoodsReservation(datas, datas.GoodsData().MaxKeyData, true)

	cost := &resdata.Cost{}
	cost.Gold = 100
	//cost.Food = 101
	//cost.Wood = 102
	cost.Stone = 103
	cost.Yuanbao = 104
	cost.GuildContributionCoin = 105
	cost.Dianquan = 106
	cost.Yinliang = 107
	cost.Goods = []*goods.GoodsData{datas.GoodsData().MinKeyData, datas.GoodsData().MaxKeyData}
	cost.GoodsCount = []uint64{2, 6}

	testCostReservation(datas, cost, false)
	testCostReservation(datas, cost, true)

	testReservationExpired(datas, cost)
}

func testGoodsReservation(datas iface.ConfigDatas, data *goods.GoodsData, success bool) {
	dep := ifacemock.ServiceDep
	dep.Mock(dep.Tlog, func() iface.TlogService {
		return ifacemock.TlogService
	})
	dep.Mock(dep.Datas, func() iface.ConfigDatas {
		return datas
	})
	hctx := NewContext(dep, operate_type.TlogIgnore)

	ctime := ifacemock.TimeService.CurrentTime()
	hero := entity.NewHero(1, "hero1", datas.HeroInitData(), ctime)

	// 没有物品可以预约
	Ω(hero.Depot().HasEnoughGoods(data.Id, 1)).Should(BeFalse())
	hasEnoughGoods, reserveResult := ReserveGoods(hctx, hero, result, data, 1, ctime)
	Ω(hasEnoughGoods).Should(BeFalse())
	Ω(reserveResult).Should(BeNil())

	// 有物品可以预约
	hero.Depot().AddGoods(data.Id, 1)
	Ω(hero.Depot().HasEnoughGoods(data.Id, 1)).Should(BeTrue())

	Ω(hero.Depot().HasEnoughGoods(data.Id, 1)).Should(BeTrue())
	hasEnoughGoods, reserveResult = ReserveGoods(hctx, hero, result, data, 1, ctime)
	Ω(hasEnoughGoods).Should(BeTrue())
	Ω(reserveResult).ShouldNot(BeNil())

	goodsIds, goodsCounts := reserveResult.GetGoodsIdCounts()
	Ω(goodsIds).Should(Equal([]uint64{data.Id}))
	Ω(goodsCounts).Should(Equal([]uint64{1}))

	Ω(hero.Depot().HasEnoughGoods(data.Id, 1)).Should(BeFalse())
	Ω(hero.Reservation().IsEmpty()).Should(BeFalse())

	// confirm success = false
	ConfirmHeroReserveResult(hctx, hero, result, reserveResult, success)

	Ω(hero.Depot().HasEnoughGoods(data.Id, 1)).Should(Equal(!success))
	Ω(hero.Reservation().IsEmpty()).Should(BeTrue())

}

func testReservationExpired(datas iface.ConfigDatas, cost *resdata.Cost) {
	dep := ifacemock.ServiceDep
	dep.Mock(dep.Tlog, func() iface.TlogService {
		return ifacemock.TlogService
	})
	hctx := NewContext(dep, operate_type.TlogIgnore)

	ctime := ifacemock.TimeService.CurrentTime()
	hero := entity.NewHero(1, "hero1", datas.HeroInitData(), ctime)

	// 有物品可以预约
	SetHeroEnoughCost(hero, cost)
	testHeroHasEnoughCost(hero, cost)
	Ω(hero.Reservation().IsEmpty()).Should(BeTrue())

	hasEnoughCost, reserveResult := ReserveCost(hctx, hero, result, cost, ctime)
	Ω(hasEnoughCost).Should(BeTrue())
	Ω(reserveResult).ShouldNot(BeNil())

	Ω(hero.Reservation().IsEmpty()).Should(BeFalse())

	// 东西都被扣掉了
	testHeroZeroCost(hero, cost)

	tickAddBackReservation(hctx, hero, result, ctime)
	testHeroZeroCost(hero, cost)
	Ω(hero.Reservation().IsEmpty()).Should(BeFalse())

	tickAddBackReservation(hctx, hero, result, ctime.Add(time.Hour).Add(-1))
	testHeroZeroCost(hero, cost)
	Ω(hero.Reservation().IsEmpty()).Should(BeFalse())

	// 时间到，加回来
	tickAddBackReservation(hctx, hero, result, ctime.Add(time.Hour))
	testHeroHasEnoughCost(hero, cost)
	Ω(hero.Reservation().IsEmpty()).Should(BeTrue())
}

func testCostReservation(datas iface.ConfigDatas, cost *resdata.Cost, success bool) {
	dep := ifacemock.ServiceDep
	dep.Mock(dep.Tlog, func() iface.TlogService {
		return ifacemock.TlogService
	})
	hctx := NewContext(dep, operate_type.TlogIgnore)

	ctime := ifacemock.TimeService.CurrentTime()
	hero := entity.NewHero(1, "hero1", datas.HeroInitData(), ctime)

	// 没有物品可以预约
	hasEnoughCost, reserveResult := ReserveCost(hctx, hero, result, cost, ctime)
	Ω(hasEnoughCost).Should(BeFalse())
	Ω(reserveResult).Should(BeNil())

	// 设置物品，达到物品条件
	SetHeroEnoughCost(hero, cost)
	testHeroHasEnoughCost(hero, cost)

	// 有物品可以预约
	hasEnoughCost, reserveResult = ReserveCost(hctx, hero, result, cost, ctime)
	Ω(hasEnoughCost).Should(BeTrue())
	Ω(reserveResult).ShouldNot(BeNil())

	gold, food, wood, stone := reserveResult.GetResource()
	Ω(gold).Should(Equal(cost.Gold))
	Ω(food).Should(Equal(cost.Food))
	Ω(wood).Should(Equal(cost.Wood))
	Ω(stone).Should(Equal(cost.Stone))

	Ω(reserveResult.GetYuanbao()).Should(Equal(cost.Yuanbao))
	Ω(reserveResult.GetDianquan()).Should(Equal(cost.Dianquan))
	Ω(reserveResult.GetYinliang()).Should(Equal(cost.Yinliang))
	Ω(reserveResult.GetGuildContributionCoin()).Should(Equal(cost.GuildContributionCoin))

	goodsIds, goodsCounts := reserveResult.GetGoodsIdCounts()
	Ω(goodsIds).Should(Equal(goods.GetGoodsDataKeyArray(cost.Goods)))
	Ω(goodsCounts).Should(Equal(cost.GoodsCount))

	// 东西都被扣掉了
	testHeroZeroCost(hero, cost)

	Ω(hero.Reservation().IsEmpty()).Should(BeFalse())

	// confirm success = false
	ConfirmHeroReserveResult(hctx, hero, result, reserveResult, success)

	if success {
		testHeroZeroCost(hero, cost)
	} else {
		testHeroHasEnoughCost(hero, cost)
	}

	Ω(hero.Reservation().IsEmpty()).Should(BeTrue())

}

func testHeroHasEnoughCost(hero *entity.Hero, cost *resdata.Cost) {
	Ω(hero.GetUnsafeResource().Gold()).Should(Equal(cost.Gold))
	Ω(hero.GetUnsafeResource().Food()).Should(Equal(cost.Food))
	Ω(hero.GetUnsafeResource().Wood()).Should(Equal(cost.Wood))
	Ω(hero.GetUnsafeResource().Stone()).Should(Equal(cost.Stone))

	Ω(hero.GetYuanbao()).Should(Equal(cost.Yuanbao))
	Ω(hero.GetDianquan()).Should(Equal(cost.Dianquan))
	Ω(hero.GetYinliang()).Should(Equal(cost.Yinliang))
	Ω(hero.GetGuildContributionCoin()).Should(Equal(cost.GuildContributionCoin))

	for i, v := range cost.Goods {
		Ω(hero.Depot().GetGoodsCount(v.Id)).Should(Equal(cost.GoodsCount[i]))
	}
}

func testHeroZeroCost(hero *entity.Hero, cost *resdata.Cost) {
	Ω(hero.GetUnsafeResource().Gold()).Should(Equal(uint64(0)))
	Ω(hero.GetUnsafeResource().Food()).Should(Equal(uint64(0)))
	Ω(hero.GetUnsafeResource().Wood()).Should(Equal(uint64(0)))
	Ω(hero.GetUnsafeResource().Stone()).Should(Equal(uint64(0)))

	Ω(hero.GetYuanbao()).Should(Equal(uint64(0)))
	Ω(hero.GetDianquan()).Should(Equal(uint64(0)))
	Ω(hero.GetGuildContributionCoin()).Should(Equal(uint64(0)))

	for _, v := range cost.Goods {
		Ω(hero.Depot().GetGoodsCount(v.Id)).Should(Equal(uint64(0)))
	}
}

func SetHeroEnoughCost(hero *entity.Hero, cost *resdata.Cost) {

	if cost.IsEmpty() {
		return
	}

	hero.GetUnsafeResource().AddResource(cost.GetResource())
	hero.AddYuanbao(cost.Yuanbao)
	hero.AddDianquan(cost.Dianquan)
	hero.AddYinliang(cost.Yinliang)
	hero.AddGuildContributionCoin(cost.GuildContributionCoin)
	for i, g := range cost.Goods {
		hero.Depot().AddGoods(g.Id, cost.GoodsCount[i])
	}
}

var result = &s{}

type s struct {
}

func (s *s) Add(pbutil.Buffer) {
}

func (s *s) AddBroadcast(pbutil.Buffer) {
}

func (s *s) AddTlog(func()) {
}

func (s *s) AddFunc(func() pbutil.Buffer) {
}

func (s *s) Changed() {
}

func (s *s) Ok() {
}
