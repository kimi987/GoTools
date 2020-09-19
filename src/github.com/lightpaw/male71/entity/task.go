package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/male7/gen/pb/task"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
)

func newTaskList(initData *heroinit.HeroInitData) *TaskList {
	t := &TaskList{}
	t.mainTask = newMainTask()
	t.mainTask.setTaskData(initData.FirstMainTask)
	t.mainTaskSequence = t.mainTask.data.Sequence

	t.branchTaskMap = make(map[uint64]*BranchTask)

	t.achieveTaskList = newAchieveTaskList(initData.AchieveTaskDatas)
	t.activeDegreeTaskList = newActiveDegreeTaskList(initData.ActiveDegreeTaskDatas)
	t.bwzlTaskList = newBwzlTaskList(initData.BwzlTaskDatas)

	t.baYeStage = NewBaYeStage(initData.FirstBaYeStageData)

	t.titleTaskList = newTitleTaskList(initData.FirstTitleData)
	t.activityTaskListModeMap = make(map[uint64]*ActivityTaskListMode)

	t.taskMonsterMap = make(map[uint64]bool)

	return t
}

type TaskList struct {
	mainTask         *MainTask
	mainTaskSequence uint64

	branchTaskMap map[uint64]*BranchTask

	achieveTaskList      *AchieveTaskList
	activeDegreeTaskList *ActiveDegreeTaskList
	bwzlTaskList         *BwzlTaskList

	collectTaskBoxId uint64

	baYeStage     *BaYeStage // 霸业目标
	lastBaYeStage uint64     // 上次完成的霸业目标stage

	titleTaskList *TitleTaskList

	activityTaskListModeMap map[uint64]*ActivityTaskListMode // 各种列表示活动任务

	taskMonsterMap map[uint64]bool

	nextCheckRealmInfoTime time.Time // 下次检查时间
}

func (t *TaskList) GetNextCheckRealmInfoTime() time.Time {
	return t.nextCheckRealmInfoTime
}

func (t *TaskList) SetNextCheckRealmInfoTime(toSet time.Time) {
	t.nextCheckRealmInfoTime = toSet
}

func (t *TaskList) GetTaskMonster(monsterId uint64) (create, exist bool) {
	create, exist = t.taskMonsterMap[monsterId]
	return
}

func (t *TaskList) SetTaskMonster(monsterId uint64, create bool) {
	t.taskMonsterMap[monsterId] = create
}

func (t *TaskList) RemoveTaskMonster(monsterId uint64) {
	delete(t.taskMonsterMap, monsterId)
}

func (t *TaskList) HasTaskMonster() bool {
	return len(t.taskMonsterMap) > 0
}

func (t *TaskList) RangeTaskMonsterId(f func(monsterId uint64, created bool)) {

	if t.HasTaskMonster() {
		for monsterId, created := range t.taskMonsterMap {
			f(monsterId, created)
		}
	}

}

func (t *TaskList) WalkAllTask(walkFunc func(task Task) (endedWalk bool)) {
	if mainTask := t.MainTask(); mainTask != nil && walkFunc(mainTask) {
		return
	}

	if t.WalkBranchTask(func(taskId uint64, task *BranchTask) bool { return walkFunc(task) }) {
		return
	}

	if baYeStage := t.BaYeStage(); baYeStage != nil && baYeStage.Walk(func(task *BaYeTask) bool { return walkFunc(task) }) {
		return
	}

	if t.achieveTaskList.Walk(func(taskId uint64, task *AchieveTask) bool { return walkFunc(task) }) {
		return
	}

	if t.activeDegreeTaskList.Walk(func(taskId uint64, task *ActiveDegreeTask) bool { return walkFunc(task) }) {
		return
	}

	if t.bwzlTaskList.Walk(func(taskId uint64, task *BwzlTask) bool { return walkFunc(task) }) {
		return
	}

	if t.titleTaskList.Walk(func(taskId uint64, task *TitleTask) bool { return walkFunc(task) }) {
		return
	}

	for _, a := range t.activityTaskListModeMap {
		if a.Walk(func(taskId uint64, task *ActivityTask) bool { return walkFunc(task) }) {
			return
		}
	}
}

func (t *TaskList) Task(taskId uint64) Task {
	if task := t.MainTask(); task != nil && task.Data().Sequence == taskId {
		return task
	}

	if task := t.GetBranchTask(taskId); task != nil {
		return task
	}

	if task := t.achieveTaskList.Task(taskId); task != nil {
		return task
	}

	if task := t.activeDegreeTaskList.Task(taskId); task != nil {
		return task
	}

	if task := t.bwzlTaskList.Task(taskId); task != nil {
		return task
	}

	if baYeStage := t.BaYeStage(); baYeStage != nil {
		if task := baYeStage.GetTask(taskId); task != nil {
			return task
		}
	}

	if task := t.titleTaskList.Task(taskId); task != nil {
		return task
	}

	for _, a := range t.activityTaskListModeMap {
		if task := a.tasks[taskId]; task != nil {
			return task
		}
	}

	return nil
}

func (t *TaskList) CollectedBoxId() uint64 {
	return t.collectTaskBoxId
}

func (t *TaskList) IncreseCollectedBoxId() {
	t.collectTaskBoxId++
}

func (t *TaskList) MainTask() *MainTask {
	return t.mainTask
}

func (t *TaskList) GmSetMainTaskData(data *taskdata.MainTaskData) {
	if t.mainTask == nil {
		t.mainTask = newMainTask()
	}

	t.mainTask.setTaskData(data)
}

func (t *TaskList) CompletedMainTaskSequence() uint64 {
	if t.mainTask == nil {
		return t.mainTaskSequence
	}
	return t.mainTask.data.Sequence
}

func (t *TaskList) GetBranchTask(id uint64) *BranchTask {
	return t.branchTaskMap[id]
}

