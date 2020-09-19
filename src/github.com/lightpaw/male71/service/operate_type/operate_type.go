package operate_type

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/gen/pb/bai_zhan"
	"github.com/lightpaw/male7/gen/pb/captain_soul"
	"github.com/lightpaw/male7/gen/pb/chat"
	"github.com/lightpaw/male7/gen/pb/country"
	"github.com/lightpaw/male7/gen/pb/depot"
	"github.com/lightpaw/male7/gen/pb/dianquan"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/dungeon"
	"github.com/lightpaw/male7/gen/pb/equipment"
	"github.com/lightpaw/male7/gen/pb/farm"
	"github.com/lightpaw/male7/gen/pb/fishing"
	"github.com/lightpaw/male7/gen/pb/garden"
	"github.com/lightpaw/male7/gen/pb/gem"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/guizu"
	"github.com/lightpaw/male7/gen/pb/hebi"
	"github.com/lightpaw/male7/gen/pb/mail"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/gen/pb/mingc_war"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/promotion"
	"github.com/lightpaw/male7/gen/pb/question"
	"github.com/lightpaw/male7/gen/pb/random_event"
	"github.com/lightpaw/male7/gen/pb/red_packet"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/gen/pb/secret_tower"
	"github.com/lightpaw/male7/gen/pb/shop"
	"github.com/lightpaw/male7/gen/pb/strategy"
	"github.com/lightpaw/male7/gen/pb/task"
	"github.com/lightpaw/male7/gen/pb/teach"
	"github.com/lightpaw/male7/gen/pb/tower"
	"github.com/lightpaw/male7/gen/pb/vip"
	"github.com/lightpaw/male7/gen/pb/xiongnu"
	"github.com/lightpaw/male7/gen/pb/xuanyuan"
	"github.com/lightpaw/male7/gen/pb/zhanjiang"
	"github.com/lightpaw/male7/gen/pb/zhengwu"
	"github.com/lightpaw/male7/gen/pb/activity"
)

