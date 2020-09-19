package equipment

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/depot"
	"github.com/lightpaw/male7/gen/pb/equipment"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/golang-lru"
	"github.com/lightpaw/pbutil"
)

const EquipProtoCacheCapacity = 2000

func NewEquipmentModule(dep iface.ServiceDep) *EquipmentModule {
	m := &EquipmentModule{}
	m.dep = dep
	m.datas = dep.Datas()
	m.timeService = dep.Time()

	cache, err := lru.New(EquipProtoCacheCapacity)
	if err != nil {
		logrus.WithError(err).Panic("NewEquipmentModule:new cache error")
	}
	m.equipMsgCache = cache

	heromodule.RegisterAddGoodsListener(m)

	return m
}

//gogen:iface
type EquipmentModule struct {
	dep         iface.ServiceDep
	datas       iface.ConfigDatas
	timeService iface.TimeService

	// 装备（各种升级、强化）的属性缓存
	equipMsgCache *lru.Cache
}

type EquipMsg struct {
	msg pbutil.Buffer
}

func (m EquipMsg) Version() uint64 {
	return 0
}

func genEquipChacheId(dataId, level, refined uint64) int64 {
	return int64(refined + level * 100 + dataId * 100000)
}

// 点击查看装备
//gogen:iface
func (m *EquipmentModule) ProcessViewChatEquip(proto *equipment.C2SViewChatEquipProto, hc iface.HeroController) {
	dataId := u64.FromInt32(proto.GetDataId())
	level := u64.FromInt32(proto.GetLevel())
	refined := u64.FromInt32(proto.GetRefined())

	id := genEquipChacheId(dataId, level, refined)
	if result, has := m.equipMsgCache.Get(id); has {
		hc.Send(result.(*EquipMsg).msg)
	} else {
		equipData := m.datas.GetEquipmentData(dataId)
		if equipData == nil {
			logrus.Debugf("点击查看装备，无效的装备id")
			hc.Send(equipment.ERR_VIEW_CHAT_EQUIP_FAIL_INVALID)
			return
		}
		refinedData := m.datas.GetEquipmentRefinedData(refined)
		if refinedData == nil {
			logrus.Debugf("点击查看装备，无效的装备强化等级")
			hc.Send(equipment.ERR_VIEW_CHAT_EQUIP_FAIL_INVALID)
			return
		}
		equip := entity.NewEquipment(0, equipData)
		equip.SetLevelData(equipData.Quality.MustLevel(level))
		equip.SetRefinedData(refinedData)
		equip.CalculateProperties()
		msg := equipment.NewS2cViewChatEquipMarshalMsg(equip.EncodeClient()).Static()
		m.equipMsgCache.Add(id, &EquipMsg { msg: msg })
		hc.Send(msg)
	}
}

func (m *EquipmentModule) OnAddGoodsEvent(hero *entity.Hero, result herolock.LockResult, goodsId uint64, addCount uint64) {
	combineDatas := m.datas.EquipCombineDatas()
	openCombineEquip := hero.OpenCombineEquip()

	if openCombineEquip.IsAllOpen(combineDatas) {
		// 都解锁了
		return
	}

	canOpenCombineData := combineDatas.GetOpenData(goodsId)
	if canOpenCombineData == nil {
		// 不影响开启数据
		return
	}

	if !openCombineEquip.Open(canOpenCombineData) {
		// 已经开启过了
		return
	}

	// 开启成功，发送消息
	result.Add(canOpenCombineData.OpenMsg)
	result.Changed()
}

