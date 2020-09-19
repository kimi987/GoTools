package realm

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/event"
	"github.com/lightpaw/male7/util/lock"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/module/realm/realmerr"
	"github.com/lightpaw/male7/util/compress"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
)

//gogen:iface
type RealmService struct {
	services    *services
	tickService iface.TickerService

	bigMap *Realm
}

type services struct {
	dep                 iface.ServiceDep
	datas               iface.ConfigDatas
	dbService           iface.DbService
	seasonService       iface.SeasonService
	timeService         iface.TimeService
	world               iface.WorldService
	heroDataService     iface.HeroDataService
	fightModule         iface.FightXService
	mail                iface.MailModule
	heroSnapshotService iface.HeroSnapshotService
	guildService        iface.GuildService
	reminderService     iface.ReminderService
	farmService         iface.FarmService
	xiongNuService      iface.XiongNuService
	extraTimesService   iface.ExtraTimesService
	pushService         iface.PushService
	tickService         iface.TickerService
	tlogService         iface.TlogService
	buffService         iface.BuffService
	baiZhanService      iface.BaiZhanService

	otherRealmEventQueue *event.FuncQueue

	homeNpcBaseIdGen npcid.NpcIdGen
}

func NewRealmService(dep iface.ServiceDep, seasonService iface.SeasonService, tlogService iface.TlogService,
	fightModule iface.FightXService, mail iface.MailModule, extraTimesService iface.ExtraTimesService,
	dbService iface.DbService, tickService iface.TickerService, pushService iface.PushService,
	reminderService iface.ReminderService, xiongNuService iface.XiongNuService, farmService iface.FarmService,
	buffService iface.BuffService, baiZhanService iface.BaiZhanService) *RealmService {

	otherRealmEventQueue := event.NewFuncQueue(1024, "RealmService.OtherRealmEventQueue")
	result := &RealmService{}

	// 获取各个场景的最大个数
	var protoBytes []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		protoBytes, err = dbService.LoadKey(ctx, server_proto.Key_Region)
		return
	})
	if err != nil {
		logrus.WithError(err).Panicf("加载场景数据失败，dbService.LoadKey(server_proto.Key_Region)")
	}

	proto := &server_proto.RegionModuleProto{}
	if len(protoBytes) > 0 {
		if newBytes, err := compress.GzipUncompress(protoBytes); err == nil && len(newBytes) > 0 {
			protoBytes = newBytes
		}

		err = proto.Unmarshal(protoBytes)
		if err != nil {
			logrus.WithError(err).Panicf("加载场景数据失败，proto.Unmarshal失败")
		}
	}

	service := &services{
		dep:                  dep,
		datas:                dep.Datas(),
		dbService:            dbService,
		timeService:          dep.Time(),
		world:                dep.World(),
		heroDataService:      dep.HeroData(),
		fightModule:          fightModule,
		mail:                 mail,
		heroSnapshotService:  dep.HeroSnapshot(),
		guildService:         dep.Guild(),
		reminderService:      reminderService,
		otherRealmEventQueue: otherRealmEventQueue,
		xiongNuService:       xiongNuService,
		extraTimesService:    extraTimesService,
		homeNpcBaseIdGen:     npcid.NewNpcIdGen(proto.HomeNpcBaseSequence, npcid.NpcType_HomeNpc),
		farmService:          farmService,
		seasonService:        seasonService,
		pushService:          pushService,
		tickService:          tickService,
		tlogService:          tlogService,
		buffService:          buffService,
		baiZhanService:       baiZhanService,
	}

	// 大地图地区
	levelData := dep.Datas().RegionData().MinKeyData // 暂时只有一种地形
	if levelData.RegionType != shared_proto.RegionType_HOME {
		logrus.WithError(err).Panic("地区等级数据配置的第一条数据必须是主城地区数据")
	}

	realmId := realmface.GetRealmId(1, 0, shared_proto.RegionType_HOME)
	result.bigMap = newRealm(realmId, service, levelData) // 暂时只有一种地形

	radius := u64.Min(u64.Max(result.bigMap.GetRadius(), proto.Radius), levelData.Block.MaxRadius())
	result.bigMap.blockManager.radius.Store(radius)

	// 初始化Npc
	result.bigMap.addMultiLevelMonsterIfAbsent(0, radius)
	result.bigMap.lastRefreshBaoZangNpcTime = timeutil.Unix64(proto.BaoZangNpcRefreshTime)
	result.bigMap.lastRefreshJunTuanNpcTime = timeutil.Unix64(proto.JunTuanNpcRefreshTime)

	result.bigMap.astDefendLogs = loadAstDefendLog(dbService)

	result.services = service
	result.tickService = tickService

	for _, ruinsProto := range proto.RealmRuins {
		for _, info := range ruinsProto.Infos {
			result.bigMap.ruinsBasePosInfoMap.AddRuinsBaseProto(result.bigMap, info, timeutil.Unix64(info.Time))
		}
	}

	tryAddPos := func(realm *Realm, baseX, baseY int, name string) bool {
		if !realm.GetMapData().IsValidHomePosition(baseX, baseY) {
			logrus.Errorf("开服发现%s坐标无效，删掉", name)
			return false
		}

		if !realm.IsPosOpened(baseX, baseY) {
			logrus.Errorf("开服发现%s坐标尚未开放，删掉", name)
			return false
		}

		if realm.IsEdgeNotHomePos(baseX, baseY) {
			logrus.Errorf("开服发现%s坐标在边缘，删掉", name)
			return false
		}

		if realm.conflict.addBaseIfCanAdd(baseX, baseY) {
			return true
		} else {
			logrus.Errorf("开服发现%s坐标距离其他主城太近，删掉", name)
			return false
		}
	}

	// 联盟工坊
	dep.Guild().Func(func(guilds sharedguilddata.Guilds) {
		guilds.Walk(func(g *sharedguilddata.Guild) {
			w := g.GetWorkshop()
			if w == nil {
				return
			}

			guildId := g.Id()
			if w.GetProsperity() <= 0 {
				if w.IsComplete() {
					logrus.WithField("guildId", guildId).Error("联盟工坊的繁荣度 == 0")
					g.SetWorkshop(nil)
					return
				}

				w.SetProsperity(dep.Datas().GuildGenConfig().WorkshopProsperityCapcity)
			}

			r := result.getRealm()

			baseX, baseY := w.GetBaseXY()
			if !tryAddPos(r, baseX, baseY, "联盟工坊") {
				g.SetWorkshop(nil)
				return
			}

			baseId := npcid.NewGuildWorkshopId(guildId)
			startTime, endTime := w.GetTime()

			base := r.newGuildWorkshopBase(baseId, dep.Datas().GuildGenConfig().WorkshopBase, guildId, baseX, baseY, startTime, endTime, w.IsComplete())
			base.ReduceProsperityDontKill(u64.Sub(base.Prosperity(), w.GetProsperity()))

			r.addBaseToMap(base)
		})
	})

	ctime := dep.Time().CurrentTime()
	ctimeUnix := ctime.Unix()
	for _, npcBaseProto := range proto.NpcBaseList {
		if npcid.IsXiongNuNpcId(npcBaseProto.BaseId) {

			info := xiongNuService.GetRInfo(npcBaseProto.GuildId)
			if info == nil {
				logrus.WithField("guild", npcBaseProto.GuildId).Errorf("匈奴主城没找到匈奴战斗数据")
				continue
			}

			r := result.getRealm()

			baseX, baseY := int(npcBaseProto.BaseX), int(npcBaseProto.BaseY)
			if !tryAddPos(r, baseX, baseY, "匈奴主城") {
				continue
			}

			base := r.newXiongNuNpcBase(npcBaseProto.BaseId, info, baseX, baseY)
			r.addBaseToMap(base)

			if b := GetXiongNuBase(base); b != nil {
				b.defenser.unmarshalCaptainSoldier(npcBaseProto.CaptainIndex, npcBaseProto.CaptainSoldier)
				b.defenser.onChanged()

				for _, troopProto := range npcBaseProto.DefendingTroop {
					mmd := dep.Datas().GetMonsterMasterData(troopProto.DataId)
					if mmd == nil {
						continue
					}

					troop := r.newDefendingNpcTroop(troopProto.Id, base, base.BaseLevel(), base, base.BaseLevel(), mmd)
					troop.unmarshalCaptainSoldier(troopProto.CaptainIndex, troopProto.CaptainSoldier)
					troop.onChanged()

					// 构建成功
					base.selfTroops[troop.Id()] = troop
					base.targetingTroops[troop.Id()] = troop
				}
			}
			continue
		}

		// 宝藏Npc
		if npcid.IsBaoZangNpcId(npcBaseProto.BaseId) {
			dataId := npcid.GetNpcDataId(npcBaseProto.BaseId)
			if baozNpcData := service.datas.GetBaozNpcData(dataId); baozNpcData != nil {

				if npcBaseProto.HeroId == 0 {
					if index := npcid.GetBaoZangIndex(npcBaseProto.BaseId); index >= baozNpcData.KeepCount {
						continue
					}
				}

				if npcBaseProto.HeroType == HeroTypeCreater && npcBaseProto.HeroEndTime > 0 && npcBaseProto.HeroEndTime < ctimeUnix {
					// 已过期
					continue
				}

				r := result.getRealm()

				baseX, baseY := int(npcBaseProto.GetBaseX()), int(npcBaseProto.GetBaseY())
				if npcid.GetBaoZangBlock(npcBaseProto.BaseId) != regdata.BlockSequence(r.mapData.MustBlockByPos(baseX, baseY)) {
					logrus.Warn("宝藏Npc对应的Block错误，删除宝藏Npc怪物")
					continue
				}

				if !tryAddPos(r, baseX, baseY, "宝藏Npc") {
					continue
				}

				var createrId, createEndTime int64
				if npcBaseProto.HeroType == HeroTypeCreater {
					createrId = npcBaseProto.HeroId
					createEndTime = npcBaseProto.HeroEndTime
				}

				base := r.newBaozNpcBase(npcBaseProto.BaseId, baozNpcData, baseX, baseY,
					createrId, int32(createEndTime))
				if npcBaseProto.Prosperity > 0 {
					toReduce := u64.Sub(base.Prosperity(), npcBaseProto.Prosperity)
					if toReduce > 0 {
						base.ReduceProsperity(toReduce)
					}
				}

				if b := GetBaoZangBase(base); b != nil && b.defenser != nil {
					var totalSoldier int32
					for _, c := range b.defenser.Captains() {
						var soldier int32
						for i, idx := range npcBaseProto.CaptainIndex {
							if int(idx) == c.Index() {
								soldier = npcBaseProto.CaptainSoldier[i]
								break
							}
						}

						c.Proto().Soldier = i32.Min(soldier, c.Proto().TotalSoldier)
						c.Proto().FightAmount = data.ProtoFightAmount(c.Proto().TotalStat, c.Proto().Soldier, c.Proto().SpellFightAmountCoef)

						totalSoldier += soldier
					}

					b.defenser.onChanged()

					if npcBaseProto.HeroType == HeroTypeKiller {
						b.heroId = npcBaseProto.HeroId
						b.heroEndTime = int32(npcBaseProto.HeroEndTime)
					}

					// 更新一下士兵数，战力
					b.updateRoBase(base.getRoBase())
				}

				r.addBaseToMap(base)
			}

			continue
		}

		// 军团怪
		if npcid.IsJunTuanNpcId(npcBaseProto.BaseId) {
			dataId := npcid.GetNpcDataId(npcBaseProto.BaseId)
			if npcData := service.datas.GetJunTuanNpcData(dataId); npcData != nil {
				r := result.getRealm()

				baseX, baseY := int(npcBaseProto.GetBaseX()), int(npcBaseProto.GetBaseY())
				if npcid.GetJunTuanBlock(npcBaseProto.BaseId) != regdata.BlockSequence(r.mapData.MustBlockByPos(baseX, baseY)) {
					logrus.Warn("军团Npc对应的Block错误，删除Npc怪物")
					continue
				}

				if !tryAddPos(r, baseX, baseY, "军团Npc") {
					continue
				}

				base := r.newJunTuanNpcBase(npcBaseProto.BaseId, npcData, baseX, baseY)

				if b := GetJunTuanBase(base); b != nil && len(npcBaseProto.DefendingTroop) > 0 {
					for _, troopProto := range npcBaseProto.DefendingTroop {
						troop := r.newDefendingNpcTroop(troopProto.Id, base, base.BaseLevel(), base, base.BaseLevel(), b.data.Npc.Npc)
						troop.unmarshalCaptainSoldier(troopProto.CaptainIndex, troopProto.CaptainSoldier)
						troop.onChanged()

						// 构建成功
						base.selfTroops[troop.Id()] = troop
						base.targetingTroops[troop.Id()] = troop
					}
				}
				base.updateRoBase()
				r.addBaseToMap(base)
			}

			continue
		}
	}

	// 只将存在主城的数据Load上来

	var regionHeros []*entity.Hero
	err = ctxfunc.Timeout1m(func(ctx context.Context) (err error) {
		regionHeros, err = service.dbService.LoadAllHeroData(ctx)
		return
	})
	if err != nil {
		logrus.WithError(err).Panic("RealmService db.LoadAllRegionHero出错")
	}
	logrus.WithField("count", len(regionHeros)).Info("加载场景英雄数据")

	//if len(regionHeros) <= 0 {
	//	// 一个人都没有，可能是之前数据没更新，全部load上来一次
	//	var allHeros []*entity.Hero
	//	err = ctxfunc.Timeout1m(func(ctx context.Context) (err error) {
	//		allHeros, err = service.dbService.LoadAllHeroData(ctx)
	//		return
	//	})
	//	if err != nil {
	//		logrus.WithError(err).Panic("RealmService db.LoadedAllHeroData出错")
	//	}
	//
	//	regionHeros = allHeros
	//}

	heroIds := make([]int64, 0, len(regionHeros))
	for _, hero := range regionHeros {
		if hero.BaseLevel() <= 0 && hero.Prosperity() <= 0 {
			continue
		}

		heroIds = append(heroIds, hero.Id())
		err := dep.HeroData().Put(hero)
		if err != nil {
			if err != lock.ErrCreateExist {
				logrus.WithError(err).Panicf("初始化野外场景主城行营，put英雄失败")
			}
		}
	}
	regionHeros = nil

	deleteCount := 0
	for _, heroId := range heroIds {
		dep.HeroData().Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {

			if err != nil {
				logrus.WithError(err).Panicf("初始化野外场景主城行营，lock英雄失败")
				return
			}

			if hero.Prosperity() <= 0 {
				// 主城存在，但是繁荣度空，给一点繁荣度
				hero.SetProsperity(u64.Max(hero.ProsperityCapcity()/3, 1))
				heroChanged = true

				logrus.Errorf("开服发现玩家主城繁荣度为0，设置繁荣度到1/3")
			}

			if hero.BaseLevel() <= 0 {
				hero.SetBaseLevel(1)
				heroChanged = true

				logrus.Errorf("开服发现玩家主城等级为0，设置到1级")
			}

			realm := result.getRealm()
			if realm == nil || realm.levelData.RegionType != shared_proto.RegionType_HOME {
				logrus.Error("开服发现玩家主城地图找不到，删掉主城和行营",
					realmface.ParseRegionType(hero.BaseRegion()),
					realmface.ParseRegionSequence(hero.BaseRegion()),
					realm == nil,
					hero.BaseRegion())
				//logrus.Errorf("开服发现玩家主城地图找不到，删掉主城和行营")
				hero.ClearBase()

				heroChanged = true
				deleteCount++
				return
			}

			if !realm.GetMapData().IsValidHomePosition(hero.BaseX(), hero.BaseY()) {
				if hero.BaseRegion() != 0 {
					logrus.Errorf("开服发现玩家主城坐标无效，删掉主城和行营")
				}
				hero.ClearBase()

				heroChanged = true
				deleteCount++
				return
			}

			if hero.BaseRegion() == 0 {
				hero.SetBaseXY(realm.Id(), hero.BaseX(), hero.BaseY())
				heroChanged = true

				logrus.WithField("id", realm.Id()).Errorf("开服发现玩家主城地区ID为0")
			}

			if !realm.IsPosOpened(hero.BaseX(), hero.BaseY()) {
				logrus.Errorf("开服发现玩家主城坐标尚未开放，删掉主城和行营")
				hero.ClearBase()

				heroChanged = true
				deleteCount++
				return
			}

			if realm.IsEdgeNotHomePos(hero.BaseX(), hero.BaseY()) {
				logrus.Errorf("开服发现玩家主城坐标在边缘，删掉主城和行营")
				hero.ClearBase()

				heroChanged = true
				deleteCount++
				return
			}

			base := realm.newBase(hero)
			if realm.conflict.addBaseIfCanAdd(base.BaseX(), base.BaseY()) {
				realm.addBaseToMap(base)
			} else {
				logrus.Errorf("开服发现玩家主城坐标距离其他主城太近，删掉主城和行营")
				hero.ClearBase()

				heroChanged = true
				deleteCount++
				return
			}

			// 添加Npc基地
			if home := GetHomeBase(base); home != nil {
				homeNpcBase := realm.newAllHomeNpcBase(hero, base.BaseX(), base.BaseY())
				for _, v := range homeNpcBase {
					home.homeNpcBase[v.Id()] = v
				}
			}

			return
		})
	}
	logrus.WithField("count", len(heroIds)).WithField("delete", deleteCount).Info("初始化玩家主城")

	// 处理宝藏怪物击杀者数据
	ctimeUnix32 := int32(ctimeUnix)
	result.getRealm().rangeBases(func(base *baseWithData) (toContinue bool) {
		if baoz := GetBaoZangBase(base); baoz != nil && baoz.heroType == HeroTypeKiller && baoz.heroId != 0 {
			if killer := result.getRealm().getBase(baoz.heroId); killer != nil {
				if home := GetHomeBase(killer); home != nil {
					if hero := dep.HeroSnapshot().Get(home.Id()); hero != nil {
						baoz.heroBytes = hero.EncodeBasic4ClientBytes()

						if ctimeUnix32 < baoz.heroEndTime {
							home.AddKeepBaozMap(baoz.Id(), baoz.heroEndTime)
						}
						return true
					}
				}
			}

			baoz.heroId = 0
			baoz.heroEndTime = 0
		}

		return true
	})

	// 刷新宝藏怪物
	result.getRealm().GmRefreshBaoZangNpc()

	// Npc部队
	for _, troopProto := range proto.TroopList {
		err := result.getRealm().rebuildTroop(troopProto, ctime)
		if err != nil {
			logrus.WithField("realmid", troopProto.RealmId).WithField("troop", troopProto).WithError(err).Error("开服时, 重构Npc troop出错")
			continue
		}
	}

	// 读取玩家部队
	var joinAssemblyHeroIds []int64
	for _, heroId := range heroIds {
		dep.HeroData().Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {

			if err != nil {
				logrus.WithError(err).Panicf("初始化野外场景部队，lock英雄失败")
				return
			}

			if hero.BaseRegion() <= 0 {
				// 不该有部队
				return
			}

			hasJoinAssemblyTroop := false
			hero.IterTroop(func(troop *entity.TroopInvaseInfo) (troopValid bool) {

				if troop.AssemblyId() == 0 || troop.AssemblyId() == troop.Id() {
					if err := result.getRealm().rebuildHeroTroop(troop, hero, ctime); err != nil {
						logrus.WithField("realmid", troop.RegionID()).WithField("heroid", hero.Id()).WithField("troopid", troop.Id()).WithError(err).Error("开服时, 重构troop出错")
						return false
					}
				} else {
					// 这里先不处理加入集结的队伍
					hasJoinAssemblyTroop = true
				}
				return true
			})

			if hasJoinAssemblyTroop {
				joinAssemblyHeroIds = append(joinAssemblyHeroIds, hero.Id())
			}

			return false
		})
	}

	for _, heroId := range joinAssemblyHeroIds {
		dep.HeroData().Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {

			if err != nil {
				logrus.WithError(err).Panicf("初始化野外场景部队，lock英雄失败")
				return
			}

			hero.IterTroop(func(troop *entity.TroopInvaseInfo) (troopValid bool) {

				if troop.AssemblyId() != 0 && troop.AssemblyId() != troop.Id() {
					// 加入别人的集结
					if err := result.getRealm().rebuildHeroTroop(troop, hero, ctime); err != nil {
						logrus.WithField("realmid", troop.RegionID()).WithField("heroid", hero.Id()).WithField("troopid", troop.Id()).WithError(err).Error("开服时, 重构加入集结troop出错")
						return false
					}
				}
				return true
			})

			return false
		})
	}

	result.getRealm().rangeBases(func(base *baseWithData) (toContinue bool) {
		for _, t := range base.selfTroops {
			if t.assembly != nil && t.assembly.self == t {
				t.assembly.updateAddedStat(result.getRealm(), 0)
			}
		}
		return true
	})

	// 初始化添加资源点占用情况
	result.walkRealmsUnderLock(func(realm *Realm) {
		realm.rangeBases(func(base *baseWithData) (toContinue bool) {
			realm.addHomeResourcePointBlock(base, false)
			return true
		})

	})

	// 所有地图开始loop
	result.walkRealmsUnderLock(func(realm *Realm) {
		go call.CatchLoopPanic(realm.loop, "Realm.loop")
	})

	go call.CatchLoopPanic(result.loop, "RealmService.loop")

	dep.Guild().RegisterCallback(result)
	heromodule.RegisterHeroEventHandler("RealmModule.HeroEventHandler", result.handleHeroEvent)

	return result
}

