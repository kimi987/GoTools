// message worker
package worker

import (
	"context"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/login"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/service"
	service1 "github.com/lightpaw/male7/service"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/service/unmarshal"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/timer"
	"github.com/lightpaw/pbutil"
	"github.com/prometheus/client_golang/prometheus"
	"runtime/debug"
	"sync/atomic"
	"time"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/gamelogs"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity/heroid"
	"github.com/lightpaw/male7/util/i64"
)

const (
	// 每个英雄的事件队列长度
	eventChanSize = 10 // 足够在处理消息的时候, 解析下面要处理的消息了
	loginTimeout  = 30 * time.Second

	autoSaveInterval = 10 * time.Minute // 定期保存的间隔

)

var (
	loginTimeoutWheel = timer.NewTimingWheel(500*time.Millisecond, 64) // 最大timeout是32秒. 英雄每秒update也在用这个
)

type MessageWorker struct {
	// 类似netty中的channel
	session Conn

	closeChan          chan struct{} // 通知主线程退出
	mainLoopExitNotify chan struct{} // 主线程已退出

	closeFlag int32

	uid  int64
	user *service1.ConnectedUser

	lastReceiveClientMsgTime time.Time

	isRobot    bool
	pf         uint32
	clientIp   string
	clientIp32 uint32
}

func NewMessageWorker(session Conn) {
	heroId := heroid.NewHeroId(session.GetLoginToken().UserSelfID, session.GetLoginToken().GameServerID)
	logrus.WithField("addr", session.GetLoginToken().GameServerAddr).WithField("heroId", heroId).Debug("收到新连接")
	result := &MessageWorker{
		session:            session,
		closeChan:          make(chan struct{}),
		mainLoopExitNotify: make(chan struct{}),
		uid:                heroId,
		isRobot:            IsRobotConn(session),
		pf:                 GetPf(session),
		clientIp:           GetClientIp(session),
		clientIp32:         GetConnClientIp32(session),
	}

	result.loop()
}

func IsRobotConn(conn Conn) bool {
	return conn.GetLoginToken().Reserved&1 == 1
}

func GetPf(conn Conn) uint32 {
	return uint32((conn.GetLoginToken().Reserved & 255) >> 1)
}

func GetClientIp(conn Conn) string {
	addr := conn.GetLoginToken().GameServerAddr
	return fmt.Sprintf("%d.%d.%d.%d", addr[0], addr[1], addr[2], addr[3])
}

func GetConnClientIp32(conn Conn) uint32 {
	return util.ToU32Ip(conn.GetLoginToken().GameServerAddr)
}

func (m *MessageWorker) ClientIp() string {
	return m.clientIp
}

func (m *MessageWorker) Id() int64 {
	return m.uid
}

// 发送消息.
func (m *MessageWorker) Send(msg pbutil.Buffer) {
	if msg != nil {
		printS2CLog(msg.Buffer())
		m.session.Write(msg.Buffer())
	}
}

// 发送在线路繁忙时可以被丢掉的消息
func (m *MessageWorker) SendIfFree(msg pbutil.Buffer) {
	if msg != nil {
		printS2CLog(msg.Buffer())
		m.session.WriteIfFree(msg.Buffer())
	}
}

// 全部发送
func (m *MessageWorker) SendAll(msgs []pbutil.Buffer) {
	for _, s := range msgs {
		m.Send(s)
	}

	return
}

// 关闭这个连接. 可以被多次调用只会处理一次下线逻辑.
// 连接断开时也会触发这个方法
func (m *MessageWorker) Close() {
	if !atomic.CompareAndSwapInt32(&m.closeFlag, 0, 1) {
		return
	}

	close(m.closeChan) // notify hero loop and receive loop
	m.session.Close()
}

// 关闭并等待, 不会再有新的消息被执行, 也已经保存完毕
func (m *MessageWorker) CloseAndWait() {
	m.Close()
	<-m.mainLoopExitNotify
}

func (m *MessageWorker) Disconnect(err msg.ErrMsg) {

	if err != misc.ErrDisconectReasonFailGm {
		logrus.WithField("id", m.Id()).WithError(err).Errorf("玩家踢下线")
	}

	if err != nil {
		m.Send(err.ErrMsg())

		if err.ErrMsg().Buffer() != nil && len(err.ErrMsg().Buffer()) >= 4 {
			m.user.SetLogoutType(uint64(err.ErrMsg().Buffer()[3]))
		}
	}

	m.Close()
}

