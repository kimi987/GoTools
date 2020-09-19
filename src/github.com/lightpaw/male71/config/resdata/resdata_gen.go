// AUTO_GEN, DONT MODIFY!!!
package resdata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
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

// start with AmountShowSortData ----------------------------------

func LoadAmountShowSortData(gos *config.GameObjects) (map[uint64]*AmountShowSortData, map[*AmountShowSortData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.AmountShowSortDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*AmountShowSortData, len(lIsT))
	pArSeRmAp := make(map[*AmountShowSortData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrAmountShowSortData) {
			continue
		}

		dAtA, err := NewAmountShowSortData(fIlEnAmE, pArSeR)
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

func SetRelatedAmountShowSortData(dAtAmAp map[*AmountShowSortData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.AmountShowSortDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetAmountShowSortDataKeyArray(datas []*AmountShowSortData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewAmountShowSortData(fIlEnAmE string, pArSeR *config.ObjectParser) (*AmountShowSortData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrAmountShowSortData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &AmountShowSortData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// skip field: TypeList

	return dAtA, nil
}

var vAlIdAtOrAmountShowSortData = map[string]*config.Validator{

	"id":   config.ParseValidator("int>0", "", false, nil, nil),
	"name": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *AmountShowSortData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *AmountShowSortData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *AmountShowSortData) Encode() *shared_proto.AmountShowSortProto {
	out := &shared_proto.AmountShowSortProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.TypeList = dAtA.TypeList

	return out
}

func ArrayEncodeAmountShowSortData(datas []*AmountShowSortData) []*shared_proto.AmountShowSortProto {

	out := make([]*shared_proto.AmountShowSortProto, 0, len(datas))
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

func (dAtA *AmountShowSortData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with BaowuData ----------------------------------

func LoadBaowuData(gos *config.GameObjects) (map[uint64]*BaowuData, map[*BaowuData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BaowuDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BaowuData, len(lIsT))
	pArSeRmAp := make(map[*BaowuData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBaowuData) {
			continue
		}

		dAtA, err := NewBaowuData(fIlEnAmE, pArSeR)
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

func SetRelatedBaowuData(dAtAmAp map[*BaowuData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BaowuDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBaowuDataKeyArray(datas []*BaowuData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBaowuData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BaowuData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBaowuData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BaowuData{}

	dAtA.Group = pArSeR.Uint64("group")
	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	// releated field: Icon
	// skip field: Quality
	// releated field: GoodsQuality
	dAtA.UnlockDuration, err = config.ParseDuration(pArSeR.String("unlock_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[unlock_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("unlock_duration"), dAtA)
	}

	// releated field: PlunderPrize
	dAtA.UpgradeNeedCount = pArSeR.Uint64("upgrade_need_count")
	dAtA.DecomposeGold = pArSeR.Uint64("decompose_gold")
	dAtA.DecomposeStone = pArSeR.Uint64("decompose_stone")
	dAtA.MiaoDuration, err = config.ParseDuration(pArSeR.String("miao_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[miao_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("miao_duration"), dAtA)
	}

	dAtA.CantRob = false
	if pArSeR.KeyExist("cant_rob") {
		dAtA.CantRob = pArSeR.Bool("cant_rob")
	}

	dAtA.Prestige = pArSeR.Uint64("prestige")

	// calculate fields
	dAtA.Id = BaoDataId(dAtA.Group, dAtA.Level)

	return dAtA, nil
}

var vAlIdAtOrBaowuData = map[string]*config.Validator{

	"group":              config.ParseValidator("int", "", false, nil, nil),
	"level":              config.ParseValidator("int>0", "", false, nil, nil),
	"name":               config.ParseValidator("string", "", false, nil, nil),
	"desc":               config.ParseValidator("string", "", false, nil, nil),
	"icon":               config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"goods_quality":      config.ParseValidator("string", "", false, nil, nil),
	"unlock_duration":    config.ParseValidator("string", "", false, nil, nil),
	"plunder_prize":      config.ParseValidator("string", "", false, nil, nil),
	"upgrade_need_count": config.ParseValidator("int>0", "", false, nil, nil),
	"decompose_gold":     config.ParseValidator("uint", "", false, nil, nil),
	"decompose_stone":    config.ParseValidator("uint", "", false, nil, nil),
	"miao_duration":      config.ParseValidator("string", "", false, nil, nil),
	"cant_rob":           config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"prestige":           config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *BaowuData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BaowuData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BaowuData) Encode() *shared_proto.BaowuDataProto {
	out := &shared_proto.BaowuDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Group = config.U64ToI32(dAtA.Group)
	out.Level = config.U64ToI32(dAtA.Level)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	if dAtA.Icon != nil {
		out.Icon = dAtA.Icon.Id
	}
	out.Quality = dAtA.Quality
	if dAtA.GoodsQuality != nil {
		out.GoodsQuality = config.U64ToI32(dAtA.GoodsQuality.Level)
	}
	out.UnlockDuration = config.Duration2I32Seconds(dAtA.UnlockDuration)
	if dAtA.PlunderPrize != nil {
		out.PlunderPrize = dAtA.PlunderPrize.Prize.PrizeProto()
	}
	out.UpgradeNeedCount = config.U64ToI32(dAtA.UpgradeNeedCount)
	out.DecomposeGold = config.U64ToI32(dAtA.DecomposeGold)
	out.DecomposeStone = config.U64ToI32(dAtA.DecomposeStone)
	out.MiaoDuration = config.Duration2I32Seconds(dAtA.MiaoDuration)
	out.Prestige = config.U64ToI32(dAtA.Prestige)

	return out
}

func ArrayEncodeBaowuData(datas []*BaowuData) []*shared_proto.BaowuDataProto {

	out := make([]*shared_proto.BaowuDataProto, 0, len(datas))
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

func (dAtA *BaowuData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.PlunderPrize = cOnFigS.GetPlunderPrize(pArSeR.Uint64("plunder_prize"))
	if dAtA.PlunderPrize == nil {
		return errors.Errorf("%s 配置的关联字段[plunder_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder_prize"), *pArSeR)
	}

	return nil
}

// start with ConditionPlunder ----------------------------------

func LoadConditionPlunder(gos *config.GameObjects) (map[uint64]*ConditionPlunder, map[*ConditionPlunder]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ConditionPlunderPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ConditionPlunder, len(lIsT))
	pArSeRmAp := make(map[*ConditionPlunder]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrConditionPlunder) {
			continue
		}

		dAtA, err := NewConditionPlunder(fIlEnAmE, pArSeR)
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

func SetRelatedConditionPlunder(dAtAmAp map[*ConditionPlunder]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ConditionPlunderPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetConditionPlunderKeyArray(datas []*ConditionPlunder) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewConditionPlunder(fIlEnAmE string, pArSeR *config.ObjectParser) (*ConditionPlunder, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrConditionPlunder)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ConditionPlunder{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: DefPlunder
	// releated field: CondItem

	return dAtA, nil
}

var vAlIdAtOrConditionPlunder = map[string]*config.Validator{

	"id":          config.ParseValidator("int>0", "", false, nil, nil),
	"def_plunder": config.ParseValidator("string", "", false, nil, nil),
	"cond_item":   config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *ConditionPlunder) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.DefPlunder = cOnFigS.GetPlunder(pArSeR.Uint64("def_plunder"))
	if dAtA.DefPlunder == nil {
		return errors.Errorf("%s 配置的关联字段[def_plunder] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("def_plunder"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("cond_item", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetConditionPlunderItem(v)
		if obj != nil {
			dAtA.CondItem = append(dAtA.CondItem, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[cond_item] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cond_item"), *pArSeR)
		}
	}

	return nil
}

// start with ConditionPlunderItem ----------------------------------

func LoadConditionPlunderItem(gos *config.GameObjects) (map[uint64]*ConditionPlunderItem, map[*ConditionPlunderItem]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ConditionPlunderItemPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ConditionPlunderItem, len(lIsT))
	pArSeRmAp := make(map[*ConditionPlunderItem]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrConditionPlunderItem) {
			continue
		}

		dAtA, err := NewConditionPlunderItem(fIlEnAmE, pArSeR)
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

func SetRelatedConditionPlunderItem(dAtAmAp map[*ConditionPlunderItem]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ConditionPlunderItemPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetConditionPlunderItemKeyArray(datas []*ConditionPlunderItem) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewConditionPlunderItem(fIlEnAmE string, pArSeR *config.ObjectParser) (*ConditionPlunderItem, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrConditionPlunderItem)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ConditionPlunderItem{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.HeroLevel, err = data.ParseCompareCondition(pArSeR.String("hero_level"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[hero_level] 解析失败(data.ParseCompareCondition)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("hero_level"), dAtA)
	}

	// releated field: Plunder

	return dAtA, nil
}

var vAlIdAtOrConditionPlunderItem = map[string]*config.Validator{

	"id":         config.ParseValidator("int>0", "", false, nil, nil),
	"hero_level": config.ParseValidator("string", "", false, nil, nil),
	"plunder":    config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ConditionPlunderItem) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Plunder = cOnFigS.GetPlunder(pArSeR.Uint64("plunder"))
	if dAtA.Plunder == nil {
		return errors.Errorf("%s 配置的关联字段[plunder] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder"), *pArSeR)
	}

	return nil
}

// start with Cost ----------------------------------

func LoadCost(gos *config.GameObjects) (map[int]*Cost, map[*Cost]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CostPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[int]*Cost, len(lIsT))
	pArSeRmAp := make(map[*Cost]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCost) {
			continue
		}

		dAtA, err := NewCost(fIlEnAmE, pArSeR)
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

func SetRelatedCost(dAtAmAp map[*Cost]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CostPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCostKeyArray(datas []*Cost) []int {

	out := make([]int, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCost(fIlEnAmE string, pArSeR *config.ObjectParser) (*Cost, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCost)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &Cost{}

	dAtA.Id = pArSeR.Int("id")
	dAtA.Gold = pArSeR.Uint64("gold")
	dAtA.Food = pArSeR.Uint64("food")
	dAtA.Wood = pArSeR.Uint64("wood")
	dAtA.Stone = pArSeR.Uint64("stone")
	dAtA.Yuanbao = pArSeR.Uint64("yuanbao")
	dAtA.Dianquan = pArSeR.Uint64("dianquan")
	dAtA.Yinliang = pArSeR.Uint64("yinliang")
	dAtA.GuildContributionCoin = pArSeR.Uint64("guild_contribution_coin")
	dAtA.Jade = pArSeR.Uint64("jade")
	dAtA.JadeOre = pArSeR.Uint64("jade_ore")
	// releated field: Goods
	dAtA.GoodsCount = pArSeR.Uint64Array("goods_count", "", false)
	// releated field: Gem
	dAtA.GemCount = pArSeR.Uint64Array("gem_count", "", false)
	// skip field: IsNotEmpty
	// skip field: IsOnlyResource

	return dAtA, nil
}

var vAlIdAtOrCost = map[string]*config.Validator{

	"id":                      config.ParseValidator("int>0", "", false, nil, nil),
	"gold":                    config.ParseValidator("uint", "", false, nil, nil),
	"food":                    config.ParseValidator("uint", "", false, nil, nil),
	"wood":                    config.ParseValidator("uint", "", false, nil, nil),
	"stone":                   config.ParseValidator("uint", "", false, nil, nil),
	"yuanbao":                 config.ParseValidator("uint", "", false, nil, nil),
	"dianquan":                config.ParseValidator("uint", "", false, nil, nil),
	"yinliang":                config.ParseValidator("uint", "", false, nil, nil),
	"guild_contribution_coin": config.ParseValidator("uint", "", false, nil, nil),
	"jade":        config.ParseValidator("uint", "", false, nil, nil),
	"jade_ore":    config.ParseValidator("uint", "", false, nil, nil),
	"goods_id":    config.ParseValidator("int,duplicate", "", true, nil, nil),
	"goods_count": config.ParseValidator("int,duplicate", "", true, nil, nil),
	"gem_id":      config.ParseValidator("int,duplicate", "", true, nil, nil),
	"gem_count":   config.ParseValidator("int,duplicate", "", true, nil, nil),
}

func (dAtA *Cost) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *Cost) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *Cost) Encode() *shared_proto.CostProto {
	out := &shared_proto.CostProto{}
	out.Gold = config.U64ToI32(dAtA.Gold)
	out.Food = config.U64ToI32(dAtA.Food)
	out.Wood = config.U64ToI32(dAtA.Wood)
	out.Stone = config.U64ToI32(dAtA.Stone)
	out.Yuanbao = config.U64ToI32(dAtA.Yuanbao)
	out.Dianquan = config.U64ToI32(dAtA.Dianquan)
	out.Yinliang = config.U64ToI32(dAtA.Yinliang)
	out.GuildContributionCoin = config.U64ToI32(dAtA.GuildContributionCoin)
	out.Jade = config.U64ToI32(dAtA.Jade)
	out.JadeOre = config.U64ToI32(dAtA.JadeOre)
	if dAtA.Goods != nil {
		out.GoodsId = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.Goods))
	}
	out.GoodsCount = config.U64a2I32a(dAtA.GoodsCount)
	if dAtA.Gem != nil {
		out.GemId = config.U64a2I32a(goods.GetGemDataKeyArray(dAtA.Gem))
	}
	out.GemCount = config.U64a2I32a(dAtA.GemCount)
	out.IsNotEmpty = dAtA.IsNotEmpty

	return out
}

func ArrayEncodeCost(datas []*Cost) []*shared_proto.CostProto {

	out := make([]*shared_proto.CostProto, 0, len(datas))
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

func (dAtA *Cost) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("goods_id", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetGoodsData(v)
		if obj != nil {
			dAtA.Goods = append(dAtA.Goods, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[goods_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("goods_id"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("gem_id", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetGemData(v)
		if obj != nil {
			dAtA.Gem = append(dAtA.Gem, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[gem_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("gem_id"), *pArSeR)
		}
	}

	return nil
}

// start with ExchangeData ----------------------------------

func NewExchangeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ExchangeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrExchangeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ExchangeData{}

	// releated field: Cost
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrExchangeData = map[string]*config.Validator{

	"cost":  config.ParseValidator("string", "", false, nil, nil),
	"prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ExchangeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ExchangeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ExchangeData) Encode() *shared_proto.ExchangeDataProto {
	out := &shared_proto.ExchangeDataProto{}
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeExchangeData(datas []*ExchangeData) []*shared_proto.ExchangeDataProto {

	out := make([]*shared_proto.ExchangeDataProto, 0, len(datas))
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

func (dAtA *ExchangeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GuildLevelPrize ----------------------------------

func LoadGuildLevelPrize(gos *config.GameObjects) (map[uint64]*GuildLevelPrize, map[*GuildLevelPrize]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildLevelPrizePath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildLevelPrize, len(lIsT))
	pArSeRmAp := make(map[*GuildLevelPrize]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildLevelPrize) {
			continue
		}

		dAtA, err := NewGuildLevelPrize(fIlEnAmE, pArSeR)
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

func SetRelatedGuildLevelPrize(dAtAmAp map[*GuildLevelPrize]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildLevelPrizePath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildLevelPrizeKeyArray(datas []*GuildLevelPrize) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildLevelPrize(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildLevelPrize, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildLevelPrize)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildLevelPrize{}

	dAtA.GroupId = pArSeR.Uint64("group_id")
	dAtA.GuildLevel = pArSeR.Uint64("guild_level")
	// releated field: Prize

	// calculate fields
	dAtA.Id = GenGuildLevelPrizeId(dAtA.GroupId, dAtA.GuildLevel)

	return dAtA, nil
}

var vAlIdAtOrGuildLevelPrize = map[string]*config.Validator{

	"group_id":    config.ParseValidator("int>0", "", false, nil, nil),
	"guild_level": config.ParseValidator("int>0", "", false, nil, nil),
	"prize":       config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *GuildLevelPrize) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildLevelPrize) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildLevelPrize) Encode() *shared_proto.GuildLevelPrizeProto {
	out := &shared_proto.GuildLevelPrizeProto{}
	out.GroupId = config.U64ToI32(dAtA.GroupId)
	out.GuildLevel = config.U64ToI32(dAtA.GuildLevel)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeGuildLevelPrize(datas []*GuildLevelPrize) []*shared_proto.GuildLevelPrizeProto {

	out := make([]*shared_proto.GuildLevelPrizeProto, 0, len(datas))
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

func (dAtA *GuildLevelPrize) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with Plunder ----------------------------------

func LoadPlunder(gos *config.GameObjects) (map[uint64]*Plunder, map[*Plunder]*config.ObjectParser, error) {
	fIlEnAmE := confpath.PlunderPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*Plunder, len(lIsT))
	pArSeRmAp := make(map[*Plunder]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrPlunder) {
			continue
		}

		dAtA, err := NewPlunder(fIlEnAmE, pArSeR)
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

func SetRelatedPlunder(dAtAmAp map[*Plunder]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.PlunderPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetPlunderKeyArray(datas []*Plunder) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewPlunder(fIlEnAmE string, pArSeR *config.ObjectParser) (*Plunder, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrPlunder)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &Plunder{}

	dAtA.Id = pArSeR.Uint64("id")
	if pArSeR.KeyExist("unsafe_gold") {
		dAtA.UnsafeGold, err = data.ParseRandAmount(pArSeR.String("unsafe_gold"))
	} else {
		dAtA.UnsafeGold, err = data.ParseRandAmount("")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[unsafe_gold] 解析失败(data.ParseRandAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("unsafe_gold"), dAtA)
	}

	if pArSeR.KeyExist("unsafe_food") {
		dAtA.UnsafeFood, err = data.ParseRandAmount(pArSeR.String("unsafe_food"))
	} else {
		dAtA.UnsafeFood, err = data.ParseRandAmount("")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[unsafe_food] 解析失败(data.ParseRandAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("unsafe_food"), dAtA)
	}

	if pArSeR.KeyExist("unsafe_wood") {
		dAtA.UnsafeWood, err = data.ParseRandAmount(pArSeR.String("unsafe_wood"))
	} else {
		dAtA.UnsafeWood, err = data.ParseRandAmount("")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[unsafe_wood] 解析失败(data.ParseRandAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("unsafe_wood"), dAtA)
	}

	if pArSeR.KeyExist("unsafe_stone") {
		dAtA.UnsafeStone, err = data.ParseRandAmount(pArSeR.String("unsafe_stone"))
	} else {
		dAtA.UnsafeStone, err = data.ParseRandAmount("")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[unsafe_stone] 解析失败(data.ParseRandAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("unsafe_stone"), dAtA)
	}

	if pArSeR.KeyExist("safe_gold") {
		dAtA.SafeGold, err = data.ParseRandAmount(pArSeR.String("safe_gold"))
	} else {
		dAtA.SafeGold, err = data.ParseRandAmount("")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[safe_gold] 解析失败(data.ParseRandAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("safe_gold"), dAtA)
	}

	if pArSeR.KeyExist("safe_food") {
		dAtA.SafeFood, err = data.ParseRandAmount(pArSeR.String("safe_food"))
	} else {
		dAtA.SafeFood, err = data.ParseRandAmount("")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[safe_food] 解析失败(data.ParseRandAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("safe_food"), dAtA)
	}

	if pArSeR.KeyExist("safe_wood") {
		dAtA.SafeWood, err = data.ParseRandAmount(pArSeR.String("safe_wood"))
	} else {
		dAtA.SafeWood, err = data.ParseRandAmount("")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[safe_wood] 解析失败(data.ParseRandAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("safe_wood"), dAtA)
	}

	if pArSeR.KeyExist("safe_stone") {
		dAtA.SafeStone, err = data.ParseRandAmount(pArSeR.String("safe_stone"))
	} else {
		dAtA.SafeStone, err = data.ParseRandAmount("")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[safe_stone] 解析失败(data.ParseRandAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("safe_stone"), dAtA)
	}

	dAtA.HeroExp, err = data.ParseRandAmount(pArSeR.String("hero_exp"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[hero_exp] 解析失败(data.ParseRandAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("hero_exp"), dAtA)
	}

	dAtA.CaptainExp, err = data.ParseRandAmount(pArSeR.String("captain_exp"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[captain_exp] 解析失败(data.ParseRandAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("captain_exp"), dAtA)
	}

	// releated field: Item
	stringKeys = pArSeR.StringArray("item_rate", "", false)
	dAtA.ItemRate = make([]*data.Rate, 0, len(stringKeys))
	for _, v := range stringKeys {
		obj, err := data.ParseRate(v)
		if err != nil {
			return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[item_rate] 解析失败(data.ParseRate)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("item_rate"), dAtA)
		}
		dAtA.ItemRate = append(dAtA.ItemRate, obj)
	}

	// releated field: Group
	stringKeys = pArSeR.StringArray("group_rate", "", false)
	dAtA.GroupRate = make([]*data.Rate, 0, len(stringKeys))
	for _, v := range stringKeys {
		obj, err := data.ParseRate(v)
		if err != nil {
			return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[group_rate] 解析失败(data.ParseRate)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("group_rate"), dAtA)
		}
		dAtA.GroupRate = append(dAtA.GroupRate, obj)
	}

	return dAtA, nil
}

var vAlIdAtOrPlunder = map[string]*config.Validator{

	"id":           config.ParseValidator("int>0", "", false, nil, nil),
	"unsafe_gold":  config.ParseValidator("string", "", false, nil, []string{""}),
	"unsafe_food":  config.ParseValidator("string", "", false, nil, []string{""}),
	"unsafe_wood":  config.ParseValidator("string", "", false, nil, []string{""}),
	"unsafe_stone": config.ParseValidator("string", "", false, nil, []string{""}),
	"safe_gold":    config.ParseValidator("string", "", false, nil, []string{""}),
	"safe_food":    config.ParseValidator("string", "", false, nil, []string{""}),
	"safe_wood":    config.ParseValidator("string", "", false, nil, []string{""}),
	"safe_stone":   config.ParseValidator("string", "", false, nil, []string{""}),
	"hero_exp":     config.ParseValidator("string", "", false, nil, nil),
	"captain_exp":  config.ParseValidator("string", "", false, nil, nil),
	"item":         config.ParseValidator("string", "", true, nil, nil),
	"item_rate":    config.ParseValidator("string,duplicate", "", true, nil, nil),
	"group":        config.ParseValidator("string", "", true, nil, nil),
	"group_rate":   config.ParseValidator("string,duplicate", "", true, nil, nil),
}

func (dAtA *Plunder) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("item", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetPlunderItem(v)
		if obj != nil {
			dAtA.Item = append(dAtA.Item, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[item] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("item"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("group", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetPlunderGroup(v)
		if obj != nil {
			dAtA.Group = append(dAtA.Group, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[group] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("group"), *pArSeR)
		}
	}

	return nil
}

// start with PlunderGroup ----------------------------------

func LoadPlunderGroup(gos *config.GameObjects) (map[uint64]*PlunderGroup, map[*PlunderGroup]*config.ObjectParser, error) {
	fIlEnAmE := confpath.PlunderGroupPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*PlunderGroup, len(lIsT))
	pArSeRmAp := make(map[*PlunderGroup]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrPlunderGroup) {
			continue
		}

		dAtA, err := NewPlunderGroup(fIlEnAmE, pArSeR)
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

func SetRelatedPlunderGroup(dAtAmAp map[*PlunderGroup]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.PlunderGroupPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetPlunderGroupKeyArray(datas []*PlunderGroup) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewPlunderGroup(fIlEnAmE string, pArSeR *config.ObjectParser) (*PlunderGroup, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrPlunderGroup)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &PlunderGroup{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Item
	dAtA.Weight = pArSeR.Uint64Array("weight", "", false)

	return dAtA, nil
}

var vAlIdAtOrPlunderGroup = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"item":   config.ParseValidator("string", "", true, nil, nil),
	"weight": config.ParseValidator(",duplicate", "", true, nil, nil),
}

func (dAtA *PlunderGroup) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("item", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetPlunderItem(v)
		if obj != nil {
			dAtA.Item = append(dAtA.Item, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[item] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("item"), *pArSeR)
		}
	}

	return nil
}

// start with PlunderItem ----------------------------------

func LoadPlunderItem(gos *config.GameObjects) (map[uint64]*PlunderItem, map[*PlunderItem]*config.ObjectParser, error) {
	fIlEnAmE := confpath.PlunderItemPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*PlunderItem, len(lIsT))
	pArSeRmAp := make(map[*PlunderItem]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrPlunderItem) {
			continue
		}

		dAtA, err := NewPlunderItem(fIlEnAmE, pArSeR)
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

func SetRelatedPlunderItem(dAtAmAp map[*PlunderItem]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.PlunderItemPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetPlunderItemKeyArray(datas []*PlunderItem) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewPlunderItem(fIlEnAmE string, pArSeR *config.ObjectParser) (*PlunderItem, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrPlunderItem)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &PlunderItem{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Goods
	// releated field: Equipment
	// releated field: Gem
	// releated field: Baowu
	// releated field: Captain
	dAtA.Count = pArSeR.Uint64("count")

	return dAtA, nil
}

var vAlIdAtOrPlunderItem = map[string]*config.Validator{

	"id":        config.ParseValidator("int>0", "", false, nil, nil),
	"goods":     config.ParseValidator("string", "", false, nil, nil),
	"equipment": config.ParseValidator("string", "", false, nil, nil),
	"gem":       config.ParseValidator("string", "", false, nil, nil),
	"baowu":     config.ParseValidator("string", "", false, nil, nil),
	"captain":   config.ParseValidator("string", "", false, nil, nil),
	"count":     config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *PlunderItem) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Goods = cOnFigS.GetGoodsData(pArSeR.Uint64("goods"))
	if dAtA.Goods == nil && pArSeR.Uint64("goods") != 0 {
		return errors.Errorf("%s 配置的关联字段[goods] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("goods"), *pArSeR)
	}

	dAtA.Equipment = cOnFigS.GetEquipmentData(pArSeR.Uint64("equipment"))
	if dAtA.Equipment == nil && pArSeR.Uint64("equipment") != 0 {
		return errors.Errorf("%s 配置的关联字段[equipment] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("equipment"), *pArSeR)
	}

	dAtA.Gem = cOnFigS.GetGemData(pArSeR.Uint64("gem"))
	if dAtA.Gem == nil && pArSeR.Uint64("gem") != 0 {
		return errors.Errorf("%s 配置的关联字段[gem] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("gem"), *pArSeR)
	}

	dAtA.Baowu = cOnFigS.GetBaowuData(pArSeR.Uint64("baowu"))
	if dAtA.Baowu == nil && pArSeR.Uint64("baowu") != 0 {
		return errors.Errorf("%s 配置的关联字段[baowu] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("baowu"), *pArSeR)
	}

	dAtA.Captain = cOnFigS.GetResCaptainData(pArSeR.Uint64("captain"))
	if dAtA.Captain == nil && pArSeR.Uint64("captain") != 0 {
		return errors.Errorf("%s 配置的关联字段[captain] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain"), *pArSeR)
	}

	return nil
}

// start with PlunderPrize ----------------------------------

func LoadPlunderPrize(gos *config.GameObjects) (map[uint64]*PlunderPrize, map[*PlunderPrize]*config.ObjectParser, error) {
	fIlEnAmE := confpath.PlunderPrizePath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*PlunderPrize, len(lIsT))
	pArSeRmAp := make(map[*PlunderPrize]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrPlunderPrize) {
			continue
		}

		dAtA, err := NewPlunderPrize(fIlEnAmE, pArSeR)
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

func SetRelatedPlunderPrize(dAtAmAp map[*PlunderPrize]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.PlunderPrizePath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetPlunderPrizeKeyArray(datas []*PlunderPrize) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewPlunderPrize(fIlEnAmE string, pArSeR *config.ObjectParser) (*PlunderPrize, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrPlunderPrize)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &PlunderPrize{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Plunder
	dAtA.GuildLevelPrizeGroupId = pArSeR.Uint64("guild_level_prize_group_id")
	// skip field: GuildLevelPrizes
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrPlunderPrize = map[string]*config.Validator{

	"id":                         config.ParseValidator("int>0", "", false, nil, nil),
	"plunder":                    config.ParseValidator("string", "", false, nil, nil),
	"guild_level_prize_group_id": config.ParseValidator("uint", "", false, nil, nil),
	"prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *PlunderPrize) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Plunder = cOnFigS.GetPlunder(pArSeR.Uint64("plunder"))
	if dAtA.Plunder == nil && pArSeR.Uint64("plunder") != 0 {
		return errors.Errorf("%s 配置的关联字段[plunder] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder"), *pArSeR)
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with Prize ----------------------------------

func LoadPrize(gos *config.GameObjects) (map[int]*Prize, map[*Prize]*config.ObjectParser, error) {
	fIlEnAmE := confpath.PrizePath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[int]*Prize, len(lIsT))
	pArSeRmAp := make(map[*Prize]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrPrize) {
			continue
		}

		dAtA, err := NewPrize(fIlEnAmE, pArSeR)
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

func SetRelatedPrize(dAtAmAp map[*Prize]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.PrizePath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetPrizeKeyArray(datas []*Prize) []int {

	out := make([]int, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewPrize(fIlEnAmE string, pArSeR *config.ObjectParser) (*Prize, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrPrize)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &Prize{}

	dAtA.Id = pArSeR.Int("id")
	dAtA.SafeGold = pArSeR.Uint64("safe_gold")
	dAtA.SafeFood = pArSeR.Uint64("safe_food")
	dAtA.SafeWood = pArSeR.Uint64("safe_wood")
	dAtA.SafeStone = pArSeR.Uint64("safe_stone")
	dAtA.UnsafeGold = pArSeR.Uint64("unsafe_gold")
	dAtA.UnsafeFood = pArSeR.Uint64("unsafe_food")
	dAtA.UnsafeWood = pArSeR.Uint64("unsafe_wood")
	dAtA.UnsafeStone = pArSeR.Uint64("unsafe_stone")
	// releated field: Goods
	dAtA.GoodsCount = pArSeR.Uint64Array("goods_count", "", false)
	// releated field: Equipment
	dAtA.EquipmentCount = pArSeR.Uint64Array("equipment_count", "", false)
	// releated field: Gem
	dAtA.GemCount = pArSeR.Uint64Array("gem_count", "", false)
	// releated field: Baowu
	dAtA.BaowuCount = pArSeR.Uint64Array("baowu_count", "", false)
	// releated field: Captain
	dAtA.CaptainCount = pArSeR.Uint64Array("captain_count", "", false)
	dAtA.CaptainExp = pArSeR.Uint64("captain_exp")
	dAtA.HeroExp = pArSeR.Uint64("hero_exp")
	dAtA.Prosperity = 0
	if pArSeR.KeyExist("prosperity") {
		dAtA.Prosperity = pArSeR.Uint64("prosperity")
	}

	dAtA.Yuanbao = 0
	if pArSeR.KeyExist("yuanbao") {
		dAtA.Yuanbao = pArSeR.Uint64("yuanbao")
	}

	dAtA.Dianquan = 0
	if pArSeR.KeyExist("dianquan") {
		dAtA.Dianquan = pArSeR.Uint64("dianquan")
	}

	dAtA.Yinliang = 0
	if pArSeR.KeyExist("yinliang") {
		dAtA.Yinliang = pArSeR.Uint64("yinliang")
	}

	dAtA.GuildContributionCoin = 0
	if pArSeR.KeyExist("guild_contribution_coin") {
		dAtA.GuildContributionCoin = pArSeR.Uint64("guild_contribution_coin")
	}

	dAtA.Jade = pArSeR.Uint64("jade")
	dAtA.JadeOre = pArSeR.Uint64("jade_ore")
	dAtA.VipExp = 0
	if pArSeR.KeyExist("vip_exp") {
		dAtA.VipExp = pArSeR.Uint64("vip_exp")
	}

	dAtA.Sp = 0
	if pArSeR.KeyExist("sp") {
		dAtA.Sp = pArSeR.Uint64("sp")
	}

	// skip field: IsNotEmpty
	// releated field: Sort

	// calculate fields
	dAtA.Gold = dAtA.SafeGold + dAtA.UnsafeGold
	// calculate fields
	dAtA.Food = dAtA.SafeFood + dAtA.UnsafeFood
	// calculate fields
	dAtA.Wood = dAtA.SafeWood + dAtA.UnsafeWood
	// calculate fields
	dAtA.Stone = dAtA.SafeStone + dAtA.UnsafeStone

	return dAtA, nil
}

var vAlIdAtOrPrize = map[string]*config.Validator{

	"id":                      config.ParseValidator("int>0", "", false, nil, nil),
	"safe_gold":               config.ParseValidator("uint", "", false, nil, nil),
	"safe_food":               config.ParseValidator("uint", "", false, nil, nil),
	"safe_wood":               config.ParseValidator("uint", "", false, nil, nil),
	"safe_stone":              config.ParseValidator("uint", "", false, nil, nil),
	"unsafe_gold":             config.ParseValidator("uint", "", false, nil, nil),
	"unsafe_food":             config.ParseValidator("uint", "", false, nil, nil),
	"unsafe_wood":             config.ParseValidator("uint", "", false, nil, nil),
	"unsafe_stone":            config.ParseValidator("uint", "", false, nil, nil),
	"goods_id":                config.ParseValidator("string", "", true, nil, nil),
	"goods_count":             config.ParseValidator("int,duplicate,", "", true, nil, nil),
	"equipment_id":            config.ParseValidator("int,duplicate", "", true, nil, nil),
	"equipment_count":         config.ParseValidator("int,duplicate", "", true, nil, nil),
	"gem_id":                  config.ParseValidator("int,duplicate", "", true, nil, nil),
	"gem_count":               config.ParseValidator("int,duplicate", "", true, nil, nil),
	"baowu_id":                config.ParseValidator("int,duplicate", "", true, nil, nil),
	"baowu_count":             config.ParseValidator("int,duplicate", "", true, nil, nil),
	"captain_id":              config.ParseValidator("int,duplicate", "", true, nil, nil),
	"captain_count":           config.ParseValidator("int,duplicate", "", true, nil, nil),
	"captain_exp":             config.ParseValidator("uint", "", false, nil, nil),
	"hero_exp":                config.ParseValidator("uint", "", false, nil, nil),
	"prosperity":              config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"yuanbao":                 config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"dianquan":                config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"yinliang":                config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"guild_contribution_coin": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"jade":     config.ParseValidator("uint", "", false, nil, nil),
	"jade_ore": config.ParseValidator("uint", "", false, nil, nil),
	"vip_exp":  config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"sp":       config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"sort":     config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *Prize) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *Prize) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *Prize) Encode() *shared_proto.PrizeProto {
	out := &shared_proto.PrizeProto{}
	out.Gold = config.U64ToI32(dAtA.Gold)
	out.Food = config.U64ToI32(dAtA.Food)
	out.Wood = config.U64ToI32(dAtA.Wood)
	out.Stone = config.U64ToI32(dAtA.Stone)
	out.SafeGold = config.U64ToI32(dAtA.SafeGold)
	out.SafeFood = config.U64ToI32(dAtA.SafeFood)
	out.SafeWood = config.U64ToI32(dAtA.SafeWood)
	out.SafeStone = config.U64ToI32(dAtA.SafeStone)
	if dAtA.Goods != nil {
		out.GoodsId = config.U64a2I32a(goods.GetGoodsDataKeyArray(dAtA.Goods))
	}
	out.GoodsCount = config.U64a2I32a(dAtA.GoodsCount)
	if dAtA.Equipment != nil {
		out.EquipmentId = config.U64a2I32a(goods.GetEquipmentDataKeyArray(dAtA.Equipment))
	}
	out.EquipmentCount = config.U64a2I32a(dAtA.EquipmentCount)
	if dAtA.Gem != nil {
		out.GemId = config.U64a2I32a(goods.GetGemDataKeyArray(dAtA.Gem))
	}
	out.GemCount = config.U64a2I32a(dAtA.GemCount)
	if dAtA.Baowu != nil {
		out.BaowuId = config.U64a2I32a(GetBaowuDataKeyArray(dAtA.Baowu))
	}
	out.BaowuCount = config.U64a2I32a(dAtA.BaowuCount)
	if dAtA.Captain != nil {
		out.CaptainId = config.U64a2I32a(GetResCaptainDataKeyArray(dAtA.Captain))
	}
	out.CaptainCount = config.U64a2I32a(dAtA.CaptainCount)
	out.CaptainExp = config.U64ToI32(dAtA.CaptainExp)
	out.HeroExp = config.U64ToI32(dAtA.HeroExp)
	out.Prosperity = config.U64ToI32(dAtA.Prosperity)
	out.Yuanbao = config.U64ToI32(dAtA.Yuanbao)
	out.Dianquan = config.U64ToI32(dAtA.Dianquan)
	out.Yinliang = config.U64ToI32(dAtA.Yinliang)
	out.GuildContributionCoin = config.U64ToI32(dAtA.GuildContributionCoin)
	out.Jade = config.U64ToI32(dAtA.Jade)
	out.JadeOre = config.U64ToI32(dAtA.JadeOre)
	out.VipExp = config.U64ToI32(dAtA.VipExp)
	out.Sp = config.U64ToI32(dAtA.Sp)
	out.IsNotEmpty = dAtA.IsNotEmpty
	if dAtA.Sort != nil {
		out.AmountShowSortId = config.U64ToI32(dAtA.Sort.Id)
	}

	return out
}

func ArrayEncodePrize(datas []*Prize) []*shared_proto.PrizeProto {

	out := make([]*shared_proto.PrizeProto, 0, len(datas))
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

func (dAtA *Prize) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("goods_id", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetGoodsData(v)
		if obj != nil {
			dAtA.Goods = append(dAtA.Goods, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[goods_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("goods_id"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("equipment_id", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetEquipmentData(v)
		if obj != nil {
			dAtA.Equipment = append(dAtA.Equipment, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[equipment_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("equipment_id"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("gem_id", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetGemData(v)
		if obj != nil {
			dAtA.Gem = append(dAtA.Gem, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[gem_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("gem_id"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("baowu_id", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetBaowuData(v)
		if obj != nil {
			dAtA.Baowu = append(dAtA.Baowu, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[baowu_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("baowu_id"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("captain_id", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetResCaptainData(v)
		if obj != nil {
			dAtA.Captain = append(dAtA.Captain, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[captain_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain_id"), *pArSeR)
		}
	}

	dAtA.Sort = cOnFigS.GetAmountShowSortData(pArSeR.Uint64("sort"))
	if dAtA.Sort == nil && pArSeR.Uint64("sort") != 0 {
		return errors.Errorf("%s 配置的关联字段[sort] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("sort"), *pArSeR)
	}

	return nil
}

// start with ResCaptainData ----------------------------------

func LoadResCaptainData(gos *config.GameObjects) (map[uint64]*ResCaptainData, map[*ResCaptainData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ResCaptainDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ResCaptainData, len(lIsT))
	pArSeRmAp := make(map[*ResCaptainData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrResCaptainData) {
			continue
		}

		dAtA, err := NewResCaptainData(fIlEnAmE, pArSeR)
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

func SetRelatedResCaptainData(dAtAmAp map[*ResCaptainData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ResCaptainDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetResCaptainDataKeyArray(datas []*ResCaptainData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewResCaptainData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ResCaptainData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrResCaptainData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ResCaptainData{}

	dAtA.Id = pArSeR.Uint64("id")

	return dAtA, nil
}

var vAlIdAtOrResCaptainData = map[string]*config.Validator{

	"id": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *ResCaptainData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetAmountShowSortData(uint64) *AmountShowSortData
	GetBaowuData(uint64) *BaowuData
	GetConditionPlunderItem(uint64) *ConditionPlunderItem
	GetCost(int) *Cost
	GetEquipmentData(uint64) *goods.EquipmentData
	GetGemData(uint64) *goods.GemData
	GetGoodsData(uint64) *goods.GoodsData
	GetGoodsQuality(uint64) *goods.GoodsQuality
	GetIcon(string) *icon.Icon
	GetPlunder(uint64) *Plunder
	GetPlunderGroup(uint64) *PlunderGroup
	GetPlunderItem(uint64) *PlunderItem
	GetPlunderPrize(uint64) *PlunderPrize
	GetPrize(int) *Prize
	GetResCaptainData(uint64) *ResCaptainData
}
