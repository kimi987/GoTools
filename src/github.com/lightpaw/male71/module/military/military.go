package military

import (
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/captain"
	captaindata "github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7combatx"
	"github.com/lightpaw/pbutil"
	"math/rand"
	"strconv"
	"github.com/lightpaw/male7/gen/pb/gem"
)

func NewMilitaryModule(dep iface.ServiceDep, datas *config.ConfigDatas, individualServerConfig iface.IndividualServerConfig,
	fightModule iface.FightService, fightXService iface.FightXService,
	realmService iface.RealmService, tssClient iface.TssClient) *MilitaryModule {

	m := &MilitaryModule{}
	m.dep = dep
	m.datas = datas
	m.individualServerConfig = individualServerConfig
	m.fightModule = fightModule
	m.fightXService = fightXService
	m.timeService = dep.Time()
	m.realmService = realmService
	m.worldService = dep.World()
	m.broadcast = dep.Broadcast()
	m.tssClient = tssClient

	return m
}

//gogen:iface
type MilitaryModule struct {
	dep                    iface.ServiceDep
	datas                  *config.ConfigDatas
	fightModule            iface.FightService
	fightXService          iface.FightXService
	individualServerConfig iface.IndividualServerConfig
	timeService            iface.TimeService
	realmService           iface.RealmService
	worldService           iface.WorldService
	broadcast              iface.BroadcastService
	tssClient              iface.TssClient
}

func (m *MilitaryModule) config() *race.RaceConfig {
	return m.datas.RaceConfig()
}

//gogen:iface
func (m *MilitaryModule) ProcessC2SCaptainChangeSoldierMsg(proto *military.C2SCaptainChangeSoldierProto, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		captainId := u64.FromInt32(proto.GetId())
		toChange := proto.GetCount()

		heroMilitary := hero.Military()

		captain := heroMilitary.Captain(captainId)
		if captain == nil {
			logrus.Debug("补兵，你没有这个武将")
			result.Add(military.ERR_CAPTAIN_CHANGE_SOLDIER_FAIL_NOT_OWNER)
			return
		}

		if captain.IsOutSide() {
			logrus.Debugf("补兵，武将出征中")
			result.Add(military.ERR_CAPTAIN_CHANGE_SOLDIER_FAIL_OUTSIDE)
			return
		}
		result.Changed()

		ctime := m.timeService.CurrentTime()
		if toChange >= 0 {
			toAddMaxCount := u64.Sub(captain.SoldierCapcity(), captain.Soldier())
			toAdd := u64.FromInt32(toChange)
			if toAdd == 0 || toAdd > toAddMaxCount {
				toAdd = toAddMaxCount
			}

			freeSoldier := heroMilitary.FreeSoldier(ctime)
			if toAdd > freeSoldier {
				toAdd = freeSoldier
			}

			heroMilitary.ReduceFreeSoldier(toAdd, ctime)
			captain.AddSoldier(toAdd)
		} else {
			toReduce := u64.FromInt32(-toChange)
			if toReduce > captain.Soldier() {
				toReduce = captain.Soldier()
			}

			captain.ReduceSoldier(toReduce)
			heroMilitary.AddFreeSoldier(toReduce, ctime)
		}

		newSoldier := captain.Soldier()
		fightAmount := captain.FightAmount()

		result.Add(military.NewS2cCaptainChangeSoldierMsg(u64.Int32(captainId), u64.Int32(newSoldier), u64.Int32(fightAmount), u64.Int32(heroMilitary.FreeSoldier(ctime))))
		result.Add(hero.Military().NewUpdateFreeSoldierMsg())

		if toChange >= 0 {
			heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ADD_CAPTAIN_SOLDIER)
		}

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_RECOVER_SOLDIER) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_RECOVER_SOLDIER)))
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ADD_CAPTAIN_SOLDIER)
		}

		result.Ok()
		return
	})

}

//gogen:iface
func (m *MilitaryModule) ProcessC2SCaptainFullSoldierMsg(proto *military.C2SCaptainFullSoldierProto, hc iface.HeroController) {

	if check.Int32Duplicate(proto.GetId()) {
		logrus.Debug("一键补兵，Id重复")
		hc.Send(military.ERR_CAPTAIN_FULL_SOLDIER_FAIL_DUPLICATE)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroMilitary := hero.Military()

		ctime := m.timeService.CurrentTime()
		freeSoldier := heroMilitary.FreeSoldier(ctime)
		if freeSoldier <= 0 {
			logrus.Debug("一键补兵，空闲士兵为0")
			result.Add(military.ERR_CAPTAIN_FULL_SOLDIER_FAIL_EMPTY_SOLDIER)
			return
		}

		var updateCaptains []*entity.Captain
		for _, v := range proto.GetId() {
			captain := heroMilitary.Captain(u64.FromInt32(v))
			if captain == nil {
				logrus.Debug("一键补兵，你没有这个武将")
				result.Add(military.ERR_CAPTAIN_FULL_SOLDIER_FAIL_CAPTAIN_NOT_EXIST)
				return
			}

			if captain.IsOutSide() {
				logrus.Debugf("一键补兵，武将出征中")
				result.Add(military.ERR_CAPTAIN_FULL_SOLDIER_FAIL_OUTSIDE)
				return
			}

			if captain.Soldier() >= captain.SoldierCapcity() {
				continue
			}

			updateCaptains = append(updateCaptains, captain)
		}

		if len(updateCaptains) <= 0 {
			result.Add(military.NewS2cCaptainFullSoldierMsg(nil, nil, nil, u64.Int32(freeSoldier)))
			return
		}

		originFreeSoldier := freeSoldier

		var ids, newSoldiers, fightAmounts []int32
		for _, captain := range updateCaptains {
			toAdd := u64.Sub(captain.SoldierCapcity(), captain.Soldier())
			if toAdd <= 0 {
				continue
			}

			if toAdd > freeSoldier {
				toAdd = freeSoldier
			}

			result.Changed()
			heroMilitary.ReduceFreeSoldier(toAdd, ctime)
			captain.AddSoldier(toAdd)

			ids = append(ids, u64.Int32(captain.Id()))
			newSoldiers = append(newSoldiers, u64.Int32(captain.Soldier()))
			fightAmounts = append(fightAmounts, u64.Int32(captain.FightAmount()))

			freeSoldier = u64.Sub(freeSoldier, toAdd)
			if freeSoldier <= 0 {
				break
			}
		}

		result.Add(military.NewS2cCaptainFullSoldierMsg(ids, newSoldiers, fightAmounts, u64.Int32(freeSoldier)))

		if originFreeSoldier != freeSoldier {
			result.Add(hero.Military().NewUpdateFreeSoldierMsg())
		}

		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ADD_CAPTAIN_SOLDIER)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_RECOVER_SOLDIER) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_RECOVER_SOLDIER)))
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ADD_CAPTAIN_SOLDIER)
		}

		result.Ok()
		return
	})
}

//gogen:iface c2s_force_add_soldier
func (m *MilitaryModule) ProcessForceAddSoldierMsg(hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryForceAddSoldier)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		building := hero.Domestic().GetBuilding(shared_proto.BuildingType_JUN_YING)
		if building == nil {
			logrus.Debug("强征，军营建筑未解锁")
			result.Add(military.ERR_FORCE_ADD_SOLDIER_FAIL_LOCKED)
			return
		}
		forceSoldier := hero.BuildingEffect().GetForceSoldier()
		if forceSoldier <= 0 {
			logrus.Debug("强征，forceSoldier <= 0")
			result.Add(military.ERR_FORCE_ADD_SOLDIER_FAIL_LOCKED)
			return
		}

		currentTimes := hero.Military().GetForceAddSoldierTimes()
		if maxTimes := m.datas.JunYingMiscData().ForceAddSoldierMaxTimes; maxTimes > 0 && currentTimes >= maxTimes {

			logrus.Debug("强征，强征次数已达上限")
			result.Add(military.ERR_FORCE_ADD_SOLDIER_FAIL_TIMES_LIMIT)
			return
		}

		// 士兵已满
		ctime := m.timeService.CurrentTime()
		freeSoldier := hero.Military().FreeSoldier(ctime)
		if freeSoldier >= hero.BuildingEffect().SoldierCapcity() {
			logrus.Debug("强征，当前士兵已满")
			result.Add(military.ERR_FORCE_ADD_SOLDIER_FAIL_FULL_SOLDIER)
			return
		}

		// 消耗
		cost := m.datas.JunYingMiscData().GetForceAddSoldierCost(currentTimes)
		if !heromodule.TryReduceCost(hctx, hero, result, cost) {
			logrus.Debug("强征，消耗不足")
			result.Add(military.ERR_FORCE_ADD_SOLDIER_FAIL_COST_NOT_ENOUGH)
			return
		}

		// 加强征次数
		newTimes := hero.Military().IncreseForceAddSoldierTimes()
		result.Add(military.NewS2cForceAddSoldierMsg(u64.Int32(newTimes)))

		// 加士兵
		hero.Military().AddFreeSoldier(forceSoldier, ctime)
		result.Add(hero.Military().NewUpdateFreeSoldierMsg())
	})
}

//gogen:iface c2s_upgrade_soldier_level
func (m *MilitaryModule) ProcessC2SUpgradeSoldierMsg(hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryUpgradeSoldier)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroMilitary := hero.Military()
		soldierLevelData := heroMilitary.SoldierLevelData()

		nextLevel := m.datas.GetSoldierLevelData(soldierLevelData.Level + 1)
		if nextLevel == nil {
			logrus.Debugf("升级士兵等级，已经是最高级了")
			result.Add(military.ERR_UPGRADE_SOLDIER_LEVEL_FAIL_MAX_LEVEL)
			return
		}

		junYing := hero.Domestic().GetBuilding(shared_proto.BuildingType_JUN_YING)
		if junYing == nil || junYing.Level < nextLevel.JunYingLevel {
			logrus.WithField("军营", fmt.Sprintf("%+v", junYing)).Debugf("升级士兵等级，军营等级不够")
			result.Add(military.ERR_UPGRADE_SOLDIER_LEVEL_FAIL_JUN_YING_LEVEL_TOO_LOW)
			return
		}

		if !heromodule.TryReduceCost(hctx, hero, result, nextLevel.UpgradeCost) {
			logrus.Debug("升级士兵等级，资源不足")
			result.Add(military.ERR_UPGRADE_SOLDIER_LEVEL_FAIL_RES_NOT_ENOUGH)
			return
		}
		result.Changed()

		heroMilitary.SetSoldierLevelData(nextLevel)

		result.Add(military.NewS2cUpgradeSoldierLevelMsg(u64.Int32(nextLevel.Level)))

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_SOLDIER_LEVEL)

		// 同步数据
		heroMilitary.SetSoldierLevelData(nextLevel)

		// 更新所有的武将数据
		for _, captain := range heroMilitary.Captains() {
			captain.CalculateProperties()
			result.Add(captain.NewUpdateCaptainStatMsg())
		}
		heromodule.UpdateAllTroopFightAmount(hero, result)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_UPGRADE_SOLDIER) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_UPGRADE_SOLDIER)))
		}

		result.Ok()

		// tlog
		hctx.Tlog().TlogResearchFlow(hero, operate_type.ResearchSolider, 0, soldierLevelData.Level, nextLevel.Level)

		return
	})
}

