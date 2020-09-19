// AUTO_GEN, DONT MODIFY!!!
package domestic_data

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/icon"
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

// start with BaseLevelData ----------------------------------

func LoadBaseLevelData(gos *config.GameObjects) (map[uint64]*BaseLevelData, map[*BaseLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BaseLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BaseLevelData, len(lIsT))
	pArSeRmAp := make(map[*BaseLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBaseLevelData) {
			continue
		}

		dAtA, err := NewBaseLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedBaseLevelData(dAtAmAp map[*BaseLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BaseLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBaseLevelDataKeyArray(datas []*BaseLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewBaseLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BaseLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBaseLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BaseLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Prosperity = pArSeR.Uint64("upgrade_prosperity")
	// skip field: UpgradeProsperity
	// skip field: UnlockGuanFuLevel
	dAtA.UnlockPowerRange = ""
	if pArSeR.KeyExist("unlock_power_range") {
		dAtA.UnlockPowerRange = pArSeR.String("unlock_power_range")
	}

	dAtA.AppearanceRes = ""
	if pArSeR.KeyExist("appearance_res") {
		dAtA.AppearanceRes = pArSeR.String("appearance_res")
	}

	dAtA.TriggerCountdownPrizeProsperity = pArSeR.Uint64("trigger_countdown_prize_prosperity")
	dAtA.AddCountdownPrizeProsperity = pArSeR.Uint64("add_countdown_prize_prosperity")
	// releated field: AddCountdownPrizeDesc

	return dAtA, nil
}

var vAlIdAtOrBaseLevelData = map[string]*config.Validator{

	"level":                              config.ParseValidator("int>0", "", false, nil, nil),
	"upgrade_prosperity":                 config.ParseValidator("int>0", "", false, nil, nil),
	"unlock_power_range":                 config.ParseValidator("string", "", false, nil, []string{""}),
	"appearance_res":                     config.ParseValidator("string", "", false, nil, []string{""}),
	"trigger_countdown_prize_prosperity": config.ParseValidator("uint", "", false, nil, nil),
	"add_countdown_prize_prosperity":     config.ParseValidator("uint", "", false, nil, nil),
	"add_countdown_prize_desc":           config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *BaseLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BaseLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BaseLevelData) Encode() *shared_proto.BaseLevelProto {
	out := &shared_proto.BaseLevelProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.UpgradeProsperity = config.U64ToI32(dAtA.UpgradeProsperity)
	out.UnlockGuanFuLevel = config.U64ToI32(dAtA.UnlockGuanFuLevel)
	out.UnlockPowerRange = dAtA.UnlockPowerRange
	out.AppearanceRes = dAtA.AppearanceRes

	return out
}

func ArrayEncodeBaseLevelData(datas []*BaseLevelData) []*shared_proto.BaseLevelProto {

	out := make([]*shared_proto.BaseLevelProto, 0, len(datas))
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

func (dAtA *BaseLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.AddCountdownPrizeDesc = cOnFigS.GetCountdownPrizeDescData(pArSeR.Uint64("add_countdown_prize_desc"))
	if dAtA.AddCountdownPrizeDesc == nil {
		return errors.Errorf("%s 配置的关联字段[add_countdown_prize_desc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("add_countdown_prize_desc"), *pArSeR)
	}

	return nil
}

// start with BuildingData ----------------------------------

func LoadBuildingData(gos *config.GameObjects) (map[uint64]*BuildingData, map[*BuildingData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BuildingDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BuildingData, len(lIsT))
	pArSeRmAp := make(map[*BuildingData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBuildingData) {
			continue
		}

		dAtA, err := NewBuildingData(fIlEnAmE, pArSeR)
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

func SetRelatedBuildingData(dAtAmAp map[*BuildingData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BuildingDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBuildingDataKeyArray(datas []*BuildingData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBuildingData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BuildingData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBuildingData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BuildingData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Type = shared_proto.BuildingType(shared_proto.BuildingType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.BuildingType(i)
	}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.WorkTime, err = config.ParseDuration(pArSeR.String("work_time"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[work_time] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("work_time"), dAtA)
	}

	dAtA.Prosperity = pArSeR.Uint64("prosperity")
	dAtA.HeroExp = pArSeR.Uint64("hero_exp")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Tips = ""
	if pArSeR.KeyExist("tips") {
		dAtA.Tips = pArSeR.String("tips")
	}

	// releated field: Icon
	dAtA.Model = ""
	if pArSeR.KeyExist("model") {
		dAtA.Model = pArSeR.String("model")
	}

	dAtA.EffectDesc = ""
	if pArSeR.KeyExist("effect_desc") {
		dAtA.EffectDesc = pArSeR.String("effect_desc")
	}

	// releated field: Effect
	dAtA.Notice = ""
	if pArSeR.KeyExist("notice") {
		dAtA.Notice = pArSeR.String("notice")
	}

	// releated field: Cost
	// releated field: RequireIds
	// releated field: BaseLevel

	return dAtA, nil
}

var vAlIdAtOrBuildingData = map[string]*config.Validator{

	"id":          config.ParseValidator("int>0", "", false, nil, nil),
	"type":        config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.BuildingType_value, 0), nil),
	"level":       config.ParseValidator("int>0", "", false, nil, nil),
	"work_time":   config.ParseValidator("string", "", false, nil, nil),
	"prosperity":  config.ParseValidator("uint", "", false, nil, nil),
	"hero_exp":    config.ParseValidator("uint", "", false, nil, nil),
	"desc":        config.ParseValidator("string", "", false, nil, nil),
	"tips":        config.ParseValidator("string", "", false, nil, []string{""}),
	"icon":        config.ParseValidator("string", "", false, nil, []string{"小卒"}),
	"model":       config.ParseValidator("string", "", false, nil, []string{""}),
	"effect_desc": config.ParseValidator("string", "", false, nil, []string{""}),
	"effect":      config.ParseValidator("string", "", false, nil, nil),
	"notice":      config.ParseValidator("string", "", false, nil, []string{""}),
	"cost":        config.ParseValidator("string", "", false, nil, nil),
	"require_ids": config.ParseValidator("string", "", true, nil, nil),
	"base_level":  config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *BuildingData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BuildingData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BuildingData) Encode() *shared_proto.BuildingDataProto {
	out := &shared_proto.BuildingDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Type = dAtA.Type
	out.Level = config.U64ToI32(dAtA.Level)
	out.WorkTime = config.Duration2I32Seconds(dAtA.WorkTime)
	out.Prosperity = config.U64ToI32(dAtA.Prosperity)
	out.HeroExp = config.U64ToI32(dAtA.HeroExp)
	out.Desc = dAtA.Desc
	out.Tips = dAtA.Tips
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.Model = dAtA.Model
	out.EffectDesc = dAtA.EffectDesc
	if dAtA.Effect != nil {
		out.Effect = dAtA.Effect.Encode()
	}
	out.Notice = dAtA.Notice
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	if dAtA.RequireIds != nil {
		out.RequireIds = config.U64a2I32a(GetBuildingDataKeyArray(dAtA.RequireIds))
	}
	if dAtA.BaseLevel != nil {
		out.BaseLevel = config.U64ToI32(dAtA.BaseLevel.Level)
	}

	return out
}

func ArrayEncodeBuildingData(datas []*BuildingData) []*shared_proto.BuildingDataProto {

	out := make([]*shared_proto.BuildingDataProto, 0, len(datas))
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

func (dAtA *BuildingData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
		dAtA.Icon = cOnFigS.GetIcon("小卒")
	}
	if dAtA.Icon == nil {
		return errors.Errorf("%s 配置的关联字段[icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("icon"), *pArSeR)
	}

	dAtA.Effect = cOnFigS.GetBuildingEffectData(pArSeR.Int("effect"))
	if dAtA.Effect == nil && pArSeR.Int("effect") != 0 {
		return errors.Errorf("%s 配置的关联字段[effect] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("effect"), *pArSeR)
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("require_ids", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetBuildingData(v)
		if obj != nil {
			dAtA.RequireIds = append(dAtA.RequireIds, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[require_ids] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("require_ids"), *pArSeR)
		}
	}

	dAtA.BaseLevel = cOnFigS.GetBaseLevelData(pArSeR.Uint64("base_level"))
	if dAtA.BaseLevel == nil && pArSeR.Uint64("base_level") != 0 {
		return errors.Errorf("%s 配置的关联字段[base_level] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("base_level"), *pArSeR)
	}

	return nil
}

// start with BuildingLayoutData ----------------------------------

func LoadBuildingLayoutData(gos *config.GameObjects) (map[uint64]*BuildingLayoutData, map[*BuildingLayoutData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BuildingLayoutDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BuildingLayoutData, len(lIsT))
	pArSeRmAp := make(map[*BuildingLayoutData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBuildingLayoutData) {
			continue
		}

		dAtA, err := NewBuildingLayoutData(fIlEnAmE, pArSeR)
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

func SetRelatedBuildingLayoutData(dAtAmAp map[*BuildingLayoutData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BuildingLayoutDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBuildingLayoutDataKeyArray(datas []*BuildingLayoutData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBuildingLayoutData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BuildingLayoutData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBuildingLayoutData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BuildingLayoutData{}

	dAtA.Id = pArSeR.Uint64("id")
	for _, v := range pArSeR.StringArray("building_type", "", false) {
		x := shared_proto.BuildingType(shared_proto.BuildingType_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.BuildingType(i)
		}
		dAtA.BuildingType = append(dAtA.BuildingType, x)
	}

	dAtA.RequireBaseLevel = pArSeR.Uint64("require_base_level")
	dAtA.RegionOffsetX = pArSeR.Int("region_offset_x")
	dAtA.RegionOffsetY = pArSeR.Int("region_offset_y")
	dAtA.IgnoreConflict = false
	if pArSeR.KeyExist("ignore_conflict") {
		dAtA.IgnoreConflict = pArSeR.Bool("ignore_conflict")
	}

	return dAtA, nil
}

var vAlIdAtOrBuildingLayoutData = map[string]*config.Validator{

	"id":                 config.ParseValidator("int>0", "", false, nil, nil),
	"building_type":      config.ParseValidator("string,notAllNil", "", true, config.EnumMapKeys(shared_proto.BuildingType_value, 0), nil),
	"require_base_level": config.ParseValidator("int>0", "", false, nil, nil),
	"region_offset_x":    config.ParseValidator("int", "", false, nil, nil),
	"region_offset_y":    config.ParseValidator("int", "", false, nil, nil),
	"ignore_conflict":    config.ParseValidator("bool", "", false, nil, []string{"false"}),
}

func (dAtA *BuildingLayoutData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BuildingLayoutData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BuildingLayoutData) Encode() *shared_proto.BuildingLayoutProto {
	out := &shared_proto.BuildingLayoutProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Building = dAtA.BuildingType
	out.RequireBaseLevel = config.U64ToI32(dAtA.RequireBaseLevel)
	out.RegionOffsetX = int32(dAtA.RegionOffsetX)
	out.RegionOffsetY = int32(dAtA.RegionOffsetY)
	out.IgnoreConflict = dAtA.IgnoreConflict

	return out
}

func ArrayEncodeBuildingLayoutData(datas []*BuildingLayoutData) []*shared_proto.BuildingLayoutProto {

	out := make([]*shared_proto.BuildingLayoutProto, 0, len(datas))
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

func (dAtA *BuildingLayoutData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with BuildingLayoutMiscData ----------------------------------

func LoadBuildingLayoutMiscData(gos *config.GameObjects) (*BuildingLayoutMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.BuildingLayoutMiscDataPath
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

	dAtA, err := NewBuildingLayoutMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedBuildingLayoutMiscData(gos *config.GameObjects, dAtA *BuildingLayoutMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BuildingLayoutMiscDataPath
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

func NewBuildingLayoutMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BuildingLayoutMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBuildingLayoutMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BuildingLayoutMiscData{}

	dAtA.Gold, err = data.ParseAmount(pArSeR.String("gold"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[gold] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("gold"), dAtA)
	}

	dAtA.Food, err = data.ParseAmount(pArSeR.String("food"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[food] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("food"), dAtA)
	}

	dAtA.Wood, err = data.ParseAmount(pArSeR.String("wood"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[wood] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("wood"), dAtA)
	}

	dAtA.Stone, err = data.ParseAmount(pArSeR.String("stone"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[stone] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("stone"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrBuildingLayoutMiscData = map[string]*config.Validator{

	"gold":  config.ParseValidator("uint", "", false, nil, nil),
	"food":  config.ParseValidator("uint", "", false, nil, nil),
	"wood":  config.ParseValidator("uint", "", false, nil, nil),
	"stone": config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *BuildingLayoutMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with BuildingUnlockData ----------------------------------

func LoadBuildingUnlockData(gos *config.GameObjects) (map[uint64]*BuildingUnlockData, map[*BuildingUnlockData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BuildingUnlockDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BuildingUnlockData, len(lIsT))
	pArSeRmAp := make(map[*BuildingUnlockData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBuildingUnlockData) {
			continue
		}

		dAtA, err := NewBuildingUnlockData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.IntBuildingType
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[IntBuildingType], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedBuildingUnlockData(dAtAmAp map[*BuildingUnlockData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BuildingUnlockDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBuildingUnlockDataKeyArray(datas []*BuildingUnlockData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.IntBuildingType)
		}
	}

	return out
}

func NewBuildingUnlockData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BuildingUnlockData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBuildingUnlockData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BuildingUnlockData{}

	dAtA.BuildingType = shared_proto.BuildingType(shared_proto.BuildingType_value[strings.ToUpper(pArSeR.String("building_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("building_type"), 10, 32); err == nil {
		dAtA.BuildingType = shared_proto.BuildingType(i)
	}

	dAtA.Desc = ""
	if pArSeR.KeyExist("desc") {
		dAtA.Desc = pArSeR.String("desc")
	}

	dAtA.Icon = "icon"
	if pArSeR.KeyExist("icon") {
		dAtA.Icon = pArSeR.String("icon")
	}

	dAtA.NotifyOrder = 1
	if pArSeR.KeyExist("notify_order") {
		dAtA.NotifyOrder = pArSeR.Uint64("notify_order")
	}

	// releated field: GuanFuLevel
	dAtA.HeroLevel = pArSeR.Uint64("hero_level")
	dAtA.MainTaskSequence = pArSeR.Uint64("main_task_sequence")
	dAtA.BaYeStage = pArSeR.Uint64("ba_ye_stage")
	// skip field: UnlockBuildingData

	// calculate fields
	dAtA.IntBuildingType = uint64(dAtA.BuildingType)

	return dAtA, nil
}

var vAlIdAtOrBuildingUnlockData = map[string]*config.Validator{

	"building_type":      config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.BuildingType_value, 0), nil),
	"desc":               config.ParseValidator("string", "", false, nil, []string{""}),
	"icon":               config.ParseValidator("string", "", false, nil, []string{"icon"}),
	"notify_order":       config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"guan_fu_level":      config.ParseValidator("string", "", false, nil, nil),
	"hero_level":         config.ParseValidator("uint", "", false, nil, nil),
	"main_task_sequence": config.ParseValidator("uint", "", false, nil, nil),
	"ba_ye_stage":        config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *BuildingUnlockData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BuildingUnlockData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BuildingUnlockData) Encode() *shared_proto.BuildingUnlockDataProto {
	out := &shared_proto.BuildingUnlockDataProto{}
	out.BuildingType = dAtA.BuildingType
	out.Desc = dAtA.Desc
	out.Icon = dAtA.Icon
	out.NotifyOrder = config.U64ToI32(dAtA.NotifyOrder)
	if dAtA.GuanFuLevel != nil {
		out.GuanFuLevel = config.U64ToI32(dAtA.GuanFuLevel.Level)
	}
	out.HeroLevel = config.U64ToI32(dAtA.HeroLevel)
	out.MainTaskSequence = config.U64ToI32(dAtA.MainTaskSequence)
	out.BaYeStage = config.U64ToI32(dAtA.BaYeStage)

	return out
}

func ArrayEncodeBuildingUnlockData(datas []*BuildingUnlockData) []*shared_proto.BuildingUnlockDataProto {

	out := make([]*shared_proto.BuildingUnlockDataProto, 0, len(datas))
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

func (dAtA *BuildingUnlockData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.GuanFuLevel = cOnFigS.GetBuildingData(pArSeR.Uint64("guan_fu_level"))
	if dAtA.GuanFuLevel == nil && pArSeR.Uint64("guan_fu_level") != 0 {
		return errors.Errorf("%s 配置的关联字段[guan_fu_level] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guan_fu_level"), *pArSeR)
	}

	return nil
}

// start with CityEventData ----------------------------------

func LoadCityEventData(gos *config.GameObjects) (map[uint64]*CityEventData, map[*CityEventData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CityEventDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CityEventData, len(lIsT))
	pArSeRmAp := make(map[*CityEventData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCityEventData) {
			continue
		}

		dAtA, err := NewCityEventData(fIlEnAmE, pArSeR)
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

func SetRelatedCityEventData(dAtAmAp map[*CityEventData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CityEventDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCityEventDataKeyArray(datas []*CityEventData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCityEventData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CityEventData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCityEventData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CityEventData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Weight = 1
	if pArSeR.KeyExist("weight") {
		dAtA.Weight = pArSeR.Uint64("weight")
	}

	// releated field: Cost
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrCityEventData = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"desc":   config.ParseValidator("string>0", "", false, nil, nil),
	"weight": config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"cost":   config.ParseValidator("string", "", false, nil, nil),
	"prize":  config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *CityEventData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CityEventData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CityEventData) Encode() *shared_proto.CityEventDataProto {
	out := &shared_proto.CityEventDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Desc = dAtA.Desc
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeCityEventData(datas []*CityEventData) []*shared_proto.CityEventDataProto {

	out := make([]*shared_proto.CityEventDataProto, 0, len(datas))
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

func (dAtA *CityEventData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Cost = cOnFigS.GetCombineCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with CityEventLevelData ----------------------------------

func LoadCityEventLevelData(gos *config.GameObjects) (map[uint64]*CityEventLevelData, map[*CityEventLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CityEventLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CityEventLevelData, len(lIsT))
	pArSeRmAp := make(map[*CityEventLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCityEventLevelData) {
			continue
		}

		dAtA, err := NewCityEventLevelData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.BaseLevel
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[BaseLevel], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedCityEventLevelData(dAtAmAp map[*CityEventLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CityEventLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCityEventLevelDataKeyArray(datas []*CityEventLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.BaseLevel)
		}
	}

	return out
}

func NewCityEventLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CityEventLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCityEventLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CityEventLevelData{}

	dAtA.BaseLevel = pArSeR.Uint64("base_level")
	// skip field: BaseLevelData
	// releated field: EventDatas

	return dAtA, nil
}

var vAlIdAtOrCityEventLevelData = map[string]*config.Validator{

	"base_level":  config.ParseValidator("int>0", "", false, nil, nil),
	"event_datas": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *CityEventLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("event_datas", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetCityEventData(v)
		if obj != nil {
			dAtA.EventDatas = append(dAtA.EventDatas, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[event_datas] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("event_datas"), *pArSeR)
		}
	}

	return nil
}

// start with CityEventMiscData ----------------------------------

func LoadCityEventMiscData(gos *config.GameObjects) (*CityEventMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.CityEventMiscDataPath
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

	dAtA, err := NewCityEventMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedCityEventMiscData(gos *config.GameObjects, dAtA *CityEventMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CityEventMiscDataPath
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

func NewCityEventMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CityEventMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCityEventMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CityEventMiscData{}

	dAtA.MaxTimes = pArSeR.Uint64("max_times")
	dAtA.RecoverDuration, err = config.ParseDuration(pArSeR.String("recover_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[recover_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("recover_duration"), dAtA)
	}

	// skip field: UnlockBaseLevel

	return dAtA, nil
}

var vAlIdAtOrCityEventMiscData = map[string]*config.Validator{

	"max_times":        config.ParseValidator("int>0", "", false, nil, nil),
	"recover_duration": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *CityEventMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CityEventMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CityEventMiscData) Encode() *shared_proto.CityEventMiscProto {
	out := &shared_proto.CityEventMiscProto{}
	out.MaxTimes = config.U64ToI32(dAtA.MaxTimes)
	out.RecoverDuration = config.Duration2I32Seconds(dAtA.RecoverDuration)
	out.UnlockBaseLevel = config.U64ToI32(dAtA.UnlockBaseLevel)

	return out
}

func ArrayEncodeCityEventMiscData(datas []*CityEventMiscData) []*shared_proto.CityEventMiscProto {

	out := make([]*shared_proto.CityEventMiscProto, 0, len(datas))
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

func (dAtA *CityEventMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with CombineCost ----------------------------------

func LoadCombineCost(gos *config.GameObjects) (map[int]*CombineCost, map[*CombineCost]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CombineCostPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[int]*CombineCost, len(lIsT))
	pArSeRmAp := make(map[*CombineCost]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCombineCost) {
			continue
		}

		dAtA, err := NewCombineCost(fIlEnAmE, pArSeR)
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

func SetRelatedCombineCost(dAtAmAp map[*CombineCost]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CombineCostPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCombineCostKeyArray(datas []*CombineCost) []int {

	out := make([]int, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCombineCost(fIlEnAmE string, pArSeR *config.ObjectParser) (*CombineCost, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCombineCost)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CombineCost{}

	dAtA.Id = pArSeR.Int("id")
	// releated field: Cost
	if pArSeR.KeyExist("building_worker_time") {
		dAtA.BuildingWorkerTime, err = config.ParseDuration(pArSeR.String("building_worker_time"))
	} else {
		dAtA.BuildingWorkerTime, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[building_worker_time] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("building_worker_time"), dAtA)
	}

	if pArSeR.KeyExist("tech_worker_time") {
		dAtA.TechWorkerTime, err = config.ParseDuration(pArSeR.String("tech_worker_time"))
	} else {
		dAtA.TechWorkerTime, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tech_worker_time] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tech_worker_time"), dAtA)
	}

	dAtA.Soldier = pArSeR.Uint64("soldier")
	dAtA.InvadeTimes = pArSeR.Uint64("invade_times")

	return dAtA, nil
}

var vAlIdAtOrCombineCost = map[string]*config.Validator{

	"id":                   config.ParseValidator("int>0", "", false, nil, nil),
	"cost":                 config.ParseValidator("string", "", false, nil, nil),
	"building_worker_time": config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"tech_worker_time":     config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"soldier":              config.ParseValidator("uint", "", false, nil, nil),
	"invade_times":         config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *CombineCost) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CombineCost) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CombineCost) Encode() *shared_proto.CombineCostProto {
	out := &shared_proto.CombineCostProto{}
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	out.BuildingWorkerTime = config.Duration2I32Seconds(dAtA.BuildingWorkerTime)
	out.TechWorkerTime = config.Duration2I32Seconds(dAtA.TechWorkerTime)
	out.Soldier = config.U64ToI32(dAtA.Soldier)
	out.InvadeTimes = config.U64ToI32(dAtA.InvadeTimes)

	return out
}

func ArrayEncodeCombineCost(datas []*CombineCost) []*shared_proto.CombineCostProto {

	out := make([]*shared_proto.CombineCostProto, 0, len(datas))
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

func (dAtA *CombineCost) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	return nil
}

// start with CountdownPrizeData ----------------------------------

func LoadCountdownPrizeData(gos *config.GameObjects) (map[uint64]*CountdownPrizeData, map[*CountdownPrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CountdownPrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CountdownPrizeData, len(lIsT))
	pArSeRmAp := make(map[*CountdownPrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCountdownPrizeData) {
			continue
		}

		dAtA, err := NewCountdownPrizeData(fIlEnAmE, pArSeR)
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

func SetRelatedCountdownPrizeData(dAtAmAp map[*CountdownPrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CountdownPrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCountdownPrizeDataKeyArray(datas []*CountdownPrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCountdownPrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CountdownPrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCountdownPrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CountdownPrizeData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Plunder
	dAtA.WaitDuration, err = config.ParseDuration(pArSeR.String("wait_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[wait_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("wait_duration"), dAtA)
	}

	// releated field: Descs

	return dAtA, nil
}

var vAlIdAtOrCountdownPrizeData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"plunder":       config.ParseValidator("string", "", false, nil, nil),
	"wait_duration": config.ParseValidator("string", "", false, nil, nil),
	"descs":         config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *CountdownPrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Plunder = cOnFigS.GetPlunder(pArSeR.Uint64("plunder"))
	if dAtA.Plunder == nil {
		return errors.Errorf("%s 配置的关联字段[plunder] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("descs", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetCountdownPrizeDescData(v)
		if obj != nil {
			dAtA.Descs = append(dAtA.Descs, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[descs] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("descs"), *pArSeR)
		}
	}

	return nil
}

// start with CountdownPrizeDescData ----------------------------------

func LoadCountdownPrizeDescData(gos *config.GameObjects) (map[uint64]*CountdownPrizeDescData, map[*CountdownPrizeDescData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CountdownPrizeDescDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CountdownPrizeDescData, len(lIsT))
	pArSeRmAp := make(map[*CountdownPrizeDescData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCountdownPrizeDescData) {
			continue
		}

		dAtA, err := NewCountdownPrizeDescData(fIlEnAmE, pArSeR)
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

func SetRelatedCountdownPrizeDescData(dAtAmAp map[*CountdownPrizeDescData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CountdownPrizeDescDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCountdownPrizeDescDataKeyArray(datas []*CountdownPrizeDescData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCountdownPrizeDescData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CountdownPrizeDescData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCountdownPrizeDescData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CountdownPrizeDescData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Desc = pArSeR.String("desc")

	return dAtA, nil
}

var vAlIdAtOrCountdownPrizeDescData = map[string]*config.Validator{

	"id":   config.ParseValidator("int>0", "", false, nil, nil),
	"desc": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *CountdownPrizeDescData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CountdownPrizeDescData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CountdownPrizeDescData) Encode() *shared_proto.CountdownPrizeDescDataProto {
	out := &shared_proto.CountdownPrizeDescDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Desc = dAtA.Desc

	return out
}

func ArrayEncodeCountdownPrizeDescData(datas []*CountdownPrizeDescData) []*shared_proto.CountdownPrizeDescDataProto {

	out := make([]*shared_proto.CountdownPrizeDescDataProto, 0, len(datas))
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

func (dAtA *CountdownPrizeDescData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GuanFuLevelData ----------------------------------

func LoadGuanFuLevelData(gos *config.GameObjects) (map[uint64]*GuanFuLevelData, map[*GuanFuLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuanFuLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuanFuLevelData, len(lIsT))
	pArSeRmAp := make(map[*GuanFuLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuanFuLevelData) {
			continue
		}

		dAtA, err := NewGuanFuLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedGuanFuLevelData(dAtAmAp map[*GuanFuLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuanFuLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuanFuLevelDataKeyArray(datas []*GuanFuLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewGuanFuLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuanFuLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuanFuLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuanFuLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.RestoreProsperity = 10
	if pArSeR.KeyExist("restore_prosperity") {
		dAtA.RestoreProsperity = pArSeR.Uint64("restore_prosperity")
	}

	if pArSeR.KeyExist("move_base_restore_home_prosperity_duration") {
		dAtA.MoveBaseRestoreHomeProsperityDuration, err = config.ParseDuration(pArSeR.String("move_base_restore_home_prosperity_duration"))
	} else {
		dAtA.MoveBaseRestoreHomeProsperityDuration, err = config.ParseDuration("4h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[move_base_restore_home_prosperity_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("move_base_restore_home_prosperity_duration"), dAtA)
	}

	dAtA.MoveBaseRestoreHomeProsperity = 100
	if pArSeR.KeyExist("move_base_restore_home_prosperity") {
		dAtA.MoveBaseRestoreHomeProsperity = pArSeR.Uint64("move_base_restore_home_prosperity")
	}

	if pArSeR.KeyExist("buy_prosperity_restore_duration_with1_yuanbao") {
		dAtA.BuyProsperityRestoreDurationWith1Yuanbao, err = config.ParseDuration(pArSeR.String("buy_prosperity_restore_duration_with1_yuanbao"))
	} else {
		dAtA.BuyProsperityRestoreDurationWith1Yuanbao, err = config.ParseDuration("1m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[buy_prosperity_restore_duration_with1_yuanbao] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("buy_prosperity_restore_duration_with1_yuanbao"), dAtA)
	}

	if pArSeR.KeyExist("buy_prosperity_restore_duration_with1_cost") {
		dAtA.BuyProsperityRestoreDurationWith1Cost, err = config.ParseDuration(pArSeR.String("buy_prosperity_restore_duration_with1_cost"))
	} else {
		dAtA.BuyProsperityRestoreDurationWith1Cost, err = config.ParseDuration("1m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[buy_prosperity_restore_duration_with1_cost] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("buy_prosperity_restore_duration_with1_cost"), dAtA)
	}

	// releated field: WorkshopPrize

	return dAtA, nil
}

var vAlIdAtOrGuanFuLevelData = map[string]*config.Validator{

	"level":                                         config.ParseValidator("int>0", "", false, nil, nil),
	"restore_prosperity":                            config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"move_base_restore_home_prosperity_duration":    config.ParseValidator("string", "", false, nil, []string{"4h"}),
	"move_base_restore_home_prosperity":             config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"buy_prosperity_restore_duration_with1_yuanbao": config.ParseValidator("string", "", false, nil, []string{"1m"}),
	"buy_prosperity_restore_duration_with1_cost":    config.ParseValidator("string", "", false, nil, []string{"1m"}),
	"workshop_prize":                                config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *GuanFuLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuanFuLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuanFuLevelData) Encode() *shared_proto.GuanFuLevelProto {
	out := &shared_proto.GuanFuLevelProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.RestoreProsperity = config.U64ToI32(dAtA.RestoreProsperity)
	out.MoveBaseRestoreHomeProsperityDuration = config.Duration2I32Seconds(dAtA.MoveBaseRestoreHomeProsperityDuration)
	out.MoveBaseRestoreHomeProsperity = config.U64ToI32(dAtA.MoveBaseRestoreHomeProsperity)
	out.BuyProsperityRestoreDurationWith1Yuanbao = config.Duration2I32Seconds(dAtA.BuyProsperityRestoreDurationWith1Yuanbao)
	out.BuyProsperityRestoreDurationWith1Cost = config.Duration2I32Seconds(dAtA.BuyProsperityRestoreDurationWith1Cost)
	if dAtA.WorkshopPrize != nil {
		out.WorkshopPrize = dAtA.WorkshopPrize.Encode()
	}

	return out
}

func ArrayEncodeGuanFuLevelData(datas []*GuanFuLevelData) []*shared_proto.GuanFuLevelProto {

	out := make([]*shared_proto.GuanFuLevelProto, 0, len(datas))
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

func (dAtA *GuanFuLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.WorkshopPrize = cOnFigS.GetPrize(pArSeR.Int("workshop_prize"))
	if dAtA.WorkshopPrize == nil {
		return errors.Errorf("%s 配置的关联字段[workshop_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("workshop_prize"), *pArSeR)
	}

	return nil
}

// start with MainCityMiscData ----------------------------------

func LoadMainCityMiscData(gos *config.GameObjects) (*MainCityMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.MainCityMiscDataPath
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

	dAtA, err := NewMainCityMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedMainCityMiscData(gos *config.GameObjects, dAtA *MainCityMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MainCityMiscDataPath
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

func NewMainCityMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MainCityMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMainCityMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MainCityMiscData{}

	for _, v := range pArSeR.StringArray("main_city_building_types", "", false) {
		x := shared_proto.BuildingType(shared_proto.BuildingType_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.BuildingType(i)
		}
		dAtA.MainCityBuildingTypes = append(dAtA.MainCityBuildingTypes, x)
	}

	return dAtA, nil
}

var vAlIdAtOrMainCityMiscData = map[string]*config.Validator{

	"main_city_building_types": config.ParseValidator("string,notAllNil", "", true, config.EnumMapKeys(shared_proto.BuildingType_value, 0), nil),
}

func (dAtA *MainCityMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with OuterCityBuildingData ----------------------------------

func LoadOuterCityBuildingData(gos *config.GameObjects) (map[uint64]*OuterCityBuildingData, map[*OuterCityBuildingData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.OuterCityBuildingDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*OuterCityBuildingData, len(lIsT))
	pArSeRmAp := make(map[*OuterCityBuildingData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrOuterCityBuildingData) {
			continue
		}

		dAtA, err := NewOuterCityBuildingData(fIlEnAmE, pArSeR)
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

func SetRelatedOuterCityBuildingData(dAtAmAp map[*OuterCityBuildingData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.OuterCityBuildingDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetOuterCityBuildingDataKeyArray(datas []*OuterCityBuildingData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewOuterCityBuildingData(fIlEnAmE string, pArSeR *config.ObjectParser) (*OuterCityBuildingData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrOuterCityBuildingData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &OuterCityBuildingData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: LockIcon
	// releated field: UnlockIcon
	dAtA.BuildingType = shared_proto.BuildingType(shared_proto.BuildingType_value[strings.ToUpper(pArSeR.String("building_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("building_type"), 10, 32); err == nil {
		dAtA.BuildingType = shared_proto.BuildingType(i)
	}

	dAtA.Level = pArSeR.Uint64("level")
	// skip field: BuildingData
	dAtA.Desc = pArSeR.String("desc")

	return dAtA, nil
}

var vAlIdAtOrOuterCityBuildingData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"lock_icon":     config.ParseValidator("string", "", false, nil, []string{"LockIcon"}),
	"unlock_icon":   config.ParseValidator("string", "", false, nil, []string{"UnlockIcon"}),
	"building_type": config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.BuildingType_value, 0), nil),
	"level":         config.ParseValidator("int>0", "", false, nil, nil),
	"desc":          config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *OuterCityBuildingData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *OuterCityBuildingData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *OuterCityBuildingData) Encode() *shared_proto.OuterCityBuildingDataProto {
	out := &shared_proto.OuterCityBuildingDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.LockIcon != nil {
		out.LockIconId = dAtA.LockIcon.Id
	}
	if dAtA.UnlockIcon != nil {
		out.UnlockIconId = dAtA.UnlockIcon.Id
	}
	if dAtA.BuildingData != nil {
		out.BuildingId = config.U64ToI32(dAtA.BuildingData.Id)
	}
	out.Desc = dAtA.Desc

	return out
}

func ArrayEncodeOuterCityBuildingData(datas []*OuterCityBuildingData) []*shared_proto.OuterCityBuildingDataProto {

	out := make([]*shared_proto.OuterCityBuildingDataProto, 0, len(datas))
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

func (dAtA *OuterCityBuildingData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("lock_icon") {
		dAtA.LockIcon = cOnFigS.GetIcon(pArSeR.String("lock_icon"))
	} else {
		dAtA.LockIcon = cOnFigS.GetIcon("LockIcon")
	}
	if dAtA.LockIcon == nil {
		return errors.Errorf("%s 配置的关联字段[lock_icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("lock_icon"), *pArSeR)
	}

	if pArSeR.KeyExist("unlock_icon") {
		dAtA.UnlockIcon = cOnFigS.GetIcon(pArSeR.String("unlock_icon"))
	} else {
		dAtA.UnlockIcon = cOnFigS.GetIcon("UnlockIcon")
	}
	if dAtA.UnlockIcon == nil {
		return errors.Errorf("%s 配置的关联字段[unlock_icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("unlock_icon"), *pArSeR)
	}

	return nil
}

// start with OuterCityBuildingDescData ----------------------------------

func LoadOuterCityBuildingDescData(gos *config.GameObjects) (map[uint64]*OuterCityBuildingDescData, map[*OuterCityBuildingDescData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.OuterCityBuildingDescDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*OuterCityBuildingDescData, len(lIsT))
	pArSeRmAp := make(map[*OuterCityBuildingDescData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrOuterCityBuildingDescData) {
			continue
		}

		dAtA, err := NewOuterCityBuildingDescData(fIlEnAmE, pArSeR)
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

func SetRelatedOuterCityBuildingDescData(dAtAmAp map[*OuterCityBuildingDescData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.OuterCityBuildingDescDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetOuterCityBuildingDescDataKeyArray(datas []*OuterCityBuildingDescData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewOuterCityBuildingDescData(fIlEnAmE string, pArSeR *config.ObjectParser) (*OuterCityBuildingDescData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrOuterCityBuildingDescData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &OuterCityBuildingDescData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// releated field: Icon
	dAtA.Desc = pArSeR.String("desc")

	return dAtA, nil
}

var vAlIdAtOrOuterCityBuildingDescData = map[string]*config.Validator{

	"id":   config.ParseValidator("int>0", "", false, nil, nil),
	"name": config.ParseValidator("string", "", false, nil, nil),
	"icon": config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"desc": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *OuterCityBuildingDescData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *OuterCityBuildingDescData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *OuterCityBuildingDescData) Encode() *shared_proto.OuterCityBuildingDescDataProto {
	out := &shared_proto.OuterCityBuildingDescDataProto{}
	out.Name = dAtA.Name
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.Desc = dAtA.Desc

	return out
}

func ArrayEncodeOuterCityBuildingDescData(datas []*OuterCityBuildingDescData) []*shared_proto.OuterCityBuildingDescDataProto {

	out := make([]*shared_proto.OuterCityBuildingDescDataProto, 0, len(datas))
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

func (dAtA *OuterCityBuildingDescData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with OuterCityData ----------------------------------

func LoadOuterCityData(gos *config.GameObjects) (map[uint64]*OuterCityData, map[*OuterCityData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.OuterCityDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*OuterCityData, len(lIsT))
	pArSeRmAp := make(map[*OuterCityData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrOuterCityData) {
			continue
		}

		dAtA, err := NewOuterCityData(fIlEnAmE, pArSeR)
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

func SetRelatedOuterCityData(dAtAmAp map[*OuterCityData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.OuterCityDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetOuterCityDataKeyArray(datas []*OuterCityData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewOuterCityData(fIlEnAmE string, pArSeR *config.ObjectParser) (*OuterCityData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrOuterCityData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &OuterCityData{}

	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	// releated field: LockIcon
	// releated field: UnlockIcon
	// skip field: FirstLevelLayoutDatas
	// releated field: Descs
	dAtA.RegionModelRes = pArSeR.String("region_model_res")
	dAtA.Diraction = shared_proto.PosDiraction(shared_proto.PosDiraction_value[strings.ToUpper(pArSeR.String("id"))])
	if i, err := strconv.ParseInt(pArSeR.String("id"), 10, 32); err == nil {
		dAtA.Diraction = shared_proto.PosDiraction(i)
	}

	dAtA.UnlockBeforeImage = pArSeR.String("unlock_before_image")
	dAtA.UnlockAfterImage = pArSeR.String("unlock_after_image")
	dAtA.UnlockDesc = pArSeR.StringArray("unlock_desc", "", false)
	dAtA.RecommandType = 0
	if pArSeR.KeyExist("recommand_type") {
		dAtA.RecommandType = pArSeR.Uint64("recommand_type")
	}

	// calculate fields
	dAtA.Id = uint64(dAtA.Diraction)

	return dAtA, nil
}

var vAlIdAtOrOuterCityData = map[string]*config.Validator{

	"name":             config.ParseValidator("string", "", false, nil, nil),
	"desc":             config.ParseValidator("string", "", false, nil, nil),
	"lock_icon":        config.ParseValidator("string", "", false, nil, []string{"LockIcon"}),
	"unlock_icon":      config.ParseValidator("string", "", false, nil, []string{"UnlockIcon"}),
	"descs":            config.ParseValidator("string", "", true, nil, nil),
	"region_model_res": config.ParseValidator("string", "", false, nil, nil),
	"id":               config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.PosDiraction_value, 0), nil),
	"unlock_before_image": config.ParseValidator("string", "", false, nil, nil),
	"unlock_after_image":  config.ParseValidator("string", "", false, nil, nil),
	"unlock_desc":         config.ParseValidator("string,duplicate,notAllNil", "", true, nil, nil),
	"recommand_type":      config.ParseValidator("int", "", false, nil, []string{"0"}),
}

func (dAtA *OuterCityData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *OuterCityData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *OuterCityData) Encode() *shared_proto.OuterCityDataProto {
	out := &shared_proto.OuterCityDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	if dAtA.LockIcon != nil {
		out.LockIconId = dAtA.LockIcon.Id
	}
	if dAtA.UnlockIcon != nil {
		out.UnlockIconId = dAtA.UnlockIcon.Id
	}
	if dAtA.FirstLevelLayoutDatas != nil {
		out.FirstLevelLayoutDatas = config.U64a2I32a(GetOuterCityLayoutDataKeyArray(dAtA.FirstLevelLayoutDatas))
	}
	if dAtA.Descs != nil {
		out.Descs = ArrayEncodeOuterCityBuildingDescData(dAtA.Descs)
	}
	out.RegionModelRes = dAtA.RegionModelRes
	out.UnlockBeforeImage = dAtA.UnlockBeforeImage
	out.UnlockAfterImage = dAtA.UnlockAfterImage
	out.UnlockDesc = dAtA.UnlockDesc

	return out
}

func ArrayEncodeOuterCityData(datas []*OuterCityData) []*shared_proto.OuterCityDataProto {

	out := make([]*shared_proto.OuterCityDataProto, 0, len(datas))
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

func (dAtA *OuterCityData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("lock_icon") {
		dAtA.LockIcon = cOnFigS.GetIcon(pArSeR.String("lock_icon"))
	} else {
		dAtA.LockIcon = cOnFigS.GetIcon("LockIcon")
	}
	if dAtA.LockIcon == nil {
		return errors.Errorf("%s 配置的关联字段[lock_icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("lock_icon"), *pArSeR)
	}

	if pArSeR.KeyExist("unlock_icon") {
		dAtA.UnlockIcon = cOnFigS.GetIcon(pArSeR.String("unlock_icon"))
	} else {
		dAtA.UnlockIcon = cOnFigS.GetIcon("UnlockIcon")
	}
	if dAtA.UnlockIcon == nil {
		return errors.Errorf("%s 配置的关联字段[unlock_icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("unlock_icon"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("descs", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetOuterCityBuildingDescData(v)
		if obj != nil {
			dAtA.Descs = append(dAtA.Descs, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[descs] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("descs"), *pArSeR)
		}
	}

	return nil
}

// start with OuterCityLayoutData ----------------------------------

func LoadOuterCityLayoutData(gos *config.GameObjects) (map[uint64]*OuterCityLayoutData, map[*OuterCityLayoutData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.OuterCityLayoutDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*OuterCityLayoutData, len(lIsT))
	pArSeRmAp := make(map[*OuterCityLayoutData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrOuterCityLayoutData) {
			continue
		}

		dAtA, err := NewOuterCityLayoutData(fIlEnAmE, pArSeR)
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

func SetRelatedOuterCityLayoutData(dAtAmAp map[*OuterCityLayoutData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.OuterCityLayoutDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetOuterCityLayoutDataKeyArray(datas []*OuterCityLayoutData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewOuterCityLayoutData(fIlEnAmE string, pArSeR *config.ObjectParser) (*OuterCityLayoutData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrOuterCityLayoutData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &OuterCityLayoutData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Layout = pArSeR.Uint64("layout")
	// releated field: OuterCity
	// releated field: ChangeTypeCost
	// releated field: UpgradeRequireLayout
	// releated field: UpgradeRequireIds
	// releated field: PrevLevel
	// skip field: NextLevel
	dAtA.DefaultUnlocked = true
	if pArSeR.KeyExist("default_unlocked") {
		dAtA.DefaultUnlocked = pArSeR.Bool("default_unlocked")
	}

	// skip field: Level
	// releated field: ChangeTypeLevel
	// releated field: MilitaryBuilding
	// releated field: EconomicBuilding
	// skip field: IsMain

	return dAtA, nil
}

var vAlIdAtOrOuterCityLayoutData = map[string]*config.Validator{

	"id":                     config.ParseValidator("int>0", "", false, nil, nil),
	"layout":                 config.ParseValidator("uint", "", false, nil, nil),
	"outer_city":             config.ParseValidator("string", "", false, nil, nil),
	"change_type_cost":       config.ParseValidator("string", "", false, nil, nil),
	"upgrade_require_layout": config.ParseValidator("string", "", false, nil, nil),
	"upgrade_require_ids":    config.ParseValidator("string", "", true, nil, nil),
	"prev_level":             config.ParseValidator("string", "", false, nil, nil),
	"default_unlocked":       config.ParseValidator("bool", "", false, nil, []string{"true"}),
	"change_type_id":         config.ParseValidator("string", "", false, nil, nil),
	"military_building":      config.ParseValidator("string", "", false, nil, nil),
	"economic_building":      config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *OuterCityLayoutData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *OuterCityLayoutData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *OuterCityLayoutData) Encode() *shared_proto.OuterCityLayoutDataProto {
	out := &shared_proto.OuterCityLayoutDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Layout = config.U64ToI32(dAtA.Layout)
	if dAtA.ChangeTypeCost != nil {
		out.ChangeTypeCost = dAtA.ChangeTypeCost.Encode()
	}
	if dAtA.UpgradeRequireLayout != nil {
		out.UpgradeRequireLayoutId = config.U64ToI32(dAtA.UpgradeRequireLayout.Id)
	}
	if dAtA.UpgradeRequireIds != nil {
		out.UpgradeRequireIds = config.U64a2I32a(GetBuildingDataKeyArray(dAtA.UpgradeRequireIds))
	}
	if dAtA.NextLevel != nil {
		out.NextLevel = config.U64ToI32(dAtA.NextLevel.Id)
	}
	if dAtA.MilitaryBuilding != nil {
		out.MilitaryBuilding = config.U64ToI32(dAtA.MilitaryBuilding.Id)
	}
	if dAtA.EconomicBuilding != nil {
		out.EconomicBuilding = config.U64ToI32(dAtA.EconomicBuilding.Id)
	}

	return out
}

func ArrayEncodeOuterCityLayoutData(datas []*OuterCityLayoutData) []*shared_proto.OuterCityLayoutDataProto {

	out := make([]*shared_proto.OuterCityLayoutDataProto, 0, len(datas))
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

func (dAtA *OuterCityLayoutData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.OuterCity = cOnFigS.GetOuterCityData(pArSeR.Uint64("outer_city"))
	if dAtA.OuterCity == nil {
		return errors.Errorf("%s 配置的关联字段[outer_city] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("outer_city"), *pArSeR)
	}

	dAtA.ChangeTypeCost = cOnFigS.GetCost(pArSeR.Int("change_type_cost"))
	if dAtA.ChangeTypeCost == nil && pArSeR.Int("change_type_cost") != 0 {
		return errors.Errorf("%s 配置的关联字段[change_type_cost] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("change_type_cost"), *pArSeR)
	}

	dAtA.UpgradeRequireLayout = cOnFigS.GetOuterCityLayoutData(pArSeR.Uint64("upgrade_require_layout"))
	if dAtA.UpgradeRequireLayout == nil && pArSeR.Uint64("upgrade_require_layout") != 0 {
		return errors.Errorf("%s 配置的关联字段[upgrade_require_layout] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("upgrade_require_layout"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("upgrade_require_ids", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetBuildingData(v)
		if obj != nil {
			dAtA.UpgradeRequireIds = append(dAtA.UpgradeRequireIds, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[upgrade_require_ids] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("upgrade_require_ids"), *pArSeR)
		}
	}

	dAtA.PrevLevel = cOnFigS.GetOuterCityLayoutData(pArSeR.Uint64("prev_level"))
	if dAtA.PrevLevel == nil && pArSeR.Uint64("prev_level") != 0 {
		return errors.Errorf("%s 配置的关联字段[prev_level] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prev_level"), *pArSeR)
	}

	dAtA.ChangeTypeLevel = cOnFigS.GetOuterCityLayoutData(pArSeR.Uint64("change_type_id"))
	if dAtA.ChangeTypeLevel == nil && pArSeR.Uint64("change_type_id") != 0 {
		return errors.Errorf("%s 配置的关联字段[change_type_id] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("change_type_id"), *pArSeR)
	}

	dAtA.MilitaryBuilding = cOnFigS.GetOuterCityBuildingData(pArSeR.Uint64("military_building"))
	if dAtA.MilitaryBuilding == nil {
		return errors.Errorf("%s 配置的关联字段[military_building] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("military_building"), *pArSeR)
	}

	dAtA.EconomicBuilding = cOnFigS.GetOuterCityBuildingData(pArSeR.Uint64("economic_building"))
	if dAtA.EconomicBuilding == nil {
		return errors.Errorf("%s 配置的关联字段[economic_building] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("economic_building"), *pArSeR)
	}

	return nil
}

// start with ProsperityDamageBuffData ----------------------------------

func LoadProsperityDamageBuffData(gos *config.GameObjects) (map[uint64]*ProsperityDamageBuffData, map[*ProsperityDamageBuffData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ProsperityDamageBuffDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ProsperityDamageBuffData, len(lIsT))
	pArSeRmAp := make(map[*ProsperityDamageBuffData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrProsperityDamageBuffData) {
			continue
		}

		dAtA, err := NewProsperityDamageBuffData(fIlEnAmE, pArSeR)
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

func SetRelatedProsperityDamageBuffData(dAtAmAp map[*ProsperityDamageBuffData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ProsperityDamageBuffDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetProsperityDamageBuffDataKeyArray(datas []*ProsperityDamageBuffData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewProsperityDamageBuffData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ProsperityDamageBuffData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrProsperityDamageBuffData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ProsperityDamageBuffData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.MinPercent = pArSeR.Uint64("min_percent")
	dAtA.MaxPercent = pArSeR.Uint64("max_percent")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.BuffId = pArSeR.Uint64("buff_id")
	// skip field: BuffData

	return dAtA, nil
}

var vAlIdAtOrProsperityDamageBuffData = map[string]*config.Validator{

	"id":          config.ParseValidator("int>0", "", false, nil, nil),
	"min_percent": config.ParseValidator("uint", "", false, nil, nil),
	"max_percent": config.ParseValidator("int>0", "", false, nil, nil),
	"desc":        config.ParseValidator("string", "", false, nil, nil),
	"buff_id":     config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *ProsperityDamageBuffData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ProsperityDamageBuffData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ProsperityDamageBuffData) Encode() *shared_proto.ProsperityDamageBuffDataProto {
	out := &shared_proto.ProsperityDamageBuffDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.MinPercent = config.U64ToI32(dAtA.MinPercent)
	out.MaxPercent = config.U64ToI32(dAtA.MaxPercent)
	out.Desc = dAtA.Desc
	out.BuffId = config.U64ToI32(dAtA.BuffId)

	return out
}

func ArrayEncodeProsperityDamageBuffData(datas []*ProsperityDamageBuffData) []*shared_proto.ProsperityDamageBuffDataProto {

	out := make([]*shared_proto.ProsperityDamageBuffDataProto, 0, len(datas))
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

func (dAtA *ProsperityDamageBuffData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with SoldierLevelData ----------------------------------

func LoadSoldierLevelData(gos *config.GameObjects) (map[uint64]*SoldierLevelData, map[*SoldierLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.SoldierLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*SoldierLevelData, len(lIsT))
	pArSeRmAp := make(map[*SoldierLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrSoldierLevelData) {
			continue
		}

		dAtA, err := NewSoldierLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedSoldierLevelData(dAtAmAp map[*SoldierLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SoldierLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetSoldierLevelDataKeyArray(datas []*SoldierLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewSoldierLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SoldierLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSoldierLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SoldierLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Load = pArSeR.Uint64("load")
	// releated field: RecruitCost
	// releated field: WoundedCost
	// releated field: UpgradeCost
	dAtA.JunYingLevel = pArSeR.Uint64("jun_ying_level")
	// releated field: BaseStat
	// skip field: TotalStatSum
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Models = pArSeR.Uint64Array("models", "", false)

	return dAtA, nil
}

var vAlIdAtOrSoldierLevelData = map[string]*config.Validator{

	"level":          config.ParseValidator("int>0", "", false, nil, nil),
	"load":           config.ParseValidator("int>0", "", false, nil, nil),
	"recruit_cost":   config.ParseValidator("string", "", false, nil, nil),
	"wounded_cost":   config.ParseValidator("string", "", false, nil, nil),
	"upgrade_cost":   config.ParseValidator("string", "", false, nil, nil),
	"jun_ying_level": config.ParseValidator("int>0", "", false, nil, nil),
	"base_stat":      config.ParseValidator("int>0,count=5,notNil,duplicate", "", true, nil, nil),
	"desc":           config.ParseValidator("string", "", false, nil, nil),
	"models":         config.ParseValidator("string,count=5", "", true, nil, nil),
}

func (dAtA *SoldierLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SoldierLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SoldierLevelData) Encode() *shared_proto.SoldierLevelProto {
	out := &shared_proto.SoldierLevelProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Load = config.U64ToI32(dAtA.Load)
	if dAtA.RecruitCost != nil {
		out.RecruitCost = dAtA.RecruitCost.Encode()
	}
	if dAtA.WoundedCost != nil {
		out.WoundedCost = dAtA.WoundedCost.Encode()
	}
	if dAtA.UpgradeCost != nil {
		out.UpgradeCost = dAtA.UpgradeCost.Encode()
	}
	out.JunYingLevel = config.U64ToI32(dAtA.JunYingLevel)
	if dAtA.BaseStat != nil {
		out.BaseStat = data.ArrayEncodeSpriteStat(dAtA.BaseStat)
	}
	out.TotalStatSum = config.U64ToI32(dAtA.TotalStatSum)
	out.Desc = dAtA.Desc
	out.Models = config.U64a2I32a(dAtA.Models)

	return out
}

func ArrayEncodeSoldierLevelData(datas []*SoldierLevelData) []*shared_proto.SoldierLevelProto {

	out := make([]*shared_proto.SoldierLevelProto, 0, len(datas))
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

func (dAtA *SoldierLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.RecruitCost = cOnFigS.GetCost(pArSeR.Int("recruit_cost"))
	if dAtA.RecruitCost == nil {
		return errors.Errorf("%s 配置的关联字段[recruit_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("recruit_cost"), *pArSeR)
	}

	dAtA.WoundedCost = cOnFigS.GetCost(pArSeR.Int("wounded_cost"))
	if dAtA.WoundedCost == nil {
		return errors.Errorf("%s 配置的关联字段[wounded_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("wounded_cost"), *pArSeR)
	}

	dAtA.UpgradeCost = cOnFigS.GetCost(pArSeR.Int("upgrade_cost"))
	if dAtA.UpgradeCost == nil {
		return errors.Errorf("%s 配置的关联字段[upgrade_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("upgrade_cost"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("base_stat", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetSpriteStat(v)
		if obj != nil {
			dAtA.BaseStat = append(dAtA.BaseStat, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[base_stat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("base_stat"), *pArSeR)
		}
	}

	return nil
}

// start with TechnologyData ----------------------------------

func LoadTechnologyData(gos *config.GameObjects) (map[uint64]*TechnologyData, map[*TechnologyData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TechnologyDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TechnologyData, len(lIsT))
	pArSeRmAp := make(map[*TechnologyData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTechnologyData) {
			continue
		}

		dAtA, err := NewTechnologyData(fIlEnAmE, pArSeR)
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

func SetRelatedTechnologyData(dAtAmAp map[*TechnologyData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TechnologyDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTechnologyDataKeyArray(datas []*TechnologyData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTechnologyData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TechnologyData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTechnologyData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TechnologyData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Type = shared_proto.TechType(shared_proto.TechType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.TechType(i)
	}

	dAtA.Sequence = pArSeR.Uint64("sequence")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.IntIcon = 0
	if pArSeR.KeyExist("int_icon") {
		dAtA.IntIcon = pArSeR.Uint64("int_icon")
	}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.IsBigTech = false
	if pArSeR.KeyExist("is_big_tech") {
		dAtA.IsBigTech = pArSeR.Bool("is_big_tech")
	}

	// releated field: Effect
	// releated field: RequireBuildingIds
	// releated field: RequireTechIds
	// releated field: Cost
	dAtA.WorkTime, err = config.ParseDuration(pArSeR.String("work_time"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[work_time] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("work_time"), dAtA)
	}

	// skip field: NextLevel
	// skip field: MaxLevel

	// calculate fields
	dAtA.Group = uint64(dAtA.Type)*10000 + dAtA.Sequence

	return dAtA, nil
}

var vAlIdAtOrTechnologyData = map[string]*config.Validator{

	"id":                   config.ParseValidator("int>0", "", false, nil, nil),
	"name":                 config.ParseValidator("string", "", false, nil, nil),
	"desc":                 config.ParseValidator("string", "", false, nil, nil),
	"type":                 config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.TechType_value, 0), nil),
	"sequence":             config.ParseValidator("int>0", "", false, nil, nil),
	"icon":                 config.ParseValidator("string", "", false, nil, nil),
	"int_icon":             config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"level":                config.ParseValidator("int>0", "", false, nil, nil),
	"is_big_tech":          config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"effect":               config.ParseValidator("string", "", false, nil, nil),
	"require_building_ids": config.ParseValidator("string", "", true, nil, nil),
	"require_tech_ids":     config.ParseValidator("string", "", true, nil, nil),
	"cost":                 config.ParseValidator("string", "", false, nil, nil),
	"work_time":            config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *TechnologyData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TechnologyData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TechnologyData) Encode() *shared_proto.TechnologyDataProto {
	out := &shared_proto.TechnologyDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.Type = dAtA.Type
	out.Sequence = config.U64ToI32(dAtA.Sequence)
	out.Icon = dAtA.Icon
	out.IntIcon = config.U64ToI32(dAtA.IntIcon)
	out.Group = config.U64ToI32(dAtA.Group)
	out.Level = config.U64ToI32(dAtA.Level)
	out.IsBigTech = dAtA.IsBigTech
	if dAtA.Effect != nil {
		out.Effect = dAtA.Effect.Encode()
	}
	if dAtA.RequireBuildingIds != nil {
		out.RequireBuildingIds = config.U64a2I32a(GetBuildingDataKeyArray(dAtA.RequireBuildingIds))
	}
	if dAtA.RequireTechIds != nil {
		out.RequireTechIds = config.U64a2I32a(GetTechnologyDataKeyArray(dAtA.RequireTechIds))
	}
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	out.WorkTime = config.Duration2I32Seconds(dAtA.WorkTime)
	if dAtA.NextLevel != nil {
		out.NextLevelId = config.U64ToI32(dAtA.NextLevel.Id)
	}
	out.MaxLevel = config.U64ToI32(dAtA.MaxLevel)

	return out
}

func ArrayEncodeTechnologyData(datas []*TechnologyData) []*shared_proto.TechnologyDataProto {

	out := make([]*shared_proto.TechnologyDataProto, 0, len(datas))
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

func (dAtA *TechnologyData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Effect = cOnFigS.GetBuildingEffectData(pArSeR.Int("effect"))
	if dAtA.Effect == nil {
		return errors.Errorf("%s 配置的关联字段[effect] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("effect"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("require_building_ids", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetBuildingData(v)
		if obj != nil {
			dAtA.RequireBuildingIds = append(dAtA.RequireBuildingIds, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[require_building_ids] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("require_building_ids"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("require_tech_ids", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetTechnologyData(v)
		if obj != nil {
			dAtA.RequireTechIds = append(dAtA.RequireTechIds, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[require_tech_ids] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("require_tech_ids"), *pArSeR)
		}
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	return nil
}

// start with TieJiangPuLevelData ----------------------------------

func LoadTieJiangPuLevelData(gos *config.GameObjects) (map[uint64]*TieJiangPuLevelData, map[*TieJiangPuLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TieJiangPuLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TieJiangPuLevelData, len(lIsT))
	pArSeRmAp := make(map[*TieJiangPuLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTieJiangPuLevelData) {
			continue
		}

		dAtA, err := NewTieJiangPuLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedTieJiangPuLevelData(dAtAmAp map[*TieJiangPuLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TieJiangPuLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTieJiangPuLevelDataKeyArray(datas []*TieJiangPuLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewTieJiangPuLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TieJiangPuLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTieJiangPuLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TieJiangPuLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.MaxForgingTimes = pArSeR.Uint64("max_forging_times")
	dAtA.RecoveryForgingDuration, err = config.ParseDuration(pArSeR.String("recovery_forging_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[recovery_forging_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("recovery_forging_duration"), dAtA)
	}

	dAtA.CanForgingEquipPos = pArSeR.Uint64Array("can_forging_equip_pos", "", false)
	// releated field: CanForgingEquip
	// skip field: LockedCanForgingEquipPos
	// skip field: LockedCanForgingEquip
	// skip field: LockedEquipNeedLevel
	dAtA.CanOneKeyForging = pArSeR.Bool("can_one_key_forging")

	return dAtA, nil
}

var vAlIdAtOrTieJiangPuLevelData = map[string]*config.Validator{

	"level":                     config.ParseValidator("int>0", "", false, nil, nil),
	"max_forging_times":         config.ParseValidator("int>0", "", false, nil, nil),
	"recovery_forging_duration": config.ParseValidator("string", "", false, nil, nil),
	"can_forging_equip_pos":     config.ParseValidator("int>0,notAllNil", "", true, nil, nil),
	"can_forging_equip_id":      config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
	"can_one_key_forging":       config.ParseValidator("bool", "", false, nil, nil),
}

func (dAtA *TieJiangPuLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TieJiangPuLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TieJiangPuLevelData) Encode() *shared_proto.TieJiangPuLevelProto {
	out := &shared_proto.TieJiangPuLevelProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.MaxForgingTimes = config.U64ToI32(dAtA.MaxForgingTimes)
	out.RecoveryForgingDuration = config.Duration2I32Seconds(dAtA.RecoveryForgingDuration)
	out.CanForgingEquipPos = config.U64a2I32a(dAtA.CanForgingEquipPos)
	if dAtA.CanForgingEquip != nil {
		out.CanForgingEquip = config.U64a2I32a(goods.GetEquipmentDataKeyArray(dAtA.CanForgingEquip))
	}
	out.LockedCanForgingEquipPos = config.U64a2I32a(dAtA.LockedCanForgingEquipPos)
	if dAtA.LockedCanForgingEquip != nil {
		out.LockedCanForgingEquip = config.U64a2I32a(goods.GetEquipmentDataKeyArray(dAtA.LockedCanForgingEquip))
	}
	out.LockedEquipNeedLevel = config.U64a2I32a(dAtA.LockedEquipNeedLevel)
	out.CanOneKeyForging = dAtA.CanOneKeyForging

	return out
}

func ArrayEncodeTieJiangPuLevelData(datas []*TieJiangPuLevelData) []*shared_proto.TieJiangPuLevelProto {

	out := make([]*shared_proto.TieJiangPuLevelProto, 0, len(datas))
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

func (dAtA *TieJiangPuLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("can_forging_equip_id", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetEquipmentData(v)
		if obj != nil {
			dAtA.CanForgingEquip = append(dAtA.CanForgingEquip, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[can_forging_equip_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("can_forging_equip_id"), *pArSeR)
		}
	}

	return nil
}

// start with WorkshopDuration ----------------------------------

func LoadWorkshopDuration(gos *config.GameObjects) (map[uint64]*WorkshopDuration, map[*WorkshopDuration]*config.ObjectParser, error) {
	fIlEnAmE := confpath.WorkshopDurationPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*WorkshopDuration, len(lIsT))
	pArSeRmAp := make(map[*WorkshopDuration]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrWorkshopDuration) {
			continue
		}

		dAtA, err := NewWorkshopDuration(fIlEnAmE, pArSeR)
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

func SetRelatedWorkshopDuration(dAtAmAp map[*WorkshopDuration]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.WorkshopDurationPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetWorkshopDurationKeyArray(datas []*WorkshopDuration) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewWorkshopDuration(fIlEnAmE string, pArSeR *config.ObjectParser) (*WorkshopDuration, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrWorkshopDuration)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &WorkshopDuration{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Duration, err = config.ParseDuration(pArSeR.String("duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("duration"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrWorkshopDuration = map[string]*config.Validator{

	"id":       config.ParseValidator("int>0", "", false, nil, nil),
	"duration": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *WorkshopDuration) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with WorkshopLevelData ----------------------------------

func LoadWorkshopLevelData(gos *config.GameObjects) (map[uint64]*WorkshopLevelData, map[*WorkshopLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.WorkshopLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*WorkshopLevelData, len(lIsT))
	pArSeRmAp := make(map[*WorkshopLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrWorkshopLevelData) {
			continue
		}

		dAtA, err := NewWorkshopLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedWorkshopLevelData(dAtAmAp map[*WorkshopLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.WorkshopLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetWorkshopLevelDataKeyArray(datas []*WorkshopLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewWorkshopLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*WorkshopLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrWorkshopLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &WorkshopLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	// releated field: Group

	return dAtA, nil
}

var vAlIdAtOrWorkshopLevelData = map[string]*config.Validator{

	"level": config.ParseValidator("int>0", "", false, nil, nil),
	"group": config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
}

func (dAtA *WorkshopLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("group", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetPlunderGroup(v)
		if obj != nil {
			dAtA.Group = append(dAtA.Group, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[group] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("group"), *pArSeR)
		}
	}

	return nil
}

// start with WorkshopRefreshCost ----------------------------------

func LoadWorkshopRefreshCost(gos *config.GameObjects) (map[uint64]*WorkshopRefreshCost, map[*WorkshopRefreshCost]*config.ObjectParser, error) {
	fIlEnAmE := confpath.WorkshopRefreshCostPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*WorkshopRefreshCost, len(lIsT))
	pArSeRmAp := make(map[*WorkshopRefreshCost]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrWorkshopRefreshCost) {
			continue
		}

		dAtA, err := NewWorkshopRefreshCost(fIlEnAmE, pArSeR)
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

func SetRelatedWorkshopRefreshCost(dAtAmAp map[*WorkshopRefreshCost]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.WorkshopRefreshCostPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetWorkshopRefreshCostKeyArray(datas []*WorkshopRefreshCost) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewWorkshopRefreshCost(fIlEnAmE string, pArSeR *config.ObjectParser) (*WorkshopRefreshCost, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrWorkshopRefreshCost)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &WorkshopRefreshCost{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Cost

	return dAtA, nil
}

var vAlIdAtOrWorkshopRefreshCost = map[string]*config.Validator{

	"id":   config.ParseValidator("int>0", "", false, nil, nil),
	"cost": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *WorkshopRefreshCost) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *WorkshopRefreshCost) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *WorkshopRefreshCost) Encode() *shared_proto.WorkshopRefreshCostProto {
	out := &shared_proto.WorkshopRefreshCostProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}

	return out
}

func ArrayEncodeWorkshopRefreshCost(datas []*WorkshopRefreshCost) []*shared_proto.WorkshopRefreshCostProto {

	out := make([]*shared_proto.WorkshopRefreshCostProto, 0, len(datas))
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

func (dAtA *WorkshopRefreshCost) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetBaseLevelData(uint64) *BaseLevelData
	GetBuildingData(uint64) *BuildingData
	GetBuildingEffectData(int) *sub.BuildingEffectData
	GetCityEventData(uint64) *CityEventData
	GetCombineCost(int) *CombineCost
	GetCost(int) *resdata.Cost
	GetCountdownPrizeDescData(uint64) *CountdownPrizeDescData
	GetEquipmentData(uint64) *goods.EquipmentData
	GetIcon(string) *icon.Icon
	GetOuterCityBuildingData(uint64) *OuterCityBuildingData
	GetOuterCityBuildingDescData(uint64) *OuterCityBuildingDescData
	GetOuterCityData(uint64) *OuterCityData
	GetOuterCityLayoutData(uint64) *OuterCityLayoutData
	GetPlunder(uint64) *resdata.Plunder
	GetPlunderGroup(uint64) *resdata.PlunderGroup
	GetPrize(int) *resdata.Prize
	GetSpriteStat(uint64) *data.SpriteStat
	GetTechnologyData(uint64) *TechnologyData
}
