package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

type hero_reservation struct {
	gold  uint64
	food  uint64
	wood  uint64
	stone uint64

	jade    uint64
	jadeOre uint64

	yuanbao  uint64
	dianquan uint64
	yinliang uint64

	guildContributionCoin uint64

	goodsMap map[uint64]uint64 // key是id，value是个数

	lastReseveTime time.Time // 最后一次预约时间
}

func (r *hero_reservation) IsEmpty() bool {
	return r.gold <= 0 && r.food <= 0 && r.wood <= 0 && r.stone <= 0 &&
		r.jade <= 0 &&
		r.yuanbao <= 0 &&
		r.dianquan <= 0 &&
		r.yinliang <= 0 &&
		r.guildContributionCoin <= 0 &&
		len(r.goodsMap) <= 0
}

func (r *hero_reservation) TryClearAndGetBackReserveResult(ctime time.Time) *ReserveResult {

	if timeutil.Marshal64(r.lastReseveTime) <= 0 {
		return nil
	}

	if ctime.Before(r.lastReseveTime) {
		return nil
	}

	r.lastReseveTime = time.Time{}

	if r.IsEmpty() {

		// 没有要返还的
		return nil
	}

	result := &ReserveResult{}

	// 资源
	result.reserveResource(r.gold, r.food, r.wood, r.stone)
	r.gold = 0
	r.food = 0
	r.wood = 0
	r.stone = 0

	result.jade += r.jade
	r.jade = 0

	result.jadeOre += r.jadeOre
	r.jadeOre = 0

	// 元宝
	result.reserveYuanbao(r.yuanbao)
	r.yuanbao = 0

	// 点券
	result.reserveDianquan(r.dianquan)
	r.dianquan = 0

	// 银两
	result.reserveYinliang(r.yinliang)
	r.yinliang = 0

	// 帮贡币
	result.guildContributionCoin += r.guildContributionCoin
	r.guildContributionCoin = 0

	// 物品
	for k, v := range r.goodsMap {
		result.reserveGoods(k, v)
		delete(r.goodsMap, k)
	}

	return result
}

func (r *hero_reservation) encode() *server_proto.HeroReservationProto {
	proto := &server_proto.HeroReservationProto{}

	proto.Gold = r.gold
	proto.Food = r.food
	proto.Wood = r.wood
	proto.Stone = r.stone

	proto.Jade = r.jade

	proto.Yuanbao = r.yuanbao
	proto.Dianquan = r.dianquan
	proto.Yinliang = r.yinliang

	proto.GuildContributionCoin = r.guildContributionCoin

	proto.Goods = r.goodsMap

	proto.LastReserveTime = timeutil.Marshal64(r.lastReseveTime)

	return proto
}

func (r *hero_reservation) unmarshal(proto *server_proto.HeroReservationProto) {
	if proto == nil {
		return
	}

	r.gold = proto.Gold
	r.food = proto.Food
	r.wood = proto.Wood
	r.stone = proto.Stone

	r.jade = proto.Jade
	r.jadeOre = proto.JadeOre

	r.yuanbao = proto.Yuanbao
	r.dianquan = proto.Dianquan
	r.yinliang = proto.Yinliang

	r.guildContributionCoin = proto.GuildContributionCoin

	u64.CopyMapTo(r.goodsMap, proto.GetGoods())

	r.lastReseveTime = timeutil.Unix64(proto.LastReserveTime)
}

func (hero *hero_reservation) ReserveCost(cost *resdata.Cost, ctime time.Time) *ReserveResult {

	result := &ReserveResult{}
	result.reserveCost(cost)

	hero.reserveResult(result, ctime)

	return result
}

func (hero *hero_reservation) ReserveGoods(goodsId, goodsCount uint64, ctime time.Time) *ReserveResult {

	result := &ReserveResult{}
	result.reserveGoods(goodsId, goodsCount)

	hero.reserveResult(result, ctime)

	return result
}

func (hero *hero_reservation) ReserveGoodsOrBuy(goodsId, goodsCount, yinliang, dianquan, yuanbao uint64, ctime time.Time) *ReserveResult {

	result := &ReserveResult{}
	if goodsCount > 0 {
		result.reserveGoods(goodsId, goodsCount)
	}

	if yinliang > 0 {
		result.reserveDianquan(yinliang)
	}

	if dianquan > 0 {
		result.reserveDianquan(dianquan)
	}

	if yuanbao > 0 {
		result.reserveYuanbao(yuanbao)
	}

	hero.reserveResult(result, ctime)

	return result
}

func (hero *hero_reservation) ReserveYuanbao(yuanbao uint64, ctime time.Time) *ReserveResult {

	result := &ReserveResult{}
	result.reserveYuanbao(yuanbao)

	hero.reserveResult(result, ctime)

	return result
}

func (hero *hero_reservation) ReserveDianquan(diqnquan uint64, ctime time.Time) *ReserveResult {

	result := &ReserveResult{}
	result.reserveDianquan(diqnquan)

	hero.reserveResult(result, ctime)

	return result
}

func (hero *hero_reservation) ReserveYinliang(yinliang uint64, ctime time.Time) *ReserveResult {

	result := &ReserveResult{}
	result.reserveYinliang(yinliang)

	hero.reserveResult(result, ctime)

	return result
}

func (r *hero_reservation) reserveResult(result *ReserveResult, ctime time.Time) {

	r.gold += result.gold
	r.food += result.food
	r.wood += result.wood
	r.stone += result.stone

	r.jade += result.jade
	r.jadeOre += result.jadeOre

	r.yuanbao += result.yuanbao
	r.dianquan += result.dianquan
	r.yinliang += result.yinliang

	r.guildContributionCoin += result.guildContributionCoin

	n := imath.Min(len(result.goodsIds), len(result.goodsCounts))
	for i := 0; i < n; i++ {
		r.goodsMap[result.goodsIds[i]] += result.goodsCounts[i]
	}

	r.lastReseveTime = ctime.Add(time.Hour) // 一小时后返还
}

