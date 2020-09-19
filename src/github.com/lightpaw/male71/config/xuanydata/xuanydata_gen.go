// AUTO_GEN, DONT MODIFY!!!
package xuanydata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
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

// start with XuanyuanMiscData ----------------------------------

func LoadXuanyuanMiscData(gos *config.GameObjects) (*XuanyuanMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.XuanyuanMiscDataPath
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

	dAtA, err := NewXuanyuanMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedXuanyuanMiscData(gos *config.GameObjects, dAtA *XuanyuanMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.XuanyuanMiscDataPath
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

func NewXuanyuanMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*XuanyuanMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrXuanyuanMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &XuanyuanMiscData{}

	dAtA.ChallengeTimesLimit = pArSeR.Uint64("challenge_times_limit")
	dAtA.DailyMaxLostScore = pArSeR.Uint64("daily_max_lost_score")
	dAtA.RankCount = 10000
	if pArSeR.KeyExist("rank_count") {
		dAtA.RankCount = pArSeR.Uint64("rank_count")
	}

	dAtA.RecordBatchCount = 10
	if pArSeR.KeyExist("record_batch_count") {
		dAtA.RecordBatchCount = pArSeR.Uint64("record_batch_count")
	}

	dAtA.InitScore = 200
	if pArSeR.KeyExist("init_score") {
		dAtA.InitScore = pArSeR.Uint64("init_score")
	}

	return dAtA, nil
}

var vAlIdAtOrXuanyuanMiscData = map[string]*config.Validator{

	"challenge_times_limit": config.ParseValidator("int>0", "", false, nil, nil),
	"daily_max_lost_score":  config.ParseValidator("int>0", "", false, nil, nil),
	"rank_count":            config.ParseValidator("int>0", "", false, nil, []string{"10000"}),
	"record_batch_count":    config.ParseValidator("int>0", "", false, nil, []string{"10"}),
	"init_score":            config.ParseValidator("int>0", "", false, nil, []string{"200"}),
}

func (dAtA *XuanyuanMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *XuanyuanMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *XuanyuanMiscData) Encode() *shared_proto.XuanyuanMiscDataProto {
	out := &shared_proto.XuanyuanMiscDataProto{}
	out.ChallengeTimesLimit = config.U64ToI32(dAtA.ChallengeTimesLimit)

	return out
}

func ArrayEncodeXuanyuanMiscData(datas []*XuanyuanMiscData) []*shared_proto.XuanyuanMiscDataProto {

	out := make([]*shared_proto.XuanyuanMiscDataProto, 0, len(datas))
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

func (dAtA *XuanyuanMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with XuanyuanRangeData ----------------------------------

func LoadXuanyuanRangeData(gos *config.GameObjects) (map[uint64]*XuanyuanRangeData, map[*XuanyuanRangeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.XuanyuanRangeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*XuanyuanRangeData, len(lIsT))
	pArSeRmAp := make(map[*XuanyuanRangeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrXuanyuanRangeData) {
			continue
		}

		dAtA, err := NewXuanyuanRangeData(fIlEnAmE, pArSeR)
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

func SetRelatedXuanyuanRangeData(dAtAmAp map[*XuanyuanRangeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.XuanyuanRangeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetXuanyuanRangeDataKeyArray(datas []*XuanyuanRangeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewXuanyuanRangeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*XuanyuanRangeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrXuanyuanRangeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &XuanyuanRangeData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.LowRank = pArSeR.Int("low_rank")
	dAtA.HighRank = pArSeR.Int("high_rank")
	dAtA.WinScore = pArSeR.Uint64("win_score")
	dAtA.LoseScore = pArSeR.Uint64("lose_score")
	dAtA.DefenseLostScore = pArSeR.Uint64("defense_lost_score")
	// releated field: CombatScene

	return dAtA, nil
}

var vAlIdAtOrXuanyuanRangeData = map[string]*config.Validator{

	"id":                 config.ParseValidator("int>0", "", false, nil, nil),
	"low_rank":           config.ParseValidator("int", "", false, nil, nil),
	"high_rank":          config.ParseValidator("int", "", false, nil, nil),
	"win_score":          config.ParseValidator("int>0", "", false, nil, nil),
	"lose_score":         config.ParseValidator("int", "", false, nil, nil),
	"defense_lost_score": config.ParseValidator("int", "", false, nil, nil),
	"combat_scene":       config.ParseValidator("string", "", false, nil, []string{"CombatScene"}),
}

func (dAtA *XuanyuanRangeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *XuanyuanRangeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *XuanyuanRangeData) Encode() *shared_proto.XuanyuanRangeDataProto {
	out := &shared_proto.XuanyuanRangeDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.LowRank = int32(dAtA.LowRank)
	out.HighRank = int32(dAtA.HighRank)
	out.WinScore = config.U64ToI32(dAtA.WinScore)
	out.LoseScore = config.U64ToI32(dAtA.LoseScore)
	out.DefenseLostScore = config.U64ToI32(dAtA.DefenseLostScore)
	if dAtA.CombatScene != nil {
		out.CombatScene = dAtA.CombatScene.Id
	}

	return out
}

func ArrayEncodeXuanyuanRangeData(datas []*XuanyuanRangeData) []*shared_proto.XuanyuanRangeDataProto {

	out := make([]*shared_proto.XuanyuanRangeDataProto, 0, len(datas))
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

func (dAtA *XuanyuanRangeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	return nil
}

// start with XuanyuanRankPrizeData ----------------------------------

func LoadXuanyuanRankPrizeData(gos *config.GameObjects) (map[uint64]*XuanyuanRankPrizeData, map[*XuanyuanRankPrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.XuanyuanRankPrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*XuanyuanRankPrizeData, len(lIsT))
	pArSeRmAp := make(map[*XuanyuanRankPrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrXuanyuanRankPrizeData) {
			continue
		}

		dAtA, err := NewXuanyuanRankPrizeData(fIlEnAmE, pArSeR)
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

func SetRelatedXuanyuanRankPrizeData(dAtAmAp map[*XuanyuanRankPrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.XuanyuanRankPrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetXuanyuanRankPrizeDataKeyArray(datas []*XuanyuanRankPrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewXuanyuanRankPrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*XuanyuanRankPrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrXuanyuanRankPrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &XuanyuanRankPrizeData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Rank = pArSeR.Uint64("rank")
	// skip field: Prize
	// releated field: PlunderPrize
	// releated field: ShowPrize

	return dAtA, nil
}

var vAlIdAtOrXuanyuanRankPrizeData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"rank":          config.ParseValidator("int>0", "", false, nil, nil),
	"plunder_prize": config.ParseValidator("string", "", false, nil, nil),
	"show_prize":    config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *XuanyuanRankPrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *XuanyuanRankPrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *XuanyuanRankPrizeData) Encode() *shared_proto.XuanyuanRankPrizeDataProto {
	out := &shared_proto.XuanyuanRankPrizeDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Rank = config.U64ToI32(dAtA.Rank)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.Encode()
	}

	return out
}

func ArrayEncodeXuanyuanRankPrizeData(datas []*XuanyuanRankPrizeData) []*shared_proto.XuanyuanRankPrizeDataProto {

	out := make([]*shared_proto.XuanyuanRankPrizeDataProto, 0, len(datas))
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

func (dAtA *XuanyuanRankPrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.PlunderPrize = cOnFigS.GetPlunderPrize(pArSeR.Uint64("plunder_prize"))
	if dAtA.PlunderPrize == nil {
		return errors.Errorf("%s 配置的关联字段[plunder_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder_prize"), *pArSeR)
	}

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil {
		return errors.Errorf("%s 配置的关联字段[show_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetCombatScene(string) *scene.CombatScene
	GetPlunderPrize(uint64) *resdata.PlunderPrize
	GetPrize(int) *resdata.Prize
}