func (m *MessageWorker) DisconnectAndWait(err msg.ErrMsg) {

	if err != misc.ErrDisconectReasonFailClose && err != misc.ErrDisconectReasonFailKick {
		logrus.WithField("id", m.Id()).WithError(err).Errorf("玩家踢下线（等待）")
	}
	if err != nil {
		m.Send(err.ErrMsg())

		if err.ErrMsg().Buffer() != nil && len(err.ErrMsg().Buffer()) >= 4 {
			m.user.SetLogoutType(uint64(err.ErrMsg().Buffer()[3]))
		}
	}

	m.CloseAndWait()
}

// 连接是否已断开. 已断开可能worker还在处理event
func (m *MessageWorker) IsClosed() bool {
	return atomic.LoadInt32(&m.closeFlag) == 1
}

//func (m *MessageWorker) doSave() {
//	if m.user != nil && !m.isRobot {
//		hc := m.user.GetHeroController()
//		if hc != nil {
//			service.DbService.SaveHero(hc.Hero())
//		}
//	}
//}

// 英雄线程. 单独goroutine执行
func (m *MessageWorker) loop() {
	defer func() {
		if r := recover(); r != nil {
			// 严重错误. 英雄线程这里不能panic
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Error("Worker.loop recovered from panic!!! SERIOUS PROBLEM")
			m.Close()

			metrics.IncPanic()
		}

		close(m.mainLoopExitNotify) // 通知所有事件都已处理完
		logrus.Debug("Worker.loop exited")
	}()

	mustLoginTimeout := loginTimeoutWheel.After(loginTimeout)
WAIT_LOGIN:
	select {
	case <-mustLoginTimeout:
		logrus.Debug("长时间不登录, 断开连接")
		m.Close()
		return
	case <-m.session.ClosedNotify():
		m.Close()
		return
	case <-m.closeChan:
		return
	case msg := <-m.session.MsgChan():
		msgData := msg.(*service.MsgData)
		printC2SLog(msgData)
		switch msgData.ModuleID {

		case login.MODULE_ID:
			var heroId int64
			switch msgData.SequenceID {
			case login.C2S_ROBOT_LOGIN:
				if !m.isRobot {
					logrus.Error("不是robot的登陆token，但是收到RobotLogin消息")
					m.Close()
					return
				}

				if !service.IndividualServerConfig.GetIsAllowRobot() {
					logrus.Error("收到RobotLogin消息, 但是不是AllowRobot模式")
					m.Close()
					return
				}

				if proto, ok := msgData.Proto.(*login.C2SRobotLoginProto); ok {
					heroId = int64(proto.Id)
				}

				if heroId == 0 {
					logrus.Error("收到RobotLogin消息, 但是HeroId == 0")
					m.Close()
					return
				}
				fallthrough
			case login.C2S_LOGIN:
				if m.processLoginMsg(heroId, msgData) {
					loginTime := service.TimeService.CurrentTime()
					defer func() {
						// HeroDataService统一管理
						//m.doSave()

						call.CatchPanic(func() {
							if tencentInfo := m.user.TencentInfo(); tencentInfo != nil {
								var tlogHero *entity.TlogHeroInfo
								if hc := m.user.GetHeroController(); hc != nil {
									hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
										if hero != nil {
											tlogHero = hero.BuildFullTlogHeroInfo(service.TimeService.CurrentTime())
										}
										return false
									})
								}

								if tlogHero == nil {
									tlogHero = entity.NewSimpleTlogHeroInfo(m.user.Id(), "")
								}

								tlogLogout(tlogHero, tencentInfo, m.user.LogoutType())
							}
						}, "tlogLogout")

						if ok := service.WorldService.RemoveUserIfSame(m.user); !ok {
							logrus.Error("下线时要删除worldService里的user, 竟然不是我自己")
						}

						// 保存玩家杂项
						if m.user.MiscNeedOfflineSave() {
							ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
								if err := service.DbService.UpdateUserMisc(ctx, m.user.Id(), m.user.Misc()); err != nil {
									logrus.WithError(err).Error("下线时要保存玩家杂项报错")
								}
								return nil
							})
						}

						// 登陆日志（离线记录一下）
						ctime := service.TimeService.CurrentTime()
						loginDuration := i64.Max(int64(ctime.Sub(loginTime)/time.Second), 0)
						gamelogs.HeroOnlineLog(constants.PID, m.user.Sid(), m.user.Id(), loginDuration)

						// 玩家离线时候，没有完成新手引导，记录一下新手引导的状态
						if !m.user.Misc().IsTutorialComplete {
							gamelogs.NewGuideLog(constants.PID, m.user.Sid(), m.user.Id(), m.user.Misc().TutorialProgress, false)
						}
					}()

					if m.user.GetHeroController() == nil {
						goto CREATE_HERO
					} else {
						goto WAIT_LOADED
					}
				} else {
					return
				}

				//case login.C2S_ROBOT_LOGIN:
				//	if m.processRobotLogin(msgData) {
				//		goto RECEIVE_MSG
				//	} else {
				//		return
				//	}

			default:
				logrus.WithField("sequenceID", msgData.SequenceID).Error("登录的第一条消息只能是C2S_LOGIN或者C2S_REBOT_LOGIN")
				m.Send(misc.ERR_DISCONECT_REASON_FAIL_MUST_LOGIN)
				m.Close()
				return
			}

		case misc.MODULE_ID:
			if m.processWaitMiscMsg(msgData) {
				goto WAIT_LOGIN
			} else {
				logrus.WithField("sequenceID", msgData.SequenceID).Error("登录之前的misc模块只能发C2S_CONFIG")
				m.Send(misc.ERR_DISCONECT_REASON_FAIL_MUST_LOGIN)
				m.Close()
				return
			}

		default:
			logrus.WithField("moduleID", msgData.ModuleID).WithField("sequenceID", msgData.SequenceID).Error("登录的第一条消息只能是登录模块的")
			m.Send(misc.ERR_DISCONECT_REASON_FAIL_MUST_LOGIN)
			m.Close()
			return
		}
	}

