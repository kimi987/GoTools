package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/promdata"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
)

func newTimeLimitGiftInfo(groupId, index uint64, endTime time.Time) *timeLimitGiftInfo {
	return &timeLimitGiftInfo {
		groupId: groupId,
		index: index,
		endTime: endTime,
	}
}

type timeLimitGiftInfo struct {
	groupId       uint64
	index         uint64
	endTime       time.Time
}

func newHeroPromotion(maps *hero_maps) *hero_promotion {
	promotion := &hero_promotion{}
	promotion.loginDayPrizeCollected = maps.getOrCreateKeys(server_proto.HeroMapCategory_login_day_prize, false)
	promotion.heroLevelFundCollected = maps.getOrCreateKeys(server_proto.HeroMapCategory_hero_level_fund, false)
	promotion.dailySpCollected = maps.getOrCreateKeys(server_proto.HeroMapCategory_daily_sp_collection, true)
	promotion.dailyFreeGiftCollected = maps.getOrCreateKeys(server_proto.HeroMapCategory_daily_free_gift_collection, true)
	promotion.freeGiftCollected = maps.getOrCreateKeys(server_proto.HeroMapCategory_free_gift_collection, false)
	promotion.dailyEventLimitGift = maps.getOrCreateKeys(server_proto.HeroMapCategory_daily_event_limit_gift, true)
	promotion.timeLimitGifts = make(map[uint64]*timeLimitGiftInfo)
	promotion.eventLimitGifts = make(map[uint64]*eventLimitGift)
	promotion.timeLimitGiftBuyTimes = maps.getOrCreateMap(server_proto.HeroMapCategory_time_limit_gift_buy_times, false)
	promotion.eventLimitGiftBuyTimes = maps.getOrCreateMap(server_proto.HeroMapCategory_event_limit_gift_buy_times, false)

	return promotion
}

type hero_promotion struct {
	loginDayPrizeCollected *herokeys

	heroLevelFundCollected *herokeys

	dailySpCollected       *herokeys // 每日领取过的阶段体力值
	dailyFreeGiftCollected *herokeys // 日常免费礼包
	freeGiftCollected      *herokeys // 一次性免费礼包
	dailyEventLimitGift    *herokeys // 日常事件时限礼包

	timeLimitGifts         map[uint64]*timeLimitGiftInfo // 时限礼包的当前购买信息
	eventLimitGifts        map[uint64]*eventLimitGift // <礼包id, 整个eventLimitGift对象>

	timeLimitGiftBuyTimes  *heromap // 每种时限礼包的购买次数
	eventLimitGiftBuyTimes *heromap // 每种事件礼包的购买次数
}

func (hero *Hero) Promotion() *hero_promotion {
	return hero.promotion
}

func (p *hero_promotion) IsLoginDayCollected(day uint64) bool {
	return p.loginDayPrizeCollected.Exist(day)
}

func (p *hero_promotion) CollectLoginDay(day uint64) {
	p.loginDayPrizeCollected.Add(day)
}

func (p *hero_promotion) IsHeroLevelFundCollected(level uint64) bool {
	return p.heroLevelFundCollected.Exist(level)
}

func (p *hero_promotion) CollectHeroLevelFund(level uint64) {
	p.heroLevelFundCollected.Add(level)
}

func (p *hero_promotion) IsDailySpCollected(id uint64) bool {
	return p.dailySpCollected.Exist(id)
}

func (p *hero_promotion) CollectDailySp(id uint64) {
	p.dailySpCollected.Add(id)
}

func (p *hero_promotion) IsDailyEventLimitGiftAppeared(id uint64) bool {
	return p.dailyEventLimitGift.Exist(id)
}

func (p *hero_promotion) SetDailyEventLimitGiftAppeared(id uint64) {
	p.dailyEventLimitGift.Add(id)
}

func (p *hero_promotion) IsFreeGiftCollected(data *promdata.FreeGiftData) bool {
	if data.Daily {
		return p.dailyFreeGiftCollected.Exist(data.Id)
	}
	return p.freeGiftCollected.Exist(data.Id)
}

func (p *hero_promotion) CollectFreeGift(data *promdata.FreeGiftData) {
	if data.Daily {
		p.dailyFreeGiftCollected.Add(data.Id)
	}
	p.freeGiftCollected.Add(data.Id)
}

