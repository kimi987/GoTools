package heromodule

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/depot"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/equipment"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"time"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/sortkeys"
	"sort"
	"sync"
	"math/rand"
	"github.com/lightpaw/male7/service/operate_type"
	captain2 "github.com/lightpaw/male7/config/captain"
)

func AddUnsafeResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gold, food, wood, stone uint64) {
	addResource(hctx, hero, result, gold, food, wood, stone, false)
}

func AddSafeResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gold, food, wood, stone uint64) {
	addResource(hctx, hero, result, gold, food, wood, stone, true)
}

func addResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gold, food, wood, stone uint64, isSafeResource bool) {
	if gold+food+wood+stone <= 0 {
		return
	}

	heroResource := hero.GetUnsafeResource()
	if isSafeResource {
		heroResource = hero.GetSafeResource()
	}

	newGold, newFood, newWood, newStone := heroResource.AddResource(gold, food, wood, stone)
	result.AddFunc(func() pbutil.Buffer {
		return domestic.NewS2cResourceUpdateMsg(
			u64.Int32(newGold),
			u64.Int32(newFood),
			u64.Int32(newWood),
			u64.Int32(newStone),
			isSafeResource,
		)
	})

	result.Changed()

	tlogResource(hctx, hero, result, gold, food, wood, stone, true)
}

func ChangeResTypeToAmountType(t shared_proto.ResType) shared_proto.AmountType {
	switch t {
	case shared_proto.ResType_GOLD:
		return shared_proto.AmountType_PT_GOLD
	case shared_proto.ResType_STONE:
		return shared_proto.AmountType_PT_STONE
	case shared_proto.ResType_FOOD:
		return shared_proto.AmountType_PT_FOOD
	case shared_proto.ResType_WOOD:
		return shared_proto.AmountType_PT_WOOD
	}
	return shared_proto.AmountType_InvalidAmountType
}

func AddUnsafeSingleResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, resType shared_proto.ResType, toAdd uint64) (succ bool) {
	succ = addSingleResource(hero.GetUnsafeResource(), result, resType, toAdd)

	amountType := ChangeResTypeToAmountType(resType)

	// tlog
	hctx.Tlog().TlogMoneyFlow(hero, uint64(amountType), hero.GetUnsafeResource().GetRes(resType), toAdd, hctx.operateType.Id(), operate_type.Add)
	return
}

func AddSafeSingleResource(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, resType shared_proto.ResType, toAdd uint64) (succ bool) {
	succ = addSingleResource(hero.GetSafeResource(), result, resType, toAdd)

	amountType := ChangeResTypeToAmountType(resType)

	// tlog
	hctx.Tlog().TlogMoneyFlow(hero, uint64(amountType), hero.GetSafeResource().GetRes(resType), toAdd, hctx.operateType.Id(), operate_type.Add)

	return
}

func addSingleResource(heroResource *entity.ResourceStorage, result herolock.LockResult, resType shared_proto.ResType, toAdd uint64) bool {
	if toAdd <= 0 {
		return true
	}

	newRes := heroResource.AddRes(resType, toAdd)
	result.Add(domestic.NewS2cResourceUpdateSingleMsg(int32(resType), u64.Int32(newRes), heroResource.IsSafeResource()))

	result.Changed()
	return true
}

func AddPrize(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, prize *resdata.Prize, ctime time.Time) {
	addPrize(hctx, hero, result, prize, true, ctime)
}

func addPrize(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, prize *resdata.Prize, addCaptain bool, ctime time.Time) {

	// 资源
	gold, food, wood, stone := prize.GetSafeResource()
	addResource(hctx, hero, result, gold, food, wood, stone, true)

	gold, food, wood, stone = prize.GetUnsafeResource()
	addResource(hctx, hero, result, gold, food, wood, stone, false)

	// 玉璧
	AddJade(hctx, hero, result, prize.Jade)

	// 玉石矿
	AddJadeOre(hctx, hero, result, prize.JadeOre)

	// 经验
	AddExp(hctx, hero, result, prize.HeroExp, ctime)

	// 加元宝
	AddYuanbao(hctx, hero, result, prize.Yuanbao)

	// 加点券
	AddDianquan(hctx, hero, result, prize.Dianquan)

	// 加银两
	AddYinliang(hctx, hero, result, prize.Yinliang)

	// 加体力值
	AddSp(hctx, hero, result, prize.Sp)

	// 加联盟贡献币
	AddGuildContributionCoin(hctx, hero, result, prize.GuildContributionCoin)

	AddGoodsArray(hctx, hero, result, prize.Goods, prize.GoodsCount)
	AddEquipDataArray(hctx, hero, result, prize.Equipment, prize.EquipmentCount, ctime)
	AddGemArray(hctx, hero, result, prize.Gem, prize.GemCount, true)

	AddBaowuArray(hctx, hero, result, prize.Baowu, prize.BaowuCount, ctime)

	if addCaptain {
		AddResCaptainArray(hctx, hero, result, prize.Captain, prize.CaptainCount, ctime)
	}
}

