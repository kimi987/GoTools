package xuanyuan

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/module/xuanyuan/xym"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/gen/pb/xuanyuan"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/pbutil"
	"context"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/config/xuanydata"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"sync/atomic"
	"github.com/lightpaw/male7/gen/pb/misc"
	"sync"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/module/rank/ranklist"
	"github.com/lightpaw/male7/service/operate_type"
)

// 轩辕会武

func NewXuanyuanModule(datas iface.ConfigDatas, time iface.TimeService, heroService iface.HeroDataService,
	heroSnapshotService iface.HeroSnapshotService, guildSnapshotService iface.GuildSnapshotService,
	fightService iface.FightXService, dbService iface.DbService, tickService iface.TickerService,
	world iface.WorldService, dep iface.ServiceDep, rankModule iface.RankModule) *XuanyuanModule {

	var data []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		data, err = dbService.LoadKey(ctx, server_proto.Key_Xuanyuan)
		return
	})
	if err != nil {
		logrus.WithError(err).Panic("加载轩辕会武模块数据失败")
	}

	dailyTick := tickService.GetDailyTickTime()

	mgr := xym.NewManager(datas.XuanyuanMiscData().RankCount)
	if len(data) > 0 {
		proto := &server_proto.XuanyuanModuleProto{}
		if err := proto.Unmarshal(data); err != nil {
			logrus.WithError(err).Panic("解析轩辕会武模块数据失败，server_proto.XuanyuanModuleProto")
		}

		mgr.Unmarshal(proto)

		if r := mgr.Get(); r != nil && !timeutil.IsZero(r.GetUpdateTime()) {
			mgr.Update(dailyTick.GetPrevTickTime(), false)
		}

	}

	m := &XuanyuanModule{
		datas:                datas,
		time:                 time,
		heroService:          heroService,
		heroSnapshotService:  heroSnapshotService,
		guildSnapshotService: guildSnapshotService,
		fightService:         fightService,
		dbService:            dbService,
		world:                world,
		dep:                  dep,
		tickService:          tickService,
		rankModule:           rankModule,
		manager:              mgr,
	}

	m.resetTickTimeRef.Store(tickdata.Copy(dailyTick))

	m.stopSaveLoop = tickService.TickPer10Minute("定时保存轩辕会武模块数据", func(tick tickdata.TickTime) {
		ctxfunc.Timeout3s(func(ctx context.Context) (err error) {
			m.save(ctx)
			return nil
		})
	})

	m.stopResetLoop = tickService.TickPerDay("轩辕会武每日重置", func(tick tickdata.TickTime) {
		m.resetDaily(tick, false)
	})

	heromodule.RegisterHeroOnlineListener(m)

	return m
}

//gogen:iface
type XuanyuanModule struct {
	datas iface.ConfigDatas

	time iface.TimeService

	heroService         iface.HeroDataService
	heroSnapshotService iface.HeroSnapshotService

	guildSnapshotService iface.GuildSnapshotService

	fightService iface.FightXService

	dbService iface.DbService

	world iface.WorldService

	dep iface.ServiceDep

	tickService iface.TickerService

	rankModule iface.RankModule

	stopSaveLoop  iface.Func
	stopResetLoop iface.Func
	resetLocker   sync.Mutex

	resetTickTimeRef atomic.Value

	manager *xym.XuanyuanManager
}

var (
	rankIsEmptyTrue = xuanyuan.NewS2cRankIsEmptyMsg(true).Static()
)

func (m *XuanyuanModule) OnHeroOnline(hc iface.HeroController) {

	if m.getRankCount() <= 0 {
		hc.Send(rankIsEmptyTrue)
		return
	}
}

func (m *XuanyuanModule) getRankCount() int {
	ro := m.manager.Get()
	if ro != nil {
		return ro.RankCount()
	}
	return 0
}

func (m *XuanyuanModule) Close() {
	if m.stopResetLoop != nil {
		m.stopResetLoop()
	}

	if m.stopSaveLoop != nil {
		m.stopSaveLoop()
	}

	ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		m.save(ctx)
		return
	})
}

