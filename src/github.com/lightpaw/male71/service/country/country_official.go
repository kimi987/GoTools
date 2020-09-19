package country

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/pb/country"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
)

func (m *CountryService) HeroOfficial(countryId uint64, heroId int64) (t shared_proto.CountryOfficialType) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.CountryService, 找不到国家：%v", countryId)
		return
	}

	return c.HeroOfficial(heroId)
}

func (m *CountryService) OfficialAppoint(selfId, targetId int64, countryId uint64, official shared_proto.CountryOfficialType, pos int32) (succMsg, errMsg pbutil.Buffer) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.OfficialAppoint, 找不到国家：%v", countryId)
		errMsg = country.ERR_OFFICIAL_APPOINT_FAIL_INVALID_COUNTRY
		return
	}

	if c.OfficialFull(official) {
		errMsg = country.ERR_OFFICIAL_APPOINT_FAIL_OFFICIAL_FULL
		return
	}

	if selfOfficial := c.HeroOfficial(selfId); !c.IsSubOfficial(selfOfficial, official) {
		logrus.Debugf("国家任职, 没有权限。%v, %v", selfOfficial, official)
		errMsg = country.ERR_OFFICIAL_APPOINT_FAIL_DENY
		return
	}

	if t := c.HeroOfficial(targetId); t != shared_proto.CountryOfficialType_COT_NO_OFFICIAL {
		errMsg = country.ERR_OFFICIAL_APPOINT_FAIL_HERO_IS_OFFICIAL
		return
	}

	if succ := m.officialAppoint(c, targetId, official, m.time.CurrentTime(), pos); !succ {
		errMsg = country.ERR_OFFICIAL_APPOINT_FAIL_SERVER_ERR
		return
	}

	succMsg = country.OFFICIAL_APPOINT_S2C

	ctime := m.time.CurrentTime()
	if d := m.datas.MailHelp().CountryOfficialAppoint; d != nil {
		if offData := m.datas.CountryOfficialData().Must(int(official)); offData != nil {
			proto := d.NewTextMail(shared_proto.MailType_MailNormal)
			proto.Title = d.Title.New().WithOfficial(offData.Name).JsonString()
			proto.Text = d.Text.New().WithOfficial(offData.Name).WithHero(m.CountryFlagHeroName(selfId)).JsonString()
			m.mail.SendProtoMail(targetId, proto, ctime)
		}
	}

	return
}

func (m *CountryService) OfficialDepose(selfId, targetId int64, countryId uint64) (oldType shared_proto.CountryOfficialType, succMsg, errMsg pbutil.Buffer) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.OfficialDepose, 找不到国家：%v", countryId)
		errMsg = country.ERR_OFFICIAL_DEPOSE_FAIL_INVALID_COUNTRY
		return
	}

	targetOffType := c.HeroOfficial(targetId)
	if targetOffType == shared_proto.CountryOfficialType_COT_NO_OFFICIAL {
		errMsg = country.ERR_OFFICIAL_DEPOSE_FAIL_HERO_NOT_IN_POST
		return
	}

	if selfOffType := c.HeroOfficial(selfId); !c.IsSubOfficial(selfOffType, targetOffType) {
		logrus.Debugf("国家免职, 没有权限。%v, %v", selfOffType, targetOffType)
		errMsg = country.ERR_OFFICIAL_DEPOSE_FAIL_DENY
		return
	}

	if succ := m.officialDepose(c, targetId, m.time.CurrentTime()); !succ {
		errMsg = country.ERR_OFFICIAL_DEPOSE_FAIL_IN_CD
		return
	}

	succMsg = country.OFFICIAL_DEPOSE_S2C

	ctime := m.time.CurrentTime()
	if d := m.datas.MailHelp().CountryOfficialDepose; d != nil {
		if offData := m.datas.CountryOfficialData().Must(int(oldType)); offData != nil {
			proto := d.NewTextMail(shared_proto.MailType_MailNormal)
			proto.Title = d.Title.New().WithOfficial(offData.Name).JsonString()
			proto.Text = d.Text.New().WithOfficial(offData.Name).WithHero(m.CountryFlagHeroName(selfId)).JsonString()
			m.mail.SendProtoMail(targetId, proto, ctime)
		}
	}

	return
}

