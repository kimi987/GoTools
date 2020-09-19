package mail

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/mail"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/util/event"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/service/operate_type"
	"strconv"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/util/call"
)

//gogen:iface
type MailModule struct {
	time          iface.TimeService
	world         iface.WorldService
	broadcast     iface.BroadcastService
	datas         iface.ConfigDatas
	tlog          iface.TlogService
	heroData      iface.HeroDataService
	guildSnapshot iface.GuildSnapshotService
	db            iface.DbService
	chat          iface.ChatService

	queue *event.FuncQueue

	idGen *atomic.Int64
}

func NewMailModule(
	time iface.TimeService,
	world iface.WorldService,
	broadcast iface.BroadcastService,
	datas iface.ConfigDatas,
	tlog iface.TlogService,
	heroData iface.HeroDataService,
	guildSnapshot iface.GuildSnapshotService,
	db iface.DbService,
	chat iface.ChatService) *MailModule {
	m := &MailModule{}
	m.time = time
	m.world = world
	m.db = db
	m.chat = chat
	m.broadcast = broadcast
	m.datas = datas
	m.tlog = tlog
	m.heroData = heroData
	m.guildSnapshot = guildSnapshot
	m.queue = event.NewFuncQueue(1024, "mail")

	var maxId uint64
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		maxId, err = db.MaxMailId(ctx)
		return
	})
	if err != nil {
		logrus.WithError(err).Panicf("获取邮件MaxID 失败")
	}

	m.idGen = atomic.NewInt64(int64(maxId))

	heromodule.RegisterHeroOnlineListener(m)

	return m
}

func (m *MailModule) OnHeroOnline(hc iface.HeroController) {
	m.notifyMailCountWithSender(hc)
}

func (m *MailModule) newMailId() uint64 {
	return uint64(m.idGen.Inc())
}

// 过期函数@AlbertFan
// 发邮件
func (m *MailModule) SendMail(target int64, icon uint64, title, text string, keep bool,
	report *shared_proto.FightReportProto, prize *shared_proto.PrizeProto, sendTime time.Time) bool {

	proto := &shared_proto.MailProto{}
	proto.Icon = u64.Int32(icon)
	proto.Title = title
	proto.Text = text
	proto.Keep = keep
	proto.Report = report
	proto.Prize = prize

	return m.SendProtoMail(target, proto, sendTime)
}

func (m *MailModule) SendReportMail(target int64, proto *shared_proto.MailProto, sendTime time.Time) bool {
	if target == 0 || npcid.IsNpcId(target) {
		return false
	}
	time := sendTime
	mailId := m.newMailId()
	proto.Id = i64.ToBytesU64(mailId)
	proto.SendTime = timeutil.Marshal32(time)
	proto.HasPrize = proto.Prize != nil
	proto.HasReport = true

	data, err := proto.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("发送Report邮件，proto marshal报错")
		return false
	}

	keep := proto.Keep
	hasReport := proto.HasReport
	hasPrize := proto.Prize != nil
	tag := proto.ReportTag

	m.queue.MustFunc(func() {
		err = ctxfunc.Timeout3s(func(ctx context.Context) error {
			return m.db.CreateMail(ctx, mailId, target, data, keep, hasReport, hasPrize, tag, timeutil.Marshal64(time))
		})
		if err != nil {
			logrus.WithError(err).Errorf("发送邮件，DB报错")
			return
		}
		m.world.SendFunc(target, func() pbutil.Buffer {
			return mail.NewS2cReceiveMailMsg(data)
		})
		m.notifyMailCount(target)

		m.tlog.TlogMailFlowById(target, operate_type.MailSend, strconv.FormatUint(mailId, 10))
	})
	return true
}

func (m *MailModule) SendProtoMail(target int64, proto *shared_proto.MailProto, sendTime time.Time) bool {
	if target == 0 || npcid.IsNpcId(target) {
		return false
	}

	time := sendTime

	mailId := m.newMailId()

	proto.Id = i64.ToBytesU64(mailId)
	proto.SendTime = timeutil.Marshal32(time)
	proto.HasPrize = proto.Prize != nil

	data, err := proto.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("发送Proto邮件，proto marshal报错")
		return false
	}
	keep := proto.Keep
	hasReport := proto.HasReport
	hasPrize := proto.Prize != nil
	m.queue.MustFunc(func() {
		err = ctxfunc.Timeout3s(func(ctx context.Context) error {

			return m.db.CreateMail(ctx, mailId, target, data, keep, hasReport, hasPrize, 0, timeutil.Marshal64(time))
		})
		if err != nil {
			logrus.WithError(err).Errorf("发送邮件，DB报错")
			return
		}

		m.world.SendFunc(target, func() pbutil.Buffer {
			return mail.NewS2cReceiveMailMsg(data)
		})

		m.notifyMailCount(target)

		go call.CatchPanic(func() {
			m.tlog.TlogMailFlowById(target, operate_type.MailSend, strconv.FormatUint(mailId, 10))
		}, "tlog 发邮件")
	})

	return true
}

