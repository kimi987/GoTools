package buff

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/farm"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
)

func newBuffEffectManager(heroData iface.HeroDataService, world iface.WorldService, datas iface.ConfigDatas, time iface.TimeService, farm iface.FarmService) *BuffEffectManager {
	m := &BuffEffectManager{}
	srv := &services{
		heroData: heroData,
		world:    world,
		datas:    datas,
		time:     time,
	}
	m.srv = srv

	m.buffMap = make(map[shared_proto.BuffEffectType]BuffEffect, len(shared_proto.BuffEffectType_name))
	for id := range shared_proto.BuffEffectType_name {
		t := shared_proto.BuffEffectType(id)
		switch t {
		case shared_proto.BuffEffectType_Buff_ET_captain_train:
			m.buffMap[t] = NewBuffEffectCaptainTrain(srv)
		case shared_proto.BuffEffectType_Buff_ET_farm_harvest:
			m.buffMap[t] = NewBuffEffectFarmHarvest(srv, farm)
		case shared_proto.BuffEffectType_Buff_ET_sprite_stat:
			m.buffMap[t] = NewBuffEffectSpriteStat(srv)
		case shared_proto.BuffEffectType_Buff_ET_tax:
			m.buffMap[t] = NewBuffEffectTax(srv)
		case shared_proto.BuffEffectType_Buff_ET_battle_mian:
			m.buffMap[t] = NewBuffEffectBattleMian(srv)
		}
	}
	return m
}

type services struct {
	heroData iface.HeroDataService
	world    iface.WorldService
	datas    iface.ConfigDatas
	time     iface.TimeService
}

type BuffEffectManager struct {
	srv     *services
	buffMap map[shared_proto.BuffEffectType]BuffEffect
}

type BuffEffect interface {
	afterAdd(heroId int64, buff *entity.BuffInfo)
	afterDel(heroId int64, buff *entity.BuffInfo)
	afterUpdate(heroId int64, oldBuff, newBuff *entity.BuffInfo)
	updatePerSecond(hero *entity.Hero, result herolock.LockResult, buff *entity.BuffInfo)
}

type BuffEffectSuper struct {
	srv *services
}

func NewBuffEffectSuper(srv *services) *BuffEffectSuper {
	b := &BuffEffectSuper{}
	b.srv = srv

	return b
}

func (m *BuffEffectSuper) updatePerSecond(hero *entity.Hero, result herolock.LockResult, buff *entity.BuffInfo) {}

// Buff_ET_captain_train
type BuffEffectCaptainTrain struct {
	*BuffEffectSuper
}

func NewBuffEffectCaptainTrain(srv *services) *BuffEffectCaptainTrain {
	b := &BuffEffectCaptainTrain{}
	b.BuffEffectSuper = NewBuffEffectSuper(srv)
	return b
}

func (m *BuffEffectCaptainTrain) afterAdd(heroId int64, buff *entity.BuffInfo) {}

func (m *BuffEffectCaptainTrain) afterDel(heroId int64, buff *entity.BuffInfo) {
	m.srv.heroData.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		buildingData := hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN)
		if buildingData == nil {
			return
		}

		toAdd := heromodule.CalcTranBuffExp(hero, buff, m.srv.datas, m.srv.time.CurrentTime())
		hero.Military().AddReservedExp(toAdd)
		result.Changed()
	})

	return
}

func (m *BuffEffectCaptainTrain) afterUpdate(heroId int64, oldBuff, newBuff *entity.BuffInfo) {
	m.srv.heroData.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		buildingData := hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN)
		if buildingData == nil {
			return
		}
		if oldBuff == nil {
			return
		}

		toAdd := heromodule.CalcTranBuffExp(hero, oldBuff, m.srv.datas, m.srv.time.CurrentTime())
		hero.Military().AddReservedExp(toAdd)
		result.Changed()
	})
}

// Buff_ET_sprite_stat
type BuffEffectSpriteStat struct {
	*BuffEffectSuper
}

