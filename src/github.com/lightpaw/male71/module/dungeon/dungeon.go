package dungeon

import (
	"github.com/lightpaw/logrus"
	dungeon2 "github.com/lightpaw/male7/config/dungeon"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/dungeon"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/gamelogs"
)

func NewDungeonModule(dep iface.ServiceDep, fightService iface.FightXService, guildSnapshotService iface.GuildSnapshotService) *DungeonModule {
	m := &DungeonModule{
		dep:                  dep,
		configDatas:          dep.Datas(),
		miscData:             dep.Datas().DungeonMiscData(),
		timeService:          dep.Time(),
		fightService:         fightService,
		guildSnapshotService: guildSnapshotService,
	}
	return m
}

//gogen:iface
type DungeonModule struct {
	dep                  iface.ServiceDep
	configDatas          iface.ConfigDatas
	miscData             *dungeon2.DungeonMiscData
	timeService          iface.TimeService
	fightService         iface.FightXService
	guildSnapshotService iface.GuildSnapshotService
}

//gogen:iface c2s_challenge
func (m *DungeonModule) ProcessChallenge(proto *dungeon.C2SChallengeProto, hc iface.HeroController) {
	dungeonData := m.configDatas.GetDungeonData(u64.FromInt32(proto.GetId()))
	if dungeonData == nil {
		logrus.Debugf("推图副本，副本没找到")
		hc.Send(dungeon.ERR_CHALLENGE_FAIL_NOT_FOUND)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Level() < dungeonData.UnlockHeroLevel {
			logrus.Debugf("推图副本，君主等级不足")
			hc.Send(dungeon.ERR_CHALLENGE_FAIL_LEVEL_NOT_ENOUGH)
			return
		}

		ctime := m.timeService.CurrentTime()

		heroDungeon := hero.Dungeon()
		// 用体力值取代原有挑战次数限制
		if dungeonData.Sp > 0 && !hero.HasEnoughSp(dungeonData.Sp) {
			logrus.Debugf("推图副本，体力值不足")
			hc.Send(dungeon.ERR_CHALLENGE_FAIL_SP_NOT_ENOUGH)
			return
		}

		if !heroDungeon.IsUnlockPreDungeon(dungeonData) {
			logrus.Debugf("推图副本，前置副本未通关")
			hc.Send(dungeon.ERR_CHALLENGE_FAIL_PREV_NOT_PASS)
			return
		}

		// 霸业前置条件限制
		if hero.TaskList().GetCompletedBaYeStage() < dungeonData.UnlockBayeStage {
			logrus.Debugf("推图副本，前置霸业阶段未通关")
			hc.Send(dungeon.ERR_CHALLENGE_FAIL_BAYE_NOT_PASS)
			return
		}

		// 有些关卡配有每日挑战次数限制
		if heroDungeon.IsPassLimit(dungeonData, 1) {
			logrus.Debugf("推图副本，当日挑战次数上限")
			hc.Send(dungeon.ERR_CHALLENGE_FAIL_PASS_LIMIT)
			return
		}

		// 打架
		troops, failType := hero.BuildPveCaptainProtosWithHeroAndMonster(true, shared_proto.PveTroopType_DUNGEON, dungeonData.YuanJunData, true)
		switch failType {
		case entity.SUCCESS:
			break
		case entity.SERVER_ERROR:
			result.Add(dungeon.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		case entity.CAPTAIN_COUNT_NOT_ENOUGH:
			result.Add(dungeon.ERR_CHALLENGE_FAIL_CAPTAIN_NOT_FULL)
			return
		default:
			logrus.Errorf("推图副本，未处理的错误类型 %v", failType)
			result.Add(dungeon.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		attacker := hero.GenCombatPlayerProtoWithCaptainProtos(troops, m.guildSnapshotService.GetSnapshot)
		if attacker == nil {
			logrus.Errorf("推图副本，构建Player出错")
			result.Add(dungeon.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		var totalCaptainCount, totalSoldier int32
		for _, t := range attacker.Troops {
			totalCaptainCount++
			totalSoldier += t.Captain.TotalSoldier
		}

		tfctx := entity.NewTlogFightContext(operate_type.BattleHUANJING, dungeonData.Id, 0, 0)
		response := m.fightService.SendFightRequest(tfctx, dungeonData.CombatScene, hc.Id(), dungeonData.Monster.GetNpcId(), attacker, dungeonData.Monster.GetPlayer())
		if response == nil {
			logrus.Errorf("推图副本，response==nil")
			result.Add(dungeon.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		if response.ReturnCode != 0 {
			logrus.Errorf("推图副本，，战斗计算发生错误，%s", response.ReturnMsg)
			result.Add(dungeon.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		passCostSeconds := u64.DivideTimes(u64.FromInt32(response.TotalFrame), m.configDatas.CombatConfig().FramePerSecond)

		if !response.AttackerWin {
			// 挑战失败
			result.Add(dungeon.NewS2cChallengeMsg(proto.GetId(), response.Link, must.Marshal(response.AttackerShare),
				[]byte{}, false, false, []bool{},
				u64.Int32(passCostSeconds),
				u64.Int32(heroDungeon.GetChapterStar(dungeonData.GetChapterStarDungeonIds())),
				u64.Int32(heroDungeon.GetLimitDungeonPassTimes(dungeonData.Id)), false))

			gamelogs.DungeonChallengeLog(hero.Pid(), hero.Sid(), hero.Id(), dungeonData.ChapterId, dungeonData.Id, uint64(dungeonData.Type), 0, false)
			return
		}

		// 通关奖励
		isPass := heroDungeon.IsPass(dungeonData)
		prize := dungeonData.Plunder.Try()
		if !isPass { // 第一次挑战
			if dungeonData.FirstPassPrize != nil {
				prize = resdata.NewPrizeBuilder().Add(prize).Add(dungeonData.FirstPassPrize).Build()
			}
			if m.configDatas.MiscConfig().CloseFightGuideDungeonId == dungeonData.Id {
				if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_CLOSE_FIGHT_GUIDE) {
					result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_CLOSE_FIGHT_GUIDE)))
				}
			}
		}
		// 扣除体力值
		hero.ReduceSp(dungeonData.Sp)
		result.Add(domestic.NewS2cUpdateSpMsg(u64.Int32(hero.GetSp())))

		hctx := heromodule.NewContext(m.dep, operate_type.DungeonChallenge)
		heromodule.AddPveTroopPrize(hctx, hero, result, prize, shared_proto.PveTroopType_DUNGEON, ctime)

		// 先临时发给测试数据给前端(剩余士兵百分比：80，通关消耗了：100秒，武将死亡个数：0)
		var aliveCaptainCount, aliveSoldier int32
		for id, soldier := range response.AttackerAliveSoldier {
			if id != 0 {
				if soldier > 0 {
					aliveCaptainCount++
					aliveSoldier += soldier
				}
			}
		}

		resultHpPercent := u64.Max(u64.FromInt32(aliveSoldier*100/totalSoldier), 1) // 血量最低1%
		captainDeathCount := u64.FromInt32(totalCaptainCount - aliveCaptainCount)

		enabledStars, starCount := dungeonData.CalculateEnabledStars(resultHpPercent, passCostSeconds, captainDeathCount)
		refreshed := heroDungeon.Pass(dungeonData, enabledStars)
		result.Add(dungeon.NewS2cChallengeMsg(proto.GetId(), response.Link, must.Marshal(response.AttackerShare),
			must.Marshal(prize.Encode()), true, !isPass, enabledStars, u64.Int32(passCostSeconds),
			u64.Int32(heroDungeon.GetChapterStar(dungeonData.GetChapterStarDungeonIds())),
			u64.Int32(heroDungeon.GetLimitDungeonPassTimes(dungeonData.Id)), refreshed))

		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_DUNGEON)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HAS_CHALLENGE_DUNGEON)

		heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_DUNGEON, dungeonData.Id)

		gamelogs.DungeonChallengeLog(hero.Pid(), hero.Sid(), hero.Id(), dungeonData.ChapterId, dungeonData.Id, uint64(dungeonData.Type), starCount, true)

		result.Changed()
		result.Ok()
	})

}

//gogen:iface c2s_collect_chapter_prize
func (m *DungeonModule) ProcessCollectChapterPrize(proto *dungeon.C2SCollectChapterPrizeProto, hc iface.HeroController) {
	chapterData := m.configDatas.GetDungeonChapterData(u64.FromInt32(proto.GetId()))
	if chapterData == nil {
		logrus.WithField("id", proto.GetId()).Debugln("难度数据没找到")
		hc.Send(dungeon.ERR_COLLECT_CHAPTER_PRIZE_FAIL_NOT_FOUND)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDungeon := hero.Dungeon()

		if !heroDungeon.IsPass(chapterData.LastDungeon) {
			logrus.WithField("id", proto.GetId()).Debugln("最后一节没通关")
			hc.Send(dungeon.ERR_COLLECT_CHAPTER_PRIZE_FAIL_NOT_PASS)
			return
		}

		if heroDungeon.IsCollectChapterPrize(chapterData) {
			logrus.WithField("id", proto.GetId()).Debugln("奖励已经领取了")
			hc.Send(dungeon.ERR_COLLECT_CHAPTER_PRIZE_FAIL_COLLECTED)
			return
		}

		heroDungeon.CollectChapterPrize(chapterData)

		result.Add(dungeon.NewS2cCollectChapterPrizeMsg(proto.GetId()))

		hctx := heromodule.NewContext(m.dep, operate_type.DungeonCollectChapterPrize)
		heromodule.AddPrize(hctx, hero, result, chapterData.PassPrize, m.timeService.CurrentTime())

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DungeonModule) ProcessCollectPassDungeonPrize(proto *dungeon.C2SCollectPassDungeonPrizeProto, hc iface.HeroController) {

	data := m.configDatas.GetDungeonData(u64.FromInt32(proto.Id))
	if data == nil {
		logrus.WithField("id", proto.GetId()).Debugln("领取推图通关奖励，数据没找到")
		hc.Send(dungeon.ERR_COLLECT_PASS_DUNGEON_PRIZE_FAIL_INVALID_ID)
		return
	}

	if data.PassPrize == nil {
		logrus.WithField("id", proto.GetId()).Debugln("领取推图通关奖励，该副本没有配置通关奖励")
		hc.Send(dungeon.ERR_COLLECT_PASS_DUNGEON_PRIZE_FAIL_NOT_FOUND)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDungeon := hero.Dungeon()

		if !heroDungeon.IsPass(data) {
			logrus.WithField("id", proto.GetId()).Debugln("最后一节没通关")
			hc.Send(dungeon.ERR_COLLECT_PASS_DUNGEON_PRIZE_FAIL_NOT_PASS)
			return
		}

		if heroDungeon.IsCollectPassDungeonPrize(data) {
			logrus.WithField("id", proto.GetId()).Debugln("奖励已经领取了")
			hc.Send(dungeon.ERR_COLLECT_PASS_DUNGEON_PRIZE_FAIL_COLLECTED)
			return
		}

		heroDungeon.CollectPassDungeonPrize(data)

		result.Add(dungeon.NewS2cCollectPassDungeonPrizeMsg(proto.GetId()))

		hctx := heromodule.NewContext(m.dep, operate_type.DungeonCollectPassDungeonPrize)
		heromodule.AddPrize(hctx, hero, result, data.PassPrize, m.timeService.CurrentTime())

		result.Changed()
		result.Ok()
	})
}

//gogen:iface c2s_auto_challenge
func (m *DungeonModule) ProcessAutoChallenge(proto *dungeon.C2SAutoChallengeProto, hc iface.HeroController) {
	data := m.configDatas.GetDungeonData(u64.FromInt32(proto.GetId()))
	if data == nil {
		logrus.WithField("id", proto.GetId()).Debugln("副本没找到")
		hc.Send(dungeon.ERR_AUTO_CHALLENGE_FAIL_NOT_FOUND)
		return
	}

	if proto.GetTimes() <= 0 {
		logrus.WithField("id", proto.GetId()).WithField("times", proto.GetTimes()).Debugln("次数非法")
		hc.Send(dungeon.ERR_AUTO_CHALLENGE_FAIL_INVALID_AUTO_TIMES)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDungeon := hero.Dungeon()

		if !heroDungeon.IsPass(data) {
			logrus.WithField("id", proto.GetId()).WithField("times", proto.GetTimes()).Debugln("副本未通关")
			hc.Send(dungeon.ERR_AUTO_CHALLENGE_FAIL_NOT_PASS)
			return
		}

		if !heroDungeon.IsCanAutoPass(data) {
			logrus.WithField("id", proto.GetId()).Debugln("无法扫荡未满星通关的副本")
			hc.Send(dungeon.ERR_AUTO_CHALLENGE_FAIL_NOT_FULL_STAR)
			return
		}

		ctime := m.timeService.CurrentTime()

		times := u64.Min(u64.FromInt32(proto.Times), m.miscData.MaxAutoTimes)
		if data.Sp <= 0 {
			logrus.WithField("关卡配置sp", data.Sp).Debugln("不大于0")
			hc.Send(dungeon.ERR_AUTO_CHALLENGE_FAIL_SP_NOT_ENOUGH)
			return
		}
		times = u64.Min(hero.GetSp()/data.Sp, times)
		if times <= 0 {
			logrus.WithField("id", proto.GetId()).Debugln("玩家没有体力了")
			hc.Send(dungeon.ERR_AUTO_CHALLENGE_FAIL_SP_NOT_ENOUGH)
			return
		}

		if heroDungeon.IsPassLimit(data, times) {
			logrus.WithField("id", proto.GetId()).Debugln("扫荡次数超出每日通关上限")
			hc.Send(dungeon.ERR_AUTO_CHALLENGE_FAIL_PASS_LIMIT)
			return
		}

		heroDungeon.AddLimitDungeonPassTimes(data.Id, times)

		// 到这里先扣体力值，然后加奖励
		hero.ReduceSp(data.Sp * times)
		result.Add(domestic.NewS2cUpdateSpMsg(u64.Int32(hero.GetSp())))

		prizeBytes := make([][]byte, 0, times)
		prizeBuidler := resdata.NewPrizeBuilder()
		for i := uint64(0); i < times; i++ {
			prize := data.Plunder.Try()
			prizeBytes = append(prizeBytes, must.Marshal(prize.Encode()))
			prizeBuidler.Add(prize)
		}

		result.Add(dungeon.NewS2cAutoChallengeMsg(proto.Id, prizeBytes, u64.Int32(heroDungeon.GetLimitDungeonPassTimes(data.Id))))

		hctx := heromodule.NewContext(m.dep, operate_type.DungeonAutoChallenge)
		prize := prizeBuidler.Build()
		heromodule.AddPveTroopPrize(hctx, hero, result, prize, shared_proto.PveTroopType_DUNGEON, ctime)

		hero.HistoryAmount().Increase(server_proto.HistoryAmountType_AccumDungeonAutoComplete, times)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_AUTO_DUNGEON)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_AUTO_DUNGEON) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_AUTO_DUNGEON)))
		}

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DungeonModule) ProcessCollectChapterStarPrize(proto *dungeon.C2SCollectChapterStarPrizeProto, hc iface.HeroController) {
	chapterData := m.configDatas.GetDungeonChapterData(u64.FromInt32(proto.GetId()))
	if chapterData == nil {
		logrus.WithField("id", proto.GetId()).Debugln("领取章节星数奖励，无效的章节ID")
		hc.Send(dungeon.ERR_COLLECT_CHAPTER_STAR_PRIZE_FAIL_INVALID_ID)
		return
	}
	collectN := int(proto.GetCollectN())
	if collectN <= 0 || collectN > len(chapterData.Star) {
		logrus.WithField("index", collectN).Debugln("领取章节星数奖励，无效的奖励下标")
		hc.Send(dungeon.ERR_COLLECT_CHAPTER_STAR_PRIZE_FAIL_INVALID_PRIZE_INDEX)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDungeon := hero.Dungeon()
		if chapterStar := heroDungeon.GetChapterStar(chapterData.GetStarDungeonIds()); chapterStar < chapterData.Star[collectN-1] {
			logrus.WithField("star", chapterStar).Debugln("领取章节星数奖励，星数不足")
			hc.Send(dungeon.ERR_COLLECT_CHAPTER_STAR_PRIZE_FAIL_STAR_NOT_ENOUGH)
			return
		}
		if !heroDungeon.TryCollectChapterStarPrize(chapterData.Id, collectN) {
			logrus.WithField("index", collectN).Debugln("领取章节星数奖励，奖励已经领取了")
			hc.Send(dungeon.ERR_COLLECT_CHAPTER_STAR_PRIZE_FAIL_COLLECTED)
			return
		}
		prize := chapterData.StarPrize[collectN-1]
		hctx := heromodule.NewContext(m.dep, operate_type.DungeonCollectChapterStarPrize)
		heromodule.AddPrize(hctx, hero, result, prize, m.timeService.CurrentTime())

		result.Add(dungeon.NewS2cCollectChapterStarPrizeMsg(proto.Id, proto.CollectN, must.Marshal(prize.Encode())))

		result.Changed()
		result.Ok()
	})
}
