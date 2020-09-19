// AUTO_GEN, DONT MODIFY!!!
package captain

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/spell"
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

// start with CaptainAbilityData ----------------------------------

func LoadCaptainAbilityData(gos *config.GameObjects) (map[uint64]*CaptainAbilityData, map[*CaptainAbilityData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CaptainAbilityDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CaptainAbilityData, len(lIsT))
	pArSeRmAp := make(map[*CaptainAbilityData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCaptainAbilityData) {
			continue
		}

		dAtA, err := NewCaptainAbilityData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Ability
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Ability], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedCaptainAbilityData(dAtAmAp map[*CaptainAbilityData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CaptainAbilityDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCaptainAbilityDataKeyArray(datas []*CaptainAbilityData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Ability)
		}
	}

	return out
}

func NewCaptainAbilityData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CaptainAbilityData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCaptainAbilityData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CaptainAbilityData{}

	dAtA.Ability = pArSeR.Uint64("ability")
	dAtA.UpgradeExp = pArSeR.Uint64("upgrade_exp")
	// releated field: SellPrice
	// releated field: FirePrice
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Quality = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("quality"))])
	if i, err := strconv.ParseInt(pArSeR.String("quality"), 10, 32); err == nil {
		dAtA.Quality = shared_proto.Quality(i)
	}

	dAtA.Title = pArSeR.String("title")
	// skip field: MaxLevel
	dAtA.UnlockSpellCount = pArSeR.Uint64("unlock_spell_count")

	return dAtA, nil
}

var vAlIdAtOrCaptainAbilityData = map[string]*config.Validator{

	"ability":            config.ParseValidator("int>0", "", false, nil, nil),
	"upgrade_exp":        config.ParseValidator("int>0", "", false, nil, nil),
	"sell_price":         config.ParseValidator("string", "", false, nil, nil),
	"fire_price":         config.ParseValidator("string", "", false, nil, nil),
	"desc":               config.ParseValidator("string", "", false, nil, nil),
	"quality":            config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), nil),
	"title":              config.ParseValidator("string", "", false, nil, nil),
	"unlock_spell_count": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *CaptainAbilityData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CaptainAbilityData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CaptainAbilityData) Encode() *shared_proto.CaptainAbilityProto {
	out := &shared_proto.CaptainAbilityProto{}
	out.Ability = config.U64ToI32(dAtA.Ability)
	out.UpgradeExp = config.U64ToI32(dAtA.UpgradeExp)
	if dAtA.SellPrice != nil {
		out.SellPrice = dAtA.SellPrice.Encode()
	}
	if dAtA.FirePrice != nil {
		out.FirePrice = dAtA.FirePrice.Encode()
	}
	out.Desc = dAtA.Desc
	out.Quality = dAtA.Quality
	out.UnlockSpellCount = config.U64ToI32(dAtA.UnlockSpellCount)

	return out
}

