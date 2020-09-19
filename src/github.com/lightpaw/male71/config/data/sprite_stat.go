package data

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"math"
)

var (
	emptyStat      = &SpriteStat{}
	emptyStatProto = &shared_proto.SpriteStatProto{}
)

func EmptyStat() *SpriteStat {
	return emptyStat
}

func Nil2EmptyStat(s *SpriteStat) *SpriteStat {
	if s != nil {
		return s
	}

	return emptyStat
}

func EmptyStatProto() *shared_proto.SpriteStatProto {
	return emptyStatProto
}

// 职业属性
func NewRaceStat(race shared_proto.Race, spriteStat *SpriteStat) *RaceStat {
	return &RaceStat{
		Race:       race,
		SpriteStat: spriteStat,
	}
}

type RaceStat struct {
	Race       shared_proto.Race
	SpriteStat *SpriteStat
}

//gogen:config
type SpriteStat struct {
	_                struct{} `file:"杂项/属性.txt"`
	_                struct{} `proto:"shared_proto.SpriteStatProto"`
	Id               uint64   `protofield:"-"`
	Attack           uint64   `validator:"uint"` // 攻击
	Defense          uint64   `validator:"uint"` // 防御
	Strength         uint64   `validator:"uint"` // 体力
	Dexterity        uint64   `validator:"uint"` // 敏捷
	SoldierCapcity   uint64   `validator:"uint"` // 统帅
	DamageIncrePer   uint64   `validator:"int"`  // 伤害增加万分比
	DamageDecrePer   uint64   `validator:"int"`  // 伤害减少万分比
	BeenHurtIncrePer uint64   `validator:"int"`  // 易伤增加万分比
	BeenHurtDecrePer uint64   `validator:"int"`  // 易伤减少万分比
}

func (s *SpriteStat) IsEmpty() bool {
	return s.Attack|s.Defense|s.Strength|s.Dexterity|s.SoldierCapcity |
		s.DamageIncrePer | s.DamageDecrePer | s.BeenHurtIncrePer | s.BeenHurtDecrePer == 0
}

func IsEqualsStat(a, b *SpriteStat) bool {
	if a == b {
		return true
	}

	if a == nil || b == nil {
		return false
	}
	return a.Attack == b.Attack &&
		a.Defense == b.Defense &&
		a.Strength == b.Strength &&
		a.Dexterity == b.Dexterity &&
		a.SoldierCapcity == b.SoldierCapcity &&
		a.DamageIncrePer == b.DamageIncrePer &&
		a.DamageDecrePer == b.DamageDecrePer &&
		a.BeenHurtIncrePer == b.BeenHurtIncrePer &&
		a.BeenHurtDecrePer == b.BeenHurtDecrePer
}

func (s *SpriteStat) Encode4Init() *shared_proto.SpriteStatProto {
	var i interface{}
	i = s

	m, ok := i.(interface {
		Encode() *shared_proto.SpriteStatProto
	})
	if !ok {
		logrus.Panicf("SpriteStat.Encode4Init() cast type fail")
	}

	return m.Encode()
}

func (s *SpriteStat) Sum4D() uint64 {
	return s.Attack + s.Defense + s.Strength + s.Dexterity
}

func (s *SpriteStat) FightAmount(soldier, spellCoef uint64) uint64 {
	return CalFightAmount(s.Attack, s.Defense, s.Strength, s.Dexterity, s.DamageIncrePer, s.DamageDecrePer, s.BeenHurtIncrePer, s.BeenHurtDecrePer, soldier, spellCoef)
}

func CalRate4DStat(point, act, def, str, dex uint64) (newAct, newDef, newStr, newDex uint64) {

	total := act + def + str + dex
	if total <= 0 {
		return
	}

	newAct = point * act / total
	newDef = point * def / total
	newStr = point * str / total
	newDex = point * dex / total

	newTotal := newAct + newDef + newStr + newDex
	if newTotal < point {
		n := u64.Sub(point, newTotal)
		toAdd := uint64(0)
		for i := toAdd; i < n; i++ {
			if act != 0 {
				newAct++

				toAdd++
				if toAdd >= n {
					break
				}
			}

			if def != 0 {
				newDef++

				toAdd++
				if toAdd >= n {
					break
				}
			}

			if str != 0 {
				newStr++

				toAdd++
				if toAdd >= n {
					break
				}
			}

			if dex != 0 {
				newDex++

				toAdd++
				if toAdd >= n {
					break
				}
			}
		}
	}
	return
}

