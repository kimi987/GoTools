package rank

import (
	"context"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/rank_data"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/module/rank/ranklist"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"sort"
	"sync"
	"time"
	"github.com/lightpaw/male7/gen/pb/rank"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/gen/pb/guild"
)

// 排行榜
// 分成单排行榜类型和多类型排行榜
// 玩家等级排行榜，属于单排行榜类型

func NewRankModule(configDatas iface.ConfigDatas, dbService iface.DbService, guildService iface.GuildService, heroSnapshotService iface.HeroSnapshotService,
	timeService iface.TimeService, serverStartStopTimeService iface.ServerStartStopTimeService, baiZhanService iface.BaiZhanService,
	tssClient iface.TssClient) *RankModule {
	m := &RankModule{}

	m.miscData = configDatas.RankMiscData()
	m.configDatas = configDatas
	m.dbService = dbService
	m.guildService = guildService
	m.heroSnapshotService = heroSnapshotService
	m.timeService = timeService
	m.serverStartStopTimeService = serverStartStopTimeService
	m.baiZhanService = baiZhanService
	m.tssClient = tssClient
	m.holders = make([]ranklist.RankHolder, len(shared_proto.RankType_name))

	m.loopExitNotify = make(chan struct{})
	m.closeNotify = make(chan struct{})

	m.heroNameQuery = func(name string) (heroId int64) {
		ctxfunc.NetTimeout2s(func(ctx context.Context) (err error) {
			heroId, err = m.dbService.HeroId(ctx, name)
			if err != nil {
				logrus.WithError(err).Debugln("通过玩家名字获取玩家id报错")
				return
			}

			return
		})

		return
	}

	// 国家排行榜
	countryRankHolder := ranklist.NewCountryRankHolder(m.miscData.MaxRankCount, func(name string) (key int64) {
		for _, c := range configDatas.GetCountryDataArray() {
			if c.Name == name {
				return int64(c.Id)
			}
		}
		return 0
	})
	m.addRankHolder(shared_proto.RankType_Country, countryRankHolder)

	// 联盟排行榜
	guildRankHolder := ranklist.NewGuildRankHolder(configDatas, m.miscData.MaxRankCount, func(name string) (key int64) {
		return m.guildService.GetGuildIdByName(name)
	})
	m.addRankHolder(shared_proto.RankType_Guild, guildRankHolder)

	// 千重楼排行榜
	towerRankHolder := ranklist.NewTowerRankHolder(m.miscData.MaxRankCount, m.heroNameQuery)
	m.addRankHolder(shared_proto.RankType_Tower, towerRankHolder)

	// 百战千军排行榜
	baiZhanRankHolder := ranklist.NewBaiZhanRankHolder(m.configDatas, m.heroNameQuery)
	m.addRankHolder(shared_proto.RankType_BaiZhan, baiZhanRankHolder)

	// 成就星数排行榜
	starTaskRankHolder := ranklist.NewStarTaskRankHolder(m.miscData.MaxRankCount, m.heroNameQuery)
	m.addRankHolder(shared_proto.RankType_RankStarTask, starTaskRankHolder)

	for _, c := range configDatas.GetCountryDataArray() {
		m.AddOrUpdateRankObj(ranklist.NewCountryRankObj(m.guildService.GetSnapshot, m.heroSnapshotService.Get, 0, c.Name, 0, 0, time.Time{}))
	}

	var guildRankObjs []rankface.RankObj
	m.guildService.Func(func(guilds sharedguilddata.Guilds) {
		guilds.Walk(func(g *sharedguilddata.Guild) {
			npcTemplate := g.GetNpcTemplate()
			if npcTemplate != nil && npcTemplate.RejectUserJoin {
				// A类联盟不加入联盟排行榜
				return
			}

			obj := ranklist.NewGuildRankObj(guildService.GetSnapshot, heroSnapshotService.Get, g)

			guildRankObjs = append(guildRankObjs, obj)
		})
	})
	// 排好序，速度快一点
	sort.Sort(rankface.RankObjSlice(guildRankObjs))
	for _, obj := range guildRankObjs {
		m.AddOrUpdateRankObj(obj)
	}

	m.load()

	go call.CatchLoopPanic(m.loop, "排行榜 loop")

	return m
}

func (m *RankModule) loop() {
	// 10分钟保存一次
	saveTick := time.NewTicker(10 * time.Minute)

	defer close(m.loopExitNotify)

	for {
		select {
		case <-saveTick.C:
			m.save()
		case <-m.closeNotify:
			m.save()
			return // quit loop
		}
	}
}