func ArrayEncodeCaptainAbilityData(datas []*CaptainAbilityData) []*shared_proto.CaptainAbilityProto {

	out := make([]*shared_proto.CaptainAbilityProto, 0, len(datas))
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

func (dAtA *CaptainAbilityData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.SellPrice = cOnFigS.GetPrize(pArSeR.Int("sell_price"))
	if dAtA.SellPrice == nil {
		return errors.Errorf("%s 配置的关联字段[sell_price] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("sell_price"), *pArSeR)
	}

	dAtA.FirePrice = cOnFigS.GetPrize(pArSeR.Int("fire_price"))
	if dAtA.FirePrice == nil {
		return errors.Errorf("%s 配置的关联字段[fire_price] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fire_price"), *pArSeR)
	}

	return nil
}

// start with CaptainData ----------------------------------

func LoadCaptainData(gos *config.GameObjects) (map[uint64]*CaptainData, map[*CaptainData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CaptainDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CaptainData, len(lIsT))
	pArSeRmAp := make(map[*CaptainData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCaptainData) {
			continue
		}

		dAtA, err := NewCaptainData(fIlEnAmE, pArSeR)
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

func SetRelatedCaptainData(dAtAmAp map[*CaptainData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CaptainDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCaptainDataKeyArray(datas []*CaptainData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCaptainData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CaptainData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCaptainData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CaptainData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Rarity
	dAtA.Name = pArSeR.String("name")
	// releated field: Icon
	dAtA.Spine = ""
	if pArSeR.KeyExist("spine") {
		dAtA.Spine = pArSeR.String("spine")
	}

	dAtA.Desc = pArSeR.String("desc")
	// releated field: Race
	// releated field: PrizeIfHas
	// releated field: ResObject
	dAtA.ObtainWays = pArSeR.Uint64Array("obtain_ways", "", false)
	dAtA.FishingObtain = pArSeR.Bool("fishing_obtain")
	dAtA.Sound = ""
	if pArSeR.KeyExist("sound") {
		dAtA.Sound = pArSeR.String("sound")
	}

	// releated field: BaseSpell
	// skip field: Star
	dAtA.InitRage = 0
	if pArSeR.KeyExist("init_rage") {
		dAtA.InitRage = pArSeR.Uint64("init_rage")
	}

	// skip field: GiftData

	return dAtA, nil
}

var vAlIdAtOrCaptainData = map[string]*config.Validator{

	"id":             config.ParseValidator("int>0", "", false, nil, nil),
	"rarity":         config.ParseValidator("string", "", false, nil, nil),
	"name":           config.ParseValidator("string", "", false, nil, nil),
	"icon":           config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"spine":          config.ParseValidator("string", "", false, nil, []string{""}),
	"desc":           config.ParseValidator("string", "", false, nil, nil),
	"race":           config.ParseValidator("string", "", false, nil, nil),
	"prize_if_has":   config.ParseValidator("string", "", false, nil, nil),
	"obtain_ways":    config.ParseValidator("int", "", true, nil, nil),
	"fishing_obtain": config.ParseValidator("bool", "", false, nil, nil),
	"sound":          config.ParseValidator("string", "", false, nil, []string{""}),
	"base_spell":     config.ParseValidator("string", "", false, nil, nil),
	"init_rage":      config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *CaptainData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CaptainData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CaptainData) Encode() *shared_proto.CaptainDataProto {
	out := &shared_proto.CaptainDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.Rarity != nil {
		out.Rarity = config.U64ToI32(dAtA.Rarity.Id)
	}
	out.Name = dAtA.Name
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.Spine = dAtA.Spine
	out.Desc = dAtA.Desc
	if dAtA.Race != nil {
		out.Race = dAtA.Race.Race
	}
	if dAtA.PrizeIfHas != nil {
		out.PrizeIfHas = dAtA.PrizeIfHas.Encode()
	}
	out.ObtainWays = config.U64a2I32a(dAtA.ObtainWays)
	out.FishingObtain = dAtA.FishingObtain
	out.Sound = dAtA.Sound
	if dAtA.BaseSpell != nil {
		out.BaseSpell = config.U64ToI32(dAtA.BaseSpell.Id)
	}
	if dAtA.Star != nil {
		out.Star = ArrayEncodeCaptainStarData(dAtA.Star)
	}
	out.InitRage = config.U64ToI32(dAtA.InitRage)

	return out
}

func ArrayEncodeCaptainData(datas []*CaptainData) []*shared_proto.CaptainDataProto {

	out := make([]*shared_proto.CaptainDataProto, 0, len(datas))
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

func (dAtA *CaptainData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Rarity = cOnFigS.GetCaptainRarityData(pArSeR.Uint64("rarity"))
	if dAtA.Rarity == nil {
		return errors.Errorf("%s 配置的关联字段[rarity] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("rarity"), *pArSeR)
	}

	if pArSeR.KeyExist("icon") {
		dAtA.Icon = cOnFigS.GetIcon(pArSeR.String("icon"))
	} else {
		dAtA.Icon = cOnFigS.GetIcon("Icon")
	}
	if dAtA.Icon == nil {
		return errors.Errorf("%s 配置的关联字段[icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("icon"), *pArSeR)
	}

	dAtA.Race = cOnFigS.GetRaceData(pArSeR.Int("race"))
	if dAtA.Race == nil {
		return errors.Errorf("%s 配置的关联字段[race] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("race"), *pArSeR)
	}

	dAtA.PrizeIfHas = cOnFigS.GetPrize(pArSeR.Int("prize_if_has"))
	if dAtA.PrizeIfHas == nil {
		return errors.Errorf("%s 配置的关联字段[prize_if_has] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize_if_has"), *pArSeR)
	}

	dAtA.ResObject = cOnFigS.GetResCaptainData(pArSeR.Uint64("id"))
	if dAtA.ResObject == nil {
		return errors.Errorf("%s 配置的关联字段[id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("id"), *pArSeR)
	}

	dAtA.BaseSpell = cOnFigS.GetSpellFacadeData(pArSeR.Uint64("base_spell"))
	if dAtA.BaseSpell == nil {
		return errors.Errorf("%s 配置的关联字段[base_spell] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("base_spell"), *pArSeR)
	}

	return nil
}

// start with CaptainFriendshipData ----------------------------------

func LoadCaptainFriendshipData(gos *config.GameObjects) (map[uint64]*CaptainFriendshipData, map[*CaptainFriendshipData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CaptainFriendshipDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CaptainFriendshipData, len(lIsT))
	pArSeRmAp := make(map[*CaptainFriendshipData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCaptainFriendshipData) {
			continue
		}

		dAtA, err := NewCaptainFriendshipData(fIlEnAmE, pArSeR)
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

func SetRelatedCaptainFriendshipData(dAtAmAp map[*CaptainFriendshipData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CaptainFriendshipDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCaptainFriendshipDataKeyArray(datas []*CaptainFriendshipData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCaptainFriendshipData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CaptainFriendshipData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCaptainFriendshipData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CaptainFriendshipData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Tips = pArSeR.String("tips")
	// releated field: Captains
	dAtA.MoveSpeedRate = pArSeR.Float64("move_speed_rate")
	// releated field: AllStat
	for _, v := range pArSeR.StringArray("race", "", false) {
		x := shared_proto.Race(shared_proto.Race_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.Race(i)
		}
		dAtA.Race = append(dAtA.Race, x)
	}

	// releated field: RaceStat
	dAtA.EffectDesc = pArSeR.StringArray("effect_desc", "", false)
	dAtA.EffectAmount = pArSeR.Uint64Array("effect_amount", "", false)

	return dAtA, nil
}

var vAlIdAtOrCaptainFriendshipData = map[string]*config.Validator{

	"id":              config.ParseValidator("int>0", "", false, nil, nil),
	"name":            config.ParseValidator("string", "", false, nil, nil),
	"desc":            config.ParseValidator("string", "", false, nil, nil),
	"tips":            config.ParseValidator("string", "", false, nil, nil),
	"captains":        config.ParseValidator("string", "", true, nil, nil),
	"move_speed_rate": config.ParseValidator("float64", "", false, nil, nil),
	"all_stat":        config.ParseValidator("string", "", false, nil, nil),
	"race":            config.ParseValidator("string", "", true, config.EnumMapKeys(shared_proto.Race_value, 0), nil),
	"race_stat":       config.ParseValidator("uint,duplicate", "", true, nil, nil),
	"effect_desc":     config.ParseValidator("string", "", true, nil, nil),
	"effect_amount":   config.ParseValidator("uint,duplicate", "", true, nil, nil),
}

func (dAtA *CaptainFriendshipData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CaptainFriendshipData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CaptainFriendshipData) Encode() *shared_proto.CaptainFriendshipDataProto {
	out := &shared_proto.CaptainFriendshipDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.Tips = dAtA.Tips
	if dAtA.Captains != nil {
		out.Captains = config.U64a2I32a(GetCaptainDataKeyArray(dAtA.Captains))
	}
	out.EffectDesc = dAtA.EffectDesc
	out.EffectAmount = config.U64a2I32a(dAtA.EffectAmount)

	return out
}

func ArrayEncodeCaptainFriendshipData(datas []*CaptainFriendshipData) []*shared_proto.CaptainFriendshipDataProto {

	out := make([]*shared_proto.CaptainFriendshipDataProto, 0, len(datas))
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

func (dAtA *CaptainFriendshipData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("captains", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetCaptainData(v)
		if obj != nil {
			dAtA.Captains = append(dAtA.Captains, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[captains] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captains"), *pArSeR)
		}
	}

	dAtA.AllStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("all_stat"))
	if dAtA.AllStat == nil && pArSeR.Uint64("all_stat") != 0 {
		return errors.Errorf("%s 配置的关联字段[all_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("all_stat"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("race_stat", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetSpriteStat(v)
		if obj != nil {
			dAtA.RaceStat = append(dAtA.RaceStat, obj)
		} else if v != 0 {
			return errors.Errorf("%s 配置的关联字段[race_stat] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("race_stat"), *pArSeR)
		}
	}

	return nil
}

// start with CaptainLevelData ----------------------------------

func LoadCaptainLevelData(gos *config.GameObjects) (map[uint64]*CaptainLevelData, map[*CaptainLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CaptainLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CaptainLevelData, len(lIsT))
	pArSeRmAp := make(map[*CaptainLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCaptainLevelData) {
			continue
		}

		dAtA, err := NewCaptainLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedCaptainLevelData(dAtAmAp map[*CaptainLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CaptainLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCaptainLevelDataKeyArray(datas []*CaptainLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCaptainLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CaptainLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCaptainLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CaptainLevelData{}

	dAtA.Rebirth = pArSeR.Uint64("rebirth")
	dAtA.Level = pArSeR.Uint64("level")
	dAtA.SoldierCapcity = pArSeR.Uint64("soldier_capcity")
	dAtA.AbilityLimit = pArSeR.Uint64("ability_limit")
	dAtA.UpgradeExp = pArSeR.Uint64("upgrade_exp")
	dAtA.GemSlotCount = pArSeR.Uint64("gem_slot_count")
	// skip field: HasNewGemSlot

	// calculate fields
	dAtA.Id = CaptainLevelId(dAtA.Rebirth, dAtA.Level)

	return dAtA, nil
}

var vAlIdAtOrCaptainLevelData = map[string]*config.Validator{

	"rebirth":         config.ParseValidator("uint", "", false, nil, nil),
	"level":           config.ParseValidator("int>0", "", false, nil, nil),
	"soldier_capcity": config.ParseValidator("int>0", "", false, nil, nil),
	"ability_limit":   config.ParseValidator("uint", "", false, nil, nil),
	"upgrade_exp":     config.ParseValidator("int>0", "", false, nil, nil),
	"gem_slot_count":  config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *CaptainLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CaptainLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CaptainLevelData) Encode() *shared_proto.CaptainLevelProto {
	out := &shared_proto.CaptainLevelProto{}
	out.Rebirth = config.U64ToI32(dAtA.Rebirth)
	out.Level = config.U64ToI32(dAtA.Level)
	out.AbilityLimit = config.U64ToI32(dAtA.AbilityLimit)
	out.UpgradeExp = config.U64ToI32(dAtA.UpgradeExp)
	out.GemSlotCount = config.U64ToI32(dAtA.GemSlotCount)
	out.HasNewGemSlot = dAtA.HasNewGemSlot

	return out
}

func ArrayEncodeCaptainLevelData(datas []*CaptainLevelData) []*shared_proto.CaptainLevelProto {

	out := make([]*shared_proto.CaptainLevelProto, 0, len(datas))
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

func (dAtA *CaptainLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with CaptainOfficialCountData ----------------------------------

func LoadCaptainOfficialCountData(gos *config.GameObjects) (map[uint64]*CaptainOfficialCountData, map[*CaptainOfficialCountData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CaptainOfficialCountDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CaptainOfficialCountData, len(lIsT))
	pArSeRmAp := make(map[*CaptainOfficialCountData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCaptainOfficialCountData) {
			continue
		}

		dAtA, err := NewCaptainOfficialCountData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.HeroLevel
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[HeroLevel], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedCaptainOfficialCountData(dAtAmAp map[*CaptainOfficialCountData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CaptainOfficialCountDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCaptainOfficialCountDataKeyArray(datas []*CaptainOfficialCountData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.HeroLevel)
		}
	}

	return out
}

func NewCaptainOfficialCountData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CaptainOfficialCountData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCaptainOfficialCountData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CaptainOfficialCountData{}

	dAtA.HeroLevel = pArSeR.Uint64("hero_level")
	dAtA.OfficialId = pArSeR.Uint64Array("official_id", "", false)
	dAtA.MaxCount = pArSeR.Uint64Array("max_count", "", false)

	return dAtA, nil
}

var vAlIdAtOrCaptainOfficialCountData = map[string]*config.Validator{

	"hero_level":  config.ParseValidator("int>0", "", false, nil, nil),
	"official_id": config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
	"max_count":   config.ParseValidator("uint,duplicate,notAllNil", "", true, nil, nil),
}

func (dAtA *CaptainOfficialCountData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with CaptainOfficialData ----------------------------------

func LoadCaptainOfficialData(gos *config.GameObjects) (map[uint64]*CaptainOfficialData, map[*CaptainOfficialData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CaptainOfficialDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CaptainOfficialData, len(lIsT))
	pArSeRmAp := make(map[*CaptainOfficialData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCaptainOfficialData) {
			continue
		}

		dAtA, err := NewCaptainOfficialData(fIlEnAmE, pArSeR)
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

func SetRelatedCaptainOfficialData(dAtAmAp map[*CaptainOfficialData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CaptainOfficialDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCaptainOfficialDataKeyArray(datas []*CaptainOfficialData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCaptainOfficialData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CaptainOfficialData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCaptainOfficialData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CaptainOfficialData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.OfficialName = pArSeR.String("official_name")
	// releated field: SpriteStat
	dAtA.NeedGongxun = pArSeR.Uint64("need_gongxun")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Desc = pArSeR.String("desc")

	return dAtA, nil
}

var vAlIdAtOrCaptainOfficialData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"official_name": config.ParseValidator("string>0", "", false, nil, nil),
	"sprite_stat":   config.ParseValidator("string", "", false, nil, nil),
	"need_gongxun":  config.ParseValidator("int>0", "", false, nil, nil),
	"icon":          config.ParseValidator("string", "", false, nil, nil),
	"desc":          config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *CaptainOfficialData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CaptainOfficialData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CaptainOfficialData) Encode() *shared_proto.CaptainOfficialProto {
	out := &shared_proto.CaptainOfficialProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.OfficialName = dAtA.OfficialName
	if dAtA.SpriteStat != nil {
		out.SpriteStat = dAtA.SpriteStat.Encode()
	}
	out.NeedGongxun = config.U64ToI32(dAtA.NeedGongxun)
	out.Icon = dAtA.Icon
	out.Desc = dAtA.Desc

	return out
}

func ArrayEncodeCaptainOfficialData(datas []*CaptainOfficialData) []*shared_proto.CaptainOfficialProto {

	out := make([]*shared_proto.CaptainOfficialProto, 0, len(datas))
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

func (dAtA *CaptainOfficialData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.SpriteStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("sprite_stat"))
	if dAtA.SpriteStat == nil {
		return errors.Errorf("%s 配置的关联字段[sprite_stat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("sprite_stat"), *pArSeR)
	}

	return nil
}

// start with CaptainRarityData ----------------------------------

func LoadCaptainRarityData(gos *config.GameObjects) (map[uint64]*CaptainRarityData, map[*CaptainRarityData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CaptainRarityDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CaptainRarityData, len(lIsT))
	pArSeRmAp := make(map[*CaptainRarityData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCaptainRarityData) {
			continue
		}

		dAtA, err := NewCaptainRarityData(fIlEnAmE, pArSeR)
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

func SetRelatedCaptainRarityData(dAtAmAp map[*CaptainRarityData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CaptainRarityDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCaptainRarityDataKeyArray(datas []*CaptainRarityData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCaptainRarityData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CaptainRarityData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCaptainRarityData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CaptainRarityData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Color = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("color"))])
	if i, err := strconv.ParseInt(pArSeR.String("color"), 10, 32); err == nil {
		dAtA.Color = shared_proto.Quality(i)
	}

	dAtA.Coef = pArSeR.Float64("coef")

	return dAtA, nil
}

var vAlIdAtOrCaptainRarityData = map[string]*config.Validator{

	"id":    config.ParseValidator("int>0", "", false, nil, nil),
	"name":  config.ParseValidator("string", "", false, nil, nil),
	"color": config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), []string{"3"}),
	"coef":  config.ParseValidator("float64>0", "", false, nil, nil),
}

func (dAtA *CaptainRarityData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CaptainRarityData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CaptainRarityData) Encode() *shared_proto.CaptainRarityDataProto {
	out := &shared_proto.CaptainRarityDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Coef = config.F64ToI32X1000(dAtA.Coef)

	return out
}

func ArrayEncodeCaptainRarityData(datas []*CaptainRarityData) []*shared_proto.CaptainRarityDataProto {

	out := make([]*shared_proto.CaptainRarityDataProto, 0, len(datas))
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

func (dAtA *CaptainRarityData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with CaptainRebirthLevelData ----------------------------------

func LoadCaptainRebirthLevelData(gos *config.GameObjects) (map[uint64]*CaptainRebirthLevelData, map[*CaptainRebirthLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CaptainRebirthLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CaptainRebirthLevelData, len(lIsT))
	pArSeRmAp := make(map[*CaptainRebirthLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCaptainRebirthLevelData) {
			continue
		}

		dAtA, err := NewCaptainRebirthLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedCaptainRebirthLevelData(dAtAmAp map[*CaptainRebirthLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CaptainRebirthLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCaptainRebirthLevelDataKeyArray(datas []*CaptainRebirthLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewCaptainRebirthLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CaptainRebirthLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCaptainRebirthLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CaptainRebirthLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Cd, err = config.ParseDuration(pArSeR.String("cd"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[cd] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("cd"), dAtA)
	}

	dAtA.RequiredCaptainLevel = pArSeR.Uint64("required_captain_level")
	// skip field: CaptainLevelLimit
	// skip field: FirstCaptainLevel
	dAtA.SpriteStatPoint = pArSeR.Uint64("sprite_stat_point")
	dAtA.SoldierCapcity = pArSeR.Uint64("soldier_capcity")
	dAtA.AbilityLimit = pArSeR.Uint64("ability_limit")
	dAtA.AbilityExp = pArSeR.Uint64("ability_exp")
	dAtA.BeforeRebirthLevel = 0
	if pArSeR.KeyExist("before_rebirth_level") {
		dAtA.BeforeRebirthLevel = pArSeR.Uint64("before_rebirth_level")
	}

	dAtA.HeroLevelLimit = 0
	if pArSeR.KeyExist("hero_level_limit") {
		dAtA.HeroLevelLimit = pArSeR.Uint64("hero_level_limit")
	}

	return dAtA, nil
}

var vAlIdAtOrCaptainRebirthLevelData = map[string]*config.Validator{

	"level": config.ParseValidator("uint", "", false, nil, nil),
	"cd":    config.ParseValidator("string", "", false, nil, nil),
	"required_captain_level": config.ParseValidator("int>0", "", false, nil, nil),
	"sprite_stat_point":      config.ParseValidator("uint", "", false, nil, nil),
	"soldier_capcity":        config.ParseValidator("uint", "", false, nil, nil),
	"ability_limit":          config.ParseValidator("int>0", "", false, nil, nil),
	"ability_exp":            config.ParseValidator("int>0", "", false, nil, nil),
	"before_rebirth_level":   config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"hero_level_limit":       config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *CaptainRebirthLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CaptainRebirthLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CaptainRebirthLevelData) Encode() *shared_proto.CaptainRebirthLevelProto {
	out := &shared_proto.CaptainRebirthLevelProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Cd = config.Duration2I32Seconds(dAtA.Cd)
	if dAtA.CaptainLevelLimit != nil {
		out.CaptainLevelLimit = config.U64ToI32(dAtA.CaptainLevelLimit.Level)
	}
	out.SpriteStatPoint = config.U64ToI32(dAtA.SpriteStatPoint)
	out.SoldierCapcity = config.U64ToI32(dAtA.SoldierCapcity)
	out.AbilityLimit = config.U64ToI32(dAtA.AbilityLimit)
	out.AbilityExp = config.U64ToI32(dAtA.AbilityExp)
	out.BeforeRebirthLevel = config.U64ToI32(dAtA.BeforeRebirthLevel)
	out.HeroLevelLimit = config.U64ToI32(dAtA.HeroLevelLimit)

	return out
}

func ArrayEncodeCaptainRebirthLevelData(datas []*CaptainRebirthLevelData) []*shared_proto.CaptainRebirthLevelProto {

	out := make([]*shared_proto.CaptainRebirthLevelProto, 0, len(datas))
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

func (dAtA *CaptainRebirthLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with CaptainStarData ----------------------------------

func LoadCaptainStarData(gos *config.GameObjects) (map[uint64]*CaptainStarData, map[*CaptainStarData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CaptainStarDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CaptainStarData, len(lIsT))
	pArSeRmAp := make(map[*CaptainStarData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCaptainStarData) {
			continue
		}

		dAtA, err := NewCaptainStarData(fIlEnAmE, pArSeR)
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

func SetRelatedCaptainStarData(dAtAmAp map[*CaptainStarData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CaptainStarDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCaptainStarDataKeyArray(datas []*CaptainStarData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCaptainStarData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CaptainStarData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCaptainStarData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CaptainStarData{}

	dAtA.CaptainId = pArSeR.Uint64("captain_id")
	dAtA.Star = pArSeR.Uint64("star")
	dAtA.Coef = pArSeR.Float64("coef")
	// releated field: SpriteStat
	// skip field: AddedStat
	// releated field: Cost
	// releated field: Spell

	// calculate fields
	dAtA.Id = CalculateCaptainStarId(dAtA.CaptainId, dAtA.Star)

	return dAtA, nil
}

var vAlIdAtOrCaptainStarData = map[string]*config.Validator{

	"captain_id":  config.ParseValidator("int>0", "", false, nil, nil),
	"star":        config.ParseValidator("int>0", "", false, nil, nil),
	"coef":        config.ParseValidator("float64>0", "", false, nil, nil),
	"sprite_stat": config.ParseValidator("string", "", false, nil, nil),
	"cost":        config.ParseValidator("string", "", false, nil, nil),
	"spell":       config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *CaptainStarData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CaptainStarData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CaptainStarData) Encode() *shared_proto.CaptainStarDataProto {
	out := &shared_proto.CaptainStarDataProto{}
	out.Star = config.U64ToI32(dAtA.Star)
	out.Coef = config.F64ToI32X1000(dAtA.Coef)
	if dAtA.AddedStat != nil {
		out.AddedStat = dAtA.AddedStat.Encode()
	}
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	if dAtA.Spell != nil {
		out.Spell = config.U64a2I32a(spell.GetSpellFacadeDataKeyArray(dAtA.Spell))
	}

	return out
}

func ArrayEncodeCaptainStarData(datas []*CaptainStarData) []*shared_proto.CaptainStarDataProto {

	out := make([]*shared_proto.CaptainStarDataProto, 0, len(datas))
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

func (dAtA *CaptainStarData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.SpriteStat = cOnFigS.GetSpriteStat(pArSeR.Uint64("sprite_stat"))
	if dAtA.SpriteStat == nil {
		return errors.Errorf("%s 配置的关联字段[sprite_stat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("sprite_stat"), *pArSeR)
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("spell", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetSpellFacadeData(v)
		if obj != nil {
			dAtA.Spell = append(dAtA.Spell, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[spell] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("spell"), *pArSeR)
		}
	}

	return nil
}

// start with NamelessCaptainData ----------------------------------

func LoadNamelessCaptainData(gos *config.GameObjects) (map[uint64]*NamelessCaptainData, map[*NamelessCaptainData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.NamelessCaptainDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*NamelessCaptainData, len(lIsT))
	pArSeRmAp := make(map[*NamelessCaptainData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrNamelessCaptainData) {
			continue
		}

		dAtA, err := NewNamelessCaptainData(fIlEnAmE, pArSeR)
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

func SetRelatedNamelessCaptainData(dAtAmAp map[*NamelessCaptainData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.NamelessCaptainDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetNamelessCaptainDataKeyArray(datas []*NamelessCaptainData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewNamelessCaptainData(fIlEnAmE string, pArSeR *config.ObjectParser) (*NamelessCaptainData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrNamelessCaptainData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &NamelessCaptainData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// releated field: Icon
	dAtA.Spine = ""
	if pArSeR.KeyExist("spine") {
		dAtA.Spine = pArSeR.String("spine")
	}

	dAtA.Male = pArSeR.Bool("male")
	dAtA.Race = shared_proto.Race(shared_proto.Race_value[strings.ToUpper(pArSeR.String("race"))])
	if i, err := strconv.ParseInt(pArSeR.String("race"), 10, 32); err == nil {
		dAtA.Race = shared_proto.Race(i)
	}

	// releated field: BaseSpell
	// releated field: Spell
	dAtA.InitRage = 0
	if pArSeR.KeyExist("init_rage") {
		dAtA.InitRage = pArSeR.Uint64("init_rage")
	}

	return dAtA, nil
}

var vAlIdAtOrNamelessCaptainData = map[string]*config.Validator{

	"id":         config.ParseValidator("int>0", "", false, nil, nil),
	"name":       config.ParseValidator("string", "", false, nil, nil),
	"icon":       config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"spine":      config.ParseValidator("string", "", false, nil, []string{""}),
	"male":       config.ParseValidator("bool", "", false, nil, nil),
	"race":       config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Race_value, 0), nil),
	"base_spell": config.ParseValidator("string", "", false, nil, nil),
	"spell":      config.ParseValidator("string", "", true, nil, nil),
	"init_rage":  config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *NamelessCaptainData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *NamelessCaptainData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *NamelessCaptainData) Encode() *shared_proto.NamelessCaptainDataProto {
	out := &shared_proto.NamelessCaptainDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	out.Spine = dAtA.Spine
	out.Male = dAtA.Male
	out.Race = dAtA.Race
	if dAtA.BaseSpell != nil {
		out.BaseSpell = config.U64ToI32(dAtA.BaseSpell.Id)
	}
	if dAtA.Spell != nil {
		out.Spell = config.U64a2I32a(spell.GetSpellFacadeDataKeyArray(dAtA.Spell))
	}
	out.InitRage = config.U64ToI32(dAtA.InitRage)

	return out
}

func ArrayEncodeNamelessCaptainData(datas []*NamelessCaptainData) []*shared_proto.NamelessCaptainDataProto {

	out := make([]*shared_proto.NamelessCaptainDataProto, 0, len(datas))
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

func (dAtA *NamelessCaptainData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.BaseSpell = cOnFigS.GetSpellFacadeData(pArSeR.Uint64("base_spell"))
	if dAtA.BaseSpell == nil {
		return errors.Errorf("%s 配置的关联字段[base_spell] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("base_spell"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("spell", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetSpellFacadeData(v)
		if obj != nil {
			dAtA.Spell = append(dAtA.Spell, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[spell] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("spell"), *pArSeR)
		}
	}

	return nil
}

type related_configs interface {
	GetCaptainData(uint64) *CaptainData
	GetCaptainRarityData(uint64) *CaptainRarityData
	GetCost(int) *resdata.Cost
	GetIcon(string) *icon.Icon
	GetPrize(int) *resdata.Prize
	GetRaceData(int) *race.RaceData
	GetResCaptainData(uint64) *resdata.ResCaptainData
	GetSpellFacadeData(uint64) *spell.SpellFacadeData
	GetSpriteStat(uint64) *data.SpriteStat
}