func loadAstDefendLog(dbService iface.DbService) (proto *server_proto.AllAstDefendLogProto) {
	proto = &server_proto.AllAstDefendLogProto{Logs: make(map[int64]*server_proto.AstDefendLogListProto)}

	var protoBytes []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		protoBytes, err = dbService.LoadKey(ctx, server_proto.Key_AstDefendLog)
		return
	})
	if err != nil {
		logrus.WithError(err).Errorf("加载盟友驻扎日志数据失败，dbService.LoadKey(server_proto.Key_Region)")
		return
	}

	if len(protoBytes) > 0 {
		err = proto.Unmarshal(protoBytes)
		if err != nil {
			logrus.WithError(err).Errorf("加载盟友驻扎数据失败，proto.Unmarshal失败")
			return
		}
	}

	return proto
}

func (s *RealmService) Close() {

	// 如果要关闭场景，需要先关闭service.otherRealmEventQueue
	s.services.otherRealmEventQueue.Close()

	// close 每一个场景
	s.walkRealmsUnderLock(func(realm *Realm) {
		realm.stop()
	})

	ctxfunc.Timeout30s(func(ctx context.Context) (err error) {
		s.save(ctx, true)
		return
	})

}

func (s *RealmService) save(ctx context.Context, isClosed bool) {
	proto := &server_proto.RegionModuleProto{}
	proto.HomeNpcBaseSequence = s.services.homeNpcBaseIdGen.Sequence()
	proto.Radius = s.GetBigMap().GetRadius()

	s.walkRealms(func(realm *Realm) {
		ruinsProto := realm.ruinsBasePosInfoMap.Encode()
		if ruinsProto != nil {
			proto.RealmRuins = append(proto.RealmRuins, ruinsProto)
		}

		if isClosed {
			realm.saveRealmNpcBases(proto)
		} else {
			realm.queueSaveRealmNpcBases(proto)
		}
	})

	if err := s.services.dbService.SaveKey(ctx, server_proto.Key_Region, compress.GzipCompress(must.Marshal(proto))); err != nil {
		logrus.WithError(err).Errorf("RealmService.Save() key:Key_Region")
	}

	if err := s.services.dbService.SaveKey(ctx, server_proto.Key_AstDefendLog, must.Marshal(s.GetBigMap().GetAstDefendLogs())); err != nil {
		logrus.WithError(err).Errorf("RealmService.Save() key:Key_AstDefendLog", )
	}
}

