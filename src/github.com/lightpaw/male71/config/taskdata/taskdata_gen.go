// AUTO_GEN, DONT MODIFY!!!
package taskdata

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

// start with AchieveTaskData ----------------------------------

func LoadAchieveTaskData(gos *config.GameObjects) (map[uint64]*AchieveTaskData, map[*AchieveTaskData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.AchieveTaskDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*AchieveTaskData, len(lIsT))
	pArSeRmAp := make(map[*AchieveTaskData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrAchieveTaskData) {
			continue
		}

		dAtA, err := NewAchieveTaskData(fIlEnAmE, pArSeR)
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

func SetRelatedAchieveTaskData(dAtAmAp map[*AchieveTaskData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.AchieveTaskDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetAchieveTaskDataKeyArray(datas []*AchieveTaskData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewAchieveTaskData(fIlEnAmE string, pArSeR *config.ObjectParser) (*AchieveTaskData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrAchieveTaskData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &AchieveTaskData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Text = pArSeR.String("text")
	dAtA.Icon = pArSeR.String("icon")
	// releated field: Target
	// releated field: Prize
	dAtA.AchieveType = pArSeR.Uint64("achieve_type")
	dAtA.Next = pArSeR.Uint64("next")
	dAtA.Star = pArSeR.Uint64("star")
	dAtA.TotalStar = pArSeR.Uint64("total_star")
	dAtA.Quality = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("quality"))])
	if i, err := strconv.ParseInt(pArSeR.String("quality"), 10, 32); err == nil {
		dAtA.Quality = shared_proto.Quality(i)
	}

	dAtA.Order = 0
	if pArSeR.KeyExist("order") {
		dAtA.Order = pArSeR.Uint64("order")
	}

	// skip field: PrevTask

	return dAtA, nil
}

var vAlIdAtOrAchieveTaskData = map[string]*config.Validator{

	"id":           config.ParseValidator("int>0", "", false, nil, nil),
	"name":         config.ParseValidator("string", "", false, nil, nil),
	"text":         config.ParseValidator("string", "", false, nil, nil),
	"icon":         config.ParseValidator("string", "", false, nil, nil),
	"target":       config.ParseValidator("string", "", false, nil, nil),
	"prize":        config.ParseValidator("string", "", false, nil, nil),
	"achieve_type": config.ParseValidator("uint", "", false, nil, nil),
	"next":         config.ParseValidator("uint", "", false, nil, nil),
	"star":         config.ParseValidator("int>0", "", false, nil, nil),
	"total_star":   config.ParseValidator("int>0", "", false, nil, nil),
	"quality":      config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), nil),
	"order":        config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *AchieveTaskData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *AchieveTaskData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *AchieveTaskData) Encode() *shared_proto.TaskDataProto {
	out := &shared_proto.TaskDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Text = dAtA.Text
	out.Icon = dAtA.Icon
	if dAtA.Target != nil {
		out.Target = dAtA.Target.Encode()
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	out.AchieveType = config.U64ToI32(dAtA.AchieveType)
	out.Star = config.U64ToI32(dAtA.Star)
	out.TotalStar = config.U64ToI32(dAtA.TotalStar)
	out.Quality = dAtA.Quality
	out.Order = config.U64ToI32(dAtA.Order)
	if dAtA.PrevTask != nil {
		out.PrevTask = config.U64ToI32(dAtA.PrevTask.Id)
	}

	return out
}

func ArrayEncodeAchieveTaskData(datas []*AchieveTaskData) []*shared_proto.TaskDataProto {

	out := make([]*shared_proto.TaskDataProto, 0, len(datas))
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

func (dAtA *AchieveTaskData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Target = cOnFigS.GetTaskTargetData(pArSeR.Uint64("target"))
	if dAtA.Target == nil {
		return errors.Errorf("%s 配置的关联字段[target] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("target"), *pArSeR)
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with AchieveTaskStarPrizeData ----------------------------------

func LoadAchieveTaskStarPrizeData(gos *config.GameObjects) (map[uint64]*AchieveTaskStarPrizeData, map[*AchieveTaskStarPrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.AchieveTaskStarPrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*AchieveTaskStarPrizeData, len(lIsT))
	pArSeRmAp := make(map[*AchieveTaskStarPrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrAchieveTaskStarPrizeData) {
			continue
		}

		dAtA, err := NewAchieveTaskStarPrizeData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Star
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Star], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedAchieveTaskStarPrizeData(dAtAmAp map[*AchieveTaskStarPrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.AchieveTaskStarPrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetAchieveTaskStarPrizeDataKeyArray(datas []*AchieveTaskStarPrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Star)
		}
	}

	return out
}

func NewAchieveTaskStarPrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*AchieveTaskStarPrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrAchieveTaskStarPrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &AchieveTaskStarPrizeData{}

	dAtA.Star = pArSeR.Uint64("star")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Desc = pArSeR.String("desc")
	// releated field: Prize
	// releated field: Plunder

	return dAtA, nil
}

var vAlIdAtOrAchieveTaskStarPrizeData = map[string]*config.Validator{

	"star":    config.ParseValidator("int>0", "", false, nil, nil),
	"icon":    config.ParseValidator("string", "", false, nil, nil),
	"desc":    config.ParseValidator("string", "", false, nil, nil),
	"prize":   config.ParseValidator("string", "", false, nil, nil),
	"plunder": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *AchieveTaskStarPrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *AchieveTaskStarPrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *AchieveTaskStarPrizeData) Encode() *shared_proto.AchieveTaskStarPrizeProto {
	out := &shared_proto.AchieveTaskStarPrizeProto{}
	out.Star = config.U64ToI32(dAtA.Star)
	out.Icon = dAtA.Icon
	out.Desc = dAtA.Desc
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeAchieveTaskStarPrizeData(datas []*AchieveTaskStarPrizeData) []*shared_proto.AchieveTaskStarPrizeProto {

	out := make([]*shared_proto.AchieveTaskStarPrizeProto, 0, len(datas))
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

func (dAtA *AchieveTaskStarPrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.Plunder = cOnFigS.GetPlunder(pArSeR.Uint64("plunder"))
	if dAtA.Plunder == nil && pArSeR.Uint64("plunder") != 0 {
		return errors.Errorf("%s 配置的关联字段[plunder] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder"), *pArSeR)
	}

	return nil
}

// start with ActiveDegreePrizeData ----------------------------------

func LoadActiveDegreePrizeData(gos *config.GameObjects) (map[uint64]*ActiveDegreePrizeData, map[*ActiveDegreePrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ActiveDegreePrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ActiveDegreePrizeData, len(lIsT))
	pArSeRmAp := make(map[*ActiveDegreePrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrActiveDegreePrizeData) {
			continue
		}

		dAtA, err := NewActiveDegreePrizeData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Degree
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Degree], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedActiveDegreePrizeData(dAtAmAp map[*ActiveDegreePrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ActiveDegreePrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetActiveDegreePrizeDataKeyArray(datas []*ActiveDegreePrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Degree)
		}
	}

	return out
}

func NewActiveDegreePrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ActiveDegreePrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrActiveDegreePrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ActiveDegreePrizeData{}

	dAtA.Degree = pArSeR.Uint64("degree")
	// releated field: Prize
	// releated field: Plunder

	return dAtA, nil
}

var vAlIdAtOrActiveDegreePrizeData = map[string]*config.Validator{

	"degree":  config.ParseValidator("int>0", "", false, nil, nil),
	"prize":   config.ParseValidator("string", "", false, nil, nil),
	"plunder": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ActiveDegreePrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ActiveDegreePrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ActiveDegreePrizeData) Encode() *shared_proto.ActiveDegreePrizeProto {
	out := &shared_proto.ActiveDegreePrizeProto{}
	out.Degree = config.U64ToI32(dAtA.Degree)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeActiveDegreePrizeData(datas []*ActiveDegreePrizeData) []*shared_proto.ActiveDegreePrizeProto {

	out := make([]*shared_proto.ActiveDegreePrizeProto, 0, len(datas))
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

func (dAtA *ActiveDegreePrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.Plunder = cOnFigS.GetPlunder(pArSeR.Uint64("plunder"))
	if dAtA.Plunder == nil && pArSeR.Uint64("plunder") != 0 {
		return errors.Errorf("%s 配置的关联字段[plunder] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("plunder"), *pArSeR)
	}

	return nil
}

// start with ActiveDegreeTaskData ----------------------------------

func LoadActiveDegreeTaskData(gos *config.GameObjects) (map[uint64]*ActiveDegreeTaskData, map[*ActiveDegreeTaskData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ActiveDegreeTaskDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ActiveDegreeTaskData, len(lIsT))
	pArSeRmAp := make(map[*ActiveDegreeTaskData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrActiveDegreeTaskData) {
			continue
		}

		dAtA, err := NewActiveDegreeTaskData(fIlEnAmE, pArSeR)
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

func SetRelatedActiveDegreeTaskData(dAtAmAp map[*ActiveDegreeTaskData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ActiveDegreeTaskDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetActiveDegreeTaskDataKeyArray(datas []*ActiveDegreeTaskData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewActiveDegreeTaskData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ActiveDegreeTaskData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrActiveDegreeTaskData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ActiveDegreeTaskData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Text = pArSeR.String("text")
	dAtA.Icon = pArSeR.String("icon")
	// releated field: Target
	dAtA.AddDegree = pArSeR.Uint64("add_degree")

	return dAtA, nil
}

var vAlIdAtOrActiveDegreeTaskData = map[string]*config.Validator{

	"id":         config.ParseValidator("int>0", "", false, nil, nil),
	"name":       config.ParseValidator("string", "", false, nil, nil),
	"text":       config.ParseValidator("string", "", false, nil, nil),
	"icon":       config.ParseValidator("string", "", false, nil, nil),
	"target":     config.ParseValidator("string", "", false, nil, nil),
	"add_degree": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *ActiveDegreeTaskData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ActiveDegreeTaskData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ActiveDegreeTaskData) Encode() *shared_proto.TaskDataProto {
	out := &shared_proto.TaskDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Text = dAtA.Text
	out.Icon = dAtA.Icon
	if dAtA.Target != nil {
		out.Target = dAtA.Target.Encode()
	}
	out.AddDegree = config.U64ToI32(dAtA.AddDegree)

	return out
}

func ArrayEncodeActiveDegreeTaskData(datas []*ActiveDegreeTaskData) []*shared_proto.TaskDataProto {

	out := make([]*shared_proto.TaskDataProto, 0, len(datas))
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

func (dAtA *ActiveDegreeTaskData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Target = cOnFigS.GetTaskTargetData(pArSeR.Uint64("target"))
	if dAtA.Target == nil {
		return errors.Errorf("%s 配置的关联字段[target] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("target"), *pArSeR)
	}

	return nil
}

// start with ActivityTaskData ----------------------------------

func LoadActivityTaskData(gos *config.GameObjects) (map[uint64]*ActivityTaskData, map[*ActivityTaskData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ActivityTaskDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ActivityTaskData, len(lIsT))
	pArSeRmAp := make(map[*ActivityTaskData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrActivityTaskData) {
			continue
		}

		dAtA, err := NewActivityTaskData(fIlEnAmE, pArSeR)
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

func SetRelatedActivityTaskData(dAtAmAp map[*ActivityTaskData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ActivityTaskDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetActivityTaskDataKeyArray(datas []*ActivityTaskData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewActivityTaskData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ActivityTaskData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrActivityTaskData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ActivityTaskData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Text = pArSeR.String("text")
	// releated field: Target
	dAtA.Link = pArSeR.String("link")
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrActivityTaskData = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"name":   config.ParseValidator("string", "", false, nil, nil),
	"text":   config.ParseValidator("string", "", false, nil, nil),
	"target": config.ParseValidator("string", "", false, nil, nil),
	"link":   config.ParseValidator("string", "", false, nil, nil),
	"prize":  config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ActivityTaskData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ActivityTaskData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ActivityTaskData) Encode() *shared_proto.ActivityTaskDataProto {
	out := &shared_proto.ActivityTaskDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Text = dAtA.Text
	if dAtA.Target != nil {
		out.Target = dAtA.Target.Encode()
	}
	out.Link = dAtA.Link
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeActivityTaskData(datas []*ActivityTaskData) []*shared_proto.ActivityTaskDataProto {

	out := make([]*shared_proto.ActivityTaskDataProto, 0, len(datas))
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

func (dAtA *ActivityTaskData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Target = cOnFigS.GetTaskTargetData(pArSeR.Uint64("target"))
	if dAtA.Target == nil {
		return errors.Errorf("%s 配置的关联字段[target] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("target"), *pArSeR)
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with BaYeStageData ----------------------------------

func LoadBaYeStageData(gos *config.GameObjects) (map[uint64]*BaYeStageData, map[*BaYeStageData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BaYeStageDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BaYeStageData, len(lIsT))
	pArSeRmAp := make(map[*BaYeStageData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBaYeStageData) {
			continue
		}

		dAtA, err := NewBaYeStageData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Stage
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Stage], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedBaYeStageData(dAtAmAp map[*BaYeStageData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BaYeStageDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBaYeStageDataKeyArray(datas []*BaYeStageData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Stage)
		}
	}

	return out
}

func NewBaYeStageData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BaYeStageData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBaYeStageData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BaYeStageData{}

	dAtA.Stage = pArSeR.Uint64("stage")
	dAtA.StageName = ""
	if pArSeR.KeyExist("stage_name") {
		dAtA.StageName = pArSeR.String("stage_name")
	}

	dAtA.Name = pArSeR.String("name")
	// releated field: Tasks
	// releated field: Prize
	// skip field: Prev
	// skip field: Next

	return dAtA, nil
}

var vAlIdAtOrBaYeStageData = map[string]*config.Validator{

	"stage":      config.ParseValidator("int>0", "", false, nil, nil),
	"stage_name": config.ParseValidator("string", "", false, nil, []string{""}),
	"name":       config.ParseValidator("string>0", "", false, nil, nil),
	"tasks":      config.ParseValidator("int", "", true, nil, nil),
	"prize":      config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *BaYeStageData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BaYeStageData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BaYeStageData) Encode() *shared_proto.BaYeStageDataProto {
	out := &shared_proto.BaYeStageDataProto{}
	out.Stage = config.U64ToI32(dAtA.Stage)
	out.StageName = dAtA.StageName
	out.Name = dAtA.Name
	if dAtA.Tasks != nil {
		out.Tasks = config.U64a2I32a(GetBaYeTaskDataKeyArray(dAtA.Tasks))
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	if dAtA.Next != nil {
		out.Next = config.U64ToI32(dAtA.Next.Stage)
	}

	return out
}

func ArrayEncodeBaYeStageData(datas []*BaYeStageData) []*shared_proto.BaYeStageDataProto {

	out := make([]*shared_proto.BaYeStageDataProto, 0, len(datas))
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

func (dAtA *BaYeStageData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("tasks", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetBaYeTaskData(v)
		if obj != nil {
			dAtA.Tasks = append(dAtA.Tasks, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[tasks] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("tasks"), *pArSeR)
		}
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with BaYeTaskData ----------------------------------

func LoadBaYeTaskData(gos *config.GameObjects) (map[uint64]*BaYeTaskData, map[*BaYeTaskData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BaYeTaskDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BaYeTaskData, len(lIsT))
	pArSeRmAp := make(map[*BaYeTaskData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBaYeTaskData) {
			continue
		}

		dAtA, err := NewBaYeTaskData(fIlEnAmE, pArSeR)
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

func SetRelatedBaYeTaskData(dAtAmAp map[*BaYeTaskData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BaYeTaskDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBaYeTaskDataKeyArray(datas []*BaYeTaskData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBaYeTaskData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BaYeTaskData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBaYeTaskData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BaYeTaskData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Text = pArSeR.String("text")
	dAtA.Icon = pArSeR.String("icon")
	// releated field: Target
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrBaYeTaskData = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"name":   config.ParseValidator("string", "", false, nil, nil),
	"text":   config.ParseValidator("string", "", false, nil, nil),
	"icon":   config.ParseValidator("string", "", false, nil, nil),
	"target": config.ParseValidator("string", "", false, nil, nil),
	"prize":  config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *BaYeTaskData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BaYeTaskData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BaYeTaskData) Encode() *shared_proto.TaskDataProto {
	out := &shared_proto.TaskDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Text = dAtA.Text
	out.Icon = dAtA.Icon
	if dAtA.Target != nil {
		out.Target = dAtA.Target.Encode()
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeBaYeTaskData(datas []*BaYeTaskData) []*shared_proto.TaskDataProto {

	out := make([]*shared_proto.TaskDataProto, 0, len(datas))
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

func (dAtA *BaYeTaskData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Target = cOnFigS.GetTaskTargetData(pArSeR.Uint64("target"))
	if dAtA.Target == nil {
		return errors.Errorf("%s 配置的关联字段[target] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("target"), *pArSeR)
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with BranchTaskData ----------------------------------

func LoadBranchTaskData(gos *config.GameObjects) (map[uint64]*BranchTaskData, map[*BranchTaskData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BranchTaskDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BranchTaskData, len(lIsT))
	pArSeRmAp := make(map[*BranchTaskData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBranchTaskData) {
			continue
		}

		dAtA, err := NewBranchTaskData(fIlEnAmE, pArSeR)
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

func SetRelatedBranchTaskData(dAtAmAp map[*BranchTaskData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BranchTaskDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBranchTaskDataKeyArray(datas []*BranchTaskData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBranchTaskData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BranchTaskData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBranchTaskData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BranchTaskData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Next = pArSeR.Uint64("next")
	dAtA.Name = pArSeR.String("name")
	dAtA.Text = pArSeR.String("text")
	// releated field: Target
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrBranchTaskData = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"next":   config.ParseValidator("uint", "", false, nil, nil),
	"name":   config.ParseValidator("string", "", false, nil, nil),
	"text":   config.ParseValidator("string", "", false, nil, nil),
	"target": config.ParseValidator("string", "", false, nil, nil),
	"prize":  config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *BranchTaskData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BranchTaskData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BranchTaskData) Encode() *shared_proto.TaskDataProto {
	out := &shared_proto.TaskDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Text = dAtA.Text
	if dAtA.Target != nil {
		out.Target = dAtA.Target.Encode()
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeBranchTaskData(datas []*BranchTaskData) []*shared_proto.TaskDataProto {

	out := make([]*shared_proto.TaskDataProto, 0, len(datas))
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

func (dAtA *BranchTaskData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Target = cOnFigS.GetTaskTargetData(pArSeR.Uint64("target"))
	if dAtA.Target == nil {
		return errors.Errorf("%s 配置的关联字段[target] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("target"), *pArSeR)
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with BwzlPrizeData ----------------------------------

func LoadBwzlPrizeData(gos *config.GameObjects) (map[uint64]*BwzlPrizeData, map[*BwzlPrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BwzlPrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BwzlPrizeData, len(lIsT))
	pArSeRmAp := make(map[*BwzlPrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBwzlPrizeData) {
			continue
		}

		dAtA, err := NewBwzlPrizeData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.CollectPrizeTaskCount
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[CollectPrizeTaskCount], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedBwzlPrizeData(dAtAmAp map[*BwzlPrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BwzlPrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBwzlPrizeDataKeyArray(datas []*BwzlPrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.CollectPrizeTaskCount)
		}
	}

	return out
}

func NewBwzlPrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BwzlPrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBwzlPrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BwzlPrizeData{}

	dAtA.CollectPrizeTaskCount = pArSeR.Uint64("collect_prize_task_count")
	dAtA.Icon = pArSeR.String("icon")
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrBwzlPrizeData = map[string]*config.Validator{

	"collect_prize_task_count": config.ParseValidator("int>0", "", false, nil, nil),
	"icon":  config.ParseValidator("string", "", false, nil, nil),
	"prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *BwzlPrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BwzlPrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BwzlPrizeData) Encode() *shared_proto.BwzlPrizeDataProto {
	out := &shared_proto.BwzlPrizeDataProto{}
	out.CollectPrizeTaskCount = config.U64ToI32(dAtA.CollectPrizeTaskCount)
	out.Icon = dAtA.Icon
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeBwzlPrizeData(datas []*BwzlPrizeData) []*shared_proto.BwzlPrizeDataProto {

	out := make([]*shared_proto.BwzlPrizeDataProto, 0, len(datas))
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

func (dAtA *BwzlPrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with BwzlTaskData ----------------------------------

func LoadBwzlTaskData(gos *config.GameObjects) (map[uint64]*BwzlTaskData, map[*BwzlTaskData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.BwzlTaskDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*BwzlTaskData, len(lIsT))
	pArSeRmAp := make(map[*BwzlTaskData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrBwzlTaskData) {
			continue
		}

		dAtA, err := NewBwzlTaskData(fIlEnAmE, pArSeR)
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

func SetRelatedBwzlTaskData(dAtAmAp map[*BwzlTaskData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.BwzlTaskDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetBwzlTaskDataKeyArray(datas []*BwzlTaskData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewBwzlTaskData(fIlEnAmE string, pArSeR *config.ObjectParser) (*BwzlTaskData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrBwzlTaskData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &BwzlTaskData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Text = pArSeR.String("text")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Day = pArSeR.Uint64("day")
	// releated field: Target
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrBwzlTaskData = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"name":   config.ParseValidator("string", "", false, nil, nil),
	"text":   config.ParseValidator("string", "", false, nil, nil),
	"icon":   config.ParseValidator("string", "", false, nil, nil),
	"day":    config.ParseValidator("int>0", "", false, nil, nil),
	"target": config.ParseValidator("string", "", false, nil, nil),
	"prize":  config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *BwzlTaskData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *BwzlTaskData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *BwzlTaskData) Encode() *shared_proto.TaskDataProto {
	out := &shared_proto.TaskDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Text = dAtA.Text
	out.Icon = dAtA.Icon
	out.Day = config.U64ToI32(dAtA.Day)
	if dAtA.Target != nil {
		out.Target = dAtA.Target.Encode()
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeBwzlTaskData(datas []*BwzlTaskData) []*shared_proto.TaskDataProto {

	out := make([]*shared_proto.TaskDataProto, 0, len(datas))
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

func (dAtA *BwzlTaskData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Target = cOnFigS.GetTaskTargetData(pArSeR.Uint64("target"))
	if dAtA.Target == nil {
		return errors.Errorf("%s 配置的关联字段[target] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("target"), *pArSeR)
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with MainTaskData ----------------------------------

func LoadMainTaskData(gos *config.GameObjects) (map[uint64]*MainTaskData, map[*MainTaskData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.MainTaskDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*MainTaskData, len(lIsT))
	pArSeRmAp := make(map[*MainTaskData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrMainTaskData) {
			continue
		}

		dAtA, err := NewMainTaskData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Sequence
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Sequence], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedMainTaskData(dAtAmAp map[*MainTaskData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.MainTaskDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetMainTaskDataKeyArray(datas []*MainTaskData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Sequence)
		}
	}

	return out
}

func NewMainTaskData(fIlEnAmE string, pArSeR *config.ObjectParser) (*MainTaskData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrMainTaskData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &MainTaskData{}

	dAtA.Sequence = pArSeR.Uint64("sequence")
	dAtA.Name = pArSeR.String("name")
	dAtA.Text = pArSeR.String("text")
	// releated field: Target
	// releated field: Prize
	// releated field: BranchTask

	return dAtA, nil
}

var vAlIdAtOrMainTaskData = map[string]*config.Validator{

	"sequence":    config.ParseValidator("int>0", "", false, nil, nil),
	"name":        config.ParseValidator("string", "", false, nil, nil),
	"text":        config.ParseValidator("string", "", false, nil, nil),
	"target":      config.ParseValidator("string", "", false, nil, nil),
	"prize":       config.ParseValidator("string", "", false, nil, nil),
	"branch_task": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *MainTaskData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *MainTaskData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *MainTaskData) Encode() *shared_proto.TaskDataProto {
	out := &shared_proto.TaskDataProto{}
	out.Id = config.U64ToI32(dAtA.Sequence)
	out.Name = dAtA.Name
	out.Text = dAtA.Text
	if dAtA.Target != nil {
		out.Target = dAtA.Target.Encode()
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeMainTaskData(datas []*MainTaskData) []*shared_proto.TaskDataProto {

	out := make([]*shared_proto.TaskDataProto, 0, len(datas))
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

func (dAtA *MainTaskData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Target = cOnFigS.GetTaskTargetData(pArSeR.Uint64("target"))
	if dAtA.Target == nil {
		return errors.Errorf("%s 配置的关联字段[target] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("target"), *pArSeR)
	}

	dAtA.Prize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("branch_task", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetBranchTaskData(v)
		if obj != nil {
			dAtA.BranchTask = append(dAtA.BranchTask, obj)
		} else if v != 0 {
			return errors.Errorf("%s 配置的关联字段[branch_task] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("branch_task"), *pArSeR)
		}
	}

	return nil
}

// start with TaskBoxData ----------------------------------

func LoadTaskBoxData(gos *config.GameObjects) (map[uint64]*TaskBoxData, map[*TaskBoxData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TaskBoxDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TaskBoxData, len(lIsT))
	pArSeRmAp := make(map[*TaskBoxData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTaskBoxData) {
			continue
		}

		dAtA, err := NewTaskBoxData(fIlEnAmE, pArSeR)
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

func SetRelatedTaskBoxData(dAtAmAp map[*TaskBoxData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TaskBoxDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTaskBoxDataKeyArray(datas []*TaskBoxData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTaskBoxData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TaskBoxData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTaskBoxData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TaskBoxData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Count = pArSeR.Uint64("count")
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrTaskBoxData = map[string]*config.Validator{

	"id":    config.ParseValidator("int>0", "", false, nil, nil),
	"count": config.ParseValidator("int>0", "", false, nil, nil),
	"prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *TaskBoxData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TaskBoxData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TaskBoxData) Encode() *shared_proto.TaskBoxProto {
	out := &shared_proto.TaskBoxProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Count = config.U64ToI32(dAtA.Count)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeTaskBoxData(datas []*TaskBoxData) []*shared_proto.TaskBoxProto {

	out := make([]*shared_proto.TaskBoxProto, 0, len(datas))
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

func (dAtA *TaskBoxData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with TaskMiscData ----------------------------------

func LoadTaskMiscData(gos *config.GameObjects) (*TaskMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.TaskMiscDataPath
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

	dAtA, err := NewTaskMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedTaskMiscData(gos *config.GameObjects, dAtA *TaskMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TaskMiscDataPath
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

func NewTaskMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TaskMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTaskMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TaskMiscData{}

	dAtA.MaxShowAchieveCount = 3
	if pArSeR.KeyExist("max_show_achieve_count") {
		dAtA.MaxShowAchieveCount = pArSeR.Uint64("max_show_achieve_count")
	}

	dAtA.BwzlBgImg = pArSeR.String("bwzl_bg_img")
	if pArSeR.KeyExist("task_monster_arrive_offset") {
		dAtA.TaskMonsterArriveOffset, err = config.ParseDuration(pArSeR.String("task_monster_arrive_offset"))
	} else {
		dAtA.TaskMonsterArriveOffset, err = config.ParseDuration("3s")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[task_monster_arrive_offset] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("task_monster_arrive_offset"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrTaskMiscData = map[string]*config.Validator{

	"max_show_achieve_count":     config.ParseValidator("int>0", "", false, nil, []string{"3"}),
	"bwzl_bg_img":                config.ParseValidator("string", "", false, nil, nil),
	"task_monster_arrive_offset": config.ParseValidator("string", "", false, nil, []string{"3s"}),
}

func (dAtA *TaskMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TaskMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TaskMiscData) Encode() *shared_proto.TaskMiscDataProto {
	out := &shared_proto.TaskMiscDataProto{}
	out.MaxShowAchieveCount = config.U64ToI32(dAtA.MaxShowAchieveCount)
	out.BwzlBgImg = dAtA.BwzlBgImg

	return out
}

func ArrayEncodeTaskMiscData(datas []*TaskMiscData) []*shared_proto.TaskMiscDataProto {

	out := make([]*shared_proto.TaskMiscDataProto, 0, len(datas))
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

func (dAtA *TaskMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with TaskTargetData ----------------------------------

func LoadTaskTargetData(gos *config.GameObjects) (map[uint64]*TaskTargetData, map[*TaskTargetData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TaskTargetDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TaskTargetData, len(lIsT))
	pArSeRmAp := make(map[*TaskTargetData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTaskTargetData) {
			continue
		}

		dAtA, err := NewTaskTargetData(fIlEnAmE, pArSeR)
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

func SetRelatedTaskTargetData(dAtAmAp map[*TaskTargetData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TaskTargetDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTaskTargetDataKeyArray(datas []*TaskTargetData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTaskTargetData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TaskTargetData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTaskTargetData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TaskTargetData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Type = shared_proto.TaskTargetType(shared_proto.TaskTargetType_value[strings.ToUpper(pArSeR.String("type"))])
	if i, err := strconv.ParseInt(pArSeR.String("type"), 10, 32); err == nil {
		dAtA.Type = shared_proto.TaskTargetType(i)
	}

	dAtA.X = pArSeR.Uint64("x")
	dAtA.Y = pArSeR.String("y")
	dAtA.Z = pArSeR.String("z")
	dAtA.Daily = false
	if pArSeR.KeyExist("daily") {
		dAtA.Daily = pArSeR.Bool("daily")
	}

	// skip field: SubType
	// skip field: SubTypeId
	dAtA.DontUpdateProgress = pArSeR.Bool("dont_update_progress")
	// skip field: TotalProgress
	// skip field: ShowProgress
	// skip field: TechGroup
	// skip field: BuildingType
	// skip field: BuildingType2
	// skip field: CaptainRebirth
	// skip field: CaptainLevel
	// skip field: CaptainQuality
	// skip field: WearEquipmentCount
	// skip field: WearEquipmentLevel
	// skip field: WearEquipmentRefineLevel
	// skip field: WearEquipmentQuality
	// skip field: CaptainSoulQuality
	// skip field: CaptainSoulLevel
	// skip field: CaptainStar
	// skip field: TrainingLevel
	// skip field: ResBuildingType
	// skip field: RegionData
	// skip field: MonsterLevel
	// skip field: InvasionMonster
	// skip field: Tutor
	// skip field: PassDungeon
	// skip field: PassSecretTower
	// skip field: TentRegion
	// skip field: PassZhanjiangGuanqia
	// skip field: CaptainOfficial
	// skip field: CaptainAbilityExp
	// skip field: GemType
	// skip field: GemLevel
	// skip field: KillHomeNpcData
	// skip field: MultiLevelNpcType

	return dAtA, nil
}

var vAlIdAtOrTaskTargetData = map[string]*config.Validator{

	"id":                   config.ParseValidator("int>0", "", false, nil, nil),
	"type":                 config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.TaskTargetType_value, 0), nil),
	"x":                    config.ParseValidator("int>0", "", false, nil, nil),
	"y":                    config.ParseValidator("string", "", false, nil, nil),
	"z":                    config.ParseValidator("string", "", false, nil, nil),
	"daily":                config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"dont_update_progress": config.ParseValidator("bool", "", false, nil, nil),
}

func (dAtA *TaskTargetData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TaskTargetData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TaskTargetData) Encode() *shared_proto.TaskTargetProto {
	out := &shared_proto.TaskTargetProto{}
	out.Type = dAtA.Type
	out.SubType = config.U64ToI32(dAtA.SubType)
	out.SubTypeId = config.U64ToI32(dAtA.SubTypeId)
	out.TotalProgress = config.U64ToI32(dAtA.ShowProgress)

	return out
}

func ArrayEncodeTaskTargetData(datas []*TaskTargetData) []*shared_proto.TaskTargetProto {

	out := make([]*shared_proto.TaskTargetProto, 0, len(datas))
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

func (dAtA *TaskTargetData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with TitleData ----------------------------------

func LoadTitleData(gos *config.GameObjects) (map[uint64]*TitleData, map[*TitleData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TitleDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TitleData, len(lIsT))
	pArSeRmAp := make(map[*TitleData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTitleData) {
			continue
		}

		dAtA, err := NewTitleData(fIlEnAmE, pArSeR)
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

func SetRelatedTitleData(dAtAmAp map[*TitleData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TitleDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTitleDataKeyArray(datas []*TitleData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTitleData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TitleData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTitleData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TitleData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Image = ""
	if pArSeR.KeyExist("image") {
		dAtA.Image = pArSeR.String("image")
	}

	dAtA.Quality = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("quality"))])
	if i, err := strconv.ParseInt(pArSeR.String("quality"), 10, 32); err == nil {
		dAtA.Quality = shared_proto.Quality(i)
	}

	// releated field: SpriteStat
	// releated field: Task
	// releated field: TitleCost
	// skip field: TotalStat
	dAtA.PoltId = pArSeR.Uint64("polt_id")
	dAtA.CountryChangeNameVoteCount = pArSeR.Uint64("country_change_name_vote_count")

	return dAtA, nil
}

var vAlIdAtOrTitleData = map[string]*config.Validator{

	"id":                             config.ParseValidator("int>0", "", false, nil, nil),
	"name":                           config.ParseValidator("string", "", false, nil, nil),
	"desc":                           config.ParseValidator("string", "", false, nil, nil),
	"image":                          config.ParseValidator("string", "", false, nil, []string{""}),
	"quality":                        config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), nil),
	"sprite_stat":                    config.ParseValidator("string", "", false, nil, nil),
	"task":                           config.ParseValidator("string", "", true, nil, nil),
	"cost":                           config.ParseValidator("string", "", false, nil, nil),
	"polt_id":                        config.ParseValidator("int>0", "", false, nil, nil),
	"country_change_name_vote_count": config.ParseValidator("uint", "", false, nil, nil),
}

func (dAtA *TitleData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TitleData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TitleData) Encode() *shared_proto.TitleDataProto {
	out := &shared_proto.TitleDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.Image = dAtA.Image
	out.Quality = dAtA.Quality
	if dAtA.SpriteStat != nil {
		out.SpriteStat = dAtA.SpriteStat.Encode()
	}
	if dAtA.Task != nil {
		out.Task = config.U64a2I32a(GetTitleTaskDataKeyArray(dAtA.Task))
	}
	if dAtA.TitleCost != nil {
		out.TitleCost = dAtA.TitleCost.Encode()
	}
	if dAtA.TotalStat != nil {
		out.TotalStat = dAtA.TotalStat.Encode()
	}
	out.PoltId = config.U64ToI32(dAtA.PoltId)
	out.CountryChangeNameVoteCount = config.U64ToI32(dAtA.CountryChangeNameVoteCount)

	return out
}

func ArrayEncodeTitleData(datas []*TitleData) []*shared_proto.TitleDataProto {

	out := make([]*shared_proto.TitleDataProto, 0, len(datas))
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

func (dAtA *TitleData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	uint64Keys = pArSeR.Uint64Array("task", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetTitleTaskData(v)
		if obj != nil {
			dAtA.Task = append(dAtA.Task, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[task] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("task"), *pArSeR)
		}
	}

	dAtA.TitleCost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.TitleCost == nil && pArSeR.Int("cost") != 0 {
		return errors.Errorf("%s 配置的关联字段[cost] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	return nil
}

// start with TitleTaskData ----------------------------------

func LoadTitleTaskData(gos *config.GameObjects) (map[uint64]*TitleTaskData, map[*TitleTaskData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TitleTaskDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TitleTaskData, len(lIsT))
	pArSeRmAp := make(map[*TitleTaskData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTitleTaskData) {
			continue
		}

		dAtA, err := NewTitleTaskData(fIlEnAmE, pArSeR)
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

func SetRelatedTitleTaskData(dAtAmAp map[*TitleTaskData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TitleTaskDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTitleTaskDataKeyArray(datas []*TitleTaskData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTitleTaskData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TitleTaskData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTitleTaskData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TitleTaskData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	// releated field: Target

	return dAtA, nil
}

var vAlIdAtOrTitleTaskData = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"name":   config.ParseValidator("string", "", false, nil, nil),
	"target": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *TitleTaskData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TitleTaskData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TitleTaskData) Encode() *shared_proto.TitleTaskDataProto {
	out := &shared_proto.TitleTaskDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	if dAtA.Target != nil {
		out.Target = dAtA.Target.Encode()
	}

	return out
}

func ArrayEncodeTitleTaskData(datas []*TitleTaskData) []*shared_proto.TitleTaskDataProto {

	out := make([]*shared_proto.TitleTaskDataProto, 0, len(datas))
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

func (dAtA *TitleTaskData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Target = cOnFigS.GetTaskTargetData(pArSeR.Uint64("target"))
	if dAtA.Target == nil {
		return errors.Errorf("%s 配置的关联字段[target] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("target"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetBaYeTaskData(uint64) *BaYeTaskData
	GetBranchTaskData(uint64) *BranchTaskData
	GetCost(int) *resdata.Cost
	GetPlunder(uint64) *resdata.Plunder
	GetPrize(int) *resdata.Prize
	GetSpriteStat(uint64) *data.SpriteStat
	GetTaskTargetData(uint64) *TaskTargetData
	GetTitleTaskData(uint64) *TitleTaskData
}
