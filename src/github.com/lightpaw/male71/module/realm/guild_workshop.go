package realm

import (
	"github.com/lightpaw/male7/entity/npcid"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/pb/server_proto"
	"fmt"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/util/imath"
)

func (r *Realm) tickGuildWorkshop(ctime time.Time) {

	ctimeUnix32 := timeutil.Marshal32(ctime)

	r.rangeBases(func(base *baseWithData) (toContinue bool) {
		if gws := GetGuildWorkshopBase(base); gws != nil {
			if gws.isComplete {
				// 已经建造完成，定时扣繁荣度
				if gws.nextReduceProsperityTime <= ctimeUnix32 {
					duration := timeutil.DurationMarshal32(r.services.datas.GuildGenConfig().WorkshopReduceProsperityDuration)
					gws.nextReduceProsperityTime = ctimeUnix32 + duration

					r.reduceGuildWorkshopProsperity(base, gws, r.services.datas.GuildGenConfig().WorkshopReduceProsperity, ctime, false)
				}

			} else {
				// 还没建造完成，看时间是否到了，到了就更新一下
				if gws.endTime <= ctimeUnix32 {
					r.completeGuildWorkshop(base, gws, ctime)
				}
			}
		}
		return true
	})
}

var showWorkshopTrueMsg  = guild.NewS2cShowWorkshopNotExistMsg(true).Static()

func (r *Realm) reduceGuildWorkshopProsperity(base *baseWithData, gws *guildWorkshopBase, toReduce uint64, ctime time.Time, addBeenHurtTimes bool) {

	newProsperity := base.ReduceProsperity(toReduce)

	// 更新联盟工坊数据
	var memberIds []int64
	r.services.guildService.FuncGuild(gws.guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.WithField("guild", gws.guildId).Error("联盟工坊扣繁荣度，联盟不存在")
			return
		}

		if newProsperity <= 0 {
			// 移除联盟工坊
			g.SetWorkshop(nil)
			memberIds = g.AllUserMemberIds()
			return
		}

		workshop := g.GetWorkshop()
		if workshop == nil {
			// 重新弄一个
			ctime := r.services.timeService.CurrentTime()
			endTime := timeutil.Marshal32(ctime)
			startTime := endTime - timeutil.DurationMarshal32(r.services.datas.GuildGenConfig().WorkshopBuildDuration)
			workshop = sharedguilddata.NewDefaultWorkshop(startTime, endTime, base.BaseX(), base.BaseY(), newProsperity)

			g.SetWorkshop(workshop)
			g.Complete(base.Prosperity())
		}

		workshop.SetProsperity(newProsperity)

		if addBeenHurtTimes {
			g.IncTodayWorkshopBeenHurtTimes()
		}
	})

	if len(memberIds) > 0 {
		r.services.world.MultiSend(memberIds, showWorkshopTrueMsg)
	}

	if newProsperity > 0 {
		gws.ClearUpdateBaseInfoMsg()
		// 通知所有正在看这张地图的人
		r.broadcastToCared(base, gws.newUpdateProgressBarMsg(), 0)
		return
	}

	// 移除这个联盟工坊
	r.removeRealmBaseNoReason(base, removeBaseTypeBroken, ctime)
}

