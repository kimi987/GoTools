// AUTO_GEN, DONT MODIFY!!!
package dungeon

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/captain"
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

// start with DungeonChapterData ----------------------------------

func LoadDungeonChapterData(gos *config.GameObjects) (map[uint64]*DungeonChapterData, map[*DungeonChapterData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.DungeonChapterDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*DungeonChapterData, len(lIsT))
	pArSeRmAp := make(map[*DungeonChapterData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrDungeonChapterData) {
			continue
		}

		dAtA, err := NewDungeonChapterData(fIlEnAmE, pArSeR)
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

func SetRelatedDungeonChapterData(dAtAmAp map[*DungeonChapterData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.DungeonChapterDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetDungeonChapterDataKeyArray(datas []*DungeonChapterData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewDungeonChapterData(fIlEnAmE string, pArSeR *config.ObjectParser) (*DungeonChapterData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrDungeonChapterData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &DungeonChapterData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.ChapterName = pArSeR.String("chapter_name")
	dAtA.ChapterDesc = pArSeR.String("chapter_desc")
	// releated field: Captain
	dAtA.Type = shared_proto.DifficultType(shared_proto.DifficultType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.DifficultType(i)
	}

	dAtA.BgImg = pArSeR.String("bg_img")
	// skip field: DungeonDatas
	// releated field: PassPrize
	// releated field: StarPrize
	dAtA.Star = pArSeR.Uint64Array("star", "", false)
	// skip field: FirstDungeon
	// skip field: LastDungeon

	return dAtA, nil
}

var vAlIdAtOrDungeonChapterData = map[string]*config.Validator{

	"id":           config.ParseValidator("int>0", "", false, nil, nil),
	"chapter_name": config.ParseValidator("string>0", "", false, nil, nil),
	"chapter_desc": config.ParseValidator("string>0", "", false, nil, nil),
	"captain":      config.ParseValidator("string", "", false, nil, nil),
	"type":         config.ParseValidator("int", "", false, config.EnumMapKeys(shared_proto.DifficultType_value, 0), nil),
	"bg_img":       config.ParseValidator("string>0", "", false, nil, nil),
	"pass_prize":   config.ParseValidator("string", "", false, nil, nil),
	"star_prize":   config.ParseValidator("string", "", true, nil, nil),
	"star":         config.ParseValidator("uint", "", true, nil, nil),
}

func (dAtA *DungeonChapterData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *DungeonChapterData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *DungeonChapterData) Encode() *shared_proto.DungeonChapterProto {
	out := &shared_proto.DungeonChapterProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.ChapterName = dAtA.ChapterName
	out.ChapterDesc = dAtA.ChapterDesc
	if dAtA.Captain != nil {
		out.CaptainSoul = config.U64ToI32(dAtA.Captain.Id)
	}
	out.Type = dAtA.Type
	out.BgImg = dAtA.BgImg
	if dAtA.DungeonDatas != nil {
		out.DungeonDatas = ArrayEncodeDungeonData(dAtA.DungeonDatas)
	}
	if dAtA.PassPrize != nil {
		out.PassPrize = dAtA.PassPrize.Encode()
	}
	if dAtA.StarPrize != nil {
		out.StarPrize = resdata.ArrayEncodePrize(dAtA.StarPrize)
	}
	out.Star = config.U64a2I32a(dAtA.Star)

	return out
}

func ArrayEncodeDungeonChapterData(datas []*DungeonChapterData) []*shared_proto.DungeonChapterProto {

	out := make([]*shared_proto.DungeonChapterProto, 0, len(datas))
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

func (dAtA *DungeonChapterData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Captain = cOnFigS.GetCaptainData(pArSeR.Uint64("captain"))
	if dAtA.Captain == nil {
		return errors.Errorf("%s 配置的关联字段[captain] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("captain"), *pArSeR)
	}

	dAtA.PassPrize = cOnFigS.GetPrize(pArSeR.Int("pass_prize"))
	if dAtA.PassPrize == nil {
		return errors.Errorf("%s 配置的关联字段[pass_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("pass_prize"), *pArSeR)
	}

	intKeys = pArSeR.IntArray("star_prize", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetPrize(v)
		if obj != nil {
			dAtA.StarPrize = append(dAtA.StarPrize, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[star_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("star_prize"), *pArSeR)
		}
	}

	return nil
}

// start with DungeonData ----------------------------------

func LoadDungeonData(gos *config.GameObjects) (map[uint64]*DungeonData, map[*DungeonData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.DungeonDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*DungeonData, len(lIsT))
	pArSeRmAp := make(map[*DungeonData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrDungeonData) {
			continue
		}

		dAtA, err := NewDungeonData(fIlEnAmE, pArSeR)
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

func SetRelatedDungeonData(dAtAmAp map[*DungeonData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.DungeonDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetDungeonDataKeyArray(datas []*DungeonData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewDungeonData(fIlEnAmE string, pArSeR *config.ObjectParser) (*DungeonData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrDungeonData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &DungeonData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.UnlockHeroLevel = pArSeR.Uint64("unlock_hero_level")
	// releated field: UnlockPassDungeon
	dAtA.UnlockBayeStage = 0
	if pArSeR.KeyExist("unlock_baye_stage") {
		dAtA.UnlockBayeStage = pArSeR.Uint64("unlock_baye_stage")
	}

	dAtA.ChapterId = pArSeR.Uint64("chapter_id")
	dAtA.Type = shared_proto.DifficultType(shared_proto.DifficultType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.DifficultType(i)
	}

	dAtA.Star = 0
	if pArSeR.KeyExist("star") {
		dAtA.Star = pArSeR.Uint64("star")
	}

	dAtA.StarCondition = pArSeR.Uint64Array("star_condition", "", false)
	dAtA.StarConditionValue = pArSeR.Uint64Array("star_condition_value", "", false)
	dAtA.PassLimit = 0
	if pArSeR.KeyExist("pass_limit") {
		dAtA.PassLimit = pArSeR.Uint64("pass_limit")
	}

	dAtA.Sp = 0
	if pArSeR.KeyExist("sp") {
		dAtA.Sp = pArSeR.Uint64("sp")
	}

	dAtA.UnlockPuzzleIndex = pArSeR.Uint64("unlock_puzzle_index")
	dAtA.StoryId = pArSeR.Uint64("story_id")
	dAtA.DialogId = 0
	if pArSeR.KeyExist("dialog_id") {
		dAtA.DialogId = pArSeR.Uint64("dialog_id")
	}

	dAtA.PreBattleDialogId = 0
	if pArSeR.KeyExist("pre_battle_dialog_id") {
		dAtA.PreBattleDialogId = pArSeR.Uint64("pre_battle_dialog_id")
	}

	dAtA.AfterBattleDialogId = 0
	if pArSeR.KeyExist("after_battle_dialog_id") {
		dAtA.AfterBattleDialogId = pArSeR.Uint64("after_battle_dialog_id")
	}

	dAtA.BallonToolTip = "ballon_tool_tip"
	if pArSeR.KeyExist("ballon_tool_tip") {
		dAtA.BallonToolTip = pArSeR.String("ballon_tool_tip")
	}

	dAtA.NpcName = "npc_name"
	if pArSeR.KeyExist("npc_name") {
		dAtA.NpcName = pArSeR.String("npc_name")
	}

	dAtA.NpcIcon = "npc_icon"
	if pArSeR.KeyExist("npc_icon") {
		dAtA.NpcIcon = pArSeR.String("npc_icon")
	}

	// releated field: FirstPassPrize
	// releated field: PassPrize
	// releated field: Plunder
	// releated field: ShowPrize
	// releated field: Monster
	// releated field: CombatScene
	// skip field: CombatSceneRes
	// skip field: Prev
	// skip field: Next
	// releated field: YuanJunData
	dAtA.PlotIdx = pArSeR.Uint64Array("plot_idx", "", false)
	dAtA.PlotId = pArSeR.Uint64Array("plot_id", "", false)
	// skip field: GuideTroop

	return dAtA, nil
}

var vAlIdAtOrDungeonData = map[string]*config.Validator{

	"id":                   config.ParseValidator("int>0", "", false, nil, nil),
	"name":                 config.ParseValidator("string>0", "", false, nil, nil),
	"desc":                 config.ParseValidator("string>0", "", false, nil, nil),
	"unlock_hero_level":    config.ParseValidator("uint", "", false, nil, nil),
	"unlock_pass_dungeon":  config.ParseValidator("string", "", true, nil, nil),
	"unlock_baye_stage":    config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"chapter_id":           config.ParseValidator("int>0", "", false, nil, nil),
	"type":                 config.ParseValidator("int", "", false, config.EnumMapKeys(shared_proto.DifficultType_value, 0), nil),
	"star":                 config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"star_condition":       config.ParseValidator("uint", "", true, nil, nil),
	"star_condition_value": config.ParseValidator("uint,duplicate", "", true, nil, nil),
	"pass_limit":           config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"sp":                   config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"unlock_puzzle_index":    config.ParseValidator("uint", "", false, nil, nil),
	"story_id":               config.ParseValidator("uint", "", false, nil, nil),
	"dialog_id":              config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"pre_battle_dialog_id":   config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"after_battle_dialog_id": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"ballon_tool_tip":        config.ParseValidator("string", "", false, nil, []string{"ballon_tool_tip"}),
	"npc_name":               config.ParseValidator("string", "", false, nil, []string{"npc_name"}),
	"npc_icon":               config.ParseValidator("string", "", false, nil, []string{"npc_icon"}),
	"first_pass_prize":       config.ParseValidator("string", "", false, nil, nil),
	"pass_prize":             config.ParseValidator("string", "", false, nil, nil),
	"plunder":                config.ParseValidator("string", "", false, nil, nil),
	"show_prize":             config.ParseValidator("string", "", false, nil, nil),
	"monster":                config.ParseValidator("string", "", false, nil, nil),
	"combat_scene":           config.ParseValidator("string", "", false, nil, []string{"CombatScene"}),
	"yuan_jun_id":            config.ParseValidator("string", "", true, nil, nil),
	"plot_idx":               config.ParseValidator("uint", "", true, nil, nil),
	"plot_id":                config.ParseValidator("uint", "", true, nil, nil),
}

func (dAtA *DungeonData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *DungeonData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *DungeonData) Encode() *shared_proto.DungeonDataProto {
	out := &shared_proto.DungeonDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.UnlockHeroLevel = config.U64ToI32(dAtA.UnlockHeroLevel)
	if dAtA.UnlockPassDungeon != nil {
		out.UnlockPassDungeon = config.U64a2I32a(GetDungeonDataKeyArray(dAtA.UnlockPassDungeon))
	}
	out.UnlockBayeStage = config.U64ToI32(dAtA.UnlockBayeStage)
	out.ChapterId = config.U64ToI32(dAtA.ChapterId)
	out.Type = dAtA.Type
	out.Star = config.U64ToI32(dAtA.Star)
	out.StarCondition = config.U64a2I32a(dAtA.StarCondition)
	out.StarConditionValue = config.U64a2I32a(dAtA.StarConditionValue)
	out.PassLimit = config.U64ToI32(dAtA.PassLimit)
	out.Sp = config.U64ToI32(dAtA.Sp)
	out.UnlockPuzzleIndex = config.U64ToI32(dAtA.UnlockPuzzleIndex)
	out.StoryId = config.U64ToI32(dAtA.StoryId)
	out.DialogId = config.U64ToI32(dAtA.DialogId)
	out.PreBattleDialogId = config.U64ToI32(dAtA.PreBattleDialogId)
	out.AfterBattleDialogId = config.U64ToI32(dAtA.AfterBattleDialogId)
	out.BallonToolTip = dAtA.BallonToolTip
	out.NpcName = dAtA.NpcName
	out.NpcIcon = dAtA.NpcIcon
	if dAtA.FirstPassPrize != nil {
		out.FirstPassPrize = dAtA.FirstPassPrize.Encode()
	}
	if dAtA.PassPrize != nil {
		out.PassPrize = dAtA.PassPrize.Encode()
	}
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.Encode()
	}
	if dAtA.Monster != nil {
		out.Monster = dAtA.Monster.Encode()
	}
	out.CombatSceneRes = dAtA.CombatSceneRes
	if dAtA.Prev != nil {
		out.Prev = config.U64ToI32(dAtA.Prev.Id)
	}
	if dAtA.Next != nil {
		out.Next = config.U64ToI32(dAtA.Next.Id)
	}
	if dAtA.YuanJunData != nil {
		out.YuanJunData = monsterdata.ArrayEncodeMonsterCaptainData(dAtA.YuanJunData)
	}
	out.PlotIdx = config.U64a2I32a(dAtA.PlotIdx)
	out.PlotId = config.U64a2I32a(dAtA.PlotId)
	if dAtA.GuideTroop != nil {
		out.GuideTroop = dAtA.GuideTroop.Encode()
	}

	return out
}

func ArrayEncodeDungeonData(datas []*DungeonData) []*shared_proto.DungeonDataProto {

	out := make([]*shared_proto.DungeonDataProto, 0, len(datas))
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

func (dAtA *DungeonData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("unlock_pass_dungeon", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetDungeonData(v)
		if obj != nil {
			dAtA.UnlockPassDungeon = append(dAtA.UnlockPassDungeon, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[unlock_pass_dungeon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("unlock_pass_dungeon"), *pArSeR)
		}
	}

	dAtA.FirstPassPrize = cOnFigS.GetPrize(pArSeR.Int("first_pass_prize"))
	if dAtA.FirstPassPrize == nil && pArSeR.Int("first_pass_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[first_pass_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_pass_prize"), *pArSeR)
	}

	dAtA.PassPrize = cOnFigS.GetPrize(pArSeR.Int("pass_prize"))
	if dAtA.PassPrize == nil && pArSeR.Int("pass_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[pass_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("pass_prize"), *pArSeR)
	}

	dAtA.Plunder = cOnFigS.GetPlunder(pArSeR.Uint64("plunder"))
	if dAtA.Plunder == nil {
		return errors.Errorf("%s 配置的关联字段[plunder] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder"), *pArSeR)
	}

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil {
		return errors.Errorf("%s 配置的关联字段[show_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
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

	uint64Keys = pArSeR.Uint64Array("yuan_jun_id", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetMonsterCaptainData(v)
		if obj != nil {
			dAtA.YuanJunData = append(dAtA.YuanJunData, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[yuan_jun_id] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("yuan_jun_id"), *pArSeR)
		}
	}

	return nil
}

// start with DungeonGuideTroopData ----------------------------------

func LoadDungeonGuideTroopData(gos *config.GameObjects) (map[uint64]*DungeonGuideTroopData, map[*DungeonGuideTroopData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.DungeonGuideTroopDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*DungeonGuideTroopData, len(lIsT))
	pArSeRmAp := make(map[*DungeonGuideTroopData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrDungeonGuideTroopData) {
			continue
		}

		dAtA, err := NewDungeonGuideTroopData(fIlEnAmE, pArSeR)
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

func SetRelatedDungeonGuideTroopData(dAtAmAp map[*DungeonGuideTroopData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.DungeonGuideTroopDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetDungeonGuideTroopDataKeyArray(datas []*DungeonGuideTroopData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewDungeonGuideTroopData(fIlEnAmE string, pArSeR *config.ObjectParser) (*DungeonGuideTroopData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrDungeonGuideTroopData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &DungeonGuideTroopData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.NotFirst = pArSeR.Bool("not_first")
	dAtA.Captain = pArSeR.Uint64Array("captain", "", false)
	dAtA.SrcPos = pArSeR.Uint64Array("src_pos", "", false)
	dAtA.SrcPosX = pArSeR.Uint64Array("src_pos_x", "", false)
	dAtA.DstPos = pArSeR.Uint64Array("dst_pos", "", false)
	dAtA.DstPosX = pArSeR.Uint64Array("dst_pos_x", "", false)

	return dAtA, nil
}

var vAlIdAtOrDungeonGuideTroopData = map[string]*config.Validator{

	"id":        config.ParseValidator("int>0", "", false, nil, nil),
	"not_first": config.ParseValidator("bool", "", false, nil, nil),
	"captain":   config.ParseValidator("uint", "", true, nil, nil),
	"src_pos":   config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
	"src_pos_x": config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
	"dst_pos":   config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
	"dst_pos_x": config.ParseValidator("int>0,duplicate,notAllNil", "", true, nil, nil),
}

func (dAtA *DungeonGuideTroopData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *DungeonGuideTroopData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *DungeonGuideTroopData) Encode() *shared_proto.DungeonGuideTroopDataProto {
	out := &shared_proto.DungeonGuideTroopDataProto{}
	out.NotFirst = dAtA.NotFirst
	out.Captain = config.U64a2I32a(dAtA.Captain)
	out.SrcPos = config.U64a2I32a(dAtA.SrcPos)
	out.SrcPosX = config.U64a2I32a(dAtA.SrcPosX)
	out.DstPos = config.U64a2I32a(dAtA.DstPos)
	out.DstPosX = config.U64a2I32a(dAtA.DstPosX)

	return out
}

func ArrayEncodeDungeonGuideTroopData(datas []*DungeonGuideTroopData) []*shared_proto.DungeonGuideTroopDataProto {

	out := make([]*shared_proto.DungeonGuideTroopDataProto, 0, len(datas))
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

func (dAtA *DungeonGuideTroopData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with DungeonMiscData ----------------------------------

func LoadDungeonMiscData(gos *config.GameObjects) (*DungeonMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.DungeonMiscDataPath
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

	dAtA, err := NewDungeonMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedDungeonMiscData(gos *config.GameObjects, dAtA *DungeonMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.DungeonMiscDataPath
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

func NewDungeonMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*DungeonMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrDungeonMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &DungeonMiscData{}

	dAtA.MaxAutoTimes = 100
	if pArSeR.KeyExist("max_auto_times") {
		dAtA.MaxAutoTimes = pArSeR.Uint64("max_auto_times")
	}

	if pArSeR.KeyExist("recover_auto_duration") {
		dAtA.RecoverAutoDuration, err = config.ParseDuration(pArSeR.String("recover_auto_duration"))
	} else {
		dAtA.RecoverAutoDuration, err = config.ParseDuration("10m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[recover_auto_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("recover_auto_duration"), dAtA)
	}

	dAtA.DefaultAutoTimes = 50
	if pArSeR.KeyExist("default_auto_times") {
		dAtA.DefaultAutoTimes = pArSeR.Uint64("default_auto_times")
	}

	dAtA.AutoPerTimes = 5
	if pArSeR.KeyExist("auto_per_times") {
		dAtA.AutoPerTimes = pArSeR.Uint64("auto_per_times")
	}

	return dAtA, nil
}

var vAlIdAtOrDungeonMiscData = map[string]*config.Validator{

	"max_auto_times":        config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"recover_auto_duration": config.ParseValidator("string", "", false, nil, []string{"10m"}),
	"default_auto_times":    config.ParseValidator("int>0", "", false, nil, []string{"50"}),
	"auto_per_times":        config.ParseValidator("int>0", "", false, nil, []string{"5"}),
}

func (dAtA *DungeonMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *DungeonMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *DungeonMiscData) Encode() *shared_proto.DungeonMiscProto {
	out := &shared_proto.DungeonMiscProto{}
	out.MaxAutoTimes = config.U64ToI32(dAtA.MaxAutoTimes)
	out.RecoverAutoDuration = config.Duration2I32Seconds(dAtA.RecoverAutoDuration)
	out.AutoPerTimes = config.U64ToI32(dAtA.AutoPerTimes)

	return out
}

func ArrayEncodeDungeonMiscData(datas []*DungeonMiscData) []*shared_proto.DungeonMiscProto {

	out := make([]*shared_proto.DungeonMiscProto, 0, len(datas))
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

func (dAtA *DungeonMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetCaptainData(uint64) *captain.CaptainData
	GetCombatScene(string) *scene.CombatScene
	GetDungeonData(uint64) *DungeonData
	GetMonsterCaptainData(uint64) *monsterdata.MonsterCaptainData
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
	GetPlunder(uint64) *resdata.Plunder
	GetPrize(int) *resdata.Prize
}
