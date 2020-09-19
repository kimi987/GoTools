package entity

import (
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"time"
)

func newCountryOfficial(d *country.CountryOfficialData) *countryOfficial {
	c := &countryOfficial{}
	c.typeData = d
	c.heros = make(map[int64]*countryOfficialHero)

	return c
}

// 在 Country.lock 中使用
type countryOfficial struct {
	typeData *country.CountryOfficialData
	heros    map[int64]*countryOfficialHero // 玩家:任命时间
}

func NewCountryOfficialHero(id int64, appointTime time.Time, pos int32) *countryOfficialHero {
	h := &countryOfficialHero{}
	h.heroId = id
	h.appointTime = appointTime
	h.pos = pos
	return h
}

type countryOfficialHero struct {
	heroId      int64
	appointTime time.Time
	pos         int32
}

func (c *countryOfficialHero) copy() *countryOfficialHero {
	return NewCountryOfficialHero(c.heroId, c.appointTime, c.pos)
}

func (c *countryOfficialHero) encode(srv interface {
	GetBasicProto(int64) *shared_proto.HeroBasicProto
}) *shared_proto.CountryOfficialHeroProto {
	p := &shared_proto.CountryOfficialHeroProto{}
	p.Hero = srv.GetBasicProto(c.heroId)
	p.AppointTime = timeutil.Marshal32(c.appointTime)
	p.Pos = c.pos

	return p
}

func (c *countryOfficialHero) encodeServer() *server_proto.CountryOfficialHeroServerProto {
	p := &server_proto.CountryOfficialHeroServerProto{}
	p.HeroId = c.heroId
	p.AppointTime = timeutil.Marshal64(c.appointTime)
	p.Pos = c.pos

	return p
}

func (c *countryOfficialHero) unmarshal(p *server_proto.CountryOfficialHeroServerProto) {
	c.heroId = p.HeroId
	c.appointTime = timeutil.Unix64(p.AppointTime)
	c.pos = p.Pos
}

func encodeCountryOfficial(copyOfficial *countryOfficial, srv interface {
	GetBasicProto(int64) *shared_proto.HeroBasicProto
}) *shared_proto.CountryOfficialProto {
	p := &shared_proto.CountryOfficialProto{}
	p.Type = copyOfficial.typeData.OfficialType

	for _, h := range copyOfficial.heros {
		p.Heros = append(p.Heros, h.encode(srv))
	}

	return p
}

func (c *countryOfficial) encodeServer() *server_proto.CountryOfficialServerProto {
	p := &server_proto.CountryOfficialServerProto{}
	p.Type = c.typeData.OfficialType
	p.Heros = make(map[int64]*server_proto.CountryOfficialHeroServerProto)
	for hid, h := range c.heros {
		p.Heros[hid] = h.encodeServer()
	}

	return p
}

func (c *countryOfficial) unmarshal(p *server_proto.CountryOfficialServerProto, datas interface {
	GetCountryOfficialData(int) *country.CountryOfficialData
}) {
	if p == nil {
		return
	}

	c.typeData = datas.GetCountryOfficialData(int(p.Type))
	for hid, hp := range p.Heros {
		c.heros[hid] = NewCountryOfficialHero(hp.HeroId, timeutil.Unix64(hp.AppointTime), hp.Pos)
	}
}

func (c *countryOfficial) officialType() shared_proto.CountryOfficialType {
	return c.typeData.OfficialType
}

func (c *Country) NextAppointTime(heroId int64) (t time.Time) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	off := c.heroOfficial(heroId)
	if off == nil {
		return
	}

	t = off.nextDeposeTime(heroId)
	return
}

func (c *countryOfficial) nextDeposeTime(heroId int64) time.Time {
	if h := c.heros[heroId]; h == nil {
		return time.Time{}
	} else if timeutil.IsZero(h.appointTime) {
		return time.Time{}
	} else {
		return h.appointTime.Add(c.typeData.Cd)
	}
}

