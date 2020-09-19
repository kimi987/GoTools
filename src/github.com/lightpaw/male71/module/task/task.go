package task

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/task"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/module/rank/ranklist"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/config/basedata"
)

func NewTaskModule(dep iface.ServiceDep, guildSnapshotService iface.GuildSnapshotService, rankModule iface.RankModule,
	realmService iface.RealmService) *TaskModule {
	m := &TaskModule{}
	m.dep = dep
	m.datas = dep.Datas()
	m.timeService = dep.Time()
	m.guildSnapshotService = guildSnapshotService
	m.heroDataService = dep.HeroData()
	m.heroSnapshotService = dep.HeroSnapshot()
	m.rankModule = rankModule
	m.realmService = realmService

	heromodule.RegisterHeroOnlineListener(m)

	entity.RegisterUpdateTaskProgressFunc(shared_proto.TaskTargetType_TASK_TARGET_GUILD_LEVEL, m.updateGuildLevelTaskType)

	return m

}

//gogen:iface
type TaskModule struct {
	dep                  iface.ServiceDep
	datas                iface.ConfigDatas
	timeService          iface.TimeService
	guildSnapshotService iface.GuildSnapshotService
	heroDataService      iface.HeroDataService
	broadcast            iface.BroadcastService
	rankModule           iface.RankModule
	heroSnapshotService  iface.HeroSnapshotService
	realmService         iface.RealmService
}

func (m *TaskModule) OnHeroOnline(hc iface.HeroController) {

	var homeNpcDatas []*basedata.HomeNpcBaseData
	checkTaskMonster := false
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroTaskList := hero.TaskList()

		var activeDegreeChanged bool
		var hasHomeNpcTask bool
		heroTaskList.WalkAllTask(func(task entity.Task) (endedWalk bool) {
			if task.Progress().UpdateTaskProgress(hero) {
				result.Changed()
				result.Add(task.NewUpdateTaskProgressMsg())

				if task.Type() == entity.ACTIVE_DEGREE_TASK {
					activeDegreeChanged = true
				}
			}
			if !task.Progress().IsCompleted() {
				if task.Progress().TargetType() == shared_proto.TaskTargetType_TASK_TARGET_KILL_HOME_NPC {
					hasHomeNpcTask = true
				}

				if monster := task.Progress().Target().InvasionMonster; monster != nil {
					_, exist := heroTaskList.GetTaskMonster(monster.Id)
					if !exist {
						heroTaskList.SetTaskMonster(monster.Id, false)

						if hero.BaseRegion() != 0 {
							checkTaskMonster = true
						}
					}
				}
			}

			return false
		})

		degree := hero.TaskList().ActiveDegreeTaskList().Degree()
		if hero.HistoryAmount().Amount(server_proto.HistoryAmountType_MaxActiveDegree) < degree {
			hero.HistoryAmount().Set(server_proto.HistoryAmountType_MaxActiveDegree, degree)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACTIVE_DEGREE)
		}

		if activeDegreeChanged {
			heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_ACTIVITY, degree)
		}

		if hero.BaseRegion() != 0 && hasHomeNpcTask {
			heroBaYeStage := hero.TaskList().GetCompletedBaYeStage() + 1
			// 检查一下是不是有该刷的野怪没刷新出来
			for _, data := range m.datas.GetHomeNpcBaseDataArray() {
				if data.BaYeStage > 0 && data.BaYeStage <= heroBaYeStage && !hero.HasCreateHomeNpcBase(data.Id) {
					homeNpcDatas = append(homeNpcDatas, data)
				}
			}
		}

		result.Ok()
	})

	if len(homeNpcDatas) > 0 {
		process, err := m.realmService.GetBigMap().AddHomeNpc(hc, homeNpcDatas)
		if !process || err != nil {
			// 极端情况，把玩家踢下线？
			logrus.WithField("reason", err).Error("登陆时候发现新的HomeNpc需要刷新，但是刷新失败")
		}
	}

	if checkTaskMonster {
		m.realmService.GetBigMap().UpdateHeroRealmInfo(hc.Id(), false, true, false)
	}
}

