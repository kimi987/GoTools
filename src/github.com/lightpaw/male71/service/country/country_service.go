package country

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/country"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/pbutil"
)

func NewCountryService(datas iface.ConfigDatas, db iface.DbService, ticker iface.TickerService,
	timeSrv iface.TimeService, heroSnapshot iface.HeroSnapshotService, heroData iface.HeroDataService,
	buffSrv iface.BuffService, world iface.WorldService, mail iface.MailModule, broadcast iface.BroadcastService) *CountryService {
	m := &CountryService{}
	m.db = db
	m.datas = datas
	m.ticker = ticker
	m.time = timeSrv
	m.heroData = heroData
	m.heroSnapshot = heroSnapshot
	m.buffSrv = buffSrv
	m.world = world
	m.mail = mail
	m.broadcast = broadcast
	m.msg = NewCountryMsg(timeSrv)
	m.changeNameTickers = make(map[uint64]*changeNameTicker)

	m.countries = make(map[uint64]*entity.Country)
	for _, c := range m.datas.GetCountryDataArray() {
		m.countries[c.Id] = entity.NewCountry(c, datas)
	}

	m.init()

	for _, c := range m.countries {
		if !c.IsChangeNameVoteCompleted() {
			m.AfterChangeNameVoteStart(c)
		}
	}

	m.stopFunc = ticker.TickPer10Minute("定时保存国家数据", func(tick tickdata.TickTime) {
		ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			m.save(ctx)
			return nil
		})
	})

	return m
}

//gogen:iface
type CountryService struct {
	db           iface.DbService
	datas        iface.ConfigDatas
	ticker       iface.TickerService
	time         iface.TimeService
	heroData     iface.HeroDataService
	heroSnapshot iface.HeroSnapshotService
	buffSrv      iface.BuffService
	world        iface.WorldService
	mail         iface.MailModule
	broadcast    iface.BroadcastService

	msg *CountryMsg

	stopFunc func()

	countries map[uint64]*entity.Country

	changeNameTickers map[uint64]*changeNameTicker
}

func (m *CountryService) init() {

	var bytes []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		bytes, err = m.db.LoadKey(ctx, server_proto.Key_Country)
		return
	})
	if err != nil {
		logrus.WithError(err).Panic("加载国家模块数据失败")
	}

	if len(bytes) > 0 {
		proto := &server_proto.CountriesServerProto{}
		if err := proto.Unmarshal(bytes); err != nil {
			logrus.WithError(err).Panic("解析国家模块数据失败")
		}

		for _, cp := range proto.Country {
			if c := m.countries[cp.Id]; c != nil {
				c.Unmarshal(cp, m.datas)
			}
		}
	}

	m.updateCountriesMsg()
}

func (m *CountryService) Close() {
	if m.stopFunc != nil {
		m.stopFunc()
	}

	// 下线保存一次
	ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		m.save(ctx)
		return nil
	})
}

func (m *CountryService) save(ctx context.Context) {
	if err := m.db.SaveKey(ctx, server_proto.Key_Country, must.Marshal(m.encodeServer())); err != nil {
		logrus.WithError(err).Panic("保存国家模块数据失败")
	}
}

func (m *CountryService) encode() *shared_proto.CountriesProto {
	proto := &shared_proto.CountriesProto{}
	for _, c := range m.countries {
		proto.Countries = append(proto.Countries, c.EncodeBasic(m.heroSnapshot))
	}

	return proto
}

func (m *CountryService) encodeServer() *server_proto.CountriesServerProto {
	proto := &server_proto.CountriesServerProto{}
	for _, c := range m.countries {
		proto.Country = append(proto.Country, c.EncodeServer())
	}

	return proto
}

func (m *CountryService) Country(id uint64) *entity.Country {
	return m.countries[id]
}

func (m *CountryService) Countries() (cs []*entity.Country) {
	for _, c := range m.countries {
		cs = append(cs, c)
	}
	return
}

func (m *CountryService) TutorialCountriesProto() *shared_proto.CountriesProto {
	return m.msg.TutorialCountriesProto()
}

func (m *CountryService) updateCountriesMsg() {
	m.msg.updateCountriesMsg(m.encode())

	m.world.Broadcast(m.msg.countriesNoticeMsg())
}

