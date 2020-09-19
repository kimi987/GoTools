package gm

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/task"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
)

func (m *GmModule) newTaskGmGroup() *gm_group {
	return &gm_group{
		tab: "任务",
		handler: []*gm_handler{
			newHeroStringHandler("完成主线任务", "", m.completeMainTask),
			newHeroIntHandler("设置主线任务", "1", m.setMainTask),
			newHeroStringHandler("完成支线任务", "", m.completeBranchTask),
			newHeroStringHandler("完成霸业目标", "", m.completeBaYeTarget),
			newStringHandler("完成霸业目标并领取", "", m.completeBaYeTargetAndCollect),
			newHeroStringHandler("完成活跃度任务", "", m.completeActiveDegreeTask),
			newHeroStringHandler("完成成就任务", "", m.completeAchieveTask),
			newHeroStringHandler("更新任务", "", m.updateTask),
			newHeroStringHandler("完成霸王之路", "", m.completeBawangTask),
			newHeroStringHandler("完成爵位任务", "", m.completeTitleTask),
			newHeroStringHandler("提升爵位", "", m.upgradeTitleTask),
		},
	}
}

func (m *GmModule) completeMainTask(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	if t := hero.TaskList().MainTask(); t != nil && !t.Progress().IsCompleted() {
		t.Progress().GmComplete()
		hc.Send(t.NewUpdateTaskProgressMsg())
	}
}

func (m *GmModule) setMainTask(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	task := m.datas.GetMainTaskData(u64.FromInt64(amount))
	if task == nil {
		return
	}

	hero.TaskList().GmSetMainTaskData(task)

	mainTask := hero.TaskList().MainTask()
	if mainTask != nil {
		mainTask.Progress().UpdateTaskProgress(hero)
		result.Add(mainTask.NewTaskMsg())
	}
}

func (m *GmModule) completeBranchTask(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	hero.TaskList().WalkBranchTask(func(taskId uint64, t *entity.BranchTask) bool {
		if !t.Progress().IsCompleted() {
			t.Progress().GmComplete()
			hc.Send(t.NewUpdateTaskProgressMsg())
		}
		return false
	})
}

func (m *GmModule) completeBaYeTarget(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	if hero.TaskList().BaYeStage() != nil {
		hero.TaskList().BaYeStage().Walk(func(t *entity.BaYeTask) (endWalk bool) {
			if !t.Progress().IsCompleted() {
				t.Progress().GmComplete()
				hc.Send(t.NewUpdateTaskProgressMsg())
			}
			return false
		})
	}
}

func (m *GmModule) completeBaYeTargetAndCollectTo(toStage int64, hc iface.HeroController) {

	var heroBayeStage uint64
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		heroBayeStage = hero.TaskList().GetCompletedBaYeStage()
		return false
	})

	for i := int64(heroBayeStage); i < toStage; i++ {
		m.completeBaYeTargetAndCollect("", hc)
	}

}

func (m *GmModule) completeBaYeTargetAndCollect(input string, hc iface.HeroController) {

	var taskIds []uint64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.TaskList().BaYeStage() != nil {
			hero.TaskList().BaYeStage().Walk(func(t *entity.BaYeTask) (endWalk bool) {
				if !t.Progress().IsCompleted() {
					t.Progress().GmComplete()
					hc.Send(t.NewUpdateTaskProgressMsg())
				}

				if !t.IsCollectPrize() {
					taskIds = append(taskIds, t.Data().Id)
				}
				return false
			})
		}
	})

	for _, tid := range taskIds {
		m.modules.TaskModule().(interface {
			ProcessCollectTaskPrize(*task.C2SCollectTaskPrizeProto, iface.HeroController)
		}).ProcessCollectTaskPrize(&task.C2SCollectTaskPrizeProto{
			Id: u64.Int32(tid),
		}, hc)
	}

	m.modules.TaskModule().(interface {
		ProcessCollectBaYeStagePrize(iface.HeroController)
	}).ProcessCollectBaYeStagePrize(hc)
}

func (m *GmModule) completeActiveDegreeTask(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	hero.TaskList().ActiveDegreeTaskList().Walk(func(taskId uint64, t *entity.ActiveDegreeTask) (endWalk bool) {
		if !t.Progress().IsCompleted() {
			t.Progress().GmComplete()
			hc.Send(t.NewUpdateTaskProgressMsg())
		}
		return
	})
}

func (m *GmModule) completeAchieveTask(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	hero.TaskList().AchieveTaskList().Walk(func(taskId uint64, t *entity.AchieveTask) (endWalk bool) {
		if !t.Progress().IsCompleted() {
			t.Progress().GmComplete()
			hc.Send(t.NewUpdateTaskProgressMsg())
		}
		return
	})

	heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_ACTIVITY, hero.TaskList().ActiveDegreeTaskList().Degree())
}

