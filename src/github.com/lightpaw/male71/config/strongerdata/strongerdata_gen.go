// AUTO_GEN, DONT MODIFY!!!
package strongerdata

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

// start with StrongerData ----------------------------------

func LoadStrongerData(gos *config.GameObjects) (map[uint64]*StrongerData, map[*StrongerData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.StrongerDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*StrongerData, len(lIsT))
	pArSeRmAp := make(map[*StrongerData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrStrongerData) {
			continue
		}

		dAtA, err := NewStrongerData(fIlEnAmE, pArSeR)
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

func SetRelatedStrongerData(dAtAmAp map[*StrongerData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.StrongerDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetStrongerDataKeyArray(datas []*StrongerData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewStrongerData(fIlEnAmE string, pArSeR *config.ObjectParser) (*StrongerData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrStrongerData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &StrongerData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Type = pArSeR.Uint64("type")
	dAtA.X = pArSeR.Uint64("x")
	dAtA.Y = pArSeR.Uint64("y")
	dAtA.Z = pArSeR.Uint64("z")

	// calculate fields
	dAtA.Id = GetStrongerDataId(dAtA.Level, dAtA.Type)

	return dAtA, nil
}

var vAlIdAtOrStrongerData = map[string]*config.Validator{

	"level": config.ParseValidator("int>0", "", false, nil, nil),
	"type":  config.ParseValidator("int>0", "", false, nil, nil),
	"x":     config.ParseValidator("int>0", "", false, nil, nil),
	"y":     config.ParseValidator("uint", "", false, nil, nil),
	"z":     config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *StrongerData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *StrongerData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *StrongerData) Encode() *shared_proto.StrongerDataProto {
	out := &shared_proto.StrongerDataProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Type = config.U64ToI32(dAtA.Type)
	out.X = config.U64ToI32(dAtA.X)
	out.Y = config.U64ToI32(dAtA.Y)
	out.Z = config.U64ToI32(dAtA.Z)

	return out
}

func ArrayEncodeStrongerData(datas []*StrongerData) []*shared_proto.StrongerDataProto {

	out := make([]*shared_proto.StrongerDataProto, 0, len(datas))
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

func (dAtA *StrongerData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
