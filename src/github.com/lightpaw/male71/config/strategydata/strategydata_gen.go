// AUTO_GEN, DONT MODIFY!!!
package strategydata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/domestic_data"
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

// start with StrategyData ----------------------------------

func LoadStrategyData(gos *config.GameObjects) (map[uint64]*StrategyData, map[*StrategyData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.StrategyDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*StrategyData, len(lIsT))
	pArSeRmAp := make(map[*StrategyData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrStrategyData) {
			continue
		}

		dAtA, err := NewStrategyData(fIlEnAmE, pArSeR)
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

func SetRelatedStrategyData(dAtAmAp map[*StrategyData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.StrategyDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetStrategyDataKeyArray(datas []*StrategyData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewStrategyData(fIlEnAmE string, pArSeR *config.ObjectParser) (*StrategyData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrStrategyData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &StrategyData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Type = pArSeR.Uint64("type")
	dAtA.Target = pArSeR.Uint64("target")
	dAtA.Name = pArSeR.String("name")
	dAtA.Sp = 0
	if pArSeR.KeyExist("sp") {
		dAtA.Sp = pArSeR.Uint64("sp")
	}

	dAtA.UnlockHeroLevel = pArSeR.Uint64("unlock_hero_level")
	dAtA.Cd, err = config.ParseDuration(pArSeR.String("cd"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("cd"), dAtA)
	}

	dAtA.TodayLimit = pArSeR.Uint64("today_limit")
	// releated field: Icon
	dAtA.Desc = pArSeR.String("desc")
	// skip field: EffectMap

	return dAtA, nil
}

var vAlIdAtOrStrategyData = map[string]*config.Validator{

	"id":                config.ParseValidator("int>0", "", false, nil, nil),
	"type":              config.ParseValidator("int>0", "", false, nil, nil),
	"target":            config.ParseValidator("int>0", "", false, nil, nil),
	"name":              config.ParseValidator("string", "", false, nil, nil),
	"sp":                config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"unlock_hero_level": config.ParseValidator("int>0", "", false, nil, nil),
	"cd":                config.ParseValidator("string", "", false, nil, nil),
	"today_limit":       config.ParseValidator("int>0", "", false, nil, nil),
	"icon":              config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"desc":              config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *StrategyData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *StrategyData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *StrategyData) Encode() *shared_proto.StrategyDataProto {
	out := &shared_proto.StrategyDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Type = config.U64ToI32(dAtA.Type)
	out.Target = config.U64ToI32(dAtA.Target)
	out.Name = dAtA.Name
	out.Sp = config.U64ToI32(dAtA.Sp)
	out.UnlockHeroLevel = config.U64ToI32(dAtA.UnlockHeroLevel)
	out.Cd = config.Duration2I32Seconds(dAtA.Cd)
	out.TodayLimit = config.U64ToI32(dAtA.TodayLimit)
	if dAtA.Icon != nil {
		out.Icon = dAtA.Icon.Id
	}
	out.Desc = dAtA.Desc

	return out
}

func ArrayEncodeStrategyData(datas []*StrategyData) []*shared_proto.StrategyDataProto {

	out := make([]*shared_proto.StrategyDataProto, 0, len(datas))
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

func (dAtA *StrategyData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with StrategyEffectData ----------------------------------

func LoadStrategyEffectData(gos *config.GameObjects) (map[uint64]*StrategyEffectData, map[*StrategyEffectData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.StrategyEffectDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*StrategyEffectData, len(lIsT))
	pArSeRmAp := make(map[*StrategyEffectData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrStrategyEffectData) {
			continue
		}

		dAtA, err := NewStrategyEffectData(fIlEnAmE, pArSeR)
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

func SetRelatedStrategyEffectData(dAtAmAp map[*StrategyEffectData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.StrategyEffectDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetStrategyEffectDataKeyArray(datas []*StrategyEffectData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewStrategyEffectData(fIlEnAmE string, pArSeR *config.ObjectParser) (*StrategyEffectData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrStrategyEffectData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &StrategyEffectData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.StrategyId = pArSeR.Uint64("strategy_id")
	dAtA.IntEffectType = pArSeR.Int("int_effect_type")
	// skip field: EffectType
	dAtA.HeroLevel = pArSeR.Uint64("hero_level")
	// releated field: Cost
	// releated field: Prize
	dAtA.FarmFastHarvestDuration, err = config.ParseDuration(pArSeR.String("farm_fast_harvest_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[farm_fast_harvest_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("farm_fast_harvest_duration"), dAtA)
	}

	dAtA.TargetReduceSolider = pArSeR.Uint64("target_reduce_solider")

	return dAtA, nil
}

var vAlIdAtOrStrategyEffectData = map[string]*config.Validator{

	"id":              config.ParseValidator("int>0", "", false, nil, nil),
	"strategy_id":     config.ParseValidator("int>0", "", false, nil, nil),
	"int_effect_type": config.ParseValidator("int>0", "", false, nil, nil),
	"hero_level":      config.ParseValidator("int>0", "", false, nil, nil),
	"cost":            config.ParseValidator("string", "", false, nil, nil),
	"prize":           config.ParseValidator("string", "", false, nil, nil),
	"farm_fast_harvest_duration": config.ParseValidator("string", "", false, nil, nil),
	"target_reduce_solider":      config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *StrategyEffectData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *StrategyEffectData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *StrategyEffectData) Encode() *shared_proto.StrategyEffectDataProto {
	out := &shared_proto.StrategyEffectDataProto{}
	out.StrategyId = config.U64ToI32(dAtA.StrategyId)
	out.EffectType = dAtA.EffectType
	out.HeroLevel = config.U64ToI32(dAtA.HeroLevel)
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	out.FarmFastHarvestDuration = config.Duration2I32Seconds(dAtA.FarmFastHarvestDuration)
	out.TargetReduceSolider = config.U64ToI32(dAtA.TargetReduceSolider)

	return out
}

func ArrayEncodeStrategyEffectData(datas []*StrategyEffectData) []*shared_proto.StrategyEffectDataProto {

	out := make([]*shared_proto.StrategyEffectDataProto, 0, len(datas))
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

func (dAtA *StrategyEffectData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

type related_configs interface {
	GetCombineCost(int) *domestic_data.CombineCost
	GetIcon(string) *icon.Icon
	GetPrize(int) *resdata.Prize
}
