package bai_zhan

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/bai_zhan_data"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/bai_zhan"
	"github.com/lightpaw/male7/module/bai_zhan/bai_zhan_objs"
	"github.com/lightpaw/male7/module/rank/ranklist"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/call"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"sort"
	"sync"
	"time"
	"github.com/lightpaw/male7/service/operate_type"
)

func NewBaiZhanModule(dep iface.ServiceDep, dbService iface.DbService, fightService iface.FightXService,
	tickerService iface.TickerService, guildSnapshotService iface.GuildSnapshotService,
	baiZhanService iface.BaiZhanService, rankModule iface.RankModule, mailModule iface.MailModule) *BaiZhanModule {

	m := &BaiZhanModule{}

	m.dep = dep
	m.dbService = dbService
	m.timeService = dep.Time()
	m.miscData = dep.Datas().BaiZhanMiscData()
	m.configDatas = dep.Datas()
	m.fightService = fightService
	m.tickerService = tickerService
	m.heroDataService = dep.HeroData()
	m.worldService = dep.World()
	m.baiZhanService = baiZhanService
	m.rankModule = rankModule
	m.mailModule = mailModule
	m.heroSnapshotService = dep.HeroSnapshot()

	m.guildSnapshotGetter = guildSnapshotService.GetSnapshot

	m.loopExitNotify = make(chan struct{})
	m.closeNotify = make(chan struct{})

	m.queryRankFunc = func(id int64) uint64 {
		obj := baiZhanService.GetBaiZhanObj(id)
		if obj == nil {
			return 0
		}

		return obj.Rank()
	}

	m.baiZhanLevelDatas = newBaiZhanLevelDatas(m.configDatas.GetJunXianLevelDataArray())
	m.resetPointRankList()

	m.loadData()

	go call.CatchLoopPanic(m.loop, "BaiZhanModule.loop()")

	return m
}

func (m *BaiZhanModule) loadData() {
	var data []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		data, err = m.dbService.LoadKey(ctx, server_proto.Key_BaiZhan)
		return
	})
	if err != nil {
		logrus.WithError(err).Panicln("加载DB中的百战千军模块数据失败")
	}

	if len(data) <= 0 {
		return
	}

	proto := server_proto.BaiZhanServerProto{}
	if err := proto.Unmarshal(data); err != nil {
		logrus.WithError(err).Panicln("BaiZhanServerProto.Unmarshal, 百战千军模块数据失败")
	}

	m.lastResetTime = timeutil.Unix64(proto.LastResetTime)

	m.baiZhanService.Func(func(objs *bai_zhan_objs.BaiZhanObjs) {
		ctime := m.timeService.CurrentTime()

		inRankListObjs := make([]*bai_zhan_objs.HeroBaiZhanObj, 0, 32)
		for _, p := range proto.BaiZhanObjs {
			obj := bai_zhan_objs.NewBaiZhanObj(p.GetId(), m.heroSnapshotService.Get, len(m.configDatas.GetJunXianLevelDataArray()), m.configDatas.JunXianLevelData().MinKeyData, ctime)
			obj.Unmarshall(p, m.configDatas.JunXianLevelData())
			objs.AddBaiZhanObj(obj)

			// 没被移除了的
			if obj.CombatMirror() != nil && !obj.IsJunXianBeenRemoved() {
				// 加镜像
				m.mustBaiZhanLevelData(obj.LevelData()).addMirror(obj)
			}

			if obj.Point() > 0 {
				inRankListObjs = append(inRankListObjs, obj)
			}
		}

		sort.Sort(hero_bai_zhan_obj_slice(inRankListObjs))

		for _, obj := range inRankListObjs {
			m.mustPointRankList(obj.LevelData()).AddOrUpdate(obj)
		}
	})
}

