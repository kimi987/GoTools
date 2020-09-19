package guildsnapshotdata

import (
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/u64"
	"sync"
	"github.com/lightpaw/male7/util/idbytes"
)

type GuildSnapshot struct {
	Id int64

	// 帮派名字
	Name string

	// 旗号
	FlagName string

	FlagType uint64

	// 帮派等级
	GuildLevel *guild_data.GuildLevelData

	Country  *country.CountryData // 国家
	Prestige uint64               // 声望

	// 盟主id
	LeaderId              int64
	LeaderSnapshotIfIsNpc *shared_proto.HeroBasicSnapshotProto

	IsNpcGuild bool

	MemberCount   uint64  // 帮派成员个数（包含Npc）
	UserMemberIds []int64 // 玩家成员Id列表

	ResistXiongNuDefenders []int64 // 匈奴防守者

	Technologys []*guild_data.GuildTechnologyData // 联盟科技

	TotalPrestigeDaily uint64 // 总的每日声望

	// 入盟条件
	RejectAutoJoin        bool   // false表示达到条件直接入盟，true表示需要申请才能加入
	RequiredHeroLevel     uint64 // 君主等级
	RequiredJunXianLevel  uint64 // 百战军衔
	RequiredTowerMaxFloor uint64 // 需要的最大千重楼层数

	Text string // 对外公告

	basicProto     *shared_proto.GuildBasicProto
	basicProtoOnce sync.Once

	proto     *shared_proto.GuildSnapshotProto
	protoOnce sync.Once

	heroGuildProto      *shared_proto.HeroGuildProto
	heroGuildEncodeOnce sync.Once
}

func (g *GuildSnapshot) GetCountryId() uint64 {
	if g.Country != nil {
		return g.Country.Id
	}
	return 0
}

func (g *GuildSnapshot) Encode(heroSnapshotGetter func(id int64) *shared_proto.HeroBasicSnapshotProto) *shared_proto.GuildSnapshotProto {
	g.protoOnce.Do(func() {
		g.encode(heroSnapshotGetter)
	})
	return g.proto
}

func (g *GuildSnapshot) encode(heroSnapshotGetter func(id int64) *shared_proto.HeroBasicSnapshotProto) {
	proto := &shared_proto.GuildSnapshotProto{}

	proto.Id = i64.Int32(g.Id)
	proto.Name = g.Name
	proto.FlagName = g.FlagName
	proto.FlagType = u64.Int32(g.FlagType)
	proto.Level = u64.Int32(g.GuildLevel.Level)
	proto.MemberCount = u64.Int32(g.MemberCount)

	if g.LeaderSnapshotIfIsNpc != nil {
		proto.Leader = g.LeaderSnapshotIfIsNpc
	} else {
		proto.Leader = heroSnapshotGetter(g.LeaderId)
	}

	proto.Text = g.Text

	proto.RejectAutoJoin = g.RejectAutoJoin
	proto.RequiredHeroLevel = u64.Int32(g.RequiredHeroLevel)
	proto.RequiredJunXianLevel = u64.Int32(g.RequiredJunXianLevel)
	proto.RequiredTowerMaxFloor = u64.Int32(g.RequiredTowerMaxFloor)

	if g.Country != nil {
		proto.PrestigeTarget = u64.Int32(g.Country.Id)
	}

	proto.Prestige = u64.Int32(g.Prestige)

	g.proto = proto
}

func (g *GuildSnapshot) BasicProto() *shared_proto.GuildBasicProto {
	g.basicProtoOnce.Do(func() {
		g.encodeBasicProto()
	})
	return g.basicProto
}

func (g *GuildSnapshot) encodeBasicProto() {
	proto := &shared_proto.GuildBasicProto{}
	proto.Id = i64.Int32(g.Id)
	proto.Name = g.Name
	proto.FlagName = g.FlagName
	proto.Level = u64.Int32(g.GuildLevel.Level)
	if g.Country != nil {
		proto.Country = u64.Int32(g.Country.Id)
	}
	g.basicProto = proto
}

func (g *GuildSnapshot) HeroGuildProto() *shared_proto.HeroGuildProto {
	g.heroGuildEncodeOnce.Do(func() {
		g.encodeHeroGuildProto()
	})
	return g.heroGuildProto
}

func (g *GuildSnapshot) encodeHeroGuildProto() {
	proto := &shared_proto.HeroGuildProto{}
	proto.Id = i64.Int32(g.Id)
	proto.Name = g.Name
	proto.FlagName = g.FlagName
	proto.Level = u64.Int32(g.GuildLevel.Level)
	if g.Country != nil {
		proto.Country = u64.Int32(g.Country.Id)
	}
	proto.Leader = idbytes.ToBytes(g.LeaderId)

	g.heroGuildProto = proto
}

type Getter func(int64) *GuildSnapshot

type Callback interface {
	OnGuildSnapshotUpdated(origin, update *GuildSnapshot)
	OnGuildSnapshotRemoved(id int64)
}
