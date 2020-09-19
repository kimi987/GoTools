// AUTO_GEN, DONT MODIFY!!!
package farm

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

// start with FarmMaxStealConfig ----------------------------------

func LoadFarmMaxStealConfig(gos *config.GameObjects) (map[uint64]*FarmMaxStealConfig, map[*FarmMaxStealConfig]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FarmMaxStealConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*FarmMaxStealConfig, len(lIsT))
	pArSeRmAp := make(map[*FarmMaxStealConfig]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFarmMaxStealConfig) {
			continue
		}

		dAtA, err := NewFarmMaxStealConfig(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.GuanFuLevel
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[GuanFuLevel], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedFarmMaxStealConfig(dAtAmAp map[*FarmMaxStealConfig]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FarmMaxStealConfigPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFarmMaxStealConfigKeyArray(datas []*FarmMaxStealConfig) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.GuanFuLevel)
		}
	}

	return out
}

func NewFarmMaxStealConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*FarmMaxStealConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFarmMaxStealConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FarmMaxStealConfig{}

	dAtA.GuanFuLevel = pArSeR.Uint64("guan_fu_level")
	dAtA.MaxDailyStealGoldAmount = pArSeR.Uint64("max_daily_steal_gold_amount")
	dAtA.MaxDailyStealStoneAmount = pArSeR.Uint64("max_daily_steal_stone_amount")

	return dAtA, nil
}

