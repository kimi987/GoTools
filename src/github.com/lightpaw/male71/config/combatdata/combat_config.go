package combatdata

import (
	"time"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/spell"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/check"
)

//gogen:config
type CombatMiscConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"战斗/杂项.txt"`
	_ struct{} `protogen:"true"`

	Spell          []uint64 `head:"-"`
	SpellAnimation []uint64 `head:"-"`

	PassiveSpell          []uint64 `head:"-"`
	PassiveSpellAnimation []uint64 `head:"-"`

	State          []uint64 `head:"-"`
	StateAnimation []uint64 `head:"-"`

	// 最大怒气值
	MaxRage uint64

	// 城墙攻速，这个值要除以1000，5600表示5.6
	WallAttackSpeed float64 `head:"-"`
}

func (s *CombatMiscConfig) Init(filename string, configs interface {
	GetSpellDataArray() []*spell.SpellData
	GetPassiveSpellDataArray() []*spell.PassiveSpellData
	GetStateDataArray() []*spell.StateData
}) {

	idMap := make(map[uint64]int)
	for _, data := range configs.GetSpellDataArray() {
		check.PanicNotTrue(idMap[data.Id] == 0, "技能id冲突，id: %d 类型（1-主动技能 2-被动技能 3-状态） %d %d", data.Id, idMap[data.Id], 1)
		idMap[data.Id] = 1
	}
	for _, data := range configs.GetPassiveSpellDataArray() {
		check.PanicNotTrue(idMap[data.Id] == 0, "技能id冲突，id: %d 类型（1-主动技能 2-被动技能 3-状态） %d %d", data.Id, idMap[data.Id], 2)
		idMap[data.Id] = 2
	}
	for _, data := range configs.GetStateDataArray() {
		check.PanicNotTrue(idMap[data.Id] == 0, "技能id冲突，id: %d 类型（1-主动技能 2-被动技能 3-状态） %d %d", data.Id, idMap[data.Id], 3)
		idMap[data.Id] = 3
	}

	s.Spell = nil
	s.SpellAnimation = nil
	for _, data := range configs.GetSpellDataArray() {
		if data.Animation > 0 {
			s.Spell = append(s.Spell, data.Id)
			s.SpellAnimation = append(s.SpellAnimation, data.Animation)
		}
	}

	s.PassiveSpell = nil
	s.PassiveSpellAnimation = nil
	for _, data := range configs.GetPassiveSpellDataArray() {
		if data.Animation > 0 {
			s.PassiveSpell = append(s.PassiveSpell, data.Id)
			s.PassiveSpellAnimation = append(s.PassiveSpellAnimation, data.Animation)
		}
	}

	s.State = nil
	s.StateAnimation = nil
	for _, data := range configs.GetStateDataArray() {
		if data.Animation > 0 {
			s.State = append(s.State, data.Id)
			s.StateAnimation = append(s.StateAnimation, data.Animation)
		}
	}
}

//gogen:config
type CombatConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"战斗/杂项.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	// 技能数据
	Spell []*spell.SpellData `head:"-"`

	// 状态数据
	State []*spell.StateData `head:"-"`

	// 被动技能
	PassiveSpell []*spell.PassiveSpellData `head:"-"`

	// 技能id映射
	SpellIdMap []*shared_proto.SpellIdProto `head:"-" protofield:",%s,SpellIdProto"`

	// 武将数据
	Captain []*captain.CaptainData `head:"-"`

	NamelessCaptain []*captain.NamelessCaptainData `head:"-"`

	// 职业数据
	Race []*race.RaceData `head:"-"`

	// 每秒帧数
	FramePerSecond uint64 `default:"10"`

	// 系数的分母
	ConfigDenominator uint64 `default:"1000"`

	// 最小攻速
	MinAttackDuration time.Duration `protofield:",int32(%s / time.Millisecond),int32"`

	// 最大攻速
	MaxAttackDuration time.Duration `protofield:",int32(%s / time.Millisecond),int32"`

	// 最小移动速度（每秒）
	MinMoveSpeed uint64

	// 最大移动速度（每秒）
	MaxMoveSpeed uint64

	// 最小战斗属性（避免属性被减到负数）
	MinStat *data.SpriteStat

	// 每场战斗持续时间
	MaxDuration time.Duration

	// 战斗结果分数比例
	ScorePercent []uint64 `head:"-"`

	CheckMoveDuration time.Duration `default:"500ms" protofield:",int32(%s / time.Millisecond),int32"`

	CritRate float64

	Coef float64

	CellLen uint64 `default:"100"`

	// 最大怒气值
	MaxRage uint64

	// 每次命中增加的怒气值
	AddRagePerHint uint64

	// 每损失1%的血量增加的怒气值
	AddRageLost1Percent uint64

	// 怒气恢复速度（每秒）
	RageRecoverSpeed uint64

	// 地图大小
	CombatXLen uint64 `default:"1000" protofield:"-"`
	CombatYLen uint64 `default:"500" protofield:"-"`

	InitAttackerX uint64 `validator:"uint" default:"0"`
	InitDefenserX uint64 `validator:"uint" default:"1000"`
	InitWallX     uint64 `validator:"uint" default:"1200"`

	// 城墙攻击等待间隔（过了这么久之后才能攻击）
	WallWaitDuration time.Duration `defalut:"3s" protofield:",int32(%s / time.Millisecond),int32"`

	// 城墙附加固定伤害次数
	WallAttackFixDamageTimes uint64 `default:"3"`

	// 每次城墙受到攻击损失最大值
	WallBeenHurtLostMaxPercent float64 `default:"0"`

	WallSpell *spell.SpellData `protofield:",config.U64ToI32(%s.Id),int32"`

	// 城墙攻击最小飞行时间
	WallFlyMinDuration time.Duration `defalut:"500ms" protofield:",int32(%s / time.Millisecond),int32"`

	// 短移动距离
	ShortMoveDistance uint64 `validator:"uint" default:"100"`
}

func (c *CombatConfig) Init(filename string, configs interface {
	GetSpellDataArray() []*spell.SpellData
	GetStateDataArray() []*spell.StateData
	GetPassiveSpellDataArray() []*spell.PassiveSpellData
	GetRaceDataArray() []*race.RaceData
	GetCaptainDataArray() []*captain.CaptainData
	GetNamelessCaptainDataArray() []*captain.NamelessCaptainData
	MilitaryConfig() *singleton.MilitaryConfig
	GetSpellFacadeDataArray() []*spell.SpellFacadeData
	CombatMiscConfig() *CombatMiscConfig
}) {
	c.Spell = configs.GetSpellDataArray()
	c.State = configs.GetStateDataArray()
	c.PassiveSpell = configs.GetPassiveSpellDataArray()

	c.Race = configs.GetRaceDataArray()

	c.Captain = configs.GetCaptainDataArray()
	c.NamelessCaptain = configs.GetNamelessCaptainDataArray()

	c.ScorePercent = configs.MilitaryConfig().CombatScorePercent

	var spellIdMap []*shared_proto.SpellIdProto
	for _, s := range configs.GetSpellFacadeDataArray() {
		proto := &shared_proto.SpellIdProto{
			Id: u64.Int32(s.Id),
		}

		if s.Spell != nil {
			proto.Spell = u64.Int32(s.Spell.Id)
		}

		if len(s.PassiveSpell) > 0 {
			for _, ps := range s.PassiveSpell {
				proto.PassiveSpell = append(proto.PassiveSpell, u64.Int32(ps.Id))
			}
		}

		spellIdMap = append(spellIdMap, proto)
	}
	c.SpellIdMap = spellIdMap

	if c.WallSpell != nil {
		configs.CombatMiscConfig().WallAttackSpeed = c.WallSpell.Cooldown.Seconds()
	}
}
