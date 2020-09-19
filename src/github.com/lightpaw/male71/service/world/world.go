package world

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/pbutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/gamelogs"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/util/call"
)

//gogen:iface
type WorldService struct {
	serverConfig iface.IndividualServerConfig
	ticker       iface.TickerService

	connectedUsers *usermap
	userCounter    *atomic.Uint64
}

func NewWorldService(register iface.MetricsRegister, serverConfig iface.IndividualServerConfig, ticker iface.TickerService) *WorldService {

	userMap := Newusermap()
	userCounter := atomic.NewUint64(0)

	if register.EnableOnlineCountMetrics() {
		onlineCounter := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: metrics.Namespace,
			Subsystem: "online_user",
			Name:      "count",
			Help:      "Online user count.",
		})

		register.Register(onlineCounter)
		register.RegisterFunc(func() {
			onlineCounter.Set(float64(userCounter.Load()))
		})
	}

	w := &WorldService{
		serverConfig:   serverConfig,
		ticker:         ticker,
		connectedUsers: userMap,
		userCounter:    userCounter,
	}

	go call.CatchLoopPanic(w.loop, "world.loop()")

	return w
}

func (w *WorldService) loop() {

	minuteTick := w.ticker.GetPerMinuteTickTime()

	for {
		select {
		case <-minuteTick.Tick():
			minuteTick = w.ticker.GetPerMinuteTickTime()

			// 定时记录在线数据
			gamelogs.ServerOnlineLog(constants.PID, uint32(w.serverConfig.GetServerID()), w.userCounter.Load())
		}
	}

}

func (w *WorldService) Close() {
	logrus.Info("WorldService开始保存所有玩家")
	entries := w.connectedUsers.Iter()
	for ch := range entries {
		for _, user := range ch {
			user.Val.DisconnectAndWait(misc.ErrDisconectReasonFailClose)
		}
	}

	if count := w.connectedUsers.Count(); count != 0 {
		logrus.WithField("count", count).Error("WorldService.Close了之后, 里面竟然还有人...")
	}
}

// 尝试放入用户, 如果用户已在里面, 则返回旧的用户和false, 如果用户放入成功, 则返回nil, ok
func (w *WorldService) PutConnectedUserIfAbsent(c iface.ConnectedUser) (iface.ConnectedUser, bool) {
	u, ok := w.connectedUsers.SetIfAbsent(c.Id(), c)
	if ok {
		w.userCounter.Inc()
	}
	return u, ok
}

// 删除用户, 如果是同一个对象的话
func (w *WorldService) RemoveUserIfSame(c iface.ConnectedUser) bool {
	ok := w.connectedUsers.RemoveIfSame(c.Id(), c)
	if ok {
		w.userCounter.Dec()
	}
	return ok
}

func (w *WorldService) IsOnline(id int64) bool {
	u, _ := w.connectedUsers.Get(id)
	return u != nil
}

func (w *WorldService) IsDontPush(id int64) bool {
	u, _ := w.connectedUsers.Get(id)
	if u == nil {
		return false
	}

	if hc := u.GetHeroController(); hc == nil {
		// 还没建号，推给他（这种情况应该很少发生）
		return false
	} else {
		return !hc.GetIsInBackgroud()
	}
}

func (w *WorldService) GetUserSender(id int64) sender.Sender {
	u, _ := w.connectedUsers.Get(id)
	return u
}

func (w *WorldService) GetUserCloseSender(id int64) sender.ClosableSender {
	u, _ := w.connectedUsers.Get(id)
	return u
}

func (w *WorldService) GetTencentInfo(id int64) *shared_proto.TencentInfoProto {
	if u, ok := w.connectedUsers.Get(id); ok {
		return u.TencentInfo()
	}
	return nil
}

func (w *WorldService) SendFunc(id int64, msgFunc iface.MsgFunc) {
	if isValidHeroId(id) && msgFunc != nil {
		u, _ := w.connectedUsers.Get(id)
		if u != nil && u.IsLoaded() {
			if msg := msgFunc(); msg != nil {
				u.Send(msg)
			}
		}
	}
}

func (w *WorldService) Send(id int64, msg pbutil.Buffer) {
	if isValidHeroId(id) && msg != nil {
		u, _ := w.connectedUsers.Get(id)
		if u != nil && u.IsLoaded() {
			u.Send(msg)
		}
	}
}

func (w *WorldService) SendMsgs(id int64, msgs []pbutil.Buffer) {

	if isValidHeroId(id) && len(msgs) > 0 {
		u, _ := w.connectedUsers.Get(id)
		if u != nil && u.IsLoaded() {
			for _, msg := range msgs {
				if msg != nil {
					u.Send(msg)
				}
			}
		}
	}
}

func isValidHeroId(id int64) bool {
	return id != 0 && !npcid.IsNpcId(id)
}

func (w *WorldService) MultiSend(ids []int64, msg pbutil.Buffer) {
	w.MultiSendIgnore(ids, msg, 0)
}

func (w *WorldService) MultiSendIgnore(ids []int64, msg pbutil.Buffer, dontSend int64) {
	if msg == nil {
		return
	}

	msg = msg.Static()
	for _, id := range ids {
		if id != dontSend {
			u, _ := w.connectedUsers.Get(id)
			if u != nil && u.IsLoaded() {
				u.Send(msg)
			}
		}
	}
}

func (w *WorldService) MultiSendMsgs(ids []int64, msgs []pbutil.Buffer) {
	w.MultiSendMsgsIgnore(ids, msgs, 0)
}

func (w *WorldService) MultiSendMsgsIgnore(ids []int64, msgs []pbutil.Buffer, dontSend int64) {
	if len(msgs) <= 0 {
		return
	}

	// 转换一下
	for i, msg := range msgs {
		msgs[i] = msg.Static()
	}

	for _, id := range ids {
		if id != dontSend {
			u, _ := w.connectedUsers.Get(id)
			if u != nil && u.IsLoaded() {
				for _, msg := range msgs {
					u.Send(msg)
				}
			}
		}
	}
}

// 广播消息. 同一个[]byte会发送给每一个人. 调用了之后不可以再修改[]byte中的内容
// 消耗大, 会一个个shard锁住user map
func (w *WorldService) Broadcast(msg pbutil.Buffer) {
	msg = msg.Static()
	w.connectedUsers.IterCb(func(uid int64, user iface.ConnectedUser) {
		if user.IsLoaded() {
			user.Send(msg)
		}
	})
}

func (w *WorldService) BroadcastIgnore(msg pbutil.Buffer, dontSend int64) {
	msg = msg.Static()
	w.connectedUsers.IterCb(func(uid int64, user iface.ConnectedUser) {
		if uid != dontSend && user.IsLoaded() {
			user.Send(msg)
		}
	})
}

func (w *WorldService) WalkHero(walk iface.HeroWalker) {
	w.connectedUsers.IterCb(func(uid int64, user iface.ConnectedUser) {
		hero := user.GetHeroController()
		if hero != nil {
			walk(uid, hero)
		}
	})
}

func (w *WorldService) FuncHero(id int64, walk iface.HeroWalker) bool {

	if isValidHeroId(id) && walk != nil {
		u, _ := w.connectedUsers.Get(id)
		if u != nil {
			hero := u.GetHeroController()
			if hero != nil {
				walk(id, hero)
				return true
			}
		}
	}

	return false
}

func (w *WorldService) WalkUser(walk iface.UserWalker) {
	w.connectedUsers.IterCb(IterCbusermap(walk))
}