func (m *XuanyuanModule) GetResetTickTime() tickdata.TickTime {
	return m.getResetTickData()
}

func (m *XuanyuanModule) getResetTickData() *tickdata.TickData {
	return m.resetTickTimeRef.Load().(*tickdata.TickData)
}

func (m *XuanyuanModule) resetDaily(nextTickTime tickdata.TickTime, gmReset bool) {

	oldRankCount := m.getRankCount()

	func() {
		m.resetLocker.Lock()
		defer m.resetLocker.Unlock()

		// 更新数据
		if m.manager.Update(nextTickTime.GetPrevTickTime(), gmReset) {
			prevTickData := m.getResetTickData()
			m.resetTickTimeRef.Store(tickdata.Copy(nextTickTime))
			prevTickData.Close()
		}
	}()

	// 更新排行榜
	if ro := m.manager.Get(); ro != nil && ro.RankCount() > 0 {
		var rankArray []rankface.RankObj
		ro.Range(func(hero *xym.XyRankHero) (toContinue bool) {
			rank := hero.Rank()
			if rank > 0 {
				rankArray = append(rankArray,
					ranklist.NewXuanyRankObj(m.heroSnapshotService.Get, hero.Id(),
						hero.GetScore(), hero.GetWin(), hero.GetLose(), uint64(rank)))
			}
			return true
		})

		m.rankModule.UpdateXuanyRankList(rankArray)
	}

	if oldRankCount == 0 && m.getRankCount() > 0 {
		m.world.Broadcast(xuanyuan.NewS2cRankIsEmptyMsg(false))
	}
}

func (m *XuanyuanModule) save(ctx context.Context) {

	proto := &server_proto.XuanyuanModuleProto{}
	m.manager.Encode(proto)

	if err := m.dbService.SaveKey(ctx, server_proto.Key_Xuanyuan, must.Marshal(proto)); err != nil {
		logrus.WithError(err).Error("保存轩辕会武模块数据失败")
	}
}

func (m *XuanyuanModule) GmReset() {
	m.resetDaily(m.tickService.GetDailyTickTime(), true)
}

func (m *XuanyuanModule) AddChallenger(heroId int64, player *shared_proto.CombatPlayerProto, score uint64) {
	m.manager.AddChallenger(heroId, score, 0, 0, player)
}

var emptySelfMsg = xuanyuan.NewS2cSelfInfoMsg(0, 0, 0, 0, 1, 0, nil).Static()

//gogen:iface c2s_self_info
func (m *XuanyuanModule) ProcessSelfInfo(hc iface.HeroController) {

	ro := m.manager.Get()
	if ro == nil {
		logrus.Debug("轩辕会武获取自己数据，ro == nil")
		hc.Send(emptySelfMsg)
		return
	}

	rankCount := ro.RankCount()
	if rankCount <= 0 {
		logrus.Debug("轩辕会武获取自己数据，ro.RankCount() == 0")
		hc.Send(emptySelfMsg)
		return
	}

	toSendProto := &xuanyuan.S2CSelfInfoProto{}
	self := ro.GetHero(hc.Id())
	selfRank := rankCount + 1
	if self != nil {
		selfRank = self.Rank()
	}

	for _, rangeData := range m.datas.GetXuanyuanRangeDataArray() {
		startPos := imath.Max(selfRank-rangeData.HighRank, 1)
		endPos := imath.Min(selfRank-rangeData.LowRank, rankCount)

		// 有挑战目标
		for i := startPos; i <= endPos; i++ {
			target := ro.GetHeroByRank(i)
			if target.Id() == hc.Id() {
				continue
			}

			if toSendProto.FirstTargetRank == 0 {
				toSendProto.RangeId = u64.Int32(rangeData.Id)
				toSendProto.FirstTargetRank = int32(i)
			}

			toSendProto.Targets = append(toSendProto.Targets, target.EncodeTarget(m.heroSnapshotService.Get))
		}

		if toSendProto.RangeId > 0 {
			break
		}
	}

	if self != nil {
		toSendProto.Rank = int32(self.Rank())
		toSendProto.Score = u64.Int32(self.GetScore())
		toSendProto.Win = u64.Int32(self.GetWin())
		toSendProto.Lose = u64.Int32(self.GetLose())
	} else {
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			toSendProto.Score = u64.Int32(hero.Xuanyuan().GetScore())
			toSendProto.Win = u64.Int32(hero.Xuanyuan().GetWin())
			toSendProto.Lose = u64.Int32(hero.Xuanyuan().GetLose())

			result.Ok()
		})
	}

	hc.Send(xuanyuan.NewS2cSelfInfoProtoMsg(toSendProto))
}

