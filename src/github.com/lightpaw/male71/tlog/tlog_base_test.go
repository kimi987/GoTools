package tlog

import (
	"testing"
	"fmt"
	"time"
	. "github.com/onsi/gomega"
	"github.com/lightpaw/male7/gen/ifacemock"
	"sync"
	"math/rand"
	"github.com/lightpaw/male7/config/kv"
)

func TestTlogService(t *testing.T) {
	RegisterTestingT(t)

	testCalcLatestNewFileTime()
	//testWriteTlog()
}

func testWriteTlog() {
	s := newMockTlogBaseService()

	crntCount := 10000
	wait := sync.WaitGroup{}
	falseCount := 0
	d, _ := time.ParseDuration("100ms") // cache 65536, 1w人，10w条/s，7mb/s, 本地 ssd
	for j := 0; j < crntCount; j++ {
		wait.Add(1)
		go func(id int) {
			defer wait.Done()
			dt, _ := time.ParseDuration(fmt.Sprintf("%vms", rand.Intn(10)))
			time.Sleep(dt)
			for i := 1; i <= 10; i++ {
				if succ := s.WriteTlog(fmt.Sprintf("tlog-%v", i*1000000+id) + "\n"); !succ {
					falseCount++
				}
				time.Sleep(d)
			}
		}(j)

	}
	wait.Wait()
	s.Close()
	Ω(falseCount).Should(Equal(0))
}

func newMockTlogBaseService() *TlogBaseService {
	timeService := ifacemock.TimeService
	timeService.Mock(timeService.CurrentTime, time.Now)

	config := kv.NewIndividualServerConfig()
	s := NewTlogBaseService(timeService, ifacemock.KafkaService, config)
	return s
}

func testCalcLatestNewFileTime() {
	East8 := time.FixedZone("East-8", 8*60*60)
	ctime := time.Date(2000, 1, 1, 0, 0, 0, 0, East8)
	duration, _ := time.ParseDuration("1m")
	ntime := calcNextNewFileTime(ctime, duration)
	result := !ntime.After(ctime.Add(duration)) && !ntime.Before(ctime)
	Ω(result).Should(BeTrue())
}

func TestNewFileName(t *testing.T) {

	ctime := time.Now()
	fmt.Println(ctime)
	fmt.Println(calcNextNewFileTime(ctime, time.Hour))
}
