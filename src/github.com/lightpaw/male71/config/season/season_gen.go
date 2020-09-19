// AUTO_GEN, DONT MODIFY!!!
package season

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
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

// start with SeasonData ----------------------------------

func LoadSeasonData(gos *config.GameObjects) (map[uint64]*SeasonData, map[*SeasonData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.SeasonDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*SeasonData, len(lIsT))
	pArSeRmAp := make(map[*SeasonData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrSeasonData) {
			continue
		}

		dAtA, err := NewSeasonData(fIlEnAmE, pArSeR)
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

func SetRelatedSeasonData(dAtAmAp map[*SeasonData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SeasonDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetSeasonDataKeyArray(datas []*SeasonData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewSeasonData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SeasonData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSeasonData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SeasonData{}

	dAtA.Season = shared_proto.Season(shared_proto.Season_value[strings.ToUpper(pArSeR.String("season"))])
	if i, err := strconv.ParseInt(pArSeR.String("season"), 10, 32); err == nil {
		dAtA.Season = shared_proto.Season(i)
	}

	dAtA.Name = pArSeR.String("name")
	dAtA.BgImg = pArSeR.String("bg_img")
	// releated field: ShowPrize
	// releated field: Prize
	dAtA.WorkerCdr = pArSeR.Float64("worker_cdr")
	dAtA.SecretTowerTimes = pArSeR.Uint64("secret_tower_times")
	dAtA.FarmBaseInc = pArSeR.Float64("farm_base_inc")
	dAtA.AddMultiMonsterTimes = pArSeR.Uint64("add_multi_monster_times")
	dAtA.DecTroopSpeedRate = pArSeR.Float64("dec_troop_speed_rate")
	dAtA.IncProsperityMultiple = pArSeR.Float64("inc_prosperity_multiple")
	// skip field: PrevSeason

	// calculate fields
	dAtA.Id = uint64(dAtA.Season)

	return dAtA, nil
}

var vAlIdAtOrSeasonData = map[string]*config.Validator{

	"season":                  config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Season_value, 0), nil),
	"name":                    config.ParseValidator("string", "", false, nil, nil),
	"bg_img":                  config.ParseValidator("string", "", false, nil, nil),
	"show_prize":              config.ParseValidator("string", "", false, nil, nil),
	"prize":                   config.ParseValidator("string", "", false, nil, nil),
	"worker_cdr":              config.ParseValidator("float64>=0", "", false, nil, nil),
	"secret_tower_times":      config.ParseValidator("uint", "", false, nil, nil),
	"farm_base_inc":           config.ParseValidator("float64>=0", "", false, nil, nil),
	"add_multi_monster_times": config.ParseValidator("uint", "", false, nil, nil),
	"dec_troop_speed_rate":    config.ParseValidator("float64>=0", "", false, nil, nil),
	"inc_prosperity_multiple": config.ParseValidator("float64>=0", "", false, nil, nil),
}

func (dAtA *SeasonData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SeasonData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SeasonData) Encode() *shared_proto.SeasonDataProto {
	out := &shared_proto.SeasonDataProto{}
	out.Season = dAtA.Season
	out.BgImg = dAtA.BgImg
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.PrizeProto()
	}
	out.WorkerCdr = config.F64ToI32X1000(dAtA.WorkerCdr)
	out.SecretTowerTimes = config.U64ToI32(dAtA.SecretTowerTimes)
	out.FarmBaseInc = config.F64ToI32X1000(dAtA.FarmBaseInc)
	out.AddMultiMonsterTimes = config.U64ToI32(dAtA.AddMultiMonsterTimes)
	out.DecTroopSpeedRate = config.F64ToI32X1000(dAtA.DecTroopSpeedRate)
	out.IncProsperityMultiple = config.F64ToI32X1000(dAtA.IncProsperityMultiple)

	return out
}

func ArrayEncodeSeasonData(datas []*SeasonData) []*shared_proto.SeasonDataProto {

	out := make([]*shared_proto.SeasonDataProto, 0, len(datas))
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

func (dAtA *SeasonData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil {
		return errors.Errorf("%s 配置的关联字段[show_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with SeasonMiscData ----------------------------------

func LoadSeasonMiscData(gos *config.GameObjects) (*SeasonMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.SeasonMiscDataPath
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

	dAtA, err := NewSeasonMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedSeasonMiscData(gos *config.GameObjects, dAtA *SeasonMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SeasonMiscDataPath
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

func NewSeasonMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SeasonMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSeasonMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SeasonMiscData{}

	// skip field: SeasonDuration
	if pArSeR.KeyExist("season_switch_duration") {
		dAtA.SeasonSwitchDuration, err = config.ParseDuration(pArSeR.String("season_switch_duration"))
	} else {
		dAtA.SeasonSwitchDuration, err = config.ParseDuration("10s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[season_switch_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("season_switch_duration"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrSeasonMiscData = map[string]*config.Validator{

	"season_switch_duration": config.ParseValidator("string", "", false, nil, []string{"10s"}),
}

func (dAtA *SeasonMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SeasonMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SeasonMiscData) Encode() *shared_proto.SeasonMiscProto {
	out := &shared_proto.SeasonMiscProto{}
	out.SeasonDuration = config.Duration2I32Seconds(dAtA.SeasonDuration)
	out.SeasonSwitchDuration = config.Duration2I32Seconds(dAtA.SeasonSwitchDuration)

	return out
}

func ArrayEncodeSeasonMiscData(datas []*SeasonMiscData) []*shared_proto.SeasonMiscProto {

	out := make([]*shared_proto.SeasonMiscProto, 0, len(datas))
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

func (dAtA *SeasonMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetPrize(int) *resdata.Prize
}
