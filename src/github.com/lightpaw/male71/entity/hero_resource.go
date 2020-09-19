package entity

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

type hero_resource struct {
	unsafeResource *ResourceStorage // 非安全资源
	safeResource   *ResourceStorage // 安全资源（不会被抢的部分 ）

	jade            uint64 // 玉璧
	jadeOre         uint64 // 玉石矿，定时会转化成玉璧
	historyJade     uint64 // 历史获得的玉璧总数
	todayObtainJade uint64 // 今日获得的玉璧总数

	nextResDecayTime time.Time // 下次超出保护上限的资源衰减时间
}

func (hero *hero_resource) GetGold() uint64 {
	return hero.unsafeResource.gold + hero.safeResource.gold
}

func (hero *hero_resource) GetStone() uint64 {
	return hero.unsafeResource.stone + hero.safeResource.stone
}

func (hero *hero_resource) GetUnsafeResource() *ResourceStorage {
	return hero.unsafeResource
}

func (hero *hero_resource) GetSafeResource() *ResourceStorage {
	return hero.safeResource
}

type ResourceStorage struct {
	gold  uint64
	food  uint64
	wood  uint64
	stone uint64

	isSafeResource bool
}

func (r *ResourceStorage) IsSafeResource() bool {
	return r.isSafeResource
}

func (hero *ResourceStorage) AddResource(toAddGold, toAddFood, toAddWood, toAddStone uint64) (newGold, newFood, newWood, newStone uint64) {
	newGold = hero.AddGold(toAddGold)
	newFood = hero.AddFood(toAddFood)
	newWood = hero.AddWood(toAddWood)
	newStone = hero.AddStone(toAddStone)
	return
}

func (hero *ResourceStorage) ReduceResource(toReduceGold, toReduceFood, toReduceWood, toReduceStone uint64) (newGold, newFood, newWood, newStone uint64) {
	newGold = hero.ReduceGold(toReduceGold)
	newFood = hero.ReduceFood(toReduceFood)
	newWood = hero.ReduceWood(toReduceWood)
	newStone = hero.ReduceStone(toReduceStone)
	return
}

func (hero *ResourceStorage) GetRes(resType shared_proto.ResType) uint64 {
	switch resType {
	case shared_proto.ResType_GOLD:
		return hero.Gold()
	case shared_proto.ResType_FOOD:
		return hero.Food()
	case shared_proto.ResType_WOOD:
		return hero.Wood()
	case shared_proto.ResType_STONE:
		return hero.Stone()
	default:
		return 0
	}
}

func (hero *ResourceStorage) AddRes(resType shared_proto.ResType, toAdd uint64) uint64 {
	switch resType {
	case shared_proto.ResType_GOLD:
		return hero.AddGold(toAdd)
	case shared_proto.ResType_FOOD:
		return hero.AddFood(toAdd)
	case shared_proto.ResType_WOOD:
		return hero.AddWood(toAdd)
	case shared_proto.ResType_STONE:
		return hero.AddStone(toAdd)
	default:
		return 0
	}
}

func (hero *ResourceStorage) HasEnoughRes(resType shared_proto.ResType, toReduce uint64) bool {
	switch resType {
	case shared_proto.ResType_GOLD:
		return hero.gold >= toReduce
	case shared_proto.ResType_FOOD:
		return hero.food >= toReduce
	case shared_proto.ResType_WOOD:
		return hero.wood >= toReduce
	case shared_proto.ResType_STONE:
		return hero.stone >= toReduce
	default:
		return false
	}
}

func (hero *ResourceStorage) ReduceRes(resType shared_proto.ResType, toReduce uint64) uint64 {
	switch resType {
	case shared_proto.ResType_GOLD:
		return hero.ReduceGold(toReduce)
	case shared_proto.ResType_FOOD:
		return hero.ReduceFood(toReduce)
	case shared_proto.ResType_WOOD:
		return hero.ReduceWood(toReduce)
	case shared_proto.ResType_STONE:
		return hero.ReduceStone(toReduce)
	default:
		return 0
	}
}

func (hero *ResourceStorage) Gold() uint64 {
	return hero.gold
}

func (hero *ResourceStorage) AddGold(toAdd uint64) uint64 {
	hero.gold += toAdd
	return hero.gold
}

func (hero *ResourceStorage) ReduceGold(toReduce uint64) uint64 {
	hero.gold = u64.Sub(hero.gold, toReduce)
	return hero.gold
}

func (hero *ResourceStorage) Food() uint64 {
	return hero.food
}

