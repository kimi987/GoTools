package mingc_war

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/chat"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/event"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timer"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"runtime/debug"
	"time"
)

//gogen:iface
type MingcWarService struct {
	dep iface.ServiceDep

	queue *event.EventQueue

	mingcWar *MingcWar

	closeNotify    chan struct{}
	loopExitNotify chan struct{}

	viewMsgVer   *atomic.Uint64
	viewMsg      pbutil.Buffer
	viewEmptyMsg pbutil.Buffer
}

func NewMingcWarService(dep iface.ServiceDep) *MingcWarService {
	m := &MingcWarService{
		dep:            dep,
		closeNotify:    make(chan struct{}),
		loopExitNotify: make(chan struct{}),
	}

	m.queue = event.NewEventQueue(2048, 5*time.Second, "MingcWarEvent")

	ctime := m.dep.Time().CurrentTime()
	m.load()

	m.viewMsgVer = atomic.NewUint64(1)
	m.viewEmptyMsg = mingc_war.NewS2cViewMcWarMsg(0, nil).Static()

	warFirstStartTime := m.dep.SvrConf().GetServerStartTime().Add(m.dep.Datas().MingcMiscData().StartAfterServerOpen)
	if m.mingcWar == nil {
		nilNewStartTime := timeutil.Max(warFirstStartTime, ctime.Add(-time.Hour*24*7))
		m.startNewWar(nilNewStartTime)
	} else if m.mingcWar.state == shared_proto.MingcWarState_MC_T_NOT_START {
		// 开服时间改了
		if warFirstStartTime.After(m.mingcWar.startTime) {
			m.startNewWar(ctime)
		}
	} else {
		m.updateMsg()
	}

	go call.CatchLoopPanic(m.loop, "MingcWarService.Loop")

	return m
}

func (m *MingcWarService) CurrMcWarStage() (state int32, start, end time.Time) {
	state = int32(m.mingcWar.state)
	start = m.mingcWar.currStage().stageStartTime()
	end = m.mingcWar.currStage().stageEndTime()
	return
}

func (m *MingcWarService) McWarStartEndTime() (start, end time.Time) {
	m.doFunc("McWarStartEndTime", func() {
		start = m.mingcWar.startTime
		end = m.mingcWar.endTime
	})

	return
}

func (m *MingcWarService) ViewMsg(ver uint64) pbutil.Buffer {
	if ver != 0 && ver == m.viewMsgVer.Load() {
		return m.viewEmptyMsg
	}
	if m.viewMsg == nil {
		m.updateMsg()
	}
	return m.viewMsg
}

func (m *MingcWarService) ViewSelfGuildProto(gid int64) *shared_proto.McWarGuildProto {
	var p *shared_proto.McWarGuildProto
	m.doFunc("ViewSelfGuildProto", func() {
		p = m.mingcWar.encodeGuild(gid, m.dep.GuildSnapshot().GetSnapshot, m.dep.HeroSnapshot().GetBasicSnapshotProto)
	})

	return p
}

func (m *MingcWarService) ViewMcWarSceneMsg(mcId uint64) (succMsg, errMsg pbutil.Buffer) {
	if mc := m.mingcWar.mingcs[mcId]; mc != nil && mc.scene != nil {
		if mc.scene.viewMsg == nil {
			mc.scene.updateViewMsg()
		}
		succMsg = mc.scene.viewMsg
	} else {
		errMsg = mingc_war.ERR_VIEW_MC_WAR_SCENE_FAIL_INVALID_MCID
	}
	return
}

func (m *MingcWarService) ViewMcWarMcMsg(mcId uint64) (succMsg, errMsg pbutil.Buffer) {
	if mc := m.mingcWar.mingcs[mcId]; mc != nil {
		succMsg = mingc_war.NewS2cViewMingcWarMcMsg(mc.encode(m.dep.GuildSnapshot().GetSnapshot))
	} else {
		errMsg = mingc_war.ERR_VIEW_MINGC_WAR_MC_FAIL_INVALID_MCID
	}
	return
}
func (m *MingcWarService) UpdateMsg() {
	m.updateMsg()
}

