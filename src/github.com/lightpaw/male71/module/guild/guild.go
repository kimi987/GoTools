package guild

import (
	"bytes"
	context2 "context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/event"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"golang.org/x/net/context"
	"math/rand"
	"strings"
	"time"
)

func NewGuildModule(dep iface.ServiceDep, datas *config.ConfigDatas, db iface.DbService, hebi iface.HebiModule, mail iface.MailModule,
	pushService iface.PushService, realmService iface.RealmService, xiongNuModule iface.XiongNuModule, chat iface.ChatService,
	xiongNuService iface.XiongNuService, rankModule iface.RankModule, tickService iface.TickerService, country iface.CountryService,
	tssClient iface.TssClient, baizhanService iface.BaiZhanService, mingcWarService iface.MingcWarService) *GuildModule {

	module := &GuildModule{}
	module.dep = dep
	module.guildService = dep.Guild()
	module.tickService = tickService
	module.realmService = realmService
	module.pushService = pushService

	module.tssClient = tssClient
	module.eventPrizeQueue = event.NewFuncQueue(1024, "GuildModule.eventPrizeQueue")
	guildService := dep.Guild()

	var templates []*guild_template
	var maxId int64
	if !guildService.Func(func(guilds sharedguilddata.Guilds) {
		var array []*sharedguilddata.Guild
		guilds.Walk(func(g *sharedguilddata.Guild) {
			array = append(array, g)
			maxId = i64.Max(maxId, g.Id())
		})

		templates = newGuildTemplateArray(datas.GetNpcGuildTemplateArray(), array)
	}) {
		logrus.Panicf("GuildModule初始化guild_template 出错")
	}

	module.guild_func = newGuildFunc(dep, datas, hebi, db, mail,
		realmService, xiongNuModule, xiongNuService, country, baizhanService, rankModule, pushService,
		mingcWarService, maxId, guildService.ClearSelfGuildMsgCache, guildService.UpdateSnapshot,
		guildService.GetSnapshot, guildService.RemoveSnapshot, chat, guildService.AddLog,
		templates)

	guildService.SetGuildRankFunc(module.guild_func.getRank, module.generateRankMsgByCountry)

	guildService.Func(func(guilds sharedguilddata.Guilds) {
		guilds.Walk(func(g *sharedguilddata.Guild) {
			module.addNotFullGuild(g)
		})
	})

	heromodule.RegisterHeroEventWithSubTypeHandler("GuildModule.HeroEventHandler", module.handleHeroEventWithSubType)

	go call.CatchLoopPanic(module.loop, "GuildModule.loop")

	return module
}

//gogen:iface
type GuildModule struct {
	*guild_func

	dep          iface.ServiceDep
	guildService iface.GuildService
	tickService  iface.TickerService
	realmService iface.RealmService
	pushService  iface.PushService
	tssClient    iface.TssClient

	eventPrizeQueue *event.FuncQueue
}

func (m *GuildModule) OnHeroOnline(hc iface.HeroController, guildId int64) {
	if guildId <= 0 {
		// 加入推荐列表
		m.guildService.AddRecommendInviteHeros(hc.Id())
		return
	}

	ctime := m.time.CurrentTime()
	listGuildEventPrizeProto := &guild.S2CListGuildEventPrizeProto{}
	var heroGuildTechnologys []*guild_data.GuildTechnologyData
	var heroGuildResetWeekly map[uint64][]int
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		// 联盟礼包
		defaultId := idbytes.ToBytes(-int64(hero.CountryId()))
		var defaultName string
		if country := m.datas.GetCountryData(hero.CountryId()); country != nil {
			defaultName = country.Name
		} else {
			defaultName = m.datas.TextHelp().MysticAlly.New().JsonString()
		}
		hero.WalkGuildEventPrize(func(p *entity.HeroGuildEventPrize) {
			if p.ExpireTime.Before(ctime) {
				hero.RemoveGuildEventPrize(p.Id)
				heroChanged = true
			} else {
				listGuildEventPrizeProto.Id = append(listGuildEventPrizeProto.Id, p.Id)
				listGuildEventPrizeProto.DataId = append(listGuildEventPrizeProto.DataId, u64.Int32(p.Data.Id))
				listGuildEventPrizeProto.ExpireTime = append(listGuildEventPrizeProto.ExpireTime, timeutil.Marshal32(p.ExpireTime))
				if p.HideGiver && p.SendHeroId != hero.Id() {
					listGuildEventPrizeProto.HeroId = append(listGuildEventPrizeProto.HeroId, defaultId)
					listGuildEventPrizeProto.HeroName = append(listGuildEventPrizeProto.HeroName, m.datas.TextHelp().MysticAlly.New().JsonString())
				} else if p.SendHeroId < 0 && p.SendHeroId == -int64(hero.CountryId()) {
					listGuildEventPrizeProto.HeroId = append(listGuildEventPrizeProto.HeroId, defaultId)
					listGuildEventPrizeProto.HeroName = append(listGuildEventPrizeProto.HeroName, defaultName)
				} else {
					listGuildEventPrizeProto.HeroId = append(listGuildEventPrizeProto.HeroId, idbytes.ToBytes(p.SendHeroId))
					if p.SendHeroId != hero.Id() {
						if snapshot := m.heroSnapshotService.Get(p.SendHeroId); snapshot != nil {
							listGuildEventPrizeProto.HeroName = append(listGuildEventPrizeProto.HeroName, snapshot.GetName())
						} else {
							listGuildEventPrizeProto.HeroName = append(listGuildEventPrizeProto.HeroName, idbytes.PlayerName(p.SendHeroId))
						}
					} else {
						listGuildEventPrizeProto.HeroName = append(listGuildEventPrizeProto.HeroName, hero.Name())
					}
				}
			}
		})

		// 联盟科技更新
		heroGuildTechnologys = hero.Domestic().GetGuildTechnology()
		// 拷贝出联盟任务领取情况
		heroGuildResetWeekly = hero.CopyGuildResetWeekly()
		return
	})

	if len(listGuildEventPrizeProto.Id) > 0 {
		hc.Send(guild.NewS2cListGuildEventPrizeProtoMsg(listGuildEventPrizeProto))
	}

	var allMemberIds []int64
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {

		var self *sharedguilddata.GuildMember
		if g != nil {
			self = g.GetMember(hc.Id())
		}

		if self == nil {
			// 玩家有联盟，但是又不在联盟里面
			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
				// 联盟不对的bug
				if hero.GuildId() != guildId {
					return
				}

				logrus.Errorf("玩家不在联盟，但是Hero中却有联盟id!%d", guildId)

				// 更新成玩家没有联盟了
				if g == nil {
					m.updateHeroGuild(hero, result, nil, nil, true)
				} else {
					m.doRemoveHeroGuild(g, hero.Id(), 0, hero, result, false)
				}

				result.Add(guild.LEAVE_GUILD_S2C)

				guildId = 0

				result.Ok()
			})
			return
		}

		// 上线联盟数据异常，同步更新一下DB hero.guild_id字段
		if guildId == 0 {
			if err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
				return m.db.UpdateHeroGuildId(ctx, hc.Id(), 0)
			}); err != nil {
				logrus.WithError(err).Errorf("保存玩家联盟ID，更新超时 heroId:%v guildId:%v", hc.Id(), 0)
			}
		}

		if g.LeaderId() == hc.Id() {
			changed := false
			if g.GetImpeachLeader() != nil {
				g.SetImpeachLeader(nil) // 取消联盟
				changed = true
				// 记录盟主提前上线导致弹劾提前结束的日志
				if d := m.datas.GuildLogHelp().TerminateImpeach; d != nil {
					if hero := m.heroSnapshotService.Get(self.Id()); hero != nil {
						proto := d.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
						proto.Text = d.Text.New().WithHeroName(hero.Name).JsonString()
						m.addGuildLog(g.Id(), proto)
					}
				}
			}
			g.UpdateLeaderOfflineTime(time.Time{}) // 更新最后在线时间
			if g.TryUpdateTarget(m.datas.GuildConfig(), ctime, 0) {
				changed = true
			}

			if changed {
				g.SetChanged()
				allMemberIds = g.AllUserMemberIds()
			}
		}

		// 发送联盟求助（这个数据变化很快，就不缓存了）
		if msg := m.getListGuildSeekHelpMsg(g); msg != nil {
			hc.Send(msg)
		}

		if msg := g.StatueCacheMsg(); msg != nil {
			hc.Send(msg)
		}

		// 联盟大宝箱
		if data := g.GetFullBigBoxData(); data != nil && g.IsFullBigBoxMember(hc.Id()) {
			hc.Send(guild.NewS2cUpdateFullBigBoxMsg(u64.Int32(data.Id), true, u64.Int32(g.GetBigBoxEnergy())))
		} else {
			hc.Send(guild.NewS2cUpdateFullBigBoxMsg(u64.Int32(g.GetBigBoxData().Id), false, u64.Int32(g.GetBigBoxEnergy())))
		}

		newTechnology := g.GetEffectTechnology()
		diff := guild_data.GetDiffTechnology(heroGuildTechnologys, newTechnology)
		if len(diff) > 0 {
			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
				hero.Domestic().SetGuildTechnology(newTechnology)
				heromodule.UpdateBuildingEffect(hero, result, m.datas, ctime, guild_data.GetTechnologyEffects(diff)...)
				result.Ok()
			})
		}

		// 联盟科技协助
		if g.GetTechUpgradeData() != nil && self.GetIsTechHelpable() {
			hc.Send(pushTechHelpableTrue)
		}

		// 联盟标识
		g.RangeMarkMsg(func(markMsg pbutil.Buffer) bool {
			hc.Send(markMsg)
			return true
		})

		// 联盟任务
		if len(heroGuildResetWeekly) > 0 {
			for id, stages := range heroGuildResetWeekly {
				maxStage := g.GetGuildTaskStageIndex(server_proto.GuildTaskType(u64.Int32(id)))
				if maxStage <= 0 { // 这种情况只可能玩家换过联盟了，否则不可能联盟未到进度而自己领过了
					continue
				}
				collectedCount := 0
				for _, stage := range stages {
					if stage <= maxStage { // 不能直接取stages的长度，因为有可能是换过联盟，领取过大于maxStage的奖励，只能一个个判定
						collectedCount++
					}
				}
				if collectedCount != maxStage { // 奖励可以跳着领的，只要小于maxStage的就直接break，就取这条任务发给前端
					hc.Send(guild.NewS2cNoticeTaskStageUpdateMsg(u64.Int32(id), int32(maxStage)))
					break
				}
			}
		}

		// 联盟昨日本国声望排名
		hc.Send(guild.NewS2cViewLastGuildRankMsg(u64.Int32(g.GetLastPrestigeRank())))

		if g.GetWorkshop() == nil && !self.GetShowWorkshopNotExist() {
			hc.Send(showWorkshopTrueMsg)
		}

		hc.Send(guild.NewS2cUpdateSelfClassLevelMsg(u64.Int32(self.ClassLevelData().Level)))
	})
	m.guildService.SelfGuildMsgCache().Clear(guildId)
	m.world.MultiSend(allMemberIds, guild.SELF_GUILD_CHANGED_S2C)
}

func (m *GuildModule) loop() {
	// 每小时重置
	hourlyTickTime := m.tickService.GetPerHourTickTime()
	// 每周重置
	weeklyTickTime := m.tickService.GetWeeklyTickTime()
	m.doResetWeekly(weeklyTickTime.GetPrevTickTime())
	// 每日重置
	dailyTickTime := m.tickService.GetDailyTickTime()
	logrus.Debug("联盟重置刷新时间", dailyTickTime.GetPrevTickTime())
	m.doResetDaily(dailyTickTime.GetPrevTickTime())

	// 定时保存
	saveTicker := time.NewTicker(5 * time.Minute)

	// ticker
	t := time.NewTicker(5 * time.Second)

	// 下一次NPC调整联盟职位的时间
	ctime := m.time.CurrentTime()
	nextNpcSetClassLevelTime := m.datas.GuildConfig().GetNextNpcSetClassLevelTime(ctime)

	minuteTickTime := ctime.Add(time.Minute)

	for {
		select {
		case <-hourlyTickTime.Tick():
			hourlyTickTime = m.tickService.GetPerHourTickTime()
			m.doResetHourly()

		case <-dailyTickTime.Tick():
			dailyTickTime = m.tickService.GetDailyTickTime()
			// 每日重置
			m.doResetDaily(dailyTickTime.GetPrevTickTime())

		case <-saveTicker.C:
			// 保存变化的帮派
			m.guildService.SaveChangedGuild()

		case <-weeklyTickTime.Tick():
			weeklyTickTime = m.tickService.GetWeeklyTickTime()
			// 每周重置
			m.doResetWeekly(weeklyTickTime.GetPrevTickTime())

		case <-t.C:

			var oldCountries []*country.CountryData
			var changeCountry []*country.CountryData
			var changeCountryHeroIds [][]int64
			m.guildService.TimeoutFunc(func(guilds sharedguilddata.Guilds) {
				// 所有定时任务
				ctime := m.time.CurrentTime()

				setNpcMemberClassLevel := false
				if nextNpcSetClassLevelTime.Before(ctime) {
					nextNpcSetClassLevelTime = nextNpcSetClassLevelTime.Add(timeutil.Day)

					// 设置联盟职位
					setNpcMemberClassLevel = true
				}

				isTickMinute := false
				if minuteTickTime.Before(ctime) {
					minuteTickTime = minuteTickTime.Add(time.Minute)
					isTickMinute = true
				}

				guilds.Walk(func(g *sharedguilddata.Guild) {

					// 踢人每日限制
					// npc帮派不允许修改入盟条件

					m.tickRemoveExpiredJoinRequest(g, ctime)
					m.tickRemoveExpiredInvateRequest(g, ctime)
					m.tickUpdateImpeachLeader(g, ctime)
					m.tickUpdateTransferLeader(g, ctime)
					m.tickUpdateUpgradeLevel(g, ctime)
					m.tickUpdateUpgradeTechnology(g, ctime)
					m.tickUpdateGuildTarget(g, ctime)

					// 更新联盟转国，
					oldCountry := g.Country()
					if m.tickUpdateGuildChangeCountry(g, ctime) {
						oldCountries = append(oldCountries, oldCountry)
						changeCountry = append(changeCountry, g.Country())
						changeCountryHeroIds = append(changeCountryHeroIds, g.AllUserMemberIds())
					}

					if isTickMinute {
						m.tickUpdateMemberLeaveTime(g, ctime)
					}

					if g.IsNpcLeader() {
						if setNpcMemberClassLevel {
							// 设置联盟职位
							m.npcSetMemberClass(g) // 每日一次
						}

						m.npcTryUpgradeLevel(g)
						m.npcTryKickMember(g, ctime)
					}

					// 更新联盟最大
					m.tickUpdateHistoryMaxPrestige(g)
				})

				// 尝试创建Npc帮派
				m.tryKeepFreeNpcGuild(guilds)
			})

			// 更新跟着联盟一起转国的英雄
			if n := imath.Min(len(changeCountry), len(changeCountryHeroIds)); n > 0 {
				for i := 0; i < n; i++ {
					oldCountry := oldCountries[i]
					newCountry := changeCountry[i]
					heroIds := changeCountryHeroIds[i]

					// 转国邮件
					var toSendMailProto *shared_proto.MailProto
					if mailData := m.datas.MailHelp().GuildChangeCountry; mailData != nil {
						toSendMailProto = mailData.NewTextMail(shared_proto.MailType_MailNormal)
						toSendMailProto.Text = mailData.NewTextFields().WithOldName(oldCountry.Name).WithNewName(newCountry.Name).JsonString()
					}

					for _, heroId := range heroIds {
						m.heroService.FuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
							heromodule.DoChangeCountry(m.dep, hero, result, newCountry.Id)
						})
						m.dep.Country().ForceOfficialDepose(oldCountry.Id, heroId)

						if toSendMailProto != nil {
							m.dep.Mail().SendProtoMail(heroId, toSendMailProto, ctime)
						}
					}
				}
			}
		}
	}

}

func (m *GuildModule) doResetWeekly(resetTime time.Time) {
	m.guildService.Func(func(guilds sharedguilddata.Guilds) {
		guilds.Walk(func(g *sharedguilddata.Guild) {
			m.resetWeekly(g, resetTime)
		})
	})
}

func (m *GuildModule) doResetHourly() {
	m.guildService.Func(func(guilds sharedguilddata.Guilds) {
		guilds.Walk(func(g *sharedguilddata.Guild) {
			g.ResetHourly()
		})
	})
}

func (m *GuildModule) doResetDaily(resetTime time.Time) {
	var gids []int64
	m.guildService.Func(func(guilds sharedguilddata.Guilds) {
		guilds.Walk(func(g *sharedguilddata.Guild) {
			m.resetDaily(g, resetTime)
		})
		gids = guilds.RefreshPrestigeRank()
	})
	for _, gid := range gids {
		var ids []int64
		var rank uint64
		m.guildService.FuncGuild(gid, func(g *sharedguilddata.Guild) {
			if g != nil {
				ids = g.AllUserMemberIds()
				rank = g.GetLastPrestigeRank()
			}
		})
		if len(ids) > 0 {
			m.world.MultiSend(ids, guild.NewS2cViewLastGuildRankMsg(u64.Int32(rank)))
		}
	}
}

func (m *GuildModule) processFuncMsg(handlerName string, serverErrorMsg pbutil.Buffer,
	hc iface.HeroController, f sharedguilddata.Funcs) {

	if !m.guildService.TimeoutFunc(f) {
		logrus.Debugf("%s，超时", handlerName)
		hc.Send(serverErrorMsg)
		return
	}
}

var listGuildEmptyMsg = guild.NewS2cListGuildMsg(nil).Static()

