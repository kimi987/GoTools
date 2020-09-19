// AUTO_GEN, DONT MODIFY!!!
package towerdata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
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

// start with SecretTowerData ----------------------------------

func LoadSecretTowerData(gos *config.GameObjects) (map[uint64]*SecretTowerData, map[*SecretTowerData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.SecretTowerDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*SecretTowerData, len(lIsT))
	pArSeRmAp := make(map[*SecretTowerData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrSecretTowerData) {
			continue
		}

		dAtA, err := NewSecretTowerData(fIlEnAmE, pArSeR)
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

func SetRelatedSecretTowerData(dAtAmAp map[*SecretTowerData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SecretTowerDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetSecretTowerDataKeyArray(datas []*SecretTowerData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewSecretTowerData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SecretTowerData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSecretTowerData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SecretTowerData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// releated field: FirstPassPrize
	// releated field: Plunder
	// releated field: ShowPrize
	dAtA.Image = pArSeR.String("image")
	// releated field: SuperPlunder
	dAtA.SuperPlunderRate = pArSeR.Uint64("super_plunder_rate")
	// releated field: SuperShowPrize
	dAtA.GuildHelpContribution = pArSeR.Uint64("guild_help_contribution")
	dAtA.StartProtectDuration, err = config.ParseDuration(pArSeR.String("start_protect_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[start_protect_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("start_protect_duration"), dAtA)
	}

	dAtA.MaxAttackerCount = pArSeR.Uint64("max_attacker_count")
	dAtA.MinAttackerCount = pArSeR.Uint64("min_attacker_count")
	dAtA.ConcurrentFightCount = pArSeR.Uint64("concurrent_fight_count")
	dAtA.MaxAttackerContinuewWinTimes = pArSeR.Uint64("max_attacker_continuew_win_times")
	dAtA.MaxDefenserContinuewWinTimes = 0
	if pArSeR.KeyExist("max_defenser_continuew_win_times") {
		dAtA.MaxDefenserContinuewWinTimes = pArSeR.Uint64("max_defenser_continuew_win_times")
	}

	dAtA.Desc = pArSeR.String("desc")
	dAtA.MonsterLeaderId = pArSeR.Uint64("monster_leader_id")
	// releated field: Monster
	// releated field: CombatScene
	// releated field: UnlockTowerData
	if pArSeR.KeyExist("team_expire_duration") {
		dAtA.TeamExpireDuration, err = config.ParseDuration(pArSeR.String("team_expire_duration"))
	} else {
		dAtA.TeamExpireDuration, err = config.ParseDuration("15m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[team_expire_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("team_expire_duration"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrSecretTowerData = map[string]*config.Validator{

	"id":                               config.ParseValidator("int>0", "", false, nil, nil),
	"name":                             config.ParseValidator("string>0", "", false, nil, nil),
	"first_pass_prize":                 config.ParseValidator("string", "", false, nil, nil),
	"plunder":                          config.ParseValidator("string", "", false, nil, nil),
	"show_prize":                       config.ParseValidator("string", "", false, nil, nil),
	"image":                            config.ParseValidator("string", "", false, nil, nil),
	"super_plunder":                    config.ParseValidator("string", "", false, nil, nil),
	"super_plunder_rate":               config.ParseValidator("int>0", "", false, nil, nil),
	"super_show_prize":                 config.ParseValidator("string", "", false, nil, nil),
	"guild_help_contribution":          config.ParseValidator("int>0", "", false, nil, nil),
	"start_protect_duration":           config.ParseValidator("string", "", false, nil, nil),
	"max_attacker_count":               config.ParseValidator("int>0", "", false, nil, nil),
	"min_attacker_count":               config.ParseValidator("uint", "", false, nil, nil),
	"concurrent_fight_count":           config.ParseValidator("int>0", "", false, nil, nil),
	"max_attacker_continuew_win_times": config.ParseValidator("uint", "", false, nil, nil),
	"max_defenser_continuew_win_times": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"desc":                 config.ParseValidator("string", "", false, nil, nil),
	"monster_leader_id":    config.ParseValidator("int>0", "", false, nil, nil),
	"monster":              config.ParseValidator("string", "", true, nil, nil),
	"combat_scene":         config.ParseValidator("string", "", false, nil, []string{"CombatScene"}),
	"unlock_tower_data":    config.ParseValidator("string", "", false, nil, nil),
	"team_expire_duration": config.ParseValidator("string", "", false, nil, []string{"15m"}),
}

func (dAtA *SecretTowerData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SecretTowerData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SecretTowerData) Encode() *shared_proto.SecretTowerDataProto {
	out := &shared_proto.SecretTowerDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	if dAtA.FirstPassPrize != nil {
		out.FirstPassPrize = dAtA.FirstPassPrize.Encode()
	}
	if dAtA.ShowPrize != nil {
		out.Prize = dAtA.ShowPrize.Encode()
	}
	out.Image = dAtA.Image
	if dAtA.SuperShowPrize != nil {
		out.SuperPrize = dAtA.SuperShowPrize.Encode()
	}
	out.GuildHelpContribution = config.U64ToI32(dAtA.GuildHelpContribution)
	out.MaxAttackerCount = config.U64ToI32(dAtA.MaxAttackerCount)
	out.MinAttackerCount = config.U64ToI32(dAtA.MinAttackerCount)
	out.MaxAttackerContinuewWinTimes = config.U64ToI32(dAtA.MaxAttackerContinuewWinTimes)
	out.Desc = dAtA.Desc
	out.MonsterLeaderId = config.U64ToI32(dAtA.MonsterLeaderId)
	if dAtA.Monster != nil {
		out.Monster = monsterdata.ArrayEncodeMonsterMasterData(dAtA.Monster)
	}
	if dAtA.UnlockTowerData != nil {
		out.UnlockTowerData = config.U64ToI32(dAtA.UnlockTowerData.Floor)
	}

	return out
}

func ArrayEncodeSecretTowerData(datas []*SecretTowerData) []*shared_proto.SecretTowerDataProto {

	out := make([]*shared_proto.SecretTowerDataProto, 0, len(datas))
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

func (dAtA *SecretTowerData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.FirstPassPrize = cOnFigS.GetPrize(pArSeR.Int("first_pass_prize"))
	if dAtA.FirstPassPrize == nil {
		return errors.Errorf("%s 配置的关联字段[first_pass_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_pass_prize"), *pArSeR)
	}

	dAtA.Plunder = cOnFigS.GetPlunder(pArSeR.Uint64("plunder"))
	if dAtA.Plunder == nil {
		return errors.Errorf("%s 配置的关联字段[plunder] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder"), *pArSeR)
	}

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil {
		return errors.Errorf("%s 配置的关联字段[show_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
	}

	dAtA.SuperPlunder = cOnFigS.GetPlunder(pArSeR.Uint64("super_plunder"))
	if dAtA.SuperPlunder == nil && pArSeR.Uint64("super_plunder") != 0 {
		return errors.Errorf("%s 配置的关联字段[super_plunder] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("super_plunder"), *pArSeR)
	}

	dAtA.SuperShowPrize = cOnFigS.GetPrize(pArSeR.Int("super_show_prize"))
	if dAtA.SuperShowPrize == nil && pArSeR.Int("super_show_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[super_show_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("super_show_prize"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("monster", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetMonsterMasterData(v)
		if obj != nil {
			dAtA.Monster = append(dAtA.Monster, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[monster] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("monster"), *pArSeR)
		}
	}

	if pArSeR.KeyExist("combat_scene") {
		dAtA.CombatScene = cOnFigS.GetCombatScene(pArSeR.String("combat_scene"))
	} else {
		dAtA.CombatScene = cOnFigS.GetCombatScene("CombatScene")
	}
	if dAtA.CombatScene == nil {
		return errors.Errorf("%s 配置的关联字段[combat_scene] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("combat_scene"), *pArSeR)
	}

	dAtA.UnlockTowerData = cOnFigS.GetTowerData(pArSeR.Uint64("unlock_tower_data"))
	if dAtA.UnlockTowerData == nil {
		return errors.Errorf("%s 配置的关联字段[unlock_tower_data] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("unlock_tower_data"), *pArSeR)
	}

	return nil
}

// start with SecretTowerMiscData ----------------------------------

func LoadSecretTowerMiscData(gos *config.GameObjects) (*SecretTowerMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.SecretTowerMiscDataPath
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

	dAtA, err := NewSecretTowerMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedSecretTowerMiscData(gos *config.GameObjects, dAtA *SecretTowerMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SecretTowerMiscDataPath
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

func NewSecretTowerMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SecretTowerMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSecretTowerMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SecretTowerMiscData{}

	dAtA.MaxTimes = pArSeR.Uint64("max_times")
	dAtA.MaxHelpTimes = pArSeR.Uint64("max_help_times")
	dAtA.MaxGuildContribution = 100
	if pArSeR.KeyExist("max_guild_contribution") {
		dAtA.MaxGuildContribution = pArSeR.Uint64("max_guild_contribution")
	}

	dAtA.MaxRecord = 10
	if pArSeR.KeyExist("max_record") {
		dAtA.MaxRecord = pArSeR.Uint64("max_record")
	}

	return dAtA, nil
}

var vAlIdAtOrSecretTowerMiscData = map[string]*config.Validator{

	"max_times":              config.ParseValidator("int>0", "", false, nil, nil),
	"max_help_times":         config.ParseValidator("int>0", "", false, nil, nil),
	"max_guild_contribution": config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"max_record":             config.ParseValidator("int>0", "", false, nil, []string{"10"}),
}

func (dAtA *SecretTowerMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SecretTowerMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SecretTowerMiscData) Encode() *shared_proto.SecretTowerMiscProto {
	out := &shared_proto.SecretTowerMiscProto{}
	out.MaxTimes = config.U64ToI32(dAtA.MaxTimes)
	out.MaxHelpTimes = config.U64ToI32(dAtA.MaxHelpTimes)
	out.MaxRecord = config.U64ToI32(dAtA.MaxRecord)

	return out
}

func ArrayEncodeSecretTowerMiscData(datas []*SecretTowerMiscData) []*shared_proto.SecretTowerMiscProto {

	out := make([]*shared_proto.SecretTowerMiscProto, 0, len(datas))
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

func (dAtA *SecretTowerMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with SecretTowerWordsData ----------------------------------

func LoadSecretTowerWordsData(gos *config.GameObjects) (map[uint64]*SecretTowerWordsData, map[*SecretTowerWordsData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.SecretTowerWordsDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*SecretTowerWordsData, len(lIsT))
	pArSeRmAp := make(map[*SecretTowerWordsData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrSecretTowerWordsData) {
			continue
		}

		dAtA, err := NewSecretTowerWordsData(fIlEnAmE, pArSeR)
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

func SetRelatedSecretTowerWordsData(dAtAmAp map[*SecretTowerWordsData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SecretTowerWordsDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetSecretTowerWordsDataKeyArray(datas []*SecretTowerWordsData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewSecretTowerWordsData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SecretTowerWordsData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSecretTowerWordsData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SecretTowerWordsData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Words = pArSeR.String("words")

	return dAtA, nil
}

var vAlIdAtOrSecretTowerWordsData = map[string]*config.Validator{

	"id":    config.ParseValidator("int>0", "", false, nil, nil),
	"words": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *SecretTowerWordsData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SecretTowerWordsData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SecretTowerWordsData) Encode() *shared_proto.SecretTowerWordsDataProto {
	out := &shared_proto.SecretTowerWordsDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Words = dAtA.Words

	return out
}

func ArrayEncodeSecretTowerWordsData(datas []*SecretTowerWordsData) []*shared_proto.SecretTowerWordsDataProto {

	out := make([]*shared_proto.SecretTowerWordsDataProto, 0, len(datas))
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

func (dAtA *SecretTowerWordsData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with TowerData ----------------------------------

func LoadTowerData(gos *config.GameObjects) (map[uint64]*TowerData, map[*TowerData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TowerDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TowerData, len(lIsT))
	pArSeRmAp := make(map[*TowerData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTowerData) {
			continue
		}

		dAtA, err := NewTowerData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Floor
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Floor], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedTowerData(dAtAmAp map[*TowerData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TowerDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTowerDataKeyArray(datas []*TowerData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Floor)
		}
	}

	return out
}

func NewTowerData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TowerData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTowerData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TowerData{}

	dAtA.Floor = pArSeR.Uint64("floor")
	// releated field: FirstPassPrize
	// skip field: FirstPassPrizeProto
	// skip field: FirstPassPrizeBytes
	// releated field: ShowPrize
	// releated field: Plunder
	// releated field: BoxPrize
	dAtA.CheckPoint = pArSeR.Bool("check_point")
	// releated field: Monster
	dAtA.Desc = pArSeR.String("desc")
	// releated field: CombatScene
	// skip field: UnlockSecretTower

	return dAtA, nil
}

var vAlIdAtOrTowerData = map[string]*config.Validator{

	"floor":            config.ParseValidator("int>0", "", false, nil, nil),
	"first_pass_prize": config.ParseValidator("string", "", false, nil, nil),
	"show_prize":       config.ParseValidator("string", "", false, nil, nil),
	"plunder":          config.ParseValidator("string", "", false, nil, nil),
	"box_prize":        config.ParseValidator("string", "", false, nil, nil),
	"check_point":      config.ParseValidator("bool", "", false, nil, nil),
	"monster":          config.ParseValidator("string", "", false, nil, nil),
	"desc":             config.ParseValidator("string", "", false, nil, nil),
	"combat_scene":     config.ParseValidator("string", "", false, nil, []string{"CombatScene"}),
}

func (dAtA *TowerData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TowerData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TowerData) Encode() *shared_proto.TowerDataProto {
	out := &shared_proto.TowerDataProto{}
	out.Floor = config.U64ToI32(dAtA.Floor)
	if dAtA.FirstPassPrize != nil {
		out.FirstPassPrize = dAtA.FirstPassPrize.Encode()
	}
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.Encode()
	}
	if dAtA.BoxPrize != nil {
		out.BoxPrize = dAtA.BoxPrize.Encode()
	}
	if dAtA.Monster != nil {
		out.Monster = dAtA.Monster.Encode()
	}
	out.Desc = dAtA.Desc
	if dAtA.UnlockSecretTower != nil {
		out.UnlockSecretTowerId = config.U64ToI32(dAtA.UnlockSecretTower.Id)
	}

	return out
}

func ArrayEncodeTowerData(datas []*TowerData) []*shared_proto.TowerDataProto {

	out := make([]*shared_proto.TowerDataProto, 0, len(datas))
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

func (dAtA *TowerData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.FirstPassPrize = cOnFigS.GetPrize(pArSeR.Int("first_pass_prize"))
	if dAtA.FirstPassPrize == nil {
		return errors.Errorf("%s 配置的关联字段[first_pass_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_pass_prize"), *pArSeR)
	}

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil {
		return errors.Errorf("%s 配置的关联字段[show_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
	}

	dAtA.Plunder = cOnFigS.GetPlunder(pArSeR.Uint64("plunder"))
	if dAtA.Plunder == nil {
		return errors.Errorf("%s 配置的关联字段[plunder] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder"), *pArSeR)
	}

	dAtA.BoxPrize = cOnFigS.GetPrize(pArSeR.Int("box_prize"))
	if dAtA.BoxPrize == nil && pArSeR.Int("box_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[box_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("box_prize"), *pArSeR)
	}

	dAtA.Monster = cOnFigS.GetMonsterMasterData(pArSeR.Uint64("monster"))
	if dAtA.Monster == nil {
		return errors.Errorf("%s 配置的关联字段[monster] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("monster"), *pArSeR)
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

type related_configs interface {
	GetCombatScene(string) *scene.CombatScene
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
	GetPlunder(uint64) *resdata.Plunder
	GetPrize(int) *resdata.Prize
	GetTowerData(uint64) *TowerData
}
