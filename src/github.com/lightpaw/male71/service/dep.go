package service

import (
	"github.com/lightpaw/male7/gen/iface"
)

func NewServiceDep(datas iface.ConfigDatas,
	svrConf iface.IndividualServerConfig,
	guild iface.GuildService,
	guildSnapshot iface.GuildSnapshotService,
	world iface.WorldService,
	broadcast iface.BroadcastService,
	heroSnapshot iface.HeroSnapshotService,
	heroData iface.HeroDataService,
	timeSrv iface.TimeService,
	push iface.PushService,
	chat iface.ChatService,
	db iface.DbService,
	country iface.CountryService,
	mingc iface.MingcService,
	fight iface.FightService,
	fightX         iface.FightXService,
	tlog iface.TlogService,
	mail iface.MailModule,
) *ServiceDep {
	m := &ServiceDep{}
	m.db = db
	m.svrConf = svrConf
	m.datas = datas
	m.guild = guild
	m.guildSnapshot = guildSnapshot
	m.world = world
	m.broadcast = broadcast
	m.timeSrv = timeSrv
	m.heroData = heroData
	m.heroSnapshot = heroSnapshot
	m.push = push
	m.chat = chat
	m.country = country
	m.mingc = mingc
	m.fight = fight
	m.fightX = fightX
	m.tlog = tlog
	m.mail = mail

	return m
}

//gogen:iface
type ServiceDep struct {
	db            iface.DbService
	datas         iface.ConfigDatas
	svrConf       iface.IndividualServerConfig
	timeSrv       iface.TimeService
	guild         iface.GuildService
	guildSnapshot iface.GuildSnapshotService
	world         iface.WorldService
	broadcast     iface.BroadcastService
	heroSnapshot  iface.HeroSnapshotService
	heroData      iface.HeroDataService
	push          iface.PushService
	chat          iface.ChatService
	country       iface.CountryService
	mingc         iface.MingcService
	fight         iface.FightService
	fightX         iface.FightXService
	tlog          iface.TlogService
	mail          iface.MailModule
}

func (m *ServiceDep) Mail() iface.MailModule {
	return m.mail
}

func (m *ServiceDep) Fight() iface.FightService {
	return m.fight
}

func (m *ServiceDep) FightX() iface.FightXService {
	return m.fightX
}

func (m *ServiceDep) SvrConf() iface.IndividualServerConfig {
	return m.svrConf
}

func (m *ServiceDep) Mingc() iface.MingcService {
	return m.mingc
}

func (m *ServiceDep) Country() iface.CountryService {
	return m.country
}

func (m *ServiceDep) Db() iface.DbService {
	return m.db
}

func (m *ServiceDep) Chat() iface.ChatService {
	return m.chat
}

func (m *ServiceDep) Push() iface.PushService {
	return m.push
}

func (m *ServiceDep) HeroSnapshot() iface.HeroSnapshotService {
	return m.heroSnapshot
}

func (m *ServiceDep) HeroData() iface.HeroDataService {
	return m.heroData
}

func (m *ServiceDep) Time() iface.TimeService {
	return m.timeSrv
}

func (m *ServiceDep) Datas() iface.ConfigDatas {
	return m.datas
}

func (m *ServiceDep) Guild() iface.GuildService {
	return m.guild
}

func (m *ServiceDep) GuildSnapshot() iface.GuildSnapshotService {
	return m.guildSnapshot
}

func (m *ServiceDep) World() iface.WorldService {
	return m.world
}

func (m *ServiceDep) Broadcast() iface.BroadcastService {
	return m.broadcast
}

func (m *ServiceDep) Tlog() iface.TlogService {
	return m.tlog
}
