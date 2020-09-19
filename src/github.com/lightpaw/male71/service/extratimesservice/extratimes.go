package extratimesservice

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/extratimesservice/extratimesface"
)

// 额外次数

func NewExtraTimesService(seasonService iface.SeasonService, timeService iface.TimeService) *ExtraTimesService {
	multiLevelNpcMaxTimes := make(extraMaxTimes)
	multiLevelNpcMaxTimes[shared_proto.ExtraTimesType_Ett_Season] = func() (times uint64) { return seasonService.Season().AddMultiMonsterTimes }

	return &ExtraTimesService{
		multiLevelNpcMaxTimes: multiLevelNpcMaxTimes,
	}
}

//gogen:iface
type ExtraTimesService struct {
	multiLevelNpcMaxTimes extratimesface.ExtraMaxTimes
}

func (s *ExtraTimesService) MultiLevelNpcMaxTimes() extratimesface.ExtraMaxTimes {
	return s.multiLevelNpcMaxTimes
}

type extraMaxTimes map[shared_proto.ExtraTimesType]func() (times uint64)

func (et extraMaxTimes) Walk(walkFunc func(extraTimesType shared_proto.ExtraTimesType, maxTimes uint64) (endWalk bool)) {
	if len(et) <= 0 {
		return
	}

	for extraTimesType, maxTimesFunc := range et {
		if walkFunc(extraTimesType, maxTimesFunc()) {
			return
		}
	}
}

func (et extraMaxTimes) MaxTimesByTime(extraTimesType shared_proto.ExtraTimesType) (maxTimes uint64) {
	f := et[extraTimesType]
	if f == nil {
		return
	}

	return f()
}

func (et extraMaxTimes) TotalTimes() (times uint64) {
	if len(et) <= 0 {
		return
	}

	for _, maxTimesFunc := range et {
		times += maxTimesFunc()
	}

	return
}