//gogen:iface c2s_list_guild
func (m *GuildModule) ProcessListGuild(hc iface.HeroController) {
	heroSnapshot := m.heroSnapshotService.Get(hc.Id())
	if heroSnapshot == nil {
		logrus.Debugf("没找到玩家镜像数据")
		hc.Send(guild.ERR_LIST_GUILD_FAIL_SERVER_ERROR)
		return
	}
	guildIds := m.guildService.RecommendGuildList(heroSnapshot)
	size := len(guildIds)
	if size <= 0 {
		hc.Send(listGuildEmptyMsg)
	}

	list := make([]*shared_proto.GuildSnapshotProto, 0, size)
	for _, gid := range guildIds {
		if guildSnapshot := m.guildService.GetSnapshot(gid); guildSnapshot != nil {
			list = append(list, guildSnapshot.Encode(m.dep.HeroSnapshot().GetBasicSnapshotProto))
		}
	}
	hc.Send(guild.NewS2cListGuildMsg(list))
}

var searchEmptyMsg = guild.NewS2cSearchGuildMsg(nil, nil, nil).Static()

//gogen:iface
func (m *GuildModule) ProcessSearchGuild(proto *guild.C2SSearchGuildProto, hc iface.HeroController) {
	if proto.Num < 0 {
		logrus.Debugf("搜索联盟列表，页数无效")
		hc.Send(guild.ERR_SEARCH_GUILD_FAIL_INVALID_NUM)
		return
	}

	name := strings.TrimSpace(proto.Name)
	if len := util.GetCharLen(name); len <= 0 || len > 14 {
		logrus.Debugf("搜索联盟列表，搜索名字无效")
		hc.Send(guild.ERR_SEARCH_GUILD_FAIL_INVALID_NAME)
		return
	}

	if !m.tssClient.TryCheckName("搜索联盟列表", hc, name, guild.ERR_CREATE_GUILD_FAIL_SENSITIVE_WORDS, guild.ERR_SEARCH_GUILD_FAIL_SERVER_ERROR) {
		return
	}

	page := u64.FromInt32(proto.Num)

	startIndex := page * m.datas.GuildConfig().GuildNumPerPage
	endIndex := startIndex + m.datas.GuildConfig().GuildNumPerPage

	searcherId, _ := hc.LockGetGuildId()

	// TODO 缓存
	//var guildDatas [][]byte
	list := make([]*shared_proto.GuildSnapshotProto, 0)
	yinliangMap := make(map[int64]*shared_proto.GuildYinliangSendProto, 0)
	guildIds := make([]int64, 0)
	m.processFuncMsg("搜索联盟列表",
		guild.ERR_SEARCH_GUILD_FAIL_SERVER_ERROR,
		hc, func(guilds sharedguilddata.Guilds) {
			guildArray := guilds.IdRankArray()
			endIndex = u64.Min(endIndex, uint64(len(guildArray)))

			if startIndex >= endIndex {
				// 没有数据了
				logrus.Debugf("搜索联盟列表，空列表，page: %v numPerPage: %v", proto.Num, m.datas.GuildConfig().GuildNumPerPage)
				hc.Send(searchEmptyMsg)
				return
			}

			findCount := uint64(0)
			for _, g := range guildArray {
				if strings.Contains(g.Name(), name) || strings.Contains(g.FlagName(), name) {

					if !proto.ShowSelfGuild {
						// 不显示自己
						if g.Id() == searcherId {
							continue
						}
					}

					if findCount < startIndex {
						findCount++
						continue
					}

					//proto := g.EncodeClient(false, m.dep.Mingc(), m.heroSnapshotService, m.getRank, m.xiongNuService.IsTodayStarted)
					//guildDatas = append(guildDatas, must.Marshal(proto))
					yinliangMap[g.Id()] = g.SendYinliangToMe[searcherId]
					guildIds = append(guildIds, g.Id())
					findCount++
					if findCount >= endIndex {
						break
					}
				}
			}

		})

	var yinliangList []*shared_proto.GuildYinliangSendProto

	for _, gid := range guildIds {
		g := m.guildSnapshotGetter(gid)
		if g == nil {
			continue
		}
		p := g.Encode(m.dep.HeroSnapshot().GetBasicSnapshotProto)
		if p != nil {
			list = append(list, p)
			// 赠送银两数据
			yinliangProto := yinliangMap[g.Id]
			if yinliangProto == nil {
				yinliangProto = &shared_proto.GuildYinliangSendProto{}
			}
			yinliangList = append(yinliangList, yinliangProto)
		}
	}

	hc.Send(guild.NewS2cSearchGuildMsg(nil, list, yinliangList))
}

//gogen:iface
func (m *GuildModule) ProcessCreateGuild(proto *guild.C2SCreateGuildProto, hc iface.HeroController) {
	// 先检查参数
	name := strings.TrimSpace(proto.Name)
	if c := util.GetCharLen(name); c <= 0 || c > 14 {
		logrus.Debugf("创建联盟，联盟名字长度不对")
		hc.Send(guild.ERR_CREATE_GUILD_FAIL_INVALID_NAME_LEN)
		return
	}

	// 先检查参数
	flagName := strings.TrimSpace(proto.FlagName)
	if c := util.GetCharLen(flagName); c <= 0 || c > 4 {
		logrus.Debugf("创建联盟，联盟旗号长度不对")
		hc.Send(guild.ERR_CREATE_GUILD_FAIL_INVALID_FLAG_NAME_LEN)
		return
	}

	if !m.tssClient.TryCheckName("创建联盟-名字", hc, name, guild.ERR_CREATE_GUILD_FAIL_SENSITIVE_WORDS, guild.ERR_CREATE_GUILD_FAIL_SERVER_ERROR) {
		return
	}

	if !m.tssClient.TryCheckName("创建联盟-旗号", hc, flagName, guild.ERR_CREATE_GUILD_FAIL_SENSITIVE_WORDS, guild.ERR_CREATE_GUILD_FAIL_SERVER_ERROR) {
		return
	}

	if m.datas.GuildConfig().ExistName(name) {
		logrus.Debugf("创建联盟，联盟名字跟Npc帮派名字重复")
		hc.Send(guild.ERR_CREATE_GUILD_FAIL_NAME_DUPLICATE)
		return
	}

	if m.datas.GuildConfig().ExistFlagName(flagName) {
		logrus.Debugf("创建联盟，联盟旗号跟Npc帮派旗号重复")
		hc.Send(guild.ERR_CREATE_GUILD_FAIL_FLAG_NAME_DUPLICATE)
		return
	}

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.GuildId() != 0 {
			logrus.Debugf("创建联盟，已经有联盟了")
			result.Add(guild.ERR_CREATE_GUILD_FAIL_IN_THE_GUILD)
			return
		}

		if !heromodule.HasEnoughCost(hero, m.datas.GuildConfig().CreateGuildCost) {
			logrus.Debugf("创建联盟，消耗不足")
			result.Add(guild.ERR_CREATE_GUILD_FAIL_COST_NOT_ENOUGH)
			return
		}

		result.Ok()
	}) {
		return
	}

	m.processFuncMsg("创建联盟", guild.ERR_CREATE_GUILD_FAIL_SERVER_ERROR, hc,
		func(guilds sharedguilddata.Guilds) {
			errMsg := m.createGuild(hc, name, flagName, guilds)
			if errMsg != nil {
				logrus.WithField("reason", errMsg).Debugf("创建联盟失败")
				hc.Send(errMsg.ErrMsg())
			}
		})
}

//gogen:iface c2s_self_guild
func (m *GuildModule) ProcessSelfGuild(proto *guild.C2SSelfGuildProto, hc iface.HeroController) {

	guildId := int64(proto.GuildId)
	if guildId == 0 {
		if heroGuild, ok := hc.LockGetGuildId(); !ok {
			logrus.Debugf("请求自己的联盟数据，获取联盟失败")
			hc.Send(guild.ERR_SELF_GUILD_FAIL_SERVER_ERROR)
			return
		} else {
			guildId = heroGuild
		}

		if guildId == 0 {
			logrus.Debugf("请求自己的联盟数据，你没有联盟")
			hc.Send(guild.ERR_SELF_GUILD_FAIL_NOT_IN_GUILD)
			return
		}
	}

	version, msg, err := m.guildService.SelfGuildMsgCache().GetI32VersionBuffer(guildId)
	if err != nil {
		if err == sharedguilddata.ErrNotExist {
			logrus.Errorf("请求自己的联盟数据，g == nil")
			hc.Send(guild.ERR_SELF_GUILD_FAIL_NOT_IN_GUILD)
			return
		}

		logrus.WithError(err).Errorf("请求自己的联盟数据，timeout")
		hc.Send(guild.ERR_SELF_GUILD_FAIL_SERVER_ERROR)
		return
	}

	if version != 0 && version == proto.Version {
		hc.Send(guild.SELF_GUILD_SAME_VERSION_S2C)
		return
	}

	hc.Send(msg)
}

//gogen:iface c2s_leave_guild
func (m *GuildModule) ProcessLeaveGuild(hc iface.HeroController) {

	m.processHeroInGuildMsg("退出联盟", hc,
		guild.ERR_LEAVE_GUILD_FAIL_NOT_IN_GUILD,
		guild.ERR_LEAVE_GUILD_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.leaveGuild(hc, guilds, g, self)
		})
}

//gogen:iface
func (m *GuildModule) ProcessKickOther(proto *guild.C2SKickOtherProto, hc iface.HeroController) {

	if len(proto.Id) <= 0 {
		logrus.Debugf("联盟踢人，目标id为空")
		hc.Send(guild.ERR_KICK_OTHER_FAIL_TARGET_NOT_IN_GUILD)
		return
	}

	kickTargetId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("联盟踢人，目标id解析失败")
		hc.Send(guild.ERR_KICK_OTHER_FAIL_TARGET_NOT_IN_GUILD)
		return
	}

	if hc.Id() == kickTargetId {
		logrus.Debugf("联盟踢人，要踢自己?")
		hc.Send(guild.ERR_KICK_OTHER_FAIL_DENY)
		return
	}

	if !m.processHeroInGuildMsg("联盟踢人", hc,
		guild.ERR_KICK_OTHER_FAIL_NOT_IN_GUILD,
		guild.ERR_KICK_OTHER_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.kickOther(g, self, kickTargetId)
		}) {
		// 有错误
		return
	}

}

//gogen:iface
func (m *GuildModule) ProcessUpdateText(proto *guild.C2SUpdateTextProto, hc iface.HeroController) {

	if u64.FromInt(util.GetCharLen(proto.Text)) > m.datas.GuildConfig().TextLimitChar {
		logrus.Debugf("更新联盟宣言，文本太长")
		hc.Send(guild.ERR_UPDATE_TEXT_FAIL_TEXT_TOO_LONG)
		return
	}

	if !m.tssClient.TryCheckName("更新联盟宣言", hc, proto.Text, guild.ERR_UPDATE_TEXT_FAIL_SENSITIVE_WORDS, guild.ERR_UPDATE_TEXT_FAIL_SERVER_ERROR) {
		return
	}

	m.processHeroInGuildMsg("更新联盟宣言", hc,
		guild.ERR_UPDATE_TEXT_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_TEXT_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateText(g, self, proto.Text)
		})
}

//gogen:iface
func (m *GuildModule) ProcessUpdateInternalText(proto *guild.C2SUpdateInternalTextProto, hc iface.HeroController) {

	if u64.FromInt(util.GetCharLen(proto.Text)) > m.datas.GuildConfig().InternalTextLimitChar {
		logrus.Debugf("更新联盟内部宣言，文本太长")
		hc.Send(guild.ERR_UPDATE_INTERNAL_TEXT_FAIL_TEXT_TOO_LONG)
		return
	}

	if !m.tssClient.TryCheckName("更新联盟内部宣言", hc, proto.Text, guild.ERR_UPDATE_INTERNAL_TEXT_FAIL_SENSITIVE_WORDS, guild.ERR_UPDATE_INTERNAL_TEXT_FAIL_SERVER_ERROR) {
		return
	}

	m.processHeroInGuildMsg("更新联盟内部宣言", hc,
		guild.ERR_UPDATE_INTERNAL_TEXT_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_INTERNAL_TEXT_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateInternalText(g, self, proto.Text)
		})
}

//gogen:iface
func (m *GuildModule) ProcessUpdateLabels(proto *guild.C2SUpdateGuildLabelProto, hc iface.HeroController) {

	if len(proto.Label) > int(m.datas.GuildConfig().GuildLabelLimitCount) {
		logrus.Debugf("更新联盟标签，个数超出上限")
		hc.Send(guild.ERR_UPDATE_GUILD_LABEL_FAIL_COUNT_LIMIT)
		return
	}

	for i, n := range proto.Label {
		name := strings.TrimSpace(n)
		if c := util.GetCharLen(name); c <= 0 || c > int(m.datas.GuildConfig().GuildLabelLimitChar) {
			logrus.Debugf("更新联盟标签，名字长度无效")
			hc.Send(guild.ERR_UPDATE_GUILD_LABEL_FAIL_CHAR_LIMIT)
			return
		}

		proto.Label[i] = name
	}

	if check.StringNilOrDuplicate(proto.Label) {
		logrus.Debugf("更新联盟标签，重名或者空值无效")
		hc.Send(guild.ERR_UPDATE_GUILD_LABEL_FAIL_DUPLICATE)
		return
	}

	for _, n := range proto.Label {
		if !m.tssClient.TryCheckName("更新联盟标签", hc, n, guild.ERR_UPDATE_GUILD_LABEL_FAIL_SENSITIVE_WORDS, guild.ERR_UPDATE_GUILD_LABEL_FAIL_SERVER_ERROR) {
			return
		}
	}

	m.processHeroInGuildMsg("更新联盟标签", hc,
		guild.ERR_UPDATE_GUILD_LABEL_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_GUILD_LABEL_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateLabel(g, self, proto.Label)
		})
}

//gogen:iface
func (m *GuildModule) ProcessUpdateClassNames(proto *guild.C2SUpdateClassNamesProto, hc iface.HeroController) {

	if len(proto.Name) != len(m.datas.GuildClassLevelData().Array) {
		logrus.Debugf("更新阶级名称，个数无效")
		hc.Send(guild.ERR_UPDATE_CLASS_NAMES_FAIL_INVALID_COUNT)
		return
	}

	for i, n := range proto.Name {
		name := strings.TrimSpace(n)
		if c := util.GetCharLen(name); c <= 0 || c > 10 {
			logrus.Debugf("更新阶级名称，名字长度无效")
			hc.Send(guild.ERR_UPDATE_CLASS_NAMES_FAIL_INVALID_COUNT)
		}

		proto.Name[i] = name
	}

	if check.StringNilOrDuplicate(proto.Name) {
		logrus.Debugf("更新阶级名称，重名或者空值无效")
		hc.Send(guild.ERR_UPDATE_CLASS_NAMES_FAIL_INVALID_DUPLICATE)
		return
	}

	for _, n := range proto.Name {
		if !m.tssClient.TryCheckName("更新阶级名称", hc, n, guild.ERR_UPDATE_CLASS_NAMES_FAIL_SENSITIVE_WORDS, guild.ERR_UPDATE_CLASS_NAMES_FAIL_SERVER_ERROR) {
			return
		}
	}

	m.processHeroInGuildMsg("更新阶级名称", hc,
		guild.ERR_UPDATE_CLASS_NAMES_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_CLASS_NAMES_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateClassNames(g, self, proto.Name)
		})
}

//gogen:iface
func (m *GuildModule) ProcessUpdateClassTitle(proto *guild.C2SUpdateClassTitleProto, hc iface.HeroController) {

	titleProto := &shared_proto.GuildClassTitleProto{}
	if err := titleProto.Unmarshal(proto.Proto); err != nil {
		logrus.Debugf("更新帮派职称，无效的proto")
		hc.Send(guild.ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_PROTO)
		return
	}

	// 自定义职称个数超出限制
	if uint64(len(titleProto.CustomClassTitleName)) > m.datas.GuildConfig().GuildClassTitleMaxCount {
		logrus.Debugf("更新职称，自定义职称个数超出上限")
		hc.Send(guild.ERR_UPDATE_CLASS_TITLE_FAIL_COUNT_LIMIT)
		return
	}

	if len(titleProto.SystemClassTitleId) != len(titleProto.SystemClassTitleMemberId) {
		logrus.Debugf("更新职称，系统职称id个数不等于系统职称成员id个数")
		hc.Send(guild.ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_MEMBER_ID)
		return
	}

	if check.Int32Duplicate(titleProto.SystemClassTitleId) {
		logrus.Debugf("更新职称，系统职称id重复")
		hc.Send(guild.ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_TITLE_ID)
		return
	}

	var heroIds []int64
	var titleDatas []*guild_data.GuildClassTitleData
	for i, id := range titleProto.SystemClassTitleMemberId {
		heroId, ok := idbytes.ToId(id)
		if !ok {
			logrus.Debugf("更新职称，系统成员id解析失败")
			hc.Send(guild.ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_MEMBER_ID)
			return
		}

		d := m.datas.GetGuildClassTitleData(u64.FromInt32(titleProto.SystemClassTitleId[i]))
		if d == nil {
			logrus.Debugf("更新职称，系统职称没找到，id: %v", titleProto.SystemClassTitleId[i])
			hc.Send(guild.ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_TITLE_ID)
			return
		}

		heroIds = append(heroIds, heroId)
		titleDatas = append(titleDatas, d)

	}

	for i, n := range titleProto.CustomClassTitleName {
		name := strings.TrimSpace(n)
		if c := uint64(util.GetCharLen(name)); c <= 0 || c > m.datas.GuildConfig().GuildClassTitleMaxCharCount {
			logrus.Debugf("更新职称，名字长度无效")
			hc.Send(guild.ERR_UPDATE_CLASS_TITLE_FAIL_INVALID_TITLE_ID)
			return
		}

		titleProto.CustomClassTitleName[i] = name
	}

	if check.StringNilOrDuplicate(titleProto.CustomClassTitleName) {
		logrus.Debugf("更新职称，自定义职称名字重复")
		hc.Send(guild.ERR_UPDATE_CLASS_TITLE_FAIL_NAME_EXIST)
		return
	}

	for _, n := range titleProto.CustomClassTitleName {
		if !m.tssClient.TryCheckName("更新职称", hc, n, guild.ERR_UPDATE_CLASS_TITLE_FAIL_SENSITIVE_WORDS, guild.ERR_UPDATE_CLASS_TITLE_FAIL_SERVER_ERROR) {
			return
		}
	}

	m.processHeroInGuildMsg("更新帮派职称", hc,
		guild.ERR_UPDATE_CLASS_TITLE_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_CLASS_TITLE_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateClassTitle(g, self, titleProto, heroIds, titleDatas)
		})
}