func (s *RealmService) walkRealms(walkFunc func(realm *Realm)) {
	walkFunc(s.bigMap)
}

func (s *RealmService) walkRealmsUnderLock(walkFunc func(realm *Realm)) {
	walkFunc(s.bigMap)
}

func (s *RealmService) loop() {

	tickPerMinuteTime := s.tickService.GetPerMinuteTickTime()
	tickPer10MinuteTime := s.tickService.GetPer10MinuteTickTime()
	for {

		select {
		case <-tickPerMinuteTime.Tick():
			tickPerMinuteTime = s.tickService.GetPerMinuteTickTime()

		case <-tickPer10MinuteTime.Tick():
			tickPer10MinuteTime = s.tickService.GetPer10MinuteTickTime()

			ctxfunc.Timeout3s(func(ctx context.Context) error {
				s.save(ctx, false)
				return nil
			})
		}

	}
}

func (s *RealmService) GetRealm(id int64) iface.Realm {
	if r := s.getRealm(); id != 0 && r != nil {
		return r
	}
	return nil
}

func (s *RealmService) GetBigMap() iface.Realm {
	return s.bigMap
}

func (s *RealmService) getRealm() *Realm {
	return s.bigMap
}

// 坐标是已经占座占好了的, 直接调用realm.AddBase就可以了
func (s *RealmService) ReserveNewHeroHomePos(country uint64) (r iface.Realm, x, y int) {

	r = s.GetBigMap()
	var ok bool
	ok, x, y = r.ReserveNewHeroHomePos(country)
	if !ok {
		// 找不到，那就随机一个吧
		x, y = r.RandomBasePos()
	}

	return
}

