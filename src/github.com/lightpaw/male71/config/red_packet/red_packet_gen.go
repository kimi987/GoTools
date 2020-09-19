// AUTO_GEN, DONT MODIFY!!!
package red_packet

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/icon"
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

// start with RedPacketData ----------------------------------

func LoadRedPacketData(gos *config.GameObjects) (map[uint64]*RedPacketData, map[*RedPacketData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.RedPacketDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*RedPacketData, len(lIsT))
	pArSeRmAp := make(map[*RedPacketData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrRedPacketData) {
			continue
		}

		dAtA, err := NewRedPacketData(fIlEnAmE, pArSeR)
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

func SetRelatedRedPacketData(dAtAmAp map[*RedPacketData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.RedPacketDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetRedPacketDataKeyArray(datas []*RedPacketData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewRedPacketData(fIlEnAmE string, pArSeR *config.ObjectParser) (*RedPacketData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrRedPacketData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &RedPacketData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Icon
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.DefaultText = pArSeR.String("default_text")
	// releated field: Cost
	dAtA.AmountType = shared_proto.AmountType(shared_proto.AmountType_value[strings.ToUpper(pArSeR.String("amount_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("amount_type"), 10, 32); err == nil {
		dAtA.AmountType = shared_proto.AmountType(i)
	}

	dAtA.Money = pArSeR.Uint64("money")
	// releated field: AllGarbbedPrize
	if pArSeR.KeyExist("expired_duration") {
		dAtA.ExpiredDuration, err = config.ParseDuration(pArSeR.String("expired_duration"))
	} else {
		dAtA.ExpiredDuration, err = config.ParseDuration("24h")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[expired_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("expired_duration"), dAtA)
	}

	dAtA.MinPartMoney = 1
	if pArSeR.KeyExist("min_part_money") {
		dAtA.MinPartMoney = pArSeR.Uint64("min_part_money")
	}

	dAtA.MinCount = 1
	if pArSeR.KeyExist("min_count") {
		dAtA.MinCount = pArSeR.Uint64("min_count")
	}

	dAtA.MaxCount = 100
	if pArSeR.KeyExist("max_count") {
		dAtA.MaxCount = pArSeR.Uint64("max_count")
	}

	dAtA.MaxTextLen = 10
	if pArSeR.KeyExist("max_text_len") {
		dAtA.MaxTextLen = pArSeR.Uint64("max_text_len")
	}

	return dAtA, nil
}

var vAlIdAtOrRedPacketData = map[string]*config.Validator{

	"id":                config.ParseValidator("int>0", "", false, nil, nil),
	"icon":              config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"name":              config.ParseValidator("string", "", false, nil, nil),
	"desc":              config.ParseValidator("string", "", false, nil, nil),
	"default_text":      config.ParseValidator("string", "", false, nil, nil),
	"cost":              config.ParseValidator("string", "", false, nil, nil),
	"amount_type":       config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.AmountType_value, 0), nil),
	"money":             config.ParseValidator("int>0", "", false, nil, nil),
	"all_garbbed_prize": config.ParseValidator("string", "", false, nil, nil),
	"expired_duration":  config.ParseValidator("string", "", false, nil, []string{"24h"}),
	"min_part_money":    config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"min_count":         config.ParseValidator("int>0", "", false, nil, []string{"1"}),
	"max_count":         config.ParseValidator("int>0", "", false, nil, []string{"100"}),
	"max_text_len":      config.ParseValidator("int>0", "", false, nil, []string{"10"}),
}

func (dAtA *RedPacketData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *RedPacketData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *RedPacketData) Encode() *shared_proto.RedPacketDataProto {
	out := &shared_proto.RedPacketDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.Icon != nil {
		out.Icon = dAtA.Icon.Encode()
	}
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.DefaultText = dAtA.DefaultText
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	out.AmountType = dAtA.AmountType
	out.Money = config.U64ToI32(dAtA.Money)
	if dAtA.AllGarbbedPrize != nil {
		out.AllGarbbedPrize = dAtA.AllGarbbedPrize.Encode()
	}
	out.ExpiredDuration = config.Duration2I32Seconds(dAtA.ExpiredDuration)
	out.MinCount = config.U64ToI32(dAtA.MinCount)
	out.MaxCount = config.U64ToI32(dAtA.MaxCount)
	out.MaxTextLen = config.U64ToI32(dAtA.MaxTextLen)

	return out
}

func ArrayEncodeRedPacketData(datas []*RedPacketData) []*shared_proto.RedPacketDataProto {

	out := make([]*shared_proto.RedPacketDataProto, 0, len(datas))
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

func (dAtA *RedPacketData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("icon") {
		dAtA.Icon = cOnFigS.GetIcon(pArSeR.String("icon"))
	} else {
		dAtA.Icon = cOnFigS.GetIcon("Icon")
	}
	if dAtA.Icon == nil {
		return errors.Errorf("%s 配置的关联字段[icon] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("icon"), *pArSeR)
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	dAtA.AllGarbbedPrize = cOnFigS.GetPrize(pArSeR.Int("all_garbbed_prize"))
	if dAtA.AllGarbbedPrize == nil {
		return errors.Errorf("%s 配置的关联字段[all_garbbed_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("all_garbbed_prize"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetCost(int) *resdata.Cost
	GetIcon(string) *icon.Icon
	GetPrize(int) *resdata.Prize
}
