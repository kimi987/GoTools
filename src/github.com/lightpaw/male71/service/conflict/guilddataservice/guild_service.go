package guilddataservice

import (
	"context"
	"fmt"
	"github.com/lightpaw/golang-lru"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/util/concurrent"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/event"
	"github.com/lightpaw/male7/util/i64"
	i64_map "github.com/lightpaw/male7/util/i64/concurrent"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/random"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"runtime/debug"
	"time"
	"github.com/lightpaw/male7/config/guild_data"
)

var (
	recommendInviteHeroVersionedValue = &RecommendInviteHeroValue{}
)

type RecommendInviteHeroValue struct{}

func (*RecommendInviteHeroValue) Version() uint64 {
	return 0
}

// 所有帮派内应该处理的事情都在这里处理
func NewGuildService(datas *config.ConfigDatas, timeService iface.TimeService, db iface.DbService, country iface.CountryService,
	heroService iface.HeroDataService, heroSnapshotService iface.HeroSnapshotService, xiongNuService iface.XiongNuService,
	guildSnapshotService iface.GuildSnapshotService, world iface.WorldService, chat iface.ChatService) *GuildService {
	m := &GuildService{}
	m.db = db
	m.heroService = heroService
	m.world = world
	m.heroSnapshotService = heroSnapshotService
	m.guildSnapshotService = guildSnapshotService
	m.xiongNuService = xiongNuService
	m.chat = chat
	m.datas = datas
	m.timeService = timeService
	m.guilds = sharedguilddata.NewGuilds(datas)

	m.selfGuildMsgCache = i64_map.NewVersionBufferCacheMapBuilder(m.genSelfGuildMsg).ExpireAfter(time.Minute).Build()

	// 加载数据
	var guilds []*sharedguilddata.Guild
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		guilds, err = db.LoadAllGuild(ctx)
		return
	})
	if err != nil {
		logrus.WithError(err).Panicf("加载帮派数据失败")
	}

	// 推荐入盟玩家列表
	recommendCacheSize := datas.GuildConfig().RecommendInviteHeroCount * 10
	m.initRecommendInviteHeroCache(u64.Int(recommendCacheSize))

	var maxGuildId int64
	for _, g := range guilds {
		m.guilds.Add(g)

		maxGuildId = i64.Max(maxGuildId, g.Id())
	}

	if !m.initGuildMemberCountry {
		m.initGuildMemberCountry = true

		// 将所有有联盟的玩家，玩家国家id改成跟联盟国家id一致
		count := 0
		m.guilds.Walk(func(g *sharedguilddata.Guild) {
			g.WalkMember(func(member *sharedguilddata.GuildMember) {
				if member.IsNpc() {
					return
				}
				heroService.FuncNotError(member.Id(), func(hero *entity.Hero) (heroChanged bool) {
					if hero.GuildId() != g.Id() {
						return false
					}

					if hero.CountryId() == g.CountryId() {
						return false
					}

					hero.SetCountryId(g.CountryId())
					count++

					return true
				})
			})
		})

		logrus.WithField("count", count).Info("处理联盟成员国家不一致的情况")
	}

	heroJoinGuildMap := make(map[int64]*sharedguilddata.Guild, len(guilds))
	m.guilds.Walk(func(g *sharedguilddata.Guild) {
		g.WalkMember(func(member *sharedguilddata.GuildMember) {
			if member.IsNpc() {
				return
			}

			g1 := heroJoinGuildMap[member.Id()]
			if g1 == nil {
				heroJoinGuildMap[member.Id()] = g
			} else {
				logrus.WithField("guild1", fmt.Sprintf("%d-%v", g1.Id(), g1.Name())).
					WithField("guild2", fmt.Sprintf("%d-%v", g.Id(), g.Name())).
					WithField("member", member.Id()).
					WithField("proto", fmt.Sprintf("%v", member.NpcProto())).
					Errorln("存在有玩家同时出现在两个联盟的严重bug!")

				heroService.FuncNotError(member.Id(), func(hero *entity.Hero) (heroChanged bool) {
					if hero == nil {
						logrus.WithField("hero_id", member.Id()).Error("玩家数据没找到")
						return
					}

					logrus.WithField("hero_id", hero.Id()).
						WithField("hero_name", hero.Name()).
						WithField("guild_id", hero.GuildId()).Error("玩家联盟数据")

					if hero.GuildId() != g.Id() {
						g.RemoveMember(member.Id(), 0)
						delete(heroJoinGuildMap, member.Id())
					}

					if hero.GuildId() != g1.Id() {
						g1.RemoveMember(member.Id(), 0)
						delete(heroJoinGuildMap, member.Id())
					}

					switch hero.GuildId() {
					case g.Id():
						heroJoinGuildMap[member.Id()] = g
					case g1.Id():
						heroJoinGuildMap[member.Id()] = g1
					case 0:
						// nothing
					default:
						heroGuild := m.guilds.Get(hero.GuildId())
						if heroGuild.GetMember(hero.Id()) == nil {
							hero.SetGuild(0)
							heroChanged = true
						}
					}
					return
				})
			}
		})
	})

	guildDailyRankLimit := u64.FromInt(len(datas.GetGuildRankPrizeDataArray()))
	m.guilds.Walk(func(g *sharedguilddata.Guild) {

		if g.MemberCount() <= 0 {
			logrus.WithField("guild", fmt.Sprintf("%d-%v", g.Id(), g.Name())).Errorln("存在联盟中人数为0的严重bug，删除联盟")
			m.guilds.Remove(g.Id())
			if err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
				return db.DeleteGuild(ctx, g.Id())
			}); err != nil {
				logrus.WithError(err).Error("删除联盟人数为0的联盟，删除失败")
			}

			return
		}

		if g.GetMember(g.LeaderId()) == nil {
			g.ResetLeader()
		}

		m.UpdateSnapshot(g)
		country.AddPrestige(g.Country().Id, g.GetPrestige())

		if g.GetLastPrestigeRank() > guildDailyRankLimit {
			logrus.WithField("guild", fmt.Sprintf("%d-%v", g.Id(), g.Name())).Errorln("有联盟昨日声望排名超出配置，清零")
			g.SetLastPrestigeRank(0)
		}
	})

	if len(guilds) <= 0 {
		// 创建A类帮派

		ctime := timeService.CurrentTime()
		for _, t := range datas.GetNpcGuildTemplateArray() {
			if t.RejectUserJoin {

				maxGuildId++
				newGuildId := maxGuildId

				newGuild := sharedguilddata.NewNpcGuild(t, newGuildId, t.Name, t.FlagName, ctime, datas)

				// 创建帮派
				guildBytes, err := newGuild.Marshal()
				if err != nil {
					logrus.WithError(err).Panicf("开服创建A类Npc联盟，Guild.Marshal 报错")
				}

				err = ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
					return m.db.CreateGuild(ctx, newGuildId, guildBytes)
				})
				if err != nil {
					logrus.WithError(err).Errorf("开服创建A类Npc联盟，DB.CreateGuild 报错")
				}

				// 创建成功，添加到guilds
				m.guilds.Add(newGuild)
				m.UpdateSnapshot(newGuild)
			}
		}
	}

	m.queue = event.NewEventQueue(2048, 3*time.Second, "GuildEvent")
	m.logQueue = event.NewFuncQueue(1024, "GuildLog")

	return m
}

