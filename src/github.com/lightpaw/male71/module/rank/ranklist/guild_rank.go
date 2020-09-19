package ranklist

import (
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

//var guildRankObjPool = sync.Pool{New: func() interface{} {
//	return &GuildRankObj{}
//}}

func NewGuildRankHolder(configDatas iface.ConfigDatas, maxRankCount uint64, nameQueryFunc rankface.NameQueryFunc) SubTypeRankHolder {
	h := &guild_rank_holder{
		rankList:     make(map[uint64]*guild_rank_list),
		maxRankCount: maxRankCount,
	}

	h.sub_type_rank_holder = sub_type_rank_holder{
		rank_holder:              rank_holder{nameQueryFunc: nameQueryFunc},
		getRankListBySubTypeFunc: h.getRankListByCountryId,
		getRankListByObjFunc:     h.getRankListByObj,
		rangeRankListFunc:        h.rangeList,
		addOrUpdateFunc:          h.addOrUpdateObj,
		removeFunc:               h.removeKeys,
	}

	h.globalRankList = newGuildRankList(nil, h.maxRankCount)
	h.rankList[0] = h.globalRankList
	for _, c := range configDatas.GetCountryDataArray() {
		h.rankList[c.Id] = newGuildRankList(c, h.maxRankCount)
	}

	return h
}

type guild_rank_holder struct {
	sub_type_rank_holder
	rankList       map[uint64]*guild_rank_list
	globalRankList *guild_rank_list

	maxRankCount uint64
}

func (h *guild_rank_holder) SelfKey(hc iface.HeroController) int64 {
	guildId, _ := hc.LockGetGuildId()
	return guildId
}

func (h *guild_rank_holder) getRankListByCountryId(countryId uint64) rankface.RankList {
	return h.rankList[countryId]
}

func (h *guild_rank_holder) getRankListByObj(obj rankface.RankObj) rankface.RankList {
	guildObj, ok := obj.(*GuildRankObj)
	if !ok {
		logrus.WithField("obj", fmt.Sprintf("%+v", obj)).Errorln("guild_rank_holder.RankListByObj obj竟然不是 GuildRankObj 类型")
		return nil
	}

	return h.getRankListByCountryId(guildObj.country.Id)
}

func (h *guild_rank_holder) rangeList(f func(list rankface.RankList)) {
	for _, rl := range h.rankList {
		f(rl)
	}
}

func (h *guild_rank_holder) addOrUpdateObj(objs ...rankface.RankObj) {

	for _, obj := range objs {
		guildObj, ok := obj.(*GuildRankObj)
		if !ok {
			logrus.WithField("obj", fmt.Sprintf("%+v", obj)).Errorln("guild_rank_holder.addOrUpdateObj obj竟然不是 GuildRankObj 类型")
			continue
		}

		// 更新对应国家的表
		rankList := h.getRankListByCountryId(guildObj.country.Id)
		if rankList != nil {
			rankList.AddOrUpdate(guildObj.copy())
		}

		// 总表更新
		h.globalRankList.AddOrUpdate(obj)
	}
}

func (h *guild_rank_holder) removeKeys(keys ...int64) {
	h.rangeList(func(list rankface.RankList) {
		for _, key := range keys {
			list.Remove(key)
		}
	})
}

func newGuildRankList(country *country.CountryData, maxRankCount uint64) *guild_rank_list {
	return &guild_rank_list{
		country:  country,
		RankList: NewRankList(shared_proto.RankType_Guild, maxRankCount),
	}
}

type guild_rank_list struct {
	country *country.CountryData
	rankface.RankList
}

func (rl *guild_rank_list) EncodeClient(startRank, rankCountPerPage uint64) *shared_proto.RankProto {
	proto := rl.RankList.EncodeClient(startRank, rankCountPerPage)

	if rl.country != nil {
		proto.SubType = u64.Int32(rl.country.Id)
	}

	return proto
}

func NewGuildRankObj(guildSnapshotGetter func(int64) *guildsnapshotdata.GuildSnapshot, heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot, guild *sharedguilddata.Guild) *GuildRankObj {
	//obj := guildRankObjPool.Get().(*GuildRankObj)
	obj := &GuildRankObj{}

	obj.rank_obj = newRankObj(guild.Id(), shared_proto.RankType_Guild)
	obj.guildSnapshotGetter = guildSnapshotGetter
	obj.heroSnapshotGetter = heroSnapshotGetter
	obj.prestige = guild.GetPrestige()
	obj.country = guild.Country()
	obj.level = guild.LevelData().Level
	obj.cityCount = 0 // 名城数量 TODO 做名城战的时候加上
	obj.upgradeEndTime = guild.GetUpgradeEndTime()
	obj.time = guild.UpdateBuildingAmountTime()
	obj.SetRank(0)

	return obj
}

type GuildRankObj struct {
	*rank_obj
	guildSnapshotGetter func(int64) *guildsnapshotdata.GuildSnapshot // 联盟镜像数据获得方法
	heroSnapshotGetter  func(int64) *snapshotdata.HeroSnapshot       // 玩家镜像数据获得方法
	prestige            uint64                                       // 联盟声望
	country             *country.CountryData                         // 联盟国家
	level               uint64                                       // 联盟等级
	cityCount           uint64                                       // 名城数量
	upgradeEndTime      time.Time                                    // 升级结束时间，0表示当前没升级
	time                time.Time                                    // 时间
}

func (o *GuildRankObj) copy() *GuildRankObj {
	newObj := &GuildRankObj{}
	*newObj = *o

	newObj.rank_obj = o.CopyObj()

	return newObj
}

func (o *GuildRankObj) Less(obj rankface.RankObj) bool {
	guildObj, ok := obj.(*GuildRankObj)
	if !ok {
		logrus.Errorf("联盟排行榜里面放的数据竟然不是 GuildRankObj!%+v", obj)
		return true
	}

	// 声望高的在前面
	if o.prestige != guildObj.prestige {
		return o.prestige > guildObj.prestige
	}

	// 等级高的在前面
	if o.level != guildObj.level {
		return o.level > guildObj.level
	}

	// 名城多的在前面
	if o.cityCount != guildObj.cityCount {
		return o.cityCount > guildObj.cityCount
	}

	// 升级结束早的在前面
	if o.upgradeEndTime != guildObj.upgradeEndTime {
		return o.upgradeEndTime.Before(guildObj.upgradeEndTime)
	}

	// 相同的，时间小的在前面
	if o.time != guildObj.time {
		return o.time.Before(guildObj.time)
	}

	return o.key < obj.Key()
}

func (o *GuildRankObj) EncodeClient(proto *shared_proto.RankProto) {
	guildSnapshot := o.guildSnapshotGetter(o.Key())
	if guildSnapshot == nil {
		logrus.WithField("guild id", o.Key()).Errorln("没有取到联盟的镜像数据")
		return
	}

	var leaderBasic *shared_proto.HeroBasicProto
	if guildSnapshot.LeaderSnapshotIfIsNpc != nil {
		leaderBasic = guildSnapshot.LeaderSnapshotIfIsNpc.GetBasic()
	} else {
		heroSnapshot := o.heroSnapshotGetter(guildSnapshot.LeaderId)
		if heroSnapshot == nil {
			logrus.WithField("hero id", guildSnapshot.LeaderId).Errorln("没有取到玩家的镜像数据")
			return
		}

		leaderBasic = heroSnapshot.EncodeBasic4Client()
	}

	if int64(leaderBasic.GuildId) != guildSnapshot.Id ||
		leaderBasic.GuildName != guildSnapshot.Name ||
		leaderBasic.GuildFlagName != guildSnapshot.FlagName {
		// 比如NPC联盟的盟主
		// 或者服务器一定的同步问题
		leaderBasic = &shared_proto.HeroBasicProto{
			Id:            leaderBasic.Id,
			Name:          leaderBasic.Name,
			Head:          leaderBasic.Head,
			Body:          leaderBasic.Body,
			Level:         leaderBasic.Level,
			Male:          leaderBasic.Male,
			GuildId:       i64.Int32(guildSnapshot.Id),
			GuildName:     guildSnapshot.Name,
			GuildFlagName: guildSnapshot.FlagName,
			VipLevel:      leaderBasic.VipLevel,
		}
	}

	rankProto := &shared_proto.GuildRankProto{
		Leader:         leaderBasic,
		Level:          u64.Int32(o.level),
		UpgradeEndTime: timeutil.Marshal32(o.upgradeEndTime),
		Prestige:       u64.Int32(o.prestige),
		Country:        u64.Int32(o.country.Id),
		CityCount:      u64.Int32(o.cityCount),
		MemberCount:    u64.Int32(guildSnapshot.MemberCount),
	}

	proto.Guild = append(proto.Guild, rankProto)
}

func (o *GuildRankObj) EncodeHeroSnapshotProto() (proto *shared_proto.HeroBasicSnapshotProto) {
	heroSnapshot := o.heroSnapshotGetter(o.Key())
	if heroSnapshot == nil {
		logrus.WithField("hero id", o.Key()).Errorln("没有取到玩家的镜像数据")
		return
	}
	proto = heroSnapshot.EncodeClient()
	return
}