func (t *TaskList) WalkBranchTask(f func(taskId uint64, data *BranchTask) bool) (isEndedWalk bool) {
	if len(t.branchTaskMap) <= 0 {
		return
	}

	for k, v := range t.branchTaskMap {
		if f(k, v) {
			return true
		}
	}

	return
}

func (t *TaskList) AchieveTaskList() *AchieveTaskList {
	return t.achieveTaskList
}

func (t *TaskList) ActiveDegreeTaskList() *ActiveDegreeTaskList {
	return t.activeDegreeTaskList
}

func (t *TaskList) BwzlTaskList() *BwzlTaskList {
	return t.bwzlTaskList
}

func (t *TaskList) BaYeStage() *BaYeStage {
	return t.baYeStage
}

// 完成了的霸业阶段
func (t *TaskList) GetCompletedBaYeStage() uint64 {
	if t.baYeStage != nil && t.baYeStage.Data().Prev != nil {
		return t.baYeStage.Data().Prev.Stage
	}
	return t.lastBaYeStage
}

func (t *TaskList) CompleteBaYeStage() {
	t.lastBaYeStage = t.baYeStage.Data().Stage
	next := t.baYeStage.Data().Next
	if next == nil {
		t.baYeStage = nil
	} else {
		t.baYeStage = NewBaYeStage(next)
	}
}

func (t *TaskList) CompleteMainTask() {
	if t.mainTask != nil {
		// 加支线任务
		for _, b := range t.mainTask.data.BranchTask {
			t.branchTaskMap[b.Id] = newBranchTask(b)
		}

		// 下一环主线
		nextTask := t.mainTask.data.NextTask()
		if nextTask != nil {
			t.mainTask.setTaskData(nextTask)
		} else {
			t.mainTask = nil
		}
		t.mainTaskSequence++
	}
}

func (t *TaskList) CompleteBranchTask(branchTask *BranchTask) {

	if branchTask != nil {
		delete(t.branchTaskMap, branchTask.data.Id)

		nextTask := branchTask.data.NextTask()
		if nextTask != nil {
			t.branchTaskMap[nextTask.Id] = newBranchTask(nextTask)
		}
	}
}

func (t *TaskList) ResetDaily(hero *Hero) {
	t.activeDegreeTaskList.ResetDaily(hero) // 重置活跃度任务
}

func (t *TaskList) unmarshal(proto *server_proto.HeroTaskServerProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	t.mainTaskSequence = proto.MainTask
	mainTaskData := datas.MainTaskData().Must(t.mainTaskSequence)
	if mainTaskData.Sequence >= t.mainTaskSequence {
		t.mainTask.setTaskData(mainTaskData)
		t.mainTask.progress.progress = proto.MainTaskProgress
	} else {
		t.mainTask = nil
	}

	n := imath.Min(len(proto.BranchTask), len(proto.BranchTaskProgress))
	for i := 0; i < n; i++ {
		data := datas.GetBranchTaskData(proto.BranchTask[i])
		if data != nil {
			bt := newBranchTask(data)
			bt.progress.progress = proto.BranchTaskProgress[i]
			t.branchTaskMap[data.Id] = bt
		}
	}

	t.achieveTaskList.unmarshal(proto.AchieveTaskList, datas)
	t.activeDegreeTaskList.unmarshal(proto.ActiveDegreeTaskList, datas)
	t.bwzlTaskList.unmarshal(proto.BwzlTaskList, datas)

	t.collectTaskBoxId = proto.CollectTaskBoxId

	// 先清空，这里已经不是创建新号了
	t.baYeStage = nil
	t.lastBaYeStage = proto.GetLastCompleteBaYeStage()
	baYeStageProto := proto.GetBaYeStage()
	if baYeStageProto != nil {
		data := datas.GetBaYeStageData(baYeStageProto.GetStage())
		if data != nil {
			t.baYeStage = NewBaYeStage(data)
			t.baYeStage.unmarshal(baYeStageProto, datas)
		}
	}

	if t.baYeStage == nil {
		// 没有霸业目标，看下是不是因为没找到或者没有过
		stageData := datas.GetBaYeStageData(t.lastBaYeStage + 1)
		if stageData != nil {
			t.baYeStage = NewBaYeStage(stageData)
		}
	}

	t.titleTaskList.unmarshal(proto.TitleTaskList, datas)

	n = imath.Min(len(proto.TaskMonster), len(proto.TaskMonsterCreate))
	for i := 0; i < n; i++ {
		t.taskMonsterMap[proto.TaskMonster[i]] = proto.TaskMonsterCreate[i]
	}
}

func (t *TaskList) encodeServer() *server_proto.HeroTaskServerProto {
	proto := &server_proto.HeroTaskServerProto{}
	proto.MainTask = t.mainTaskSequence
	if t.mainTask != nil {
		proto.MainTaskProgress = t.mainTask.progress.progress
	}

	for _, b := range t.branchTaskMap {
		proto.BranchTask = append(proto.BranchTask, b.data.Id)
		proto.BranchTaskProgress = append(proto.BranchTaskProgress, b.progress.progress)
	}

	proto.AchieveTaskList = t.achieveTaskList.encodeServer()
	proto.ActiveDegreeTaskList = t.activeDegreeTaskList.encodeServer()
	proto.BwzlTaskList = t.bwzlTaskList.encodeServer()

	proto.CollectTaskBoxId = t.collectTaskBoxId

	if t.baYeStage != nil {
		proto.BaYeStage = t.baYeStage.encodeServer()
	}

	proto.LastCompleteBaYeStage = t.lastBaYeStage

	proto.TitleTaskList = t.titleTaskList.encodeServer()
	if len(t.activityTaskListModeMap) > 0 {
		proto.ActivityTaskListModeMap = make(map[uint64]*server_proto.ActivityTaskListModeServerProto)
		for id, a := range t.activityTaskListModeMap {
			proto.ActivityTaskListModeMap[id] = a.encodeServer()
		}
	}

	for k, v := range t.taskMonsterMap {
		proto.TaskMonster = append(proto.TaskMonster, k)
		proto.TaskMonsterCreate = append(proto.TaskMonsterCreate, v)
	}

	return proto
}