//gogen:iface
func (m *EquipmentModule) ProcessWearEquipment(proto *equipment.C2SWearEquipmentProto, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		captainId := u64.FromInt32(proto.CaptainId)

		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("装备穿戴，武将id不存在")
			result.Add(equipment.ERR_WEAR_EQUIPMENT_FAIL_INVALID_CAPTAIN_ID)
			return
		}

		equipmentId := u64.FromInt32(proto.EquipmentId)
		var upId, downId uint64

		ctime := m.timeService.CurrentTime()
		heroDepot := hero.Depot()

		if proto.Down {
			// 出征武将不能卸下装备
			if captain.IsOutSide() {
				logrus.Debugf("装备穿戴，出征武将不能卸下装备")
				result.Add(equipment.ERR_WEAR_EQUIPMENT_FAIL_CAPTAIN_OUTSIDE)
				return
			}

			if !heroDepot.HasEnoughGenIdGoodsCapacity(goods.EQUIPMENT, 1) {
				logrus.Debugf("装备穿戴，空间不够")
				result.Add(equipment.ERR_WEAR_EQUIPMENT_FAIL_DEPOT_EQUIP_FULL)
				return
			}

			// 脱装备
			e := captain.RemoveEquipment(equipmentId)
			if e == nil {
				logrus.Debugf("装备穿戴，脱的装备id不存在")
				result.Add(equipment.ERR_WEAR_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID)
				return
			}

			heroDepot.AddGenIdGoods(e, ctime)

			downId = e.Id()

			m.dep.Tlog().TlogPlayerEquipFlow(hero, captainId, captain.Level(), e.Data().Id, e.Data().Quality.GoodsQuality.Level, uint64(e.Data().Type), 0)
		} else {
			// 穿装备
			g, haveExpireTime := heroDepot.GetNotExpiredGenIdGoods(equipmentId, ctime)
			if g == nil {
				logrus.Debugf("装备穿戴，穿的装备id不存在")
				result.Add(equipment.ERR_WEAR_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID)
				return
			}

			e, ok := g.(*entity.Equipment)
			if !ok {
				logrus.Debugf("装备穿戴，穿的不是个装备")
				result.Add(equipment.ERR_WEAR_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID)
				return
			}

			if haveExpireTime && !heroDepot.HasEnoughGenIdGoodsCapacity(goods.EQUIPMENT, 1) {
				logrus.Debugf("装备穿戴，没有足够的空间")
				result.Add(equipment.ERR_WEAR_EQUIPMENT_FAIL_DEPOT_EQUIP_FULL)
				return
			}

			// 出征武将
			if captain.IsOutSide() {
				if oldEqu := captain.GetEquipment(e.Data().Type); oldEqu != nil {
					if !proto.Inhert {
						logrus.Debugf("装备穿戴，出征武将替换装备必须继承")
						result.Add(equipment.ERR_WEAR_EQUIPMENT_FAIL_CAPTAIN_OUTSIDE_MUST_INHERIT)
						return
					}

					if e.Data().Quality.GoodsQuality.Level <= oldEqu.Data().Quality.GoodsQuality.Level {
						logrus.Debugf("装备穿戴，出征武将不能替换装备")
						result.Add(equipment.ERR_WEAR_EQUIPMENT_FAIL_CAPTAIN_OUTSIDE_QUALITY_ERR)
						return
					}
				}
			}

			heroDepot.RemoveGenIdGoods(equipmentId)
			old := captain.SetEquipment(e)

			if old != nil {
				heroDepot.AddGenIdGoods(old, ctime)
				downId = old.Id()
				// 继承
				if proto.Inhert {
					m.inherit(old, e, hero, result, captain)
				}

				m.dep.Tlog().TlogPlayerEquipFlow(hero, captainId, captain.Level(), old.Data().Id, old.Data().Quality.GoodsQuality.Level, uint64(old.Data().Type), 0)
			} else {
				moveToDepotIds := heroDepot.MoveTmpGoodsToDepotIfDepotHaveSlot(goods.EQUIPMENT, ctime)
				if len(moveToDepotIds) > 0 {
					result.Add(depot.NewS2cGoodsExpireTimeRemoveMsg(u64.Int32Array(moveToDepotIds)))
				}
			}
			upId = equipmentId

			m.dep.Tlog().TlogPlayerEquipFlow(hero, captainId, captain.Level(), e.Data().Id, e.Data().Quality.GoodsQuality.Level, uint64(e.Data().Type), 1)

			// 更新任务进度
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_EQUIPMENT)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_LEVEL_Y)
			heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_QUALITY_Y)

			if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_WEAR_EQUIPMENT) {
				result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_WEAR_EQUIPMENT)))
			}
		}
		// 套装
		oldTaozLv := captain.TaozLevel()
		captain.UpdateMorale(m.datas.EquipmentTaozConfig())

		if taozLv := captain.TaozLevel(); taozLv != oldTaozLv {
			// 套装系统广播
			hctx := heromodule.NewContext(m.dep, operate_type.EquipmentTaozUpgradeQuality)
			if d := hctx.BroadcastHelp().TaozLevel; d != nil {
				hctx.AddBroadcast(d, hero, result, 0, taozLv, func() *i18n.Fields {
					text := d.NewTextFields()
					text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
					text.WithFields(data.KeyNum, captain.TaozStar())
					return text
				})
			}
		}

		// 更新武将属性，战斗力
		captain.CalculateProperties()
		result.Add(captain.NewUpdateCaptainStatMsg())
		heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)

		result.Changed()

		result.Add(equipment.NewS2cWearEquipmentMsg(u64.Int32(captainId), u64.Int32(upId), u64.Int32(downId), u64.Int32(captain.TaozLevel())))

		result.Ok()
	})
}