// 坐标是已经占座占好了的, 直接调用realm.AddBase就可以了
func (s *RealmService) ReserveRandomHomePos(t realmface.RandomPointType) (r iface.Realm, x, y int) {

	r = s.GetBigMap()
	var ok bool
	ok, x, y = r.ReserveRandomHomePos(t)
	if !ok {
		// 找不到，那就随机一个吧
		x, y = r.RandomBasePos()
	}

	return
}

func (s *RealmService) StartCareMilitary(hc iface.HeroController) {
	s.bigMap.StartCareMilitary(hc)
}

// callback
func (s *RealmService) OnGuildSnapshotUpdated(origin, update *guildsnapshotdata.GuildSnapshot) {
	if origin == nil {
		// 新建帮派，什么都不干
		return
	}

	if origin.Name == update.Name && origin.FlagName == update.FlagName {
		// 不是改名或者改旗号，什么都不干
		return
	}

	s.onGuildChanged(update.Id, update.Name, update.FlagName, false)
}

func (s *RealmService) OnGuildSnapshotRemoved(id int64) {
	s.onGuildChanged(id, "", "", true)
}

func (s *RealmService) onGuildChanged(id int64, name, flagName string, isRemoved bool) {

	s.services.otherRealmEventQueue.MustFunc(func() {
		s.walkRealms(func(realm *Realm) {
			realm.onGuildChanged(id, name, flagName, isRemoved)
		})
	})
}

