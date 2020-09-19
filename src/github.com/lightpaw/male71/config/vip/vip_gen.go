// AUTO_GEN, DONT MODIFY!!!
package vip

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/resdata"
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

// start with VipContinueDaysData ----------------------------------

func LoadVipContinueDaysData(gos *config.GameObjects) (map[uint64]*VipContinueDaysData, map[*VipContinueDaysData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.VipContinueDaysDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*VipContinueDaysData, len(lIsT))
	pArSeRmAp := make(map[*VipContinueDaysData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrVipContinueDaysData) {
			continue
		}

		dAtA, err := NewVipContinueDaysData(fIlEnAmE, pArSeR)
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

func SetRelatedVipContinueDaysData(dAtAmAp map[*VipContinueDaysData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.VipContinueDaysDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetVipContinueDaysDataKeyArray(datas []*VipContinueDaysData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewVipContinueDaysData(fIlEnAmE string, pArSeR *config.ObjectParser) (*VipContinueDaysData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrVipContinueDaysData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &VipContinueDaysData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Days = pArSeR.Uint64Array("days", "", false)
	dAtA.Exp = pArSeR.Uint64Array("exp", "", false)

	return dAtA, nil
}

var vAlIdAtOrVipContinueDaysData = map[string]*config.Validator{

	"level": config.ParseValidator("uint", "", false, nil, nil),
	"days":  config.ParseValidator(",duplicate", "", true, nil, nil),
	"exp":   config.ParseValidator("uint,duplicate", "", true, nil, nil),
}

func (dAtA *VipContinueDaysData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *VipContinueDaysData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *VipContinueDaysData) Encode() *shared_proto.VipContinueDaysDataProto {
	out := &shared_proto.VipContinueDaysDataProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Days = config.U64a2I32a(dAtA.Days)
	out.Exp = config.U64a2I32a(dAtA.Exp)

	return out
}

func ArrayEncodeVipContinueDaysData(datas []*VipContinueDaysData) []*shared_proto.VipContinueDaysDataProto {

	out := make([]*shared_proto.VipContinueDaysDataProto, 0, len(datas))
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

func (dAtA *VipContinueDaysData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with VipLevelData ----------------------------------

func LoadVipLevelData(gos *config.GameObjects) (map[uint64]*VipLevelData, map[*VipLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.VipLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*VipLevelData, len(lIsT))
	pArSeRmAp := make(map[*VipLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrVipLevelData) {
			continue
		}

		dAtA, err := NewVipLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedVipLevelData(dAtAmAp map[*VipLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.VipLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetVipLevelDataKeyArray(datas []*VipLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewVipLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*VipLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrVipLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &VipLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Name = pArSeR.String("name")
	// releated field: Icon
	dAtA.UpgradeExp = pArSeR.Uint64("upgrade_exp")
	dAtA.DailyExp = pArSeR.Uint64("daily_exp")
	// releated field: DailyPrize
	// releated field: ShowDailyPrizeCost
	// releated field: LevelPrize
	// releated field: LevelPrizeCost
	// releated field: ShowLevelPrizeCost
	// skip field: ContinuteDays
	// skip field: NextLevelData
	dAtA.BuyProsperity = pArSeR.Bool("buy_prosperity")
	dAtA.JiuGuanAutoMax = pArSeR.Bool("jiu_guan_auto_max")
	dAtA.JiuGuanCostRefreshCount = pArSeR.Uint64("jiu_guan_cost_refresh_count")
	dAtA.JiuGuanCostRefreshInfinite = pArSeR.Bool("jiu_guan_cost_refresh_infinite")
	dAtA.JiuGuanQuickConsult = pArSeR.Bool("jiu_guan_quick_consult")
	dAtA.CaptainTrainCoef = pArSeR.Float64("captain_train_coef")
	dAtA.CaptainTrainCapacity, err = config.ParseDuration(pArSeR.String("captain_train_capacity"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[captain_train_capacity] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("captain_train_capacity"), dAtA)
	}

	dAtA.WallAutoFullSoldier = pArSeR.Bool("wall_auto_full_soldier")
	dAtA.BuySpMaxTimes = pArSeR.Uint64("buy_sp_max_times")
	dAtA.DungeonMaxCostTimesLimit = pArSeR.Uint64("dungeon_max_cost_times_limit")
	dAtA.WorkerUnlockPos = pArSeR.Uint64("worker_unlock_pos")
	dAtA.InvadeMultiLevelMonsterOnceCount = pArSeR.Uint64("invade_multi_level_monster_once_count")
	dAtA.GuildPrizeOneKeyCollect = pArSeR.Bool("guild_prize_one_key_collect")
	dAtA.AddBlackMarketRefreshTimes = pArSeR.Uint64("add_black_market_refresh_times")
	dAtA.FishingCaptainProbability = pArSeR.Bool("fishing_captain_probability")
	dAtA.ShowRegionHome = false
	if pArSeR.KeyExist("show_region_home") {
		dAtA.ShowRegionHome = pArSeR.Bool("show_region_home")
	}

	dAtA.ShowRegionSign = false
	if pArSeR.KeyExist("show_region_sign") {
		dAtA.ShowRegionSign = pArSeR.Bool("show_region_sign")
	}

	dAtA.ShowRegionTitle = false
	if pArSeR.KeyExist("show_region_title") {
		dAtA.ShowRegionTitle = pArSeR.Bool("show_region_title")
	}

	dAtA.ShowHeadFrame = false
	if pArSeR.KeyExist("show_head_frame") {
		dAtA.ShowHeadFrame = pArSeR.Bool("show_head_frame")
	}

	dAtA.ZhengWuAutoCompleted = false
	if pArSeR.KeyExist("zheng_wu_auto_completed") {
		dAtA.ZhengWuAutoCompleted = pArSeR.Bool("zheng_wu_auto_completed")
	}

	dAtA.WorkshopAutoCompleted = false
	if pArSeR.KeyExist("workshop_auto_completed") {
		dAtA.WorkshopAutoCompleted = pArSeR.Bool("workshop_auto_completed")
	}

	dAtA.CollectDailySp = false
	if pArSeR.KeyExist("collect_daily_sp") {
		dAtA.CollectDailySp = pArSeR.Bool("collect_daily_sp")
	}

	return dAtA, nil
}

var vAlIdAtOrVipLevelData = map[string]*config.Validator{

	"level":                                 config.ParseValidator("uint", "", false, nil, nil),
	"name":                                  config.ParseValidator("string", "", false, nil, nil),
	"icon":                                  config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"upgrade_exp":                           config.ParseValidator("int>0", "", false, nil, nil),
	"daily_exp":                             config.ParseValidator("uint", "", false, nil, nil),
	"daily_prize":                           config.ParseValidator("string", "", false, nil, nil),
	"show_daily_prize_cost":                 config.ParseValidator("string", "", false, nil, nil),
	"level_prize":                           config.ParseValidator("string", "", false, nil, nil),
	"level_prize_cost":                      config.ParseValidator("string", "", false, nil, nil),
	"show_level_prize_cost":                 config.ParseValidator("string", "", false, nil, nil),
	"buy_prosperity":                        config.ParseValidator("bool", "", false, nil, nil),
	"jiu_guan_auto_max":                     config.ParseValidator("bool", "", false, nil, nil),
	"jiu_guan_cost_refresh_count":           config.ParseValidator("uint", "", false, nil, nil),
	"jiu_guan_cost_refresh_infinite":        config.ParseValidator("bool", "", false, nil, nil),
	"jiu_guan_quick_consult":                config.ParseValidator("bool", "", false, nil, nil),
	"captain_train_coef":                    config.ParseValidator("float64", "", false, nil, nil),
	"captain_train_capacity":                config.ParseValidator("string", "", false, nil, nil),
	"wall_auto_full_soldier":                config.ParseValidator("bool", "", false, nil, nil),
	"buy_sp_max_times":                      config.ParseValidator("uint", "", false, nil, nil),
	"dungeon_max_cost_times_limit":          config.ParseValidator("uint", "", false, nil, nil),
	"worker_unlock_pos":                     config.ParseValidator("uint", "", false, nil, nil),
	"invade_multi_level_monster_once_count": config.ParseValidator("uint", "", false, nil, nil),
	"guild_prize_one_key_collect":           config.ParseValidator("bool", "", false, nil, nil),
	"add_black_market_refresh_times":        config.ParseValidator("uint", "", false, nil, nil),
	"fishing_captain_probability":           config.ParseValidator("bool", "", false, nil, nil),
	"show_region_home":                      config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"show_region_sign":                      config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"show_region_title":                     config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"show_head_frame":                       config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"zheng_wu_auto_completed":               config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"workshop_auto_completed":               config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"collect_daily_sp":                      config.ParseValidator("bool", "", false, nil, []string{"false"}),
}

func (dAtA *VipLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *VipLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *VipLevelData) Encode() *shared_proto.VipLevelDataProto {
	out := &shared_proto.VipLevelDataProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Name = dAtA.Name
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.UpgradeExp = config.U64ToI32(dAtA.UpgradeExp)
	out.DailyExp = config.U64ToI32(dAtA.DailyExp)
	if dAtA.DailyPrize != nil {
		out.DailyPrize = dAtA.DailyPrize.Encode()
	}
	if dAtA.ShowDailyPrizeCost != nil {
		out.ShowDailyPrizeCost = dAtA.ShowDailyPrizeCost.Encode()
	}
	if dAtA.LevelPrize != nil {
		out.LevelPrize = dAtA.LevelPrize.Encode()
	}
	if dAtA.LevelPrizeCost != nil {
		out.LevelPrizeCost = dAtA.LevelPrizeCost.Encode()
	}
	if dAtA.ShowLevelPrizeCost != nil {
		out.ShowLevelPrizeCost = dAtA.ShowLevelPrizeCost.Encode()
	}
	out.BuyProsperity = dAtA.BuyProsperity
	out.JiuGuanAutoMax = dAtA.JiuGuanAutoMax
	out.JiuGuanCostRefreshCount = config.U64ToI32(dAtA.JiuGuanCostRefreshCount)
	out.JiuGuanCostRefreshInfinite = dAtA.JiuGuanCostRefreshInfinite
	out.JiuGuanQuickConsult = dAtA.JiuGuanQuickConsult
	out.CaptainTrainCoef = config.F64ToI32X1000(dAtA.CaptainTrainCoef)
	out.CaptainTrainCapacity = config.Duration2I32Seconds(dAtA.CaptainTrainCapacity)
	out.WallAutoFullSoldier = dAtA.WallAutoFullSoldier
	out.BuySpMaxTimes = config.U64ToI32(dAtA.BuySpMaxTimes)
	out.DungeonMaxCostTimesLimit = config.U64ToI32(dAtA.DungeonMaxCostTimesLimit)
	out.WorkerUnlockPos = config.U64ToI32(dAtA.WorkerUnlockPos)
	out.InvadeMultiLevelMonsterOnceCount = config.U64ToI32(dAtA.InvadeMultiLevelMonsterOnceCount)
	out.GuildPrizeOneKeyCollect = dAtA.GuildPrizeOneKeyCollect
	out.AddBlackMarketRefreshTimes = config.U64ToI32(dAtA.AddBlackMarketRefreshTimes)
	out.FishingCaptainProbability = dAtA.FishingCaptainProbability
	out.ShowRegionHome = dAtA.ShowRegionHome
	out.ShowRegionSign = dAtA.ShowRegionSign
	out.ShowRegionTitle = dAtA.ShowRegionTitle
	out.ShowHeadFrame = dAtA.ShowHeadFrame
	out.ZhengWuAutoCompleted = dAtA.ZhengWuAutoCompleted
	out.WorkshopAutoCompleted = dAtA.WorkshopAutoCompleted
	out.CollectDailySp = dAtA.CollectDailySp

	return out
}

func ArrayEncodeVipLevelData(datas []*VipLevelData) []*shared_proto.VipLevelDataProto {

	out := make([]*shared_proto.VipLevelDataProto, 0, len(datas))
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

func (dAtA *VipLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("icon") {
		dAtA.Icon = cOnFigS.GetIcon(pArSeR.String("icon"))
	} else {
		dAtA.Icon = cOnFigS.GetIcon("Icon")
	}
	if dAtA.Icon == nil {
		return errors.Errorf("%s 配置的关联字段[icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("icon"), *pArSeR)
	}

	dAtA.DailyPrize = cOnFigS.GetPrize(pArSeR.Int("daily_prize"))
	if dAtA.DailyPrize == nil && pArSeR.Int("daily_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[daily_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("daily_prize"), *pArSeR)
	}

	dAtA.ShowDailyPrizeCost = cOnFigS.GetCost(pArSeR.Int("show_daily_prize_cost"))
	if dAtA.ShowDailyPrizeCost == nil && pArSeR.Int("show_daily_prize_cost") != 0 {
		return errors.Errorf("%s 配置的关联字段[show_daily_prize_cost] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_daily_prize_cost"), *pArSeR)
	}

	dAtA.LevelPrize = cOnFigS.GetPrize(pArSeR.Int("level_prize"))
	if dAtA.LevelPrize == nil && pArSeR.Int("level_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[level_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("level_prize"), *pArSeR)
	}

	dAtA.LevelPrizeCost = cOnFigS.GetCost(pArSeR.Int("level_prize_cost"))
	if dAtA.LevelPrizeCost == nil && pArSeR.Int("level_prize_cost") != 0 {
		return errors.Errorf("%s 配置的关联字段[level_prize_cost] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("level_prize_cost"), *pArSeR)
	}

	dAtA.ShowLevelPrizeCost = cOnFigS.GetCost(pArSeR.Int("show_level_prize_cost"))
	if dAtA.ShowLevelPrizeCost == nil && pArSeR.Int("show_level_prize_cost") != 0 {
		return errors.Errorf("%s 配置的关联字段[show_level_prize_cost] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_level_prize_cost"), *pArSeR)
	}

	return nil
}

// start with VipMiscData ----------------------------------

func LoadVipMiscData(gos *config.GameObjects) (*VipMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.VipMiscDataPath
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

	dAtA, err := NewVipMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedVipMiscData(gos *config.GameObjects, dAtA *VipMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.VipMiscDataPath
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

func NewVipMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*VipMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrVipMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &VipMiscData{}

	dAtA.CollectVipDailyExpMinHeroLevel = pArSeR.Uint64("collect_vip_daily_exp_min_hero_level")
	// releated field: DungeonTimesCost
	dAtA.DungeonTimesEachBuy = pArSeR.Uint64Array("dungeon_times_each_buy", "", false)

	return dAtA, nil
}

var vAlIdAtOrVipMiscData = map[string]*config.Validator{

	"collect_vip_daily_exp_min_hero_level": config.ParseValidator("int>0", "", false, nil, nil),
	"dungeon_times_cost":                   config.ParseValidator("string", "", true, nil, nil),
	"dungeon_times_each_buy":               config.ParseValidator("uint,duplicate", "", true, nil, nil),
}

func (dAtA *VipMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *VipMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *VipMiscData) Encode() *shared_proto.VipMiscDataProto {
	out := &shared_proto.VipMiscDataProto{}
	out.CollectVipDailyExpMinHeroLevel = config.U64ToI32(dAtA.CollectVipDailyExpMinHeroLevel)
	if dAtA.DungeonTimesCost != nil {
		out.DungeonTimesCost = resdata.ArrayEncodeCost(dAtA.DungeonTimesCost)
	}
	out.DungeonTimesEachBuy = config.U64a2I32a(dAtA.DungeonTimesEachBuy)

	return out
}

func ArrayEncodeVipMiscData(datas []*VipMiscData) []*shared_proto.VipMiscDataProto {

	out := make([]*shared_proto.VipMiscDataProto, 0, len(datas))
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

func (dAtA *VipMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	intKeys = pArSeR.IntArray("dungeon_times_cost", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetCost(v)
		if obj != nil {
			dAtA.DungeonTimesCost = append(dAtA.DungeonTimesCost, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[dungeon_times_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("dungeon_times_cost"), *pArSeR)
		}
	}

	return nil
}

type related_configs interface {
	GetCost(int) *resdata.Cost
	GetIcon(string) *icon.Icon
	GetPrize(int) *resdata.Prize
}
