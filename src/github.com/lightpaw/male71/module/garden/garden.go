package garden

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/garden"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/service/operate_type"
)

func NewGardenModule(dep iface.ServiceDep, serverConfig *kv.IndividualServerConfig, seasonService iface.SeasonService) *GardenModule {
	m := &GardenModule{
		dep:                 dep,
		datas:               dep.Datas(),
		serverConfig:        serverConfig,
		timeService:         dep.Time(),
		seasonService:       seasonService,
		heroDataService:     dep.HeroData(),
		heroSnapshotService: dep.HeroSnapshot(),
		guildService:        dep.Guild(),
		world:               dep.World(),
	}
	return m
}

//gogen:iface
type GardenModule struct {
	dep           iface.ServiceDep
	datas         iface.ConfigDatas
	serverConfig  *kv.IndividualServerConfig
	timeService   iface.TimeService
	seasonService iface.SeasonService

	heroDataService     iface.HeroDataService
	heroSnapshotService iface.HeroSnapshotService

	guildService iface.GuildService

	world iface.WorldService
}

var emptyListTreasuryTreeHeroMsg = garden.NewS2cListTreasuryTreeHeroProtoMsg(&garden.S2CListTreasuryTreeHeroProto{}).Static()

// 这个还是要保留，之后删掉一些无用字段。帮助记录用 ProcessListHelpMe
//gogen:iface c2s_list_treasury_tree_hero
func (m *GardenModule) ProcessListTreasuryTreeHero(hc iface.HeroController) {

	var guildId int64
	var helpMeHeroIds []int64
	var helpMeSeasons []shared_proto.Season
	hasError := hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		guildId = hero.GuildId()

		// 帮助列表
		helpMeHeroIds, helpMeSeasons = hero.TreasuryTree().HelpMeInfo()

		result.Ok()
	})

	if hasError {
		logrus.Error("获取摇钱树玩家列表，获取玩家联盟id失败")
		hc.Send(emptyListTreasuryTreeHeroMsg)
		return
	}

	s2cProto := &garden.S2CListTreasuryTreeHeroProto{}

	if guildId != 0 {
		g := m.guildService.GetSnapshot(guildId)
		if g != nil {
			careWaterTimesMap := make(map[int64]uint64)
			for _, memberId := range g.UserMemberIds {
				if memberId != hc.Id() {
					member := m.heroSnapshotService.Get(memberId)
					if member != nil && member.GuildId == guildId {
						s2cProto.HeroId = append(s2cProto.HeroId, member.IdBytes)
						s2cProto.HeroName = append(s2cProto.HeroName, member.Name)
						s2cProto.HeroHead = append(s2cProto.HeroHead, member.Head)
						s2cProto.HeroGuild = append(s2cProto.HeroGuild, i64.Int32(guildId))
						s2cProto.HeroFriend = append(s2cProto.HeroFriend, false) // 目前都不是好友
						s2cProto.HeroWaterTimes = append(s2cProto.HeroWaterTimes, u64.Int32(member.TreasuryTreeWaterTimes))

						careWaterTimesMap[member.Id] = member.TreasuryTreeWaterTimes
					}
				}
			}
			hc.SetCareWaterTimesMap(careWaterTimesMap)

			s2cProto.GuildId = append(s2cProto.GuildId, i64.Int32(guildId))
			s2cProto.FlagName = append(s2cProto.FlagName, g.FlagName)
		}
	}

	var otherGuildIds []int64
	for i, helpMeHeroId := range helpMeHeroIds {

		if helpMeHeroId == hc.Id() {
			season := shared_proto.Season_SUMMER
			if i < len(helpMeSeasons) {
				season = helpMeSeasons[i]
			}

			s2cProto.HelpMeHeroId = append(s2cProto.HelpMeHeroId, hc.IdBytes())
			s2cProto.HelpMeHeroName = append(s2cProto.HelpMeHeroName, "")
			s2cProto.HelpMeHeroGuild = append(s2cProto.HelpMeHeroGuild, 0)
			s2cProto.HelpMeHeroSeason = append(s2cProto.HelpMeHeroSeason, int32(season))
		} else {
			helpMeHero := m.heroSnapshotService.Get(helpMeHeroId)
			if helpMeHero != nil {
				season := shared_proto.Season_SUMMER
				if i < len(helpMeSeasons) {
					season = helpMeSeasons[i]
				}

				s2cProto.HelpMeHeroId = append(s2cProto.HelpMeHeroId, helpMeHero.IdBytes)
				s2cProto.HelpMeHeroName = append(s2cProto.HelpMeHeroName, helpMeHero.Name)
				s2cProto.HelpMeHeroGuild = append(s2cProto.HelpMeHeroGuild, i64.Int32(guildId))
				s2cProto.HelpMeHeroSeason = append(s2cProto.HelpMeHeroSeason, int32(season))

				if helpMeHero.GuildId != 0 && helpMeHero.GuildId != guildId && !i64.Contains(otherGuildIds, helpMeHero.GuildId) {
					otherGuildIds = append(otherGuildIds, helpMeHero.GuildId)
				}
			}
		}
	}

	for _, guildId := range otherGuildIds {
		g := m.guildService.GetSnapshot(guildId)
		if g != nil {
			s2cProto.GuildId = append(s2cProto.GuildId, i64.Int32(guildId))
			s2cProto.FlagName = append(s2cProto.FlagName, g.FlagName)
		}
	}

	hc.Send(garden.NewS2cListTreasuryTreeHeroProtoMsg(s2cProto))
}

