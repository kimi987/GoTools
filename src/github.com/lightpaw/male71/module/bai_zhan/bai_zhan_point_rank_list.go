package bai_zhan

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/bai_zhan_data"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/bai_zhan"
	"github.com/lightpaw/male7/module/bai_zhan/bai_zhan_objs"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"sort"
	"sync"
)

type RankChangeType int32

const (
	NoChange RankChangeType = iota // 排名没变化
	Up                             // 排名上升
	Down                           // 排名下降
)

func NewBaiZhanPointRankLists(configDatas iface.ConfigDatas) (lists []*bai_zhan_point_rank_list) {
	lists = make([]*bai_zhan_point_rank_list, 0, len(configDatas.GetJunXianLevelDataArray()))

	for _, levelData := range configDatas.GetJunXianLevelDataArray() {
		list := &bai_zhan_point_rank_list{
			levelData:    levelData,
			maxRankCount: configDatas.RankMiscData().MaxRankCount,
		}

		lists = append(lists, list)
	}

	return lists
}

type bai_zhan_point_rank_list struct {
	levelData    *bai_zhan_data.JunXianLevelData
	maxRankCount uint64                          // 最大排行人数
	rankObjArray []*bai_zhan_objs.HeroBaiZhanObj // 排行对象Array
	sync.RWMutex
}

func (r *bai_zhan_point_rank_list) RankCount() uint64 {
	return uint64(len(r.rankObjArray))
}

// 加锁操作
func (r *bai_zhan_point_rank_list) AddOrUpdate(obj *bai_zhan_objs.HeroBaiZhanObj) {
	r.Lock()
	defer r.Unlock()

	var rankChangeType = NoChange
	var oldRank uint64

	var oldObj *bai_zhan_objs.HeroBaiZhanObj
	if obj.Rank() != 0 {
		oldObj = obj
	}

	if oldObj == nil {
		// 不存在
		if r.RankCount() >= r.maxRankCount && !obj.Less(r.rankObjArray[r.RankCount()-1]) {
			// 已经满了且最后一个都不可以干掉
			return
		}

		// 加在尾巴上面
		r.rankObjArray = append(r.rankObjArray, obj)
		oldRank = r.RankCount()
		obj.SetRank(oldRank)

		// 是不是只有一个或者比原来最后一个还要排名靠后
		if r.RankCount() <= 1 || r.rankObjArray[oldRank-2].Less(obj) {
			// 搞啥搞，排名没变化
			return
		}

		// 刚刚就找了有人比我小了
		rankChangeType = Up
	} else {
		oldRank = oldObj.Rank()
		if oldRank <= 0 || oldRank > r.RankCount() || r.rankObjArray[oldObj.Rank()-1] != oldObj {
			// 重排序
			logrus.WithField("oldRank", oldRank).WithField("len", r.RankCount()).Errorln("百战排行榜存在玩家的排名跟数组中的位置不匹配的问题")
			r.resort()

			oldRank = oldObj.Rank()

			if oldRank <= 0 || oldRank > r.RankCount() || r.rankObjArray[oldObj.Rank()-1] != oldObj {
				logrus.WithField("oldRank", oldRank).WithField("len", r.RankCount()).Errorln("百战排行榜存在玩家的排名跟数组中的位置不匹配的问题，重新排序依旧没解决")
				return
			}
		}

		if oldRank > 1 && !r.rankObjArray[oldRank-2].Less(obj) {
			// 上升了
			rankChangeType = Up
		} else if oldRank < r.RankCount() && !obj.Less(r.rankObjArray[oldRank]) {
			// 下降了
			rankChangeType = Down
		} else {
			// 排名没变化
			return
		}
	}

	if rankChangeType == Up {
		// 上升了，前面列表里面一定有比我小的
		array := r.rankObjArray[0: oldRank-1]
		setIndex := sort.Search(len(array), func(i int) bool {
			return obj.Less(array[i])
		})

		var newRank uint64 = uint64(setIndex + 1)

		if newRank != oldRank {
			// 往后copy
			copy(r.rankObjArray[newRank:oldRank], r.rankObjArray[newRank-1:oldRank-1])
		}
		r.rankObjArray[newRank-1] = obj

		// 重新排名
		for k := newRank; k <= oldRank; k++ {
			r.rankObjArray[k-1].SetRank(k)
		}
	} else {
		// 下降了，后面列表里面一定有比我大的
		array := r.rankObjArray[oldRank:]

		setIndex := sort.Search(len(array), func(i int) bool {
			return obj.Less(array[i])
		})

		// 因为本来就要往上移一位，所以不加1
		var newRank uint64 = oldRank + uint64(setIndex)

		if newRank != oldRank {
			// 往前copy
			copy(r.rankObjArray[oldRank-1:newRank-1], r.rankObjArray[oldRank:newRank])
		}
		r.rankObjArray[newRank-1] = obj

		// 重新排名
		for k := oldRank; k <= newRank; k++ {
			r.rankObjArray[k-1].SetRank(k)
		}
	}

	// 检查长度是否超出了
	if r.RankCount() > r.maxRankCount {
		toRemoveObj := r.rankObjArray[r.RankCount()-1]
		r.rankObjArray = r.rankObjArray[:r.RankCount()-1]
		toRemoveObj.SetRank(0)

	}
}