//gogen:iface
type BaiZhanModule struct {
	dep                 iface.ServiceDep
	dbService           iface.DbService
	timeService         iface.TimeService
	miscData            *bai_zhan_data.BaiZhanMiscData
	configDatas         iface.ConfigDatas
	fightService        iface.FightXService
	tickerService       iface.TickerService
	heroDataService     iface.HeroDataService
	worldService        iface.WorldService
	baiZhanService      iface.BaiZhanService
	rankModule          iface.RankModule
	mailModule          iface.MailModule
	heroSnapshotService iface.HeroSnapshotService

	guildSnapshotGetter func(int64) *guildsnapshotdata.GuildSnapshot

	baiZhanLevelDatas     []*bai_zhan_level_data      // 百战等级数据
	baiZhanPointRankLists []*bai_zhan_point_rank_list // 百战等级排行榜数据
	queryRankFunc         func(id int64) uint64

	lastResetTime time.Time

	loopExitNotify chan struct{}
	closeNotify    chan struct{}
	closeOnce      sync.Once
}

func (m *BaiZhanModule) Close() {
	m.closeOnce.Do(func() {
		close(m.closeNotify)
	})
	<-m.loopExitNotify
}

func (m *BaiZhanModule) loop() {
	// daily reset
	dailyTickTime := m.tickerService.GetDailyTickTime()
	if dailyTickTime.GetPrevTickTime().After(m.lastResetTime) {
		m.tryResetDaily(dailyTickTime.GetPrevTickTime())
	}

	// 10分钟保存一次
	saveTick := time.NewTicker(10 * time.Minute)

	defer close(m.loopExitNotify)
	defer saveTick.Stop()

	for {
		select {
		case <-saveTick.C:
			ctxfunc.Timeout3s(func(ctx context.Context) error {
				m.baiZhanService.Func(func(objs *bai_zhan_objs.BaiZhanObjs) {
					m.save(ctx, objs)
				})
				return nil
			})
		case <-dailyTickTime.Tick():
			m.tryResetDaily(dailyTickTime.GetTickTime())
			dailyTickTime = m.tickerService.GetDailyTickTime()
		case <-m.closeNotify:
			m.baiZhanService.Stop(func(objs *bai_zhan_objs.BaiZhanObjs) {
				ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
					m.save(ctx, objs)
					return
				})

			})
			return // quit loop
		}
	}
}

// 保存数据
func (m *BaiZhanModule) save(ctx context.Context, objs *bai_zhan_objs.BaiZhanObjs) {
	serverProto := &server_proto.BaiZhanServerProto{}
	serverProto.BaiZhanObjs = make([]*server_proto.BaiZhanObjServerProto, 0, objs.Count())

	objs.Walk(func(obj *bai_zhan_objs.HeroBaiZhanObj) {
		serverProto.BaiZhanObjs = append(serverProto.BaiZhanObjs, obj.EncodeServer())
	})

	serverProto.LastResetTime = timeutil.Marshal64(m.lastResetTime)

	err := m.dbService.SaveKey(ctx, server_proto.Key_BaiZhan, must.Marshal(serverProto))
	if err != nil {
		logrus.WithError(err).Errorf("保存百战千军数据出错")
	}
}

func (m *BaiZhanModule) mustBaiZhanLevelData(levelData *bai_zhan_data.JunXianLevelData) (data *bai_zhan_level_data) {
	return m.baiZhanLevelDatas[levelData.Level-1]
}

func (m *BaiZhanModule) resetPointRankList() {
	m.baiZhanPointRankLists = NewBaiZhanPointRankLists(m.configDatas)
}

func (m *BaiZhanModule) mustPointRankList(levelData *bai_zhan_data.JunXianLevelData) (data *bai_zhan_point_rank_list) {
	return m.baiZhanPointRankLists[levelData.Level-1]
}