func AddJade(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount <= 0 {
		return
	}

	newJade := hero.AddJade(amount)
	result.AddFunc(func() pbutil.Buffer {
		return domestic.NewS2cUpdateJadeMsg(u64.Int32(newJade), u64.Int32(hero.HistoryJade()), u64.Int32(hero.TodayObtainJade()))
	})

	hero.HistoryAmount().Increase(server_proto.HistoryAmountType_Jade, amount)
	UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_JADE)

	jade := hero.Jade()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_JADE), jade, amount, hctx.operateType.Id(), operate_type.Add)
}

func AddJadeOre(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount <= 0 {
		return
	}

	newJadeOre := hero.AddJadeOre(amount)
	result.AddFunc(func() pbutil.Buffer {
		return domestic.NewS2cUpdateJadeOreMsg(u64.Int32(newJadeOre))
	})

	jadeOre := hero.JadeOre()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_JADE_ORE), jadeOre, amount, hctx.operateType.Id(), operate_type.Add)
}

func AddYuanbao(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount <= 0 {
		return
	}

	hero.AddYuanbao(amount)
	result.Add(domestic.NewS2cUpdateYuanbaoMsg(u64.Int32(hero.GetYuanbao())))

	result.Changed()

	yuanbao := hero.GetYuanbao()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_YUANBAO), yuanbao, amount, hctx.operateType.Id(), operate_type.Add)
}

func AddYuanbaoGiftLimit(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount <= 0 {
		return
	}

	hero.Product().AddYuanbaoGiftLimit(amount)
	result.Add(domestic.NewS2cUpdateYuanbaoGiftLimitMsg(u64.Int32(hero.Product().GetYuanbaoGiftLimit())))

	result.Changed()
}

func ReduceYuanbaoGiftLimit(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount <= 0 {
		return
	}

	hero.Product().SubYuanbaoGiftLimit(amount)
	result.Add(domestic.NewS2cUpdateYuanbaoGiftLimitMsg(u64.Int32(hero.Product().GetYuanbaoGiftLimit())))

	result.Changed()
}

func AddDianquan(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount <= 0 {
		return
	}

	hero.AddDianquan(amount)
	result.Add(domestic.NewS2cUpdateDianquanMsg(u64.Int32(hero.GetDianquan())))

	result.Changed()

	dianquan := hero.GetDianquan()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_DIANQUAN), dianquan, amount, hctx.operateType.Id(), operate_type.Add)
}

func AddYinliang(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount <= 0 {
		return
	}
	hero.AddYinliang(amount)
	result.Add(domestic.NewS2cUpdateYinliangMsg(u64.Int32(hero.GetYinliang())))

	result.Changed()

	yinliang := hero.GetYinliang()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_YINLIANG), yinliang, amount, hctx.operateType.Id(), operate_type.Add)
}

func AddSp(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, amount uint64) {
	if amount <= 0 {
		return
	}
	hero.AddSp(amount)
	result.Add(domestic.NewS2cUpdateSpMsg(u64.Int32(hero.GetSp())))
	result.Changed()
}

// 加繁荣度
func AddProsperity(realmService iface.RealmService, heroId int64, baseRegion int64, toAdd uint64) {
	updateProsperity(realmService, heroId, baseRegion, true, toAdd)
}

func ReduceProsperity(realmService iface.RealmService, heroId int64, baseRegion int64, toAdd uint64) {
	updateProsperity(realmService, heroId, baseRegion, false, toAdd)
}

func updateProsperity(realmService iface.RealmService, heroId int64, baseRegion int64, isAdd bool, toChange uint64) {
	if toChange <= 0 || baseRegion == 0 {
		return
	}

	r := realmService.GetBigMap()
	if r != nil {
		if isAdd {
			r.AddProsperity(heroId, toChange)
		} else {
			r.ReduceProsperity(heroId, toChange)
		}
	}
}

