package hebi

import (
	"time"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/hebi"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/timer"
	"github.com/lightpaw/male7/util/event"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/pbutil"
	hebiMsg "github.com/lightpaw/male7/gen/pb/hebi"
	"github.com/lightpaw/logrus"
	"math"
	"github.com/lightpaw/male7/pb/server_proto"
	"context"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/idbytes"
	"runtime/debug"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"sync"
)

func NewHebiManager(dep iface.ServiceDep, mail iface.MailModule, ticker iface.TickerService) *HebiManager {
	m := &HebiManager{
		closeNotify:    make(chan struct{}),
		loopExitNotify: make(chan struct{}),
	}

	m.dep = dep
	m.db = dep.Db()
	m.time = dep.Time()
	m.datas = dep.Datas()
	m.world = dep.World()
	m.mail = mail
	m.ticker = ticker

	m.rooms = make([]*HebiRoom, dep.Datas().HebiMiscData().RoomsMaxSize)
	m.heroRecordMsgs = &sync.Map{}
	m.init()

	m.queue = event.NewEventQueue(2048, 5*time.Second, "HebiEvent")

	go call.CatchLoopPanic(m.loop, "HebiManager.Loop")

	m.msgVersion = atomic.NewUint64(1)

	ctime := dep.Time().CurrentTime()
	m.UpdateMsg(ctime)

	return m
}

type HebiManager struct {
	dep    iface.ServiceDep
	db     iface.DbService
	time   iface.TimeService
	datas  iface.ConfigDatas
	world  iface.WorldService
	ticker iface.TickerService
	mail   iface.MailModule

	rooms []*HebiRoom

	queue *event.EventQueue

	msgVersion *atomic.Uint64
	emptyMsg   pbutil.Buffer
	msg        pbutil.Buffer

	heroRecords         map[int64]*shared_proto.HebiHeroRecordProto
	heroRecordMsgs      *sync.Map
	lastResetRecordTime time.Time

	closeNotify    chan struct{}
	loopExitNotify chan struct{}
}

var (
	emptyFunc = func() {}
)

func (m *HebiManager) loop() {
	defer close(m.loopExitNotify)

	dailyTickTime := m.ticker.GetDailyTickTime()
	m.resetDaily(dailyTickTime)

	loopWheel := timer.NewTimingWheel(500*time.Millisecond, 32)
	secondTick := loopWheel.After(time.Second)

	for {
		select {
		case <-secondTick:
			secondTick = loopWheel.After(time.Second)
			m.Func("hebi.loop", m.update, emptyFunc)
		case <-dailyTickTime.Tick():
			m.resetDaily(dailyTickTime)
			dailyTickTime = m.ticker.GetDailyTickTime()
		case <-m.closeNotify:
			return
		}
	}
}

func (m *HebiManager) resetDaily(tickTime tickdata.TickTime) {
	ctime := m.time.CurrentTime()
	if m.lastResetRecordTime.Before(tickTime.GetPrevTickTime()) {
		m.lastResetRecordTime = ctime
		m.heroRecords = make(map[int64]*shared_proto.HebiHeroRecordProto)
	}
}

func (m *HebiManager) init() {
	var bytes []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		bytes, err = m.db.LoadKey(ctx, server_proto.Key_NewHebi)
		return
	})

	if err != nil {
		logrus.WithError(err).Panic("加载天命合璧模块数据失败")
	}

	if len(bytes) <= 0 {
		return
	}

	p := &server_proto.HebiServerProto{}
	if err := p.Unmarshal(bytes); err != nil {
		logrus.WithError(err).Panic("解析天命合璧模块数据失败")
	}
	m.heroRecords = p.Records
	if m.heroRecords == nil {
		m.heroRecords = make(map[int64]*shared_proto.HebiHeroRecordProto)
	}
	for heroId, record := range m.heroRecords {
		m.updateHeroRecordMsg(heroId, record)
	}

	m.lastResetRecordTime = timeutil.Unix64(p.LastResetRecordTime)

	m.unmarshalRooms(p.Info, m.dep.Datas())
}

