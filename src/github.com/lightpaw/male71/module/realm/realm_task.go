package realm

import (
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/logrus"
	"math/rand"
	"github.com/lightpaw/male7/gen/pb/region"
)

func (r *Realm) UpdateHeroRealmInfo(heroId int64, checkHomeNpcTask, checkMonsterTask, checkTroopOutside bool) {
	if base := r.getBase(heroId); base == nil {
		return
	}

	r.queueFunc(false, func() {

		base := r.getBase(heroId)
		if base == nil {
			return
		}

		home := GetHomeBase(base)
		if home == nil {
			return
		}

		var addFigntNpcFunc func()
		r.heroFuncWithSend(heroId, func(hero *entity.Hero, result herolock.LockResult) {
			addFigntNpcFunc = r.updateHeroRealmInfo(hero, result, base, home, checkHomeNpcTask, checkMonsterTask, checkTroopOutside)

		})

		if addFigntNpcFunc != nil {
			addFigntNpcFunc()
		}

	})
}

func (r *Realm) updateHeroRealmInfo(hero *entity.Hero, result herolock.LockResult, base *baseWithData, home *heroHome,
	updateHomeNpc, isCheckTaskMonster, checkTroopOutside bool) (addFigntNpcFunc func()) {

	if updateHomeNpc {
		r.updateHomeNpcTask(hero, result, home)
	}

	if isCheckTaskMonster {
		addFigntNpcFunc = r.checkTaskMonster(hero, result, base)
	}

	if checkTroopOutside {
		r.updateTroopOutside(hero, result, base)
	}

	return
}

func (r *Realm) updateHomeNpcTask(hero *entity.Hero, result herolock.LockResult, home *heroHome) {

	heroTaskList := hero.TaskList()

	update := false
	heroTaskList.WalkAllTask(func(task entity.Task) (endedWalk bool) {
		if task.Progress().IsCompleted() {
			return false
		}

		if task.Progress().TargetType() == shared_proto.TaskTargetType_TASK_TARGET_KILL_HOME_NPC {
			if npcData := task.Progress().Target().KillHomeNpcData; npcData != nil {
				targetId := npcid.NewHomeNpcId(npcData.Id)
				if hero.GetHomeNpcBase(targetId) != nil && home.homeNpcBase[targetId] == nil {
					hero.RemoveHomeNpcBase(targetId)

					update = true
				}
			}
		}

		return false
	})

	if update {
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_KILL_HOME_NPC)
	}
}

func (r *Realm) checkTaskMonster(hero *entity.Hero, result herolock.LockResult, base *baseWithData) (addFigntNpcFunc func()) {

	var monIds []uint64
	for _, t := range base.targetingTroops {
		if t.mmd == nil {
			continue
		}

		if t.mmd.GetSpec() != monsterdata.InvasionTask {
			continue
		}

		if !t.State().IsInvateState() {
			continue
		}

		monIds = append(monIds, t.mmd.Id)
	}

	blockX, blockY := r.mapData.GetBlockByPos(base.BaseX(), base.BaseY())
	sequence := regdata.BlockSequence(blockX, blockY)
	var startBase []*baseWithData
	for _, data := range r.levelData.GetMultiLevelMonsters() {
		id := npcid.GetNpcId(sequence, data.Id, npcid.NpcType_MultiLevelMonster)

		base := r.getBase(id)
		if base == nil {
			logrus.WithField("x", blockX).WithField("y", blockY).
				WithField("data", data.Id).
				Error("检查任务野怪，根据id找不到野怪")
			continue
		}

		// 怪物太多了，换一个
		if len(base.selfTroops) >= npcid.TroopMaxSequence {
			continue
		}

		startBase = append(startBase, base)
	}

	var troopIds []int64
	var startingBases []*baseWithData
	var monsterDatas []*monsterdata.MonsterMasterData

	update := false
	hero.TaskList().RangeTaskMonsterId(func(monsterId uint64, created bool) {

		if created {
			// 已经创建，检查是否已经完成
			if !u64.Contains(monIds, monsterId) {
				// 怪已经不见了，完成任务

				// 移除任务怪
				hero.TaskList().RemoveTaskMonster(monsterId)

				// 添加任务怪
				hero.HistoryAmount().IncreaseOneWithSubType(server_proto.HistoryAmountType_ExpelFightMonster, monsterId)

				update = true
			}

		} else {
			// 未创建，这里创建出来吧
			monsterData := r.services.datas.GetMonsterMasterData(monsterId)
			if monsterData == nil {
				logrus.Error("英雄检查任务怪物时候，monsterData == nil")
				hero.TaskList().RemoveTaskMonster(monsterId)
				return
			}

			if u64.Contains(monIds, monsterId) {
				logrus.Error("英雄检查任务怪物时候，create=false，但是已经有这个怪")
				hero.TaskList().SetTaskMonster(monsterId, true)
				return
			}

			n := len(startBase)
			if n > 0 {
				offset := rand.Int()
				for i := 0; i < n; i++ {
					idx := (i + offset) % n
					startingBase := startBase[idx]

					newTroopId := startingBase.nextNpcTroopId()
					if newTroopId != 0 {
						if !i64.Contains(troopIds, newTroopId) {
							// 添加野怪
							hero.TaskList().SetTaskMonster(monsterId, true)

							troopIds = append(troopIds, newTroopId)
							startingBases = append(startingBases, startingBase)
							monsterDatas = append(monsterDatas, monsterData)
							break
						}
					}
				}
			}
		}

	})

	if update {
		// 更新任务
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_EXPEL_FIGHT_MONSTER)
	}

	if len(troopIds) > 0 {
		addFigntNpcFunc = func() {
			for i, newTroopId := range troopIds {
				r.addInvasionFightNpcTroop(newTroopId, startingBases[i], base, 1, monsterDatas[i],
					r.services.datas.TaskMiscData().TaskMonsterArriveOffset)
			}
		}
	}

	return
}

func (r *Realm) updateTroopOutside(hero *entity.Hero, result herolock.LockResult, base *baseWithData) {

	// 检查一下部队状态，如果队伍状态跟里面的状态不一致，更新玩家数据
	for _, t := range hero.Troops() {
		if t == nil {
			continue
		}

		ii := t.GetInvateInfo()
		selfTroop := base.selfTroops[t.Id()]

		if ii != nil {
			if selfTroop == nil {

				var targetName interface{} = ii.OriginTargetID()
				if target := r.getBase(ii.OriginTargetID()); target != nil {
					targetName = target.Base().HeroName()
				}

				// 状态不对，打印日志
				logrus.WithField("name", hero.Name()).Errorf("检查玩家部队出征状态异常 ii != nil，target: %v state:%v time: %v", targetName, ii.State(), ii.CreateTime())

				// 部队出征状态不对，移除数据
				t.RemoveInvateInfo()
				result.Add(region.NewS2cUpdateSelfTroopsOutsideMsg(u64.Int32(entity.GetTroopIndex(t.Id()))+1, false))
			}
		} else {
			if selfTroop != nil {
				var targetName interface{} = selfTroop.originTargetId
				if selfTroop.targetBase != nil {
					targetName = selfTroop.targetBase.Base().HeroName()
				} else {
					if target := r.getBase(selfTroop.originTargetId); target != nil {
						targetName = target.Base().HeroName()
					}
				}

				// 状态不对，打印日志
				logrus.WithField("name", hero.Name()).Errorf("检查玩家部队出征状态异常 ii == nil，target: %v state:%v time: %v", targetName, selfTroop.State(), selfTroop.CreateTime())
			}
		}
	}

}
