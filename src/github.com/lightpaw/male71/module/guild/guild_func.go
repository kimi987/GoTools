package guild

import (
	"bytes"
	"context"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/pushdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/module/rank/ranklist"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/lock"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"runtime/debug"
	"sync"
	"time"
)

func newGuildFunc(dep iface.ServiceDep, datas *config.ConfigDatas, hebi iface.HebiModule, db iface.DbService, mail iface.MailModule,
	realmService iface.RealmService, xiongNuModule iface.XiongNuModule, xiongNuService iface.XiongNuService, country iface.CountryService,
	baizhanService iface.BaiZhanService, rankModule iface.RankModule, pushService iface.PushService, mingcWarService iface.MingcWarService,
	maxGuildId int64, clearSelfGuildMsgCache func(int64), updateSnapshot func(g *sharedguilddata.Guild) *guildsnapshotdata.GuildSnapshot,
	guildSnapshotGetter func(int64) *guildsnapshotdata.GuildSnapshot, removeSnapshot func(int64), chat iface.ChatService,
	addGuildLog func(guildId int64, proto *shared_proto.GuildLogProto), templates []*guild_template) *guild_func {
	m := &guild_func{}
	m.dep = dep
	m.datas = datas
	m.hebiModule = hebi
	m.time = dep.Time()
	m.db = db
	m.world = dep.World()
	m.mail = mail
	m.chat = chat
	m.country = country
	m.mingcWarService = mingcWarService
	m.heroService = dep.HeroData()
	m.heroSnapshotService = dep.HeroSnapshot()
	m.realmService = realmService
	m.xiongNuService = xiongNuService
	m.xiongNuModule = xiongNuModule
	m.rankModule = rankModule
	m.baizhanService = baizhanService
	m.pushService = pushService
	m.guildIdGen = atomic.NewInt64(maxGuildId)
	m.clearSelfGuildMsgCache = clearSelfGuildMsgCache
	m.updateSnapshot = updateSnapshot
	m.guildSnapshotGetter = guildSnapshotGetter
	m.removeSnapshot = removeSnapshot
	m.addGuildLog = addGuildLog
	m.templates = templates

	m.freeNpcGuildChanged = true

	return m
}

type guild_func struct {
	dep iface.ServiceDep

	datas *config.ConfigDatas

	time iface.TimeService

	db iface.DbService

	world iface.WorldService

	mail iface.MailModule

	chat iface.ChatService

	heroService         iface.HeroDataService
	heroSnapshotService iface.HeroSnapshotService
	realmService        iface.RealmService
	xiongNuService      iface.XiongNuService
	xiongNuModule       iface.XiongNuModule
	baizhanService      iface.BaiZhanService
	pushService         iface.PushService
	mingcWarService     iface.MingcWarService
	country             iface.CountryService

	rankModule iface.RankModule
	hebiModule iface.HebiModule

	guildIdGen             *atomic.Int64
	clearSelfGuildMsgCache func(int64)

	// 更新联盟的镜像
	updateSnapshot      func(g *sharedguilddata.Guild) *guildsnapshotdata.GuildSnapshot
	guildSnapshotGetter func(id int64) *guildsnapshotdata.GuildSnapshot
	removeSnapshot      func(int64)
	addGuildLog         func(guildId int64, proto *shared_proto.GuildLogProto)

	templates []*guild_template

	// 出现以下情况，停止创建NPC联盟
	// 1、帮派id超出NPC联盟id上限
	// 2、所有Npc帮派的后缀都用完了
	// 这个值在第一次发现不能再创建npc帮派时候设置
	dontCreateNpcGuild bool

	// 空闲的Npc帮派
	freeNpcGuildChanged bool

	notFullGuildList  []int64      // 没有满的联盟id列表
	notFullGuildMutex sync.RWMutex // 没有满的联盟id列表的锁
}

func (m *guild_func) getRank(guildId int64, country *country.CountryData) (rank, rankByCountry uint64) {
	m.rankModule.SubTypeRRankListFunc(shared_proto.RankType_Guild, 0, func(rankList rankface.RRankList) {
		obj := rankList.GetRankObj(guildId)
		if obj != nil {
			rank = obj.Rank()
		}
	})

	m.rankModule.SubTypeRRankListFunc(shared_proto.RankType_Guild, country.Id, func(rankList rankface.RRankList) {
		obj := rankList.GetRankObj(guildId)
		if obj != nil {
			rankByCountry = obj.Rank()
		}
	})

	return
}

func (m *guild_func) newGuildId() int64 {
	return m.guildIdGen.Inc()
}

func (m *guild_func) walkNotFullGuild(walkFunc func(id int64) (endWalk bool)) {
	m.notFullGuildMutex.RLock()
	defer m.notFullGuildMutex.RUnlock()

	for _, id := range m.notFullGuildList {
		if walkFunc(id) {
			return
		}
	}
}

func (m *guild_func) addNotFullGuild(g *sharedguilddata.Guild) {
	if g.IsFull() {
		// 满了
		return
	}

	if npcGuildTemplate := g.GetNpcTemplate(); npcGuildTemplate != nil && npcGuildTemplate.RejectUserJoin {
		// npc联盟
		return
	}

	m.notFullGuildMutex.Lock()
	defer m.notFullGuildMutex.Unlock()
	m.notFullGuildList = i64.AddIfAbsent(m.notFullGuildList, g.Id())
}

func (m *guild_func) removeNotFullGuild(id int64) {
	m.notFullGuildMutex.Lock()
	defer m.notFullGuildMutex.Unlock()
	m.notFullGuildList = i64.RemoveIfPresent(m.notFullGuildList, id)
}

func (m *guild_func) updateGuildRankObj(guild *sharedguilddata.Guild) {
	m.rankModule.AddOrUpdateRankObj(ranklist.NewGuildRankObj(m.guildSnapshotGetter, m.heroSnapshotService.Get, guild))
}

func (m *guild_func) createGuild(hc iface.HeroController, guildName, flagName string, guilds sharedguilddata.Guilds) (errMsg msg.ErrMsg) {
	if guilds.GetIdByName(guildName) != 0 {
		logrus.Debugf("创建联盟，联盟名字重复")
		errMsg = guild.ErrCreateGuildFailNameDuplicate
		return
	}

	if guilds.GetIdByFlagName(flagName) != 0 {
		logrus.Debugf("创建联盟，联盟旗号重复")
		errMsg = guild.ErrCreateGuildFailFlagNameDuplicate
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.GuildCreate)

	var newGuildId int64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.GuildId() != 0 {
			logrus.Debugf("创建联盟，已经有联盟了")
			errMsg = guild.ErrCreateGuildFailInTheGuild
			return
		}

		if !heromodule.HasEnoughCost(hero, m.datas.GuildConfig().CreateGuildCost) {
			logrus.Debugf("创建联盟，消耗不足")
			errMsg = guild.ErrCreateGuildFailCostNotEnough
			return
		}

		ctime := m.time.CurrentTime()
		countryData := m.datas.GetCountryData(hero.CountryId())
		if countryData == nil {
			errMsg = guild.ErrCreateGuildFailHeroNoCountry
			return
		}

		newGuildId = m.newGuildId()
		newGuild := sharedguilddata.NewGuild(newGuildId, guildName, flagName, m.datas, ctime)
		newGuild.SetLeader(hc.Id())
		newGuild.SetCountry(countryData)
		member := sharedguilddata.NewMember(hc.Id(), m.datas.GuildClassLevelData().MaxKeyData, ctime)
		newGuild.AddMember(member)

		// 设置联盟目标
		newGuild.TryUpdateTarget(m.datas.GuildConfig(), ctime, 0)

		// 创建帮派
		guildBytes, err := newGuild.Marshal()
		if err != nil {
			logrus.WithError(err).Errorf("创建联盟，Guild.Marshal 报错")
			errMsg = guild.ErrCreateGuildFailServerError
			return
		}

		err = ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			return m.db.CreateGuild(ctx, newGuildId, guildBytes)
		})
		if err != nil {
			logrus.WithError(err).Errorf("创建联盟，DB.CreateGuild 报错")
			errMsg = guild.ErrCreateGuildFailServerError
			return
		}

		// 创建成功，添加到guilds
		m.updateSnapshot(newGuild)
		guilds.Add(newGuild)

		m.updateGuildRankObj(newGuild)
		m.addNotFullGuild(newGuild)

		// 扣钱，使用reduce anyway 有多少扣多少
		heromodule.ReduceCostAnyway(hctx, hero, result, m.datas.GuildConfig().CreateGuildCost)

		// 设置帮派信息
		m.updateHeroGuild(hero, result, newGuild, member, true)

		result.Add(guild.NewS2cCreateGuildMsg(must.Marshal(newGuild.EncodeClient(true, m.heroSnapshotService, m.getRank, m.xiongNuService.IsTodayStarted))))

		// 加入帮派任务
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_JOIN_GUILD)

		if data := m.datas.GuildLogHelp().CreateGuild; data != nil {
			proto := data.NewHeroLogProto(ctime, hero.IdBytes(), hero.Head())
			proto.Text = data.Text.New().WithHeroName(hero.Name()).JsonString()

			m.addGuildLog(newGuild.Id(), proto)
		}

		// 系统广播
		hctx := heromodule.NewContext(m.dep, operate_type.GuildCreate)
		if d := hctx.BroadcastHelp().GuildCreate; d != nil {
			hctx.AddBroadcast(d, hero, result, 0, 0, func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickGuildFields(data.KeyGuild, hctx.GetFlagGuildName(newGuild), newGuild.Id())
				text.WithClickHeroFields(data.KeyName, hero.Name(), hero.Id())
				return text
			})
		}

		// 初次加入联盟邮件
		m.trySendFirstJoinGuildMail(hero, false)

		result.Add(guild.NewS2cUpdateHeroGuildMsg(0, newGuild.NewHeroGuildProto()))

		result.Changed()
		result.Ok()

		// tlog
		hctx.Tlog().TlogGuildFlow(hero, operate_type.GuildCreate.Id(), u64.FromInt64(newGuild.Id()), newGuild.LevelData().Level, 1)

	}) {
		return
	}

	if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		return m.db.UpdateHeroGuildId(ctx, hc.Id(), newGuildId)
	}); err != nil {
		logrus.WithError(err).Errorf("保存玩家联盟ID，更新超时 heroId:%v guildId:%v", hc.Id(), newGuildId)
	}

	m.hebiModule.UpdateGuildInfo(hc.Id(), newGuildId)

	return
}

func (m *guild_func) leaveGuild(hc iface.HeroController, guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	heroId := self.Id()

	if g.LeaderId() == heroId && g.MemberCount() > 1 {
		logrus.Debugf("退出联盟，你是盟主，而且联盟还有人")
		errMsg = guild.ErrLeaveGuildFailLeader
		return
	}

	// 名城战
	mcId, t := m.mingcWarService.GuildMcWarType(g.Id())
	if t == shared_proto.MingcWarGuildType_MC_G_ATK && g.MemberCount() == 1 {
		errMsg = guild.ErrLeaveGuildFailIsMcWarAtk
		return
	}

	if mc := m.dep.Mingc().Mingc(mcId); mc != nil {
		if mc.HostGuildId() == g.Id() && g.MemberCount() == 1 {
			errMsg = guild.ErrLeaveGuildFailIsMcWarDef
			return
		}

		state, _, _ := m.mingcWarService.CurrMcWarStage()
		if state == int32(shared_proto.MingcWarState_MC_T_FIGHT) {
			errMsg = guild.ErrLeaveGuildFailInMcWarFight
			return
		}
	}

	// 集结

	if g.IsResistXiongNuDefenders(heroId) && m.xiongNuService.IsStarted(g.Id()) {
		logrus.Debugf("退出联盟，自己在联盟匈奴入侵防守队伍中，且联盟开启了匈奴入侵")
		errMsg = guild.ErrLeaveGuildFailXiongNuDefender
		return
	}

	if errMsg = m.tryRemoveHeroGuild(g, heroId, heroId, false,
		guild.ErrLeaveGuildFailAssembly, guild.ErrLeaveGuildFailServerError); errMsg != nil {
		return
	}

	if g.MemberCount() <= 0 {
		guilds.Remove(g.Id())
		m.removeNotFullGuild(g.Id())
		m.removeSnapshot(g.Id())
		m.mingcWarService.CleanOnGuildRemoved(g.Id())

		ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
			if err := m.db.DeleteGuild(ctx, g.Id()); err != nil {
				// 出错了，只能打个日志出来，到下一次开服的时候，清掉所有的没有成员的帮派
				logrus.WithError(err).Errorf("退出联盟，联盟没人了，删除出错")
			}
			return nil
		})

		// 调用删除联盟
		m.rankModule.RemoveRankObj(shared_proto.RankType_Guild, g.Id())
	}

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.LEAVE_GUILD_S2C
	broadcastChanged = true

	// 加入推荐列表
	m.dep.Guild().AddRecommendInviteHeros(heroId)

	m.hebiModule.UpdateGuildInfo(hc.Id(), 0)

	m.dep.Tlog().TlogGuildFlowById(self.Id(), operate_type.GuildLeave.Id(), u64.FromInt64(g.Id()), g.LevelData().Level, u64.FromInt(g.MemberCount()))
	if g.MemberCount() <= 0 {
		m.dep.Tlog().TlogGuildFlowById(self.Id(), operate_type.GuildDestroyed.Id(), u64.FromInt64(g.Id()), g.LevelData().Level, u64.FromInt(g.MemberCount()))
	}
	return
}

func (m *guild_func) kickOther(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, kickTargetId int64) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("联盟踢人，Npc联盟不允许踢人")
		errMsg = guild.ErrKickOtherFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.KickLowerMember
	}) {
		logrus.Debugf("联盟踢人，没有权限")
		errMsg = guild.ErrKickOtherFailDeny
		return
	}

	kickMember := g.GetMember(kickTargetId)
	if kickMember == nil {
		logrus.Debugf("联盟踢人，other == nil")
		errMsg = guild.ErrKickOtherFailNotInGuild
		return
	}

	if self.ClassLevelData().Level <= kickMember.ClassLevelData().Level {
		logrus.Debugf("联盟踢人，你级别没有他高")
		errMsg = guild.ErrKickOtherFailDeny
		return
	}

	// 帮主弹劾期间，不踢人
	if g.GetImpeachLeader() != nil {
		logrus.Debugf("联盟踢人，帮主弹劾期间不能踢人")
		errMsg = guild.ErrKickOtherFailImpeachLeader
		return
	}

	// 超出上限
	if !g.IsFull() && g.GetKickMemberCount() >= m.datas.GuildConfig().DailyMaxKickCount {
		logrus.Debugf("联盟踢人，超出每日上限")
		errMsg = guild.ErrKickOtherFailLimit
		return
	}

	// 名城战
	state, _, _ := m.mingcWarService.CurrMcWarStage()
	if state == int32(shared_proto.MingcWarState_MC_T_FIGHT) {
		mcId, _ := m.mingcWarService.GuildMcWarType(g.Id())
		if m.dep.Mingc().Mingc(mcId) != nil {
			errMsg = guild.ErrKickOtherFailInMcWarFight
			return
		}
	}

	if g.IsResistXiongNuDefenders(kickMember.Id()) && m.xiongNuService.IsStarted(g.Id()) {
		logrus.Debugf("联盟踢人，目标在联盟匈奴入侵防守队伍中，且对方联盟开启了匈奴入侵")
		errMsg = guild.ErrKickOtherFailXiongNuDefender
		return
	}

	if errMsg = m.tryRemoveHeroGuild(g, kickMember.Id(), self.Id(), true,
		guild.ErrKickOtherFailAssembly, guild.ErrKickOtherFailServerError); errMsg != nil {
		return
	}

	g.IncKickMemberCount()

	m.clearSelfGuildMsgCache(g.Id())

	var kickMemberName string
	kickHero := m.heroSnapshotService.Get(kickTargetId)
	if kickHero == nil {
		kickMemberName = idbytes.PlayerName(kickTargetId)
	} else {
		kickMemberName = kickHero.Name
	}

	successMsg = guild.NewS2cKickOtherMsg(kickMember.IdBytes(), kickMemberName)
	broadcastChanged = true

	// 加入推荐列表
	m.dep.Guild().AddRecommendInviteHeros(kickMember.Id())

	m.hebiModule.UpdateGuildInfo(kickMember.Id(), 0)

	// 被踢邮件
	leader := m.heroSnapshotService.Get(g.LeaderId())
	if leader != nil {
		leaderName := m.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(g.FlagName(), leader.Name)
		countryName := g.Country().Name
		proto := m.datas.MailHelp().GuildBeKickedOut.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Text = m.datas.MailHelp().GuildBeKickedOut.NewTextFields().WithCountry(countryName).WithLeader(leaderName).JsonString()
		ctime := m.time.CurrentTime()
		m.mail.SendProtoMail(kickMember.Id(), proto, ctime)
	}

	m.dep.Tlog().TlogGuildFlowById(self.Id(), operate_type.GuildDismiss.Id(), u64.FromInt64(g.Id()), g.LevelData().Level, u64.FromInt(g.MemberCount()))
	return
}

