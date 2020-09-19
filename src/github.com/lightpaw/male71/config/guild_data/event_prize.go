package guild_data

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/promdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"time"
)

// 联盟事件礼包
//gogen:config
type GuildEventPrizeData struct {
	_ struct{} `file:"联盟/联盟盟友礼包.txt"`
	_ struct{} `proto:"shared_proto.GuildEventPrizeDataProto"`
	_ struct{} `protoconfig:"guild_event_prize"`

	Id uint64

	Name string

	Desc string

	Quality shared_proto.Quality

	Icon *icon.Icon `protofield:"IconId,%s.Id"` // 图标

	Prize *resdata.PlunderPrize `protofield:"Prize,%s.Prize.PrizeProto()"`

	ExipreDuration time.Duration `protofield:"-"` // 过期时间

	FromShop bool `default:"false"` // 是否商城购得

	Energy uint64 // 大礼包能量

	DailyLimit uint64 `validator:"uint" protofield:"-"` // 每日领取上限

	GuildLevelPrizeGroupId uint64                     `validator:"uint" default:"0" protofield:"-"` // 联盟等级奖励组id
	GuildLevelPrizes       []*resdata.GuildLevelPrize `head:"-"`                                    // 联盟等级奖励组

	// 英雄事件触发规则
	TriggerEvent           shared_proto.HeroEvent `white:"0" protofield:"-"`
	TriggerEventCondition  *data.CompareCondition `protofield:"-"`
	TriggerEventTimes      uint64                 `validator:"uint" protofield:"-"`
	TriggerEventDailyReset bool                   `protofield:"-"`
}

func (data *GuildEventPrizeData) InitAll(filename string, configs interface {
	GetGuildEventPrizeData(uint64) *GuildEventPrizeData
	GetEventLimitGiftDataArray() []*promdata.EventLimitGiftData
	GetTimeLimitGiftDataArray() []*promdata.TimeLimitGiftData
}) {
	for _, d := range configs.GetEventLimitGiftDataArray() {
		if d.GuildEventPrizeId > 0 {
			check.PanicNotTrue(configs.GetGuildEventPrizeData(d.GuildEventPrizeId) != nil, "")
		}
	}
	for _, d := range configs.GetTimeLimitGiftDataArray() {
		if d.GuildEventPrizeId > 0 {
			check.PanicNotTrue(configs.GetGuildEventPrizeData(d.GuildEventPrizeId) != nil, "")
		}
	}
}

func (d *GuildEventPrizeData) Init(configs interface {
	GetGuildLevelPrizeArray() []*resdata.GuildLevelPrize
}) {
	d.GuildLevelPrizes = resdata.GuildLevelPrizeGroup(d.GuildLevelPrizeGroupId, configs)
}

//gogen:config
type GuildBigBoxData struct {
	_ struct{} `file:"联盟/联盟大宝箱.txt"`
	_ struct{} `proto:"shared_proto.GuildBigBoxDataProto"`
	_ struct{} `protoconfig:"guild_big_box"`

	Id uint64

	PlunderPrize *resdata.PlunderPrize `protofield:"Prize,%s.Prize.PrizeProto()"` // 掉落奖励

	UnlockEnergy uint64 // 解锁能量

	GuildLevelPrizeGroupId uint64                     `validator:"uint" default:"0" protofield:"-"` // 联盟等级奖励组id
	GuildLevelPrizes       []*resdata.GuildLevelPrize `head:"-"`                                    // 联盟等级奖励组

	TechLevel uint64 `head:"-"`
}

func (d *GuildBigBoxData) Init(configs interface {
	GetGuildLevelPrizeArray() []*resdata.GuildLevelPrize
}) {
	d.GuildLevelPrizes = resdata.GuildLevelPrizeGroup(d.GuildLevelPrizeGroupId, configs)
}
