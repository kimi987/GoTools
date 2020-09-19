package mingc_war

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func buildMcWarId(startTime time.Time) int32 {
	return timeutil.Marshal32(startTime)
}

func newMingcWar(ctime time.Time, dep iface.ServiceDep, d *mingcdata.MingcTimeData) *MingcWar {
	schedule, state := d.GetNextSchedule(ctime)

	w := NewDefaultMingcWar(dep)
	w.state = state
	w.id = buildMcWarId(schedule[shared_proto.MingcWarState_MC_T_APPLY_ATK][0])

	dep.Mingc().WalkMingcs(func(c *entity.Mingc) {
		if dep.Datas().GetMingcWarSceneData(c.Id()) == nil {
			return
		}
		w.mingcs[c.Id()] = newMingcObj(c.Id(), w.id, c.HostGuildId(), 0, dep, schedule[shared_proto.MingcWarState_MC_T_APPLY_ATK][1])
	})

	w.stages[shared_proto.MingcWarState_MC_T_NOT_START] = newStage(shared_proto.MingcWarState_MC_T_NOT_START, ctime, schedule[shared_proto.MingcWarState_MC_T_APPLY_ATK][0])
	w.stages[shared_proto.MingcWarState_MC_T_APPLY_ATK] = newApplyAtkObj(shared_proto.MingcWarState_MC_T_APPLY_ATK, schedule[shared_proto.MingcWarState_MC_T_APPLY_ATK][0], schedule[shared_proto.MingcWarState_MC_T_APPLY_ATK][1])
	w.stages[shared_proto.MingcWarState_MC_T_APPLY_AST] = newApplyAstObj(shared_proto.MingcWarState_MC_T_APPLY_AST, schedule[shared_proto.MingcWarState_MC_T_APPLY_AST][0], schedule[shared_proto.MingcWarState_MC_T_APPLY_AST][1])
	w.stages[shared_proto.MingcWarState_MC_T_FIGHT] = newFightObj(shared_proto.MingcWarState_MC_T_FIGHT, schedule[shared_proto.MingcWarState_MC_T_FIGHT][0], schedule[shared_proto.MingcWarState_MC_T_FIGHT][1])
	w.stages[shared_proto.MingcWarState_MC_T_FIGHT_END] = newStage(shared_proto.MingcWarState_MC_T_FIGHT_END, schedule[shared_proto.MingcWarState_MC_T_FIGHT_END][0], schedule[shared_proto.MingcWarState_MC_T_FIGHT_END][1])

	w.startTime = schedule[shared_proto.MingcWarState_MC_T_APPLY_ATK][0]
	w.endTime = schedule[shared_proto.MingcWarState_MC_T_FIGHT_END][1]

	return w
}

type MingcWar struct {
	id int32

	datas iface.ConfigDatas

	state shared_proto.MingcWarState

	startTime time.Time
	endTime   time.Time

	stages map[shared_proto.MingcWarState]MingcWarStageable // 每个阶段

	mingcs map[uint64]*MingcObj // 所有名城

	joinedHeros map[int64]uint64 // 参战玩家。一定要在队列里操作
}

func NewDefaultMingcWar(dep iface.ServiceDep) *MingcWar {
	w := &MingcWar{}
	w.datas = dep.Datas()
	w.stages = make(map[shared_proto.MingcWarState]MingcWarStageable)
	w.mingcs = make(map[uint64]*MingcObj)
	w.joinedHeros = make(map[int64]uint64)

	return w
}

func (c *MingcWar) stage(state shared_proto.MingcWarState) MingcWarStageable {
	return c.stages[state]
}

func (c *MingcWar) currStage() MingcWarStageable {
	return c.stages[c.state]
}

func (c *MingcWar) changeState(ctime time.Time) (changed bool, oldState shared_proto.MingcWarState) {
	oldState = c.state
	if between(ctime, c.currStage().stageStartTime(), c.currStage().stageEndTime()) {
		return
	}

	for k, v := range c.stages {
		if !between(ctime, v.stageStartTime(), v.stageEndTime()) {
			continue
		}

		c.state = k
		changed = true
		return
	}
	return
}

func (w *MingcWar) bidValue(gid int64) uint64 {
	stage := w.stage(shared_proto.MingcWarState_MC_T_APPLY_ATK).(*ApplyAtkObj)

	mcId := stage.applicants[gid]
	if mcId <= 0 {
		return 0
	}

	if mc := w.mingcs[mcId]; mc != nil {
		return mc.bid.GetBid(gid)
	}

	return 0
}

func between(ctime time.Time, startTime time.Time, endTime time.Time) bool {
	return !ctime.Before(startTime) && ctime.Before(endTime)
}

func (c *MingcWar) encodeServer() *server_proto.MingcWarServerProto {
	p := &server_proto.MingcWarServerProto{}
	p.Id = c.id
	p.State = c.state
	p.StartTime = timeutil.Marshal64(c.startTime)
	p.EndTime = timeutil.Marshal64(c.endTime)

	p.NotStart = c.stage(shared_proto.MingcWarState_MC_T_NOT_START).(*Stage).encodeServer()
	p.ApplyAtk = c.stage(shared_proto.MingcWarState_MC_T_APPLY_ATK).(*ApplyAtkObj).encodeServer()
	p.ApplyAst = c.stage(shared_proto.MingcWarState_MC_T_APPLY_AST).(*ApplyAstObj).encodeServer()
	p.Fight = c.stage(shared_proto.MingcWarState_MC_T_FIGHT).(*FightObj).encodeServer()
	p.FightEnd = c.stage(shared_proto.MingcWarState_MC_T_FIGHT_END).(*Stage).encodeServer()

	p.Mingc = make(map[uint64]*server_proto.MingcObjProto)
	for k, v := range c.mingcs {
		p.Mingc[k] = v.encodeServer()
	}

	return p
}