// 重新排序
func (r *bai_zhan_point_rank_list) resort() {
	sort.Sort(hero_bai_zhan_obj_slice(r.rankObjArray))
	for index, obj := range r.rankObjArray {
		obj.SetRank(uint64(index + 1))
	}
}

func (r *bai_zhan_point_rank_list) walk(f func(obj *bai_zhan_objs.HeroBaiZhanObj)) {
	for _, obj := range r.rankObjArray {
		f(obj)
	}
}

// 军衔等级上升的人数
// upNeedPoint 升级需要的积分
func (r *bai_zhan_point_rank_list) levelUpCount() (count uint64, upNeedPoint uint64) {
	if r.levelData.NextLevel == nil {
		return
	}

	upNeedPoint = r.levelData.LevelUpPoint

	rankCount := r.RankCount()
	if rankCount <= 0 {
		// 一个人都没有
		return
	}

	if r.rankObjArray[0].Point() < r.levelData.LevelUpPoint {
		// 第一个就不符合
		return
	}

	maxUpCount := u64.Min(rankCount*r.levelData.LevelUpPercent/100+1, rankCount)
	if maxUpCount <= 0 {
		// 一个可能上升的都没有(人数太少)
		return
	}

	if r.rankObjArray[maxUpCount-1].Point() >= r.levelData.LevelUpPoint {
		// 最后那个也是符合的
		return maxUpCount, r.rankObjArray[maxUpCount-1].Point() + 1
	}

	count = uint64(sort.Search(int(maxUpCount), func(i int) bool {
		return r.rankObjArray[i].Point() <= r.levelData.LevelUpPoint
	}))

	return
}

func (r *bai_zhan_point_rank_list) LevelUpAndDownRankAndPoint() (levelUpMaxRank, levelUpNeedMinPoint, levelDownMinRank, levelKeepNeedPoint uint64) {
	rankCount := r.RankCount()

	levelUpMaxRank, levelUpNeedMinPoint = r.levelUpCount()

	// levelDownMinRank 这个是不准确的
	levelDownMinRank = u64.Sub(rankCount, r.levelData.LevelDownPercent*rankCount/100+1)

	if r.levelData.PrevLevel == nil {
		// 没有上一级
		return
	}

	// 保级的最小的排名，这个是不准确的
	keepRank := u64.Max(levelUpMaxRank+r.levelData.MinKeepLevelCount, levelDownMinRank-1)
	if keepRank > rankCount {
		// 新增一个人也保级，等于不行
		return
	}

	// 有的多，就可能有可以掉级的
	levelKeepNeedPoint = u64.Min(r.rankObjArray[keepRank-1].Point(), r.levelData.LevelDownPoint) + 1

	return
}

func (r *bai_zhan_point_rank_list) RankCache(self bool, startRank, count uint64) (oc pbutil.Buffer, err error) {

	r.RLock()
	defer r.RUnlock()

	rankCount := r.RankCount()

	var levelUpNeedMinPoint = r.levelData.LevelUpPoint
	var levelKeepNeedPoint = uint64(0)

	data := [][]byte{}
	if rankCount > 0 {
		if startRank+count > rankCount {
			// 溢出了
			startRank = u64.Sub(rankCount, count)
		}

		startRank = u64.Max(startRank, 1)

		array := r.rankObjArray[u64.Sub(startRank, 1):u64.Min(rankCount, startRank+count-1)]
		data = make([][]byte, 0, len(array))

		// 获得升级的最低等级
		var levelUpMaxRank, levelDownMinRank uint64

		levelUpMaxRank, levelUpNeedMinPoint, levelDownMinRank, levelKeepNeedPoint = r.LevelUpAndDownRankAndPoint()
		for _, obj := range array {
			proto := obj.Encode4Rank(levelUpMaxRank, levelDownMinRank)
			if proto == nil {
				continue
			}

			data = append(data, must.Marshal(proto))
		}
	}

	oc = bai_zhan.NewS2cRequestRankMsg(self, u64.Int32(r.levelData.Level), u64.Int32(startRank), u64.Int32(rankCount), u64.Int32(levelUpNeedMinPoint), u64.Int32(levelKeepNeedPoint), data)

	return
}

type hero_bai_zhan_obj_slice []*bai_zhan_objs.HeroBaiZhanObj

func (p hero_bai_zhan_obj_slice) Len() int           { return len(p) }
func (p hero_bai_zhan_obj_slice) Less(i, j int) bool { return p[i].Less(p[j]) }
func (p hero_bai_zhan_obj_slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