//gogen:iface
func (m *XuanyuanModule) ProcessListTarget(proto *xuanyuan.C2SListTargetProto, hc iface.HeroController) {

	rangeData := m.datas.GetXuanyuanRangeData(u64.FromInt32(proto.RangeId))
	if rangeData == nil {
		logrus.Debug("获取轩辕会武挑战目标，无效的range")
		hc.Send(xuanyuan.ERR_LIST_TARGET_FAIL_INVALID_RANGE)
		return
	}

	ro := m.manager.Get()
	if ro == nil {
		logrus.Debug("获取轩辕会武挑战目标，ro == nil")
		hc.Send(xuanyuan.NewS2cListTargetMsg(proto.RangeId, 0, nil))
		return
	}

	rankCount := ro.RankCount()
	if rankCount <= 0 {
		logrus.Debug("获取轩辕会武挑战目标，rankCount == 0")
		hc.Send(xuanyuan.NewS2cListTargetMsg(proto.RangeId, 0, nil))
		return
	}

	// 获取自己的排名，然后根据排名以及类型，获取到排名区间，找对手，lock英雄获取数据，发送数据
	self := ro.GetHero(hc.Id())
	myPos := rankCount + 1
	if self != nil {
		myPos = self.Rank()
	}

	startPos := imath.Max(myPos-rangeData.HighRank, 1)
	endPos := imath.Min(myPos-rangeData.LowRank, rankCount)

	toSendProto := &xuanyuan.S2CListTargetProto{
		RangeId: proto.RangeId,
	}
	for i := startPos; i <= endPos; i++ {
		rankHero := ro.GetHeroByRank(i)
		if rankHero == nil {
			logrus.WithField("rank", i).Warn("获取轩辕会武挑战目标，根据排名找不到英雄id")
			continue
		}

		if rankHero.Id() == hc.Id() {
			// 防御性
			continue
		}

		if toSendProto.FirstTargetRank == 0 {
			toSendProto.FirstTargetRank = int32(i)
		}
		toSendProto.Targets = append(toSendProto.Targets, rankHero.EncodeTarget(m.heroSnapshotService.Get))
	}

	hc.Send(xuanyuan.NewS2cListTargetProtoMsg(toSendProto))
}

//gogen:iface
func (m *XuanyuanModule) ProcessQueryTargetTroop(proto *xuanyuan.C2SQueryTargetTroopProto, hc iface.HeroController) {

	targetId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debug("轩辕会武获取部队，无效的id")
		hc.Send(xuanyuan.ERR_QUERY_TARGET_TROOP_FAIL_INVALID_ID)
		return
	}

	ro := m.manager.Get()
	if ro == nil {
		logrus.Debug("轩辕会武获取部队，ro == nil")
		hc.Send(xuanyuan.ERR_QUERY_TARGET_TROOP_FAIL_INVALID_ID)
		return
	}

	target := ro.GetHero(targetId)
	if target == nil {
		logrus.Debug("轩辕会武获取部队，target == nil")
		hc.Send(xuanyuan.ERR_QUERY_TARGET_TROOP_FAIL_INVALID_ID)
		return
	}

	hc.Send(target.GetQueryTargetTroopMsg())
}

