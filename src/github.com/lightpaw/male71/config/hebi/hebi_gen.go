// AUTO_GEN, DONT MODIFY!!!
package hebi

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/goods"
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

// start with HebiMiscData ----------------------------------

func LoadHebiMiscData(gos *config.GameObjects) (*HebiMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.HebiMiscDataPath
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

	dAtA, err := NewHebiMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedHebiMiscData(gos *config.GameObjects, dAtA *HebiMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.HebiMiscDataPath
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

func NewHebiMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*HebiMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrHebiMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &HebiMiscData{}

	dAtA.RoomsMaxSize = 100
	if pArSeR.KeyExist("rooms_max_size") {
		dAtA.RoomsMaxSize = pArSeR.Uint64("rooms_max_size")
	}

	dAtA.DailyRobCount = 3
	if pArSeR.KeyExist("daily_rob_count") {
		dAtA.DailyRobCount = pArSeR.Uint64("daily_rob_count")
	}

	dAtA.RobCdDuration, err = config.ParseDuration(pArSeR.String("rob_cd_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[rob_cd_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("rob_cd_duration"), dAtA)
	}

	if pArSeR.KeyExist("rob_pos_cd_duration") {
		dAtA.RobPosCdDuration, err = config.ParseDuration(pArSeR.String("rob_pos_cd_duration"))
	} else {
		dAtA.RobPosCdDuration, err = config.ParseDuration("10s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[rob_pos_cd_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("rob_pos_cd_duration"), dAtA)
	}

	dAtA.RobProtectDuration, err = config.ParseDuration(pArSeR.String("rob_protect_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[rob_protect_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("rob_protect_duration"), dAtA)
	}

	dAtA.HebiDuration, err = config.ParseDuration(pArSeR.String("hebi_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[hebi_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("hebi_duration"), dAtA)
	}

	dAtA.HeShiBiType = shared_proto.HebiType(shared_proto.HebiType_value[strings.ToUpper(pArSeR.String("he_shi_bi_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("he_shi_bi_type"), 10, 32); err == nil {
		dAtA.HeShiBiType = shared_proto.HebiType(i)
	}

	if pArSeR.KeyExist("room_wait_expired_duration") {
		dAtA.RoomWaitExpiredDuration, err = config.ParseDuration(pArSeR.String("room_wait_expired_duration"))
	} else {
		dAtA.RoomWaitExpiredDuration, err = config.ParseDuration("3h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[room_wait_expired_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("room_wait_expired_duration"), dAtA)
	}

	dAtA.HebiHeroRecordMaxSize = 50
	if pArSeR.KeyExist("hebi_hero_record_max_size") {
		dAtA.HebiHeroRecordMaxSize = pArSeR.Uint64("hebi_hero_record_max_size")
	}

	// releated field: CombatScene
	// releated field: CopySelfGoods
	// skip field: HebiGoods

	return dAtA, nil
}

var vAlIdAtOrHebiMiscData = map[string]*config.Validator{

	"rooms_max_size":             config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"daily_rob_count":            config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"rob_cd_duration":            config.ParseValidator("string", "", false, nil, nil),
	"rob_pos_cd_duration":        config.ParseValidator("string", "", false, nil, []string{"10s"}),
	"rob_protect_duration":       config.ParseValidator("string", "", false, nil, nil),
	"hebi_duration":              config.ParseValidator("string", "", false, nil, nil),
	"he_shi_bi_type":             config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.HebiType_value, 0), nil),
	"room_wait_expired_duration": config.ParseValidator("string", "", false, nil, []string{"3h"}),
	"hebi_hero_record_max_size":  config.ParseValidator("int>0", "", false, nil, []string{"50"}),
	"combat_scene":               config.ParseValidator("string", "", false, nil, []string{"CombatScene"}),
	"copy_self_goods":            config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *HebiMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *HebiMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *HebiMiscData) Encode() *shared_proto.HebiMiscProto {
	out := &shared_proto.HebiMiscProto{}
	out.RoomsMaxSize = config.U64ToI32(dAtA.RoomsMaxSize)
	out.DailyRobCount = config.U64ToI32(dAtA.DailyRobCount)
	out.RobCdDuration = config.Duration2I32Seconds(dAtA.RobCdDuration)
	out.RobPosCdDuration = config.Duration2I32Seconds(dAtA.RobPosCdDuration)
	out.RobProtectDuration = config.Duration2I32Seconds(dAtA.RobProtectDuration)
	out.HebiDuration = config.Duration2I32Seconds(dAtA.HebiDuration)
	out.HeShiBiType = dAtA.HeShiBiType
	if dAtA.CombatScene != nil {
		out.CombatScene = dAtA.CombatScene.Id
	}
	if dAtA.CopySelfGoods != nil {
		out.CopySelfGoods = config.U64ToI32(dAtA.CopySelfGoods.Id)
	}

	return out
}

func ArrayEncodeHebiMiscData(datas []*HebiMiscData) []*shared_proto.HebiMiscProto {

	out := make([]*shared_proto.HebiMiscProto, 0, len(datas))
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

func (dAtA *HebiMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("combat_scene") {
		dAtA.CombatScene = cOnFigS.GetCombatScene(pArSeR.String("combat_scene"))
	} else {
		dAtA.CombatScene = cOnFigS.GetCombatScene("CombatScene")
	}
	if dAtA.CombatScene == nil {
		return errors.Errorf("%s 配置的关联字段[combat_scene] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("combat_scene"), *pArSeR)
	}

	dAtA.CopySelfGoods = cOnFigS.GetGoodsData(pArSeR.Uint64("copy_self_goods"))
	if dAtA.CopySelfGoods == nil {
		return errors.Errorf("%s 配置的关联字段[copy_self_goods] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("copy_self_goods"), *pArSeR)
	}

	return nil
}

// start with HebiPrizeData ----------------------------------

func LoadHebiPrizeData(gos *config.GameObjects) (map[uint64]*HebiPrizeData, map[*HebiPrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.HebiPrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*HebiPrizeData, len(lIsT))
	pArSeRmAp := make(map[*HebiPrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrHebiPrizeData) {
			continue
		}

		dAtA, err := NewHebiPrizeData(fIlEnAmE, pArSeR)
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

func SetRelatedHebiPrizeData(dAtAmAp map[*HebiPrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.HebiPrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetHebiPrizeDataKeyArray(datas []*HebiPrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewHebiPrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*HebiPrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrHebiPrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &HebiPrizeData{}

	dAtA.HeroLevel = pArSeR.Uint64("hero_level")
	dAtA.HebiType = shared_proto.HebiType(shared_proto.HebiType_value[strings.ToUpper(pArSeR.String("hebi_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("hebi_type"), 10, 32); err == nil {
		dAtA.HebiType = shared_proto.HebiType(i)
	}

	dAtA.Quality = pArSeR.Uint64("quality")
	// releated field: Prize
	// releated field: AmountPrize
	// releated field: PlunderPrize

	// calculate fields
	dAtA.Id = GenHebiPrizeId(dAtA.HeroLevel, dAtA.HebiType, dAtA.Quality)

	return dAtA, nil
}

var vAlIdAtOrHebiPrizeData = map[string]*config.Validator{

	"hero_level":    config.ParseValidator("int>0", "", false, nil, nil),
	"hebi_type":     config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.HebiType_value, 0), nil),
	"quality":       config.ParseValidator("int>0", "", false, nil, nil),
	"prize":         config.ParseValidator("string", "", false, nil, nil),
	"amount_prize":  config.ParseValidator("string", "", false, nil, nil),
	"plunder_prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *HebiPrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *HebiPrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *HebiPrizeData) Encode() *shared_proto.HebiPrizeProto {
	out := &shared_proto.HebiPrizeProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.HeroLevel = config.U64ToI32(dAtA.HeroLevel)
	out.HebiType = dAtA.HebiType
	out.Quality = config.U64ToI32(dAtA.Quality)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeHebiPrizeData(datas []*HebiPrizeData) []*shared_proto.HebiPrizeProto {

	out := make([]*shared_proto.HebiPrizeProto, 0, len(datas))
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

func (dAtA *HebiPrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.AmountPrize = cOnFigS.GetPrize(pArSeR.Int("amount_prize"))
	if dAtA.AmountPrize == nil {
		return errors.Errorf("%s 配置的关联字段[amount_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("amount_prize"), *pArSeR)
	}

	dAtA.PlunderPrize = cOnFigS.GetPlunderPrize(pArSeR.Uint64("plunder_prize"))
	if dAtA.PlunderPrize == nil {
		return errors.Errorf("%s 配置的关联字段[plunder_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder_prize"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetCombatScene(string) *scene.CombatScene
	GetGoodsData(uint64) *goods.GoodsData
	GetPlunderPrize(uint64) *resdata.PlunderPrize
	GetPrize(int) *resdata.Prize
}
