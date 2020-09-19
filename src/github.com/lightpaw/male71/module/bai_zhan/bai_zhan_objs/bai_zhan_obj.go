package bai_zhan_objs

import (
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/bai_zhan_data"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

type RHeroBaiZhanObj interface {
	Id() int64
	ChallengeTimes() uint64
	RecordVsn() (version int32)
	LevelData() *bai_zhan_data.JunXianLevelData
	Point() (curPoint uint64)
	Rank() uint64
	IsCollectDailySalary() bool
	HasCollectedJunXianPrize(data *bai_zhan_data.JunXianPrizeData) (hasCollected bool)
	LastCollectJunXianPrizeId() uint64
	LevelChangeType(levelUpMaxRank, levelDownMinRank uint64) shared_proto.LevelChangeType
	HistoryMaxJunXianLevelData() *bai_zhan_data.JunXianLevelData
	HistoryMaxPoint(levelData *bai_zhan_data.JunXianLevelData) uint64
}

func NewBaiZhanObj(id int64, heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot, maxJunXianLevel int, minLevelData *bai_zhan_data.JunXianLevelData, ctime time.Time) *HeroBaiZhanObj {
	return &HeroBaiZhanObj{
		id:                         id,                 // 玩家id
		heroSnapshotGetter:         heroSnapshotGetter, // 玩家信息获取
		junXianLevelData:           minLevelData,       // 玩家军衔等级默认一级
		historyMaxJunXianLevelData: minLevelData,       // 历史最高军衔等级，默认一级
		historyMaxPoints:           make([]uint64, maxJunXianLevel),
		recordVsn:                  1, // 默认是1，因为客户端第一次发送上来的是0
		lastPointChangeTime:        ctime,
	}
}

// 玩家百战
type HeroBaiZhanObj struct {
	id int64 // 玩家id

	heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot // 玩家镜像数据获得方法

	// 军衔等级
	junXianLevelData           *bai_zhan_data.JunXianLevelData // 军衔等级数据
	lastJunXianLevelData       *bai_zhan_data.JunXianLevelData // 上次的军衔等级(每天第一次进百战看完进阶动画发这个,干掉客户端会有问题,因此下面新增字段oldjunXianLevelData用于排行榜)
	oldjunXianLevelData        *bai_zhan_data.JunXianLevelData
	isJunXianBeenRemoved       bool                            // 军衔等级是否被移除
	historyMaxJunXianLevelData *bai_zhan_data.JunXianLevelData // 历史最高军衔等级

	challengeTimes      uint64    // 挑战次数
	point               uint64    // 今日积分
	lastPointChangeTime time.Time // 最近一次积分变更的时间

	historyMaxPoints []uint64 // 每个军衔等级对应的历史最高积分

	rank uint64 // 排名

	// 奖励
	collectedDailySalary        bool   // 是否领取了今日的俸禄
	lastCollectedJunXianPrizeId uint64 // 最后领取的军衔奖励id

	// 镜像
	combatMirror *shared_proto.CombatPlayerProto // 玩家的镜像，可能为空

	// 镜像战斗力
	combatMirrorFightAmount uint64

	// 百战版本号
	recordVsn int32 // 记录的版本号
}

func (h *HeroBaiZhanObj) Id() int64 {
	return h.id
}

func (h *HeroBaiZhanObj) CombatMirror() *shared_proto.CombatPlayerProto {
	return h.combatMirror
}

func (h *HeroBaiZhanObj) SetCombatMirror(mirror *shared_proto.CombatPlayerProto) {
	h.combatMirror = mirror
	h.combatMirrorFightAmount = u64.FromInt32(mirror.TotalFightAmount)
}

func (h *HeroBaiZhanObj) CombatMirrorFightAmount() uint64 {
	return h.combatMirrorFightAmount
}

func (h *HeroBaiZhanObj) LevelData() *bai_zhan_data.JunXianLevelData {
	return h.junXianLevelData
}

func (o *HeroBaiZhanObj) Rank() uint64 {
	return o.rank
}

func (o *HeroBaiZhanObj) SetRank(toSet uint64) {
	o.rank = toSet
}

func (o *HeroBaiZhanObj) LastPointChangeTime() time.Time {
	return o.lastPointChangeTime
}

func (o *HeroBaiZhanObj) SetLastPointChangeTime(toSet time.Time) {
	o.lastPointChangeTime = toSet
}

func (o *HeroBaiZhanObj) Less(obj *HeroBaiZhanObj) bool {
	// 层级高的在前面
	if o.junXianLevelData != obj.junXianLevelData {
		return o.junXianLevelData.Level > obj.junXianLevelData.Level
	}

	// 层级相同的，时间小的在前面
	if o.Point() != obj.Point() {
		return o.Point() > obj.Point()
	}

	if o.lastPointChangeTime != obj.lastPointChangeTime {
		return o.lastPointChangeTime.Before(obj.lastPointChangeTime)
	}

	return o.id < obj.id
}

func (h *HeroBaiZhanObj) checkMaxJunXianLevel(jxData *bai_zhan_data.JunXianLevelData) (historyLevelChanged bool) {
	if historyLevelChanged = jxData.Level > h.historyMaxJunXianLevelData.Level; historyLevelChanged {
		h.historyMaxJunXianLevelData = jxData
		//if s := h.heroSnapshotGetter(h.id); s != nil {
		//	s.BaiZhanJunXianLevel = jxData.Level
		//	s.EncodeClient().JunXianLevel = int32(jxData.Level)
		//}
	}
	return
}

func (h *HeroBaiZhanObj) ResetLevelData(toSet *bai_zhan_data.JunXianLevelData, newPointChangeTime time.Time) (levelChanged bool, historyLevelChanged bool) {
	historyLevelChanged = h.checkMaxJunXianLevel(toSet)

	h.lastJunXianLevelData = h.junXianLevelData
	h.oldjunXianLevelData = h.junXianLevelData
	if h.junXianLevelData == toSet {
		return
	}

	h.junXianLevelData = toSet
	h.lastPointChangeTime = newPointChangeTime
	levelChanged = true
	return
}

func (h *HeroBaiZhanObj) LastJunXianLevelData() *bai_zhan_data.JunXianLevelData {
	return h.lastJunXianLevelData
}

func (h *HeroBaiZhanObj) OldJunXianLevelData() *bai_zhan_data.JunXianLevelData {
	return h.oldjunXianLevelData
}

func (h *HeroBaiZhanObj) ClearLastJunXianLevelData() {
	h.lastJunXianLevelData = nil
}

func (h *HeroBaiZhanObj) IsJunXianBeenRemoved() bool {
	return h.isJunXianBeenRemoved
}

func (h *HeroBaiZhanObj) RemoveJunXian() {
	if h.isJunXianBeenRemoved {
		// 已经被移除了
		return
	}

	h.isJunXianBeenRemoved = true

	h.lastJunXianLevelData = h.junXianLevelData
	h.oldjunXianLevelData = h.junXianLevelData
	h.junXianLevelData = h.junXianLevelData.ReaddJunXianLevelData

	h.combatMirror = nil // 镜像也要清掉
	h.point = 0          // 积分也清掉
}

// 重新把被移除的军衔加回来
func (h *HeroBaiZhanObj) ReaddRemovedJunXian(ctime time.Time) (readdSuccess bool) {
	if !h.IsJunXianBeenRemoved() {
		return
	}

	h.isJunXianBeenRemoved = false
	h.lastPointChangeTime = ctime

	// 再清一遍
	h.combatMirror = nil
	h.point = 0

	return true
}

func (h *HeroBaiZhanObj) ChallengeTimes() uint64 {
	return h.challengeTimes
}

func (h *HeroBaiZhanObj) IncreChallengeTimes() (newChallengeTimes uint64) {
	h.challengeTimes++
	return h.ChallengeTimes()
}

func (h *HeroBaiZhanObj) AddPoint(toAdd uint64, time time.Time) (newPoint uint64) {
	h.point += toAdd
	h.lastPointChangeTime = time
	h.historyMaxPoints[h.LevelData().Level-1] = u64.Max(h.HistoryMaxPoint(h.LevelData()), h.point)
	return h.point
}

func (h *HeroBaiZhanObj) Point() (curPoint uint64) {
	return h.point
}

func (h *HeroBaiZhanObj) ClearPoint() {
	h.point = 0
}

func (h *HeroBaiZhanObj) HistoryMaxPoint(levelData *bai_zhan_data.JunXianLevelData) uint64 {
	return h.historyMaxPoints[levelData.Level-1]
}

func (h *HeroBaiZhanObj) HistoryMaxJunXianLevelData() *bai_zhan_data.JunXianLevelData {
	return h.historyMaxJunXianLevelData
}

func (h *HeroBaiZhanObj) IsCollectDailySalary() bool {
	return h.collectedDailySalary
}

func (h *HeroBaiZhanObj) CollectDailySalary() {
	h.collectedDailySalary = true
}

func (h *HeroBaiZhanObj) CollectJunXianPrize(data *bai_zhan_data.JunXianPrizeData) {
	h.lastCollectedJunXianPrizeId = data.Id
}

func (h *HeroBaiZhanObj) LastCollectJunXianPrizeId() uint64 {
	return h.lastCollectedJunXianPrizeId
}

func (h *HeroBaiZhanObj) HasCollectedJunXianPrize(data *bai_zhan_data.JunXianPrizeData) (hasCollected bool) {
	return h.lastCollectedJunXianPrizeId >= data.Id
}

func (h *HeroBaiZhanObj) RecordVsn() (version int32) {
	return h.recordVsn
}

func (h *HeroBaiZhanObj) IncreRecordVsn() (version int32) {
	h.recordVsn++
	return h.recordVsn
}

func (h *HeroBaiZhanObj) ResetChallengeTimes() {
	// 次数置空
	h.challengeTimes = 0
}

func (h *HeroBaiZhanObj) ResetDaily() {
	// 次数置空
	h.challengeTimes = 0
	// 积分置空
	h.point = 0
	// 重置以后，积分为0，没有上榜
	h.rank = 0
	// 俸禄没有领取
	h.collectedDailySalary = false
}

func (o *HeroBaiZhanObj) LevelChangeType(levelUpMaxRank, levelDownMinRank uint64) shared_proto.LevelChangeType {
	if o.junXianLevelData.NextLevel != nil && o.rank <= levelUpMaxRank && o.point >= o.junXianLevelData.LevelUpPoint {
		// 升级
		return shared_proto.LevelChangeType_LEVEL_UP
	}

	if o.junXianLevelData.PrevLevel == nil || o.rank <= levelUpMaxRank + o.junXianLevelData.MinKeepLevelCount {
		// 一定是保级
		return shared_proto.LevelChangeType_LEVEL_KEEP
	}

	if o.point <= o.junXianLevelData.LevelDownPoint && o.rank > levelDownMinRank {
		// 一定是降级
		return shared_proto.LevelChangeType_LEVEL_DOWN
	}

	// 一定是保级
	return shared_proto.LevelChangeType_LEVEL_KEEP
}

func (h *HeroBaiZhanObj) EncodeClient() *shared_proto.BaiZhanObjProto {
	proto := &shared_proto.BaiZhanObjProto{}

	proto.ChallengeTimes = u64.Int32(h.challengeTimes)
	proto.Point = u64.Int32(h.point)
	proto.IsCollectSalary = h.collectedDailySalary
	proto.LastCollectedJunXianPrizeId = u64.Int32(h.lastCollectedJunXianPrizeId)
	proto.JunXianLevel = u64.Int32(h.junXianLevelData.Level)
	if h.lastJunXianLevelData != nil {
		proto.LastJunXianLevel = u64.Int32(h.lastJunXianLevelData.Level)
	}

	proto.HistoryMaxJunXianLevel = u64.Int32(h.historyMaxJunXianLevelData.Level)
	proto.HistoryMaxPoints = u64.Int32Array(h.historyMaxPoints)

	return proto
}

func (o *HeroBaiZhanObj) Encode4Rank(levelUpMaxRank, levelDownMinRank uint64) *shared_proto.BaiZhanRankObjProto {
	proto := &shared_proto.BaiZhanRankObjProto{}
	snapshot := o.heroSnapshotGetter(o.id)
	if snapshot == nil {
		return nil
	}

	proto.Basic = snapshot.EncodeBasic4Client()
	proto.Point = u64.Int32(o.Point())
	proto.LevelChangeType = o.LevelChangeType(levelUpMaxRank, levelDownMinRank)
	proto.FightAmount = u64.Int32(o.combatMirrorFightAmount)

	return proto
}

func (h *HeroBaiZhanObj) EncodeServer() *server_proto.BaiZhanObjServerProto {
	proto := &server_proto.BaiZhanObjServerProto{}

	proto.Id = h.id
	proto.ChallengeTimes = h.challengeTimes
	proto.Point = h.point
	proto.LastPointChangeTime = timeutil.Marshal64(h.lastPointChangeTime)
	proto.IsCollectSalary = h.collectedDailySalary

	proto.JunXianLevel = h.junXianLevelData.Level
	proto.IsJunXianBeenRemoved = h.IsJunXianBeenRemoved()

	if h.lastJunXianLevelData != nil {
		proto.LastJunXianLevel = h.lastJunXianLevelData.Level
	}

	proto.HistoryMaxJunXianLevel = h.historyMaxJunXianLevelData.Level
	proto.HistoryMaxPoints = h.historyMaxPoints

	proto.LastCollectedJunXianPrizeId = h.lastCollectedJunXianPrizeId

	proto.Mirror = h.combatMirror

	proto.MirrorFightAmount = h.combatMirrorFightAmount

	return proto
}

func (h *HeroBaiZhanObj) Unmarshall(proto *server_proto.BaiZhanObjServerProto, junXianLevelData *config.JunXianLevelDataConfig) {
	h.id = proto.GetId()
	h.challengeTimes = proto.GetChallengeTimes()
	h.lastPointChangeTime = timeutil.Unix64(proto.GetLastPointChangeTime())
	h.collectedDailySalary = proto.GetIsCollectSalary()
	h.lastCollectedJunXianPrizeId = proto.GetLastCollectedJunXianPrizeId()

	h.junXianLevelData = junXianLevelData.Must(proto.GetJunXianLevel())
	h.isJunXianBeenRemoved = proto.GetIsJunXianBeenRemoved()

	if !h.isJunXianBeenRemoved {
		// 没被移除军衔
		h.combatMirror = proto.GetMirror()
		h.point = proto.GetPoint()
	}

	h.combatMirrorFightAmount = proto.MirrorFightAmount

	if proto.GetLastJunXianLevel() > 0 {
		h.lastJunXianLevelData = junXianLevelData.Must(proto.GetLastJunXianLevel())
		h.oldjunXianLevelData = h.lastJunXianLevelData
	}

	h.historyMaxJunXianLevelData = junXianLevelData.Must(proto.GetHistoryMaxJunXianLevel())
	h.checkMaxJunXianLevel(h.junXianLevelData)

	copy(h.historyMaxPoints, proto.HistoryMaxPoints)
}
