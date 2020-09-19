// AUTO_GEN, DONT MODIFY!!!
package settings

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
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

// start with PrivacySettingData ----------------------------------

func LoadPrivacySettingData(gos *config.GameObjects) (map[uint64]*PrivacySettingData, map[*PrivacySettingData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.PrivacySettingDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*PrivacySettingData, len(lIsT))
	pArSeRmAp := make(map[*PrivacySettingData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrPrivacySettingData) {
			continue
		}

		dAtA, err := NewPrivacySettingData(fIlEnAmE, pArSeR)
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

func SetRelatedPrivacySettingData(dAtAmAp map[*PrivacySettingData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.PrivacySettingDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetPrivacySettingDataKeyArray(datas []*PrivacySettingData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewPrivacySettingData(fIlEnAmE string, pArSeR *config.ObjectParser) (*PrivacySettingData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrPrivacySettingData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &PrivacySettingData{}

	dAtA.Id = pArSeR.Uint64("id")
	// skip field: SettingType
	dAtA.Name = pArSeR.String("name")
	dAtA.NameType = pArSeR.String("name_type")
	dAtA.DefaultOpen = false
	if pArSeR.KeyExist("default_open") {
		dAtA.DefaultOpen = pArSeR.Bool("default_open")
	}

	dAtA.RuleTitle = pArSeR.String("rule_title")
	dAtA.RuleDesc = pArSeR.String("rule_desc")

	return dAtA, nil
}

var vAlIdAtOrPrivacySettingData = map[string]*config.Validator{

	"id":           config.ParseValidator("int>0", "", false, nil, nil),
	"name":         config.ParseValidator("string", "", false, nil, nil),
	"name_type":    config.ParseValidator("string", "", false, nil, nil),
	"default_open": config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"rule_title":   config.ParseValidator("string", "", false, nil, nil),
	"rule_desc":    config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *PrivacySettingData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *PrivacySettingData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *PrivacySettingData) Encode() *shared_proto.PrivacySettingDataProto {
	out := &shared_proto.PrivacySettingDataProto{}
	out.SettingType = dAtA.SettingType
	out.Name = dAtA.Name
	out.NameType = dAtA.NameType
	out.RuleTitle = dAtA.RuleTitle
	out.RuleDesc = dAtA.RuleDesc

	return out
}

func ArrayEncodePrivacySettingData(datas []*PrivacySettingData) []*shared_proto.PrivacySettingDataProto {

	out := make([]*shared_proto.PrivacySettingDataProto, 0, len(datas))
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

func (dAtA *PrivacySettingData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with SettingMiscData ----------------------------------

func LoadSettingMiscData(gos *config.GameObjects) (*SettingMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.SettingMiscDataPath
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

	dAtA, err := NewSettingMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedSettingMiscData(gos *config.GameObjects, dAtA *SettingMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SettingMiscDataPath
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

func NewSettingMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SettingMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSettingMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SettingMiscData{}

	for _, v := range pArSeR.StringArray("default_settings", "", false) {
		x := shared_proto.SettingType(shared_proto.SettingType_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.SettingType(i)
		}
		dAtA.DefaultSettings = append(dAtA.DefaultSettings, x)
	}

	return dAtA, nil
}

var vAlIdAtOrSettingMiscData = map[string]*config.Validator{

	"default_settings": config.ParseValidator("string,notAllNil", "", true, config.EnumMapKeys(shared_proto.SettingType_value, 0), nil),
}

func (dAtA *SettingMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
}