//gogen:iface c2s_fight
func (m *MilitaryModule) ProcessC2SFightMsg(proto *military.C2SFightProto, hc iface.HeroController) {
	if !m.individualServerConfig.GetIsDebug() {
		return
	}

	// TODO
	config := m.datas.MilitaryConfig()

	attackerId, attacker := genFightObj(m.datas.GetCaptainDataArray(), true, true)
	defenserId, defenser := genFightObj(m.datas.GetCaptainDataArray(), false, true)
	if proto.Wall {
		defenser.WallStat = &shared_proto.SpriteStatProto{Attack: 100, Defense: 100, Strength: 10}
		defenser.WallFixDamage = 10
		defenser.TotalWallLife = i32.Max(defenser.WallStat.Strength, 1)
		defenser.WallLevel = 1
	}

	tfctx := entity.NewTlogFightContext(operate_type.BattleSingle, 0, 0, 0)
	response := m.fightModule.SendFightRequestReturnResult(tfctx, config.CombatRes, attackerId, defenserId, attacker, defenser, true)
	logrus.Debugf(response.String())
	if response.ReturnCode != 0 {
		logrus.Errorf("战斗创建失败, %s", response.ReturnMsg)
		return
	}

	if response.Result == nil {
		logrus.Error("response.Result == nil")
		return
	}

	combatProto := response.Result
	combatBytes, err := combatProto.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("战斗结果Marshal失败")
		return
	}

	logrus.Debugf("战斗结果: attackerWin: %v, rounds: %v", combatProto.GetAttackerWin(), len(combatProto.Rounds))
	//PrintResult(combatProto)

	hc.Send(military.NewS2cFightMsg(combatBytes))
}

//gogen:iface c2s_multi_fight
func (m *MilitaryModule) ProcessC2SMultiFightMsg(hc iface.HeroController) {
	if !m.individualServerConfig.GetIsDebug() {
		return
	}

	// TODO
	config := m.datas.MilitaryConfig()

	proto := &server_proto.MultiCombatRequestServerProto{}

	for i := 0; i < 6; i++ {
		id, obj := genFightObj(m.datas.GetCaptainDataArray(), true, true)

		proto.AttackerId = append(proto.AttackerId, id)
		proto.Attacker = append(proto.Attacker, obj)
	}

	for i := 0; i < 6; i++ {
		id, obj := genFightObj(m.datas.GetCaptainDataArray(), false, true)

		proto.DefenserId = append(proto.DefenserId, id)
		proto.Defenser = append(proto.Defenser, obj)
	}

	tfctx := entity.NewTlogFightContext(operate_type.BattleMulti, 0, 0, 0)
	response := m.fightModule.SendMultiFightRequestReturnResult(tfctx, config.CombatRes, proto.AttackerId, proto.DefenserId,
		proto.Attacker, proto.Defenser, 3, 3, 3, true)

	logrus.Debugf(response.String())
	if response.ReturnCode != 0 {
		logrus.Errorf("战斗创建失败, %s", response.ReturnMsg)
		return
	}

	combatProto := response.Result
	combatBytes, err := combatProto.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("战斗结果Marshal失败")
		return
	}

	logrus.Debugf("战斗结果: attackerWin: %v, fightCount: %v", combatProto.GetAttackerWin(), len(combatProto.GetCombats()))
	//PrintResult(combatProto)

	hc.Send(military.NewS2cMultiFightMsg(combatBytes))
}

//gogen:iface
func (m *MilitaryModule) ProcessC2SFightxMsg(proto *military.C2SFightxProto, hc iface.HeroController) {
	if !m.individualServerConfig.GetIsDebug() {
		return
	}

	config := m.datas.MilitaryConfig()

	hasAttackerCaptain := false
	attackerCaptains := m.datas.GetCaptainDataArray()
	if len(proto.Attacker) > 0 {
		captain := make([]*captain.CaptainData, len(proto.Attacker))
		for i, v := range proto.Attacker {
			data := m.datas.GetCaptainData(u64.FromInt32(v))
			if data != nil {
				captain[i] = data
				hasAttackerCaptain = true
			}
		}

		if hasAttackerCaptain {
			attackerCaptains = captain
		}
	}

	hasDefenserCaptain := false
	defenserCaptains := m.datas.GetCaptainDataArray()
	if len(proto.Defenser) > 0 {
		captain := make([]*captain.CaptainData, len(proto.Defenser))
		for i, v := range proto.Defenser {
			data := m.datas.GetCaptainData(u64.FromInt32(v))
			if data != nil {
				captain[i] = data
				hasDefenserCaptain = true
			}
		}

		if hasDefenserCaptain {
			defenserCaptains = captain
		}
	}

	attackerId, attacker := genFightObj(attackerCaptains, true, !hasAttackerCaptain)
	defenserId, defenser := genFightObj(defenserCaptains, false, !hasDefenserCaptain)
	if proto.Wall {
		defenser.WallStat = &shared_proto.SpriteStatProto{Attack: 100, Defense: 100, Strength: 10}
		defenser.WallFixDamage = 10
		defenser.TotalWallLife = i32.Max(defenser.WallStat.Strength, 1)
		defenser.WallLevel = 1
	}

	tfctx := entity.NewTlogFightContext(operate_type.BattleSingle, 0, 0, 0)
	response := m.fightXService.SendFightRequestReturnResult(tfctx, config.CombatRes, attackerId, defenserId, attacker, defenser, true)
	logrus.Debugf(response.String())
	if response.ReturnCode != 0 {
		logrus.Errorf("战斗创建失败, %s", response.ReturnMsg)
		return
	}

	if response.Result == nil {
		logrus.Error("response.Result == nil")
		return
	}

	combatProto := response.Result
	combatBytes, err := combatProto.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("战斗结果Marshal失败")
		return
	}

	logrus.Debugf("战斗结果: attackerWin: %v, maxFrame: %v", combatProto.GetAttackerWin(), combatProto.MaxFrame)
	combatx.PrintResult(combatProto)

	hc.Send(military.NewS2cFightxMsg(combatBytes))
}

var playerId int64 = 0

func genFightObj(captain []*captain.CaptainData, isAttacker, isRandom bool) (id int64, obj *shared_proto.CombatPlayerProto) {
	playerId++

	player := &shared_proto.CombatPlayerProto{}
	player.Hero = &shared_proto.HeroBasicProto{}

	player.Hero.Id = idbytes.ToBytes(playerId)
	player.Hero.Head = "1"
	player.Hero.Name = strconv.FormatInt(id, 10)
	player.Hero.Level = 1
	if isAttacker {
		player.Hero.GuildName = "进攻联盟"
	} else {
		player.Hero.GuildName = "防守联盟"
	}
	tfa := data.NewTroopFightAmount()
	if isRandom {
		for index := int32(1); index <= 5; index++ {
			troops := randomTroops(captain)

			tt := &shared_proto.CombatTroopsProto{
				FightIndex: index,
				Captain:    troops,
			}

			player.Troops = append(player.Troops, tt)
			tfa.AddInt32(tt.Captain.FightAmount)
		}
	} else {
		for i, c := range captain {
			if c == nil {
				continue
			}

			index := int32(i + 1)

			troops := genTroops0(u64.Int32(c.Id), c)
			tt := &shared_proto.CombatTroopsProto{
				FightIndex: index,
				Captain:    troops,
			}

			player.Troops = append(player.Troops, tt)
			tfa.AddInt32(tt.Captain.FightAmount)
		}
	}

	player.TotalFightAmount = tfa.ToI32()

	return playerId, player
}

var captainId int32 = 0

func randomTroops(captain []*captain.CaptainData) (troops *shared_proto.CaptainInfoProto) {
	captainId++
	id := (captainId % 10) + 1

	c := captain[rand.Intn(len(captain))]

	return genTroops0(id, c)
}

func genTroops0(id int32, captain *captain.CaptainData) (troops *shared_proto.CaptainInfoProto) {

	troops = &shared_proto.CaptainInfoProto{}
	troops.Id = id
	troops.Name = &shared_proto.CaptainNameProto{Name: fmt.Sprintf("Captain-%d", id)}
	troops.Race = captain.Race.Race
	troops.Quality = shared_proto.Quality_PURPLE
	troops.IconId = "wujiang8"

	troops.LifePerSoldier = rand.Int31n(5) + 5000
	troops.Soldier = rand.Int31n(20) + 10000
	troops.TotalSoldier = troops.Soldier

	totalStat := &data.SpriteStat{}
	totalStat.Attack = uint64(500 + rand.Int31n(30))
	totalStat.Defense = uint64(1000 + rand.Int31n(30))
	totalStat.Strength = uint64(1000 + rand.Int31n(30))
	totalStat.Dexterity = uint64(1000 + rand.Int31n(30))
	totalStat.SoldierCapcity = uint64(troops.TotalSoldier)

	troops.TotalStat = totalStat.Encode()

	troops.FightAmount = u64.Int32(totalStat.FightAmount(u64.FromInt32(troops.Soldier), 0))

	troops.Model = rand.Int31n(3)*5 + int32(troops.Race)

	troops.CaptainId = u64.Int32(captain.Id)

	star := captain.Star[rand.Intn(len(captain.Star))]
	troops.Star = int32(star.Star)

	troops.UnlockSpellCount = int32(len(star.Spell))

	return troops
}

//gogen:iface
func (m *MilitaryModule) ProcessSetDefenseTroop(proto *military.C2SSetDefenseTroopProto, hc iface.HeroController) {

	troopIndex := u64.FromInt32(proto.TroopIndex)

	var updateBaseDefenserFunc func()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if troopIndex > 0 {
			t := hero.GetTroopByIndex(troopIndex - 1)
			if t == nil {
				logrus.Debugf("设置防守队伍，无效的队伍编号")
				result.Add(military.ERR_SET_DEFENSE_TROOP_FAIL_INVALID_TROOP_INDEX)
				return
			}
		}

		if hero.GetHomeDefenseTroopIndex() == troopIndex {
			result.Add(military.NewS2cSetDefenseTroopMsg(proto.IsTent, proto.TroopIndex))
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_DEFENSER_FIGHTING)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_SET_DEFENSER)
			return
		}

		hero.SetHomeDefenseTroopIndex(troopIndex)

		result.Add(military.NewS2cSetDefenseTroopMsg(proto.IsTent, proto.TroopIndex))

		if update, newAmount := hero.TryUpdateHomeDefenserFightAmount(); update {
			result.Add(domestic.NewS2cUpdateHeroFightAmountMsg(u64.Int32(newAmount)))
		}

		hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_SET_DEFENSER)

		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_DEFENSER_FIGHTING)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_SET_DEFENSER)

		result.Changed()
		result.Ok()
		return
	})

	if updateBaseDefenserFunc != nil {
		updateBaseDefenserFunc()
	}
}

var (
	dontFalseMsg = military.NewS2cSetDefenserAutoFullSoldierMsg(false).Static()
	dontTrueMsg  = military.NewS2cSetDefenserAutoFullSoldierMsg(true).Static()
)

//gogen:iface
func (m *MilitaryModule) ProcessSetDefenserAutoFullSoldier(proto *military.C2SSetDefenserAutoFullSoldierProto, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if !m.datas.GetVipLevelData(hero.VipLevel()).WallAutoFullSoldier {
			result.Add(military.ERR_SET_DEFENSER_AUTO_FULL_SOLDIER_FAIL_VIP_LEVEL_LIMIT)
			return
		}

		hero.MiscData().SetDefenserDontAutoFullSoldier(proto.Dont)

		if proto.Dont {
			result.Add(dontTrueMsg)
		} else {
			result.Add(dontFalseMsg)
		}
		result.Ok()
	})
}

