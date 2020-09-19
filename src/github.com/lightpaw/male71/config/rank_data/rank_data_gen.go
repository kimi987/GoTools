// AUTO_GEN, DONT MODIFY!!!
package rank_data

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

// start with RankMiscData ----------------------------------

func LoadRankMiscData(gos *config.GameObjects) (*RankMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.RankMiscDataPath
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

	dAtA, err := NewRankMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedRankMiscData(gos *config.GameObjects, dAtA *RankMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RankMiscDataPath
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

func NewRankMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*RankMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRankMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RankMiscData{}

	dAtA.RankCountPerPage = 5
	if pArSeR.KeyExist("rank_count_per_page") {
		dAtA.RankCountPerPage = pArSeR.Uint64("rank_count_per_page")
	}

	dAtA.MaxRankCount = 10000
	if pArSeR.KeyExist("max_rank_count") {
		dAtA.MaxRankCount = pArSeR.Uint64("max_rank_count")
	}

	return dAtA, nil
}

var vAlIdAtOrRankMiscData = map[string]*config.Validator{

	"rank_count_per_page": config.ParseValidator("int>0", "", false, nil, []string{"5"}),
	"max_rank_count":      config.ParseValidator("int>0", "", false, nil, []string{"10000"}),
}

func (dAtA *RankMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *RankMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *RankMiscData) Encode() *shared_proto.RankMiscProto {
	out := &shared_proto.RankMiscProto{}
	out.RankCountPerPage = config.U64ToI32(dAtA.RankCountPerPage)

	return out
}

func ArrayEncodeRankMiscData(datas []*RankMiscData) []*shared_proto.RankMiscProto {

	out := make([]*shared_proto.RankMiscProto, 0, len(datas))
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

func (dAtA *RankMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
