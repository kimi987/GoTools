package combatx

import (
	"fmt"
	"github.com/lightpaw/male7/pb/shared_proto"
	"sort"
)

func PrintResult(proto *shared_proto.CombatXProto) {

	fmt.Println("攻方胜利：", proto.AttackerWin, "最大帧：", proto.MaxFrame)

	player := make(map[int32]int32)
	fmt.Println("攻方信息：")
	for i, t := range proto.Attacker.Troops {
		data := proto.AttackerTroopData[i]
		idx := data.Index
		fmt.Println(idx, t.FightIndex, t.Captain.CaptainId, t.Captain.Name, t.Captain.Soldier*t.Captain.LifePerSoldier, t.Captain.Race)
		fmt.Println(getIndentation(1), "坐标", data.X, data.Y, "怒气", data.Rage, data.Recover)
		player[idx] = t.Captain.Soldier * t.Captain.LifePerSoldier
	}

	fmt.Println("防方信息：")
	for i, t := range proto.Defenser.Troops {
		data := proto.DefenserTroopData[i]
		idx := data.Index
		fmt.Println(idx, t.FightIndex, t.Captain.CaptainId, t.Captain.Name, t.Captain.Soldier*t.Captain.LifePerSoldier, t.Captain.Race)
		fmt.Println(getIndentation(1), "坐标", data.X, data.Y, "怒气", data.Rage, data.Recover)
		player[idx] = t.Captain.Soldier * t.Captain.LifePerSoldier
	}

	actions := make([]*shared_proto.CombatFrameProto, len(proto.Frame))
	copy(actions, proto.Frame)
	sort.Sort(FrameSlice(actions))

	for _, act := range actions {
		frame := act.Frame
		fmt.Println("帧：", act.Frame, len(act.Action))

		for _, act := range act.Action {

			switch act.ActionType {
			case ActionTypeMove, ActionTypeShortMove:
				if act.Move == nil {
					fmt.Println("act.Move == nil")
					break
				}

				index := act.Index
				act := act.Move
				if act.MoveEndFrame != 0 {
					fmt.Println(index, "移动", act.MoveStartX, act.MoveStartY, "->", act.MoveEndX, act.MoveEndY, "到达帧", frame, "->", act.MoveEndFrame)
				}

			case ActionTypeReleaseSpell:
				if act.ReleaseSpell == nil {
					fmt.Println("act.ReleaseSpell == nil")
					break
				}

				index := act.Index
				act := act.ReleaseSpell
				fmt.Println(index, "释放技能", act.SpellId, "施法类型", act.ReleaseType, "结束帧", act.EndFrame, "施法坐标", act.SpellX, act.SpellY)

				for _, target := range act.Target {
					fmt.Println(getIndentation(1), target.HurtType, "目标", target.TargetIndex, "扣盾", target.ChangeShield, "死兵", target.ChangeSoldier, "剩余士兵：", target.Soldier, "坐标", target.TargetX, target.TargetY)
					if len(target.AddState) > 0 || len(target.RemoveState) > 0 {
						fmt.Println(getIndentation(1), "目标", target.TargetIndex, "状态加减", target.AddState, target.RemoveState)
					}
					if target.Rage != nil {
						fmt.Println(getIndentation(1), "目标", target.TargetIndex, "怒气更新", target.Rage.Rage, target.Rage.Recover)
					}
				}

			case ActionTypeSpellEffect:
				if act.SpellEffect == nil {
					fmt.Println("act.SpellEffect == nil")
					break
				}

				index := act.Index
				act := act.SpellEffect
				fmt.Println(index, "技能延时生效，技能:", act.Id, " 打", index, "扣盾", act.ChangeShield, "死兵", act.ChangeSoldier, "剩余士兵：", act.Soldier)
				if len(act.RemoveState) > 0 {
					fmt.Println(getIndentation(1), "目标", index, "状态加减", act.RemoveState)
				}
				if act.Rage != nil {
					fmt.Println(getIndentation(1), "目标", index, "怒气更新", act.Rage.Rage, act.Rage.Recover)
				}

			case ActionTypeStateEffect:
				if act.StateEffect == nil {
					fmt.Println("act.StateEffect == nil")
					break
				}

				index := act.Index
				act := act.StateEffect
				fmt.Println(index, "状态生效，状态:", act.Id, " 打", index, "扣盾", act.ChangeShield, "死兵", act.ChangeSoldier, "剩余士兵：", act.Soldier)
				if len(act.RemoveState) > 0 {
					fmt.Println(getIndentation(1), "目标", index, "状态加减", act.RemoveState)
				}
				if act.Rage != nil {
					fmt.Println(getIndentation(1), "目标", index, "怒气更新", act.Rage.Rage, act.Rage.Recover)
				}

			case ActionTypeTriggerPassiveSpell:
				if act.TriggerPassiveSpell == nil {
					fmt.Println("act.TriggerPassiveSpell == nil")
					break
				}

				index := act.Index
				act := act.TriggerPassiveSpell
				fmt.Println(index, "触发", act.SelfIndex, act.TargetIndex, "的被动:", act.PassiveSpellId)

				if len(act.SelfAddState) > 0 {
					fmt.Println(getIndentation(1), "目标", act.SelfIndex, "添加self状态", act.SelfAddState)
				}

				if len(act.TargetAddState) > 0 {
					fmt.Println(getIndentation(1), "目标", act.TargetIndex, "添加target状态", act.TargetAddState)
				}

				if len(act.TargetStateEffect) > 0 {
					fmt.Println(getIndentation(1), act.TargetIndex, "持续技能立即生效")
					for _, effect := range act.TargetStateEffect {
						fmt.Println(getIndentation(2), "Id", effect.Id, "扣盾", effect.ChangeShield, "死兵", effect.ChangeSoldier, "剩余士兵", effect.Soldier)

						if len(effect.RemoveState) > 0 {
							fmt.Println(getIndentation(2), "目标", act.TargetIndex, "状态加减", effect.RemoveState)
						}
						if effect.Rage != nil {
							fmt.Println(getIndentation(2), "目标", act.TargetIndex, "怒气更新", effect.Rage.Rage, effect.Rage.Recover)
						}
					}
				}

			default:
				fmt.Println("invalid action type", act.ActionType)
			}

		}

	}

	fmt.Println("存活士兵：")
	for _, s := range proto.AliveSolider {
		fmt.Println(s.Key, s.Value)
	}
}