//gogen:iface
type GuildService struct {
	db    iface.DbService
	world iface.WorldService

	heroService          iface.HeroDataService
	heroSnapshotService  iface.HeroSnapshotService
	guildSnapshotService iface.GuildSnapshotService
	xiongNuService       iface.XiongNuService
	chat                 iface.ChatService
	datas                iface.ConfigDatas
	timeService          iface.TimeService

	guilds sharedguilddata.Guilds

	selfGuildMsgCache i64_map.I64BufferMap

	queue *event.EventQueue

	logQueue *event.FuncQueue

	rankFunc sharedguilddata.GetGuildRankFunc

	dailyRankMsgFunc sharedguilddata.GenerateRankMsgFunc

	recommendInviteHeros *lru.Cache

	initGuildMemberCountry bool

	// 国家声望
}

func (m *GuildService) SetGuildRankFunc(f1 sharedguilddata.GetGuildRankFunc, f2 sharedguilddata.GenerateRankMsgFunc) {
	m.rankFunc = f1
	m.dailyRankMsgFunc = f2
}

func (m *GuildService) initRecommendInviteHeroCache(recommendCacheSize int) {
	// 邀请推荐列表
	recommendCache, err := lru.New(recommendCacheSize)
	if err != nil {
		logrus.WithError(err).Panic("new guild recommendInviteHeros err")
	}

	m.recommendInviteHeros = recommendCache
	var recommendProtoBytes []byte
	err = ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		recommendProtoBytes, err = m.db.LoadKey(ctx, server_proto.Key_GuildRcmdHeros)
		return
	})
	if err != nil {
		logrus.WithError(err).Panic("guild recommendInviteHeros db err")
	}

	recommendProto := &server_proto.InviteGuildRecommendHerosProto{}
	if len(recommendProtoBytes) > 0 {
		err = recommendProto.Unmarshal(recommendProtoBytes)
		if err != nil {
			logrus.WithError(err).Panic("guild recommendInviteHeros Unmarshal err")
		}
	}
	for _, heroId := range recommendProto.HeroId {
		m.recommendInviteHeros.Add(heroId, recommendInviteHeroVersionedValue)
	}
	m.initGuildMemberCountry = recommendProto.InitGuildMemberCountry
}

