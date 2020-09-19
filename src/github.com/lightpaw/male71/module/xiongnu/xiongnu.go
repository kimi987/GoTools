package xiongnu

import (
	"github.com/lightpaw/logrus"
	xiongnu2 "github.com/lightpaw/male7/config/xiongnu"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/xiongnu"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuface"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuinfo"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timer"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/module/rank/ranklist"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/gen/pb/guild"
)

var (
	updateWheel = timer.NewTimingWheel(500*time.Millisecond, 32) // 最大timeout是16秒. 英雄每秒update也在用这个
)

func NewXiongNuModule(dep iface.ServiceDep,
	configDatas *config.ConfigDatas, // 配置
	realmService iface.RealmService, // 场景数据
	mailModule iface.MailModule,     // 邮件服务
	tickerService iface.TickerService,
	xiongNuService iface.XiongNuService,
	pushService iface.PushService,
	rankModule iface.RankModule,
) *XiongNuModule {

	m := &XiongNuModule{
		dep:                 dep,
		timeService:         dep.Time(),
		configDatas:         configDatas,
		guildService:        dep.Guild(),
		realmService:        realmService,
		heroSnapshotService: dep.HeroSnapshot(),
		mailModule:          mailModule,
		worldService:        dep.World(),
		tickerService:       tickerService,
		xiongNuService:      xiongNuService,
		pushService:         pushService,
		rankModule:          rankModule,
		closeNotify:         make(chan struct{}),
		loopExitNotify:      make(chan struct{}),
	}

	heromodule.RegisterHeroOnlineListener(m)

	// 检查匈奴的城存不存在
	m.xiongNuService.WalkInfo(func(info xiongnuface.ResistXiongNuInfo) {
		if info.BaseId() == 0 {
			m.xiongNuService.RemoveInfo(info)
			return
		}

		basePosInfo := m.realmService.GetBigMap().GetRoBase(info.BaseId())
		if basePosInfo == nil {
			logrus.Errorf("一上线，发现匈奴的主城找不到")
			m.xiongNuService.RemoveInfo(info)
			return
		}
	})

	go call.CatchPanic(m.loop, "匈奴刷新")

	return m
}

//gogen:iface
type XiongNuModule struct {
	dep                 iface.ServiceDep
	timeService         iface.TimeService         // 时间
	configDatas         *config.ConfigDatas       // 配置
	guildService        iface.GuildService        // 联盟数据
	realmService        iface.RealmService        // 场景数据
	heroSnapshotService iface.HeroSnapshotService // 玩家镜像
	mailModule          iface.MailModule          // 邮件服务
	worldService        iface.WorldService        // 世界服务
	tickerService       iface.TickerService       // 定时事件
	xiongNuService      iface.XiongNuService
	pushService         iface.PushService
	rankModule          iface.RankModule
	closeNotify         chan struct{}
	loopExitNotify      chan struct{}
}

func (m *XiongNuModule) OnHeroOnline(hc iface.HeroController) {
	// 发送消息
	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		return
	}

	msg := m.xiongNuService.XiongNuInfoMsg(guildId)
	if msg != nil {
		hc.Send(msg)
	}
}

// 启动一个定时任务，刷新
func (m *XiongNuModule) Close() {
	close(m.closeNotify)
	<-m.loopExitNotify
}

func (m *XiongNuModule) loop() {
	defer close(m.loopExitNotify)

	dailyTickTime := m.tickerService.GetDailyTickTime()
	saveTickTime := m.tickerService.GetPer10MinuteTickTime()
	secondTick := updateWheel.After(time.Second)

	for {
		select {
		case <-secondTick:
			secondTick = updateWheel.After(time.Second)
			call.CatchPanic(m.updatePerSecond, "每秒更新匈奴")
		case <-saveTickTime.Tick():
			saveTickTime = m.tickerService.GetPer10MinuteTickTime()
			call.CatchPanic(m.xiongNuService.Save, "保存匈奴数据")
		case <-m.closeNotify:
			call.CatchPanic(m.xiongNuService.Save, "保存匈奴数据")
			return
		case <-dailyTickTime.Tick():
			// 每日重置
			m.xiongNuService.ResetDaily(dailyTickTime.GetTickTime())
			dailyTickTime = m.tickerService.GetDailyTickTime()
		}
	}
}