// 单线程调用
func (m *BaiZhanModule) getOrCreateBaiZhanObj(objs *bai_zhan_objs.BaiZhanObjs, id int64) *bai_zhan_objs.HeroBaiZhanObj {
	baiZhanObj := objs.GetBaiZhanObj(id)
	if baiZhanObj == nil {
		// 创建一个新的
		baiZhanObj = bai_zhan_objs.NewBaiZhanObj(id, m.heroSnapshotService.Get, len(m.configDatas.GetJunXianLevelDataArray()), m.configDatas.JunXianLevelData().MinKeyData, m.timeService.CurrentTime())

		objs.AddBaiZhanObj(baiZhanObj)

		m.rankModule.AddOrUpdateRankObj(m.newRankObj(baiZhanObj))
	}

	// 清掉上次的百战等级数据
	if baiZhanObj.IsJunXianBeenRemoved() {
		baiZhanObj.ReaddRemovedJunXian(m.timeService.CurrentTime())
		m.rankModule.AddOrUpdateRankObj(m.newRankObj(baiZhanObj))
	}

	return baiZhanObj
}

func (m *BaiZhanModule) newRankObj(obj *bai_zhan_objs.HeroBaiZhanObj) *ranklist.BaiZhanRankObj {
	lastJunXianLevelData := obj.OldJunXianLevelData()
	if lastJunXianLevelData == nil {
		lastJunXianLevelData = obj.LevelData()
	}
	return ranklist.NewBaiZhanRankObj(m.heroSnapshotService.Get, m.baiZhanService.GetPoint, obj.Id(),
		obj.LevelData(), lastJunXianLevelData, obj.LastPointChangeTime(), obj.Point(), obj.CombatMirrorFightAmount())
}

func (m *BaiZhanModule) processFunc(handlerName string, hc iface.HeroController, serverErrMsg pbutil.Buffer, f func(objs *bai_zhan_objs.BaiZhanObjs)) {
	if !m.baiZhanService.TimeOutFunc(func(objs *bai_zhan_objs.BaiZhanObjs) {
		f(objs)
	}) {
		logrus.Debugf("百战 %s，超时", handlerName)
		hc.Send(serverErrMsg)
		return
	}
}

func (m *BaiZhanModule) GmSetJunXian(junXian int64, hc iface.HeroController) {
	levelData := m.configDatas.JunXianLevelData().Must(u64.FromInt64(junXian))
	m.processFunc("gm变更军衔", hc, bai_zhan.ERR_BAI_ZHAN_CHALLENGE_FAIL_SERVER_ERROR, func(objs *bai_zhan_objs.BaiZhanObjs) {
		baiZhanObj := m.getOrCreateBaiZhanObj(objs, hc.Id())
		levelChange, historyLevelChanged := baiZhanObj.ResetLevelData(levelData, baiZhanObj.LastPointChangeTime())

		if levelChange {
			m.onJunXianLevelChanged(baiZhanObj, historyLevelChanged)
			m.mustPointRankList(levelData).AddOrUpdate(baiZhanObj)

			hc.Send(bai_zhan.RESET_S2C)
		}
	})
}

