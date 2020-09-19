package bootstrap

import (
	"fmt"
	"github.com/lightpaw/eventlog"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/build"
	"github.com/lightpaw/male7/gen/service"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/worker"
	"math/rand"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
	"github.com/lightpaw/male7/util/grmon"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/timer"
	"os/signal"
	"syscall"
)

func Start(pprofAddr string, version bool, createMuxListener iface.ServeListener) {

	rand.Seed(time.Now().UnixNano())

	logrus.Info("main")
	logrus.Info("server: ", build.GetVersion())
	logrus.Info("config: ", build.GetConfigVersion())
	logrus.Info("buildtime: ", build.GetBuildTime())
	logrus.Info("buildunix: ", build.GetBuildUnixTime())
	logrus.Info("tag: ", build.GitTag)

	if version {
		return
	}

	if len(pprofAddr) > 0 {
		grmon.Start(pprofAddr)
	}

	// 系统初始化
	service.Init()

	// 监听自定义信号
	listenSigs()

	if !service.AwsService.InitFirehoseEventLog() {
		eventlog.Start(eventlog.NewLocalFileDestination("log/event.log"), 1*time.Second)
	}

	// 开启监控
	metricsHandler, err := service.GameExporter.Start()
	if err != nil {
		logrus.WithError(err).Panic("开启监控报错")
	}

	// TODO 端口监听失败，这里可能把开着的服的数据改写了，导致bug，检查如果端口被监听了，直接报错

	// 开启http服务器
	startHttpServer(metricsHandler)

	service.ServerStartStopTimeService.SaveStartTime()

	// tlog 服务器状态
	tlogGameSvrState()

	service.GameServer.Serve(createMuxListener, worker.NewMessageWorker)

	// 开始关闭服务器，保存数据等等
	call.CatchPanic(service.ClusterService.Close, "close ClusterService")
	call.CatchPanic(service.WorldService.Close, "close WorldService")
	call.CatchPanic(service.RealmService.Close, "close RealmService") // 早点把野外关掉
	call.CatchPanic(service.XiongNuModule.Close, "close XiongNuModule")
	call.CatchPanic(service.TowerModule.Close, "close TowerModule")
	call.CatchPanic(service.BaiZhanModule.Close, "close BaiZhanModule")
	call.CatchPanic(service.RankModule.Close, "close RankModule")
	call.CatchPanic(service.FarmService.Close, "close FarmService")
	call.CatchPanic(service.XuanyuanModule.Close, "close XuanyuanModule")
	call.CatchPanic(service.HebiModule.Close, "close HebiModule")
	call.CatchPanic(service.MingcWarService.Close, "close MingcWarService")
	call.CatchPanic(service.TimeLimitGiftService.Close, "close TimeLimitGiftService")
	call.CatchPanic(service.RedPacketService.Close, "close RedPacketService")
	call.CatchPanic(service.ActivityModule.Close, "close ActivityModule")

	// 被其他模块依赖的部分数据
	call.CatchPanic(service.MingcService.Close, "close MingcService")
	call.CatchPanic(service.CountryService.Close, "close CountryService")
	call.CatchPanic(service.GuildService.Close, "close GuildService")
	call.CatchPanic(service.HeroDataService.Close, "close HeroDataService")

	// 对其他模块没有依赖的，最后单独关闭
	call.CatchPanic(service.BroadcastService.Close, "close BroadcastService")
	call.CatchPanic(service.ServerStartStopTimeService.SaveStopTime, "close ServerStartStopTimeService.SaveStopTime")
	call.CatchPanic(service.TickerService.Close, "close TickerService")

	// 其他进程模块
	call.CatchPanic(service.FightService.Close, "close FightService")

	call.CatchLoopPanic(service.TlogService.Close, "close TlogService")
	call.CatchLoopPanic(service.KafkaService.Close, "close KafkaService")

	call.CatchPanic(service.TssClient.Close, "close TssClient")
	call.CatchPanic(eventlog.Flush, "close eventlog.Flush")

	// 最后关闭DB
	call.CatchPanic(func() {
		if err := service.DbService.Close(); err != nil {
			logrus.WithError(err).Error("service.DbService.Close() fail")
		}
	}, "close cluster")

	logrus.Info("服务器已关闭")
}