func (m *HebiManager) close() {
	close(m.closeNotify)
	<-m.loopExitNotify

	m.queue.Stop()

	if err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		return m.db.SaveKey(ctx, server_proto.Key_NewHebi, must.Marshal(m.encodeServer()))
	}); err != nil {
		logrus.WithError(err).Error("保存天命合璧模块数据失败")
	}
}

func (m *HebiManager) Func(handlerName string, f, sendErrMsg func()) (succ bool) {
	return m.doFunc(handlerName, f, sendErrMsg, true)
}

func (m *HebiManager) FuncNoWait(handlerName string, f, sendErrMsg func()) (succ bool) {
	return m.doFunc(handlerName, f, sendErrMsg, false)
}

func (m *HebiManager) doFunc(handlerName string, f, sendErrMsg func(), wait bool) (succ bool) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("stack", string(debug.Stack())).WithField("err", r).Errorf("%s recovered from panic. SEVERE!!!", handlerName)
			metrics.IncPanic()
			sendErrMsg()
		}
	}()

	ok := m.queue.TimeoutFunc(wait, f)
	if !ok {
		sendErrMsg()
	} else {
		succ = true
	}

	return
}

func (m *HebiManager) update() {
	var changed bool
	ctime := m.time.CurrentTime()
	for i, room := range m.rooms {
		if room == nil {
			continue
		}

		switch room.state {
		case shared_proto.HebiRoomState_HebiRoomEmpty:
			m.rooms[i] = nil
		case shared_proto.HebiRoomState_HebiRoomWait:
			if !timeutil.IsZero(room.hebiWaitExpiredTime) &&
				ctime.After(room.hebiWaitExpiredTime) {
				// 一直没有合璧，踢出去
				hostId := room.hostId
				ok, hostGoodsData := room.Leave()
				if ok {
					hctx := heromodule.NewContext(m.dep, operate_type.HebiKick)
					changed = true
					m.dep.World().Send(hostId, hebiMsg.NewS2cLeaveRoomMsg(u64.Int32(room.roomId)))
					// 物品补回去
					m.dep.HeroData().FuncWithSend(hostId, func(hero *entity.Hero, result herolock.LockResult) {
						heromodule.AddGoods(hctx, hero, result, hostGoodsData, 1)
						result.Changed()
						result.Ok()
					})
				}
			}
		case shared_proto.HebiRoomState_HebiRoomRobProtect:
			if ctime.After(room.hebiRobProtectEndTime) {
				// 超过保护时间还没进，房间准备删掉
				room.state = shared_proto.HebiRoomState_HebiRoomEmpty
				changed = true
			}
		case shared_proto.HebiRoomState_HebiRoomHebiRunning:
			if ctime.After(room.hebiCompleteTime) {
				hostId, guestId := room.hostId, room.guestId
				hostGoods, guestGoods := room.hostGoodsData, room.guestGoodsData
				copySelf := room.copySelf
				succ, prizeId := room.complete()
				if succ {
					m.GiveCompletePrize(hostId, guestId, hostGoods, guestGoods, prizeId, ctime, copySelf)
				} else {
					logrus.Debugf("合璧完成，但是 succ == false")
				}

				m.dep.World().Send(hostId, hebiMsg.NewS2cCompleteMsg(u64.Int32(room.roomId)))
				m.dep.World().Send(guestId, hebiMsg.NewS2cCompleteMsg(u64.Int32(room.roomId)))

				changed = true
			}
		}
	}
	if changed {
		m.UpdateMsg(ctime)
	}
}

func (m *HebiManager) UpdateMsg(ctime time.Time) {
	m.msgVersion.Inc()
	m.msg = hebiMsg.NewS2cRoomListMsg(u64.Int32(m.msgVersion.Load()), m.encode()).Static()
	m.emptyMsg = hebiMsg.NewS2cRoomListMsg(u64.Int32(m.msgVersion.Load()), nil).Static()
}

func (m *HebiManager) encode() *shared_proto.HebiInfoProto {
	p := &shared_proto.HebiInfoProto{}
	for _, r := range m.rooms {
		if r != nil {
			p.Room = append(p.Room, r.encode())
		}
	}

	return p
}