func AddGoodsArray(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsArray []*goods.GoodsData, goodsCount []uint64) {
	// 物品
	n := imath.Min(len(goodsArray), len(goodsCount))

	if n <= 0 {
		return
	}

	heroDepot := hero.Depot()

	oldCount := make([]uint64, n)
	newCount := make([]uint64, n)
	for i := 0; i < n; i++ {
		count := goodsCount[i]
		if count > 0 {
			// 合璧物品轮流给
			g, ok := getPartnerGoodsIfHebiType(hero.GuildId(), goodsArray[i])
			if ok && g != nil {
				goodsArray[i] = g
			}

			oldCount[i] = heroDepot.GetGoodsCount(goodsArray[i].Id)
			newCount[i] = heroDepot.AddGoods(goodsArray[i].Id, count)

			onAddGoodsEvent(hero, result, goodsArray[i].Id, count)
		}
	}

	result.Add(depot.NewS2cUpdateMultiGoodsMsg(u64.Int32Array(goods.GetGoodsDataKeyArray(goodsArray)), u64.Int32Array(newCount)))

	result.Changed()

	// tlog
	for i := 0; i < n; i++ {
		gd, count, oldCnt, newCnt := goodsArray[i], goodsCount[i], oldCount[i], newCount[i]
		if count <= 0 {
			continue
		}
		hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeConsumable, uint64(gd.EffectType)), gd.Id, oldCnt, newCnt, hctx.operateType.Id(), operate_type.Add, gd.GoodsQuality.Level, gd.Id, int64(count))
	}
}

var hebiGoodsCounters = sync.Map{}

// 合璧物品要左右轮流给
func getPartnerGoodsIfHebiType(guildId int64, goodsData *goods.GoodsData) (g *goods.GoodsData, ok bool) {
	g = goodsData
	if goodsData.SpecType != shared_proto.GoodsSpecType_GAT_HEBI {
		return
	}

	if g.HebiType != shared_proto.HebiType_HeShiBi {
		guildId = 0
	}

	counter, ok := hebiGoodsCounters.Load(guildId)
	if !ok {
		counter = atomic.NewUint64(0)
		hebiGoodsCounters.Store(guildId, counter)
	}

	var i uint64
	if c, ok := counter.(*atomic.Uint64); ok {
		i = c.Inc()
	} else {
		i = rand.Uint64()
	}

	if i%2 == 0 {
		g = goodsData.PartnerHebiGoodsData
		ok = true
		return
	}

	return
}

func AddGoodsIdArray(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsIdArray, goodsCount []uint64) {
	// 物品
	n := imath.Min(len(goodsIdArray), len(goodsCount))

	if n <= 0 {
		return
	}

	heroDepot := hero.Depot()

	oldCount := make([]uint64, n)
	newCount := make([]uint64, n)
	for i := 0; i < n; i++ {
		count := goodsCount[i]
		if count > 0 {
			goodsId := goodsIdArray[i]

			oldCount[i] = heroDepot.GetGoodsCount(goodsId)
			newCount[i] = heroDepot.AddGoods(goodsId, count)

			onAddGoodsEvent(hero, result, goodsId, count)
		}
	}

	result.Add(depot.NewS2cUpdateMultiGoodsMsg(u64.Int32Array(goodsIdArray), u64.Int32Array(newCount)))

	result.Changed()

	// tlog
	for i := 0; i < n; i++ {
		goodsId, count, oldCnt, newCnt := goodsIdArray[i], goodsCount[i], oldCount[i], newCount[i]
		if count <= 0 {
			continue
		}
		if g := hctx.Datas().GetGoodsData(goodsId); g != nil {
			hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeConsumable, uint64(g.EffectType)), g.Id, oldCnt, newCnt, hctx.operateType.Id(), operate_type.Add, g.GoodsQuality.Level, g.Id, int64(count))
		} else if g := hctx.Datas().GetGemData(goodsId); g != nil {
			hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeGem, g.GemType), g.Id, oldCnt, newCnt, hctx.operateType.Id(), operate_type.Add, uint64(g.Quality), g.Id, int64(count))
		}
	}
}