func (c *MingcWar) encode(guildGetter func(id int64) *guildsnapshotdata.GuildSnapshot) *shared_proto.McWarProto {
	p := &shared_proto.McWarProto{}
	p.State = c.state
	p.StartTime = timeutil.Marshal32(c.startTime)
	p.EndTime = timeutil.Marshal32(c.endTime)

	for _, v := range c.mingcs {
		p.Mc = append(p.Mc, v.encode(guildGetter))
	}

	p.Atk = c.stage(shared_proto.MingcWarState_MC_T_APPLY_ATK).(*ApplyAtkObj).encode()
	p.Ast = c.stage(shared_proto.MingcWarState_MC_T_APPLY_AST).(*ApplyAstObj).encode()
	p.Fight = c.stage(shared_proto.MingcWarState_MC_T_FIGHT).(*FightObj).encode()

	return p
}

func (c *MingcWar) encodeGuild(gid int64, guildGetter func(id int64) *guildsnapshotdata.GuildSnapshot, heroGetter func(id int64) *shared_proto.HeroBasicSnapshotProto) *shared_proto.McWarGuildProto {
	p := &shared_proto.McWarGuildProto{}
	atkStage := c.stage(shared_proto.MingcWarState_MC_T_APPLY_ATK).(*ApplyAtkObj)
	if mcId, ok := atkStage.applicants[gid]; ok {
		p.ApplyMcId = u64.Int32(mcId)
		p.ApplyMcValue = u64.Int32(c.mingcs[mcId].bid.GetBid(gid))
	}

	for _, mc := range c.mingcs {
		astListProto := &shared_proto.McWarGuildApplyAstListProto{McId: u64.Int32(mc.id)}
		for id, isAtk := range mc.applyAsts {
			g := guildGetter(id)
			if g == nil {
				continue
			}
			if mc.atkId == gid && isAtk {
				astListProto.Guild = append(astListProto.Guild, g.Encode(heroGetter))
			} else if mc.defId == gid && !isAtk {
				astListProto.Guild = append(astListProto.Guild, g.Encode(heroGetter))
			} else if id == gid {
				if isAtk {
					p.ApplyAstAtkMcId = append(p.ApplyAstAtkMcId, u64.Int32(mc.id))
				} else {
					p.ApplyAstDefMcId = append(p.ApplyAstDefMcId, u64.Int32(mc.id))
				}
			}
		}
		p.ReqAstGuild = append(p.ReqAstGuild, astListProto)
	}

	return p
}

func (c *MingcWar) unmarshal(p *server_proto.MingcWarServerProto, dep iface.ServiceDep) {
	c.id = p.Id
	c.state = p.State
	c.startTime = timeutil.Unix64(p.StartTime)
	c.endTime = timeutil.Unix64(p.EndTime)

	for k, v := range p.Mingc {
		if dep.Datas().GetMingcBaseData(v.Id) == nil {
			logrus.Warnf("MingcWar.unmarshal 名城：%v 配置被删掉了, 忽略这座名城的 MingcObj.Unmarshal", v.Id)
			continue
		}

		mc := newDefaultMingcObj()
		mc.unmarshal(v, dep)
		mc.warId = c.id
		c.mingcs[k] = mc
		if mc.scene != nil {
			for _, t := range mc.scene.atkTroops {
				c.joinedHeros[t.heroId] = mc.id
			}
			for _, t := range mc.scene.defTroops {
				c.joinedHeros[t.heroId] = mc.id
			}
		}
	}

	nstart := newStage(shared_proto.MingcWarState_MC_T_NOT_START, time.Time{}, time.Time{})
	nstart.unmarshal(shared_proto.MingcWarState_MC_T_NOT_START, p.NotStart)
	c.stages[shared_proto.MingcWarState_MC_T_NOT_START] = nstart

	atk := newApplyAtkObj(shared_proto.MingcWarState_MC_T_APPLY_ATK, time.Time{}, time.Time{})
	atk.unmarshal(p.ApplyAtk, c.mingcs)
	c.stages[shared_proto.MingcWarState_MC_T_APPLY_ATK] = atk

	ast := newApplyAstObj(shared_proto.MingcWarState_MC_T_APPLY_AST, time.Time{}, time.Time{})
	ast.unmarshal(p.ApplyAst, c.mingcs)
	c.stages[shared_proto.MingcWarState_MC_T_APPLY_AST] = ast

	fgt := newFightObj(shared_proto.MingcWarState_MC_T_FIGHT, time.Time{}, time.Time{})
	fgt.unmarshal(p.Fight)
	c.stages[shared_proto.MingcWarState_MC_T_FIGHT] = fgt

	fgte := newStage(shared_proto.MingcWarState_MC_T_FIGHT_END, time.Time{}, time.Time{})
	fgte.unmarshal(shared_proto.MingcWarState_MC_T_FIGHT_END, p.FightEnd)
	c.stages[shared_proto.MingcWarState_MC_T_FIGHT_END] = fgte
}

func (c *MingcWar) applyAtkNotice(gid int64) (ok bool) {
	if c.state != shared_proto.MingcWarState_MC_T_APPLY_ATK {
		return
	}

	stage := c.stage(shared_proto.MingcWarState_MC_T_APPLY_ATK).(*ApplyAtkObj)
	if _, exist := stage.applicants[gid]; !exist {
		return true
	}
	return
}

func (c *MingcWar) applyAstNotice(gid int64) (ok bool) {
	if c.state != shared_proto.MingcWarState_MC_T_APPLY_AST {
		return
	}

	for _, mc := range c.mingcs {
		if mc.atkId != gid && mc.defId != gid {
			continue
		}

		for _, isAtk := range mc.applyAsts {
			if mc.atkId == gid {
				if isAtk {
					return true
				}
			} else {
				if !isAtk {
					return true
				}
			}
		}
	}

	return
}

