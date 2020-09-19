package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/lightpaw/eventlog"
	"github.com/lightpaw/logrus"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"
)

var (
	logPath = flag.String("path", "testlog", "log file path")

	startTimeStr = flag.String("start_time", "2017-10-01", "日志开始时间")

	logDay = flag.Int("log_day", 15, "日志天数")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// 数据生成到今天为止
	os.RemoveAll(*logPath)

	startTime, err := time.ParseInLocation("2006-01-02", *startTimeStr, time.UTC)
	if err != nil {
		logrus.WithError(err).Panicf("解析start_server失败")
	}

	if *logDay > 0 {
		startTime = time.Now().Add(time.Duration(-*logDay) * 24 * time.Hour).UTC().Truncate(24 * time.Hour)
	}

	eventlog.Start(eventlog.NewMsgpFileDestination("testlog/event/", time.Hour), 10*time.Second)

	serverInfos := []*server_info{
		newServer(1, 1, 0,
			200, 10, 0,
			[]int{0, 14, 10, 7, 6, 5, 4, 3, 2, 1},                         // aliveDays
			[]float64{0.2, 0.3, 0.35, 0.4, 0.45, 0.5, 0.55, 0.6, 0.65, 1}, // aliveRates
			[][]string{{"ios", "andriod"}, {"andriod"}, {"ios"}},          // osTypes
			[]float64{0.1, 0.4, 1},                                        // osTypeRates
			[]int{5000, 4000, 3000, 2000, 1000, 0},                        // progress
			[]float64{0.7, 0.75, 0.8, 0.85, 0.9, 1},                       // progressRates
			[]int{1, 2, 3},                                                // firstRechargeDays
			[]float64{0.2, 0.3, 0.35},                                     // firstRechargeDayRates
		),
		newServer(1, 2, 2,
			300, 30, 0,
			[]int{0, 14, 10, 7, 6, 5, 4, 3, 2, 1},                          // aliveDays
			[]float64{0.25, 0.3, 0.35, 0.4, 0.45, 0.5, 0.55, 0.6, 0.65, 1}, // aliveRates
			[][]string{{"ios", "andriod"}, {"andriod"}, {"ios"}},           // osTypes
			[]float64{0.1, 0.4, 1},                                         // osTypeRates
			[]int{5000, 4000, 3000, 2000, 1000, 0},                         // progress
			[]float64{0.6, 0.75, 0.8, 0.85, 0.9, 1},                        // progressRates
			[]int{1, 2, 3},                                                 // firstRechargeDays
			[]float64{0.25, 0.3, 0.35},                                     // firstRechargeDayRates
		),
		newServer(2, 1, 1,
			250, 20, 0,
			[]int{0, 14, 10, 7, 6, 5, 4, 3, 2, 1},                        // aliveDays
			[]float64{0.1, 0.2, 0.3, 0.35, 0.4, 0.45, 0.5, 0.55, 0.6, 1}, // aliveRates
			[][]string{{"ios", "andriod"}, {"andriod"}, {"ios"}},         // osTypes
			[]float64{0.1, 0.3, 1},                                       // osTypeRates
			[]int{5000, 4000, 3000, 2000, 1000, 0},                       // progress
			[]float64{0.7, 0.75, 0.8, 0.85, 0.9, 1},                      // progressRates
			[]int{1, 2, 3},                                               // firstRechargeDays
			[]float64{0.1, 0.2, 0.3},                                     // firstRechargeDayRates
		),
	}

	infos := []log_type{
		&create_hero_log{},
		&online_log{},
		&new_guide_log{},
		&device_login_log{},
		&recharge_log{},
		&yuanbao_log{},
	}

	formats := []*format{
		newCsvFormat(),
		newMysqlFormat(),
		newEventFormat(),
	}

	now := time.Now().UTC()
	for _, s := range serverInfos {
		if err := s.generateFormatLogData(*logPath, startTime, now, infos, formats); err != nil {
			logrus.WithError(err).Panicf("生成数据失败")
		}
	}

	eventlog.Flush()
}

func newCsvFormat() *format {
	f := &format{}
	f.name = "csv"
	f.writeDailyFile = f.writeDtTimeFile

	f.getRecordeWriter = func(t log_type) RecordWriter {
		return t.getCsvRecordWriter()
	}

	return f
}

func newMysqlFormat() *format {
	f := &format{}
	f.name = "sql"
	f.writeAllFile = f.writeSingleFile

	f.getRecordeWriter = func(t log_type) RecordWriter {
		return t.getMysqlRecordWriter()
	}

	return f
}

