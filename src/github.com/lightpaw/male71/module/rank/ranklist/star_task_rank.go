package ranklist

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"sync"
	"time"
)

var starTaskRankObjPool = sync.Pool{New: func() interface{} {
	return &StarTaskRankObj{}
}}

func NewStarTaskRankHolder(maxRankCount uint64, nameQueryFunc rankface.NameQueryFunc) SingleRankHolder {
	rankObjQueryFunc, addOrUpdateFunc, removeFunc := newRankListLstrFuncs()

	newAddOrUpdateFunc := addOrUpdateFunc.AndThen(func(oldObj, newObj rankface.RankObj) {
		if oldObj != nil {
			starTaskRankObjPool.Put(oldObj)
		}
	})

	newRemoveFunc := removeFunc.AndThen(func(obj rankface.RankObj) {
		if obj != nil {
			starTaskRankObjPool.Put(obj)
		}
	})

	h := &starTask_rank_holder{}

	h.rank_holder = rank_holder{nameQueryFunc: nameQueryFunc}
	h.rankList = NewRankListWithFunc(shared_proto.RankType_RankStarTask, maxRankCount, rankObjQueryFunc, newAddOrUpdateFunc, newRemoveFunc)

	return h
}

type starTask_rank_holder struct {
	single_rank_holder
}

func (h *starTask_rank_holder) SelfKey(hc iface.HeroController) int64 {
	return hc.Id()
}

func NewStarTaskRankObj(heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot, heroId int64, star uint64, time time.Time) *StarTaskRankObj {
	obj := starTaskRankObjPool.Get().(*StarTaskRankObj)

	obj.rank_obj = newRankObj(heroId, shared_proto.RankType_RankStarTask)
	obj.heroSnapshotGetter = heroSnapshotGetter
	obj.star = star
	obj.time = time
	obj.SetRank(0)

	return obj
}

type StarTaskRankObj struct {
	*rank_obj
	heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot // 玩家镜像数据获得方法
	star               uint64                                 // 星数
	time               time.Time                              // 时间
}

func (o *StarTaskRankObj) Less(obj rankface.RankObj) bool {
	starTaskObj, ok := obj.(*StarTaskRankObj)
	if !ok {
		logrus.Errorf("成就星数排行榜里面放的数据竟然不是 StarTaskRankObj!%+v", obj)
		return true
	}

	// 层级高的在前面
	if o.star != starTaskObj.star {
		return o.star > starTaskObj.star
	}

	// 层级相同的，时间小的在前面
	if o.time != starTaskObj.time {
		return o.time.Before(starTaskObj.time)
	}

	return o.key < obj.Key()
}

func (o *StarTaskRankObj) EncodeClient(proto *shared_proto.RankProto) {
	heroSnapshot := o.heroSnapshotGetter(o.Key())
	if heroSnapshot == nil {
		logrus.WithField("hero id", o.Key()).Errorln("没有取到玩家的镜像数据")
		return
	}

	rankProto := &shared_proto.StarTaskRankProto{
		Hero: heroSnapshot.EncodeBasic4Client(),
		Star: u64.Int32(o.star),
	}

	proto.StarTask = append(proto.StarTask, rankProto)
}

func (o *StarTaskRankObj) EncodeHeroSnapshotProto() (proto *shared_proto.HeroBasicSnapshotProto) {
	heroSnapshot := o.heroSnapshotGetter(o.Key())
	if heroSnapshot == nil {
		logrus.WithField("hero id", o.Key()).Errorln("没有取到玩家的镜像数据")
		return
	}
	proto = heroSnapshot.EncodeClient()
	return
}

func (o *StarTaskRankObj) EncodeServer(proto *server_proto.RankServerProto) {
	proto.StarTask = append(proto.StarTask, &server_proto.StarTaskRankServerProto{
		HeroId: o.Key(),
		Star:   o.star,
		Time:   timeutil.Marshal64(o.time),
	})
}
