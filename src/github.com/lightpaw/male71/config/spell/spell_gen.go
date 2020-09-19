// AUTO_GEN, DONT MODIFY!!!
package spell

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

var _ = strings.ToUpper("")      // import strings
var _ = strconv.IntSize          // import strconv
var _ = shared_proto.Int32Pair{} // import shared_proto
var _ = errors.Errorf("")        // import errors
var _ = time.Second              // import time

// start with PassiveSpellData ----------------------------------

func LoadPassiveSpellData(gos *config.GameObjects) (map[uint64]*PassiveSpellData, map[*PassiveSpellData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.PassiveSpellDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*PassiveSpellData, len(lIsT))
	pArSeRmAp := make(map[*PassiveSpellData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrPassiveSpellData) {
			continue
		}

		dAtA, err := NewPassiveSpellData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Id
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Id], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedPassiveSpellData(dAtAmAp map[*PassiveSpellData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.PassiveSpellDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetPassiveSpellDataKeyArray(datas []*PassiveSpellData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewPassiveSpellData(fIlEnAmE string, pArSeR *config.ObjectParser) (*PassiveSpellData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrPassiveSpellData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &PassiveSpellData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Animation = pArSeR.Uint64("animation")
	dAtA.TriggerRate = pArSeR.Float64("trigger_rate")
	dAtA.TriggerType = shared_proto.SpellTriggerType(shared_proto.SpellTriggerType_value[strings.ToUpper(pArSeR.String("trigger_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("trigger_type"), 10, 32); err == nil {
		dAtA.TriggerType = shared_proto.SpellTriggerType(i)
	}

	dAtA.TriggerTarget, err = NewSpellTargetData(fIlEnAmE, pArSeR)
	if err != nil {
		return nil, err
	}
	dAtA.TriggerHit = pArSeR.Uint64("trigger_hit")
	if pArSeR.KeyExist("target_cooldown") {
		dAtA.TargetCooldown, err = config.ParseDuration(pArSeR.String("target_cooldown"))
	} else {
		dAtA.TargetCooldown, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[target_cooldown] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("target_cooldown"), dAtA)
	}

	// releated field: SelfState
	// releated field: TargetState
	// releated field: Spell
	dAtA.ExciteEffectType = pArSeR.Uint64("excite_effect_type")
	dAtA.Rage = pArSeR.Uint64("rage")
	// releated field: SpriteStat
	dAtA.BeenHurtEffectIncType = pArSeR.Uint64Array("been_hurt_effect_inc_type", "", false)
	dAtA.BeenHurtEffectInc = pArSeR.Float64Array("been_hurt_effect_inc", "", false)
	dAtA.BeenHurtEffectDecType = pArSeR.Uint64Array("been_hurt_effect_dec_type", "", false)
	dAtA.BeenHurtEffectDec = pArSeR.Float64Array("been_hurt_effect_dec", "", false)
	dAtA.RelivePercent = 0
	if pArSeR.KeyExist("relive_percent") {
		dAtA.RelivePercent = pArSeR.Float64("relive_percent")
	}

	return dAtA, nil
}

var vAlIdAtOrPassiveSpellData = map[string]*config.Validator{

	"id":                        config.ParseValidator("int>0", "", false, nil, nil),
	"animation":                 config.ParseValidator("uint", "", false, nil, nil),
	"trigger_rate":              config.ParseValidator("float64>0", "", false, nil, nil),
	"trigger_type":              config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.SpellTriggerType_value), nil),
	"trigger_hit":               config.ParseValidator("uint", "", false, nil, nil),
	"target_cooldown":           config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"self_state":                config.ParseValidator("string", "", true, nil, nil),
	"target_state":              config.ParseValidator("string", "", true, nil, nil),
	"spell":                     config.ParseValidator("string", "", false, nil, nil),
	"excite_effect_type":        config.ParseValidator("uint", "", false, nil, nil),
	"rage":                      config.ParseValidator("uint", "", false, nil, nil),
	"sprite_stat":               config.ParseValidator("string", "", false, nil, nil),
	"been_hurt_effect_inc_type": config.ParseValidator("uint", "", true, nil, nil),
	"been_hurt_effect_inc":      config.ParseValidator("float64", "", true, nil, nil),
	"been_hurt_effect_dec_type": config.ParseValidator("uint", "", true, nil, nil),
	"been_hurt_effect_dec":      config.ParseValidator("float64", "", true, nil, nil),
	"relive_percent":            config.ParseValidator("float64", "", false, nil, []string{"0"}),
}

func (dAtA *PassiveSpellData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *PassiveSpellData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *PassiveSpellData) Encode() *shared_proto.PassiveSpellDataProto {
	out := &shared_proto.PassiveSpellDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Animation = config.U64ToI32(dAtA.Animation)
	out.TriggerRate = config.F64ToI32X1000(dAtA.TriggerRate)
	out.TriggerType = dAtA.TriggerType
	if dAtA.TriggerTarget != nil {
		out.TriggerTarget = dAtA.TriggerTarget.Encode()
	}
	out.TriggerHit = config.U64ToI32(dAtA.TriggerHit)
	out.TargetCooldown = int32(dAtA.TargetCooldown / time.Millisecond)
	if dAtA.SelfState != nil {
		out.SelfState = config.U64a2I32a(GetStateDataKeyArray(dAtA.SelfState))
	}
	if dAtA.TargetState != nil {
		out.TargetState = config.U64a2I32a(GetStateDataKeyArray(dAtA.TargetState))
	}
	if dAtA.Spell != nil {
		out.Spell = config.U64ToI32(dAtA.Spell.Id)
	}
	out.ExciteEffectType = config.U64ToI32(dAtA.ExciteEffectType)
	out.Rage = config.U64ToI32(dAtA.Rage)
	if dAtA.SpriteStat != nil {
		out.SpriteStat = dAtA.SpriteStat.Encode()
	}
	out.BeenHurtEffectIncType = config.U64a2I32a(dAtA.BeenHurtEffectIncType)
	out.BeenHurtEffectInc = config.F64a2I32aX1000(dAtA.BeenHurtEffectInc)
	out.BeenHurtEffectDecType = config.U64a2I32a(dAtA.BeenHurtEffectDecType)
	out.BeenHurtEffectDec = config.F64a2I32aX1000(dAtA.BeenHurtEffectDec)
	out.RelivePercent = config.F64ToI32X1000(dAtA.RelivePercent)

	return out
}

func ArrayEncodePassiveSpellData(datas []*PassiveSpellData) []*shared_proto.PassiveSpellDataProto {

	out := make([]*shared_proto.PassiveSpellDataProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *PassiveSpellData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if err := dAtA.TriggerTarget.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS0); err != nil {
		return err
	}

	uint64Keys = pArSeR.Uint64Array("self_state", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetStateData(v)
		if obj != nil {
			dAtA.SelfState = append(dAtA.SelfState, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[self_state] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("self_state"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("target_state", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetStateData(v)
		if obj != nil {
			dAtA.TargetState = append(dAtA.TargetState, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[target_state] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("target_state"), *pArSeR)
		}
	}

	dAtA.Spell = cOnFigS.GetSpellData(pArSeR.Uint64("spell"))
	if dAtA.Spell == nil && pArSeR.Uint64("spell") != 0 {
		return errors.Errorf("%s 配置的关联字段[spell] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("spell"), *pArSeR)
	}

	dAtA.SpriteStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("sprite_stat"))
	if dAtA.SpriteStat == nil && pArSeR.Uint64("sprite_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[sprite_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("sprite_stat"), *pArSeR)
	}

	return nil
}

// start with Spell ----------------------------------

func LoadSpell(gos *config.GameObjects) (map[uint64]*Spell, map[*Spell]*config.ObjectParser, error) {
	fIlEnAmE := confpath.SpellPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*Spell, len(lIsT))
	pArSeRmAp := make(map[*Spell]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrSpell) {
			continue
		}

		dAtA, err := NewSpell(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Id
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Id], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedSpell(dAtAmAp map[*Spell]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SpellPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetSpellKeyArray(datas []*Spell) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewSpell(fIlEnAmE string, pArSeR *config.ObjectParser) (*Spell, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSpell)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &Spell{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// releated field: Icon
	dAtA.Desc = pArSeR.String("desc")

	return dAtA, nil
}

var vAlIdAtOrSpell = map[string]*config.Validator{

	"id":   config.ParseValidator("int>0", "", false, nil, nil),
	"name": config.ParseValidator("string>0", "", false, nil, nil),
	"icon": config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"desc": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *Spell) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *Spell) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *Spell) Encode() *shared_proto.SpellProto {
	out := &shared_proto.SpellProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.Desc = dAtA.Desc

	return out
}

func ArrayEncodeSpell(datas []*Spell) []*shared_proto.SpellProto {

	out := make([]*shared_proto.SpellProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *Spell) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("icon") {
		dAtA.Icon = cOnFigS.GetIcon(pArSeR.String("icon"))
	} else {
		dAtA.Icon = cOnFigS.GetIcon("Icon")
	}
	if dAtA.Icon == nil {
		return errors.Errorf("%s 配置的关联字段[icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("icon"), *pArSeR)
	}

	return nil
}

// start with SpellData ----------------------------------

func LoadSpellData(gos *config.GameObjects) (map[uint64]*SpellData, map[*SpellData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.SpellDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*SpellData, len(lIsT))
	pArSeRmAp := make(map[*SpellData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrSpellData) {
			continue
		}

		dAtA, err := NewSpellData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Id
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Id], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedSpellData(dAtAmAp map[*SpellData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SpellDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetSpellDataKeyArray(datas []*SpellData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewSpellData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SpellData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSpellData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SpellData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Animation = pArSeR.Uint64("animation")
	dAtA.Cooldown, err = config.ParseDuration(pArSeR.String("cooldown"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[cooldown] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("cooldown"), dAtA)
	}

	dAtA.StrongeDuration, err = config.ParseDuration(pArSeR.String("stronge_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[stronge_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("stronge_duration"), dAtA)
	}

	dAtA.RageSpell = pArSeR.Bool("rage_spell")
	dAtA.KeepMove = pArSeR.Bool("keep_move")
	dAtA.FriendSpell = pArSeR.Bool("friend_spell")
	dAtA.SelfAsTarget = pArSeR.Bool("self_as_target")
	dAtA.TargetSubType = shared_proto.SpellTargetSubType(shared_proto.SpellTargetSubType_value[strings.ToUpper(pArSeR.String("target_sub_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("target_sub_type"), 10, 32); err == nil {
		dAtA.TargetSubType = shared_proto.SpellTargetSubType(i)
	}

	dAtA.Target, err = NewSpellTargetData(fIlEnAmE, pArSeR)
	if err != nil {
		return nil, err
	}
	dAtA.ReleaseRange = pArSeR.Uint64("release_range")
	dAtA.HurtRange = pArSeR.Uint64("hurt_range")
	dAtA.HurtCount = pArSeR.Uint64("hurt_count")
	dAtA.HurtType = 0
	if pArSeR.KeyExist("hurt_type") {
		dAtA.HurtType = pArSeR.Uint64("hurt_type")
	}

	dAtA.Coef = pArSeR.Float64("coef")
	dAtA.FlySpeed = pArSeR.Uint64("fly_speed")
	dAtA.EffectType = pArSeR.Uint64("effect_type")
	// releated field: SelfState
	dAtA.SelfStateRate = pArSeR.Float64Array("self_state_rate", "", false)
	// releated field: TargetState
	dAtA.TargetStateRate = pArSeR.Float64Array("target_state_rate", "", false)
	dAtA.SelfRage = pArSeR.Int("self_rage")
	dAtA.TargetRage = pArSeR.Int("target_rage")

	return dAtA, nil
}

var vAlIdAtOrSpellData = map[string]*config.Validator{

	"id":                config.ParseValidator("int>0", "", false, nil, nil),
	"name":              config.ParseValidator("string>0", "", false, nil, nil),
	"animation":         config.ParseValidator("uint", "", false, nil, nil),
	"cooldown":          config.ParseValidator("string", "", false, nil, nil),
	"stronge_duration":  config.ParseValidator("string", "", false, nil, nil),
	"rage_spell":        config.ParseValidator("bool", "", false, nil, nil),
	"keep_move":         config.ParseValidator("bool", "", false, nil, nil),
	"friend_spell":      config.ParseValidator("bool", "", false, nil, nil),
	"self_as_target":    config.ParseValidator("bool", "", false, nil, nil),
	"target_sub_type":   config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.SpellTargetSubType_value), nil),
	"release_range":     config.ParseValidator("uint", "", false, nil, nil),
	"hurt_range":        config.ParseValidator("uint", "", false, nil, nil),
	"hurt_count":        config.ParseValidator("uint", "", false, nil, nil),
	"hurt_type":         config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"coef":              config.ParseValidator("float64", "", false, nil, nil),
	"fly_speed":         config.ParseValidator("uint", "", false, nil, nil),
	"effect_type":       config.ParseValidator("uint", "", false, nil, nil),
	"self_state":        config.ParseValidator("string", "", true, nil, nil),
	"self_state_rate":   config.ParseValidator("float64,duplicate", "", true, nil, nil),
	"target_state":      config.ParseValidator("string", "", true, nil, nil),
	"target_state_rate": config.ParseValidator("float64,duplicate", "", true, nil, nil),
	"self_rage":         config.ParseValidator("uint", "", false, nil, nil),
	"target_rage":       config.ParseValidator("int", "", false, nil, nil),
}

func (dAtA *SpellData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SpellData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SpellData) Encode() *shared_proto.SpellDataProto {
	out := &shared_proto.SpellDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Animation = config.U64ToI32(dAtA.Animation)
	out.Cooldown = int32(dAtA.Cooldown / time.Millisecond)
	out.StrongeDuration = int32(dAtA.StrongeDuration / time.Millisecond)
	out.RageSpell = dAtA.RageSpell
	out.KeepMove = dAtA.KeepMove
	out.FriendSpell = dAtA.FriendSpell
	out.SelfAsTarget = dAtA.SelfAsTarget
	out.TargetSubType = dAtA.TargetSubType
	if dAtA.Target != nil {
		out.Target = dAtA.Target.Encode()
	}
	out.ReleaseRange = config.U64ToI32(dAtA.ReleaseRange)
	out.HurtRange = config.U64ToI32(dAtA.HurtRange)
	out.HurtCount = config.U64ToI32(dAtA.HurtCount)
	out.HurtType = config.U64ToI32(dAtA.HurtType)
	out.Coef = config.F64ToI32X1000(dAtA.Coef)
	out.FlySpeed = config.U64ToI32(dAtA.FlySpeed)
	out.EffectType = config.U64ToI32(dAtA.EffectType)
	if dAtA.SelfState != nil {
		out.SelfState = config.U64a2I32a(GetStateDataKeyArray(dAtA.SelfState))
	}
	out.SelfStateRate = config.F64a2I32aX1000(dAtA.SelfStateRate)
	if dAtA.TargetState != nil {
		out.TargetState = config.U64a2I32a(GetStateDataKeyArray(dAtA.TargetState))
	}
	out.TargetStateRate = config.F64a2I32aX1000(dAtA.TargetStateRate)
	out.SelfRage = int32(dAtA.SelfRage)
	out.TargetRage = int32(dAtA.TargetRage)

	return out
}

func ArrayEncodeSpellData(datas []*SpellData) []*shared_proto.SpellDataProto {

	out := make([]*shared_proto.SpellDataProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *SpellData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if err := dAtA.Target.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS0); err != nil {
		return err
	}

	uint64Keys = pArSeR.Uint64Array("self_state", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetStateData(v)
		if obj != nil {
			dAtA.SelfState = append(dAtA.SelfState, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[self_state] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("self_state"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("target_state", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetStateData(v)
		if obj != nil {
			dAtA.TargetState = append(dAtA.TargetState, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[target_state] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("target_state"), *pArSeR)
		}
	}

	return nil
}

// start with SpellFacadeData ----------------------------------

func LoadSpellFacadeData(gos *config.GameObjects) (map[uint64]*SpellFacadeData, map[*SpellFacadeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.SpellFacadeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*SpellFacadeData, len(lIsT))
	pArSeRmAp := make(map[*SpellFacadeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrSpellFacadeData) {
			continue
		}

		dAtA, err := NewSpellFacadeData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Id
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Id], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedSpellFacadeData(dAtAmAp map[*SpellFacadeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SpellFacadeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetSpellFacadeDataKeyArray(datas []*SpellFacadeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewSpellFacadeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SpellFacadeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSpellFacadeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SpellFacadeData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.SubDesc = ""
	if pArSeR.KeyExist("sub_desc") {
		dAtA.SubDesc = pArSeR.String("sub_desc")
	}

	dAtA.FightAmountCoef = 0
	if pArSeR.KeyExist("fight_amount_coef") {
		dAtA.FightAmountCoef = pArSeR.Uint64("fight_amount_coef")
	}

	dAtA.Group = pArSeR.Uint64("group")
	dAtA.Level = pArSeR.Uint64("level")
	// skip field: SpellType
	// releated field: Spell
	// releated field: PassiveSpell
	// releated field: BuildingEffect

	return dAtA, nil
}

var vAlIdAtOrSpellFacadeData = map[string]*config.Validator{

	"id":                config.ParseValidator("int>0", "", false, nil, nil),
	"name":              config.ParseValidator("string>0", "", false, nil, nil),
	"icon":              config.ParseValidator("string", "", false, nil, nil),
	"desc":              config.ParseValidator("string", "", false, nil, nil),
	"sub_desc":          config.ParseValidator("string", "", false, nil, []string{""}),
	"fight_amount_coef": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"group":             config.ParseValidator("int>0", "", false, nil, nil),
	"level":             config.ParseValidator("int>0", "", false, nil, nil),
	"spell":             config.ParseValidator("string", "", false, nil, nil),
	"passive_spell":     config.ParseValidator("string", "", true, nil, nil),
	"building_effect":   config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *SpellFacadeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SpellFacadeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SpellFacadeData) Encode() *shared_proto.SpellFacadeDataProto {
	out := &shared_proto.SpellFacadeDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.IconId = dAtA.Icon
	out.Desc = dAtA.Desc
	out.SubDesc = dAtA.SubDesc
	out.Group = config.U64ToI32(dAtA.Group)
	out.Level = config.U64ToI32(dAtA.Level)
	out.SpellType = config.U64ToI32(dAtA.SpellType)

	return out
}

func ArrayEncodeSpellFacadeData(datas []*SpellFacadeData) []*shared_proto.SpellFacadeDataProto {

	out := make([]*shared_proto.SpellFacadeDataProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *SpellFacadeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Spell = cOnFigS.GetSpellData(pArSeR.Uint64("spell"))
	if dAtA.Spell == nil && pArSeR.Uint64("spell") != 0 {
		return errors.Errorf("%s 配置的关联字段[spell] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("spell"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("passive_spell", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetPassiveSpellData(v)
		if obj != nil {
			dAtA.PassiveSpell = append(dAtA.PassiveSpell, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[passive_spell] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("passive_spell"), *pArSeR)
		}
	}

	dAtA.BuildingEffect = cOnFigS.GetBuildingEffectData(pArSeR.Int("building_effect"))
	if dAtA.BuildingEffect == nil && pArSeR.Int("building_effect") != 0 {
		return errors.Errorf("%s 配置的关联字段[building_effect] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("building_effect"), *pArSeR)
	}

	return nil
}

// start with SpellTargetData ----------------------------------

func NewSpellTargetData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SpellTargetData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSpellTargetData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SpellTargetData{}

	for _, v := range pArSeR.StringArray("target_race", "", false) {
		x := shared_proto.Race(shared_proto.Race_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.Race(i)
		}
		dAtA.TargetRace = append(dAtA.TargetRace, x)
	}

	dAtA.TargetEffectType = pArSeR.Uint64("target_effect_type")
	dAtA.TargetUnmovable = pArSeR.Bool("target_unmovable")
	dAtA.TargetNotAttackable = pArSeR.Bool("target_not_attackable")
	dAtA.TargetSilence = pArSeR.Bool("target_silence")
	dAtA.TargetStun = pArSeR.Bool("target_stun")

	return dAtA, nil
}

var vAlIdAtOrSpellTargetData = map[string]*config.Validator{

	"target_race":           config.ParseValidator("string", "", true, config.EnumMapKeys(shared_proto.Race_value, 0), nil),
	"target_effect_type":    config.ParseValidator("uint", "", false, nil, nil),
	"target_unmovable":      config.ParseValidator("bool", "", false, nil, nil),
	"target_not_attackable": config.ParseValidator("bool", "", false, nil, nil),
	"target_silence":        config.ParseValidator("bool", "", false, nil, nil),
	"target_stun":           config.ParseValidator("bool", "", false, nil, nil),
}

func (dAtA *SpellTargetData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SpellTargetData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SpellTargetData) Encode() *shared_proto.SpellTargetDataProto {
	out := &shared_proto.SpellTargetDataProto{}
	out.TargetRace = dAtA.TargetRace
	out.TargetEffectType = config.U64ToI32(dAtA.TargetEffectType)
	out.TargetUnmovable = dAtA.TargetUnmovable
	out.TargetNotAttackable = dAtA.TargetNotAttackable
	out.TargetSilence = dAtA.TargetSilence
	out.TargetStun = dAtA.TargetStun

	return out
}

func ArrayEncodeSpellTargetData(datas []*SpellTargetData) []*shared_proto.SpellTargetDataProto {

	out := make([]*shared_proto.SpellTargetDataProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *SpellTargetData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	return nil
}

// start with StateData ----------------------------------

func LoadStateData(gos *config.GameObjects) (map[uint64]*StateData, map[*StateData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.StateDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*StateData, len(lIsT))
	pArSeRmAp := make(map[*StateData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrStateData) {
			continue
		}

		dAtA, err := NewStateData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Id
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Id], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedStateData(dAtAmAp map[*StateData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.StateDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetStateDataKeyArray(datas []*StateData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewStateData(fIlEnAmE string, pArSeR *config.ObjectParser) (*StateData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrStateData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &StateData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Animation = pArSeR.Uint64("animation")
	dAtA.StackType = shared_proto.StateStackType(shared_proto.StateStackType_value[strings.ToUpper(pArSeR.String("stack_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("stack_type"), 10, 32); err == nil {
		dAtA.StackType = shared_proto.StateStackType(i)
	}

	dAtA.StackMaxTimes = pArSeR.Uint64("stack_max_times")
	dAtA.TickTimes = pArSeR.Uint64("tick_times")
	dAtA.TickDuration, err = config.ParseDuration(pArSeR.String("tick_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tick_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tick_duration"), dAtA)
	}

	// releated field: ChangeStat
	dAtA.IsAddStat = pArSeR.Bool("is_add_stat")
	dAtA.MoveSpeedRate = pArSeR.Float64("move_speed_rate")
	dAtA.AttackSpeedRate = pArSeR.Float64("attack_speed_rate")
	dAtA.ShieldRate = pArSeR.Float64("shield_rate")
	dAtA.ShieldEffectRate = 1
	if pArSeR.KeyExist("shield_effect_rate") {
		dAtA.ShieldEffectRate = pArSeR.Float64("shield_effect_rate")
	}

	dAtA.Unmovable = pArSeR.Bool("unmovable")
	dAtA.NotAttackable = pArSeR.Bool("not_attackable")
	dAtA.Silence = pArSeR.Bool("silence")
	dAtA.Stun = pArSeR.Bool("stun")
	dAtA.EffectType = pArSeR.Uint64("effect_type")
	dAtA.BeenHurtEffectIncType = pArSeR.Uint64Array("been_hurt_effect_inc_type", "", false)
	dAtA.BeenHurtEffectInc = pArSeR.Float64Array("been_hurt_effect_inc", "", false)
	dAtA.BeenHurtEffectDecType = pArSeR.Uint64Array("been_hurt_effect_dec_type", "", false)
	dAtA.BeenHurtEffectDec = pArSeR.Float64Array("been_hurt_effect_dec", "", false)
	dAtA.DamageCoef = pArSeR.Float64("damage_coef")
	dAtA.DamageHurtType = 0
	if pArSeR.KeyExist("damage_hurt_type") {
		dAtA.DamageHurtType = pArSeR.Uint64("damage_hurt_type")
	}

	dAtA.Rage = pArSeR.Int("rage")
	dAtA.RageRecoverRate = pArSeR.Float64("rage_recover_rate")

	return dAtA, nil
}

var vAlIdAtOrStateData = map[string]*config.Validator{

	"id":                        config.ParseValidator("int>0", "", false, nil, nil),
	"name":                      config.ParseValidator("string", "", false, nil, nil),
	"animation":                 config.ParseValidator("uint", "", false, nil, nil),
	"stack_type":                config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.StateStackType_value), nil),
	"stack_max_times":           config.ParseValidator("int>0", "", false, nil, nil),
	"tick_times":                config.ParseValidator("int>0", "", false, nil, nil),
	"tick_duration":             config.ParseValidator("string", "", false, nil, nil),
	"change_stat":               config.ParseValidator("string", "", false, nil, nil),
	"is_add_stat":               config.ParseValidator("bool", "", false, nil, nil),
	"move_speed_rate":           config.ParseValidator("float64", "", false, nil, nil),
	"attack_speed_rate":         config.ParseValidator("float64", "", false, nil, nil),
	"shield_rate":               config.ParseValidator("float64", "", false, nil, nil),
	"shield_effect_rate":        config.ParseValidator("float64", "", false, nil, []string{"1"}),
	"unmovable":                 config.ParseValidator("bool", "", false, nil, nil),
	"not_attackable":            config.ParseValidator("bool", "", false, nil, nil),
	"silence":                   config.ParseValidator("bool", "", false, nil, nil),
	"stun":                      config.ParseValidator("bool", "", false, nil, nil),
	"effect_type":               config.ParseValidator("uint", "", false, nil, nil),
	"been_hurt_effect_inc_type": config.ParseValidator("uint", "", true, nil, nil),
	"been_hurt_effect_inc":      config.ParseValidator("float64", "", true, nil, nil),
	"been_hurt_effect_dec_type": config.ParseValidator("uint", "", true, nil, nil),
	"been_hurt_effect_dec":      config.ParseValidator("float64", "", true, nil, nil),
	"damage_coef":               config.ParseValidator("float64", "", false, nil, nil),
	"damage_hurt_type":          config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"rage":                      config.ParseValidator("int", "", false, nil, nil),
	"rage_recover_rate":         config.ParseValidator("float64", "", false, nil, nil),
}

func (dAtA *StateData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *StateData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *StateData) Encode() *shared_proto.StateDataProto {
	out := &shared_proto.StateDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Animation = config.U64ToI32(dAtA.Animation)
	out.StackType = dAtA.StackType
	out.StackMaxTimes = config.U64ToI32(dAtA.StackMaxTimes)
	out.TickTimes = config.U64ToI32(dAtA.TickTimes)
	out.TickDuration = int32(dAtA.TickDuration / time.Millisecond)
	if dAtA.ChangeStat != nil {
		out.ChangeStat = dAtA.ChangeStat.Encode()
	}
	out.IsAddStat = dAtA.IsAddStat
	out.MoveSpeedRate = config.F64ToI32X1000(dAtA.MoveSpeedRate)
	out.AttackSpeedRate = config.F64ToI32X1000(dAtA.AttackSpeedRate)
	out.ShieldRate = config.F64ToI32X1000(dAtA.ShieldRate)
	out.ShieldEffectRate = config.F64ToI32X1000(dAtA.ShieldEffectRate)
	out.Unmovable = dAtA.Unmovable
	out.NotAttackable = dAtA.NotAttackable
	out.Silence = dAtA.Silence
	out.Stun = dAtA.Stun
	out.EffectType = config.U64ToI32(dAtA.EffectType)
	out.BeenHurtEffectIncType = config.U64a2I32a(dAtA.BeenHurtEffectIncType)
	out.BeenHurtEffectInc = config.F64a2I32aX1000(dAtA.BeenHurtEffectInc)
	out.BeenHurtEffectDecType = config.U64a2I32a(dAtA.BeenHurtEffectDecType)
	out.BeenHurtEffectDec = config.F64a2I32aX1000(dAtA.BeenHurtEffectDec)
	out.DamageCoef = config.F64ToI32X1000(dAtA.DamageCoef)
	out.DamageHurtType = config.U64ToI32(dAtA.DamageHurtType)
	out.Rage = int32(dAtA.Rage)
	out.RageRecoverRate = config.F64ToI32X1000(dAtA.RageRecoverRate)

	return out
}

func ArrayEncodeStateData(datas []*StateData) []*shared_proto.StateDataProto {

	out := make([]*shared_proto.StateDataProto, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			o := d.Encode()
			if o != nil {
				out = append(out, o)
			}
		}
	}

	return out
}

func (dAtA *StateData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.ChangeStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("change_stat"))
	if dAtA.ChangeStat == nil && pArSeR.Uint64("change_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[change_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("change_stat"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetBuildingEffectData(int) *sub.BuildingEffectData
	GetIcon(string) *icon.Icon
	GetPassiveSpellData(uint64) *PassiveSpellData
	GetSpellData(uint64) *SpellData
	GetSpriteStat(uint64) *data.SpriteStat
	GetStateData(uint64) *StateData
}
