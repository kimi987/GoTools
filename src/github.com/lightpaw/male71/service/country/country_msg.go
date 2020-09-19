package country

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/country"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"sync"
	"time"
)

func NewCountryMsg(timeSrv iface.TimeService) *CountryMsg {
	c := &CountryMsg{}
	c.msgVer = atomic.NewUint64(1)
	c.emptyPrestigeMsg = country.NewS2cRequestCountryPrestigeMsg(0, nil, nil).Static()

	c.countriesMsgCache = msg.NewMsgCache(60*time.Second, timeSrv)
	c.countriesNoticeMsgCache = msg.NewMsgCache(60*time.Second, timeSrv)
	c.countryDetailMsgCache = msg.NewMsgCache(60*time.Second, timeSrv)

	return c
}

type CountryMsg struct {
	lock   sync.RWMutex
	msgVer *atomic.Uint64 // 版本号没有用了

	prestigeMsg      pbutil.Buffer
	emptyPrestigeMsg pbutil.Buffer

	tutorialCountriesProto *shared_proto.CountriesProto
	countriesProto         *shared_proto.CountriesProto
	mcWarProto             *shared_proto.McWarProto
	mingcsProto            *shared_proto.MingcsProto

	countriesMsgCache       *msg.MsgCache
	countriesNoticeMsgCache *msg.MsgCache
	countryDetailMsgCache   *msg.MsgCache
}

// 调用方不要改返回的 proto
func (c *CountryMsg) TutorialCountriesProto() *shared_proto.CountriesProto {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.tutorialCountriesProto
}

func (c *CountryMsg) updateCountriesMsg(countriesProto *shared_proto.CountriesProto) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.countriesProto = countriesProto

	t := &shared_proto.CountriesProto{}
	for _, cc := range countriesProto.Countries {
		t.Countries = append(t.Countries, &shared_proto.CountryProto{Id: cc.Id, Name: cc.Name})
	}
	c.tutorialCountriesProto = t

	c.countriesMsgCache.Disable(0)
	c.countriesNoticeMsgCache.Disable(0)

	// 旧版本的消息不知道客户端还用不用了
	pp := &country.S2CRequestCountryPrestigeProto{}
	pp.Vsn = u64.Int32(c.msgVer.Load())
	for _, c := range c.countriesProto.Countries {
		pp.Ids = append(pp.Ids, c.Id)
		pp.Prestige = append(pp.Prestige, c.Prestige)
	}
	c.prestigeMsg = country.NewS2cRequestCountryPrestigeProtoMsg(pp).Static()
}

func (c *CountryMsg) updateMcWarMsg(mcWarProto *shared_proto.McWarProto) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.mcWarProto = mcWarProto

	c.countriesMsgCache.Disable(0)
}

func (c *CountryMsg) updateMingcsMsg(mingcsProto *shared_proto.MingcsProto) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.mingcsProto = mingcsProto

	c.countriesMsgCache.Disable(0)
}

func (m *CountryMsg) getCountriesMsg(ver uint64) pbutil.Buffer {
	m.lock.RLock()
	defer m.lock.RUnlock()

	message := m.countriesMsgCache.GetOrUpdate(0, func() (result pbutil.Buffer) {
		return country.NewS2cRequestCountriesMsg(0, m.countriesProto, m.mcWarProto, m.mingcsProto)
	})

	return message
}

func (m *CountryMsg) countriesNoticeMsg() pbutil.Buffer {
	m.lock.RLock()
	defer m.lock.RUnlock()

	message := m.countriesNoticeMsgCache.GetOrUpdate(0, func() (result pbutil.Buffer) {
		return country.NewS2cCountriesUpdateNoticeMsg(m.countriesProto)
	})

	return message
}

func (m *CountryMsg) countryPrestigeMsg(ver uint64) pbutil.Buffer {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if ver != 0 && ver == m.msgVer.Load() {
		return m.emptyPrestigeMsg
	}
	return m.prestigeMsg
}