func (m *TaskModule) updateGuildLevelTaskType(hero *entity.Hero, t *entity.TaskProgress) (newProgress uint64) {
	if hero.GuildId() <= 0 {
		return
	}

	snapshot := m.guildSnapshotService.GetSnapshot(hero.GuildId())
	if snapshot != nil {
		return snapshot.GuildLevel.Level
	}

	return 0
}

var completeMainTask = task.NewS2cNewTaskMsg(0, 0, false, true)

//gogen:iface
func (m *TaskModule) ProcessCollectTaskPrize(proto *task.C2SCollectTaskPrizeProto, hc iface.HeroController) {

	var starRankFunc func()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroTaskList := hero.TaskList()

		taskId := u64.FromInt32(proto.Id)

		t := heroTaskList.Task(taskId)
		if t == nil {
			// 没有找到目标任务
			logrus.Debugf("领取任务奖励，目标任务没找到")
			result.Add(task.ERR_COLLECT_TASK_PRIZE_FAIL_INVALID_ID)
			return
		}

		if !t.Progress().IsCompleted() {
			logrus.Debugf("领取任务奖励，任务还未完成")
			result.Add(task.ERR_COLLECT_TASK_PRIZE_FAIL_NOT_COMPLETE)
			return
		}

		ohctx := heromodule.NewContext(m.dep, operate_type.TaskCollectTaskPrize)
		switch tk := t.(type) {
		case *entity.MainTask:
			// 给奖励
			hctx := ohctx.Copy(operate_type.TaskCollectMainTaskPrize)
			heromodule.AddPrize(hctx, hero, result, tk.Data().Prize, m.timeService.CurrentTime())

			data := tk.Data()

			heroTaskList.CompleteMainTask()

			result.Add(task.NewS2cCollectTaskPrizeMsg(proto.Id))
			//result.Add(data.RemoveTaskMsg)

			// 新任务
			if heroTaskList.MainTask() != nil {
				heroTaskList.MainTask().Progress().UpdateTaskProgress(hero)
				result.Add(heroTaskList.MainTask().NewTaskMsg())
			} else {
				result.Add(completeMainTask)
			}

			for _, b := range data.BranchTask {
				bt := heroTaskList.GetBranchTask(b.Id)
				if bt != nil {
					bt.Progress().UpdateTaskProgress(hero)
					result.Add(bt.NewTaskMsg())
				}
			}

			heromodule.CheckFuncsOpened(hero, result)
		case *entity.BranchTask:
			// 给奖励
			hctx := ohctx.Copy(operate_type.TaskCollectBranchTaskPrize)
			heromodule.AddPrize(hctx, hero, result, tk.Data().Prize, m.timeService.CurrentTime())

			data := tk.Data()

			// 完成任务
			heroTaskList.CompleteBranchTask(tk)

			result.Add(task.NewS2cCollectTaskPrizeMsg(proto.Id))
			//result.Add(data.RemoveTaskMsg)

			// 新任务
			if data.NextTask() != nil {
				bt := heroTaskList.GetBranchTask(data.NextTask().Id)
				if bt != nil {
					bt.Progress().UpdateTaskProgress(hero)
					result.Add(bt.NewTaskMsg())
				}
			}
		case *entity.AchieveTask:
			if tk.IsCollectPrize() {
				logrus.Debugf("领取任务奖励，奖励已经领取")
				result.Add(task.ERR_COLLECT_TASK_PRIZE_FAIL_IS_COLLECTED)
				return
			}

			tk.CollectPrize()

			ctime := m.timeService.CurrentTime()
			// 给奖励
			hctx := ohctx.Copy(operate_type.TaskCollectAchieveTaskPrize)
			heromodule.AddPrize(hctx, hero, result, tk.Data().Prize, ctime)

			// 完成任务
			newTask := heroTaskList.AchieveTaskList().Complete(tk, ctime)

			result.Add(task.NewS2cCollectTaskPrizeMsg(proto.Id))

			// 新任务
			if newTask != nil {
				result.Add(tk.Data().RemoveTaskMsg)
				newTask.Progress().UpdateTaskProgress(hero)
				result.Add(newTask.NewTaskMsg())

				result.Add(task.NewS2cAchieveReachMsg(u64.Int32(newTask.Data().Id), timeutil.Marshal32(newTask.ReachTime())))
			} else {
				result.Add(task.NewS2cAchieveReachMsg(u64.Int32(tk.Data().Id), timeutil.Marshal32(tk.ReachTime())))
			}

			// 设置展示成就
			if uint64(hero.TaskList().AchieveTaskList().ShowCount()) < m.datas.TaskMiscData().MaxShowAchieveCount {
				if hero.TaskList().AchieveTaskList().AddSelectShowAchieveType(tk.Data().AchieveType) {
					result.Add(task.NewS2cChangeSelectShowAchieveMsg(u64.Int32(tk.Data().AchieveType), true))
				}
			}

			heroId := hero.Id()
			star := hero.TaskList().AchieveTaskList().TotalStar()
			starRankFunc = func() {
				// 成就排行榜
				m.rankModule.AddOrUpdateRankObj(ranklist.NewStarTaskRankObj(m.heroSnapshotService.Get, heroId, star, ctime))
			}
		case *entity.BaYeTask:
			if tk.IsCollectPrize() {
				logrus.Debugf("领取任务奖励，奖励已经领取")
				result.Add(task.ERR_COLLECT_TASK_PRIZE_FAIL_IS_COLLECTED)
				return
			}

			// 给奖励
			tk.CollectPrize()
			hctx := ohctx.Copy(operate_type.TaskCollectBaYeTaskPrize)
			heromodule.AddPrize(hctx, hero, result, tk.Data().Prize, m.timeService.CurrentTime())

			// 完成任务
			result.Add(task.NewS2cCollectTaskPrizeMsg(proto.Id))

			if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_BAYE_PRIZE) {
				result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_BAYE_PRIZE)))
			}
		case *entity.BwzlTask:
			if tk.IsCollectPrize() {
				logrus.Debugf("领取任务奖励，奖励已经领取")
				result.Add(task.ERR_COLLECT_TASK_PRIZE_FAIL_IS_COLLECTED)
				return
			}

			day := m.timeService.CurrentTime().Sub(timeutil.DailyTime.PrevTime(hero.CreateTime()))/timeutil.Day + 1
			if day < 0 || tk.Data().Day > uint64(day) {
				logrus.Debugf("领取霸王之路任务奖励，没到可以领取的天数")
				result.Add(task.ERR_COLLECT_TASK_PRIZE_FAIL_CAN_NOT_COLLECT_THIS_DAY)
				return
			}

			// 给奖励
			tk.CollectPrize()
			hctx := ohctx.Copy(operate_type.TaskCollectBwzlTaskPrize)
			heromodule.AddPrize(hctx, hero, result, tk.Data().Prize, m.timeService.CurrentTime())

			// 完成任务
			result.Add(task.NewS2cCollectTaskPrizeMsg(proto.Id))

			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BWZL_COMPLETE)

		case *entity.ActiveDegreeTask:
			// 没有找到目标任务
			logrus.Debugf("领取任务奖励，活跃度任务没有奖励: %d-%s", tk.Data().Id, tk.Data().Name)
			result.Add(task.ERR_COLLECT_TASK_PRIZE_FAIL_INVALID_ID)
			return
		case *entity.ActivityTask:
			if tk.IsCollectPrize() {
				logrus.Debugf("领取任务奖励，奖励已经领取")
				result.Add(task.ERR_COLLECT_TASK_PRIZE_FAIL_IS_COLLECTED)
				return
			}
			// 给奖励
			tk.CollectPrize()
			hctx := ohctx.Copy(operate_type.TaskCollectActivityTaskPrize)
			heromodule.AddPrize(hctx, hero, result, tk.Data().Prize, m.timeService.CurrentTime())
			// 完成任务
			result.Add(task.NewS2cCollectTaskPrizeMsg(proto.Id))

		default:
			// 没有找到目标任务
			logrus.Errorf("领取任务奖励，未处理的任务类型: %+v", tk)
			result.Add(task.ERR_COLLECT_TASK_PRIZE_FAIL_SERVER_ERROR)
			return
		}

		// 历史累计数量
		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_TaskCount)

		// tlog
		var taskType uint64
		switch t.(type) {
		case *entity.BaYeTask:
			taskType = operate_type.TaskTypeBaYe
		case *entity.BwzlTask:
			taskType = operate_type.TaskTypeBwzl
		case *entity.AchieveTask:
			taskType = operate_type.TaskTypeAchieve
		case *entity.ActiveDegreeTask:
			taskType = operate_type.TaskTypeActiveDegree
		case *entity.MainTask:
			taskType = operate_type.TaskTypeMain
		case *entity.BranchTask:
			taskType = operate_type.TaskTypeBranch
		case *entity.ActivityTask:
			taskType = operate_type.TaskTypeActivity
		}

		m.dep.Tlog().TlogTaskFlow(hero, taskType, taskId, operate_type.TaskOperTypeCollect)

		result.Changed()
		result.Ok()
	})

	if starRankFunc != nil {
		starRankFunc()
	}
}

