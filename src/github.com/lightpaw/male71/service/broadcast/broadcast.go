package broadcast

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/timer"
	"time"
)

var (
	secondTimeWheel = timer.NewTimingWheel(500*time.Millisecond, 32)
)

func NewBroadcastService(guild iface.GuildSnapshotService, world iface.WorldService, chat iface.ChatService,
	heroSnapshot iface.HeroSnapshotService, datas iface.ConfigDatas, timeService iface.TimeService) *BroadcastService {
	m := &BroadcastService{
		closeNotify:    make(chan struct{}),
		loopExitNotify: make(chan struct{}),
	}
	m.guild = guild
	m.world = world
	m.datas = datas
	m.chat = chat
	m.timeService = timeService
	m.heroSnapshot = heroSnapshot

	go call.CatchLoopPanic(m.loop, "BroadcastService.loop")
	return m
}

//gogen:iface
type BroadcastService struct {
	guild        iface.GuildSnapshotService
	world        iface.WorldService
	heroSnapshot iface.HeroSnapshotService
	datas        iface.ConfigDatas
	timeService  iface.TimeService
	chat         iface.ChatService

	closeNotify    chan struct{}
	loopExitNotify chan struct{}
}

// ========== 系统定时广播 =============

func (m *BroadcastService) loop() {
	defer close(m.loopExitNotify)

	secondTick := secondTimeWheel.After(3 * time.Second)
	for {
		select {
		case <-secondTick:
			secondTick = secondTimeWheel.After(3 * time.Second)
			call.CatchPanic(m.sendAllTimingBroadcast, "broadcast.sendAllTimingBroadcast")
		case <-m.closeNotify:
			return
		}
	}
}

func (m *BroadcastService) Close() {
	close(m.closeNotify)
	<-m.loopExitNotify
}

func (m *BroadcastService) sendAllTimingBroadcast() {
	ctime := m.timeService.CurrentTime()
	for _, d := range m.datas.GetBroadcastDataArray() {
		if d.BcType != shared_proto.BCType_BC_TIMING_SYS {
			continue
		}
		if nextTime := d.GetNextTime(ctime); nextTime.Before(ctime) {
			m.world.Broadcast(d.TimingBroadcastMsg)
		}
	}
}

func (m *BroadcastService) GetEquipText(equipData *goods.EquipmentData) string {
	qid := equipData.Quality.GoodsQuality.Quality
	return m.getColorText(qid, equipData.Name)
}

func (m *BroadcastService) GetCaptainText(captain *entity.Captain) string {
	return m.getColorText(captain.Quality(), captain.Name())
}

func (m *BroadcastService) getQualityText(qid uint64) string {
	d := m.datas.ColorData().Must(qid)
	return m.getColorText(shared_proto.Quality(qid), d.ColorName)
}

func (m *BroadcastService) getColorText(qid shared_proto.Quality, text string) string {
	d := m.datas.GetColorData(uint64(qid))
	if d == nil {
		logrus.Debugf("系统广播，getColorText 收到未知物品品质：%v", qid)
		return text
	}

	return "[color=" + d.ColorCode + "]" + text + "[/color]"
}

func (m *BroadcastService) Broadcast(text string, sendChat bool) {
	if sendChat {
		m.chat.BroadcastSystemChat(text)
	} else {
		m.world.Broadcast(misc.NewS2cSysBroadcastMsg(text).Static())
	}
}
