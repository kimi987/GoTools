package xuanydata

import (
	"testing"
	. "github.com/onsi/gomega"
)

func TestGetRankPrizeDataByRank(t *testing.T) {
	RegisterTestingT(t)

	var datas []*XuanyuanRankPrizeData
	for _, v := range []uint64{1, 2, 6, 11, 51, 101} {
		datas = append(datas, &XuanyuanRankPrizeData{
			Rank: v,
		})
	}

	Ω(GetRankPrizeDataByRank(datas, 0)).Should(BeEquivalentTo(datas[5]))

	Ω(GetRankPrizeDataByRank(datas, 1)).Should(BeEquivalentTo(datas[0]))

	Ω(GetRankPrizeDataByRank(datas, 2)).Should(BeEquivalentTo(datas[1]))
	Ω(GetRankPrizeDataByRank(datas, 3)).Should(BeEquivalentTo(datas[1]))
	Ω(GetRankPrizeDataByRank(datas, 5)).Should(BeEquivalentTo(datas[1]))

	Ω(GetRankPrizeDataByRank(datas, 6)).Should(BeEquivalentTo(datas[2]))
	Ω(GetRankPrizeDataByRank(datas, 8)).Should(BeEquivalentTo(datas[2]))
	Ω(GetRankPrizeDataByRank(datas, 10)).Should(BeEquivalentTo(datas[2]))

	Ω(GetRankPrizeDataByRank(datas, 11)).Should(BeEquivalentTo(datas[3]))
	Ω(GetRankPrizeDataByRank(datas, 30)).Should(BeEquivalentTo(datas[3]))
	Ω(GetRankPrizeDataByRank(datas, 50)).Should(BeEquivalentTo(datas[3]))

	Ω(GetRankPrizeDataByRank(datas, 51)).Should(BeEquivalentTo(datas[4]))
	Ω(GetRankPrizeDataByRank(datas, 99)).Should(BeEquivalentTo(datas[4]))
	Ω(GetRankPrizeDataByRank(datas, 100)).Should(BeEquivalentTo(datas[4]))

	Ω(GetRankPrizeDataByRank(datas, 101)).Should(BeEquivalentTo(datas[5]))
	Ω(GetRankPrizeDataByRank(datas, 1000)).Should(BeEquivalentTo(datas[5]))
}
