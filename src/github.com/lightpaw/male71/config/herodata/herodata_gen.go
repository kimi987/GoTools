// AUTO_GEN, DONT MODIFY!!!
package herodata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/military_data"
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

// start with HeroLevelData ----------------------------------

func LoadHeroLevelData(gos *config.GameObjects) (map[uint64]*HeroLevelData, map[*HeroLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.HeroLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*HeroLevelData, len(lIsT))
	pArSeRmAp := make(map[*HeroLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrHeroLevelData) {
			continue
		}

		dAtA, err := NewHeroLevelData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Level
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Level], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedHeroLevelData(dAtAmAp map[*HeroLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.HeroLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetHeroLevelDataKeyArray(datas []*HeroLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewHeroLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*HeroLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrHeroLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &HeroLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	// skip field: Sub
	// releated field: CaptainTrainingLevel

	return dAtA, nil
}

var vAlIdAtOrHeroLevelData = map[string]*config.Validator{

	"level":                  config.ParseValidator("int>0", "", false, nil, nil),
	"captain_training_level": config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
}

func (dAtA *HeroLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("captain_training_level", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetTrainingLevelData(v)
		if obj != nil {
			dAtA.CaptainTrainingLevel = append(dAtA.CaptainTrainingLevel, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[captain_training_level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain_training_level"), *pArSeR)
		}
	}

	return nil
}

type related_configs interface {
	GetTrainingLevelData(uint64) *military_data.TrainingLevelData
}
