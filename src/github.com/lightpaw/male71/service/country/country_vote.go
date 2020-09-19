package country

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/country"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func (m *CountryService) IsOnChangeNameVote(countryId uint64) bool {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.IsOnChangeNameVote, 找不到国家：%v", countryId)
		return false
	}

	return c.InChangeNameVoteDuration(m.time.CurrentTime())
}

func (m *CountryService) AfterUpgradeTitle(heroId int64, countryId, newCount uint64) {
	if !m.IsOnChangeNameVote(countryId) {
		return
	}

	c := m.Country(countryId)

	var oldCount uint64
	var agree bool
	var id int32
	if m.heroData.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		id, oldCount, agree = hero.CountryMisc().ClearNameVote()
		hero.CountryMisc().NameVote(id, newCount, agree)
		result.Changed()
		result.Ok()
	}) {
		return
	}

	if oldCount == newCount {
		return
	}

	c.UpdateChangeNameVoteCount(id, int(newCount-oldCount), agree)
	m.world.Send(heroId, country.NewS2cHeroChangeNameVoteCountUpdateNoticeMsg(u64.Int32(newCount)))

	m.DetailMsgCacheDisable(countryId)
}

func (m *CountryService) AfterChangeNameVoteStart(country *entity.Country) {
	t := NewChangeNameTicker(country, m)
	m.changeNameTickers[country.Id()] = t
	go call.CatchPanic(func() {
		logrus.Debugf("创建国家:%v 改名ticker: id:%v name:%v to %v", country.Id(), country.VoteId(), country.Name(), country.VoteNewName())
		t.changeNameTick(m.time.CurrentTime(), country.VoteEndTime())
	}, "国家改名")
}

func NewChangeNameTicker(c *entity.Country, countrySrv *CountryService) *changeNameTicker {
	m := &changeNameTicker{}
	m.country = c
	m.countrySrv = countrySrv
	return m
}

type changeNameTicker struct {
	country    *entity.Country
	countrySrv *CountryService
}

func (c *changeNameTicker) changeNameTick(ctime, endTime time.Time) {
	select {
	case <-time.After(endTime.Sub(ctime)):
		c.onChangeNameVoteEnd()
	}
}

func (m *changeNameTicker) onChangeNameVoteEnd() {
	oldName := m.country.Name()
	ended, changed := m.country.OnChangeNameVoteEnd()
	logrus.Debugf("执行国家:%v 改名ticker: id:%v name:%v to %v", m.country.Id(), m.country.VoteId(), oldName, m.country.Name())
	if !ended {
		return
	}

	m.countrySrv.BroadcastCountry(country.NewS2cChangeNameSuccNoticeMsg(changed, u64.Int32(m.country.Id()), m.country.Name()), m.country.Id())

	if !changed {
		if d := m.countrySrv.datas.MailHelp().CountryChangeNameFail; d != nil {
			proto := d.NewTextMail(shared_proto.MailType_MailNormal)
			proto.Text = d.Text.New().WithHero(m.countrySrv.CountryFlagHeroName(m.country.King())).JsonString()

			m.countrySrv.WalkCountryOnlineHero(m.country.Id(), func(id int64, hc iface.HeroController) {
				m.countrySrv.mail.SendProtoMail(id, proto, m.countrySrv.time.CurrentTime())
			})
		}
		return
	}

	m.countrySrv.updateCountriesMsg()

	newName := m.country.Name()
	if d := m.countrySrv.datas.MailHelp().CountryChangeNameSucc; d != nil {
		proto := d.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Text = d.Text.New().WithHero(m.countrySrv.CountryFlagHeroName(m.country.King())).WithOldName(oldName).WithNewName(newName).JsonString()

		m.countrySrv.WalkCountryOnlineHero(m.country.Id(), func(id int64, hc iface.HeroController) {
			m.countrySrv.mail.SendProtoMail(id, proto, m.countrySrv.time.CurrentTime())
		})
	}

	if d := m.countrySrv.datas.BroadcastHelp().CountryChangeName; d != nil {
		m.countrySrv.broadcast.Broadcast(d.Text.New().WithOldName(oldName).WithNewName(newName).JsonString(), d.SendChat)
	}
}
