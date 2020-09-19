// AUTO_GEN, DONT MODIFY!!!
package teach

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/monsterdata"
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

// start with TeachChapterData ----------------------------------

func LoadTeachChapterData(gos *config.GameObjects) (map[uint64]*TeachChapterData, map[*TeachChapterData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TeachChapterDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TeachChapterData, len(lIsT))
	pArSeRmAp := make(map[*TeachChapterData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTeachChapterData) {
			continue
		}

		dAtA, err := NewTeachChapterData(fIlEnAmE, pArSeR)
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

func SetRelatedTeachChapterData(dAtAmAp map[*TeachChapterData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TeachChapterDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTeachChapterDataKeyArray(datas []*TeachChapterData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTeachChapterData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TeachChapterData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTeachChapterData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TeachChapterData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.MinHeroLevel = pArSeR.Uint64("min_hero_level")
	dAtA.PassBaYeTaskStage = pArSeR.Uint64("pass_ba_ye_task_stage")
	dAtA.PassDungeonId = pArSeR.Uint64("pass_dungeon_id")
	dAtA.Title = pArSeR.String("title")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Image = pArSeR.String("image")
	// releated field: Prize
	// releated field: AtkStartMonster
	// releated field: AtkEndMonster
	// releated field: DefMonster
	// skip field: PrevData

	return dAtA, nil
}

var vAlIdAtOrTeachChapterData = map[string]*config.Validator{

	"id":                    config.ParseValidator("int>0", "", false, nil, nil),
	"min_hero_level":        config.ParseValidator("int>0", "", false, nil, nil),
	"pass_ba_ye_task_stage": config.ParseValidator("uint", "", false, nil, nil),
	"pass_dungeon_id":       config.ParseValidator("uint", "", false, nil, nil),
	"title":                 config.ParseValidator("string", "", false, nil, nil),
	"desc":                  config.ParseValidator("string", "", false, nil, nil),
	"image":                 config.ParseValidator("string", "", false, nil, nil),
	"prize":                 config.ParseValidator("string", "", false, nil, nil),
	"atk_start_monster":     config.ParseValidator("string", "", false, nil, nil),
	"atk_end_monster":       config.ParseValidator("string", "", false, nil, nil),
	"def_monster":           config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *TeachChapterData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TeachChapterData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TeachChapterData) Encode() *shared_proto.TeachChapterDataProto {
	out := &shared_proto.TeachChapterDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.MinHeroLevel = config.U64ToI32(dAtA.MinHeroLevel)
	out.PassBaYeTaskStage = config.U64ToI32(dAtA.PassBaYeTaskStage)
	out.PassDungeonId = config.U64ToI32(dAtA.PassDungeonId)
	out.Title = dAtA.Title
	out.Desc = dAtA.Desc
	out.Image = dAtA.Image
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	if dAtA.AtkStartMonster != nil {
		out.AtkStartMonster = dAtA.AtkStartMonster.Encode()
	}
	if dAtA.AtkEndMonster != nil {
		out.AtkEndMonster = dAtA.AtkEndMonster.Encode()
	}
	if dAtA.DefMonster != nil {
		out.DefMonster = dAtA.DefMonster.Encode()
	}

	return out
}

func ArrayEncodeTeachChapterData(datas []*TeachChapterData) []*shared_proto.TeachChapterDataProto {

	out := make([]*shared_proto.TeachChapterDataProto, 0, len(datas))
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

func (dAtA *TeachChapterData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	dAtA.AtkStartMonster = cOnFigS.GetMonsterMasterData(pArSeR.Uint64("atk_start_monster"))
	if dAtA.AtkStartMonster == nil {
		return errors.Errorf("%s 配置的关联字段[atk_start_monster] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("atk_start_monster"), *pArSeR)
	}

	dAtA.AtkEndMonster = cOnFigS.GetMonsterMasterData(pArSeR.Uint64("atk_end_monster"))
	if dAtA.AtkEndMonster == nil {
		return errors.Errorf("%s 配置的关联字段[atk_end_monster] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("atk_end_monster"), *pArSeR)
	}

	dAtA.DefMonster = cOnFigS.GetMonsterMasterData(pArSeR.Uint64("def_monster"))
	if dAtA.DefMonster == nil {
		return errors.Errorf("%s 配置的关联字段[def_monster] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("def_monster"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
	GetPrize(int) *resdata.Prize
}