//gogen:iface
func (m *TaskModule) ProcessCollectTaskBoxPrize(proto *task.C2SCollectTaskBoxPrizeProto, hc iface.HeroController) {

	box := m.datas.GetTaskBoxData(u64.FromInt32(proto.Id))

	if box == nil {
		logrus.Debugf("领取任务宝箱，宝箱id无效")
		hc.Send(task.ERR_COLLECT_TASK_BOX_PRIZE_FAIL_INVALID_ID)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.TaskList().CollectedBoxId()+1 < box.Id {
			logrus.Debugf("领取任务宝箱，前面还有未领取的宝箱")
			result.Add(task.ERR_COLLECT_TASK_BOX_PRIZE_FAIL_PREV_NOT_COLLECTED)
			return
		}

		if hero.TaskList().CollectedBoxId() >= box.Id {
			logrus.Debugf("领取任务宝箱，宝箱已领取")
			result.Add(task.ERR_COLLECT_TASK_BOX_PRIZE_FAIL_COLLECTED)
			return
		}

		if hero.HistoryAmount().Amount(server_proto.HistoryAmountType_TaskCount) < box.Count {
			logrus.Debugf("领取任务宝箱，累积完成任务个数不足")
			result.Add(task.ERR_COLLECT_TASK_BOX_PRIZE_FAIL_TASK_NOT_ENOUGH)
			return
		}

		hctx := heromodule.NewContext(m.dep, operate_type.TaskCollectTaskBoxPrize)
		heromodule.AddPrize(hctx, hero, result, box.Prize, m.timeService.CurrentTime())

		hero.TaskList().IncreseCollectedBoxId()
		result.Add(task.NewS2cCollectTaskBoxPrizeMsg(proto.Id))

		result.Changed()
		result.Ok()
	})
}

