package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/config/promdata"
)

func NewHeroMisc(maps *hero_maps) *HeroMisc {
	misc := &HeroMisc{}
	misc.broadcastSendType = maps.getOrCreateKeys(server_proto.HeroMapCategory_sys_broadcast_send_type, false)
	misc.collectedChargePrizes = maps.getOrCreateKeys(server_proto.HeroMapCategory_collected_charge_prizes, false)
	misc.boughtBargainTimes = maps.getOrCreateMap(server_proto.HeroMapCategory_daily_bought_bargain_times, true)
	misc.collectedDurationCardDailyPrizes = maps.getOrCreateKeys(server_proto.HeroMapCategory_daily_collected_duration_card_prizes, true)
	misc.durationCards = make(map[uint64]time.Time)
	misc.firstChargedObjs = maps.getOrCreateKeys(server_proto.HeroMapCategory_first_charged_objs, false)
	misc.nextTime4Operations = make(map[server_proto.OperationCDType]time.Time)

	return misc
}

// 零碎的额外数据
type HeroMisc struct {
	FarmUnlockAnimationBaseLevel uint64

	broadcastSendType                *herokeys
	// 领过的充值奖励，主键对应 充值奖励.txt
	collectedChargePrizes            *herokeys
	// 当前充值总金额
	chargeAmount                     uint64
	// 最近一次充值的时间
	lastChargeTime                   int64
	// 今日特惠已购(领)次数
	boughtBargainTimes               *heromap
	// 尊享卡结束时间（如果是永久卡则是记录开卡时间）
	durationCards                    map[uint64]time.Time
	// 尊享卡当天奖励领取记录
	collectedDurationCardDailyPrizes *herokeys

	// 首充过的充值项，主键对应 充值项.txt，该列表会根据服务器的首充重置而重置
	firstChargedObjs                 *herokeys

	// 不需要存库<可操作编号，下次允许操作的时间>
	nextTime4Operations              map[server_proto.OperationCDType]time.Time
}

func (misc *HeroMisc) GetNextOperationTime(t server_proto.OperationCDType) time.Time {
	return misc.nextTime4Operations[t]
}

func (misc *HeroMisc) SetNextOperationTime(t server_proto.OperationCDType, time time.Time) {
	misc.nextTime4Operations[t] = time
}

func (misc *HeroMisc) EncodeServer() *server_proto.HeroMiscServerProto {
	proto := &server_proto.HeroMiscServerProto{}
	proto.FarmShowUnlockAnimationBaseLevel = misc.FarmUnlockAnimationBaseLevel
	proto.ChargeAmount = misc.chargeAmount
	proto.LastChargeTime = misc.lastChargeTime
	proto.DurationCards = make(map[uint64]int64)
	for k, v := range misc.durationCards {
		proto.DurationCards[k] = timeutil.Marshal64(v)
	}
	return proto
}

func (misc *HeroMisc) Unmarshal(proto *server_proto.HeroMiscServerProto) {
	if proto == nil {
		return
	}
	misc.FarmUnlockAnimationBaseLevel = proto.FarmShowUnlockAnimationBaseLevel
	misc.chargeAmount = proto.ChargeAmount
	misc.lastChargeTime = proto.LastChargeTime
	if len(proto.DurationCards) > 0 {
		for k, v := range proto.DurationCards {
			misc.durationCards[k] = timeutil.Unix64(v)
		}
	}
}

func (misc *HeroMisc) AddBroadcastSendType(bcType, subType, condition uint64) {
	misc.broadcastSendType.Add(combineBroadcastType(bcType, subType, condition))
}

func combineBroadcastType(bcType, subType, cond uint64) uint64 {
	return bcType<<56 | subType<<26 | cond
}

func (misc *HeroMisc) BroadcastSendTypeExist(bcType, subType, condition uint64) bool {
	return misc.broadcastSendType.Exist(combineBroadcastType(bcType, subType, condition))
}

