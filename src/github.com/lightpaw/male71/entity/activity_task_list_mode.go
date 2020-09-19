package entity

import (
	"time"
	"github.com/lightpaw/male7/config/activitydata"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/gen/pb/activity"
	"github.com/lightpaw/male7/util/timeutil"
)

// 列表式任务活动单例
type ActivityListModeTasks struct {
	data      *activitydata.ActivityTaskListModeData
	startTime time.Time
	endTime   time.Time
}

func (a *ActivityListModeTasks) equal(al *ActivityListModeTasks) bool {
	return a.data.Equal(al.data) && a.startTime.Equal(al.startTime) && a.endTime.Equal(al.endTime)
}

func (a *ActivityListModeTasks) encode() *shared_proto.ActivityTaskListModeProto {
	return &shared_proto.ActivityTaskListModeProto {
		Data: a.data.Encode(),
		EndTime: timeutil.Marshal32(a.endTime),
	}
}

// 一堆列表式任务活动
func NewActivityListModeTasksList(datas []*activitydata.ActivityTaskListModeData, serverStartTime, ctime time.Time) *ActivityListModeTasksList {
	m := &ActivityListModeTasksList{}
	for _, data := range datas {
		ok, start, end := data.GetRecentActivityTime(serverStartTime, ctime)
		if !ok || !timeutil.Between(ctime, start, end) {
			continue
		}
		m.list = append(m.list, &ActivityListModeTasks {
			data: data,
			startTime: start,
			endTime: end,
		})
	}
	return m
}

type ActivityListModeTasksList struct {
	list  []*ActivityListModeTasks
	msg   pbutil.Buffer
	toMap map[uint64]*ActivityListModeTasks
}

func (a *ActivityListModeTasksList) Msg() pbutil.Buffer {
	return a.msg
}

func (a *ActivityListModeTasksList) GetMap() map[uint64]*ActivityListModeTasks {
	if a.toMap == nil {
		a.toMap = make(map[uint64]*ActivityListModeTasks, len(a.list))
		for _, m := range a.list {
			a.toMap[m.data.Id] = m
		}
	}
	return a.toMap
}

func (a *ActivityListModeTasksList) Equal(al *ActivityListModeTasksList) bool {
	size := len(a.list)
	if size != len(al.list) {
		return false
	}
	for i := 0; i < size; i++ {
		if !a.list[i].equal(al.list[i]) {
			return false
		}
	}
	return true
}

func (a *ActivityListModeTasksList) RefreshMsg() {
	if size := len(a.list); size > 0 {
		protos := make([]*shared_proto.ActivityTaskListModeProto, size)
		for i := 0; i < size; i++ {
			protos[i] = a.list[i].encode()
		}
		a.msg = activity.NewS2cNoticeTaskListModeMsg(protos).Static()
	}
}
