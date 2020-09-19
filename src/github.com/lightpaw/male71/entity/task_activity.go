package entity

import (
	"github.com/lightpaw/male7/config/activitydata"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/task"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/pb/server_proto"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/taskdata"
)

func (t *TaskList) unmarshalActivityListMode(proto map[uint64]*server_proto.ActivityTaskListModeServerProto, datas *config.ConfigDatas, ctime time.Time) {
	if proto == nil {
		return
	}
	for id, a := range proto {
		activityData := datas.GetActivityTaskListModeData(id)
		if activityData == nil {
			continue
		}
		startTime, endTime := timeutil.Unix64(a.StartTime), timeutil.Unix64(a.EndTime)
		if !timeutil.Between(ctime, startTime, endTime) {
			continue
		}
		m := &ActivityTaskListMode{
			tasks: make(map[uint64]*ActivityTask),
		}
		m.startTime = startTime
		m.endTime = endTime
		for _, task := range a.Tasks {
			data := datas.GetActivityTaskData(task.Id)
			if data == nil {
				return
			}
			activityTask := newActivityTask(data)
			activityTask.isCollectPrize = task.IsCollected
			activityTask.progress.progress = task.Progress
			m.tasks[task.Id] = activityTask
		}
		t.activityTaskListModeMap[id] = m
	}
}

func (t *TaskList) WalkActivityTaskListMode(f func(taskId uint64, task *ActivityTask) (endWalk bool)) (isEndedWalk bool) {
	for _, a := range t.activityTaskListModeMap {
		if a.Walk(f) {
			return true
		}
	}
	return
}

func (t *TaskList) TryCorrectActivityTaskListMode(m map[uint64]*ActivityListModeTasks) (corrected bool) {
	for id, _ := range t.activityTaskListModeMap {
		if m[id] == nil {
			delete(t.activityTaskListModeMap, id)
			corrected = true
		}
	}
	for id, a := range m {
		info := t.activityTaskListModeMap[id]
		if info == nil || !info.startTime.Equal(a.startTime) || !info.endTime.Equal(a.endTime) {
			t.activityTaskListModeMap[id] = newActivityTaskListMode(a.data, a.startTime, a.endTime)
			corrected = true
		}
	}
	return
}

func (t *TaskList) EncodeActivityTaskListMode() []*shared_proto.HeroActivityTaskListModeProto {
	p := []*shared_proto.HeroActivityTaskListModeProto{}
	for id, a := range t.activityTaskListModeMap {
		proto := &shared_proto.HeroActivityTaskListModeProto{
			Id:       u64.Int32(id),
			Progress: a.EncodeClient(),
		}
		for _, t := range a.tasks {
			if t.isCollectPrize {
				proto.Collected = append(proto.Collected, u64.Int32(t.data.Id))
			}
		}
		p = append(p, proto)
	}
	return p
}

func newActivityTaskListMode(data *activitydata.ActivityTaskListModeData, startTime, endTime time.Time) *ActivityTaskListMode {
	m := &ActivityTaskListMode{
		tasks:     make(map[uint64]*ActivityTask),
		startTime: startTime,
		endTime:   endTime,
	}
	for _, task := range data.Tasks {
		m.tasks[task.Id] = newActivityTask(task)
	}
	return m
}

// 列表式任务活动
type ActivityTaskListMode struct {
	tasks     map[uint64]*ActivityTask
	startTime time.Time
	endTime   time.Time
}

func (m *ActivityTaskListMode) EncodeClient() []*shared_proto.Int32Pair {
	p := []*shared_proto.Int32Pair{}
	for _, task := range m.tasks {
		p = append(p, task.encodeClient())
	}
	return p
}

func (m *ActivityTaskListMode) encodeServer() *server_proto.ActivityTaskListModeServerProto {
	p := &server_proto.ActivityTaskListModeServerProto{}
	p.StartTime = timeutil.Marshal64(m.startTime)
	p.EndTime = timeutil.Marshal64(m.endTime)

	for _, task := range m.tasks {
		p.Tasks = append(p.Tasks, task.encodeServer())
	}
	return p
}

func (m *ActivityTaskListMode) Walk(f func(taskId uint64, task *ActivityTask) (endWalk bool)) (isEndedWalk bool) {
	if len(m.tasks) <= 0 {
		return
	}
	for id, task := range m.tasks {
		if f(id, task) {
			return true
		}
	}
	return
}

// 活动任务样板
func newActivityTask(data *taskdata.ActivityTaskData) *ActivityTask {
	t := &ActivityTask{}
	t.data = data
	t.progress = &TaskProgress{}
	t.progress.target = data.Target
	return t
}

type ActivityTask struct {
	data           *taskdata.ActivityTaskData
	progress       *TaskProgress
	isCollectPrize bool
}

func (t *ActivityTask) Id() uint64 {
	return t.data.Id
}

func (t *ActivityTask) Data() *taskdata.ActivityTaskData {
	return t.data
}

func (t *ActivityTask) IsCollectPrize() bool {
	return t.isCollectPrize
}

func (t *ActivityTask) CollectPrize() {
	t.isCollectPrize = true
}

func (t *ActivityTask) Type() task_type {
	return ACTIVITY_TASK
}

func (t *ActivityTask) NewTaskMsg() pbutil.Buffer {
	return task.NewS2cNewTaskMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.encodeClient()), t.progress.IsCompleted(), false)
}

func (t *ActivityTask) Progress() *TaskProgress {
	return t.progress
}

func (t *ActivityTask) NewUpdateTaskProgressMsg() pbutil.Buffer {
	return task.NewS2cUpdateTaskProgressMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.GetProgress()), t.progress.IsCompleted())
}

func (t *ActivityTask) encodeClient() *shared_proto.Int32Pair {
	return &shared_proto.Int32Pair{
		Key:   u64.Int32(t.data.Id),
		Value: u64.Int32(t.progress.encodeClient()),
	}
}

func (t *ActivityTask) encodeServer() *server_proto.ActivityTaskServerProto {
	proto := &server_proto.ActivityTaskServerProto{}
	proto.Id = t.data.Id
	proto.Progress = t.progress.encodeClient()
	proto.IsCollected = t.isCollectPrize
	return proto
}
