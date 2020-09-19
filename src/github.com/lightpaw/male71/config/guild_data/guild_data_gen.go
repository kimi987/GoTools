// AUTO_GEN, DONT MODIFY!!!
package guild_data

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/monsterdata"
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

// start with GuildBigBoxData ----------------------------------

func LoadGuildBigBoxData(gos *config.GameObjects) (map[uint64]*GuildBigBoxData, map[*GuildBigBoxData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildBigBoxDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildBigBoxData, len(lIsT))
	pArSeRmAp := make(map[*GuildBigBoxData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildBigBoxData) {
			continue
		}

		dAtA, err := NewGuildBigBoxData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildBigBoxData(dAtAmAp map[*GuildBigBoxData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildBigBoxDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildBigBoxDataKeyArray(datas []*GuildBigBoxData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildBigBoxData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildBigBoxData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildBigBoxData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildBigBoxData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: PlunderPrize
	dAtA.UnlockEnergy = pArSeR.Uint64("unlock_energy")
	dAtA.GuildLevelPrizeGroupId = 0
	if pArSeR.KeyExist("guild_level_prize_group_id") {
		dAtA.GuildLevelPrizeGroupId = pArSeR.Uint64("guild_level_prize_group_id")
	}

	// skip field: GuildLevelPrizes
	// skip field: TechLevel

	return dAtA, nil
}

var vAlIdAtOrGuildBigBoxData = map[string]*config.Validator{

	"id":                         config.ParseValidator("int>0", "", false, nil, nil),
	"plunder_prize":              config.ParseValidator("string", "", false, nil, nil),
	"unlock_energy":              config.ParseValidator("int>0", "", false, nil, nil),
	"guild_level_prize_group_id": config.ParseValidator("uint", "", false, nil, []string{"0"}),
}

func (dAtA *GuildBigBoxData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildBigBoxData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildBigBoxData) Encode() *shared_proto.GuildBigBoxDataProto {
	out := &shared_proto.GuildBigBoxDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	if dAtA.PlunderPrize != nil {
		out.Prize = dAtA.PlunderPrize.Prize.PrizeProto()
	}
	out.UnlockEnergy = config.U64ToI32(dAtA.UnlockEnergy)
	if dAtA.GuildLevelPrizes != nil {
		out.GuildLevelPrizes = resdata.ArrayEncodeGuildLevelPrize(dAtA.GuildLevelPrizes)
	}
	out.TechLevel = config.U64ToI32(dAtA.TechLevel)

	return out
}

func ArrayEncodeGuildBigBoxData(datas []*GuildBigBoxData) []*shared_proto.GuildBigBoxDataProto {

	out := make([]*shared_proto.GuildBigBoxDataProto, 0, len(datas))
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

func (dAtA *GuildBigBoxData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	return nil
}

// start with GuildClassLevelData ----------------------------------

func LoadGuildClassLevelData(gos *config.GameObjects) (map[uint64]*GuildClassLevelData, map[*GuildClassLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildClassLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildClassLevelData, len(lIsT))
	pArSeRmAp := make(map[*GuildClassLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildClassLevelData) {
			continue
		}

		dAtA, err := NewGuildClassLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildClassLevelData(dAtAmAp map[*GuildClassLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildClassLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildClassLevelDataKeyArray(datas []*GuildClassLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewGuildClassLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildClassLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildClassLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildClassLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Name = pArSeR.String("name")
	dAtA.CorePrestige = pArSeR.Bool("core_prestige")
	dAtA.VoteScore = pArSeR.Uint64("vote_score")
	dAtA.Permission, err = NewGuildPermissionData(fIlEnAmE, pArSeR)
	if err != nil {
		return nil, err
	}

	return dAtA, nil
}

var vAlIdAtOrGuildClassLevelData = map[string]*config.Validator{

	"level":         config.ParseValidator("int>0", "", false, nil, nil),
	"name":          config.ParseValidator("string", "", false, nil, nil),
	"core_prestige": config.ParseValidator("bool", "", false, nil, nil),
	"vote_score":    config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *GuildClassLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildClassLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildClassLevelData) Encode() *shared_proto.GuildClassLevelProto {
	out := &shared_proto.GuildClassLevelProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.Name = dAtA.Name
	out.VoteScore = config.U64ToI32(dAtA.VoteScore)
	if dAtA.Permission != nil {
		out.Permission = dAtA.Permission.Encode()
	}

	return out
}

func ArrayEncodeGuildClassLevelData(datas []*GuildClassLevelData) []*shared_proto.GuildClassLevelProto {

	out := make([]*shared_proto.GuildClassLevelProto, 0, len(datas))
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

func (dAtA *GuildClassLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if err := dAtA.Permission.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS0); err != nil {
		return err
	}

	return nil
}

// start with GuildClassTitleData ----------------------------------

func LoadGuildClassTitleData(gos *config.GameObjects) (map[uint64]*GuildClassTitleData, map[*GuildClassTitleData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildClassTitleDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildClassTitleData, len(lIsT))
	pArSeRmAp := make(map[*GuildClassTitleData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildClassTitleData) {
			continue
		}

		dAtA, err := NewGuildClassTitleData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildClassTitleData(dAtAmAp map[*GuildClassTitleData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildClassTitleDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildClassTitleDataKeyArray(datas []*GuildClassTitleData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildClassTitleData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildClassTitleData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildClassTitleData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildClassTitleData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Permission, err = NewGuildPermissionData(fIlEnAmE, pArSeR)
	if err != nil {
		return nil, err
	}

	return dAtA, nil
}

var vAlIdAtOrGuildClassTitleData = map[string]*config.Validator{

	"id":   config.ParseValidator("int>0", "", false, nil, nil),
	"name": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *GuildClassTitleData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildClassTitleData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildClassTitleData) Encode() *shared_proto.GuildClassTitleDataProto {
	out := &shared_proto.GuildClassTitleDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	if dAtA.Permission != nil {
		out.Permission = dAtA.Permission.Encode()
	}

	return out
}

func ArrayEncodeGuildClassTitleData(datas []*GuildClassTitleData) []*shared_proto.GuildClassTitleDataProto {

	out := make([]*shared_proto.GuildClassTitleDataProto, 0, len(datas))
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

func (dAtA *GuildClassTitleData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if err := dAtA.Permission.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS0); err != nil {
		return err
	}

	return nil
}

// start with GuildDonateData ----------------------------------

func LoadGuildDonateData(gos *config.GameObjects) (map[uint64]*GuildDonateData, map[*GuildDonateData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildDonateDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildDonateData, len(lIsT))
	pArSeRmAp := make(map[*GuildDonateData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildDonateData) {
			continue
		}

		dAtA, err := NewGuildDonateData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildDonateData(dAtAmAp map[*GuildDonateData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildDonateDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildDonateDataKeyArray(datas []*GuildDonateData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildDonateData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildDonateData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildDonateData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildDonateData{}

	dAtA.Sequence = pArSeR.Uint64("sequence")
	dAtA.Times = pArSeR.Uint64("times")
	// releated field: Cost
	dAtA.GuildBuildingAmount = pArSeR.Uint64("guild_building_amount")
	dAtA.ContributionAmount = pArSeR.Uint64("contribution_amount")
	dAtA.DonationAmount = pArSeR.Uint64("donation_amount")
	dAtA.ContributionCoin = pArSeR.Uint64("contribution_coin")
	dAtA.RecommandGuanfuLevel = 1
	if pArSeR.KeyExist("recommand_guanfu_level") {
		dAtA.RecommandGuanfuLevel = pArSeR.Uint64("recommand_guanfu_level")
	}

	// calculate fields
	dAtA.Id = DonateId(dAtA.Sequence, dAtA.Times)

	return dAtA, nil
}

var vAlIdAtOrGuildDonateData = map[string]*config.Validator{

	"sequence":               config.ParseValidator("int>0", "", false, nil, nil),
	"times":                  config.ParseValidator("int>0", "", false, nil, nil),
	"cost":                   config.ParseValidator("string", "", false, nil, nil),
	"guild_building_amount":  config.ParseValidator("int>0", "", false, nil, nil),
	"contribution_amount":    config.ParseValidator("int>0", "", false, nil, nil),
	"donation_amount":        config.ParseValidator("int>0", "", false, nil, nil),
	"contribution_coin":      config.ParseValidator("int>0", "", false, nil, nil),
	"recommand_guanfu_level": config.ParseValidator("int>0", "", false, nil, []string{"1"}),
}

func (dAtA *GuildDonateData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildDonateData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildDonateData) Encode() *shared_proto.GuildDonateProto {
	out := &shared_proto.GuildDonateProto{}
	out.Sequence = config.U64ToI32(dAtA.Sequence)
	out.Times = config.U64ToI32(dAtA.Times)
	if dAtA.Cost != nil {
		out.Cost = dAtA.Cost.Encode()
	}
	out.GuildBuildingAmount = config.U64ToI32(dAtA.GuildBuildingAmount)
	out.ContributionAmount = config.U64ToI32(dAtA.ContributionAmount)
	out.DonationAmount = config.U64ToI32(dAtA.DonationAmount)
	out.ContributionCoin = config.U64ToI32(dAtA.ContributionCoin)
	out.RecommandGuanfuLevel = config.U64ToI32(dAtA.RecommandGuanfuLevel)

	return out
}

func ArrayEncodeGuildDonateData(datas []*GuildDonateData) []*shared_proto.GuildDonateProto {

	out := make([]*shared_proto.GuildDonateProto, 0, len(datas))
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

func (dAtA *GuildDonateData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Cost = cOnFigS.GetCost(pArSeR.Int("cost"))
	if dAtA.Cost == nil {
		return errors.Errorf("%s 配置的关联字段[cost] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("cost"), *pArSeR)
	}

	return nil
}

// start with GuildEventPrizeData ----------------------------------

func LoadGuildEventPrizeData(gos *config.GameObjects) (map[uint64]*GuildEventPrizeData, map[*GuildEventPrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildEventPrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildEventPrizeData, len(lIsT))
	pArSeRmAp := make(map[*GuildEventPrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildEventPrizeData) {
			continue
		}

		dAtA, err := NewGuildEventPrizeData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildEventPrizeData(dAtAmAp map[*GuildEventPrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildEventPrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildEventPrizeDataKeyArray(datas []*GuildEventPrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildEventPrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildEventPrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildEventPrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildEventPrizeData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Quality = shared_proto.Quality(shared_proto.Quality_value[strings.ToUpper(pArSeR.String("quality"))])
	if i, err := strconv.ParseInt(pArSeR.String("quality"), 10, 32); err == nil {
		dAtA.Quality = shared_proto.Quality(i)
	}

	// releated field: Icon
	// releated field: Prize
	dAtA.ExipreDuration, err = config.ParseDuration(pArSeR.String("exipre_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[exipre_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("exipre_duration"), dAtA)
	}

	dAtA.FromShop = false
	if pArSeR.KeyExist("from_shop") {
		dAtA.FromShop = pArSeR.Bool("from_shop")
	}

	dAtA.Energy = pArSeR.Uint64("energy")
	dAtA.DailyLimit = pArSeR.Uint64("daily_limit")
	dAtA.GuildLevelPrizeGroupId = 0
	if pArSeR.KeyExist("guild_level_prize_group_id") {
		dAtA.GuildLevelPrizeGroupId = pArSeR.Uint64("guild_level_prize_group_id")
	}

	// skip field: GuildLevelPrizes
	dAtA.TriggerEvent = shared_proto.HeroEvent(shared_proto.HeroEvent_value[strings.ToUpper(pArSeR.String("trigger_event"))])
	if i, err := strconv.ParseInt(pArSeR.String("trigger_event"), 10, 32); err == nil {
		dAtA.TriggerEvent = shared_proto.HeroEvent(i)
	}

	dAtA.TriggerEventCondition, err = data.ParseCompareCondition(pArSeR.String("trigger_event_condition"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[trigger_event_condition] 解析失败(data.ParseCompareCondition)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("trigger_event_condition"), dAtA)
	}

	dAtA.TriggerEventTimes = pArSeR.Uint64("trigger_event_times")
	dAtA.TriggerEventDailyReset = pArSeR.Bool("trigger_event_daily_reset")

	return dAtA, nil
}

var vAlIdAtOrGuildEventPrizeData = map[string]*config.Validator{

	"id":                         config.ParseValidator("int>0", "", false, nil, nil),
	"name":                       config.ParseValidator("string", "", false, nil, nil),
	"desc":                       config.ParseValidator("string", "", false, nil, nil),
	"quality":                    config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.Quality_value, 0), nil),
	"icon":                       config.ParseValidator("string", "", false, nil, []string{"Icon"}),
	"prize":                      config.ParseValidator("string", "", false, nil, nil),
	"exipre_duration":            config.ParseValidator("string", "", false, nil, nil),
	"from_shop":                  config.ParseValidator("bool", "", false, nil, []string{"false"}),
	"energy":                     config.ParseValidator("int>0", "", false, nil, nil),
	"daily_limit":                config.ParseValidator("uint", "", false, nil, nil),
	"guild_level_prize_group_id": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"trigger_event":              config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.HeroEvent_value), nil),
	"trigger_event_condition":    config.ParseValidator("string", "", false, nil, nil),
	"trigger_event_times":        config.ParseValidator("uint", "", false, nil, nil),
	"trigger_event_daily_reset":  config.ParseValidator("bool", "", false, nil, nil),
}

func (dAtA *GuildEventPrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildEventPrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildEventPrizeData) Encode() *shared_proto.GuildEventPrizeDataProto {
	out := &shared_proto.GuildEventPrizeDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.Quality = dAtA.Quality
	if dAtA.Icon != nil {
		out.IconId = dAtA.Icon.Id
	}
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Prize.PrizeProto()
	}
	out.FromShop = dAtA.FromShop
	out.Energy = config.U64ToI32(dAtA.Energy)
	if dAtA.GuildLevelPrizes != nil {
		out.GuildLevelPrizes = resdata.ArrayEncodeGuildLevelPrize(dAtA.GuildLevelPrizes)
	}

	return out
}

func ArrayEncodeGuildEventPrizeData(datas []*GuildEventPrizeData) []*shared_proto.GuildEventPrizeDataProto {

	out := make([]*shared_proto.GuildEventPrizeDataProto, 0, len(datas))
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

func (dAtA *GuildEventPrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.Prize = cOnFigS.GetPlunderPrize(pArSeR.Uint64("prize"))
	if dAtA.Prize == nil {
		return errors.Errorf("%s 配置的关联字段[prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prize"), *pArSeR)
	}

	return nil
}

// start with GuildLevelCdrData ----------------------------------

func LoadGuildLevelCdrData(gos *config.GameObjects) (map[uint64]*GuildLevelCdrData, map[*GuildLevelCdrData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildLevelCdrDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildLevelCdrData, len(lIsT))
	pArSeRmAp := make(map[*GuildLevelCdrData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildLevelCdrData) {
			continue
		}

		dAtA, err := NewGuildLevelCdrData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildLevelCdrData(dAtAmAp map[*GuildLevelCdrData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildLevelCdrDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildLevelCdrDataKeyArray(datas []*GuildLevelCdrData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildLevelCdrData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildLevelCdrData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildLevelCdrData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildLevelCdrData{}

	dAtA.Group = 0
	if pArSeR.KeyExist("group") {
		dAtA.Group = pArSeR.Uint64("group")
	}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.Times = pArSeR.Uint64("times")
	dAtA.Cost = pArSeR.Uint64("cost")
	dAtA.CDR, err = config.ParseDuration(pArSeR.String("cdr"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[cdr] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("cdr"), dAtA)
	}

	// calculate fields
	dAtA.Id = GuildLevelCdrId(dAtA.Group, dAtA.Level, dAtA.Times)

	return dAtA, nil
}

var vAlIdAtOrGuildLevelCdrData = map[string]*config.Validator{

	"group": config.ParseValidator("uint", "", false, nil, []string{"0"}),
	"level": config.ParseValidator("int>0", "", false, nil, nil),
	"times": config.ParseValidator("int>0", "", false, nil, nil),
	"cost":  config.ParseValidator("int>0", "", false, nil, nil),
	"cdr":   config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *GuildLevelCdrData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildLevelCdrData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildLevelCdrData) Encode() *shared_proto.GuildLevelCdrProto {
	out := &shared_proto.GuildLevelCdrProto{}
	out.Times = config.U64ToI32(dAtA.Times)
	out.Cost = config.U64ToI32(dAtA.Cost)
	out.Cdr = config.Duration2I32Seconds(dAtA.CDR)

	return out
}

func ArrayEncodeGuildLevelCdrData(datas []*GuildLevelCdrData) []*shared_proto.GuildLevelCdrProto {

	out := make([]*shared_proto.GuildLevelCdrProto, 0, len(datas))
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

func (dAtA *GuildLevelCdrData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GuildLevelData ----------------------------------

func LoadGuildLevelData(gos *config.GameObjects) (map[uint64]*GuildLevelData, map[*GuildLevelData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildLevelDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildLevelData, len(lIsT))
	pArSeRmAp := make(map[*GuildLevelData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildLevelData) {
			continue
		}

		dAtA, err := NewGuildLevelData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildLevelData(dAtAmAp map[*GuildLevelData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildLevelDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildLevelDataKeyArray(datas []*GuildLevelData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Level)
		}
	}

	return out
}

func NewGuildLevelData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildLevelData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildLevelData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildLevelData{}

	dAtA.Level = pArSeR.Uint64("level")
	dAtA.MemberCount = pArSeR.Uint64("member_count")
	dAtA.ClassMemberCount = pArSeR.Uint64Array("class_member_count", "", false)
	dAtA.UpgradeBuilding = pArSeR.Uint64("upgrade_building")
	dAtA.UpgradeDuration, err = config.ParseDuration(pArSeR.String("upgrade_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[upgrade_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("upgrade_duration"), dAtA)
	}

	// skip field: Cdrs

	return dAtA, nil
}

var vAlIdAtOrGuildLevelData = map[string]*config.Validator{

	"level":              config.ParseValidator("int>0", "", false, nil, nil),
	"member_count":       config.ParseValidator("int>0", "", false, nil, nil),
	"class_member_count": config.ParseValidator("uint,duplicate", "", true, nil, nil),
	"upgrade_building":   config.ParseValidator("int>0", "", false, nil, nil),
	"upgrade_duration":   config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *GuildLevelData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildLevelData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildLevelData) Encode() *shared_proto.GuildLevelProto {
	out := &shared_proto.GuildLevelProto{}
	out.Level = config.U64ToI32(dAtA.Level)
	out.MemberCount = config.U64ToI32(dAtA.MemberCount)
	out.ClassMemberCount = config.U64a2I32a(dAtA.ClassMemberCount)
	out.UpgradeBuilding = config.U64ToI32(dAtA.UpgradeBuilding)
	out.UpgradeDuration = config.Duration2I32Seconds(dAtA.UpgradeDuration)
	if dAtA.Cdrs != nil {
		out.Cdrs = ArrayEncodeGuildLevelCdrData(dAtA.Cdrs)
	}

	return out
}

func ArrayEncodeGuildLevelData(datas []*GuildLevelData) []*shared_proto.GuildLevelProto {

	out := make([]*shared_proto.GuildLevelProto, 0, len(datas))
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

func (dAtA *GuildLevelData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GuildLogData ----------------------------------

func LoadGuildLogData(gos *config.GameObjects) (map[string]*GuildLogData, map[*GuildLogData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildLogDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[string]*GuildLogData, len(lIsT))
	pArSeRmAp := make(map[*GuildLogData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildLogData) {
			continue
		}

		dAtA, err := NewGuildLogData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildLogData(dAtAmAp map[*GuildLogData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildLogDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildLogDataKeyArray(datas []*GuildLogData) []string {

	out := make([]string, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildLogData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildLogData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildLogData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildLogData{}

	dAtA.Id = pArSeR.String("id")
	dAtA.LogType = shared_proto.GuildLogType(shared_proto.GuildLogType_value[strings.ToUpper(pArSeR.String("log_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("log_type"), 10, 32); err == nil {
		dAtA.LogType = shared_proto.GuildLogType(i)
	}

	dAtA.Icon = pArSeR.String("icon")
	dAtA.SendChat = false
	if pArSeR.KeyExist("send_chat") {
		dAtA.SendChat = pArSeR.Bool("send_chat")
	}

	// i18n fields
	dAtA.Text = i18n.NewI18nRef(fIlEnAmE, "text", dAtA.Id, pArSeR.String("text"))

	return dAtA, nil
}

var vAlIdAtOrGuildLogData = map[string]*config.Validator{

	"id":        config.ParseValidator("string", "", false, nil, nil),
	"log_type":  config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.GuildLogType_value, 0), nil),
	"icon":      config.ParseValidator("string", "", false, nil, nil),
	"text":      config.ParseValidator("string", "", false, nil, nil),
	"send_chat": config.ParseValidator("bool", "", false, nil, []string{"false"}),
}

func (dAtA *GuildLogData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GuildLogHelp ----------------------------------

func LoadGuildLogHelp(gos *config.GameObjects) (*GuildLogHelp, *config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildLogHelpPath
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

	dAtA, err := NewGuildLogHelp(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedGuildLogHelp(gos *config.GameObjects, dAtA *GuildLogHelp, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildLogHelpPath
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

func NewGuildLogHelp(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildLogHelp, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildLogHelp)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildLogHelp{}

	// releated field: JoinGuild
	// releated field: LeaveGuild
	// releated field: ReplyJoinGuild
	// releated field: KickLeaveGuild
	// releated field: UpdateMemberClass
	// releated field: UpgradeTechnology
	// releated field: PrestigePrize
	// releated field: InvaseMonster
	// releated field: CollectSalary
	// releated field: CreateGuild
	// releated field: NewLeaderImpeach
	// releated field: NewLeaderDemise
	// releated field: StartImpeach
	// releated field: TerminateImpeach
	// releated field: UpdateName
	// releated field: UpdateFlagName
	// releated field: UpdateCountry
	// releated field: UpgradeLevel
	// releated field: UpdateInternalText
	// releated field: FAttSucc
	// releated field: FDefFail
	// releated field: FAttDestroy
	// releated field: FDefDestroy
	// releated field: StartXiongNu
	// releated field: WipeOutXiongNuTroop
	// releated field: UnlockXiongNu
	// releated field: ResistXiongNuAddPrestige
	// releated field: XiongNuBaseDestroy
	// releated field: YinliangMingcHost
	// releated field: YinliangMcWarAtkWin
	// releated field: YinliangMcWarAtkFail
	// releated field: YinliangGuildReceive
	// releated field: YinliangGuildSend
	// releated field: YinliangGuildSendMember
	// releated field: YinliangGuildPaySalary
	// releated field: McWarApplyAtkFail
	// releated field: McWarApplyAtkSucc
	// releated field: McWarAtkWin
	// releated field: McWarDefWin
	// releated field: HufuAdd
	// releated field: BaowuAddPrestige
	// releated field: UpdateMark
	// releated field: AssemblyWin

	return dAtA, nil
}

var vAlIdAtOrGuildLogHelp = map[string]*config.Validator{

	"join_guild":                   config.ParseValidator("string", "", false, nil, []string{"JoinGuild"}),
	"leave_guild":                  config.ParseValidator("string", "", false, nil, []string{"LeaveGuild"}),
	"reply_join_guild":             config.ParseValidator("string", "", false, nil, []string{"ReplyJoinGuild"}),
	"kick_leave_guild":             config.ParseValidator("string", "", false, nil, []string{"KickLeaveGuild"}),
	"update_member_class":          config.ParseValidator("string", "", false, nil, []string{"UpdateMemberClass"}),
	"upgrade_technology":           config.ParseValidator("string", "", false, nil, []string{"UpgradeTechnology"}),
	"prestige_prize":               config.ParseValidator("string", "", false, nil, []string{"PrestigePrize"}),
	"invase_monster":               config.ParseValidator("string", "", false, nil, []string{"InvaseMonster"}),
	"collect_salary":               config.ParseValidator("string", "", false, nil, []string{"CollectSalary"}),
	"create_guild":                 config.ParseValidator("string", "", false, nil, []string{"CreateGuild"}),
	"new_leader_impeach":           config.ParseValidator("string", "", false, nil, []string{"NewLeaderImpeach"}),
	"new_leader_demise":            config.ParseValidator("string", "", false, nil, []string{"NewLeaderDemise"}),
	"start_impeach":                config.ParseValidator("string", "", false, nil, []string{"StartImpeach"}),
	"terminate_impeach":            config.ParseValidator("string", "", false, nil, []string{"TerminateImpeach"}),
	"update_name":                  config.ParseValidator("string", "", false, nil, []string{"UpdateName"}),
	"update_flag_name":             config.ParseValidator("string", "", false, nil, []string{"UpdateFlagName"}),
	"update_country":               config.ParseValidator("string", "", false, nil, []string{"UpdateCountry"}),
	"upgrade_level":                config.ParseValidator("string", "", false, nil, []string{"UpgradeLevel"}),
	"update_internal_text":         config.ParseValidator("string", "", false, nil, []string{"UpdateInternalText"}),
	"fatt_succ":                    config.ParseValidator("string", "", false, nil, []string{"FAttSucc"}),
	"fdef_fail":                    config.ParseValidator("string", "", false, nil, []string{"FDefFail"}),
	"fatt_destroy":                 config.ParseValidator("string", "", false, nil, []string{"FAttDestroy"}),
	"fdef_destroy":                 config.ParseValidator("string", "", false, nil, []string{"FDefDestroy"}),
	"start_xiong_nu":               config.ParseValidator("string", "", false, nil, []string{"StartXiongNu"}),
	"wipe_out_xiong_nu_troop":      config.ParseValidator("string", "", false, nil, []string{"WipeOutXiongNuTroop"}),
	"unlock_xiong_nu":              config.ParseValidator("string", "", false, nil, []string{"UnlockXiongNu"}),
	"resist_xiong_nu_add_prestige": config.ParseValidator("string", "", false, nil, []string{"ResistXiongNuAddPrestige"}),
	"xiong_nu_base_destroy":        config.ParseValidator("string", "", false, nil, []string{"XiongNuBaseDestroy"}),
	"yinliang_mingc_host":          config.ParseValidator("string", "", false, nil, []string{"YinliangMingcHost"}),
	"yinliang_mc_war_atk_win":      config.ParseValidator("string", "", false, nil, []string{"YinliangMcWarAtkWin"}),
	"yinliang_mc_war_atk_fail":     config.ParseValidator("string", "", false, nil, []string{"YinliangMcWarAtkFail"}),
	"yinliang_guild_receive":       config.ParseValidator("string", "", false, nil, []string{"YinliangGuildReceive"}),
	"yinliang_guild_send":          config.ParseValidator("string", "", false, nil, []string{"YinliangGuildSend"}),
	"yinliang_guild_send_member":   config.ParseValidator("string", "", false, nil, []string{"YinliangGuildSendMember"}),
	"yinliang_guild_pay_salary":    config.ParseValidator("string", "", false, nil, []string{"YinliangGuildPaySalary"}),
	"mc_war_apply_atk_fail":        config.ParseValidator("string", "", false, nil, []string{"McWarApplyAtkFail"}),
	"mc_war_apply_atk_succ":        config.ParseValidator("string", "", false, nil, []string{"McWarApplyAtkSucc"}),
	"mc_war_atk_win":               config.ParseValidator("string", "", false, nil, []string{"McWarAtkWin"}),
	"mc_war_def_win":               config.ParseValidator("string", "", false, nil, []string{"McWarDefWin"}),
	"hufu_add":                     config.ParseValidator("string", "", false, nil, []string{"HufuAdd"}),
	"baowu_add_prestige":           config.ParseValidator("string", "", false, nil, []string{"BaowuAddPrestige"}),
	"update_mark":                  config.ParseValidator("string", "", false, nil, []string{"UpdateMark"}),
	"assembly_win":                 config.ParseValidator("string", "", false, nil, []string{"AssemblyWin"}),
}

func (dAtA *GuildLogHelp) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	if pArSeR.KeyExist("join_guild") {
		dAtA.JoinGuild = cOnFigS.GetGuildLogData(pArSeR.String("join_guild"))
	} else {
		dAtA.JoinGuild = cOnFigS.GetGuildLogData("JoinGuild")
	}
	if dAtA.JoinGuild == nil {
		return errors.Errorf("%s 配置的关联字段[join_guild] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("join_guild"), *pArSeR)
	}

	if pArSeR.KeyExist("leave_guild") {
		dAtA.LeaveGuild = cOnFigS.GetGuildLogData(pArSeR.String("leave_guild"))
	} else {
		dAtA.LeaveGuild = cOnFigS.GetGuildLogData("LeaveGuild")
	}
	if dAtA.LeaveGuild == nil {
		return errors.Errorf("%s 配置的关联字段[leave_guild] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("leave_guild"), *pArSeR)
	}

	if pArSeR.KeyExist("reply_join_guild") {
		dAtA.ReplyJoinGuild = cOnFigS.GetGuildLogData(pArSeR.String("reply_join_guild"))
	} else {
		dAtA.ReplyJoinGuild = cOnFigS.GetGuildLogData("ReplyJoinGuild")
	}
	if dAtA.ReplyJoinGuild == nil {
		return errors.Errorf("%s 配置的关联字段[reply_join_guild] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("reply_join_guild"), *pArSeR)
	}

	if pArSeR.KeyExist("kick_leave_guild") {
		dAtA.KickLeaveGuild = cOnFigS.GetGuildLogData(pArSeR.String("kick_leave_guild"))
	} else {
		dAtA.KickLeaveGuild = cOnFigS.GetGuildLogData("KickLeaveGuild")
	}
	if dAtA.KickLeaveGuild == nil {
		return errors.Errorf("%s 配置的关联字段[kick_leave_guild] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("kick_leave_guild"), *pArSeR)
	}

	if pArSeR.KeyExist("update_member_class") {
		dAtA.UpdateMemberClass = cOnFigS.GetGuildLogData(pArSeR.String("update_member_class"))
	} else {
		dAtA.UpdateMemberClass = cOnFigS.GetGuildLogData("UpdateMemberClass")
	}
	if dAtA.UpdateMemberClass == nil {
		return errors.Errorf("%s 配置的关联字段[update_member_class] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("update_member_class"), *pArSeR)
	}

	if pArSeR.KeyExist("upgrade_technology") {
		dAtA.UpgradeTechnology = cOnFigS.GetGuildLogData(pArSeR.String("upgrade_technology"))
	} else {
		dAtA.UpgradeTechnology = cOnFigS.GetGuildLogData("UpgradeTechnology")
	}
	if dAtA.UpgradeTechnology == nil {
		return errors.Errorf("%s 配置的关联字段[upgrade_technology] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("upgrade_technology"), *pArSeR)
	}

	if pArSeR.KeyExist("prestige_prize") {
		dAtA.PrestigePrize = cOnFigS.GetGuildLogData(pArSeR.String("prestige_prize"))
	} else {
		dAtA.PrestigePrize = cOnFigS.GetGuildLogData("PrestigePrize")
	}
	if dAtA.PrestigePrize == nil {
		return errors.Errorf("%s 配置的关联字段[prestige_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prestige_prize"), *pArSeR)
	}

	if pArSeR.KeyExist("invase_monster") {
		dAtA.InvaseMonster = cOnFigS.GetGuildLogData(pArSeR.String("invase_monster"))
	} else {
		dAtA.InvaseMonster = cOnFigS.GetGuildLogData("InvaseMonster")
	}
	if dAtA.InvaseMonster == nil {
		return errors.Errorf("%s 配置的关联字段[invase_monster] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("invase_monster"), *pArSeR)
	}

	if pArSeR.KeyExist("collect_salary") {
		dAtA.CollectSalary = cOnFigS.GetGuildLogData(pArSeR.String("collect_salary"))
	} else {
		dAtA.CollectSalary = cOnFigS.GetGuildLogData("CollectSalary")
	}
	if dAtA.CollectSalary == nil {
		return errors.Errorf("%s 配置的关联字段[collect_salary] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("collect_salary"), *pArSeR)
	}

	if pArSeR.KeyExist("create_guild") {
		dAtA.CreateGuild = cOnFigS.GetGuildLogData(pArSeR.String("create_guild"))
	} else {
		dAtA.CreateGuild = cOnFigS.GetGuildLogData("CreateGuild")
	}
	if dAtA.CreateGuild == nil {
		return errors.Errorf("%s 配置的关联字段[create_guild] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("create_guild"), *pArSeR)
	}

	if pArSeR.KeyExist("new_leader_impeach") {
		dAtA.NewLeaderImpeach = cOnFigS.GetGuildLogData(pArSeR.String("new_leader_impeach"))
	} else {
		dAtA.NewLeaderImpeach = cOnFigS.GetGuildLogData("NewLeaderImpeach")
	}
	if dAtA.NewLeaderImpeach == nil {
		return errors.Errorf("%s 配置的关联字段[new_leader_impeach] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("new_leader_impeach"), *pArSeR)
	}

	if pArSeR.KeyExist("new_leader_demise") {
		dAtA.NewLeaderDemise = cOnFigS.GetGuildLogData(pArSeR.String("new_leader_demise"))
	} else {
		dAtA.NewLeaderDemise = cOnFigS.GetGuildLogData("NewLeaderDemise")
	}
	if dAtA.NewLeaderDemise == nil {
		return errors.Errorf("%s 配置的关联字段[new_leader_demise] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("new_leader_demise"), *pArSeR)
	}

	if pArSeR.KeyExist("start_impeach") {
		dAtA.StartImpeach = cOnFigS.GetGuildLogData(pArSeR.String("start_impeach"))
	} else {
		dAtA.StartImpeach = cOnFigS.GetGuildLogData("StartImpeach")
	}
	if dAtA.StartImpeach == nil {
		return errors.Errorf("%s 配置的关联字段[start_impeach] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("start_impeach"), *pArSeR)
	}

	if pArSeR.KeyExist("terminate_impeach") {
		dAtA.TerminateImpeach = cOnFigS.GetGuildLogData(pArSeR.String("terminate_impeach"))
	} else {
		dAtA.TerminateImpeach = cOnFigS.GetGuildLogData("TerminateImpeach")
	}
	if dAtA.TerminateImpeach == nil {
		return errors.Errorf("%s 配置的关联字段[terminate_impeach] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("terminate_impeach"), *pArSeR)
	}

	if pArSeR.KeyExist("update_name") {
		dAtA.UpdateName = cOnFigS.GetGuildLogData(pArSeR.String("update_name"))
	} else {
		dAtA.UpdateName = cOnFigS.GetGuildLogData("UpdateName")
	}
	if dAtA.UpdateName == nil {
		return errors.Errorf("%s 配置的关联字段[update_name] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("update_name"), *pArSeR)
	}

	if pArSeR.KeyExist("update_flag_name") {
		dAtA.UpdateFlagName = cOnFigS.GetGuildLogData(pArSeR.String("update_flag_name"))
	} else {
		dAtA.UpdateFlagName = cOnFigS.GetGuildLogData("UpdateFlagName")
	}
	if dAtA.UpdateFlagName == nil {
		return errors.Errorf("%s 配置的关联字段[update_flag_name] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("update_flag_name"), *pArSeR)
	}

	if pArSeR.KeyExist("update_country") {
		dAtA.UpdateCountry = cOnFigS.GetGuildLogData(pArSeR.String("update_country"))
	} else {
		dAtA.UpdateCountry = cOnFigS.GetGuildLogData("UpdateCountry")
	}
	if dAtA.UpdateCountry == nil {
		return errors.Errorf("%s 配置的关联字段[update_country] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("update_country"), *pArSeR)
	}

	if pArSeR.KeyExist("upgrade_level") {
		dAtA.UpgradeLevel = cOnFigS.GetGuildLogData(pArSeR.String("upgrade_level"))
	} else {
		dAtA.UpgradeLevel = cOnFigS.GetGuildLogData("UpgradeLevel")
	}
	if dAtA.UpgradeLevel == nil {
		return errors.Errorf("%s 配置的关联字段[upgrade_level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("upgrade_level"), *pArSeR)
	}

	if pArSeR.KeyExist("update_internal_text") {
		dAtA.UpdateInternalText = cOnFigS.GetGuildLogData(pArSeR.String("update_internal_text"))
	} else {
		dAtA.UpdateInternalText = cOnFigS.GetGuildLogData("UpdateInternalText")
	}
	if dAtA.UpdateInternalText == nil {
		return errors.Errorf("%s 配置的关联字段[update_internal_text] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("update_internal_text"), *pArSeR)
	}

	if pArSeR.KeyExist("fatt_succ") {
		dAtA.FAttSucc = cOnFigS.GetGuildLogData(pArSeR.String("fatt_succ"))
	} else {
		dAtA.FAttSucc = cOnFigS.GetGuildLogData("FAttSucc")
	}
	if dAtA.FAttSucc == nil {
		return errors.Errorf("%s 配置的关联字段[fatt_succ] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fatt_succ"), *pArSeR)
	}

	if pArSeR.KeyExist("fdef_fail") {
		dAtA.FDefFail = cOnFigS.GetGuildLogData(pArSeR.String("fdef_fail"))
	} else {
		dAtA.FDefFail = cOnFigS.GetGuildLogData("FDefFail")
	}
	if dAtA.FDefFail == nil {
		return errors.Errorf("%s 配置的关联字段[fdef_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fdef_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("fatt_destroy") {
		dAtA.FAttDestroy = cOnFigS.GetGuildLogData(pArSeR.String("fatt_destroy"))
	} else {
		dAtA.FAttDestroy = cOnFigS.GetGuildLogData("FAttDestroy")
	}
	if dAtA.FAttDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[fatt_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fatt_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("fdef_destroy") {
		dAtA.FDefDestroy = cOnFigS.GetGuildLogData(pArSeR.String("fdef_destroy"))
	} else {
		dAtA.FDefDestroy = cOnFigS.GetGuildLogData("FDefDestroy")
	}
	if dAtA.FDefDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[fdef_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("fdef_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("start_xiong_nu") {
		dAtA.StartXiongNu = cOnFigS.GetGuildLogData(pArSeR.String("start_xiong_nu"))
	} else {
		dAtA.StartXiongNu = cOnFigS.GetGuildLogData("StartXiongNu")
	}
	if dAtA.StartXiongNu == nil {
		return errors.Errorf("%s 配置的关联字段[start_xiong_nu] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("start_xiong_nu"), *pArSeR)
	}

	if pArSeR.KeyExist("wipe_out_xiong_nu_troop") {
		dAtA.WipeOutXiongNuTroop = cOnFigS.GetGuildLogData(pArSeR.String("wipe_out_xiong_nu_troop"))
	} else {
		dAtA.WipeOutXiongNuTroop = cOnFigS.GetGuildLogData("WipeOutXiongNuTroop")
	}
	if dAtA.WipeOutXiongNuTroop == nil {
		return errors.Errorf("%s 配置的关联字段[wipe_out_xiong_nu_troop] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("wipe_out_xiong_nu_troop"), *pArSeR)
	}

	if pArSeR.KeyExist("unlock_xiong_nu") {
		dAtA.UnlockXiongNu = cOnFigS.GetGuildLogData(pArSeR.String("unlock_xiong_nu"))
	} else {
		dAtA.UnlockXiongNu = cOnFigS.GetGuildLogData("UnlockXiongNu")
	}
	if dAtA.UnlockXiongNu == nil {
		return errors.Errorf("%s 配置的关联字段[unlock_xiong_nu] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("unlock_xiong_nu"), *pArSeR)
	}

	if pArSeR.KeyExist("resist_xiong_nu_add_prestige") {
		dAtA.ResistXiongNuAddPrestige = cOnFigS.GetGuildLogData(pArSeR.String("resist_xiong_nu_add_prestige"))
	} else {
		dAtA.ResistXiongNuAddPrestige = cOnFigS.GetGuildLogData("ResistXiongNuAddPrestige")
	}
	if dAtA.ResistXiongNuAddPrestige == nil {
		return errors.Errorf("%s 配置的关联字段[resist_xiong_nu_add_prestige] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("resist_xiong_nu_add_prestige"), *pArSeR)
	}

	if pArSeR.KeyExist("xiong_nu_base_destroy") {
		dAtA.XiongNuBaseDestroy = cOnFigS.GetGuildLogData(pArSeR.String("xiong_nu_base_destroy"))
	} else {
		dAtA.XiongNuBaseDestroy = cOnFigS.GetGuildLogData("XiongNuBaseDestroy")
	}
	if dAtA.XiongNuBaseDestroy == nil {
		return errors.Errorf("%s 配置的关联字段[xiong_nu_base_destroy] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("xiong_nu_base_destroy"), *pArSeR)
	}

	if pArSeR.KeyExist("yinliang_mingc_host") {
		dAtA.YinliangMingcHost = cOnFigS.GetGuildLogData(pArSeR.String("yinliang_mingc_host"))
	} else {
		dAtA.YinliangMingcHost = cOnFigS.GetGuildLogData("YinliangMingcHost")
	}
	if dAtA.YinliangMingcHost == nil {
		return errors.Errorf("%s 配置的关联字段[yinliang_mingc_host] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("yinliang_mingc_host"), *pArSeR)
	}

	if pArSeR.KeyExist("yinliang_mc_war_atk_win") {
		dAtA.YinliangMcWarAtkWin = cOnFigS.GetGuildLogData(pArSeR.String("yinliang_mc_war_atk_win"))
	} else {
		dAtA.YinliangMcWarAtkWin = cOnFigS.GetGuildLogData("YinliangMcWarAtkWin")
	}
	if dAtA.YinliangMcWarAtkWin == nil {
		return errors.Errorf("%s 配置的关联字段[yinliang_mc_war_atk_win] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("yinliang_mc_war_atk_win"), *pArSeR)
	}

	if pArSeR.KeyExist("yinliang_mc_war_atk_fail") {
		dAtA.YinliangMcWarAtkFail = cOnFigS.GetGuildLogData(pArSeR.String("yinliang_mc_war_atk_fail"))
	} else {
		dAtA.YinliangMcWarAtkFail = cOnFigS.GetGuildLogData("YinliangMcWarAtkFail")
	}
	if dAtA.YinliangMcWarAtkFail == nil {
		return errors.Errorf("%s 配置的关联字段[yinliang_mc_war_atk_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("yinliang_mc_war_atk_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("yinliang_guild_receive") {
		dAtA.YinliangGuildReceive = cOnFigS.GetGuildLogData(pArSeR.String("yinliang_guild_receive"))
	} else {
		dAtA.YinliangGuildReceive = cOnFigS.GetGuildLogData("YinliangGuildReceive")
	}
	if dAtA.YinliangGuildReceive == nil {
		return errors.Errorf("%s 配置的关联字段[yinliang_guild_receive] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("yinliang_guild_receive"), *pArSeR)
	}

	if pArSeR.KeyExist("yinliang_guild_send") {
		dAtA.YinliangGuildSend = cOnFigS.GetGuildLogData(pArSeR.String("yinliang_guild_send"))
	} else {
		dAtA.YinliangGuildSend = cOnFigS.GetGuildLogData("YinliangGuildSend")
	}
	if dAtA.YinliangGuildSend == nil {
		return errors.Errorf("%s 配置的关联字段[yinliang_guild_send] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("yinliang_guild_send"), *pArSeR)
	}

	if pArSeR.KeyExist("yinliang_guild_send_member") {
		dAtA.YinliangGuildSendMember = cOnFigS.GetGuildLogData(pArSeR.String("yinliang_guild_send_member"))
	} else {
		dAtA.YinliangGuildSendMember = cOnFigS.GetGuildLogData("YinliangGuildSendMember")
	}
	if dAtA.YinliangGuildSendMember == nil {
		return errors.Errorf("%s 配置的关联字段[yinliang_guild_send_member] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("yinliang_guild_send_member"), *pArSeR)
	}

	if pArSeR.KeyExist("yinliang_guild_pay_salary") {
		dAtA.YinliangGuildPaySalary = cOnFigS.GetGuildLogData(pArSeR.String("yinliang_guild_pay_salary"))
	} else {
		dAtA.YinliangGuildPaySalary = cOnFigS.GetGuildLogData("YinliangGuildPaySalary")
	}
	if dAtA.YinliangGuildPaySalary == nil {
		return errors.Errorf("%s 配置的关联字段[yinliang_guild_pay_salary] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("yinliang_guild_pay_salary"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_apply_atk_fail") {
		dAtA.McWarApplyAtkFail = cOnFigS.GetGuildLogData(pArSeR.String("mc_war_apply_atk_fail"))
	} else {
		dAtA.McWarApplyAtkFail = cOnFigS.GetGuildLogData("McWarApplyAtkFail")
	}
	if dAtA.McWarApplyAtkFail == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_apply_atk_fail] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_apply_atk_fail"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_apply_atk_succ") {
		dAtA.McWarApplyAtkSucc = cOnFigS.GetGuildLogData(pArSeR.String("mc_war_apply_atk_succ"))
	} else {
		dAtA.McWarApplyAtkSucc = cOnFigS.GetGuildLogData("McWarApplyAtkSucc")
	}
	if dAtA.McWarApplyAtkSucc == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_apply_atk_succ] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_apply_atk_succ"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_atk_win") {
		dAtA.McWarAtkWin = cOnFigS.GetGuildLogData(pArSeR.String("mc_war_atk_win"))
	} else {
		dAtA.McWarAtkWin = cOnFigS.GetGuildLogData("McWarAtkWin")
	}
	if dAtA.McWarAtkWin == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_atk_win] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_atk_win"), *pArSeR)
	}

	if pArSeR.KeyExist("mc_war_def_win") {
		dAtA.McWarDefWin = cOnFigS.GetGuildLogData(pArSeR.String("mc_war_def_win"))
	} else {
		dAtA.McWarDefWin = cOnFigS.GetGuildLogData("McWarDefWin")
	}
	if dAtA.McWarDefWin == nil {
		return errors.Errorf("%s 配置的关联字段[mc_war_def_win] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("mc_war_def_win"), *pArSeR)
	}

	if pArSeR.KeyExist("hufu_add") {
		dAtA.HufuAdd = cOnFigS.GetGuildLogData(pArSeR.String("hufu_add"))
	} else {
		dAtA.HufuAdd = cOnFigS.GetGuildLogData("HufuAdd")
	}
	if dAtA.HufuAdd == nil {
		return errors.Errorf("%s 配置的关联字段[hufu_add] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("hufu_add"), *pArSeR)
	}

	if pArSeR.KeyExist("baowu_add_prestige") {
		dAtA.BaowuAddPrestige = cOnFigS.GetGuildLogData(pArSeR.String("baowu_add_prestige"))
	} else {
		dAtA.BaowuAddPrestige = cOnFigS.GetGuildLogData("BaowuAddPrestige")
	}
	if dAtA.BaowuAddPrestige == nil {
		return errors.Errorf("%s 配置的关联字段[baowu_add_prestige] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("baowu_add_prestige"), *pArSeR)
	}

	if pArSeR.KeyExist("update_mark") {
		dAtA.UpdateMark = cOnFigS.GetGuildLogData(pArSeR.String("update_mark"))
	} else {
		dAtA.UpdateMark = cOnFigS.GetGuildLogData("UpdateMark")
	}
	if dAtA.UpdateMark == nil {
		return errors.Errorf("%s 配置的关联字段[update_mark] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("update_mark"), *pArSeR)
	}

	if pArSeR.KeyExist("assembly_win") {
		dAtA.AssemblyWin = cOnFigS.GetGuildLogData(pArSeR.String("assembly_win"))
	} else {
		dAtA.AssemblyWin = cOnFigS.GetGuildLogData("AssemblyWin")
	}
	if dAtA.AssemblyWin == nil {
		return errors.Errorf("%s 配置的关联字段[assembly_win] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("assembly_win"), *pArSeR)
	}

	return nil
}

// start with GuildPermissionData ----------------------------------

func NewGuildPermissionData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildPermissionData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildPermissionData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildPermissionData{}

	dAtA.InvateOther = pArSeR.Bool("invate_other")
	dAtA.AgreeJoin = pArSeR.Bool("agree_join")
	dAtA.UpdateText = pArSeR.Bool("update_text")
	dAtA.UpdateInternalText = pArSeR.Bool("update_internal_text")
	dAtA.UpdateLabel = pArSeR.Bool("update_label")
	dAtA.UpdateFriendGuild = pArSeR.Bool("update_friend_guild")
	dAtA.UpdateEnemyGuild = pArSeR.Bool("update_enemy_guild")
	dAtA.UpdateClassName = pArSeR.Bool("update_class_name")
	dAtA.UpdateFlagType = pArSeR.Bool("update_flag_type")
	dAtA.UpdateLowerMemberClassLevel = pArSeR.Bool("update_lower_member_class_level")
	dAtA.UpdateClassTitle = pArSeR.Bool("update_class_title")
	dAtA.KickLowerMember = pArSeR.Bool("kick_lower_member")
	dAtA.UpdateJoinCondition = pArSeR.Bool("update_join_condition")
	dAtA.ImpeachNpcLeader = pArSeR.Bool("impeach_npc_leader")
	dAtA.UpgradeLevel = pArSeR.Bool("upgrade_level")
	dAtA.UpgradeLevelCdr = pArSeR.Bool("upgrade_level_cdr")
	dAtA.UpgradeBuilding = pArSeR.Bool("upgrade_building")
	dAtA.UpgradeTechnology = pArSeR.Bool("upgrade_technology")
	dAtA.UpdatePrestigeTarget = pArSeR.Bool("update_prestige_target")
	dAtA.OpenResistXiongNu = pArSeR.Bool("open_resist_xiong_nu")
	dAtA.SendToAllMembers = pArSeR.Bool("send_to_all_members")
	dAtA.UpgradeTechnologyCdr = pArSeR.Bool("upgrade_technology_cdr")
	dAtA.UpdateName = pArSeR.Bool("update_name")
	dAtA.UpdateFlagName = pArSeR.Bool("update_flag_name")
	dAtA.LeaveGuild = pArSeR.Bool("leave_guild")
	dAtA.DismissGuild = pArSeR.Bool("dismiss_guild")
	dAtA.ChangeLeader = pArSeR.Bool("change_leader")
	dAtA.ChangeYinliang = pArSeR.Bool("change_yinliang")
	dAtA.UpdateMark = pArSeR.Bool("update_mark")
	dAtA.ConveneMember = pArSeR.Bool("convene_member")
	dAtA.GetOnlineInfo = pArSeR.Bool("get_online_info")
	dAtA.Workshop = pArSeR.Bool("workshop")
	dAtA.RecommendMcBuild = pArSeR.Bool("recommend_mc_build")

	return dAtA, nil
}

var vAlIdAtOrGuildPermissionData = map[string]*config.Validator{

	"invate_other":                    config.ParseValidator("bool", "", false, nil, nil),
	"agree_join":                      config.ParseValidator("bool", "", false, nil, nil),
	"update_text":                     config.ParseValidator("bool", "", false, nil, nil),
	"update_internal_text":            config.ParseValidator("bool", "", false, nil, nil),
	"update_label":                    config.ParseValidator("bool", "", false, nil, nil),
	"update_friend_guild":             config.ParseValidator("bool", "", false, nil, nil),
	"update_enemy_guild":              config.ParseValidator("bool", "", false, nil, nil),
	"update_class_name":               config.ParseValidator("bool", "", false, nil, nil),
	"update_flag_type":                config.ParseValidator("bool", "", false, nil, nil),
	"update_lower_member_class_level": config.ParseValidator("bool", "", false, nil, nil),
	"update_class_title":              config.ParseValidator("bool", "", false, nil, nil),
	"kick_lower_member":               config.ParseValidator("bool", "", false, nil, nil),
	"update_join_condition":           config.ParseValidator("bool", "", false, nil, nil),
	"impeach_npc_leader":              config.ParseValidator("bool", "", false, nil, nil),
	"upgrade_level":                   config.ParseValidator("bool", "", false, nil, nil),
	"upgrade_level_cdr":               config.ParseValidator("bool", "", false, nil, nil),
	"upgrade_building":                config.ParseValidator("bool", "", false, nil, nil),
	"upgrade_technology":              config.ParseValidator("bool", "", false, nil, nil),
	"update_prestige_target":          config.ParseValidator("bool", "", false, nil, nil),
	"open_resist_xiong_nu":            config.ParseValidator("bool", "", false, nil, nil),
	"send_to_all_members":             config.ParseValidator("bool", "", false, nil, nil),
	"upgrade_technology_cdr":          config.ParseValidator("bool", "", false, nil, nil),
	"update_name":                     config.ParseValidator("bool", "", false, nil, nil),
	"update_flag_name":                config.ParseValidator("bool", "", false, nil, nil),
	"leave_guild":                     config.ParseValidator("bool", "", false, nil, nil),
	"dismiss_guild":                   config.ParseValidator("bool", "", false, nil, nil),
	"change_leader":                   config.ParseValidator("bool", "", false, nil, nil),
	"change_yinliang":                 config.ParseValidator("bool", "", false, nil, nil),
	"update_mark":                     config.ParseValidator("bool", "", false, nil, nil),
	"convene_member":                  config.ParseValidator("bool", "", false, nil, nil),
	"get_online_info":                 config.ParseValidator("bool", "", false, nil, nil),
	"workshop":                        config.ParseValidator("bool", "", false, nil, nil),
	"recommend_mc_build":              config.ParseValidator("bool", "", false, nil, nil),
}

func (dAtA *GuildPermissionData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildPermissionData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildPermissionData) Encode() *shared_proto.GuildPermissionProto {
	out := &shared_proto.GuildPermissionProto{}
	out.InvateOther = dAtA.InvateOther
	out.AgreeJoin = dAtA.AgreeJoin
	out.UpdateText = dAtA.UpdateText
	out.UpdateInternalText = dAtA.UpdateInternalText
	out.UpdateLabel = dAtA.UpdateLabel
	out.UpdateFriendGuild = dAtA.UpdateFriendGuild
	out.UpdateEnemyGuild = dAtA.UpdateEnemyGuild
	out.UpdateClassName = dAtA.UpdateClassName
	out.UpdateFlagType = dAtA.UpdateFlagType
	out.UpdateLowerMemberClassLevel = dAtA.UpdateLowerMemberClassLevel
	out.UpdateClassTitle = dAtA.UpdateClassTitle
	out.KickLowerMember = dAtA.KickLowerMember
	out.UpdateJoinCondition = dAtA.UpdateJoinCondition
	out.ImpeachNpcLeader = dAtA.ImpeachNpcLeader
	out.UpgradeLevel = dAtA.UpgradeLevel
	out.UpgradeLevelCdr = dAtA.UpgradeLevelCdr
	out.UpgradeBuilding = dAtA.UpgradeBuilding
	out.UpgradeTechnology = dAtA.UpgradeTechnology
	out.UpdatePrestigeTarget = dAtA.UpdatePrestigeTarget
	out.OpenResistXiongNu = dAtA.OpenResistXiongNu
	out.SendToAllMembers = dAtA.SendToAllMembers
	out.UpgradeTechnologyCdr = dAtA.UpgradeTechnologyCdr
	out.UpdateName = dAtA.UpdateName
	out.UpdateFlagName = dAtA.UpdateFlagName
	out.LeaveGuild = dAtA.LeaveGuild
	out.DismissGuild = dAtA.DismissGuild
	out.ChangeLeader = dAtA.ChangeLeader
	out.ChangeYinliang = dAtA.ChangeYinliang
	out.UpdateMark = dAtA.UpdateMark
	out.ConveneMember = dAtA.ConveneMember
	out.GetOnlineInfo = dAtA.GetOnlineInfo
	out.Workshop = dAtA.Workshop
	out.RecommendMcBuild = dAtA.RecommendMcBuild

	return out
}

func ArrayEncodeGuildPermissionData(datas []*GuildPermissionData) []*shared_proto.GuildPermissionProto {

	out := make([]*shared_proto.GuildPermissionProto, 0, len(datas))
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

func (dAtA *GuildPermissionData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GuildPermissionShowData ----------------------------------

func LoadGuildPermissionShowData(gos *config.GameObjects) (map[uint64]*GuildPermissionShowData, map[*GuildPermissionShowData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildPermissionShowDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildPermissionShowData, len(lIsT))
	pArSeRmAp := make(map[*GuildPermissionShowData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildPermissionShowData) {
			continue
		}

		dAtA, err := NewGuildPermissionShowData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildPermissionShowData(dAtAmAp map[*GuildPermissionShowData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildPermissionShowDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildPermissionShowDataKeyArray(datas []*GuildPermissionShowData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildPermissionShowData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildPermissionShowData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildPermissionShowData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildPermissionShowData{}

	dAtA.PermType = shared_proto.GuildPermissionType(shared_proto.GuildPermissionType_value[strings.ToUpper(pArSeR.String("perm_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("perm_type"), 10, 32); err == nil {
		dAtA.PermType = shared_proto.GuildPermissionType(i)
	}

	dAtA.IsShow = pArSeR.Bool("is_show")
	// skip field: ClassLevel

	// calculate fields
	dAtA.Id = uint64(dAtA.PermType)

	// i18n fields
	dAtA.Name = i18n.NewI18nRef(fIlEnAmE, "name", dAtA.Id, pArSeR.String("name"))

	return dAtA, nil
}

var vAlIdAtOrGuildPermissionShowData = map[string]*config.Validator{

	"perm_type": config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.GuildPermissionType_value, 0), nil),
	"name":      config.ParseValidator("string", "", false, nil, nil),
	"is_show":   config.ParseValidator("bool", "", false, nil, nil),
}

func (dAtA *GuildPermissionShowData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildPermissionShowData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildPermissionShowData) Encode() *shared_proto.GuildPermissionShowProto {
	out := &shared_proto.GuildPermissionShowProto{}
	out.PermType = dAtA.PermType
	if dAtA.Name != nil {
		out.Name = dAtA.Name.Encode()
	}
	out.IsShow = dAtA.IsShow
	out.ClassLevel = config.U64a2I32a(dAtA.ClassLevel)

	return out
}

func ArrayEncodeGuildPermissionShowData(datas []*GuildPermissionShowData) []*shared_proto.GuildPermissionShowProto {

	out := make([]*shared_proto.GuildPermissionShowProto, 0, len(datas))
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

func (dAtA *GuildPermissionShowData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GuildPrestigeEventData ----------------------------------

func LoadGuildPrestigeEventData(gos *config.GameObjects) (map[uint64]*GuildPrestigeEventData, map[*GuildPrestigeEventData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildPrestigeEventDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildPrestigeEventData, len(lIsT))
	pArSeRmAp := make(map[*GuildPrestigeEventData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildPrestigeEventData) {
			continue
		}

		dAtA, err := NewGuildPrestigeEventData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildPrestigeEventData(dAtAmAp map[*GuildPrestigeEventData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildPrestigeEventDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildPrestigeEventDataKeyArray(datas []*GuildPrestigeEventData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildPrestigeEventData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildPrestigeEventData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildPrestigeEventData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildPrestigeEventData{}

	dAtA.TriggerEvent = shared_proto.HeroEvent(shared_proto.HeroEvent_value[strings.ToUpper(pArSeR.String("trigger_event"))])
	if i, err := strconv.ParseInt(pArSeR.String("trigger_event"), 10, 32); err == nil {
		dAtA.TriggerEvent = shared_proto.HeroEvent(i)
	}

	dAtA.TriggerEventCondition, err = data.ParseCompareCondition(pArSeR.String("trigger_event_condition"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[trigger_event_condition] 解析失败(data.ParseCompareCondition)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("trigger_event_condition"), dAtA)
	}

	dAtA.TriggerEventTimes = pArSeR.Uint64("trigger_event_times")
	dAtA.Prestige = pArSeR.Uint64("prestige")
	dAtA.Hufu = pArSeR.Uint64("hufu")
	dAtA.IgnoreMemberLimit = pArSeR.Bool("ignore_member_limit")

	// calculate fields
	dAtA.Id = GetPrestigeEventId(dAtA.TriggerEvent, dAtA.TriggerEventCondition.Amount)

	return dAtA, nil
}

var vAlIdAtOrGuildPrestigeEventData = map[string]*config.Validator{

	"trigger_event":           config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.HeroEvent_value), nil),
	"trigger_event_condition": config.ParseValidator("string", "", false, nil, nil),
	"trigger_event_times":     config.ParseValidator("uint", "", false, nil, nil),
	"prestige":                config.ParseValidator("int>0", "", false, nil, nil),
	"hufu":                    config.ParseValidator("uint", "", false, nil, nil),
	"ignore_member_limit":     config.ParseValidator("bool", "", false, nil, nil),
}

func (dAtA *GuildPrestigeEventData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GuildPrestigePrizeData ----------------------------------

func LoadGuildPrestigePrizeData(gos *config.GameObjects) (map[uint64]*GuildPrestigePrizeData, map[*GuildPrestigePrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildPrestigePrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildPrestigePrizeData, len(lIsT))
	pArSeRmAp := make(map[*GuildPrestigePrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildPrestigePrizeData) {
			continue
		}

		dAtA, err := NewGuildPrestigePrizeData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Prestige
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Prestige], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedGuildPrestigePrizeData(dAtAmAp map[*GuildPrestigePrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildPrestigePrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildPrestigePrizeDataKeyArray(datas []*GuildPrestigePrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Prestige)
		}
	}

	return out
}

func NewGuildPrestigePrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildPrestigePrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildPrestigePrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildPrestigePrizeData{}

	dAtA.Prestige = pArSeR.Uint64("prestige")
	// releated field: EventPrize
	dAtA.BuildingAmount = pArSeR.Uint64("building_amount")
	dAtA.Hufu = pArSeR.Uint64("hufu")

	return dAtA, nil
}

var vAlIdAtOrGuildPrestigePrizeData = map[string]*config.Validator{

	"prestige":        config.ParseValidator("int>0", "", false, nil, nil),
	"event_prize":     config.ParseValidator("string", "", false, nil, nil),
	"building_amount": config.ParseValidator("int>0", "", false, nil, nil),
	"hufu":            config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *GuildPrestigePrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildPrestigePrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildPrestigePrizeData) Encode() *shared_proto.GuildPrestigePrizeDataProto {
	out := &shared_proto.GuildPrestigePrizeDataProto{}
	out.Prestige = config.U64ToI32(dAtA.Prestige)
	if dAtA.EventPrize != nil {
		out.EventPrize = config.U64ToI32(dAtA.EventPrize.Id)
	}
	out.BuildingAmount = config.U64ToI32(dAtA.BuildingAmount)
	out.Hufu = config.U64ToI32(dAtA.Hufu)

	return out
}

func ArrayEncodeGuildPrestigePrizeData(datas []*GuildPrestigePrizeData) []*shared_proto.GuildPrestigePrizeDataProto {

	out := make([]*shared_proto.GuildPrestigePrizeDataProto, 0, len(datas))
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

func (dAtA *GuildPrestigePrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.EventPrize = cOnFigS.GetGuildEventPrizeData(pArSeR.Uint64("event_prize"))
	if dAtA.EventPrize == nil {
		return errors.Errorf("%s 配置的关联字段[event_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("event_prize"), *pArSeR)
	}

	return nil
}

// start with GuildRankPrizeData ----------------------------------

func LoadGuildRankPrizeData(gos *config.GameObjects) (map[uint64]*GuildRankPrizeData, map[*GuildRankPrizeData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildRankPrizeDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildRankPrizeData, len(lIsT))
	pArSeRmAp := make(map[*GuildRankPrizeData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildRankPrizeData) {
			continue
		}

		dAtA, err := NewGuildRankPrizeData(fIlEnAmE, pArSeR)
		if err != nil {
			return nil, nil, err
		}

		key := dAtA.Rank
		if dAtAmAp[key] != nil {
			return nil, nil, errors.Errorf("%s 表中存在重复的Key字段[Rank], key: %s", fIlEnAmE, key)
		}

		dAtAmAp[key] = dAtA
		pArSeRmAp[dAtA] = pArSeR
	}

	return dAtAmAp, pArSeRmAp, nil
}

func SetRelatedGuildRankPrizeData(dAtAmAp map[*GuildRankPrizeData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildRankPrizeDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildRankPrizeDataKeyArray(datas []*GuildRankPrizeData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Rank)
		}
	}

	return out
}

func NewGuildRankPrizeData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildRankPrizeData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildRankPrizeData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildRankPrizeData{}

	dAtA.Rank = pArSeR.Uint64("rank")
	// releated field: Prize
	// releated field: CountryDestroyPrize

	return dAtA, nil
}

var vAlIdAtOrGuildRankPrizeData = map[string]*config.Validator{

	"rank":                  config.ParseValidator("int>0", "", false, nil, nil),
	"prize":                 config.ParseValidator("string", "", false, nil, nil),
	"country_destroy_prize": config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *GuildRankPrizeData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildRankPrizeData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildRankPrizeData) Encode() *shared_proto.GuildRankPrizeDataProto {
	out := &shared_proto.GuildRankPrizeDataProto{}
	out.Rank = config.U64ToI32(dAtA.Rank)
	if dAtA.Prize != nil {
		out.Prize = dAtA.Prize.Encode()
	}
	if dAtA.CountryDestroyPrize != nil {
		out.CountryDestroyPrize = dAtA.CountryDestroyPrize.Encode()
	}

	return out
}

func ArrayEncodeGuildRankPrizeData(datas []*GuildRankPrizeData) []*shared_proto.GuildRankPrizeDataProto {

	out := make([]*shared_proto.GuildRankPrizeDataProto, 0, len(datas))
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

func (dAtA *GuildRankPrizeData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

	dAtA.CountryDestroyPrize = cOnFigS.GetPrize(pArSeR.Int("country_destroy_prize"))
	if dAtA.CountryDestroyPrize == nil {
		return errors.Errorf("%s 配置的关联字段[country_destroy_prize] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("country_destroy_prize"), *pArSeR)
	}

	return nil
}

// start with GuildTarget ----------------------------------

func LoadGuildTarget(gos *config.GameObjects) (map[uint64]*GuildTarget, map[*GuildTarget]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildTargetPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildTarget, len(lIsT))
	pArSeRmAp := make(map[*GuildTarget]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildTarget) {
			continue
		}

		dAtA, err := NewGuildTarget(fIlEnAmE, pArSeR)
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

func SetRelatedGuildTarget(dAtAmAp map[*GuildTarget]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildTargetPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildTargetKeyArray(datas []*GuildTarget) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildTarget(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildTarget, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildTarget)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildTarget{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Group = pArSeR.Uint64("group")
	dAtA.ButtonText = pArSeR.String("button_text")
	dAtA.TargetType = shared_proto.GuildTargetType(shared_proto.GuildTargetType_value[strings.ToUpper(pArSeR.String("target_type"))])
	if i, err := strconv.ParseInt(pArSeR.String("target_type"), 10, 32); err == nil {
		dAtA.TargetType = shared_proto.GuildTargetType(i)
	}

	dAtA.Target = pArSeR.Uint64("target")
	dAtA.Order = pArSeR.Uint64("order")
	// skip field: OrderAmount

	return dAtA, nil
}

var vAlIdAtOrGuildTarget = map[string]*config.Validator{

	"id":          config.ParseValidator("int>0", "", false, nil, nil),
	"name":        config.ParseValidator("string>0", "", false, nil, nil),
	"desc":        config.ParseValidator("string>0", "", false, nil, nil),
	"icon":        config.ParseValidator("string", "", false, nil, nil),
	"group":       config.ParseValidator("int>0", "", false, nil, nil),
	"button_text": config.ParseValidator("string", "", false, nil, nil),
	"target_type": config.ParseValidator("string,notAllNil", "", false, config.EnumMapKeys(shared_proto.GuildTargetType_value, 0), nil),
	"target":      config.ParseValidator("uint", "", false, nil, nil),
	"order":       config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *GuildTarget) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildTarget) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildTarget) Encode() *shared_proto.GuildTargetProto {
	out := &shared_proto.GuildTargetProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.Icon = dAtA.Icon
	out.ButtonText = dAtA.ButtonText
	out.TargetType = dAtA.TargetType
	out.Target = config.U64ToI32(dAtA.Target)

	return out
}

func ArrayEncodeGuildTarget(datas []*GuildTarget) []*shared_proto.GuildTargetProto {

	out := make([]*shared_proto.GuildTargetProto, 0, len(datas))
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

func (dAtA *GuildTarget) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with GuildTaskData ----------------------------------

func LoadGuildTaskData(gos *config.GameObjects) (map[uint64]*GuildTaskData, map[*GuildTaskData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildTaskDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildTaskData, len(lIsT))
	pArSeRmAp := make(map[*GuildTaskData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildTaskData) {
			continue
		}

		dAtA, err := NewGuildTaskData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildTaskData(dAtAmAp map[*GuildTaskData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildTaskDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildTaskDataKeyArray(datas []*GuildTaskData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildTaskData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildTaskData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildTaskData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildTaskData{}

	dAtA.Id = pArSeR.Uint64("id")
	// skip field: TaskType
	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.String("desc")
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Stages = pArSeR.Uint64Array("stages", "", false)
	// releated field: Prizes

	return dAtA, nil
}

var vAlIdAtOrGuildTaskData = map[string]*config.Validator{

	"id":     config.ParseValidator("int>0", "", false, nil, nil),
	"name":   config.ParseValidator("string", "", false, nil, nil),
	"desc":   config.ParseValidator("string", "", false, nil, nil),
	"icon":   config.ParseValidator("string", "", false, nil, nil),
	"stages": config.ParseValidator("uint,notAllNil", "", true, nil, nil),
	"prizes": config.ParseValidator("uint,notAllNil,duplicate", "", true, nil, nil),
}

func (dAtA *GuildTaskData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildTaskData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildTaskData) Encode() *shared_proto.GuildTaskDataProto {
	out := &shared_proto.GuildTaskDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.Icon = dAtA.Icon
	out.Stages = config.U64a2I32a(dAtA.Stages)
	if dAtA.Prizes != nil {
		out.Prizes = resdata.ArrayEncodePrize(dAtA.Prizes)
	}

	return out
}

func ArrayEncodeGuildTaskData(datas []*GuildTaskData) []*shared_proto.GuildTaskDataProto {

	out := make([]*shared_proto.GuildTaskDataProto, 0, len(datas))
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

func (dAtA *GuildTaskData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	intKeys = pArSeR.IntArray("prizes", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetPrize(v)
		if obj != nil {
			dAtA.Prizes = append(dAtA.Prizes, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[prizes] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prizes"), *pArSeR)
		}
	}

	return nil
}

// start with GuildTaskEvaluateData ----------------------------------

func LoadGuildTaskEvaluateData(gos *config.GameObjects) (map[uint64]*GuildTaskEvaluateData, map[*GuildTaskEvaluateData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildTaskEvaluateDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildTaskEvaluateData, len(lIsT))
	pArSeRmAp := make(map[*GuildTaskEvaluateData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildTaskEvaluateData) {
			continue
		}

		dAtA, err := NewGuildTaskEvaluateData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildTaskEvaluateData(dAtAmAp map[*GuildTaskEvaluateData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildTaskEvaluateDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildTaskEvaluateDataKeyArray(datas []*GuildTaskEvaluateData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildTaskEvaluateData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildTaskEvaluateData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildTaskEvaluateData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildTaskEvaluateData{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.Complete = pArSeR.Uint64("complete")
	// releated field: Prizes

	return dAtA, nil
}

var vAlIdAtOrGuildTaskEvaluateData = map[string]*config.Validator{

	"id":       config.ParseValidator("int>0", "", false, nil, nil),
	"name":     config.ParseValidator("string", "", false, nil, nil),
	"complete": config.ParseValidator("int>0", "", false, nil, nil),
	"prizes":   config.ParseValidator("uint,notAllNil,duplicate", "", true, nil, nil),
}

func (dAtA *GuildTaskEvaluateData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildTaskEvaluateData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildTaskEvaluateData) Encode() *shared_proto.GuildTaskEvaluateDataProto {
	out := &shared_proto.GuildTaskEvaluateDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Complete = config.U64ToI32(dAtA.Complete)
	if dAtA.Prizes != nil {
		out.Prizes = resdata.ArrayEncodePrize(dAtA.Prizes)
	}

	return out
}

func ArrayEncodeGuildTaskEvaluateData(datas []*GuildTaskEvaluateData) []*shared_proto.GuildTaskEvaluateDataProto {

	out := make([]*shared_proto.GuildTaskEvaluateDataProto, 0, len(datas))
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

func (dAtA *GuildTaskEvaluateData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	intKeys = pArSeR.IntArray("prizes", "", false)
	for _, v := range intKeys {
		obj := cOnFigS.GetPrize(v)
		if obj != nil {
			dAtA.Prizes = append(dAtA.Prizes, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[prizes] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("prizes"), *pArSeR)
		}
	}

	return nil
}

// start with GuildTechnologyData ----------------------------------

func LoadGuildTechnologyData(gos *config.GameObjects) (map[uint64]*GuildTechnologyData, map[*GuildTechnologyData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.GuildTechnologyDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*GuildTechnologyData, len(lIsT))
	pArSeRmAp := make(map[*GuildTechnologyData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrGuildTechnologyData) {
			continue
		}

		dAtA, err := NewGuildTechnologyData(fIlEnAmE, pArSeR)
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

func SetRelatedGuildTechnologyData(dAtAmAp map[*GuildTechnologyData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.GuildTechnologyDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetGuildTechnologyDataKeyArray(datas []*GuildTechnologyData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewGuildTechnologyData(fIlEnAmE string, pArSeR *config.ObjectParser) (*GuildTechnologyData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrGuildTechnologyData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &GuildTechnologyData{}

	dAtA.Name = pArSeR.String("name")
	dAtA.Desc = pArSeR.StringArray("desc", "", false)
	dAtA.Icon = pArSeR.String("icon")
	dAtA.Group = pArSeR.Uint64("group")
	dAtA.Level = pArSeR.Uint64("level")
	dAtA.RequireGuildLevel = pArSeR.Uint64("require_guild_level")
	dAtA.UpgradeBuilding = pArSeR.Uint64("upgrade_building")
	dAtA.UpgradeDuration, err = config.ParseDuration(pArSeR.String("upgrade_duration"))
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[upgrade_duration] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("upgrade_duration"), dAtA)
	}

	// skip field: Cdrs
	if pArSeR.KeyExist("help_cdr") {
		dAtA.HelpCdr, err = config.ParseDuration(pArSeR.String("help_cdr"))
	} else {
		dAtA.HelpCdr, err = config.ParseDuration("1m")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "%s (行数: %s) 配置的字段[help_cdr] 解析失败(config.ParseDuration)，%s, %s", fIlEnAmE, pArSeR.Line(), pArSeR.OriginStringArray("help_cdr"), dAtA)
	}

	// releated field: Effect
	// releated field: BigBox

	// calculate fields
	dAtA.Id = GetTechnologyDataId(dAtA.Group, dAtA.Level)

	return dAtA, nil
}

var vAlIdAtOrGuildTechnologyData = map[string]*config.Validator{

	"name":                config.ParseValidator("string", "", false, nil, nil),
	"desc":                config.ParseValidator("string,notAllNil,duplicate", "", true, nil, nil),
	"icon":                config.ParseValidator("string", "", false, nil, nil),
	"group":               config.ParseValidator("int>0", "", false, nil, nil),
	"level":               config.ParseValidator("int>0", "", false, nil, nil),
	"require_guild_level": config.ParseValidator("int>0", "", false, nil, nil),
	"upgrade_building":    config.ParseValidator("int>0", "", false, nil, nil),
	"upgrade_duration":    config.ParseValidator("string", "", false, nil, nil),
	"help_cdr":            config.ParseValidator("string", "", false, nil, []string{"1m"}),
	"effect":              config.ParseValidator("string", "", false, nil, nil),
	"big_box":             config.ParseValidator("string", "", false, nil, nil),
}

func (dAtA *GuildTechnologyData) Marshal() ([]byte, error) {
	return dAtA.Encode().Marshal()
}

func (dAtA *GuildTechnologyData) MarshalTo(data []byte) (int, error) {
	return dAtA.Encode().MarshalTo(data)
}

func (dAtA *GuildTechnologyData) Encode() *shared_proto.GuildTechnologyDataProto {
	out := &shared_proto.GuildTechnologyDataProto{}
	out.Id = config.U64ToI32(dAtA.Id)
	out.Name = dAtA.Name
	out.Desc = dAtA.Desc
	out.Icon = dAtA.Icon
	out.Group = config.U64ToI32(dAtA.Group)
	out.Level = config.U64ToI32(dAtA.Level)
	out.RequireGuildLevel = config.U64ToI32(dAtA.RequireGuildLevel)
	out.UpgradeBuilding = config.U64ToI32(dAtA.UpgradeBuilding)
	out.UpgradeDuration = config.Duration2I32Seconds(dAtA.UpgradeDuration)
	if dAtA.Cdrs != nil {
		out.Cdrs = ArrayEncodeGuildLevelCdrData(dAtA.Cdrs)
	}
	out.HelpCdr = config.Duration2I32Seconds(dAtA.HelpCdr)
	if dAtA.Effect != nil {
		out.Effect = dAtA.Effect.Encode()
	}
	if dAtA.BigBox != nil {
		out.BigBox = config.U64ToI32(dAtA.BigBox.Id)
	}

	return out
}

func ArrayEncodeGuildTechnologyData(datas []*GuildTechnologyData) []*shared_proto.GuildTechnologyDataProto {

	out := make([]*shared_proto.GuildTechnologyDataProto, 0, len(datas))
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

func (dAtA *GuildTechnologyData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Effect = cOnFigS.GetBuildingEffectData(pArSeR.Int("effect"))
	if dAtA.Effect == nil && pArSeR.Int("effect") != 0 {
		return errors.Errorf("%s 配置的关联字段[effect] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("effect"), *pArSeR)
	}

	dAtA.BigBox = cOnFigS.GetGuildBigBoxData(pArSeR.Uint64("big_box"))
	if dAtA.BigBox == nil && pArSeR.Uint64("big_box") != 0 {
		return errors.Errorf("%s 配置的关联字段[big_box] 填的值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("big_box"), *pArSeR)
	}

	return nil
}

// start with NpcGuildSuffixName ----------------------------------

func LoadNpcGuildSuffixName(gos *config.GameObjects) (*NpcGuildSuffixName, *config.ObjectParser, error) {
	fIlEnAmE := confpath.NpcGuildSuffixNamePath
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

	dAtA, err := NewNpcGuildSuffixName(fIlEnAmE, pArSeR)
	return dAtA, pArSeR, err
}

func SetRelatedNpcGuildSuffixName(gos *config.GameObjects, dAtA *NpcGuildSuffixName, cOnFigS interface{}) error {
	fIlEnAmE := confpath.NpcGuildSuffixNamePath
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

func NewNpcGuildSuffixName(fIlEnAmE string, pArSeR *config.ObjectParser) (*NpcGuildSuffixName, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrNpcGuildSuffixName)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &NpcGuildSuffixName{}

	dAtA.Name = pArSeR.StringArray("name", "", false)

	return dAtA, nil
}

var vAlIdAtOrNpcGuildSuffixName = map[string]*config.Validator{

	"name": config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *NpcGuildSuffixName) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
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

// start with NpcGuildTemplate ----------------------------------

func LoadNpcGuildTemplate(gos *config.GameObjects) (map[uint64]*NpcGuildTemplate, map[*NpcGuildTemplate]*config.ObjectParser, error) {
	fIlEnAmE := confpath.NpcGuildTemplatePath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*NpcGuildTemplate, len(lIsT))
	pArSeRmAp := make(map[*NpcGuildTemplate]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrNpcGuildTemplate) {
			continue
		}

		dAtA, err := NewNpcGuildTemplate(fIlEnAmE, pArSeR)
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

func SetRelatedNpcGuildTemplate(dAtAmAp map[*NpcGuildTemplate]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.NpcGuildTemplatePath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetNpcGuildTemplateKeyArray(datas []*NpcGuildTemplate) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewNpcGuildTemplate(fIlEnAmE string, pArSeR *config.ObjectParser) (*NpcGuildTemplate, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrNpcGuildTemplate)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &NpcGuildTemplate{}

	dAtA.Id = pArSeR.Uint64("id")
	dAtA.Name = pArSeR.String("name")
	dAtA.FlagName = pArSeR.String("flag_name")
	dAtA.Text = pArSeR.String("text")
	dAtA.InternalText = pArSeR.String("internal_text")
	dAtA.Labels = pArSeR.StringArray("labels", "", false)
	// releated field: Level
	// releated field: Country
	dAtA.RejectUserJoin = pArSeR.Bool("reject_user_join")
	dAtA.NpcLeaderVote = pArSeR.Uint64("npc_leader_vote")
	// releated field: Leader
	// releated field: Members

	return dAtA, nil
}

var vAlIdAtOrNpcGuildTemplate = map[string]*config.Validator{

	"id":               config.ParseValidator("int>0", "", false, nil, nil),
	"name":             config.ParseValidator("string", "", false, nil, nil),
	"flag_name":        config.ParseValidator("string", "", false, nil, nil),
	"text":             config.ParseValidator("string", "", false, nil, nil),
	"internal_text":    config.ParseValidator("string", "", false, nil, nil),
	"labels":           config.ParseValidator("string", "", true, nil, nil),
	"level":            config.ParseValidator("string", "", false, nil, nil),
	"country":          config.ParseValidator("string", "", false, nil, nil),
	"reject_user_join": config.ParseValidator("bool", "", false, nil, nil),
	"npc_leader_vote":  config.ParseValidator("int>0", "", false, nil, nil),
	"leader":           config.ParseValidator("string", "", false, nil, nil),
	"members":          config.ParseValidator("string", "", true, nil, nil),
}

func (dAtA *NpcGuildTemplate) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Level = cOnFigS.GetGuildLevelData(pArSeR.Uint64("level"))
	if dAtA.Level == nil {
		return errors.Errorf("%s 配置的关联字段[level] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("level"), *pArSeR)
	}

	dAtA.Country = cOnFigS.GetCountryData(pArSeR.Uint64("country"))
	if dAtA.Country == nil {
		return errors.Errorf("%s 配置的关联字段[country] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("country"), *pArSeR)
	}

	dAtA.Leader = cOnFigS.GetNpcMemberData(pArSeR.Uint64("leader"))
	if dAtA.Leader == nil {
		return errors.Errorf("%s 配置的关联字段[leader] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("leader"), *pArSeR)
	}

	uint64Keys = pArSeR.Uint64Array("members", "", false)
	for _, v := range uint64Keys {
		obj := cOnFigS.GetNpcMemberData(v)
		if obj != nil {
			dAtA.Members = append(dAtA.Members, obj)
		} else {
			return errors.Errorf("%s 配置的关联字段[members] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("members"), *pArSeR)
		}
	}

	return nil
}

// start with NpcMemberData ----------------------------------

func LoadNpcMemberData(gos *config.GameObjects) (map[uint64]*NpcMemberData, map[*NpcMemberData]*config.ObjectParser, error) {
	fIlEnAmE := confpath.NpcMemberDataPath
	lIsT, err := gos.LoadFile(fIlEnAmE)
	if err != nil {
		return nil, nil, err
	}

	if len(lIsT) <= 0 {
		return nil, nil, errors.Errorf("%s 表中没有数据", fIlEnAmE)
	}

	dAtAmAp := make(map[uint64]*NpcMemberData, len(lIsT))
	pArSeRmAp := make(map[*NpcMemberData]*config.ObjectParser, len(lIsT))
	for _, pArSeR := range lIsT {
		if pArSeR.IsEmpty(vAlIdAtOrNpcMemberData) {
			continue
		}

		dAtA, err := NewNpcMemberData(fIlEnAmE, pArSeR)
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

func SetRelatedNpcMemberData(dAtAmAp map[*NpcMemberData]*config.ObjectParser, cOnFigS interface{}) error {
	fIlEnAmE := confpath.NpcMemberDataPath
	for dAtA, pArSeR := range dAtAmAp {
		if err := dAtA.SetRelatedObject(fIlEnAmE, pArSeR, cOnFigS); err != nil {
			return err
		}
	}

	return nil
}

func GetNpcMemberDataKeyArray(datas []*NpcMemberData) []uint64 {

	out := make([]uint64, 0, len(datas))
	for _, d := range datas {
		if d != nil {
			out = append(out, d.Id)
		}
	}

	return out
}

func NewNpcMemberData(fIlEnAmE string, pArSeR *config.ObjectParser) (*NpcMemberData, error) {

	err := pArSeR.Validate(fIlEnAmE, vAlIdAtOrNpcMemberData)
	if err != nil {
		return nil, err
	}

	var stringKeys []string
	if len(stringKeys) > 0 {
	}

	dAtA := &NpcMemberData{}

	dAtA.Id = pArSeR.Uint64("id")
	// releated field: Master
	dAtA.ContributionAmount = pArSeR.Uint64("contribution_amount")
	dAtA.ContributionAmount7 = pArSeR.Uint64("contribution_amount7")
	dAtA.TotalContributionAmount = pArSeR.Uint64("total_contribution_amount")

	return dAtA, nil
}

var vAlIdAtOrNpcMemberData = map[string]*config.Validator{

	"id":                        config.ParseValidator("int>0", "", false, nil, nil),
	"master":                    config.ParseValidator("string", "", false, nil, nil),
	"contribution_amount":       config.ParseValidator("int>0", "", false, nil, nil),
	"contribution_amount7":      config.ParseValidator("int>0", "", false, nil, nil),
	"total_contribution_amount": config.ParseValidator("int>0", "", false, nil, nil),
}

func (dAtA *NpcMemberData) SetRelatedObject(fIlEnAmE string, pArSeR *config.ObjectParser, cOnFigS0 interface{}) error {
	cOnFigS := cOnFigS0.(related_configs)
	if cOnFigS == nil {
	}

	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtA.Master = cOnFigS.GetMonsterMasterData(pArSeR.Uint64("master"))
	if dAtA.Master == nil {
		return errors.Errorf("%s 配置的关联字段[master] 必填，没有填值或者值在关联表中没找到，填的值是[%v]，这行数据: %v", fIlEnAmE, pArSeR.OriginStringArray("master"), *pArSeR)
	}

	return nil
}

type related_configs interface {
	GetBuildingEffectData(int) *sub.BuildingEffectData
	GetCost(int) *resdata.Cost
	GetCountryData(uint64) *country.CountryData
	GetGuildBigBoxData(uint64) *GuildBigBoxData
	GetGuildEventPrizeData(uint64) *GuildEventPrizeData
	GetGuildLevelData(uint64) *GuildLevelData
	GetGuildLogData(string) *GuildLogData
	GetIcon(string) *icon.Icon
	GetMonsterMasterData(uint64) *monsterdata.MonsterMasterData
	GetNpcMemberData(uint64) *NpcMemberData
	GetPlunderPrize(uint64) *resdata.PlunderPrize
	GetPrize(int) *resdata.Prize
}