func (hero *ResourceStorage) AddFood(toAdd uint64) uint64 {
	//if toAdd >= 1000000{
	//	logrus.WithField("Food", string(debug.Stack())).Error("单次添加资源，超过百万，预警，可能是大bug")
	//}
	//hero.food += toAdd
	return hero.food
}

func (hero *ResourceStorage) ReduceFood(toReduce uint64) uint64 {
	hero.food = u64.Sub(hero.food, toReduce)
	return hero.food
}

func (hero *ResourceStorage) Wood() uint64 {
	return hero.wood
}

func (hero *ResourceStorage) AddWood(toAdd uint64) uint64 {
	//if toAdd >= 1000000{
	//	logrus.WithField("Wood", string(debug.Stack())).Error("单次添加资源，超过百万，预警，可能是大bug")
	//}
	//hero.wood += toAdd
	return hero.wood
}

func (hero *ResourceStorage) ReduceWood(toReduce uint64) uint64 {
	hero.wood = u64.Sub(hero.wood, toReduce)
	return hero.wood
}

func (hero *ResourceStorage) Stone() uint64 {
	return hero.stone
}

func (hero *ResourceStorage) AddStone(toAdd uint64) uint64 {
	hero.stone += toAdd
	return hero.stone
}

func (hero *ResourceStorage) ReduceStone(toReduce uint64) uint64 {
	hero.stone = u64.Sub(hero.stone, toReduce)
	return hero.stone
}

func (hero *Hero) ProtectedCapcity() uint64 {
	return hero.buildingEffect.protectedCapcity
}

func (d *hero_resource) HasEnoughResource(gold, food, wood, stone uint64) bool {
	unsafe := d.GetUnsafeResource()
	safe := d.GetSafeResource()

	if unsafe.Gold() < gold {
		gold -= unsafe.Gold()
		if safe.Gold() < gold {
			return false
		}
	}

	if unsafe.Food() < food {
		food -= unsafe.Food()
		if safe.Food() < food {
			return false
		}
	}

	if unsafe.Wood() < wood {
		wood -= unsafe.Wood()
		if safe.Wood() < wood {
			return false
		}
	}

	if unsafe.Stone() < stone {
		stone -= unsafe.Stone()
		if safe.Stone() < stone {
			return false
		}
	}

	return true
}

func (d *hero_resource) GetResourceMultiple(gold, food, wood, stone uint64) uint64 {

	unsafe := d.GetUnsafeResource()
	safe := d.GetSafeResource()

	var n uint64 = 1
	if gold > 0 {
		n = (unsafe.gold + safe.gold) / gold
		if n <= 0 {
			return 0
		}
	}

	if food > 0 {
		n = u64.Min(n, (unsafe.food+safe.food)/food)
		if n <= 0 {
			return 0
		}
	}

	if wood > 0 {
		n = u64.Min(n, (unsafe.wood+safe.wood)/wood)
		if n <= 0 {
			return 0
		}
	}

	if stone > 0 {
		n = u64.Min(n, (unsafe.stone+safe.stone)/stone)
		if n <= 0 {
			return 0
		}
	}

	return n
}

func (hero *hero_resource) SetNextResDecayTime(t time.Time) {
	hero.nextResDecayTime = t
}

func (hero *hero_resource) NextResDecayTime() time.Time {
	return hero.nextResDecayTime
}

// 玉石矿

func (hero *hero_resource) JadeOre() uint64 {
	return hero.jadeOre
}

func (hero *hero_resource) AddJadeOre(toAdd uint64) uint64 {
	hero.jadeOre += toAdd
	return hero.jadeOre
}

func (hero *hero_resource) ReduceJadeOre(toReduce uint64) uint64 {
	hero.jadeOre = u64.Sub(hero.jadeOre, toReduce)
	return hero.jadeOre
}

func (hero *hero_resource) Jade() uint64 {
	return hero.jade
}

func (hero *hero_resource) AddJade(toAdd uint64) uint64 {
	hero.jade += toAdd
	hero.historyJade += toAdd
	hero.todayObtainJade += toAdd
	return hero.jade
}

func (hero *hero_resource) ReduceJade(toReduce uint64) uint64 {
	hero.jade = u64.Sub(hero.jade, toReduce)
	return hero.jade
}

func (hero *hero_resource) HistoryJade() uint64 {
	return hero.historyJade
}

func (hero *hero_resource) TodayObtainJade() uint64 {
	return hero.todayObtainJade
}