// 每秒更新
func (m *XiongNuModule) updatePerSecond() {
	ctime := m.timeService.CurrentTime()

	m.xiongNuService.WalkInfo(func(info xiongnuface.ResistXiongNuInfo) {
		if info.BaseId() == 0 {
			return
		}

		ended := m.update(info, ctime)

		if ended {
			m.xiongNuService.RemoveInfo(info)
		} else if info.NeedSyncChange() {
			// 同步下数据
			snapshot := m.guildService.GetSnapshot(info.GuildId())
			if snapshot != nil {
				m.worldService.MultiSend(snapshot.UserMemberIds, xiongnu.NewS2cInfoBroadcastMsg(must.Marshal(info.EncodeClient())))
			}
		}
	})
}

// 加入联盟后的处理
func (m *XiongNuModule) JoinGuild(heroId, guildId int64) {
	if m.xiongNuService.IsTodayStarted(heroId) {
		return
	}

	info := m.xiongNuService.GetInfo(guildId)
	if info == nil {
		return
	}

	if i64.Contains(info.GivePrizeMembers(), heroId) {
		// 在里面了
		return
	}

	info.AddGivePrizeMember(heroId)
	m.xiongNuService.SetTodayStarted(heroId) // TODO 极端情况下会出现该玩家奖励没给到的情况

	m.worldService.SendFunc(heroId, func() pbutil.Buffer { return xiongnu.NewS2cInfoBroadcastMsg(must.Marshal(info.EncodeClient())) })
}

// 设置防守者
//gogen:iface c2s_set_defender
func (m *XiongNuModule) ProcessSetDefender(proto *xiongnu.C2SSetDefenderProto, hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.Debugf("没有加入联盟")
		hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_NO_GUILD)
		return
	}

	if m.xiongNuService.IsStarted(guildId) {
		logrus.Debugf("联盟当前已经开启了")
		hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_STARTED)
		return
	}

	targetId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("客户端发送的id非法")
		hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_TARGET_NOT_FOUND)
		return
	}

	// 检查是否有权限
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Errorf("没有找到联盟数据: %d", guildId)
			hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_NO_GUILD)
			return
		}

		member := g.GetMember(hc.Id())
		if member == nil {
			logrus.Debugf("玩家没有加入联盟")
			hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_NO_GUILD)
			return
		}

		if m.xiongNuService.IsStarted(guildId) {
			logrus.Debugf("联盟当前已经开启了")
			hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_STARTED)
			return
		}

		if !member.ClassLevelData().Permission.OpenResistXiongNu {
			logrus.Debugf("没有开启抗击匈奴的权限")
			hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_NO_PERMISISON)
			return
		}

		// 检查是不是有足够的防守的队伍
		toSet := g.GetMember(targetId)
		if toSet == nil {
			logrus.Debugf("没找到要设置防御的成员")
			hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_TARGET_NOT_MEMBER)
			return
		}

		exist := i64.Contains(g.ResistXiongNuDefenders(), targetId)
		if proto.ToSet {
			if uint64(len(g.ResistXiongNuDefenders())) > m.configDatas.ResistXiongNuMisc().DefenseMemberCount {
				logrus.Debugf("设置防御成员已满，请稍后再试")
				hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_FULL)
				return
			}

			// 上
			targetSnapshot := m.heroSnapshotService.Get(targetId)
			if targetSnapshot == nil {
				logrus.Errorf("没找到玩家镜像数据: %d", targetId)
				hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_SERVER_ERROR)
				return
			}

			if targetSnapshot.GuildId != guildId {
				logrus.Errorf("要设置防御的成员不是我们联盟的")
				hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_SERVER_ERROR)
				return
			}

			if targetSnapshot.BaseRegion == 0 || targetSnapshot.BaseLevel == 0 {
				logrus.Debugf("要设置防御的成员当前处于流亡状态")
				hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_HOME_NOT_ALIVE)
				return
			}

			if m.xiongNuService.IsTodayStarted(targetId) {
				logrus.Debugf("要设置防御的成员今天已经参加了")
				hc.Send(xiongnu.ERR_SET_DEFENDER_FAIL_TODAY_START)
				return
			}

			if !exist {
				g.AddResistXiongNuDefender(targetId)
				m.guildService.ClearSelfGuildMsgCache(g.Id()) // 清掉缓存

				m.dep.Tlog().TlogGuildFlowById(hc.Id(), operate_type.XiongNuSetDef.Id(), u64.FromInt64(g.Id()), g.LevelData().Level, u64.FromInt(g.MemberCount()))
			}
		} else {
			// 下
			if exist {
				g.RemoveResistXiongNuDefender(targetId)
				m.guildService.ClearSelfGuildMsgCache(g.Id()) // 清掉缓存
			}
		}

		g.SetChanged()

		// 发送协议，告诉他，设置成功了
		hc.Send(xiongnu.NewS2cSetDefenderMsg(proto.Id, proto.ToSet))

		m.guildService.ClearSelfGuildMsgCache(guildId)
	})
}

