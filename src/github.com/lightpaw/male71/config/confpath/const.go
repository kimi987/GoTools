package confpath

const (
	Import = 0

	// activitydata.ActivityCollectionData
	ActivityCollectionDataPath = "活动/收集活动.txt"
	ActivityCollectionDataKey  = "id"

	// activitydata.ActivityShowData
	ActivityShowDataPath = "活动/活动展示.txt"
	ActivityShowDataKey  = "id"

	// activitydata.ActivityTaskListModeData
	ActivityTaskListModeDataPath = "活动/列表式任务活动.txt"
	ActivityTaskListModeDataKey  = "id"

	// activitydata.CollectionExchangeData
	CollectionExchangeDataPath = "活动/收集兑换.txt"
	CollectionExchangeDataKey  = "id"

	// bai_zhan_data.JunXianLevelData
	JunXianLevelDataPath = "百战千军/等级.txt"
	JunXianLevelDataKey  = "level"

	// bai_zhan_data.JunXianPrizeData
	JunXianPrizeDataPath = "百战千军/等级奖励.txt"
	JunXianPrizeDataKey  = "id"

	// basedata.HomeNpcBaseData
	HomeNpcBaseDataPath = "地图/玩家主城野怪.txt"
	HomeNpcBaseDataKey  = "id"

	// basedata.NpcBaseData
	NpcBaseDataPath = "地图/野怪基础数据.txt"
	NpcBaseDataKey  = "id"

	// blockdata.BlockData
	BlockDataPath = "地图阻挡/阻挡块.txt"
	BlockDataKey  = "id"

	// body.BodyData
	BodyDataPath = "杂项/形象.txt"
	BodyDataKey  = "id"

	// buffer.BufferData
	BufferDataPath = "策略/增益.txt"
	BufferDataKey  = "id"

	// buffer.BufferTypeData
	BufferTypeDataPath = "策略/增益类型.txt"
	BufferTypeDataKey  = "id"

	// captain.CaptainAbilityData
	CaptainAbilityDataPath = "武将/武将成长.txt"
	CaptainAbilityDataKey  = "ability"

	// captain.CaptainData
	CaptainDataPath = "武将/武将.txt"
	CaptainDataKey  = "id"

	// captain.CaptainFriendshipData
	CaptainFriendshipDataPath = "武将/武将羁绊.txt"
	CaptainFriendshipDataKey  = "id"

	// captain.CaptainLevelData
	CaptainLevelDataPath = "武将/武将等级.txt"
	CaptainLevelDataKey  = "-"

	// captain.CaptainOfficialCountData
	CaptainOfficialCountDataPath = "武将/武将官职数量.txt"
	CaptainOfficialCountDataKey  = "hero_level"

	// captain.CaptainOfficialData
	CaptainOfficialDataPath = "武将/武将官职.txt"
	CaptainOfficialDataKey  = "id"

	// captain.CaptainRarityData
	CaptainRarityDataPath = "武将/稀有度.txt"
	CaptainRarityDataKey  = "id"

	// captain.CaptainRebirthLevelData
	CaptainRebirthLevelDataPath = "武将/武将转生.txt"
	CaptainRebirthLevelDataKey  = "level"

	// captain.CaptainStarData
	CaptainStarDataPath = "武将/武将星数.txt"
	CaptainStarDataKey  = "-"

	// captain.NamelessCaptainData
	NamelessCaptainDataPath = "武将/无名武将.txt"
	NamelessCaptainDataKey  = "id"

	// charge.ChargeObjData
	ChargeObjDataPath = "充值/充值项.txt"
	ChargeObjDataKey  = "id"

	// charge.ChargePrizeData
	ChargePrizeDataPath = "充值/充值奖励.txt"
	ChargePrizeDataKey  = "id"

	// charge.ProductData
	ProductDataPath = "充值/收费.txt"
	ProductDataKey  = "id"

	// combine.EquipCombineData
	EquipCombineDataPath = "合成/装备合成.txt"
	EquipCombineDataKey  = "id"

	// combine.GoodsCombineData
	GoodsCombineDataPath = "合成/物品合成.txt"
	GoodsCombineDataKey  = "id"

	// country.CountryData
	CountryDataPath = "国家/国家.txt"
	CountryDataKey  = "id"

	// country.CountryOfficialData
	CountryOfficialDataPath = "国家/官职.txt"
	CountryOfficialDataKey  = "id"

	// country.CountryOfficialNpcData
	CountryOfficialNpcDataPath = "国家/官职npc.txt"
	CountryOfficialNpcDataKey  = "id"

	// country.FamilyNameData
	FamilyNameDataPath = "国家/姓氏荐国.txt"
	FamilyNameDataKey  = "id"

	// data.BroadcastData
	BroadcastDataPath = "文字/广播.txt"
	BroadcastDataKey  = "id"

	// data.BuffEffectData
	BuffEffectDataPath = "杂项/buff.txt"
	BuffEffectDataKey  = "id"

	// data.ColorData
	ColorDataPath = "文字/品质颜色.txt"
	ColorDataKey  = "quality_key"

	// data.FamilyName
	FamilyNamePath = "武将名字/姓.txt"
	FamilyNameKey  = "family_name"

	// data.FemaleGivenName
	FemaleGivenNamePath = "武将名字/女名.txt"
	FemaleGivenNameKey  = "name"

	// data.HeroLevelSubData
	HeroLevelSubDataPath = "内政/君主等级.txt"
	HeroLevelSubDataKey  = "level"

	// data.MaleGivenName
	MaleGivenNamePath = "武将名字/男名.txt"
	MaleGivenNameKey  = "name"

	// data.SpriteStat
	SpriteStatPath = "杂项/属性.txt"
	SpriteStatKey  = "id"

	// data.Text
	TextPath = "文字/文本.txt"
	TextKey  = "id"

	// data.TimeRuleData
	TimeRuleDataPath = "杂项/时间规则.txt"
	TimeRuleDataKey  = "id"

	// domestic_data.BaseLevelData
	BaseLevelDataPath = "内政/主城等级.txt"
	BaseLevelDataKey  = "level"

	// domestic_data.BuildingData
	BuildingDataPath = "内政/建筑.txt"
	BuildingDataKey  = "id"

	// domestic_data.BuildingLayoutData
	BuildingLayoutDataPath = "内政/建筑布局.txt"
	BuildingLayoutDataKey  = "id"

	// domestic_data.BuildingUnlockData
	BuildingUnlockDataPath = "内政/建筑解锁.txt"
	BuildingUnlockDataKey  = "-"

	// domestic_data.CityEventData
	CityEventDataPath = "内政/城内事件.txt"
	CityEventDataKey  = "id"

	// domestic_data.CityEventLevelData
	CityEventLevelDataPath = "内政/城内事件等级.txt"
	CityEventLevelDataKey  = "base_level"

	// domestic_data.CombineCost
	CombineCostPath = "杂项/组合消耗.txt"
	CombineCostKey  = "id"

	// domestic_data.CountdownPrizeData
	CountdownPrizeDataPath = "内政/倒计时礼包.txt"
	CountdownPrizeDataKey  = "id"

	// domestic_data.CountdownPrizeDescData
	CountdownPrizeDescDataPath = "内政/倒计时礼包描述.txt"
	CountdownPrizeDescDataKey  = "id"

	// domestic_data.GuanFuLevelData
	GuanFuLevelDataPath = "内政/官府等级.txt"
	GuanFuLevelDataKey  = "level"

	// domestic_data.OuterCityBuildingData
	OuterCityBuildingDataPath = "内政/外城建筑.txt"
	OuterCityBuildingDataKey  = "id"

	// domestic_data.OuterCityBuildingDescData
	OuterCityBuildingDescDataPath = "内政/外城描述.txt"
	OuterCityBuildingDescDataKey  = "id"

	// domestic_data.OuterCityData
	OuterCityDataPath = "内政/外城.txt"
	OuterCityDataKey  = "-"

	// domestic_data.OuterCityLayoutData
	OuterCityLayoutDataPath = "内政/外城布局.txt"
	OuterCityLayoutDataKey  = "id"

	// domestic_data.ProsperityDamageBuffData
	ProsperityDamageBuffDataPath = "内政/繁荣度buff.txt"
	ProsperityDamageBuffDataKey  = "id"

	// domestic_data.SoldierLevelData
	SoldierLevelDataPath = "军事/士兵等级.txt"
	SoldierLevelDataKey  = "level"

	// domestic_data.TechnologyData
	TechnologyDataPath = "内政/科技.txt"
	TechnologyDataKey  = "id"

	// domestic_data.TieJiangPuLevelData
	TieJiangPuLevelDataPath = "内政/铁匠铺等级.txt"
	TieJiangPuLevelDataKey  = "level"

	// domestic_data.WorkshopDuration
	WorkshopDurationPath = "内政/装备作坊时间.txt"
	WorkshopDurationKey  = "id"

	// domestic_data.WorkshopLevelData
	WorkshopLevelDataPath = "内政/装备作坊.txt"
	WorkshopLevelDataKey  = "level"

	// domestic_data.WorkshopRefreshCost
	WorkshopRefreshCostPath = "内政/装备作坊刷新消耗.txt"
	WorkshopRefreshCostKey  = "id"

	// dungeon.DungeonChapterData
	DungeonChapterDataPath = "推图副本/推图副本章节.txt"
	DungeonChapterDataKey  = "id"

	// dungeon.DungeonData
	DungeonDataPath = "推图副本/推图副本.txt"
	DungeonDataKey  = "id"

	// dungeon.DungeonGuideTroopData
	DungeonGuideTroopDataPath = "推图副本/引导布阵.txt"
	DungeonGuideTroopDataKey  = "id"

	// farm.FarmMaxStealConfig
	FarmMaxStealConfigPath = "农场/农场偷菜上限.txt"
	FarmMaxStealConfigKey  = "guan_fu_level"

	// farm.FarmOneKeyConfig
	FarmOneKeyConfigPath = "农场/农场一键种植.txt"
	FarmOneKeyConfigKey  = "base_level"

	// farm.FarmResConfig
	FarmResConfigPath = "农场/农场资源.txt"
	FarmResConfigKey  = "id"

	// fishing_data.FishData
	FishDataPath = "钓鱼/钓鱼数据.txt"
	FishDataKey  = "id"

	// fishing_data.FishingCaptainProbabilityData
	FishingCaptainProbabilityDataPath = "钓鱼/金杆钓.txt"
	FishingCaptainProbabilityDataKey  = "captain_id"

	// fishing_data.FishingCostData
	FishingCostDataPath = "钓鱼/钓鱼消耗.txt"
	FishingCostDataKey  = "-"

	// fishing_data.FishingShowData
	FishingShowDataPath = "钓鱼/钓鱼展示.txt"
	FishingShowDataKey  = "id"

	// function.FunctionOpenData
	FunctionOpenDataPath = "功能开启/功能开启.txt"
	FunctionOpenDataKey  = "function_type"

	// gardendata.TreasuryTreeData
	TreasuryTreeDataPath = "花园/摇钱树.txt"
	TreasuryTreeDataKey  = "-"

	// goods.EquipmentData
	EquipmentDataPath = "物品/装备.txt"
	EquipmentDataKey  = "id"

	// goods.EquipmentLevelData
	EquipmentLevelDataPath = "物品/装备等级.txt"
	EquipmentLevelDataKey  = "level"

	// goods.EquipmentQualityData
	EquipmentQualityDataPath = "物品/装备品质.txt"
	EquipmentQualityDataKey  = "id"

	// goods.EquipmentRefinedData
	EquipmentRefinedDataPath = "物品/装备强化.txt"
	EquipmentRefinedDataKey  = "level"

	// goods.EquipmentTaozData
	EquipmentTaozDataPath = "物品/装备套装.txt"
	EquipmentTaozDataKey  = "level"

	// goods.GemData
	GemDataPath = "物品/宝石.txt"
	GemDataKey  = "id"

	// goods.GoodsData
	GoodsDataPath = "物品/物品.txt"
	GoodsDataKey  = "id"

	// goods.GoodsQuality
	GoodsQualityPath = "物品/物品品质.txt"
	GoodsQualityKey  = "level"

	// guild_data.GuildBigBoxData
	GuildBigBoxDataPath = "联盟/联盟大宝箱.txt"
	GuildBigBoxDataKey  = "id"

	// guild_data.GuildClassLevelData
	GuildClassLevelDataPath = "联盟/联盟职位.txt"
	GuildClassLevelDataKey  = "level"

	// guild_data.GuildClassTitleData
	GuildClassTitleDataPath = "联盟/联盟职称.txt"
	GuildClassTitleDataKey  = "id"

	// guild_data.GuildDonateData
	GuildDonateDataPath = "联盟/联盟捐献.txt"
	GuildDonateDataKey  = "-"

	// guild_data.GuildEventPrizeData
	GuildEventPrizeDataPath = "联盟/联盟盟友礼包.txt"
	GuildEventPrizeDataKey  = "id"

	// guild_data.GuildLevelCdrData
	GuildLevelCdrDataPath = "联盟/联盟升级加速.txt"
	GuildLevelCdrDataKey  = "-"

	// guild_data.GuildLevelData
	GuildLevelDataPath = "联盟/联盟等级.txt"
	GuildLevelDataKey  = "level"

	// guild_data.GuildLogData
	GuildLogDataPath = "联盟/联盟日志.txt"
	GuildLogDataKey  = "id"

	// guild_data.GuildPermissionShowData
	GuildPermissionShowDataPath = "联盟/联盟权限.txt"
	GuildPermissionShowDataKey  = "-"

	// guild_data.GuildPrestigeEventData
	GuildPrestigeEventDataPath = "联盟/联盟声望事件.txt"
	GuildPrestigeEventDataKey  = "-"

	// guild_data.GuildPrestigePrizeData
	GuildPrestigePrizeDataPath = "联盟/联盟声望礼包.txt"
	GuildPrestigePrizeDataKey  = "prestige"

	// guild_data.GuildRankPrizeData
	GuildRankPrizeDataPath = "联盟/联盟排行奖励.txt"
	GuildRankPrizeDataKey  = "rank"

	// guild_data.GuildTarget
	GuildTargetPath = "联盟/联盟目标.txt"
	GuildTargetKey  = "id"

	// guild_data.GuildTaskData
	GuildTaskDataPath = "联盟/联盟任务.txt"
	GuildTaskDataKey  = "id"

	// guild_data.GuildTaskEvaluateData
	GuildTaskEvaluateDataPath = "联盟/联盟任务评价.txt"
	GuildTaskEvaluateDataKey  = "id"

	// guild_data.GuildTechnologyData
	GuildTechnologyDataPath = "联盟/联盟科技.txt"
	GuildTechnologyDataKey  = "-"

	// guild_data.NpcGuildTemplate
	NpcGuildTemplatePath = "联盟/联盟Npc模板.txt"
	NpcGuildTemplateKey  = "id"

	// guild_data.NpcMemberData
	NpcMemberDataPath = "联盟/联盟Npc成员.txt"
	NpcMemberDataKey  = "id"

	// head.HeadData
	HeadDataPath = "杂项/头像.txt"
	HeadDataKey  = "id"

	// hebi.HebiPrizeData
	HebiPrizeDataPath = "合璧/合璧奖励.txt"
	HebiPrizeDataKey  = "-"

	// herodata.HeroLevelData
	HeroLevelDataPath = "内政/君主等级.txt"
	HeroLevelDataKey  = "level"

	// i18n.I18nData
	I18nDataPath = "i18n/语言.txt"
	I18nDataKey  = "id"

	// icon.Icon
	IconPath = "杂项/图标.txt"
	IconKey  = "id"

	// location.LocationData
	LocationDataPath = "杂项/省市.txt"
	LocationDataKey  = "id"

	// maildata.MailData
	MailDataPath = "文字/邮件.txt"
	MailDataKey  = "id"

	// military_data.JiuGuanData
	JiuGuanDataPath = "军事/酒馆.txt"
	JiuGuanDataKey  = "level"

	// military_data.JunYingLevelData
	JunYingLevelDataPath = "军事/军营.txt"
	JunYingLevelDataKey  = "level"

	// military_data.TrainingLevelData
	TrainingLevelDataPath = "军事/修炼馆等级.txt"
	TrainingLevelDataKey  = "level"

	// military_data.TutorData
	TutorDataPath = "军事/酒馆导师.txt"
	TutorDataKey  = "id"

	// mingcdata.McBuildAddSupportData
	McBuildAddSupportDataPath = "名城营建/增加民心.txt"
	McBuildAddSupportDataKey  = "bai_zhan_level"

	// mingcdata.McBuildGuildMemberPrizeData
	McBuildGuildMemberPrizeDataPath = "名城营建/联盟成员奖励.txt"
	McBuildGuildMemberPrizeDataKey  = "id"

	// mingcdata.McBuildMcSupportData
	McBuildMcSupportDataPath = "名城营建/名城民心.txt"
	McBuildMcSupportDataKey  = "level"

	// mingcdata.MingcBaseData
	MingcBaseDataPath = "名城战/名城.txt"
	MingcBaseDataKey  = "id"

	// mingcdata.MingcTimeData
	MingcTimeDataPath = "名城战/时间.txt"
	MingcTimeDataKey  = "id"

	// mingcdata.MingcWarBuildingData
	MingcWarBuildingDataPath = "名城战/据点.txt"
	MingcWarBuildingDataKey  = "id"

	// mingcdata.MingcWarDrumStatData
	MingcWarDrumStatDataPath = "名城战/鼓舞加成.txt"
	MingcWarDrumStatDataKey  = "id"

	// mingcdata.MingcWarMapData
	MingcWarMapDataPath = "名城战/地图.txt"
	MingcWarMapDataKey  = "id"

	// mingcdata.MingcWarMultiKillData
	MingcWarMultiKillDataPath = "名城战/连斩.txt"
	MingcWarMultiKillDataKey  = "multi_kill"

	// mingcdata.MingcWarNpcData
	MingcWarNpcDataPath = "名城战/初始城主.txt"
	MingcWarNpcDataKey  = "id"

	// mingcdata.MingcWarNpcGuildData
	MingcWarNpcGuildDataPath = "名城战/初始城主联盟.txt"
	MingcWarNpcGuildDataKey  = "id"

	// mingcdata.MingcWarSceneData
	MingcWarSceneDataPath = "名城战/场景.txt"
	MingcWarSceneDataKey  = "id"

	// mingcdata.MingcWarTouShiBuildingTargetData
	MingcWarTouShiBuildingTargetDataPath = "名城战/投石机目标.txt"
	MingcWarTouShiBuildingTargetDataKey  = "id"

	// mingcdata.MingcWarTroopLastBeatWhenFailData
	MingcWarTroopLastBeatWhenFailDataPath = "名城战/舍命一击.txt"
	MingcWarTroopLastBeatWhenFailDataKey  = "bai_zhan_level"

	// monsterdata.MonsterCaptainData
	MonsterCaptainDataPath = "怪物/怪物武将.txt"
	MonsterCaptainDataKey  = "id"

	// monsterdata.MonsterMasterData
	MonsterMasterDataPath = "怪物/怪物君主.txt"
	MonsterMasterDataKey  = "id"

	// promdata.DailyBargainData
	DailyBargainDataPath = "福利/每日特惠.txt"
	DailyBargainDataKey  = "id"

	// promdata.DurationCardData
	DurationCardDataPath = "福利/尊享卡.txt"
	DurationCardDataKey  = "id"

	// promdata.EventLimitGiftData
	EventLimitGiftDataPath = "福利/事件时限礼包.txt"
	EventLimitGiftDataKey  = "id"

	// promdata.FreeGiftData
	FreeGiftDataPath = "福利/免费礼包.txt"
	FreeGiftDataKey  = "id"

	// promdata.HeroLevelFundData
	HeroLevelFundDataPath = "福利/君主等级基金.txt"
	HeroLevelFundDataKey  = "level"

	// promdata.LoginDayData
	LoginDayDataPath = "福利/7日登陆奖励.txt"
	LoginDayDataKey  = "day"

	// promdata.SpCollectionData
	SpCollectionDataPath = "福利/体力领取.txt"
	SpCollectionDataKey  = "id"

	// promdata.TimeLimitGiftData
	TimeLimitGiftDataPath = "福利/定时时限礼包.txt"
	TimeLimitGiftDataKey  = "id"

	// promdata.TimeLimitGiftGroupData
	TimeLimitGiftGroupDataPath = "福利/定时时限礼包组.txt"
	TimeLimitGiftGroupDataKey  = "id"

	// pushdata.PushData
	PushDataPath = "杂项/推送.txt"
	PushDataKey  = "-"

	// pvetroop.PveTroopData
	PveTroopDataPath = "杂项/pve部队.txt"
	PveTroopDataKey  = "-"

	// question.QuestionData
	QuestionDataPath = "答题/答题问题.txt"
	QuestionDataKey  = "id"

	// question.QuestionPrizeData
	QuestionPrizeDataPath = "答题/答题奖励.txt"
	QuestionPrizeDataKey  = "score"

	// question.QuestionSayingData
	QuestionSayingDataPath = "答题/答题名言.txt"
	QuestionSayingDataKey  = "id"

	// race.RaceData
	RaceDataPath = "武将/职业.txt"
	RaceDataKey  = "id"

	// random_event.EventOptionData
	EventOptionDataPath = "随机事件/选项.txt"
	EventOptionDataKey  = "id"

	// random_event.EventPosition
	EventPositionPath = "随机事件/事件坐标.txt"
	EventPositionKey  = "id"

	// random_event.OptionPrize
	OptionPrizePath = "随机事件/选项奖励.txt"
	OptionPrizeKey  = "id"

	// random_event.RandomEventData
	RandomEventDataPath = "随机事件/随机事件.txt"
	RandomEventDataKey  = "id"

	// red_packet.RedPacketData
	RedPacketDataPath = "杂项/红包.txt"
	RedPacketDataKey  = "id"

	// regdata.AreaData
	AreaDataPath = "地图/区块链.txt"
	AreaDataKey  = "id"

	// regdata.AssemblyData
	AssemblyDataPath = "地图/集结.txt"
	AssemblyDataKey  = "id"

	// regdata.BaozNpcData
	BaozNpcDataPath = "地图/宝藏怪物.txt"
	BaozNpcDataKey  = "id"

	// regdata.JunTuanNpcData
	JunTuanNpcDataPath = "地图/军团怪物.txt"
	JunTuanNpcDataKey  = "id"

	// regdata.JunTuanNpcPlaceData
	JunTuanNpcPlaceDataPath = "地图/军团怪物刷新.txt"
	JunTuanNpcPlaceDataKey  = "id"

	// regdata.RegionAreaData
	RegionAreaDataPath = "地图/地区区域带.txt"
	RegionAreaDataKey  = "id"

	// regdata.RegionData
	RegionDataPath = "地图/地区.txt"
	RegionDataKey  = "-"

	// regdata.RegionMonsterData
	RegionMonsterDataPath = "地图/地区定点野怪.txt"
	RegionMonsterDataKey  = "id"

	// regdata.RegionMultiLevelNpcData
	RegionMultiLevelNpcDataPath = "地图/多等级野怪.txt"
	RegionMultiLevelNpcDataKey  = "id"

	// regdata.RegionMultiLevelNpcLevelData
	RegionMultiLevelNpcLevelDataPath = "地图/多等级野怪等级.txt"
	RegionMultiLevelNpcLevelDataKey  = "-"

	// regdata.RegionMultiLevelNpcTypeData
	RegionMultiLevelNpcTypeDataPath = "地图/多等级野怪类型.txt"
	RegionMultiLevelNpcTypeDataKey  = "-"

	// regdata.TroopDialogueData
	TroopDialogueDataPath = "地图/部队对话.txt"
	TroopDialogueDataKey  = "-"

	// regdata.TroopDialogueTextData
	TroopDialogueTextDataPath = "地图/部队对话文字.txt"
	TroopDialogueTextDataKey  = "id"

	// resdata.AmountShowSortData
	AmountShowSortDataPath = "杂项/展示排序.txt"
	AmountShowSortDataKey  = "id"

	// resdata.BaowuData
	BaowuDataPath = "物品/宝物.txt"
	BaowuDataKey  = "-"

	// resdata.ConditionPlunder
	ConditionPlunderPath = "杂项/条件掉落.txt"
	ConditionPlunderKey  = "id"

	// resdata.ConditionPlunderItem
	ConditionPlunderItemPath = "杂项/条件掉落项.txt"
	ConditionPlunderItemKey  = "id"

	// resdata.Cost
	CostPath = "杂项/消耗.txt"
	CostKey  = "id"

	// resdata.GuildLevelPrize
	GuildLevelPrizePath = "杂项/联盟等级奖励.txt"
	GuildLevelPrizeKey  = "-"

	// resdata.Plunder
	PlunderPath = "杂项/掉落.txt"
	PlunderKey  = "id"

	// resdata.PlunderGroup
	PlunderGroupPath = "杂项/掉落组.txt"
	PlunderGroupKey  = "id"

	// resdata.PlunderItem
	PlunderItemPath = "杂项/掉落项.txt"
	PlunderItemKey  = "id"

	// resdata.PlunderPrize
	PlunderPrizePath = "杂项/掉落_奖励.txt"
	PlunderPrizeKey  = "id"

	// resdata.Prize
	PrizePath = "杂项/奖励.txt"
	PrizeKey  = "id"

	// resdata.ResCaptainData
	ResCaptainDataPath = "武将/武将.txt"
	ResCaptainDataKey  = "id"

	// scene.CombatScene
	CombatScenePath = "地图/战斗场景.txt"
	CombatSceneKey  = "id"

	// season.SeasonData
	SeasonDataPath = "季节/季节.txt"
	SeasonDataKey  = "-"

	// settings.PrivacySettingData
	PrivacySettingDataPath = "设置/隐私设置.txt"
	PrivacySettingDataKey  = "id"

	// shop.BlackMarketData
	BlackMarketDataPath = "商店/黑市.txt"
	BlackMarketDataKey  = "id"

	// shop.BlackMarketGoodsData
	BlackMarketGoodsDataPath = "商店/黑市商品.txt"
	BlackMarketGoodsDataKey  = "id"

	// shop.BlackMarketGoodsGroupData
	BlackMarketGoodsGroupDataPath = "商店/黑市商品分组.txt"
	BlackMarketGoodsGroupDataKey  = "id"

	// shop.DiscountColorData
	DiscountColorDataPath = "商店/折扣颜色.txt"
	DiscountColorDataKey  = "discount"

	// shop.NormalShopGoods
	NormalShopGoodsPath = "商店/普通商品.txt"
	NormalShopGoodsKey  = "id"

	// shop.Shop
	ShopPath = "商店/商店.txt"
	ShopKey  = "type"

	// shop.ZhenBaoGeShopGoods
	ZhenBaoGeShopGoodsPath = "商店/珍宝阁商品.txt"
	ZhenBaoGeShopGoodsKey  = "id"

	// spell.PassiveSpellData
	PassiveSpellDataPath = "战斗/被动技能.txt"
	PassiveSpellDataKey  = "id"

	// spell.Spell
	SpellPath = "军事/技能.txt"
	SpellKey  = "id"

	// spell.SpellData
	SpellDataPath = "战斗/技能.txt"
	SpellDataKey  = "id"

	// spell.SpellFacadeData
	SpellFacadeDataPath = "战斗/技能盒.txt"
	SpellFacadeDataKey  = "id"

	// spell.StateData
	StateDataPath = "战斗/状态.txt"
	StateDataKey  = "id"

	// strategydata.StrategyData
	StrategyDataPath = "策略/策略.txt"
	StrategyDataKey  = "id"

	// strategydata.StrategyEffectData
	StrategyEffectDataPath = "策略/策略效果.txt"
	StrategyEffectDataKey  = "id"

	// strongerdata.StrongerData
	StrongerDataPath = "杂项/变强.txt"
	StrongerDataKey  = "-"

	// sub.BuildingEffectData
	BuildingEffectDataPath = "内政/建筑效果.txt"
	BuildingEffectDataKey  = "id"

	// survey.SurveyData
	SurveyDataPath = "问卷调查/问卷调查.txt"
	SurveyDataKey  = "id"

	// taskdata.AchieveTaskData
	AchieveTaskDataPath = "任务/成就任务.txt"
	AchieveTaskDataKey  = "id"

	// taskdata.AchieveTaskStarPrizeData
	AchieveTaskStarPrizeDataPath = "任务/成就星数奖励.txt"
	AchieveTaskStarPrizeDataKey  = "star"

	// taskdata.ActiveDegreePrizeData
	ActiveDegreePrizeDataPath = "任务/活跃度任务奖励.txt"
	ActiveDegreePrizeDataKey  = "degree"

	// taskdata.ActiveDegreeTaskData
	ActiveDegreeTaskDataPath = "任务/活跃度任务.txt"
	ActiveDegreeTaskDataKey  = "id"

	// taskdata.ActivityTaskData
	ActivityTaskDataPath = "活动/活动任务.txt"
	ActivityTaskDataKey  = "id"

	// taskdata.BaYeStageData
	BaYeStageDataPath = "任务/霸业任务阶段.txt"
	BaYeStageDataKey  = "stage"

	// taskdata.BaYeTaskData
	BaYeTaskDataPath = "任务/霸业任务.txt"
	BaYeTaskDataKey  = "id"

	// taskdata.BranchTaskData
	BranchTaskDataPath = "任务/支线任务.txt"
	BranchTaskDataKey  = "id"

	// taskdata.BwzlPrizeData
	BwzlPrizeDataPath = "任务/霸王之路奖励.txt"
	BwzlPrizeDataKey  = "collect_prize_task_count"

	// taskdata.BwzlTaskData
	BwzlTaskDataPath = "任务/霸王之路.txt"
	BwzlTaskDataKey  = "id"

	// taskdata.MainTaskData
	MainTaskDataPath = "任务/主线任务.txt"
	MainTaskDataKey  = "sequence"

	// taskdata.TaskBoxData
	TaskBoxDataPath = "任务/任务宝箱.txt"
	TaskBoxDataKey  = "id"

	// taskdata.TaskTargetData
	TaskTargetDataPath = "任务/任务目标.txt"
	TaskTargetDataKey  = "id"

	// taskdata.TitleData
	TitleDataPath = "任务/称号.txt"
	TitleDataKey  = "id"

	// taskdata.TitleTaskData
	TitleTaskDataPath = "任务/称号任务.txt"
	TitleTaskDataKey  = "id"

	// teach.TeachChapterData
	TeachChapterDataPath = "教学/关卡.txt"
	TeachChapterDataKey  = "id"

	// towerdata.SecretTowerData
	SecretTowerDataPath = "系统模块/重楼密室.txt"
	SecretTowerDataKey  = "id"

	// towerdata.SecretTowerWordsData
	SecretTowerWordsDataPath = "系统模块/重楼密室聊天.txt"
	SecretTowerWordsDataKey  = "id"

	// towerdata.TowerData
	TowerDataPath = "系统模块/千重楼.txt"
	TowerDataKey  = "floor"

	// vip.VipContinueDaysData
	VipContinueDaysDataPath = "vip/连续登录奖励.txt"
	VipContinueDaysDataKey  = "level"

	// vip.VipLevelData
	VipLevelDataPath = "vip/vip等级.txt"
	VipLevelDataKey  = "level"

	// xiongnu.ResistXiongNuData
	ResistXiongNuDataPath = "抗击匈奴/难度.txt"
	ResistXiongNuDataKey  = "level"

	// xiongnu.ResistXiongNuScoreData
	ResistXiongNuScoreDataPath = "抗击匈奴/评分.txt"
	ResistXiongNuScoreDataKey  = "level"

	// xiongnu.ResistXiongNuWaveData
	ResistXiongNuWaveDataPath = "抗击匈奴/攻城波次.txt"
	ResistXiongNuWaveDataKey  = "id"

	// xuanydata.XuanyuanRangeData
	XuanyuanRangeDataPath = "轩辕会武/积分区间.txt"
	XuanyuanRangeDataKey  = "id"

	// xuanydata.XuanyuanRankPrizeData
	XuanyuanRankPrizeDataPath = "轩辕会武/排名奖励.txt"
	XuanyuanRankPrizeDataKey  = "id"

	// zhanjiang.ZhanJiangChapterData
	ZhanJiangChapterDataPath = "过关斩将/章节.txt"
	ZhanJiangChapterDataKey  = "chapter_id"

	// zhanjiang.ZhanJiangData
	ZhanJiangDataPath = "过关斩将/小关卡.txt"
	ZhanJiangDataKey  = "id"

	// zhanjiang.ZhanJiangGuanQiaData
	ZhanJiangGuanQiaDataPath = "过关斩将/关卡.txt"
	ZhanJiangGuanQiaDataKey  = "id"

	// zhengwu.ZhengWuCompleteData
	ZhengWuCompleteDataPath = "政务/完成消耗.txt"
	ZhengWuCompleteDataKey  = "-"

	// zhengwu.ZhengWuData
	ZhengWuDataPath = "政务/政务.txt"
	ZhengWuDataKey  = "id"

	// zhengwu.ZhengWuRefreshData
	ZhengWuRefreshDataPath = "政务/刷新消耗.txt"
	ZhengWuRefreshDataKey  = "times"

	// bai_zhan_data.BaiZhanMiscData
	BaiZhanMiscDataPath = "百战千军/杂项.txt"

	// combatdata.CombatConfig
	CombatConfigPath = "战斗/杂项.txt"

	// combatdata.CombatMiscConfig
	CombatMiscConfigPath = "战斗/杂项.txt"

	// combine.EquipCombineDatas
	EquipCombineDatasPath = "合成/装备合成.txt"

	// country.CountryMiscData
	CountryMiscDataPath = "国家/国家杂项.txt"

	// data.BroadcastHelp
	BroadcastHelpPath = "文字/广播.txt"

	// data.TextHelp
	TextHelpPath = "文字/文本.txt"

	// dianquan.ExchangeMiscData
	ExchangeMiscDataPath = "点券/点券杂项.txt"

	// domestic_data.BuildingLayoutMiscData
	BuildingLayoutMiscDataPath = "内政/建筑布局杂项.txt"

	// domestic_data.CityEventMiscData
	CityEventMiscDataPath = "内政/城内事件杂项.txt"

	// domestic_data.MainCityMiscData
	MainCityMiscDataPath = "内政/主城杂项.txt"

	// dungeon.DungeonMiscData
	DungeonMiscDataPath = "推图副本/推图副本杂项.txt"

	// farm.FarmMiscConfig
	FarmMiscConfigPath = "农场/农场杂项.txt"

	// fishing_data.FishRandomer
	FishRandomerPath = "钓鱼/钓鱼数据.txt"

	// gardendata.GardenConfig
	GardenConfigPath = "花园/杂项.txt"

	// goods.EquipmentTaozConfig
	EquipmentTaozConfigPath = "物品/taoz.txt"

	// goods.GemDatas
	GemDatasPath = "物品/宝石.txt"

	// goods.GoodsCheck
	GoodsCheckPath = "物品/物品检查.txt"

	// guild_data.GuildLogHelp
	GuildLogHelpPath = "联盟/联盟日志.txt"

	// guild_data.NpcGuildSuffixName
	NpcGuildSuffixNamePath = "联盟/联盟Npc后缀名.txt"

	// hebi.HebiMiscData
	HebiMiscDataPath = "合璧/合璧杂项.txt"

	// heroinit.HeroCreateData
	HeroCreateDataPath = "杂项/创建英雄基础数据.txt"

	// heroinit.HeroInitData
	HeroInitDataPath = "singleton/hero_init_data.txt"

	// maildata.MailHelp
	MailHelpPath = "文字/邮件.txt"

	// military_data.JiuGuanMiscData
	JiuGuanMiscDataPath = "军事/酒馆杂项.txt"

	// military_data.JunYingMiscData
	JunYingMiscDataPath = "军事/军营杂项.txt"

	// mingcdata.McBuildMiscData
	McBuildMiscDataPath = "名城营建/营建杂项.txt"

	// mingcdata.MingcMiscData
	MingcMiscDataPath = "名城战/杂项.txt"

	// promdata.EventLimitGiftConfig
	EventLimitGiftConfigPath = "福利/事件时限礼包.txt"

	// promdata.PromotionMiscData
	PromotionMiscDataPath = "福利/福利杂项.txt"

	// question.QuestionMiscData
	QuestionMiscDataPath = "答题/答题杂项.txt"

	// race.RaceConfig
	RaceConfigPath = "杂项/职业杂项.txt"

	// random_event.RandomEventDataDictionary
	RandomEventDataDictionaryPath = "随机事件/随机事件.txt"

	// random_event.RandomEventPositionDictionary
	RandomEventPositionDictionaryPath = "随机事件/事件坐标.txt"

	// rank_data.RankMiscData
	RankMiscDataPath = "杂项/排行榜杂项.txt"

	// regdata.JunTuanNpcPlaceConfig
	JunTuanNpcPlaceConfigPath = "地图/军团怪物.txt"

	// season.SeasonMiscData
	SeasonMiscDataPath = "季节/杂项.txt"

	// settings.SettingMiscData
	SettingMiscDataPath = "设置/默认设置.txt"

	// shop.ShopMiscData
	ShopMiscDataPath = "商店/商店杂项.txt"

	// singleton.GoodsConfig
	GoodsConfigPath = "物品/物品杂项.txt"

	// singleton.GuildConfig
	GuildConfigPath = "联盟/联盟杂项.txt"

	// singleton.GuildGenConfig
	GuildGenConfigPath = "联盟/联盟杂项.txt"

	// singleton.MilitaryConfig
	MilitaryConfigPath = "军事/军事杂项.txt"

	// singleton.MiscConfig
	MiscConfigPath = "杂项/杂项.txt"

	// singleton.MiscGenConfig
	MiscGenConfigPath = "杂项/杂项.txt"

	// singleton.RegionConfig
	RegionConfigPath = "地图/地区杂项.txt"

	// singleton.RegionGenConfig
	RegionGenConfigPath = "地图/地区杂项.txt"

	// tag.TagMiscData
	TagMiscDataPath = "杂项/标签杂项.txt"

	// taskdata.TaskMiscData
	TaskMiscDataPath = "任务/杂项.txt"

	// towerdata.SecretTowerMiscData
	SecretTowerMiscDataPath = "系统模块/重楼密室杂项.txt"

	// vip.VipMiscData
	VipMiscDataPath = "vip/vip杂项.txt"

	// xiongnu.ResistXiongNuMisc
	ResistXiongNuMiscPath = "抗击匈奴/杂项.txt"

	// xuanydata.XuanyuanMiscData
	XuanyuanMiscDataPath = "轩辕会武/杂项.txt"

	// zhanjiang.ZhanJiangMiscData
	ZhanJiangMiscDataPath = "过关斩将/其他.txt"

	// zhengwu.ZhengWuMiscData
	ZhengWuMiscDataPath = "政务/其他.txt"

	// zhengwu.ZhengWuRandomData
	ZhengWuRandomDataPath = "政务/其他.txt"
)