func (r *Realm) completeGuildWorkshop(base *baseWithData, gws *guildWorkshopBase, ctime time.Time) {
	gws.isComplete = true
	gws.nextReduceProsperityTime = timeutil.Marshal32(ctime.Add(r.services.datas.GuildGenConfig().WorkshopReduceProsperityDuration))

	// 加满繁荣度
	gws.AddProsperity(gws.prosperityCapcity)

	gws.ClearUpdateBaseInfoMsg()

	// 通知所有正在看这张地图的人
	r.broadcastToCared(base, gws.newUpdateProgressBarMsg(), 0)

	// 更新联盟工坊数据
	var leaderId int64
	r.services.guildService.FuncGuild(gws.guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.WithField("guild", gws.guildId).Error("更新联盟工坊完成建设，联盟不存在")
			return
		}

		workshop := g.GetWorkshop()
		if workshop == nil {
			// 重新弄一个
			ctime := r.services.timeService.CurrentTime()
			endTime := timeutil.Marshal32(ctime)
			startTime := endTime - timeutil.DurationMarshal32(r.services.datas.GuildGenConfig().WorkshopBuildDuration)
			workshop = sharedguilddata.NewDefaultWorkshop(startTime, endTime, base.BaseX(), base.BaseY(), base.Prosperity())

			g.SetWorkshop(workshop)
		}

		todayCompleted := g.GetWorkshopTodayCompleted()

		g.Complete(base.Prosperity())

		// add log
		if d := r.dep.Datas().TextHelp().GuildGongfangBuildComplete; d != nil {
			workshop.AddLog(d.New().JsonString(), ctime)
		}

		// 每日给奖励上限
		if !todayCompleted {
			// 给联盟中的每个人都加一个奖励
			g.WalkMember(func(member *sharedguilddata.GuildMember) {
				member.SetWorkshopPrizeCount(r.services.datas.GuildGenConfig().WorkshopPrizeInitCount)

				// 消息
				count := member.GetWorkshopPrizeCount()
				r.services.world.SendFunc(member.Id(), func() pbutil.Buffer {
					return region.NewS2cUpdateGuildWorkshopPrizeCountMsg(u64.Int32(count))
				})
			})
		}

		leaderId = g.LeaderId()
	})

	if data := r.services.datas.TextHelp().GuildWorkshopCompletedChat; data != nil {
		showJson := fmt.Sprintf(`{"workshop_id":%d,"pos_x":%d,"pos_y":%d}`, gws.id, gws.baseX, gws.baseY)
		r.dep.Chat().SysChat(leaderId, gws.guildId, shared_proto.ChatType_ChatGuild, showJson, shared_proto.ChatMsgType_ChatMsgGuildWorkshopCompleted, false, true, false, true)
	}
}

func (r *Realm) AddGuildWorkshop(guildId int64, x, y int, startTime, endTime int32) (processed bool) {

	var updateGuildFunc func(g *sharedguilddata.Guild)
	baseExist := false
	processed = r.queueFunc(true, func() {

		baseId := npcid.NewGuildWorkshopId(guildId)
		if base := r.getBase(baseId); base != nil {
			baseExist = true

			// 更新联盟工坊数据
			if gwb := GetGuildWorkshopBase(base); gwb != nil {
				updateGuildFunc = getUpdateGuildWorkshopFunc(guildId, base.BaseX(), base.BaseY(), gwb.startTime, gwb.endTime,
					gwb.isComplete, gwb.Prosperity())
			}
			return
		}

		data := r.services.datas.GuildGenConfig().WorkshopBase
		base := r.newGuildWorkshopBase(baseId, data, guildId, x, y, startTime, endTime, false)
		r.addBaseToMap(base)

		// 通知所有正在看这张地图的人
		r.broadcastBaseInfoToCared(base, addBaseTypeNewHero, 0)

		if gwb := GetGuildWorkshopBase(base); gwb != nil {
			updateGuildFunc = getUpdateGuildWorkshopFunc(guildId, base.BaseX(), base.BaseY(), gwb.startTime, gwb.endTime,
				gwb.isComplete, gwb.Prosperity())
		}
	})

	// 更新联盟工坊数据
	if updateGuildFunc != nil {
		r.services.guildService.FuncGuild(guildId, updateGuildFunc)

		if g := r.services.guildService.GetSnapshot(guildId); g != nil {
			r.services.guildService.SelfGuildMsgCache().Clear(g.Id)
			r.services.world.MultiSend(g.UserMemberIds, guild.SELF_GUILD_CHANGED_S2C)
		}
	}

	if !processed || baseExist {
		r.CancelReservedPos(x, y)
	}
	return
}