func NewBuffEffectSpriteStat(srv *services) *BuffEffectSpriteStat {
	b := &BuffEffectSpriteStat{}
	b.BuffEffectSuper = NewBuffEffectSuper(srv)
	return b
}

func (m *BuffEffectSpriteStat) afterAdd(heroId int64, buff *entity.BuffInfo) {
	m.srv.heroData.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		hero.Buff().Update()

		if !buff.EffectData.PvpBuff {
			heromodule.UpdateOnBuffChanged(hero, result)
		}

		result.Changed()
	})
}

func (m *BuffEffectSpriteStat) afterDel(heroId int64, buff *entity.BuffInfo) {
	m.srv.heroData.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		hero.Buff().Update()

		if !buff.EffectData.PvpBuff {
			heromodule.UpdateOnBuffChanged(hero, result)
		}

		result.Changed()
	})
}

func (m *BuffEffectSpriteStat) afterUpdate(heroId int64, oldBuff, newBuff *entity.BuffInfo) {
	m.srv.heroData.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		hero.Buff().Update()

		if !newBuff.EffectData.PvpBuff {
			heromodule.UpdateOnBuffChanged(hero, result)
		}

		result.Changed()
	})
}

// Buff_ET_farm_harvest
type BuffEffectFarmHarvest struct {
	*BuffEffectSuper
	srv  *services
	farm iface.FarmService
}

func NewBuffEffectFarmHarvest(srv *services, farm iface.FarmService) *BuffEffectFarmHarvest {
	b := &BuffEffectFarmHarvest{}
	b.BuffEffectSuper = NewBuffEffectSuper(srv)
	b.srv = srv
	b.farm = farm
	return b
}

func (m *BuffEffectFarmHarvest) afterAdd(heroId int64, buff *entity.BuffInfo) {
	m.farm.ReduceRipeTimePercent(heroId, nil, buff)
	m.srv.world.Send(heroId, farm.FARM_IS_UPDATE_S2C)
}

func (m *BuffEffectFarmHarvest) afterDel(heroId int64, buff *entity.BuffInfo) {
	m.farm.ReduceRipeTimePercent(heroId, buff, nil)
	m.srv.world.Send(heroId, farm.FARM_IS_UPDATE_S2C)
}

func (m *BuffEffectFarmHarvest) afterUpdate(heroId int64, oldBuff, newBuff *entity.BuffInfo) {
	m.farm.ReduceRipeTimePercent(heroId, oldBuff, newBuff)
	m.srv.world.Send(heroId, farm.FARM_IS_UPDATE_S2C)
}

// Buff_ET_tax
type BuffEffectTax struct {
	*BuffEffectSuper
}

func NewBuffEffectTax(srv *services) *BuffEffectTax {
	b := &BuffEffectTax{}
	b.BuffEffectSuper = NewBuffEffectSuper(srv)
	return b
}

func (m *BuffEffectTax) afterAdd(heroId int64, buff *entity.BuffInfo) {}

func (m *BuffEffectTax) afterDel(heroId int64, buff *entity.BuffInfo) {}

func (m *BuffEffectTax) afterUpdate(heroId int64, oldBuff, newBuff *entity.BuffInfo) {}

// Buff_ET_battle_mian
type BuffEffectBattleMian struct {
	*BuffEffectSuper
}

func NewBuffEffectBattleMian(srv *services) *BuffEffectBattleMian {
	b := &BuffEffectBattleMian{}
	b.BuffEffectSuper = NewBuffEffectSuper(srv)
	return b
}

func (m *BuffEffectBattleMian) afterAdd(heroId int64, buff *entity.BuffInfo) {}

func (m *BuffEffectBattleMian) afterDel(heroId int64, buff *entity.BuffInfo) {}

func (m *BuffEffectBattleMian) afterUpdate(heroId int64, oldBuff, newBuff *entity.BuffInfo) {}