func (m *HebiManager) encodeServer() *server_proto.HebiServerProto {
	p := &server_proto.HebiServerProto{}
	p.Records = m.heroRecords
	p.LastResetRecordTime = timeutil.Marshal64(m.lastResetRecordTime)

	info := &shared_proto.HebiInfoProto{}
	for _, r := range m.rooms {
		if r != nil {
			info.Room = append(info.Room, r.encode())
		}
	}
	p.Info = info

	return p
}

func (m *HebiManager) unmarshalRooms(p *shared_proto.HebiInfoProto, datas iface.ConfigDatas) {
	for _, rp := range p.Room {
		r := &HebiRoom{}
		r.unmarshal(rp, datas)
		m.rooms[r.roomId] = r
	}
}

func (m *HebiManager) GiveCompletePrize(hostId, guestId int64, hostGoods, guestGoods *goods.GoodsData,
	prizeId uint64, ctime time.Time, copySelf bool) {
	m.updateGuildTaskProgress(hostId, guestId, m.datas.GuildConfig().HebiCompleteTaskProgress)

	prize := m.datas.GetHebiPrizeData(prizeId)
	if prize == nil {
		logrus.Debugf("合璧完成，但找不到奖励 %v", prizeId)
		return
	}

	if !copySelf {
		// 恭喜您与[color=#6Baaf0]{{name}}[/color]成功合成[color=#e2e2e2]【{{goods}}】[/color]，获得了丰富奖励。
		if mailData := m.datas.MailHelp().HebiCompletePrize; mailData != nil {
			mailProto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
			mailProto.Prize = resdata.AppendPrize(prize.PlunderPrize.GetPrize(), prize.AmountPrize).Encode()
			// 发给自己
			mailProto.Text = mailData.NewTextFields().WithName(m.dep.HeroSnapshot().GetFlagHeroName(guestId)).WithGoods(goods.GetGoodsName(hostGoods)).JsonString()
			m.mail.SendProtoMail(hostId, mailProto, ctime)
			// 发给对方
			mailProto.Text = mailData.NewTextFields().WithName(m.dep.HeroSnapshot().GetFlagHeroName(hostId)).WithGoods(goods.GetGoodsName(guestGoods)).JsonString()
			m.mail.SendProtoMail(guestId, mailProto, ctime)
		}
	} else {
		if mailData := m.datas.MailHelp().HebiCopyCompletePrize; mailData != nil {
			mailProto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
			mailProto.Prize = resdata.AppendPrize(prize.PlunderPrize.GetPrize(), prize.AmountPrize).Encode()
			mailProto.Text = mailData.NewTextFields().WithGoods(goods.GetGoodsName(hostGoods)).JsonString()
			m.mail.SendProtoMail(hostId, mailProto, ctime)
		}
	}
}

func (m *HebiManager) updateGuildTaskProgress(hostId, guestId int64, progress uint64) {
	data := m.datas.GetGuildTaskData(u64.FromInt32(int32(server_proto.GuildTaskType_HeBi)))
	if host := m.dep.HeroSnapshot().Get(hostId); host != nil {
		m.dep.Guild().AddGuildTaskProgress(host.GuildId, data, progress)
	}
	if guestId > 0 && guestId != hostId {
		if guest := m.dep.HeroSnapshot().Get(guestId); guest != nil {
			m.dep.Guild().AddGuildTaskProgress(guest.GuildId, data, progress)
		}
	}
}

func calcRobPrizeCode(attackerFightAmount, defenderFightAmount uint64) float64 {
	// 抢夺系数=max(min((被抢者战力/抢夺者战力)^3,1.1),0.35)
	return math.Max(math.Min(math.Pow(
		u64.Division2Float64(defenderFightAmount, attackerFightAmount), 3),
		1.1), 0.35)
}

func (m *HebiManager) GetRobPrize(prizeId uint64, ctime time.Time, attackerFightAmount, defenderFightAmount uint64) *resdata.Prize {
	prize := m.datas.GetHebiPrizeData(prizeId)
	if prize == nil {
		logrus.Debugf("合璧抢夺成功，但找不到奖励")
		return nil
	}

	code := calcRobPrizeCode(attackerFightAmount, defenderFightAmount)

	// int(7/5*抢夺系数*合璧奖励)
	robMulti := math.Min(1, math.Max(0, (7*code)/5))

	return prize.AmountPrize.MultiCoef(robMulti)
}

