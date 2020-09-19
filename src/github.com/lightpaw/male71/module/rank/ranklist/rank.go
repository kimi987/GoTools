package ranklist

import (
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
)

func newRankObj(key int64, rankType shared_proto.RankType) *rank_obj {
	return &rank_obj{
		key:      key,
		rankType: rankType,
	}
}

// 排行对象
type rank_obj struct {
	key      int64                 // key
	rankType shared_proto.RankType // 所在的排行榜
	rank     uint64                // 排名
}

func (o *rank_obj) Key() int64 {
	return o.key
}

func (o *rank_obj) RankType() shared_proto.RankType {
	return o.rankType
}

func (o *rank_obj) Rank() uint64 {
	return o.rank
}

func (o *rank_obj) SetRank(toSet uint64) {
	o.rank = toSet
}

func (o *rank_obj) CopyObj() *rank_obj {
	r := newRankObj(o.key, o.rankType)
	r.rank = o.rank
	return r
}

// 默认给个空的方法，别的地方就不用再继承了
func (o *rank_obj) EncodeServer(proto *server_proto.RankServerProto) {}
