package rankface

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
)

type RankList interface {
	RRankList                // 只可读排行榜
	AddOrUpdate(obj RankObj) // 添加/更新
	Remove(key int64)        // 移除
	Walk(func(obj RankObj))
	EncodeServer(proto *server_proto.RankServerProto) // 序列化给服务器
}

type RRankListFunc func(rankList RRankList)

// 只可读排行榜
type RRankList interface {
	RankType() shared_proto.RankType                                         // 排行榜类型
	RankCount() uint64                                                       // 排行榜人数
	GetRankObj(key int64) RankObj                                            // 请求排行对象
	EncodeClient(startRank, rankCountPerPage uint64) *shared_proto.RankProto // 把从某个排名开始的数据序列化
	RankKeys(startCount, rankCountPerPage uint64) []int64
	Range(func(obj RankObj) bool)
}

type RankObj interface {
	Key() int64                                       // key
	RankType() shared_proto.RankType                  // 排行榜类型
	Rank() uint64                                     // 获得排行
	SetRank(toSet uint64)                             // 设置排行
	Less(obj RankObj) bool                            // 是不是排在obj前面
	EncodeClient(proto *shared_proto.RankProto)       // 序列化
	EncodeServer(proto *server_proto.RankServerProto) // 序列化给服务器
	EncodeHeroSnapshotProto() (proto *shared_proto.HeroBasicSnapshotProto)
}

type RankObjSlice []RankObj

func (p RankObjSlice) Len() int           { return len(p) }
func (p RankObjSlice) Less(i, j int) bool { return p[i].Less(p[j]) }
func (p RankObjSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type NameQueryFunc func(name string) (key int64)    // 通过 名字 查询 key 的方法
type RankObjQueryFunc func(key int64) (obj RankObj) // 通过 key 查询 排行对象 的方法
type AddOrUpdateFunc func(oldObj, newObj RankObj)   // 添加/更新排行对象方法
type RemoveFunc func(obj RankObj)                   // 移除排行对象方法

func (beforeFunc AddOrUpdateFunc) AndThen(thenFunc func(oldObj, newObj RankObj)) AddOrUpdateFunc {
	return func(oldObj, newObj RankObj) {
		beforeFunc(oldObj, newObj)
		thenFunc(oldObj, newObj)
	}
}

func (beforeFunc RemoveFunc) AndThen(thenFunc func(obj RankObj)) RemoveFunc {
	return func(obj RankObj) {
		beforeFunc(obj)
		thenFunc(obj)
	}
}
