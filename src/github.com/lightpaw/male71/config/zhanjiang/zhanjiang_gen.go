// AUTO_GEN, DONT MODIFY!!!
package zhanjiang

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

// start with ZhanJiangChapterData ----------------------------------

func LoadZhanJiangChapterData(gos *config.GameObjects) (map[uint64]*ZhanJiangChapterData, map[*ZhanJiangChapterData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ZhanJiangChapterDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ZhanJiangChapterData, len(lIsT))
	pArSeRmAp := make(map[*ZhanJiangChapterData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrZhanJiangChapterData) {
			continue
		}

		dAtA, err := NewZhanJiangChapterData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.ChapterId
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[ChapterId], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedZhanJiangChapterData(dAtAmAp map[*ZhanJiangChapterData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ZhanJiangChapterDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetZhanJiangChapterDataKeyArray(datas []*ZhanJiangChapterData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.ChapterId)
		}
	}

	return out
}

func NewZhanJiangChapterData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ZhanJiangChapterData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrZhanJiangChapterData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ZhanJiangChapterData{}

	dAtA.ChapterId = pArSeR.Uint64("chapter_id")
	dAtA.ChapterName = pArSeR.String("chapter_name")
	dAtA.ChapterDesc = pArSeR.String("chapter_desc")
	dAtA.BgImg = pArSeR.String("bg_img")
	// releated field: ZhanJiangDatas
	// skip field: PreChapter

	return dAtA, nil
}

var vAlIdAtOrZhanJiangChapterData = map[string]*config.Validator{

	"chapter_id":   config.ParseValidator("int>0", "", false, nil, nil),
	"chapter_name": config.ParseValidator("string>0", "", false, nil, nil),
	"chapter_desc": config.ParseValidator("string>0", "", false, nil, nil),
	"bg_img":       config.ParseValidator("string>0", "", false, nil, nil),
	"guan_qia":     config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *ZhanJiangChapterData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ZhanJiangChapterData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ZhanJiangChapterData) Encode() *shared_proto.ZhanJiangChapterProto {
	out := &shared_proto.ZhanJiangChapterProto{}
	out.ChapterId = config.U64ToI32(dAtA.ChapterId)
	out.ChapterName = dAtA.ChapterName
	out.ChapterDesc = dAtA.ChapterDesc
	out.BgImg = dAtA.BgImg
	if dAtA.ZhanJiangDatas != nil {
		out.GuanQia = ArrayEncodeZhanJiangGuanQiaData(dAtA.ZhanJiangDatas)
	}

	return out
}

