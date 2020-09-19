// AUTO_GEN, DONT MODIFY!!!
package country

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/body"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/head"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/mingcdata"
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

// start with CountryData ----------------------------------

func LoadCountryData(gos *config.GameObjects) (map[uint64]*CountryData, map[*CountryData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CountryDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CountryData, len(lIsT))
	pArSeRmAp := make(map[*CountryData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCountryData) {
			continue
		}

		dAtA, err := NewCountryData(fIlEnAmE, pArSeR)
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

func SetRelatedCountryData(dAtAmAp map[*CountryData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CountryDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCountryDataKeyArray(datas []*CountryData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCountryData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CountryData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCountryData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CountryData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = ""
	if pArSeR.KeyExist("desc") {
		dAtA.Desc = pArSeR.String("desc")
	}

	dAtA.DefaultPrestige = 0
	if pArSeR.KeyExist("default_prestige") {
		dAtA.DefaultPrestige = pArSeR.Uint64("default_prestige")
	}

	// releated field: Capital
	dAtA.BornCenterX = pArSeR.Uint64("born_center_x")
	dAtA.BornCenterY = pArSeR.Uint64("born_center_y")
	dAtA.BornRadiusX = pArSeR.Uint64("born_radius_x")
	dAtA.BornRadiusY = pArSeR.Uint64("born_radius_y")
	for _, v := range pArSeR.StringArray("npc_official", "", false) {
		x := shared_proto.CountryOfficialType(shared_proto.CountryOfficialType_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.CountryOfficialType(i)
		}
		dAtA.NpcOfficial = append(dAtA.NpcOfficial, x)
	}

	dAtA.NpcId = pArSeR.Uint64Array("npc_id", "", false)

	return dAtA, nil
}

var vAlIdAtOrCountryData = map[string]*config.Validator{

	"id":               config.ParseValidator("int>0", "", false, nil, nil),
	"name":             config.ParseValidator("string", "", false, nil, nil),
	"desc":             config.ParseValidator("string", "", false, nil, []string{""}),
	"default_prestige": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"capital":          config.ParseValidator("string", "", false, nil, nil),
	"born_center_x":    config.ParseValidator("uint", "", false, nil, nil),
	"born_center_y":    config.ParseValidator("uint", "", false, nil, nil),
	"born_radius_x":    config.ParseValidator("uint", "", false, nil, nil),
	"born_radius_y":    config.ParseValidator("uint", "", false, nil, nil),
	"npc_official":     config.ParseValidator(",duplicate", "", true, config.EnumMapKeys(shared_proto.CountryOfficialType_value, 0), nil),
	"npc_id":           config.ParseValidator("uint,duplicate", "", true, nil, nil),
}

func (dAtA *CountryData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CountryData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CountryData) Encode() *shared_proto.CountryDataProto {
	out := &shared_proto.CountryDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.DefaultPrestige = config.U64ToI32(dAtA.DefaultPrestige)
	if dAtA.Capital != nil {
		out.Capital = config.U64ToI32(dAtA.Capital.Id)
	}
	out.NpcOfficial = dAtA.NpcOfficial
	out.NpcId = config.U64a2I32a(dAtA.NpcId)

	return out
}

func ArrayEncodeCountryData(datas []*CountryData) []*shared_proto.CountryDataProto {

	out := make([]*shared_proto.CountryDataProto, 0, len(datas))
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

func (dAtA *CountryData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Capital = cOnFigS.GetMingcBaseData(pArSeR.Uint64("capital"))
	if dAtA.Capital == nil {
		return errors.Errorf("%s 配置的关联字段[capital] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("capital"), *pArSeR)
	}

	return nil
}

// start with CountryMiscData ----------------------------------

func LoadCountryMiscData(gos *config.GameObjects) (*CountryMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.CountryMiscDataPath
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

	dAtA, err := NewCountryMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedCountryMiscData(gos *config.GameObjects, dAtA *CountryMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CountryMiscDataPath
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

func NewCountryMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CountryMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCountryMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CountryMiscData{}

	// releated field: NormalChangeCountryGoods
	dAtA.NewHeroChangeCountryCd, err = config.ParseDuration(pArSeR.String("new_hero_change_country_cd"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[new_hero_change_country_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("new_hero_change_country_cd"), dAtA)
	}

	dAtA.NormalChangeCountryCd, err = config.ParseDuration(pArSeR.String("normal_change_country_cd"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[normal_change_country_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("normal_change_country_cd"), dAtA)
	}

	dAtA.NewHeroMaxLevel = pArSeR.Uint64("new_hero_max_level")
	dAtA.ChangeNameVoteDuration, err = config.ParseDuration(pArSeR.String("change_name_vote_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[change_name_vote_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("change_name_vote_duration"), dAtA)
	}

	// releated field: ChangeNameCost
	if pArSeR.KeyExist("change_name_cd") {
		dAtA.ChangeNameCd, err = config.ParseDuration(pArSeR.String("change_name_cd"))
	} else {
		dAtA.ChangeNameCd, err = config.ParseDuration("168h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[change_name_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("change_name_cd"), dAtA)
	}

	dAtA.MaxSearchHeroDefaultCount = 20
	if pArSeR.KeyExist("max_search_hero_default_count") {
		dAtA.MaxSearchHeroDefaultCount = pArSeR.Int("max_search_hero_default_count")
	}

	dAtA.MaxSearchHeroByNameCount = 200
	if pArSeR.KeyExist("max_search_hero_by_name_count") {
		dAtA.MaxSearchHeroByNameCount = pArSeR.Int("max_search_hero_by_name_count")
	}

	return dAtA, nil
}

var vAlIdAtOrCountryMiscData = map[string]*config.Validator{

	"normal_change_country_goods":   config.ParseValidator("string", "", false, nil, nil),
	"new_hero_change_country_cd":    config.ParseValidator("string", "", false, nil, nil),
	"normal_change_country_cd":      config.ParseValidator("string", "", false, nil, nil),
	"new_hero_max_level":            config.ParseValidator("int>0", "", false, nil, nil),
	"change_name_vote_duration":     config.ParseValidator("string", "", false, nil, nil),
	"change_name_cost":              config.ParseValidator("string", "", false, nil, nil),
	"change_name_cd":                config.ParseValidator("string", "", false, nil, []string{"168h"}),
	"max_search_hero_default_count": config.ParseValidator("int>0", "", false, nil, []string{"20"}),
	"max_search_hero_by_name_count": config.ParseValidator("int>0", "", false, nil, []string{"200"}),
}

func (dAtA *CountryMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CountryMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CountryMiscData) Encode() *shared_proto.CountryMiscDataProto {
	out := &shared_proto.CountryMiscDataProto{}
	if dAtA.NormalChangeCountryGoods != nil {
		out.NormalChangeCountryGoods = config.U64ToI32(dAtA.NormalChangeCountryGoods.Id)
	}
	out.NewHeroChangeCountryCd = config.Duration2I32Seconds(dAtA.NewHeroChangeCountryCd)
	out.NormalChangeCountryCd = config.Duration2I32Seconds(dAtA.NormalChangeCountryCd)
	out.NewHeroMaxLevel = config.U64ToI32(dAtA.NewHeroMaxLevel)
	out.ChangeNameVoteDuration = config.Duration2I32Seconds(dAtA.ChangeNameVoteDuration)
	if dAtA.ChangeNameCost != nil {
		out.ChangeNameCost = dAtA.ChangeNameCost.Encode()
	}
	out.ChangeNameCd = config.Duration2I32Seconds(dAtA.ChangeNameCd)
	out.MaxSearchHeroDefaultCount = int32(dAtA.MaxSearchHeroDefaultCount)
	out.MaxSearchHeroByNameCount = int32(dAtA.MaxSearchHeroByNameCount)

	return out
}

func ArrayEncodeCountryMiscData(datas []*CountryMiscData) []*shared_proto.CountryMiscDataProto {

	out := make([]*shared_proto.CountryMiscDataProto, 0, len(datas))
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

func (dAtA *CountryMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.NormalChangeCountryGoods = cOnFigS.GetGoodsData(pArSeR.Uint64("normal_change_country_goods"))
	if dAtA.NormalChangeCountryGoods == nil {
		return errors.Errorf("%s 配置的关联字段[normal_change_country_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("normal_change_country_goods"), *pArSeR)
	}

	dAtA.ChangeNameCost = cOnFigS.GetCost(pArSeR.Int("change_name_cost"))
	if dAtA.ChangeNameCost == nil {
		return errors.Errorf("%s 配置的关联字段[change_name_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("change_name_cost"), *pArSeR)
	}

	return nil
}

// start with CountryOfficialData ----------------------------------

func LoadCountryOfficialData(gos *config.GameObjects) (map[int]*CountryOfficialData, map[*CountryOfficialData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CountryOfficialDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[int]*CountryOfficialData, len(lIsT))
	pArSeRmAp := make(map[*CountryOfficialData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCountryOfficialData) {
			continue
		}

		dAtA, err := NewCountryOfficialData(fIlEnAmE, pArSeR)
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

func SetRelatedCountryOfficialData(dAtAmAp map[*CountryOfficialData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CountryOfficialDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCountryOfficialDataKeyArray(datas []*CountryOfficialData) []int {

	out := make([]int, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCountryOfficialData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CountryOfficialData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCountryOfficialData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CountryOfficialData{}

	dAtA.Id = pArSeR.Int("id")
	// skip field: OfficialType
	dAtA.Name = pArSeR.String("name")
	// releated field: BuildingEffect
	// releated field: Buff
	dAtA.Count = pArSeR.Int("count")
	// releated field: Salary
	// releated field: ShowSalary
	// releated field: Icon
	// releated field: Head
	// releated field: Body
	dAtA.Cd, err = config.ParseDuration(pArSeR.String("cd"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("cd"), dAtA)
	}

	dAtA.EffectDesc = pArSeR.String("effect_desc")
	for _, v := range pArSeR.StringArray("sub_officials", "", false) {
		x := shared_proto.CountryOfficialType(shared_proto.CountryOfficialType_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.CountryOfficialType(i)
		}
		dAtA.SubOfficials = append(dAtA.SubOfficials, x)
	}

	return dAtA, nil
}

var vAlIdAtOrCountryOfficialData = map[string]*config.Validator{

	"id":              config.ParseValidator("int>0", "", false, nil, nil),
	"name":            config.ParseValidator("string", "", false, nil, nil),
	"building_effect": config.ParseValidator("string", "", false, nil, nil),
	"buff":            config.ParseValidator("string", "", false, nil, nil),
	"count":           config.ParseValidator("int>0", "", false, nil, nil),
	"salary":          config.ParseValidator("string", "", false, nil, nil),
	"show_salary":     config.ParseValidator("string", "", false, nil, nil),
	"icon":            config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"head":            config.ParseValidator("string", "", false, nil, nil),
	"body":            config.ParseValidator("string", "", false, nil, nil),
	"cd":              config.ParseValidator("string", "", false, nil, nil),
	"effect_desc":     config.ParseValidator("string", "", false, nil, nil),
	"sub_officials":   config.ParseValidator(",duplicate", "", true, config.EnumMapKeys(shared_proto.CountryOfficialType_value, 0), nil),
}

func (dAtA *CountryOfficialData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CountryOfficialData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CountryOfficialData) Encode() *shared_proto.CountryOfficialDataProto {
	out := &shared_proto.CountryOfficialDataProto{}
	out.OfficialType = dAtA.OfficialType
	out.Name = dAtA.Name
	if dAtA.BuildingEffect != nil {
		out.BuildingEffect = dAtA.BuildingEffect.Encode()
	}
	if dAtA.Buff != nil {
		out.Buff = dAtA.Buff.Encode()
	}
	out.Count = int32(dAtA.Count)
	if dAtA.ShowSalary != nil {
		out.ShowSalary = dAtA.ShowSalary.Encode()
	}
	if dAtA.Icon != nil {
		out.Icon = dAtA.Icon.Id
	}
	if dAtA.Head != nil {
		out.Head = dAtA.Head.Id
	}
	if dAtA.Body != nil {
		out.Body = config.U64ToI32(dAtA.Body.Id)
	}
	out.Cd = config.Duration2I32Seconds(dAtA.Cd)
	out.EffectDesc = dAtA.EffectDesc
	out.SubOfficials = dAtA.SubOfficials

	return out
}

func ArrayEncodeCountryOfficialData(datas []*CountryOfficialData) []*shared_proto.CountryOfficialDataProto {

	out := make([]*shared_proto.CountryOfficialDataProto, 0, len(datas))
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

func (dAtA *CountryOfficialData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.BuildingEffect = cOnFigS.GetBuildingEffectData(pArSeR.Int("building_effect"))
	if dAtA.BuildingEffect == nil && pArSeR.Int("building_effect") != 0 {
		return errors.Errorf("%s 配置的关联字段[building_effect] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("building_effect"), *pArSeR)
	}

	dAtA.Buff = cOnFigS.GetBuffEffectData(pArSeR.Uint64("buff"))
	if dAtA.Buff == nil && pArSeR.Uint64("buff") != 0 {
		return errors.Errorf("%s 配置的关联字段[buff] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("buff"), *pArSeR)
	}

	dAtA.Salary = cOnFigS.GetPlunder(pArSeR.Uint64("salary"))
	if dAtA.Salary == nil {
		return errors.Errorf("%s 配置的关联字段[salary] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("salary"), *pArSeR)
	}

	dAtA.ShowSalary = cOnFigS.GetPrize(pArSeR.Int("show_salary"))
	if dAtA.ShowSalary == nil {
		return errors.Errorf("%s 配置的关联字段[show_salary] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_salary"), *pArSeR)
	}

	if pArSeR.KeyExist("icon") {
		dAtA.Icon = cOnFigS.GetIcon(pArSeR.String("icon"))
	} else {
		dAtA.Icon = cOnFigS.GetIcon("Icon")
	}
	if dAtA.Icon == nil {
		return errors.Errorf("%s 配置的关联字段[icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("icon"), *pArSeR)
	}

	dAtA.Head = cOnFigS.GetHeadData(pArSeR.String("head"))
	if dAtA.Head == nil && pArSeR.String("head") != "" {
		return errors.Errorf("%s 配置的关联字段[head] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("head"), *pArSeR)
	}

	dAtA.Body = cOnFigS.GetBodyData(pArSeR.Uint64("body"))
	if dAtA.Body == nil && pArSeR.Uint64("body") != 0 {
		return errors.Errorf("%s 配置的关联字段[body] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("body"), *pArSeR)
	}

	return nil
}

// start with CountryOfficialNpcData ----------------------------------

func LoadCountryOfficialNpcData(gos *config.GameObjects) (map[uint64]*CountryOfficialNpcData, map[*CountryOfficialNpcData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CountryOfficialNpcDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CountryOfficialNpcData, len(lIsT))
	pArSeRmAp := make(map[*CountryOfficialNpcData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCountryOfficialNpcData) {
			continue
		}

		dAtA, err := NewCountryOfficialNpcData(fIlEnAmE, pArSeR)
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

func SetRelatedCountryOfficialNpcData(dAtAmAp map[*CountryOfficialNpcData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CountryOfficialNpcDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCountryOfficialNpcDataKeyArray(datas []*CountryOfficialNpcData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCountryOfficialNpcData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CountryOfficialNpcData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCountryOfficialNpcData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CountryOfficialNpcData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// releated field: Head
	// releated field: Body

	return dAtA, nil
}

var vAlIdAtOrCountryOfficialNpcData = map[string]*config.Validator{

	"id":   config.ParseValidator("int>0", "", false, nil, nil),
	"name": config.ParseValidator("string", "", false, nil, nil),
	"head": config.ParseValidator("string", "", false, nil, []string{"Head"}),
	"body": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *CountryOfficialNpcData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CountryOfficialNpcData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CountryOfficialNpcData) Encode() *shared_proto.CountryOfficialNpcDataProto {
	out := &shared_proto.CountryOfficialNpcDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	if dAtA.Head != nil {
		out.Head = dAtA.Head.Id
	}
	if dAtA.Body != nil {
		out.Body = config.U64ToI32(dAtA.Body.Id)
	}

	return out
}

func ArrayEncodeCountryOfficialNpcData(datas []*CountryOfficialNpcData) []*shared_proto.CountryOfficialNpcDataProto {

	out := make([]*shared_proto.CountryOfficialNpcDataProto, 0, len(datas))
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

func (dAtA *CountryOfficialNpcData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("head") {
		dAtA.Head = cOnFigS.GetHeadData(pArSeR.String("head"))
	} else {
		dAtA.Head = cOnFigS.GetHeadData("Head")
	}
	if dAtA.Head == nil {
		return errors.Errorf("%s 配置的关联字段[head] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("head"), *pArSeR)
	}

	dAtA.Body = cOnFigS.GetBodyData(pArSeR.Uint64("body"))
	if dAtA.Body == nil {
		return errors.Errorf("%s 配置的关联字段[body] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("body"), *pArSeR)
	}

	return nil
}

// start with FamilyNameData ----------------------------------

func LoadFamilyNameData(gos *config.GameObjects) (map[uint64]*FamilyNameData, map[*FamilyNameData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FamilyNameDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*FamilyNameData, len(lIsT))
	pArSeRmAp := make(map[*FamilyNameData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFamilyNameData) {
			continue
		}

		dAtA, err := NewFamilyNameData(fIlEnAmE, pArSeR)
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

func SetRelatedFamilyNameData(dAtAmAp map[*FamilyNameData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FamilyNameDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFamilyNameDataKeyArray(datas []*FamilyNameData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewFamilyNameData(fIlEnAmE string, pArSeR *config.ObjectParser) (*FamilyNameData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFamilyNameData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FamilyNameData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	// releated field: RecommendCountry

	return dAtA, nil
}

var vAlIdAtOrFamilyNameData = map[string]*config.Validator{

	"id":                config.ParseValidator("int>0", "", false, nil, nil),
	"name":              config.ParseValidator("string", "", false, nil, nil),
	"desc":              config.ParseValidator("string", "", false, nil, nil),
	"recommend_country": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *FamilyNameData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *FamilyNameData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *FamilyNameData) Encode() *shared_proto.FamilyNameDataProto {
	out := &shared_proto.FamilyNameDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	if dAtA.RecommendCountry != nil {
		out.RecommendCountry = config.U64a2I32a(GetCountryDataKeyArray(dAtA.RecommendCountry))
	}

	return out
}

func ArrayEncodeFamilyNameData(datas []*FamilyNameData) []*shared_proto.FamilyNameDataProto {

	out := make([]*shared_proto.FamilyNameDataProto, 0, len(datas))
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

func (dAtA *FamilyNameData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("recommend_country", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetCountryData(v)
		if obj != nil {
			dAtA.RecommendCountry = append(dAtA.RecommendCountry, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[recommend_country] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("recommend_country"), *pArSeR)
		}
	}

	return nil
}

type related_configs interface {
	GetBodyData(uint64) *body.BodyData
	GetBuffEffectData(uint64) *data.BuffEffectData
	GetBuildingEffectData(int) *sub.BuildingEffectData
	GetCost(int) *resdata.Cost
	GetCountryData(uint64) *CountryData
	GetGoodsData(uint64) *goods.GoodsData
	GetHeadData(string) *head.HeadData
	GetIcon(string) *icon.Icon
	GetMingcBaseData(uint64) *mingcdata.MingcBaseData
	GetPlunder(uint64) *resdata.Plunder
	GetPrize(int) *resdata.Prize
}
