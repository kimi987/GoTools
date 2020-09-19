package military

import (
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/util/random"
)

var maxRecruitSoldierMsg = military.NewS2cGetMaxRecruitSoldierMsg(0).Static()

//gogen:iface c2s_get_max_recruit_soldier
func (m *MilitaryModule) ProcessGetMaxRecruitSoldier_deprecated(hc iface.HeroController) {
	hc.Send(maxRecruitSoldierMsg)

	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	heroMilitary := hero.Military()
	//
	//	// 新兵数量
	//	ctime := m.timeService.CurrentTime()
	//	newSoldierCount := heroMilitary.NewSoldierCount(ctime)
	//
	//	var amount uint64
	//	if newSoldierCount > 0 {
	//		// 剩余可招募士兵数
	//		maxRecruitCount := u64.Sub(heroMilitary.SoldierCapcity(), heroMilitary.FreeSoldier()+heroMilitary.FightSoldier())
	//		if maxRecruitCount <= 0 {
	//			return
	//		}
	//
	//		costCount := heromodule.GetCostCount(hero, heroMilitary.SoldierLevelData().RecruitCost)
	//		if costCount <= 0 {
	//			return
	//		}
	//
	//		amount = u64.Min(newSoldierCount, u64.Min(maxRecruitCount, costCount))
	//	}
	//
	//	result.Add(military.NewS2cGetMaxRecruitSoldierMsg(u64.Int32(amount)))
	//
	//	result.Ok()
	//	return
	//})
}

//gogen:iface
func (m *MilitaryModule) ProcessC2SRecruitSoldierMsg_deprecated(proto *military.C2SRecruitSoldierProto, hc iface.HeroController) {

	hc.Send(military.ERR_RECRUIT_SOLDIER_FAIL_INVALID_COUNT)
	return

	//recruitCount := u64.FromInt32(proto.GetCount())
	//if recruitCount <= 0 {
	//	logrus.Debug("招募新兵，新兵数<=0")
	//	hc.Send(military.ERR_RECRUIT_SOLDIER_FAIL_INVALID_COUNT)
	//	return
	//}
	//
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	// 先设置为错误，后面再设回来
	//
	//	heroMilitary := hero.Military()
	//
	//	// 新兵数量
	//	ctime := m.timeService.CurrentTime()
	//	newSoldierCount := heroMilitary.NewSoldierCount(ctime)
	//
	//	if newSoldierCount < recruitCount {
	//		logrus.Debug("招募新兵，新兵没这么多")
	//		result.Add(military.ERR_RECRUIT_SOLDIER_FAIL_SOLDIER_NOT_ENOUGH)
	//		return
	//	}
	//
	//	// 新兵招募消耗资源
	//	soldierLevelData := heroMilitary.SoldierLevelData()
	//	totalCost := soldierLevelData.RecruitCost.Multiple(recruitCount)
	//
	//	// 可招募士兵数
	//
	//	maxRecruitCount := u64.Sub(heroMilitary.SoldierCapcity(), heroMilitary.FreeSoldier()+heroMilitary.FightSoldier())
	//	if recruitCount > maxRecruitCount {
	//		logrus.Debug("招募新兵，超出士兵上限")
	//		result.Add(military.ERR_RECRUIT_SOLDIER_FAIL_SOLDIER_CAPCITY_OVERFLOW)
	//		return
	//	}
	//
	//	if !heromodule.TryReduceCost(hero, result, totalCost) {
	//		logrus.Debug("招募新兵，资源不足")
	//		result.Add(military.ERR_RECRUIT_SOLDIER_FAIL_RESOURCE_NOT_ENOUGH)
	//		return
	//	}
	//	result.Changed()
	//
	//	// 开始加人
	//	currentNewSoldierCount := u64.Sub(newSoldierCount, recruitCount)
	//	heroMilitary.SetNewSoldierCount(currentNewSoldierCount, ctime)
	//	newFreeSoldier := heroMilitary.AddFreeSoldier(recruitCount)
	//
	//	result.Add(military.NewS2cRecruitSoldierMsg(u64.Int32(currentNewSoldierCount), u64.Int32(newFreeSoldier)))
	//
	//	// 历史累计数量
	//	hero.HistoryAmount().Increase(server_proto.HistoryAmountType_RecruitSoldierCount, recruitCount)
	//
	//	// 更新任务进度
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_RECRUIT_SOLDIER_COUNT)
	//	heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_RECRUIT_SOLDIER)
	//
	//	result.Ok()
	//	result.Changed()
	//	return
	//})
}

//gogen:iface c2s_recruit_soldier_v2
func (m *MilitaryModule) ProcessC2SRecruitSoldierV2Msg_deprecated(proto *military.C2SRecruitSoldierV2Proto, hc iface.HeroController) {
	hc.Send(military.ERR_RECRUIT_SOLDIER_V2_FAIL_NO_TIMES)

	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	heroMilitary := hero.Military()
	//
	//	building := hero.Domestic().GetBuilding(shared_proto.BuildingType_JUN_YING)
	//	if building == nil {
	//		logrus.Debug("招募新兵V2，没有军营")
	//		result.Add(military.ERR_RECRUIT_SOLDIER_V2_FAIL_NO_JUN_YING)
	//		return
	//	}
	//
	//	// 计算出需要恢复的士兵数量
	//	needRecoverSolidier := heroMilitary.NeedRecoverSoldier()
	//
	//	if needRecoverSolidier <= 0 && heroMilitary.FreeSoldier() >= heroMilitary.SoldierCapcity() {
	//		logrus.Debug("招募新兵V2，超出士兵上限")
	//		result.Add(military.ERR_RECRUIT_SOLDIER_V2_FAIL_SOLDIER_CAPCITY_OVERFLOW)
	//		return
	//	}
	//
	//	if heroMilitary.FreeSoldier() < heroMilitary.SoldierCapcity()+needRecoverSolidier {
	//		// 士兵数量不足以自动征兵后还超出上限
	//
	//		ctime := m.timeService.CurrentTime()
	//		recruitTimes := heroMilitary.RecruitTimes()
	//		times := recruitTimes.Times(ctime)
	//		if times <= 0 {
	//			logrus.Debug("招募新兵V2，没有次数了")
	//			result.Add(military.ERR_RECRUIT_SOLDIER_V2_FAIL_NO_TIMES)
	//			return
	//		}
	//
	//		var reduceTimes uint64
	//
	//		// 怎么都要补兵
	//		if proto.All {
	//			// 差这么多兵
	//			diff := heroMilitary.SoldierCapcity() + needRecoverSolidier - heroMilitary.FreeSoldier()
	//			reduceTimes = u64.Min(u64.DivideTimes(diff, heroMilitary.NewSoldierRecruitCount()), times)
	//		} else {
	//			reduceTimes = 1
	//		}
	//
	//		// 减少次数
	//		recruitTimes.ReduceTimes(reduceTimes, ctime)
	//
	//		recruitSoldierCount := heroMilitary.NewSoldierRecruitCount() * reduceTimes
	//
	//		// 增加士兵数量
	//		heroMilitary.AddFreeSoldier(recruitSoldierCount)
	//
	//		result.Add(military.NewS2cUpdateFreeSoldierMsg(u64.Int32(heroMilitary.FreeSoldier())))
	//		result.Add(military.NewS2cRecruitSoldierV2Msg(u64.Int32(recruitSoldierCount)))
	//		result.Add(military.NewS2cRecruitSoldierTimesChangedMsg(timeutil.Marshal32(recruitTimes.StartRecoveryTime())))
	//
	//		// 历史累计数量
	//		hero.HistoryAmount().Increase(server_proto.HistoryAmountType_RecruitSoldierCount, recruitSoldierCount)
	//
	//		// 更新任务进度
	//		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_RECRUIT_SOLDIER_COUNT)
	//		heromodule.IncreTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_RECRUIT_SOLDIER, reduceTimes)
	//	}
	//
	//	if needRecoverSolidier > 0 {
	//		// 需要补兵
	//		heromodule.AutoRecoverCaptainSoldier(hero, result)
	//	}
	//
	//	result.Ok()
	//	result.Changed()
	//	return
	//})
}

var maxHealSoldierMsg = military.NewS2cGetMaxHealSoldierMsg(0).Static()

//gogen:iface c2s_get_max_heal_soldier
func (m *MilitaryModule) ProcessGetMaxHealSoldier_deprecated(hc iface.HeroController) {
	hc.Send(maxHealSoldierMsg)

	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	heroMilitary := hero.Military()
	//
	//	amount := hero.Military().WoundedSoldier()
	//	if amount > 0 {
	//
	//		// 剩余可招募士兵数
	//		maxRecruitCount := u64.Sub(heroMilitary.SoldierCapcity(), heroMilitary.FreeSoldier()+heroMilitary.FightSoldier())
	//		if maxRecruitCount <= 0 {
	//			return
	//		}
	//
	//		costCount := heromodule.GetCostCount(hero, heroMilitary.SoldierLevelData().WoundedCost)
	//		if costCount <= 0 {
	//			return
	//		}
	//
	//		amount = u64.Min(amount, u64.Min(maxRecruitCount, costCount))
	//	}
	//	result.Add(military.NewS2cGetMaxHealSoldierMsg(u64.Int32(amount)))
	//	return
	//})
}

