// AUTO_GEN, DONT MODIFY!!!
package race

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/spell"
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

// start with RaceConfig ----------------------------------

func LoadRaceConfig(gos *config.GameObjects) (*RaceConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.RaceConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	dAtA, err := NewRaceConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedRaceConfig(gos *config.GameObjects, dAtA *RaceConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RaceConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return err
	}

	var pArSeR *config.ObjectParser
	if len(lIsT) <= 0 {
		pArSeR = config.NewObjectParser(nil, nil, 0)
	} else {
		pArSeR = lIsT[0]
	}

	return dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS)
}

func NewRaceConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*RaceConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRaceConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RaceConfig{}

	return dAtA, nil
}

var vAlIdAtOrRaceConfig = map[string]*config.Validator{}

func (dAtA *RaceConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with RaceData ----------------------------------

func LoadRaceData(gos *config.GameObjects) (map[int]*RaceData, map[*RaceData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.RaceDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[int]*RaceData, len(lIsT))
	pArSeRmAp := make(map[*RaceData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrRaceData) {
			continue
		}

		dAtA, err := NewRaceData(fIlEnAmE, pArSeR)
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

func SetRelatedRaceData(dAtAmAp map[*RaceData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RaceDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetRaceDataKeyArray(datas []*RaceData) []int {

	out := make([]int, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewRaceData(fIlEnAmE string, pArSeR *config.ObjectParser) (*RaceData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRaceData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RaceData{}

	dAtA.Id = pArSeR.Int("id")
	dAtA.Race = shared_proto.Race(shared_proto.Race_value[strings.ToUpper(pArSeR.String("race"))])
	if i, err := strconv.ParseInt(pArSeR.String("race"), 10, 32); err == nil {
		dAtA.Race = shared_proto.Race(i)
	}

	dAtA.Name = "name"
	if pArSeR.KeyExist("name") {
		dAtA.Name = pArSeR.String("name")
	}

	dAtA.IsFar = false
	if pArSeR.KeyExist("is_far") {
		dAtA.IsFar = pArSeR.Bool("is_far")
	}

	dAtA.AttackRange = pArSeR.Uint64("attack_range")
	dAtA.MoveTimesPerRound = pArSeR.Uint64("move_times_per_round")
	dAtA.MoveSpeed = pArSeR.Uint64("move_speed")
	dAtA.ViewRange = pArSeR.Uint64("view_range")
	for _, v := range pArSeR.StringArray("priority", "", false) {
		x := shared_proto.Race(shared_proto.Race_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.Race(i)
		}
		dAtA.Priority = append(dAtA.Priority, x)
	}

	dAtA.RaceCoef = pArSeR.Uint64Array("race_coef", "", false)
	dAtA.WallCoef = pArSeR.Uint64("wall_coef")
	dAtA.AbilityRate = pArSeR.Float64Array("ability_rate", "", false)
	for _, v := range pArSeR.StringArray("restraint_race", "", false) {
		x := shared_proto.Race(shared_proto.Race_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.Race(i)
		}
		dAtA.RestraintRace = append(dAtA.RestraintRace, x)
	}

	dAtA.RestraintRoundType = shared_proto.RestraintRoundType(shared_proto.RestraintRoundType_value[strings.ToUpper(pArSeR.String("restraint_round_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("restraint_round_type"), 10, 32); err == nil {
		dAtA.RestraintRoundType = shared_proto.RestraintRoundType(i)
	}

	// releated field: RestraintSpell
	dAtA.UnlockRestraintSpellNeedAbility = pArSeR.Uint64("unlock_restraint_spell_need_ability")
	// releated field: NormalSpell
	// releated field: SoldierSpell
	dAtA.SoldierAttackSpeed = 5
	if pArSeR.KeyExist("soldier_attack_speed") {
		dAtA.SoldierAttackSpeed = pArSeR.Float64("soldier_attack_speed")
	}

	dAtA.SoldierAttackRange = 30
	if pArSeR.KeyExist("soldier_attack_range") {
		dAtA.SoldierAttackRange = pArSeR.Uint64("soldier_attack_range")
	}

	dAtA.GemTypes = pArSeR.Uint64Array("gem_types", "", false)

	return dAtA, nil
}

var vAlIdAtOrRaceData = map[string]*config.Validator{

	"id":                                  config.ParseValidator("int>0", "", false, nil, nil),
	"race":                                config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Race_value, 0), nil),
	"name":                                config.ParseValidator("string", "", false, nil, []string{"name"}),
	"is_far":                              config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"attack_range":                        config.ParseValidator("int>0", "", false, nil, nil),
	"move_times_per_round":                config.ParseValidator("int>0", "", false, nil, nil),
	"move_speed":                          config.ParseValidator("int>0", "", false, nil, nil),
	"view_range":                          config.ParseValidator("int>0", "", false, nil, nil),
	"priority":                            config.ParseValidator("string,count=5,notNil", "", true, config.EnumMapKeys(shared_proto.Race_value, 0), nil),
	"race_coef":                           config.ParseValidator("int>0,count=5,notNil,duplicate", "", true, nil, nil),
	"wall_coef":                           config.ParseValidator("int>0", "", false, nil, nil),
	"ability_rate":                        config.ParseValidator("float64>0,count=4,notNil,duplicate,sum=1", "", true, nil, nil),
	"restraint_race":                      config.ParseValidator("string,notAllNil", "", true, config.EnumMapKeys(shared_proto.Race_value, 0), nil),
	"restraint_round_type":                config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.RestraintRoundType_value, 0), nil),
	"restraint_spell_id":                  config.ParseValidator("string", "", false, nil, nil),
	"unlock_restraint_spell_need_ability": config.ParseValidator("uint", "", false, nil, nil),
	"normal_spell_id":                     config.ParseValidator("string", "", false, nil, nil),
	"soldier_spell":                       config.ParseValidator("string", "", true, nil, nil),
	"soldier_attack_speed":                config.ParseValidator("float64>0", "", false, nil, []string{"5"}),
	"soldier_attack_range":                config.ParseValidator("int>0", "", false, nil, []string{"30"}),
	"gem_types":                           config.ParseValidator("int>0,count=9,notNil,duplicate", "", true, nil, nil),
}

func (dAtA *RaceData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *RaceData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *RaceData) Encode() *shared_proto.RaceDataProto {
	out := &shared_proto.RaceDataProto{}
	out.Id = int32(dAtA.Id)
	out.Race = dAtA.Race
	out.IsFar = dAtA.IsFar
	out.AttackRange = config.U64ToI32(dAtA.AttackRange)
	out.MoveTimesPerRound = config.U64ToI32(dAtA.MoveTimesPerRound)
	out.MoveSpeed = config.U64ToI32(dAtA.MoveSpeed)
	out.ViewRange = config.U64ToI32(dAtA.ViewRange)
	out.Priority = dAtA.Priority
	out.RaceCoef = config.U64a2I32a(dAtA.RaceCoef)
	out.WallCoef = config.U64ToI32(dAtA.WallCoef)
	out.AbilityRate = config.F64a2I32aX1000(dAtA.AbilityRate)
	out.RestraintRace = dAtA.RestraintRace
	out.RestraintRoundType = dAtA.RestraintRoundType
	if dAtA.RestraintSpell != nil {
		out.RestraintSpellId = config.U64ToI32(dAtA.RestraintSpell.Id)
	}
	out.UnlockRestraintSpellNeedAbility = config.U64ToI32(dAtA.UnlockRestraintSpellNeedAbility)
	if dAtA.NormalSpell != nil {
		out.NormalSpellId = config.U64ToI32(dAtA.NormalSpell.Id)
	}
	if dAtA.SoldierSpell != nil {
		out.SoldierSpell = config.U64a2I32a(spell.GetSpellFacadeDataKeyArray(dAtA.SoldierSpell))
	}
	out.SoldierAttackSpeed = config.F64ToI32X1000(dAtA.SoldierAttackSpeed)
	out.SoldierAttackRange = config.U64ToI32(dAtA.SoldierAttackRange)
	out.GemTypes = config.U64a2I32a(dAtA.GemTypes)

	return out
}

func ArrayEncodeRaceData(datas []*RaceData) []*shared_proto.RaceDataProto {

	out := make([]*shared_proto.RaceDataProto, 0, len(datas))
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

func (dAtA *RaceData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.RestraintSpell = cOnFigS.GetSpell(pArSeR.Uint64("restraint_spell_id"))
	if dAtA.RestraintSpell == nil {
		return errors.Errorf("%s 配置的关联字段[restraint_spell_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("restraint_spell_id"), *pArSeR)
	}

	dAtA.NormalSpell = cOnFigS.GetSpell(pArSeR.Uint64("normal_spell_id"))
	if dAtA.NormalSpell == nil {
		return errors.Errorf("%s 配置的关联字段[normal_spell_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("normal_spell_id"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("soldier_spell", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetSpellFacadeData(v)
		if obj != nil {
			dAtA.SoldierSpell = append(dAtA.SoldierSpell, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[soldier_spell] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("soldier_spell"), *pArSeR)
		}
	}

	return nil
}

type related_configs interface {
	GetSpell(uint64) *spell.Spell
	GetSpellFacadeData(uint64) *spell.SpellFacadeData
}
