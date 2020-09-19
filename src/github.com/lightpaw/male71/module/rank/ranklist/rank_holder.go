package ranklist

import (
	"github.com/lightpaw/male7/module/rank/rankface"
	"sync"
)

type rank_holder struct {
	sync.RWMutex
	nameQueryFunc rankface.NameQueryFunc
}

func (h *rank_holder) NameQuery(name string) (key int64) {
	return h.nameQueryFunc(name)
}

type single_rank_holder struct {
	rank_holder
	rankList rankface.RankList
}

func (h *single_rank_holder) RankListByObj(obj rankface.RankObj) rankface.RankList {
	return h.rankList
}

func (h *single_rank_holder) RRankListByObj(obj rankface.RankObj) rankface.RRankList {
	list := h.RankListByObj(obj)
	if list == nil {
		return nil
	}

	return list
}

func (h *single_rank_holder) Walk(f func(list rankface.RankList)) {
	f(h.rankList)
}

func (h *single_rank_holder) RRankList() rankface.RRankList {
	return h.rankList
}

func (h *single_rank_holder) RankObjQuery(key int64) (obj rankface.RankObj) {
	return h.rankList.GetRankObj(key)
}

func (h *single_rank_holder) AddOrUpdate(objs ...rankface.RankObj) {
	for _, obj := range objs {
		h.rankList.AddOrUpdate(obj)
	}
}

func (h *single_rank_holder) Remove(keys ...int64) {
	for _, key := range keys {
		h.rankList.Remove(key)
	}
}

func (h *single_rank_holder) RLockFunc(f func(h RLockedRankHolder)) {
	h.RLock()
	defer h.RUnlock()
	f(h)
}

func (h *single_rank_holder) LockFunc(f func(h LockedRankHolder)) {
	h.RLock()
	defer h.RUnlock()
	f(h)
}

type sub_type_rank_holder struct {
	rank_holder

	getRankListBySubTypeFunc func(subType uint64) rankface.RankList
	getRankListByObjFunc     func(obj rankface.RankObj) rankface.RankList
	rangeRankListFunc        func(f func(list rankface.RankList))
	addOrUpdateFunc          func(objs ...rankface.RankObj)
	removeFunc               func(keys ...int64)
}

func (h *sub_type_rank_holder) RankListByObj(obj rankface.RankObj) rankface.RankList {
	return h.getRankListByObjFunc(obj)
}

func (h *sub_type_rank_holder) RRankListByObj(obj rankface.RankObj) rankface.RRankList {
	list := h.RankListByObj(obj)
	if list == nil {
		return nil
	}

	return list
}

func (h *sub_type_rank_holder) Walk(f func(list rankface.RankList)) {
	h.rangeRankListFunc(f)
}

func (h *sub_type_rank_holder) RRankList(subType uint64) rankface.RRankList {
	return h.getRankListBySubTypeFunc(subType)
}

func (h *sub_type_rank_holder) RankObjQuery(subType uint64, key int64) (rankface.RRankList, rankface.RankObj) {
	if l := h.RRankList(subType); l != nil {
		return l, l.GetRankObj(key)
	}
	return nil, nil
}

func (h *sub_type_rank_holder) AddOrUpdate(objs ...rankface.RankObj) {
	h.addOrUpdateFunc(objs...)
}

func (h *sub_type_rank_holder) Remove(keys ...int64) {
	h.removeFunc(keys...)
}

func (h *sub_type_rank_holder) RLockFunc(f func(h RLockedRankHolder)) {
	h.RLock()
	defer h.RUnlock()
	f(h)
}

func (h *sub_type_rank_holder) LockFunc(f func(h LockedRankHolder)) {
	h.RLock()
	defer h.RUnlock()
	f(h)
}
