package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"sort"
	"time"
)

func newMcBuild(minLevel uint64) *mcBuild {
	m := &mcBuild{}
	m.level = minLevel
	m.guildInfos = make(map[int64]*mcBuildGuildInfo)
	return m
}

type mcBuild struct {
	level             uint64
	support           uint64
	dailyAddedSupport uint64
	guildInfos        map[int64]*mcBuildGuildInfo
}

func (m *mcBuild) encodeMcBuildServer() *server_proto.McBuildProto {
	p := &server_proto.McBuildProto{}
	p.Level = m.level
	p.Support = m.support
	p.DailyAddedSupport = m.dailyAddedSupport

	p.GuildInfos = make(map[int64]*server_proto.McBuildGuildInfoProto)
	for gid, g := range m.guildInfos {
		p.GuildInfos[gid] = g.encodeServer()
	}

	return p
}

func (m *mcBuild) unmarshal(p *server_proto.McBuildProto, datas interface {
	McBuildMcSupportData() *config.McBuildMcSupportDataConfig
}) {
	if p == nil {
		return
	}
	m.level = datas.McBuildMcSupportData().Must(p.Level).Level
	if m.level != p.Level {
		logrus.Errorf("unmarshal mcBuild, 找不到 supportLevel:%v", p.Level)
	}

	m.support = p.Support
	m.dailyAddedSupport = p.DailyAddedSupport

	if p.GuildInfos != nil {
		for gid, gp := range p.GuildInfos {
			g := newMcBuildGuildInfo()
			g.unmarshal(gp)
			m.guildInfos[gid] = g
		}
	}
}

func (m *mcBuild) addSupport(toAdd uint64, gid, heroId int64) {
	m.support += toAdd
	m.dailyAddedSupport += toAdd

	g := m.guildInfos[gid]
	if g == nil {
		g = newMcBuildGuildInfo()
		m.guildInfos[gid] = g
	}
	g.addSupport(toAdd, heroId)
}

func (m *mcBuild) reduceSupport(toReduce uint64, supportData *mingcdata.McBuildMcSupportData) (levelReduced bool) {
	if toReduce <= 0 {
		return
	}

	if supportData == nil || supportData.PrevLevelData == nil {
		return
	}

	allSupport := m.support
	newData := supportData
	for prevData := newData.PrevLevelData; allSupport < toReduce; prevData = prevData.PrevLevelData {
		if prevData == nil {
			break
		}
		allSupport += prevData.UpgradeSupport
		newData = prevData

		levelReduced = true
	}

	m.support = u64.Sub(allSupport, toReduce)
	m.level = newData.Level

	return
}

func newMcBuildGuildInfo() *mcBuildGuildInfo {
	g := &mcBuildGuildInfo{}
	g.heroInfos = make(map[int64]*mcBuildHeroInfo)
	return g
}

type mcBuildGuildInfo struct {
	support    uint64
	buildCount uint64
	heroInfos  map[int64]*mcBuildHeroInfo
}

func (m *mcBuildGuildInfo) encodeServer() *server_proto.McBuildGuildInfoProto {
	p := &server_proto.McBuildGuildInfoProto{}
	p.Support = m.support
	p.BuildCount = m.buildCount

	p.HeroInfos = make(map[int64]*server_proto.McBuildHeroInfoProto)
	for hid, h := range m.heroInfos {
		p.HeroInfos[hid] = h.encodeServer()
	}

	return p
}

func (m *mcBuildGuildInfo) unmarshal(p *server_proto.McBuildGuildInfoProto) {
	m.support = p.Support
	m.buildCount = p.BuildCount

	if p.HeroInfos != nil {
		for hid, hp := range p.HeroInfos {
			h := newMcBuildHeroInfo()
			h.unmarshal(hp)
			m.heroInfos[hid] = h
		}
	}
}

func (g *mcBuildGuildInfo) addSupport(toAdd uint64, heroId int64) {
	g.support += toAdd
	g.buildCount++

	h := g.heroInfos[heroId]
	if h == nil {
		h = newMcBuildHeroInfo()
		g.heroInfos[heroId] = h
	}
	h.support += toAdd
	h.buildCount++
}

func newMcBuildHeroInfo() *mcBuildHeroInfo {
	m := &mcBuildHeroInfo{}
	return m
}

type mcBuildHeroInfo struct {
	support    uint64
	buildCount uint64
}

func (m *mcBuildHeroInfo) encodeServer() *server_proto.McBuildHeroInfoProto {
	p := &server_proto.McBuildHeroInfoProto{}
	p.Support = m.support
	p.BuildCount = m.buildCount

	return p
}

func (m *mcBuildHeroInfo) unmarshal(p *server_proto.McBuildHeroInfoProto) {
	m.support = p.Support
	m.buildCount = p.BuildCount
}

func NewHeroMcBuild() *HeroMcBuild {
	h := &HeroMcBuild{}
	h.NextTime = time.Time{}
	return h
}

