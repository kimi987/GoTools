// AUTO_GEN, DONT MODIFY!!!
package scene

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

// start with CombatScene ----------------------------------

func LoadCombatScene(gos *config.GameObjects) (map[string]*CombatScene, map[*CombatScene]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CombatScenePath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*CombatScene, len(lIsT))
	pArSeRmAp := make(map[*CombatScene]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCombatScene) {
			continue
		}

		dAtA, err := NewCombatScene(fIlEnAmE, pArSeR)
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

func SetRelatedCombatScene(dAtAmAp map[*CombatScene]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CombatScenePath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCombatSceneKeyArray(datas []*CombatScene) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCombatScene(fIlEnAmE string, pArSeR *config.ObjectParser) (*CombatScene, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCombatScene)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CombatScene{}

	dAtA.Id = pArSeR.String("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.MapRes = pArSeR.String("map_res")
	dAtA.WallMapRes = pArSeR.String("wall_map_res")

	return dAtA, nil
}

var vAlIdAtOrCombatScene = map[string]*config.Validator{

	"id":           config.ParseValidator("string>0", "", false, nil, nil),
	"name":         config.ParseValidator("string>0", "", false, nil, nil),
	"map_res":      config.ParseValidator("string>0", "", false, nil, nil),
	"wall_map_res": config.ParseValidator("string>0", "", false, nil, nil),
}

func (dAtA *CombatScene) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
