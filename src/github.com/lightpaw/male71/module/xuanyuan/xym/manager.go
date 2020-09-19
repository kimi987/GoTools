package xym

import (
	"sync"
	"github.com/lightpaw/male7/pb/shared_proto"
	"sync/atomic"
	"sort"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/sortkeys"
	atomic2 "github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/male7/gen/pb/xuanyuan"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/pb/server_proto"
	"time"
	"github.com/lightpaw/male7/util/timeutil"
)

func NewManager(rankCount uint64) *XuanyuanManager {
	return &XuanyuanManager{
		rankCount:     rankCount,
		challengerMap: make(map[int64]*XyHero),
	}
}

type XuanyuanManager struct {
	sync.RWMutex

	rankCount uint64

	// 排行榜数据，在排行榜中的玩家都在上面
	rrl RoRankList

	// 挑战者数据
	challengerMap map[int64]*XyHero
}

func (m *XuanyuanManager) Encode(proto *server_proto.XuanyuanModuleProto) {
	r := m.Get()
	if r != nil {
		proto.UpdateTime = timeutil.Marshal64(r.updateTime)

		proto.RankHero = make([]*server_proto.XuanyuanRankHeroProto, 0, len(r.rankHeros))
		for _, v := range r.rankHeros {
			_, mirror := v.GetMirror()
			proto.RankHero = append(proto.RankHero, &server_proto.XuanyuanRankHeroProto{
				HeroId:    v.heroId,
				Score:     v.score.Load(),
				RankScore: v.rankScore,
				Win:       v.win.Load(),
				Lose:      v.lose.Load(),
				Mirror:    mirror,
			})
		}
	}

	m.RLock()
	defer m.RUnlock()
	for _, v := range m.challengerMap {
		proto.Challenger = append(proto.Challenger, &server_proto.XuanyuanRankHeroProto{
			HeroId: v.heroId,
			Score:  v.score,
			Win:    v.win,
			Lose:   v.lose,
			Mirror: v.combatMirror,
		})
	}
}

func (m *XuanyuanManager) Unmarshal(proto *server_proto.XuanyuanModuleProto) {
	if proto == nil {
		return
	}

	n := len(proto.RankHero)
	newRo := &RoRank{
		heroMap:    make(map[int64]*XyRankHero, n),
		rankHeros:  make([]*XyRankHero, 0, n),
		updateTime: timeutil.Unix64(proto.UpdateTime),
	}
	for i, v := range proto.RankHero {
		rank := i + 1
		newHero := newRankHero(v.HeroId,
			u64.FromInt64(int64(v.Score)),
			u64.FromInt64(int64(v.RankScore)),
			v.Win, v.Lose, rank, v.Mirror)
		newRo.heroMap[newHero.heroId] = newHero
		newRo.rankHeros = append(newRo.rankHeros, newHero)
	}
	m.rrl.set(newRo)

	m.Lock()
	defer m.Unlock()
	for _, v := range proto.Challenger {
		m.addChallenger(v.HeroId, v.Score, v.Win, v.Lose, v.Mirror)
	}
}

func (m *XuanyuanManager) Get() *RoRank {
	return m.rrl.Get()
}

func (m *XuanyuanManager) getAndClearChallenger() map[int64]*XyHero {
	m.Lock()
	defer m.Unlock()

	toReturn := m.challengerMap
	m.challengerMap = make(map[int64]*XyHero)
	return toReturn
}

func (m *XuanyuanManager) Update(updateTime time.Time, gmReset bool) bool {

	prev := m.Get()
	if !gmReset && prev != nil && !prev.updateTime.Before(updateTime) {
		return false
	}

	heroMap := m.getAndClearChallenger()
	m.rrl.update(heroMap, int(m.rankCount), updateTime, prev)
	return true
}

func (m *XuanyuanManager) AddChallenger(heroId int64, score, win, lose uint64, player *shared_proto.CombatPlayerProto) {
	m.Lock()
	defer m.Unlock()

	m.addChallenger(heroId, score, win, lose, player)
}

func (m *XuanyuanManager) addChallenger(heroId int64, score, win, lose uint64, player *shared_proto.CombatPlayerProto) {
	m.challengerMap[heroId] = &XyHero{
		heroId:       heroId,
		score:        score,
		win:          win,
		lose:         lose,
		combatMirror: player,
	}
}

type RoRankList struct {
	v atomic.Value
}

func (r *RoRankList) Get() *RoRank {
	if rank := r.v.Load(); rank != nil {
		return rank.(*RoRank)
	}
	return nil
}