func newEventFormat() *format {
	f := &format{}
	f.name = "event"

	f.getRecordeWriter = func(t log_type) RecordWriter {
		return t.getEventRecordWriter()
	}

	return f
}

type format struct {
	name string

	writeDailyFile func(basePath string, t time.Time, b *bytes.Buffer, pid, sid int) error
	writeAllFile   func(basePath string, b *bytes.Buffer, pid, sid int) error

	getRecordeWriter func(t log_type) RecordWriter
}

type log_type interface {
	getName() string
	getCsvRecordWriter() RecordWriter
	getMysqlRecordWriter() RecordWriter
	getEventRecordWriter() RecordWriter
}

type RecordWriter func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero)

func (s *server_info) generateFormatLogData(basePath string, startTime, now time.Time,
	logs []log_type, formats []*format) error {

	s.initHero(startTime, now)

	b := &bytes.Buffer{}

	for _, f := range formats {
		for _, v := range logs {
			writeRecord := f.getRecordeWriter(v)
			if writeRecord == nil {
				return errors.Errorf("没找到日志[%s]对应的格式生成器, %s", v.getName(), f.name)
			}

			err := s.generateLogs(path.Join(basePath, f.name, v.getName()), startTime, now, b, writeRecord, f.writeDailyFile)
			if err != nil {
				return errors.Wrapf(err, "日志[%s]，格式[%s]，生成错误", v.getName(), f.name)
			}
		}

		if f.writeAllFile != nil {
			err := f.writeAllFile(path.Join(basePath, f.name), b, s.pid, s.sid)
			if err != nil {
				return errors.Wrapf(err, "All日志，格式[%s]，生成错误", f.name)
			}
		}
	}

	s.heroMap = nil

	return nil
}

func (s *server_info) generateLogs(basePath string, startTime, now time.Time, b *bytes.Buffer,
	writeRecords func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero),
	writeDailyFile func(basePath string, t time.Time, b *bytes.Buffer, pid, sid int) error,
) error {
	startServerTime := startTime.Add(time.Duration(s.startDelayDay) * Day)
	diff := now.Sub(startServerTime)

	if diff < 0 {
		return nil
	}

	d := int((diff + Day - 1) / Day)
	for i := 0; i < d; i++ {
		t := startServerTime.Add(time.Duration(i) * Day)

		writeRecords(s, b, startServerTime, t, s.heroMap)

		if writeDailyFile != nil {
			// 写入文件
			if err := writeDailyFile(basePath, t, b, s.pid, s.sid); err != nil {
				return err
			}
		}
	}

	return nil
}

func newServer(pid int, sid int, startDelayDay int,
	firstDayUser int, otherDayUserMin int, otherDayUserRandom int,
	aliveDays []int,
	aliveRates []float64,
	osTypes [][]string,
	osTypeRates []float64,
	progress []int,
	progressRates []float64,
	firstRechargeDays []int,
	firstRechargeDayRates []float64,
) *server_info {
	return &server_info{
		pid:                   pid,
		sid:                   sid,
		startDelayDay:         startDelayDay,
		firstDayUser:          firstDayUser,
		otherDayUserMin:       otherDayUserMin,
		otherDayUserRandom:    otherDayUserRandom,
		aliveDays:             aliveDays,
		aliveRates:            aliveRates,
		osTypes:               osTypes,
		osTypeRates:           osTypeRates,
		progress:              progress,
		progressRates:         progressRates,
		firstRechargeDays:     firstRechargeDays,
		firstRechargeDayRates: firstRechargeDayRates,
		heroMap:               make(map[int64]*hero),
	}
}

type server_info struct {
	pid int // 平台id
	sid int // 区服id

	startDelayDay int // 开服延迟几天

	firstDayUser int // 第一天导量

	otherDayUserMin    int // 其他时间新进最少（总量=最少+随机）
	otherDayUserRandom int // 其他时间新进随机

	// 留存数据
	aliveDays  []int
	aliveRates []float64

	// 设备占比
	osTypes     [][]string
	osTypeRates []float64

	// 新手引导占比
	progress      []int
	progressRates []float64

	// 首冲天数占比
	firstRechargeDays     []int
	firstRechargeDayRates []float64

	// 新用户id
	idCounter int64

	heroMap map[int64]*hero
}

