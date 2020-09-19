package singleton

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/config/shop"
	"github.com/lightpaw/male7/util/sortkeys"
	"sort"
)

//gogen:config
type GoodsConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"物品/物品杂项.txt"`
	_ struct{} `proto:"shared_proto.GoodsConfigProto"`
	_ struct{} `protoconfig:"GoodsConfig"`

	EquipmentUpgradeGoods *goods.GoodsData   `protofield:",config.U64ToI32(%s.Id)"`
	EquipmentRefinedGoods *goods.GoodsData   `protofield:",config.U64ToI32(%s.Id)"`
	CaptainRefinedGoods   []*goods.GoodsData `protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`

	ChangeCaptainRaceGoods *goods.GoodsData `protofield:",config.U64ToI32(%s.Id)"`

	CaptainRebirthGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`

	MoveBaseGoods       *goods.GoodsData `head:"-" protofield:",config.U64ToI32(%s.Id)"`
	RandomMoveBaseGoods *goods.GoodsData `head:"-" protofield:",config.U64ToI32(%s.Id)"`
	GuildMoveBaseGoods  *goods.GoodsData `head:"-" protofield:",config.U64ToI32(%s.Id)"`

	MianGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`

	JiuGuanExpCaptainRefinedGoods []*goods.GoodsData `validator:"int>0,count=4" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"` // 酒馆请教获得的物品id
	ExpCaptainSoulUpgradeGoods    []*goods.GoodsData `validator:"int>0,count=4" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"` // 将魂升级的物品id

	SpeedUpGoods                 []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`
	TrainExpGoods                []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`
	TrainExpGoods4CaptainUpgrade []*goods.GoodsData `head:"-" protofield:"-"`

	BuildingCdrGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`
	TechCdrGoods     []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`
	WorkshopCdrGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`

	SpecGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`

	GoldGoods           []uint64 `head:"-"`
	GoldNormalShopGoods []uint64 `head:"-"`

	StoneGoods           []uint64 `head:"-"`
	StoneNormalShopGoods []uint64 `head:"-"`

	GongXunGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`

	MultiLevelNpcTimesGoods *goods.GoodsData `head:"-" protofield:",config.U64ToI32(%s.Id)"`

	InvaseHeroTimesGoods *goods.GoodsData `head:"-" protofield:",config.U64ToI32(%s.Id)"`

	JunTuanNpcTimesGoods *goods.GoodsData `head:"-" protofield:",config.U64ToI32(%s.Id)"`

	// 鱼饵
	FishGoods *goods.GoodsData `protofield:",config.U64ToI32(%s.Id)"`

	CopyDefenserGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`

	// 虎符
	TigerGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`
	// 过关斩将道具
	ZhanjiangGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`
	// 千重楼道具
	TowerGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`

	// buff 道具
	BuffGoods []*goods.GoodsData `head:"-" protofield:",config.U64a2I32a(goods.GetGoodsDataKeyArray(%s))"`

	// 等级最高的前n个宝物不可分解
	IndecomposableBaowuCount uint64
}

func (c *GoodsConfig) Init(filename string, configs interface {
	GetGoodsDataArray() []*goods.GoodsData
	GetShopArray() []*shop.Shop
}) {
	check.PanicNotTrue(len(c.CaptainRefinedGoods) > 0, "%s 没有配置武将强化符", filename)

	for _, g := range c.CaptainRefinedGoods {
		check.PanicNotTrue(g.GoodsEffect != nil && g.GoodsEffect.ExpType == shared_proto.GoodsExpEffectType_EXP_CAPTAIN_REFINED,
			"%s 配置的武将强化符物品，必须是经验物品，而且是武将强化经验物品, %d-%s", filename, g.Id, g.Name)
	}

	// 武将转生使用的物品就是武将强化物品
	c.CaptainRebirthGoods = c.CaptainRefinedGoods

	for _, g := range configs.GetGoodsDataArray() {
		if g.SpecType > 0 {
			c.SpecGoods = append(c.SpecGoods, g)
		}

		if g.GoodsEffect == nil {
			continue
		}

		switch g.GoodsEffect.MoveBaseType {
		case shared_proto.GoodsMoveBaseType_MOVE_BASE_POINT:
			c.MoveBaseGoods = panicOriginExist("高级迁城", c.MoveBaseGoods, g)
		case shared_proto.GoodsMoveBaseType_MOVE_BASE_RANDOM:
			c.RandomMoveBaseGoods = panicOriginExist("随机迁城", c.RandomMoveBaseGoods, g)
		case shared_proto.GoodsMoveBaseType_MOVE_BASE_GUILD:
			c.GuildMoveBaseGoods = panicOriginExist("联盟随机迁城", c.GuildMoveBaseGoods, g)
		}

		if g.GoodsEffect.MianDuration > 0 {
			c.MianGoods = append(c.MianGoods, g)
		}

		if g.GoodsEffect.TroopSpeedUpRate > 0 {
			c.SpeedUpGoods = append(c.SpeedUpGoods, g)
		}

		if g.GoodsEffect.ExpType == shared_proto.GoodsExpEffectType_EXP_CAPTAIN {
			c.TrainExpGoods = append(c.TrainExpGoods, g)
		}

		if g.EffectType == shared_proto.GoodsEffectType_EFFECT_BUFF || g.EffectType == shared_proto.GoodsEffectType_EFFECT_MIAN {
			c.BuffGoods = append(c.BuffGoods, g)
		}

		if g.GoodsEffect.Cdr > 0 {
			if g.GoodsEffect.BuildingCdr {
				c.BuildingCdrGoods = append(c.BuildingCdrGoods, g)
			}

			if g.GoodsEffect.TechCdr {
				c.TechCdrGoods = append(c.TechCdrGoods, g)
			}

			if g.GoodsEffect.WorkshopCdr {
				c.WorkshopCdrGoods = append(c.WorkshopCdrGoods, g)
			}
		}

		if g.GoodsEffect.Exp > 0 {
			switch g.GoodsEffect.ExpType {
			case shared_proto.GoodsExpEffectType_EXP_GONG_XUN:
				c.GongXunGoods = append(c.GongXunGoods, g)
			}
		}

		if g.GoodsEffect.AddMultiLevelNpcTimes {
			c.MultiLevelNpcTimesGoods = panicOriginExist("讨伐令", c.MultiLevelNpcTimesGoods, g)
		}

		if g.GoodsEffect.AddInvaseHeroTimes {
			c.InvaseHeroTimesGoods = panicOriginExist("攻城令", c.InvaseHeroTimesGoods, g)
		}

		if g.GoodsEffect.Duration > 0 {
			switch g.GoodsEffect.DurationType {
			case shared_proto.GoodsDurationEffectType_COPY_DEFENSER:
				c.CopyDefenserGoods = append(c.CopyDefenserGoods, g)
			}
		}

		if g.GoodsEffect.Amount > 0 {
			switch g.GoodsEffect.AmountType {
			case shared_proto.GoodsAmountType_AmountHufu:
				c.TigerGoods = append(c.TigerGoods, g)
			case shared_proto.GoodsAmountType_AmountZhanjiang:
				c.ZhanjiangGoods = append(c.ZhanjiangGoods, g)
			case shared_proto.GoodsAmountType_AmountTower:
				c.TowerGoods = append(c.TowerGoods, g)
			case shared_proto.GoodsAmountType_AmountJunTuan:
				c.JunTuanNpcTimesGoods = panicOriginExist("军团令", c.JunTuanNpcTimesGoods, g)
			}
		}
	}

	check.PanicNotTrue(c.RandomMoveBaseGoods != nil, "请配置随机迁城物品")
	check.PanicNotTrue(c.MoveBaseGoods != nil, "请配置高级迁城物品")
	check.PanicNotTrue(c.RandomMoveBaseGoods != nil, "请配置随机迁行营物品")

	check.PanicNotTrue(len(c.MianGoods) > 0, "请至少配置一个免战牌物品")
	check.PanicNotTrue(len(c.SpeedUpGoods) > 0, "请至少配置一个行军加速物品")
	check.PanicNotTrue(len(c.TrainExpGoods) > 0, "请至少配置一个修炼馆经验丹物品")

	check.PanicNotTrue(len(c.BuildingCdrGoods) > 0, "请至少配置一个减建筑CD物品")
	check.PanicNotTrue(len(c.TechCdrGoods) > 0, "请至少配置一个减科技CD物品")
	check.PanicNotTrue(len(c.WorkshopCdrGoods) > 0, "请至少配置一个减装备锻造CD物品")

	check.PanicNotTrue(len(c.GongXunGoods) > 0, "请至少配置一个功勋令牌物品")

	// 排序便于武将一键升级使用
	var kv_s []*sortkeys.U64KV
	for _, t := range c.TrainExpGoods {
		kv_s = append(kv_s, sortkeys.NewU64KV(t.GoodsEffect.Exp, t))
	}
	sort.Sort(sort.Reverse(sortkeys.U64KVSlice(kv_s)))
	for _, kv := range kv_s {
		c.TrainExpGoods4CaptainUpgrade = append(c.TrainExpGoods4CaptainUpgrade, kv.V.(*goods.GoodsData))
	}

	check.PanicNotTrue(c.MultiLevelNpcTimesGoods != nil, "请配置讨伐令物品")

	check.PanicNotTrue(c.InvaseHeroTimesGoods != nil, "请配置攻城令物品")

	check.PanicNotTrue(c.JunTuanNpcTimesGoods != nil, "请配置军团令物品")

	check.PanicNotTrue(len(c.CopyDefenserGoods) > 0, "请至少配置一个驻防镜像物品")
	check.PanicNotTrue(len(c.TigerGoods) > 0, "请至少配置一个虎符物品")
	check.PanicNotTrue(len(c.ZhanjiangGoods) > 0, "请至少配置一个过关斩将物品")
	check.PanicNotTrue(len(c.TowerGoods) > 0, "请至少配置一个千重楼物品")

	for _, data := range c.ExpCaptainSoulUpgradeGoods {
		check.PanicNotTrue(data.GoodsEffect != nil && data.GoodsEffect.ExpType == shared_proto.GoodsExpEffectType_EXP_CAPTAIN_SOUL_UPGRADE,
			"%s 将魂升级的物品id必须配置将魂升级加经验丹:%+v", filename, data)
	}

	for k := range shared_proto.GoodsSpecType_name {
		if k > 0 {
			t := shared_proto.GoodsSpecType(k)

			exist := false
			for _, g := range c.SpecGoods {
				if g.SpecType == t {
					exist = true
					break
				}
			}

			if !exist {
				logrus.Errorf("开服发现没有配置特殊物品 %v = %d", t, k)
			}
			//check.PanicNotTrue(exist, "请配置特殊物品 %v = %d", t, k)
		}
	}

	buildEffectResourceGoodsList(c, configs)
}

func panicOriginExist(goodsType string, origin, toSet *goods.GoodsData) *goods.GoodsData {
	if origin != nil {
		logrus.Panicf("%s道具配置了多个, goods1:%d-%s goods2:%d-%s", goodsType, origin.Id, origin.Name, toSet.Id, toSet.Name)
	}

	return toSet
}

// 构造资源物品-商品列表
func buildEffectResourceGoodsList(c *GoodsConfig, datas interface {
	GetGoodsDataArray() []*goods.GoodsData
	GetShopArray() []*shop.Shop
}) {
	// 防止重复写
	if len(c.GoldGoods) > 0 || len(c.GoldNormalShopGoods) > 0 ||
		len(c.StoneGoods) > 0 || len(c.StoneNormalShopGoods) > 0 {
		logrus.Debug("shop.buildEffectResourceGoodsList 构造资源物品-商品列表 重复调用了！")
		return
	}

	// 元宝商店里的普通资源物品
	for _, s := range datas.GetShopArray() {
		if s.Type != shop.YuanbaoShopType {
			continue
		}

	OUT:
		for _, g := range datas.GetGoodsDataArray() {
			if g.EffectType != shared_proto.GoodsEffectType_EFFECT_RESOURCE {
				continue
			}

			if g.GoodsEffect.ResourceTypeCount() != 1 {
				continue
			}

			for _, shopGoods := range s.NormalGoods {

				prize := shopGoods.PlunderPrize.GetPrize()
				if prize == nil {
					continue
				}

				if prize.TypeCount() != 1 {
					continue
				}

				if prize.Gold > 0 && prize.Gold == g.GoodsEffect.Gold && prize.UnsafeGold == 0 {
					c.GoldGoods = append(c.GoldGoods, g.Id)
					c.GoldNormalShopGoods = append(c.GoldNormalShopGoods, shopGoods.Id)
					continue OUT
				}
				if prize.Stone > 0 && prize.Stone == g.GoodsEffect.Stone && prize.UnsafeStone == 0 {
					c.StoneGoods = append(c.StoneGoods, g.Id)
					c.StoneNormalShopGoods = append(c.StoneNormalShopGoods, shopGoods.Id)
					continue OUT
				}
			}

			// 商店里买不到这个物品
			if g.GoodsEffect.Gold > 0 {
				c.GoldGoods = append(c.GoldGoods, g.Id)
				c.GoldNormalShopGoods = append(c.GoldNormalShopGoods, 0)
			} else if g.GoodsEffect.Stone > 0 {
				c.StoneGoods = append(c.StoneGoods, g.Id)
				c.StoneNormalShopGoods = append(c.StoneNormalShopGoods, 0)
			}
		}

	}

}
