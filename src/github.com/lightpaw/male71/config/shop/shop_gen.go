// AUTO_GEN, DONT MODIFY!!!
package shop

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/guild_data"
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

// start with BlackMarketData ----------------------------------

func LoadBlackMarketData(gos *config.GameObjects) (map[uint64]*BlackMarketData, map[*BlackMarketData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BlackMarketDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BlackMarketData, len(lIsT))
	pArSeRmAp := make(map[*BlackMarketData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBlackMarketData) {
			continue
		}

		dAtA, err := NewBlackMarketData(fIlEnAmE, pArSeR)
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

func SetRelatedBlackMarketData(dAtAmAp map[*BlackMarketData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BlackMarketDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBlackMarketDataKeyArray(datas []*BlackMarketData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBlackMarketData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BlackMarketData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBlackMarketData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BlackMarketData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Group
	dAtA.MustDiscount = pArSeR.Float64Array("must_discount", "", false)
	dAtA.MustDiscountCount = pArSeR.Uint64Array("must_discount_count", "", false)
	dAtA.Discount = pArSeR.Float64Array("discount", "", false)
	dAtA.DiscountWeight = pArSeR.Uint64Array("discount_weight", "", false)

	return dAtA, nil
}

var vAlIdAtOrBlackMarketData = map[string]*config.Validator{

	"id":                  config.ParseValidator("int>0", "", false, nil, nil),
	"group":               config.ParseValidator("string", "", true, nil, nil),
	"must_discount":       config.ParseValidator("float64>0", "", true, nil, nil),
	"must_discount_count": config.ParseValidator(",duplicate", "", true, nil, nil),
	"discount":            config.ParseValidator("float64>0", "", true, nil, nil),
	"discount_weight":     config.ParseValidator(",duplicate", "", true, nil, nil),
}

func (dAtA *BlackMarketData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
		obj := cOnFigS.GetBlackMarketGoodsGroupData(v)
		if obj != nil {
			dAtA.Group = append(dAtA.Group, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[group] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("group"), *pArSeR)
		}
	}

	return nil
}

// start with BlackMarketGoodsData ----------------------------------

func LoadBlackMarketGoodsData(gos *config.GameObjects) (map[uint64]*BlackMarketGoodsData, map[*BlackMarketGoodsData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BlackMarketGoodsDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BlackMarketGoodsData, len(lIsT))
	pArSeRmAp := make(map[*BlackMarketGoodsData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBlackMarketGoodsData) {
			continue
		}

		dAtA, err := NewBlackMarketGoodsData(fIlEnAmE, pArSeR)
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

func SetRelatedBlackMarketGoodsData(dAtAmAp map[*BlackMarketGoodsData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BlackMarketGoodsDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBlackMarketGoodsDataKeyArray(datas []*BlackMarketGoodsData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBlackMarketGoodsData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BlackMarketGoodsData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBlackMarketGoodsData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BlackMarketGoodsData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Prize
	// releated field: ShowPrize
	// releated field: Cost
	dAtA.Quality = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("quality"))])
	if i, err := strconv.ParseInt(pArSeR.String("quality"), 10, 32); err == nil {
		dAtA.Quality = shared_proto.Quality(i)
	}

	dAtA.RequiredHeroLevel = 0
	if pArSeR.KeyExist("required_hero_level") {
		dAtA.RequiredHeroLevel = pArSeR.Uint64("required_hero_level")
	}

	return dAtA, nil
}

var vAlIdAtOrBlackMarketGoodsData = map[string]*config.Validator{

	"id":                  config.ParseValidator("int>0", "", false, nil, nil),
	"prize":               config.ParseValidator("string", "", false, nil, nil),
	"show_prize":          config.ParseValidator("string", "", false, nil, nil),
	"cost":                config.ParseValidator("string", "", false, nil, nil),
	"quality":             config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), nil),
	"required_hero_level": config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *BlackMarketGoodsData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BlackMarketGoodsData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BlackMarketGoodsData) Encode() *shared_proto.BlackMarketGoodsDataProto {
	out := &shared_proto.BlackMarketGoodsDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.Encode()
	}
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	out.Quality = dAtA.Quality

	return out
}

func ArrayEncodeBlackMarketGoodsData(datas []*BlackMarketGoodsData) []*shared_proto.BlackMarketGoodsDataProto {

	out := make([]*shared_proto.BlackMarketGoodsDataProto, 0, len(datas))
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

func (dAtA *BlackMarketGoodsData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil {
		return errors.Errorf("%s 配置的关联字段[show_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	return nil
}

// start with BlackMarketGoodsGroupData ----------------------------------

func LoadBlackMarketGoodsGroupData(gos *config.GameObjects) (map[uint64]*BlackMarketGoodsGroupData, map[*BlackMarketGoodsGroupData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BlackMarketGoodsGroupDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BlackMarketGoodsGroupData, len(lIsT))
	pArSeRmAp := make(map[*BlackMarketGoodsGroupData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBlackMarketGoodsGroupData) {
			continue
		}

		dAtA, err := NewBlackMarketGoodsGroupData(fIlEnAmE, pArSeR)
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

func SetRelatedBlackMarketGoodsGroupData(dAtAmAp map[*BlackMarketGoodsGroupData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BlackMarketGoodsGroupDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBlackMarketGoodsGroupDataKeyArray(datas []*BlackMarketGoodsGroupData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBlackMarketGoodsGroupData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BlackMarketGoodsGroupData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBlackMarketGoodsGroupData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BlackMarketGoodsGroupData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Count = pArSeR.Uint64("count")
	// releated field: Goods
	dAtA.Weight = pArSeR.Uint64Array("weight", "", false)

	return dAtA, nil
}

var vAlIdAtOrBlackMarketGoodsGroupData = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"count":  config.ParseValidator("int>0", "", false, nil, nil),
	"goods":  config.ParseValidator("string", "", true, nil, nil),
	"weight": config.ParseValidator(",duplicate", "", true, nil, nil),
}

func (dAtA *BlackMarketGoodsGroupData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("goods", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetBlackMarketGoodsData(v)
		if obj != nil {
			dAtA.Goods = append(dAtA.Goods, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("goods"), *pArSeR)
		}
	}

	return nil
}

// start with DiscountColorData ----------------------------------

func LoadDiscountColorData(gos *config.GameObjects) (map[uint64]*DiscountColorData, map[*DiscountColorData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.DiscountColorDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*DiscountColorData, len(lIsT))
	pArSeRmAp := make(map[*DiscountColorData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrDiscountColorData) {
			continue
		}

		dAtA, err := NewDiscountColorData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Discount
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Discount], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedDiscountColorData(dAtAmAp map[*DiscountColorData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.DiscountColorDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetDiscountColorDataKeyArray(datas []*DiscountColorData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Discount)
		}
	}

	return out
}

func NewDiscountColorData(fIlEnAmE string, pArSeR *config.ObjectParser) (*DiscountColorData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrDiscountColorData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &DiscountColorData{}

	dAtA.Discount = pArSeR.Uint64("discount")
	dAtA.Color = pArSeR.String("color")

	return dAtA, nil
}

var vAlIdAtOrDiscountColorData = map[string]*config.Validator{

	"discount": config.ParseValidator("int>0", "", false, nil, nil),
	"color":    config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *DiscountColorData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *DiscountColorData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *DiscountColorData) Encode() *shared_proto.DiscountColorDataProto {
	out := &shared_proto.DiscountColorDataProto{}
	out.Discount = config.U64ToI32(dAtA.Discount)
	out.Color = dAtA.Color

	return out
}

func ArrayEncodeDiscountColorData(datas []*DiscountColorData) []*shared_proto.DiscountColorDataProto {

	out := make([]*shared_proto.DiscountColorDataProto, 0, len(datas))
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

func (dAtA *DiscountColorData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with NormalShopGoods ----------------------------------

func LoadNormalShopGoods(gos *config.GameObjects) (map[uint64]*NormalShopGoods, map[*NormalShopGoods]*config.ObjectParser, error) {
	fIlEnAmE := confpath.NormalShopGoodsPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*NormalShopGoods, len(lIsT))
	pArSeRmAp := make(map[*NormalShopGoods]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrNormalShopGoods) {
			continue
		}

		dAtA, err := NewNormalShopGoods(fIlEnAmE, pArSeR)
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

func SetRelatedNormalShopGoods(dAtAmAp map[*NormalShopGoods]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.NormalShopGoodsPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetNormalShopGoodsKeyArray(datas []*NormalShopGoods) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewNormalShopGoods(fIlEnAmE string, pArSeR *config.ObjectParser) (*NormalShopGoods, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrNormalShopGoods)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &NormalShopGoods{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.CountLimit = pArSeR.Uint64("count_limit")
	dAtA.UnlockCondition, err = data.NewUnlockCondition(fIlEnAmE, pArSeR)
	if err != nil {
		return nil, err
	}
	dAtA.FreeTimes = pArSeR.Uint64("free_times")
	// releated field: UseImmediatelyGoods
	// skip field: ShopType
	// releated field: Cost
	// releated field: ShowPrize
	// releated field: PlunderPrize
	dAtA.Tag = pArSeR.String("tag")
	dAtA.ShowSale = 0
	if pArSeR.KeyExist("show_sale") {
		dAtA.ShowSale = pArSeR.Uint64("show_sale")
	}

	// releated field: ShowOriginCost
	dAtA.VipDailyMaxCount = 0
	if pArSeR.KeyExist("vip_daily_max_count") {
		dAtA.VipDailyMaxCount = pArSeR.Uint64("vip_daily_max_count")
	}

	// releated field: GuildEventPrize

	return dAtA, nil
}

var vAlIdAtOrNormalShopGoods = map[string]*config.Validator{

	"id":                    config.ParseValidator("int>0", "", false, nil, nil),
	"count_limit":           config.ParseValidator("uint", "", false, nil, nil),
	"free_times":            config.ParseValidator("uint", "", false, nil, nil),
	"use_immediately_goods": config.ParseValidator("string", "", true, nil, nil),
	"cost":                  config.ParseValidator("string", "", false, nil, nil),
	"show_prize":            config.ParseValidator("string", "", false, nil, nil),
	"plunder_prize":         config.ParseValidator("string", "", false, nil, nil),
	"tag":                   config.ParseValidator("string", "", false, nil, nil),
	"show_sale":             config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"show_origin_cost":      config.ParseValidator("string", "", false, nil, nil),
	"vip_daily_max_count":   config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"guild_event_prize":     config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *NormalShopGoods) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *NormalShopGoods) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *NormalShopGoods) Encode() *shared_proto.NormalShopGoodsProto {
	out := &shared_proto.NormalShopGoodsProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.CountLimit = config.U64ToI32(dAtA.CountLimit)
	if dAtA.UnlockCondition != nil {
		out.UnlockCondition = dAtA.UnlockCondition.Encode()
	}
	out.FreeTimes = config.U64ToI32(dAtA.FreeTimes)
	if dAtA.UseImmediatelyGoods != nil {
		out.UseImmediatelyGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.UseImmediatelyGoods))
	}
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.Encode()
	}
	if dAtA.PlunderPrize != nil {
		out.Prize = dAtA.PlunderPrize.Prize.PrizeProto()
	}
	out.Tag = shared_proto.ShopGoodsTag_value[dAtA.Tag]
	out.ShowSale = config.U64ToI32(dAtA.ShowSale)
	if dAtA.ShowOriginCost != nil {
		out.ShowOriginCost = dAtA.ShowOriginCost.Encode()
	}
	out.VipDailyMaxCount = config.U64ToI32(dAtA.VipDailyMaxCount)
	if dAtA.GuildEventPrize != nil {
		out.GuildEventPrize = config.U64ToI32(dAtA.GuildEventPrize.Id)
	}

	return out
}

func ArrayEncodeNormalShopGoods(datas []*NormalShopGoods) []*shared_proto.NormalShopGoodsProto {

	out := make([]*shared_proto.NormalShopGoodsProto, 0, len(datas))
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

func (dAtA *NormalShopGoods) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if err := dAtA.UnlockCondition.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS0); err != nil {
		return err
	}

	uint64Keys = pArSeR.Uint64Array("use_immediately_goods", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetGoodsData(v)
		if obj != nil {
			dAtA.UseImmediatelyGoods = append(dAtA.UseImmediatelyGoods, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[use_immediately_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("use_immediately_goods"), *pArSeR)
		}
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil {
		return errors.Errorf("%s 配置的关联字段[show_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
	}

	dAtA.PlunderPrize = cOnFigS.GetPlunderPrize(pArSeR.Uint64("plunder_prize"))
	if dAtA.PlunderPrize == nil {
		return errors.Errorf("%s 配置的关联字段[plunder_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder_prize"), *pArSeR)
	}

	dAtA.ShowOriginCost = cOnFigS.GetCost(pArSeR.Int("show_origin_cost"))
	if dAtA.ShowOriginCost == nil && pArSeR.Int("show_origin_cost") != 0 {
		return errors.Errorf("%s 配置的关联字段[show_origin_cost] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_origin_cost"), *pArSeR)
	}

	dAtA.GuildEventPrize = cOnFigS.GetGuildEventPrizeData(pArSeR.Uint64("guild_event_prize"))
	if dAtA.GuildEventPrize == nil && pArSeR.Uint64("guild_event_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[guild_event_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_event_prize"), *pArSeR)
	}

	return nil
}

// start with Shop ----------------------------------

func LoadShop(gos *config.GameObjects) (map[uint64]*Shop, map[*Shop]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ShopPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*Shop, len(lIsT))
	pArSeRmAp := make(map[*Shop]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrShop) {
			continue
		}

		dAtA, err := NewShop(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Type
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Type], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedShop(dAtAmAp map[*Shop]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ShopPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetShopKeyArray(datas []*Shop) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Type)
		}
	}

	return out
}

