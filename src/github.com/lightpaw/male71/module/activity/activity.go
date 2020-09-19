package activity

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/activity"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/service/operate_type"
	"sync/atomic"
)

func NewActivityModule(dep iface.ServiceDep, tickerService iface.TickerService) *ActivityModule {
	m := &ActivityModule{
		dep:            dep,
		tickerService:  tickerService,
		closeNotify:    make(chan struct{}),
		loopExitNotify: make(chan struct{}),
	}
	ctime := dep.Time().CurrentTime()

	serverStartTime := dep.SvrConf().GetServerStartTime()

	// 活动展示
	show := entity.NewActivityShowList(dep.Datas().GetActivityShowDataArray(), serverStartTime, ctime)
	show.RefreshMsg()
	m.activityShowRef.Store(show)
	// 收集活动
	collection := entity.NewActivityCollectionList(dep.Datas().GetActivityCollectionDataArray(), serverStartTime, ctime)
	collection.RefreshMsg()
	m.collectionRef.Store(collection)
	// 列表式任务活动
	taskListMode := entity.NewActivityListModeTasksList(dep.Datas().GetActivityTaskListModeDataArray(), serverStartTime, ctime)
	taskListMode.RefreshMsg()
	m.taskListModeRef.Store(taskListMode)

	heromodule.RegisterHeroOnlineListener(m)

	go call.CatchPanic(m.loop, "活动loop")
	return m
}

//gogen:iface
type ActivityModule struct {
	dep           iface.ServiceDep
	tickerService iface.TickerService

	// 活动展示
	activityShowRef atomic.Value

	// 收集活动
	collectionRef atomic.Value

	// 列表式任务活动
	taskListModeRef atomic.Value

	// 多切页列表式活动

	closeNotify    chan struct{}
	loopExitNotify chan struct{}
}

func (m *ActivityModule) getActivityShow() *entity.ActivityShowList {
	return m.activityShowRef.Load().(*entity.ActivityShowList)
}

func (m *ActivityModule) setActivityShow(a *entity.ActivityShowList) {
	m.activityShowRef.Store(a)
}

func (m *ActivityModule) getActivityCollection() *entity.ActivityCollectionList {
	return m.collectionRef.Load().(*entity.ActivityCollectionList)
}

func (m *ActivityModule) setActivityCollection(a *entity.ActivityCollectionList) {
	m.collectionRef.Store(a)
}

func (m *ActivityModule) getActivityTaskListMode() *entity.ActivityListModeTasksList {
	return m.taskListModeRef.Load().(*entity.ActivityListModeTasksList)
}

func (m *ActivityModule) setActivityTaskListMode(a *entity.ActivityListModeTasksList) {
	m.taskListModeRef.Store(a)
}

func (m *ActivityModule) OnHeroOnline(hc iface.HeroController) {
	// 活动展示
	show := m.getActivityShow()
	if show.Msg() != nil {
		hc.Send(show.Msg())
	}
	isLockHero := false
	// 活动
	collection := m.getActivityCollection()
	isCollectionOpen := collection.Msg() != nil
	if isCollectionOpen {
		hc.Send(collection.Msg())
		isLockHero = true
	}
	taskListMode := m.getActivityTaskListMode()
	isTaskListModeOpen := taskListMode.Msg() != nil
	if isTaskListModeOpen {
		hc.Send(taskListMode.Msg())
		isLockHero = true
	}
	if !isLockHero {
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if isCollectionOpen {
			hero.Activity().TryCorrectCollection(collection.GetMap())
			result.Add(activity.NewS2cNoticeCollectionCountsMsg(hero.Activity().EncodeClient()))
		}
		if isTaskListModeOpen {
			taskList := hero.TaskList()
			taskList.TryCorrectActivityTaskListMode(taskListMode.GetMap())
			taskList.WalkActivityTaskListMode(func(taskId uint64, task *entity.ActivityTask) (endWalk bool) {
				task.Progress().UpdateTaskProgress(hero)
				return
			})
			result.Add(activity.NewS2cNoticeTaskListModeProgressMsg(taskList.EncodeActivityTaskListMode()))
		}
		result.Ok()
	})
}

func (m *ActivityModule) Close() {
	close(m.closeNotify)
	<-m.loopExitNotify
}

