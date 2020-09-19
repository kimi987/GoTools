// AUTO_GEN, DONT MODIFY!!!
package sub

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
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

// start with BuildingEffectData ----------------------------------

func LoadBuildingEffectData(gos *config.GameObjects) (map[int]*BuildingEffectData, map[*BuildingEffectData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BuildingEffectDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[int]*BuildingEffectData, len(lIsT))
	pArSeRmAp := make(map[*BuildingEffectData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBuildingEffectData) {
			continue
		}

		dAtA, err := NewBuildingEffectData(fIlEnAmE, pArSeR)
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

func SetRelatedBuildingEffectData(dAtAmAp map[*BuildingEffectData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BuildingEffectDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBuildingEffectDataKeyArray(datas []*BuildingEffectData) []int {

	out := make([]int, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBuildingEffectData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BuildingEffectData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBuildingEffectData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BuildingEffectData{}

	dAtA.Id = pArSeR.Int("id")
	stringKeys = pArSeR.StringArray("capcity", "", false)
	dAtA.Capcity = make([]*data.Amount, 0, len(stringKeys))
	for _, v := range stringKeys {
		obj, err := data.ParseAmount(v)
		if err != nil {
			return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[capcity] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("capcity"), dAtA)
		}
		dAtA.Capcity = append(dAtA.Capcity, obj)
	}

	dAtA.ProtectedCapcity, err = data.ParseAmount(pArSeR.String("protected_capcity"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[protected_capcity] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("protected_capcity"), dAtA)
	}

	dAtA.OutputType = shared_proto.ResType(shared_proto.ResType_value[strings.ToUpper(pArSeR.String("output_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("output_type"), 10, 32); err == nil {
		dAtA.OutputType = shared_proto.ResType(i)
	}

	dAtA.Output, err = data.ParseAmount(pArSeR.String("output"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[output] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("output"), dAtA)
	}

	dAtA.OutputCapcity, err = data.ParseAmount(pArSeR.String("output_capcity"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[output_capcity] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("output_capcity"), dAtA)
	}

	dAtA.SoldierCapcity = pArSeR.Uint64("soldier_capcity")
	dAtA.SoldierOutput = pArSeR.Uint64("soldier_output")
	dAtA.ForceSoldier = pArSeR.Uint64("force_soldier")
	dAtA.NewSoldierOutput = pArSeR.Uint64("new_soldier_output")
	dAtA.NewSoldierCapcity = pArSeR.Uint64("new_soldier_capcity")
	dAtA.WoundedSoldierCapcity = pArSeR.Uint64("wounded_soldier_capcity")
	dAtA.RecruitSoldierCount = pArSeR.Uint64("recruit_soldier_count")
	// releated field: FarStat
	// releated field: CloseStat
	for _, v := range pArSeR.StringArray("soldier_race", "", false) {
		x := shared_proto.Race(shared_proto.Race_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.Race(i)
		}
		dAtA.SoldierRace = append(dAtA.SoldierRace, x)
	}

	// releated field: SoldierStat
	// releated field: AllSoldierStat
	dAtA.SoldierLoad = pArSeR.Uint64("soldier_load")
	dAtA.TrainOutput = pArSeR.Float64("train_output")
	dAtA.TrainCapcity = pArSeR.Float64("train_capcity")
	dAtA.TrainCoef = pArSeR.Float64("train_coef")
	dAtA.TrainExpPerHour = pArSeR.Uint64("train_exp_per_hour")
	dAtA.BuildingWorkerCdr = pArSeR.Float64("building_worker_cdr")
	if pArSeR.KeyExist("building_worker_fatigue_duration") {
		dAtA.BuildingWorkerFatigueDuration, err = config.ParseDuration(pArSeR.String("building_worker_fatigue_duration"))
	} else {
		dAtA.BuildingWorkerFatigueDuration, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[building_worker_fatigue_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("building_worker_fatigue_duration"), dAtA)
	}

	dAtA.TechWorkerCdr = pArSeR.Float64("tech_worker_cdr")
	if pArSeR.KeyExist("tech_worker_fatigue_duration") {
		dAtA.TechWorkerFatigueDuration, err = config.ParseDuration(pArSeR.String("tech_worker_fatigue_duration"))
	} else {
		dAtA.TechWorkerFatigueDuration, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tech_worker_fatigue_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tech_worker_fatigue_duration"), dAtA)
	}

	if pArSeR.KeyExist("seek_help_cdr") {
		dAtA.SeekHelpCdr, err = config.ParseDuration(pArSeR.String("seek_help_cdr"))
	} else {
		dAtA.SeekHelpCdr, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[seek_help_cdr] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("seek_help_cdr"), dAtA)
	}

	dAtA.SeekHelpMaxTimes = 0
	if pArSeR.KeyExist("seek_help_max_times") {
		dAtA.SeekHelpMaxTimes = pArSeR.Uint64("seek_help_max_times")
	}

	dAtA.GuildDonateTimes = 0
	if pArSeR.KeyExist("guild_donate_times") {
		dAtA.GuildDonateTimes = pArSeR.Uint64("guild_donate_times")
	}

	// releated field: HomeWallStat
	dAtA.HomeWallFixDamage = pArSeR.Uint64("home_wall_fix_damage")
	dAtA.FarmOutputType = shared_proto.ResType(shared_proto.ResType_value[strings.ToUpper(pArSeR.String("farm_output_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("farm_output_type"), 10, 32); err == nil {
		dAtA.FarmOutputType = shared_proto.ResType(i)
	}

	dAtA.FarmOutput, err = data.ParseAmount(pArSeR.String("farm_output"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[farm_output] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("farm_output"), dAtA)
	}

	stringKeys = pArSeR.StringArray("tax", "", false)
	dAtA.Tax = make([]*data.Amount, 0, len(stringKeys))
	for _, v := range stringKeys {
		obj, err := data.ParseAmount(v)
		if err != nil {
			return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tax] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tax"), dAtA)
		}
		dAtA.Tax = append(dAtA.Tax, obj)
	}

	// releated field: AddedDefenseStat
	// releated field: AddedAssistStat
	// releated field: AddedCopyDefenseStat
	dAtA.BuildingCostReduceCoef = 0
	if pArSeR.KeyExist("building_cost_reduce_coef") {
		dAtA.BuildingCostReduceCoef = pArSeR.Float64("building_cost_reduce_coef")
	}

	dAtA.TechCostReduceCoef = 0
	if pArSeR.KeyExist("tech_cost_reduce_coef") {
		dAtA.TechCostReduceCoef = pArSeR.Float64("tech_cost_reduce_coef")
	}

	return dAtA, nil
}

var vAlIdAtOrBuildingEffectData = map[string]*config.Validator{

	"id":                               config.ParseValidator("int>0", "", false, nil, nil),
	"capcity":                          config.ParseValidator("string,count=4,allNilOrNot,duplicate", "", true, nil, nil),
	"protected_capcity":                config.ParseValidator("string", "", false, nil, nil),
	"output_type":                      config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.ResType_value, 0), nil),
	"output":                           config.ParseValidator("string", "", false, nil, nil),
	"output_capcity":                   config.ParseValidator("string", "", false, nil, nil),
	"soldier_capcity":                  config.ParseValidator("uint", "", false, nil, nil),
	"soldier_output":                   config.ParseValidator("uint", "", false, nil, nil),
	"force_soldier":                    config.ParseValidator("uint", "", false, nil, nil),
	"new_soldier_output":               config.ParseValidator("uint", "", false, nil, nil),
	"new_soldier_capcity":              config.ParseValidator("uint", "", false, nil, nil),
	"wounded_soldier_capcity":          config.ParseValidator("uint", "", false, nil, nil),
	"recruit_soldier_count":            config.ParseValidator("uint", "", false, nil, nil),
	"far_stat":                         config.ParseValidator("string", "", false, nil, nil),
	"close_stat":                       config.ParseValidator("string", "", false, nil, nil),
	"soldier_race":                     config.ParseValidator("string", "", true, config.EnumMapKeys(shared_proto.Race_value, 0), nil),
	"soldier_stat":                     config.ParseValidator("string", "", false, nil, nil),
	"all_soldier_stat":                 config.ParseValidator("string", "", false, nil, nil),
	"soldier_load":                     config.ParseValidator("uint", "", false, nil, nil),
	"train_output":                     config.ParseValidator("float64", "", false, nil, nil),
	"train_capcity":                    config.ParseValidator("float64", "", false, nil, nil),
	"train_coef":                       config.ParseValidator("float64", "", false, nil, nil),
	"train_exp_per_hour":               config.ParseValidator("uint", "", false, nil, nil),
	"building_worker_cdr":              config.ParseValidator("float64", "", false, nil, nil),
	"building_worker_fatigue_duration": config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"tech_worker_cdr":                  config.ParseValidator("float64", "", false, nil, nil),
	"tech_worker_fatigue_duration":     config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"seek_help_cdr":                    config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"seek_help_max_times":              config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"guild_donate_times":               config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"home_wall_stat":                   config.ParseValidator("string", "", false, nil, nil),
	"home_wall_fix_damage":             config.ParseValidator("uint", "", false, nil, nil),
	"farm_output_type":                 config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.ResType_value, 0), nil),
	"farm_output":                      config.ParseValidator("string", "", false, nil, nil),
	"tax":                              config.ParseValidator("string,count=4,allNilOrNot,duplicate", "", true, nil, nil),
	"added_defense_stat":               config.ParseValidator("string", "", false, nil, nil),
	"added_assist_stat":                config.ParseValidator("string", "", false, nil, nil),
	"added_copy_defense_stat":          config.ParseValidator("string", "", false, nil, nil),
	"building_cost_reduce_coef":        config.ParseValidator("float64", "", false, nil, []string{"0"}),
	"tech_cost_reduce_coef":            config.ParseValidator("float64", "", false, nil, []string{"0"}),
}

func (dAtA *BuildingEffectData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BuildingEffectData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BuildingEffectData) Encode() *shared_proto.DomesticEffectProto {
	out := &shared_proto.DomesticEffectProto{}
	if dAtA.Capcity != nil {
		out.Capcity = data.ArrayEncodeAmount(dAtA.Capcity)
	}
	if dAtA.ProtectedCapcity != nil {
		out.ProtectedCapcity = dAtA.ProtectedCapcity.Encode()
	}
	out.OutputType = dAtA.OutputType
	if dAtA.Output != nil {
		out.Output = dAtA.Output.Encode()
	}
	if dAtA.OutputCapcity != nil {
		out.OutputCapcity = dAtA.OutputCapcity.Encode()
	}
	out.SoldierCapcity = config.U64ToI32(dAtA.SoldierCapcity)
	out.SoldierOutput = config.U64ToI32(dAtA.SoldierOutput)
	out.ForceSoldier = config.U64ToI32(dAtA.ForceSoldier)
	out.NewSoldierOutput = config.U64ToI32(dAtA.NewSoldierOutput)
	out.NewSoldierCapcity = config.U64ToI32(dAtA.NewSoldierCapcity)
	out.WoundedSoldierCapcity = config.U64ToI32(dAtA.WoundedSoldierCapcity)
	out.RecruitSoldierCount = config.U64ToI32(dAtA.RecruitSoldierCount)
	if dAtA.FarStat != nil {
		out.FarStat = dAtA.FarStat.Encode()
	}
	if dAtA.CloseStat != nil {
		out.CloseStat = dAtA.CloseStat.Encode()
	}
	out.SoldierRace = dAtA.SoldierRace
	if dAtA.SoldierStat != nil {
		out.SoldierStat = dAtA.SoldierStat.Encode()
	}
	if dAtA.AllSoldierStat != nil {
		out.AllSoldierStat = dAtA.AllSoldierStat.Encode()
	}
	out.SoldierLoad = config.U64ToI32(dAtA.SoldierLoad)
	out.TrainOutput = config.F64ToI32X1000(dAtA.TrainOutput)
	out.TrainCapcity = config.F64ToI32X1000(dAtA.TrainCapcity)
	out.TrainCoef = config.F64ToI32X1000(dAtA.TrainCoef)
	out.TrainExpPerHour = config.U64ToI32(dAtA.TrainExpPerHour)
	out.BuildingWorkerCoef = config.F64ToI32X1000(dAtA.BuildingWorkerCdr)
	out.BuildingWorkerFatigueDuration = config.Duration2I32Seconds(dAtA.BuildingWorkerFatigueDuration)
	out.TechWorkerCoef = config.F64ToI32X1000(dAtA.TechWorkerCdr)
	out.TechWorkerFatigueDuration = config.Duration2I32Seconds(dAtA.TechWorkerFatigueDuration)
	out.SeekHelpCdr = config.Duration2I32Seconds(dAtA.SeekHelpCdr)
	out.SeekHelpMaxTimes = config.U64ToI32(dAtA.SeekHelpMaxTimes)
	out.GuildDonateTimes = config.U64ToI32(dAtA.GuildDonateTimes)
	if dAtA.HomeWallStat != nil {
		out.HomeWallStat = dAtA.HomeWallStat.Encode()
	}
	out.HomeWallFixDamage = config.U64ToI32(dAtA.HomeWallFixDamage)
	out.FarmOutputType = dAtA.FarmOutputType
	if dAtA.FarmOutput != nil {
		out.FarmOutput = dAtA.FarmOutput.Encode()
	}
	if dAtA.Tax != nil {
		out.Tax = data.ArrayEncodeAmount(dAtA.Tax)
	}
	if dAtA.AddedDefenseStat != nil {
		out.AddedDefenseStat = dAtA.AddedDefenseStat.Encode()
	}
	if dAtA.AddedAssistStat != nil {
		out.AddedAssistStat = dAtA.AddedAssistStat.Encode()
	}
	if dAtA.AddedCopyDefenseStat != nil {
		out.AddedCopyDefenseStat = dAtA.AddedCopyDefenseStat.Encode()
	}
	out.BuildingCostReduceCoef = config.F64ToI32X1000(dAtA.BuildingCostReduceCoef)
	out.TechCostReduceCoef = config.F64ToI32X1000(dAtA.TechCostReduceCoef)

	return out
}

func ArrayEncodeBuildingEffectData(datas []*BuildingEffectData) []*shared_proto.DomesticEffectProto {

	out := make([]*shared_proto.DomesticEffectProto, 0, len(datas))
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

func (dAtA *BuildingEffectData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.FarStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("far_stat"))
	if dAtA.FarStat == nil && pArSeR.Uint64("far_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[far_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("far_stat"), *pArSeR)
	}

	dAtA.CloseStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("close_stat"))
	if dAtA.CloseStat == nil && pArSeR.Uint64("close_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[close_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("close_stat"), *pArSeR)
	}

	dAtA.SoldierStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("soldier_stat"))
	if dAtA.SoldierStat == nil && pArSeR.Uint64("soldier_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[soldier_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("soldier_stat"), *pArSeR)
	}

	dAtA.AllSoldierStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("all_soldier_stat"))
	if dAtA.AllSoldierStat == nil && pArSeR.Uint64("all_soldier_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[all_soldier_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("all_soldier_stat"), *pArSeR)
	}

	dAtA.HomeWallStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("home_wall_stat"))
	if dAtA.HomeWallStat == nil && pArSeR.Uint64("home_wall_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[home_wall_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("home_wall_stat"), *pArSeR)
	}

	dAtA.AddedDefenseStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("added_defense_stat"))
	if dAtA.AddedDefenseStat == nil && pArSeR.Uint64("added_defense_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[added_defense_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("added_defense_stat"), *pArSeR)
	}

	dAtA.AddedAssistStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("added_assist_stat"))
	if dAtA.AddedAssistStat == nil && pArSeR.Uint64("added_assist_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[added_assist_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("added_assist_stat"), *pArSeR)
	}

	dAtA.AddedCopyDefenseStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("added_copy_defense_stat"))
	if dAtA.AddedCopyDefenseStat == nil && pArSeR.Uint64("added_copy_defense_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[added_copy_defense_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("added_copy_defense_stat"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetSpriteStat(uint64) *data.SpriteStat
}
