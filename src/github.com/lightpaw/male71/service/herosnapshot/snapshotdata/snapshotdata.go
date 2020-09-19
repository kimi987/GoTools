package snapshotdata

import (
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"sync"
	"time"
)

// -------- snapshot object -------

// 只读, 不要修改, 谁改谁傻逼
// 里面也不能放hero里的类似Depot或Military对象. 必须全部是单独copy出来的只读副本
type HeroSnapshot struct {
	INTERNAL_VERSION         uint64
	clientBasicProto         *shared_proto.HeroBasicProto
	clientBasicProtoMux      sync.RWMutex
	clientBasicProtoBytes    []byte
	clientBasicProtoBytesMux sync.RWMutex
	clientProto              *shared_proto.HeroBasicSnapshotProto
	clientProtoMux           sync.RWMutex
	clientProtoBytes         []byte
	clientProtoBytesMux      sync.RWMutex
	snapshotCreateTime       int32

	getGuildSnapshot guildsnapshotdata.Getter

	Id         int64
	IdBytes    []byte
	Name       string
	Male       bool
	Head       string
	Body       uint64
	Level      uint64
	CreateTime time.Time

	BaseRegion   int64
	BaseLevel    uint64
	BaseX, BaseY int

	FightAmount uint64
	Prosperity  uint64

	GuildId         int64
	CountryId       uint64
	countryOfficial shared_proto.CountryOfficialType

	// 千重楼的最大层数
	TowerMaxFloor uint64

	// 百战千军军衔等级
	BaiZhanJunXianLevel uint64

	// 玩家摇钱树次数
	TreasuryTreeWaterTimes uint64

	// 密室最大的开启了的塔的id
	SecretTowerMaxOpenId uint64
	// 密室挑战次数
	SecretTowerChallengeTimes uint64
	// 密室协助次数
	SecretTowerHelpTimes uint64

	// 玩家摇钱树浇水记录
	HelpMeHeroIds []int64
	HelpMeSeasons []shared_proto.Season

	// 上次离线时间
	LastOfflineTime time.Time

	// 好友数
	FriendsCount uint64

	// vip等级
	VipLevel uint64

	// 设置
	Settings uint64

	// 称号
	Title uint64

	// 故乡
	Location uint64
}

// 必须由HeroSnapshotService.NewSnapshot调用
func NewSnapshot(version uint64, hero *entity.Hero, getGuildSnapshot guildsnapshotdata.Getter, junXianGetter func(id int64) uint64) *HeroSnapshot {
	s := &HeroSnapshot{
		INTERNAL_VERSION:   version,
		snapshotCreateTime: int32(time.Now().Unix()),
		getGuildSnapshot:   getGuildSnapshot,

		Id:         hero.Id(),
		IdBytes:    hero.IdBytes(),
		Name:       hero.Name(),
		Male:       hero.Male(),
		Head:       hero.Head(),
		Body:       hero.Body(),
		Level:      hero.Level(),
		CreateTime: hero.CreateTime(),

		BaseRegion: hero.BaseRegion(),
		BaseLevel:  hero.BaseLevel(),
		BaseX:      hero.BaseX(),
		BaseY:      hero.BaseY(),

		FightAmount: hero.GetHomeDefenserFightAmount(),
		Prosperity:  hero.Prosperity(),

		GuildId:         hero.GuildId(),
		CountryId:       hero.CountryId(),
		countryOfficial: hero.CountryMisc().ShowOfficialType(),

		TowerMaxFloor: hero.Tower().HistoryMaxFloor(),

		BaiZhanJunXianLevel: junXianGetter(hero.Id()),

		SecretTowerMaxOpenId:      hero.SecretTower().MaxOpenSecretTowerId(),
		SecretTowerChallengeTimes: hero.SecretTower().ChallengeTimes(),
		SecretTowerHelpTimes:      hero.SecretTower().HelpTimes(),

		TreasuryTreeWaterTimes: hero.TreasuryTree().WaterTimes(),

		LastOfflineTime: hero.LastOfflineTime(),

		Settings: hero.Settings().EncodeToUint64(),

		Title: hero.TaskList().TitleId(),

		Location: hero.Location(),

		VipLevel: hero.VipLevel(),

		FriendsCount: hero.FriendsCount(),
	}

	s.HelpMeHeroIds, s.HelpMeSeasons = hero.TreasuryTree().HelpMeInfo()

	return s
}