//gogen:iface
func (m *MilitaryModule) ProcessSetMultiCaptainIndex(proto *military.C2SSetMultiCaptainIndexProto, hc iface.HeroController) {

	if check.Int32DuplicateIgnoreZero(proto.Id) {
		logrus.Debugf("设置武将编队（多），武将id重复")
		hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_INVALID_ID)
		return
	}

	if len(proto.Id) != len(proto.XIndex) {
		logrus.Debugf("设置武将编队（多），武将id跟xindex长度不一致")
		hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_INVALID_ID)
		return
	}

	count := len(proto.Id)
	if count != len(proto.XIndex) {
		logrus.Debugf("设置武将编队（多），武将id重复")
		hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_INVALID_ID)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		tcc := hero.Military().TroopCaptainCount()

		captains := make([]*entity.Captain, count)
		for i, id := range proto.Id {
			if id == 0 {
				continue
			}

			if tcc > 0 && hero.LevelData().Sub.TroopsCaptainCount < tcc {
				idx := uint64(i) % tcc
				if idx >= hero.LevelData().Sub.TroopsCaptainCount {
					logrus.Debugf("设置武将编队（多），武将位置还未解锁")
					hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_INVALID_ID)
					return
				}
			}

			captain := hero.Military().Captain(u64.FromInt32(id))
			if captain == nil {
				logrus.Debugf("设置武将编队（多），无效的武将id")
				hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_INVALID_ID)
				return
			}

			captains[i] = captain
		}

		if proto.Index > 0 {
			t := hero.GetTroopByIndex(u64.FromInt32(proto.Index - 1))
			if t == nil {
				logrus.Debugf("设置武将编队（多），无效的队伍编号")
				hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_INVALID_INDEX)
				return
			}

			if t.IsOutside() {
				logrus.Debugf("设置武将编队（多），出征武将不能修改")
				hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_OUTSIDE)
				return
			}

			if count != int(tcc) {
				logrus.Debugf("设置武将编队（多），无效的队伍编号")
				hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_INVALID_ID)
				return
			}

			for _, captain := range captains {
				if captain != nil && captain.IsOutSide() {
					logrus.Debugf("设置武将编队（多），出征武将不能设置上阵")
					hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_OUTSIDE)
					return
				}
			}

			t.Set(captains, proto.XIndex)

		} else {

			if count < int(tcc)*len(hero.Troops()) {
				logrus.Debugf("设置武将编队（多），全部设置时，发送的id个数不对")
				hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_INVALID_ID)
				return
			}

			// 如果队伍出征，那么里面的人不能动
			for idx, t := range hero.Troops() {
				startIndex := idx * int(tcc)
				if t.IsOutside() {
					// 出征中的不能改变
					for i, pos := range t.Pos() {
						index := startIndex + i
						if pos.Captain() != captains[index] {
							logrus.Debugf("设置武将编队（多），出征武将不能修改")
							hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_OUTSIDE)
							return
						}
					}
				} else {
					// 不是出征队伍，设置的武将不能是出征状态的
					for i := 0; i < int(tcc); i++ {
						index := startIndex + i

						captain := captains[index]
						if captain != nil && captain.IsOutSide() {
							logrus.Debugf("设置武将编队（多），武将在别的出征队伍中")
							hc.Send(military.ERR_SET_MULTI_CAPTAIN_INDEX_FAIL_OUTSIDE)
							return
						}
					}
				}
			}

			for idx, t := range hero.Troops() {
				if !t.IsOutside() {
					startIndex := idx * int(tcc)
					endIndex := startIndex + int(tcc)
					t.Set(captains[startIndex:endIndex], proto.XIndex[startIndex:endIndex])
				}
			}

		}

		result.Add(military.NewS2cSetMultiCaptainIndexMsg(proto.Index, proto.Id, proto.XIndex))

		heromodule.UpdateAllTroopFightAmount(hero, result)

		result.Changed()
		result.Ok()
		return
	})
}

//gogen:iface c2s_set_pve_captain
func (m *MilitaryModule) ProcessSetPveCaptain(proto *military.C2SSetPveCaptainProto, hc iface.HeroController) {

	if check.Int32DuplicateIgnoreZero(proto.Id) {
		logrus.Debugf("设置pve编队，武将id重复")
		hc.Send(military.ERR_SET_PVE_CAPTAIN_FAIL_DUP_CAPTAIN_ID)
		return
	}

	if check.Int32CountIgnoreZero(proto.Id) <= 0 {
		logrus.Debugf("设置pve编队，竟然一个武将都不设置，不可以")
		hc.Send(military.ERR_SET_PVE_CAPTAIN_FAIL_NO_CAPTAIN)
		return
	}

	troopData := m.datas.GetPveTroopData(u64.FromInt32(proto.PveType))
	if troopData == nil {
		logrus.Debugf("设置pve编队，未知的pve类型")
		hc.Send(military.ERR_SET_PVE_CAPTAIN_FAIL_INVALID_PVE_TYPE)
		return
	}

	if len(proto.Id) != len(proto.XIndex) {
		logrus.Debugf("设置pve编队，len(proto.Id) != len(proto.XIndex)")
		hc.Send(military.ERR_SET_PVE_CAPTAIN_FAIL_INVALID_X_INDEX)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		t := hero.PveTroop(troopData.PveTroopType)
		if t == nil {
			logrus.Errorf("设置pve编队，pve队伍没找到")
			hc.Send(military.ERR_SET_PVE_CAPTAIN_FAIL_INVALID_PVE_TYPE)
			return
		}

		if len(proto.Id) != len(t.Captains()) {
			logrus.Errorf("设置pve编队，客户端发送过来的队伍人数不匹配")
			hc.Send(military.ERR_SET_PVE_CAPTAIN_FAIL_INVALID_CAPTAIN_COUNT)
			return
		}

		heroMilitary := hero.Military()

		newArray := make([]*entity.Captain, len(proto.Id))
		for index, cid := range proto.Id {
			if cid == 0 {
				continue
			}

			captain := heroMilitary.Captain(u64.FromInt32(cid))
			if captain == nil {
				logrus.Debugf("设置武将编队（多），同一位置上有多个元素")
				hc.Send(military.ERR_SET_PVE_CAPTAIN_FAIL_INVALID_ID)
				return
			}

			newArray[index] = captain
		}

		t.SetCaptain(newArray, proto.XIndex)

		result.Add(military.NewS2cSetPveCaptainMsg(must.Marshal(t.Encode())))

		// 轩辕会武未开启之前，使用千重楼的队伍
		if troopData.PveTroopType == shared_proto.PveTroopType_TOWER &&
			!hero.Bools().Get(shared_proto.HeroBoolType_BOOL_XUAN_YUAN) {
			if t := hero.PveTroop(shared_proto.PveTroopType_PVE_XUAN_YUAN); t != nil {
				t.SetCaptain(newArray, proto.XIndex)
				result.Add(military.NewS2cSetPveCaptainMsg(must.Marshal(t.Encode())))
			}
		}

		result.Changed()
		result.Ok()
		return
	})

	heromodule.OnHeroPveTroopChange(hc.Id(), troopData.PveTroopType)
}

//gogen:iface
func (m *MilitaryModule) ProcessCaptainEnhance(proto *military.C2SCaptainEnhanceProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captainId := u64.FromInt32(proto.Captain)
		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("武将成长，武将不存在")
			result.Add(military.ERR_CAPTAIN_ENHANCE_FAIL_INVALID_CAPTAIN)
			return
		}
		if captain.Ability() >= captain.RebirthAbilityLimit() {
			logrus.Debugf("武将成长，成长已达转生上限")
			result.Add(military.ERR_CAPTAIN_ENHANCE_FAIL_REBIRTH_LIMIT)
			return
		}
		if captain.IsMaxAbility() {
			logrus.Debugf("武将成长，已经最高成长值了")
			result.Add(military.ERR_CAPTAIN_ENHANCE_FAIL_ABILITY_MAX)
			return
		}
		goodsLen := len(proto.GoodsId)
		if goodsLen <= 0 {
			logrus.Debugf("武将成长，物品id个数==0")
			result.Add(military.ERR_CAPTAIN_ENHANCE_FAIL_INVALID_GOODS)
			return
		}
		if goodsLen != len(proto.Count) {
			logrus.Debugf("武将成长，物品id和count的个数不一致")
			result.Add(military.ERR_CAPTAIN_ENHANCE_FAIL_INVALID_GOODS)
			return
		}
		goodsCountMap := make(map[*goods.GoodsData]uint64, goodsLen)
		for i, v := range proto.GoodsId {
			goodsId := u64.FromInt32(v)
			goods := m.datas.GetGoodsData(goodsId)
			if goods == nil {
				logrus.Debugf("武将成长，物品不存在")
				result.Add(military.ERR_CAPTAIN_ENHANCE_FAIL_INVALID_GOODS)
				return
			}
			if goods.GoodsEffect == nil || goods.GoodsEffect.ExpType != shared_proto.GoodsExpEffectType_EXP_CAPTAIN_REFINED {
				logrus.Debugf("武将成长，物品不是加武将强化经验的")
				result.Add(military.ERR_CAPTAIN_ENHANCE_FAIL_INVALID_GOODS)
				return
			}
			count := u64.FromInt32(proto.Count[i])
			if count <= 0 {
				logrus.Debugf("武将成长，物品个数无效")
				result.Add(military.ERR_CAPTAIN_ENHANCE_FAIL_INVALID_COUNT)
				return
			}
			goodsCountMap[goods] += count
		}
		heroDepot := hero.Depot()
		for goods, count := range goodsCountMap {
			if !heroDepot.HasEnoughGoods(goods.Id, count) {
				logrus.Debugf("武将成长，物品个数不足")
				result.Add(military.ERR_CAPTAIN_ENHANCE_FAIL_INVALID_COUNT)
				return
			}
		}
		hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCaptainRefined)
		toAddExp := uint64(0)
		for goods, count := range goodsCountMap {
			heromodule.ReduceGoodsAnyway(hctx, hero, result, goods, count)
			toAddExp += goods.GoodsEffect.Exp * count
		}
		ctime := m.timeService.CurrentTime()
		heromodule.AddCaptainAbilityExp(hctx, hero, result, captain, toAddExp, ctime)
		result.Add(military.NewS2cCaptainEnhanceMsg(u64.Int32(captainId), u64.Int32(captain.Ability()), u64.Int32(captain.AbilityExp()), int32(captain.Quality())))
		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_CaptainRefinedTimes)

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_REFINED_TIMES)
		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_REFINED_CAPTAIN) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_REFINED_CAPTAIN)))
		}

		result.Changed()
		result.Ok()
		return
	})
}