// 扣掉英雄预约的数据，更新实际扣了多少（返还的时候，不会多给了）
func (r *hero_reservation) ConfirmResult(result *ReserveResult) {

	// 更新实际能扣多少，之后最多加回这么多
	result.gold = u64.Min(result.gold, r.gold)
	result.food = u64.Min(result.food, r.food)
	result.wood = u64.Min(result.wood, r.wood)
	result.stone = u64.Min(result.stone, r.stone)

	result.jade = u64.Min(result.jade, r.jade)
	result.jadeOre = u64.Min(result.jadeOre, r.jadeOre)

	result.yuanbao = u64.Min(result.yuanbao, r.yuanbao)
	result.dianquan = u64.Min(result.dianquan, r.dianquan)
	result.yinliang = u64.Min(result.yinliang, r.yinliang)

	result.guildContributionCoin = u64.Min(result.guildContributionCoin, r.guildContributionCoin)

	n := imath.Min(len(result.goodsIds), len(result.goodsCounts))
	for i := 0; i < n; i++ {
		goodsId := result.goodsIds[i]
		result.goodsCounts[i] = u64.Min(result.goodsCounts[i], r.goodsMap[goodsId])
	}

	// 实际从预约对象中扣掉
	r.gold = u64.Sub(r.gold, result.gold)
	r.food = u64.Sub(r.food, result.food)
	r.wood = u64.Sub(r.wood, result.wood)
	r.stone = u64.Sub(r.stone, result.stone)

	r.jade = u64.Sub(r.jade, result.jade)
	r.jadeOre = u64.Sub(r.jadeOre, result.jadeOre)

	r.yuanbao = u64.Sub(r.yuanbao, result.yuanbao)
	r.dianquan = u64.Sub(r.dianquan, result.dianquan)
	r.yinliang = u64.Sub(r.yinliang, result.yinliang)

	r.guildContributionCoin = u64.Sub(r.guildContributionCoin, result.guildContributionCoin)

	for i := 0; i < n; i++ {
		toReduce := result.goodsCounts[i]
		if toReduce <= 0 {
			continue
		}

		goodsId := result.goodsIds[i]

		count, exist := r.goodsMap[goodsId]
		if exist {
			if count <= toReduce {
				delete(r.goodsMap, goodsId)
			} else {
				r.goodsMap[goodsId] = u64.Sub(count, toReduce)
			}
		}
	}

}

type ReserveResult struct {
	gold  uint64
	food  uint64
	wood  uint64
	stone uint64

	jade    uint64
	jadeOre uint64

	yuanbao  uint64
	dianquan uint64
	yinliang uint64

	guildContributionCoin uint64

	goodsIds    []uint64
	goodsCounts []uint64
}

func (r *ReserveResult) GetResource() (uint64, uint64, uint64, uint64) {
	return r.gold, r.food, r.wood, r.stone
}

func (r *ReserveResult) GetYuanbao() uint64 {
	return r.yuanbao
}

func (r *ReserveResult) GetDianquan() uint64 {
	return r.dianquan
}

func (r *ReserveResult) GetYinliang() uint64 {
	return r.yinliang
}

func (r *ReserveResult) GetGuildContributionCoin() uint64 {
	return r.guildContributionCoin
}

func (r *ReserveResult) GetGoodsIdCounts() ([]uint64, []uint64) {
	return r.goodsIds, r.goodsCounts
}

func (result *ReserveResult) reserveCost(cost *resdata.Cost) {

	result.reserveResource(cost.GetResource())

	// 玉璧
	result.jade += cost.Jade
	result.jadeOre += cost.JadeOre

	result.reserveYuanbao(cost.Yuanbao)
	result.reserveDianquan(cost.Dianquan)
	result.reserveYinliang(cost.Yinliang)

	result.guildContributionCoin += cost.GuildContributionCoin

	n := imath.Min(len(cost.Goods), len(cost.GoodsCount))
	for i := 0; i < n; i++ {
		if g := cost.Goods[i]; g != nil {
			result.reserveGoods(g.Id, cost.GoodsCount[i])
		}
	}

}

func (result *ReserveResult) reserveResource(gold, food, wood, stone uint64) {
	result.gold += gold
	result.food += food
	result.wood += wood
	result.stone += stone
}

func (result *ReserveResult) reserveYuanbao(yuanbao uint64) {
	result.yuanbao += yuanbao
}

func (result *ReserveResult) reserveDianquan(dianquan uint64) {
	result.dianquan += dianquan
}

func (result *ReserveResult) reserveYinliang(yinliang uint64) {
	result.yinliang += yinliang
}

func (result *ReserveResult) reserveGoods(goodsId, goodsCount uint64) {
	if goodsCount > 0 {
		result.goodsIds = append(result.goodsIds, goodsId)
		result.goodsCounts = append(result.goodsCounts, goodsCount)
	}
}

func (r *ReserveResult) Print(msg string) {
	logrus.WithField("gold", r.gold).
		WithField("food", r.food).
		WithField("wood", r.wood).
		WithField("stone", r.stone).
		WithField("jade", r.jade).
		WithField("jadeOre", r.jadeOre).
		WithField("yuanbao", r.yuanbao).
		WithField("dianquan", r.dianquan).
		WithField("yinliang", r.yinliang).
		WithField("guildContributionCoin", r.guildContributionCoin).
		WithField("goodsIds", r.goodsIds).
		WithField("goodsCounts", r.goodsCounts).
		Error(msg)
}
