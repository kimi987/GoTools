package realm

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/pushdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/sortkeys"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"strconv"
	"time"
	"math"
)

func (t *troop) assistArrivedEvent(r *Realm) {
	t.event.RemoveFromQueue()
	t.event = nil

	logrus.WithField("troopid", t.Id()).Debug("处理出征到达")

	if t.state != realmface.MovingToAssist {
		logrus.WithField("troopid", t.Id()).WithField("state", t.state).Error("realm竟然触发了assistArrivedEvent, 但是队伍状态不是MovingToAssist")
		return
	}

	defenserBase := t.targetBase
	if defenserBase == nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了assistArrivedEvent, 但是troop的targetBase为nil")
		return
	}

	ctime := r.services.timeService.CurrentTime()

	assembly := t.GetAssembly()
	if assembly != nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了assistArrivedEvent, 但是 assembly != nil")
		// 简单处理，回家
		t.backHome(r, ctime, true, true)
		return
	}

	// 看下有没有人正在打劫, 有的话就搞
	fighted := false
	for {
		robbing, hasMore := t.targetBase.getARobbingTroopWithLowestFightAmount()

		if robbing == nil {
			break
		}

		isAttackerWin := false // 这里attacker指的是援助方
		if assembly := robbing.GetAssembly(); assembly != nil {
			fightErr, isAttackerAlive, isDefenserAlive, _ := r.fightAssembly(robbing, defenserBase, []*troop{t}, nil, fightTypeAssistArrive)
			if fightErr {
				t.event = r.newEvent(r.services.timeService.CurrentTime().Add(1*time.Second), t.assistArrivedEvent) // 1秒后再来一次
				return
			}
			if isDefenserAlive {
				// 进攻方输了，看下进攻方是否有兵活着，如果有，回家
				if isAttackerAlive {
					// 集结到达，解散集结，返回自己家
					assembly.destroyAndTroopBackHome(r, defenserBase.Id(), defenserBase.BaseX(), defenserBase.BaseY(), ctime, true)
					return
				}
			}

			isAttackerWin = true
		} else {
			// 打一架
			ctx := &fightContext{}
			fightErr, attackerWin, _, _ := r.fight(ctx, t, robbing, fightTypeAssistArrive, hasMore)
			if fightErr {
				t.event = r.newEvent(r.services.timeService.CurrentTime().Add(1*time.Second), t.assistArrivedEvent) // 1秒后再来一次
				return
			}
			isAttackerWin = attackerWin
		}

		if !isAttackerWin {
			logrus.WithField("troopid", t.Id()).Debug("已到达目的地, 援助失败，回家")
			return
		}
		fighted = true
	}

	// 到这里是部队还活着, 看下驻扎人数是否满了
	if t.targetBase.getDefendingTroopCount() >= r.config().MaxAssist {
		logrus.WithField("troopid", t.Id()).Debug("帮忙的人来了, 但是没位子了, 回家")

		t.backHome(r, ctime, true, true)

		r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().RealmAssistArrivedBack.New().
					WithTroopIndex(entity.GetTroopIndex(t.Id()) + 1).
					JsonString())
		})
		return
	}

	// 援助到达

	r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
		if fighted {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().RealmAssistSuccess.New().
					WithTroopIndex(entity.GetTroopIndex(t.Id()) + 1).
					JsonString())
		} else {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().RealmAssistArrivedStay.New().
					WithTroopIndex(entity.GetTroopIndex(t.Id()) + 1).
					JsonString())
		}
	})

	// 驻扎下来吧
	logrus.WithField("troopid", t.Id()).Debug("已到达目的地, 驻扎防守")
	t.state = realmface.Defending
	//t.moveStartTime = time.Time{}
	//t.moveArriveTime = time.Time{}

	// 没有event
	//actionType, moveType := stateToActionMoveType(t.state)
	//t.proto.Action = actionType.Int32()
	//t.proto.MoveType = moveType.Int32()
	//t.proto.MoveStartTime, t.proto.MoveArrivedTime = 0, 0

	t.UpdateAssistDefendStartTime(ctime)

	t.onChanged()
	r.broadcastMaToCared(t, addTroopTypeUpdate, 0)

	// lock hero
	r.heroBaseFuncNotError(t.startingBase.Base(), func(hero *entity.Hero) (heroChanged bool) {
		hero.UpdateTroop(t, false)
		return true
	})
}

