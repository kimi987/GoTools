package gem

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/gem"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"sort"
	"github.com/pkg/errors"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/service/operate_type"
)

func NewGemModule(dep iface.ServiceDep) *GemModule {
	return &GemModule{
		dep:         dep,
		configDatas: dep.Datas(),
		timeService: dep.Time(),
	}
}

// 宝石模块
//gogen:iface
type GemModule struct {
	dep         iface.ServiceDep
	configDatas iface.ConfigDatas
	timeService iface.TimeService
}

// 镶嵌宝石
//gogen:iface
func (m *GemModule) ProcessUseGem(proto *gem.C2SUseGemProto, hc iface.HeroController) {
	//hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
	//	captainId := u64.FromInt32(proto.CaptainId)
	//
	//	captain := hero.Military().Captain(captainId)
	//	if captain == nil {
	//		logrus.Debugf("镶嵌宝石，武将id不存在")
	//		result.Add(gem.ERR_USE_GEM_FAIL_INVALID_CAPTAIN_ID)
	//		return
	//	}
	//
	//	slotIdx := u64.FromInt32(proto.GetSlotIdx())
	//	slotDatas := m.configDatas.GetGemSlotDataArray()
	//	if slotIdx < 0 || slotIdx >= uint64(len(slotDatas)) {
	//		logrus.Debugf("镶嵌宝石，宝石槽位非法")
	//		result.Add(gem.ERR_USE_GEM_FAIL_INVALID_SLOT_IDX)
	//		return
	//	}
	//
	//	slotData := slotDatas[slotIdx]
	//	if slotData == nil {
	//		logrus.Debugf("镶嵌宝石，宝石槽位数据竟然没找到")
	//		result.Add(gem.ERR_USE_GEM_FAIL_INVALID_SLOT_IDX)
	//		return
	//	}
	//
	//	if captain.Ability() < slotData.NeedAbility {
	//		logrus.Debugf("镶嵌宝石，成长值不够")
	//		result.Add(gem.ERR_USE_GEM_FAIL_ABILITY_NOT_ENOUGH)
	//		return
	//	}
	//
	//	heroDepot := hero.Depot()
	//
	//	var oldId, newId uint64
	//
	//	hctx := heromodule.NewContext(m.dep, operate_type.GemUseGem)
	//
	//	if proto.Down {
	//		// 出征武将不能卸下宝石
	//		if captain.IsOutSide() {
	//			logrus.Debugf("镶嵌宝石，出征武将不能卸下宝石")
	//			result.Add(gem.ERR_USE_GEM_FAIL_CAPTAIN_OUTSIDE)
	//			return
	//		}
	//
	//		// 脱宝石
	//		downGem := captain.GetGem(slotIdx)
	//		if downGem == nil {
	//			logrus.Debugf("镶嵌宝石，脱的宝石不存在")
	//			result.Add(gem.ERR_USE_GEM_FAIL_INVALID_GEM_ID)
	//			return
	//		}
	//
	//		captain.RemoveGem(slotIdx)
	//
	//		heromodule.AddGem(hctx, hero, result, downGem, 1, false)
	//
	//		oldId = downGem.Id
	//		newId = 0
	//	} else {
	//		gemId := u64.FromInt32(proto.GetGemId())
	//
	//		// 穿装备
	//		upGem := m.configDatas.GetGemData(gemId)
	//		if upGem == nil {
	//			logrus.Debugf("镶嵌宝石，镶嵌的宝石id不存在")
	//			result.Add(gem.ERR_USE_GEM_FAIL_INVALID_GEM_ID)
	//			return
	//		}
	//
	//		if !heroDepot.HasEnoughGoods(gemId, 1) {
	//			logrus.Debugf("镶嵌宝石，宝石数量不够")
	//			result.Add(gem.ERR_USE_GEM_FAIL_GEM_NOT_ENOUGH)
	//			return
	//		}
	//
	//		for _, data := range slotData.GetSameTypeSlots() {
	//			if data == slotData {
	//				// 同一个位置
	//				continue
	//			}
	//			captainGem := captain.GetGem(data.SlotIdx)
	//			if captainGem != nil && captainGem.GemType == upGem.GemType {
	//				logrus.Debugf("镶嵌宝石，同一个部件不可以镶嵌两个相同类型的宝石")
	//				result.Add(gem.ERR_USE_GEM_FAIL_HAVE_SAME_TYPE_GEM_IN_THIS_TYPE)
	//				return
	//			}
	//		}
	//
	//		old := captain.GetGem(slotIdx)
	//		// 替换
	//		if old != nil && captain.IsOutSide() {
	//			if upGem.Level <= old.Level {
	//				logrus.Debugf("镶嵌宝石，必须替换为更高品质的宝石")
	//				result.Add(gem.ERR_USE_GEM_FAIL_CAPTAIN_OUTSIDE_QUALITY_ERR)
	//				return
	//			}
	//		}
	//
	//		heromodule.ReduceGemAnyway(hctx, hero, result, upGem, 1)
	//		old = captain.SetGem(slotIdx, upGem)
	//
	//		if old != nil {
	//			heromodule.AddGem(hctx, hero, result, old, 1, false)
	//			oldId = old.Id
	//		}
	//
	//		newId = gemId
	//	}
	//
	//	// 更新武将属性，战斗力
	//	captain.CalculateProperties()
	//	result.Add(captain.NewUpdateCaptainStatMsg())
	//	heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)
	//
	//	result.Add(gem.NewS2cUseGemMsg(u64.Int32(captainId), u64.Int32(slotIdx), u64.Int32(newId)))
	//
	//	if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_USE_GEM) {
	//		result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_USE_GEM)))
	//	}
	//
	//	result.Changed()
	//	result.Ok()
	//	// TODO 更新任务进度 宝石
	//
	//	hctx.Tlog().TlogEquipmentAddStarFlow(hero, captainId, captain.Level(), oldId, newId, slotIdx)
	//})
}