// gm 命令开启
func (m *XiongNuModule) GmStart(hc iface.HeroController, guildId int64) (suc bool) {
	if m.xiongNuService.GetRInfo(guildId) != nil {
		logrus.Debugf("匈奴已经开启了，还开启个P啊")
		return
	}

	var heroBaseX, heroBaseY int
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		heroBaseX, heroBaseY = hero.BaseX(), hero.BaseY()
		return
	})

	var memberIds []int64
	var info xiongnuface.ResistXiongNuInfo

	// 检查是否有权限
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Debugf("没有找到联盟数据: %d", guildId)
			return
		}

		if g.LevelData().Level < m.configDatas.ResistXiongNuMisc().OpenNeedGuildLevel {
			logrus.Debugf("联盟等级不足，无法开启: %d", guildId)
			return
		}

		if m.xiongNuService.IsStarted(guildId) {
			logrus.Debugf("联盟当前已经开启了")
			return
		}

		//if g.IsStartResistXiongNuToday() {
		//	logrus.Debugf("今天已经开启了，无法再开启")
		//	return
		//}

		memberIds = g.AllUserMemberIds()
		//if uint64(len(memberIds)) < m.configDatas.ResistXiongNuMisc().DefenseMemberCount {
		//	logrus.Debugf("人数不够")
		//	return
		//}

		for _, memberId := range memberIds {
			if i64.Contains(g.ResistXiongNuDefenders(), memberId) {
				continue
			}

			count := uint64(len(g.ResistXiongNuDefenders()))
			if count > 0 && count >= m.configDatas.ResistXiongNuMisc().DefenseMemberCount {
				break
			}

			g.AddResistXiongNuDefender(memberId)
		}

		//if uint64(len(g.ResistXiongNuDefenders())) < m.configDatas.ResistXiongNuMisc().DefenseMemberCount {
		//	logrus.Debugf("人数不够")
		//	return
		//}

		g.SetChanged()

		// 开启

		info = xiongnuinfo.NewResistXiongNuInfo(
			m.configDatas,
			guildId,
			g.UnlockResistXiongNuData(),
			g.ResistXiongNuDefenders(),
			m.timeService.CurrentTime())

		m.xiongNuService.AddInfo(info)

		g.SetIsStartResistXiongNuToday(true)
	})

	// 预备开启，没开启
	if info == nil {
		return
	}

	baseId, baseX, baseY := m.realmService.GetBigMap().AddXiongNuBase(info, heroBaseX, heroBaseY, m.configDatas.ResistXiongNuMisc().GetMinRange(), m.configDatas.ResistXiongNuMisc().GetMaxRange())
	if baseId == 0 {
		// 设置为今天没开启，必须处理
		m.guildService.Func(func(guilds sharedguilddata.Guilds) {
			g := guilds.Get(guildId)
			if g == nil {
				logrus.Debugf("处理失败，但是想要处理回来，竟然没找到联盟了")
				return
			}

			g.SetChanged()
			g.SetIsStartResistXiongNuToday(false)
		})

		m.xiongNuService.RemoveInfo(info)

		return
	}

	for _, id := range memberIds {
		info.AddGivePrizeMember(id)
	}

	logrus.Debugf("GM开启抗击匈奴成功: %v", info)

	// 成功了
	info.SetBase(baseId, baseX, baseY)
	hc.Send(xiongnu.NewS2cStartMsg(i64.Int32(baseId), baseX, baseY))

	if data := m.configDatas.GuildLogHelp().StartXiongNu; data != nil {
		proto := data.NewLogProto(m.timeService.CurrentTime())
		proto.Text = data.Text.New().WithHeroName("GM命令开启").JsonString()
		m.guildService.AddLog(info.GuildId(), proto)
	}

	m.worldService.MultiSendMsgs(memberIds, []pbutil.Buffer{
		xiongnu.NewS2cBroadcastStartMsg("GM命令开启"),
		xiongnu.NewS2cInfoBroadcastMsg(must.Marshal(info.EncodeClient())),
	})

	m.pushService.MultiPush(shared_proto.SettingType_ST_GUILD_XIONG_NU, memberIds, 0)

	return true
}