var (
	checkMap = make(map[uint64]struct{})

	TlogIgnore = newOperateType(0, 0, "", checkMap)

	GM    = newOperateType(1000, 1, "GM 接口", checkMap)
	GMCmd = newOperateType(1000, 2, "GM 命令", checkMap)

	HeroAddBackReservation = newOperateType(1001, 1, "玩家-预付退回", checkMap)
	HeroExchange           = newOperateType(1001, 2, "玩家-交换", checkMap)
	HeroUpgradeBaseLevel   = newOperateType(1001, 3, "玩家-升级主城等级", checkMap)
	HeroUpgradeLevel       = newOperateType(1001, 4, "玩家-升级君主等级", checkMap)
	HeroUpdatePerSecond    = newOperateType(1001, 5, "玩家-每秒更新", checkMap)
	HeroLoginLoad          = newOperateType(1001, 6, "玩家-登录加载", checkMap)

	RealmRobNpc         = newOperateType(1002, 1, "野外-抢夺 NPC", checkMap)
	RealmInvasionNpc    = newOperateType(1002, 2, "野外-入侵 NPC", checkMap)
	RealmKillInvadeNpc  = newOperateType(1002, 3, "野外-击杀出征Npc", checkMap)
	RealmKillDefNpc     = newOperateType(1002, 7, "野外-击杀防守Npc", checkMap)
	RealmInvasion       = newOperateType(1002, 4, "野外-入侵玩家", checkMap)
	RealmRobHeroHome    = newOperateType(1002, 5, "野外-抢夺玩家", checkMap)
	RealmReturnedToBase = newOperateType(1002, 6, "野外-回家", checkMap)

	BaoWuBeenRobbed = newOperateType(1004, 1, "宝物-被抢", checkMap)
	BaoWuRobHero    = newOperateType(1004, 2, "宝物-抢夺玩家", checkMap)
	BaoWuRobNpc     = newOperateType(1004, 3, "宝物-抢夺 npc", checkMap)
	BaoWuDecompose  = newOperateType(1004, 4, "宝物-分解", checkMap)

	ActivityCollect = newOperateType(activity.MODULE_ID, 1, "活动-收集", checkMap)

	BaiZhanCollectSalary       = newOperateType(bai_zhan.MODULE_ID, 1, "百战千军-领俸禄", checkMap)
	BaiZhanCollectJunXianPrize = newOperateType(bai_zhan.MODULE_ID, 2, "百战千军-领军衔奖励", checkMap)

	CaptainSoulCollectFettersPrize = newOperateType(captain_soul.MODULE_ID, 1, "名士-羁绊奖励", checkMap)
	CaptainSoulInherit             = newOperateType(captain_soul.MODULE_ID, 2, "名士-继承", checkMap)
	CaptainSoulUpgrade             = newOperateType(captain_soul.MODULE_ID, 3, "名士-升级", checkMap)
	CaptainSoulUpgradeV2           = newOperateType(captain_soul.MODULE_ID, 4, "名士-升级v2", checkMap)
	CaptainSoulUnlockSpell         = newOperateType(captain_soul.MODULE_ID, 5, "名士-解锁技能", checkMap)
	CaptainSoulFuShen              = newOperateType(captain_soul.MODULE_ID, 6, "名士-附身", checkMap)
	CaptainSoulReborn              = newOperateType(captain_soul.MODULE_ID, 7, "名士-重生", checkMap)

	ChatLaba = newOperateType(chat.MODULE_ID, 1, "聊天-喇叭广播", checkMap)

	CountryChangeCountry         = newOperateType(country.MODULE_ID, 1, "国家-更改国家", checkMap)
	CountryCollectOfficialSalary = newOperateType(country.MODULE_ID, 2, "国家-领取官职俸禄", checkMap)
	CountryChangeName            = newOperateType(country.MODULE_ID, 3, "国家-改名", checkMap)

	DepotUseGoods           = newOperateType(depot.MODULE_ID, 1, "背包-使用物品", checkMap)
	DepotGoodsCombine       = newOperateType(depot.MODULE_ID, 2, "背包-合成", checkMap)
	DepotGoodsPartCombine   = newOperateType(depot.MODULE_ID, 3, "背包-碎片合成", checkMap)
	DepotRemoveExpiredGoods = newOperateType(depot.MODULE_ID, 4, "背包-删除过期物品", checkMap)
	DepotUseCdrGoods        = newOperateType(depot.MODULE_ID, 5, "背包-使用 CD 物品", checkMap)
	DepotCollectBaowu       = newOperateType(depot.MODULE_ID, 6, "背包-领取宝物", checkMap)

	DianquanExchange = newOperateType(dianquan.MODULE_ID, 1, "点券兑换", checkMap)

	DomesticCreateResourceBuilding   = newOperateType(domestic.MODULE_ID, 1, "内政-解锁资源建筑", checkMap)
	DomesticUpgradeResourceBuilding  = newOperateType(domestic.MODULE_ID, 2, "内政-升级资源建筑", checkMap)
	DomesticRebuildBuilding          = newOperateType(domestic.MODULE_ID, 3, "内政-改建建筑", checkMap)
	DomesticCollectResource          = newOperateType(domestic.MODULE_ID, 4, "内政-收集资源", checkMap)
	DomesticCollectResourceV2        = newOperateType(domestic.MODULE_ID, 5, "内政-收集资源v2", checkMap)
	DomesticRequestResourceConflict  = newOperateType(domestic.MODULE_ID, 6, "", checkMap)
	DomesticUpgradeStableBuilding    = newOperateType(domestic.MODULE_ID, 7, "内政-升级城内建筑", checkMap)
	DomesticUpgradeOuterCityBuilding = newOperateType(domestic.MODULE_ID, 8, "内政-升级外城建筑", checkMap)
	DomesticChangeHeroName           = newOperateType(domestic.MODULE_ID, 9, "内政-君主改名", checkMap)
	DomesticMiaoBuildingWorkerCd     = newOperateType(domestic.MODULE_ID, 10, "内政-建筑秒CD", checkMap)
	DomesticMiaoTechWorkerCd         = newOperateType(domestic.MODULE_ID, 11, "内政-科技秒CD", checkMap)
	DomesticForgingEquip             = newOperateType(domestic.MODULE_ID, 12, "", checkMap)
	DomesticWorkshopStartForge       = newOperateType(domestic.MODULE_ID, 13, "", checkMap)
	DomesticWorkshopCollect          = newOperateType(domestic.MODULE_ID, 14, "内政-铁匠铺领取装备", checkMap)
	DomesticCityEventExchange        = newOperateType(domestic.MODULE_ID, 15, "内政-请求城内事件兑换", checkMap)
	DomesticCollectCountdownPrize    = newOperateType(domestic.MODULE_ID, 16, "内政-领取马车奖励", checkMap)
	DomesticUnlockStableBuilding     = newOperateType(domestic.MODULE_ID, 17, "内政-解锁城内建筑", checkMap)
	DomesticCollectSeasonPrize       = newOperateType(domestic.MODULE_ID, 18, "内政-领取季节奖励", checkMap)
	DomesticRefreshWorkshop          = newOperateType(domestic.MODULE_ID, 19, "内政-铁匠铺刷新", checkMap)
	DomesticUpdateBuildingLevel      = newOperateType(domestic.MODULE_ID, 20, "内政-建筑升级", checkMap)
	DomesticLearnTechnology          = newOperateType(domestic.MODULE_ID, 21, "内政-学习科技", checkMap)
	DomesticWorkshopMiaoCd           = newOperateType(domestic.MODULE_ID, 22, "内政-铁匠铺秒CD", checkMap)
	DomesticUpdateOutcityType        = newOperateType(domestic.MODULE_ID, 23, "内政-外城改建", checkMap)
	DomesticBuySp                    = newOperateType(domestic.MODULE_ID, 24, "内政-购买体力值", checkMap)
	DomesticUseBufEffect             = newOperateType(domestic.MODULE_ID, 25, "内政-使用增益", checkMap)
	DomesticWorkerUnlock             = newOperateType(domestic.MODULE_ID, 26, "内政-解锁建筑队", checkMap)

	DungeonChallenge               = newOperateType(dungeon.MODULE_ID, 1, "推图-挑战", checkMap)
	DungeonCollectChapterPrize     = newOperateType(dungeon.MODULE_ID, 2, "推图-领取章节奖励", checkMap)
	DungeonCollectPassDungeonPrize = newOperateType(dungeon.MODULE_ID, 3, "推图-领取通关奖励", checkMap)
	DungeonAutoChallenge           = newOperateType(dungeon.MODULE_ID, 4, "推图-领取扫荡奖励", checkMap)
	DungeonCollectChapterStarPrize = newOperateType(dungeon.MODULE_ID, 5, "推图-领取章节星数奖励", checkMap)

	EquipmentUpgrade            = newOperateType(equipment.MODULE_ID, 1, "装备-升级", checkMap)
	EquipmentInheritUpgrade     = newOperateType(equipment.MODULE_ID, 2, "装备-继承升级", checkMap)
	EquipmentInheritRefine      = newOperateType(equipment.MODULE_ID, 3, "装备-继承强化", checkMap)
	EquipmentUpgradeAll         = newOperateType(equipment.MODULE_ID, 4, "装备-升级所有", checkMap)
	EquipmentRefined            = newOperateType(equipment.MODULE_ID, 5, "装备-强化", checkMap)
	EquipmentSmelt              = newOperateType(equipment.MODULE_ID, 6, "装备-熔炼", checkMap)
	EquipmentRebuild            = newOperateType(equipment.MODULE_ID, 7, "装备-重铸", checkMap)
	EquipmentForge              = newOperateType(equipment.MODULE_ID, 8, "装备-锻造", checkMap)
	EquipmentTaozUpgradeQuality = newOperateType(equipment.MODULE_ID, 9, "装备-套装升级品质", checkMap)

	FarmHarvest       = newOperateType(farm.MODULE_ID, 1, "农场-收获", checkMap)
	FarmChange        = newOperateType(farm.MODULE_ID, 2, "农场-改建", checkMap)
	FarmOneKeyHarvest = newOperateType(farm.MODULE_ID, 3, "农场-一键收获", checkMap)
	FarmOneKeyReset   = newOperateType(farm.MODULE_ID, 4, "农场-一键重置", checkMap)
	FarmSteal         = newOperateType(farm.MODULE_ID, 5, "农场-偷菜", checkMap)
	FarmOneKeySteal   = newOperateType(farm.MODULE_ID, 6, "农场-一键偷菜", checkMap)

	FishingFishing           = newOperateType(fishing.MODULE_ID, 1, "钓鱼", checkMap)
	FishingFishPointExchange = newOperateType(fishing.MODULE_ID, 2, "钓鱼-积分兑换", checkMap)

	GardenCollectTreasureTreePrize = newOperateType(garden.MODULE_ID, 1, "摇钱树-领奖", checkMap)
	GardenWaterTreasuryTree        = newOperateType(garden.MODULE_ID, 2, "摇钱树-浇水", checkMap)

	GemOneKeyCombineDepotGem = newOperateType(gem.MODULE_ID, 1, "宝石-背包一键合成", checkMap)
	GemOneKeyCombineGem      = newOperateType(gem.MODULE_ID, 2, "宝石-一键合成", checkMap)
	GemOneKeyUseGem          = newOperateType(gem.MODULE_ID, 3, "宝石-一键使用", checkMap)
	GemUseGem                = newOperateType(gem.MODULE_ID, 4, "宝石-使用", checkMap)
	GemCombineGem            = newOperateType(gem.MODULE_ID, 5, "宝石-合成", checkMap)

	GuildCreate                 = newOperateType(guild.MODULE_ID, 1, "联盟-创建", checkMap)
	GuildCollectGuildEventPrize = newOperateType(guild.MODULE_ID, 2, "联盟-领取事件奖励", checkMap)
	GuildCollectFullBigBox      = newOperateType(guild.MODULE_ID, 3, "联盟-领取大宝箱", checkMap)
	GuildDonate                 = newOperateType(guild.MODULE_ID, 4, "联盟-捐献", checkMap)
	GuildUpdateGuildName        = newOperateType(guild.MODULE_ID, 5, "联盟-改名", checkMap)
	GuildRequestJoin            = newOperateType(guild.MODULE_ID, 6, "联盟-请求加入", checkMap)
	GuildUpgradeLevel           = newOperateType(guild.MODULE_ID, 7, "联盟-升级", checkMap)
	GuildJoin                   = newOperateType(guild.MODULE_ID, 8, "联盟-加入", checkMap)
	GuildDismiss                = newOperateType(guild.MODULE_ID, 9, "联盟-踢人", checkMap)
	GuildLeave                  = newOperateType(guild.MODULE_ID, 10, "联盟-离开", checkMap)
	GuildImpeachLeader          = newOperateType(guild.MODULE_ID, 11, "联盟-弹劾", checkMap)
	GuildChangePrestigeTarget   = newOperateType(guild.MODULE_ID, 12, "联盟-修改成就目标", checkMap)
	GuildDestroyed              = newOperateType(guild.MODULE_ID, 13, "联盟-解散", checkMap)
	GuildVote                   = newOperateType(guild.MODULE_ID, 14, "", checkMap)
	GuildSendYinliangToMember   = newOperateType(guild.MODULE_ID, 15, "", checkMap)
	GuildSendYinliangToGuild    = newOperateType(guild.MODULE_ID, 16, "", checkMap)
	GuildPaySalary              = newOperateType(guild.MODULE_ID, 17, "", checkMap)
	GuildHelpMember             = newOperateType(guild.MODULE_ID, 18, "联盟-帮助", checkMap)
	GuildCollectRankPrize       = newOperateType(guild.MODULE_ID, 19, "联盟-领取排行奖励", checkMap)
	GuildCollectWorkshopPrize   = newOperateType(guild.MODULE_ID, 20, "联盟-领取工坊奖励", checkMap)
	GuildCollectTaskPrize       = newOperateType(guild.MODULE_ID, 21, "联盟-领取任务阶段奖励", checkMap)
	GuildChangeCountry          = newOperateType(guild.MODULE_ID, 22, "联盟-转国", checkMap)

	GuiZuCollectLevelPrize = newOperateType(guizu.MODULE_ID, 1, "贵族-领取等级奖励", checkMap)

	HebiRob         = newOperateType(hebi.MODULE_ID, 1, "合璧-抢夺", checkMap)
	HebiLeave       = newOperateType(hebi.MODULE_ID, 2, "合璧-离开", checkMap)
	HebiKick        = newOperateType(hebi.MODULE_ID, 3, "合璧-超时踢出", checkMap)
	HebiCheckInRoom = newOperateType(hebi.MODULE_ID, 4, "合璧-创建房间", checkMap)
	HebiJoinRoom    = newOperateType(hebi.MODULE_ID, 5, "合璧-加入房间", checkMap)
	HebiCopySelf    = newOperateType(hebi.MODULE_ID, 6, "合璧-分身", checkMap)
	HebiBeRobPos    = newOperateType(hebi.MODULE_ID, 7, "合璧-被抢位", checkMap)

	MailCollectMailPrize = newOperateType(mail.MODULE_ID, 1, "邮件-领取奖励", checkMap)
	MailReadMulti        = newOperateType(mail.MODULE_ID, 2, "邮件-批量读取", checkMap)

	MilitaryForceAddSoldier           = newOperateType(military.MODULE_ID, 1, "军事-强征", checkMap)
	MilitaryUpgradeSoldier            = newOperateType(military.MODULE_ID, 2, "军事-升级士兵", checkMap)
	MilitarySellSeekCaptain           = newOperateType(military.MODULE_ID, 3, "军事-招募武将", checkMap)
	MilitaryChangeCaptainName         = newOperateType(military.MODULE_ID, 4, "军事-武将改名", checkMap)
	MilitaryChangeCaptainRace         = newOperateType(military.MODULE_ID, 5, "军事-武将转职", checkMap)
	MilitaryCaptainRefined            = newOperateType(military.MODULE_ID, 6, "军事-武将强化", checkMap)
	MilitaryCaptainRebirthMiaoCd      = newOperateType(military.MODULE_ID, 7, "军事-武将强化秒 CD", checkMap)
	MilitaryJiuGuanConsult            = newOperateType(military.MODULE_ID, 8, "军事-酒馆请教", checkMap)
	MilitaryUseTrainingExpGoods       = newOperateType(military.MODULE_ID, 9, "军事-使用突飞令", checkMap)
	MilitaryJiuGuanRefresh            = newOperateType(military.MODULE_ID, 10, "军事-酒馆刷新", checkMap)
	MilitaryCaptainUpgradeAbilityExp  = newOperateType(military.MODULE_ID, 11, "军事-武将成长", checkMap)
	MilitaryCaptainRebirth            = newOperateType(military.MODULE_ID, 12, "军事-武将转生", checkMap)
	MilitaryCollectCaptainTrainingExp = newOperateType(military.MODULE_ID, 13, "军事-武将训练", checkMap)
	MilitaryUpdateCaptainOfficial     = newOperateType(military.MODULE_ID, 14, "军事-武将册封", checkMap)
	MilitaryLeaveCaptainOfficial      = newOperateType(military.MODULE_ID, 15, "军事-武将解职", checkMap)
	MilitaryGongXunGoods              = newOperateType(military.MODULE_ID, 16, "军事-使用功勋物品", checkMap)
	MilitaryCopyDefenserGoods         = newOperateType(military.MODULE_ID, 17, "军事-使用复制驻防物品", checkMap)
	MilitaryCaptainBorn               = newOperateType(military.MODULE_ID, 18, "军事-武将生成", checkMap)
	MilitaryCaptainUpstar             = newOperateType(military.MODULE_ID, 19, "军事-武将升星", checkMap)
	MilitaryCaptainLevelExchange      = newOperateType(military.MODULE_ID, 20, "军事-武将经验传承", checkMap)
	MilitaryCaptainReset              = newOperateType(military.MODULE_ID, 21, "军事-武将重生", checkMap)

	MingcWarSpeedUp = newOperateType(mingc_war.MODULE_ID, 1, "名城战加速", checkMap)

	MiscCollectChargePrize            = newOperateType(misc.MODULE_ID, 1, "领取充值奖励", checkMap)
	MiscCollectDailyBargain           = newOperateType(misc.MODULE_ID, 2, "领取每日特惠", checkMap)
	MiscActivateDurationCard          = newOperateType(misc.MODULE_ID, 3, "激活尊享卡", checkMap)
	MiscCollectDurationCardDailyPrize = newOperateType(misc.MODULE_ID, 4, "领取尊享卡每日奖励", checkMap)
	MiscActivateChargeObj             = newOperateType(misc.MODULE_ID, 5, "激活充值项", checkMap)

	PromotionCollectLoginPrize = newOperateType(promotion.MODULE_ID, 1, "领取累积登陆奖励", checkMap)
	PromotionCollectLevelFund  = newOperateType(promotion.MODULE_ID, 2, "领取君主等级基金", checkMap)
	PromotionBuyHeroLevelFund  = newOperateType(promotion.MODULE_ID, 3, "购买君主等级基金", checkMap)
	PromotionCollectFreeGift   = newOperateType(promotion.MODULE_ID, 4, "领取免费礼包", checkMap)
	PromotionBuyTimeLimitGift  = newOperateType(promotion.MODULE_ID, 5, "购买时限礼包", checkMap)
	PromotionCollectDailySp    = newOperateType(promotion.MODULE_ID, 6, "体力值领取", checkMap)

	QuestionPrize = newOperateType(question.MODULE_ID, 1, "", checkMap)

	RandomEventChooseOption = newOperateType(random_event.MODULE_ID, 1, "随机事件-选项", checkMap)

	RedPacketBuy             = newOperateType(red_packet.MODULE_ID, 1, "红包-购买", checkMap)
	RedPacketGrab            = newOperateType(red_packet.MODULE_ID, 2, "红包-抢", checkMap)
	RedPacketAllGrabbedPrize = newOperateType(red_packet.MODULE_ID, 3, "红包-抢光奖励", checkMap)

	RegionInvestigate                = newOperateType(region.MODULE_ID, 1, "野外-瞭望", checkMap)
	RegionFastMoveBase               = newOperateType(region.MODULE_ID, 2, "野外-使用迁城令", checkMap)
	RegionUseMianGoods               = newOperateType(region.MODULE_ID, 3, "野外-使用免战", checkMap)
	RegionSpeedUp                    = newOperateType(region.MODULE_ID, 4, "野外-加速", checkMap)
	RegionBuyProsperity              = newOperateType(region.MODULE_ID, 5, "野外-购买繁荣度", checkMap)
	RegionChallenge                  = newOperateType(region.MODULE_ID, 6, "", checkMap)
	RegionUseMultiLevelNpcTimesGoods = newOperateType(region.MODULE_ID, 7, "野外-使用讨伐令", checkMap)
	RegionUseInvaseHeroTimesGoods    = newOperateType(region.MODULE_ID, 8, "野外-使用攻城令", checkMap)
	RegionUseJunTuanNpcTimesGoods    = newOperateType(region.MODULE_ID, 9, "野外-使用军团令", checkMap)

	ShopBuyZhenBaoGeGoods       = newOperateType(shop.MODULE_ID, 1, "商店-银两商店购买", checkMap)
	ShopBuyNormalGoods          = newOperateType(shop.MODULE_ID, 2, "商店-普通商店购买", checkMap)
	ShopBuyBlackMarketGoods     = newOperateType(shop.MODULE_ID, 3, "商店-黑市购买", checkMap)
	ShopRefreshBlackMarketGoods = newOperateType(shop.MODULE_ID, 4, "商店-刷新黑市", checkMap)

	SecretTowerPrize    = newOperateType(secret_tower.MODULE_ID, 1, "密室-领取奖励", checkMap)
	SecretStartCallenge = newOperateType(secret_tower.MODULE_ID, 2, "", checkMap)

	StrategyUse = newOperateType(strategy.MODULE_ID, 1, "使用君主策略", checkMap)

	TaskCollectTaskPrize          = newOperateType(task.MODULE_ID, 1, "任务-领取任务奖励", checkMap)
	TaskCollectTaskBoxPrize       = newOperateType(task.MODULE_ID, 2, "任务-", checkMap)
	TaskCollectBaYeStagePrize     = newOperateType(task.MODULE_ID, 3, "任务-领取霸业阶段奖励", checkMap)
	TaskCollectBaYeTaskPrize      = newOperateType(task.MODULE_ID, 4, "任务-领取霸业奖励", checkMap)
	TaskCollectActiveDegreePrize  = newOperateType(task.MODULE_ID, 5, "任务-领取活跃度奖励", checkMap)
	TaskCollectAchieveStarPrize   = newOperateType(task.MODULE_ID, 6, "任务-领取成就星级奖励", checkMap)
	TaskCollectAchieveTaskPrize   = newOperateType(task.MODULE_ID, 7, "任务-领取成就任务奖励", checkMap)
	TaskCollectBwzlTaskPrize      = newOperateType(task.MODULE_ID, 8, "任务-领取霸王之路奖励", checkMap)
	TaskCollectBwzlTaskStagePrize = newOperateType(task.MODULE_ID, 9, "任务-领取霸王之路阶段奖励", checkMap)
	TaskCollectMainTaskPrize      = newOperateType(task.MODULE_ID, 10, "任务-领取主线任务奖励", checkMap)
	TaskCollectBranchTaskPrize    = newOperateType(task.MODULE_ID, 11, "任务-领取支线任务奖励", checkMap)
	TaskUpgradeTitle              = newOperateType(task.MODULE_ID, 12, "任务-升级任务称号", checkMap)
	TaskCollectActivityTaskPrize  = newOperateType(task.MODULE_ID, 13, "任务-领取活动任务奖励", checkMap)

	TeachCollectPrize = newOperateType(teach.MODULE_ID, 1, "教学关卡-领奖", checkMap)

	TowerChallenge     = newOperateType(tower.MODULE_ID, 1, "千重楼-挑战", checkMap)
	TowerAutoChallenge = newOperateType(tower.MODULE_ID, 2, "千重楼-扫荡", checkMap)

	VipCollectDailyPrize = newOperateType(vip.MODULE_ID, 1, "vip-每日礼包", checkMap)
	VipCollectLevelPrize = newOperateType(vip.MODULE_ID, 2, "vip-专属礼包", checkMap)
	VipBuyDungeonTimes   = newOperateType(vip.MODULE_ID, 3, "vip-购买推图次数", checkMap)

	XiongNuSucc   = newOperateType(xiongnu.MODULE_ID, 1, "匈奴-胜利", checkMap)
	XiongNuJoin   = newOperateType(xiongnu.MODULE_ID, 2, "匈奴-参战", checkMap)
	XiongNuSetDef = newOperateType(xiongnu.MODULE_ID, 3, "匈奴-设置防守", checkMap)
	XiongNuStart  = newOperateType(xiongnu.MODULE_ID, 4, "匈奴-开始", checkMap)

	XuanYuanCollectRankPrize = newOperateType(xuanyuan.MODULE_ID, 1, "轩辕会武-领取排名奖励", checkMap)

	ZhanJiangChallenge = newOperateType(zhanjiang.MODULE_ID, 1, "过关斩将-挑战", checkMap)

	ZhengWuCollect         = newOperateType(zhengwu.MODULE_ID, 1, "政务-领奖", checkMap)
	ZhengWuYuanBaoComplete = newOperateType(zhengwu.MODULE_ID, 2, "政务-元宝完成", checkMap)
	ZhengWuYuanBaoRefresh  = newOperateType(zhengwu.MODULE_ID, 3, "政务-元宝刷新", checkMap)
)