func (s *server_info) firstRechargeDay(rate float64) (day, rechargeAmount, level int, rechargeType string) {

	for i, v := range s.firstRechargeDayRates {
		if rate <= v {
			return s.firstRechargeDays[i], (2*i + 1) * 10, 2*i + 1, fmt.Sprintf("充值-%d", i)
		}
	}

	return 0, 0, 0, ""
}

func (s *server_info) aliveDay(rate float64) int {

	for i, v := range s.aliveRates {
		if rate <= v {
			return s.aliveDays[i]
		}
	}

	return s.aliveDays[len(s.aliveDays)-1]
}

func (s *server_info) osType(rate float64) []string {

	for i, v := range s.osTypeRates {
		if rate <= v {
			return s.osTypes[i]
		}
	}

	return s.osTypes[len(s.osTypes)-1]
}

func (s *server_info) guideProgress(rate float64) (int, bool) {

	for i, v := range s.progressRates {
		if rate <= v {
			return s.progress[i], i <= 0
		}
	}

	return s.progress[len(s.progressRates)-1], false
}

func (s *server_info) initHero(startTime, now time.Time) {

	startServerTime := startTime.Add(time.Duration(s.startDelayDay) * Day)
	diff := now.Sub(startServerTime)

	if diff < 0 {
		return
	}

	d := int((diff + Day - 1) / Day)
	for i := 0; i < d; i++ {
		t := startServerTime.Add(time.Duration(i) * Day)

		newUserCount := s.firstDayUser // 首日
		if i > 0 {
			// 其他时间
			newUserCount = s.otherDayUserMin
			if s.otherDayUserRandom > 0 {
				newUserCount += rand.Intn(s.otherDayUserRandom)
			}
		}

		// 进入新建英雄
		for i := 0; i < newUserCount; i++ {
			id := s.newHeroId()
			createTime := t.Add(time.Duration(i%SecondsPerDay) * time.Second)

			// 登陆X日
			rate := float64(i) / float64(newUserCount)

			// 设备
			osTypes := s.osType(rate)
			var deviceIds []string
			for i := 0; i < len(osTypes); i++ {
				deviceIds = append(deviceIds, uuid.New().String())
			}

			// 新手引导进度
			progress, completed := s.guideProgress(rate)

			// 首冲
			firstRechargeDay, firstRechargeAmount, firstRechargeLevel, rechargeType := s.firstRechargeDay(rate)

			s.heroMap[id] = &hero{
				id:               id,
				name:             fmt.Sprintf("name-%d", id),
				createTime:       createTime.Unix(),
				onlineTotalDay:   s.aliveDay(rate),
				onlineOffsetDay:  1,
				progress:         progress,
				completed:        completed,
				deviceIds:        deviceIds,
				osTypes:          osTypes,
				firstRechargeDay: firstRechargeDay,
				rechargeAmount:   firstRechargeAmount,
				rechargeLevel:    firstRechargeLevel,
				rechargeType:     rechargeType,
			}

		}

	}
}

type hero struct {
	id         int64
	name       string
	createTime int64

	onlineTotalDay  int // 在线几天，包含登陆当天，0表示不流失
	onlineOffsetDay int // 下次登陆间隔

	progress  int  // 新手完成度
	completed bool // 是否完成新手

	// 登陆设备号
	deviceIds []string
	osTypes   []string

	// 充值（第几天，充值金额，充值等级）
	firstRechargeDay int
	rechargeAmount   int
	rechargeLevel    int
	rechargeType     string

	// 消费，默认充值当天就消费一半
}

func (h *hero) getRemainYuanbao(createDay int) int64 {
	if createDay >= h.firstRechargeDay {
		return int64(h.rechargeAmount) / 2
	}

	return 0
}

func (s *server_info) newHeroId() int64 {
	s.idCounter++
	return s.idCounter
}

const Day = 24 * time.Hour
const SecondsPerDay = int(Day / time.Second)
const Int64SecondsPerDay = int64(SecondsPerDay)

func (l *create_hero_log) writeData(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, v *hero,
	recordWriter DataWriter_create_hero_log) bool {
	startServiceUnixTime := startServerTime.Unix()

	if v.createTime >= t.Unix()+Int64SecondsPerDay {
		// 还没建号
		return false
	}

	createDay := (int(t.Unix()+Int64SecondsPerDay-v.createTime) / SecondsPerDay) + 1
	if createDay != 1 {
		return false
	}

	_time := t.Add(time.Duration(rand.Intn(SecondsPerDay)) * time.Second)
	sinceDuration := v.createTime - startServiceUnixTime
	recordWriter(b, _time.Unix(), s.pid, s.sid, v.id, v.createTime, sinceDuration)

	return true
}

