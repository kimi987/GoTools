package realm

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/util/i64"
	. "github.com/onsi/gomega"
	"testing"
	"github.com/lightpaw/male7/util/u64"
	"math"
	"fmt"
)

func TestFirstCalculate(t *testing.T) {
	RegisterTestingT(t)
	r, _ := newMockRealm()

	transformResource := CalculateFirstRobTransformResource

	var soldier uint64 = 1000

	heroRes := &entity.ResourceStorage{}
	heroRes.AddResource(0, 10000, 20000, 100000, )

	var weakCoef float64 = 0.5

	// 还可以抢资源
	lostGold := u64.Min(u64.Sub(heroRes.Gold(), 10000),
		transformResource(soldier, 1, heroRes.Gold(), 10000, weakCoef))
	lostFood := u64.Min(u64.Sub(heroRes.Food(), 10000),
		transformResource(soldier, 1, heroRes.Food(), 10000, weakCoef))
	lostWood := u64.Min(u64.Sub(heroRes.Wood(), 10000),
		transformResource(soldier, 1, heroRes.Wood(), 10000, weakCoef))
	lostStone := u64.Min(u64.Sub(heroRes.Stone(), 10000),
		transformResource(soldier, 1, heroRes.Stone(), 10000, weakCoef))

	fmt.Println(lostGold, lostFood, lostWood, lostStone)

	coef := math.Min(r.config().RobberCoef, 1)

	toAddGold := u64.MultiF64(lostGold, coef)
	toAddFood := u64.MultiF64(lostFood, coef)
	toAddWood := u64.MultiF64(lostWood, coef)
	toAddStone := u64.MultiF64(lostStone, coef)

	fmt.Println(toAddGold, toAddFood, toAddWood, toAddStone)

}

func TestRecurringCalculate(t *testing.T) {
	RegisterTestingT(t)
	r, _ := newMockRealm()

	transformResource := CalculateRecurringRobTransformResource

	var soldier uint64 = 1000

	heroRes := &entity.ResourceStorage{}
	heroRes.AddResource(100000, 100000, 100000, 100000, )

	//var weakCoef float64 = 0.5

	var defProsperityCapcity uint64 = 51000
	var actProsperityCapcity uint64 = 31000

	weakCoef := math.Atan2(
		u64.Sub2Float64(defProsperityCapcity, actProsperityCapcity),
		float64(defProsperityCapcity+actProsperityCapcity),
	) / math.Pi + 1

	// 还可以抢资源
	lostGold := u64.Min(u64.Sub(heroRes.Gold(), 10000),
		transformResource(soldier, 1, heroRes.Gold(), 10000, weakCoef))
	lostFood := u64.Min(u64.Sub(heroRes.Food(), 10000),
		transformResource(soldier, 1, heroRes.Food(), 10000, weakCoef))
	lostWood := u64.Min(u64.Sub(heroRes.Wood(), 10000),
		transformResource(soldier, 1, heroRes.Wood(), 10000, weakCoef))
	lostStone := u64.Min(u64.Sub(heroRes.Stone(), 10000),
		transformResource(soldier, 1, heroRes.Stone(), 10000, weakCoef))

	coef := math.Min(r.config().RobberCoef, 1)

	toAddGold := u64.MultiF64(lostGold, coef)
	toAddFood := u64.MultiF64(lostFood, coef)
	toAddWood := u64.MultiF64(lostWood, coef)
	toAddStone := u64.MultiF64(lostStone, coef)

	fmt.Println(toAddGold, toAddFood, toAddWood, toAddStone)

}

