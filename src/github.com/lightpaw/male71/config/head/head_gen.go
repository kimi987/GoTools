// AUTO_GEN, DONT MODIFY!!!
package head

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/icon"
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

// start with HeadData ----------------------------------

func LoadHeadData(gos *config.GameObjects) (map[string]*HeadData, map[*HeadData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.HeadDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*HeadData, len(lIsT))
	pArSeRmAp := make(map[*HeadData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrHeadData) {
			continue
		}

		dAtA, err := NewHeadData(fIlEnAmE, pArSeR)
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

func SetRelatedHeadData(dAtAmAp map[*HeadData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.HeadDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetHeadDataKeyArray(datas []*HeadData) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewHeadData(fIlEnAmE string, pArSeR *config.ObjectParser) (*HeadData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrHeadData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &HeadData{}

	dAtA.Id = pArSeR.String("id")
	// releated field: Icon
	// releated field: UnlockNeedCaptain
	// releated field: UnlockNeedHeroLevel
	dAtA.DefaultHead = true
	if pArSeR.KeyExist("default_head") {
		dAtA.DefaultHead = pArSeR.Bool("default_head")
	}

	dAtA.CountryOfficial = 0
	if pArSeR.KeyExist("country_official") {
		dAtA.CountryOfficial = pArSeR.Uint64("country_official")
	}

	// skip field: CountryOfficialType

	return dAtA, nil
}

var vAlIdAtOrHeadData = map[string]*config.Validator{

	"id":                     config.ParseValidator("string", "", false, nil, nil),
	"icon":                   config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"unlock_need_captain":    config.ParseValidator("string", "", false, nil, nil),
	"unlock_need_hero_level": config.ParseValidator("string", "", false, nil, nil),
	"default_head":           config.ParseValidator("bool", "", false, nil, []string{"true"}),
	"country_official":       config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *HeadData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *HeadData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *HeadData) Encode() *shared_proto.HeadProto {
	out := &shared_proto.HeadProto{}
	out.Id = dAtA.Id
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	if dAtA.UnlockNeedCaptain != nil {
		out.UnlockNeedCaptainSoul = config.U64ToI32(dAtA.UnlockNeedCaptain.Id)
	}
	if dAtA.UnlockNeedHeroLevel != nil {
		out.UnlockNeedHeroLevel = config.U64ToI32(dAtA.UnlockNeedHeroLevel.Level)
	}

	return out
}

func ArrayEncodeHeadData(datas []*HeadData) []*shared_proto.HeadProto {

	out := make([]*shared_proto.HeadProto, 0, len(datas))
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

func (dAtA *HeadData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.UnlockNeedCaptain = cOnFigS.GetCaptainData(pArSeR.Uint64("unlock_need_captain"))
	if dAtA.UnlockNeedCaptain == nil && pArSeR.Uint64("unlock_need_captain") != 0 {
		return errors.Errorf("%s 配置的关联字段[unlock_need_captain] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("unlock_need_captain"), *pArSeR)
	}

	dAtA.UnlockNeedHeroLevel = cOnFigS.GetHeroLevelData(pArSeR.Uint64("unlock_need_hero_level"))
	if dAtA.UnlockNeedHeroLevel == nil && pArSeR.Uint64("unlock_need_hero_level") != 0 {
		return errors.Errorf("%s 配置的关联字段[unlock_need_hero_level] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("unlock_need_hero_level"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetCaptainData(uint64) *captain.CaptainData
	GetHeroLevelData(uint64) *herodata.HeroLevelData
	GetIcon(string) *icon.Icon
}
