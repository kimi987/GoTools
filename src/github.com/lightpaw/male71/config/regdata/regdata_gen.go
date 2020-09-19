// AUTO_GEN, DONT MODIFY!!!
package regdata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/monsterdata"
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

// start with AreaData ----------------------------------

func LoadAreaData(gos *config.GameObjects) (map[uint64]*AreaData, map[*AreaData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.AreaDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*AreaData, len(lIsT))
	pArSeRmAp := make(map[*AreaData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrAreaData) {
			continue
		}

		dAtA, err := NewAreaData(fIlEnAmE, pArSeR)
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

func SetRelatedAreaData(dAtAmAp map[*AreaData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.AreaDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetAreaDataKeyArray(datas []*AreaData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewAreaData(fIlEnAmE string, pArSeR *config.ObjectParser) (*AreaData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrAreaData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &AreaData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.CenterX = pArSeR.Int("center_x")
	dAtA.CenterY = pArSeR.Int("center_y")
	dAtA.MinRadius = pArSeR.Int("min_radius")
	dAtA.MaxRadius = pArSeR.Int("max_radius")
	dAtA.IncludeX = pArSeR.IntArray("include_x", "", false)
	dAtA.IncludeY = pArSeR.IntArray("include_y", "", false)
	dAtA.ExcludeX = pArSeR.IntArray("exclude_x", "", false)
	dAtA.ExcludeY = pArSeR.IntArray("exclude_y", "", false)

	return dAtA, nil
}

var vAlIdAtOrAreaData = map[string]*config.Validator{

	"id":         config.ParseValidator("int>0", "", false, nil, nil),
	"center_x":   config.ParseValidator("int>=0", "", false, nil, nil),
	"center_y":   config.ParseValidator("int>=0", "", false, nil, nil),
	"min_radius": config.ParseValidator("int>=0", "", false, nil, nil),
	"max_radius": config.ParseValidator("int>=0", "", false, nil, nil),
	"include_x":  config.ParseValidator("uint", "", true, nil, nil),
	"include_y":  config.ParseValidator("uint", "", true, nil, nil),
	"exclude_x":  config.ParseValidator("uint", "", true, nil, nil),
	"exclude_y":  config.ParseValidator("uint", "", true, nil, nil),
}

func (dAtA *AreaData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *AreaData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *AreaData) Encode() *shared_proto.AreaDataProto {
	out := &shared_proto.AreaDataProto{}
	out.CenterX = int32(dAtA.CenterX)
	out.CenterY = int32(dAtA.CenterY)
	out.MinRadius = int32(dAtA.MinRadius)
	out.MaxRadius = int32(dAtA.MaxRadius)

	return out
}

func ArrayEncodeAreaData(datas []*AreaData) []*shared_proto.AreaDataProto {

	out := make([]*shared_proto.AreaDataProto, 0, len(datas))
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

func (dAtA *AreaData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with AssemblyData ----------------------------------

func LoadAssemblyData(gos *config.GameObjects) (map[uint64]*AssemblyData, map[*AssemblyData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.AssemblyDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*AssemblyData, len(lIsT))
	pArSeRmAp := make(map[*AssemblyData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrAssemblyData) {
			continue
		}

		dAtA, err := NewAssemblyData(fIlEnAmE, pArSeR)
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

func SetRelatedAssemblyData(dAtAmAp map[*AssemblyData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.AssemblyDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetAssemblyDataKeyArray(datas []*AssemblyData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewAssemblyData(fIlEnAmE string, pArSeR *config.ObjectParser) (*AssemblyData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrAssemblyData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &AssemblyData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.MemberCount = pArSeR.Uint64("member_count")
	stringKeys = pArSeR.StringArray("wait_duration", "", false)
	dAtA.WaitDuration = make([]time.Duration, 0, len(stringKeys))
	for _, v := range stringKeys {
		obj, err := config.ParseDuration(v)
		if err != nil {
			return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[wait_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("wait_duration"), dAtA)
		}
		dAtA.WaitDuration = append(dAtA.WaitDuration, obj)
	}

	return dAtA, nil
}

var vAlIdAtOrAssemblyData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"member_count":  config.ParseValidator("int>0", "", false, nil, nil),
	"wait_duration": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *AssemblyData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *AssemblyData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *AssemblyData) Encode() *shared_proto.AssemblyDataProto {
	out := &shared_proto.AssemblyDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.MemberCount = config.U64ToI32(dAtA.MemberCount)
	out.WaitDuration = config.DurationArr2I32Seconds(dAtA.WaitDuration)

	return out
}

func ArrayEncodeAssemblyData(datas []*AssemblyData) []*shared_proto.AssemblyDataProto {

	out := make([]*shared_proto.AssemblyDataProto, 0, len(datas))
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

func (dAtA *AssemblyData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with BaozNpcData ----------------------------------

func LoadBaozNpcData(gos *config.GameObjects) (map[uint64]*BaozNpcData, map[*BaozNpcData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BaozNpcDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BaozNpcData, len(lIsT))
	pArSeRmAp := make(map[*BaozNpcData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBaozNpcData) {
			continue
		}

		dAtA, err := NewBaozNpcData(fIlEnAmE, pArSeR)
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

func SetRelatedBaozNpcData(dAtAmAp map[*BaozNpcData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BaozNpcDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBaozNpcDataKeyArray(datas []*BaozNpcData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBaozNpcData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BaozNpcData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBaozNpcData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BaozNpcData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Npc
	dAtA.KeepCount = pArSeR.Uint64("keep_count")
	dAtA.RequiredHeroLevel = pArSeR.Uint64("required_hero_level")
	dAtA.RareBaowuIds = pArSeR.Uint64Array("rare_baowu_ids", "", false)

	return dAtA, nil
}

var vAlIdAtOrBaozNpcData = map[string]*config.Validator{

	"id":                  config.ParseValidator("int>0", "", false, nil, nil),
	"npc":                 config.ParseValidator("string", "", false, nil, nil),
	"keep_count":          config.ParseValidator("int>0", "", false, nil, nil),
	"required_hero_level": config.ParseValidator("int>0", "", false, nil, nil),
	"rare_baowu_ids":      config.ParseValidator("uint", "", true, nil, nil),
}

func (dAtA *BaozNpcData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BaozNpcData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BaozNpcData) Encode() *shared_proto.BaozNpcDataProto {
	out := &shared_proto.BaozNpcDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.Npc != nil {
		out.Npc = config.U64ToI32(dAtA.Npc.Id)
	}
	out.RequiredHeroLevel = config.U64ToI32(dAtA.RequiredHeroLevel)

	return out
}

func ArrayEncodeBaozNpcData(datas []*BaozNpcData) []*shared_proto.BaozNpcDataProto {

	out := make([]*shared_proto.BaozNpcDataProto, 0, len(datas))
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

func (dAtA *BaozNpcData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Npc = cOnFigS.GetNpcBaseData(pArSeR.Uint64("npc"))
	if dAtA.Npc == nil {
		return errors.Errorf("%s 配置的关联字段[npc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("npc"), *pArSeR)
	}

	return nil
}

// start with JunTuanNpcData ----------------------------------

func LoadJunTuanNpcData(gos *config.GameObjects) (map[uint64]*JunTuanNpcData, map[*JunTuanNpcData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.JunTuanNpcDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*JunTuanNpcData, len(lIsT))
	pArSeRmAp := make(map[*JunTuanNpcData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrJunTuanNpcData) {
			continue
		}

		dAtA, err := NewJunTuanNpcData(fIlEnAmE, pArSeR)
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

func SetRelatedJunTuanNpcData(dAtAmAp map[*JunTuanNpcData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.JunTuanNpcDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetJunTuanNpcDataKeyArray(datas []*JunTuanNpcData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewJunTuanNpcData(fIlEnAmE string, pArSeR *config.ObjectParser) (*JunTuanNpcData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrJunTuanNpcData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &JunTuanNpcData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Npc
	dAtA.TroopCount = pArSeR.Uint64("troop_count")
	dAtA.Group = pArSeR.Uint64("group")
	dAtA.Level = pArSeR.Uint64("level")
	dAtA.RequiredHeroLevel = pArSeR.Uint64("required_hero_level")

	return dAtA, nil
}

var vAlIdAtOrJunTuanNpcData = map[string]*config.Validator{

	"id":                  config.ParseValidator("int>0", "", false, nil, nil),
	"npc":                 config.ParseValidator("string", "", false, nil, nil),
	"troop_count":         config.ParseValidator("int>0", "", false, nil, nil),
	"group":               config.ParseValidator("int>0", "", false, nil, nil),
	"level":               config.ParseValidator("int>0", "", false, nil, nil),
	"required_hero_level": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *JunTuanNpcData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *JunTuanNpcData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *JunTuanNpcData) Encode() *shared_proto.JunTuanNpcDataProto {
	out := &shared_proto.JunTuanNpcDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.Npc != nil {
		out.Npc = config.U64ToI32(dAtA.Npc.Id)
	}
	out.TroopCount = config.U64ToI32(dAtA.TroopCount)
	out.Group = config.U64ToI32(dAtA.Group)
	out.Level = config.U64ToI32(dAtA.Level)
	out.RequiredHeroLevel = config.U64ToI32(dAtA.RequiredHeroLevel)

	return out
}

func ArrayEncodeJunTuanNpcData(datas []*JunTuanNpcData) []*shared_proto.JunTuanNpcDataProto {

	out := make([]*shared_proto.JunTuanNpcDataProto, 0, len(datas))
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

func (dAtA *JunTuanNpcData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Npc = cOnFigS.GetNpcBaseData(pArSeR.Uint64("npc"))
	if dAtA.Npc == nil {
		return errors.Errorf("%s 配置的关联字段[npc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("npc"), *pArSeR)
	}

	return nil
}

// start with JunTuanNpcPlaceConfig ----------------------------------

func LoadJunTuanNpcPlaceConfig(gos *config.GameObjects) (*JunTuanNpcPlaceConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.JunTuanNpcPlaceConfigPath
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

	dAtA, err := NewJunTuanNpcPlaceConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedJunTuanNpcPlaceConfig(gos *config.GameObjects, dAtA *JunTuanNpcPlaceConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.JunTuanNpcPlaceConfigPath
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

func NewJunTuanNpcPlaceConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*JunTuanNpcPlaceConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrJunTuanNpcPlaceConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &JunTuanNpcPlaceConfig{}

	return dAtA, nil
}

var vAlIdAtOrJunTuanNpcPlaceConfig = map[string]*config.Validator{}

func (dAtA *JunTuanNpcPlaceConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with JunTuanNpcPlaceData ----------------------------------

func LoadJunTuanNpcPlaceData(gos *config.GameObjects) (map[uint64]*JunTuanNpcPlaceData, map[*JunTuanNpcPlaceData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.JunTuanNpcPlaceDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*JunTuanNpcPlaceData, len(lIsT))
	pArSeRmAp := make(map[*JunTuanNpcPlaceData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrJunTuanNpcPlaceData) {
			continue
		}

		dAtA, err := NewJunTuanNpcPlaceData(fIlEnAmE, pArSeR)
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

func SetRelatedJunTuanNpcPlaceData(dAtAmAp map[*JunTuanNpcPlaceData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.JunTuanNpcPlaceDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetJunTuanNpcPlaceDataKeyArray(datas []*JunTuanNpcPlaceData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewJunTuanNpcPlaceData(fIlEnAmE string, pArSeR *config.ObjectParser) (*JunTuanNpcPlaceData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrJunTuanNpcPlaceData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &JunTuanNpcPlaceData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Day = pArSeR.Uint64("day")
	dAtA.Group = pArSeR.Uint64("group")
	// releated field: Area
	dAtA.KeepCount = pArSeR.Uint64("keep_count")

	return dAtA, nil
}

var vAlIdAtOrJunTuanNpcPlaceData = map[string]*config.Validator{

	"id":         config.ParseValidator("int>0", "", false, nil, nil),
	"day":        config.ParseValidator("int>0", "", false, nil, nil),
	"group":      config.ParseValidator("int>0", "", false, nil, nil),
	"area":       config.ParseValidator("string", "", false, nil, nil),
	"keep_count": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *JunTuanNpcPlaceData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Area = cOnFigS.GetAreaData(pArSeR.Uint64("area"))
	if dAtA.Area == nil && pArSeR.Uint64("area") != 0 {
		return errors.Errorf("%s 配置的关联字段[area] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("area"), *pArSeR)
	}

	return nil
}

// start with RegionAreaData ----------------------------------

func LoadRegionAreaData(gos *config.GameObjects) (map[uint64]*RegionAreaData, map[*RegionAreaData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.RegionAreaDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*RegionAreaData, len(lIsT))
	pArSeRmAp := make(map[*RegionAreaData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrRegionAreaData) {
			continue
		}

		dAtA, err := NewRegionAreaData(fIlEnAmE, pArSeR)
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

func SetRelatedRegionAreaData(dAtAmAp map[*RegionAreaData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RegionAreaDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetRegionAreaDataKeyArray(datas []*RegionAreaData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewRegionAreaData(fIlEnAmE string, pArSeR *config.ObjectParser) (*RegionAreaData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRegionAreaData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RegionAreaData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// releated field: Area
	dAtA.WorkshopPrizeCoef = pArSeR.Float64("workshop_prize_coef")

	return dAtA, nil
}

var vAlIdAtOrRegionAreaData = map[string]*config.Validator{

	"id":                  config.ParseValidator("int>0", "", false, nil, nil),
	"name":                config.ParseValidator("string", "", false, nil, nil),
	"area":                config.ParseValidator("string", "", false, nil, nil),
	"workshop_prize_coef": config.ParseValidator("float64", "", false, nil, nil),
}

func (dAtA *RegionAreaData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *RegionAreaData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *RegionAreaData) Encode() *shared_proto.RegionAreaDataProto {
	out := &shared_proto.RegionAreaDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	if dAtA.Area != nil {
		out.Area = dAtA.Area.Encode()
	}
	out.WorkshopPrizeCoef = config.F64ToI32X1000(dAtA.WorkshopPrizeCoef)

	return out
}

func ArrayEncodeRegionAreaData(datas []*RegionAreaData) []*shared_proto.RegionAreaDataProto {

	out := make([]*shared_proto.RegionAreaDataProto, 0, len(datas))
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

func (dAtA *RegionAreaData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Area = cOnFigS.GetAreaData(pArSeR.Uint64("area"))
	if dAtA.Area == nil {
		return errors.Errorf("%s 配置的关联字段[area] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("area"), *pArSeR)
	}

	return nil
}

// start with RegionData ----------------------------------

func LoadRegionData(gos *config.GameObjects) (map[uint64]*RegionData, map[*RegionData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.RegionDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*RegionData, len(lIsT))
	pArSeRmAp := make(map[*RegionData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrRegionData) {
			continue
		}

		dAtA, err := NewRegionData(fIlEnAmE, pArSeR)
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

func SetRelatedRegionData(dAtAmAp map[*RegionData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RegionDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetRegionDataKeyArray(datas []*RegionData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewRegionData(fIlEnAmE string, pArSeR *config.ObjectParser) (*RegionData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRegionData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RegionData{}

	dAtA.RegionType = shared_proto.RegionType(shared_proto.RegionType_value[strings.ToUpper(pArSeR.String("region_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("region_type"), 10, 32); err == nil {
		dAtA.RegionType = shared_proto.RegionType(i)
	}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.BlockXLen = pArSeR.Uint64("block_xlen")
	dAtA.BlockYLen = pArSeR.Uint64("block_ylen")
	dAtA.BlockId = pArSeR.Uint64("block")
	dAtA.CenterBlockX = 3
	if pArSeR.KeyExist("center_block_x") {
		dAtA.CenterBlockX = pArSeR.Uint64("center_block_x")
	}

	dAtA.CenterBlockY = 3
	if pArSeR.KeyExist("center_block_y") {
		dAtA.CenterBlockY = pArSeR.Uint64("center_block_y")
	}

	dAtA.InitRadius = 3
	if pArSeR.KeyExist("init_radius") {
		dAtA.InitRadius = pArSeR.Uint64("init_radius")
	}

	dAtA.RandomMinRadius = 3
	if pArSeR.KeyExist("random_min_radius") {
		dAtA.RandomMinRadius = pArSeR.Uint64("random_min_radius")
	}

	dAtA.RandomMaxRadius = 3
	if pArSeR.KeyExist("random_max_radius") {
		dAtA.RandomMaxRadius = pArSeR.Uint64("random_max_radius")
	}

	// skip field: Block
	// skip field: SubBlockXLen
	// skip field: SubBlockYLen
	dAtA.BornMinRadius = 3
	if pArSeR.KeyExist("born_min_radius") {
		dAtA.BornMinRadius = pArSeR.Uint64("born_min_radius")
	}

	dAtA.BornMaxRadius = 3
	if pArSeR.KeyExist("born_max_radius") {
		dAtA.BornMaxRadius = pArSeR.Uint64("born_max_radius")
	}

	dAtA.GuildMoveBaseMinRadius = 1
	if pArSeR.KeyExist("guild_move_base_min_radius") {
		dAtA.GuildMoveBaseMinRadius = pArSeR.Uint64("guild_move_base_min_radius")
	}

	dAtA.GuildMoveBaseMaxRadius = 32
	if pArSeR.KeyExist("guild_move_base_max_radius") {
		dAtA.GuildMoveBaseMaxRadius = pArSeR.Uint64("guild_move_base_max_radius")
	}

	// skip field: Area

	// calculate fields
	dAtA.Id = RegionDataID(dAtA.RegionType, dAtA.Level)

	return dAtA, nil
}

var vAlIdAtOrRegionData = map[string]*config.Validator{

	"region_type":                config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.RegionType_value, 0), nil),
	"level":                      config.ParseValidator("int>0", "", false, nil, nil),
	"block_xlen":                 config.ParseValidator("int>0", "", false, nil, nil),
	"block_ylen":                 config.ParseValidator("int>0", "", false, nil, nil),
	"block":                      config.ParseValidator("int>0", "", false, nil, nil),
	"center_block_x":             config.ParseValidator("uint", "", false, nil, []string{"3"}),
	"center_block_y":             config.ParseValidator("uint", "", false, nil, []string{"3"}),
	"init_radius":                config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"random_min_radius":          config.ParseValidator("uint", "", false, nil, []string{"3"}),
	"random_max_radius":          config.ParseValidator("uint", "", false, nil, []string{"3"}),
	"born_min_radius":            config.ParseValidator("uint", "", false, nil, []string{"3"}),
	"born_max_radius":            config.ParseValidator("uint", "", false, nil, []string{"3"}),
	"guild_move_base_min_radius": config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"guild_move_base_max_radius": config.ParseValidator("int>0", "", false, nil, []string{"32"}),
}

func (dAtA *RegionData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *RegionData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *RegionData) Encode() *shared_proto.RegionDataProto {
	out := &shared_proto.RegionDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.RegionType = dAtA.RegionType
	out.Level = config.U64ToI32(dAtA.Level)
	out.SubBlockXLen = config.U64ToI32(dAtA.SubBlockXLen)
	out.SubBlockYLen = config.U64ToI32(dAtA.SubBlockYLen)
	out.GuildMoveBaseMinRadius = config.U64ToI32(dAtA.GuildMoveBaseMinRadius)
	out.GuildMoveBaseMaxRadius = config.U64ToI32(dAtA.GuildMoveBaseMaxRadius)

	return out
}

func ArrayEncodeRegionData(datas []*RegionData) []*shared_proto.RegionDataProto {

	out := make([]*shared_proto.RegionDataProto, 0, len(datas))
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

func (dAtA *RegionData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with RegionMonsterData ----------------------------------

func LoadRegionMonsterData(gos *config.GameObjects) (map[uint64]*RegionMonsterData, map[*RegionMonsterData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.RegionMonsterDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*RegionMonsterData, len(lIsT))
	pArSeRmAp := make(map[*RegionMonsterData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrRegionMonsterData) {
			continue
		}

		dAtA, err := NewRegionMonsterData(fIlEnAmE, pArSeR)
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

func SetRelatedRegionMonsterData(dAtAmAp map[*RegionMonsterData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RegionMonsterDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetRegionMonsterDataKeyArray(datas []*RegionMonsterData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewRegionMonsterData(fIlEnAmE string, pArSeR *config.ObjectParser) (*RegionMonsterData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRegionMonsterData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RegionMonsterData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Base
	dAtA.BaseX = pArSeR.Int("base_x")
	dAtA.BaseY = pArSeR.Int("base_y")
	dAtA.RegionId = pArSeR.Uint64("region_id")

	return dAtA, nil
}

var vAlIdAtOrRegionMonsterData = map[string]*config.Validator{

	"id":        config.ParseValidator("int>0", "", false, nil, nil),
	"base":      config.ParseValidator("string", "", false, nil, nil),
	"base_x":    config.ParseValidator("int>0", "", false, nil, nil),
	"base_y":    config.ParseValidator("int>0", "", false, nil, nil),
	"region_id": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *RegionMonsterData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Base = cOnFigS.GetNpcBaseData(pArSeR.Uint64("base"))
	if dAtA.Base == nil {
		return errors.Errorf("%s 配置的关联字段[base] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("base"), *pArSeR)
	}

	return nil
}

// start with RegionMultiLevelNpcData ----------------------------------

func LoadRegionMultiLevelNpcData(gos *config.GameObjects) (map[uint64]*RegionMultiLevelNpcData, map[*RegionMultiLevelNpcData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.RegionMultiLevelNpcDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*RegionMultiLevelNpcData, len(lIsT))
	pArSeRmAp := make(map[*RegionMultiLevelNpcData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrRegionMultiLevelNpcData) {
			continue
		}

		dAtA, err := NewRegionMultiLevelNpcData(fIlEnAmE, pArSeR)
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

func SetRelatedRegionMultiLevelNpcData(dAtAmAp map[*RegionMultiLevelNpcData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RegionMultiLevelNpcDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetRegionMultiLevelNpcDataKeyArray(datas []*RegionMultiLevelNpcData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewRegionMultiLevelNpcData(fIlEnAmE string, pArSeR *config.ObjectParser) (*RegionMultiLevelNpcData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRegionMultiLevelNpcData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RegionMultiLevelNpcData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: TypeData
	// skip field: LevelBases
	dAtA.OffsetBaseX = pArSeR.Uint64("offset_base_x")
	dAtA.OffsetBaseY = pArSeR.Uint64("offset_base_y")
	dAtA.RegionType = shared_proto.RegionType(shared_proto.RegionType_value[strings.ToUpper(pArSeR.String("region_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("region_type"), 10, 32); err == nil {
		dAtA.RegionType = shared_proto.RegionType(i)
	}

	dAtA.RegionLevel = pArSeR.Uint64("region_level")

	// calculate fields
	dAtA.RegionId = RegionDataID(dAtA.RegionType, dAtA.RegionLevel)

	return dAtA, nil
}

var vAlIdAtOrRegionMultiLevelNpcData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"type_data":     config.ParseValidator("string", "", false, nil, nil),
	"offset_base_x": config.ParseValidator("int>0", "", false, nil, nil),
	"offset_base_y": config.ParseValidator("int>0", "", false, nil, nil),
	"region_type":   config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.RegionType_value, 0), nil),
	"region_level":  config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *RegionMultiLevelNpcData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *RegionMultiLevelNpcData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *RegionMultiLevelNpcData) Encode() *shared_proto.RegionMultiLevelNpcDataProto {
	out := &shared_proto.RegionMultiLevelNpcDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.TypeData != nil {
		out.Type = dAtA.TypeData.Type
	}
	if dAtA.LevelBases != nil {
		out.LevelBaseId = encodeMultiLevelNpcLevelData(dAtA.LevelBases)
	}
	out.OffsetBaseX = config.U64ToI32(dAtA.OffsetBaseX)
	out.OffsetBaseY = config.U64ToI32(dAtA.OffsetBaseY)

	return out
}

func ArrayEncodeRegionMultiLevelNpcData(datas []*RegionMultiLevelNpcData) []*shared_proto.RegionMultiLevelNpcDataProto {

	out := make([]*shared_proto.RegionMultiLevelNpcDataProto, 0, len(datas))
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

func (dAtA *RegionMultiLevelNpcData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.TypeData = cOnFigS.GetRegionMultiLevelNpcTypeData(pArSeR.Int("type_data"))
	if dAtA.TypeData == nil {
		return errors.Errorf("%s 配置的关联字段[type_data] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("type_data"), *pArSeR)
	}

	return nil
}

// start with RegionMultiLevelNpcLevelData ----------------------------------

func LoadRegionMultiLevelNpcLevelData(gos *config.GameObjects) (map[uint64]*RegionMultiLevelNpcLevelData, map[*RegionMultiLevelNpcLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.RegionMultiLevelNpcLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*RegionMultiLevelNpcLevelData, len(lIsT))
	pArSeRmAp := make(map[*RegionMultiLevelNpcLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrRegionMultiLevelNpcLevelData) {
			continue
		}

		dAtA, err := NewRegionMultiLevelNpcLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedRegionMultiLevelNpcLevelData(dAtAmAp map[*RegionMultiLevelNpcLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RegionMultiLevelNpcLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetRegionMultiLevelNpcLevelDataKeyArray(datas []*RegionMultiLevelNpcLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewRegionMultiLevelNpcLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*RegionMultiLevelNpcLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRegionMultiLevelNpcLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RegionMultiLevelNpcLevelData{}

	dAtA.MultiLevelNpcId = pArSeR.Uint64("multi_level_npc_id")
	dAtA.Level = pArSeR.Uint64("level")
	// releated field: Npc
	dAtA.HateTickDuration, err = config.ParseDuration(pArSeR.String("hate_tick_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[hate_tick_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("hate_tick_duration"), dAtA)
	}

	dAtA.TickHate = pArSeR.Uint64("tick_hate")
	dAtA.FirstHate = pArSeR.Uint64("first_hate")
	dAtA.FailHate = pArSeR.Int("fail_hate")
	// releated field: FightNpc

	// calculate fields
	dAtA.Id = dAtA.MultiLevelNpcId*10000 + dAtA.Level

	return dAtA, nil
}

var vAlIdAtOrRegionMultiLevelNpcLevelData = map[string]*config.Validator{

	"multi_level_npc_id": config.ParseValidator("int>0", "", false, nil, nil),
	"level":              config.ParseValidator("int>0", "", false, nil, nil),
	"npc":                config.ParseValidator("string", "", false, nil, nil),
	"hate_tick_duration": config.ParseValidator("string", "", false, nil, nil),
	"tick_hate":          config.ParseValidator("uint", "", false, nil, nil),
	"first_hate":         config.ParseValidator("uint", "", false, nil, nil),
	"fail_hate":          config.ParseValidator("int", "", false, nil, nil),
	"fight_npc":          config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *RegionMultiLevelNpcLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Npc = cOnFigS.GetNpcBaseData(pArSeR.Uint64("npc"))
	if dAtA.Npc == nil {
		return errors.Errorf("%s 配置的关联字段[npc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("npc"), *pArSeR)
	}

	dAtA.FightNpc = cOnFigS.GetMonsterMasterData(pArSeR.Uint64("fight_npc"))
	if dAtA.FightNpc == nil {
		return errors.Errorf("%s 配置的关联字段[fight_npc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fight_npc"), *pArSeR)
	}

	return nil
}

// start with RegionMultiLevelNpcTypeData ----------------------------------

func LoadRegionMultiLevelNpcTypeData(gos *config.GameObjects) (map[int]*RegionMultiLevelNpcTypeData, map[*RegionMultiLevelNpcTypeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.RegionMultiLevelNpcTypeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[int]*RegionMultiLevelNpcTypeData, len(lIsT))
	pArSeRmAp := make(map[*RegionMultiLevelNpcTypeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrRegionMultiLevelNpcTypeData) {
			continue
		}

		dAtA, err := NewRegionMultiLevelNpcTypeData(fIlEnAmE, pArSeR)
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

func SetRelatedRegionMultiLevelNpcTypeData(dAtAmAp map[*RegionMultiLevelNpcTypeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RegionMultiLevelNpcTypeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetRegionMultiLevelNpcTypeDataKeyArray(datas []*RegionMultiLevelNpcTypeData) []int {

	out := make([]int, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewRegionMultiLevelNpcTypeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*RegionMultiLevelNpcTypeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRegionMultiLevelNpcTypeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RegionMultiLevelNpcTypeData{}

	dAtA.Type = shared_proto.MultiLevelNpcType(shared_proto.MultiLevelNpcType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.MultiLevelNpcType(i)
	}

	dAtA.InitHate = pArSeR.Uint64("init_hate")
	dAtA.MaxHate = pArSeR.Uint64("max_hate")
	dAtA.FightHate = pArSeR.Uint64("fight_hate")
	dAtA.FightReduceHate = pArSeR.Uint64("fight_reduce_hate")
	dAtA.FightMustDistance = pArSeR.Uint64("fight_must_distance")
	if pArSeR.KeyExist("fight_delay") {
		dAtA.FightDelay, err = config.ParseDuration(pArSeR.String("fight_delay"))
	} else {
		dAtA.FightDelay, err = config.ParseDuration("30s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[fight_delay] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("fight_delay"), dAtA)
	}

	// calculate fields
	dAtA.Id = int(dAtA.Type)

	// i18n fields
	dAtA.Name = i18n.NewI18nRef(fIlEnAmE, "name", dAtA.Id, pArSeR.String("name"))

	return dAtA, nil
}

var vAlIdAtOrRegionMultiLevelNpcTypeData = map[string]*config.Validator{

	"name":                config.ParseValidator("string", "", false, nil, nil),
	"type":                config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.MultiLevelNpcType_value, 0), nil),
	"init_hate":           config.ParseValidator("uint", "", false, nil, nil),
	"max_hate":            config.ParseValidator("int>0", "", false, nil, nil),
	"fight_hate":          config.ParseValidator("int>0", "", false, nil, nil),
	"fight_reduce_hate":   config.ParseValidator("int>0", "", false, nil, nil),
	"fight_must_distance": config.ParseValidator("int>0", "", false, nil, nil),
	"fight_delay":         config.ParseValidator("string", "", false, nil, []string{"30s"}),
}

func (dAtA *RegionMultiLevelNpcTypeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *RegionMultiLevelNpcTypeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *RegionMultiLevelNpcTypeData) Encode() *shared_proto.RegionMultiLevelNpcTypeProto {
	out := &shared_proto.RegionMultiLevelNpcTypeProto{}
	if dAtA.Name != nil {
		out.Name = dAtA.Name.Encode()
	}
	out.Type = dAtA.Type
	out.MaxHate = config.U64ToI32(dAtA.MaxHate)
	out.FightHate = config.U64ToI32(dAtA.FightHate)

	return out
}

func ArrayEncodeRegionMultiLevelNpcTypeData(datas []*RegionMultiLevelNpcTypeData) []*shared_proto.RegionMultiLevelNpcTypeProto {

	out := make([]*shared_proto.RegionMultiLevelNpcTypeProto, 0, len(datas))
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

func (dAtA *RegionMultiLevelNpcTypeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with TroopDialogueData ----------------------------------

func LoadTroopDialogueData(gos *config.GameObjects) (map[uint64]*TroopDialogueData, map[*TroopDialogueData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TroopDialogueDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TroopDialogueData, len(lIsT))
	pArSeRmAp := make(map[*TroopDialogueData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTroopDialogueData) {
			continue
		}

		dAtA, err := NewTroopDialogueData(fIlEnAmE, pArSeR)
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

func SetRelatedTroopDialogueData(dAtAmAp map[*TroopDialogueData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TroopDialogueDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTroopDialogueDataKeyArray(datas []*TroopDialogueData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTroopDialogueData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TroopDialogueData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTroopDialogueData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TroopDialogueData{}

	dAtA.BaseTargetType = shared_proto.BaseTargetType(shared_proto.BaseTargetType_value[strings.ToUpper(pArSeR.String("base_target_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("base_target_type"), 10, 32); err == nil {
		dAtA.BaseTargetType = shared_proto.BaseTargetType(i)
	}

	dAtA.BaseTargetSubType = pArSeR.Uint64("base_target_sub_type")
	dAtA.BayeStage = pArSeR.Uint64("baye_stage")
	dAtA.HeroLevel = pArSeR.Uint64("hero_level")
	dAtA.FirstDelay, err = config.ParseDuration(pArSeR.String("first_delay"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[first_delay] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("first_delay"), dAtA)
	}

	dAtA.NextDelay, err = config.ParseDuration(pArSeR.String("next_delay"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[next_delay] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("next_delay"), dAtA)
	}

	dAtA.RandomText = pArSeR.Bool("random_text")
	// releated field: Texts

	// calculate fields
	dAtA.Id = GetTroopDialogueId(dAtA.BaseTargetType, dAtA.BaseTargetSubType)

	return dAtA, nil
}

var vAlIdAtOrTroopDialogueData = map[string]*config.Validator{

	"base_target_type":     config.ParseValidator("string", "", false, config.EnumMapKeys(shared_proto.BaseTargetType_value), nil),
	"base_target_sub_type": config.ParseValidator("uint", "", false, nil, nil),
	"baye_stage":           config.ParseValidator("uint", "", false, nil, nil),
	"hero_level":           config.ParseValidator("uint", "", false, nil, nil),
	"first_delay":          config.ParseValidator("string", "", false, nil, nil),
	"next_delay":           config.ParseValidator("string", "", false, nil, nil),
	"random_text":          config.ParseValidator("bool", "", false, nil, nil),
	"texts":                config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *TroopDialogueData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TroopDialogueData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TroopDialogueData) Encode() *shared_proto.TroopDialogueDataProto {
	out := &shared_proto.TroopDialogueDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.BaseTargetType = dAtA.BaseTargetType
	out.BaseTargetSubType = config.U64ToI32(dAtA.BaseTargetSubType)
	out.FirstDelay = config.Duration2I32Seconds(dAtA.FirstDelay)
	out.NextDelay = config.Duration2I32Seconds(dAtA.NextDelay)
	out.RandomText = dAtA.RandomText
	if dAtA.Texts != nil {
		out.Texts = config.U64a2I32a(GetTroopDialogueTextDataKeyArray(dAtA.Texts))
	}

	return out
}

func ArrayEncodeTroopDialogueData(datas []*TroopDialogueData) []*shared_proto.TroopDialogueDataProto {

	out := make([]*shared_proto.TroopDialogueDataProto, 0, len(datas))
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

func (dAtA *TroopDialogueData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("texts", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetTroopDialogueTextData(v)
		if obj != nil {
			dAtA.Texts = append(dAtA.Texts, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[texts] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("texts"), *pArSeR)
		}
	}

	return nil
}

// start with TroopDialogueTextData ----------------------------------

func LoadTroopDialogueTextData(gos *config.GameObjects) (map[uint64]*TroopDialogueTextData, map[*TroopDialogueTextData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TroopDialogueTextDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TroopDialogueTextData, len(lIsT))
	pArSeRmAp := make(map[*TroopDialogueTextData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTroopDialogueTextData) {
			continue
		}

		dAtA, err := NewTroopDialogueTextData(fIlEnAmE, pArSeR)
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

func SetRelatedTroopDialogueTextData(dAtAmAp map[*TroopDialogueTextData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TroopDialogueTextDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTroopDialogueTextDataKeyArray(datas []*TroopDialogueTextData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTroopDialogueTextData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TroopDialogueTextData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTroopDialogueTextData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TroopDialogueTextData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Head = pArSeR.String("head")
	dAtA.Text = pArSeR.String("text")
	dAtA.Duration, err = config.ParseDuration(pArSeR.String("duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("duration"), dAtA)
	}

	dAtA.Direction = pArSeR.Uint64("direction")

	return dAtA, nil
}

var vAlIdAtOrTroopDialogueTextData = map[string]*config.Validator{

	"id":        config.ParseValidator("int>0", "", false, nil, nil),
	"head":      config.ParseValidator("string", "", false, nil, nil),
	"text":      config.ParseValidator("string", "", false, nil, nil),
	"duration":  config.ParseValidator("string", "", false, nil, nil),
	"direction": config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *TroopDialogueTextData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TroopDialogueTextData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TroopDialogueTextData) Encode() *shared_proto.TroopDialogueTextDataProto {
	out := &shared_proto.TroopDialogueTextDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Head = dAtA.Head
	out.Text = dAtA.Text
	out.Duration = config.Duration2I32Seconds(dAtA.Duration)
	out.Direction = config.U64ToI32(dAtA.Direction)

	return out
}

func ArrayEncodeTroopDialogueTextData(datas []*TroopDialogueTextData) []*shared_proto.TroopDialogueTextDataProto {

	out := make([]*shared_proto.TroopDialogueTextDataProto, 0, len(datas))
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

func (dAtA *TroopDialogueTextData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetAreaData(uint64) *AreaData
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
	GetNpcBaseData(uint64) *basedata.NpcBaseData
	GetRegionMultiLevelNpcTypeData(int) *RegionMultiLevelNpcTypeData
	GetTroopDialogueTextData(uint64) *TroopDialogueTextData
}
