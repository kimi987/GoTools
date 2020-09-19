// AUTO_GEN, DONT MODIFY!!!
package question

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
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

// start with QuestionData ----------------------------------

func LoadQuestionData(gos *config.GameObjects) (map[uint64]*QuestionData, map[*QuestionData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.QuestionDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*QuestionData, len(lIsT))
	pArSeRmAp := make(map[*QuestionData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrQuestionData) {
			continue
		}

		dAtA, err := NewQuestionData(fIlEnAmE, pArSeR)
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

func SetRelatedQuestionData(dAtAmAp map[*QuestionData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.QuestionDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetQuestionDataKeyArray(datas []*QuestionData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewQuestionData(fIlEnAmE string, pArSeR *config.ObjectParser) (*QuestionData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrQuestionData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &QuestionData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Question = pArSeR.String("question")
	dAtA.RightAnswer = pArSeR.String("right_answer")
	dAtA.WrongAnswer = pArSeR.StringArray("wrong_answer", "", false)

	return dAtA, nil
}

var vAlIdAtOrQuestionData = map[string]*config.Validator{

	"id":           config.ParseValidator("int>0", "", false, nil, nil),
	"question":     config.ParseValidator("string>0", "", false, nil, nil),
	"right_answer": config.ParseValidator("string>0", "", false, nil, nil),
	"wrong_answer": config.ParseValidator("string>0", "", true, nil, nil),
}

func (dAtA *QuestionData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *QuestionData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *QuestionData) Encode() *shared_proto.QuestionProto {
	out := &shared_proto.QuestionProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Question = dAtA.Question
	out.RightAnswer = dAtA.RightAnswer
	out.WrongAnswer = dAtA.WrongAnswer

	return out
}

func ArrayEncodeQuestionData(datas []*QuestionData) []*shared_proto.QuestionProto {

	out := make([]*shared_proto.QuestionProto, 0, len(datas))
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

func (dAtA *QuestionData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with QuestionMiscData ----------------------------------

func LoadQuestionMiscData(gos *config.GameObjects) (*QuestionMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.QuestionMiscDataPath
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

	dAtA, err := NewQuestionMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedQuestionMiscData(gos *config.GameObjects, dAtA *QuestionMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.QuestionMiscDataPath
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

func NewQuestionMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*QuestionMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrQuestionMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &QuestionMiscData{}

	dAtA.MaxTimes = pArSeR.Uint64("max_times")
	dAtA.QuestionCount = pArSeR.Uint64("question_count")

	return dAtA, nil
}

var vAlIdAtOrQuestionMiscData = map[string]*config.Validator{

	"max_times":      config.ParseValidator("int>0", "", false, nil, nil),
	"question_count": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *QuestionMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *QuestionMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *QuestionMiscData) Encode() *shared_proto.QuestionMiscProto {
	out := &shared_proto.QuestionMiscProto{}
	out.MaxTimes = config.U64ToI32(dAtA.MaxTimes)
	out.QuestionCount = config.U64ToI32(dAtA.QuestionCount)

	return out
}

func ArrayEncodeQuestionMiscData(datas []*QuestionMiscData) []*shared_proto.QuestionMiscProto {

	out := make([]*shared_proto.QuestionMiscProto, 0, len(datas))
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

func (dAtA *QuestionMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with QuestionPrizeData ----------------------------------

func LoadQuestionPrizeData(gos *config.GameObjects) (map[uint64]*QuestionPrizeData, map[*QuestionPrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.QuestionPrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*QuestionPrizeData, len(lIsT))
	pArSeRmAp := make(map[*QuestionPrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrQuestionPrizeData) {
			continue
		}

		dAtA, err := NewQuestionPrizeData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Score
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Score], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedQuestionPrizeData(dAtAmAp map[*QuestionPrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.QuestionPrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetQuestionPrizeDataKeyArray(datas []*QuestionPrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Score)
		}
	}

	return out
}

func NewQuestionPrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*QuestionPrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrQuestionPrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &QuestionPrizeData{}

	dAtA.Score = pArSeR.Uint64("score")
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrQuestionPrizeData = map[string]*config.Validator{

	"score": config.ParseValidator("int", "", false, nil, nil),
	"prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *QuestionPrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *QuestionPrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *QuestionPrizeData) Encode() *shared_proto.QuestionPrizeProto {
	out := &shared_proto.QuestionPrizeProto{}
	out.Score = config.U64ToI32(dAtA.Score)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeQuestionPrizeData(datas []*QuestionPrizeData) []*shared_proto.QuestionPrizeProto {

	out := make([]*shared_proto.QuestionPrizeProto, 0, len(datas))
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

func (dAtA *QuestionPrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	return nil
}

// start with QuestionSayingData ----------------------------------

func LoadQuestionSayingData(gos *config.GameObjects) (map[uint64]*QuestionSayingData, map[*QuestionSayingData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.QuestionSayingDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*QuestionSayingData, len(lIsT))
	pArSeRmAp := make(map[*QuestionSayingData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrQuestionSayingData) {
			continue
		}

		dAtA, err := NewQuestionSayingData(fIlEnAmE, pArSeR)
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

func SetRelatedQuestionSayingData(dAtAmAp map[*QuestionSayingData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.QuestionSayingDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetQuestionSayingDataKeyArray(datas []*QuestionSayingData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewQuestionSayingData(fIlEnAmE string, pArSeR *config.ObjectParser) (*QuestionSayingData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrQuestionSayingData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &QuestionSayingData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Content = pArSeR.String("content")
	dAtA.Author = pArSeR.String("author")

	return dAtA, nil
}

var vAlIdAtOrQuestionSayingData = map[string]*config.Validator{

	"id":      config.ParseValidator("int>0", "", false, nil, nil),
	"content": config.ParseValidator("string>0", "", false, nil, nil),
	"author":  config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *QuestionSayingData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *QuestionSayingData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *QuestionSayingData) Encode() *shared_proto.QuestionSayingProto {
	out := &shared_proto.QuestionSayingProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Content = dAtA.Content
	out.Author = dAtA.Author

	return out
}

func ArrayEncodeQuestionSayingData(datas []*QuestionSayingData) []*shared_proto.QuestionSayingProto {

	out := make([]*shared_proto.QuestionSayingProto, 0, len(datas))
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

func (dAtA *QuestionSayingData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
	GetPrize(int) *resdata.Prize
}