// ********** stage **********

type MingcWarStageable interface {
	warState() shared_proto.MingcWarState
	stageStartTime() time.Time
	stageEndTime() time.Time
	onStart(w *MingcWar, dep iface.ServiceDep)
	onEnd(w *MingcWar, dep iface.ServiceDep)
	onUpdate(w *MingcWar, dep iface.ServiceDep) (changed bool)
}

type Stage struct {
	state     shared_proto.MingcWarState
	startTime time.Time
	endTime   time.Time
}

func (c *Stage) warState() shared_proto.MingcWarState                      { return c.state }
func (c *Stage) stageStartTime() time.Time                                 { return c.startTime }
func (c *Stage) stageEndTime() time.Time                                   { return c.endTime }
func (c *Stage) onStart(w *MingcWar, dep iface.ServiceDep)                 {}
func (c *Stage) onEnd(w *MingcWar, dep iface.ServiceDep)                   {}
func (c *Stage) onUpdate(w *MingcWar, dep iface.ServiceDep) (changed bool) { return false }

type ApplyAtkObj struct {
	*Stage

	applicants map[int64]uint64
	winners    map[uint64]int64 // 竞拍成功者
}

func newApplyAtkObj(state shared_proto.MingcWarState, startTime, endTime time.Time) *ApplyAtkObj {
	c := &ApplyAtkObj{Stage: &Stage{}}
	c.state = state
	c.startTime = startTime
	c.endTime = endTime

	c.applicants = make(map[int64]uint64)
	c.winners = make(map[uint64]int64)

	return c
}

func (c *ApplyAtkObj) onStart(w *MingcWar, dep iface.ServiceDep) {
	for _, mc := range w.mingcs {
		mc.bid = entity.NewBidInfo(c.endTime)
	}
	allInOneGuildId := dep.Mingc().AllInOneGuild()
	dep.Guild().Func(func(guilds sharedguilddata.Guilds) {
		guilds.Walk(func(g *sharedguilddata.Guild) {
			if g.LevelData().Level <= dep.Datas().MingcMiscData().RedPointMinGuildLevel {
				return
			}
			if g.Id() == allInOneGuildId {
				return
			}
			dep.World().Send(g.LeaderId(), mingc_war.RED_POINT_NOTICE_S2C)
		})
	})
}

func (c *ApplyAtkObj) onEnd(w *MingcWar, dep iface.ServiceDep) {
	ctime := dep.Time().CurrentTime()
	for mcId, mc := range w.mingcs {
		// 成功通知
		if winner, _, succ := mc.bid.Winner(ctime); succ {
			c.winners[mcId] = winner
			mc.atkId = winner
			if g := dep.GuildSnapshot().GetSnapshot(winner); g != nil {
				dep.World().Send(g.LeaderId, mingc_war.NewS2cApplyAtkSuccMsg(u64.Int32(mcId)))

				// 联盟全员发邮件
				if d := dep.Datas().MailHelp().GuildMcWarApplyAtkSucc; d != nil {
					if mc := dep.Datas().GetMingcBaseData(mcId); mc != nil {
						proto := d.NewTextMail(shared_proto.MailType_MailNormal)
						proto.Text = d.NewTextFields().WithMingc(mc.Name).JsonString()
						for _, heroId := range g.UserMemberIds {
							dep.Mail().SendProtoMail(heroId, proto, ctime)
						}
					}
				}

			}
			// 联盟事件
			if d := dep.Datas().GuildLogHelp().McWarApplyAtkSucc; d != nil {
				if mc := dep.Datas().GetMingcBaseData(mcId); mc != nil {
					proto := d.NewLogProto(ctime)
					proto.Text = d.Text.New().WithMingc(mc.Name).JsonString()
					dep.Guild().AddLog(winner, proto)
				}
			}
		}

		// 还虎符
		if costBacks, succ := mc.bid.Losers(ctime); !succ {
			logrus.Errorf("计算竞拍失败者失败2, 时间错误，ctime:%v", ctime)
		} else {
			for gid, cost := range costBacks {
				asyncCostBack(dep, mcId, gid, cost)
			}
		}
	}

	// 构造 MingcObj
	mcHosts := make(map[uint64]int64)
	dep.Mingc().WalkMingcs(func(c *entity.Mingc) {
		if dep.Datas().GetMingcWarSceneData(c.Id()) == nil {
			return
		}
		if c.HostGuildId() > 0 {
			mcHosts[c.Id()] = c.HostGuildId()
		}
	})

	for mcId, gid := range c.winners {
		if _, ok := w.mingcs[mcId]; !ok {
			w.mingcs[mcId] = newMingcObj(mcId, w.id, 0, gid, dep, c.endTime)
		}
	}
}

func asyncCostBack(dep iface.ServiceDep, mcId uint64, gid int64, cost uint64) {
	ctime := dep.Time().CurrentTime()

	go call.CatchPanic(func() {
		var leaderId int64
		dep.Guild().FuncGuild(gid, func(g *sharedguilddata.Guild) {
			if g == nil {
				return
			}

			g.AddHufu(cost)
			leaderId = g.LeaderId()
		})
		dep.World().Send(leaderId, mingc_war.NewS2cApplyAtkFailMsg(u64.Int32(mcId), u64.Int32(cost)))

		if mc := dep.Datas().GetMingcBaseData(mcId); mc != nil {
			// 邮件
			if d := dep.Datas().MailHelp().GuildMcWarApplyAtkFail; d != nil {
				proto := d.NewTextMail(shared_proto.MailType_MailNormal)
				proto.Text = d.NewTextFields().WithMingc(mc.Name).WithAmount(cost).JsonString()
				dep.Mail().SendProtoMail(leaderId, proto, ctime)
			}
			// 联盟事件
			if d := dep.Datas().GuildLogHelp().McWarApplyAtkFail; d != nil {
				proto := d.NewLogProto(ctime)
				proto.Text = d.Text.New().WithMingc(mc.Name).WithAmount(cost).JsonString()
				dep.Guild().AddLog(gid, proto)
			}
		}

	}, "MingcWarService.asyncCostBack")
}

