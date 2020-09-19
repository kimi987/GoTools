// AUTO_GEN, DONT MODIFY!!!
package data

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/i18n"
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

// start with Amount ----------------------------------

func NewAmount(fIlEnAmE string, pArSeR *config.ObjectParser) (*Amount, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrAmount)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &Amount{}

	dAtA.Amount = pArSeR.Uint64("amount")
	dAtA.Percent = pArSeR.Uint64("percent")

	return dAtA, nil
}

var vAlIdAtOrAmount = map[string]*config.Validator{

	"amount":  config.ParseValidator("uint", "", false, nil, nil),
	"percent": config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *Amount) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *Amount) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *Amount) Encode() *shared_proto.AmountProto {
	out := &shared_proto.AmountProto{}
	out.Amount = config.U64ToI32(dAtA.Amount)
	out.Percent = config.U64ToI32(dAtA.Percent)

	return out
}

func ArrayEncodeAmount(datas []*Amount) []*shared_proto.AmountProto {

	out := make([]*shared_proto.AmountProto, 0, len(datas))
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

func (dAtA *Amount) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with BroadcastData ----------------------------------

func LoadBroadcastData(gos *config.GameObjects) (map[string]*BroadcastData, map[*BroadcastData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BroadcastDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*BroadcastData, len(lIsT))
	pArSeRmAp := make(map[*BroadcastData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBroadcastData) {
			continue
		}

		dAtA, err := NewBroadcastData(fIlEnAmE, pArSeR)
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

func SetRelatedBroadcastData(dAtAmAp map[*BroadcastData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BroadcastDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBroadcastDataKeyArray(datas []*BroadcastData) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBroadcastData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BroadcastData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBroadcastData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BroadcastData{}

	dAtA.Id = pArSeR.String("id")
	dAtA.Sequence = pArSeR.Uint64("sequence")
	dAtA.Condition = pArSeR.Uint64Array("condition", "", false)
	dAtA.BcType = shared_proto.BCType(shared_proto.BCType_value[strings.ToUpper(pArSeR.String("bc_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("bc_type"), 10, 32); err == nil {
		dAtA.BcType = shared_proto.BCType(i)
	}

	dAtA.SendHourMinute = []uint64{1200, 1600, 2000}
	if pArSeR.KeyExist("send_hour_minute") {
		dAtA.SendHourMinute = pArSeR.Uint64Array("send_hour_minute", "", false)
	}

	// skip field: SendDuration
	dAtA.SendChat = false
	if pArSeR.KeyExist("send_chat") {
		dAtA.SendChat = pArSeR.Bool("send_chat")
	}

	dAtA.OnlySendOnce = false
	if pArSeR.KeyExist("only_send_once") {
		dAtA.OnlySendOnce = pArSeR.Bool("only_send_once")
	}

	// i18n fields
	dAtA.Text = i18n.NewI18nRef(fIlEnAmE, "text", dAtA.Id, pArSeR.String("text"))

	return dAtA, nil
}

var vAlIdAtOrBroadcastData = map[string]*config.Validator{

	"id":               config.ParseValidator("string", "", false, nil, nil),
	"sequence":         config.ParseValidator("int>0", "", false, nil, nil),
	"text":             config.ParseValidator("string", "", false, nil, nil),
	"condition":        config.ParseValidator("uint", "", true, nil, nil),
	"bc_type":          config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.BCType_value, 0), nil),
	"send_hour_minute": config.ParseValidator("uint", "", true, nil, []string{"1200", "1600", "2000"}),
	"send_chat":        config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"only_send_once":   config.ParseValidator("bool", "", false, nil, []string{"false"}),
}

func (dAtA *BroadcastData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with BroadcastHelp ----------------------------------

func LoadBroadcastHelp(gos *config.GameObjects) (*BroadcastHelp, *config.ObjectParser, error) {
	fIlEnAmE := confpath.BroadcastHelpPath
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

	dAtA, err := NewBroadcastHelp(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedBroadcastHelp(gos *config.GameObjects, dAtA *BroadcastHelp, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BroadcastHelpPath
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

func NewBroadcastHelp(fIlEnAmE string, pArSeR *config.ObjectParser) (*BroadcastHelp, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBroadcastHelp)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BroadcastHelp{}

	// releated field: BaiZhanTouFang
	// releated field: BaoZangDestory
	// releated field: BaseLevel
	// releated field: CaptainAbilityExp
	// releated field: CaptainReBrith
	// releated field: CaptainSoulLevel
	// releated field: CaptainSoulUnlockSpell
	// releated field: CountryChangeKing
	// releated field: CountryChangeName
	// releated field: CountryDestroy
	// releated field: DungeonPrizeCaptainSoul
	// releated field: EquipLevel
	// releated field: FishCaptainSoul
	// releated field: FishEquip
	// releated field: FishGem
	// releated field: GemGet
	// releated field: GuildChangeCountry
	// releated field: GuildChangeFlag
	// releated field: GuildChangeName
	// releated field: GuildCreate
	// releated field: GuildLevel
	// releated field: GuildTanHe
	// releated field: HeroLevel
	// releated field: JiuGuanBaoJi
	// releated field: JiuGuanQingJiao
	// releated field: McWarAtkWin
	// releated field: MixEquip
	// releated field: NpcMonsterSucc
	// releated field: TaozLevel
	// releated field: TaskBwzlPrize
	// releated field: TaskChengJiuLevel
	// releated field: Title
	// releated field: TowerFloor
	// releated field: TowerPrizeEquip
	// releated field: XiongNuSucc
	// releated field: ZhanJiangZhangJieComplete
	// releated field: ZhenBaoGeBaoJi

	return dAtA, nil
}

var vAlIdAtOrBroadcastHelp = map[string]*config.Validator{

	"bai_zhan_tou_fang":             config.ParseValidator("string", "", false, nil, []string{"BaiZhanTouFang"}),
	"bao_zang_destory":              config.ParseValidator("string", "", false, nil, []string{"BaoZangDestory"}),
	"base_level":                    config.ParseValidator("string", "", false, nil, []string{"BaseLevel"}),
	"captain_ability_exp":           config.ParseValidator("string", "", false, nil, []string{"CaptainAbilityExp"}),
	"captain_re_brith":              config.ParseValidator("string", "", false, nil, []string{"CaptainReBrith"}),
	"captain_soul_level":            config.ParseValidator("string", "", false, nil, []string{"CaptainSoulLevel"}),
	"captain_soul_unlock_spell":     config.ParseValidator("string", "", false, nil, []string{"CaptainSoulUnlockSpell"}),
	"country_change_king":           config.ParseValidator("string", "", false, nil, []string{"CountryChangeKing"}),
	"country_change_name":           config.ParseValidator("string", "", false, nil, []string{"CountryChangeName"}),
	"country_destroy":               config.ParseValidator("string", "", false, nil, []string{"CountryDestroy"}),
	"dungeon_prize_captain_soul":    config.ParseValidator("string", "", false, nil, []string{"DungeonPrizeCaptainSoul"}),
	"equip_level":                   config.ParseValidator("string", "", false, nil, []string{"EquipLevel"}),
	"fish_captain_soul":             config.ParseValidator("string", "", false, nil, []string{"FishCaptainSoul"}),
	"fish_equip":                    config.ParseValidator("string", "", false, nil, []string{"FishEquip"}),
	"fish_gem":                      config.ParseValidator("string", "", false, nil, []string{"FishGem"}),
	"gem_get":                       config.ParseValidator("string", "", false, nil, []string{"GemGet"}),
	"guild_change_country":          config.ParseValidator("string", "", false, nil, []string{"GuildChangeCountry"}),
	"guild_change_flag":             config.ParseValidator("string", "", false, nil, []string{"GuildChangeFlag"}),
	"guild_change_name":             config.ParseValidator("string", "", false, nil, []string{"GuildChangeName"}),
	"guild_create":                  config.ParseValidator("string", "", false, nil, []string{"GuildCreate"}),
	"guild_level":                   config.ParseValidator("string", "", false, nil, []string{"GuildLevel"}),
	"guild_tan_he":                  config.ParseValidator("string", "", false, nil, []string{"GuildTanHe"}),
	"hero_level":                    config.ParseValidator("string", "", false, nil, []string{"HeroLevel"}),
	"jiu_guan_bao_ji":               config.ParseValidator("string", "", false, nil, []string{"JiuGuanBaoJi"}),
	"jiu_guan_qing_jiao":            config.ParseValidator("string", "", false, nil, []string{"JiuGuanQingJiao"}),
	"mc_war_atk_win":                config.ParseValidator("string", "", false, nil, []string{"McWarAtkWin"}),
	"mix_equip":                     config.ParseValidator("string", "", false, nil, []string{"MixEquip"}),
	"npc_monster_succ":              config.ParseValidator("string", "", false, nil, []string{"NpcMonsterSucc"}),
	"taoz_level":                    config.ParseValidator("string", "", false, nil, []string{"TaozLevel"}),
	"task_bwzl_prize":               config.ParseValidator("string", "", false, nil, []string{"TaskBwzlPrize"}),
	"task_cheng_jiu_level":          config.ParseValidator("string", "", false, nil, []string{"TaskChengJiuLevel"}),
	"title":                         config.ParseValidator("string", "", false, nil, []string{"Title"}),
	"tower_floor":                   config.ParseValidator("string", "", false, nil, []string{"TowerFloor"}),
	"tower_prize_equip":             config.ParseValidator("string", "", false, nil, []string{"TowerPrizeEquip"}),
	"xiong_nu_succ":                 config.ParseValidator("string", "", false, nil, []string{"XiongNuSucc"}),
	"zhan_jiang_zhang_jie_complete": config.ParseValidator("string", "", false, nil, []string{"ZhanJiangZhangJieComplete"}),
	"zhen_bao_ge_bao_ji":            config.ParseValidator("string", "", false, nil, []string{"ZhenBaoGeBaoJi"}),
}

func (dAtA *BroadcastHelp) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("bai_zhan_tou_fang") {
		dAtA.BaiZhanTouFang = cOnFigS.GetBroadcastData(pArSeR.String("bai_zhan_tou_fang"))
	} else {
		dAtA.BaiZhanTouFang = cOnFigS.GetBroadcastData("BaiZhanTouFang")
	}
	if dAtA.BaiZhanTouFang == nil {
		return errors.Errorf("%s 配置的关联字段[bai_zhan_tou_fang] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("bai_zhan_tou_fang"), *pArSeR)
	}

	if pArSeR.KeyExist("bao_zang_destory") {
		dAtA.BaoZangDestory = cOnFigS.GetBroadcastData(pArSeR.String("bao_zang_destory"))
	} else {
		dAtA.BaoZangDestory = cOnFigS.GetBroadcastData("BaoZangDestory")
	}
	if dAtA.BaoZangDestory == nil {
		return errors.Errorf("%s 配置的关联字段[bao_zang_destory] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("bao_zang_destory"), *pArSeR)
	}

	if pArSeR.KeyExist("base_level") {
		dAtA.BaseLevel = cOnFigS.GetBroadcastData(pArSeR.String("base_level"))
	} else {
		dAtA.BaseLevel = cOnFigS.GetBroadcastData("BaseLevel")
	}
	if dAtA.BaseLevel == nil {
		return errors.Errorf("%s 配置的关联字段[base_level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("base_level"), *pArSeR)
	}

	if pArSeR.KeyExist("captain_ability_exp") {
		dAtA.CaptainAbilityExp = cOnFigS.GetBroadcastData(pArSeR.String("captain_ability_exp"))
	} else {
		dAtA.CaptainAbilityExp = cOnFigS.GetBroadcastData("CaptainAbilityExp")
	}
	if dAtA.CaptainAbilityExp == nil {
		return errors.Errorf("%s 配置的关联字段[captain_ability_exp] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain_ability_exp"), *pArSeR)
	}

	if pArSeR.KeyExist("captain_re_brith") {
		dAtA.CaptainReBrith = cOnFigS.GetBroadcastData(pArSeR.String("captain_re_brith"))
	} else {
		dAtA.CaptainReBrith = cOnFigS.GetBroadcastData("CaptainReBrith")
	}
	if dAtA.CaptainReBrith == nil {
		return errors.Errorf("%s 配置的关联字段[captain_re_brith] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain_re_brith"), *pArSeR)
	}

	if pArSeR.KeyExist("captain_soul_level") {
		dAtA.CaptainSoulLevel = cOnFigS.GetBroadcastData(pArSeR.String("captain_soul_level"))
	} else {
		dAtA.CaptainSoulLevel = cOnFigS.GetBroadcastData("CaptainSoulLevel")
	}
	if dAtA.CaptainSoulLevel == nil {
		return errors.Errorf("%s 配置的关联字段[captain_soul_level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain_soul_level"), *pArSeR)
	}

	if pArSeR.KeyExist("captain_soul_unlock_spell") {
		dAtA.CaptainSoulUnlockSpell = cOnFigS.GetBroadcastData(pArSeR.String("captain_soul_unlock_spell"))
	} else {
		dAtA.CaptainSoulUnlockSpell = cOnFigS.GetBroadcastData("CaptainSoulUnlockSpell")
	}
	if dAtA.CaptainSoulUnlockSpell == nil {
		return errors.Errorf("%s 配置的关联字段[captain_soul_unlock_spell] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain_soul_unlock_spell"), *pArSeR)
	}

	if pArSeR.KeyExist("country_change_king") {
		dAtA.CountryChangeKing = cOnFigS.GetBroadcastData(pArSeR.String("country_change_king"))
	} else {
		dAtA.CountryChangeKing = cOnFigS.GetBroadcastData("CountryChangeKing")
	}
	if dAtA.CountryChangeKing == nil {
		return errors.Errorf("%s 配置的关联字段[country_change_king] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("country_change_king"), *pArSeR)
	}

	if pArSeR.KeyExist("country_change_name") {
		dAtA.CountryChangeName = cOnFigS.GetBroadcastData(pArSeR.String("country_change_name"))
	} else {
		dAtA.CountryChangeName = cOnFigS.GetBroadcastData("CountryChangeName")
	}
	if dAtA.CountryChangeName == nil {
		return errors.Errorf("%s 配置的关联字段[country_change_name] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("country_change_name"), *pArSeR)
	}

	if pArSeR.KeyExist("country_destroy") {
		dAtA.CountryDestroy = cOnFigS.GetBroadcastData(pArSeR.String("country_destroy"))
	} else {
		dAtA.CountryDestroy = cOnFigS.GetBroadcastData("CountryDestroy")
	}
	if dAtA.CountryDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[country_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("country_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("dungeon_prize_captain_soul") {
		dAtA.DungeonPrizeCaptainSoul = cOnFigS.GetBroadcastData(pArSeR.String("dungeon_prize_captain_soul"))
	} else {
		dAtA.DungeonPrizeCaptainSoul = cOnFigS.GetBroadcastData("DungeonPrizeCaptainSoul")
	}
	if dAtA.DungeonPrizeCaptainSoul == nil {
		return errors.Errorf("%s 配置的关联字段[dungeon_prize_captain_soul] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("dungeon_prize_captain_soul"), *pArSeR)
	}

	if pArSeR.KeyExist("equip_level") {
		dAtA.EquipLevel = cOnFigS.GetBroadcastData(pArSeR.String("equip_level"))
	} else {
		dAtA.EquipLevel = cOnFigS.GetBroadcastData("EquipLevel")
	}
	if dAtA.EquipLevel == nil {
		return errors.Errorf("%s 配置的关联字段[equip_level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("equip_level"), *pArSeR)
	}

	if pArSeR.KeyExist("fish_captain_soul") {
		dAtA.FishCaptainSoul = cOnFigS.GetBroadcastData(pArSeR.String("fish_captain_soul"))
	} else {
		dAtA.FishCaptainSoul = cOnFigS.GetBroadcastData("FishCaptainSoul")
	}
	if dAtA.FishCaptainSoul == nil {
		return errors.Errorf("%s 配置的关联字段[fish_captain_soul] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fish_captain_soul"), *pArSeR)
	}

	if pArSeR.KeyExist("fish_equip") {
		dAtA.FishEquip = cOnFigS.GetBroadcastData(pArSeR.String("fish_equip"))
	} else {
		dAtA.FishEquip = cOnFigS.GetBroadcastData("FishEquip")
	}
	if dAtA.FishEquip == nil {
		return errors.Errorf("%s 配置的关联字段[fish_equip] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fish_equip"), *pArSeR)
	}

	if pArSeR.KeyExist("fish_gem") {
		dAtA.FishGem = cOnFigS.GetBroadcastData(pArSeR.String("fish_gem"))
	} else {
		dAtA.FishGem = cOnFigS.GetBroadcastData("FishGem")
	}
	if dAtA.FishGem == nil {
		return errors.Errorf("%s 配置的关联字段[fish_gem] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fish_gem"), *pArSeR)
	}

	if pArSeR.KeyExist("gem_get") {
		dAtA.GemGet = cOnFigS.GetBroadcastData(pArSeR.String("gem_get"))
	} else {
		dAtA.GemGet = cOnFigS.GetBroadcastData("GemGet")
	}
	if dAtA.GemGet == nil {
		return errors.Errorf("%s 配置的关联字段[gem_get] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("gem_get"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_change_country") {
		dAtA.GuildChangeCountry = cOnFigS.GetBroadcastData(pArSeR.String("guild_change_country"))
	} else {
		dAtA.GuildChangeCountry = cOnFigS.GetBroadcastData("GuildChangeCountry")
	}
	if dAtA.GuildChangeCountry == nil {
		return errors.Errorf("%s 配置的关联字段[guild_change_country] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_change_country"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_change_flag") {
		dAtA.GuildChangeFlag = cOnFigS.GetBroadcastData(pArSeR.String("guild_change_flag"))
	} else {
		dAtA.GuildChangeFlag = cOnFigS.GetBroadcastData("GuildChangeFlag")
	}
	if dAtA.GuildChangeFlag == nil {
		return errors.Errorf("%s 配置的关联字段[guild_change_flag] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_change_flag"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_change_name") {
		dAtA.GuildChangeName = cOnFigS.GetBroadcastData(pArSeR.String("guild_change_name"))
	} else {
		dAtA.GuildChangeName = cOnFigS.GetBroadcastData("GuildChangeName")
	}
	if dAtA.GuildChangeName == nil {
		return errors.Errorf("%s 配置的关联字段[guild_change_name] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_change_name"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_create") {
		dAtA.GuildCreate = cOnFigS.GetBroadcastData(pArSeR.String("guild_create"))
	} else {
		dAtA.GuildCreate = cOnFigS.GetBroadcastData("GuildCreate")
	}
	if dAtA.GuildCreate == nil {
		return errors.Errorf("%s 配置的关联字段[guild_create] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_create"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_level") {
		dAtA.GuildLevel = cOnFigS.GetBroadcastData(pArSeR.String("guild_level"))
	} else {
		dAtA.GuildLevel = cOnFigS.GetBroadcastData("GuildLevel")
	}
	if dAtA.GuildLevel == nil {
		return errors.Errorf("%s 配置的关联字段[guild_level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_level"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_tan_he") {
		dAtA.GuildTanHe = cOnFigS.GetBroadcastData(pArSeR.String("guild_tan_he"))
	} else {
		dAtA.GuildTanHe = cOnFigS.GetBroadcastData("GuildTanHe")
	}
	if dAtA.GuildTanHe == nil {
		return errors.Errorf("%s 配置的关联字段[guild_tan_he] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_tan_he"), *pArSeR)
	}

	if pArSeR.KeyExist("hero_level") {
		dAtA.HeroLevel = cOnFigS.GetBroadcastData(pArSeR.String("hero_level"))
	} else {
		dAtA.HeroLevel = cOnFigS.GetBroadcastData("HeroLevel")
	}
	if dAtA.HeroLevel == nil {
		return errors.Errorf("%s 配置的关联字段[hero_level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("hero_level"), *pArSeR)
	}

	if pArSeR.KeyExist("jiu_guan_bao_ji") {
		dAtA.JiuGuanBaoJi = cOnFigS.GetBroadcastData(pArSeR.String("jiu_guan_bao_ji"))
	} else {
		dAtA.JiuGuanBaoJi = cOnFigS.GetBroadcastData("JiuGuanBaoJi")
	}
	if dAtA.JiuGuanBaoJi == nil {
		return errors.Errorf("%s 配置的关联字段[jiu_guan_bao_ji] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("jiu_guan_bao_ji"), *pArSeR)
	}

	if pArSeR.KeyExist("jiu_guan_qing_jiao") {
		dAtA.JiuGuanQingJiao = cOnFigS.GetBroadcastData(pArSeR.String("jiu_guan_qing_jiao"))
	} else {
		dAtA.JiuGuanQingJiao = cOnFigS.GetBroadcastData("JiuGuanQingJiao")
	}
	if dAtA.JiuGuanQingJiao == nil {
		return errors.Errorf("%s 配置的关联字段[jiu_guan_qing_jiao] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("jiu_guan_qing_jiao"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_atk_win") {
		dAtA.McWarAtkWin = cOnFigS.GetBroadcastData(pArSeR.String("mc_war_atk_win"))
	} else {
		dAtA.McWarAtkWin = cOnFigS.GetBroadcastData("McWarAtkWin")
	}
	if dAtA.McWarAtkWin == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_atk_win] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_atk_win"), *pArSeR)
	}

	if pArSeR.KeyExist("mix_equip") {
		dAtA.MixEquip = cOnFigS.GetBroadcastData(pArSeR.String("mix_equip"))
	} else {
		dAtA.MixEquip = cOnFigS.GetBroadcastData("MixEquip")
	}
	if dAtA.MixEquip == nil {
		return errors.Errorf("%s 配置的关联字段[mix_equip] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mix_equip"), *pArSeR)
	}

	if pArSeR.KeyExist("npc_monster_succ") {
		dAtA.NpcMonsterSucc = cOnFigS.GetBroadcastData(pArSeR.String("npc_monster_succ"))
	} else {
		dAtA.NpcMonsterSucc = cOnFigS.GetBroadcastData("NpcMonsterSucc")
	}
	if dAtA.NpcMonsterSucc == nil {
		return errors.Errorf("%s 配置的关联字段[npc_monster_succ] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("npc_monster_succ"), *pArSeR)
	}

	if pArSeR.KeyExist("taoz_level") {
		dAtA.TaozLevel = cOnFigS.GetBroadcastData(pArSeR.String("taoz_level"))
	} else {
		dAtA.TaozLevel = cOnFigS.GetBroadcastData("TaozLevel")
	}
	if dAtA.TaozLevel == nil {
		return errors.Errorf("%s 配置的关联字段[taoz_level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("taoz_level"), *pArSeR)
	}

	if pArSeR.KeyExist("task_bwzl_prize") {
		dAtA.TaskBwzlPrize = cOnFigS.GetBroadcastData(pArSeR.String("task_bwzl_prize"))
	} else {
		dAtA.TaskBwzlPrize = cOnFigS.GetBroadcastData("TaskBwzlPrize")
	}
	if dAtA.TaskBwzlPrize == nil {
		return errors.Errorf("%s 配置的关联字段[task_bwzl_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("task_bwzl_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("task_cheng_jiu_level") {
		dAtA.TaskChengJiuLevel = cOnFigS.GetBroadcastData(pArSeR.String("task_cheng_jiu_level"))
	} else {
		dAtA.TaskChengJiuLevel = cOnFigS.GetBroadcastData("TaskChengJiuLevel")
	}
	if dAtA.TaskChengJiuLevel == nil {
		return errors.Errorf("%s 配置的关联字段[task_cheng_jiu_level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("task_cheng_jiu_level"), *pArSeR)
	}

	if pArSeR.KeyExist("title") {
		dAtA.Title = cOnFigS.GetBroadcastData(pArSeR.String("title"))
	} else {
		dAtA.Title = cOnFigS.GetBroadcastData("Title")
	}
	if dAtA.Title == nil {
		return errors.Errorf("%s 配置的关联字段[title] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("title"), *pArSeR)
	}

	if pArSeR.KeyExist("tower_floor") {
		dAtA.TowerFloor = cOnFigS.GetBroadcastData(pArSeR.String("tower_floor"))
	} else {
		dAtA.TowerFloor = cOnFigS.GetBroadcastData("TowerFloor")
	}
	if dAtA.TowerFloor == nil {
		return errors.Errorf("%s 配置的关联字段[tower_floor] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("tower_floor"), *pArSeR)
	}

	if pArSeR.KeyExist("tower_prize_equip") {
		dAtA.TowerPrizeEquip = cOnFigS.GetBroadcastData(pArSeR.String("tower_prize_equip"))
	} else {
		dAtA.TowerPrizeEquip = cOnFigS.GetBroadcastData("TowerPrizeEquip")
	}
	if dAtA.TowerPrizeEquip == nil {
		return errors.Errorf("%s 配置的关联字段[tower_prize_equip] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("tower_prize_equip"), *pArSeR)
	}

	if pArSeR.KeyExist("xiong_nu_succ") {
		dAtA.XiongNuSucc = cOnFigS.GetBroadcastData(pArSeR.String("xiong_nu_succ"))
	} else {
		dAtA.XiongNuSucc = cOnFigS.GetBroadcastData("XiongNuSucc")
	}
	if dAtA.XiongNuSucc == nil {
		return errors.Errorf("%s 配置的关联字段[xiong_nu_succ] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("xiong_nu_succ"), *pArSeR)
	}

	if pArSeR.KeyExist("zhan_jiang_zhang_jie_complete") {
		dAtA.ZhanJiangZhangJieComplete = cOnFigS.GetBroadcastData(pArSeR.String("zhan_jiang_zhang_jie_complete"))
	} else {
		dAtA.ZhanJiangZhangJieComplete = cOnFigS.GetBroadcastData("ZhanJiangZhangJieComplete")
	}
	if dAtA.ZhanJiangZhangJieComplete == nil {
		return errors.Errorf("%s 配置的关联字段[zhan_jiang_zhang_jie_complete] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("zhan_jiang_zhang_jie_complete"), *pArSeR)
	}

	if pArSeR.KeyExist("zhen_bao_ge_bao_ji") {
		dAtA.ZhenBaoGeBaoJi = cOnFigS.GetBroadcastData(pArSeR.String("zhen_bao_ge_bao_ji"))
	} else {
		dAtA.ZhenBaoGeBaoJi = cOnFigS.GetBroadcastData("ZhenBaoGeBaoJi")
	}
	if dAtA.ZhenBaoGeBaoJi == nil {
		return errors.Errorf("%s 配置的关联字段[zhen_bao_ge_bao_ji] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("zhen_bao_ge_bao_ji"), *pArSeR)
	}

	return nil
}

// start with BuffEffectData ----------------------------------

func LoadBuffEffectData(gos *config.GameObjects) (map[uint64]*BuffEffectData, map[*BuffEffectData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BuffEffectDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BuffEffectData, len(lIsT))
	pArSeRmAp := make(map[*BuffEffectData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBuffEffectData) {
			continue
		}

		dAtA, err := NewBuffEffectData(fIlEnAmE, pArSeR)
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

func SetRelatedBuffEffectData(dAtAmAp map[*BuffEffectData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BuffEffectDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBuffEffectDataKeyArray(datas []*BuffEffectData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBuffEffectData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BuffEffectData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBuffEffectData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BuffEffectData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Group = pArSeR.Uint64("group")
	dAtA.EffectType = shared_proto.BuffEffectType(shared_proto.BuffEffectType_value[strings.ToUpper(pArSeR.String("effect_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("effect_type"), 10, 32); err == nil {
		dAtA.EffectType = shared_proto.BuffEffectType(i)
	}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.KeepDuration, err = config.ParseDuration(pArSeR.String("keep_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[keep_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("keep_duration"), dAtA)
	}

	dAtA.NoDuration = pArSeR.Bool("no_duration")
	dAtA.PvpBuff = pArSeR.Bool("pvp_buff")
	// releated field: StatBuff
	dAtA.CaptainTrain, err = ParseAmount(pArSeR.String("captain_train"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[captain_train] 解析失败(ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("captain_train"), dAtA)
	}

	dAtA.FarmHarvest, err = ParseAmount(pArSeR.String("farm_harvest"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[farm_harvest] 解析失败(ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("farm_harvest"), dAtA)
	}

	dAtA.Tax, err = ParseAmount(pArSeR.String("tax"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[tax] 解析失败(ParseAmount)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("tax"), dAtA)
	}

	dAtA.AdvantageId = pArSeR.Uint64("advantage_id")

	return dAtA, nil
}

var vAlIdAtOrBuffEffectData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"name":          config.ParseValidator("string", "", false, nil, nil),
	"desc":          config.ParseValidator("string", "", false, nil, nil),
	"group":         config.ParseValidator("int>0", "", false, nil, nil),
	"effect_type":   config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.BuffEffectType_value, 0), nil),
	"level":         config.ParseValidator("int>0", "", false, nil, nil),
	"keep_duration": config.ParseValidator("string", "", false, nil, nil),
	"no_duration":   config.ParseValidator("bool", "", false, nil, nil),
	"pvp_buff":      config.ParseValidator("bool", "", false, nil, nil),
	"stat_buff":     config.ParseValidator("string", "", false, nil, nil),
	"captain_train": config.ParseValidator("string", "", false, nil, nil),
	"farm_harvest":  config.ParseValidator("string", "", false, nil, nil),
	"tax":           config.ParseValidator("string", "", false, nil, nil),
	"advantage_id":  config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *BuffEffectData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BuffEffectData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BuffEffectData) Encode() *shared_proto.BuffEffectDataProto {
	out := &shared_proto.BuffEffectDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.Group = config.U64ToI32(dAtA.Group)
	out.EffectType = dAtA.EffectType
	out.Level = config.U64ToI32(dAtA.Level)
	out.KeepDuration = config.Duration2I32Seconds(dAtA.KeepDuration)
	out.NoDuration = dAtA.NoDuration
	out.PvpBuff = dAtA.PvpBuff
	if dAtA.StatBuff != nil {
		out.StatBuff = dAtA.StatBuff.Encode()
	}
	if dAtA.CaptainTrain != nil {
		out.CaptainTrain = dAtA.CaptainTrain.Encode()
	}
	if dAtA.FarmHarvest != nil {
		out.FarmHarvest = dAtA.FarmHarvest.Encode()
	}
	if dAtA.Tax != nil {
		out.Tax = dAtA.Tax.Encode()
	}
	out.AdvantageId = config.U64ToI32(dAtA.AdvantageId)

	return out
}

func ArrayEncodeBuffEffectData(datas []*BuffEffectData) []*shared_proto.BuffEffectDataProto {

	out := make([]*shared_proto.BuffEffectDataProto, 0, len(datas))
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

func (dAtA *BuffEffectData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.StatBuff = cOnFigS.GetSpriteStat(pArSeR.Uint64("stat_buff"))
	if dAtA.StatBuff == nil && pArSeR.Uint64("stat_buff") != 0 {
		return errors.Errorf("%s 配置的关联字段[stat_buff] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("stat_buff"), *pArSeR)
	}

	return nil
}

// start with ColorData ----------------------------------

func LoadColorData(gos *config.GameObjects) (map[uint64]*ColorData, map[*ColorData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ColorDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ColorData, len(lIsT))
	pArSeRmAp := make(map[*ColorData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrColorData) {
			continue
		}

		dAtA, err := NewColorData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.QualityKey
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[QualityKey], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedColorData(dAtAmAp map[*ColorData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ColorDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetColorDataKeyArray(datas []*ColorData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.QualityKey)
		}
	}

	return out
}

func NewColorData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ColorData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrColorData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ColorData{}

	dAtA.QualityKey = pArSeR.Uint64("quality_key")
	// skip field: Quality
	dAtA.ColorCode = pArSeR.String("color_code")
	dAtA.ColorName = pArSeR.String("color_name")
	// skip field: QualityColorText
	dAtA.CaptainSoulQualityName = "S"
	if pArSeR.KeyExist("captain_soul_quality_name") {
		dAtA.CaptainSoulQualityName = pArSeR.String("captain_soul_quality_name")
	}

	return dAtA, nil
}

var vAlIdAtOrColorData = map[string]*config.Validator{

	"quality_key":               config.ParseValidator("uint", "", false, nil, nil),
	"color_code":                config.ParseValidator("string", "", false, nil, nil),
	"color_name":                config.ParseValidator("string", "", false, nil, nil),
	"captain_soul_quality_name": config.ParseValidator("string", "", false, nil, []string{"S"}),
}

func (dAtA *ColorData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with CompareCondition ----------------------------------

func NewCompareCondition(fIlEnAmE string, pArSeR *config.ObjectParser) (*CompareCondition, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCompareCondition)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CompareCondition{}

	dAtA.Greater = pArSeR.Bool("greater")
	dAtA.Less = pArSeR.Bool("less")
	dAtA.Equal = pArSeR.Bool("equal")
	dAtA.Amount = pArSeR.Uint64("amount")

	return dAtA, nil
}

var vAlIdAtOrCompareCondition = map[string]*config.Validator{

	"greater": config.ParseValidator("bool", "", false, nil, nil),
	"less":    config.ParseValidator("bool", "", false, nil, nil),
	"equal":   config.ParseValidator("bool", "", false, nil, nil),
	"amount":  config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *CompareCondition) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with FamilyName ----------------------------------

func LoadFamilyName(gos *config.GameObjects) (map[string]*FamilyName, map[*FamilyName]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FamilyNamePath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*FamilyName, len(lIsT))
	pArSeRmAp := make(map[*FamilyName]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFamilyName) {
			continue
		}

		dAtA, err := NewFamilyName(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.FamilyName
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[FamilyName], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedFamilyName(dAtAmAp map[*FamilyName]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FamilyNamePath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFamilyNameKeyArray(datas []*FamilyName) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.FamilyName)
		}
	}

	return out
}

func NewFamilyName(fIlEnAmE string, pArSeR *config.ObjectParser) (*FamilyName, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFamilyName)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FamilyName{}

	dAtA.FamilyName = pArSeR.String("family_name")

	return dAtA, nil
}

var vAlIdAtOrFamilyName = map[string]*config.Validator{

	"family_name": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *FamilyName) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with FemaleGivenName ----------------------------------

func LoadFemaleGivenName(gos *config.GameObjects) (map[string]*FemaleGivenName, map[*FemaleGivenName]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FemaleGivenNamePath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*FemaleGivenName, len(lIsT))
	pArSeRmAp := make(map[*FemaleGivenName]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFemaleGivenName) {
			continue
		}

		dAtA, err := NewFemaleGivenName(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Name
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Name], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedFemaleGivenName(dAtAmAp map[*FemaleGivenName]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FemaleGivenNamePath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFemaleGivenNameKeyArray(datas []*FemaleGivenName) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Name)
		}
	}

	return out
}

func NewFemaleGivenName(fIlEnAmE string, pArSeR *config.ObjectParser) (*FemaleGivenName, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFemaleGivenName)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FemaleGivenName{}

	dAtA.Name = pArSeR.String("name")

	return dAtA, nil
}

var vAlIdAtOrFemaleGivenName = map[string]*config.Validator{

	"name": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *FemaleGivenName) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with HeroLevelSubData ----------------------------------

func LoadHeroLevelSubData(gos *config.GameObjects) (map[uint64]*HeroLevelSubData, map[*HeroLevelSubData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.HeroLevelSubDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*HeroLevelSubData, len(lIsT))
	pArSeRmAp := make(map[*HeroLevelSubData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrHeroLevelSubData) {
			continue
		}

		dAtA, err := NewHeroLevelSubData(fIlEnAmE, pArSeR)
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

func SetRelatedHeroLevelSubData(dAtAmAp map[*HeroLevelSubData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.HeroLevelSubDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetHeroLevelSubDataKeyArray(datas []*HeroLevelSubData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewHeroLevelSubData(fIlEnAmE string, pArSeR *config.ObjectParser) (*HeroLevelSubData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrHeroLevelSubData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &HeroLevelSubData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.UpgradeExp = pArSeR.Uint64("upgrade_exp")
	dAtA.AddSoldierCapacity = 0
	if pArSeR.KeyExist("add_soldier_capacity") {
		dAtA.AddSoldierCapacity = pArSeR.Uint64("add_soldier_capacity")
	}

	dAtA.EquipmentLevelLimit = pArSeR.Uint64("equipment_level_limit")
	dAtA.CaptainSoulLevelLimit = pArSeR.Uint64("captain_soul_level_limit")
	dAtA.CaptainLevelLimit = pArSeR.Uint64("captain_level_limit")
	for _, v := range pArSeR.StringArray("unlocked_races", "", false) {
		x := shared_proto.Race(shared_proto.Race_value[strings.ToUpper(v)])
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			x = shared_proto.Race(i)
		}
		dAtA.UnlockedRaces = append(dAtA.UnlockedRaces, x)
	}

	dAtA.CaptainTrainingLevel = pArSeR.Uint64Array("captain_training_level", "", false)
	dAtA.CaptainTrainingLevelLimit = pArSeR.Uint64Array("captain_training_level_limit", "", false)
	dAtA.StrategyLimit = pArSeR.Uint64("strategy_limit")
	dAtA.SpLimit = pArSeR.Uint64("sp_limit")
	dAtA.TroopsCount = 3
	if pArSeR.KeyExist("troops_count") {
		dAtA.TroopsCount = pArSeR.Uint64("troops_count")
	}

	dAtA.TroopsCaptainCount = 2
	if pArSeR.KeyExist("troops_captain_count") {
		dAtA.TroopsCaptainCount = pArSeR.Uint64("troops_captain_count")
	}

	// skip field: CaptainOfficialId
	// skip field: CaptainOfficialCount
	// skip field: CaptainOfficialIdCount

	return dAtA, nil
}

var vAlIdAtOrHeroLevelSubData = map[string]*config.Validator{

	"level":                        config.ParseValidator("int>0", "", false, nil, nil),
	"upgrade_exp":                  config.ParseValidator("int>0", "", false, nil, nil),
	"add_soldier_capacity":         config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"equipment_level_limit":        config.ParseValidator("int>0", "", false, nil, nil),
	"captain_soul_level_limit":     config.ParseValidator("int>0", "", false, nil, nil),
	"captain_level_limit":          config.ParseValidator("int>0", "", false, nil, nil),
	"unlocked_races":               config.ParseValidator("string,notAllNil", "", true, config.EnumMapKeys(shared_proto.Race_value, 0), nil),
	"captain_training_level":       config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
	"captain_training_level_limit": config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
	"strategy_limit":               config.ParseValidator("int>0", "", false, nil, nil),
	"sp_limit":                     config.ParseValidator("int>0", "", false, nil, nil),
	"troops_count":                 config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"troops_captain_count":         config.ParseValidator("int>0", "", false, nil, []string{"2"}),
}

func (dAtA *HeroLevelSubData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *HeroLevelSubData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *HeroLevelSubData) Encode() *shared_proto.HeroLevelProto {
	out := &shared_proto.HeroLevelProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.UpgradeExp = config.U64ToI32(dAtA.UpgradeExp)
	out.AddSoldierCapacity = config.U64ToI32(dAtA.AddSoldierCapacity)
	out.EquipmentLevelLimit = config.U64ToI32(dAtA.EquipmentLevelLimit)
	out.CaptainSoulLevelLimit = config.U64ToI32(dAtA.CaptainSoulLevelLimit)
	out.CaptainLevelLimit = config.U64ToI32(dAtA.CaptainLevelLimit)
	out.UnlockedRaces = dAtA.UnlockedRaces
	out.StrategyLimit = config.U64ToI32(dAtA.StrategyLimit)
	out.SpLimit = config.U64ToI32(dAtA.SpLimit)
	out.TroopsCount = config.U64ToI32(dAtA.TroopsCount)
	out.TroopsCaptainCount = config.U64ToI32(dAtA.TroopsCaptainCount)
	out.CaptainOfficialId = config.U64a2I32a(dAtA.CaptainOfficialId)
	out.CaptainOfficialCount = config.U64a2I32a(dAtA.CaptainOfficialCount)

	return out
}

func ArrayEncodeHeroLevelSubData(datas []*HeroLevelSubData) []*shared_proto.HeroLevelProto {

	out := make([]*shared_proto.HeroLevelProto, 0, len(datas))
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

func (dAtA *HeroLevelSubData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with MaleGivenName ----------------------------------

func LoadMaleGivenName(gos *config.GameObjects) (map[string]*MaleGivenName, map[*MaleGivenName]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MaleGivenNamePath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*MaleGivenName, len(lIsT))
	pArSeRmAp := make(map[*MaleGivenName]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMaleGivenName) {
			continue
		}

		dAtA, err := NewMaleGivenName(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Name
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Name], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedMaleGivenName(dAtAmAp map[*MaleGivenName]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MaleGivenNamePath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMaleGivenNameKeyArray(datas []*MaleGivenName) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Name)
		}
	}

	return out
}

func NewMaleGivenName(fIlEnAmE string, pArSeR *config.ObjectParser) (*MaleGivenName, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMaleGivenName)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MaleGivenName{}

	dAtA.Name = pArSeR.String("name")

	return dAtA, nil
}

var vAlIdAtOrMaleGivenName = map[string]*config.Validator{

	"name": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *MaleGivenName) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with RandAmount ----------------------------------

func NewRandAmount(fIlEnAmE string, pArSeR *config.ObjectParser) (*RandAmount, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRandAmount)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RandAmount{}

	// skip field: Min
	// skip field: Max

	return dAtA, nil
}

var vAlIdAtOrRandAmount = map[string]*config.Validator{}

func (dAtA *RandAmount) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with Rate ----------------------------------

func NewRate(fIlEnAmE string, pArSeR *config.ObjectParser) (*Rate, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRate)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &Rate{}

	// skip field: Numerator
	// skip field: Denominator

	return dAtA, nil
}

var vAlIdAtOrRate = map[string]*config.Validator{}

func (dAtA *Rate) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with ResAmount ----------------------------------

func NewResAmount(fIlEnAmE string, pArSeR *config.ObjectParser) (*ResAmount, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrResAmount)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ResAmount{}

	dAtA.Type = shared_proto.ResType(shared_proto.ResType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.ResType(i)
	}

	dAtA.Amount = pArSeR.Uint64("amount")
	dAtA.Percent = pArSeR.Uint64("percent")

	return dAtA, nil
}

var vAlIdAtOrResAmount = map[string]*config.Validator{

	"type":    config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.ResType_value, 0), nil),
	"amount":  config.ParseValidator("uint", "", false, nil, nil),
	"percent": config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *ResAmount) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ResAmount) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ResAmount) Encode() *shared_proto.ResProto {
	out := &shared_proto.ResProto{}
	out.Type = dAtA.Type
	out.Amount = config.U64ToI32(dAtA.Amount)
	out.Percent = config.U64ToI32(dAtA.Percent)

	return out
}

func ArrayEncodeResAmount(datas []*ResAmount) []*shared_proto.ResProto {

	out := make([]*shared_proto.ResProto, 0, len(datas))
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

func (dAtA *ResAmount) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with SpriteStat ----------------------------------

func LoadSpriteStat(gos *config.GameObjects) (map[uint64]*SpriteStat, map[*SpriteStat]*config.ObjectParser, error) {
	fIlEnAmE := confpath.SpriteStatPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*SpriteStat, len(lIsT))
	pArSeRmAp := make(map[*SpriteStat]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrSpriteStat) {
			continue
		}

		dAtA, err := NewSpriteStat(fIlEnAmE, pArSeR)
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

func SetRelatedSpriteStat(dAtAmAp map[*SpriteStat]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SpriteStatPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetSpriteStatKeyArray(datas []*SpriteStat) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewSpriteStat(fIlEnAmE string, pArSeR *config.ObjectParser) (*SpriteStat, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSpriteStat)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SpriteStat{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Attack = pArSeR.Uint64("attack")
	dAtA.Defense = pArSeR.Uint64("defense")
	dAtA.Strength = pArSeR.Uint64("strength")
	dAtA.Dexterity = pArSeR.Uint64("dexterity")
	dAtA.SoldierCapcity = pArSeR.Uint64("soldier_capcity")
	dAtA.DamageIncrePer = pArSeR.Uint64("damage_incre_per")
	dAtA.DamageDecrePer = pArSeR.Uint64("damage_decre_per")
	dAtA.BeenHurtIncrePer = pArSeR.Uint64("been_hurt_incre_per")
	dAtA.BeenHurtDecrePer = pArSeR.Uint64("been_hurt_decre_per")

	return dAtA, nil
}

var vAlIdAtOrSpriteStat = map[string]*config.Validator{

	"id":                  config.ParseValidator("int>0", "", false, nil, nil),
	"attack":              config.ParseValidator("uint", "", false, nil, nil),
	"defense":             config.ParseValidator("uint", "", false, nil, nil),
	"strength":            config.ParseValidator("uint", "", false, nil, nil),
	"dexterity":           config.ParseValidator("uint", "", false, nil, nil),
	"soldier_capcity":     config.ParseValidator("uint", "", false, nil, nil),
	"damage_incre_per":    config.ParseValidator("int", "", false, nil, nil),
	"damage_decre_per":    config.ParseValidator("int", "", false, nil, nil),
	"been_hurt_incre_per": config.ParseValidator("int", "", false, nil, nil),
	"been_hurt_decre_per": config.ParseValidator("int", "", false, nil, nil),
}

func (dAtA *SpriteStat) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SpriteStat) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SpriteStat) Encode() *shared_proto.SpriteStatProto {
	out := &shared_proto.SpriteStatProto{}
	out.Attack = config.U64ToI32(dAtA.Attack)
	out.Defense = config.U64ToI32(dAtA.Defense)
	out.Strength = config.U64ToI32(dAtA.Strength)
	out.Dexterity = config.U64ToI32(dAtA.Dexterity)
	out.SoldierCapcity = config.U64ToI32(dAtA.SoldierCapcity)
	out.DamageIncrePer = config.U64ToI32(dAtA.DamageIncrePer)
	out.DamageDecrePer = config.U64ToI32(dAtA.DamageDecrePer)
	out.BeenHurtIncrePer = config.U64ToI32(dAtA.BeenHurtIncrePer)
	out.BeenHurtDecrePer = config.U64ToI32(dAtA.BeenHurtDecrePer)

	return out
}

func ArrayEncodeSpriteStat(datas []*SpriteStat) []*shared_proto.SpriteStatProto {

	out := make([]*shared_proto.SpriteStatProto, 0, len(datas))
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

func (dAtA *SpriteStat) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with Text ----------------------------------

func LoadText(gos *config.GameObjects) (map[string]*Text, map[*Text]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TextPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*Text, len(lIsT))
	pArSeRmAp := make(map[*Text]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrText) {
			continue
		}

		dAtA, err := NewText(fIlEnAmE, pArSeR)
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

func SetRelatedText(dAtAmAp map[*Text]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TextPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTextKeyArray(datas []*Text) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewText(fIlEnAmE string, pArSeR *config.ObjectParser) (*Text, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrText)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &Text{}

	dAtA.Id = pArSeR.String("id")

	// i18n fields
	dAtA.Text = i18n.NewI18nRef(fIlEnAmE, "text", dAtA.Id, pArSeR.String("text"))

	return dAtA, nil
}

var vAlIdAtOrText = map[string]*config.Validator{

	"id":   config.ParseValidator("string", "", false, nil, nil),
	"text": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *Text) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with Text2 ----------------------------------

func NewText2(fIlEnAmE string, pArSeR *config.ObjectParser) (*Text2, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrText2)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &Text2{}

	return dAtA, nil
}

var vAlIdAtOrText2 = map[string]*config.Validator{}

func (dAtA *Text2) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with TextFormatter ----------------------------------

func NewTextFormatter(fIlEnAmE string, pArSeR *config.ObjectParser) (*TextFormatter, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTextFormatter)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TextFormatter{}

	dAtA.Text = pArSeR.StringArray("text", "", false)
	// skip field: OneText

	return dAtA, nil
}

var vAlIdAtOrTextFormatter = map[string]*config.Validator{

	"text": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *TextFormatter) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with TextHelp ----------------------------------

func LoadTextHelp(gos *config.GameObjects) (*TextHelp, *config.ObjectParser, error) {
	fIlEnAmE := confpath.TextHelpPath
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

	dAtA, err := NewTextHelp(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedTextHelp(gos *config.GameObjects, dAtA *TextHelp, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TextHelpPath
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

func NewTextHelp(fIlEnAmE string, pArSeR *config.ObjectParser) (*TextHelp, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTextHelp)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TextHelp{}

	// releated field: BannerBaseFail
	// releated field: BannerBaseSuccess
	// releated field: BannerTroopFail
	// releated field: BannerTroopSuccess
	// releated field: GuildGongfangBuildComplete
	// releated field: GuildGongfangBuildTimeAdd
	// releated field: GuildGongfangBuildTimeReduce
	// releated field: GuildGongfangBuilding
	// releated field: GuildGongfangEfficiencyAdd
	// releated field: GuildGongfangEfficiencyReduce
	// releated field: GuildGongfangPrizeSend
	// releated field: GuildPleaseHelpMe
	// releated field: GuildSayHi
	// releated field: GuildWorkshopBuildShow
	// releated field: GuildWorkshopCompletedChat
	// releated field: GuildWorkshopCreatedChat
	// releated field: GuildWorkshopPrizeShow
	// releated field: GuildWorkshopProdShow
	// releated field: MRDRAttBroken4a
	// releated field: MRDRAttBroken4d
	// releated field: MRDRAttMoveBase4a
	// releated field: MRDRAttMoveBase4d
	// releated field: MRDRBroken4a
	// releated field: MRDRBroken4d
	// releated field: MRDRExpelled4a
	// releated field: MRDRExpelled4d
	// releated field: MRDRFinished4a
	// releated field: MRDRFinished4d
	// releated field: MRDRMian4a
	// releated field: MRDRMian4d
	// releated field: MRDRMoveBase4a
	// releated field: MRDRMoveBase4d
	// releated field: MRDRProtected4a
	// releated field: MRDRProtected4d
	// releated field: MRDRRecall4a
	// releated field: MRDRRecall4d
	// releated field: MailAddHateDesc
	// releated field: MailLostProsperityDesc
	// releated field: MailReduceProsperityDesc
	// releated field: MailReportLoserDesc
	// releated field: MailReportWinnerDesc
	// releated field: McWarAtkFail
	// releated field: McWarChatEnemyAtkBuildingDestroy
	// releated field: McWarChatEnemyDefBuildingDestroy
	// releated field: McWarChatOurAtkBuildingDestroy
	// releated field: McWarChatOurDefBuildingDestroy
	// releated field: McWarDefFail
	// releated field: McWarReliveSucc
	// releated field: MiscHelpGuildMember
	// releated field: MultiMonsterPrizeCount
	// releated field: MysticAlly
	// releated field: RealmAssemblyArrived
	// releated field: RealmAssemblyBackTimeout
	// releated field: RealmAssemblyCanceled
	// releated field: RealmAssemblyMemberCountNotEnough
	// releated field: RealmAssemblyMemberDestroy
	// releated field: RealmAssemblyStart
	// releated field: RealmAssemblyTargetDestroy
	// releated field: RealmAssistArrivedBack
	// releated field: RealmAssistArrivedStay
	// releated field: RealmAssistDefFail
	// releated field: RealmAssistDefSuccess
	// releated field: RealmAssistFail
	// releated field: RealmAssistSuccess
	// releated field: RealmExpelledFail
	// releated field: RealmExpelledSuccess
	// releated field: RealmInvadeFail
	// releated field: RealmInvadeFinished
	// releated field: RealmInvadeSuccess
	// releated field: RealmInvadeSuccessBack
	// releated field: RealmInvadeTargetLost
	// releated field: RealmMonsterInvasion
	// releated field: RealmTroopBacked
	// releated field: RealmTroopBaozRepatriate
	// releated field: RealmTroopRecall
	// releated field: RealmTroopRepatriate
	// releated field: RealmTroopSpeedUp
	// releated field: SysChatFriendAdded
	// releated field: SysChatSenderName
	// releated field: XiongNuResistBroadcast

	return dAtA, nil
}

var vAlIdAtOrTextHelp = map[string]*config.Validator{

	"banner_base_fail":                       config.ParseValidator("string", "", false, nil, []string{"BannerBaseFail"}),
	"banner_base_success":                    config.ParseValidator("string", "", false, nil, []string{"BannerBaseSuccess"}),
	"banner_troop_fail":                      config.ParseValidator("string", "", false, nil, []string{"BannerTroopFail"}),
	"banner_troop_success":                   config.ParseValidator("string", "", false, nil, []string{"BannerTroopSuccess"}),
	"guild_gongfang_build_complete":          config.ParseValidator("string", "", false, nil, []string{"GuildGongfangBuildComplete"}),
	"guild_gongfang_build_time_add":          config.ParseValidator("string", "", false, nil, []string{"GuildGongfangBuildTimeAdd"}),
	"guild_gongfang_build_time_reduce":       config.ParseValidator("string", "", false, nil, []string{"GuildGongfangBuildTimeReduce"}),
	"guild_gongfang_building":                config.ParseValidator("string", "", false, nil, []string{"GuildGongfangBuilding"}),
	"guild_gongfang_efficiency_add":          config.ParseValidator("string", "", false, nil, []string{"GuildGongfangEfficiencyAdd"}),
	"guild_gongfang_efficiency_reduce":       config.ParseValidator("string", "", false, nil, []string{"GuildGongfangEfficiencyReduce"}),
	"guild_gongfang_prize_send":              config.ParseValidator("string", "", false, nil, []string{"GuildGongfangPrizeSend"}),
	"guild_please_help_me":                   config.ParseValidator("string", "", false, nil, []string{"GuildPleaseHelpMe"}),
	"guild_say_hi":                           config.ParseValidator("string", "", false, nil, []string{"GuildSayHi"}),
	"guild_workshop_build_show":              config.ParseValidator("string", "", false, nil, []string{"GuildWorkshopBuildShow"}),
	"guild_workshop_completed_chat":          config.ParseValidator("string", "", false, nil, []string{"GuildWorkshopCompletedChat"}),
	"guild_workshop_created_chat":            config.ParseValidator("string", "", false, nil, []string{"GuildWorkshopCreatedChat"}),
	"guild_workshop_prize_show":              config.ParseValidator("string", "", false, nil, []string{"GuildWorkshopPrizeShow"}),
	"guild_workshop_prod_show":               config.ParseValidator("string", "", false, nil, []string{"GuildWorkshopProdShow"}),
	"mrdratt_broken4a":                       config.ParseValidator("string", "", false, nil, []string{"MRDRAttBroken4a"}),
	"mrdratt_broken4d":                       config.ParseValidator("string", "", false, nil, []string{"MRDRAttBroken4d"}),
	"mrdratt_move_base4a":                    config.ParseValidator("string", "", false, nil, []string{"MRDRAttMoveBase4a"}),
	"mrdratt_move_base4d":                    config.ParseValidator("string", "", false, nil, []string{"MRDRAttMoveBase4d"}),
	"mrdrbroken4a":                           config.ParseValidator("string", "", false, nil, []string{"MRDRBroken4a"}),
	"mrdrbroken4d":                           config.ParseValidator("string", "", false, nil, []string{"MRDRBroken4d"}),
	"mrdrexpelled4a":                         config.ParseValidator("string", "", false, nil, []string{"MRDRExpelled4a"}),
	"mrdrexpelled4d":                         config.ParseValidator("string", "", false, nil, []string{"MRDRExpelled4d"}),
	"mrdrfinished4a":                         config.ParseValidator("string", "", false, nil, []string{"MRDRFinished4a"}),
	"mrdrfinished4d":                         config.ParseValidator("string", "", false, nil, []string{"MRDRFinished4d"}),
	"mrdrmian4a":                             config.ParseValidator("string", "", false, nil, []string{"MRDRMian4a"}),
	"mrdrmian4d":                             config.ParseValidator("string", "", false, nil, []string{"MRDRMian4d"}),
	"mrdrmove_base4a":                        config.ParseValidator("string", "", false, nil, []string{"MRDRMoveBase4a"}),
	"mrdrmove_base4d":                        config.ParseValidator("string", "", false, nil, []string{"MRDRMoveBase4d"}),
	"mrdrprotected4a":                        config.ParseValidator("string", "", false, nil, []string{"MRDRProtected4a"}),
	"mrdrprotected4d":                        config.ParseValidator("string", "", false, nil, []string{"MRDRProtected4d"}),
	"mrdrrecall4a":                           config.ParseValidator("string", "", false, nil, []string{"MRDRRecall4a"}),
	"mrdrrecall4d":                           config.ParseValidator("string", "", false, nil, []string{"MRDRRecall4d"}),
	"mail_add_hate_desc":                     config.ParseValidator("string", "", false, nil, []string{"MailAddHateDesc"}),
	"mail_lost_prosperity_desc":              config.ParseValidator("string", "", false, nil, []string{"MailLostProsperityDesc"}),
	"mail_reduce_prosperity_desc":            config.ParseValidator("string", "", false, nil, []string{"MailReduceProsperityDesc"}),
	"mail_report_loser_desc":                 config.ParseValidator("string", "", false, nil, []string{"MailReportLoserDesc"}),
	"mail_report_winner_desc":                config.ParseValidator("string", "", false, nil, []string{"MailReportWinnerDesc"}),
	"mc_war_atk_fail":                        config.ParseValidator("string", "", false, nil, []string{"McWarAtkFail"}),
	"mc_war_chat_enemy_atk_building_destroy": config.ParseValidator("string", "", false, nil, []string{"McWarChatEnemyAtkBuildingDestroy"}),
	"mc_war_chat_enemy_def_building_destroy": config.ParseValidator("string", "", false, nil, []string{"McWarChatEnemyDefBuildingDestroy"}),
	"mc_war_chat_our_atk_building_destroy":   config.ParseValidator("string", "", false, nil, []string{"McWarChatOurAtkBuildingDestroy"}),
	"mc_war_chat_our_def_building_destroy":   config.ParseValidator("string", "", false, nil, []string{"McWarChatOurDefBuildingDestroy"}),
	"mc_war_def_fail":                        config.ParseValidator("string", "", false, nil, []string{"McWarDefFail"}),
	"mc_war_relive_succ":                     config.ParseValidator("string", "", false, nil, []string{"McWarReliveSucc"}),
	"misc_help_guild_member":                 config.ParseValidator("string", "", false, nil, []string{"MiscHelpGuildMember"}),
	"multi_monster_prize_count":              config.ParseValidator("string", "", false, nil, []string{"MultiMonsterPrizeCount"}),
	"mystic_ally":                            config.ParseValidator("string", "", false, nil, []string{"MysticAlly"}),
	"realm_assembly_arrived":                 config.ParseValidator("string", "", false, nil, []string{"RealmAssemblyArrived"}),
	"realm_assembly_back_timeout":            config.ParseValidator("string", "", false, nil, []string{"RealmAssemblyBackTimeout"}),
	"realm_assembly_canceled":                config.ParseValidator("string", "", false, nil, []string{"RealmAssemblyCanceled"}),
	"realm_assembly_member_count_not_enough": config.ParseValidator("string", "", false, nil, []string{"RealmAssemblyMemberCountNotEnough"}),
	"realm_assembly_member_destroy":          config.ParseValidator("string", "", false, nil, []string{"RealmAssemblyMemberDestroy"}),
	"realm_assembly_start":                   config.ParseValidator("string", "", false, nil, []string{"RealmAssemblyStart"}),
	"realm_assembly_target_destroy":          config.ParseValidator("string", "", false, nil, []string{"RealmAssemblyTargetDestroy"}),
	"realm_assist_arrived_back":              config.ParseValidator("string", "", false, nil, []string{"RealmAssistArrivedBack"}),
	"realm_assist_arrived_stay":              config.ParseValidator("string", "", false, nil, []string{"RealmAssistArrivedStay"}),
	"realm_assist_def_fail":                  config.ParseValidator("string", "", false, nil, []string{"RealmAssistDefFail"}),
	"realm_assist_def_success":               config.ParseValidator("string", "", false, nil, []string{"RealmAssistDefSuccess"}),
	"realm_assist_fail":                      config.ParseValidator("string", "", false, nil, []string{"RealmAssistFail"}),
	"realm_assist_success":                   config.ParseValidator("string", "", false, nil, []string{"RealmAssistSuccess"}),
	"realm_expelled_fail":                    config.ParseValidator("string", "", false, nil, []string{"RealmExpelledFail"}),
	"realm_expelled_success":                 config.ParseValidator("string", "", false, nil, []string{"RealmExpelledSuccess"}),
	"realm_invade_fail":                      config.ParseValidator("string", "", false, nil, []string{"RealmInvadeFail"}),
	"realm_invade_finished":                  config.ParseValidator("string", "", false, nil, []string{"RealmInvadeFinished"}),
	"realm_invade_success":                   config.ParseValidator("string", "", false, nil, []string{"RealmInvadeSuccess"}),
	"realm_invade_success_back":              config.ParseValidator("string", "", false, nil, []string{"RealmInvadeSuccessBack"}),
	"realm_invade_target_lost":               config.ParseValidator("string", "", false, nil, []string{"RealmInvadeTargetLost"}),
	"realm_monster_invasion":                 config.ParseValidator("string", "", false, nil, []string{"RealmMonsterInvasion"}),
	"realm_troop_backed":                     config.ParseValidator("string", "", false, nil, []string{"RealmTroopBacked"}),
	"realm_troop_baoz_repatriate":            config.ParseValidator("string", "", false, nil, []string{"RealmTroopBaozRepatriate"}),
	"realm_troop_recall":                     config.ParseValidator("string", "", false, nil, []string{"RealmTroopRecall"}),
	"realm_troop_repatriate":                 config.ParseValidator("string", "", false, nil, []string{"RealmTroopRepatriate"}),
	"realm_troop_speed_up":                   config.ParseValidator("string", "", false, nil, []string{"RealmTroopSpeedUp"}),
	"sys_chat_friend_added":                  config.ParseValidator("string", "", false, nil, []string{"SysChatFriendAdded"}),
	"sys_chat_sender_name":                   config.ParseValidator("string", "", false, nil, []string{"SysChatSenderName"}),
	"xiong_nu_resist_broadcast":              config.ParseValidator("string", "", false, nil, []string{"XiongNuResistBroadcast"}),
}

func (dAtA *TextHelp) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("banner_base_fail") {
		dAtA.BannerBaseFail = cOnFigS.GetText(pArSeR.String("banner_base_fail"))
	} else {
		dAtA.BannerBaseFail = cOnFigS.GetText("BannerBaseFail")
	}
	if dAtA.BannerBaseFail == nil {
		return errors.Errorf("%s 配置的关联字段[banner_base_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("banner_base_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("banner_base_success") {
		dAtA.BannerBaseSuccess = cOnFigS.GetText(pArSeR.String("banner_base_success"))
	} else {
		dAtA.BannerBaseSuccess = cOnFigS.GetText("BannerBaseSuccess")
	}
	if dAtA.BannerBaseSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[banner_base_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("banner_base_success"), *pArSeR)
	}

	if pArSeR.KeyExist("banner_troop_fail") {
		dAtA.BannerTroopFail = cOnFigS.GetText(pArSeR.String("banner_troop_fail"))
	} else {
		dAtA.BannerTroopFail = cOnFigS.GetText("BannerTroopFail")
	}
	if dAtA.BannerTroopFail == nil {
		return errors.Errorf("%s 配置的关联字段[banner_troop_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("banner_troop_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("banner_troop_success") {
		dAtA.BannerTroopSuccess = cOnFigS.GetText(pArSeR.String("banner_troop_success"))
	} else {
		dAtA.BannerTroopSuccess = cOnFigS.GetText("BannerTroopSuccess")
	}
	if dAtA.BannerTroopSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[banner_troop_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("banner_troop_success"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_gongfang_build_complete") {
		dAtA.GuildGongfangBuildComplete = cOnFigS.GetText(pArSeR.String("guild_gongfang_build_complete"))
	} else {
		dAtA.GuildGongfangBuildComplete = cOnFigS.GetText("GuildGongfangBuildComplete")
	}
	if dAtA.GuildGongfangBuildComplete == nil {
		return errors.Errorf("%s 配置的关联字段[guild_gongfang_build_complete] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_gongfang_build_complete"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_gongfang_build_time_add") {
		dAtA.GuildGongfangBuildTimeAdd = cOnFigS.GetText(pArSeR.String("guild_gongfang_build_time_add"))
	} else {
		dAtA.GuildGongfangBuildTimeAdd = cOnFigS.GetText("GuildGongfangBuildTimeAdd")
	}
	if dAtA.GuildGongfangBuildTimeAdd == nil {
		return errors.Errorf("%s 配置的关联字段[guild_gongfang_build_time_add] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_gongfang_build_time_add"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_gongfang_build_time_reduce") {
		dAtA.GuildGongfangBuildTimeReduce = cOnFigS.GetText(pArSeR.String("guild_gongfang_build_time_reduce"))
	} else {
		dAtA.GuildGongfangBuildTimeReduce = cOnFigS.GetText("GuildGongfangBuildTimeReduce")
	}
	if dAtA.GuildGongfangBuildTimeReduce == nil {
		return errors.Errorf("%s 配置的关联字段[guild_gongfang_build_time_reduce] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_gongfang_build_time_reduce"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_gongfang_building") {
		dAtA.GuildGongfangBuilding = cOnFigS.GetText(pArSeR.String("guild_gongfang_building"))
	} else {
		dAtA.GuildGongfangBuilding = cOnFigS.GetText("GuildGongfangBuilding")
	}
	if dAtA.GuildGongfangBuilding == nil {
		return errors.Errorf("%s 配置的关联字段[guild_gongfang_building] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_gongfang_building"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_gongfang_efficiency_add") {
		dAtA.GuildGongfangEfficiencyAdd = cOnFigS.GetText(pArSeR.String("guild_gongfang_efficiency_add"))
	} else {
		dAtA.GuildGongfangEfficiencyAdd = cOnFigS.GetText("GuildGongfangEfficiencyAdd")
	}
	if dAtA.GuildGongfangEfficiencyAdd == nil {
		return errors.Errorf("%s 配置的关联字段[guild_gongfang_efficiency_add] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_gongfang_efficiency_add"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_gongfang_efficiency_reduce") {
		dAtA.GuildGongfangEfficiencyReduce = cOnFigS.GetText(pArSeR.String("guild_gongfang_efficiency_reduce"))
	} else {
		dAtA.GuildGongfangEfficiencyReduce = cOnFigS.GetText("GuildGongfangEfficiencyReduce")
	}
	if dAtA.GuildGongfangEfficiencyReduce == nil {
		return errors.Errorf("%s 配置的关联字段[guild_gongfang_efficiency_reduce] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_gongfang_efficiency_reduce"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_gongfang_prize_send") {
		dAtA.GuildGongfangPrizeSend = cOnFigS.GetText(pArSeR.String("guild_gongfang_prize_send"))
	} else {
		dAtA.GuildGongfangPrizeSend = cOnFigS.GetText("GuildGongfangPrizeSend")
	}
	if dAtA.GuildGongfangPrizeSend == nil {
		return errors.Errorf("%s 配置的关联字段[guild_gongfang_prize_send] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_gongfang_prize_send"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_please_help_me") {
		dAtA.GuildPleaseHelpMe = cOnFigS.GetText(pArSeR.String("guild_please_help_me"))
	} else {
		dAtA.GuildPleaseHelpMe = cOnFigS.GetText("GuildPleaseHelpMe")
	}
	if dAtA.GuildPleaseHelpMe == nil {
		return errors.Errorf("%s 配置的关联字段[guild_please_help_me] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_please_help_me"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_say_hi") {
		dAtA.GuildSayHi = cOnFigS.GetText(pArSeR.String("guild_say_hi"))
	} else {
		dAtA.GuildSayHi = cOnFigS.GetText("GuildSayHi")
	}
	if dAtA.GuildSayHi == nil {
		return errors.Errorf("%s 配置的关联字段[guild_say_hi] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_say_hi"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_workshop_build_show") {
		dAtA.GuildWorkshopBuildShow = cOnFigS.GetText(pArSeR.String("guild_workshop_build_show"))
	} else {
		dAtA.GuildWorkshopBuildShow = cOnFigS.GetText("GuildWorkshopBuildShow")
	}
	if dAtA.GuildWorkshopBuildShow == nil {
		return errors.Errorf("%s 配置的关联字段[guild_workshop_build_show] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_workshop_build_show"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_workshop_completed_chat") {
		dAtA.GuildWorkshopCompletedChat = cOnFigS.GetText(pArSeR.String("guild_workshop_completed_chat"))
	} else {
		dAtA.GuildWorkshopCompletedChat = cOnFigS.GetText("GuildWorkshopCompletedChat")
	}
	if dAtA.GuildWorkshopCompletedChat == nil {
		return errors.Errorf("%s 配置的关联字段[guild_workshop_completed_chat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_workshop_completed_chat"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_workshop_created_chat") {
		dAtA.GuildWorkshopCreatedChat = cOnFigS.GetText(pArSeR.String("guild_workshop_created_chat"))
	} else {
		dAtA.GuildWorkshopCreatedChat = cOnFigS.GetText("GuildWorkshopCreatedChat")
	}
	if dAtA.GuildWorkshopCreatedChat == nil {
		return errors.Errorf("%s 配置的关联字段[guild_workshop_created_chat] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_workshop_created_chat"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_workshop_prize_show") {
		dAtA.GuildWorkshopPrizeShow = cOnFigS.GetText(pArSeR.String("guild_workshop_prize_show"))
	} else {
		dAtA.GuildWorkshopPrizeShow = cOnFigS.GetText("GuildWorkshopPrizeShow")
	}
	if dAtA.GuildWorkshopPrizeShow == nil {
		return errors.Errorf("%s 配置的关联字段[guild_workshop_prize_show] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_workshop_prize_show"), *pArSeR)
	}

	if pArSeR.KeyExist("guild_workshop_prod_show") {
		dAtA.GuildWorkshopProdShow = cOnFigS.GetText(pArSeR.String("guild_workshop_prod_show"))
	} else {
		dAtA.GuildWorkshopProdShow = cOnFigS.GetText("GuildWorkshopProdShow")
	}
	if dAtA.GuildWorkshopProdShow == nil {
		return errors.Errorf("%s 配置的关联字段[guild_workshop_prod_show] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guild_workshop_prod_show"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdratt_broken4a") {
		dAtA.MRDRAttBroken4a = cOnFigS.GetText(pArSeR.String("mrdratt_broken4a"))
	} else {
		dAtA.MRDRAttBroken4a = cOnFigS.GetText("MRDRAttBroken4a")
	}
	if dAtA.MRDRAttBroken4a == nil {
		return errors.Errorf("%s 配置的关联字段[mrdratt_broken4a] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdratt_broken4a"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdratt_broken4d") {
		dAtA.MRDRAttBroken4d = cOnFigS.GetText(pArSeR.String("mrdratt_broken4d"))
	} else {
		dAtA.MRDRAttBroken4d = cOnFigS.GetText("MRDRAttBroken4d")
	}
	if dAtA.MRDRAttBroken4d == nil {
		return errors.Errorf("%s 配置的关联字段[mrdratt_broken4d] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdratt_broken4d"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdratt_move_base4a") {
		dAtA.MRDRAttMoveBase4a = cOnFigS.GetText(pArSeR.String("mrdratt_move_base4a"))
	} else {
		dAtA.MRDRAttMoveBase4a = cOnFigS.GetText("MRDRAttMoveBase4a")
	}
	if dAtA.MRDRAttMoveBase4a == nil {
		return errors.Errorf("%s 配置的关联字段[mrdratt_move_base4a] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdratt_move_base4a"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdratt_move_base4d") {
		dAtA.MRDRAttMoveBase4d = cOnFigS.GetText(pArSeR.String("mrdratt_move_base4d"))
	} else {
		dAtA.MRDRAttMoveBase4d = cOnFigS.GetText("MRDRAttMoveBase4d")
	}
	if dAtA.MRDRAttMoveBase4d == nil {
		return errors.Errorf("%s 配置的关联字段[mrdratt_move_base4d] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdratt_move_base4d"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrbroken4a") {
		dAtA.MRDRBroken4a = cOnFigS.GetText(pArSeR.String("mrdrbroken4a"))
	} else {
		dAtA.MRDRBroken4a = cOnFigS.GetText("MRDRBroken4a")
	}
	if dAtA.MRDRBroken4a == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrbroken4a] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrbroken4a"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrbroken4d") {
		dAtA.MRDRBroken4d = cOnFigS.GetText(pArSeR.String("mrdrbroken4d"))
	} else {
		dAtA.MRDRBroken4d = cOnFigS.GetText("MRDRBroken4d")
	}
	if dAtA.MRDRBroken4d == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrbroken4d] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrbroken4d"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrexpelled4a") {
		dAtA.MRDRExpelled4a = cOnFigS.GetText(pArSeR.String("mrdrexpelled4a"))
	} else {
		dAtA.MRDRExpelled4a = cOnFigS.GetText("MRDRExpelled4a")
	}
	if dAtA.MRDRExpelled4a == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrexpelled4a] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrexpelled4a"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrexpelled4d") {
		dAtA.MRDRExpelled4d = cOnFigS.GetText(pArSeR.String("mrdrexpelled4d"))
	} else {
		dAtA.MRDRExpelled4d = cOnFigS.GetText("MRDRExpelled4d")
	}
	if dAtA.MRDRExpelled4d == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrexpelled4d] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrexpelled4d"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrfinished4a") {
		dAtA.MRDRFinished4a = cOnFigS.GetText(pArSeR.String("mrdrfinished4a"))
	} else {
		dAtA.MRDRFinished4a = cOnFigS.GetText("MRDRFinished4a")
	}
	if dAtA.MRDRFinished4a == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrfinished4a] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrfinished4a"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrfinished4d") {
		dAtA.MRDRFinished4d = cOnFigS.GetText(pArSeR.String("mrdrfinished4d"))
	} else {
		dAtA.MRDRFinished4d = cOnFigS.GetText("MRDRFinished4d")
	}
	if dAtA.MRDRFinished4d == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrfinished4d] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrfinished4d"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrmian4a") {
		dAtA.MRDRMian4a = cOnFigS.GetText(pArSeR.String("mrdrmian4a"))
	} else {
		dAtA.MRDRMian4a = cOnFigS.GetText("MRDRMian4a")
	}
	if dAtA.MRDRMian4a == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrmian4a] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrmian4a"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrmian4d") {
		dAtA.MRDRMian4d = cOnFigS.GetText(pArSeR.String("mrdrmian4d"))
	} else {
		dAtA.MRDRMian4d = cOnFigS.GetText("MRDRMian4d")
	}
	if dAtA.MRDRMian4d == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrmian4d] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrmian4d"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrmove_base4a") {
		dAtA.MRDRMoveBase4a = cOnFigS.GetText(pArSeR.String("mrdrmove_base4a"))
	} else {
		dAtA.MRDRMoveBase4a = cOnFigS.GetText("MRDRMoveBase4a")
	}
	if dAtA.MRDRMoveBase4a == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrmove_base4a] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrmove_base4a"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrmove_base4d") {
		dAtA.MRDRMoveBase4d = cOnFigS.GetText(pArSeR.String("mrdrmove_base4d"))
	} else {
		dAtA.MRDRMoveBase4d = cOnFigS.GetText("MRDRMoveBase4d")
	}
	if dAtA.MRDRMoveBase4d == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrmove_base4d] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrmove_base4d"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrprotected4a") {
		dAtA.MRDRProtected4a = cOnFigS.GetText(pArSeR.String("mrdrprotected4a"))
	} else {
		dAtA.MRDRProtected4a = cOnFigS.GetText("MRDRProtected4a")
	}
	if dAtA.MRDRProtected4a == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrprotected4a] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrprotected4a"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrprotected4d") {
		dAtA.MRDRProtected4d = cOnFigS.GetText(pArSeR.String("mrdrprotected4d"))
	} else {
		dAtA.MRDRProtected4d = cOnFigS.GetText("MRDRProtected4d")
	}
	if dAtA.MRDRProtected4d == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrprotected4d] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrprotected4d"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrrecall4a") {
		dAtA.MRDRRecall4a = cOnFigS.GetText(pArSeR.String("mrdrrecall4a"))
	} else {
		dAtA.MRDRRecall4a = cOnFigS.GetText("MRDRRecall4a")
	}
	if dAtA.MRDRRecall4a == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrrecall4a] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrrecall4a"), *pArSeR)
	}

	if pArSeR.KeyExist("mrdrrecall4d") {
		dAtA.MRDRRecall4d = cOnFigS.GetText(pArSeR.String("mrdrrecall4d"))
	} else {
		dAtA.MRDRRecall4d = cOnFigS.GetText("MRDRRecall4d")
	}
	if dAtA.MRDRRecall4d == nil {
		return errors.Errorf("%s 配置的关联字段[mrdrrecall4d] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mrdrrecall4d"), *pArSeR)
	}

	if pArSeR.KeyExist("mail_add_hate_desc") {
		dAtA.MailAddHateDesc = cOnFigS.GetText(pArSeR.String("mail_add_hate_desc"))
	} else {
		dAtA.MailAddHateDesc = cOnFigS.GetText("MailAddHateDesc")
	}
	if dAtA.MailAddHateDesc == nil {
		return errors.Errorf("%s 配置的关联字段[mail_add_hate_desc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mail_add_hate_desc"), *pArSeR)
	}

	if pArSeR.KeyExist("mail_lost_prosperity_desc") {
		dAtA.MailLostProsperityDesc = cOnFigS.GetText(pArSeR.String("mail_lost_prosperity_desc"))
	} else {
		dAtA.MailLostProsperityDesc = cOnFigS.GetText("MailLostProsperityDesc")
	}
	if dAtA.MailLostProsperityDesc == nil {
		return errors.Errorf("%s 配置的关联字段[mail_lost_prosperity_desc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mail_lost_prosperity_desc"), *pArSeR)
	}

	if pArSeR.KeyExist("mail_reduce_prosperity_desc") {
		dAtA.MailReduceProsperityDesc = cOnFigS.GetText(pArSeR.String("mail_reduce_prosperity_desc"))
	} else {
		dAtA.MailReduceProsperityDesc = cOnFigS.GetText("MailReduceProsperityDesc")
	}
	if dAtA.MailReduceProsperityDesc == nil {
		return errors.Errorf("%s 配置的关联字段[mail_reduce_prosperity_desc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mail_reduce_prosperity_desc"), *pArSeR)
	}

	if pArSeR.KeyExist("mail_report_loser_desc") {
		dAtA.MailReportLoserDesc = cOnFigS.GetText(pArSeR.String("mail_report_loser_desc"))
	} else {
		dAtA.MailReportLoserDesc = cOnFigS.GetText("MailReportLoserDesc")
	}
	if dAtA.MailReportLoserDesc == nil {
		return errors.Errorf("%s 配置的关联字段[mail_report_loser_desc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mail_report_loser_desc"), *pArSeR)
	}

	if pArSeR.KeyExist("mail_report_winner_desc") {
		dAtA.MailReportWinnerDesc = cOnFigS.GetText(pArSeR.String("mail_report_winner_desc"))
	} else {
		dAtA.MailReportWinnerDesc = cOnFigS.GetText("MailReportWinnerDesc")
	}
	if dAtA.MailReportWinnerDesc == nil {
		return errors.Errorf("%s 配置的关联字段[mail_report_winner_desc] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mail_report_winner_desc"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_atk_fail") {
		dAtA.McWarAtkFail = cOnFigS.GetText(pArSeR.String("mc_war_atk_fail"))
	} else {
		dAtA.McWarAtkFail = cOnFigS.GetText("McWarAtkFail")
	}
	if dAtA.McWarAtkFail == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_atk_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_atk_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_chat_enemy_atk_building_destroy") {
		dAtA.McWarChatEnemyAtkBuildingDestroy = cOnFigS.GetText(pArSeR.String("mc_war_chat_enemy_atk_building_destroy"))
	} else {
		dAtA.McWarChatEnemyAtkBuildingDestroy = cOnFigS.GetText("McWarChatEnemyAtkBuildingDestroy")
	}
	if dAtA.McWarChatEnemyAtkBuildingDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_chat_enemy_atk_building_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_chat_enemy_atk_building_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_chat_enemy_def_building_destroy") {
		dAtA.McWarChatEnemyDefBuildingDestroy = cOnFigS.GetText(pArSeR.String("mc_war_chat_enemy_def_building_destroy"))
	} else {
		dAtA.McWarChatEnemyDefBuildingDestroy = cOnFigS.GetText("McWarChatEnemyDefBuildingDestroy")
	}
	if dAtA.McWarChatEnemyDefBuildingDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_chat_enemy_def_building_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_chat_enemy_def_building_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_chat_our_atk_building_destroy") {
		dAtA.McWarChatOurAtkBuildingDestroy = cOnFigS.GetText(pArSeR.String("mc_war_chat_our_atk_building_destroy"))
	} else {
		dAtA.McWarChatOurAtkBuildingDestroy = cOnFigS.GetText("McWarChatOurAtkBuildingDestroy")
	}
	if dAtA.McWarChatOurAtkBuildingDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_chat_our_atk_building_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_chat_our_atk_building_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_chat_our_def_building_destroy") {
		dAtA.McWarChatOurDefBuildingDestroy = cOnFigS.GetText(pArSeR.String("mc_war_chat_our_def_building_destroy"))
	} else {
		dAtA.McWarChatOurDefBuildingDestroy = cOnFigS.GetText("McWarChatOurDefBuildingDestroy")
	}
	if dAtA.McWarChatOurDefBuildingDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_chat_our_def_building_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_chat_our_def_building_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_def_fail") {
		dAtA.McWarDefFail = cOnFigS.GetText(pArSeR.String("mc_war_def_fail"))
	} else {
		dAtA.McWarDefFail = cOnFigS.GetText("McWarDefFail")
	}
	if dAtA.McWarDefFail == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_def_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_def_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_relive_succ") {
		dAtA.McWarReliveSucc = cOnFigS.GetText(pArSeR.String("mc_war_relive_succ"))
	} else {
		dAtA.McWarReliveSucc = cOnFigS.GetText("McWarReliveSucc")
	}
	if dAtA.McWarReliveSucc == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_relive_succ] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_relive_succ"), *pArSeR)
	}

	if pArSeR.KeyExist("misc_help_guild_member") {
		dAtA.MiscHelpGuildMember = cOnFigS.GetText(pArSeR.String("misc_help_guild_member"))
	} else {
		dAtA.MiscHelpGuildMember = cOnFigS.GetText("MiscHelpGuildMember")
	}
	if dAtA.MiscHelpGuildMember == nil {
		return errors.Errorf("%s 配置的关联字段[misc_help_guild_member] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("misc_help_guild_member"), *pArSeR)
	}

	if pArSeR.KeyExist("multi_monster_prize_count") {
		dAtA.MultiMonsterPrizeCount = cOnFigS.GetText(pArSeR.String("multi_monster_prize_count"))
	} else {
		dAtA.MultiMonsterPrizeCount = cOnFigS.GetText("MultiMonsterPrizeCount")
	}
	if dAtA.MultiMonsterPrizeCount == nil {
		return errors.Errorf("%s 配置的关联字段[multi_monster_prize_count] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("multi_monster_prize_count"), *pArSeR)
	}

	if pArSeR.KeyExist("mystic_ally") {
		dAtA.MysticAlly = cOnFigS.GetText(pArSeR.String("mystic_ally"))
	} else {
		dAtA.MysticAlly = cOnFigS.GetText("MysticAlly")
	}
	if dAtA.MysticAlly == nil {
		return errors.Errorf("%s 配置的关联字段[mystic_ally] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mystic_ally"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assembly_arrived") {
		dAtA.RealmAssemblyArrived = cOnFigS.GetText(pArSeR.String("realm_assembly_arrived"))
	} else {
		dAtA.RealmAssemblyArrived = cOnFigS.GetText("RealmAssemblyArrived")
	}
	if dAtA.RealmAssemblyArrived == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assembly_arrived] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assembly_arrived"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assembly_back_timeout") {
		dAtA.RealmAssemblyBackTimeout = cOnFigS.GetText(pArSeR.String("realm_assembly_back_timeout"))
	} else {
		dAtA.RealmAssemblyBackTimeout = cOnFigS.GetText("RealmAssemblyBackTimeout")
	}
	if dAtA.RealmAssemblyBackTimeout == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assembly_back_timeout] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assembly_back_timeout"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assembly_canceled") {
		dAtA.RealmAssemblyCanceled = cOnFigS.GetText(pArSeR.String("realm_assembly_canceled"))
	} else {
		dAtA.RealmAssemblyCanceled = cOnFigS.GetText("RealmAssemblyCanceled")
	}
	if dAtA.RealmAssemblyCanceled == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assembly_canceled] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assembly_canceled"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assembly_member_count_not_enough") {
		dAtA.RealmAssemblyMemberCountNotEnough = cOnFigS.GetText(pArSeR.String("realm_assembly_member_count_not_enough"))
	} else {
		dAtA.RealmAssemblyMemberCountNotEnough = cOnFigS.GetText("RealmAssemblyMemberCountNotEnough")
	}
	if dAtA.RealmAssemblyMemberCountNotEnough == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assembly_member_count_not_enough] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assembly_member_count_not_enough"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assembly_member_destroy") {
		dAtA.RealmAssemblyMemberDestroy = cOnFigS.GetText(pArSeR.String("realm_assembly_member_destroy"))
	} else {
		dAtA.RealmAssemblyMemberDestroy = cOnFigS.GetText("RealmAssemblyMemberDestroy")
	}
	if dAtA.RealmAssemblyMemberDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assembly_member_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assembly_member_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assembly_start") {
		dAtA.RealmAssemblyStart = cOnFigS.GetText(pArSeR.String("realm_assembly_start"))
	} else {
		dAtA.RealmAssemblyStart = cOnFigS.GetText("RealmAssemblyStart")
	}
	if dAtA.RealmAssemblyStart == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assembly_start] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assembly_start"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assembly_target_destroy") {
		dAtA.RealmAssemblyTargetDestroy = cOnFigS.GetText(pArSeR.String("realm_assembly_target_destroy"))
	} else {
		dAtA.RealmAssemblyTargetDestroy = cOnFigS.GetText("RealmAssemblyTargetDestroy")
	}
	if dAtA.RealmAssemblyTargetDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assembly_target_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assembly_target_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assist_arrived_back") {
		dAtA.RealmAssistArrivedBack = cOnFigS.GetText(pArSeR.String("realm_assist_arrived_back"))
	} else {
		dAtA.RealmAssistArrivedBack = cOnFigS.GetText("RealmAssistArrivedBack")
	}
	if dAtA.RealmAssistArrivedBack == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assist_arrived_back] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assist_arrived_back"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assist_arrived_stay") {
		dAtA.RealmAssistArrivedStay = cOnFigS.GetText(pArSeR.String("realm_assist_arrived_stay"))
	} else {
		dAtA.RealmAssistArrivedStay = cOnFigS.GetText("RealmAssistArrivedStay")
	}
	if dAtA.RealmAssistArrivedStay == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assist_arrived_stay] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assist_arrived_stay"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assist_def_fail") {
		dAtA.RealmAssistDefFail = cOnFigS.GetText(pArSeR.String("realm_assist_def_fail"))
	} else {
		dAtA.RealmAssistDefFail = cOnFigS.GetText("RealmAssistDefFail")
	}
	if dAtA.RealmAssistDefFail == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assist_def_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assist_def_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assist_def_success") {
		dAtA.RealmAssistDefSuccess = cOnFigS.GetText(pArSeR.String("realm_assist_def_success"))
	} else {
		dAtA.RealmAssistDefSuccess = cOnFigS.GetText("RealmAssistDefSuccess")
	}
	if dAtA.RealmAssistDefSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assist_def_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assist_def_success"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assist_fail") {
		dAtA.RealmAssistFail = cOnFigS.GetText(pArSeR.String("realm_assist_fail"))
	} else {
		dAtA.RealmAssistFail = cOnFigS.GetText("RealmAssistFail")
	}
	if dAtA.RealmAssistFail == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assist_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assist_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_assist_success") {
		dAtA.RealmAssistSuccess = cOnFigS.GetText(pArSeR.String("realm_assist_success"))
	} else {
		dAtA.RealmAssistSuccess = cOnFigS.GetText("RealmAssistSuccess")
	}
	if dAtA.RealmAssistSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[realm_assist_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_assist_success"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_expelled_fail") {
		dAtA.RealmExpelledFail = cOnFigS.GetText(pArSeR.String("realm_expelled_fail"))
	} else {
		dAtA.RealmExpelledFail = cOnFigS.GetText("RealmExpelledFail")
	}
	if dAtA.RealmExpelledFail == nil {
		return errors.Errorf("%s 配置的关联字段[realm_expelled_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_expelled_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_expelled_success") {
		dAtA.RealmExpelledSuccess = cOnFigS.GetText(pArSeR.String("realm_expelled_success"))
	} else {
		dAtA.RealmExpelledSuccess = cOnFigS.GetText("RealmExpelledSuccess")
	}
	if dAtA.RealmExpelledSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[realm_expelled_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_expelled_success"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_invade_fail") {
		dAtA.RealmInvadeFail = cOnFigS.GetText(pArSeR.String("realm_invade_fail"))
	} else {
		dAtA.RealmInvadeFail = cOnFigS.GetText("RealmInvadeFail")
	}
	if dAtA.RealmInvadeFail == nil {
		return errors.Errorf("%s 配置的关联字段[realm_invade_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_invade_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_invade_finished") {
		dAtA.RealmInvadeFinished = cOnFigS.GetText(pArSeR.String("realm_invade_finished"))
	} else {
		dAtA.RealmInvadeFinished = cOnFigS.GetText("RealmInvadeFinished")
	}
	if dAtA.RealmInvadeFinished == nil {
		return errors.Errorf("%s 配置的关联字段[realm_invade_finished] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_invade_finished"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_invade_success") {
		dAtA.RealmInvadeSuccess = cOnFigS.GetText(pArSeR.String("realm_invade_success"))
	} else {
		dAtA.RealmInvadeSuccess = cOnFigS.GetText("RealmInvadeSuccess")
	}
	if dAtA.RealmInvadeSuccess == nil {
		return errors.Errorf("%s 配置的关联字段[realm_invade_success] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_invade_success"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_invade_success_back") {
		dAtA.RealmInvadeSuccessBack = cOnFigS.GetText(pArSeR.String("realm_invade_success_back"))
	} else {
		dAtA.RealmInvadeSuccessBack = cOnFigS.GetText("RealmInvadeSuccessBack")
	}
	if dAtA.RealmInvadeSuccessBack == nil {
		return errors.Errorf("%s 配置的关联字段[realm_invade_success_back] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_invade_success_back"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_invade_target_lost") {
		dAtA.RealmInvadeTargetLost = cOnFigS.GetText(pArSeR.String("realm_invade_target_lost"))
	} else {
		dAtA.RealmInvadeTargetLost = cOnFigS.GetText("RealmInvadeTargetLost")
	}
	if dAtA.RealmInvadeTargetLost == nil {
		return errors.Errorf("%s 配置的关联字段[realm_invade_target_lost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_invade_target_lost"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_monster_invasion") {
		dAtA.RealmMonsterInvasion = cOnFigS.GetText(pArSeR.String("realm_monster_invasion"))
	} else {
		dAtA.RealmMonsterInvasion = cOnFigS.GetText("RealmMonsterInvasion")
	}
	if dAtA.RealmMonsterInvasion == nil {
		return errors.Errorf("%s 配置的关联字段[realm_monster_invasion] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_monster_invasion"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_troop_backed") {
		dAtA.RealmTroopBacked = cOnFigS.GetText(pArSeR.String("realm_troop_backed"))
	} else {
		dAtA.RealmTroopBacked = cOnFigS.GetText("RealmTroopBacked")
	}
	if dAtA.RealmTroopBacked == nil {
		return errors.Errorf("%s 配置的关联字段[realm_troop_backed] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_troop_backed"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_troop_baoz_repatriate") {
		dAtA.RealmTroopBaozRepatriate = cOnFigS.GetText(pArSeR.String("realm_troop_baoz_repatriate"))
	} else {
		dAtA.RealmTroopBaozRepatriate = cOnFigS.GetText("RealmTroopBaozRepatriate")
	}
	if dAtA.RealmTroopBaozRepatriate == nil {
		return errors.Errorf("%s 配置的关联字段[realm_troop_baoz_repatriate] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_troop_baoz_repatriate"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_troop_recall") {
		dAtA.RealmTroopRecall = cOnFigS.GetText(pArSeR.String("realm_troop_recall"))
	} else {
		dAtA.RealmTroopRecall = cOnFigS.GetText("RealmTroopRecall")
	}
	if dAtA.RealmTroopRecall == nil {
		return errors.Errorf("%s 配置的关联字段[realm_troop_recall] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_troop_recall"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_troop_repatriate") {
		dAtA.RealmTroopRepatriate = cOnFigS.GetText(pArSeR.String("realm_troop_repatriate"))
	} else {
		dAtA.RealmTroopRepatriate = cOnFigS.GetText("RealmTroopRepatriate")
	}
	if dAtA.RealmTroopRepatriate == nil {
		return errors.Errorf("%s 配置的关联字段[realm_troop_repatriate] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_troop_repatriate"), *pArSeR)
	}

	if pArSeR.KeyExist("realm_troop_speed_up") {
		dAtA.RealmTroopSpeedUp = cOnFigS.GetText(pArSeR.String("realm_troop_speed_up"))
	} else {
		dAtA.RealmTroopSpeedUp = cOnFigS.GetText("RealmTroopSpeedUp")
	}
	if dAtA.RealmTroopSpeedUp == nil {
		return errors.Errorf("%s 配置的关联字段[realm_troop_speed_up] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("realm_troop_speed_up"), *pArSeR)
	}

	if pArSeR.KeyExist("sys_chat_friend_added") {
		dAtA.SysChatFriendAdded = cOnFigS.GetText(pArSeR.String("sys_chat_friend_added"))
	} else {
		dAtA.SysChatFriendAdded = cOnFigS.GetText("SysChatFriendAdded")
	}
	if dAtA.SysChatFriendAdded == nil {
		return errors.Errorf("%s 配置的关联字段[sys_chat_friend_added] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("sys_chat_friend_added"), *pArSeR)
	}

	if pArSeR.KeyExist("sys_chat_sender_name") {
		dAtA.SysChatSenderName = cOnFigS.GetText(pArSeR.String("sys_chat_sender_name"))
	} else {
		dAtA.SysChatSenderName = cOnFigS.GetText("SysChatSenderName")
	}
	if dAtA.SysChatSenderName == nil {
		return errors.Errorf("%s 配置的关联字段[sys_chat_sender_name] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("sys_chat_sender_name"), *pArSeR)
	}

	if pArSeR.KeyExist("xiong_nu_resist_broadcast") {
		dAtA.XiongNuResistBroadcast = cOnFigS.GetText(pArSeR.String("xiong_nu_resist_broadcast"))
	} else {
		dAtA.XiongNuResistBroadcast = cOnFigS.GetText("XiongNuResistBroadcast")
	}
	if dAtA.XiongNuResistBroadcast == nil {
		return errors.Errorf("%s 配置的关联字段[xiong_nu_resist_broadcast] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("xiong_nu_resist_broadcast"), *pArSeR)
	}

	return nil
}

// start with TimeRuleData ----------------------------------

func LoadTimeRuleData(gos *config.GameObjects) (map[uint64]*TimeRuleData, map[*TimeRuleData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TimeRuleDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TimeRuleData, len(lIsT))
	pArSeRmAp := make(map[*TimeRuleData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTimeRuleData) {
			continue
		}

		dAtA, err := NewTimeRuleData(fIlEnAmE, pArSeR)
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

func SetRelatedTimeRuleData(dAtAmAp map[*TimeRuleData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TimeRuleDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTimeRuleDataKeyArray(datas []*TimeRuleData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTimeRuleData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TimeRuleData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTimeRuleData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TimeRuleData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.RuleType = 0
	if pArSeR.KeyExist("rule_type") {
		dAtA.RuleType = pArSeR.Uint64("rule_type")
	}

	dAtA.Rule = pArSeR.String("rule")
	dAtA.Time = pArSeR.String("time")
	dAtA.TimeDuration, err = config.ParseDuration(pArSeR.String("time_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[time_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("time_duration"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrTimeRuleData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"rule_type":     config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"rule":          config.ParseValidator("string", "", false, nil, nil),
	"time":          config.ParseValidator("string", "", false, nil, nil),
	"time_duration": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *TimeRuleData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with UnlockCondition ----------------------------------

func NewUnlockCondition(fIlEnAmE string, pArSeR *config.ObjectParser) (*UnlockCondition, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrUnlockCondition)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &UnlockCondition{}

	dAtA.RequiredHeroLevel = 0
	if pArSeR.KeyExist("required_hero_level") {
		dAtA.RequiredHeroLevel = pArSeR.Uint64("required_hero_level")
	}

	dAtA.RequiredBaseLevel = 0
	if pArSeR.KeyExist("required_base_level") {
		dAtA.RequiredBaseLevel = pArSeR.Uint64("required_base_level")
	}

	dAtA.RequiredGuildLevel = 0
	if pArSeR.KeyExist("required_guild_level") {
		dAtA.RequiredGuildLevel = pArSeR.Uint64("required_guild_level")
	}

	dAtA.RequiredVipLevel = 0
	if pArSeR.KeyExist("required_vip_level") {
		dAtA.RequiredVipLevel = pArSeR.Uint64("required_vip_level")
	}

	return dAtA, nil
}

var vAlIdAtOrUnlockCondition = map[string]*config.Validator{

	"required_hero_level":  config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"required_base_level":  config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"required_guild_level": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"required_vip_level":   config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *UnlockCondition) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *UnlockCondition) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *UnlockCondition) Encode() *shared_proto.UnlockConditionProto {
	out := &shared_proto.UnlockConditionProto{}
	out.RequiredHeroLevel = config.U64ToI32(dAtA.RequiredHeroLevel)
	out.RequiredBaseLevel = config.U64ToI32(dAtA.RequiredBaseLevel)
	out.RequiredGuildLevel = config.U64ToI32(dAtA.RequiredGuildLevel)
	out.RequiredVipLevel = config.U64ToI32(dAtA.RequiredVipLevel)

	return out
}

func ArrayEncodeUnlockCondition(datas []*UnlockCondition) []*shared_proto.UnlockConditionProto {

	out := make([]*shared_proto.UnlockConditionProto, 0, len(datas))
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

func (dAtA *UnlockCondition) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetBroadcastData(string) *BroadcastData
	GetSpriteStat(uint64) *SpriteStat
	GetText(string) *Text
}