func (m *guild_func) tryRemoveHeroGuild(g *sharedguilddata.Guild, removeHeroId, operateHeroId int64, kick bool, assemblyErrMsg, serverErrorMsg msg.ErrMsg) (errMsg msg.ErrMsg) {
	hasError, err := m.heroService.FuncWithSendError(removeHeroId, func(hero *entity.Hero, result herolock.LockResult) {
		if hero.HasAssemblyTroop() {
			errMsg = assemblyErrMsg
			return
		}

		m.doRemoveHeroGuild(g, hero.Id(), operateHeroId, hero, result, true)

		if kick {
			result.Add(guild.SELF_BEEN_KICKED_S2C)
		}

		result.Changed()
		result.Ok()
		return
	})

	if errMsg != nil {
		return
	}

	if hasError {
		if err == lock.ErrEmpty {
			// 没有这个英雄？？？
			logrus.WithField("hero_id", removeHeroId).Error("移除英雄联盟，居然英雄不存在？？？Bug，这里还是删掉吧")

			m.doRemoveHeroGuild(g, removeHeroId, operateHeroId, nil, nil, false)
		}

		return serverErrorMsg
	} else {
		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			return m.db.UpdateHeroGuildId(ctx, removeHeroId, 0)
		}); err != nil {
			logrus.WithError(err).Errorf("保存玩家联盟ID，更新超时 heroId:%v guildId:%v", removeHeroId, 0)
		}
	}

	return nil
}

var cleanGuildSeekHelpList = guild.NewS2cListGuildSeekHelpMsg(nil).Static()

func (m *guild_func) doRemoveHeroGuild(g *sharedguilddata.Guild, heroId, operateHeroId int64,
	hero *entity.Hero, result herolock.LockResult, notifyHeroGuildChanged bool) {

	ctime := m.time.CurrentTime()

	// 玩家删除
	if hero != nil {
		if hero.GuildId() == g.Id() {
			m.updateHeroGuild(hero, result, nil, nil, notifyHeroGuildChanged)
			result.Add(cleanGuildSeekHelpList)

			// 原来联盟的逼格宝箱没有领取，发邮件处理
			if fullBigBoxData := g.GetFullBigBoxData(); fullBigBoxData != nil {
				if g.RemoveFullBigBoxMemberId(hero.Id()) {
					proto := m.datas.MailHelp().GuildBigBox.NewTextMail(shared_proto.MailType_MailNormal)
					proto.Prize = fullBigBoxData.PlunderPrize.GetPrize().PrizeProto()

					m.mail.SendProtoMail(hero.Id(), proto, ctime)
				}
			}
		} else {
			logrus.WithField("stack", string(debug.Stack())).Errorf("guild_func.doRemoveHeroGuild 存在玩家联盟id不在跟原联盟id不相同的严重bug!%d,%d", hero.GuildId(), g.Id())
		}
	}

	guildOldIsFull := g.IsFull()

	// 帮派删除
	g.RemoveMember(heroId, ctime.Unix())

	// 联盟日志
	if hero != nil {
		if operateHeroId == hero.Id() {
			if d := m.datas.GuildLogHelp().LeaveGuild; d != nil {
				proto := d.NewHeroLogProto(ctime, hero.IdBytes(), hero.Head())
				proto.Text = d.Text.New().WithHeroName(hero.Name()).JsonString()
				m.dep.Guild().AddLogWithMemberIds(g.Id(), g.AllUserMemberIds(), proto)
			}
		} else {
			if d := m.datas.GuildLogHelp().KickLeaveGuild; d != nil {
				var operatorName string
				if operator := m.heroSnapshotService.Get(operateHeroId); operator != nil {
					operatorName = operator.Name
				} else {
					operatorName = idbytes.PlayerName(operateHeroId)
				}

				proto := d.NewHeroLogProto(ctime, hero.IdBytes(), hero.Head())
				proto.Text = d.Text.New().WithHeroName(hero.Name()).
					WithLeaderName(operatorName).JsonString()
				m.dep.Guild().AddLogWithMemberIds(g.Id(), g.AllUserMemberIds(), proto)
			}
		}

	}

	if guildOldIsFull {
		m.addNotFullGuild(g) // 没有到解散
	}

	allUserMemberIds := g.AllUserMemberIds()

	//if hero != nil {
	//	m.world.MultiSendIgnore(allUserMemberIds, guild.NewS2cLeaveGuildForOtherMsg(hero.IdBytes(), hero.Name(), hero.Head()), operateHeroId)
	//}

	// 遍历联盟求助，把离开联盟的家伙的求助删掉
	g.RangeSeekHelp(func(proto *shared_proto.GuildSeekHelpProto) (isContinue bool) {
		if bytes.Equal(idbytes.ToBytes(heroId), proto.HeroId) {
			g.RemoveSeekHelp(proto.Id)

			m.world.MultiSend(allUserMemberIds, guild.NewS2cRemoveGuildSeekHelpMsg(proto.Id))
		}
		return true
	})

	oldLeaderMember := g.GetMember(g.LeaderId())

	isNpcLeader := g.IsNpcLeader()
	if g.TryImpeachLeader(
		m.datas.GuildClassLevelData().MinKeyData,
		m.datas.GuildClassLevelData().MaxKeyData) {
		// 弹劾成功，广播 TODO

		m.afterImpeachLeader(g, oldLeaderMember, ctime)

		m.country.ChangeKing(g.CountryId(), g.LeaderId())

		if isNpcLeader {
			// Npc被弹劾掉，刷新新的Npc联盟
			m.doUpdateFreeNpcGuildChanged()
		}
	}

	m.updateSnapshot(g)
}

func (m *guild_func) updateText(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, text string) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
	if g.IsNpcLeader() {
		logrus.Debugf("更新联盟宣言，Npc联盟不允许踢人")
		errMsg = guild.ErrUpdateTextFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateText
	}) {
		logrus.Debugf("更新联盟宣言，没有权限")
		errMsg = guild.ErrUpdateTextFailDeny
		return
	}

	g.SetText(text)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.NewS2cUpdateTextMsg(text)
	return
}

func (m *guild_func) updateInternalText(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, text string) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("更新联盟内部宣言，Npc联盟不允许操作")
		errMsg = guild.ErrUpdateInternalTextFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateInternalText
	}) {
		logrus.Debugf("更新联盟内部宣言，没有权限")
		errMsg = guild.ErrUpdateInternalTextFailDeny
		return
	}

	g.SetInternalText(text)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.NewS2cUpdateInternalTextMsg(text)

	if data := m.datas.GuildLogHelp().UpdateInternalText; data != nil {
		if hero := m.heroSnapshotService.Get(self.Id()); hero != nil {
			ctime := m.time.CurrentTime()
			proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
			proto.Text = data.Text.New().WithHeroName(hero.Name).JsonString()

			m.addGuildLog(g.Id(), proto)
		}
	}

	return
}

func (m *guild_func) updateLabel(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, label []string) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("更新联盟标签，Npc联盟不允许操作")
		errMsg = guild.ErrUpdateGuildLabelFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateClassName
	}) {
		logrus.Debugf("更新联盟标签，没有权限")
		errMsg = guild.ErrUpdateGuildLabelFailDeny
		return
	}

	g.SetLabels(label)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.NewS2cUpdateGuildLabelMsg(label)
	return
}

func (m *guild_func) updateClassNames(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, names []string) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("更新阶级名称，Npc联盟不允许操作")
		errMsg = guild.ErrUpdateClassNamesFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateClassName
	}) {
		logrus.Debugf("更新阶级名称，没有权限")
		errMsg = guild.ErrUpdateClassNamesFailDeny
		return
	}

	// 职称名字已经被使用
	if m.checkDuplicateCustomClassTitle(names, g) ||
		m.checkDuplicateSystemClassTitle(names) {
		logrus.Debugf("更新职称，自定义职称名字已经被使用了")
		errMsg = guild.ErrUpdateClassNamesFailInvalidDuplicate
		return
	}

	g.SetClassNames(names)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.NewS2cUpdateClassNamesMsg(names)
	return
}

func (m *guild_func) updateClassTitle(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, proto *shared_proto.GuildClassTitleProto,
	heroIds []int64, setClassTitleData []*guild_data.GuildClassTitleData) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("更新职称，Npc联盟不允许操作")
		errMsg = guild.ErrUpdateClassTitleFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateClassTitle
	}) {
		logrus.Debugf("更新职称，没有权限")
		errMsg = guild.ErrUpdateClassTitleFailDeny
		return
	}

	// 职称名字已经被使用
	if m.checkDuplicateClassNames(proto.CustomClassTitleName, g) ||
		m.checkDuplicateSystemClassTitle(proto.CustomClassTitleName) {
		logrus.Debugf("更新职称，自定义职称名字已经被使用了")
		errMsg = guild.ErrUpdateClassTitleFailNameExist
		return
	}

	n := imath.Min(len(heroIds), len(setClassTitleData))
	var members []*sharedguilddata.GuildMember
	for i := 0; i < n; i++ {
		member := g.GetMember(heroIds[i])
		if member == nil {
			logrus.Debugf("更新职称，英雄不存在, %v", heroIds[i])
			errMsg = guild.ErrUpdateClassTitleFailInvalidMemberId
			return
		}

		members = append(members, member)
	}

	g.WalkMember(func(member *sharedguilddata.GuildMember) {
		// 删掉所有职称
		member.SetClassTitleData(nil)
	})

	for i, member := range members {
		member.SetClassTitleData(setClassTitleData[i])
	}

	g.SetClassTitle(proto)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.UPDATE_CLASS_TITLE_S2C
	broadcastChanged = true
	return
}

func (m *guild_func) checkDuplicateClassNames(names []string, g *sharedguilddata.Guild) bool {

	// 各个值内部是不重复的
	if len(g.GetClassNames()) > 0 {
		return util.StringArrayContainsAny(names, g.GetClassNames())
	}

	for _, d := range m.datas.GetGuildClassLevelDataArray() {
		if util.StringArrayContains(names, d.Name) {
			return true
		}
	}

	return false
}

func (m *guild_func) checkDuplicateCustomClassTitle(names []string, g *sharedguilddata.Guild) bool {
	return util.StringArrayContainsAny(names, g.GetCustomClassTitle())
}

func (m *guild_func) checkDuplicateSystemClassTitle(names []string) bool {

	for _, d := range m.datas.GetGuildClassTitleDataArray() {
		if util.StringArrayContains(names, d.Name) {
			return true
		}
	}

	return false
}

func (m *guild_func) updateFlagType(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, flagType uint64) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
	if g.IsNpcLeader() {
		logrus.Debugf("更新联盟旗帜，Npc联盟不允许操作")
		errMsg = guild.ErrUpdateFlagTypeFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateFlagType
	}) {
		logrus.Debugf("更新联盟旗帜，没有权限")
		errMsg = guild.ErrUpdateFlagTypeFailDeny
		return
	}

	g.SetFlagType(flagType)

	m.updateSnapshot(g)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.NewS2cUpdateFlagTypeMsg(u64.Int32(flagType))

	return
}

func (m *guild_func) updateMemberClassLevel(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, targetId int64, newClassLevelData *guild_data.GuildClassLevelData) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("修改成员阶级，Npc联盟不允许修改")
		errMsg = guild.ErrUpdateMemberClassLevelFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateLowerMemberClassLevel
	}) {
		logrus.Debugf("修改成员阶级，没有权限")
		errMsg = guild.ErrUpdateMemberClassLevelFailDeny
		return
	}

	target := g.GetMember(targetId)
	if target == nil {
		logrus.Debugf("修改成员阶级，目标玩家不在你的联盟")
		errMsg = guild.ErrUpdateMemberClassLevelFailTargetNotInGuild
		return
	}

	oldClassLevelData := target.ClassLevelData()

	// 特殊情况，帮主转给别人
	if newClassLevelData == m.datas.GuildClassLevelData().MaxKeyData {
		if self.ClassLevelData() != m.datas.GuildClassLevelData().MaxKeyData {
			logrus.Debugf("修改成员阶级，你不是帮主，要给别人帮主?")
			errMsg = guild.ErrUpdateMemberClassLevelFailDeny
			return
		}

		if g.MemberCount() >= int(m.datas.GuildConfig().ChangeLeaderCountdownMemberCount) {
			// 进入倒计时
			ctime := m.time.CurrentTime()
			countdownTime := ctime.Add(m.datas.GuildConfig().ChangeLeaderCountdownDuration)
			g.ChangeLeaderCountDown(targetId, countdownTime)

			m.clearSelfGuildMsgCache(g.Id())

			successMsg = nil
			broadcastChanged = true
			m.world.MultiSend(g.AllUserMemberIds(), guild.NewS2cUpdateMemberClassLevelMsg(target.IdBytes(), int32(newClassLevelData.Level)))
			return
		}
		if g.HasChangeLeaderCountDown() {
			g.CancelChangeLeader()
		}

		// 帮主变成平民
		self.SetClassLevelData(m.datas.GuildClassLevelData().MinKeyData)

		g.SetLeader(targetId)

		m.afterDemiseLeader(g, self.Id(), targetId)
		m.updateSnapshot(g)

		// 换国王
		m.country.ChangeKing(g.CountryId(), targetId)

	} else {
		if self.ClassLevelData().Level <= newClassLevelData.Level {
			logrus.Debugf("修改成员阶级，只能赋予比你自己更低级的职位")
			errMsg = guild.ErrUpdateMemberClassLevelFailDeny
			return
		}

		if self.ClassLevelData().Level <= target.ClassLevelData().Level {
			logrus.Debugf("修改成员阶级，目标玩家不比你的等级低")
			errMsg = guild.ErrUpdateMemberClassLevelFailDenyTarget
			return
		}

		if target.ClassLevelData() == newClassLevelData {
			// 就发给你自己好了
			errMsg = guild.ErrUpdateMemberClassLevelFailDenyTarget
			return
		}

		// 目标职位已满
		if g.IsClassFull(newClassLevelData) {
			logrus.Debugf("修改成员阶级，目标职位已经满员")
			errMsg = guild.ErrUpdateMemberClassLevelFailClassFull
			return
		}
	}

	target.SetClassLevelData(newClassLevelData)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = nil
	broadcastChanged = true
	m.world.MultiSend(g.AllUserMemberIds(), guild.NewS2cUpdateMemberClassLevelMsg(target.IdBytes(), int32(newClassLevelData.Level)))

	if data := m.datas.GuildLogHelp().UpdateMemberClass; data != nil {
		if hero := m.heroSnapshotService.Get(targetId); hero != nil {
			ctime := m.time.CurrentTime()
			proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
			proto.Text = data.Text.New().WithHeroName(hero.Name).WithClassName(newClassLevelData.Name).JsonString()

			m.addGuildLog(g.Id(), proto)
		}
	}

	// 邮件通知
	m.sendUpdateMemberClassLevelMail(self, targetId, oldClassLevelData, newClassLevelData)
	return
}

func (m *guild_func) sendUpdateMemberClassLevelMail(operator *sharedguilddata.GuildMember,
	targetId int64, oldClassLevelData, newClassLevelData *guild_data.GuildClassLevelData) {
	if newClassLevelData == m.datas.GuildClassLevelData().MaxKeyData {
		// 盟主的不在这里处理
		return
	}

	// 邮件通知
	var selfName string
	if hero := m.heroSnapshotService.Get(operator.Id()); hero != nil {
		selfName = hero.Name
	} else {
		selfName = idbytes.PlayerName(operator.Id())
	}

	if newClassLevelData.Level > oldClassLevelData.Level {
		mailData := m.datas.MailHelp().GuildClassLevelUp
		proto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Text = mailData.NewTextFields().
			WithNewClassLevel(newClassLevelData.Name).
			WithClassLevel(operator.ClassLevelData().Name).
			WithName(selfName).JsonString()
		m.mail.SendProtoMail(targetId, proto, m.time.CurrentTime())
	} else {
		mailData := m.datas.MailHelp().GuildClassLevelDown
		proto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Text = mailData.NewTextFields().
			WithOldClassLevel(oldClassLevelData.Name).
			WithNewClassLevel(newClassLevelData.Name).
			WithClassLevel(operator.ClassLevelData().Name).
			WithName(selfName).JsonString()
		m.mail.SendProtoMail(targetId, proto, m.time.CurrentTime())
	}
}

func (m *guild_func) cancelChangeLeader(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if self.ClassLevelData() != m.datas.GuildClassLevelData().MaxKeyData {
		logrus.Debugf("取消变更帮主，你不是帮主")
		errMsg = guild.ErrCancelChangeLeaderFailDeny
		return
	}

	if !g.HasChangeLeaderCountDown() {
		logrus.Debugf("取消变更帮主，当前没有变更帮主")
		errMsg = guild.ErrCancelChangeLeaderFailNotChangeLeader
		return
	}

	g.CancelChangeLeader()
	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.CANCEL_CHANGE_LEADER_S2C
	broadcastChanged = true
	return
}