// 新版镶嵌(摘除)宝石
//gogen:iface
func (m *GemModule) ProcessInlayGem(proto *gem.C2SInlayGemProto, hc iface.HeroController) {
	captainId := u64.FromInt32(proto.CaptainId)
	captainData := m.configDatas.GetCaptainData(captainId)
	if captainData == nil {
		logrus.Debugf("镶嵌(摘除)，没有该武将配置")
		hc.Send(gem.ERR_INLAY_GEM_FAIL_NO_CAPTAIN)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("镶嵌(摘除)，并未获得该武将")
			result.Add(gem.ERR_INLAY_GEM_FAIL_NO_CAPTAIN)
			return
		}
		slotIdx := int(proto.GetSlotIdx())
		if slotIdx < 0 || slotIdx >= captain.GemSlotCapacity() {
			logrus.Debugf("镶嵌(摘除)，宝石槽位非法")
			result.Add(gem.ERR_INLAY_GEM_FAIL_WRONG_SLOT_IDX)
			return
		}

		oldGem := captain.GetGem(slotIdx)
		hctx := heromodule.NewContext(m.dep, operate_type.GemUseGem)
		gemId := u64.FromInt32(proto.GemId)
		if gemId != 0 { // 镶嵌新的宝石
			gemData := m.configDatas.GetGemData(gemId)
			if gemData == nil {
				logrus.Debugf("镶嵌(摘除)，没有该宝石配置")
				hc.Send(gem.ERR_INLAY_GEM_FAIL_NO_GEM)
				return
			}
			if !hero.Depot().HasEnoughGoods(gemId, 1) {
				logrus.Debugf("镶嵌(摘除)，背包里没有该宝石")
				result.Add(gem.ERR_INLAY_GEM_FAIL_NO_GEM)
				return
			}
			if gemData.GemType != captainData.Race.GemTypes[slotIdx] {
				logrus.Debugf("镶嵌(摘除)，宝石类型不匹配")
				hc.Send(gem.ERR_INLAY_GEM_FAIL_WRONG_GEM)
				return
			}
			if oldGem != nil && oldGem.Id == gemData.Id {
				logrus.Debugf("镶嵌(摘除)，已经有相同宝石")
				hc.Send(gem.ERR_INLAY_GEM_FAIL_HAS_SAME_GEM)
				return
			}
			// 扣除背包里的宝石
			heromodule.ReduceGemAnyway(hctx, hero, result, gemData, 1)
			captain.SetGem(slotIdx, gemData)
		} else { // 摘除旧的宝石
			if oldGem == nil {
				logrus.Debugf("镶嵌(摘除)，摘空气吗")
				hc.Send(gem.ERR_INLAY_GEM_FAIL_NO_GEM)
				return
			}
			captain.SetGem(slotIdx, nil)
		}
		if oldGem != nil { // 把原来的宝石丢回仓库
			heromodule.AddGem(hctx, hero, result, oldGem, 1, false)
		}
		// 更新武将属性，战斗力
		captain.CalculateProperties()
		result.Add(captain.NewUpdateCaptainStatMsg())
		heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)

		result.Add(gem.NewS2cInlayGemMsg(proto.CaptainId, proto.SlotIdx, proto.GemId))

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_USE_GEM) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_USE_GEM)))
		}

		result.Changed()
		result.Ok()
		// TODO 更新任务进度 宝石

		oldGemId := uint64(0)
		if oldGem != nil {
			oldGemId = oldGem.Id
		}
		hctx.Tlog().TlogEquipmentAddStarFlow(hero, captainId, captain.Level(), oldGemId, gemId, u64.FromInt(slotIdx))
	})
}

