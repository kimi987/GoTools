package zhengwu

import (
	"github.com/lightpaw/logrus"
	zhengwu_data "github.com/lightpaw/male7/config/zhengwu"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/zhengwu"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/operate_type"
)

func NewZhengWuModule(dep iface.ServiceDep) *ZhengWuModule {
	m := &ZhengWuModule{}

	m.dep = dep
	m.datas = dep.Datas()
	m.broadcast = dep.Broadcast()
	m.timeService = dep.Time()

	return m

}

//gogen:iface
type ZhengWuModule struct {
	dep         iface.ServiceDep
	datas       iface.ConfigDatas
	timeService iface.TimeService
	broadcast   iface.BroadcastService
}

//gogen:iface
func (m *ZhengWuModule) ProcessStart(proto *zhengwu.C2SStartProto, hc iface.HeroController) {
	data := m.datas.GetZhengWuData(uint64(proto.Id))
	if data == nil {
		logrus.Debugf("客户端发送上来要开始的政务没找到: %d", proto.Id)
		hc.Send(zhengwu.ERR_START_FAIL_NOT_FOUND)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		ctime := m.timeService.CurrentTime()

		heroZhengWu := hero.ZhengWu()
		doingZhengWu := heroZhengWu.Doing()
		if doingZhengWu != nil {
			if doingZhengWu.IsCompleted(ctime) {
				logrus.Debugf("当前正在做的政务已经完成了")
				result.Add(zhengwu.ERR_START_FAIL_NEED_COLLECT)
				return
			}

			logrus.Debugf("当前有正在做的政务")
			result.Add(zhengwu.ERR_START_FAIL_HAVE_DOING)
			return
		}

		if !heroZhengWu.RemoveToDo(data) {
			logrus.Debugf("要开始的政务不在todo列表里面")
			result.Add(zhengwu.ERR_START_FAIL_NOT_FOUND)
			return
		}

		result.Changed()
		result.Ok()

		newDoing := entity.NewZhengWu(data, ctime.Add(data.Duration))
		heroZhengWu.SetDoing(newDoing)

		result.Add(zhengwu.NewS2cStartMsg(must.Marshal(newDoing.EncodeClient())))
	})
}

//gogen:iface
func (m *ZhengWuModule) ProcessVipCollect(proto *zhengwu.C2SVipCollectProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if vipData := m.dep.Datas().GetVipLevelData(hero.VipLevel()); vipData == nil || !vipData.ZhengWuAutoCompleted {
			result.Add(zhengwu.ERR_VIP_COLLECT_FAIL_VIP_LIMIT)
			return
		}

		id := u64.FromInt32(proto.Id)
		ctime := m.dep.Time().CurrentTime()

		heroZhengWu := hero.ZhengWu()
		doingZhengWu := heroZhengWu.Doing()

		var data *zhengwu_data.ZhengWuData
		if doingZhengWu != nil && doingZhengWu.Data().Id == id {
			data = doingZhengWu.Data()
			doingZhengWu.Complete(ctime)
			heroZhengWu.Complete()
		} else {
			for _, d := range heroZhengWu.ToDoList() {
				if d.Id == id {
					data = d
					heroZhengWu.RemoveToDo(d)
					break
				}
			}
		}
		if data == nil {
			result.Add(zhengwu.ERR_VIP_COLLECT_FAIL_NOT_IN_LIST)
			return
		}

		hctx := heromodule.NewContext(m.dep, operate_type.ZhengWuCollect)
		heromodule.AddPrize(hctx, hero, result, data.Prize, ctime)

		result.Add(zhengwu.NewS2cVipCollectMsg(proto.Id))

		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_CompleteZhengWu)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_COMPLETE_ZHENG_WU)

		result.Changed()
		result.Ok()
	})
}

//gogen:iface c2s_collect
func (m *ZhengWuModule) ProcessCollect(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroZhengWu := hero.ZhengWu()
		doingZhengWu := heroZhengWu.Doing()
		if doingZhengWu == nil {
			logrus.Debugf("当前没有正在做的政务")
			result.Add(zhengwu.ERR_COLLECT_FAIL_NO_DOING_ZHENG_WU)
			return
		}

		ctime := m.timeService.CurrentTime()
		if !doingZhengWu.IsCompleted(ctime) {
			logrus.Debugf("当前政务没有完成")
			result.Add(zhengwu.ERR_COLLECT_FAIL_NOT_COMPLETE)
			return
		}

		heroZhengWu.Complete()

		hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_COMPELTED_FIRST_ZHENGWU)

		data := doingZhengWu.Data()

		hctx := heromodule.NewContext(m.dep, operate_type.ZhengWuCollect)
		heromodule.AddPrize(hctx, hero, result, data.Prize, ctime)

		result.Add(zhengwu.COLLECT_S2C)

		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_CompleteZhengWu)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_COMPLETE_ZHENG_WU)

		result.Changed()
		result.Ok()
	})
}