func (c *ApplyAtkObj) encode() *shared_proto.McApplyAtkProto {
	p := &shared_proto.McApplyAtkProto{}
	p.StartTime = timeutil.Marshal32(c.startTime)
	p.EndTime = timeutil.Marshal32(c.endTime)
	return p
}

func (c *ApplyAtkObj) encodeServer() *server_proto.MingcApplyAtkProto {
	p := &server_proto.MingcApplyAtkProto{}
	p.StartTime = timeutil.Marshal64(c.startTime)
	p.EndTime = timeutil.Marshal64(c.endTime)
	p.Applicant = c.applicants
	p.Winners = c.winners

	return p
}

func (c *ApplyAtkObj) unmarshal(p *server_proto.MingcApplyAtkProto, mcs map[uint64]*MingcObj) {
	c.startTime = timeutil.Unix64(p.StartTime)
	c.endTime = timeutil.Unix64(p.EndTime)
	if p.Applicant != nil {
		c.applicants = p.Applicant
	}
	if p.Winners != nil {
		c.winners = p.Winners
	}
}

type ApplyAstObj struct {
	*Stage

	ApplyAstGuilds map[int64][]uint64 // 申请协助的联盟 k:申请人 v:名城id
	AstGuilds      map[int64][]uint64 // 确认协助的联盟 k:协助人 v:名城id
}

func newApplyAstObj(state shared_proto.MingcWarState, startTime, endTime time.Time) *ApplyAstObj {
	c := &ApplyAstObj{Stage: &Stage{}}
	c.state = state
	c.startTime = startTime
	c.endTime = endTime

	c.AstGuilds = make(map[int64][]uint64)
	c.ApplyAstGuilds = make(map[int64][]uint64)

	return c
}

func (c *ApplyAstObj) encode() *shared_proto.McApplyAstProto {
	p := &shared_proto.McApplyAstProto{}
	p.StartTime = timeutil.Marshal32(c.startTime)
	p.EndTime = timeutil.Marshal32(c.endTime)
	return p
}

func (c *ApplyAstObj) encodeServer() *server_proto.MingcApplyAstProto {
	p := &server_proto.MingcApplyAstProto{}
	p.StartTime = timeutil.Marshal64(c.startTime)
	p.EndTime = timeutil.Marshal64(c.endTime)

	return p
}

func (c *ApplyAstObj) unmarshal(p *server_proto.MingcApplyAstProto, mcs map[uint64]*MingcObj) {
	c.startTime = timeutil.Unix64(p.StartTime)
	c.endTime = timeutil.Unix64(p.EndTime)

	for _, mc := range mcs {
		for gid := range mc.applyAsts {
			c.ApplyAstGuilds[gid] = append(c.ApplyAstGuilds[gid], mc.id)
		}
		for _, gid := range mc.astAtkList {
			c.AstGuilds[gid] = append(c.AstGuilds[gid], mc.id)
		}
		for _, gid := range mc.astDefList {
			c.AstGuilds[gid] = append(c.AstGuilds[gid], mc.id)
		}
	}
}

func (c *ApplyAstObj) addApplyAstGuilds(gid int64, mcId uint64) {
	if applyMcIds, ok := c.ApplyAstGuilds[gid]; ok {
		applyMcIds = append(applyMcIds, mcId)
	} else {
		applyMcIds = []uint64{mcId}
		c.ApplyAstGuilds[gid] = applyMcIds
	}
}

func (c *ApplyAstObj) delApplyAstGuilds(gid int64, mcId uint64) {
	if applyMcIds, ok := c.ApplyAstGuilds[gid]; ok {
		c.ApplyAstGuilds[gid] = u64.RemoveIfPresent(applyMcIds, mcId)
	}
}

func (c *ApplyAstObj) addAstGuilds(gid int64, mcId uint64) {
	if asts, ok := c.AstGuilds[gid]; ok {
		asts = append(asts, mcId)
	} else {
		asts = []uint64{mcId}
		c.AstGuilds[gid] = asts
	}
}

type FightObj struct {
	*Stage
	fightStartNoticed bool
}

func newFightObj(state shared_proto.MingcWarState, startTime, endTime time.Time) *FightObj {
	c := &FightObj{Stage: &Stage{}}
	c.state = state
	c.startTime = startTime
	c.endTime = endTime

	return c
}

func (c *FightObj) encode() *shared_proto.McFightProto {
	p := &shared_proto.McFightProto{}
	p.StartTime = timeutil.Marshal32(c.startTime)
	p.EndTime = timeutil.Marshal32(c.endTime)
	return p
}

func (c *FightObj) encodeServer() *server_proto.MingcFightProto {
	p := &server_proto.MingcFightProto{}
	p.StartTime = timeutil.Marshal64(c.startTime)
	p.EndTime = timeutil.Marshal64(c.endTime)

	return p
}

func (c *FightObj) unmarshal(p *server_proto.MingcFightProto) {
	c.startTime = timeutil.Unix64(p.StartTime)
	c.endTime = timeutil.Unix64(p.EndTime)
}