// 合成宝石
//gogen:iface
func (m *GemModule) ProcessCombineGem(proto *gem.C2SCombineGemProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.GemCombineGem)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captainId := u64.FromInt32(proto.CaptainId)

		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("合成宝石，武将id不存在")
			result.Add(gem.ERR_COMBINE_GEM_FAIL_INVALID_CAPTAIN_ID)
			return
		}

		// 出征武将不能操作宝石
		//if captain.IsOutSide() {
		//	logrus.Debugf("合成宝石，出征武将不能合成宝石")
		//	result.Add(gem.ERR_COMBINE_GEM_FAIL_CAPTAIN_OUTSIDE)
		//	return
		//}

		slotIdx := int(proto.GetSlotIdx())

		// 目标宝石
		targetGem := captain.GetGem(slotIdx)
		if targetGem == nil {
			logrus.Debugf("合成宝石，没有镶嵌宝石在上面")
			result.Add(gem.ERR_COMBINE_GEM_FAIL_INVALID_SLOT_IDX)
			return
		}

		if targetGem.NextLevel == nil {
			logrus.Debugf("合成宝石，宝石没有下一级")
			result.Add(gem.ERR_COMBINE_GEM_FAIL_LEVEL_MAX)
			return
		}

		heroDepot := hero.Depot()

		if !heroDepot.HasEnoughGoods(targetGem.Id, targetGem.UpgradeNeedCount-1) {
			logrus.Debugf("合成宝石，没有足够的物品")
			result.Add(gem.ERR_COMBINE_GEM_FAIL_NOT_ENOUGH)
			return
		}

		heromodule.ReduceGemAnyway(hctx, hero, result, targetGem, targetGem.UpgradeNeedCount-1)
		captain.SetGem(slotIdx, targetGem.NextLevel)

		// 更新武将属性，战斗力
		captain.CalculateProperties()
		result.Add(captain.NewUpdateCaptainStatMsg())
		heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)

		result.Add(gem.NewS2cCombineGemMsg(proto.GetCaptainId(), proto.GetSlotIdx(), u64.Int32(targetGem.NextLevel.Id)))

		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_GEM)

		// 系统广播
		if d := hctx.BroadcastHelp().GemGet; d != nil {
			heromodule.AddGemBroadcast(hctx, hero, targetGem, result)
		}

		result.Changed()
		result.Ok()
	})
}

// 一键镶嵌宝石
//gogen:iface
func (m *GemModule) ProcessOneKeyUseGem(proto *gem.C2SOneKeyUseGemProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captainId := u64.FromInt32(proto.CaptainId)

		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("镶嵌宝石，武将id不存在")
			result.Add(gem.ERR_ONE_KEY_USE_GEM_FAIL_INVALID_CAPTAIN_ID)
			return
		}

		//// 出征武将不能操作宝石
		//if captain.IsOutSide() {
		//	logrus.Debugf("镶嵌宝石，出征武将不能镶嵌宝石")
		//	result.Add(gem.ERR_ONE_KEY_USE_GEM_FAIL_CAPTAIN_OUTSIDE)
		//	return
		//}

		hctx := heromodule.NewContext(m.dep, operate_type.GemOneKeyUseGem)

		gemsChanged := false

		var newId []int32
		if proto.GetDownAll() {
			// 出征武将不能卸下宝石
			if captain.IsOutSide() {
				logrus.Debugf("镶嵌宝石，出征武将不能卸下宝石")
				result.Add(gem.ERR_ONE_KEY_USE_GEM_FAIL_CAPTAIN_OUTSIDE)
				return
			}

			// 脱宝石
			if captain.GetGemCount() <= 0 {
				logrus.Debugf("镶嵌宝石，武将身上没宝石")
				result.Add(gem.ERR_ONE_KEY_USE_GEM_FAIL_NO_GEM_TO_DOWN_ALL)
				return
			}

			proto.EquipType = 0

			gems := captain.RemoveAllGem()
			heromodule.AddGemArrayGive1(hctx, hero, result, gems, false)

			gemsChanged = len(gems) > 0
		} else {

			newGemsArray, errMsg := m.calcOneKeyUseGemCost(hero, captain)
			if errMsg != nil {
				result.Add(errMsg)
				return
			}

			newId = make([]int32, len(newGemsArray))
			for slotIdx, newGem := range newGemsArray {
				if newGem == nil {
					continue
				}

				// 新的宝石id
				newId[slotIdx] = u64.Int32(newGem.Id)

				if captain.GetGem(slotIdx) == newGem {
					// 相同的宝石
					continue
				}

				gemsChanged = true

				heromodule.ReduceGemAnyway(hctx, hero, result, newGem, 1)

				// 设置了新的宝石
				oldGem := captain.SetGem(slotIdx, newGem)
				var oldId uint64
				if oldGem != nil {
					heromodule.AddGem(hctx, hero, result, oldGem, 1, false)
					oldId = oldGem.Id
				}

				// tlog
				hctx.Tlog().TlogEquipmentAddStarFlow(hero, captainId, captain.Level(), oldId, newGem.Id, u64.FromInt(slotIdx))

			}
		}

		// 更新武将属性，战斗力
		if gemsChanged {
			captain.CalculateProperties()
			result.Add(captain.NewUpdateCaptainStatMsg())
			heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)
		}

		result.Add(gem.NewS2cOneKeyUseGemMsg(u64.Int32(captainId), proto.GetDownAll(), newId, proto.EquipType))

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_USE_GEM) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_USE_GEM)))
		}

		result.Changed()
		result.Ok()
		// TODO 更新任务进度 宝石
	})
}

