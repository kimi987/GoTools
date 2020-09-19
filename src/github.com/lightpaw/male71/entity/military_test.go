package entity

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

//func TestTraining_ResetTrainTimeWhenOutputChanged(t *testing.T) {
//	RegisterTestingT(t)
//
//	var training *Training = NewTraining(10, &military_data.TrainingLevelData{Level: 1, Name: "NAME", Desc: "DESC", Coef: 1, Cost: nil},
//		10, &building_effect{trainOutput: 10, trainCapcity: 40})
//
//	ctime := time.Now()
//
//	// 操作失败，时间为1970xxxx
//	var canCollectExp uint64 = 0
//	training.ResetTrainTimeWhenOutputChanged(canCollectExp, ctime)
//	Ω(timeutil.IsZero(training.startTime)).Should(Equal(true))
//
//	// 操作成功，时间为当前时间-30/10*小时
//	canCollectExp = 30
//	training.ResetTrainTimeWhenOutputChanged(canCollectExp, ctime)
//	Ω(training.startTime).Should(Equal(ctime.Add(-time.Hour * 3)))
//
//	Ω(training.GetExp(ctime)).Should(Equal(canCollectExp))
//
//	// 操作成功，时间为当前时间-30/20*小时
//	training.buildingEffect.trainOutput = 20 // 每小时产出20
//	training.ResetTrainTimeWhenOutputChanged(canCollectExp, ctime)
//	Ω(training.startTime).Should(Equal(ctime.Add(-time.Hour * 30 / 20)))
//
//	// 操作成功，时间为当前时间-30/25*小时
//	training.buildingEffect.trainOutput = 25 // 每小时产出25
//	training.ResetTrainTimeWhenOutputChanged(canCollectExp, ctime)
//	Ω(training.startTime).Should(Equal(ctime.Add(-time.Hour * 30 / 25)))
//
//	// 操作成功，时间为当前时间-30/4*小时
//	training.buildingEffect.trainOutput = 4 // 每小时产出4
//	training.ResetTrainTimeWhenOutputChanged(canCollectExp, ctime)
//	Ω(training.startTime).Should(Equal(ctime.Add(-time.Hour * 30 / 4)))
//
//	canCollectExp = 300
//	// capcity最多为40，时间为当前时间-40/10*小时
//	training.buildingEffect.trainOutput = 10 // 每小时产出10
//	training.ResetTrainTimeWhenOutputChanged(canCollectExp, ctime)
//	Ω(training.startTime).Should(Equal(ctime.Add(-time.Hour * 4)))
//
//	// capcity最多为40，经验要相等
//	var capcity uint64 = 40
//	Ω(training.GetExp(ctime)).Should(Equal(capcity))
//}

func TestTroopId(t *testing.T) {
	RegisterTestingT(t)

	tid := NewHeroTroopId(0, 0)
	Ω(GetTroopHeroId(tid)).Should(Equal(int64(0)))
	Ω(GetTroopIndex(tid)).Should(Equal(uint64(0)))

	tid = NewHeroTroopId(101, 8)
	Ω(GetTroopHeroId(tid)).Should(Equal(int64(101)))
	Ω(GetTroopIndex(tid)).Should(Equal(uint64(8)))

	tid = NewHeroTroopId(1<<56-1, 15)
	Ω(GetTroopHeroId(tid)).Should(Equal(int64(1<<56 - 1)))
	Ω(GetTroopIndex(tid)).Should(Equal(uint64(15)))
}

