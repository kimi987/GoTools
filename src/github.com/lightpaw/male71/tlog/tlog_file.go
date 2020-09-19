package tlog

import (
	"path"
	"os"
	"github.com/lightpaw/logrus"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"bufio"
	"github.com/lightpaw/male7/gen/iface"
	"path/filepath"
	"github.com/lightpaw/male7/util/call"
)

/*
	========= tlog file =============
 */

type TlogFile struct {
	timeService iface.TimeService

	logChan   chan string
	closeChan chan struct{}
	loopChan  chan struct{}

	fileCurrentPath string
	fileArchivePath string
	fileNamePrefix  string

	rotateDuration time.Duration
	rotateSize     int
	writeBuffer    int

	nextNewFileTime time.Time
	currentFileName string
	currentFileSize int
	tlogFile        *os.File
	writer          *bufio.Writer
}

func NewTlogFile(timeService iface.TimeService, fileCurrentDir, fileArchiveDir, filenamePrefix string,
	rotateDuration time.Duration, rotateSize, writeBuffer, chanSize int) *TlogFile {
	s := &TlogFile{}

	s.timeService = timeService
	s.logChan = make(chan string, chanSize)
	s.closeChan = make(chan struct{})
	s.loopChan = make(chan struct{})

	s.fileCurrentPath = fileCurrentDir
	s.fileArchivePath = fileArchiveDir
	s.fileNamePrefix = filenamePrefix

	s.rotateDuration = rotateDuration
	s.rotateSize = rotateSize
	s.writeBuffer = writeBuffer

	ctime := timeService.CurrentTime()
	s.nextNewFileTime = calcNextNewFileTime(ctime, rotateDuration)

	var err error

	// current dir
	s.fileCurrentPath, err = filepath.Abs(fileCurrentDir)
	if err != nil {
		logrus.WithError(err).Panic("TlogFile file_current_base_dir 配置错误")
	}
	os.MkdirAll(s.fileCurrentPath, os.ModePerm)

	// archive dir
	s.fileArchivePath, err = filepath.Abs(fileArchiveDir)
	if err != nil {
		logrus.WithError(err).Panic("TlogFile file_archive_base_dir 配置错误")
	}
	os.MkdirAll(s.fileArchivePath, os.ModePerm)

	// tlog file
	//s.fileNamePrefix = fmt.Sprintf("male7_%v_%v_", config.GetPlatformID(), config.GetServerID())

	s.switchFile(ctime)

	// TODO 看下有没有文件需要移动到归档路径

	go call.CatchLoopPanic(s.loop, "tlog_file")

	return s
}

func (s *TlogFile) loop() {

	defer close(s.loopChan)

	// 每2秒检查尝试flush
	tick := time.NewTicker(2 * time.Second)

	for {
		select {
		case logStr := <-s.logChan:
			// 将数据写入文件
			n, _ := s.writer.WriteString(logStr)
			s.currentFileSize += n

		case <-tick.C:
			s.tick()

		case <-s.closeChan:
			return
		}
	}
}

func (s *TlogFile) tick() {
	ctime := s.timeService.CurrentTime()

	// 检查时间是否到了要更新文件
	if s.rotateDuration > 0 && ctime.After(s.nextNewFileTime) {
		s.nextNewFileTime = calcNextNewFileTime(ctime, s.rotateDuration)
		s.switchFile(ctime)
		return
	}

	// 检查文件大小是否超出限制
	if s.rotateSize > 0 && s.currentFileSize >= s.rotateSize {
		s.switchFile(ctime)
		return
	}

	// flush一下
	s.writer.Flush()
}

func (s *TlogFile) switchFile(ctime time.Time) bool {

	// 新的文件名
	filename := s.newFileName(ctime)

	if s.currentFileName == filename {
		logrus.Error("创建新的日志文件，发现新的文件名跟当前文件名一样")
		return false
	}

	// 创建新文件
	newFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		logrus.WithError(err).WithField("file", filename).Errorf("tlogBaseService 打开新的 tlog 文件失败")
		return false
	}

	// 停止写入旧文件，将旧文件移动到归档目录
	if s.writer != nil {
		s.writer.Flush()
	}

	if s.tlogFile != nil {
		s.tlogFile.Close()

		// 移动到归档目录
		s.doMove2ArchiveDir(s.currentFileName)
	}

	s.currentFileName = filename
	s.currentFileSize = 0
	s.tlogFile = newFile
	s.writer = bufio.NewWriterSize(newFile, s.writeBuffer)
	return true
}

func (s *TlogFile) doMove2ArchiveDir(filePath string) {
	filename := path.Base(filePath)

	newFilePath := path.Join(s.fileArchivePath, filename)

	if err := os.Rename(filePath, newFilePath); err != nil {
		logrus.WithError(err).Error("将日志文件移动到归档目录失败")
	}

}

func (s *TlogFile) Close() {
	close(s.closeChan)
	<-s.loopChan

	// 检查所有
	close(s.logChan)

	if len(s.logChan) > 0 {
		for v := range s.logChan {
			s.writer.WriteString(v)
		}
	}

	if s.writer != nil {
		s.writer.Flush()
	}

	if s.tlogFile != nil {
		s.tlogFile.Close()

		s.doMove2ArchiveDir(s.currentFileName)
	}

	logrus.Debugf("TlogFile 关闭成功")
}

func (s *TlogFile) AddLog(content string) {

	select {
	case s.logChan <- content:
		return
	case <-time.After(100 * time.Millisecond):
		logrus.WithField("content", content).Info("添加tlog超时")
		return
	}
}

func calcNextNewFileTime(ctime time.Time, duration time.Duration) time.Time {
	midnightTime := timeutil.DailyTime.PrevTime(ctime)
	diff := ctime.Sub(midnightTime)
	return ctime.Add(duration - (diff % duration))
}

// /prefix_{yyyyMMddHHmmss}.log
func (s *TlogFile) newFileName(ctime time.Time) string {
	simpleName := s.fileNamePrefix + ctime.Format("20060102150405") + ".log"
	return path.Join(s.fileCurrentPath, simpleName)
}
