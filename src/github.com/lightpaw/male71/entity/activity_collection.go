package entity

import (
	"time"
	"github.com/lightpaw/male7/config/activitydata"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/gen/pb/activity"
)

// 收集活动单例
type ActivityCollection struct {
	data      *activitydata.ActivityCollectionData
	startTime time.Time
	endTime   time.Time
}

func (a *ActivityCollection) Data() *activitydata.ActivityCollectionData {
	return a.data
}

func (a *ActivityCollection) Id() uint64 {
	return a.data.Id
}

func (a *ActivityCollection) equal(ac *ActivityCollection) bool {
	return a.data.Equal(ac.data) && a.startTime.Equal(ac.startTime) && a.endTime.Equal(ac.endTime)
}

func (a *ActivityCollection) encode() *shared_proto.ActivityCollectionProto {
	return &shared_proto.ActivityCollectionProto{
		Data:    a.data.Encode(),
		EndTime: timeutil.Marshal32(a.endTime),
	}
}

// 一堆收集活动
func NewActivityCollectionList(datas []*activitydata.ActivityCollectionData, serverStartTime, ctime time.Time) *ActivityCollectionList {
	m := &ActivityCollectionList{}
	for _, data := range datas {
		ok, start, end := data.GetRecentActivityTime(serverStartTime, ctime)
		if !ok || !timeutil.Between(ctime, start, end) {
			continue
		}
		m.list = append(m.list, &ActivityCollection{
			data:      data,
			startTime: start,
			endTime:   end,
		})
	}
	return m
}

type ActivityCollectionList struct {
	list  []*ActivityCollection
	msg   pbutil.Buffer
	toMap map[uint64]*ActivityCollection
}

func (a *ActivityCollectionList) Msg() pbutil.Buffer {
	return a.msg
}

func (a *ActivityCollectionList) GetActivity(id uint64) *ActivityCollection {
	return a.GetMap()[id]
}

func (a *ActivityCollectionList) GetMap() map[uint64]*ActivityCollection {
	if a.toMap == nil {
		a.toMap = make(map[uint64]*ActivityCollection, len(a.list))
		for _, c := range a.list {
			a.toMap[c.data.Id] = c
		}
	}
	return a.toMap
}

func (a *ActivityCollectionList) Equal(acl *ActivityCollectionList) bool {
	size := len(a.list)
	if size != len(acl.list) {
		return false
	}
	for i := 0; i < size; i++ {
		if !a.list[i].equal(acl.list[i]) {
			return false
		}
	}
	return true
}

func (a *ActivityCollectionList) RefreshMsg() {
	if size := len(a.list); size > 0 {
		protos := make([]*shared_proto.ActivityCollectionProto, size)
		for i := 0; i < size; i++ {
			protos[i] = a.list[i].encode()
		}
		a.msg = activity.NewS2cNoticeCollectionMsg(protos).Static()
	}
}