// 开启
//gogen:iface c2s_start
func (m *XiongNuModule) ProcessStart(proto *xiongnu.C2SStartProto, hc iface.HeroController) {
	ctime := m.timeService.CurrentTime()
	if ctime.Before(m.dep.SvrConf().GetServerStartTime().Add(m.configDatas.ResistXiongNuMisc().StartAfterServerOpen)) {
		logrus.Debugf("没有到开服匈奴开启时间")
		hc.Send(xiongnu.ERR_START_FAIL_START_TIME_LIMIT)
		return
	}

	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.Debugf("没有加入联盟")
		hc.Send(xiongnu.ERR_START_FAIL_NO_GUILD)
		return
	}

	if m.xiongNuService.IsStarted(guildId) {
		logrus.Debugf("联盟当前已经开启了")
		hc.Send(xiongnu.ERR_START_FAIL_STARTED)
		return
	}

	data := m.configDatas.GetResistXiongNuData(u64.FromInt32(proto.Level))
	if data == nil {
		logrus.Debugf("等级数据没找到")
		hc.Send(xiongnu.ERR_START_FAIL_INVALID_LEVEL)
		return
	}

	snapshot := m.heroSnapshotService.Get(hc.Id())
	if snapshot == nil {
		logrus.Errorf("没找到玩家镜像数据")
		hc.Send(xiongnu.ERR_START_FAIL_SERVER_ERROR)
		return
	}

	if snapshot.BaseLevel == 0 {
		logrus.Debugf("当前处于流亡状态")
		hc.Send(xiongnu.ERR_START_FAIL_BASE_DEAD)
		return
	}

	name := snapshot.Name

	var memberIds []int64
	var info xiongnuface.ResistXiongNuInfo

	var glevel uint64
	// 检查是否有权限
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Debugf("没有找到联盟数据: %d", guildId)
			hc.Send(xiongnu.ERR_START_FAIL_NO_GUILD)
			return
		}

		if g.LevelData().Level < m.configDatas.ResistXiongNuMisc().OpenNeedGuildLevel {
			logrus.Debugf("联盟等级不足，无法开启: %d", guildId)
			hc.Send(xiongnu.ERR_START_FAIL_LEVEL_NOT_ENOUGH)
			return
		}

		member := g.GetMember(hc.Id())
		if member == nil {
			logrus.Debugf("玩家没有加入联盟")
			hc.Send(xiongnu.ERR_START_FAIL_NO_GUILD)
			return
		}

		if m.xiongNuService.IsStarted(guildId) {
			logrus.Debugf("联盟当前已经开启了")
			hc.Send(xiongnu.ERR_START_FAIL_STARTED)
			return
		}

		if !member.ClassLevelData().Permission.OpenResistXiongNu {
			logrus.Debugf("没有开启抗击匈奴的权限")
			hc.Send(xiongnu.ERR_START_FAIL_NO_PERMISISON)
			return
		}

		if g.UnlockResistXiongNuData().Level < data.Level {
			logrus.Debugf("未解锁，无法开启")
			hc.Send(xiongnu.ERR_START_FAIL_LOCK_LEVEL)
			return
		}

		if g.IsStartResistXiongNuToday() {
			logrus.Debugf("今天已经开启了，无法再开启")
			hc.Send(xiongnu.ERR_START_FAIL_STARTED)
			return
		}

		if uint64(len(g.ResistXiongNuDefenders())) < m.configDatas.ResistXiongNuMisc().DefenseMemberCount {
			logrus.Debugf("未解锁，防守人数不够")
			hc.Send(xiongnu.ERR_START_FAIL_NOT_ENOUGH)
			return
		}

		g.SetChanged()

		// 开启

		memberIds = g.AllUserMemberIds()

		info = xiongnuinfo.NewResistXiongNuInfo(
			m.configDatas,
			guildId,
			data,
			g.ResistXiongNuDefenders(),
			m.timeService.CurrentTime())

		m.xiongNuService.AddInfo(info)

		g.SetIsStartResistXiongNuToday(true)

		glevel = g.LevelData().Level
	})

	// 预备开启，没开启
	if info == nil {
		return
	}

	baseId, baseX, baseY := m.realmService.GetBigMap().AddXiongNuBase(info, snapshot.BaseX, snapshot.BaseY, m.configDatas.ResistXiongNuMisc().GetMinRange(), m.configDatas.ResistXiongNuMisc().GetMaxRange())
	if baseId == 0 {
		// 设置为今天没开启，必须处理
		m.guildService.Func(func(guilds sharedguilddata.Guilds) {
			g := guilds.Get(guildId)
			if g == nil {
				logrus.Debugf("处理失败，但是想要处理回来，竟然没找到联盟了")
				return
			}

			g.SetChanged()
			g.SetIsStartResistXiongNuToday(false)
		})

		m.xiongNuService.RemoveInfo(info)

		// 发送服务器繁忙，请稍后再试
		hc.Send(xiongnu.ERR_START_FAIL_SERVER_ERROR)

		return
	}

	m.dep.Tlog().TlogGuildFlowById(hc.Id(), operate_type.XiongNuStart.Id(), u64.FromInt64(guildId), glevel, u64.FromInt(len(memberIds)))

	todayJoinMap := m.xiongNuService.TodayJoinMap()
	for _, id := range memberIds {
		_, setSuc := todayJoinMap.SetIfAbsent(id, id)
		if setSuc {
			info.AddGivePrizeMember(id)

			// 参与匈奴成就任务
			m.dep.HeroData().FuncWithSend(id, func(hero *entity.Hero, result herolock.LockResult) {
				hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumXiongNuStart)
				hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_AccumXiongNuStart, data.Level)
				heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_START_XIONGNU)

				m.dep.Tlog().TlogGuildFlow(hero, operate_type.XiongNuJoin.Id(), u64.FromInt64(guildId), glevel, u64.FromInt(len(memberIds)))
			})
		}

	}

	// 成功了
	info.SetBase(baseId, baseX, baseY)

	hc.Send(xiongnu.NewS2cStartMsg(i64.Int32(baseId), baseX, baseY))

	if data := m.configDatas.GuildLogHelp().StartXiongNu; data != nil {
		proto := data.NewHeroLogProto(ctime, snapshot.IdBytes, snapshot.Head)
		proto.Text = data.Text.New().WithHeroName(snapshot.Name).JsonString()
		m.guildService.AddLog(info.GuildId(), proto)
	}

	// 广播
	m.worldService.MultiSendMsgs(memberIds, []pbutil.Buffer{
		xiongnu.NewS2cBroadcastStartMsg(name),
		xiongnu.NewS2cInfoBroadcastMsg(must.Marshal(info.EncodeClient())),
	})

	m.pushService.MultiPush(shared_proto.SettingType_ST_GUILD_XIONG_NU, memberIds, 0)
}