func (m *GuildService) Close() {

	// 停止帮派线程
	m.queue.Stop()

	// 保存数据
	m.guilds.Walk(func(g *sharedguilddata.Guild) {
		ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
			m.save(ctx, g)
			return
		})
	})

	// 保存联盟入盟推荐列表
	ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		m.saveRecommendInviteHero(ctx)
		return
	})
}

func (m *GuildService) SaveChangedGuild() {

	logrus.Debugf("定时保存变化的帮派")

	// 获取所有changed的帮派对象，保存到数据库
	m.Func(func(guilds sharedguilddata.Guilds) {
		guilds.Walk(func(g *sharedguilddata.Guild) {
			if g.SetFalseIfChanged() {
				ctxfunc.Timeout3s(func(ctx context.Context) error {
					m.save(ctx, g)
					return nil
				})
			}
		})
	})

	ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
		m.saveRecommendInviteHero(ctx)
		return
	})

}

func (m *GuildService) save(ctx context.Context, g *sharedguilddata.Guild) {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("err", r).WithField("stack", string(debug.Stack())).Errorf("GuildService.saveObject recovered from panic. SEVERE!!!")
			metrics.IncPanic()
		}
	}()

	data, err := g.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("GuildService guild.marshal 失败")
		return
	}

	err = m.db.SaveGuild(ctx, g.Id(), data)
	if err != nil {
		logrus.WithError(err).Errorf("GuildService guild.save 失败")
		return
	}
}

// 保存邀请入盟推荐列表
func (m *GuildService) saveRecommendInviteHero(ctx context.Context) {
	recommendProto := &server_proto.InviteGuildRecommendHerosProto{}
	recommendProto.HeroId = m.recommendInviteHeros.Keys()
	recommendProto.InitGuildMemberCountry = m.initGuildMemberCountry

	recommendProtoBytes, err := recommendProto.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("GuildService guild.save recommendInviteHeros Marshal 失败")
		return
	}

	err = m.db.SaveKey(ctx, server_proto.Key_GuildRcmdHeros, recommendProtoBytes)
	if err != nil {
		logrus.WithError(err).Errorf("GuildService guild.save recommendInviteHeros 失败")
		return
	}
}

func (m *GuildService) GetGuildIdByName(name string) int64 {
	return m.guilds.GetIdByName(name)
}

func (m *GuildService) GetGuildIdByFlagName(flagName string) int64 {
	return m.guilds.GetIdByFlagName(flagName)
}

func (m *GuildService) genSelfGuildMsg(guildId int64, version uint64) (toSend pbutil.Buffer, err error) {
	if !m.TimeoutFunc(func(guilds sharedguilddata.Guilds) {

		g := guilds.Get(guildId)
		if g == nil {
			err = sharedguilddata.ErrNotExist
			return
		}

		proto := g.EncodeClient(true, m.heroSnapshotService, m.rankFunc, m.xiongNuService.IsTodayStarted)

		toSend = guild.NewS2cSelfGuildMsg(concurrent.I32Version(version), must.Marshal(proto))
	}) {
		// timeout
		err = sharedguilddata.ErrTimeout
	}

	return
}

func (m *GuildService) ClearSelfGuildMsgCache(guildId int64) {
	m.selfGuildMsgCache.Clear(guildId)
}