//gogen:iface
func (m *MilitaryModule) ProcessC2SHealWoundedSoldierMsg_deprecated(proto *military.C2SHealWoundedSoldierProto, hc iface.HeroController) {

	hc.Send(military.ERR_HEAL_WOUNDED_SOLDIER_FAIL_INVALID_COUNT)

	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	healCount := u64.FromInt32(proto.GetCount())
	//	if healCount <= 0 {
	//		logrus.Debug("治疗伤兵，治疗数<=0")
	//		result.Add(military.ERR_HEAL_WOUNDED_SOLDIER_FAIL_INVALID_COUNT)
	//		return
	//	}
	//
	//	heroMilitary := hero.Military()
	//
	//	// 伤兵数量
	//	woundedSoldierCount := heroMilitary.WoundedSoldier()
	//
	//	if woundedSoldierCount < healCount {
	//		logrus.Debug("治疗伤兵，伤兵没这么多")
	//		result.Add(military.ERR_HEAL_WOUNDED_SOLDIER_FAIL_SOLDIER_NOT_ENOUGH)
	//		return
	//	}
	//
	//	// 治疗士兵数
	//	toAddMaxCount := u64.Sub(heroMilitary.SoldierCapcity(), heroMilitary.FreeSoldier()+heroMilitary.FightSoldier())
	//	if healCount > toAddMaxCount {
	//		logrus.Debug("治疗伤兵，超出士兵上限")
	//		result.Add(military.ERR_HEAL_WOUNDED_SOLDIER_FAIL_SOLDIER_CAPCITY_OVERFLOW)
	//		return
	//	}
	//
	//	// 新兵招募消耗资源
	//	soldierLevelData := heroMilitary.SoldierLevelData()
	//	totalCost := soldierLevelData.WoundedCost.Multiple(healCount)
	//
	//	if !heromodule.TryReduceCost(hero, result, totalCost) {
	//		logrus.Debug("治疗伤兵，资源不足")
	//		result.Add(military.ERR_HEAL_WOUNDED_SOLDIER_FAIL_RESOURCE_NOT_ENOUGH)
	//		return
	//	}
	//	result.Changed()
	//
	//	// 开始加人
	//	heroMilitary.ReduceWoundedSoldier(healCount)
	//
	//	heroMilitary.AddFreeSoldier(healCount)
	//	result.Add(military.NewS2cHealWoundedSoldierMsg(u64.Int32(healCount)))
	//
	//	// 历史累计数量
	//	hero.HistoryAmount().Increase(server_proto.HistoryAmountType_HealSoldierCount, healCount)
	//
	//	// 更新任务进度
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_HEAL_SOLDIER_COUNT)
	//
	//	result.Ok()
	//	return
	//})
}

//func PrintResult(proto *shared_proto.CombatProto) {
//
//	player := make(map[int32]int32)
//	fmt.Println("攻方信息：")
//	i := int32(0)
//	for _, t := range proto.Attacker.Troops {
//		fmt.Println(t.Index, t.Captain.Name, t.Captain.Soldier*t.Captain.LifePerSoldier)
//
//		i++
//		player[i] = t.Captain.Soldier * t.Captain.LifePerSoldier
//	}
//
//	fmt.Println("防方信息：")
//	for _, t := range proto.Attacker.Troops {
//		fmt.Println(t.Index, t.Captain.Name, t.Captain.Soldier*t.Captain.LifePerSoldier)
//		i++
//		player[i] = t.Captain.Soldier * t.Captain.LifePerSoldier
//	}
//
//	fmt.Println("攻方胜利：", proto.AttackerWin)
//
//	for i, round := range proto.Rounds {
//		fmt.Println("轮次：", i)
//
//		for _, act := range round.Actions {
//			if act.MoveDirection == shared_proto.Direction_ORIGIN {
//				player[act.TargetIndex] = player[act.TargetIndex] - act.Damage
//				fmt.Println(act.Index, "打", act.TargetIndex, act.HurtType, act.Damage, "剩余血量：", player[act.TargetIndex])
//			} else {
//				fmt.Println(act.Index, "移动", act.MoveDirection)
//			}
//		}
//	}
//
//	fmt.Println("存活士兵：")
//	for _, s := range proto.AliveSolider {
//		fmt.Println(s.Key, s.Value)
//	}
//}

//gogen:iface c2s_recruit_captain_v2
func (m *MilitaryModule) ProcessRecruitCaptainV2_deprecated(hc iface.HeroController) {
	//var pveTroopChangedFuncs []func()
	//
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	heroMilitary := hero.Military()
	//
	//	// 查看一下是否有队伍空位可以加入
	//	troop, index := hero.GetRecruitCaptainTroop()
	//	if troop == nil {
	//		logrus.Debugf("招募寻访武将，没有空闲的队伍可以加")
	//		result.Add(military.ERR_RECRUIT_CAPTAIN_V2_FAIL_NO_FREE_TROOP)
	//		return
	//	}
	//
	//	captainGenerator := m.datas.GetCaptainGenerator(uint64(heroMilitary.CaptainCount()))
	//	if captainGenerator == nil {
	//		logrus.Debugf("招募寻访武将，没有武将可以招募")
	//		result.Add(military.ERR_RECRUIT_CAPTAIN_V2_FAIL_NO_CAPTAIN_CAN_RECRUIT)
	//		return
	//	}
	//
	//	if hero.Level() < captainGenerator.NeedHeroLevel {
	//		logrus.Debugf("招募寻访武将，君主等级不够")
	//		result.Add(military.ERR_RECRUIT_CAPTAIN_V2_FAIL_HERO_LEVEL_TOO_LOW)
	//		return
	//	}
	//
	//	ability := captainGenerator.RandomAbility()
	//	raceData := captainGenerator.RandomRaceData()
	//
	//	s := GenerateSeekCaptain(ability, raceData, m.datas)
	//
	//	ctime := m.timeService.CurrentTime()
	//
	//	captain := heroMilitary.RecruitCaptainWithCaptainProto(s, hero.LevelData, raceData, ctime, m.datas, hero.TaskList().GetTitleData)
	//	if captain == nil {
	//		logrus.Debugf("招募寻访武将，heroMilitary.RecruitCaptain 执行失败，%s", hero.Id())
	//		result.Add(military.ERR_RECRUIT_CAPTAIN_V2_FAIL_NO_CAPTAIN_CAN_RECRUIT)
	//		return
	//	}
	//	result.Changed()
	//
	//	captainProto := captain.EncodeClient()
	//
	//	troop.SetCaptainIfAbsent(index, captain)
	//
	//	sendIndex := troop.Sequence()*uint64(len(troop.Captains())) + index
	//	result.Add(military.NewS2cRecruitCaptainV2MarshalMsg(captainProto, u64.Int32(sendIndex+1)))
	//	heromodule.UpdateTroopFightAmount(hero, troop, result)
	//
	//	// 更新任务进度
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_COUNT)
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_LEVEL_COUNT)
	//
	//	result.Ok()
	//
	//	hero.WalkPveTroop(func(troop *entity.PveTroop) (endWalk bool) {
	//		if troop.AddCaptain(captain) {
	//			result.Add(military.NewS2cSetPveCaptainMsg(must.Marshal(troop.Encode())))
	//			pveTroopChangedFuncs = append(pveTroopChangedFuncs, func() { heromodule.OnHeroPveTroopChange(hc.Id(), troop.TroopData().PveTroopType) })
	//		}
	//		return
	//	})
	//
	//	return
	//})
	//
	//for _, changedFunc := range pveTroopChangedFuncs {
	//	changedFunc()
	//}
}

//func GenerateSeekCaptain(ability uint64, raceData *race.RaceData, datas *config.ConfigDatas) *shared_proto.CaptainSeekerProto {
//
//	s := &shared_proto.CaptainSeekerProto{}
//
//	s.Male = rand.Intn(2) == 1
//	s.IconId = datas.CaptainHeads().RandomHead(s.Male).Id
//
//	familyNameArray := datas.FamilyName().Array
//	s.FamilyName = familyNameArray[rand.Intn(len(familyNameArray))].FamilyName
//
//	if s.Male {
//		givenNameArray := datas.MaleGivenName().Array
//		s.GivenName = givenNameArray[rand.Intn(len(givenNameArray))].Name
//	} else {
//		givenNameArray := datas.FemaleGivenName().Array
//		s.GivenName = givenNameArray[rand.Intn(len(givenNameArray))].Name
//	}
//
//	s.Race = raceData.Race
//	s.Ability = u64.Int32(ability)
//
//	abilityData := datas.CaptainAbilityData().Must(ability)
//	s.Quality = abilityData.Quality
//
//	s.Name = military_data.CaptainName(s.FamilyName, s.GivenName, abilityData, raceData,
//		datas.CaptainLevelData().MinKeyData)
//
//	return s
//}