func (m *GemModule) calcOneKeyUseGemCost(hero *entity.Hero, captain *entity.Captain) (newGemsArray []*goods.GemData, errMsg pbutil.Buffer) {

	heroDepot := hero.Depot()

	// 获取武将的宝石槽位

	// 排序，空位排在最前

	// 每个宝石位置，根据宝石类型，宝石等级从大到小遍历

	mayUseOrUpgradeSlotArray := make([]*gemAndSlot, 0, captain.GemSlotCapacity())
	for idx, gem := range captain.Gems() {
		if gem != nil && gem.NextLevel == nil {
			// 最后一级满级了，不需要升级了
			continue
		}

		g := &gemAndSlot{
			gem:     gem,
			slotIdx: idx,
		}

		mayUseOrUpgradeSlotArray = append(mayUseOrUpgradeSlotArray, g)
	}

	if len(mayUseOrUpgradeSlotArray) <= 0 {
		logrus.Debugln("一键镶嵌宝石，一个可以用的槽位都没有!")
		errMsg = gem.ERR_ONE_KEY_USE_GEM_FAIL_NO_CAN_UPGRADE_GEM
		return
	}

	// 空的排在前面了
	sort.Sort(gemAndSlotSlice(mayUseOrUpgradeSlotArray))

	// 新的宝石队列
	newGemsArray = captain.CopyGems()

	// 使用了的宝石数量
	usedGemCountMap := map[uint64]uint64{}

	// 排序里面是从空->低级->高级的排序的，看有没有足够的可以替换
	for _, gemAndslotData := range mayUseOrUpgradeSlotArray {
		if gemAndslotData.slotIdx >= len(captain.Race().GemTypes) {
			continue
		}

		gemType := captain.Race().GemTypes[gemAndslotData.slotIdx]
		array := m.configDatas.GemDatas().GetLevelDescArrayByGemType(gemType)

		// 从高等级宝石到低等级宝石镶嵌宝石
		for _, gemData := range array {
			if gemAndslotData.gem != nil && gemData.Level <= gemAndslotData.gem.Level {
				// 该类型的宝石等级低的不可鞥替换高的
				break
			}

			// 看看背包中有没有这个宝石
			depotGemCount := heroDepot.GetGoodsCount(gemData.Id)
			if depotGemCount <= 0 {
				continue
			}

			// 看下够不够
			usedGemCount := usedGemCountMap[gemData.Id]
			leftGemCount := u64.Sub(depotGemCount, usedGemCount)
			if leftGemCount <= 0 {
				// 推荐下一颗等级的宝石
				continue
			}

			usedGemCountMap[gemData.Id] = usedGemCount + 1

			newGemsArray[gemAndslotData.slotIdx] = gemData
			// 镶嵌好了，下一个槽位
			break
		}

		// 执行到这里说明材料不够，或者都有镶嵌的有
	}

	for id, usedCount := range usedGemCountMap {
		if !heroDepot.HasEnoughGoods(id, usedCount) {
			logrus.Errorf("一键镶嵌宝石，物品不够，服务器怎么判断的？%v, %v", newGemsArray, usedGemCountMap)
			errMsg = gem.ERR_ONE_KEY_USE_GEM_FAIL_SERVER_ERROR
			return
		}

		// 多消耗一点，也总好过于出错吧
		usedCountCheck := uint64(0)
		for idx, newGem := range newGemsArray {
			if newGem == nil {
				continue
			}

			if newGem.Id != id {
				continue
			}

			oldGem := captain.GetGem(idx)
			if oldGem != nil && oldGem.Id == id {
				// 旧的已经是这个宝石了
				continue
			}

			usedCountCheck++
		}

		if usedCount != usedCountCheck {
			logrus.Errorf("一键镶嵌宝石，物品不够，服务器怎么判断的？%v, %v", newGemsArray, usedGemCountMap)
			errMsg = gem.ERR_ONE_KEY_USE_GEM_FAIL_SERVER_ERROR
			return
		}
	}

	return
}

