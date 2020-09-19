package mingc

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/mingc"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/ticker"
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"sync"
	"time"
)

func NewMingcService(datas iface.ConfigDatas, db iface.DbService, tickerService iface.TickerService,
	guild iface.GuildService, guildSnapshot iface.GuildSnapshotService, countryService iface.CountryService,
	timeService iface.TimeService, srvConf iface.IndividualServerConfig, mail iface.MailModule, chat iface.ChatService,
	heroService iface.HeroDataService, broadcast iface.BroadcastService, world iface.WorldService, tlog iface.TlogService) *MingcService {
	m := &MingcService{}
	m.db = db
	m.datas = datas
	m.heroService = heroService
	m.tickerService = tickerService
	m.guildSnapshot = guildSnapshot
	m.guild = guild
	m.time = timeService
	m.countryService = countryService
	m.svrConf = srvConf
	m.broadcast = broadcast
	m.world = world
	m.tlog = tlog
	m.mail = mail
	m.chat = chat

	ctime := timeService.CurrentTime()
	nextResetTime := singleton.GetNextResetDailyTime(ctime, datas.MingcMiscData().DailyUpdateMingcTime)
	m.dailyTicker = ticker.NewTicker(ctime, nextResetTime.Sub(ctime), 24*time.Hour)

	m.mingcs = make(map[uint64]*entity.Mingc)

	m.msgVer = atomic.NewUint64(1)
	m.emptyMingcsMsg = mingc.NewS2cMingcListMsg(0, nil)

	m.msgCache = msg.NewMsgCache(60*time.Second, timeService)

	m.init()

	m.stopFunc = tickerService.TickPer10Minute("定时保存名城数据", func(tick tickdata.TickTime) {
		ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			m.save(ctx)
			return nil
		})
	})

	go call.CatchLoopPanic(m.loop, "MingcService.loop")

	return m
}

//gogen:iface
type MingcService struct {
	db             iface.DbService
	datas          iface.ConfigDatas
	tickerService  iface.TickerService // 定时保存数据的 tiker
	dailyTicker    *ticker.Ticker      // 每日更新银两的 tiker
	heroService    iface.HeroDataService
	guild          iface.GuildService
	guildSnapshot  iface.GuildSnapshotService
	time           iface.TimeService
	countryService iface.CountryService
	svrConf        iface.IndividualServerConfig
	broadcast      iface.BroadcastService
	world          iface.WorldService
	tlog           iface.TlogService
	mail           iface.MailModule
	chat           iface.ChatService

	msgCache *msg.MsgCache

	stopFunc func()

	mingcs map[uint64]*entity.Mingc

	allInOneGuildId int64

	lock           sync.RWMutex
	msgVer         *atomic.Uint64
	mingcsMsg      pbutil.Buffer
	emptyMingcsMsg pbutil.Buffer
}

func (m *MingcService) DisableMcBuildLogCache(mcId uint64) {
	m.msgCache.Disable(mcId)
}

func (m *MingcService) McBuildLogMsg(mc *entity.Mingc) (msg pbutil.Buffer) {
	msg = m.msgCache.Get(mc.Id())
	if msg == nil {
		m.msgCache.Update(mc.Id(), func() (result pbutil.Buffer) {
			logs := m.updateMcBuildLog(mc)
			msg = mingc.NewS2cMcBuildLogMsg(logs)
			return msg
		})
	}

	return
}

func (m *MingcService) updateMcBuildLog(mc *entity.Mingc) (logs *shared_proto.GuildMcBuildProto) {
	logs = &shared_proto.GuildMcBuildProto{}

	mc.WalkSupportGuilds(func(gid int64, buildCount, support uint64) {
		log := &shared_proto.SingleGuildMcBuildProto{}
		g := m.guildSnapshot.GetGuildBasicProto(gid)
		if g == nil {
			logrus.Warnf("名城营建记录，找不到联盟:%v, 可能是解散了", gid)
			return
		}

		log.Guild = g
		log.BuildCount = u64.Int32(buildCount)
		log.Support = u64.Int32(support)

		logs.Log = append(logs.Log, log)
	})

	// 排序
	entity.SortDescGuildMcBuildLogs(logs)
	maxLen := m.datas.McBuildMiscData().MaxMcBuildLogCount
	if len(logs.Log) > maxLen {
		logs.Log = logs.Log[:maxLen]
	}

	return logs
}

func (m *MingcService) WalkMingcs(f entity.MingcFunc) {
	for _, c := range m.mingcs {
		f(c)
	}
}