func (c *FightObj) onStart(w *MingcWar, dep iface.ServiceDep) {
	w.joinedHeros = make(map[int64]uint64)
	for _, mc := range w.mingcs {
		if mc.atkId > 0 {
			mc.scene = newMcWarScene(mc, dep, c.startTime)
		}
	}

	go call.CatchPanic(func() {
		fightStartTime := c.startTime.Add(dep.Datas().MingcMiscData().FightPrepareDuration)
		dep.World().Broadcast(mingc_war.NewS2cMingcWarFightPrepareStartMsg(timeutil.Marshal32(c.startTime), timeutil.Marshal32(fightStartTime)))
	}, "mingcWarPrepareStartNotice")
}

func (c *FightObj) onEnd(w *MingcWar, dep iface.ServiceDep) {
	ctime := dep.Time().CurrentTime()
	for _, mc := range w.mingcs {
		if mc.scene == nil || mc.scene.ended {
			continue
		}
		// 结算未完成的名城战h
		mc.scene.onEnd()
		mc.onSceneEnd(true, w, dep, ctime)
	}
	w.joinedHeros = make(map[int64]uint64)

	// 删除 db 中的旧记录
	go call.CatchPanic(func() {
		days := time.Duration(dep.Datas().MingcMiscData().SaveHeroRecordMaxDays)
		earliestWarId := buildMcWarId(ctime.Add(-24 * days * time.Hour))
		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			return dep.Db().DelMcWarHeroRecord(ctx, earliestWarId)
		}); err != nil {
			logrus.Warnf("名城战结束删除旧的玩家详细记录失败")
		}
	}, "DelMcWarHeroRecord")
}

func (c *FightObj) onUpdate(w *MingcWar, dep iface.ServiceDep) (changed bool) {
	ctime := dep.Time().CurrentTime()

	if !c.fightStartNoticed {
		fightStartTime := c.startTime.Add(dep.Datas().MingcMiscData().FightPrepareDuration)
		if ctime.After(fightStartTime) {
			dep.World().Broadcast(mingc_war.NewS2cMingcWarFightStartMsg(timeutil.Marshal32(fightStartTime), timeutil.Marshal32(c.endTime)))
			c.fightStartNoticed = true
		}
	}

	for _, mc := range w.mingcs {
		if mc.scene == nil || mc.scene.ended {
			continue
		}

		mc.scene.onUpdate(ctime)
		if mc.scene.ended && !mc.scene.endedNoticed {
			mc.onSceneEnd(false, w, dep, ctime)
			changed = true
		}
	}

	return
}

func (mc *MingcObj) onSceneEnd(isTimeout bool, war *MingcWar, dep iface.ServiceDep, ctime time.Time) {
	defer func() {
		mc.scene.endedNoticed = true
	}()

	msg := mingc_war.NewS2cSceneWarEndMsg(u64.Int32(mc.id))
	mc.scene.broadcast(msg)

	// 设置名城归属
	mc.updateMingcHost(dep)

	// 分钱
	toReduce := mc.sendYinliangOnEnded(dep, ctime)

	// 排行
	mc.scene.sort4TroopRank()

	// 退出战场
	for _, t := range mc.scene.atkTroops {
		delete(war.joinedHeros, t.heroId)
	}
	for _, t := range mc.scene.defTroops {
		delete(war.joinedHeros, t.heroId)
	}

	// 结算界面
	mc.scene.record.atkYinliang = toReduce
	record, troop := mc.scene.record.encode(mc, isTimeout, ctime, getGuildBasicProto, dep)
	mc.scene.broadcast(mingc_war.NewS2cMcWarEndRecordMsg(war.id, record).Static())
	mc.saveRecordIntoDBOnEnded(dep, record, troop)

	// 联盟历史战绩列表
	mc.addGuildMcWarRecord(dep, toReduce, ctime)

	mcData := mc.dep.Datas().MingcBaseData().Get(mc.id)
	if mc.scene.atkWin {
		mc.broadcastAtkWinOnEnded(dep, mcData)
		mc.addAtkWinGuildLogOnEnded(dep, mcData, ctime)
	} else {
		mc.addDefWinGuildLogOnEnded(dep, mcData, ctime)
	}
}

// 设置名城归属
func (mc *MingcObj) updateMingcHost(dep iface.ServiceDep) {
	if !mc.scene.atkWin {
		return
	}

	mc.dep.Mingc().SetHostGuild(mc.id, mc.atkId)
	var newHostProto *shared_proto.GuildBasicProto
	var mcCount uint64
	var memberIds []int64
	dep.Guild().FuncGuild(mc.atkId, func(g *sharedguilddata.Guild) {
		if g != nil {
			mcCount = g.AddHostMingc(mc.id)
			newHostProto = g.NewBasicProto()
			memberIds = g.AllUserMemberIds()
		}
	})
	// 更新联盟坐拥名城数成就
	for _, heroId := range memberIds {
		dep.HeroData().FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
			if hero.HistoryAmount().SetIfGreater(server_proto.HistoryAmountType_McOccupy, mcCount) {
				heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_GUILD_LEVEL)
			}
			result.Ok()
		})
	}
	if mc.defId > 0 {
		dep.Guild().FuncGuild(mc.defId, func(g *sharedguilddata.Guild) {
			if g != nil {
				g.DelHostMingc(mc.id)
			}
		})
	}
	if newHostProto != nil {
		dep.World().Broadcast(mingc_war.NewS2cMingcHostUpdateNoticeMsg(u64.Int32(mc.id), newHostProto))
	}

	mc.onCapitalUpdate(dep)
}

