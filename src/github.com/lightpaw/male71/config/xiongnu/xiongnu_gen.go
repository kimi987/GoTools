// AUTO_GEN, DONT MODIFY!!!
package xiongnu

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
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

// start with ResistXiongNuData ----------------------------------

func LoadResistXiongNuData(gos *config.GameObjects) (map[uint64]*ResistXiongNuData, map[*ResistXiongNuData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ResistXiongNuDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ResistXiongNuData, len(lIsT))
	pArSeRmAp := make(map[*ResistXiongNuData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrResistXiongNuData) {
			continue
		}

		dAtA, err := NewResistXiongNuData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Level
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Level], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedResistXiongNuData(dAtAmAp map[*ResistXiongNuData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ResistXiongNuDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetResistXiongNuDataKeyArray(datas []*ResistXiongNuData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewResistXiongNuData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ResistXiongNuData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrResistXiongNuData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ResistXiongNuData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Name = pArSeR.String("name")
	// releated field: NpcBaseData
	// releated field: AssistMonsters
	// releated field: ResistWaves
	// skip field: FirstResistWave
	// skip field: ScorePrizes
	// releated field: ScorePlunderPrizes
	dAtA.ScorePrestiges = pArSeR.Uint64Array("score_prestiges", "", false)
	// releated field: ResistSucPrize
	dAtA.ResistSucPrestige = pArSeR.Uint64("resist_suc_prestige")
	// skip field: NextLevel
	// skip field: TotalMonsterCount
	// skip field: MaxFightAmount
	// releated field: ShowPrizes

	return dAtA, nil
}

var vAlIdAtOrResistXiongNuData = map[string]*config.Validator{

	"level":                config.ParseValidator("int>0", "", false, nil, nil),
	"name":                 config.ParseValidator("string", "", false, nil, nil),
	"npc_base_data":        config.ParseValidator("string", "", false, nil, nil),
	"assist_monsters":      config.ParseValidator("int,duplicate", "", true, nil, nil),
	"resist_waves":         config.ParseValidator("string", "", true, nil, nil),
	"score_plunder_prizes": config.ParseValidator("int,duplicate", "", true, nil, nil),
	"score_prestiges":      config.ParseValidator("int>0,duplicate", "", true, nil, nil),
	"resist_suc_prize":     config.ParseValidator("string", "", false, nil, nil),
	"resist_suc_prestige":  config.ParseValidator("int>0", "", false, nil, nil),
	"show_prizes":          config.ParseValidator("int,duplicate", "", true, nil, nil),
}

func (dAtA *ResistXiongNuData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ResistXiongNuData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ResistXiongNuData) Encode() *shared_proto.ResistXiongNuDataProto {
	out := &shared_proto.ResistXiongNuDataProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Name = dAtA.Name
	if dAtA.NpcBaseData != nil {
		out.Npc = config.U64ToI32(dAtA.NpcBaseData.Id)
	}
	if dAtA.AssistMonsters != nil {
		out.AssistMonsters = monsterdata.ArrayEncodeMonsterMasterData(dAtA.AssistMonsters)
	}
	if dAtA.ScorePrizes != nil {
		out.ScorePrizes = resdata.ArrayEncodePrize(dAtA.ScorePrizes)
	}
	out.ScorePrestiges = config.U64a2I32a(dAtA.ScorePrestiges)
	if dAtA.ResistSucPrize != nil {
		out.ResistSucPrize = dAtA.ResistSucPrize.Encode()
	}
	out.ResistSucPrestige = config.U64ToI32(dAtA.ResistSucPrestige)
	out.TotalMonsterCount = config.U64ToI32(dAtA.TotalMonsterCount)
	out.MaxFightAmount = config.U64ToI32(dAtA.MaxFightAmount)
	if dAtA.ShowPrizes != nil {
		out.ShowPrizes = resdata.ArrayEncodePrize(dAtA.ShowPrizes)
	}

	return out
}

func ArrayEncodeResistXiongNuData(datas []*ResistXiongNuData) []*shared_proto.ResistXiongNuDataProto {

	out := make([]*shared_proto.ResistXiongNuDataProto, 0, len(datas))
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

func (dAtA *ResistXiongNuData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.NpcBaseData = cOnFigS.GetNpcBaseData(pArSeR.Uint64("npc_base_data"))
	if dAtA.NpcBaseData == nil {
		return errors.Errorf("%s 配置的关联字段[npc_base_data] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("npc_base_data"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("assist_monsters", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetMonsterMasterData(v)
		if obj != nil {
			dAtA.AssistMonsters = append(dAtA.AssistMonsters, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[assist_monsters] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assist_monsters"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("resist_waves", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetResistXiongNuWaveData(v)
		if obj != nil {
			dAtA.ResistWaves = append(dAtA.ResistWaves, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[resist_waves] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("resist_waves"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("score_plunder_prizes", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetPlunderPrize(v)
		if obj != nil {
			dAtA.ScorePlunderPrizes = append(dAtA.ScorePlunderPrizes, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[score_plunder_prizes] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("score_plunder_prizes"), *pArSeR)
		}
	}

	dAtA.ResistSucPrize = cOnFigS.GetPrize(pArSeR.Int("resist_suc_prize"))
	if dAtA.ResistSucPrize == nil {
		return errors.Errorf("%s 配置的关联字段[resist_suc_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("resist_suc_prize"), *pArSeR)
	}

	intKeys = pArSeR.IntArray("show_prizes", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetPrize(v)
		if obj != nil {
			dAtA.ShowPrizes = append(dAtA.ShowPrizes, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[show_prizes] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prizes"), *pArSeR)
		}
	}

	return nil
}

// start with ResistXiongNuMisc ----------------------------------

func LoadResistXiongNuMisc(gos *config.GameObjects) (*ResistXiongNuMisc, *config.ObjectParser, error) {
	fIlEnAmE := confpath.ResistXiongNuMiscPath
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

	dAtA, err := NewResistXiongNuMisc(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedResistXiongNuMisc(gos *config.GameObjects, dAtA *ResistXiongNuMisc, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ResistXiongNuMiscPath
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

func NewResistXiongNuMisc(fIlEnAmE string, pArSeR *config.ObjectParser) (*ResistXiongNuMisc, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrResistXiongNuMisc)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ResistXiongNuMisc{}

	dAtA.DefenseMemberCount = pArSeR.Uint64("defense_member_count")
	dAtA.InvadeDuration, err = config.ParseDuration(pArSeR.String("invade_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[invade_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("invade_duration"), dAtA)
	}

	dAtA.ResistDuration, err = config.ParseDuration(pArSeR.String("resist_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[resist_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("resist_duration"), dAtA)
	}

	stringKeys = pArSeR.StringArray("invade_wave_duration", "", false)
	dAtA.InvadeWaveDuration = make([]time.Duration, 0, len(stringKeys))
	for _, v := range stringKeys {
		obj, err := config.ParseDuration(v)
		if err != nil {
			return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[invade_wave_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("invade_wave_duration"), dAtA)
		}
		dAtA.InvadeWaveDuration = append(dAtA.InvadeWaveDuration, obj)
	}

	dAtA.TroopMoveVelocityPerSecond = 0.25
	if pArSeR.KeyExist("troop_move_velocity_per_second") {
		dAtA.TroopMoveVelocityPerSecond = pArSeR.Float64("troop_move_velocity_per_second")
	}

	dAtA.MinMoveDuration, err = config.ParseDuration(pArSeR.String("min_move_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[min_move_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("min_move_duration"), dAtA)
	}

	dAtA.MaxMoveDuration, err = config.ParseDuration(pArSeR.String("max_move_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[max_move_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("max_move_duration"), dAtA)
	}

	dAtA.RobbingDuration, err = config.ParseDuration(pArSeR.String("robbing_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[robbing_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("robbing_duration"), dAtA)
	}

	dAtA.MaxMorale = pArSeR.Uint64("max_morale")
	dAtA.WipeOutReduceMorale = pArSeR.Uint64("wipe_out_reduce_morale")
	dAtA.OneMoraleReduceSoldierPer = pArSeR.Uint64("one_morale_reduce_soldier_per")
	if pArSeR.KeyExist("attack_duration") {
		dAtA.AttackDuration, err = config.ParseDuration(pArSeR.String("attack_duration"))
	} else {
		dAtA.AttackDuration, err = config.ParseDuration("3s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[attack_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("attack_duration"), dAtA)
	}

	dAtA.OpenNeedGuildLevel = pArSeR.Uint64("open_need_guild_level")
	dAtA.MaxCanOpenTimesPerDay = 1
	if pArSeR.KeyExist("max_can_open_times_per_day") {
		dAtA.MaxCanOpenTimesPerDay = pArSeR.Uint64("max_can_open_times_per_day")
	}

	dAtA.MinBaseLevel = 4
	if pArSeR.KeyExist("min_base_level") {
		dAtA.MinBaseLevel = pArSeR.Uint64("min_base_level")
	}

	dAtA.MaxDistance = 500
	if pArSeR.KeyExist("max_distance") {
		dAtA.MaxDistance = pArSeR.Uint64("max_distance")
	}

	if pArSeR.KeyExist("start_after_server_open") {
		dAtA.StartAfterServerOpen, err = config.ParseDuration(pArSeR.String("start_after_server_open"))
	} else {
		dAtA.StartAfterServerOpen, err = config.ParseDuration("24h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[start_after_server_open] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("start_after_server_open"), dAtA)
	}

	dAtA.BaseMinRange = 20
	if pArSeR.KeyExist("base_min_range") {
		dAtA.BaseMinRange = pArSeR.Uint64("base_min_range")
	}

	dAtA.BaseMaxRange = 100
	if pArSeR.KeyExist("base_max_range") {
		dAtA.BaseMaxRange = pArSeR.Uint64("base_max_range")
	}

	return dAtA, nil
}

var vAlIdAtOrResistXiongNuMisc = map[string]*config.Validator{

	"defense_member_count":           config.ParseValidator("int>0", "", false, nil, nil),
	"invade_duration":                config.ParseValidator("string", "", false, nil, nil),
	"resist_duration":                config.ParseValidator("string", "", false, nil, nil),
	"invade_wave_duration":           config.ParseValidator("string", "", true, nil, nil),
	"troop_move_velocity_per_second": config.ParseValidator("float64>0", "", false, nil, []string{"0.25"}),
	"min_move_duration":              config.ParseValidator("string", "", false, nil, nil),
	"max_move_duration":              config.ParseValidator("string", "", false, nil, nil),
	"robbing_duration":               config.ParseValidator("string", "", false, nil, nil),
	"max_morale":                     config.ParseValidator("int>0", "", false, nil, nil),
	"wipe_out_reduce_morale":         config.ParseValidator("int>0", "", false, nil, nil),
	"one_morale_reduce_soldier_per":  config.ParseValidator("int>0", "", false, nil, nil),
	"attack_duration":                config.ParseValidator("string", "", false, nil, []string{"3s"}),
	"open_need_guild_level":          config.ParseValidator("int>0", "", false, nil, nil),
	"max_can_open_times_per_day":     config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"min_base_level":                 config.ParseValidator("int>0", "", false, nil, []string{"4"}),
	"max_distance":                   config.ParseValidator("int>0", "", false, nil, []string{"500"}),
	"start_after_server_open":        config.ParseValidator("string", "", false, nil, []string{"24h"}),
	"base_min_range":                 config.ParseValidator("int>0", "", false, nil, []string{"20"}),
	"base_max_range":                 config.ParseValidator("int>0", "", false, nil, []string{"100"}),
}

func (dAtA *ResistXiongNuMisc) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ResistXiongNuMisc) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ResistXiongNuMisc) Encode() *shared_proto.ResistXiongNuMiscProto {
	out := &shared_proto.ResistXiongNuMiscProto{}
	out.DefenseMemberCount = config.U64ToI32(dAtA.DefenseMemberCount)
	out.InvadeDuration = config.Duration2I32Seconds(dAtA.InvadeDuration)
	out.ResistDuration = config.Duration2I32Seconds(dAtA.ResistDuration)
	out.InvadeWaveDuration = timeutil.DurationArrayToSecondArray(dAtA.InvadeWaveDuration)
	out.MaxMorale = config.U64ToI32(dAtA.MaxMorale)
	out.WipeOutReduceMorale = config.U64ToI32(dAtA.WipeOutReduceMorale)
	out.OneMoraleReduceSoldierPer = config.U64ToI32(dAtA.OneMoraleReduceSoldierPer)
	out.OpenNeedGuildLevel = config.U64ToI32(dAtA.OpenNeedGuildLevel)
	out.MaxCanOpenTimesPerDay = config.U64ToI32(dAtA.MaxCanOpenTimesPerDay)
	out.MinBaseLevel = config.U64ToI32(dAtA.MinBaseLevel)
	out.MaxDistance = config.U64ToI32(dAtA.MaxDistance)
	out.StartAfterServerOpen = config.Duration2I32Seconds(dAtA.StartAfterServerOpen)

	return out
}

func ArrayEncodeResistXiongNuMisc(datas []*ResistXiongNuMisc) []*shared_proto.ResistXiongNuMiscProto {

	out := make([]*shared_proto.ResistXiongNuMiscProto, 0, len(datas))
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

func (dAtA *ResistXiongNuMisc) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with ResistXiongNuScoreData ----------------------------------

func LoadResistXiongNuScoreData(gos *config.GameObjects) (map[uint64]*ResistXiongNuScoreData, map[*ResistXiongNuScoreData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ResistXiongNuScoreDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ResistXiongNuScoreData, len(lIsT))
	pArSeRmAp := make(map[*ResistXiongNuScoreData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrResistXiongNuScoreData) {
			continue
		}

		dAtA, err := NewResistXiongNuScoreData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Level
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Level], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedResistXiongNuScoreData(dAtAmAp map[*ResistXiongNuScoreData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ResistXiongNuScoreDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetResistXiongNuScoreDataKeyArray(datas []*ResistXiongNuScoreData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewResistXiongNuScoreData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ResistXiongNuScoreData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrResistXiongNuScoreData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ResistXiongNuScoreData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Name = pArSeR.String("name")
	dAtA.WipeOutInvadeMonsterCount = pArSeR.Uint64("wipe_out_invade_monster_count")
	dAtA.UnlockNextLevel = pArSeR.Bool("unlock_next_level")

	return dAtA, nil
}

var vAlIdAtOrResistXiongNuScoreData = map[string]*config.Validator{

	"level": config.ParseValidator("int>0", "", false, nil, nil),
	"name":  config.ParseValidator("string", "", false, nil, nil),
	"wipe_out_invade_monster_count": config.ParseValidator("uint", "", false, nil, nil),
	"unlock_next_level":             config.ParseValidator("bool", "", false, nil, nil),
}

func (dAtA *ResistXiongNuScoreData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ResistXiongNuScoreData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ResistXiongNuScoreData) Encode() *shared_proto.ResistXiongNuScoreProto {
	out := &shared_proto.ResistXiongNuScoreProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Name = dAtA.Name
	out.WipeOutInvadeMonsterCount = config.U64ToI32(dAtA.WipeOutInvadeMonsterCount)
	out.UnlockNextLevel = dAtA.UnlockNextLevel

	return out
}

func ArrayEncodeResistXiongNuScoreData(datas []*ResistXiongNuScoreData) []*shared_proto.ResistXiongNuScoreProto {

	out := make([]*shared_proto.ResistXiongNuScoreProto, 0, len(datas))
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

func (dAtA *ResistXiongNuScoreData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with ResistXiongNuWaveData ----------------------------------

func LoadResistXiongNuWaveData(gos *config.GameObjects) (map[uint64]*ResistXiongNuWaveData, map[*ResistXiongNuWaveData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ResistXiongNuWaveDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ResistXiongNuWaveData, len(lIsT))
	pArSeRmAp := make(map[*ResistXiongNuWaveData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrResistXiongNuWaveData) {
			continue
		}

		dAtA, err := NewResistXiongNuWaveData(fIlEnAmE, pArSeR)
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

func SetRelatedResistXiongNuWaveData(dAtAmAp map[*ResistXiongNuWaveData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ResistXiongNuWaveDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetResistXiongNuWaveDataKeyArray(datas []*ResistXiongNuWaveData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewResistXiongNuWaveData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ResistXiongNuWaveData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrResistXiongNuWaveData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ResistXiongNuWaveData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// releated field: Monsters
	// skip field: Next

	return dAtA, nil
}

var vAlIdAtOrResistXiongNuWaveData = map[string]*config.Validator{

	"id":       config.ParseValidator("int>0", "", false, nil, nil),
	"name":     config.ParseValidator("string", "", false, nil, nil),
	"monsters": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *ResistXiongNuWaveData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("monsters", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetMonsterMasterData(v)
		if obj != nil {
			dAtA.Monsters = append(dAtA.Monsters, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[monsters] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("monsters"), *pArSeR)
		}
	}

	return nil
}

type related_configs interface {
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
	GetNpcBaseData(uint64) *basedata.NpcBaseData
	GetPlunderPrize(uint64) *resdata.PlunderPrize
	GetPrize(int) *resdata.Prize
	GetResistXiongNuWaveData(uint64) *ResistXiongNuWaveData
}
