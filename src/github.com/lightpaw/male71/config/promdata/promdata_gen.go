// AUTO_GEN, DONT MODIFY!!!
package promdata

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

// start with DailyBargainData ----------------------------------

func LoadDailyBargainData(gos *config.GameObjects) (map[uint64]*DailyBargainData, map[*DailyBargainData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.DailyBargainDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*DailyBargainData, len(lIsT))
	pArSeRmAp := make(map[*DailyBargainData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrDailyBargainData) {
			continue
		}

		dAtA, err := NewDailyBargainData(fIlEnAmE, pArSeR)
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

func SetRelatedDailyBargainData(dAtAmAp map[*DailyBargainData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.DailyBargainDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetDailyBargainDataKeyArray(datas []*DailyBargainData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewDailyBargainData(fIlEnAmE string, pArSeR *config.ObjectParser) (*DailyBargainData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrDailyBargainData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &DailyBargainData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.GiveYuanbao = 0
	if pArSeR.KeyExist("give_yuanbao") {
		dAtA.GiveYuanbao = pArSeR.Uint64("give_yuanbao")
	}

	dAtA.ShowYuanbao = 0
	if pArSeR.KeyExist("show_yuanbao") {
		dAtA.ShowYuanbao = pArSeR.Uint64("show_yuanbao")
	}

	dAtA.Limit = pArSeR.Uint64("limit")
	dAtA.ChargeAmount = pArSeR.Uint64("charge_amount")
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrDailyBargainData = map[string]*config.Validator{

	"id":            config.ParseValidator("int>0", "", false, nil, nil),
	"name":          config.ParseValidator("string", "", false, nil, nil),
	"give_yuanbao":  config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"show_yuanbao":  config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"limit":         config.ParseValidator("int>0", "", false, nil, nil),
	"charge_amount": config.ParseValidator("int>0", "", false, nil, nil),
	"prize":         config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *DailyBargainData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *DailyBargainData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *DailyBargainData) Encode() *shared_proto.DailyBargainDataProto {
	out := &shared_proto.DailyBargainDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.GiveYuanbao = config.U64ToI32(dAtA.GiveYuanbao)
	out.ShowYuanbao = config.U64ToI32(dAtA.ShowYuanbao)
	out.Limit = config.U64ToI32(dAtA.Limit)
	out.ChargeAmount = config.U64ToI32(dAtA.ChargeAmount)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeDailyBargainData(datas []*DailyBargainData) []*shared_proto.DailyBargainDataProto {

	out := make([]*shared_proto.DailyBargainDataProto, 0, len(datas))
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

func (dAtA *DailyBargainData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with DurationCardData ----------------------------------

func LoadDurationCardData(gos *config.GameObjects) (map[uint64]*DurationCardData, map[*DurationCardData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.DurationCardDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*DurationCardData, len(lIsT))
	pArSeRmAp := make(map[*DurationCardData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrDurationCardData) {
			continue
		}

		dAtA, err := NewDurationCardData(fIlEnAmE, pArSeR)
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

func SetRelatedDurationCardData(dAtAmAp map[*DurationCardData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.DurationCardDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetDurationCardDataKeyArray(datas []*DurationCardData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewDurationCardData(fIlEnAmE string, pArSeR *config.ObjectParser) (*DurationCardData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrDurationCardData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &DurationCardData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Duration, err = config.ParseDuration(pArSeR.String("duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("duration"), dAtA)
	}

	dAtA.ChargeAmount = pArSeR.Uint64("charge_amount")
	// releated field: Prize
	// releated field: DailyPrize
	dAtA.BeforePromptDuration, err = config.ParseDuration(pArSeR.String("before_prompt_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[before_prompt_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("before_prompt_duration"), dAtA)
	}

	return dAtA, nil
}

var vAlIdAtOrDurationCardData = map[string]*config.Validator{

	"id":                     config.ParseValidator("int>0", "", false, nil, nil),
	"name":                   config.ParseValidator("string", "", false, nil, nil),
	"icon":                   config.ParseValidator("string", "", false, nil, nil),
	"desc":                   config.ParseValidator("string", "", false, nil, nil),
	"duration":               config.ParseValidator("string", "", false, nil, nil),
	"charge_amount":          config.ParseValidator("int>0", "", false, nil, nil),
	"prize":                  config.ParseValidator("string", "", false, nil, nil),
	"daily_prize":            config.ParseValidator("string", "", false, nil, nil),
	"before_prompt_duration": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *DurationCardData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *DurationCardData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *DurationCardData) Encode() *shared_proto.DurationCardDataProto {
	out := &shared_proto.DurationCardDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Icon = dAtA.Icon
	out.Desc = dAtA.Desc
	out.Duration = config.Duration2I32Seconds(dAtA.Duration)
	out.ChargeAmount = config.U64ToI32(dAtA.ChargeAmount)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	if dAtA.DailyPrize != nil {
		out.DailyPrize = dAtA.DailyPrize.Encode()
	}
	out.BeforePromptDuration = config.Duration2I32Seconds(dAtA.BeforePromptDuration)

	return out
}

func ArrayEncodeDurationCardData(datas []*DurationCardData) []*shared_proto.DurationCardDataProto {

	out := make([]*shared_proto.DurationCardDataProto, 0, len(datas))
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

func (dAtA *DurationCardData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.DailyPrize = cOnFigS.GetPrize(pArSeR.Int("daily_prize"))
	if dAtA.DailyPrize == nil {
		return errors.Errorf("%s 配置的关联字段[daily_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("daily_prize"), *pArSeR)
	}

	return nil
}

// start with EventLimitGiftConfig ----------------------------------

func LoadEventLimitGiftConfig(gos *config.GameObjects) (*EventLimitGiftConfig, *config.ObjectParser, error) {
	fIlEnAmE := confpath.EventLimitGiftConfigPath
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

	dAtA, err := NewEventLimitGiftConfig(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedEventLimitGiftConfig(gos *config.GameObjects, dAtA *EventLimitGiftConfig, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EventLimitGiftConfigPath
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

func NewEventLimitGiftConfig(fIlEnAmE string, pArSeR *config.ObjectParser) (*EventLimitGiftConfig, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEventLimitGiftConfig)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EventLimitGiftConfig{}

	return dAtA, nil
}

var vAlIdAtOrEventLimitGiftConfig = map[string]*config.Validator{}

func (dAtA *EventLimitGiftConfig) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with EventLimitGiftData ----------------------------------

func LoadEventLimitGiftData(gos *config.GameObjects) (map[uint64]*EventLimitGiftData, map[*EventLimitGiftData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.EventLimitGiftDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*EventLimitGiftData, len(lIsT))
	pArSeRmAp := make(map[*EventLimitGiftData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrEventLimitGiftData) {
			continue
		}

		dAtA, err := NewEventLimitGiftData(fIlEnAmE, pArSeR)
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

func SetRelatedEventLimitGiftData(dAtAmAp map[*EventLimitGiftData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.EventLimitGiftDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetEventLimitGiftDataKeyArray(datas []*EventLimitGiftData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewEventLimitGiftData(fIlEnAmE string, pArSeR *config.ObjectParser) (*EventLimitGiftData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrEventLimitGiftData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &EventLimitGiftData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Image = pArSeR.String("image")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.YuanbaoPrice = pArSeR.Uint64("yuanbao_price")
	dAtA.OldPrice = pArSeR.Uint64("old_price")
	dAtA.SignIcon = pArSeR.String("sign_icon")
	dAtA.SignName = pArSeR.String("sign_name")
	dAtA.DiscountIcon = pArSeR.String("discount_icon")
	dAtA.Discount = pArSeR.String("discount")
	dAtA.Dianquan = pArSeR.Uint64("dianquan")
	dAtA.Priority = 0
	if pArSeR.KeyExist("priority") {
		dAtA.Priority = pArSeR.Uint64("priority")
	}

	if pArSeR.KeyExist("time_duration") {
		dAtA.TimeDuration, err = config.ParseDuration(pArSeR.String("time_duration"))
	} else {
		dAtA.TimeDuration, err = config.ParseDuration("1h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[time_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("time_duration"), dAtA)
	}

	dAtA.MinHeroLevel = 0
	if pArSeR.KeyExist("min_hero_level") {
		dAtA.MinHeroLevel = pArSeR.Uint64("min_hero_level")
	}

	dAtA.MaxHeroLevel = 0
	if pArSeR.KeyExist("max_hero_level") {
		dAtA.MaxHeroLevel = pArSeR.Uint64("max_hero_level")
	}

	dAtA.MinGuanfuLevel = 0
	if pArSeR.KeyExist("min_guanfu_level") {
		dAtA.MinGuanfuLevel = pArSeR.Uint64("min_guanfu_level")
	}

	dAtA.MaxGuanfuLevel = 0
	if pArSeR.KeyExist("max_guanfu_level") {
		dAtA.MaxGuanfuLevel = pArSeR.Uint64("max_guanfu_level")
	}

	dAtA.BuyLimit = 0
	if pArSeR.KeyExist("buy_limit") {
		dAtA.BuyLimit = pArSeR.Uint64("buy_limit")
	}

	// releated field: Prize
	dAtA.Condition = pArSeR.Uint64("condition")
	dAtA.ConditionValue = 0
	if pArSeR.KeyExist("condition_value") {
		dAtA.ConditionValue = pArSeR.Uint64("condition_value")
	}

	// releated field: ShowPrize
	dAtA.GuildEventPrizeId = 0
	if pArSeR.KeyExist("guild_event_prize_id") {
		dAtA.GuildEventPrizeId = pArSeR.Uint64("guild_event_prize_id")
	}

	return dAtA, nil
}

var vAlIdAtOrEventLimitGiftData = map[string]*config.Validator{

	"id":                   config.ParseValidator("int>0", "", false, nil, nil),
	"name":                 config.ParseValidator("string", "", false, nil, nil),
	"icon":                 config.ParseValidator("string", "", false, nil, nil),
	"image":                config.ParseValidator("string", "", false, nil, nil),
	"desc":                 config.ParseValidator("string", "", false, nil, nil),
	"yuanbao_price":        config.ParseValidator("int>0", "", false, nil, nil),
	"old_price":            config.ParseValidator("int>0", "", false, nil, nil),
	"sign_icon":            config.ParseValidator("string", "", false, nil, nil),
	"sign_name":            config.ParseValidator("string", "", false, nil, nil),
	"discount_icon":        config.ParseValidator("string", "", false, nil, nil),
	"discount":             config.ParseValidator("string", "", false, nil, nil),
	"dianquan":             config.ParseValidator("int>0", "", false, nil, nil),
	"priority":             config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"time_duration":        config.ParseValidator("string", "", false, nil, []string{"1h"}),
	"min_hero_level":       config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"max_hero_level":       config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"min_guanfu_level":     config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"max_guanfu_level":     config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"buy_limit":            config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"prize":                config.ParseValidator("string", "", false, nil, nil),
	"condition":            config.ParseValidator("int>0", "", false, nil, nil),
	"condition_value":      config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"show_prize":           config.ParseValidator("string", "", false, nil, nil),
	"guild_event_prize_id": config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *EventLimitGiftData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *EventLimitGiftData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *EventLimitGiftData) Encode() *shared_proto.EventLimitGiftDataProto {
	out := &shared_proto.EventLimitGiftDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Icon = dAtA.Icon
	out.Image = dAtA.Image
	out.Desc = dAtA.Desc
	out.YuanbaoPrice = config.U64ToI32(dAtA.YuanbaoPrice)
	out.OldPrice = config.U64ToI32(dAtA.OldPrice)
	out.SignIcon = dAtA.SignIcon
	out.SignName = dAtA.SignName
	out.DiscountIcon = dAtA.DiscountIcon
	out.Discount = dAtA.Discount
	out.Dianquan = config.U64ToI32(dAtA.Dianquan)
	out.Priority = config.U64ToI32(dAtA.Priority)
	out.BuyLimit = config.U64ToI32(dAtA.BuyLimit)
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.Encode()
	}
	out.GuildEventPrizeId = config.U64ToI32(dAtA.GuildEventPrizeId)

	return out
}

func ArrayEncodeEventLimitGiftData(datas []*EventLimitGiftData) []*shared_proto.EventLimitGiftDataProto {

	out := make([]*shared_proto.EventLimitGiftDataProto, 0, len(datas))
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

func (dAtA *EventLimitGiftData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil {
		return errors.Errorf("%s 配置的关联字段[show_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
	}

	return nil
}

// start with FreeGiftData ----------------------------------

func LoadFreeGiftData(gos *config.GameObjects) (map[uint64]*FreeGiftData, map[*FreeGiftData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.FreeGiftDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*FreeGiftData, len(lIsT))
	pArSeRmAp := make(map[*FreeGiftData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrFreeGiftData) {
			continue
		}

		dAtA, err := NewFreeGiftData(fIlEnAmE, pArSeR)
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

func SetRelatedFreeGiftData(dAtAmAp map[*FreeGiftData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.FreeGiftDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetFreeGiftDataKeyArray(datas []*FreeGiftData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewFreeGiftData(fIlEnAmE string, pArSeR *config.ObjectParser) (*FreeGiftData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrFreeGiftData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &FreeGiftData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.GiftType = shared_proto.GiftType(shared_proto.GiftType_value[strings.ToUpper(pArSeR.String("gift_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("gift_type"), 10, 32); err == nil {
		dAtA.GiftType = shared_proto.GiftType(i)
	}

	dAtA.Daily = false
	if pArSeR.KeyExist("daily") {
		dAtA.Daily = pArSeR.Bool("daily")
	}

	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrFreeGiftData = map[string]*config.Validator{

	"id":        config.ParseValidator("int>0", "", false, nil, nil),
	"name":      config.ParseValidator("string", "", false, nil, nil),
	"gift_type": config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.GiftType_value, 0), nil),
	"daily":     config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"prize":     config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *FreeGiftData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *FreeGiftData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *FreeGiftData) Encode() *shared_proto.FreeGiftDataProto {
	out := &shared_proto.FreeGiftDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.GiftType = dAtA.GiftType
	out.Daily = dAtA.Daily
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeFreeGiftData(datas []*FreeGiftData) []*shared_proto.FreeGiftDataProto {

	out := make([]*shared_proto.FreeGiftDataProto, 0, len(datas))
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

func (dAtA *FreeGiftData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with HeroLevelFundData ----------------------------------

func LoadHeroLevelFundData(gos *config.GameObjects) (map[uint64]*HeroLevelFundData, map[*HeroLevelFundData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.HeroLevelFundDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*HeroLevelFundData, len(lIsT))
	pArSeRmAp := make(map[*HeroLevelFundData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrHeroLevelFundData) {
			continue
		}

		dAtA, err := NewHeroLevelFundData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Level
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Level], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedHeroLevelFundData(dAtAmAp map[*HeroLevelFundData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.HeroLevelFundDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetHeroLevelFundDataKeyArray(datas []*HeroLevelFundData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewHeroLevelFundData(fIlEnAmE string, pArSeR *config.ObjectParser) (*HeroLevelFundData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrHeroLevelFundData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &HeroLevelFundData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Rebate = pArSeR.Uint64("rebate")
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrHeroLevelFundData = map[string]*config.Validator{

	"level":  config.ParseValidator("int>0", "", false, nil, nil),
	"rebate": config.ParseValidator("int>0", "", false, nil, nil),
	"prize":  config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *HeroLevelFundData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *HeroLevelFundData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *HeroLevelFundData) Encode() *shared_proto.HeroLevelFundDataProto {
	out := &shared_proto.HeroLevelFundDataProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Rebate = config.U64ToI32(dAtA.Rebate)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeHeroLevelFundData(datas []*HeroLevelFundData) []*shared_proto.HeroLevelFundDataProto {

	out := make([]*shared_proto.HeroLevelFundDataProto, 0, len(datas))
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

func (dAtA *HeroLevelFundData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with LoginDayData ----------------------------------

func LoadLoginDayData(gos *config.GameObjects) (map[uint64]*LoginDayData, map[*LoginDayData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.LoginDayDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*LoginDayData, len(lIsT))
	pArSeRmAp := make(map[*LoginDayData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrLoginDayData) {
			continue
		}

		dAtA, err := NewLoginDayData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Day
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Day], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedLoginDayData(dAtAmAp map[*LoginDayData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.LoginDayDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetLoginDayDataKeyArray(datas []*LoginDayData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Day)
		}
	}

	return out
}

func NewLoginDayData(fIlEnAmE string, pArSeR *config.ObjectParser) (*LoginDayData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrLoginDayData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &LoginDayData{}

	dAtA.Day = pArSeR.Uint64("day")
	// releated field: Prize

	return dAtA, nil
}

var vAlIdAtOrLoginDayData = map[string]*config.Validator{

	"day":   config.ParseValidator("int>0", "", false, nil, nil),
	"prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *LoginDayData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *LoginDayData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *LoginDayData) Encode() *shared_proto.LoginDayDataProto {
	out := &shared_proto.LoginDayDataProto{}
	out.Day = config.U64ToI32(dAtA.Day)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}

	return out
}

func ArrayEncodeLoginDayData(datas []*LoginDayData) []*shared_proto.LoginDayDataProto {

	out := make([]*shared_proto.LoginDayDataProto, 0, len(datas))
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

func (dAtA *LoginDayData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with PromotionMiscData ----------------------------------

func LoadPromotionMiscData(gos *config.GameObjects) (*PromotionMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.PromotionMiscDataPath
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

	dAtA, err := NewPromotionMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedPromotionMiscData(gos *config.GameObjects, dAtA *PromotionMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.PromotionMiscDataPath
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

func NewPromotionMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*PromotionMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrPromotionMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &PromotionMiscData{}

	// releated field: HeroLevelFundCost

	return dAtA, nil
}

var vAlIdAtOrPromotionMiscData = map[string]*config.Validator{

	"hero_level_fund_cost": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *PromotionMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *PromotionMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *PromotionMiscData) Encode() *shared_proto.PromotionMiscDataProto {
	out := &shared_proto.PromotionMiscDataProto{}
	if dAtA.HeroLevelFundCost != nil {
		out.HeroLevelFundCost = dAtA.HeroLevelFundCost.Encode()
	}

	return out
}

func ArrayEncodePromotionMiscData(datas []*PromotionMiscData) []*shared_proto.PromotionMiscDataProto {

	out := make([]*shared_proto.PromotionMiscDataProto, 0, len(datas))
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

func (dAtA *PromotionMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.HeroLevelFundCost = cOnFigS.GetCost(pArSeR.Int("hero_level_fund_cost"))
	if dAtA.HeroLevelFundCost == nil {
		return errors.Errorf("%s 配置的关联字段[hero_level_fund_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("hero_level_fund_cost"), *pArSeR)
	}

	return nil
}

// start with SpCollectionData ----------------------------------

func LoadSpCollectionData(gos *config.GameObjects) (map[uint64]*SpCollectionData, map[*SpCollectionData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.SpCollectionDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*SpCollectionData, len(lIsT))
	pArSeRmAp := make(map[*SpCollectionData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrSpCollectionData) {
			continue
		}

		dAtA, err := NewSpCollectionData(fIlEnAmE, pArSeR)
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

func SetRelatedSpCollectionData(dAtAmAp map[*SpCollectionData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.SpCollectionDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetSpCollectionDataKeyArray(datas []*SpCollectionData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewSpCollectionData(fIlEnAmE string, pArSeR *config.ObjectParser) (*SpCollectionData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrSpCollectionData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &SpCollectionData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.TimeShow = pArSeR.String("time_show")
	dAtA.StartDuration, err = config.ParseDuration(pArSeR.String("start_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[start_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("start_duration"), dAtA)
	}

	dAtA.EndDuration, err = config.ParseDuration(pArSeR.String("end_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[end_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("end_duration"), dAtA)
	}

	dAtA.Sp = pArSeR.Uint64("sp")
	dAtA.RepairVip = pArSeR.Uint64("repair_vip")
	// releated field: SpPrize
	// releated field: RepairCost

	return dAtA, nil
}

var vAlIdAtOrSpCollectionData = map[string]*config.Validator{

	"id":             config.ParseValidator("int>0", "", false, nil, nil),
	"name":           config.ParseValidator("string", "", false, nil, nil),
	"icon":           config.ParseValidator("string", "", false, nil, nil),
	"time_show":      config.ParseValidator("string", "", false, nil, nil),
	"start_duration": config.ParseValidator("string", "", false, nil, nil),
	"end_duration":   config.ParseValidator("string", "", false, nil, nil),
	"sp":             config.ParseValidator("int>0", "", false, nil, nil),
	"repair_vip":     config.ParseValidator("int>0", "", false, nil, nil),
	"prize":          config.ParseValidator("string", "", false, nil, nil),
	"repair_cost":    config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *SpCollectionData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *SpCollectionData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *SpCollectionData) Encode() *shared_proto.SpCollectionDataProto {
	out := &shared_proto.SpCollectionDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Icon = dAtA.Icon
	out.TimeShow = dAtA.TimeShow
	out.StartDuration = config.Duration2I32Seconds(dAtA.StartDuration)
	out.EndDuration = config.Duration2I32Seconds(dAtA.EndDuration)
	out.Sp = config.U64ToI32(dAtA.Sp)
	out.RepairVip = config.U64ToI32(dAtA.RepairVip)
	if dAtA.SpPrize != nil {
		out.SpPrize = dAtA.SpPrize.Encode()
	}
	if dAtA.RepairCost != nil {
		out.RepairCost = dAtA.RepairCost.Encode()
	}

	return out
}

func ArrayEncodeSpCollectionData(datas []*SpCollectionData) []*shared_proto.SpCollectionDataProto {

	out := make([]*shared_proto.SpCollectionDataProto, 0, len(datas))
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

func (dAtA *SpCollectionData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.SpPrize = cOnFigS.GetPrize(pArSeR.Int("prize"))
	if dAtA.SpPrize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	dAtA.RepairCost = cOnFigS.GetCost(pArSeR.Int("repair_cost"))
	if dAtA.RepairCost == nil {
		return errors.Errorf("%s 配置的关联字段[repair_cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("repair_cost"), *pArSeR)
	}

	return nil
}

// start with TimeLimitGiftData ----------------------------------

func LoadTimeLimitGiftData(gos *config.GameObjects) (map[uint64]*TimeLimitGiftData, map[*TimeLimitGiftData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TimeLimitGiftDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TimeLimitGiftData, len(lIsT))
	pArSeRmAp := make(map[*TimeLimitGiftData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTimeLimitGiftData) {
			continue
		}

		dAtA, err := NewTimeLimitGiftData(fIlEnAmE, pArSeR)
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

func SetRelatedTimeLimitGiftData(dAtAmAp map[*TimeLimitGiftData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TimeLimitGiftDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTimeLimitGiftDataKeyArray(datas []*TimeLimitGiftData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTimeLimitGiftData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TimeLimitGiftData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTimeLimitGiftData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TimeLimitGiftData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Image = pArSeR.String("image")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.YuanbaoPrice = pArSeR.Uint64("yuanbao_price")
	dAtA.OldPrice = pArSeR.Uint64("old_price")
	dAtA.SignIcon = pArSeR.String("sign_icon")
	dAtA.SignName = pArSeR.String("sign_name")
	dAtA.DiscountIcon = pArSeR.String("discount_icon")
	dAtA.Discount = pArSeR.String("discount")
	dAtA.Dianquan = pArSeR.Uint64("dianquan")
	dAtA.Priority = 0
	if pArSeR.KeyExist("priority") {
		dAtA.Priority = pArSeR.Uint64("priority")
	}

	// releated field: Prize
	dAtA.BuyLimit = 0
	if pArSeR.KeyExist("buy_limit") {
		dAtA.BuyLimit = pArSeR.Uint64("buy_limit")
	}

	dAtA.Sort = 0
	if pArSeR.KeyExist("sort") {
		dAtA.Sort = pArSeR.Uint64("sort")
	}

	// releated field: ShowPrize
	dAtA.GuildEventPrizeId = 0
	if pArSeR.KeyExist("guild_event_prize_id") {
		dAtA.GuildEventPrizeId = pArSeR.Uint64("guild_event_prize_id")
	}

	return dAtA, nil
}

var vAlIdAtOrTimeLimitGiftData = map[string]*config.Validator{

	"id":                   config.ParseValidator("int>0", "", false, nil, nil),
	"name":                 config.ParseValidator("string", "", false, nil, nil),
	"icon":                 config.ParseValidator("string", "", false, nil, nil),
	"image":                config.ParseValidator("string", "", false, nil, nil),
	"desc":                 config.ParseValidator("string", "", false, nil, nil),
	"yuanbao_price":        config.ParseValidator("int>0", "", false, nil, nil),
	"old_price":            config.ParseValidator("int>0", "", false, nil, nil),
	"sign_icon":            config.ParseValidator("string", "", false, nil, nil),
	"sign_name":            config.ParseValidator("string", "", false, nil, nil),
	"discount_icon":        config.ParseValidator("string", "", false, nil, nil),
	"discount":             config.ParseValidator("string", "", false, nil, nil),
	"dianquan":             config.ParseValidator("int>0", "", false, nil, nil),
	"priority":             config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"prize":                config.ParseValidator("string", "", false, nil, nil),
	"buy_limit":            config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"sort":                 config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"show_prize":           config.ParseValidator("string", "", false, nil, nil),
	"guild_event_prize_id": config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *TimeLimitGiftData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TimeLimitGiftData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TimeLimitGiftData) Encode() *shared_proto.TimeLimitGiftDataProto {
	out := &shared_proto.TimeLimitGiftDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Icon = dAtA.Icon
	out.Image = dAtA.Image
	out.Desc = dAtA.Desc
	out.YuanbaoPrice = config.U64ToI32(dAtA.YuanbaoPrice)
	out.OldPrice = config.U64ToI32(dAtA.OldPrice)
	out.SignIcon = dAtA.SignIcon
	out.SignName = dAtA.SignName
	out.DiscountIcon = dAtA.DiscountIcon
	out.Discount = dAtA.Discount
	out.Dianquan = config.U64ToI32(dAtA.Dianquan)
	out.Priority = config.U64ToI32(dAtA.Priority)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	out.BuyLimit = config.U64ToI32(dAtA.BuyLimit)
	out.Sort = config.U64ToI32(dAtA.Sort)
	if dAtA.ShowPrize != nil {
		out.ShowPrize = dAtA.ShowPrize.Encode()
	}
	out.GuildEventPrizeId = config.U64ToI32(dAtA.GuildEventPrizeId)

	return out
}

func ArrayEncodeTimeLimitGiftData(datas []*TimeLimitGiftData) []*shared_proto.TimeLimitGiftDataProto {

	out := make([]*shared_proto.TimeLimitGiftDataProto, 0, len(datas))
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

func (dAtA *TimeLimitGiftData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.ShowPrize = cOnFigS.GetPrize(pArSeR.Int("show_prize"))
	if dAtA.ShowPrize == nil {
		return errors.Errorf("%s 配置的关联字段[show_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("show_prize"), *pArSeR)
	}

	return nil
}

// start with TimeLimitGiftGroupData ----------------------------------

func LoadTimeLimitGiftGroupData(gos *config.GameObjects) (map[uint64]*TimeLimitGiftGroupData, map[*TimeLimitGiftGroupData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.TimeLimitGiftGroupDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*TimeLimitGiftGroupData, len(lIsT))
	pArSeRmAp := make(map[*TimeLimitGiftGroupData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrTimeLimitGiftGroupData) {
			continue
		}

		dAtA, err := NewTimeLimitGiftGroupData(fIlEnAmE, pArSeR)
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

func SetRelatedTimeLimitGiftGroupData(dAtAmAp map[*TimeLimitGiftGroupData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.TimeLimitGiftGroupDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetTimeLimitGiftGroupDataKeyArray(datas []*TimeLimitGiftGroupData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewTimeLimitGiftGroupData(fIlEnAmE string, pArSeR *config.ObjectParser) (*TimeLimitGiftGroupData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrTimeLimitGiftGroupData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &TimeLimitGiftGroupData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: TimeRule
	// releated field: Gifts
	dAtA.MinHeroLevel = 0
	if pArSeR.KeyExist("min_hero_level") {
		dAtA.MinHeroLevel = pArSeR.Uint64("min_hero_level")
	}

	dAtA.MaxHeroLevel = 0
	if pArSeR.KeyExist("max_hero_level") {
		dAtA.MaxHeroLevel = pArSeR.Uint64("max_hero_level")
	}

	dAtA.MinGuanfuLevel = 0
	if pArSeR.KeyExist("min_guanfu_level") {
		dAtA.MinGuanfuLevel = pArSeR.Uint64("min_guanfu_level")
	}

	dAtA.MaxGuanfuLevel = 0
	if pArSeR.KeyExist("max_guanfu_level") {
		dAtA.MaxGuanfuLevel = pArSeR.Uint64("max_guanfu_level")
	}

	return dAtA, nil
}

var vAlIdAtOrTimeLimitGiftGroupData = map[string]*config.Validator{

	"id":               config.ParseValidator("int>0", "", false, nil, nil),
	"time_rule":        config.ParseValidator("string", "", false, nil, nil),
	"gifts":            config.ParseValidator("uint,notAllNil,duplicate", "", true, nil, nil),
	"min_hero_level":   config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"max_hero_level":   config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"min_guanfu_level": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"max_guanfu_level": config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *TimeLimitGiftGroupData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *TimeLimitGiftGroupData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *TimeLimitGiftGroupData) Encode() *shared_proto.TimeLimitGiftGroupDataProto {
	out := &shared_proto.TimeLimitGiftGroupDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.Gifts != nil {
		out.GiftIds = config.U64a2I32a(GetTimeLimitGiftDataKeyArray(dAtA.Gifts))
	}
	out.MinHeroLevel = config.U64ToI32(dAtA.MinHeroLevel)
	out.MaxHeroLevel = config.U64ToI32(dAtA.MaxHeroLevel)
	out.MinGuanfuLevel = config.U64ToI32(dAtA.MinGuanfuLevel)
	out.MaxGuanfuLevel = config.U64ToI32(dAtA.MaxGuanfuLevel)

	return out
}

func ArrayEncodeTimeLimitGiftGroupData(datas []*TimeLimitGiftGroupData) []*shared_proto.TimeLimitGiftGroupDataProto {

	out := make([]*shared_proto.TimeLimitGiftGroupDataProto, 0, len(datas))
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

func (dAtA *TimeLimitGiftGroupData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	uint64Keys = pArSeR.Uint64Array("gifts", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetTimeLimitGiftData(v)
		if obj != nil {
			dAtA.Gifts = append(dAtA.Gifts, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[gifts] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("gifts"), *pArSeR)
		}
	}

	return nil
}

type related_configs interface {
	GetCost(int) *resdata.Cost
	GetPrize(int) *resdata.Prize
	GetTimeLimitGiftData(uint64) *TimeLimitGiftData
	GetTimeRuleData(uint64) *data.TimeRuleData
}
