package gm

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/u64"
)

func (m *GmModule) newMingcWarGmGroup() *gm_group {
	return &gm_group{
		tab: "名城战",
		handler: []*gm_handler{
			newIntHandler("联盟加虎符（负数减）", "1000", m.addGuildHufu),
			newIntHandler("名城归属自己", "1", m.setMingcHost),
			newIntHandler("重置名城归属", "1", m.resetMingcHost),
			newStringHandler("新开名城战", "", m.newMcWar),
			newStringHandler("开启申请攻打", "", m.startApplyAtk),
			newStringHandler("开启申请助攻", "", m.startApplyAst),
			newStringHandler("开始战斗", "", m.startFight),
			newStringHandler("结束战斗", "", m.startFightEnd),
			newIntHandler("申请名城攻方(申请攻打阶段有效)", "1", m.setSelfAtkGuild),
			newIntHandler("成为名城守方", "1", m.setSelfDefGuild),
			newIntHandler("成为名城助攻", "1", m.addSelfAstAtkGuild),
			newIntHandler("成为名城协防", "1", m.addSelfAstDefGuild),
			newIntHandler("Id:X城,攻方去死", "1", m.atkDie),
			newIntHandler("Id:X城,守方去死", "1", m.defDie),
		},
	}
}

func (m *GmModule) setMingcHost(mcId int64, hc iface.HeroController) {
	if gid, ok := hc.LockGetGuildId(); ok {
		if gid <= 0 || m.sharedGuildService.GetSnapshot(gid) == nil {
			return
		}

		if mc := m.mingcService.Mingc(u64.FromInt64(mcId)); mc != nil {
			mc.SetHostGuildId(gid)
			m.mingcService.UpdateMsg()
			m.mingcWarService.GmSetDefGuild(u64.FromInt64(mcId), gid)
			m.dep.World().Broadcast(mingc_war.NewS2cMingcHostUpdateNoticeMsg(i64.Int32(mcId), m.dep.GuildSnapshot().GetGuildBasicProto(gid)))
		}
	}
}

func (m *GmModule) resetMingcHost(mcId int64, hc iface.HeroController) {
	if mc := m.mingcService.Mingc(u64.FromInt64(mcId)); mc != nil {
		mc.SetHostGuildId(0)
		m.mingcService.UpdateMsg()
		m.mingcWarService.GmSetDefGuild(u64.FromInt64(mcId), 0)

		m.dep.World().Broadcast(mingc_war.NewS2cMingcHostUpdateNoticeMsg(i64.Int32(mcId), &shared_proto.GuildBasicProto{}))
	}
}

func (m *GmModule) newMcWar(input string, hc iface.HeroController) {
	m.mingcWarService.GmNewMingcWar()
}

func (m *GmModule) startApplyAtk(input string, hc iface.HeroController) {
	m.mingcWarService.GmChangeStage(shared_proto.MingcWarState_MC_T_APPLY_ATK, hc, m.time.CurrentTime())
}

func (m *GmModule) startApplyAst(input string, hc iface.HeroController) {
	m.mingcWarService.GmChangeStage(shared_proto.MingcWarState_MC_T_APPLY_AST, hc, m.time.CurrentTime())
}

func (m *GmModule) startFight(input string, hc iface.HeroController) {
	m.mingcWarService.GmChangeStage(shared_proto.MingcWarState_MC_T_FIGHT, hc, m.time.CurrentTime())
}

func (m *GmModule) startFightEnd(input string, hc iface.HeroController) {
	m.mingcWarService.GmChangeStage(shared_proto.MingcWarState_MC_T_FIGHT_END, hc, m.time.CurrentTime())
}

func (m *GmModule) setAtkGuild(input int64, hc iface.HeroController) {
	m.mingcWarService.GmApplyAtkGuild(1, input)
}

func (m *GmModule) setDefGuild(input int64, hc iface.HeroController) {
	m.mingcWarService.GmSetDefGuild(1, input)
}

func (m *GmModule) addAstAtkGuild(input int64, hc iface.HeroController) {
	m.mingcWarService.GmSetAstAtkGuild(1, input)
}

