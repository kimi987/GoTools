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

var towerRankObjPool = sync.Pool{New: func() interface{} {
	return &TowerRankObj{}
}}

func NewTowerRankHolder(maxRankCount uint64, nameQueryFunc rankface.NameQueryFunc) SingleRankHolder {
	rankObjQueryFunc, addOrUpdateFunc, removeFunc := newRankListLstrFuncs()

	newAddOrUpdateFunc := addOrUpdateFunc.AndThen(func(oldObj, newObj rankface.RankObj) {
		if oldObj != nil {
			towerRankObjPool.Put(oldObj)
		}
	})

	newRemoveFunc := removeFunc.AndThen(func(obj rankface.RankObj) {
		if obj != nil {
			towerRankObjPool.Put(obj)
		}
	})

	h := &tower_rank_holder{}

	h.rank_holder = rank_holder{nameQueryFunc: nameQueryFunc}
	h.rankList = NewRankListWithFunc(shared_proto.RankType_Tower, maxRankCount, rankObjQueryFunc, newAddOrUpdateFunc, newRemoveFunc)

	return h
}

type tower_rank_holder struct {
	single_rank_holder
}

func (h *tower_rank_holder) SelfKey(hc iface.HeroController) int64 {
	return hc.Id()
}

func NewTowerRankObj(heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot, heroId int64, maxFloor uint64, time time.Time) *TowerRankObj {
	obj := towerRankObjPool.Get().(*TowerRankObj)

	obj.rank_obj = newRankObj(heroId, shared_proto.RankType_Tower)
	obj.heroSnapshotGetter = heroSnapshotGetter
	obj.maxFloor = maxFloor
	obj.time = time
	obj.SetRank(0)

	return obj
}

type TowerRankObj struct {
	*rank_obj
	heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot // 玩家镜像数据获得方法
	maxFloor           uint64                                 // 最大层数
	time               time.Time                              // 时间
}

func (o *TowerRankObj) Less(obj rankface.RankObj) bool {
	towerObj, ok := obj.(*TowerRankObj)
	if !ok {
		logrus.Errorf("千重楼排行榜里面放的数据竟然不是 TowerRankObj!%+v", obj)
		return true
	}

	// 层级高的在前面
	if o.maxFloor != towerObj.maxFloor {
		return o.maxFloor > towerObj.maxFloor
	}

	// 层级相同的，时间小的在前面
	if o.time != towerObj.time {
		return o.time.Before(towerObj.time)
	}

	return o.key < obj.Key()
}

func (o *TowerRankObj) EncodeClient(proto *shared_proto.RankProto) {
	heroSnapshot := o.heroSnapshotGetter(o.Key())
	if heroSnapshot == nil {
		logrus.WithField("hero id", o.Key()).Errorln("没有取到玩家的镜像数据")
		return
	}

	rankProto := &shared_proto.TowerRankProto{
		Hero:  heroSnapshot.EncodeBasic4Client(),
		Floor: u64.Int32(o.maxFloor),
		Time:  timeutil.Marshal32(o.time),
	}

	proto.Tower = append(proto.Tower, rankProto)
}

func (o *TowerRankObj) EncodeHeroSnapshotProto() (proto *shared_proto.HeroBasicSnapshotProto) {
	heroSnapshot := o.heroSnapshotGetter(o.Key())
	if heroSnapshot == nil {
		logrus.WithField("hero id", o.Key()).Errorln("没有取到玩家的镜像数据")
		return
	}
	proto = heroSnapshot.EncodeClient()
	return
}

func (o *TowerRankObj) EncodeServer(proto *server_proto.RankServerProto) {
	proto.Tower = append(proto.Tower, &server_proto.TowerRankServerProto{
		HeroId:   o.Key(),
		MaxFloor: o.maxFloor,
		Time:     timeutil.Marshal64(o.time),
	})
}
