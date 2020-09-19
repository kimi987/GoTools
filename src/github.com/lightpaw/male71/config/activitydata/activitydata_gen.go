// AUTO_GEN, DONT MODIFY!!!
package activitydata

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/combine"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/taskdata"
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

// start with ActivityCollectionData ----------------------------------

func LoadActivityCollectionData(gos *config.GameObjects) (map[uint64]*ActivityCollectionData, map[*ActivityCollectionData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ActivityCollectionDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ActivityCollectionData, len(lIsT))
	pArSeRmAp := make(map[*ActivityCollectionData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrActivityCollectionData) {
			continue
		}

		dAtA, err := NewActivityCollectionData(fIlEnAmE, pArSeR)
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

func SetRelatedActivityCollectionData(dAtAmAp map[*ActivityCollectionData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ActivityCollectionDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetActivityCollectionDataKeyArray(datas []*ActivityCollectionData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewActivityCollectionData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ActivityCollectionData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrActivityCollectionData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ActivityCollectionData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.NameIcon = pArSeR.String("name_icon")
	dAtA.TimeShow = pArSeR.String("time_show")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.IconSelect = pArSeR.String("icon_select")
	dAtA.Image = pArSeR.String("image")
	dAtA.Sort = pArSeR.Uint64("sort")
	// releated field: Exchanges
	// releated field: TimeRule

	return dAtA, nil
}

var vAlIdAtOrActivityCollectionData = map[string]*config.Validator{

	"id":          config.ParseValidator("int>0", "", false, nil, nil),
	"name":        config.ParseValidator("string", "", false, nil, nil),
	"name_icon":   config.ParseValidator("string", "", false, nil, nil),
	"time_show":   config.ParseValidator("string", "", false, nil, nil),
	"desc":        config.ParseValidator("string", "", false, nil, nil),
	"icon":        config.ParseValidator("string", "", false, nil, nil),
	"icon_select": config.ParseValidator("string", "", false, nil, nil),
	"image":       config.ParseValidator("string", "", false, nil, nil),
	"sort":        config.ParseValidator("int>0", "", false, nil, nil),
	"exchanges":   config.ParseValidator("string", "", true, nil, nil),
	"time_rule":   config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ActivityCollectionData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ActivityCollectionData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ActivityCollectionData) Encode() *shared_proto.ActivityCollectionDataProto {
	out := &shared_proto.ActivityCollectionDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.NameIcon = dAtA.NameIcon
	out.TimeShow = dAtA.TimeShow
	out.Desc = dAtA.Desc
	out.Icon = dAtA.Icon
	out.IconSelect = dAtA.IconSelect
	out.Image = dAtA.Image
	out.Sort = config.U64ToI32(dAtA.Sort)
	if dAtA.Exchanges != nil {
		out.Exchanges = ArrayEncodeCollectionExchangeData(dAtA.Exchanges)
	}

	return out
}

func ArrayEncodeActivityCollectionData(datas []*ActivityCollectionData) []*shared_proto.ActivityCollectionDataProto {

	out := make([]*shared_proto.ActivityCollectionDataProto, 0, len(datas))
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

func (dAtA *ActivityCollectionData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	uint64Keys = pArSeR.Uint64Array("exchanges", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetCollectionExchangeData(v)
		if obj != nil {
			dAtA.Exchanges = append(dAtA.Exchanges, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[exchanges] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("exchanges"), *pArSeR)
		}
	}

	dAtA.TimeRule = cOnFigS.GetTimeRuleData(pArSeR.Uint64("time_rule"))
	if dAtA.TimeRule == nil {
		return errors.Errorf("%s 配置的关联字段[time_rule] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("time_rule"), *pArSeR)
	}

	return nil
}

// start with ActivityShowData ----------------------------------

func LoadActivityShowData(gos *config.GameObjects) (map[uint64]*ActivityShowData, map[*ActivityShowData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ActivityShowDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ActivityShowData, len(lIsT))
	pArSeRmAp := make(map[*ActivityShowData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrActivityShowData) {
			continue
		}

		dAtA, err := NewActivityShowData(fIlEnAmE, pArSeR)
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

func SetRelatedActivityShowData(dAtAmAp map[*ActivityShowData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ActivityShowDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetActivityShowDataKeyArray(datas []*ActivityShowData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewActivityShowData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ActivityShowData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrActivityShowData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ActivityShowData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.SpineId = pArSeR.Uint64("spine_id")
	dAtA.Name = pArSeR.String("name")
	dAtA.NameIcon = pArSeR.String("name_icon")
	dAtA.TimeShow = pArSeR.String("time_show")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.PrizeDesc = pArSeR.String("prize_desc")
	// releated field: TimeRule
	dAtA.Icon = pArSeR.String("icon")
	dAtA.IconSelect = pArSeR.String("icon_select")
	dAtA.ShowCountdown = false
	if pArSeR.KeyExist("show_countdown") {
		dAtA.ShowCountdown = pArSeR.Bool("show_countdown")
	}

	dAtA.Sort = pArSeR.Uint64("sort")
	dAtA.Image = pArSeR.String("image")
	dAtA.ImagePos = pArSeR.Uint64("image_pos")
	dAtA.LinkName = pArSeR.String("link_name")
	// releated field: LinkTaskTarget
	// skip field: LinkTargetType
	// skip field: LinkTargetSubType
	// skip field: LinkTargetSubTypeId

	return dAtA, nil
}

var vAlIdAtOrActivityShowData = map[string]*config.Validator{

	"id":               config.ParseValidator("int>0", "", false, nil, nil),
	"spine_id":         config.ParseValidator("uint", "", false, nil, nil),
	"name":             config.ParseValidator("string", "", false, nil, nil),
	"name_icon":        config.ParseValidator("string", "", false, nil, nil),
	"time_show":        config.ParseValidator("string", "", false, nil, nil),
	"desc":             config.ParseValidator("string", "", false, nil, nil),
	"prize_desc":       config.ParseValidator("string", "", false, nil, nil),
	"time_rule":        config.ParseValidator("string", "", false, nil, nil),
	"icon":             config.ParseValidator("string", "", false, nil, nil),
	"icon_select":      config.ParseValidator("string", "", false, nil, nil),
	"show_countdown":   config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"sort":             config.ParseValidator("int>0", "", false, nil, nil),
	"image":            config.ParseValidator("string", "", false, nil, nil),
	"image_pos":        config.ParseValidator("int>0", "", false, nil, nil),
	"link_name":        config.ParseValidator("string", "", false, nil, nil),
	"link_task_target": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *ActivityShowData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ActivityShowData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ActivityShowData) Encode() *shared_proto.ActivityShowDataProto {
	out := &shared_proto.ActivityShowDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.SpineId = config.U64ToI32(dAtA.SpineId)
	out.Name = dAtA.Name
	out.NameIcon = dAtA.NameIcon
	out.TimeShow = dAtA.TimeShow
	out.Desc = dAtA.Desc
	out.PrizeDesc = dAtA.PrizeDesc
	out.Icon = dAtA.Icon
	out.IconSelect = dAtA.IconSelect
	out.ShowCountdown = dAtA.ShowCountdown
	out.Sort = config.U64ToI32(dAtA.Sort)
	out.Image = dAtA.Image
	out.ImagePos = config.U64ToI32(dAtA.ImagePos)
	out.LinkName = dAtA.LinkName
	out.LinkTargetType = dAtA.LinkTargetType
	out.LinkTargetSubType = config.U64ToI32(dAtA.LinkTargetSubType)
	out.LinkTargetSubTypeId = config.U64ToI32(dAtA.LinkTargetSubTypeId)

	return out
}

func ArrayEncodeActivityShowData(datas []*ActivityShowData) []*shared_proto.ActivityShowDataProto {

	out := make([]*shared_proto.ActivityShowDataProto, 0, len(datas))
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

func (dAtA *ActivityShowData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.TimeRule = cOnFigS.GetTimeRuleData(pArSeR.Uint64("time_rule"))
	if dAtA.TimeRule == nil {
		return errors.Errorf("%s 配置的关联字段[time_rule] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("time_rule"), *pArSeR)
	}

	dAtA.LinkTaskTarget = cOnFigS.GetTaskTargetData(pArSeR.Uint64("link_task_target"))
	if dAtA.LinkTaskTarget == nil && pArSeR.Uint64("link_task_target") != 0 {
		return errors.Errorf("%s 配置的关联字段[link_task_target] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("link_task_target"), *pArSeR)
	}

	return nil
}

// start with ActivityTaskListModeData ----------------------------------

func LoadActivityTaskListModeData(gos *config.GameObjects) (map[uint64]*ActivityTaskListModeData, map[*ActivityTaskListModeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.ActivityTaskListModeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*ActivityTaskListModeData, len(lIsT))
	pArSeRmAp := make(map[*ActivityTaskListModeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrActivityTaskListModeData) {
			continue
		}

		dAtA, err := NewActivityTaskListModeData(fIlEnAmE, pArSeR)
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

func SetRelatedActivityTaskListModeData(dAtAmAp map[*ActivityTaskListModeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ActivityTaskListModeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetActivityTaskListModeDataKeyArray(datas []*ActivityTaskListModeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewActivityTaskListModeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ActivityTaskListModeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrActivityTaskListModeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ActivityTaskListModeData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.NameIcon = pArSeR.String("name_icon")
	dAtA.TimeShow = pArSeR.String("time_show")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.IconSelect = pArSeR.String("icon_select")
	dAtA.Image = pArSeR.String("image")
	dAtA.Sort = pArSeR.Uint64("sort")
	// releated field: TimeRule
	// releated field: Tasks

	return dAtA, nil
}

var vAlIdAtOrActivityTaskListModeData = map[string]*config.Validator{

	"id":          config.ParseValidator("int>0", "", false, nil, nil),
	"name":        config.ParseValidator("string", "", false, nil, nil),
	"name_icon":   config.ParseValidator("string", "", false, nil, nil),
	"time_show":   config.ParseValidator("string", "", false, nil, nil),
	"desc":        config.ParseValidator("string", "", false, nil, nil),
	"icon":        config.ParseValidator("string", "", false, nil, nil),
	"icon_select": config.ParseValidator("string", "", false, nil, nil),
	"image":       config.ParseValidator("string", "", false, nil, nil),
	"sort":        config.ParseValidator("int>0", "", false, nil, nil),
	"time_rule":   config.ParseValidator("string", "", false, nil, nil),
	"tasks":       config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *ActivityTaskListModeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ActivityTaskListModeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ActivityTaskListModeData) Encode() *shared_proto.ActivityTaskListModeDataProto {
	out := &shared_proto.ActivityTaskListModeDataProto{}
	out.Name = dAtA.Name
	out.NameIcon = dAtA.NameIcon
	out.TimeShow = dAtA.TimeShow
	out.Desc = dAtA.Desc
	out.Icon = dAtA.Icon
	out.IconSelect = dAtA.IconSelect
	out.Image = dAtA.Image
	out.Sort = config.U64ToI32(dAtA.Sort)
	if dAtA.Tasks != nil {
		out.Tasks = taskdata.ArrayEncodeActivityTaskData(dAtA.Tasks)
	}

	return out
}

func ArrayEncodeActivityTaskListModeData(datas []*ActivityTaskListModeData) []*shared_proto.ActivityTaskListModeDataProto {

	out := make([]*shared_proto.ActivityTaskListModeDataProto, 0, len(datas))
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

func (dAtA *ActivityTaskListModeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.TimeRule = cOnFigS.GetTimeRuleData(pArSeR.Uint64("time_rule"))
	if dAtA.TimeRule == nil {
		return errors.Errorf("%s 配置的关联字段[time_rule] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("time_rule"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("tasks", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetActivityTaskData(v)
		if obj != nil {
			dAtA.Tasks = append(dAtA.Tasks, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[tasks] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("tasks"), *pArSeR)
		}
	}

	return nil
}

// start with CollectionExchangeData ----------------------------------

func LoadCollectionExchangeData(gos *config.GameObjects) (map[uint64]*CollectionExchangeData, map[*CollectionExchangeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.CollectionExchangeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*CollectionExchangeData, len(lIsT))
	pArSeRmAp := make(map[*CollectionExchangeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrCollectionExchangeData) {
			continue
		}

		dAtA, err := NewCollectionExchangeData(fIlEnAmE, pArSeR)
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

func SetRelatedCollectionExchangeData(dAtAmAp map[*CollectionExchangeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.CollectionExchangeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetCollectionExchangeDataKeyArray(datas []*CollectionExchangeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewCollectionExchangeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*CollectionExchangeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrCollectionExchangeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &CollectionExchangeData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Combine
	dAtA.Limit = 0
	if pArSeR.KeyExist("limit") {
		dAtA.Limit = pArSeR.Uint64("limit")
	}

	return dAtA, nil
}

var vAlIdAtOrCollectionExchangeData = map[string]*config.Validator{

	"id":      config.ParseValidator("int>0", "", false, nil, nil),
	"combine": config.ParseValidator("string", "", false, nil, nil),
	"limit":   config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *CollectionExchangeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *CollectionExchangeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *CollectionExchangeData) Encode() *shared_proto.CollectionExchangeDataProto {
	out := &shared_proto.CollectionExchangeDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.Combine != nil {
		out.Combine = dAtA.Combine.Encode()
	}
	out.Limit = config.U64ToI32(dAtA.Limit)

	return out
}

func ArrayEncodeCollectionExchangeData(datas []*CollectionExchangeData) []*shared_proto.CollectionExchangeDataProto {

	out := make([]*shared_proto.CollectionExchangeDataProto, 0, len(datas))
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

func (dAtA *CollectionExchangeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Combine = cOnFigS.GetGoodsCombineData(pArSeR.Uint64("combine"))
	if dAtA.Combine == nil {
		return errors.Errorf("%s 配置的关联字段[combine] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("combine"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetActivityTaskData(uint64) *taskdata.ActivityTaskData
	GetCollectionExchangeData(uint64) *CollectionExchangeData
	GetGoodsCombineData(uint64) *combine.GoodsCombineData
	GetTaskTargetData(uint64) *taskdata.TaskTargetData
	GetTimeRuleData(uint64) *data.TimeRuleData
}
