// AUTO_GEN, DONT MODIFY!!!
package zhengwu

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

// start with ZhengWuCompleteData ----------------------------------

func LoadZhengWuCompleteData(gos *config.GameObjects) (map[uint64]*ZhengWuCompleteData, map[*ZhengWuCompleteData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ZhengWuCompleteDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ZhengWuCompleteData, len(lIsT))
	pArSeRmAp := make(map[*ZhengWuCompleteData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrZhengWuCompleteData) {
			continue
		}

		dAtA, err := NewZhengWuCompleteData(fIlEnAmE, pArSeR)
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

func SetRelatedZhengWuCompleteData(dAtAmAp map[*ZhengWuCompleteData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ZhengWuCompleteDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetZhengWuCompleteDataKeyArray(datas []*ZhengWuCompleteData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewZhengWuCompleteData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ZhengWuCompleteData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrZhengWuCompleteData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ZhengWuCompleteData{}

	dAtA.Quality = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("quality"))])
	if i, err := strconv.ParseInt(pArSeR.String("quality"), 10, 32); err == nil {
		dAtA.Quality = shared_proto.Quality(i)
	}

	dAtA.Cost = pArSeR.Uint64("cost")

	// calculate fields
	dAtA.Id = uint64(dAtA.Quality)

	return dAtA, nil
}

var vAlIdAtOrZhengWuCompleteData = map[string]*config.Validator{

	"quality": config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), nil),
	"cost":    config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *ZhengWuCompleteData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with ZhengWuData ----------------------------------