//gogen:iface
func (m *GuildModule) ProcessUpdateFlagType(proto *guild.C2SUpdateFlagTypeProto, hc iface.HeroController) {
	if proto.FlagType < 0 || proto.FlagType > 5 {
		logrus.Debugf("更新联盟旗帜，个数无效")
		hc.Send(guild.ERR_UPDATE_FLAG_TYPE_FAIL_INVALID_TYPE)
		return
	}

	flagType := u64.FromInt32(proto.FlagType)

	m.processHeroInGuildMsg("更新联盟旗帜", hc,
		guild.ERR_UPDATE_FLAG_TYPE_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_FLAG_TYPE_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateFlagType(g, self, flagType)
		})
}

//gogen:iface
func (m *GuildModule) ProcessUpdateMemberClassLevel(proto *guild.C2SUpdateMemberClassLevelProto, hc iface.HeroController) {

	if len(proto.Id) <= 0 {
		logrus.Debugf("修改成员阶级，id为空")
		hc.Send(guild.ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_NOT_IN_GUILD)
		return
	}

	targetId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("修改成员阶级，目标id解析失败")
		hc.Send(guild.ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_NOT_IN_GUILD)
		return
	}

	if hc.Id() == targetId {
		logrus.Debugf("修改成员阶级，不能改自己的阶级")
		hc.Send(guild.ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_DENY)
		return
	}

	classLevel := u64.FromInt32(proto.ClassLevel)
	newClassLevelData := m.datas.GetGuildClassLevelData(classLevel)
	if newClassLevelData == nil {
		logrus.Debugf("修改成员阶级，无效的阶级等级")
		hc.Send(guild.ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_INVALID_CLASS_LEVEL)
		return
	}

	ctime := m.time.CurrentTime()
	countryId := hc.LockHeroCountry()
	if m.dep.Country().King(countryId) == hc.Id() {
		if c := m.dep.Country().Country(countryId); c != nil {
			if ctime.Before(c.NextAppointTime(hc.Id())) {
				hc.Send(guild.ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_CHANGE_KING_IN_CD)
				return
			}
		}
	}

	m.processHeroInGuildMsg("修改成员阶级", hc,
		guild.ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_MEMBER_CLASS_LEVEL_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateMemberClassLevel(g, self, targetId, newClassLevelData)
		})
}

//gogen:iface c2s_cancel_change_leader
func (m *GuildModule) ProcessCancelChangeLeader(hc iface.HeroController) {

	m.processHeroInGuildMsg("取消变更帮主", hc,
		guild.ERR_CANCEL_CHANGE_LEADER_FAIL_NOT_IN_GUILD,
		guild.ERR_CANCEL_CHANGE_LEADER_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

			return m.cancelChangeLeader(g, self)
		})
}

//gogen:iface
func (m *GuildModule) ProcessDonation(proto *guild.C2SDonateProto, hc iface.HeroController) {

	seq := u64.FromInt32(proto.Sequence)

	m.processHeroInGuildMsg("联盟捐献", hc,
		guild.ERR_DONATE_FAIL_NOT_IN_GUILD,
		guild.ERR_DONATE_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.donation(hc, g, self, seq)
		})
}

//gogen:iface c2s_upgrade_level
func (m *GuildModule) ProcessUpgradeLevel(hc iface.HeroController) {

	m.processHeroInGuildMsg("联盟升级", hc,
		guild.ERR_UPGRADE_LEVEL_FAIL_NOT_IN_GUILD,
		guild.ERR_UPGRADE_LEVEL_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.upgradeLevel(g, self)
		})
}

//gogen:iface
func (m *GuildModule) ProcessReduceUpgradeLevelCd(proto *guild.C2SReduceUpgradeLevelCdProto, hc iface.HeroController) {

	m.processHeroInGuildMsg("联盟升级加速", hc,
		guild.ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_NOT_IN_GUILD,
		guild.ERR_REDUCE_UPGRADE_LEVEL_CD_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.reduceUpgradeLevelCd(g, self)
		})
}

//gogen:iface c2s_impeach_leader
func (m *GuildModule) ProcessImpeachLeader(hc iface.HeroController) {

	m.processHeroInGuildMsg("联盟弹劾", hc,
		guild.ERR_IMPEACH_LEADER_FAIL_NOT_IN_GUILD,
		guild.ERR_IMPEACH_LEADER_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.impeachLeader(g, self)
		})
}

//gogen:iface
func (m *GuildModule) ProcessImpeachLeaderVote(proto *guild.C2SImpeachLeaderVoteProto, hc iface.HeroController) {

	if len(proto.Target) <= 0 {
		logrus.Debugf("联盟弹劾投票，目标id为空")
		hc.Send(guild.ERR_IMPEACH_LEADER_VOTE_FAIL_INVALID_TARGET)
		return
	}

	voteTargetId, ok := idbytes.ToId(proto.Target)
	if !ok {
		logrus.Debugf("联盟弹劾投票，目标id解析失败")
		hc.Send(guild.ERR_IMPEACH_LEADER_VOTE_FAIL_INVALID_TARGET)
		return
	}

	m.processHeroInGuildMsg("联盟弹劾投票", hc,
		guild.ERR_IMPEACH_LEADER_VOTE_FAIL_NOT_IN_GUILD,
		guild.ERR_IMPEACH_LEADER_VOTE_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.impeachLeaderVote(g, self, voteTargetId)
		})
}

func (m *GuildModule) processHeroInGuildMsg(handlerName string, hc iface.HeroController, notInGuildMsg, serverErrorMsg pbutil.Buffer,
	f func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool)) (success bool) {
	guildId, ok := hc.LockGetGuildId()
	if !ok {
		logrus.Debugf("%s，获取联盟失败", handlerName)
		hc.Send(serverErrorMsg)
		return
	}

	heroId := hc.Id()

	return m.processHeroInGuildMsg0(handlerName, heroId, guildId, hc, notInGuildMsg, serverErrorMsg, f)
}

func (m *GuildModule) processHeroInGuildMsg0(handlerName string, heroId, guildId int64, sender sender.Sender, notInGuildMsg, serverErrorMsg pbutil.Buffer,
	f func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool)) (success bool) {

	if guildId == 0 {
		logrus.Debugf("%s，你没有联盟", handlerName)
		sender.Send(notInGuildMsg)
		return
	}

	if !m.guildService.TimeoutFunc(func(guilds sharedguilddata.Guilds) {

		g := guilds.Get(guildId)
		if g == nil {
			logrus.Errorf("%s，g == nil", handlerName)
			sender.Send(notInGuildMsg)
			return
		}

		self := g.GetMember(heroId)
		if self == nil {
			// TODO
			logrus.Errorf("%s，self == nil", handlerName)
			sender.Send(notInGuildMsg)
			return
		}

		successMsg, errMsg, broadcastChanged := f(guilds, g, self)
		if errMsg != nil {
			logrus.WithField("reason", errMsg).Debugf(handlerName)
			sender.Send(errMsg.ErrMsg())
			return
		}

		sender.Send(successMsg)

		// 广播给其他人
		if broadcastChanged {
			m.clearSelfGuildMsgCache(guildId)
			m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)
		}
	}) {
		logrus.Debugf("%s，超时", handlerName)
		sender.Send(serverErrorMsg)
		return
	}

	return true
}

////gogen:iface
//func (m *GuildModule) ProcessJoinGuildV01(proto *guild.C2SJoinGuildV01Proto, hc iface.HeroController) {
//
//	guildId := int64(proto.Id)
//	logrus.Debugf("加入帮派，%v", proto.Id)
//	if guildId <= 0 {
//		logrus.Debugf("加入帮派，获取联盟id失败")
//		hc.Send(guild.ERR_JOIN_GUILD_V01_FAIL_INVALID_ID)
//		return
//	}
//
//	var heroName string
//	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
//
//		if hero.GuildId() != 0 {
//			logrus.Debugf("加入帮派，已经在联盟中")
//			result.Add(guild.ERR_JOIN_GUILD_V01_FAIL_IN_THE_GUILD)
//			return
//		}
//
//		heroName = hero.Name()
//
//		result.Ok()
//	}) {
//		return
//	}
//
//	toAdd := sharedguilddata.NewMember(hc.Id(), m.config.leaderGuildClass)
//
//	var errMsg pbutil.Buffer
//	var memberIds []int64
//
//	updateGuildId := int64(0)
//	updateGuildName := ""
//	updateGuildFlagName := ""
//	m.guildService.TimeoutFunc(guildId, func(g *sharedguilddata.Guild, err error) {
//		if err != nil {
//			if err == sharedguilddata.ErrEmpty {
//				logrus.Errorf("加入帮派，g == nil")
//				errMsg = guild.ERR_JOIN_GUILD_V01_FAIL_INVALID_ID
//				return
//			}
//
//			logrus.WithError(err).Errorf("加入帮派，lock guild 报错")
//			errMsg = guild.ERR_JOIN_GUILD_V01_FAIL_SERVER_ERROR
//			return
//		}
//
//		// 看下帮派是否满员
//		if g.IsFull() {
//			logrus.Debugf("加入帮派，已经满员")
//			errMsg = guild.ERR_JOIN_GUILD_V01_FAIL_FULL
//			return
//		}
//
//		memberIds = g.AllUserMemberIds()
//		g.AddMember(toAdd)
//
//		updateGuildId = g.Id()
//		updateGuildName = g.Name()
//		updateGuildFlagName = g.FlagName()
//
//		m.clearSelfGuildMsgCache(g.Id())
//	})
//
//	if errMsg != nil {
//		hc.Send(errMsg)
//		return
//	}
//
//	if !m.casHeroGuildIdName(hc.Id(), 0, updateGuildId, updateGuildName, updateGuildFlagName, true) {
//		logrus.Debugf("加入帮派，cas帮派失败")
//		hc.Send(guild.ERR_JOIN_GUILD_V01_FAIL_IN_THE_GUILD)
//
//		if m.guildService.TimeoutFuncNotError(guildId, func(g *sharedguilddata.Guild) {
//			g.RemoveMember(hc.Id())
//		}) {
//			logrus.Errorf("加入帮派，hero cas 失败，删除帮派成员时候失败")
//		}
//		return
//	}
//
//	hc.Send(guild.NewS2cJoinGuildV01Msg(int32(guildId), updateGuildName, updateGuildFlagName))
//
//	// 广播给其他人
//	m.world.MultiSend(memberIds, guild.NewS2cAddGuildMemberMsg(hc.IdBytes(), heroName).Static())
//
//	// 更新任务进度
//	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
//		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_JOIN_GUILD)
//		result.Ok()
//	})
//
//}

// ---------------加入联盟，玩家部分

//gogen:iface
func (m *GuildModule) ProcessUpdateJoinCondition(proto *guild.C2SUpdateJoinConditionProto, hc iface.HeroController) {

	m.processHeroInGuildMsg("更新入盟条件", hc,
		guild.ERR_UPDATE_JOIN_CONDITION_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_JOIN_CONDITION_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateJoinCondition(g, self, proto.RejectAutoJoin,
				u64.FromInt32(proto.RequiredHeroLevel), u64.FromInt32(proto.RequiredJunXianLevel), u64.FromInt32(proto.RequiredTowerMaxFloor))
		})
}

//gogen:iface
func (m *GuildModule) ProcessUpdateGuildName(proto *guild.C2SUpdateGuildNameProto, hc iface.HeroController) {

	newName := strings.TrimSpace(proto.Name)
	if c := util.GetCharLen(newName); c <= 0 || c > 14 {
		logrus.Debugf("修改帮派名称，无效的联盟名字")
		hc.Send(guild.ERR_UPDATE_GUILD_NAME_FAIL_INVALID_NAME)
		return
	}

	newFlagName := strings.TrimSpace(proto.FlagName)
	if c := util.GetCharLen(newFlagName); c <= 0 || c > 4 {
		logrus.Debugf("修改帮派名称，无效的联盟旗号")
		hc.Send(guild.ERR_UPDATE_GUILD_NAME_FAIL_INVALID_FLAG_NAME)
		return
	}

	if m.datas.GuildConfig().ExistName(newName) {
		logrus.Debugf("修改帮派名称，联盟名字跟Npc帮派名字重复")
		hc.Send(guild.ERR_UPDATE_GUILD_NAME_FAIL_EXIST_NAME)
		return
	}

	if m.datas.GuildConfig().ExistFlagName(newFlagName) {
		logrus.Debugf("修改帮派名称，联盟旗号跟Npc帮派旗号重复")
		hc.Send(guild.ERR_UPDATE_GUILD_NAME_FAIL_EXIST_FLAG_NAME)
		return
	}

	if !m.tssClient.TryCheckName("修改帮派名称", hc, newName, guild.ERR_UPDATE_GUILD_NAME_FAIL_SENSITIVE_WORDS, guild.ERR_UPDATE_GUILD_NAME_FAIL_SERVER_ERROR) {
		return
	}

	if !m.tssClient.TryCheckName("修改帮派旗号", hc, newFlagName, guild.ERR_UPDATE_GUILD_NAME_FAIL_SENSITIVE_WORDS, guild.ERR_UPDATE_GUILD_NAME_FAIL_SERVER_ERROR) {
		return
	}

	var guildId int64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		guildId = hero.GuildId()
		if guildId == 0 {
			logrus.Debugf("修改帮派名称，无效的联盟旗号")
			result.Add(guild.ERR_UPDATE_GUILD_NAME_FAIL_INVALID_FLAG_NAME)
			return
		}

		// 先检查下够不够钱
		if !heromodule.HasEnoughCost(hero, m.datas.GuildConfig().ChangeGuildNameCost) {
			logrus.Debugf("修改帮派名称，消耗不足")
			result.Add(guild.ERR_UPDATE_GUILD_NAME_FAIL_COST_NOT_ENOUGH)
			return
		}

		result.Ok()
	}) {
		return
	}

	m.processHeroInGuildMsg0("修改帮派名称", hc.Id(), guildId, hc,
		guild.ERR_UPDATE_GUILD_NAME_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_GUILD_NAME_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateGuildName(hc, guilds, g, self, newName, newFlagName)
		})
}

//gogen:iface
func (m *GuildModule) ProcessListGuildByIds(proto *guild.C2SListGuildByIdsProto, hc iface.HeroController) {

	if check.Int32AnyZero(proto.Ids) {
		logrus.Debugf("ids查联盟列表，id中有0")
		hc.Send(guild.ERR_LIST_GUILD_BY_IDS_FAIL_INVALID_ID)
		return
	}

	if check.Int32Duplicate(proto.Ids) {
		logrus.Debugf("ids查联盟列表，id中有重复值")
		hc.Send(guild.ERR_LIST_GUILD_BY_IDS_FAIL_INVALID_ID)
		return
	}

	// 个数限制
	if uint64(len(proto.Ids)) > m.datas.GuildConfig().GuildNumPerPage {
		logrus.Debugf("ids查联盟列表，个数超出限制")
		hc.Send(guild.ERR_LIST_GUILD_BY_IDS_FAIL_INVALID_COUNT)
		return
	}

	m.processFuncMsg("ids查联盟列表",
		guild.ERR_LIST_GUILD_BY_IDS_FAIL_SERVER_ERROR,
		hc, func(guilds sharedguilddata.Guilds) {

			// TODO 缓存
			guildDatas := make([][]byte, 0, len(proto.Ids))
			for _, gid := range proto.Ids {
				g := guilds.Get(int64(gid))
				if g != nil {
					guildDatas = append(guildDatas, must.Marshal(g.EncodeClient(false, m.heroSnapshotService, m.getRank, m.xiongNuService.IsTodayStarted)))
				}
			}

			hc.Send(guild.NewS2cListGuildByIdsMsg(guildDatas))
		})
}

// ---- 申请加入 ----

//gogen:iface
func (m *GuildModule) ProcessUserRequestJoin(proto *guild.C2SUserRequestJoinProto, hc iface.HeroController) {

	toJoin := int64(proto.Id)
	if toJoin == 0 {
		logrus.Debugf("申请加入帮派，无效的帮派id")
		hc.Send(guild.ERR_USER_REQUEST_JOIN_FAIL_INVALID_ID)
		return
	}

	g := m.guildService.GetSnapshot(toJoin)
	if g == nil {
		logrus.Debugf("申请加入帮派，帮派不存在")
		hc.Send(guild.ERR_USER_REQUEST_JOIN_FAIL_INVALID_ID)
		return
	}

	// 先判断是否超出自己的申请上限
	var selfGuildId int64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if !g.RejectAutoJoin {
			joinGuildIds := hero.GetJoinGuildIds()
			if uint64(len(joinGuildIds)) >= m.datas.GuildConfig().GuildMaxJoinRequestCount {
				logrus.Debugf("申请加入帮派，超出申请上限")
				result.Add(guild.ERR_USER_REQUEST_JOIN_FAIL_SELF_FULL)
				return
			}

			if i64.Contains(joinGuildIds, toJoin) {
				logrus.Debugf("申请加入帮派，已经申请了这个联盟")
				result.Add(guild.ERR_USER_REQUEST_JOIN_FAIL_DUPLICATE)
				return
			}
		}

		selfGuildId = hero.GuildId()
		if selfGuildId == toJoin {
			logrus.Debugf("申请加入帮派，不能申请自己的联盟")
			result.Add(guild.ERR_USER_REQUEST_JOIN_FAIL_SELF_GUILD)
			return
		}

		result.Ok()
	}) {
		return
	}

	m.processFuncMsg("申请加入帮派",
		guild.ERR_USER_REQUEST_JOIN_FAIL_SERVER_ERROR,
		hc, func(guilds sharedguilddata.Guilds) {
			errMsg := m.userRequestJoin(hc, guilds, toJoin)
			if errMsg != nil {
				logrus.WithField("reason", errMsg).Debugf("申请加入帮派，处理消息失败")
				hc.Send(errMsg.ErrMsg())
				return
			}
		})
}

