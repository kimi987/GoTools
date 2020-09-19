// AUTO_GEN, DONT MODIFY!!!
package function

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/dungeon"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/male7/config/towerdata"
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

// start with FunctionOpenData ----------------------------------

func LoadFunctionOpenData(gos *config.GameObjects) (map[uint64]*FunctionOpenData, map[*FunctionOpenData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FunctionOpenDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*FunctionOpenData, len(lIsT))
	pArSeRmAp := make(map[*FunctionOpenData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFunctionOpenData) {
			continue
		}

		dAtA, err := NewFunctionOpenData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.FunctionType
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[FunctionType], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedFunctionOpenData(dAtAmAp map[*FunctionOpenData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FunctionOpenDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFunctionOpenDataKeyArray(datas []*FunctionOpenData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.FunctionType)
		}
	}

	return out
}

func NewFunctionOpenData(fIlEnAmE string, pArSeR *config.ObjectParser) (*FunctionOpenData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFunctionOpenData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FunctionOpenData{}

	dAtA.FunctionType = pArSeR.Uint64("function_type")
	dAtA.Desc = ""
	if pArSeR.KeyExist("desc") {
		dAtA.Desc = pArSeR.String("desc")
	}

	dAtA.Icon = ""
	if pArSeR.KeyExist("icon") {
		dAtA.Icon = pArSeR.String("icon")
	}

	dAtA.NotifyOrder = 0
	if pArSeR.KeyExist("notify_order") {
		dAtA.NotifyOrder = pArSeR.Uint64("notify_order")
	}

	// releated field: GuanFuLevel
	// releated field: Building
	// releated field: HeroLevel
	// releated field: MainTask
	// releated field: BaYeStage
	// releated field: TowerFloor
	// releated field: Dungeon

	return dAtA, nil
}

var vAlIdAtOrFunctionOpenData = map[string]*config.Validator{

	"function_type": config.ParseValidator("int>0", "", false, nil, nil),
	"desc":          config.ParseValidator("string", "", false, nil, []string{""}),
	"icon":          config.ParseValidator("string", "", false, nil, []string{""}),
	"notify_order":  config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"guan_fu_level": config.ParseValidator("string", "", false, nil, nil),
	"building":      config.ParseValidator("string", "", false, nil, nil),
	"hero_level":    config.ParseValidator("string", "", false, nil, nil),
	"main_task":     config.ParseValidator("string", "", false, nil, nil),
	"ba_ye_stage":   config.ParseValidator("string", "", false, nil, nil),
	"tower_floor":   config.ParseValidator("string", "", false, nil, nil),
	"dungeon":       config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *FunctionOpenData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *FunctionOpenData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *FunctionOpenData) Encode() *shared_proto.FunctionOpenDataProto {
	out := &shared_proto.FunctionOpenDataProto{}
	out.FunctionType = config.U64ToI32(dAtA.FunctionType)
	out.Desc = dAtA.Desc
	out.Icon = dAtA.Icon
	out.NotifyOrder = config.U64ToI32(dAtA.NotifyOrder)
	if dAtA.GuanFuLevel != nil {
		out.GuanFuLevel = config.U64ToI32(dAtA.GuanFuLevel.Level)
	}
	if dAtA.Building != nil {
		out.BuildingId = config.U64ToI32(dAtA.Building.Level)
	}
	if dAtA.HeroLevel != nil {
		out.HeroLevel = config.U64ToI32(dAtA.HeroLevel.Level)
	}
	if dAtA.MainTask != nil {
		out.MainTask = config.U64ToI32(dAtA.MainTask.Sequence)
	}
	if dAtA.BaYeStage != nil {
		out.BaYeStage = config.U64ToI32(dAtA.BaYeStage.Stage)
	}
	if dAtA.TowerFloor != nil {
		out.TowerFloor = config.U64ToI32(dAtA.TowerFloor.Floor)
	}
	if dAtA.Dungeon != nil {
		out.Dungeon = config.U64ToI32(dAtA.Dungeon.Id)
	}

	return out
}

func ArrayEncodeFunctionOpenData(datas []*FunctionOpenData) []*shared_proto.FunctionOpenDataProto {

	out := make([]*shared_proto.FunctionOpenDataProto, 0, len(datas))
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

func (dAtA *FunctionOpenData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.GuanFuLevel = cOnFigS.GetBuildingData(pArSeR.Uint64("guan_fu_level"))
	if dAtA.GuanFuLevel == nil && pArSeR.Uint64("guan_fu_level") != 0 {
		return errors.Errorf("%s 配置的关联字段[guan_fu_level] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guan_fu_level"), *pArSeR)
	}

	dAtA.Building = cOnFigS.GetBuildingData(pArSeR.Uint64("building"))
	if dAtA.Building == nil && pArSeR.Uint64("building") != 0 {
		return errors.Errorf("%s 配置的关联字段[building] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("building"), *pArSeR)
	}

	dAtA.HeroLevel = cOnFigS.GetHeroLevelData(pArSeR.Uint64("hero_level"))
	if dAtA.HeroLevel == nil && pArSeR.Uint64("hero_level") != 0 {
		return errors.Errorf("%s 配置的关联字段[hero_level] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("hero_level"), *pArSeR)
	}

	dAtA.MainTask = cOnFigS.GetMainTaskData(pArSeR.Uint64("main_task"))
	if dAtA.MainTask == nil && pArSeR.Uint64("main_task") != 0 {
		return errors.Errorf("%s 配置的关联字段[main_task] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("main_task"), *pArSeR)
	}

	dAtA.BaYeStage = cOnFigS.GetBaYeStageData(pArSeR.Uint64("ba_ye_stage"))
	if dAtA.BaYeStage == nil && pArSeR.Uint64("ba_ye_stage") != 0 {
		return errors.Errorf("%s 配置的关联字段[ba_ye_stage] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("ba_ye_stage"), *pArSeR)
	}

	dAtA.TowerFloor = cOnFigS.GetTowerData(pArSeR.Uint64("tower_floor"))
	if dAtA.TowerFloor == nil && pArSeR.Uint64("tower_floor") != 0 {
		return errors.Errorf("%s 配置的关联字段[tower_floor] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("tower_floor"), *pArSeR)
	}

	dAtA.Dungeon = cOnFigS.GetDungeonData(pArSeR.Uint64("dungeon"))
	if dAtA.Dungeon == nil && pArSeR.Uint64("dungeon") != 0 {
		return errors.Errorf("%s 配置的关联字段[dungeon] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("dungeon"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetBaYeStageData(uint64) *taskdata.BaYeStageData
	GetBuildingData(uint64) *domestic_data.BuildingData
	GetDungeonData(uint64) *dungeon.DungeonData
	GetHeroLevelData(uint64) *herodata.HeroLevelData
	GetMainTaskData(uint64) *taskdata.MainTaskData
	GetTowerData(uint64) *towerdata.TowerData
}
