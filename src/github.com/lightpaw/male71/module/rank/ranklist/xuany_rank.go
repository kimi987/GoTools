package ranklist

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/u64"
	"sync"
	"time"
)

var xuanyRankObjPool = sync.Pool{New: func() interface{} {
	return &XuanyRankObj{}
}}

func NewXuanyRankHolder(maxRankCount uint64, nameQueryFunc rankface.NameQueryFunc, rankObjArray []rankface.RankObj) SingleRankHolder {
	rankObjQueryFunc, addOrUpdateFunc, removeFunc := newRankListLstrFuncs()

	newAddOrUpdateFunc := addOrUpdateFunc.AndThen(func(oldObj, newObj rankface.RankObj) {
		if oldObj != nil {
			xuanyRankObjPool.Put(oldObj)
		}
	})

	newRemoveFunc := removeFunc.AndThen(func(obj rankface.RankObj) {
		if obj != nil {
			xuanyRankObjPool.Put(obj)
		}
	})

	h := &xuany_rank_holder{}

	h.rank_holder = rank_holder{nameQueryFunc: nameQueryFunc}
	rankList := newRankListWithFunc(shared_proto.RankType_RankXuanyuan, maxRankCount, rankObjQueryFunc, newAddOrUpdateFunc, newRemoveFunc)
	rankList.rankObjArray = rankObjArray
	for _, v := range rankObjArray {
		rankList.addOrUpdateFunc(nil, v)
	}

	h.rankList = rankList

	return h
}

type xuany_rank_holder struct {
	single_rank_holder
}

func (h *xuany_rank_holder) SelfKey(hc iface.HeroController) int64 {
	return hc.Id()
}

func NewXuanyRankObj(heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot, heroId int64, point, win, lose, rank uint64) *XuanyRankObj {
	obj := xuanyRankObjPool.Get().(*XuanyRankObj)

	obj.rank_obj = newRankObj(heroId, shared_proto.RankType_RankXuanyuan)
	obj.heroSnapshotGetter = heroSnapshotGetter
	obj.point = point
	obj.win = win
	obj.lose = lose
	obj.SetRank(rank)

	return obj
}

type XuanyRankObj struct {
	*rank_obj
	heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot // 玩家镜像数据获得方法
	point              uint64                                 // 星数
	time               time.Time                              // 时间
	win                uint64                                 // 胜利次数
	lose               uint64                                 // 失败次数
}

func (o *XuanyRankObj) Less(obj rankface.RankObj) bool {
	xuanyObj, ok := obj.(*XuanyRankObj)
	if !ok {
		logrus.Errorf("轩辕会武排行榜里面放的数据竟然不是 XuanyRankObj!%+v", obj)
		return true
	}

	// 层级高的在前面
	if o.point != xuanyObj.point {
		return o.point > xuanyObj.point
	}

	return o.key < obj.Key()
}

func (o *XuanyRankObj) EncodeClient(proto *shared_proto.RankProto) {
	heroSnapshot := o.heroSnapshotGetter(o.Key())
	if heroSnapshot == nil {
		logrus.WithField("hero id", o.Key()).Errorln("没有取到玩家的镜像数据")
		return
	}

	rankProto := &shared_proto.XuanyRankProto{
		Hero:  heroSnapshot.EncodeBasic4Client(),
		Point: u64.Int32(o.point),
		Win: u64.Int32(o.win),
		Lose: u64.Int32(o.lose),
	}

	proto.Xuanyuan = append(proto.Xuanyuan, rankProto)
}

func (o *XuanyRankObj) EncodeHeroSnapshotProto() (proto *shared_proto.HeroBasicSnapshotProto) {
	heroSnapshot := o.heroSnapshotGetter(o.Key())
	if heroSnapshot == nil {
		logrus.WithField("hero id", o.Key()).Errorln("没有取到玩家的镜像数据")
		return
	}
	proto = heroSnapshot.EncodeClient()
	return
}

func (o *XuanyRankObj) EncodeServer(proto *server_proto.RankServerProto) {
	proto.Xuanyuan = append(proto.Xuanyuan, &server_proto.XuanyRankServerProto{
		HeroId: o.Key(),
		Point:  o.point,
		Win:    o.win,
		Lose:   o.lose,
	})
}
