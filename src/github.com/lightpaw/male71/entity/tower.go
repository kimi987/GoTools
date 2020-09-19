package entity

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func newTower() *Tower {
	return &Tower{}
}

type Tower struct {
	// 挑战次数，最大挑战次数，从配置中读取
	challengeTimes         uint64
	nextResetChallengeTime time.Time

	// 今日打到第几层，0表示还没开始挑战
	currentFloor uint64

	// 历史通关的最高楼层
	historyMaxFloor     uint64
	historyMaxFloorTime time.Time // 历史通关的最高楼层的时间

	// 扫荡最高楼层
	autoMaxFloor uint64

	collectedBoxFloors []uint64
}

func (t *Tower) resetDaily() {
	t.currentFloor = 0
	t.challengeTimes = 0
	t.nextResetChallengeTime = time.Time{}
}

func (t *Tower) ChallengeTimes() uint64 {
	return t.challengeTimes
}

func (t *Tower) IncreseChallengeTimes(ctime time.Time) uint64 {
	if t.nextResetChallengeTime.Before(ctime) {
		t.challengeTimes = 0
	}

	t.challengeTimes++
	return t.challengeTimes
}

func (t *Tower) SetNextResetChallengeTime(toSet time.Time) {
	t.nextResetChallengeTime = toSet
}

func (t *Tower) GetNextResetChallengeTime() time.Time {
	return t.nextResetChallengeTime
}

func (t *Tower) CurrentFloor() uint64 {
	return t.currentFloor
}

func (t *Tower) IncreseCurrentFloor(ctime time.Time, reduceAutoFloor uint64) (historyMaxFloorChanged bool) {
	t.currentFloor++
	if t.historyMaxFloor < t.currentFloor {
		t.historyMaxFloor = t.currentFloor
		t.historyMaxFloorTime = ctime
		historyMaxFloorChanged = true

		if f := u64.Sub(t.currentFloor, reduceAutoFloor); f > t.autoMaxFloor {
			t.autoMaxFloor = f
		}
	}

	return
}

func (t *Tower) SetCurrentFloorToAuto() {
	if t.currentFloor < t.autoMaxFloor {
		t.currentFloor = t.autoMaxFloor
	}
}

func (t *Tower) SetCurrentFloorToZero() {
	t.currentFloor = 0
}

func (t *Tower) HistoryMaxFloor() uint64 {
	return t.historyMaxFloor
}

func (t *Tower) HistoryMaxFloorTime() time.Time {
	return t.historyMaxFloorTime
}

func (t *Tower) AutoMaxFloor() uint64 {
	return t.autoMaxFloor
}

//func (t *Tower) IncreaseAutoFloor() {
//	if t.autoMaxFloor+1 == t.currentFloor {
//		t.autoMaxFloor = t.currentFloor
//	}
//}

func (t *Tower) IsBoxCollected(floor uint64) bool {
	return u64.Contains(t.collectedBoxFloors, floor)
}

func (t *Tower) CollectBox(floor uint64) {
	t.collectedBoxFloors = u64.AddIfAbsent(t.collectedBoxFloors, floor)
}

func (t *Tower) EncodeClient() *shared_proto.HeroTowerProto {
	proto := &shared_proto.HeroTowerProto{}

	proto.ChallengeTimes = u64.Int32(t.challengeTimes)
	proto.NextResetChallengeTime = timeutil.Marshal32(t.nextResetChallengeTime)
	proto.CurrentFloor = u64.Int32(t.currentFloor)
	proto.HistoryMaxFloor = u64.Int32(t.historyMaxFloor)
	proto.AutoMaxFloor = u64.Int32(t.autoMaxFloor)
	proto.CollectedBoxFloors = u64.Int32Array(t.collectedBoxFloors)

	return proto
}

func (t *Tower) encodeServer() *server_proto.HeroTowerServerProto {
	proto := &server_proto.HeroTowerServerProto{}

	proto.ChallengeTimes = t.challengeTimes
	proto.NextResetChallengeTime = timeutil.Marshal64(t.nextResetChallengeTime)
	proto.CurrentFloor = t.currentFloor
	proto.HistoryMaxFloor = t.historyMaxFloor
	proto.HistoryMaxFloorTime = timeutil.Marshal64(t.historyMaxFloorTime)
	proto.AutoMaxFloor = t.autoMaxFloor
	proto.CollectedBoxFloors = t.collectedBoxFloors

	return proto
}

func (t *Tower) unmarshal(proto *server_proto.HeroTowerServerProto) {
	if proto == nil {
		return
	}

	t.challengeTimes = proto.ChallengeTimes
	t.nextResetChallengeTime = timeutil.Unix64(proto.NextResetChallengeTime)
	t.currentFloor = proto.CurrentFloor
	t.historyMaxFloor = proto.HistoryMaxFloor
	t.historyMaxFloorTime = timeutil.Unix64(proto.HistoryMaxFloorTime)
	t.autoMaxFloor = proto.AutoMaxFloor
	t.collectedBoxFloors = proto.CollectedBoxFloors
}