//gogen:iface
func (m *EquipmentModule) ProcessUpgradeEquipment(proto *equipment.C2SUpgradeEquipmentProto, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		captainId := u64.FromInt32(proto.CaptainId)

		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("装备升级，武将id不存在")
			result.Add(equipment.ERR_UPGRADE_EQUIPMENT_FAIL_INVALID_CAPTAIN_ID)
			return
		}

		equipmentId := u64.FromInt32(proto.EquipmentId)
		// 穿装备
		e := captain.GetEquipmentById(equipmentId)
		if e == nil {
			logrus.Debugf("装备升级，装备id不存在")
			result.Add(equipment.ERR_UPGRADE_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID)
			return
		}

		if e.Data().Quality.GoodsQuality.Quality < shared_proto.Quality_PURPLE {
			logrus.Debugf("装备升级，紫色装备以下的装备不能升级")
			result.Add(equipment.ERR_UPGRADE_EQUIPMENT_FAIL_LEVEL_LIMIT)
			return
		}

		if e.Level() >= hero.LevelData().Sub.EquipmentLevelLimit {
			logrus.Debugf("装备升级，装备等级已达上限")
			result.Add(equipment.ERR_UPGRADE_EQUIPMENT_FAIL_LEVEL_LIMIT)
			return
		}

		nextLevel := e.LevelData().NextLevel()
		if nextLevel == nil {
			logrus.Debugf("装备升级，装备等级已达上限")
			result.Add(equipment.ERR_UPGRADE_EQUIPMENT_FAIL_LEVEL_LIMIT)
			return
		}

		newIntLevel := u64.Min(e.Level()+u64.FromInt32(proto.UpgradeTimes), hero.LevelData().Sub.EquipmentLevelLimit)

		newLevel := nextLevel
		if e.Level() < newIntLevel {
			newLevel = e.Data().Quality.MustLevel(newIntLevel)
		}

		// 检查消耗
		costCount := u64.Sub(newLevel.UpgradeLevelCost, e.LevelData().UpgradeLevelCost)
		costGoodsId := m.datas.GoodsConfig().EquipmentUpgradeGoods.Id
		heroDepot := hero.Depot()
		if !heroDepot.HasEnoughGoods(costGoodsId, costCount) {
			// 不能一次性升满，接下来一个个检查，看下能升几次

			current := e.LevelData()
			for i := e.Level(); i < newLevel.Level; i++ {
				nextLevel := current.NextLevel()
				if nextLevel == nil {
					logrus.Errorf("装备升级多次，current.NextLevel() == nil")
					result.Add(equipment.ERR_UPGRADE_EQUIPMENT_FAIL_SERVER_ERROR)
					return
				}

				costCount = u64.Sub(nextLevel.UpgradeLevelCost, e.LevelData().UpgradeLevelCost)
				if !heroDepot.HasEnoughGoods(costGoodsId, costCount) {
					newLevel = current
					break
				}

				current = current.NextLevel()
			}

			costCount = u64.Sub(newLevel.UpgradeLevelCost, e.LevelData().UpgradeLevelCost)
			if !heroDepot.HasEnoughGoods(costGoodsId, costCount) {
				logrus.Debugf("装备升级多次，消耗不足")
				result.Add(equipment.ERR_UPGRADE_EQUIPMENT_FAIL_COST_NOT_ENOUGH)
				return
			}
		}

		// 出征武将不能操作
		//if captain.IsOutSide() {
		//	logrus.Debugf("装备升级，出征武将不能操作")
		//	result.Add(equipment.ERR_UPGRADE_EQUIPMENT_FAIL_CAPTAIN_OUTSIDE)
		//	return
		//}

		hctx := heromodule.NewContext(m.dep, operate_type.EquipmentUpgrade)
		oldLevel := e.Level()

		// 扣除消耗
		newCount := heroDepot.RemoveGoods(costGoodsId, costCount)
		result.Changed()

		e.SetLevelData(newLevel)
		e.AddUpgradeCostCount(costCount)
		e.CalculateProperties()

		// 更新武将属性，战斗力
		captain.CalculateProperties()
		result.Add(captain.NewUpdateCaptainStatMsg())
		heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)

		result.Add(depot.NewS2cUpdateGoodsMsg(u64.Int32(costGoodsId), u64.Int32(newCount)))
		result.Add(equipment.NewS2cUpdateEquipmentMarshalMsg(e.EncodeClient()))

		result.Add(equipment.NewS2cUpgradeEquipmentMsg(u64.Int32(captainId), u64.Int32(equipmentId), u64.Int32(newLevel.Level)))

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_EQUIPMENT)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_LEVEL_Y)

		result.Changed()
		result.Ok()

		// tlog
		hctx.Tlog().TlogStrenghEquipmentFlow(hero, captainId, captain.Level(), e.Data().Id, operate_type.EquipUpgrade, operate_type.EquipNoInherit, oldLevel, e.Level(), e.RefinedLevel(), e.RefinedLevel(), uint64(e.Data().Type))

	})
}