func (m *HebiManager) GiveBeRobbedPrize(hostId, guestId int64, hostGoods, guestGoods *goods.GoodsData,
	prizeId uint64, attackerFlagName string, attackerFightAmount, defenderFightAmount uint64, ctime time.Time, copySelf bool) {
	m.updateGuildTaskProgress(hostId, guestId, m.datas.GuildConfig().HebiCompleteTaskProgress)

	prize := m.datas.GetHebiPrizeData(prizeId)
	if prize == nil {
		logrus.Debugf("合璧抢夺成功，但找不到奖励")
		return
	}

	code := calcRobPrizeCode(attackerFightAmount, defenderFightAmount)

	// int( (1-7/10*抢夺系数)*合璧奖励)
	beRobbedMulti := math.Min(1, math.Max(0, 1-(code*7/10)))
	if !copySelf {
		// 您在与[color=#6Baaf0]{{name}}[/color]合成[color=#e2e2e2]【{{goods}}】[/color]的过程中，遭到[color=#ff304e]{{attacker}}[/color]干扰，失去了部分奖励。
		if mailData := m.datas.MailHelp().HebiBeRobbedPrize; mailData != nil {
			mailProto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
			realPrize := resdata.AppendPrize(prize.PlunderPrize.GetPrize(), prize.AmountPrize.MultiCoef(beRobbedMulti))
			mailProto.Prize = realPrize.Encode()
			// 发给自己
			mailProto.Text = mailData.NewTextFields().WithName(m.dep.HeroSnapshot().GetFlagHeroName(guestId)).WithGoods(goods.GetGoodsName(hostGoods)).WithAttacker(attackerFlagName).JsonString()
			m.mail.SendProtoMail(hostId, mailProto, ctime)
			// 发给对方
			mailProto.Text = mailData.NewTextFields().WithName(m.dep.HeroSnapshot().GetFlagHeroName(hostId)).WithGoods(goods.GetGoodsName(guestGoods)).WithAttacker(attackerFlagName).JsonString()
			m.mail.SendProtoMail(guestId, mailProto, ctime)
		}
	} else {
		if mailData := m.datas.MailHelp().HebiCopyBeRobbedPrize; mailData != nil {
			mailProto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
			realPrize := resdata.AppendPrize(prize.PlunderPrize.GetPrize(), prize.AmountPrize.MultiCoef(beRobbedMulti))
			mailProto.Prize = realPrize.Encode()

			mailProto.Text = mailData.NewTextFields().WithGoods(goods.GetGoodsName(hostGoods)).WithAttacker(attackerFlagName).JsonString()
			m.mail.SendProtoMail(hostId, mailProto, ctime)
		}
	}
}

func (m *HebiManager) GetRoom(roomId uint64) *HebiRoom {
	if roomId >= m.datas.HebiMiscData().RoomsMaxSize {
		return nil
	}
	return m.rooms[roomId]
}

func (m *HebiManager) CheckInHebiRoom(roomId uint64, heroId int64, hero *shared_proto.HeroBasicProto, captain *shared_proto.HebiCaptainProto, yubi *goods.GoodsData, ctime time.Time) *HebiRoom {
	if roomId >= m.datas.HebiMiscData().RoomsMaxSize {
		return nil
	}

	if r := m.rooms[roomId]; r != nil {
		if r.state != shared_proto.HebiRoomState_HebiRoomEmpty && r.state != shared_proto.HebiRoomState_HebiRoomRobProtect {
			return nil
		}
	}

	r := &HebiRoom{}
	r.miscData = m.datas.HebiMiscData()

	r.roomId = roomId
	r.hostId = heroId
	r.host = hero
	r.hostCaptain = captain
	r.hostGoodsData = yubi
	r.hebiWaitExpiredTime = ctime.Add(m.datas.HebiMiscData().RoomWaitExpiredDuration)

	r.state = shared_proto.HebiRoomState_HebiRoomWait
	m.rooms[roomId] = r

	return r
}