// 处理消息

func invalidParam(p int32) bool {
	return p < 0 || p > 2
}

func (m *MailModule) notifyMailCount(targetId int64) {
	m.world.FuncHero(targetId, func(id int64, hc iface.HeroController) {
		m.notifyMailCountWithSender(hc)
	})

}

func (m *MailModule) notifyMailCountWithSender(hc iface.HeroController) {
	if hc == nil {
		return
	}

	targetId := hc.Id()

	f := func() bool {
		if m.db.CallingTimes() > constants.DBBusyCallingTimes {
			return false
		}

		var hasPrizeNotCollectedCount, hasReportNotReadedCount, hasYwReportNotReadedCount, hasBzReportNotReadedCount, noReportNotReadedCount int
		err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			hasPrizeNotCollectedCount, err = m.db.LoadMailCountHasPrizeNotCollected(ctx, targetId)
			if err != nil {
				return
			}

			hasReportNotReadedCount, err = m.db.LoadMailCountHasReportNotReaded(ctx, targetId, 0)
			if err != nil {
				return
			}

			hasYwReportNotReadedCount, err = m.db.LoadMailCountHasReportNotReaded(ctx, targetId, 1)
			if err != nil {
				return
			}

			hasBzReportNotReadedCount, err = m.db.LoadMailCountHasReportNotReaded(ctx, targetId, 2)
			if err != nil {
				return
			}

			noReportNotReadedCount, err = m.db.LoadMailCountNoReportNotReaded(ctx, targetId)
			if err != nil {
				return
			}

			return
		})

		if err != nil {
			logrus.WithError(err).Errorf("notifyMailCount报错了")
			return true
		}

		hc.Send(mail.NewS2cNotifyMailCountMsg(int32(hasPrizeNotCollectedCount), int32(hasReportNotReadedCount), int32(hasYwReportNotReadedCount), int32(hasBzReportNotReadedCount), int32(noReportNotReadedCount)))

		return true
	}

	if !f() {
		hc.AddTickFunc(f)
	}
}

//gogen:iface
func (m *MailModule) ListMail(proto *mail.C2SListMailProto, hc iface.HeroController) {

	minId, ok := tryParseMailId("请求邮件列表", proto.MinId, true)
	if !ok {
		hc.Send(mail.ERR_LIST_MAIL_FAIL_INVALID_ID)
		return
	}

	if invalidParam(proto.Read) ||
		invalidParam(proto.Keep) ||
		invalidParam(proto.Report) ||
		invalidParam(proto.ReportTag) ||
		invalidParam(proto.HasPrize) ||
		invalidParam(proto.Collected) {
		logrus.WithField("read", proto.Read).
			WithField("keep", proto.Keep).
			WithField("report", proto.Report).
			WithField("report_tag", proto.ReportTag).
			WithField("has_prize", proto.HasPrize).
			WithField("collected", proto.Collected).
			Debugf("请求邮件列表，无效的参数")
		hc.Send(mail.ERR_LIST_MAIL_FAIL_INVALID_ID)
		return
	}

	if proto.Report == 2 && proto.HasPrize == 2 {
		// 又有战报，又有奖励的邮件是不存在的，返回空列表

		hc.Send(mail.NewS2cListMailMsg(proto.Read, proto.Keep, proto.Report, proto.ReportTag, proto.HasPrize, proto.Collected, nil))
		return
	}

	var mails []*shared_proto.MailProto
	err := ctxfunc.Timeout3s(func(ctx context.Context) error {
		var err error

		mails, err = m.db.LoadHeroMailList(ctx, hc.Id(), minId, proto.Keep, proto.Read, proto.Report, proto.HasPrize, proto.Collected, proto.ReportTag, u64.FromInt32(proto.Count))
		return err
	})
	if err != nil {
		logrus.WithError(err).Errorf("请求邮件列表，DB报错")
		hc.Send(mail.ERR_LIST_MAIL_FAIL_SERVER_ERROR)
		return
	}

	datas := make([][]byte, 0, len(mails))
	for _, mail := range mails {
		datas = append(datas, must.Marshal(mail))
	}

	hc.Send(mail.NewS2cListMailMsg(proto.Read, proto.Keep, proto.Report, proto.ReportTag, proto.HasPrize, proto.Collected, datas))
}

func tryParseMailId(name string, bid []byte, allowZero bool) (uint64, bool) {
	id, ok := i64.FromBytesU64(bid)
	if !ok {
		logrus.Debug(name + "，parse bid fail")
		return 0, false
	}

	if allowZero || id != 0 {
		return id, true
	} else {
		logrus.Debug(name + "，id == 0")
		return 0, false
	}
}