//gogen:iface
type RankModule struct {
	dbService                  iface.DbService
	guildService               iface.GuildService
	heroSnapshotService        iface.HeroSnapshotService
	timeService                iface.TimeService
	serverStartStopTimeService iface.ServerStartStopTimeService
	baiZhanService             iface.BaiZhanService
	tssClient                  iface.TssClient
	miscData                   *rank_data.RankMiscData
	configDatas                iface.ConfigDatas
	holders                    []ranklist.RankHolder
	heroNameQuery              rankface.NameQueryFunc

	loopExitNotify chan struct{}
	closeNotify    chan struct{}
	closeOnce      sync.Once
}

func (m *RankModule) Close() {
	m.closeOnce.Do(func() {
		close(m.closeNotify)
	})
	<-m.loopExitNotify
}

func (m *RankModule) save() {
	proto := &server_proto.RankServerProto{}

	for _, holder := range m.holders {
		holder.LockFunc(func(h ranklist.LockedRankHolder) {
			h.Walk(func(list rankface.RankList) {
				list.EncodeServer(proto)
			})
		})
	}

	// 保存
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		return m.dbService.SaveKey(ctx, server_proto.Key_Rank, must.Marshal(proto))
	})
	if err != nil {
		logrus.WithError(err).Errorf("保存排行榜数据出错")
	}
}

func (m *RankModule) load() {
	// 是正常关服
	var data []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		data, err = m.dbService.LoadKey(ctx, server_proto.Key_Rank)
		return
	})
	if err != nil {
		logrus.WithError(err).Panicf("RankModule.load 读取数据报错")
	}

	proto := &server_proto.RankServerProto{}

	if len(data) > 0 {
		err = proto.Unmarshal(data)
		if err != nil {
			logrus.WithError(err).Panicf("RankModule.load 解析Proto出错")
		}
	}

	if m.serverStartStopTimeService.IsNormalStop() {
		initRank := func(rankType shared_proto.RankType, parse func(proto *server_proto.RankServerProto) []rankface.RankObj) {
			objs := parse(proto)

			sort.Sort(rankface.RankObjSlice(objs))
			holder := m.rankHolder(rankType)
			holder.LockFunc(func(h ranklist.LockedRankHolder) {
				h.AddOrUpdate(objs...)
			})
		}

		// 千重楼数据
		initRank(shared_proto.RankType_Tower, func(proto *server_proto.RankServerProto) []rankface.RankObj {
			var rankObjs []rankface.RankObj
			for _, p := range proto.Tower {
				rankObjs = append(rankObjs, ranklist.NewTowerRankObj(m.heroSnapshotService.Get, p.HeroId, p.MaxFloor, timeutil.Unix64(p.Time)))
			}
			return rankObjs
		})

		// 成就星数数据
		initRank(shared_proto.RankType_RankStarTask, func(proto *server_proto.RankServerProto) []rankface.RankObj {
			var rankObjs []rankface.RankObj
			for _, p := range proto.StarTask {
				rankObjs = append(rankObjs, ranklist.NewStarTaskRankObj(m.heroSnapshotService.Get, p.HeroId, p.Star, timeutil.Unix64(p.Time)))
			}
			return rankObjs
		})

	} else {
		logrus.Debugln("服务器非正常关闭，重新加载玩家千重楼数据")

		// 不是正常关服
		var heroes []*entity.Hero
		err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
			heroes, err = m.dbService.LoadAllHeroData(ctx)
			return
		})

		if err != nil {
			logrus.WithError(err).Panicf("RankModule.loadAbnormal 读取数据报错")
		}

		initRank := func(rankType shared_proto.RankType, parse func(heros []*entity.Hero) []rankface.RankObj) {
			objs := parse(heroes)

			sort.Sort(rankface.RankObjSlice(objs))
			holder := m.rankHolder(rankType)
			holder.LockFunc(func(h ranklist.LockedRankHolder) {
				h.AddOrUpdate(objs...)
			})
		}

		// 千重楼数据
		initRank(shared_proto.RankType_Tower, func(heros []*entity.Hero) []rankface.RankObj {
			var rankObjs []rankface.RankObj
			for _, hero := range heroes {
				tower := hero.Tower()
				if tower.HistoryMaxFloor() <= 0 {
					continue
				}

				rankObjs = append(rankObjs, ranklist.NewTowerRankObj(m.heroSnapshotService.Get, hero.Id(), tower.HistoryMaxFloor(), tower.HistoryMaxFloorTime()))
			}
			return rankObjs
		})

		// 成就星数数据
		initRank(shared_proto.RankType_RankStarTask, func(heros []*entity.Hero) []rankface.RankObj {
			var rankObjs []rankface.RankObj
			for _, hero := range heroes {
				star := hero.TaskList().AchieveTaskList().TotalStar()
				if star <= 0 {
					continue
				}

				rankObjs = append(rankObjs, ranklist.NewStarTaskRankObj(m.heroSnapshotService.Get, hero.Id(), star, time.Time{}))
			}
			return rankObjs
		})

	}

	// 百战数据
	baiZhanObjs := make([]rankface.RankObj, 0, len(proto.BaiZhan))
	for _, p := range proto.BaiZhan {
		baiZhanObjs = append(baiZhanObjs, ranklist.NewBaiZhanRankObj(m.heroSnapshotService.Get, m.baiZhanService.GetPoint, p.HeroId,
			m.configDatas.JunXianLevelData().Must(p.JunXianLevel), m.configDatas.JunXianLevelData().Must(p.LastJunXianLevel),
			timeutil.Unix64(p.Time), p.Point, p.FightAmount))
	}
	sort.Sort(rankface.RankObjSlice(baiZhanObjs))
	baiZhanHolder := m.rankHolder(shared_proto.RankType_BaiZhan)
	baiZhanHolder.LockFunc(func(h ranklist.LockedRankHolder) {
		h.AddOrUpdate(baiZhanObjs...)
	})

	// 轩辕会武数据
	xuanyObjs := make([]rankface.RankObj, 0, len(proto.Xuanyuan))
	for i, p := range proto.Xuanyuan {
		rank := uint64(i + 1)
		xuanyObjs = append(xuanyObjs, ranklist.NewXuanyRankObj(
			m.heroSnapshotService.Get, p.HeroId, p.Point, p.Win, p.Lose, rank))
	}
	xuanyHolder := ranklist.NewXuanyRankHolder(m.miscData.MaxRankCount, m.heroNameQuery, xuanyObjs)
	m.addRankHolder(shared_proto.RankType_RankXuanyuan, xuanyHolder)
}