type gemAndSlot struct {
	gem     *goods.GemData
	slotIdx int
}

type gemAndSlotSlice []*gemAndSlot

func (p gemAndSlotSlice) Len() int { return len(p) }
func (p gemAndSlotSlice) Less(i, j int) bool {
	pi, pj := p[i], p[j]
	if pi.gem == nil {
		if pj.gem == nil {
			return i < j
		}

		return true
	}

	if pj.gem == nil {
		return false
	}

	// 都有宝石
	if pi.gem.GemType != pj.gem.GemType {
		// 小的宝石类型放前面
		return pi.gem.GemType < pj.gem.GemType
	}

	// 类型相同小的等级放前面
	return pi.gem.Level < pj.gem.Level
}
func (p gemAndSlotSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// 一键合成宝石
//gogen:iface
func (m *GemModule) ProcessOneKeyCombineGem(proto *gem.C2SOneKeyCombineGemProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.GemOneKeyCombineGem)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		captainId := u64.FromInt32(proto.CaptainId)

		captain := hero.Military().Captain(captainId)
		if captain == nil {
			logrus.Debugf("合成宝石，武将id不存在")
			result.Add(gem.ERR_ONE_KEY_COMBINE_GEM_FAIL_INVALID_CAPTAIN_ID)
			return
		}

		// 出征武将不能操作宝石
		//if captain.IsOutSide() {
		//	logrus.Debugf("合成宝石，出征武将不能合成宝石")
		//	result.Add(gem.ERR_ONE_KEY_COMBINE_GEM_FAIL_CAPTAIN_OUTSIDE)
		//	return
		//}

		slotIdx := int(proto.GetSlotIdx())

		// 目标宝石
		targetGem := captain.GetGem(slotIdx)
		if targetGem == nil {
			logrus.Debugf("合成宝石，没有镶嵌宝石在上面")
			result.Add(gem.ERR_ONE_KEY_COMBINE_GEM_FAIL_INVALID_SLOT_IDX)
			return
		}

		newGem := targetGem.NextLevel
		if newGem == nil {
			logrus.Debugf("合成宝石，宝石没有下一级")
			result.Add(gem.ERR_ONE_KEY_COMBINE_GEM_FAIL_LEVEL_MAX)
			return
		}

		heroDepot := hero.Depot()

		costGemArray, costGemCountArray, err := calcOneKeyCombineGemCost(
			m.configDatas.GemDatas().GetLevelAscArrayByGemType(newGem.GemType),
			heroDepot, newGem, 1, targetGem, 1)
		if err != nil {
			switch err {
			case errCombineGemNotEnough:
				if !proto.Buy {
					result.Add(gem.ERR_ONE_KEY_COMBINE_GEM_FAIL_NOT_ENOUGH)
					return
				}

				firstLevelGem := m.configDatas.GetGemData(goods.GetGemId(newGem.GemType, 1))
				if firstLevelGem == nil || firstLevelGem.YuanbaoPrice <= 0 {
					logrus.Debug("一键合成宝石，未开放购买")
					result.Add(gem.ERR_ONE_KEY_COMBINE_GEM_FAIL_CANT_BUY)
					return
				}

				// 不够的花钱购买
				var totalFirstLevelCount uint64
				totalFirstLevelCount += targetGem.UpgradeToThisNeedFirstLevelCount
				for i, g := range costGemArray {
					count := costGemCountArray[i]
					totalFirstLevelCount += g.UpgradeToThisNeedFirstLevelCount * count
				}

				needFirstLevelCount := newGem.UpgradeToThisNeedFirstLevelCount
				buyCount := u64.Sub(needFirstLevelCount, totalFirstLevelCount)
				if buyCount <= 0 {
					// 不是说数量不够吗，为毛这里又够了
					logrus.Error("一键合成宝石，不是说数量不够吗，为毛这里又够了")
					result.Add(gem.ERR_ONE_KEY_COMBINE_GEM_FAIL_NOT_ENOUGH)
					return
				}

				totalCost := firstLevelGem.YuanbaoPrice * buyCount
				if !heromodule.ReduceYuanbao(hctx, hero, result, totalCost) {
					logrus.Debug("一键合成宝石（不足购买），消耗不足")
					result.Add(gem.ERR_ONE_KEY_COMBINE_GEM_FAIL_COST_NOT_ENOUGH)
					return
				}
			default:
				result.Add(gem.ERR_ONE_KEY_COMBINE_GEM_FAIL_SERVER_ERROR)
				return
			}
		}

		heromodule.ReduceGemArray(hctx, hero, result, costGemArray, costGemCountArray)
		captain.SetGem(slotIdx, newGem)

		// 更新武将属性，战斗力
		captain.CalculateProperties()
		result.Add(captain.NewUpdateCaptainStatMsg())
		heromodule.UpdateTroopFightAmount(hero, captain.GetTroop(), result)

		result.Add(gem.NewS2cOneKeyCombineGemMsg(proto.GetCaptainId(), proto.GetSlotIdx(), u64.Int32(newGem.Id)))

		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_GEM)

		// 系统广播
		if d := hctx.BroadcastHelp().GemGet; d != nil {
			heromodule.AddGemBroadcast(hctx, hero, newGem, result)
		}

		result.Changed()
		result.Ok()
	})
}

