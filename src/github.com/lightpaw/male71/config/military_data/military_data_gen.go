// AUTO_GEN, DONT MODIFY!!!
package military_data

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

// start with JiuGuanData ----------------------------------

func LoadJiuGuanData(gos *config.GameObjects) (map[uint64]*JiuGuanData, map[*JiuGuanData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.JiuGuanDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*JiuGuanData, len(lIsT))
	pArSeRmAp := make(map[*JiuGuanData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrJiuGuanData) {
			continue
		}

		dAtA, err := NewJiuGuanData(fIlEnAmE, pArSeR)
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

func SetRelatedJiuGuanData(dAtAmAp map[*JiuGuanData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.JiuGuanDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetJiuGuanDataKeyArray(datas []*JiuGuanData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewJiuGuanData(fIlEnAmE string, pArSeR *config.ObjectParser) (*JiuGuanData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrJiuGuanData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &JiuGuanData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.MaxTimes = pArSeR.Uint64("max_times")
	if pArSeR.KeyExist("recovery_duration") {
		dAtA.RecoveryDuration, err = config.ParseDuration(pArSeR.String("recovery_duration"))
	} else {
		dAtA.RecoveryDuration, err = config.ParseDuration("3m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[recovery_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("recovery_duration"), dAtA)
	}

	// releated field: TutorDatas
	dAtA.RandomTutorWeight = pArSeR.Uint64Array("random_tutor_weight", "", false)
	dAtA.InitTutorWeight = pArSeR.Uint64Array("init_tutor_weight", "", false)
	dAtA.BroadcastContent = pArSeR.String("broadcast_content")
	dAtA.BroadcastMinCritMulti = pArSeR.Uint64("broadcast_min_crit_multi")

	return dAtA, nil
}

var vAlIdAtOrJiuGuanData = map[string]*config.Validator{

	"level":                    config.ParseValidator("uint", "", false, nil, nil),
	"max_times":                config.ParseValidator("int>0", "", false, nil, nil),
	"recovery_duration":        config.ParseValidator("string", "", false, nil, []string{"3m"}),
	"tutor_datas":              config.ParseValidator("string", "", true, nil, nil),
	"random_tutor_weight":      config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
	"init_tutor_weight":        config.ParseValidator("int,duplicate,notAllNil", "", true, nil, nil),
	"broadcast_content":        config.ParseValidator("string>0", "", false, nil, nil),
	"broadcast_min_crit_multi": config.ParseValidator("int", "", false, nil, nil),
}

func (dAtA *JiuGuanData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *JiuGuanData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *JiuGuanData) Encode() *shared_proto.JiuGuanDataProto {
	out := &shared_proto.JiuGuanDataProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.MaxTimes = config.U64ToI32(dAtA.MaxTimes)
	out.RecoveryDuration = config.Duration2I32Seconds(dAtA.RecoveryDuration)
	if dAtA.TutorDatas != nil {
		out.TutorDatas = ArrayEncodeTutorData(dAtA.TutorDatas)
	}
	out.BroadcastContent = dAtA.BroadcastContent

	return out
}

func ArrayEncodeJiuGuanData(datas []*JiuGuanData) []*shared_proto.JiuGuanDataProto {

	out := make([]*shared_proto.JiuGuanDataProto, 0, len(datas))
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

func (dAtA *JiuGuanData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("tutor_datas", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetTutorData(v)
		if obj != nil {
			dAtA.TutorDatas = append(dAtA.TutorDatas, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[tutor_datas] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("tutor_datas"), *pArSeR)
		}
	}

	return nil
}

// start with JiuGuanMiscData ----------------------------------

func LoadJiuGuanMiscData(gos *config.GameObjects) (*JiuGuanMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.JiuGuanMiscDataPath
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

	dAtA, err := NewJiuGuanMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedJiuGuanMiscData(gos *config.GameObjects, dAtA *JiuGuanMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.JiuGuanMiscDataPath
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

func NewJiuGuanMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*JiuGuanMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrJiuGuanMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &JiuGuanMiscData{}

	dAtA.RefreshCostYuanBao = []uint64{0, 10, 20, 30, 40, 50}
	if pArSeR.KeyExist("refresh_cost_yuan_bao") {
		dAtA.RefreshCostYuanBao = pArSeR.Uint64Array("refresh_cost_yuan_bao", "", false)
	}

	// releated field: RefreshCost
	dAtA.RecoveryDuration, err = config.ParseDuration(pArSeR.String("recovery_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[recovery_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("recovery_duration"), dAtA)
	}

	dAtA.FirstRefreshIndex = 2
	if pArSeR.KeyExist("first_refresh_index") {
		dAtA.FirstRefreshIndex = pArSeR.Uint64("first_refresh_index")
	}

	return dAtA, nil
}

var vAlIdAtOrJiuGuanMiscData = map[string]*config.Validator{

	"refresh_cost_yuan_bao": config.ParseValidator("uint", "", true, nil, []string{"0", "10", "20", "30", "40", "50"}),
	"refresh_cost":          config.ParseValidator("string", "", true, nil, nil),
	"recovery_duration":     config.ParseValidator("string", "", false, nil, nil),
	"first_refresh_index":   config.ParseValidator("int>0", "", false, nil, []string{"2"}),
}

func (dAtA *JiuGuanMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *JiuGuanMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *JiuGuanMiscData) Encode() *shared_proto.JiuGuanMiscDataProto {
	out := &shared_proto.JiuGuanMiscDataProto{}
	out.RefreshCostYuanBao = config.U64a2I32a(dAtA.RefreshCostYuanBao)
	if dAtA.RefreshCost != nil {
		out.RefreshCost = resdata.ArrayEncodeCost(dAtA.RefreshCost)
	}
	out.RecoveryDuration = config.Duration2I32Seconds(dAtA.RecoveryDuration)

	return out
}

func ArrayEncodeJiuGuanMiscData(datas []*JiuGuanMiscData) []*shared_proto.JiuGuanMiscDataProto {

	out := make([]*shared_proto.JiuGuanMiscDataProto, 0, len(datas))
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

func (dAtA *JiuGuanMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	intKeys = pArSeR.IntArray("refresh_cost", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetCost(v)
		if obj != nil {
			dAtA.RefreshCost = append(dAtA.RefreshCost, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[refresh_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("refresh_cost"), *pArSeR)
		}
	}

	return nil
}

// start with JunYingLevelData ----------------------------------

func LoadJunYingLevelData(gos *config.GameObjects) (map[uint64]*JunYingLevelData, map[*JunYingLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.JunYingLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*JunYingLevelData, len(lIsT))
	pArSeRmAp := make(map[*JunYingLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrJunYingLevelData) {
			continue
		}

		dAtA, err := NewJunYingLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedJunYingLevelData(dAtAmAp map[*JunYingLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.JunYingLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetJunYingLevelDataKeyArray(datas []*JunYingLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewJunYingLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*JunYingLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrJunYingLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &JunYingLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.MaxTimes = pArSeR.Uint64("max_times")
	dAtA.RecoveryDuration, err = config.ParseDuration(pArSeR.String("recovery_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[recovery_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("recovery_duration"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrJunYingLevelData = map[string]*config.Validator{

	"level":             config.ParseValidator("uint", "", false, nil, nil),
	"max_times":         config.ParseValidator("uint", "", false, nil, nil),
	"recovery_duration": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *JunYingLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *JunYingLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *JunYingLevelData) Encode() *shared_proto.JunYingLevelDataProto {
	out := &shared_proto.JunYingLevelDataProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.MaxTimes = config.U64ToI32(dAtA.MaxTimes)
	out.RecoveryDuration = config.Duration2I32Seconds(dAtA.RecoveryDuration)

	return out
}

func ArrayEncodeJunYingLevelData(datas []*JunYingLevelData) []*shared_proto.JunYingLevelDataProto {

	out := make([]*shared_proto.JunYingLevelDataProto, 0, len(datas))
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

func (dAtA *JunYingLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with JunYingMiscData ----------------------------------

func LoadJunYingMiscData(gos *config.GameObjects) (*JunYingMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.JunYingMiscDataPath
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

	dAtA, err := NewJunYingMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedJunYingMiscData(gos *config.GameObjects, dAtA *JunYingMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.JunYingMiscDataPath
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

func NewJunYingMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*JunYingMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrJunYingMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &JunYingMiscData{}

	dAtA.DefaultTimes = pArSeR.Uint64("default_times")
	dAtA.ForceAddSoldierMaxTimes = pArSeR.Uint64("force_add_soldier_max_times")
	dAtA.ForceAddSoldierCost = []uint64{10, 50, 100, 200, 300, 500}
	if pArSeR.KeyExist("force_add_soldier_cost") {
		dAtA.ForceAddSoldierCost = pArSeR.Uint64Array("force_add_soldier_cost", "", false)
	}

	// releated field: ForceAddSoldierNewCost

	return dAtA, nil
}

var vAlIdAtOrJunYingMiscData = map[string]*config.Validator{

	"default_times":               config.ParseValidator("uint", "", false, nil, nil),
	"force_add_soldier_max_times": config.ParseValidator("uint", "", false, nil, nil),
	"force_add_soldier_cost":      config.ParseValidator("uint,duplicate", "", true, nil, []string{"10", "50", "100", "200", "300", "500"}),
	"force_add_soldier_new_cost":  config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *JunYingMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *JunYingMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *JunYingMiscData) Encode() *shared_proto.JunYingMiscProto {
	out := &shared_proto.JunYingMiscProto{}
	out.ForceAddSoldierMaxTimes = config.U64ToI32(dAtA.ForceAddSoldierMaxTimes)
	out.ForceAddSoldierCost = config.U64a2I32a(dAtA.ForceAddSoldierCost)
	if dAtA.ForceAddSoldierNewCost != nil {
		out.ForceAddSoldierNewCost = resdata.ArrayEncodeCost(dAtA.ForceAddSoldierNewCost)
	}

	return out
}

func ArrayEncodeJunYingMiscData(datas []*JunYingMiscData) []*shared_proto.JunYingMiscProto {

	out := make([]*shared_proto.JunYingMiscProto, 0, len(datas))
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

func (dAtA *JunYingMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	intKeys = pArSeR.IntArray("force_add_soldier_new_cost", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetCost(v)
		if obj != nil {
			dAtA.ForceAddSoldierNewCost = append(dAtA.ForceAddSoldierNewCost, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[force_add_soldier_new_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("force_add_soldier_new_cost"), *pArSeR)
		}
	}

	return nil
}

// start with TrainingLevelData ----------------------------------

func LoadTrainingLevelData(gos *config.GameObjects) (map[uint64]*TrainingLevelData, map[*TrainingLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TrainingLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TrainingLevelData, len(lIsT))
	pArSeRmAp := make(map[*TrainingLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTrainingLevelData) {
			continue
		}

		dAtA, err := NewTrainingLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedTrainingLevelData(dAtAmAp map[*TrainingLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TrainingLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTrainingLevelDataKeyArray(datas []*TrainingLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewTrainingLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TrainingLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTrainingLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TrainingLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Coef = pArSeR.Float64("coef")
	// releated field: Cost

	return dAtA, nil
}

var vAlIdAtOrTrainingLevelData = map[string]*config.Validator{

	"level": config.ParseValidator("int>0", "", false, nil, nil),
	"name":  config.ParseValidator("string", "", false, nil, nil),
	"desc":  config.ParseValidator("string", "", false, nil, nil),
	"coef":  config.ParseValidator("float64>0", "", false, nil, nil),
	"cost":  config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *TrainingLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TrainingLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TrainingLevelData) Encode() *shared_proto.TrainingLevelProto {
	out := &shared_proto.TrainingLevelProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.Coef = config.F64ToI32X1000(dAtA.Coef)
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}

	return out
}

func ArrayEncodeTrainingLevelData(datas []*TrainingLevelData) []*shared_proto.TrainingLevelProto {

	out := make([]*shared_proto.TrainingLevelProto, 0, len(datas))
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

func (dAtA *TrainingLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with TutorData ----------------------------------

func LoadTutorData(gos *config.GameObjects) (map[uint64]*TutorData, map[*TutorData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TutorDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TutorData, len(lIsT))
	pArSeRmAp := make(map[*TutorData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTutorData) {
			continue
		}

		dAtA, err := NewTutorData(fIlEnAmE, pArSeR)
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

func SetRelatedTutorData(dAtAmAp map[*TutorData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TutorDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTutorDataKeyArray(datas []*TutorData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTutorData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TutorData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTutorData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TutorData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Image = pArSeR.String("image")
	dAtA.Weight = pArSeR.Uint64Array("weight", "", false)
	dAtA.Crit = pArSeR.Uint64Array("crit", "", false)
	dAtA.CritImgIndex = pArSeR.IntArray("crit_img_index", "", false)
	// releated field: Prize
	dAtA.ChatContent = pArSeR.String("chat_content")
	// releated field: RefreshMaxCost

	return dAtA, nil
}

var vAlIdAtOrTutorData = map[string]*config.Validator{

	"id":               config.ParseValidator("uint", "", false, nil, nil),
	"name":             config.ParseValidator("string>0", "", false, nil, nil),
	"image":            config.ParseValidator("string>0", "", false, nil, nil),
	"weight":           config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
	"crit":             config.ParseValidator("int>0,notAllNil", "", true, nil, nil),
	"crit_img_index":   config.ParseValidator("int,notAllNil", "", true, nil, nil),
	"prize":            config.ParseValidator("string", "", false, nil, nil),
	"chat_content":     config.ParseValidator("string", "", false, nil, nil),
	"refresh_max_cost": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *TutorData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TutorData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TutorData) Encode() *shared_proto.TutorDataProto {
	out := &shared_proto.TutorDataProto{}
	out.Name = dAtA.Name
	out.Image = dAtA.Image
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	out.ChatContent = dAtA.ChatContent
	if dAtA.RefreshMaxCost != nil {
		out.RefreshMaxCost = dAtA.RefreshMaxCost.Encode()
	}

	return out
}

func ArrayEncodeTutorData(datas []*TutorData) []*shared_proto.TutorDataProto {

	out := make([]*shared_proto.TutorDataProto, 0, len(datas))
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

func (dAtA *TutorData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	dAtA.RefreshMaxCost = cOnFigS.GetCost(pArSeR.Int("refresh_max_cost"))
	if dAtA.RefreshMaxCost == nil && pArSeR.Int("refresh_max_cost") != 0 {
		return errors.Errorf("%s 配置的关联字段[refresh_max_cost] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("refresh_max_cost"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetCost(int) *resdata.Cost
	GetPrize(int) *resdata.Prize
	GetTutorData(uint64) *TutorData
}