func (l *online_log) writeData(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, v *hero,
	recordWriter DataWriter_online_log) bool {

	if v.createTime >= t.Unix()+Int64SecondsPerDay {
		// 还没建号
		return false
	}

	createDay := (int(t.Unix()+Int64SecondsPerDay-v.createTime) / SecondsPerDay) + 1
	if v.onlineTotalDay != 0 {
		// 建号第X天

		if v.onlineTotalDay < createDay {
			// 超出在线天数
			return false
		}

		// 今天不登陆
		if (createDay-1)%v.onlineOffsetDay != 0 {
			return false
		}
	}

	// 记录3次日志
	write := func(state string, onlineSeconds int) {

		_time := t
		onlineDuration := 0
		offlineDuration := 0
		switch state {
		case "login":
			_time = t.Add(time.Duration(SecondsPerDay-onlineSeconds) / 2 * time.Second)
			onlineDuration = 0
			offlineDuration = SecondsPerDay - onlineSeconds
		case "logout":
			_time = t.Add(time.Duration(SecondsPerDay+onlineSeconds) / 2 * time.Second)
			onlineDuration = onlineSeconds
			offlineDuration = 0
		case "online":
			_time = t.Add(time.Duration(SecondsPerDay/2) * time.Second)
			onlineDuration = onlineSeconds / 2
			offlineDuration = 0
		default:
			logrus.Panicf("unkown state, %s", state)
		}

		recordWriter(b, _time.Unix(), s.pid, s.sid,
			v.id, fmt.Sprintf("name-%v", v.id), v.createTime, v.getRemainYuanbao(createDay), onlineDuration, offlineDuration, state)
	}

	onlineSeconds := rand.Intn(SecondsPerDay)
	write("login", onlineSeconds)
	write("online", onlineSeconds)
	write("logout", onlineSeconds)

	return true
}

func (l *new_guide_log) writeData(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, v *hero,
	recordWriter DataWriter_new_guide_log) bool {

	diff := v.createTime - t.Unix()
	if diff < 0 || diff >= Int64SecondsPerDay {
		// 不是建号当天
		return false
	}

	// 记录3次日志
	write := func(progress int, isComplete bool) {

		_time := t.Add(time.Duration(rand.Intn(SecondsPerDay)) * time.Second)

		c := 0
		if isComplete {
			c = 1
		}

		recordWriter(b, _time.Unix(), s.pid, s.sid,
			v.id, v.name, progress, c)

	}

	write(0, false)
	if v.progress != 0 || v.completed {
		write(v.progress, v.completed)
	}

	return true
}

func (l *device_login_log) writeData(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, v *hero,
	recordWriter DataWriter_device_login_log) bool {

	if v.createTime >= t.Unix()+Int64SecondsPerDay {
		// 还没建号
		return false
	}

	if v.onlineTotalDay != 0 {
		// 建号第X天
		createDay := (int(t.Unix()+Int64SecondsPerDay-v.createTime) / SecondsPerDay) + 1

		if v.onlineTotalDay < createDay {
			// 超出在线天数
			return false
		}

		// 今天不登陆
		if (createDay-1)%v.onlineOffsetDay != 0 {
			return false
		}
	}

	// 记录3次日志
	write := func(deviceId, osType string) {
		_time := t.Add(time.Duration(rand.Intn(SecondsPerDay)) * time.Second)
		recordWriter(b, _time.Unix(), s.pid, s.sid,
			v.id, deviceId, osType)
	}

	for i, deviceId := range v.deviceIds {
		write(deviceId, v.osTypes[i])
	}

	return true
}

func (l *recharge_log) writeData(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, v *hero,
	recordWriter DataWriter_recharge_log) bool {

	if v.firstRechargeDay == 0 {
		return false
	}

	if v.createTime >= t.Unix()+Int64SecondsPerDay {
		// 还没建号
		return false
	}

	// 建号第X天
	createDay := (int(t.Unix()+Int64SecondsPerDay-v.createTime) / SecondsPerDay) + 1

	if v.firstRechargeDay != createDay {
		return false
	}

	_time := t.Add(time.Duration(rand.Intn(SecondsPerDay)) * time.Second)
	recordWriter(b, _time.Unix(), s.pid, s.sid,
		uuid.New().String(), v.rechargeType, v.rechargeAmount, v.id, v.name, v.createTime, v.rechargeLevel)

	return true
}