// 请求一键合成宝石消耗
//gogen:iface
func (m *GemModule) ProcessRequestOneKeyCombineGemCost(proto *gem.C2SRequestOneKeyCombineCostProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		var newGem, excludeGem *goods.GemData
		var excludeCount uint64

		captainId := u64.FromInt32(proto.CaptainId)
		if captainId == 0 {
			targetGem := m.configDatas.GetGemData(u64.FromInt32(proto.GemId))
			if targetGem == nil {
				logrus.Debugf("合成宝石消耗，宝石id不存在")
				result.Add(gem.ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_INVALID_SLOT_IDX)
				return
			}
			newGem = targetGem.NextLevel
		} else {
			captain := hero.Military().Captain(captainId)
			if captain == nil {
				logrus.Debugf("合成宝石消耗，武将id不存在")
				result.Add(gem.ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_INVALID_CAPTAIN_ID)
				return
			}

			// 出征武将不能操作宝石
			//if captain.IsOutSide() {
			//	logrus.Debugf("合成宝石，出征武将不能合成宝石")
			//	result.Add(gem.ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_CAPTAIN_OUTSIDE)
			//	return
			//}

			slotIdx := int(proto.GetSlotIdx())

			// 目标宝石
			targetGem := captain.GetGem(slotIdx)
			if targetGem == nil {
				logrus.Debugf("合成宝石消耗，没有镶嵌宝石在上面")
				result.Add(gem.ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_INVALID_SLOT_IDX)
				return
			}

			newGem = targetGem.NextLevel
			excludeGem = targetGem
			excludeCount = 1
		}

		if newGem == nil {
			logrus.Debugf("合成宝石消耗，宝石没有下一级")
			result.Add(gem.ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_LEVEL_MAX)
			return
		}

		heroDepot := hero.Depot()

		costGemArray, costGemCountArray, errMsg := calcOneKeyCombineGemCost(
			m.configDatas.GemDatas().GetLevelAscArrayByGemType(newGem.GemType),
			heroDepot, newGem, 1, excludeGem, excludeCount)
		if errMsg != nil {
			switch errMsg {
			case errCombineGemNotEnough:
				firstLevelGem := m.configDatas.GetGemData(goods.GetGemId(newGem.GemType, 1))
				if firstLevelGem == nil || firstLevelGem.YuanbaoPrice <= 0 {
					logrus.Debug("合成宝石消耗，未开放购买")
					result.Add(gem.ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_INVALID_SLOT_IDX)
					return
				}

				// 不够的花钱购买
				var totalFirstLevelCount uint64
				if excludeGem != nil {
					totalFirstLevelCount += excludeGem.UpgradeToThisNeedFirstLevelCount * excludeCount
				}

				for i, g := range costGemArray {
					count := costGemCountArray[i]
					totalFirstLevelCount += g.UpgradeToThisNeedFirstLevelCount * count
				}

				needFirstLevelCount := newGem.UpgradeToThisNeedFirstLevelCount
				buyCount := u64.Sub(needFirstLevelCount, totalFirstLevelCount)
				if buyCount <= 0 {
					// 不是说数量不够吗，为毛这里又够了
					logrus.Error("合成宝石消耗，不是说数量不够吗，为毛这里又够了")
				}

				totalCost := firstLevelGem.YuanbaoPrice * buyCount

				result.Add(gem.NewS2cRequestOneKeyCombineCostMsg(proto.GetCaptainId(), proto.GetSlotIdx(), false, nil, nil, u64.Int32(buyCount), u64.Int32(totalCost)))
			default:
				result.Add(gem.ERR_REQUEST_ONE_KEY_COMBINE_COST_FAIL_INVALID_SLOT_IDX)
				return
			}
		} else {
			//result.Add(gem.NewS2cRequestOneKeyCombineCostMsg(proto.GetCaptainId(), proto.GetSlotIdx(), true, u64.Int32Array(goods.GetGemDataKeyArray(costGemArray)), u64.Int32Array(costGemCountArray)))
			result.Add(gem.NewS2cRequestOneKeyCombineCostMsg(proto.GetCaptainId(), proto.GetSlotIdx(), true, nil, nil, 0, 0))
		}

		result.Ok()
	})
}

