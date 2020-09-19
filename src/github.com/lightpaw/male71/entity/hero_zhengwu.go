package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/zhengwu"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func NewHeroZhengWu(nextRefreshTime time.Time, defaultZhengWu []*zhengwu.ZhengWuData) *HeroZhengWu {
	return &HeroZhengWu{nextRefreshTime: nextRefreshTime, toDoList: defaultZhengWu}
}

// 英雄的所有政务
type HeroZhengWu struct {
	refreshTimes    uint64                 // 刷新次数
	nextRefreshTime time.Time              // 上次刷新时间
	doing           *ZhengWu               // 正在做的政务
	toDoList        []*zhengwu.ZhengWuData // 要做的列表
}

// 获得上次刷新时间
func (h *HeroZhengWu) NextRefreshTime() time.Time {
	return h.nextRefreshTime
}

// 设置上次刷新时间
func (h *HeroZhengWu) SetNextRefreshTime(toSet time.Time) {
	h.nextRefreshTime = toSet
}

func (h *HeroZhengWu) Doing() *ZhengWu {
	return h.doing
}

func (h *HeroZhengWu) SetDoing(toSet *ZhengWu) {
	h.doing = toSet
}

func (h *HeroZhengWu) Complete() {
	h.doing = nil
}

func (h *HeroZhengWu) RefreshTimes() uint64 {
	return h.refreshTimes
}

func (h *HeroZhengWu) IncRefreshTimes() uint64 {
	h.refreshTimes++
	return h.refreshTimes
}

func (h *HeroZhengWu) ToDoList() []*zhengwu.ZhengWuData {
	return h.toDoList
}

func (h *HeroZhengWu) SetToDoList(toSet []*zhengwu.ZhengWuData) {
	h.toDoList = toSet
}

func (h *HeroZhengWu) RemoveToDo(data *zhengwu.ZhengWuData) bool {
	for idx, toDoData := range h.toDoList {
		if toDoData != data {
			// 移除掉
			continue
		}

		if idx+1 != len(h.toDoList) {
			// 交换
			h.toDoList[idx] = h.toDoList[len(h.toDoList)-1]
		}

		h.toDoList = h.toDoList[:len(h.toDoList)-1]

		return true
	}

	return false
}

func (h *HeroZhengWu) ForceCompleteAllZhengWu() {

}

func (h *HeroZhengWu) ResetDaily() {
	h.refreshTimes = 0
}

func (h *HeroZhengWu) EncodeClient() *shared_proto.HeroZhengWuProto {
	proto := &shared_proto.HeroZhengWuProto{}

	proto.RefreshTimes = u64.Int32(h.refreshTimes)

	proto.NextRefreshTime = timeutil.Marshal32(h.nextRefreshTime)

	if h.doing != nil {
		proto.Doing = h.doing.EncodeClient()
	}

	proto.ToDoList = make([]*shared_proto.ZhengWuDataProto, 0, len(h.toDoList))
	for _, todo := range h.toDoList {
		proto.ToDoList = append(proto.ToDoList, todo.Proto)
	}

	return proto
}

func (h *HeroZhengWu) EncodeServer() *server_proto.HeroZhengWuServerProto {
	proto := &server_proto.HeroZhengWuServerProto{}

	proto.RefreshTimes = h.refreshTimes

	proto.NextRefreshTime = timeutil.Marshal64(h.nextRefreshTime)

	if h.doing != nil {
		proto.Doing = h.doing.EncodeServer()
	}

	proto.ToDoList = make([]uint64, 0, len(h.toDoList))
	for _, todo := range h.toDoList {
		proto.ToDoList = append(proto.ToDoList, todo.Id)
	}

	return proto
}

func (h *HeroZhengWu) unmarshal(proto *server_proto.HeroZhengWuServerProto, datas *config.ConfigDatas, ctime time.Time) {
	if proto == nil {
		return
	}

	h.refreshTimes = proto.RefreshTimes

	h.nextRefreshTime = timeutil.Unix64(proto.NextRefreshTime)

	if proto.Doing != nil {
		doingData := datas.GetZhengWuData(proto.Doing.Data)
		if doingData != nil {
			h.doing = NewZhengWu(doingData, timeutil.Unix64(proto.Doing.EndTime))
		} else {
			logrus.Errorf("玩家正在做的政务不见了: %#v", proto.Doing)
		}
	}

	h.toDoList = make([]*zhengwu.ZhengWuData, 0, len(proto.ToDoList))
	for _, todoId := range proto.ToDoList {
		data := datas.GetZhengWuData(todoId)
		if data == nil {
			logrus.Errorf("玩家的政务列表中有配置找不到了: %v", proto.ToDoList)
		} else {
			h.toDoList = append(h.toDoList, data)
		}
	}
}

func NewZhengWu(
	data *zhengwu.ZhengWuData, // 政务
	endTime time.Time, // 开始时间
) *ZhengWu {
	return &ZhengWu{
		data:    data,
		endTime: endTime,
	}
}

// 政务
type ZhengWu struct {
	data    *zhengwu.ZhengWuData // 政务
	endTime time.Time            // 结束时间
}

func (zw *ZhengWu) Data() *zhengwu.ZhengWuData {
	return zw.data
}

func (zw *ZhengWu) StartTime() time.Time {
	return zw.endTime.Add(-zw.data.Duration)
}

func (zw *ZhengWu) EndTime() time.Time {
	return zw.endTime
}

func (zw *ZhengWu) Complete(ctime time.Time) {
	zw.endTime = ctime
}

func (zw *ZhengWu) IsCompleted(ctime time.Time) bool {
	return zw.endTime.Before(ctime.Add(time.Second)) // 减少1秒
}

func (zw *ZhengWu) EncodeClient() *shared_proto.ZhengWuProto {
	proto := &shared_proto.ZhengWuProto{}

	proto.Data = zw.data.Proto
	proto.EndTime = timeutil.Marshal32(zw.endTime)

	return proto
}

func (zw *ZhengWu) EncodeServer() *server_proto.ZhengWuServerProto {
	proto := &server_proto.ZhengWuServerProto{}

	proto.Data = zw.data.Id
	proto.EndTime = timeutil.Marshal64(zw.endTime)

	return proto
}
