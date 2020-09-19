package goods

import (
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"time"
	"github.com/lightpaw/male7/config/data"
)

func GetGoodsName(g *GoodsData) string {
	if g != nil {
		return g.Name
	}
	return ""
}

//gogen:config
type GoodsData struct {
	_  struct{} `file:"物品/物品.txt"`
	_  struct{} `proto:"shared_proto.GoodsDataProto"`
	_  struct{} `protoconfig:"goods"`
	Id uint64

	Name string

	Desc string

	Icon *icon.Icon `protofield:"IconId,%s.Id"` // 图标

	ObtainWays []string // 获得的途径

	Quality shared_proto.Quality `head:"-"` // GoodsQuality

	GoodsQuality *GoodsQuality `protofield:",config.U64ToI32(%s.Level)"`

	YuanbaoPrice  uint64 `validator:"uint" default:"0"`
	DianquanPrice uint64 `validator:"uint" default:"0"`
	YinliangPrice uint64 `validator:"uint" default:"0"`

	Cd      time.Duration `default:"0s"`
	Captain uint64        `validator:"uint" default:"0"`

	// 效果
	EffectType  shared_proto.GoodsEffectType `head:"-" type:"enum"`
	GoodsEffect *GoodsEffect                 `type:"sub"`

	SpecType shared_proto.GoodsSpecType `validator:"string" type:"enum" white:"0"` // 特殊类型物品

	// 合璧类型
	HebiType shared_proto.HebiType `validator:"string" type:"enum" white:"0"`
	// 合璧子类型
	HebiSubType shared_proto.HebiSubType `validator:"string" type:"enum" white:"0"`

	// 匹配玉璧
	PartnerHebiGoods uint64 `validator:"uint" default:"0"`

	PartnerHebiGoodsData *GoodsData `head:"-" protofield:"-"`
}

func (g *GoodsData) Init(filename string, configs interface {
	GetGoodsData(uint64) *GoodsData
	GetBuffEffectData(uint64) *data.BuffEffectData
}) {
	g.Quality = g.GoodsQuality.Quality

	g.GoodsEffect.Init(filename)
	g.EffectType = g.GoodsEffect.EffectType
	if g.EffectType == shared_proto.GoodsEffectType_Normal {
		g.GoodsEffect = nil
	}

	check.PanicNotTrue(!(g.YuanbaoPrice > 0 && g.DianquanPrice > 0), "%v, 元宝价格或点券价格只能配一个, id:%v", filename, g.Id)

	if g.SpecType == shared_proto.GoodsSpecType_GAT_HEBI {
		check.PanicNotTrue(g.HebiType > 0, "%v, 合璧物品 %v 必须配置 hebi_type", filename, g.Id)
		check.PanicNotTrue(g.HebiSubType > 0, "%v, 合璧物品 %v 必须配置 hebi_sub_type", filename, g.Id)

		g.PartnerHebiGoodsData = configs.GetGoodsData(g.PartnerHebiGoods)
		check.PanicNotTrue(g.PartnerHebiGoodsData != nil, "%v, 合璧物品 %v partner_hebi_goods:%v 不存在", filename, g.Id, g.PartnerHebiGoods)
	}

	if g.GoodsEffect != nil && g.GoodsEffect.BuffId > 0 {
		check.PanicNotTrue(configs.GetBuffEffectData(g.GoodsEffect.BuffId) != nil, "%v, buff 物品 %v buff_id：%v 不存在。", filename, g.Id, g.GoodsEffect.BuffId)
	}
}

//gogen:config
type GoodsQuality struct {
	_ struct{} `file:"物品/物品品质.txt"`
	_ struct{} `proto:"shared_proto.GoodsQualityProto"`
	_ struct{} `protoconfig:"goods_quality"`

	Level   uint64 `key:"true" validator:"uint"`
	Quality shared_proto.Quality
}

//gogen:config
type GoodsEffect struct {
	_ struct{} `proto:"shared_proto.GoodsEffectProto"`

	Id int `protofield:"-"`

	EffectType shared_proto.GoodsEffectType `head:"-" protofield:"-"`

	// 加资源类型效果
	Gold  uint64 `validator:"uint"`
	Food  uint64 `validator:"uint"`
	Wood  uint64 `validator:"uint"`
	Stone uint64 `validator:"uint"`

	// 加速建筑效果
	BuildingCdr bool // true表示可以用于建筑减cd
	TechCdr     bool // true表示可以用于科技减cd
	WorkshopCdr bool // true表示可以用于装备锻造减cd
	// 减cd时间
	Cdr time.Duration

	// 迁移主城
	MoveBaseType shared_proto.GoodsMoveBaseType `head:"-"` // 迁城类型

	MoveBase      bool // true表示迁移主城
	RandomPos     bool // true表示随机位置
	GuildMoveBase bool // true表示联盟主城迁城令

	MonsterMoveSubLevel uint64 `validator:"uint" protofield:"-"` // 迁移行营随机等级

	// 经验丹
	ExpType shared_proto.GoodsExpEffectType `validator:"string" type:"enum"`
	Exp     uint64                          `validator:"uint"`

	// 免战
	MianDuration time.Duration

	// 行军加速百分比
	TroopSpeedUpRate float64 `validator:"float64"`

	// 修炼馆经验丹
	TrainDuration time.Duration

	// 合成物品
	PartsCombineCount   uint64                     `validator:"uint"`
	PartsPlunderPrizeId []uint64                   `validator:"uint" protofield:"-"`
	PartsShowPrize      []*shared_proto.PrizeProto `head:"-" protofield:",%s"`
	PartsShowType       uint64                     `validator:"uint" default:"0"`

	// 讨伐令
	AddMultiLevelNpcTimes bool // true表示添加讨伐次数

	// 攻城令
	AddInvaseHeroTimes bool // true表示添加攻城次数

	// 驻防镜像时间
	DurationType shared_proto.GoodsDurationEffectType `validator:"string" type:"enum"`
	Duration     time.Duration

	// 加数据类型效果
	AmountType shared_proto.GoodsAmountType `validator:"string" type:"enum"`
	Amount     uint64                       `validator:"uint"`

	// buff id
	BuffId uint64 `validator:"uint"`
}

