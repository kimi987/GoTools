package guild_data

import (
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/u64"
)

// 联盟任务评价奖励
//gogen:config
type GuildTaskEvaluateData struct {
	_ struct{} `file:"联盟/联盟任务评价.txt"`
	_ struct{} `proto:"shared_proto.GuildTaskEvaluateDataProto"`
	_ struct{} `protoconfig:"guild_task_evaluate"`

	Id        uint64
	Name      string
	Complete  uint64
	Prizes    []*resdata.Prize `validator:"uint,notAllNil,duplicate"`
}

// 联盟任务
//gogen:config
type GuildTaskData struct {
	_ struct{} `file:"联盟/联盟任务.txt"`
	_ struct{} `proto:"shared_proto.GuildTaskDataProto"`
	_ struct{} `protoconfig:"guild_task"`

	Id        uint64
	TaskType  server_proto.GuildTaskType `head:"-" protofield:"-"`
	Name      string
	Desc      string
	Icon      string
	Stages    []uint64         `validator:"uint,notAllNil"`
	Prizes    []*resdata.Prize `validator:"uint,notAllNil,duplicate"`
}

func (d *GuildTaskData) Init(filename string) {
	check.PanicNotTrue(len(d.Stages) == len(d.Prizes), "%s, 阶段和阶段奖励数组长度不一致", filename)
	i := 0
	max := len(d.Stages)-1
	for ; i < max; i++ {
		check.PanicNotTrue(d.Stages[i] < d.Stages[i+1], "%s, 前一个阶段必须比后一个阶段数字小", filename)
	}
	t := u64.Int32(d.Id)
	_, ok := server_proto.GuildTaskType_name[t]
	check.PanicNotTrue(ok, "%s, 未知的任务类型：%v", filename, d.Id)
	d.TaskType = server_proto.GuildTaskType(t)
}

func (d *GuildTaskData) GetCompletedStageCount(progress uint64) uint64 {
	count := uint64(0)
	for _, stage := range d.Stages {
		if progress < stage {
			break
		}
		count++
	}
	return count
}

func (d *GuildTaskData) GetStageIndex(progress uint64, fromIndex int) (newIndex int) {
	newIndex = fromIndex
	for l := len(d.Stages); newIndex < l; newIndex++ {
		if progress < d.Stages[newIndex] {
			break
		}
	}
	return
}