func (m *MingcWarService) updateMsg() {
	mcWarProto := m.mingcWar.encode(m.dep.GuildSnapshot().GetSnapshot)
	m.viewMsg = mingc_war.NewS2cViewMcWarMsg(u64.Int32(m.viewMsgVer.Inc()), mcWarProto).Static()
	m.dep.Country().UpdateMcWarMsg(mcWarProto)
}

func (m *MingcWarService) loop() {
	defer close(m.loopExitNotify)

	loopWheel := timer.NewTimingWheel(mingcdata.McWarLoopDuration, 32)
	secondTick := loopWheel.After(time.Second)
	wheel := timer.NewTimingWheel(30*time.Second, 32)
	minuteTick := wheel.After(time.Minute)

	for {
		select {
		case <-secondTick:
			secondTick = loopWheel.After(time.Second)
			m.doFunc("mingc.loop", m.update)
		case <-minuteTick:
			minuteTick = wheel.After(time.Minute)
			m.doFunc("mingc.refresh", m.refresh)
		case <-m.closeNotify:
			return
		}
	}
}

func (m *MingcWarService) load() {

	var bytes []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		bytes, err = m.dep.Db().LoadKey(ctx, server_proto.Key_MingcWar)
		return
	})
	if err != nil {
		logrus.WithError(err).Panic("加载名城战数据失败")
	}

	if len(bytes) <= 0 {
		return
	}

	proto := &server_proto.MingcWarServerProto{}
	if err := proto.Unmarshal(bytes); err != nil {
		logrus.WithError(err).Panic("解析名城战数据失败")
	}

	m.mingcWar = NewDefaultMingcWar(m.dep)
	m.mingcWar.unmarshal(proto, m.dep)

	logrus.Debugf("加载名城战数据成功")
}

func (m *MingcWarService) refresh() {
	if m.mingcWar.state == shared_proto.MingcWarState_MC_T_FIGHT {
		for _, mingc := range m.mingcWar.mingcs {
			if mingc.scene != nil && mingc.scene.fightState == shared_proto.MingcWarFightState_MC_F_FIGHT_RUNNING {
				mingc.scene.sort4TroopRank()
			}
		}
	}
}

func (m *MingcWarService) update() {
	ctime := m.dep.Time().CurrentTime()

	if ctime.Before(m.mingcWar.startTime) {
		return
	}

	if ctime.After(m.mingcWar.endTime) {
		m.mingcWar.currStage().onEnd(m.mingcWar, m.dep)
		m.onEnd()
		m.startNewWar(m.mingcWar.endTime.Add(-time.Minute)) // 前移1分钟，保证ctime < 下场开始时间
		m.updateMsg()
		return
	}

	if changed, oldState := m.mingcWar.changeState(ctime); changed {
		oldStage := m.mingcWar.stage(oldState)
		if oldStage != nil {
			oldStage.onEnd(m.mingcWar, m.dep)
		}
		newStage := m.mingcWar.currStage()
		if newStage != nil {
			newStage.onStart(m.mingcWar, m.dep)
		}
		m.updateMsg()
	}

	currStage := m.mingcWar.currStage()
	if currStage.onUpdate(m.mingcWar, m.dep) {
		m.updateMsg()
	}
}

func (m *MingcWarService) Close() {
	close(m.closeNotify)
	<-m.loopExitNotify

	m.queue.Stop()

	if err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		return m.dep.Db().SaveKey(ctx, server_proto.Key_MingcWar, must.Marshal(m.mingcWar.encodeServer()))
	}); err != nil {
		logrus.WithError(err).Error("保存名城战数据失败")
		return
	}
	logrus.Debugf("保存名城战数据成功")
}

func (m *MingcWarService) doFunc(handlerName string, f func()) (succ bool) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", r).Errorf("%s recovered from panic. SEVERE!!!", handlerName)
			metrics.IncPanic()
		}
	}()

	return m.queue.TimeoutFunc(true, f)
}

func (m *MingcWarService) newMingcWar(ctime time.Time) (w *MingcWar) {
	warFirstStartTime := m.dep.SvrConf().GetServerStartTime().Add(m.dep.Datas().MingcMiscData().StartAfterServerOpen)
	if ctime.Before(warFirstStartTime) {
		ctime = warFirstStartTime
	}

	var timeData *mingcdata.MingcTimeData
	var start time.Time
	for _, d := range m.dep.Datas().GetMingcTimeDataArray() {
		nextStart, _ := d.WarTime.NextTime(ctime)
		if timeutil.IsZero(start) || nextStart.Before(start) {
			start, timeData = nextStart, d
		}
	}

	return newMingcWar(ctime, m.dep, timeData)
}