func (m *BaiZhanModule) Challenge(hc iface.HeroController) (err msg.ErrMsg) {
	m.processFunc("挑战", hc, bai_zhan.ERR_BAI_ZHAN_CHALLENGE_FAIL_SERVER_ERROR, func(objs *bai_zhan_objs.BaiZhanObjs) {
		baiZhanObj := m.getOrCreateBaiZhanObj(objs, hc.Id())

		ctime := m.timeService.CurrentTime()

		if baiZhanObj.ChallengeTimes() >= m.miscData.GetCanChallengeTimes(
			m.configDatas.MiscConfig().DailyResetTime.Duration(ctime),
			m.configDatas.MiscConfig().DailyResetDuration,
		) {
			logrus.Debugf("挑战百战千军，没有挑战次数了")
			err = bai_zhan.ErrBaiZhanChallengeFailNoChallengeTiems
			return
		}

		var attacker *shared_proto.CombatPlayerProto

		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			// 打架

			player, failType := hero.GenCombatPlayerProto(true, shared_proto.PveTroopType_DUNGEON, m.guildSnapshotGetter)
			switch failType {
			case entity.SUCCESS:
				break
			case entity.SERVER_ERROR:
				err = bai_zhan.ErrBaiZhanChallengeFailServerError
				return
			case entity.CAPTAIN_COUNT_NOT_ENOUGH:
				err = bai_zhan.ErrBaiZhanChallengeFailCaptainNotFull
				return
			default:
				logrus.Error("挑战百战千军，未处理的错误类型 %v", failType)
				err = bai_zhan.ErrBaiZhanChallengeFailServerError
				return
			}

			if player == nil {
				logrus.Errorf("挑战百战千军，构建Player出错")
				err = bai_zhan.ErrBaiZhanChallengeFailServerError
				return
			}

			attacker = player

			result.Ok()
		}) {
			// hasError 有错误
			if err == nil {
				err = bai_zhan.ErrBaiZhanChallengeFailServerError
			}
			return
		}

		if attacker == nil {
			// 里面有处理
			return
		}

		levelData := m.mustBaiZhanLevelData(baiZhanObj.LevelData())
		targetId, targetProto, isDefenderNpc := levelData.randomTarget(baiZhanObj.Id(), baiZhanObj.ChallengeTimes())

		tfctx := entity.NewTlogFightContext(operate_type.BattleBaiZhan, baiZhanObj.LevelData().Level,0, 0)
		response := m.fightService.SendFightRequest(tfctx, baiZhanObj.LevelData().CombatScene, hc.Id(), targetId, attacker, targetProto)
		if response == nil {
			logrus.Errorf("挑战百战千军，response==nil")
			err = bai_zhan.ErrBaiZhanChallengeFailServerError
			return
		}

		if response.ReturnCode != 0 {
			logrus.Errorf("挑战百战千军，战斗计算发生错误，%s", response.ReturnMsg)
			err = bai_zhan.ErrBaiZhanChallengeFailServerError
			return
		}

		newChallengeTimes := baiZhanObj.IncreChallengeTimes()
		newPoint := baiZhanObj.AddPoint(m.miscData.GetAddPoint(response.AttackerWin), m.timeService.CurrentTime())
		hc.Send(bai_zhan.NewS2cBaiZhanChallengeMsg(response.AttackerWin, u64.Int32(newChallengeTimes), response.Link, must.Marshal(response.AttackerShare), u64.Int32(newPoint), u64.Int32(baiZhanObj.HistoryMaxPoint(baiZhanObj.LevelData()))))

		if newPoint > 0 {
			// 更新积分榜
			m.mustPointRankList(baiZhanObj.LevelData()).AddOrUpdate(baiZhanObj)
		}

		m.addBaiZhanReplay(response, isDefenderNpc, attacker, targetProto)

		// 存镜像
		oldHasMirror := baiZhanObj.CombatMirror() != nil
		baiZhanObj.SetCombatMirror(attacker)
		if !oldHasMirror {
			// 此前没有加到任何镜像里面去
			levelData.addMirror(baiZhanObj)
		}

		baiZhanObj.IncreRecordVsn()
		if !isDefenderNpc {
			targetBaiZhan := objs.GetBaiZhanObj(targetId)
			if targetBaiZhan != nil {
				targetBaiZhan.IncreRecordVsn()
			} else {
				logrus.WithField("targetId", targetId).Debugln("刚刚跟目标打完，竟然目标的百战数据找不到了")
			}

			targetSender := m.worldService.GetUserSender(targetId)
			if targetSender != nil {
				targetSender.Send(bai_zhan.SELF_DEFENCE_RECORD_CHANGED_S2C)
			}
		}

		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_BAI_ZHAN)
			hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumBaiZhan)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BAI_ZHAN)
			result.Ok()
		})
	})

	return
}

