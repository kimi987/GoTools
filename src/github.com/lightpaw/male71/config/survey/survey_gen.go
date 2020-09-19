// AUTO_GEN, DONT MODIFY!!!
package survey

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
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

// start with SurveyData ----------------------------------

func LoadSurveyData(gos *config.GameObjects) (map[string]*SurveyData, map[*SurveyData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.SurveyDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*SurveyData, len(lIsT))
	pArSeRmAp := make(map[*SurveyData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrSurveyData) {
			continue
		}

		dAtA, err := NewSurveyData(fIlEnAmE, pArSeR)
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

func SetRelatedSurveyData(dAtAmAp map[*SurveyData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SurveyDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetSurveyDataKeyArray(datas []*SurveyData) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewSurveyData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SurveyData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSurveyData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SurveyData{}

	dAtA.Id = pArSeR.String("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Url = pArSeR.String("url")
	dAtA.Condition, err = data.NewUnlockCondition(fIlEnAmE, pArSeR)
	if err != nil {
		return nil, err
	}
	// releated field: Prize
	// skip field: PrizeProto

	return dAtA, nil
}

var vAlIdAtOrSurveyData = map[string]*config.Validator{

	"id":    config.ParseValidator("string", "", false, nil, nil),
	"name":  config.ParseValidator("string", "", false, nil, nil),
	"icon":  config.ParseValidator("string", "", false, nil, nil),
	"url":   config.ParseValidator("string", "", false, nil, nil),
	"prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *SurveyData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SurveyData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SurveyData) Encode() *shared_proto.SurveyDataProto {
	out := &shared_proto.SurveyDataProto{}
	out.Id = dAtA.Id
	out.Name = dAtA.Name
	out.Icon = dAtA.Icon
	out.Url = dAtA.Url
	if dAtA.Condition != nil {
		out.Condition = dAtA.Condition.Encode()
	}
	if dAtA.PrizeProto != nil {
		out.Prize = dAtA.PrizeProto
	}

	return out
}

func ArrayEncodeSurveyData(datas []*SurveyData) []*shared_proto.SurveyDataProto {

	out := make([]*shared_proto.SurveyDataProto, 0, len(datas))
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

func (dAtA *SurveyData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if err := dAtA.Condition.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS0); err != nil {
		return err
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetPrize(int) *resdata.Prize
}