// 继承
func (m *EquipmentModule) inherit(oldEqu *entity.Equipment, newEqu *entity.Equipment, hero *entity.Hero, result herolock.LockResult, captain *entity.Captain) {
	if oldEqu == nil || newEqu == nil {
		return
	}
	var isUpdate bool
	if oldEqu.Level() > newEqu.Level() {
		isUpdate = m.inheritUpgrade(oldEqu, newEqu, hero, result, captain) || isUpdate
	}

	if oldEqu.RefinedLevel() > newEqu.RefinedLevel() {
		isUpdate = m.inheritRefine(oldEqu, newEqu, hero, result, captain) || isUpdate
	}

	if isUpdate {
		result.Add(equipment.NewS2cUpdateEquipmentMarshalMsg(newEqu.EncodeClient()))
	}
}

func (m *EquipmentModule) inheritUpgrade(oldEqu *entity.Equipment, newEqu *entity.Equipment, hero *entity.Hero, result herolock.LockResult, captain *entity.Captain) (isUpdate bool) {
	oldEquOldLevel := oldEqu.Level()

	goodsCount := oldEqu.RebuildUpgrade()
	if goodsCount <= 0 {
		return
	}
	result.Add(equipment.NewS2cUpdateEquipmentMarshalMsg(oldEqu.EncodeClient()))
	result.Changed()

	newLevel := newEqu.LevelData()
	currentCost := newLevel.UpgradeLevelCost
	maxForSize := u64.Int(u64.Sub(hero.LevelData().Sub.EquipmentLevelLimit, newLevel.Level))
	for i := 0; i <= maxForSize; i++ {
		nextLevel := newLevel.NextLevel()
		if nextLevel == nil {
			break
		}
		if goodsCount < u64.Sub(nextLevel.UpgradeLevelCost, currentCost) {
			break
		}
		if newLevel.Level >= hero.LevelData().Sub.EquipmentLevelLimit {
			break
		}
		newLevel = nextLevel
	}

	costCount := u64.Sub(newLevel.UpgradeLevelCost, currentCost)
	backCount := u64.Sub(goodsCount, costCount)

	hctx := heromodule.NewContext(m.dep, operate_type.EquipmentInheritUpgrade)
	if costCount > 0 {
		newEquOldLevel := newEqu.Level()

		newEqu.SetLevelData(newLevel)
		newEqu.AddUpgradeCostCount(costCount)
		newEqu.CalculateProperties()

		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_EQUIPMENT)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_REFINE_LEVEL_Y)
		result.Changed()

		isUpdate = true

		// tlog
		if captain != nil {
			hctx.Tlog().TlogStrenghEquipmentFlow(hero, captain.Id(), captain.Level(), oldEqu.Data().Id, operate_type.EquipUpgrade, operate_type.EquipInherit, oldEquOldLevel, oldEqu.Level(), oldEqu.RefinedLevel(), oldEqu.RefinedLevel(), uint64(oldEqu.Data().Type))
			hctx.Tlog().TlogStrenghEquipmentFlow(hero, captain.Id(), captain.Level(), newEqu.Data().Id, operate_type.EquipUpgrade, operate_type.EquipInherit, newEquOldLevel, newEqu.Level(), newEqu.RefinedLevel(), newEqu.RefinedLevel(), uint64(newEqu.Data().Type))
		} else {
			hctx.Tlog().TlogStrenghEquipmentFlow(hero, 0, 0, oldEqu.Data().Id, operate_type.EquipUpgrade, operate_type.EquipInherit, oldEquOldLevel, oldEqu.Level(), oldEqu.RefinedLevel(), oldEqu.RefinedLevel(), uint64(oldEqu.Data().Type))
			hctx.Tlog().TlogStrenghEquipmentFlow(hero, 0, 0, newEqu.Data().Id, operate_type.EquipUpgrade, operate_type.EquipInherit, newEquOldLevel, newEqu.Level(), newEqu.RefinedLevel(), newEqu.RefinedLevel(), uint64(newEqu.Data().Type))
		}

	}
	if backCount > 0 {
		heromodule.AddGoods(hctx, hero, result, m.datas.GoodsConfig().EquipmentUpgradeGoods, backCount)
	}
	return
}

