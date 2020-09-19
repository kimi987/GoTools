package gm

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/gm"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/xuanyuan"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/timeutil"
	"time"
)

func (m *GmModule) newResetGmGroup() *gm_group {
	return &gm_group{
		tab: "重置",
		handler: []*gm_handler{
			newHeroStringHandler("每日重置", "", m.resetDaily),
			newHeroStringHandler("每日0点重置", "", m.resetDailyZero),
			newHeroStringHandler("每日22点重置", "", m.resetDailyMc),
			newStringHandler("重置联盟", "", m.resetGuildDaily),
			newStringHandler("清除联盟重置时间", "", m.setGuildResetTime),
			newStringHandler("重置百战", "", m.resetBaiZhan),
			newStringHandler("重置百战挑战次数", "", m.resetBaiZhanChallengeTimes),
			newStringHandler("重置轩辕会武", "", m.resetXuanyuan),
			newStringHandler("重置轩辕会武挑战次数", "", m.resetXuanyuanChallengeTimes),
			newStringHandler("重置四季", "", m.resetSeason),
			newStringHandler("给10次酒馆次数", "", m.addJiuGuanTimes),
			//newStringHandler("重置修炼馆等级", "", m.resetXiuLianGuanLevel),
			newStringHandler("给10次征兵次数", "", m.addZhengBingTimes),
			newStringHandler("刷新宝藏野怪", "", m.refreshBaoZangNpc),
		},
	}
}

func (m *GmModule) resetDaily(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	ctime := m.time.CurrentTime()
	hero.ResetDaily(ctime, m.datas)
	result.Add(misc.RESET_DAILY_S2C)

	hero.TaskList().ActiveDegreeTaskList().Walk(func(taskId uint64, task *entity.ActiveDegreeTask) bool {
		result.Add(task.NewUpdateTaskProgressMsg())
		return false
	})
}

func (m *GmModule) resetDailyZero(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	ctime := m.time.CurrentTime()
	hero.ResetDailyZero(ctime, m.datas)
	result.Add(misc.RESET_DAILY_ZERO_S2C)
	result.Changed()

	// vip
	heromodule.GmVipResetDailyZero(m.hctx, hero, result, m.datas, ctime)

	hero.TaskList().ActiveDegreeTaskList().Walk(func(taskId uint64, task *entity.ActiveDegreeTask) bool {
		result.Add(task.NewUpdateTaskProgressMsg())
		return false
	})
}

func (m *GmModule) resetDailyMc(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	ctime := m.time.CurrentTime()
	hero.ResetDailyMc(ctime, m.datas)
	result.Add(misc.RESET_DAILY_MC_S2C)
	result.Changed()

	hero.TaskList().ActiveDegreeTaskList().Walk(func(taskId uint64, task *entity.ActiveDegreeTask) bool {
		result.Add(task.NewUpdateTaskProgressMsg())
		return false
	})
}

func (m *GmModule) resetGuildDaily(input string, hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		hc.Send(gm.NewS2cGmMsg("没有加入联盟"))
		return
	}

	resetTime := m.tick.GetDailyTickTime().GetPrevTickTime()
	m.sharedGuildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		g.GmResetDaily(m.datas.GuildConfig().ContributionDay, resetTime)
	})
}

func (m *GmModule) setGuildResetTime(input string, hc iface.HeroController) {
	resetTime := m.tick.GetDailyTickTime().GetPrevTickTime()
	m.sharedGuildService.Func(func(guilds sharedguilddata.Guilds) {
		guilds.Walk(func(g *sharedguilddata.Guild) {
			g.GmSetResetTime(resetTime)
		})
	})
}

func (m *GmModule) resetBaiZhan(input string, hc iface.HeroController) {
	m.modules.BaiZhanModule().GmResetDaily()
}

func (m *GmModule) resetBaiZhanChallengeTimes(input string, hc iface.HeroController) {
	m.modules.BaiZhanModule().GmResetChallengeTimes(hc.Id())
}

func (m *GmModule) resetXuanyuan(input string, hc iface.HeroController) {
	m.modules.XuanyuanModule().GmReset()
}

func (m *GmModule) resetXuanyuanChallengeTimes(input string, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.Xuanyuan().Reset()
		result.Add(xuanyuan.RESET_S2C)
	})
}

func (m *GmModule) resetSeason(input string, hc iface.HeroController) {
	ctime := m.time.CurrentTime()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.SetSeasonResetTime(ctime.Add(-24 * time.Hour))
	})

	curSeason := m.seasonService.SeasonByTime(ctime)
	heromodule.TryResetSeason(hc, ctime, timeutil.DailyTime.PrevTime(ctime), curSeason)
}

func (m *GmModule) refreshBaoZangNpc(input string, hc iface.HeroController) {
	m.realmService.GetBigMap().GmRefreshBaoZangNpc()
}

// 给10次征兵次数
func (m *GmModule) addZhengBingTimes(input string, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		mil := hero.Military()
		mil.RecruitTimes().AddTimes(10, m.time.CurrentTime())

		result.Changed()
		result.Ok()

		result.Add(military.NewS2cRecruitSoldierTimesChangedMsg(mil.RecruitTimes().StartTimeUnix32()))
	})
}

func (m *GmModule) addJiuGuanTimes(input string, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		mil := hero.Military()
		mil.JiuGuanTimes().SetTimes(0)
		mil.JiuGuanTimes().SetNextTime(time.Time{})

		result.Changed()
		result.Ok()

		result.Add(military.NewS2cJiuGuanTimesChangedMsg(0, 0))
	})
}

//func (m *GmModule) resetXiuLianGuanLevel(input string, hc iface.HeroController) {
//	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
//		training := hero.Military().GetTraining(hero.LevelData())
//
//		ctime := m.time.CurrentTime()
//
//		minLevelData := m.datas.TrainingLevelData().MinKeyData
//
//		for idx, t := range training {
//			// 结算经验
//			if t.Captain() > 0 {
//				// 增加经验
//				toAddExp := t.GetExp(ctime)
//				heromodule.AddCaptainExp(hero, result, t.Captain(), toAddExp)
//			}
//
//			// 设置修炼位等级
//			t.SetData(minLevelData, ctime)
//			result.Add(military.NewS2cUpgradeTrainingMsg(int32(idx), u64.Int32(t.Data().Level)))
//
//			// 更新修炼位产出
//			result.Add(military.NewS2cUpdateTrainingOutputMsg([]int32{int32(idx)}, []int32{u64.Int32(t.Output())}, []int32{u64.Int32(t.Capcity())}, []int32{u64.Int32(t.GetExp(ctime))}))
//
//			// 更新任务进度
//			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_TRAINING_LEVEL_COUNT)
//		}
//
//		result.Changed()
//		result.Ok()
//	})
//}