func (mc *MingcObj) onCapitalUpdate(dep iface.ServiceDep) {
	mcData := dep.Datas().GetMingcBaseData(mc.id)
	if mcData.Type != shared_proto.MincType_MC_DU {
		return
	}

	g := dep.Guild().GetSnapshot(mc.atkId)
	if g == nil {
		logrus.Debugf("MingcObj.onCapitalUpdate mc:%v 找不到胜利联盟:%v。", mc.id, mc.atkId)
		return
	}

	if g.Country.Id == mcData.Country {
		dep.Country().CancelCountryDestroy(mcData.Country)
		dep.Country().ChangeCountryHost(mcData.Country, g.LeaderId)
	} else {
		dep.Country().CountryDestroy(mcData.Country)

		if d := dep.Datas().BroadcastHelp().CountryDestroy; d != nil {
			guildCountryName := dep.Country().CountryName(g.Country.Id)
			countryName := dep.Country().CountryName(mcData.Country)
			guildName := dep.Datas().MiscConfig().FlagHeroName.FormatIgnoreEmpty(g.FlagName, g.Name)
			dep.Broadcast().Broadcast(d.Text.New().WithGuildCountry(guildCountryName).WithGuild(guildName).WithCountry(countryName).JsonString(), d.SendChat)
		}
	}
}

// 分钱
func (mc *MingcObj) sendYinliangOnEnded(dep iface.ServiceDep, ctime time.Time) (toReduce uint64) {
	mingc := dep.Mingc().Mingc(mc.id)
	var percent uint64
	if mc.scene.atkWin {
		percent = 100
	} else {
		code := float64(mc.calcAtkDestroyed()) / float64(dep.Datas().GetMingcWarSceneData(mc.id).DefFullProsperity)
		percent = u64.MultiF64(100, code)
	}

	if mc.scene.atkWin {
		// 名城换主人，清空溢出银两
		mingc.CleanExtraYinliang()
	}
	_, toReduce = mingc.ReducePercentYinliang(percent)

	dep.Mingc().UpdateMsg()
	dep.Guild().FuncGuild(mc.atkId, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		g.AddYinliang(toReduce)
	})
	mc.addGuildYinliangRecord(dep, toReduce, percent, ctime)

	return
}

// 保存各种历史战绩到 DB
func (mc *MingcObj) saveRecordIntoDBOnEnded(dep iface.ServiceDep, record *shared_proto.McWarFightRecordProto, troop *shared_proto.McWarTroopsInfoProto) {
	// 保存部队战斗明细
	for _, t := range mc.scene.atkTroops {
		mc.saveTroopRecord(mc.warId, t)
	}
	for _, t := range mc.scene.defTroops {
		mc.saveTroopRecord(mc.warId, t)
	}

	// 保存结算战绩
	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		err = mc.dep.Db().AddMcWarRecord(ctx, u64.FromInt32(mc.warId), mc.id, record)
		return
	}); err != nil {
		logrus.WithError(err).Errorf("名城战结束保存战报失败 mcId:%v", mc.id)
	}

	// 保存公会部队战绩
	guildIdMap := make(map[int64]*shared_proto.McWarTroopsInfoProto)
	for _, t := range troop.TroopsInfo {
		guildTroop, ok := guildIdMap[int64(t.Hero.GuildId)]
		if !ok {
			guildTroop = &shared_proto.McWarTroopsInfoProto{}
			guildIdMap[int64(t.Hero.GuildId)] = guildTroop
		}
		guildTroop.TroopsInfo = append(guildTroop.TroopsInfo, t)
	}
	data := mc.dep.Datas().GetGuildTaskData(u64.FromInt32(int32(server_proto.GuildTaskType_McWar)))
	for guildId, troop := range guildIdMap {
		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			err = mc.dep.Db().AddMcWarGuildRecord(ctx, u64.FromInt32(mc.warId), mc.id, guildId, troop)
			return
		}); err != nil {
			logrus.WithError(err).Errorf("保存公会部队战绩失败 mcId:%v", mc.id)
		}
		// 更新各个联盟的任务进度
		mc.dep.Guild().AddGuildTaskProgress(guildId, data, u64.FromInt(len(troop.TroopsInfo)))
	}
}

// 攻方胜利系统广播
func (mc *MingcObj) broadcastAtkWinOnEnded(dep iface.ServiceDep, mcData *mingcdata.MingcBaseData) {
	d := dep.Datas().BroadcastHelp().McWarAtkWin
	if d == nil {
		return
	}

	g := dep.GuildSnapshot().GetSnapshot(mc.atkId)
	if g == nil {
		return
	}

	guildName := dep.Datas().MiscConfig().FlagHeroName.FormatIgnoreEmpty(g.FlagName, g.Name)
	text := d.NewTextFields().WithClickGuildFields(data.KeyGuild, guildName, g.Id)
	text.WithClickMingcFields(data.KeyMingc, mcData.Name, mcData.Id)
	if d.SendChat {
		dep.Chat().BroadcastSystemChat(text.JsonString())
	} else {
		dep.World().Broadcast(misc.NewS2cSysBroadcastMsg(text.JsonString()))
	}
}

// 攻方胜利联盟日志
func (mc *MingcObj) addAtkWinGuildLogOnEnded(dep iface.ServiceDep, mcData *mingcdata.MingcBaseData, ctime time.Time) {
	d := dep.Datas().GuildLogHelp().McWarAtkWin
	if d == nil {
		return
	}

	var guildName string
	var countryName string
	if g := dep.GuildSnapshot().GetSnapshot(mc.defId); g != nil {
		guildName = dep.Datas().MiscConfig().FlagHeroName.FormatIgnoreEmpty(g.FlagName, g.Name)
		countryName = g.Country.Name
	} else {
		if country := dep.Datas().GetCountryData(mcData.Country); country != nil {
			countryName = country.Name
		}
	}

	proto := d.NewLogProto(ctime)
	proto.Text = d.Text.New().WithMingc(mcData.Name).WithCountry(countryName).WithGuild(guildName).JsonString()
	dep.Guild().AddLog(mc.atkId, proto)
}