var emptyListTreasuryTreeWaterTimesMsg = garden.NewS2cListTreasuryTreeTimesMsg(nil, nil).Static()

//gogen:iface
func (m *GardenModule) ProcessListHelpMe(proto *garden.C2SListHelpMeProto, hc iface.HeroController) {
	targetId, ok := idbytes.ToId(proto.TargetId)
	if !ok {
		hc.Send(garden.ERR_LIST_HELP_ME_FAIL_INVALID_TARGET)
		return
	}

	var helpMeHeroIds []int64
	var helpMeSeasons []shared_proto.Season

	heroSnapshot := m.heroSnapshotService.Get(targetId)
	if heroSnapshot == nil {
		logrus.Debugf("摇钱树-ProcessListHelpMe，herosnapshot 里找不到这个人 id:%v", targetId)
		hc.Send(garden.ERR_LIST_HELP_ME_FAIL_INVALID_TARGET)
		return
	}

	helpMeHeroIds = heroSnapshot.HelpMeHeroIds
	helpMeSeasons = heroSnapshot.HelpMeSeasons

	s2cProto := &garden.S2CListHelpMeProto{}
	s2cProto.TargetId = proto.TargetId
	for i, heroId := range helpMeHeroIds {
		helpMeHero := m.heroSnapshotService.Get(heroId)
		if helpMeHero == nil {
			continue
		}

		season := shared_proto.Season_SUMMER
		if i < len(helpMeSeasons) {
			season = helpMeSeasons[i]
		}
		s2cProto.HelpMeHeroId = append(s2cProto.HelpMeHeroId, helpMeHero.IdBytes)
		s2cProto.HelpMeHeroName = append(s2cProto.HelpMeHeroName, helpMeHero.Name)
		s2cProto.HelpMeHeroGuild = append(s2cProto.HelpMeHeroGuild, i64.Int32(helpMeHero.GuildId))
		s2cProto.HelpMeHeroSeason = append(s2cProto.HelpMeHeroSeason, int32(season))

		flagName := ""
		if helpMeHero.GuildId != 0 {
			g := m.guildService.GetSnapshot(helpMeHero.GuildId)
			if g != nil {
				flagName = g.FlagName
			}
		}
		s2cProto.HelpMeHeroFlagName = append(s2cProto.HelpMeHeroFlagName, flagName)
	}

	hc.Send(garden.NewS2cListHelpMeProtoMsg(s2cProto))
}