type OperateType struct {
	id         uint64
	operate    uint64
	subOperate uint64
	desc       string
}

func (t *OperateType) Oper() uint64 {
	return t.operate
}

func (t *OperateType) SubOper() uint64 {
	return t.subOperate
}

func (t *OperateType) Id() uint64 {
	return t.id
}

func newOperateType(operate uint64, subOperate uint64, desc string, checkMap map[uint64]struct{}) *OperateType {
	k := operate*100000 + subOperate
	if _, ok := checkMap[k]; ok {
		logrus.Panicf("newOperateType 操作类型重复。operate:%v subOperate:%v", operate, subOperate)
	}

	checkMap[k] = struct{}{}
	return &OperateType{id: k, operate: operate, subOperate: subOperate}
}

const (
	EquipUpgrade = iota + 1
	EquipRefine
)

const (
	EquipNoInherit = iota + 1
	EquipInherit
)

const (
	EquipOperTypeRebuild = iota + 1
	EquipOperTypeSmelt
)

const (
	MailSend     = iota + 1
	MailReceived
	MailRead
	MailDel
	MailCollect
)

const (
	GardenCareSelf    = iota + 1
	GardenCareFriends
)

const (
	RefreshAuto  = iota + 1
	RefreshMoney
)

const (
	SpeedUpMoney = iota + 1
	SpeedUpItem
	SpeedDonated
)

