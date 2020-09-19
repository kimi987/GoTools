package entity

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/extratimesservice/extratimesface"
	"github.com/lightpaw/male7/util/u64"
)

type UsedExtraTimes interface {
	Times(maxTimes extratimesface.ExtraMaxTimes) uint64
	UsedTimes() uint64
	ReduceOneTimes(maxTimes extratimesface.ExtraMaxTimes) (suc bool, extraTimesType shared_proto.ExtraTimesType, newUsedTimes uint64)
	UsedExtraTimes(extraTimesType shared_proto.ExtraTimesType) uint64
	SetUsedExtraTimes(extraTimesType shared_proto.ExtraTimesType, toSet uint64)
	Reset(extraTimesType shared_proto.ExtraTimesType)
	Clear()
	Encode() *shared_proto.ExtraTimesListProto
}

func NewUsedExtraTimes() UsedExtraTimes {
	return &usedExtraTimes{}
}

// 可恢复的次数跟额外的次数
type usedExtraTimes struct {
	usedTimes []*shared_proto.ExtraTimesProto
}

func (r *usedExtraTimes) UsedTimes() (result uint64) {
	for _, proto := range r.usedTimes {
		result = u64.Plus(result, u64.FromInt32(proto.UsedTimes))
	}
	return
}

func (r *usedExtraTimes) Times(maxTimes extratimesface.ExtraMaxTimes) uint64 {
	if len(r.usedTimes) <= 0 {
		return maxTimes.TotalTimes()
	}

	var times uint64

	maxTimes.Walk(func(extraTimesType shared_proto.ExtraTimesType, maxTimes uint64) (endWalk bool) {
		if maxTimes <= 0 {
			return
		}

		usedTimes := r.UsedExtraTimes(extraTimesType)
		times += u64.Sub(maxTimes, usedTimes)
		return
	})

	return times
}

func (r *usedExtraTimes) ReduceOneTimes(maxTimes extratimesface.ExtraMaxTimes) (suc bool, extraTimesType shared_proto.ExtraTimesType, newUsedTimes uint64) {
	maxTimes.Walk(func(extraTimesType shared_proto.ExtraTimesType, maxTimes uint64) (endWalk bool) {
		if maxTimes <= 0 {
			return
		}

		usedTimes := r.UsedExtraTimes(extraTimesType)
		if usedTimes >= maxTimes {
			return
		}

		suc = true
		newUsedTimes = usedTimes + 1
		r.SetUsedExtraTimes(extraTimesType, newUsedTimes)
		return true
	})

	return
}

func (r *usedExtraTimes) UsedExtraTimes(extraTimesType shared_proto.ExtraTimesType) uint64 {
	for _, value := range r.usedTimes {
		if value.Type == extraTimesType {
			return u64.FromInt32(value.UsedTimes)
		}
	}
	return 0
}

func (r *usedExtraTimes) SetUsedExtraTimes(extraTimesType shared_proto.ExtraTimesType, toSet uint64) {
	if toSet <= 0 {
		return
	}

	for _, value := range r.usedTimes {
		if value.Type == extraTimesType {
			value.UsedTimes = u64.Int32(toSet)
			return
		}
	}

	r.usedTimes = append(r.usedTimes, &shared_proto.ExtraTimesProto{Type: extraTimesType, UsedTimes: u64.Int32(toSet)})
}

func (r *usedExtraTimes) Reset(extraTimesType shared_proto.ExtraTimesType) {
	r.SetUsedExtraTimes(extraTimesType, 0)
}

func (r *usedExtraTimes) Clear() {
	r.usedTimes = r.usedTimes[:0]
}

func (r *usedExtraTimes) Encode() *shared_proto.ExtraTimesListProto {
	return &shared_proto.ExtraTimesListProto{TypeList: r.usedTimes}
}
