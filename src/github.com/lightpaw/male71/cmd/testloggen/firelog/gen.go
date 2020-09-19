package main

//import (
//	"fmt"
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/service/firehose"
//	"github.com/google/uuid"
//	"github.com/lightpaw/eventlog"
//	"github.com/lightpaw/male7/gamelogs"
//	"math/rand"
//	"time"
//)
//
//func main() {
//
//	session := session.Must(session.NewSession())
//	f := firehose.New(session, aws.NewConfig().WithRegion("us-west-2"))
//
//	eventlog.Start(eventlog.NewFirehoseDestination("male7-firelog", f), 10*time.Second)
//
//	ctime := time.Now().Unix()
//
//	pid := 1
//	for i := 0; i < 100; i++ {
//		sid := (i % 3) + 1
//
//		id := int64(i + 10000)
//		name := fmt.Sprintf("player_%v", id)
//		createTime := ctime - rand.Int63n(60*60)
//
//		gamelogs.CreateHeroLog(pid, sid, id, createTime, rand.Int63n(100000))
//		gamelogs.OnlineLog(pid, sid, id, name, createTime, 20,
//			rand.Int63n(100000), rand.Int63n(100000), "login")
//
//		gamelogs.YuanbaoLog(pid, sid, id, name, createTime, "MA", 10, 10)
//
//		gamelogs.DeviceLoginLog(pid, sid, id, uuid.New().String(), "ios")
//
//		progress := (rand.Int31n(10) + 1) * 10000
//		gamelogs.NewGuideLog(pid, sid, id, name, progress, progress >= 100000)
//
//		gamelogs.RechargeLog(pid, sid, uuid.New().String(), "YB", 10, id, name, createTime, 10)
//
//	}
//
//	eventlog.Flush()
//
//}