func NewShop(fIlEnAmE string, pArSeR *config.ObjectParser) (*Shop, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrShop)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &Shop{}

	dAtA.Type = pArSeR.Uint64("type")
	// releated field: NormalGoods
	// releated field: ZhenBaoGeGoods

	return dAtA, nil
}

var vAlIdAtOrShop = map[string]*config.Validator{

	"type":              config.ParseValidator("int>0", "", false, nil, nil),
	"normal_goods":      config.ParseValidator("string", "", true, nil, nil),
	"zhen_bao_ge_goods": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *Shop) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *Shop) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *Shop) Encode() *shared_proto.ShopProto {
	out := &shared_proto.ShopProto{}
	out.Type = config.U64ToI32(dAtA.Type)
	if dAtA.NormalGoods != nil {
		out.NormalGoods = ArrayEncodeNormalShopGoods(dAtA.NormalGoods)
	}
	if dAtA.ZhenBaoGeGoods != nil {
		out.ZhenBaoGeGoods = ArrayEncodeZhenBaoGeShopGoods(dAtA.ZhenBaoGeGoods)
	}

	return out
}

func ArrayEncodeShop(datas []*Shop) []*shared_proto.ShopProto {

	out := make([]*shared_proto.ShopProto, 0, len(datas))
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

func (dAtA *Shop) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("normal_goods", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetNormalShopGoods(v)
		if obj != nil {
			dAtA.NormalGoods = append(dAtA.NormalGoods, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[normal_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("normal_goods"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("zhen_bao_ge_goods", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetZhenBaoGeShopGoods(v)
		if obj != nil {
			dAtA.ZhenBaoGeGoods = append(dAtA.ZhenBaoGeGoods, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[zhen_bao_ge_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("zhen_bao_ge_goods"), *pArSeR)
		}
	}

	return nil
}

// start with ShopMiscData ----------------------------------

func LoadShopMiscData(gos *config.GameObjects) (*ShopMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.ShopMiscDataPath
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

	dAtA, err := NewShopMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedShopMiscData(gos *config.GameObjects, dAtA *ShopMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ShopMiscDataPath
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

func NewShopMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ShopMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrShopMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ShopMiscData{}

	// releated field: RefreshBlackMarketCost
	stringKeys = pArSeR.StringArray("auto_refresh_black_market_duration", "", false)
	dAtA.AutoRefreshBlackMarketDuration = make([]time.Duration, 0, len(stringKeys))
	for _, v := range stringKeys {
		obj, err := config.ParseDuration(v)
		if err != nil {
			return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[auto_refresh_black_market_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("auto_refresh_black_market_duration"), dAtA)
		}
		dAtA.AutoRefreshBlackMarketDuration = append(dAtA.AutoRefreshBlackMarketDuration, obj)
	}

	return dAtA, nil
}

var vAlIdAtOrShopMiscData = map[string]*config.Validator{

	"refresh_black_market_cost":          config.ParseValidator(",duplicate", "", true, nil, nil),
	"auto_refresh_black_market_duration": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *ShopMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ShopMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ShopMiscData) Encode() *shared_proto.ShopMiscDataProto {
	out := &shared_proto.ShopMiscDataProto{}
	if dAtA.RefreshBlackMarketCost != nil {
		out.RefreshBlackMarketCost = resdata.ArrayEncodeCost(dAtA.RefreshBlackMarketCost)
	}
	out.AutoRefreshBlackMarketDuration = config.DurationArr2I32Seconds(dAtA.AutoRefreshBlackMarketDuration)

	return out
}

func ArrayEncodeShopMiscData(datas []*ShopMiscData) []*shared_proto.ShopMiscDataProto {

	out := make([]*shared_proto.ShopMiscDataProto, 0, len(datas))
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

func (dAtA *ShopMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	intKeys = pArSeR.IntArray("refresh_black_market_cost", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetCost(v)
		if obj != nil {
			dAtA.RefreshBlackMarketCost = append(dAtA.RefreshBlackMarketCost, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[refresh_black_market_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("refresh_black_market_cost"), *pArSeR)
		}
	}

	return nil
}

// start with ZhenBaoGeShopGoods ----------------------------------

func LoadZhenBaoGeShopGoods(gos *config.GameObjects) (map[uint64]*ZhenBaoGeShopGoods, map[*ZhenBaoGeShopGoods]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ZhenBaoGeShopGoodsPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ZhenBaoGeShopGoods, len(lIsT))
	pArSeRmAp := make(map[*ZhenBaoGeShopGoods]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrZhenBaoGeShopGoods) {
			continue
		}

		dAtA, err := NewZhenBaoGeShopGoods(fIlEnAmE, pArSeR)
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

func SetRelatedZhenBaoGeShopGoods(dAtAmAp map[*ZhenBaoGeShopGoods]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ZhenBaoGeShopGoodsPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetZhenBaoGeShopGoodsKeyArray(datas []*ZhenBaoGeShopGoods) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewZhenBaoGeShopGoods(fIlEnAmE string, pArSeR *config.ObjectParser) (*ZhenBaoGeShopGoods, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrZhenBaoGeShopGoods)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ZhenBaoGeShopGoods{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.CountLimit = pArSeR.Uint64("count_limit")
	dAtA.UnlockCondition, err = data.NewUnlockCondition(fIlEnAmE, pArSeR)
	if err != nil {
		return nil, err
	}
	dAtA.FreeTimes = pArSeR.Uint64("free_times")
	// releated field: UseImmediatelyGoods
	// skip field: ShopType
	// releated field: BuyCosts
	dAtA.MaxBuyCountPerTimes = 0
	if pArSeR.KeyExist("max_buy_count_per_times") {
		dAtA.MaxBuyCountPerTimes = pArSeR.Uint64("max_buy_count_per_times")
	}

	// releated field: ShowPrizes
	// releated field: Prizes
	// releated field: Plunders
	dAtA.Levels = pArSeR.Uint64Array("levels", "", false)
	dAtA.CritWeight = pArSeR.Uint64Array("crit_weight", "", false)
	dAtA.CritMulti = pArSeR.Uint64Array("crit_multi", "", false)
	dAtA.BroadcastMinMulti = pArSeR.Uint64("broadcast_min_multi")
	dAtA.Tag = pArSeR.String("tag")

	return dAtA, nil
}

var vAlIdAtOrZhenBaoGeShopGoods = map[string]*config.Validator{

	"id":                      config.ParseValidator("int>0", "", false, nil, nil),
	"count_limit":             config.ParseValidator("uint", "", false, nil, nil),
	"free_times":              config.ParseValidator("uint", "", false, nil, nil),
	"use_immediately_goods":   config.ParseValidator("string", "", true, nil, nil),
	"buy_costs":               config.ParseValidator("string", "", true, nil, nil),
	"max_buy_count_per_times": config.ParseValidator("int", "", false, nil, []string{"0"}),
	"show_prizes":             config.ParseValidator("string", "", true, nil, nil),
	"prizes":                  config.ParseValidator("string", "", true, nil, nil),
	"plunders":                config.ParseValidator("string", "", true, nil, nil),
	"levels":                  config.ParseValidator("int", "", true, nil, nil),
	"crit_weight":             config.ParseValidator("int", "", true, nil, nil),
	"crit_multi":              config.ParseValidator("int", "", true, nil, nil),
	"broadcast_min_multi":     config.ParseValidator("uint", "", false, nil, nil),
	"tag": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ZhenBaoGeShopGoods) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ZhenBaoGeShopGoods) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ZhenBaoGeShopGoods) Encode() *shared_proto.ZhenBaoGeShopGoodsProto {
	out := &shared_proto.ZhenBaoGeShopGoodsProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.CountLimit = config.U64ToI32(dAtA.CountLimit)
	if dAtA.UnlockCondition != nil {
		out.UnlockCondition = dAtA.UnlockCondition.Encode()
	}
	out.FreeTimes = config.U64ToI32(dAtA.FreeTimes)
	if dAtA.UseImmediatelyGoods != nil {
		out.UseImmediatelyGoods = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.UseImmediatelyGoods))
	}
	if dAtA.BuyCosts != nil {
		out.BuyCosts = resdata.ArrayEncodeCost(dAtA.BuyCosts)
	}
	if dAtA.ShowPrizes != nil {
		out.Prizes = resdata.ArrayEncodePrize(dAtA.ShowPrizes)
	}
	out.Levels = config.U64a2I32a(dAtA.Levels)
	out.Tag = shared_proto.ShopGoodsTag_value[dAtA.Tag]

	return out
}

func ArrayEncodeZhenBaoGeShopGoods(datas []*ZhenBaoGeShopGoods) []*shared_proto.ZhenBaoGeShopGoodsProto {

	out := make([]*shared_proto.ZhenBaoGeShopGoodsProto, 0, len(datas))
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

func (dAtA *ZhenBaoGeShopGoods) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if err := dAtA.UnlockCondition.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS0); err != nil {
		return err
	}

	uint64Keys = pArSeR.Uint64Array("use_immediately_goods", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetGoodsData(v)
		if obj != nil {
			dAtA.UseImmediatelyGoods = append(dAtA.UseImmediatelyGoods, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[use_immediately_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("use_immediately_goods"), *pArSeR)
		}
	}

	intKeys = pArSeR.IntArray("buy_costs", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetCost(v)
		if obj != nil {
			dAtA.BuyCosts = append(dAtA.BuyCosts, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[buy_costs] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("buy_costs"), *pArSeR)
		}
	}

	intKeys = pArSeR.IntArray("show_prizes", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetPrize(v)
		if obj != nil {
			dAtA.ShowPrizes = append(dAtA.ShowPrizes, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[show_prizes] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prizes"), *pArSeR)
		}
	}

	intKeys = pArSeR.IntArray("prizes", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetPrize(v)
		if obj != nil {
			dAtA.Prizes = append(dAtA.Prizes, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[prizes] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prizes"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("plunders", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetPlunder(v)
		if obj != nil {
			dAtA.Plunders = append(dAtA.Plunders, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[plunders] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunders"), *pArSeR)
		}
	}

	return nil
}

type related_configs interface {
	GetBlackMarketGoodsData(uint64) *BlackMarketGoodsData
	GetBlackMarketGoodsGroupData(uint64) *BlackMarketGoodsGroupData
	GetCost(int) *resdata.Cost
	GetGoodsData(uint64) *goods.GoodsData
	GetGuildEventPrizeData(uint64) *guild_data.GuildEventPrizeData
	GetNormalShopGoods(uint64) *NormalShopGoods
	GetPlunder(uint64) *resdata.Plunder
	GetPlunderPrize(uint64) *resdata.PlunderPrize
	GetPrize(int) *resdata.Prize
	GetZhenBaoGeShopGoods(uint64) *ZhenBaoGeShopGoods
}
