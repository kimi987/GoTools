package entity

import (
	"github.com/lightpaw/male7/config/promdata"
	"time"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/timeutil"
)

func newEventLimitGift(data *promdata.EventLimitGiftData, endTime time.Time) *eventLimitGift {
	m := &eventLimitGift{}
	m.data = data
	m.endTime = endTime
	return m
}

type eventLimitGift struct {
	data         *promdata.EventLimitGiftData
	endTime      time.Time
}

func (g *eventLimitGift) Data() *promdata.EventLimitGiftData {
	return g.data
}

func (g *eventLimitGift) Encode() *shared_proto.EventLimitGiftProto {
	return &shared_proto.EventLimitGiftProto {
		Id: u64.Int32(g.data.Id),
		EndTime: timeutil.Marshal32(g.endTime),
		Prize: g.data.Prize.Encode(),
	}
}

func (g *eventLimitGift) EncodeServer() *shared_proto.EventLimitGiftProto {
	return &shared_proto.EventLimitGiftProto {
		Id: u64.Int32(g.data.Id),
		EndTime: timeutil.Marshal32(g.endTime),
	}
}