func (m *HebiManager) InAnyRoom(heroId int64) bool {
	// 是否已经在房间里
	_, exist := m.HeroRoom(heroId)
	return exist
}

func (m *HebiManager) HeroRoom(heroId int64) (*HebiRoom, bool) {
	for _, r := range m.rooms {
		if r == nil {
			continue
		}

		if r.host != nil {
			if r.hostId == heroId {
				return r, true
			}
		}

		if r.guest != nil {
			if r.guestId == heroId {
				return r, true
			}
		}
	}
	return nil, false
}

func (m *HebiManager) updateGuildInfo(heroId int64, g *guildsnapshotdata.GuildSnapshot) bool {
	ctime := m.time.CurrentTime()

	// 是否已经在房间里
	for _, r := range m.rooms {
		if r == nil {
			continue
		}

		if r.host != nil {
			if r.hostId == heroId {
				updateHeroGuild(r.host, g)
				m.UpdateMsg(ctime)
				return true
			}
		}

		if r.guest != nil {
			if r.guestId == heroId {
				updateHeroGuild(r.guest, g)
				m.UpdateMsg(ctime)
				return true
			}
		}
	}

	return false
}

func updateHeroGuild(hero *shared_proto.HeroBasicProto, g *guildsnapshotdata.GuildSnapshot) {
	if g == nil {
		hero.GuildId = 0
		hero.GuildFlagName = ""
		hero.GuildName = ""
	} else {
		hero.GuildId = int32(g.Id)
		hero.GuildFlagName = g.FlagName
		hero.GuildName = g.Name
	}
}

type HebiRoom struct {
	roomId uint64
	state  shared_proto.HebiRoomState

	hostId        int64
	host          *shared_proto.HeroBasicProto
	hostCaptain   *shared_proto.HebiCaptainProto
	hostGoodsData *goods.GoodsData

	guestId        int64
	guest          *shared_proto.HeroBasicProto
	guestCaptain   *shared_proto.HebiCaptainProto
	guestGoodsData *goods.GoodsData

	prizeId uint64

	hebiWaitExpiredTime   time.Time
	hebiCompleteTime      time.Time
	hebiRobProtectEndTime time.Time

	hebiRobProtectId int64

	copySelf bool

	miscData *hebi.HebiMiscData
}

func (r *HebiRoom) GetHostId() int64 {
	return r.hostId
}

func (r *HebiRoom) GetHostGoods() *goods.GoodsData {
	return r.hostGoodsData
}

func (r *HebiRoom) GetGuestId() int64 {
	return r.guestId
}

func (r *HebiRoom) GetGuestGoods() *goods.GoodsData {
	return r.guestGoodsData
}

func (r *HebiRoom) encode() *shared_proto.HebiRoomProto {
	p := &shared_proto.HebiRoomProto{}
	p.RoomId = u64.Int32(r.roomId)
	p.State = r.state

	p.Host = r.host
	p.HostCaptain = r.hostCaptain
	if r.hostGoodsData != nil {
		p.HostGoodsId = u64.Int32(r.hostGoodsData.Id)
	}
	p.Guest = r.guest
	p.GuestCatpain = r.guestCaptain
	if r.guestGoodsData != nil {
		p.GuestGoodsId = u64.Int32(r.guestGoodsData.Id)
	}

	p.CopySelf = r.copySelf
	p.HebiWaitExpiredTime = timeutil.Marshal32(r.hebiWaitExpiredTime)
	p.HebiCompleteTime = timeutil.Marshal32(r.hebiCompleteTime)
	p.HebiRobProtectedTime = timeutil.Marshal32(r.hebiRobProtectEndTime)

	p.PrizeId = u64.Int32(r.prizeId)

	return p
}

