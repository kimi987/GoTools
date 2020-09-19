package entity

import (
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/towerdata"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/eapache/queue"
)

func newHeroSecretTower() *HeroSecretTower {
	return &HeroSecretTower{
		firstPassPrizeMap: make(map[uint64]struct{}),
		records:           queue.New(),
	}
}

type HeroSecretTower struct {
	challengeTimes                   uint64 // 挑战次数
	helpTimes                        uint64 // 协助次数
	maxOpenSecretTowerId             uint64 // 最大的开启了的塔的id
	firstPassPrizeMap                map[uint64]struct{}
	todayAddCollectGuildContribution uint64 // 今日获得的最大的联盟贡献
	historyChallengeTimes            uint64 // 历史挑战次数

	// 密室记录
	records *queue.Queue
}

func (t *HeroSecretTower) AddRecord(record *shared_proto.SecretRecordProto, recordLimit uint64) {
	if recordLimit <= 0 {
		return
	}

	if uint64(t.records.Length()) >= recordLimit {
		t.records.Remove()
	}

	t.records.Add(record)
}

func (t *HeroSecretTower) GetRecords() (records []*shared_proto.SecretRecordProto) {
	for i := 0; i < t.records.Length(); i++ {
		r := t.records.Get(i)
		if r != nil {
			if record, ok := r.(*shared_proto.SecretRecordProto); ok {
				records = append(records, record)
			}
		}
	}
	return
}

func (t *HeroSecretTower) TryOpenAndGiveDefaultTimes(unlockSecretTower *towerdata.SecretTowerData) (openNewSecretTower bool) {
	if t.maxOpenSecretTowerId < unlockSecretTower.Id {
		t.maxOpenSecretTowerId = unlockSecretTower.Id
		openNewSecretTower = true
	}

	return
}

func (t *HeroSecretTower) HasOpen(secretTowerId uint64) bool {
	return t.maxOpenSecretTowerId >= secretTowerId
}

func (t *HeroSecretTower) MaxOpenSecretTowerId() uint64 {
	return t.maxOpenSecretTowerId
}

func (t *HeroSecretTower) ChallengeTimes() uint64 {
	return t.challengeTimes
}

func (t *HeroSecretTower) HasEnoughHelpTimes(maxHelpTimes uint64) bool {
	return t.helpTimes < maxHelpTimes
}

func (t *HeroSecretTower) HelpTimes() uint64 {
	return t.helpTimes
}

func (t *HeroSecretTower) IncreHelpTimes() (newTimes uint64) {
	t.helpTimes++
	return t.helpTimes
}

func (t *HeroSecretTower) IncreChallengeTimes() uint64 {
	t.challengeTimes++
	return t.challengeTimes
}

func (t *HeroSecretTower) IncreHistoryChallengeTimes() uint64 {
	t.historyChallengeTimes++
	return t.historyChallengeTimes
}

func (t *HeroSecretTower) HistoryChallengeTimes() uint64 {
	return t.historyChallengeTimes
}

func (t *HeroSecretTower) HasFirstPass(towerId uint64) bool {
	_, ok := t.firstPassPrizeMap[towerId]
	return ok
}

func (t *HeroSecretTower) GiveFirstPassPrize(towerId uint64) {
	t.firstPassPrizeMap[towerId] = struct{}{}
}

//func (h *HeroSecretTower) AddTodayGuildContributionAmount(maxGuildContribution, guildContribution uint64) (addAmount uint64) {
//	amount := u64.Sub(maxGuildContribution, h.todayAddCollectGuildContribution)
//
//	if amount < guildContribution {
//		addAmount = amount
//	} else {
//		addAmount = guildContribution
//	}
//
//	h.todayAddCollectGuildContribution += addAmount
//
//	return
//}

func (h *HeroSecretTower) AddTodayGuildContributionAmount(guildContribution uint64) uint64 {
	h.todayAddCollectGuildContribution += guildContribution
	return h.todayAddCollectGuildContribution
}

func (t *HeroSecretTower) ResetDaily() {
	t.challengeTimes = 0
	t.helpTimes = 0
}

func (t *HeroSecretTower) EncodeClient() *shared_proto.HeroSecretTowerProto {
	proto := &shared_proto.HeroSecretTowerProto{}

	proto.ChallengeTimes = u64.Int32(t.challengeTimes)
	proto.HelpTimes = u64.Int32(t.helpTimes)
	proto.MaxOpenTowerId = u64.Int32(t.maxOpenSecretTowerId)
	proto.HistoryChallengeTimes = u64.Int32(t.historyChallengeTimes)

	if len(t.firstPassPrizeMap) > 0 {
		proto.CollectedFirstPrizeId = make([]int32, 0, len(t.firstPassPrizeMap))
		for id := range t.firstPassPrizeMap {
			proto.CollectedFirstPrizeId = append(proto.CollectedFirstPrizeId, u64.Int32(id))
		}
	}

	return proto
}

func (t *HeroSecretTower) encodeServer() *server_proto.HeroSecretTowerServerProto {
	proto := &server_proto.HeroSecretTowerServerProto{}

	proto.ChallengeTimes = t.challengeTimes
	proto.HelpTimes = t.helpTimes
	proto.MaxOpenSecretTowerId = t.maxOpenSecretTowerId
	proto.TodayAddCollectGuildContribution = t.todayAddCollectGuildContribution
	proto.HistoryChallengeTimes = t.historyChallengeTimes

	proto.FirstPassSecretTowerId = make([]uint64, 0, len(t.firstPassPrizeMap))
	for id, _ := range t.firstPassPrizeMap {
		proto.FirstPassSecretTowerId = append(proto.FirstPassSecretTowerId, id)
	}

	proto.Records = t.GetRecords()

	return proto
}

func (t *HeroSecretTower) unmarshal(proto *server_proto.HeroSecretTowerServerProto, heroTower *Tower, configDatas *config.ConfigDatas, ctime time.Time) {
	if proto != nil {
		t.challengeTimes = proto.GetChallengeTimes()
		t.helpTimes = proto.GetHelpTimes()
		t.maxOpenSecretTowerId = proto.MaxOpenSecretTowerId
		t.todayAddCollectGuildContribution = proto.GetTodayAddCollectGuildContribution()
		t.historyChallengeTimes = proto.GetHistoryChallengeTimes()

		for _, id := range proto.FirstPassSecretTowerId {
			t.GiveFirstPassPrize(id)
		}

		for _, v := range proto.Records {
			t.records.Add(v)
		}
	}
}
