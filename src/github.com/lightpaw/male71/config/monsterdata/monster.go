package monsterdata

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/config/captain"
)

type SpecType uint8

const (
	None         SpecType = 0
	InvasionTask SpecType = 1
)

//gogen:config
type MonsterMasterData struct {
	_  struct{} `file:"怪物/怪物君主.txt"`
	_  struct{} `proto:"shared_proto.MonsterMasterDataProto"`
	Id uint64

	Name string

	Icon *icon.Icon `protofield:"IconId,%s.Id"` // 图标

	Male bool

	Level uint64

	Captains    []*MonsterCaptainData
	FightAmount uint64 `head:"-" protofield:"-"`

	WallLevel     uint64           `default:"1" protofield:"-"`        // 城墙等级
	WallStat      *data.SpriteStat `default:"nullable" protofield:"-"` // 城墙属性
	WallFixDamage uint64           `validator:"uint" protofield:"-"`   // 城墙固定伤害

	// 驱逐奖励
	InvadePrize *resdata.PlunderPrize `default:"nullable" protofield:"-"`

	// 击杀奖励
	BeenKillPrize *resdata.PlunderPrize `default:"nullable" protofield:"-"`

	player *shared_proto.CombatPlayerProto `head:"-" protofield:"-"`

	spec SpecType

	npcId int64
}

func (m *MonsterMasterData) Init(filename string, configs interface{
	GetMonsterCaptainDataArray() []*MonsterCaptainData
}) {
	m.npcId = npcid.NewMonsterId(m.Id)
	m.FightAmount = m.CalculateFightAmount()

	var idxs []uint64
	for i, c := range m.Captains {
		idx := c.Index
		if c.Index == 0 {
			idx = uint64(i + 1)
		}

		check.PanicNotTrue(!u64.Contains(idxs, idx), "%s %s 配置的武将顺序重叠，", m.Id, m.Name, idx)
		idxs = append(idxs, c.Index)
	}

	m.player = m.newPlayer()
}

func (m *MonsterMasterData) InitSpec(toSet SpecType) {
	check.PanicNotTrue(m.spec == None, "%s %s 初始化特殊类型时候，发现已经设置过特殊类型了, m.spec=%v toSet=%v", m.Id, m.Name, m.spec, toSet)
	m.spec = toSet
}

func (m *MonsterMasterData) GetSpec() SpecType {
	return m.spec
}

func (m *MonsterMasterData) GetNpcId() int64 {
	return m.npcId
}

func (m *MonsterMasterData) CalculateFightAmount() uint64 {
	tfa := data.NewTroopFightAmount()
	for _, c := range m.Captains {
		tfa.Add(c.calculateFightAmount())
	}

	return tfa.ToU64()
}

func (m *MonsterMasterData) EncodeSnapshot(id int64) *shared_proto.HeroBasicSnapshotProto {
	s := &shared_proto.HeroBasicSnapshotProto{}

	s.Basic = m.EncodeHeroBasicProto(idbytes.ToBytes(id))
	s.FightAmount = u64.Int32(m.FightAmount)

	return s
}

func (m *MonsterMasterData) EncodeHeroBasicProto(id []byte) *shared_proto.HeroBasicProto {
	proto := &shared_proto.HeroBasicProto{}

	proto.Id = id
	proto.Name = m.Name
	proto.Head = m.Icon.Id
	proto.Level = u64.Int32(m.Level)
	proto.Male = m.Male

	return proto
}

func (m *MonsterMasterData) GetPlayer() *shared_proto.CombatPlayerProto {
	return m.player
}

func (m *MonsterMasterData) newPlayer() *shared_proto.CombatPlayerProto {
	player := &shared_proto.CombatPlayerProto{}
	player.Hero = &shared_proto.HeroBasicProto{}
	player.Hero.Id = idbytes.ToBytes(npcid.NewMonsterId(m.Id))
	player.Hero.Name = m.Name
	player.Hero.Head = m.Icon.Id
	player.Hero.Male = m.Male
	player.Hero.Level = u64.Int32(m.Level)

	tfa := data.NewTroopFightAmount()
	for i, c := range m.Captains {
		proto := &shared_proto.CombatTroopsProto{}
		if c.Index > 0 && c.Index <= 5 {
			proto.FightIndex = u64.Int32(c.Index)
		} else {
			proto.FightIndex = imath.Int32(i + 1)
		}
		proto.Captain = c.EncodeCaptainInfo()

		player.Troops = append(player.Troops, proto)

		tfa.AddInt32(proto.Captain.FightAmount)
	}
	player.TotalFightAmount = tfa.ToI32()

	if m.WallStat != nil {
		player.WallStat = m.WallStat.Encode4Init()
		player.WallFixDamage = u64.Int32(m.WallFixDamage)
		player.TotalWallLife = i32.Max(player.WallStat.Strength, 1)
		player.WallLevel = u64.Int32(m.WallLevel)
	}

	return player
}

