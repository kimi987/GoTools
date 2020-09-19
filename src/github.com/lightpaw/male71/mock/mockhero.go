package mock

import (
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/service/heromodule"
	"time"
)

var heroMap = map[int64]*entity.Hero{}

func init() {
	s := ifacemock.HeroDataService

	s.Mock(s.Create, func(a0 *entity.Hero) error {
		heroMap[a0.Id()] = a0
		return nil
	})

	s.Mock(s.Exist, func(a0 int64) (bool, error) {
		_, exist := heroMap[a0]
		return exist, nil
	})

	s.Mock(s.Func, func(a0 int64, a1 herolock.Func) {
		hero := heroMap[a0]
		if hero != nil {
			a1(hero, nil)
		}
	})

	s.Mock(s.FuncNotError, func(a0 int64, a1 herolock.FuncNotError) bool {
		hero := heroMap[a0]
		if hero != nil {
			a1(hero)
		}
		return false
	})

	s.Mock(s.FuncWithSend, func(a0 int64, a1 herolock.SendFunc) bool {
		hero := heroMap[a0]
		if hero != nil {
			a1(hero, LockResult)
		}
		return false
	})

	s.Mock(s.FuncWithSendError, func(a0 int64, a1 herolock.SendFunc) (bool, error) {
		hero := heroMap[a0]
		if hero != nil {
			a1(hero, LockResult)
		}
		return false, nil
	})
}

func DefaultHero(hero *entity.Hero) {
	SetHero(ifacemock.HeroController, hero)
}

func SetHero(hc *ifacemock.MockHeroController, hero *entity.Hero) {

	hc.Mock(hc.Id, func() int64 {
		return hero.Id()
	})

	hc.Mock(hc.IdBytes(), func() []byte {
		return hero.IdBytes()
	})

	hc.Mock(hc.FuncWithSend, func(a0 herolock.SendFunc) bool {
		a0(hero, LockResult)
		return false
	})

	hc.Mock(hc.Func, func(a0 herolock.Func) {
		a0(hero, nil)
	})

	hc.Mock(hc.FuncNotError, func(a0 herolock.FuncNotError) bool {
		a0(hero)
		return false
	})

	ifacemock.HeroDataService.Create(hero)

	LockResult.Reset()
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

func SetHeroNotEnoughCost(hero *entity.Hero, cost *resdata.Cost) {
	if cost.IsEmpty() {
		panic("mock.SetHeroNotEnoughCost cost.IsEmpty()")
	}

}

func SetHeroTroopCaptain(hero *entity.Hero, datas iface.ConfigDatas) {

	ctime := time.Now()

	// 解锁所有的武将
	result := newLockResult()
	for _, data := range datas.GetCaptainDataArray() {
		heromodule.TryAddCaptain(hero, result, data, ctime)

		if captain := hero.Military().Captain(data.Id); captain != nil {
			captain.SetSoldier(captain.SoldierCapcity())

			troop, index := hero.GetRecruitCaptainTroop()
			if troop != nil {
				troop.SetCaptainIfAbsent(index, captain, 0)
			}
		}
	}
}
