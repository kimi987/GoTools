// AUTO_GEN, DONT MODIFY!!!
package dianquan

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
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

// start with ExchangeMiscData ----------------------------------

func LoadExchangeMiscData(gos *config.GameObjects) (*ExchangeMiscData, *config.ObjectParser, error) {
	fIlEnAmE := confpath.ExchangeMiscDataPath
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

	dAtA, err := NewExchangeMiscData(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedExchangeMiscData(gos *config.GameObjects, dAtA *ExchangeMiscData, cOnFigS interface{}) error {
	fIlEnAmE := confpath.ExchangeMiscDataPath
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

func NewExchangeMiscData(fIlEnAmE string, pArSeR *config.ObjectParser) (*ExchangeMiscData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrExchangeMiscData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &ExchangeMiscData{}

	dAtA.ExchangeBaseYuanbao = pArSeR.Uint64("exchange_base_yuanbao")
	dAtA.ExchangeBaseDianquan = pArSeR.Uint64("exchange_base_dianquan")

	return dAtA, nil
}

var vAlIdAtOrExchangeMiscData = map[string]*config.Validator{

	"exchange_base_yuanbao":  config.ParseValidator("int>0", "", false, nil, nil),
	"exchange_base_dianquan": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *ExchangeMiscData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *ExchangeMiscData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *ExchangeMiscData) Encode() *shared_proto.DianquanMiscProto {
	out := &shared_proto.DianquanMiscProto{}
	out.ExchangeBaseYuanbao = config.U64ToI32(dAtA.ExchangeBaseYuanbao)
	out.ExchangeBaseDianquan = config.U64ToI32(dAtA.ExchangeBaseDianquan)

	return out
}

func ArrayEncodeExchangeMiscData(datas []*ExchangeMiscData) []*shared_proto.DianquanMiscProto {

	out := make([]*shared_proto.DianquanMiscProto, 0, len(datas))
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

func (dAtA *ExchangeMiscData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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
}