//gogen:iface c2s_list_treasury_tree_times
func (m *GardenModule) ProcessListTreasuryTreeTimes(hc iface.HeroController) {

	var guildId int64
	var friendIds []int64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		// 我的好友
		friendIds = hero.Relation().FriendIds()

		// 我的联盟
		guildId = hero.GuildId()

		result.Ok()
	})

	if guildId != 0 {
		g := m.guildService.GetSnapshot(guildId)
		if g != nil {
			for _, memberId := range g.UserMemberIds {
				if memberId != hc.Id() {
					member := m.heroSnapshotService.Get(memberId)
					if member != nil && member.GuildId == guildId {
						friendIds = i64.AddIfAbsent(friendIds, member.Id)
					}
				}
			}
		}
	}

	var heroIdBytes [][]byte
	var waterTimesArray []int32
	if len(friendIds) > 0 {
		heroWaterTimesMap := hc.GetCareWaterTimesMap()

		for _, heroId := range friendIds {
			hero := m.heroSnapshotService.Get(heroId)
			if hero == nil {
				continue
			}

			waterTimes, exist := heroWaterTimesMap[heroId]
			if !exist || waterTimes != hero.TreasuryTreeWaterTimes {
				heroWaterTimesMap[heroId] = hero.TreasuryTreeWaterTimes

				heroIdBytes = append(heroIdBytes, hero.IdBytes)
				waterTimesArray = append(waterTimesArray, u64.Int32(hero.TreasuryTreeWaterTimes))
			}
		}
	}

	if len(heroIdBytes) <= 0 {
		hc.Send(emptyListTreasuryTreeWaterTimesMsg)
		return
	}

	hc.Send(garden.NewS2cListTreasuryTreeTimesMsg(heroIdBytes, waterTimesArray))
}