//侦察到达事件
func (t *troop) investigateArrivedEvent(r *Realm) {
	t.event.RemoveFromQueue()
	t.event = nil

	ctime := r.services.timeService.CurrentTime()
	//最后要执行返回操作
	defer t.backHome(r, ctime, true, true)

	logrus.WithField("troopid", t.Id()).Debug("处理侦察到达")

	if t.state != realmface.MovingToInvesigate {
		logrus.WithField("troopid", t.Id()).WithField("state", t.state).Error("realm竟然触发了investigateArrivedEvent, 但是队伍状态不是MovingToInvesigate")
		return
	}

	defenseBase := t.targetBase
	if defenseBase == nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了investigateArrivedEvent, 但是troop的targetBase为nil")
		return
	}

	// 免战，回家
	if t.startingBase.isHeroBaseType() {
		if ctime.Unix() < int64(defenseBase.MianDisappearTime()) {
			return
		}
	}
	var selfGuildId int64
	var attackerId int64
	// var selfBaseX, selfBaseY int
	var actProsperityCapcity uint64
	var attacker *shared_proto.ReportHeroProto

	if r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {

		selfGuildId = hero.GuildId()
		attackerId = hero.Id()
		// selfBaseX, selfBaseY = hero.BaseX(), hero.BaseY()
		actProsperityCapcity = hero.ProsperityCapcity()

		attacker = &shared_proto.ReportHeroProto{
			Id:                hero.IdBytes(),
			Name:              hero.Name(),
			Level:             int32(hero.Level()),
			Head:              hero.Head(),
			BaseRegion:        i64.Int32(hero.BaseRegion()),
			BaseX:             imath.Int32(hero.BaseX()),
			BaseY:             imath.Int32(hero.BaseY()),
			Prosperity:        u64.Int32(hero.Prosperity()),
			ProsperityCapcity: u64.Int32(hero.ProsperityCapcity()),
			Country:           u64.Int32(hero.CountryId()),
		}
		result.Ok()
	}) {
		return
	}

	report := &shared_proto.FightReportProto{}

	var defenserId int64

	if r.heroBaseFuncWithSend(defenseBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
		if selfGuildId != 0 && selfGuildId == hero.GuildId() {
			logrus.Debug("侦查城池，不能侦查盟友")
			result.Add(region.ERR_INVESTIGATE_INVADE_FAIL_SAME_GUILD)
			return
		}

		if ctime.Before(hero.GetMianDisappearTime()) {
			logrus.Debug("侦查城池，目标免战中")
			result.Add(region.ERR_INVESTIGATE_INVADE_FAIL_TARGET_MIAN)
			return
		}

		// 收税收
		heromodule.TryUpdateTax(hero, result, ctime, r.services.datas.MiscGenConfig().TaxDuration, r.services.datas.GetBuffEffectData)

		// defenserGuildId = hero.GuildId()
		defenserId = hero.Id()

		unsafe := hero.GetUnsafeResource()

		report.AttackerSide = true

		report.Defenser = &shared_proto.ReportHeroProto{
			Id:                hero.IdBytes(),
			Name:              hero.Name(),
			Level:             int32(hero.Level()),
			Head:              hero.Head(),
			BaseRegion:        i64.Int32(hero.BaseRegion()),
			BaseX:             imath.Int32(hero.BaseX()),
			BaseY:             imath.Int32(hero.BaseY()),
			Prosperity:        u64.Int32(hero.Prosperity()),
			ProsperityCapcity: u64.Int32(hero.ProsperityCapcity()),
			Country:           u64.Int32(hero.CountryId()),
		}

		report.Defenser.BaseLevel = u64.Int32(hero.BaseLevel())

		if bd := hero.Domestic().GetBuilding(shared_proto.BuildingType_CHENG_QIANG); bd != nil {
			report.Defenser.WallLevel = u64.Int32(bd.Level)
		}

		// 防守阵容
		if t := hero.GetHomeDefenser(); t != nil && !t.IsOutside() && t.HasSoldier() {
			for _, pos := range t.Pos() {
				c := pos.Captain()
				race := shared_proto.Race_InvalidRace
				if c != nil {
					race = c.Race().Race
				}
				report.Defenser.Race = append(report.Defenser.Race, race)
			}

			report.Defenser.TotalFightAmount = u64.Int32(t.CalDefenseFightAmount(hero))
		} else if copyDefenser := hero.GetCopyDefenser(); copyDefenser != nil {

			var captainFightAmounts []uint64
			for _, cis := range copyDefenser.GetCaptains() {
				race := shared_proto.Race_InvalidRace
				if cis != nil {
					c := hero.Military().Captain(cis.GetId())
					if c != nil {
						race = c.Race().Race

						if cis.GetSoldier() > 0 {
							captainFightAmounts = append(captainFightAmounts, cis.GetFightAmount())
						}
					}
				}
				report.Defenser.Race = append(report.Defenser.Race, race)
			}

			report.Defenser.TotalFightAmount = u64.Int32(data.TroopFightAmount(captainFightAmounts...))
		}

		// 衰减系数= arctan（（M守－M攻）/（M守＋M攻））/ π ＋1    π为圆周率
		defProsperityCapcity := hero.ProsperityCapcity()
		weakCoef := math.Atan2(
			u64.Sub2Float64(defProsperityCapcity, actProsperityCapcity),
			float64(defProsperityCapcity+actProsperityCapcity),
		)/math.Pi + 1
		//systemCoef := math.Min(m.datas.RegionConfig().RobberCoef, 1) // 系统损耗

		// 可掠夺资源 =（该资源当前储量-仓库保护值）*衰减系数*系统损耗
		coef := weakCoef //* systemCoef
		report.ShowPrize = &shared_proto.PrizeProto{}
		report.ShowPrize.Gold = i32.MultiF64(u64.Sub(unsafe.Gold(), hero.BuildingEffect().ProtectedCapcity()), coef)
		report.ShowPrize.Food = i32.MultiF64(u64.Sub(unsafe.Food(), hero.BuildingEffect().ProtectedCapcity()), coef)
		report.ShowPrize.Wood = i32.MultiF64(u64.Sub(unsafe.Wood(), hero.BuildingEffect().ProtectedCapcity()), coef)
		report.ShowPrize.Stone = i32.MultiF64(u64.Sub(unsafe.Stone(), hero.BuildingEffect().ProtectedCapcity()), coef)

		// 宝物
		topN := r.services.datas.RegionConfig().InvestigationBaowuCount
		if topN > 0 {

			top := sortkeys.NewU64TopN(topN)

			hero.Depot().RangeBaowu(func(id, count uint64) (toContinue bool) {
				if count <= 0 {
					return true
				}

				data := r.services.datas.GetBaowuData(id)
				if data == nil {
					return true
				}

				if data.CantRob {
					return true
				}

				top.Add(data.Level, data)

				return true
			})

			for _, kv := range top.SortDesc() {
				data := kv.V.(*resdata.BaowuData)

				count := hero.Depot().GetBaowuCount(data.Id)
				report.ShowPrize.BaowuId = append(report.ShowPrize.BaowuId, u64.Int32(data.Id))
				report.ShowPrize.BaowuCount = append(report.ShowPrize.BaowuCount, u64.Int32(count))
			}
		}
		resdata.SetPrizeProtoIsNotEmpty(report.ShowPrize)

		// 加被瞭望次数
		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_BeenInverstigation)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_BEEN_INVESTIGATION)

		result.Changed()
		result.Ok()
	}) {
		return
	}

	attackerFlagName := ""
	if attacker != nil {

		attackerFlagName = r.services.heroSnapshotService.GetFlagHeroName(attackerId)
	}

	defenserFlagName := ""
	if report.Defenser != nil {

		defenserFlagName = r.services.heroSnapshotService.GetFlagHeroName(defenserId)
	}

	var attackerMailId uint64
	if data := r.services.datas.MailHelp().ReportWatchAttacker; data != nil {
		mailProto := data.NewTextMail(shared_proto.MailType_MailInvestigation)
		mailProto.SubTitle = data.NewSubTitleFields().WithFields("attacker", attackerFlagName).WithFields("defenser", defenserFlagName).JsonString()
		mailProto.Text = data.NewTextFields().WithFields("attacker", attackerFlagName).WithFields("defenser", defenserFlagName).JsonString()
		mailProto.Report = report
		r.services.mail.SendReportMail(attackerId, mailProto, ctime)

		attackerMailId, _ = i64.FromBytesU64(mailProto.Id)
	}

	if data := r.services.datas.MailHelp().ReportWatchDefenser; data != nil {
		mailProto := data.NewTextMail(shared_proto.MailType_MailBeenInvestigation)
		mailProto.SubTitle = data.NewSubTitleFields().WithFields("attacker", attackerFlagName).WithFields("defenser", defenserFlagName).JsonString()
		mailProto.Text = data.NewTextFields().WithFields("attacker", attackerFlagName).WithFields("defenser", defenserFlagName).JsonString()
		mailProto.Report = &shared_proto.FightReportProto{
			AttackerSide: false,
			Attacker:     attacker,
		}
		r.services.mail.SendReportMail(defenserId, mailProto, ctime)
	}

	if r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
		if attackerMailId > 0 {
			expireTime := ctime.Add(r.services.datas.RegionConfig().InvestigationExpireDuration)
			hero.AddInvestigation(defenserId, expireTime, attackerMailId)
			if hero.GetInvestigationCount() >= r.services.datas.RegionConfig().InvestigationLimit {
				hero.ClearExpiredInvestigation(ctime, true)
			}

			result.Add(region.NewS2cGetPrevInvestigateMsg(idbytes.ToBytes(defenserId), i64.ToBytesU64(attackerMailId), timeutil.Marshal32(expireTime)))
		}

		result.Ok()
	}) {
		return
	}
}

