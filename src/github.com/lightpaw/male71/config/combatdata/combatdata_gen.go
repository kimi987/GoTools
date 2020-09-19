// AUTO_GEN, DONT MODIFY!!!
package combatdata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/race"
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

// start with CombatConfig ----------------------------------

func LoadCombatConfig(gos *config.GameObjects) (*CombatConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.CombatConfigPath
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

	dAtA, err := NewCombatConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedCombatConfig(gos *config.GameObjects, dAtA *CombatConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CombatConfigPath
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

func NewCombatConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*CombatConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCombatConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CombatConfig{}

	// skip field: Spell
	// skip field: State
	// skip field: PassiveSpell
	// skip field: SpellIdMap
	// skip field: Captain
	// skip field: NamelessCaptain
	// skip field: Race
	dAtA.FramePerSecond = 10
	if pArSeR.KeyExist("frame_per_second") {
		dAtA.FramePerSecond = pArSeR.Uint64("frame_per_second")
	}

	dAtA.ConfigDenominator = 1000
	if pArSeR.KeyExist("config_denominator") {
		dAtA.ConfigDenominator = pArSeR.Uint64("config_denominator")
	}

	dAtA.MinAttackDuration, err = config.ParseDuration(pArSeR.String("min_attack_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[min_attack_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("min_attack_duration"), dAtA)
	}

	dAtA.MaxAttackDuration, err = config.ParseDuration(pArSeR.String("max_attack_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[max_attack_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("max_attack_duration"), dAtA)
	}

	dAtA.MinMoveSpeed = pArSeR.Uint64("min_move_speed")
	dAtA.MaxMoveSpeed = pArSeR.Uint64("max_move_speed")
	// releated field: MinStat
	dAtA.MaxDuration, err = config.ParseDuration(pArSeR.String("max_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[max_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("max_duration"), dAtA)
	}

	// skip field: ScorePercent
	if pArSeR.KeyExist("check_move_duration") {
		dAtA.CheckMoveDuration, err = config.ParseDuration(pArSeR.String("check_move_duration"))
	} else {
		dAtA.CheckMoveDuration, err = config.ParseDuration("500ms")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[check_move_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("check_move_duration"), dAtA)
	}

	dAtA.CritRate = pArSeR.Float64("crit_rate")
	dAtA.Coef = pArSeR.Float64("coef")
	dAtA.CellLen = 100
	if pArSeR.KeyExist("cell_len") {
		dAtA.CellLen = pArSeR.Uint64("cell_len")
	}

	dAtA.MaxRage = pArSeR.Uint64("max_rage")
	dAtA.AddRagePerHint = pArSeR.Uint64("add_rage_per_hint")
	dAtA.AddRageLost1Percent = pArSeR.Uint64("add_rage_lost1_percent")
	dAtA.RageRecoverSpeed = pArSeR.Uint64("rage_recover_speed")
	dAtA.CombatXLen = 1000
	if pArSeR.KeyExist("combat_xlen") {
		dAtA.CombatXLen = pArSeR.Uint64("combat_xlen")
	}

	dAtA.CombatYLen = 500
	if pArSeR.KeyExist("combat_ylen") {
		dAtA.CombatYLen = pArSeR.Uint64("combat_ylen")
	}

	dAtA.InitAttackerX = 0
	if pArSeR.KeyExist("init_attacker_x") {
		dAtA.InitAttackerX = pArSeR.Uint64("init_attacker_x")
	}

	dAtA.InitDefenserX = 1000
	if pArSeR.KeyExist("init_defenser_x") {
		dAtA.InitDefenserX = pArSeR.Uint64("init_defenser_x")
	}

	dAtA.InitWallX = 1200
	if pArSeR.KeyExist("init_wall_x") {
		dAtA.InitWallX = pArSeR.Uint64("init_wall_x")
	}

	dAtA.WallWaitDuration, err = config.ParseDuration(pArSeR.String("wall_wait_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[wall_wait_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("wall_wait_duration"), dAtA)
	}

	dAtA.WallAttackFixDamageTimes = 3
	if pArSeR.KeyExist("wall_attack_fix_damage_times") {
		dAtA.WallAttackFixDamageTimes = pArSeR.Uint64("wall_attack_fix_damage_times")
	}

	dAtA.WallBeenHurtLostMaxPercent = 0
	if pArSeR.KeyExist("wall_been_hurt_lost_max_percent") {
		dAtA.WallBeenHurtLostMaxPercent = pArSeR.Float64("wall_been_hurt_lost_max_percent")
	}

	// releated field: WallSpell
	dAtA.WallFlyMinDuration, err = config.ParseDuration(pArSeR.String("wall_fly_min_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[wall_fly_min_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("wall_fly_min_duration"), dAtA)
	}

	dAtA.ShortMoveDistance = 100
	if pArSeR.KeyExist("short_move_distance") {
		dAtA.ShortMoveDistance = pArSeR.Uint64("short_move_distance")
	}

	return dAtA, nil
}

var vAlIdAtOrCombatConfig = map[string]*config.Validator{

	"frame_per_second":                config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"config_denominator":              config.ParseValidator("int>0", "", false, nil, []string{"1000"}),
	"min_attack_duration":             config.ParseValidator("string", "", false, nil, nil),
	"max_attack_duration":             config.ParseValidator("string", "", false, nil, nil),
	"min_move_speed":                  config.ParseValidator("int>0", "", false, nil, nil),
	"max_move_speed":                  config.ParseValidator("int>0", "", false, nil, nil),
	"min_stat":                        config.ParseValidator("string", "", false, nil, nil),
	"max_duration":                    config.ParseValidator("string", "", false, nil, nil),
	"check_move_duration":             config.ParseValidator("string", "", false, nil, []string{"500ms"}),
	"crit_rate":                       config.ParseValidator("float64>0", "", false, nil, nil),
	"coef":                            config.ParseValidator("float64>0", "", false, nil, nil),
	"cell_len":                        config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"max_rage":                        config.ParseValidator("int>0", "", false, nil, nil),
	"add_rage_per_hint":               config.ParseValidator("int>0", "", false, nil, nil),
	"add_rage_lost1_percent":          config.ParseValidator("int>0", "", false, nil, nil),
	"rage_recover_speed":              config.ParseValidator("int>0", "", false, nil, nil),
	"combat_xlen":                     config.ParseValidator("int>0", "", false, nil, []string{"1000"}),
	"combat_ylen":                     config.ParseValidator("int>0", "", false, nil, []string{"500"}),
	"init_attacker_x":                 config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"init_defenser_x":                 config.ParseValidator("uint", "", false, nil, []string{"1000"}),
	"init_wall_x":                     config.ParseValidator("uint", "", false, nil, []string{"1200"}),
	"wall_wait_duration":              config.ParseValidator("string", "", false, nil, nil),
	"wall_attack_fix_damage_times":    config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"wall_been_hurt_lost_max_percent": config.ParseValidator("float64>0", "", false, nil, []string{"0"}),
	"wall_spell":                      config.ParseValidator("string", "", false, nil, nil),
	"wall_fly_min_duration":           config.ParseValidator("string", "", false, nil, nil),
	"short_move_distance":             config.ParseValidator("uint", "", false, nil, []string{"100"}),
}

func (dAtA *CombatConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CombatConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CombatConfig) Encode() *shared_proto.CombatConfigProto {
	out := &shared_proto.CombatConfigProto{}
	if dAtA.Spell != nil {
		out.Spell = spell.ArrayEncodeSpellData(dAtA.Spell)
	}
	if dAtA.State != nil {
		out.State = spell.ArrayEncodeStateData(dAtA.State)
	}
	if dAtA.PassiveSpell != nil {
		out.PassiveSpell = spell.ArrayEncodePassiveSpellData(dAtA.PassiveSpell)
	}
	if dAtA.SpellIdMap != nil {
		out.SpellIdMap = dAtA.SpellIdMap
	}
	if dAtA.Captain != nil {
		out.Captain = captain.ArrayEncodeCaptainData(dAtA.Captain)
	}
	if dAtA.NamelessCaptain != nil {
		out.NamelessCaptain = captain.ArrayEncodeNamelessCaptainData(dAtA.NamelessCaptain)
	}
	if dAtA.Race != nil {
		out.Race = race.ArrayEncodeRaceData(dAtA.Race)
	}
	out.FramePerSecond = config.U64ToI32(dAtA.FramePerSecond)
	out.ConfigDenominator = config.U64ToI32(dAtA.ConfigDenominator)
	out.MinAttackDuration = int32(dAtA.MinAttackDuration / time.Millisecond)
	out.MaxAttackDuration = int32(dAtA.MaxAttackDuration / time.Millisecond)
	out.MinMoveSpeed = config.U64ToI32(dAtA.MinMoveSpeed)
	out.MaxMoveSpeed = config.U64ToI32(dAtA.MaxMoveSpeed)
	if dAtA.MinStat != nil {
		out.MinStat = dAtA.MinStat.Encode()
	}
	out.MaxDuration = config.Duration2I32Seconds(dAtA.MaxDuration)
	out.ScorePercent = config.U64a2I32a(dAtA.ScorePercent)
	out.CheckMoveDuration = int32(dAtA.CheckMoveDuration / time.Millisecond)
	out.CritRate = config.F64ToI32X1000(dAtA.CritRate)
	out.Coef = config.F64ToI32X1000(dAtA.Coef)
	out.CellLen = config.U64ToI32(dAtA.CellLen)
	out.MaxRage = config.U64ToI32(dAtA.MaxRage)
	out.AddRagePerHint = config.U64ToI32(dAtA.AddRagePerHint)
	out.AddRageLost1Percent = config.U64ToI32(dAtA.AddRageLost1Percent)
	out.RageRecoverSpeed = config.U64ToI32(dAtA.RageRecoverSpeed)
	out.InitAttackerX = config.U64ToI32(dAtA.InitAttackerX)
	out.InitDefenserX = config.U64ToI32(dAtA.InitDefenserX)
	out.InitWallX = config.U64ToI32(dAtA.InitWallX)
	out.WallWaitDuration = int32(dAtA.WallWaitDuration / time.Millisecond)
	out.WallAttackFixDamageTimes = config.U64ToI32(dAtA.WallAttackFixDamageTimes)
	out.WallBeenHurtLostMaxPercent = config.F64ToI32X1000(dAtA.WallBeenHurtLostMaxPercent)
	if dAtA.WallSpell != nil {
		out.WallSpell = config.U64ToI32(dAtA.WallSpell.Id)
	}
	out.WallFlyMinDuration = int32(dAtA.WallFlyMinDuration / time.Millisecond)
	out.ShortMoveDistance = config.U64ToI32(dAtA.ShortMoveDistance)

	return out
}

func ArrayEncodeCombatConfig(datas []*CombatConfig) []*shared_proto.CombatConfigProto {

	out := make([]*shared_proto.CombatConfigProto, 0, len(datas))
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

func (dAtA *CombatConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.MinStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("min_stat"))
	if dAtA.MinStat == nil {
		return errors.Errorf("%s 配置的关联字段[min_stat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("min_stat"), *pArSeR)
	}

	dAtA.WallSpell = cOnFigS.GetSpellData(pArSeR.Uint64("wall_spell"))
	if dAtA.WallSpell == nil {
		return errors.Errorf("%s 配置的关联字段[wall_spell] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("wall_spell"), *pArSeR)
	}

	return nil
}

// start with CombatMiscConfig ----------------------------------

func LoadCombatMiscConfig(gos *config.GameObjects) (*CombatMiscConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.CombatMiscConfigPath
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

	dAtA, err := NewCombatMiscConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedCombatMiscConfig(gos *config.GameObjects, dAtA *CombatMiscConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CombatMiscConfigPath
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

func NewCombatMiscConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*CombatMiscConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCombatMiscConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CombatMiscConfig{}

	// skip field: Spell
	// skip field: SpellAnimation
	// skip field: PassiveSpell
	// skip field: PassiveSpellAnimation
	// skip field: State
	// skip field: StateAnimation
	dAtA.MaxRage = pArSeR.Uint64("max_rage")
	// skip field: WallAttackSpeed

	return dAtA, nil
}

var vAlIdAtOrCombatMiscConfig = map[string]*config.Validator{

	"max_rage": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *CombatMiscConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CombatMiscConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CombatMiscConfig) Encode() *shared_proto.CombatMiscConfigProto {
	out := &shared_proto.CombatMiscConfigProto{}
	out.Spell = config.U64a2I32a(dAtA.Spell)
	out.SpellAnimation = config.U64a2I32a(dAtA.SpellAnimation)
	out.PassiveSpell = config.U64a2I32a(dAtA.PassiveSpell)
	out.PassiveSpellAnimation = config.U64a2I32a(dAtA.PassiveSpellAnimation)
	out.State = config.U64a2I32a(dAtA.State)
	out.StateAnimation = config.U64a2I32a(dAtA.StateAnimation)
	out.MaxRage = config.U64ToI32(dAtA.MaxRage)
	out.WallAttackSpeed = config.F64ToI32X1000(dAtA.WallAttackSpeed)

	return out
}

func ArrayEncodeCombatMiscConfig(datas []*CombatMiscConfig) []*shared_proto.CombatMiscConfigProto {

	out := make([]*shared_proto.CombatMiscConfigProto, 0, len(datas))
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

func (dAtA *CombatMiscConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

type related_configs interface {
	GetSpellData(uint64) *spell.SpellData
	GetSpriteStat(uint64) *data.SpriteStat
}
