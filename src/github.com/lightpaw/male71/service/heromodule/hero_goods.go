package heromodule

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/config/goods"
	"time"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/pb/zhanjiang"
	"github.com/lightpaw/male7/gen/pb/tower"
	"github.com/lightpaw/male7/gen/pb/region"
)

type UseGoodsFunc func(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *goods.GoodsData, count uint64, ctime time.Time) (useCount uint64)

func GetUseEffectGoodsFunc(hctx *HeroContext, hero *entity.Hero, goodsData *goods.GoodsData, count uint64, ctime time.Time) (maxCanUseCount uint64, f UseGoodsFunc) {
	switch goodsData.EffectType {
	case shared_proto.GoodsEffectType_EFFECT_AMOUNT:

		switch goodsData.GoodsEffect.AmountType {
		case shared_proto.GoodsAmountType_AmountZhanjiang:
			// 过关斩将（只有次数消耗完了才可以使用，因此无论如何每次只会消耗1个物品）
			canAddTimes := hero.ZhanJiang().OpenTimes()
			maxUseCount := u64.DivideTimes(canAddTimes, goodsData.GoodsEffect.Amount)

			return u64.Min(count, maxUseCount), useZhanJiangTimesGoods
		case shared_proto.GoodsAmountType_AmountTower:
			// 千重楼（只有进行过扫荡才可以使用重置，因此无论如何每次只会消耗1个物品）
			if hero.Tower().CurrentFloor() >= hero.Tower().AutoMaxFloor() {
				return 1, useTowerResetGoods
			}
			return 0, useTowerResetGoods
		case shared_proto.GoodsAmountType_AmountJunTuan:
			// 加军团怪次数
			canAddTimes := u64.Sub(hero.GetJunTuanNpcTimes().MaxTimes(0), hero.GetJunTuanNpcTimes().Times(ctime, 0))
			maxUseCount := u64.DivideTimes(canAddTimes, goodsData.GoodsEffect.Amount)

			return u64.Min(count, maxUseCount), useJunTuanTimesGoods
		}

	case shared_proto.GoodsEffectType_EFFECT_RESOURCE:
		// 使用资源物品
		return count, useResourceGoods
	}

	return 0, nil
}

func useResourceGoods(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *goods.GoodsData, count uint64, ctime time.Time) (useCount uint64) {
	toAddGold := goodsData.GoodsEffect.Gold
	toAddFood := goodsData.GoodsEffect.Food
	toAddWood := goodsData.GoodsEffect.Wood
	toAddStone := goodsData.GoodsEffect.Stone

	AddSafeResource(hctx, hero, result, toAddGold*count, toAddFood*count, toAddWood*count, toAddStone*count)
	return count
}

func useZhanJiangTimesGoods(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *goods.GoodsData, count uint64, ctime time.Time) (useCount uint64) {
	canAddTimes := hero.ZhanJiang().OpenTimes()
	maxUseCount := u64.DivideTimes(canAddTimes, goodsData.GoodsEffect.Amount)
	useCount = u64.Min(count, maxUseCount)

	//if hero.ZhanJiang().OpenTimes() < hctx.Datas().ZhanJiangMiscData().DefaultTimes {
	//	result.Add(depot.ERR_USE_GOODS_FAIL_OPEN_TIMES_NOT_ZERO)
	//	return
	//}
	if useCount > 0 {
		hero.ZhanJiang().ReduceOpenTimes(goodsData.GoodsEffect.Amount * useCount)
		result.Add(zhanjiang.NewS2cUpdateOpenTimesMsg(u64.Int32(hero.ZhanJiang().OpenTimes())))
	}
	return useCount
}

func useTowerResetGoods(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *goods.GoodsData, count uint64, ctime time.Time) (useCount uint64) {
	//if hero.Tower().CurrentFloor() < hero.Tower().AutoMaxFloor() {
	//	result.Add(depot.ERR_USE_GOODS_FAIL_FLOOR_LESS_THAN_MAX)
	//	return
	//}

	if hero.Tower().CurrentFloor() >= hero.Tower().AutoMaxFloor() {
		hero.Tower().SetCurrentFloorToZero()
		result.Add(tower.NewS2cUpdateCurrentFloorMsg(u64.Int32(hero.Tower().CurrentFloor())))
		return 1
	}

	return 0
}

func useJunTuanTimesGoods(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *goods.GoodsData, count uint64, ctime time.Time) (useCount uint64) {
	canAddTimes := u64.Sub(hero.GetJunTuanNpcTimes().MaxTimes(0), hero.GetJunTuanNpcTimes().Times(ctime, 0))
	maxUseCount := u64.DivideTimes(canAddTimes, goodsData.GoodsEffect.Amount)
	useCount = u64.Min(count, maxUseCount)

	if useCount > 0 {
		hero.GetJunTuanNpcTimes().AddTimes(goodsData.GoodsEffect.Amount*useCount, ctime, 0)
		result.Add(region.NewS2cUpdateJunTuanNpcTimesMsg(hero.GetJunTuanNpcTimes().StartTimeUnix32(), nil))
	}
	return useCount
}
