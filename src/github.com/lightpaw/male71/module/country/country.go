package country

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/country"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"strings"
	"time"
	"unicode/utf8"
)

func NewCountryModule(dep iface.ServiceDep, buffSrv iface.BuffService, rank iface.RankModule) *CountryModule {
	m := &CountryModule{}
	m.dep = dep
	m.buffSrv = buffSrv
	m.countrySrv = dep.Country()
	m.rank = rank

	m.toAppointDefaultHerosMsgCache = msg.NewMsgCache(60*time.Second, m.dep.Time())

	return m
}

//gogen:iface
type CountryModule struct {
	dep        iface.ServiceDep
	countrySrv iface.CountryService
	buffSrv    iface.BuffService
	rank       iface.RankModule

	toAppointDefaultHerosMsgCache *msg.MsgCache
}

//gogen:iface
func (m *CountryModule) ProcessRequestCountryPrestige(proto *country.C2SRequestCountryPrestigeProto, hc iface.HeroController) {
	ver := proto.Vsn
	msg := m.countrySrv.CountryPrestigeMsg(u64.FromInt32(ver))
	hc.Send(msg)
}

//gogen:iface
func (m *CountryModule) ProcessRequestCountries(proto *country.C2SRequestCountriesProto, hc iface.HeroController) {
	ver := proto.Vsn
	msg := m.countrySrv.CountriesMsg(u64.FromInt32(ver))
	hc.Send(msg)
}

// 转国
//gogen:iface
func (m *CountryModule) ProcessHeroChangeCountry(proto *country.C2SHeroChangeCountryProto, hc iface.HeroController) {

	newCountry := u64.FromInt32(proto.NewCountry)
	if m.dep.Country().Country(newCountry) == nil {
		hc.Send(country.ERR_HERO_CHANGE_COUNTRY_FAIL_INVALID_COUNTRY)
		return
	}

	oldCountry := hc.LockHeroCountry()
	if m.countrySrv.HeroOfficial(oldCountry, hc.Id()) == shared_proto.CountryOfficialType_COT_KING {
		hc.Send(country.ERR_HERO_CHANGE_COUNTRY_FAIL_IS_KING)
		return
	}

	ctime := m.dep.Time().CurrentTime()

	var errMsg msg.ErrMsg
	var voteCount uint64
	var agree bool
	var voteId int32
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		errMsg = heromodule.ChangeCountry(m.dep, hero, result, ctime, newCountry, proto.Buy)
		voteId, voteCount, agree = hero.CountryMisc().ClearNameVote()
		result.Changed()
		result.Ok()
	})
	if errMsg != nil {
		hc.Send(errMsg.ErrMsg())
		return
	}

	m.countrySrv.AfterChangeCountry(hc.Id(), oldCountry, newCountry, voteId, u64.Int(voteCount), agree)
}

//gogen:iface
func (m *CountryModule) ProcessCountryDetail(proto *country.C2SCountryDetailProto, hc iface.HeroController) {
	hc.Send(m.countrySrv.CountryDetailMsg(u64.FromInt32(proto.CountryId)))
}