//gogen:iface
func (m *GuildModule) ProcessUserCancelJoinRequest(proto *guild.C2SUserCancelJoinRequestProto, hc iface.HeroController) {
	cancelGuildId := int64(proto.Id)

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if !i64.Contains(hero.GetJoinGuildIds(), cancelGuildId) {
			logrus.Debugf("取消申请入帮，没有申请过这个帮派")
			result.Add(guild.ERR_USER_CANCEL_JOIN_REQUEST_FAIL_INVALID_ID)
			return
		}

		result.Ok()
	}) {
		return
	}

	m.processFuncMsg("取消申请入帮",
		guild.ERR_USER_CANCEL_JOIN_REQUEST_FAIL_SERVER_ERROR,
		hc, func(guilds sharedguilddata.Guilds) {
			errMsg := m.userCancelRequestJoin(hc, guilds, cancelGuildId)
			if errMsg != nil {
				logrus.WithField("reason", errMsg).Debugf("取消申请入帮，处理消息失败")
				hc.Send(errMsg.ErrMsg())
			}
		})
}

//gogen:iface
func (m *GuildModule) ProcessGuildReplyJoinRequest(proto *guild.C2SGuildReplyJoinRequestProto, hc iface.HeroController) {

	if len(proto.Id) <= 0 {
		logrus.Debugf("审批加入联盟，目标id为空")
		hc.Send(guild.ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_INVALID_REQUEST)
		return
	}

	targetId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("审批加入联盟，目标id解析失败")
		hc.Send(guild.ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_INVALID_REQUEST)
		return
	}

	if hc.Id() == targetId {
		logrus.Debugf("审批加入联盟，目标是自己?")
		hc.Send(guild.ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_INVALID_REQUEST)
		return
	}

	var heroName string
	hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		heroName = hero.Name()
		return
	})

	gctx := &sharedguilddata.GuildContext{
		OperType:     sharedguilddata.ReplyJoinGuild,
		OperatorId:   hc.Id(),
		OperatorName: heroName,
	}

	m.processHeroInGuildMsg("审批加入联盟", hc,
		guild.ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_NO_GUILD,
		guild.ERR_GUILD_REPLY_JOIN_REQUEST_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.guildReplyJoinRequest(gctx, guilds, g, self, targetId, proto.Agree)
		})

}

// ---- End 申请加入 ----
// ---- 邀请加入 ----

//gogen:iface
func (m *GuildModule) ProcessInvateOtherRequest(proto *guild.C2SGuildInvateOtherProto, hc iface.HeroController) {

	if len(proto.Id) <= 0 {
		logrus.Debugf("邀请加入联盟，目标id为空")
		hc.Send(guild.ERR_GUILD_INVATE_OTHER_FAIL_INVALID_ID)
		return
	}

	targetId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("邀请加入联盟，目标id解析失败")
		hc.Send(guild.ERR_GUILD_INVATE_OTHER_FAIL_INVALID_ID)
		return
	}

	if hc.Id() == targetId {
		logrus.Debugf("邀请加入联盟，目标是自己?")
		hc.Send(guild.ERR_GUILD_INVATE_OTHER_FAIL_INVALID_ID)
		return
	}

	if npcid.IsNpcId(targetId) {
		logrus.Debugf("邀请加入联盟，目标是Npc?")
		hc.Send(guild.ERR_GUILD_INVATE_OTHER_FAIL_INVALID_ID)
		return
	}

	m.processHeroInGuildMsg("邀请加入联盟", hc,
		guild.ERR_GUILD_INVATE_OTHER_FAIL_NOT_IN_GUILD,
		guild.ERR_GUILD_INVATE_OTHER_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.guildInvateOther(guilds, g, self, targetId)
		})

}

//gogen:iface
func (m *GuildModule) ProcessCancelInvateOtherRequest(proto *guild.C2SGuildCancelInvateOtherProto, hc iface.HeroController) {
	if len(proto.Id) <= 0 {
		logrus.Debugf("取消邀请加入联盟，目标id为空")
		hc.Send(guild.ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_INVALID_ID)
		return
	}

	targetId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debugf("取消邀请加入联盟，目标id解析失败")
		hc.Send(guild.ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_INVALID_ID)
		return
	}

	if hc.Id() == targetId {
		logrus.Debugf("取消邀请加入联盟，目标是自己?")
		hc.Send(guild.ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_INVALID_ID)
		return
	}

	m.processHeroInGuildMsg("取消邀请加入联盟", hc,
		guild.ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_NOT_IN_GUILD,
		guild.ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.guildCancelInvateOther(g, self, targetId)
		})
}

//gogen:iface
func (m *GuildModule) ProcessUserReplyInvateRequest(proto *guild.C2SUserReplyInvateRequestProto, hc iface.HeroController) {

	targetGuildId := int64(proto.Id)

	var heroName string
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroName = hero.Name()

		if !i64.Contains(hero.GetBeenInvateGuildIds(), targetGuildId) {
			logrus.Debugf("回复帮派邀请，帮派没有邀请你")
			result.Add(guild.ERR_USER_REPLY_INVATE_REQUEST_FAIL_INVALID_ID)
			return
		}

		result.Ok()
	}) {
		logrus.Debugf("回复帮派邀请，退出")
		return
	}

	gctx := &sharedguilddata.GuildContext{
		OperType:     sharedguilddata.InviteJoinGuild,
		OperatorId:   hc.Id(),
		OperatorName: heroName,
	}

	m.processFuncMsg("取消邀请加入联盟",
		guild.ERR_GUILD_CANCEL_INVATE_OTHER_FAIL_SERVER_ERROR,
		hc, func(guilds sharedguilddata.Guilds) {

			var errMsg msg.ErrMsg
			if proto.Agree {
				errMsg = m.userAgreeInvateRequest(gctx, hc, guilds, targetGuildId)
			} else {
				errMsg = m.userRejectInvateRequest(hc, guilds, targetGuildId)
			}

			if errMsg != nil {
				logrus.WithField("reason", errMsg).Debugf("回复邀请帮派")
				hc.Send(errMsg.ErrMsg())
			}
		})
}

//gogen:iface c2s_list_invite_me_guild
func (m *GuildModule) ProcessListInviteMeGuild(hc iface.HeroController) {
	var guildIds []int64
	hc.Func(func(hero *entity.Hero, err error) (heroChanged bool) {
		guildIds = hero.GetBeenInvateGuildIds()
		return
	})

	s2cProto := &guild.S2CListInviteMeGuildProto{}
	for _, gid := range guildIds {
		g := m.guildSnapshotGetter(gid)
		if g == nil {
			logrus.Debugf("ProcessListInviteMeGuild hero.GetBeenInvateGuildIds() guild:%v 不存在snapshot。", gid)
			continue
		}
		s2cProto.GuildList = append(s2cProto.GuildList, g.Encode(m.dep.HeroSnapshot().GetBasicSnapshotProto))
	}

	hc.Send(guild.NewS2cListInviteMeGuildProtoMsg(s2cProto))
}

//gogen:iface
func (m *GuildModule) ProcessUpdateFriendGuild(proto *guild.C2SUpdateFriendGuildProto, hc iface.HeroController) {

	if u64.FromInt(util.GetCharLen(proto.Text)) > m.datas.GuildConfig().FriendGuildTextLimitChar {
		logrus.Debugf("更新友盟，文本太长")
		hc.Send(guild.ERR_UPDATE_FRIEND_GUILD_FAIL_TEXT_TOO_LONG)
		return
	}

	if !m.tssClient.TryCheckName("更新友盟", hc, proto.Text, guild.ERR_UPDATE_FRIEND_GUILD_FAIL_SENSITIVE_WORDS, guild.ERR_UPDATE_FRIEND_GUILD_FAIL_SERVER_ERROR) {
		return
	}

	m.processHeroInGuildMsg("更新友盟", hc,
		guild.ERR_UPDATE_FRIEND_GUILD_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_FRIEND_GUILD_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateFriendGuild(g, self, proto.Text)
		})
}

//gogen:iface
func (m *GuildModule) ProcessUpdateEnemyGuild(proto *guild.C2SUpdateEnemyGuildProto, hc iface.HeroController) {

	if u64.FromInt(util.GetCharLen(proto.Text)) > m.datas.GuildConfig().EnemyGuildTextLimitChar {
		logrus.Debugf("更新敌盟，文本太长")
		hc.Send(guild.ERR_UPDATE_ENEMY_GUILD_FAIL_TEXT_TOO_LONG)
		return
	}

	if !m.tssClient.TryCheckName("更新敌盟", hc, proto.Text, guild.ERR_UPDATE_ENEMY_GUILD_FAIL_SENSITIVE_WORDS, guild.ERR_UPDATE_ENEMY_GUILD_FAIL_SERVER_ERROR) {
		return
	}

	m.processHeroInGuildMsg("更新敌盟", hc,
		guild.ERR_UPDATE_ENEMY_GUILD_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_ENEMY_GUILD_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateEnemyGuild(g, self, proto.Text)
		})
}

//gogen:iface c2s_update_guild_prestige
func (m *GuildModule) ProcessUpdateGuildPrestige(proto *guild.C2SUpdateGuildPrestigeProto, hc iface.HeroController) {

	// 不允许修改联盟目标
	hc.Send(guild.ERR_UPDATE_GUILD_PRESTIGE_FAIL_TARGET_NOT_FOUND)

	//target := m.datas.GetCountryData(uint64(proto.GetTarget()))
	//if target == nil {
	//	logrus.Debugf("更新声望目标，目标没找到")
	//	hc.Send(guild.ERR_UPDATE_GUILD_PRESTIGE_FAIL_TARGET_NOT_FOUND)
	//	return
	//}
	//
	//m.processHeroInGuildMsg("更新声望目标", hc,
	//	guild.ERR_UPDATE_ENEMY_GUILD_FAIL_NOT_IN_GUILD,
	//	guild.ERR_UPDATE_ENEMY_GUILD_FAIL_SERVER_ERROR,
	//	func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
	//		return m.updateGuildPrestige(hc, g, self, target)
	//	})
}

//gogen:iface c2s_place_guild_statue
func (m *GuildModule) ProcessPlaceGuildStatue(proto *guild.C2SPlaceGuildStatueProto, hc iface.HeroController) {
	realm := m.realmService.GetBigMap()
	if realm == nil {
		logrus.WithField("realmId", proto.GetRealmId()).Debugln("放置雕像")
		hc.Send(guild.ERR_PLACE_GUILD_STATUE_FAIL_MAP_NOT_FOUND)
		return
	}

	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.WithField("guildId", guildId).Debugln("没有联盟")
		hc.Send(guild.ERR_PLACE_GUILD_STATUE_FAIL_NO_GUILD)
		return
	}

	snapshot := m.guildService.GetSnapshot(guildId)
	if snapshot == nil {
		logrus.WithField("guildId", guildId).Errorln("联盟snapshot没找到")
		hc.Send(guild.ERR_PLACE_GUILD_STATUE_FAIL_SERVER_ERROR)
		return
	}

	if snapshot.LeaderId != hc.Id() {
		logrus.WithField("leaderId", snapshot.LeaderId).WithField("id", hc.Id()).Debugln("玩家不是盟主")
		hc.Send(guild.ERR_PLACE_GUILD_STATUE_FAIL_NOT_LEADER)
		return
	}

	m.processHeroInGuildMsg("放置雕像", hc,
		guild.ERR_PLACE_GUILD_STATUE_FAIL_NO_GUILD,
		guild.ERR_PLACE_GUILD_STATUE_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.placeGuildStatue(g, self, realm)
		})
}

//gogen:iface c2s_take_back_guild_statue
func (m *GuildModule) ProcessTakeBackGuildStatue(hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.WithField("guildId", guildId).Debugln("没有联盟")
		hc.Send(guild.ERR_TAKE_BACK_GUILD_STATUE_FAIL_NO_GUILD)
		return
	}

	snapshot := m.guildService.GetSnapshot(guildId)
	if snapshot == nil {
		logrus.WithField("guildId", guildId).Errorln("联盟snapshot没找到")
		hc.Send(guild.ERR_TAKE_BACK_GUILD_STATUE_FAIL_SERVER_ERROR)
		return
	}

	if snapshot.LeaderId != hc.Id() {
		logrus.WithField("leaderId", snapshot.LeaderId).WithField("id", hc.Id()).Debugln("玩家不是盟主")
		hc.Send(guild.ERR_TAKE_BACK_GUILD_STATUE_FAIL_NOT_LEADER)
		return
	}

	m.processHeroInGuildMsg("放置雕像", hc,
		guild.ERR_TAKE_BACK_GUILD_STATUE_FAIL_NO_GUILD,
		guild.ERR_TAKE_BACK_GUILD_STATUE_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.takeBackGuildStatue(g, self)
		})
}

//gogen:iface c2s_collect_first_join_guild_prize
func (m *GuildModule) ProcessCollectFirstJoinGuildPrize(hc iface.HeroController) {
	hc.Send(guild.ERR_COLLECT_FIRST_JOIN_GUILD_PRIZE_FAIL_COLLECTED)
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	if hero.GuildId() == 0 {
	//		logrus.Debugln("领取首次加入联盟奖励，没有联盟")
	//		result.Add(guild.ERR_COLLECT_FIRST_JOIN_GUILD_PRIZE_FAIL_NO_GUILD)
	//		return
	//	}
	//
	//	if hero.IsFirstJoinGuildPrizeCollected() {
	//		logrus.Debugln("领取首次加入联盟奖励，首次加入联盟奖励已经领取了")
	//		result.Add(guild.ERR_COLLECT_FIRST_JOIN_GUILD_PRIZE_FAIL_COLLECTED)
	//		return
	//	}
	//
	//	hero.CollectFirstJoinGuildPrize()
	//
	//	result.Add(guild.COLLECT_FIRST_JOIN_GUILD_PRIZE_S2C)
	//
	//	heromodule.AddPrize(hero, result, m.datas.GuildConfig().FirstJoinGuildPrize, m.time.CurrentTime())
	//
	//	result.Changed()
	//	result.Ok()
	//})
}

//gogen:iface
func (m *GuildModule) ProcessSeekHelp(proto *guild.C2SSeekHelpProto, hc iface.HeroController) {

	var guildId int64
	var seekHelpProto *shared_proto.GuildSeekHelpProto
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.GuildId() == 0 {
			logrus.Debugf("请求联盟帮助，自己没有联盟")
			result.Add(guild.ERR_SEEK_HELP_FAIL_NOT_IN_GUILD)
			return
		}

		building := hero.Domestic().GetBuilding(shared_proto.BuildingType_WAI_SHI_YUAN)
		if building == nil {
			logrus.Debugf("请求联盟帮助，自己还没有外使院建筑")
			result.Add(guild.ERR_SEEK_HELP_FAIL_WAI_SHI_YUAN)
			return
		}

		if proto.HelpType == constants.SeekTypeWorker {
			// 设置请求状态为false，返回成功消息
			if !hero.Domestic().UnsetWorkerSeekHelpIfTrue(int(proto.WorkerPos)) {
				logrus.WithField("workerPos", proto.WorkerPos).Debugf("请求联盟帮助，建筑队当前不能请求帮助")
				result.Add(guild.ERR_SEEK_HELP_FAIL_DISABLE)
				return
			}

		} else {
			// 设置请求状态为false，返回成功消息
			if !hero.Domestic().UnsetTechSeekHelpIfTrue(int(proto.WorkerPos)) {
				logrus.WithField("workerPos", proto.WorkerPos).Debugf("请求联盟帮助，科研队当前不能请求帮助")
				result.Add(guild.ERR_SEEK_HELP_FAIL_DISABLE)
				return
			}
		}

		// 返回成功消息
		result.Add(guild.NewS2cSeekHelpMsg(proto.HelpType, proto.WorkerPos))

		guildId = hero.GuildId()
		seekHelpProto = &shared_proto.GuildSeekHelpProto{
			HeroId:                hero.IdBytes(),
			HeroName:              hero.Name(),
			HeroHead:              hero.Head(),
			HelpType:              proto.HelpType,
			WorkerPos:             proto.WorkerPos,
			ReduceSecondsPerCount: int32(building.Effect.SeekHelpCdr / time.Second),
			HelpMaxHeroCount:      u64.Int32(building.Effect.SeekHelpMaxTimes),
		}
	})

	if guildId != 0 && seekHelpProto != nil {
		var memberIds []int64
		m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
			if g == nil {
				logrus.Debugf("请求联盟帮助，联盟不存在（之前检查有的...）")
				return
			}

			m := g.GetMember(hc.Id())
			if m == nil {
				logrus.Debugf("请求联盟帮助，自己不是这个联盟的人")
				return
			}

			// 新增一条请求，然后广播给联盟里的人
			g.AddSeekHelp(seekHelpProto)
			memberIds = g.AllUserMemberIds()
		})

		if len(memberIds) > 0 {
			m.world.MultiSend(memberIds, guild.NewS2cAddGuildSeekHelpMarshalMsg(seekHelpProto))
		}
	}

}

