package entity

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"time"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
)

type hero_treasury_tree struct {
	// 摇钱树已浇水次数
	waterTimes uint64

	// 可领取的奖励季节
	collectSeason shared_proto.Season
	collectTime   time.Time

	// 帮助我的玩家列表（多线程读取，copyOnWrite)
	helpMeHeroIds []int64
	helpMeSeason  []shared_proto.Season

	// 今日浇水列表
	waterHeroIds []int64
}

func (hero *Hero) TreasuryTree() *hero_treasury_tree {
	return hero.treasuryTree
}

func (t *hero_treasury_tree) WaterTimes() uint64 {
	return t.waterTimes
}

func (t *hero_treasury_tree) IncreseWaterTimes() uint64 {
	t.waterTimes++
	return t.waterTimes
}

func (t *hero_treasury_tree) IsWatered(heroId int64) bool {
	return i64.Contains(t.waterHeroIds, heroId)
}

func (t *hero_treasury_tree) AddWaterHero(toAdd int64) {
	t.waterHeroIds = append(t.waterHeroIds, toAdd)
}

func (t *hero_treasury_tree) CollectInfo() (collectTime time.Time, collectSeason shared_proto.Season) {
	return t.collectTime, t.collectSeason
}

func (t *hero_treasury_tree) SetCollectInfo(collectTime time.Time, collectSeason shared_proto.Season) {
	t.collectTime = collectTime
	t.collectSeason = collectSeason
}

func (hero *Hero) CollectTreasuryTreePrize() {
	hero.treasuryTree.collectPrize(hero.Id())
}

func (t *hero_treasury_tree) collectPrize(selfId int64) {
	t.waterTimes = 0
	t.SetCollectInfo(time.Time{}, 0)

	// 移除自己今日浇水状态，领完可以立马浇水一次
	t.waterHeroIds = i64.RemoveIfPresent(t.waterHeroIds, selfId)
}

func (t *hero_treasury_tree) AddHelpMeInfo(helpMeHeroId int64, season shared_proto.Season, maxLogCount uint64) {
	if maxLogCount <= 0 {
		return
	}

	// copy on write
	n := uint64(imath.Min(len(t.helpMeHeroIds), len(t.helpMeSeason)))
	oldHelpMeHeroIds := t.helpMeHeroIds
	oldHelpMeSeason := t.helpMeSeason
	if n >= maxLogCount {
		n = maxLogCount - 1
		oldHelpMeHeroIds = oldHelpMeHeroIds[1:]
		oldHelpMeSeason = oldHelpMeSeason[1:]
	}

	// 直接copy出来，设置到最后
	t.helpMeHeroIds = make([]int64, n+1)
	copy(t.helpMeHeroIds, oldHelpMeHeroIds)
	t.helpMeHeroIds[n] = helpMeHeroId

	t.helpMeSeason = make([]shared_proto.Season, n+1)
	copy(t.helpMeSeason, oldHelpMeSeason)
	t.helpMeSeason[n] = season
}

func (t *hero_treasury_tree) HelpMeInfo() (heroIds []int64, seasons []shared_proto.Season) {
	return t.helpMeHeroIds, t.helpMeSeason
}

func (t *hero_treasury_tree) resetDaily() {
	t.waterHeroIds = nil
}

func (t *hero_treasury_tree) unmarshal(proto *server_proto.HeroTreasuryTreeServerProto) {
	if proto == nil {
		return
	}

	t.waterTimes = proto.WaterTimes
	t.collectSeason = proto.CollectSession
	t.collectTime = timeutil.Unix64(proto.CollectTime)
	t.helpMeHeroIds = i64.Copy(proto.HelpMeHeroIds)
	t.waterHeroIds = i64.Copy(proto.WaterHeroIds)
}

func (t *hero_treasury_tree) encode() *server_proto.HeroTreasuryTreeServerProto {
	proto := &server_proto.HeroTreasuryTreeServerProto{}

	proto.WaterTimes = t.waterTimes
	proto.CollectSession = t.collectSeason
	proto.CollectTime = timeutil.Marshal64(t.collectTime)
	proto.HelpMeHeroIds = t.helpMeHeroIds
	proto.WaterHeroIds = t.waterHeroIds

	return proto
}

func (t *hero_treasury_tree) encodeClient() *shared_proto.HeroTreasuryTreeProto {
	proto := &shared_proto.HeroTreasuryTreeProto{}

	proto.WaterTimes = u64.Int32(t.waterTimes)
	proto.CollectSession = t.collectSeason
	proto.CollectTime = timeutil.Marshal32(t.collectTime)
	proto.WaterHeroIds = idbytes.ToBytesArray(t.waterHeroIds)

	return proto
}
