package entity

import (
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"sync"
	"time"
)

type MingcFunc func(c *Mingc)

func NewMingc(d *mingcdata.MingcBaseData, minSupprotLevel uint64) *Mingc {
	c := &Mingc{}
	c.id = d.Id
	c.yinliang = d.DefaultYinliang
	c.lastResetTime = time.Time{}

	c.mcBuild = newMcBuild(minSupprotLevel)

	return c
}

type Mingc struct {
	id uint64

	*mcBuild

	hostGuildId int64

	yinliang uint64 // 仓库正常银两，不会超过仓库容量

	extraYinliang uint64 // 仓库爆仓银两

	lastResetTime time.Time

	lock sync.RWMutex
}

func (c *Mingc) Unmarshal(p *server_proto.MingcServerProto, datas interface {
	McBuildMcSupportData() *config.McBuildMcSupportDataConfig
}) {
	c.id = p.Id
	c.hostGuildId = p.HostGuildId
	c.yinliang = p.Yinliang
	c.extraYinliang = p.ExtraYinliang
	c.lastResetTime = timeutil.Unix64(p.LastResetTime)

	minLevel := datas.McBuildMcSupportData().MinKeyData.Level
	c.mcBuild = newMcBuild(minLevel)
	c.mcBuild.unmarshal(p.McBuild, datas)
}

func (c *Mingc) EncodeServer() *server_proto.MingcServerProto {
	c.lock.RLock()
	defer c.lock.RUnlock()

	p := &server_proto.MingcServerProto{}
	p.Id = c.id
	p.HostGuildId = c.hostGuildId
	p.Yinliang = c.yinliang
	p.ExtraYinliang = c.extraYinliang
	p.LastResetTime = timeutil.Marshal64(c.lastResetTime)

	p.McBuild = c.encodeMcBuildServer()

	return p
}

func (c *Mingc) Encode(d *mingcdata.MingcBaseData, guildGetter func(gid int64) *shared_proto.GuildBasicProto) *shared_proto.MingcProto {
	c.lock.RLock()
	defer c.lock.RUnlock()

	p := &shared_proto.MingcProto{}
	p.Id = u64.Int32(c.id)
	p.Yinliang = u64.Int32(c.yinliang)
	p.ExtraYinliang = u64.Int32(c.extraYinliang)
	p.HostExtraYinliang = u64.Int32(u64.Sub(c.yinliang+c.extraYinliang+d.DailyAddYinliang, d.MaxYinliang))
	p.HostGuild = guildGetter(c.hostGuildId)

	p.Level = u64.Int32(c.level)
	p.Support = u64.Int32(c.support)
	p.DailyAddedSupport = u64.Int32(c.dailyAddedSupport)

	return p
}

func (c *Mingc) Id() uint64 {
	return c.id
}

func (c *Mingc) LastResetTime() time.Time {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.lastResetTime
}

func (c *Mingc) ResetDaily(ctime time.Time, d *mingcdata.MingcBaseData, supportData *mingcdata.McBuildMcSupportData) bool {
	c.lock.Lock()
	c.lastResetTime = ctime
	c.lock.Unlock()

	if d == nil {
		return false
	}

	addYinliang := d.DailyAddYinliang + supportData.AddDailyYinliang
	maxYinliang := d.MaxYinliang + supportData.AddMaxYinliang
	c.AddYinliang(addYinliang, maxYinliang)

	return true
}

func (c *Mingc) ResetDailyMcBuild(ctime time.Time, miscData *mingcdata.McBuildMiscData, supportData *mingcdata.McBuildMcSupportData) bool {
	c.resetDailyMcBuild(miscData, supportData)
	return true
}

func (c *Mingc) HostGuildId() int64 {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.hostGuildId
}

func (c *Mingc) SetHostGuildId(gid int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.hostGuildId = gid
}

func (c *Mingc) CleanExtraYinliang() (oldExtra uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	oldExtra = c.extraYinliang
	c.extraYinliang = 0

	return
}

func (c *Mingc) Yinliang() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.yinliang
}

func (c *Mingc) ExtraYinliang() uint64 {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.extraYinliang
}

func (c *Mingc) AddYinliang(toAdd uint64, max uint64) (new, extra uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	all := c.yinliang + toAdd
	if all > max {
		c.yinliang = max
		c.extraYinliang += u64.Sub(all, max)
	} else {
		c.yinliang += toAdd
	}

	return c.yinliang, c.extraYinliang
}

func (c *Mingc) ReducePercentYinliang(percent uint64) (newYinliang, toReduce uint64) {
	if percent == 0 {
		return c.yinliang, 0
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	toReduce = c.yinliang * percent / 100
	c.yinliang = u64.Sub(c.yinliang, toReduce)

	return c.yinliang, toReduce
}

type JoinedMcWarIds struct {
	WarMcIds map[int32][]int32
}

func NewJoinedMcWarIds() *JoinedMcWarIds {
	c := &JoinedMcWarIds{}
	c.WarMcIds = make(map[int32][]int32)
	return c
}

