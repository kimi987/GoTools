// AUTO_GEN, DONT MODIFY!!!
package combine

import (
	"github.com/lightpaw/config"
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

// start with EquipCombineData ----------------------------------

func LoadEquipCombineData(gos *config.GameObjects) (map[uint64]*EquipCombineData, map[*EquipCombineData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.EquipCombineDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*EquipCombineData, len(lIsT))
	pArSeRmAp := make(map[*EquipCombineData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrEquipCombineData) {
			continue
		}

		dAtA, err := NewEquipCombineData(fIlEnAmE, pArSeR)
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

func SetRelatedEquipCombineData(dAtAmAp map[*EquipCombineData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EquipCombineDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetEquipCombineDataKeyArray(datas []*EquipCombineData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewEquipCombineData(fIlEnAmE string, pArSeR *config.ObjectParser) (*EquipCombineData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEquipCombineData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EquipCombineData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.GroupName = pArSeR.String("group_name")
	// skip field: CostGoods
	// skip field: CombineEquip
	// releated field: CombineData

	return dAtA, nil
}

var vAlIdAtOrEquipCombineData = map[string]*config.Validator{

	"id":           config.ParseValidator("int>0", "", false, nil, nil),
	"group_name":   config.ParseValidator("string>0", "", false, nil, nil),
	"combine_data": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *EquipCombineData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *EquipCombineData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *EquipCombineData) Encode() *shared_proto.EquipCombineDataProto {
	out := &shared_proto.EquipCombineDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.GroupName = dAtA.GroupName
	if dAtA.CostGoods != nil {
		out.CostGoodsId = config.U64ToI32(dAtA.CostGoods.Id)
	}
	if dAtA.CombineEquip != nil {
		out.CombineEquipId = config.U64a2I32a(goods.GetEquipmentDataKeyArray(dAtA.CombineEquip))
	}
	if dAtA.CombineData != nil {
		out.CombineData = ArrayEncodeGoodsCombineData(dAtA.CombineData)
	}

	return out
}

func ArrayEncodeEquipCombineData(datas []*EquipCombineData) []*shared_proto.EquipCombineDataProto {

	out := make([]*shared_proto.EquipCombineDataProto, 0, len(datas))
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

func (dAtA *EquipCombineData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("combine_data", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetGoodsCombineData(v)
		if obj != nil {
			dAtA.CombineData = append(dAtA.CombineData, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[combine_data] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("combine_data"), *pArSeR)
		}
	}

	return nil
}

// start with EquipCombineDatas ----------------------------------

func LoadEquipCombineDatas(gos *config.GameObjects) (*EquipCombineDatas, *config.ObjectParser, error) {
	fIlEnAmE := confpath.EquipCombineDatasPath
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

	dAtA, err := NewEquipCombineDatas(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedEquipCombineDatas(gos *config.GameObjects, dAtA *EquipCombineDatas, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EquipCombineDatasPath
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

func NewEquipCombineDatas(fIlEnAmE string, pArSeR *config.ObjectParser) (*EquipCombineDatas, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEquipCombineDatas)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EquipCombineDatas{}

	return dAtA, nil
}

var vAlIdAtOrEquipCombineDatas = map[string]*config.Validator{}

func (dAtA *EquipCombineDatas) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GoodsCombineData ----------------------------------

func LoadGoodsCombineData(gos *config.GameObjects) (map[uint64]*GoodsCombineData, map[*GoodsCombineData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GoodsCombineDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GoodsCombineData, len(lIsT))
	pArSeRmAp := make(map[*GoodsCombineData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGoodsCombineData) {
			continue
		}

		dAtA, err := NewGoodsCombineData(fIlEnAmE, pArSeR)
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

func SetRelatedGoodsCombineData(dAtAmAp map[*GoodsCombineData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GoodsCombineDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGoodsCombineDataKeyArray(datas []*GoodsCombineData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGoodsCombineData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GoodsCombineData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGoodsCombineData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GoodsCombineData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Cost
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrGoodsCombineData = map[string]*config.Validator{

	"id":    config.ParseValidator("int>0", "", false, nil, nil),
	"cost":  config.ParseValidator("string", "", false, nil, nil),
	"prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *GoodsCombineData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GoodsCombineData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GoodsCombineData) Encode() *shared_proto.GoodsCombineDataProto {
	out := &shared_proto.GoodsCombineDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeGoodsCombineData(datas []*GoodsCombineData) []*shared_proto.GoodsCombineDataProto {

	out := make([]*shared_proto.GoodsCombineDataProto, 0, len(datas))
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

func (dAtA *GoodsCombineData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetCost(int) *resdata.Cost
	GetGoodsCombineData(uint64) *GoodsCombineData
	GetPrize(int) *resdata.Prize
}
