// AUTO_GEN, DONT MODIFY!!!
package tag

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

// start with TagMiscData ----------------------------------

func LoadTagMiscData(gos *config.GameObjects) (*TagMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.TagMiscDataPath
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

	dAtA, err := NewTagMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedTagMiscData(gos *config.GameObjects, dAtA *TagMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TagMiscDataPath
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

func NewTagMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TagMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTagMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TagMiscData{}

	dAtA.MaxCount = 50
	if pArSeR.KeyExist("max_count") {
		dAtA.MaxCount = pArSeR.Uint64("max_count")
	}

	dAtA.MaxCharCount = 5
	if pArSeR.KeyExist("max_char_count") {
		dAtA.MaxCharCount = pArSeR.Uint64("max_char_count")
	}

	dAtA.MaxRecordCount = 50
	if pArSeR.KeyExist("max_record_count") {
		dAtA.MaxRecordCount = pArSeR.Uint64("max_record_count")
	}

	dAtA.MaxShowForViewCount = 5
	if pArSeR.KeyExist("max_show_for_view_count") {
		dAtA.MaxShowForViewCount = pArSeR.Uint64("max_show_for_view_count")
	}

	dAtA.MaxTagColorType = 4
	if pArSeR.KeyExist("max_tag_color_type") {
		dAtA.MaxTagColorType = pArSeR.Uint64("max_tag_color_type")
	}

	return dAtA, nil
}

var vAlIdAtOrTagMiscData = map[string]*config.Validator{

	"max_count":               config.ParseValidator("int>0", "", false, nil, []string{"50"}),
	"max_char_count":          config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"max_record_count":        config.ParseValidator("int>0", "", false, nil, []string{"50"}),
	"max_show_for_view_count": config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"max_tag_color_type":      config.ParseValidator("int>0", "", false, nil, []string{"4"}),
}

func (dAtA *TagMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TagMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TagMiscData) Encode() *shared_proto.TagMiscProto {
	out := &shared_proto.TagMiscProto{}
	out.MaxCount = config.U64ToI32(dAtA.MaxCount)
	out.MaxCharCount = config.U64ToI32(dAtA.MaxCharCount)
	out.MaxRecordCount = config.U64ToI32(dAtA.MaxRecordCount)
	out.MaxShowForViewCount = config.U64ToI32(dAtA.MaxShowForViewCount)

	return out
}

func ArrayEncodeTagMiscData(datas []*TagMiscData) []*shared_proto.TagMiscProto {

	out := make([]*shared_proto.TagMiscProto, 0, len(datas))
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

func (dAtA *TagMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