func (t *troop) invasionArrivedEvent(r *Realm) {
	t.event.RemoveFromQueue()
	t.event = nil

	logrus.WithField("troopid", t.Id()).Debug("处理出征到达")

	if t.state != realmface.MovingToInvade {
		logrus.WithField("troopid", t.Id()).WithField("state", t.state).Error("realm竟然触发了invasionArrivedEvent, 但是队伍状态不是MovingToInvade")
		return
	}

	defenseBase := t.targetBase
	if defenseBase == nil {
		logrus.WithField("troopid", t.Id()).Error("realm竟然触发了invasionArrivedEvent, 但是troop的targetBase为nil")
		return
	}

	ctime := r.services.timeService.CurrentTime()

	// 免战，回家
	if t.startingBase.isHeroBaseType() {
		if ctime.Unix() < int64(defenseBase.MianDisappearTime()) {
			t.backHome(r, ctime, true, true)
			return
		}
	}

	if npcid.IsBaoZangNpcId(defenseBase.Id()) {
		// 探索殷墟到了就算
		t.walkAll(func(st *troop) (toContinue bool) {
			r.heroBaseFuncWithSend(st.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
				if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_INVASE_BAOZ) {
					heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_INVASE_BAOZ)
				}
			})
			return true
		})
	}

	// 看下有没有人防守, 有的话就搞
	targetDefenser := defenseBase.getBaseDefenser(r, t.targetBaseLevel)
	defenseBaseInfo := defenseBase.internalBase.getBaseInfoByLevel(t.targetBaseLevel)

	var defenseReport *shared_proto.FightReportProto
	var score int32
	if assembly := t.GetAssembly(); assembly != nil {
		// 集结战斗
		fightErr, isAttackerAlive, isDefenserAlive, report := r.fightAssembly(t, defenseBase, defenseBase.getAssisterTroops(), targetDefenser, fightTypeInvadeArrive)
		if fightErr {
			t.event = r.newEvent(r.services.timeService.CurrentTime().Add(1*time.Second), t.invasionArrivedEvent)
			return
		}

		if isDefenserAlive {
			// 进攻方输了，看下进攻方是否有兵活着，如果有，回家
			if isAttackerAlive {
				// 集结到达，解散集结，返回自己家
				assembly.destroyAndTroopBackHome(r, defenseBase.Id(), defenseBase.BaseX(), defenseBase.BaseY(), ctime, true)
			}
			// 前线补给礼包
			if giftData := r.services.datas.EventLimitGiftConfig().GetSupplyGift(); giftData != nil {
				assembly.self.walkAll(func(st *troop) (toContinue bool) {
					r.heroBaseFuncWithSend(st.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
						if !hero.Promotion().IsDailyEventLimitGiftAppeared(giftData.Condition) {
							hero.Promotion().SetDailyEventLimitGiftAppeared(giftData.Condition)
							heromodule.ActivateEventLimitGift(hero, result, giftData, ctime)
						}
						result.Ok()
					})
					return true
				})
			}
			return
		}

		defenseReport = report
	} else {
		// 普通战斗
		ctx := &fightContext{}

		for {
			defending, hasMore := defenseBase.getADefendingTroopWithLowestFightAmount(targetDefenser)

			if defending == nil {
				break
			}

			// 打一架
			fightErr, attackerWin, response, report := r.fight(ctx, t, defending, fightTypeInvadeArrive, hasMore)
			if fightErr {
				t.event = r.newEvent(r.services.timeService.CurrentTime().Add(1*time.Second), t.invasionArrivedEvent)
				return
			}

			if !attackerWin {
				logrus.WithField("troopid", t.Id()).Debug("已到达目的地, 出征失败，回家")

				// 怪物攻城，没破城，攻城怪物升级
				baseInfo := t.getStartingBaseInfo()
				if hateData := baseInfo.getHateData(); hateData != nil && defenseBase != nil {
					r.heroBaseFuncWithSend(defenseBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
						// 怪物攻城升级
						info := hero.GetOrCreateNpcTypeInfo(hateData.TypeData().Type)
						if info.GetRevengeLevel() < baseInfo.GetBaseLevel() {
							info.SetRevengeLevel(baseInfo.GetBaseLevel())
						}
						result.Ok()
					})
				}
				// 前线补给礼包
				if giftData := r.services.datas.EventLimitGiftConfig().GetSupplyGift(); giftData != nil {
					r.heroBaseFuncWithSend(t.StartingBase(), func(hero *entity.Hero, result herolock.LockResult) {
						if !hero.Promotion().IsDailyEventLimitGiftAppeared(giftData.Condition) {
							hero.Promotion().SetDailyEventLimitGiftAppeared(giftData.Condition)
							heromodule.ActivateEventLimitGift(hero, result, giftData, ctime)
						}
						result.Ok()
					})
				}
				return
			}

			if defending == targetDefenser {
				// 打完最后一架了
				defenseReport = report
				score = response.GetScore()
				break
			}
		}
	}

	// 怪物攻城，破城，怪物攻城降级
	baseInfo := t.getStartingBaseInfo()
	if hateData := baseInfo.getHateData(); hateData != nil && defenseBase != nil {
		r.heroBaseFuncWithSend(defenseBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
			// 怪物攻城降级
			info := hero.GetOrCreateNpcTypeInfo(hateData.TypeData().Type)

			if info.GetRevengeLevel() > 0 {
				info.SetRevengeLevel(u64.Sub(info.GetRevengeLevel(), 1))
			}

			result.Ok()
		})
	}

	var robResult robResult
	var robMaxDuration time.Duration

	isXiongNuDefense := npcid.IsXiongNuNpcId(defenseBase.Id())
	hasRobBaozTroop := false

	switch {
	case isXiongNuDefense:
		// nothing
	case npcid.IsBaoZangNpcId(defenseBase.Id()):
		// 更新宝藏怪物的防守部队

		if b := GetBaoZangBase(defenseBase); b != nil {
			if targetDefenser != nil {
				// 设置击杀者
				if b.heroType == 0 {
					// 玩家创建的不设置
					if home := GetHomeBase(t.startingBase); home != nil {
						if hero := r.services.heroSnapshotService.Get(t.startingBase.Id()); hero != nil {
							var keepEndTime int32
							if util.IsInRange(home.BaseX(), home.BaseY(), b.BaseX(), b.BaseY(), r.config().KeepBaozMaxDistance) {
								keepEndTime = timeutil.Marshal32(ctime.Add(r.config().KeepBaozDuration))
							}

							b.setKiller(hero.Id, hero.EncodeBasic4ClientBytes(), keepEndTime)
							if keepEndTime > 0 {
								home.AddKeepBaozMap(b.Id(), keepEndTime)
							}
						}
					}
				}

				b.ClearUpdateBaseInfoMsg()
				r.broadcastBaseInfoToCared(defenseBase, addBaseTypeUpdate, 0)

				defenseBase.updateRoBase()
			}

			// 同时只能有一只队伍在进行宝藏持续掠夺
			if t.startingBase.isHeroHomeBase() {
				for _, t := range t.startingBase.selfTroops {
					if t.State() == realmface.Robbing && npcid.IsBaoZangNpcId(t.originTargetId) {
						hasRobBaozTroop = true
						break
					}
				}
			}

		}

		if targetDefenser == nil {
			if hasRobBaozTroop {
				// 已经有一只队伍在进行持续掠夺，这只队伍就回家吧
				break
			}

			// 持续掠夺
			robMaxDuration = t.getRobDuration(r, defenseBaseInfo)
			break
		}
		fallthrough
	default:

		// 到这里部队还活着, 所有防守人都死光了
		// 先抢一次初始资源
		robResult = r.initialRob(defenseBase, t, ctime)

		if !hasRobBaozTroop {
			robMaxDuration = t.getRobDuration(r, defenseBaseInfo)
		}

	}

	tall := t.getAllTroops()

	// 回家
	if robberCount, robberMaxCount := defenseBase.getRobbingTroopCount(), defenseBaseInfo.BeenRobMaxCount(r);
		robResult.hasError || robResult.robberIsFull ||
			(robResult.targetWontLoseAnyProsperity && robResult.targetWontLoseAnyMoreResource) ||
			(robberMaxCount > 0 && robberCount >= robberMaxCount) ||
			robMaxDuration <= 0 || defenseBaseInfo.IsDestroyWhenLose() {
		// 回去算了
		logrus.WithField("troopid", t.Id()).
			WithField("err", robResult.hasError).
			WithField("robberIsFull", robResult.robberIsFull).
			WithField("wontLostAnyProsperity", robResult.targetWontLoseAnyProsperity).
			WithField("wontLostAnyMoreResource", robResult.targetWontLoseAnyMoreResource).
			WithField("robberCount", robberCount).
			WithField("robMaxDuration", robMaxDuration).
			WithField("destroyWhenLose", defenseBaseInfo.IsDestroyWhenLose()).
			Debug("抢了第一次之后, 不符合继续掠夺条件, 回家")

		// 解散后，里面没有集结的信息了
		t.backHome(r, ctime, true, false)

		jbase := GetJunTuanBase(defenseBase)
		if t.startingBase.isHeroBaseType() {
			for _, t := range tall {
				// lock hero
				var heroId int64
				var baseX, baseY int
				var baseLevel uint64
				var npcOffsets []cb.Cube
				r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
					heroId = hero.Id()
					baseX = hero.BaseX()
					baseY = hero.BaseY()
					baseLevel = hero.BaseLevel()
					hero.UpdateTroop(t, false)

					if defenseBaseInfo.IsDestroyWhenLose() {
						// 主城野怪Npc处理
						target := defenseBase
						if hero.GetHomeNpcBase(target.Id()) != nil {
							// 移除野怪Npc
							hero.RemoveHomeNpcBase(target.Id())
							r.updateHeroResourcePointConflict([]cb.Cube{cb.XYCube(target.BaseX(), target.BaseY())}, target.BaseX(), target.BaseY(), hero, result, ctime)
							npcOffset := hexagon.EvenOffsetBetween(hero.BaseX(), hero.BaseY(), target.BaseX(), target.BaseY())
							npcOffsets = append(npcOffsets, npcOffset)

							heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_KILL_HOME_NPC)
						}
					}

					// 胜利仇恨变化
					if hateData := defenseBaseInfo.getHateData(); hateData != nil {
						if hero.GetMultiLevelNpcPassLevel() < hateData.Level {
							hero.SetMultiLevelNpcPassLevel(hateData.Level)
							result.Add(region.NewS2cUpdateMultiLevelNpcPassLevelMsg(int32(hateData.TypeData().Type), u64.Int32(hateData.Level)))
						}
					}

					if npcid.IsMultiLevelMonsterNpcId(defenseBase.Id()) {
						// 扣野怪讨伐次数
						multiLevelNpcTimes := hero.GetMultiLevelNpcTimes()
						multiLevelNpcTimes.ReduceTimes(t.npcTimes, ctime, r.services.extraTimesService.MultiLevelNpcMaxTimes().TotalTimes())
						result.Add(region.NewS2cUpdateMultiLevelNpcTimesMsg(multiLevelNpcTimes.StartTimeUnix32(), nil))
						heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_KILL_MONSTER, t.targetBaseLevel)

						hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_WinMultiLevelMonster, 0)
						hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_WinMultiLevelMonster, taskdata.SubTypeMultiLevelNpc(0, defenseBaseInfo.GetBaseLevel()))

						if data := r.services.datas.GetRegionMultiLevelNpcData(npcid.GetNpcDataId(defenseBase.Id())); data != nil {
							hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_WinMultiLevelMonster, taskdata.SubTypeMultiLevelNpc(data.TypeData.Type, 0))
							hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_WinMultiLevelMonster, taskdata.SubTypeMultiLevelNpc(data.TypeData.Type, defenseBaseInfo.GetBaseLevel()))
						}
						heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_WIN_MULTI_LEVEL_MONSTER)

						// 系统广播:恭喜[color=#FFD633]{{self}}[/color]首次击破[color=#FFD633]{{num}}级{{name}}[/color]，缴获玄冥装备碎片，三军士气大振
						hctx := heromodule.NewContext(r.dep, operate_type.RealmInvasionNpc)
						if d := hctx.BroadcastHelp().NpcMonsterSucc; d != nil {
							hctx.AddBroadcast(d, hero, result, 0, defenseBaseInfo.GetBaseLevel(), func() *i18n.Fields {
								text := d.NewTextFields()
								text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
								text.WithNum(defenseBaseInfo.GetBaseLevel() * 10)
								text.WithName(defenseBase.internalBase.HeroName())
								return text
							})
						}
					}

					if jbase != nil {
						// 军团怪，扣次数
						hero.GetJunTuanNpcTimes().ReduceOneTimes(ctime, 0)
						result.Add(region.NewS2cUpdateJunTuanNpcTimesMsg(hero.GetJunTuanNpcTimes().StartTimeUnix32(), nil))

						hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_KillJunTuanTimes, 0)
						hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_KillJunTuanTimes, jbase.data.Level)
						heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_KILL_JUN_TUAN)
					}

					result.Changed()
					result.Ok()
				})

				// 胜利飘字
				r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
					return misc.NewS2cScreenShowWordsMsg(
						r.getTextHelp().RealmInvadeSuccessBack.New().
							WithTroopIndex(entity.GetTroopIndex(t.Id()) + 1).
							JsonString())
				})

				// 农场野怪冲突
				ctime := r.services.timeService.CurrentTime()
				r.updateFarmWithNpc(heroId, baseLevel, baseX, baseY, npcOffsets, false, ctime)
			}
		}

		// 第一次抢完就回家 邮件
		r.sendInvadeSuccessMail(t, tall, defenseBase, targetDefenser, score,
			defenseReport, ctime, true)

		// 摸一下就爆了（匈奴大营，军团怪，摸一下就爆）
		if defenseBaseInfo.IsDestroyWhenLose() || isXiongNuDefense {
			r.removeRealmBaseNoReason(defenseBase, removeBaseTypeBroken, ctime)
		}

		gid := t.startingBase.GuildId()
		if gid != 0 {
			if isXiongNuDefense { // 匈奴
				if data := r.services.datas.GuildLogHelp().XiongNuBaseDestroy; data != nil {
					if hero := r.services.heroSnapshotService.Get(t.startingBase.Id()); hero != nil {
						if hero.GuildId == gid {
							proto := data.NewHeroLogProto(ctime, hero.IdBytes, hero.Head)
							proto.Text = data.Text.New().WithHeroName(hero.Name).JsonString()
							proto.FightX = int32(defenseBase.BaseX())
							proto.FightY = int32(defenseBase.BaseY())

							r.services.guildService.AddLog(hero.GuildId, proto)
						}
					}
				}
			}
			if jbase != nil { // 军团怪
				// 更新联盟周任务进度
				data := r.services.datas.GetGuildTaskData(u64.FromInt32(int32(server_proto.GuildTaskType_QuanRong)))
				r.services.guildService.AddGuildTaskProgress(gid, data, 1)
				// 触发联盟声望事件
				r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
					// 这个hero是集结发起者
					heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_ASSEMBLY_KILL_MONSTER, jbase.data.Npc.Id)
					result.Ok()
				})
			}
		}
		return
	}

	// 持续掠夺

	// 胜利飘字
	r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
		return misc.NewS2cScreenShowWordsMsg(
			r.getTextHelp().RealmInvadeSuccess.New().
				WithTroopIndex(entity.GetTroopIndex(t.Id()) + 1).
				JsonString())
	})

	if defenseBase.isHeroHomeBase() {
		attackerName := t.getStartingBaseFlagHeroName(r)

		r.services.pushService.PushFunc(shared_proto.SettingType_ST_BEEN_ROBBING, defenseBase.Id(), func(d *pushdata.PushData) (title, content string) {
			return d.Title, d.ReplaceContent("{{attacker}}", attackerName)
		})

		if snapshot := r.services.guildService.GetSnapshot(defenseBase.GuildId()); snapshot != nil {
			friendName := r.toBaseFlagHeroName(defenseBase, defenseBase.BaseLevel())

			r.services.pushService.MultiPushFunc(shared_proto.SettingType_ST_GUILD_MEMBER_BEEN_ROBBING, snapshot.UserMemberIds, defenseBase.Id(), func(d *pushdata.PushData) (title, content string) {
				return d.Title, d.ReplaceContent("{{attacker}}", attackerName, "{{friend}}", friendName)
			})
		}
	}

	// 开始持续掠夺
	logrus.WithField("troopid", t.Id()).WithField("target", defenseBase).Debug("开始持续掠夺")
	t.state = realmface.Robbing
	if t.targetBase != nil {
		t.targetBase.remindAttackOrRobCountChanged(r)
	}
	//t.moveStartTime = time.Time{}
	//t.moveArriveTime = time.Time{}

	// 设置掠夺时间上限, 还要修改proto, 保存在数据库
	robbingEndTime := ctime.Add(robMaxDuration)
	nextTickTime := t.InitRobbing(robbingEndTime, r, defenseBaseInfo, ctime)
	t.event = r.newEvent(nextTickTime, t.doRobEvent)
	//t.proto = t.doEncode(r) // TODO 只修改proto中对应的值
	//t.refreshMsg()

	t.onChanged()
	r.broadcastMaToCared(t, addTroopTypeUpdate, 0)

	// 持续掠夺 邮件
	r.sendInvadeSuccessMail(t, tall, defenseBase, targetDefenser, score,
		defenseReport, ctime, false)

	// lock hero
	var isRemoveWhiteFlag bool
	r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
		hero.UpdateTroop(t, false)

		// 移除白旗
		if hero.GetWhiteFlagGuildId() != 0 {
			isRemoveWhiteFlag = hero.GetWhiteFlagGuildId() == defenseBase.GuildId()
			if isRemoveWhiteFlag {
				hero.RemoveWhiteFlag()
			}
		}

		// 胜利仇恨变化
		if hateData := defenseBaseInfo.getHateData(); hateData != nil {
			if hero.GetMultiLevelNpcPassLevel() < hateData.Level {
				hero.SetMultiLevelNpcPassLevel(hateData.Level)
				result.Add(region.NewS2cUpdateMultiLevelNpcPassLevelMsg(int32(hateData.TypeData().Type), u64.Int32(hateData.Level)))
			}
		}

		if npcid.IsMultiLevelMonsterNpcId(defenseBase.Id()) {
			// 扣野怪讨伐次数
			multiLevelNpcTimes := hero.GetMultiLevelNpcTimes()
			multiLevelNpcTimes.ReduceTimes(t.npcTimes, ctime, r.services.extraTimesService.MultiLevelNpcMaxTimes().TotalTimes())
			result.Add(region.NewS2cUpdateMultiLevelNpcTimesMsg(multiLevelNpcTimes.StartTimeUnix32(), nil))
			heromodule.OnHeroEventWithSubType(hero, result, shared_proto.HeroEvent_HERO_EVENT_KILL_MONSTER, t.targetBaseLevel)

			hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_WinMultiLevelMonster, 0)
			hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_WinMultiLevelMonster, taskdata.SubTypeMultiLevelNpc(0, defenseBaseInfo.GetBaseLevel()))

			if data := r.services.datas.GetRegionMultiLevelNpcData(npcid.GetNpcDataId(defenseBase.Id())); data != nil {
				hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_WinMultiLevelMonster, taskdata.SubTypeMultiLevelNpc(data.TypeData.Type, 0))
				hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_WinMultiLevelMonster, taskdata.SubTypeMultiLevelNpc(data.TypeData.Type, defenseBaseInfo.GetBaseLevel()))
			}
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_WIN_MULTI_LEVEL_MONSTER)

		}

		result.Changed()
		result.Ok()
	})

	if isRemoveWhiteFlag {
		r.removeHeroWhiteFlagIfSame(t.startingBase, defenseBase.GuildId(), true)
	}

	if assembly := t.GetAssembly(); assembly != nil {
		assembly.broadcastChanged(r)
	}

}