func (m *guild_func) donation(hc iface.HeroController, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, seq uint64) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
	hctx := heromodule.NewContext(m.dep, operate_type.GuildDonate)
	ctime := m.time.CurrentTime()

	var heroName string
	var donateData *guild_data.GuildDonateData
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		times, ok := hero.GetDonateTimes(seq)
		if !ok {
			logrus.Debugf("联盟捐献，无效的sequence")
			errMsg = guild.ErrDonateFailInvalidSequence
			return
		}

		b := hero.Domestic().GetBuilding(domestic_data.DonateBuilding)
		if b == nil || times >= b.Effect.GuildDonateTimes {
			logrus.Debugf("联盟捐献，已经达到外使院捐献次数上限")
			errMsg = guild.ErrDonateFailDonateTimesLimit
			return
		}

		if m.datas.GuildConfig().GuildDonateNeedHeroLevel > hero.Level() {
			logrus.Debugf("联盟捐献，君主等级不够")
			errMsg = guild.ErrDonateFailLevelNotEnough
			return
		}

		donateData = m.datas.GetGuildDonateData(guild_data.DonateId(seq, times+1))
		if donateData == nil {
			logrus.Debugf("联盟捐献，已经是最大捐献次数")
			errMsg = guild.ErrDonateFailMaxTimes
			return
		}

		if !heromodule.TryReduceCost(hctx, hero, result, donateData.Cost) {
			logrus.Debugf("联盟捐献，捐献消耗不足")
			errMsg = guild.ErrDonateFailCostNotEnough
			return
		}

		heroName = hero.Name()

		// 加捐献次数
		hero.AddDonateTimes(seq)

		// 加捐献币
		heromodule.AddGuildContributionCoin(hctx, hero, result, donateData.ContributionCoin)

		result.Changed()
		result.Ok()

		// 处理帮派相关
		if uint64(g.GetDonateHeroCount()) < g.LevelData().MemberCount || g.IsDonate(hero.Id()) {
			// 加帮派建设值
			g.AddDonateBuildingAmount(donateData.GuildBuildingAmount, ctime, hero.Id())
		}

		// 加自己的贡献值
		self.AddContribution(donateData.ContributionAmount)

		// 加自己的捐献值
		self.AddDonation(donateData.DonationAmount)

		// 加自己的捐献元宝
		self.AddDonateYuanbao(donateData.Cost.Yuanbao)

		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_GUILD_DONATE)
		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_GuildDonate)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_GUILD_DONATE)

		m.dep.Tlog().TlogGuildFlow(hero, operate_type.GuildDonate.Id(), u64.FromInt64(g.Id()), g.LevelData().Level, u64.FromInt(g.MemberCount()))

		return
	}) {
		if errMsg == nil {
			errMsg = guild.ErrDonateFailServerError
		}
		return
	}

	if errMsg != nil {
		return
	}

	// 增加联盟任务进度
	if g.LevelData().Level >= m.datas.GuildGenConfig().TaskOpenLevel {
		if data := m.datas.GetGuildTaskData(u64.FromInt32(int32(server_proto.GuildTaskType_Donate))); data != nil {
			if g.AddGuildTaskProgress(data, 1) {
				m.world.MultiSend(g.AllUserMemberIds(), guild.NewS2cNoticeTaskStageUpdateMsg(int32(data.TaskType), int32(g.GetGuildTaskStageIndex(data.TaskType))))
			}
		}
	}

	g.AddDonateRecord(&shared_proto.GuildDonateRecordProto{
		Name:       heroName,
		Sequence:   u64.Int32(donateData.Sequence),
		Times:      u64.Int32(donateData.Times),
		DonateTime: timeutil.Marshal32(ctime),
	}, m.datas.GuildConfig().GuildMaxDonateRecordCount)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.NewS2cDonateMsg(u64.Int32(seq), u64.Int32(donateData.Times), u64.Int32(donateData.Id), u64.Int32(g.GetBuildingAmount()),
		u64.Int32(self.ContributionAmount()), u64.Int32(self.ContributionTotalAmount()),
		u64.Int32(self.ContributionAmount7()), u64.Int32(self.DonationAmount()),
		u64.Int32(self.DonationTotalAmount()), u64.Int32(self.DonationAmount7()),
		u64.Int32(self.DonateTotalYuanbao()))

	broadcastChanged = true

	return
}

func (m *guild_func) upgradeLevel(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("联盟升级，Npc联盟不允许操作")
		errMsg = guild.ErrUpgradeLevelFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpgradeLevel
	}) {
		logrus.Debugf("联盟升级，权限不足")
		errMsg = guild.ErrUpgradeLevelFailDeny
		return
	}

	if !timeutil.IsZero(g.GetUpgradeEndTime()) {
		logrus.Debugf("联盟升级，当前正在升级中")
		errMsg = guild.ErrUpgradeLevelFailUpgrading
		return
	}

	if g.LevelData().NextLevel() == nil {
		logrus.Debugf("联盟升级，当前已经是最高等级")
		errMsg = guild.ErrUpgradeLevelFailMaxLevel
		return
	}

	if g.GetBuildingAmount() < g.LevelData().UpgradeBuilding {
		logrus.Debugf("联盟升级，建设值不足")
		errMsg = guild.ErrUpgradeLevelFailCostNotEnough
		return
	}

	m.startUpgradeLevel(g)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.UPGRADE_LEVEL_S2C
	broadcastChanged = true
	return
}

func (m *guild_func) startUpgradeLevel(g *sharedguilddata.Guild) {
	ctime := m.time.CurrentTime()

	// 扣建设值，开始升级
	g.ReduceBuildingAmount(g.LevelData().UpgradeBuilding, ctime)

	g.SetUpgradeEndTime(ctime.Add(g.LevelData().UpgradeDuration))
	m.updateGuildRankObj(g)

	// 看下能不能升级
	m.tryUpgradeGuildLevel(g, ctime)
}

func (m *guild_func) tryUpgradeGuildLevel(g *sharedguilddata.Guild, ctime time.Time) bool {

	if g.TryUpgradeLevel(ctime) {

		g.TryUpdateTarget(m.datas.GuildConfig(), ctime, shared_proto.GuildTargetType_GuildLevelUp)

		m.tryUpdateFreeNpcGuildChanged(g)

		m.updateGuildRankObj(g)

		g.WalkMember(func(member *sharedguilddata.GuildMember) {
			if member.IsNpc() {
				return
			}

			memSender := m.world.GetUserCloseSender(member.Id())
			if memSender == nil {
				return
			}

			m.heroService.FuncWithSend(member.Id(), func(hero *entity.Hero, result herolock.LockResult) {
				heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_GUILD_LEVEL)
				result.Ok()
			})
		})

		// 大事记
		if logData := m.datas.GuildLogHelp().UpgradeLevel; logData != nil {
			proto := logData.NewLogProto(ctime)
			proto.Text = logData.Text.New().WithLevel(g.LevelData().Level).JsonString()

			m.addGuildLog(g.Id(), proto)
		}

		hctx := heromodule.NewContext(m.dep, operate_type.GuildUpgradeLevel)
		if d := hctx.BroadcastHelp().GuildLevel; d != nil {
			hctx.AddGuildBroadcast(d, g.Id(), 0, g.LevelData().Level, func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickGuildFields(data.KeyGuild, hctx.GetFlagGuildName(g), g.Id())
				text.WithFields(data.KeyNum, g.LevelData().Level)
				return text
			})
		}

		m.updateGuildMemberHeroGuild(g, HeroGuildLevel)

		return true
	}

	return false
}

func (m *guild_func) reduceUpgradeLevelCd(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("联盟升级加速，Npc联盟不允许操作")
		errMsg = guild.ErrReduceUpgradeLevelCdFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpgradeLevelCdr
	}) {
		logrus.Debugf("联盟升级加速，权限不足")
		errMsg = guild.ErrReduceUpgradeLevelCdFailDeny
		return
	}

	upgradeEndTime := g.GetUpgradeEndTime()
	if timeutil.IsZero(upgradeEndTime) {
		logrus.Debugf("联盟升级加速，当前没有在升级")
		errMsg = guild.ErrReduceUpgradeLevelCdFailNoUpgrading
		return
	}

	cdr := g.LevelData().GetCdr(g.GetCdrTimes() + 1)
	if cdr == nil {
		logrus.Debugf("联盟升级加速，已经达到最大加速次数")
		errMsg = guild.ErrReduceUpgradeLevelCdFailMaxTimes
		return
	}

	if g.GetBuildingAmount() < cdr.Cost {
		logrus.Debugf("联盟升级加速，建设值不足")
		errMsg = guild.ErrReduceUpgradeLevelCdFailCostNotEnough
		return
	}

	ctime := m.time.CurrentTime()

	// 扣建设值，开始升级
	g.ReduceBuildingAmount(cdr.Cost, ctime)
	g.IncCdrTimes()

	newUpgradeEndTime := upgradeEndTime.Add(-cdr.CDR)
	g.SetUpgradeEndTime(newUpgradeEndTime)
	m.updateGuildRankObj(g)

	// 立即升级
	m.tickUpdateUpgradeLevel(g, ctime)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.REDUCE_UPGRADE_LEVEL_CD_S2C
	broadcastChanged = true
	return
}

func (m *guild_func) impeachLeader(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.HasChangeLeaderCountDown() {
		logrus.Errorf("联盟弹劾，盟主禅让中")
		errMsg = guild.ErrImpeachLeaderFailChangingLeader
		return
	}

	if g.GetImpeachLeader() != nil {
		logrus.Debugf("联盟弹劾，正在弹劾中")
		errMsg = guild.ErrImpeachLeaderFailImpeachExist
		return
	}

	ctime := m.time.CurrentTime()
	isNpcLeader := g.IsNpcLeader()
	var leaderName string
	if isNpcLeader {
		//if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		//	return permission.ImpeachNpcLeader
		//}) {
		//	logrus.Debugf("联盟弹劾，权限不足")
		//	errMsg = guild.ErrImpeachLeaderFailDeny
		//	return
		//}

		// 新增，Npc联盟不能弹劾，定时开启
		logrus.Debugf("联盟弹劾，Npc联盟不能弹劾")
		errMsg = guild.ErrImpeachLeaderFailDeny
		return
	} else {
		if g.LeaderId() == self.Id() {
			logrus.Debugf("联盟弹劾，弹劾自己？")
			errMsg = guild.ErrImpeachLeaderFailConditionNotReach
			return
		}

		if g.MemberCount() < m.datas.GuildConfig().ImpeachUserLeaderMemberCount {
			logrus.Debugf("联盟弹劾，人数少于10人")
			errMsg = guild.ErrImpeachLeaderFailConditionNotReach
			return
		}

		leaderMember := g.GetMember(g.LeaderId())
		if leaderMember == nil {
			logrus.Errorf("联盟弹劾，leaderMember == nil")
			errMsg = guild.ErrImpeachLeaderFailServerError
			return
		}

		// 离线时间
		leaderSnapshot := m.heroSnapshotService.Get(leaderMember.Id())
		if leaderSnapshot != nil {
			offlineTime := leaderSnapshot.LastOfflineTime
			if timeutil.IsZero(offlineTime) || ctime.Before(offlineTime.Add(m.datas.GuildConfig().ImpeachUserLeaderOffline)) {
				logrus.Debugf("联盟弹劾，盟主离线时间不足")
				errMsg = guild.ErrImpeachLeaderFailConditionNotReach
				return
			}

			leaderName = leaderSnapshot.Name
		} else {
			// 没找到，可以弹劾
			logrus.Errorf("联盟弹劾，盟主snapshot没找到")
		}
	}

	// 弹劾时间
	endTime := m.datas.GuildConfig().GetImpeachLeaderTime(ctime, isNpcLeader)
	if endTime.Before(ctime) {
		logrus.Debugf("联盟弹劾，今日弹劾时间已过")
		errMsg = guild.ErrImpeachLeaderFailInvalidTime
		return
	}

	g.StartImpeachLeader(self.Id(), ctime, endTime, m.datas.GuildConfig().ImpeachExtraCandidateCount)

	g.TryUpdateTarget(m.datas.GuildConfig(), ctime, 0)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.IMPEACH_LEADER_S2C
	broadcastChanged = true

	var selfName string
	if data := m.datas.GuildLogHelp().StartImpeach; data != nil {
		if hero := m.heroSnapshotService.Get(self.Id()); hero != nil {
			proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
			proto.Text = data.Text.New().WithHeroName(hero.Name).JsonString()

			m.addGuildLog(g.Id(), proto)

			selfName = hero.Name
		}
	}

	if len(selfName) > 0 && len(leaderName) > 0 {
		m.pushService.MultiPushFunc(shared_proto.SettingType_ST_GUILD_IMPEACH, g.AllUserMemberIds(), self.Id(), func(d *pushdata.PushData) (title, content string) {
			return d.Title, d.ReplaceContent("{{self}}", selfName, "{{leader}}", leaderName)
		})
	}

	return
}

var voteImpeachend = guild.NewS2cImpeachLeaderVoteMsg(true, []byte{}).Static()

func (m *guild_func) impeachLeaderVote(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, voteTargetId int64) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	impeachLeader := g.GetImpeachLeader()
	if impeachLeader == nil {
		logrus.Debugf("联盟弹劾投票，没有弹劾")
		errMsg = guild.ErrImpeachLeaderVoteFailImpeachNotExist
		return
	}

	if voteTargetId == g.LeaderId() {
		logrus.Debugf("联盟弹劾投票，无法投票给原盟主")
		errMsg = guild.ErrImpeachLeaderVoteFailOldLeader
		return
	}

	if !impeachLeader.IsValidCandidate(voteTargetId) {
		logrus.Debugf("联盟弹劾投票，投票目标不是候选人")
		errMsg = guild.ErrImpeachLeaderVoteFailInvalidTarget
		return
	}

	g.TryVote(self.Id(), voteTargetId)

	m.clearSelfGuildMsgCache(g.Id())

	guildImpeach := g.EncodeImpeachLeader()
	if guildImpeach != nil {
		successMsg = guild.NewS2cImpeachLeaderVoteMsg(false, must.Marshal(guildImpeach))
	} else {
		successMsg = voteImpeachend
	}

	broadcastChanged = true

	m.dep.Tlog().TlogGuildFlowById(self.Id(), operate_type.GuildImpeachLeader.Id(), u64.FromInt64(g.Id()), g.LevelData().Level, u64.FromInt(g.MemberCount()))
	return
}

func (m *guild_func) updateJoinCondition(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember,
	rejectAutoJoin bool, requiredHeroLevel, requiredJunXianLevel, requiredTowerMaxFloor uint64) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("更新入盟条件，Npc联盟不允许修改")
		errMsg = guild.ErrUpdateJoinConditionFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateJoinCondition
	}) {
		logrus.Debugf("更新入盟条件，没有权限")
		errMsg = guild.ErrUpdateJoinConditionFailDeny
		return
	}

	g.SetJoinCondition(rejectAutoJoin, requiredHeroLevel, requiredJunXianLevel, requiredTowerMaxFloor)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.NewS2cUpdateJoinConditionMsg(rejectAutoJoin, u64.Int32(requiredHeroLevel), u64.Int32(requiredJunXianLevel), u64.Int32(requiredTowerMaxFloor))

	broadcastChanged = true
	return
}

func (m *guild_func) updateGuildName(hc iface.HeroController, guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember,
	newName, newFlagName string) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("修改帮派名称，Npc联盟不允许操作")
		errMsg = guild.ErrUpdateGuildNameFailNpc
		return
	}

	if self.ClassLevelData() != m.datas.GuildClassLevelData().MaxKeyData {
		logrus.Debugf("修改帮派名称，没有权限")
		errMsg = guild.ErrUpdateGuildNameFailDeny
		return
	}

	updateName := g.Name() != newName
	updateFlagName := g.FlagName() != newFlagName
	if !updateName && !updateFlagName {
		logrus.Debugf("修改帮派名称，帮派名字和旗号跟现在的一样")
		errMsg = guild.ErrUpdateGuildNameFailExistName
		return
	}

	if updateName && guilds.GetIdByName(newName) != 0 {
		logrus.Debugf("修改帮派名称，帮派名字存在")
		errMsg = guild.ErrUpdateGuildNameFailExistName
		return
	}

	if updateFlagName && guilds.GetIdByFlagName(newFlagName) != 0 {
		logrus.Debugf("修改帮派名称，帮旗名字存在")
		errMsg = guild.ErrUpdateGuildNameFailExistFlagName
		return
	}

	ctime := m.time.CurrentTime()

	// CD中
	if ctime.Before(g.NextChangeNameTime()) {
		logrus.Debugf("修改帮派名称，CD中")
		errMsg = guild.ErrUpdateGuildNameFailCooldown
		return
	}

	//// 查看帮派名字和帮旗名字是否重复
	//for _, otherGuild := range guilds.IdRankArray() {
	//	if g == otherGuild {
	//		continue
	//	}
	//
	//	if otherGuild.Name() == newName {
	//		logrus.Debugf("创建联盟，联盟名字重复")
	//		errMsg = guild.ErrUpdateGuildNameFailExistName
	//		return
	//	}
	//
	//	if otherGuild.FlagName() == newFlagName {
	//		logrus.Debugf("创建联盟，联盟旗号重复")
	//		errMsg = guild.ErrUpdateGuildNameFailExistFlagName
	//		return
	//	}
	//}

	hctx := heromodule.NewContext(m.dep, operate_type.GuildUpdateGuildName)

	// 扣钱
	if !g.IsFreeChangeName() {
		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			if !heromodule.TryReduceCost(hctx, hero, result, m.datas.GuildConfig().ChangeGuildNameCost) {
				logrus.Debugf("修改帮派名称，消耗不足")
				errMsg = guild.ErrUpdateGuildNameFailCostNotEnough
				return
			}

			result.Ok()
		}) {
			return
		}
	}

	oldName := g.Name()
	oldFlagName := g.FlagName()
	guilds.ChangeName(g.Name(), newName, g.Id())
	guilds.ChangeFlagName(g.FlagName(), newFlagName, g.Id())
	nextChangeNameTime := ctime.Add(m.datas.GuildConfig().ChangeGuildNameDuration)
	g.ChangeName(newName, newFlagName, nextChangeNameTime)

	m.updateSnapshot(g)

	m.clearSelfGuildMsgCache(g.Id())

	guildNameChangedMsg := guild.NewS2cUpdateGuildNameBroadcastMsg(i64.Int32(g.Id()), newName, newFlagName).Static()
	m.world.Broadcast(guildNameChangedMsg)

	successMsg = guild.NewS2cUpdateGuildNameMsg(timeutil.Marshal32(g.NextChangeNameTime()))
	//broadcastChanged = true

	if updateName {
		if data := m.datas.GuildLogHelp().UpdateName; data != nil {
			if hero := m.heroSnapshotService.Get(self.Id()); hero != nil {
				proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
				proto.Text = data.Text.New().WithHeroName(hero.Name).WithGuildName(newName).JsonString()

				m.addGuildLog(g.Id(), proto)
			}
		}

		if d := hctx.BroadcastHelp().GuildChangeName; d != nil {
			hctx.AddGuildBroadcast(d, g.Id(), 0, 0, func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithFields(data.KeyText, hctx.GetFlagName(oldFlagName, oldName))
				text.WithClickGuildFields(data.KeyGuild, hctx.GetFlagName(g.FlagName(), g.Name()), g.Id())
				return text
			})
		}
	}

	if updateFlagName {
		if data := m.datas.GuildLogHelp().UpdateFlagName; data != nil {
			if hero := m.heroSnapshotService.Get(self.Id()); hero != nil {
				proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
				proto.Text = data.Text.New().WithHeroName(hero.Name).WithFlagName(newFlagName).JsonString()

				m.addGuildLog(g.Id(), proto)
			}
		}

		if d := hctx.BroadcastHelp().GuildChangeFlag; d != nil {
			hctx.AddGuildBroadcast(d, g.Id(), 0, 0, func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithFields(data.KeyText, hctx.GetFlagName(oldFlagName, oldName))
				text.WithClickGuildFields(data.KeyGuild, hctx.GetFlagName(g.FlagName(), g.Name()), g.Id())
				return text
			})
		}
	}

	m.hebiModule.UpdateGuildInfoBatch(g.AllUserMemberIds(), g.Id())

	m.dep.Guild().UpdateGuildHeroSnapshot(g)

	m.updateGuildRankObj(g)

	return
}

