// AUTO_GEN, DONT MODIFY!!!
package mingcdata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/pb/server_proto"
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

// start with McBuildAddSupportData ----------------------------------

func LoadMcBuildAddSupportData(gos *config.GameObjects) (map[uint64]*McBuildAddSupportData, map[*McBuildAddSupportData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.McBuildAddSupportDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*McBuildAddSupportData, len(lIsT))
	pArSeRmAp := make(map[*McBuildAddSupportData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMcBuildAddSupportData) {
			continue
		}

		dAtA, err := NewMcBuildAddSupportData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.BaiZhanLevel
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[BaiZhanLevel], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedMcBuildAddSupportData(dAtAmAp map[*McBuildAddSupportData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.McBuildAddSupportDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMcBuildAddSupportDataKeyArray(datas []*McBuildAddSupportData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.BaiZhanLevel)
		}
	}

	return out
}

func NewMcBuildAddSupportData(fIlEnAmE string, pArSeR *config.ObjectParser) (*McBuildAddSupportData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMcBuildAddSupportData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &McBuildAddSupportData{}

	dAtA.BaiZhanLevel = pArSeR.Uint64("bai_zhan_level")
	dAtA.AddSupport = pArSeR.Uint64("add_support")

	return dAtA, nil
}

var vAlIdAtOrMcBuildAddSupportData = map[string]*config.Validator{

	"bai_zhan_level": config.ParseValidator("int>0", "", false, nil, nil),
	"add_support":    config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *McBuildAddSupportData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *McBuildAddSupportData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *McBuildAddSupportData) Encode() *shared_proto.McBuildAddSupportDataProto {
	out := &shared_proto.McBuildAddSupportDataProto{}
	out.BaiZhanLevel = config.U64ToI32(dAtA.BaiZhanLevel)
	out.AddSupport = config.U64ToI32(dAtA.AddSupport)

	return out
}

func ArrayEncodeMcBuildAddSupportData(datas []*McBuildAddSupportData) []*shared_proto.McBuildAddSupportDataProto {

	out := make([]*shared_proto.McBuildAddSupportDataProto, 0, len(datas))
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

func (dAtA *McBuildAddSupportData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with McBuildGuildMemberPrizeData ----------------------------------

func LoadMcBuildGuildMemberPrizeData(gos *config.GameObjects) (map[uint64]*McBuildGuildMemberPrizeData, map[*McBuildGuildMemberPrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.McBuildGuildMemberPrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*McBuildGuildMemberPrizeData, len(lIsT))
	pArSeRmAp := make(map[*McBuildGuildMemberPrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMcBuildGuildMemberPrizeData) {
			continue
		}

		dAtA, err := NewMcBuildGuildMemberPrizeData(fIlEnAmE, pArSeR)
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

func SetRelatedMcBuildGuildMemberPrizeData(dAtAmAp map[*McBuildGuildMemberPrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.McBuildGuildMemberPrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMcBuildGuildMemberPrizeDataKeyArray(datas []*McBuildGuildMemberPrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMcBuildGuildMemberPrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*McBuildGuildMemberPrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMcBuildGuildMemberPrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &McBuildGuildMemberPrizeData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.MinBuildCount = pArSeR.Uint64("min_build_count")
	dAtA.MaxBuildCount = pArSeR.Uint64("max_build_count")
	// skip field: BuildCountRage
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrMcBuildGuildMemberPrizeData = map[string]*config.Validator{

	"id":              config.ParseValidator("int>0", "", false, nil, nil),
	"min_build_count": config.ParseValidator("int>0", "", false, nil, nil),
	"max_build_count": config.ParseValidator("int>0", "", false, nil, nil),
	"prize":           config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *McBuildGuildMemberPrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *McBuildGuildMemberPrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *McBuildGuildMemberPrizeData) Encode() *shared_proto.McBuildGuildMemberPrizeDataProto {
	out := &shared_proto.McBuildGuildMemberPrizeDataProto{}
	out.MinBuildCount = config.U64ToI32(dAtA.MinBuildCount)
	out.MaxBuildCount = config.U64ToI32(dAtA.MaxBuildCount)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeMcBuildGuildMemberPrizeData(datas []*McBuildGuildMemberPrizeData) []*shared_proto.McBuildGuildMemberPrizeDataProto {

	out := make([]*shared_proto.McBuildGuildMemberPrizeDataProto, 0, len(datas))
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

func (dAtA *McBuildGuildMemberPrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with McBuildMcSupportData ----------------------------------

func LoadMcBuildMcSupportData(gos *config.GameObjects) (map[uint64]*McBuildMcSupportData, map[*McBuildMcSupportData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.McBuildMcSupportDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*McBuildMcSupportData, len(lIsT))
	pArSeRmAp := make(map[*McBuildMcSupportData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMcBuildMcSupportData) {
			continue
		}

		dAtA, err := NewMcBuildMcSupportData(fIlEnAmE, pArSeR)
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

func SetRelatedMcBuildMcSupportData(dAtAmAp map[*McBuildMcSupportData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.McBuildMcSupportDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMcBuildMcSupportDataKeyArray(datas []*McBuildMcSupportData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewMcBuildMcSupportData(fIlEnAmE string, pArSeR *config.ObjectParser) (*McBuildMcSupportData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMcBuildMcSupportData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &McBuildMcSupportData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.UpgradeSupport = pArSeR.Uint64("upgrade_support")
	dAtA.AddDailyYinliang = pArSeR.Uint64("add_daily_yinliang")
	dAtA.AddMaxYinliang = pArSeR.Uint64("add_max_yinliang")
	dAtA.AddHostDailyYinliang = pArSeR.Uint64("add_host_daily_yinliang")
	// skip field: PrevLevelData
	// skip field: NextLevelData

	return dAtA, nil
}

var vAlIdAtOrMcBuildMcSupportData = map[string]*config.Validator{

	"level":                   config.ParseValidator("int>0", "", false, nil, nil),
	"upgrade_support":         config.ParseValidator("int>0", "", false, nil, nil),
	"add_daily_yinliang":      config.ParseValidator("int>0", "", false, nil, nil),
	"add_max_yinliang":        config.ParseValidator("int>0", "", false, nil, nil),
	"add_host_daily_yinliang": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *McBuildMcSupportData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *McBuildMcSupportData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *McBuildMcSupportData) Encode() *shared_proto.McBuildMcSupportDataProto {
	out := &shared_proto.McBuildMcSupportDataProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.UpgradeSupport = config.U64ToI32(dAtA.UpgradeSupport)
	out.AddDailyYinliang = config.U64ToI32(dAtA.AddDailyYinliang)
	out.AddMaxYinliang = config.U64ToI32(dAtA.AddMaxYinliang)
	out.AddHostDailyYinliang = config.U64ToI32(dAtA.AddHostDailyYinliang)

	return out
}

func ArrayEncodeMcBuildMcSupportData(datas []*McBuildMcSupportData) []*shared_proto.McBuildMcSupportDataProto {

	out := make([]*shared_proto.McBuildMcSupportDataProto, 0, len(datas))
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

func (dAtA *McBuildMcSupportData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with McBuildMiscData ----------------------------------

func LoadMcBuildMiscData(gos *config.GameObjects) (*McBuildMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.McBuildMiscDataPath
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

	dAtA, err := NewMcBuildMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedMcBuildMiscData(gos *config.GameObjects, dAtA *McBuildMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.McBuildMiscDataPath
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

func NewMcBuildMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*McBuildMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMcBuildMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &McBuildMiscData{}

	dAtA.MaxDailyAddSupport = pArSeR.Uint64("max_daily_add_support")
	dAtA.DailyReduceSupport = pArSeR.Uint64("daily_reduce_support")
	dAtA.BuildCd, err = config.ParseDuration(pArSeR.String("build_cd"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[build_cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("build_cd"), dAtA)
	}

	dAtA.DailyBuildMaxCount = pArSeR.Uint64("daily_build_max_count")
	dAtA.BuildMinHeroLevel = pArSeR.Uint64("build_min_hero_level")
	dAtA.MaxRecommendMcGuildCount = pArSeR.Int("max_recommend_mc_guild_count")
	dAtA.MaxMcBuildLogCount = pArSeR.Int("max_mc_build_log_count")

	return dAtA, nil
}

var vAlIdAtOrMcBuildMiscData = map[string]*config.Validator{

	"max_daily_add_support":        config.ParseValidator("int>0", "", false, nil, nil),
	"daily_reduce_support":         config.ParseValidator("int>0", "", false, nil, nil),
	"build_cd":                     config.ParseValidator("string", "", false, nil, nil),
	"daily_build_max_count":        config.ParseValidator("int>0", "", false, nil, nil),
	"build_min_hero_level":         config.ParseValidator("int>0", "", false, nil, nil),
	"max_recommend_mc_guild_count": config.ParseValidator("int>0", "", false, nil, nil),
	"max_mc_build_log_count":       config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *McBuildMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *McBuildMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *McBuildMiscData) Encode() *shared_proto.McBuildMiscDataProto {
	out := &shared_proto.McBuildMiscDataProto{}
	out.MaxDailyAddSupport = config.U64ToI32(dAtA.MaxDailyAddSupport)
	out.DailyReduceSupport = config.U64ToI32(dAtA.DailyReduceSupport)
	out.BuildCd = config.Duration2I32Seconds(dAtA.BuildCd)
	out.DailyBuildMaxCount = config.U64ToI32(dAtA.DailyBuildMaxCount)
	out.BuildMinHeroLevel = config.U64ToI32(dAtA.BuildMinHeroLevel)
	out.MaxRecommendMcGuildCount = int32(dAtA.MaxRecommendMcGuildCount)
	out.MaxMcBuildLogCount = int32(dAtA.MaxMcBuildLogCount)

	return out
}

func ArrayEncodeMcBuildMiscData(datas []*McBuildMiscData) []*shared_proto.McBuildMiscDataProto {

	out := make([]*shared_proto.McBuildMiscDataProto, 0, len(datas))
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

func (dAtA *McBuildMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with MingcBaseData ----------------------------------

func LoadMingcBaseData(gos *config.GameObjects) (map[uint64]*MingcBaseData, map[*MingcBaseData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcBaseDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcBaseData, len(lIsT))
	pArSeRmAp := make(map[*MingcBaseData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcBaseData) {
			continue
		}

		dAtA, err := NewMingcBaseData(fIlEnAmE, pArSeR)
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

func SetRelatedMingcBaseData(dAtAmAp map[*MingcBaseData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcBaseDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcBaseDataKeyArray(datas []*MingcBaseData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMingcBaseData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcBaseData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcBaseData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcBaseData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Model = pArSeR.String("model")
	dAtA.BaseX = pArSeR.Uint64("base_x")
	dAtA.BaseY = pArSeR.Uint64("base_y")
	dAtA.Radius = pArSeR.Uint64("radius")
	dAtA.Type = shared_proto.MincType(shared_proto.MincType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.MincType(i)
	}

	dAtA.ZhouCaptain = pArSeR.Bool("zhou_captain")
	dAtA.DefaultYinliang = pArSeR.Uint64("default_yinliang")
	dAtA.DailyAddYinliang = pArSeR.Uint64("daily_add_yinliang")
	dAtA.MaxYinliang = pArSeR.Uint64("max_yinliang")
	dAtA.HostDailyAddYinliang = pArSeR.Uint64("host_daily_add_yinliang")
	dAtA.Country = pArSeR.Uint64("country")
	// releated field: WarIcon
	dAtA.AtkMinHufu = pArSeR.Uint64("atk_min_hufu")
	dAtA.AtkMinGuildLevel = pArSeR.Uint64("atk_min_guild_level")
	dAtA.AstMaxGuild = pArSeR.Uint64("ast_max_guild")
	dAtA.BaseMinDistance = pArSeR.Uint64("base_min_distance")

	return dAtA, nil
}

var vAlIdAtOrMingcBaseData = map[string]*config.Validator{

	"id":                      config.ParseValidator("int>0", "", false, nil, nil),
	"name":                    config.ParseValidator("string", "", false, nil, nil),
	"model":                   config.ParseValidator("string", "", false, nil, nil),
	"base_x":                  config.ParseValidator("int>0", "", false, nil, nil),
	"base_y":                  config.ParseValidator("int>0", "", false, nil, nil),
	"radius":                  config.ParseValidator("int>0", "", false, nil, nil),
	"type":                    config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.MincType_value, 0), nil),
	"zhou_captain":            config.ParseValidator("bool", "", false, nil, nil),
	"default_yinliang":        config.ParseValidator("int>0", "", false, nil, nil),
	"daily_add_yinliang":      config.ParseValidator("int>0", "", false, nil, nil),
	"max_yinliang":            config.ParseValidator("int>0", "", false, nil, nil),
	"host_daily_add_yinliang": config.ParseValidator("int>0", "", false, nil, nil),
	"country":                 config.ParseValidator("uint", "", false, nil, nil),
	"war_icon":                config.ParseValidator("string", "", false, nil, []string{"WarIcon"}),
	"atk_min_hufu":            config.ParseValidator("int>0", "", false, nil, nil),
	"atk_min_guild_level":     config.ParseValidator("int>0", "", false, nil, nil),
	"ast_max_guild":           config.ParseValidator("int>0", "", false, nil, nil),
	"base_min_distance":       config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *MingcBaseData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MingcBaseData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MingcBaseData) Encode() *shared_proto.MingcBaseDataProto {
	out := &shared_proto.MingcBaseDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Model = dAtA.Model
	out.BaseX = config.U64ToI32(dAtA.BaseX)
	out.BaseY = config.U64ToI32(dAtA.BaseY)
	out.Radius = config.U64ToI32(dAtA.Radius)
	out.Type = dAtA.Type
	out.ZhouCaptain = dAtA.ZhouCaptain
	out.DefaultYinliang = config.U64ToI32(dAtA.DefaultYinliang)
	out.DailyAddYinliang = config.U64ToI32(dAtA.DailyAddYinliang)
	out.MaxYinliang = config.U64ToI32(dAtA.MaxYinliang)
	out.HostDailyAddYinliang = config.U64ToI32(dAtA.HostDailyAddYinliang)
	out.Country = config.U64ToI32(dAtA.Country)
	if dAtA.WarIcon != nil {
		out.WarIcon = dAtA.WarIcon.Id
	}
	out.AtkMinHufu = config.U64ToI32(dAtA.AtkMinHufu)
	out.AtkMinGuildLevel = config.U64ToI32(dAtA.AtkMinGuildLevel)
	out.AstMaxGuild = config.U64ToI32(dAtA.AstMaxGuild)
	out.BaseMinDistance = config.U64ToI32(dAtA.BaseMinDistance)

	return out
}

func ArrayEncodeMingcBaseData(datas []*MingcBaseData) []*shared_proto.MingcBaseDataProto {

	out := make([]*shared_proto.MingcBaseDataProto, 0, len(datas))
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

func (dAtA *MingcBaseData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("war_icon") {
		dAtA.WarIcon = cOnFigS.GetIcon(pArSeR.String("war_icon"))
	} else {
		dAtA.WarIcon = cOnFigS.GetIcon("WarIcon")
	}
	if dAtA.WarIcon == nil {
		return errors.Errorf("%s 配置的关联字段[war_icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("war_icon"), *pArSeR)
	}

	return nil
}

// start with MingcMiscData ----------------------------------

func LoadMingcMiscData(gos *config.GameObjects) (*MingcMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcMiscDataPath
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

	dAtA, err := NewMingcMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedMingcMiscData(gos *config.GameObjects, dAtA *MingcMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcMiscDataPath
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

func NewMingcMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcMiscData{}

	if pArSeR.KeyExist("fight_prepare_duration") {
		dAtA.FightPrepareDuration, err = config.ParseDuration(pArSeR.String("fight_prepare_duration"))
	} else {
		dAtA.FightPrepareDuration, err = config.ParseDuration("30m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[fight_prepare_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("fight_prepare_duration"), dAtA)
	}

	if pArSeR.KeyExist("join_fight_duration") {
		dAtA.JoinFightDuration, err = config.ParseDuration(pArSeR.String("join_fight_duration"))
	} else {
		dAtA.JoinFightDuration, err = config.ParseDuration("180s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[join_fight_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("join_fight_duration"), dAtA)
	}

	dAtA.ApplyAstLimit = 3
	if pArSeR.KeyExist("apply_ast_limit") {
		dAtA.ApplyAstLimit = pArSeR.Uint64("apply_ast_limit")
	}

	dAtA.StartAfterServerOpen, err = config.ParseDuration(pArSeR.String("start_after_server_open"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[start_after_server_open] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("start_after_server_open"), dAtA)
	}

	dAtA.StartSelfCapitalAfterServerOpen, err = config.ParseDuration(pArSeR.String("start_self_capital_after_server_open"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[start_self_capital_after_server_open] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("start_self_capital_after_server_open"), dAtA)
	}

	dAtA.StartOtherCapitalAfterServerOpen, err = config.ParseDuration(pArSeR.String("start_other_capital_after_server_open"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[start_other_capital_after_server_open] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("start_other_capital_after_server_open"), dAtA)
	}

	dAtA.StartSelfCapitalNoticeAfterServerOpen, err = config.ParseDuration(pArSeR.String("start_self_capital_notice_after_server_open"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[start_self_capital_notice_after_server_open] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("start_self_capital_notice_after_server_open"), dAtA)
	}

	dAtA.StartOtherCapitalNoticeAfterServerOpen, err = config.ParseDuration(pArSeR.String("start_other_capital_notice_after_server_open"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[start_other_capital_notice_after_server_open] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("start_other_capital_notice_after_server_open"), dAtA)
	}

	dAtA.DestroyProsperityMaxTroop = 5
	if pArSeR.KeyExist("destroy_prosperity_max_troop") {
		dAtA.DestroyProsperityMaxTroop = pArSeR.Uint64("destroy_prosperity_max_troop")
	}

	dAtA.PerDestroyProsperity = pArSeR.Uint64("per_destroy_prosperity")
	dAtA.DestroyProsperityDuration, err = config.ParseDuration(pArSeR.String("destroy_prosperity_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[destroy_prosperity_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("destroy_prosperity_duration"), dAtA)
	}

	dAtA.ReliveDuration, err = config.ParseDuration(pArSeR.String("relive_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[relive_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("relive_duration"), dAtA)
	}

	// releated field: WallStat
	dAtA.WallFixDamage = 0
	if pArSeR.KeyExist("wall_fix_damage") {
		dAtA.WallFixDamage = pArSeR.Uint64("wall_fix_damage")
	}

	dAtA.WallLevel = 1
	if pArSeR.KeyExist("wall_level") {
		dAtA.WallLevel = pArSeR.Uint64("wall_level")
	}

	dAtA.SceneRecordMaxLen = 50
	if pArSeR.KeyExist("scene_record_max_len") {
		dAtA.SceneRecordMaxLen = pArSeR.Uint64("scene_record_max_len")
	}

	dAtA.Speed = 1
	if pArSeR.KeyExist("speed") {
		dAtA.Speed = pArSeR.Float64("speed")
	}

	dAtA.CloseDuCheng = true
	if pArSeR.KeyExist("close_du_cheng") {
		dAtA.CloseDuCheng = pArSeR.Bool("close_du_cheng")
	}

	dAtA.JoinFightHeroMinLevel = 10
	if pArSeR.KeyExist("join_fight_hero_min_level") {
		dAtA.JoinFightHeroMinLevel = pArSeR.Uint64("join_fight_hero_min_level")
	}

	dAtA.RedPointMinGuildLevel = 3
	if pArSeR.KeyExist("red_point_min_guild_level") {
		dAtA.RedPointMinGuildLevel = pArSeR.Uint64("red_point_min_guild_level")
	}

	dAtA.SaveHeroRecordMaxDays = 7
	if pArSeR.KeyExist("save_hero_record_max_days") {
		dAtA.SaveHeroRecordMaxDays = pArSeR.Uint64("save_hero_record_max_days")
	}

	if pArSeR.KeyExist("daily_update_mingc_time") {
		dAtA.DailyUpdateMingcTime, err = config.ParseDuration(pArSeR.String("daily_update_mingc_time"))
	} else {
		dAtA.DailyUpdateMingcTime, err = config.ParseDuration("22h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[daily_update_mingc_time] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("daily_update_mingc_time"), dAtA)
	}

	dAtA.FreeTankSpeed = 0.7
	if pArSeR.KeyExist("free_tank_speed") {
		dAtA.FreeTankSpeed = pArSeR.Float64("free_tank_speed")
	}

	dAtA.FreeTankPerDestroyProsperity = 300
	if pArSeR.KeyExist("free_tank_per_destroy_prosperity") {
		dAtA.FreeTankPerDestroyProsperity = pArSeR.Uint64("free_tank_per_destroy_prosperity")
	}

	if pArSeR.KeyExist("tou_shi_building_turn_duration") {
		dAtA.TouShiBuildingTurnDuration, err = config.ParseDuration(pArSeR.String("tou_shi_building_turn_duration"))
	} else {
		dAtA.TouShiBuildingTurnDuration, err = config.ParseDuration("5s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tou_shi_building_turn_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tou_shi_building_turn_duration"), dAtA)
	}

	if pArSeR.KeyExist("tou_shi_building_prepare_duration") {
		dAtA.TouShiBuildingPrepareDuration, err = config.ParseDuration(pArSeR.String("tou_shi_building_prepare_duration"))
	} else {
		dAtA.TouShiBuildingPrepareDuration, err = config.ParseDuration("5s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tou_shi_building_prepare_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tou_shi_building_prepare_duration"), dAtA)
	}

	dAtA.TouShiBuildingDestroyProsperity = 300
	if pArSeR.KeyExist("tou_shi_building_destroy_prosperity") {
		dAtA.TouShiBuildingDestroyProsperity = pArSeR.Uint64("tou_shi_building_destroy_prosperity")
	}

	dAtA.TouShiBuildingBaseHurt = 100
	if pArSeR.KeyExist("tou_shi_building_base_hurt") {
		dAtA.TouShiBuildingBaseHurt = pArSeR.Uint64("tou_shi_building_base_hurt")
	}

	if pArSeR.KeyExist("tou_shi_building_hurt_percent") {
		dAtA.TouShiBuildingHurtPercent, err = data.ParseAmount(pArSeR.String("tou_shi_building_hurt_percent"))
	} else {
		dAtA.TouShiBuildingHurtPercent, err = data.ParseAmount("10%")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tou_shi_building_hurt_percent] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tou_shi_building_hurt_percent"), dAtA)
	}

	dAtA.TouShiBuildingBaseHurtMaxTroop = 5
	if pArSeR.KeyExist("tou_shi_building_base_hurt_max_troop") {
		dAtA.TouShiBuildingBaseHurtMaxTroop = pArSeR.Uint64("tou_shi_building_base_hurt_max_troop")
	}

	if pArSeR.KeyExist("tou_shi_building_bomb_fly_duration") {
		dAtA.TouShiBuildingBombFlyDuration, err = config.ParseDuration(pArSeR.String("tou_shi_building_bomb_fly_duration"))
	} else {
		dAtA.TouShiBuildingBombFlyDuration, err = config.ParseDuration("4s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tou_shi_building_bomb_fly_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tou_shi_building_bomb_fly_duration"), dAtA)
	}

	if pArSeR.KeyExist("durm_duration") {
		dAtA.DurmDuration, err = config.ParseDuration(pArSeR.String("durm_duration"))
	} else {
		dAtA.DurmDuration, err = config.ParseDuration("300s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[durm_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("durm_duration"), dAtA)
	}

	if pArSeR.KeyExist("drum_stop_duration") {
		dAtA.DrumStopDuration, err = config.ParseDuration(pArSeR.String("drum_stop_duration"))
	} else {
		dAtA.DrumStopDuration, err = config.ParseDuration("10s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[drum_stop_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("drum_stop_duration"), dAtA)
	}

	dAtA.DrumMinBaiZhanLevel = 2
	if pArSeR.KeyExist("drum_min_bai_zhan_level") {
		dAtA.DrumMinBaiZhanLevel = pArSeR.Uint64("drum_min_bai_zhan_level")
	}

	return dAtA, nil
}

var vAlIdAtOrMingcMiscData = map[string]*config.Validator{

	"fight_prepare_duration":                       config.ParseValidator("string", "", false, nil, []string{"30m"}),
	"join_fight_duration":                          config.ParseValidator("string", "", false, nil, []string{"180s"}),
	"apply_ast_limit":                              config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"start_after_server_open":                      config.ParseValidator("string", "", false, nil, nil),
	"start_self_capital_after_server_open":         config.ParseValidator("string", "", false, nil, nil),
	"start_other_capital_after_server_open":        config.ParseValidator("string", "", false, nil, nil),
	"start_self_capital_notice_after_server_open":  config.ParseValidator("string", "", false, nil, nil),
	"start_other_capital_notice_after_server_open": config.ParseValidator("string", "", false, nil, nil),
	"destroy_prosperity_max_troop":                 config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"per_destroy_prosperity":                       config.ParseValidator("int>0", "", false, nil, nil),
	"destroy_prosperity_duration":                  config.ParseValidator("string", "", false, nil, nil),
	"relive_duration":                              config.ParseValidator("string", "", false, nil, nil),
	"wall_stat":                                    config.ParseValidator("string", "", false, nil, nil),
	"wall_fix_damage":                              config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"wall_level":                                   config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"scene_record_max_len":                         config.ParseValidator("int>0", "", false, nil, []string{"50"}),
	"speed":                                        config.ParseValidator("float64>0", "", false, nil, []string{"1"}),
	"close_du_cheng":                               config.ParseValidator("bool", "", false, nil, []string{"true"}),
	"join_fight_hero_min_level":                    config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"red_point_min_guild_level":                    config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"save_hero_record_max_days":                    config.ParseValidator("int>0", "", false, nil, []string{"7"}),
	"daily_update_mingc_time":                      config.ParseValidator("string", "", false, nil, []string{"22h"}),
	"free_tank_speed":                              config.ParseValidator("float64>0", "", false, nil, []string{"0.7"}),
	"free_tank_per_destroy_prosperity":             config.ParseValidator("int>0", "", false, nil, []string{"300"}),
	"tou_shi_building_turn_duration":               config.ParseValidator("string", "", false, nil, []string{"5s"}),
	"tou_shi_building_prepare_duration":            config.ParseValidator("string", "", false, nil, []string{"5s"}),
	"tou_shi_building_destroy_prosperity":          config.ParseValidator("int>0", "", false, nil, []string{"300"}),
	"tou_shi_building_base_hurt":                   config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"tou_shi_building_hurt_percent":                config.ParseValidator("string", "", false, nil, []string{"10%"}),
	"tou_shi_building_base_hurt_max_troop":         config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"tou_shi_building_bomb_fly_duration":           config.ParseValidator("string", "", false, nil, []string{"4s"}),
	"durm_duration":                                config.ParseValidator("string", "", false, nil, []string{"300s"}),
	"drum_stop_duration":                           config.ParseValidator("string", "", false, nil, []string{"10s"}),
	"drum_min_bai_zhan_level":                      config.ParseValidator("int>0", "", false, nil, []string{"2"}),
}

func (dAtA *MingcMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MingcMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MingcMiscData) Encode() *shared_proto.MingcMiscDataProto {
	out := &shared_proto.MingcMiscDataProto{}
	out.FightPrepareDuration = config.Duration2I32Seconds(dAtA.FightPrepareDuration)
	out.JoinFightDuration = config.Duration2I32Seconds(dAtA.JoinFightDuration)
	out.ApplyAstLimit = config.U64ToI32(dAtA.ApplyAstLimit)
	out.StartAfterServerOpen = config.Duration2I32Seconds(dAtA.StartAfterServerOpen)
	out.StartSelfCapitalAfterServerOpen = config.Duration2I32Seconds(dAtA.StartSelfCapitalAfterServerOpen)
	out.StartOtherCapitalAfterServerOpen = config.Duration2I32Seconds(dAtA.StartOtherCapitalAfterServerOpen)
	out.StartSelfCapitalNoticeAfterServerOpen = config.Duration2I32Seconds(dAtA.StartSelfCapitalNoticeAfterServerOpen)
	out.StartOtherCapitalNoticeAfterServerOpen = config.Duration2I32Seconds(dAtA.StartOtherCapitalNoticeAfterServerOpen)
	out.DestroyProsperityMaxTroop = config.U64ToI32(dAtA.DestroyProsperityMaxTroop)
	out.PerDestroyProsperity = config.U64ToI32(dAtA.PerDestroyProsperity)
	out.DestroyProsperityDuration = config.Duration2I32Seconds(dAtA.DestroyProsperityDuration)
	out.ReliveDuration = config.Duration2I32Seconds(dAtA.ReliveDuration)
	out.WallLevel = config.U64ToI32(dAtA.WallLevel)
	out.Speed = config.F64ToI32X1000(dAtA.Speed)
	out.CloseDuCheng = dAtA.CloseDuCheng
	out.JoinFightHeroMinLevel = config.U64ToI32(dAtA.JoinFightHeroMinLevel)
	out.DailyUpdateMingcTime = config.Duration2I32Seconds(dAtA.DailyUpdateMingcTime)
	out.FreeTankSpeed = config.F64ToI32X1000(dAtA.FreeTankSpeed)
	out.FreeTankPerDestroyProsperity = config.U64ToI32(dAtA.FreeTankPerDestroyProsperity)
	out.TouShiBuildingTurnDuration = config.Duration2I32Seconds(dAtA.TouShiBuildingTurnDuration)
	out.TouShiBuildingPrepareDuration = config.Duration2I32Seconds(dAtA.TouShiBuildingPrepareDuration)
	out.TouShiBuildingDestroyProsperity = config.U64ToI32(dAtA.TouShiBuildingDestroyProsperity)
	out.TouShiBuildingBaseHurt = config.U64ToI32(dAtA.TouShiBuildingBaseHurt)
	if dAtA.TouShiBuildingHurtPercent != nil {
		out.TouShiBuildingHurtPercent = dAtA.TouShiBuildingHurtPercent.Encode()
	}
	out.TouShiBuildingBaseHurtMaxTroop = config.U64ToI32(dAtA.TouShiBuildingBaseHurtMaxTroop)
	out.TouShiBuildingBombFlyDuration = config.Duration2I32Seconds(dAtA.TouShiBuildingBombFlyDuration)
	out.DurmDuration = config.Duration2I32Seconds(dAtA.DurmDuration)
	out.DrumStopDuration = config.Duration2I32Seconds(dAtA.DrumStopDuration)
	out.DrumMinBaiZhanLevel = config.U64ToI32(dAtA.DrumMinBaiZhanLevel)

	return out
}

func ArrayEncodeMingcMiscData(datas []*MingcMiscData) []*shared_proto.MingcMiscDataProto {

	out := make([]*shared_proto.MingcMiscDataProto, 0, len(datas))
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

func (dAtA *MingcMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.WallStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("wall_stat"))
	if dAtA.WallStat == nil && pArSeR.Uint64("wall_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[wall_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("wall_stat"), *pArSeR)
	}

	return nil
}

// start with MingcTimeData ----------------------------------

func LoadMingcTimeData(gos *config.GameObjects) (map[uint64]*MingcTimeData, map[*MingcTimeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcTimeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcTimeData, len(lIsT))
	pArSeRmAp := make(map[*MingcTimeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcTimeData) {
			continue
		}

		dAtA, err := NewMingcTimeData(fIlEnAmE, pArSeR)
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

func SetRelatedMingcTimeData(dAtAmAp map[*MingcTimeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcTimeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcTimeDataKeyArray(datas []*MingcTimeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMingcTimeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcTimeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcTimeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcTimeData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.ApplyAtkStart = pArSeR.String("apply_atk_start")
	dAtA.ApplyAtkDuration, err = config.ParseDuration(pArSeR.String("apply_atk_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[apply_atk_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("apply_atk_duration"), dAtA)
	}

	dAtA.ApplyAstStart = pArSeR.String("apply_ast_start")
	dAtA.ApplyAstDuration, err = config.ParseDuration(pArSeR.String("apply_ast_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[apply_ast_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("apply_ast_duration"), dAtA)
	}

	dAtA.FightStart = pArSeR.String("fight_start")
	dAtA.FightDuration, err = config.ParseDuration(pArSeR.String("fight_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[fight_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("fight_duration"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrMingcTimeData = map[string]*config.Validator{

	"id":                 config.ParseValidator("int>0", "", false, nil, nil),
	"desc":               config.ParseValidator("string", "", false, nil, nil),
	"apply_atk_start":    config.ParseValidator("string", "", false, nil, nil),
	"apply_atk_duration": config.ParseValidator("string", "", false, nil, nil),
	"apply_ast_start":    config.ParseValidator("string", "", false, nil, nil),
	"apply_ast_duration": config.ParseValidator("string", "", false, nil, nil),
	"fight_start":        config.ParseValidator("string", "", false, nil, nil),
	"fight_duration":     config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *MingcTimeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MingcTimeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MingcTimeData) Encode() *shared_proto.MingcTimeDataProto {
	out := &shared_proto.MingcTimeDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Desc = dAtA.Desc
	out.ApplyAtkStart = dAtA.ApplyAtkStart
	out.ApplyAtkDuration = config.Duration2I32Seconds(dAtA.ApplyAtkDuration)
	out.ApplyAstStart = dAtA.ApplyAstStart
	out.ApplyAstDuration = config.Duration2I32Seconds(dAtA.ApplyAstDuration)
	out.FightStart = dAtA.FightStart
	out.FightDuration = config.Duration2I32Seconds(dAtA.FightDuration)

	return out
}

func ArrayEncodeMingcTimeData(datas []*MingcTimeData) []*shared_proto.MingcTimeDataProto {

	out := make([]*shared_proto.MingcTimeDataProto, 0, len(datas))
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

func (dAtA *MingcTimeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with MingcWarBuildingData ----------------------------------

func LoadMingcWarBuildingData(gos *config.GameObjects) (map[uint64]*MingcWarBuildingData, map[*MingcWarBuildingData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcWarBuildingDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcWarBuildingData, len(lIsT))
	pArSeRmAp := make(map[*MingcWarBuildingData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcWarBuildingData) {
			continue
		}

		dAtA, err := NewMingcWarBuildingData(fIlEnAmE, pArSeR)
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

func SetRelatedMingcWarBuildingData(dAtAmAp map[*MingcWarBuildingData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcWarBuildingDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcWarBuildingDataKeyArray(datas []*MingcWarBuildingData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMingcWarBuildingData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcWarBuildingData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcWarBuildingData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcWarBuildingData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Type = shared_proto.MingcWarBuildingType(shared_proto.MingcWarBuildingType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.MingcWarBuildingType(i)
	}

	dAtA.Prosperity = pArSeR.Uint64("prosperity")
	dAtA.AtkModel = pArSeR.String("atk_model")
	dAtA.DefModel = pArSeR.String("def_model")
	dAtA.WallAtk = pArSeR.Bool("wall_atk")
	dAtA.CanBeAtked = pArSeR.Bool("can_be_atked")
	// releated field: CombatScene

	return dAtA, nil
}

var vAlIdAtOrMingcWarBuildingData = map[string]*config.Validator{

	"id":           config.ParseValidator("int>0", "", false, nil, nil),
	"type":         config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.MingcWarBuildingType_value, 0), nil),
	"prosperity":   config.ParseValidator("uint", "", false, nil, nil),
	"atk_model":    config.ParseValidator("string", "", false, nil, nil),
	"def_model":    config.ParseValidator("string", "", false, nil, nil),
	"wall_atk":     config.ParseValidator("bool", "", false, nil, nil),
	"can_be_atked": config.ParseValidator("bool", "", false, nil, nil),
	"combat_scene": config.ParseValidator("string", "", false, nil, []string{"CombatScene"}),
}

func (dAtA *MingcWarBuildingData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MingcWarBuildingData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MingcWarBuildingData) Encode() *shared_proto.MingcWarBuildingDataProto {
	out := &shared_proto.MingcWarBuildingDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Type = dAtA.Type
	out.Prosperity = config.U64ToI32(dAtA.Prosperity)
	out.AtkModel = dAtA.AtkModel
	out.DefModel = dAtA.DefModel
	out.WallAtk = dAtA.WallAtk
	out.CanBeAtked = dAtA.CanBeAtked
	if dAtA.CombatScene != nil {
		out.CombatScene = dAtA.CombatScene.Id
	}

	return out
}

func ArrayEncodeMingcWarBuildingData(datas []*MingcWarBuildingData) []*shared_proto.MingcWarBuildingDataProto {

	out := make([]*shared_proto.MingcWarBuildingDataProto, 0, len(datas))
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

func (dAtA *MingcWarBuildingData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("combat_scene") {
		dAtA.CombatScene = cOnFigS.GetCombatScene(pArSeR.String("combat_scene"))
	} else {
		dAtA.CombatScene = cOnFigS.GetCombatScene("CombatScene")
	}
	if dAtA.CombatScene == nil {
		return errors.Errorf("%s 配置的关联字段[combat_scene] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("combat_scene"), *pArSeR)
	}

	return nil
}

// start with MingcWarDrumStatData ----------------------------------

func LoadMingcWarDrumStatData(gos *config.GameObjects) (map[uint64]*MingcWarDrumStatData, map[*MingcWarDrumStatData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcWarDrumStatDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcWarDrumStatData, len(lIsT))
	pArSeRmAp := make(map[*MingcWarDrumStatData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcWarDrumStatData) {
			continue
		}

		dAtA, err := NewMingcWarDrumStatData(fIlEnAmE, pArSeR)
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

func SetRelatedMingcWarDrumStatData(dAtAmAp map[*MingcWarDrumStatData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcWarDrumStatDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcWarDrumStatDataKeyArray(datas []*MingcWarDrumStatData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMingcWarDrumStatData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcWarDrumStatData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcWarDrumStatData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcWarDrumStatData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.DrumMin = pArSeR.Uint64("drum_min")
	dAtA.DrumMax = pArSeR.Uint64("drum_max")
	// skip field: DrumRange
	dAtA.DrumDamageIncreDesc = pArSeR.String("drum_damage_incre_desc")
	dAtA.DrumDamageIncreMin = pArSeR.Uint64("drum_damage_incre_min")
	dAtA.DrumDamageIncreMax = pArSeR.Uint64("drum_damage_incre_max")
	// skip field: DrumDamageIncreRange
	dAtA.DrumDamageDecreDesc = pArSeR.String("drum_damage_decre_desc")
	dAtA.DrumDamageDecreMin = pArSeR.Uint64("drum_damage_decre_min")
	dAtA.DrumDamageDecreMax = pArSeR.Uint64("drum_damage_decre_max")
	// skip field: DrumDamageDecreRange

	return dAtA, nil
}

var vAlIdAtOrMingcWarDrumStatData = map[string]*config.Validator{

	"id":                     config.ParseValidator("int>0", "", false, nil, nil),
	"drum_min":               config.ParseValidator("int>0", "", false, nil, nil),
	"drum_max":               config.ParseValidator("int>0", "", false, nil, nil),
	"drum_damage_incre_desc": config.ParseValidator("string", "", false, nil, nil),
	"drum_damage_incre_min":  config.ParseValidator("int>0", "", false, nil, nil),
	"drum_damage_incre_max":  config.ParseValidator("int>0", "", false, nil, nil),
	"drum_damage_decre_desc": config.ParseValidator("string", "", false, nil, nil),
	"drum_damage_decre_min":  config.ParseValidator("int>0", "", false, nil, nil),
	"drum_damage_decre_max":  config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *MingcWarDrumStatData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with MingcWarMapData ----------------------------------

func LoadMingcWarMapData(gos *config.GameObjects) (map[uint64]*MingcWarMapData, map[*MingcWarMapData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcWarMapDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcWarMapData, len(lIsT))
	pArSeRmAp := make(map[*MingcWarMapData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcWarMapData) {
			continue
		}

		dAtA, err := NewMingcWarMapData(fIlEnAmE, pArSeR)
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

func SetRelatedMingcWarMapData(dAtAmAp map[*MingcWarMapData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcWarMapDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcWarMapDataKeyArray(datas []*MingcWarMapData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMingcWarMapData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcWarMapData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcWarMapData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcWarMapData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Mingc
	dAtA.StartX = pArSeR.Int("start_x")
	dAtA.StartY = pArSeR.Int("start_y")
	dAtA.DestX = pArSeR.IntArray("dest_x", "", false)
	dAtA.DestY = pArSeR.IntArray("dest_y", "", false)

	return dAtA, nil
}

var vAlIdAtOrMingcWarMapData = map[string]*config.Validator{

	"id":      config.ParseValidator("int>0", "", false, nil, nil),
	"mingc":   config.ParseValidator("string", "", false, nil, nil),
	"start_x": config.ParseValidator("int>0", "", false, nil, nil),
	"start_y": config.ParseValidator("int>0", "", false, nil, nil),
	"dest_x":  config.ParseValidator("int,duplicate", "", true, nil, nil),
	"dest_y":  config.ParseValidator("int,duplicate", "", true, nil, nil),
}

func (dAtA *MingcWarMapData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MingcWarMapData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MingcWarMapData) Encode() *shared_proto.MingcWarMapDataProto {
	out := &shared_proto.MingcWarMapDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.Mingc != nil {
		out.Mingc = config.U64ToI32(dAtA.Mingc.Id)
	}
	out.StartX = int32(dAtA.StartX)
	out.StartY = int32(dAtA.StartY)
	out.DestX = config.Ia2I32a(dAtA.DestX)
	out.DestY = config.Ia2I32a(dAtA.DestY)

	return out
}

func ArrayEncodeMingcWarMapData(datas []*MingcWarMapData) []*shared_proto.MingcWarMapDataProto {

	out := make([]*shared_proto.MingcWarMapDataProto, 0, len(datas))
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

func (dAtA *MingcWarMapData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Mingc = cOnFigS.GetMingcBaseData(pArSeR.Uint64("mingc"))
	if dAtA.Mingc == nil {
		return errors.Errorf("%s 配置的关联字段[mingc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mingc"), *pArSeR)
	}

	return nil
}

// start with MingcWarMultiKillData ----------------------------------

func LoadMingcWarMultiKillData(gos *config.GameObjects) (map[uint64]*MingcWarMultiKillData, map[*MingcWarMultiKillData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcWarMultiKillDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcWarMultiKillData, len(lIsT))
	pArSeRmAp := make(map[*MingcWarMultiKillData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcWarMultiKillData) {
			continue
		}

		dAtA, err := NewMingcWarMultiKillData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.MultiKill
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[MultiKill], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedMingcWarMultiKillData(dAtAmAp map[*MingcWarMultiKillData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcWarMultiKillDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcWarMultiKillDataKeyArray(datas []*MingcWarMultiKillData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.MultiKill)
		}
	}

	return out
}

func NewMingcWarMultiKillData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcWarMultiKillData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcWarMultiKillData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcWarMultiKillData{}

	dAtA.MultiKill = pArSeR.Uint64("multi_kill")

	return dAtA, nil
}

var vAlIdAtOrMingcWarMultiKillData = map[string]*config.Validator{

	"multi_kill": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *MingcWarMultiKillData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with MingcWarNpcData ----------------------------------

func LoadMingcWarNpcData(gos *config.GameObjects) (map[uint64]*MingcWarNpcData, map[*MingcWarNpcData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcWarNpcDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcWarNpcData, len(lIsT))
	pArSeRmAp := make(map[*MingcWarNpcData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcWarNpcData) {
			continue
		}

		dAtA, err := NewMingcWarNpcData(fIlEnAmE, pArSeR)
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

func SetRelatedMingcWarNpcData(dAtAmAp map[*MingcWarNpcData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcWarNpcDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcWarNpcDataKeyArray(datas []*MingcWarNpcData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMingcWarNpcData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcWarNpcData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcWarNpcData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcWarNpcData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Mingc
	// releated field: Npc
	dAtA.Def = pArSeR.Bool("def")
	dAtA.AstDef = pArSeR.Bool("ast_def")
	// releated field: Guild
	dAtA.AiType = server_proto.McWarAiType(server_proto.McWarAiType_value[strings.ToUpper(pArSeR.String("ai_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("ai_type"), 10, 32); err == nil {
		dAtA.AiType = server_proto.McWarAiType(i)
	}

	dAtA.BaiZhanLevel = pArSeR.Uint64("bai_zhan_level")

	return dAtA, nil
}

var vAlIdAtOrMingcWarNpcData = map[string]*config.Validator{

	"id":             config.ParseValidator("int>0", "", false, nil, nil),
	"mingc":          config.ParseValidator("string", "", false, nil, nil),
	"npc":            config.ParseValidator("string", "", false, nil, nil),
	"def":            config.ParseValidator("bool", "", false, nil, nil),
	"ast_def":        config.ParseValidator("bool", "", false, nil, nil),
	"guild":          config.ParseValidator("string", "", false, nil, nil),
	"ai_type":        config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(server_proto.McWarAiType_value, 0), nil),
	"bai_zhan_level": config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *MingcWarNpcData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Mingc = cOnFigS.GetMingcBaseData(pArSeR.Uint64("mingc"))
	if dAtA.Mingc == nil {
		return errors.Errorf("%s 配置的关联字段[mingc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mingc"), *pArSeR)
	}

	dAtA.Npc = cOnFigS.GetMonsterMasterData(pArSeR.Uint64("npc"))
	if dAtA.Npc == nil {
		return errors.Errorf("%s 配置的关联字段[npc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("npc"), *pArSeR)
	}

	dAtA.Guild = cOnFigS.GetMingcWarNpcGuildData(pArSeR.Uint64("guild"))
	if dAtA.Guild == nil {
		return errors.Errorf("%s 配置的关联字段[guild] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild"), *pArSeR)
	}

	return nil
}

// start with MingcWarNpcGuildData ----------------------------------

func LoadMingcWarNpcGuildData(gos *config.GameObjects) (map[uint64]*MingcWarNpcGuildData, map[*MingcWarNpcGuildData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcWarNpcGuildDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcWarNpcGuildData, len(lIsT))
	pArSeRmAp := make(map[*MingcWarNpcGuildData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcWarNpcGuildData) {
			continue
		}

		dAtA, err := NewMingcWarNpcGuildData(fIlEnAmE, pArSeR)
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

func SetRelatedMingcWarNpcGuildData(dAtAmAp map[*MingcWarNpcGuildData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcWarNpcGuildDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcWarNpcGuildDataKeyArray(datas []*MingcWarNpcGuildData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMingcWarNpcGuildData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcWarNpcGuildData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcWarNpcGuildData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcWarNpcGuildData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.FlagName = pArSeR.String("flag_name")
	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Country = pArSeR.Uint64("country")

	return dAtA, nil
}

var vAlIdAtOrMingcWarNpcGuildData = map[string]*config.Validator{

	"id":        config.ParseValidator("int>0", "", false, nil, nil),
	"name":      config.ParseValidator("string", "", false, nil, nil),
	"flag_name": config.ParseValidator("string", "", false, nil, nil),
	"level":     config.ParseValidator("int>0", "", false, nil, nil),
	"country":   config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *MingcWarNpcGuildData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with MingcWarSceneData ----------------------------------

func LoadMingcWarSceneData(gos *config.GameObjects) (map[uint64]*MingcWarSceneData, map[*MingcWarSceneData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcWarSceneDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcWarSceneData, len(lIsT))
	pArSeRmAp := make(map[*MingcWarSceneData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcWarSceneData) {
			continue
		}

		dAtA, err := NewMingcWarSceneData(fIlEnAmE, pArSeR)
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

func SetRelatedMingcWarSceneData(dAtAmAp map[*MingcWarSceneData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcWarSceneDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcWarSceneDataKeyArray(datas []*MingcWarSceneData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMingcWarSceneData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcWarSceneData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcWarSceneData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcWarSceneData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.AtkReliveName = pArSeR.String("atk_relive_name")
	dAtA.AtkRelivePosX = pArSeR.Int("atk_relive_pos_x")
	dAtA.AtkRelivePosY = pArSeR.Int("atk_relive_pos_y")
	dAtA.AtkHomeName = pArSeR.String("atk_home_name")
	dAtA.AtkHomePosX = pArSeR.Int("atk_home_pos_x")
	dAtA.AtkHomePosY = pArSeR.Int("atk_home_pos_y")
	dAtA.AtkCastleName = pArSeR.StringArray("atk_castle_name", "", false)
	dAtA.AtkCastlePosX = pArSeR.IntArray("atk_castle_pos_x", "", false)
	dAtA.AtkCastlePosY = pArSeR.IntArray("atk_castle_pos_y", "", false)
	dAtA.AtkGateName = pArSeR.StringArray("atk_gate_name", "", false)
	dAtA.AtkGatePosX = pArSeR.IntArray("atk_gate_pos_x", "", false)
	dAtA.AtkGatePosY = pArSeR.IntArray("atk_gate_pos_y", "", false)
	dAtA.DefReliveName = pArSeR.String("def_relive_name")
	dAtA.DefRelivePosX = pArSeR.Int("def_relive_pos_x")
	dAtA.DefRelivePosY = pArSeR.Int("def_relive_pos_y")
	dAtA.DefHomeName = pArSeR.String("def_home_name")
	dAtA.DefHomePosX = pArSeR.Int("def_home_pos_x")
	dAtA.DefHomePosY = pArSeR.Int("def_home_pos_y")
	dAtA.DefCastleName = pArSeR.StringArray("def_castle_name", "", false)
	dAtA.DefCastlePosX = pArSeR.IntArray("def_castle_pos_x", "", false)
	dAtA.DefCastlePosY = pArSeR.IntArray("def_castle_pos_y", "", false)
	dAtA.DefGateName = pArSeR.StringArray("def_gate_name", "", false)
	dAtA.DefGatePosX = pArSeR.IntArray("def_gate_pos_x", "", false)
	dAtA.DefGatePosY = pArSeR.IntArray("def_gate_pos_y", "", false)
	dAtA.AtkTouShiName = pArSeR.StringArray("atk_tou_shi_name", "", false)
	dAtA.AtkTouShiPosX = pArSeR.IntArray("atk_tou_shi_pos_x", "", false)
	dAtA.AtkTouShiPosY = pArSeR.IntArray("atk_tou_shi_pos_y", "", false)
	dAtA.DefTouShiName = pArSeR.StringArray("def_tou_shi_name", "", false)
	dAtA.DefTouShiPosX = pArSeR.IntArray("def_tou_shi_pos_x", "", false)
	dAtA.DefTouShiPosY = pArSeR.IntArray("def_tou_shi_pos_y", "", false)
	// skip field: AtkFullProsperity
	// skip field: DefFullProsperity

	return dAtA, nil
}

var vAlIdAtOrMingcWarSceneData = map[string]*config.Validator{

	"id":                config.ParseValidator("int>0", "", false, nil, nil),
	"desc":              config.ParseValidator("string", "", false, nil, nil),
	"atk_relive_name":   config.ParseValidator("string", "", false, nil, nil),
	"atk_relive_pos_x":  config.ParseValidator("int>0", "", false, nil, nil),
	"atk_relive_pos_y":  config.ParseValidator("int>0", "", false, nil, nil),
	"atk_home_name":     config.ParseValidator("string", "", false, nil, nil),
	"atk_home_pos_x":    config.ParseValidator("int>0", "", false, nil, nil),
	"atk_home_pos_y":    config.ParseValidator("int>0", "", false, nil, nil),
	"atk_castle_name":   config.ParseValidator("string", "", true, nil, nil),
	"atk_castle_pos_x":  config.ParseValidator("int,duplicate", "", true, nil, nil),
	"atk_castle_pos_y":  config.ParseValidator("int,duplicate", "", true, nil, nil),
	"atk_gate_name":     config.ParseValidator("string", "", true, nil, nil),
	"atk_gate_pos_x":    config.ParseValidator("int,duplicate", "", true, nil, nil),
	"atk_gate_pos_y":    config.ParseValidator("int,duplicate", "", true, nil, nil),
	"def_relive_name":   config.ParseValidator("string", "", false, nil, nil),
	"def_relive_pos_x":  config.ParseValidator("int>0", "", false, nil, nil),
	"def_relive_pos_y":  config.ParseValidator("int>0", "", false, nil, nil),
	"def_home_name":     config.ParseValidator("string", "", false, nil, nil),
	"def_home_pos_x":    config.ParseValidator("int>0", "", false, nil, nil),
	"def_home_pos_y":    config.ParseValidator("int>0", "", false, nil, nil),
	"def_castle_name":   config.ParseValidator("string", "", true, nil, nil),
	"def_castle_pos_x":  config.ParseValidator("int,duplicate", "", true, nil, nil),
	"def_castle_pos_y":  config.ParseValidator("int,duplicate", "", true, nil, nil),
	"def_gate_name":     config.ParseValidator("string", "", true, nil, nil),
	"def_gate_pos_x":    config.ParseValidator("int,duplicate", "", true, nil, nil),
	"def_gate_pos_y":    config.ParseValidator("int,duplicate", "", true, nil, nil),
	"atk_tou_shi_name":  config.ParseValidator("string,duplicate", "", true, nil, nil),
	"atk_tou_shi_pos_x": config.ParseValidator("int,duplicate", "", true, nil, nil),
	"atk_tou_shi_pos_y": config.ParseValidator("int,duplicate", "", true, nil, nil),
	"def_tou_shi_name":  config.ParseValidator("string,duplicate", "", true, nil, nil),
	"def_tou_shi_pos_x": config.ParseValidator("int,duplicate", "", true, nil, nil),
	"def_tou_shi_pos_y": config.ParseValidator("int,duplicate", "", true, nil, nil),
}

func (dAtA *MingcWarSceneData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MingcWarSceneData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MingcWarSceneData) Encode() *shared_proto.MingcWarSceneDataProto {
	out := &shared_proto.MingcWarSceneDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Desc = dAtA.Desc
	out.AtkReliveName = dAtA.AtkReliveName
	out.AtkRelivePosX = int32(dAtA.AtkRelivePosX)
	out.AtkRelivePosY = int32(dAtA.AtkRelivePosY)
	out.AtkHomeName = dAtA.AtkHomeName
	out.AtkHomePosX = int32(dAtA.AtkHomePosX)
	out.AtkHomePosY = int32(dAtA.AtkHomePosY)
	out.AtkCastleName = dAtA.AtkCastleName
	out.AtkCastlePosX = config.Ia2I32a(dAtA.AtkCastlePosX)
	out.AtkCastlePosY = config.Ia2I32a(dAtA.AtkCastlePosY)
	out.AtkGateName = dAtA.AtkGateName
	out.AtkGatePosX = config.Ia2I32a(dAtA.AtkGatePosX)
	out.AtkGatePosY = config.Ia2I32a(dAtA.AtkGatePosY)
	out.DefReliveName = dAtA.DefReliveName
	out.DefRelivePosX = int32(dAtA.DefRelivePosX)
	out.DefRelivePosY = int32(dAtA.DefRelivePosY)
	out.DefHomeName = dAtA.DefHomeName
	out.DefHomePosX = int32(dAtA.DefHomePosX)
	out.DefHomePosY = int32(dAtA.DefHomePosY)
	out.DefCastleName = dAtA.DefCastleName
	out.DefCastlePosX = config.Ia2I32a(dAtA.DefCastlePosX)
	out.DefCastlePosY = config.Ia2I32a(dAtA.DefCastlePosY)
	out.DefGateName = dAtA.DefGateName
	out.DefGatePosX = config.Ia2I32a(dAtA.DefGatePosX)
	out.DefGatePosY = config.Ia2I32a(dAtA.DefGatePosY)
	out.AtkTouShiName = dAtA.AtkTouShiName
	out.AtkTouShiPosX = config.Ia2I32a(dAtA.AtkTouShiPosX)
	out.AtkTouShiPosY = config.Ia2I32a(dAtA.AtkTouShiPosY)
	out.DefTouShiName = dAtA.DefTouShiName
	out.DefTouShiPosX = config.Ia2I32a(dAtA.DefTouShiPosX)
	out.DefTouShiPosY = config.Ia2I32a(dAtA.DefTouShiPosY)
	out.AtkFullProsperity = config.U64ToI32(dAtA.AtkFullProsperity)
	out.DefFullProsperity = config.U64ToI32(dAtA.DefFullProsperity)

	return out
}

func ArrayEncodeMingcWarSceneData(datas []*MingcWarSceneData) []*shared_proto.MingcWarSceneDataProto {

	out := make([]*shared_proto.MingcWarSceneDataProto, 0, len(datas))
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

func (dAtA *MingcWarSceneData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with MingcWarTouShiBuildingTargetData ----------------------------------

func LoadMingcWarTouShiBuildingTargetData(gos *config.GameObjects) (map[uint64]*MingcWarTouShiBuildingTargetData, map[*MingcWarTouShiBuildingTargetData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcWarTouShiBuildingTargetDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcWarTouShiBuildingTargetData, len(lIsT))
	pArSeRmAp := make(map[*MingcWarTouShiBuildingTargetData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcWarTouShiBuildingTargetData) {
			continue
		}

		dAtA, err := NewMingcWarTouShiBuildingTargetData(fIlEnAmE, pArSeR)
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

func SetRelatedMingcWarTouShiBuildingTargetData(dAtAmAp map[*MingcWarTouShiBuildingTargetData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcWarTouShiBuildingTargetDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcWarTouShiBuildingTargetDataKeyArray(datas []*MingcWarTouShiBuildingTargetData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMingcWarTouShiBuildingTargetData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcWarTouShiBuildingTargetData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcWarTouShiBuildingTargetData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcWarTouShiBuildingTargetData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Mingc
	dAtA.PosX = pArSeR.Int("pos_x")
	dAtA.PosY = pArSeR.Int("pos_y")
	dAtA.TargetX = pArSeR.IntArray("target_x", "", false)
	dAtA.TargetY = pArSeR.IntArray("target_y", "", false)

	return dAtA, nil
}

var vAlIdAtOrMingcWarTouShiBuildingTargetData = map[string]*config.Validator{

	"id":       config.ParseValidator("int>0", "", false, nil, nil),
	"mingc":    config.ParseValidator("string", "", false, nil, nil),
	"pos_x":    config.ParseValidator("int>0", "", false, nil, nil),
	"pos_y":    config.ParseValidator("int>0", "", false, nil, nil),
	"target_x": config.ParseValidator("int,duplicate", "", true, nil, nil),
	"target_y": config.ParseValidator("int,duplicate", "", true, nil, nil),
}

func (dAtA *MingcWarTouShiBuildingTargetData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MingcWarTouShiBuildingTargetData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MingcWarTouShiBuildingTargetData) Encode() *shared_proto.MingcWarTouShiBuildingTargetDataProto {
	out := &shared_proto.MingcWarTouShiBuildingTargetDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.Mingc != nil {
		out.Mingc = config.U64ToI32(dAtA.Mingc.Id)
	}
	out.PosX = int32(dAtA.PosX)
	out.PosY = int32(dAtA.PosY)
	out.TargetX = config.Ia2I32a(dAtA.TargetX)
	out.TargetY = config.Ia2I32a(dAtA.TargetY)

	return out
}

func ArrayEncodeMingcWarTouShiBuildingTargetData(datas []*MingcWarTouShiBuildingTargetData) []*shared_proto.MingcWarTouShiBuildingTargetDataProto {

	out := make([]*shared_proto.MingcWarTouShiBuildingTargetDataProto, 0, len(datas))
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

func (dAtA *MingcWarTouShiBuildingTargetData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Mingc = cOnFigS.GetMingcBaseData(pArSeR.Uint64("mingc"))
	if dAtA.Mingc == nil {
		return errors.Errorf("%s 配置的关联字段[mingc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mingc"), *pArSeR)
	}

	return nil
}

// start with MingcWarTroopLastBeatWhenFailData ----------------------------------

func LoadMingcWarTroopLastBeatWhenFailData(gos *config.GameObjects) (map[uint64]*MingcWarTroopLastBeatWhenFailData, map[*MingcWarTroopLastBeatWhenFailData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MingcWarTroopLastBeatWhenFailDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MingcWarTroopLastBeatWhenFailData, len(lIsT))
	pArSeRmAp := make(map[*MingcWarTroopLastBeatWhenFailData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMingcWarTroopLastBeatWhenFailData) {
			continue
		}

		dAtA, err := NewMingcWarTroopLastBeatWhenFailData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.BaiZhanLevel
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[BaiZhanLevel], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedMingcWarTroopLastBeatWhenFailData(dAtAmAp map[*MingcWarTroopLastBeatWhenFailData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MingcWarTroopLastBeatWhenFailDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMingcWarTroopLastBeatWhenFailDataKeyArray(datas []*MingcWarTroopLastBeatWhenFailData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.BaiZhanLevel)
		}
	}

	return out
}

func NewMingcWarTroopLastBeatWhenFailData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MingcWarTroopLastBeatWhenFailData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMingcWarTroopLastBeatWhenFailData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MingcWarTroopLastBeatWhenFailData{}

	dAtA.BaiZhanLevel = pArSeR.Uint64("bai_zhan_level")
	dAtA.SoliderAmount = pArSeR.Uint64("solider_amount")
	dAtA.HurtPercent, err = data.ParseAmount(pArSeR.String("hurt_percent"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[hurt_percent] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("hurt_percent"), dAtA)
	}

	dAtA.AtkBackHurtPercent, err = data.ParseAmount(pArSeR.String("atk_back_hurt_percent"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[atk_back_hurt_percent] 解析失败(data.ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("atk_back_hurt_percent"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrMingcWarTroopLastBeatWhenFailData = map[string]*config.Validator{

	"bai_zhan_level":        config.ParseValidator("int>0", "", false, nil, nil),
	"solider_amount":        config.ParseValidator("uint", "", false, nil, nil),
	"hurt_percent":          config.ParseValidator("string", "", false, nil, nil),
	"atk_back_hurt_percent": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *MingcWarTroopLastBeatWhenFailData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MingcWarTroopLastBeatWhenFailData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MingcWarTroopLastBeatWhenFailData) Encode() *shared_proto.MingcWarTroopLastBeatWhenFailDataProto {
	out := &shared_proto.MingcWarTroopLastBeatWhenFailDataProto{}
	out.BaiZhanLevel = config.U64ToI32(dAtA.BaiZhanLevel)
	out.SoliderAmount = config.U64ToI32(dAtA.SoliderAmount)
	if dAtA.HurtPercent != nil {
		out.HurtPercent = dAtA.HurtPercent.Encode()
	}
	if dAtA.AtkBackHurtPercent != nil {
		out.AtkBackHurtPercent = dAtA.AtkBackHurtPercent.Encode()
	}

	return out
}

func ArrayEncodeMingcWarTroopLastBeatWhenFailData(datas []*MingcWarTroopLastBeatWhenFailData) []*shared_proto.MingcWarTroopLastBeatWhenFailDataProto {

	out := make([]*shared_proto.MingcWarTroopLastBeatWhenFailDataProto, 0, len(datas))
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

func (dAtA *MingcWarTroopLastBeatWhenFailData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetCombatScene(string) *scene.CombatScene
	GetIcon(string) *icon.Icon
	GetMingcBaseData(uint64) *MingcBaseData
	GetMingcWarNpcGuildData(uint64) *MingcWarNpcGuildData
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
	GetPrize(int) *resdata.Prize
	GetSpriteStat(uint64) *data.SpriteStat
}