// 消灭匈奴npc
func (r *Realm) onWipeOutXiongNuNpcTroop(baseWithData *baseWithData, xiongNuNpcBase *xiongNuNpcBase, toReduceMorale uint64) {
	if toReduceMorale > 0 {
		// 减少士兵数量
		reducePer := float64(r.services.datas.ResistXiongNuMisc().OneMoraleReduceSoldierPer*toReduceMorale) / 10000

		for _, troop := range baseWithData.selfTroops {
			if troop.State() != realmface.Defending {
				continue
			}

			// 协助
			hasAlive := false

			for _, c := range troop.Captains() {
				if c.Proto().Soldier <= 0 {
					continue
				}

				toSet := c.Proto().Soldier - int32(float64(c.Proto().TotalSoldier)*reducePer)
				if toSet < 0 {
					toSet = 0
				}

				c.GetAndSetSoldierCount(uint64(toSet))

				if toSet > 0 {
					hasAlive = true
				}
			}

			// 看看是不是所有的都死了
			troop.onChanged()

			if !hasAlive {
				// 都死了，移除掉
				troop.removeWithoutReturnCaptain(r)
				r.broadcastRemoveMaToCared(troop)
				troop.clearMsgs()
			}
		}

		if troop := xiongNuNpcBase.Defenser(); troop != nil {
			// 防守
			for _, c := range troop.Captains() {
				if c.Proto().Soldier <= 0 {
					continue
				}

				toSet := c.Proto().Soldier - int32(float64(c.Proto().TotalSoldier)*reducePer)
				if toSet < 0 {
					toSet = 0
				}

				c.GetAndSetSoldierCount(uint64(toSet))
			}

			// 看看是不是所有的都死了
			troop.onChanged()
		}
	}

	// 干掉这个主城
	defenser := xiongNuNpcBase.Defenser()
	if defenser != nil && defenser.AliveSoldier() > 0 {
		return
	}

	r.removeRealmBaseNoReason(baseWithData, removeBaseTypeBroken, r.services.timeService.CurrentTime())
}