func (r *HebiRoom) unmarshal(p *shared_proto.HebiRoomProto, datas iface.ConfigDatas) {
	r.miscData = datas.HebiMiscData()

	r.roomId = u64.FromInt32(p.RoomId)
	r.state = p.State
	r.host = p.Host
	if r.host != nil {
		r.hostId, _ = idbytes.ToId(r.host.Id)
	}
	r.hostCaptain = p.HostCaptain
	r.hostGoodsData = datas.GetGoodsData(u64.FromInt32(p.HostGoodsId))
	r.guest = p.Guest
	if r.guest != nil {
		r.guestId, _ = idbytes.ToId(r.guest.Id)
	}
	r.guestCaptain = p.GuestCatpain
	r.guestGoodsData = datas.GetGoodsData(u64.FromInt32(p.GuestGoodsId))

	r.copySelf = p.CopySelf
	r.hebiWaitExpiredTime = timeutil.Unix32(p.HebiWaitExpiredTime)
	r.hebiCompleteTime = timeutil.Unix32(p.HebiCompleteTime)
	r.hebiRobProtectEndTime = timeutil.Unix32(p.HebiRobProtectedTime)

	r.prizeId = u64.FromInt32(p.PrizeId)
}

func (r *HebiRoom) HostPos() shared_proto.HebiSubType {
	if r.hostGoodsData == nil {
		return shared_proto.HebiSubType_HebiInvalidSubType
	}
	return r.hostGoodsData.HebiSubType
}

func (r *HebiRoom) ChangeCaptain(heroId int64, c *shared_proto.HebiCaptainProto) bool {
	if r.state != shared_proto.HebiRoomState_HebiRoomWait {
		return false
	}

	if r.hostId == heroId {
		r.hostCaptain = c
		return true
	}
	if r.guestId == heroId {
		r.guestCaptain = c
		return true
	}
	return false
}

func (r *HebiRoom) Join(heroId int64, hero *shared_proto.HeroBasicProto, captain *shared_proto.HebiCaptainProto, yubi *goods.GoodsData, ctime time.Time) bool {
	if r.state != shared_proto.HebiRoomState_HebiRoomWait {
		return false
	}
	if r.hostGoodsData.HebiSubType == yubi.HebiSubType {
		return false
	}

	r.guestId = heroId
	r.guest = hero
	r.guestCaptain = captain
	r.guestGoodsData = yubi

	return r.startHebi(ctime)
}

func (r *HebiRoom) CopySelf(ctime time.Time, captain *shared_proto.HebiCaptainProto, goodsData *goods.GoodsData) bool {
	if r.state != shared_proto.HebiRoomState_HebiRoomWait {
		return false
	}
	r.copySelf = true
	r.guestId = r.hostId
	r.guest = r.host
	r.guestGoodsData = goodsData
	if captain != nil {
		r.guestCaptain = captain
		r.hostCaptain = captain
	} else {
		r.guestCaptain = r.hostCaptain
	}

	return r.startHebi(ctime)
}

func (r *HebiRoom) startHebi(ctime time.Time) bool {
	if r.state != shared_proto.HebiRoomState_HebiRoomWait {
		return false
	}
	r.hebiCompleteTime = ctime.Add(r.miscData.HebiDuration)

	r.state = shared_proto.HebiRoomState_HebiRoomHebiRunning

	r.prizeId = hebi.GenHebiPrizeId(u64.Max(u64.FromInt32(r.host.Level), u64.FromInt32(r.guest.Level)), r.hostGoodsData.HebiType, u64.Max(r.hostGoodsData.GoodsQuality.Level, r.guestGoodsData.GoodsQuality.Level))

	return true
}

func (r *HebiRoom) complete() (succ bool, prizeId uint64) {
	if r.state != shared_proto.HebiRoomState_HebiRoomHebiRunning {
		return
	}

	prizeId = r.prizeId
	r.reset()
	r.state = shared_proto.HebiRoomState_HebiRoomEmpty
	succ = true
	return
}

// 只有房主一个人在时才能离开，客人加入后会直接进行合璧
func (r *HebiRoom) Leave() (succ bool, yubi *goods.GoodsData) {
	if r.state != shared_proto.HebiRoomState_HebiRoomWait && r.state != shared_proto.HebiRoomState_HebiRoomRobProtect {
		return
	}

	yubi = r.hostGoodsData
	r.reset()
	r.state = shared_proto.HebiRoomState_HebiRoomEmpty
	succ = true

	return
}