// 开启
//gogen:iface c2s_troop_info
func (m *XiongNuModule) ProcessTroopInfo(hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.Debugf("没有加入联盟")
		hc.Send(xiongnu.ERR_TROOP_INFO_FAIL_NO_GUILD)
		return
	}

	info := m.xiongNuService.GetInfo(guildId)
	if info == nil {
		logrus.Debugf("联盟当前没有开启")
		hc.Send(xiongnu.ERR_TROOP_INFO_FAIL_NOT_STARTED)
		return
	}

	if info.BaseId() == 0 {
		logrus.Debugf("联盟当前没有开启")
		hc.Send(xiongnu.ERR_TROOP_INFO_FAIL_NOT_STARTED)
		return
	}

	proto := m.realmService.GetBigMap().GetXiongNuTroopInfo(info.BaseId(), info.GuildId())
	if proto == nil {
		logrus.Debugf("没有找到匈奴队伍信息")
		hc.Send(xiongnu.ERR_TROOP_INFO_FAIL_NOT_STARTED)
		return
	}

	hc.Send(xiongnu.NewS2cTroopInfoMsg(must.Marshal(proto), u64.Int32(info.Morale())))
}

//gogen:iface
func (m *XiongNuModule) ProcessGetXiongNuNpcBaseInfo(proto *xiongnu.C2SGetXiongNuNpcBaseInfoProto, hc iface.HeroController) {
	info := m.xiongNuService.GetInfo(int64(proto.GuildId))
	if info == nil {
		logrus.Debugf("联盟没有开启抗击匈奴")
		hc.Send(xiongnu.ERR_GET_XIONG_NU_NPC_BASE_INFO_FAIL_NOT_STARTED)
		return
	}

	suc, fightingAmount := m.realmService.GetBigMap().GetMaxXiongNuTroopFightingAmount(info.BaseId(), info.GuildId())
	if !suc {
		logrus.Debugf("取匈奴最高战斗力失败")
		hc.Send(xiongnu.ERR_GET_XIONG_NU_NPC_BASE_INFO_FAIL_NOT_STARTED)
		return
	}

	var guildName, guildFlag string
	snapshot := m.guildService.GetSnapshot(info.GuildId())
	if snapshot != nil {
		guildName, guildFlag = snapshot.Name, snapshot.FlagName
	}

	hc.Send(xiongnu.NewS2cGetXiongNuNpcBaseInfoMsg(proto.GuildId, guildName, guildFlag, u64.Int32(info.Morale()), timeutil.Marshal32(info.StartTime()), u64.Int32(fightingAmount)))
}