CREATE_HERO:

	select {
	case <-m.closeChan:
		return
	case <-m.session.ClosedNotify():
		m.Close()
		return
	case msg := <-m.session.MsgChan():
		msgData := msg.(*service.MsgData)
		printC2SLog(msgData)

		switch msgData.ModuleID {
		case login.MODULE_ID:
			switch msgData.SequenceID {
			case login.C2S_CREATE_HERO:
				if m.processCreateHero(msgData) {
					goto WAIT_LOADED
				} else {
					m.Close()
					return
				}
			case login.C2S_SET_TUTORIAL_PROGRESS:
				// 设置新手教程进度
				userMisc := m.user.Misc()
				if userMisc.IsTutorialComplete {
					goto CREATE_HERO
				}

				if proto, ok := msgData.Proto.(*login.C2SSetTutorialProgressProto); ok {
					userMisc.TutorialProgress = proto.Progress
					userMisc.IsTutorialComplete = proto.IsComplete
					if userMisc.IsTutorialComplete {
						m.Send(getNotHeroLoginMsg())
					}
					// tlog
					if tencentInfo := m.user.TencentInfo(); tencentInfo != nil {
						tlogHero := entity.NewSimpleTlogHeroInfo(m.user.Id(), "")
						data := service.TlogService.BuildGuideFlow(tlogHero, tencentInfo, u64.FromInt32(proto.Progress), 1)
						service.TlogService.WriteLog(data)
					}

					if userMisc.IsTutorialComplete {
						gamelogs.NewGuideLog(constants.PID, m.user.Sid(), m.user.Id(), m.user.Misc().TutorialProgress, true)
					}

					goto CREATE_HERO
				} else {
					m.Close()
					return
				}

			default:
				logrus.WithField("moduleID", msgData.ModuleID).WithField("sequenceID", msgData.SequenceID).Error("等待接收login的C2S_CREATE_HERO, 但是收到了别的消息")
				m.Close()
				return
			}
		case misc.MODULE_ID:
			if m.processWaitMiscMsg(msgData) {
				goto CREATE_HERO
			} else {
				logrus.WithField("moduleID", msgData.ModuleID).WithField("sequenceID", msgData.SequenceID).Error("等待接收login的C2S_CREATE_HERO, 但是收到了别的消息")
				m.Close()
				return
			}
		default:
			logrus.WithField("moduleID", msgData.ModuleID).WithField("sequenceID", msgData.SequenceID).Error("等待接收login的C2S_CREATE_HERO, 但是收到了别的消息")
			m.Close()
			return
		}
	}