//gogen:iface
func (m *CountryModule) ProcessOfficialAppoint(proto *country.C2SOfficialAppointProto, hc iface.HeroController) {
	if proto.OfficialType == shared_proto.CountryOfficialType_COT_NO_OFFICIAL {
		hc.Send(country.ERR_OFFICIAL_APPOINT_FAIL_INVALID_OFFICIAL)
		return
	}

	targetId, ok := idbytes.ToId(proto.HeroId)
	if !ok {
		hc.Send(country.ERR_OFFICIAL_APPOINT_FAIL_INVALID_HERO)
		return
	}

	var targetGuildId int64
	var targetCountryId uint64
	if m.dep.HeroData().FuncNotError(targetId, func(hero *entity.Hero) (heroChanged bool) {
		targetGuildId = hero.GuildId()
		targetCountryId = hero.CountryId()
		return
	}) {
		hc.Send(country.ERR_OFFICIAL_APPOINT_FAIL_INVALID_HERO)
		return
	}

	var inGuildChangeCountry bool

	ctime := m.dep.Time().CurrentTime()
	if targetGuildId > 0 {
		m.dep.Guild().FuncGuild(targetGuildId, func(g *sharedguilddata.Guild) {
			if g != nil && ctime.Unix() < g.GetChangeCountryWaitEndTime() {
				inGuildChangeCountry = true
			}
		})
	}
	if inGuildChangeCountry {
		hc.Send(country.ERR_OFFICIAL_APPOINT_FAIL_HERO_IN_GUILD_CHANGE_COUNTRY)
		return
	}

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.CountryId() != targetCountryId {
			result.Add(country.ERR_OFFICIAL_APPOINT_FAIL_HERO_NOT_SAME_COUNTRY)
			return
		}
		result.Ok()
	}) {
		return
	}

	succMsg, errMsg := m.countrySrv.OfficialAppoint(hc.Id(), targetId, targetCountryId, proto.OfficialType, proto.Pos)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface
func (m *CountryModule) ProcessOfficialDepose(proto *country.C2SOfficialDeposeProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.HeroId)
	if !ok {
		hc.Send(country.ERR_OFFICIAL_DEPOSE_FAIL_INVALID_HERO)
		return
	}

	target := m.dep.HeroSnapshot().Get(targetId)
	if target == nil {
		hc.Send(country.ERR_OFFICIAL_DEPOSE_FAIL_INVALID_HERO)
		return
	}

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.CountryId() != target.CountryId {
			result.Add(country.ERR_OFFICIAL_DEPOSE_FAIL_HERO_NOT_SAME_COUNTRY)
			return
		}
		result.Ok()
	}) {
		return
	}

	_, succMsg, errMsg := m.countrySrv.OfficialDepose(hc.Id(), targetId, target.CountryId)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface c2s_official_leave
func (m *CountryModule) ProcessOfficialLeave(hc iface.HeroController) {
	countryId := hc.LockHeroCountry()
	if m.countrySrv.Country(countryId) == nil {
		hc.Send(country.ERR_OFFICIAL_LEAVE_FAIL_INVALID_COUNTRY)
		return
	}

	_, succMsg, errMsg := m.countrySrv.OfficialLeave(hc.Id(), countryId)
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(succMsg)
}

//gogen:iface c2s_collect_official_salary
func (m *CountryModule) ProcessCollectOfficialSalary(hc iface.HeroController) {
	countryId := hc.LockHeroCountry()
	c := m.countrySrv.Country(countryId)
	if c == nil {
		hc.Send(country.ERR_COLLECT_OFFICIAL_SALARY_FAIL_NO_COUNTRY)
		return
	}

	offType := m.countrySrv.HeroOfficial(countryId, hc.Id())
	if offType == shared_proto.CountryOfficialType_COT_NO_OFFICIAL {
		hc.Send(country.ERR_COLLECT_OFFICIAL_SALARY_FAIL_NO_OFFICIAL)
		return
	}

	d := m.dep.Datas().GetCountryOfficialData(int(offType))
	if d == nil {
		logrus.Debugf("领取国家官职俸禄，找不到官职 data. type:%v", offType)
		hc.Send(country.ERR_COLLECT_OFFICIAL_SALARY_FAIL_NO_OFFICIAL)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Country().GetAppointOnSameDay() {
			result.Add(country.ERR_COLLECT_OFFICIAL_SALARY_FAIL_APPOINT_ON_SAME_DAY)
			return
		}
		if hero.Country().GetDailySalaryCollected() {
			result.Add(country.ERR_COLLECT_OFFICIAL_SALARY_FAIL_ALREADY_COLLECTED)
			return
		}

		prize := d.Salary.Try()
		hero.Country().SetDailySalaryCollected(true)

		hctx := heromodule.NewContext(m.dep, operate_type.CountryCollectOfficialSalary)
		heromodule.AddPrize(hctx, hero, result, prize, m.dep.Time().CurrentTime())

		result.Add(country.NewS2cCollectOfficialSalaryMsg(prize.Encode()))

		result.Changed()
		result.Ok()
	})

}

