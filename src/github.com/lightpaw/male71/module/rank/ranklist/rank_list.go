package ranklist

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"sort"
)

type RankChangeType int32

const (
	NoChange RankChangeType = iota // 排名没变化
	Up                             // 排名上升
	Down                           // 排名下降
)

func newRankListLstrFuncs() (rankObjQueryFunc rankface.RankObjQueryFunc, addOrUpdateFunc rankface.AddOrUpdateFunc, removeFunc rankface.RemoveFunc) {
	objMap := Newrank_obj_map()

	addOrUpdateFunc = func(oldObj, newObj rankface.RankObj) {
		objMap.Set(newObj.Key(), newObj)
	}

	removeFunc = func(obj rankface.RankObj) {
		objMap.Remove(obj.Key())
	}

	rankObjQueryFunc = func(key int64) rankface.RankObj {
		obj, _ := objMap.Get(key)
		return obj
	}

	return
}

func NewRankList(rankType shared_proto.RankType, maxRankCount uint64) rankface.RankList {
	rankObjQueryFunc, addOrUpdateFunc, removeFunc := newRankListLstrFuncs()
	return NewRankListWithFunc(rankType, maxRankCount, rankObjQueryFunc, addOrUpdateFunc, removeFunc)
}

func NewRankListWithFunc(rankType shared_proto.RankType, maxRankCount uint64, queryFunc rankface.RankObjQueryFunc, addOrUpdateFunc rankface.AddOrUpdateFunc, removeFunc rankface.RemoveFunc) rankface.RankList {
	return newRankListWithFunc(rankType, maxRankCount, queryFunc, addOrUpdateFunc, removeFunc)
}

func newRankListWithFunc(rankType shared_proto.RankType, maxRankCount uint64, queryFunc rankface.RankObjQueryFunc, addOrUpdateFunc rankface.AddOrUpdateFunc, removeFunc rankface.RemoveFunc) *rank_list {
	return &rank_list{
		rankType:         rankType,
		maxRankCount:     maxRankCount,
		rankObjArray:     make([]rankface.RankObj, 0, maxRankCount>>4), // 先初始化一定长度的，比如10000人，默认长度625，只用4次扩展就可以扩展回来
		rankObjQueryFunc: queryFunc,
		addOrUpdateFunc:  addOrUpdateFunc,
		removeFunc:       removeFunc,
	}
}

type rank_list struct {
	rankType         shared_proto.RankType     // 排行榜类型
	rankObjArray     []rankface.RankObj        // 排行对象Array
	maxRankCount     uint64                    // 最大排行人数
	rankObjQueryFunc rankface.RankObjQueryFunc // 排行对象请求
	addOrUpdateFunc  rankface.AddOrUpdateFunc  // 在新增或者更新后的调用
	removeFunc       rankface.RemoveFunc       // 在删除后的调用
}

func (r *rank_list) RankCount() uint64 {
	return uint64(len(r.rankObjArray))
}

func (r *rank_list) RankType() shared_proto.RankType {
	return r.rankType
}

