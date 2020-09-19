// AUTO_GEN, DONT MODIFY!!!
package icon

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

// start with Icon ----------------------------------

func LoadIcon(gos *config.GameObjects) (map[string]*Icon, map[*Icon]*config.ObjectParser, error) {
	fIlEnAmE := confpath.IconPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*Icon, len(lIsT))
	pArSeRmAp := make(map[*Icon]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrIcon) {
			continue
		}

		dAtA, err := NewIcon(fIlEnAmE, pArSeR)
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

func SetRelatedIcon(dAtAmAp map[*Icon]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.IconPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetIconKeyArray(datas []*Icon) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewIcon(fIlEnAmE string, pArSeR *config.ObjectParser) (*Icon, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrIcon)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &Icon{}

	dAtA.Id = pArSeR.String("id")
	dAtA.DefaultIcon = pArSeR.String("icon")
	dAtA.MiddleIcon = ""
	if pArSeR.KeyExist("middle_icon") {
		dAtA.MiddleIcon = pArSeR.String("middle_icon")
	}

	dAtA.BigIcon = ""
	if pArSeR.KeyExist("big_icon") {
		dAtA.BigIcon = pArSeR.String("big_icon")
	}

	dAtA.HeadIcon = ""
	if pArSeR.KeyExist("head_icon") {
		dAtA.HeadIcon = pArSeR.String("head_icon")
	}

	dAtA.TailIcon = ""
	if pArSeR.KeyExist("tail_icon") {
		dAtA.TailIcon = pArSeR.String("tail_icon")
	}

	dAtA.SuperBigIcon = ""
	if pArSeR.KeyExist("super_big_icon") {
		dAtA.SuperBigIcon = pArSeR.String("super_big_icon")
	}

	dAtA.CaptainHead = false
	if pArSeR.KeyExist("captain_head") {
		dAtA.CaptainHead = pArSeR.Bool("captain_head")
	}

	dAtA.Tab = 0
	if pArSeR.KeyExist("tab") {
		dAtA.Tab = pArSeR.Uint64("tab")
	}

	dAtA.Text = ""
	if pArSeR.KeyExist("text") {
		dAtA.Text = pArSeR.String("text")
	}

	return dAtA, nil
}

var vAlIdAtOrIcon = map[string]*config.Validator{

	"id":             config.ParseValidator("string", "", false, nil, nil),
	"icon":           config.ParseValidator("string>0", "", false, nil, nil),
	"middle_icon":    config.ParseValidator("string", "", false, nil, []string{""}),
	"big_icon":       config.ParseValidator("string", "", false, nil, []string{""}),
	"head_icon":      config.ParseValidator("string", "", false, nil, []string{""}),
	"tail_icon":      config.ParseValidator("string", "", false, nil, []string{""}),
	"super_big_icon": config.ParseValidator("string", "", false, nil, []string{""}),
	"captain_head":   config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"tab":            config.ParseValidator("int", "", false, nil, []string{"0"}),
	"text":           config.ParseValidator("string", "", false, nil, []string{""}),
}

func (dAtA *Icon) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *Icon) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *Icon) Encode() *shared_proto.IconProto {
	out := &shared_proto.IconProto{}
	out.Id = dAtA.Id
	out.Icon = dAtA.DefaultIcon
	out.MiddleIcon = dAtA.MiddleIcon
	out.BigIcon = dAtA.BigIcon
	out.HeadIcon = dAtA.HeadIcon
	out.TailIcon = dAtA.TailIcon
	out.SuperBigIcon = dAtA.SuperBigIcon
	out.Tab = config.U64ToI32(dAtA.Tab)
	out.Text = dAtA.Text

	return out
}

func ArrayEncodeIcon(datas []*Icon) []*shared_proto.IconProto {

	out := make([]*shared_proto.IconProto, 0, len(datas))
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

func (dAtA *Icon) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
