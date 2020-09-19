package sharedguilddata

import "testing"
import (
	. "github.com/onsi/gomega"
	"sort"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/entity/daily_amount"
	"fmt"
	"time"
)

func TestUtil(t *testing.T) {
	RegisterTestingT(t)

	g1 := &Guild{id: 1}
	g2 := &Guild{id: 2}
	g3 := &Guild{id: 3}
	g4 := &Guild{id: 4}
	g5 := &Guild{id: 5}
	g6 := &Guild{id: 6}
	g7 := &Guild{id: 7}

	array := []*Guild{g1, g3, g6, g2, g5, g4, g7}
	sort.Sort(idRankSlice(array))
	Ω(array).Should(Equal([]*Guild{g1, g2, g3, g4, g5, g6, g7}))

	Ω(GetLast(array)).Should(Equal(g7))

	array = RemoveAndLeftShift(array, 0)
	Ω(array).Should(Equal([]*Guild{g1, g2, g3, g4, g5, g6, g7}))

	array = RemoveAndLeftShift(array, 1)
	Ω(array).Should(Equal([]*Guild{g2, g3, g4, g5, g6, g7}))

	array = RemoveAndLeftShift(array, 8)
	Ω(array).Should(Equal([]*Guild{g2, g3, g4, g5, g6, g7}))

	array = RemoveAndLeftShift(array, 7)
	Ω(array).Should(Equal([]*Guild{g2, g3, g4, g5, g6}))
}

func TestGuilds(t *testing.T) {
	RegisterTestingT(t)

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	if err != nil {
		return
	}
	guilds := NewGuilds(datas)

	g1 := &Guild{id: 1}
	g2 := &Guild{id: 2}
	g3 := &Guild{id: 3}
	g4 := &Guild{id: 4}
	g5 := &Guild{id: 5}
	g6 := &Guild{id: 6}
	g7 := &Guild{id: 7}

	Ω(guilds.Get(1)).Should(BeNil())
	Ω(guilds.IdRankArray()).Should(BeEmpty())

	// 添加 g1
	guilds.Add(g1)
	Ω(guilds.Get(1)).Should(Equal(g1))
	Ω(guilds.IdRankArray()).Should(Equal([]*Guild{g1}))

	// 重复添加
	fakeG1 := &Guild{id: 1}
	guilds.Add(fakeG1)
	Ω(guilds.Get(1)).Should(Equal(g1))

	// 删掉 g1
	Ω(guilds.Remove(1)).Should(Equal(g1))
	Ω(guilds.Get(1)).Should(BeNil())

	// 添加 g2 g4 g3
	guilds.Add(g2)
	Ω(guilds.Get(2)).Should(Equal(g2))
	Ω(guilds.IdRankArray()).Should(Equal([]*Guild{g2}))

	guilds.Add(g4)
	Ω(guilds.Get(4)).Should(Equal(g4))
	Ω(guilds.IdRankArray()).Should(Equal([]*Guild{g2, g4}))

	guilds.Add(g3)
	Ω(guilds.Get(3)).Should(Equal(g3))
	Ω(guilds.IdRankArray()).Should(Equal([]*Guild{g2, g3, g4}))

	// 删除 g3
	Ω(guilds.Remove(3)).Should(Equal(g3))
	Ω(guilds.Get(3)).Should(BeNil())
	Ω(guilds.IdRankArray()).Should(Equal([]*Guild{g2, g4}))

	// 添加 g5 g7 g6
	guilds.Add(g5)
	Ω(guilds.Get(5)).Should(Equal(g5))
	Ω(guilds.IdRankArray()).Should(Equal([]*Guild{g2, g4, g5}))

	guilds.Add(g7)
	Ω(guilds.Get(7)).Should(Equal(g7))
	Ω(guilds.IdRankArray()).Should(Equal([]*Guild{g2, g4, g5, g7}))

	guilds.Add(g6)
	Ω(guilds.Get(6)).Should(Equal(g6))
	Ω(guilds.IdRankArray()).Should(Equal([]*Guild{g2, g4, g5, g6, g7}))

	// 删除 不存在的
	Ω(guilds.Remove(8)).Should(BeNil())
	Ω(guilds.IdRankArray()).Should(Equal([]*Guild{g2, g4, g5, g6, g7}))

	// 全部删掉
	array := make([]*Guild, len(guilds.IdRankArray()))
	copy(array, guilds.IdRankArray())
	for i := 0; i < len(array); i++ {
		v := array[i]
		Ω(guilds.Remove(v.id)).Should(Equal(v))
	}

	Ω(guilds.IdRankArray()).Should(BeEmpty())
}

func TestRefreshPrestigeRank(t *testing.T) {
	RegisterTestingT(t)

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	if err != nil {
		return
	}
	guilds := NewGuilds(datas)

	g1 := &Guild{id: 1}
	g2 := &Guild{id: 2}
	g3 := &Guild{id: 3}

	g1.country = datas.CountryData().MinKeyData
	g2.country = datas.CountryData().MinKeyData
	g3.country = datas.CountryData().MinKeyData

	g1.prestigeDaily = daily_amount.NewDailyAmount(datas.GuildConfig().KeepDailyPrestigeCount)
	g2.prestigeDaily = daily_amount.NewDailyAmount(datas.GuildConfig().KeepDailyPrestigeCount)
	g3.prestigeDaily = daily_amount.NewDailyAmount(datas.GuildConfig().KeepDailyPrestigeCount)

	g1.AddPrestige(1)
	g2.AddPrestige(2)
	g3.AddPrestige(3)

	guilds.Add(g1)
	guilds.Add(g2)
	guilds.Add(g3)

	guilds.Walk(func(g *Guild) {
		g.GmResetDaily(1, time.Now())
	})
	guilds.RefreshPrestigeRank()

	fmt.Println(g1.GetLastPrestigeRank())
	fmt.Println(g2.GetLastPrestigeRank())
	fmt.Println(g3.GetLastPrestigeRank())

	g1.AddPrestige(3)
	g2.AddPrestige(2)
	g3.AddPrestige(1)
	guilds.Walk(func(g *Guild) {
		g.GmResetDaily(1, time.Now())
	})
	guilds.RefreshPrestigeRank()

	fmt.Println(g1.GetLastPrestigeRank())
	fmt.Println(g2.GetLastPrestigeRank())
	fmt.Println(g3.GetLastPrestigeRank())
}