func (r *HebiRoom) RobPos(heroId int64, ctime time.Time) (yubi *goods.GoodsData, succ bool) {
	if r.state != shared_proto.HebiRoomState_HebiRoomWait {
		return
	}

	yubi = r.hostGoodsData

	r.reset()
	r.hebiRobProtectId = heroId
	r.hebiRobProtectEndTime = ctime.Add(r.miscData.RobProtectDuration)
	r.state = shared_proto.HebiRoomState_HebiRoomRobProtect

	succ = true
	return
}

func (r *HebiRoom) Rob() (succ bool, prizeId uint64) {
	if r.state != shared_proto.HebiRoomState_HebiRoomHebiRunning {
		return
	}

	prizeId = r.prizeId
	r.reset()
	r.state = shared_proto.HebiRoomState_HebiRoomEmpty
	succ = true
	return
}

func (r *HebiRoom) resetHost() {
	r.hostId = 0
	r.host = nil
	r.hostGoodsData = nil
	r.hostCaptain = nil
}

func (r *HebiRoom) resetGuest() {
	r.guestId = 0
	r.guest = nil
	r.guestGoodsData = nil
	r.guestCaptain = nil
}

func (r *HebiRoom) reset() {
	r.resetHost()
	r.resetGuest()

	r.prizeId = 0
	r.hebiWaitExpiredTime = time.Time{}
	r.hebiCompleteTime = time.Time{}
	r.hebiRobProtectEndTime = time.Time{}
	r.hebiRobProtectId = 0
	r.copySelf = false
}

func (m *HebiManager) AddHeroRecord(selfId, partnerId int64, self, target *shared_proto.HeroBasicProto, selfCaptain, targetCaptain *shared_proto.HebiCaptainProto, isRob, isAtk bool, hebiType shared_proto.HebiType, fightNum uint64, resp *server_proto.CombatResponseServerProto, ctime time.Time) {
	p := &shared_proto.HebiHeroSingleRecordProto{}
	p.Time = timeutil.Marshal32(ctime)
	p.Self = self
	p.Target = target
	p.SelfCaptain = selfCaptain
	p.TargetCaptain = targetCaptain
	p.IsRob = isRob
	p.HebiType = hebiType
	p.FightNum = u64.Int32(fightNum)
	p.IsAtk = isAtk
	if isAtk {
		p.Win = resp.AttackerWin
	} else {
		p.Win = !resp.AttackerWin
	}

	p.Combat = &shared_proto.CombatShareProto{}
	p.Combat.Link = resp.Link
	p.Combat.Type = shared_proto.CombatType_SINGLE
	p.Combat.IsAttacker = isAtk

	m.addHeroRecord0(selfId, p)
	if partnerId > 0 && partnerId != selfId {
		m.addHeroRecord0(partnerId, p)
	}
}

func (m *HebiManager) addHeroRecord0(selfId int64, p *shared_proto.HebiHeroSingleRecordProto) {
	records, ok := m.heroRecords[selfId]
	if !ok || records == nil {
		records = &shared_proto.HebiHeroRecordProto{}
		m.heroRecords[selfId] = records
	}

	if u64.FromInt(len(records.Record)) >= m.dep.Datas().HebiMiscData().HebiHeroRecordMaxSize {
		records.Record = records.Record[1:]
	}

	records.Record = append(records.Record, p)

	m.updateHeroRecordMsg(selfId, records)
}

func (m *HebiManager) updateHeroRecordMsg(heroId int64, p *shared_proto.HebiHeroRecordProto) {
	if p == nil {
		p = &shared_proto.HebiHeroRecordProto{}
	}
	m.heroRecordMsgs.Store(heroId, hebiMsg.NewS2cHeroRecordListMsg(p))
}

var EmptyHeroRecordMsg = hebiMsg.NewS2cHeroRecordListMsg(&shared_proto.HebiHeroRecordProto{})

func (m *HebiManager) GetHeroRecordMsg(heroId int64) (msg pbutil.Buffer) {
	if value, ok := m.heroRecordMsgs.Load(heroId); !ok || value == nil {
		msg = EmptyHeroRecordMsg
	} else if msg, ok = value.(pbutil.Buffer); !ok {
		msg = EmptyHeroRecordMsg
	}
	return
}
