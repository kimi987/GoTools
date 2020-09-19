// AUTO_GEN, DONT MODIFY!!!
package pvetroop

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

// start with PveTroopData ----------------------------------

func LoadPveTroopData(gos *config.GameObjects) (map[uint64]*PveTroopData, map[*PveTroopData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.PveTroopDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*PveTroopData, len(lIsT))
	pArSeRmAp := make(map[*PveTroopData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrPveTroopData) {
			continue
		}

		dAtA, err := NewPveTroopData(fIlEnAmE, pArSeR)
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

func SetRelatedPveTroopData(dAtAmAp map[*PveTroopData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.PveTroopDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetPveTroopDataKeyArray(datas []*PveTroopData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewPveTroopData(fIlEnAmE string, pArSeR *config.ObjectParser) (*PveTroopData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrPveTroopData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &PveTroopData{}

	dAtA.PveTroopType = shared_proto.PveTroopType(shared_proto.PveTroopType_value[strings.ToUpper(pArSeR.String("pve_troop_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("pve_troop_type"), 10, 32); err == nil {
		dAtA.PveTroopType = shared_proto.PveTroopType(i)
	}

	dAtA.Capacity = 5
	if pArSeR.KeyExist("capacity") {
		dAtA.Capacity = pArSeR.Uint64("capacity")
	}

	dAtA.MinCaptainCount = 5
	if pArSeR.KeyExist("min_captain_count") {
		dAtA.MinCaptainCount = pArSeR.Uint64("min_captain_count")
	}

	// calculate fields
	dAtA.Id = uint64(dAtA.PveTroopType)

	return dAtA, nil
}

var vAlIdAtOrPveTroopData = map[string]*config.Validator{

	"pve_troop_type":    config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.PveTroopType_value, 0), nil),
	"capacity":          config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"min_captain_count": config.ParseValidator("int>0", "", false, nil, []string{"5"}),
}

func (dAtA *PveTroopData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *PveTroopData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *PveTroopData) Encode() *shared_proto.PveTroopDataProto {
	out := &shared_proto.PveTroopDataProto{}
	out.PveTroopType = dAtA.PveTroopType
	out.Capacity = config.U64ToI32(dAtA.Capacity)
	out.MinCaptainCount = config.U64ToI32(dAtA.MinCaptainCount)

	return out
}

func ArrayEncodePveTroopData(datas []*PveTroopData) []*shared_proto.PveTroopDataProto {

	out := make([]*shared_proto.PveTroopDataProto, 0, len(datas))
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

func (dAtA *PveTroopData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