func (m *MingcWarService) onEnd() {}

func (m *MingcWarService) startNewWar(ctime time.Time) {
	m.mingcWar = m.newMingcWar(ctime)
	m.mingcWar.currStage().onStart(m.mingcWar, m.dep)
	logrus.Debugf("新开名城战：start:%v end:%v", m.mingcWar.startTime, m.mingcWar.endTime)
	m.updateMsg()
}

func (m *MingcWarService) ApplyAtk(gid int64, mcData *mingcdata.MingcBaseData, cost uint64) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("ApplyAtk", func() {
		ctime := m.dep.Time().CurrentTime()
		if errMsg = applyAttack(m.mingcWar, gid, mcData, cost, ctime); errMsg != nil {
			return
		}
		succMsg = mingc_war.NewS2cApplyAtkMsg(u64.Int32(mcData.Id), u64.Int32(cost))
		m.updateMsg()
	}) {
		errMsg = mingc_war.ERR_APPLY_ATK_FAIL_SERVER_ERR
	}

	return
}

func (m *MingcWarService) ApplyAst(gid int64, isAtk bool, mcData *mingcdata.MingcBaseData) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("ApplyAst", func() {
		var recvGId int64
		if recvGId, errMsg = applyAssist(m.mingcWar, gid, isAtk, mcData, m.dep.Datas().MingcMiscData()); errMsg != nil {
			return
		}
		succMsg = mingc_war.NewS2cApplyAstMsg(u64.Int32(mcData.Id), isAtk)
		m.updateMsg()

		if g := m.dep.GuildSnapshot().GetSnapshot(recvGId); g != nil {
			m.dep.World().Send(g.LeaderId, mingc_war.NewS2cReceiveApplyAstMsg(u64.Int32(mcData.Id)))
			m.dep.World().Send(g.LeaderId, mingc_war.RED_POINT_NOTICE_S2C)
		}
	}) {
		errMsg = mingc_war.ERR_APPLY_ATK_FAIL_SERVER_ERR
	}

	return
}

func (m *MingcWarService) CancelApplyAst(gid int64, mcData *mingcdata.MingcBaseData) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("CancelApplyAst", func() {
		var recvGId int64
		if recvGId, errMsg = cancelApplyAst(m.mingcWar, gid, mcData); errMsg != nil {
			return
		}
		succMsg = mingc_war.NewS2cCancelApplyAstMsg(u64.Int32(mcData.Id))
		m.updateMsg()

		if g := m.dep.GuildSnapshot().GetSnapshot(recvGId); g != nil {
			m.dep.World().Send(g.LeaderId, mingc_war.NewS2cReceiveCancelApplyAstMsg(u64.Int32(mcData.Id)))
		}
	}) {
		errMsg = mingc_war.ERR_APPLY_ATK_FAIL_SERVER_ERR
	}

	return
}

func (m *MingcWarService) ReplyApplyAst(operId, gid int64, mcData *mingcdata.MingcBaseData, agree bool) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("replyApplyAst", func() {
		if errMsg = replyApplyAst(m.mingcWar, operId, gid, mcData, agree, m.dep.Datas().MingcMiscData()); errMsg != nil {
			return
		}
		succMsg = mingc_war.NewS2cReplyApplyAstMsg(u64.Int32(mcData.Id), i64.Int32(gid), agree)
		m.updateMsg()

		if g := m.dep.GuildSnapshot().GetSnapshot(gid); g != nil {
			m.dep.World().Send(g.LeaderId, mingc_war.NewS2cApplyAstPassMsg(u64.Int32(mcData.Id), agree))
		}
	}) {
		errMsg = mingc_war.ERR_APPLY_ATK_FAIL_SERVER_ERR
	}

	return
}

