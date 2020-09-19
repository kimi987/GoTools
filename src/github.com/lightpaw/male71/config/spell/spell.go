package spell

import (
	"github.com/lightpaw/male7/config/icon"
	"time"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/domestic_data/sub"
)

// 技能
//gogen:config
type Spell struct {
	_ struct{} `file:"军事/技能.txt"`
	_ struct{} `proto:"shared_proto.SpellProto"`
	_ struct{} `protoconfig:"SpellConfig"`

	Id   uint64     `validator:"int>0"`         // 技能id
	Name string     `validator:"string>0"`      // 名字
	Icon *icon.Icon `protofield:"IconId,%s.Id"` // 图标
	Desc string                                 // 描述
}

// 新版技能
//gogen:config
type SpellFacadeData struct {
	_ struct{} `file:"战斗/技能盒.txt"`
	_ struct{} `protogen:"true"`

	Id      uint64 `validator:"int>0"`    // 技能id
	Name    string `validator:"string>0"` // 名字
	Icon    string `protofield:"IconId"`  // 图标
	Desc    string
	SubDesc string `default:" "`

	// 技能战斗力系数，跟属性一样，分母1W
	FightAmountCoef uint64 `validator:"uint" default:"0" protofield:"-"`

	// 分组
	Group uint64

	// 等级
	Level uint64

	// 技能类型 0-战 1-怒 2-政
	SpellType uint64 `head:"-"`

	// 主动技能
	Spell *SpellData `default:"nullable" protofield:"-"`

	// 被动技能
	PassiveSpell []*PassiveSpellData `protofield:"-"`

	// 内政效果
	BuildingEffect *sub.BuildingEffectData `validator:"uint" default:"nullable" protofield:"-"`
}

func (s *SpellFacadeData) Init(filename string) {

	if s.Spell != nil && s.Spell.RageSpell {
		s.SpellType = 1
	} else if s.BuildingEffect != nil {
		s.SpellType = 2
	}

}

//gogen:config
type PassiveSpellData struct {
	_ struct{} `file:"战斗/被动技能.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	Id uint64

	// 动画特效
	Animation uint64 `validator:"uint"`

	// 触发概率
	TriggerRate float64

	// 触发类型，0-无 1-开局释放 2-首次攻击触发 3-攻击N次触发 4-被攻击触发 5-护盾被击破触发 6-己方获得护盾 7-己方护盾存续期间 8-自身存活期间
	TriggerType shared_proto.SpellTriggerType `validator:"string" white:"0" type:"enum"`

	// 触发目标
	TriggerTarget *SpellTargetData `type:"sub"`

	// 连续攻击N次触发
	TriggerHit uint64 `validator:"uint"`

	// 目标触发CD
	TargetCooldown time.Duration `default:"0s" protofield:",int32(%s / time.Millisecond),int32"`

	// 己方有单位获得护盾时候触发，给自己也上一个

	// 己方护盾存续期间，加攻速 TODO

	// 自身存活，延长状态持续时间 TODO

	// 下面是被动技能效果

	// 给自己加状态
	SelfState []*StateData `protofield:",config.U64a2I32a(GetStateDataKeyArray(%s)),int32"`

	// 给目标加状态
	TargetState []*StateData `protofield:",config.U64a2I32a(GetStateDataKeyArray(%s)),int32"`

	// 触发技能
	Spell *SpellData `default:"nullable" protofield:",config.U64ToI32(%s.Id),int32"`

	// 持续效果立即生效
	ExciteEffectType uint64 `validator:"uint"`

	// 添加额外初始怒气
	Rage uint64 `validator:"uint"`

	// 加属性
	SpriteStat *data.SpriteStat `default:"nullable"`

	// 别人打我，打的更疼
	BeenHurtEffectIncType []uint64
	BeenHurtEffectInc     []float64 `validator:"float64"`

	// 别人打我，打的更轻
	BeenHurtEffectDecType []uint64
	BeenHurtEffectDec     []float64 `validator:"float64"`

	// 原地复活，拥有兵力百分比
	RelivePercent float64 `validator:"float64" default:"0"`

	// -------------

	// 战斗开始，给自己上buff（释放技能）
}

//gogen:config
type SpellData struct {
	_ struct{} `file:"战斗/技能.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"spell.proto"`
	_ struct{} `protoconfig:"-"`

	Id   uint64 `validator:"int>0"`    // 技能id
	Name string `validator:"string>0"` // 名字

	// 动画特效
	Animation uint64 `validator:"uint"`

	Cooldown        time.Duration `protofield:",int32(%s / time.Millisecond),int32"`
	StrongeDuration time.Duration `protofield:",int32(%s / time.Millisecond),int32"`

	RageSpell bool

	KeepMove bool

	// true表示友方技能，否则表示敌方技能
	FriendSpell bool

	SelfAsTarget bool

	// 0-无 1-随机目标 2-血量比例最高 3-血量比例最低 4-怒气最高 5-怒气最低
	TargetSubType shared_proto.SpellTargetSubType `validator:"string" white:"0" type:"enum"`

	// 施法目标
	Target *SpellTargetData `type:"sub"`

	// 施法距离
	ReleaseRange uint64 `validator:"uint"`

	// 伤害范围
	HurtRange uint64 `validator:"uint"`

	// 伤害个数
	HurtCount uint64 `validator:"uint"`

	// 伤害类型 0-伤害公式 1-兵力*系数 2-平均4维伤害公式
	HurtType uint64 `validator:"uint" default:"0"`

	// 技能系数X
	Coef float64 `validator:"float64"`

	// 飞行速度
	FlySpeed uint64 `validator:"uint"`

	//// 延时生效
	//DamageDelay time.Duration `protofield:",int32(%s / time.Millisecond),int32"`

	// 技能效果分类 0-无 1-燃烧 2-中毒 3-流血
	EffectType uint64 `validator:"uint"`

	// 状态列表
	SelfState     []*StateData `protofield:",config.U64a2I32a(GetStateDataKeyArray(%s)),int32"`
	SelfStateRate []float64    `validator:"float64,duplicate"` // 触发概率

	TargetState     []*StateData `protofield:",config.U64a2I32a(GetStateDataKeyArray(%s)),int32"`
	TargetStateRate []float64    `validator:"float64,duplicate"` // 触发概率

	// 技能加怒气
	SelfRage   int `validator:"uint"`
	TargetRage int `validator:"int"`

	// 技能目标 ------------ TODO

	// 根据是否有护盾，添加护盾或者属性 TODO
}

