// AUTO_GEN, DONT MODIFY!!!
package fishing_data

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/goods"
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

// start with FishData ----------------------------------

func LoadFishData(gos *config.GameObjects) (map[uint64]*FishData, map[*FishData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FishDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*FishData, len(lIsT))
	pArSeRmAp := make(map[*FishData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFishData) {
			continue
		}

		dAtA, err := NewFishData(fIlEnAmE, pArSeR)
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

func SetRelatedFishData(dAtAmAp map[*FishData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FishDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFishDataKeyArray(datas []*FishData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewFishData(fIlEnAmE string, pArSeR *config.ObjectParser) (*FishData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFishData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FishData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Prize
	dAtA.IsShow = pArSeR.Bool("is_show")
	dAtA.Weight = pArSeR.Uint64("weight")
	dAtA.FishType = pArSeR.Uint64("fish_type")
	dAtA.FishTimes = 0
	if pArSeR.KeyExist("fish_times") {
		dAtA.FishTimes = pArSeR.Uint64("fish_times")
	}

	dAtA.IsPriority = false
	if pArSeR.KeyExist("is_priority") {
		dAtA.IsPriority = pArSeR.Bool("is_priority")
	}

	return dAtA, nil
}

var vAlIdAtOrFishData = map[string]*config.Validator{

	"id":          config.ParseValidator("int>0", "", false, nil, nil),
	"prize":       config.ParseValidator("string", "", false, nil, nil),
	"is_show":     config.ParseValidator("bool", "", false, nil, nil),
	"weight":      config.ParseValidator("int>0", "", false, nil, nil),
	"fish_type":   config.ParseValidator("int", "", false, nil, nil),
	"fish_times":  config.ParseValidator("int", "", false, nil, []string{"0"}),
	"is_priority": config.ParseValidator("bool", "", false, nil, []string{"false"}),
}

func (dAtA *FishData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	return nil
}

// start with FishRandomer ----------------------------------

func LoadFishRandomer(gos *config.GameObjects) (*FishRandomer, *config.ObjectParser, error) {
	fIlEnAmE := confpath.FishRandomerPath
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

	dAtA, err := NewFishRandomer(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedFishRandomer(gos *config.GameObjects, dAtA *FishRandomer, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FishRandomerPath
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

func NewFishRandomer(fIlEnAmE string, pArSeR *config.ObjectParser) (*FishRandomer, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFishRandomer)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FishRandomer{}

	return dAtA, nil
}

var vAlIdAtOrFishRandomer = map[string]*config.Validator{}

func (dAtA *FishRandomer) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with FishingCaptainProbabilityData ----------------------------------

func LoadFishingCaptainProbabilityData(gos *config.GameObjects) (map[uint64]*FishingCaptainProbabilityData, map[*FishingCaptainProbabilityData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FishingCaptainProbabilityDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*FishingCaptainProbabilityData, len(lIsT))
	pArSeRmAp := make(map[*FishingCaptainProbabilityData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFishingCaptainProbabilityData) {
			continue
		}

		dAtA, err := NewFishingCaptainProbabilityData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.CaptainId
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[CaptainId], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedFishingCaptainProbabilityData(dAtAmAp map[*FishingCaptainProbabilityData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FishingCaptainProbabilityDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFishingCaptainProbabilityDataKeyArray(datas []*FishingCaptainProbabilityData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.CaptainId)
		}
	}

	return out
}

func NewFishingCaptainProbabilityData(fIlEnAmE string, pArSeR *config.ObjectParser) (*FishingCaptainProbabilityData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFishingCaptainProbabilityData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FishingCaptainProbabilityData{}

	dAtA.CaptainId = pArSeR.Uint64("captain_id")
	dAtA.Multiple = pArSeR.Uint64("multiple")

	return dAtA, nil
}

var vAlIdAtOrFishingCaptainProbabilityData = map[string]*config.Validator{

	"captain_id": config.ParseValidator("uint", "", false, nil, nil),
	"multiple":   config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *FishingCaptainProbabilityData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *FishingCaptainProbabilityData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *FishingCaptainProbabilityData) Encode() *shared_proto.FishingCaptainProbabilityDataProto {
	out := &shared_proto.FishingCaptainProbabilityDataProto{}
	out.CaptainId = config.U64ToI32(dAtA.CaptainId)
	out.Multiple = config.U64ToI32(dAtA.Multiple)

	return out
}

func ArrayEncodeFishingCaptainProbabilityData(datas []*FishingCaptainProbabilityData) []*shared_proto.FishingCaptainProbabilityDataProto {

	out := make([]*shared_proto.FishingCaptainProbabilityDataProto, 0, len(datas))
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

func (dAtA *FishingCaptainProbabilityData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with FishingCostData ----------------------------------

func LoadFishingCostData(gos *config.GameObjects) (map[uint64]*FishingCostData, map[*FishingCostData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FishingCostDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*FishingCostData, len(lIsT))
	pArSeRmAp := make(map[*FishingCostData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFishingCostData) {
			continue
		}

		dAtA, err := NewFishingCostData(fIlEnAmE, pArSeR)
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

func SetRelatedFishingCostData(dAtAmAp map[*FishingCostData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FishingCostDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFishingCostDataKeyArray(datas []*FishingCostData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewFishingCostData(fIlEnAmE string, pArSeR *config.ObjectParser) (*FishingCostData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFishingCostData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FishingCostData{}

	dAtA.Times = pArSeR.Uint64("times")
	dAtA.FreeTimes = 0
	if pArSeR.KeyExist("free_times") {
		dAtA.FreeTimes = pArSeR.Uint64("free_times")
	}

	dAtA.DiscountTimes = 0
	if pArSeR.KeyExist("discount_times") {
		dAtA.DiscountTimes = pArSeR.Uint64("discount_times")
	}

	// releated field: DiscountCost
	// releated field: Cost
	dAtA.FishType = pArSeR.Uint64("fish_type")
	dAtA.DailyTimes = pArSeR.Uint64("daily_times")
	dAtA.FreeCountdown, err = config.ParseDuration(pArSeR.String("free_countdown"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[free_countdown] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("free_countdown"), dAtA)
	}

	// calculate fields
	dAtA.Id = FishId(dAtA.FishType, dAtA.Times)

	return dAtA, nil
}

var vAlIdAtOrFishingCostData = map[string]*config.Validator{

	"times":          config.ParseValidator("int>0", "", false, nil, nil),
	"free_times":     config.ParseValidator("int", "", false, nil, []string{"0"}),
	"discount_times": config.ParseValidator("int", "", false, nil, []string{"0"}),
	"discount_cost":  config.ParseValidator("string", "", false, nil, nil),
	"cost":           config.ParseValidator("string", "", false, nil, nil),
	"fish_type":      config.ParseValidator("int", "", false, nil, nil),
	"daily_times":    config.ParseValidator("int", "", false, nil, nil),
	"free_countdown": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *FishingCostData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *FishingCostData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *FishingCostData) Encode() *shared_proto.FishingCostProto {
	out := &shared_proto.FishingCostProto{}
	out.Times = config.U64ToI32(dAtA.Times)
	out.FreeTimes = config.U64ToI32(dAtA.FreeTimes)
	out.DiscountTimes = config.U64ToI32(dAtA.DiscountTimes)
	if dAtA.DiscountCost != nil {
		out.DiscountCost = dAtA.DiscountCost.Encode()
	}
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	out.FishType = config.U64ToI32(dAtA.FishType)
	out.DailyTimes = config.U64ToI32(dAtA.DailyTimes)

	return out
}

func ArrayEncodeFishingCostData(datas []*FishingCostData) []*shared_proto.FishingCostProto {

	out := make([]*shared_proto.FishingCostProto, 0, len(datas))
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

func (dAtA *FishingCostData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.DiscountCost = cOnFigS.GetCost(pArSeR.Int("discount_cost"))
	if dAtA.DiscountCost == nil {
		return errors.Errorf("%s 配置的关联字段[discount_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("discount_cost"), *pArSeR)
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	return nil
}

// start with FishingShowData ----------------------------------

func LoadFishingShowData(gos *config.GameObjects) (map[uint64]*FishingShowData, map[*FishingShowData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FishingShowDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*FishingShowData, len(lIsT))
	pArSeRmAp := make(map[*FishingShowData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFishingShowData) {
			continue
		}

		dAtA, err := NewFishingShowData(fIlEnAmE, pArSeR)
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

func SetRelatedFishingShowData(dAtAmAp map[*FishingShowData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FishingShowDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFishingShowDataKeyArray(datas []*FishingShowData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewFishingShowData(fIlEnAmE string, pArSeR *config.ObjectParser) (*FishingShowData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFishingShowData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FishingShowData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: GoodsData
	// releated field: GemData
	// releated field: EquipmentData
	// releated field: CaptainData
	dAtA.Desc = pArSeR.String("desc")
	dAtA.FishType = pArSeR.Uint64("fish_type")
	dAtA.Out = pArSeR.Bool("out")
	dAtA.ShowType = 1
	if pArSeR.KeyExist("show_type") {
		dAtA.ShowType = pArSeR.Uint64("show_type")
	}

	return dAtA, nil
}

var vAlIdAtOrFishingShowData = map[string]*config.Validator{

	"id":             config.ParseValidator("int>0", "", false, nil, nil),
	"goods_data":     config.ParseValidator("string", "", false, nil, nil),
	"gem_data":       config.ParseValidator("string", "", false, nil, nil),
	"equipment_data": config.ParseValidator("string", "", false, nil, nil),
	"captain_data":   config.ParseValidator("string", "", false, nil, nil),
	"desc":           config.ParseValidator("string", "", false, nil, nil),
	"fish_type":      config.ParseValidator("int", "", false, nil, nil),
	"out":            config.ParseValidator("bool", "", false, nil, nil),
	"show_type":      config.ParseValidator("int>0", "", false, nil, []string{"1"}),
}

func (dAtA *FishingShowData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *FishingShowData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *FishingShowData) Encode() *shared_proto.FishingShowProto {
	out := &shared_proto.FishingShowProto{}
	if dAtA.GoodsData != nil {
		out.GoodsData = config.U64ToI32(dAtA.GoodsData.Id)
	}
	if dAtA.GemData != nil {
		out.GemData = config.U64ToI32(dAtA.GemData.Id)
	}
	if dAtA.EquipmentData != nil {
		out.EquipmentData = config.U64ToI32(dAtA.EquipmentData.Id)
	}
	if dAtA.CaptainData != nil {
		out.CaptainSoulData = config.U64ToI32(dAtA.CaptainData.Id)
	}
	out.Desc = dAtA.Desc
	out.FishType = config.U64ToI32(dAtA.FishType)
	out.Out = dAtA.Out
	out.ShowType = config.U64ToI32(dAtA.ShowType)

	return out
}

func ArrayEncodeFishingShowData(datas []*FishingShowData) []*shared_proto.FishingShowProto {

	out := make([]*shared_proto.FishingShowProto, 0, len(datas))
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

func (dAtA *FishingShowData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.GoodsData = cOnFigS.GetGoodsData(pArSeR.Uint64("goods_data"))
	if dAtA.GoodsData == nil && pArSeR.Uint64("goods_data") != 0 {
		return errors.Errorf("%s 配置的关联字段[goods_data] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("goods_data"), *pArSeR)
	}

	dAtA.GemData = cOnFigS.GetGemData(pArSeR.Uint64("gem_data"))
	if dAtA.GemData == nil && pArSeR.Uint64("gem_data") != 0 {
		return errors.Errorf("%s 配置的关联字段[gem_data] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("gem_data"), *pArSeR)
	}

	dAtA.EquipmentData = cOnFigS.GetEquipmentData(pArSeR.Uint64("equipment_data"))
	if dAtA.EquipmentData == nil && pArSeR.Uint64("equipment_data") != 0 {
		return errors.Errorf("%s 配置的关联字段[equipment_data] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("equipment_data"), *pArSeR)
	}

	dAtA.CaptainData = cOnFigS.GetCaptainData(pArSeR.Uint64("captain_data"))
	if dAtA.CaptainData == nil && pArSeR.Uint64("captain_data") != 0 {
		return errors.Errorf("%s 配置的关联字段[captain_data] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain_data"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetCaptainData(uint64) *captain.CaptainData
	GetCost(int) *resdata.Cost
	GetEquipmentData(uint64) *goods.EquipmentData
	GetGemData(uint64) *goods.GemData
	GetGoodsData(uint64) *goods.GoodsData
	GetPrize(int) *resdata.Prize
}
