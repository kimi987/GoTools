package eventlog

import (
	"bytes"
	"compress/gzip"
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"github.com/tinylib/msgp/msgp"
	"os"
	"path"
	"sync"
	"time"
)

type msgpFileDestination struct {
	dir string

	round time.Duration

	bufPool sync.Pool
}

func NewMsgpFileDestination(dir string, round time.Duration) *msgpFileDestination {
	pool := sync.Pool{New: func() interface{} {
		buf := bytes.NewBuffer(make([]byte, 0, 65536))
		return &msgpWriter{buf: buf, msgpWriter: msgp.NewWriter(buf)}
	}}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logrus.WithError(err).Panic("无法创建log目录")
	}

	if round < time.Minute || round > 24*time.Hour {
		logrus.Panic("目录分隔的时间，必须")
	}
	return &msgpFileDestination{dir: dir, round: round, bufPool: pool}
}

func (d *msgpFileDestination) Write(log *EventLog) (*EventLog, error) {
	writer := d.bufPool.Get().(*msgpWriter)
	buf := writer.buf
	defer func() {
		writer.buf.Reset()
		writer.msgpWriter.Reset(writer.buf)
		d.bufPool.Put(writer)
	}()

	for event := log; event != nil; event = event.next {
		if err := data(event.fields).EncodeMsg(writer.msgpWriter); err != nil {
			return nil, errors.Wrap(err, "msgpack Encode log 出错")
		}

		//if writer.buf.Len() > single_record_limit_minus_buf {
		//	if err := writer.msgpWriter.Flush(); err != nil {
		//		return nil, errors.Wrap(err, "msgpack Flush log 出错")
		//	}
		//
		//	if err := f.doSend(writer.buf.Bytes()); err != nil {
		//		return nil, errors.Wrap(err, "发送log出错")
		//	}
		//	writer.buf.Reset()
		//}
	}

	if err := writer.msgpWriter.Flush(); err != nil {
		return nil, errors.Wrap(err, "msgpack Flush log 出错")
	}

	t := time.Now()
	fileName := path.Join(d.dir, t.Format("2006/01/02"), t.Round(d.round).Format("2006-01-02-15-04")+".gz")
	if err := os.MkdirAll(path.Dir(fileName), os.ModePerm); err != nil {
		logrus.WithError(err).WithField("file", fileName).Panic("无法创建log目录")
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		logrus.WithError(err).Error("打开log文件失败")
		return log, nil
	}
	defer file.Close()

	gfile := gzip.NewWriter(file)
	defer gfile.Close()

	if _, err := gfile.Write(buf.Bytes()); err != nil {
		logrus.WithError(err).Error("写入log文件失败")
		return log, nil
	}

	return nil, nil
}

type msgpWriter struct {
	buf        *bytes.Buffer
	msgpWriter *msgp.Writer
}