func LoadZhengWuData(gos *config.GameObjects) (map[uint64]*ZhengWuData, map[*ZhengWuData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ZhengWuDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ZhengWuData, len(lIsT))
	pArSeRmAp := make(map[*ZhengWuData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrZhengWuData) {
			continue
		}

		dAtA, err := NewZhengWuData(fIlEnAmE, pArSeR)
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

func SetRelatedZhengWuData(dAtAmAp map[*ZhengWuData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ZhengWuDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetZhengWuDataKeyArray(datas []*ZhengWuData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewZhengWuData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ZhengWuData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrZhengWuData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ZhengWuData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Quality = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("quality"))])
	if i, err := strconv.ParseInt(pArSeR.String("quality"), 10, 32); err == nil {
		dAtA.Quality = shared_proto.Quality(i)
	}

	// skip field: Cost
	// releated field: Prize
	dAtA.Duration, err = config.ParseDuration(pArSeR.String("duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("duration"), dAtA)
	}

	// skip field: Proto
	// skip field: ProtoBytes

	return dAtA, nil
}

var vAlIdAtOrZhengWuData = map[string]*config.Validator{

	"id":       config.ParseValidator("int>0", "", false, nil, nil),
	"name":     config.ParseValidator("string", "", false, nil, nil),
	"icon":     config.ParseValidator("string", "", false, nil, nil),
	"quality":  config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), nil),
	"prize":    config.ParseValidator("string", "", false, nil, nil),
	"duration": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ZhengWuData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ZhengWuData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ZhengWuData) Encode() *shared_proto.ZhengWuDataProto {
	out := &shared_proto.ZhengWuDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Icon = dAtA.Icon
	out.Quality = dAtA.Quality
	if dAtA.Cost != nil {
		out.Cost = config.U64ToI32(dAtA.Cost.Cost)
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	out.Duration = config.Duration2I32Seconds(dAtA.Duration)

	return out
}

func ArrayEncodeZhengWuData(datas []*ZhengWuData) []*shared_proto.ZhengWuDataProto {

	out := make([]*shared_proto.ZhengWuDataProto, 0, len(datas))
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

func (dAtA *ZhengWuData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with ZhengWuMiscData ----------------------------------

func LoadZhengWuMiscData(gos *config.GameObjects) (*ZhengWuMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.ZhengWuMiscDataPath
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

	dAtA, err := NewZhengWuMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedZhengWuMiscData(gos *config.GameObjects, dAtA *ZhengWuMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ZhengWuMiscDataPath
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

func NewZhengWuMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ZhengWuMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrZhengWuMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ZhengWuMiscData{}

	stringKeys = pArSeR.StringArray("auto_refresh_duration", "", false)
	dAtA.AutoRefreshDuration = make([]time.Duration, 0, len(stringKeys))
	for _, v := range stringKeys {
		obj, err := config.ParseDuration(v)
		if err != nil {
			return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[auto_refresh_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("auto_refresh_duration"), dAtA)
		}
		dAtA.AutoRefreshDuration = append(dAtA.AutoRefreshDuration, obj)
	}

	dAtA.RandomCount = pArSeR.Uint64("random_count")
	// releated field: FirstZhengWu

	return dAtA, nil
}

var vAlIdAtOrZhengWuMiscData = map[string]*config.Validator{

	"auto_refresh_duration": config.ParseValidator("string", "", true, nil, nil),
	"random_count":          config.ParseValidator("int>0", "", false, nil, nil),
	"first_zheng_wu":        config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ZhengWuMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ZhengWuMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ZhengWuMiscData) Encode() *shared_proto.ZhengWuMiscProto {
	out := &shared_proto.ZhengWuMiscProto{}
	out.AutoRefreshDuration = config.DurationArr2I32Seconds(dAtA.AutoRefreshDuration)

	return out
}

func ArrayEncodeZhengWuMiscData(datas []*ZhengWuMiscData) []*shared_proto.ZhengWuMiscProto {

	out := make([]*shared_proto.ZhengWuMiscProto, 0, len(datas))
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

func (dAtA *ZhengWuMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.FirstZhengWu = cOnFigS.GetZhengWuData(pArSeR.Uint64("first_zheng_wu"))
	if dAtA.FirstZhengWu == nil && pArSeR.Uint64("first_zheng_wu") != 0 {
		return errors.Errorf("%s 配置的关联字段[first_zheng_wu] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("first_zheng_wu"), *pArSeR)
	}

	return nil
}

// start with ZhengWuRandomData ----------------------------------

func LoadZhengWuRandomData(gos *config.GameObjects) (*ZhengWuRandomData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.ZhengWuRandomDataPath
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

	dAtA, err := NewZhengWuRandomData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedZhengWuRandomData(gos *config.GameObjects, dAtA *ZhengWuRandomData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ZhengWuRandomDataPath
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

func NewZhengWuRandomData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ZhengWuRandomData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrZhengWuRandomData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ZhengWuRandomData{}

	return dAtA, nil
}

var vAlIdAtOrZhengWuRandomData = map[string]*config.Validator{}

func (dAtA *ZhengWuRandomData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with ZhengWuRefreshData ----------------------------------

func LoadZhengWuRefreshData(gos *config.GameObjects) (map[uint64]*ZhengWuRefreshData, map[*ZhengWuRefreshData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ZhengWuRefreshDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ZhengWuRefreshData, len(lIsT))
	pArSeRmAp := make(map[*ZhengWuRefreshData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrZhengWuRefreshData) {
			continue
		}

		dAtA, err := NewZhengWuRefreshData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Times
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Times], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedZhengWuRefreshData(dAtAmAp map[*ZhengWuRefreshData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ZhengWuRefreshDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetZhengWuRefreshDataKeyArray(datas []*ZhengWuRefreshData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Times)
		}
	}

	return out
}

func NewZhengWuRefreshData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ZhengWuRefreshData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrZhengWuRefreshData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ZhengWuRefreshData{}

	dAtA.Times = pArSeR.Uint64("times")
	dAtA.Cost = 1000
	if pArSeR.KeyExist("cost") {
		dAtA.Cost = pArSeR.Uint64("cost")
	}

	// releated field: NewCost

	return dAtA, nil
}

var vAlIdAtOrZhengWuRefreshData = map[string]*config.Validator{

	"times":    config.ParseValidator("int>0", "", false, nil, nil),
	"cost":     config.ParseValidator("uint", "", false, nil, []string{"1000"}),
	"new_cost": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ZhengWuRefreshData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ZhengWuRefreshData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ZhengWuRefreshData) Encode() *shared_proto.ZhengWuRefreshDataProto {
	out := &shared_proto.ZhengWuRefreshDataProto{}
	out.Times = config.U64ToI32(dAtA.Times)
	out.Cost = config.U64ToI32(dAtA.Cost)
	if dAtA.NewCost != nil {
		out.NewCost = dAtA.NewCost.Encode()
	}

	return out
}

func ArrayEncodeZhengWuRefreshData(datas []*ZhengWuRefreshData) []*shared_proto.ZhengWuRefreshDataProto {

	out := make([]*shared_proto.ZhengWuRefreshDataProto, 0, len(datas))
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

func (dAtA *ZhengWuRefreshData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.NewCost = cOnFigS.GetCost(pArSeR.Int("new_cost"))
	if dAtA.NewCost == nil {
		return errors.Errorf("%s 配置的关联字段[new_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("new_cost"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetCost(int) *resdata.Cost
	GetPrize(int) *resdata.Prize
	GetZhengWuData(uint64) *ZhengWuData
}
