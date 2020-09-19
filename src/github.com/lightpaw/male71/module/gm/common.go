package gm

import (
	"context"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/build"
	"github.com/lightpaw/male7/config/fishing_data"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/chat"
	"github.com/lightpaw/male7/gen/pb/dungeon"
	"github.com/lightpaw/male7/gen/pb/login"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/random_event"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/module/rank/ranklist"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"math/rand"
	"time"
)

func (m *GmModule) newCommonGmGroup() *gm_group {
	return &gm_group{
		tab: "常用",
		handler: []*gm_handler{
			newStringHandler("功能开启", "", m.openFunction),
			newHeroStringHandler("每日重置", "", m.resetDaily),
			newHeroIntHandler("服务器踢人啦", "", m.disconnect),
			newHeroStringHandler("加物品、装备、宝石、元宝(1m)、资源(10m)", "", func(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
				m.addYuanbao(1000000, hero, result, hc)
				m.addDianquan(1000000, hero, result, hc)
				m.addYinliang(1000000, hero, result, hc)
				m.addResources(10000000, hero, result, hc)
				m.addGoods(input, hero, result, hc)
				m.addGem(input, hero, result, hc)
				m.addEquip(input, hero, result, hc)
			}),
			newHeroStringHandler("清物品、装备、元宝、资源", "", func(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
				m.addYuanbao(-1000000000, hero, result, hc)
				m.addDianquan(-1000000000, hero, result, hc)
				m.addYinliang(-1000000000, hero, result, hc)
				m.clearResources(0, hero, result, hc)
				m.clearGoods(input, hero, result, hc)
				m.clearGems(input, hero, result, hc)
				m.clearEquips(input, hero, result, hc)
			}),
			newHeroIntHandler("加君主等级", "1", m.addHeroLevel),
			newIntHandler("加主城等级", "1", m.addHomeLevel),
			newIntHandler("加/减繁荣度(不会改等级)", "1000", m.addProsperity),
			newHeroIntHandler("清空资源", "", m.clearResources),
			newHeroIntHandler("资源全加", "1000000", m.addResources),
			newHeroIntHandler("加元宝(负数表示减)", "100000", m.addYuanbao),
			newHeroIntHandler("加点券(负数表示减)", "100000", m.addDianquan),
			newHeroIntHandler("给人加银两(负数表示减)", "100000", m.addYinliang),
			newHeroIntHandler("加武将等级(不含外出的)", "100", m.addCaptainLevel),
			newIntHandler("群嘲", "0", m.chaoFeng),
			newIntHandler("护驾", "0", m.huJia),
			newHeroIntHandler("千重楼爬塔", "1", m.towerUp),
			newHeroStringHandler("飘字", "", func(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
				result.AddFunc(func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(m.datas.TextHelp().RealmInvadeSuccess.New().
						WithTroopIndex(1).
						JsonString())
				})
			}),
			newStringHandler("行军加速", "", func(input string, hc iface.HeroController) {
				// 找到第一只队伍，加速
				var troopId int64
				hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
					for _, v := range hero.Troops() {
						ii := v.GetInvateInfo()
						if ii != nil && ii.State().IsMoving() {
							troopId = ii.Id()
							break
						}
					}

					return false
				})

				r := m.realmService.GetBigMap()
				if r != nil {
					r.SpeedUp(hc, troopId, 0, 0.5, 0)
				}

			}),
			newHeroStringHandler("发小礼包", "", func(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {

				if hero.GuildId() == 0 {
					result.AddFunc(func() pbutil.Buffer {
						return misc.NewS2cScreenShowWordsMsg("你没有联盟，发毛小礼包啊")
					})
					return
				}

				data := m.datas.GuildEventPrizeData().Array[rand.Intn(len(m.datas.GuildEventPrizeData().Array))]
				m.modules.GuildModule().GmGiveGuildEventPrize(hero, result, []*guild_data.GuildEventPrizeData{data})
			}),
			newHeroStringHandler("发很多小礼包", "", func(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
				if hero.GuildId() == 0 {
					result.AddFunc(func() pbutil.Buffer {
						return misc.NewS2cScreenShowWordsMsg("你没有联盟，发毛小礼包啊")
					})
					return
				}

				var datas []*guild_data.GuildEventPrizeData
				for i := 0; i < 100; i++ {
					data := m.datas.GuildEventPrizeData().Array[rand.Intn(len(m.datas.GuildEventPrizeData().Array))]
					datas = append(datas, data)
				}

				m.modules.GuildModule().GmGiveGuildEventPrize(hero, result, datas)
			}),
			newIntHandler("联盟宝箱充能", "0", func(amount int64, hc iface.HeroController) {
				guildId, ok := hc.LockGetGuildId()
				if !ok || guildId == 0 {
					hc.Send(misc.NewS2cScreenShowWordsMsg("你没有联盟，发毛小礼包啊"))
					return
				}

				m.modules.GuildModule().GmAddBigBoxEnergy(guildId, u64.FromInt64(amount))
			}),
			newStringHandler("怪物出征", "", func(input string, hc iface.HeroController) {

				hc.FuncNotError(func(hero *entity.Hero) (heroChanged bool) {
					for _, v := range m.datas.GetRegionMultiLevelNpcTypeDataArray() {
						hero.GetOrCreateNpcTypeInfo(v.Type).SetHate(v.FightHate)
					}
					return
				})
			}),
			newStringHandler("加速打我", "", func(input string, hc iface.HeroController) {
				m.realmService.GetBigMap().GmSpeedUpFightMe(hc.Id())
			}),
			newHeroStringHandler("加推图次数", "", func(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {

				ctime := m.time.CurrentTime()

				heroDungeon := hero.Dungeon()
				heroDungeon.ChallengeTimes().SetTimes(heroDungeon.ChallengeTimes().MaxTimes(), ctime)
				result.Add(dungeon.NewS2cUpdateChallengeTimesMsg(heroDungeon.ChallengeTimes().StartTimeUnix32()))
			}),
			newIntHandler("加联盟日志", "10", func(input int64, hc iface.HeroController) {
				guildId, _ := hc.LockGetGuildId()
				if guildId != 0 {
					seconds := timeutil.Marshal32(m.time.CurrentTime())

					hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
						for k := range shared_proto.GuildLogType_name {
							if k == 0 {
								continue
							}

							t := shared_proto.GuildLogType(k)
							for i := 0; i < int(input); i++ {
								seconds += 60

								proto := &shared_proto.GuildLogProto{
									Icon: hero.Head(),
									Time: seconds,
									Text: "GM生成的文字",
									Type: t,
								}

								switch rand.Int31n(3) {
								case 0:
									proto.HeroId = hero.IdBytes()
								case 1:
									proto.FightX = int32(hero.BaseX())
									proto.FightY = int32(hero.BaseY())
								}

								m.sharedGuildService.AddLog(guildId, proto)
							}
						}
					})
				}
			}),
			newIntHandler("联盟升级", "1", m.upgradeGuildLevel),
			newIntHandler("联盟加银两", "1", m.addGuildYinliang),
			newIntHandler("联盟加建设值", "10000", func(input int64, hc iface.HeroController) {
				m.addGuildBuildAmount(input, hc)
			}),
			newIntHandler("联盟秒科技cd", "1", m.miaoGuildTechCd),
			newIntHandler("联盟加声望", "100", m.addPrestige),
			newHeroStringHandler("联盟开启弹劾", "", func(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
				m.openImpeachLeader(hero)
			}),
			newStringHandler("清除NPC联盟", "", func(input string, hc iface.HeroController) {
				m.modules.GuildModule().GmRemoveNpcGuild()
			}),
			newHeroIntHandler("恢复讨伐次数", "1", func(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
				ctime := m.time.CurrentTime()
				multiLevelNpcTimes := hero.GetMultiLevelNpcTimes()
				multiLevelNpcTimes.AddTimes(u64.FromInt64(amount), ctime, 0)
				result.Add(region.NewS2cUpdateMultiLevelNpcTimesMsg(multiLevelNpcTimes.StartTimeUnix32(), nil))
			}),
			newStringHandler("热更新版本", build.GetClientVersion(), func(input string, hc iface.HeroController) {
				m.clusterService.GMUpdateClientVersion(input)
				hc.Send(misc.NewS2cScreenShowWordsMsg("热更新版本号修改为 " + input))
			}),
			newStringHandler("服务器推送测试", "", func(input string, hc iface.HeroController) {
				m.pushService.GmPush(hc.Id())
			}),
			newIntHandler("钓鱼概率测试", "10000", func(input int64, hc iface.HeroController) {
				m.modules.FishingModule().GmFishingRate(hc.Id(), fishing_data.FishTypeYuanbao, u64.Min(uint64(input), 100000))
			}),
			newIntHandler("驿站(东施)概率测试", "10000", func(input int64, hc iface.HeroController) {
				m.modules.MilitaryModule().GmRate(hc.Id(), 1, u64.Min(uint64(input), 100000))
			}),
			newIntHandler("驿站(韩非)概率测试", "10000", func(input int64, hc iface.HeroController) {
				m.modules.MilitaryModule().GmRate(hc.Id(), 2, u64.Min(uint64(input), 100000))
			}),
			newIntHandler("驿站(西施)概率测试", "10000", func(input int64, hc iface.HeroController) {
				m.modules.MilitaryModule().GmRate(hc.Id(), 3, u64.Min(uint64(input), 100000))
			}),
			newIntHandler("驿站(白起)概率测试", "10000", func(input int64, hc iface.HeroController) {
				m.modules.MilitaryModule().GmRate(hc.Id(), 4, u64.Min(uint64(input), 100000))
			}),
			newHeroIntHandler("禁言(秒)", "60", func(input int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
				hero.MiscData().SetBanChatEndTime(m.time.CurrentTime().Add(time.Second * time.Duration(input)))
				hc.Send(chat.NewS2cBanChatMsg(timeutil.Marshal32(hero.MiscData().GetBanChatEndTime())))
			}),
			newHeroStringHandler("禁言解除", "", func(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
				hero.MiscData().SetBanChatEndTime(m.time.CurrentTime())
				hc.Send(chat.NewS2cBanChatMsg(timeutil.Marshal32(hero.MiscData().GetBanChatEndTime())))
			}),
			newIntHandler("封号(秒)", "60", func(input int64, hc iface.HeroController) {
				hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
					addTime := time.Second * time.Duration(input)
					hero.MiscData().SetBanLoginEndTime(m.time.CurrentTime().Add(addTime))
					result.Changed()
					result.Add(login.NewS2cBanLoginMsg(timeutil.DurationMarshal32(addTime)))
					result.Ok()
				})
				hc.Disconnect(misc.ErrDisconectReasonFailGm) //@Albert 封号
			}),
			newHeroIntHandler("创建新号(个数)", "1", func(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
				count := u64.Min(uint64(amount), 100)
				id := int64(u64.Min(uint64(rand.Uint32()), 9223372036854775807-0xFFFF))
				for i := uint64(0); i < count; i++ {
					name := fmt.Sprintf("君主_%d", id)
					exist := false
					err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
						exist, err = m.db.HeroNameExist(ctx, name)
						return
					})
					if err != nil {
						logrus.Error(err.Error())
						return
					}
					if exist {
						count++
						id++
						continue
					}
					countryId := rand.Uint64()%7 + 1
					ctime := m.time.CurrentTime()
					hero := entity.NewHero(id, name, m.datas.HeroInitData(), ctime)
					hero.SetMale(true)
					hero.SetCountryId(countryId)
					// 初始化数据
					data := m.datas.HeroCreateData()
					heroMilitary := hero.Military()
					heroMilitary.SetNewSoldierCount(data.NewSoldier, ctime)
					hero.SetDailyResetTime(ctime)
					hero.SetDailyZeroResetTime(ctime)
					hero.SetSeasonResetTime(ctime)
					heroMilitary.AddFreeSoldier(data.NewSoldier, ctime)
					hero.GetSafeResource().AddGold(data.Gold)
					hero.GetSafeResource().AddFood(data.Food)
					hero.GetSafeResource().AddWood(data.Wood)
					hero.GetSafeResource().AddStone(data.Stone)
					// 加初始仇恨
					for _, t := range m.datas.GetRegionMultiLevelNpcTypeDataArray() {
						if t.InitHate > 0 {
							hero.GetOrCreateNpcTypeInfo(t.Type).SetHate(t.InitHate)
						}
					}
					// 初始入侵野怪
					for _, v := range m.datas.GetHomeNpcBaseDataArray() {
						if v.HomeBaseLevel == 1 || v.BaYeStage == 1 {
							hero.CreateHomeNpcBase(v)
						}
					}
					// 加武将
					for _, data := range m.datas.HeroCreateData().Captain {
						captain := hero.NewCaptain(data, ctime)
						captain.CalculateProperties()
						captain.AddSoldier(u64.Sub(captain.SoldierCapcity(), captain.Soldier()))

						hero.Military().AddCaptain(captain)

						hero.WalkPveTroop(func(troop *entity.PveTroop) (endWalk bool) {
							troop.AddCaptain(captain)
							return
						})

						troop, index := hero.GetRecruitCaptainTroop()
						if troop != nil {
							troop.SetCaptainIfAbsent(index, captain, 0)
						}
					}
					// 修炼馆开始时间
					hero.Military().SetGlobalTrainStartTime(ctime)
					for _, t := range hero.Troops() {
						t.UpdateFightAmountIfChanged()
					}
					// 轩辕会武初始积分
					hero.Xuanyuan().AddScore(m.datas.XuanyuanMiscData().InitScore)
					if m.config.GetSkipHeader() {
						for k := range shared_proto.HeroBoolType_name {
							bt := shared_proto.HeroBoolType(k)
							switch bt {
							case shared_proto.HeroBoolType_BOOL_XUAN_YUAN:
							default:
								hero.Bools().SetTrue(bt)
							}
						}
					}
					// 功能开启
					for _, data := range hero.Function().FunctionOpenDataArray {
						if heromodule.GetIsFuncOpened(hero, data) {
							hero.Function().OpenFunction(data.FunctionType)
						}
					}
					// 任务
					hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumLoginDays)
					taskTypes := []shared_proto.TaskTargetType{
						shared_proto.TaskTargetType_TASK_TARGET_ACCUM_LOGIN_DAY,
					}
					if len(taskTypes) > 0 {
						hero.TaskList().WalkAllTask(func(task entity.Task) (endedWalk bool) {
							for _, t := range taskTypes {
								task.Progress().UpdateTaskTypeProgress(t, hero)
							}
							return false
						})
					}
					// 初始不自动补兵
					hero.MiscData().SetDefenserDontAutoFullSoldier(true)
					// 创建
					if e := m.heroDataService.Create(hero); e != nil {
						logrus.Error(err.Error())
						return
					}
					// 创建成功
					m.gameExporter.GetRegisterCounter().Inc()
					// 创建主城
					realm, baseX, baseY := m.realmService.ReserveNewHeroHomePos(hero.CountryId())
					//logrus.WithField("x", baseX).WithField("y", baseY).Debugf("初始化玩家主城")
					processed, err := realm.AddBase(hc.Id(), baseX, baseY, realmface.AddBaseHomeNewHero)
					if !processed || err != nil {
						logrus.WithError(err).WithField("processed", processed).Error("创建英雄时, realm.AddBase失败")
					}
					id++
				}
			}),
			newIntHandler("充值", "100", m.charge),
			newStringHandler("随机事件填满老家区域", "", func(input string, hc iface.HeroController) {
				hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
					ctime := m.time.CurrentTime()
					e := hero.RandomEvent()
					e.SetBigRefreshTime(ctime.Add(m.datas.MiscGenConfig().RandomEventBigRefreshDuration))
					e.SetSmallRefreshTime(ctime.Add(m.datas.MiscGenConfig().RandomEventSmallRefreshDuration))
					e.ClearAllEvents()
					cubes := m.datas.RandomEventPositionDictionary().GetRandomPositions(m.datas.MiscGenConfig().RandomEventNum)
					e.PutEvents(cubes)
					var arrX, arrY []int32
					for _, c := range cubes {
						x, y := c.XYI32()
						arrX = append(arrX, x)
						arrY = append(arrY, y)
					}
					if cubes = m.datas.RandomEventPositionDictionary().CheckAndCatch4Block(cubes, len(m.datas.GetEventPositionArray()), hero.BaseX(), hero.BaseY()); cubes != nil {
						e.PutEvents(cubes)
						for _, c := range cubes {
							x, y := c.XYI32()
							arrX = append(arrX, x)
							arrY = append(arrY, y)
						}
					}
					result.Add(random_event.NewS2cNewEventMsg(arrX, arrY))
				})
			}),
		},
	}
}