func (m *MingcWarService) Watch(hc iface.HeroController, mcId uint64) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("Watch", func() {
		succMsg, errMsg = watch(m.mingcWar, mcId, hc.Id())
	}) {
		errMsg = mingc_war.ERR_WATCH_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) QuitWatch(hc iface.HeroController, mcId uint64) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("QuitWatch", func() {
		succMsg, errMsg = quitWatch(m.mingcWar, mcId, hc.Id())
	}) {
		errMsg = mingc_war.ERR_QUIT_WATCH_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) JoinFight(hc iface.HeroController, mcId uint64, capIds []uint64, xIndexs []int32) (succMsg, errMsg pbutil.Buffer) {
	var gid int64
	var heroProto *shared_proto.HeroBasicProto
	var caps []*shared_proto.CaptainInfoProto
	if !hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		gid = hero.GuildId()
		heroProto = hero.EncodeBasicProto(m.dep.GuildSnapshot().GetSnapshot)

		var capExisted bool
		for i, cid := range capIds {
			if cid <= 0 {
				caps = append(caps, nil)
				continue
			}

			if c := hero.Military().Captain(cid); c != nil {
				var xIndex int32
				if i < len(xIndexs) {
					xIndex = xIndexs[i]
				}

				caps = append(caps, c.EncodeInvaseCaptainInfo(hero, true, xIndex))
				capExisted = true
			} else {
				errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_INVALID_CAPTAIN_ID
				return
			}
		}

		if !capExisted {
			logrus.Warnf("名城战，加入战斗，一个武将都没有。heroId:%v", hero.Id())
			errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_INVALID_CAPTAIN_ID
			return
		}
	}) {
		errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_SERVER_ERR
		return
	}

	var bornPos cb.Cube
	if mc := m.mingcWar.mingcs[mcId]; mc != nil {
		if mc.isAtk(gid) {
			bornPos = m.dep.Datas().MingcWarSceneData().Get(mcId).AtkRelivePos
		} else if mc.isDef(gid) {
			bornPos = m.dep.Datas().MingcWarSceneData().Get(mcId).DefRelivePos
		} else {
			errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_NOT_APPLY
			return
		}
	} else {
		errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_INVALID_MCID
		return
	}

	var joinGuildTime time.Time
	m.dep.Guild().FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		joinGuildTime = g.GetMember(hc.Id()).GetCreateTime()
	})

	if !m.doFunc("joinFight", func() {
		joinGuildMaxTime := m.mingcWar.currStage().stageStartTime()
		if joinGuildTime.After(joinGuildMaxTime) {
			errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_JOIN_GUILD_TOO_LATE
			return
		}
		succMsg, errMsg = joinFight(m.mingcWar, mcId, gid, hc.Id(), heroProto, caps, bornPos, m.dep.Time().CurrentTime())
	}) {
		errMsg = mingc_war.ERR_JOIN_FIGHT_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) QuitFight(heroId int64) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("quitFight", func() {
		succMsg, errMsg = quitFight(m.mingcWar, heroId)
	}) {
		errMsg = mingc_war.ERR_QUIT_FIGHT_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) SceneMove(heroId int64, dest cb.Cube) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("sceneMove", func() {
		ctime := m.dep.Time().CurrentTime()
		succMsg, errMsg = sceneMove(m.mingcWar, heroId, dest, ctime)
	}) {
		errMsg = mingc_war.ERR_SCENE_MOVE_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) SceneBack(heroId int64) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("SceneBack", func() {
		ctime := m.dep.Time().CurrentTime()
		succMsg, errMsg = sceneBack(m.mingcWar, heroId, ctime)
	}) {
		errMsg = mingc_war.ERR_SCENE_BACK_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) SceneSpeedUp(heroId int64, speedUpRate float64) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("sceneSpeedUp", func() {
		ctime := m.dep.Time().CurrentTime()
		succMsg, errMsg = sceneSpeedUp(m.mingcWar, heroId, speedUpRate, ctime)
	}) {
		errMsg = mingc_war.ERR_SCENE_SPEED_UP_FAIL_SERVER_ERROR
	}
	return
}