func (m *EquipmentModule) inheritRefine(oldEqu *entity.Equipment, newEqu *entity.Equipment, hero *entity.Hero, result herolock.LockResult, captain *entity.Captain) (isUpdate bool) {
	oldEquOldRefined := oldEqu.RefinedLevel()

	oldRefinedData := oldEqu.RebuildRefine()
	if oldRefinedData == nil {
		return
	}
	result.Add(equipment.NewS2cUpdateEquipmentMarshalMsg(oldEqu.EncodeClient()))
	result.Changed()

	goodsCount := oldRefinedData.TotalCostCount
	curRefinedData := newEqu.RefinedData()

	hctx := heromodule.NewContext(m.dep, operate_type.EquipmentInheritRefine)

	var newLevel *goods.EquipmentRefinedData
	var currentCost uint64
	if curRefinedData == nil {
		newLevel = m.datas.EquipmentRefinedData().MinKeyData
		currentCost = 0
	} else {
		newLevel = curRefinedData.NextLevel()
		currentCost = curRefinedData.TotalCostCount
	}
	if newLevel == nil || newLevel.Level > newEqu.Data().Quality.RefinedLevelLimit {
		heromodule.AddGoods(hctx, hero, result, m.datas.GoodsConfig().EquipmentRefinedGoods, goodsCount)
		return
	}

	maxForSize := u64.Int(u64.Sub(newEqu.Data().Quality.RefinedLevelLimit, newLevel.Level))
	for i := 0; i <= maxForSize; i++ {
		nextLevel := newLevel.NextLevel()
		if nextLevel == nil {
			break
		}
		if goodsCount < u64.Sub(nextLevel.TotalCostCount, currentCost) {
			break
		}
		if newLevel.Level >= newEqu.Data().Quality.RefinedLevelLimit {
			break
		}
		newLevel = nextLevel
	}

	costCount := u64.Sub(newLevel.TotalCostCount, currentCost)
	backCount := u64.Sub(goodsCount, costCount)
	if costCount > 0 {
		newEquOldRefined := newEqu.RefinedLevel()

		newEqu.SetRefinedData(newLevel)
		newEqu.CalculateProperties()

		result.Add(equipment.NewS2cRefinedEquipmentMsg(u64.Int32(captain.Id()), u64.Int32(newEqu.Id()), u64.Int32(newLevel.Level), u64.Int32(captain.TaozLevel())))
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_REFINE_LEVEL_Y)
		result.Changed()

		isUpdate = true

		// tlog
		hctx.Tlog().TlogStrenghEquipmentFlow(hero, captain.Id(), captain.Level(), oldEqu.Data().Id, operate_type.EquipRefine, operate_type.EquipInherit, oldEqu.Level(), oldEqu.Level(), oldEquOldRefined, oldEqu.RefinedLevel(), uint64(oldEqu.Data().Type))
		hctx.Tlog().TlogStrenghEquipmentFlow(hero, captain.Id(), captain.Level(), newEqu.Data().Id, operate_type.EquipRefine, operate_type.EquipInherit, newEqu.Level(), newEqu.Level(), newEquOldRefined, newEqu.RefinedLevel(), uint64(newEqu.Data().Type))
	}
	if backCount > 0 {
		heromodule.AddGoods(hctx, hero, result, m.datas.GoodsConfig().EquipmentRefinedGoods, backCount)
	}
	return
}