//func TestResourcePointConflict5_6(t *testing.T) {
//	RegisterTestingT(t)
//
//	r, _ := newMockRealm()
//	ctime := r.services.timeService.CurrentTime()
//
//	// 第一个英雄
//	hero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
//	mock.DefaultHero(hero)
//	hero.SetBaseLevel(5)
//	hero.SetProsperity(100)
//
//	x, y := 75, 69
//	ok := r.ReservePos(x, y)
//	Ω(ok).Should(BeTrue())
//	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())
//
//	processed, err := r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeTransfer)
//	Ω(processed).Should(BeTrue())
//	Ω(err).Should(Succeed())
//
//	base := r.getBase(hero.Id())
//	Ω(base).ShouldNot(BeNil())
//
//	checkResourcePointConflict(r, base, hero, 0)
//
//	// 第二个英雄
//	hero2 := entity.NewHero(2, "hero2", r.services.datas.HeroInitData(), ctime)
//	mock.DefaultHero(hero2)
//	hero2.SetBaseLevel(6)
//	hero2.SetProsperity(100)
//
//	//x, y = 90, 80
//
//	x, y = 81, 70
//	ok = r.ReservePos(x, y)
//	Ω(ok).Should(BeTrue())
//	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())
//
//	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeTransfer)
//	Ω(processed).Should(BeTrue())
//	Ω(err).Should(Succeed())
//
//	base2 := r.getBase(hero2.Id())
//	Ω(base2).ShouldNot(BeNil())
//
//	conflictCube1 := cb.XYCube(78, 69)
//	conflictCube2 := cb.XYCube(78, 70)
//
//	Ω(r.resourceConflictHeroMap[conflictCube1]).Should(ConsistOf([]int64{1, 2}))
//	Ω(r.resourceConflictHeroMap[conflictCube2]).Should(ConsistOf([]int64{1, 2}))
//
//	checkPointConflict(r, hero, 78, 69, true)
//	checkPointConflict(r, hero, 78, 70, true)
//
//	checkPointConflict(r, hero2, 78, 69, true)
//	checkPointConflict(r, hero2, 78, 70, true)
//
//}
//
//func TestResourcePointConflict5_6Move(t *testing.T) {
//	RegisterTestingT(t)
//
//	r, _ := newMockRealm()
//	ctime := r.services.timeService.CurrentTime()
//
//	// 第一个英雄
//	hero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
//	mock.DefaultHero(hero)
//	hero.SetBaseLevel(5)
//	hero.SetProsperity(100)
//
//	x, y := 75, 69
//	ok := r.ReservePos(x, y)
//	Ω(ok).Should(BeTrue())
//	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())
//
//	processed, err := r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeTransfer)
//	Ω(processed).Should(BeTrue())
//	Ω(err).Should(Succeed())
//
//	base := r.getBase(hero.Id())
//	Ω(base).ShouldNot(BeNil())
//
//	checkResourcePointConflict(r, base, hero, 0)
//
//	// 第二个英雄
//	hero2 := entity.NewHero(2, "hero2", r.services.datas.HeroInitData(), ctime)
//	mock.DefaultHero(hero2)
//	hero2.SetBaseLevel(6)
//	hero2.SetProsperity(100)
//
//	//_, x, y = r.randomHomePosAndAdd()
//	//fmt.Println(x, y)
//
//	x, y = 75, 68
//	ok = r.ReservePos(x, y)
//	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())
//	Ω(ok).Should(BeFalse())
//
//	x, y = 59, 65
//	ok = r.ReservePos(x, y)
//	Ω(ok).Should(BeTrue())
//	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())
//
//	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeTransfer)
//	Ω(processed).Should(BeTrue())
//	Ω(err).Should(Succeed())
//
//	base2 := r.getBase(hero2.Id())
//	Ω(base2).ShouldNot(BeNil())
//
//	// 移动到这个位置
//	oldX, oldY := x, y
//	x, y = 60, 65
//	ok = r.ReservePosForMoveBase(oldX, oldY, x, y)
//	Ω(ok).Should(BeTrue())
//	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())
//
//	processed, err = r.MoveBase(ifacemock.HeroController, oldX, oldY, x, y)
//	Ω(processed).Should(BeTrue())
//	Ω(err).Should(Succeed())
//
//	base2 = r.getBase(hero2.Id())
//	Ω(base2).ShouldNot(BeNil())
//
//	// 移动到这个位置
//	oldX, oldY = x, y
//	x, y = 81, 70
//	ok = r.ReservePosForMoveBase(oldX, oldY, x, y)
//	Ω(ok).Should(BeTrue())
//	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())
//
//	processed, err = r.MoveBase(ifacemock.HeroController, oldX, oldY, x, y)
//	Ω(processed).Should(BeTrue())
//	Ω(err).Should(Succeed())
//
//	base2 = r.getBase(hero2.Id())
//	Ω(base2).ShouldNot(BeNil())
//
//	conflictCube1 := cb.XYCube(78, 69)
//	conflictCube2 := cb.XYCube(78, 70)
//
//	Ω(r.resourceConflictHeroMap[conflictCube1]).Should(ConsistOf([]int64{1, 2}))
//	Ω(r.resourceConflictHeroMap[conflictCube2]).Should(ConsistOf([]int64{1, 2}))
//
//	checkPointConflict(r, hero, 78, 69, true)
//	checkPointConflict(r, hero, 78, 70, true)
//
//	checkPointConflict(r, hero2, 78, 69, true)
//	checkPointConflict(r, hero2, 78, 70, true)
//
//}

func checkPointConflict(r *Realm, hero *entity.Hero, x, y int, conflict bool) {
	offset := hexagon.EvenOffsetBetween(hero.BaseX(), hero.BaseY(), x, y)
	layoutData := r.services.datas.RegionConfig().GetLayoutDataByEvenOffset(offset)
	Ω(layoutData).ShouldNot(BeNil())
	Ω(hero.IsHeroConflictResourcePoint(layoutData)).Should(Equal(conflict))
}