func (m *GuildService) SelfGuildMsgCache() i64_map.I64BufferMap {
	return m.selfGuildMsgCache
}

func (m *GuildService) RegisterCallback(callback guildsnapshotdata.Callback) {
	m.guildSnapshotService.RegisterCallback(callback)
}

func (m *GuildService) GetSnapshot(id int64) *guildsnapshotdata.GuildSnapshot {
	return m.guildSnapshotService.GetSnapshot(id)
}

func (m *GuildService) UpdateSnapshot(g *sharedguilddata.Guild) *guildsnapshotdata.GuildSnapshot {
	s := g.NewSnapshot()
	m.guildSnapshotService.UpdateSnapshot(s)
	return s
}

func (m *GuildService) RemoveSnapshot(id int64) {
	m.guildSnapshotService.RemoveSnapshot(id)
}

func (m *GuildService) FuncGuild(id int64, f sharedguilddata.Func) {
	m.queue.TimeoutFunc(true, func() {
		g := m.guilds.Get(id)
		f(g)
	})

	return
}

func (m *GuildService) TimeoutFunc(f sharedguilddata.Funcs) (ok bool) {
	ok = m.queue.TimeoutFunc(true, func() {
		f(m.guilds)
	})

	return
}

//func (m *GuildService) FuncDontWait(f sharedguilddata.Funcs) (ok bool) {
//	return m.funcs(false, f)
//}
//
//func (m *GuildService) funcs(waitResult bool, f sharedguilddata.Funcs) (ok bool) {
//	ok = m.queue.TryTimeoutFunc(waitResult, func() {
//		f(m.guilds)
//	})
//
//	return
//}

func (m *GuildService) Func(f sharedguilddata.Funcs) (ok bool) {
	ok = m.queue.Func(true, func() {
		f(m.guilds)
	})

	return
}

func (m *GuildService) AddLog(guildId int64, proto *shared_proto.GuildLogProto) {

	if guildId == 0 {
		return
	}

	g := m.guildSnapshotService.GetSnapshot(guildId)
	if g == nil {
		return
	}

	m.logQueue.MustFunc(func() {
		err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			return m.db.InsertGuildLog(ctx, guildId, proto)
		})
		if err != nil {
			logrus.WithError(err).Error("添加联盟日志到DB，添加失败")
			return
		}

		// 成功后，将proto广播给联盟
		m.world.MultiSend(g.UserMemberIds, guild.NewS2cAddGuildLogMarshalMsg(proto))

		if d := m.datas.GetGuildLogData(proto.DataId); d != nil {
			if d.SendChat {
				// 同步联盟聊天
				m.chat.SysChat(0, guildId, shared_proto.ChatType_ChatGuild, proto.Text, shared_proto.ChatMsgType_ChatMsgGuildLog, true, true, true, false)
			}
		}

	})

	return
}

func (m *GuildService) AddLogWithMemberIds(guildId int64, memberIds []int64, proto *shared_proto.GuildLogProto) {

	if guildId == 0 {
		return
	}

	m.logQueue.MustFunc(func() {
		err := ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			return m.db.InsertGuildLog(ctx, guildId, proto)
		})
		if err != nil {
			logrus.WithError(err).Error("添加联盟日志到DB，添加失败")
			return
		}

		// 成功后，将proto广播给联盟
		m.world.MultiSend(memberIds, guild.NewS2cAddGuildLogMarshalMsg(proto))
	})

	return
}

func (m *GuildService) CheckAndAddRecommendInviteHeros(heroId int64) {
	// 条件， TODO 优先级
	hero := m.heroSnapshotService.Get(heroId)
	if hero == nil {
		return
	}
	if hero.GuildId > 0 {
		return
	}

	m.recommendInviteHeros.Add(heroId, recommendInviteHeroVersionedValue)
}

func (m *GuildService) AddRecommendInviteHeros(heroId int64) {
	m.recommendInviteHeros.Add(heroId, recommendInviteHeroVersionedValue)
}

func (m *GuildService) RemoveRecommendInviteHero(heroId int64) {
	m.recommendInviteHeros.Remove(heroId)
}