func (t *troop) returnedToBaseEvent(r *Realm) {
	t.event.RemoveFromQueue()
	t.event = nil

	logrus.WithField("troopid", t.Id()).Debug("处理回家到达")

	if t.assembly != nil && t.assembly.self == t {
		ctime := r.services.timeService.CurrentTime()
		t.assembly.destroyAndTroopBackHome(r, t.StartingBase().Id(), t.startingBase.BaseX(), t.startingBase.BaseY(), ctime, true)
	} else {
		t.doReturnedToBase(r, true)
	}
}

func (t *troop) doReturnedToBase(r *Realm, isBackHome bool) {
	// 回到家处理
	if t.startingBase.isHeroBaseType() {
		hasError := r.heroTroopReturnedToBase(t, isBackHome)
		if hasError {
			// 有错误，下一秒再来一次
			t.event = r.newEvent(r.services.timeService.CurrentTime().Add(1*time.Second), t.returnedToBaseEvent)
			return
		}
	}
	t.removeWithoutReturnCaptain(r)
	r.broadcastRemoveMaToCared(t)
	t.clearMsgs()
}

func (r *Realm) heroTroopReturnedToBase(t realmface.Troop, isBackHome bool) (hasError bool) {

	troopIndex := entity.GetTroopIndex(t.Id()) + 1

	hasError = r.heroBaseFuncWithSend(t.StartingBase(), func(hero *entity.Hero, result herolock.LockResult) {
		var showText *data.Text
		if isBackHome {
			showText = r.services.datas.TextHelp().RealmTroopBacked
		}

		r.doRemoveHeroTroop(hero, result, t, showText)
		result.Ok()
		return
	})

	if hasError {
		return
	}

	if isBackHome {
		r.services.pushService.PushFunc(shared_proto.SettingType_ST_TROOP_BACK_HOME, t.StartingBase().Id(), func(d *pushdata.PushData) (title, content string) {
			return d.Title, d.ReplaceContent("{{troop_index}}", strconv.FormatUint(troopIndex, 10))
		})
	}

	return

}

