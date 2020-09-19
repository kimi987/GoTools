// AUTO_GEN, DONT MODIFY!!!
package pushdata

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

// start with PushData ----------------------------------

func LoadPushData(gos *config.GameObjects) (map[uint64]*PushData, map[*PushData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.PushDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*PushData, len(lIsT))
	pArSeRmAp := make(map[*PushData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrPushData) {
			continue
		}

		dAtA, err := NewPushData(fIlEnAmE, pArSeR)
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

func SetRelatedPushData(dAtAmAp map[*PushData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.PushDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetPushDataKeyArray(datas []*PushData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewPushData(fIlEnAmE string, pArSeR *config.ObjectParser) (*PushData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrPushData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &PushData{}

	dAtA.Type = shared_proto.SettingType(shared_proto.SettingType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.SettingType(i)
	}

	dAtA.Title = pArSeR.String("title")
	dAtA.Content = pArSeR.String("content")
	dAtA.TickTime = pArSeR.String("tick_time")

	// calculate fields
	dAtA.Id = uint64(dAtA.Type)

	return dAtA, nil
}

var vAlIdAtOrPushData = map[string]*config.Validator{

	"type":      config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.SettingType_value, 0), nil),
	"title":     config.ParseValidator("string", "", false, nil, nil),
	"content":   config.ParseValidator("string", "", false, nil, nil),
	"tick_time": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *PushData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *PushData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *PushData) Encode() *shared_proto.PushDataProto {
	out := &shared_proto.PushDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Type = dAtA.Type
	out.Title = dAtA.Title
	out.Content = dAtA.Content
	out.TickTime = dAtA.TickTime

	return out
}

func ArrayEncodePushData(datas []*PushData) []*shared_proto.PushDataProto {

	out := make([]*shared_proto.PushDataProto, 0, len(datas))
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

func (dAtA *PushData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