func (l *yuanbao_log) writeData(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, v *hero,
	recordWrite DataWriter_yuanbao_log) bool {

	if v.firstRechargeDay == 0 {
		return false
	}

	if v.createTime >= t.Unix()+Int64SecondsPerDay {
		// 还没建号
		return false
	}

	// 建号第X天
	createDay := (int(t.Unix()+Int64SecondsPerDay-v.createTime) / SecondsPerDay) + 1

	if v.firstRechargeDay != createDay {
		return false
	}

	write := func(opType string, changeAmount, remainAmount int) {
		_time := t.Add(time.Duration(rand.Intn(SecondsPerDay)) * time.Second)

		recordWrite(b, _time.Unix(), s.pid, s.sid, v.id, v.name, v.createTime, opType, changeAmount, remainAmount)
	}

	// 充值
	write("充值", v.firstRechargeDay, v.firstRechargeDay)

	// 消费
	write("消费", -v.firstRechargeDay/2, v.firstRechargeDay-v.firstRechargeDay/2)

	return true
}

func (f *format) writeDtTimeFile(basePath string, today time.Time, b *bytes.Buffer, pid, sid int) error {
	// 生成的文件路径
	filePath := path.Join(basePath, "dt="+today.Format("2006-01-02"), fmt.Sprintf("%d-%d.%s", pid, sid, f.name))

	// 写入文件
	if err := WriteFile(filePath, b.Bytes()); err != nil {
		return errors.Wrapf(err, "写入文件失败, %s", filePath)
	}

	b.Reset()

	return nil
}

func (f *format) writeSingleFile(basePath string, b *bytes.Buffer, pid, sid int) error {
	// 生成的文件路径
	filePath := path.Join(basePath, fmt.Sprintf("%d-%d.%s", pid, sid, f.name))

	// 写入文件
	if err := WriteFile(filePath, b.Bytes()); err != nil {
		return errors.Wrapf(err, "写入文件失败, %s", filePath)
	}

	b.Reset()

	return nil
}

func WriteFile(filename string, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := os.MkdirAll(path.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, os.ModePerm)
}

// -------------- auto gen 分割线 --------------------

type create_hero_log struct{}

func (l *create_hero_log) getName() string { return "create_hero_log" }

func (l *create_hero_log) getCsvRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeCsvData_create_hero_log)
		}
	}
}

func (l *create_hero_log) getMysqlRecordWriter() RecordWriter {
	return newRecordWriter_create_hero_log(l.getName(), l.writeData)
}

func (l *create_hero_log) getEventRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeEventData_create_hero_log)
		}
	}
}

type DataWriter_create_hero_log func(b *bytes.Buffer, _time int64, pid, sid, id, create_time, since_duration interface{})

func writeCsvData_create_hero_log(b *bytes.Buffer, _time int64, pid, sid,
	id, create_time, since_duration interface{}) {

	b.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v\n", pid, sid, _time,
		id, create_time, since_duration,
	))
}

func writeMysqlData_create_hero_log(b *bytes.Buffer, _time int64, pid, sid,
	id, create_time, since_duration interface{}) {
	b.WriteString(fmt.Sprintf(
		"('%v','%v','%v','%v','%v','%v','%v')",
		pid, sid, _time,
		id, create_time, since_duration,
		time.Unix(_time, 0).Format("2006-01-02")))
}

func newRecordWriter_create_hero_log(
	tableName string,
	f func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, h *hero, w DataWriter_create_hero_log) bool,
) func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {

	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		b.WriteString("begin;\n")

		idx := 0
		for _, v := range heroMap {

			f(s, b, startServerTime, t, v, func(b *bytes.Buffer, _time int64, pid, sid, id, create_time, since_duration interface{}) {
				if idx%10000 == 0 {
					b.WriteString("insert into ")
					b.WriteString(tableName)
					b.WriteString(" values ")
				} else {
					b.WriteString(", ")
				}

				writeMysqlData_create_hero_log(b, _time, pid, sid, id, create_time, since_duration)

				idx++
				if idx%10000 == 0 {
					b.WriteString(";\n")
				}
			})
		}

		if idx%10000 != 0 {
			b.WriteString(";\n")
		}

		b.WriteString("commit;\n")
	}
}

func writeEventData_create_hero_log(b *bytes.Buffer, _time int64, pid, sid,
	id, create_time, since_duration interface{}) {
	eventlog.Commit(eventlog.NewEvent("create_hero_log", pid.(uint32), sid.(uint32)).
		With("id", id).
		With("create_time", create_time).
		With("since_duration", since_duration).
		WithTime(time.Unix(_time, 0)))
}

