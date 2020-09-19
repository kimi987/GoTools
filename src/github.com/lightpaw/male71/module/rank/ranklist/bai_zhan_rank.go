package ranklist

import (
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/bai_zhan_data"
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

var baiZhanRankObjPool = sync.Pool{New: func() interface{} {
	return &BaiZhanRankObj{}
}}

func NewBaiZhanRankHolder(configDatas iface.ConfigDatas, nameQueryFunc rankface.NameQueryFunc) SubTypeRankHolder {
	queryFunc, addOrUpdateFunc, removeFunc := newRankListLstrFuncs()

	newAddOrUpdateFunc := addOrUpdateFunc.AndThen(func(oldObj, newObj rankface.RankObj) {
		if oldObj != nil {
			baiZhanRankObjPool.Put(oldObj)
		}
	})

	newRemoveFunc := removeFunc.AndThen(func(obj rankface.RankObj) {
		if obj != nil {
			baiZhanRankObjPool.Put(obj)
		}
	})

	h := &bai_zhan_rank_holder{
		rankObjQuery: queryFunc,
	}
	h.sub_type_rank_holder = sub_type_rank_holder{
		rank_holder:              rank_holder{nameQueryFunc: nameQueryFunc},
		getRankListBySubTypeFunc: h.getRankListByLevel,
		getRankListByObjFunc:     h.getRankListByObj,
		rangeRankListFunc:        h.rangeList,
		addOrUpdateFunc:          h.addOrUpdateObj,
		removeFunc:               h.removeKeys,
	}

	for _, junXianLevelData := range configDatas.GetJunXianLevelDataArray() {
		h.rankList = append(h.rankList, &bai_zhan_rank_list{
			junXianLevel: junXianLevelData.Level,
			RankList:     NewRankListWithFunc(shared_proto.RankType_BaiZhan, configDatas.RankMiscData().MaxRankCount, queryFunc, newAddOrUpdateFunc, newRemoveFunc),
		})
	}

	return h
}

type bai_zhan_rank_holder struct {
	sub_type_rank_holder
	rankList []*bai_zhan_rank_list

	rankObjQuery rankface.RankObjQueryFunc
}

func (h *bai_zhan_rank_holder) SelfKey(hc iface.HeroController) int64 {
	return hc.Id()
}

func (h *bai_zhan_rank_holder) getRankListByLevel(junXianLevel uint64) rankface.RankList {
	if junXianLevel <= 0 {
		return h.rankList[0]
	} else if junXianLevel > uint64(len(h.rankList)) {
		return h.rankList[len(h.rankList)-1]
	} else {
		return h.rankList[junXianLevel-1]
	}
}

func (h *bai_zhan_rank_holder) getRankListByObj(obj rankface.RankObj) rankface.RankList {
	baiZhanObj, ok := obj.(*BaiZhanRankObj)
	if !ok {
		logrus.WithField("obj", fmt.Sprintf("%+v", obj)).Errorln("bai_zhan_rank_holder.RankListByObj obj竟然不是 BaiZhanRankObj 类型")
		return nil
	}

	return h.getRankListByLevel(baiZhanObj.JunXianLevel().Level)
}

func (h *bai_zhan_rank_holder) rangeList(f func(list rankface.RankList)) {
	for _, rl := range h.rankList {
		f(rl)
	}
}

func (h *bai_zhan_rank_holder) addOrUpdateObj(objs ...rankface.RankObj) {
	for _, obj := range objs {
		rankList := h.getRankListByObj(obj)
		if rankList != nil {
			rankList.AddOrUpdate(obj)
		}
	}
}

func (h *bai_zhan_rank_holder) removeKeys(keys ...int64) {
	for _, key := range keys {
		obj := h.rankObjQuery(key)
		if obj != nil {
			rankList := h.getRankListByObj(obj)
			if rankList != nil {
				rankList.Remove(key)
			}
		}
	}
}

type bai_zhan_rank_list struct {
	junXianLevel uint64 // 百战等级
	rankface.RankList
}

func (rl *bai_zhan_rank_list) EncodeClient(startRank, rankCountPerPage uint64) *shared_proto.RankProto {
	proto := rl.RankList.EncodeClient(startRank, rankCountPerPage)

	proto.JunXianLevel = u64.Int32(rl.junXianLevel)
	proto.SubType = u64.Int32(rl.junXianLevel)

	return proto
}

func NewBaiZhanRankObj(heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot, baiZhanPointGetter func(id int64) uint64, heroId int64,
	curJunXianLevel, lastJunXianLevel *bai_zhan_data.JunXianLevelData, time time.Time, point, fightAmount uint64) *BaiZhanRankObj {
	obj := baiZhanRankObjPool.Get().(*BaiZhanRankObj)

	obj.rank_obj = newRankObj(heroId, shared_proto.RankType_BaiZhan)
	obj.heroSnapshotGetter = heroSnapshotGetter
	obj.baiZhanPointGetter = baiZhanPointGetter
	obj.curJunXianLevel = curJunXianLevel
	obj.lastJunXianLevel = lastJunXianLevel
	obj.point = point
	obj.time = time
	obj.fightAmount = fightAmount
	obj.SetRank(0)

	return obj
}

type BaiZhanRankObj struct {
	*rank_obj
	heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot // 玩家镜像数据获得方法
	baiZhanPointGetter func(id int64) uint64                  // 百战积分获得
	curJunXianLevel    *bai_zhan_data.JunXianLevelData        // 当前百战军衔等级
	lastJunXianLevel   *bai_zhan_data.JunXianLevelData        // 昨日百战军衔等级
	point              uint64                                 // 积分
	time               time.Time                              // 时间
	fightAmount        uint64
}

func (o *BaiZhanRankObj) Less(obj rankface.RankObj) bool {
	baiZhanObj, ok := obj.(*BaiZhanRankObj)
	if !ok {
		logrus.Errorf("百战排行榜里面放的数据竟然不是 BaiZhanRankObj!%+v", obj)
		return true
	}

	// 层级高的在前面
	if o.lastJunXianLevel != baiZhanObj.lastJunXianLevel {
		return o.lastJunXianLevel.Level > baiZhanObj.lastJunXianLevel.Level
	}

	// 层级高的在前面
	if o.curJunXianLevel != baiZhanObj.curJunXianLevel {
		return o.curJunXianLevel.Level > baiZhanObj.curJunXianLevel.Level
	}

	if o.point != baiZhanObj.point {
		return o.point > baiZhanObj.point
	}

	// 层级相同的，时间小的在前面
	if o.time != baiZhanObj.time {
		return o.time.Before(baiZhanObj.time)
	}

	return o.key < obj.Key()
}

func (o *BaiZhanRankObj) JunXianLevel() *bai_zhan_data.JunXianLevelData {
	return o.lastJunXianLevel
}

func (o *BaiZhanRankObj) EncodeClient(proto *shared_proto.RankProto) {
	heroSnapshot := o.heroSnapshotGetter(o.Key())
	if heroSnapshot == nil {
		logrus.WithField("hero id", o.Key()).Errorln("没有取到玩家的镜像数据")
		return
	}

	levelChangeType := shared_proto.LevelChangeType_LEVEL_KEEP
	if o.curJunXianLevel != o.lastJunXianLevel {
		if o.curJunXianLevel.Level < o.lastJunXianLevel.Level {
			levelChangeType = shared_proto.LevelChangeType_LEVEL_DOWN
		} else {
			levelChangeType = shared_proto.LevelChangeType_LEVEL_UP
		}
	}

	rankProto := &shared_proto.BaiZhanRankProto{
		Hero:            heroSnapshot.EncodeBasic4Client(),
		Level:           u64.Int32(o.lastJunXianLevel.Level),
		Point:           u64.Int32(o.point),
		FightAmount:     u64.Int32(o.fightAmount),
		LevelChangeType: levelChangeType,
	}

	proto.BaiZhan = append(proto.BaiZhan, rankProto)
}

func (o *BaiZhanRankObj) EncodeServer(proto *server_proto.RankServerProto) {
	proto.BaiZhan = append(proto.BaiZhan, &server_proto.BaiZhanRankServerProto{
		HeroId:           o.Key(),
		JunXianLevel:     o.curJunXianLevel.Level,
		LastJunXianLevel: o.lastJunXianLevel.Level,
		Time:             timeutil.Marshal64(o.time),
		Point:            o.point,
		FightAmount:      o.fightAmount,
	})
}

func (o *BaiZhanRankObj) EncodeHeroSnapshotProto() (proto *shared_proto.HeroBasicSnapshotProto) {
	heroSnapshot := o.heroSnapshotGetter(o.Key())
	if heroSnapshot == nil {
		logrus.WithField("hero id", o.Key()).Errorln("没有取到玩家的镜像数据")
		return
	}
	proto = heroSnapshot.EncodeClient()
	return
}
