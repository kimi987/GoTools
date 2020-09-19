// AUTO_GEN, DONT MODIFY!!!
package gardendata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/i18n"
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

// start with GardenConfig ----------------------------------

func LoadGardenConfig(gos *config.GameObjects) (*GardenConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.GardenConfigPath
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

	dAtA, err := NewGardenConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedGardenConfig(gos *config.GameObjects, dAtA *GardenConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GardenConfigPath
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

func NewGardenConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*GardenConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGardenConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GardenConfig{}

	dAtA.TreasuryTreeFullTimes = pArSeR.Uint64("treasury_tree_full_times")
	if pArSeR.KeyExist("treasury_tree_collect_duration") {
		dAtA.TreasuryTreeCollectDuration, err = config.ParseDuration(pArSeR.String("treasury_tree_collect_duration"))
	} else {
		dAtA.TreasuryTreeCollectDuration, err = config.ParseDuration("24h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[treasury_tree_collect_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("treasury_tree_collect_duration"), dAtA)
	}

	dAtA.TreasuryTreeHelpMeLogCount = 5
	if pArSeR.KeyExist("treasury_tree_help_me_log_count") {
		dAtA.TreasuryTreeHelpMeLogCount = pArSeR.Uint64("treasury_tree_help_me_log_count")
	}

	return dAtA, nil
}

var vAlIdAtOrGardenConfig = map[string]*config.Validator{

	"treasury_tree_full_times":        config.ParseValidator("int>0", "", false, nil, nil),
	"treasury_tree_collect_duration":  config.ParseValidator("string", "", false, nil, []string{"24h"}),
	"treasury_tree_help_me_log_count": config.ParseValidator("int>0", "", false, nil, []string{"5"}),
}

func (dAtA *GardenConfig) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GardenConfig) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GardenConfig) Encode() *shared_proto.GradonConfigProto {
	out := &shared_proto.GradonConfigProto{}
	out.TreasuryTreeFullTimes = config.U64ToI32(dAtA.TreasuryTreeFullTimes)
	out.TreasuryTreeHelpMeLogCount = config.U64ToI32(dAtA.TreasuryTreeHelpMeLogCount)

	return out
}

func ArrayEncodeGardenConfig(datas []*GardenConfig) []*shared_proto.GradonConfigProto {

	out := make([]*shared_proto.GradonConfigProto, 0, len(datas))
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

func (dAtA *GardenConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with TreasuryTreeData ----------------------------------

func LoadTreasuryTreeData(gos *config.GameObjects) (map[uint64]*TreasuryTreeData, map[*TreasuryTreeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TreasuryTreeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TreasuryTreeData, len(lIsT))
	pArSeRmAp := make(map[*TreasuryTreeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTreasuryTreeData) {
			continue
		}

		dAtA, err := NewTreasuryTreeData(fIlEnAmE, pArSeR)
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

func SetRelatedTreasuryTreeData(dAtAmAp map[*TreasuryTreeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TreasuryTreeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTreasuryTreeDataKeyArray(datas []*TreasuryTreeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTreasuryTreeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TreasuryTreeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTreasuryTreeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TreasuryTreeData{}

	dAtA.Season = shared_proto.Season(shared_proto.Season_value[strings.ToUpper(pArSeR.String("season"))])
	if i, err := strconv.ParseInt(pArSeR.String("season"), 10, 32); err == nil {
		dAtA.Season = shared_proto.Season(i)
	}

	// releated field: Prize

	// calculate fields
	dAtA.Id = uint64(dAtA.Season)

	// i18n fields
	dAtA.Desc = i18n.NewI18nRef(fIlEnAmE, "desc", dAtA.Id, pArSeR.String("desc"))

	return dAtA, nil
}

var vAlIdAtOrTreasuryTreeData = map[string]*config.Validator{

	"season": config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Season_value, 0), nil),
	"prize":  config.ParseValidator("string", "", false, nil, nil),
	"desc":   config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *TreasuryTreeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TreasuryTreeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TreasuryTreeData) Encode() *shared_proto.TreasuryTreeDataProto {
	out := &shared_proto.TreasuryTreeDataProto{}
	out.Season = dAtA.Season
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	if dAtA.Desc != nil {
		out.Desc = dAtA.Desc.Encode()
	}

	return out
}

func ArrayEncodeTreasuryTreeData(datas []*TreasuryTreeData) []*shared_proto.TreasuryTreeDataProto {

	out := make([]*shared_proto.TreasuryTreeDataProto, 0, len(datas))
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

func (dAtA *TreasuryTreeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

type related_configs interface {
	GetPrize(int) *resdata.Prize
}