func TestResourcePoint(t *testing.T) {
	RegisterTestingT(t)

	r, _ := newMockRealm()
	ctime := r.services.timeService.CurrentTime()

	// 第一个英雄
	hero := entity.NewHero(1, "hero1", r.services.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	hero.SetBaseLevel(1)

	ok, x, y := r.randomBasePos()
	Ω(ok).Should(BeTrue())
	Ω(r.mapData.IsValidHomePosition(x, y)).Should(BeTrue())

	processed, err := r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeNewHero)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	base := r.getBase(hero.Id())
	Ω(base).ShouldNot(BeNil())

	checkResourcePointConflict(r, base, hero, 0)

	processed, err, originX, originY := r.RemoveBase(hero.Id(), true, r.getTextHelp().MRDRAttMoveBase4a.Text, r.getTextHelp().MRDRAttMoveBase4d.Text)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())
	Ω(x).Should(Equal(originX))
	Ω(y).Should(Equal(originY))

	Ω(r.resourceConflictHeroMap).Should(BeEmpty())

	// 重新加入
	processed, err = r.AddBase(ifacemock.HeroController.Id(), x, y, realmface.AddBaseHomeTransfer)
	Ω(processed).Should(BeTrue())
	Ω(err).Should(Succeed())

	checkResourcePointConflict(r, base, hero, 0)

	// 第二个英雄
	hero2 := entity.NewHero(2, "hero2", r.services.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero2)

	hero2.SetBaseLevel(2)
	hero2.SetBaseXY(r.Id(), x, y)
	base2 := r.newBase(hero2)
	setBaseForTest(r, base2)
	r.addHomeResourcePointBlock(base2,true)

	checkResourcePointConflict(r, base, hero, 1)
	checkResourcePointConflict(r, base2, hero2, 1)
}

func checkNpcResourceConflict(r *Realm, base *baseWithData, hero *entity.Hero) {

	Ω(base.BaseX()).Should(Equal(hero.BaseX()))
	Ω(base.BaseY()).Should(Equal(hero.BaseY()))

	// 看等级
	for _, v := range r.services.datas.GetHomeNpcBaseDataArray() {
		if v.HomeBaseLevel > 0 && v.HomeBaseLevel <= base.BaseLevel() {

			// 这里应该是有一个npc行营的
			npcBasePos := hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), v.EvenOffsetX, v.EvenOffsetY)
			Ω(r.conflict.baseConflictCount[npcBasePos] <= 1).Should(BeTrue())

			offset := cb.XYCube(v.EvenOffsetX, v.EvenOffsetY)
			layoutData := r.services.datas.RegionConfig().GetLayoutDataByEvenOffset(offset)
			if layoutData == nil {
				Ω(r.resourceConflictHeroMap[npcBasePos]).Should(BeEmpty())
				continue
			}

			Ω(hero.IsNpcConflictResourcePoint(layoutData)).Should(BeTrue())

			if layoutData.RequireBaseLevel <= base.BaseLevel() {
				Ω(r.resourceConflictHeroMap[npcBasePos]).Should(Equal([]int64{hero.Id()}))
			} else {
				Ω(r.resourceConflictHeroMap[npcBasePos]).Should(BeEmpty())
			}

		}
	}

	home := GetHomeBase(base)

	count := 1
	hero.WalkHomeNpcBase(func(heroNpcBase *entity.HomeNpcBase) bool {

		//count++

		npcBasePos := hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), heroNpcBase.GetData().EvenOffsetX, heroNpcBase.GetData().EvenOffsetY)

		//npcBase := r.getBase(heroNpcBase.Id())
		npcBase := home.homeNpcBase[heroNpcBase.Id()]

		Ω(npcBase).ShouldNot(BeNil())
		x, y := npcBasePos.XY()
		Ω(npcBase.BaseX()).Should(Equal(x))
		Ω(npcBase.BaseY()).Should(Equal(y))

		homeNpc := GetHomeNpcBase(npcBase)
		Ω(homeNpc.ownerHeroId).Should(Equal(hero.Id()))

		return false
	})

	Ω(baseLen(r)).Should(Equal(count))
}

func baseLen(r *Realm) int {
	count := 0
	r.rangeBases(func(base *baseWithData) (toContinue bool) {
		count++
		return true
	})
	return count
}

func checkResourcePointConflict(r *Realm, base *baseWithData, hero *entity.Hero, conflictBaseLevel uint64) {

	if base.BaseLevel() <= 0 {
		return
	}

	for _, targetCube := range hexagon.SpiralRing(base.BaseX(), base.BaseY(), 1) {
		Ω(i64.Contains(r.resourceConflictHeroMap[targetCube], base.Id())).Should(BeTrue())
	}

	for _, layoutData := range r.services.datas.GetBuildingLayoutDataArray() {
		own := base.BaseLevel() >= layoutData.RequireBaseLevel

		targetCube := hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), layoutData.RegionOffsetX, layoutData.RegionOffsetY)
		Ω(i64.Contains(r.resourceConflictHeroMap[targetCube], base.Id())).Should(Equal(own))

		Ω(hero.IsHeroConflictResourcePoint(layoutData)).Should(Equal(layoutData.RequireBaseLevel <= conflictBaseLevel))
	}

}
