package xuanydata

import (
	"github.com/lightpaw/male7/util/check"
	"math"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/util/u64"
)

//gogen:config
type XuanyuanRangeData struct {
	_ struct{} `file:"轩辕会武/积分区间.txt"`
	_ struct{} `protogen:"true"`

	Id uint64

	// 排名区间低
	LowRank int `validator:"int" desc:"排名区间低"`

	// 排名区间高
	HighRank int `validator:"int" desc:"排名区间高"`

	// 挑战获胜积分
	WinScore uint64 `desc:"挑战获胜积分"`

	// 挑战失败积分
	LoseScore uint64 `validator:"int" desc:"挑战失败积分"`

	// 防守失败（被挑战成功），损失积分
	DefenseLostScore uint64 `validator:"int" desc:"防守失败（被挑战成功），损失积分"`

	// 战斗场景
	CombatScene *scene.CombatScene `protofield:",%s.Id,string" desc:"战斗场景id"`
}

func (d *XuanyuanRangeData) Init(filename string, dataMap map[uint64]*XuanyuanRangeData) {

	check.PanicNotTrue(d.LowRank < d.HighRank, "%s Id=%d 配置必须满足 low_rank < high_rank", filename, d.Id)
	if d.Id > 1 {
		prev := dataMap[d.Id-1]
		check.PanicNotTrue(prev != nil, "%s，Id必须从1开始连续配置", filename)

		check.PanicNotTrue(prev.LowRank > d.HighRank, "%s，前一级的low_rank > 后一级的high_rank", filename)
	}
}

//gogen:config
type XuanyuanRankPrizeData struct {
	_ struct{} `file:"轩辕会武/排名奖励.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"base.proto"`

	Id uint64

	Rank uint64 `desc:"排名"`

	Prize *resdata.Prize `desc:"排名奖励" head:"-"`
	PlunderPrize *resdata.PlunderPrize `protofield:"-"`
	ShowPrize *resdata.Prize `desc:"展示奖励"`
}

func (d *XuanyuanRankPrizeData) Init(filename string, dataMap map[uint64]*XuanyuanRankPrizeData, configs interface {
	GetGuildLevelPrizeArray() []*resdata.GuildLevelPrize
}) {

	if d.Id > 1 {
		prev := dataMap[d.Id-1]
		check.PanicNotTrue(prev != nil, "%s，Id必须从1开始连续配置", filename)

		check.PanicNotTrue(prev.Rank < d.Rank, "%s，rank字段必须从小到大配置", filename)
	}

	d.Prize = d.PlunderPrize.Prize
}

func GetRankPrizeDataByRank(datas []*XuanyuanRankPrizeData, rank int) *XuanyuanRankPrizeData {
	if rank <= 0 {
		return datas[len(datas)-1]
	}

	u64rank := u64.FromInt(rank)
	n := len(datas)
	for i := n - 1; i >= 0; i-- {
		d := datas[i]
		if u64rank >= d.Rank {
			return d
		}
	}

	return datas[0]
}

//gogen:config
type XuanyuanMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"轩辕会武/杂项.txt"`
	_ struct{} `protogen:"true"`

	// 每日挑战次数
	ChallengeTimesLimit uint64 `desc:"每日挑战次数上限"`

	// 每日最少损失积分
	DailyMaxLostScore uint64 `protofield:"-"`

	RankCount uint64 `default:"10000" protofield:"-"`

	RecordBatchCount uint64 `default:"10" protofield:"-"`

	InitScore uint64 `default:"200" protofield:"-"`

	posRangeMap map[int]*XuanyuanRangeData
}

func (d *XuanyuanMiscData) Init(filename string, configs interface {
	GetXuanyuanRangeDataArray() []*XuanyuanRangeData
}) {

	rangeMap := make(map[int]*XuanyuanRangeData)
	min := math.MaxInt64
	max := math.MinInt64
	for _, v := range configs.GetXuanyuanRangeDataArray() {
		min = imath.Min(min, v.LowRank)
		max = imath.Max(max, v.HighRank)

		for i := v.LowRank; i <= v.HighRank; i++ {
			check.PanicNotTrue(rangeMap[i] == nil, "轩辕会武积分区间配置存在重叠，重叠排名 %d", i)
			rangeMap[i] = v
		}
	}

	for i := min; i < max; i++ {
		check.PanicNotTrue(i == 0 || rangeMap[i] != nil, "轩辕会武积分区间配置必须连续，非连续排名 %d", i)
	}

	d.posRangeMap = rangeMap
}

func (d *XuanyuanMiscData) GetRangeByDiff(diff int) *XuanyuanRangeData {
	return d.posRangeMap[diff]
}