func (m *XiongNuModule) update(info xiongnuface.ResistXiongNuInfo, ctime time.Time) (ended bool) {
	m.updateWave(info, ctime)
	return m.updateEnd(info, ctime)
}

func (m *XiongNuModule) updateWave(info xiongnuface.ResistXiongNuInfo, ctime time.Time) {

	if info.Wave() > info.Data().MaxWave() {
		// 已经刷完所有的怪
		return
	}

	var curWaveData *xiongnu2.ResistXiongNuWaveData
	if info.Wave() == 0 && !ctime.Before(info.NextWaveTime()) {
		curWaveData = info.Data().FirstResistWave
		info.IncWave()
	} else {
		curWaveData = info.Data().WaveData(info.Wave())
		if curWaveData == nil {
			logrus.Errorf("没找到当前回合的数据: %d", info.Wave(), info.Data().Name)
			return
		}

		if uint64(len(curWaveData.Monsters)) <= info.AddMonsterCount() {
			// 检查是不是到了下一波了
			if ctime.Before(info.NextWaveTime()) {
				// 下一波了
				return
			}

			curWaveData = curWaveData.Next

			info.IncWave()
			if curWaveData == nil {
				// 刷怪结束了
				logrus.Debugf("刷怪结束了")
				return
			}
		} else if ctime.Before(info.NextRefreshMonsterTime()) {
			// 还没到下次刷新时间
			return
		}
	}

	// 刷怪
	defenders := info.Defenders()
	if len(defenders) <= 0 {
		info.SetAddMonsterCount(uint64(len(curWaveData.Monsters)))
		logrus.Errorf("防守人员数量竟然为0")
		return
	}

	var validDefenders []int64

	for _, defenderId := range defenders {
		defender := m.heroSnapshotService.Get(defenderId)
		if defender == nil {
			continue
		}

		if defender.GuildId != info.GuildId() {
			continue
		}

		if defender.BaseRegion == 0 {
			continue
		}

		if defender.BaseLevel == 0 {
			continue
		}

		validDefenders = append(validDefenders, defenderId)
	}

	if len(validDefenders) <= 0 {
		info.SetAddMonsterCount(uint64(len(curWaveData.Monsters)))
		logrus.Debugf("没有防守英雄")
		return
	}

	toAddTroops := curWaveData.Monsters[info.AddMonsterCount():]
	if waveLeftMonsterCount := len(toAddTroops); len(validDefenders) > waveLeftMonsterCount {
		// 防守人数超过怪物数量，随机下
		i64.Shuffle(validDefenders)
		validDefenders = validDefenders[:len(toAddTroops)]
	} else {
		// 防守人数小于怪物数量
		toAddTroops = toAddTroops[:len(validDefenders)]
	}

	logrus.Debugf("可以被攻击的防守城市: %v, 所有防守城市: %v, 怪物数量: %d", validDefenders, info.Defenders(), len(toAddTroops))

	processed := m.realmService.GetBigMap().AddXiongNuTroop(info.BaseId(), validDefenders, toAddTroops) // 添加匈奴队伍
	if !processed {
		// 处理失败，下一秒再来
		return
	}

	info.IncAddMonsterCount(uint64(len(toAddTroops)))
	info.SetNextRefreshMonsterTime(ctime.Add(m.configDatas.ResistXiongNuMisc().AttackDuration)) // 更新下次刷新时间
}