// 添加百战回放记录
func (m *BaiZhanModule) addBaiZhanReplay(result *server_proto.CombatXResponseServerProto, isDefenderNpc bool, attacker, defenser *shared_proto.CombatPlayerProto) {
	replay := &shared_proto.BaiZhanReplayProto{}

	timeUnix := timeutil.Marshal64(m.timeService.CurrentTime())

	replay.Link = result.Link
	replay.IsDefenderNpc = isDefenderNpc
	replay.IsAttackerWin = result.GetAttackerWin()
	replay.Time = i64.Int32(timeUnix)

	replay.Attacker = m.buildCombatObj(attacker)
	replay.AttackerShare = result.AttackerShare
	replay.Defender = m.buildCombatObj(defenser)
	replay.DefenderShare = result.DefenserShare

	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		return m.dbService.InsertBaiZhanReplay(ctx, result.GetAttackerId(), result.GetDefenserId(), replay, isDefenderNpc, timeUnix)
	})
	if err != nil {
		logrus.WithError(err).Error("添加百战千军回放记录失败")
	}

}

func (m *BaiZhanModule) buildCombatObj(proto *shared_proto.CombatPlayerProto) *shared_proto.CombatObjProto {
	result := &shared_proto.CombatObjProto{}
	result.Hero = proto.Hero

	//result.Id = proto.GetId()
	//result.Name = proto.GetName()
	//result.Guild = proto.GetGuild()
	//result.Head = proto.GetHead()
	//result.GuildId = proto.GetGuildId()
	//result.GuildFlagName = proto.GetGuildFlagName()
	//result.Country = proto.CountryId
	//result.Male = proto.GetMale()
	//result.Level = proto.GetLevel()

	//tfa := data.NewTroopFightAmount()
	//if len(proto.Troops) > 0 {
	//	for _, t := range proto.Troops {
	//		if t.Captain != nil {
	//			tfa.AddInt32(t.Captain.FightAmount)
	//		}
	//	}
	//}
	//result.FightAmount = tfa.ToI32()
	result.FightAmount = proto.TotalFightAmount

	races := make([]shared_proto.Race, 5)
	for i, v := range proto.Troops {
		idx := i + 1
		if v.FightIndex > 0 {
			idx = int(v.FightIndex)
		}

		if idx <= len(races) {
			races[idx-1] = v.Captain.Race
		}
	}
	result.Race = races

	return result
}

func (m *BaiZhanModule) CollectSalary(hc iface.HeroController) (collectSuccess bool, err msg.ErrMsg) {
	obj := m.baiZhanService.GetBaiZhanObj(hc.Id())
	if obj != nil && obj.IsCollectDailySalary() {
		logrus.Debugln("领取俸禄时，玩家此前已经领取了")
		err = bai_zhan.ErrCollectSalaryFailSalaryCollect
		return
	}

	var heroName string
	var heroHead string
	var guildId int64
	m.processFunc("领取俸禄", hc, bai_zhan.ERR_COLLECT_SALARY_FAIL_SERVER_ERROR, func(objs *bai_zhan_objs.BaiZhanObjs) {
		baiZhanObj := m.getOrCreateBaiZhanObj(objs, hc.Id())

		if baiZhanObj.IsCollectDailySalary() {
			logrus.Debugln("领取俸禄时，玩家此前已经领取了")
			err = bai_zhan.ErrCollectSalaryFailSalaryCollect
			return
		}

		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			baiZhanObj.CollectDailySalary()
			hctx := heromodule.NewContext(m.dep, operate_type.BaiZhanCollectSalary)
			heromodule.AddPrize(hctx, hero, result, baiZhanObj.LevelData().DailySalary, m.timeService.CurrentTime())

			heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_COLLECT_BAI_ZHAN_SALARY, baiZhanObj.LevelData().Level)

			guildId = hero.GuildId()

			result.Changed()
			result.Ok()
		}) {
			// hasError 报错了
			err = bai_zhan.ErrCollectSalaryFailServerError
			return
		}

		collectSuccess = true

		toAddHufu := baiZhanObj.LevelData().DailyHufu
		if guildId > 0 && toAddHufu > 0 {
			m.dep.Guild().AddHufu(baiZhanObj.LevelData().DailyHufu, hc.Id(), guildId, heroName, heroHead)
		}

		hc.Send(bai_zhan.COLLECT_SALARY_S2C)
	})

	return
}