var errCombineGemServerError = errors.Errorf("calculate one key combine gem cost server error")
var errCombineGemNotEnough = errors.Errorf("calculate one key combine gem cost server error")

// 计算一键合成消耗的宝石
func calcOneKeyCombineGemCost(levelAscArray []*goods.GemData,
	heroDepot interface {
		GetGoodsCount(goodsId uint64) uint64
	}, targetGem *goods.GemData, targetCount uint64,
	excludeGem *goods.GemData, excludeCount uint64) (costGemArray []*goods.GemData, costGemCountArray []uint64, err error) {
	if targetCount <= 0 {
		logrus.Error("合成宝石，targetCount <= 0")
		err = errCombineGemServerError
		return
	}

	if excludeGem != nil {
		if targetGem.Level <= excludeGem.Level {
			logrus.WithField("target", targetGem.Id).WithField("target level", targetGem.Level).
				WithField("exlude", excludeGem.Id).WithField("exclude level", excludeGem.Level).
				Errorln("有排除掉的宝石，但是排除掉的宝石等级>=目标宝石等级")
			err = errCombineGemServerError
			return
		}

		if excludeCount <= 0 {
			logrus.WithField("target", targetGem.Id).WithField("target level", targetGem.Level).
				WithField("exlude", excludeGem.Id).WithField("exclude count", excludeCount).
				Errorln("有排除掉的宝石，但是排除掉的宝石数量<=0")
			err = errCombineGemServerError
			return
		}

		if targetGem.UpgradeToThisNeedFirstLevelCount*targetCount <= excludeGem.UpgradeToThisNeedFirstLevelCount*excludeCount {
			logrus.WithField("target", targetGem.Id).WithField("target level", targetGem.Level).
				WithField("exlude", excludeGem.Id).WithField("exclude count", excludeCount).
				Errorln("有排除掉的宝石，但是排除掉的宝石数量超出了升级到目标宝石的数量")
			err = errCombineGemServerError
			return
		}
	}

	if targetGem.PrevLevel == nil {
		logrus.WithField("target", targetGem.Id).WithField("target level", targetGem.Level).
			Errorln("要合成的目标宝石没有前置宝石")
		err = errCombineGemServerError
		return
	}

	costGemArray = levelAscArray
	costGemArray = costGemArray[:targetGem.PrevLevel.Level]

	needFirstLevelCount := targetGem.UpgradeToThisNeedFirstLevelCount * targetCount

	var totalFirstLevelCount uint64
	if excludeGem != nil {
		totalFirstLevelCount += excludeGem.UpgradeToThisNeedFirstLevelCount * excludeCount

		if totalFirstLevelCount >= needFirstLevelCount {
			// 在检查一下
			logrus.WithField("target", targetGem.Id).WithField("target level", targetGem.Level).
				WithField("exlude", excludeGem.Id).WithField("exclude count", excludeCount).
				Errorln("有排除掉的宝石，但是排除掉的宝石数量超出了升级到目标宝石的数量")
			err = errCombineGemServerError
			return
		}
	}

	costGemCountArray = make([]uint64, len(costGemArray))
	for idx, gem := range costGemArray {
		// 这里呢有这么多
		count := heroDepot.GetGoodsCount(gem.Id)

		firstLevelCount := gem.UpgradeToThisNeedFirstLevelCount * count
		if totalFirstLevelCount+firstLevelCount >= needFirstLevelCount {
			// 数量足够了，加到刚好满，不要多
			canAddCount := u64.Sub(needFirstLevelCount, totalFirstLevelCount)
			toAddCount := u64.DivideTimes(canAddCount, gem.UpgradeToThisNeedFirstLevelCount)

			costGemCountArray[idx] += toAddCount
			totalFirstLevelCount += gem.UpgradeToThisNeedFirstLevelCount * toAddCount
			break
		}
		costGemCountArray[idx] += count
		totalFirstLevelCount += firstLevelCount
	}

	if totalFirstLevelCount < needFirstLevelCount {
		// 合并宝石数量不足
		logrus.Debugf("合成宝石，宝石不够")
		err = errCombineGemNotEnough

		// 对1级宝石做处理，1级宝石只提供可以合成2级的数量，不提供额外数量
		// 比如合成3级宝石，还差1个2级宝石，但是1级宝石只有1个，不能合成2级宝石，那么这个1级宝石就不用在合成操作
		//if len(costGemCountArray) > 0 {
		//	if upgradeNeedCount := levelAscArray[0].UpgradeNeedCount; upgradeNeedCount > 0 {
		//		costGemCountArray[0] = costGemCountArray[0] / upgradeNeedCount * upgradeNeedCount
		//	}
		//}
		return
	}

	// 超出所需个数，从最高级的那一档开始减少个数
	if totalFirstLevelCount > needFirstLevelCount {
		diff := u64.Sub(totalFirstLevelCount, needFirstLevelCount)

		for i := len(costGemCountArray) - 1; i >= 0; i-- {
			g := costGemArray[i]
			count := costGemCountArray[i]

			if count > 0 && diff >= g.UpgradeToThisNeedFirstLevelCount {
				toReduceCount := diff / g.UpgradeToThisNeedFirstLevelCount
				toReduceCount = u64.Min(toReduceCount, count)

				diff = u64.Sub(diff, g.UpgradeToThisNeedFirstLevelCount*toReduceCount)
				costGemCountArray[i] = u64.Sub(count, toReduceCount)

				if diff <= 0 {
					break
				}
			}
		}

		if diff > 0 {
			logrus.Error("合成宝石，diff > 0")
			err = errCombineGemServerError
			return
		}

	}

	return
}