func (m *MingcWarService) SceneTroopRelive(heroId int64) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("sceneRelive", func() {
		ctime := m.dep.Time().CurrentTime()
		succMsg, errMsg = sceneRelive(m.mingcWar, heroId, ctime)
	}) {
		errMsg = mingc_war.ERR_SCENE_TROOP_RELIVE_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) SceneChangeMode(heroId int64, newMode shared_proto.MingcWarModeType) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("changeMode", func() {
		succMsg, errMsg = sceneChangeMode(m.mingcWar, heroId, newMode)
	}) {
		errMsg = mingc_war.ERR_SCENE_CHANGE_MODE_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) SceneTouShiBuildingTurnTo(heroId int64, pos cb.Cube, left bool) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("sceneTouShiBuildingTurnTo", func() {
		ctime := m.dep.Time().CurrentTime()
		succMsg, errMsg = sceneTouShiBuildingTurnTo(m.mingcWar, heroId, pos, left, ctime)
	}) {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_TURN_TO_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) SceneTouShiBuildingFire(heroId int64, pos cb.Cube) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("sceneTouShiBuildingFire", func() {
		ctime := m.dep.Time().CurrentTime()
		succMsg, errMsg = sceneTouShiBuildingFire(m.mingcWar, heroId, pos, ctime)
	}) {
		errMsg = mingc_war.ERR_SCENE_TOU_SHI_BUILDING_FIRE_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) SceneDrum(heroId int64) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("sceneDrum", func() {
		ctime := m.dep.Time().CurrentTime()
		succMsg, errMsg = sceneDrum(m.mingcWar, heroId, ctime)
	}) {
		errMsg = mingc_war.ERR_SCENE_DRUM_FAIL_SERVER_ERR
	}
	return
}

func (m *MingcWarService) JoiningFightMingc(heroId int64) (mcId uint64, joined bool) {
	if m.mingcWar.state != shared_proto.MingcWarState_MC_T_FIGHT {
		return
	}
	if id, ok := m.mingcWar.joinedHeros[heroId]; !ok {
		return
	} else {
		if mc, ok := m.mingcWar.mingcs[id]; !ok {
			return
		} else {
			if mc.scene.ended {
				return
			}
		}
		return id, ok
	}
	return
}

func (m *MingcWarService) BuildFightStartMsg(ctime time.Time) (msg pbutil.Buffer, isFightState bool) {
	// 名城战在战斗阶段，推送展示图标
	state, start, end := m.CurrMcWarStage()
	if state != int32(shared_proto.MingcWarState_MC_T_FIGHT) {
		return
	}

	isFightState = true
	fightStart := start.Add(m.dep.Datas().MingcMiscData().FightPrepareDuration)
	if ctime.Before(fightStart) {
		msg = mingc_war.NewS2cMingcWarFightPrepareStartMsg(timeutil.Marshal32(start), timeutil.Marshal32(fightStart))
	} else {
		msg = mingc_war.NewS2cMingcWarFightStartMsg(timeutil.Marshal32(fightStart), timeutil.Marshal32(end))
	}

	return
}

func (m *MingcWarService) GuildMcWarType(gid int64) (mcId uint64, t shared_proto.MingcWarGuildType) {
	for _, mc := range m.mingcWar.mingcs {
		if mc.atkId == gid {
			mcId = mc.id
			t = shared_proto.MingcWarGuildType_MC_G_ATK
			return
		}
		if mc.defId == gid {
			mcId = mc.id
			t = shared_proto.MingcWarGuildType_MC_G_DEF
			return
		}
		for _, id := range mc.astAtkList {
			if id == gid {
				mcId = mc.id
				t = shared_proto.MingcWarGuildType_MC_G_AST_ATK
				return
			}
		}
		for _, id := range mc.astDefList {
			if id == gid {
				mcId = mc.id
				t = shared_proto.MingcWarGuildType_MC_G_AST_DEF
				return
			}
		}
	}
	return
}

func (m *MingcWarService) CleanOnGuildRemoved(gid int64) {
	switch m.mingcWar.state {
	case shared_proto.MingcWarState_MC_T_APPLY_ATK:
		stage := m.mingcWar.stages[shared_proto.MingcWarState_MC_T_APPLY_ATK].(*ApplyAtkObj)
		delete(stage.applicants, gid)
	case shared_proto.MingcWarState_MC_T_APPLY_AST:
		stage := m.mingcWar.stages[shared_proto.MingcWarState_MC_T_APPLY_AST].(*ApplyAstObj)

		apyAstMcIds := stage.ApplyAstGuilds[gid]
		for _, mcId := range apyAstMcIds {
			if mc := m.mingcWar.mingcs[mcId]; mc != nil {
				delete(mc.applyAsts, gid)
			}
		}
		delete(stage.ApplyAstGuilds, gid)

		astMcIds := stage.AstGuilds[gid]
		for _, mcId := range astMcIds {
			if mc := m.mingcWar.mingcs[mcId]; mc != nil {
				mc.astAtkList = i64.RemoveIfPresent(mc.astAtkList, gid)
				mc.astDefList = i64.RemoveIfPresent(mc.astDefList, gid)
			}
		}
		delete(stage.AstGuilds, gid)
	}

	m.updateMsg()
}

