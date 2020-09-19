// AUTO_GEN, DONT MODIFY!!!
package i18n

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

// start with I18nData ----------------------------------

func LoadI18nData(gos *config.GameObjects) (map[string]*I18nData, map[*I18nData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.I18nDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*I18nData, len(lIsT))
	pArSeRmAp := make(map[*I18nData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrI18nData) {
			continue
		}

		dAtA, err := NewI18nData(fIlEnAmE, pArSeR)
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

func SetRelatedI18nData(dAtAmAp map[*I18nData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.I18nDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetI18nDataKeyArray(datas []*I18nData) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewI18nData(fIlEnAmE string, pArSeR *config.ObjectParser) (*I18nData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrI18nData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &I18nData{}

	dAtA.Id = pArSeR.String("id")
	dAtA.Language = pArSeR.String("language")
	dAtA.Display = pArSeR.Bool("display")
	// skip field: Pair

	return dAtA, nil
}

var vAlIdAtOrI18nData = map[string]*config.Validator{

	"id":       config.ParseValidator("string", "", false, nil, nil),
	"language": config.ParseValidator("string", "", false, nil, nil),
	"display":  config.ParseValidator("bool", "", false, nil, nil),
}

func (dAtA *I18nData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *I18nData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *I18nData) Encode() *shared_proto.I18NDataProto {
	out := &shared_proto.I18NDataProto{}
	out.Id = dAtA.Id
	out.Language = dAtA.Language
	out.Display = dAtA.Display
	if dAtA.Pair != nil {
		out.Pair = ArrayEncodeI18nPair(dAtA.Pair)
	}

	return out
}

func ArrayEncodeI18nData(datas []*I18nData) []*shared_proto.I18NDataProto {

	out := make([]*shared_proto.I18NDataProto, 0, len(datas))
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

func (dAtA *I18nData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with I18nPair ----------------------------------

func NewI18nPair(fIlEnAmE string, pArSeR *config.ObjectParser) (*I18nPair, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrI18nPair)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &I18nPair{}

	dAtA.Key = pArSeR.StringArray("key", "", false)
	dAtA.Value = pArSeR.String("value")

	return dAtA, nil
}

var vAlIdAtOrI18nPair = map[string]*config.Validator{

	"key":   config.ParseValidator("string", "", true, nil, nil),
	"value": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *I18nPair) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *I18nPair) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *I18nPair) Encode() *shared_proto.I18NPairProto {
	out := &shared_proto.I18NPairProto{}
	out.Key = dAtA.Key
	out.Value = dAtA.Value

	return out
}

func ArrayEncodeI18nPair(datas []*I18nPair) []*shared_proto.I18NPairProto {

	out := make([]*shared_proto.I18NPairProto, 0, len(datas))
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

func (dAtA *I18nPair) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
