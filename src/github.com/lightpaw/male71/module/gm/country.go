package gm

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"strconv"
)

func (m *GmModule) newCountryGmGroup() *gm_group {
	return &gm_group{
		tab: "国家",
		handler: []*gm_handler{
			newStringHandler("所有人免职", "", m.countryOfficialDeposeAll),
			newStringHandler("任命国王", "", m.countryOfficialAppointKing),
			newStringHandler("任命官职", "2", m.countryOfficialAppointOfficial),
			newStringHandler("国王免职", "", m.countryOfficialDeposeKing),
			//newIntHandler("国家all", "0", m.countryAll),
			//newStringHandler("国家detail", "3", m.countryDetail),
			//newStringHandler("国家任职", "3 3 14797", m.countryAppoint),
			//newStringHandler("国家免职", "3 14797", m.countryDepose),
			newStringHandler("国家卸任", "", m.countryLeave),
			//newStringHandler("领取官职俸禄", "3", m.collect_official_salary),
			//newStringHandler("发起国家改名", "新", m.change_name_start),
			//newStringHandler("国家改名投票", "true", m.change_name_vote),
			//newStringHandler("国家任命玩家默认列表", " ", m.default_to_appoint_hero_list),
			//newStringHandler("国家任命玩家搜索列表", " ", m.search_to_appoint_hero_list),
		},
	}
}

func (m *GmModule) countryAll(cid int64, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleCountry_c2s_request_countries(amount string, hc iface.HeroController)
	}); ok {
		im.handleCountry_c2s_request_countries(strconv.FormatInt(0, 10), hc)
	}
}

func (m *GmModule) countryDetail(cid string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleCountry_c2s_country_detail(amount string, hc iface.HeroController)
	}); ok {
		im.handleCountry_c2s_country_detail(cid, hc)
	}
}

func (m *GmModule) countryAppoint(cid string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleCountry_c2s_official_appoint(amount string, hc iface.HeroController)
	}); ok {
		im.handleCountry_c2s_official_appoint(cid, hc)
	}
}

func (m *GmModule) countryDepose(cid string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleCountry_c2s_official_depose(amount string, hc iface.HeroController)
	}); ok {
		im.handleCountry_c2s_official_depose(cid, hc)
	}
}

func (m *GmModule) countryLeave(cid string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleCountry_c2s_official_leave(amount string, hc iface.HeroController)
	}); ok {
		im.handleCountry_c2s_official_leave(cid, hc)
	}
}

func (m *GmModule) collect_official_salary(cid string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleCountry_c2s_collect_official_salary(amount string, hc iface.HeroController)
	}); ok {
		im.handleCountry_c2s_collect_official_salary(cid, hc)
	}
}

func (m *GmModule) change_name_start(cid string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleCountry_c2s_change_name_start(amount string, hc iface.HeroController)
	}); ok {
		im.handleCountry_c2s_change_name_start(cid, hc)
	}
}

func (m *GmModule) change_name_vote(cid string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleCountry_c2s_change_name_vote(amount string, hc iface.HeroController)
	}); ok {
		im.handleCountry_c2s_change_name_vote(cid, hc)
	}
}

func (m *GmModule) default_to_appoint_hero_list(cid string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleCountry_c2s_default_to_appoint_hero_list(amount string, hc iface.HeroController)
	}); ok {
		im.handleCountry_c2s_default_to_appoint_hero_list(cid, hc)
	}
}

func (m *GmModule) search_to_appoint_hero_list(cid string, hc iface.HeroController) {
	var i interface{} = m
	if im, ok := i.(interface {
		handleCountry_c2s_search_to_appoint_hero_list(amount string, hc iface.HeroController)
	}); ok {
		im.handleCountry_c2s_search_to_appoint_hero_list(cid, hc)
	}
}

func (m *GmModule) countryOfficialAppointOfficial(official string, hc iface.HeroController) {
	m.country.GmAppointOfficial(hc.LockHeroCountry(), hc.Id(), shared_proto.CountryOfficialType(parseInt32(official)))
}

func (m *GmModule) countryOfficialAppointKing(str string, hc iface.HeroController) {
	m.country.GmAppointKing(hc.LockHeroCountry(), hc.Id())
}

func (m *GmModule) countryOfficialDeposeKing(str string, hc iface.HeroController) {
	m.country.GmDeposeKing(hc.LockHeroCountry())
}

func (m *GmModule) countryOfficialDeposeAll(str string, hc iface.HeroController) {
	m.country.GmOfficialDeposeAll(hc.LockHeroCountry())
}