//gogen:iface c2s_random_captain_head
func (m *MilitaryModule) ProcessRandomCaptainHead_deprecated(hc iface.HeroController) {
	//var currHeads []string
	//var randHeads []string
	//if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	currHeads = hero.Military().GetCaptainHeads()
	//	randHeads = hero.Military().GetCandidateCaptainHeads()
	//	result.Ok()
	//}) {
	//	logrus.Debugf("获得玩家当前武将头像失败，就从所有头像中随机")
	//	currHeads = make([]string, 0)
	//	randHeads = make([]string, 0)
	//}
	//
	//if len(randHeads) > 0 {
	//	hc.Send(military.NewS2cRandomCaptainHeadMsg(randHeads))
	//	return
	//}
	//
	//allCaptainHeads := m.datas.MilitaryConfig().AllCaptainHeads()
	//randSize := u64.Int(m.datas.MilitaryConfig().CaptainSeekerCandidateCount)
	//
	//randHeads = randHead(allCaptainHeads, currHeads, randSize)
	//hc.Send(military.NewS2cRandomCaptainHeadMsg(randHeads))
	//
	//hc.Func(func(hero *entity.Hero, err error) (heroChanged bool) {
	//	if err != nil {
	//		return
	//	}
	//	hero.Military().SetCandidateCaptainHeads(randHeads)
	//	return true
	//})
}

func randHead(allCaptainHeads []string, currHeads []string, randSize int) (randHeads []string) {
	headPool := make([]string, 0)
	for _, h := range allCaptainHeads {
		if !strContains(currHeads, h) {
			headPool = append(headPool, h)
		}
	}

	random.MixStrArray(headPool)

	if len(headPool) >= randSize {
		randHeads = headPool[:randSize]
	} else {
		randHeads = headPool
		random.MixStrArray(allCaptainHeads)
		size := randSize - len(headPool)
		randHeads = append(randHeads, allCaptainHeads[:size]...)
	}

	return
}

func strContains(strs []string, str string) bool {
	if strs == nil || len(strs) <= 0 {
		return false
	}
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

//gogen:iface
func (m *MilitaryModule) ProcessRecruitCaptainSeeker_deprecated(proto *military.C2SRecruitCaptainSeekerProto, hc iface.HeroController) {

	//var pveTroopChangedFuncs []func()
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	heroMilitary := hero.Military()
	//
	//	data := m.datas.GetCaptainSeekerData(military_data.CaptainSeekerDataId(uint64(proto.Index), uint64(heroMilitary.CaptainCount()+1)))
	//	if data == nil {
	//		logrus.Debugf("招募寻访武将，招募武将数据没找到")
	//		result.Add(military.ERR_RECRUIT_CAPTAIN_SEEKER_FAIL_NO_FREE_TROOP)
	//		return
	//	}
	//
	//	// 查看一下是否有队伍空位可以加入
	//	troop, index := hero.GetRecruitCaptainTroop()
	//	if troop == nil {
	//		logrus.Debugf("招募寻访武将，没有空闲的队伍可以加")
	//		result.Add(military.ERR_RECRUIT_CAPTAIN_SEEKER_FAIL_NO_FREE_TROOP)
	//		return
	//	}
	//
	//	if hero.Level() < data.RequiredHeroLevel {
	//		logrus.Debugf("招募寻访武将，君主等级不够")
	//		result.Add(military.ERR_RECRUIT_CAPTAIN_SEEKER_FAIL_HERO_LEVEL_TOO_LOW)
	//		return
	//	}
	//
	//	ctime := m.timeService.CurrentTime()
	//
	//	seeker := data.GetSeeker()
	//	seeker.IconId = proto.Head
	//	captain := heroMilitary.RecruitCaptainWithCaptainProto(seeker, hero.LevelData, data.GetRaceData(), ctime, m.datas, hero.TaskList().GetTitleData)
	//	if captain == nil {
	//		logrus.Debugf("招募寻访武将，heroMilitary.RecruitCaptain 执行失败，%s", hero.Id())
	//		result.Add(military.ERR_RECRUIT_CAPTAIN_SEEKER_FAIL_NO_CAPTAIN_CAN_RECRUIT)
	//		return
	//	}
	//	result.Changed()
	//
	//	captainProto := captain.EncodeClient()
	//
	//	troop.SetCaptainIfAbsent(index, captain)
	//
	//	sendIndex := troop.Sequence()*uint64(len(troop.Captains())) + index
	//	result.Add(military.NewS2cRecruitCaptainSeekerMarshalMsg(proto.Index, captainProto, u64.Int32(sendIndex+1)))
	//	heromodule.UpdateTroopFightAmount(hero, troop, result)
	//
	//	// 更新任务进度
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_COUNT)
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_LEVEL_COUNT)
	//
	//	if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_SEEK_CATPAIN) {
	//		result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_SEEK_CATPAIN)))
	//	}
	//
	//	// 清除招募武将头像缓存
	//	hero.Military().ResetCandidateCaptainHeads()
	//
	//	result.Ok()
	//
	//	hero.WalkPveTroop(func(troop *entity.PveTroop) (endWalk bool) {
	//		if troop.AddCaptain(captain) {
	//			result.Add(military.NewS2cSetPveCaptainMsg(must.Marshal(troop.Encode())))
	//			pveTroopChangedFuncs = append(pveTroopChangedFuncs, func() { heromodule.OnHeroPveTroopChange(hc.Id(), troop.TroopData().PveTroopType) })
	//		}
	//		return
	//	})
	//
	//	return
	//})
	//
	//for _, changedFunc := range pveTroopChangedFuncs {
	//	changedFunc()
	//}
}

//gogen:iface
func (m *MilitaryModule) ProcessSellSeekCaptain_deprecated(proto *military.C2SSellSeekCaptainProto, hc iface.HeroController) {
	//
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	heroMilitary := hero.Military()
	//
	//	if heroMilitary.SeekerCaptainCount() <= 0 {
	//		logrus.Debugf("出售寻访武将，寻访武将列表为空")
	//		result.Add(military.ERR_SELL_SEEK_CAPTAIN_FAIL_EMPTY)
	//		return
	//	}
	//
	//	hctx := heromodule.NewContext(m.dep, operate_type.MilitarySellSeekCaptain)
	//	if proto.Index == 0 {
	//		// 将所有的都兑换成强化符
	//		seekers := heroMilitary.Seekers()
	//
	//		builder := resdata.NewPrizeBuilder()
	//		for _, s := range seekers {
	//			data := m.datas.CaptainAbilityData().Must(u64.FromInt32(s.Ability))
	//			builder.Add(data.SellPrice)
	//		}
	//		toAdd := builder.Build()
	//
	//		// 给东西
	//		heromodule.AddPrize(hctx, hero, result, toAdd, m.timeService.CurrentTime())
	//
	//		heroMilitary.ClearSeekCaptain()
	//
	//	} else {
	//		removeIndex, seekCaptain := heroMilitary.SeekerCaptain(proto.Index)
	//
	//		if seekCaptain == nil {
	//			logrus.Debugf("出售寻访武将，Index无效")
	//			result.Add(military.ERR_SELL_SEEK_CAPTAIN_FAIL_INVALID_INDEX)
	//			return
	//		}
	//
	//		data := m.datas.CaptainAbilityData().Must(u64.FromInt32(seekCaptain.Ability))
	//		toAdd := data.SellPrice
	//
	//		// 给东西
	//		heromodule.AddPrize(hctx, hero, result, toAdd, m.timeService.CurrentTime())
	//
	//		heroMilitary.RemoveSeekerCaptain(removeIndex, proto.Index)
	//	}
	//	result.Changed()
	//
	//	result.Add(military.NewS2cSellSeekCaptainMsg(proto.Index))
	//
	//	result.Ok()
	//	return
	//})
}