func (m *RankModule) addRankHolder(rankType shared_proto.RankType, rankHolder ranklist.RankHolder) {
	m.holders[rankType] = rankHolder
}

func (m *RankModule) rankHolder(rankType shared_proto.RankType) ranklist.RankHolder {
	return m.holders[rankType]
}

func (m *RankModule) UpdateBaiZhanRankList(objs []rankface.RankObj) {
	sort.Sort(rankface.RankObjSlice(objs))
	holder := ranklist.NewBaiZhanRankHolder(m.configDatas, m.heroNameQuery)

	holder.LockFunc(func(h ranklist.LockedRankHolder) {
		h.AddOrUpdate(objs...)
	})

	// 直接替换掉
	m.addRankHolder(shared_proto.RankType_BaiZhan, holder)
}

func (m *RankModule) UpdateXuanyRankList(rankArray []rankface.RankObj) {
	holder := ranklist.NewXuanyRankHolder(m.miscData.MaxRankCount, m.heroNameQuery, rankArray)
	// 直接替换掉
	m.addRankHolder(shared_proto.RankType_RankXuanyuan, holder)
}

func (m *RankModule) AddOrUpdateRankObj(obj rankface.RankObj) {
	holder := m.rankHolder(obj.RankType())
	if holder == nil {
		logrus.WithField("obj", fmt.Sprintf("%+v", obj)).Errorln("没有找到排行榜")
		return
	}

	holder.LockFunc(func(h ranklist.LockedRankHolder) {
		h.AddOrUpdate(obj)
	})
}

func (m *RankModule) RemoveRankObj(rankType shared_proto.RankType, key int64) {
	holder := m.rankHolder(rankType)
	if holder == nil {
		logrus.WithField("rank type", rankType).Errorln("没有找到排行榜")
		return
	}

	holder.LockFunc(func(h ranklist.LockedRankHolder) {
		h.Remove(key)
	})
}

// 百战千军类型的特殊排行榜不能从这里获取
func (m *RankModule) SingleRRankListFunc(rankType shared_proto.RankType, f rankface.RRankListFunc) bool {
	holder := m.rankHolder(rankType)
	switch tp := holder.(type) {
	case ranklist.SingleRankHolder:
		f(tp.RRankList())
		return true
	default:
		logrus.WithField("rank type", rankType).Errorln("没找到对应的排行榜类型数据")
		return false
	}
}