func (m *CountryService) OfficialLeave(heroId int64, countryId uint64) (oldType shared_proto.CountryOfficialType, succMsg, errMsg pbutil.Buffer) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.OfficialLeave, 找不到国家：%v", countryId)
		errMsg = country.ERR_OFFICIAL_LEAVE_FAIL_INVALID_COUNTRY
		return
	}

	var offType shared_proto.CountryOfficialType
	if currOffType := c.HeroOfficial(heroId); currOffType == shared_proto.CountryOfficialType_COT_NO_OFFICIAL {
		errMsg = country.ERR_OFFICIAL_LEAVE_FAIL_NO_OFFICIAL
		return
	} else {
		offType = currOffType
	}

	if offType == shared_proto.CountryOfficialType_COT_KING {
		errMsg = country.ERR_OFFICIAL_LEAVE_FAIL_IS_KING
		return
	}

	if succ := m.officialDepose(c, heroId, m.time.CurrentTime()); !succ {
		errMsg = country.ERR_OFFICIAL_LEAVE_FAIL_IN_CD
		return
	}

	succMsg = country.OFFICIAL_LEAVE_S2C

	return
}

func (m *CountryService) King(countryId uint64) (heroId int64) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.King, 找不到国家：%v", countryId)
		return
	}

	return c.King()
}

func (m *CountryService) ChangeCountryHost(countryId uint64, newKingId int64) (succ bool) {
	if succ = m.changeKing(countryId, newKingId); !succ {
		return
	}

	m.world.Broadcast(country.NewS2cCountryHostChangedNoticeMsg(u64.Int32(countryId)))

	if d := m.datas.BroadcastHelp().CountryChangeKing; d != nil {
		m.broadcast.Broadcast(d.Text.New().WithHero(m.CountryFlagHeroName(newKingId)).WithCountry(m.Country(countryId).Name()).JsonString(), d.SendChat)
	}

	return
}

func (m *CountryService) IsCountryDestroyed(countryId uint64) bool {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.IsCountryDestroyed, 找不到国家：%v", countryId)
		return false
	}
	return c.IsDestroyed()
}

func (m *CountryService) CancelCountryDestroy(countryId uint64) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.CancelCountryDestroy, 找不到国家：%v", countryId)
		return
	}

	c.CancelDestroy()
	m.MsgCacheDisable(countryId)
}

func (m *CountryService) CountryDestroy(countryId uint64) (succ bool) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.CountryDestroy, 找不到国家：%v", countryId)
		return
	}

	ctime := m.time.CurrentTime()
	m.officialDeposeAll(c, ctime)
	c.Destroy()
	m.world.Broadcast(country.NewS2cCountryDestroyNoticeMsg(u64.Int32(countryId)))
	succ = true

	m.MsgCacheDisable(countryId)

	return
}

func (m *CountryService) officialDeposeAll(c *entity.Country, ctime time.Time) {
	oldOfficials := c.OfficialDeposeAll(m.datas)
	m.MsgCacheDisable(c.Id())

	for t, heros := range oldOfficials {
		d := m.datas.GetCountryOfficialData(int(t))
		for _, hid := range heros {
			m.heroData.FuncWithSend(hid, func(hero *entity.Hero, result herolock.LockResult) {
				heromodule.CountryOfficialDepose(hero, result, d, m.datas, ctime)
			})
			if d.Buff != nil {
				m.buffSrv.CancelGroup(hid, d.Buff.Group)
			}
		}
	}
}

func (m *CountryService) ChangeKing(countryId uint64, newKingId int64) (succ bool) {
	if succ = m.changeKing(countryId, newKingId); !succ {
		return
	}

	m.world.Broadcast(country.NewS2cKingChangedNoticeMsg(u64.Int32(countryId)))

	return
}

func (m *CountryService) changeKing(countryId uint64, newKingId int64) (succ bool) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.ChangeKing, 找不到国家：%v", countryId)
		return
	}

	if exist, err := m.heroData.Exist(newKingId); !exist || err != nil {
		logrus.Debugf("CountryService.ChangeKing 新国君不是真实玩家。%v", newKingId)
		return
	}

	ctime := m.time.CurrentTime()
	m.officialDeposeByOfficial(c, shared_proto.CountryOfficialType_COT_KING, ctime)

	if succ = m.officialAppoint(c, newKingId, shared_proto.CountryOfficialType_COT_KING, ctime, 0); !succ {
		return
	}

	m.MsgCacheDisable(countryId)

	succ = true
	return
}

func (m *CountryService) AfterChangeCountry(heroId int64, oldCountryId, newCountryId uint64, voteId int32, voteCount int, agree bool) {
	c := m.Country(oldCountryId)
	if c == nil {
		logrus.Debugf("CountryService.AfterChangeCountry, 找不到国家：%v", oldCountryId)
		return
	}

	m.officialDepose(c, heroId, m.time.CurrentTime())
	c.UpdateChangeNameVoteCount(voteId, -voteCount, agree)

	if m.IsOnChangeNameVote(newCountryId) {
		m.world.Send(heroId, country.CHANGE_NAME_START_NOTICE_S2C)
	}
}