func (misc *HeroMisc) EncodeClient() *shared_proto.ChargeProto {
	p := &shared_proto.ChargeProto{}
	p.ChargeAmount = u64.Int32(misc.chargeAmount)
	for k := range misc.collectedChargePrizes.internalMap {
		p.CollectedChargePrizes = append(p.CollectedChargePrizes, u64.Int32(k))
	}
	for k, v := range misc.boughtBargainTimes.internalMap {
		p.BoughtBargainTimes = append(p.BoughtBargainTimes, &shared_proto.BroughtBargainProto {
			Id: u64.Int32(k),
			Times: u64.Int32(v),
		})
	}
	for k := range misc.collectedDurationCardDailyPrizes.internalMap {
		p.CollectedDurationCardDailyPrizes = append(p.CollectedDurationCardDailyPrizes, u64.Int32(k))
	}
	for k, v := range misc.durationCards {
		p.DurationCards = append(p.DurationCards, &shared_proto.DurationCardProto {
			Id: u64.Int32(k),
			EndTime: timeutil.Marshal32(v),
		})
	}
	for k := range misc.firstChargedObjs.internalMap {
		p.FirstChargedObjs = append(p.FirstChargedObjs, u64.Int32(k))
	}
	return p
}

func (misc *HeroMisc) IsCollectChargePrize(id uint64) bool {
	return misc.collectedChargePrizes.Exist(id)
}

func (misc *HeroMisc) TryCollectChargePrize(id uint64) bool {
	if misc.IsCollectChargePrize(id) {
		return false
	}
	misc.collectedChargePrizes.Add(id)
	return true
}

func (misc *HeroMisc) AddChargeAmount(amount uint64, ctime time.Time) {
	misc.chargeAmount += amount
	misc.lastChargeTime = timeutil.Marshal64(ctime)
}

func (misc *HeroMisc) ChargeAmount() uint64 {
	return misc.chargeAmount
}

func (misc *HeroMisc) LastChargeTime() int64 {
	return misc.lastChargeTime
}

func (misc *HeroMisc) IncreaseBoughtBargainTimes(id uint64) {
	 misc.boughtBargainTimes.Increse(id)
}

func (misc *HeroMisc) GetBoughtBargainTimes(id uint64) uint64 {
	return misc.boughtBargainTimes.Get(id)
}

// 如果是新激活或延期则以当天零点为基准设置endTime，返回endTime；如果激活的是永久卡，返回的是startTime
func (misc *HeroMisc) ActivateDurationCard(data *promdata.DurationCardData, ctime time.Time) time.Time {
	var time time.Time
	if data.Duration > 0 { // 非永久卡记录截止时间
		t, ok := misc.durationCards[data.Id]
		if !ok || t.Before(ctime) {
			midnightTime := timeutil.DailyTime.PrevTime(ctime)
			time = midnightTime.Add(data.Duration)
		} else {
			time = t.Add(data.Duration)
		}
	} else {
		time = ctime // 永久卡记录开卡时间
	}
	misc.durationCards[data.Id] = time
	return time
}

func (misc *HeroMisc) IsCollectDurationCardDailyPrize(id uint64) bool {
	return misc.collectedDurationCardDailyPrizes.Exist(id)
}

func (misc *HeroMisc) IsDurationCardTimeEnd(data *promdata.DurationCardData, ctime time.Time) bool {
	endTime, ok := misc.durationCards[data.Id]
	if !ok {
		return true
	}
	if data.Duration <= 0 {
		return false
	}
	return ctime.After(endTime)
}

func (misc *HeroMisc) SetCollectedDurationCardDailyPrize(id uint64) {
	misc.collectedDurationCardDailyPrizes.Add(id)
}

func (misc *HeroMisc) IsDurationCardActive(id uint64) bool {
	_, ok := misc.durationCards[id]
	return ok
}

func (misc *HeroMisc) SetFirstChargedObj(id uint64) {
	misc.firstChargedObjs.Add(id)
}

func (misc *HeroMisc) IsFirstChargedObj(id uint64) bool {
	return !misc.firstChargedObjs.Exist(id)
}