//gogen:iface
func (m *EquipmentModule) ProcessUpgradeEquipmentAll(proto *equipment.C2SUpgradeEquipmentAllProto, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		captainId := u64.FromInt32(proto.CaptainId)

		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("装备升级，武将id不存在")
			result.Add(equipment.ERR_UPGRADE_EQUIPMENT_ALL_FAIL_INVALID_CAPTAIN_ID)
			return
		}

		costGoodsId := m.datas.GoodsConfig().EquipmentUpgradeGoods.Id
		heroDepot := hero.Depot()

		// 穿装备
		// 出征武将不能操作
		//if captain.IsOutSide() {
		//	logrus.Debugf("装备升级，出征武将不能操作")
		//	result.Add(equipment.ERR_UPGRADE_EQUIPMENT_ALL_FAIL_CAPTAIN_OUTSIDE)
		//	return
		//}

		hctx := heromodule.NewContext(m.dep, operate_type.EquipmentUpgradeAll)

		var newCount uint64
		var hasEquipment bool
		for _, t := range goods.TypeArray {
			e := captain.GetEquipment(t)
			if e == nil {
				continue
			}

			if e.Data().Quality.GoodsQuality.Quality < shared_proto.Quality_PURPLE {
				continue
			}

			if e.Level() >= hero.LevelData().Sub.EquipmentLevelLimit {
				continue
			}

			newLevel := e.LevelData().NextLevel()
			if newLevel == nil {
				continue
			}

			// 检查消耗
			costCount := newLevel.CurrentUpgradeLevelCost
			if !heroDepot.HasEnoughGoods(costGoodsId, costCount) {
				continue
			}
			hasEquipment = true

			// 扣除消耗
			newCount = heroDepot.RemoveGoods(costGoodsId, costCount)
			result.Changed()

			oldLevel := e.Level()

			// 升级
			e.UpgradeLevel()
			e.AddUpgradeCostCount(costCount)
			e.CalculateProperties()

			hctx.Tlog().TlogStrenghEquipmentFlow(hero, captainId, captain.Level(), e.Data().Id, operate_type.EquipUpgrade, operate_type.EquipNoInherit, oldLevel, e.Level(), e.RefinedLevel(), e.RefinedLevel(), uint64(e.Data().Type))
		}

		if !hasEquipment {
			logrus.Debugf("装备升级，消耗不足")
			result.Add(equipment.ERR_UPGRADE_EQUIPMENT_ALL_FAIL_COST_NOT_ENOUGH)
			return
		}

		// 更新武将属性，战斗力
		captain.CalculateProperties()
		result.Add(captain.NewUpdateCaptainStatMsg())
		heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)

		// 扣物品消息
		result.Add(depot.NewS2cUpdateGoodsMsg(u64.Int32(costGoodsId), u64.Int32(newCount)))

		level := make([]int32, 0, len(goods.TypeArray))
		eqs := make([][]byte, 0, len(goods.TypeArray))
		for _, t := range goods.TypeArray {
			e := captain.GetEquipment(t)
			if e == nil {
				level = append(level, 0)
			} else {
				level = append(level, u64.Int32(e.Level()))
				eqs = append(eqs, must.Marshal(e.EncodeClient()))
			}
		}

		result.Add(equipment.NewS2cUpdateMultiEquipmentMsg(eqs))
		result.Add(equipment.NewS2cUpgradeEquipmentAllMsg(proto.CaptainId, level))

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_EQUIPMENT)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_LEVEL_Y)

		result.Changed()
		result.Ok()
	})

}

// 装备升星
//gogen:iface
func (m *EquipmentModule) ProcessRefinedEquipment(proto *equipment.C2SRefinedEquipmentProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captainId := u64.FromInt32(proto.CaptainId)

		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("装备强化，武将id不存在")
			result.Add(equipment.ERR_REFINED_EQUIPMENT_FAIL_INVALID_CAPTAIN_ID)
			return
		}

		equipmentId := u64.FromInt32(proto.EquipmentId)
		e := captain.GetEquipmentById(equipmentId)
		if e == nil {
			logrus.Debugf("装备强化，装备id不存在")
			result.Add(equipment.ERR_REFINED_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID)
			return
		}

		refinedData := e.RefinedData()
		nextLevel := m.datas.EquipmentRefinedData().MinKeyData
		if refinedData != nil {
			nextLevel = m.datas.GetEquipmentRefinedData(refinedData.Level + 1)
		}

		if nextLevel == nil {
			logrus.Debugf("装备强化，当前是强化最大等级")
			result.Add(equipment.ERR_REFINED_EQUIPMENT_FAIL_LEVEL_LIMIT)
			return
		}

		if nextLevel.Level > e.Data().Quality.RefinedLevelLimit {
			logrus.Debugf("装备强化，装备强化等级已达上限")
			result.Add(equipment.ERR_REFINED_EQUIPMENT_FAIL_LEVEL_LIMIT)
			return
		}

		if nextLevel.HeroLevelLimit > hero.Level() {
			logrus.Debugf("装备强化，君主等级不足")
			result.Add(equipment.ERR_REFINED_EQUIPMENT_FAIL_HERO_LEVEL_LIMIT)
			return
		}

		// 检查消耗
		costGoodsId := m.datas.GoodsConfig().EquipmentRefinedGoods.Id
		heroDepot := hero.Depot()
		if !heroDepot.HasEnoughGoods(costGoodsId, nextLevel.CostCount) {
			logrus.Debugf("装备强化，消耗不足")
			result.Add(equipment.ERR_REFINED_EQUIPMENT_FAIL_COST_NOT_ENOUGH)
			return
		}

		var newCount uint64

		// 出征武将不能操作
		//if captain.IsOutSide() {
		//	logrus.Debugf("装备强化，出征武将不能操作")
		//	result.Add(equipment.ERR_REFINED_EQUIPMENT_FAIL_CAPTAIN_OUTSIDE)
		//	return
		//}

		hctx := heromodule.NewContext(m.dep, operate_type.EquipmentRefined)
		oldRefined := e.RefinedLevel()

		// 扣除消耗
		newCount = heroDepot.RemoveGoods(costGoodsId, nextLevel.CostCount)

		e.SetRefinedData(nextLevel)
		e.CalculateProperties()

		// 套装
		captain.UpdateMorale(m.datas.EquipmentTaozConfig())

		// 套装系统广播
		hc := heromodule.NewContext(m.dep, operate_type.EquipmentTaozUpgradeQuality)
		if d := hc.BroadcastHelp().TaozLevel; d != nil {
			hc.AddBroadcast(d, hero, result, 0, captain.TaozLevel(), func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hc.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, captain.TaozStar())
				return text
			})
		}

		// 更新武将属性，战斗力
		captain.CalculateProperties()
		result.Add(captain.NewUpdateCaptainStatMsg())
		heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)

		result.Add(depot.NewS2cUpdateGoodsMsg(u64.Int32(costGoodsId), u64.Int32(newCount)))

		result.Add(equipment.NewS2cUpdateEquipmentMarshalMsg(e.EncodeClient()))

		result.Add(equipment.NewS2cRefinedEquipmentMsg(u64.Int32(captainId), u64.Int32(equipmentId), u64.Int32(nextLevel.Level), u64.Int32(captain.TaozLevel())))

		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_X_EQIUP_REFINE_LEVEL_Y)

		if d := hctx.BroadcastHelp().EquipLevel; d != nil {
			hctx.AddBroadcast(d, hero, result, 0, e.RefinedLevel(), func() *i18n.Fields {
				text := d.NewTextFields()
				text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
				text.WithFields(data.KeyNum, e.RefinedStar())
				return text
			})
		}

		result.Changed()
		result.Ok()

		hctx.Tlog().TlogStrenghEquipmentFlow(hero, captainId, captain.Level(), e.Data().Id, operate_type.EquipRefine, operate_type.EquipNoInherit, e.Level(), e.Level(), oldRefined, e.RefinedLevel(), uint64(e.Data().Type))
	})
}

