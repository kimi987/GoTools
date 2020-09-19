package guildsnapshot

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/util/event"
	"runtime/debug"
	"github.com/lightpaw/male7/pb/shared_proto"
)

func NewGuildSnapshotService() *GuildSnapshotService {
	s := &GuildSnapshotService{
		dataMap:       NewSnapshotMap(),
		callbackQueue: event.NewFuncQueue(1024, "GuildSnapshotService.loop()"),
	}

	return s
}

//gogen:iface
type GuildSnapshotService struct {
	dataMap *SnapshotMap

	callbacks     []guildsnapshotdata.Callback
	callbackQueue *event.FuncQueue
}

// 注册监听snapshot变化callback
func (s *GuildSnapshotService) RegisterCallback(callback guildsnapshotdata.Callback) {
	s.callbacks = append(s.callbacks, callback)
}

func  (s *GuildSnapshotService) GetGuildBasicProto(gid int64) *shared_proto.GuildBasicProto {
	if g := s.GetSnapshot(gid); g != nil {
		return g.BasicProto()
	} else {
		return nil
	}
}

func (m *GuildSnapshotService) GetSnapshot(id int64) *guildsnapshotdata.GuildSnapshot {
	if id != 0 {
		s, _ := m.dataMap.Get(id)
		return s
	}
	return nil
}

func (m *GuildSnapshotService) GetGuildLevel(id int64) uint64 {
	if g := m.GetSnapshot(id); g != nil {
		return g.GuildLevel.Level
	}
	return 0
}

func (m *GuildSnapshotService) UpdateSnapshot(g *guildsnapshotdata.GuildSnapshot) {
	origin := m.dataMap.Set(g.Id, g)

	m.callbackQueue.TryFunc(func() {
		for _, callback := range m.callbacks {
			safeUpdateCallback(callback, origin, g)
		}
	})
}

func (m *GuildSnapshotService) RemoveSnapshot(id int64) {
	m.dataMap.Remove(id)

	m.callbackQueue.TryFunc(func() {
		for _, callback := range m.callbacks {
			safeRemoveCallback(callback, id)
		}
	})
}

func safeUpdateCallback(callback guildsnapshotdata.Callback, origin, update *guildsnapshotdata.GuildSnapshot) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("GuildSnapshot.safeUpdateCallback recovered from panic. SEVERE!!!")
			metrics.IncPanic()
		}
	}()

	callback.OnGuildSnapshotUpdated(origin, update)
}

func safeRemoveCallback(callback guildsnapshotdata.Callback, id int64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("GuildSnapshot.safeRemoveCallback recovered from panic. SEVERE!!!")
			metrics.IncPanic()
		}
	}()

	callback.OnGuildSnapshotRemoved(id)
}