//gogen:iface
func (m *MilitaryModule) ProcessCaptainRevirthPreview(proto *military.C2SCaptainRebirthPreviewProto, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		cid := u64.FromInt32(proto.Id)
		captain := hero.Military().Captain(cid)
		if captain == nil {
			logrus.Debugf("武将转生预览，武将不存在，%v", proto.Id)
			result.Add(military.ERR_CAPTAIN_REBIRTH_PREVIEW_FAIL_INVALID_ID)
			return
		}

		nextRebirthLevel := captain.RebirthLevelData().GetNextLevel()
		if nextRebirthLevel == nil {
			logrus.Debugf("武将转生预览，武将已转生到最高等级")
			result.Add(military.ERR_CAPTAIN_REBIRTH_PREVIEW_FAIL_MAX_LEVEL)
			return
		}

		// 计算出新的成长值
		newAbilityData, _ := entity.UpgradeAbility(captain.AbilityData(), captain.AbilityExp()+nextRebirthLevel.AbilityExp)
		newAbility := newAbilityData.Ability

		levelSoldierCapcity := u64.Sub(nextRebirthLevel.FirstCaptainLevel.SoldierCapcity, captain.LevelData().SoldierCapcity)
		rebirthSoldierCapcity := u64.Sub(nextRebirthLevel.SoldierCapcity, captain.RebirthLevelData().SoldierCapcity)
		newSoldierCapcity := captain.SoldierCapcity() + levelSoldierCapcity + rebirthSoldierCapcity

		// 原来的成长值
		oldStatPoint := entity.CalculateCaptainStatPoint(captain.Ability(), captain.Level(), captain.RebirthLevelData(), captain.GetRarityStarCoef())
		newStatPoint := entity.CalculateCaptainStatPoint(newAbility, 1, nextRebirthLevel, captain.GetRarityStarCoef())

		diffStatPoint := u64.Sub(newStatPoint, oldStatPoint)
		toAddStatProto := data.New4DStatProto(entity.CalculateCaptain4DStatPoint(diffStatPoint, captain.Race()))

		result.Add(military.NewS2cCaptainRebirthPreviewMsg(
			proto.Id,
			nil,
			u64.Int32(nextRebirthLevel.Level),
			int32(newAbilityData.Quality),
			u64.Int32(newAbility),
			u64.Int32(nextRebirthLevel.AbilityLimit),
			u64.Int32(nextRebirthLevel.SpriteStatPoint),
			u64.Int32(newSoldierCapcity),
			must.Marshal(toAddStatProto),
		))
		result.Ok()
		return
	})
}

//gogen:iface
func (m *MilitaryModule) ProcessCaptainProgress(proto *military.C2SCaptainProgressProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCaptainRebirth)
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		cid := u64.FromInt32(proto.Id)
		captain := hero.Military().Captain(cid)
		if captain == nil {
			logrus.Debugf("武将转生，武将不存在，%v", proto.Id)
			result.Add(military.ERR_CAPTAIN_PROGRESS_FAIL_INVALID_ID)
			return
		}
		nextRebirthData := captain.NextRebirthData()
		if nextRebirthData == nil {
			logrus.Debugf("武将转生，武将已转生到最高等级")
			result.Add(military.ERR_CAPTAIN_PROGRESS_FAIL_MAX_LEVEL)
			return
		}
		if nextRebirthData.HeroLevelLimit > hero.Level() {
			logrus.Debugf("武将转生，君主等级不足")
			result.Add(military.ERR_CAPTAIN_PROGRESS_FAIL_LEVEL_NOT_ENOUGH)
			return
		}
		if !captain.NextLevelIsReBirth() {
			logrus.Debugf("武将转生，未达到转生等级")
			result.Add(military.ERR_CAPTAIN_PROGRESS_FAIL_LEVEL_NOT_ENOUGH)
			return
		}
		ctime := m.timeService.CurrentTime()
		if captain.IsInRebrithing(ctime) {
			if !proto.Miao {
				logrus.Debugf("武将转生，还在转生 CD 中")
				result.Add(military.ERR_CAPTAIN_PROGRESS_FAIL_IN_REBIRTHING)
				return
			}
			d := m.datas.MiscConfig().MiaoCaptainRebirthDuration
			if d <= 0 {
				logrus.Debugln("秒武将转生cd，功能未开放")
				result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_NOT_IN_REBIRTH)
				return
			}
			endTime := captain.RebirthEndTime()
			multi := int64((endTime.Sub(ctime) + d - 1) / d)
			if multi <= 0 {
				logrus.Errorln("秒武将转生cd，计算出来的multi <= 0")
				result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_COST_NOT_ENOUGH)
				return
			}
			cost := m.datas.MiscConfig().MiaoCaptainRebirthCost.Multiple(u64.FromInt64(multi))
			if !heromodule.TryReduceCost(hctx, hero, result, cost) {
				logrus.Debugln("秒武将转生cd，资源不足")
				result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_COST_NOT_ENOUGH)
				return
			}
			captain.SetRebirthEndTime(ctime)
		}
		oldLevel := captain.Level()
		oldAbility := captain.Ability()
		// 转生
		if !captain.ProcessRebirth(ctime) {
			logrus.Debugf("武将转生，前面验证通过还失败")
			result.Add(military.ERR_CAPTAIN_PROGRESS_FAIL_LEVEL_NOT_ENOUGH)
			return
		}
		// 转生赠送成长经验
		heromodule.AddCaptainAbilityExp(hctx, hero, result, captain, nextRebirthData.AbilityExp, ctime)
		captain.CalculateProperties()
		var toReduce uint64
		if captain.Soldier() > captain.SoldierCapcity() {
			// 减掉一些士兵
			toReduce := u64.Sub(captain.Soldier(), captain.SoldierCapcity())
			captain.ReduceSoldier(toReduce)
		}
		heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)
		if toReduce > 0 {
			ctime := m.timeService.CurrentTime()
			hero.Military().AddFreeSoldier(toReduce, ctime)
			result.Add(hero.Military().NewUpdateFreeSoldierMsg())
		}
		result.Add(military.NewS2cCaptainProgressMsg(
			proto.Id,
			u64.Int32(captain.RebirthLevel()),
			0,
			int32(captain.Quality()),
			u64.Int32(captain.Ability()),
			u64.Int32(captain.AbilityExp()),
			u64.Int32(nextRebirthData.AbilityLimit),
			u64.Int32(captain.Soldier()),
			u64.Int32(captain.SoldierCapcity()),
			must.Marshal(captain.GetTotalStat()),
			u64.Int32(captain.FightAmount()),
			u64.Int32(captain.FullSoldierFightAmount())))
		// 系统广播
		hctx = heromodule.NewContext(m.dep, operate_type.MilitaryCaptainRebirth)
		if d := hctx.BroadcastHelp().CaptainReBrith; d != nil {
			hctx.AddBroadcast(d, hero, result, 0, captain.RebirthLevel(), func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, captain.RebirthLevel())
				return text
			})
		}

		result.Changed()
		result.Ok()
		// tlog
		hctx.Tlog().TlogPlayerCultivateFlow(hero, captain.Id(), operate_type.CaptainOperTypeUpgrade, oldLevel, captain.Level(), oldAbility, captain.Ability(), u64.FromInt(captain.RaceId()), u64.FromInt(captain.RaceId()), captain.OfficialId(), captain.OfficialId(), hctx.OperId())
		return
	})
}

//gogen:iface
func (m *MilitaryModule) ProcessCaptainRebirthMiaoCd(proto *military.C2SCaptainRebirthMiaoCdProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCaptainRebirthMiaoCd)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captainId := u64.FromInt32(proto.Id)
		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("武将转生，元宝减CD，武将不存在:%s", captainId)
			result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_INVALID_ID)
			return
		}

		nextLevel := captain.RebirthLevelData().GetNextLevel()
		if nextLevel == nil {
			logrus.Debugln("已经转生到顶级了，还触发了转生元宝减时间逻辑")
			return
		}

		endTime := captain.RebirthEndTime()
		ctime := m.timeService.CurrentTime()
		if !captain.IsInRebrithing(ctime) {
			logrus.Debugln("武将转生，元宝减CD，没在转生 CD 中。")
			result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_NOT_IN_REBIRTH)
			return
		}

		d := m.datas.MiscConfig().MiaoCaptainRebirthDuration
		if d <= 0 {
			logrus.Debugln("秒武将转生cd，功能未开放")
			result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_NOT_IN_REBIRTH)
			return
		}

		multi := int64((endTime.Sub(ctime) + d - 1) / d)
		if multi <= 0 {
			logrus.Errorln("秒武将转生cd，计算出来的multi <= 0")
			result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_COST_NOT_ENOUGH)
			return
		}

		cost := m.datas.MiscConfig().MiaoCaptainRebirthCost.Multiple(u64.FromInt64(multi))

		if !heromodule.TryReduceCost(hctx, hero, result, cost) {
			logrus.Debugln("秒武将转生cd，资源不足")
			result.Add(military.ERR_CAPTAIN_REBIRTH_MIAO_CD_FAIL_COST_NOT_ENOUGH)
			return
		}

		captain.SetRebirthEndTime(ctime)
		result.Add(military.NewS2cCaptainRebirthMiaoCdMsg(u64.Int32(captain.Id())))

		result.Changed()
		result.Ok()
	})
}

//gogen:iface c2s_collect_captain_training_exp
func (m *MilitaryModule) ProcessCollectCaptainTrainingExp(hc iface.HeroController) {

	// 一键领取修炼经验
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCollectCaptainTrainingExp)

		building := hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN)
		if building == nil {
			logrus.Debug("一键领取修炼馆经验，修炼馆未开放")
			result.Add(military.ERR_COLLECT_CAPTAIN_TRAINING_EXP_FAIL_BUILDING_LOCKED)
			return
		}

		ctime := m.timeService.CurrentTime()
		trainingExp := heromodule.CalcTrainAllExp(hero, m.datas, ctime)

		var captainIds, oldCaptainLevels, newCaptainLevels []uint64
		logrus.Debugf("修炼馆领取武将经验：%v", trainingExp)

		for _, t := range hero.Troops() {
			if t == nil {
				continue
			}

			for _, pos := range t.Pos() {
				c := pos.Captain()
				if c == nil {
					continue
				}

				toAddExp := trainingExp + c.GetTrainAccExp()
				if toAddExp > 0 {
					captainIds = append(captainIds, c.Id())
					oldCaptainLevels = append(oldCaptainLevels, c.Level())
					heromodule.AddCaptainExp(hctx, hero, result, c, toAddExp, ctime) // 加经验
					newCaptainLevels = append(newCaptainLevels, c.Level())

					// 清空数据
					c.SetTrainAccExp(0)
				}
			}
		}

		hero.Military().SetGlobalTrainStartTime(ctime)
		hero.Military().ClearReservedExp()

		result.Add(military.NewS2cCollectCaptainTrainingExpMsg(timeutil.Marshal32(ctime)))

		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_COLLECT_XIU_LIAN_EXP)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_COLLECT_TRAIN_EXP) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_COLLECT_TRAIN_EXP)))
		}

		for i, id := range captainIds {
			m.dep.Tlog().TlogPlayerExpDrugFlow(hero, id, oldCaptainLevels[i], newCaptainLevels[i])
		}
	})
}

// 废弃，跟 collect_captain_training_exp 一模一样
//gogen:iface c2s_captain_train_exp
func (m *MilitaryModule) ProcessCaptainTrainExp(hc iface.HeroController) {
	hc.Send(military.ERR_CAPTAIN_TRAIN_EXP_FAIL_BUILDING_LOCKED)
}