//gogen:iface
func (m *XuanyuanModule) ProcessChallenge(proto *xuanyuan.C2SChallengeProto, hc iface.HeroController) {

	targetId, ok := idbytes.ToId(proto.Id)
	if !ok {
		logrus.Debug("轩辕会武挑战，无效的id")
		hc.Send(xuanyuan.ERR_CHALLENGE_FAIL_INVALID_ID)
		return
	}

	if targetId == hc.Id() {
		logrus.Debug("轩辕会武挑战，挑战的目标是自己")
		hc.Send(xuanyuan.ERR_CHALLENGE_FAIL_INVALID_ID)
		return
	}

	ro := m.manager.Get()
	if ro == nil {
		logrus.Debug("轩辕会武挑战，ro == nil")
		hc.Send(xuanyuan.ERR_CHALLENGE_FAIL_INVALID_ID)
		return
	}

	rankCount := ro.RankCount()
	if rankCount <= 0 {
		logrus.Debug("轩辕会武挑战，榜单为空")
		hc.Send(xuanyuan.ERR_CHALLENGE_FAIL_INVALID_ID)
		return
	}

	target := ro.GetHero(targetId)
	if target == nil {
		logrus.Debug("轩辕会武挑战，目标不存在")
		hc.Send(xuanyuan.ERR_CHALLENGE_FAIL_INVALID_ID)
		return
	}

	version, targetMirror := target.GetMirror()
	if version != int64(proto.Version) {
		logrus.Debug("轩辕会武挑战，镜像数据已刷新")
		hc.Send(xuanyuan.ERR_CHALLENGE_FAIL_VERSION)

		// 再发送数据，刷新目标
		hc.Send(target.GetQueryTargetTroopMsg())
		return
	}

	self := ro.GetHero(hc.Id())
	selfRank := ro.RankCount() + 1
	if self != nil {
		selfRank = self.Rank()
	}

	diff := selfRank - target.Rank()
	rangeData := m.datas.XuanyuanMiscData().GetRangeByDiff(diff)
	if rangeData == nil {
		logrus.WithField("selfRank", selfRank).WithField("challengeRank", target.Rank()).Debug("轩辕会武挑战，目标名次无效")
		hc.Send(xuanyuan.ERR_CHALLENGE_FAIL_INVALID_ID)
		return
	}

	// lock英雄获取镜像，打架
	var attacker *shared_proto.CombatPlayerProto
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		// 挑战次数
		if hero.Xuanyuan().GetChallengeTargetLen() >= u64.Int(m.datas.XuanyuanMiscData().ChallengeTimesLimit) {
			logrus.Debug("轩辕会武挑战，挑战次数不足")
			result.Add(xuanyuan.ERR_CHALLENGE_FAIL_TIMES_LIMIT)
			return
		}

		if hero.Xuanyuan().ContainsChallengeTarget(target.Id()) {
			logrus.Debug("轩辕会武挑战，目标已经挑战过了")
			result.Add(xuanyuan.ERR_CHALLENGE_FAIL_CHALLENGED)
			return
		}

		player, failType := hero.GenCombatPlayerProto(true, shared_proto.PveTroopType_DUNGEON, m.guildSnapshotService.GetSnapshot)
		switch failType {
		case entity.SUCCESS:
			break
		case entity.SERVER_ERROR:
			logrus.Debug("轩辕会武挑战，目标已经挑战过了")
			result.Add(xuanyuan.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		case entity.CAPTAIN_COUNT_NOT_ENOUGH:
			logrus.Debug("轩辕会武挑战，上阵武将个数不足")
			result.Add(xuanyuan.ERR_CHALLENGE_FAIL_CAPTAIN_NOT_ENOUGH)
			return
		default:
			logrus.Error("轩辕会武挑战，未处理的错误类型 %v", failType)
			result.Add(xuanyuan.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		attacker = player

		result.Ok()
	}) {
		return
	}

	tfctx := entity.NewTlogFightContext(operate_type.BattleXuanYuan, 0, 0, 0)
	response := m.fightService.SendFightRequest(tfctx, rangeData.CombatScene, hc.Id(), target.Id(), attacker, targetMirror)
	if response == nil {
		logrus.Errorf("轩辕会武挑战，response==nil")
		hc.Send(xuanyuan.ERR_CHALLENGE_FAIL_SERVER_ERROR)
		return
	}

	if response.ReturnCode != 0 {
		logrus.Errorf("轩辕会武挑战，战斗计算发生错误，%s", response.ReturnMsg)
		hc.Send(xuanyuan.ERR_CHALLENGE_FAIL_SERVER_ERROR)
		return
	}

	toAddScore := rangeData.LoseScore
	if response.AttackerWin {
		toAddScore = rangeData.WinScore
	}

	ctime := m.time.CurrentTime()

	// 更新自己数据
	var newScore, newWin, newLose uint64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		hero.Xuanyuan().AddChallengeTarget(target.Id())
		if response.AttackerWin {
			// 挑战成功
			hero.Xuanyuan().IncWin()
		} else {
			// 挑战失败，自己加积分
			hero.Xuanyuan().IncLose()
		}

		hero.Xuanyuan().AddScore(toAddScore)

		newScore = hero.Xuanyuan().GetScore()
		newWin = hero.Xuanyuan().GetWin()
		newLose = hero.Xuanyuan().GetLose()

		result.Add(xuanyuan.NewS2cUpdateXyInfoMsg(u64.Int32(newScore), u64.Int32(newWin), u64.Int32(newLose)))

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_XUAN_YUAN) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_XUAN_YUAN)))
		}

		if self != nil {
			if response.AttackerWin {
				// 挑战成功
				self.IncWin()
			} else {
				// 挑战失败，自己加积分
				self.IncLose()
			}
			self.SetScore(newScore)

			// 设置挑战镜像
			self.SetMirror(attacker, timeutil.Marshal64(ctime))
		}

		if hero.HistoryAmount().GetAmount(false, server_proto.HistoryAmountType_XuanyuanHisMaxScore) < newScore {
			hero.HistoryAmount().SetAmount(false, server_proto.HistoryAmountType_XuanyuanHisMaxScore, newScore)
		}
		if hero.HistoryAmount().GetAmount(true, server_proto.HistoryAmountType_XuanyuanHisMaxScore) < newScore {
			hero.HistoryAmount().SetAmount(true, server_proto.HistoryAmountType_XuanyuanHisMaxScore, newScore)
		}
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_XUANYUAN_SCORE)

		result.Changed()
		result.Ok()
	}) {
		return
	}

	if self == nil {
		// 一次性数据保存进去
		m.manager.AddChallenger(hc.Id(), newScore, newWin, newLose, attacker)
	}

	hc.Send(xuanyuan.NewS2cChallengeMsg(proto.Id, response.Link, u64.Int32(toAddScore)))

	// 更新对手数据（失败也不管）
	var toSub uint64
	m.heroService.FuncWithSend(target.Id(), func(hero *entity.Hero, result herolock.LockResult) {
		if response.AttackerWin {
			// 被挑战成功
			hero.Xuanyuan().IncLose()

			// 每天最多掉X分
			toSub = u64.Sub(m.datas.XuanyuanMiscData().DailyMaxLostScore, hero.Xuanyuan().GetLostScore())
			toSub = u64.Min(toSub, rangeData.DefenseLostScore)
			if toSub > 0 {
				hero.Xuanyuan().SubScore(toSub)
				hero.Xuanyuan().AddLostScore(toSub)
			}
		} else {
			// 被挑战失败，
			hero.Xuanyuan().IncWin()
		}

		newScore := hero.Xuanyuan().GetScore()
		newWin := hero.Xuanyuan().GetWin()
		newLose := hero.Xuanyuan().GetLose()
		result.AddFunc(func() pbutil.Buffer {
			return xuanyuan.NewS2cUpdateXyInfoMsg(u64.Int32(newScore), u64.Int32(newWin), u64.Int32(newLose))
		})

		if response.AttackerWin {
			// 被挑战成功
			target.IncLose()
			if toSub > 0 {
				target.SetScore(newScore)
			}
		} else {
			// 被挑战失败
			target.IncWin()
		}

		result.Ok()
	})

	toRacordHero := func(player *shared_proto.CombatPlayerProto, scoreChaged uint64) *shared_proto.XuanyuanRecordHeroProto {
		races := make([]shared_proto.Race, 5)
		for i, v := range player.Troops {
			idx := i + 1
			if v.FightIndex > 0 {
				idx = int(v.FightIndex)
			}

			if idx <= len(races) {
				races[idx-1] = v.Captain.Race
			}
		}

		return &shared_proto.XuanyuanRecordHeroProto{
			Hero:         player.Hero,
			FightAmount:  player.TotalFightAmount,
			Race:         races,
			ScoreChanged: u64.Int32(scoreChaged),
		}
	}

	record := &shared_proto.XuanyuanRecordProto{
		IsAttackerWin: response.AttackerWin,
		Time:          timeutil.Marshal32(ctime),
		Link:          response.Link,
		Attacker:      toRacordHero(attacker, toAddScore),
		Defender:      toRacordHero(targetMirror, toSub),
	}

	recordBytes := must.Marshal(record)

	// 战斗记录
	var recordId int64
	ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		if recordId, err = m.dbService.InsertXuanyRecord(ctx, hc.Id(), recordBytes); err != nil {
			logrus.WithError(err).Error("保存轩辕会武战斗记录失败")
		}
		return
	})
	hc.Send(xuanyuan.NewS2cAddRecordMsg(i64.Int32(recordId), recordBytes))

	ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		if recordId, err = m.dbService.InsertXuanyRecord(ctx, targetId, recordBytes); err != nil {
			logrus.WithError(err).Error("保存轩辕会武战斗记录失败")
		}
		return
	})
	m.world.SendFunc(targetId, func() pbutil.Buffer {
		return xuanyuan.NewS2cAddRecordMsg(i64.Int32(recordId), recordBytes)
	})

}