func (t *TaskList) EncodeClient() *shared_proto.HeroTaskProto {
	proto := &shared_proto.HeroTaskProto{}
	if t.mainTask != nil {
		proto.MainTaskId = u64.Int32(t.mainTask.data.Sequence)
		proto.MainTaskProgress = u64.Int32(t.mainTask.progress.encodeClient())
	}

	for _, b := range t.branchTaskMap {
		proto.BranchTaskId = append(proto.BranchTaskId, u64.Int32(b.data.Id))
		proto.BranchTaskProgress = append(proto.BranchTaskProgress, u64.Int32(b.progress.encodeClient()))
	}

	proto.AchieveTaskList = t.achieveTaskList.encodeClient()
	proto.ActiveDegreeTaskList = t.activeDegreeTaskList.encodeClient()
	proto.BwzlTaskList = t.bwzlTaskList.encodeClient()

	proto.CollectTaskBoxId = u64.Int32(t.collectTaskBoxId)

	if t.baYeStage != nil {
		proto.BaYeStage = t.baYeStage.EncodeClient()
	}

	proto.TitleTaskList = t.titleTaskList.encodeClient()

	return proto
}

type task_type uint8

const (
	MAIN_TASK          task_type = iota
	BRANCH_TASK
	BAYE_TASK
	ACHIEVE_TASK
	ACTIVE_DEGREE_TASK
	BWZL_TASK
	TITLE_TASK
	ACTIVITY_TASK
)

func (t *MainTask) Type() task_type {
	return MAIN_TASK
}

func (t *BranchTask) Type() task_type {
	return BRANCH_TASK
}

func (t *BaYeTask) Type() task_type {
	return BAYE_TASK
}

func (t *AchieveTask) Type() task_type {
	return ACHIEVE_TASK
}

func (t *ActiveDegreeTask) Type() task_type {
	return ACTIVE_DEGREE_TASK
}

type Task interface {
	NewTaskMsg() pbutil.Buffer
	NewUpdateTaskProgressMsg() pbutil.Buffer
	Progress() *TaskProgress
	Type() task_type
	Id() uint64
}

// main
func newMainTask() *MainTask {
	t := &MainTask{}
	t.progress = &TaskProgress{}
	return t
}

type MainTask struct {
	data *taskdata.MainTaskData

	progress *TaskProgress
}

func (t *MainTask) NewTaskMsg() pbutil.Buffer {
	return task.NewS2cNewTaskMsg(u64.Int32(t.data.Sequence), u64.Int32(t.progress.encodeClient()), t.progress.IsCompleted(), true)
}

func (t *MainTask) Data() *taskdata.MainTaskData {
	return t.data
}

func (t *MainTask) Progress() *TaskProgress {
	return t.progress
}

func (t *MainTask) Id() uint64 {
	return t.data.Sequence
}

func (t *MainTask) NewUpdateTaskProgressMsg() pbutil.Buffer {
	return task.NewS2cUpdateTaskProgressMsg(u64.Int32(t.data.Sequence), u64.Int32(t.progress.GetProgress()), t.progress.IsCompleted())
}

func (t *MainTask) setTaskData(data *taskdata.MainTaskData) {
	t.data = data
	t.progress.target = data.Target
	t.progress.progress = 0
}

// branch
func newBranchTask(data *taskdata.BranchTaskData) *BranchTask {
	t := &BranchTask{}
	t.data = data
	t.progress = &TaskProgress{}
	t.progress.target = data.Target

	return t
}

type BranchTask struct {
	data *taskdata.BranchTaskData

	progress *TaskProgress
}

func (t *BranchTask) NewTaskMsg() pbutil.Buffer {
	return task.NewS2cNewTaskMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.encodeClient()), t.progress.IsCompleted(), false)
}

func (t *BranchTask) Data() *taskdata.BranchTaskData {
	return t.data
}

func (t *BranchTask) Progress() *TaskProgress {
	return t.progress
}

func (t *BranchTask) Id() uint64 {
	return t.data.Id
}

func (t *BranchTask) NewUpdateTaskProgressMsg() pbutil.Buffer {
	return task.NewS2cUpdateTaskProgressMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.GetProgress()), t.progress.IsCompleted())
}

// achieve
func newAchieveTask(data *taskdata.AchieveTaskData) *AchieveTask {
	t := &AchieveTask{}
	t.data = data
	t.progress = &TaskProgress{}
	t.progress.target = data.Target

	return t
}

type AchieveTask struct {
	data *taskdata.AchieveTaskData

	progress *TaskProgress

	isCollectPrize bool

	reachTime time.Time
}

func (t *AchieveTask) NewTaskMsg() pbutil.Buffer {
	return task.NewS2cNewTaskMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.encodeClient()), t.progress.IsCompleted(), false)
}

func (t *AchieveTask) Data() *taskdata.AchieveTaskData {
	return t.data
}

func (t *AchieveTask) Progress() *TaskProgress {
	return t.progress
}

func (t *AchieveTask) Id() uint64 {
	return t.data.Id
}

func (t *AchieveTask) SetReachTime(toSet time.Time) {
	t.reachTime = toSet
}

func (t *AchieveTask) ReachTime() time.Time {
	return t.reachTime
}

func (t *AchieveTask) NewUpdateTaskProgressMsg() pbutil.Buffer {
	return task.NewS2cUpdateTaskProgressMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.GetProgress()), t.progress.IsCompleted())
}