func (m *MingcService) loop() {

	//dailyTickTime := m.tickerService.GetDailyTickTime()
	dailyTickTime := m.dailyTicker.GetTickTime()
	m.resetDaily(dailyTickTime.GetPrevTickTime())

	for {
		select {
		case <-dailyTickTime.Tick():
			m.resetDaily(dailyTickTime.GetTickTime())
			dailyTickTime = m.dailyTicker.GetTickTime()
		}
	}
}

func (m *MingcService) resetDaily(resetTime time.Time) {
	// 名城营建奖励
	m.resetDailyMcBuild(resetTime)

	mcWarFirstStartTime := m.svrConf.GetServerStartTime().Add(m.datas.MingcMiscData().StartAfterServerOpen)
	if resetTime.Before(mcWarFirstStartTime) {
		return
	}

	for _, mc := range m.mingcs {
		if !mc.LastResetTime().Before(resetTime) {
			continue
		}

		// 在名城每日重置之前执行
		m.updateMcBuildCountInAllGuild(mc, true)

		supportData := m.datas.GetMcBuildMcSupportData(mc.Level())
		mc.ResetDaily(resetTime, m.datas.GetMingcBaseData(mc.Id()), supportData)

		m.msgCache.Disable(mc.Id())

		if mc.HostGuildId() > 0 {
			// 加名城占领盟银两
			m.guild.FuncGuild(mc.HostGuildId(), func(g *sharedguilddata.Guild) {
				if g == nil {
					return
				}

				hostExtra := mc.CleanExtraYinliang()
				toAdd := m.datas.GetMingcBaseData(mc.Id()).HostDailyAddYinliang + hostExtra + supportData.AddHostDailyYinliang
				g.AddYinliang(toAdd)

				if d := m.datas.GuildLogHelp().YinliangMingcHost; d != nil {
					mingcName := m.datas.GetMingcBaseData(mc.Id()).Name
					text := d.Text.New().WithMingc(mingcName).WithYinliang(toAdd).JsonString()
					g.AddYinliangRecord(text, resetTime)
					if d.SendChat {
						m.chat.SysChat(0, g.Id(), shared_proto.ChatType_ChatGuild, text, shared_proto.ChatMsgType_ChatMsgGuildLog, true, true, true, false)
					}
				}
			})
		}
	}
	m.UpdateMsg()
	m.world.Broadcast(mingc.RESET_DAILY_MC_S2C)
}

type tmpHeroMcBuildInfo struct {
	heroId int64
	gid    int64
	mcId   uint64
	mcName string
	count  uint64
}

func (m *MingcService) resetDailyMcBuild(resetTime time.Time) {
	heros := make(map[int64]map[uint64]*tmpHeroMcBuildInfo)

	for _, mc := range m.mingcs {
		if !mc.LastResetTime().Before(resetTime) {
			continue
		}

		mcName := m.datas.MingcBaseData().Get(mc.Id()).Name

		// 给参与营建的人发奖励，每人发自己营建的名城中，联盟营建次数最高的奖励
		mc.WalkSupportHeros(func(heroId, gid int64, guildBuildCount uint64) {
			mcc, ok := heros[heroId]
			if !ok {
				mcc = make(map[uint64]*tmpHeroMcBuildInfo)
				heros[heroId] = mcc
			}
			if guildBuildCount > 0 {
				mcc[mc.Id()] = &tmpHeroMcBuildInfo{heroId: heroId, gid: gid, mcId: mc.Id(), mcName: mcName, count: guildBuildCount}
			}
		})

		// 每日重置名城营建
		supportData := m.datas.GetMcBuildMcSupportData(mc.Level())
		mc.ResetDailyMcBuild(resetTime, m.datas.McBuildMiscData(), supportData)
	}

	for heroId, infos := range heros {
		var maxBuildCountInfo *tmpHeroMcBuildInfo
		var guildId int64
		m.heroService.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
			if hero.GuildId() <= 0 {
				return
			}

			guildId = hero.GuildId()

			for _, info := range infos {
				if info.gid != hero.GuildId() {
					continue
				}
				if maxBuildCountInfo == nil {
					maxBuildCountInfo = info
					continue
				}
				if info.count <= maxBuildCountInfo.count {
					continue
				}
				maxBuildCountInfo = info
			}
		})

		if maxBuildCountInfo != nil {
			if prize := m.GetMcBuildGuildMemberPrize(maxBuildCountInfo.count); prize != nil {
				if d := m.datas.MailHelp().McBuildGuildMemberPrize; d != nil {
					guildName := m.guild.GetGuildFlagName(guildId)
					proto := d.NewTextMail(shared_proto.MailType_MailNormal)
					proto.Text = d.NewTextFields().WithGuild(guildName).WithMingc(maxBuildCountInfo.mcName).WithCount(maxBuildCountInfo.count).JsonString()
					proto.Prize = prize.PrizeProto()
					m.mail.SendProtoMail(heroId, proto, resetTime)
				}
			}
		}
	}
}