//gogen:iface
func (m *MailModule) DeleteMail(proto *mail.C2SDeleteMailProto, hc iface.HeroController) {

	id, ok := tryParseMailId("删除邮件", proto.Id, false)
	if !ok {
		hc.Send(mail.ERR_DELETE_MAIL_FAIL_INVALID_ID)
		return
	}

	var collectable bool
	if err := ctxfunc.Timeout3s(func(ctx context.Context) error {
		var err error
		collectable, err = m.db.IsCollectableMail(ctx, id)
		if err != nil {
			return err
		}

		if !collectable {
			return m.db.DeleteMail(ctx, id, hc.Id())
		}
		return nil
	}); err != nil {
		logrus.WithError(err).Errorf("删除邮件，DB报错")
		hc.Send(mail.ERR_DELETE_MAIL_FAIL_SERVER_ERROR)
		return
	}

	if collectable {
		logrus.Errorf("删除邮件，邮件有奖励可以领取")
		hc.Send(mail.ERR_DELETE_MAIL_FAIL_NOT_EMPTY)
		return
	}

	hc.Send(mail.NewS2cDeleteMailMsg(proto.Id))

	m.notifyMailCountWithSender(hc)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.RemoveInvestigationMail(id)
		result.Ok()
	})

	m.tlog.TlogMailFlowById(hc.Id(), operate_type.MailDel, strconv.FormatUint(id, 10))
}

//gogen:iface
func (m *MailModule) KeepMail(proto *mail.C2SKeepMailProto, hc iface.HeroController) {

	id, ok := tryParseMailId("收藏邮件", proto.Id, false)
	if !ok {
		hc.Send(mail.ERR_KEEP_MAIL_FAIL_INVALID_ID)
		return
	}

	if err := ctxfunc.Timeout3s(func(ctx context.Context) error {
		return m.db.UpdateMailKeep(ctx, id, hc.Id(), proto.Keep)
	}); err != nil {
		logrus.WithError(err).Errorf("收藏邮件，DB报错")
		hc.Send(mail.ERR_KEEP_MAIL_FAIL_SERVER_ERROR)
		return
	}

	hc.Send(mail.NewS2cKeepMailMsg(proto.Id, proto.Keep))
}

//gogen:iface
func (m *MailModule) ReadMail(proto *mail.C2SReadMailProto, hc iface.HeroController) {

	id, ok := tryParseMailId("已读邮件", proto.Id, false)
	if !ok {
		hc.Send(mail.ERR_READ_MAIL_FAIL_INVALID_ID)
		return
	}

	if err := ctxfunc.Timeout3s(func(ctx context.Context) error {
		return m.db.UpdateMailRead(ctx, id, hc.Id(), true)
	}); err != nil {
		logrus.WithError(err).Errorf("已读邮件，DB报错")
		hc.Send(mail.ERR_READ_MAIL_FAIL_SERVER_ERROR)
		return
	}

	hc.Send(mail.NewS2cReadMailMsg(proto.Id))

	m.notifyMailCountWithSender(hc)

	m.tlog.TlogMailFlowById(hc.Id(), operate_type.MailRead, strconv.FormatUint(id, 10))
}

//gogen:iface
func (m *MailModule) ProcessCollectMailPrize(proto *mail.C2SCollectMailPrizeProto, hc iface.HeroController) {

	id, ok := tryParseMailId("领取邮件奖励", proto.Id, false)
	if !ok {
		hc.Send(mail.ERR_COLLECT_MAIL_PRIZE_FAIL_INVALID_ID)
		return
	}

	var prize *resdata.Prize
	err := ctxfunc.Timeout3s(func(ctx context.Context) error {
		var err error
		prize, err = m.db.LoadCollectMailPrize(ctx, id, hc.Id())
		return err
	})
	if err != nil {
		logrus.WithError(err).Errorf("领取邮件奖励，DB报错")
		hc.Send(mail.ERR_COLLECT_MAIL_PRIZE_FAIL_SERVER_ERROR)
		return
	}

	if prize == nil {
		logrus.Debugf("领取邮件奖励，邮件中没有奖励或者已经被领取")
		hc.Send(mail.ERR_COLLECT_MAIL_PRIZE_FAIL_NOT_PRIZE)
		return
	}

	if err := ctxfunc.Timeout3s(func(ctx context.Context) error {
		return m.db.UpdateMailCollected(ctx, id, hc.Id(), true)
	}); err != nil {
		logrus.Errorf("领取邮件奖励，更新领取状态失败")
		hc.Send(mail.ERR_COLLECT_MAIL_PRIZE_FAIL_SERVER_ERROR)
		return
	}

	// 给玩家加奖励
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hctx := heromodule.NewContext2(m.datas, m.broadcast, m.guildSnapshot, m.world, m.chat, m.tlog, operate_type.MailCollectMailPrize)
		heromodule.AddPrize(hctx, hero, result, prize, m.time.CurrentTime())

		m.tlog.TlogMailFlow(hero, operate_type.MailCollect, strconv.FormatUint(id, 10))
		result.Ok()
		return
	}) {
		logrus.Debugf("领取邮件奖励，领取失败")

		// 将奖励加回去
		if err := ctxfunc.Timeout3s(func(ctx context.Context) error {
			return m.db.UpdateMailCollected(ctx, id, hc.Id(), false)
		}); err != nil {
			logrus.WithField("hero_id", hc.Id()).
				WithField("prize", prize.Encode().String()).
				WithError(err).Errorf("领取邮件奖励失败，但是加回去的时候也失败了，玩家丢失邮件奖励")
		}
		return
	}

	hc.Send(mail.NewS2cCollectMailPrizeMsg(proto.Id))

	m.notifyMailCountWithSender(hc)
}

