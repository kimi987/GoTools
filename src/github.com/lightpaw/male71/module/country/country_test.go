package country

import (
	country2 "github.com/lightpaw/male7/service/country"
	"testing"
	. "github.com/onsi/gomega"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/pb/country"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/timeutil"
)

func TestCountryModule_ProcessHeroChangeCountry(t *testing.T) {
	RegisterTestingT(t)

	dep := mock.NewMockDep2()
	dep.Mock(dep.Country, func() iface.CountryService {
		return country2.NewCountryService(dep.Datas(), dep.Db(), ifacemock.TickerService, dep.Time(), dep.HeroSnapshot())
	})

	ctime := dep.Time().CurrentTime()
	hero := entity.NewHero(1, "hero1", dep.Datas().HeroInitData(), ctime)
	mock.DefaultHero(hero)

	m := NewCountryModule(dep, ifacemock.BuffService)
	c2s := &country.C2SHeroChangeCountryProto{NewCountry: 3}
	m.ProcessHeroChangeCountry(c2s, ifacemock.HeroController)

	Ω(hero.Country().GetCountryId()).Should(Equal(uint64(3)))

	msgs := mock.ReadMsgList(mock.LockResult)
	Ω(msgs[0]).Should(Equal(country.NewS2cHeroChangeCountryMsg(c2s.NewCountry,
		timeutil.Marshal32(hero.Country().GetNewUserExpiredTime()),
		timeutil.Marshal32(hero.Country().GetNormalExpiredTime()))))

	mock.LockResult.Reset()
}