//gogen:iface c2s_captain_can_collect_exp
func (m *MilitaryModule) ProcessCaptainCanCollectExp(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN) == nil {
			result.Add(military.ERR_CAPTAIN_CAN_COLLECT_EXP_FAIL_UNLOCK)
			return
		}

		var coefs []int32
		hero.Buff().Walk(func(buff *entity.BuffInfo) {
			if buff.EffectData.EffectType == shared_proto.BuffEffectType_Buff_ET_captain_train {
				coefs = append(coefs, u64.Int32(buff.EffectData.CaptainTrain.Percent))
			}
		})

		allExp := heromodule.CalcTrainAllExp(hero, m.datas, m.timeService.CurrentTime())
		dur := heromodule.GetTrainingMaxDuration(hero, m.datas)
		result.Add(military.NewS2cCaptainCanCollectExpMsg(u64.Int32(allExp), coefs, timeutil.DurationMarshal32(dur)))
	})
}

// 新新版武将使用经验书升级（使用1个道具）
//gogen:iface
func (m *MilitaryModule) ProcessUseLevelExpGoods2(proto *military.C2SUseLevelExpGoods2Proto, hc iface.HeroController) {
	// 新新版武将使用经验书升级，经验书道具不再和修炼馆等级挂钩，而是直接配固定经验。 有4个档次的经验书 可配置。
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		building := hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN)
		if building == nil {
			logrus.Debug("使用经验书，修炼馆未开放")
			result.Add(military.ERR_USE_LEVEL_EXP_GOODS2_FAIL_BUILDING_LOCKED)
			return
		}

		captainId := u64.FromInt32(proto.Captain)
		c := hero.Military().Captain(captainId)
		if c == nil {
			logrus.Debugf("使用经验书，武将不存在")
			result.Add(military.ERR_USE_LEVEL_EXP_GOODS2_FAIL_NO_CAPTAIN)
			return
		}
		if c.IsMaxLevel() {
			logrus.Debug("使用经验书，武将已达最大等级")
			result.Add(military.ERR_USE_LEVEL_EXP_GOODS2_FAIL_CAPTAIN_MAX_LEVEL)
			return
		}
		if c.IsLevelLimit(hero.LevelData().Sub.CaptainLevelLimit, true) {
			logrus.Debug("使用经验书，武将等级受限，请提升君主等级")
			result.Add(military.ERR_USE_LEVEL_EXP_GOODS2_FAIL_CAPTAIN_LEVEL_LIMIT)
			return
		}
		ctime := m.timeService.CurrentTime()
		// 转生 CD 中，不能加经验
		if c.IsInRebrithing(ctime) {
			logrus.Debug("使用经验书，武将正在转生中")
			result.Add(military.ERR_USE_LEVEL_EXP_GOODS2_FAIL_IN_REBIRTHING)
			return
		}
		maxCanAddExp := c.GetMaxCanAddExp(hero.LevelData().Sub.CaptainLevelLimit)
		if maxCanAddExp <= 0 {
			logrus.Debug("使用经验书，已经不能再加经验了")
			result.Add(military.ERR_USE_LEVEL_EXP_GOODS2_FAIL_CAPTAIN_LEVEL_LIMIT)
			return
		}
		goodsId := u64.FromInt32(proto.GoodsId)
		goods := m.datas.GetGoodsData(goodsId)
		if goods == nil {
			logrus.Debugf("使用经验书，物品不存在")
			result.Add(military.ERR_USE_LEVEL_EXP_GOODS2_FAIL_NOT_EXP_GOODS)
			return
		}
		if goods.GoodsEffect == nil || goods.GoodsEffect.ExpType != shared_proto.GoodsExpEffectType_EXP_CAPTAIN {
			logrus.Debugf("使用经验书，物品不是加武将经验的")
			result.Add(military.ERR_USE_LEVEL_EXP_GOODS2_FAIL_NOT_EXP_GOODS)
			return
		}
		if !hero.Depot().HasEnoughGoods(goods.Id, 1) {
			logrus.Debugf("使用经验书，经验书不足")
			result.Add(military.ERR_USE_LEVEL_EXP_GOODS2_FAIL_NOT_ENOUTH_GOODS)
			return
		}
		hctx := heromodule.NewContext(m.dep, operate_type.MilitaryUseTrainingExpGoods)
		heromodule.ReduceGoodsAnyway(hctx, hero, result, goods, 1)

		upgrade := heromodule.AddCaptainExp(hctx, hero, result, c, goods.GoodsEffect.Exp, ctime) // 加经验
		result.Add(military.NewS2cUseLevelExpGoods2Msg(proto.Captain, u64.Int32(c.Level()), u64.Int32(c.Exp()), upgrade))

		hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CaptainExpGoodsUsed, 1)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_EXP_GOODS_USE)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_USE_TU_FEI_GOODS) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_USE_TU_FEI_GOODS)))

			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_USE_TU_FEI_GOODS)
		}

		result.Changed()
		result.Ok()
		return
	})
}

// 新新版武将使用经验书升级（升1级），自动扣除经验书
//gogen:iface
func (m *MilitaryModule) ProcessAutoUseGoodsUntilCaptainLevelup(proto *military.C2SAutoUseGoodsUntilCaptainLevelupProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		building := hero.Domestic().GetBuilding(shared_proto.BuildingType_XIU_LIAN_GUAN)
		if building == nil {
			logrus.Debug("武将升一级，修炼馆未开放")
			result.Add(military.ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_BUILDING_LOCKED)
			return
		}

		captainId := u64.FromInt32(proto.Captain)
		c := hero.Military().Captain(captainId)
		if c == nil {
			logrus.Debugf("武将升一级，武将不存在")
			result.Add(military.ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_NO_CAPTAIN)
			return
		}
		if c.IsMaxLevel() {
			logrus.Debug("武将升一级，武将已达最大等级")
			result.Add(military.ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_CAPTAIN_MAX_LEVEL)
			return
		}
		if c.IsLevelLimit(hero.LevelData().Sub.CaptainLevelLimit, true) {
			logrus.Debug("武将升一级，武将等级受限，请提升君主等级")
			result.Add(military.ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_CAPTAIN_LEVEL_LIMIT)
			return
		}
		ctime := m.timeService.CurrentTime()
		// 转生 CD 中，不能加经验
		if c.IsInRebrithing(ctime) {
			logrus.Debug("武将升一级，武将正在转生中")
			result.Add(military.ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_IN_REBIRTHING)
			return
		}
		maxCanAddExp := c.GetMaxCanAddExp(u64.Min(hero.LevelData().Sub.CaptainLevelLimit, c.Level()+1))
		if maxCanAddExp <= 0 {
			logrus.Debug("武将升一级，已经不能再加经验了")
			result.Add(military.ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_CAPTAIN_LEVEL_LIMIT)
			return
		}
		// 记录所有经验书可以提升的总经验
		var totalAddExp uint64
		// 记录背包中当前所有数量
		var useGoods []*goods.GoodsData
		var useCounts []uint64
		var totalUseCount uint64
		for _, goods := range m.datas.GoodsConfig().TrainExpGoods4CaptainUpgrade {
			if count := hero.Depot().GetGoodsCount(goods.Id); count > 0 {

				canAddExp := u64.Sub(maxCanAddExp, totalAddExp)
				useCount := u64.DivideTimes(canAddExp, goods.GoodsEffect.Exp)
				useCount = u64.Min(useCount, count)

				if useCount > 0 {
					totalUseCount += useCount
					totalAddExp += goods.GoodsEffect.Exp * useCount

					useGoods = append(useGoods, goods)
					useCounts = append(useCounts, useCount)

					if totalAddExp >= maxCanAddExp {
						break
					}
				}
			}
		}

		if totalAddExp <= 0 {
			logrus.Debug("武将升一级，没有任何经验书")
			result.Add(military.ERR_AUTO_USE_GOODS_UNTIL_CAPTAIN_LEVELUP_FAIL_NO_EXP_GOODS)
			return
		}
		hctx := heromodule.NewContext(m.dep, operate_type.MilitaryUseTrainingExpGoods)

		// 扣物品
		for i, goods := range useGoods {
			if count := useCounts[i]; count > 0 {
				heromodule.ReduceGoodsAnyway(hctx, hero, result, goods, count)
			}
		}

		// 加经验
		upgrade := heromodule.AddCaptainExp(hctx, hero, result, c, totalAddExp, ctime) // 加经验
		result.Add(military.NewS2cAutoUseGoodsUntilCaptainLevelupMsg(proto.Captain, u64.Int32(c.Level()), u64.Int32(c.Exp()), upgrade))

		hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CaptainExpGoodsUsed, totalUseCount)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_EXP_GOODS_USE)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_USE_TU_FEI_GOODS) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_USE_TU_FEI_GOODS)))

			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_USE_TU_FEI_GOODS)
		}

		result.Changed()
		result.Ok()
		return
	})
}

//gogen:iface c2s_jiu_guan_consult
func (m *MilitaryModule) ProcessJiuGuanConsult(hc iface.HeroController) {
	var broadcastFunc func()

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroMilitary := hero.Military()

		building := hero.Domestic().GetBuilding(shared_proto.BuildingType_JIU_GUAN)
		if building == nil {
			logrus.Debugln("没有酒馆")
			result.Add(military.ERR_JIU_GUAN_CONSULT_FAIL_NO_JIU_GUAN)
			return
		}

		jiuGuanData := m.datas.JiuGuanData().Must(building.Level)
		if jiuGuanData == nil {
			logrus.WithField("level", building.Level).Errorln("酒馆等级数据没找到")
			result.Add(military.ERR_JIU_GUAN_CONSULT_FAIL_SERVER_BUSY)
			return
		}

		ctime := m.timeService.CurrentTime()

		jiuGuanTimes := heroMilitary.JiuGuanTimes()
		if jiuGuanTimes.Times() >= jiuGuanData.MaxTimes {
			logrus.Debugln("没有酒馆次数了")
			result.Add(military.ERR_JIU_GUAN_CONSULT_FAIL_NO_TIMES)
			return
		}

		if !m.datas.GetVipLevelData(hero.VipLevel()).JiuGuanQuickConsult {
			if ctime.Before(jiuGuanTimes.NextTime()) {
				logrus.Debugln("下次请教时间未到")
				result.Add(military.ERR_JIU_GUAN_CONSULT_FAIL_COUNTDOWN)
				return
			}
		}

		originIndex := heroMilitary.JiuGuanTutorIndex()
		tutor := jiuGuanData.MustTutorData(originIndex)

		critMulti, critImgIndex := tutor.RandomCritMulti()

		prize := tutor.Prize
		if critMulti > 1 {
			prize = resdata.NewPrizeBuilder().AddMultiple(tutor.Prize, critMulti).Build()
		}

		nextIndex := jiuGuanData.InitRefresh()
		heroMilitary.JiuGuanConsult(ctime.Add(jiuGuanData.RecoveryDuration), nextIndex) // 请教一次

		hctx := heromodule.NewContext(m.dep, operate_type.MilitaryJiuGuanConsult)
		heromodule.AddPrize(hctx, hero, result, prize, ctime)

		result.Add(military.NewS2cJiuGuanConsultMsg(must.Marshal(prize.Encode()), u64.Int32(critMulti), int32(critImgIndex), u64.Int32(originIndex), u64.Int32(nextIndex)))
		result.Add(military.NewS2cJiuGuanTimesChangedMsg(jiuGuanTimes.TimesAndNextTime()))

		if critMulti >= jiuGuanData.BroadcastMinCritMulti {
			// 广播
			broadcastMsg := military.NewS2cJiuGuanConsultBroadcastMsg(u64.Int32(jiuGuanData.Level), u64.Int32(critMulti), hero.Name()).Static()

			broadcastFunc = func() {
				m.worldService.Broadcast(broadcastMsg)
			}
		}

		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_JIU_GUAN_CONSULT)
		hero.HistoryAmount().Increase(server_proto.HistoryAmountType_Consult, 1)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_CONSULT)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_JIUGUAN_CONSULT) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_JIUGUAN_CONSULT)))
		}

		// 系统广播
		if d := hctx.BroadcastHelp().JiuGuanQingJiao; d != nil {
			hctx.AddBroadcast(d, hero, result, 0, tutor.Id, func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNpc, tutor.Name)
				return text
			})
		}
		if d := hctx.BroadcastHelp().JiuGuanBaoJi; d != nil {
			hctx.AddBroadcast(d, hero, result, 0, critMulti, func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, critMulti)
				return text
			})
		}

		result.Changed()
		result.Ok()
		return
	})

	if broadcastFunc != nil {
		broadcastFunc()
	}
}

