package combat

import (
	"bytes"
	"fmt"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/idbytes"
)

func PrintResult(proto *shared_proto.CombatProto) {

	fmt.Println("攻方胜利：", proto.AttackerWin, "轮次：", len(proto.Rounds))

	player := make(map[int32]int32)
	fmt.Println("攻方信息：")
	for i, t := range proto.Attacker.Troops {
		idx := proto.AttackerTroopPos[i].Index
		fmt.Println(idx, t.Captain.Name, t.Captain.Soldier*t.Captain.LifePerSoldier)
		player[idx] = t.Captain.Soldier * t.Captain.LifePerSoldier
	}

	fmt.Println("防方信息：")
	for i, t := range proto.Defenser.Troops {
		idx := proto.DefenserTroopPos[i].Index
		fmt.Println(idx, t.Captain.Name, t.Captain.Soldier*t.Captain.LifePerSoldier)
		player[idx] = t.Captain.Soldier * t.Captain.LifePerSoldier
	}

	for i, round := range proto.Rounds {
		fmt.Println("轮次：", i)

		for _, act := range round.Actions {
			if act.MoveDirection != shared_proto.Direction_ORIGIN {
				fmt.Println(act.Index, "移动", act.MoveDirection)
			}

			if act.TargetIndex > 0 {
				player[act.TargetIndex] = player[act.TargetIndex] - act.Damage
				fmt.Println(act.Index, "打", act.TargetIndex, act.HurtType, act.Damage, "剩余血量：", player[act.TargetIndex])
			}
		}
	}

	fmt.Println("存活士兵：")
	for _, s := range proto.AliveSolider {
		fmt.Println(s.Key, s.Value)
	}
}

func PrintMultiResult(result *shared_proto.MultiCombatProto) {

	leftSoldierLifeMap := make(map[int32]int32)
	leftSoldierMap := make(map[int32]int32)
	attackerIdMap := make(map[int64]int32)
	defenserIdMap := make(map[int64]int32)

	for _, c := range result.Combats {
		id, _ := idbytes.ToId(c.Attacker.Hero.Id)
		if _, exist := attackerIdMap[id]; !exist {
			attackerIdMap[id] = 1
			for i, t := range c.Attacker.Troops {
				idx := c.AttackerTroopPos[i].Index
				leftSoldierLifeMap[idx] = t.Captain.Soldier * t.Captain.LifePerSoldier
				leftSoldierMap[idx] = t.Captain.Soldier
			}
		}

		id, _ = idbytes.ToId(c.Defenser.Hero.Id)
		if _, exist := defenserIdMap[id]; !exist {
			defenserIdMap[id] = 1
			for i, t := range c.Defenser.Troops {
				idx := c.DefenserTroopPos[i].Index
				leftSoldierLifeMap[idx] = t.Captain.Soldier * t.Captain.LifePerSoldier
				leftSoldierMap[idx] = t.Captain.Soldier
			}
		}
	}

	for _, f := range result.GetCombats() {
		PrintFight(leftSoldierLifeMap, leftSoldierMap, result, f)

		attackerId, _ := idbytes.ToId(f.Attacker.Hero.Id)
		defenserId, _ := idbytes.ToId(f.Defenser.Hero.Id)

		if f.AttackerWin {
			delete(defenserIdMap, defenserId)
			if f.IsWinnerContinueWinLeave {
				delete(attackerIdMap, attackerId)
			}
		} else {
			delete(attackerIdMap, attackerId)
			if f.IsWinnerContinueWinLeave {
				delete(defenserIdMap, defenserId)
			}
		}
	}

	fmt.Println("---------------------------")
	fmt.Println("最终胜利方是进攻方: ", result.AttackerWin)
	fmt.Println("进攻方剩余部队数: ", len(attackerIdMap))
	fmt.Println("防守方剩余部队数: ", len(defenserIdMap))
}

func PrintFight(leftSoldierLifeMap, leftSoldierMap map[int32]int32, combat *shared_proto.MultiCombatProto, proto *shared_proto.CombatProto) {
	fmt.Println("------------------------------------------")
	fmt.Println("攻方id:", proto.Attacker.Hero.Id)
	fmt.Println("防方id:", proto.Defenser.Hero.Id)
	fmt.Println("攻方胜利：", proto.AttackerWin)
	fmt.Println("轮次：", len(proto.Rounds))
	fmt.Println("战斗所在队列", proto.GetFightQueue())

	fmt.Println("攻方信息：")

	for _, attacker := range combat.Attacker {
		if bytes.Equal(attacker.Hero.Id, proto.Attacker.Hero.Id) {
			for i, t := range attacker.Troops {
				idx := proto.AttackerTroopPos[i].Index
				fmt.Println(idx, t.Captain.Name, leftSoldierMap[idx], leftSoldierLifeMap[idx])
			}
		}
	}

	fmt.Println("防方信息：")
	for _, defenser := range combat.Defenser {
		if bytes.Equal(defenser.Hero.Id, proto.Defenser.Hero.Id) {
			for i, t := range defenser.Troops {
				idx := proto.DefenserTroopPos[i].Index
				fmt.Println(idx, t.Captain.Name, leftSoldierMap[idx], leftSoldierLifeMap[idx])
			}
		}
	}

	for i, round := range proto.Rounds {
		fmt.Println("轮次：", i+1)

		for _, act := range round.Actions {
			if act.MoveDirection != shared_proto.Direction_ORIGIN {
				fmt.Println("\t", act.Index, "移动", act.MoveDirection)
			}

			if act.TargetIndex > 0 {
				oldLife := leftSoldierLifeMap[act.TargetIndex]
				leftSoldierLifeMap[act.TargetIndex] = oldLife - act.Damage
				fmt.Println("\t", act.Index, "打", act.TargetIndex, act.HurtType, act.Damage, "此前血量：", oldLife, "剩余血量：", leftSoldierLifeMap[act.TargetIndex])
			}
		}
	}

	if proto.AttackerWin {
		fmt.Printf("进攻方 [%s] 存活士兵: \n", proto.Attacker.Hero.Id)
	} else {
		fmt.Printf("防守方 [%s] 存活士兵: \n", proto.Defenser.Hero.Id)
	}

	for _, s := range proto.AliveSolider {
		fmt.Println(s.Key, s.Value, leftSoldierLifeMap[s.Key])

		leftSoldierMap[s.Key] = s.Value
	}

	fmt.Println("胜利方连胜次数: ", proto.WinnerContinueWinTimes, " 胜利方是否连胜离场: ", proto.IsWinnerContinueWinLeave)
}
