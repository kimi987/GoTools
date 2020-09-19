// AUTO_GEN, DONT MODIFY!!!
package charge

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

// start with ChargeObjData ----------------------------------

func LoadChargeObjData(gos *config.GameObjects) (map[uint64]*ChargeObjData, map[*ChargeObjData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ChargeObjDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ChargeObjData, len(lIsT))
	pArSeRmAp := make(map[*ChargeObjData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrChargeObjData) {
			continue
		}

		dAtA, err := NewChargeObjData(fIlEnAmE, pArSeR)
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

func SetRelatedChargeObjData(dAtAmAp map[*ChargeObjData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ChargeObjDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetChargeObjDataKeyArray(datas []*ChargeObjData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewChargeObjData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ChargeObjData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrChargeObjData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ChargeObjData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Image = pArSeR.String("image")
	// releated field: Product
	dAtA.ChargeAmount = pArSeR.Uint64("charge_amount")
	dAtA.Yuanbao = pArSeR.Uint64("yuanbao")
	dAtA.YuanbaoAddition = 0
	if pArSeR.KeyExist("yuanbao_addition") {
		dAtA.YuanbaoAddition = pArSeR.Uint64("yuanbao_addition")
	}

	dAtA.FirstChargeYuanbao = 0
	if pArSeR.KeyExist("first_charge_yuanbao") {
		dAtA.FirstChargeYuanbao = pArSeR.Uint64("first_charge_yuanbao")
	}

	dAtA.VipExp = pArSeR.Uint64("vip_exp")

	return dAtA, nil
}

var vAlIdAtOrChargeObjData = map[string]*config.Validator{

	"id":                   config.ParseValidator("int>0", "", false, nil, nil),
	"name":                 config.ParseValidator("string", "", false, nil, nil),
	"icon":                 config.ParseValidator("string", "", false, nil, nil),
	"image":                config.ParseValidator("string", "", false, nil, nil),
	"product":              config.ParseValidator("string", "", false, nil, nil),
	"charge_amount":        config.ParseValidator("int>0", "", false, nil, nil),
	"yuanbao":              config.ParseValidator("int>0", "", false, nil, nil),
	"yuanbao_addition":     config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"first_charge_yuanbao": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"vip_exp":              config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *ChargeObjData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ChargeObjData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ChargeObjData) Encode() *shared_proto.ChargeObjDataProto {
	out := &shared_proto.ChargeObjDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Icon = dAtA.Icon
	out.Image = dAtA.Image
	if dAtA.Product != nil {
		out.ProductId = int32(dAtA.Product.Id)
	}
	out.ChargeAmount = config.U64ToI32(dAtA.ChargeAmount)
	out.Yuanbao = config.U64ToI32(dAtA.Yuanbao)
	out.YuanbaoAddition = config.U64ToI32(dAtA.YuanbaoAddition)
	out.FirstChargeYuanbao = config.U64ToI32(dAtA.FirstChargeYuanbao)
	out.VipExp = config.U64ToI32(dAtA.VipExp)

	return out
}

func ArrayEncodeChargeObjData(datas []*ChargeObjData) []*shared_proto.ChargeObjDataProto {

	out := make([]*shared_proto.ChargeObjDataProto, 0, len(datas))
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

func (dAtA *ChargeObjData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Product = cOnFigS.GetProductData(pArSeR.Uint64("product"))
	if dAtA.Product == nil && pArSeR.Uint64("product") != 0 {
		return errors.Errorf("%s 配置的关联字段[product] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("product"), *pArSeR)
	}

	return nil
}

// start with ChargePrizeData ----------------------------------

func LoadChargePrizeData(gos *config.GameObjects) (map[uint64]*ChargePrizeData, map[*ChargePrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ChargePrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ChargePrizeData, len(lIsT))
	pArSeRmAp := make(map[*ChargePrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrChargePrizeData) {
			continue
		}

		dAtA, err := NewChargePrizeData(fIlEnAmE, pArSeR)
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

func SetRelatedChargePrizeData(dAtAmAp map[*ChargePrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ChargePrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetChargePrizeDataKeyArray(datas []*ChargePrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewChargePrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ChargePrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrChargePrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ChargePrizeData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Image = pArSeR.String("image")
	dAtA.Amount = pArSeR.Uint64("amount")
	dAtA.Value = pArSeR.Uint64("value")
	// releated field: Prize
	dAtA.Desc = pArSeR.String("desc")

	return dAtA, nil
}

var vAlIdAtOrChargePrizeData = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"name":   config.ParseValidator("string", "", false, nil, nil),
	"image":  config.ParseValidator("string", "", false, nil, nil),
	"amount": config.ParseValidator("int>0", "", false, nil, nil),
	"value":  config.ParseValidator("int>0", "", false, nil, nil),
	"prize":  config.ParseValidator("string", "", false, nil, nil),
	"desc":   config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ChargePrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ChargePrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ChargePrizeData) Encode() *shared_proto.ChargePrizeDataProto {
	out := &shared_proto.ChargePrizeDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Image = dAtA.Image
	out.Amount = config.U64ToI32(dAtA.Amount)
	out.Value = config.U64ToI32(dAtA.Value)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	out.Desc = dAtA.Desc

	return out
}

func ArrayEncodeChargePrizeData(datas []*ChargePrizeData) []*shared_proto.ChargePrizeDataProto {

	out := make([]*shared_proto.ChargePrizeDataProto, 0, len(datas))
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

func (dAtA *ChargePrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with ProductData ----------------------------------

func LoadProductData(gos *config.GameObjects) (map[uint64]*ProductData, map[*ProductData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ProductDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ProductData, len(lIsT))
	pArSeRmAp := make(map[*ProductData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrProductData) {
			continue
		}

		dAtA, err := NewProductData(fIlEnAmE, pArSeR)
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

func SetRelatedProductData(dAtAmAp map[*ProductData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ProductDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetProductDataKeyArray(datas []*ProductData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewProductData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ProductData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrProductData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ProductData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.ProductId = pArSeR.String("product_id")
	dAtA.ProductName = pArSeR.String("product_name")
	dAtA.Price = pArSeR.Uint64("price")

	return dAtA, nil
}

var vAlIdAtOrProductData = map[string]*config.Validator{

	"id":           config.ParseValidator("int>0", "", false, nil, nil),
	"product_id":   config.ParseValidator("string", "", false, nil, nil),
	"product_name": config.ParseValidator("string", "", false, nil, nil),
	"price":        config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *ProductData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetProductData(uint64) *ProductData
}