type online_log struct{}

func (l *online_log) getName() string { return "online_log" }

func (l *online_log) getCsvRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeCsvData_online_log)
		}
	}
}

func (l *online_log) getMysqlRecordWriter() RecordWriter {
	return newRecordWriter_online_log(l.getName(), l.writeData)
}

func (l *online_log) getEventRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeEventData_online_log)
		}
	}
}

type DataWriter_online_log func(b *bytes.Buffer, _time int64, pid, sid, id, name, create_time, remain_yuanbao, online_duration, offline_duration, state interface{})

func writeCsvData_online_log(b *bytes.Buffer, _time int64, pid, sid,
	id, name, create_time, remain_yuanbao, online_duration, offline_duration, state interface{}) {

	b.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n", pid, sid, _time,
		id, name, create_time, remain_yuanbao, online_duration, offline_duration, state,
	))
}

func writeMysqlData_online_log(b *bytes.Buffer, _time int64, pid, sid,
	id, name, create_time, remain_yuanbao, online_duration, offline_duration, state interface{}) {
	b.WriteString(fmt.Sprintf(
		"('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v')",
		pid, sid, _time,
		id, name, create_time, remain_yuanbao, online_duration, offline_duration, state,
		time.Unix(_time, 0).Format("2006-01-02")))
}

func newRecordWriter_online_log(
	tableName string,
	f func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, h *hero, w DataWriter_online_log) bool,
) func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {

	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		b.WriteString("begin;\n")

		idx := 0
		for _, v := range heroMap {

			f(s, b, startServerTime, t, v, func(b *bytes.Buffer, _time int64, pid, sid, id, name, create_time, remain_yuanbao, online_duration, offline_duration, state interface{}) {
				if idx%10000 == 0 {
					b.WriteString("insert into ")
					b.WriteString(tableName)
					b.WriteString(" values ")
				} else {
					b.WriteString(", ")
				}

				writeMysqlData_online_log(b, _time, pid, sid, id, name, create_time, remain_yuanbao, online_duration, offline_duration, state)

				idx++
				if idx%10000 == 0 {
					b.WriteString(";\n")
				}
			})
		}

		if idx%10000 != 0 {
			b.WriteString(";\n")
		}

		b.WriteString("commit;\n")
	}
}

func writeEventData_online_log(b *bytes.Buffer, _time int64, pid, sid,
	id, name, create_time, remain_yuanbao, online_duration, offline_duration, state interface{}) {
	eventlog.Commit(eventlog.NewEvent("online_log", pid.(uint32), sid.(uint32)).
		With("id", id).
		With("name", name).
		With("create_time", create_time).
		With("remain_yuanbao", remain_yuanbao).
		With("online_duration", online_duration).
		With("offline_duration", offline_duration).
		With("state", state).
		WithTime(time.Unix(_time, 0)))
}

type new_guide_log struct{}

func (l *new_guide_log) getName() string { return "new_guide_log" }

func (l *new_guide_log) getCsvRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeCsvData_new_guide_log)
		}
	}
}

func (l *new_guide_log) getMysqlRecordWriter() RecordWriter {
	return newRecordWriter_new_guide_log(l.getName(), l.writeData)
}

func (l *new_guide_log) getEventRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeEventData_new_guide_log)
		}
	}
}

type DataWriter_new_guide_log func(b *bytes.Buffer, _time int64, pid, sid, id, name, progress, completed interface{})

func writeCsvData_new_guide_log(b *bytes.Buffer, _time int64, pid, sid,
	id, name, progress, completed interface{}) {

	b.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v\n", pid, sid, _time,
		id, name, progress, completed,
	))
}

func writeMysqlData_new_guide_log(b *bytes.Buffer, _time int64, pid, sid,
	id, name, progress, completed interface{}) {
	b.WriteString(fmt.Sprintf(
		"('%v','%v','%v','%v','%v','%v','%v','%v')",
		pid, sid, _time,
		id, name, progress, completed,
		time.Unix(_time, 0).Format("2006-01-02")))
}

func newRecordWriter_new_guide_log(
	tableName string,
	f func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, h *hero, w DataWriter_new_guide_log) bool,
) func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {

	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		b.WriteString("begin;\n")

		idx := 0
		for _, v := range heroMap {

			f(s, b, startServerTime, t, v, func(b *bytes.Buffer, _time int64, pid, sid, id, name, progress, completed interface{}) {
				if idx%10000 == 0 {
					b.WriteString("insert into ")
					b.WriteString(tableName)
					b.WriteString(" values ")
				} else {
					b.WriteString(", ")
				}

				writeMysqlData_new_guide_log(b, _time, pid, sid, id, name, progress, completed)

				idx++
				if idx%10000 == 0 {
					b.WriteString(";\n")
				}
			})
		}

		if idx%10000 != 0 {
			b.WriteString(";\n")
		}

		b.WriteString("commit;\n")
	}
}