//gogen:iface
func (m *CountryModule) ProcessChangeNameStart(proto *country.C2SChangeNameStartProto, hc iface.HeroController) {
	newName := strings.TrimSpace(proto.NewName)
	if utf8.RuneCountInString(newName) != 1 {
		hc.Send(country.ERR_CHANGE_NAME_START_FAIL_NEW_NAME_LEN_LIMIT)
		return
	}

	countryId := hc.LockHeroCountry()

	for _, d := range m.dep.Datas().GetCountryDataArray() {
		if d.Id != countryId && newName == d.Name {
			hc.Send(country.ERR_CHANGE_NAME_START_FAIL_NEW_NAME_IN_DEFAULT_NAME)
			return
		}
	}

	for _, c := range m.dep.Country().Countries() {
		if newName == c.Name() {
			hc.Send(country.ERR_CHANGE_NAME_START_FAIL_NEW_NAME_IN_CURRENT_NAME)
			return
		}
	}

	c := m.countrySrv.Country(countryId)
	if c == nil || c.HeroOfficial(hc.Id()) != shared_proto.CountryOfficialType_COT_KING {
		hc.Send(country.ERR_CHANGE_NAME_START_FAIL_NOT_KING)
		return
	}

	ctime := m.dep.Time().CurrentTime()
	if ctime.Before(c.VoteEndTime().Add(m.dep.Datas().CountryMiscData().ChangeNameCd)) {
		hc.Send(country.ERR_CHANGE_NAME_START_FAIL_IN_CD)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.CountryChangeName)

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !heromodule.TryReduceCost(hctx, hero, result, m.dep.Datas().CountryMiscData().ChangeNameCost) {
			result.Add(country.ERR_CHANGE_NAME_START_FAIL_COST_NOT_ENOUGH)
			return
		}
		result.Changed()
		result.Ok()
	}) {
		return
	}

	endTime := ctime.Add(m.dep.Datas().CountryMiscData().ChangeNameVoteDuration)
	c.StartChangeNameVote(newName, endTime)
	m.countrySrv.AfterChangeNameVoteStart(c)

	hc.Send(country.CHANGE_NAME_START_S2C)

	m.countrySrv.BroadcastCountry(country.CHANGE_NAME_START_NOTICE_S2C, countryId)
	m.countrySrv.DetailMsgCacheDisable(countryId)
}

//gogen:iface
func (m *CountryModule) ProcessChangeNameVote(proto *country.C2SChangeNameVoteProto, hc iface.HeroController) {
	countryId := hc.LockHeroCountry()
	c := m.countrySrv.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryModule.ProcessChangeNameVote,hero:%v 找不到国家:%v", hc.Id(), countryId)
		hc.Send(country.ERR_CHANGE_NAME_VOTE_FAIL_NOT_IN_VOTE)
		return
	}

	if !c.InChangeNameVoteDuration(m.dep.Time().CurrentTime()) {
		hc.Send(country.ERR_CHANGE_NAME_VOTE_FAIL_NOT_IN_VOTE)
		return
	}

	var oldAgree bool
	var oldVoteCount uint64

	var voteCount uint64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.CountryMisc().IsNameVoted(c.VoteId(), proto.Agree) {
			result.Add(country.ERR_CHANGE_NAME_VOTE_FAIL_VOTED)
			return
		}

		voteCount = heromodule.GetCountryChangeNameVoteCount(hero)
		if voteCount <= 0 {
			result.Add(country.ERR_CHANGE_NAME_VOTE_FAIL_NO_VOTE_COUNT)
			return
		}

		// 之前投过的票
		oldVoteCount, oldAgree = hero.CountryMisc().GetNameVote(c.VoteId())
		hero.CountryMisc().NameVote(c.VoteId(), voteCount, proto.Agree)

		result.Add(country.NewS2cChangeNameVoteMsg(proto.Agree, u64.Int32(voteCount), u64.Int32(oldVoteCount)))
		result.Changed()
		result.Ok()
	}) {
		return
	}

	if oldVoteCount > 0 {
		c.UpdateChangeNameVoteCount(c.VoteId(), -u64.Int(oldVoteCount), oldAgree)
	}
	c.UpdateChangeNameVoteCount(c.VoteId(), u64.Int(voteCount), proto.Agree)
	m.countrySrv.DetailMsgCacheDisable(countryId)
}

