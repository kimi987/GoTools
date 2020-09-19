package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func NewHeroCountryMisc() *HeroCountryMisc {
	h := &HeroCountryMisc{}
	h.vote = NewHeroCountryChangeNameVote()
	return h
}

type HeroCountryMisc struct {
	showOfficialType shared_proto.CountryOfficialType // 只做展示用
	showAppointTime  time.Time
	effect           *sub.BuildingEffectData

	vote *heroCountryChangeNameVote
}

func (h *HeroCountryMisc) ShowOfficialType() shared_proto.CountryOfficialType {
	return h.showOfficialType
}

func (h *HeroCountryMisc) Appoint(d *country.CountryOfficialData, ctime time.Time) {
	h.showOfficialType = d.OfficialType
	h.effect = d.BuildingEffect
	h.showAppointTime = ctime
}

func (h *HeroCountryMisc) Depose() {
	h.showOfficialType = shared_proto.CountryOfficialType_COT_NO_OFFICIAL
	h.showAppointTime = time.Time{}
	h.effect = nil
}

func (h *HeroCountryMisc) encode() *shared_proto.HeroCountryMiscProto {
	p := &shared_proto.HeroCountryMiscProto{}
	p.OfficialType = h.showOfficialType
	p.AppointTime = timeutil.Marshal32(h.showAppointTime)
	p.Vote = h.vote.encode()
	return p
}

func (h *HeroCountryMisc) encodeServer() *server_proto.HeroCountryMiscServerProto {
	p := &server_proto.HeroCountryMiscServerProto{}
	p.OfficialType = h.showOfficialType
	p.AppointTime = timeutil.Marshal64(h.showAppointTime)
	p.Vote = h.vote.encodeServer()
	return p
}

func (h *HeroCountryMisc) unmarshal(p *server_proto.HeroCountryMiscServerProto, conf interface {
	GetCountryOfficialData(int) *country.CountryOfficialData
}) {
	if p == nil {
		return
	}

	if p.OfficialType == shared_proto.CountryOfficialType_COT_NO_OFFICIAL {
		return
	}

	d := conf.GetCountryOfficialData(int(p.OfficialType))
	if d == nil {
		logrus.Warnf("HeroCountryMisc.unmarshal 国家官职:%v 不在配置中。", p.OfficialType)
		return
	}

	h.showOfficialType = p.OfficialType
	h.showAppointTime = timeutil.Unix64(p.AppointTime)
	h.effect = d.BuildingEffect

	h.vote.unmarshal(p.Vote)
}

func NewHeroCountryChangeNameVote() *heroCountryChangeNameVote {
	h := &heroCountryChangeNameVote{}
	return h
}

type heroCountryChangeNameVote struct {
	id    int32
	votes uint64
	agree bool
}

func (h *heroCountryChangeNameVote) encode() *shared_proto.HeroCountryChangeNameVoteProto {
	p := &shared_proto.HeroCountryChangeNameVoteProto{}
	p.Id = h.id
	p.Votes = u64.Int32(h.votes)
	p.Agree = h.agree
	return p
}

func (h *heroCountryChangeNameVote) encodeServer() *server_proto.HeroCountryChangeNameVoteServerProto {
	p := &server_proto.HeroCountryChangeNameVoteServerProto{}
	p.Id = h.id
	p.Votes = h.votes
	p.Agree = h.agree
	return p
}

func (h *heroCountryChangeNameVote) unmarshal(p *server_proto.HeroCountryChangeNameVoteServerProto) {
	if p == nil {
		return
	}

	h.id = p.Id
	h.votes = p.Votes
	h.agree = p.Agree
}

func (h *HeroCountryMisc) NameVote(id int32, votes uint64, agree bool) {
	h.vote.id = id
	h.vote.votes = votes
	h.vote.agree = agree
}

func (h *HeroCountryMisc) ClearNameVote() (id int32, votes uint64, agree bool) {
	id = h.vote.id
	votes = h.vote.votes
	agree = h.vote.agree
	h.vote = NewHeroCountryChangeNameVote()
	return
}

func (h *HeroCountryMisc) IsNameVoted(id int32, agree bool) bool {
	return h.vote.id == id && h.vote.votes > 0 && agree == h.vote.agree
}

func (h *HeroCountryMisc) GetNameVote(id int32) (votes uint64, agree bool) {
	if id != h.vote.id {
		return
	}
	votes = h.vote.votes
	agree = h.vote.agree
	return
}