func AddGoods(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, g *goods.GoodsData, count uint64) {
	if count <= 0 {
		return
	}

	oldCount := hero.Depot().GetGoodsCount(g.Id)
	newCount := hero.Depot().AddGoods(g.Id, count)

	onAddGoodsEvent(hero, result, g.Id, count)

	result.Add(depot.NewS2cUpdateGoodsMsg(u64.Int32(g.Id), u64.Int32(newCount)))

	result.Changed()

	// tlog
	hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeConsumable, uint64(g.EffectType)), g.Id, oldCount, newCount, hctx.operateType.Id(), operate_type.Add, g.GoodsQuality.Level, g.Id, int64(count))
}

func AddEquipDataArray(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, equipArray []*goods.EquipmentData, equipCount []uint64, ctime time.Time) {
	// 物品
	n := imath.Min(len(equipArray), len(equipCount))

	if n <= 0 {
		return
	}

	for i := 0; i < n; i++ {
		AddEquipData(hctx, hero, result, equipArray[i], equipCount[i], ctime)
	}
}

func AddEquipData(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, equipData *goods.EquipmentData, equipCount uint64, ctime time.Time) {
	heroDepot := hero.Depot()
	for j := uint64(0); j < equipCount; j++ {
		e := entity.NewEquipment(heroDepot.NewId(), equipData)
		AddEquip(hctx, hero, result, e, ctime)
	}
}

func AddEquipArray(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, equipArray []*entity.Equipment, ctime time.Time) {
	if len(equipArray) <= 0 {
		return
	}

	for _, e := range equipArray {
		AddEquip(hctx, hero, result, e, ctime)
	}
}

func AddEquip(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, equip *entity.Equipment, ctime time.Time) {
	expireTime := hero.Depot().AddGenIdGoods(equip, ctime)
	if expireTime != 0 {
		result.Add(equipment.NewS2cAddEquipmentWithExpireTimeMarshalMsg(equip.EncodeClient(), i64.Int32(expireTime)))
	} else {
		result.Add(equipment.NewS2cAddEquipmentMarshalMsg(equip.EncodeClient()))
	}
	result.Changed()

	addEquipBroadcast(hctx, hero, equip, result)

	// tlog
	hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeEquip, uint64(equip.Data().Type)), equip.Data().Id, 0, 1, hctx.operateType.Id(), operate_type.Add, equip.Data().Quality.GoodsQuality.Level, equip.Data().Id, 1)
}

func AddGemArray(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gemArray []*goods.GemData, gemCount []uint64, isNewGem bool) {
	// 物品
	n := imath.Min(len(gemArray), len(gemCount))

	if n <= 0 {
		return
	}

	heroDepot := hero.Depot()

	oldCount := make([]uint64, n)
	newCount := make([]uint64, n)
	for i := 0; i < n; i++ {
		count := gemCount[i]
		if count > 0 {
			gemData := gemArray[i]

			oldCount[i] = heroDepot.GetGoodsCount(gemData.Id)
			newCount[i] = heroDepot.AddGoods(gemData.Id, count)

			onAddGoodsEvent(hero, result, gemData.Id, count)
		}
	}

	result.Add(depot.NewS2cUpdateMultiGoodsMsg(u64.Int32Array(goods.GetGemDataKeyArray(gemArray)), u64.Int32Array(newCount)))

	result.Changed()

	if isNewGem {
		// 获得了新宝石，更新任务
		UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_GEM)
	}

	// tlog
	for i := 0; i < n; i++ {
		gemData, count, oldCnt, newCnt := gemArray[i], gemCount[i], oldCount[i], newCount[i]
		if count <= 0 {
			continue
		}
		hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeGem, gemData.GemType), gemData.Id, oldCnt, newCnt, hctx.operateType.Id(), operate_type.Add, uint64(gemData.Quality), gemData.Id, int64(count))
	}
}

// 加一个宝石数组，每一个都加1
func AddGemArrayGive1(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gemArray []*goods.GemData, isNewGem bool) {
	if len(gemArray) <= 0 {
		return
	}

	// 物品
	heroDepot := hero.Depot()

	oldCount := make([]uint64, len(gemArray))
	newCount := make([]uint64, len(gemArray))
	for i := 0; i < len(gemArray); i++ {
		gemData := gemArray[i]

		oldCount[i] = heroDepot.GetGoodsCount(gemData.Id)
		newCount[i] = heroDepot.AddGoods(gemData.Id, 1)

		onAddGoodsEvent(hero, result, gemData.Id, 1)
	}

	result.Add(depot.NewS2cUpdateMultiGoodsMsg(u64.Int32Array(goods.GetGemDataKeyArray(gemArray)), u64.Int32Array(newCount)))

	if isNewGem {
		// 获得了新宝石，更新任务
		UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_GEM)
	}

	result.Changed()

	// tlog
	for i := 0; i < len(gemArray); i++ {
		gemData, oldCnt, newCnt := gemArray[i], oldCount[i], newCount[i]
		hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeGem, gemData.GemType), gemData.Id, oldCnt, newCnt, hctx.operateType.Id(), operate_type.Add, uint64(gemData.Quality), gemData.Id, 1)
	}
}