//gogen:iface
func (m *GuildModule) ProcessHelpGuildMember(proto *guild.C2SHelpGuildMemberProto, hc iface.HeroController) {

	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.Debugf("帮助盟友求助，自己不在联盟中")
		hc.Send(guild.ERR_HELP_GUILD_MEMBER_FAIL_NOT_IN_GUILD)
		return
	}

	var allMemberIds []int64
	var broadcastMsgs []pbutil.Buffer
	var helpProto *shared_proto.GuildSeekHelpProto
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Error("帮助盟友求助，g == nil")
			hc.Send(guild.ERR_HELP_GUILD_MEMBER_FAIL_NOT_IN_GUILD)
			return
		}

		self := g.GetMember(hc.Id())
		if self == nil {
			logrus.Error("帮助盟友求助，member == nil")
			hc.Send(guild.ERR_HELP_GUILD_MEMBER_FAIL_NOT_IN_GUILD)
			return
		}

		helpProto = g.GetSeekHelp(proto.Id)
		if helpProto == nil {
			logrus.Debugf("帮助盟友求助，helpProto == nil")
			hc.Send(guild.ERR_HELP_GUILD_MEMBER_FAIL_ID_NOT_FOUND)
			return
		}

		if len(helpProto.HelpHeroIds) >= int(helpProto.HelpMaxHeroCount) {
			// 清理掉这条求助
			logrus.Errorf("帮助盟友求助，求助次数已满，为什么没有移除掉")
			hc.Send(guild.ERR_HELP_GUILD_MEMBER_FAIL_ID_NOT_FOUND)

			// 清理掉这条数据
			g.RemoveSeekHelp(proto.Id)
			allMemberIds = g.AllUserMemberIds()
			broadcastMsgs = append(broadcastMsgs, guild.NewS2cRemoveGuildSeekHelpMsg(proto.Id))
			return
		}

		// 自己已经帮助过这条求助了
		for _, v := range helpProto.HelpHeroIds {
			if bytes.Equal(hc.IdBytes(), v) {
				logrus.Debugf("帮助盟友求助，这个求助已经帮助过了")
				hc.Send(guild.ERR_HELP_GUILD_MEMBER_FAIL_HELPED)
				return
			}
		}

		helpProto.HelpHeroIds = append(helpProto.HelpHeroIds, hc.IdBytes())

		allMemberIds = g.AllUserMemberIds()

		if len(helpProto.HelpHeroIds) >= int(helpProto.HelpMaxHeroCount) {
			g.RemoveSeekHelp(proto.Id)
			broadcastMsgs = append(broadcastMsgs, guild.NewS2cRemoveGuildSeekHelpMsg(proto.Id))
		} else {
			broadcastMsgs = append(broadcastMsgs, guild.NewS2cAddGuildSeekHelpHeroIdsMsg(proto.Id, hc.IdBytes()))
		}

		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

			// 帮助成功
			result.Add(guild.NewS2cHelpGuildMemberMsg(proto.Id))

			// 加帮助盟友贡献
			if hero.Domestic().GetDailyHelpMemberTimes() < m.datas.GuildConfig().ContributionMaxCountPerDay {
				times := hero.Domestic().IncDailyHelpMemberTimes()
				result.Add(guild.NewS2cUpdateHelpMemberTimesMsg(u64.Int32(times)))

				hctx := heromodule.NewContext(m.dep, operate_type.GuildHelpMember)
				heromodule.AddGuildContributionCoin(hctx, hero, result, m.datas.GuildConfig().ContributionPerHelp)
			}

			heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HELP_GUILD_MEMBER)
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_HelpGuildMember)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_HELP_GUILD_MEMBER)

			result.Changed()
			result.Ok()
		})
	})

	for _, msg := range broadcastMsgs {
		m.world.MultiSend(allMemberIds, msg)
	}

	if helpProto != nil {

		helpMemberId, _ := idbytes.ToId(helpProto.HeroId)

		// 被帮助飘字
		m.world.SendFunc(helpMemberId, func() pbutil.Buffer {

			var memberName string
			hs := m.heroSnapshotService.Get(hc.Id())
			if hs != nil {
				memberName = hs.Name
			} else {
				memberName = idbytes.PlayerName(hc.Id())
			}

			return misc.NewS2cScreenShowWordsMsg(m.datas.TextHelp().MiscHelpGuildMember.New().
				WithMemberName(memberName).
				WithReduceMinute(i32.Max(helpProto.ReduceSecondsPerCount/60, 1)).
				JsonString())
		})

		// 被帮助的人，减少CD
		m.heroService.FuncWithSend(helpMemberId, func(hero *entity.Hero, result herolock.LockResult) {

			workerPos := int(helpProto.WorkerPos)
			ctime := m.time.CurrentTime()
			duration := -time.Duration(helpProto.ReduceSecondsPerCount) * time.Second

			if helpProto.HelpType == constants.SeekTypeWorker {
				// 如果有CD可以减少，则减CD
				ok, t := hero.Domestic().GetWorkerRestEndTime(workerPos)
				if ok && ctime.Before(t) {
					workerRestEndTime, _ := hero.Domestic().AddWorkerRestEndTime(workerPos, ctime, duration)
					result.AddFunc(func() pbutil.Buffer {
						return domestic.NewS2cBuildingWorkerTimeChangedMsg(int32(workerPos), timeutil.Marshal32(workerRestEndTime))
					})
				}
			} else {
				// 如果有CD可以减少，则减CD
				ok, t := hero.Domestic().GetTechWorkerRestEndTime(workerPos)
				if ok && ctime.Before(t) {
					workerRestEndTime, _ := hero.Domestic().AddTechWorkerRestEndTime(workerPos, ctime, duration)
					result.AddFunc(func() pbutil.Buffer {
						return domestic.NewS2cTechWorkerTimeChangedMsg(int32(workerPos), timeutil.Marshal32(workerRestEndTime))
					})
				}
			}
		})

	}

}

//gogen:iface c2s_help_all_guild_member
func (m *GuildModule) ProcessHelpAllGuildMember(hc iface.HeroController) {

	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.Debugf("一键帮助盟友求助，自己不在联盟中")
		hc.Send(guild.ERR_HELP_ALL_GUILD_MEMBER_FAIL_NOT_IN_GUILD)
		return
	}

	var allMemberIds []int64
	var broadcastMsgs []pbutil.Buffer
	var helpProtos []*shared_proto.GuildSeekHelpProto
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Error("一键帮助盟友求助，g == nil")
			hc.Send(guild.ERR_HELP_ALL_GUILD_MEMBER_FAIL_NOT_IN_GUILD)
			return
		}

		self := g.GetMember(hc.Id())
		if self == nil {
			logrus.Error("一键帮助盟友求助，member == nil")
			hc.Send(guild.ERR_HELP_ALL_GUILD_MEMBER_FAIL_NOT_IN_GUILD)
			return
		}

		allMemberIds = g.AllUserMemberIds()

		helpTimes := uint64(0)
		g.RangeSeekHelp(func(helpProto *shared_proto.GuildSeekHelpProto) (isContinue bool) {
			isContinue = true

			if len(helpProto.HelpHeroIds) >= int(helpProto.HelpMaxHeroCount) {
				// 清理掉这条求助
				logrus.Errorf("一键帮助盟友求助，求助次数已满，为什么没有移除掉")

				// 清理掉这条数据
				g.RemoveSeekHelp(helpProto.Id)
				broadcastMsgs = append(broadcastMsgs, guild.NewS2cRemoveGuildSeekHelpMsg(helpProto.Id))
				return
			}

			// 自己的求助不帮助
			if bytes.Equal(hc.IdBytes(), helpProto.HeroId) {
				return
			}

			// 自己已经帮助过这条求助了
			for _, v := range helpProto.HelpHeroIds {
				if bytes.Equal(hc.IdBytes(), v) {
					return
				}
			}

			helpProtos = append(helpProtos, helpProto)

			helpProto.HelpHeroIds = append(helpProto.HelpHeroIds, hc.IdBytes())

			if len(helpProto.HelpHeroIds) >= int(helpProto.HelpMaxHeroCount) {
				g.RemoveSeekHelp(helpProto.Id)
				broadcastMsgs = append(broadcastMsgs, guild.NewS2cRemoveGuildSeekHelpMsg(helpProto.Id))
			} else {
				broadcastMsgs = append(broadcastMsgs, guild.NewS2cAddGuildSeekHelpHeroIdsMsg(helpProto.Id, hc.IdBytes()))
			}

			helpTimes++

			return
		})

		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

			// 帮助成功
			var helpIds []string
			for _, v := range helpProtos {
				helpIds = append(helpIds, v.Id)
			}
			result.Add(guild.NewS2cHelpAllGuildMemberMsg(helpIds))

			// 加帮助盟友贡献
			toAddTimes := u64.Min(helpTimes,
				u64.Sub(m.datas.GuildConfig().ContributionMaxCountPerDay, hero.Domestic().GetDailyHelpMemberTimes()))
			if toAddTimes > 0 {
				times := hero.Domestic().AddDailyHelpMemberTimes(toAddTimes)
				result.Add(guild.NewS2cUpdateHelpMemberTimesMsg(u64.Int32(times)))

				hctx := heromodule.NewContext(m.dep, operate_type.GuildHelpMember)
				heromodule.AddGuildContributionCoin(hctx, hero, result, m.datas.GuildConfig().ContributionPerHelp*toAddTimes)
			}

			heromodule.IncreTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HELP_GUILD_MEMBER, u64.FromInt(len(helpIds)))
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_HelpGuildMember, u64.FromInt(len(helpIds)))
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_HELP_GUILD_MEMBER)

			result.Changed()
			result.Ok()
		})
	})

	for _, msg := range broadcastMsgs {
		m.world.MultiSend(allMemberIds, msg)
	}

	if len(helpProtos) > 0 {
		var memberName string
		if hs := m.heroSnapshotService.Get(hc.Id()); hs != nil {
			memberName = hs.Name
		} else {
			memberName = idbytes.PlayerName(hc.Id())
		}

		// 相同英雄，只调用一次
		heroIdProtoMap := make(map[int64][]*shared_proto.GuildSeekHelpProto)
		for _, helpProto := range helpProtos {
			heroId, _ := idbytes.ToId(helpProto.HeroId)
			heroIdProtoMap[heroId] = append(heroIdProtoMap[heroId], helpProto)
		}

		for helpMemberId, helpProtos := range heroIdProtoMap {

			// 被帮助飘字
			for _, helpProto := range helpProtos {
				m.world.SendFunc(helpMemberId, func() pbutil.Buffer {

					return misc.NewS2cScreenShowWordsMsg(
						m.datas.TextHelp().MiscHelpGuildMember.New().
							WithMemberName(memberName).
							WithReduceMinute(i32.Max(helpProto.ReduceSecondsPerCount/60, 1)).
							JsonString())
				})
			}

			// 被帮助的人，减少CD
			m.heroService.FuncWithSend(helpMemberId, func(hero *entity.Hero, result herolock.LockResult) {

				for _, helpProto := range helpProtos {
					workerPos := int(helpProto.WorkerPos)
					ctime := m.time.CurrentTime()
					duration := -time.Duration(helpProto.ReduceSecondsPerCount) * time.Second

					if helpProto.HelpType == constants.SeekTypeWorker {
						// 如果有CD可以减少，则减CD
						ok, t := hero.Domestic().GetWorkerRestEndTime(workerPos)
						if ok && ctime.Before(t) {
							workerRestEndTime, _ := hero.Domestic().AddWorkerRestEndTime(workerPos, ctime, duration)
							result.AddFunc(func() pbutil.Buffer {
								return domestic.NewS2cBuildingWorkerTimeChangedMsg(int32(workerPos), timeutil.Marshal32(workerRestEndTime))
							})
						}
					} else {
						// 如果有CD可以减少，则减CD
						ok, t := hero.Domestic().GetTechWorkerRestEndTime(workerPos)
						if ok && ctime.Before(t) {
							workerRestEndTime, _ := hero.Domestic().AddTechWorkerRestEndTime(workerPos, ctime, duration)
							result.AddFunc(func() pbutil.Buffer {
								return domestic.NewS2cTechWorkerTimeChangedMsg(int32(workerPos), timeutil.Marshal32(workerRestEndTime))
							})
						}
					}
				}
			})
		}
	}
}

func (m *GuildModule) handleAddPrestigeEvent(heroId, guildId int64, datas []*guild_data.GuildPrestigeEventData, subType uint64) {
	if len(datas) <= 0 {
		return
	}

	m.eventPrizeQueue.MustFunc(func() {

		ctime := m.time.CurrentTime()

		var toAddHufu uint64
		var triggerDatas []*guild_data.GuildPrestigeEventData
		m.guildService.TimeoutFunc(func(guilds sharedguilddata.Guilds) {
			g := guilds.Get(guildId)
			if g == nil {
				return
			}

			var toAddPrestige uint64
			for _, data := range datas {
				if !data.IgnoreMemberLimit && heroId != 0 {

					if !g.IsPrestigeHero(heroId) {
						if g.GetPrestigeHeroCount() < g.MemberCount() {
							g.PutPrestigeHero(heroId)
						} else {
							// 加声望，有人数限制
							continue
						}
					}
				}
				toAddPrestige += data.Prestige
				toAddHufu += data.Hufu
				triggerDatas = append(triggerDatas, data)
			}

			// 加联盟声望
			if member := g.GetMember(heroId); member != nil && member.ClassLevelData().CorePrestige {
				g.AddPrestigeCore(toAddPrestige)
			}
			g.AddPrestige(toAddPrestige)
			// 加国家声望
			m.dep.Country().AddPrestige(g.Country().Id, toAddPrestige)

			// 尝试更新历史最大声望
			m.tryUpdateGuildHistoryMaxPrestige(g)

			g.TryUpdateTarget(m.datas.GuildConfig(), ctime, shared_proto.GuildTargetType_PrestigeUp)

			m.updateGuildRankObj(g)
		})

		hero := m.dep.HeroSnapshot().Get(heroId)
		if hero == nil {
			logrus.Errorf("handleAddPrestigeEvent，找不到 heroSnapshot。heroId:%v", heroId)
			return
		}

		// 加虎符
		m.guildService.AddHufu(toAddHufu, heroId, guildId, hero.Name, hero.Head)

		// 加联盟声望事件日志
		for _, data := range triggerDatas {
			switch data.TriggerEvent {
			case shared_proto.HeroEvent_HERO_EVENT_KILL_MONSTER:
				if d := m.datas.GuildLogHelp().InvaseMonster; data != nil {
					monsterLevel := subType
					proto := d.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
					proto.Text = d.Text.New().WithHeroName(hero.Name).WithLevel(monsterLevel).WithPrestige(data.Prestige).JsonString()
					m.addGuildLog(guildId, proto)
				}
			case shared_proto.HeroEvent_HERO_EVENT_COLLECT_BAI_ZHAN_SALARY:
				if dd := m.datas.GetJunXianLevelData(subType); dd != nil {
					if d := m.datas.GuildLogHelp().CollectSalary; d != nil {
						proto := d.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
						proto.Text = d.Text.New().WithHeroName(hero.Name).WithPrestige(data.Prestige).WithJunXian(dd.Name).WithHufu(data.Hufu).JsonString()
						m.addGuildLog(guildId, proto)
					}
				}
			case shared_proto.HeroEvent_HERO_EVENT_ASSEMBLY_KILL_MONSTER:
				if dd := m.datas.GetNpcBaseData(subType); dd != nil {
					if d := m.datas.GuildLogHelp().AssemblyWin; data != nil {
						proto := d.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
						proto.Text = d.Text.New().WithHeroName(hero.Name).WithText(dd.Name).WithAmount(data.Prestige).JsonString()
						m.addGuildLog(guildId, proto)
					}
				}
			}
		}

		m.clearSelfGuildMsgCache(guildId)

	})
}

// 帮派科技最大声望更新
func (m *GuildModule) tickUpdateHistoryMaxPrestige(g *sharedguilddata.Guild) {

	if m.tryUpdateGuildHistoryMaxPrestige(g) {
		// 升级成功

		m.clearSelfGuildMsgCache(g.Id())
		m.world.MultiSend(g.AllUserMemberIds(), guild.SELF_GUILD_CHANGED_S2C)

		m.updateSnapshot(g)
		g.SetChanged()
	}

}

func (m *GuildModule) tryUpdateGuildHistoryMaxPrestige(g *sharedguilddata.Guild) bool {

	historyMaxPretige := g.GetHistoryMaxPrestige()
	if historyMaxPretige >= g.GetPrestige() {
		return false
	}

	newPrestige := g.RefreshHistoryMaxPrestige(g.GetPrestige())
	if historyMaxPretige >= m.datas.GuildPrestigePrizeData().MaxKeyData.Prestige {
		return false
	}

	ctime := m.time.CurrentTime()

	var prizeDatas []*guild_data.GuildEventPrizeData
	oldHufu := g.GetHufu()
	for _, d := range m.datas.GuildPrestigePrizeData().Array {
		if historyMaxPretige >= d.Prestige {
			continue
		}

		if newPrestige < d.Prestige {
			break
		}

		// 加建设值
		g.AddBuildingAmount(d.BuildingAmount, ctime)

		// 加虎符值
		g.AddHufu(d.Hufu)

		// 发国家礼包
		prizeDatas = append(prizeDatas, d.EventPrize)

		if data := m.datas.GuildLogHelp().PrestigePrize; data != nil {
			proto := data.NewLogProto(ctime)
			proto.Text = data.Text.New().WithPrestige(d.Prestige).WithBuildingAmount(d.BuildingAmount).WithAmount(d.Hufu).JsonString()

			m.addGuildLog(g.Id(), proto)
		}
	}

	memberIds := g.AllUserMemberIds()
	if oldHufu != g.GetHufu() {
		m.world.MultiSend(memberIds, guild.NewS2cUpdateHufuMsg(int32(g.GetHufu())))
	}

	if len(prizeDatas) > 0 {
		country := g.Country()
		santaId := -int64(country.Id)
		santaName := country.Name
		addPrizeFunc := m.newAddHeroGuildEventPrizeFunc(santaId, idbytes.ToBytes(santaId), santaName, g.Id(), prizeDatas, false)
		for _, memberId := range memberIds {
			m.heroService.FuncWithSend(memberId, addPrizeFunc)
		}
	}

	return true
}