func (m *GmModule) charge(amount int64, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.Misc().AddChargeAmount(u64.FromInt64(amount), m.time.CurrentTime())
		m.addYuanbao(amount, hero, result, hc)
		result.Add(misc.NewS2cUpdateChargeAmountMsg(u64.Int32(hero.Misc().ChargeAmount())))
	})
}

func (m *GmModule) addGuildYinliang(amount int64, hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	m.modules.GuildModule().GmAddGuildYinliang(amount, guildId)
}

func (m *GmModule) upgradeGuildLevel(amount int64, hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	m.modules.GuildModule().GmUpgradeGuildLevel(guildId)
}

func (m *GmModule) addGuildBuildAmount(amount int64, hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	m.modules.GuildModule().GmAddGuildBuildAmount(guildId, uint64(amount))
}

func (m *GmModule) miaoGuildTechCd(amount int64, hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	m.modules.GuildModule().GmMiaoGuildTechCd(guildId)
}

func (m *GmModule) addPrestige(amount int64, hc iface.HeroController) {
	guildId, _ := hc.LockGetGuildId()
	if guildId <= 0 {
		return
	}
	m.sharedGuildService.FuncGuild(guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			return
		}

		var newAmount uint64
		if amount < 0 {
			newAmount = u64.Sub(g.GetPrestige(), u64.FromInt64(i64.Abs(amount)))
		} else {
			newAmount = g.GetPrestige() + u64.FromInt64(amount)
		}
		g.SetPrestige(newAmount)

		m.modules.RankModule().AddOrUpdateRankObj(ranklist.NewGuildRankObj(
			m.sharedGuildService.GetSnapshot, m.heroSnapshotService.Get, g))
	})
}

func (m *GmModule) openImpeachLeader(hero *entity.Hero) {
	gid := hero.GuildId()
	if gid > 0 {
		m.modules.GuildModule().GmOpenImpeachLeader(gid)
	}
}

func (m *GmModule) disconnect(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	hc.Disconnect(misc.ErrDisconectReasonFailGm)
}

func (m *GmModule) openFunction(input string, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		// 君主等级到16级
		m.addHeroLevelTo(16, hero, result, hc)

		// 爬千重楼到20层
		m.towerUpTo(20, hero, result, hc)

		// 官府等级11级
		m.upgradeBuildToLevelX(11, hero, result, hc)
	})

	// 霸业任务到第10章节
	m.completeBaYeTargetAndCollectTo(10, hc)
}