//gogen:iface
func (m *MilitaryModule) ProcessChangeCaptainName_deprecated(proto *military.C2SChangeCaptainNameProto, hc iface.HeroController) {

	//name := strings.TrimSpace(proto.Name)
	//if c := util.GetCharLen(name); c <= 0 || c > 14 {
	//	logrus.Debugf("武将改名，新名字长度无效 %s", proto.Name)
	//	hc.Send(military.ERR_CHANGE_CAPTAIN_NAME_FAIL_INVALID_NAME)
	//	return
	//}
	//
	//if !util.IsValidName(name) {
	//	logrus.Debugf("武将改名，名字包含非法字符，newName: %s", proto.Name)
	//	hc.Send(military.ERR_CHANGE_CAPTAIN_NAME_FAIL_INVALID_NAME)
	//	return
	//}
	//
	//// 通过TSS 查询名字是否可用
	//if !m.tssClient.TryCheckName("武将改名", hc, name, military.ERR_CHANGE_CAPTAIN_NAME_FAIL_SENSITIVE_WORDS, military.ERR_CHANGE_CAPTAIN_NAME_FAIL_SERVER_ERROR) {
	//	return
	//}
	//
	//hctx := heromodule.NewContext(m.dep, operate_type.MilitaryChangeCaptainName)
	//
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	cid := u64.FromInt32(proto.Id)
	//	captain := hero.Military().Captain(cid)
	//	if captain == nil {
	//		logrus.Debugf("武将改名，武将id不存在，%v %s", proto.Id, proto.Name)
	//		result.Add(military.ERR_CHANGE_CAPTAIN_NAME_FAIL_INVALID_ID)
	//		return
	//	}
	//
	//	if captain.OriginName() == name {
	//		logrus.Debugf("武将改名，新名字跟原来名字一样，%v %s", proto.Id, proto.Name)
	//		result.Add(military.ERR_CHANGE_CAPTAIN_NAME_FAIL_SAME_NAME)
	//		return
	//	}
	//
	//	// 也不能跟其他武将重名
	//	for _, c := range hero.Military().Captains() {
	//		if c != nil && c.OriginName() == name {
	//			if c == captain {
	//				logrus.Debugf("武将改名，新名字跟原来名字一样，%v %s", proto.Id, proto.Name)
	//				result.Add(military.ERR_CHANGE_CAPTAIN_NAME_FAIL_SAME_NAME)
	//				return
	//			} else {
	//				logrus.Debugf("武将改名，跟其他武将重名，%v %s", proto.Id, proto.Name)
	//				result.Add(military.ERR_CHANGE_CAPTAIN_NAME_FAIL_DUPLICATE_NAME)
	//				return
	//			}
	//		}
	//	}
	//
	//	if !heromodule.TryReduceCost(hctx, hero, result, m.datas.MiscConfig().ChangeCaptainNameCost) {
	//		logrus.Debugf("武将改名，消耗不足")
	//		result.Add(military.ERR_CHANGE_CAPTAIN_NAME_FAIL_COST_NOT_ENOUGH)
	//		return
	//	}
	//
	//	result.Changed()
	//
	//	// 扣完钱，开始操作
	//	captain.SetName(name)
	//
	//	// 结束
	//	result.Add(military.NewS2cChangeCaptainNameMsg(proto.Id, name))
	//
	//	result.Ok()
	//	return
	//})
}

//gogen:iface
func (m *MilitaryModule) ProcessFireCaptain_deprecated(proto *military.C2SFireCaptainProto, hc iface.HeroController) {
	hc.Send(military.ERR_FIRE_CAPTAIN_FAIL_SERVER_BUSY)

	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	heroMilitary := hero.Military()
	//
	//	id := u64.FromInt32(proto.Id)
	//
	//	captain := heroMilitary.Captain(id)
	//	if captain == nil {
	//		logrus.Debugf("武将解雇，武将id没找到, %d", id)
	//		result.Add(military.ERR_FIRE_CAPTAIN_FAIL_ID_NOT_FOUND)
	//		return
	//	}
	//
	//	if captain.Level() > m.datas.MilitaryConfig().FireLevelLimit {
	//		logrus.Debugf("武将解雇，武将等级太高")
	//		result.Add(military.ERR_FIRE_CAPTAIN_FAIL_LEVEL_LIMIT)
	//		return
	//	}
	//
	//	if captain.IsOutSide() {
	//		logrus.Debugf("武将解雇，武将出征中", hero.Id())
	//		result.Add(military.ERR_FIRE_CAPTAIN_FAIL_OUTSIDE)
	//		return
	//	}
	//
	//	depot := hero.Depot()
	//	if !depot.HasEnoughGenIdGoodsCapacity(goods.EQUIPMENT, captain.GetEquipmentCount()) {
	//		logrus.Debugf("武将解雇，装备背包空间不足")
	//		result.Add(military.ERR_FIRE_CAPTAIN_FAIL_DEPOT_EQUIPMENT_FULL)
	//		return
	//	}
	//
	//	result.Changed()
	//
	//	// 开始操作，将武将移除
	//	// 同时将防守阵容和武将编队移除
	//	heroMilitary.RemoveCaptain(id)
	//
	//	if captain.GetFuShenCaptainSoul() != nil {
	//		hero.CaptainSoul().CancelFuShen(captain.GetFuShenCaptainSoul().Id())
	//		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_FU_SHEN)
	//	}
	//
	//	ctime := m.timeService.CurrentTime()
	//
	//	// 脱装备
	//	equipments := captain.RemoveAllEquipment()
	//	heromodule.AddEquipArray(hero, result, equipments, ctime)
	//
	//	// 脱宝石
	//	heromodule.AddGemArrayGive1(hero, result, captain.RemoveAllGem())
	//
	//	// 还士兵
	//	if captain.Soldier() > 0 {
	//		newFreeSoldier := heroMilitary.AddFreeSoldier(captain.Soldier())
	//		result.Add(military.NewS2cUpdateFreeSoldierMsg(u64.Int32(newFreeSoldier)))
	//	}
	//
	//	result.Add(military.NewS2cFireCaptainMsg(u64.Int32(id)))
	//
	//	// 返还物品
	//	heromodule.AddPrize(hero, result, captain.AbilityData().FirePrice, ctime)
	//
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_LEVEL_Y)
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_EQUIPMENT)
	//
	//	result.Ok()
	//	return
	//})
}

//gogen:iface
func (m *MilitaryModule) ProcessCaptainRefined_deprecated(proto *military.C2SCaptainRefinedProto, hc iface.HeroController) {
	hc.Send(military.ERR_CAPTAIN_REFINED_FAIL_INVALID_CAPTAIN)
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	captainId := u64.FromInt32(proto.Captain)
	//	captain := hero.Military().Captain(captainId)
	//	if captain == nil {
	//		logrus.Debugf("武将成长，武将不存在")
	//		result.Add(military.ERR_CAPTAIN_REFINED_FAIL_INVALID_CAPTAIN)
	//		return
	//	}
	//
	//	//if captain.IsOutSide() {
	//	//	logrus.Debugf("武将成长，武将外出了")
	//	//	result.Add(military.ERR_CAPTAIN_REFINED_FAIL_OUTSIDE)
	//	//	return
	//	//}
	//
	//	if captain.AbilityData().Ability >= captain.RebirthLevelData().AbilityLimit {
	//		logrus.Debugf("武将成长，成长已达转生上限")
	//		result.Add(military.ERR_CAPTAIN_REFINED_FAIL_REBIRTH_LIMIT)
	//		return
	//	}
	//
	//	if captain.AbilityData().NextLevel() == nil {
	//		logrus.Debugf("武将成长，已经最高成长值了")
	//		result.Add(military.ERR_CAPTAIN_REFINED_FAIL_INVALID_CAPTAIN)
	//		return
	//	}
	//
	//	if len(proto.GoodsId) <= 0 {
	//		logrus.Debugf("武将成长，物品id个数==0")
	//		result.Add(military.ERR_CAPTAIN_REFINED_FAIL_INVALID_GOODS)
	//		return
	//	}
	//
	//	if len(proto.GoodsId) != len(proto.Count) {
	//		logrus.Debugf("武将成长，物品id和count的个数不一致")
	//		result.Add(military.ERR_CAPTAIN_REFINED_FAIL_INVALID_GOODS)
	//		return
	//	}
	//
	//	goodsCountMap := make(map[*goods.GoodsData]uint64, len(proto.GoodsId))
	//	for i, v := range proto.GoodsId {
	//		goodsId := u64.FromInt32(v)
	//		goods := m.datas.GetGoodsData(goodsId)
	//		if goods == nil {
	//			logrus.Debugf("武将成长，物品不存在")
	//			result.Add(military.ERR_CAPTAIN_REFINED_FAIL_INVALID_GOODS)
	//			return
	//		}
	//
	//		if goods.GoodsEffect == nil || goods.GoodsEffect.ExpType != shared_proto.GoodsExpEffectType_EXP_CAPTAIN_REFINED {
	//			logrus.Debugf("武将成长，物品不是加武将强化经验的")
	//			result.Add(military.ERR_CAPTAIN_REFINED_FAIL_INVALID_GOODS)
	//			return
	//		}
	//
	//		count := u64.FromInt32(proto.Count[i])
	//		if count <= 0 {
	//			logrus.Debugf("武将成长，物品个数无效")
	//			result.Add(military.ERR_CAPTAIN_REFINED_FAIL_INVALID_COUNT)
	//			return
	//		}
	//
	//		goodsCountMap[goods] += count
	//	}
	//
	//	heroDepot := hero.Depot()
	//	for goods, count := range goodsCountMap {
	//		if !heroDepot.HasEnoughGoods(goods.Id, count) {
	//			logrus.Debugf("武将成长，物品个数不足")
	//			result.Add(military.ERR_CAPTAIN_REFINED_FAIL_INVALID_COUNT)
	//			return
	//		}
	//	}
	//
	//	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCaptainUpgradeAbilityExp)
	//	result.Changed()
	//
	//	toAddExp := uint64(0)
	//	for goods, count := range goodsCountMap {
	//		// 开始操作，扣物品，加经验，升级
	//		heromodule.ReduceGoodsAnyway(hctx, hero, result, goods, count)
	//
	//		toAddExp += goods.GoodsEffect.Exp * count
	//	}
	//
	//	heromodule.AddCaptainAbilityExp(hctx, hero, result, captain, toAddExp)
	//	result.Add(military.NewS2cCaptainRefinedMsg(u64.Int32(captainId), u64.Int32(captain.AbilityExp())))
	//
	//	hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_CaptainRefinedTimes)
	//
	//	// 更新任务进度
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_REFINED_TIMES)
	//
	//	if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_REFINED_CAPTAIN) {
	//		result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_REFINED_CAPTAIN)))
	//	}
	//
	//	result.Ok()
	//	return
	//})
}