func (r *RoRankList) set(toSet *RoRank) {
	r.v.Store(toSet)
}

func (r *RoRankList) update(newHeroMap map[int64]*XyHero, rankCount int, updateTime time.Time, prev *RoRank) {
	// 单线程更新
	//if len(newHeroMap) <= 0 && prev != nil {
	//	// 改个时间
	//	r.set(&RoRank{
	//		heroMap:    prev.heroMap,
	//		rankHeros:  prev.rankHeros,
	//		updateTime: updateTime,
	//	})
	//	return
	//}

	pa := make([]*sortkeys.U64K2V, 0, len(newHeroMap))
	if prev != nil {
		for heroId, v := range prev.heroMap {
			// 如果在榜单中，已榜单为准
			delete(newHeroMap, heroId)

			// 积分 + 战力
			_, m := v.GetMirror()
			var fightAmount uint64
			if m != nil {
				fightAmount = u64.FromInt32(m.TotalFightAmount)
			}

			pa = append(pa, &sortkeys.U64K2V{
				K1: v.score.Load(),
				K2: fightAmount,
				V:  heroId,
			})
		}
	}

	for heroId, v := range newHeroMap {
		var fightAmount uint64
		if v.combatMirror != nil {
			fightAmount = u64.FromInt32(v.combatMirror.TotalFightAmount)
		}

		pa = append(pa, &sortkeys.U64K2V{
			K1: v.score,
			K2: fightAmount,
			V:  heroId,
		})
	}

	sort.Sort(sort.Reverse(sortkeys.U64K2VSlice(pa)))

	n := imath.Min(len(pa), rankCount)
	newRo := &RoRank{
		heroMap:    make(map[int64]*XyRankHero, n),
		rankHeros:  make([]*XyRankHero, 0, n),
		updateTime: updateTime,
	}

	for i := 0; i < n; i++ {
		rank := i + 1

		p := pa[i]
		score := p.K1
		heroId := p.I64Value()

		var newHero *XyRankHero
		if prev != nil {
			prevHero := prev.GetHero(heroId)
			if prevHero != nil {
				newHero = prevHero.copy(score, rank)
			}
		}

		if newHero == nil {
			challenger := newHeroMap[heroId]
			newHero = challenger.newRankHero(rank)
		}

		newRo.heroMap[heroId] = newHero
		newRo.rankHeros = append(newRo.rankHeros, newHero)
	}

	r.set(newRo)
}

type RoRank struct {
	heroMap map[int64]*XyRankHero

	rankHeros []*XyRankHero

	updateTime time.Time
}

func (m *RoRank) GetUpdateTime() time.Time {
	return m.updateTime
}

func (m *RoRank) RankCount() int {
	return len(m.rankHeros)
}

func (m *RoRank) GetHero(heroId int64) *XyRankHero {
	return m.heroMap[heroId]
}

func (m *RoRank) GetHeroByRank(rank int) *XyRankHero {
	if rank > 0 && rank <= len(m.rankHeros) {
		return m.rankHeros[rank-1]
	}
	return nil
}

func (m *RoRank) Range(f func(hero *XyRankHero) (toContinue bool)) {
	for _, v := range m.rankHeros {
		if !f(v) {
			break
		}
	}
}

func newRankHero(heroId int64, score, rankScore, win, lose uint64, rank int, combatMirror *shared_proto.CombatPlayerProto) *XyRankHero {
	newHero := &XyRankHero{
		heroId:          heroId,
		score:           atomic2.NewUint64(score),
		rankScore:       rankScore,
		rank:            rank,
		win:             atomic2.NewUint64(win),
		lose:            atomic2.NewUint64(lose),
		combatMirrorRef: &atomic.Value{},
	}

	newHero.SetMirror(combatMirror, int64(rank))

	return newHero
}

type XyRankHero struct {
	// 玩家id
	heroId int64

	// 当前积分
	score *atomic2.Uint64

	// 排名积分
	rankScore uint64

	// 名次
	rank int

	// 胜利次数
	win *atomic2.Uint64

	// 失败次数
	lose *atomic2.Uint64

	// 挑战镜像
	combatMirrorRef *atomic.Value

	targetBytesCache atomic.Value
}

func (hero *XyRankHero) copy(rankScore uint64, rank int) *XyRankHero {
	newHero := &XyRankHero{
		heroId:          hero.heroId,
		score:           hero.score,
		rankScore:       rankScore,
		rank:            rank,
		win:             hero.win,
		lose:            hero.lose,
		combatMirrorRef: hero.combatMirrorRef,
	}

	return newHero
}