func (m *CountryService) OnHeroOnline(hc iface.HeroController, countryId uint64) {
	hc.Send(m.msg.countriesNoticeMsg())
}

func (m *CountryService) UpdateMcWarMsg(p *shared_proto.McWarProto) {
	m.msg.updateMcWarMsg(p)
}

func (m *CountryService) UpdateMingcsMsg(p *shared_proto.MingcsProto) {
	m.msg.updateMingcsMsg(p)
}

func (m *CountryService) CountryPrestigeMsg(ver uint64) pbutil.Buffer {
	return m.msg.countryPrestigeMsg(ver)
}

func (m *CountryService) CountriesMsg(ver uint64) pbutil.Buffer {
	return m.msg.getCountriesMsg(ver)
}

func (m *CountryService) CountryDetailMsg(countryId uint64) (msg pbutil.Buffer) {
	c := m.Country(countryId)
	if c == nil {
		msg = country.ERR_COUNTRY_DETAIL_FAIL_INVALID_ID
		return
	}

	msg = m.msg.countryDetailMsgCache.GetOrUpdate(countryId, func() (result pbutil.Buffer) {
		return country.NewS2cCountryDetailMsg(c.Encode(m.heroSnapshot, m.datas.CountryMiscData().ChangeNameCd))
	})
	return
}

func (m *CountryService) AddPrestige(id uint64, toAdd uint64) (newPrestige uint64, ok bool) {
	if c, ok := m.countries[id]; ok {
		newPrestige = c.AddPrestige(toAdd)
	}

	m.updateCountriesMsg()
	m.DetailMsgCacheDisable(id)
	return
}

func (m *CountryService) ReducePrestige(id uint64, toReduce uint64) (newPrestige uint64, ok bool) {
	if c, ok := m.countries[id]; ok {
		newPrestige = c.AddPrestige(toReduce)
	}

	m.updateCountriesMsg()
	m.DetailMsgCacheDisable(id)
	return
}

func (m *CountryService) HeroCountry(heroId int64) (countryId uint64) {
	hero := m.heroSnapshot.Get(heroId)
	if hero == nil {
		return
	}
	return hero.CountryId
}

func (m *CountryService) BroadcastCountry(msg pbutil.Buffer, countryId uint64) {
	msg = msg.Static()
	m.world.WalkHero(func(id int64, hc iface.HeroController) {
		if hc.LockHeroCountry() == countryId {
			hc.Send(msg)
		}
	})
}

func (m *CountryService) WalkCountryOnlineHero(countryId uint64, walkFunc iface.CountryHeroWalker) {
	m.world.WalkHero(func(id int64, hc iface.HeroController) {
		if m.HeroCountry(hc.Id()) == countryId {
			walkFunc(id, hc)
		}
	})
}

func (m *CountryService) MsgCacheDisable(countryId uint64) {
	m.updateCountriesMsg()
	m.DetailMsgCacheDisable(countryId)
}

func (m *CountryService) DetailMsgCacheDisable(countryId uint64) {
	m.msg.countryDetailMsgCache.Disable(countryId)
}

func (m *CountryService) LockHeroCapital(heroId int64) (mcId uint64) {
	countryId := m.LockHeroCountry(heroId)
	d := m.datas.GetCountryData(countryId)
	if d == nil {
		return
	}
	return d.Capital.Id
}

func (m *CountryService) LockHeroCountry(heroId int64) (countryId uint64) {
	m.heroData.FuncNotError(heroId, func(hero *entity.Hero) (heroChanged bool) {
		countryId = hero.CountryId()
		return
	})
	return
}

func (m *CountryService) CountryName(countryId uint64) (name string) {
	if c := m.Country(countryId); c != nil {
		return c.Name()
	}
	return
}

func (m *CountryService) CountryFlagHeroName(id int64) string {
	hero := m.heroSnapshot.Get(id)
	if hero == nil {
		return idbytes.PlayerName(id)
	}

	return m.datas.MiscConfig().CountryFlagHeroName.FormatIgnoreEmpty(m.CountryName(hero.CountryId), hero.GuildFlagName(), hero.Name)
}