func (m *GuildModule) handleHeroEventWithSubType(hero *entity.Hero, result herolock.LockResult, event shared_proto.HeroEvent, subType uint64) {

	guildId := hero.GuildId()
	if guildId == 0 {
		return
	}

	triggerData := hero.GetDailyGuildPrestigeEventTriggerTimesData()

	data := m.datas.GetGuildPrestigeEventData(guild_data.GetPrestigeEventId(event, subType))
	var triggerDatas []*guild_data.GuildPrestigeEventData
	if data != nil && data.TriggerEventCondition.Equal {
		// 给联盟加声望(条件相等判断时候直接使用)
		if data.TriggerEventTimes > 0 {
			if triggerData.Get(data.Id) < data.TriggerEventTimes {
				triggerData.Increse(data.Id)
				triggerDatas = append(triggerDatas, data)
			}
		} else {
			triggerDatas = append(triggerDatas, data)
		}
	}

	// 条件不相等的事件触发
	if datas := m.datas.GuildConfig().GetPrestigeEvent(event); len(datas) > 0 {
		for _, data := range datas {
			if data.TriggerEventCondition.Compare(subType) {
				// 给联盟加声望(条件相等判断时候直接使用)
				if data.TriggerEventTimes > 0 {
					if triggerData.Get(data.Id) >= data.TriggerEventTimes {
						continue
					}
					triggerData.Increse(data.Id)
				}

				triggerDatas = append(triggerDatas, data)
			}
		}
	}

	m.handleAddPrestigeEvent(hero.Id(), guildId, triggerDatas, subType)

	datas := m.datas.GuildConfig().GetEventPrizes(event)
	if len(datas) > 0 {
		m.handleGiveGuildEventPrize(hero, result, guildId, datas, subType)
	}
}

func (m *GuildModule) GmGiveGuildEventPrize(hero *entity.Hero, result herolock.LockResult, prizeDatas []*guild_data.GuildEventPrizeData) {
	g := m.guildService.GetSnapshot(hero.GuildId())
	if g == nil {
		return
	}

	memberIds := g.UserMemberIds
	if !i64.Contains(memberIds, hero.Id()) {
		memberIds = append(memberIds, hero.Id())
	}

	m.asyncGiveGuildEventPrize(memberIds, hero, result, prizeDatas)
}

func (m *GuildModule) HandleGiveGuildEventPrize(hero *entity.Hero, result herolock.LockResult, guildId int64, datas []*guild_data.GuildEventPrizeData, subType uint64) {
	m.handleGiveGuildEventPrize(hero, result, guildId, datas, subType)
}

func (m *GuildModule) handleGiveGuildEventPrize(hero *entity.Hero, result herolock.LockResult, guildId int64, datas []*guild_data.GuildEventPrizeData, subType uint64) {
	g := m.guildService.GetSnapshot(guildId)
	if g == nil {
		logrus.Error("触发英雄事件，但是找不到联盟快照，guildId", guildId)
		return
	}

	memberIds := g.UserMemberIds
	if len(memberIds) <= 0 || !i64.Contains(memberIds, hero.Id()) {
		logrus.Error("触发英雄事件，联盟快照中的成员列表不包含玩家自己，guildId", guildId)

		// 简单处理的话，只给玩家一个人发
		memberIds = []int64{hero.Id()}
	}

	// 根据HeroEvent查找联盟事件礼包，有的话，异步调用联盟线程，进行发礼包操作

	var triggerDatas []*guild_data.GuildEventPrizeData
	for _, data := range datas {

		if !data.TriggerEventCondition.Compare(subType) {
			continue
		}

		if data.TriggerEventTimes > 0 {
			triggerData := hero.GetGuildEventPrizeTriggerTimesData(data.TriggerEventDailyReset)

			if triggerData.Get(data.Id) >= data.TriggerEventTimes {
				continue
			}

			triggerData.Increse(data.Id)
		}

		// 触发成功
		triggerDatas = append(triggerDatas, data)

	}

	m.asyncGiveGuildEventPrize(memberIds, hero, result, triggerDatas)
}

func (m *GuildModule) asyncGiveGuildEventPrize(memberIds []int64, santa *entity.Hero, santaResult herolock.LockResult, prizeDatas []*guild_data.GuildEventPrizeData) {

	// 给联盟中的所有人发礼包
	if len(prizeDatas) <= 0 {
		return
	}

	// 异步调用
	heroId := santa.Id()
	heroIdBytes := santa.IdBytes()
	heroName := santa.Name()
	guildId := santa.GuildId()

	// 处理函数
	hideSanta := !santa.Settings().IsPrivacySettingOpen(shared_proto.PrivacySettingType_PST_SHARE_GUILD_GIFT_GIVER)
	addPrizeFunc := m.newAddHeroGuildEventPrizeFunc(heroId, heroIdBytes, heroName, guildId, prizeDatas, hideSanta)

	// 给玩家自己加
	addPrizeFunc(santa, santaResult)

	if len(memberIds) > 1 {
		// 给联盟成员加
		m.eventPrizeQueue.MustFunc(func() {
			for _, memberId := range memberIds {
				if memberId != heroId {
					// 给其他盟友加
					m.heroService.FuncWithSend(memberId, addPrizeFunc)
				}
			}
		})
	}

}

func (m *GuildModule) newAddHeroGuildEventPrizeFunc(santaId int64, santaIdBytes []byte, santaName string,
	guildId int64, prizeDatas []*guild_data.GuildEventPrizeData, hideSanta bool) herolock.SendFunc {

	ctime := m.time.CurrentTime()

	prizeCount := len(prizeDatas)
	expireTime := make([]time.Time, prizeCount)
	intExpireTime := make([]int32, prizeCount)
	prizeDataId := make([]int32, prizeCount)
	for i, prizeData := range prizeDatas {
		t := ctime.Add(prizeData.ExipreDuration)
		expireTime[i] = t
		intExpireTime[i] = timeutil.Marshal32(t)
		prizeDataId[i] = u64.Int32(prizeData.Id)
	}

	return func(hero *entity.Hero, result herolock.LockResult) {
		if hero.GuildId() != guildId {
			// 只给盟友发
			return
		}

		for i, prizeData := range prizeDatas {
			if prizeData.DailyLimit > 0 {
				collectData := hero.GetGuildEventPrizeCollectTimesData()
				if collectData.Get(prizeData.Id) >= prizeData.DailyLimit {
					// 已经到达次数上限，不给
					continue
				}

				// 次数加1
				collectData.Increse(prizeData.Id)
			}

			// 加礼包
			newId := hero.AddGuildEventPrize(prizeData, santaId, expireTime[i], hideSanta)

			result.AddFunc(func() pbutil.Buffer {
				if hideSanta && hero.Id() != santaId {
					santaName = m.datas.TextHelp().MysticAlly.New().JsonString()
					santaIdBytes = idbytes.ToBytes(-int64(hero.CountryId()))
				}
				return guild.NewS2cAddGuildEventPrizeMsg(newId, prizeDataId[i], intExpireTime[i], santaIdBytes, santaName)
			})
		}

		// 超出最大上限，替换老的礼包，遍历找到最小id的N个礼包，移除掉
		count := uint64(hero.GetGuildEventPrizeCount())
		if count > m.datas.GuildConfig().EventPrizeMaxCount {
			prizes := hero.SortGuildEventPrizeExpireSlice()

			removeCount := count - m.datas.GuildConfig().EventPrizeMaxCount
			for i := uint64(0); i < removeCount; i++ {
				toRemove := prizes[i]
				hero.RemoveGuildEventPrize(toRemove.Id)

				if ctime.Before(toRemove.ExpireTime) {
					result.AddFunc(func() pbutil.Buffer {
						return guild.NewS2cRemoveGuildEventPrizeMsg(toRemove.Id)
					})
				}
			}
		}

		result.Changed()
		result.Ok()
	}
}

//gogen:iface
func (m *GuildModule) ProcessCollectGuildEventPrize(proto *guild.C2SCollectGuildEventPrizeProto, hc iface.HeroController) {

	var guildId int64
	var toAddEnergy uint64
	var toAddPrizeBytes []byte
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		guildId = hero.GuildId()
		if guildId == 0 {
			logrus.Debug("领取联盟礼包，自己没有联盟")
			result.Add(guild.ERR_COLLECT_GUILD_EVENT_PRIZE_FAIL_NOT_IN_GUILD)
			return
		}

		ctime := m.time.CurrentTime()
		hctx := heromodule.NewContext(m.dep, operate_type.GuildCollectGuildEventPrize)

		if proto.Id == 0 {
			// 一键领取
			builder := resdata.NewPrizeBuilder()
			hero.WalkGuildEventPrize(func(p *entity.HeroGuildEventPrize) {
				if vipData := m.dep.Datas().GetVipLevelData(hero.VipLevel()); vipData == nil || !vipData.GuildPrizeOneKeyCollect {
					result.Add(guild.ERR_COLLECT_GUILD_EVENT_PRIZE_FAIL_VIP_LIMIT)
					return
				}

				hero.RemoveGuildEventPrize(p.Id)
				if p.ExpireTime.Before(ctime) {
					return
				}

				builder.Add(p.Data.Prize.GetPrize())
				toAddEnergy += p.Data.Energy
			})

			toAddPrize := builder.Build()
			heromodule.AddPrize(hctx, hero, result, toAddPrize, ctime)

			toAddPrizeBytes = must.Marshal(toAddPrize.Encode())
		} else {
			eventPrize := hero.GetGuildEventPrize(proto.Id)
			if eventPrize == nil {
				logrus.Debug("领取联盟礼包，礼包不存在")
				result.Add(guild.ERR_COLLECT_GUILD_EVENT_PRIZE_FAIL_ID_NOT_FOUND)
				return
			}

			hero.RemoveGuildEventPrize(eventPrize.Id)
			if eventPrize.ExpireTime.Before(ctime) {
				logrus.Debug("领取联盟礼包，礼包已过期")
				result.Add(guild.ERR_COLLECT_GUILD_EVENT_PRIZE_FAIL_ID_NOT_FOUND)
				return
			}

			toAdd := eventPrize.Data.Prize.GetPrize()
			heromodule.AddPrize(hctx, hero, result, toAdd, ctime)

			toAddEnergy = eventPrize.Data.Energy
			toAddPrizeBytes = must.Marshal(toAdd.Encode())
		}

		result.Ok()
	})

	if toAddEnergy > 0 {
		// 加宝箱经验
		newEnergy := m.addBigBoxEnergy(guildId, toAddEnergy)

		hc.Send(guild.NewS2cCollectGuildEventPrizeMsg(proto.Id, u64.Int32(newEnergy), u64.Int32(toAddEnergy), toAddPrizeBytes))
	}
}

func (m *GuildModule) GmAddBigBoxEnergy(guildId int64, toAddEnergy uint64) {

	if toAddEnergy <= 0 {
		m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
			if g == nil {
				return
			}

			toAddEnergy = u64.Sub(g.GetBigBoxData().UnlockEnergy, g.GetBigBoxEnergy())
		})

		if toAddEnergy <= 0 {
			return
		}
	}

	m.addBigBoxEnergy(guildId, toAddEnergy)

}

func (m *GuildModule) addBigBoxEnergy(guildId int64, toAddEnergy uint64) uint64 {
	// 加宝箱经验
	var newEnergy uint64
	var sendFullBigBoxMailFunc func()
	var broadcastFullBigBoxMsgFunc func()
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.WithField("guild_id", guildId).Error("领取联盟礼包，加宝箱经验，g == nil")
			return
		}

		var full bool
		newEnergy, full = g.AddBigBoxEnergy(toAddEnergy)
		if full {
			// 加完之后，满了

			// 上一个宝箱，还未领取的玩家，给玩家发邮件
			fullBigBoxData, fullBigBoxMemberIds := g.ClearFullBigBox()
			if fullBigBoxData != nil && len(fullBigBoxMemberIds) > 0 {
				sendFullBigBoxMailFunc = func() {
					// 原来联盟的逼格宝箱没有领取，发邮件处理
					for _, memberId := range fullBigBoxMemberIds {
						proto := m.datas.MailHelp().GuildBigBox.NewTextMail(shared_proto.MailType_MailNormal)
						proto.Prize = fullBigBoxData.PlunderPrize.GetPrize().PrizeProto()

						ctime := m.time.CurrentTime()
						m.mail.SendProtoMail(memberId, proto, ctime)
					}
				}
			}

			memberIds := g.AllUserMemberIds()

			// 更新下一个逼格宝箱
			originBigBoxData := g.GetBigBoxData()
			nextBigBoxData := originBigBoxData // 下一个宝箱还是同一个
			g.SetNextBigBox(nextBigBoxData, memberIds)

			// 发送给联盟内所有的人
			newEnergy = g.GetBigBoxEnergy()

			if len(memberIds) > 0 {
				broadcastFullBigBoxMsgFunc = func() {
					m.world.MultiSend(memberIds, guild.NewS2cUpdateFullBigBoxMsg(u64.Int32(originBigBoxData.Id), true, u64.Int32(newEnergy)))

					m.pushService.MultiPush(shared_proto.SettingType_ST_GUILD_BIG_BOX, memberIds, 0)
				}
			}
		}
	})

	// 清掉缓存
	m.guildService.ClearSelfGuildMsgCache(guildId)

	if broadcastFullBigBoxMsgFunc != nil {
		broadcastFullBigBoxMsgFunc()
	}

	if sendFullBigBoxMailFunc != nil {
		sendFullBigBoxMailFunc()
	}

	return newEnergy
}

//gogen:iface c2s_collect_full_big_box
func (m *GuildModule) ProcessCollectFullBigBox(hc iface.HeroController) {

	m.processHeroInGuildMsg("领取联盟大宝箱", hc,
		guild.ERR_COLLECT_FULL_BIG_BOX_FAIL_NOT_IN_GUILD,
		guild.ERR_COLLECT_FULL_BIG_BOX_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {

			fullBigBoxData := g.GetFullBigBoxData()
			if fullBigBoxData == nil {
				logrus.Debug("领取联盟大宝箱，没有可以领取的宝箱")
				errMsg = guild.ErrCollectFullBigBoxFailLocked
				return
			}

			ctime := m.time.CurrentTime()

			// 加入联盟时间不足，不能领取宝箱
			if ctime.Before(self.GetCreateTime().Add(m.datas.GuildConfig().BigBoxCollectableDuration)) {
				logrus.Debug("领取联盟大宝箱，加入联盟时间不足，不能领取")
				errMsg = guild.ErrCollectFullBigBoxFailTimeNotEnough
				return
			}

			if !g.RemoveFullBigBoxMemberId(self.Id()) {
				logrus.Debug("领取联盟大宝箱，自己不在领取列表中")
				errMsg = guild.ErrCollectFullBigBoxFailLocked
				return
			}

			hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
				hctx := heromodule.NewContext(m.dep, operate_type.GuildCollectFullBigBox)
				heromodule.AddPrize(hctx, hero, result, fullBigBoxData.PlunderPrize.GetPrize(), ctime)
				result.Add(guild.NewS2cUpdateFullBigBoxMsg(u64.Int32(g.GetBigBoxData().Id), false, u64.Int32(g.GetBigBoxEnergy())))

				m.dep.Tlog().TlogGuildFlow(hero, operate_type.GuildCollectFullBigBox.Id(), u64.FromInt64(g.Id()), g.LevelData().Level, u64.FromInt(g.MemberCount()))
				result.Ok()
			})

			successMsg = guild.COLLECT_FULL_BIG_BOX_S2C

			return
		})
}

//gogen:iface
func (m *GuildModule) ProcessUpgradeTechnology(proto *guild.C2SUpgradeTechnologyProto, hc iface.HeroController) {
	firstLevel := m.datas.GetGuildTechnologyData(guild_data.GetTechnologyDataId(u64.FromInt32(proto.Group), 1))
	if firstLevel == nil {
		logrus.Debug("升级联盟科技，Group不存在")
		hc.Send(guild.ERR_UPGRADE_TECHNOLOGY_FAIL_INVALID_GROUP)
		return
	}

	m.processHeroInGuildMsg("科技升级", hc,
		guild.ERR_UPGRADE_TECHNOLOGY_FAIL_NOT_IN_GUILD,
		guild.ERR_UPGRADE_TECHNOLOGY_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.upgradeTechnology(hc, g, self, firstLevel)
		})
}

//gogen:iface
func (m *GuildModule) ProcessReduceTechnologyCd(proto *guild.C2SReduceTechnologyCdProto, hc iface.HeroController) {
	m.processHeroInGuildMsg("科技升级加速", hc,
		guild.ERR_REDUCE_TECHNOLOGY_CD_FAIL_NOT_IN_GUILD,
		guild.ERR_REDUCE_TECHNOLOGY_CD_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.reduceUpgradeTechnologyCd(g, self)
		})
}

//gogen:iface c2s_help_tech
func (m *GuildModule) ProcessHelpTech(hc iface.HeroController) {

	m.processHeroInGuildMsg("联盟科技协助", hc,
		guild.ERR_HELP_TECH_FAIL_NOT_IN_GUILD,
		guild.ERR_HELP_TECH_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.helpTech(g, self)
		})
}

var emptyGuildLogsMsg = guild.NewS2cListGuildLogsMsg(nil).Static()

//gogen:iface
func (m *GuildModule) ProcessListGuildLogs(proto *guild.C2SListGuildLogsProto, hc iface.HeroController) {

	if proto.LogType == 0 {
		logrus.Debug("请求联盟日志，无效的日志类型 0")
		hc.Send(emptyGuildLogsMsg)
		return
	}

	if _, exist := shared_proto.GuildLogType_name[proto.LogType]; !exist {
		logrus.Debug("请求联盟日志，无效的日志类型", proto.LogType)
		hc.Send(emptyGuildLogsMsg)
		return
	}

	guildId, _ := hc.LockGetGuildId()
	if guildId == 0 {
		logrus.Debug("请求联盟日志，玩家没有联盟")
		hc.Send(emptyGuildLogsMsg)
		return
	}

	var logs []*shared_proto.GuildLogProto
	err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		logs, err = m.db.LoadGuildLogs(ctx, guildId, shared_proto.GuildLogType(proto.LogType), int64(proto.MinId), u64.FromInt32(proto.Count))
		return
	})
	if err != nil {
		logrus.WithError(err).Error("请求联盟日志，报错")
		hc.Send(emptyGuildLogsMsg)
		return
	}

	datas := make([][]byte, 0, len(logs))
	for _, d := range logs {
		datas = append(datas, must.Marshal(d))
	}
	hc.Send(guild.NewS2cListGuildLogsMsg(datas))
}

var noRecommendGuildMsg = guild.NewS2cRequestRecommendGuildMsg(false, 0, 0, "", "", 0).Static()