func (m *GmModule) addAstDefGuild(input int64, hc iface.HeroController) {
	m.mingcWarService.GmSetAstAtkGuild(1, input)
}

func (m *GmModule) setSelfAtkGuild(input int64, hc iface.HeroController) {
	gid, _ := hc.LockGetGuildId()
	if gid <= 0 {
		return
	}
	m.mingcWarService.GmApplyAtkGuild(u64.FromInt64(input), gid)
}

func (m *GmModule) setSelfDefGuild(input int64, hc iface.HeroController) {
	gid, _ := hc.LockGetGuildId()
	if gid <= 0 {
		return
	}
	m.mingcWarService.GmSetDefGuild(u64.FromInt64(input), gid)
}

func (m *GmModule) addSelfAstAtkGuild(input int64, hc iface.HeroController) {
	gid, _ := hc.LockGetGuildId()
	if gid <= 0 {
		return
	}
	m.mingcWarService.GmSetAstAtkGuild(u64.FromInt64(input), gid)
}

func (m *GmModule) addSelfAstDefGuild(input int64, hc iface.HeroController) {
	gid, _ := hc.LockGetGuildId()
	if gid <= 0 {
		return
	}
	m.mingcWarService.GmSetAstDefGuild(u64.FromInt64(input), gid)
}

func (m *GmModule) addGuildHufu(amount int64, hc iface.HeroController) {
	gid, _ := hc.LockGetGuildId()
	if gid <= 0 {
		return
	}
	m.sharedGuildService.FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		var newAmount uint64
		if amount < 0 {
			newAmount = u64.Sub(g.GetHufu(), u64.FromInt64(i64.Abs(amount)))
		} else {
			newAmount = g.GetHufu() + u64.FromInt64(amount)
		}
		g.SetHufu(newAmount)
	})
	m.sharedGuildService.ClearSelfGuildMsgCache(gid)
}

func (m *GmModule) joinFight(mcId int64, hc iface.HeroController) {
	var caps []uint64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		var size int
		for _, c := range hero.Military().Captains() {
			if size >= 5 {
				break
			}
			caps = append(caps, c.Id())
		}
	})

	xIndex := make([]int32, len(caps))
	succMsg, errMsg := m.mingcWarService.JoinFight(hc, u64.FromInt64(mcId), caps, xIndex)
	if errMsg != nil {
		hc.Send(errMsg)
	} else {
		hc.Send(succMsg)
	}
}

func (m *GmModule) quitFight(mcId int64, hc iface.HeroController) {
	succMsg, errMsg := m.mingcWarService.QuitFight(hc.Id())
	if errMsg != nil {
		hc.Send(errMsg)
	} else {
		hc.Send(succMsg)
	}
}

func (m *GmModule) sceneMove(i int64, hc iface.HeroController) {

	mcId, joined := m.mingcWarService.JoiningFightMingc(hc.Id())
	if !joined {
		hc.Send(mingc_war.ERR_SCENE_MOVE_FAIL_NOT_IN_SCENE)
		return
	}
	pos := cb.XYCube(int(i/100), int(i%100))
	dests := m.datas.GetMingcWarSceneData(mcId).GetDest(pos)
	if len(dests) <= 0 {
		hc.Send(mingc_war.ERR_SCENE_MOVE_FAIL_ALREADY_ON_DEST_POS)
		return
	}
	succMsg, errMsg := m.mingcWarService.SceneMove(hc.Id(), dests[0])
	if errMsg != nil {
		hc.Send(errMsg)
	} else {
		hc.Send(succMsg)
	}
}

func (m *GmModule) viewScene(mcId int64, hc iface.HeroController) {
	succMsg, errMsg := m.mingcWarService.ViewMcWarSceneMsg(u64.FromInt64(mcId))
	if errMsg != nil {
		hc.Send(errMsg)
	} else {
		hc.Send(succMsg)
	}
}

func (m *GmModule) atkDie(mcId int64, hc iface.HeroController) {
	m.mingcWarService.GmCampFail(u64.FromInt64(mcId), true)
}

func (m *GmModule) defDie(mcId int64, hc iface.HeroController) {
	m.mingcWarService.GmCampFail(u64.FromInt64(mcId), false)
}