func getUpdateGuildWorkshopFunc(guildId int64, x, y int, startTime, endTime int32, isComplete bool, prosperity uint64) func(g *sharedguilddata.Guild) {
	return func(g *sharedguilddata.Guild) {
		if g == nil {
			// 联盟不存在？
			logrus.WithField("guild", guildId).Error("更新联盟工坊完成建设，联盟不存在")
			return
		}

		workshop := g.GetWorkshop()
		if workshop != nil {
			logrus.WithField("guild", guildId).Errorf("添加联盟工坊到野外，城池已经存在，更新下数据")
			workshop.SetData(x, y, startTime, endTime, isComplete, prosperity)
			return
		}

		workshop = sharedguilddata.NewDefaultWorkshop(startTime, endTime, x, y, prosperity)
		g.SetWorkshop(workshop)

		if isComplete {
			g.Complete(prosperity)
		}
	}
}

// 破坏联盟工坊
func (r *Realm) HurtGuildWorkshop(hc iface.HeroController, targetGuildId int64) (processed, baseNotExist bool) {
	processed = r.queueFunc(true, func() {

		baseId := npcid.NewGuildWorkshopId(targetGuildId)
		base := r.getBase(baseId)
		if base == nil {
			logrus.Debug("破坏联盟工坊，base == nil")
			baseNotExist = true
			return
		}

		gws := GetGuildWorkshopBase(base)
		if gws == nil {
			logrus.Debug("破坏联盟工坊，gws == nil")
			baseNotExist = true
			return
		}

		ctime := r.services.timeService.CurrentTime()
		if gws.isComplete {
			// 扣繁荣度
			r.reduceGuildWorkshopProsperity(base, gws, r.services.datas.GuildGenConfig().WorkshopHurtProsperity, ctime, true)
		} else {
			// 加建设完成时间
			toAdd := timeutil.DurationMarshal32(r.services.datas.GuildGenConfig().WorkshopHurtDuration)
			gws.startTime += toAdd
			gws.endTime += toAdd

			gws.ClearUpdateBaseInfoMsg()
			// 通知所有正在看这张地图的人
			r.broadcastToCared(base, gws.newUpdateProgressBarMsg(), 0)

			// 更新联盟工坊数据
			r.services.guildService.FuncGuild(gws.guildId, func(g *sharedguilddata.Guild) {
				if g == nil {
					logrus.WithField("guild", gws.guildId).Error("破坏联盟工坊，联盟不存在（region）")
					return
				}

				workshop := g.GetWorkshop()
				if workshop == nil {
					logrus.WithField("guild", gws.guildId).Error("破坏联盟工坊，workshop == nil")
					return
				}

				workshop.SetTime(gws.startTime, gws.endTime)
				g.IncTodayWorkshopBeenHurtTimes()
			})
		}
	})
	return
}

