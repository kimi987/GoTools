package entity

import (
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/pb/server_proto"
	"sort"
	"time"
)

func newHeroEventMap(maps *hero_maps) *guild_event_prize_map {

	return &guild_event_prize_map{
		historyGuildEventPrizeTriggerTimes:  maps.getOrCreateMap(server_proto.HeroMapCategory_history_guild_event_prize_trigger_times, false),
		dailyGuildEventPrizeTriggerTimes:    maps.getOrCreateMap(server_proto.HeroMapCategory_daily_guild_event_prize_trigger_times, true),
		dailyGuildEventPrizeCollectTimes:    maps.getOrCreateMap(server_proto.HeroMapCategory_daily_guild_event_prize_collect_times, true),
		dailyGuildPrestigeEventTriggerTimes: maps.getOrCreateMap(server_proto.HeroMapCategory_daily_guild_prestige_event_trigger_times, true),
		guildEventPrizes:                    make(map[int32]*HeroGuildEventPrize),
	}
}

type guild_event_prize_map struct {
	historyGuildEventPrizeTriggerTimes  *heromap
	dailyGuildEventPrizeTriggerTimes    *heromap
	dailyGuildEventPrizeCollectTimes    *heromap
	dailyGuildPrestigeEventTriggerTimes *heromap

	guildEventPrizeIdGen int32
	guildEventPrizes     map[int32]*HeroGuildEventPrize
}

func (m *guild_event_prize_map) resetDaily(resetTime time.Time) {
	for k, p := range m.guildEventPrizes {
		if p.ExpireTime.Before(resetTime) {
			delete(m.guildEventPrizes, k)
		}
	}
}

func (hero *Hero) AddGuildEventPrize(data *guild_data.GuildEventPrizeData, heroId int64, expireTime time.Time, hideGiver bool) int32 {
	hero.eventMap.guildEventPrizeIdGen++
	id := hero.eventMap.guildEventPrizeIdGen

	hero.eventMap.guildEventPrizes[id] = &HeroGuildEventPrize{
		Id:         id,
		Data:       data,
		SendHeroId: heroId,
		ExpireTime: expireTime,
		HideGiver:  hideGiver,
	}
	return id
}

func (hero *Hero) GetGuildEventPrize(id int32) *HeroGuildEventPrize {
	return hero.eventMap.guildEventPrizes[id]
}

func (hero *Hero) RemoveGuildEventPrize(id int32) {
	delete(hero.eventMap.guildEventPrizes, id)
}

func (hero *Hero) GetGuildEventPrizeCount() int {
	return len(hero.eventMap.guildEventPrizes)
}

func (hero *Hero) WalkGuildEventPrize(f func(p *HeroGuildEventPrize)) {
	for _, p := range hero.eventMap.guildEventPrizes {
		f(p)
	}
}

func (hero *Hero) SortGuildEventPrizeExpireSlice() (slice []*HeroGuildEventPrize) {
	if len(hero.eventMap.guildEventPrizes) > 0 {
		for _, p := range hero.eventMap.guildEventPrizes {
			slice = append(slice, p)
		}
		sort.Sort(heroGuildEventPrizeExpireSlice(slice))
	}
	return
}

type heroGuildEventPrizeExpireSlice []*HeroGuildEventPrize

func (p heroGuildEventPrizeExpireSlice) Len() int           { return len(p) }
func (p heroGuildEventPrizeExpireSlice) Less(i, j int) bool { return p[i].ExpireTime.Before(p[j].ExpireTime) }
func (p heroGuildEventPrizeExpireSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// 可以领取的联盟礼包
type HeroGuildEventPrize struct {
	Id int32

	Data *guild_data.GuildEventPrizeData

	// 谁发的
	SendHeroId int64

	// 过期时间
	ExpireTime time.Time

	// 是否神秘发送
	HideGiver bool
}

func (hero *Hero) GetGuildEventPrizeTriggerTimesData(isDaily bool) *heromap {
	if isDaily {
		return hero.eventMap.dailyGuildEventPrizeTriggerTimes
	} else {
		return hero.eventMap.historyGuildEventPrizeTriggerTimes
	}
}

func (hero *Hero) GetGuildEventPrizeCollectTimesData() *heromap {
	return hero.eventMap.dailyGuildEventPrizeCollectTimes
}

func (hero *Hero) GetDailyGuildPrestigeEventTriggerTimesData() *heromap {
	return hero.eventMap.dailyGuildPrestigeEventTriggerTimes
}

//type event_map struct {
//	category HeroEventCategory
//
//	eventMap *heromap
//}
//
//func (m *event_map) Amount(eventType shared_proto.HeroEvent) uint64 {
//	return m.AmountWithSubType(eventType, 0)
//}
//
//func (m *event_map) IncreaseOne(eventType shared_proto.HeroEvent) {
//	m.IncreaseOneWithSubType(eventType, 0)
//}
//
//func (m *event_map) Increase(eventType shared_proto.HeroEvent, toAdd uint64) {
//	m.IncreaseWithSubType(eventType, 0, toAdd)
//}
//
//func (m *event_map) Set(eventType shared_proto.HeroEvent, toSet uint64) {
//	m.SetWithSubType(eventType, 0, toSet)
//}
//
//func (m *event_map) AmountWithSubType(eventType shared_proto.HeroEvent, subType uint64) uint64 {
//	return m.eventMap.internalMap[combineHeroEventKey(m.category, eventType, subType)]
//}
//
//func (m *event_map) IncreaseOneWithSubType(eventType shared_proto.HeroEvent, subType uint64) {
//	m.IncreaseWithSubType(eventType, subType, 1)
//}
//
//func (m *event_map) IncreaseWithSubType(eventType shared_proto.HeroEvent, subType, toAdd uint64) {
//	m.eventMap.internalMap[combineHeroEventKey(m.category, eventType, subType)] += toAdd
//}
//
//func (m *event_map) SetWithSubType(eventType shared_proto.HeroEvent, subType, toSet uint64) {
//	m.eventMap.internalMap[combineHeroEventKey(m.category, eventType, subType)] = toSet
//}
//
//type HeroEventCategory uint8
//
//const (
//	GuildEventPrize HeroEventCategory = 0
//)
//
//func combineHeroEventKey(category HeroEventCategory, eventType shared_proto.HeroEvent, key uint64) uint64 {
//	return uint64(category) | (uint64(eventType) << 4) | (key << 16)
//}
