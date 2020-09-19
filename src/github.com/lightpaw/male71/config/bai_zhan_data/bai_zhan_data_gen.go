// AUTO_GEN, DONT MODIFY!!!
package bai_zhan_data

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
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

// start with BaiZhanMiscData ----------------------------------

func LoadBaiZhanMiscData(gos *config.GameObjects) (*BaiZhanMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.BaiZhanMiscDataPath
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

	dAtA, err := NewBaiZhanMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedBaiZhanMiscData(gos *config.GameObjects, dAtA *BaiZhanMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BaiZhanMiscDataPath
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

func NewBaiZhanMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BaiZhanMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBaiZhanMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BaiZhanMiscData{}

	dAtA.WinPoint = pArSeR.Uint64("win_point")
	dAtA.FailPoint = pArSeR.Uint64("fail_point")
	dAtA.InitTimes = 0
	if pArSeR.KeyExist("init_times") {
		dAtA.InitTimes = pArSeR.Uint64("init_times")
	}

	stringKeys = pArSeR.StringArray("recover_times_time", "", false)
	dAtA.RecoverTimesTime = make([]time.Duration, 0, len(stringKeys))
	for _, v := range stringKeys {
		obj, err := config.ParseDuration(v)
		if err != nil {
			return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[recover_times_time] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("recover_times_time"), dAtA)
		}
		dAtA.RecoverTimesTime = append(dAtA.RecoverTimesTime, obj)
	}

	dAtA.RecoverTimes = pArSeR.Uint64Array("recover_times", "", false)
	dAtA.MaxRecord = 20
	if pArSeR.KeyExist("max_record") {
		dAtA.MaxRecord = pArSeR.Uint64("max_record")
	}

	dAtA.ShowRankCount = 6
	if pArSeR.KeyExist("show_rank_count") {
		dAtA.ShowRankCount = pArSeR.Uint64("show_rank_count")
	}

	return dAtA, nil
}

var vAlIdAtOrBaiZhanMiscData = map[string]*config.Validator{

	"win_point":          config.ParseValidator("int>0", "", false, nil, nil),
	"fail_point":         config.ParseValidator("int>0", "", false, nil, nil),
	"init_times":         config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"recover_times_time": config.ParseValidator("string", "", true, nil, nil),
	"recover_times":      config.ParseValidator("int>0,duplicate", "", true, nil, nil),
	"max_record":         config.ParseValidator("int>0", "", false, nil, []string{"20"}),
	"show_rank_count":    config.ParseValidator("int>0", "", false, nil, []string{"6"}),
}

func (dAtA *BaiZhanMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BaiZhanMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BaiZhanMiscData) Encode() *shared_proto.BaiZhanMiscProto {
	out := &shared_proto.BaiZhanMiscProto{}
	out.WinPoint = config.U64ToI32(dAtA.WinPoint)
	out.FailPoint = config.U64ToI32(dAtA.FailPoint)
	out.RecoverTimesTime = timeutil.DurationArrayToSecondArray(dAtA.RecoverTimesTime)
	out.RecoverTimes = config.U64a2I32a(dAtA.RecoverTimes)
	out.MaxRecord = config.U64ToI32(dAtA.MaxRecord)
	out.ShowRankCount = config.U64ToI32(dAtA.ShowRankCount)

	return out
}