//gogen:config
type MonsterCaptainData struct {
	_  struct{} `file:"怪物/怪物武将.txt"`
	_  struct{} `proto:"shared_proto.MonsterCaptainDataProto"`
	Id uint64

	Captain *captain.CaptainData `default:"nullable" protofield:"-"`

	Star uint64 `validator:"uint"`

	NamelessCaptain *captain.NamelessCaptainData `default:"nullable" protofield:"-"`

	CaptainId uint64 `head:"-"`

	IsNameless bool `head:"-"`

	UnlockSpellCount uint64 `validator:"uint"`

	Quality shared_proto.Quality `type:"enum"`

	Soldier uint64

	FightAmount uint64 `head:"-"`

	Morale uint64 `protofield:"-"`

	Level uint64

	SoldierLevel uint64

	RebirthLevel uint64 `validator:"uint" default:"0"`

	TotalStat *data.SpriteStat

	Label uint64 `validator:"uint"`

	Model uint64 `default:"1" protofield:"-"`

	// 是否可以触发克制技，默认不触发
	CanTriggerRestraintSpell bool `default:"false" protofield:"-"`

	Index uint64 `validator:"uint" default:"0"`

	XIndex uint64 `validator:"uint" default:"0"`

	YuanJun bool `default:"false" protofield:"-"`
}

func (c *MonsterCaptainData) Init(filename string) {
	check.PanicNotTrue(c.Id > 1000000, "%v id 必须 >1000000. id:%v", filename, c.Id)

	c.FightAmount = c.calculateFightAmount()

	hasCaptain := c.Captain != nil
	hasNamelessCaptain := c.NamelessCaptain != nil
	check.PanicNotTrue(hasCaptain != hasNamelessCaptain, "%s 配置武将和无名武将必须配置，并且只能配置一个，不能同时配置", filename)

	c.CaptainId = c.getCaptainId()
	c.IsNameless = hasNamelessCaptain

}

func (c *MonsterCaptainData) getCaptainId() uint64 {
	if c.Captain != nil {
		return c.Captain.Id
	} else if c.NamelessCaptain != nil {
		return c.NamelessCaptain.Id
	}
	return 0
}

func (c *MonsterCaptainData) getRace() shared_proto.Race {
	if c.Captain != nil {
		return c.Captain.Race.Race
	} else if c.NamelessCaptain != nil {
		return c.NamelessCaptain.Race
	}
	return 0
}

func (c *MonsterCaptainData) getSpellFightAmountCoef() uint64 {
	if c.Captain != nil {
		return c.Captain.GetStar(c.Star).GetSpellFightAmountCoef(c.UnlockSpellCount)
	} else if c.NamelessCaptain != nil {
		return c.NamelessCaptain.GetSpellFightAmountCoef(c.UnlockSpellCount)
	}
	return 0
}

func (c *MonsterCaptainData) calculateFightAmount() uint64 {
	return c.TotalStat.FightAmount(c.Soldier, c.getSpellFightAmountCoef())
}

// 因为很多地方会直接使用CaptainInfoProto，改变里面的soldier值，所以不要做缓存
func (c *MonsterCaptainData) EncodeCaptainInfo() *shared_proto.CaptainInfoProto {
	out := &shared_proto.CaptainInfoProto{}

	out.Id = u64.Int32(c.Id)

	out.Race = c.getRace()
	out.Quality = c.Quality
	out.TotalSoldier = u64.Int32(c.Soldier)
	out.Soldier = u64.Int32(c.Soldier)
	out.SpellFightAmountCoef = u64.Int32(c.getSpellFightAmountCoef())
	out.FightAmount = u64.Int32(c.calculateFightAmount())

	out.CanTriggerRestraintSpell = c.CanTriggerRestraintSpell

	out.Morale = u64.Int32(c.Morale)
	out.Level = u64.Int32(c.Level)
	out.SoldierLevel = u64.Int32(c.SoldierLevel)

	out.TotalStat = c.TotalStat.Encode4Init()
	out.LifePerSoldier = u64.Int32(c.TotalStat.Life())

	out.Model = u64.Int32(c.Model)

	out.RebirthLevel = u64.Int32(c.RebirthLevel)

	out.YuanJun = c.YuanJun

	out.XIndex = u64.Int32(c.XIndex)

	out.CaptainId = u64.Int32(c.getCaptainId())
	if c.NamelessCaptain != nil {
		out.IsNameless = true
	}

	out.Star = u64.Int32(c.Star)
	out.UnlockSpellCount = u64.Int32(c.UnlockSpellCount)

	return out
}
