package status

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"time"
	"context"
	"github.com/lightpaw/male7/util/ctxfunc"
)

func NewServerStartStopTimeService(db iface.DbService, timeService iface.TimeService) *ServerStartStopTimeService {
	return &ServerStartStopTimeService{
		db:            db,
		timeService:   timeService,
		lastStartTime: loadTime(db, server_proto.Key_ServerStartTime),
		lastStopTime:  loadTime(db, server_proto.Key_ServerStopTime),
	}
}

func loadTime(db iface.DbService, key server_proto.Key) time.Time {

	var startBytes []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		startBytes, err = db.LoadKey(ctx, key)
		return
	})
	if err != nil {
		logrus.WithError(err).Panicln("开启服务器 NewServerStartStopTimeService 报错")
	}

	if len(startBytes) > 0 {
		startTimeProto := server_proto.TimeProto{}
		startTimeProto.Unmarshal(startBytes)
		return timeutil.Unix64(startTimeProto.Time)
	}

	return time.Time{}
}

//gogen:iface
type ServerStartStopTimeService struct {
	db            iface.DbService
	timeService   iface.TimeService
	lastStartTime time.Time // 上次开启的时间
	lastStopTime  time.Time // 上次关闭的时间
}

func (s *ServerStartStopTimeService) IsNormalStop() bool {
	return s.lastStopTime.After(s.lastStartTime)
}

func (s *ServerStartStopTimeService) SaveStartTime() {
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		return s.db.SaveKey(ctx, server_proto.Key_ServerStartTime, must.Marshal(&server_proto.TimeProto{
			Time: timeutil.Marshal64(s.timeService.CurrentTime()),
		}))
	})
	if err != nil {
		logrus.WithError(err).Errorf("保存开服时间数据出错")
	}
}

func (s *ServerStartStopTimeService) SaveStopTime() {
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		return s.db.SaveKey(ctx, server_proto.Key_ServerStopTime, must.Marshal(&server_proto.TimeProto{
			Time: timeutil.Marshal64(s.timeService.CurrentTime()),
		}))
	})
	if err != nil {
		logrus.WithError(err).Errorf("保存关服时间数据出错")
	}
}