//gogen:iface c2s_jiu_guan_refresh
func (m *MilitaryModule) ProcessJiuGuanRefresh(proto *military.C2SJiuGuanRefreshProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryJiuGuanRefresh)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroMilitary := hero.Military()

		building := hero.Domestic().GetBuilding(shared_proto.BuildingType_JIU_GUAN)
		if building == nil {
			logrus.Debugln("没有酒馆")
			result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_NO_JIU_GUAN)
			return
		}

		data := m.datas.JiuGuanData().Must(building.Level)
		if data == nil {
			logrus.WithField("level", building.Level).Errorln("酒馆等级数据没找到")
			result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_SERVER_BUSY)
			return
		}

		jiuGuanTimes := heroMilitary.JiuGuanTimes()
		if jiuGuanTimes.Times() >= data.MaxTimes {
			logrus.Debugln("酒馆刷新，没有酒馆请教次数了")
			result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_NO_TIMES)
			return
		}

		vipData := m.datas.GetVipLevelData(hero.VipLevel())
		if !vipData.JiuGuanCostRefreshInfinite && !proto.AutoMax {
			if heroMilitary.JiuGuanRefreshTimes() >= vipData.JiuGuanCostRefreshCount {
				logrus.Debugln("没有酒馆刷新次数了")
				result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_NO_REFRESH_TIMES)
				return
			}
		}

		ctime := m.timeService.CurrentTime()
		if !vipData.JiuGuanQuickConsult {
			if ctime.Before(jiuGuanTimes.NextTime()) {
				logrus.Debugln("刷新导师，下次请教倒计时未结束")
				result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_COUNTDOWN)
				return
			}
		}

		if heroMilitary.JiuGuanTutorIndex() >= uint64(len(data.TutorDatas)) {
			logrus.Debugln("已经是最好品质的导师了，不需要再换了")
			result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_TUTOR_NO_NEED)
			return
		}

		var costId int
		var newTutorIndex uint64
		isFirstRefresh := !hero.Bools().Get(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)
		if proto.AutoMax {
			tutor := data.MustTutorData(heroMilitary.JiuGuanTutorIndex())
			if tutor.RefreshMaxCost == nil {
				logrus.Debugln("刷新导师(一键最高)，未开放")
				result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_AUTO_MAX_NOT_OPEN)
				return
			}

			if !vipData.JiuGuanAutoMax {
				logrus.Debugln("刷新导师(一键最高)，vip 等级不够")
				result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_AUTO_MAX_NOT_OPEN)
				return
			}

			costId = tutor.RefreshMaxCost.Id
			// 要消耗了
			if !heromodule.TryReduceCost(hctx, hero, result, tutor.RefreshMaxCost) {
				logrus.Debugln("刷新导师(一键最高)，消耗不够")
				result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_YUANBAO_NOT_ENOUGH)
				return
			}

			newTutorIndex = uint64(len(m.datas.GetTutorDataArray()))
			newTutorIndex = u64.Sub(newTutorIndex, 1)

		} else {
			if !isFirstRefresh {
				index, consultSuc := data.Refresh(heroMilitary.JiuGuanTutorIndex())
				if !consultSuc {
					logrus.Error("已经是最好品质的导师了（前面不是已经判断过了吗）")
					result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_TUTOR_NO_NEED)
					return
				}
				newTutorIndex = index
			} else {
				newTutorIndex = m.datas.JiuGuanMiscData().FirstRefreshIndex
			}

			cost := m.datas.JiuGuanMiscData().GetRefreshCostDianquan(heroMilitary.JiuGuanRefreshTimes())
			if cost != nil {
				costId = cost.Id
				// 要消耗点券了
				if !heromodule.TryReduceCost(hctx, hero, result, cost) {
					logrus.Debugln("刷新导师，消耗不够")
					result.Add(military.ERR_JIU_GUAN_REFRESH_FAIL_YUANBAO_NOT_ENOUGH)
					return
				}
			}
		}

		if isFirstRefresh {
			hero.Bools().SetTrue(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_JIUGUAN_REFRESH)))
		}

		heroMilitary.IncreJiuGuanRefreshTimes()
		heroMilitary.SetJiuGuanTutorIndex(newTutorIndex)

		result.Add(military.NewS2cJiuGuanRefreshMsg(u64.Int32(newTutorIndex), u64.Int32(heroMilitary.JiuGuanRefreshTimes()), proto.AutoMax))

		result.Changed()
		result.Ok()

		m.dep.Tlog().TlogRefreshFlow(hero, uint64(shared_proto.BuildingType_JIU_GUAN), operate_type.RefreshMoney, u64.FromInt(costId))
		return
	})
}

var clearHomeTroopDefeatedMailMsg = military.NewS2cClearDefenseTroopDefeatedMailMsg(false).Static()

// 清掉防守失败邮件
//gogen:iface c2s_clear_oop_defeated_mail
func (m *MilitaryModule) ProcessClearDefenseTroopDefeatedMail(proto *military.C2SClearDefenseTroopDefeatedMailProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.SetTroopDefeatedMailProto(nil) {
			result.Changed()
		}

		result.Add(clearHomeTroopDefeatedMailMsg)

		result.Ok()
		return
	})
}

//gogen:iface
func (m *MilitaryModule) ProcessCaptainStatDetails(proto *military.C2SCaptainStatDetailsProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captainId := u64.FromInt32(proto.Captain)
		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("查看武将详细属性，武将不存在")
			result.Add(military.ERR_CAPTAIN_STAT_DETAILS_FAIL_CAPTAIN_NOT_FOUND)
			return
		}
		result.Add(military.NewS2cCaptainStatDetailsMarshalMsg(proto.Captain, captain.GetDetailStat()))
		result.Ok()
		return
	})
}

func (m *MilitaryModule) setCaptainOfficialAnyway(hero *entity.Hero, result herolock.LockResult, captain *entity.Captain, officialId, officialIdx uint64) {

	oldLevel := captain.Level()
	oldAbility := captain.Ability()
	oldRace := captain.RaceId()
	oldOfficial := captain.OfficialId()

	// 清掉该武将的老坑标记
	if captain.OfficialId() != 0 {
		hero.Military().SetOfficialAnyway(captain.OfficialId(), captain.OfficialIndex(), 0)
	}
	// 新的坑要是有老武将就把老武将身上的印记干掉
	if oldCaptainId, ok := hero.Military().TryGetOfficialCaptainId(officialId, officialIdx); ok {
		oldCaptain := hero.Military().Captain(oldCaptainId)
		if oldCaptain != nil {
			oldCaptain.SetOfficial(captaindata.EmptyOfficialData, 0)
			oldCaptain.CalculateProperties()
			result.Add(oldCaptain.NewUpdateCaptainStatMsg())
		}
	}
	// 给该武将身上标注新的印记
	captain.SetOfficial(m.datas.CaptainOfficialData().Get(officialId), officialIdx)
	// 给新的坑标记
	hero.Military().SetOfficialAnyway(officialId, officialIdx, captain.Id())
	captain.CalculateProperties()
	result.Add(captain.NewUpdateCaptainStatMsg())

	// 减掉一些士兵
	var toReduce uint64
	if captain.Soldier() > captain.SoldierCapcity() {
		toReduce := u64.Sub(captain.Soldier(), captain.SoldierCapcity())
		captain.ReduceSoldier(toReduce)
	}
	if toReduce > 0 {
		ctime := m.timeService.CurrentTime()
		hero.Military().AddFreeSoldier(toReduce, ctime)
		result.Add(hero.Military().NewUpdateFreeSoldierMsg())
	}
	heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)
	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_OFFICIAL_UPDATE)
	if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_CAPTAIN_UPGRADE_OFFICIAL) {
		result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_CAPTAIN_UPGRADE_OFFICIAL)))
	}
	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryUpdateCaptainOfficial)
	// tlog
	hctx.Tlog().TlogPlayerCultivateFlow(hero, captain.Id(), operate_type.CaptainOperTypeOfficial, oldLevel, captain.Level(), oldAbility, captain.Ability(), u64.FromInt(oldRace), u64.FromInt(captain.RaceId()), oldOfficial, captain.OfficialId(), hctx.OperId())
}

func (m *MilitaryModule) cancelCaptainOfficialAnyway(hero *entity.Hero, result herolock.LockResult, captain *entity.Captain) {

	oldLevel := captain.Level()
	oldAbility := captain.Ability()
	oldRace := captain.RaceId()
	oldOfficial := captain.OfficialId()

	hero.Military().SetOfficialAnyway(captain.OfficialId(), captain.OfficialIndex(), 0)
	captain.SetOfficial(captaindata.EmptyOfficialData, 0)
	captain.CalculateProperties()
	result.Add(captain.NewUpdateCaptainStatMsg())

	// 减掉一些士兵
	var toReduce uint64
	if captain.Soldier() > captain.SoldierCapcity() {
		toReduce := u64.Sub(captain.Soldier(), captain.SoldierCapcity())
		captain.ReduceSoldier(toReduce)
	}
	if toReduce > 0 {
		ctime := m.timeService.CurrentTime()
		hero.Military().AddFreeSoldier(toReduce, ctime)
		result.Add(hero.Military().NewUpdateFreeSoldierMsg())
	}
	heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)

	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryLeaveCaptainOfficial)
	// tlog
	hctx.Tlog().TlogPlayerCultivateFlow(hero, captain.Id(), operate_type.CaptainOperTypeOfficial, oldLevel, captain.Level(), oldAbility, captain.Ability(), u64.FromInt(oldRace), u64.FromInt(captain.RaceId()), oldOfficial, captain.OfficialId(), hctx.OperId())

}