// 请求推荐联盟
//gogen:iface c2s_request_recommend_guild
func (m *GuildModule) ProcessRequestRecommendGuild(hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	if guildId != 0 {
		logrus.Debugf("客户端加入了联盟了，就不要来请求推荐联盟了")
		hc.Send(noRecommendGuildMsg)
		return
	}

	heroSnapshot := m.heroSnapshotService.Get(hc.Id())
	if heroSnapshot == nil {
		logrus.Debugf("请求推荐联盟，取不到玩家镜像数据")
		hc.Send(noRecommendGuildMsg)
		return
	}

	var nextNotifyGuildTime time.Time
	if hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		nextNotifyGuildTime = hero.NextNotifyGuildTime()
		return
	}) {
		logrus.Debugf("请求推荐联盟，服务器出错")
		hc.Send(noRecommendGuildMsg)
		return
	}

	if nextNotifyGuildTime.After(m.time.CurrentTime()) {
		logrus.Debugf("请求推荐联盟，cd没到呢")
		hc.Send(noRecommendGuildMsg)
		return
	}

	snapshots := make([]*guildsnapshotdata.GuildSnapshot, 0, m.datas.GuildConfig().NotifyJoinGuildMaxPrestigeRank)

	var minPrestige uint64

	// 找推荐联盟
	//m.guildService.GetSnapshot()
	m.walkNotFullGuild(func(id int64) (endWalk bool) {
		snapshot := m.guildService.GetSnapshot(id)
		if snapshot == nil {
			return
		}

		if snapshot.MemberCount >= snapshot.GuildLevel.MemberCount {
			// 满了
			return
		}

		if snapshot.GetCountryId() != heroSnapshot.CountryId {
			// 国家不同
			return
		}

		// 入盟条件
		if snapshot.RejectAutoJoin {
			return
		}

		if snapshot.RequiredHeroLevel > heroSnapshot.Level {
			// 君主等级不够
			return
		}

		if snapshot.RequiredJunXianLevel > heroSnapshot.BaiZhanJunXianLevel {
			// 百战军衔不够
			return
		}

		if snapshot.RequiredTowerMaxFloor > heroSnapshot.TowerMaxFloor {
			// 最大千重楼层数不够
			return
		}

		if uint64(len(snapshots)) >= m.datas.GuildConfig().NotifyJoinGuildMaxPrestigeRank {
			// 满了
			if minPrestige >= snapshot.TotalPrestigeDaily {
				// 最小的都比你大了
				return
			}
		}

		if minPrestige == 0 {
			minPrestige = snapshot.TotalPrestigeDaily
		} else {
			minPrestige = u64.Min(minPrestige, snapshot.TotalPrestigeDaily)
		}

		snapshots = append(snapshots, snapshot)

		return
	})

	// 找到xx-xx名
	if len(snapshots) <= 0 {
		// 没得可以推荐的
		hc.Send(noRecommendGuildMsg)
		return
	}

	// 随机一下
	snapshot := snapshots[rand.Intn(len(snapshots))]
	// 发送推荐这个
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		result.Changed()
		result.Ok()

		nextNotifyGuildTime := m.time.CurrentTime().Add(m.datas.GuildConfig().NotifyJoinGuildDuration)
		hero.SetNextNotifyGuildTime(nextNotifyGuildTime)
		result.Add(guild.NewS2cRequestRecommendGuildMsg(true, timeutil.Marshal32(nextNotifyGuildTime), i64.Int32(snapshot.Id), snapshot.Name, snapshot.FlagName, u64.Int32(snapshot.GetCountryId())))
	})
}

//gogen:iface c2s_recommend_invite_heros
func (m *GuildModule) ProcessRecommendInviteHeros(hc iface.HeroController) {
	guildId, ok := hc.LockGetGuildId()
	var invateHeroIds []int64
	var countryData *country.CountryData
	if ok {
		m.guildService.Func(func(guilds sharedguilddata.Guilds) {
			if g := guilds.Get(guildId); g != nil {
				invateHeroIds = i64.Copy(g.GetInvateHeroIds())
				countryData = g.Country()
			}
		})
	}

	if countryData == nil {
		logrus.Error("获取联盟邀请玩家推荐列表，countryData == nil")
		hc.Send(guild.NewS2cRecommendInviteHerosMsg(nil))
		return
	}

	heros := m.guildService.RecommendInviteHeroList(countryData.Id,
		m.datas.GuildConfig().RecommendInviteHeroCount, invateHeroIds)

	proto := &shared_proto.GuildRecommendInviteHeros{}
	for _, hero := range heros {
		proto.Hero = append(proto.Hero, hero.EncodeClient())
	}
	protoBytes, err := proto.Marshal()
	if err != nil {
		hc.Send(guild.ERR_RECOMMEND_INVITE_HEROS_FAIL_SERVER_ERROR)
		return
	}

	hc.Send(guild.NewS2cRecommendInviteHerosMsg(protoBytes))
}

//gogen:iface
func (m *GuildModule) ProcessSearchNoGuildHeros(proto *guild.C2SSearchNoGuildHerosProto, hc iface.HeroController) {
	name := strings.TrimSpace(proto.Name)
	if len(name) <= 0 {
		logrus.Debugf("联盟模糊搜索玩家，关键字错误 name: %v", name)
		hc.Send(guild.ERR_SEARCH_NO_GUILD_HEROS_FAIL_INVALID_ARG)
		return
	}
	if proto.Page < 0 {
		logrus.Debugf("联盟模糊搜索玩家，页数错误 page: %v", proto.Page)
		hc.Send(guild.ERR_SEARCH_NO_GUILD_HEROS_FAIL_INVALID_ARG)
		return
	}

	ctime := m.time.CurrentTime()
	if ctime.Before(hc.NextSearchNoGuildHeros()) {
		logrus.Debugf("联盟模糊搜索玩家，请求太频繁")
		hc.Send(guild.ERR_SEARCH_NO_GUILD_HEROS_FAIL_TOO_FAST)
		return
	}
	hc.SetNextSearchNoGuildHeros(ctime.Add(m.datas.GuildConfig().SearchNoGuildHerosDuration))

	size := m.datas.GuildConfig().SearchNoGuildHerosPerPageSize
	if size <= 0 {
		logrus.Debugf("联盟模糊搜索玩家，size <= 0")
		hc.Send(guild.ERR_SEARCH_NO_GUILD_HEROS_FAIL_INVALID_ARG)
		return
	}

	if !m.tssClient.TryCheckName("联盟模糊搜索玩家", hc, name, guild.ERR_CREATE_GUILD_FAIL_SENSITIVE_WORDS, guild.ERR_SEARCH_NO_GUILD_HEROS_FAIL_SERVER_ERROR) {
		return
	}

	page := u64.FromInt32(proto.Page)

	startIndex := page * size

	var heros []*entity.Hero
	err := ctxfunc.Timeout3s(func(ctx context2.Context) (err error) {
		heros, err = m.db.LoadNoGuildHeroListByName(ctx, name, startIndex, size)
		return
	})
	if err != nil {
		logrus.WithError(err).Debugf("联盟模糊搜索玩家，页数错误 page: %v", proto.Page)
		hc.Send(guild.ERR_SEARCH_NO_GUILD_HEROS_FAIL_SERVER_ERROR)
		return
	}

	msgProto := &shared_proto.GuildRecommendInviteHeros{}
	for _, hero := range heros {
		if hero == nil {
			logrus.Debugf("联盟模糊搜索玩家，hero == nil")
			continue
		}

		heroSnapshot := m.heroSnapshotService.GetFromCache(hero.Id())

		if heroSnapshot == nil {
			if hero.GuildId() != 0 {
				continue
			}
			heroSnapshot = m.heroSnapshotService.NewSnapshot(hero)
		}

		if heroSnapshot.GuildId != 0 {
			continue
		}

		msgProto.Hero = append(msgProto.Hero, heroSnapshot.EncodeClient())
	}

	protoBytes, err := msgProto.Marshal()
	if err != nil {
		logrus.WithError(err).Debugf("联盟模糊搜索玩家，proto.Marshal() err")
		hc.Send(guild.ERR_SEARCH_NO_GUILD_HEROS_FAIL_SERVER_ERROR)
		return
	}

	hc.Send(guild.NewS2cSearchNoGuildHerosMsg(protoBytes))
}

//gogen:iface
func (m *GuildModule) ProcessUpdateGuildMark(proto *guild.C2SUpdateGuildMarkProto, hc iface.HeroController) {

	//if proto == nil {
	//	logrus.Debug("更新联盟标记，proto == nil")
	//	hc.Send(guild.ERR_UPDATE_GUILD_MARK_FAIL_INVALID_INDEX)
	//	return
	//}

	idx := u64.FromInt32(proto.Index)
	if idx <= 0 || idx > m.datas.GuildGenConfig().GuildMarkCount {
		logrus.Debug("更新联盟标记，无效的序号")
		hc.Send(guild.ERR_UPDATE_GUILD_MARK_FAIL_INVALID_INDEX)
		return
	}

	if proto.PosX < 0 || proto.PosY < 0 {
		logrus.Debug("更新联盟标记，无效的坐标")
		hc.Send(guild.ERR_UPDATE_GUILD_MARK_FAIL_INVALID_POS)
		return
	}

	proto.Msg = strings.TrimSpace(proto.Msg)
	if len(proto.Msg) > 0 && !m.tssClient.TryCheckName("更新联盟标记", hc, proto.Msg, guild.ERR_UPDATE_GUILD_MARK_FAIL_SENSITIVE_WORDS, guild.ERR_UPDATE_GUILD_MARK_FAIL_SERVER_ERROR) {
		return
	}

	mark := &shared_proto.GuildMarkProto{
		Index: proto.Index,
		PosX:  proto.PosX,
		PosY:  proto.PosY,
		Msg:   proto.Msg,
	}

	m.processHeroInGuildMsg("更新联盟标记", hc,
		guild.ERR_UPDATE_GUILD_MARK_FAIL_NOT_IN_GUILD,
		guild.ERR_UPDATE_GUILD_MARK_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.updateMark(g, self, mark)
		})
}

// gm模块

func (m *GuildModule) GmUpgradeGuildLevel(guildId int64) {

	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		ctime := m.time.CurrentTime()
		g.SetUpgradeEndTime(ctime.Add(-1))
		m.tryUpgradeGuildLevel(g, ctime)

		m.clearSelfGuildMsgCache(g.Id())

		m.updateSnapshot(g)
	})
}

func (m *GuildModule) GmAddGuildYinliang(amount int64, guildId int64) {

	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		if amount > 0 {
			g.AddYinliang(u64.FromInt64(amount))
		} else {
			g.ReduceYinliang(u64.FromInt64(amount))
		}

		m.clearSelfGuildMsgCache(g.Id())

		m.updateSnapshot(g)
	})
}

func (m *GuildModule) GmAddGuildBuildAmount(guildId int64, toAdd uint64) {

	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		ctime := m.time.CurrentTime()
		g.AddBuildingAmount(toAdd, ctime)

		m.clearSelfGuildMsgCache(g.Id())
	})
}

func (m *GuildModule) GmMiaoGuildTechCd(guildId int64) {

	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		ctime := m.time.CurrentTime()
		if timeutil.IsZero(g.GetTechUpgradeEndTime()) {
			return
		}

		g.SetTechUpgradeEndTime(ctime.Add(-1))
		m.tryUpgradeTechnology(g, ctime)

		m.clearSelfGuildMsgCache(g.Id())
	})
}

func (m *GuildModule) GmOpenImpeachLeader(guildId int64) {
	ctime := m.time.CurrentTime()
	var memIds []int64
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g.IsNpcLeader() {
			endTime := m.datas.GuildConfig().GetNextNpcSetClassLevelTime(ctime)
			g.StartImpeachLeader(g.LeaderId(), ctime, endTime, m.datas.GuildConfig().ImpeachExtraCandidateCount)
		}
		memIds = g.AllUserMemberIds()
	})
	m.guildService.SelfGuildMsgCache().Clear(guildId)
	m.world.MultiSend(memIds, guild.SELF_GUILD_CHANGED_S2C)
}

func (m *GuildModule) GmRemoveNpcGuild() {

	m.guildService.Func(func(guilds sharedguilddata.Guilds) {
		guilds.Walk(func(g *sharedguilddata.Guild) {
			if !g.IsNpcLeader() {
				return
			}

			g.WalkMember(func(member *sharedguilddata.GuildMember) {
				if member.IsNpc() {
					g.RemoveMember(member.Id(), m.time.CurrentTime().Unix())
					return
				}

				if !m.guild_func.GmTryKickMember(g, member.Id()) {
					logrus.Debugf("GM 删除 NPC 联盟 id:%v name:%v 失败，成员id:%v 没有被删除", g.Id(), g.Name(), member.Id())
				} else {
					m.world.Send(member.Id(), guild.LEAVE_GUILD_S2C)
					m.dep.Guild().AddRecommendInviteHeros(member.Id())
				}
			})

			if g.MemberCount() <= 0 {
				guilds.Remove(g.Id())
				m.removeNotFullGuild(g.Id())
				m.removeSnapshot(g.Id())

				if err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
					return m.db.DeleteGuild(ctx, g.Id())
				}); err != nil {
					// 出错了，只能打个日志出来，到下一次开服的时候，清掉所有的没有成员的帮派
					logrus.WithError(err).Errorf("退出联盟，联盟没人了，删除出错")
				}

				// 调用删除联盟
				m.rankModule.RemoveRankObj(shared_proto.RankType_Guild, g.Id())
			}

			m.clearSelfGuildMsgCache(g.Id())

			if guilds.Get(g.Id()) == nil {
				logrus.Debugf("GM 删除 NPC 联盟 id:%v name:%v 成功", g.Id(), g.Name())
			} else {
				logrus.Debugf("GM 删除 NPC 联盟 id:%v name:%v 失败", g.Id(), g.Name())
			}
		})
	})
}

//gogen:iface c2s_view_mc_war_record
func (m *GuildModule) ProcessViewMcWarRecord(hc iface.HeroController) {
	if gid, ok := hc.LockGetGuildId(); ok {
		var p *shared_proto.McWarAllRecordProto
		m.guildService.FuncGuild(gid, func(g *sharedguilddata.Guild) {
			if g == nil {
				hc.Send(guild.ERR_VIEW_MC_WAR_RECORD_FAIL_NO_GUILD)
				return
			}
			p = g.McWarRecord()
		})
		hc.Send(m.buildMcWarRecordMsg(hc, p))
	} else {
		hc.Send(guild.ERR_VIEW_MC_WAR_RECORD_FAIL_NO_GUILD)
	}
}

func (m *GuildModule) buildMcWarRecordMsg(hc iface.HeroController, allRecords *shared_proto.McWarAllRecordProto) pbutil.Buffer {

	var warIdObj *entity.JoinedMcWarIds
	if err := ctxfunc.Timeout3s(func(ctx context2.Context) (err error) {
		warIdObj, err = m.dep.Db().LoadJoinedMcWarId(ctx, hc.Id())
		return
	}); err != nil {
		warIdObj = entity.NewJoinedMcWarIds()
	}

	proto := &shared_proto.McWarAllRecordWithJoinedProto{}
	for _, p := range allRecords.Record {
		mcIds := warIdObj.WarMcIds[p.WarId]
		var isJoined bool
		for _, mcId := range mcIds {
			if mcId == p.McId {
				isJoined = true
				break
			}
		}
		proto.Record = append(proto.Record, p)
		proto.IsJoined = append(proto.IsJoined, isJoined)
	}

	return guild.NewS2cViewMcWarRecordMsg(allRecords, proto)
}

//gogen:iface
func (m *GuildModule) ProcessSendYinliangToOtherGuild(proto *guild.C2SSendYinliangToOtherGuildProto, hc iface.HeroController) {
	amount := u64.FromInt32(proto.Amount)
	if amount <= 0 {
		hc.Send(guild.ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_NOT_ENOUGH)
		return
	}

	receiverId := int64(proto.Gid)
	m.processHeroInGuildMsg("SendYinliangToOtherGuild", hc,
		guild.ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_DENY,
		guild.ERR_SEND_YINLIANG_TO_OTHER_GUILD_FAIL_SERVER_ERR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.sendYinliangToGuild(guilds, g, self, receiverId, amount)
		})
}

//gogen:iface
func (m *GuildModule) ProcessSendYinliangToMember(proto *guild.C2SSendYinliangToMemberProto, hc iface.HeroController) {
	amount := u64.FromInt32(proto.Amount)
	if amount <= 0 {
		hc.Send(guild.ERR_SEND_YINLIANG_TO_MEMBER_FAIL_INVALID_AMOUNT)
		return
	}

	receiverId, ok := idbytes.ToId(proto.MemId)
	if !ok {
		hc.Send(guild.ERR_SEND_YINLIANG_TO_MEMBER_FAIL_NO_MEMBER)
		return
	}
	m.processHeroInGuildMsg("SendYinliangToMember", hc,
		guild.ERR_SEND_YINLIANG_TO_MEMBER_FAIL_DENY,
		guild.ERR_SEND_YINLIANG_TO_MEMBER_FAIL_SERVER_ERR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.sendYinliangToMem(g, self, receiverId, amount)
		})
}

//gogen:iface c2s_pay_salary
func (m *GuildModule) ProcessPaySalary(hc iface.HeroController) {
	m.processHeroInGuildMsg("PaySalary", hc,
		guild.ERR_PAY_SALARY_FAIL_DENY,
		guild.ERR_PAY_SALARY_FAIL_SERVER_ERR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.paySalary(g, self)
		})
}

//gogen:iface
func (m *GuildModule) ProcessSetSalary(proto *guild.C2SSetSalaryProto, hc iface.HeroController) {
	amount := u64.FromInt32(proto.Salary)
	if amount < 0 {
		hc.Send(guild.ERR_SET_SALARY_FAIL_INVALID_SALARY)
		return
	}

	targetId, ok := idbytes.ToId(proto.MemId)
	if !ok {
		hc.Send(guild.ERR_SET_SALARY_FAIL_NO_MEMBER)
		return
	}

	m.processHeroInGuildMsg("SetSalary", hc,
		guild.ERR_SET_SALARY_FAIL_DENY,
		guild.ERR_SET_SALARY_FAIL_SERVER_ERR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.SetSalary(g, self, targetId, amount)
		})
}

