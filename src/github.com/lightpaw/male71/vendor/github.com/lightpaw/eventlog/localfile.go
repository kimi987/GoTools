package eventlog

import (
	"bytes"
	"fmt"
	"github.com/lightpaw/logrus"
	"os"
	"path"
	"sync"
	"time"
)

type localFileDestination struct {
	file string

	bufPool sync.Pool
}

func NewLocalFileDestination(file string) *localFileDestination {
	pool := sync.Pool{New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 65536))
	}}

	if dir := path.Dir(file); dir != "" {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logrus.WithError(err).Panic("无法创建log目录")
		}
	}

	if file, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm); err != nil {
		logrus.WithError(err).Panic("无法打开log文件")
	} else {
		file.Close()
	}

	return &localFileDestination{file: file, bufPool: pool}
}

func (d *localFileDestination) Write(log *EventLog) (*EventLog, error) {
	buf := d.bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		d.bufPool.Put(buf)
	}()

	for event := log; event != nil; event = event.next {
		eventType := event.fields["_type"]
		eventTime := event.fields["_time"]
		t := time.Unix(eventTime.(int64), 0)
		timeFormat := t.Format("2006-01-02 15:04:05")

		fmt.Fprintf(buf, "%10s %s ", eventType, timeFormat)

		for k, v := range event.fields {
			switch k {
			case "_time":
				fallthrough
			case "_type":
				fallthrough
			case "_uuid":
				fallthrough
			case "_pid":
				fallthrough
			case "_sid":
				continue
			default:
				fmt.Fprintf(buf, "%s: %v ", k, v)
			}
		}

		buf.WriteString("\n\n")
	}

	file, err := os.OpenFile(d.file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		logrus.WithError(err).Error("打开log文件失败")
		return log, nil
	}

	defer file.Close()
	if _, err := file.Write(buf.Bytes()); err != nil {
		logrus.WithError(err).Error("写入log文件失败")
		return log, nil
	}

	return nil, nil
}