func AddGem(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, gemData *goods.GemData, gemCount uint64, isNewGem bool) {
	if gemCount <= 0 {
		return
	}

	oldCount := hero.Depot().GetGoodsCount(gemData.Id)
	newCount := hero.Depot().AddGoods(gemData.Id, gemCount)

	onAddGoodsEvent(hero, result, gemData.Id, gemCount)

	result.Add(depot.NewS2cUpdateGoodsMsg(u64.Int32(gemData.Id), u64.Int32(newCount)))

	if isNewGem {
		// 获得了新宝石，更新任务
		UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_GEM)
	}

	// 系统广播
	AddGemBroadcast(hctx, hero, gemData, result)

	result.Changed()

	// tlog
	hctx.Tlog().TlogItemFlow(hero, operate_type.BuildGoodsType(operate_type.GoodsTypeGem, gemData.GemType), gemData.Id, oldCount, newCount, hctx.operateType.Id(), operate_type.Add, uint64(gemData.Quality), gemData.Id, int64(gemCount))
}

func AddResCaptainArray(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, captainArray []*resdata.ResCaptainData, captainCount []uint64, ctime time.Time) {
	// 物品
	n := imath.Min(len(captainArray), len(captainCount))

	if n <= 0 {
		return
	}

	for i := 0; i < n; i++ {
		AddResCaptain(hctx, hero, result, captainArray[i], captainCount[i], ctime)
	}
}

func AddResCaptain(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, resCaptainData *resdata.ResCaptainData,
	captainCount uint64, ctime time.Time) {
	if captainCount <= 0 {
		return
	}

	captainData := captain2.GetCaptainData(resCaptainData)
	if captainData == nil {
		logrus.Errorf("AddResCaptain captainData == nil, id: %v", resCaptainData.Id)
		return
	}

	if TryAddCaptain(hero, result, captainData, ctime) {
		// 没解锁，给一个
		captainCount--

		UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_SOUL_COUNT)
		UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_SOUL_QUALITY_COUNT)
		UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_SOUL_LEVEL_COUNT)

		addCaptainBroadcast(hctx, hero, captainData, result)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_GET_CAPTAIN_SOUL) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_GET_CAPTAIN_SOUL)))
		}

		if hctx.OperType() != operate_type.FishingFishing {
			// 排除钓鱼，钓鱼逻辑有特殊处理
			result.Add(military.NewS2cShowPrizeCaptainMsg(u64.Int32(captainData.Id), false))
		}
	}

	if captainCount > 0 {
		var prize *resdata.Prize

		if captainCount == 1 {
			prize = captainData.PrizeIfHas
		} else {
			prize = resdata.NewPrizeBuilder().AddMultiple(captainData.PrizeIfHas, captainCount).Build()
		}

		// 不要再触发加将魂了，防止死循环
		addPrize(hctx, hero, result, prize, false, ctime)

		if hctx.OperType() != operate_type.FishingFishing {
			// 排除钓鱼，钓鱼逻辑有特殊处理
			result.Add(military.NewS2cShowPrizeCaptainMsg(u64.Int32(captainData.Id), true))
		}
	}

	result.Changed()
}

func AddPveTroopPrize(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, prize *resdata.Prize, pveTroopType shared_proto.PveTroopType, ctime time.Time) {
	AddPrize(hctx, hero, result, prize, ctime)

	pveTroop := hero.PveTroop(pveTroopType)
	if pveTroop == nil {
		logrus.Errorf("AddPveTroopPrize 时，未找到类型: %v 的队伍", pveTroopType)
		return
	}

	for _, pos := range pveTroop.Captains() {
		c := pos.Captain()
		if c != nil {
			AddCaptainExp(hctx, hero, result, c, prize.CaptainExp, ctime)
		}
	}
}