func writeEventData_new_guide_log(b *bytes.Buffer, _time int64, pid, sid,
	id, name, progress, completed interface{}) {
	eventlog.Commit(eventlog.NewEvent("new_guide_log", pid.(uint32), sid.(uint32)).
		With("id", id).
		With("name", name).
		With("progress", progress).
		With("completed", completed).
		WithTime(time.Unix(_time, 0)))
}

type device_login_log struct{}

func (l *device_login_log) getName() string { return "device_login_log" }

func (l *device_login_log) getCsvRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeCsvData_device_login_log)
		}
	}
}

func (l *device_login_log) getMysqlRecordWriter() RecordWriter {
	return newRecordWriter_device_login_log(l.getName(), l.writeData)
}

func (l *device_login_log) getEventRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeEventData_device_login_log)
		}
	}
}

type DataWriter_device_login_log func(b *bytes.Buffer, _time int64, pid, sid, id, device_id, os_type interface{})

func writeCsvData_device_login_log(b *bytes.Buffer, _time int64, pid, sid,
	id, device_id, os_type interface{}) {

	b.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v\n", pid, sid, _time,
		id, device_id, os_type,
	))
}

func writeMysqlData_device_login_log(b *bytes.Buffer, _time int64, pid, sid,
	id, device_id, os_type interface{}) {
	b.WriteString(fmt.Sprintf(
		"('%v','%v','%v','%v','%v','%v','%v')",
		pid, sid, _time,
		id, device_id, os_type,
		time.Unix(_time, 0).Format("2006-01-02")))
}

func newRecordWriter_device_login_log(
	tableName string,
	f func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, h *hero, w DataWriter_device_login_log) bool,
) func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {

	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		b.WriteString("begin;\n")

		idx := 0
		for _, v := range heroMap {

			f(s, b, startServerTime, t, v, func(b *bytes.Buffer, _time int64, pid, sid, id, device_id, os_type interface{}) {
				if idx%10000 == 0 {
					b.WriteString("insert into ")
					b.WriteString(tableName)
					b.WriteString(" values ")
				} else {
					b.WriteString(", ")
				}

				writeMysqlData_device_login_log(b, _time, pid, sid, id, device_id, os_type)

				idx++
				if idx%10000 == 0 {
					b.WriteString(";\n")
				}
			})
		}

		if idx%10000 != 0 {
			b.WriteString(";\n")
		}

		b.WriteString("commit;\n")
	}
}

func writeEventData_device_login_log(b *bytes.Buffer, _time int64, pid, sid,
	id, device_id, os_type interface{}) {
	eventlog.Commit(eventlog.NewEvent("device_login_log", pid.(uint32), sid.(uint32)).
		With("id", id).
		With("device_id", device_id).
		With("os_type", os_type).
		WithTime(time.Unix(_time, 0)))
}

type recharge_log struct{}

func (l *recharge_log) getName() string { return "recharge_log" }

func (l *recharge_log) getCsvRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeCsvData_recharge_log)
		}
	}
}

func (l *recharge_log) getMysqlRecordWriter() RecordWriter {
	return newRecordWriter_recharge_log(l.getName(), l.writeData)
}

func (l *recharge_log) getEventRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeEventData_recharge_log)
		}
	}
}

type DataWriter_recharge_log func(b *bytes.Buffer, _time int64, pid, sid, order_id, order_type, order_amount, id, name, create_time, level interface{})

func writeCsvData_recharge_log(b *bytes.Buffer, _time int64, pid, sid,
	order_id, order_type, order_amount, id, name, create_time, level interface{}) {

	b.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n", pid, sid, _time,
		order_id, order_type, order_amount, id, name, create_time, level,
	))
}

func writeMysqlData_recharge_log(b *bytes.Buffer, _time int64, pid, sid,
	order_id, order_type, order_amount, id, name, create_time, level interface{}) {
	b.WriteString(fmt.Sprintf(
		"('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v','%v')",
		pid, sid, _time,
		order_id, order_type, order_amount, id, name, create_time, level,
		time.Unix(_time, 0).Format("2006-01-02")))
}