type HeroMcBuild struct {
	BuildCount uint64
	NextTime   time.Time
}

func (h *HeroMcBuild) Encode() *shared_proto.HeroMcBuildProto {
	p := &shared_proto.HeroMcBuildProto{}
	p.McBuildCount = u64.Int32(h.BuildCount)
	p.McBuildNextTime = timeutil.Marshal32(h.NextTime)
	return p
}

func (h *HeroMcBuild) encodeServer() *server_proto.HeroMcBuildServerProto {
	p := &server_proto.HeroMcBuildServerProto{}
	p.McBuildCount = h.BuildCount
	p.McBuildNextTime = timeutil.Marshal64(h.NextTime)
	return p
}

func (h *HeroMcBuild) unmarshal(p *server_proto.HeroMcBuildServerProto) {
	if p == nil {
		return
	}

	h.BuildCount = p.McBuildCount
	h.NextTime = timeutil.Unix64(p.McBuildNextTime)
}

func (h *HeroMcBuild) ResetDaily() {
	h.BuildCount = 0
	h.NextTime = time.Time{}
}

func (h *HeroMcBuild) Build(newTime time.Time) {
	h.BuildCount++
	h.NextTime = newTime
}

func (m *Mingc) AddSupport(toAdd uint64, gid, heroId int64, conf interface {
	GetMcBuildMcSupportData(uint64) *mingcdata.McBuildMcSupportData
}) (upgraded bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	// toAdd == 0 时也要继续执行，后面会加营建次数
	if toAdd < 0 {
		return
	}

	currData := conf.GetMcBuildMcSupportData(m.level)
	if currData == nil {
		return
	}

	m.addSupport(toAdd, gid, heroId)

	if currData.NextLevelData == nil {
		m.support = u64.Min(currData.UpgradeSupport, m.support)
		return
	}

	// 升级
	for d := currData; m.support >= d.UpgradeSupport; d = conf.GetMcBuildMcSupportData(m.level) {
		nextLevel := d.NextLevelData
		if nextLevel == nil {
			m.support = u64.Min(d.UpgradeSupport, m.support)
			break
		}
		m.support = u64.Sub(m.support, d.UpgradeSupport)
		m.level = nextLevel.Level

		upgraded = true
	}

	return
}

type tempMcBuildInfo struct {
	heroId          int64
	gid             int64
	guildBuildCount uint64
	support         uint64
}

func (m *Mingc) WalkSupportHeros(f func(heroId, gid int64, guildBuildCount uint64)) {
	for _, h := range m.copySupportHeros() {
		f(h.heroId, h.gid, h.guildBuildCount)
	}
}

func (m *Mingc) copySupportHeros() (heros []*tempMcBuildInfo) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for gid, g := range m.guildInfos {
		if g == nil {
			continue
		}
		for hid := range g.heroInfos {
			heros = append(heros, &tempMcBuildInfo{heroId: hid, gid: gid, guildBuildCount: g.buildCount})
		}
	}
	return
}

func (m *Mingc) WalkSupportGuilds(f func(gid int64, guildBuildCount, support uint64)) {
	for _, h := range m.copySupportGuilds() {
		f(h.gid, h.guildBuildCount, h.support)
	}
}

func (m *Mingc) copySupportGuilds() (guilds []*tempMcBuildInfo) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for gid, g := range m.guildInfos {
		if g == nil {
			continue
		}
		guilds = append(guilds, &tempMcBuildInfo{gid: gid, guildBuildCount: g.buildCount, support: g.support})
	}

	return
}

func (m *Mingc) GuildBuildCount(gid int64) (count, support uint64) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if g := m.mcBuild.guildInfos[gid]; g != nil {
		return g.buildCount, g.support
	}
	return
}

func (m *Mingc) resetDailyMcBuild(miscData *mingcdata.McBuildMiscData, supportData *mingcdata.McBuildMcSupportData) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.dailyAddedSupport = 0
	m.reduceSupport(miscData.DailyReduceSupport, supportData)
	m.guildInfos = make(map[int64]*mcBuildGuildInfo)
}

func (m *Mingc) Level() uint64 {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.level
}

func (m *Mingc) DailyAddedSupport() uint64 {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.dailyAddedSupport
}

func (m *Mingc) Support() uint64 {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.support
}

func SortDescGuildMcBuildLogs(logs *shared_proto.GuildMcBuildProto) {
	sort.Sort(GuildMcBuildProtoSorter(logs.Log))
}

type GuildMcBuildProtoSorter []*shared_proto.SingleGuildMcBuildProto

func (p GuildMcBuildProtoSorter) Less(i, j int) bool {
	if p[i].Support > p[j].Support {
		return true
	} else if p[i].Support == p[j].Support {
		if p[i].BuildCount > p[j].BuildCount {
			return true
		}
	}
	return false
}

func (p GuildMcBuildProtoSorter) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p GuildMcBuildProtoSorter) Len() int      { return len(p) }
