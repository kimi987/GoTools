package worker

import (
	"github.com/lightpaw/male7/config/season"
	"github.com/lightpaw/male7/gen/service"
	"github.com/lightpaw/male7/service/heromodule"
	"time"
)

// 时间不准, 1-2秒都有可能
func (m *MessageWorker) updatePerSeconds() {
	if m.user != nil {
		hc := m.user.GetHeroController()
		if hc != nil {
			ctime := service.TimeService.CurrentTime()

			//if !m.isRobot && !service.IndividualServerConfig.GetIgnoreHeartBeat() {
			//	if ctime.After(m.lastReceiveClientMsgTime.Add(90 * time.Second)) {
			//		if !timeutil.IsZero(m.lastReceiveClientMsgTime) {
			//			// 30秒没收到任何消息，踢掉
			//			m.Disconnect(misc.ErrDisconectReasonFailHeartBeat)
			//			return
			//		}
			//
			//		m.lastReceiveClientMsgTime = ctime
			//	}
			//}

			services := heromodule.GetService(service.ServiceDep, service.ConfigDatas, service.RealmService, service.GuildService, service.WorldService, service.MailModule, service.BuffService)
			heromodule.TryUpdatePerSeconds(hc, ctime, services)
		}
	}
}

func (m *MessageWorker) updatePerMinute() {
	if m.user != nil {
		hc := m.user.GetHeroController()
		if hc != nil {
			ctime := service.TimeService.CurrentTime()

			//if !m.isRobot && !service.IndividualServerConfig.GetIgnoreHeartBeat() {
			//	if ctime.After(m.lastReceiveClientMsgTime.Add(90 * time.Second)) {
			//		if !timeutil.IsZero(m.lastReceiveClientMsgTime) {
			//			// 30秒没收到任何消息，踢掉
			//			m.Disconnect(misc.ErrDisconectReasonFailHeartBeat)
			//			return
			//		}
			//
			//		m.lastReceiveClientMsgTime = ctime
			//	}
			//}

			services := heromodule.GetService(service.ServiceDep, service.ConfigDatas, service.RealmService, service.GuildService, service.WorldService, service.MailModule, service.BuffService)
			heromodule.TryUpdatePerMinute(hc, ctime, services)

			// 登陆日志（每60分钟记录一次）
			//if hc.TryNextWriteOnlineLogTime(ctime, time.Hour) {
			//	gamelogs.HeroOnlineLog(constants.PID, hc.Sid(), m.user.Id())
			//}
		}
	}
}

func (m *MessageWorker) tryResetDailyMc(resetTime time.Time) {
	if m.user != nil {
		hc := m.user.GetHeroController()
		if hc != nil {
			ctime := service.TimeService.CurrentTime()

			services := heromodule.GetService(service.ServiceDep, service.ConfigDatas, service.RealmService, service.GuildService, service.WorldService, service.MailModule, service.BuffService)
			heromodule.TryResetDailyMc(hc, ctime, resetTime, services)
		}
	}
}

func (m *MessageWorker) tryResetWeekly(resetTime time.Time) {
	if m.user != nil {
		hc := m.user.GetHeroController()
		if hc != nil {
			ctime := service.TimeService.CurrentTime()
			heromodule.TryResetWeekly(hc, ctime, resetTime)
		}
	}
}

func (m *MessageWorker) tryResetDailyZero(resetTime time.Time) {
	if m.user != nil {
		hc := m.user.GetHeroController()
		if hc != nil {
			ctime := service.TimeService.CurrentTime()

			services := heromodule.GetService(service.ServiceDep, service.ConfigDatas, service.RealmService, service.GuildService, service.WorldService, service.MailModule, service.BuffService)
			heromodule.TryResetDailyZero(hc, ctime, resetTime, services)
		}
	}
}

func (m *MessageWorker) tryResetDaily(resetTime time.Time) {
	if m.user != nil {
		hc := m.user.GetHeroController()
		if hc != nil {
			ctime := service.TimeService.CurrentTime()

			services := heromodule.GetService(service.ServiceDep, service.ConfigDatas, service.RealmService, service.GuildService, service.WorldService, service.MailModule, service.BuffService)
			heromodule.TryResetDaily(hc, ctime, resetTime, services)
		}
	}
}

func (m *MessageWorker) tryResetXuanyuan(resetTime time.Time) {
	if m.user != nil {
		hc := m.user.GetHeroController()
		if hc != nil {
			heromodule.TryResetXuanyuan(hc, resetTime)
		}
	}
}

func (m *MessageWorker) tryResetSeason(resetTime time.Time, seasonData *season.SeasonData) {
	if m.user != nil {
		hc := m.user.GetHeroController()
		if hc != nil {
			ctime := service.TimeService.CurrentTime()

			heromodule.TryResetSeason(hc, ctime, resetTime, seasonData)
		}
	}
}

//func (m *MessageWorker) tryResetRandomEvent(seasonData *season.SeasonData, endTime time.Time)  {
//	if m.user != nil {
//		hc := m.user.GetHeroController()
//		if hc != nil {
//			heromodule.TryResetRandomEvent(hc, seasonData, endTime)
//		}
//	}
//}
