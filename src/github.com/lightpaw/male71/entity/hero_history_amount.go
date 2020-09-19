package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
)

func NewHeroHistoryAmount(maps *hero_maps) *HeroHistoryAmount {
	h := &HeroHistoryAmount{}
	h.amountMap = maps.getOrCreateMap(server_proto.HeroMapCategory_history_amount, false)
	h.dailyMap = maps.getOrCreateMap(server_proto.HeroMapCategory_daily_amount, true)
	return h
}

type HeroHistoryAmount struct {
	amountMap *heromap
	dailyMap  *heromap
}

func (hero *HeroHistoryAmount) Amount(key server_proto.HistoryAmountType) uint64 {
	return hero.AmountWithSubType(key, 0)
}

func (hero *HeroHistoryAmount) DailyAmount(key server_proto.HistoryAmountType) uint64 {
	return hero.DailyAmountWithSubType(key, 0)
}

func (hero *HeroHistoryAmount) GetAmount(daily bool, key server_proto.HistoryAmountType) uint64 {
	if daily {
		return hero.DailyAmount(key)
	} else {
		return hero.Amount(key)
	}
}

func (hero *HeroHistoryAmount) IncreaseOne(key server_proto.HistoryAmountType) {
	hero.IncreaseOneWithSubType(key, 0)
}

func (hero *HeroHistoryAmount) Increase(key server_proto.HistoryAmountType, toAdd uint64) {
	if toAdd > 0 {
		hero.IncreaseWithSubType(key, 0, toAdd)
	}
}

func (hero *HeroHistoryAmount) SetIfGreater(key server_proto.HistoryAmountType, toSet uint64) bool {
	changed := hero.SetIfGreaterWithSubType(false, key, 0, toSet)
	return hero.SetIfGreaterWithSubType(true, key, 0, toSet) || changed
}

func (hero *HeroHistoryAmount) Set(key server_proto.HistoryAmountType, toSet uint64) {
	hero.SetWithSubType(false, key, 0, toSet)
	hero.SetWithSubType(true, key, 0, toSet)
}

func (hero *HeroHistoryAmount) SetAmount(daily bool, key server_proto.HistoryAmountType, toSet uint64) {
	hero.SetWithSubType(daily, key, 0, toSet)
}

func (hero *HeroHistoryAmount) AmountWithSubType(key server_proto.HistoryAmountType, subType uint64) uint64 {
	return hero.amountMap.internalMap[combineKey(key, subType)]
}

func (hero *HeroHistoryAmount) DailyAmountWithSubType(key server_proto.HistoryAmountType, subType uint64) uint64 {
	return hero.dailyMap.internalMap[combineKey(key, subType)]
}

func (hero *HeroHistoryAmount) GetAmountWithSubType(daily bool, key server_proto.HistoryAmountType, subType uint64) uint64 {
	if daily {
		return hero.DailyAmountWithSubType(key, subType)
	} else {
		return hero.AmountWithSubType(key, subType)
	}
}

func (hero *HeroHistoryAmount) IncreaseOneWithSubType(key server_proto.HistoryAmountType, subType uint64) {
	hero.IncreaseWithSubType(key, subType, 1)
}

func (hero *HeroHistoryAmount) IncreaseWithSubType(key server_proto.HistoryAmountType, subType, toAdd uint64) {
	if toAdd > 0 {
		hero.amountMap.internalMap[combineKey(key, subType)] += toAdd
		hero.dailyMap.internalMap[combineKey(key, subType)] += toAdd
	}
}

func (hero *HeroHistoryAmount) SetWithSubType(daily bool, key server_proto.HistoryAmountType, subType, toSet uint64) {
	if daily {
		hero.dailyMap.internalMap[combineKey(key, subType)] = toSet
	} else {
		hero.amountMap.internalMap[combineKey(key, subType)] = toSet
	}
}

func (hero *HeroHistoryAmount) SetIfGreaterWithSubType(daily bool, key server_proto.HistoryAmountType, subType, toSet uint64) (changed bool) {
	t := combineKey(key, subType)
	if daily {
		if hero.dailyMap.internalMap[t] < toSet {
			hero.dailyMap.internalMap[t] = toSet
			changed = true
		}
	} else {
		if hero.amountMap.internalMap[t] < toSet {
			hero.amountMap.internalMap[t] = toSet
			changed = true
		}
	}
	return
}

func combineKey(key server_proto.HistoryAmountType, subType uint64) uint64 {
	return subType<<10 | uint64(key)
}