func (m *guild_func) userRequestJoin(hc iface.HeroController, guilds sharedguilddata.Guilds, toJoin int64) (errMsg msg.ErrMsg) {

	toJoinGuild := guilds.Get(toJoin)
	if toJoinGuild == nil {
		logrus.Debugf("申请加入帮派，目标帮派不存在")
		errMsg = guild.ErrUserRequestJoinFailInvalidId
		return
	}

	if toJoinGuild.GetNpcTemplate() != nil && toJoinGuild.GetNpcTemplate().RejectUserJoin {
		logrus.Debugf("申请加入帮派，纯Npc联盟不允许加入")
		errMsg = guild.ErrUserRequestJoinFailNpc
		return
	}

	if toJoinGuild.IsFull() {
		logrus.Debugf("申请加入帮派，目标帮派已满员")
		errMsg = guild.ErrUserRequestJoinFailFull
		return
	}

	ctime := m.time.CurrentTime()
	if m.isInLeaveMemberCd(toJoinGuild, hc.Id(), ctime) {
		logrus.Debugf("申请加入帮派，离开帮派不满4小时，不能加入")
		errMsg = guild.ErrUserRequestJoinFailLeaveCd
		return
	}

	if toJoinGuild.IsRejectAutoJoin() {
		// 不允许自动加

		requestJoinHeroIds := toJoinGuild.GetRequestJoinHeroIds()
		if i64.Contains(requestJoinHeroIds, hc.Id()) {
			logrus.Debugf("申请加入帮派，已经申请过这个联盟")
			errMsg = guild.ErrUserRequestJoinFailDuplicate
			return
		}

		if uint64(len(requestJoinHeroIds)) >= m.datas.GuildConfig().GuildMaxJoinRequestCount {
			// 超出上限，拒绝最老的一个
			if len(requestJoinHeroIds) > 0 {
				removeHeroId := toJoinGuild.RemoveFirstRequestJoinHero()

				// 这里如果锁英雄失败，没改到，则在英雄每日更新的时候改一次
				if !m.heroService.FuncNotError(removeHeroId, func(hero *entity.Hero) (heroChanged bool) {

					hero.RemoveJoinGuildIds(toJoin)
					heroChanged = true
					return
				}) {
					// 成功，发送消息
					m.world.Send(removeHeroId, guild.NewS2cUserRemoveJoinRequestMsg(i64.Int32(toJoin)))
				}
			}

		}
	}

	junXianLevel := m.baizhanService.GetJunXianLevel(hc.Id())
	if junXianLevel < toJoinGuild.GetRequiredJunXianLevel() {
		logrus.Debugf("申请加入帮派，军衔未满足入盟条件")
		errMsg = guild.ErrUserRequestJoinFailCondition
		return
	}

	var joinSucc bool
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.Level() < toJoinGuild.GetRequiredHeroLevel() {
			logrus.Debugf("申请加入帮派，英雄等级未满足入盟条件")
			errMsg = guild.ErrUserRequestJoinFailCondition
			return
		}

		if hero.Tower().HistoryMaxFloor() < toJoinGuild.GetRequiredTowerMaxFloor() {
			logrus.Debugf("申请加入帮派，千重楼层数不够")
			errMsg = guild.ErrUserRequestJoinFailCondition
			return
		}

		if !toJoinGuild.IsRejectAutoJoin() {
			// 允许自动加入帮派
			gctx := &sharedguilddata.GuildContext{
				OperType:     sharedguilddata.JoinGuild,
				OperatorId:   hero.Id(),
				OperatorName: hero.Name(),
			}
			if hero.CountryId() != toJoinGuild.CountryId() { // 如果国家不一样，就先转国
				// 能否转国
				if errMsg = heromodule.CheckChangeCountry(hero, m.datas, ctime); errMsg != nil {
					logrus.Debugf("申请加入帮派，无法转国")
					return
				}
				// 转国前，如果有联盟就必须先退出联盟
				if !m.leaveGuildIfNotLeader(hero, result, guilds, false) {
					logrus.Debugf("申请加入帮派，自己是盟主")
					errMsg = guild.ErrUserRequestJoinFailLeader
					return
				}
				// 转国
				heromodule.ChangeCountryAnyway(m.dep, hero, result, ctime, toJoinGuild.CountryId())
				// 加入新的联盟
				m.joinNewGuild(gctx, hero, result, guilds, toJoinGuild)
				joinSucc = true

			} else if joinSucc = m.joinNewGuildIfNotLeader(gctx, hero, result, guilds, toJoinGuild); !joinSucc {
				logrus.Debugf("申请加入帮派，自己是盟主")
				errMsg = guild.ErrUserRequestJoinFailLeader
				return
			}

			m.dep.Tlog().TlogGuildFlow(hero, operate_type.GuildJoin.Id(), u64.FromInt64(toJoin), toJoinGuild.LevelData().Level, u64.FromInt(toJoinGuild.MemberCount()))
		} else {
			if selfGuild := guilds.Get(hero.GuildId()); selfGuild != nil && selfGuild.IsResistXiongNuDefenders(hero.Id()) && m.xiongNuService.IsStarted(hero.GuildId()) {
				logrus.Debugf("申请加入联盟，自己当前在联盟匈奴入侵防守队伍中，且对方联盟开启了匈奴入侵")
				errMsg = guild.ErrUserRequestJoinFailXiongNuDefender
				return
			}
			if hero.CountryId() != toJoinGuild.CountryId() {
				logrus.Debugf("接受帮派邀请，玩家国家跟联盟不一致")
				errMsg = guild.ErrUserRequestJoinFailCountry

				return
			}

			// 添加申请
			joinGuildRequests := hero.GetJoinGuildIds()
			if uint64(len(joinGuildRequests)) >= m.datas.GuildConfig().UserMaxJoinRequestCount {
				logrus.Debugf("申请加入帮派，达到申请上限")
				errMsg = guild.ErrUserRequestJoinFailSelfFull
				return
			}

			selfGuildId := hero.GuildId()
			if selfGuildId == toJoin {
				logrus.Debugf("申请加入帮派，不能申请自己的联盟")
				errMsg = guild.ErrUserRequestJoinFailSelfGuild
				return
			}

			if i64.Contains(joinGuildRequests, toJoin) {
				hero.RemoveJoinGuildIds(toJoin)
			}

			expiredTime := ctime.Add(m.datas.GuildConfig().JoinRequestDuration)

			// 加给自己
			hero.AddJoinGuildIds(toJoin)
			result.Add(guild.NewS2cUserRequestJoinMsg(i64.Int32(toJoin)))

			// 加到帮派
			toJoinGuild.AddRequestJoinHeroId(hero.Id(), expiredTime)
			m.updateSnapshot(toJoinGuild)

			// 联盟更新
			m.clearSelfGuildMsgCache(toJoin)
			m.world.MultiSend(toJoinGuild.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)

			m.dep.Tlog().TlogGuildFlow(hero, operate_type.GuildRequestJoin.Id(), u64.FromInt64(toJoin), toJoinGuild.LevelData().Level, u64.FromInt(toJoinGuild.MemberCount()))
		}

		result.Changed()
		result.Ok()
	}) {
		if errMsg == nil {
			errMsg = guild.ErrUserRequestJoinFailServerError
		}
		return
	}

	if joinSucc {
		m.trySendJoinChat(hc.Id(), toJoin)

		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			return m.db.UpdateHeroGuildId(ctx, hc.Id(), toJoin)
		}); err != nil {
			logrus.WithError(err).Errorf("保存玩家联盟ID，更新超时 heroId:%v guildId:%v", hc.Id(), toJoin)
		}

		m.hebiModule.UpdateGuildInfo(hc.Id(), toJoin)
	}

	return
}

func (m *guild_func) leaveGuildIfNotLeader(hero *entity.Hero, result herolock.LockResult, guilds sharedguilddata.Guilds, notifyHeroGuildChanged bool) bool {
	if hero.GuildId() != 0 {
		selfGuild := guilds.Get(hero.GuildId())
		if selfGuild != nil {

			if selfGuild.LeaderId() == hero.Id() {
				return false
			}

			m.doRemoveHeroGuild(selfGuild, hero.Id(), hero.Id(), hero, result, notifyHeroGuildChanged)

			m.clearSelfGuildMsgCache(selfGuild.Id())

			m.world.MultiSend(selfGuild.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)
		}
	}

	return true
}

func (m *guild_func) joinNewGuild(gctx *sharedguilddata.GuildContext, hero *entity.Hero, result herolock.LockResult, guilds sharedguilddata.Guilds, toJoinGuild *sharedguilddata.Guild) {
	if hero.GuildId() != 0 {
		return
	}

	memberIds := toJoinGuild.AllUserMemberIds()
	ctime := m.time.CurrentTime()
	newMember := sharedguilddata.NewMember(hero.Id(), m.datas.GuildClassLevelData().MinKeyData, ctime)
	toJoinGuild.AddMember(newMember)

	m.updateSnapshot(toJoinGuild)

	if toJoinGuild.StatueCacheMsg() != nil {
		result.Add(toJoinGuild.StatueCacheMsg())
	}

	// 自己加入新联盟
	m.updateHeroGuild(hero, result, toJoinGuild, newMember, true)

	result.Add(guild.NewS2cUserJoinedMsg(i64.Int32(toJoinGuild.Id()), toJoinGuild.Name(), toJoinGuild.FlagName(), u64.Int32(toJoinGuild.Country().Id)))

	// 清掉自己的所有申请联盟数据
	joinGuildIds := hero.ClearJoinGuildIds()
	if len(joinGuildIds) > 0 {
		result.Add(guild.USER_CLEAR_JOIN_REQUEST_S2C)

		for _, gid := range joinGuildIds {
			g0 := guilds.Get(gid)
			if g0 != nil {
				g0.RemoveRequestJoinHeroId(hero.Id())
			}
		}
	}

	// 加入帮派任务
	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_JOIN_GUILD)

	// 联盟有新人加入
	m.world.MultiSend(memberIds, guild.NewS2cAddGuildMemberMsg(hero.IdBytes(), hero.Name(), hero.Head()).Static())

	// 联盟更新
	m.clearSelfGuildMsgCache(toJoinGuild.Id())
	m.world.MultiSend(memberIds, guild.SELF_GUILD_CHANGED_S2C)

	// 更新空闲npc联盟数量
	m.tryUpdateFreeNpcGuildChanged(toJoinGuild)

	// 联盟日志
	if gctx.OperType == sharedguilddata.ReplyJoinGuild {
		ctime := m.time.CurrentTime()
		if d := m.datas.GuildLogHelp().ReplyJoinGuild; d != nil {
			proto := d.NewHeroLogProto(ctime, hero.IdBytes(), hero.Head())
			proto.Text = d.Text.New().WithHeroName(hero.Name()).WithLeaderName(gctx.OperatorName).JsonString()
			m.addGuildLog(toJoinGuild.Id(), proto)
		}
	} else {
		if d := m.datas.GuildLogHelp().JoinGuild; d != nil {
			proto := d.NewHeroLogProto(ctime, hero.IdBytes(), hero.Head())
			proto.Text = d.Text.New().WithHeroName(hero.Name()).JsonString()

			m.addGuildLog(toJoinGuild.Id(), proto)
		}
	}

	// 发送联盟求助（这个数据变化很快，就不缓存了）
	if msg := m.getListGuildSeekHelpMsg(toJoinGuild); msg != nil {
		result.Add(msg)
	}

	// 加入联盟邮件
	m.trySendFirstJoinGuildMail(hero, toJoinGuild.IsNpcLeader())

	m.clearSelfGuildMsgCache(toJoinGuild.Id())
}

func (m *guild_func) joinNewGuildIfNotLeader(gctx *sharedguilddata.GuildContext, hero *entity.Hero, result herolock.LockResult, guilds sharedguilddata.Guilds, toJoinGuild *sharedguilddata.Guild) bool {
	if !m.leaveGuildIfNotLeader(hero, result, guilds, false) {
		return false
	}
	// 已经离开了原来的联盟
	// 直接加入新的联盟
	m.joinNewGuild(gctx, hero, result, guilds, toJoinGuild)
	return true
}

func (m *guild_func) trySendJoinChat(heroId, guildId int64) {
	text := m.dep.Datas().TextHelp().GuildSayHi.New().JsonString()
	m.chat.SysChat(heroId, guildId, shared_proto.ChatType_ChatGuild, text, shared_proto.ChatMsgType_ChatMsgText, false, true, true, false)
}

func (m *guild_func) trySendFirstJoinGuildMail(hero *entity.Hero, isNpcGuild bool) {
	ctime := m.time.CurrentTime()
	if !hero.IsFirstJoinGuildPrizeCollected() {
		hero.CollectFirstJoinGuildPrize()

		proto := m.datas.MailHelp().GuildFirstJoin.NewTextMail(shared_proto.MailType_MailNormal)
		//proto.Prize = m.datas.GuildConfig().FirstJoinGuildPrizeProto

		m.mail.SendProtoMail(hero.Id(), proto, ctime)
	}

	if isNpcGuild {
		toSendFunc := heromodule.TrySendMailFunc(m.mail, hero, shared_proto.HeroBoolType_BOOL_JOIN_NPC_GUILD,
			m.datas.MailHelp().GuildFirstJoinNpc, ctime)
		if toSendFunc != nil {
			toSendFunc()
		}
	}
}

func (m *guild_func) getListGuildSeekHelpMsg(g *sharedguilddata.Guild) pbutil.Buffer {
	var data [][]byte
	g.RangeSeekHelp(func(proto *shared_proto.GuildSeekHelpProto) (isContinue bool) {
		heroId, _ := idbytes.ToId(proto.HeroId)
		if snapshot := m.heroSnapshotService.Get(heroId); snapshot != nil {
			proto.HeroName = snapshot.Name
			proto.HeroHead = snapshot.Head
		}

		data = append(data, must.Marshal(proto))
		return true
	})
	if len(data) > 0 {
		return guild.NewS2cListGuildSeekHelpMsg(data)
	}
	return nil
}