func newRecordWriter_recharge_log(
	tableName string,
	f func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, h *hero, w DataWriter_recharge_log) bool,
) func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {

	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		b.WriteString("begin;\n")

		idx := 0
		for _, v := range heroMap {

			f(s, b, startServerTime, t, v, func(b *bytes.Buffer, _time int64, pid, sid, order_id, order_type, order_amount, id, name, create_time, level interface{}) {
				if idx%10000 == 0 {
					b.WriteString("insert into ")
					b.WriteString(tableName)
					b.WriteString(" values ")
				} else {
					b.WriteString(", ")
				}

				writeMysqlData_recharge_log(b, _time, pid, sid, order_id, order_type, order_amount, id, name, create_time, level)

				idx++
				if idx%10000 == 0 {
					b.WriteString(";\n")
				}
			})
		}

		if idx%10000 != 0 {
			b.WriteString(";\n")
		}

		b.WriteString("commit;\n")
	}
}

func writeEventData_recharge_log(b *bytes.Buffer, _time int64, pid, sid,
	order_id, order_type, order_amount, id, name, create_time, level interface{}) {
	eventlog.Commit(eventlog.NewEvent("recharge_log", pid.(uint32), sid.(uint32)).
		With("order_id", order_id).
		With("order_type", order_type).
		With("order_amount", order_amount).
		With("id", id).
		With("name", name).
		With("create_time", create_time).
		With("level", level).
		WithTime(time.Unix(_time, 0)))
}

type yuanbao_log struct{}

func (l *yuanbao_log) getName() string { return "yuanbao_log" }

func (l *yuanbao_log) getCsvRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeCsvData_yuanbao_log)
		}
	}
}

func (l *yuanbao_log) getMysqlRecordWriter() RecordWriter {
	return newRecordWriter_yuanbao_log(l.getName(), l.writeData)
}

func (l *yuanbao_log) getEventRecordWriter() RecordWriter {
	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		for _, v := range heroMap {
			l.writeData(s, b, startServerTime, t, v, writeEventData_yuanbao_log)
		}
	}
}

type DataWriter_yuanbao_log func(b *bytes.Buffer, _time int64, pid, sid, id, name, create_time, op_type, change_amount, remain_amount interface{})

func writeCsvData_yuanbao_log(b *bytes.Buffer, _time int64, pid, sid,
	id, name, create_time, op_type, change_amount, remain_amount interface{}) {

	b.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v\n", pid, sid, _time,
		id, name, create_time, op_type, change_amount, remain_amount,
	))
}

func writeMysqlData_yuanbao_log(b *bytes.Buffer, _time int64, pid, sid,
	id, name, create_time, op_type, change_amount, remain_amount interface{}) {
	b.WriteString(fmt.Sprintf(
		"('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v')",
		pid, sid, _time,
		id, name, create_time, op_type, change_amount, remain_amount,
		time.Unix(_time, 0).Format("2006-01-02")))
}

func newRecordWriter_yuanbao_log(
	tableName string,
	f func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, h *hero, w DataWriter_yuanbao_log) bool,
) func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {

	return func(s *server_info, b *bytes.Buffer, startServerTime, t time.Time, heroMap map[int64]*hero) {
		b.WriteString("begin;\n")

		idx := 0
		for _, v := range heroMap {

			f(s, b, startServerTime, t, v, func(b *bytes.Buffer, _time int64, pid, sid, id, name, create_time, op_type, change_amount, remain_amount interface{}) {
				if idx%10000 == 0 {
					b.WriteString("insert into ")
					b.WriteString(tableName)
					b.WriteString(" values ")
				} else {
					b.WriteString(", ")
				}

				writeMysqlData_yuanbao_log(b, _time, pid, sid, id, name, create_time, op_type, change_amount, remain_amount)

				idx++
				if idx%10000 == 0 {
					b.WriteString(";\n")
				}
			})
		}

		if idx%10000 != 0 {
			b.WriteString(";\n")
		}

		b.WriteString("commit;\n")
	}
}

func writeEventData_yuanbao_log(b *bytes.Buffer, _time int64, pid, sid,
	id, name, create_time, op_type, change_amount, remain_amount interface{}) {
	eventlog.Commit(eventlog.NewEvent("yuanbao_log", pid.(uint32), sid.(uint32)).
		With("id", id).
		With("name", name).
		With("create_time", create_time).
		With("op_type", op_type).
		With("change_amount", change_amount).
		With("remain_amount", remain_amount).
		WithTime(time.Unix(_time, 0)))
}
