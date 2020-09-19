package ranklist

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestRank(t *testing.T) {
	RegisterTestingT(t)

	snapshotGetter := func(id int64) *snapshotdata.HeroSnapshot {
		return &snapshotdata.HeroSnapshot{Id: id}
	}

	rk := NewRankList(shared_proto.RankType_Tower, 10)

	for i := 1; i <= 10; i++ {
		rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(i), uint64(i), time.Now()))
	}

	for i := 10; i >= 1; i-- {
		Ω(rk.GetRankObj(int64(i)).Rank()).Should(BeEquivalentTo(11 - i))
	}

	for i := 11; i <= 20; i++ {
		rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(i), uint64(i), time.Now()))
		Ω(rk.GetRankObj(int64(i)).Rank()).Should(BeEquivalentTo(1))
	}

	for i := 30; i >= 21; i-- {
		rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(i), uint64(i), time.Now()))
		Ω(rk.GetRankObj(int64(i)).Rank()).Should(BeEquivalentTo(30 - i + 1))
	}

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(31), uint64(22), time.Now()))
	Ω(rk.GetRankObj(int64(31)).Rank()).Should(BeEquivalentTo(10))

	// 25积分有个相同的，但是时间比我早的
	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(31), uint64(25), time.Now()))
	Ω(rk.GetRankObj(int64(31)).Rank()).Should(BeEquivalentTo(7))

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(31), uint64(27), time.Now()))
	Ω(rk.GetRankObj(int64(31)).Rank()).Should(BeEquivalentTo(5))

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(31), uint64(33), time.Now()))
	Ω(rk.GetRankObj(int64(31)).Rank()).Should(BeEquivalentTo(1))

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(31), uint64(26), time.Now()))
	Ω(rk.GetRankObj(int64(31)).Rank()).Should(BeEquivalentTo(6))

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(31), uint64(23), time.Now()))
	Ω(rk.GetRankObj(int64(31)).Rank()).Should(BeEquivalentTo(9))

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(31), uint64(24), time.Now()))
	Ω(rk.GetRankObj(int64(31)).Rank()).Should(BeEquivalentTo(8))

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(31), uint64(22), time.Now()))
	Ω(rk.GetRankObj(int64(31)).Rank()).Should(BeEquivalentTo(10))

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(31), uint64(2), time.Now()))
	Ω(rk.GetRankObj(int64(31)).Rank()).Should(BeEquivalentTo(10))

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(32), uint64(3), time.Now()))
	Ω(rk.GetRankObj(int64(32)).Rank()).Should(BeEquivalentTo(10))

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(31), uint64(2), time.Now()))
	Ω(rk.GetRankObj(int64(31))).Should(BeNil())

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(33), uint64(33), time.Now()))
	Ω(rk.GetRankObj(int64(33)).Rank()).Should(BeEquivalentTo(1))

	rk.AddOrUpdate(NewTowerRankObj(snapshotGetter, int64(32), uint64(3), time.Now()))
	Ω(rk.GetRankObj(int64(32))).Should(BeNil())

	Ω(rk.GetRankObj(int64(30)).Rank()).Should(BeEquivalentTo(2))
	Ω(rk.GetRankObj(int64(29)).Rank()).Should(BeEquivalentTo(3))
	rk.Remove(33)
	Ω(rk.GetRankObj(int64(30)).Rank()).Should(BeEquivalentTo(1))
	Ω(rk.GetRankObj(int64(29)).Rank()).Should(BeEquivalentTo(2))

	Ω(rk.RankType()).Should(BeEquivalentTo(shared_proto.RankType_Tower))
}