func (s *RealmService) handleHeroEvent(hero *entity.Hero, result herolock.LockResult, event shared_proto.HeroEvent) {
	//switch event {
	//case heromodule.JadeOreChanged:
	//	heroId := hero.Id()
	//	homeBaseRegion := hero.BaseRegion()
	//	if homeBaseRegion == 0 {
	//		return
	//	}
	//
	//	newJadeOre := hero.JadeOre()
	//	s.services.otherRealmEventQueue.TryFunc(func() {
	//		if homeBaseRegion != 0 {
	//			if r := s.getRealm(); r != nil {
	//				r.UpdateBaseJadeOre(heroId, newJadeOre)
	//			}
	//		}
	//	})
	//}
}

func (s *RealmService) AddProsperityFunc(heroId, baseRegion int64, toAddProsperity uint64, opName string) iface.Func {
	if toAddProsperity > 0 {
		realm := s.GetBigMap()
		if realm != nil {
			return func() {
				processed, err := realm.AddProsperity(heroId, toAddProsperity)
				if err != nil {
					logrus.WithError(err).Errorf("%s，加繁荣度出错", opName)
				} else if !processed {
					logrus.Errorf("%s，加繁荣度超时", opName)
				}
			}
		}
	}

	return nil
}

func (m *RealmService) CheckCanMoveBase(heroId int64, oldX, oldY int, removeSelfTroop bool) (succ bool) {
	return m.GetBigMap().CheckCanMoveBase(heroId, oldX, oldY, removeSelfTroop) == nil
}

