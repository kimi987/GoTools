package vip

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	vipMsg "github.com/lightpaw/male7/gen/pb/vip"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/u64"
)

func NewVipModule(dep iface.ServiceDep) *VipModule {
	m := &VipModule{}
	m.dep = dep
	return m
}

//gogen:iface
type VipModule struct {
	dep iface.ServiceDep
}

//gogen:iface
func (m *VipModule) ProcessVipCollectDailyPrize(proto *vipMsg.C2SVipCollectDailyPrizeProto, hc iface.HeroController) {
	vipLevel := u64.FromInt32(proto.VipLevel)
	d := m.dep.Datas().GetVipLevelData(vipLevel)
	if d == nil {
		hc.Send(vipMsg.ERR_VIP_COLLECT_DAILY_PRIZE_FAIL_INVALID_LEVEL)
		return
	}

	if d.DailyPrize == nil {
		hc.Send(vipMsg.ERR_VIP_COLLECT_DAILY_PRIZE_FAIL_NO_DAILY_PRIZE)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		vip := hero.Vip()
		if vipLevel > vip.Level() {
			hc.Send(vipMsg.ERR_VIP_COLLECT_DAILY_PRIZE_FAIL_LEVEL_NOT_ENOUGH)
			return
		}

		if !vip.CanCollectDailyPrize(vipLevel) {
			hc.Send(vipMsg.ERR_VIP_COLLECT_DAILY_PRIZE_FAIL_ALREADY_COLLECTED)
			return
		}

		hctx := heromodule.NewContext(m.dep, operate_type.VipCollectDailyPrize)
		vip.CollectDailyPrize(vipLevel)
		heromodule.AddPrize(hctx, hero, result, d.DailyPrize, m.dep.Time().CurrentTime())

		result.Add(vipMsg.NewS2cVipCollectDailyPrizeMsg(u64.Int32(vipLevel)))
		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *VipModule) ProcessVipCollectLevelPrize(proto *vipMsg.C2SVipCollectLevelPrizeProto, hc iface.HeroController) {
	vipLevel := u64.FromInt32(proto.VipLevel)
	d := m.dep.Datas().GetVipLevelData(vipLevel)
	if d == nil {
		hc.Send(vipMsg.ERR_VIP_COLLECT_LEVEL_PRIZE_FAIL_INVALID_LEVEL)
		return
	}

	if d.LevelPrize == nil {
		hc.Send(vipMsg.ERR_VIP_COLLECT_LEVEL_PRIZE_FAIL_NO_LEVEL_PRIZE)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		vip := hero.Vip()
		if vipLevel > vip.Level() {
			hc.Send(vipMsg.ERR_VIP_COLLECT_LEVEL_PRIZE_FAIL_LEVEL_NOT_ENOUGH)
			return
		}

		if !vip.CanCollectLevelPrize(vipLevel) {
			hc.Send(vipMsg.ERR_VIP_COLLECT_LEVEL_PRIZE_FAIL_ALREADY_COLLECTED)
			return
		}

		hctx := heromodule.NewContext(m.dep, operate_type.VipCollectLevelPrize)
		if d.LevelPrizeCost != nil {
			if !heromodule.TryReduceCost(hctx, hero, result, d.LevelPrizeCost) {
				hc.Send(vipMsg.ERR_VIP_COLLECT_LEVEL_PRIZE_FAIL_COST_NOT_ENOUGH)
				return
			}
		}

		vip.CollectLevelPrize(vipLevel)
		heromodule.AddPrize(hctx, hero, result, d.LevelPrize, m.dep.Time().CurrentTime())

		result.Add(vipMsg.NewS2cVipCollectLevelPrizeMsg(u64.Int32(vipLevel)))
		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *VipModule) ProcessVipBuyDungeonTimes(proto *vipMsg.C2SVipBuyDungeonTimesProto, hc iface.HeroController) {
	dungeonId := u64.FromInt32(proto.DungeonId)
	data := m.dep.Datas().GetDungeonData(dungeonId)
	if data == nil {
		logrus.Debugf("ProcessVipBuyDungeonTimes, DungeonId：%v 不存在", proto.DungeonId)
		hc.Send(vipMsg.ERR_VIP_BUY_DUNGEON_TIMES_FAIL_DUNGEON_ID_INVALID)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		vipData := m.dep.Datas().GetVipLevelData(hero.VipLevel())
		vipMiscData := m.dep.Datas().VipMiscData()

		times := hero.Dungeon().VipAddedPassBoughtTimes(dungeonId)
		if vipData == nil || times >= vipData.DungeonMaxCostTimesLimit {
			result.Add(vipMsg.ERR_VIP_BUY_DUNGEON_TIMES_FAIL_VIP_LEVEL_LIMIT)
			return
		}

		cost, toAdd := vipMiscData.VipDungeonTimesCost(times + 1)
		hctx := heromodule.NewContext(m.dep, operate_type.VipBuyDungeonTimes)
		if !heromodule.TryReduceCost(hctx, hero, result, cost) {
			result.Add(vipMsg.ERR_VIP_BUY_DUNGEON_TIMES_FAIL_COST_NOT_ENOUGH)
			return
		}

		hero.Dungeon().AddVipAddedPassLimit(dungeonId, toAdd)

		result.Add(vipMsg.NewS2cVipBuyDungeonTimesMsg(proto.DungeonId, u64.Int32(hero.Dungeon().VipAddedPassBoughtTimes(dungeonId))))
		result.Changed()
		result.Ok()
	})
}
