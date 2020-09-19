package guild_data

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
)

func GetPrestigeEventId(t shared_proto.HeroEvent, subType uint64) uint64 {
	return uint64(t) | subType<<8
}

// 联盟声望事件
//gogen:config
type GuildPrestigeEventData struct {
	_ struct{} `file:"联盟/联盟声望事件.txt"`

	Id uint64 `head:"-,GetPrestigeEventId(%s.TriggerEvent%c %s.TriggerEventCondition.Amount)"`

	// 英雄事件触发规则
	TriggerEvent          shared_proto.HeroEvent `white:"0" protofield:"-"`
	TriggerEventCondition *data.CompareCondition `protofield:"-"`
	TriggerEventTimes     uint64                 `validator:"uint" protofield:"-"`
	//TriggerEventDailyReset bool                   `protofield:"-"`

	Prestige uint64

	Hufu uint64 `validator:"uint"`

	IgnoreMemberLimit bool
}

func (d *GuildPrestigeEventData) Init(filename string) {

	if d.TriggerEventCondition.Greater || d.TriggerEventCondition.Less {
		check.PanicNotTrue(d.TriggerEventTimes > 0, "%s 配置的id[%d]的触发次数必须>0", filename, d.Id)
	}

}

// 联盟声望礼包奖励
//gogen:config
type GuildPrestigePrizeData struct {
	_ struct{} `file:"联盟/联盟声望礼包.txt"`
	_ struct{} `proto:"shared_proto.GuildPrestigePrizeDataProto"`
	_ struct{} `protoconfig:"guild_prestige_prize"`

	Prestige uint64 `key:"true"`

	EventPrize *GuildEventPrizeData `protofield:",config.U64ToI32(%s.Id)"` // 奖励的礼包

	BuildingAmount uint64 // 奖励的联盟建设值

	Hufu uint64 // 虎符
}

func (d *GuildPrestigePrizeData) Init(filename string) {
	check.PanicNotTrue(d.EventPrize.TriggerEvent == shared_proto.HeroEvent_InvalidHeroEvent, "%s 配置的声望 %d的奖励礼包%d-%s，触发类型必须是0", filename, d.Prestige, d.EventPrize.Id, d.EventPrize.Name)
}

// 联盟排行礼包
//gogen:config
type GuildRankPrizeData struct {
	_ struct{} `file:"联盟/联盟排行奖励.txt"`
	_ struct{} `protogen:"true"`

	Rank                uint64 `key:"true"`
	Prize               *resdata.Prize
	CountryDestroyPrize *resdata.Prize
}

func (*GuildRankPrizeData) InitAll(filename string, configs interface {
	GetGuildRankPrizeDataArray() []*GuildRankPrizeData
}) {
	for i, data := range configs.GetGuildRankPrizeDataArray() {
		check.PanicNotTrue(data.Rank == u64.FromInt(i+1), "%s 联盟排名奖励的主键Rank只能够从1开始自增")
	}
}
