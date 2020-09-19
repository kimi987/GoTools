package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func NewHeroVip() *HeroVip {
	m := &HeroVip{}
	m.ResetContinueDays()
	m.dailyPrizeCanCollectLevel = make(map[uint64]bool)
	m.AddDailyPrizeCanCollectLevel(m.level)
	m.levelPrizeCollectedLevel = make(map[uint64]bool)
	m.vipShopGoodsBoughtCount = make(map[uint64]uint64)
	return m
}

type HeroVip struct {
	level uint64
	exp   uint64

	dailyFirstLoginTime time.Time
	continueDays        uint64

	dailyPrizeCanCollectLevel map[uint64]bool
	levelPrizeCollectedLevel  map[uint64]bool

	dailyPrizeCollectedLevel []uint64

	vipShopGoodsBoughtCount map[uint64]uint64
}

func (h *HeroVip) ResetDaily(ctime time.Time) {
	h.dailyPrizeCanCollectLevel = make(map[uint64]bool)
	h.AddDailyPrizeCanCollectLevel(h.level)
	h.vipShopGoodsBoughtCount = make(map[uint64]uint64)
	h.dailyPrizeCollectedLevel = []uint64{}
}

func (h *HeroVip) Level() uint64 {
	return h.level
}

func (h *HeroVip) Exp() uint64 {
	return h.exp
}

func (h *HeroVip) SetLevel(toSet uint64) {
	h.level = toSet
}

func (h *HeroVip) AddExp(toAdd uint64) uint64 {
	h.exp += toAdd
	return h.exp
}

func (h *HeroVip) SetExp(toSet uint64) {
	h.exp = toSet
}

func (h *HeroVip) DailyFirstLoginTime() time.Time {
	return h.dailyFirstLoginTime
}

func (h *HeroVip) SetDailyFirstLoginTime(toSet time.Time) {
	h.dailyFirstLoginTime = toSet
}

func (h *HeroVip) ContinueDays() uint64 {
	return h.continueDays
}

func (h *HeroVip) IncrContinueDays() uint64 {
	h.continueDays++
	return h.continueDays
}

func (h *HeroVip) ResetContinueDays() uint64 {
	h.continueDays = 1
	return h.continueDays
}

func (h *HeroVip) CanCollectDailyPrize(level uint64) bool {
	_, ok := h.dailyPrizeCanCollectLevel[level]
	return ok
}

func (h *HeroVip) AddDailyPrizeCanCollectLevel(level uint64) {
	h.dailyPrizeCanCollectLevel[level] = true
}

func (h *HeroVip) CollectDailyPrize(level uint64) {
	delete(h.dailyPrizeCanCollectLevel, level)
	h.dailyPrizeCollectedLevel = append(h.dailyPrizeCollectedLevel, level)
}

func (h *HeroVip) CanCollectLevelPrize(level uint64) bool {
	_, ok := h.levelPrizeCollectedLevel[level]
	return !ok
}

func (h *HeroVip) CollectLevelPrize(level uint64) {
	h.levelPrizeCollectedLevel[level] = true
}

func (h *HeroVip) AddVipShopGoodsBoughtCount(shopGoodsId, amount uint64) {
	h.vipShopGoodsBoughtCount[shopGoodsId] = h.vipShopGoodsBoughtCount[shopGoodsId] + amount
}

func (h *HeroVip) VipShopGoodsBoughtCount(shopGoodsId uint64) uint64 {
	return h.vipShopGoodsBoughtCount[shopGoodsId]
}

func (h *HeroVip) Encode() *shared_proto.HeroVipProto {
	p := &shared_proto.HeroVipProto{}
	p.Level = u64.Int32(h.level)
	p.Exp = u64.Int32(h.level)
	p.ContinueDays = u64.Int32(h.continueDays)

	for k := range h.dailyPrizeCanCollectLevel {
		p.DailyPrizeCanCollectLevel = append(p.DailyPrizeCanCollectLevel, u64.Int32(k))
	}
	for k := range h.levelPrizeCollectedLevel {
		p.LevelPrizeCollectedLevel = append(p.LevelPrizeCollectedLevel, u64.Int32(k))
	}
	for k := range h.vipShopGoodsBoughtCount {
		p.VipShopBoughtId = append(p.VipShopBoughtId, u64.Int32(k))
	}
	p.DailyPrizeCollectedLevel = u64.Int32Array(h.dailyPrizeCollectedLevel)

	return p
}

func (h *HeroVip) encodeServer() *server_proto.HeroVipServerProto {
	p := &server_proto.HeroVipServerProto{}
	p.Level = h.level
	p.Exp = h.exp
	p.DailyFirstLoginTime = timeutil.Marshal64(h.dailyFirstLoginTime)
	p.ContinueDays = h.continueDays

	p.DailyPrizeCanCollectLevel = h.dailyPrizeCanCollectLevel
	p.LevelPrizeCollectedLevel = h.levelPrizeCollectedLevel
	p.VipShopBoughtCount = h.vipShopGoodsBoughtCount
	p.DailyPrizeCollectedLevel = h.dailyPrizeCollectedLevel

	return p
}

func (h *HeroVip) unmarshal(p *server_proto.HeroVipServerProto, datas interface {
	VipLevelData() *config.VipLevelDataConfig
}) {
	if p == nil {
		return
	}

	d := datas.VipLevelData().Must(p.Level)
	if d.Level != p.Level {
		logrus.Errorf("unmarshal heroVip, 找不到 vipLevel：%v, 改为：%v", p.Level, d.Level)
	}
	h.level = d.Level

	h.exp = p.Exp
	h.dailyFirstLoginTime = timeutil.Unix64(p.DailyFirstLoginTime)
	h.continueDays = p.ContinueDays

	if len(p.DailyPrizeCanCollectLevel) > 0 {
		h.dailyPrizeCanCollectLevel = p.DailyPrizeCanCollectLevel
	} else {
		h.dailyPrizeCanCollectLevel = make(map[uint64]bool)
	}
	if p.LevelPrizeCollectedLevel != nil {
		h.levelPrizeCollectedLevel = p.LevelPrizeCollectedLevel
	}
	if p.VipShopBoughtCount != nil {
		h.vipShopGoodsBoughtCount = p.VipShopBoughtCount
	}
	h.dailyPrizeCollectedLevel = p.DailyPrizeCollectedLevel
}

func (h *HeroVip) UpdateContinueDays(ctime time.Time) {
	// 同一天
	if timeutil.IsSameDay(h.dailyFirstLoginTime, ctime) {
		return
	}

	defer h.SetDailyFirstLoginTime(ctime)

	// 第一次登录
	if timeutil.IsZero(h.dailyFirstLoginTime) {
		h.ResetContinueDays()
		return
	}

	// 第二天
	if timeutil.Midnight(h.dailyFirstLoginTime.Add(time.Duration(24 * time.Hour))).Equal(timeutil.Midnight(ctime)) {
		h.IncrContinueDays()
	} else {
		h.ResetContinueDays()
	}
}