func (s *HeroSnapshot) Version() uint64 {
	return s.INTERNAL_VERSION
}

func (s *HeroSnapshot) SnapshotCreateTime() int32 {
	return s.snapshotCreateTime
}

func (s *HeroSnapshot) Guild() *guildsnapshotdata.GuildSnapshot {
	if s.GuildId != 0 && s.getGuildSnapshot != nil {
		return s.getGuildSnapshot(s.GuildId)
	}

	return nil
}

func (s *HeroSnapshot) GuildFlagName() string {
	if g := s.Guild(); g != nil {
		return g.FlagName
	}

	return ""
}

func (s *HeroSnapshot) ClearProto() {

	s.clientProtoBytesMux.Lock()
	s.clientProtoBytes = nil
	s.clientProtoBytesMux.Unlock()

	s.clientProtoMux.Lock()
	s.clientProto = nil
	s.clientProtoMux.Unlock()

	s.clientBasicProtoBytesMux.Lock()
	s.clientBasicProtoBytes = nil
	s.clientBasicProtoBytesMux.Unlock()

	s.clientBasicProtoMux.Lock()
	s.clientBasicProto = nil
	s.clientBasicProtoMux.Unlock()
}

func (s *HeroSnapshot) EncodeClient() *shared_proto.HeroBasicSnapshotProto {
	if proto := s.getClientProto(); proto != nil {
		return proto
	}

	s.clientProtoMux.Lock()
	defer s.clientProtoMux.Unlock()
	return s.encodeClient()
}

func (s *HeroSnapshot) getClientProto() *shared_proto.HeroBasicSnapshotProto {
	s.clientProtoMux.RLock()
	defer s.clientProtoMux.RUnlock()

	return s.clientProto
}

func (s *HeroSnapshot) encodeClient() *shared_proto.HeroBasicSnapshotProto {
	proto := &shared_proto.HeroBasicSnapshotProto{}

	proto.Basic = s.EncodeBasic4Client()

	proto.BaseRegion = i64.Int32(s.BaseRegion)
	proto.BaseLevel = u64.Int32(s.BaseLevel)
	proto.BaseX = imath.Int32(s.BaseX)
	proto.BaseY = imath.Int32(s.BaseY)

	proto.FightAmount = u64.Int32(s.FightAmount)
	proto.Prosperity = u64.Int32(s.Prosperity)

	proto.TowerMaxFloor = u64.Int32(s.TowerMaxFloor)

	proto.SecretTowerMaxOpenId = u64.Int32(s.SecretTowerMaxOpenId)
	proto.SecretTowerChallengeTimes = u64.Int32(s.SecretTowerChallengeTimes)
	proto.SecretTowerHelpTimes = u64.Int32(s.SecretTowerHelpTimes)

	proto.JunXianLevel = u64.Int32(s.BaiZhanJunXianLevel)

	proto.LastOfflineTime = timeutil.Marshal32(s.LastOfflineTime)

	s.clientProto = proto

	return proto
}

func (s *HeroSnapshot) EncodeClientBytes() []byte {
	if protoBytes := s.getClientProtoBytes(); len(protoBytes) > 0 {
		return protoBytes
	}
	s.clientProtoBytesMux.Lock()
	defer s.clientProtoBytesMux.Unlock()

	return s.encodeClientBytes()
}

func (s *HeroSnapshot) getClientProtoBytes() []byte {
	s.clientProtoBytesMux.RLock()
	defer s.clientProtoBytesMux.RUnlock()
	return s.clientProtoBytes
}

