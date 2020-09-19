package tlog

import (
	"github.com/lightpaw/logrus"
	"time"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/pb/mlogpb/mlog"
	"github.com/lightpaw/male7/util/atomic"
	"fmt"
	"github.com/lightpaw/male7/config/kv"
	"path"
)

const (
	sep  = "|"
	line = "\n"
)

//gogen:iface
type TlogBaseService struct {
	timeService iface.TimeService
	config      *kv.IndividualServerConfig
	kafka       iface.KafkaService

	closed *atomic.Bool // 服务已关闭

	nextNewFileTime time.Time
	newFileDuration time.Duration

	tlogFile *TlogFile

	closeNotify    chan struct{}
	loopExitNotify chan struct{}
}

func NewTlogBaseService(timeService iface.TimeService, kafka iface.KafkaService, config *kv.IndividualServerConfig) *TlogBaseService {
	s := &TlogBaseService{
		closeNotify:    make(chan struct{}),
		loopExitNotify: make(chan struct{}),
	}

	s.config = config
	s.timeService = timeService
	s.kafka = kafka
	s.closed = atomic.NewBool(false)

	if config.GetTlogStart() {
		pidStr := fmt.Sprintf("%d", config.GetPlatformID())
		sidStr := fmt.Sprintf("%d", config.GetServerID())
		fileCurrentDir := path.Join(config.TlogCurrentBaseDir, pidStr, sidStr)
		fileArchiveDir := path.Join(config.TlogArchiveBaseDir, pidStr, sidStr)

		filenamePrefix := "tlog_"

		s.tlogFile = NewTlogFile(timeService, fileCurrentDir, fileArchiveDir, filenamePrefix,
			config.TlogRotateDuration, config.TlogRotateSize, config.TlogWriteBufSize, config.TlogCacheSize)
	}

	return s
}

func (s *TlogBaseService) Close() {
	s.closed.Store(true)

	//  tlogFile 由 loop 切到主线程 flush
	if s.tlogFile != nil {
		s.tlogFile.Close()
	}

	logrus.Debugf("TlogBaseService 关闭成功")
}

/*
	========= tlog file =============
 */

func (s *TlogBaseService) WriteTlog(content string) bool {
	if s.closed.Load() {
		return true
	}

	if len(content) <= 0 {
		logrus.Warnf("TlogBaseService len(content) == 0")
		return true
	}

	if s.config.TlogStart {
		s.tlogFile.AddLog(content)
	}

	if s.config.GetKafkaStart() {
		v, err := s.buildMlogProto(content)
		if err != nil {
			logrus.WithError(err).Errorf("MlogProto.Marshal() 异常")
			return false
		}

		msg := s.kafka.NewProducerMsg(entity.Topic(s.config.GetTlogTopic()), v)
		s.kafka.SendAsync(msg)
		//if err := s.kafka.SendSync(msg); err != nil {
		//	return false
		//}

		// 如果写文件队列满了，这里会阻塞到超时
	}

	return true
}

func (s *TlogBaseService) buildMlogProto(content string) (v []byte, err error) {
	p := &mlog.MlogProto{}
	p.Pid = int64(s.config.GetPlatformID())
	p.Sid = int64(s.config.GetServerID())
	p.Content = []byte(content)
	return p.Marshal()
}