// 一键合成宝石
//gogen:iface
func (m *GemModule) ProcessOneKeyCombineDepotGem(proto *gem.C2SOneKeyCombineDepotGemProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.GemOneKeyCombineDepotGem)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		// 目标宝石
		targetGem := m.configDatas.GetGemData(uint64(proto.GemId))
		if targetGem == nil {
			logrus.Debugf("一键合成背包宝石，宝石id不存在")
			result.Add(gem.ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_INVALID_GEM_ID)
			return
		}

		newGem := targetGem.NextLevel
		if newGem == nil {
			logrus.Debugf("一键合成背包宝石，宝石没有下一级")
			result.Add(gem.ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_LEVEL_MAX)
			return
		}

		newGemCount := u64.Max(u64.FromInt32(proto.NewGemCount), 1)

		heroDepot := hero.Depot()

		costGemArray, costGemCountArray, err := calcOneKeyCombineGemCost(
			m.configDatas.GemDatas().GetLevelAscArrayByGemType(newGem.GemType),
			heroDepot, newGem, newGemCount, nil, 0)
		// 如果个数不够，并且
		if err != nil {
			switch err {
			case errCombineGemNotEnough:
				if !proto.Buy {
					// 数量不够，也不买，返回错误码
					result.Add(gem.ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_NOT_ENOUGH)
					return
				}

				if newGemCount > 1 {
					logrus.Debug("一键合成背包宝石，不允许购买并合成多个")
					result.Add(gem.ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_NOT_ENOUGH)
					return
				}

				firstLevelGem := m.configDatas.GetGemData(goods.GetGemId(newGem.GemType, 1))
				if firstLevelGem == nil || firstLevelGem.YuanbaoPrice <= 0 {
					logrus.Debug("一键合成背包宝石，未开放购买")
					result.Add(gem.ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_CANT_BUY)
					return
				}

				// 不够的花钱购买
				var totalFirstLevelCount uint64
				for i, g := range costGemArray {
					count := costGemCountArray[i]
					totalFirstLevelCount += g.UpgradeToThisNeedFirstLevelCount * count
				}

				needFirstLevelCount := newGem.UpgradeToThisNeedFirstLevelCount * newGemCount
				buyCount := u64.Sub(needFirstLevelCount, totalFirstLevelCount)
				if buyCount <= 0 {
					// 不是说数量不够吗，为毛这里又够了
					logrus.Error("一键合成背包宝石，不是说数量不够吗，为毛这里又够了")
					result.Add(gem.ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_NOT_ENOUGH)
					return
				}

				totalCost := firstLevelGem.YuanbaoPrice * buyCount
				if !heromodule.ReduceYuanbao(hctx, hero, result, totalCost) {
					logrus.Debug("一键合成背包宝石（不足购买），消耗不足")
					result.Add(gem.ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_COST_NOT_ENOUGH)
					return
				}

			default:
				result.Add(gem.ERR_ONE_KEY_COMBINE_DEPOT_GEM_FAIL_SERVER_ERROR)
				return
			}
		}

		heromodule.ReduceGemArray(hctx, hero, result, costGemArray, costGemCountArray)
		heromodule.AddGem(hctx, hero, result, newGem, newGemCount, true)

		result.Add(gem.NewS2cOneKeyCombineDepotGemMsg(u64.Int32(newGem.Id), proto.NewGemCount))

		result.Changed()
		result.Ok()
	})
}
