// AUTO_GEN, DONT MODIFY!!!
package buffer

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
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

// start with BufferData ----------------------------------

func LoadBufferData(gos *config.GameObjects) (map[uint64]*BufferData, map[*BufferData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BufferDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BufferData, len(lIsT))
	pArSeRmAp := make(map[*BufferData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBufferData) {
			continue
		}

		dAtA, err := NewBufferData(fIlEnAmE, pArSeR)
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

func SetRelatedBufferData(dAtAmAp map[*BufferData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BufferDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBufferDataKeyArray(datas []*BufferData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBufferData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BufferData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBufferData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BufferData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.NameDesc = pArSeR.String("name_desc")
	// releated field: Icon
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Type = pArSeR.Uint64("type")
	// skip field: TypeData
	dAtA.ShowKeepDuration, err = config.ParseDuration(pArSeR.String("show_keep_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[show_keep_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("show_keep_duration"), dAtA)
	}

	dAtA.ShowLevel = pArSeR.Uint64("show_level")
	dAtA.BuffGoodsId = pArSeR.Uint64("buff_goods_id")
	// skip field: BuffGoodsData

	return dAtA, nil
}

var vAlIdAtOrBufferData = map[string]*config.Validator{

	"id":                 config.ParseValidator("int>0", "", false, nil, nil),
	"name":               config.ParseValidator("string", "", false, nil, nil),
	"name_desc":          config.ParseValidator("string", "", false, nil, nil),
	"icon":               config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"desc":               config.ParseValidator("string", "", false, nil, nil),
	"type":               config.ParseValidator("int>0", "", false, nil, nil),
	"show_keep_duration": config.ParseValidator("string", "", false, nil, nil),
	"show_level":         config.ParseValidator("int>0", "", false, nil, nil),
	"buff_goods_id":      config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *BufferData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BufferData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BufferData) Encode() *shared_proto.BufferDataProto {
	out := &shared_proto.BufferDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.NameDesc = dAtA.NameDesc
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.Desc = dAtA.Desc
	out.Type = config.U64ToI32(dAtA.Type)
	out.ShowKeepDuration = config.Duration2I32Seconds(dAtA.ShowKeepDuration)
	out.ShowLevel = config.U64ToI32(dAtA.ShowLevel)
	out.BuffGoodsId = config.U64ToI32(dAtA.BuffGoodsId)

	return out
}

func ArrayEncodeBufferData(datas []*BufferData) []*shared_proto.BufferDataProto {

	out := make([]*shared_proto.BufferDataProto, 0, len(datas))
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

func (dAtA *BufferData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	return nil
}

// start with BufferTypeData ----------------------------------

func LoadBufferTypeData(gos *config.GameObjects) (map[uint64]*BufferTypeData, map[*BufferTypeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BufferTypeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BufferTypeData, len(lIsT))
	pArSeRmAp := make(map[*BufferTypeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBufferTypeData) {
			continue
		}

		dAtA, err := NewBufferTypeData(fIlEnAmE, pArSeR)
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

func SetRelatedBufferTypeData(dAtAmAp map[*BufferTypeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BufferTypeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBufferTypeDataKeyArray(datas []*BufferTypeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBufferTypeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BufferTypeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBufferTypeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BufferTypeData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.TypeName = pArSeR.String("type_name")
	dAtA.TypeDesc = pArSeR.String("type_desc")
	dAtA.IsMian = pArSeR.Bool("is_mian")
	dAtA.BuffGroup = pArSeR.Uint64("buff_group")
	// releated field: Icon
	dAtA.Sort = 0
	if pArSeR.KeyExist("sort") {
		dAtA.Sort = pArSeR.Uint64("sort")
	}

	return dAtA, nil
}

var vAlIdAtOrBufferTypeData = map[string]*config.Validator{

	"id":         config.ParseValidator("int>0", "", false, nil, nil),
	"type_name":  config.ParseValidator("string", "", false, nil, nil),
	"type_desc":  config.ParseValidator("string", "", false, nil, nil),
	"is_mian":    config.ParseValidator("bool", "", false, nil, nil),
	"buff_group": config.ParseValidator("uint", "", false, nil, nil),
	"icon":       config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"sort":       config.ParseValidator("int>0", "", false, nil, []string{"0"}),
}

func (dAtA *BufferTypeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BufferTypeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BufferTypeData) Encode() *shared_proto.BufferTypeDataProto {
	out := &shared_proto.BufferTypeDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.TypeName = dAtA.TypeName
	out.TypeDesc = dAtA.TypeDesc
	out.IsMian = dAtA.IsMian
	out.BuffGroup = config.U64ToI32(dAtA.BuffGroup)
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.Sort = config.U64ToI32(dAtA.Sort)

	return out
}

func ArrayEncodeBufferTypeData(datas []*BufferTypeData) []*shared_proto.BufferTypeDataProto {

	out := make([]*shared_proto.BufferTypeDataProto, 0, len(datas))
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

func (dAtA *BufferTypeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	return nil
}

type related_configs interface {
	GetIcon(string) *icon.Icon
}