//func startServer(addr string) error {
//	config := &nonmux.Config{CloseWaitTime: 30 * time.Second, Unmarshaller: unmarshal.NewProtoUnmarshaller(), DontEncrypt: service.IndividualServerConfig.GetDontEncrypt()}
//	listener, err := nonmux.Listen(addr, config)
//	if err != nil {
//		return errors.Wrap(err, "监听rmux失败")
//	}
//
//	defer func() {
//		if r := recover(); r != nil {
//			// 严重错误. 这里不能panic，后面需要关闭服务器
//			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Error("game server recovered from panic. SEVERE!!!")
//			metrics.IncPanic()
//		}
//
//		// 服务器断开连接了，把所有英雄T下线
//		logrus.Info("worker server exited, start kick user")
//	}()
//
//	go call.CatchLoopPanic(func() {
//		sigs := make(chan os.Signal, 1)
//		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
//
//		select {
//		case <-sigs:
//			listener.Close()
//		}
//	}, "main.sigs")
//
//	logrus.WithField("addr", addr).Info("开始监听连接")
//	for {
//		conn, err := listener.Accept()
//		if err != nil {
//			return err
//		}
//
//		if worker.IsRobotConn(conn) && !service.IndividualServerConfig.GetIsAllowRobot() {
//			logrus.Debug("收到机器人登陆，服务器不允许机器人登陆，踢掉")
//			conn.Close()
//			continue
//		}
//
//		go worker.NewMessageWorker(conn)
//	}
//}

func startHttpServer(metricsHandler http.Handler) {
	root := "temp/"
	os.MkdirAll(path.Dir(root), os.ModePerm)

	go call.CatchLoopPanic(func() {
		logrus.Infof("启动HttpServer[%d], 回放本地地址: %v", service.IndividualServerConfig.GetHttpPort(), root)

		// 回放服务器
		prefix := "/" + service.IndividualServerConfig.GetReplayPrefix() // like /replay/
		http.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(root))))

		// 监控服务器，promethus metrics
		http.Handle("/metrics", metricsHandler)

		err := http.ListenAndServe(fmt.Sprintf(":%d", service.IndividualServerConfig.GetHttpPort()), nil)
		logrus.WithError(err).Errorf("启动HttpServer失败")
	}, "main.startHttpServer")
}

func listenSigs() {
	c := make(chan os.Signal, 1)
	//signal.Notify(c, syscall.SIGUSR1)
	signal.Notify(c, syscall.SIGHUP)
	go call.CatchPanic(func() {
		for range c {
			doDumpStacks()
		}
	}, "main.listenSigs")
}

func doDumpStacks() {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, true)]
	logrus.Errorf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===", buf)
}

func tlogGameSvrState() {
	go func() {
		loopWheel := timer.NewTimingWheel(time.Minute, 32)
		min5Tick := loopWheel.After(5 * time.Minute)

		//loopWheel := timer.NewTimingWheel(500*time.Millisecond, 32)
		//secondTick := loopWheel.After(time.Second)

		service.TlogService.TlogGameSvrState()

		for {
			select {

			case <-min5Tick:
				logrus.Debugf("========== exec 5 tlogGameSvrState")
				min5Tick = loopWheel.After(5 * time.Minute)
				service.TlogService.TlogGameSvrState()

				//case <-secondTick:
				//	logrus.Debugf("========== exec 1 tlogGameSvrState")
				//	secondTick = loopWheel.After(time.Second)
				//	service.TlogService.TlogGameSvrStateDetail(0)

			}
		}

	}()
}
