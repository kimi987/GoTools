package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
)

func newHeroActivity() *hero_activity {
	p := &hero_activity{}
	p.collections = make(map[uint64]*collection_info)
	return p
}

type collection_info struct {
	exchangeTimes map[uint64]uint64 // <兑换id, 兑换次数>

	startTime time.Time
	endTime   time.Time
}

func (info *collection_info) encodeServer() *server_proto.HeroCollectionServerProto {
	return &server_proto.HeroCollectionServerProto {
		ExchangeTimes: info.exchangeTimes,
		StartTime: timeutil.Marshal64(info.startTime),
		EndTime: timeutil.Marshal64(info.endTime),
	}
}

func (info *collection_info) encodeClient() []*shared_proto.Int32Pair {
	p := []*shared_proto.Int32Pair{}
	for id, times := range info.exchangeTimes {
		p = append(p, &shared_proto.Int32Pair {
			Key: u64.Int32(id),
			Value: u64.Int32(times),
		})
	}
	return p
}

type hero_activity struct {
	collections map[uint64]*collection_info // <活动id, 兑换信息>
}

func (hero *Hero) Activity() *hero_activity {
	return hero.activity
}

func (h *hero_activity) GetCollectionTimes(id, exchangeId uint64) uint64 {
	if info := h.collections[id]; info != nil {
		return info.exchangeTimes[exchangeId]
	}
	return 0
}

func (h *hero_activity) IncCollectionTimes(id, exchangeId uint64) {
	if info := h.collections[id]; info != nil {
		info.exchangeTimes[exchangeId]++
	}
}

func (h *hero_activity) TryCorrectCollection(collections map[uint64]*ActivityCollection) (corrected bool) {
	for id, _ := range h.collections {
		if collections[id] == nil {
			delete(h.collections, id)
			corrected = true
		}
	}
	for id, c := range collections {
		info := h.collections[id]
		if info != nil {
			if !info.startTime.Equal(c.startTime) || !info.endTime.Equal(c.endTime) {
				if len(info.exchangeTimes) > 0 {
					info.exchangeTimes = make(map[uint64]uint64)
				}
				info.startTime = c.startTime
				info.endTime = c.endTime
				corrected = true
			}
		} else {
			h.collections[id] = &collection_info {
				exchangeTimes: make(map[uint64]uint64),
				startTime: c.startTime,
				endTime: c.endTime,
			}
			corrected = true
		}
	}
	return
}

func (h *hero_activity) EncodeClient() *shared_proto.HeroAllCollectionProto {
	p := &shared_proto.HeroAllCollectionProto{}
	for id, c := range h.collections {
		p.Collections = append(p.Collections, &shared_proto.HeroCollectionProto {
			Id: u64.Int32(id),
			ExchangeTimes: c.encodeClient(),
		})
	}
	return p
}

func (h *hero_activity) unmarshal(proto *server_proto.HeroActivityServerProto, ctime time.Time) {
	if proto == nil {
		return
	}
	for id, c := range proto.Collections {
		start := timeutil.Unix64(c.StartTime)
		end := timeutil.Unix64(c.EndTime)
		if !timeutil.Between(ctime, start, end) {
			continue
		}
		info := &collection_info {
			exchangeTimes: c.ExchangeTimes,
			startTime: start,
			endTime: end,
		}
		if len(c.ExchangeTimes) > 0 {
			info.exchangeTimes = c.ExchangeTimes
		} else {
			info.exchangeTimes = make(map[uint64]uint64)
		}
		h.collections[id] = info
	}
}

func (h *hero_activity) encodeServer() *server_proto.HeroActivityServerProto {
	proto := &server_proto.HeroActivityServerProto {}
	if len(h.collections) > 0 {
		proto.Collections = make(map[uint64]*server_proto.HeroCollectionServerProto)
		for id, c := range h.collections {
			proto.Collections[id] = c.encodeServer()
		}
	}
	return proto
}
