package ranklist

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func NewCountryRankHolder(maxRankCount uint64, nameQueryFunc rankface.NameQueryFunc) SingleRankHolder {
	h := &country_rank_holder{}

	h.rankList = NewRankList(shared_proto.RankType_Country, maxRankCount)
	h.rank_holder = rank_holder{nameQueryFunc: nameQueryFunc}

	return h
}

type country_rank_holder struct {
	single_rank_holder
}

func (h *country_rank_holder) SelfKey(hc iface.HeroController) int64 {
	guildId, _ := hc.LockGetGuildId()
	return guildId
}

func NewCountryRankObj(guildSnapshotGetter func(int64) *guildsnapshotdata.GuildSnapshot, heroSnapshotGetter func(int64) *snapshotdata.HeroSnapshot,
	guildId int64, name string, cityCount uint64, output uint64, time time.Time) *CountryRankObj {
	return &CountryRankObj{
		rank_obj:            newRankObj(guildId, shared_proto.RankType_Country),
		guildSnapshotGetter: guildSnapshotGetter,
		heroSnapshotGetter:  heroSnapshotGetter,
		name:                name,
		//colorIndex:          colorIndex,
		cityCount: cityCount,
		output:    output,
		time:      time,
	}
}

type CountryRankObj struct {
	*rank_obj
	guildSnapshotGetter func(int64) *guildsnapshotdata.GuildSnapshot // 联盟镜像数据获得方法
	heroSnapshotGetter  func(int64) *snapshotdata.HeroSnapshot       // 玩家镜像数据获得方法
	name                string                                       // 国家名字
	colorIndex          uint64                                       // 国家颜色下标
	cityCount           uint64                                       // 名城数量
	output              uint64                                       // 产出
	time                time.Time                                    // 时间
}

func (o *CountryRankObj) Less(obj rankface.RankObj) bool {
	countryObj, ok := obj.(*CountryRankObj)
	if !ok {
		logrus.Errorf("国家排行榜里面放的数据竟然不是 CountryRankObj!%+v", obj)
		return true
	}

	// 名城数量多的在前面
	if o.cityCount != countryObj.cityCount {
		return o.cityCount > countryObj.cityCount
	}

	// 产出高的在前面
	if o.output != countryObj.output {
		return o.output > countryObj.output
	}

	// 相同的，时间小的在前面
	if o.time != countryObj.time {
		return o.time.Before(countryObj.time)
	}

	return o.key < obj.Key()
}

func (o *CountryRankObj) EncodeClient(proto *shared_proto.RankProto) {
	guildSnapshot := o.guildSnapshotGetter(o.Key())
	if guildSnapshot == nil {
		logrus.WithField("guild id", o.Key()).Errorln("没有取到联盟的镜像数据")
		return
	}

	var leader *shared_proto.HeroBasicProto
	if guildSnapshot.LeaderSnapshotIfIsNpc != nil {
		leader = guildSnapshot.LeaderSnapshotIfIsNpc.GetBasic()
	} else {
		heroSnapshot := o.heroSnapshotGetter(guildSnapshot.LeaderId)
		if heroSnapshot == nil {
			logrus.WithField("hero id", guildSnapshot.LeaderId).Errorln("没有取到玩家的镜像数据")
			return
		}

		leader = heroSnapshot.EncodeBasic4Client()
	}

	if int64(leader.GuildId) != guildSnapshot.Id {
		// 比如NPC联盟的盟主
		// 或者服务器一定的同步问题
		leader = &shared_proto.HeroBasicProto{
			Id:            leader.Id,
			Name:          leader.Name,
			Head:          leader.Head,
			Body:          leader.Body,
			Level:         leader.Level,
			Male:          leader.Male,
			GuildId:       i64.Int32(guildSnapshot.Id),
			GuildName:     guildSnapshot.Name,
			GuildFlagName: guildSnapshot.FlagName,
			VipLevel:      leader.VipLevel,
		}
	}

	rankProto := &shared_proto.CountryRankProto{
		Leader:     leader,
		Name:       o.name,
		ColorIndex: u64.Int32(o.colorIndex),
		CityCount:  u64.Int32(o.cityCount),
		CityOutput: u64.Int32(o.output),
	}

	proto.Country = append(proto.Country, rankProto)
}

func (o *CountryRankObj) EncodeHeroSnapshotProto() (proto *shared_proto.HeroBasicSnapshotProto) {
	heroSnapshot := o.heroSnapshotGetter(o.Key())
	if heroSnapshot == nil {
		logrus.WithField("hero id", o.Key()).Errorln("没有取到玩家的镜像数据")
		return
	}
	proto = heroSnapshot.EncodeClient()
	return
}
