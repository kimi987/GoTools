package mingc_war

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"sync/atomic"
	"github.com/lightpaw/pbutil"
	"sort"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
)

func newMcWarTroopRankObject(heroId int64, hero *shared_proto.HeroBasicProto, atk bool) *McWarTroopRankObject {
	p := &McWarTroopRankObject{}
	p.heroId = heroId
	p.heroProto = hero
	p.isAtk = atk
	return p
}

type McWarTroopRankObject struct {
	McWarRankData
	multiKill uint64 // 连斩
}

func (o *McWarTroopRankObject) copyData() *McWarRankData {
	return &McWarRankData{
		heroProto:    o.heroProto,
		wounded:      o.wounded,
		kill:         o.kill,
		destroy:      o.destroy,
		winTimes:     o.winTimes,
		loseTimes:    o.loseTimes,
		drumTimes:    o.drumTimes,
		maxMultiKill: o.maxMultiKill,
		heroId:       o.heroId,
		isAtk:        o.isAtk,
	}
}

func (o *McWarTroopRankObject) refresh4FightResult(kill, wounded uint64, win bool) (multiKillRefreshed bool) {
	o.kill += kill
	o.wounded += wounded
	if win {
		o.winTimes++
		o.multiKill++
		multiKillRefreshed = o.multiKill > 1
		if multiKillRefreshed && o.multiKill > o.maxMultiKill {
			o.maxMultiKill = o.multiKill
		}
	} else {
		o.loseTimes++
		o.multiKill = 0
	}
	return
}

func (o *McWarTroopRankObject) addDestroy(destroy uint64) {
	o.destroy += destroy
}

func (o *McWarTroopRankObject) increaseDrumTimes() {
	o.drumTimes++
}

func (o *McWarTroopRankObject) resetMultiKill() {
	o.multiKill = 0
}

func newMcWarTroopsRank() *McWarTroopsRank {
	p := &McWarTroopsRank{
		troopRankMap: make(map[int64]*McWarTroopRankObject),
		rankRef:      &atomic.Value{},
	}
	p.rankRef.Store(&RankRef{})
	return p
}

type RankRef struct {
	dataMap map[int64]*McWarRankData // 所有的部队排名<heroId, McWarRankData>

	rankMsg     pbutil.Buffer // 前100名的数据静态消息
	sortVersion uint64
}

type McWarTroopsRank struct {
	troopRankMap map[int64]*McWarTroopRankObject

	rankRef     *atomic.Value
	rankVersion uint64
	needSort    bool
}

func (r *McWarTroopsRank) loadRef() *RankRef {
	if r.needSort {
		r.needSort = false
		r.sort()
	}
	return r.rankRef.Load().(*RankRef)
}

func (r *McWarTroopsRank) getRankObj(heroId int64) *McWarTroopRankObject {
	if o, ok := r.troopRankMap[heroId]; ok {
		return o
	}
	return nil
}

func (r *McWarTroopsRank) putRankObj(rankOjb *McWarTroopRankObject) {
	r.troopRankMap[rankOjb.heroId] = rankOjb
	if len(r.troopRankMap) <= mc_rank_max_num {
		r.needSort = false
		r.sort()
	}
}

func (r *McWarTroopsRank) sort() {
	r.rankVersion++

	var slice TroopSlice
	for _, rankObj := range r.troopRankMap {
		slice = append(slice, rankObj.copyData())
	}
	sort.Sort(slice)

	dataMap := make(map[int64]*McWarRankData)
	for i, data := range slice {
		data.rank = i + 1
		dataMap[data.heroId] = data
	}

	arrLen := len(slice)
	if arrLen > mc_rank_max_num {
		arrLen = mc_rank_max_num
	}
	infos := make([]*shared_proto.McWarTroopRankProto, arrLen)
	for i := 0; i < arrLen; i++ {
		data := slice[i]
		infos[i] = data.encode4Rank()
	}

	ref := &RankRef{
		dataMap:     dataMap,
		rankMsg:     mingc_war.NewS2cApplyRefreshRankMsg(u64.Int32(r.rankVersion), &shared_proto.McWarTroopsRankProto{infos}).Static(),
		sortVersion: r.rankVersion,
	}
	r.rankRef.Store(ref)
}

type McWarRankData struct {
	heroProto *shared_proto.HeroBasicProto

	wounded      uint64 // 损兵
	kill         uint64 // 歼兵
	destroy      uint64 // 破坏
	winTimes     uint64 // 击败
	loseTimes    uint64 // 被击败
	drumTimes    uint64 // 击鼓数
	maxMultiKill uint64 // 连斩
	heroId       int64

	rank  int
	isAtk bool
}

func (d *McWarRankData) encode4Rank() *shared_proto.McWarTroopRankProto {
	p := &shared_proto.McWarTroopRankProto{}
	p.IsAtk = d.isAtk
	p.Rank = int32(d.rank)
	p.Info = d.encode4Info()
	return p
}

func (d *McWarRankData) encode4Info() *shared_proto.McWarTroopInfoProto {
	return &shared_proto.McWarTroopInfoProto{
		Hero:  d.heroProto,
		Score: d.encode(),
	}
}

func (d *McWarRankData) encode() *shared_proto.McWarTroopScoreProto {
	p := &shared_proto.McWarTroopScoreProto{}
	p.KillAmount = u64.Int32(d.kill)
	p.DestroyAmount = u64.Int32(d.destroy)
	p.WinTimes = u64.Int32(d.winTimes)
	p.MultiKill = u64.Int32(d.maxMultiKill)
	p.DrumTimes = u64.Int32(d.drumTimes)
	p.LoseTimes = u64.Int32(d.loseTimes)
	return p
}

type TroopSlice []*McWarRankData

func (p TroopSlice) Len() int {
	return len(p)
}

func (p TroopSlice) Less(i, j int) bool {
	o1, o2 := p[i], p[j]
	if o1.kill != o2.kill {
		return o1.kill > o2.kill
	}
	if o1.destroy != o2.destroy {
		return o1.destroy > o2.destroy
	}
	if o1.winTimes != o2.winTimes {
		return o1.winTimes > o2.winTimes
	}
	if o1.maxMultiKill != o2.maxMultiKill {
		return o1.maxMultiKill > o2.maxMultiKill
	}
	if o1.drumTimes != o2.drumTimes {
		return o1.drumTimes > o2.drumTimes
	}
	return o1.heroId < o2.heroId
}

func (p TroopSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
