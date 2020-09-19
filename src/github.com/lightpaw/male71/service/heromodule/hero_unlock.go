package heromodule

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
)

func IsHeroLocked(heroId int64, heroFunc func(id int64, f herolock.Func), guildSnapshot guildsnapshotdata.Getter,
	unlockData *data.UnlockCondition) (locked bool) {
	return IsLocked(func(f herolock.Func) {
		heroFunc(heroId, f)
	}, guildSnapshot, unlockData)
}

func IsLocked(heroFunc func(f herolock.Func), guildSnapshot guildsnapshotdata.Getter,
	unlockData *data.UnlockCondition) (locked bool) {

	if unlockData.IsEmptyCondition() {
		return false
	}

	var guildId int64
	heroFunc(func(hero *entity.Hero, err error) (heroChanged bool) {
		if err != nil {
			locked = true
			return
		}

		if hero.Level() < unlockData.RequiredHeroLevel {
			locked = true
			return
		}

		if hero.BaseLevel() < unlockData.RequiredBaseLevel {
			locked = true
			return
		}

		if hero.VipLevel() < unlockData.RequiredVipLevel {
			locked = true
			return
		}

		guildId = hero.GuildId()

		return
	})

	if locked {
		return
	}

	if unlockData.RequiredGuildLevel > 0 {
		if guildId == 0 {
			locked = true
			return
		}

		g := guildSnapshot(guildId)
		if g == nil || g.GuildLevel.Level < unlockData.RequiredGuildLevel {
			locked = true
			return
		}
	}

	return
}