var noBaYeStageMsgCache = task.NewS2cCollectBaYeStagePrizeMsg([]byte{}).Static()

//gogen:iface c2s_collect_ba_ye_stage_prize
func (m *TaskModule) ProcessCollectBaYeStagePrize(hc iface.HeroController) {

	var homeNpcDatas []*basedata.HomeNpcBaseData
	checkTaskMonster := false
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		taskList := hero.TaskList()
		baYeStage := taskList.BaYeStage()
		if baYeStage == nil {
			logrus.Debugf("领取霸业目标奖励，当前没有霸业目标")
			hc.Send(task.ERR_COLLECT_BA_YE_STAGE_PRIZE_FAIL_NO_BA_YE_STAGE)
			return
		}

		hasTaskNotComplete, hasTaskNotCollectPrize := baYeStage.CanCollectStagePrize()
		if hasTaskNotComplete {
			logrus.Debugf("领取霸业目标奖励，有任务没有完成")
			hc.Send(task.ERR_COLLECT_BA_YE_STAGE_PRIZE_FAIL_HAS_TASK_NOT_COMPLETE)
			return
		}

		if hasTaskNotCollectPrize {
			logrus.Debugf("领取霸业目标奖励，有任务没有领取奖励")
			hc.Send(task.ERR_COLLECT_BA_YE_STAGE_PRIZE_FAIL_HAS_TASK_NOT_COLLECT_PRIZE)
			return
		}

		taskList.CompleteBaYeStage()

		hctx := heromodule.NewContext(m.dep, operate_type.TaskCollectBaYeStagePrize)
		heromodule.AddPrize(hctx, hero, result, baYeStage.Data().Prize, m.timeService.CurrentTime())

		newBaYeStage := taskList.BaYeStage()
		if newBaYeStage == nil {
			// 没有下一个了
			result.Add(noBaYeStageMsgCache)
		} else {
			// 在下面的协议里面发了
			newBaYeStage.Walk(func(t *entity.BaYeTask) (endWalk bool) {
				if !t.Progress().UpdateTaskProgress(hero) {
					if t.Data().Target.InvasionMonster != nil {
						taskList.SetTaskMonster(t.Data().Target.InvasionMonster.Id, false)
						checkTaskMonster = true
					}
				}
				return false
			})

			result.Add(task.NewS2cCollectBaYeStagePrizeMsg(must.Marshal(newBaYeStage.EncodeClient())))
		}

		heromodule.CheckFuncsOpened(hero, result)

		// 刷新玩家主城野怪
		if newBaYeStage != nil {
			for _, data := range m.datas.GetHomeNpcBaseDataArray() {
				if data.BaYeStage > 0 && data.BaYeStage == newBaYeStage.Data().Stage && !hero.HasCreateHomeNpcBase(data.Id) {
					homeNpcDatas = append(homeNpcDatas, data)
				}
			}
		}

		result.Changed()
		result.Ok()
	})

	if len(homeNpcDatas) > 0 {
		process, err := m.realmService.GetBigMap().AddHomeNpc(hc, homeNpcDatas)
		if !process || err != nil {
			// 极端情况，把玩家踢下线？
			logrus.WithField("reason", err).Error("领取霸业奖励，刷新HomeNpc失败")
		}
	}

	if checkTaskMonster {
		m.realmService.GetBigMap().UpdateHeroRealmInfo(hc.Id(), false, true, false)
	}
}