func (t *AchieveTask) IsCollectPrize() bool {
	return t.isCollectPrize
}

func (t *AchieveTask) CollectPrize() {
	t.isCollectPrize = true
}

func (t *AchieveTask) encodeClient() *shared_proto.AchieveTaskProto {
	proto := &shared_proto.AchieveTaskProto{}

	proto.Id = u64.Int32(t.data.Id)
	proto.Progress = u64.Int32(t.progress.encodeClient())
	proto.IsCollected = t.IsCollectPrize()
	proto.ReachTime = timeutil.Marshal32(t.reachTime)

	return proto
}

func (t *AchieveTask) encodeServer() *server_proto.AchieveTaskServerProto {
	proto := &server_proto.AchieveTaskServerProto{}

	proto.Id = t.data.Id
	proto.Progress = t.progress.encodeClient()
	proto.IsCollected = t.IsCollectPrize()
	proto.ReachTime = timeutil.Marshal64(t.reachTime)

	return proto
}

// 活跃度
func newActiveDegreeTask(data *taskdata.ActiveDegreeTaskData) *ActiveDegreeTask {
	t := &ActiveDegreeTask{}
	t.data = data
	t.progress = &TaskProgress{}
	t.progress.target = data.Target

	return t
}

type ActiveDegreeTask struct {
	data *taskdata.ActiveDegreeTaskData

	progress *TaskProgress
}

func (t *ActiveDegreeTask) NewTaskMsg() pbutil.Buffer {
	return task.NewS2cNewTaskMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.encodeClient()), t.progress.IsCompleted(), false)
}

func (t *ActiveDegreeTask) Data() *taskdata.ActiveDegreeTaskData {
	return t.data
}

func (t *ActiveDegreeTask) Progress() *TaskProgress {
	return t.progress
}

func (t *ActiveDegreeTask) Id() uint64 {
	return t.data.Id
}

func (t *ActiveDegreeTask) NewUpdateTaskProgressMsg() pbutil.Buffer {
	return task.NewS2cUpdateTaskProgressMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.GetProgress()), t.progress.IsCompleted())
}

func (t *ActiveDegreeTask) encodeClient() *shared_proto.ActiveDegreeTaskProto {
	proto := &shared_proto.ActiveDegreeTaskProto{}

	proto.Id = u64.Int32(t.data.Id)
	proto.Progress = u64.Int32(t.progress.encodeClient())

	return proto
}

func (t *ActiveDegreeTask) encodeServer() *server_proto.ActiveDegreeTaskServerProto {
	proto := &server_proto.ActiveDegreeTaskServerProto{}

	proto.Id = t.data.Id
	proto.Progress = t.progress.encodeClient()

	return proto
}

func NewBaYeStage(data *taskdata.BaYeStageData) *BaYeStage {
	s := &BaYeStage{data: data}

	s.tasks = make([]*BaYeTask, 0, len(data.Tasks))

	for _, d := range data.Tasks {
		s.tasks = append(s.tasks, newBaYeTask(d))
	}

	return s
}

type BaYeStage struct {
	data  *taskdata.BaYeStageData // 阶段数据
	tasks []*BaYeTask             // 任务
}

func (s *BaYeStage) Data() *taskdata.BaYeStageData {
	return s.data
}

func (s *BaYeStage) Walk(walkFunc func(t *BaYeTask) (endWalk bool)) (isEndedWalk bool) {
	if len(s.tasks) <= 0 {
		return
	}

	for _, t := range s.tasks {
		if walkFunc(t) {
			return true
		}
	}

	return
}

func (s *BaYeStage) GetTask(id uint64) *BaYeTask {
	for _, t := range s.tasks {
		if t.Data().Id == id {
			return t
		}
	}

	return nil
}

// 能否领取阶段奖励
func (s *BaYeStage) CanCollectStagePrize() (hasTaskNotComplete, hasTaskNotCollectPrize bool) {
	for _, t := range s.tasks {
		if !t.Progress().IsCompleted() {
			hasTaskNotComplete = true
			return
		}
	}

	for _, t := range s.tasks {
		if !t.IsCollectPrize() {
			hasTaskNotCollectPrize = true
			return
		}
	}

	return
}

func (s *BaYeStage) EncodeClient() *shared_proto.HeroBaYeStageProto {
	proto := &shared_proto.HeroBaYeStageProto{}

	proto.Stage = u64.Int32(s.data.Stage)
	proto.TaskId = make([]int32, len(s.tasks))
	proto.TaskProgress = make([]int32, len(s.tasks))
	proto.IsCollected = make([]bool, len(s.tasks))
	for idx, t := range s.tasks {
		proto.TaskId[idx] = u64.Int32(t.Data().Id)
		proto.TaskProgress[idx] = u64.Int32(t.Progress().GetProgress())
		proto.IsCollected[idx] = t.IsCollectPrize()
	}

	return proto
}

func (s *BaYeStage) encodeServer() *server_proto.HeroBaYeStageServerProto {
	proto := &server_proto.HeroBaYeStageServerProto{}

	proto.Stage = s.data.Stage
	proto.TaskId = make([]uint64, len(s.tasks))
	proto.TaskProgress = make([]uint64, len(s.tasks))
	proto.IsCollected = make([]bool, len(s.tasks))
	for idx, t := range s.tasks {
		proto.TaskId[idx] = t.Data().Id
		proto.TaskProgress[idx] = t.Progress().GetProgress()
		proto.IsCollected[idx] = t.IsCollectPrize()
	}

	return proto
}