func TestSwapTroopsCaptain(t *testing.T) {
	//RegisterTestingT(t)
	//
	//t1 := &Troop{}
	//t1.captains = make([]*Captain, 5)
	//c1 := &Captain{totalStat: &data.SpriteStat{}}
	//c1.troop = t1
	//t1.captains[3] = c1
	//
	//t2 := &Troop{}
	//t2.captains = make([]*Captain, 5)
	//c2 := &Captain{totalStat: &data.SpriteStat{}}
	//c2.troop = t2
	//t2.captains[3] = c2
	//
	//SwapTroopsCaptain(t1, t2, 0, 0)
	//Ω(t1.captains).Should(Equal([]*Captain{nil, nil, nil, c1, nil}))
	//Ω(t2.captains).Should(Equal([]*Captain{nil, nil, nil, c2, nil}))
	//Ω(c1.troop).Should(Equal(t1))
	//Ω(c2.troop).Should(Equal(t2))
	//
	//SwapTroopsCaptain(t1, t2, 3, 0)
	//Ω(t1.captains).Should(Equal([]*Captain{nil, nil, nil, nil, nil}))
	//Ω(t2.captains).Should(Equal([]*Captain{c1, nil, nil, c2, nil}))
	//fmt.Println(c1.troop)
	//fmt.Println(c2.troop)
	//Ω(c1.troop).Should(Equal(t2))
	//Ω(c2.troop).Should(Equal(t2))
	//
	//SwapTroopsCaptain(t1, t2, 3, 0)
	//Ω(t1.captains).Should(Equal([]*Captain{nil, nil, nil, c1, nil}))
	//Ω(t2.captains).Should(Equal([]*Captain{nil, nil, nil, c2, nil}))
	//Ω(c1.troop).Should(Equal(t1))
	//Ω(c2.troop).Should(Equal(t2))
	//
	//SwapTroopsCaptain(t1, t2, 2, 3)
	//Ω(t1.captains).Should(Equal([]*Captain{nil, nil, c2, c1, nil}))
	//Ω(t2.captains).Should(Equal([]*Captain{nil, nil, nil, nil, nil}))
	//Ω(c1.troop).Should(Equal(t1))
	//Ω(c2.troop).Should(Equal(t1))
	//
	//SwapTroopsCaptain(t1, t1, 2, 3)
	//Ω(t1.captains).Should(Equal([]*Captain{nil, nil, c1, c2, nil}))
	//Ω(t2.captains).Should(Equal([]*Captain{nil, nil, nil, nil, nil}))
	//Ω(c1.troop).Should(Equal(t1))
	//Ω(c2.troop).Should(Equal(t1))
	//
	//SwapTroopsCaptain(t1, t1, 1, 3)
	//Ω(t1.captains).Should(Equal([]*Captain{nil, c2, c1, nil, nil}))
	//Ω(t2.captains).Should(Equal([]*Captain{nil, nil, nil, nil, nil}))
	//Ω(c1.troop).Should(Equal(t1))
	//Ω(c2.troop).Should(Equal(t1))

}

func TestSoldier(t *testing.T) {
	RegisterTestingT(t)

	ctime := time.Now()

	m := &Military{
		freeSoldier:         NewRecoverableTimes(ctime, time.Second, 100),
		overflowFreeSoldier: 0,
		buildingEffect:newBuildingEffect(),
	}

	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(0)))

	m.AddFreeSoldier(10, ctime)
	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(10)))

	m.ReduceFreeSoldier(5, ctime)
	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(5)))

	m.ReduceFreeSoldier(100, ctime)
	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(0)))

	m.AddFreeSoldier(111, ctime)
	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(111)))

	m.ReduceFreeSoldier(5, ctime)
	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(106)))

	m.ReduceFreeSoldier(10, ctime)
	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(96)))

	m.AddFreeSoldier(10, ctime)
	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(106)))
	Ω(m.overflowFreeSoldier).Should(Equal(uint64(6)))

	// buildingEffect.soldierCapcity == 0
	m.freeSoldier.ChangeMaxTimes(200, ctime)
	m.SetFreeSoldier(106, ctime)
	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(106)))
	Ω(m.overflowFreeSoldier).Should(Equal(uint64(0)))

	m.SetFreeSoldier(206, ctime)
	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(206)))
	Ω(m.overflowFreeSoldier).Should(Equal(uint64(6)))

	m.SetFreeSoldier(106, ctime)
	Ω(m.FreeSoldier(ctime)).Should(Equal(uint64(106)))
	Ω(m.overflowFreeSoldier).Should(Equal(uint64(0)))
}