func ArrayEncodeBaiZhanMiscData(datas []*BaiZhanMiscData) []*shared_proto.BaiZhanMiscProto {

	out := make([]*shared_proto.BaiZhanMiscProto, 0, len(datas))
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

func (dAtA *BaiZhanMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with JunXianLevelData ----------------------------------

func LoadJunXianLevelData(gos *config.GameObjects) (map[uint64]*JunXianLevelData, map[*JunXianLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.JunXianLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*JunXianLevelData, len(lIsT))
	pArSeRmAp := make(map[*JunXianLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrJunXianLevelData) {
			continue
		}

		dAtA, err := NewJunXianLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedJunXianLevelData(dAtAmAp map[*JunXianLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.JunXianLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetJunXianLevelDataKeyArray(datas []*JunXianLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewJunXianLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*JunXianLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrJunXianLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &JunXianLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	// releated field: Icon
	// releated field: StrongMatchNpcGuardCaptain
	// releated field: NpcGuardCaptain
	// skip field: PrevLevel
	// skip field: NextLevel
	// releated field: ReaddJunXianLevelData
	dAtA.LevelUpPercent = pArSeR.Uint64("level_up_percent")
	dAtA.LevelUpPoint = pArSeR.Uint64("level_up_point")
	dAtA.LevelDownPercent = pArSeR.Uint64("level_down_percent")
	dAtA.LevelDownPoint = pArSeR.Uint64("level_down_point")
	// releated field: DailySalary
	dAtA.DailyHufu = pArSeR.Uint64("daily_hufu")
	// releated field: AssemblyStat
	dAtA.MinKeepLevelCount = pArSeR.Uint64("min_keep_level_count")
	// releated field: CombatScene

	// i18n fields
	dAtA.Name = i18n.NewI18nRef(fIlEnAmE, "name", dAtA.Level, pArSeR.String("name"))

	return dAtA, nil
}

var vAlIdAtOrJunXianLevelData = map[string]*config.Validator{

	"level": config.ParseValidator("int>0", "", false, nil, nil),
	"name":  config.ParseValidator("string>0", "", false, nil, nil),
	"icon":  config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"strong_match_npc_guard_captain": config.ParseValidator("string", "", true, nil, nil),
	"npc_guard_captain":              config.ParseValidator("string", "", true, nil, nil),
	"readd_jun_xian_level_data":      config.ParseValidator("string", "", false, nil, nil),
	"level_up_percent":               config.ParseValidator("uint", "", false, nil, nil),
	"level_up_point":                 config.ParseValidator("uint", "", false, nil, nil),
	"level_down_percent":             config.ParseValidator("uint", "", false, nil, nil),
	"level_down_point":               config.ParseValidator("uint", "", false, nil, nil),
	"daily_salary":                   config.ParseValidator("string", "", false, nil, nil),
	"daily_hufu":                     config.ParseValidator("uint", "", false, nil, nil),
	"assembly_stat":                  config.ParseValidator("string", "", false, nil, nil),
	"min_keep_level_count":           config.ParseValidator("int>0", "", false, nil, nil),
	"combat_scene":                   config.ParseValidator("string", "", false, nil, []string{"CombatScene"}),
}

func (dAtA *JunXianLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *JunXianLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *JunXianLevelData) Encode() *shared_proto.JunXianLevelDataProto {
	out := &shared_proto.JunXianLevelDataProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	if dAtA.Name != nil {
		out.Name = dAtA.Name.Encode()
	}
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.LevelUpPercent = config.U64ToI32(dAtA.LevelUpPercent)
	out.LevelUpPoint = config.U64ToI32(dAtA.LevelUpPoint)
	out.LevelDownPercent = config.U64ToI32(dAtA.LevelDownPercent)
	out.LevelDownPoint = config.U64ToI32(dAtA.LevelDownPoint)
	if dAtA.DailySalary != nil {
		out.DailySalary = dAtA.DailySalary.Encode()
	}
	out.DailyHufu = config.U64ToI32(dAtA.DailyHufu)
	if dAtA.AssemblyStat != nil {
		out.AssemblyStat = dAtA.AssemblyStat.Encode()
	}

	return out
}

func ArrayEncodeJunXianLevelData(datas []*JunXianLevelData) []*shared_proto.JunXianLevelDataProto {

	out := make([]*shared_proto.JunXianLevelDataProto, 0, len(datas))
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

func (dAtA *JunXianLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	uint64Keys = pArSeR.Uint64Array("strong_match_npc_guard_captain", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetMonsterMasterData(v)
		if obj != nil {
			dAtA.StrongMatchNpcGuardCaptain = append(dAtA.StrongMatchNpcGuardCaptain, obj)
		} else if v != 0 {
			return errors.Errorf("%s 配置的关联字段[strong_match_npc_guard_captain] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("strong_match_npc_guard_captain"), *pArSeR)
		}
	}

	uint64Keys = pArSeR.Uint64Array("npc_guard_captain", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetMonsterMasterData(v)
		if obj != nil {
			dAtA.NpcGuardCaptain = append(dAtA.NpcGuardCaptain, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[npc_guard_captain] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("npc_guard_captain"), *pArSeR)
		}
	}

	dAtA.ReaddJunXianLevelData = cOnFigS.GetJunXianLevelData(pArSeR.Uint64("readd_jun_xian_level_data"))
	if dAtA.ReaddJunXianLevelData == nil {
		return errors.Errorf("%s 配置的关联字段[readd_jun_xian_level_data] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("readd_jun_xian_level_data"), *pArSeR)
	}

	dAtA.DailySalary = cOnFigS.GetPrize(pArSeR.Int("daily_salary"))
	if dAtA.DailySalary == nil {
		return errors.Errorf("%s 配置的关联字段[daily_salary] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("daily_salary"), *pArSeR)
	}

	dAtA.AssemblyStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("assembly_stat"))
	if dAtA.AssemblyStat == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_stat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_stat"), *pArSeR)
	}

	if pArSeR.KeyExist("combat_scene") {
		dAtA.CombatScene = cOnFigS.GetCombatScene(pArSeR.String("combat_scene"))
	} else {
		dAtA.CombatScene = cOnFigS.GetCombatScene("CombatScene")
	}
	if dAtA.CombatScene == nil {
		return errors.Errorf("%s 配置的关联字段[combat_scene] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("combat_scene"), *pArSeR)
	}

	return nil
}

// start with JunXianPrizeData ----------------------------------

func LoadJunXianPrizeData(gos *config.GameObjects) (map[uint64]*JunXianPrizeData, map[*JunXianPrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.JunXianPrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*JunXianPrizeData, len(lIsT))
	pArSeRmAp := make(map[*JunXianPrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrJunXianPrizeData) {
			continue
		}

		dAtA, err := NewJunXianPrizeData(fIlEnAmE, pArSeR)
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

func SetRelatedJunXianPrizeData(dAtAmAp map[*JunXianPrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.JunXianPrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetJunXianPrizeDataKeyArray(datas []*JunXianPrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewJunXianPrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*JunXianPrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrJunXianPrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &JunXianPrizeData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: LevelData
	dAtA.Point = pArSeR.Uint64("point")
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrJunXianPrizeData = map[string]*config.Validator{

	"id":    config.ParseValidator("int>0", "", false, nil, nil),
	"level": config.ParseValidator("string", "", false, nil, nil),
	"point": config.ParseValidator("uint", "", false, nil, nil),
	"prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *JunXianPrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *JunXianPrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *JunXianPrizeData) Encode() *shared_proto.JunXianLevelPrizeProto {
	out := &shared_proto.JunXianLevelPrizeProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.LevelData != nil {
		out.Level = config.U64ToI32(dAtA.LevelData.Level)
	}
	out.Point = config.U64ToI32(dAtA.Point)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeJunXianPrizeData(datas []*JunXianPrizeData) []*shared_proto.JunXianLevelPrizeProto {

	out := make([]*shared_proto.JunXianLevelPrizeProto, 0, len(datas))
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

func (dAtA *JunXianPrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.LevelData = cOnFigS.GetJunXianLevelData(pArSeR.Uint64("level"))
	if dAtA.LevelData == nil {
		return errors.Errorf("%s 配置的关联字段[level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("level"), *pArSeR)
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetCombatScene(string) *scene.CombatScene
	GetIcon(string) *icon.Icon
	GetJunXianLevelData(uint64) *JunXianLevelData
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
	GetPrize(int) *resdata.Prize
	GetSpriteStat(uint64) *data.SpriteStat
}
