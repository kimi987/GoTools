package data

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"time"
)

//gogen:config
type BuffEffectData struct {
	_ struct{} `file:"杂项/buff.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"buff.proto"`

	Id           uint64
	Name         string                      `desc:"增益名称"`
	Desc         string                      `desc:"描述"`
	Group        uint64                      `desc:"buff 分组"`
	EffectType   shared_proto.BuffEffectType `desc:"buff效果类型"`
	Level        uint64                      `desc:"同类型的等级"`
	KeepDuration time.Duration               `desc:"持续时间"`
	NoDuration   bool                        `desc:"没有持续时间"`
	PvpBuff      bool                        `desc:"pvp buff"`
	StatBuff     *SpriteStat                 `default:"nullable"` // 加buff
	CaptainTrain *Amount                     `desc:"buff 效果，千分比"`
	FarmHarvest  *Amount                     `desc:"buff 效果，千分比"`
	Tax          *Amount                     `desc:"税收，千分比"`
	AdvantageId  uint64                      `desc:"触发它的增益 Id. BufferDataProto.Id" validator:"uint"`
}

func (buffEffectData *BuffEffectData) InitAll(filename string, conf interface {
	GetBuffEffectDataArray() []*BuffEffectData
}) {
	groupMap := make(map[uint64]*BuffEffectData)
	for _, d := range conf.GetBuffEffectDataArray() {
		d.defaultInit()
		prevData, ok := groupMap[d.Group]
		if !ok {
			check.PanicNotTrue(d.Level == 1, "%v, 相同 group 的 level 必须从1开始递增。id:%v level: %v", filename, d.Id, d.Level)
			groupMap[d.Group] = d
			continue
		}

		check.PanicNotTrue(prevData.EffectType == d.EffectType, "%v, 相同 group 的 type 必须相同。id:%v", filename, d.Id)
		check.PanicNotTrue(prevData.Level+1 == d.Level, "%v, 相同 group 的 level 必须从1开始递增。id:%v level: %v preLevel:%v", filename, d.Id, d.Level, prevData.Level)
		groupMap[d.Group] = d

		if !d.NoDuration {
			check.PanicNotTrue(d.KeepDuration > 0, "%v, buff id:%v 的 keep_duration 必须配置。", filename, d.Id)
		}
	}
}

func (d *BuffEffectData) defaultInit() {
	if d.CaptainTrain == nil {
		d.CaptainTrain = NewAmountBuilder().Amount()
	}
	if d.FarmHarvest == nil {
		d.FarmHarvest = NewAmountBuilder().Amount()
	}
	if d.Tax == nil {
		d.Tax = NewAmountBuilder().Amount()
	}
}

////gogen:config
//type StatBuff struct {
//	_ struct{} `proto:"shared_proto.StatBuffProto"`
//
//	Id int `protofield:"-"`
//
//	Attack         *Amount `desc:"sprite_stat, 千分比"`
//	Defense        *Amount `desc:"sprite_stat, 千分比"`
//	Strength       *Amount `desc:"sprite_stat, 千分比"`
//	Dexterity      *Amount `desc:"sprite_stat, 千分比"`
//	SoldierCapcity *Amount `desc:"sprite_stat, 千分比"`
//	DamageIncrePer int     `desc:"sprite_stat, 万分比" validator:"int"`
//	DamageDecrePer int     `desc:"sprite_stat, 万分比" validator:"int"`
//}
//
//func (d *StatBuff) defaultInit() {
//	if d.Attack == nil {
//		d.Attack = NewAmountBuilder().Amount()
//	}
//	if d.Defense == nil {
//		d.Defense = NewAmountBuilder().Amount()
//	}
//	if d.Strength == nil {
//		d.Strength = NewAmountBuilder().Amount()
//	}
//	if d.Dexterity == nil {
//		d.Dexterity = NewAmountBuilder().Amount()
//	}
//	if d.SoldierCapcity == nil {
//		d.SoldierCapcity = NewAmountBuilder().Amount()
//	}
//}
//
//func TotalStatBuff(stats ...*StatBuff) *StatBuff {
//	newBuff := &StatBuff{}
//
//	attackBuilder := NewAmountBuilder()
//	defenseBuilder := NewAmountBuilder()
//	strengthBuilder := NewAmountBuilder()
//	dexterityBuilder := NewAmountBuilder()
//	soldierCapcityBuilder := NewAmountBuilder()
//	var damageIncrePerBuilder int
//	var damageDecrePerBuilder int
//
//	for _, s := range stats {
//		attackBuilder.AddAmount(s.Attack)
//		defenseBuilder.AddAmount(s.Defense)
//		strengthBuilder.AddAmount(s.Strength)
//		defenseBuilder.AddAmount(s.Defense)
//		soldierCapcityBuilder.AddAmount(s.SoldierCapcity)
//		damageIncrePerBuilder += s.DamageIncrePer
//		damageDecrePerBuilder += s.DamageDecrePer
//	}
//
//	newBuff.Attack = attackBuilder.Amount()
//	newBuff.Defense = defenseBuilder.Amount()
//	newBuff.Strength = strengthBuilder.Amount()
//	newBuff.Dexterity = dexterityBuilder.Amount()
//	newBuff.SoldierCapcity = soldierCapcityBuilder.Amount()
//	newBuff.DamageIncrePer = damageIncrePerBuilder
//	newBuff.DamageDecrePer = damageDecrePerBuilder
//
//	return newBuff
//}
//
//func CalcStatBuff(stat *SpriteStat, buff *StatBuff) *SpriteStat {
//	b := &SpriteStatBuilder{}
//	b.attack = buff.Attack.Calc(stat.Attack)
//	b.defense = buff.Defense.Calc(stat.Defense)
//	b.strength = buff.Strength.Calc(stat.Strength)
//	b.dexterity = buff.Dexterity.Calc(stat.Dexterity)
//	b.soldierCapcity = buff.SoldierCapcity.Calc(stat.SoldierCapcity)
//	b.damageIncrePer = buff.DamageIncrePer + stat.DamageIncrePer
//	b.damageDecrePer = buff.DamageDecrePer + stat.DamageDecrePer
//
//	return b.Build()
//}
//
//func CalcStatProtoBuff(stat *shared_proto.SpriteStatProto, buff *StatBuff) *shared_proto.SpriteStatProto {
//	b := &SpriteStatBuilder{}
//	b.attack = buff.Attack.Calc(u64.FromInt32(stat.Attack))
//	b.defense = buff.Defense.Calc(u64.FromInt32(stat.Defense))
//	b.strength = buff.Strength.Calc(u64.FromInt32(stat.Strength))
//	b.dexterity = buff.Dexterity.Calc(u64.FromInt32(stat.Dexterity))
//	b.soldierCapcity = buff.SoldierCapcity.Calc(u64.FromInt32(stat.SoldierCapcity))
//	b.damageIncrePer = buff.DamageIncrePer + int(stat.DamageIncrePer)
//	b.damageDecrePer = buff.DamageDecrePer + int(stat.DamageDecrePer)
//
//	return b.Build().Encode4Init()
//}