func (m *GuildService) RecommendInviteHeroList(countryId, size uint64, ignores []int64) []*snapshotdata.HeroSnapshot {
	sizeInt := u64.Int(size)
	heroIds := m.recommendInviteHeros.Keys()
	random.MixI64Array(heroIds)

	reserveHeros := make([]*snapshotdata.HeroSnapshot, 0)
	heros := make([]*snapshotdata.HeroSnapshot, 0)
	for _, heroId := range heroIds {
		// 过滤掉已经邀请的列表
		if i64.Contains(ignores, heroId) {
			continue
		}

		hero := m.heroSnapshotService.Get(heroId)
		if hero == nil || hero.CountryId != countryId || hero.GuildId != 0 {
			continue
		}

		// 优先条件
		if m.world.IsOnline(heroId) {
			heros = append(heros, hero)
		} else {
			reserveHeros = append(reserveHeros, hero)
		}

		if len(heros) >= sizeInt {
			return heros
		}
	}

	if len(heros) < sizeInt {
		if needCount := sizeInt - len(heros); needCount > 0 {
			if n := imath.Min(len(reserveHeros), needCount); n > 0 {
				heros = append(heros, reserveHeros[:n]...)
			}
		}
	}

	return heros
}

func (m *GuildService) UpdateGuildHeroSnapshot(g *sharedguilddata.Guild) {
	for _, heroId := range g.AllUserMemberIds() {
		heroSnapshot := m.heroSnapshotService.GetFromCache(heroId)
		if heroSnapshot == nil {
			continue
		}

		heroSnapshot.ClearProto()
	}
}

func (m *GuildService) AddHufu(addHufu uint64, heroId, gid int64, heroName, heroHead string) (succ bool) {
	if gid <= 0 {
		return
	}
	var newHufu uint64
	var ids []int64
	m.FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}
		newHufu = g.AddHufu(addHufu)
		if member := g.GetMember(heroId); member != nil {
			member.AddHufu(addHufu)
		}
		g.SetChanged()
		ids = g.AllUserMemberIds()
		succ = true
	})
	if !succ {
		return
	}

	// 清掉缓存
	m.ClearSelfGuildMsgCache(gid)
	if len(ids) > 0 {
		// 发送给联盟内所有的成员
		m.world.MultiSend(ids, guild.NewS2cUpdateHufuMsg(u64.Int32(newHufu)))
	}

	return
}

// 根据国家id获取声望排名列表的消息
func (m *GuildService) GetGuildPrestigeRankMsg(cid uint64, ctime time.Time) pbutil.Buffer {
	return m.guilds.PrestigeRankMsg(cid, ctime, m.dailyRankMsgFunc)
}

func (m *GuildService) Broadcast(gid int64, msg pbutil.Buffer) {
	m.FuncGuild(gid, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}
		for _, hid := range g.AllUserMemberIds() {
			m.world.Send(hid, msg)
		}
	})
}

// 增加联盟任务进度
func (m *GuildService) AddGuildTaskProgress(guildId int64, data *guild_data.GuildTaskData, addProgress uint64) {
	if data == nil {
		return
	}
	var allMemberIds []int64
	var stageIndex int
	m.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}
		if g.LevelData().Level < m.datas.GuildGenConfig().TaskOpenLevel {
			return
		}
		if g.AddGuildTaskProgress(data, addProgress) {
			allMemberIds = g.AllUserMemberIds()
			stageIndex = g.GetGuildTaskStageIndex(data.TaskType)
		}
	})
	if len(allMemberIds) > 0 {
		m.world.MultiSend(allMemberIds, guild.NewS2cNoticeTaskStageUpdateMsg(int32(data.TaskType), int32(stageIndex)))
	}
}


func (m *GuildService) GetGuildFlagName(gid int64) string {
	if g := m.guildSnapshotService.GetSnapshot(gid); g != nil {
		return m.datas.MiscConfig().FlagHeroName.FormatIgnoreEmpty(g.FlagName, g.Name)
	}
	return ""
}

// 筛选推荐联盟
func (m *GuildService) RecommendGuildList(hero *snapshotdata.HeroSnapshot) (guildIds []int64) {
	guildList := m.guilds.ListGuild(hero.Level, hero.BaiZhanJunXianLevel, hero.TowerMaxFloor, hero.CountryId)
	for _, guild := range guildList {
		guildIds = append(guildIds, guild.Id())
	}
	return
}
