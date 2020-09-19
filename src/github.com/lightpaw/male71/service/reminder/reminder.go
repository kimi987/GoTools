package reminder

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/pbutil"
	"sync"
)

// 紧急提醒
func NewReminderService(guildSnapshotService iface.GuildSnapshotService, worldService iface.WorldService) *ReminderService {
	rs := &ReminderService{
		guildSnapshotService:    guildSnapshotService,
		worldService:            worldService,
		heroBeenAttackOrRobMap:  Newhero_been_attack_or_rob_status_map(),
		guildBeenAttackOrRobMap: Newguild_been_attack_or_rob_map(),
	}

	heromodule.RegisterHeroOnlineListener(rs)

	return rs
}

//gogen:iface
type ReminderService struct {
	guildSnapshotService    iface.GuildSnapshotService
	worldService            iface.WorldService
	heroBeenAttackOrRobMap  *hero_been_attack_or_rob_status_map
	guildBeenAttackOrRobMap *guild_been_attack_or_rob_map
}

func (rs *ReminderService) OnHeroOnline(hc iface.HeroController) {
	status := rs.getStatus(hc.Id())
	if status != nil {
		if attackCount, robCount := status.AttackAndRobCount(); attackCount != 0 || robCount != 0 {
			hc.Send(region.NewS2cSelfBeenAttackRobChangedMsg(i64.Int32(attackCount), i64.Int32(robCount)))
		}
	}

	guildId, _ := hc.LockGetGuildId()
	if guildId != 0 {
		// 0 就不要发了
		totalCount, _ := rs.guildBeenAttackOrRobMap.Get(guildId)
		if totalCount > 0 {
			hc.Send(region.NewS2cGuildBeenAttackRobChangedMsg(i64.Int32(totalCount)))
		}
	}
}

var noGuildBeenAttackRobChangedMsg = region.NewS2cGuildBeenAttackRobChangedMsg(0).Static()

func (rs *ReminderService) ChangeAttackOrRobCount(heroId, beenAttackCount, beenRobCount, newGuildId int64, isHome bool) {
	logrus.Debugf("ReminderService.ChangeAttackOrRobCount %d, %d, %d, %t", heroId, beenAttackCount, beenRobCount, isHome)

	var status *hero_been_attack_or_rob_status
	if beenAttackCount == 0 && beenRobCount == 0 && newGuildId == 0 {
		status = rs.getStatus(heroId)
		if status == nil {
			// 以前没有被打，没有联盟，现在还是没有被打没有联盟
			return
		}
	} else {
		status = rs.getOrCreateStatus(heroId)
	}

	oldTotal, curTotal, oldGuildId, changed := status.ChangeAttackOrRobCount(beenAttackCount, beenRobCount, newGuildId, isHome)
	logrus.Debugf("ReminderService.ChangeAttackOrRobCount %d, oldTotal: %d, curTotal: %d, guildId: %d, changed: %t info!", heroId, oldTotal, curTotal, oldGuildId, changed)
	if !changed {
		// 没有任何变化
		return
	}

	// 有变化
	rs.worldService.SendFunc(heroId, func() pbutil.Buffer {
		attackCount, robCount := status.AttackAndRobCount()
		return region.NewS2cSelfBeenAttackRobChangedMsg(i64.Int32(attackCount), i64.Int32(robCount))
	})

	if oldGuildId == newGuildId {
		// 同一个联盟
		if oldGuildId == 0 {
			return
		}

		cur := rs.guildBeenAttackOrRobMap.Upsert(newGuildId, curTotal-oldTotal, func(exist bool, valueInMap int64, newValue int64) int64 {
			return valueInMap + newValue
		})

		if cur == 0 {
			rs.guildBeenAttackOrRobMap.RemoveIfSame(newGuildId, cur)
		}

		logrus.Debugf("ReminderService.ChangeAttackOrRobCount %d, %d, %d, %t self broadcast!", heroId, beenAttackCount, beenRobCount, isHome)
		rs.broadcastGuildReminderChanged(newGuildId, cur)
	} else {
		// 不是同一个联盟
		if oldGuildId != 0 {
			cur := rs.guildBeenAttackOrRobMap.Upsert(oldGuildId, -oldTotal, func(exist bool, valueInMap int64, newValue int64) int64 {
				return valueInMap + newValue
			})

			if cur == 0 {
				rs.guildBeenAttackOrRobMap.RemoveIfSame(oldGuildId, cur)
			}

			logrus.Debugf("ReminderService.ChangeAttackOrRobCount %d, %d, %d, %t self old guild changed broadcast!", heroId, beenAttackCount, beenRobCount, isHome)
			rs.broadcastGuildReminderChanged(oldGuildId, cur)
		}

		if newGuildId != 0 {
			cur := rs.guildBeenAttackOrRobMap.Upsert(newGuildId, curTotal, func(exist bool, valueInMap int64, newValue int64) int64 {
				return valueInMap + newValue
			})

			if cur == 0 {
				rs.guildBeenAttackOrRobMap.RemoveIfSame(newGuildId, cur)
			}

			logrus.Debugf("ReminderService.ChangeAttackOrRobCount %d, %d, %d, %t self new guild changed broadcast!", heroId, beenAttackCount, beenRobCount, isHome)
			rs.broadcastGuildReminderChanged(newGuildId, cur)
		} else {
			logrus.Debugf("ReminderService.ChangeAttackOrRobCount %d, %d, %d, %t self not guild now!", heroId, beenAttackCount, beenRobCount, isHome)
			rs.worldService.Send(heroId, noGuildBeenAttackRobChangedMsg)
		}
	}
}

