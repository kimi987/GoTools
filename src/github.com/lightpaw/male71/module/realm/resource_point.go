package realm

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/entity/hexagon"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/util/must"
)

// 处理下面几种情况导致的资源点冲突，添加阻挡，移除阻挡
// 进场景
func (r *Realm) addHomeResourcePointBlock(base *baseWithData, updateFarm bool) {
	if !base.isHeroHomeBase() {
		return
	}

	evenOffsetCubes := r.services.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(base.BaseLevel())
	if len(evenOffsetCubes) <= 0 {
		return
	}

	var targetCubes []cb.Cube
	for _, offset := range evenOffsetCubes {
		evenOffsetX, evenOffsetY := offset.XY()
		targetCube := hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), evenOffsetX, evenOffsetY)
		targetCubes = append(targetCubes, targetCube)
	}

	r.updateHomeResourcePointBlock(base.Id(), true, targetCubes, updateFarm)
}

// 出场景
func (r *Realm) removeHomeResourcePointBlock(base *baseWithData) {
	if !base.isHeroHomeBase() {
		return
	}

	evenOffsetCubes := r.services.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(base.BaseLevel())
	if len(evenOffsetCubes) <= 0 {
		return
	}

	var targetCubes []cb.Cube

	for _, offset := range evenOffsetCubes {
		evenOffsetX, evenOffsetY := offset.XY()
		targetCube := hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), evenOffsetX, evenOffsetY)
		targetCubes = append(targetCubes, targetCube)
	}

	r.updateHomeResourcePointBlock(base.Id(), false, targetCubes, true)
}

// 场景内迁移
func (r *Realm) updateHomeResourcePointBlockWhenPosChanged(base realmface.Base, originX, originY int) {
	if base.BaseType() != realmface.BaseTypeHome {
		return
	}

	if base.BaseX() == originX && base.BaseY() == originY {
		logrus.Error("Realm.updateHomeResourcePointBlockWhenPosChanged base.BaseX() == originX && base.BaseY() == originY")
		return
	}

	evenOffsetCubes := r.services.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(base.GetBaseLevel())
	if len(evenOffsetCubes) <= 0 {
		return
	}

	var toAdds, toRemoves []cb.Cube
	for _, offset := range evenOffsetCubes {
		evenOffsetX, evenOffsetY := offset.XY()

		toRemove := hexagon.ShiftEvenOffset(originX, originY, evenOffsetX, evenOffsetY)
		toRemoves = append(toRemoves, toRemove)

		toAdd := hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), evenOffsetX, evenOffsetY)
		toAdds = append(toAdds, toAdd)
	}

	r.removeAddHomeResourcePointBlock(base.Id(), toAdds, toRemoves, true)
}

// 升级/降级
func (r *Realm) updateHomeResourcePointBlockWhenLevelChanged(base *baseWithData, originLevel uint64) {
	if !base.isHeroHomeBase() {
		return
	}

	if base.BaseLevel() == originLevel {
		logrus.Error("Realm.updateHomeResourcePointBlockWhenLevelChanged base.BaseLevel() == originLevel")
		return
	}

	minLevel := u64.Min(base.BaseLevel(), originLevel)
	maxLevel := u64.Max(base.BaseLevel(), originLevel)

	isAdd := originLevel < base.BaseLevel()

	var evenOffsetCubes []cb.Cube
	for i := minLevel + 1; i <= maxLevel; i++ {
		cubes := r.services.datas.RegionConfig().GetEvenOffsetCubesOnlyCurrentLevel(i)
		evenOffsetCubes = append(evenOffsetCubes, cubes...)
	}

	if len(evenOffsetCubes) <= 0 {
		return
	}

	var targetCubes []cb.Cube
	for _, offset := range evenOffsetCubes {
		evenOffsetX, evenOffsetY := offset.XY()
		targetCube := hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), evenOffsetX, evenOffsetY)
		targetCubes = append(targetCubes, targetCube)
	}

	r.updateHomeResourcePointBlock(base.Id(), isAdd, targetCubes, true)
}

// 联盟变更，导致地块冲突变更
func (r *Realm) updateHomeResourcePointBlockWhenGuildChanged(base *baseWithData) {
	if !base.isHeroHomeBase() {
		return
	}

	var evenOffsetCubes = r.services.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(base.BaseLevel())

	if len(evenOffsetCubes) <= 0 {
		return
	}

	var targetCubes []cb.Cube
	for _, offset := range evenOffsetCubes {
		evenOffsetX, evenOffsetY := offset.XY()
		targetCube := hexagon.ShiftEvenOffset(base.BaseX(), base.BaseY(), evenOffsetX, evenOffsetY)
		targetCubes = append(targetCubes, targetCube)
	}

	r.updateHomeResourcePointBlock(base.Id(), true, targetCubes, true)
}

func (r *Realm) updateHomeResourcePointBlock(heroId int64, isAdd bool, targetCubes []cb.Cube, updateFarm bool) {
	if isAdd {
		r.removeAddHomeResourcePointBlock(heroId, targetCubes, nil, updateFarm)
	} else {
		r.removeAddHomeResourcePointBlock(heroId, nil, targetCubes, updateFarm)
	}
}