//gogen:iface
func (m *MilitaryModule) ProcessSetCaptainOfficial(proto *military.C2SSetCaptainOfficialProto, hc iface.HeroController) {

	captainIds := proto.GetCaptain()
	officialIds := proto.GetOfficial()
	officialIndexs := proto.GetOfficialIdx()
	arrLen := len(captainIds)
	if arrLen <= 0 || arrLen != len(officialIds) || arrLen != len(officialIndexs) {
		logrus.Debugf("武将封官(卸任)，数组长度错误")
		hc.Send(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_ERR_DATA)
		return
	}
	if check.Int32AnyZero(captainIds) {
		logrus.Debugf("武将封官(卸任)，武将ID有0")
		hc.Send(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_INVALID_CAPTAIN)
		return
	}
	if check.Int32Duplicate(captainIds) {
		logrus.Debugf("武将封官(卸任)，武将ID有重复")
		hc.Send(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_ERR_DATA)
		return
	}
	if check.Int32AnyLt0(officialIds) {
		logrus.Debugf("武将封官(卸任)，官职ID有负数")
		hc.Send(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_ERR_DATA)
		return
	}
	if check.Int32AnyLt0(officialIndexs) {
		logrus.Debugf("武将封官(卸任)，官职位置下标有负数")
		hc.Send(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_ERR_DATA)
		return
	}
	if arrLen > 1 {
		arrIdxMap := make(map[int32][]int32)
		for i := 0; i < arrLen; i++ {
			arrIdx := arrIdxMap[officialIds[i]]
			arrIdx = append(arrIdx, officialIndexs[i])
			arrIdxMap[officialIds[i]] = arrIdx
		}
		for offid, arr := range arrIdxMap {
			if len(arr) < 2 {
				continue
			}
			if check.Int32Duplicate(arr) {
				logrus.Debugf("武将封官(卸任)，官职 official id:%v 的位置下标有重复", offid)
				hc.Send(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_ERR_DATA)
				return
			}
		}
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captains := make([]*entity.Captain, arrLen)
		for i := 0; i < arrLen; i++ {
			captainId := u64.FromInt32(captainIds[i])
			captain := hero.Military().Captain(captainId)
			if captain == nil {
				logrus.Debugf("武将封官(卸任)，武将不存在 id:%v", captainId)
				result.Add(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_INVALID_CAPTAIN)
				return
			}
			officialId := u64.FromInt32(officialIds[i])
			if officialId != 0 { // 是上任的
				official := m.datas.CaptainOfficialData().Get(officialId)
				if official == nil {
					logrus.Debugf("武将封官(卸任)，官职不存在 official id:%v", officialId)
					result.Add(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_INVALID_OFFICIAL)
					return
				}
				if captain.GongXun() < official.NeedGongxun {
					logrus.Debugf("武将封官(卸任)，功勋不够")
					result.Add(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_GONGXUN_NOT_ENOUGH)
					return
				}
				officialIdx := u64.FromInt32(officialIndexs[i])
				if cid, ok := hero.Military().TryGetOfficialCaptainId(officialId, officialIdx); !ok {
					logrus.Debugf("武将封官(卸任)，超出官职位置下标")
					result.Add(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_INVALID_OFFICIAL)
					return
				} else if cid == captainId {
					logrus.Debugf("武将封官(卸任)，重复设置")
					result.Add(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_ALREADY_ON_OFFICIAL)
					return
				}
				// 如果降职
				if officialId < captain.OfficialId() && captain.IsOutSide() {
					logrus.Debugf("武将封官(卸任)，武将在外面")
					result.Add(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_CAPTAIN_IS_OUTSIDE)
					return
				}
			} else { // 是卸任的
				if captain.OfficialId() == 0 {
					logrus.Debugf("武将封官(卸任)，武将已经没有官职了")
					result.Add(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_ALREADY_NO_OFFICIAL)
					return
				}
				if captain.IsOutSide() {
					logrus.Debugf("武将封官(卸任)，武将在外面")
					result.Add(military.ERR_SET_CAPTAIN_OFFICIAL_FAIL_CAPTAIN_IS_OUTSIDE)
					return
				}
			}
			captains[i] = captain
		}

		// 开始匹量处理
		for i := 0; i < arrLen; i++ {
			if officialIds[i] != 0 {
				m.setCaptainOfficialAnyway(hero, result, captains[i], u64.FromInt32(officialIds[i]), u64.FromInt32(officialIndexs[i]))
			} else {
				m.cancelCaptainOfficialAnyway(hero, result, captains[i])
			}
		}

		result.Add(military.NewS2cSetCaptainOfficialMsg(captainIds, officialIds, officialIndexs))

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *MilitaryModule) ProcessUseGongXunGoods(proto *military.C2SUseGongXunGoodsProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryGongXunGoods)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		captain := hero.Military().Captain(u64.FromInt32(proto.Captain))
		if captain == nil {
			logrus.Debug("使用功勋物品，无效的武将id")
			result.Add(military.ERR_USE_GONG_XUN_GOODS_FAIL_INVALID_CAPTAIN)
			return
		}

		maxGongXunOfficialData := hero.LevelData().GetCaptainOfficialCountData().GetMaxGongXunOfficialData()
		if maxGongXunOfficialData == nil {
			logrus.Debug("使用功勋物品，maxGongXunOfficialData == nil")
			result.Add(military.ERR_USE_GONG_XUN_GOODS_FAIL_OFFICIAL_LIMIT)
			return
		}

		canAddGongXun := u64.Sub(maxGongXunOfficialData.NeedGongxun, captain.GongXun())
		if canAddGongXun <= 0 {
			logrus.Debug("使用功勋物品，已达最高功勋职位")
			result.Add(military.ERR_USE_GONG_XUN_GOODS_FAIL_OFFICIAL_LIMIT)
			return
		}

		var toAddExp uint64 // 加的功勋

		goodsCountMap := make(map[*goods.GoodsData]uint64, len(proto.GoodsId))
		for i, v := range proto.GoodsId {
			goodsId := u64.FromInt32(v)
			goods := m.datas.GetGoodsData(goodsId)
			if goods == nil {
				logrus.Debugf("使用功勋物品，物品不存在")
				result.Add(military.ERR_USE_GONG_XUN_GOODS_FAIL_INVALID_GOODS)
				return
			}

			if goods.GoodsEffect == nil ||
				goods.GoodsEffect.ExpType != shared_proto.GoodsExpEffectType_EXP_GONG_XUN {
				logrus.Debugf("使用功勋物品，物品不是加武将功勋的")
				result.Add(military.ERR_USE_GONG_XUN_GOODS_FAIL_INVALID_GOODS)
				return
			}

			count := u64.FromInt32(proto.Count[i])
			if count <= 0 {
				logrus.Debugf("使用功勋物品，物品个数无效")
				result.Add(military.ERR_USE_GONG_XUN_GOODS_FAIL_INVALID_COUNT)
				return
			}

			if toAddExp+goods.GoodsEffect.Exp*count >= canAddGongXun {
				sub := u64.Sub(canAddGongXun, toAddExp)
				useCount := u64.DivideTimes(sub, goods.GoodsEffect.Exp)

				goodsCountMap[goods] += useCount
				toAddExp += goods.GoodsEffect.Exp * useCount
				break
			}

			goodsCountMap[goods] += count
			toAddExp += goods.GoodsEffect.Exp * count
		}

		heroDepot := hero.Depot()
		for goods, count := range goodsCountMap {
			if !heroDepot.HasEnoughGoods(goods.Id, count) {
				logrus.Debugf("使用功勋物品，物品个数不足")
				result.Add(military.ERR_USE_GONG_XUN_GOODS_FAIL_INVALID_COUNT)
				return
			}
		}

		// 扣物品
		for goods, count := range goodsCountMap {
			// 开始操作，扣物品，加经验，升级
			heromodule.ReduceGoodsAnyway(hctx, hero, result, goods, count)
		}

		// 给武将加功勋
		newGongXun := captain.AddGongXun(toAddExp)

		result.Add(military.NewS2cUseGongXunGoodsMsg(proto.Captain, u64.Int32(newGongXun)))

		result.Changed()
		result.Ok()
	})

}

//gogen:iface
func (m *MilitaryModule) ProcessCloseFightGuide(proto *military.C2SCloseFightGuideProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		result.Add(military.NewS2cCloseFightGuideMsg(proto.Close))

		if proto.Close {
			hero.Bools().SetTrue(shared_proto.HeroBoolType_BOOL_CLOSE_FIGHT_GUIDE)
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_CLOSE_FIGHT_GUIDE)))
		} else {
			hero.Bools().SetFalse(shared_proto.HeroBoolType_BOOL_CLOSE_FIGHT_GUIDE)
			result.Add(misc.NewS2cResetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_CLOSE_FIGHT_GUIDE)))
		}

		result.Changed()
		result.Ok()
	})

}

//gogen:iface
func (m *MilitaryModule) ProcessViewOtherHeroCaptain(proto *military.C2SViewOtherHeroCaptainProto, hc iface.HeroController) {

	heroId, ok := idbytes.ToId(proto.HeroId)
	if !ok {
		logrus.Debug("查看其它玩家武将，英雄id无效")
		hc.Send(military.ERR_VIEW_OTHER_HERO_CAPTAIN_FAIL_INVALID_HERO_ID)
		return
	}

	var toSend pbutil.Buffer
	if m.dep.HeroData().FuncNotError(heroId, func(hero *entity.Hero) (heroChanged bool) {

		captain := hero.Military().Captain(u64.FromInt32(proto.CaptainId))
		if captain == nil {
			logrus.Debug("查看其它玩家武将，captain == nil")
			toSend = military.ERR_VIEW_OTHER_HERO_CAPTAIN_FAIL_NOT_FOUND
			return
		}

		toSend = military.NewS2cViewOtherHeroCaptainMsg(hero.IdBytes(), hero.Name(), must.Marshal(captain.EncodeOtherView()))
		return
	}) {
		hc.Send(military.ERR_VIEW_OTHER_HERO_CAPTAIN_FAIL_INVALID_HERO_ID)
		return
	}

	if toSend == nil {
		toSend = military.ERR_VIEW_OTHER_HERO_CAPTAIN_FAIL_NOT_FOUND
	}

	hc.Send(toSend)
}

//gogen:iface
func (m *MilitaryModule) ProcessCopyDefenserGoods(proto *military.C2SUseCopyDefenserGoodsProto, hc iface.HeroController) {

	goodsData := m.datas.GetGoodsData(u64.FromInt32(proto.Goods))
	if goodsData == nil {
		logrus.Debug("使用复制驻防物品，物品id无效")
		hc.Send(military.ERR_USE_COPY_DEFENSER_GOODS_FAIL_INVALID_GOODS)
		return
	}

	if goodsData.GoodsEffect == nil || goodsData.GoodsEffect.DurationType != shared_proto.GoodsDurationEffectType_COPY_DEFENSER {
		logrus.Debug("使用复制驻防物品，物品不是复制驻防的物品")
		hc.Send(military.ERR_USE_COPY_DEFENSER_GOODS_FAIL_INVALID_GOODS)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCopyDefenserGoods)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		t := hero.GetTroopByIndex(u64.FromInt32(proto.TroopIndex - 1))
		if t == nil {
			logrus.Debug("使用复制驻防物品，队伍序号无效")
			result.Add(military.ERR_USE_COPY_DEFENSER_GOODS_FAIL_INVALID_TROOP)
			return
		}

		captains := make([]*entity.CaptainIdSoldier, len(t.Pos()))
		captainCount := 0
		for i, c := range t.Pos() {
			if c.Captain() == nil {
				continue
			}

			soldier := c.Captain().SoldierCapcity()
			captainStat := c.Captain().GetCopyDefenseStat(hero)
			cis := entity.NewCaptainIdSoldier(c.Captain().Id(), c.XIndex(), captainStat.Encode(), soldier, soldier, c.Captain().GetSpellFightAmountCoef())
			captains[i] = cis

			captainCount ++
		}

		if captainCount <= 0 {
			logrus.Debug("使用复制驻防物品，captainCount <= 0")
			result.Add(military.ERR_USE_COPY_DEFENSER_GOODS_FAIL_INVALID_TROOP)
			return
		}

		if !heromodule.TryReduceOrBuyGoods(hctx, hero, result, goodsData, 1, proto.AutoBuy) {
			if proto.AutoBuy {
				logrus.Debug("使用复制驻防物品，自动购买，钱不够")
				result.Add(military.ERR_USE_COPY_DEFENSER_GOODS_FAIL_COST_NOT_ENOUGH)
			} else {
				logrus.Debug("使用复制驻防物品，物品个数不足")
				result.Add(military.ERR_USE_COPY_DEFENSER_GOODS_FAIL_COUNT_NOT_ENOUGH)
			}
			return
		}

		ctime := m.timeService.CurrentTime()
		endTime := ctime.Add(goodsData.GoodsEffect.Duration)
		hero.SetCopyDefenser(t.Sequence()+1, captains, endTime)

		result.Add(military.NewS2cUseCopyDefenserGoodsMsg(proto.TroopIndex, timeutil.Marshal32(endTime)))

		result.Ok()
	})

}