//gogen:iface
func (m *GardenModule) ProcessWaterTreasuryTree(proto *garden.C2SWaterTreasuryTreeProto, hc iface.HeroController) {

	targetId, ok := idbytes.ToId(proto.Target)
	if !ok {
		logrus.Debug("摇钱树浇水，无效的目标")
		hc.Send(garden.ERR_WATER_TREASURY_TREE_FAIL_INVALID_TARGET)
		return
	}

	if npcid.IsNpcId(targetId) {
		logrus.Debug("摇钱树浇水，目标是个npc?")
		hc.Send(garden.ERR_WATER_TREASURY_TREE_FAIL_INVALID_TARGET)
		return
	}

	if targetId == 0 {
		targetId = hc.Id()
	}

	if targetId != hc.Id() {
		target := m.heroSnapshotService.Get(targetId)
		if target == nil {
			logrus.Debug("摇钱树浇水，加载目标的HeroSnapshot失败")
			hc.Send(garden.ERR_WATER_TREASURY_TREE_FAIL_INVALID_TARGET)
			return
		}

		if target.TreasuryTreeWaterTimes >= m.datas.GardenConfig().TreasuryTreeFullTimes {
			logrus.Debug("摇钱树浇水，目标的摇钱树已经满了")
			hc.Send(garden.ERR_WATER_TREASURY_TREE_FAIL_FULL)
			return
		}
	}

	ctime := m.timeService.CurrentTime()

	season := m.seasonService.Season().Season
	var helpMeHeroName string
	var helpMeGuildId int64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		helpMeHeroName = hero.Name()
		helpMeGuildId = hero.GuildId()

		if hero.TreasuryTree().IsWatered(targetId) {
			logrus.Debug("摇钱树浇水，这个人已经浇过水了")
			result.Add(garden.ERR_WATER_TREASURY_TREE_FAIL_INVALID_TARGET)
			return
		}

		if targetId == hero.Id() {

			// 从来都没给自己浇过水，第一次可以随便浇水
			isFirstWater := false
			if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_WATER) {
				result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_WATER)))
				isFirstWater = true
			}

			// 给自己浇水
			oldTimes := hero.TreasuryTree().WaterTimes()
			newTimes := oldTimes
			if hero.TreasuryTree().WaterTimes() >= m.datas.GardenConfig().TreasuryTreeFullTimes {
				// 第一次浇水，就算满了，也给他浇
				if !isFirstWater {
					logrus.Debug("摇钱树浇水，自己的摇钱树已经满了")
					result.Add(garden.ERR_WATER_TREASURY_TREE_FAIL_FULL)
					return
				}
			} else {
				newTimes = hero.TreasuryTree().IncreseWaterTimes()
			}

			if newTimes >= m.datas.GardenConfig().TreasuryTreeFullTimes {
				// 加完次数已经满了，设置领取时间
				collectTime := ctime.Add(m.datas.GardenConfig().TreasuryTreeCollectDuration)
				hero.TreasuryTree().SetCollectInfo(collectTime, season)

				result.Add(garden.NewS2cUpdateSelfTreasuryTreeFullMsg(int32(season), timeutil.Marshal32(collectTime)))
			}

			result.Add(garden.NewS2cWaterTreasuryTreeMsg(hero.IdBytes(), u64.Int32(newTimes)))
			result.Changed()

			hero.TreasuryTree().AddHelpMeInfo(targetId, season, m.datas.GardenConfig().TreasuryTreeHelpMeLogCount)

			m.dep.Tlog().TlogCareFlow(hero, operate_type.GardenCareSelf, u64.FromInt64(targetId), hero.Name(), oldTimes, newTimes)
		} else {
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumHelpWater)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_HELP_WATER)
		}

		hero.TreasuryTree().AddWaterHero(targetId)

		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_TreasuryTree)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_TREASURY_TREE)

		result.Ok()
	}) {
		return
	}

	if targetId == hc.Id() {
		return
	}

	var helpMeGuildFlagName string
	if g := m.guildService.GetSnapshot(helpMeGuildId); g != nil {
		helpMeGuildFlagName = g.FlagName
	} else {
		helpMeGuildId = 0
	}

	// 给别人浇水
	helpMeHeroId := hc.Id()
	helpMeHeroIdBytes := hc.IdBytes()
	var newTimes uint64
	m.heroDataService.FuncWithSend(targetId, func(hero *entity.Hero, result herolock.LockResult) {
		oldTimes := hero.TreasuryTree().WaterTimes()
		newTimes = oldTimes
		if hero.TreasuryTree().WaterTimes() < m.datas.GardenConfig().TreasuryTreeFullTimes {
			newTimes = hero.TreasuryTree().IncreseWaterTimes()
			result.Changed()

			// 加被帮助日志
			hero.TreasuryTree().AddHelpMeInfo(helpMeHeroId, season, m.datas.GardenConfig().TreasuryTreeHelpMeLogCount)
			result.Add(garden.NewS2cUpdateSelfTreasuryTreeTimesMsg(u64.Int32(newTimes), helpMeHeroIdBytes, helpMeHeroName, i64.Int32(helpMeGuildId), helpMeGuildFlagName))

			if newTimes >= m.datas.GardenConfig().TreasuryTreeFullTimes {
				// 加完次数已经满了，设置领取时间
				collectTime := ctime.Add(m.datas.GardenConfig().TreasuryTreeCollectDuration)
				hero.TreasuryTree().SetCollectInfo(collectTime, season)

				result.Add(garden.NewS2cUpdateSelfTreasuryTreeFullMsg(int32(season), timeutil.Marshal32(collectTime)))
			}

		}

		result.Ok()

		m.dep.Tlog().TlogCareFlow(hero, operate_type.GardenCareFriends, u64.FromInt64(targetId), hero.Name(), oldTimes, newTimes)
	})

	hc.Send(garden.NewS2cWaterTreasuryTreeMsg(idbytes.ToBytes(targetId), u64.Int32(newTimes)))
}

//gogen:iface c2s_collect_treasury_tree_prize
func (m *GardenModule) ProcessCollectTreasureTreePrize(hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.TreasuryTree().WaterTimes() < m.datas.GardenConfig().TreasuryTreeFullTimes {
			logrus.Debug("摇钱树领奖，摇钱树还未满")
			result.Add(garden.ERR_COLLECT_TREASURY_TREE_PRIZE_FAIL_NO_FULL)
			return
		}

		collectTime, collectSeason := hero.TreasuryTree().CollectInfo()
		ctime := m.timeService.CurrentTime()
		if ctime.Before(collectTime) {
			logrus.Debug("摇钱树领奖，领奖时间未到")
			result.Add(garden.ERR_COLLECT_TREASURY_TREE_PRIZE_FAIL_COUNTDOWN)
			return
		}

		data := m.datas.TreasuryTreeData().Must(uint64(collectSeason))

		hero.CollectTreasuryTreePrize()

		result.Add(garden.COLLECT_TREASURY_TREE_PRIZE_S2C)

		hctx := heromodule.NewContext(m.dep, operate_type.GardenCollectTreasureTreePrize)
		heromodule.AddPrize(hctx, hero, result, data.Prize, ctime)

		result.Changed()
		result.Ok()
	})

}