func (m *MingcWarService) ApplyAtkNotice(heroId, gid int64) (ok bool) {
	g := m.dep.GuildSnapshot().GetSnapshot(gid)
	if g == nil || g.LeaderId != heroId {
		return
	}
	if g.GuildLevel.Level <= m.dep.Datas().MingcMiscData().RedPointMinGuildLevel {
		return
	}
	if g.Id == m.dep.Mingc().AllInOneGuild() {
		return
	}

	m.doFunc("applyAtkNotice", func() {
		ok = m.mingcWar.applyAtkNotice(gid)
	})
	return
}

func (m *MingcWarService) ApplyAstNotice(heroId, gid int64) (ok bool) {
	g := m.dep.GuildSnapshot().GetSnapshot(gid)
	if g == nil || g.LeaderId != heroId {
		return
	}

	m.doFunc("ApplyAstNotice", func() {
		ok = m.mingcWar.applyAstNotice(gid)
	})
	return
}

func (m *MingcWarService) ViewSceneTroopRecord(heroId int64) (succMsg, errMsg pbutil.Buffer) {
	if !m.doFunc("ViewSceneTroopRecord", func() {
		if m.mingcWar.state != shared_proto.MingcWarState_MC_T_FIGHT {
			errMsg = mingc_war.ERR_VIEW_SCENE_TROOP_RECORD_FAIL_INVALID_TIME
			return
		}

		var scene *McWarScene
		if mcId, ok := m.mingcWar.joinedHeros[heroId]; !ok {
			errMsg = mingc_war.ERR_VIEW_SCENE_TROOP_RECORD_FAIL_NOT_IN_SCENE
			return
		} else if mc, ok := m.mingcWar.mingcs[mcId]; !ok || mc == nil {
			logrus.Warnf("ViewSceneTroopRecord。找不到名城.mcId:%v", mcId)
			errMsg = mingc_war.ERR_VIEW_SCENE_TROOP_RECORD_FAIL_NOT_IN_SCENE
			return
		} else {
			if mc.scene == nil {
				logrus.Warnf("ViewSceneTroopRecord。名城战时找不到场景.mcId:%v", mcId)
				errMsg = mingc_war.ERR_VIEW_SCENE_TROOP_RECORD_FAIL_SERVER_ERR
				return
			}
			scene = mc.scene
		}

		records := scene.viewTroopRecord(heroId)
		var data [][]byte
		for _, record := range records {
			data = append(data, must.Marshal(record))
		}

		succMsg = mingc_war.NewS2cViewSceneTroopRecordMsg(data)
	}) {
		errMsg = mingc_war.ERR_VIEW_SCENE_TROOP_RECORD_FAIL_SERVER_ERR
	}
	return
}

//func (m *MingcWarService) SceneChat(hc iface.HeroController, word string) (errMsg pbutil.Buffer) {
//	if m.mingcWar.state != shared_proto.MingcWarState_MC_T_FIGHT {
//		errMsg = mingc_war.ERR_SCENE_CHAT_FAIL_NOT_JOIN_FIGHT
//		return
//	}
//
//	gid, ok := hc.LockGetGuildId()
//	if !ok {
//		errMsg = mingc_war.ERR_SCENE_CHAT_FAIL_NOT_JOIN_FIGHT
//		return
//	}
//
//	hero := m.dep.HeroSnapshot().GetBasicProto(hc.Id())
//
//	if !m.doFunc("SceneChat", func() {
//		var mc *MingcObj
//		if mcId, ok := m.mingcWar.joinedHeros[hc.Id()]; !ok {
//			errMsg = mingc_war.ERR_SCENE_CHAT_FAIL_NOT_JOIN_FIGHT
//			return
//		} else {
//			mc = m.mingcWar.mingcs[mcId]
//		}
//
//		msg := mingc_war.NewS2cSceneOtherChatMsg(hc.IdBytes(), hero, word)
//		mc.scene.broadcastCampExclude(msg, mc.isAtk(gid), hc.Id())
//	}) {
//		errMsg = mingc_war.ERR_SCENE_CHAT_FAIL_NOT_JOIN_FIGHT
//	}
//	return
//}