//gogen:iface c2s_collect_active_degree_prize
func (m *TaskModule) ProcessCollectActiveDegreePrize(proto *task.C2SCollectActiveDegreePrizeProto, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		activeDegreeTaskList := hero.TaskList().ActiveDegreeTaskList()

		prizeDataConfig := m.datas.ActiveDegreePrizeData()
		if proto.CollectIndex < 0 || uint64(proto.CollectIndex) >= uint64(len(prizeDataConfig.Array)) {
			logrus.Debugf("领取活跃度奖励，奖励全部领取了")
			hc.Send(task.ERR_COLLECT_ACTIVE_DEGREE_PRIZE_FAIL_COLLECTED)
			return
		}

		if activeDegreeTaskList.IsPrizeCollected(proto.CollectIndex) {
			logrus.Debugf("领取活跃度奖励，奖励领取了")
			hc.Send(task.ERR_COLLECT_ACTIVE_DEGREE_PRIZE_FAIL_COLLECTED)
			return
		}

		prizeData := prizeDataConfig.Array[proto.CollectIndex]
		if prizeData.Degree > activeDegreeTaskList.Degree() {
			logrus.Debugf("领取活跃度奖励，活跃度不够")
			hc.Send(task.ERR_COLLECT_ACTIVE_DEGREE_PRIZE_FAIL_DEGREE_NOT_ENOUGH)
			return
		}

		prize := prizeData.Prize
		if prizeData.Plunder != nil {
			prize = prizeData.Plunder.Try()
		}

		activeDegreeTaskList.CollectPrize(proto.CollectIndex)
		hctx := heromodule.NewContext(m.dep, operate_type.TaskCollectActiveDegreePrize)
		heromodule.AddPrize(hctx, hero, result, prize, m.timeService.CurrentTime())

		result.Add(task.NewS2cCollectActiveDegreePrizeMsg(proto.CollectIndex))

		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_CollectActiveBox)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_COLLECT_ACTIVE_BOX)

		result.Changed()
		result.Ok()
	})
}