//gogen:iface
func (m *EquipmentModule) ProcessSmeltEquipment(proto *equipment.C2SSmeltEquipmentProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.EquipmentSmelt)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		ctime := m.timeService.CurrentTime()
		ids := u64.FromInt32Array(proto.EquipmentId)

		heroDepot := hero.Depot()

		smeltEquipcount := uint64(0)
		var totalBackUpgradeCount, totalBackRefinedCostCount uint64
		var equDatas []*goods.EquipmentData
		for _, id := range ids {
			g, _ := heroDepot.GetNotExpiredGenIdGoods(id, ctime)
			if g == nil {
				continue
			}

			e, ok := g.(*entity.Equipment)
			if !ok {
				continue
			}

			equDatas = append(equDatas, e.Data())

			// 删掉
			heroDepot.RemoveGenIdGoods(id)
			result.Changed()

			// 计算返还的物品
			upgradeCount, refinedData := e.Rebuild()
			totalBackUpgradeCount += upgradeCount + e.Data().Quality.SmeltBackCount

			smeltEquipcount++

			if refinedData != nil {
				totalBackRefinedCostCount += refinedData.TotalCostCount
			}
		}

		if smeltEquipcount <= 0 {
			result.Add(equipment.ERR_SMELT_EQUIPMENT_FAIL_INVALID_EQUIPMENT_ID)
			return
		}

		moveToDepotIds := heroDepot.MoveTmpGoodsToDepotIfDepotHaveSlot(goods.EQUIPMENT, ctime)
		if len(moveToDepotIds) > 0 {
			result.Add(depot.NewS2cGoodsExpireTimeRemoveMsg(u64.Int32Array(moveToDepotIds)))
		}

		if totalBackUpgradeCount > 0 {
			heromodule.AddGoods(hctx, hero, result, m.datas.GoodsConfig().EquipmentUpgradeGoods, totalBackUpgradeCount)
		}

		if totalBackRefinedCostCount > 0 {
			heromodule.AddGoods(hctx, hero, result, m.datas.GoodsConfig().EquipmentRefinedGoods, totalBackRefinedCostCount)
		}

		result.Add(equipment.NewS2cSmeltEquipmentMsg(proto.EquipmentId))

		heromodule.IncreTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_SMELT_EQUIP, smeltEquipcount)

		hero.HistoryAmount().Increase(server_proto.HistoryAmountType_SmeltEquip, smeltEquipcount)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_SMELT_EQUIP)

		result.Ok()

		for _, equ := range equDatas {
			hctx.Tlog().TlogMountRefreshFlow(hero, operate_type.EquipOperTypeSmelt, operate_type.EquipNoInherit, equ.Id, equ.Quality.GoodsQuality.Level, uint64(equ.Type))
		}
	})
}

