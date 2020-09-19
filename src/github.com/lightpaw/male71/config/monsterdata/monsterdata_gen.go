// AUTO_GEN, DONT MODIFY!!!
package monsterdata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
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

// start with MonsterCaptainData ----------------------------------

func LoadMonsterCaptainData(gos *config.GameObjects) (map[uint64]*MonsterCaptainData, map[*MonsterCaptainData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MonsterCaptainDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MonsterCaptainData, len(lIsT))
	pArSeRmAp := make(map[*MonsterCaptainData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMonsterCaptainData) {
			continue
		}

		dAtA, err := NewMonsterCaptainData(fIlEnAmE, pArSeR)
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

func SetRelatedMonsterCaptainData(dAtAmAp map[*MonsterCaptainData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MonsterCaptainDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMonsterCaptainDataKeyArray(datas []*MonsterCaptainData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMonsterCaptainData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MonsterCaptainData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMonsterCaptainData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MonsterCaptainData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Captain
	dAtA.Star = pArSeR.Uint64("star")
	// releated field: NamelessCaptain
	// skip field: CaptainId
	// skip field: IsNameless
	dAtA.UnlockSpellCount = pArSeR.Uint64("unlock_spell_count")
	dAtA.Quality = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("quality"))])
	if i, err := strconv.ParseInt(pArSeR.String("quality"), 10, 32); err == nil {
		dAtA.Quality = shared_proto.Quality(i)
	}

	dAtA.Soldier = pArSeR.Uint64("soldier")
	// skip field: FightAmount
	dAtA.Morale = pArSeR.Uint64("morale")
	dAtA.Level = pArSeR.Uint64("level")
	dAtA.SoldierLevel = pArSeR.Uint64("soldier_level")
	dAtA.RebirthLevel = 0
	if pArSeR.KeyExist("rebirth_level") {
		dAtA.RebirthLevel = pArSeR.Uint64("rebirth_level")
	}

	// releated field: TotalStat
	dAtA.Label = pArSeR.Uint64("label")
	dAtA.Model = 1
	if pArSeR.KeyExist("model") {
		dAtA.Model = pArSeR.Uint64("model")
	}

	dAtA.CanTriggerRestraintSpell = false
	if pArSeR.KeyExist("can_trigger_restraint_spell") {
		dAtA.CanTriggerRestraintSpell = pArSeR.Bool("can_trigger_restraint_spell")
	}

	dAtA.Index = 0
	if pArSeR.KeyExist("index") {
		dAtA.Index = pArSeR.Uint64("index")
	}

	dAtA.XIndex = 0
	if pArSeR.KeyExist("xindex") {
		dAtA.XIndex = pArSeR.Uint64("xindex")
	}

	dAtA.YuanJun = false
	if pArSeR.KeyExist("yuan_jun") {
		dAtA.YuanJun = pArSeR.Bool("yuan_jun")
	}

	return dAtA, nil
}

var vAlIdAtOrMonsterCaptainData = map[string]*config.Validator{

	"id":                 config.ParseValidator("int>0", "", false, nil, nil),
	"captain":            config.ParseValidator("string", "", false, nil, nil),
	"star":               config.ParseValidator("uint", "", false, nil, nil),
	"nameless_captain":   config.ParseValidator("string", "", false, nil, nil),
	"unlock_spell_count": config.ParseValidator("uint", "", false, nil, nil),
	"quality":            config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), nil),
	"soldier":            config.ParseValidator("int>0", "", false, nil, nil),
	"morale":             config.ParseValidator("int>0", "", false, nil, nil),
	"level":              config.ParseValidator("int>0", "", false, nil, nil),
	"soldier_level":      config.ParseValidator("int>0", "", false, nil, nil),
	"rebirth_level":      config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"total_stat":         config.ParseValidator("string", "", false, nil, nil),
	"label":              config.ParseValidator("uint", "", false, nil, nil),
	"model":              config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"can_trigger_restraint_spell": config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"index":    config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"xindex":   config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"yuan_jun": config.ParseValidator("bool", "", false, nil, []string{"false"}),
}

func (dAtA *MonsterCaptainData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MonsterCaptainData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MonsterCaptainData) Encode() *shared_proto.MonsterCaptainDataProto {
	out := &shared_proto.MonsterCaptainDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Star = config.U64ToI32(dAtA.Star)
	out.CaptainId = config.U64ToI32(dAtA.CaptainId)
	out.IsNameless = dAtA.IsNameless
	out.UnlockSpellCount = config.U64ToI32(dAtA.UnlockSpellCount)
	out.Quality = dAtA.Quality
	out.Soldier = config.U64ToI32(dAtA.Soldier)
	out.FightAmount = config.U64ToI32(dAtA.FightAmount)
	out.Level = config.U64ToI32(dAtA.Level)
	out.SoldierLevel = config.U64ToI32(dAtA.SoldierLevel)
	out.RebirthLevel = config.U64ToI32(dAtA.RebirthLevel)
	if dAtA.TotalStat != nil {
		out.TotalStat = dAtA.TotalStat.Encode()
	}
	out.Label = config.U64ToI32(dAtA.Label)
	out.Index = config.U64ToI32(dAtA.Index)
	out.XIndex = config.U64ToI32(dAtA.XIndex)

	return out
}