//gogen:iface c2s_collect_achieve_star_prize
func (m *TaskModule) ProcessCollectAchieveStarPrize(proto *task.C2SCollectAchieveStarPrizeProto, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		achieveTaskList := hero.TaskList().AchieveTaskList()

		starData := m.datas.GetAchieveTaskStarPrizeData(uint64(proto.StarCount))
		if starData == nil {
			logrus.Debugf("领取成就星级奖励，没找到该星级的奖励")
			result.Add(task.ERR_COLLECT_ACHIEVE_STAR_PRIZE_FAIL_PRIZE_NOT_FOUND)
			return
		}

		if achieveTaskList.IsPrizeCollected(starData.Star) {
			logrus.Debugf("领取成就星级奖励，奖励领取了")
			result.Add(task.ERR_COLLECT_ACHIEVE_STAR_PRIZE_FAIL_COLLECTED)
			return
		}

		if achieveTaskList.TotalStar() < starData.Star {
			logrus.Debugf("领取成就星级奖励，星数不够")
			result.Add(task.ERR_COLLECT_ACHIEVE_STAR_PRIZE_FAIL_STAR_NOT_ENOUGH)
			return
		}

		prize := starData.Prize
		if starData.Plunder != nil {
			prize = starData.Plunder.Try()
		}

		result.Changed()
		result.Ok()

		achieveTaskList.CollectPrize(starData.Star)
		hctx := heromodule.NewContext(m.dep, operate_type.TaskCollectAchieveStarPrize)
		heromodule.AddPrize(hctx, hero, result, prize, m.timeService.CurrentTime())

		result.Add(task.NewS2cCollectAchieveStarPrizeMsg(proto.StarCount))

		// 系统广播
		if d := hctx.BroadcastHelp().TaskChengJiuLevel; d != nil {
			hctx.AddBroadcast(d, hero, result, 0, starData.Star, func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, starData.Star)
				return text
			})
		}
	})
}

//gogen:iface c2s_change_select_show_achieve
func (m *TaskModule) ProcessChangeSelectShowAchieve(proto *task.C2SChangeSelectShowAchieveProto, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		achieveTaskList := hero.TaskList().AchieveTaskList()

		achieveType := u64.FromInt32(proto.AchieveType)

		if proto.AddOrRemove {
			// 展示
			if uint64(achieveTaskList.ShowCount()) >= m.datas.TaskMiscData().MaxShowAchieveCount {
				logrus.Debugf("当前已经展示很多了")
				result.Add(task.ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_SHOW_TO_MANY)
				return
			}

			if achieveTaskList.IsAchieveTypeSelectShow(achieveType) {
				logrus.Debugf("该成就当前正在展示中")
				result.Add(task.ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_SHOWING)
				return
			}

			achieveTask := achieveTaskList.GetTaskByAchieveType(achieveType)
			if achieveTask == nil {
				logrus.Debugf("没有完成该成就，无法展示")
				result.Add(task.ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_NOT_FOUND)
				return
			}

			if achieveTask.Data().PrevTask == nil && !achieveTask.IsCollectPrize() {
				logrus.Debugf("没有完成该成就")
				result.Add(task.ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_LOCKED)
				return
			}

			achieveTaskList.AddSelectShowAchieveType(achieveType)
		} else {
			if !achieveTaskList.IsAchieveTypeSelectShow(achieveType) {
				logrus.Debugf("该成就当前没在展示中")
				result.Add(task.ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_NOT_SHOWING)
				return
			}

			// 不展示
			achieveTaskList.RemoveSelectShowAchieveType(achieveType)
		}

		result.Changed()
		result.Ok()

		result.Add(task.NewS2cChangeSelectShowAchieveMsg(proto.AchieveType, proto.AddOrRemove))
	})
}