func (c *countryOfficial) full() bool {
	return len(c.heros) >= c.typeData.Count
}

func (c *countryOfficial) heroOn(heroId int64, ctime time.Time, pos int32) (succ bool) {
	if c.full() {
		return
	}

	c.heros[heroId] = NewCountryOfficialHero(heroId, ctime, pos)
	succ = true
	return
}

func (c *countryOfficial) heroOff(heroId int64) {
	delete(c.heros, heroId)
}

func (c *Country) heroOfficial(heroId int64) (official *countryOfficial) {
	for _, official := range c.officials {
		if _, ok := official.heros[heroId]; ok {
			return official
		}
	}

	return
}

func (c *Country) OfficialHeros(t shared_proto.CountryOfficialType) (heroIds []int64) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	o := c.officials[t]
	if o == nil {
		return
	}

	for hid := range o.heros {
		heroIds = append(heroIds, hid)
	}

	return
}

// 调用方判断权限
func (c *Country) OfficialAppoint(t shared_proto.CountryOfficialType, heroId int64, ctime time.Time, pos int32) (oldType, newType shared_proto.CountryOfficialType, succ bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	newOfficial, ok := c.officials[t]
	if !ok || newOfficial == nil {
		return
	}

	if newOfficial.full() {
		return
	}

	if currOfficial := c.heroOfficial(heroId); currOfficial != nil {
		return
	}

	newOfficial.heroOn(heroId, ctime, pos)

	succ = true
	newType = newOfficial.officialType()

	return
}

func (c *Country) OfficialNextDeposeTime(t shared_proto.CountryOfficialType, heroId int64) (nextTime time.Time) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	off, ok := c.officials[t]
	if !ok || off == nil {
		return
	}

	return off.nextDeposeTime(heroId)
}

func (c *Country) OfficialFull(t shared_proto.CountryOfficialType) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	off, ok := c.officials[t]
	if !ok || off == nil {
		return false
	}

	return off.full()
}

func (c *Country) OfficialDeposeByOfficial(t shared_proto.CountryOfficialType) (heroIds []int64, succ bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	off := c.officials[t]
	if off == nil {
		return
	}

	for hid := range off.heros {
		heroIds = append(heroIds, hid)
	}
	off.heros = make(map[int64]*countryOfficialHero)

	succ = true
	return
}

func (c *Country) OfficialDepose(heroId int64, ctime time.Time, force bool) (oldType shared_proto.CountryOfficialType, succ bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	currOfficial := c.heroOfficial(heroId)
	if currOfficial == nil {
		return
	}

	if !force && ctime.Before(currOfficial.nextDeposeTime(heroId)) {
		return
	}

	oldType = currOfficial.officialType()
	currOfficial.heroOff(heroId)

	succ = true
	return
}

func (c *Country) OfficialDeposeAll(datas interface {
	GetCountryOfficialDataArray() []*country.CountryOfficialData
}) (officials map[shared_proto.CountryOfficialType][]int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	officials = make(map[shared_proto.CountryOfficialType][]int64)
	for t, heros := range c.officials {
		for hid := range heros.heros {
			officials[t] = append(officials[t], hid)
		}
	}
	c.officials = c.newOfficials(datas)

	return officials
}

func (c *Country) IsSubOfficial(t, subT shared_proto.CountryOfficialType) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	official, ok := c.officials[t]
	if !ok || official == nil {
		return false
	}
	return official.typeData.IsSubOfficial(subT)
}

func (c *Country) King() (kingId int64) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if king := c.officials[shared_proto.CountryOfficialType_COT_KING]; king != nil {
		for id := range king.heros {
			kingId = id
			return
		}
	}
	return
}

func (c *Country) HeroOfficial(heroId int64) (t shared_proto.CountryOfficialType) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if off := c.heroOfficial(heroId); off != nil {
		return off.officialType()
	}

	return
}