//gogen:iface c2s_view_yinliang_record
func (m *GuildModule) ProcessViewYinliangRecord(hc iface.HeroController) {
	if gid, ok := hc.LockGetGuildId(); ok {
		m.guildService.FuncGuild(gid, func(g *sharedguilddata.Guild) {
			if g == nil {
				hc.Send(guild.ERR_VIEW_YINLIANG_RECORD_FAIL_NO_GUILD)
				return
			}

			msg := g.YinliangRecordMsg()
			if msg == nil {
				msg = g.BuildYinliangRecordMsg()
			}
			hc.Send(msg)
		})
	} else {
		hc.Send(guild.ERR_VIEW_YINLIANG_RECORD_FAIL_NO_GUILD)
	}
}

//gogen:iface c2s_view_send_yinliang_to_guild
func (m *GuildModule) ProcessViewSendYinliangToGuild(hc iface.HeroController) {
	if gid, ok := hc.LockGetGuildId(); ok {
		m.guildService.FuncGuild(gid, func(g *sharedguilddata.Guild) {
			if g == nil {
				hc.Send(guild.ERR_VIEW_SEND_YINLIANG_TO_GUILD_FAIL_SERVER_ERR)
				return
			}

			toSend := &shared_proto.GuildAllYinliangSendToGuildProto{}
			g.RangeYinliangSendToGuildRecord(func(r *shared_proto.GuildYinliangSendToGuildProto) (toContinue bool) {
				toSend.Guilds = append(toSend.Guilds, r)
				return true
			})

			hc.Send(guild.NewS2cViewSendYinliangToGuildMsg(toSend))
		})
	} else {
		hc.Send(guild.ERR_VIEW_SEND_YINLIANG_TO_GUILD_FAIL_SERVER_ERR)
	}
}

////gogen:iface c2s_ask_for_help
//func (m *GuildModule) ProcessAskForHelp(hc iface.HeroController)  {
//
//	ctime := m.time.CurrentTime();
//	if ctime.Sub(hc.LastClickTime()) < 5 * time.Second {
//		hc.Send(guild.ERR_ASK_FOR_HELP_FAIL_IN_CD)
//		return
//	}
//	hc.SetLastClickTime(ctime)
//
//	if gid, ok := hc.LockGetGuildId(); ok {
//		m.guildService.FuncGuild(gid, func(g *sharedguilddata.Guild) {
//			if g == nil {
//				hc.Send(guild.ERR_ASK_FOR_HELP_FAIL_NO_GUILD)
//				return
//			}
//		})
//	} else {
//		hc.Send(guild.ERR_ASK_FOR_HELP_FAIL_IN_CD)
//	}
//
//	m.dep.HeroData().FuncWithSend(hc.Id(), func(hero *entity.Hero, result herolock.LockResult) {
//		if realm := m.realmService.GetRealm(hero.BaseRegion()); realm == nil || !realm.CheckIsFucked(hc.Id()) {
//			hc.Send(guild.ERR_ASK_FOR_HELP_FAIL_NO_DANGER)
//			return
//		}
//		result.Add(guild.NewS2cAskForHelpMsg(i64.Int32(hero.BaseRegion()), int32(hero.BaseX()), int32(hero.BaseY())))
//	})
//}

//gogen:iface
func (m *GuildModule) ProcessConvene(proto *guild.C2SConveneProto, hc iface.HeroController) {

	var targetId int64
	if len(proto.Target) > 0 {
		ok := false
		targetId, ok = idbytes.ToId(proto.Target)
		if !ok {
			logrus.Debug("联盟官员召集，目标id无效")
			hc.Send(guild.ERR_CONVENE_FAIL_INVALID_TARGET)
			return
		}

		if targetId == hc.Id() {
			logrus.Debug("联盟官员召集，召集自己?")
			hc.Send(guild.ERR_CONVENE_FAIL_INVALID_TARGET)
			return
		}
	}

	var heroName string
	var baseX, baseY int
	if hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
		heroName = hero.Name()
		baseX, baseY = hero.BaseX(), hero.BaseY()
		return
	}) {
		hc.Send(guild.ERR_CONVENE_FAIL_SERVER_ERROR)
		return
	}

	m.processHeroInGuildMsg("Convene", hc,
		guild.ERR_CONVENE_FAIL_NOT_IN_GUILD,
		guild.ERR_CONVENE_FAIL_SERVER_ERROR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.convene(g, self, targetId, heroName, baseX, baseY)
		})
}

//gogen:iface c2s_collect_daily_guild_rank_prize
func (m *GuildModule) ProcessCollectDailyGuildRankPrize(hc iface.HeroController) {
	gid, ok := hc.LockGetGuildId();
	if !ok || gid == 0 {
		logrus.Debug("领取联盟日常排行奖励，没有联盟")
		hc.Send(guild.ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_NO_GUILD)
		return
	}
	var rank uint64
	var countryId uint64
	m.guildService.FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g == nil {
			hc.Send(guild.ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_NO_GUILD)
			return
		}
		rank = g.GetLastPrestigeRank()
		countryId = g.CountryId()
	})
	if rank <= 0 {
		logrus.Debug("领取联盟日常排行奖励，联盟没有上榜")
		hc.Send(guild.ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_NO_GUILD_RANK)
		return
	}
	data := m.datas.GetGuildRankPrizeData(rank)
	if data == nil {
		logrus.Debug("领取联盟日常排行奖励，配置错误")
		hc.Send(guild.ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_SERVER_ERROR)
		return
	}

	countryDestroyed := m.country.Country(countryId).IsDestroyed()

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.CollectedDailyGuildRankPrize() {
			logrus.Debug("领取联盟日常排行奖励，无法领取")
			result.Add(guild.ERR_COLLECT_DAILY_GUILD_RANK_PRIZE_FAIL_COLLECTED)
			return
		}

		prize := data.Prize
		if countryDestroyed {
			prize = data.CountryDestroyPrize
		}

		hctx := heromodule.NewContext(m.dep, operate_type.GuildCollectRankPrize)
		heromodule.AddPrize(hctx, hero, result, prize, m.time.CurrentTime())
		result.Add(guild.NewS2cCollectDailyGuildRankPrizeMsg(must.Marshal(prize.Encode())))

		hero.SetCollectedDailyGuildRankPrize(true)

		result.Ok()
	})
}

//gogen:iface c2s_view_daily_guild_rank
func (m *GuildModule) ProcessViewDailyGuildRank(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		cid := hero.CountryId()
		if cid == 0 {
			logrus.Debug("查看联盟日常排行，没有国家")
			result.Add(guild.ERR_VIEW_DAILY_GUILD_RANK_FAIL_NO_COUNTRY)
			return
		}
		result.Add(m.guildService.GetGuildPrestigeRankMsg(cid, m.time.CurrentTime()))

		result.Ok()
	})
}

//gogen:iface
func (m *GuildModule) ProcessAddRecommendMcBuild(proto *guild.C2SAddRecommendMcBuildProto, hc iface.HeroController) {
	mcId := u64.FromInt32(proto.McId)
	if m.dep.Mingc().Mingc(mcId) == nil {
		hc.Send(guild.ERR_ADD_RECOMMEND_MC_BUILD_FAIL_INVALID_MC_ID)
		return
	}

	m.processHeroInGuildMsg("AddRecommendMcBuild", hc,
		guild.ERR_ADD_RECOMMEND_MC_BUILD_FAIL_NO_GUILD,
		guild.ERR_ADD_RECOMMEND_MC_BUILD_FAIL_SERVER_ERR,
		func(guilds sharedguilddata.Guilds, g *sharedguilddata.Guild, self *sharedguilddata.GuildMember) (successMsg pbutil.Buffer, errMsg msg.ErrMsg, broadcastChanged bool) {
			return m.addRecommendMcBuild(g, self, mcId)
		})
}

//gogen:iface
func (m *GuildModule) ProcessViewTaskProgress(proto *guild.C2SViewTaskProgressProto, hc iface.HeroController) {
	gid, ok := hc.LockGetGuildId()
	if !ok || gid == 0 {
		logrus.Debug("查看联盟周任务进度，没有联盟")
		hc.Send(guild.ERR_VIEW_TASK_PROGRESS_FAIL_NO_GUILD)
		return
	}
	m.guildService.FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g == nil {
			hc.Send(guild.ERR_VIEW_TASK_PROGRESS_FAIL_NO_GUILD)
			return
		}
		if g.LevelData().Level < m.datas.GuildGenConfig().TaskOpenLevel {
			hc.Send(guild.ERR_VIEW_TASK_PROGRESS_FAIL_GUILD_LEVEL_LIMIT)
			return
		}
		if proto.GetVersion() == g.GetGuildTaskVersion() {
			hc.Send(guild.NewS2cViewTaskProgressMsg(g.GetGuildTaskVersion(), []*shared_proto.Int32Pair{}))
			return
		}
		hc.Send(g.GetWeeklyTasksMsg())
	})
}

//gogen:iface
func (m *GuildModule) ProcessCollectTaskPrizeProgress(proto *guild.C2SCollectTaskPrizeProto, hc iface.HeroController) {
	data := m.datas.GetGuildTaskData(u64.FromInt32(proto.GetTaskId()))
	if data == nil {
		logrus.Debug("领取联盟周任务进度奖励，没有该任务")
		hc.Send(guild.ERR_COLLECT_TASK_PRIZE_FAIL_INVALID_VALUE)
		return
	}
	stageIndex := int(proto.GetStage()) - 1
	if stageIndex < 0 || stageIndex >= len(data.Stages) {
		logrus.Debug("领取联盟周任务进度奖励，没有该任务进度")
		hc.Send(guild.ERR_COLLECT_TASK_PRIZE_FAIL_INVALID_VALUE)
		return
	}
	gid, ok := hc.LockGetGuildId()
	if !ok || gid == 0 {
		logrus.Debug("领取联盟周任务进度奖励，没有联盟")
		hc.Send(guild.ERR_COLLECT_TASK_PRIZE_FAIL_NO_GUILD)
		return
	}
	var errMsg pbutil.Buffer
	m.guildService.FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Debug("领取联盟周任务进度奖励，没有联盟")
			errMsg = guild.ERR_COLLECT_TASK_PRIZE_FAIL_NO_GUILD
			return
		}
		if g.GetGuildTaskProgress(server_proto.GuildTaskType(proto.TaskId)) < data.Stages[stageIndex] {
			logrus.Debug("领取联盟周任务进度奖励，阶段奖励未激活")
			errMsg = guild.ERR_COLLECT_TASK_PRIZE_FAIL_NO_PRIZE
			return
		}
	})
	if errMsg != nil {
		hc.Send(errMsg)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !hero.TrySetCollectedTaskStage(data.Id, proto.Stage) {
			logrus.Debug("领取联盟周任务进度奖励，阶段奖励已领取")
			result.Add(guild.ERR_COLLECT_TASK_PRIZE_FAIL_COLLECTED)
			return
		}
		hctx := heromodule.NewContext(m.dep, operate_type.GuildCollectTaskPrize)
		heromodule.AddPrize(hctx, hero, result, data.Prizes[stageIndex], m.time.CurrentTime())

		result.Add(guild.NewS2cCollectTaskPrizeMsg(proto.TaskId, proto.Stage))
		result.Ok()
	})
}

//gogen:iface
func (m *GuildModule) ProcessGuildChangeCountry(proto *guild.C2SGuildChangeCountryProto, hc iface.HeroController) {

	target := m.datas.GetCountryData(u64.FromInt32(proto.Country))
	if target == nil {
		logrus.Debug("联盟转国，target == nil")
		hc.Send(guild.ERR_GUILD_CHANGE_COUNTRY_FAIL_INVALID_COUNTRY)
		return
	}

	guildId, ok := hc.LockGetGuildId()
	if !ok {
		logrus.Debug("联盟转国，获取英雄的联盟id失败")
		hc.Send(guild.ERR_GUILD_CHANGE_COUNTRY_FAIL_NOT_IN_GUILD)
		return
	}

	if guildId == 0 {
		logrus.Debug("联盟转国，你没有联盟")
		hc.Send(guild.ERR_GUILD_CHANGE_COUNTRY_FAIL_NOT_IN_GUILD)
		return
	}

	ctime := m.time.CurrentTime()
	ctimeUnix := ctime.Unix()

	hctx := heromodule.NewContext(m.dep, operate_type.GuildChangeCountry)

	var kingGuild int64
	m.dep.HeroData().FuncNotError(m.dep.Country().King(hc.LockHeroCountry()), func(hero *entity.Hero) (heroChanged bool) {
		kingGuild = hero.GuildId()
		return
	})
	if kingGuild == guildId {
		hc.Send(guild.ERR_GUILD_CHANGE_COUNTRY_FAIL_IS_KING)
		return
	}

	var errMsg pbutil.Buffer
	var allMemberIds []int64
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Error("联盟转国，g == nil")
			errMsg = guild.ERR_GUILD_CHANGE_COUNTRY_FAIL_NOT_IN_GUILD
			return
		}

		if g.Country() == target {
			logrus.Debug("联盟转国，要转的国家跟现在国家一样")
			errMsg = guild.ERR_GUILD_CHANGE_COUNTRY_FAIL_SAME_COUNTRY
			return
		}

		if g.LeaderId() != hc.Id() {
			logrus.Debug("联盟转国，你不是盟主")
			errMsg = guild.ERR_GUILD_CHANGE_COUNTRY_FAIL_NOT_LEARDER
			return
		}

		if ctimeUnix < g.GetChangeCountryNextTime() {
			logrus.Debug("联盟转国，转国CD中")
			errMsg = guild.ERR_GUILD_CHANGE_COUNTRY_FAIL_COOLDOWN
			return
		}

		if g.GetChangeCountryWaitEndTime() > 0 {
			logrus.Debug("联盟转国，当前正在转国中")
			errMsg = guild.ERR_GUILD_CHANGE_COUNTRY_FAIL_EXIST
			return
		}

		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

			if !heromodule.TryReduceCost(hctx, hero, result, m.datas.GuildGenConfig().GuildChangeCountryCost) {
				logrus.Debugf("修改声望目标，消耗不足")
				errMsg = guild.ERR_GUILD_CHANGE_COUNTRY_FAIL_COST_NOT_ENOUGH
				return
			}
			result.Ok()
		}) {
			return
		}

		// 转国
		waitTime := ctime.Add(m.datas.GuildGenConfig().GuildChangeCountryWaitDuration)
		nextTime := ctime.Add(m.datas.GuildGenConfig().GuildChangeCountryCooldown)
		g.SetChangeCountry(target, waitTime.Unix(), nextTime.Unix())

		// 设置事件
		g.TryUpdateTarget(m.datas.GuildConfig(), ctime, shared_proto.GuildTargetType_GuildChangeCountry)

		g.SetChanged()
		allMemberIds = g.AllUserMemberIds()

		m.clearSelfGuildMsgCache(g.Id())
	})

	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(guild.NewS2cGuildChangeCountryMsg(proto.Country))

	if len(allMemberIds) > 0 {
		m.world.MultiSend(allMemberIds, guild.SELF_GUILD_CHANGED_S2C)
	}
}

//gogen:iface c2s_cancel_guild_change_country
func (m *GuildModule) ProcessCancelGuildChangeCountry(hc iface.HeroController) {

	guildId, ok := hc.LockGetGuildId()
	if !ok {
		logrus.Debug("取消联盟转国，获取英雄的联盟id失败")
		hc.Send(guild.ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_IN_GUILD)
		return
	}

	if guildId == 0 {
		logrus.Debug("取消联盟转国，你没有联盟")
		hc.Send(guild.ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_IN_GUILD)
		return
	}

	ctime := m.time.CurrentTime()

	var errMsg pbutil.Buffer
	var allMemberIds []int64
	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Error("取消联盟转国，g == nil")
			errMsg = guild.ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_IN_GUILD
			return
		}

		if g.LeaderId() != hc.Id() {
			logrus.Debug("取消联盟转国，你不是盟主")
			errMsg = guild.ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_LEARDER
			return
		}

		if g.GetChangeCountryTarget() == nil && g.GetChangeCountryWaitEndTime() == 0 {
			logrus.Debug("取消联盟转国，联盟没有处于转国中")
			errMsg = guild.ERR_CANCEL_GUILD_CHANGE_COUNTRY_FAIL_NOT_EXIST
			return
		}

		// 取消转国
		g.CancelChangeCountry()

		// 设置事件
		g.TryUpdateTarget(m.datas.GuildConfig(), ctime, shared_proto.GuildTargetType_GuildChangeCountry)

		g.SetChanged()
		allMemberIds = g.AllUserMemberIds()

		m.clearSelfGuildMsgCache(g.Id())
	})

	if errMsg != nil {
		hc.Send(errMsg)
		return
	}

	hc.Send(guild.CANCEL_GUILD_CHANGE_COUNTRY_S2C)

	if len(allMemberIds) > 0 {
		m.world.MultiSend(allMemberIds, guild.SELF_GUILD_CHANGED_S2C)
	}
}

var (
	showWorkshopFalseMsg = guild.NewS2cShowWorkshopNotExistMsg(false).Static()
	showWorkshopTrueMsg  = guild.NewS2cShowWorkshopNotExistMsg(true).Static()
)

//gogen:iface c2s_show_workshop_not_exist
func (m *GuildModule) ProcessShowWorkshopNotExist(hc iface.HeroController) {

	guildId, ok := hc.LockGetGuildId()
	if !ok {
		logrus.Debug("更新联盟工坊首次提醒，获取联盟失败")
		hc.Send(showWorkshopFalseMsg)
		return
	}

	if guildId == 0 {
		logrus.Debug("更新联盟工坊首次提醒，玩家没有联盟")
		hc.Send(showWorkshopFalseMsg)
		return
	}

	m.guildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.Errorf("更新联盟工坊首次提醒，g == nil")
			return
		}

		if member := g.GetMember(hc.Id()); member != nil {
			member.SetShowWorkshopNotExist(true)
		}
	})

	hc.Send(showWorkshopFalseMsg)
}