//gogen:iface c2s_collect_bwzl_prize
func (m *TaskModule) ProcessCollectBwzlPrize(proto *task.C2SCollectBwzlPrizeProto, hc iface.HeroController) {
	bwzlData := m.datas.BwzlPrizeData().Get(u64.FromInt32(proto.CompleteCount))
	if bwzlData == nil {
		logrus.Debugf("霸王之路任务奖励没找到")
		hc.Send(task.ERR_COLLECT_BWZL_PRIZE_FAIL_NOT_FOUND)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		bwzlTaskList := hero.TaskList().BwzlTaskList()

		if bwzlTaskList.IsCollectedPrizes(bwzlData) {
			logrus.Debugf("霸王之路任务奖励已经领取")
			hc.Send(task.ERR_COLLECT_BWZL_PRIZE_FAIL_COLLECTED)
			return
		}

		if bwzlTaskList.CollectPrizeCount() < bwzlData.CollectPrizeTaskCount {
			logrus.Debugf("霸王之路任务奖励无法领取")
			hc.Send(task.ERR_COLLECT_BWZL_PRIZE_FAIL_NOT_REACH)
			return
		}

		result.Changed()
		result.Ok()

		bwzlTaskList.CollectedPrize(bwzlData)

		hctx := heromodule.NewContext(m.dep, operate_type.TaskCollectBwzlTaskStagePrize)
		heromodule.AddPrize(hctx, hero, result, bwzlData.Prize, m.timeService.CurrentTime())

		result.Add(bwzlData.CollectPrizeMsg)

		if d := hctx.BroadcastHelp().TaskBwzlPrize; d != nil {
			hctx.AddBroadcast(d, hero, result, 0, bwzlData.CollectPrizeTaskCount, func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				return text
			})
		}
	})
}

//gogen:iface c2s_view_other_achieve_task_list
func (m *TaskModule) ProcessViewOtherAchieveTaskList(proto *task.C2SViewOtherAchieveTaskListProto, hc iface.HeroController) {
	id, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("请求成就列表，解析id失败")
		hc.Send(task.ERR_VIEW_OTHER_ACHIEVE_TASK_LIST_FAIL_INVALID_ID)
		return
	}

	if id == hc.Id() {
		logrus.Debugf("请求成就列表，自己的请不要走这里")
		hc.Send(task.ERR_VIEW_OTHER_ACHIEVE_TASK_LIST_FAIL_CANT_REQUEST_SELF)
		return
	}

	m.heroDataService.Func(id, func(hero *entity.Hero, err error) (heroChanged bool) {
		if err != nil {
			logrus.WithError(err).Debugf("请求他人的成就列表出错")
			hc.Send(task.ERR_VIEW_OTHER_ACHIEVE_TASK_LIST_FAIL_SERVER_ERROR)
			return
		}

		achieveTaskList := hero.TaskList().AchieveTaskList()

		hc.Send(task.NewS2cViewOtherAchieveTaskListMsg(proto.Id, achieveTaskList.EncodeOther()))

		return
	})
}