//gogen:iface
func (m *EquipmentModule) ProcessRebuildEquipment(proto *equipment.C2SRebuildEquipmentProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.EquipmentRebuild)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		ctime := m.timeService.CurrentTime()

		ids := u64.FromInt32Array(proto.EquipmentId)

		heroDepot := hero.Depot()

		var totalBackUpgradeCount, totalBackRefinedCostCount uint64
		eqs := make([][]byte, 0, len(ids))
		var equDatas []*goods.EquipmentData
		for _, id := range ids {
			g, _ := heroDepot.GetNotExpiredGenIdGoods(id, ctime)
			if g == nil {
				continue
			}

			e, ok := g.(*entity.Equipment)
			if !ok {
				continue
			}

			if e.Level() <= 1 && e.RefinedData() == nil {
				continue
			}

			equDatas = append(equDatas, e.Data())

			upgradeLevelCost := e.LevelData().UpgradeLevelCost

			// 重新变回1级数据
			backUpgradeCount, refinedData := e.Rebuild()
			result.Changed()

			if e.Level() > 1 && backUpgradeCount <= 0 {
				backUpgradeCount = upgradeLevelCost
			}

			// 计算返还的物品
			totalBackUpgradeCount += backUpgradeCount

			if refinedData != nil {
				totalBackRefinedCostCount += refinedData.TotalCostCount
			}

			eqs = append(eqs, must.Marshal(e.EncodeClient()))
		}

		if totalBackUpgradeCount > 0 {
			heromodule.AddGoods(hctx, hero, result, m.datas.GoodsConfig().EquipmentUpgradeGoods, totalBackUpgradeCount)
		}

		if totalBackRefinedCostCount > 0 {
			heromodule.AddGoods(hctx, hero, result, m.datas.GoodsConfig().EquipmentRefinedGoods, totalBackRefinedCostCount)
		}

		result.Add(equipment.NewS2cRebuildEquipmentMsg(proto.EquipmentId))

		if len(eqs) > 0 {
			result.Add(equipment.NewS2cUpdateMultiEquipmentMsg(eqs))
		}

		result.Ok()

		for _, equ := range equDatas {
			hctx.Tlog().TlogMountRefreshFlow(hero, operate_type.EquipOperTypeRebuild, operate_type.EquipNoInherit, equ.Id, equ.Quality.GoodsQuality.Level, uint64(equ.Type))
		}
	})
}

// 一键卸装备
//gogen:iface c2s_one_key_take_off
func (m *EquipmentModule) ProcessOneKeyTakeOff(proto *equipment.C2SOneKeyTakeOffProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captain := hero.Military().Captain(u64.FromInt32(proto.GetCaptainId()))
		if captain == nil {
			logrus.Debugf("一键卸装，没有该武将")
			result.Add(equipment.ERR_ONE_KEY_TAKE_OFF_FAIL_NO_CAPTAIN)
			return
		}
		if captain.IsOutSide() {
			logrus.Debugf("一键卸装，武将出征中")
			result.Add(equipment.ERR_ONE_KEY_TAKE_OFF_FAIL_OUTSIDE)
			return
		}
		equipmentCount := captain.GetEquipmentCount()
		if equipmentCount <= 0 {
			logrus.Debugf("一键卸装，武将身上没有任何装备")
			result.Add(equipment.ERR_ONE_KEY_TAKE_OFF_FAIL_NO_EQUIPMENT)
			return
		}
		if !hero.Depot().HasEnoughGenIdGoodsCapacity(goods.EQUIPMENT, equipmentCount) {
			logrus.Debugf("一键卸装，背包空间不够")
			result.Add(equipment.ERR_ONE_KEY_TAKE_OFF_FAIL_DEPOT_SPACE_NOT_ENOUGH)
			return
		}
		// 脱光装备
		removedEquipments := captain.RemoveAllEquipment()
		if len(removedEquipments) <= 0 {
			logrus.Debugf("一键卸装，武将身上没有任何装备")
			result.Add(equipment.ERR_ONE_KEY_TAKE_OFF_FAIL_NO_EQUIPMENT)
			return
		}
		ctime := m.timeService.CurrentTime()
		for _, equip := range removedEquipments {
			hero.Depot().AddGenIdGoods(equip, ctime)
			m.dep.Tlog().TlogPlayerEquipFlow(hero, captain.Id(), captain.Level(), equip.Data().Id, equip.Data().Quality.GoodsQuality.Level, uint64(equip.Data().Type), 0)
		}
		// 套装
		captain.UpdateMorale(m.datas.EquipmentTaozConfig())
		// 更新武将属性，战斗力
		captain.CalculateProperties()
		result.Add(captain.NewUpdateCaptainStatMsg())
		heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)

		result.Changed()

		result.Add(equipment.NewS2cOneKeyTakeOffMsg(proto.CaptainId))
		result.Ok()
	})
}