func (m *BaiZhanModule) CollectJunXianPrize(id int32, hc iface.HeroController) (collectSuccess bool, err msg.ErrMsg) {
	prizeData := m.configDatas.GetJunXianPrizeData(u64.FromInt32(id))
	if prizeData == nil {
		logrus.WithField("id", id).Debugln("军衔奖励数据没找到!")
		err = bai_zhan.ErrCollectJunXianPrizeFailPrizeNotFound
		return
	}

	obj := m.baiZhanService.GetBaiZhanObj(hc.Id())
	if obj != nil {
		// 可读数据做些检查
		if obj.HasCollectedJunXianPrize(prizeData) {
			logrus.WithField("id", id).Debugln("该军衔等级奖励已经领取了!")
			err = bai_zhan.ErrCollectJunXianPrizeFailPrizeCollected
			return
		}

		if obj.HistoryMaxJunXianLevelData().Level <= prizeData.LevelData.Level && obj.HistoryMaxPoint(prizeData.LevelData) < prizeData.Point {
			logrus.WithField("curPoint", obj.Point()).
				WithField("needPoint", prizeData.Point).
				Debugln("今日军衔积分不够，无法领取!")
			err = bai_zhan.ErrCollectJunXianPrizeFailPointTooLow
			return
		}

		if prizeData.Id != obj.LastCollectJunXianPrizeId()+1 {
			logrus.WithField("id", id).WithField("此前已经领取的id", obj.LastCollectJunXianPrizeId()).Debugln("领取军衔奖励数据id必须在上次的基础上+1!")
			err = bai_zhan.ErrCollectJunXianPrizeFailInvalidPrize
			return
		}
	}

	m.processFunc("领取军衔奖励", hc, bai_zhan.ERR_COLLECT_JUN_XIAN_PRIZE_FAIL_SERVER_ERROR, func(objs *bai_zhan_objs.BaiZhanObjs) {
		baiZhanObj := m.getOrCreateBaiZhanObj(objs, hc.Id())

		if baiZhanObj.HistoryMaxJunXianLevelData().Level <= prizeData.LevelData.Level && baiZhanObj.HistoryMaxPoint(prizeData.LevelData) < prizeData.Point {
			logrus.WithField("Point", baiZhanObj.HistoryMaxPoint(prizeData.LevelData)).
				WithField("needPoint", prizeData.Point).
				Debugln("军衔积分不够，无法领取!")
			err = bai_zhan.ErrCollectJunXianPrizeFailPointTooLow
			return
		}

		if baiZhanObj.HasCollectedJunXianPrize(prizeData) {
			logrus.WithField("id", id).Debugln("该军衔等级奖励已经领取了!")
			err = bai_zhan.ErrCollectJunXianPrizeFailPrizeCollected
			return
		}

		if prizeData.Id != baiZhanObj.LastCollectJunXianPrizeId()+1 {
			logrus.WithField("id", id).WithField("此前已经领取的id", baiZhanObj.LastCollectJunXianPrizeId()).Debugln("领取军衔奖励数据id必须在上次的基础上+1!")
			err = bai_zhan.ErrCollectJunXianPrizeFailInvalidPrize
			return
		}

		if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			// 标记领取奖励了
			baiZhanObj.CollectJunXianPrize(prizeData)

			hctx := heromodule.NewContext(m.dep, operate_type.BaiZhanCollectJunXianPrize)
			heromodule.AddPrize(hctx, hero, result, prizeData.Prize, m.timeService.CurrentTime())

			result.Changed()
			result.Ok()
		}) {
			// hasError 报错了
			err = bai_zhan.ErrCollectJunXianPrizeFailServerError
			return
		}

		collectSuccess = true

		hc.Send(prizeData.CollectedMsg)
	})

	return
}

