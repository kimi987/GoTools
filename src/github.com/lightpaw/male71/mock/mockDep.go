package mock

import (
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/service/herosnapshot"
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"time"
)

var initDep = false

func MockDep() *ifacemock.MockServiceDep {

	if initDep {
		return ifacemock.ServiceDep
	}
	initDep=true

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	if err != nil {
		logrus.WithError(err).Panic("mockDep err!")
	}

	s := ifacemock.ServiceDep
	s.Mock(s.Datas, func() iface.ConfigDatas {
		return datas
	})

	s.Mock(s.Time, func() iface.TimeService {
		ts := ifacemock.TimeService
		//ts.Mock(ts.CurrentTime, func() time.Time {
		//	return time.Now()
		//})
		return ts
	})
	s.Mock(s.HeroSnapshot, func() iface.HeroSnapshotService {
		return ifacemock.HeroSnapshotService
	})
	s.Mock(s.HeroData, func() iface.HeroDataService {
		return ifacemock.HeroDataService
	})
	s.Mock(s.Guild, func() iface.GuildService {
		return ifacemock.GuildService
	})
	s.Mock(s.World, func() iface.WorldService {
		return ifacemock.WorldService
	})
	s.Mock(s.Broadcast, func() iface.BroadcastService {
		return ifacemock.BroadcastService
	})
	s.Mock(s.Db, func() iface.DbService {
		return ifacemock.DbService
	})
	s.Mock(s.Country, func() iface.CountryService {
		return ifacemock.CountryService
	})
	s.Mock(s.Mingc, func() iface.MingcService {
		return ifacemock.MingcService
	})
	s.Mock(s.GuildSnapshot, func() iface.GuildSnapshotService {
		return ifacemock.GuildSnapshotService
	})
	s.Mock(s.Push, func() iface.PushService {
		return ifacemock.PushService
	})
	s.Mock(s.Chat, func() iface.ChatService {
		return ifacemock.ChatService
	})
	s.Mock(s.Tlog, func() iface.TlogService {
		return ifacemock.TlogService
	})

	return s
}

func initMockDep2(s *MockDep2) {
	s.Mock(s.Time, func() iface.TimeService {
		return ifacemock.TimeService
	})
	//s.Mock(s.HeroSnapshot, func() iface.HeroSnapshotService {
	//	return ifacemock.HeroSnapshotService
	//})
	s.Mock(s.HeroData, func() iface.HeroDataService {
		return ifacemock.HeroDataService
	})
	s.Mock(s.Guild, func() iface.GuildService {
		return ifacemock.GuildService
	})
	s.Mock(s.World, func() iface.WorldService {
		return ifacemock.WorldService
	})
	s.Mock(s.Broadcast, func() iface.BroadcastService {
		return ifacemock.BroadcastService
	})
	s.Mock(s.Db, func() iface.DbService {
		return ifacemock.DbService
	})
	s.Mock(s.Country, func() iface.CountryService {
		return ifacemock.CountryService
	})
	s.Mock(s.Mingc, func() iface.MingcService {
		return ifacemock.MingcService
	})
	s.Mock(s.GuildSnapshot, func() iface.GuildSnapshotService {
		return ifacemock.GuildSnapshotService
	})
	s.Mock(s.Push, func() iface.PushService {
		return ifacemock.PushService
	})
	s.Mock(s.Chat, func() iface.ChatService {
		return ifacemock.ChatService
	})
	s.Mock(s.Tlog, func() iface.TlogService {
		return ifacemock.TlogService
	})
	s.Mock(s.Mail, func() iface.MailModule {
		return ifacemock.MailModule
	})
}

type MockDep2 struct {
	ifacemock.MockServiceDep
	datas        iface.ConfigDatas
	heroSnapshot iface.HeroSnapshotService
}

var mockDep2 *MockDep2

func NewMockDep2() *MockDep2 {
	if mockDep2 != nil {
		return mockDep2
	}

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	if err != nil {
		logrus.WithError(err).Panic("mockDep err!")
	}

	hss := herosnapshot.NewHeroSnapshotService(ifacemock.DbService, datas, ifacemock.GuildSnapshotService, ifacemock.BaiZhanService, &kv.IndividualServerConfig{})
	m := &MockDep2{datas: datas, heroSnapshot: hss}

	initMockDep2(m)

	mockDep2 = m
	return m
}

func (m *MockDep2) Datas() iface.ConfigDatas {
	return m.datas
}

func (m *MockDep2) HeroSnapshot() iface.HeroSnapshotService {
	return m.heroSnapshot
}

func MockTick() *ifacemock.MockTickerService {
	s := &ifacemock.MockTickerService{}

	s.Mock(s.GetDailyTickTime, func() tickdata.TickTime {
		return tickdata.New(time.Now(), 1*time.Hour)
	})

	s.Mock(s.GetPer10MinuteTickTime, s.GetDailyTickTime)
	s.Mock(s.GetPer30MinuteTickTime, s.GetDailyTickTime)
	s.Mock(s.GetPerHourTickTime, s.GetDailyTickTime)
	s.Mock(s.GetPerMinuteTickTime, s.GetDailyTickTime)

	return s
}