var readEmptyMsg = mail.NewS2cReadMultiMsg(nil, false, nil).Static()

//gogen:iface
func (m *MailModule) ProcessReadMulti(proto *mail.C2SReadMultiProto, hc iface.HeroController) {

	ids, ok := i64.FromBytesArrayU64(proto.Ids)
	if !ok || len(ids) == 0 {
		logrus.Debug("邮件一键已读，无效的id")
		hc.Send(readEmptyMsg)
		return
	}

	var prize *resdata.Prize
	err := ctxfunc.Timeout3s(func(ctx context.Context) error {
		var err error
		prize, err = m.db.ReadMultiMail(ctx, hc.Id(), ids, proto.Report)
		return err
	})
	if err != nil {
		logrus.WithError(err).Errorf("邮件一键已读，DB报错")
		hc.Send(mail.ERR_COLLECT_MAIL_PRIZE_FAIL_SERVER_ERROR)
		return
	}

	// 给玩家加奖励
	if prize != nil {
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			hctx := heromodule.NewContext2(m.datas, m.broadcast, m.guildSnapshot, m.world, m.chat, m.tlog, operate_type.MailReadMulti)
			heromodule.AddPrize(hctx, hero, result, prize, m.time.CurrentTime())

			result.Ok()
			return
		}) {
			logrus.Debugf("邮件一键已读，领取失败")
			return
		}
	}

	var prizeBytes []byte
	if prize != nil {
		prizeBytes = util.SafeMarshal(prize.Encode())
	}
	hc.Send(mail.NewS2cReadMultiMsg(proto.Ids, proto.Report, prizeBytes))

	m.notifyMailCountWithSender(hc)
}

var deleteEmptyMsg = mail.NewS2cDeleteMultiMsg(nil, false).Static()

//gogen:iface
func (m *MailModule) ProcessDeleteMulti(proto *mail.C2SDeleteMultiProto, hc iface.HeroController) {

	ids, ok := i64.FromBytesArrayU64(proto.Ids)
	if !ok || len(ids) == 0 {
		logrus.Debug("邮件删除已读，无效的id")
		hc.Send(deleteEmptyMsg)
		return
	}

	err := ctxfunc.Timeout3s(func(ctx context.Context) error {
		return m.db.DeleteMultiMail(ctx, hc.Id(), ids, proto.Report)
	})
	if err != nil {
		logrus.WithError(err).Errorf("邮件删除已读，DB报错")
		return
	}

	hc.Send(mail.NewS2cDeleteMultiMsg(proto.Ids, proto.Report))

	m.notifyMailCountWithSender(hc)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		for _, id := range ids {
			hero.RemoveInvestigationMail(id)
		}
		result.Ok()
	})
}

//gogen:iface
func (m *MailModule) ProcessGetMail(proto *mail.C2SGetMailProto, hc iface.HeroController) {

	id, ok := i64.FromBytesU64(proto.Bid)
	if !ok {
		logrus.Debug("获取邮件ByID，id错误")
		hc.Send(mail.ERR_GET_MAIL_FAIL_MAIL_NOT_FOUND)
		return
	}

	var data []byte
	var err error
	ctxfunc.Timeout3s(func(ctx context.Context) error {
		data, err = m.db.LoadMail(ctx, id)
		return nil
	})
	if err != nil {
		logrus.WithError(err).Errorf("获取邮件列表，DB报错")
		hc.Send(mail.ERR_GET_MAIL_FAIL_MAIL_NOT_FOUND)
		return
	}

	if len(data) <= 0 {
		logrus.Debug("获取邮件ByID，邮件不存在")
		hc.Send(mail.ERR_GET_MAIL_FAIL_MAIL_NOT_FOUND)
		return
	}

	hc.Send(mail.NewS2cGetMailMsg(data))
}
