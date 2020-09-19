// AUTO_GEN, DONT MODIFY!!!
package goods

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
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

// start with EquipmentData ----------------------------------

func LoadEquipmentData(gos *config.GameObjects) (map[uint64]*EquipmentData, map[*EquipmentData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.EquipmentDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*EquipmentData, len(lIsT))
	pArSeRmAp := make(map[*EquipmentData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrEquipmentData) {
			continue
		}

		dAtA, err := NewEquipmentData(fIlEnAmE, pArSeR)
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

func SetRelatedEquipmentData(dAtAmAp map[*EquipmentData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EquipmentDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetEquipmentDataKeyArray(datas []*EquipmentData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewEquipmentData(fIlEnAmE string, pArSeR *config.ObjectParser) (*EquipmentData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEquipmentData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EquipmentData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	// releated field: Icon
	dAtA.Type = shared_proto.EquipmentType(shared_proto.EquipmentType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.EquipmentType(i)
	}

	// releated field: Quality
	// releated field: BaseStat
	// skip field: BaseStatProto

	return dAtA, nil
}

var vAlIdAtOrEquipmentData = map[string]*config.Validator{

	"id":        config.ParseValidator("int>0", "", false, nil, nil),
	"name":      config.ParseValidator("string", "", false, nil, nil),
	"desc":      config.ParseValidator("string", "", false, nil, nil),
	"icon":      config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"type":      config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.EquipmentType_value, 0), nil),
	"quality":   config.ParseValidator("string", "", false, nil, nil),
	"base_stat": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *EquipmentData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *EquipmentData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *EquipmentData) Encode() *shared_proto.EquipmentDataProto {
	out := &shared_proto.EquipmentDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.Type = dAtA.Type
	if dAtA.Quality != nil {
		out.Quality = config.U64ToI32(dAtA.Quality.Id)
	}
	if dAtA.BaseStat != nil {
		out.BaseStat = dAtA.BaseStat.Encode()
	}

	return out
}

func ArrayEncodeEquipmentData(datas []*EquipmentData) []*shared_proto.EquipmentDataProto {

	out := make([]*shared_proto.EquipmentDataProto, 0, len(datas))
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

func (dAtA *EquipmentData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.Quality = cOnFigS.GetEquipmentQualityData(pArSeR.Uint64("quality"))
	if dAtA.Quality == nil {
		return errors.Errorf("%s 配置的关联字段[quality] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("quality"), *pArSeR)
	}

	dAtA.BaseStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("base_stat"))
	if dAtA.BaseStat == nil {
		return errors.Errorf("%s 配置的关联字段[base_stat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("base_stat"), *pArSeR)
	}

	return nil
}

// start with EquipmentLevelData ----------------------------------

func LoadEquipmentLevelData(gos *config.GameObjects) (map[uint64]*EquipmentLevelData, map[*EquipmentLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.EquipmentLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*EquipmentLevelData, len(lIsT))
	pArSeRmAp := make(map[*EquipmentLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrEquipmentLevelData) {
			continue
		}

		dAtA, err := NewEquipmentLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedEquipmentLevelData(dAtAmAp map[*EquipmentLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EquipmentLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetEquipmentLevelDataKeyArray(datas []*EquipmentLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewEquipmentLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*EquipmentLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEquipmentLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EquipmentLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.UpgradeCostCoef = pArSeR.Float64("upgrade_cost_coef")

	return dAtA, nil
}

var vAlIdAtOrEquipmentLevelData = map[string]*config.Validator{

	"level":             config.ParseValidator("int>0", "", false, nil, nil),
	"upgrade_cost_coef": config.ParseValidator("float64>0", "", false, nil, nil),
}

func (dAtA *EquipmentLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with EquipmentQualityData ----------------------------------

func LoadEquipmentQualityData(gos *config.GameObjects) (map[uint64]*EquipmentQualityData, map[*EquipmentQualityData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.EquipmentQualityDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*EquipmentQualityData, len(lIsT))
	pArSeRmAp := make(map[*EquipmentQualityData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrEquipmentQualityData) {
			continue
		}

		dAtA, err := NewEquipmentQualityData(fIlEnAmE, pArSeR)
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

func SetRelatedEquipmentQualityData(dAtAmAp map[*EquipmentQualityData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EquipmentQualityDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetEquipmentQualityDataKeyArray(datas []*EquipmentQualityData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewEquipmentQualityData(fIlEnAmE string, pArSeR *config.ObjectParser) (*EquipmentQualityData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEquipmentQualityData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EquipmentQualityData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Level = 1
	if pArSeR.KeyExist("level") {
		dAtA.Level = pArSeR.Uint64("level")
	}

	// releated field: GoodsQuality
	// releated field: FirstLevelRefined
	dAtA.RefinedLevelLimit = pArSeR.Uint64("refined_level_limit")
	dAtA.SmeltBackCount = pArSeR.Uint64("smelt_back_count")
	dAtA.LevelCostCoef = pArSeR.Float64("level_cost_coef")
	dAtA.LevelStatCoef = pArSeR.Float64("level_stat_coef")
	// skip field: LevelDatas

	return dAtA, nil
}

var vAlIdAtOrEquipmentQualityData = map[string]*config.Validator{

	"id":                  config.ParseValidator("int>0", "", false, nil, nil),
	"level":               config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"goods_quality":       config.ParseValidator("string", "", false, nil, nil),
	"first_level_refined": config.ParseValidator("string", "", false, nil, []string{"1"}),
	"refined_level_limit": config.ParseValidator("int>0", "", false, nil, nil),
	"smelt_back_count":    config.ParseValidator("int>0", "", false, nil, nil),
	"level_cost_coef":     config.ParseValidator("float64>0", "", false, nil, nil),
	"level_stat_coef":     config.ParseValidator("float64>0", "", false, nil, nil),
}

func (dAtA *EquipmentQualityData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *EquipmentQualityData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *EquipmentQualityData) Encode() *shared_proto.EquipmentQualityProto {
	out := &shared_proto.EquipmentQualityProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Level = config.U64ToI32(dAtA.Level)
	if dAtA.GoodsQuality != nil {
		out.GoodsQuality = config.U64ToI32(dAtA.GoodsQuality.Level)
	}
	out.RefinedLevelLimit = config.U64ToI32(dAtA.RefinedLevelLimit)
	out.SmeltBackCount = config.U64ToI32(dAtA.SmeltBackCount)

	return out
}

func ArrayEncodeEquipmentQualityData(datas []*EquipmentQualityData) []*shared_proto.EquipmentQualityProto {

	out := make([]*shared_proto.EquipmentQualityProto, 0, len(datas))
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

func (dAtA *EquipmentQualityData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.GoodsQuality = cOnFigS.GetGoodsQuality(pArSeR.Uint64("goods_quality"))
	if dAtA.GoodsQuality == nil {
		return errors.Errorf("%s 配置的关联字段[goods_quality] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("goods_quality"), *pArSeR)
	}

	if pArSeR.KeyExist("first_level_refined") {
		dAtA.FirstLevelRefined = cOnFigS.GetEquipmentRefinedData(pArSeR.Uint64("first_level_refined"))
	} else {
		dAtA.FirstLevelRefined = cOnFigS.GetEquipmentRefinedData(1)
	}
	if dAtA.FirstLevelRefined == nil {
		return errors.Errorf("%s 配置的关联字段[first_level_refined] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_level_refined"), *pArSeR)
	}

	return nil
}

// start with EquipmentRefinedData ----------------------------------

func LoadEquipmentRefinedData(gos *config.GameObjects) (map[uint64]*EquipmentRefinedData, map[*EquipmentRefinedData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.EquipmentRefinedDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*EquipmentRefinedData, len(lIsT))
	pArSeRmAp := make(map[*EquipmentRefinedData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrEquipmentRefinedData) {
			continue
		}

		dAtA, err := NewEquipmentRefinedData(fIlEnAmE, pArSeR)
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

func SetRelatedEquipmentRefinedData(dAtAmAp map[*EquipmentRefinedData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EquipmentRefinedDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetEquipmentRefinedDataKeyArray(datas []*EquipmentRefinedData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewEquipmentRefinedData(fIlEnAmE string, pArSeR *config.ObjectParser) (*EquipmentRefinedData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEquipmentRefinedData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EquipmentRefinedData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.StatCoef = pArSeR.Float64("stat_coef")
	// skip field: CurrentStatCoef
	dAtA.CostCount = pArSeR.Uint64("cost_count")
	dAtA.HeroLevelLimit = pArSeR.Uint64("hero_level_limit")
	// skip field: TotalCostCount

	// calculate fields
	dAtA.Star = dAtA.Level / 2

	return dAtA, nil
}

var vAlIdAtOrEquipmentRefinedData = map[string]*config.Validator{

	"level":            config.ParseValidator("int>0", "", false, nil, nil),
	"stat_coef":        config.ParseValidator("float64>0", "", false, nil, nil),
	"cost_count":       config.ParseValidator("int>0", "", false, nil, nil),
	"hero_level_limit": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *EquipmentRefinedData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *EquipmentRefinedData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *EquipmentRefinedData) Encode() *shared_proto.EquipmentRefinedProto {
	out := &shared_proto.EquipmentRefinedProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.CostCount = config.U64ToI32(dAtA.CostCount)
	out.HeroLevelLimit = config.U64ToI32(dAtA.HeroLevelLimit)
	out.TotalCostCount = config.U64ToI32(dAtA.TotalCostCount)

	return out
}

func ArrayEncodeEquipmentRefinedData(datas []*EquipmentRefinedData) []*shared_proto.EquipmentRefinedProto {

	out := make([]*shared_proto.EquipmentRefinedProto, 0, len(datas))
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

func (dAtA *EquipmentRefinedData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with EquipmentTaozConfig ----------------------------------

func LoadEquipmentTaozConfig(gos *config.GameObjects) (*EquipmentTaozConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.EquipmentTaozConfigPath
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

	dAtA, err := NewEquipmentTaozConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedEquipmentTaozConfig(gos *config.GameObjects, dAtA *EquipmentTaozConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EquipmentTaozConfigPath
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

func NewEquipmentTaozConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*EquipmentTaozConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEquipmentTaozConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EquipmentTaozConfig{}

	return dAtA, nil
}

var vAlIdAtOrEquipmentTaozConfig = map[string]*config.Validator{}

func (dAtA *EquipmentTaozConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with EquipmentTaozData ----------------------------------

func LoadEquipmentTaozData(gos *config.GameObjects) (map[uint64]*EquipmentTaozData, map[*EquipmentTaozData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.EquipmentTaozDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*EquipmentTaozData, len(lIsT))
	pArSeRmAp := make(map[*EquipmentTaozData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrEquipmentTaozData) {
			continue
		}

		dAtA, err := NewEquipmentTaozData(fIlEnAmE, pArSeR)
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

func SetRelatedEquipmentTaozData(dAtAmAp map[*EquipmentTaozData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EquipmentTaozDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetEquipmentTaozDataKeyArray(datas []*EquipmentTaozData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewEquipmentTaozData(fIlEnAmE string, pArSeR *config.ObjectParser) (*EquipmentTaozData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEquipmentTaozData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EquipmentTaozData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Count = pArSeR.Uint64("count")
	dAtA.RefinedLevel = pArSeR.Uint64("refined_level")
	dAtA.Morale = pArSeR.Uint64("morale")
	// releated field: SpriteStat

	// calculate fields
	dAtA.Star = dAtA.Level / 2

	return dAtA, nil
}

var vAlIdAtOrEquipmentTaozData = map[string]*config.Validator{

	"level":         config.ParseValidator("int>0", "", false, nil, nil),
	"count":         config.ParseValidator("int>0", "", false, nil, nil),
	"refined_level": config.ParseValidator("int>0", "", false, nil, nil),
	"morale":        config.ParseValidator("int>0", "", false, nil, nil),
	"sprite_stat":   config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *EquipmentTaozData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *EquipmentTaozData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *EquipmentTaozData) Encode() *shared_proto.EquipmentTaozProto {
	out := &shared_proto.EquipmentTaozProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Count = config.U64ToI32(dAtA.Count)
	out.RefinedLevel = config.U64ToI32(dAtA.RefinedLevel)
	out.Morale = config.U64ToI32(dAtA.Morale)
	if dAtA.SpriteStat != nil {
		out.SpriteStat = dAtA.SpriteStat.Encode()
	}

	return out
}

func ArrayEncodeEquipmentTaozData(datas []*EquipmentTaozData) []*shared_proto.EquipmentTaozProto {

	out := make([]*shared_proto.EquipmentTaozProto, 0, len(datas))
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

func (dAtA *EquipmentTaozData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.SpriteStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("sprite_stat"))
	if dAtA.SpriteStat == nil {
		return errors.Errorf("%s 配置的关联字段[sprite_stat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("sprite_stat"), *pArSeR)
	}

	return nil
}

// start with GemData ----------------------------------

func LoadGemData(gos *config.GameObjects) (map[uint64]*GemData, map[*GemData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GemDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GemData, len(lIsT))
	pArSeRmAp := make(map[*GemData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGemData) {
			continue
		}

		dAtA, err := NewGemData(fIlEnAmE, pArSeR)
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

func SetRelatedGemData(dAtAmAp map[*GemData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GemDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGemDataKeyArray(datas []*GemData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGemData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GemData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGemData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GemData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	// releated field: Icon
	dAtA.Quality = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("quality"))])
	if i, err := strconv.ParseInt(pArSeR.String("quality"), 10, 32); err == nil {
		dAtA.Quality = shared_proto.Quality(i)
	}

	// releated field: BaseStat
	dAtA.GemType = pArSeR.Uint64("gem_type")
	dAtA.Level = pArSeR.Uint64("level")
	dAtA.YuanbaoPrice = 0
	if pArSeR.KeyExist("yuanbao_price") {
		dAtA.YuanbaoPrice = pArSeR.Uint64("yuanbao_price")
	}

	dAtA.UpgradeNeedCount = pArSeR.Uint64("upgrade_need_count")
	// skip field: UpgradeToThisNeedFirstLevelCount
	// skip field: PrevLevel
	// skip field: NextLevel

	return dAtA, nil
}

var vAlIdAtOrGemData = map[string]*config.Validator{

	"id":                 config.ParseValidator("int>0", "", false, nil, nil),
	"name":               config.ParseValidator("string", "", false, nil, nil),
	"desc":               config.ParseValidator("string", "", false, nil, nil),
	"icon":               config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"quality":            config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), nil),
	"base_stat":          config.ParseValidator("string", "", false, nil, nil),
	"gem_type":           config.ParseValidator("int>0", "", false, nil, nil),
	"level":              config.ParseValidator("int>0", "", false, nil, nil),
	"yuanbao_price":      config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"upgrade_need_count": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *GemData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GemData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GemData) Encode() *shared_proto.GemDataProto {
	out := &shared_proto.GemDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.Quality = dAtA.Quality
	if dAtA.BaseStat != nil {
		out.BaseStat = dAtA.BaseStat.Encode()
	}
	out.GemType = config.U64ToI32(dAtA.GemType)
	out.Level = config.U64ToI32(dAtA.Level)
	out.YuanbaoPrice = config.U64ToI32(dAtA.YuanbaoPrice)
	out.UpgradeNeedCount = config.U64ToI32(dAtA.UpgradeNeedCount)
	if dAtA.PrevLevel != nil {
		out.PrevLevel = config.U64ToI32(dAtA.PrevLevel.Id)
	}
	if dAtA.NextLevel != nil {
		out.NextLevel = config.U64ToI32(dAtA.NextLevel.Id)
	}

	return out
}

func ArrayEncodeGemData(datas []*GemData) []*shared_proto.GemDataProto {

	out := make([]*shared_proto.GemDataProto, 0, len(datas))
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

func (dAtA *GemData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.BaseStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("base_stat"))
	if dAtA.BaseStat == nil {
		return errors.Errorf("%s 配置的关联字段[base_stat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("base_stat"), *pArSeR)
	}

	return nil
}

// start with GemDatas ----------------------------------

func LoadGemDatas(gos *config.GameObjects) (*GemDatas, *config.ObjectParser, error) {
	fIlEnAmE := confpath.GemDatasPath
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

	dAtA, err := NewGemDatas(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedGemDatas(gos *config.GameObjects, dAtA *GemDatas, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GemDatasPath
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

func NewGemDatas(fIlEnAmE string, pArSeR *config.ObjectParser) (*GemDatas, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGemDatas)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GemDatas{}

	return dAtA, nil
}

var vAlIdAtOrGemDatas = map[string]*config.Validator{}

func (dAtA *GemDatas) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GoodsCheck ----------------------------------

func LoadGoodsCheck(gos *config.GameObjects) (*GoodsCheck, *config.ObjectParser, error) {
	fIlEnAmE := confpath.GoodsCheckPath
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

	dAtA, err := NewGoodsCheck(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedGoodsCheck(gos *config.GameObjects, dAtA *GoodsCheck, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GoodsCheckPath
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

func NewGoodsCheck(fIlEnAmE string, pArSeR *config.ObjectParser) (*GoodsCheck, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGoodsCheck)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GoodsCheck{}

	return dAtA, nil
}

var vAlIdAtOrGoodsCheck = map[string]*config.Validator{}

func (dAtA *GoodsCheck) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GoodsData ----------------------------------

func LoadGoodsData(gos *config.GameObjects) (map[uint64]*GoodsData, map[*GoodsData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GoodsDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GoodsData, len(lIsT))
	pArSeRmAp := make(map[*GoodsData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGoodsData) {
			continue
		}

		dAtA, err := NewGoodsData(fIlEnAmE, pArSeR)
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

func SetRelatedGoodsData(dAtAmAp map[*GoodsData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GoodsDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGoodsDataKeyArray(datas []*GoodsData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGoodsData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GoodsData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGoodsData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GoodsData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	// releated field: Icon
	dAtA.ObtainWays = pArSeR.StringArray("obtain_ways", "", false)
	// skip field: Quality
	// releated field: GoodsQuality
	dAtA.YuanbaoPrice = 0
	if pArSeR.KeyExist("yuanbao_price") {
		dAtA.YuanbaoPrice = pArSeR.Uint64("yuanbao_price")
	}

	dAtA.DianquanPrice = 0
	if pArSeR.KeyExist("dianquan_price") {
		dAtA.DianquanPrice = pArSeR.Uint64("dianquan_price")
	}

	dAtA.YinliangPrice = 0
	if pArSeR.KeyExist("yinliang_price") {
		dAtA.YinliangPrice = pArSeR.Uint64("yinliang_price")
	}

	if pArSeR.KeyExist("cd") {
		dAtA.Cd, err = config.ParseDuration(pArSeR.String("cd"))
	} else {
		dAtA.Cd, err = config.ParseDuration("0s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("cd"), dAtA)
	}

	dAtA.Captain = 0
	if pArSeR.KeyExist("captain") {
		dAtA.Captain = pArSeR.Uint64("captain")
	}

	// skip field: EffectType
	dAtA.GoodsEffect, err = NewGoodsEffect(fIlEnAmE, pArSeR)
	if err != nil {
		return nil, err
	}
	dAtA.SpecType = shared_proto.GoodsSpecType(shared_proto.GoodsSpecType_value[strings.ToUpper(pArSeR.String("spec_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("spec_type"), 10, 32); err == nil {
		dAtA.SpecType = shared_proto.GoodsSpecType(i)
	}

	dAtA.HebiType = shared_proto.HebiType(shared_proto.HebiType_value[strings.ToUpper(pArSeR.String("hebi_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("hebi_type"), 10, 32); err == nil {
		dAtA.HebiType = shared_proto.HebiType(i)
	}

	dAtA.HebiSubType = shared_proto.HebiSubType(shared_proto.HebiSubType_value[strings.ToUpper(pArSeR.String("hebi_sub_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("hebi_sub_type"), 10, 32); err == nil {
		dAtA.HebiSubType = shared_proto.HebiSubType(i)
	}

	dAtA.PartnerHebiGoods = 0
	if pArSeR.KeyExist("partner_hebi_goods") {
		dAtA.PartnerHebiGoods = pArSeR.Uint64("partner_hebi_goods")
	}

	// skip field: PartnerHebiGoodsData

	return dAtA, nil
}

var vAlIdAtOrGoodsData = map[string]*config.Validator{

	"id":                 config.ParseValidator("int>0", "", false, nil, nil),
	"name":               config.ParseValidator("string", "", false, nil, nil),
	"desc":               config.ParseValidator("string", "", false, nil, nil),
	"icon":               config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"obtain_ways":        config.ParseValidator("string", "", true, nil, nil),
	"goods_quality":      config.ParseValidator("string", "", false, nil, nil),
	"yuanbao_price":      config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"dianquan_price":     config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"yinliang_price":     config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"cd":                 config.ParseValidator("string", "", false, nil, []string{"0s"}),
	"captain":            config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"spec_type":          config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.GoodsSpecType_value), nil),
	"hebi_type":          config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.HebiType_value), nil),
	"hebi_sub_type":      config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.HebiSubType_value), nil),
	"partner_hebi_goods": config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *GoodsData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GoodsData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GoodsData) Encode() *shared_proto.GoodsDataProto {
	out := &shared_proto.GoodsDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.ObtainWays = dAtA.ObtainWays
	out.Quality = dAtA.Quality
	if dAtA.GoodsQuality != nil {
		out.GoodsQuality = config.U64ToI32(dAtA.GoodsQuality.Level)
	}
	out.YuanbaoPrice = config.U64ToI32(dAtA.YuanbaoPrice)
	out.DianquanPrice = config.U64ToI32(dAtA.DianquanPrice)
	out.YinliangPrice = config.U64ToI32(dAtA.YinliangPrice)
	out.Cd = config.Duration2I32Seconds(dAtA.Cd)
	out.Captain = config.U64ToI32(dAtA.Captain)
	out.EffectType = dAtA.EffectType
	if dAtA.GoodsEffect != nil {
		out.GoodsEffect = dAtA.GoodsEffect.Encode()
	}
	out.SpecType = dAtA.SpecType
	out.HebiType = dAtA.HebiType
	out.HebiSubType = dAtA.HebiSubType
	out.PartnerHebiGoods = config.U64ToI32(dAtA.PartnerHebiGoods)

	return out
}

func ArrayEncodeGoodsData(datas []*GoodsData) []*shared_proto.GoodsDataProto {

	out := make([]*shared_proto.GoodsDataProto, 0, len(datas))
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

func (dAtA *GoodsData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.GoodsQuality = cOnFigS.GetGoodsQuality(pArSeR.Uint64("goods_quality"))
	if dAtA.GoodsQuality == nil {
		return errors.Errorf("%s 配置的关联字段[goods_quality] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("goods_quality"), *pArSeR)
	}

	if err := dAtA.GoodsEffect.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS0); err != nil {
		return err
	}

	return nil
}

// start with GoodsEffect ----------------------------------

func NewGoodsEffect(fIlEnAmE string, pArSeR *config.ObjectParser) (*GoodsEffect, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGoodsEffect)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GoodsEffect{}

	dAtA.Id = pArSeR.Int("id")
	// skip field: EffectType
	dAtA.Gold = pArSeR.Uint64("gold")
	dAtA.Food = pArSeR.Uint64("food")
	dAtA.Wood = pArSeR.Uint64("wood")
	dAtA.Stone = pArSeR.Uint64("stone")
	dAtA.BuildingCdr = pArSeR.Bool("building_cdr")
	dAtA.TechCdr = pArSeR.Bool("tech_cdr")
	dAtA.WorkshopCdr = pArSeR.Bool("workshop_cdr")
	dAtA.Cdr, err = config.ParseDuration(pArSeR.String("cdr"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[cdr] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("cdr"), dAtA)
	}

	// skip field: MoveBaseType
	dAtA.MoveBase = pArSeR.Bool("move_base")
	dAtA.RandomPos = pArSeR.Bool("random_pos")
	dAtA.GuildMoveBase = pArSeR.Bool("guild_move_base")
	dAtA.MonsterMoveSubLevel = pArSeR.Uint64("monster_move_sub_level")
	dAtA.ExpType = shared_proto.GoodsExpEffectType(shared_proto.GoodsExpEffectType_value[strings.ToUpper(pArSeR.String("exp_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("exp_type"), 10, 32); err == nil {
		dAtA.ExpType = shared_proto.GoodsExpEffectType(i)
	}

	dAtA.Exp = pArSeR.Uint64("exp")
	dAtA.MianDuration, err = config.ParseDuration(pArSeR.String("mian_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[mian_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("mian_duration"), dAtA)
	}

	dAtA.TroopSpeedUpRate = pArSeR.Float64("troop_speed_up_rate")
	dAtA.TrainDuration, err = config.ParseDuration(pArSeR.String("train_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[train_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("train_duration"), dAtA)
	}

	dAtA.PartsCombineCount = pArSeR.Uint64("parts_combine_count")
	dAtA.PartsPlunderPrizeId = pArSeR.Uint64Array("parts_plunder_prize_id", "", false)
	// skip field: PartsShowPrize
	dAtA.PartsShowType = 0
	if pArSeR.KeyExist("parts_show_type") {
		dAtA.PartsShowType = pArSeR.Uint64("parts_show_type")
	}

	dAtA.AddMultiLevelNpcTimes = pArSeR.Bool("add_multi_level_npc_times")
	dAtA.AddInvaseHeroTimes = pArSeR.Bool("add_invase_hero_times")
	dAtA.DurationType = shared_proto.GoodsDurationEffectType(shared_proto.GoodsDurationEffectType_value[strings.ToUpper(pArSeR.String("duration_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("duration_type"), 10, 32); err == nil {
		dAtA.DurationType = shared_proto.GoodsDurationEffectType(i)
	}

	dAtA.Duration, err = config.ParseDuration(pArSeR.String("duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("duration"), dAtA)
	}

	dAtA.AmountType = shared_proto.GoodsAmountType(shared_proto.GoodsAmountType_value[strings.ToUpper(pArSeR.String("amount_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("amount_type"), 10, 32); err == nil {
		dAtA.AmountType = shared_proto.GoodsAmountType(i)
	}

	dAtA.Amount = pArSeR.Uint64("amount")
	dAtA.BuffId = pArSeR.Uint64("buff_id")

	return dAtA, nil
}

var vAlIdAtOrGoodsEffect = map[string]*config.Validator{

	"id":                        config.ParseValidator("int>0", "", false, nil, nil),
	"gold":                      config.ParseValidator("uint", "", false, nil, nil),
	"food":                      config.ParseValidator("uint", "", false, nil, nil),
	"wood":                      config.ParseValidator("uint", "", false, nil, nil),
	"stone":                     config.ParseValidator("uint", "", false, nil, nil),
	"building_cdr":              config.ParseValidator("bool", "", false, nil, nil),
	"tech_cdr":                  config.ParseValidator("bool", "", false, nil, nil),
	"workshop_cdr":              config.ParseValidator("bool", "", false, nil, nil),
	"cdr":                       config.ParseValidator("string", "", false, nil, nil),
	"move_base":                 config.ParseValidator("bool", "", false, nil, nil),
	"random_pos":                config.ParseValidator("bool", "", false, nil, nil),
	"guild_move_base":           config.ParseValidator("bool", "", false, nil, nil),
	"monster_move_sub_level":    config.ParseValidator("uint", "", false, nil, nil),
	"exp_type":                  config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.GoodsExpEffectType_value, 0), nil),
	"exp":                       config.ParseValidator("uint", "", false, nil, nil),
	"mian_duration":             config.ParseValidator("string", "", false, nil, nil),
	"troop_speed_up_rate":       config.ParseValidator("float64", "", false, nil, nil),
	"train_duration":            config.ParseValidator("string", "", false, nil, nil),
	"parts_combine_count":       config.ParseValidator("uint", "", false, nil, nil),
	"parts_plunder_prize_id":    config.ParseValidator("uint", "", true, nil, nil),
	"parts_show_type":           config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"add_multi_level_npc_times": config.ParseValidator("bool", "", false, nil, nil),
	"add_invase_hero_times":     config.ParseValidator("bool", "", false, nil, nil),
	"duration_type":             config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.GoodsDurationEffectType_value, 0), nil),
	"duration":                  config.ParseValidator("string", "", false, nil, nil),
	"amount_type":               config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.GoodsAmountType_value, 0), nil),
	"amount":                    config.ParseValidator("uint", "", false, nil, nil),
	"buff_id":                   config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *GoodsEffect) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GoodsEffect) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GoodsEffect) Encode() *shared_proto.GoodsEffectProto {
	out := &shared_proto.GoodsEffectProto{}
	out.Gold = config.U64ToI32(dAtA.Gold)
	out.Food = config.U64ToI32(dAtA.Food)
	out.Wood = config.U64ToI32(dAtA.Wood)
	out.Stone = config.U64ToI32(dAtA.Stone)
	out.BuildingCdr = dAtA.BuildingCdr
	out.TechCdr = dAtA.TechCdr
	out.WorkshopCdr = dAtA.WorkshopCdr
	out.Cdr = config.Duration2I32Seconds(dAtA.Cdr)
	out.MoveBaseType = dAtA.MoveBaseType
	out.MoveBase = dAtA.MoveBase
	out.RandomPos = dAtA.RandomPos
	out.GuildMoveBase = dAtA.GuildMoveBase
	out.ExpType = dAtA.ExpType
	out.Exp = config.U64ToI32(dAtA.Exp)
	out.MianDuration = config.Duration2I32Seconds(dAtA.MianDuration)
	out.TroopSpeedUpRate = config.F64ToI32X1000(dAtA.TroopSpeedUpRate)
	out.TrainDuration = config.Duration2I32Seconds(dAtA.TrainDuration)
	out.PartsCombineCount = config.U64ToI32(dAtA.PartsCombineCount)
	if dAtA.PartsShowPrize != nil {
		out.PartsShowPrize = dAtA.PartsShowPrize
	}
	out.PartsShowType = config.U64ToI32(dAtA.PartsShowType)
	out.AddMultiLevelNpcTimes = dAtA.AddMultiLevelNpcTimes
	out.AddInvaseHeroTimes = dAtA.AddInvaseHeroTimes
	out.DurationType = dAtA.DurationType
	out.Duration = config.Duration2I32Seconds(dAtA.Duration)
	out.AmountType = dAtA.AmountType
	out.Amount = config.U64ToI32(dAtA.Amount)
	out.BuffId = config.U64ToI32(dAtA.BuffId)

	return out
}

func ArrayEncodeGoodsEffect(datas []*GoodsEffect) []*shared_proto.GoodsEffectProto {

	out := make([]*shared_proto.GoodsEffectProto, 0, len(datas))
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

func (dAtA *GoodsEffect) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GoodsQuality ----------------------------------

func LoadGoodsQuality(gos *config.GameObjects) (map[uint64]*GoodsQuality, map[*GoodsQuality]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GoodsQualityPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GoodsQuality, len(lIsT))
	pArSeRmAp := make(map[*GoodsQuality]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGoodsQuality) {
			continue
		}

		dAtA, err := NewGoodsQuality(fIlEnAmE, pArSeR)
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

func SetRelatedGoodsQuality(dAtAmAp map[*GoodsQuality]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GoodsQualityPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGoodsQualityKeyArray(datas []*GoodsQuality) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewGoodsQuality(fIlEnAmE string, pArSeR *config.ObjectParser) (*GoodsQuality, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGoodsQuality)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GoodsQuality{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Quality = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("quality"))])
	if i, err := strconv.ParseInt(pArSeR.String("quality"), 10, 32); err == nil {
		dAtA.Quality = shared_proto.Quality(i)
	}

	return dAtA, nil
}

var vAlIdAtOrGoodsQuality = map[string]*config.Validator{

	"level":   config.ParseValidator("uint", "", false, nil, nil),
	"quality": config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), nil),
}

func (dAtA *GoodsQuality) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GoodsQuality) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GoodsQuality) Encode() *shared_proto.GoodsQualityProto {
	out := &shared_proto.GoodsQualityProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Quality = dAtA.Quality

	return out
}

func ArrayEncodeGoodsQuality(datas []*GoodsQuality) []*shared_proto.GoodsQualityProto {

	out := make([]*shared_proto.GoodsQualityProto, 0, len(datas))
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

func (dAtA *GoodsQuality) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetEquipmentQualityData(uint64) *EquipmentQualityData
	GetEquipmentRefinedData(uint64) *EquipmentRefinedData
	GetGoodsQuality(uint64) *GoodsQuality
	GetIcon(string) *icon.Icon
	GetSpriteStat(uint64) *data.SpriteStat
}