func (s *BaYeStage) unmarshal(proto *server_proto.HeroBaYeStageServerProto, datas *config.ConfigDatas) {
	for idx, task := range s.tasks {
		if idx >= len(proto.GetTaskId()) {
			break
		}

		// 多给次奖励就多给次
		if task.Data().Id != proto.GetTaskId()[idx] {
			continue
		}

		task.progress.setProgress(proto.GetTaskProgress()[idx])
		task.isCollectPrize = proto.GetIsCollected()[idx]
	}
}

// 霸业目标
func newBaYeTask(data *taskdata.BaYeTaskData) *BaYeTask {
	t := &BaYeTask{}
	t.data = data
	t.progress = &TaskProgress{}
	t.progress.target = data.Target

	return t
}

type BaYeTask struct {
	data *taskdata.BaYeTaskData

	progress *TaskProgress

	isCollectPrize bool
}

func (t *BaYeTask) NewTaskMsg() pbutil.Buffer {
	return task.NewS2cNewTaskMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.encodeClient()), t.progress.IsCompleted(), false)
}

func (t *BaYeTask) Data() *taskdata.BaYeTaskData {
	return t.data
}

func (t *BaYeTask) Progress() *TaskProgress {
	return t.progress
}

func (t *BaYeTask) Id() uint64 {
	return t.data.Id
}

func (t *BaYeTask) IsCollectPrize() bool {
	return t.isCollectPrize
}

func (t *BaYeTask) CollectPrize() {
	t.isCollectPrize = true
}

func (t *BaYeTask) NewUpdateTaskProgressMsg() pbutil.Buffer {
	return task.NewS2cUpdateTaskProgressMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.GetProgress()), t.progress.IsCompleted())
}

func newAchieveTaskList(achieveTaskDatas map[uint64]*taskdata.AchieveTaskData) *AchieveTaskList {
	taskList := &AchieveTaskList{tasks: make(map[uint64]*AchieveTask)}

	for _, data := range achieveTaskDatas {
		taskList.Give(data)
	}

	return taskList
}

// 成就任务
type AchieveTaskList struct {
	tasks                  map[uint64]*AchieveTask // 成就任务列表
	collectStarCount       []uint64                // 领取了的星级奖励数量
	selectShowAchieveTypes []uint64                // 被设置来展示的成就类型
}

func (t *AchieveTaskList) Count() int {
	return len(t.tasks)
}

func (t *AchieveTaskList) Task(id uint64) *AchieveTask {
	return t.tasks[id]
}

func (t *AchieveTaskList) TotalStar() (result uint64) {
	t.Walk(func(taskId uint64, task *AchieveTask) (endWalk bool) {
		if task.IsCollectPrize() {
			result += task.Data().TotalStar
		} else if prevTask := task.Data().PrevTask; prevTask != nil {
			result += prevTask.TotalStar
		}
		return
	})

	return
}

func (t *AchieveTaskList) Walk(f func(taskId uint64, task *AchieveTask) (endWalk bool)) (isEndedWalk bool) {
	if len(t.tasks) <= 0 {
		return
	}

	for k, v := range t.tasks {
		if f(k, v) {
			return true
		}
	}

	return
}

func (t *AchieveTaskList) Complete(achieveTask *AchieveTask, ctime time.Time) (nextAchieveTask *AchieveTask) {
	if achieveTask == nil {
		return
	}

	nextTask := achieveTask.data.NextTask()
	if nextTask == nil {
		// 成就任务如果没有下一个任务，则不删除
		achieveTask.SetReachTime(ctime)
		return
	}

	delete(t.tasks, achieveTask.data.Id)
	nextAchieveTask = newAchieveTask(nextTask)
	t.tasks[nextTask.Id] = nextAchieveTask
	nextAchieveTask.SetReachTime(ctime)

	return
}

func (t *AchieveTaskList) Give(achieveTaskData *taskdata.AchieveTaskData) (newTask *AchieveTask) {
	newTask = newAchieveTask(achieveTaskData)
	t.tasks[achieveTaskData.Id] = newTask
	return
}

// 创建获得检查成就是否存在的func
func (t *AchieveTaskList) IsAchieveTypeSelectShow(achieveType uint64) bool {
	return u64.Contains(t.selectShowAchieveTypes, achieveType)
}

// 展示数量
func (t *AchieveTaskList) ShowCount() int {
	return len(t.selectShowAchieveTypes)
}

// 添加
func (t *AchieveTaskList) AddSelectShowAchieveType(achieveType uint64) (isAdd bool) {
	isAdd, t.selectShowAchieveTypes = u64.TryAdd(t.selectShowAchieveTypes, achieveType)
	return
}

// 移除
func (t *AchieveTaskList) RemoveSelectShowAchieveType(achieveType uint64) {
	t.selectShowAchieveTypes = u64.RemoveIfPresent(t.selectShowAchieveTypes, achieveType)
}

func (t *AchieveTaskList) GetTaskByAchieveType(achieveType uint64) *AchieveTask {
	if len(t.tasks) <= 0 {
		return nil
	}

	for _, task := range t.tasks {
		if achieveType == task.data.AchieveType {
			return task
		}
	}

	return nil
}

func (t *AchieveTaskList) GetTasksByAchieveType(achieveTypes []uint64) (result []*AchieveTask) {
	if len(achieveTypes) <= 0 {
		return
	}

	if len(t.tasks) <= 0 {
		return
	}

	for _, task := range t.tasks {
		if u64.Contains(achieveTypes, task.data.AchieveType) {
			result = append(result, task)
		}
	}

	return
}

