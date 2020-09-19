package extratimesface

import "github.com/lightpaw/male7/pb/shared_proto"

type ExtraMaxTimes interface {
	Walk(walkFunc func(extraTimesType shared_proto.ExtraTimesType, maxTimes uint64) (endWalk bool))
	MaxTimesByTime(extraTimesType shared_proto.ExtraTimesType) (maxTimes uint64)
	TotalTimes() (times uint64)
}
