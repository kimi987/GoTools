package kafka

import (
	"testing"
	"github.com/lightpaw/logrus"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/lightpaw/male7/util/call"
	"sync"
	"time"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/util/atomic"
)


func TestSend(t *testing.T) {
	//send_test()
}

func send_test() {
	s := KafkaService{}
	conf := conf()
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, conf)
	if err != nil {
		logrus.WithError(err).Panicf("")
	}
	s.syncProducer = producer
	defer s.syncProducer.Close()

	asyncProducer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, conf)
	if err != nil || !conf.Producer.Return.Successes {
		logrus.WithError(err).Panicf("sarama.NewAsyncProducer 失败")
	}
	s.asyncProducer = asyncProducer
	defer s.asyncProducer.Close()
	go call.CatchLoopPanic(s.asyncProducerResultLoop, "KafkaService.asyncProducerResultLoop")

	topic1 := entity.Topic("test1")
	topic2 := entity.Topic("test2")
	topic3 := entity.Topic("test3")
	topics := []entity.Topic{topic1, topic2, topic3}

	value := atomic.NewInt64(time.Now().Unix())

	wg1 := sync.WaitGroup{}
	for i := 0; i < 12; i++ {
		wg1.Add(1)
		idx := i % len(topics)
		topic := topics[idx]
		go func() {
			value := fmt.Sprintf("topic %v sync value:%v", topic, value.Inc())
			msg := s.NewProducerMsg(topic, []byte(value))
			s.SendSync(msg)
			fmt.Printf("sync send %+v %+v\n", msg.Topic, msg.Value)
			wg1.Done()
		}()
	}
	wg1.Wait()
	fmt.Println("sync send producer closed")

	wg2 := sync.WaitGroup{}
	for i := 0; i < 12; i++ {
		wg2.Add(1)
		idx := i % len(topics)
		go func() {
			topic := topics[idx]
			value := fmt.Sprintf("topic %v async value:%v", topic, value.Inc())
			msg := s.NewProducerMsg(topic, []byte(value))
			s.SendAsync(msg)
			wg2.Done()
		}()
	}
	wg2.Wait()
	fmt.Println("async send producer closed")

	// 等回调
	d, _ := time.ParseDuration("1s")
	time.Sleep(d)
}
