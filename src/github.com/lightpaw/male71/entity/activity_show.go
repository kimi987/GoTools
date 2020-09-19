package entity

import (
	"time"
	"github.com/lightpaw/male7/config/activitydata"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/gen/pb/activity"
)

// 活动展示
func NewActivityShowList(datas []*activitydata.ActivityShowData, serverStartTime, ctime time.Time) *ActivityShowList {
	m := &ActivityShowList{}
	for _, data := range datas {
		ok, startT, endT := data.GetRecentShowTime(serverStartTime, ctime)
		if !ok || !timeutil.Between(ctime, startT, endT) {
			continue
		}
		m.shows = append(m.shows, &ActivityShow{
			data:      data,
			startTime: startT,
			endTime:   endT,
		})
	}
	return m
}

type ActivityShowList struct {
	shows []*ActivityShow // 所有的展示
	msg   pbutil.Buffer
}

func (sl *ActivityShowList) Msg() pbutil.Buffer {
	return sl.msg
}

func (sl *ActivityShowList) Equal(list *ActivityShowList) bool {
	sLen := len(sl.shows)
	if sLen != len(list.shows) {
		return false
	}
	for i := 0; i < sLen; i++ {
		if !sl.shows[i].equal(list.shows[i]) {
			return false
		}
	}
	return true
}

func (s *ActivityShowList) RefreshMsg() {
	if size := len(s.shows); size > 0 {
		protos := make([]*shared_proto.ActiviyShowProto, size)
		for i := 0; i < size; i++ {
			protos[i] = s.shows[i].encode()
		}
		s.msg = activity.NewS2cNoticeActivityShowMsg(protos).Static()
	}
}

type ActivityShow struct {
	data      *activitydata.ActivityShowData
	startTime time.Time
	endTime   time.Time
}

func (s *ActivityShow) equal(show *ActivityShow) bool {
	return s.data.Equal(show.data) && s.startTime.Equal(show.startTime) && s.endTime.Equal(show.endTime)
}

func (s *ActivityShow) encode() *shared_proto.ActiviyShowProto {
	return &shared_proto.ActiviyShowProto{
		Data:    s.data.Encode(),
		EndTime: timeutil.Marshal32(s.endTime),
	}
}