func (r *Realm) doRemoveHeroTroop(hero *entity.Hero, result herolock.LockResult, t realmface.Troop, showText *data.Text) {

	troopIndex := entity.GetTroopIndex(t.Id()) + 1

	// 删掉部队
	hero.RemoveTroop(t, true)

	hctx := heromodule.NewContext(r.dep, operate_type.RealmReturnedToBase)
	// 武将回家升级
	for _, c := range t.Captains() {
		heromodule.TryUpgradeCaptain(hctx, hero, result, c.Id(), r.services.timeService.CurrentTime())
	}

	// 武将到家消息
	result.AddFunc(func() pbutil.Buffer {
		return region.NewS2cUpdateSelfTroopsOutsideMsg(u64.Int32(troopIndex), false)
	})

	if showText != nil {
		result.AddFunc(func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(showText.New().
				WithTroopIndex(troopIndex).
				JsonString())
		})
	}

	result.Changed()

}

func (t *troop) doRobEvent(r *Realm) {

	// 持续掠夺
	//logrus.WithField("troopid", t.Id()).Debug("处理持续掠夺")

	// 更新下一次tick时间
	ctime := r.services.timeService.CurrentTime()

	targetBase := t.targetBase
	if targetBase == nil {
		logrus.WithField("troopid", t.Id()).Error("处理持续掠夺，targetBase == nil")

		r.trySendTroopDoneMail(t, r.getTextHelp().MRDRFinished4a.Text, r.getTextHelp().MRDRFinished4d.Text, ctime)
		t.backHome(r, ctime, false, true)
		return
	}

	targetBaseInfo := targetBase.internalBase.getBaseInfoByLevel(t.targetBaseLevel)
	// 加一次tick时间

	robbingStartTime := t.robbingEndTime.Add(-t.getRobDuration(r, targetBaseInfo))

	nextTickTime := t.robbingEndTime

	addPrize := false
	if !t.startingBase.isNpcBase() {
		if d := targetBaseInfo.BeenRobTickDuration(r); d > 0 {
			addPrize = t.nextAddPrizeTime.Before(ctime)
			if addPrize {
				t.nextAddPrizeTime = timeutil.NextTickTime(robbingStartTime, ctime, d)
			}
			nextTickTime = timeutil.Min(nextTickTime, t.nextAddPrizeTime)
		}
	}

	reduceProsperity := false
	if t.mmd == nil || t.mmd.GetSpec() != monsterdata.InvasionTask {
		// 任务怪不会打掉繁荣度

		if d := targetBaseInfo.BeenRobLostProsperityDuration(r); d > 0 {
			reduceProsperity = t.nextReduceProsperityTime.Before(ctime)
			if reduceProsperity {
				// 扣对方繁荣度
				t.nextReduceProsperityTime = timeutil.NextTickTime(robbingStartTime, ctime, d)
			}
			nextTickTime = timeutil.Min(nextTickTime, t.nextReduceProsperityTime)
		}
	}

	robBaowu := false
	if t.startingBase.isHeroHomeBase() && targetBase.isHeroHomeBase() {
		if d := r.services.datas.RegionConfig().RobBaowuTickDuration; d > 0 {
			// 人打人才会抢宝物
			robBaowu = t.nextRobBaowuTime.Before(ctime)
			if robBaowu {
				t.nextRobBaowuTime = timeutil.NextTickTime(robbingStartTime, ctime, d)
			}
			nextTickTime = timeutil.Min(nextTickTime, t.nextRobBaowuTime)
		}
	}

	addHate := false
	if !t.startingBase.isNpcBase() {
		if hateData := targetBaseInfo.getHateData(); hateData != nil && hateData.HateTickDuration > 0 {
			addHate = t.nextAddHateTime.Before(ctime)
			if addHate {
				// 加仇恨
				t.nextAddHateTime = timeutil.NextTickTime(robbingStartTime, ctime, hateData.HateTickDuration)
			}
			nextTickTime = timeutil.Min(nextTickTime, t.nextAddHateTime)
		}
	}

	t.event.UpdateTime(nextTickTime)

	robResult := r.recurringRob(targetBase, t, addPrize, reduceProsperity, robBaowu, addHate, ctime)
	if robResult.targetDestroyed {
		// 对方基地爆了，里面处理完了
		if npcid.IsBaoZangNpcId(targetBase.Id()) {
			// 尾刀击破广播
			r.heroBaseFuncWithSend(t.startingBase.Base(), func(hero *entity.Hero, result herolock.LockResult) {
				hctx := heromodule.NewContext(r.dep, operate_type.RealmRobNpc)
				if d := hctx.BroadcastHelp().BaoZangDestory; d != nil {
					hctx.AddBroadcast(d, hero, result, 0, targetBase.BaseLevel(), func() *i18n.Fields {
						text := d.NewTextFields()
						text.WithClickHeroFields(data.KeySelf, hero.Name(), hero.Id())
						text.WithFields(data.KeyNum, targetBase.BaseLevel())
						return text
					})
				}
			})
		}
		return
	}

	if robResult.hasError || robResult.robberIsFull ||
		(robResult.targetWontLoseAnyProsperity && robResult.targetWontLoseAnyMoreResource) ||
		t.robbingEndTime.Before(ctime) {
		// 回去算了
		logrus.WithField("troopid", t.Id()).
			WithField("err", robResult.hasError).
			WithField("robberIsFull", robResult.robberIsFull).
			WithField("wontLostAnyProsperity", robResult.targetWontLoseAnyProsperity).
			WithField("wontLostAnyMoreResource", robResult.targetWontLoseAnyMoreResource).
			WithField("robDuration", t.robbingEndTime.Before(ctime)).
			Debug("持续掠夺, 不符合继续掠夺条件, 回家")

		r.trySetWhiteFlag(t)
		r.trySendTroopDoneMail(t, r.getTextHelp().MRDRFinished4a.Text, r.getTextHelp().MRDRFinished4d.Text, ctime)
		t.backHome(r, ctime, true, true)

		// 破坏结束飘字
		r.services.world.SendFunc(t.startingBase.Id(), func() pbutil.Buffer {
			return misc.NewS2cScreenShowWordsMsg(
				r.getTextHelp().RealmInvadeFinished.New().
					WithTroopIndex(entity.GetTroopIndex(t.Id()) + 1).
					JsonString())
		})
		return
	}

	//if !robResult.targetWontLoseAnyMoreResource {
	//	// 抢到了东西
	//	t.proto = t.doEncode(r) // TODO 只修改proto中对应的值
	//	t.refreshMsg()
	//	r.broadcastMaToCared(t, 0)
	//}

	// lock hero
	r.heroBaseFuncNotError(t.startingBase.Base(), func(hero *entity.Hero) (heroChanged bool) {
		hero.UpdateTroop(t, false)
		return true
	})
}