func (m *RankModule) SubTypeRRankListFunc(rankType shared_proto.RankType, subType uint64, f rankface.RRankListFunc) bool {
	holder := m.rankHolder(rankType)
	switch tp := holder.(type) {
	case ranklist.SubTypeRankHolder:
		list := tp.RRankList(subType)
		if list == nil {
			logrus.WithField("rank type", rankType).WithField("subType", subType).Errorln("没找到对应的Sub排行榜类型数据")
			return false
		}

		f(list)
		return true
	default:
		logrus.WithField("rank type", rankType).Errorln("没找到对应的Sub排行榜类型数据")
		return false
	}
}

//gogen:iface c2s_request_rank
func (m *RankModule) ProcessRequestRank(proto *rank.C2SRequestRankProto, hc iface.HeroController) {
	if _, ok := shared_proto.RankType_name[proto.GetRankType()]; !ok {
		logrus.WithField("proto", proto).Debugln("未知的排行榜类型")
		hc.Send(rank.ERR_REQUEST_RANK_FAIL_UNKNOWN_RANK_TYPE)
		return
	}

	rankType := shared_proto.RankType(proto.GetRankType())

	holder := m.rankHolder(rankType)
	if holder == nil {
		// 发消息
		logrus.WithField("rank type", rankType).Debugln("不存在榜单")
		hc.Send(rank.ERR_REQUEST_RANK_FAIL_SERVER_ERROR)
		return
	}

	var queryKey int64

	if len(proto.Name) > 0 {
		if !m.tssClient.TryCheckName("排行榜搜索", hc, proto.Name, guild.ERR_CREATE_GUILD_FAIL_SENSITIVE_WORDS, rank.ERR_REQUEST_RANK_FAIL_SERVER_ERROR) {
			return
		}

		// 查询
		queryKey = holder.NameQuery(proto.Name)

		if queryKey == 0 {
			logrus.WithField("rank type", rankType).WithField("proto", proto).Debugln("没有根据名字找到排行榜的目标")
			hc.Send(rank.ERR_REQUEST_RANK_FAIL_TARGET_NOT_FOUND)
			return
		}

		proto.Self = false
	} else if proto.Self {
		queryKey = holder.SelfKey(hc)
	}

	subType := proto.JunXianLevel
	if proto.SubType != 0 {
		subType = proto.SubType
	}

	holder.RLockFunc(func(h ranklist.RLockedRankHolder) {
		var obj rankface.RankObj

		// 找
		var list rankface.RRankList
		if queryKey != 0 {
			switch tp := holder.(type) {
			case ranklist.SingleRankHolder:
				obj = tp.RankObjQuery(queryKey)
				list = tp.RRankList()
			case ranklist.SubTypeRankHolder:
				list, obj = tp.RankObjQuery(u64.FromInt32(subType), queryKey)
			default:
				logrus.WithField("holder", fmt.Sprintf("%+v", holder)).Errorln("未处理的排行榜类型")
				hc.Send(rank.ERR_REQUEST_RANK_FAIL_SERVER_ERROR)
				return
			}
		}

		var startRank = uint64(proto.StartCount)
		if obj == nil {
			if len(proto.Name) > 0 {
				logrus.WithField("rank type", rankType).WithField("proto", proto).Debugln("目标没在榜")
				hc.Send(rank.ERR_REQUEST_RANK_FAIL_TARGET_NOT_IN_RANK_LIST)
				return
			}

			// 处理咯
			switch tp := holder.(type) {
			case ranklist.SingleRankHolder:
				list = tp.RRankList()
			case ranklist.SubTypeRankHolder:
				list = tp.RRankList(u64.FromInt32(subType))
			default:
				logrus.WithField("holder", fmt.Sprintf("%+v", holder)).Errorln("未处理的排行榜类型")
				hc.Send(rank.ERR_REQUEST_RANK_FAIL_SERVER_ERROR)
				return
			}

			proto.Self = false
		} else {
			if list == nil {
				list = h.RRankListByObj(obj)
			}
			startRank = u64.Sub(obj.Rank(), m.miscData.RankCountPerPage>>1) // 把自己居中显示
		}

		if list == nil {
			logrus.WithField("rank type", rankType).WithField("proto", proto).Debugln("目标没在榜")
			hc.Send(rank.ERR_REQUEST_RANK_FAIL_SERVER_ERROR)
			return
		}

		result := list.EncodeClient(startRank, m.miscData.RankCountPerPage)
		result.Self = proto.Self

		logrus.Debugf("排行榜返回: \n%+v", result)

		hc.Send(rank.NewS2cRequestRankMsg(must.Marshal(result)))
	})
}