//gogen:iface c2s_yuanbao_complete
// 元宝完成
func (m *ZhengWuModule) ProcessYuanBaoComplete(hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.ZhengWuYuanBaoComplete)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroZhengWu := hero.ZhengWu()
		doingZhengWu := heroZhengWu.Doing()
		if doingZhengWu == nil {
			logrus.Debugf("当前没有正在做的政务")
			result.Add(zhengwu.ERR_YUANBAO_COMPLETE_FAIL_NOT_DOING)
			return
		}

		ctime := m.timeService.CurrentTime()
		if doingZhengWu.IsCompleted(ctime) {
			logrus.Debugf("当前政务已经完成")
			result.Add(zhengwu.ERR_YUANBAO_COMPLETE_FAIL_COMPLETE)
			return
		}

		leftTime := doingZhengWu.EndTime().Sub(ctime)
		if leftTime <= time.Second {
			logrus.Debugf("当前政务快完成")
			result.Add(zhengwu.ERR_YUANBAO_COMPLETE_FAIL_COMPLETE)
			return
		}

		cost := uint64(float64(leftTime) / float64(doingZhengWu.Data().Duration) * float64(doingZhengWu.Data().Cost.Cost))
		cost = u64.Max(cost, 1)

		if !heromodule.ReduceDianquan(hctx, hero, result, cost) {
			logrus.Debugf("没有足够的点券")
			result.Add(zhengwu.ERR_YUANBAO_COMPLETE_FAIL_NOT_ENOUGH_YUANBAO)
			return
		}

		result.Changed()
		result.Ok()

		doingZhengWu.Complete(ctime)

		result.Add(zhengwu.YUANBAO_COMPLETE_S2C)
	})
}

//gogen:iface c2s_yuanbao_refresh
// 元宝刷新
func (m *ZhengWuModule) ProcessYuanBaoRefresh(hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.ZhengWuYuanBaoRefresh)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroZhengWu := hero.ZhengWu()

		ctime := m.timeService.CurrentTime()

		doingZhengWu := heroZhengWu.Doing()
		if doingZhengWu != nil && doingZhengWu.IsCompleted(ctime) {
			logrus.Debugf("当前有政务已经完成，请先领取奖励再刷新")
			result.Add(zhengwu.ERR_YUANBAO_REFRESH_FAIL_COMPLETE)
			return
		}

		var refreshData *zhengwu_data.ZhengWuRefreshData

		refreshTimes := heroZhengWu.RefreshTimes()
		for _, checkRefreshData := range m.datas.GetZhengWuRefreshDataArray() {
			if refreshTimes < checkRefreshData.Times {
				refreshData = checkRefreshData
				break
			}
		}

		if refreshData == nil {
			logrus.Debugf("次数不够，无法刷新")
			result.Add(zhengwu.ERR_YUANBAO_REFRESH_FAIL_NOT_ENOUGH_TIMES)
			return
		}

		if refreshData.NewCost != nil {
			if !heromodule.TryReduceCost(hctx, hero, result, refreshData.NewCost) {
				logrus.Debugf("没有足够的点券")
				result.Add(zhengwu.ERR_YUANBAO_REFRESH_FAIL_NOT_ENOUGH_COST)
				return
			}
		}

		result.Changed()
		result.Ok()

		heroZhengWu.IncRefreshTimes()

		randomCount := m.datas.ZhengWuMiscData().RandomCount

		completedFirstZhengWu := hero.Bools().Get(shared_proto.HeroBoolType_BOOL_COMPELTED_FIRST_ZHENGWU)
		heroZhengWu.SetToDoList(m.datas.ZhengWuRandomData().Random(randomCount, completedFirstZhengWu))

		// 发送新的政务列表
		result.Add(zhengwu.NewS2cRefreshMsg(must.Marshal(heroZhengWu.EncodeClient())))

		result.Add(zhengwu.YUANBAO_REFRESH_S2C)
	})
}
