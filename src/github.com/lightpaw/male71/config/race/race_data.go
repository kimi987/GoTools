package race

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/spell"
	"github.com/lightpaw/male7/pb/shared_proto"
)

//gogen:config
type RaceData struct {
	_  struct{} `file:"武将/职业.txt"`
	_  struct{} `proto:"shared_proto.RaceDataProto"`
	_  struct{} `protoconfig:"race_data"`
	Id int

	Race shared_proto.Race `type:"enum"`
	Name string            `default:"name" protofield:"-"`

	IsFar bool `default:"false"` // 是否是远程

	AttackRange       uint64
	MoveTimesPerRound uint64
	MoveSpeed         uint64
	ViewRange         uint64

	// 攻击优先级
	Priority []shared_proto.Race `type:"enum" validator:"string,count=5,notNil"`

	// 兵种克制系数
	RaceCoef []uint64 `validator:"int>0,count=5,notNil,duplicate"`
	WallCoef uint64

	// 四维属性比例
	AbilityRate []float64 `validator:"float64>0,count=4,notNil,duplicate,sum=1"`

	RestraintRace                   []shared_proto.Race                                                                           // 克制的职业
	RestraintRoundType              shared_proto.RestraintRoundType                                                               // 触发克制技能的轮数的类型
	RestraintSpell                  *spell.Spell `head:"restraint_spell_id" protofield:"RestraintSpellId,config.U64ToI32(%s.Id)"` // 克制技
	UnlockRestraintSpellNeedAbility uint64       `validator:"uint"`                                                               // 解锁克制技能需要的成长值

	// 普通技能
	NormalSpell *spell.Spell `head:"normal_spell_id" protofield:"NormalSpellId,config.U64ToI32(%s.Id)"`

	// 兵种技能
	SoldierSpell []*spell.SpellFacadeData `protofield:",config.U64a2I32a(spell.GetSpellFacadeDataKeyArray(%s))"`

	// 兵种攻速
	SoldierAttackSpeed float64 `default:"5"`

	// 兵种攻击距离
	SoldierAttackRange uint64 `default:"30"`

	GemTypes []uint64 `validator:"int>0,count=9,notNil,duplicate"` // 实际装配宝石的槽位类型列表（新版镶嵌）

	proto *shared_proto.RaceDataProto
}

func (r *RaceData) Init() {
	if int(r.Race) != r.Id {
		logrus.Panicf("加载RaceData，职业ID与职业名称对应不上，职业[%s] 的ID应该是 %d, 实际配置的是 %d", r.Race, int(r.Race), r.Id)
	}

	var i interface{}
	i = r
	m, ok := i.(interface {
		Encode() *shared_proto.RaceDataProto
	})
	if !ok {
		logrus.Panicf("RaceData.Encode4Init() cast type fail")
	}

	r.proto = m.Encode()
}

func (r *RaceData) GetProto() *shared_proto.RaceDataProto {
	return r.proto
}

func (r *RaceData) GetAbilityRate(statType shared_proto.StatType) float64 {
	index := int(statType) - 1
	if index >= 0 && index < len(r.AbilityRate) {
		return r.AbilityRate[index]
	}
	return 0
}