func ProtoFightAmount(s *shared_proto.SpriteStatProto, soldier, spellCoef int32) int32 {
	return u64.Int32(CalFightAmount(
		u64.FromInt32(s.Attack),
		u64.FromInt32(s.Defense),
		u64.FromInt32(s.Strength),
		u64.FromInt32(s.Dexterity),
		u64.FromInt32(s.DamageIncrePer),
		u64.FromInt32(s.DamageDecrePer),
		u64.FromInt32(s.BeenHurtIncrePer),
		u64.FromInt32(s.BeenHurtDecrePer),
		u64.FromInt32(soldier),
		u64.FromInt32(spellCoef),
	))
}

func CalFightAmount(attack, defense, strength, dexterity, damageIncrePer, damageDecrePer, beenHurtIncrePer, beenHurtDecrePer, soldier, spellCoef uint64) uint64 {
	if soldier <= 0 {
		return 0
	}

	sum4D := attack + defense + strength + dexterity

	//S=（（1+加伤1+加伤2…）（1+减伤1+减伤2+…）（1+武将技能参数1+武将技能参数2…）/（1+伤害减弱1+伤害减弱2…）/（1+伤害加深1+伤害加深2…））^0.5
	sc := math.Sqrt(
		(float64(damageIncrePer+constants.Iw) / constants.Iw) *
			(float64(damageDecrePer+constants.Iw) / constants.Iw) *
			(float64(spellCoef+constants.Iw) / constants.Iw) /
			(float64(beenHurtIncrePer+constants.Iw) / constants.Iw) /
			(float64(beenHurtDecrePer+constants.Iw) / constants.Iw))
	scbr := math.Cbrt(float64(soldier) * sc)

	// FP=W * (N*S)^(1/3)
	return u64.MultiF64(sum4D, scbr)
}

func NewTroopFightAmount() troop_fight_amounts {
	return nil
}

type troop_fight_amounts []uint64

func (tfa *troop_fight_amounts) Add(amt uint64) {
	*tfa = append(*tfa, amt)
}

func (tfa *troop_fight_amounts) AddInt32(amt int32) {
	*tfa = append(*tfa, u64.FromInt32(amt))
}

func (tfa *troop_fight_amounts) ToU64() uint64 {
	return TroopFightAmount(*tfa...)
}

func (tfa *troop_fight_amounts) ToI32() int32 {
	return u64.Int32(TroopFightAmount(*tfa...))
}

func TroopFightAmount(captainFightAmount ...uint64) uint64 {
	//AFP=((FP1/10000)^3+(FP2/10000)^3+(FP3/10000)^3+(FP4/10000)^3+(FP5/10000)^3)^(1/3)*10000

	n := len(captainFightAmount)
	switch n {
	case 0:
		return 0
	case 1:
		return captainFightAmount[0]
	}

	var total float64
	for _, amt := range captainFightAmount {
		if amt > 0 {
			total += math.Pow(float64(amt)/10000, 3)
		}
	}

	if total <= 0 {
		return 0
	}

	return uint64(math.Pow(total, 1.0/3.0) * 10000)
}

// 计算单兵血量
// 单兵血量 = 体力 * 5
func (s *SpriteStat) Life() uint64 {
	return getLifePerSoldier(s.Strength)
}

func ProtoLife(proto *shared_proto.SpriteStatProto) uint64 {
	return getLifePerSoldier(u64.FromInt32(proto.Strength))
}

func getLifePerSoldier(strength uint64) uint64 {
	return u64.Max(strength*5, 1)
}

func New4DStat(Attack uint64, Defense uint64, Strength uint64, Dexterity uint64) *SpriteStat {
	ss := &SpriteStat{}
	ss.Attack = Attack
	ss.Defense = Defense
	ss.Strength = Strength
	ss.Dexterity = Dexterity

	return ss
}

func New4DStatProto(Attack uint64, Defense uint64, Strength uint64, Dexterity uint64) *shared_proto.SpriteStatProto {
	ss := &shared_proto.SpriteStatProto{}
	ss.Attack = u64.Int32(Attack)
	ss.Defense = u64.Int32(Defense)
	ss.Strength = u64.Int32(Strength)
	ss.Dexterity = u64.Int32(Dexterity)

	return ss
}

func AppendSpriteStat(appends ...*SpriteStat) *SpriteStat {
	b := &SpriteStatBuilder{}
	for _, a := range appends {
		if a != nil {
			b.Add(a)
		}
	}

	return b.Build()
}

func AppendSpriteStatProto(appends ...*shared_proto.SpriteStatProto) *shared_proto.SpriteStatProto {
	b := &SpriteStatBuilder{}
	for _, a := range appends {
		if a != nil {
			b.AddProto(a)
		}
	}

	return b.Build().Encode4Init()
}

func DiffStat(t, s *SpriteStat) *SpriteStat {
	b := &SpriteStatBuilder{}
	b.Add(t)
	b.Sub(s)
	return b.Build()
}

func NewSpriteStatBuilder() *SpriteStatBuilder {
	return &SpriteStatBuilder{}
}