func (m *XiongNuModule) updateEnd(info xiongnuface.ResistXiongNuInfo, ctime time.Time) (ended bool) {
	if ctime.Before(info.EndTime()) {
		// 检查主城是否被打爆了
		if ctime.Before(info.ResistTime()) {
			// 反击前
			return
		}
		// 设置resist并且发送联盟公告
		if !info.IsResist() {
			info.SetResist()
			m.dep.Chat().SysChat(0, info.GuildId(), shared_proto.ChatType_ChatGuild, m.configDatas.TextHelp().XiongNuResistBroadcast.New().JsonString(), shared_proto.ChatMsgType_ChatMsgText, true, true, true, false)
		}

		basePosInfo := m.realmService.GetBigMap().GetRoBase(info.BaseId())
		if basePosInfo != nil && basePosInfo.GetBaseType() == realmface.BaseTypeNpc && npcid.IsXiongNuNpcId(basePosInfo.GetId()) {
			return
		}

		info.SetDefeated()

		// 消失了
	} else {
		// 干掉主城
		processed, _, _, _ := m.realmService.GetBigMap().RemoveBase(info.BaseId(), false, nil, nil)
		if !processed {
			return // 再来一次
		}
	}

	// 设置联盟信息
	m.guildService.Func(func(guilds sharedguilddata.Guilds) {
		g := guilds.Get(info.GuildId())
		if g == nil {
			return
		}

		m.onEnd(info, g, ctime)
	})

	logrus.Debugf("%d 抗击匈奴结束了", info.GuildId())

	return true
}

func (m *XiongNuModule) onEnd(info xiongnuface.ResistXiongNuInfo, g *sharedguilddata.Guild, ctime time.Time) {
	g.SetChanged()

	// 设置联盟状态，算次数什么的
	// 摧毁了，联盟广播结束什么的

	scoreData := m.configDatas.ResistXiongNuScoreData().MinKeyData

	scorePrestige := uint64(0)
	prize := info.Data().ScorePlunderPrizes[0].GetGuildPrize(g.LevelData().Level)

	for idx, data := range m.configDatas.GetResistXiongNuScoreDataArray() {
		if info.WipeOutMonsterCount() >= data.WipeOutInvadeMonsterCount {
			scoreData = data
			scorePrestige = info.Data().ScorePrestiges[idx]
			prize = info.Data().ScorePlunderPrizes[idx].GetGuildPrize(g.LevelData().Level)
		}
	}

	unlockNextLevel := false
	// 检查是不是升级了
	// 处理
	if scoreData.UnlockNextLevel && info.Data().NextLevel != nil {
		unlockNextLevel = true
		g.SetUnlockResistXiongNuData(info.Data().NextLevel)

		if data := m.configDatas.GuildLogHelp().UnlockXiongNu; data != nil {
			proto := data.NewLogProto(ctime)
			proto.Text = data.Text.New().WithResistXiongNuName(info.Data().Name).WithUnlockXiongNuName(info.Data().NextLevel.Name).JsonString()
			m.guildService.AddLog(info.GuildId(), proto)
		}
	}

	addPrestige := scorePrestige
	if info.IsDefeated() {
		// 打败了
		addPrestige = u64.Plus(addPrestige, info.Data().ResistSucPrestige)
	}

	if data := m.configDatas.GuildLogHelp().ResistXiongNuAddPrestige; data != nil {
		proto := data.NewLogProto(ctime)
		proto.Text = data.Text.New().WithResistXiongNuName(info.Data().Name).WithScoreName(scoreData.Name).WithPrestige(addPrestige).JsonString()
		m.guildService.AddLog(info.GuildId(), proto)
	}
	allMemberIds := g.AllUserMemberIds()
	if addPrestige > 0 {
		g.AddPrestige(addPrestige)
		m.dep.Country().AddPrestige(g.Country().Id, addPrestige)
		// 更新联盟任务进度
		if g.LevelData().Level >= m.configDatas.GuildGenConfig().TaskOpenLevel {
			if data := m.configDatas.GetGuildTaskData(u64.FromInt32(int32(server_proto.GuildTaskType_XiongNv))); data != nil {
				if g.AddGuildTaskProgress(data, addPrestige) {
					m.worldService.MultiSend(allMemberIds, guild.NewS2cNoticeTaskStageUpdateMsg(int32(data.TaskType), int32(g.GetGuildTaskStageIndex(data.TaskType))))
				}
			}
		}
	}

	m.rankModule.AddOrUpdateRankObj(ranklist.NewGuildRankObj(
		m.guildService.GetSnapshot, m.heroSnapshotService.Get, g))

	if mailData := m.configDatas.MailHelp().XiongNuScore; mailData != nil && prize != nil {
		scoreMailProto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
		scoreMailProto.Prize = prize.Encode()
		scoreMailProto.Text = mailData.NewTextFields().WithWipeOutCount(info.WipeOutMonsterCount()).WithScoreName(scoreData.Name).WithPrestige(scorePrestige).JsonString()

		for _, id := range info.GivePrizeMembers() {
			m.mailModule.SendProtoMail(id, scoreMailProto, ctime)
		}
	}

	if mailData := m.configDatas.MailHelp().XiongNuResistSuc; info.IsDefeated() && mailData != nil {
		scoreMailProto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
		scoreMailProto.Prize = info.Data().ResistSucPrize.Encode()
		scoreMailProto.Text = mailData.NewTextFields().WithPrestige(info.Data().ResistSucPrestige).JsonString()

		for _, id := range info.GivePrizeMembers() {
			m.mailModule.SendProtoMail(id, scoreMailProto, ctime)
		}
	}

	// 设置最新的一场战斗的记录
	lastResist := info.EncodeLast(scoreData.Level, m.heroSnapshotService.Get)
	g.SetLastResistXiongNuProto(lastResist)

	fightProto := info.EncodeFight(m.heroSnapshotService.Get)
	g.SetLastResistXiongNuFightProto(fightProto)

	toBroadcast := xiongnu.NewS2cEndBroadcastMsg(i64.Int32(g.Id()), must.Marshal(lastResist), unlockNextLevel, u64.Int32(addPrestige))
	m.worldService.MultiSend(allMemberIds, toBroadcast)

	hctx := heromodule.NewContext(m.dep, operate_type.XiongNuSucc)
	if d := hctx.BroadcastHelp().XiongNuSucc; d != nil {
		hctx.AddGuildBroadcast(d, g.Id(), info.Data().Level, scoreData.Level, func() *i18n.Fields {
			gName, ok := hctx.GetFlagGuildNameFromSnapshot(g.Id())
			if !ok {
				return nil
			}
			text := d.NewTextFields()
			text.WithClickGuildFields(data.KeyGuild, gName, g.Id())
			text.WithFields(data.KeyText, info.Data().Name).WithFields(data.KeyScore, scoreData.Name)
			return text
		})

	}
}