//gogen:iface
func (m *MilitaryModule) ProcessChangeCaptainRace_deprecated(proto *military.C2SChangeCaptainRaceProto, hc iface.HeroController) {

	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryChangeCaptainRace)
	//
	//	cid := u64.FromInt32(proto.Id)
	//	captain := hero.Military().Captain(cid)
	//	if captain == nil {
	//		logrus.Debugf("武将转职，武将id不存在，%v %s", proto.Id)
	//		result.Add(military.ERR_CHANGE_CAPTAIN_RACE_FAIL_INVALID_ID)
	//		return
	//	}
	//
	//	if captain.RebirthLevel() == 0 && captain.Level() < m.datas.MiscConfig().ChangeCaptainRaceLevel {
	//		logrus.Debugf("武将转职，武将等级不足")
	//		result.Add(military.ERR_CHANGE_CAPTAIN_RACE_FAIL_LEVEL_NOT_ENOUGH)
	//		return
	//	}
	//
	//	ctime := m.timeService.CurrentTime()
	//	if ctime.Before(captain.GetRaceCdEndTime()) {
	//		logrus.Debugf("武将转职，cd中")
	//		result.Add(military.ERR_CHANGE_CAPTAIN_RACE_FAIL_COOLDOWN)
	//		return
	//	}
	//
	//	newRaceData := m.datas.GetRaceData(int(proto.Race))
	//	if newRaceData == nil {
	//		logrus.Debugf("武将转职，无效的职业类型")
	//		result.Add(military.ERR_CHANGE_CAPTAIN_RACE_FAIL_INVALID_RACE)
	//		return
	//	}
	//
	//	if captain.Race() == newRaceData {
	//		logrus.Debugf("武将转职，转职职业跟当前职业相同")
	//		result.Add(military.ERR_CHANGE_CAPTAIN_RACE_FAIL_SAME_RACE)
	//		return
	//	}
	//
	//	if captain.IsOutSide() {
	//		logrus.Debugf("武将转职，武将出征中")
	//		result.Add(military.ERR_CHANGE_CAPTAIN_RACE_FAIL_OUTSIDE)
	//		return
	//	}
	//
	//	// 首次转职免费
	//	if !timeutil.IsZero(captain.GetRaceCdEndTime()) {
	//		costGoods := m.datas.GoodsConfig().ChangeCaptainRaceGoods
	//		if proto.Money {
	//			if costGoods.DianquanPrice <= 0 {
	//				logrus.Debugf("武将转职，不支持点券购买")
	//				result.Add(military.ERR_CHANGE_CAPTAIN_RACE_FAIL_NOT_SUPPORT_YUANBAO)
	//				return
	//			}
	//
	//			if !heromodule.ReduceDianquan(hctx, hero, result, costGoods.DianquanPrice) {
	//				logrus.Debugf("武将转职，点券购买，点券不足")
	//				result.Add(military.ERR_CHANGE_CAPTAIN_RACE_FAIL_COST_NOT_ENOUGH)
	//				return
	//			}
	//		} else {
	//
	//			heroDepot := hero.Depot()
	//			if !heroDepot.HasEnoughGoods(costGoods.Id, 1) {
	//				logrus.Debugf("武将转职，消耗物品不足")
	//				result.Add(military.ERR_CHANGE_CAPTAIN_RACE_FAIL_COST_NOT_ENOUGH)
	//				return
	//			}
	//
	//			heromodule.ReduceGoodsAnyway(hctx, hero, result, costGoods, 1)
	//		}
	//	}
	//
	//	oldLevel := captain.Level()
	//	oldAbility := captain.Ability()
	//	oldRace := captain.Race().Id
	//	oldOfficial := captain.Official().Id
	//
	//	captain.SetRace(newRaceData, hero.Military().GetOrCreateSoldierData(newRaceData))
	//	cdEndTime := ctime.Add(m.datas.MiscConfig().ChangeCaptainRaceDuration)
	//	captain.SetRaceCdEndTime(cdEndTime)
	//
	//	result.Add(military.NewS2cChangeCaptainRaceMsg(proto.Id, proto.Race, timeutil.Marshal32(cdEndTime), must.Marshal(captain.Name())))
	//
	//	// 更新武将属性
	//
	//	captain.CalculateProperties()
	//	result.Add(captain.NewUpdateCaptainStatMsg())
	//	heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)
	//
	//	// 更新任务进度
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_QUALITY_COUNT)
	//
	//	// tlog
	//	hctx.Tlog().TlogPlayerCultivateFlow(hero, captain.Id(), operate_type.CaptainOperTypeChangeRace, oldLevel, captain.Level(), oldAbility, captain.Ability(), u64.FromInt(oldRace), u64.FromInt(captain.Race().Id), oldOfficial, captain.Official().Id, hctx.OperId())
	//
	//	result.Changed()
	//	result.Ok()
	//	return
	//})
}

