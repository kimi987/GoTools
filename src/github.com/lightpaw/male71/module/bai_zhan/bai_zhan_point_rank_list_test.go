package bai_zhan

import (
	"github.com/lightpaw/male7/config/bai_zhan_data"
	"github.com/lightpaw/male7/config/rank_data"
	"github.com/lightpaw/male7/gen/ifacemock"
	"github.com/lightpaw/male7/module/bai_zhan/bai_zhan_objs"
	"github.com/lightpaw/male7/pb/shared_proto"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestBai_zhan_point_rank_list_AddOrUpdate(t *testing.T) {
	RegisterTestingT(t)

	ifacemock.ConfigDatas.Mock(ifacemock.ConfigDatas.RankMiscData, func() *rank_data.RankMiscData {
		return &rank_data.RankMiscData{
			MaxRankCount: 10000,
		}
	})

	ifacemock.ConfigDatas.Mock(ifacemock.ConfigDatas.BaiZhanMiscData, func() *bai_zhan_data.BaiZhanMiscData {
		return &bai_zhan_data.BaiZhanMiscData{
			ShowRankCount: 100,
		}
	})

	level1 := &bai_zhan_data.JunXianLevelData{
		Level:             1,
		LevelUpPoint:      33,
		LevelUpPercent:    30,
		LevelDownPoint:    22,
		LevelDownPercent:  30,
		MinKeepLevelCount: 1,
	}
	level2 := &bai_zhan_data.JunXianLevelData{
		Level:             2,
		LevelUpPoint:      33,
		LevelUpPercent:    30,
		LevelDownPoint:    22,
		LevelDownPercent:  30,
		MinKeepLevelCount: 1,
	}
	level3 := &bai_zhan_data.JunXianLevelData{
		Level:             3,
		LevelUpPoint:      33,
		LevelUpPercent:    30,
		LevelDownPoint:    22,
		LevelDownPercent:  30,
		MinKeepLevelCount: 1,
	}
	level1.NextLevel = level2
	level2.PrevLevel = level1
	level2.NextLevel = level3
	level3.PrevLevel = level2

	ifacemock.ConfigDatas.Mock(ifacemock.ConfigDatas.GetJunXianLevelDataArray, func() []*bai_zhan_data.JunXianLevelData {
		return []*bai_zhan_data.JunXianLevelData{
			level1,
			level2,
			level3,
		}
	})

	lists := NewBaiZhanPointRankLists(ifacemock.ConfigDatas)
	list := lists[level2.Level-1]

	obj1 := bai_zhan_objs.NewBaiZhanObj(1, ifacemock.HeroSnapshotService.Get, 2, level2, time.Now())
	obj1.AddPoint(10, time.Now())
	list.AddOrUpdate(obj1)

	Ω(obj1.Rank()).Should(BeEquivalentTo(1))

	obj2 := bai_zhan_objs.NewBaiZhanObj(2, ifacemock.HeroSnapshotService.Get, 2, level2, time.Now())
	obj2.AddPoint(11, time.Now())
	list.AddOrUpdate(obj2)
	Ω(obj1.Rank()).Should(BeEquivalentTo(2))
	Ω(obj2.Rank()).Should(BeEquivalentTo(1))

	obj2.AddPoint(11, time.Now())
	list.AddOrUpdate(obj2)
	Ω(obj1.Rank()).Should(BeEquivalentTo(2))
	Ω(obj2.Rank()).Should(BeEquivalentTo(1))

	obj1.AddPoint(22, time.Now())
	list.AddOrUpdate(obj1)
	Ω(obj1.Rank()).Should(BeEquivalentTo(1))
	Ω(obj2.Rank()).Should(BeEquivalentTo(2))

	obj3 := bai_zhan_objs.NewBaiZhanObj(3, ifacemock.HeroSnapshotService.Get, 2, level2, time.Now())
	obj3.AddPoint(23, time.Now())
	list.AddOrUpdate(obj3)
	Ω(obj1.Rank()).Should(BeEquivalentTo(1))
	Ω(obj2.Rank()).Should(BeEquivalentTo(3))
	Ω(obj3.Rank()).Should(BeEquivalentTo(2))

	obj4 := bai_zhan_objs.NewBaiZhanObj(4, ifacemock.HeroSnapshotService.Get, 2, level2, time.Now())
	obj4.AddPoint(40, time.Now())
	list.AddOrUpdate(obj4)
	Ω(obj1.Rank()).Should(BeEquivalentTo(2))
	Ω(obj2.Rank()).Should(BeEquivalentTo(4))
	Ω(obj3.Rank()).Should(BeEquivalentTo(3))
	Ω(obj4.Rank()).Should(BeEquivalentTo(1))

	levelUpMaxRank, levelUpNeedMinPoint, levelDownMinRank, levelKeepNeedPoint := list.LevelUpAndDownRankAndPoint()
	t.Log(levelUpMaxRank, levelUpNeedMinPoint, levelDownMinRank, levelKeepNeedPoint)

	Ω(obj4.LevelChangeType(levelUpMaxRank, levelDownMinRank)).Should(BeEquivalentTo(shared_proto.LevelChangeType_LEVEL_UP))
	Ω(obj1.LevelChangeType(levelUpMaxRank, levelDownMinRank)).Should(BeEquivalentTo(shared_proto.LevelChangeType_LEVEL_KEEP))
	Ω(obj3.LevelChangeType(levelUpMaxRank, levelDownMinRank)).Should(BeEquivalentTo(shared_proto.LevelChangeType_LEVEL_KEEP))
	Ω(obj2.LevelChangeType(levelUpMaxRank, levelDownMinRank)).Should(BeEquivalentTo(shared_proto.LevelChangeType_LEVEL_DOWN))

	t.Log(obj1.Point(), obj1.LevelChangeType(levelUpMaxRank, levelDownMinRank))
	t.Log(obj2.Point(), obj2.LevelChangeType(levelUpMaxRank, levelDownMinRank))
	t.Log(obj3.Point(), obj3.LevelChangeType(levelUpMaxRank, levelDownMinRank))
	t.Log(obj4.Point(), obj4.LevelChangeType(levelUpMaxRank, levelDownMinRank))
}