func (m *RealmService) DoMoveBase(moveType shared_proto.GoodsMoveBaseType, newRealm iface.Realm, hc iface.HeroController, originBaseX, originBaseY, newX, newY int, removeSelfTroop bool) (succ bool) {
	// 同一张地图
	processed, err := newRealm.MoveBase(hc, originBaseX, originBaseY, newX, newY, removeSelfTroop)
	if !processed {
		logrus.Debug("迁移基地, 超时")
		hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_SERVER_ERROR)
		return
	}

	switch err {
	case nil:
		// 成功了
		succ = true

		// tlog
		m.services.tlogService.TlogMoveCitylFlowById(hc.Id(), uint64(moveType), int64(originBaseX), int64(originBaseY), int64(newX), int64(newY))

	case realmerr.ErrFastMoveBaseSelfNoBase:
		hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_BASE_NOT_EXIST)
	case realmerr.ErrFastMoveBasePosChanged:
		hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_SERVER_ERROR)
	case realmerr.ErrFastMoveBaseOutside:
		hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_CAPTAIN_OUT_SIDE)
	case realmerr.ErrLockHeroErr:
		hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_SERVER_ERROR)
	default:
		hc.Send(region.ERR_FAST_MOVE_BASE_FAIL_SERVER_ERROR)
		logrus.WithError(err).Errorf("迁移基地, 遇到未知的错误码")
	}

	return
}