//gogen:iface c2s_default_to_appoint_hero_list
func (m *CountryModule) ProcessDefaultToAppointHeroList(hc iface.HeroController) {
	countryId := hc.LockHeroCountry()
	c := m.countrySrv.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryModule.ProcessDefaultToAppointHeroList, hero:%v 找不到国家:%v", hc.Id(), countryId)
		hc.Send(country.ERR_DEFAULT_TO_APPOINT_HERO_LIST_FAIL_SERVER_ERR)
		return
	}

	msg := m.toAppointDefaultHerosMsgCache.GetOrUpdate(countryId, func() (result pbutil.Buffer) {
		maxCount := m.dep.Datas().CountryMiscData().MaxSearchHeroDefaultCount

		heros := m.rank.CountryOfficial(maxCount, "", countryId, shared_proto.CountryOfficialType_COT_NO_OFFICIAL)
		return country.NewS2cDefaultToAppointHeroListMsg(heros)
	})

	hc.Send(msg)
}

//gogen:iface
func (m *CountryModule) ProcessSearchToAppointHeroList(proto *country.C2SSearchToAppointHeroListProto, hc iface.HeroController) {
	countryId := hc.LockHeroCountry()
	c := m.countrySrv.Country(countryId)
	if c == nil {
		logrus.Debugf("CountryModule.ProcessSearchToAppointHeroList, hero:%v 找不到国家:%v", hc.Id(), countryId)
		hc.Send(country.ERR_SEARCH_TO_APPOINT_HERO_LIST_FAIL_SERVER_ERR)
		return
	}

	name := strings.TrimSpace(proto.Name)
	if name == "" {
		hc.Send(country.ERR_SEARCH_TO_APPOINT_HERO_LIST_FAIL_NAME_IS_EMPTY)
		return
	}

	maxCount := m.dep.Datas().CountryMiscData().MaxSearchHeroByNameCount

	heros := m.rank.CountryOfficial(maxCount, name, countryId, shared_proto.CountryOfficialType_COT_NO_OFFICIAL)
	if len(heros) >= maxCount {
		hc.Send(country.NewS2cSearchToAppointHeroListMsg(heros))
		return
	}

	// 从 db 中补全
	var dbHeros []*entity.Hero
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		dbHeros, err = m.dep.Db().LoadHeroListByNameAndCountry(ctx, name, countryId, 0, u64.FromInt(maxCount))
		return
	})
	if err != nil {
		logrus.WithError(err).Debugf("国家官职任命，搜索玩家错误")
		hc.Send(guild.ERR_SEARCH_NO_GUILD_HEROS_FAIL_SERVER_ERROR)
		return
	}

	currHeroIds := make(map[int64]struct{})
	for _, h := range heros {
		if id, ok := idbytes.ToId(h.Basic.Id); ok {
			currHeroIds[id] = struct{}{}
		}
	}

	for _, h := range dbHeros {
		if _, ok := currHeroIds[h.Id()]; ok {
			continue
		}
		if m.countrySrv.HeroOfficial(countryId, h.Id()) != shared_proto.CountryOfficialType_COT_NO_OFFICIAL {
			continue
		}

		heroSnapshot := m.dep.HeroSnapshot().Get(h.Id())
		heros = append(heros, heroSnapshot.EncodeClient())
	}

	hc.Send(country.NewS2cSearchToAppointHeroListMsg(heros))
}
