package tower

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/tower"
	"github.com/lightpaw/male7/module/rank/ranklist"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/concurrent"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"sort"
	"sync"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/gamelogs"
)

func NewTowerModule(dep iface.ServiceDep, fightService iface.FightXService, db iface.DbService,
	rankModule iface.RankModule, xuanyModule iface.XuanyuanModule, mailModule iface.MailModule) *TowerModule {
	m := &TowerModule{}
	m.dep = dep
	m.datas = dep.Datas()
	m.time = dep.Time()
	m.db = db
	m.heroSnapshotService = dep.HeroSnapshot()
	m.rankModule = rankModule
	m.xuanyModule = xuanyModule
	m.mailModule = mailModule
	m.fightService = fightService
	m.guildService = dep.Guild()

	m.floorMap = make(map[uint64]*tower_floor_replay)
	for _, data := range m.datas.TowerData().Array {
		r := &tower_floor_replay{}
		r.floor = data.Floor
		r.replayMsg = concurrent.NewBufferCache(r.buildReplayMsg)

		m.floorMap[r.floor] = r
	}

	var data []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		data, err = db.LoadKey(ctx, server_proto.Key_Tower)
		return
	})
	if err != nil {
		logrus.WithError(err).Panicf("加载DB中的千重楼模块数据失败")
	} else {
		if len(data) > 0 {
			proto := server_proto.TowerModuleProto{}
			if err := proto.Unmarshal(data); err != nil {
				logrus.WithError(err).Panicf("TowerModuleProto.Unmarshal, 千重楼模块数据失败")
			}

			for _, r := range proto.Replay {
				data := m.floorMap[u64.FromInt32(r.Floor)]
				if data != nil {
					data.add(r, m.datas.MiscConfig().TowerReplayCount)
				}

			}
		}
	}

	return m
}

//gogen:iface
type TowerModule struct {
	dep iface.ServiceDep

	datas iface.ConfigDatas

	time iface.TimeService

	db iface.DbService

	heroSnapshotService iface.HeroSnapshotService

	rankModule iface.RankModule

	xuanyModule iface.XuanyuanModule

	mailModule iface.MailModule

	fightService iface.FightXService

	guildService iface.GuildService

	floorMap map[uint64]*tower_floor_replay
}

func (t *TowerModule) Close() {

	proto := &server_proto.TowerModuleProto{}
	for _, r := range t.floorMap {
		if r.lowestReplay != nil {
			proto.Replay = append(proto.Replay, r.lowestReplay)
		}

		for _, v := range r.replay {
			if v != nil {
				proto.Replay = append(proto.Replay, v)
			}
		}
	}

	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		return t.db.SaveKey(ctx, server_proto.Key_Tower, must.Marshal(proto))
	})
	if err != nil {
		logrus.WithError(err).Errorf("保存千重楼数据出错")
	}
}

func (m *TowerModule) getFloorReplay(floor uint64) *tower_floor_replay {
	return m.floorMap[floor]
}

type tower_floor_replay struct {
	sync.Mutex

	floor uint64

	replay []*shared_proto.TowerReplayProto

	lowestReplay *shared_proto.TowerReplayProto

	// 缓存数据
	replayMsg concurrent.BufferCache
}

func (f *tower_floor_replay) getReplayMsg() pbutil.Buffer {
	return must.Msg(f.replayMsg.Get())
}

func (f *tower_floor_replay) buildReplayMsg() (pbutil.Buffer, error) {
	f.Lock()
	defer f.Unlock()

	proto := &shared_proto.TowerFloorReplayProto{}
	proto.Floor = u64.Int32(f.floor)
	proto.LowestReplay = f.lowestReplay
	proto.Replay = f.replay
	return tower.NewS2cListPassReplayMarshalMsg(u64.Int32(f.floor), proto).Static(), nil
}

type replaySlice []*shared_proto.TowerReplayProto