//gogen:iface
func (m *MilitaryModule) ProcessUseLevelExpGoods_deprecated(proto *military.C2SUseLevelExpGoodsProto, hc iface.HeroController) {
	hc.Send(military.ERR_USE_LEVEL_EXP_GOODS_FAIL_INVALID_GOODS)
	//// 这个旧版突飞令不用了，用 ProcessUseLevelExpGoods2
	//goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.GoodsId))
	//if goodsData == nil {
	//	logrus.WithField("goods", proto.GoodsId).Debug("使用突飞令，物品不存在")
	//	hc.Send(military.ERR_USE_LEVEL_EXP_GOODS_FAIL_INVALID_GOODS)
	//	return
	//}
	//if goodsData.EffectType != shared_proto.GoodsEffectType_EFFECT_TRAIN_EXP {
	//	logrus.WithField("goods", goodsData.Name).Debug("使用突飞令，物品不是突飞令")
	//	hc.Send(military.ERR_USE_LEVEL_EXP_GOODS_FAIL_INVALID_GOODS)
	//	return
	//}
	//useCount := u64.Max(u64.FromInt32(proto.Count), 1)
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	building := hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN)
	//	if building == nil {
	//		logrus.Debug("使用突飞令，修炼馆未开放")
	//		result.Add(military.ERR_USE_LEVEL_EXP_GOODS_FAIL_BUILDING_LOCKED)
	//		return
	//	}
	//	captain := hero.Military().Captain(u64.FromInt32(proto.CaptainId))
	//	if captain == nil {
	//		logrus.WithField("captain", proto.CaptainId).Debug("使用突飞令，武将不存在")
	//		hc.Send(military.ERR_USE_LEVEL_EXP_GOODS_FAIL_INVALID_GOODS)
	//		return
	//	}
	//	if captain.IsMaxLevel() {
	//		logrus.Debug("使用突飞令，武将已达最大等级")
	//		hc.Send(military.ERR_USE_LEVEL_EXP_GOODS_FAIL_CAPTAIN_MAX_LEVEL)
	//		return
	//	}
	//	if captain.IsLevelLimit(hero.LevelData().Sub.CaptainLevelLimit, true) {
	//		logrus.Debug("使用突飞令，武将等级受限，请提升君主等级")
	//		hc.Send(military.ERR_USE_LEVEL_EXP_GOODS_FAIL_CAPTAIN_LEVEL_LIMIT)
	//		return
	//	}
	//	ctime := m.timeService.CurrentTime()
	//	// 转生 CD 中，不能加经验
	//	if captain.IsInRebrithing(ctime) {
	//		logrus.Debug("使用突飞令，武将正在转生中")
	//		hc.Send(military.ERR_USE_LEVEL_EXP_GOODS_FAIL_IN_REBIRTHING)
	//		return
	//	}
	//	maxCanAddExp := captain.GetMaxCanAddExp(hero.LevelData().Sub.CaptainLevelLimit)
	//	if maxCanAddExp <= 0 {
	//		logrus.Debug("使用突飞令，已经不能再加经验了")
	//		hc.Send(military.ERR_USE_LEVEL_EXP_GOODS_FAIL_CAPTAIN_LEVEL_LIMIT)
	//		return
	//	}
	//	// 计算具体可以使用几个丹
	//	toAddExpPerCount := u64.Multi(building.Effect.TrainExpPerHour,
	//		goodsData.GoodsEffect.TrainDuration.Hours()*(1+hero.BuildingEffect().GetTrainCoef()))
	//	realUseCount := u64.DivideTimes(maxCanAddExp, toAddExpPerCount)
	//	realUseCount = u64.Min(realUseCount, useCount)
	//	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryUseTrainingExpGoods)
	//	if !heromodule.TryReduceGoods(hctx, hero, result, goodsData, realUseCount) {
	//		logrus.Debug("使用突飞令，物品个数不足")
	//		hc.Send(military.ERR_USE_LEVEL_EXP_GOODS_FAIL_GOODS_NOT_ENOUGH)
	//		return
	//	}
	//	toAddExp := toAddExpPerCount * realUseCount
	//	upgrade := heromodule.AddCaptainExp(hctx, hero, result, captain, toAddExp, ctime) // 加经验
	//	result.Add(military.NewS2cUseLevelExpGoodsMsg(proto.CaptainId, proto.GoodsId, u64.Int32(realUseCount), upgrade))
	//
	//	hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_UseTuFeiGoodsCount)
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_USE_TU_FEI_GOODS)
	//
	//	if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_USE_TU_FEI_GOODS) {
	//		result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_USE_TU_FEI_GOODS)))
	//	}
	//
	//	result.Changed()
	//	result.Ok()
	//})

}

// 解锁克制技
//gogen:iface
func (m *MilitaryModule) ProcessUnlockCaptainRestraintSpell_deprecated(proto *military.C2SUnlockCaptainRestraintSpellProto, hc iface.HeroController) {
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	captainId := u64.FromInt32(proto.Captain)
	//	captain := hero.Military().Captain(captainId)
	//	if captain == nil {
	//		logrus.Debugf("解锁武将技，武将不存在")
	//		result.Add(military.ERR_UNLOCK_CAPTAIN_RESTRAINT_SPELL_FAIL_CAPTAIN_NOT_FOUND)
	//		return
	//	}
	//
	//	//if captain.IsOutSide() {
	//	//	logrus.Debugf("解锁武将技，武将外出了")
	//	//	result.Add(military.ERR_UNLOCK_CAPTAIN_RESTRAINT_SPELL_FAIL_OUT_SIDE)
	//	//	return
	//	//}
	//
	//	if !captain.RestraintSpellUnlocked() {
	//		if captain.Ability() < captain.Race().UnlockRestraintSpellNeedAbility {
	//			logrus.Debugf("解锁武将技，武将成长值不够")
	//			result.Add(military.ERR_UNLOCK_CAPTAIN_RESTRAINT_SPELL_FAIL_ABILITY_NOT_ENOUGH)
	//			return
	//		}
	//
	//		captain.UnlockRestraintSpell()
	//		result.Changed()
	//	}
	//
	//	result.Add(military.NewS2cUnlockCaptainRestraintSpellMsg(proto.Captain))
	//
	//	result.Ok()
	//	return
	//})
}

// 武将详细属性
//gogen:iface
func (m *MilitaryModule) ProcessGetCaptainStatDetails_deprecated(proto *military.C2SGetCaptainStatDetailsProto, hc iface.HeroController) {
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	captainId := u64.FromInt32(proto.Captain)
	//	captain := hero.Military().Captain(captainId)
	//	if captain == nil {
	//		logrus.Debugf("查看武将详细属性，武将不存在")
	//		result.Add(military.ERR_GET_CAPTAIN_STAT_DETAILS_FAIL_CAPTAIN_NOT_FOUND)
	//		return
	//	}
	//
	//	result.Add(military.NewS2cGetCaptainStatDetailsMarshalMsg(proto.Captain, captain.GetDetailStat()))
	//	result.Ok()
	//	return
	//})
}

// 武将封官(老版的，没用了)

//gogen:iface
func (m *MilitaryModule) ProcessUpdateCaptainOfficial_deprecated(proto *military.C2SUpdateCaptainOfficialProto, hc iface.HeroController) {
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	captainId := u64.FromInt32(proto.Captain)
	//	captain := hero.Military().Captain(captainId)
	//	if captain == nil {
	//		logrus.Debugf("武将封官，武将不存在 id:%v", captainId)
	//		result.Add(military.ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_INVALID_CAPTAIN)
	//		return
	//	}
	//
	//	officialId := u64.FromInt32(proto.Official)
	//	official := m.datas.CaptainOfficialData().Get(officialId)
	//	if official == nil {
	//		logrus.Debugf("武将封官，官职不存在 official id:%v", officialId)
	//		result.Add(military.ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_INVALID_OFFICIAL)
	//		return
	//	}
	//
	//	if captain.Official().Id == officialId {
	//		logrus.Debugf("武将封官，已在这个官职上。official id:%v", officialId)
	//		result.Add(military.ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_ALREADY_ON_OFFICIAL)
	//		return
	//	}
	//
	//	if captain.GongXun() < official.NeedGongxun {
	//		logrus.Debugf("武将封官，功勋不够")
	//		result.Add(military.ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_GONGXUN_NOT_ENOUGH)
	//		return
	//	}
	//
	//	// 如果降职
	//	if officialId < captain.Official().Id {
	//		if captain.IsOutSide() {
	//			logrus.Debugf("武将封官，武将在外面")
	//			result.Add(military.ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_CAPTAIN_IS_OUTSIDE)
	//			return
	//		}
	//	}
	//
	//	maxCount := heroOfficialCount(hero.LevelData(), officialId)
	//	if hero.Military().GetOfficialCount(officialId) >= maxCount {
	//		logrus.Debugf("武将封官，册封数已达到君主等级限制。maxcount:%v", maxCount)
	//		result.Add(military.ERR_UPDATE_CAPTAIN_OFFICIAL_FAIL_MAX_COUNT)
	//		return
	//	}
	//
	//	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryUpdateCaptainOfficial)
	//	oldLevel := captain.Level()
	//	oldAbility := captain.Ability()
	//	oldRace := captain.Race().Id
	//	oldOfficial := captain.Official().Id
	//
	//	updateCaptainOffice(hero, captain, official)
	//	result.Add(military.NewS2cUpdateCaptainOfficialMsg(u64.Int32(captainId), u64.Int32(officialId)))
	//	result.Add(captain.NewUpdateCaptainStatMsg())
	//
	//	// 减掉一些士兵
	//	var toReduce uint64
	//	if captain.Soldier() > captain.SoldierCapcity() {
	//		toReduce := u64.Sub(captain.Soldier(), captain.SoldierCapcity())
	//		captain.ReduceSoldier(toReduce)
	//	}
	//	if toReduce > 0 {
	//		ctime := m.timeService.CurrentTime()
	//		hero.Military().AddFreeSoldier(toReduce, ctime)
	//		result.Add(hero.Military().NewUpdateFreeSoldierMsg())
	//	}
	//	heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)
	//
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_OFFICIAL_UPDATE)
	//
	//	if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_CAPTAIN_UPGRADE_OFFICIAL) {
	//		result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_CAPTAIN_UPGRADE_OFFICIAL)))
	//	}
	//
	//	// tlog
	//	hctx.Tlog().TlogPlayerCultivateFlow(hero, captain.Id(), operate_type.CaptainOperTypeOfficial, oldLevel, captain.Level(), oldAbility, captain.Ability(), u64.FromInt(oldRace), u64.FromInt(captain.Race().Id), oldOfficial, captain.Official().Id, hctx.OperId())
	//
	//	result.Changed()
	//	result.Ok()
	//})
}