func (m *MilitaryModule) GmRate(heroId int64, tutorid, times uint64) {
	// GM 测试驿站概率
	prizeBuilder := resdata.NewPrizeBuilder()
	tutor := m.datas.JiuGuanData().Must(100).MustTutorData(tutorid)
	for i := uint64(0); i < times; i++ {
		critMulti, _ := tutor.RandomCritMulti()
		prize := tutor.Prize
		if critMulti > 1 {
			prize = resdata.NewPrizeBuilder().AddMultiple(tutor.Prize, critMulti).Build()
		}
		prizeBuilder.Add(prize)
	}
	text := fmt.Sprintf("驿站 %+v", prizeBuilder.Build().Encode())
	m.dep.Chat().SysChat(0, heroId, shared_proto.ChatType_ChatWorld, text, shared_proto.ChatMsgType_ChatMsgSys, true, false, true, false)
}

//gogen:iface
func (m *MilitaryModule) ProcessCaptainBorn(proto *military.C2SCaptainBornProto, hc iface.HeroController) {
	captainId := i32.Uint64(proto.CaptainId)
	captainData := m.datas.GetCaptainData(captainId)
	if captainData == nil {
		logrus.Debugf("武将生成，没有该武将")
		hc.Send(military.ERR_CAPTAIN_BORN_FAIL_NO_CAPTAIN)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if !heromodule.HasEnoughCost(hero, captainData.GetFirstStar().Cost) {
			logrus.Debugf("武将生成，将魂数量不足")
			result.Add(military.ERR_CAPTAIN_BORN_FAIL_ITEM_NOT_ENOUGH)
			return
		}
		if !heromodule.TryAddCaptain(hero, result, captainData, m.timeService.CurrentTime()) {
			logrus.Debugf("武将生成，已经获得该武将")
			result.Add(military.ERR_CAPTAIN_BORN_FAIL_EXISTED)
			return
		}
		hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCaptainBorn)
		heromodule.ReduceCostAnyway(hctx, hero, result, captainData.GetFirstStar().Cost)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_CAPTAIN_BORN) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_CAPTAIN_BORN)))
		}

		result.Changed()
		result.Ok()

		return
	})

}

//gogen:iface
func (m *MilitaryModule) ProcessCaptainUpstar(proto *military.C2SCaptainUpstarProto, hc iface.HeroController) {
	captainId := i32.Uint64(proto.CaptainId)
	captainData := m.datas.GetCaptainData(captainId)
	if captainData == nil {
		logrus.Debugf("武将升星，没有该武将")
		hc.Send(military.ERR_CAPTAIN_UPSTAR_FAIL_NO_CAPTAIN)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroMilitary := hero.Military()
		captain := heroMilitary.Captain(captainId)
		if captain == nil {
			logrus.Debugf("武将升星，未获得该武将")
			result.Add(military.ERR_CAPTAIN_UPSTAR_FAIL_NOT_GAIN)
			return
		}

		nextStar := captain.StarData().GetNextStar()
		if nextStar == nil {
			logrus.Debugf("武将升星，已经升至最高星")
			result.Add(military.ERR_CAPTAIN_UPSTAR_FAIL_MAX_STAR)
			return
		}

		hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCaptainUpstar)
		if !heromodule.TryReduceCost(hctx, hero, result, nextStar.Cost) {
			logrus.Debugf("武将升星，将魂数量不足")
			result.Add(military.ERR_CAPTAIN_UPSTAR_FAIL_ITEM_NOT_ENOUGH)
			return
		}

		captain.Upstar()
		result.Add(military.NewS2cCaptainUpstarMsg(u64.Int32(captainId), u64.Int32(captain.Star())))

		captain.CalculateProperties()
		result.Add(captain.NewUpdateCaptainStatMsg())

		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_UPSTAR)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_CAPTAIN_UP_STAR) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_CAPTAIN_UP_STAR)))
		}

		// 武将内政技能
		if nextStar.HasBuildingEffectSpell() {
			ctime := m.timeService.CurrentTime()
			heromodule.UpdateBuildingEffect(hero, result, m.datas, ctime,
				nextStar.GetBuildingEffectSpell(captain.AbilityData().UnlockSpellCount)...)
		}

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *MilitaryModule) ProcessCaptainExchange(proto *military.C2SCaptainExchangeProto, hc iface.HeroController) {
	cid1 := u64.FromInt32(proto.Cap1Id)
	cid2 := u64.FromInt32(proto.Cap2Id)

	if cid1 == cid2 {
		logrus.Debugf("武将经验传承，两个武将相同")
		hc.Send(military.ERR_CAPTAIN_EXCHANGE_FAIL_INVALID_CAPTAIN)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.MilitaryCaptainLevelExchange)
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		c1 := hero.Military().Captain(cid1)
		if c1 == nil {
			logrus.Debugf("武将经验传承，1没有该武将")
			result.Add(military.ERR_CAPTAIN_EXCHANGE_FAIL_INVALID_CAPTAIN)
			return
		}
		if c1.IsOutSide() {
			result.Add(military.ERR_CAPTAIN_EXCHANGE_FAIL_CAPTAIN_OUTSIDE)
			return
		}
		c2 := hero.Military().Captain(cid2)
		if c2 == nil {
			logrus.Debugf("武将经验传承，2没有该武将")
			result.Add(military.ERR_CAPTAIN_EXCHANGE_FAIL_INVALID_CAPTAIN)
			return
		}
		if c2.IsOutSide() {
			result.Add(military.ERR_CAPTAIN_EXCHANGE_FAIL_CAPTAIN_OUTSIDE)
			return
		}
		if !heromodule.TryReduceCost(hctx, hero, result, m.datas.MiscGenConfig().CaptainResetCost) {
			result.Add(military.ERR_CAPTAIN_EXCHANGE_FAIL_NOT_ENOUGH_COST)
			return
		}

		// 新版传承：等级、成长、官职（功勋）、装备互换，2个武将宝石全部卸下
		c1.ExchangeLevel(c2)
		c1.ExchangeAbility(c2)
		c1.ExchangeOfficial(c2)
		if c1.OfficialId() > 0 {
			hero.Military().SetOfficialAnyway(c1.OfficialId(), c1.OfficialIndex(), c1.Id())
		}
		if c2.OfficialId() > 0 {
			hero.Military().SetOfficialAnyway(c2.OfficialId(), c2.OfficialIndex(), c2.Id())
		}
		c1.ExchangeEquipments(c2)
		gems := c1.RemoveAllGem()
		if len(gems) > 0 {
			heromodule.AddGemArrayGive1(hctx, hero, result, gems, false)
			result.Add(gem.NewS2cOneKeyUseGemMsg(u64.Int32(c1.Id()), true, []int32{}, 0))
		}
		gems = c2.RemoveAllGem()
		if len(gems) > 0 {
			heromodule.AddGemArrayGive1(hctx, hero, result, gems, false)
			result.Add(gem.NewS2cOneKeyUseGemMsg(u64.Int32(c2.Id()), true, []int32{}, 0))
		}
		c1.CalculateProperties()
		heromodule.UpdateTroopFightAmount(hero, c1.GetTroop(), result)
		c2.CalculateProperties()
		heromodule.UpdateTroopFightAmount(hero, c2.GetTroop(), result)

		result.Add(military.NewS2cCaptainExchangeMarshalMsg(c1.EncodeClient(), c2.EncodeClient()))
		result.Changed()
		result.Ok()
	})

}

//gogen:iface c2s_notice_captain_has_viewed
func (m *MilitaryModule) ProcessNoticeCaptainHasViewed(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captains := hero.Military().Captains()
		changed := false
		for _, captain := range captains {
			if captain.TrySetViewedTrue() {
				changed = true
			}
		}
		if !changed {
			logrus.Debugf("去除武将新标识，没有任何武将需要设置")
			result.Add(military.ERR_NOTICE_CAPTAIN_HAS_VIEWED_FAIL_NO_CAPTAIN_NO_VIEWED)
			return
		}
		result.Add(military.NOTICE_CAPTAIN_HAS_VIEWED_S2C)
		result.Ok()
	})
}

// 激活武将羁绊
//gogen:iface
func (m *MilitaryModule) ProcessActivateCaptainFetter(proto *military.C2SActivateCaptainFriendshipProto, hc iface.HeroController) {
	data := m.datas.GetCaptainFriendshipData(u64.FromInt32(proto.GetId()))
	if data == nil {
		logrus.Debugf("激活武将羁绊，无效的羁绊id")
		hc.Send(military.ERR_ACTIVATE_CAPTAIN_FRIENDSHIP_FAIL_INVALID_ID)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroMilitary := hero.Military()
		for _, c := range data.Captains {
			if nil == heroMilitary.Captain(c.Id) {
				logrus.Debugf("激活武将羁绊，条件不足")
				result.Add(military.ERR_ACTIVATE_CAPTAIN_FRIENDSHIP_FAIL_NO_FETTER)
				return
			}
		}
		if !heroMilitary.TrySetCaptainFriendship(data) {
			logrus.Debugf("激活武将羁绊，已激活")
			result.Add(military.ERR_ACTIVATE_CAPTAIN_FRIENDSHIP_FAIL_ACTIVATED)
			return
		}

		result.Add(military.NewS2cActivateCaptainFriendshipMsg(proto.Id))

		// 更新武将属性
		if data.AllStat != nil || len(data.Race) > 0 {
			for _, c := range hero.Military().Captains() {
				if data.IsValidRace(c.Race().Race) {
					c.CalculateProperties()
					result.Add(c.NewUpdateCaptainStatMsg())
				}
			}

			heromodule.UpdateAllTroopFightAmount(hero, result)
		}

		result.Ok()
	})
}

// 告诉服务器已经预览过的官职槽位（去除红点标识）
//gogen:iface
func (m *MilitaryModule) ProcessNoticeOfficialHasViewed(proto *military.C2SNoticeOfficialHasViewedProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		exist, ok := hero.Military().TrySetOfficialViewed(u64.FromInt32(proto.GetOfficialId()), proto.GetOfficialIdx())
		if !exist {
			logrus.Debugf("去除官职红点标识，官职不存在")
			result.Add(military.ERR_NOTICE_OFFICIAL_HAS_VIEWED_FAIL_NO_OFFICIAL)
			return
		}
		if !ok {
			logrus.Debugf("去除官职红点标识，已经设置过")
			result.Add(military.ERR_NOTICE_OFFICIAL_HAS_VIEWED_FAIL_VIEWED)
			return
		}
		result.Add(military.NewS2cNoticeOfficialHasViewedMsg(proto.OfficialId, proto.OfficialIdx))

		result.Ok()
	})
}