func (m *guild_func) updateHeroGuild(hero *entity.Hero, result herolock.LockResult, newGuild *sharedguilddata.Guild, selfMember *sharedguilddata.GuildMember, notifyHeroGuildChanged bool) {

	if hero.GetGuildEventPrizeCount() > 0 {
		result.AddFunc(func() pbutil.Buffer {
			return guild.NewS2cRemoveGuildEventPrizeMsg(0)
		})

		var prizeBuilder *resdata.PrizeBuilder

		hero.WalkGuildEventPrize(func(p *entity.HeroGuildEventPrize) {
			hero.RemoveGuildEventPrize(p.Id)

			if prizeBuilder == nil {
				prizeBuilder = resdata.NewPrizeBuilder()
			}

			prizeBuilder.Add(p.Data.Prize.GetPrize())
		})

		if prizeBuilder != nil {
			proto := m.datas.MailHelp().GuildEventPrize.NewTextMail(shared_proto.MailType_MailNormal)
			proto.Prize = prizeBuilder.Build().PrizeProto()
			m.mail.SendProtoMail(hero.Id(), proto, m.time.CurrentTime())
		}
	}

	if newGuild != nil {
		hero.SetGuild(newGuild.Id())
		hero.SetCollectedDailyGuildRankPrize(true)

		if selfMember != nil {
			joinGuildTime := selfMember.GetCreateTime()
			hero.SetJoinGuildTime(joinGuildTime)

			result.AddFunc(func() pbutil.Buffer {
				return guild.NewS2cUpdateHeroJoinGuildTimeMsg(timeutil.Marshal32(joinGuildTime))
			})
		}

		result.Add(guild.NewS2cUpdateFullBigBoxMsg(u64.Int32(newGuild.GetBigBoxData().Id), false, u64.Int32(newGuild.GetBigBoxEnergy())))

		xiongNuInfoMsg := m.xiongNuService.XiongNuInfoMsg(newGuild.Id())
		if xiongNuInfoMsg != nil {
			result.Add(xiongNuInfoMsg)
		}

		m.xiongNuModule.JoinGuild(hero.Id(), newGuild.Id())

		if newGuild.IsFull() {
			m.removeNotFullGuild(newGuild.Id())
		}

		result.Add(guild.NewS2cUpdateHeroGuildMsg(0, newGuild.NewHeroGuildProto()))

		if newGuild.GetWorkshop() == nil {
			result.Add(showWorkshopTrueMsg)
		}

		result.Add(guild.NewS2cUpdateSelfClassLevelMsg(u64.Int32(selfMember.ClassLevelData().Level)))
	} else {
		hero.SetGuild(0)

		ctime := m.time.CurrentTime()
		old := hero.Domestic().SetGuildTechnology(nil)
		heromodule.UpdateBuildingEffect(hero, result, m.datas, ctime, guild_data.GetTechnologyEffects(old)...)

		result.Add(removeHeroGuildMsg)
	}

	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_GUILD_LEVEL)

	if notifyHeroGuildChanged {
		// 注意，这里是持有英雄锁
		if hero.BaseRegion() > 0 {
			r := m.realmService.GetBigMap()
			if r != nil {
				// TODO 优化掉英雄锁的情况
				r.UpdateHeroBasicInfoNoBlock(hero.Id())
			}
		}
	}
}

func (m *guild_func) userCancelRequestJoin(hc iface.HeroController, guilds sharedguilddata.Guilds, cancelGuildId int64) (errMsg msg.ErrMsg) {

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		// 删掉自己的数据
		if !hero.RemoveJoinGuildIds(cancelGuildId) {
			logrus.Debugf("取消申请联盟，申请联盟id不存在")
			errMsg = guild.ErrUserCancelJoinRequestFailInvalidId
			return
		}

		// 删掉帮派中的数据
		g := guilds.Get(cancelGuildId)
		if g != nil {
			if g.RemoveRequestJoinHeroId(hero.Id()) {
				m.clearSelfGuildMsgCache(cancelGuildId)
				m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)
			}
		}

		result.Add(guild.NewS2cUserCancelJoinRequestMsg(i64.Int32(cancelGuildId)))

		result.Changed()
		result.Ok()
	}) {
		if errMsg == nil {
			errMsg = guild.ErrUserCancelJoinRequestFailServerError
		}
		return
	}

	return
}

func (m *guild_func) guildReplyJoinRequest(gctx *sharedguilddata.GuildContext, guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember,
	targetId int64, agree bool) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.AgreeJoin
	}) {
		logrus.Debugf("审批加入联盟，权限不足")
		errMsg = guild.ErrGuildReplyJoinRequestFailDeny
		return
	}

	if g.IsFull() {
		logrus.Debugf("审批加入联盟，联盟已经满员")
		errMsg = guild.ErrGuildReplyJoinRequestFailFullMember
		return
	}

	if g.GetMember(targetId) != nil {
		logrus.Debugf("审批加入联盟，玩家已经在联盟中了")
		errMsg = guild.ErrGuildReplyJoinRequestFailInvalidRequest
		return
	}

	requestJoinHeroIds := g.GetRequestJoinHeroIds()
	if !i64.Contains(requestJoinHeroIds, targetId) {
		logrus.Debugf("审批加入联盟，目标id不在邀请列表中")
		errMsg = guild.ErrGuildReplyJoinRequestFailInvalidRequest
		return
	}

	ctime := m.time.CurrentTime()
	if m.isInLeaveMemberCd(g, targetId, ctime) {
		logrus.Debugf("审批加入联盟，离开帮派不满4小时，不能加入")
		errMsg = guild.ErrGuildReplyJoinRequestFailLeaveCd
		return
	}

	var joinSucc bool
	if m.heroService.FuncWithSend(targetId, func(hero *entity.Hero, result herolock.LockResult) {

		if agree {
			if hero.GuildId() != 0 {
				logrus.Debugf("审批加入联盟，对方已经有联盟")
				errMsg = guild.ErrGuildReplyJoinRequestFailInvalidRequest

				// 既然是盟主，清掉自己的申请帮派数据
				hero.ClearJoinGuildIds()

				g.RemoveRequestJoinHeroId(targetId)
				return
			}

			if hero.CountryId() != g.CountryId() {
				logrus.Debugf("审批加入联盟，玩家国家跟联盟不一致")
				errMsg = guild.ErrGuildReplyJoinRequestFailInvalidRequest

				// 修改掉那个对象
				if hero.RemoveJoinGuildIds(g.Id()) {
					result.Add(guild.NewS2cUserRemoveJoinRequestMsg(i64.Int32(g.Id())))
				}
				g.RemoveRequestJoinHeroId(targetId)

				return
			}

			//if targetGuild := guilds.Get(targetGuildId); targetGuild != nil && targetGuild.IsResistXiongNuDefenders(hero.Id()) && m.xiongNuService.IsStarted(targetGuildId) {
			//	logrus.Debugf("审批加入联盟，目标在联盟匈奴入侵防守队伍中，且对方联盟开启了匈奴入侵")
			//	errMsg = guild.ErrGuildReplyJoinRequestFailXiongNuDefender
			//	return
			//}

			// 同意加入
			if joinSucc = m.joinNewGuildIfNotLeader(gctx, hero, result, guilds, g); !joinSucc {
				logrus.Debugf("审批加入联盟，对方是盟主")
				errMsg = guild.ErrGuildReplyJoinRequestFailInvalidRequest

				// 既然是盟主，清掉自己的申请帮派数据
				hero.ClearJoinGuildIds()

				g.RemoveRequestJoinHeroId(targetId)
				return
			}

			m.dep.Tlog().TlogGuildFlow(hero, operate_type.GuildJoin.Id(), u64.FromInt64(g.Id()), g.LevelData().Level, u64.FromInt(g.MemberCount()))

		} else {
			// 拒绝加入，修改掉那个对象
			if hero.RemoveJoinGuildIds(g.Id()) {
				result.Add(guild.NewS2cUserRemoveJoinRequestMsg(i64.Int32(g.Id())))
			}

			g.RemoveRequestJoinHeroId(targetId)
		}

		successMsg = guild.NewS2cGuildReplyJoinRequestMsg(hero.IdBytes(), agree)
		broadcastChanged = true

		result.Changed()
		result.Ok()
		return
	}) {
		if errMsg == nil {
			errMsg = guild.ErrGuildReplyJoinRequestFailServerError
		}
		return
	}

	if joinSucc {
		m.updateSnapshot(g)

		m.trySendJoinChat(targetId, g.Id())

		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			return m.db.UpdateHeroGuildId(ctx, targetId, g.Id())
		}); err != nil {
			logrus.WithError(err).Errorf("保存玩家联盟ID，更新超时 heroId:%v guildId:%v", targetId, g.Id())
		}

		m.hebiModule.UpdateGuildInfo(targetId, g.Id())
	}

	return
}

func (m *guild_func) guildInvateOther(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember,
	targetId int64) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.InvateOther
	}) {
		logrus.Debugf("邀请玩家加入联盟，权限不足")
		errMsg = guild.ErrGuildInvateOtherFailDeny
		return
	}

	if g.IsFull() {
		logrus.Debugf("邀请玩家加入联盟，联盟已经满员")
		errMsg = guild.ErrGuildInvateOtherFailFullMember
	}

	if g.GetMember(targetId) != nil {
		logrus.Debugf("邀请玩家加入联盟，玩家已经在联盟中了")
		errMsg = guild.ErrGuildInvateOtherFailGuildMember
		return
	}

	ctime := m.time.CurrentTime()

	invateHeroIds := g.GetInvateHeroIds()
	if i64.Contains(invateHeroIds, targetId) {
		logrus.Debugf("邀请玩家加入联盟，目标id已经在邀请列表中")
		errMsg = guild.ErrGuildInvateOtherFailInvated
		return
	}

	if m.heroService.FuncWithSend(targetId, func(hero *entity.Hero, result herolock.LockResult) {

		hero.AddBeenInvateGuildId(g.Id())

		if uint64(len(hero.GetBeenInvateGuildIds())) > m.datas.GuildConfig().UserMaxBeenInvateCount {
			// 将第一个弹出来
			toRemoveId := hero.RemoveFirstBeenInvateGuildId()
			g0 := guilds.Get(toRemoveId)
			if g0 != nil {
				g0.RemoveInvateHero(targetId)
			}
			result.Add(guild.NewS2cUserRemoveBeenInvateGuildMsg(i64.Int32(toRemoveId)))
		}

		result.Add(guild.NewS2cUserAddBeenInvateGuildMsg(i64.Int32(g.Id())))

		m.updateSnapshot(g)
		successMsg = guild.NewS2cGuildInvateOtherMsg(hero.IdBytes())
		broadcastChanged = true

		result.Changed()
		result.Ok()
		return
	}) {
		if errMsg == nil {
			errMsg = guild.ErrGuildInvateOtherFailServerError
		}
		return
	}

	g.AddInvateHero(targetId, ctime.Add(m.datas.GuildConfig().InvateDuration))
	if uint64(len(g.GetInvateHeroIds())) > m.datas.GuildConfig().GuildMaxInvateCount {
		// 超出最大邀请数量，将最早的那个移除掉
		toRemoveHeroId := g.RemoveFirstInvateHeroId()
		rejectGuildId := g.Id()
		m.world.SendFunc(toRemoveHeroId, func() pbutil.Buffer {
			return guild.NewS2cUserReplyInvateRequestMsg(i64.Int32(rejectGuildId), false)
		})
	}

	return
}

func (m *guild_func) guildCancelInvateOther(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember,
	targetId int64) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.InvateOther
	}) {
		logrus.Debugf("取消邀请玩家加入联盟，权限不足")
		errMsg = guild.ErrGuildCancelInvateOtherFailDeny
		return
	}

	invateHeroIds := g.GetInvateHeroIds()
	if !i64.Contains(invateHeroIds, targetId) {
		logrus.Debugf("取消邀请玩家加入联盟，目标id不在邀请列表中")
		errMsg = guild.ErrGuildCancelInvateOtherFailIdNotExist
		return
	}

	if m.heroService.FuncWithSend(targetId, func(hero *entity.Hero, result herolock.LockResult) {

		if hero.RemoveBeenInvateGuildId(g.Id()) {
			result.Add(guild.NewS2cUserRemoveBeenInvateGuildMsg(i64.Int32(g.Id())))
			result.Changed()
		}

		g.RemoveInvateHero(hero.Id())
		m.updateSnapshot(g)
		successMsg = guild.NewS2cGuildInvateOtherMsg(hero.IdBytes())
		broadcastChanged = true

		result.Ok()
		return
	}) {
		if errMsg == nil {
			errMsg = guild.ErrGuildCancelInvateOtherFailServerError
		}
		return
	}

	return
}

func (m *guild_func) userRejectInvateRequest(hc iface.HeroController, guilds sharedguilddata.Guilds,
	rejectGuildId int64) (errMsg msg.ErrMsg) {

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !hero.RemoveBeenInvateGuildId(rejectGuildId) {
			logrus.Debugf("拒绝帮派邀请，帮派id无效")
			errMsg = guild.ErrUserReplyInvateRequestFailInvalidId
			return
		}

		g := guilds.Get(rejectGuildId)
		if g != nil {
			g.RemoveInvateHero(hero.Id())
		}

		result.Add(guild.NewS2cUserReplyInvateRequestMsg(i64.Int32(rejectGuildId), false))

		result.Changed()
		result.Ok()
		return
	}) {
		if errMsg == nil {
			errMsg = guild.ErrUserReplyInvateRequestFailServerError
		}
		return
	}

	return
}

func (m *guild_func) userAgreeInvateRequest(gctx *sharedguilddata.GuildContext, hc iface.HeroController, guilds sharedguilddata.Guilds,
	agreeGuildId int64) (errMsg msg.ErrMsg) {

	var joinSucc bool
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if !i64.Contains(hero.GetBeenInvateGuildIds(), agreeGuildId) {
			logrus.Debugf("接受帮派邀请，帮派没有邀请你")
			errMsg = guild.ErrUserReplyInvateRequestFailInvalidId
			return
		}

		toJoinGuild := guilds.Get(agreeGuildId)
		if toJoinGuild == nil {
			logrus.Debugf("接受帮派邀请，帮派id无效")
			errMsg = guild.ErrUserReplyInvateRequestFailInvalidId

			if hero.RemoveBeenInvateGuildId(agreeGuildId) {
				result.Add(guild.NewS2cUserRemoveBeenInvateGuildMsg(i64.Int32(agreeGuildId)))
				result.Changed()
			}

			return
		}

		if toJoinGuild.IsFull() {
			logrus.Debugf("接受帮派邀请，帮派已满员")
			errMsg = guild.ErrUserReplyInvateRequestFailFullMember

			if hero.RemoveBeenInvateGuildId(agreeGuildId) {
				result.Add(guild.NewS2cUserRemoveBeenInvateGuildMsg(i64.Int32(agreeGuildId)))
				result.Changed()
			}

			toJoinGuild.RemoveInvateHero(hero.Id())
			return
		}

		if toJoinGuild.GetNpcTemplate() != nil && toJoinGuild.GetNpcTemplate().RejectUserJoin {
			logrus.Debugf("接受帮派邀请，不能加入纯Npc帮派")
			errMsg = guild.ErrUserReplyInvateRequestFailInvalidId

			if hero.RemoveBeenInvateGuildId(agreeGuildId) {
				result.Add(guild.NewS2cUserRemoveBeenInvateGuildMsg(i64.Int32(agreeGuildId)))
				result.Changed()
			}

			toJoinGuild.RemoveInvateHero(hero.Id())
			return
		}

		if toJoinGuild.GetMember(hc.Id()) != nil {
			logrus.Debugf("接受帮派邀请，你已经是这个帮派的成员")
			errMsg = guild.ErrUserReplyInvateRequestFailInvalidId

			if hero.RemoveBeenInvateGuildId(agreeGuildId) {
				result.Add(guild.NewS2cUserRemoveBeenInvateGuildMsg(i64.Int32(agreeGuildId)))
				result.Changed()
			}

			toJoinGuild.RemoveInvateHero(hero.Id())
			return
		}

		ctime := m.time.CurrentTime()
		if m.isInLeaveMemberCd(toJoinGuild, hc.Id(), ctime) {
			logrus.Debugf("接受帮派邀请，离开帮派不满4小时，不能加入")
			errMsg = guild.ErrUserReplyInvateRequestFailLeaveCd
			return
		}

		if hero.GuildId() != 0 {
			logrus.Debugf("接受帮派邀请，玩家已经有帮派了")
			errMsg = guild.ErrUserReplyInvateRequestFailInvalidId
			return
		}

		if hero.CountryId() != toJoinGuild.CountryId() { // 玩家国家跟联盟不一致，先转国
			if errMsg = heromodule.ChangeCountry(m.dep, hero, result, ctime, toJoinGuild.CountryId(), true); errMsg != nil {
				logrus.Debugf("接受帮派邀请，转国条件不足")
				return
			}
		}

		joinSucc = true
		m.joinNewGuild(gctx, hero, result, guilds, toJoinGuild)

		hero.RemoveBeenInvateGuildId(agreeGuildId)
		toJoinGuild.RemoveInvateHero(hero.Id())

		result.Add(guild.NewS2cUserReplyInvateRequestMsg(i64.Int32(agreeGuildId), true))

		result.Changed()
		result.Ok()
		return
	}) {
		if errMsg == nil {
			errMsg = guild.ErrUserReplyInvateRequestFailServerError
		}
		return
	}

	if joinSucc {
		m.trySendJoinChat(hc.Id(), agreeGuildId)
		// 从推荐列表中删除
		m.dep.Guild().RemoveRecommendInviteHero(hc.Id())

		if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			return m.db.UpdateHeroGuildId(ctx, hc.Id(), agreeGuildId)
		}); err != nil {
			logrus.WithError(err).Errorf("保存玩家联盟ID，更新超时 heroId:%v guildId:%v", hc.Id(), agreeGuildId)
		}

		m.hebiModule.UpdateGuildInfo(hc.Id(), agreeGuildId)
	}

	return
}

// 定时触发的内容

// 帮派申请过期
func (m *guild_func) tickRemoveExpiredJoinRequest(g *sharedguilddata.Guild, ctime time.Time) {

	removedHeroIds := g.RemoveExpiredRequestJoinHero(ctime)

	if len(removedHeroIds) > 0 {
		toSend := guild.NewS2cUserRemoveJoinRequestMsg(i64.Int32(g.Id())).Static()

		// 这里如果锁英雄失败，没改到，则在英雄每日更新的时候改一次
		for _, removeHeroId := range removedHeroIds {
			if !m.heroService.FuncNotError(removeHeroId, func(hero *entity.Hero) (heroChanged bool) {

				hero.RemoveJoinGuildIds(g.Id())
				heroChanged = true
				return
			}) {
				// 成功，发送消息
				m.world.Send(removeHeroId, toSend)
			}
		}

		m.clearSelfGuildMsgCache(g.Id())

		g.SetChanged()
	}
}

