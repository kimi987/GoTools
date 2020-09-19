package country

import (
	"github.com/lightpaw/male7/config/body"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/head"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"time"
)

// 默认国家配置: 七雄+周
//gogen:config
type CountryData struct {
	_ struct{} `file:"国家/国家.txt"`
	_ struct{} `protogen:"true"`

	Id uint64 `desc:"国家id"`

	Name string `desc:"国家名字"`

	Desc string `default:" "` // 国家描述

	DefaultPrestige uint64 `validator:"uint" default:"0"` // 初始声望

	Capital *mingcdata.MingcBaseData `desc:"都城 MingcBaseDataProto.Id" protofield:",config.U64ToI32(%s.Id),int32"`

	BornCenterX uint64 `validator:"uint" protofield:"-"`

	BornCenterY uint64 `validator:"uint" protofield:"-"`

	BornRadiusX uint64 `validator:"uint" protofield:"-"`

	BornRadiusY uint64 `validator:"uint" protofield:"-"`

	NpcOfficial []shared_proto.CountryOfficialType `desc:"Npc 官职" validator:",duplicate"`
	NpcId       []uint64                           `desc:"Npc id CountryOfficialNpcData.Id，与 NpcOfficial 一一对应" validator:"uint,duplicate"`
}

func (*CountryData) InitAll(filename string, configs interface {
	GetCountryDataArray() []*CountryData
	GetCountryData(uint64) *CountryData
	GetMingcBaseDataArray() []*mingcdata.MingcBaseData
}) {
	for idx, c := range configs.GetCountryDataArray() {
		check.PanicNotTrue(c.Id == uint64(idx+1), "%s 国家的id只能够从1开始每次加1: %d-%s", filename, c.Id, c.Name)
		check.PanicNotTrue(c.Capital.Country == c.Id, "%s 国家的都城配置错误: %d-%s", filename, c.Id, c.Capital.Id)
		check.PanicNotTrue(c.Capital.Type == shared_proto.MincType_MC_DU, "%s 国家的都城配置错误: %d-%s", filename, c.Id, c.Capital.Id)
	}

	// 名城配置验证
	var luoyangCount uint64
	for _, mc := range configs.GetMingcBaseDataArray() {
		check.PanicNotTrue(mc.Id <= npcid.NpcDataMask, "%s npc城池的配置数据的id最大不能超过 %d, id: %d", filename, npcid.NpcDataMask, mc.Id)

		if mc.ZhouCaptain {
			check.PanicNotTrue(mc.Country == 0, "%v，zhou_captain:1 的 country 才能配成0.country:%v", filename, mc.Country)
			luoyangCount++
		} else {
			check.PanicNotTrue(configs.GetCountryData(mc.Country) != nil, "%v，找不到 id:%v 的 country:%v", filename, mc.Id, mc.Country)
		}
	}
	check.PanicNotTrue(luoyangCount == 1, "%v，zhou_captain:1 只能配一个", filename)
}

//gogen:config
type CountryMiscData struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"国家/国家杂项.txt"`
	_ struct{} `protogen:"true"`

	// 转国消耗物品
	NormalChangeCountryGoods *goods.GoodsData `protofield:",config.U64ToI32(%s.Id),int32"`

	// 新手转国 cd
	NewHeroChangeCountryCd time.Duration

	// 普通转国 cd
	NormalChangeCountryCd time.Duration

	// 新手最大君主等级
	NewHeroMaxLevel uint64

	ChangeNameVoteDuration time.Duration `desc:"更改名字的投票时间"`
	ChangeNameCost         *resdata.Cost `desc:"更改名字的消耗"`
	ChangeNameCd           time.Duration `desc:"改国号CD" default:"168h"`

	MaxSearchHeroDefaultCount int `desc:"任命玩家默认列表最大长度" default:"20"`
	MaxSearchHeroByNameCount  int `desc:"搜索任命玩家列表最大长度" default:"200"`
}

//gogen:config
type CountryOfficialData struct {
	_ struct{} `file:"国家/官职.txt"`
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"country.proto"`
	_ struct{} `protoimport:"base.proto"`
	_ struct{} `protoimport:"task.proto"`

	Id             int                              `protofield:"-"`
	OfficialType   shared_proto.CountryOfficialType `desc:"key。官职类型" head:"-"`
	Name           string                           `desc:"名字"`
	BuildingEffect *sub.BuildingEffectData          `desc:"内政技能" default:"nullable"`
	Buff           *data.BuffEffectData             `desc:"战斗buff" default:"nullable"`
	Count          int                              `desc:"数量"`
	Salary         *resdata.Plunder                 `desc:"俸禄" protofield:"-"`
	ShowSalary     *resdata.Prize                   `desc:"展示俸禄"`
	Icon           *icon.Icon                       `desc:"图标 IconProto.Id" protofield:",%s.Id,string"`
	Head           *head.HeadData                   `desc:"外显头像 HeadDataProto.Id。空表示没有特殊外显" default:"nullable" protofield:",%s.Id,string"`
	Body           *body.BodyData                   `desc:"外显形象 BodyDataProto.Id。0表示没有特殊外显" default:"nullable" protofield:",config.U64ToI32(%s.Id),int32"`
	Cd             time.Duration                    `desc:"任命CD,国王行表示禅让CD "`
	EffectDesc     string                           `desc:"官职效果描述" default:""`

	SubOfficials []shared_proto.CountryOfficialType `desc:"下属类型" validator:",duplicate"`
}

func (data *CountryOfficialData) InitAll(filename string, conf interface {
	GetCountryOfficialDataArray() []*CountryOfficialData
	GetCountryOfficialData(int) *CountryOfficialData
}) {
	// 每个 type 都要配置
	for id := range shared_proto.CountryOfficialType_name {
		if id <= 0 {
			continue
		}
		check.PanicNotTrue(conf.GetCountryOfficialData(int(id)) != nil, "%v, 国家官职类型 t:%v 没有配置", filename, id)
	}

	for _, d := range conf.GetCountryOfficialDataArray() {
		d.OfficialType = shared_proto.CountryOfficialType(d.Id)
		switch d.OfficialType {
		case shared_proto.CountryOfficialType_COT_KING, shared_proto.CountryOfficialType_COT_QUEEN:
			check.PanicNotTrue(d.Count == 1, "%v, 官职类型:%v 的 count 必须 == 1.", filename, d.OfficialType)
		}
		for _, subType := range d.SubOfficials {
			check.PanicNotTrue(subType > d.OfficialType, "%v, 官职类型:%v 的下属:%v 必须比自己官职低。", filename, d.OfficialType, subType)
		}
	}
}

func (d *CountryOfficialData) IsSubOfficial(t shared_proto.CountryOfficialType) bool {
	for _, subT := range d.SubOfficials {
		if t == subT {
			return true
		}
	}
	return false
}

//gogen:config
type CountryOfficialNpcData struct {
	_ struct{} `file:"国家/官职npc.txt"`
	_ struct{} `protogen:"true"`

	Id   uint64
	Name string         `desc:"名字"`
	Head *head.HeadData `desc:"外显头像 HeadDataProto.Id" protofield:",%s.Id,string"`
	Body *body.BodyData `desc:"外显形象 BodyDataProto.Id" protofield:",config.U64ToI32(%s.Id),int32"`
}
