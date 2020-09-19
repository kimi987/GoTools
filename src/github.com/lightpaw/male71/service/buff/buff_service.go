package buff

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
)

func NewBuffService(heroData iface.HeroDataService, world iface.WorldService, datas iface.ConfigDatas, time iface.TimeService, farm iface.FarmService) *BuffService {
	m := &BuffService{}
	m.heroData = heroData
	m.world = world
	m.datas = datas
	m.time = time

	m.manager = newBuffEffectManager(heroData, world, datas, time, farm)
	return m
}

//gogen:iface
type BuffService struct {
	heroData iface.HeroDataService
	world    iface.WorldService
	datas    iface.ConfigDatas
	time     iface.TimeService

	manager *BuffEffectManager
}

func (m *BuffService) Cancel(heroId int64, buffs []*entity.BuffInfo) (succ bool) {
	for _, b := range buffs {
		m.afterCancel(heroId, b)
		m.world.Send(heroId, misc.NewS2cUpdateBuffNoticeMsg(int32(b.EffectData.Group), nil))
	}
	succ = true
	return
}

func (m *BuffService) CancelGroup(heroId int64, group uint64) (succ bool) {
	var buff *entity.BuffInfo
	m.heroData.FuncNotError(heroId, func(hero *entity.Hero) (heroChanged bool) {
		buff = hero.Buff().Buff(group)
		hero.Buff().Del(group, m.time.CurrentTime())
		return
	})
	if buff == nil {
		return
	}
	return m.afterCancel(heroId, buff)
}

func (m *BuffService) afterCancel(heroId int64, buff *entity.BuffInfo) (succ bool) {
	if effect, ok := m.manager.buffMap[buff.EffectData.EffectType]; ok {
		if buff != nil {
			effect.afterDel(heroId, buff)
		}
		succ = true
	} else {
		logrus.Errorf("BuffService.Cancel, 找不到类型：%v", buff.EffectData.EffectType)
	}

	return
}

// 新增或替换buff
func (m *BuffService) AddBuffToSelf(buff *data.BuffEffectData, heroId int64) (succ bool) {
	ctime := m.time.CurrentTime()
	var newBuff, oldBuff *entity.BuffInfo
	if m.heroData.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		succ, newBuff, oldBuff = hero.Buff().Add(buff, heroId, ctime)
		if succ {
			result.Changed()
			result.Ok()
		}
	}) {
		return
	}

	if newBuff == nil {
		return
	}

	if effect, ok := m.manager.buffMap[buff.EffectType]; ok {
		if oldBuff != nil {
			effect.afterUpdate(heroId, oldBuff, newBuff)
		} else {
			effect.afterAdd(heroId, newBuff)
		}
	} else {
		logrus.Errorf("BuffService.AddBuffToSelf 找不到 manager. type:%v", buff.EffectType)
		return
	}

	m.world.Send(heroId, misc.NewS2cUpdateBuffNoticeMsg(int32(buff.Group), newBuff.Encode()))

	return true
}

// 新增或替换buff
func (m *BuffService) AddBuffToTarget(buff *data.BuffEffectData, heroId, targetId int64) (succ bool) {
	ctime := m.time.CurrentTime()
	var newBuff, oldBuff *entity.BuffInfo
	if m.heroData.FuncWithSend(targetId, func(hero *entity.Hero, result herolock.LockResult) {
		succ, newBuff, oldBuff = hero.Buff().Add(buff, heroId, ctime)
		if succ {
			result.Changed()
			result.Ok()
		}
	}) {
		return
	}

	if newBuff == nil {
		return
	}

	if effect, ok := m.manager.buffMap[buff.EffectType]; ok {
		if oldBuff != nil {
			effect.afterUpdate(targetId, oldBuff, newBuff)
		} else {
			effect.afterAdd(targetId, newBuff)
		}
	} else {
		logrus.Errorf("BuffService.AddBuffToTarget 找不到 manager. type:%v", buff.EffectType)
		return
	}
	m.world.Send(targetId, misc.NewS2cUpdateBuffNoticeMsg(int32(buff.Group), newBuff.Encode()))

	return true
}

func (m *BuffService) UpdatePerSecond(hero *entity.Hero, result herolock.LockResult, buff *entity.BuffInfo) {
	if buff == nil {
		logrus.Debugf("BuffService.UpdatePerSecond, buff == nil")
		return
	}

	if effect, ok := m.manager.buffMap[buff.EffectData.EffectType]; ok {
		effect.updatePerSecond(hero, result, buff)
	} else {
		logrus.Debugf("BuffService.UpdatePerSecond, effect == nil")
	}
}