// 帮派邀请过期
func (m *guild_func) tickRemoveExpiredInvateRequest(g *sharedguilddata.Guild, ctime time.Time) {

	removedHeroIds := g.RemoveExpiredInvateHero(ctime)

	if len(removedHeroIds) > 0 {
		toSend := guild.NewS2cUserRemoveBeenInvateGuildMsg(i64.Int32(g.Id())).Static()

		// 这里如果锁英雄失败，没改到，则在英雄每日更新的时候改一次
		for _, removeHeroId := range removedHeroIds {
			if !m.heroService.FuncNotError(removeHeroId, func(hero *entity.Hero) (heroChanged bool) {

				hero.RemoveBeenInvateGuildId(g.Id())
				heroChanged = true
				return
			}) {
				// 成功，发送消息
				m.world.Send(removeHeroId, toSend)
			}
		}

		m.clearSelfGuildMsgCache(g.Id())

		g.SetChanged()
	}

}

// 帮主弹劾
func (m *guild_func) tickUpdateImpeachLeader(g *sharedguilddata.Guild, ctime time.Time) bool {
	oldLeaderMember := g.GetMember(g.LeaderId())

	changed, success := g.TryTickImpeachLeader(ctime, m.datas.GuildClassLevelData().MinKeyData, m.datas.GuildClassLevelData().MaxKeyData)
	if !changed {
		return success
	}

	if success {
		// TODO 弹劾成功，诏书什么的

		// 这里不用更新这个，因为Npc盟主的弹劾，这里不会成功的
		//m.doUpdateFreeNpcGuildChanged()

		// 防御性，取消一下
		g.CancelChangeLeader()

		m.afterImpeachLeader(g, oldLeaderMember, ctime)
		m.updateSnapshot(g)

		// 系统广播
		hctx := heromodule.NewContext(m.dep, operate_type.GuildImpeachLeader)
		if d := hctx.BroadcastHelp().GuildTanHe; d != nil {
			hctx.AddGuildBroadcast(d, g.Id(), 0, 0, func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickGuildFields(data.KeyGuild, hctx.GetFlagGuildName(g), g.Id())

				var heroName string
				heroSnapshot := m.heroSnapshotService.Get(g.LeaderId())
				if heroSnapshot == nil {
					heroName = idbytes.PlayerName(g.LeaderId())
				} else {
					heroName = heroSnapshot.Name
				}
				text.WithClickHeroFields(data.KeyName, heroName, heroSnapshot.Id)
				return text
			})
		}

		m.country.ChangeKing(g.CountryId(), g.LeaderId())
	}

	m.clearSelfGuildMsgCache(g.Id())
	m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)

	g.SetChanged()

	return success
}

// 帮主转让
func (m *guild_func) tickUpdateTransferLeader(g *sharedguilddata.Guild, ctime time.Time) {
	oldLeaderId := g.LeaderId()
	if g.TryTickChangeLeader(ctime, m.datas.GuildClassLevelData().MinKeyData, m.datas.GuildClassLevelData().MaxKeyData) {
		// 帮主转让成功
		// TODO ，诏书什么的

		m.afterDemiseLeader(g, oldLeaderId, g.LeaderId())
		m.updateSnapshot(g)

		m.clearSelfGuildMsgCache(g.Id())
		m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)

		g.SetChanged()

		m.country.ChangeKing(g.CountryId(), g.LeaderId())
	}
}

// 帮派每周重置
func (m *guild_func) resetWeekly(g *sharedguilddata.Guild, resetTime time.Time) {
	compleatedCount := g.GetAllGuildTasksCompletedStageCount(m.datas.GetGuildTaskDataArray())
	var evaluateData *guild_data.GuildTaskEvaluateData
	for _, data := range m.datas.GetGuildTaskEvaluateDataArray() {
		if compleatedCount < data.Complete {
			break
		}
		evaluateData = data
	}
	if g.ResetWeekly(resetTime) {
		m.clearSelfGuildMsgCache(g.Id())
		m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)
		if evaluateData != nil {
			// 您的联盟上周完成联盟任务数量{{num}}，获得[color=#ff304e]{{name}}[/color]评价。您作为联盟{{class_name}}获得以下奖励！
			if d := m.datas.MailHelp().GuildTaskEvaluatePrize; d != nil {
				ctime := m.time.CurrentTime()
				g.WalkMember(func(member *sharedguilddata.GuildMember) {
					// 发联盟任务评价奖励邮件
					proto := d.NewTextMail(shared_proto.MailType_MailNormal)
					proto.Text = d.NewTextFields().WithNum(compleatedCount).WithName(evaluateData.Name).WithClassName(member.ClassLevelData().Name).JsonString()
					proto.Prize = evaluateData.Prizes[member.ClassLevelData().Level-1].Encode()
					m.mail.SendProtoMail(member.Id(), proto, ctime)
				})
			}
		}
		g.SetChanged()
	}
}

// 帮派每日重置
func (m *guild_func) resetDaily(g *sharedguilddata.Guild, resetTime time.Time) {
	if g.ResetDaily(m.datas.GuildConfig().ContributionDay, resetTime) {
		m.clearSelfGuildMsgCache(g.Id())
		m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)
		g.SetChanged()
	}
}

// 帮派升级
func (m *guild_func) tickUpdateUpgradeLevel(g *sharedguilddata.Guild, ctime time.Time) {

	if m.tryUpgradeGuildLevel(g, ctime) {
		// 升级成功
		// 诏书什么的 TODO

		m.clearSelfGuildMsgCache(g.Id())
		m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)

		m.updateSnapshot(g)
		g.SetChanged()
	}

}

type hero_guild_type uint8

const (
	HeroGuildAll    hero_guild_type = 0
	HeroGuildLevel  hero_guild_type = 1
	HeroGuildLeader hero_guild_type = 2
)

func (m *guild_func) updateGuildMemberHeroGuild(g *sharedguilddata.Guild, t hero_guild_type) {
	updateGuildMemberHeroGuild(m.world, g, t)
}

var removeHeroGuildMsg = guild.NewS2cUpdateHeroGuildMsg(0, &shared_proto.HeroGuildProto{}).Static()

func updateGuildMemberHeroGuild(world iface.WorldService, g *sharedguilddata.Guild, t hero_guild_type) {

	switch t {
	case HeroGuildLevel:
		world.MultiSend(g.AllUserMemberIds(), guild.NewS2cUpdateHeroGuildMsg(int32(t), &shared_proto.HeroGuildProto{
			Level: u64.Int32(g.LevelData().Level),
		}))

	case HeroGuildLeader:
		world.MultiSend(g.AllUserMemberIds(), guild.NewS2cUpdateHeroGuildMsg(int32(t), &shared_proto.HeroGuildProto{
			Leader: idbytes.ToBytes(g.LeaderId()),
		}))

	}

}

// 帮派科技升级
func (m *guild_func) tickUpdateUpgradeTechnology(g *sharedguilddata.Guild, ctime time.Time) {

	if m.tryUpgradeTechnology(g, ctime) {
		// 升级成功

		m.clearSelfGuildMsgCache(g.Id())
		m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)

		m.updateSnapshot(g)
		g.SetChanged()
	}

}

// 更新帮派目标
func (m *guild_func) tickUpdateGuildTarget(g *sharedguilddata.Guild, ctime time.Time) {

	shouldUpdate := false
	if !g.IsNpcLeader() && g.GetImpeachLeader() == nil {
		// 盟主超时提示事件
		offlineTime := g.GetLeaderOfflineTime()
		if timeutil.IsZero(offlineTime) {
			// 更新这个值
			offlineTime = ctime
			if leaderSnapshot := m.heroSnapshotService.Get(g.LeaderId()); leaderSnapshot != nil {
				offlineTime = leaderSnapshot.LastOfflineTime
			}

			g.UpdateLeaderOfflineTime(offlineTime)
		}

		if !timeutil.IsZero(offlineTime) && ctime.After(offlineTime.Add(m.datas.GuildConfig().ImpeachUserLeaderOffline)) {
			if g.TryUpdateTarget(m.datas.GuildConfig(), ctime, shared_proto.GuildTargetType_UserLeaderUseless) {
				shouldUpdate = true
			}
		}
	}

	if shouldUpdate {
		m.clearSelfGuildMsgCache(g.Id())
		m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)
	}

}

// 更新帮派目标
func (m *guild_func) tickUpdateGuildChangeCountry(g *sharedguilddata.Guild, ctime time.Time) bool {

	newCountry := g.GetChangeCountryTarget()
	if newCountry == nil {
		return false
	}

	if g.Country() == newCountry {
		logrus.Error("定时检查联盟转国，目标国家跟当前国家一样")
		g.CancelChangeCountry()

		m.clearSelfGuildMsgCache(g.Id())
		m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)
		return false
	}

	if ctime.Unix() < g.GetChangeCountryWaitEndTime() {
		// 时间还没到
		return false
	}

	logrus.Debug("联盟转国，开始执行")

	// 开始操作转国
	g.SetCountry(newCountry)

	g.CancelChangeCountry()

	g.TryUpdateTarget(m.datas.GuildConfig(), ctime, shared_proto.GuildTargetType_GuildChangeCountry)

	m.updateSnapshot(g)

	m.clearSelfGuildMsgCache(g.Id())
	m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)

	m.dep.Mingc().UpdateMsg()
	m.mingcWarService.UpdateMsg()

	return true
}

func (m *guild_func) isInLeaveMemberCd(g *sharedguilddata.Guild, id int64, ctime time.Time) bool {

	leaveTime := g.GetMemverLeaveMemver(id)
	if leaveTime > 0 {
		expireTime := ctime.Add(-m.datas.GuildGenConfig().LeaveAfterJoinDuration).Unix()
		if leaveTime > expireTime {
			return true
		}
	}

	return false
}

// 移除过期的成员离开时间
func (m *guild_func) tickUpdateMemberLeaveTime(g *sharedguilddata.Guild, ctime time.Time) {

	expireTime := ctime.Add(-m.datas.GuildGenConfig().LeaveAfterJoinDuration).Unix()
	g.RemoveExpireMemberLeaveTime(expireTime)

}

// NPC更新帮派职位
func (m *guild_func) npcSetMemberClass(g *sharedguilddata.Guild) {

	if !g.IsNpcLeader() {
		// 玩家帮派不处理
		return
	}

	// 如果有弹劾盟主，先处理弹劾问题
	ctime := m.time.CurrentTime()
	m.tickUpdateImpeachLeader(g, ctime)

	// 再检查一次
	if !g.IsNpcLeader() {
		return
	}

	// 获取到leader对象
	leader := g.GetMaxClassLevelMember()
	if leader == nil {
		logrus.Error("Npc联盟定时更新职位，但是找不到leader")
		return
	}

	// 到点，调整所有的帮派成员的职位
	array := g.GetContribution7Slice()

	var index int
	nextMember := func() *sharedguilddata.GuildMember {

		if index < len(array) {
			member := array[index]
			index++

			if member != leader {
				return member
			}

			if index < len(array) {
				member := array[index]
				index++

				return member
			}
		}

		return nil
	}

	hasUser := false

	// 从头开始设置各个职位
out:
	for _, classLevel := range m.datas.GuildConfig().GetNpcSetClassLevelArray() {

		count := g.LevelData().GetClassMemberCount(classLevel.Level)
		for i := uint64(0); i < count; i++ {

			member := nextMember()
			if member == nil {
				break out
			}
			hasUser = hasUser || !npcid.IsNpcId(member.Id())

			oldClassLevelData := member.ClassLevelData()
			newClassLevelData := classLevel

			if oldClassLevelData != newClassLevelData {
				member.SetClassLevelData(newClassLevelData)

				if !npcid.IsNpcId(member.Id()) {
					// 发职位变更邮件
					m.sendUpdateMemberClassLevelMail(leader, member.Id(), oldClassLevelData, newClassLevelData)
				}
			}
		}
	}

	// 处理原来是官员，后面不是官员的成员
	minClassLevelData := m.datas.GuildClassLevelData().MinKeyData
	for i := 0; i < len(array); i++ {
		member := nextMember()
		if member == nil {
			break
		}
		hasUser = hasUser || !npcid.IsNpcId(member.Id())

		if oldClassLevelData := member.ClassLevelData(); oldClassLevelData != minClassLevelData {
			member.SetClassLevelData(minClassLevelData)

			if !npcid.IsNpcId(member.Id()) {
				// 发职位变更邮件
				m.sendUpdateMemberClassLevelMail(leader, member.Id(), oldClassLevelData, minClassLevelData)
			}
		}
	}

	// 如果有玩家存在，则开启盟主弹劾投票
	if hasUser {
		endTime := m.datas.GuildConfig().GetNextNpcSetClassLevelTime(ctime)
		g.StartImpeachLeader(g.LeaderId(), ctime, endTime, m.datas.GuildConfig().ImpeachExtraCandidateCount)
	}

	g.TryUpdateTarget(m.datas.GuildConfig(), ctime, 0)

	m.clearSelfGuildMsgCache(g.Id())
	m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)

	g.SetChanged()
}

// NPC开始升级帮派
func (m *guild_func) npcTryUpgradeLevel(g *sharedguilddata.Guild) {

	// 如果当前没有升级帮派，而且帮派建设值也足够的情况下，开始升级帮派

	if g.LevelData().NextLevel() == nil {
		// 最高级联盟
		return
	}

	if !timeutil.IsZero(g.GetUpgradeEndTime()) {
		// 帮派升级中
		return
	}

	if g.GetBuildingAmount() < g.LevelData().UpgradeBuilding {
		// 帮派建设值不足
		return
	}

	// 开始升级帮派

	m.startUpgradeLevel(g)

	m.clearSelfGuildMsgCache(g.Id())

	m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)

	g.SetChanged()
}

// NPC踢人
func (m *guild_func) npcTryKickMember(g *sharedguilddata.Guild, ctime time.Time) {

	// 查看是否有符合条件的玩家要被T出帮派
	// 超过24小时不在线的玩家

	if !g.IsFull() {
		// 帮派未满，不踢人
		return
	}

	// 帮主弹劾期间，不踢人
	if g.GetImpeachLeader() != nil {
		return
	}

	// 找到离线时间最长的那个家伙
	var kickMember *sharedguilddata.GuildMember
	g.WalkMember(func(member *sharedguilddata.GuildMember) {
		if member.IsNpc() {
			return
		}

		if member.Id() == g.LeaderId() {
			// 帮主跳过
			return
		}

		if g.IsResistXiongNuDefenders(member.Id()) && m.xiongNuService.IsStarted(g.Id()) {
			// 联盟防守人员，且当前开启了
			return
		}

		snapshot := m.heroSnapshotService.Get(member.Id())
		if snapshot != nil {
			// 离线超过24小时
			offlineTime := snapshot.LastOfflineTime
			if timeutil.IsZero(offlineTime) {
				// 当前在线
				return
			}

			if ctime.Before(offlineTime.Add(m.datas.GuildConfig().NpcKickOfflineDuration)) {
				// 离线时间未达到T人标准
				return
			}
		} else {
			// snapshot 取不到？那赶紧踢了
			logrus.WithField("member id", member.Id()).Errorf("npcTryKickMember，成员 snapshot 没找到!")
		}

		if kickMember == nil {
			kickMember = member
		} else {
			// 7日贡献最低的
			if member.ContributionAmount7() < kickMember.ContributionAmount7() {
				kickMember = member
			} else if member.ContributionAmount7() == kickMember.ContributionAmount7() &&
				member.GetCreateTime().Before(kickMember.GetCreateTime()) {
				kickMember = member
			}
		}
	})

	if kickMember == nil {
		// 没有满足条件的家伙
		return
	}
	if errMsg := m.tryRemoveHeroGuild(g, kickMember.Id(), 0, true,
		guild.ErrKickOtherFailAssembly, guild.ErrKickOtherFailServerError); errMsg != nil {
		logrus.Debugf("NPC帮派定时踢人，踢人失败，" + errMsg.Error())
		return
	}

	// 踢成功
	m.clearSelfGuildMsgCache(g.Id())
	m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)

	// 被npc联盟踢，不加入联盟推荐入盟列表

	g.SetChanged()
}

func (m *guild_func) tryUpdateFreeNpcGuildChanged(g *sharedguilddata.Guild) {
	if g.GetNpcTemplate() == nil || g.GetNpcTemplate().RejectUserJoin {
		// A类和C类联盟
		return
	}

	if !g.IsNpcLeader() {
		return
	}

	m.doUpdateFreeNpcGuildChanged()
}

func (m *guild_func) doUpdateFreeNpcGuildChanged() {
	// 以下情况，更新这个值
	// Npc帮派有人加入
	// Npc帮派升级
	// Npc帮派弹劾成功
	m.freeNpcGuildChanged = true
}