var vAlIdAtOrFarmMaxStealConfig = map[string]*config.Validator{

	"guan_fu_level":                config.ParseValidator("uint", "", false, nil, nil),
	"max_daily_steal_gold_amount":  config.ParseValidator("uint", "", false, nil, nil),
	"max_daily_steal_stone_amount": config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *FarmMaxStealConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with FarmMiscConfig ----------------------------------

func LoadFarmMiscConfig(gos *config.GameObjects) (*FarmMiscConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.FarmMiscConfigPath
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

	dAtA, err := NewFarmMiscConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedFarmMiscConfig(gos *config.GameObjects, dAtA *FarmMiscConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FarmMiscConfigPath
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

func NewFarmMiscConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*FarmMiscConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFarmMiscConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FarmMiscConfig{}

	dAtA.EarlyHarvestPercent, err = data.ParseAmount(pArSeR.String("early_harvest_percent"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[early_harvest_percent] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("early_harvest_percent"), dAtA)
	}

	dAtA.RipeProtectDuration, err = config.ParseDuration(pArSeR.String("ripe_protect_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[ripe_protect_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("ripe_protect_duration"), dAtA)
	}

	// skip field: NegRipeProtectDuration
	dAtA.StealGainPercent, err = data.ParseAmount(pArSeR.String("steal_gain_percent"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[steal_gain_percent] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("steal_gain_percent"), dAtA)
	}

	dAtA.CubeStealMaxTime = pArSeR.Uint64("cube_steal_max_time")
	dAtA.StealLogMaxCount = pArSeR.Uint64("steal_log_max_count")
	dAtA.StealLogExpiredHours = pArSeR.String("steal_log_expired_hours")
	// skip field: StealLogExpiredDuration
	// skip field: OneKeyResConfig

	return dAtA, nil
}

var vAlIdAtOrFarmMiscConfig = map[string]*config.Validator{

	"early_harvest_percent":   config.ParseValidator("string", "", false, nil, nil),
	"ripe_protect_duration":   config.ParseValidator("string", "", false, nil, nil),
	"steal_gain_percent":      config.ParseValidator("string", "", false, nil, nil),
	"cube_steal_max_time":     config.ParseValidator("uint", "", false, nil, nil),
	"steal_log_max_count":     config.ParseValidator("uint", "", false, nil, nil),
	"steal_log_expired_hours": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *FarmMiscConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *FarmMiscConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *FarmMiscConfig) Encode() *shared_proto.FarmMiscConfigProto {
	out := &shared_proto.FarmMiscConfigProto{}
	if dAtA.EarlyHarvestPercent != nil {
		out.EarlyHarvestPercent = dAtA.EarlyHarvestPercent.Encode()
	}
	out.RipeProtectDuration = config.Duration2I32Seconds(dAtA.RipeProtectDuration)
	if dAtA.StealGainPercent != nil {
		out.StealGainPercent = dAtA.StealGainPercent.Encode()
	}
	out.CubeStealMaxTime = config.U64ToI32(dAtA.CubeStealMaxTime)
	out.StealLogMaxCount = config.U64ToI32(dAtA.StealLogMaxCount)

	return out
}

func ArrayEncodeFarmMiscConfig(datas []*FarmMiscConfig) []*shared_proto.FarmMiscConfigProto {

	out := make([]*shared_proto.FarmMiscConfigProto, 0, len(datas))
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

func (dAtA *FarmMiscConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with FarmOneKeyConfig ----------------------------------

func LoadFarmOneKeyConfig(gos *config.GameObjects) (map[uint64]*FarmOneKeyConfig, map[*FarmOneKeyConfig]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FarmOneKeyConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*FarmOneKeyConfig, len(lIsT))
	pArSeRmAp := make(map[*FarmOneKeyConfig]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFarmOneKeyConfig) {
			continue
		}

		dAtA, err := NewFarmOneKeyConfig(fIlEnAmE, pArSeR)
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

func SetRelatedFarmOneKeyConfig(dAtAmAp map[*FarmOneKeyConfig]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FarmOneKeyConfigPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFarmOneKeyConfigKeyArray(datas []*FarmOneKeyConfig) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.BaseLevel)
		}
	}

	return out
}

func NewFarmOneKeyConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*FarmOneKeyConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFarmOneKeyConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FarmOneKeyConfig{}

	dAtA.BaseLevel = pArSeR.Uint64("base_level")
	dAtA.StoneHopeCubeCount = pArSeR.Uint64("stone_hope_cube_count")
	dAtA.GoldHopeCubeCount = pArSeR.Uint64("gold_hope_cube_count")
	// skip field: MaxDailyStealGoldAmount
	// skip field: MaxDailyStealStoneAmount

	return dAtA, nil
}

var vAlIdAtOrFarmOneKeyConfig = map[string]*config.Validator{

	"base_level":            config.ParseValidator("uint", "", false, nil, nil),
	"stone_hope_cube_count": config.ParseValidator("uint", "", false, nil, nil),
	"gold_hope_cube_count":  config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *FarmOneKeyConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with FarmResConfig ----------------------------------

func LoadFarmResConfig(gos *config.GameObjects) (map[uint64]*FarmResConfig, map[*FarmResConfig]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FarmResConfigPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*FarmResConfig, len(lIsT))
	pArSeRmAp := make(map[*FarmResConfig]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFarmResConfig) {
			continue
		}

		dAtA, err := NewFarmResConfig(fIlEnAmE, pArSeR)
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

func SetRelatedFarmResConfig(dAtAmAp map[*FarmResConfig]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FarmResConfigPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFarmResConfigKeyArray(datas []*FarmResConfig) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewFarmResConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*FarmResConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFarmResConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FarmResConfig{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.RipeDuration, err = time.ParseDuration(pArSeR.String("ripe_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[ripe_duration] 解析失败(time.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("ripe_duration"), dAtA)
	}

	dAtA.ResType = shared_proto.ResType(shared_proto.ResType_value[strings.ToUpper(pArSeR.String("res_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("res_type"), 10, 32); err == nil {
		dAtA.ResType = shared_proto.ResType(i)
	}

	dAtA.BaseOutput, err = data.ParseAmount(pArSeR.String("base_output"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[base_output] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("base_output"), dAtA)
	}

	dAtA.Icon = pArSeR.String("icon")

	return dAtA, nil
}

var vAlIdAtOrFarmResConfig = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"ripe_duration": config.ParseValidator("string", "", false, nil, nil),
	"res_type":      config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.ResType_value, 0), nil),
	"base_output":   config.ParseValidator("string", "", false, nil, nil),
	"icon":          config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *FarmResConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *FarmResConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *FarmResConfig) Encode() *shared_proto.FarmResConfigProto {
	out := &shared_proto.FarmResConfigProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.RipeDuration = config.Duration2I32Seconds(dAtA.RipeDuration)
	out.ResType = dAtA.ResType
	if dAtA.BaseOutput != nil {
		out.BaseOutput = dAtA.BaseOutput.Encode()
	}
	out.Icon = dAtA.Icon

	return out
}

func ArrayEncodeFarmResConfig(datas []*FarmResConfig) []*shared_proto.FarmResConfigProto {

	out := make([]*shared_proto.FarmResConfigProto, 0, len(datas))
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

func (dAtA *FarmResConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
}
