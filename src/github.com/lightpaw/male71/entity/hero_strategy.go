package entity

import (
	"time"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/idbytes"
)

func newStratagem(id uint64, ctime time.Time) *stratagem {
	return &stratagem {
		id: id,
		nextUseableTime: ctime,
	}
}

// internal or military 计策
type stratagem struct {
	id              uint64    // 策略ID
	dailyUsedTimes  uint64    // 今日使用该策略的次数
	nextUseableTime time.Time // 下次可发动的时间
}

func (s *stratagem) Id() uint64 {
	return s.id
}

func (s *stratagem) DailyUsedTimes() uint64 {
	return s.dailyUsedTimes
}

func (s *stratagem) NextUseableTime() time.Time {
	return s.nextUseableTime
}

func (s *stratagem) EncodeClient() *shared_proto.StratagemProto {
	return &shared_proto.StratagemProto {
		Id: u64.Int32(s.id),
		DailyUsedTimes: u64.Int32(s.dailyUsedTimes),
		NextUseableTime: timeutil.Marshal32(s.nextUseableTime),
	}
}

// 中计
type trappedStratagem struct {
	id                 uint64  // 策略ID
	endTime            time.Time // 持续结束时间（可用特殊道具缩短时间或解除）
}

func (s *trappedStratagem) EncodeClient() *shared_proto.TrappedStratagemProto {
	return &shared_proto.TrappedStratagemProto {
		Id: u64.Int32(s.id),
		EndTime: timeutil.Marshal32(s.endTime),
	}
}

func newHeroStrategy() *HeroStrategy {
	return &HeroStrategy {
		stratagems: make(map[uint64]*stratagem),
		trappedStratagems: make(map[uint64]*trappedStratagem),
		todayTargetTimes: make(map[int64]uint64),
	}
}

type HeroStrategy struct {
	stratagems          map[uint64]*stratagem // 已经激活的策略列表
	trappedStratagems   map[uint64]*trappedStratagem // 中计列表

	todayTargetTimes    map[int64]uint64 // 今日对别的玩家使用军事计策的次数列表
	todayTrappedTimes   uint64  // 今日中计的次数
}

func (s *HeroStrategy) ResetDaily() {
	s.todayTrappedTimes = 0
	s.todayTargetTimes = make(map[int64]uint64)
	for _, v := range s.stratagems {
		v.dailyUsedTimes = 0
	}
}

func (s *HeroStrategy) IsStratagemCd(id uint64, ctime time.Time) bool {
	stratagem := s.stratagems[id]
	if stratagem == nil {
		return false
	}
	return stratagem.nextUseableTime.After(ctime)
}

func (s *HeroStrategy) GetTodayUsedTimes(id uint64) uint64 {
	stratagem := s.stratagems[id]
	if stratagem == nil {
		return 0
	}
	return stratagem.dailyUsedTimes
}

func (s *HeroStrategy) GetTodayTargetTimes(targetId int64) uint64 {
	return s.todayTargetTimes[targetId]
}

func (s *HeroStrategy) TodayTrappedTimes() uint64 {
	return s.todayTrappedTimes
}

func (s *HeroStrategy) IsTrappedStratagemEnd(id uint64, ctime time.Time) bool {
	trappedStratagem := s.trappedStratagems[id]
	if trappedStratagem == nil {
		return true
	}
	return trappedStratagem.endTime.Before(ctime)
}

// 中计
func (s *HeroStrategy) TrappedStratagem(id uint64, endTime time.Time) *trappedStratagem {
	s.todayTrappedTimes++
	stratagem := s.trappedStratagems[id]
	if stratagem != nil {
		stratagem.endTime = endTime
	} else {
		stratagem = &trappedStratagem {
			id: id,
			endTime: endTime,
		}
		s.trappedStratagems[id] = stratagem
	}
	return stratagem
}