// 定时检查创建npc帮派
func (m *guild_func) tryKeepFreeNpcGuild(guilds sharedguilddata.Guilds) {

	if m.dontCreateNpcGuild {
		// 不再创建npc帮派
		return
	}

	if !m.freeNpcGuildChanged {
		// 没有更新
		return
	}

	// 创建一个新的NPC帮派

	// 循环找出空虚的B盟使用模板个数
	countMap := make(map[uint64]uint64)
	templateUseTimes := make(map[uint64]uint64)
	var freeNpcCount uint64
	guilds.Walk(func(g *sharedguilddata.Guild) {
		t := g.GetNpcTemplate()
		if t == nil || t.RejectUserJoin {
			// A类联盟和C类联盟
			return
		}

		templateUseTimes[t.Id]++

		if m.isFreeNpcGuild(g) {
			freeNpcCount++
			countMap[t.Id]++
		}
	})

	if freeNpcCount >= m.datas.GuildConfig().FreeNpcGuildKeepCount {
		// 空闲帮派超出上限，不创建新的帮派
		m.freeNpcGuildChanged = false
		return
	}

	logrus.Debugf("创建Npc帮派")

	createCount := m.datas.GuildConfig().FreeNpcGuildKeepCount - freeNpcCount
	for i := uint64(0); i < createCount; i++ {

		// 选择模板
		t := m.nextNpcTemplate(countMap, templateUseTimes)
		if t == nil {
			logrus.Debugf("已经找不到模板使用，不再创建联盟")
			m.dontCreateNpcGuild = true
			return
		}
		template := t.template

		guildName, flagName := t.PopName()
		if len(guildName) <= 0 || len(flagName) <= 0 {
			logrus.Errorf("找到的模板有问题，返回了无效的帮派名字和旗号，%s-%s", guildName, flagName)
			m.dontCreateNpcGuild = true
			return
		}

		if m.guildIdGen.Load()+1 > npcid.NpcDataMask {
			logrus.Debugf("没有可用的NPC帮派id，不再创建联盟")
			m.dontCreateNpcGuild = true
			return
		}

		newGuildId := m.newGuildId()
		if newGuildId > npcid.NpcDataMask {
			logrus.Debugf("没有可用的NPC帮派id，不再创建联盟")
			m.dontCreateNpcGuild = true
			return
		}

		if !m.createNpcGuild(guilds, template, newGuildId, guildName, flagName) {
			logrus.Errorf("创建NPC联盟失败")
			break
		}

		// 创建联盟成功
		countMap[template.Id]++
		templateUseTimes[template.Id]++

		logrus.Debugf("创建Npc帮派成功，%s-%s", guildName, flagName)
	}

	m.freeNpcGuildChanged = false
}

func (m *guild_func) nextNpcTemplate(freeCountMap, templateUsingCount map[uint64]uint64) *guild_template {

	// 选择模板
	// 选择剩余“空虚的B盟”数量最少的模板
	// 若有并列，从中选择一个已使用次数最少的模板
	// 若有并列，从中选择顺序号最小的模板

	// 选择模板
	var template *guild_template
	var freeCount uint64  // 当前使用个数
	var usingCount uint64 // 使用次数
	for _, t := range m.templates {

		if t.template.RejectUserJoin {
			// A类联盟
			continue
		}

		// 看看是否还有模板可以用
		if !t.HasName() {
			// 没有名字可以使用了
			continue
		}

		fc := freeCountMap[t.template.Id]
		tuc := templateUsingCount[t.template.Id]
		if template == nil {
			template = t
			freeCount = fc
			usingCount = tuc
			continue
		}

		if fc < freeCount {
			template = t
			freeCount = fc
			usingCount = tuc
			continue
		}

		if fc == freeCount && tuc < usingCount {
			template = t
			freeCount = fc
			usingCount = tuc
			continue
		}
	}

	return template
}

func (m *guild_func) isFreeNpcGuild(g *sharedguilddata.Guild) bool {
	return g.IsNpcLeader() && g.EmptyMemberCount() > m.datas.GuildConfig().FreeNpcGuildEmptyCount
}

// Npc guild
func (m *guild_func) createNpcGuild(guilds sharedguilddata.Guilds, template *guild_data.NpcGuildTemplate,
	newGuildId int64, guildName, flagName string) (success bool) {

	ctime := m.time.CurrentTime()
	newGuild := sharedguilddata.NewNpcGuild(template, newGuildId, guildName, flagName, ctime, m.datas)

	newGuild.TryUpdateTarget(m.datas.GuildConfig(), ctime, 0)

	// 创建帮派
	guildBytes, err := newGuild.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("创建Npc联盟，Guild.Marshal 报错")
		return
	}

	err = ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		return m.db.CreateGuild(ctx, newGuildId, guildBytes)
	})
	if err != nil {
		logrus.WithError(err).Errorf("创建Npc联盟，DB.CreateGuild 报错")
		return
	}

	// 创建成功，添加到guilds
	m.updateSnapshot(newGuild)
	guilds.Add(newGuild)
	m.addNotFullGuild(newGuild)

	if newGuild.GetNpcTemplate() == nil || !newGuild.GetNpcTemplate().RejectUserJoin {
		m.updateGuildRankObj(newGuild)
	}

	return true

}

func (m *guild_func) updateFriendGuild(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, text string) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("修改友盟，Npc联盟不允许操作")
		errMsg = guild.ErrUpdateFriendGuildFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateFriendGuild
	}) {
		logrus.Debugf("修改友盟，没有权限")
		errMsg = guild.ErrUpdateFriendGuildFailDeny
		return
	}

	g.SetFriendGuildText(text)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.NewS2cUpdateFriendGuildMsg(text)
	return
}

func (m *guild_func) updateEnemyGuild(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, text string) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("修改敌盟，Npc联盟不允许操作")
		errMsg = guild.ErrUpdateEnemyGuildFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateEnemyGuild
	}) {
		logrus.Debugf("修改敌盟，没有权限")
		errMsg = guild.ErrUpdateEnemyGuildFailDeny
		return
	}

	g.SetEnemyGuildText(text)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.NewS2cUpdateEnemyGuildMsg(text)
	return
}

//func (m *guild_func) updateGuildPrestige(hc iface.HeroController, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, target *country.CountryData) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
//
//	if g.IsNpcLeader() {
//		logrus.Debugf("修改声望目标，Npc联盟不允许操作")
//		errMsg = guild.ErrUpdateGuildPrestigeFailNpc
//		return
//	}
//
//	// 检查权限
//	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
//		return permission.UpdatePrestigeTarget
//	}) {
//		logrus.Debugf("修改声望目标，没有权限")
//		errMsg = guild.ErrUpdateGuildPrestigeFailDeny
//		return
//	}
//
//	ctime := m.time.CurrentTime()
//	// 检查CD
//	if ctime.Before(g.GetNextUpdatePrestigeTargetTime()) {
//		logrus.Debugf("修改声望目标，CD中")
//		errMsg = guild.ErrUpdateGuildPrestigeFailCountdown
//		return
//	}
//
//	oldTarget := g.Country()
//	if oldTarget == target {
//		logrus.Debugf("修改声望目标，新目标跟当前目标一致")
//		errMsg = guild.ErrUpdateGuildPrestigeFailSameTarget
//		return
//	}
//
//	hctx := heromodule.NewContext(m.dep, operate_type.GuildChangePrestigeTarget)
//
//	// 扣钱
//	if m.datas.GuildConfig().UpdateCountryCost != nil {
//		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
//			if !heromodule.TryReduceCost(hctx, hero, result, m.datas.GuildConfig().UpdateCountryCost) {
//				logrus.Debugf("修改声望目标，消耗不足")
//				errMsg = guild.ErrUpdateGuildPrestigeFailCostNotEnough
//				return
//			}
//			result.Ok()
//		}) {
//			if errMsg == nil {
//				errMsg = guild.ErrUpdateGuildPrestigeFailServerError
//			}
//			return
//		}
//
//		if errMsg != nil {
//			return
//		}
//	}
//
//	oldTargetPrestige := g.GetPrestige()
//	g.SetCountry(target)
//
//	// 扣声望
//	newPrestige := u64.Sub(g.GetPrestige(), u64.MultiCoef(g.GetPrestige(), m.datas.GuildConfig().UpdateCountryLostPrestigeCoef))
//	g.SetPrestige(newPrestige)
//
//	// 旧退出国家减声望
//	m.dep.Country().ReducePrestige(oldTarget.Id, oldTargetPrestige)
//
//	// 新加入国家加声望
//	m.dep.Country().AddPrestige(target.Id, newPrestige)
//
//	// 加CD
//	g.SetNextUpdatePrestigeTargetTime(ctime.Add(m.datas.GuildConfig().UpdateCountryDuration))
//
//	m.clearSelfGuildMsgCache(g.Id())
//
//	successMsg = guild.NewS2cUpdateGuildPrestigeMsg(u64.Int32(oldTarget.Id), u64.Int32(target.Id))
//	broadcastChanged = true
//
//	if data := m.datas.GuildLogHelp().UpdateCountry; data != nil {
//		if hero := m.heroSnapshotService.Get(self.Id()); hero != nil {
//			proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
//			proto.Text = data.Text.New().WithHeroName(hero.Name).WithCountry(oldTarget.Name).WithCountryName(target.Name).JsonString()
//
//			m.addGuildLog(g.Id(), proto)
//		}
//	}
//
//	if d := hctx.BroadcastHelp().GuildChangeCountry; d != nil {
//		hctx.AddGuildBroadcast(d, g.Id(), 0, 0, func() *i18n.Fields {
//			text := d.NewTextFields()
//			text.WithClickGuildFields(data.KeyGuild, hctx.GetFlagGuildName(g), g.Id())
//			text.WithFields(data.KeyText, oldTarget.Name).WithFields(data.KeyCountry, target.Name)
//			return text
//		})
//	}
//
//	m.updateSnapshot(g)
//
//	m.updateGuildRankObj(g)
//	return
//}

func (m *guild_func) placeGuildStatue(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, realm iface.Realm) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.LeaderId() != self.Id() {
		logrus.Debugf("放置雕像，不是盟主")
		errMsg = guild.ErrPlaceGuildStatueFailNotLeader
		return
	}

	realmId := g.Statue()
	if realmId != 0 {
		logrus.Debugf("放置雕像，当前有雕像了")
		errMsg = guild.ErrPlaceGuildStatueFailHasPlaced
		return
	}

	g.PlaceStatue(realm.Id())

	m.updateSnapshot(g)
	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.PLACE_GUILD_STATUE_S2C

	m.world.MultiSend(g.AllUserMemberIds(), g.StatueCacheMsg())

	//snapshot := m.heroSnapshotService.Get(self.Id())
	//if snapshot != nil {
	//	level, sequence, _ := realmface.ParseRealmId(realm.Id())
	//	g.AddBigEvent(&shared_proto.GuildBigEventProto{
	//		Time: timeutil.Marshal32(m.time.CurrentTime()),
	//		Type: shared_proto.GuildBigEventType_Statue,
	//		Statue: &shared_proto.StatueProto{
	//			LeaderId:    snapshot.IdBytes,
	//			LeaderName:  snapshot.Name,
	//			RegionLevel: u64.Int32(level),
	//			RegionId:    u64.Int32(sequence),
	//		},
	//	}, m.datas.GuildConfig().GuildMaxBigEventCount)
	//}

	return
}

func (m *guild_func) takeBackGuildStatue(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.LeaderId() != self.Id() {
		logrus.Debugf("放置雕像，不是盟主")
		errMsg = guild.ErrPlaceGuildStatueFailNotLeader
		return
	}

	realmId := g.Statue()
	if realmId == 0 {
		logrus.Debugf("放置雕像，当前没有雕像")
		errMsg = guild.ErrPlaceGuildStatueFailHasPlaced
		return
	}

	g.TakeBackStatue()

	m.updateSnapshot(g)
	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.TAKE_BACK_GUILD_STATUE_S2C

	m.world.MultiSend(g.AllUserMemberIds(), guild.TAKE_BACK_GUILD_STATUE_S2C_BROADCAST)
	return
}

func (m *guild_func) afterImpeachLeader(g *sharedguilddata.Guild, oldLeaderMember *sharedguilddata.GuildMember, ctime time.Time) {
	g.TryUpdateTarget(m.datas.GuildConfig(), ctime, 0)

	m.updateGuildMemberHeroGuild(g, HeroGuildLeader)

	//if oldLeaderMember != nil {
	//var oldLeaderIdBytes []byte
	//var oldLeaderName string
	//
	//if oldLeaderMember.IsNpc() {
	//	oldLeaderSnapshot := oldLeaderMember.NpcProto()
	//	if oldLeaderSnapshot == nil {
	//		logrus.WithField("member", oldLeaderMember.Id()).Errorln("联盟成员是个npc，但是呢，他的 npc proto为空")
	//		return
	//	}
	//
	//	oldLeaderIdBytes = oldLeaderSnapshot.Basic.GetId()
	//	oldLeaderName = oldLeaderSnapshot.Basic.GetName()
	//} else {
	//	oldLeader := m.heroSnapshotService.Get(oldLeaderMember.Id())
	//	if oldLeader == nil {
	//		return
	//	}
	//
	//	oldLeaderIdBytes = oldLeader.IdBytes
	//	oldLeaderName = oldLeader.Name
	//}

	newLeader := m.heroSnapshotService.Get(g.LeaderId())
	if newLeader == nil {
		return
	}

	if data := m.datas.GuildLogHelp().NewLeaderImpeach; data != nil {
		ctime := m.time.CurrentTime()

		proto := data.NewHeroLogProto(ctime, newLeader.IdBytes, newLeader.Head)
		proto.Text = data.Text.New().WithHeroName(newLeader.Name).JsonString()

		m.addGuildLog(g.Id(), proto)
	}

	if mailData := m.datas.MailHelp().GuildLeaderGenerated; mailData != nil {
		proto := mailData.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Text = mailData.NewTextFields().WithName(newLeader.Name).JsonString()
		m.mail.SendReportMail(newLeader.Id, proto, ctime)
	}

	//}
}

func (m *guild_func) afterDemiseLeader(g *sharedguilddata.Guild, oldId, newId int64) {

	m.updateGuildMemberHeroGuild(g, HeroGuildLeader)

	// 转让只有玩家才干的出来
	oldLeader := m.heroSnapshotService.Get(oldId)
	if oldLeader == nil {
		logrus.Error("guild.afterDemiseLeader oldLeader == nil")
		return
	}

	newLeader := m.heroSnapshotService.Get(g.LeaderId())
	if newLeader == nil {
		logrus.Error("guild.afterDemiseLeader newLeader == nil")
		return
	}

	ctime := m.time.CurrentTime()
	if data := m.datas.GuildLogHelp().NewLeaderDemise; data != nil {
		proto := data.NewHeroLogProto(ctime, newLeader.IdBytes, newLeader.Head)
		proto.Text = data.Text.New().WithHeroName(oldLeader.Name).WithLeader(newLeader.Name).JsonString()

		m.addGuildLog(g.Id(), proto)
	}

	if data := m.datas.MailHelp().GuildLeaderChanged; data != nil {
		proto := data.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Text = data.NewTextFields().WithName(oldLeader.Name).WithLeader(newLeader.Name).JsonString()
		m.mail.SendProtoMail(oldId, proto, ctime)
		m.mail.SendProtoMail(newId, proto, ctime)
	}

}

func (m *guild_func) upgradeTechnology(hc iface.HeroController, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, firstLevel *guild_data.GuildTechnologyData) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	// 检查权限
	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpgradeTechnology
	}) {
		logrus.Debugf("联盟升级科技，没有权限")
		errMsg = guild.ErrUpgradeTechnologyFailDeny
		return
	}

	// 科技升级中
	if g.GetTechUpgradeData() != nil && g.GetTechUpgradeEndTime().Unix() > 0 {
		logrus.Debugf("联盟升级科技，当前有科技正在升级")
		errMsg = guild.ErrUpgradeTechnologyFailUpgrading
		return
	}

	nextLevel := firstLevel
	if currentLevel := g.GetTechnology(firstLevel.Group); currentLevel != nil {
		// 找第一级
		nextLevel = currentLevel.GetNextLevel()
		if nextLevel == nil {
			logrus.Debugf("联盟升级科技，已经是最高等级")
			errMsg = guild.ErrUpgradeTechnologyFailMaxLevel
			return
		}
	}

	if g.LevelData().Level < nextLevel.RequireGuildLevel {
		logrus.Debugf("联盟升级科技，已经是最高等级")
		errMsg = guild.ErrUpgradeTechnologyFailRequired
		return
	}

	if g.GetBuildingAmount() < nextLevel.UpgradeBuilding {
		logrus.Debugf("联盟升级科技，建设值不足")
		errMsg = guild.ErrUpgradeTechnologyFailCostNotEnough
		return
	}

	ctime := m.time.CurrentTime()
	g.ReduceBuildingAmount(nextLevel.UpgradeBuilding, ctime) // 扣建设值

	endTime := ctime.Add(nextLevel.UpgradeDuration)

	g.UpgradeTechnology(nextLevel, endTime)

	m.clearSelfGuildMsgCache(g.Id())

	// 广播所有
	m.world.MultiSend(g.AllUserMemberIds(), pushTechHelpableTrue)

	// 在联盟中的人可以协助
	g.WalkMember(func(member *sharedguilddata.GuildMember) {
		member.SetIsTechHelpable(true)
	})

	successMsg = guild.NewS2cUpgradeTechnologyMsg(u64.Int32(nextLevel.Group), timeutil.Marshal32(endTime))
	broadcastChanged = true
	return
}

