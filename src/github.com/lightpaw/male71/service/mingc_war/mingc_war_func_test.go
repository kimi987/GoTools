package mingc_war

import (
	"fmt"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/mingc"
	"github.com/lightpaw/male7/util/timeutil"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestMingcWarService(t *testing.T) {
	RegisterTestingT(t)
	dep := mock.NewMockDep2()
	dep.Mock(dep.Mingc, func() iface.MingcService {
		return mingc.NewMingcService(dep.Datas(), dep.Db(), ifacemock.TickerService, dep.Guild(), dep.GuildSnapshot(), dep.Country(), dep.Time(), dep.SvrConf(), dep.Mail(), dep.HeroData(), dep.Broadcast(), dep.World(), dep.Tlog())
	})
	dep.Mock(dep.Time, func() iface.TimeService {
		ts := ifacemock.TimeService
		ts.Mock(ts.CurrentTime, func() time.Time {
			return time.Now()
		})
		return ts
	})

	// test
	ctime1 := time.Date(2018, 8, 22, 22, 9, 0, 27, timeutil.East8)

	var timeData *mingcdata.MingcTimeData
	var start time.Time
	for _, d := range dep.Datas().GetMingcTimeDataArray() {
		nextStart, _ := d.WarTime.NextTime(ctime1)
		if timeutil.IsZero(start) || nextStart.Before(start) {
			start, timeData = nextStart, d
		}
	}

	w := newMingcWar(ctime1, dep, timeData)
	fmt.Println(ctime1)
	fmt.Println(w.startTime)
	fmt.Println(w.endTime)
	fmt.Println(w.stage(shared_proto.MingcWarState_MC_T_FIGHT).stageEndTime())

	// mock time

	ctime := dep.Time().CurrentTime()
	war := mockMingcWar(dep, ctime)
	war.changeState(ctime)
	Ω(war.currStage().warState()).Should(Equal(shared_proto.MingcWarState_MC_T_APPLY_ATK))

	var mcId uint64 = 1
	mc := war.mingcs[mcId]

	// applyAtk
	cost := dep.Datas().GetMingcBaseData(mcId).AtkMinHufu
	errMsg := applyAttack(war, 1, dep.Datas().GetMingcBaseData(mcId), cost, ctime)
	Ω(errMsg).Should(BeNil())
	errMsg = applyAttack(war, 2, dep.Datas().GetMingcBaseData(mcId), cost+100, ctime)
	Ω(errMsg).Should(BeNil())

	// applyAst
	dep.Mock(dep.Time, func() iface.TimeService {
		ts := ifacemock.TimeService
		ts.Mock(ts.CurrentTime, func() time.Time {
			return time.Now().Add(time.Hour)
		})
		return ts
	})
	ctime = dep.Time().CurrentTime()
	_, oldState := war.changeState(ctime)
	Ω(war.currStage().warState()).Should(Equal(shared_proto.MingcWarState_MC_T_APPLY_AST))
	war.stage(oldState).onEnd(war, dep)
	war.currStage().onStart(war, dep)
	Ω(mc.atkId).Should(Equal(int64(2)))

	recvGid, errMsg := applyAssist(war, 3, true, dep.Datas().GetMingcBaseData(mcId), dep.Datas().MingcMiscData())
	Ω(errMsg).Should(BeNil())
	Ω(recvGid).Should(Equal(int64(2)))
	recvGid, errMsg = applyAssist(war, 4, true, dep.Datas().GetMingcBaseData(mcId), dep.Datas().MingcMiscData())
	Ω(errMsg).Should(BeNil())
	Ω(recvGid).Should(Equal(int64(2)))

	recvGid, errMsg = cancelApplyAst(war, 4, dep.Datas().GetMingcBaseData(mcId))
	Ω(errMsg).Should(BeNil())
	Ω(recvGid).Should(Equal(int64(2)))
	_, ok := war.mingcs[mcId].applyAsts[4]
	Ω(ok).Should(BeFalse())
	mcs, ok := war.currStage().(*ApplyAstObj).ApplyAstGuilds[4]
	Ω(len(mcs)).Should(Equal(0))

	errMsg = replyApplyAst(war, 2, 3, dep.Datas().GetMingcBaseData(mcId), true, dep.Datas().MingcMiscData())
	Ω(errMsg).Should(BeNil())
	Ω(len(mc.astAtkList)).Should(Equal(1))

	actionTime := time.Now()
	// fight
	dep.Mock(dep.Time, func() iface.TimeService {
		ts := ifacemock.TimeService
		ts.Mock(ts.CurrentTime, func() time.Time {
			actionTime = actionTime.Add(2 * time.Hour)
			return actionTime
		})
		return ts
	})
	ctime = dep.Time().CurrentTime()
	_, oldState = war.changeState(ctime)
	Ω(war.currStage().warState()).Should(Equal(shared_proto.MingcWarState_MC_T_FIGHT))
	war.stage(oldState).onEnd(war, dep)
	war.currStage().onStart(war, dep)

	_, errMsg = joinFight(war, mcId, 2, 2, &shared_proto.HeroBasicProto{}, []*shared_proto.CaptainInfoProto{}, mc.scene.atkReliveBuilding.pos, ctime)
	Ω(errMsg).Should(BeNil())
	_, errMsg = joinFight(war, mcId, 3, 3, &shared_proto.HeroBasicProto{}, []*shared_proto.CaptainInfoProto{}, mc.scene.atkReliveBuilding.pos, ctime)
	Ω(errMsg).Should(BeNil())
	Ω(len(mc.scene.atkTroops)).Should(Equal(2))

	_, errMsg = quitFight(war, 3)
	Ω(errMsg).Should(BeNil())
	Ω(len(mc.scene.atkTroops)).Should(Equal(1))

	// change mode
	mc.scene.onUpdate(ctime)
	_, errMsg = sceneChangeMode(war, 2, shared_proto.MingcWarModeType_MC_MT_FREE_TANK)
	Ω(errMsg).Should(BeNil())
	Ω(mc.scene.atkTroops[2].mode).Should(Equal(shared_proto.MingcWarModeType_MC_MT_FREE_TANK))
	Ω(mc.scene.atkTroops[2].getDestroyProsperity()).Should(Equal(dep.Datas().MingcMiscData().FreeTankPerDestroyProsperity))
	Ω(mc.scene.atkTroops[2].getSpeed()).Should(Equal(dep.Datas().MingcMiscData().FreeTankSpeed))
	mc.scene.onUpdate(ctime)
	_, errMsg = sceneChangeMode(war, 2, shared_proto.MingcWarModeType_MC_MT_NORMAL)
	Ω(errMsg).Should(BeNil())
	Ω(mc.scene.atkTroops[2].mode).Should(Equal(shared_proto.MingcWarModeType_MC_MT_NORMAL))
	Ω(mc.scene.atkTroops[2].getDestroyProsperity()).Should(Equal(dep.Datas().MingcMiscData().PerDestroyProsperity))
	Ω(mc.scene.atkTroops[2].getSpeed()).Should(Equal(dep.Datas().MingcMiscData().Speed))

	// move. from relive to home
	dep.Mock(dep.Time, func() iface.TimeService {
		ts := ifacemock.TimeService
		ts.Mock(ts.CurrentTime, func() time.Time {
			actionTime = actionTime.Add(dep.Datas().MingcMiscData().FightPrepareDuration).Add(dep.Datas().MingcMiscData().JoinFightDuration)
			return actionTime
		})
		return ts
	})
	ctime = dep.Time().CurrentTime()
	mc.scene.onUpdate(ctime)
	_, errMsg = sceneMove(war, 2, mc.scene.atkHomeBuilding.pos, ctime)
	Ω(errMsg).Should(BeNil())
	Ω(mc.scene.atkTroops[2].action.getState()).Should(Equal(shared_proto.MingcWarTroopState_MC_TP_MOVING))

	dep.Mock(dep.Time, func() iface.TimeService {
		ts := ifacemock.TimeService
		ts.Mock(ts.CurrentTime, func() time.Time {
			actionTime = actionTime.Add(time.Minute)
			return actionTime
		})
		return ts
	})
	ctime = dep.Time().CurrentTime()
	mc.scene.onUpdate(ctime)
	Ω(mc.scene.atkTroops[2].action.getPos()).Should(Equal(mc.scene.atkHomeBuilding.pos))
	Ω(mc.scene.atkTroops[2].action.getState()).Should(Equal(shared_proto.MingcWarTroopState_MC_TP_STATION))

	// moving. from home to relive
	dep.Mock(dep.Time, func() iface.TimeService {
		ts := ifacemock.TimeService
		ts.Mock(ts.CurrentTime, func() time.Time {
			actionTime.Add(time.Minute)
			return actionTime
		})
		return ts
	})
	ctime = dep.Time().CurrentTime()
	mc.scene.onUpdate(ctime)
	_, errMsg = sceneMove(war, 2, mc.scene.atkReliveBuilding.pos, ctime)
	Ω(errMsg).Should(BeNil())
	Ω(mc.scene.atkTroops[2].action.getState()).Should(Equal(shared_proto.MingcWarTroopState_MC_TP_MOVING))

	dep.Mock(dep.Time, func() iface.TimeService {
		ts := ifacemock.TimeService
		ts.Mock(ts.CurrentTime, func() time.Time {
			actionTime = actionTime.Add(time.Second)
			return actionTime
		})
		return ts
	})
	ctime = dep.Time().CurrentTime()
	mc.scene.onUpdate(ctime)
	Ω(mc.scene.atkTroops[2].action.getPos()).Should(Equal(mc.scene.atkHomeBuilding.pos))
	Ω(mc.scene.atkTroops[2].action.getState()).Should(Equal(shared_proto.MingcWarTroopState_MC_TP_MOVING))

	// back.back to home
	_, errMsg = sceneBack(war, 2, ctime)
	Ω(errMsg).Should(BeNil())
	Ω(mc.scene.atkTroops[2].action.getPos()).Should(Equal(mc.scene.atkReliveBuilding.pos))
	Ω(mc.scene.atkTroops[2].action.(*MoveAction).destPos).Should(Equal(mc.scene.atkHomeBuilding.pos))
	Ω(mc.scene.atkTroops[2].action.getState()).Should(Equal(shared_proto.MingcWarTroopState_MC_TP_MOVING))

	dep.Mock(dep.Time, func() iface.TimeService {
		ts := ifacemock.TimeService
		ts.Mock(ts.CurrentTime, func() time.Time {
			actionTime = actionTime.Add(time.Minute)
			return actionTime
		})
		return ts
	})
	ctime = dep.Time().CurrentTime()
	mc.scene.onUpdate(ctime)
	Ω(errMsg).Should(BeNil())
	Ω(mc.scene.atkTroops[2].action.getState()).Should(Equal(shared_proto.MingcWarTroopState_MC_TP_STATION))
	Ω(mc.scene.atkTroops[2].action.getPos()).Should(Equal(mc.scene.atkHomeBuilding.pos))
}

func mockMingcWar(dep iface.ServiceDep, ctime time.Time) *MingcWar {
	timeData := &mingcdata.MingcTimeData{}
	t := ctime.Add((-7 * 24) * time.Hour)
	tdur := t.Hour() * int(time.Hour)
	timeData.ApplyAtkTime = &timeutil.WeekDurTime{WeekTime: &timeutil.WeekTime{Week: timeutil.ConvCnWeekday(ctime.Weekday()), Time: time.Duration(tdur)}, Dur: time.Hour}
	timeData.ApplyAstTime = &timeutil.WeekDurTime{WeekTime: &timeutil.WeekTime{Week: timeutil.ConvCnWeekday(ctime.Weekday()), Time: time.Duration(tdur + int(time.Hour))}, Dur: time.Hour}
	timeData.FightTime = &timeutil.WeekDurTime{WeekTime: &timeutil.WeekTime{Week: timeutil.ConvCnWeekday(ctime.Weekday()), Time: time.Duration(tdur + 2*int(time.Hour))}, Dur: time.Hour}
	timeData.WarTime = &timeutil.WeekDurTime{WeekTime: &timeutil.WeekTime{Week: timeutil.ConvCnWeekday(ctime.Weekday()), Time: time.Duration(tdur)}, Dur: 3 * time.Hour}
	war := newMingcWar(t, dep, timeData)
	return war
}