func (rs *ReminderService) getStatus(heroId int64) *hero_been_attack_or_rob_status {
	status, _ := rs.heroBeenAttackOrRobMap.Get(heroId)
	return status
}

func (rs *ReminderService) getOrCreateStatus(heroId int64) *hero_been_attack_or_rob_status {
	status := rs.getStatus(heroId)
	if status != nil {
		return status
	}

	status = &hero_been_attack_or_rob_status{}

	if old, setSuccess := rs.heroBeenAttackOrRobMap.SetIfAbsent(heroId, status); !setSuccess {
		// 设置失败，肯定是里面有东西
		status = old
	}

	return status
}

func (rs *ReminderService) broadcastGuildReminderChanged(guildId, newCount int64) {
	guildSnapshot := rs.guildSnapshotService.GetSnapshot(guildId)
	if guildSnapshot == nil {
		return
	}

	if len(guildSnapshot.UserMemberIds) <= 0 {
		return
	}

	if newCount == 0 {
		rs.worldService.MultiSend(guildSnapshot.UserMemberIds, noGuildBeenAttackRobChangedMsg)
	} else {
		rs.worldService.MultiSend(guildSnapshot.UserMemberIds, region.NewS2cGuildBeenAttackRobChangedMsg(i64.Int32(newCount)))
	}
	logrus.Debugf("ReminderService.broadcastGuildReminderChanged %d, %d broadcast!", guildId, newCount)
}

// 玩家被攻击、掠夺的状态
type hero_been_attack_or_rob_status struct {
	homeBeenAttackCount int64
	homeBeenRobCount    int64
	homeGuildId         int64 // 联盟id

	tentBeenAttackCount int64
	tentBeenRobCount    int64
	tentGuildId         int64 // 联盟id

	sync.Mutex // 锁
}

func (status *hero_been_attack_or_rob_status) ChangeAttackOrRobCount(beenAttackCount, beenRobCount, newGuildId int64, isHome bool) (oldTotal, curTotal, guildId int64, changed bool) {
	status.Lock()
	defer status.Unlock()

	if isHome {
		changed = status.homeBeenAttackCount != beenAttackCount || status.homeBeenRobCount != beenRobCount
		oldTotal = status.homeBeenAttackCount + status.homeBeenRobCount
		status.homeBeenAttackCount = beenAttackCount
		status.homeBeenRobCount = beenRobCount
		curTotal = status.homeBeenAttackCount + status.homeBeenRobCount

		guildId = status.homeGuildId
		status.homeGuildId = newGuildId
	} else {
		changed = status.tentBeenAttackCount != beenAttackCount || status.tentBeenRobCount != beenRobCount
		oldTotal = status.tentBeenAttackCount + status.tentBeenRobCount
		status.tentBeenAttackCount = beenAttackCount
		status.tentBeenRobCount = beenRobCount
		curTotal = status.tentBeenAttackCount + status.tentBeenRobCount

		guildId = status.tentGuildId
		status.tentGuildId = newGuildId
	}

	if !changed {
		changed = guildId != newGuildId
	}

	return
}

func (status *hero_been_attack_or_rob_status) ChangeGuild(newGuildId int64, isHome bool) (totalCount, oldGuildId int64) {
	status.Lock()
	defer status.Unlock()

	if isHome {
		oldGuildId = status.homeGuildId
		totalCount = status.homeBeenAttackCount + status.homeBeenRobCount
		status.homeGuildId = newGuildId
	} else {
		oldGuildId = status.tentGuildId
		totalCount = status.tentBeenAttackCount + status.tentBeenRobCount
		status.tentGuildId = newGuildId
	}

	return
}

func (status *hero_been_attack_or_rob_status) AttackAndRobCount() (attackCount, robCount int64) {
	// 不上锁了
	attackCount = status.homeBeenAttackCount + status.tentBeenAttackCount
	robCount = status.homeBeenRobCount + status.tentBeenRobCount
	return
}