const indentation = "  "

func getIndentation(layout int) (s string) {
	for i := 0; i < layout; i++ {
		s += indentation
	}
	return
}

//func PrintMultiResult(result *shared_proto.MultiCombatProto) {
//
//	leftSoldierLifeMap := make(map[int32]int32)
//	leftSoldierMap := make(map[int32]int32)
//	attackerIdMap := make(map[int64]int32)
//	defenserIdMap := make(map[int64]int32)
//
//	for _, c := range result.Combats {
//		id, _ := idbytes.ToId(c.Attacker.Hero.Id)
//		if _, exist := attackerIdMap[id]; !exist {
//			attackerIdMap[id] = 1
//			for i, t := range c.Attacker.Troops {
//				idx := c.AttackerTroopPos[i].Index
//				leftSoldierLifeMap[idx] = t.Captain.Soldier * t.Captain.LifePerSoldier
//				leftSoldierMap[idx] = t.Captain.Soldier
//			}
//		}
//
//		id, _ = idbytes.ToId(c.Defenser.Hero.Id)
//		if _, exist := defenserIdMap[id]; !exist {
//			defenserIdMap[id] = 1
//			for i, t := range c.Defenser.Troops {
//				idx := c.DefenserTroopPos[i].Index
//				leftSoldierLifeMap[idx] = t.Captain.Soldier * t.Captain.LifePerSoldier
//				leftSoldierMap[idx] = t.Captain.Soldier
//			}
//		}
//	}
//
//	for _, f := range result.GetCombats() {
//		PrintFight(leftSoldierLifeMap, leftSoldierMap, result, f)
//
//		attackerId, _ := idbytes.ToId(f.Attacker.Hero.Id)
//		defenserId, _ := idbytes.ToId(f.Defenser.Hero.Id)
//
//		if f.AttackerWin {
//			delete(defenserIdMap, defenserId)
//			if f.IsWinnerContinueWinLeave {
//				delete(attackerIdMap, attackerId)
//			}
//		} else {
//			delete(attackerIdMap, attackerId)
//			if f.IsWinnerContinueWinLeave {
//				delete(defenserIdMap, defenserId)
//			}
//		}
//	}
//
//	fmt.Println("---------------------------")
//	fmt.Println("最终胜利方是进攻方: ", result.AttackerWin)
//	fmt.Println("进攻方剩余部队数: ", len(attackerIdMap))
//	fmt.Println("防守方剩余部队数: ", len(defenserIdMap))
//}
//
//func PrintFight(leftSoldierLifeMap, leftSoldierMap map[int32]int32, combat *shared_proto.MultiCombatProto, proto *shared_proto.CombatProto) {
//	fmt.Println("------------------------------------------")
//	fmt.Println("攻方id:", proto.Attacker.Hero.Id)
//	fmt.Println("防方id:", proto.Defenser.Hero.Id)
//	fmt.Println("攻方胜利：", proto.AttackerWin)
//	fmt.Println("轮次：", len(proto.Rounds))
//	fmt.Println("战斗所在队列", proto.GetFightQueue())
//
//	fmt.Println("攻方信息：")
//
//	for _, attacker := range combat.Attacker {
//		if bytes.Equal(attacker.Hero.Id, proto.Attacker.Hero.Id) {
//			for i, t := range attacker.Troops {
//				idx := proto.AttackerTroopPos[i].Index
//				fmt.Println(idx, t.Captain.Name, leftSoldierMap[idx], leftSoldierLifeMap[idx])
//			}
//		}
//	}
//
//	fmt.Println("防方信息：")
//	for _, defenser := range combat.Defenser {
//		if bytes.Equal(defenser.Hero.Id, proto.Defenser.Hero.Id) {
//			for i, t := range defenser.Troops {
//				idx := proto.DefenserTroopPos[i].Index
//				fmt.Println(idx, t.Captain.Name, leftSoldierMap[idx], leftSoldierLifeMap[idx])
//			}
//		}
//	}
//
//	for i, round := range proto.Rounds {
//		fmt.Println("轮次：", i+1)
//
//		for _, act := range round.Actions {
//			if act.MoveDirection != shared_proto.Direction_ORIGIN {
//				fmt.Println("\t", act.Index, "移动", act.MoveDirection)
//			}
//
//			if act.TargetIndex > 0 {
//				oldLife := leftSoldierLifeMap[act.TargetIndex]
//				leftSoldierLifeMap[act.TargetIndex] = oldLife - act.Damage
//				fmt.Println("\t", act.Index, "打", act.TargetIndex, act.HurtType, act.Damage, "此前血量：", oldLife, "剩余血量：", leftSoldierLifeMap[act.TargetIndex])
//			}
//		}
//	}
//
//	if proto.AttackerWin {
//		fmt.Printf("进攻方 [%s] 存活士兵: \n", proto.Attacker.Hero.Id)
//	} else {
//		fmt.Printf("防守方 [%s] 存活士兵: \n", proto.Defenser.Hero.Id)
//	}
//
//	for _, s := range proto.AliveSolider {
//		fmt.Println(s.Key, s.Value, leftSoldierLifeMap[s.Key])
//
//		leftSoldierMap[s.Key] = s.Value
//	}
//
//	fmt.Println("胜利方连胜次数: ", proto.WinnerContinueWinTimes, " 胜利方是否连胜离场: ", proto.IsWinnerContinueWinLeave)
//}