func (hero *XyRankHero) Id() int64 {
	return hero.heroId
}

func (hero *XyRankHero) Rank() int {
	return hero.rank
}

func (hero *XyRankHero) GetScore() uint64 {
	return hero.score.Load()
}

func (hero *XyRankHero) SetScore(toSet uint64) {
	hero.score.Store(toSet)
}

func (hero *XyRankHero) GetWin() uint64 {
	return hero.win.Load()
}

func (hero *XyRankHero) IncWin() uint64 {
	amt := hero.win.Inc()
	hero.clearTargetBytesCache()
	return amt
}

func (hero *XyRankHero) GetLose() uint64 {
	return hero.lose.Load()
}

func (hero *XyRankHero) IncLose() uint64 {
	amt := hero.lose.Inc()
	hero.clearTargetBytesCache()
	return amt
}

func (hero *XyRankHero) EncodeTarget(getter func(int64) *snapshotdata.HeroSnapshot) []byte {
	cache := hero.targetBytesCache.Load()
	if cache != nil {
		if b, ok := cache.([]byte); ok && len(b) > 0 {
			return b
		}
	}

	proto := hero.encodeTarget(getter)
	protoBytes := must.Marshal(proto)
	hero.targetBytesCache.Store(protoBytes)
	return protoBytes
}

var emptyBytes = make([]byte, 0)

func (hero *XyRankHero) clearTargetBytesCache() {
	hero.targetBytesCache.Store(emptyBytes)
}

func (hero *XyRankHero) encodeTarget(getter func(int64) *snapshotdata.HeroSnapshot) *shared_proto.XuanyuanTargetProto {
	proto := &shared_proto.XuanyuanTargetProto{}

	heroSnapshot := getter(hero.Id())
	if heroSnapshot != nil {
		proto.Hero = heroSnapshot.EncodeBasic4Client()
	} else {
		proto.Hero = idbytes.HeroBasicProto(hero.Id())
	}
	proto.Win = u64.Int32(hero.GetWin())
	proto.Lose = u64.Int32(hero.GetLose())
	proto.Score = u64.Int32(hero.rankScore)

	ref := hero.getMirrorRef()
	proto.FightAmount = ref.combatMirror.TotalFightAmount

	return proto
}

func (hero *XyRankHero) getMirrorRef() *combatMirrorWithVersion {
	return hero.combatMirrorRef.Load().(*combatMirrorWithVersion)
}

func (hero *XyRankHero) GetMirror() (int64, *shared_proto.CombatPlayerProto) {
	ref := hero.getMirrorRef()
	return ref.version, ref.combatMirror
}

func (hero *XyRankHero) SetMirror(toSet *shared_proto.CombatPlayerProto, version int64) int64 {
	newMirror := newCombatMirror(toSet, version)
	hero.combatMirrorRef.Store(newMirror)
	hero.clearTargetBytesCache()
	return newMirror.version
}

func (hero *XyRankHero) GetQueryTargetTroopMsg() pbutil.Buffer {
	return hero.getMirrorRef().getQueryTroopMsg(hero.Id())
}

func newCombatMirror(combatMirror *shared_proto.CombatPlayerProto, version int64) *combatMirrorWithVersion {
	return &combatMirrorWithVersion{
		version:      version,
		combatMirror: combatMirror,
	}
}

type combatMirrorWithVersion struct {
	version      int64
	combatMirror *shared_proto.CombatPlayerProto

	queryTroopMsgCache atomic.Value
}

func (c *combatMirrorWithVersion) getQueryTroopMsg(heroId int64) pbutil.Buffer {
	msgRef := c.queryTroopMsgCache.Load()
	if msgRef != nil {
		return msgRef.(pbutil.Buffer)
	}

	msg := xuanyuan.NewS2cQueryTargetTroopMsg(idbytes.ToBytes(heroId), int32(c.version), must.Marshal(c.combatMirror)).Static()
	c.queryTroopMsgCache.Store(msg)
	return msg
}

type XyHero struct {
	// 玩家id
	heroId int64

	// 最新积分
	score uint64

	// 胜利次数
	win uint64

	// 失败次数
	lose uint64

	// 挑战镜像
	combatMirror *shared_proto.CombatPlayerProto
}

func (hero *XyHero) newRankHero(rank int) *XyRankHero {
	return newRankHero(hero.heroId, hero.score, hero.score, hero.win, hero.lose,
		rank, hero.combatMirror)
}