//gogen:iface
func (m *XuanyuanModule) ProcessListRecord(proto *xuanyuan.C2SListRecordProto, hc iface.HeroController) {

	var ids []int64
	var datas [][]byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		ids, datas, err = m.dbService.LoadXuanyRecord(ctx, hc.Id(), int64(proto.Id), proto.Up)
		return
	})
	if err != nil {
		logrus.WithError(err).Error("请求轩辕会武战斗记录失败")
		hc.Send(xuanyuan.NewS2cListRecordMsg(proto.Id, proto.Up, nil, nil))
		return
	}

	hc.Send(xuanyuan.NewS2cListRecordMsg(proto.Id, proto.Up, i64.Int32Array(ids), datas))
}

//gogen:iface c2s_collect_rank_prize
func (m *XuanyuanModule) ProcessCollectRankPrize(hc iface.HeroController) {
	rank := 0

	r := m.manager.Get()
	if r != nil {
		self := r.GetHero(hc.Id())
		if self != nil {
			rank = self.Rank()
		}
	}

	rankPrizeData := xuanydata.GetRankPrizeDataByRank(m.datas.GetXuanyuanRankPrizeDataArray(), rank)
	if rankPrizeData == nil {
		logrus.Debug("轩辕会武领取排名奖励，没有奖励可以领取")
		hc.Send(xuanyuan.ERR_COLLECT_RANK_PRIZE_FAIL_COLLECTED)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.Xuanyuan().GetRankPrizeCollected() {
			logrus.Debug("轩辕会武领取排名奖励，已领取过奖励")
			result.Add(xuanyuan.ERR_COLLECT_RANK_PRIZE_FAIL_COLLECTED)
			return
		}

		hero.Xuanyuan().SetRankPrizeCollected(true)

		ctime := m.time.CurrentTime()

		gid := m.guildSnapshotService.GetGuildLevel(hero.GuildId())
		prize := rankPrizeData.PlunderPrize.GetGuildPrize(gid)

		hctx := heromodule.NewContext(m.dep, operate_type.XuanYuanCollectRankPrize)
		heromodule.AddPrize(hctx, hero, result, prize, ctime)

		result.Add(xuanyuan.NewS2cCollectRankPrizeMsg(prize.Encode()))
	})

}