func (m *GmModule) completeBawangTask(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	day := uint64(m.time.CurrentTime().Sub(timeutil.DailyTime.PrevTime(hero.CreateTime()))/timeutil.Day + 1)

	hero.TaskList().BwzlTaskList().Walk(func(taskId uint64, t *entity.BwzlTask) (endWalk bool) {
		if t.Data().Day <= day {
			if !t.Progress().IsCompleted() {
				t.Progress().GmComplete()
				hc.Send(t.NewUpdateTaskProgressMsg())
			}
		}
		return
	})

	heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_ACTIVITY, hero.TaskList().ActiveDegreeTaskList().Degree())
}

func (m *GmModule) completeTitleTask(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	hero.TaskList().TitleTaskList().Walk(func(taskId uint64, t *entity.TitleTask) (endWalk bool) {
		if !t.Progress().IsCompleted() {
			t.Progress().GmComplete()
			hc.Send(t.NewUpdateTaskProgressMsg())
		}
		return
	})

	heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_ACTIVITY, hero.TaskList().ActiveDegreeTaskList().Degree())
}

func (m *GmModule) upgradeTitleTask(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	ttl := hero.TaskList().TitleTaskList()

	doingData := ttl.GetDoingTitleData()

	if doingData == nil {
		logrus.Debug("gm升级任务称号，已经是最高等级")
		return
	}

	ttl.UpgradeTitleData()
	result.Add(doingData.GetCompleteMsg())

	// 更新武将
	for _, captain := range hero.Military().Captains() {
		captain.CalculateProperties()
		result.Add(captain.NewUpdateCaptainStatMsg())
	}
	heromodule.UpdateAllTroopFightAmount(hero, result)

	ttl.Walk(func(taskId uint64, task *entity.TitleTask) (endWalk bool) {
		task.Progress().UpdateTaskProgress(hero)
		result.Add(task.NewUpdateTaskProgressMsg())
		return false
	})

	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_TITLE)
}

func (m *GmModule) updateTask(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {

	if t := hero.TaskList().MainTask(); t != nil {
		logrus.Debugf("主线任务 %s, 类型: %s，完成：%v 进度：%v/%v", t.Data().Name, t.Data().Target.Type, t.Progress().IsCompleted(), t.Progress().GetProgress(), t.Data().Target.TotalProgress)
	}

	hero.TaskList().WalkBranchTask(func(taskId uint64, t *entity.BranchTask) bool {
		logrus.Debugf("支线任务 %s, 类型: %s，完成：%v 进度：%v/%v", t.Data().Name, t.Data().Target.Type, t.Progress().IsCompleted(), t.Progress().GetProgress(), t.Data().Target.TotalProgress)
		return false
	})

	if hero.TaskList().BaYeStage() != nil {
		hero.TaskList().BaYeStage().Walk(func(t *entity.BaYeTask) bool {
			logrus.Debugf("霸业任务 %s, 类型: %s，完成：%v 进度：%v/%v", t.Data().Name, t.Data().Target.Type, t.Progress().IsCompleted(), t.Progress().GetProgress(), t.Data().Target.TotalProgress)
			return false
		})
	}

	if t := hero.TaskList().MainTask(); t != nil {
		heromodule.UpdateTaskProgress(hero, result, t.Data().Target.Type)
		logrus.Debugf("主线任务（更新后） %s, 类型: %s，完成：%v 进度：%v/%v", t.Data().Name, t.Data().Target.Type, t.Progress().IsCompleted(), t.Progress().GetProgress(), t.Data().Target.TotalProgress)
	}

	hero.TaskList().WalkBranchTask(func(taskId uint64, t *entity.BranchTask) bool {
		heromodule.UpdateTaskProgress(hero, result, t.Data().Target.Type)
		logrus.Debugf("支线任务（更新后） %s, 类型: %s，完成：%v 进度：%v/%v", t.Data().Name, t.Data().Target.Type, t.Progress().IsCompleted(), t.Progress().GetProgress(), t.Data().Target.TotalProgress)
		return false
	})

	if hero.TaskList().BaYeStage() != nil {
		hero.TaskList().BaYeStage().Walk(func(t *entity.BaYeTask) bool {
			heromodule.UpdateTaskProgress(hero, result, t.Data().Target.Type)
			logrus.Debugf("霸业任务（更新后） %s, 类型: %s，完成：%v 进度：%v/%v", t.Data().Name, t.Data().Target.Type, t.Progress().IsCompleted(), t.Progress().GetProgress(), t.Data().Target.TotalProgress)
			return false
		})
	}
}