func AddCaptainPrize(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, prize *resdata.Prize, captainIds []uint64, ctime time.Time) {

	AddPrize(hctx, hero, result, prize, ctime)

	if prize.CaptainExp > 0 && len(captainIds) > 0 {
		// 武将经验部分，用map是防止传入的captainIds存在重复的id，导致后面同一个武将加多次经验
		// 为什么不提前检查一下captainIds呢，因为大部分情况都是没有问题的
		captains := make(map[uint64]*entity.Captain, len(captainIds))
		for _, cid := range captainIds {
			c := hero.Military().Captain(cid)
			if c == nil {
				continue
			}

			captains[cid] = c
		}

		if len(captains) > 0 {
			for _, c := range captains {
				AddCaptainExp(hctx, hero, result, c, prize.CaptainExp, ctime)
			}
		}
	}
}

func AddCaptainExpById(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, captainId, toAddExp uint64, ctime time.Time) bool {

	if toAddExp <= 0 {
		return false
	}

	captain := hero.Military().Captain(captainId)
	if captain == nil {
		logrus.Debugf("添加武将经验，武将没找到")
		return false
	}

	return AddCaptainExp(hctx, hero, result, captain, toAddExp, ctime)
}

func AddCaptainExp(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, captain *entity.Captain, toAddExp uint64, ctime time.Time) bool {

	result.Changed()

	// 转生时不加经验。目前需要通知客户端的地方只有吃经验丹时
	if captain.IsInRebrithing(ctime) {
		return false
	}

	// 增加经验
	if !captain.AddExp(toAddExp, hero.LevelData().Sub.CaptainLevelLimit) {
		result.Add(military.NewS2cUpdateCaptainExpMsg(u64.Int32(captain.Id()), u64.Int32(captain.Exp())))
		return false
	}

	//if captain.IsOutSide() {
	//	result.Add(military.NewS2cUpdateCaptainExpMsg(u64.Int32(captain.Id()), u64.Int32(captain.Exp())))
	//	return false
	//}

	// 处理升级
	doUpgradeCaptain(hctx, hero, result, captain, ctime)
	return true
}

func doUpgradeCaptain(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, captain *entity.Captain, ctime time.Time) {
	oldLevel := captain.Level()
	oldAbility := captain.Ability()
	oldRace := captain.Race().Id
	oldOfficial := captain.Official().Id

	// 处理升级
	originLevel, newLevel, startRebirthCD := captain.UpgradeLevel(hero.LevelData().Sub.CaptainLevelLimit, ctime)

	if originLevel < newLevel {
		// sync
		captain.CalculateProperties()

		result.Add(military.NewS2cCaptainLevelupMsg(u64.Int32(captain.Id()), u64.Int32(captain.Exp()), u64.Int32(captain.Level()), u64.Int32(captain.SoldierCapcity())))
		result.Add(captain.NewUpdateCaptainStatMsg())
		UpdateTroopFightAmount(hero, captain.GetTroop(), result)

		// tlog
		hctx.Tlog().TlogPlayerCultivateFlow(hero, captain.Id(), operate_type.CaptainOperTypeUpgrade, oldLevel, captain.Level(), oldAbility, captain.Ability(), u64.FromInt(oldRace), u64.FromInt(captain.Race().Id), oldOfficial, captain.Official().Id, hctx.OperId())

	} else {
		// 没升级就更新经验
		result.Add(military.NewS2cUpdateCaptainExpMsg(u64.Int32(captain.Id()), u64.Int32(captain.Exp())))
	}

	if startRebirthCD {
		result.Add(military.NewS2cCaptainRebirthCdStartMsg(u64.Int32(captain.Id()), timeutil.Marshal32(captain.RebirthEndTime())))
	}
	//if originRebirthExp < newRebirthExp {
	//	result.Add(military.NewS2cUpdateCaptainExpMsg(u64.Int32(captain.Id()), u64.Int32(captain.Exp())))
	//}

	UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_CAPTAIN_LEVEL_COUNT)
}

func TryUpgradeCaptain(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, captainId uint64, ctime time.Time) {
	captain := hero.Military().Captain(captainId)
	if captain != nil && captain.IsUpgrade(hero.LevelData().Sub.CaptainLevelLimit) {
		doUpgradeCaptain(hctx, hero, result, captain, ctime)
	}
}