// 守方胜利联盟日志
func (mc *MingcObj) addDefWinGuildLogOnEnded(dep iface.ServiceDep, mcData *mingcdata.MingcBaseData, ctime time.Time) {
	if mc.defId <= 0 {
		return
	}

	d := dep.Datas().GuildLogHelp().McWarDefWin
	if d == nil {
		return
	}

	g := dep.GuildSnapshot().GetSnapshot(mc.atkId)
	if g == nil {
		return
	}

	guildName := dep.Datas().MiscConfig().FlagHeroName.FormatIgnoreEmpty(g.FlagName, g.Name)
	countryName := g.Country.Name

	proto := d.NewLogProto(ctime)
	proto.Text = d.Text.New().WithMingc(mcData.Name).WithCountry(countryName).WithGuild(guildName).JsonString()
	dep.Guild().AddLog(mc.defId, proto)

}

func (mc *MingcObj) calcAtkDestroyed() (amount uint64) {
	if g, ok := mc.scene.record.guilds[mc.atkId]; ok {
		amount += g.destroyed
	}
	for _, atk := range mc.astAtkList {
		if g, ok := mc.scene.record.guilds[atk]; ok {
			amount += g.destroyed
		}
	}
	return
}

func (mc *MingcObj) saveTroopRecord(warId int32, troop *McWarTroop) {
	if troop.heroId < 0 {
		return
	}

	var troopWin bool
	if mc.scene != nil && mc.scene.ended {
		if troop.atk {
			troopWin = mc.scene.atkWin
		} else {
			troopWin = !mc.scene.atkWin
		}
	}

	// 更新任务
	mc.dep.HeroData().FuncWithSend(troop.heroId, func(hero *entity.Hero, result herolock.LockResult) {
		changed := false

		if troop.killAmount > 0 {
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_McWarKillSolider, troop.killAmount)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_MC_WAR_KILL_SOLDIER)
			changed = true
		}

		if troop.destroyBuilding > 0 {
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_McWarDestroyBuilding, troop.destroyBuilding)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_MC_WAR_DESTROY_BUILDING)
			changed = true
		}

		if troopWin {
			// 提前退出战斗，不算胜利
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_McWarWin)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_MC_WAR_WIN)
			changed = true
		}

		if changed {
			result.Changed()
		}
		result.Ok()
	})

	if len(troop.records) > 0 {
		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			p, err := mc.dep.Db().LoadMcWarHeroRecord(ctx, u64.FromInt32(warId), mc.id, troop.heroId)
			if err != nil || p == nil {
				p = &shared_proto.McWarTroopAllRecordProto{}
			}
			p.Record = append(p.Record, troop.records...)
			return mc.dep.Db().AddMcWarHeroRecord(ctx, u64.FromInt32(warId), mc.id, troop.heroId, p)
		}); err != nil {
			logrus.WithError(err).Errorf("名城战退出战斗保存战报失败 mcId:%v heroId%v", mc.id, troop.heroId)
		}
	}
}

func (c *MingcObj) addGuildYinliangRecord(dep iface.ServiceDep, atkYinliang, percent uint64, ctime time.Time) {
	dep.Guild().FuncGuild(c.atkId, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		if c.scene.atkWin {
			if d := dep.Datas().GuildLogHelp().YinliangMcWarAtkWin; d != nil {
				mingcName := dep.Datas().GetMingcBaseData(c.id).Name
				text := d.Text.New().WithMingc(mingcName).WithYinliang(atkYinliang).JsonString()
				addYinliangRecord(g, text, ctime, d, dep.Chat())
			}
		} else {
			if d := dep.Datas().GuildLogHelp().YinliangMcWarAtkFail; d != nil {
				mingcName := dep.Datas().GetMingcBaseData(c.id).Name
				text := d.Text.New().WithMingc(mingcName).WithDestroy(percent).WithYinliang(atkYinliang).JsonString()
				addYinliangRecord(g, text, ctime, d, dep.Chat())
			}
		}

	})
}

func addYinliangRecord(g *sharedguilddata.Guild, text string, ctime time.Time, d *guild_data.GuildLogData, chat iface.ChatService) {
	g.AddYinliangRecord(text, ctime)
	if d.SendChat {
		chat.SysChat(0, g.Id(), shared_proto.ChatType_ChatGuild, text, shared_proto.ChatMsgType_ChatMsgGuildLog, true, true, true, false)
	}
}

func (c *MingcObj) addGuildMcWarRecord(dep iface.ServiceDep, atkYinliang uint64, ctime time.Time) {
	p := &shared_proto.McWarRecordProto{}
	p.WarId = c.warId
	p.Time = timeutil.Marshal32(ctime)
	p.Atk = getGuildBasicProto(c.atkId, dep)
	p.Def = getGuildBasicProto(c.defId, dep)
	p.McId = u64.Int32(c.id)
	p.McCountryId = u64.Int32(dep.Mingc().Country(c.id))
	p.AtkWin = c.scene.atkWin
	p.AtkYinliang = u64.Int32(atkYinliang)

	addGuildMcWarRecord0(c.atkId, p, shared_proto.McWarActionType_MC_WAR_A_ATK, dep)
	addGuildMcWarRecord0(c.defId, p, shared_proto.McWarActionType_MC_WAR_A_DEF, dep)
	for _, gid := range c.astAtkList {
		addGuildMcWarRecord0(gid, p, shared_proto.McWarActionType_MC_WAR_A_AST_ATK, dep)
	}
	for _, gid := range c.astDefList {
		addGuildMcWarRecord0(gid, p, shared_proto.McWarActionType_MC_WAR_A_AST_DEF, dep)
	}
}