func (m *guild_func) tryUpgradeTechnology(g *sharedguilddata.Guild, ctime time.Time) bool {

	if ok, toUpgrade := g.TryUpgradeTechnology(ctime); ok {
		if toUpgrade != nil {

			memberIds := g.AllUserMemberIds()
			if toUpgrade.Effect != nil {
				// 给联盟成员更新联盟效果，异步来做
				technologys := g.GetEffectTechnology()

				for _, id := range memberIds {
					m.world.FuncHero(id, func(id int64, hc iface.HeroController) {
						hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
							hero.Domestic().SetGuildTechnology(technologys)
							heromodule.UpdateBuildingEffect(hero, result, m.datas, ctime, toUpgrade.Effect)
							result.Ok()
						})
					})
				}
			}

			if toUpgrade.BigBox != nil {
				// 当前不能领取宝箱的人收到这个消息
				toSend := guild.NewS2cUpdateFullBigBoxMsg(u64.Int32(g.GetBigBoxData().Id), false, u64.Int32(g.GetBigBoxEnergy())).Static()
				for _, memberId := range memberIds {
					if !g.IsFullBigBoxMember(memberId) {
						// 更新当前的逼格宝箱
						m.world.Send(memberId, toSend)
					}
				}
			}

			if data := m.datas.GuildLogHelp().UpgradeTechnology; data != nil {
				ctime := m.time.CurrentTime()
				proto := data.NewLogProto(ctime)
				proto.Text = data.Text.New().WithTechName(toUpgrade.Name).WithLevel(toUpgrade.Level).JsonString()

				m.addGuildLog(g.Id(), proto)
			}
		}

		// 遍历所有成员
		g.WalkMember(func(member *sharedguilddata.GuildMember) {
			if !member.IsNpc() && member.GetIsTechHelpable() {
				member.SetIsTechHelpable(false)

				m.world.Send(member.Id(), pushTechHelpableFalse)
			}
		})

		return true
	}

	return false
}

var pushTechHelpableTrue = guild.NewS2cPushTechHelpableMsg(true).Static()
var pushTechHelpableFalse = guild.NewS2cPushTechHelpableMsg(false).Static()

func (m *guild_func) reduceUpgradeTechnologyCd(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if g.IsNpcLeader() {
		logrus.Debugf("联盟科技升级加速，Npc联盟不允许操作")
		errMsg = guild.ErrReduceTechnologyCdFailNpc
		return
	}

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpgradeTechnologyCdr
	}) {
		logrus.Debugf("联盟科技升级加速，权限不足")
		errMsg = guild.ErrReduceTechnologyCdFailDeny
		return
	}

	toUpgrade := g.GetTechUpgradeData()
	if toUpgrade == nil {
		logrus.Debugf("联盟科技升级加速，没有升级中的科技")
		errMsg = guild.ErrReduceTechnologyCdFailNoUpgrading
		return
	}

	upgradeEndTime := g.GetTechUpgradeEndTime()
	if timeutil.IsZero(upgradeEndTime) {
		logrus.Debugf("联盟科技升级加速，当前没有在升级")
		errMsg = guild.ErrReduceTechnologyCdFailNoUpgrading
		return
	}

	cdr := toUpgrade.GetCdr(g.GetTechCdrTimes() + 1)
	if cdr == nil {
		logrus.Debugf("联盟科技升级加速，已经达到最大加速次数")
		errMsg = guild.ErrReduceTechnologyCdFailMaxTimes
		return
	}

	if g.GetBuildingAmount() < cdr.Cost {
		logrus.Debugf("联盟科技升级加速，建设值不足")
		errMsg = guild.ErrReduceTechnologyCdFailCostNotEnough
		return
	}

	ctime := m.time.CurrentTime()

	// 扣建设值，开始升级
	g.ReduceBuildingAmount(cdr.Cost, ctime)

	g.IncTechCdrTimes()
	newUpgradeEndTime := upgradeEndTime.Add(-cdr.CDR)
	g.SetTechUpgradeEndTime(newUpgradeEndTime)

	// 立即升级
	m.tryUpgradeTechnology(g, ctime)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.REDUCE_TECHNOLOGY_CD_S2C
	broadcastChanged = true
	return
}

func (m *guild_func) helpTech(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if !self.GetIsTechHelpable() {
		logrus.Debugf("联盟科技协助，你不能协助")
		errMsg = guild.ErrHelpTechFailCantHelp
		return
	}

	toUpgrade := g.GetTechUpgradeData()
	if toUpgrade == nil {
		logrus.Debugf("联盟科技协助，没有升级中的科技")
		errMsg = guild.ErrHelpTechFailNoTechUpgrading
		return
	}

	if toUpgrade.HelpCdr <= 0 {
		logrus.Debugf("联盟科技协助，HelpCdr = 0")
		errMsg = guild.ErrHelpTechFailCantHelp
		return
	}

	upgradeEndTime := g.GetTechUpgradeEndTime()
	if timeutil.IsZero(upgradeEndTime) {
		logrus.Debugf("联盟科技协助，当前没有在升级")
		errMsg = guild.ErrHelpTechFailNoTechUpgrading
		return
	}

	self.SetIsTechHelpable(false)

	ctime := m.time.CurrentTime()

	newUpgradeEndTime := upgradeEndTime.Add(-toUpgrade.HelpCdr)
	g.SetTechUpgradeEndTime(newUpgradeEndTime)

	// 立即升级
	m.tryUpgradeTechnology(g, ctime)

	m.clearSelfGuildMsgCache(g.Id())

	successMsg = guild.NewS2cHelpTechMsg(u64.Int32(toUpgrade.Id))
	broadcastChanged = true
	return
}

func (m *guild_func) GmTryKickMember(g *sharedguilddata.Guild, heroId int64) bool {
	return m.tryRemoveHeroGuild(g, heroId, heroId, true, guild.ErrKickOtherFailAssembly, guild.ErrKickOtherFailServerError) == nil
}

func (m *guild_func) updateMark(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, mark *shared_proto.GuildMarkProto) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.UpdateMark
	}) {
		logrus.Debugf("更新联盟标记，权限不足")
		errMsg = guild.ErrUpdateGuildMarkFailDeny
		return
	}

	successMsg = guild.NewS2cUpdateGuildMarkMsg(mark).Static()

	isRemove := mark.PosX == 0 && mark.PosY == 0 && len(mark.Msg) == 0
	if isRemove {
		g.RemoveMark(int(mark.Index - 1))
	} else {
		g.AddMark(mark, successMsg)
		// 记录日志
		if logData := m.datas.GuildLogHelp().UpdateMark; logData != nil {
			hero := m.heroSnapshotService.Get(self.Id())
			if hero != nil {
				proto := logData.NewHeroLogProto(m.time.CurrentTime(), hero.IdBytes, hero.Head)
				proto.Text = logData.Text.New().WithHeroName(hero.Name).WithPosX(mark.PosX).WithPosY(mark.PosY).WithText(mark.Msg).JsonString()
				m.dep.Guild().AddLog(g.Id(), proto)
			}
		}
	}

	m.world.MultiSendIgnore(g.AllUserMemberIds(), successMsg, self.Id())
	return
}

func (m *guild_func) sendYinliangToMem(g *sharedguilddata.Guild, sender *sharedguilddata.GuildMember, receiverId int64, toSend uint64) (succMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
	if !sender.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.ChangeYinliang
	}) {
		errMsg = guild.ErrSendYinliangToMemberFailDeny
		return
	}

	receiver := g.GetMember(receiverId)
	if receiver == nil {
		errMsg = guild.ErrSendYinliangToMemberFailNoMember
		return
	}

	if _, reduced := g.ReduceYinliang(toSend); !reduced {
		errMsg = guild.ErrSendYinliangToMemberFailNotEnough
		return
	}

	ctime := m.dep.Time().CurrentTime()

	receiver.AddHistorySalary(toSend)

	if d := m.datas.GuildLogHelp().YinliangGuildSendMember; d != nil {
		senderName := m.dep.HeroSnapshot().GetHeroName(sender.Id())
		receiverName := m.dep.HeroSnapshot().GetHeroName(receiver.Id())
		text := d.Text.New().WithLeader(senderName).WithTarget(receiverName).WithYinliang(toSend).JsonString()
		m.addYinliangRecord(g, text, ctime, d)
	}

	succMsg = guild.NewS2cSendYinliangToMemberMsg(idbytes.ToBytes(receiverId), u64.Int32(toSend))
	broadcastChanged = true

	m.clearSelfGuildMsgCache(g.Id())

	if d := m.dep.Datas().MailHelp().GuildSendYinliangToMember; d != nil {
		mail := d.NewTextMail(shared_proto.MailType_MailNormal)
		leaderName := m.dep.HeroSnapshot().GetHeroName(sender.Id())
		mail.Text = d.NewTextFields().WithLeader(leaderName).WithAmount(toSend).JsonString()
		prize := &shared_proto.PrizeProto{}
		prize.Yinliang = u64.Int32(toSend)
		mail.Prize = prize
		m.mail.SendProtoMail(receiverId, mail, ctime)
	}

	return
}

func (m *guild_func) paySalary(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (succMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.ChangeYinliang
	}) {
		errMsg = guild.ErrPaySalaryFailDeny
		return
	}

	var allSalary uint64
	g.WalkMember(func(member *sharedguilddata.GuildMember) {
		allSalary += member.Salary()
	})
	if _, succ := g.ReduceYinliang(allSalary); !succ {
		errMsg = guild.ErrPaySalaryFailNotEnough
		return
	}

	ctime := m.dep.Time().CurrentTime()

	leaderName := m.dep.HeroSnapshot().GetHeroName(g.LeaderId())
	leaderFlagName := m.toFlagHeroName(g.FlagName(), leaderName)
	g.WalkMember(func(member *sharedguilddata.GuildMember) {
		if member.Salary() <= 0 {
			return
		}

		member.AddHistorySalary(member.Salary())

		if d := m.dep.Datas().MailHelp().GuildPaySalary; d != nil {
			mail := d.NewTextMail(shared_proto.MailType_MailNormal)
			mail.Text = d.NewTextFields().WithLeader(leaderFlagName).JsonString()
			prize := &shared_proto.PrizeProto{}
			prize.Yinliang = u64.Int32(member.Salary())
			mail.Prize = prize
			m.mail.SendProtoMail(member.Id(), mail, ctime)
		}
	})

	if d := m.datas.GuildLogHelp().YinliangGuildPaySalary; d != nil {
		text := d.Text.New().WithLeader(leaderName).WithYinliang(allSalary).JsonString()
		m.addYinliangRecord(g, text, ctime, d)
	}

	succMsg = guild.NewS2cPaySalaryMsg(u64.Int32(allSalary))
	broadcastChanged = true

	m.clearSelfGuildMsgCache(g.Id())

	return
}

func (m *guild_func) sendYinliangToGuild(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, receiverId int64, toSend uint64) (succMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.ChangeYinliang
	}) {
		errMsg = guild.ErrSendYinliangToOtherGuildFailDeny
		return
	}

	if _, reduced := g.ReduceYinliang(toSend); !reduced {
		errMsg = guild.ErrSendYinliangToOtherGuildFailNotEnough
		return
	}

	receiver := guilds.Get(receiverId)
	if receiver == nil {
		errMsg = guild.ErrSendYinliangToOtherGuildFailNoGuild
		return
	}

	receiver.AddYinliang(toSend)
	succMsg = guild.NewS2cSendYinliangToOtherGuildMsg(int32(receiverId), u64.Int32(toSend))
	broadcastChanged = true

	m.clearSelfGuildMsgCache(g.Id())
	m.clearSelfGuildMsgCache(receiverId)

	// 联盟赠送银两记录
	ctime := m.dep.Time().CurrentTime()
	sendProto := receiver.SendYinliangToMe[g.Id()]
	if sendProto == nil {
		sendProto = &shared_proto.GuildYinliangSendProto{}
		receiver.SendYinliangToMe[g.Id()] = sendProto
	}
	sendProto.Time = timeutil.Marshal32(ctime)
	sendProto.Send = u64.Int32(toSend)
	sendProto.WeeklySend += u64.Int32(toSend)
	sendProto.AllSend += u64.Int32(toSend)
	g.UpdateYinliangSendToGuild(receiver.NewBasicProto(), sendProto)

	// 联盟日志
	if d := m.datas.GuildLogHelp().YinliangGuildSend; d != nil {
		text := d.Text.New().WithTarget(m.FlagGuildName(receiver)).WithYinliang(toSend).JsonString()
		m.addYinliangRecord(g, text, ctime, d)
	}
	if d := m.datas.GuildLogHelp().YinliangGuildReceive; d != nil {
		leaderName := m.dep.HeroSnapshot().GetHeroName(g.LeaderId())
		text := d.Text.New().WithSender(leaderName).WithYinliang(toSend).JsonString()
		m.addYinliangRecord(g, text, ctime, d)
	}

	// 邮件
	if !receiver.IsNpcLeader() {
		leaderName := m.dep.HeroSnapshot().GetHeroName(g.LeaderId())
		guildName := m.FlagGuildName(g)
		proto := m.datas.MailHelp().GuildReceiveYinliangFromGuild.NewTextMail(shared_proto.MailType_MailNormal)
		proto.Text = m.datas.MailHelp().GuildReceiveYinliangFromGuild.NewTextFields().WithGuild(guildName).WithLeader(leaderName).WithAmount(toSend).JsonString()
		m.mail.SendProtoMail(receiver.LeaderId(), proto, m.time.CurrentTime())
	}

	return
}

func (m *guild_func) SetSalary(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, targetId int64, toSet uint64) (succMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.ChangeYinliang
	}) {
		errMsg = guild.ErrSetSalaryFailDeny
		return
	}

	target := g.GetMember(targetId)
	if target == nil {
		errMsg = guild.ErrSetSalaryFailNoMember
		return
	}
	target.SetSalary(toSet)
	succMsg = guild.NewS2cSetSalaryMsg(idbytes.ToBytes(targetId), u64.Int32(toSet))
	return
}

func (m *guild_func) FlagGuildName(g *sharedguilddata.Guild) string {
	return m.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(g.FlagName(), g.Name())
}

func (m *guild_func) toFlagHeroName(flagName, heroName string) string {
	return m.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(flagName, heroName)
}

var conveneAllMsg = guild.NewS2cConveneMsg(nil).Static()

func (m *guild_func) convene(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, targetId int64,
	selfName string, selfBaseX, selfBaseY int) (succMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.ConveneMember
	}) {
		errMsg = guild.ErrConveneFailDeny
		return
	}

	data := m.datas.MailHelp().GuildConveneOfficer
	if g.LeaderId() == self.Id() {
		data = m.datas.MailHelp().GuildConveneLeader
	}

	posString := fmt.Sprintf("(%d, %d)", selfBaseX, selfBaseY)

	ctime := m.time.CurrentTime()
	if targetId == 0 {

		g.WalkMember(func(member *sharedguilddata.GuildMember) {
			if member == self {
				return
			}

			if ctime.Before(member.GetNextConveneTime()) {
				return
			}
			member.SetNextConveneTime(ctime.Add(m.datas.GuildGenConfig().ConveneCooldown))

			// 发召集邮件
			if data != nil {
				proto := data.NewTextMail(shared_proto.MailType_MailGuildCall)
				proto.Text = data.NewTextFields().WithName(selfName).WithPos(posString).JsonString()
				proto.GuildCallHeroId = idbytes.ToBytes(self.Id())
				proto.GuildCallPosX = int32(selfBaseX)
				proto.GuildCallPosY = int32(selfBaseY)
				m.mail.SendProtoMail(member.Id(), proto, ctime)
			}
		})

		succMsg = conveneAllMsg
	} else {
		target := g.GetMember(targetId)
		if target == nil {
			errMsg = guild.ErrConveneFailTargetNotInGuild
			return
		}

		if ctime.Before(target.GetNextConveneTime()) {
			errMsg = guild.ErrConveneFailCooldown
			return
		}

		target.SetNextConveneTime(ctime.Add(m.datas.GuildGenConfig().ConveneCooldown))

		// 发召集邮件
		if data != nil {
			proto := data.NewTextMail(shared_proto.MailType_MailGuildCall)
			proto.Text = data.NewTextFields().WithName(selfName).WithPos(posString).JsonString()
			proto.GuildCallHeroId = idbytes.ToBytes(self.Id())
			proto.GuildCallPosX = int32(selfBaseX)
			proto.GuildCallPosY = int32(selfBaseY)
			m.mail.SendProtoMail(targetId, proto, ctime)
		}

		succMsg = guild.NewS2cConveneMsg(idbytes.ToBytes(targetId))
	}

	return
}

func (m *guild_func) generateRankMsgByCountry(countryId uint64, rankLimit int) (msg pbutil.Buffer) {
	m.rankModule.SubTypeRRankListFunc(shared_proto.RankType_Guild, countryId, func(rankList rankface.RRankList) {
		proto := rankList.EncodeClient(1, u64.FromInt(rankLimit))
		msg = guild.NewS2cViewDailyGuildRankMsg(proto).Static()
	})
	return
}

func (m *guild_func) addRecommendMcBuild(g *sharedguilddata.Guild, self *sharedguilddata.GuildMember, mcId uint64) (succMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
	if !self.HasPermission(func(permission *guild_data.GuildPermissionData) bool {
		return permission.RecommendMcBuild
	}) {
		errMsg = guild.ErrAddRecommendMcBuildFailDeny
		return
	}

	if g.RecommendMcBuildExist(mcId) {
		errMsg = guild.ErrAddRecommendMcBuildFailMcIsRecommended
		return
	}

	g.AddRecommendMcBuild(mcId, m.dep.Datas().McBuildMiscData().MaxRecommendMcGuildCount)

	succMsg = guild.NewS2cAddRecommendMcBuildMsg(u64.Int32Array(g.RecommendMcBuilds()))

	return
}

func (m *guild_func) addYinliangRecord(g *sharedguilddata.Guild, text string, ctime time.Time, d *guild_data.GuildLogData) {
	g.AddYinliangRecord(text, ctime)
	if d.SendChat {
		m.chat.SysChat(0, g.Id(), shared_proto.ChatType_ChatGuild, text, shared_proto.ChatMsgType_ChatMsgGuildLog, true, true, true, false)
	}
}
