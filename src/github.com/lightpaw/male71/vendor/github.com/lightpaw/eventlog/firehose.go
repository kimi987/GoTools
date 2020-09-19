package eventlog

import (
	"bytes"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/firehose/firehoseiface"
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"github.com/tinylib/msgp/msgp"
	"sync"
	"time"
)

const (
	single_record_limit           = 1000000
	single_record_limit_minus_buf = single_record_limit - 3000
)

type firehoseDestination struct {
	client firehoseiface.FirehoseAPI

	deliveryStreamName string

	writerPool sync.Pool
}

func NewFirehoseDestination(deliveryStreamName string, client firehoseiface.FirehoseAPI) *firehoseDestination {
	pool := sync.Pool{New: func() interface{} {
		buf := bytes.NewBuffer(make([]byte, 0, 65536))
		return &firehoseWriter{buf: buf, msgpWriter: msgp.NewWriter(buf)}
	}}
	return &firehoseDestination{client: client, deliveryStreamName: deliveryStreamName, writerPool: pool}
}

func (f *firehoseDestination) Write(log *EventLog) (newHead *EventLog, err error) {
	if log == nil {
		return
	}

	newHead = log

	writer := f.writerPool.Get().(*firehoseWriter)
	defer func() {
		writer.buf.Reset()
		writer.msgpWriter.Reset(writer.buf)
		f.writerPool.Put(writer)
	}()

	for event := log; event != nil; event = event.next {
		if err = data(event.fields).EncodeMsg(writer.msgpWriter); err != nil {
			err = errors.Wrap(err, "msgpack Encode log 出错")
			return
		}

		if writer.buf.Len() > single_record_limit_minus_buf {
			if err = writer.msgpWriter.Flush(); err != nil {
				err = errors.Wrap(err, "msgpack Flush log 出错")
				return
			}

			if err = f.doSend(writer.buf.Bytes()); err != nil {
				err = errors.Wrap(err, "发送log出错")
				return
			}
			writer.buf.Reset()
			newHead = event.next
		}
	}

	if err = writer.msgpWriter.Flush(); err != nil {
		err = errors.Wrap(err, "msgpack Flush log 出错")
		return
	}

	if writer.buf.Len() > 0 {
		if err = f.doSend(writer.buf.Bytes()); err != nil {
			err = errors.Wrap(err, "发送log出错")
		} else {
			newHead = nil
		}
	}

	return
}

type firehoseWriter struct {
	buf        *bytes.Buffer
	msgpWriter *msgp.Writer
}

func (f *firehoseDestination) doSend(data []byte) (err error) {
	sleepTime := 1 * time.Second
	for i := 0; i < 10; i++ {
		input := &firehose.PutRecordInput{DeliveryStreamName: &f.deliveryStreamName, Record: &firehose.Record{Data: data}}

		_, err = f.client.PutRecord(input)
		if err == nil {
			return nil
		}

		logrus.WithError(err).Error("Log发送失败")
		time.Sleep(sleepTime)
		sleepTime += 2 * time.Second
	}
	return err
}