func addGuildMcWarRecord0(gid int64, tmp *shared_proto.McWarRecordProto, atype shared_proto.McWarActionType, dep iface.ServiceDep) {
	if gid < 0 {
		return
	}

	p := &shared_proto.McWarRecordProto{}
	p.WarId = tmp.WarId
	p.Time = tmp.Time
	p.Atk = tmp.Atk
	p.Def = tmp.Def
	p.AtkWin = tmp.AtkWin
	p.AtkYinliang = tmp.AtkYinliang
	p.McId = tmp.McId
	p.McCountryId = tmp.McCountryId
	p.Type = atype

	dep.Guild().FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g != nil {
			g.AddMcWarRecord(p)
		}
	})
}

func newStage(state shared_proto.MingcWarState, startTime, endTime time.Time) *Stage {
	c := &Stage{}
	c.state = state
	c.startTime = startTime
	c.endTime = endTime

	return c
}

func (c *Stage) encodeServer() *server_proto.MingcStageProto {
	p := &server_proto.MingcStageProto{}
	p.StartTime = timeutil.Marshal64(c.startTime)
	p.EndTime = timeutil.Marshal64(c.endTime)
	return p
}

func (c *Stage) unmarshal(state shared_proto.MingcWarState, p *server_proto.MingcStageProto) {
	c.state = state
	c.startTime = timeutil.Unix64(p.StartTime)
	c.endTime = timeutil.Unix64(p.EndTime)
}

type MingcObj struct {
	id  uint64
	dep iface.ServiceDep

	warId int32

	bid *entity.BidInfo

	defId int64 // 防守联盟id。开始战斗时，会设置为 npc 联盟 ID
	atkId int64 // 进攻联盟id

	astDefList []int64 // 协防联盟列表
	astAtkList []int64 // 助攻联盟列表

	applyAsts map[int64]bool // 申请协助列表 true:助攻

	scene *McWarScene // 副本场景
}

func newDefaultMingcObj() *MingcObj {
	c := &MingcObj{}
	c.applyAsts = make(map[int64]bool)
	return c
}

func newMingcObj(id uint64, warId int32, hostId, atkId int64, dep iface.ServiceDep, applyAtkEndTime time.Time) *MingcObj {
	c := newDefaultMingcObj()
	c.id = id
	c.warId = warId
	c.dep = dep
	c.defId = hostId
	c.atkId = atkId

	c.bid = entity.NewBidInfo(applyAtkEndTime)

	return c
}

func (c *MingcObj) encode(guildGetter func(id int64) *guildsnapshotdata.GuildSnapshot) *shared_proto.McWarMcProto {
	p := &shared_proto.McWarMcProto{}
	p.Id = u64.Int32(c.id)

	c.bid.WalkBid(func(obj *entity.BidObj) {
		p.ApplyAtkCount++
	})

	atk := guildGetter(c.atkId)
	if atk != nil {
		p.AtkGuild = atk.BasicProto()
	}

	def := guildGetter(c.defId)
	if def != nil {
		p.DefGuild = def.BasicProto()
		p.Country = u64.Int32(def.Country.Id)
	} else {
		p.Country = u64.Int32(c.dep.Datas().GetMingcBaseData(c.id).Country)
	}

	for _, gid := range c.astAtkList {
		g := guildGetter(gid)
		if g != nil {
			p.AstAtkGuild = append(p.AstAtkGuild, g.BasicProto())
		}
	}

	for _, gid := range c.astDefList {
		g := guildGetter(gid)
		if g != nil {
			p.AstDefGuild = append(p.AstDefGuild, g.BasicProto())
		}
	}

	if c.scene != nil {
		p.Ended = c.scene.ended
		p.AtkWin = c.scene.atkWin
	}

	return p
}

func (c *MingcObj) encodeServer() *server_proto.MingcObjProto {
	p := &server_proto.MingcObjProto{}
	p.Id = c.id
	p.DefId = c.defId
	p.AtkId = c.atkId
	p.AstDef = c.astDefList
	p.AstAtk = c.astAtkList
	p.ApplyAst = c.applyAsts

	p.Bid = c.bid.Encode()
	if c.scene != nil {
		p.Scene = c.scene.encodeServer()
	}

	return p
}

func (c *MingcObj) unmarshal(proto *server_proto.MingcObjProto, dep iface.ServiceDep) {
	c.id = proto.Id
	c.dep = dep
	c.defId = proto.DefId
	c.atkId = proto.AtkId
	c.astDefList = proto.AstDef
	c.astAtkList = proto.AstAtk
	if proto.ApplyAst != nil {
		c.applyAsts = proto.ApplyAst
	}
	c.bid = entity.NewBidInfo(time.Time{})
	c.bid.Unmarshal(proto.Bid)

	if proto.Scene != nil {
		c.scene = newDefaultMcWarScene(dep.Datas())
		c.scene.unmershal(proto.Scene, dep)
	}
}

func (c *MingcObj) replyApplyAst(operId, gid int64, mcId uint64, agree bool) bool {
	if isAtk, ok := c.applyAsts[gid]; !ok {
		return false
	} else {
		if isAtk && operId == c.atkId {
			delete(c.applyAsts, gid)
			if agree {
				c.astAtkList = append(c.astAtkList, gid)
			}
			return true
		} else if !isAtk && operId == c.defId {
			delete(c.applyAsts, gid)
			if agree {
				c.astDefList = append(c.astDefList, gid)
			}
			return true
		}
	}

	return false
}

func (c *MingcObj) isAtk(gid int64) bool {
	if c.atkId == gid {
		return true
	}

	for _, a := range c.astAtkList {
		if a == gid {
			return true
		}
	}

	return false
}

func (c *MingcObj) isDef(gid int64) bool {
	if c.defId == gid {
		return true
	}

	for _, a := range c.astDefList {
		if a == gid {
			return true
		}
	}

	return false
}