// 武将卸任(老版的，没用了)
//gogen:iface
func (m *MilitaryModule) ProcessLeaveCaptainOfficial_deprecated(proto *military.C2SLeaveCaptainOfficialProto, hc iface.HeroController) {
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	captainId := u64.FromInt32(proto.Captain)
	//	captain := hero.Military().Captain(captainId)
	//	if captain == nil {
	//		logrus.Debugf("武将卸任，武将不存在 %v", captainId)
	//		result.Add(military.ERR_LEAVE_CAPTAIN_OFFICIAL_FAIL_INVALID_CAPTAIN)
	//		return
	//	}
	//
	//	if captain.Official().Id == 0 {
	//		logrus.Debugf("武将卸任，武将已经没有官职了")
	//		result.Add(military.ERR_LEAVE_CAPTAIN_OFFICIAL_FAIL_ALREADY_NO_OFFICIAL)
	//		return
	//	}
	//
	//	if captain.IsOutSide() {
	//		logrus.Debugf("武将卸任，武将在外面")
	//		result.Add(military.ERR_LEAVE_CAPTAIN_OFFICIAL_FAIL_CAPTAIN_IS_OUTSIDE)
	//		return
	//	}
	//
	//	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryLeaveCaptainOfficial)
	//	oldLevel := captain.Level()
	//	oldAbility := captain.Ability()
	//	oldRace := captain.Race().Id
	//	oldOfficial := captain.Official().Id
	//
	//	updateCaptainOffice(hero, captain, military_data.EmptyOfficialData)
	//
	//	result.Add(military.NewS2cLeaveCaptainOfficialMsg(u64.Int32(captainId)))
	//	result.Add(captain.NewUpdateCaptainStatMsg())
	//
	//	// 减掉一些士兵
	//	var toReduce uint64
	//	if captain.Soldier() > captain.SoldierCapcity() {
	//		toReduce := u64.Sub(captain.Soldier(), captain.SoldierCapcity())
	//		captain.ReduceSoldier(toReduce)
	//	}
	//	if toReduce > 0 {
	//		ctime := m.timeService.CurrentTime()
	//		hero.Military().AddFreeSoldier(toReduce, ctime)
	//		result.Add(hero.Military().NewUpdateFreeSoldierMsg())
	//	}
	//	heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)
	//
	//	// tlog
	//	hctx.Tlog().TlogPlayerCultivateFlow(hero, captain.Id(), operate_type.CaptainOperTypeOfficial, oldLevel, captain.Level(), oldAbility, captain.Ability(), u64.FromInt(oldRace), u64.FromInt(captain.Race().Id), oldOfficial, captain.Official().Id, hctx.OperId())
	//
	//	result.Changed()
	//	result.Ok()
	//})
}

//// 更新武将官职
//func updateCaptainOffice(hero *entity.Hero, captain *entity.Captain, newOfficial *military_data.CaptainOfficialData) {
//	oldOfficialId := captain.Official().Id
//	if oldOfficialId > 0 {
//		hero.Military().AddOfficialCount(oldOfficialId, -1)
//	}
//	if newOfficial.Id > 0 {
//		hero.Military().AddOfficialCount(newOfficial.Id, 1)
//	}
//	captain.SetOfficial(newOfficial)
//	captain.CalculateProperties()
//}

//gogen:iface
func (m *MilitaryModule) ProcessUseGongxunGoods_deprecated(proto *military.C2SUseGongxunGoodsProto, hc iface.HeroController) {
	hc.Send(military.ERR_USE_GONGXUN_GOODS_FAIL_INVALID_CAPTAIN)
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	captain := hero.Military().Captain(u64.FromInt32(proto.Captain))
	//	if captain == nil {
	//		logrus.Debug("使用功勋令，无效的武将id")
	//		result.Add(military.ERR_USE_GONGXUN_GOODS_FAIL_INVALID_CAPTAIN)
	//		return
	//	}
	//	maxGongXunOfficialData := hero.LevelData().GetCaptainOfficialCountData().GetMaxGongXunOfficialData()
	//	if maxGongXunOfficialData == nil {
	//		logrus.Debug("使用功勋令，maxGongXunOfficialData == nil")
	//		result.Add(military.ERR_USE_GONGXUN_GOODS_FAIL_OFFICIAL_LIMIT)
	//		return
	//	}
	//	canAddGongXun := u64.Sub(maxGongXunOfficialData.NeedGongxun, captain.GongXun())
	//	if canAddGongXun <= 0 {
	//		logrus.Debug("使用功勋令，已达最高功勋职位")
	//		result.Add(military.ERR_USE_GONGXUN_GOODS_FAIL_OFFICIAL_LIMIT)
	//		return
	//	}
	//	var toAddExp uint64 // 加的功勋
	//	goodsCountMap := make(map[*goods.GoodsData]uint64, len(proto.GoodsId))
	//	for i, v := range proto.GoodsId {
	//		goodsId := u64.FromInt32(v)
	//		goods := m.datas.GetGoodsData(goodsId)
	//		if goods == nil {
	//			logrus.Debugf("使用功勋令，物品不存在")
	//			result.Add(military.ERR_USE_GONGXUN_GOODS_FAIL_INVALID_GOODS)
	//			return
	//		}
	//		if goods.GoodsEffect == nil ||
	//			goods.GoodsEffect.ExpType != shared_proto.GoodsExpEffectType_EXP_GONG_XUN {
	//			logrus.Debugf("使用功勋令，物品不是加武将功勋的")
	//			result.Add(military.ERR_USE_GONGXUN_GOODS_FAIL_INVALID_GOODS)
	//			return
	//		}
	//		count := u64.FromInt32(proto.Count[i])
	//		if count <= 0 {
	//			logrus.Debugf("使用功勋令，物品个数无效")
	//			result.Add(military.ERR_USE_GONGXUN_GOODS_FAIL_INVALID_COUNT)
	//			return
	//		}
	//		if toAddExp+goods.GoodsEffect.Exp*count >= canAddGongXun {
	//			sub := u64.Sub(canAddGongXun, toAddExp)
	//			useCount := u64.DivideTimes(sub, goods.GoodsEffect.Exp)
	//			goodsCountMap[goods] += useCount
	//			toAddExp += goods.GoodsEffect.Exp * useCount
	//			break
	//		}
	//		goodsCountMap[goods] += count
	//		toAddExp += goods.GoodsEffect.Exp * count
	//	}
	//	heroDepot := hero.Depot()
	//	for goods, count := range goodsCountMap {
	//		if !heroDepot.HasEnoughGoods(goods.Id, count) {
	//			logrus.Debugf("使用功勋令，物品个数不足")
	//			result.Add(military.ERR_USE_GONGXUN_GOODS_FAIL_INVALID_COUNT)
	//			return
	//		}
	//	}
	//
	//	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryGongXunGoods)
	//	// 扣物品
	//	for goods, count := range goodsCountMap {
	//		// 开始操作，扣物品，加经验，升级
	//		heromodule.ReduceGoodsAnyway(hctx, hero, result, goods, count)
	//	}
	//	// 给武将加功勋
	//	captain.AddGongXun(toAddExp)
	//	result.Add(military.NewS2cUseGongxunGoodsMsg(proto.Captain, u64.Int32(captain.GongXun())))
	//
	//	result.Changed()
	//	result.Ok()
	//})

}