WAIT_LOADED:
	select {
	case <-m.closeChan:
		return
	case <-m.session.ClosedNotify():
		m.Close()
		return
	case msg := <-m.session.MsgChan():
		msgData := msg.(*service.MsgData)
		printC2SLog(msgData)

		switch msgData.ModuleID {
		case login.MODULE_ID:
			switch msgData.SequenceID {
			case login.C2S_LOADED:
				if m.processLoadedMsg(msgData) {
					// 里面会构建snapshot, 并调用HeroSnapshotService.Online
					goto RECEIVE_MSG
				} else {
					m.Close()
					return
				}
			default:
				logrus.WithField("moduleID", msgData.ModuleID).WithField("sequenceID", msgData.SequenceID).Error("等待接收login的C2S_LOADED, 但是收到了别的消息")
				m.Close()
				return
			}
		case misc.MODULE_ID:
			if m.processWaitMiscMsg(msgData) {
				goto WAIT_LOADED
			} else {
				logrus.WithField("moduleID", msgData.ModuleID).WithField("sequenceID", msgData.SequenceID).Error("等待接收login的C2S_LOADED, 但是收到了别的消息")
				m.Close()
				return
			}
		default:
			logrus.WithField("moduleID", msgData.ModuleID).WithField("sequenceID", msgData.SequenceID).Error("等待接收login的C2S_LOADED, 但是收到了别的消息")
			m.Close()
			return
		}
	}

RECEIVE_MSG:

	defer func(hc iface.HeroController) {
		hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
			hero.Offline(service.TimeService.CurrentTime())
			return true
		})
		heromodule.OnHeroOfflineEvent(hc)
		// 下线时通知snapshot service
		service.HeroSnapshotService.Offline(hc.Id())
	}(m.user.GetHeroController())

	mustLoginTimeout = nil
	//nextSaveTime := autoSaveTimeoutWheel.After(autoSaveInterval)

	// weekly reset
	weeklyTickTime := service.TickerService.GetWeeklyTickTime()

	// daily reset
	dailyTickTime := service.TickerService.GetDailyTickTime()

	// daily zero reset
	dailyZeroTickTime := service.TickerService.GetDailyZeroTickTime()

	// daily mc reset
	dailyMcTickTime := service.TickerService.GetDailyMcTickTime()

	// 轩辕会武每日重置
	xuanyDailyTickTime := service.XuanyuanModule.GetResetTickTime()

	// season reset
	seasonTickTime := service.SeasonService.GetSeasonTickTime()

	// update per seconds
	secondTick := loginTimeoutWheel.After(1 * time.Second)

	// update per minute
	minuteTick := service.TickerService.GetPerMinuteTickTime()

	// randomEvent reset
	//randomEventTick := service.TickerService.GetPer6HourTickTime()

	for {
		select {
		case msg := <-m.session.MsgChan():
			m.processMsg(msg)
		case <-m.session.ClosedNotify():
			m.Close()
			return
			//case <-nextSaveTime:
			//	m.doSave()
			//	nextSaveTime = autoSaveTimeoutWheel.After(autoSaveInterval)
		case <-weeklyTickTime.Tick():
			tickTime := weeklyTickTime.GetTickTime()
			weeklyTickTime = service.TickerService.GetWeeklyTickTime()
			m.tryResetWeekly(tickTime)

		case <-dailyTickTime.Tick():
			tickTime := dailyTickTime.GetTickTime()
			dailyTickTime = service.TickerService.GetDailyTickTime()

			m.tryResetDaily(tickTime)

		case <-dailyZeroTickTime.Tick():
			tickTime := dailyZeroTickTime.GetTickTime()
			dailyZeroTickTime = service.TickerService.GetDailyZeroTickTime()

			m.tryResetDailyZero(tickTime)

		case <-dailyMcTickTime.Tick():
			tickTime := dailyMcTickTime.GetTickTime()
			dailyMcTickTime = service.TickerService.GetDailyMcTickTime()

			m.tryResetDailyMc(tickTime)

		case <-xuanyDailyTickTime.Tick():
			tickTime := xuanyDailyTickTime.GetTickTime()
			xuanyDailyTickTime = service.XuanyuanModule.GetResetTickTime()

			m.tryResetXuanyuan(tickTime)

		case <-seasonTickTime.Tick():
			tickTime := seasonTickTime.GetTickTime()
			seasonTickTime = service.SeasonService.GetSeasonTickTime()

			curSeason := service.SeasonService.SeasonByTime(tickTime)
			m.tryResetSeason(tickTime, curSeason)

		case <-secondTick:
			secondTick = loginTimeoutWheel.After(1 * time.Second)
			m.updatePerSeconds()

		case <-minuteTick.Tick():
			minuteTick = service.TickerService.GetPerMinuteTickTime()
			m.updatePerMinute()

			//case <-randomEventTick.Tick():
			//	tickTime := randomEventTick.GetTickTime()
			//	randomEventTick = service.TickerService.GetPer6HourTickTime()
			//
			//	curSeason := service.SeasonService.SeasonByTime(tickTime)
			//	m.tryResetRandomEvent(curSeason, tickTime)

		case <-m.closeChan:
			return
		}
	}
}