func (p *hero_promotion) encodeClient() *shared_proto.HeroPromotionProto {

	proto := &shared_proto.HeroPromotionProto{}

	for day := range p.loginDayPrizeCollected.internalMap {
		proto.CollectLoginPrizeDay = append(proto.CollectLoginPrizeDay, u64.Int32(day))
	}

	for level := range p.heroLevelFundCollected.internalMap {
		proto.CollectLevelFund = append(proto.CollectLevelFund, u64.Int32(level))
	}

	for v := range p.dailySpCollected.internalMap {
		proto.CollectDailySp = append(proto.CollectDailySp, u64.Int32(v))
	}

	for v := range p.dailyFreeGiftCollected.internalMap {
		proto.CollectDailyFreeGift = append(proto.CollectDailyFreeGift, u64.Int32(v))
	}

	for v := range p.freeGiftCollected.internalMap {
		proto.CollectFreeGift = append(proto.CollectFreeGift, u64.Int32(v))
	}

	for id, info := range p.timeLimitGifts {
		proto.BuyTimeLimitGift = append(proto.BuyTimeLimitGift, &shared_proto.TimeLimitGiftBoughtProto {
			Id: u64.Int32(id),
			Index: u64.Int32(info.index),
			EndTime: timeutil.Marshal32(info.endTime),
		})
	}

	for _, gift := range p.eventLimitGifts {
		proto.EventLimitGifts = append(proto.EventLimitGifts, gift.Encode())
	}

	for id, times := range p.timeLimitGiftBuyTimes.internalMap {
		proto.TimeLimitGiftBuyTimes = append(proto.TimeLimitGiftBuyTimes, &shared_proto.TimeLimitGiftBuyTimesProto {
			Id: u64.Int32(id),
			Times: u64.Int32(times),
		})
	}

	return proto
}

func (p *hero_promotion) Unmarshal(proto *server_proto.HeroPromotionServerProto, datas *config.ConfigDatas, ctime time.Time) {
	if proto == nil {
		return
	}
	if len(proto.BuyTimeLimitGift) > 0 {
		for _, info := range proto.BuyTimeLimitGift {
			groupId := u64.FromInt32(info.Id)
			data := datas.GetTimeLimitGiftGroupData(groupId)
			if data == nil {
				logrus.Debugf("加载定时时限礼包组数据时没有配置的礼包组id:%d", groupId)
				continue
			}
			endTime := timeutil.Unix32(info.EndTime)
			if ctime.After(endTime) {
				continue
			}
			p.timeLimitGifts[groupId] = newTimeLimitGiftInfo(groupId, u64.FromInt32(info.Index), endTime)
		}
	}
	if len(proto.EventLimitGifts) > 0 {
		for _, gift := range proto.EventLimitGifts {
			gid := u64.FromInt32(gift.Id)
			data := datas.GetEventLimitGiftData(gid)
			if data == nil {
				logrus.Debugf("加载事件礼包数据时没有配置的礼包id:%d", gid)
				continue
			}
			endTime := timeutil.Unix32(gift.EndTime)
			if ctime.After(endTime) {
				continue
			}
			if p.eventLimitGiftBuyTimes.Get(gid) >= data.BuyLimit { // 策划可能修改了配置
				continue
			}
			p.eventLimitGifts[data.Id] = newEventLimitGift(data, endTime)
		}
	}
}

func (p *hero_promotion) EncodeServer() *server_proto.HeroPromotionServerProto {
	proto := &server_proto.HeroPromotionServerProto{}
	if len(p.eventLimitGifts) > 0 {
		for _, gift := range p.eventLimitGifts {
			proto.EventLimitGifts = append(proto.EventLimitGifts, gift.EncodeServer())
		}
	}
	return proto
}

func (p *hero_promotion) TrySetEventLimitGift(data *promdata.EventLimitGiftData, ctime time.Time) (*eventLimitGift, bool) {
	if gift := p.eventLimitGifts[data.Id]; gift == nil || ctime.After(gift.endTime) {
		gift = newEventLimitGift(data, ctime.Add(data.TimeDuration))
		p.eventLimitGifts[data.Id] = gift
		return gift, true
	}
	return nil, false
}

func (p *hero_promotion) GetEventLimitGiftWithCheck(id uint64, ctime time.Time) *eventLimitGift {
	gift := p.eventLimitGifts[id]
	if gift == nil {
		return nil
	}
	if ctime.After(gift.endTime) {
		delete(p.eventLimitGifts, id)
		return nil
	}
	return gift
}

func (p *hero_promotion) RemoveEventLimitGift(id uint64) {
	delete(p.eventLimitGifts, id)
}

func (p *hero_promotion) GetTimeLimitGiftBuyTimes(id uint64) uint64 {
	return p.timeLimitGiftBuyTimes.Get(id)
}

func (p *hero_promotion) IncTimeLimitGiftBuyTimes(id uint64) {
	p.timeLimitGiftBuyTimes.Increse(id)
}

func (p *hero_promotion) GetEventLimitGiftBuyTimes(id uint64) uint64 {
	return p.eventLimitGiftBuyTimes.Get(id)
}

func (p *hero_promotion) IncEventLimitGiftBuyTimes(id uint64) {
	p.eventLimitGiftBuyTimes.Increse(id)
}

func (p *hero_promotion) GetTimeLimitGiftIndex(grpId uint64, endTime time.Time) uint64 {
	g := p.timeLimitGifts[grpId]
	if g == nil {
		return 0
	}
	if !g.endTime.Equal(endTime) {
		delete(p.timeLimitGifts, grpId)
		return 0
	}
	return g.index
}

func (p *hero_promotion) SetTimeLimitGiftIndex(grpId, index uint64, endTime time.Time) {
	g := p.timeLimitGifts[grpId]
	if g == nil {
		g = newTimeLimitGiftInfo(grpId, index, endTime)
		p.timeLimitGifts[grpId] = g
	} else {
		g.index = index
	}
}
