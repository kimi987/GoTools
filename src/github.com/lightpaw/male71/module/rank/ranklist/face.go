package ranklist

import (
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/module/rank/rankface"
)

type RankHolder interface {
	NameQuery(name string) (key int64)     // 通过 名字 查询 key 的方法
	SelfKey(hc iface.HeroController) int64 // 获得自己的key
	RLockFunc(f func(h RLockedRankHolder)) // 在读锁中调用方法，注意别用错了
	LockFunc(f func(h LockedRankHolder))   // 在读锁中调用方法，注意别用错了
}

type LockedRankHolder interface {
	RLockedRankHolder
	AddOrUpdate(objs ...rankface.RankObj) // 添加/更新
	Remove(keys ...int64)                 // 移除
	Walk(f func(list rankface.RankList))
}

type RLockedRankHolder interface {
	RRankListByObj(obj rankface.RankObj) rankface.RRankList // 通过排行榜对象，获得排行榜对象所处的排行榜列表
}

type SingleRankHolder interface {
	RankHolder
	RRankList() rankface.RRankList
	RankObjQuery(key int64) (obj rankface.RankObj) // 通过 key 查询 排行对象 的方法
}

type SubTypeRankHolder interface {
	RankHolder
	RRankList(subType uint64) rankface.RRankList
	RankObjQuery(subType uint64, key int64) (list rankface.RRankList, obj rankface.RankObj) // 通过 key 查询 排行对象 的方法
}