func (m *MessageWorker) processWaitMiscMsg(msgData *service.MsgData) bool {
	switch msgData.SequenceID {
	case misc.C2S_HEART_BEAT:
		return true
	case misc.C2S_PING:
		m.Send(misc.PING_S2C)
		return true
	case misc.C2S_SYNC_TIME:
		if proto, ok := msgData.Proto.(*misc.C2SSyncTimeProto); ok {
			service.MiscModule.SyncTime(proto.ClientTime, m)
		} else {
			logrus.Error("客户端登陆之前请求C2S_SYNC_TIME，proto 错误")
		}
		return true
	case misc.C2S_CONFIG:
		if proto, ok := msgData.Proto.(*misc.C2SConfigProto); ok {
			service.MiscModule.SendConfig(proto, m)
		} else {
			logrus.Error("客户端登陆之前请求C2S_CONFIG，proto 错误")
		}
		return true
	case misc.C2S_CONFIGLUA:
		if proto, ok := msgData.Proto.(*misc.C2SConfigluaProto); ok {
			service.MiscModule.SendLuaConfig(proto, m)
		} else {
			logrus.Error("客户端登陆之前请求C2S_CONFIGLUA，proto 错误")
		}
		return true
	case misc.C2S_CLIENT_LOG:
		if proto, ok := msgData.Proto.(*misc.C2SClientLogProto); ok {
			service.MiscModule.PrintClientLog(m.Id(), "", proto.Level, proto.Text)
		} else {
			logrus.Error("客户端登陆之前请求C2S_CLIENT_LOG，proto 错误")
		}
		return true
	case misc.C2S_CLIENT_VERSION:
		var os, tag string
		if proto, ok := msgData.Proto.(*misc.C2SClientVersionProto); ok {
			os = proto.Os
			tag = proto.T
		}

		if msg := service.ClusterService.GetClientVersionMsg(os, tag); msg != nil {
			m.Send(msg)
		}

		return true
	}
	return false
}

// 处理事件. catch panic
func (m *MessageWorker) processMsg(msg interface{}) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("msg", msg).WithField("stack", string(debug.Stack())).Error("Worker.processEvent recovered from panic")

			metrics.IncPanic()
		}
	}()
	m.lastReceiveClientMsgTime = service.TimeService.CurrentTime()

	msgData := msg.(*service.MsgData)
	printC2SLog(msgData)
	m.handleMsgData(msgData, m.user.GetHeroController())
}

func (m *MessageWorker) handleMsgData(msgData *service.MsgData, hc iface.HeroController) {

	if mtc := service.GameExporter.GetMsgTimeCost(); mtc != nil {
		timer := prometheus.NewTimer(mtc.GetMsgTimeCostObserver(msgData.ModuleID, msgData.SequenceID))
		defer timer.ObserveDuration()
	}

	service.Handle(msgData, m.user.GetHeroController())
}

func printC2SLog(msg *service.MsgData) {
	switch msg.ModuleID {
	case misc.MODULE_ID:
		switch msg.SequenceID {
		case misc.C2S_SYNC_TIME, misc.C2S_HEART_BEAT, misc.C2S_PING:
			return
		}
	}

	if service.IndividualServerConfig.GetIsDebug() {
		logrus.Debugf("收到协议: %d-%d: [%v]", msg.ModuleID, msg.SequenceID, msg.Proto)
	}
}

func printS2CLog(msg []byte) {
	if !service.IndividualServerConfig.GetIsDebug() {
		return
	}

	moduleID, sequenceID, msgJson, err := unmarshal.S2cMsgString(msg)
	if err != nil {
		logrus.WithError(err).Debugf("发送协议: %d-%d，解析出错", moduleID, sequenceID)
		return
	}

	switch moduleID {
	case misc.MODULE_ID:
		switch sequenceID {
		case 9, 2, 4, 11, 77:
			return
		}
	}

	logrus.Debugf("发送协议: %d-%d: [%s]", moduleID, sequenceID, msgJson)
}