// 施计
func (s *HeroStrategy) UseStratagem(id uint64, target int64, nextUseableTime time.Time) *stratagem {
	if target != 0 {
		s.todayTargetTimes[target]++
	}
	str := s.stratagems[id]
	if str != nil {
		str.dailyUsedTimes++
		str.nextUseableTime = nextUseableTime
	} else {
		str = &stratagem {
			id:              id,
			dailyUsedTimes:  1,
			nextUseableTime: nextUseableTime,
		}
		s.stratagems[id] = str
	}
	return str
}

func (s *HeroStrategy) ReturnStratgem(id uint64, target int64, nextUseableTime time.Time) *stratagem {
	if target != 0 {
		s.todayTargetTimes[target] = u64.Sub(s.todayTargetTimes[target], 1)
	}
	str := s.stratagems[id]
	if str != nil {
		str.dailyUsedTimes = u64.Sub(str.dailyUsedTimes, 1)
		str.nextUseableTime = nextUseableTime
	}
	return str

}

func (s *HeroStrategy) unmarshal(proto *server_proto.HeroStrategyServerProto) {
	if proto == nil {
		return
	}

	if len(proto.Stratagems) > 0 {
		for _, v := range proto.Stratagems {
			s.stratagems[v.Id] = &stratagem {
				id:              v.Id,
				dailyUsedTimes:  v.DailyUsedTimes,
				nextUseableTime: timeutil.Unix64(v.NextUseableTime),
			}
		}
	}

	if len(proto.TrappedStratagems) > 0 {
		for _, v := range proto.TrappedStratagems {
			s.trappedStratagems[v.Id] = &trappedStratagem {
				id: v.Id,
				endTime: timeutil.Unix64(v.EndTime),
			}
		}
	}

	if len(proto.TodayTargetTimes) > 0 {
		for _, v := range proto.TodayTargetTimes {
			s.todayTargetTimes[v.Id] = v.Times
		}
	}

	s.todayTrappedTimes = proto.TodayTrappedTimes
}

func (s *HeroStrategy) encode() *server_proto.HeroStrategyServerProto {
	proto := &server_proto.HeroStrategyServerProto{}

	for k, v := range s.stratagems {
		proto.Stratagems = append(proto.Stratagems, &server_proto.StratagemServerProto {
			Id: k,
			DailyUsedTimes: v.dailyUsedTimes,
			NextUseableTime: timeutil.Marshal64(v.nextUseableTime),
		})
	}

	for k, v := range s.trappedStratagems {
		proto.TrappedStratagems = append(proto.TrappedStratagems, &server_proto.TrappedStratagemServerProto {
			Id: k,
			EndTime: timeutil.Marshal64(v.endTime),
		})
	}

	for k, v := range s.todayTargetTimes {
		proto.TodayTargetTimes = append(proto.TodayTargetTimes, &server_proto.StratagemTargetTimesServerProto {
			Id: k,
			Times: v,
		})
	}

	proto.TodayTrappedTimes = s.todayTrappedTimes

	return proto
}

func (s *HeroStrategy) encodeClient() *shared_proto.HeroStrategyProto {
	proto := &shared_proto.HeroStrategyProto{}

	for k, v := range s.stratagems {
		proto.Stratagems = append(proto.Stratagems, &shared_proto.StratagemProto {
			Id: u64.Int32(k),
			DailyUsedTimes: u64.Int32(v.dailyUsedTimes),
			NextUseableTime: timeutil.Marshal32(v.nextUseableTime),
		})
	}

	for k, v := range s.trappedStratagems {
		proto.TrappedStratagems = append(proto.TrappedStratagems, &shared_proto.TrappedStratagemProto {
			Id: u64.Int32(k),
			EndTime: timeutil.Marshal32(v.endTime),
		})
	}

	for k, v := range s.todayTargetTimes {
		proto.TodayTargetTimes = append(proto.TodayTargetTimes, &shared_proto.StratagemTargetTimesProto {
			Id: idbytes.ToBytes(k),
			Times: u64.Int32(v),
		})
	}

	return proto
}
