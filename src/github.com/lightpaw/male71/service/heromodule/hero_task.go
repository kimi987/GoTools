package heromodule

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/gamelogs"
)

func UpdateTaskProgress(hero *entity.Hero, result herolock.LockResult, targetType shared_proto.TaskTargetType) {

	var activeDegreeChanged bool
	hero.TaskList().WalkAllTask(func(task entity.Task) (endedWalk bool) {
		if task.Progress().UpdateTaskTypeProgress(targetType, hero) {
			result.Changed()
			result.Add(task.NewUpdateTaskProgressMsg())

			if task.Type() == entity.ACTIVE_DEGREE_TASK {
				activeDegreeChanged = true
			}

			if task.Progress().IsCompleted() {
				gamelogs.CompleteTaskLog(hero.Pid(), hero.Sid(), hero.Id(), uint64(task.Type()), task.Id())
			}
		}
		return false
	})

	if activeDegreeChanged {
		degree := hero.TaskList().ActiveDegreeTaskList().Degree()
		OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_ACTIVITY, degree)

		if hero.HistoryAmount().Amount(server_proto.HistoryAmountType_MaxActiveDegree) < degree {
			hero.HistoryAmount().Set(server_proto.HistoryAmountType_MaxActiveDegree, degree)
			UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACTIVE_DEGREE)
		}
	}
}

func UpdateTaskProgressWithFunc(hero *entity.Hero, result herolock.LockResult, targetType shared_proto.TaskTargetType, updateFunc func(hero *entity.Hero, t *entity.TaskProgress) (newProgress uint64)) {

	var activeDegreeChanged bool
	hero.TaskList().WalkAllTask(func(task entity.Task) (endedWalk bool) {
		if task.Progress().UpdateTaskTypeProgressWithFunc(targetType, hero, updateFunc) {
			result.Changed()
			result.Add(task.NewUpdateTaskProgressMsg())

			if task.Type() == entity.ACTIVE_DEGREE_TASK {
				activeDegreeChanged = true
			}

			if task.Progress().IsCompleted() {
				gamelogs.CompleteTaskLog(hero.Pid(), hero.Sid(), hero.Id(), uint64(task.Type()), task.Id())
			}
		}
		return false
	})

	if activeDegreeChanged {
		degree := hero.TaskList().ActiveDegreeTaskList().Degree()
		OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_ACTIVITY, degree)

		if hero.HistoryAmount().Amount(server_proto.HistoryAmountType_MaxActiveDegree) < degree {
			hero.HistoryAmount().Set(server_proto.HistoryAmountType_MaxActiveDegree, degree)
			UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACTIVE_DEGREE)
		}
	}
}

func IncreTaskProgressOne(hero *entity.Hero, result herolock.LockResult, targetType shared_proto.TaskTargetType) {
	IncreTaskProgress(hero, result, targetType, 1)
}

func IncreTaskProgress(hero *entity.Hero, result herolock.LockResult, targetType shared_proto.TaskTargetType, amount uint64) {

	var activeDegreeChanged bool
	hero.TaskList().WalkAllTask(func(task entity.Task) (endedWalk bool) {
		if task.Progress().IncreTaskTypeProgress(targetType, hero, amount) {
			result.Changed()
			result.Add(task.NewUpdateTaskProgressMsg())

			if task.Type() == entity.ACTIVE_DEGREE_TASK {
				activeDegreeChanged = true
			}

			if task.Progress().IsCompleted() {
				gamelogs.CompleteTaskLog(hero.Pid(), hero.Sid(), hero.Id(), uint64(task.Type()), task.Id())
			}
		}
		return false
	})

	if activeDegreeChanged {
		degree := hero.TaskList().ActiveDegreeTaskList().Degree()
		OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_ACTIVITY, degree)

		if hero.HistoryAmount().Amount(server_proto.HistoryAmountType_MaxActiveDegree) < degree {
			hero.HistoryAmount().Set(server_proto.HistoryAmountType_MaxActiveDegree, degree)
			UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACTIVE_DEGREE)
		}
	}
}

func IncreTaskProgressWithFunc(hero *entity.Hero, result herolock.LockResult, targetType shared_proto.TaskTargetType, increFunc func(hero *entity.Hero, t *entity.TaskProgress) uint64) {

	var activeDegreeChanged bool
	hero.TaskList().WalkAllTask(func(task entity.Task) (endedWalk bool) {
		if task.Progress().IncreTaskTypeProgressWithFunc(targetType, hero, increFunc) {
			result.Changed()
			result.Add(task.NewUpdateTaskProgressMsg())

			if task.Type() == entity.ACTIVE_DEGREE_TASK {
				activeDegreeChanged = true
			}

			if task.Progress().IsCompleted() {
				gamelogs.CompleteTaskLog(hero.Pid(), hero.Sid(), hero.Id(), uint64(task.Type()), task.Id())
			}
		}
		return false
	})

	if activeDegreeChanged {
		degree := hero.TaskList().ActiveDegreeTaskList().Degree()
		OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_ACTIVITY, degree)

		if hero.HistoryAmount().Amount(server_proto.HistoryAmountType_MaxActiveDegree) < degree {
			hero.HistoryAmount().Set(server_proto.HistoryAmountType_MaxActiveDegree, degree)
			UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACTIVE_DEGREE)
		}
	}
}