func ArrayEncodeMonsterCaptainData(datas []*MonsterCaptainData) []*shared_proto.MonsterCaptainDataProto {

	out := make([]*shared_proto.MonsterCaptainDataProto, 0, len(datas))
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

func (dAtA *MonsterCaptainData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Captain = cOnFigS.GetCaptainData(pArSeR.Uint64("captain"))
	if dAtA.Captain == nil && pArSeR.Uint64("captain") != 0 {
		return errors.Errorf("%s 配置的关联字段[captain] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain"), *pArSeR)
	}

	dAtA.NamelessCaptain = cOnFigS.GetNamelessCaptainData(pArSeR.Uint64("nameless_captain"))
	if dAtA.NamelessCaptain == nil && pArSeR.Uint64("nameless_captain") != 0 {
		return errors.Errorf("%s 配置的关联字段[nameless_captain] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("nameless_captain"), *pArSeR)
	}

	dAtA.TotalStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("total_stat"))
	if dAtA.TotalStat == nil {
		return errors.Errorf("%s 配置的关联字段[total_stat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("total_stat"), *pArSeR)
	}

	return nil
}

// start with MonsterMasterData ----------------------------------

func LoadMonsterMasterData(gos *config.GameObjects) (map[uint64]*MonsterMasterData, map[*MonsterMasterData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MonsterMasterDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MonsterMasterData, len(lIsT))
	pArSeRmAp := make(map[*MonsterMasterData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMonsterMasterData) {
			continue
		}

		dAtA, err := NewMonsterMasterData(fIlEnAmE, pArSeR)
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

func SetRelatedMonsterMasterData(dAtAmAp map[*MonsterMasterData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MonsterMasterDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMonsterMasterDataKeyArray(datas []*MonsterMasterData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewMonsterMasterData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MonsterMasterData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMonsterMasterData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MonsterMasterData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// releated field: Icon
	dAtA.Male = pArSeR.Bool("male")
	dAtA.Level = pArSeR.Uint64("level")
	// releated field: Captains
	// skip field: FightAmount
	dAtA.WallLevel = 1
	if pArSeR.KeyExist("wall_level") {
		dAtA.WallLevel = pArSeR.Uint64("wall_level")
	}

	// releated field: WallStat
	dAtA.WallFixDamage = pArSeR.Uint64("wall_fix_damage")
	// releated field: InvadePrize
	// releated field: BeenKillPrize

	return dAtA, nil
}

var vAlIdAtOrMonsterMasterData = map[string]*config.Validator{

	"id":              config.ParseValidator("int>0", "", false, nil, nil),
	"name":            config.ParseValidator("string", "", false, nil, nil),
	"icon":            config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"male":            config.ParseValidator("bool", "", false, nil, nil),
	"level":           config.ParseValidator("int>0", "", false, nil, nil),
	"captains":        config.ParseValidator("string", "", true, nil, nil),
	"wall_level":      config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"wall_stat":       config.ParseValidator("string", "", false, nil, nil),
	"wall_fix_damage": config.ParseValidator("uint", "", false, nil, nil),
	"invade_prize":    config.ParseValidator("string", "", false, nil, nil),
	"been_kill_prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *MonsterMasterData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MonsterMasterData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MonsterMasterData) Encode() *shared_proto.MonsterMasterDataProto {
	out := &shared_proto.MonsterMasterDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.Male = dAtA.Male
	out.Level = config.U64ToI32(dAtA.Level)
	if dAtA.Captains != nil {
		out.Captains = ArrayEncodeMonsterCaptainData(dAtA.Captains)
	}

	return out
}

func ArrayEncodeMonsterMasterData(datas []*MonsterMasterData) []*shared_proto.MonsterMasterDataProto {

	out := make([]*shared_proto.MonsterMasterDataProto, 0, len(datas))
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

func (dAtA *MonsterMasterData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	uint64Keys = pArSeR.Uint64Array("captains", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetMonsterCaptainData(v)
		if obj != nil {
			dAtA.Captains = append(dAtA.Captains, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[captains] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captains"), *pArSeR)
		}
	}

	dAtA.WallStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("wall_stat"))
	if dAtA.WallStat == nil && pArSeR.Uint64("wall_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[wall_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("wall_stat"), *pArSeR)
	}

	dAtA.InvadePrize = cOnFigS.GetPlunderPrize(pArSeR.Uint64("invade_prize"))
	if dAtA.InvadePrize == nil && pArSeR.Uint64("invade_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[invade_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("invade_prize"), *pArSeR)
	}

	dAtA.BeenKillPrize = cOnFigS.GetPlunderPrize(pArSeR.Uint64("been_kill_prize"))
	if dAtA.BeenKillPrize == nil && pArSeR.Uint64("been_kill_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[been_kill_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("been_kill_prize"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetCaptainData(uint64) *captain.CaptainData
	GetIcon(string) *icon.Icon
	GetMonsterCaptainData(uint64) *MonsterCaptainData
	GetNamelessCaptainData(uint64) *captain.NamelessCaptainData
	GetPlunderPrize(uint64) *resdata.PlunderPrize
	GetSpriteStat(uint64) *data.SpriteStat
}