func (a replaySlice) Len() int           { return len(a) }
func (a replaySlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a replaySlice) Less(i, j int) bool { return a[i].Time < a[j].Time } // 时间小的在前面

func (f *tower_floor_replay) add(proto *shared_proto.TowerReplayProto, replayCount uint64) {
	f.Lock()
	defer f.Unlock()

	defer f.replayMsg.Clear()

	if f.lowestReplay == nil {
		f.lowestReplay = proto
		return
	}

	replaceLowest := proto.FightAmount < f.lowestReplay.FightAmount
	if replaceLowest {
		proto, f.lowestReplay = f.lowestReplay, proto
	}

	// 判断有没有满，如果没满，直接塞
	if replaceLowest || u64.FromInt(len(f.replay)) < replayCount {
		f.replay = append(f.replay, proto)
		sort.Sort(replaySlice(f.replay))
		if u64.FromInt(len(f.replay)) > replayCount {
			// 超出最大上限，干掉一个
			for i := uint64(0); i < replayCount; i++ {
				f.replay[i] = f.replay[i+1]
			}

			f.replay = f.replay[:replayCount]
		}
		return
	}

	// 满了，将时间少的那个替换出去
	n := len(f.replay) - 1
	for i := 0; i < n; i++ {
		f.replay[i] = f.replay[i+1]
	}
	f.replay[n] = proto
}

func (m *TowerModule) addFirstPassReplay(floor uint64, player *shared_proto.CombatPlayerProto, link string, share *shared_proto.CombatShareProto, firstPassPrize, prize *shared_proto.PrizeProto) {
	f := m.getFloorReplay(floor)
	if f == nil {
		return
	}

	proto := &shared_proto.TowerReplayProto{}
	proto.Floor = u64.Int32(floor)
	proto.Time = timeutil.Marshal32(m.time.CurrentTime())
	proto.Link = link
	proto.Share = share
	proto.FirstPassPrize = firstPassPrize
	proto.Prize = prize

	proto.Hero = player.Hero

	curIndex := int32(1)

	for _, t := range player.Troops {
		if curIndex < t.FightIndex && t.FightIndex <= 5 {
			for i := curIndex; i < t.FightIndex; i++ {
				proto.Race = append(proto.Race, 0)
			}
			curIndex = t.FightIndex
		}

		proto.Race = append(proto.Race, int32(t.Captain.Race))
		curIndex++
	}
	proto.FightAmount = player.TotalFightAmount

	f.add(proto, m.datas.MiscConfig().TowerReplayCount)
}

