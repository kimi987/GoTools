// AUTO_GEN, DONT MODIFY!!!
package location

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/country"
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

// start with LocationData ----------------------------------

func LoadLocationData(gos *config.GameObjects) (map[uint64]*LocationData, map[*LocationData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.LocationDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*LocationData, len(lIsT))
	pArSeRmAp := make(map[*LocationData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrLocationData) {
			continue
		}

		dAtA, err := NewLocationData(fIlEnAmE, pArSeR)
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

func SetRelatedLocationData(dAtAmAp map[*LocationData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.LocationDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetLocationDataKeyArray(datas []*LocationData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewLocationData(fIlEnAmE string, pArSeR *config.ObjectParser) (*LocationData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrLocationData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &LocationData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// releated field: RecommendCountry

	return dAtA, nil
}

var vAlIdAtOrLocationData = map[string]*config.Validator{

	"id":                config.ParseValidator("int>0", "", false, nil, nil),
	"name":              config.ParseValidator("string", "", false, nil, nil),
	"recommend_country": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *LocationData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *LocationData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *LocationData) Encode() *shared_proto.LocationDataProto {
	out := &shared_proto.LocationDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	if dAtA.RecommendCountry != nil {
		out.RecommendCountry = config.U64a2I32a(country.GetCountryDataKeyArray(dAtA.RecommendCountry))
	}

	return out
}

func ArrayEncodeLocationData(datas []*LocationData) []*shared_proto.LocationDataProto {

	out := make([]*shared_proto.LocationDataProto, 0, len(datas))
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

func (dAtA *LocationData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("recommend_country", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetCountryData(v)
		if obj != nil {
			dAtA.RecommendCountry = append(dAtA.RecommendCountry, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[recommend_country] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("recommend_country"), *pArSeR)
		}
	}

	return nil
}

type related_configs interface {
	GetCountryData(uint64) *country.CountryData
}
