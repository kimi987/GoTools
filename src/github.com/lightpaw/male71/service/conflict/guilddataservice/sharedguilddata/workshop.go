package sharedguilddata

import (
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/collection"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/region"
)

const log_capacity = 10

// 联盟工坊
func NewDefaultWorkshop(startTime, endTime int32, x, y int, prosperity uint64) *Workshop {
	return &Workshop{
		startTime:  startTime,
		endTime:    endTime,
		x:          x,
		y:          y,
		prosperity: prosperity,
		log:        collection.NewRingList(log_capacity),
		logVersion: 1,
	}
}

func NewWorkshop(proto *server_proto.GuildWorkshopServerProto) *Workshop {
	p := NewDefaultWorkshop(proto.StartTime, proto.EndTime, int(proto.X), int(proto.Y), proto.Prosperity)
	p.isComplete = proto.IsComplete
	if len(proto.Log) > 0 {
		for _, v := range proto.Log {
			p.log.Add(v)
		}
	}
	p.refreshMsg()
	return p
}

type Workshop struct {
	// 开始时间
	startTime int32

	// 竣工时间
	endTime int32

	// 是否已竣工
	isComplete bool

	// 坐标
	x, y int

	// 繁荣度
	prosperity uint64

	// 联盟工坊日志
	log        *collection.RingList
	logVersion uint64 // 日志版本号
	msg        pbutil.Buffer
	emptyMsg   pbutil.Buffer // 客户端版本号相同时发的消息（内容为空）
}

func (g *Guild) GetWorkshop() *Workshop {
	return g.workshop
}

func (g *Guild) SetWorkshop(toSet *Workshop) {
	g.workshop = toSet

	if toSet == nil {
		g.WalkMember(func(member *GuildMember) {
			member.SetShowWorkshopNotExist(false)
		})
	}
}

func (g *Guild) GetTodayWorkshopBeenHurtTimes() uint64 {
	return g.workshopBeenHurtTimes
}

func (g *Guild) IncTodayWorkshopBeenHurtTimes() uint64 {
	g.workshopBeenHurtTimes++
	return g.workshopBeenHurtTimes
}

func (w *Workshop) SetData(x, y int, startTime, endTime int32, isComplete bool, prosperity uint64) {
	w.x = x
	w.y = y
	w.startTime = startTime
	w.endTime = endTime
	w.isComplete = isComplete
	w.prosperity = prosperity
}

func (w *Workshop) GetBaseXY() (int, int) {
	return w.x, w.y
}

func (w *Workshop) GetTime() (int32, int32) {
	return w.startTime, w.endTime
}

func (w *Workshop) SetTime(startTime, endTime int32) {
	w.startTime = startTime
	w.endTime = endTime
}

func (w *Workshop) IsComplete() bool {
	return w.isComplete
}

func (g *Guild) GetWorkshopTodayCompleted() bool {
	return g.workshopTodayCompleted
}

func (g *Guild) Complete(prosperity uint64) {

	g.workshopOutput = 0
	g.workshopTodayCompleted = true

	if g.workshop != nil {
		g.workshop.isComplete = true
		g.workshop.prosperity = prosperity
	}
}

func (g *Guild) GetWorkshopOutput() uint64 {
	return g.workshopOutput
}

func (g *Guild) GetWorkshopOutputPrizeCount() int {
	return u64.Int(g.workshopOutputPrizeCount)
}

func (g *Guild) AddWorkshopOutput(toAdd uint64) uint64 {
	g.workshopOutput += toAdd
	return g.workshopOutput
}

func (g *Guild) TryWorkshopOutput(maxOutput []uint64) bool {
	if g.GetWorkshopOutputPrizeCount() >= len(maxOutput) {
		return false
	}

	outputLimit := maxOutput[g.workshopOutputPrizeCount]

	if g.workshopOutput >= outputLimit {
		g.workshopOutput = 0
		g.workshopOutputPrizeCount++
		return true
	}
	return false
}

func (w *Workshop) GetProsperity() uint64 {
	return w.prosperity
}

func (w *Workshop) SetProsperity(toSet uint64) {
	w.prosperity = toSet
}

func (w *Workshop) getLog() []*shared_proto.GuildWorkshopLogProto {
	p := []*shared_proto.GuildWorkshopLogProto{}
	if w.log.Length() > 0 {
		w.log.ReverseRange(func(v interface{}) (toContinue bool) {
			log := v.(*shared_proto.GuildWorkshopLogProto)
			p = append(p, log)
			return true
		})
	}
	return p
}

func (w *Workshop) refreshMsg() {
	w.emptyMsg = region.NewS2cCatchGuildWorkshopLogsMsg(u64.Int32(w.logVersion), nil).Static()
	w.msg = region.NewS2cCatchGuildWorkshopLogsMsg(u64.Int32(w.logVersion), w.getLog()).Static()
}

func (w *Workshop) GetMsg(v uint64) pbutil.Buffer {
	if v == w.logVersion {
		return w.emptyMsg
	}
	return w.msg
}

func (w *Workshop) AddLog(text string, t time.Time) {
	w.log.Add(&shared_proto.GuildWorkshopLogProto{
		Time: timeutil.Marshal32(t),
		Text: text,
	})
	w.logVersion++
	w.refreshMsg()
}

func (w *Workshop) LogVersion() uint64 {
	return w.logVersion
}