//gogen:iface c2s_challenge
func (m *TowerModule) ProcessChallenge(proto *tower.C2SChallengeProto, hc iface.HeroController) {

	nextFloorData := m.datas.TowerData().MinKeyData
	ctime := m.time.CurrentTime()

	var attacker *shared_proto.CombatPlayerProto
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroTower := hero.Tower()

		if u64.FromInt32(proto.Floor) != heroTower.CurrentFloor() {
			logrus.Debugf("挑战千重楼，当前楼层无效")
			result.Add(tower.ERR_CHALLENGE_FAIL_INVALID_FLOOR)
			return
		}

		if heroTower.ChallengeTimes() >= m.datas.MiscConfig().TowerChallengeMaxTimes && ctime.Before(heroTower.GetNextResetChallengeTime()) {
			logrus.Debugf("挑战千重楼，没有挑战次数了")
			result.Add(tower.ERR_CHALLENGE_FAIL_MAX_CHALLENGE_TIMES)
			return
		}

		if heroTower.CurrentFloor() > 0 {
			currentFloorData := m.datas.GetTowerData(heroTower.CurrentFloor())
			if currentFloorData == nil {
				// 这是神马情况
				logrus.Errorf("挑战千重楼，找不到当前楼层数据")
				result.Add(tower.ERR_CHALLENGE_FAIL_SERVER_ERROR)
				return
			}

			nextFloorData = currentFloorData.NextFloor()
			if nextFloorData == nil {
				logrus.Debugf("挑战千重楼，已经是最高层")
				result.Add(tower.ERR_CHALLENGE_FAIL_MAX_FLOOR)
				return
			}
		}

		// 打架

		player, failType := hero.GenCombatPlayerProto(true, shared_proto.PveTroopType_DUNGEON, m.guildService.GetSnapshot)
		switch failType {
		case entity.SUCCESS:
			break
		case entity.SERVER_ERROR:
			result.Add(tower.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		case entity.CAPTAIN_COUNT_NOT_ENOUGH:
			result.Add(tower.ERR_CHALLENGE_FAIL_CAPTAIN_NOT_FULL)
			return
		default:
			logrus.Errorf("挑战千重楼，未处理的错误类型 %v", failType)
			result.Add(tower.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		attacker = player
		if attacker == nil {
			logrus.Errorf("挑战千重楼，构建Player出错")
			result.Add(tower.ERR_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		result.Ok()
	})

	if attacker == nil {
		return
	}

	tfctx := entity.NewTlogFightContext(operate_type.BattleTower, nextFloorData.Floor, 0, 0)
	response := m.fightService.SendFightRequest(tfctx, nextFloorData.CombatScene, hc.Id(), nextFloorData.Monster.GetNpcId(), attacker, nextFloorData.Monster.GetPlayer())
	if response == nil {
		logrus.Errorf("挑战千重楼，response==nil")
		hc.Send(tower.ERR_CHALLENGE_FAIL_SERVER_ERROR)
		return
	}

	if response.ReturnCode != 0 {
		logrus.Errorf("挑战千重楼，战斗计算发生错误，%s", response.ReturnMsg)
		hc.Send(tower.ERR_CHALLENGE_FAIL_SERVER_ERROR)
		return
	}

	var addChallengerFunc func()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroTower := hero.Tower()

		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CHALLENGE_TOWER)

		gamelogs.TowerChallengeLog(hero.Pid(), hero.Sid(), hero.Id(), nextFloorData.Floor, response.AttackerWin)
		if !response.AttackerWin {
			// 挑战失败
			newTimes := heroTower.IncreseChallengeTimes(ctime)
			if newTimes == 1 {
				// 加挑战重置倒计时
				heroTower.SetNextResetChallengeTime(ctime.Add(m.datas.MiscConfig().TowerResetChallengeDuration))
			}

			result.Add(tower.NewS2cFailureChallengeMsg(u64.Int32(newTimes), timeutil.Marshal32(heroTower.GetNextResetChallengeTime()), response.Link, must.Marshal(response.AttackerShare)))

			toSendFunc := heromodule.TrySendMailFunc(m.mailModule, hero, shared_proto.HeroBoolType_BOOL_TOWER_FAIL,
				m.datas.MailHelp().FirstTowerFail, ctime)
			if toSendFunc != nil {
				toSendFunc()
			}
			// 每日首次挑战失败，发放千重楼协助礼包
			if giftData := m.datas.EventLimitGiftConfig().GetTowerHelpGift(); giftData != nil {
				if !hero.Promotion().IsDailyEventLimitGiftAppeared(giftData.Condition) {
					hero.Promotion().SetDailyEventLimitGiftAppeared(giftData.Condition)
					heromodule.ActivateEventLimitGift(hero, result, giftData, m.time.CurrentTime())
				}
			}
			return
		}

		// 获取当前楼层扫荡奖励
		prize := nextFloorData.GenChallengePrize()
		prizeProto := prize.Encode()
		challengePrize := must.Marshal(prizeProto)

		toAdd := prize

		var firstPassPrize []byte
		if nextFloorData.Floor > heroTower.HistoryMaxFloor() {
			firstPassPrize = nextFloorData.FirstPassPrizeBytes
			toAdd = resdata.AppendPrize(prize, nextFloorData.FirstPassPrize)

			// 首通回放
			m.addFirstPassReplay(nextFloorData.Floor, attacker, response.Link, response.AttackerShare, nextFloorData.FirstPassPrizeProto, prizeProto)
		}

		hctx := heromodule.NewContext(m.dep, operate_type.TowerChallenge)
		heromodule.AddPveTroopPrize(hctx, hero, result, toAdd, shared_proto.PveTroopType_DUNGEON, ctime)

		historyMaxFloorChanged := heroTower.IncreseCurrentFloor(ctime, m.datas.MiscConfig().TowerAutoKeepFloor)

		//heroTower.IncreaseAutoFloor()

		result.Add(tower.NewS2cChallengeMsg(response.Link, must.Marshal(response.AttackerShare), firstPassPrize, challengePrize, u64.Int32(heroTower.AutoMaxFloor())))

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_TOWER_FLOOR)

		if nextFloorData.UnlockSecretTower != nil {
			// 激活密室
			heroSecretTower := hero.SecretTower()

			openNewSecretTower := heroSecretTower.TryOpenAndGiveDefaultTimes(nextFloorData.UnlockSecretTower)
			if openNewSecretTower {
				result.Add(nextFloorData.UnlockSecretTower.UnlockMsg)
			}
		}

		heromodule.CheckFuncsOpened(hero, result)

		if historyMaxFloorChanged {
			// 历史最高层数变更了
			m.rankModule.AddOrUpdateRankObj(ranklist.NewTowerRankObj(m.heroSnapshotService.Get, hc.Id(), heroTower.HistoryMaxFloor(), heroTower.HistoryMaxFloorTime()))
			heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_TOWER, heroTower.HistoryMaxFloor())

			if d := hctx.BroadcastHelp().TowerFloor; d != nil {
				hctx.AddBroadcast(d, hero, result, 0, heroTower.HistoryMaxFloor(), func() *i18n.Fields {
					text := d.NewTextFields()
					text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
					text.WithFields(data.KeyNum, heroTower.HistoryMaxFloor())
					return text
				})
			}
			if giftData := m.datas.EventLimitGiftConfig().GetTowerGift(heroTower.HistoryMaxFloor()); giftData != nil {
				heromodule.ActivateEventLimitGift(hero, result, giftData, ctime)
			}
		}

		if !hero.Bools().Get(shared_proto.HeroBoolType_BOOL_XUAN_YUAN) {
			score := hero.Xuanyuan().GetScore()
			addChallengerFunc = func() {
				m.xuanyModule.AddChallenger(hc.Id(), attacker, score)
			}
		}

		result.Changed()
		result.Ok()
	})

	if addChallengerFunc != nil {
		addChallengerFunc()
	}

}

