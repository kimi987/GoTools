// AUTO_GEN, DONT MODIFY!!!
package basedata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
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

// start with HomeNpcBaseData ----------------------------------

func LoadHomeNpcBaseData(gos *config.GameObjects) (map[uint64]*HomeNpcBaseData, map[*HomeNpcBaseData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.HomeNpcBaseDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*HomeNpcBaseData, len(lIsT))
	pArSeRmAp := make(map[*HomeNpcBaseData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrHomeNpcBaseData) {
			continue
		}

		dAtA, err := NewHomeNpcBaseData(fIlEnAmE, pArSeR)
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

func SetRelatedHomeNpcBaseData(dAtAmAp map[*HomeNpcBaseData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.HomeNpcBaseDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetHomeNpcBaseDataKeyArray(datas []*HomeNpcBaseData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewHomeNpcBaseData(fIlEnAmE string, pArSeR *config.ObjectParser) (*HomeNpcBaseData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrHomeNpcBaseData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &HomeNpcBaseData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Data
	dAtA.EvenOffsetX = pArSeR.Int("even_offset_x")
	dAtA.EvenOffsetY = pArSeR.Int("even_offset_y")
	dAtA.HomeBaseLevel = pArSeR.Uint64("home_base_level")
	dAtA.BaYeStage = pArSeR.Uint64("ba_ye_stage")

	return dAtA, nil
}

var vAlIdAtOrHomeNpcBaseData = map[string]*config.Validator{

	"id":              config.ParseValidator("int>0", "", false, nil, nil),
	"data":            config.ParseValidator("string", "", false, nil, nil),
	"even_offset_x":   config.ParseValidator("int", "", false, nil, nil),
	"even_offset_y":   config.ParseValidator("int", "", false, nil, nil),
	"home_base_level": config.ParseValidator("int", "", false, nil, nil),
	"ba_ye_stage":     config.ParseValidator("int", "", false, nil, nil),
}

func (dAtA *HomeNpcBaseData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Data = cOnFigS.GetNpcBaseData(pArSeR.Uint64("data"))
	if dAtA.Data == nil {
		return errors.Errorf("%s 配置的关联字段[data] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("data"), *pArSeR)
	}

	return nil
}

// start with NpcBaseData ----------------------------------

func LoadNpcBaseData(gos *config.GameObjects) (map[uint64]*NpcBaseData, map[*NpcBaseData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.NpcBaseDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*NpcBaseData, len(lIsT))
	pArSeRmAp := make(map[*NpcBaseData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrNpcBaseData) {
			continue
		}

		dAtA, err := NewNpcBaseData(fIlEnAmE, pArSeR)
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

func SetRelatedNpcBaseData(dAtAmAp map[*NpcBaseData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.NpcBaseDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetNpcBaseDataKeyArray(datas []*NpcBaseData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewNpcBaseData(fIlEnAmE string, pArSeR *config.ObjectParser) (*NpcBaseData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrNpcBaseData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &NpcBaseData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = ""
	if pArSeR.KeyExist("name") {
		dAtA.Name = pArSeR.String("name")
	}

	// releated field: Npc
	dAtA.BaseLevel = pArSeR.Uint64("base_level")
	dAtA.Model = ""
	if pArSeR.KeyExist("model") {
		dAtA.Model = pArSeR.String("model")
	}

	dAtA.DefModel = ""
	if pArSeR.KeyExist("def_model") {
		dAtA.DefModel = pArSeR.String("def_model")
	}

	dAtA.ProsperityCapcity = pArSeR.Uint64("prosperity_capcity")
	if pArSeR.KeyExist("rob_max_duration") {
		dAtA.RobMaxDuration, err = config.ParseDuration(pArSeR.String("rob_max_duration"))
	} else {
		dAtA.RobMaxDuration, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[rob_max_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("rob_max_duration"), dAtA)
	}

	if pArSeR.KeyExist("lost_prosperity_duration") {
		dAtA.LostProsperityDuration, err = config.ParseDuration(pArSeR.String("lost_prosperity_duration"))
	} else {
		dAtA.LostProsperityDuration, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[lost_prosperity_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("lost_prosperity_duration"), dAtA)
	}

	dAtA.LostProsperityPerDuration = pArSeR.Uint64("lost_prosperity_per_duration")
	dAtA.FirstLoseProsperity = 0
	if pArSeR.KeyExist("first_lose_prosperity") {
		dAtA.FirstLoseProsperity = pArSeR.Uint64("first_lose_prosperity")
	}

	if pArSeR.KeyExist("tick_duration") {
		dAtA.TickDuration, err = config.ParseDuration(pArSeR.String("tick_duration"))
	} else {
		dAtA.TickDuration, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tick_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tick_duration"), dAtA)
	}

	dAtA.TickIcon = pArSeR.StringArray("tick_icon", "", false)
	// releated field: TickPrize
	// releated field: FirstPrize
	// releated field: TickPlunder
	// releated field: FirstPlunder
	// releated field: TickConditionPlunder
	// releated field: FirstConditionPlunder
	// releated field: ShowPrize
	// releated field: ShowSubPrize
	dAtA.MaxRobbers = pArSeR.Uint64("max_robbers")
	dAtA.DestroyWhenLose = false
	if pArSeR.KeyExist("destroy_when_lose") {
		dAtA.DestroyWhenLose = pArSeR.Bool("destroy_when_lose")
	}

	return dAtA, nil
}

var vAlIdAtOrNpcBaseData = map[string]*config.Validator{

	"id":                           config.ParseValidator("int>0", "", false, nil, nil),
	"name":                         config.ParseValidator("string", "", false, nil, []string{""}),
	"npc":                          config.ParseValidator("string", "", false, nil, nil),
	"base_level":                   config.ParseValidator("int>0", "", false, nil, nil),
	"model":                        config.ParseValidator("string", "", false, nil, []string{""}),
	"def_model":                    config.ParseValidator("string", "", false, nil, []string{""}),
	"prosperity_capcity":           config.ParseValidator("int>0", "", false, nil, nil),
	"rob_max_duration":             config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"lost_prosperity_duration":     config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"lost_prosperity_per_duration": config.ParseValidator("uint", "", false, nil, nil),
	"first_lose_prosperity":        config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"tick_duration":                config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"tick_icon":                    config.ParseValidator("string", "", true, nil, nil),
	"tick_prize":                   config.ParseValidator("string", "", false, nil, nil),
	"first_prize":                  config.ParseValidator("string", "", false, nil, nil),
	"tick_plunder":                 config.ParseValidator("string", "", false, nil, nil),
	"first_plunder":                config.ParseValidator("string", "", false, nil, nil),
	"tick_condition_plunder":       config.ParseValidator("string", "", false, nil, nil),
	"first_condition_plunder":      config.ParseValidator("string", "", false, nil, nil),
	"show_prize":                   config.ParseValidator("string", "", false, nil, nil),
	"show_sub_prize":               config.ParseValidator("string", "", false, nil, nil),
	"max_robbers":                  config.ParseValidator("uint", "", false, nil, nil),
	"destroy_when_lose":            config.ParseValidator("bool", "", false, nil, []string{"false"}),
}

func (dAtA *NpcBaseData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *NpcBaseData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *NpcBaseData) Encode() *shared_proto.NpcBaseDataProto {
	out := &shared_proto.NpcBaseDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	if dAtA.Npc != nil {
		out.Npc = dAtA.Npc.Encode()
	}
	out.BaseLevel = config.U64ToI32(dAtA.BaseLevel)
	out.Model = dAtA.Model
	out.DefModel = dAtA.DefModel
	out.ProsperityCapcity = config.U64ToI32(dAtA.ProsperityCapcity)
	out.RobMaxDuration = config.Duration2I32Seconds(dAtA.RobMaxDuration)
	out.LostProsperityDuration = config.Duration2I32Seconds(dAtA.LostProsperityDuration)
	out.TickDuration = config.Duration2I32Seconds(dAtA.TickDuration)
	out.TickIcon = dAtA.TickIcon
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.Encode()
	}
	if dAtA.ShowSubPrize != nil {
		out.ShowSubPrize = dAtA.ShowSubPrize.Encode()
	}

	return out
}

func ArrayEncodeNpcBaseData(datas []*NpcBaseData) []*shared_proto.NpcBaseDataProto {

	out := make([]*shared_proto.NpcBaseDataProto, 0, len(datas))
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

func (dAtA *NpcBaseData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Npc = cOnFigS.GetMonsterMasterData(pArSeR.Uint64("npc"))
	if dAtA.Npc == nil {
		return errors.Errorf("%s 配置的关联字段[npc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("npc"), *pArSeR)
	}

	dAtA.TickPrize = cOnFigS.GetPrize(pArSeR.Int("tick_prize"))
	if dAtA.TickPrize == nil && pArSeR.Int("tick_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[tick_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("tick_prize"), *pArSeR)
	}

	dAtA.FirstPrize = cOnFigS.GetPrize(pArSeR.Int("first_prize"))
	if dAtA.FirstPrize == nil && pArSeR.Int("first_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[first_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_prize"), *pArSeR)
	}

	dAtA.TickPlunder = cOnFigS.GetPlunder(pArSeR.Uint64("tick_plunder"))
	if dAtA.TickPlunder == nil && pArSeR.Uint64("tick_plunder") != 0 {
		return errors.Errorf("%s 配置的关联字段[tick_plunder] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("tick_plunder"), *pArSeR)
	}

	dAtA.FirstPlunder = cOnFigS.GetPlunder(pArSeR.Uint64("first_plunder"))
	if dAtA.FirstPlunder == nil && pArSeR.Uint64("first_plunder") != 0 {
		return errors.Errorf("%s 配置的关联字段[first_plunder] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_plunder"), *pArSeR)
	}

	dAtA.TickConditionPlunder = cOnFigS.GetConditionPlunder(pArSeR.Uint64("tick_condition_plunder"))
	if dAtA.TickConditionPlunder == nil && pArSeR.Uint64("tick_condition_plunder") != 0 {
		return errors.Errorf("%s 配置的关联字段[tick_condition_plunder] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("tick_condition_plunder"), *pArSeR)
	}

	dAtA.FirstConditionPlunder = cOnFigS.GetConditionPlunder(pArSeR.Uint64("first_condition_plunder"))
	if dAtA.FirstConditionPlunder == nil && pArSeR.Uint64("first_condition_plunder") != 0 {
		return errors.Errorf("%s 配置的关联字段[first_condition_plunder] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_condition_plunder"), *pArSeR)
	}

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil && pArSeR.Int("show_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[show_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
	}

	dAtA.ShowSubPrize = cOnFigS.GetPrize(pArSeR.Int("show_sub_prize"))
	if dAtA.ShowSubPrize == nil && pArSeR.Int("show_sub_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[show_sub_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_sub_prize"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetConditionPlunder(uint64) *resdata.ConditionPlunder
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
	GetNpcBaseData(uint64) *NpcBaseData
	GetPlunder(uint64) *resdata.Plunder
	GetPrize(int) *resdata.Prize
}