var emptyDefenserFightAmountMsg = xiongnu.NewS2cGetDefenserFightAmountMsg(0, nil, nil, nil).Static()

//gogen:iface
func (m *XiongNuModule) ProcessGetDefenserFightAmount(proto *xiongnu.C2SGetDefenserFightAmountProto, hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		hc.Send(emptyDefenserFightAmountMsg)
		return
	}

	info := m.xiongNuService.GetInfo(guildId)
	if info == nil {
		logrus.Debugf("请求防守战力，联盟没有开启抗击匈奴")
		hc.Send(emptyDefenserFightAmountMsg)
		return
	}

	defensers := info.Defenders()
	if len(defensers) <= 0 {
		hc.Send(emptyDefenserFightAmountMsg)
		return
	}

	getEnemyCount := m.realmService.GetBigMap().GetXiongNuInvateTargetCount(info.BaseId())
	if getEnemyCount == nil {
		getEnemyCount = i64.EmptyGetU64()
	}

	ctime := m.timeService.CurrentTime()

	dfa := info.GetDefenserFightAmount()
	if dfa == nil || dfa.ExpireTime.Before(ctime) {
		// 不在有效期，先更新数据

		fightAmounts := make([]uint64, len(defensers))
		enemyCounts := make([]uint64, len(defensers))
		for i, defId := range defensers {
			if def := m.heroSnapshotService.Get(defId); def != nil {
				fightAmounts[i] = def.FightAmount
			}
			enemyCounts[i] = getEnemyCount(defId)
		}

		expireTime := ctime.Add(2 * time.Second)
		dfa = info.UpdateDefenserFightAmount(defensers, fightAmounts, enemyCounts, expireTime, dfa)
	}

	// 版本一致，并且在有效期内，直接发送数据
	if dfa.Version == proto.Version {
		hc.Send(dfa.SameVersionMsg)
	} else {
		hc.Send(dfa.DiffVersionMsg)
	}
}

var emptyXiongNuFightInfoMsg = xiongnu.NewS2cGetXiongNuFightInfoMsg(nil).Static()

//gogen:iface c2s_get_xiong_nu_fight_info
func (m *XiongNuModule) ProcessGetXiongNuFightInfo(hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.Debug("请求匈奴战斗排行榜，你没有联盟")
		hc.Send(xiongnu.ERR_GET_XIONG_NU_FIGHT_INFO_FAIL_NOT_IN_GUILD)
		return
	}

	var toSend pbutil.Buffer
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		toSend = g.GetLastResistXiongNuFightMsg()
	})

	if toSend == nil {
		toSend = emptyXiongNuFightInfoMsg
	}

	hc.Send(toSend)
}