// 到达联盟工坊事件
func (t *troop) workshopArrivedEvent(r *Realm) {
	t.event.RemoveFromQueue()
	t.event = nil

	logrus.WithField("troopid", t.Id()).Debug("处理到达联盟工坊地点")

	troopState := t.state
	if t.state != realmface.MovingToWorkshopBuild && t.state != realmface.MovingToWorkshopProd {
		logrus.WithField("troopid", t.Id()).WithField("state", t.state).Error("realm竟然触发了workshopArrivedEvent, 但是队伍状态不是MovingToWorkshopBuild或者MovingToWorkshopProd")
		return
	}

	defenserBase := t.targetBase
	if defenserBase == nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了workshopArrivedEvent, 但是troop的targetBase为nil")
		return
	}

	// 到了就到了，改一下状态，没你什么事了

	ctime := r.services.timeService.CurrentTime()
	t.backHome(r, ctime, true, false)

	gws := GetGuildWorkshopBase(defenserBase)
	if gws == nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了workshopArrivedEvent, gws == nil")
		return
	}

	var heroName string
	var heroId int64
	r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
		heroName = hero.Name()
		heroId = hero.Id()
		hero.UpdateTroop(t, false)

		if gws.isComplete {
			// 加生产次数
			hero.GetWorkshopOutputTimes().ReduceOneTimes(ctime.Add(-time.Second))
			result.Add(region.NewS2cUpdateHeroOutputWorkshopTimesMsg(hero.GetWorkshopOutputTimes().StartTimeUnix32()))
		} else {
			// 加建设次数
			hero.GuildWorkshop().IncDailyBuildTimes()
			result.Add(region.NewS2cUpdateHeroBuildWorkshopTimesMsg(u64.Int32(hero.GuildWorkshop().GetDailyBuildTimes())))
		}

		result.Ok()
	})

	// 加繁荣度，加产出次数
	if gws.isComplete {
		// 加繁荣度
		gws.AddProsperity(r.services.datas.GuildGenConfig().WorkshopAddProsperity)
	} else {
		// 减少建筑时间
		duration := timeutil.DurationMarshal32(r.services.datas.GuildGenConfig().WorkshopHeroBuildDuration)
		gws.startTime -= duration
		gws.endTime -= duration
	}

	gws.ClearUpdateBaseInfoMsg()
	// 通知所有正在看这张地图的人
	r.broadcastToCared(defenserBase, gws.newUpdateProgressBarMsg(), 0)

	// 更新联盟工坊数据
	var allMemberIds []int64
	var stageIndx int
	r.services.guildService.FuncGuild(gws.guildId, func(g *sharedguilddata.Guild) {
		if g == nil {
			logrus.WithField("guild", gws.guildId).Error("更新联盟工坊完成建设，联盟不存在")
			return
		}

		base := defenserBase

		workshop := g.GetWorkshop()
		if workshop == nil {
			// 重新弄一个
			workshop = sharedguilddata.NewDefaultWorkshop(gws.startTime, gws.endTime, base.BaseX(), base.BaseY(), base.Prosperity())
			g.SetWorkshop(workshop)
		}
		// log
		if workshop.IsComplete() {
			if d := r.services.datas.TextHelp().GuildGongfangEfficiencyAdd; d != nil {
				text := d.New().WithClickHeroFields(data.KeyName, heroName, heroId)
				text.WithNum(r.services.datas.GuildGenConfig().WorkshopAddProsperity)
				workshop.AddLog(text.JsonString(), ctime)
			}
			// 增加联盟任务进度
			if g.LevelData().Level >= r.services.datas.GuildGenConfig().TaskOpenLevel {
				data := r.services.datas.GetGuildTaskData(uint64(server_proto.GuildTaskType_Workshop))
				if g.AddGuildTaskProgress(data, 1) {
					allMemberIds = g.AllUserMemberIds()
					stageIndx = g.GetGuildTaskStageIndex(data.TaskType)
				}
			}

		} else if d := r.services.datas.TextHelp().GuildGongfangBuildTimeReduce; d != nil {
			text := d.New().WithClickHeroFields(data.KeyName, heroName, heroId)
			minutes := int64(r.services.datas.GuildGenConfig().WorkshopHeroBuildDuration / time.Minute)
			text.WithNum(minutes)
			workshop.AddLog(text.JsonString(), ctime)
		}

		isOutput := false
		if gws.isComplete {

			if g.GetWorkshopOutputPrizeCount() < len(r.services.datas.GuildGenConfig().WorkshopMaxOutput) {
				// 加产出
				g.AddWorkshopOutput(r.services.datas.GuildGenConfig().WorkshopAddOutput)

				isOutput = g.TryWorkshopOutput(r.services.datas.GuildGenConfig().WorkshopMaxOutput)
			} else {
				logrus.WithField("guildId", gws.guildId).Debug("联盟今日产出个数已达上限")
			}

		}

		workshop.SetData(base.BaseX(), base.BaseY(), gws.startTime, gws.endTime, gws.isComplete, base.Prosperity())

		if isOutput {
			// 给联盟中的每个人都加一个奖励
			g.WalkMember(func(member *sharedguilddata.GuildMember) {
				if member.GetWorkshopPrizeCount() < r.services.datas.GuildGenConfig().WorkshopPrizeMaxCount {
					member.IncWorkshopPrizeCount()

					// 消息
					count := member.GetWorkshopPrizeCount()
					r.services.world.SendFunc(member.Id(), func() pbutil.Buffer {
						return region.NewS2cUpdateGuildWorkshopPrizeCountMsg(u64.Int32(count))
					})
				}
			})
			// log
			if d := r.services.datas.TextHelp().GuildGongfangPrizeSend; d != nil {
				workshop.AddLog(d.New().JsonString(), ctime)
			}
		}
	})

	if len(allMemberIds) > 0 {
		r.dep.World().MultiSend(allMemberIds, guild.NewS2cNoticeTaskStageUpdateMsg(int32(server_proto.GuildTaskType_Workshop), int32(stageIndx)))
	}

	// 对应飘字
	switch troopState {
	case realmface.MovingToWorkshopBuild:
		minute := imath.Max(int(r.services.datas.GuildGenConfig().WorkshopHeroBuildDuration/time.Minute), 1)
		r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().GuildWorkshopBuildShow.New().
					WithNum(minute).JsonString())
		})
	case realmface.MovingToWorkshopProd:
		r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().GuildWorkshopProdShow.Text.KeysOnlyJson())
		})
	}

}

