package entity

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/fishing_data"
	"github.com/lightpaw/male7/util/i64"
)

func NewHeroFishing() *HeroFishing {
	return &HeroFishing{
		fishingTimesMap:  make(map[uint64]uint64),
		fishingCountdown: make(map[uint64]int64),
	}
}

type HeroFishing struct {
	fishingTimesMap  map[uint64]uint64 // 钓鱼1次跟钓鱼10次每种钓的次数，key是 FishingCostProto.Times, vlaue是 该钓鱼次数类型的钓鱼次数
	fishingCountdown map[uint64]int64

	captainSet uint64 // 金杆钓的武将设置
}

func (h *HeroFishing) CaptainSet() uint64 {
	return h.captainSet
}

func (h *HeroFishing) ChangeCaptainSet(set uint64) {
	h.captainSet = set
}

func (h *HeroFishing) FishingTimes(fishId uint64) uint64 {
	return h.fishingTimesMap[fishId]
}

func (h *HeroFishing) FishingCountdown(fishId uint64) int64 {
	return h.fishingCountdown[fishId]
}

func (h *HeroFishing) SetFishingCountdown(fishId uint64, toSet int64) {
	h.fishingCountdown[fishId] = toSet
}

func (h *HeroFishing) IncreFishingTimes(fishId uint64) uint64 {
	newTimes := h.fishingTimesMap[fishId] + 1
	h.fishingTimesMap[fishId] = newTimes
	return newTimes
}

func (h *HeroFishing) ResetDaily() {
	if len(h.fishingTimesMap) > 0 {
		for key := range h.fishingTimesMap {
			delete(h.fishingTimesMap, key)
		}
	}
}

func (h *HeroFishing) encode() *shared_proto.HeroFishingProto {

	proto := &shared_proto.HeroFishingProto{}
	if len(h.fishingTimesMap) > 0 {
		proto.Times = make([]int32, 0, len(h.fishingTimesMap))
		proto.FishingTimes = make([]int32, 0, len(h.fishingTimesMap))
		proto.FishType = make([]int32, 0, len(h.fishingTimesMap))
		proto.Countdown = make([]int32, 0, len(h.fishingTimesMap))

		for fishId, fishingTimes := range h.fishingTimesMap {
			fishType, times := fishing_data.FishTypeTimes(fishId)
			proto.Times = append(proto.Times, u64.Int32(times))
			proto.FishingTimes = append(proto.FishingTimes, u64.Int32(fishingTimes))
			proto.FishType = append(proto.FishType, u64.Int32(fishType))
			proto.Countdown = append(proto.Countdown, i64.Int32(h.FishingCountdown(fishId)))
		}
	}

	proto.CaptainSet = u64.Int32(h.captainSet)

	return proto
}

func (h *HeroFishing) unmarshal(proto *shared_proto.HeroFishingProto) {
	if proto == nil {
		return
	}

	n := imath.Minx(len(proto.Times), len(proto.FishingTimes), len(proto.FishType), len(proto.Countdown))
	for i := 0; i < n; i++ {
		fishId := fishing_data.FishId(u64.FromInt32(proto.FishType[i]), u64.FromInt32(proto.Times[i]))
		h.fishingTimesMap[fishId] = u64.FromInt32(proto.FishingTimes[i])

		if c := proto.Countdown[i]; c > 0 {
			h.fishingCountdown[fishId] = int64(c)
		}
	}
	h.captainSet = u64.FromInt32(proto.CaptainSet)
}