func (m *CountryService) officialAppoint(c *entity.Country, heroId int64, t shared_proto.CountryOfficialType, ctime time.Time, pos int32) (succ bool) {
	var heroCountry uint64
	if m.heroData.FuncNotError(heroId, func(hero *entity.Hero) (heroChanged bool) {
		heroCountry = hero.CountryId()
		return
	}) {
		return
	}
	if heroCountry != c.Id() {
		return
	}

	if _, _, succ = c.OfficialAppoint(t, heroId, ctime, pos); !succ {
		return
	}

	m.DetailMsgCacheDisable(c.Id())

	// 任命成功
	d := m.datas.GetCountryOfficialData(int(t))
	m.heroData.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		heromodule.CountryOfficialAppoint(hero, result, d, m.datas, ctime)
	})
	if d.Buff != nil {
		m.buffSrv.AddBuffToSelf(d.Buff, heroId)
	}

	return
}

func (m *CountryService) officialDepose(c *entity.Country, heroId int64, ctime time.Time) (succ bool) {
	return m.officialDepose0(c, heroId, ctime, false)
}

func (m *CountryService) officialDepose0(c *entity.Country, heroId int64, ctime time.Time, force bool) (succ bool) {
	var oldType shared_proto.CountryOfficialType
	if oldType, succ = c.OfficialDepose(heroId, m.time.CurrentTime(), force); !succ {
		return
	}

	m.DetailMsgCacheDisable(c.Id())

	// 免职成功
	d := m.datas.GetCountryOfficialData(int(oldType))
	m.heroData.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
		heromodule.CountryOfficialDepose(hero, result, d, m.datas, ctime)
	})
	if d.Buff != nil {
		m.buffSrv.CancelGroup(heroId, d.Buff.Group)
	}

	return
}

func (m *CountryService) officialDeposeByOfficial(c *entity.Country, t shared_proto.CountryOfficialType, ctime time.Time) (succ bool) {
	heroIds, _ := c.OfficialDeposeByOfficial(t)
	if len(heroIds) <= 0 {
		return
	}

	m.DetailMsgCacheDisable(c.Id())

	for _, heroId := range heroIds {
		// 免职成功
		d := m.datas.GetCountryOfficialData(int(t))
		m.heroData.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
			heromodule.CountryOfficialDepose(hero, result, d, m.datas, ctime)
		})
		if d.Buff != nil {
			m.buffSrv.CancelGroup(heroId, d.Buff.Group)
		}
	}

	succ = true
	return
}

func (m *CountryService) ForceOfficialAppoint(countryId uint64, heroId int64, official shared_proto.CountryOfficialType) (succ bool) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.ForceOfficialAppoint, 找不到国家：%v", countryId)
		return
	}

	m.officialAppoint(c, heroId, official, m.time.CurrentTime(), 0)
	succ = true
	return
}

func (m *CountryService) ForceOfficialDepose(countryId uint64, heroId int64) (succ bool) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.ForceOfficialDepose, 找不到国家：%v", countryId)
		return
	}

	m.officialDepose0(c, heroId, m.time.CurrentTime(), true)
	succ = true
	return
}

func (m *CountryService) GmOfficialDeposeAll(countryId uint64) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.GmOfficialDeposeAll, 找不到国家：%v", countryId)
		return
	}

	m.officialDeposeAll(c, m.time.CurrentTime())
}

func (m *CountryService) GmAppointKing(countryId uint64, heroId int64) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.GmAppointKing, 找不到国家：%v", countryId)
		return
	}

	m.officialAppoint(c, heroId, shared_proto.CountryOfficialType_COT_KING, m.time.CurrentTime(), 0)
	m.MsgCacheDisable(countryId)
}

func (m *CountryService) GmAppointOfficial(countryId uint64, heroId int64, t shared_proto.CountryOfficialType) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.GmAppointKing, 找不到国家：%v", countryId)
		return
	}

	m.officialAppoint(c, heroId, t, m.time.CurrentTime(), 0)
	m.MsgCacheDisable(countryId)
}

func (m *CountryService) GmDeposeKing(countryId uint64) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.GmDeposeKing, 找不到国家：%v", countryId)
		return
	}

	kingId := m.King(countryId)

	m.officialDepose0(c, kingId, m.time.CurrentTime(), true)
	m.MsgCacheDisable(countryId)
}

func (m *CountryService) GmDestroy(countryId uint64) {
	c := m.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryService.GmDestroy, 找不到国家：%v", countryId)
		return
	}

	kingId := m.King(countryId)

	m.officialDepose0(c, kingId, m.time.CurrentTime(), true)

	m.CountryDestroy(countryId)
}