// 加锁操作
func (r *rank_list) AddOrUpdate(obj rankface.RankObj) {
	var rankChangeType = NoChange
	var oldRank uint64

	oldObj := r.GetRankObj(obj.Key())
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
			r.addOrUpdateFunc(oldObj, obj)
			return
		}

		// 刚刚就找了有人比我小了
		rankChangeType = Up
	} else {
		oldRank = oldObj.Rank()
		if oldRank <= 0 || oldRank > r.RankCount() || r.rankObjArray[oldObj.Rank()-1] != oldObj {
			// 重排序
			logrus.WithField("oldRank", oldRank).WithField("len", r.RankCount()).WithField("rankType", r.rankType).Errorln("存在玩家的排名跟数组中的位置不匹配的问题")
			r.resort()

			oldRank = oldObj.Rank()

			if oldRank <= 0 || oldRank > r.RankCount() || r.rankObjArray[oldObj.Rank()-1] != oldObj {
				logrus.WithField("oldRank", oldRank).WithField("len", r.RankCount()).Errorln("存在玩家的排名跟数组中的位置不匹配的问题，重新排序依旧没解决")
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
			// 排名没变化，覆盖掉
			obj.SetRank(oldRank)
			r.rankObjArray[oldRank-1] = obj
			r.addOrUpdateFunc(oldObj, obj)
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

	r.addOrUpdateFunc(oldObj, obj)

	// 检查长度是否超出了
	if r.RankCount() > r.maxRankCount {
		toRemoveObj := r.rankObjArray[r.RankCount()-1]
		r.rankObjArray = r.rankObjArray[:r.RankCount()-1]
		r.removeFunc(toRemoveObj)
	}
}

// 加锁操作
func (r *rank_list) Remove(key int64) {
	obj := r.GetRankObj(key)
	if obj == nil {
		// 不存在
		return
	}

	// 移除掉
	rk := obj.Rank()
	if rk <= 0 || rk > r.RankCount() || r.rankObjArray[obj.Rank()-1] != obj {
		// 重排序
		logrus.WithField("rank", rk).WithField("len", r.RankCount()).Errorln("存在玩家的排名跟数组中的位置不匹配的问题")
		r.resort()

		rk = obj.Rank()

		if rk <= 0 || rk > r.RankCount() || r.rankObjArray[obj.Rank()-1] != obj {
			logrus.WithField("rank", rk).WithField("len", r.RankCount()).Errorln("存在玩家的排名跟数组中的位置不匹配的问题，重新排序依旧没解决")
			return
		}
	}

	r.removeFunc(obj)

	// 干掉吧
	if rk < r.RankCount() {
		copy(r.rankObjArray[rk-1:], r.rankObjArray[rk:])
	}

	r.rankObjArray = r.rankObjArray[:r.RankCount()-1]

	for k := rk; k <= r.RankCount(); k++ {
		r.rankObjArray[k-1].SetRank(k)
	}
}

// 重新排序
func (r *rank_list) resort() {
	sort.Sort(rankface.RankObjSlice(r.rankObjArray))
	for index, obj := range r.rankObjArray {
		obj.SetRank(uint64(index + 1))
	}
}

// 允许多线程操作
func (r *rank_list) GetRankObj(key int64) rankface.RankObj {
	return r.rankObjQueryFunc(key)
}

func (r *rank_list) Walk(f func(obj rankface.RankObj)) {
	for _, obj := range r.rankObjArray {
		f(obj)
	}
}

func (r *rank_list) EncodeClient(startRank, rankCountPerPage uint64) *shared_proto.RankProto {
	rankCount := r.RankCount()

	proto := &shared_proto.RankProto{
		Type:      r.RankType(),
		StartRank: 1,
		RankCount: u64.Int32(rankCount),
	}

	if rankCount > 0 {
		startIndex := u64.Sub(startRank, 1)
		if startIndex+rankCountPerPage > rankCount {
			// 超出最后一页了
			startIndex = u64.Sub(rankCount, rankCountPerPage)
		}

		for _, obj := range r.rankObjArray[startIndex:u64.Min(startIndex+rankCountPerPage, rankCount)] {
			obj.EncodeClient(proto)
		}

		proto.StartRank = u64.Int32(startIndex + 1)
	}

	return proto
}

func (r *rank_list) RankKeys(startCount, rankCountPerPage uint64) []int64 {
	rankCount := r.RankCount()

	if rankCount < startCount {
		return nil
	}

	startIndex := u64.Sub(startCount, 1)
	endIndex := u64.Min(startIndex+rankCountPerPage, rankCount)
	n := u64.Sub(endIndex, startIndex)

	array := make([]int64, n)
	for i := uint64(0); i < n; i++ {
		array[i] = r.rankObjArray[startIndex+i].Key()
	}

	return array
}

func (r *rank_list) Range(f func(o rankface.RankObj) bool) {
	for _, o := range r.rankObjArray {
		if !f(o) {
			return
		}
	}
}

func (r *rank_list) EncodeServer(proto *server_proto.RankServerProto) {
	for _, obj := range r.rankObjArray {
		obj.EncodeServer(proto)
	}
}