func (t *AchieveTaskList) unmarshal(proto *server_proto.AchieveTaskListServerProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	t.collectStarCount = proto.CollectStarCount

	for _, taskProto := range proto.List {
		data := datas.GetAchieveTaskData(taskProto.Id)
		if data == nil {
			continue
		}

		existData := datas.HeroInitData().AchieveTaskDatas[data.AchieveType]
		delete(t.tasks, existData.Id) // 删掉默认插入进去的

		bt := t.Give(data)
		bt.progress.progress = taskProto.Progress
		bt.isCollectPrize = taskProto.IsCollected
		bt.SetReachTime(timeutil.Unix64(taskProto.ReachTime))
		if bt.IsCollectPrize() && data.NextTask() != nil {
			t.Complete(bt, bt.ReachTime())
		}
	}

	t.selectShowAchieveTypes = make([]uint64, 0, len(proto.SelectShowAchieves))
	for _, achieveType := range proto.SelectShowAchieves {
		if datas.HeroInitData().AchieveTaskDatas[achieveType] != nil {
			// 存在了
			t.selectShowAchieveTypes = append(t.selectShowAchieveTypes, achieveType)
		}
	}
}

func (t *AchieveTaskList) encodeServer() *server_proto.AchieveTaskListServerProto {
	proto := &server_proto.AchieveTaskListServerProto{}

	for _, b := range t.tasks {
		proto.List = append(proto.List, b.encodeServer())
	}

	proto.CollectStarCount = t.collectStarCount

	proto.SelectShowAchieves = t.selectShowAchieveTypes

	return proto
}

func (t *AchieveTaskList) encodeClient() *shared_proto.AchieveTaskListProto {
	proto := &shared_proto.AchieveTaskListProto{List: make([]*shared_proto.AchieveTaskProto, 0, len(t.tasks))}

	for _, b := range t.tasks {
		proto.List = append(proto.List, b.encodeClient())
	}

	proto.CollectStarCount = u64.Int32Array(t.collectStarCount)

	proto.SelectShowAchieves = u64.Int32Array(t.selectShowAchieveTypes)

	return proto
}

func (t *AchieveTaskList) EncodeOther() *shared_proto.OtherAchieveTaskListProto {
	proto := &shared_proto.OtherAchieveTaskListProto{List: make([]*shared_proto.AchieveTaskProto, 0, len(t.tasks))}

	for _, b := range t.tasks {
		proto.List = append(proto.List, b.encodeClient())
	}

	proto.SelectShowAchieves = u64.Int32Array(t.selectShowAchieveTypes)

	return proto
}

func (t *AchieveTaskList) EncodeSelectShowAchieves() *shared_proto.SelectShowAchievesProto {
	proto := &shared_proto.SelectShowAchievesProto{}

	if len(t.selectShowAchieveTypes) > 0 {
		showTasks := t.GetTasksByAchieveType(t.selectShowAchieveTypes)
		proto.AchieveTaskId = make([]int32, 0, len(showTasks))
		proto.AchieveTaskReachTime = make([]int32, 0, len(showTasks))
		for _, task := range showTasks {
			if task.Data().PrevTask == nil && task.IsCollectPrize() {
				logrus.Errorf("玩家未解锁的成就在展示成就列表里面: %#v", task.Data())
				continue
			}

			if task.IsCollectPrize() {
				proto.AchieveTaskId = append(proto.AchieveTaskId, u64.Int32(task.Data().Id))
			} else {
				proto.AchieveTaskId = append(proto.AchieveTaskId, u64.Int32(task.Data().PrevTask.Id))
			}
			proto.AchieveTaskReachTime = append(proto.AchieveTaskReachTime, timeutil.Marshal32(task.ReachTime()))
		}
	}

	proto.TotalStar = u64.Int32(t.TotalStar())

	return proto
}

func (t *AchieveTaskList) IsPrizeCollected(starCount uint64) bool {
	return u64.Contains(t.collectStarCount, starCount)
}

func (t *AchieveTaskList) CollectPrize(starCount uint64) {
	t.collectStarCount = append(t.collectStarCount, starCount)
}

func newActiveDegreeTaskList(activeDegreeTaskDatas []*taskdata.ActiveDegreeTaskData) *ActiveDegreeTaskList {
	list := &ActiveDegreeTaskList{tasks: make(map[uint64]*ActiveDegreeTask)}
	for _, data := range activeDegreeTaskDatas {
		list.tasks[data.Id] = newActiveDegreeTask(data)
	}
	return list
}

// 成就任务
type ActiveDegreeTaskList struct {
	tasks               map[uint64]*ActiveDegreeTask
	collectedPrizeIndex []int32 // 领取了的活跃度奖励的下标
}

func (t *ActiveDegreeTaskList) Task(id uint64) *ActiveDegreeTask {
	return t.tasks[id]
}

func (t *ActiveDegreeTaskList) Walk(f func(taskId uint64, task *ActiveDegreeTask) (endWalk bool)) (isEndedWalk bool) {
	if len(t.tasks) <= 0 {
		return
	}

	for k, v := range t.tasks {
		if f(k, v) {
			return true
		}
	}

	return
}

// 计算活跃度
func (t *ActiveDegreeTaskList) Degree() (degree uint64) {
	t.Walk(func(taskId uint64, task *ActiveDegreeTask) bool {
		degree += task.Progress().GetProgress() * task.Data().AddDegree
		return false
	})

	return
}

// 领取活跃度奖励
func (t *ActiveDegreeTaskList) CollectPrize(index int32) {
	t.collectedPrizeIndex = append(t.collectedPrizeIndex, index)
}

// 活跃度奖励是否有领取
func (t *ActiveDegreeTaskList) IsPrizeCollected(index int32) bool {
	for _, collectedIndex := range t.collectedPrizeIndex {
		if collectedIndex == index {
			return true
		}
	}

	return false
}

func (t *ActiveDegreeTaskList) ResetDaily(hero *Hero) {
	t.collectedPrizeIndex = t.collectedPrizeIndex[:0]

	t.Walk(func(taskId uint64, task *ActiveDegreeTask) (endWalk bool) {
		task.progress.setProgress(0)
		task.Progress().UpdateTaskProgress(hero)
		return
	})
}

