// AUTO_GEN, DONT MODIFY!!!
package random_event

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

// start with EventOptionData ----------------------------------

func LoadEventOptionData(gos *config.GameObjects) (map[uint64]*EventOptionData, map[*EventOptionData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.EventOptionDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*EventOptionData, len(lIsT))
	pArSeRmAp := make(map[*EventOptionData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrEventOptionData) {
			continue
		}

		dAtA, err := NewEventOptionData(fIlEnAmE, pArSeR)
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

func SetRelatedEventOptionData(dAtAmAp map[*EventOptionData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EventOptionDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetEventOptionDataKeyArray(datas []*EventOptionData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewEventOptionData(fIlEnAmE string, pArSeR *config.ObjectParser) (*EventOptionData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEventOptionData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EventOptionData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Content = pArSeR.String("content")
	dAtA.Success = pArSeR.String("success")
	dAtA.Failed = pArSeR.String("failed")
	// releated field: Cost
	dAtA.FailedRate = 0
	if pArSeR.KeyExist("failed_rate") {
		dAtA.FailedRate = pArSeR.Uint64("failed_rate")
	}

	// releated field: SuccessPrize
	// releated field: FailedPrize

	return dAtA, nil
}

var vAlIdAtOrEventOptionData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"content":       config.ParseValidator("string>0", "", false, nil, nil),
	"success":       config.ParseValidator("string", "", false, nil, nil),
	"failed":        config.ParseValidator("string", "", false, nil, nil),
	"cost":          config.ParseValidator("string", "", false, nil, nil),
	"failed_rate":   config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"success_prize": config.ParseValidator("string", "", false, nil, nil),
	"failed_prize":  config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *EventOptionData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *EventOptionData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *EventOptionData) Encode() *shared_proto.EventOptionProto {
	out := &shared_proto.EventOptionProto{}
	out.OptionText = dAtA.Content
	out.SuccessText = dAtA.Success
	out.FailedText = dAtA.Failed
	if dAtA.Cost != nil {
		out.OptionCost = dAtA.Cost.Encode()
	}

	return out
}

func ArrayEncodeEventOptionData(datas []*EventOptionData) []*shared_proto.EventOptionProto {

	out := make([]*shared_proto.EventOptionProto, 0, len(datas))
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

func (dAtA *EventOptionData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil && pArSeR.Int("cost") != 0 {
		return errors.Errorf("%s 配置的关联字段[cost] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	dAtA.SuccessPrize = cOnFigS.GetOptionPrize(pArSeR.Uint64("success_prize"))
	if dAtA.SuccessPrize == nil && pArSeR.Uint64("success_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[success_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("success_prize"), *pArSeR)
	}

	dAtA.FailedPrize = cOnFigS.GetOptionPrize(pArSeR.Uint64("failed_prize"))
	if dAtA.FailedPrize == nil && pArSeR.Uint64("failed_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[failed_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("failed_prize"), *pArSeR)
	}

	return nil
}

// start with EventPosition ----------------------------------

func LoadEventPosition(gos *config.GameObjects) (map[uint64]*EventPosition, map[*EventPosition]*config.ObjectParser, error) {
	fIlEnAmE := confpath.EventPositionPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*EventPosition, len(lIsT))
	pArSeRmAp := make(map[*EventPosition]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrEventPosition) {
			continue
		}

		dAtA, err := NewEventPosition(fIlEnAmE, pArSeR)
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

func SetRelatedEventPosition(dAtAmAp map[*EventPosition]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EventPositionPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetEventPositionKeyArray(datas []*EventPosition) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewEventPosition(fIlEnAmE string, pArSeR *config.ObjectParser) (*EventPosition, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEventPosition)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EventPosition{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.PosX = pArSeR.Int("pos_x")
	dAtA.PosY = pArSeR.Int("pos_y")
	dAtA.TypeArea = pArSeR.Int("type_area")

	return dAtA, nil
}

var vAlIdAtOrEventPosition = map[string]*config.Validator{

	"id":        config.ParseValidator("int>0", "", false, nil, nil),
	"pos_x":     config.ParseValidator("int>0", "", false, nil, nil),
	"pos_y":     config.ParseValidator("int>0", "", false, nil, nil),
	"type_area": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *EventPosition) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with OptionPrize ----------------------------------

func LoadOptionPrize(gos *config.GameObjects) (map[uint64]*OptionPrize, map[*OptionPrize]*config.ObjectParser, error) {
	fIlEnAmE := confpath.OptionPrizePath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*OptionPrize, len(lIsT))
	pArSeRmAp := make(map[*OptionPrize]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrOptionPrize) {
			continue
		}

		dAtA, err := NewOptionPrize(fIlEnAmE, pArSeR)
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

func SetRelatedOptionPrize(dAtAmAp map[*OptionPrize]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.OptionPrizePath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetOptionPrizeKeyArray(datas []*OptionPrize) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewOptionPrize(fIlEnAmE string, pArSeR *config.ObjectParser) (*OptionPrize, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrOptionPrize)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &OptionPrize{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.GfAdd = pArSeR.Uint64("gf_add")
	// releated field: Prize
	dAtA.Weight = pArSeR.Uint64Array("weight", "", false)

	return dAtA, nil
}

var vAlIdAtOrOptionPrize = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"gf_add": config.ParseValidator("uint", "", false, nil, nil),
	"prize":  config.ParseValidator("string", "", true, nil, nil),
	"weight": config.ParseValidator("uint", "", true, nil, nil),
}

func (dAtA *OptionPrize) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	intKeys = pArSeR.IntArray("prize", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetPrize(v)
		if obj != nil {
			dAtA.Prize = append(dAtA.Prize, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
		}
	}

	return nil
}

// start with RandomEventData ----------------------------------

func LoadRandomEventData(gos *config.GameObjects) (map[uint64]*RandomEventData, map[*RandomEventData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.RandomEventDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*RandomEventData, len(lIsT))
	pArSeRmAp := make(map[*RandomEventData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrRandomEventData) {
			continue
		}

		dAtA, err := NewRandomEventData(fIlEnAmE, pArSeR)
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

func SetRelatedRandomEventData(dAtAmAp map[*RandomEventData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RandomEventDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetRandomEventDataKeyArray(datas []*RandomEventData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewRandomEventData(fIlEnAmE string, pArSeR *config.ObjectParser) (*RandomEventData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRandomEventData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RandomEventData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Title = pArSeR.String("title")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Content = pArSeR.String("content")
	dAtA.Image = pArSeR.String("image")
	dAtA.TypeSeason = shared_proto.Season(shared_proto.Season_value[strings.ToUpper(pArSeR.String("type_season"))])
	if i, err := strconv.ParseInt(pArSeR.String("type_season"), 10, 32); err == nil {
		dAtA.TypeSeason = shared_proto.Season(i)
	}

	dAtA.TypeArea = 0
	if pArSeR.KeyExist("type_area") {
		dAtA.TypeArea = pArSeR.Int("type_area")
	}

	// releated field: OptionDatas
	// skip field: OptionsProto4Send

	return dAtA, nil
}

var vAlIdAtOrRandomEventData = map[string]*config.Validator{

	"id":          config.ParseValidator("int>0", "", false, nil, nil),
	"title":       config.ParseValidator("string>0", "", false, nil, nil),
	"desc":        config.ParseValidator("string>0", "", false, nil, nil),
	"content":     config.ParseValidator("string", "", false, nil, nil),
	"image":       config.ParseValidator("string", "", false, nil, nil),
	"type_season": config.ParseValidator("int", "", false, config.EnumMapKeys(shared_proto.Season_value, 0), []string{"shared_proto.Season_InvalidSeason"}),
	"type_area":   config.ParseValidator("int", "", false, nil, []string{"0"}),
	"option":      config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *RandomEventData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *RandomEventData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *RandomEventData) Encode() *shared_proto.RandomEventDataProto {
	out := &shared_proto.RandomEventDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Title = dAtA.Title
	out.Desc = dAtA.Desc
	out.Content = dAtA.Content
	out.Image = dAtA.Image

	return out
}

func ArrayEncodeRandomEventData(datas []*RandomEventData) []*shared_proto.RandomEventDataProto {

	out := make([]*shared_proto.RandomEventDataProto, 0, len(datas))
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

func (dAtA *RandomEventData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("option", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetEventOptionData(v)
		if obj != nil {
			dAtA.OptionDatas = append(dAtA.OptionDatas, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[option] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("option"), *pArSeR)
		}
	}

	return nil
}

// start with RandomEventDataDictionary ----------------------------------

func LoadRandomEventDataDictionary(gos *config.GameObjects) (*RandomEventDataDictionary, *config.ObjectParser, error) {
	fIlEnAmE := confpath.RandomEventDataDictionaryPath
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

	dAtA, err := NewRandomEventDataDictionary(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedRandomEventDataDictionary(gos *config.GameObjects, dAtA *RandomEventDataDictionary, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RandomEventDataDictionaryPath
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

func NewRandomEventDataDictionary(fIlEnAmE string, pArSeR *config.ObjectParser) (*RandomEventDataDictionary, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRandomEventDataDictionary)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RandomEventDataDictionary{}

	return dAtA, nil
}

var vAlIdAtOrRandomEventDataDictionary = map[string]*config.Validator{}

func (dAtA *RandomEventDataDictionary) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with RandomEventPositionDictionary ----------------------------------

func LoadRandomEventPositionDictionary(gos *config.GameObjects) (*RandomEventPositionDictionary, *config.ObjectParser, error) {
	fIlEnAmE := confpath.RandomEventPositionDictionaryPath
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

	dAtA, err := NewRandomEventPositionDictionary(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedRandomEventPositionDictionary(gos *config.GameObjects, dAtA *RandomEventPositionDictionary, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RandomEventPositionDictionaryPath
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

func NewRandomEventPositionDictionary(fIlEnAmE string, pArSeR *config.ObjectParser) (*RandomEventPositionDictionary, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRandomEventPositionDictionary)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RandomEventPositionDictionary{}

	return dAtA, nil
}

var vAlIdAtOrRandomEventPositionDictionary = map[string]*config.Validator{}

func (dAtA *RandomEventPositionDictionary) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetCost(int) *resdata.Cost
	GetEventOptionData(uint64) *EventOptionData
	GetOptionPrize(uint64) *OptionPrize
	GetPrize(int) *resdata.Prize
}