type SpriteStatBuilder struct {
	attack           uint64
	defense          uint64
	strength         uint64
	dexterity        uint64
	soldierCapcity   uint64
	damageIncrePer   uint64
	damageDecrePer   uint64
	beenHurtIncrePer uint64
	beenHurtDecrePer uint64
}

func (b *SpriteStatBuilder) Build() *SpriteStat {
	ss := &SpriteStat{}
	ss.Attack = b.attack
	ss.Defense = b.defense
	ss.Strength = b.strength
	ss.Dexterity = b.dexterity
	ss.SoldierCapcity = b.soldierCapcity

	ss.DamageIncrePer = b.damageIncrePer
	ss.DamageDecrePer = b.damageDecrePer

	return ss
}

func (b *SpriteStatBuilder) Add(toAdd *SpriteStat) {
	if toAdd == nil {
		return
	}

	b.AddAttack(toAdd.Attack)
	b.AddDefense(toAdd.Defense)
	b.AddStrength(toAdd.Strength)
	b.AddDexterity(toAdd.Dexterity)
	b.AddSoldierCapcity(toAdd.SoldierCapcity)
	b.AddDamageIncrePer(toAdd.DamageIncrePer)
	b.AddDamageDecrePer(toAdd.DamageDecrePer)
	b.AddBeenHurtIncrePer(toAdd.BeenHurtIncrePer)
	b.AddBeenHurtDecrePer(toAdd.BeenHurtDecrePer)
}

func (b *SpriteStatBuilder) Sub(toSub *SpriteStat) {
	if toSub == nil {
		return
	}

	b.attack = u64.Sub(b.attack, toSub.Attack)
	b.defense = u64.Sub(b.defense, toSub.Defense)
	b.strength = u64.Sub(b.strength, toSub.Strength)
	b.dexterity = u64.Sub(b.dexterity, toSub.Dexterity)
	b.soldierCapcity = u64.Sub(b.soldierCapcity, toSub.SoldierCapcity)
	b.damageIncrePer = u64.Sub(b.damageIncrePer, toSub.DamageIncrePer)
	b.damageDecrePer = u64.Sub(b.damageDecrePer, toSub.DamageDecrePer)
	b.beenHurtIncrePer = u64.Sub(b.beenHurtIncrePer, toSub.BeenHurtIncrePer)
	b.beenHurtDecrePer = u64.Sub(b.beenHurtDecrePer, toSub.BeenHurtDecrePer)
}

func (b *SpriteStatBuilder) AddProto(toAdd *shared_proto.SpriteStatProto) {
	if toAdd == nil {
		return
	}

	b.AddAttack(u64.FromInt32(toAdd.Attack))
	b.AddDefense(u64.FromInt32(toAdd.Defense))
	b.AddStrength(u64.FromInt32(toAdd.Strength))
	b.AddDexterity(u64.FromInt32(toAdd.Dexterity))
	b.AddSoldierCapcity(u64.FromInt32(toAdd.SoldierCapcity))
	b.AddDamageIncrePer(u64.FromInt32(toAdd.DamageIncrePer))
	b.AddDamageDecrePer(u64.FromInt32(toAdd.DamageDecrePer))
	b.AddBeenHurtIncrePer(u64.FromInt32(toAdd.BeenHurtIncrePer))
	b.AddBeenHurtDecrePer(u64.FromInt32(toAdd.BeenHurtDecrePer))
}

func (b *SpriteStatBuilder) Add4D(attack, defense, strength, dexterity uint64) {
	b.AddAttack(attack)
	b.AddDefense(defense)
	b.AddStrength(strength)
	b.AddDexterity(dexterity)
}

func (b *SpriteStatBuilder) AddAttack(toAdd uint64) {
	b.attack += toAdd
}

func (b *SpriteStatBuilder) AddDefense(toAdd uint64) {
	b.defense += toAdd
}

func (b *SpriteStatBuilder) AddStrength(toAdd uint64) {
	b.strength += toAdd
}

func (b *SpriteStatBuilder) AddDexterity(toAdd uint64) {
	b.dexterity += toAdd
}

func (b *SpriteStatBuilder) AddSoldierCapcity(toAdd uint64) {
	b.soldierCapcity += toAdd
}

func (b *SpriteStatBuilder) AddDamageIncrePer(toAdd uint64) {
	b.damageIncrePer += toAdd
}

func (b *SpriteStatBuilder) AddDamageDecrePer(toAdd uint64) {
	b.damageDecrePer += toAdd
}

func (b *SpriteStatBuilder) AddBeenHurtIncrePer(toAdd uint64) {
	b.beenHurtIncrePer += toAdd
}

func (b *SpriteStatBuilder) AddBeenHurtDecrePer(toAdd uint64) {
	b.beenHurtDecrePer += toAdd
}