func (m *BaiZhanModule) SelfRecord(vsn int32, hc iface.HeroController) {
	baiZhanObj := m.baiZhanService.GetBaiZhanObj(hc.Id())
	if baiZhanObj == nil {
		// 都没创建，肯定记录都没有啊
		hc.Send(bai_zhan.SELF_RECORD_S2C_NO_CHANGE)
		return
	}

	curVersion := baiZhanObj.RecordVsn()
	if curVersion == vsn {
		hc.Send(bai_zhan.SELF_RECORD_S2C_NO_CHANGE)
		return
	}

	var result [][]byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		result, err = m.dbService.LoadBaiZhanRecord(ctx, hc.Id(), m.miscData.MaxRecord)
		return
	})
	if err != nil {
		logrus.WithError(err).Debugln("获取自己的个人挑战记录报错了!")
		hc.Send(bai_zhan.ERR_SELF_RECORD_FAIL_SERVER_ERROR)
		return
	}

	hc.Send(bai_zhan.NewS2cSelfRecordMsg(curVersion, result))
}

func (m *BaiZhanModule) QueryBaiZhanInfo(hc iface.HeroController) {
	m.processFunc("查询百战信息", hc, bai_zhan.ERR_QUERY_BAI_ZHAN_INFO_FAIL_SERVER_ERROR, func(objs *bai_zhan_objs.BaiZhanObjs) {
		baiZhanObj := m.getOrCreateBaiZhanObj(objs, hc.Id())

		hc.Send(bai_zhan.NewS2cQueryBaiZhanInfoMarshalMsg(baiZhanObj.EncodeClient()))
	})
}

func (m *BaiZhanModule) ClearLastJunXian(hc iface.HeroController) {
	m.processFunc("清除Last军衔", hc, bai_zhan.ERR_QUERY_BAI_ZHAN_INFO_FAIL_SERVER_ERROR, func(objs *bai_zhan_objs.BaiZhanObjs) {
		if baiZhanObj := objs.GetBaiZhanObj(hc.Id()); baiZhanObj != nil {
			baiZhanObj.ClearLastJunXianLevelData()
		}
	})
}

var notInRankListMsg = bai_zhan.NewS2cRequestSelfRankMsg(0, int32(shared_proto.LevelChangeType_LEVEL_KEEP), 0, 0).Static()

func (m *BaiZhanModule) RequestSelfRank(hc iface.HeroController) {
	obj := m.baiZhanService.GetBaiZhanObj(hc.Id())
	if obj == nil {
		hc.Send(notInRankListMsg)
		return
	}

	rank := obj.Rank()
	if rank <= 0 {
		hc.Send(notInRankListMsg)
		return
	}

	levelUpMaxRank, levelUpNeedMinPoint, levelDownMinRank, levelKeepNeedPoint := m.mustPointRankList(obj.LevelData()).LevelUpAndDownRankAndPoint()

	levelChangeType := obj.LevelChangeType(levelUpMaxRank, levelDownMinRank)

	hc.Send(bai_zhan.NewS2cRequestSelfRankMsg(u64.Int32(rank), int32(levelChangeType), u64.Int32(levelUpNeedMinPoint), u64.Int32(levelKeepNeedPoint)))
}

func (m *BaiZhanModule) RequestRank(self bool, queryStartRank uint64, hc iface.HeroController) {
	baiZhanObj := m.baiZhanService.GetBaiZhanObj(hc.Id())

	var levelData *bai_zhan_data.JunXianLevelData

	if baiZhanObj == nil {
		levelData = m.configDatas.JunXianLevelData().MinKeyData
	} else {
		levelData = baiZhanObj.LevelData()

		if self {
			queryStartRank = u64.Max(1, u64.Sub(baiZhanObj.Rank(), m.miscData.ShowRankCount>>1))
		}
	}

	oc, err := m.mustPointRankList(levelData).RankCache(self, queryStartRank, m.miscData.ShowRankCount)
	if err != nil {
		logrus.WithError(err).Errorln("请求排行榜数据出错")
		hc.Send(bai_zhan.ERR_REQUEST_RANK_FAIL_SERVER_ERROR)
		return
	}

	hc.Send(oc)
}