func (m *MingcWarService) GmChangeStage(newState shared_proto.MingcWarState, hc iface.HeroController, ctime time.Time) {
	if m.mingcWar.state == newState {
		return
	}

	m.doFunc("sceneMove", func() {
		oldStartTime, oldEndTime := m.mingcWar.currStage().stageStartTime(), m.mingcWar.currStage().stageEndTime()

		stage := m.mingcWar.stages[newState]
		dur := stage.stageStartTime().Sub(ctime)
		gmTransMingcWarTime(m.mingcWar, -dur)
		succ, oldState := m.mingcWar.changeState(ctime)
		if succ {
			m.mingcWar.stage(oldState).onEnd(m.mingcWar, m.dep)
			m.mingcWar.currStage().onStart(m.mingcWar, m.dep)
		}

		newStartTime, newEndTime := m.mingcWar.currStage().stageStartTime(), m.mingcWar.currStage().stageEndTime()
		logrus.Debugf("名城战改变阶段 从 state:%v start:%v end:%v \n 改为 state:%v start:%v end:%v", oldState, oldStartTime, oldEndTime, m.mingcWar.state, newStartTime, newEndTime)

		m.updateMsg()
	})
}

func (m *MingcWarService) GmApplyAtkGuild(mcId uint64, gid int64) {
	mc := m.dep.Datas().GetMingcBaseData(mcId)
	if mc == nil {
		logrus.Debugf("GmApplyAtkGuild mcId:%v 不存在", mcId)
		return
	}
	g := m.dep.GuildSnapshot().GetSnapshot(gid)
	if g == nil {
		logrus.Debugf("GmApplyAtkGuild gid:%v 不存在", gid)
		return
	}
	ctime := m.dep.Time().CurrentTime()
	applyAttack(m.mingcWar, gid, mc, u64.FromInt64(ctime.Unix()), m.dep.Time().CurrentTime())
	m.updateMsg()
}

func (m *MingcWarService) GmSetDefGuild(mcId uint64, gid int64) {
	if m.dep.Datas().GetMingcBaseData(mcId) == nil {
		logrus.Debugf("GmSetDefGuild mcId:%v 不存在", mcId)
		return
	}

	if gid > 0 && m.dep.GuildSnapshot().GetSnapshot(gid) == nil {
		logrus.Debugf("GmSetDefGuild gid:%v 不存在", gid)
		return
	}

	mc := m.mingcWar.mingcs[mcId]
	if mc.isAtk(gid) {
		logrus.Debugf("GmSetDefGuild gid:%v 已经是攻方", gid)
		return
	}

	mc.defId = gid
	m.updateMsg()
}

func (m *MingcWarService) GmSetAstAtkGuild(mcId uint64, gid int64) {
	if m.dep.Datas().GetMingcBaseData(mcId) == nil {
		logrus.Debugf("GmSetAstAtkGuild mcId:%v 不存在", mcId)
		return
	}

	if m.dep.GuildSnapshot().GetSnapshot(gid) == nil {
		logrus.Debugf("GmSetAstAtkGuild gid:%v 不存在", gid)
		return
	}

	mc := m.mingcWar.mingcs[mcId]
	if mc.isDef(gid) {
		logrus.Debugf("GmSetDefGuild gid:%v 已经是守方", gid)
		return
	}

	mc.astAtkList = append(mc.astAtkList, gid)
	m.updateMsg()
}

func (m *MingcWarService) GmSetAstDefGuild(mcId uint64, gid int64) {
	if m.dep.Datas().GetMingcBaseData(mcId) == nil {
		logrus.Debugf("GmSetAstDefGuild mcId:%v 不存在", mcId)
		return
	}

	if m.dep.GuildSnapshot().GetSnapshot(gid) == nil {
		logrus.Debugf("GmSetAstDefGuild gid:%v 不存在", gid)
		return
	}

	mc := m.mingcWar.mingcs[mcId]
	if mc.isAtk(gid) {
		logrus.Debugf("GmSetDefGuild gid:%v 已经是攻方", gid)
		return
	}

	mc.astDefList = append(mc.astDefList, gid)
	m.updateMsg()
}

