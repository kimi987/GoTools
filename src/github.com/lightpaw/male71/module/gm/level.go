package gm

import (
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/u64"
	"strconv"
)

func (m *GmModule) newLevelGmGroup() *gm_group {
	return &gm_group{
		tab: "等级",
		handler: []*gm_handler{
			newHeroIntHandler("加君主等级", "1", m.addHeroLevel),
			newIntHandler("加主城等级", "1", m.addHomeLevel),
			newHeroIntHandler("加武将等级(不含外出的)", "100", m.addCaptainLevel),
			newHeroIntHandler("加VIP经验", "100", m.addVipExp),
			//newIntHandler("领VIP每日礼包", "1", m.collectVipDailyPrize),
			//newIntHandler("领VIP专属礼包", "1", m.collectVipLevelPrize),
			newIntHandler("解锁建筑队", "1", m.unlockWorker),
			newIntHandler("购买推图次数", "20106", m.buyDungeonTimes),
			newStringHandler("推图扫荡", "20106 5", m.dungeonAutoChallenge),
			newHeroIntHandler("加体力", "100", m.addHeroSp),
		},
	}
}

func (m *GmModule) addHeroSp(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	heromodule.AddSp(m.hctx, hero, result, u64.FromInt64(amount))
}


func (m *GmModule) addVipExp(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	heromodule.AddVipExp(m.hctx, hero, result, u64.FromInt64(amount), m.time.CurrentTime())
}

func (m *GmModule) collectVipDailyPrize(level int64, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleVip_c2s_vip_collect_daily_prize(amount string, hc iface.HeroController)
	}); ok {
		im.handleVip_c2s_vip_collect_daily_prize(strconv.FormatInt(level, 10), hc)
	}

}

func (m *GmModule) collectVipLevelPrize(level int64, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleVip_c2s_vip_collect_level_prize(amount string, hc iface.HeroController)
	}); ok {
		im.handleVip_c2s_vip_collect_level_prize(strconv.FormatInt(level, 10), hc)
	}

}

func (m *GmModule) unlockWorker(pos int64, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleDomestic_c2s_worker_unlock(amount string, hc iface.HeroController)
	}); ok {
		im.handleDomestic_c2s_worker_unlock(strconv.FormatInt(pos, 10), hc)
	}

}

func (m *GmModule) addHeroLevel(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	for i := 0; i < int(amount); i++ {
		heromodule.AddExp(m.hctx, hero, result, hero.LevelData().Sub.UpgradeExp, m.time.CurrentTime())
	}
}

func (m *GmModule) addHeroLevelTo(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	for i := hero.Level(); i < uint64(amount); i++ {
		if hero.LevelData().NextLevel() == nil {
			return
		}
		heromodule.AddExp(m.hctx, hero, result, hero.LevelData().Sub.UpgradeExp, m.time.CurrentTime())
	}
}

func (m *GmModule) addHomeLevel(amount int64, hc iface.HeroController) {
	c := int64(1)
	if amount > 1 {
		c = amount
	}

	var region int64
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		region = hero.BaseRegion()
		return
	})

	if region == 0 {
		return
	}

	realm := m.realmService.GetBigMap()
	if realm == nil {
		return
	}

	for i := 0; i < int(c); i++ {
		var levelData *domestic_data.BaseLevelData
		hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
			levelData = m.datas.GetBaseLevelData(hero.BaseLevel())
			if levelData != nil {
				hero.AddProsperityCapcity(u64.Sub(levelData.UpgradeProsperity, hero.ProsperityCapcity()))
				return true
			}
			return
		})

		if levelData == nil {
			break
		}

		m.addProsperity(int64(levelData.UpgradeProsperity), hc)

		realm.UpgradeBase(hc)
	}
}

func (m *GmModule) addCaptainLevel(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	for _, c := range hero.Military().Captains() {
		if c.IsOutSide() {
			continue
		}

		if c.LevelData().NextLevel() == nil {
			continue
		}

		for i := int64(0); i < amount; i++ {
			if c.LevelData().NextLevel() == nil {
				break
			}

			heromodule.AddCaptainExp(m.hctx, hero, result, c, c.LevelData().UpgradeExp, m.time.CurrentTime())
		}
	}
}

func (m *GmModule) buyDungeonTimes(pos int64, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleVip_c2s_vip_buy_dungeon_times(amount string, hc iface.HeroController)
	}); ok {
		im.handleVip_c2s_vip_buy_dungeon_times(strconv.FormatInt(pos, 10), hc)
	}

}

func (m *GmModule) dungeonAutoChallenge(pos string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleDungeon_c2s_auto_challenge(amount string, hc iface.HeroController)
	}); ok {
		im.handleDungeon_c2s_auto_challenge(pos, hc)
	}

}