//gogen:iface c2s_auto_challenge
func (m *TowerModule) ProcessAutoChallenge(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroTower := hero.Tower()

		autoMaxFloorData := m.datas.GetTowerData(heroTower.AutoMaxFloor())
		if autoMaxFloorData == nil {
			logrus.Debugf("扫荡千重楼，没有楼层可以扫荡")
			result.Add(tower.ERR_AUTO_CHALLENGE_FAIL_AUTO_MAX)
			return
		}

		autoStartFloorData := m.datas.GetTowerData(heroTower.CurrentFloor() + 1)
		if autoStartFloorData == nil {
			logrus.Errorf("扫荡千重楼，扫荡起始楼层没有找到, %d", heroTower.CurrentFloor())
			result.Add(tower.ERR_AUTO_CHALLENGE_FAIL_AUTO_MAX)
			return
		}

		if autoStartFloorData.Floor > autoMaxFloorData.Floor {
			logrus.Debugf("扫荡千重楼，没有楼层可以扫荡")
			result.Add(tower.ERR_AUTO_CHALLENGE_FAIL_AUTO_MAX)
			return
		}

		failType := hero.CheckCanGenCombatPlayer(shared_proto.PveTroopType_DUNGEON)
		switch failType {
		case entity.SUCCESS:
			break
		case entity.SERVER_ERROR:
			logrus.Debugf("扫荡千重楼，服务器内部错误")
			result.Add(tower.ERR_AUTO_CHALLENGE_FAIL_SERVER_ERROR)
			return
		case entity.CAPTAIN_COUNT_NOT_ENOUGH:
			logrus.Debugf("扫荡千重楼，武将个数未满")
			result.Add(tower.ERR_AUTO_CHALLENGE_FAIL_CAPTAIN_NOT_FULL)
			return
		default:
			logrus.Debugf("扫荡千重楼，未处理的错误类型")
			result.Add(tower.ERR_AUTO_CHALLENGE_FAIL_SERVER_ERROR)
			return
		}

		// 获取当前楼层扫荡奖励
		b := resdata.NewPrizeBuilder()
		current := autoStartFloorData
		var prizeBytes [][]byte
		for i := autoStartFloorData.Floor; i <= autoMaxFloorData.Floor; i++ {
			toAdd := current.GenChallengePrize()
			b.Add(toAdd)
			prizeBytes = append(prizeBytes, must.Marshal(toAdd))

			if current.NextFloor() == nil {
				break
			}

			current = current.NextFloor()
		}

		hctx := heromodule.NewContext(m.dep, operate_type.TowerAutoChallenge)
		prize := b.Build()
		heromodule.AddPveTroopPrize(hctx, hero, result, prize, shared_proto.PveTroopType_DUNGEON, m.time.CurrentTime())

		heroTower.SetCurrentFloorToAuto()

		result.Add(tower.NewS2cAutoChallengeMsg(u64.Int32(heroTower.CurrentFloor()), prizeBytes))

		result.Changed()
		result.Ok()
	})
}