func (m *MingcWarService) GmNewMingcWar() {
	ctime := m.dep.Time().CurrentTime()
	m.startNewWar(ctime.Add(-time.Minute)) // 前移1分钟，保证ctime < 下场开始时间
	logrus.Debugf("GM 新开名城战 start:%v end:%v", m.mingcWar.startTime, m.mingcWar.endTime)
	m.updateMsg()
}

func (m *MingcWarService) GmCampFail(mcId uint64, atk bool) {
	if mc, ok := m.mingcWar.mingcs[mcId]; ok {
		var b *McWarBuilding
		if atk {
			b = mc.scene.atkHomeBuilding
		} else {
			b = mc.scene.defHomeBuilding
		}
		b.reduceBuildingPropsperity(b.prosperity+10000, mc.scene)
	} else {
		logrus.Warnf("GmCampFail.mcId:%v 不存在", mcId)
	}
}

func (m *MingcWarService) SendChat(sendId int64, proto *shared_proto.ChatMsgProto) (errMsg pbutil.Buffer) {
	if mcId, ok := m.mingcWar.joinedHeros[sendId]; !ok {
		errMsg = chat.ERR_SEND_CHAT_FAIL_MC_WAR_NOT_JOIN_FIGHT
		return
	} else if mc, ok := m.mingcWar.mingcs[mcId]; !ok || mc == nil || mc.scene == nil {
		errMsg = chat.ERR_SEND_CHAT_FAIL_MC_WAR_NOT_IN_FIGHT_STAGE
		return
	} else {
		var isAtk bool
		if _, ok := mc.scene.defTroops[sendId]; ok {
			isAtk = false
			mc.scene.recordDefChat(proto)
		} else if _, ok := mc.scene.atkTroops[sendId]; ok {
			isAtk = true
			mc.scene.recordAtkChat(proto)
		} else {
			logrus.Debug("发送名城战聊天，攻守都没有，玩自撸？")
			errMsg = chat.ERR_SEND_CHAT_FAIL_MC_WAR_NOT_IN_FIGHT_STAGE
			return
		}
		mc.scene.broadcastCampExclude(chat.NewS2cOtherSendChatMarshalMsg(proto), isAtk, true, sendId)
	}
	return
}

func (m *MingcWarService) CatchHistoryRecord(heroId int64, minChatId int64) (sendMsg pbutil.Buffer) {
	if mcId, ok := m.mingcWar.joinedHeros[heroId]; ok {
		if mc, ok := m.mingcWar.mingcs[mcId]; ok && mc != nil && mc.scene != nil {
			if _, ok := mc.scene.defTroops[heroId]; ok {
				sendMsg = mc.scene.catchDefChatRecord(minChatId)
			} else if _, ok := mc.scene.atkTroops[heroId]; ok {
				sendMsg = mc.scene.catchAtkChatRecord(minChatId)
			}
		}
	}
	return
}

func (m *MingcWarService) CatchTroopsRank(heroId int64, version uint64) (resultMsg, myRankMsg pbutil.Buffer) {
	if m.mingcWar.state != shared_proto.MingcWarState_MC_T_FIGHT {
		resultMsg = mingc_war.ERR_APPLY_REFRESH_RANK_FAIL_INVALID_TIME
	} else if mcId, ok := m.mingcWar.joinedHeros[heroId]; !ok {
		resultMsg = mingc_war.ERR_APPLY_REFRESH_RANK_FAIL_NOT_IN_SCENE
	} else if mc, ok := m.mingcWar.mingcs[mcId]; !ok || mc == nil || mc.scene == nil {
		resultMsg = mingc_war.ERR_APPLY_REFRESH_RANK_FAIL_SERVER_ERR
	} else {
		rankRef := mc.scene.troopsRank.loadRef()
		if rankRef.sortVersion > version {
			resultMsg = rankRef.rankMsg
			if data, ok := rankRef.dataMap[heroId]; ok {
				if data.rank > mc_rank_max_num {
					myRankMsg = mingc_war.NewS2cMyRankMsg(data.encode4Rank())
				}
			}
		} else {
			resultMsg = sameVersionRankMsg
		}
	}
	return
}
