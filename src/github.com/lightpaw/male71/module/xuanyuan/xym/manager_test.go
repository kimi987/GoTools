package xym

import (
	"testing"
	"math/rand"
	"time"
	. "github.com/onsi/gomega"
)

func TestUpdate(t *testing.T) {
	RegisterTestingT(t)

	x := &RoRankList{}

	mMapFunc := func(m map[int64]uint64) map[int64]*XyHero {
		out := make(map[int64]*XyHero)
		for k, v := range m {
			out[k] = &XyHero{
				heroId: k,
				score:  v,
			}
		}
		return out
	}

	m1 := mMapFunc(map[int64]uint64{
		1: 10,
		2: 30,
		3: 20,
	})

	rankHeroIdFunc := func(r *RoRank) (out []int64) {
		for _, v := range r.rankHeros {
			out = append(out, v.heroId)
		}
		return
	}

	heroScoreMapFunc := func(r *RoRank) map[int64]uint64 {
		m := make(map[int64]uint64)
		for _, v := range r.rankHeros {
			m[v.heroId] = v.rankScore
		}
		return m
	}

	heroRankMapFunc := func(r *RoRank) map[int64]uint64 {
		m := make(map[int64]uint64)
		for _, v := range r.rankHeros {
			m[v.heroId] = uint64(v.rank)
		}
		return m
	}

	updateFunc := func(m map[int64]*XyHero, r *RoRank) {

		for _, v := range m {
			rh := r.GetHero(v.heroId)
			if rh != nil {
				rh.score.Store(v.score)
			}
		}
		return
	}

	now := time.Now()

	// 第一次更新
	x.update(m1, 10000, now, x.Get())

	Ω(rankHeroIdFunc(x.Get())).Should(Equal([]int64{2, 3, 1}))
	Ω(heroScoreMapFunc(x.Get())).Should(Equal(map[int64]uint64{
		1: 10,
		2: 30,
		3: 20,
	}))
	Ω(heroRankMapFunc(x.Get())).Should(Equal(map[int64]uint64{
		1: 3,
		2: 1,
		3: 2,
	}))

	// 同样的数据，结果没有变化
	x.update(m1, 5, now, x.Get())
	Ω(rankHeroIdFunc(x.Get())).Should(Equal([]int64{2, 3, 1}))
	Ω(heroScoreMapFunc(x.Get())).Should(Equal(map[int64]uint64{
		1: 10,
		2: 30,
		3: 20,
	}))
	Ω(heroRankMapFunc(x.Get())).Should(Equal(map[int64]uint64{
		1: 3,
		2: 1,
		3: 2,
	}))

	// 第二次更新
	m2 := mMapFunc(map[int64]uint64{
		1: 25,
		4: 40,
	})
	updateFunc(m2, x.Get())
	x.update(m2, 5, now, x.Get())

	Ω(rankHeroIdFunc(x.Get())).Should(Equal([]int64{4, 2, 1, 3}))
	Ω(heroScoreMapFunc(x.Get())).Should(Equal(map[int64]uint64{
		1: 25,
		2: 30,
		3: 20,
		4: 40,
	}))
	Ω(heroRankMapFunc(x.Get())).Should(Equal(map[int64]uint64{
		1: 3,
		2: 2,
		3: 4,
		4: 1,
	}))

	// 第三次更新
	m3 := mMapFunc(map[int64]uint64{
		5: 5,
		6: 60,
		7: 35,
	})
	updateFunc(m3, x.Get())
	x.update(m3, 5, now, x.Get())

	Ω(rankHeroIdFunc(x.Get())).Should(Equal([]int64{6, 4, 7, 2, 1}))
	Ω(heroScoreMapFunc(x.Get())).Should(Equal(map[int64]uint64{
		1: 25,
		2: 30,
		4: 40,
		6: 60,
		7: 35,
	}))
	Ω(heroRankMapFunc(x.Get())).Should(Equal(map[int64]uint64{
		1: 5,
		2: 4,
		4: 2,
		6: 1,
		7: 3,
	}))
}

func BenchmarkUpdate(b *testing.B) {

	x := &RoRankList{}

	maps := randTestData(10, 10000)
	n := len(maps)

	updateFunc := func(m map[int64]*XyHero, r *RoRank) {
		if r == nil {
			return
		}

		for _, v := range m {
			rh := r.GetHero(v.heroId)
			if rh != nil {
				rh.score.Store(v.score)
			}
		}
		return
	}

	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := maps[i%n]
		prev := x.Get()
		updateFunc(m, prev)
		x.update(m, 10000, now, prev)
	}
}

//func BenchmarkUpdate1(b *testing.B) {
//
//	x := &RoRankList{}
//
//	maps := randTestData(10, 10000)
//	n := len(maps)
//
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		x.Update1(maps[i%n], 10000)
//	}
//}

var seed = time.Now().UnixNano()

func randTestData(mapCount, countPerMap int) []map[int64]*XyHero {
	rand.Seed(seed)

	mMapFunc := func(m map[int64]uint64) map[int64]*XyHero {
		out := make(map[int64]*XyHero)
		for k, v := range m {
			out[k] = &XyHero{
				heroId: k,
				score:  v,
			}
		}
		return out
	}

	var mm []map[int64]*XyHero
	for i := 0; i < mapCount; i++ {
		m := make(map[int64]uint64)
		for i := 0; i < countPerMap; i++ {
			m[rand.Int63()] = rand.Uint64()
		}

		mm = append(mm, mMapFunc(m))
	}

	return mm
}
