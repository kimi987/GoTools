// AUTO_GEN, DONT MODIFY!!!
package blockdata

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

// start with BlockData ----------------------------------

func LoadBlockData(gos *config.GameObjects) (map[uint64]*BlockData, map[*BlockData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BlockDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BlockData, len(lIsT))
	pArSeRmAp := make(map[*BlockData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBlockData) {
			continue
		}

		dAtA, err := NewBlockData(fIlEnAmE, pArSeR)
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

func SetRelatedBlockData(dAtAmAp map[*BlockData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BlockDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBlockDataKeyArray(datas []*BlockData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBlockData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BlockData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBlockData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BlockData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.XLen = pArSeR.Uint64("xlen")
	dAtA.YLen = pArSeR.Uint64("ylen")
	dAtA.AutoExpandBaseCount = pArSeR.Uint64("auto_expand_base_count")
	dAtA.NewHeroCrowdedCapcity = pArSeR.Uint64("new_hero_crowded_capcity")
	dAtA.BaseCountLimit = pArSeR.Uint64("base_count_limit")
	dAtA.EdgeNotHomeLen = 0
	if pArSeR.KeyExist("edge_not_home_len") {
		dAtA.EdgeNotHomeLen = pArSeR.Uint64("edge_not_home_len")
	}

	dAtA.ProtoFileName = pArSeR.String("proto_file_name")
	dAtA.Radius = []uint64{10, 30, 50}
	if pArSeR.KeyExist("radius") {
		dAtA.Radius = pArSeR.Uint64Array("radius", "", false)
	}

	// calculate fields
	dAtA.CenterX = dAtA.XLen / 2
	// calculate fields
	dAtA.CenterY = dAtA.YLen / 2

	return dAtA, nil
}

var vAlIdAtOrBlockData = map[string]*config.Validator{

	"id":   config.ParseValidator("int>0", "", false, nil, nil),
	"name": config.ParseValidator("string", "", false, nil, nil),
	"xlen": config.ParseValidator("int>0", "", false, nil, nil),
	"ylen": config.ParseValidator("int>0", "", false, nil, nil),
	"auto_expand_base_count":   config.ParseValidator("int>0", "", false, nil, nil),
	"new_hero_crowded_capcity": config.ParseValidator("int>0", "", false, nil, nil),
	"base_count_limit":         config.ParseValidator("int>0", "", false, nil, nil),
	"edge_not_home_len":        config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"proto_file_name":          config.ParseValidator("string", "", false, nil, nil),
	"radius":                   config.ParseValidator("uint", "", true, nil, []string{"10", "30", "50"}),
}

func (dAtA *BlockData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