//gogen:config
type StateData struct {
	_ struct{} `file:"战斗/状态.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	Id   uint64
	Name string

	// 动画特效
	Animation uint64 `validator:"uint"`

	// 堆叠类型 0-不堆叠 1-刷新时间 2-堆叠
	StackType shared_proto.StateStackType `validator:"string" white:"0" type:"enum"`

	// 堆叠最大层数
	StackMaxTimes uint64

	// 跳跃次数
	TickTimes uint64

	// 每次跳跃间隔
	TickDuration time.Duration `protofield:",int32(%s / time.Millisecond),int32"`

	// 附加属性
	ChangeStat *data.SpriteStat `default:"nullable"`
	IsAddStat  bool

	// 移动速度
	MoveSpeedRate float64 `validator:"float64"`

	// 攻速
	AttackSpeedRate float64 `validator:"float64"`

	// 护盾
	ShieldRate       float64 `validator:"float64"`
	ShieldEffectRate float64 `validator:"float64" default:"1"`

	// 不可走
	Unmovable bool

	// 不可攻击（普攻）
	NotAttackable bool

	// 沉默（只能普攻）
	Silence bool

	// 晕眩
	Stun bool

	// 状态类型 0-无 1-流血 2-中毒 3-燃烧
	EffectType uint64 `validator:"uint"`

	// 别人打我，打的更疼
	BeenHurtEffectIncType []uint64
	BeenHurtEffectInc     []float64 `validator:"float64"`

	// 别人打我，打的更轻
	BeenHurtEffectDecType []uint64
	BeenHurtEffectDec     []float64 `validator:"float64"`

	// 掉血，伤害系数
	DamageCoef float64 `validator:"float64"`

	// 伤害类型 0-伤害公式 1-兵力*系数 2-平均4维伤害公式
	DamageHurtType uint64 `validator:"uint" default:"0"`

	// 固定加怒气
	Rage int `validator:"int"`

	// 怒气恢复速度buff
	RageRecoverRate float64 `validator:"float64"`

	// 新增状态，有护盾破时候，爆炸伤害周围武将
}

//gogen:config
type SpellTargetData struct {
	_ struct{} `protogen:"true"`
	_ struct{} `protoconfig:"-"`

	// 触发目标职业
	TargetRace []shared_proto.Race `validator:"string" type:"enum"`

	// 触发目标有特定效果（燃烧，中毒)
	TargetEffectType uint64 `validator:"uint"`

	// 目标不可行走时触发
	TargetUnmovable bool

	// 目标不可攻击（普攻）时触发
	TargetNotAttackable bool

	// 目标沉默时触发
	TargetSilence bool

	// 目标晕眩时触发
	TargetStun bool
}
