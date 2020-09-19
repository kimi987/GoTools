package seasonservice

import (
	"github.com/lightpaw/male7/config/season"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/ticker"
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/pbutil"
	"sync"
	"time"
)

//gogen:iface
type SeasonService struct {
	timeService    iface.TimeService
	serverConfig   iface.IndividualServerConfig
	seasonMiscData *season.SeasonMiscData
	ticker         *ticker.Ticker
	tickService    iface.TickerService
	datas          iface.ConfigDatas
	once           sync.Once
	cacheTime      time.Time
	cacheMsg       pbutil.Buffer
}

func NewSeasonService(timeService iface.TimeService,
	serverConfig iface.IndividualServerConfig,
	tickService iface.TickerService,
	datas iface.ConfigDatas) *SeasonService {

	ctime := timeService.CurrentTime()
	seasonDuration := datas.SeasonMiscData().SeasonDuration
	delay := seasonDuration - ctime.Sub(serverConfig.GetServerStartTime().Add(datas.MiscConfig().DailyResetDuration))%seasonDuration

	s := &SeasonService{
		timeService:    timeService,
		serverConfig:   serverConfig,
		tickService:    tickService,
		seasonMiscData: datas.SeasonMiscData(),
		ticker:         ticker.NewTicker(ctime, delay, seasonDuration),
		datas:          datas,
	}

	heromodule.RegisterHeroOnlineListener(s)

	return s
}

func (s *SeasonService) Season() *season.SeasonData {
	return s.SeasonByTime(s.timeService.CurrentTime())
}

func (s *SeasonService) SeasonByTime(ctime time.Time) *season.SeasonData {
	serverStartDuration := ctime.Sub(s.serverConfig.GetServerStartTime().Add(s.datas.MiscConfig().DailyResetDuration))
	return s.seasonMiscData.GetCurrentSeason(serverStartDuration)
}

func (s *SeasonService) GetSeasonTickTime() tickdata.TickTime {
	//return s.ticker.GetTickTime()
	return s.tickService.GetDailyTickTime()
}

func (s *SeasonService) OnHeroOnline(hc iface.HeroController) {
	prevTickTime := s.GetSeasonTickTime().GetPrevTickTime()
	season := s.SeasonByTime(s.timeService.CurrentTime())
	hc.Send(domestic.NewS2cSeasonStartBroadcastMsg(season.Season, timeutil.Marshal32(prevTickTime), false))
}