//gogen:iface
func (m *TaskModule) ProcessGetUpgradeTitleFightAmount(proto *task.C2SGetTroopTitleFightAmountProto, hc iface.HeroController) {

	var fightAmount uint64
	if targetTitle := m.datas.GetTitleData(u64.FromInt32(proto.TitleId)); targetTitle != nil {
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

			currentTitle := hero.TaskList().GetTitleData()
			if currentTitle != nil && currentTitle.Id >= targetTitle.Id {
				return
			}

			troop := hero.GetTroopByIndex(u64.FromInt32(proto.TroopIndex - 1))
			if troop == nil {
				return
			}

			var captainFightAmount []uint64
			for _, pos := range troop.Pos() {
				c := pos.Captain()
				if c != nil {
					captainStat := c.GetPreviewStat(targetTitle)
					fm := captainStat.FightAmount(captainStat.SoldierCapcity, c.GetSpellFightAmountCoef())
					captainFightAmount = append(captainFightAmount, fm)
				}
			}
			newFightAmount := data.TroopFightAmount(captainFightAmount...)

			fightAmount = u64.Sub(newFightAmount, troop.FullFightAmount())

			result.Ok()
		})
	}

	hc.Send(task.NewS2cGetTroopTitleFightAmountMsg(proto.TroopIndex, proto.TitleId, u64.Int32(fightAmount)))

}

//gogen:iface c2s_upgrade_title
func (m *TaskModule) ProcessUpgradeTitle(hc iface.HeroController) {

	var countryId uint64
	var newVoteCount uint64
	var newTitle uint64
	hasError := hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		ttl := hero.TaskList().TitleTaskList()

		doingData := ttl.GetDoingTitleData()
		if doingData == nil {
			logrus.Debug("升级任务称号，已经是最高等级")
			result.Add(task.ERR_UPGRADE_TITLE_FAIL_MAX_LEVEL)
			return
		}

		if !ttl.IsAllTaskComplete() {
			logrus.Debug("升级任务称号，任务未完成")
			result.Add(task.ERR_UPGRADE_TITLE_FAIL_TASK_NOT_COMPLETE)
			return
		}

		hctx := heromodule.NewContext(m.dep, operate_type.TaskUpgradeTitle)
		if doingData.TitleCost != nil && !heromodule.TryReduceCost(hctx, hero, result, doingData.TitleCost) {
			logrus.Debug("升级任务称号，贡品不足")
			result.Add(task.ERR_UPGRADE_TITLE_FAIL_NOT_ENOUGH_COST)
			return
		}

		//检测是否广播橙色及以上称号
		if d := hctx.BroadcastHelp().Title; d != nil {
			hctx.AddBroadcast(d, hero, result, 0, uint64(ttl.GetDoingTitleData().Quality), func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyText, ttl.GetDoingTitleData().Name)
				return text
			})
		}

		ttl.UpgradeTitleData()
		result.Add(doingData.GetCompleteMsg())
		newTitle = ttl.TitleId()

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

		countryId = hero.CountryId()
		if t := hero.TaskList().GetTitleData(); t != nil {
			newVoteCount = t.CountryChangeNameVoteCount
		}

		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_TITLE)

		result.Changed()
		result.Ok()
	})

	if !hasError {
		// 更新野外主城称号
		m.realmService.GetBigMap().ChangeTitle(hc.Id(), newTitle)
		m.dep.Country().AfterUpgradeTitle(hc.Id(), countryId, newVoteCount)
	}

}

//gogen:iface
func (m *TaskModule) ProcessCompleteBoolTask(proto *task.C2SCompleteBoolTaskProto, hc iface.HeroController) {

	boolType := shared_proto.HeroBoolType(proto.BoolType)
	switch boolType {
	case shared_proto.HeroBoolType_BOOL_VIEW_MAP:
	default:
		logrus.Warn("完成bool类型任务，客户端发送的Bool类型不支持")
		hc.Send(task.NewS2cCompleteBoolTaskMsg(proto.BoolType))
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.Bools().TrySet(boolType) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(boolType)))
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BOOL)
		}

		result.Ok()
	})
}