func AddGuildContributionCoin(hctx *HeroContext, hero *entity.Hero, result herolock.LockResult, toAdd uint64) {
	if toAdd <= 0 {
		return
	}

	hero.AddGuildContributionCoin(toAdd)
	result.Add(guild.NewS2cUpdateContributionCoinMsg(u64.Int32(hero.GetGuildContributionCoin())))

	hero.HistoryAmount().Increase(server_proto.HistoryAmountType_GuildContributionCoin, toAdd)
	UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_GUILD_CONTRIBUTION)

	result.Changed()

	gcCoin := hero.GetGuildContributionCoin()
	hctx.Tlog().TlogMoneyFlow(hero, uint64(shared_proto.AmountType_PT_GUILD_CONTRIBUTION_COIN), gcCoin, toAdd, hctx.operateType.Id(), operate_type.Add)
}

func addEquipBroadcast(hctx *HeroContext, hero *entity.Hero, equip *entity.Equipment, result herolock.LockResult) {
	var d *data.BroadcastData
	switch hctx.operateType {
	case operate_type.FishingFishing:
		d = hctx.datas.BroadcastHelp().FishEquip
	case operate_type.TowerChallenge:
		fallthrough
	case operate_type.TowerAutoChallenge:
		d = hctx.datas.BroadcastHelp().TowerPrizeEquip
	case operate_type.DepotGoodsCombine, operate_type.DepotGoodsPartCombine:
		d = hctx.datas.BroadcastHelp().MixEquip
	default:
		return
	}

	if d != nil {
		equipData := equip.Data()
		hctx.AddBroadcast(d, hero, result, equipData.Id, uint64(equipData.Quality.GoodsQuality.Quality), func() *i18n.Fields {
			text := d.NewTextFields()
			text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
			text.WithClickEquipFields(data.KeyEquip, hctx.broadcast.GetEquipText(equipData), equipData.Id, equip.Level(), equip.RefinedLevel())
			return text
		})
	}
}

func addCaptainBroadcast(hctx *HeroContext, hero *entity.Hero, csData *captain2.CaptainData, result herolock.LockResult) {
	var d *data.BroadcastData
	switch hctx.operateType {
	case operate_type.FishingFishing:
		d = hctx.datas.BroadcastHelp().FishCaptainSoul
	case operate_type.DungeonChallenge:
		fallthrough
	case operate_type.DungeonAutoChallenge, operate_type.DungeonCollectChapterPrize, operate_type.DungeonCollectPassDungeonPrize:
		d = hctx.datas.BroadcastHelp().DungeonPrizeCaptainSoul
	default:
		return
	}

	if d != nil {
		hctx.AddBroadcast(d, hero, result, 0, csData.Rarity.Id, func() *i18n.Fields {
			text := d.NewTextFields().WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id()).WithFields(data.KeyText, csData.Rarity.Name)
			nameColorData := hctx.datas.ColorData().Must(uint64(shared_proto.Quality_ORANGE))
			text.WithClickCaptainFields(data.KeyCaptain, nameColorData.GetColorText(csData.Name), hero.Id(), csData.Id)
			return text
		})
	}
}

func AddGemBroadcast(hctx *HeroContext, hero *entity.Hero, gemData *goods.GemData, result herolock.LockResult) {
	var d *data.BroadcastData
	switch hctx.operateType {
	case operate_type.GemCombineGem:
		fallthrough
	case operate_type.GemOneKeyCombineGem:
		fallthrough
	case operate_type.GemOneKeyCombineDepotGem:
		d = hctx.datas.BroadcastHelp().GemGet
	case operate_type.FishingFishing:
		d = hctx.datas.BroadcastHelp().FishGem
	default:
		return
	}

	if d != nil {
		hctx.AddBroadcast(d, hero, result, 0, gemData.Level, func() *i18n.Fields {
			text := d.NewTextFields()
			text.WithClickHeroFields(data.KeySelf, hctx.GetFlagHeroName(hero), hero.Id())
			text.WithFields(data.KeyNum, gemData.Level)
			return text
		})
	}
}

func AddBaowuArray(ctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsArray []*resdata.BaowuData, goodsCount []uint64, ctime time.Time) {
	// 宝物
	n := imath.Min(len(goodsArray), len(goodsCount))

	if n <= 0 {
		return
	}

	heroDepot := hero.Depot()
	for i := 0; i < n; i++ {
		count := goodsCount[i]
		if count > 0 {
			if data := goodsArray[i]; data != nil {
				heroDepot.AddBaowu(data.Id, count)

				// 添加日志
				AddBaowuLog(ctx, hero, result, data, count, ctime)
			}
		}
	}

	tryAutoUpgradeBaowu(hero, result, ctime, goodsArray...)

	result.Changed()
}