func (m *MingcService) GetMcBuildGuildMemberPrize(buildCount uint64) *resdata.Prize {
	for _, d := range m.datas.GetMcBuildGuildMemberPrizeDataArray() {
		if d.BuildCountRage.ContainsClosed(buildCount) {
			return d.Prize
		}
	}
	return nil
}

func (m *MingcService) init() {
	var bytes []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		bytes, err = m.db.LoadKey(ctx, server_proto.Key_Mingc)
		return
	})
	if err != nil {
		logrus.WithError(err).Panic("加载mingc模块数据失败")
	}

	minSupportLevel := m.datas.McBuildMcSupportData().MinKeyData.Level

	if len(bytes) <= 0 {
		// 初始名城
		for _, d := range m.datas.GetMingcBaseDataArray() {
			m.mingcs[d.Id] = entity.NewMingc(d, minSupportLevel)
		}
		return
	}

	proto := &server_proto.MingcsServerProto{}
	if err := proto.Unmarshal(bytes); err != nil {
		logrus.WithError(err).Panic("解析mingc模块数据失败")
	}

	for _, p := range proto.Mcs {
		if m.datas.GetMingcBaseData(p.Id) == nil {
			logrus.Warnf("MingcService.init 名城：%v 配置被删掉了, 忽略 Unmarshal", p.Id)
			continue
		}

		mc := &entity.Mingc{}
		mc.Unmarshal(p, m.datas)
		m.mingcs[mc.Id()] = mc
	}

	// 在名城初始化完成后执行
	m.initAllMcBuildCountInGuild()

	m.updateAllInOneGuild()
	m.UpdateMsg()
}

// 更新联盟对象里的名城营建数据
func (m *MingcService) initAllMcBuildCountInGuild() {
	for _, mc := range m.mingcs {
		m.updateMcBuildCountInAllGuild(mc, false)
	}
}

// 更新联盟对象里的名城营建数据
func (m *MingcService) updateMcBuildCountInAllGuild(mc *entity.Mingc, reset bool) {
	mc.WalkSupportGuilds(func(gid int64, buildCount, support uint64) {
		if reset {
			m.updateMcBuildCountInGuild(gid, 0, 0, mc)
		} else {
			m.updateMcBuildCountInGuild(gid, buildCount, support, mc)
		}
	})
}

func (m *MingcService) updateMcBuildCountInGuild(gid int64, buildCount, support uint64, mc *entity.Mingc) {
	m.guild.FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}
		if g.McBuildCount(mc.Id()) == buildCount {
			return
		}
		g.SetMcBuildCount(mc.Id(), buildCount)
		m.guild.ClearSelfGuildMsgCache(gid)
		m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)
	})
}

func (m *MingcService) Close() {
	if m.stopFunc != nil {
		m.stopFunc()
	}

	// 下线保存一次
	ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = m.save(ctx)
		return
	})
}

func (m *MingcService) save(ctx context.Context) (err error) {
	if err = m.db.SaveKey(ctx, server_proto.Key_Mingc, must.Marshal(m.encodeServer())); err != nil {
		logrus.WithError(err).Panic("保存mingc模块数据失败")
	}
	return
}

func (m *MingcService) encodeServer() *server_proto.MingcsServerProto {
	proto := &server_proto.MingcsServerProto{}
	for _, c := range m.mingcs {
		proto.Mcs = append(proto.Mcs, c.EncodeServer())
	}

	return proto
}

func (m *MingcService) encode() *shared_proto.MingcsProto {
	p := &shared_proto.MingcsProto{}
	for _, mc := range m.mingcs {
		p.Mingcs = append(p.Mingcs, mc.Encode(m.datas.GetMingcBaseData(mc.Id()), m.guildSnapshot.GetGuildBasicProto))
	}

	return p
}

func (m *MingcService) UpdateMsg() pbutil.Buffer {

	proto := func() *shared_proto.MingcsProto {
		m.lock.Lock()
		defer m.lock.Unlock()

		ver := m.msgVer.Inc()
		proto := m.encode()
		m.mingcsMsg = mingc.NewS2cMingcListMsg(u64.Int32(ver), proto).Static()
		return proto
	}()

	m.countryService.UpdateMingcsMsg(proto)

	return m.mingcsMsg
}

