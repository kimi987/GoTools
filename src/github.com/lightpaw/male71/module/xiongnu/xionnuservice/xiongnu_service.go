package xionnuservice

import (
	"context"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/xiongnu"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuface"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuinfo"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/pbutil"
	"time"
	"github.com/lightpaw/male7/util/ctxfunc"
)

func NewXiongNuService(timeService iface.TimeService, // 时间
	configDatas *config.ConfigDatas,                  // 配置
	tickerService iface.TickerService,
	dbService iface.DbService) *XiongNuService {
	s := &XiongNuService{
		configDatas:   configDatas,
		tickerService: tickerService,
		dbService:     dbService,
		infoMap:       NewinfoMap(),
		todayJoinMap:  xiongnuinfo.NewTodayJoinMap(),
	}

	s.load()

	return s
}

//gogen:iface
type XiongNuService struct {
	configDatas        *config.ConfigDatas // 配置
	tickerService      iface.TickerService // 定时事件
	dbService          iface.DbService     // 数据库服务
	infoMap            *infoMap
	todayJoinMap       *xiongnuinfo.TodayJoinMap
	lastResetDailyTime time.Time // 上次每日重置时间
}

func (m *XiongNuService) XiongNuInfoMsg(guildId int64) pbutil.Buffer {
	info, exist := m.infoMap.Get(guildId)
	if exist && info.BaseId() != 0 {
		return xiongnu.NewS2cInfoBroadcastMsg(must.Marshal(info.EncodeClient()))
	}
	return nil
}

func (m *XiongNuService) IsTodayStarted(heroId int64) bool {
	_, todayJoin := m.todayJoinMap.Get(heroId)
	return todayJoin
}

func (m *XiongNuService) SetTodayStarted(heroId int64) {
	m.todayJoinMap.Set(heroId, heroId)
}

func (m *XiongNuService) TodayJoinMap() *xiongnuinfo.TodayJoinMap {
	return m.todayJoinMap
}

func (m *XiongNuService) GetRInfo(guildId int64) xiongnuface.RResistXiongNuInfo {
	info := m.GetInfo(guildId)
	if info == nil {
		return nil
	}
	return info
}

func (m *XiongNuService) GetInfo(guildId int64) xiongnuface.ResistXiongNuInfo {
	info, _ := m.infoMap.Get(guildId)
	if info == nil {
		return nil
	}
	return info
}

func (m *XiongNuService) load() {
	var datas []byte
	err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		datas, err = m.dbService.LoadKey(ctx, server_proto.Key_XiongNu)
		return
	})
	if err != nil {
		logrus.WithError(err).Panicf("开服加载匈奴数据报错!")
	}

	if len(datas) <= 0 {
		return
	}

	allProto := &server_proto.AllXiongNuServerProto{}
	err = allProto.Unmarshal(datas)
	if err != nil {
		logrus.WithError(err).Panicf("开服解析匈奴数据报错!")
	}

	m.lastResetDailyTime = timeutil.Unix64(allProto.LastResetDailyTime)

	for _, id := range allProto.TodayJoinIds {
		m.todayJoinMap.Set(id, id)
	}

	for _, proto := range allProto.Datas {
		data := m.configDatas.GetResistXiongNuData(proto.Level)
		check.PanicNotTrue(data != nil, "开服解析匈奴等级数据，没找到!%v", proto)

		info := xiongnuinfo.NewResistXiongNuInfo(m.configDatas,
			proto.GuildId,
			data,
			proto.Defenders,
			timeutil.Unix64(proto.StartTime))

		info.Unmarshal(proto)

		m.infoMap.Set(info.GuildId(), info)
	}

	// 看是否需要重置
	if prevTickTime := m.tickerService.GetDailyTickTime().GetPrevTickTime(); m.lastResetDailyTime.Before(prevTickTime) {
		m.lastResetDailyTime = prevTickTime
		m.todayJoinMap = xiongnuinfo.NewTodayJoinMap()
	}
}

func (m *XiongNuService) ResetDaily(newDailyTime time.Time) {
	m.todayJoinMap = xiongnuinfo.NewTodayJoinMap()
	m.lastResetDailyTime = newDailyTime
}

func (m *XiongNuService) Save() {
	if err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
		return m.dbService.SaveKey(ctx, server_proto.Key_XiongNu, must.Marshal(m.encode()))
	}); err != nil {
		logrus.WithError(err).Error("保存匈奴数据失败")
	}
}

func (m *XiongNuService) encode() *server_proto.AllXiongNuServerProto {
	proto := &server_proto.AllXiongNuServerProto{Datas: make([]*server_proto.XiongNuServerProto, 0, m.infoMap.Count())}

	m.infoMap.IterCb(func(key int64, v xiongnuface.ResistXiongNuInfo) {
		proto.Datas = append(proto.Datas, v.Encode())
	})

	for ids := range m.todayJoinMap.Keys() {
		proto.TodayJoinIds = append(proto.TodayJoinIds, ids...)
	}
	proto.LastResetDailyTime = timeutil.Marshal64(m.lastResetDailyTime)

	return proto
}

func (m *XiongNuService) IsStarted(guildId int64) bool {
	_, exist := m.infoMap.Get(guildId)
	return exist
}

func (m *XiongNuService) AddInfo(info xiongnuface.ResistXiongNuInfo) {
	m.infoMap.Set(info.GuildId(), info)
}

func (m *XiongNuService) RemoveInfo(info xiongnuface.ResistXiongNuInfo) {
	m.infoMap.RemoveIfSame(info.GuildId(), info)
}

func (m *XiongNuService) WalkInfo(walkFunc xiongnuface.WalkInfoFunc) {
	for tp := range m.infoMap.Iter() {
		for _, tupleInfo := range tp {
			walkFunc(tupleInfo.Val)
		}
	}
}