//gogen:iface
func (m *MilitaryModule) ProcessUseTrainingExpGoods_deprecated(proto *military.C2SUseTrainingExpGoodsProto, hc iface.HeroController) {

	//goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.GoodsId))
	//if goodsData == nil {
	//	logrus.WithField("goods", proto.GoodsId).Debug("使用修炼馆经验丹，物品不存在")
	//	hc.Send(military.ERR_USE_TRAINING_EXP_GOODS_FAIL_INVALID_GOODS)
	//	return
	//}
	//
	//if goodsData.EffectType != shared_proto.GoodsEffectType_EFFECT_TRAIN_EXP {
	//	logrus.WithField("goods", goodsData.Name).Debug("使用修炼馆经验丹，物品不是修炼馆经验丹")
	//	hc.Send(military.ERR_USE_TRAINING_EXP_GOODS_FAIL_INVALID_GOODS)
	//	return
	//}
	//
	//useCount := u64.Max(u64.FromInt32(proto.Count), 1)
	//
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	building := hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN)
	//	if building == nil {
	//		logrus.Debug("使用修炼馆经验丹，修炼馆未开放")
	//		result.Add(military.ERR_USE_TRAINING_EXP_GOODS_FAIL_BUILDING_LOCKED)
	//		return
	//	}
	//
	//	captain := hero.Military().Captain(u64.FromInt32(proto.CaptainId))
	//	if captain == nil {
	//		logrus.WithField("captain", proto.CaptainId).Debug("使用修炼馆经验丹，武将不存在")
	//		hc.Send(military.ERR_USE_TRAINING_EXP_GOODS_FAIL_INVALID_GOODS)
	//		return
	//	}
	//
	//	if captain.IsMaxLevel() {
	//		logrus.Debug("使用修炼馆经验丹，武将已达最大等级")
	//		hc.Send(military.ERR_USE_TRAINING_EXP_GOODS_FAIL_CAPTAIN_MAX_LEVEL)
	//		return
	//	}
	//
	//	if captain.IsLevelLimit(hero.LevelData().Sub.CaptainLevelLimit, true) {
	//		logrus.Debug("使用修炼馆经验丹，武将等级受限，请提升君主等级")
	//		hc.Send(military.ERR_USE_TRAINING_EXP_GOODS_FAIL_CAPTAIN_LEVEL_LIMIT)
	//		return
	//	}
	//
	//	ctime := m.timeService.CurrentTime()
	//	// 转生 CD 中，不能加经验
	//	if captain.IsInRebrithing(ctime) {
	//		logrus.Debug("使用修炼馆经验丹，武将正在转生中")
	//		hc.Send(military.ERR_USE_TRAINING_EXP_GOODS_FAIL_IN_REBIRTHING)
	//		return
	//	}
	//
	//	maxCanAddExp := captain.GetMaxCanAddExp(hero.LevelData().Sub.CaptainLevelLimit)
	//	if maxCanAddExp <= 0 {
	//		logrus.Debug("使用修炼馆经验丹，已经不能再加经验了")
	//		hc.Send(military.ERR_USE_TRAINING_EXP_GOODS_FAIL_CAPTAIN_LEVEL_LIMIT)
	//		return
	//	}
	//
	//	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryUseTrainingExpGoods)
	//
	//	// 计算具体可以使用几个丹
	//	toAddExpPerCount := u64.Multi(building.Effect.TrainExpPerHour,
	//		goodsData.GoodsEffect.TrainDuration.Hours()*(1+hero.BuildingEffect().GetTrainCoef()))
	//	realUseCount := u64.DivideTimes(maxCanAddExp, toAddExpPerCount)
	//	realUseCount = u64.Min(realUseCount, useCount)
	//
	//	if !heromodule.TryReduceGoods(hctx, hero, result, goodsData, realUseCount) {
	//		logrus.Debug("使用修炼馆经验丹，物品个数不足")
	//		hc.Send(military.ERR_USE_TRAINING_EXP_GOODS_FAIL_GOODS_NOT_ENOUGH)
	//		return
	//	}
	//
	//	toAddExp := toAddExpPerCount * realUseCount
	//	upgrade := heromodule.AddCaptainExp(hctx, hero, result, captain, toAddExp, ctime) // 加经验
	//
	//	result.Add(military.NewS2cUseTrainingExpGoodsMsg(proto.CaptainId, proto.GoodsId, u64.Int32(realUseCount), upgrade))
	//
	//	hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_UseTuFeiGoodsCount)
	//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_USE_TU_FEI_GOODS)
	//
	//	if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_USE_TU_FEI_GOODS) {
	//		result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_USE_TU_FEI_GOODS)))
	//	}
	//
	//	result.Changed()
	//	result.Ok()
	//})

}

//gogen:iface
func (m *MilitaryModule) ProcessCaptainRebirth_deprecated(proto *military.C2SCaptainRebirthProto, hc iface.HeroController) {
	//hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCaptainRebirth)
	//
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//
	//	cid := u64.FromInt32(proto.Id)
	//	captain := hero.Military().Captain(cid)
	//	if captain == nil {
	//		logrus.Debugf("武将转生，武将不存在，%v", proto.Id)
	//		result.Add(military.ERR_CAPTAIN_REBIRTH_FAIL_INVALID_ID)
	//		return
	//	}
	//
	//	nextLevel := captain.RebirthLevelData().GetNextLevel()
	//	if nextLevel == nil {
	//		logrus.Debugf("武将转生，武将已转生到最高等级")
	//		result.Add(military.ERR_CAPTAIN_REBIRTH_FAIL_MAX_LEVEL)
	//		return
	//	}
	//
	//	if !captain.NextLevelIsReBirth() {
	//		logrus.Debugf("武将转生，未达到转生等级")
	//		result.Add(military.ERR_CAPTAIN_REBIRTH_FAIL_LEVEL_NOT_ENOUGH)
	//		return
	//	}
	//
	//	ctime := m.timeService.CurrentTime()
	//	if captain.IsInRebrithing(ctime) {
	//		if !proto.Miao {
	//			logrus.Debugf("武将转生，还在转生 CD 中")
	//			result.Add(military.ERR_CAPTAIN_REBIRTH_FAIL_IN_REBIRTHING)
	//			return
	//		}
	//
	//		d := m.datas.MiscConfig().MiaoCaptainRebirthDuration
	//		if d <= 0 {
	//			logrus.Debugln("秒武将转生cd，功能未开放")
	//			result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_NOT_IN_REBIRTH)
	//			return
	//		}
	//
	//		endTime := captain.RebirthEndTime()
	//		multi := int64((endTime.Sub(ctime) + d - 1) / d)
	//		if multi <= 0 {
	//			logrus.Errorln("秒武将转生cd，计算出来的multi <= 0")
	//			result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_COST_NOT_ENOUGH)
	//			return
	//		}
	//
	//		cost := m.datas.MiscConfig().MiaoCaptainRebirthCost.Multiple(u64.FromInt64(multi))
	//
	//		if !heromodule.TryReduceCost(hctx, hero, result, cost) {
	//			logrus.Debugln("秒武将转生cd，资源不足")
	//			result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_COST_NOT_ENOUGH)
	//			return
	//		}
	//
	//		captain.SetRebirthEndTime(ctime)
	//	}
	//
	//	oldLevel := captain.Level()
	//	oldAbility := captain.Ability()
	//
	//	// 转生
	//	if !captain.ProcessRebirth(ctime) {
	//		logrus.Debugf("武将转生，前面验证通过还失败")
	//		result.Add(military.ERR_CAPTAIN_REBIRTH_FAIL_LEVEL_NOT_ENOUGH)
	//		return
	//	}
	//
	//	// 转生赠送成长经验
	//	heromodule.AddCaptainAbilityExp(hctx, hero, result, captain, nextLevel.AbilityExp)
	//	if captain.TryUpgradeAbility() {
	//		// 更新任务进度
	//		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_QUALITY_COUNT)
	//	}
	//
	//	captain.CalculateProperties()
	//
	//	var toReduce uint64
	//	if captain.Soldier() > captain.SoldierCapcity() {
	//		// 减掉一些士兵
	//		toReduce := u64.Sub(captain.Soldier(), captain.SoldierCapcity())
	//		captain.ReduceSoldier(toReduce)
	//	}
	//	heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)
	//
	//	if toReduce > 0 {
	//		ctime := m.timeService.CurrentTime()
	//		hero.Military().AddFreeSoldier(toReduce, ctime)
	//		result.Add(hero.Military().NewUpdateFreeSoldierMsg())
	//	}
	//
	//	soldier := captain.Soldier()
	//	soldierCapcity := captain.SoldierCapcity()
	//	totalStatBytes := must.Marshal(captain.GetTotalStat())
	//	fightAmount := captain.FightAmount()
	//	fullSoldierFightAmount := captain.FullSoldierFightAmount()
	//
	//	result.Add(military.NewS2cCaptainRebirthMsg(
	//		proto.Id,
	//		nil,
	//		u64.Int32(captain.RebirthLevel()),
	//		0,
	//		int32(captain.Quality()),
	//		u64.Int32(captain.Ability()),
	//		u64.Int32(captain.AbilityExp()),
	//		u64.Int32(nextLevel.AbilityLimit),
	//		u64.Int32(soldier),
	//		u64.Int32(soldierCapcity),
	//		totalStatBytes,
	//		u64.Int32(fightAmount),
	//		u64.Int32(fullSoldierFightAmount)))
	//
	//	// 系统广播
	//	hctx = heromodule.NewContext(m.dep, operate_type.MilitaryCaptainRebirth)
	//	if d := hctx.BroadcastHelp().CaptainReBrith; d != nil {
	//		hctx.AddBroadcast(d, hero, result, 0, captain.RebirthLevel(), func() *i18n.Fields {
	//			text := d.NewTextFields()
	//			text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
	//			text.WithFields(data.KeyNum, captain.RebirthLevel())
	//			return text
	//		})
	//	}
	//
	//	result.Changed()
	//	result.Ok()
	//
	//	// tlog
	//	hctx.Tlog().TlogPlayerCultivateFlow(hero, captain.Id(), operate_type.CaptainOperTypeUpgrade, oldLevel, captain.Level(), oldAbility, captain.Ability(), u64.FromInt(captain.Race().Id), u64.FromInt(captain.Race().Id), captain.Official().Id, captain.Official().Id, hctx.OperId())
	//
	//	return
	//})
}