func (r *Realm) removeAddHomeResourcePointBlock(heroId int64, toAdds, toRemoves []cb.Cube, updateFarm bool) {

	var affectHeroIds = []int64{heroId}
	var targetCubes []cb.Cube

	if len(toRemoves) > 0 {
		for _, targetCube := range toRemoves {
			if !cb.Contains(targetCubes, targetCube) {
				targetCubes = append(targetCubes, targetCube)
			}

			heroIds := r.resourceConflictHeroMap[targetCube]
			for _, heroId := range heroIds {
				affectHeroIds = i64.AddIfAbsent(affectHeroIds, heroId)
			}

			heroIds = i64.RemoveIfPresent(heroIds, heroId)

			if len(heroIds) > 0 {
				r.resourceConflictHeroMap[targetCube] = heroIds
			} else {
				delete(r.resourceConflictHeroMap, targetCube)
			}
		}
	}

	if len(toAdds) > 0 {
		for _, targetCube := range toAdds {
			if !cb.Contains(targetCubes, targetCube) {
				targetCubes = append(targetCubes, targetCube)
			}

			heroIds := r.resourceConflictHeroMap[targetCube]
			for _, heroId := range heroIds {
				affectHeroIds = i64.AddIfAbsent(affectHeroIds, heroId)
			}

			heroIds = i64.AddIfAbsent(heroIds, heroId)
			r.resourceConflictHeroMap[targetCube] = heroIds
		}
	}

	r.updateResourcePointConflict(affectHeroIds, targetCubes, updateFarm)
}

// 资源点冲突更新
func (r *Realm) updateResourcePointConflict(heroIds []int64, blocks []cb.Cube, updateFarm bool) {

	ctime := r.services.timeService.CurrentTime()

	// 每个人都lock一次，然后看下哪个cube是需要自己处理的
	for _, heroId := range heroIds {
		base := r.getBase(heroId)
		if base == nil || !base.isHeroHomeBase() {
			continue
		}

		var baseLevel uint64
		var npcConflictOffsets []cb.Cube
		r.heroBaseFuncWithSend(base.Base(), func(hero *entity.Hero, result herolock.LockResult) {
			baseLevel = hero.BaseLevel()
			r.updateHeroResourcePointConflict(blocks, base.BaseX(), base.BaseY(), hero, result, ctime)
			npcConflictOffsets = hero.GetNpcConflictResourcePointOffset()
		})

		if updateFarm {
			allConflictedBlocks := r.getAllConflictedBlock(heroId, blocks)
			r.updateFarm(heroId, baseLevel, blocks, allConflictedBlocks, npcConflictOffsets, base.BaseX(), base.BaseY(), ctime)
		}
	}
}

// 野怪npc造成的冲突地块
func (r *Realm) updateHeroResourcePointConflict(blocks []cb.Cube, baseX int, baseY int, hero *entity.Hero, result herolock.LockResult, ctime time.Time) {
	var layoutIds []uint64

	for _, block := range blocks {
		targetX, targetY := block.XY()
		evenOffset := hexagon.EvenOffsetBetween(baseX, baseY, targetX, targetY)
		layoutData := r.services.datas.RegionConfig().GetLayoutDataByEvenOffset(evenOffset)
		if layoutData == nil {
			continue
		}

		// 看下这个位置自己是不是有冲突
		conflicted := r.isResourcePointConflicted(hero.Id(), block)

		if hero.TrySetHeroResourcePointConflicted(layoutData, conflicted, ctime) {
			layoutIds = append(layoutIds, layoutData.Id)
		}
	}

	result.Add(domestic.RESOURCE_CONFLICT_CHANGED_S2C)
	result.Add(domestic.NewS2cResourcePointChangeV2Msg(must.Marshal(hero.EncodeResourcePointV2(r.services.datas))))
	heromodule.SendResourcePointUpdateMsgByLayoutIds(hero, result, layoutIds)
}

func (r *Realm) updateFarm(heroId int64, baseLevel uint64, blocks []cb.Cube, allConflictedBlocks, npcConflictOffsets []cb.Cube, baseX int, baseY int, ctime time.Time) {

	if heroId == 0 || npcid.IsNpcId(heroId) {
		return
	}

	// 农场
	r.services.farmService.FuncNoWait(func() {
		r.services.farmService.UpdateFarmCubes(heroId, baseLevel, blocks, allConflictedBlocks, npcConflictOffsets, baseX, baseY, ctime)
	})
}

func (r *Realm) updateFarmWithNpc(heroId int64, baseLevel uint64, baseX, baseY int, offsets []cb.Cube, isAdd bool, ctime time.Time) {
	if heroId == 0 || npcid.IsNpcId(heroId) || len(offsets) <= 0 {
		return
	}

	realOffsets := make([]cb.Cube, 0)
	for _, offset := range offsets {
		x, y := offset.XY()
		cube := hexagon.ShiftEvenOffset(baseX, baseY, x, y)
		if !isAdd && r.isResourcePointConflicted(heroId, cube) {
			// 要解除 npc 冲突的地块，还在资源冲突状态
			continue
		}
		realOffsets = append(realOffsets, offset)
	}

	r.services.farmService.FuncNoWait(func() {
		r.services.farmService.UpdateFarmCubeWithOffset(heroId, baseLevel, baseX, baseY, realOffsets, isAdd, ctime)
	})
}

func (r *Realm) getAllConflictedBlock(heroId int64, blocks []cb.Cube) (result []cb.Cube) {
	for _, block := range blocks {
		if r.isResourcePointConflicted(heroId, block) {
			result = append(result, block)
		}
	}

	return result
}