const (
	SNSAddFriend       = iota + 1
	SNSDelFriend
	SNSAddBlack
	SNSDelBlack
	SNSAddEnemy
	SNSDelEnemy
	SNSSetImportant
	SNSCancelImportant
)

const (
	Add    = iota
	Reduce
)

const (
	GoodsTypeEquip      = iota + 1
	GoodsTypeConsumable
	GoodsTypeGem
)

const (
	BattleHUANJING    = iota + 1
	BattleTower
	BattleBaiZhan
	BattleSecretTower
	BattleZhanJiang
	BattleInvade
	BattleAssist
	BattleExpel
	BattleSingle
	BattleMulti
	BattleMcWar
	BattleHebi
	BattleXuanYuan
	BattleTeach
)

const (
	BattleTypeAtk = iota + 1
	BattleTypeDef
)

const (
	BattleTypeLeader = iota + 1
	BattleTypeMember
)

const (
	ResearchTech    = iota + 1
	ResearchSolider
)

const (
	CountrySelectByName      = iota + 1
	CountrySelectBySelf
	CountrySelectByRecommend
	CountryChange
)

const (
	FarmTypePlant        = iota + 1
	FarmTypeHarvest
	FarmTypeEarlyHarvest
	FarmTypeSteal
	FarmTypeChange
	FarmTypeReset
)

const (
	FishNormal      = iota + 1
	FishYuanbao
	FishYuanbaoHalf
	FishYuanbaoTen
	FishUseGoods
)

const (
	TaskOperTypeStart    = iota + 1
	TaskOperTypeComplete
	TaskOperTypeCollect
)

const (
	TaskTypeBaYe         = iota + 1
	TaskTypeBwzl
	TaskTypeActiveDegree
	TaskTypeAchieve
	TaskTypeMain
	TaskTypeBranch
	TaskTypeActivity
)

const (
	CaptainOperTypeUpgrade    = iota
	CaptainOperTypeRefine
	CaptainOperTypeChangeRace
	CaptainOperTypeOfficial
)

func BuildGoodsType(t, subType uint64) uint64 {
	//return t*1000 + subType
	return t // 忽略子类型，直接使用大类
}