func (m *ActivityModule) loop() {
	defer close(m.loopExitNotify)

	tickTime := m.tickerService.GetPerMinuteTickTime()

	for {
		select {
		case <-tickTime.Tick():
			tickTime = m.tickerService.GetPerMinuteTickTime()
			call.CatchPanic(m.update, "活动update")
		case <-m.closeNotify:
			return
		}
	}
}

func (m *ActivityModule) update() {
	serverStartTime := m.dep.SvrConf().GetServerStartTime()
	ctime := m.dep.Time().CurrentTime()
	// 活动展示
	show := m.getActivityShow()
	newShow := entity.NewActivityShowList(m.dep.Datas().GetActivityShowDataArray(), serverStartTime, ctime)
	if !show.Equal(newShow) {
		if newShow.RefreshMsg(); newShow.Msg() != nil {
			m.dep.World().Broadcast(newShow.Msg())
		}
		m.setActivityShow(newShow)
	}
	// 收集活动
	collection := m.getActivityCollection()
	newCollection := entity.NewActivityCollectionList(m.dep.Datas().GetActivityCollectionDataArray(), serverStartTime, ctime)
	if !collection.Equal(newCollection) {
		if newCollection.RefreshMsg(); newCollection.Msg() != nil {
			m.dep.World().Broadcast(newCollection.Msg())
		}
		m.setActivityCollection(newCollection)
	}
	// 列表式任务活动
	taskListMode := m.getActivityTaskListMode()
	newTaskListMode := entity.NewActivityListModeTasksList(m.dep.Datas().GetActivityTaskListModeDataArray(), serverStartTime, ctime)
	if !taskListMode.Equal(newTaskListMode) {
		newTaskListMode.RefreshMsg()
		m.setActivityTaskListMode(newTaskListMode)
		if newTaskListMode.Msg() != nil {
			m.dep.World().Broadcast(newTaskListMode.Msg())
			m.dep.World().WalkHero(func(id int64, hc iface.HeroController) {
				hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
					if hero.TaskList().TryCorrectActivityTaskListMode(newTaskListMode.GetMap()) {
						hero.TaskList().WalkActivityTaskListMode(func(taskId uint64, task *entity.ActivityTask) (endWalk bool) {
							task.Progress().UpdateTaskProgress(hero)
							return
						})
						result.Add(activity.NewS2cNoticeTaskListModeProgressMsg(hero.TaskList().EncodeActivityTaskListMode()))
					}
					result.Ok()
				})
			})
		}
	}
}

//gogen:iface
func (m *ActivityModule) ProcessCollectCollection(proto *activity.C2SCollectCollectionProto, hc iface.HeroController) {
	collection := m.getActivityCollection()
	a := collection.GetActivity(u64.FromInt32(proto.GetId()))
	if a == nil {
		logrus.Debugf("收集兑换，不在活动期间内")
		hc.Send(activity.ERR_COLLECT_COLLECTION_FAIL_OUT_TIME)
		return
	}
	data := a.Data().GetExchange(u64.FromInt32(proto.GetExchangeId()))
	if data == nil {
		logrus.Debugf("收集兑换，活动中没有该兑换id")
		hc.Send(activity.ERR_COLLECT_COLLECTION_FAIL_INVALID_ID)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Activity().TryCorrectCollection(collection.GetMap()) {
			result.Add(activity.NewS2cNoticeCollectionCountsMsg(hero.Activity().EncodeClient()))
		}
		if data.Limit > 0 && hero.Activity().GetCollectionTimes(a.Id(), data.Id) >= data.Limit { // 有兑换次数限制
			logrus.Debugf("收集兑换，兑换上限")
			result.Add(activity.ERR_COLLECT_COLLECTION_FAIL_LIMIT)
			return
		}
		hctx := heromodule.NewContext(m.dep, operate_type.ActivityCollect)
		if !heromodule.TryReduceCost(hctx, hero, result, data.Combine.Cost) {
			logrus.Debugf("收集兑换，收集品不足")
			result.Add(activity.ERR_COLLECT_COLLECTION_FAIL_NOT_ENOUGH_COST)
			return
		}
		heromodule.AddPrize(hctx, hero, result, data.Combine.Prize, m.dep.Time().CurrentTime())
		if data.Limit > 0 {
			hero.Activity().IncCollectionTimes(a.Id(), data.Id)
		}
		result.Add(activity.NewS2cCollectCollectionMsg(proto.Id, proto.ExchangeId))
		result.Ok()
	})
}