func (m *MingcService) Mingc(id uint64) *entity.Mingc {
	return m.mingcs[id]
}

func (m *MingcService) getMingcsMsg() pbutil.Buffer {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.mingcsMsg
}

func (m *MingcService) MingcsMsg(ver uint64) pbutil.Buffer {
	if ver != 0 && ver == m.msgVer.Load() {
		return m.emptyMingcsMsg
	}
	if msg := m.getMingcsMsg(); msg != nil {
		return msg
	}
	return m.UpdateMsg()
}

func (m *MingcService) GuildMingc(gid int64) *entity.Mingc {
	for _, c := range m.mingcs {
		if c.HostGuildId() == gid {
			return c
		}
	}
	return nil
}

func (m *MingcService) Country(mcId uint64) uint64 {
	gid := m.Mingc(mcId).HostGuildId()
	if gid <= 0 {
		return m.datas.GetMingcBaseData(mcId).Country
	}

	if g := m.guildSnapshot.GetSnapshot(gid); g != nil {
		return g.Country.Id
	} else {
		return m.datas.GetMingcBaseData(mcId).Country
	}
}

func (m *MingcService) SetHostGuild(mcId uint64, gid int64) (succ bool) {
	if mc := m.Mingc(mcId); mc == nil {
		return
	} else {
		mc.SetHostGuildId(gid)
	}

	succ = true
	m.updateAllInOneGuild()
	m.UpdateMsg()

	return
}

func (m *MingcService) updateAllInOneGuild() {
	var gid int64
	for _, mc := range m.mingcs {
		if gid <= 0 {
			gid = mc.HostGuildId()
			continue
		}
		if gid != mc.HostGuildId() {
			m.allInOneGuildId = 0
			return
		}
	}
	if gid > 0 {
		m.allInOneGuildId = gid
	}
}

func (m *MingcService) AllInOneGuild() int64 {
	return m.allInOneGuildId
}

func (m *MingcService) Build(heroId, gid int64, baiZhanLevel, mcId uint64) bool {
	mc := m.mingcs[mcId]
	if mc == nil {
		return false
	}

	toAdd := m.getToAddSupport(heroId, gid, baiZhanLevel, mc)
	mc.AddSupport(toAdd, gid, heroId, m.datas)
	m.UpdateMsg()

	m.DisableMcBuildLogCache(mc.Id())

	count, support := mc.GuildBuildCount(gid)
	m.updateMcBuildCountInGuild(gid, count, support, mc)

	return true
}

func (m *MingcService) getToAddSupport(heroId, gid int64, baiZhanLevel uint64, mc *entity.Mingc) (toAdd uint64) {
	mcData := m.datas.GetMingcBaseData(mc.Id())
	if mcData == nil {
		return
	}

	g := m.guildSnapshot.GetSnapshot(gid)
	if g == nil {
		return
	}

	if g.Country.Id != mcData.Country {
		return
	}

	if d := m.datas.GetMcBuildAddSupportData(baiZhanLevel); d != nil {
		toAdd = d.AddSupport
	}

	toAdd = u64.Min(toAdd, u64.Sub(m.datas.McBuildMiscData().MaxDailyAddSupport, mc.Support()))

	return
}

// 国家当前占领的本国初始名城
func (m *MingcService) CountryHoldInitMcs(countryId uint64) (mcs []*entity.Mingc) {
	for _, mc := range m.mingcs {
		d := m.datas.GetMingcBaseData(mc.Id())
		if d.Country != countryId {
			continue
		}

		if m.isMcInInitCountry(mc, d) {
			mcs = append(mcs, mc)
		}
	}

	return
}

func (m *MingcService) isMcInInitCountry(mc *entity.Mingc, d *mingcdata.MingcBaseData) bool {
	if mc.HostGuildId() <= 0 {
		return true
	}

	g := m.guildSnapshot.GetSnapshot(mc.HostGuildId())
	if g == nil {
		return true
	}

	return g.Country.Id == d.Country
}

func (m *MingcService) CaptainHostGuild(countryId uint64) (gid int64) {
	cdata := m.datas.GetCountryData(countryId)
	if cdata == nil {
		return
	}

	gid = m.Mingc(cdata.Capital.Id).HostGuildId()

	return
}

func (m *MingcService) IsHoldCountryCapital(gid int64, countryId uint64) bool {
	country := m.datas.GetCountryData(countryId)
	if country == nil {
		return false
	}

	return m.Mingc(country.Capital.Id).HostGuildId() == gid
}