func (s *HeroSnapshot) encodeClientBytes() []byte {
	proto := s.EncodeClient()
	s.clientProtoBytes = util.SafeMarshal(proto)
	return s.clientProtoBytes
}

func (s *HeroSnapshot) EncodeBasic4Client() *shared_proto.HeroBasicProto {
	if proto := s.getBasicProto(); proto != nil {
		return proto
	}

	s.clientBasicProtoMux.Lock()
	defer s.clientBasicProtoMux.Unlock()

	return s.encodeBasic4Client()
}

func (s *HeroSnapshot) getBasicProto() *shared_proto.HeroBasicProto {
	s.clientBasicProtoMux.RLock()
	defer s.clientBasicProtoMux.RUnlock()
	return s.clientBasicProto
}

func (s *HeroSnapshot) encodeBasic4Client() *shared_proto.HeroBasicProto {
	proto := &shared_proto.HeroBasicProto{}

	proto.Id = s.IdBytes
	proto.Name = s.Name
	proto.Head = s.Head
	proto.Body = u64.Int32(s.Body)
	proto.Level = u64.Int32(s.Level)
	proto.VipLevel = u64.Int32(s.VipLevel)
	proto.Male = s.Male
	proto.Location = u64.Int32(s.Location)

	if g := s.Guild(); g != nil {
		proto.GuildId = i64.Int32(g.Id)
		proto.GuildName = g.Name
		proto.GuildFlagName = g.FlagName
	}
	proto.CountryId = u64.Int32(s.CountryId)
	proto.Official = s.countryOfficial

	s.clientBasicProto = proto
	return proto
}

func (s *HeroSnapshot) EncodeBasic4ClientBytes() []byte {
	if protoBytes := s.getBasicProtoBytes(); len(protoBytes) > 0 {
		return protoBytes
	}

	s.clientBasicProtoBytesMux.Lock()
	defer s.clientBasicProtoBytesMux.Unlock()
	return s.encodeBasic4ClientBytes()
}

func (s *HeroSnapshot) getBasicProtoBytes() []byte {
	s.clientBasicProtoBytesMux.RLock()
	defer s.clientBasicProtoBytesMux.RUnlock()
	return s.clientBasicProtoBytes
}

func (s *HeroSnapshot) encodeBasic4ClientBytes() []byte {
	proto := s.EncodeBasic4Client()
	s.clientBasicProtoBytes = util.SafeMarshal(proto)
	return s.clientBasicProtoBytes
}

type SnapshotCallback interface {
	OnHeroSnapshotUpdate(*HeroSnapshot)
}

func NewIdBasicProto(id int64) *shared_proto.HeroBasicProto {
	return &shared_proto.HeroBasicProto{
		Id:    idbytes.ToBytes(id),
		Name:  idbytes.PlayerName(id),
		Level: 1,
	}
}

func (s *HeroSnapshot) GetLevel() uint64 {
	return s.Level
}

func (s *HeroSnapshot) GetVipLevel() uint64 {
	return s.VipLevel
}

func (s *HeroSnapshot) GetName() string {
	return s.Name
}

// tlog
func (s *HeroSnapshot) GetTlogHero() *TlogHero {
	return &TlogHero{
		hero: s,
	}
}

var _ entity.TlogHero = (*TlogHero)(nil)

type TlogHero struct {
	hero *HeroSnapshot
}

func (s *TlogHero) Id() int64                      { return s.hero.Id }
func (s *TlogHero) Level() uint64                  { return s.hero.Level }
func (s *TlogHero) VipLevel() uint64               { return s.hero.VipLevel }
func (s *TlogHero) Name() string                   { return s.hero.Name }
func (s *TlogHero) FriendsCount() uint64           { return s.hero.FriendsCount }
func (s *TlogHero) TotalOnlineTime() time.Duration { return 0 } // TODO 数据这暂时获取不到
func (s *TlogHero) BaseLevel() uint64              { return s.hero.BaseLevel }
