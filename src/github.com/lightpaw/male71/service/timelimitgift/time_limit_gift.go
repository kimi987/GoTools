package timelimitgift

import (
	"time"
	"github.com/lightpaw/male7/config/promdata"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/timer"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/promotion"
)

var (
	minuteTimeWheel = timer.NewTimingWheel(time.Second * 30, 32)
)

func newTimeLimitGift(data *promdata.TimeLimitGiftGroupData, endTime time.Time) *timeLimitGift {
	m := &timeLimitGift{}
	m.data = data
	m.endTime = endTime
	return m
}

type timeLimitGift struct {
	data         *promdata.TimeLimitGiftGroupData
	endTime      time.Time
}

func (g *timeLimitGift) encode() *shared_proto.TimeLimitGiftProto {
	return &shared_proto.TimeLimitGiftProto {
		Id: u64.Int32(g.data.Id),
		EndTime: timeutil.Marshal32(g.endTime),
	}
}

func NewTimeLimitGiftService(dep iface.ServiceDep) *TimeLimitGiftService {
	t := &TimeLimitGiftService{}
	t.dep = dep
	t.closeNotify = make(chan struct{})
	t.loopExitNotify = make(chan struct{})
	t.gifts = make(map[uint64]*timeLimitGift)

	t.update()
	heromodule.RegisterHeroOnlineListener(t)

	go call.CatchLoopPanic(t.loop, "TimeLimitGiftService.Loop")
	return t
}

//gogen:iface
type TimeLimitGiftService struct {
	dep            iface.ServiceDep
	closeNotify    chan struct{}
	loopExitNotify chan struct{}
	// 时限礼包列表
	gifts          map[uint64]*timeLimitGift
	broadcastMsg   pbutil.Buffer
}

func (t *TimeLimitGiftService) update() {
	ctime := t.dep.Time().CurrentTime()
	changed := false
	for _, data := range t.dep.Datas().GetTimeLimitGiftGroupDataArray() {
		gift := t.gifts[data.Id]
		if gift == nil {
			if endTime, ok := data.IsCanOpenUp(t.dep.SvrConf().GetServerStartTime(), ctime); ok {
				t.gifts[data.Id] = newTimeLimitGift(data, endTime)
				changed = true
			}
		} else if ctime.After(gift.endTime) {
			delete(t.gifts, data.Id)
			changed = true
		}
	}
	if changed { // 生成新的msg并且广播
		t.broadcastMsg = promotion.NewS2cNoticeTimeLimitGiftsMsg(t.EncodeClient()).Static()
		t.dep.World().Broadcast(t.broadcastMsg)
	}
}

func (t *TimeLimitGiftService) loop() {
	defer close(t.loopExitNotify)

	minuteTick := minuteTimeWheel.After(time.Minute)
	for {
		select {
		case <-minuteTick:
			minuteTick = minuteTimeWheel.After(time.Minute)
			t.update()
		case <-t.closeNotify:
			return
		}
	}
}

func (t *TimeLimitGiftService) Close() {
	close(t.closeNotify)
	<-t.loopExitNotify
}

func (t *TimeLimitGiftService) EncodeClient() []*shared_proto.TimeLimitGiftProto {
	proto := []*shared_proto.TimeLimitGiftProto{}
	for _, gift := range t.gifts {
		proto = append(proto, gift.encode())
	}
	return proto
}

func (t *TimeLimitGiftService) GetGiftEndTime(id uint64) (endTime time.Time, ok bool) {
	gift := t.gifts[id]
	if gift != nil {
		endTime = gift.endTime
		ok = true
	}
	return
}

func (t *TimeLimitGiftService) OnHeroOnline(hc iface.HeroController) {
	hc.Send(t.broadcastMsg)
}