func AddBaowu(ctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *resdata.BaowuData, goodsCount uint64, ctime time.Time) {
	if goodsCount <= 0 {
		return
	}

	hero.Depot().AddBaowu(goodsData.Id, goodsCount)

	// 添加日志
	AddBaowuLog(ctx, hero, result, goodsData, goodsCount, ctime)

	tryAutoUpgradeBaowu(hero, result, ctime, goodsData)

	result.Changed()
}

func AddBaowuLog(ctx *HeroContext, hero *entity.Hero, result herolock.LockResult, goodsData *resdata.BaowuData, goodsCount uint64, ctime time.Time) {

	// 添加日志
	log := &shared_proto.BaowuLogProto{}
	log.Op = ctx.GetBaowuOp()
	log.Baowu = u64.Int32(goodsData.Id)
	log.Count = u64.Int32(goodsCount)
	log.Time = timeutil.Marshal32(ctime)

	if log.Op == shared_proto.BaowuOpType_BOTInvalid {
		log.Op = shared_proto.BaowuOpType_BOTAdd
	} else {
		log.OtherName, log.OtherX, log.OtherY = ctx.BaowuOtherInfo()
		log.IsRare = ctx.IsRareBaowu(goodsData.Id)
	}

	logBytes := must.Marshal(log)
	result.Add(depot.NewS2cAddBaowuLogMsg(logBytes))

	hero.Depot().AddBaowuLog(logBytes)
}

func tryAutoUpgradeBaowu(hero *entity.Hero, result herolock.LockResult, ctime time.Time, datas ...*resdata.BaowuData) {

	// 根据id进行排序
	n := len(datas)
	sortDatas := make([]*sortkeys.U64KV, n)
	for i, v := range datas {
		sortDatas[i] = &sortkeys.U64KV{
			K: v.Id,
			V: v,
		}
	}
	if n > 1 {
		sort.Sort(sortkeys.U64KVSlice(sortDatas))
	}

	proto := &depot.S2CUpdateMultiBaowuProto{}
	for i := 0; i < n; i++ {
		curData := sortDatas[i].V.(*resdata.BaowuData)
		curCount := hero.Depot().GetBaowuCount(curData.Id)

		var nextData *resdata.BaowuData
		if i+1 < n {
			nextData = datas[i+1]
		}

		for {
			if curData == nextData {
				break
			}

			if curData.UpgradeNeedCount <= 0 || curCount < curData.UpgradeNeedCount {
				// 添加消息
				proto.Id = append(proto.Id, u64.Int32(curData.Id))
				proto.Count = append(proto.Count, u64.Int32(curCount))
				break
			}

			nextLevel := curData.GetNextLevel()
			if nextLevel == nil || curData.Level >= nextLevel.Level {
				proto.Id = append(proto.Id, u64.Int32(curData.Id))
				proto.Count = append(proto.Count, u64.Int32(curCount))
				break
			}

			// 自动合成

			upgradeCount := curCount / curData.UpgradeNeedCount
			costCount := upgradeCount * curData.UpgradeNeedCount

			// 扣掉低级的个数
			curCount = hero.Depot().RemoveBaowu(curData.Id, costCount)
			proto.Id = append(proto.Id, u64.Int32(curData.Id))
			proto.Count = append(proto.Count, u64.Int32(curCount))

			nextlevelCount := hero.Depot().AddBaowu(nextLevel.Id, upgradeCount)

			// 合成日志
			log := &shared_proto.BaowuLogProto{}
			log.Op = shared_proto.BaowuOpType_BOTAuto
			log.Baowu = u64.Int32(curData.Id)
			log.Count = u64.Int32(upgradeCount)
			log.Time = timeutil.Marshal32(ctime)

			logBytes := must.Marshal(log)
			result.Add(depot.NewS2cAddBaowuLogMsg(logBytes))

			hero.Depot().AddBaowuLog(logBytes)

			// 进入下一个循环
			curData = nextLevel
			curCount = nextlevelCount

		}
	}

	result.Add(depot.NewS2cUpdateMultiBaowuProtoMsg(proto))
}
