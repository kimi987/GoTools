package tlog

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/config/kv"
	"sync"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/util/call"
	"time"
	"bytes"
	"strconv"
	"strings"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/monitor/metrics"
)

//gogen:iface
type TlogService struct {
	tlogBaseService     iface.TlogBaseService
	worldService        iface.WorldService
	heroSnapshotService iface.HeroSnapshotService
	config              *kv.IndividualServerConfig
	timeService         iface.TimeService

	dontGenTlog bool

	closeChan chan struct{}
	loopChan  chan struct{}

	logChan chan string

	logCacheMux sync.Mutex
	logCache    []string
	maxLogSize  int
}

func NewTlogService(worldService iface.WorldService, heroSnapshotService iface.HeroSnapshotService, tlogBaseService iface.TlogBaseService, timeService iface.TimeService, config *kv.IndividualServerConfig) *TlogService {
	s := &TlogService{}
	s.worldService = worldService
	s.tlogBaseService = tlogBaseService
	s.heroSnapshotService = heroSnapshotService
	s.config = config
	s.timeService = timeService

	s.dontGenTlog = !config.TlogStart && !config.KafkaStart

	s.closeChan = make(chan struct{})
	s.loopChan = make(chan struct{})

	s.logChan = make(chan string, 2048)
	s.maxLogSize = 8192

	go call.CatchLoopPanic(s.loop, "TlogService.loop()")

	return s
}

func (s *TlogService) DontGenTlog() bool {
	return s.dontGenTlog
}

func (s *TlogService) Close() {
	close(s.closeChan)
	<-s.loopChan

	close(s.logChan)

	if len(s.logChan) > 0 {
		for data := range s.logChan {
			s.tlogBaseService.WriteTlog(data)
		}
	}

	// 处理缓存数据
	s.handle()

	s.tlogBaseService.Close()
}

func (s *TlogService) loop() {

	defer close(s.loopChan)

	tick := time.NewTicker(2 * time.Second)
	for {
		select {
		case data := <-s.logChan:
			call.CatchPanic(func() {
				s.tlogBaseService.WriteTlog(data)
			}, "s.tlogBaseService.WriteTlog(data)")
		case <-tick.C:
			call.CatchPanic(s.handle, "s.tlogBaseService.tick")
		case <-s.closeChan:
			return
		}
	}

}

func (s *TlogService) WriteLog(data string) {
	if len(data) <= 0 {
		return
	}

	select {
	case s.logChan <- data:
	default:
		if !s.addCacheLog(data) {
			// 添加失败，打印日志
			logrus.Warn("TLogCacheFull " + data)
		}
	}
}

func (s *TlogService) handle() {
	if datas := s.popCacheLog(); len(datas) > 0 {
		for _, data := range datas {
			s.tlogBaseService.WriteTlog(data)
		}
	}
}

func (s *TlogService) popCacheLog() []string {
	s.logCacheMux.Lock()
	defer s.logCacheMux.Unlock()

	if len(s.logCache) > 0 {
		cache := s.logCache
		s.logCache = nil
		return cache
	}
	return nil
}

func (s *TlogService) addCacheLog(data string) bool {
	s.logCacheMux.Lock()
	defer s.logCacheMux.Unlock()

	n := len(s.logCache)
	if n >= cap(s.logCache) {
		if n >= s.maxLogSize {
			// 超出限制，打印日志就结束
			return false
		}

		newSize := imath.Max(imath.Min(n*2, s.maxLogSize), 32)
		newLogCache := make([]string, newSize)
		copy(newLogCache, s.logCache)
		s.logCache = newLogCache
	}

	s.logCache = append(s.logCache, data)
	return true
}

// 工具类数据

func (s *TlogService) buildLogHeroTx(heroInfo entity.TlogHero, name string,
	f func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string,
) (str string) {
	if s.dontGenTlog {
		return
	}

	if heroInfo == nil {
		logrus.Errorf("TlogService.%s heroInfo == nil", name)
		return
	}

	tencentInfo := s.worldService.GetTencentInfo(heroInfo.Id())
	if tencentInfo == nil {
		//logrus.Warnf("TlogService.%s tencentInfo == nil", name)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("err", err).Errorf("recovered from tlogService.%s panic. SEVERE!!!", name)
			metrics.IncPanic()
		}
	}()

	return f(heroInfo, tencentInfo)
}

func (s *TlogService) buildLogHero(heroInfo entity.TlogHero, name string,
	f func(heroInfo entity.TlogHero) string,
) (str string) {
	if s.dontGenTlog {
		return
	}

	if heroInfo == nil {
		logrus.Errorf("TlogService.%s heroInfo == nil", name)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("err", err).Errorf("recovered from tlogService.%s panic. SEVERE!!!", name)
			metrics.IncPanic()
		}
	}()

	return f(heroInfo)
}

func (s *TlogService) buildLogHeroIdTx(heroId int64, name string,
	f func(heroId int64, tencentInfo *shared_proto.TencentInfoProto) string,
) (str string) {
	if s.dontGenTlog {
		return
	}

	tencentInfo := s.worldService.GetTencentInfo(heroId)
	if tencentInfo == nil {
		logrus.Warnf("TlogService.%s tencentInfo == nil", name)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("err", err).Errorf("recovered from tlogService.%s panic. SEVERE!!!", name)
			metrics.IncPanic()
		}
	}()

	return f(heroId, tencentInfo)
}

func (s *TlogService) buildLogHeroId(heroId int64, name string,
	f func(heroId int64) string,
) (str string) {
	if s.dontGenTlog {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("err", err).Errorf("recovered from tlogService.%s panic. SEVERE!!!", name)
			metrics.IncPanic()
		}
	}()

	return f(heroId)
}

func (s *TlogService) buildLog(name string,
	f func() string,
) (str string) {
	if s.dontGenTlog {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logrus.WithField("err", err).Errorf("recovered from tlogService.%s panic. SEVERE!!!", name)
			metrics.IncPanic()
		}
	}()

	return f()
}

func writeU64(buf *bytes.Buffer, data uint64) {
	buf.WriteString(strconv.FormatUint(data, 10))
}

func writeI32(buf *bytes.Buffer, data int32) {
	writeI64(buf, int64(data))
}

func writeI64(buf *bytes.Buffer, data int64) {
	buf.WriteString(strconv.FormatInt(data, 10))
}

func writeInt(buf *bytes.Buffer, data int) {
	buf.WriteString(strconv.Itoa(data))
}

func writeF32(buf *bytes.Buffer, data float32) {
	writeF64(buf, float64(data))
}

func writeF64(buf *bytes.Buffer, data float64) {
	buf.WriteString(strconv.FormatFloat(data, 'f', -1, 64))
}

func writeTime(buf *bytes.Buffer, data time.Time) {
	buf.WriteString(data.Format("2006-01-02 15:04:05"))
}

func writeString(buf *bytes.Buffer, data string) {
	buf.WriteString(data)
}

func writeBool(buf *bytes.Buffer, data bool) {

	buf.WriteString(strconv.FormatBool(data))
}

func writeU64Array(buf *bytes.Buffer, array []uint64) {
	var strs []string
	for _, i := range array {
		strs = append(strs, strconv.FormatUint(i, 10))
	}
	buf.WriteString(strings.Join(strs, ","))
}