func ArrayEncodeZhanJiangChapterData(datas []*ZhanJiangChapterData) []*shared_proto.ZhanJiangChapterProto {

	out := make([]*shared_proto.ZhanJiangChapterProto, 0, len(datas))
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

func (dAtA *ZhanJiangChapterData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("guan_qia", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetZhanJiangGuanQiaData(v)
		if obj != nil {
			dAtA.ZhanJiangDatas = append(dAtA.ZhanJiangDatas, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[guan_qia] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guan_qia"), *pArSeR)
		}
	}

	return nil
}

// start with ZhanJiangData ----------------------------------

func LoadZhanJiangData(gos *config.GameObjects) (map[uint64]*ZhanJiangData, map[*ZhanJiangData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ZhanJiangDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ZhanJiangData, len(lIsT))
	pArSeRmAp := make(map[*ZhanJiangData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrZhanJiangData) {
			continue
		}

		dAtA, err := NewZhanJiangData(fIlEnAmE, pArSeR)
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

func SetRelatedZhanJiangData(dAtAmAp map[*ZhanJiangData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ZhanJiangDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetZhanJiangDataKeyArray(datas []*ZhanJiangData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewZhanJiangData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ZhanJiangData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrZhanJiangData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ZhanJiangData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Icon = pArSeR.String("icon")
	// releated field: PassPrize
	// releated field: Plunder
	// releated field: ShowPrize
	// releated field: Monster
	// releated field: CombatScene
	dAtA.GongXun = pArSeR.Uint64("gong_xun")

	return dAtA, nil
}

var vAlIdAtOrZhanJiangData = map[string]*config.Validator{

	"id":           config.ParseValidator("int>0", "", false, nil, nil),
	"name":         config.ParseValidator("string>0", "", false, nil, nil),
	"desc":         config.ParseValidator("string>0", "", false, nil, nil),
	"icon":         config.ParseValidator("string>0", "", false, nil, nil),
	"pass_prize":   config.ParseValidator("string", "", false, nil, nil),
	"plunder":      config.ParseValidator("string", "", false, nil, nil),
	"show_prize":   config.ParseValidator("string", "", false, nil, nil),
	"monster":      config.ParseValidator("string", "", false, nil, nil),
	"combat_scene": config.ParseValidator("string", "", false, nil, []string{"CombatScene"}),
	"gong_xun":     config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *ZhanJiangData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ZhanJiangData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ZhanJiangData) Encode() *shared_proto.ZhanJiangDataProto {
	out := &shared_proto.ZhanJiangDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.Icon = dAtA.Icon
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.Encode()
	}
	if dAtA.Monster != nil {
		out.Monster = dAtA.Monster.Encode()
	}
	out.GongXun = config.U64ToI32(dAtA.GongXun)

	return out
}

func ArrayEncodeZhanJiangData(datas []*ZhanJiangData) []*shared_proto.ZhanJiangDataProto {

	out := make([]*shared_proto.ZhanJiangDataProto, 0, len(datas))
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

func (dAtA *ZhanJiangData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.PassPrize = cOnFigS.GetPrize(pArSeR.Int("pass_prize"))
	if dAtA.PassPrize == nil && pArSeR.Int("pass_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[pass_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("pass_prize"), *pArSeR)
	}

	dAtA.Plunder = cOnFigS.GetPlunder(pArSeR.Uint64("plunder"))
	if dAtA.Plunder == nil && pArSeR.Uint64("plunder") != 0 {
		return errors.Errorf("%s 配置的关联字段[plunder] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder"), *pArSeR)
	}

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil && pArSeR.Int("show_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[show_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
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

// start with ZhanJiangGuanQiaData ----------------------------------

func LoadZhanJiangGuanQiaData(gos *config.GameObjects) (map[uint64]*ZhanJiangGuanQiaData, map[*ZhanJiangGuanQiaData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ZhanJiangGuanQiaDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ZhanJiangGuanQiaData, len(lIsT))
	pArSeRmAp := make(map[*ZhanJiangGuanQiaData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrZhanJiangGuanQiaData) {
			continue
		}

		dAtA, err := NewZhanJiangGuanQiaData(fIlEnAmE, pArSeR)
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

func SetRelatedZhanJiangGuanQiaData(dAtAmAp map[*ZhanJiangGuanQiaData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ZhanJiangGuanQiaDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetZhanJiangGuanQiaDataKeyArray(datas []*ZhanJiangGuanQiaData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewZhanJiangGuanQiaData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ZhanJiangGuanQiaData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrZhanJiangGuanQiaData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ZhanJiangGuanQiaData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.PositionDesc = pArSeR.String("position_desc")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.BgImg = pArSeR.String("bg_img")
	// releated field: ZhanJiangDatas
	dAtA.AbilityExp = pArSeR.Uint64("ability_exp")
	// skip field: Prev
	// skip field: Next
	dAtA.ShowGongXun = pArSeR.Uint64("show_gong_xun")
	// releated field: ShowPrize
	// skip field: ChapterData

	return dAtA, nil
}

var vAlIdAtOrZhanJiangGuanQiaData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"name":          config.ParseValidator("string>0", "", false, nil, nil),
	"position_desc": config.ParseValidator("string>0", "", false, nil, nil),
	"desc":          config.ParseValidator("string>0", "", false, nil, nil),
	"bg_img":        config.ParseValidator("string>0", "", false, nil, nil),
	"guan":          config.ParseValidator("string", "", true, nil, nil),
	"ability_exp":   config.ParseValidator("uint", "", false, nil, nil),
	"show_gong_xun": config.ParseValidator("uint", "", false, nil, nil),
	"show_prize":    config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ZhanJiangGuanQiaData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ZhanJiangGuanQiaData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ZhanJiangGuanQiaData) Encode() *shared_proto.ZhanJiangGuanQiaProto {
	out := &shared_proto.ZhanJiangGuanQiaProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.PositionDesc = dAtA.PositionDesc
	out.Desc = dAtA.Desc
	out.BgImg = dAtA.BgImg
	if dAtA.ZhanJiangDatas != nil {
		out.Guan = ArrayEncodeZhanJiangData(dAtA.ZhanJiangDatas)
	}
	out.AbilityExp = config.U64ToI32(dAtA.AbilityExp)
	if dAtA.Prev != nil {
		out.Prev = config.U64ToI32(dAtA.Prev.Id)
	}
	if dAtA.Next != nil {
		out.Next = config.U64ToI32(dAtA.Next.Id)
	}
	out.ShowGongXun = config.U64ToI32(dAtA.ShowGongXun)
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.Encode()
	}

	return out
}

func ArrayEncodeZhanJiangGuanQiaData(datas []*ZhanJiangGuanQiaData) []*shared_proto.ZhanJiangGuanQiaProto {

	out := make([]*shared_proto.ZhanJiangGuanQiaProto, 0, len(datas))
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

func (dAtA *ZhanJiangGuanQiaData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("guan", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetZhanJiangData(v)
		if obj != nil {
			dAtA.ZhanJiangDatas = append(dAtA.ZhanJiangDatas, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[guan] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("guan"), *pArSeR)
		}
	}

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil && pArSeR.Int("show_prize") != 0 {
		return errors.Errorf("%s 配置的关联字段[show_prize] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
	}

	return nil
}

// start with ZhanJiangMiscData ----------------------------------

func LoadZhanJiangMiscData(gos *config.GameObjects) (*ZhanJiangMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.ZhanJiangMiscDataPath
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

	dAtA, err := NewZhanJiangMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedZhanJiangMiscData(gos *config.GameObjects, dAtA *ZhanJiangMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ZhanJiangMiscDataPath
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

func NewZhanJiangMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ZhanJiangMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrZhanJiangMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ZhanJiangMiscData{}

	dAtA.DefaultTimes = 3
	if pArSeR.KeyExist("default_times") {
		dAtA.DefaultTimes = pArSeR.Uint64("default_times")
	}

	return dAtA, nil
}

var vAlIdAtOrZhanJiangMiscData = map[string]*config.Validator{

	"default_times": config.ParseValidator("int>0", "", false, nil, []string{"3"}),
}

func (dAtA *ZhanJiangMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ZhanJiangMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ZhanJiangMiscData) Encode() *shared_proto.ZhanJiangMiscDataProto {
	out := &shared_proto.ZhanJiangMiscDataProto{}
	out.MaxTimes = config.U64ToI32(dAtA.DefaultTimes)

	return out
}

func ArrayEncodeZhanJiangMiscData(datas []*ZhanJiangMiscData) []*shared_proto.ZhanJiangMiscDataProto {

	out := make([]*shared_proto.ZhanJiangMiscDataProto, 0, len(datas))
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

func (dAtA *ZhanJiangMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetCombatScene(string) *scene.CombatScene
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
	GetPlunder(uint64) *resdata.Plunder
	GetPrize(int) *resdata.Prize
	GetZhanJiangData(uint64) *ZhanJiangData
	GetZhanJiangGuanQiaData(uint64) *ZhanJiangGuanQiaData
}