//gogen:iface c2s_collect_box
func (m *TowerModule) ProcessCollectBox(proto *tower.C2SCollectBoxProto, hc iface.HeroController) {
	// 改成无宝箱
	hc.Send(tower.ERR_COLLECT_BOX_FAIL_NOT_BOX)

	//towerData := m.datas.GetTowerData(u64.FromInt32(proto.BoxFloor))
	//if towerData == nil {
	//	logrus.Debugf("领取重楼宝箱，没找到重楼数据")
	//	hc.Send(tower.ERR_COLLECT_BOX_FAIL_NOT_BOX)
	//	return
	//}
	//
	//if towerData.BoxPrize == nil {
	//	logrus.Debugf("领取重楼宝箱，该楼层没有宝箱")
	//	hc.Send(tower.ERR_COLLECT_BOX_FAIL_NOT_BOX)
	//	return
	//}
	//
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	heroTower := hero.Tower()
	//
	//	if heroTower.HistoryMaxFloor() < towerData.Floor {
	//		logrus.Debugf("领取重楼宝箱，当前层还无法领取")
	//		result.Add(tower.ERR_COLLECT_BOX_FAIL_COLLECTED)
	//		return
	//	}
	//
	//	if heroTower.IsBoxCollected(towerData.Floor) {
	//		logrus.Debugf("领取重楼宝箱，已经领取过了")
	//		result.Add(tower.ERR_COLLECT_BOX_FAIL_COLLECTED)
	//		return
	//	}
	//
	//	heroTower.CollectBox(towerData.Floor)
	//
	//	heromodule.AddPrize(hero, result, towerData.BoxPrize, m.time.CurrentTime())
	//
	//	result.Add(towerData.CollectBoxPrizeMsg)
	//
	//	result.Changed()
	//	result.Ok()
	//})
}

//gogen:iface
func (m *TowerModule) ProcessListPassReplay(proto *tower.C2SListPassReplayProto, hc iface.HeroController) {

	floor := u64.FromInt32(proto.Floor)

	f := m.getFloorReplay(floor)
	if f == nil {
		logrus.Debugf("请求千重楼回放数据，无效的楼层")
		hc.Send(tower.ERR_LIST_PASS_REPLAY_FAIL_INVALID_FLOOR)
		return
	}

	hc.Send(f.getReplayMsg())
}
