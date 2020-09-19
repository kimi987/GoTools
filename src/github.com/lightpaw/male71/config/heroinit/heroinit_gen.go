// AUTO_GEN, DONT MODIFY!!!
package heroinit

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/taskdata"
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

// start with HeroCreateData ----------------------------------

func LoadHeroCreateData(gos *config.GameObjects) (*HeroCreateData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.HeroCreateDataPath
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

	dAtA, err := NewHeroCreateData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedHeroCreateData(gos *config.GameObjects, dAtA *HeroCreateData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.HeroCreateDataPath
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

func NewHeroCreateData(fIlEnAmE string, pArSeR *config.ObjectParser) (*HeroCreateData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrHeroCreateData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &HeroCreateData{}

	dAtA.Gold = pArSeR.Uint64("gold")
	dAtA.Food = pArSeR.Uint64("food")
	dAtA.Wood = pArSeR.Uint64("wood")
	dAtA.Stone = pArSeR.Uint64("stone")
	dAtA.NewSoldier = pArSeR.Uint64("new_soldier")
	// skip field: Prosperity
	// releated field: Captain

	return dAtA, nil
}

var vAlIdAtOrHeroCreateData = map[string]*config.Validator{

	"gold":        config.ParseValidator("uint", "", false, nil, nil),
	"food":        config.ParseValidator("uint", "", false, nil, nil),
	"wood":        config.ParseValidator("uint", "", false, nil, nil),
	"stone":       config.ParseValidator("uint", "", false, nil, nil),
	"new_soldier": config.ParseValidator("uint", "", false, nil, nil),
	"captain":     config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *HeroCreateData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("captain", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetCaptainData(v)
		if obj != nil {
			dAtA.Captain = append(dAtA.Captain, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[captain] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain"), *pArSeR)
		}
	}

	return nil
}

// start with HeroInitData ----------------------------------

func LoadHeroInitData(gos *config.GameObjects) (*HeroInitData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.HeroInitDataPath
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

	dAtA, err := NewHeroInitData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedHeroInitData(gos *config.GameObjects, dAtA *HeroInitData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.HeroInitDataPath
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

func NewHeroInitData(fIlEnAmE string, pArSeR *config.ObjectParser) (*HeroInitData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrHeroInitData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &HeroInitData{}

	dAtA.BuildingWorkerMaxCount = 2
	if pArSeR.KeyExist("building_worker_max_count") {
		dAtA.BuildingWorkerMaxCount = pArSeR.Uint64("building_worker_max_count")
	}

	dAtA.TechWorkerMaxCount = 1
	if pArSeR.KeyExist("tech_worker_max_count") {
		dAtA.TechWorkerMaxCount = pArSeR.Uint64("tech_worker_max_count")
	}

	// releated field: FirstLevelSoldierData
	// releated field: FirstLevelHeroData
	// releated field: FirstMainTask
	// releated field: FirstBaYeStageData
	// releated field: FirstTitleData
	// releated field: FirstLevelCountdownPrize
	dAtA.CaptainIndexCount = 15
	if pArSeR.KeyExist("captain_index_count") {
		dAtA.CaptainIndexCount = pArSeR.Uint64("captain_index_count")
	}

	// skip field: Building
	// skip field: MinLevelTieJiangPuData
	// skip field: DefaultForgingTimes
	// skip field: MaxDepotEquipCapacity
	// skip field: TempDepotExpireDuration
	// skip field: GuildDonateTypeCount
	dAtA.TroopCaptainCount = 5
	if pArSeR.KeyExist("troop_captain_count") {
		dAtA.TroopCaptainCount = pArSeR.Uint64("troop_captain_count")
	}

	// skip field: MaxTagColorType
	// skip field: MaxTagRecordCount
	// skip field: MaxShowForViewCount
	// skip field: StrategyRestoreDuration
	// skip field: DungeonChallengeDefaultTimes
	// skip field: DungeonChallengeMaxTimes
	// skip field: DungeonChallengeRecoverDuration
	// skip field: JunYingRecoveryDuration
	// skip field: JunYingMaxTimes
	// skip field: JunYingDefaultTimes
	// skip field: MaxFavoritePosCount
	// skip field: ActiveDegreeTaskDatas
	// skip field: BwzlTaskDatas
	// skip field: AchieveTaskDatas
	// skip field: FunctionOpenDataArray
	// skip field: DefaultHead
	// skip field: DefaultBody
	// skip field: PveTroopDatas
	// skip field: MultiLevelNpcInitTimes
	// skip field: MultiLevelNpcMaxTimes
	// skip field: MultiLevelNpcRecoveryDuration
	// skip field: InvaseHeroInitTimes
	// skip field: InvaseHeroMaxTimes
	// skip field: InvaseHeroRecoveryDuration
	// skip field: JunTuanNpcInitTimes
	// skip field: JunTuanNpcMaxTimes
	// skip field: JunTuanNpcRecoveryDuration
	// skip field: WorkshopOutputInitTimes
	// skip field: WorkshopOutputMaxTimes
	// skip field: WorkshopOutputRecoveryDuration
	// skip field: ZhengWuMiscData
	// skip field: DefaultSettings
	// skip field: BuildingInitEffect
	// skip field: BaowuLogLimit
	// skip field: CaptainInitData

	return dAtA, nil
}

var vAlIdAtOrHeroInitData = map[string]*config.Validator{

	"building_worker_max_count":   config.ParseValidator("int>0", "", false, nil, []string{"2"}),
	"tech_worker_max_count":       config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"first_level_soldier_data":    config.ParseValidator("string", "", false, nil, []string{"1"}),
	"first_level_hero_data":       config.ParseValidator("string", "", false, nil, []string{"1"}),
	"first_main_task":             config.ParseValidator("string", "", false, nil, []string{"1"}),
	"first_ba_ye_stage_data":      config.ParseValidator("string", "", false, nil, []string{"1"}),
	"first_title_data":            config.ParseValidator("string", "", false, nil, []string{"1"}),
	"first_level_countdown_prize": config.ParseValidator("string", "", false, nil, []string{"1"}),
	"captain_index_count":         config.ParseValidator("int>0", "", false, nil, []string{"15"}),
	"troop_captain_count":         config.ParseValidator("int>0", "", false, nil, []string{"5"}),
}

func (dAtA *HeroInitData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("first_level_soldier_data") {
		dAtA.FirstLevelSoldierData = cOnFigS.GetSoldierLevelData(pArSeR.Uint64("first_level_soldier_data"))
	} else {
		dAtA.FirstLevelSoldierData = cOnFigS.GetSoldierLevelData(1)
	}
	if dAtA.FirstLevelSoldierData == nil {
		return errors.Errorf("%s 配置的关联字段[first_level_soldier_data] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_level_soldier_data"), *pArSeR)
	}

	if pArSeR.KeyExist("first_level_hero_data") {
		dAtA.FirstLevelHeroData = cOnFigS.GetHeroLevelData(pArSeR.Uint64("first_level_hero_data"))
	} else {
		dAtA.FirstLevelHeroData = cOnFigS.GetHeroLevelData(1)
	}
	if dAtA.FirstLevelHeroData == nil {
		return errors.Errorf("%s 配置的关联字段[first_level_hero_data] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_level_hero_data"), *pArSeR)
	}

	if pArSeR.KeyExist("first_main_task") {
		dAtA.FirstMainTask = cOnFigS.GetMainTaskData(pArSeR.Uint64("first_main_task"))
	} else {
		dAtA.FirstMainTask = cOnFigS.GetMainTaskData(1)
	}
	if dAtA.FirstMainTask == nil {
		return errors.Errorf("%s 配置的关联字段[first_main_task] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_main_task"), *pArSeR)
	}

	if pArSeR.KeyExist("first_ba_ye_stage_data") {
		dAtA.FirstBaYeStageData = cOnFigS.GetBaYeStageData(pArSeR.Uint64("first_ba_ye_stage_data"))
	} else {
		dAtA.FirstBaYeStageData = cOnFigS.GetBaYeStageData(1)
	}
	if dAtA.FirstBaYeStageData == nil {
		return errors.Errorf("%s 配置的关联字段[first_ba_ye_stage_data] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_ba_ye_stage_data"), *pArSeR)
	}

	if pArSeR.KeyExist("first_title_data") {
		dAtA.FirstTitleData = cOnFigS.GetTitleData(pArSeR.Uint64("first_title_data"))
	} else {
		dAtA.FirstTitleData = cOnFigS.GetTitleData(1)
	}
	if dAtA.FirstTitleData == nil {
		return errors.Errorf("%s 配置的关联字段[first_title_data] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_title_data"), *pArSeR)
	}

	if pArSeR.KeyExist("first_level_countdown_prize") {
		dAtA.FirstLevelCountdownPrize = cOnFigS.GetCountdownPrizeData(pArSeR.Uint64("first_level_countdown_prize"))
	} else {
		dAtA.FirstLevelCountdownPrize = cOnFigS.GetCountdownPrizeData(1)
	}
	if dAtA.FirstLevelCountdownPrize == nil {
		return errors.Errorf("%s 配置的关联字段[first_level_countdown_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_level_countdown_prize"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetBaYeStageData(uint64) *taskdata.BaYeStageData
	GetCaptainData(uint64) *captain.CaptainData
	GetCountdownPrizeData(uint64) *domestic_data.CountdownPrizeData
	GetHeroLevelData(uint64) *herodata.HeroLevelData
	GetMainTaskData(uint64) *taskdata.MainTaskData
	GetSoldierLevelData(uint64) *domestic_data.SoldierLevelData
	GetTitleData(uint64) *taskdata.TitleData
}
