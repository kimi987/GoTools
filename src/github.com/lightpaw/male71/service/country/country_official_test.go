package country

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/service/mingc"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestCountryService(t *testing.T) {
	RegisterTestingT(t)
	dep := mock.NewMockDep2()
	dep.Mock(dep.Time, func() iface.TimeService {
		ts := ifacemock.TimeService
		ts.Mock(ts.CurrentTime, func() time.Time {
			return time.Now()
		})
		return ts
	})
	dep.Mock(dep.Country, func() iface.CountryService {
		return NewCountryService(dep.Datas(), dep.Db(), ifacemock.TickerService, dep.Time(), dep.HeroSnapshot(), ifacemock.BuffService, dep.World(), dep.Mail(), dep.Broadcast())
	})
	dep.Mock(dep.Mingc, func() iface.MingcService {
		return mingc.NewMingcService(dep.Datas(), dep.Db(), ifacemock.TickerService, dep.Guild(), dep.GuildSnapshot(), dep.Country(), dep.Time(), dep.SvrConf(), dep.Mail(), dep.HeroData(), dep.Broadcast(), dep.World(), dep.Tlog())
	})


}
