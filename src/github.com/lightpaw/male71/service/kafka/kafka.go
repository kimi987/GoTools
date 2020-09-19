package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/pb/mlogpb/mlog"
	"github.com/lightpaw/male7/tlog"
	"fmt"
	"github.com/lightpaw/male7/config/kv"
	"time"
	"path"
)

func conf() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Compression = sarama.CompressionGZIP
	config.Producer.CompressionLevel = 2
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	return config
}

//gogen:iface
type KafkaService struct {
	dbService   iface.DbService
	timeService iface.TimeService
	config      *kv.IndividualServerConfig

	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer

	tlogFile *tlog.TlogFile

	closeNotify    chan struct{}
	loopExitNotify chan struct{}
}

func NewKafkaService(config *kv.IndividualServerConfig, dbService iface.DbService, timeService iface.TimeService) *KafkaService {
	m := &KafkaService{}
	m.dbService = dbService
	m.timeService = timeService
	m.config = config

	if !config.GetKafkaStart() {
		logrus.Info("kafka service not run")
		return m
	}
	pidStr := fmt.Sprintf("%d", config.GetPlatformID())
	sidStr := fmt.Sprintf("%d", config.GetServerID())
	fileCurrentDir := path.Join(config.TlogKafkaFailCurrentBaseDir, pidStr, sidStr)
	fileArchiveDir := path.Join(config.TlogKafkaFailArchiveBaseDir, pidStr, sidStr)

	filenamePrefix := "kafka_"

	m.tlogFile = tlog.NewTlogFile(timeService, fileCurrentDir, fileArchiveDir, filenamePrefix,
		config.TlogKafkaFailRotateDuration, config.TlogKafkaFailRotateSize, config.TlogKafkaFailWriteBufSize, config.TlogKafkaFailCacheSize)

	conf := conf()
	syncProducer, err := sarama.NewSyncProducer(config.GetKafkaBrokerAddr(), conf)
	if err != nil {
		logrus.WithError(err).Panicf("sarama.NewSyncProducer 失败")
	}
	m.syncProducer = syncProducer

	m.closeNotify = make(chan struct{})
	m.loopExitNotify = make(chan struct{})

	asyncProducer, err := sarama.NewAsyncProducer(config.GetKafkaBrokerAddr(), conf)
	if err != nil || !conf.Producer.Return.Successes {
		logrus.WithError(err).Panicf("sarama.NewAsyncProducer 失败")
	}
	m.asyncProducer = asyncProducer

	go call.CatchLoopPanic(m.asyncProducerResultLoop, "KafkaService.asyncProducerResultLoop")

	logrus.Info("kafka service is running")
	return m
}

func (m *KafkaService) Close() {
	if !m.config.GetKafkaStart() {
		return
	}

	close(m.closeNotify)
	<-m.loopExitNotify

	if m.tlogFile != nil {
		m.tlogFile.Close()
	}

	if m.syncProducer != nil {
		m.syncProducer.Close()
	}

	if m.asyncProducer != nil {
		m.asyncProducer.Close()
	}

}

func (m *KafkaService) SendSync(msg *sarama.ProducerMessage) error {
	partition, offset, err := m.syncProducer.SendMessage(msg)
	if err != nil {
		logrus.WithError(err).Errorf("kafka topic:%v partition:%v offset;%v Send message Fail. msg:%v", msg.Topic, partition, offset, msg)
		m.doIfError(msg)
		return err
	}

	logrus.Debugf("kafka Sync topic:%v Partition = %d, offset=%d\n", msg.Topic, partition, offset)

	return nil
}

func (m *KafkaService) SendAsync(msg *sarama.ProducerMessage) {

	select {
	case m.asyncProducer.Input() <- msg:
		return
	case <-time.After(100 * time.Millisecond):
		logrus.Error("异步发送kafka消息超时")
		return
	}
}

func (m *KafkaService) asyncProducerResultLoop() {
	defer close(m.loopExitNotify)

	for {
		select {
		case msg := <-m.asyncProducer.Successes():
			logrus.Debugf("kafka Async response succ. topic:%v Partition = %d, offset=%d\n", msg.Topic, msg.Partition, msg.Offset)
		case err := <-m.asyncProducer.Errors():
			logrus.WithError(err).Errorf("kafka Async Send message Fail.")
			m.doIfError(err.Msg)
		case <-m.closeNotify:
			return
		}
	}
}

func (m *KafkaService) doIfError(msg *sarama.ProducerMessage) {
	bytes, err := msg.Value.Encode()
	if err != nil {
		logrus.WithError(err).Errorf(" kafka msg.Value.Encode 失败")
		return
	}

	proto := &mlog.MlogProto{}
	err = proto.Unmarshal(bytes)
	if err != nil {
		logrus.WithError(err).Errorf(" kafka proto.Unmarshal 失败")
		return
	}

	m.tlogFile.AddLog(string(proto.Content))
}

func (m *KafkaService) NewProducerMsg(topic entity.Topic, value []byte) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic: string(topic),
		Key:   sarama.StringEncoder("key"),
		Value: sarama.ByteEncoder(value),
	}

	return msg
}