func (t *ActiveDegreeTaskList) unmarshal(proto *server_proto.ActiveDegreeTaskListServerProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	for _, taskProto := range proto.List {
		data := datas.GetActiveDegreeTaskData(taskProto.Id)
		if data == nil {
			continue
		}

		bt := newActiveDegreeTask(data)
		bt.progress.progress = taskProto.Progress
		t.tasks[bt.data.Id] = bt
	}

	t.collectedPrizeIndex = proto.CollectedPrizeIndex
}

func (t *ActiveDegreeTaskList) encodeServer() *server_proto.ActiveDegreeTaskListServerProto {
	proto := &server_proto.ActiveDegreeTaskListServerProto{}

	for _, b := range t.tasks {
		proto.List = append(proto.List, b.encodeServer())
	}

	proto.CollectedPrizeIndex = t.collectedPrizeIndex

	return proto
}

func (t *ActiveDegreeTaskList) encodeClient() *shared_proto.ActiveDegreeTaskListProto {
	proto := &shared_proto.ActiveDegreeTaskListProto{List: make([]*shared_proto.ActiveDegreeTaskProto, 0, len(t.tasks))}

	for _, b := range t.tasks {
		proto.List = append(proto.List, b.encodeClient())
	}

	proto.CollectedPrizeIndex = t.collectedPrizeIndex

	return proto
}

// achieve
func newBwzlTask(data *taskdata.BwzlTaskData) *BwzlTask {
	t := &BwzlTask{}
	t.data = data
	t.progress = &TaskProgress{}
	t.progress.target = data.Target

	return t
}

type BwzlTask struct {
	data *taskdata.BwzlTaskData

	progress *TaskProgress

	isCollectPrize bool
}

func (t *BwzlTask) Type() task_type {
	return BWZL_TASK
}

func (t *BwzlTask) NewTaskMsg() pbutil.Buffer {
	return task.NewS2cNewTaskMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.encodeClient()), t.progress.IsCompleted(), false)
}

func (t *BwzlTask) Data() *taskdata.BwzlTaskData {
	return t.data
}

func (t *BwzlTask) Progress() *TaskProgress {
	return t.progress
}

func (t *BwzlTask) Id() uint64 {
	return t.data.Id
}

func (t *BwzlTask) NewUpdateTaskProgressMsg() pbutil.Buffer {
	return task.NewS2cUpdateTaskProgressMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.GetProgress()), t.progress.IsCompleted())
}

func (t *BwzlTask) IsCollectPrize() bool {
	return t.isCollectPrize
}

func (t *BwzlTask) CollectPrize() {
	t.isCollectPrize = true
}

func (t *BwzlTask) encodeClient() *shared_proto.BwzlTaskProto {
	proto := &shared_proto.BwzlTaskProto{}

	proto.Id = u64.Int32(t.data.Id)
	proto.Progress = u64.Int32(t.progress.encodeClient())
	proto.IsCollected = t.IsCollectPrize()

	return proto
}

func (t *BwzlTask) encodeServer() *server_proto.BwzlTaskServerProto {
	proto := &server_proto.BwzlTaskServerProto{}

	proto.Id = t.data.Id
	proto.Progress = t.progress.encodeClient()
	proto.IsCollected = t.IsCollectPrize()

	return proto
}

func newBwzlTaskList(datas []*taskdata.BwzlTaskData) *BwzlTaskList {
	taskList := &BwzlTaskList{tasks: make(map[uint64]*BwzlTask, len(datas))}

	for _, data := range datas {
		taskList.Give(data)
	}

	return taskList
}

// 成就任务
type BwzlTaskList struct {
	tasks           map[uint64]*BwzlTask // 成就任务列表
	collectedPrizes []uint64             // 领取了的霸王之路的奖励
}

func (t *BwzlTaskList) Count() int {
	return len(t.tasks)
}

func (t *BwzlTaskList) Task(id uint64) *BwzlTask {
	return t.tasks[id]
}

func (t *BwzlTaskList) CollectPrizeCount() (result uint64) {
	t.Walk(func(taskId uint64, task *BwzlTask) (endWalk bool) {
		if task.IsCollectPrize() {
			result++
		}
		return
	})

	return
}

func (t *BwzlTaskList) Walk(f func(taskId uint64, task *BwzlTask) (endWalk bool)) (isEndedWalk bool) {
	if len(t.tasks) <= 0 {
		return
	}

	for k, v := range t.tasks {
		if f(k, v) {
			return true
		}
	}

	return
}

func (t *BwzlTaskList) Give(data *taskdata.BwzlTaskData) (newTask *BwzlTask) {
	newTask = newBwzlTask(data)
	t.tasks[data.Id] = newTask
	return
}

func (t *BwzlTaskList) IsCollectedPrizes(data *taskdata.BwzlPrizeData) bool {
	return u64.Contains(t.collectedPrizes, data.CollectPrizeTaskCount)
}

func (t *BwzlTaskList) CollectedPrize(data *taskdata.BwzlPrizeData) {
	t.collectedPrizes = u64.AddIfAbsent(t.collectedPrizes, data.CollectPrizeTaskCount)
}

func (t *BwzlTaskList) unmarshal(proto *server_proto.BwzlTaskListServerProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	for _, taskProto := range proto.List {
		task := t.Task(taskProto.Id)
		if task == nil {
			logrus.Errorf("玩家霸王之路任务数据丢失: %v", taskProto)
			continue
		}

		task.progress.progress = taskProto.Progress
		task.isCollectPrize = taskProto.IsCollected
	}

	t.collectedPrizes = proto.CollectedPrizes
}

func (t *BwzlTaskList) encodeServer() *server_proto.BwzlTaskListServerProto {
	proto := &server_proto.BwzlTaskListServerProto{}

	for _, b := range t.tasks {
		proto.List = append(proto.List, b.encodeServer())
	}

	proto.CollectedPrizes = t.collectedPrizes

	return proto
}