func (e *GoodsEffect) Init(filename string) {
	e.initEffectType(filename)
}

func (e *GoodsEffect) initEffectType(filename string) {

	if e.Gold+e.Food+e.Wood+e.Stone > 0 {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_RESOURCE)
	}

	if e.BuildingCdr || e.TechCdr {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_CDR)

		check.PanicNotTrue(e.Cdr > 0, "%s 物品效果配置%v 配置的减cd时间必须>0", filename, e.Id)
	}

	if e.MoveBase {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_MOVE_BASE)

		if e.RandomPos {
			if e.GuildMoveBase {
				e.MoveBaseType = shared_proto.GoodsMoveBaseType_MOVE_BASE_GUILD
			} else {
				e.MoveBaseType = shared_proto.GoodsMoveBaseType_MOVE_BASE_RANDOM
			}
		} else {
			e.MoveBaseType = shared_proto.GoodsMoveBaseType_MOVE_BASE_POINT
		}
	}

	if e.Exp > 0 {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_EXP)

		_, exist := shared_proto.GoodsExpEffectType_name[int32(e.ExpType)]
		check.PanicNotTrue(exist && e.ExpType != shared_proto.GoodsExpEffectType_InvalidGoodsExp, "%s 物品效果配置%v 配置的exp>0，但是配置的经验类型ExpType无效, %v", filename, e.Id, e.ExpType)
	} else {
		check.PanicNotTrue(e.ExpType == shared_proto.GoodsExpEffectType_InvalidGoodsExp, "%s 物品%v 配置了经验丹类型，但是没有配置经验值", filename, e.Id, e.ExpType)
	}

	if e.MianDuration > 0 {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_MIAN)
	}

	if e.TroopSpeedUpRate > 0 {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_SPEED_UP)

		check.PanicNotTrue(e.TroopSpeedUpRate <= 1, "%s 物品效果配置%v 配置的行军加速百分比必须<=1", filename, e.Id)
	}

	if e.TrainDuration > 0 {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_TRAIN_EXP)
	}

	if e.PartsCombineCount > 0 {
		//check.PanicNotTrue(e.PartsCombineId > 0, "%s 物品效果配置%v 配置的合成物品零件，合成ID必须 > 0", filename, e.Id)
		check.PanicNotTrue(e.PartsCombineCount > 0, "%s 物品效果配置%v 配置的合成物品零件，合成所需零件个数必须 > 0", filename, e.Id)
		check.PanicNotTrue(len(e.PartsPlunderPrizeId) > 0, "%s 物品效果配置%v 配置的合成物品零件，没有配置合成掉落奖励", filename, e.Id)

		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_PARTS)
	}

	if e.AddMultiLevelNpcTimes {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_TFL)
	}

	if e.AddInvaseHeroTimes {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_GCL)
	}

	if e.Duration > 0 {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_DURATION)

		_, exist := shared_proto.GoodsDurationEffectType_name[int32(e.ExpType)]
		check.PanicNotTrue(exist && e.DurationType != shared_proto.GoodsDurationEffectType_InvalidGoodsDuration, "%s 物品效果配置%v 配置的duration>0，但是配置的经验类型DurationType无效, %v", filename, e.Id, e.DurationType)
	} else {
		check.PanicNotTrue(e.DurationType == shared_proto.GoodsDurationEffectType_InvalidGoodsDuration, "%s 物品%v 配置了时长类型，但是没有配置时长", filename, e.Id, e.DurationType)
	}

	if e.Amount > 0 {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_AMOUNT)
	}

	if e.BuffId > 0 && e.EffectType != shared_proto.GoodsEffectType_EFFECT_MIAN {
		e.panicSetEffectTypeIfNormal(filename, shared_proto.GoodsEffectType_EFFECT_BUFF)
	}
}

func (e *GoodsEffect) panicSetEffectTypeIfNormal(filename string, toSet shared_proto.GoodsEffectType) {
	check.PanicNotTrue(e.EffectType == shared_proto.GoodsEffectType_Normal, "%s 物品%v 配置的效果存在多种类型，除了是%v，还是%v", filename, e.Id, e.EffectType, toSet)
	e.EffectType = toSet
}

func (e *GoodsEffect) ResourceTypeCount() (count int) {
	if e.Gold > 0 {
		count++
	}
	if e.Food > 0 {
		count++
	}
	if e.Wood > 0 {
		count++
	}
	if e.Stone > 0 {
		count++
	}
	return
}
