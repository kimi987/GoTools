package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"sync"
	"time"
)

func NewCountry(defCountry *country.CountryData, datas interface {
	GetCountryOfficialDataArray() []*country.CountryOfficialData
}) *Country {
	c := &Country{}
	c.id = defCountry.Id
	c.prestige = defCountry.DefaultPrestige
	c.name = defCountry.Name
	c.data = defCountry

	c.officials = c.newOfficials(datas)

	c.vote = NewCountryChangeNameVote()

	return c
}

func (c *Country) newOfficials(datas interface {
	GetCountryOfficialDataArray() []*country.CountryOfficialData
}) (officials map[shared_proto.CountryOfficialType]*countryOfficial) {

	officials = make(map[shared_proto.CountryOfficialType]*countryOfficial)

	for _, d := range datas.GetCountryOfficialDataArray() {
		officials[d.OfficialType] = newCountryOfficial(d)
	}

	return
}

type Country struct {
	id       uint64
	prestige uint64

	name string

	data *country.CountryData

	officials map[shared_proto.CountryOfficialType]*countryOfficial

	destroyed bool

	vote *CountryChangeNameVote

	lock sync.RWMutex
}



func (c *Country) Id() uint64 {
	return c.id
}

func (c *Country) Prestige() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.prestige
}

func (c *Country) AddPrestige(toAdd uint64) uint64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.prestige += toAdd

	return c.prestige
}

func (c *Country) ReducePrestige(toReduce uint64) uint64 {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.prestige = u64.Sub(c.prestige, toReduce)

	return c.prestige
}

func (c *Country) Name() string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.name
}

func (c *Country) EncodeBasic(srv interface {
	GetBasicProto(int64) *shared_proto.HeroBasicProto
}) *shared_proto.CountryProto {
	p := &shared_proto.CountryProto{}
	var kid int64

	func() {
		c.lock.RLock()
		defer c.lock.RUnlock()

		p.Id = int32(c.id)
		p.Prestige = u64.Int32(c.prestige)
		p.Name = c.name
		p.Destroyed = c.destroyed
		kid = c.King()
	}()

	if kid > 0 {
		p.King = srv.GetBasicProto(kid)
	}

	return p
}

func (c *Country) Encode(srv interface {
	GetBasicProto(int64) *shared_proto.HeroBasicProto
}, changeNameCd time.Duration) *shared_proto.CountryDetailProto {

	p := &shared_proto.CountryDetailProto{}
	p.Basic = c.EncodeBasic(srv)

	var tmpOffs []*countryOfficial
	func(){
		c.lock.RLock()
		defer c.lock.RUnlock()

		for _, o := range c.officials {
			tmp := newCountryOfficial(o.typeData)
			for hid, h := range o.heros {
				tmp.heros[hid] = h.copy()
			}
			tmpOffs = append(tmpOffs, tmp)
		}
		p.Vote = c.vote.encode()
		p.NextChangeNameTime = timeutil.Marshal32(c.vote.endTime.Add(changeNameCd))
	}()

	for _, o := range tmpOffs {
		p.Officials = append(p.Officials, encodeCountryOfficial(o, srv))
	}

	return p
}

func (c *Country) EncodeServer() *server_proto.CountryServerProto {
	c.lock.RLock()
	defer c.lock.RUnlock()

	p := &server_proto.CountryServerProto{}
	p.Id = c.id
	p.Prestige = c.prestige
	p.Destroyed = c.destroyed
	p.Name = c.name
	for _, o := range c.officials {
		p.Officials = append(p.Officials, o.encodeServer())
	}
	p.Vote = c.vote.encodeServer()

	return p
}

func (c *Country) Unmarshal(p *server_proto.CountryServerProto, datas interface {
	GetCountryOfficialData(int) *country.CountryOfficialData
}) {
	if p == nil {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.id = p.Id
	c.prestige = p.Prestige
	if p.Name != "" {
		c.name = p.Name
	}
	c.destroyed = p.Destroyed
	for _, op := range p.Officials {
		o := c.officials[op.Type]
		if o == nil {
			logrus.Errorf("Country.Unmarshal, 发现 proto 中保存了配置中没有的 CountryOfficialData。country:%v type:%v", p.Id, op.Type)
			continue
		}
		o.unmarshal(op, datas)
	}
	c.vote.unmarshal(p.Vote)
}




func (c *Country) CancelDestroy() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.destroyed = false
}


func (c *Country) Destroy() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.destroyed = true
}

func (c *Country) IsDestroyed() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.destroyed
}