// 到达联盟工坊领取奖励事件
func (t *troop) workshopPrizeArrivedEvent(r *Realm) {
	t.event.RemoveFromQueue()
	t.event = nil

	logrus.WithField("troopid", t.Id()).Debug("处理到达联盟工坊领取奖励地点")

	if t.state != realmface.MovingToWorkshopPrize {
		logrus.WithField("troopid", t.Id()).WithField("state", t.state).Error("realm竟然触发了workshopPrizeArrivedEvent, 但是队伍状态不是MovingToWorkshopPrize")
		return
	}

	defenserBase := t.targetBase
	if defenserBase == nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了workshopPrizeArrivedEvent, 但是troop的targetBase为nil")
		return
	}

	// 到了就到了，改一下状态，没你什么事了

	ctime := r.services.timeService.CurrentTime()
	t.backHome(r, ctime, true, true)

	// 领取奖励
	// 加繁荣度，加产出次数
	if gws := GetGuildWorkshopBase(defenserBase); gws != nil {
		// 更新联盟数据
		// 更新联盟工坊数据
		var addPrizeCount uint64
		r.services.guildService.FuncGuild(gws.guildId, func(g *sharedguilddata.Guild) {
			if g == nil {
				logrus.WithField("guild", gws.guildId).Error("领取联盟工坊奖励，联盟不存在")
				return
			}

			member := g.GetMember(t.startingBase.Id())
			if member == nil {
				logrus.WithField("guild", gws.guildId).Error("领取联盟工坊奖励，member不存在")
				return
			}

			addPrizeCount = member.ClearAndGetWorkshopPrizeCount()
			if addPrizeCount > 0 {
				// 消息
				count := member.GetWorkshopPrizeCount()
				r.services.world.SendFunc(member.Id(), func() pbutil.Buffer {
					return region.NewS2cUpdateGuildWorkshopPrizeCountMsg(u64.Int32(count))
				})
			}
		})

		if addPrizeCount > 0 {

			var addedCoef float64
			if area := r.regionData.GetAreaByPos(gws.baseX, gws.baseY); area != nil {
				addedCoef = area.WorkshopPrizeCoef
			}

			r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {

				b := hero.Domestic().GetBuilding(shared_proto.BuildingType_GUAN_FU)
				if b == nil {
					return
				}

				data := r.services.datas.GetGuanFuLevelData(b.Level)
				if data == nil {
					return
				}

				toAdd := data.WorkshopPrize
				if addPrizeCount > 1 || addedCoef > 0 {
					if addedCoef > 0 {
						multi := float64(addPrizeCount) * (1 + addedCoef)
						toAdd = toAdd.MultiCoef(multi)
					} else {
						toAdd = toAdd.Multiple(addPrizeCount)
					}
				}

				hctx := heromodule.NewContext(r.dep, operate_type.GuildCollectWorkshopPrize)
				heromodule.AddPrize(hctx, hero, result, toAdd, ctime)

				result.Ok()
			})
		}
	}

	// 对应飘字
	r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
		return misc.NewS2cScreenShowWordsMsg(
			r.getTextHelp().GuildWorkshopPrizeShow.Text.KeysOnlyJson())
	})
}