func (t *BwzlTaskList) encodeClient() *shared_proto.BwzlTaskListProto {
	proto := &shared_proto.BwzlTaskListProto{List: make([]*shared_proto.BwzlTaskProto, 0, len(t.tasks))}

	for _, b := range t.tasks {
		proto.List = append(proto.List, b.encodeClient())
	}

	proto.CollectedPrizes = u64.Int32Array(t.collectedPrizes)

	return proto
}

// 称号任务
func newTitleTaskList(firstData *taskdata.TitleData) *TitleTaskList {
	list := &TitleTaskList{
		doingData: firstData,
	}

	for _, v := range firstData.Task {
		list.tasks = append(list.tasks, newTitleTask(v))
	}

	return list
}

type TitleTaskList struct {
	doingData *taskdata.TitleData
	tasks     []*TitleTask

	completedData *taskdata.TitleData
}

func (t *TitleTaskList) encodeClient() *shared_proto.TitleTaskListProto {
	proto := &shared_proto.TitleTaskListProto{}

	if t.completedData != nil {
		proto.CompletedTitleId = u64.Int32(t.completedData.Id)
	}

	for _, v := range t.tasks {
		proto.List = append(proto.List, v.encodeClient())
	}

	return proto
}

func (t *TitleTaskList) encodeServer() *server_proto.TitleTaskListServerProto {
	proto := &server_proto.TitleTaskListServerProto{}

	if t.completedData != nil {
		proto.CompletedTitleId = t.completedData.Id
	}

	for _, v := range t.tasks {
		proto.List = append(proto.List, v.encodeServer())
	}

	return proto
}

func (t *TitleTaskList) unmarshal(proto *server_proto.TitleTaskListServerProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	data := datas.GetTitleData(proto.CompletedTitleId)
	if data != nil {
		t.completedData = data
		t.doingData = data.GetNextData()

		if t.doingData == nil {
			t.tasks = nil
		} else {
			for i, v := range t.doingData.Task {
				task := newTitleTask(v)
				t.tasks = append(t.tasks, task)

				// 先根据下标取一下（正常情况是可以取到的，如果取不到再遍历找一下）
				if i < len(proto.List) {
					if pt := proto.List[i]; pt.Id == task.data.Id {
						task.progress.progress = pt.Progress
						continue
					}
				}
				for _, pt := range proto.List {
					if pt.Id == task.data.Id {
						task.progress.progress = pt.Progress
						break
					}
				}
			}
		}
	}

}

func (t *TaskList) TitleTaskList() *TitleTaskList {
	return t.titleTaskList
}

func (t *TaskList) TitleId() uint64 {
	return t.titleTaskList.TitleId()
}

func (d *TaskList) GetTitleData() *taskdata.TitleData {
	return d.titleTaskList.completedData
}

func (d *TitleTaskList) TitleId() uint64 {
	if d.completedData != nil {
		return d.completedData.Id
	}
	return 0
}

func (d *TitleTaskList) GetDoingTitleData() *taskdata.TitleData {
	return d.doingData
}

func (d *TitleTaskList) IsAllTaskComplete() bool {
	if d.doingData == nil {
		return false
	}

	for _, v := range d.tasks {
		if !v.progress.IsCompleted() {
			return false
		}
	}
	return true
}

func (d *TitleTaskList) UpgradeTitleData() {

	if d.doingData == nil {
		return
	}

	d.completedData = d.doingData
	d.doingData = d.completedData.GetNextData()

	var tasks []*TitleTask
	if d.doingData != nil {
		for _, v := range d.doingData.Task {
			tasks = append(tasks, newTitleTask(v))
		}
	}
	d.tasks = tasks
}

func (t *TitleTaskList) Task(id uint64) *TitleTask {
	for _, v := range t.tasks {
		if v.data.Id == id {
			return v
		}
	}
	return nil
}

func (t *TitleTaskList) Walk(f func(taskId uint64, task *TitleTask) (endWalk bool)) (isEndedWalk bool) {
	if len(t.tasks) <= 0 {
		return
	}

	for _, v := range t.tasks {
		if f(v.data.Id, v) {
			return true
		}
	}

	return
}

func newTitleTask(data *taskdata.TitleTaskData) *TitleTask {
	t := &TitleTask{}
	t.data = data
	t.progress = &TaskProgress{}
	t.progress.target = data.Target

	return t
}

type TitleTask struct {
	data *taskdata.TitleTaskData

	progress *TaskProgress
}

func (t *TitleTask) Type() task_type {
	return TITLE_TASK
}

func (t *TitleTask) NewTaskMsg() pbutil.Buffer {
	return task.NewS2cNewTaskMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.encodeClient()), t.progress.IsCompleted(), false)
}

func (t *TitleTask) Progress() *TaskProgress {
	return t.progress
}

func (t *TitleTask) Id() uint64 {
	return t.data.Id
}

func (t *TitleTask) NewUpdateTaskProgressMsg() pbutil.Buffer {
	return task.NewS2cUpdateTaskProgressMsg(u64.Int32(t.data.Id), u64.Int32(t.progress.GetProgress()), t.progress.IsCompleted())
}

func (t *TitleTask) encodeClient() *shared_proto.TitleTaskProto {
	proto := &shared_proto.TitleTaskProto{}

	proto.Id = u64.Int32(t.data.Id)
	proto.Progress = u64.Int32(t.progress.encodeClient())

	return proto
}

func (t *TitleTask) encodeServer() *server_proto.TitleTaskServerProto {
	proto := &server_proto.TitleTaskServerProto{}

	proto.Id = t.data.Id
	proto.Progress = t.progress.encodeClient()

	return proto
}
