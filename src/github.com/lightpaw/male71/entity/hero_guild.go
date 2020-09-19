package entity

import (
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/util/recovtimes"
	"github.com/lightpaw/male7/config/heroinit"
)

func newHeroGuild(initData *heroinit.HeroInitData, ctime time.Time) *hero_guild {
	g := &hero_guild{
		workshopOutputTimes: NewRecoverableTimes(ctime, initData.WorkshopOutputRecoveryDuration, initData.WorkshopOutputMaxTimes),
	}

	g.workshopOutputTimes.SetTimes(initData.WorkshopOutputInitTimes, ctime)
	g.collectedTaskStages = make(map[uint64]map[int32]struct{})
	return g
}

type hero_guild struct {
	guildId       int64
	joinGuildTime time.Time // 加入联盟时间

	guildContributionCoin uint64   // 贡献币
	guildDonateTimes      []uint64 // 捐献次数

	joinGuildIds []int64 // 申请加入的帮派id

	beenInvateGuildIds []int64 // 被邀请加入的帮派id

	collectedFirstJoinGuildPrize bool // 是否有领取首次加入联盟奖励
	collectedDailyGuildRankPrize bool // 是否已经领取当日联盟排名奖励

	nextNotifyGuildTime time.Time // 下次提醒加入联盟的时间

	workshopOutputTimes *recovtimes.RecoverTimes // 联盟工坊今日生产次数

	collectedTaskStages map[uint64]map[int32]struct{} // 联盟周任务的领奖情况
}

// 万一当前帮派没有删除掉，每日更新的时候，处理一次 TODO

func (hero *hero_guild) GetJoinGuildIds() []int64 {
	return hero.joinGuildIds
}

func (hero *hero_guild) AddJoinGuildIds(joinGuildId int64) {
	hero.joinGuildIds = i64.AddIfAbsent(hero.joinGuildIds, joinGuildId)
}

func (hero *hero_guild) RemoveJoinGuildIds(joinGuildId int64) bool {
	var idx int
	hero.joinGuildIds, idx = i64.LeftShiftRemoveIfPresentReturnIndex(hero.joinGuildIds, joinGuildId)
	return idx >= 0
}

func (hero *hero_guild) GetBeenInvateGuildIds() []int64 {
	return hero.beenInvateGuildIds
}

func (hero *hero_guild) AddBeenInvateGuildId(toAdd int64) {
	hero.beenInvateGuildIds = i64.AddIfAbsent(hero.beenInvateGuildIds, toAdd)
}

func (hero *hero_guild) RemoveBeenInvateGuildId(toRemove int64) bool {
	var idx int
	hero.beenInvateGuildIds, idx = i64.LeftShiftRemoveIfPresentReturnIndex(hero.beenInvateGuildIds, toRemove)

	return idx >= 0
}

func (hero *hero_guild) RemoveFirstBeenInvateGuildId() int64 {
	if len(hero.beenInvateGuildIds) > 0 {
		first := hero.beenInvateGuildIds[0]
		hero.beenInvateGuildIds = i64.LeftShift(hero.beenInvateGuildIds, 0, 1)
		return first
	}

	return 0
}

func (hero *hero_guild) ClearJoinGuildIds() []int64 {
	ids := hero.joinGuildIds
	hero.joinGuildIds = nil
	return ids
}

func (hero *hero_guild) GuildId() int64 {
	return hero.guildId
}

func (hero *hero_guild) SetGuild(id int64) {
	hero.guildId = id
}

func (hero *hero_guild) JoinGuildTime() time.Time {
	return hero.joinGuildTime
}

func (hero *hero_guild) SetJoinGuildTime(joinGuildTime time.Time) {
	hero.joinGuildTime = joinGuildTime
}

func (hero *hero_guild) IsFirstJoinGuildPrizeCollected() bool {
	return hero.collectedFirstJoinGuildPrize
}

func (hero *hero_guild) CollectFirstJoinGuildPrize() {
	hero.collectedFirstJoinGuildPrize = true
}

func (hero *hero_guild) GetGuildContributionCoin() uint64 {
	return hero.guildContributionCoin
}

func (hero *hero_guild) AddGuildContributionCoin(toAdd uint64) {
	hero.guildContributionCoin += toAdd
}

func (hero *hero_guild) ReduceGuildContributionCoin(toReduce uint64) {
	hero.guildContributionCoin = u64.Sub(hero.guildContributionCoin, toReduce)
}

func (hero *hero_guild) GetDonateTimes(seq uint64) (uint64, bool) {
	if seq > 0 && seq <= uint64(len(hero.guildDonateTimes)) {
		return hero.guildDonateTimes[seq-1], true
	}

	return 0, false
}

func (hero *hero_guild) AddDonateTimes(seq uint64) {
	if seq > 0 && seq <= uint64(len(hero.guildDonateTimes)) {
		hero.guildDonateTimes[seq-1]++
	}
}

func (hero *hero_guild) NextNotifyGuildTime() time.Time {
	return hero.nextNotifyGuildTime
}

func (hero *hero_guild) SetNextNotifyGuildTime(toSet time.Time) {
	hero.nextNotifyGuildTime = toSet
}

func (hero *hero_guild) CollectedDailyGuildRankPrize() bool {
	return hero.collectedDailyGuildRankPrize
}

func (hero *hero_guild) SetCollectedDailyGuildRankPrize(collected bool) {
	hero.collectedDailyGuildRankPrize = collected
}

func (hero *hero_guild) guildResetDaily() {
	for idx := range hero.guildDonateTimes {
		hero.guildDonateTimes[idx] = 0
	}
	hero.collectedDailyGuildRankPrize = false
}

func (hero *hero_guild) GetWorkshopOutputTimes() *recovtimes.RecoverTimes {
	return hero.workshopOutputTimes
}

func (hero *hero_guild) TrySetCollectedTaskStage(id uint64, stage int32) bool {
	m := hero.collectedTaskStages[id]
	if m == nil {
		m = make(map[int32]struct{})
		m[stage] = struct{}{}
		hero.collectedTaskStages[id] = m
		return true
	}
	_, ok := m[stage]
	if !ok {
		m[stage] = struct{}{}
		return true
	}
	return false
}

func (hero *hero_guild) guildResetWeekly() {
	if len(hero.collectedTaskStages) > 0 {
		hero.collectedTaskStages = make(map[uint64]map[int32]struct{})
	}
}

func (hero *hero_guild) CopyGuildResetWeekly() map[uint64][]int {
	length := len(hero.collectedTaskStages)
	if length <= 0 {
		return nil
	}
	m := make(map[uint64][]int, length)
	for id, stages := range hero.collectedTaskStages {
		arr := make([]int, 0, len(stages))
		for stage, _ := range stages {
			arr = append(arr, int(stage))
		}
		m[id] = arr
	}
	return m
}
