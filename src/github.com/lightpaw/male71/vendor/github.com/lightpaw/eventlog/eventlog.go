package eventlog

import (
	"github.com/lightpaw/logrus"
	"math/rand"
	"regexp"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultDelay = 1 * time.Minute

	seconds_per_day = 60 * 60 * 24
	uuidCounterMask = 1<<25 - 1
)

var (
	head *EventLog
	tail *EventLog
	lock sync.Mutex

	// for uuid gen. 17 bit time (today second), 5 bit pid, 10 bit sid, 8 bit rand, 24 bit counter
	uuidCounter int64 = 1
	ctime             = time.Now().UnixNano()
	randCounter       = rand.Int63() & 255

	optionLock  sync.Mutex
	destination destinationFace
	maxDelay    time.Duration

	typeNameMatcher = regexp.MustCompile("^[a-zA-Z]\\w*$")
)

// 只能调用一次
func Start(dest destinationFace, mDelay time.Duration) {
	optionLock.Lock()
	defer optionLock.Unlock()
	if destination != nil || maxDelay != 0 {
		logrus.Panic("eventlog.Start不能重复调用")
	}

	if dest == nil {
		logrus.Panic("eventlog.Start必须指定destination")
	}

	destination = dest
	maxDelay = mDelay
	if maxDelay <= 0 {
		maxDelay = defaultDelay
	}

	go loop()
}

// 不等定时器, 直接将缓存着的消息发送出去.
// 关服时调用
func Flush() {
	defer func() {
		if r := recover(); r != nil {
			// 严重错误. 这里不能panic
			logrus.WithField("err", r).Error("eventlog.Flush recovered from panic!!! SERIOUS PROBLEM")
			debug.PrintStack()
		}
	}()

	if destination == nil {
		logrus.Error("eventlog还没有Start, 不能Flush")
		return
	}

	lock.Lock()
	toFlushHead := head
	head, tail = nil, nil
	lock.Unlock()

	var err error
	for {
		if toFlushHead == nil {
			return
		}
		logrus.Debug("Flush日志")
		toFlushHead, err = destination.Write(toFlushHead)
		if err != nil {
			logrus.WithError(err).Error("eventlog发送日志出错")
			time.Sleep(1 * time.Second)
		} else {
			logrus.Debug("日志Flush成功")
			return
		}
	}
}

func loop() {
	// 定时flush

	for range time.Tick(maxDelay) {
		Flush()
	}
}

// 记录日志
func Commit(event *EventLog) {
	lock.Lock()

	if head == nil {
		head = event
		tail = event
	} else {
		tail.next = event
		tail = event
	}

	lock.Unlock()
}

// 新日志, 并不发送. 发送必须调用Commit
func NewEvent(eventType string, platformID, serverID uint32) *EventLog {
	result := &EventLog{
		fields: make(map[string]interface{}, 16),
	}

	if !typeNameMatcher.MatchString(eventType) {
		logrus.WithField("type", eventType).Panic("eventlog的type类型只能是字母或数字和下划线, 且不能数字和下划线开头")
	}

	ctime := time.Now().Unix()

	result.fields["_type"] = eventType
	result.fields["_pid"] = platformID
	result.fields["_sid"] = serverID
	result.fields["_time"] = ctime
	result.fields["_uuid"] = (ctime % seconds_per_day << 47) | (int64(platformID) << 42) | (int64(serverID) << 32) | (randCounter << 24) | (atomic.AddInt64(&uuidCounter, 1) & uuidCounterMask)

	return result
}

type EventLog struct {
	next *EventLog

	fields map[string]interface{}
}

type data map[string]interface{}

// 提供修改时间的方法
func (l *EventLog) WithTime(t time.Time) *EventLog {
	l.fields["_time"] = t.Unix()
	return l
}

func (l *EventLog) With(field string, value interface{}) *EventLog {
	if !typeNameMatcher.MatchString(field) {
		logrus.WithField("field", field).Panic("eventlog的field类型只能是字母或数字和下划线, 且不能数字和下划线开头")
	}

	if field == "" {
		logrus.Panic("EventLog的字段名不能为空")
	}

	switch value.(type) {
	case int8:
	case uint8:
	case int16:
	case uint16:
	case int32:
	case uint32:
	case int64:
	case uint64:
	case int:
	case uint:
	case string:
	case float32:
	case float64:
	case bool:
	default:
		logrus.WithField("value", value).Panic("eventlog的value类型只能是数字或bool或string类型")
	}

	if _, has := l.fields[field]; has {
		logrus.WithField("field", field).Panic("EventLog的字段名重复")
	}

	l.fields[field] = value

	return l
}
