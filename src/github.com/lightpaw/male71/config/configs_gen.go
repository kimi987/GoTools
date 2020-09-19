package config

import (
	"github.com/lightpaw/config"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/activitydata"
	"github.com/lightpaw/male7/config/bai_zhan_data"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/blockdata"
	"github.com/lightpaw/male7/config/body"
	"github.com/lightpaw/male7/config/buffer"
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/charge"
	"github.com/lightpaw/male7/config/combatdata"
	"github.com/lightpaw/male7/config/combine"
	"github.com/lightpaw/male7/config/country"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/dianquan"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/config/dungeon"
	"github.com/lightpaw/male7/config/farm"
	"github.com/lightpaw/male7/config/fishing_data"
	"github.com/lightpaw/male7/config/function"
	"github.com/lightpaw/male7/config/gardendata"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/head"
	"github.com/lightpaw/male7/config/hebi"
	"github.com/lightpaw/male7/config/herodata"
	"github.com/lightpaw/male7/config/heroinit"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/location"
	"github.com/lightpaw/male7/config/maildata"
	"github.com/lightpaw/male7/config/military_data"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/promdata"
	"github.com/lightpaw/male7/config/pushdata"
	"github.com/lightpaw/male7/config/pvetroop"
	"github.com/lightpaw/male7/config/question"
	"github.com/lightpaw/male7/config/race"
	"github.com/lightpaw/male7/config/random_event"
	"github.com/lightpaw/male7/config/rank_data"
	"github.com/lightpaw/male7/config/red_packet"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/config/season"
	"github.com/lightpaw/male7/config/settings"
	"github.com/lightpaw/male7/config/shop"
	"github.com/lightpaw/male7/config/singleton"
	"github.com/lightpaw/male7/config/spell"
	"github.com/lightpaw/male7/config/strategydata"
	"github.com/lightpaw/male7/config/strongerdata"
	"github.com/lightpaw/male7/config/survey"
	"github.com/lightpaw/male7/config/tag"
	"github.com/lightpaw/male7/config/taskdata"
	"github.com/lightpaw/male7/config/teach"
	"github.com/lightpaw/male7/config/towerdata"
	"github.com/lightpaw/male7/config/vip"
	"github.com/lightpaw/male7/config/xiongnu"
	"github.com/lightpaw/male7/config/xuanydata"
	"github.com/lightpaw/male7/config/zhanjiang"
	"github.com/lightpaw/male7/config/zhengwu"
	"github.com/lightpaw/male7/pb/shared_proto"
	"sort"
)

func NewConfigDatas() *ConfigDatas {
	c, err := LoadConfigDatas("conf")
	if err != nil {
		logrus.WithError(err).Panic("读取配置文件conf出错")
	}

	return c
}

func LoadConfigDatas(folder string) (*ConfigDatas, error) {

	gos, err := config.NewConfigGameObjects(folder)
	if err != nil {
		return nil, err
	}

	c, err := ParseConfigDatas(gos)
	if err != nil {
		return nil, err
	}

	return c, nil
}

//gogen:iface
type ConfigDatas struct {
	activityCollectionData *ActivityCollectionDataConfig

	activityShowData *ActivityShowDataConfig

	activityTaskListModeData *ActivityTaskListModeDataConfig

	collectionExchangeData *CollectionExchangeDataConfig

	junXianLevelData *JunXianLevelDataConfig

	junXianPrizeData *JunXianPrizeDataConfig

	homeNpcBaseData *HomeNpcBaseDataConfig

	npcBaseData *NpcBaseDataConfig

	blockData *BlockDataConfig

	bodyData *BodyDataConfig

	bufferData *BufferDataConfig

	bufferTypeData *BufferTypeDataConfig

	captainAbilityData *CaptainAbilityDataConfig

	captainData *CaptainDataConfig

	captainFriendshipData *CaptainFriendshipDataConfig

	captainLevelData *CaptainLevelDataConfig

	captainOfficialCountData *CaptainOfficialCountDataConfig

	captainOfficialData *CaptainOfficialDataConfig

	captainRarityData *CaptainRarityDataConfig

	captainRebirthLevelData *CaptainRebirthLevelDataConfig

	captainStarData *CaptainStarDataConfig

	namelessCaptainData *NamelessCaptainDataConfig

	chargeObjData *ChargeObjDataConfig

	chargePrizeData *ChargePrizeDataConfig

	productData *ProductDataConfig

	equipCombineData *EquipCombineDataConfig

	goodsCombineData *GoodsCombineDataConfig

	countryData *CountryDataConfig

	countryOfficialData *CountryOfficialDataConfig

	countryOfficialNpcData *CountryOfficialNpcDataConfig

	familyNameData *FamilyNameDataConfig

	broadcastData *BroadcastDataConfig

	buffEffectData *BuffEffectDataConfig

	colorData *ColorDataConfig

	familyName *FamilyNameConfig

	femaleGivenName *FemaleGivenNameConfig

	heroLevelSubData *HeroLevelSubDataConfig

	maleGivenName *MaleGivenNameConfig

	spriteStat *SpriteStatConfig

	text *TextConfig

	timeRuleData *TimeRuleDataConfig

	baseLevelData *BaseLevelDataConfig

	buildingData *BuildingDataConfig

	buildingLayoutData *BuildingLayoutDataConfig

	buildingUnlockData *BuildingUnlockDataConfig

	cityEventData *CityEventDataConfig

	cityEventLevelData *CityEventLevelDataConfig

	combineCost *CombineCostConfig

	countdownPrizeData *CountdownPrizeDataConfig

	countdownPrizeDescData *CountdownPrizeDescDataConfig

	guanFuLevelData *GuanFuLevelDataConfig

	outerCityBuildingData *OuterCityBuildingDataConfig

	outerCityBuildingDescData *OuterCityBuildingDescDataConfig

	outerCityData *OuterCityDataConfig

	outerCityLayoutData *OuterCityLayoutDataConfig

	prosperityDamageBuffData *ProsperityDamageBuffDataConfig

	soldierLevelData *SoldierLevelDataConfig

	technologyData *TechnologyDataConfig

	tieJiangPuLevelData *TieJiangPuLevelDataConfig

	workshopDuration *WorkshopDurationConfig

	workshopLevelData *WorkshopLevelDataConfig

	workshopRefreshCost *WorkshopRefreshCostConfig

	dungeonChapterData *DungeonChapterDataConfig

	dungeonData *DungeonDataConfig

	dungeonGuideTroopData *DungeonGuideTroopDataConfig

	farmMaxStealConfig *FarmMaxStealConfigConfig

	farmOneKeyConfig *FarmOneKeyConfigConfig

	farmResConfig *FarmResConfigConfig

	fishData *FishDataConfig

	fishingCaptainProbabilityData *FishingCaptainProbabilityDataConfig

	fishingCostData *FishingCostDataConfig

	fishingShowData *FishingShowDataConfig

	functionOpenData *FunctionOpenDataConfig

	treasuryTreeData *TreasuryTreeDataConfig

	equipmentData *EquipmentDataConfig

	equipmentLevelData *EquipmentLevelDataConfig

	equipmentQualityData *EquipmentQualityDataConfig

	equipmentRefinedData *EquipmentRefinedDataConfig

	equipmentTaozData *EquipmentTaozDataConfig

	gemData *GemDataConfig

	goodsData *GoodsDataConfig

	goodsQuality *GoodsQualityConfig

	guildBigBoxData *GuildBigBoxDataConfig

	guildClassLevelData *GuildClassLevelDataConfig

	guildClassTitleData *GuildClassTitleDataConfig

	guildDonateData *GuildDonateDataConfig

	guildEventPrizeData *GuildEventPrizeDataConfig

	guildLevelCdrData *GuildLevelCdrDataConfig

	guildLevelData *GuildLevelDataConfig

	guildLogData *GuildLogDataConfig

	guildPermissionShowData *GuildPermissionShowDataConfig

	guildPrestigeEventData *GuildPrestigeEventDataConfig

	guildPrestigePrizeData *GuildPrestigePrizeDataConfig

	guildRankPrizeData *GuildRankPrizeDataConfig

	guildTarget *GuildTargetConfig

	guildTaskData *GuildTaskDataConfig

	guildTaskEvaluateData *GuildTaskEvaluateDataConfig

	guildTechnologyData *GuildTechnologyDataConfig

	npcGuildTemplate *NpcGuildTemplateConfig

	npcMemberData *NpcMemberDataConfig

	headData *HeadDataConfig

	hebiPrizeData *HebiPrizeDataConfig

	heroLevelData *HeroLevelDataConfig

	i18nData *I18nDataConfig

	icon *IconConfig

	locationData *LocationDataConfig

	mailData *MailDataConfig

	jiuGuanData *JiuGuanDataConfig

	junYingLevelData *JunYingLevelDataConfig

	trainingLevelData *TrainingLevelDataConfig

	tutorData *TutorDataConfig

	mcBuildAddSupportData *McBuildAddSupportDataConfig

	mcBuildGuildMemberPrizeData *McBuildGuildMemberPrizeDataConfig

	mcBuildMcSupportData *McBuildMcSupportDataConfig

	mingcBaseData *MingcBaseDataConfig

	mingcTimeData *MingcTimeDataConfig

	mingcWarBuildingData *MingcWarBuildingDataConfig

	mingcWarDrumStatData *MingcWarDrumStatDataConfig

	mingcWarMapData *MingcWarMapDataConfig

	mingcWarMultiKillData *MingcWarMultiKillDataConfig

	mingcWarNpcData *MingcWarNpcDataConfig

	mingcWarNpcGuildData *MingcWarNpcGuildDataConfig

	mingcWarSceneData *MingcWarSceneDataConfig

	mingcWarTouShiBuildingTargetData *MingcWarTouShiBuildingTargetDataConfig

	mingcWarTroopLastBeatWhenFailData *MingcWarTroopLastBeatWhenFailDataConfig

	monsterCaptainData *MonsterCaptainDataConfig

	monsterMasterData *MonsterMasterDataConfig

	dailyBargainData *DailyBargainDataConfig

	durationCardData *DurationCardDataConfig

	eventLimitGiftData *EventLimitGiftDataConfig

	freeGiftData *FreeGiftDataConfig

	heroLevelFundData *HeroLevelFundDataConfig

	loginDayData *LoginDayDataConfig

	spCollectionData *SpCollectionDataConfig

	timeLimitGiftData *TimeLimitGiftDataConfig

	timeLimitGiftGroupData *TimeLimitGiftGroupDataConfig

	pushData *PushDataConfig

	pveTroopData *PveTroopDataConfig

	questionData *QuestionDataConfig

	questionPrizeData *QuestionPrizeDataConfig

	questionSayingData *QuestionSayingDataConfig

	raceData *RaceDataConfig

	eventOptionData *EventOptionDataConfig

	eventPosition *EventPositionConfig

	optionPrize *OptionPrizeConfig

	randomEventData *RandomEventDataConfig

	redPacketData *RedPacketDataConfig

	areaData *AreaDataConfig

	assemblyData *AssemblyDataConfig

	baozNpcData *BaozNpcDataConfig

	junTuanNpcData *JunTuanNpcDataConfig

	junTuanNpcPlaceData *JunTuanNpcPlaceDataConfig

	regionAreaData *RegionAreaDataConfig

	regionData *RegionDataConfig

	regionMonsterData *RegionMonsterDataConfig

	regionMultiLevelNpcData *RegionMultiLevelNpcDataConfig

	regionMultiLevelNpcLevelData *RegionMultiLevelNpcLevelDataConfig

	regionMultiLevelNpcTypeData *RegionMultiLevelNpcTypeDataConfig

	troopDialogueData *TroopDialogueDataConfig

	troopDialogueTextData *TroopDialogueTextDataConfig

	amountShowSortData *AmountShowSortDataConfig

	baowuData *BaowuDataConfig

	conditionPlunder *ConditionPlunderConfig

	conditionPlunderItem *ConditionPlunderItemConfig

	cost *CostConfig

	guildLevelPrize *GuildLevelPrizeConfig

	plunder *PlunderConfig

	plunderGroup *PlunderGroupConfig

	plunderItem *PlunderItemConfig

	plunderPrize *PlunderPrizeConfig

	prize *PrizeConfig

	resCaptainData *ResCaptainDataConfig

	combatScene *CombatSceneConfig

	seasonData *SeasonDataConfig

	privacySettingData *PrivacySettingDataConfig

	blackMarketData *BlackMarketDataConfig

	blackMarketGoodsData *BlackMarketGoodsDataConfig

	blackMarketGoodsGroupData *BlackMarketGoodsGroupDataConfig

	discountColorData *DiscountColorDataConfig

	normalShopGoods *NormalShopGoodsConfig

	shop *ShopConfig

	zhenBaoGeShopGoods *ZhenBaoGeShopGoodsConfig

	passiveSpellData *PassiveSpellDataConfig

	spell *SpellConfig

	spellData *SpellDataConfig

	spellFacadeData *SpellFacadeDataConfig

	stateData *StateDataConfig

	strategyData *StrategyDataConfig

	strategyEffectData *StrategyEffectDataConfig

	strongerData *StrongerDataConfig

	buildingEffectData *BuildingEffectDataConfig

	surveyData *SurveyDataConfig

	achieveTaskData *AchieveTaskDataConfig

	achieveTaskStarPrizeData *AchieveTaskStarPrizeDataConfig

	activeDegreePrizeData *ActiveDegreePrizeDataConfig

	activeDegreeTaskData *ActiveDegreeTaskDataConfig

	activityTaskData *ActivityTaskDataConfig

	baYeStageData *BaYeStageDataConfig

	baYeTaskData *BaYeTaskDataConfig

	branchTaskData *BranchTaskDataConfig

	bwzlPrizeData *BwzlPrizeDataConfig

	bwzlTaskData *BwzlTaskDataConfig

	mainTaskData *MainTaskDataConfig

	taskBoxData *TaskBoxDataConfig

	taskTargetData *TaskTargetDataConfig

	titleData *TitleDataConfig

	titleTaskData *TitleTaskDataConfig

	teachChapterData *TeachChapterDataConfig

	secretTowerData *SecretTowerDataConfig

	secretTowerWordsData *SecretTowerWordsDataConfig

	towerData *TowerDataConfig

	vipContinueDaysData *VipContinueDaysDataConfig

	vipLevelData *VipLevelDataConfig

	resistXiongNuData *ResistXiongNuDataConfig

	resistXiongNuScoreData *ResistXiongNuScoreDataConfig

	resistXiongNuWaveData *ResistXiongNuWaveDataConfig

	xuanyuanRangeData *XuanyuanRangeDataConfig

	xuanyuanRankPrizeData *XuanyuanRankPrizeDataConfig

	zhanJiangChapterData *ZhanJiangChapterDataConfig

	zhanJiangData *ZhanJiangDataConfig

	zhanJiangGuanQiaData *ZhanJiangGuanQiaDataConfig

	zhengWuCompleteData *ZhengWuCompleteDataConfig

	zhengWuData *ZhengWuDataConfig

	zhengWuRefreshData *ZhengWuRefreshDataConfig

	baiZhanMiscData       *bai_zhan_data.BaiZhanMiscData
	baiZhanMiscDataParser *config.ObjectParser

	combatConfig       *combatdata.CombatConfig
	combatConfigParser *config.ObjectParser

	combatMiscConfig       *combatdata.CombatMiscConfig
	combatMiscConfigParser *config.ObjectParser

	equipCombineDatas       *combine.EquipCombineDatas
	equipCombineDatasParser *config.ObjectParser

	countryMiscData       *country.CountryMiscData
	countryMiscDataParser *config.ObjectParser

	broadcastHelp       *data.BroadcastHelp
	broadcastHelpParser *config.ObjectParser

	textHelp       *data.TextHelp
	textHelpParser *config.ObjectParser

	exchangeMiscData       *dianquan.ExchangeMiscData
	exchangeMiscDataParser *config.ObjectParser

	buildingLayoutMiscData       *domestic_data.BuildingLayoutMiscData
	buildingLayoutMiscDataParser *config.ObjectParser

	cityEventMiscData       *domestic_data.CityEventMiscData
	cityEventMiscDataParser *config.ObjectParser

	mainCityMiscData       *domestic_data.MainCityMiscData
	mainCityMiscDataParser *config.ObjectParser

	dungeonMiscData       *dungeon.DungeonMiscData
	dungeonMiscDataParser *config.ObjectParser

	farmMiscConfig       *farm.FarmMiscConfig
	farmMiscConfigParser *config.ObjectParser

	fishRandomer       *fishing_data.FishRandomer
	fishRandomerParser *config.ObjectParser

	gardenConfig       *gardendata.GardenConfig
	gardenConfigParser *config.ObjectParser

	equipmentTaozConfig       *goods.EquipmentTaozConfig
	equipmentTaozConfigParser *config.ObjectParser

	gemDatas       *goods.GemDatas
	gemDatasParser *config.ObjectParser

	goodsCheck       *goods.GoodsCheck
	goodsCheckParser *config.ObjectParser

	guildLogHelp       *guild_data.GuildLogHelp
	guildLogHelpParser *config.ObjectParser

	npcGuildSuffixName       *guild_data.NpcGuildSuffixName
	npcGuildSuffixNameParser *config.ObjectParser

	hebiMiscData       *hebi.HebiMiscData
	hebiMiscDataParser *config.ObjectParser

	heroCreateData       *heroinit.HeroCreateData
	heroCreateDataParser *config.ObjectParser

	heroInitData       *heroinit.HeroInitData
	heroInitDataParser *config.ObjectParser

	mailHelp       *maildata.MailHelp
	mailHelpParser *config.ObjectParser

	jiuGuanMiscData       *military_data.JiuGuanMiscData
	jiuGuanMiscDataParser *config.ObjectParser

	junYingMiscData       *military_data.JunYingMiscData
	junYingMiscDataParser *config.ObjectParser

	mcBuildMiscData       *mingcdata.McBuildMiscData
	mcBuildMiscDataParser *config.ObjectParser

	mingcMiscData       *mingcdata.MingcMiscData
	mingcMiscDataParser *config.ObjectParser

	eventLimitGiftConfig       *promdata.EventLimitGiftConfig
	eventLimitGiftConfigParser *config.ObjectParser

	promotionMiscData       *promdata.PromotionMiscData
	promotionMiscDataParser *config.ObjectParser

	questionMiscData       *question.QuestionMiscData
	questionMiscDataParser *config.ObjectParser

	raceConfig       *race.RaceConfig
	raceConfigParser *config.ObjectParser

	randomEventDataDictionary       *random_event.RandomEventDataDictionary
	randomEventDataDictionaryParser *config.ObjectParser

	randomEventPositionDictionary       *random_event.RandomEventPositionDictionary
	randomEventPositionDictionaryParser *config.ObjectParser

	rankMiscData       *rank_data.RankMiscData
	rankMiscDataParser *config.ObjectParser

	junTuanNpcPlaceConfig       *regdata.JunTuanNpcPlaceConfig
	junTuanNpcPlaceConfigParser *config.ObjectParser

	seasonMiscData       *season.SeasonMiscData
	seasonMiscDataParser *config.ObjectParser

	settingMiscData       *settings.SettingMiscData
	settingMiscDataParser *config.ObjectParser

	shopMiscData       *shop.ShopMiscData
	shopMiscDataParser *config.ObjectParser

	goodsConfig       *singleton.GoodsConfig
	goodsConfigParser *config.ObjectParser

	guildConfig       *singleton.GuildConfig
	guildConfigParser *config.ObjectParser

	guildGenConfig       *singleton.GuildGenConfig
	guildGenConfigParser *config.ObjectParser

	militaryConfig       *singleton.MilitaryConfig
	militaryConfigParser *config.ObjectParser

	miscConfig       *singleton.MiscConfig
	miscConfigParser *config.ObjectParser

	miscGenConfig       *singleton.MiscGenConfig
	miscGenConfigParser *config.ObjectParser

	regionConfig       *singleton.RegionConfig
	regionConfigParser *config.ObjectParser

	regionGenConfig       *singleton.RegionGenConfig
	regionGenConfigParser *config.ObjectParser

	tagMiscData       *tag.TagMiscData
	tagMiscDataParser *config.ObjectParser

	taskMiscData       *taskdata.TaskMiscData
	taskMiscDataParser *config.ObjectParser

	secretTowerMiscData       *towerdata.SecretTowerMiscData
	secretTowerMiscDataParser *config.ObjectParser

	vipMiscData       *vip.VipMiscData
	vipMiscDataParser *config.ObjectParser

	resistXiongNuMisc       *xiongnu.ResistXiongNuMisc
	resistXiongNuMiscParser *config.ObjectParser

	xuanyuanMiscData       *xuanydata.XuanyuanMiscData
	xuanyuanMiscDataParser *config.ObjectParser

	zhanJiangMiscData       *zhanjiang.ZhanJiangMiscData
	zhanJiangMiscDataParser *config.ObjectParser

	zhengWuMiscData       *zhengwu.ZhengWuMiscData
	zhengWuMiscDataParser *config.ObjectParser

	zhengWuRandomData       *zhengwu.ZhengWuRandomData
	zhengWuRandomDataParser *config.ObjectParser
}

func (dAtA *ConfigDatas) GetActivityCollectionData(key uint64) *activitydata.ActivityCollectionData {
	return dAtA.activityCollectionData.Map[key]
}

func (dAtA *ConfigDatas) GetActivityCollectionDataArray() []*activitydata.ActivityCollectionData {
	return dAtA.activityCollectionData.Array
}

func (dAtA *ConfigDatas) ActivityCollectionData() *ActivityCollectionDataConfig {
	return dAtA.activityCollectionData
}

type ActivityCollectionDataConfig struct {
	Map   map[uint64]*activitydata.ActivityCollectionData
	Array []*activitydata.ActivityCollectionData

	MinKeyData *activitydata.ActivityCollectionData
	MaxKeyData *activitydata.ActivityCollectionData

	parserMap map[*activitydata.ActivityCollectionData]*config.ObjectParser
}

func (d *ActivityCollectionDataConfig) Get(key uint64) *activitydata.ActivityCollectionData {
	return d.Map[key]
}

func (d *ActivityCollectionDataConfig) Must(key uint64) *activitydata.ActivityCollectionData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetActivityShowData(key uint64) *activitydata.ActivityShowData {
	return dAtA.activityShowData.Map[key]
}

func (dAtA *ConfigDatas) GetActivityShowDataArray() []*activitydata.ActivityShowData {
	return dAtA.activityShowData.Array
}

func (dAtA *ConfigDatas) ActivityShowData() *ActivityShowDataConfig {
	return dAtA.activityShowData
}

type ActivityShowDataConfig struct {
	Map   map[uint64]*activitydata.ActivityShowData
	Array []*activitydata.ActivityShowData

	MinKeyData *activitydata.ActivityShowData
	MaxKeyData *activitydata.ActivityShowData

	parserMap map[*activitydata.ActivityShowData]*config.ObjectParser
}

func (d *ActivityShowDataConfig) Get(key uint64) *activitydata.ActivityShowData {
	return d.Map[key]
}

func (d *ActivityShowDataConfig) Must(key uint64) *activitydata.ActivityShowData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetActivityTaskListModeData(key uint64) *activitydata.ActivityTaskListModeData {
	return dAtA.activityTaskListModeData.Map[key]
}

func (dAtA *ConfigDatas) GetActivityTaskListModeDataArray() []*activitydata.ActivityTaskListModeData {
	return dAtA.activityTaskListModeData.Array
}

func (dAtA *ConfigDatas) ActivityTaskListModeData() *ActivityTaskListModeDataConfig {
	return dAtA.activityTaskListModeData
}

type ActivityTaskListModeDataConfig struct {
	Map   map[uint64]*activitydata.ActivityTaskListModeData
	Array []*activitydata.ActivityTaskListModeData

	MinKeyData *activitydata.ActivityTaskListModeData
	MaxKeyData *activitydata.ActivityTaskListModeData

	parserMap map[*activitydata.ActivityTaskListModeData]*config.ObjectParser
}

func (d *ActivityTaskListModeDataConfig) Get(key uint64) *activitydata.ActivityTaskListModeData {
	return d.Map[key]
}

func (d *ActivityTaskListModeDataConfig) Must(key uint64) *activitydata.ActivityTaskListModeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCollectionExchangeData(key uint64) *activitydata.CollectionExchangeData {
	return dAtA.collectionExchangeData.Map[key]
}

func (dAtA *ConfigDatas) GetCollectionExchangeDataArray() []*activitydata.CollectionExchangeData {
	return dAtA.collectionExchangeData.Array
}

func (dAtA *ConfigDatas) CollectionExchangeData() *CollectionExchangeDataConfig {
	return dAtA.collectionExchangeData
}

type CollectionExchangeDataConfig struct {
	Map   map[uint64]*activitydata.CollectionExchangeData
	Array []*activitydata.CollectionExchangeData

	MinKeyData *activitydata.CollectionExchangeData
	MaxKeyData *activitydata.CollectionExchangeData

	parserMap map[*activitydata.CollectionExchangeData]*config.ObjectParser
}

func (d *CollectionExchangeDataConfig) Get(key uint64) *activitydata.CollectionExchangeData {
	return d.Map[key]
}

func (d *CollectionExchangeDataConfig) Must(key uint64) *activitydata.CollectionExchangeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetJunXianLevelData(key uint64) *bai_zhan_data.JunXianLevelData {
	return dAtA.junXianLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetJunXianLevelDataArray() []*bai_zhan_data.JunXianLevelData {
	return dAtA.junXianLevelData.Array
}

func (dAtA *ConfigDatas) JunXianLevelData() *JunXianLevelDataConfig {
	return dAtA.junXianLevelData
}

type JunXianLevelDataConfig struct {
	Map   map[uint64]*bai_zhan_data.JunXianLevelData
	Array []*bai_zhan_data.JunXianLevelData

	MinKeyData *bai_zhan_data.JunXianLevelData
	MaxKeyData *bai_zhan_data.JunXianLevelData

	parserMap map[*bai_zhan_data.JunXianLevelData]*config.ObjectParser
}

func (d *JunXianLevelDataConfig) Get(key uint64) *bai_zhan_data.JunXianLevelData {
	return d.Map[key]
}

func (d *JunXianLevelDataConfig) Must(key uint64) *bai_zhan_data.JunXianLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetJunXianPrizeData(key uint64) *bai_zhan_data.JunXianPrizeData {
	return dAtA.junXianPrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetJunXianPrizeDataArray() []*bai_zhan_data.JunXianPrizeData {
	return dAtA.junXianPrizeData.Array
}

func (dAtA *ConfigDatas) JunXianPrizeData() *JunXianPrizeDataConfig {
	return dAtA.junXianPrizeData
}

type JunXianPrizeDataConfig struct {
	Map   map[uint64]*bai_zhan_data.JunXianPrizeData
	Array []*bai_zhan_data.JunXianPrizeData

	MinKeyData *bai_zhan_data.JunXianPrizeData
	MaxKeyData *bai_zhan_data.JunXianPrizeData

	parserMap map[*bai_zhan_data.JunXianPrizeData]*config.ObjectParser
}

func (d *JunXianPrizeDataConfig) Get(key uint64) *bai_zhan_data.JunXianPrizeData {
	return d.Map[key]
}

func (d *JunXianPrizeDataConfig) Must(key uint64) *bai_zhan_data.JunXianPrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetHomeNpcBaseData(key uint64) *basedata.HomeNpcBaseData {
	return dAtA.homeNpcBaseData.Map[key]
}

func (dAtA *ConfigDatas) GetHomeNpcBaseDataArray() []*basedata.HomeNpcBaseData {
	return dAtA.homeNpcBaseData.Array
}

func (dAtA *ConfigDatas) HomeNpcBaseData() *HomeNpcBaseDataConfig {
	return dAtA.homeNpcBaseData
}

type HomeNpcBaseDataConfig struct {
	Map   map[uint64]*basedata.HomeNpcBaseData
	Array []*basedata.HomeNpcBaseData

	MinKeyData *basedata.HomeNpcBaseData
	MaxKeyData *basedata.HomeNpcBaseData

	parserMap map[*basedata.HomeNpcBaseData]*config.ObjectParser
}

func (d *HomeNpcBaseDataConfig) Get(key uint64) *basedata.HomeNpcBaseData {
	return d.Map[key]
}

func (d *HomeNpcBaseDataConfig) Must(key uint64) *basedata.HomeNpcBaseData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetNpcBaseData(key uint64) *basedata.NpcBaseData {
	return dAtA.npcBaseData.Map[key]
}

func (dAtA *ConfigDatas) GetNpcBaseDataArray() []*basedata.NpcBaseData {
	return dAtA.npcBaseData.Array
}

func (dAtA *ConfigDatas) NpcBaseData() *NpcBaseDataConfig {
	return dAtA.npcBaseData
}

type NpcBaseDataConfig struct {
	Map   map[uint64]*basedata.NpcBaseData
	Array []*basedata.NpcBaseData

	MinKeyData *basedata.NpcBaseData
	MaxKeyData *basedata.NpcBaseData

	parserMap map[*basedata.NpcBaseData]*config.ObjectParser
}

func (d *NpcBaseDataConfig) Get(key uint64) *basedata.NpcBaseData {
	return d.Map[key]
}

func (d *NpcBaseDataConfig) Must(key uint64) *basedata.NpcBaseData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBlockData(key uint64) *blockdata.BlockData {
	return dAtA.blockData.Map[key]
}

func (dAtA *ConfigDatas) GetBlockDataArray() []*blockdata.BlockData {
	return dAtA.blockData.Array
}

func (dAtA *ConfigDatas) BlockData() *BlockDataConfig {
	return dAtA.blockData
}

type BlockDataConfig struct {
	Map   map[uint64]*blockdata.BlockData
	Array []*blockdata.BlockData

	MinKeyData *blockdata.BlockData
	MaxKeyData *blockdata.BlockData

	parserMap map[*blockdata.BlockData]*config.ObjectParser
}

func (d *BlockDataConfig) Get(key uint64) *blockdata.BlockData {
	return d.Map[key]
}

func (d *BlockDataConfig) Must(key uint64) *blockdata.BlockData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBodyData(key uint64) *body.BodyData {
	return dAtA.bodyData.Map[key]
}

func (dAtA *ConfigDatas) GetBodyDataArray() []*body.BodyData {
	return dAtA.bodyData.Array
}

func (dAtA *ConfigDatas) BodyData() *BodyDataConfig {
	return dAtA.bodyData
}

type BodyDataConfig struct {
	Map   map[uint64]*body.BodyData
	Array []*body.BodyData

	MinKeyData *body.BodyData
	MaxKeyData *body.BodyData

	parserMap map[*body.BodyData]*config.ObjectParser
}

func (d *BodyDataConfig) Get(key uint64) *body.BodyData {
	return d.Map[key]
}

func (d *BodyDataConfig) Must(key uint64) *body.BodyData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBufferData(key uint64) *buffer.BufferData {
	return dAtA.bufferData.Map[key]
}

func (dAtA *ConfigDatas) GetBufferDataArray() []*buffer.BufferData {
	return dAtA.bufferData.Array
}

func (dAtA *ConfigDatas) BufferData() *BufferDataConfig {
	return dAtA.bufferData
}

type BufferDataConfig struct {
	Map   map[uint64]*buffer.BufferData
	Array []*buffer.BufferData

	MinKeyData *buffer.BufferData
	MaxKeyData *buffer.BufferData

	parserMap map[*buffer.BufferData]*config.ObjectParser
}

func (d *BufferDataConfig) Get(key uint64) *buffer.BufferData {
	return d.Map[key]
}

func (d *BufferDataConfig) Must(key uint64) *buffer.BufferData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBufferTypeData(key uint64) *buffer.BufferTypeData {
	return dAtA.bufferTypeData.Map[key]
}

func (dAtA *ConfigDatas) GetBufferTypeDataArray() []*buffer.BufferTypeData {
	return dAtA.bufferTypeData.Array
}

func (dAtA *ConfigDatas) BufferTypeData() *BufferTypeDataConfig {
	return dAtA.bufferTypeData
}

type BufferTypeDataConfig struct {
	Map   map[uint64]*buffer.BufferTypeData
	Array []*buffer.BufferTypeData

	MinKeyData *buffer.BufferTypeData
	MaxKeyData *buffer.BufferTypeData

	parserMap map[*buffer.BufferTypeData]*config.ObjectParser
}

func (d *BufferTypeDataConfig) Get(key uint64) *buffer.BufferTypeData {
	return d.Map[key]
}

func (d *BufferTypeDataConfig) Must(key uint64) *buffer.BufferTypeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCaptainAbilityData(key uint64) *captain.CaptainAbilityData {
	return dAtA.captainAbilityData.Map[key]
}

func (dAtA *ConfigDatas) GetCaptainAbilityDataArray() []*captain.CaptainAbilityData {
	return dAtA.captainAbilityData.Array
}

func (dAtA *ConfigDatas) CaptainAbilityData() *CaptainAbilityDataConfig {
	return dAtA.captainAbilityData
}

type CaptainAbilityDataConfig struct {
	Map   map[uint64]*captain.CaptainAbilityData
	Array []*captain.CaptainAbilityData

	MinKeyData *captain.CaptainAbilityData
	MaxKeyData *captain.CaptainAbilityData

	parserMap map[*captain.CaptainAbilityData]*config.ObjectParser
}

func (d *CaptainAbilityDataConfig) Get(key uint64) *captain.CaptainAbilityData {
	return d.Map[key]
}

func (d *CaptainAbilityDataConfig) Must(key uint64) *captain.CaptainAbilityData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCaptainData(key uint64) *captain.CaptainData {
	return dAtA.captainData.Map[key]
}

func (dAtA *ConfigDatas) GetCaptainDataArray() []*captain.CaptainData {
	return dAtA.captainData.Array
}

func (dAtA *ConfigDatas) CaptainData() *CaptainDataConfig {
	return dAtA.captainData
}

type CaptainDataConfig struct {
	Map   map[uint64]*captain.CaptainData
	Array []*captain.CaptainData

	MinKeyData *captain.CaptainData
	MaxKeyData *captain.CaptainData

	parserMap map[*captain.CaptainData]*config.ObjectParser
}

func (d *CaptainDataConfig) Get(key uint64) *captain.CaptainData {
	return d.Map[key]
}

func (d *CaptainDataConfig) Must(key uint64) *captain.CaptainData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCaptainFriendshipData(key uint64) *captain.CaptainFriendshipData {
	return dAtA.captainFriendshipData.Map[key]
}

func (dAtA *ConfigDatas) GetCaptainFriendshipDataArray() []*captain.CaptainFriendshipData {
	return dAtA.captainFriendshipData.Array
}

func (dAtA *ConfigDatas) CaptainFriendshipData() *CaptainFriendshipDataConfig {
	return dAtA.captainFriendshipData
}

type CaptainFriendshipDataConfig struct {
	Map   map[uint64]*captain.CaptainFriendshipData
	Array []*captain.CaptainFriendshipData

	MinKeyData *captain.CaptainFriendshipData
	MaxKeyData *captain.CaptainFriendshipData

	parserMap map[*captain.CaptainFriendshipData]*config.ObjectParser
}

func (d *CaptainFriendshipDataConfig) Get(key uint64) *captain.CaptainFriendshipData {
	return d.Map[key]
}

func (d *CaptainFriendshipDataConfig) Must(key uint64) *captain.CaptainFriendshipData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCaptainLevelData(key uint64) *captain.CaptainLevelData {
	return dAtA.captainLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetCaptainLevelDataArray() []*captain.CaptainLevelData {
	return dAtA.captainLevelData.Array
}

func (dAtA *ConfigDatas) CaptainLevelData() *CaptainLevelDataConfig {
	return dAtA.captainLevelData
}

type CaptainLevelDataConfig struct {
	Map   map[uint64]*captain.CaptainLevelData
	Array []*captain.CaptainLevelData

	MinKeyData *captain.CaptainLevelData
	MaxKeyData *captain.CaptainLevelData

	parserMap map[*captain.CaptainLevelData]*config.ObjectParser
}

func (d *CaptainLevelDataConfig) Get(key uint64) *captain.CaptainLevelData {
	return d.Map[key]
}

func (d *CaptainLevelDataConfig) Must(key uint64) *captain.CaptainLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCaptainOfficialCountData(key uint64) *captain.CaptainOfficialCountData {
	return dAtA.captainOfficialCountData.Map[key]
}

func (dAtA *ConfigDatas) GetCaptainOfficialCountDataArray() []*captain.CaptainOfficialCountData {
	return dAtA.captainOfficialCountData.Array
}

func (dAtA *ConfigDatas) CaptainOfficialCountData() *CaptainOfficialCountDataConfig {
	return dAtA.captainOfficialCountData
}

type CaptainOfficialCountDataConfig struct {
	Map   map[uint64]*captain.CaptainOfficialCountData
	Array []*captain.CaptainOfficialCountData

	MinKeyData *captain.CaptainOfficialCountData
	MaxKeyData *captain.CaptainOfficialCountData

	parserMap map[*captain.CaptainOfficialCountData]*config.ObjectParser
}

func (d *CaptainOfficialCountDataConfig) Get(key uint64) *captain.CaptainOfficialCountData {
	return d.Map[key]
}

func (d *CaptainOfficialCountDataConfig) Must(key uint64) *captain.CaptainOfficialCountData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCaptainOfficialData(key uint64) *captain.CaptainOfficialData {
	return dAtA.captainOfficialData.Map[key]
}

func (dAtA *ConfigDatas) GetCaptainOfficialDataArray() []*captain.CaptainOfficialData {
	return dAtA.captainOfficialData.Array
}

func (dAtA *ConfigDatas) CaptainOfficialData() *CaptainOfficialDataConfig {
	return dAtA.captainOfficialData
}

type CaptainOfficialDataConfig struct {
	Map   map[uint64]*captain.CaptainOfficialData
	Array []*captain.CaptainOfficialData

	MinKeyData *captain.CaptainOfficialData
	MaxKeyData *captain.CaptainOfficialData

	parserMap map[*captain.CaptainOfficialData]*config.ObjectParser
}

func (d *CaptainOfficialDataConfig) Get(key uint64) *captain.CaptainOfficialData {
	return d.Map[key]
}

func (d *CaptainOfficialDataConfig) Must(key uint64) *captain.CaptainOfficialData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCaptainRarityData(key uint64) *captain.CaptainRarityData {
	return dAtA.captainRarityData.Map[key]
}

func (dAtA *ConfigDatas) GetCaptainRarityDataArray() []*captain.CaptainRarityData {
	return dAtA.captainRarityData.Array
}

func (dAtA *ConfigDatas) CaptainRarityData() *CaptainRarityDataConfig {
	return dAtA.captainRarityData
}

type CaptainRarityDataConfig struct {
	Map   map[uint64]*captain.CaptainRarityData
	Array []*captain.CaptainRarityData

	MinKeyData *captain.CaptainRarityData
	MaxKeyData *captain.CaptainRarityData

	parserMap map[*captain.CaptainRarityData]*config.ObjectParser
}

func (d *CaptainRarityDataConfig) Get(key uint64) *captain.CaptainRarityData {
	return d.Map[key]
}

func (d *CaptainRarityDataConfig) Must(key uint64) *captain.CaptainRarityData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCaptainRebirthLevelData(key uint64) *captain.CaptainRebirthLevelData {
	return dAtA.captainRebirthLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetCaptainRebirthLevelDataArray() []*captain.CaptainRebirthLevelData {
	return dAtA.captainRebirthLevelData.Array
}

func (dAtA *ConfigDatas) CaptainRebirthLevelData() *CaptainRebirthLevelDataConfig {
	return dAtA.captainRebirthLevelData
}

type CaptainRebirthLevelDataConfig struct {
	Map   map[uint64]*captain.CaptainRebirthLevelData
	Array []*captain.CaptainRebirthLevelData

	MinKeyData *captain.CaptainRebirthLevelData
	MaxKeyData *captain.CaptainRebirthLevelData

	parserMap map[*captain.CaptainRebirthLevelData]*config.ObjectParser
}

func (d *CaptainRebirthLevelDataConfig) Get(key uint64) *captain.CaptainRebirthLevelData {
	return d.Map[key]
}

func (d *CaptainRebirthLevelDataConfig) Must(key uint64) *captain.CaptainRebirthLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCaptainStarData(key uint64) *captain.CaptainStarData {
	return dAtA.captainStarData.Map[key]
}

func (dAtA *ConfigDatas) GetCaptainStarDataArray() []*captain.CaptainStarData {
	return dAtA.captainStarData.Array
}

func (dAtA *ConfigDatas) CaptainStarData() *CaptainStarDataConfig {
	return dAtA.captainStarData
}

type CaptainStarDataConfig struct {
	Map   map[uint64]*captain.CaptainStarData
	Array []*captain.CaptainStarData

	MinKeyData *captain.CaptainStarData
	MaxKeyData *captain.CaptainStarData

	parserMap map[*captain.CaptainStarData]*config.ObjectParser
}

func (d *CaptainStarDataConfig) Get(key uint64) *captain.CaptainStarData {
	return d.Map[key]
}

func (d *CaptainStarDataConfig) Must(key uint64) *captain.CaptainStarData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetNamelessCaptainData(key uint64) *captain.NamelessCaptainData {
	return dAtA.namelessCaptainData.Map[key]
}

func (dAtA *ConfigDatas) GetNamelessCaptainDataArray() []*captain.NamelessCaptainData {
	return dAtA.namelessCaptainData.Array
}

func (dAtA *ConfigDatas) NamelessCaptainData() *NamelessCaptainDataConfig {
	return dAtA.namelessCaptainData
}

type NamelessCaptainDataConfig struct {
	Map   map[uint64]*captain.NamelessCaptainData
	Array []*captain.NamelessCaptainData

	MinKeyData *captain.NamelessCaptainData
	MaxKeyData *captain.NamelessCaptainData

	parserMap map[*captain.NamelessCaptainData]*config.ObjectParser
}

func (d *NamelessCaptainDataConfig) Get(key uint64) *captain.NamelessCaptainData {
	return d.Map[key]
}

func (d *NamelessCaptainDataConfig) Must(key uint64) *captain.NamelessCaptainData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetChargeObjData(key uint64) *charge.ChargeObjData {
	return dAtA.chargeObjData.Map[key]
}

func (dAtA *ConfigDatas) GetChargeObjDataArray() []*charge.ChargeObjData {
	return dAtA.chargeObjData.Array
}

func (dAtA *ConfigDatas) ChargeObjData() *ChargeObjDataConfig {
	return dAtA.chargeObjData
}

type ChargeObjDataConfig struct {
	Map   map[uint64]*charge.ChargeObjData
	Array []*charge.ChargeObjData

	MinKeyData *charge.ChargeObjData
	MaxKeyData *charge.ChargeObjData

	parserMap map[*charge.ChargeObjData]*config.ObjectParser
}

func (d *ChargeObjDataConfig) Get(key uint64) *charge.ChargeObjData {
	return d.Map[key]
}

func (d *ChargeObjDataConfig) Must(key uint64) *charge.ChargeObjData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetChargePrizeData(key uint64) *charge.ChargePrizeData {
	return dAtA.chargePrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetChargePrizeDataArray() []*charge.ChargePrizeData {
	return dAtA.chargePrizeData.Array
}

func (dAtA *ConfigDatas) ChargePrizeData() *ChargePrizeDataConfig {
	return dAtA.chargePrizeData
}

type ChargePrizeDataConfig struct {
	Map   map[uint64]*charge.ChargePrizeData
	Array []*charge.ChargePrizeData

	MinKeyData *charge.ChargePrizeData
	MaxKeyData *charge.ChargePrizeData

	parserMap map[*charge.ChargePrizeData]*config.ObjectParser
}

func (d *ChargePrizeDataConfig) Get(key uint64) *charge.ChargePrizeData {
	return d.Map[key]
}

func (d *ChargePrizeDataConfig) Must(key uint64) *charge.ChargePrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetProductData(key uint64) *charge.ProductData {
	return dAtA.productData.Map[key]
}

func (dAtA *ConfigDatas) GetProductDataArray() []*charge.ProductData {
	return dAtA.productData.Array
}

func (dAtA *ConfigDatas) ProductData() *ProductDataConfig {
	return dAtA.productData
}

type ProductDataConfig struct {
	Map   map[uint64]*charge.ProductData
	Array []*charge.ProductData

	MinKeyData *charge.ProductData
	MaxKeyData *charge.ProductData

	parserMap map[*charge.ProductData]*config.ObjectParser
}

func (d *ProductDataConfig) Get(key uint64) *charge.ProductData {
	return d.Map[key]
}

func (d *ProductDataConfig) Must(key uint64) *charge.ProductData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetEquipCombineData(key uint64) *combine.EquipCombineData {
	return dAtA.equipCombineData.Map[key]
}

func (dAtA *ConfigDatas) GetEquipCombineDataArray() []*combine.EquipCombineData {
	return dAtA.equipCombineData.Array
}

func (dAtA *ConfigDatas) EquipCombineData() *EquipCombineDataConfig {
	return dAtA.equipCombineData
}

type EquipCombineDataConfig struct {
	Map   map[uint64]*combine.EquipCombineData
	Array []*combine.EquipCombineData

	MinKeyData *combine.EquipCombineData
	MaxKeyData *combine.EquipCombineData

	parserMap map[*combine.EquipCombineData]*config.ObjectParser
}

func (d *EquipCombineDataConfig) Get(key uint64) *combine.EquipCombineData {
	return d.Map[key]
}

func (d *EquipCombineDataConfig) Must(key uint64) *combine.EquipCombineData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGoodsCombineData(key uint64) *combine.GoodsCombineData {
	return dAtA.goodsCombineData.Map[key]
}

func (dAtA *ConfigDatas) GetGoodsCombineDataArray() []*combine.GoodsCombineData {
	return dAtA.goodsCombineData.Array
}

func (dAtA *ConfigDatas) GoodsCombineData() *GoodsCombineDataConfig {
	return dAtA.goodsCombineData
}

type GoodsCombineDataConfig struct {
	Map   map[uint64]*combine.GoodsCombineData
	Array []*combine.GoodsCombineData

	MinKeyData *combine.GoodsCombineData
	MaxKeyData *combine.GoodsCombineData

	parserMap map[*combine.GoodsCombineData]*config.ObjectParser
}

func (d *GoodsCombineDataConfig) Get(key uint64) *combine.GoodsCombineData {
	return d.Map[key]
}

func (d *GoodsCombineDataConfig) Must(key uint64) *combine.GoodsCombineData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCountryData(key uint64) *country.CountryData {
	return dAtA.countryData.Map[key]
}

func (dAtA *ConfigDatas) GetCountryDataArray() []*country.CountryData {
	return dAtA.countryData.Array
}

func (dAtA *ConfigDatas) CountryData() *CountryDataConfig {
	return dAtA.countryData
}

type CountryDataConfig struct {
	Map   map[uint64]*country.CountryData
	Array []*country.CountryData

	MinKeyData *country.CountryData
	MaxKeyData *country.CountryData

	parserMap map[*country.CountryData]*config.ObjectParser
}

func (d *CountryDataConfig) Get(key uint64) *country.CountryData {
	return d.Map[key]
}

func (d *CountryDataConfig) Must(key uint64) *country.CountryData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCountryOfficialData(key int) *country.CountryOfficialData {
	return dAtA.countryOfficialData.Map[key]
}

func (dAtA *ConfigDatas) GetCountryOfficialDataArray() []*country.CountryOfficialData {
	return dAtA.countryOfficialData.Array
}

func (dAtA *ConfigDatas) CountryOfficialData() *CountryOfficialDataConfig {
	return dAtA.countryOfficialData
}

type CountryOfficialDataConfig struct {
	Map   map[int]*country.CountryOfficialData
	Array []*country.CountryOfficialData

	MinKeyData *country.CountryOfficialData
	MaxKeyData *country.CountryOfficialData

	parserMap map[*country.CountryOfficialData]*config.ObjectParser
}

func (d *CountryOfficialDataConfig) Get(key int) *country.CountryOfficialData {
	return d.Map[key]
}

func (d *CountryOfficialDataConfig) Must(key int) *country.CountryOfficialData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCountryOfficialNpcData(key uint64) *country.CountryOfficialNpcData {
	return dAtA.countryOfficialNpcData.Map[key]
}

func (dAtA *ConfigDatas) GetCountryOfficialNpcDataArray() []*country.CountryOfficialNpcData {
	return dAtA.countryOfficialNpcData.Array
}

func (dAtA *ConfigDatas) CountryOfficialNpcData() *CountryOfficialNpcDataConfig {
	return dAtA.countryOfficialNpcData
}

type CountryOfficialNpcDataConfig struct {
	Map   map[uint64]*country.CountryOfficialNpcData
	Array []*country.CountryOfficialNpcData

	MinKeyData *country.CountryOfficialNpcData
	MaxKeyData *country.CountryOfficialNpcData

	parserMap map[*country.CountryOfficialNpcData]*config.ObjectParser
}

func (d *CountryOfficialNpcDataConfig) Get(key uint64) *country.CountryOfficialNpcData {
	return d.Map[key]
}

func (d *CountryOfficialNpcDataConfig) Must(key uint64) *country.CountryOfficialNpcData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFamilyNameData(key uint64) *country.FamilyNameData {
	return dAtA.familyNameData.Map[key]
}

func (dAtA *ConfigDatas) GetFamilyNameDataArray() []*country.FamilyNameData {
	return dAtA.familyNameData.Array
}

func (dAtA *ConfigDatas) FamilyNameData() *FamilyNameDataConfig {
	return dAtA.familyNameData
}

type FamilyNameDataConfig struct {
	Map   map[uint64]*country.FamilyNameData
	Array []*country.FamilyNameData

	MinKeyData *country.FamilyNameData
	MaxKeyData *country.FamilyNameData

	parserMap map[*country.FamilyNameData]*config.ObjectParser
}

func (d *FamilyNameDataConfig) Get(key uint64) *country.FamilyNameData {
	return d.Map[key]
}

func (d *FamilyNameDataConfig) Must(key uint64) *country.FamilyNameData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBroadcastData(key string) *data.BroadcastData {
	return dAtA.broadcastData.Map[key]
}

func (dAtA *ConfigDatas) GetBroadcastDataArray() []*data.BroadcastData {
	return dAtA.broadcastData.Array
}

func (dAtA *ConfigDatas) BroadcastData() *BroadcastDataConfig {
	return dAtA.broadcastData
}

type BroadcastDataConfig struct {
	Map   map[string]*data.BroadcastData
	Array []*data.BroadcastData

	MinKeyData *data.BroadcastData
	MaxKeyData *data.BroadcastData

	parserMap map[*data.BroadcastData]*config.ObjectParser
}

func (d *BroadcastDataConfig) Get(key string) *data.BroadcastData {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetBuffEffectData(key uint64) *data.BuffEffectData {
	return dAtA.buffEffectData.Map[key]
}

func (dAtA *ConfigDatas) GetBuffEffectDataArray() []*data.BuffEffectData {
	return dAtA.buffEffectData.Array
}

func (dAtA *ConfigDatas) BuffEffectData() *BuffEffectDataConfig {
	return dAtA.buffEffectData
}

type BuffEffectDataConfig struct {
	Map   map[uint64]*data.BuffEffectData
	Array []*data.BuffEffectData

	MinKeyData *data.BuffEffectData
	MaxKeyData *data.BuffEffectData

	parserMap map[*data.BuffEffectData]*config.ObjectParser
}

func (d *BuffEffectDataConfig) Get(key uint64) *data.BuffEffectData {
	return d.Map[key]
}

func (d *BuffEffectDataConfig) Must(key uint64) *data.BuffEffectData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetColorData(key uint64) *data.ColorData {
	return dAtA.colorData.Map[key]
}

func (dAtA *ConfigDatas) GetColorDataArray() []*data.ColorData {
	return dAtA.colorData.Array
}

func (dAtA *ConfigDatas) ColorData() *ColorDataConfig {
	return dAtA.colorData
}

type ColorDataConfig struct {
	Map   map[uint64]*data.ColorData
	Array []*data.ColorData

	MinKeyData *data.ColorData
	MaxKeyData *data.ColorData

	parserMap map[*data.ColorData]*config.ObjectParser
}

func (d *ColorDataConfig) Get(key uint64) *data.ColorData {
	return d.Map[key]
}

func (d *ColorDataConfig) Must(key uint64) *data.ColorData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFamilyName(key string) *data.FamilyName {
	return dAtA.familyName.Map[key]
}

func (dAtA *ConfigDatas) GetFamilyNameArray() []*data.FamilyName {
	return dAtA.familyName.Array
}

func (dAtA *ConfigDatas) FamilyName() *FamilyNameConfig {
	return dAtA.familyName
}

type FamilyNameConfig struct {
	Map   map[string]*data.FamilyName
	Array []*data.FamilyName

	MinKeyData *data.FamilyName
	MaxKeyData *data.FamilyName

	parserMap map[*data.FamilyName]*config.ObjectParser
}

func (d *FamilyNameConfig) Get(key string) *data.FamilyName {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetFemaleGivenName(key string) *data.FemaleGivenName {
	return dAtA.femaleGivenName.Map[key]
}

func (dAtA *ConfigDatas) GetFemaleGivenNameArray() []*data.FemaleGivenName {
	return dAtA.femaleGivenName.Array
}

func (dAtA *ConfigDatas) FemaleGivenName() *FemaleGivenNameConfig {
	return dAtA.femaleGivenName
}

type FemaleGivenNameConfig struct {
	Map   map[string]*data.FemaleGivenName
	Array []*data.FemaleGivenName

	MinKeyData *data.FemaleGivenName
	MaxKeyData *data.FemaleGivenName

	parserMap map[*data.FemaleGivenName]*config.ObjectParser
}

func (d *FemaleGivenNameConfig) Get(key string) *data.FemaleGivenName {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetHeroLevelSubData(key uint64) *data.HeroLevelSubData {
	return dAtA.heroLevelSubData.Map[key]
}

func (dAtA *ConfigDatas) GetHeroLevelSubDataArray() []*data.HeroLevelSubData {
	return dAtA.heroLevelSubData.Array
}

func (dAtA *ConfigDatas) HeroLevelSubData() *HeroLevelSubDataConfig {
	return dAtA.heroLevelSubData
}

type HeroLevelSubDataConfig struct {
	Map   map[uint64]*data.HeroLevelSubData
	Array []*data.HeroLevelSubData

	MinKeyData *data.HeroLevelSubData
	MaxKeyData *data.HeroLevelSubData

	parserMap map[*data.HeroLevelSubData]*config.ObjectParser
}

func (d *HeroLevelSubDataConfig) Get(key uint64) *data.HeroLevelSubData {
	return d.Map[key]
}

func (d *HeroLevelSubDataConfig) Must(key uint64) *data.HeroLevelSubData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMaleGivenName(key string) *data.MaleGivenName {
	return dAtA.maleGivenName.Map[key]
}

func (dAtA *ConfigDatas) GetMaleGivenNameArray() []*data.MaleGivenName {
	return dAtA.maleGivenName.Array
}

func (dAtA *ConfigDatas) MaleGivenName() *MaleGivenNameConfig {
	return dAtA.maleGivenName
}

type MaleGivenNameConfig struct {
	Map   map[string]*data.MaleGivenName
	Array []*data.MaleGivenName

	MinKeyData *data.MaleGivenName
	MaxKeyData *data.MaleGivenName

	parserMap map[*data.MaleGivenName]*config.ObjectParser
}

func (d *MaleGivenNameConfig) Get(key string) *data.MaleGivenName {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetSpriteStat(key uint64) *data.SpriteStat {
	return dAtA.spriteStat.Map[key]
}

func (dAtA *ConfigDatas) GetSpriteStatArray() []*data.SpriteStat {
	return dAtA.spriteStat.Array
}

func (dAtA *ConfigDatas) SpriteStat() *SpriteStatConfig {
	return dAtA.spriteStat
}

type SpriteStatConfig struct {
	Map   map[uint64]*data.SpriteStat
	Array []*data.SpriteStat

	MinKeyData *data.SpriteStat
	MaxKeyData *data.SpriteStat

	parserMap map[*data.SpriteStat]*config.ObjectParser
}

func (d *SpriteStatConfig) Get(key uint64) *data.SpriteStat {
	return d.Map[key]
}

func (d *SpriteStatConfig) Must(key uint64) *data.SpriteStat {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetText(key string) *data.Text {
	return dAtA.text.Map[key]
}

func (dAtA *ConfigDatas) GetTextArray() []*data.Text {
	return dAtA.text.Array
}

func (dAtA *ConfigDatas) Text() *TextConfig {
	return dAtA.text
}

type TextConfig struct {
	Map   map[string]*data.Text
	Array []*data.Text

	MinKeyData *data.Text
	MaxKeyData *data.Text

	parserMap map[*data.Text]*config.ObjectParser
}

func (d *TextConfig) Get(key string) *data.Text {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetTimeRuleData(key uint64) *data.TimeRuleData {
	return dAtA.timeRuleData.Map[key]
}

func (dAtA *ConfigDatas) GetTimeRuleDataArray() []*data.TimeRuleData {
	return dAtA.timeRuleData.Array
}

func (dAtA *ConfigDatas) TimeRuleData() *TimeRuleDataConfig {
	return dAtA.timeRuleData
}

type TimeRuleDataConfig struct {
	Map   map[uint64]*data.TimeRuleData
	Array []*data.TimeRuleData

	MinKeyData *data.TimeRuleData
	MaxKeyData *data.TimeRuleData

	parserMap map[*data.TimeRuleData]*config.ObjectParser
}

func (d *TimeRuleDataConfig) Get(key uint64) *data.TimeRuleData {
	return d.Map[key]
}

func (d *TimeRuleDataConfig) Must(key uint64) *data.TimeRuleData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBaseLevelData(key uint64) *domestic_data.BaseLevelData {
	return dAtA.baseLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetBaseLevelDataArray() []*domestic_data.BaseLevelData {
	return dAtA.baseLevelData.Array
}

func (dAtA *ConfigDatas) BaseLevelData() *BaseLevelDataConfig {
	return dAtA.baseLevelData
}

type BaseLevelDataConfig struct {
	Map   map[uint64]*domestic_data.BaseLevelData
	Array []*domestic_data.BaseLevelData

	MinKeyData *domestic_data.BaseLevelData
	MaxKeyData *domestic_data.BaseLevelData

	parserMap map[*domestic_data.BaseLevelData]*config.ObjectParser
}

func (d *BaseLevelDataConfig) Get(key uint64) *domestic_data.BaseLevelData {
	return d.Map[key]
}

func (d *BaseLevelDataConfig) Must(key uint64) *domestic_data.BaseLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBuildingData(key uint64) *domestic_data.BuildingData {
	return dAtA.buildingData.Map[key]
}

func (dAtA *ConfigDatas) GetBuildingDataArray() []*domestic_data.BuildingData {
	return dAtA.buildingData.Array
}

func (dAtA *ConfigDatas) BuildingData() *BuildingDataConfig {
	return dAtA.buildingData
}

type BuildingDataConfig struct {
	Map   map[uint64]*domestic_data.BuildingData
	Array []*domestic_data.BuildingData

	MinKeyData *domestic_data.BuildingData
	MaxKeyData *domestic_data.BuildingData

	parserMap map[*domestic_data.BuildingData]*config.ObjectParser
}

func (d *BuildingDataConfig) Get(key uint64) *domestic_data.BuildingData {
	return d.Map[key]
}

func (d *BuildingDataConfig) Must(key uint64) *domestic_data.BuildingData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBuildingLayoutData(key uint64) *domestic_data.BuildingLayoutData {
	return dAtA.buildingLayoutData.Map[key]
}

func (dAtA *ConfigDatas) GetBuildingLayoutDataArray() []*domestic_data.BuildingLayoutData {
	return dAtA.buildingLayoutData.Array
}

func (dAtA *ConfigDatas) BuildingLayoutData() *BuildingLayoutDataConfig {
	return dAtA.buildingLayoutData
}

type BuildingLayoutDataConfig struct {
	Map   map[uint64]*domestic_data.BuildingLayoutData
	Array []*domestic_data.BuildingLayoutData

	MinKeyData *domestic_data.BuildingLayoutData
	MaxKeyData *domestic_data.BuildingLayoutData

	parserMap map[*domestic_data.BuildingLayoutData]*config.ObjectParser
}

func (d *BuildingLayoutDataConfig) Get(key uint64) *domestic_data.BuildingLayoutData {
	return d.Map[key]
}

func (d *BuildingLayoutDataConfig) Must(key uint64) *domestic_data.BuildingLayoutData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBuildingUnlockData(key uint64) *domestic_data.BuildingUnlockData {
	return dAtA.buildingUnlockData.Map[key]
}

func (dAtA *ConfigDatas) GetBuildingUnlockDataArray() []*domestic_data.BuildingUnlockData {
	return dAtA.buildingUnlockData.Array
}

func (dAtA *ConfigDatas) BuildingUnlockData() *BuildingUnlockDataConfig {
	return dAtA.buildingUnlockData
}

type BuildingUnlockDataConfig struct {
	Map   map[uint64]*domestic_data.BuildingUnlockData
	Array []*domestic_data.BuildingUnlockData

	MinKeyData *domestic_data.BuildingUnlockData
	MaxKeyData *domestic_data.BuildingUnlockData

	parserMap map[*domestic_data.BuildingUnlockData]*config.ObjectParser
}

func (d *BuildingUnlockDataConfig) Get(key uint64) *domestic_data.BuildingUnlockData {
	return d.Map[key]
}

func (d *BuildingUnlockDataConfig) Must(key uint64) *domestic_data.BuildingUnlockData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCityEventData(key uint64) *domestic_data.CityEventData {
	return dAtA.cityEventData.Map[key]
}

func (dAtA *ConfigDatas) GetCityEventDataArray() []*domestic_data.CityEventData {
	return dAtA.cityEventData.Array
}

func (dAtA *ConfigDatas) CityEventData() *CityEventDataConfig {
	return dAtA.cityEventData
}

type CityEventDataConfig struct {
	Map   map[uint64]*domestic_data.CityEventData
	Array []*domestic_data.CityEventData

	MinKeyData *domestic_data.CityEventData
	MaxKeyData *domestic_data.CityEventData

	parserMap map[*domestic_data.CityEventData]*config.ObjectParser
}

func (d *CityEventDataConfig) Get(key uint64) *domestic_data.CityEventData {
	return d.Map[key]
}

func (d *CityEventDataConfig) Must(key uint64) *domestic_data.CityEventData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCityEventLevelData(key uint64) *domestic_data.CityEventLevelData {
	return dAtA.cityEventLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetCityEventLevelDataArray() []*domestic_data.CityEventLevelData {
	return dAtA.cityEventLevelData.Array
}

func (dAtA *ConfigDatas) CityEventLevelData() *CityEventLevelDataConfig {
	return dAtA.cityEventLevelData
}

type CityEventLevelDataConfig struct {
	Map   map[uint64]*domestic_data.CityEventLevelData
	Array []*domestic_data.CityEventLevelData

	MinKeyData *domestic_data.CityEventLevelData
	MaxKeyData *domestic_data.CityEventLevelData

	parserMap map[*domestic_data.CityEventLevelData]*config.ObjectParser
}

func (d *CityEventLevelDataConfig) Get(key uint64) *domestic_data.CityEventLevelData {
	return d.Map[key]
}

func (d *CityEventLevelDataConfig) Must(key uint64) *domestic_data.CityEventLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCombineCost(key int) *domestic_data.CombineCost {
	return dAtA.combineCost.Map[key]
}

func (dAtA *ConfigDatas) GetCombineCostArray() []*domestic_data.CombineCost {
	return dAtA.combineCost.Array
}

func (dAtA *ConfigDatas) CombineCost() *CombineCostConfig {
	return dAtA.combineCost
}

type CombineCostConfig struct {
	Map   map[int]*domestic_data.CombineCost
	Array []*domestic_data.CombineCost

	MinKeyData *domestic_data.CombineCost
	MaxKeyData *domestic_data.CombineCost

	parserMap map[*domestic_data.CombineCost]*config.ObjectParser
}

func (d *CombineCostConfig) Get(key int) *domestic_data.CombineCost {
	return d.Map[key]
}

func (d *CombineCostConfig) Must(key int) *domestic_data.CombineCost {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCountdownPrizeData(key uint64) *domestic_data.CountdownPrizeData {
	return dAtA.countdownPrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetCountdownPrizeDataArray() []*domestic_data.CountdownPrizeData {
	return dAtA.countdownPrizeData.Array
}

func (dAtA *ConfigDatas) CountdownPrizeData() *CountdownPrizeDataConfig {
	return dAtA.countdownPrizeData
}

type CountdownPrizeDataConfig struct {
	Map   map[uint64]*domestic_data.CountdownPrizeData
	Array []*domestic_data.CountdownPrizeData

	MinKeyData *domestic_data.CountdownPrizeData
	MaxKeyData *domestic_data.CountdownPrizeData

	parserMap map[*domestic_data.CountdownPrizeData]*config.ObjectParser
}

func (d *CountdownPrizeDataConfig) Get(key uint64) *domestic_data.CountdownPrizeData {
	return d.Map[key]
}

func (d *CountdownPrizeDataConfig) Must(key uint64) *domestic_data.CountdownPrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCountdownPrizeDescData(key uint64) *domestic_data.CountdownPrizeDescData {
	return dAtA.countdownPrizeDescData.Map[key]
}

func (dAtA *ConfigDatas) GetCountdownPrizeDescDataArray() []*domestic_data.CountdownPrizeDescData {
	return dAtA.countdownPrizeDescData.Array
}

func (dAtA *ConfigDatas) CountdownPrizeDescData() *CountdownPrizeDescDataConfig {
	return dAtA.countdownPrizeDescData
}

type CountdownPrizeDescDataConfig struct {
	Map   map[uint64]*domestic_data.CountdownPrizeDescData
	Array []*domestic_data.CountdownPrizeDescData

	MinKeyData *domestic_data.CountdownPrizeDescData
	MaxKeyData *domestic_data.CountdownPrizeDescData

	parserMap map[*domestic_data.CountdownPrizeDescData]*config.ObjectParser
}

func (d *CountdownPrizeDescDataConfig) Get(key uint64) *domestic_data.CountdownPrizeDescData {
	return d.Map[key]
}

func (d *CountdownPrizeDescDataConfig) Must(key uint64) *domestic_data.CountdownPrizeDescData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuanFuLevelData(key uint64) *domestic_data.GuanFuLevelData {
	return dAtA.guanFuLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetGuanFuLevelDataArray() []*domestic_data.GuanFuLevelData {
	return dAtA.guanFuLevelData.Array
}

func (dAtA *ConfigDatas) GuanFuLevelData() *GuanFuLevelDataConfig {
	return dAtA.guanFuLevelData
}

type GuanFuLevelDataConfig struct {
	Map   map[uint64]*domestic_data.GuanFuLevelData
	Array []*domestic_data.GuanFuLevelData

	MinKeyData *domestic_data.GuanFuLevelData
	MaxKeyData *domestic_data.GuanFuLevelData

	parserMap map[*domestic_data.GuanFuLevelData]*config.ObjectParser
}

func (d *GuanFuLevelDataConfig) Get(key uint64) *domestic_data.GuanFuLevelData {
	return d.Map[key]
}

func (d *GuanFuLevelDataConfig) Must(key uint64) *domestic_data.GuanFuLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetOuterCityBuildingData(key uint64) *domestic_data.OuterCityBuildingData {
	return dAtA.outerCityBuildingData.Map[key]
}

func (dAtA *ConfigDatas) GetOuterCityBuildingDataArray() []*domestic_data.OuterCityBuildingData {
	return dAtA.outerCityBuildingData.Array
}

func (dAtA *ConfigDatas) OuterCityBuildingData() *OuterCityBuildingDataConfig {
	return dAtA.outerCityBuildingData
}

type OuterCityBuildingDataConfig struct {
	Map   map[uint64]*domestic_data.OuterCityBuildingData
	Array []*domestic_data.OuterCityBuildingData

	MinKeyData *domestic_data.OuterCityBuildingData
	MaxKeyData *domestic_data.OuterCityBuildingData

	parserMap map[*domestic_data.OuterCityBuildingData]*config.ObjectParser
}

func (d *OuterCityBuildingDataConfig) Get(key uint64) *domestic_data.OuterCityBuildingData {
	return d.Map[key]
}

func (d *OuterCityBuildingDataConfig) Must(key uint64) *domestic_data.OuterCityBuildingData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetOuterCityBuildingDescData(key uint64) *domestic_data.OuterCityBuildingDescData {
	return dAtA.outerCityBuildingDescData.Map[key]
}

func (dAtA *ConfigDatas) GetOuterCityBuildingDescDataArray() []*domestic_data.OuterCityBuildingDescData {
	return dAtA.outerCityBuildingDescData.Array
}

func (dAtA *ConfigDatas) OuterCityBuildingDescData() *OuterCityBuildingDescDataConfig {
	return dAtA.outerCityBuildingDescData
}

type OuterCityBuildingDescDataConfig struct {
	Map   map[uint64]*domestic_data.OuterCityBuildingDescData
	Array []*domestic_data.OuterCityBuildingDescData

	MinKeyData *domestic_data.OuterCityBuildingDescData
	MaxKeyData *domestic_data.OuterCityBuildingDescData

	parserMap map[*domestic_data.OuterCityBuildingDescData]*config.ObjectParser
}

func (d *OuterCityBuildingDescDataConfig) Get(key uint64) *domestic_data.OuterCityBuildingDescData {
	return d.Map[key]
}

func (d *OuterCityBuildingDescDataConfig) Must(key uint64) *domestic_data.OuterCityBuildingDescData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetOuterCityData(key uint64) *domestic_data.OuterCityData {
	return dAtA.outerCityData.Map[key]
}

func (dAtA *ConfigDatas) GetOuterCityDataArray() []*domestic_data.OuterCityData {
	return dAtA.outerCityData.Array
}

func (dAtA *ConfigDatas) OuterCityData() *OuterCityDataConfig {
	return dAtA.outerCityData
}

type OuterCityDataConfig struct {
	Map   map[uint64]*domestic_data.OuterCityData
	Array []*domestic_data.OuterCityData

	MinKeyData *domestic_data.OuterCityData
	MaxKeyData *domestic_data.OuterCityData

	parserMap map[*domestic_data.OuterCityData]*config.ObjectParser
}

func (d *OuterCityDataConfig) Get(key uint64) *domestic_data.OuterCityData {
	return d.Map[key]
}

func (d *OuterCityDataConfig) Must(key uint64) *domestic_data.OuterCityData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetOuterCityLayoutData(key uint64) *domestic_data.OuterCityLayoutData {
	return dAtA.outerCityLayoutData.Map[key]
}

func (dAtA *ConfigDatas) GetOuterCityLayoutDataArray() []*domestic_data.OuterCityLayoutData {
	return dAtA.outerCityLayoutData.Array
}

func (dAtA *ConfigDatas) OuterCityLayoutData() *OuterCityLayoutDataConfig {
	return dAtA.outerCityLayoutData
}

type OuterCityLayoutDataConfig struct {
	Map   map[uint64]*domestic_data.OuterCityLayoutData
	Array []*domestic_data.OuterCityLayoutData

	MinKeyData *domestic_data.OuterCityLayoutData
	MaxKeyData *domestic_data.OuterCityLayoutData

	parserMap map[*domestic_data.OuterCityLayoutData]*config.ObjectParser
}

func (d *OuterCityLayoutDataConfig) Get(key uint64) *domestic_data.OuterCityLayoutData {
	return d.Map[key]
}

func (d *OuterCityLayoutDataConfig) Must(key uint64) *domestic_data.OuterCityLayoutData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetProsperityDamageBuffData(key uint64) *domestic_data.ProsperityDamageBuffData {
	return dAtA.prosperityDamageBuffData.Map[key]
}

func (dAtA *ConfigDatas) GetProsperityDamageBuffDataArray() []*domestic_data.ProsperityDamageBuffData {
	return dAtA.prosperityDamageBuffData.Array
}

func (dAtA *ConfigDatas) ProsperityDamageBuffData() *ProsperityDamageBuffDataConfig {
	return dAtA.prosperityDamageBuffData
}

type ProsperityDamageBuffDataConfig struct {
	Map   map[uint64]*domestic_data.ProsperityDamageBuffData
	Array []*domestic_data.ProsperityDamageBuffData

	MinKeyData *domestic_data.ProsperityDamageBuffData
	MaxKeyData *domestic_data.ProsperityDamageBuffData

	parserMap map[*domestic_data.ProsperityDamageBuffData]*config.ObjectParser
}

func (d *ProsperityDamageBuffDataConfig) Get(key uint64) *domestic_data.ProsperityDamageBuffData {
	return d.Map[key]
}

func (d *ProsperityDamageBuffDataConfig) Must(key uint64) *domestic_data.ProsperityDamageBuffData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetSoldierLevelData(key uint64) *domestic_data.SoldierLevelData {
	return dAtA.soldierLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetSoldierLevelDataArray() []*domestic_data.SoldierLevelData {
	return dAtA.soldierLevelData.Array
}

func (dAtA *ConfigDatas) SoldierLevelData() *SoldierLevelDataConfig {
	return dAtA.soldierLevelData
}

type SoldierLevelDataConfig struct {
	Map   map[uint64]*domestic_data.SoldierLevelData
	Array []*domestic_data.SoldierLevelData

	MinKeyData *domestic_data.SoldierLevelData
	MaxKeyData *domestic_data.SoldierLevelData

	parserMap map[*domestic_data.SoldierLevelData]*config.ObjectParser
}

func (d *SoldierLevelDataConfig) Get(key uint64) *domestic_data.SoldierLevelData {
	return d.Map[key]
}

func (d *SoldierLevelDataConfig) Must(key uint64) *domestic_data.SoldierLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTechnologyData(key uint64) *domestic_data.TechnologyData {
	return dAtA.technologyData.Map[key]
}

func (dAtA *ConfigDatas) GetTechnologyDataArray() []*domestic_data.TechnologyData {
	return dAtA.technologyData.Array
}

func (dAtA *ConfigDatas) TechnologyData() *TechnologyDataConfig {
	return dAtA.technologyData
}

type TechnologyDataConfig struct {
	Map   map[uint64]*domestic_data.TechnologyData
	Array []*domestic_data.TechnologyData

	MinKeyData *domestic_data.TechnologyData
	MaxKeyData *domestic_data.TechnologyData

	parserMap map[*domestic_data.TechnologyData]*config.ObjectParser
}

func (d *TechnologyDataConfig) Get(key uint64) *domestic_data.TechnologyData {
	return d.Map[key]
}

func (d *TechnologyDataConfig) Must(key uint64) *domestic_data.TechnologyData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTieJiangPuLevelData(key uint64) *domestic_data.TieJiangPuLevelData {
	return dAtA.tieJiangPuLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetTieJiangPuLevelDataArray() []*domestic_data.TieJiangPuLevelData {
	return dAtA.tieJiangPuLevelData.Array
}

func (dAtA *ConfigDatas) TieJiangPuLevelData() *TieJiangPuLevelDataConfig {
	return dAtA.tieJiangPuLevelData
}

type TieJiangPuLevelDataConfig struct {
	Map   map[uint64]*domestic_data.TieJiangPuLevelData
	Array []*domestic_data.TieJiangPuLevelData

	MinKeyData *domestic_data.TieJiangPuLevelData
	MaxKeyData *domestic_data.TieJiangPuLevelData

	parserMap map[*domestic_data.TieJiangPuLevelData]*config.ObjectParser
}

func (d *TieJiangPuLevelDataConfig) Get(key uint64) *domestic_data.TieJiangPuLevelData {
	return d.Map[key]
}

func (d *TieJiangPuLevelDataConfig) Must(key uint64) *domestic_data.TieJiangPuLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetWorkshopDuration(key uint64) *domestic_data.WorkshopDuration {
	return dAtA.workshopDuration.Map[key]
}

func (dAtA *ConfigDatas) GetWorkshopDurationArray() []*domestic_data.WorkshopDuration {
	return dAtA.workshopDuration.Array
}

func (dAtA *ConfigDatas) WorkshopDuration() *WorkshopDurationConfig {
	return dAtA.workshopDuration
}

type WorkshopDurationConfig struct {
	Map   map[uint64]*domestic_data.WorkshopDuration
	Array []*domestic_data.WorkshopDuration

	MinKeyData *domestic_data.WorkshopDuration
	MaxKeyData *domestic_data.WorkshopDuration

	parserMap map[*domestic_data.WorkshopDuration]*config.ObjectParser
}

func (d *WorkshopDurationConfig) Get(key uint64) *domestic_data.WorkshopDuration {
	return d.Map[key]
}

func (d *WorkshopDurationConfig) Must(key uint64) *domestic_data.WorkshopDuration {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetWorkshopLevelData(key uint64) *domestic_data.WorkshopLevelData {
	return dAtA.workshopLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetWorkshopLevelDataArray() []*domestic_data.WorkshopLevelData {
	return dAtA.workshopLevelData.Array
}

func (dAtA *ConfigDatas) WorkshopLevelData() *WorkshopLevelDataConfig {
	return dAtA.workshopLevelData
}

type WorkshopLevelDataConfig struct {
	Map   map[uint64]*domestic_data.WorkshopLevelData
	Array []*domestic_data.WorkshopLevelData

	MinKeyData *domestic_data.WorkshopLevelData
	MaxKeyData *domestic_data.WorkshopLevelData

	parserMap map[*domestic_data.WorkshopLevelData]*config.ObjectParser
}

func (d *WorkshopLevelDataConfig) Get(key uint64) *domestic_data.WorkshopLevelData {
	return d.Map[key]
}

func (d *WorkshopLevelDataConfig) Must(key uint64) *domestic_data.WorkshopLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetWorkshopRefreshCost(key uint64) *domestic_data.WorkshopRefreshCost {
	return dAtA.workshopRefreshCost.Map[key]
}

func (dAtA *ConfigDatas) GetWorkshopRefreshCostArray() []*domestic_data.WorkshopRefreshCost {
	return dAtA.workshopRefreshCost.Array
}

func (dAtA *ConfigDatas) WorkshopRefreshCost() *WorkshopRefreshCostConfig {
	return dAtA.workshopRefreshCost
}

type WorkshopRefreshCostConfig struct {
	Map   map[uint64]*domestic_data.WorkshopRefreshCost
	Array []*domestic_data.WorkshopRefreshCost

	MinKeyData *domestic_data.WorkshopRefreshCost
	MaxKeyData *domestic_data.WorkshopRefreshCost

	parserMap map[*domestic_data.WorkshopRefreshCost]*config.ObjectParser
}

func (d *WorkshopRefreshCostConfig) Get(key uint64) *domestic_data.WorkshopRefreshCost {
	return d.Map[key]
}

func (d *WorkshopRefreshCostConfig) Must(key uint64) *domestic_data.WorkshopRefreshCost {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetDungeonChapterData(key uint64) *dungeon.DungeonChapterData {
	return dAtA.dungeonChapterData.Map[key]
}

func (dAtA *ConfigDatas) GetDungeonChapterDataArray() []*dungeon.DungeonChapterData {
	return dAtA.dungeonChapterData.Array
}

func (dAtA *ConfigDatas) DungeonChapterData() *DungeonChapterDataConfig {
	return dAtA.dungeonChapterData
}

type DungeonChapterDataConfig struct {
	Map   map[uint64]*dungeon.DungeonChapterData
	Array []*dungeon.DungeonChapterData

	MinKeyData *dungeon.DungeonChapterData
	MaxKeyData *dungeon.DungeonChapterData

	parserMap map[*dungeon.DungeonChapterData]*config.ObjectParser
}

func (d *DungeonChapterDataConfig) Get(key uint64) *dungeon.DungeonChapterData {
	return d.Map[key]
}

func (d *DungeonChapterDataConfig) Must(key uint64) *dungeon.DungeonChapterData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetDungeonData(key uint64) *dungeon.DungeonData {
	return dAtA.dungeonData.Map[key]
}

func (dAtA *ConfigDatas) GetDungeonDataArray() []*dungeon.DungeonData {
	return dAtA.dungeonData.Array
}

func (dAtA *ConfigDatas) DungeonData() *DungeonDataConfig {
	return dAtA.dungeonData
}

type DungeonDataConfig struct {
	Map   map[uint64]*dungeon.DungeonData
	Array []*dungeon.DungeonData

	MinKeyData *dungeon.DungeonData
	MaxKeyData *dungeon.DungeonData

	parserMap map[*dungeon.DungeonData]*config.ObjectParser
}

func (d *DungeonDataConfig) Get(key uint64) *dungeon.DungeonData {
	return d.Map[key]
}

func (d *DungeonDataConfig) Must(key uint64) *dungeon.DungeonData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetDungeonGuideTroopData(key uint64) *dungeon.DungeonGuideTroopData {
	return dAtA.dungeonGuideTroopData.Map[key]
}

func (dAtA *ConfigDatas) GetDungeonGuideTroopDataArray() []*dungeon.DungeonGuideTroopData {
	return dAtA.dungeonGuideTroopData.Array
}

func (dAtA *ConfigDatas) DungeonGuideTroopData() *DungeonGuideTroopDataConfig {
	return dAtA.dungeonGuideTroopData
}

type DungeonGuideTroopDataConfig struct {
	Map   map[uint64]*dungeon.DungeonGuideTroopData
	Array []*dungeon.DungeonGuideTroopData

	MinKeyData *dungeon.DungeonGuideTroopData
	MaxKeyData *dungeon.DungeonGuideTroopData

	parserMap map[*dungeon.DungeonGuideTroopData]*config.ObjectParser
}

func (d *DungeonGuideTroopDataConfig) Get(key uint64) *dungeon.DungeonGuideTroopData {
	return d.Map[key]
}

func (d *DungeonGuideTroopDataConfig) Must(key uint64) *dungeon.DungeonGuideTroopData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFarmMaxStealConfig(key uint64) *farm.FarmMaxStealConfig {
	return dAtA.farmMaxStealConfig.Map[key]
}

func (dAtA *ConfigDatas) GetFarmMaxStealConfigArray() []*farm.FarmMaxStealConfig {
	return dAtA.farmMaxStealConfig.Array
}

func (dAtA *ConfigDatas) FarmMaxStealConfig() *FarmMaxStealConfigConfig {
	return dAtA.farmMaxStealConfig
}

type FarmMaxStealConfigConfig struct {
	Map   map[uint64]*farm.FarmMaxStealConfig
	Array []*farm.FarmMaxStealConfig

	MinKeyData *farm.FarmMaxStealConfig
	MaxKeyData *farm.FarmMaxStealConfig

	parserMap map[*farm.FarmMaxStealConfig]*config.ObjectParser
}

func (d *FarmMaxStealConfigConfig) Get(key uint64) *farm.FarmMaxStealConfig {
	return d.Map[key]
}

func (d *FarmMaxStealConfigConfig) Must(key uint64) *farm.FarmMaxStealConfig {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFarmOneKeyConfig(key uint64) *farm.FarmOneKeyConfig {
	return dAtA.farmOneKeyConfig.Map[key]
}

func (dAtA *ConfigDatas) GetFarmOneKeyConfigArray() []*farm.FarmOneKeyConfig {
	return dAtA.farmOneKeyConfig.Array
}

func (dAtA *ConfigDatas) FarmOneKeyConfig() *FarmOneKeyConfigConfig {
	return dAtA.farmOneKeyConfig
}

type FarmOneKeyConfigConfig struct {
	Map   map[uint64]*farm.FarmOneKeyConfig
	Array []*farm.FarmOneKeyConfig

	MinKeyData *farm.FarmOneKeyConfig
	MaxKeyData *farm.FarmOneKeyConfig

	parserMap map[*farm.FarmOneKeyConfig]*config.ObjectParser
}

func (d *FarmOneKeyConfigConfig) Get(key uint64) *farm.FarmOneKeyConfig {
	return d.Map[key]
}

func (d *FarmOneKeyConfigConfig) Must(key uint64) *farm.FarmOneKeyConfig {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFarmResConfig(key uint64) *farm.FarmResConfig {
	return dAtA.farmResConfig.Map[key]
}

func (dAtA *ConfigDatas) GetFarmResConfigArray() []*farm.FarmResConfig {
	return dAtA.farmResConfig.Array
}

func (dAtA *ConfigDatas) FarmResConfig() *FarmResConfigConfig {
	return dAtA.farmResConfig
}

type FarmResConfigConfig struct {
	Map   map[uint64]*farm.FarmResConfig
	Array []*farm.FarmResConfig

	MinKeyData *farm.FarmResConfig
	MaxKeyData *farm.FarmResConfig

	parserMap map[*farm.FarmResConfig]*config.ObjectParser
}

func (d *FarmResConfigConfig) Get(key uint64) *farm.FarmResConfig {
	return d.Map[key]
}

func (d *FarmResConfigConfig) Must(key uint64) *farm.FarmResConfig {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFishData(key uint64) *fishing_data.FishData {
	return dAtA.fishData.Map[key]
}

func (dAtA *ConfigDatas) GetFishDataArray() []*fishing_data.FishData {
	return dAtA.fishData.Array
}

func (dAtA *ConfigDatas) FishData() *FishDataConfig {
	return dAtA.fishData
}

type FishDataConfig struct {
	Map   map[uint64]*fishing_data.FishData
	Array []*fishing_data.FishData

	MinKeyData *fishing_data.FishData
	MaxKeyData *fishing_data.FishData

	parserMap map[*fishing_data.FishData]*config.ObjectParser
}

func (d *FishDataConfig) Get(key uint64) *fishing_data.FishData {
	return d.Map[key]
}

func (d *FishDataConfig) Must(key uint64) *fishing_data.FishData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFishingCaptainProbabilityData(key uint64) *fishing_data.FishingCaptainProbabilityData {
	return dAtA.fishingCaptainProbabilityData.Map[key]
}

func (dAtA *ConfigDatas) GetFishingCaptainProbabilityDataArray() []*fishing_data.FishingCaptainProbabilityData {
	return dAtA.fishingCaptainProbabilityData.Array
}

func (dAtA *ConfigDatas) FishingCaptainProbabilityData() *FishingCaptainProbabilityDataConfig {
	return dAtA.fishingCaptainProbabilityData
}

type FishingCaptainProbabilityDataConfig struct {
	Map   map[uint64]*fishing_data.FishingCaptainProbabilityData
	Array []*fishing_data.FishingCaptainProbabilityData

	MinKeyData *fishing_data.FishingCaptainProbabilityData
	MaxKeyData *fishing_data.FishingCaptainProbabilityData

	parserMap map[*fishing_data.FishingCaptainProbabilityData]*config.ObjectParser
}

func (d *FishingCaptainProbabilityDataConfig) Get(key uint64) *fishing_data.FishingCaptainProbabilityData {
	return d.Map[key]
}

func (d *FishingCaptainProbabilityDataConfig) Must(key uint64) *fishing_data.FishingCaptainProbabilityData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFishingCostData(key uint64) *fishing_data.FishingCostData {
	return dAtA.fishingCostData.Map[key]
}

func (dAtA *ConfigDatas) GetFishingCostDataArray() []*fishing_data.FishingCostData {
	return dAtA.fishingCostData.Array
}

func (dAtA *ConfigDatas) FishingCostData() *FishingCostDataConfig {
	return dAtA.fishingCostData
}

type FishingCostDataConfig struct {
	Map   map[uint64]*fishing_data.FishingCostData
	Array []*fishing_data.FishingCostData

	MinKeyData *fishing_data.FishingCostData
	MaxKeyData *fishing_data.FishingCostData

	parserMap map[*fishing_data.FishingCostData]*config.ObjectParser
}

func (d *FishingCostDataConfig) Get(key uint64) *fishing_data.FishingCostData {
	return d.Map[key]
}

func (d *FishingCostDataConfig) Must(key uint64) *fishing_data.FishingCostData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFishingShowData(key uint64) *fishing_data.FishingShowData {
	return dAtA.fishingShowData.Map[key]
}

func (dAtA *ConfigDatas) GetFishingShowDataArray() []*fishing_data.FishingShowData {
	return dAtA.fishingShowData.Array
}

func (dAtA *ConfigDatas) FishingShowData() *FishingShowDataConfig {
	return dAtA.fishingShowData
}

type FishingShowDataConfig struct {
	Map   map[uint64]*fishing_data.FishingShowData
	Array []*fishing_data.FishingShowData

	MinKeyData *fishing_data.FishingShowData
	MaxKeyData *fishing_data.FishingShowData

	parserMap map[*fishing_data.FishingShowData]*config.ObjectParser
}

func (d *FishingShowDataConfig) Get(key uint64) *fishing_data.FishingShowData {
	return d.Map[key]
}

func (d *FishingShowDataConfig) Must(key uint64) *fishing_data.FishingShowData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFunctionOpenData(key uint64) *function.FunctionOpenData {
	return dAtA.functionOpenData.Map[key]
}

func (dAtA *ConfigDatas) GetFunctionOpenDataArray() []*function.FunctionOpenData {
	return dAtA.functionOpenData.Array
}

func (dAtA *ConfigDatas) FunctionOpenData() *FunctionOpenDataConfig {
	return dAtA.functionOpenData
}

type FunctionOpenDataConfig struct {
	Map   map[uint64]*function.FunctionOpenData
	Array []*function.FunctionOpenData

	MinKeyData *function.FunctionOpenData
	MaxKeyData *function.FunctionOpenData

	parserMap map[*function.FunctionOpenData]*config.ObjectParser
}

func (d *FunctionOpenDataConfig) Get(key uint64) *function.FunctionOpenData {
	return d.Map[key]
}

func (d *FunctionOpenDataConfig) Must(key uint64) *function.FunctionOpenData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTreasuryTreeData(key uint64) *gardendata.TreasuryTreeData {
	return dAtA.treasuryTreeData.Map[key]
}

func (dAtA *ConfigDatas) GetTreasuryTreeDataArray() []*gardendata.TreasuryTreeData {
	return dAtA.treasuryTreeData.Array
}

func (dAtA *ConfigDatas) TreasuryTreeData() *TreasuryTreeDataConfig {
	return dAtA.treasuryTreeData
}

type TreasuryTreeDataConfig struct {
	Map   map[uint64]*gardendata.TreasuryTreeData
	Array []*gardendata.TreasuryTreeData

	MinKeyData *gardendata.TreasuryTreeData
	MaxKeyData *gardendata.TreasuryTreeData

	parserMap map[*gardendata.TreasuryTreeData]*config.ObjectParser
}

func (d *TreasuryTreeDataConfig) Get(key uint64) *gardendata.TreasuryTreeData {
	return d.Map[key]
}

func (d *TreasuryTreeDataConfig) Must(key uint64) *gardendata.TreasuryTreeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetEquipmentData(key uint64) *goods.EquipmentData {
	return dAtA.equipmentData.Map[key]
}

func (dAtA *ConfigDatas) GetEquipmentDataArray() []*goods.EquipmentData {
	return dAtA.equipmentData.Array
}

func (dAtA *ConfigDatas) EquipmentData() *EquipmentDataConfig {
	return dAtA.equipmentData
}

type EquipmentDataConfig struct {
	Map   map[uint64]*goods.EquipmentData
	Array []*goods.EquipmentData

	MinKeyData *goods.EquipmentData
	MaxKeyData *goods.EquipmentData

	parserMap map[*goods.EquipmentData]*config.ObjectParser
}

func (d *EquipmentDataConfig) Get(key uint64) *goods.EquipmentData {
	return d.Map[key]
}

func (d *EquipmentDataConfig) Must(key uint64) *goods.EquipmentData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetEquipmentLevelData(key uint64) *goods.EquipmentLevelData {
	return dAtA.equipmentLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetEquipmentLevelDataArray() []*goods.EquipmentLevelData {
	return dAtA.equipmentLevelData.Array
}

func (dAtA *ConfigDatas) EquipmentLevelData() *EquipmentLevelDataConfig {
	return dAtA.equipmentLevelData
}

type EquipmentLevelDataConfig struct {
	Map   map[uint64]*goods.EquipmentLevelData
	Array []*goods.EquipmentLevelData

	MinKeyData *goods.EquipmentLevelData
	MaxKeyData *goods.EquipmentLevelData

	parserMap map[*goods.EquipmentLevelData]*config.ObjectParser
}

func (d *EquipmentLevelDataConfig) Get(key uint64) *goods.EquipmentLevelData {
	return d.Map[key]
}

func (d *EquipmentLevelDataConfig) Must(key uint64) *goods.EquipmentLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetEquipmentQualityData(key uint64) *goods.EquipmentQualityData {
	return dAtA.equipmentQualityData.Map[key]
}

func (dAtA *ConfigDatas) GetEquipmentQualityDataArray() []*goods.EquipmentQualityData {
	return dAtA.equipmentQualityData.Array
}

func (dAtA *ConfigDatas) EquipmentQualityData() *EquipmentQualityDataConfig {
	return dAtA.equipmentQualityData
}

type EquipmentQualityDataConfig struct {
	Map   map[uint64]*goods.EquipmentQualityData
	Array []*goods.EquipmentQualityData

	MinKeyData *goods.EquipmentQualityData
	MaxKeyData *goods.EquipmentQualityData

	parserMap map[*goods.EquipmentQualityData]*config.ObjectParser
}

func (d *EquipmentQualityDataConfig) Get(key uint64) *goods.EquipmentQualityData {
	return d.Map[key]
}

func (d *EquipmentQualityDataConfig) Must(key uint64) *goods.EquipmentQualityData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetEquipmentRefinedData(key uint64) *goods.EquipmentRefinedData {
	return dAtA.equipmentRefinedData.Map[key]
}

func (dAtA *ConfigDatas) GetEquipmentRefinedDataArray() []*goods.EquipmentRefinedData {
	return dAtA.equipmentRefinedData.Array
}

func (dAtA *ConfigDatas) EquipmentRefinedData() *EquipmentRefinedDataConfig {
	return dAtA.equipmentRefinedData
}

type EquipmentRefinedDataConfig struct {
	Map   map[uint64]*goods.EquipmentRefinedData
	Array []*goods.EquipmentRefinedData

	MinKeyData *goods.EquipmentRefinedData
	MaxKeyData *goods.EquipmentRefinedData

	parserMap map[*goods.EquipmentRefinedData]*config.ObjectParser
}

func (d *EquipmentRefinedDataConfig) Get(key uint64) *goods.EquipmentRefinedData {
	return d.Map[key]
}

func (d *EquipmentRefinedDataConfig) Must(key uint64) *goods.EquipmentRefinedData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetEquipmentTaozData(key uint64) *goods.EquipmentTaozData {
	return dAtA.equipmentTaozData.Map[key]
}

func (dAtA *ConfigDatas) GetEquipmentTaozDataArray() []*goods.EquipmentTaozData {
	return dAtA.equipmentTaozData.Array
}

func (dAtA *ConfigDatas) EquipmentTaozData() *EquipmentTaozDataConfig {
	return dAtA.equipmentTaozData
}

type EquipmentTaozDataConfig struct {
	Map   map[uint64]*goods.EquipmentTaozData
	Array []*goods.EquipmentTaozData

	MinKeyData *goods.EquipmentTaozData
	MaxKeyData *goods.EquipmentTaozData

	parserMap map[*goods.EquipmentTaozData]*config.ObjectParser
}

func (d *EquipmentTaozDataConfig) Get(key uint64) *goods.EquipmentTaozData {
	return d.Map[key]
}

func (d *EquipmentTaozDataConfig) Must(key uint64) *goods.EquipmentTaozData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGemData(key uint64) *goods.GemData {
	return dAtA.gemData.Map[key]
}

func (dAtA *ConfigDatas) GetGemDataArray() []*goods.GemData {
	return dAtA.gemData.Array
}

func (dAtA *ConfigDatas) GemData() *GemDataConfig {
	return dAtA.gemData
}

type GemDataConfig struct {
	Map   map[uint64]*goods.GemData
	Array []*goods.GemData

	MinKeyData *goods.GemData
	MaxKeyData *goods.GemData

	parserMap map[*goods.GemData]*config.ObjectParser
}

func (d *GemDataConfig) Get(key uint64) *goods.GemData {
	return d.Map[key]
}

func (d *GemDataConfig) Must(key uint64) *goods.GemData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGoodsData(key uint64) *goods.GoodsData {
	return dAtA.goodsData.Map[key]
}

func (dAtA *ConfigDatas) GetGoodsDataArray() []*goods.GoodsData {
	return dAtA.goodsData.Array
}

func (dAtA *ConfigDatas) GoodsData() *GoodsDataConfig {
	return dAtA.goodsData
}

type GoodsDataConfig struct {
	Map   map[uint64]*goods.GoodsData
	Array []*goods.GoodsData

	MinKeyData *goods.GoodsData
	MaxKeyData *goods.GoodsData

	parserMap map[*goods.GoodsData]*config.ObjectParser
}

func (d *GoodsDataConfig) Get(key uint64) *goods.GoodsData {
	return d.Map[key]
}

func (d *GoodsDataConfig) Must(key uint64) *goods.GoodsData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGoodsQuality(key uint64) *goods.GoodsQuality {
	return dAtA.goodsQuality.Map[key]
}

func (dAtA *ConfigDatas) GetGoodsQualityArray() []*goods.GoodsQuality {
	return dAtA.goodsQuality.Array
}

func (dAtA *ConfigDatas) GoodsQuality() *GoodsQualityConfig {
	return dAtA.goodsQuality
}

type GoodsQualityConfig struct {
	Map   map[uint64]*goods.GoodsQuality
	Array []*goods.GoodsQuality

	MinKeyData *goods.GoodsQuality
	MaxKeyData *goods.GoodsQuality

	parserMap map[*goods.GoodsQuality]*config.ObjectParser
}

func (d *GoodsQualityConfig) Get(key uint64) *goods.GoodsQuality {
	return d.Map[key]
}

func (d *GoodsQualityConfig) Must(key uint64) *goods.GoodsQuality {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildBigBoxData(key uint64) *guild_data.GuildBigBoxData {
	return dAtA.guildBigBoxData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildBigBoxDataArray() []*guild_data.GuildBigBoxData {
	return dAtA.guildBigBoxData.Array
}

func (dAtA *ConfigDatas) GuildBigBoxData() *GuildBigBoxDataConfig {
	return dAtA.guildBigBoxData
}

type GuildBigBoxDataConfig struct {
	Map   map[uint64]*guild_data.GuildBigBoxData
	Array []*guild_data.GuildBigBoxData

	MinKeyData *guild_data.GuildBigBoxData
	MaxKeyData *guild_data.GuildBigBoxData

	parserMap map[*guild_data.GuildBigBoxData]*config.ObjectParser
}

func (d *GuildBigBoxDataConfig) Get(key uint64) *guild_data.GuildBigBoxData {
	return d.Map[key]
}

func (d *GuildBigBoxDataConfig) Must(key uint64) *guild_data.GuildBigBoxData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildClassLevelData(key uint64) *guild_data.GuildClassLevelData {
	return dAtA.guildClassLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildClassLevelDataArray() []*guild_data.GuildClassLevelData {
	return dAtA.guildClassLevelData.Array
}

func (dAtA *ConfigDatas) GuildClassLevelData() *GuildClassLevelDataConfig {
	return dAtA.guildClassLevelData
}

type GuildClassLevelDataConfig struct {
	Map   map[uint64]*guild_data.GuildClassLevelData
	Array []*guild_data.GuildClassLevelData

	MinKeyData *guild_data.GuildClassLevelData
	MaxKeyData *guild_data.GuildClassLevelData

	parserMap map[*guild_data.GuildClassLevelData]*config.ObjectParser
}

func (d *GuildClassLevelDataConfig) Get(key uint64) *guild_data.GuildClassLevelData {
	return d.Map[key]
}

func (d *GuildClassLevelDataConfig) Must(key uint64) *guild_data.GuildClassLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildClassTitleData(key uint64) *guild_data.GuildClassTitleData {
	return dAtA.guildClassTitleData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildClassTitleDataArray() []*guild_data.GuildClassTitleData {
	return dAtA.guildClassTitleData.Array
}

func (dAtA *ConfigDatas) GuildClassTitleData() *GuildClassTitleDataConfig {
	return dAtA.guildClassTitleData
}

type GuildClassTitleDataConfig struct {
	Map   map[uint64]*guild_data.GuildClassTitleData
	Array []*guild_data.GuildClassTitleData

	MinKeyData *guild_data.GuildClassTitleData
	MaxKeyData *guild_data.GuildClassTitleData

	parserMap map[*guild_data.GuildClassTitleData]*config.ObjectParser
}

func (d *GuildClassTitleDataConfig) Get(key uint64) *guild_data.GuildClassTitleData {
	return d.Map[key]
}

func (d *GuildClassTitleDataConfig) Must(key uint64) *guild_data.GuildClassTitleData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildDonateData(key uint64) *guild_data.GuildDonateData {
	return dAtA.guildDonateData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildDonateDataArray() []*guild_data.GuildDonateData {
	return dAtA.guildDonateData.Array
}

func (dAtA *ConfigDatas) GuildDonateData() *GuildDonateDataConfig {
	return dAtA.guildDonateData
}

type GuildDonateDataConfig struct {
	Map   map[uint64]*guild_data.GuildDonateData
	Array []*guild_data.GuildDonateData

	MinKeyData *guild_data.GuildDonateData
	MaxKeyData *guild_data.GuildDonateData

	parserMap map[*guild_data.GuildDonateData]*config.ObjectParser
}

func (d *GuildDonateDataConfig) Get(key uint64) *guild_data.GuildDonateData {
	return d.Map[key]
}

func (d *GuildDonateDataConfig) Must(key uint64) *guild_data.GuildDonateData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildEventPrizeData(key uint64) *guild_data.GuildEventPrizeData {
	return dAtA.guildEventPrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildEventPrizeDataArray() []*guild_data.GuildEventPrizeData {
	return dAtA.guildEventPrizeData.Array
}

func (dAtA *ConfigDatas) GuildEventPrizeData() *GuildEventPrizeDataConfig {
	return dAtA.guildEventPrizeData
}

type GuildEventPrizeDataConfig struct {
	Map   map[uint64]*guild_data.GuildEventPrizeData
	Array []*guild_data.GuildEventPrizeData

	MinKeyData *guild_data.GuildEventPrizeData
	MaxKeyData *guild_data.GuildEventPrizeData

	parserMap map[*guild_data.GuildEventPrizeData]*config.ObjectParser
}

func (d *GuildEventPrizeDataConfig) Get(key uint64) *guild_data.GuildEventPrizeData {
	return d.Map[key]
}

func (d *GuildEventPrizeDataConfig) Must(key uint64) *guild_data.GuildEventPrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildLevelCdrData(key uint64) *guild_data.GuildLevelCdrData {
	return dAtA.guildLevelCdrData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildLevelCdrDataArray() []*guild_data.GuildLevelCdrData {
	return dAtA.guildLevelCdrData.Array
}

func (dAtA *ConfigDatas) GuildLevelCdrData() *GuildLevelCdrDataConfig {
	return dAtA.guildLevelCdrData
}

type GuildLevelCdrDataConfig struct {
	Map   map[uint64]*guild_data.GuildLevelCdrData
	Array []*guild_data.GuildLevelCdrData

	MinKeyData *guild_data.GuildLevelCdrData
	MaxKeyData *guild_data.GuildLevelCdrData

	parserMap map[*guild_data.GuildLevelCdrData]*config.ObjectParser
}

func (d *GuildLevelCdrDataConfig) Get(key uint64) *guild_data.GuildLevelCdrData {
	return d.Map[key]
}

func (d *GuildLevelCdrDataConfig) Must(key uint64) *guild_data.GuildLevelCdrData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildLevelData(key uint64) *guild_data.GuildLevelData {
	return dAtA.guildLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildLevelDataArray() []*guild_data.GuildLevelData {
	return dAtA.guildLevelData.Array
}

func (dAtA *ConfigDatas) GuildLevelData() *GuildLevelDataConfig {
	return dAtA.guildLevelData
}

type GuildLevelDataConfig struct {
	Map   map[uint64]*guild_data.GuildLevelData
	Array []*guild_data.GuildLevelData

	MinKeyData *guild_data.GuildLevelData
	MaxKeyData *guild_data.GuildLevelData

	parserMap map[*guild_data.GuildLevelData]*config.ObjectParser
}

func (d *GuildLevelDataConfig) Get(key uint64) *guild_data.GuildLevelData {
	return d.Map[key]
}

func (d *GuildLevelDataConfig) Must(key uint64) *guild_data.GuildLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildLogData(key string) *guild_data.GuildLogData {
	return dAtA.guildLogData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildLogDataArray() []*guild_data.GuildLogData {
	return dAtA.guildLogData.Array
}

func (dAtA *ConfigDatas) GuildLogData() *GuildLogDataConfig {
	return dAtA.guildLogData
}

type GuildLogDataConfig struct {
	Map   map[string]*guild_data.GuildLogData
	Array []*guild_data.GuildLogData

	MinKeyData *guild_data.GuildLogData
	MaxKeyData *guild_data.GuildLogData

	parserMap map[*guild_data.GuildLogData]*config.ObjectParser
}

func (d *GuildLogDataConfig) Get(key string) *guild_data.GuildLogData {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetGuildPermissionShowData(key uint64) *guild_data.GuildPermissionShowData {
	return dAtA.guildPermissionShowData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildPermissionShowDataArray() []*guild_data.GuildPermissionShowData {
	return dAtA.guildPermissionShowData.Array
}

func (dAtA *ConfigDatas) GuildPermissionShowData() *GuildPermissionShowDataConfig {
	return dAtA.guildPermissionShowData
}

type GuildPermissionShowDataConfig struct {
	Map   map[uint64]*guild_data.GuildPermissionShowData
	Array []*guild_data.GuildPermissionShowData

	MinKeyData *guild_data.GuildPermissionShowData
	MaxKeyData *guild_data.GuildPermissionShowData

	parserMap map[*guild_data.GuildPermissionShowData]*config.ObjectParser
}

func (d *GuildPermissionShowDataConfig) Get(key uint64) *guild_data.GuildPermissionShowData {
	return d.Map[key]
}

func (d *GuildPermissionShowDataConfig) Must(key uint64) *guild_data.GuildPermissionShowData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildPrestigeEventData(key uint64) *guild_data.GuildPrestigeEventData {
	return dAtA.guildPrestigeEventData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildPrestigeEventDataArray() []*guild_data.GuildPrestigeEventData {
	return dAtA.guildPrestigeEventData.Array
}

func (dAtA *ConfigDatas) GuildPrestigeEventData() *GuildPrestigeEventDataConfig {
	return dAtA.guildPrestigeEventData
}

type GuildPrestigeEventDataConfig struct {
	Map   map[uint64]*guild_data.GuildPrestigeEventData
	Array []*guild_data.GuildPrestigeEventData

	MinKeyData *guild_data.GuildPrestigeEventData
	MaxKeyData *guild_data.GuildPrestigeEventData

	parserMap map[*guild_data.GuildPrestigeEventData]*config.ObjectParser
}

func (d *GuildPrestigeEventDataConfig) Get(key uint64) *guild_data.GuildPrestigeEventData {
	return d.Map[key]
}

func (d *GuildPrestigeEventDataConfig) Must(key uint64) *guild_data.GuildPrestigeEventData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildPrestigePrizeData(key uint64) *guild_data.GuildPrestigePrizeData {
	return dAtA.guildPrestigePrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildPrestigePrizeDataArray() []*guild_data.GuildPrestigePrizeData {
	return dAtA.guildPrestigePrizeData.Array
}

func (dAtA *ConfigDatas) GuildPrestigePrizeData() *GuildPrestigePrizeDataConfig {
	return dAtA.guildPrestigePrizeData
}

type GuildPrestigePrizeDataConfig struct {
	Map   map[uint64]*guild_data.GuildPrestigePrizeData
	Array []*guild_data.GuildPrestigePrizeData

	MinKeyData *guild_data.GuildPrestigePrizeData
	MaxKeyData *guild_data.GuildPrestigePrizeData

	parserMap map[*guild_data.GuildPrestigePrizeData]*config.ObjectParser
}

func (d *GuildPrestigePrizeDataConfig) Get(key uint64) *guild_data.GuildPrestigePrizeData {
	return d.Map[key]
}

func (d *GuildPrestigePrizeDataConfig) Must(key uint64) *guild_data.GuildPrestigePrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildRankPrizeData(key uint64) *guild_data.GuildRankPrizeData {
	return dAtA.guildRankPrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildRankPrizeDataArray() []*guild_data.GuildRankPrizeData {
	return dAtA.guildRankPrizeData.Array
}

func (dAtA *ConfigDatas) GuildRankPrizeData() *GuildRankPrizeDataConfig {
	return dAtA.guildRankPrizeData
}

type GuildRankPrizeDataConfig struct {
	Map   map[uint64]*guild_data.GuildRankPrizeData
	Array []*guild_data.GuildRankPrizeData

	MinKeyData *guild_data.GuildRankPrizeData
	MaxKeyData *guild_data.GuildRankPrizeData

	parserMap map[*guild_data.GuildRankPrizeData]*config.ObjectParser
}

func (d *GuildRankPrizeDataConfig) Get(key uint64) *guild_data.GuildRankPrizeData {
	return d.Map[key]
}

func (d *GuildRankPrizeDataConfig) Must(key uint64) *guild_data.GuildRankPrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildTarget(key uint64) *guild_data.GuildTarget {
	return dAtA.guildTarget.Map[key]
}

func (dAtA *ConfigDatas) GetGuildTargetArray() []*guild_data.GuildTarget {
	return dAtA.guildTarget.Array
}

func (dAtA *ConfigDatas) GuildTarget() *GuildTargetConfig {
	return dAtA.guildTarget
}

type GuildTargetConfig struct {
	Map   map[uint64]*guild_data.GuildTarget
	Array []*guild_data.GuildTarget

	MinKeyData *guild_data.GuildTarget
	MaxKeyData *guild_data.GuildTarget

	parserMap map[*guild_data.GuildTarget]*config.ObjectParser
}

func (d *GuildTargetConfig) Get(key uint64) *guild_data.GuildTarget {
	return d.Map[key]
}

func (d *GuildTargetConfig) Must(key uint64) *guild_data.GuildTarget {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildTaskData(key uint64) *guild_data.GuildTaskData {
	return dAtA.guildTaskData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildTaskDataArray() []*guild_data.GuildTaskData {
	return dAtA.guildTaskData.Array
}

func (dAtA *ConfigDatas) GuildTaskData() *GuildTaskDataConfig {
	return dAtA.guildTaskData
}

type GuildTaskDataConfig struct {
	Map   map[uint64]*guild_data.GuildTaskData
	Array []*guild_data.GuildTaskData

	MinKeyData *guild_data.GuildTaskData
	MaxKeyData *guild_data.GuildTaskData

	parserMap map[*guild_data.GuildTaskData]*config.ObjectParser
}

func (d *GuildTaskDataConfig) Get(key uint64) *guild_data.GuildTaskData {
	return d.Map[key]
}

func (d *GuildTaskDataConfig) Must(key uint64) *guild_data.GuildTaskData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildTaskEvaluateData(key uint64) *guild_data.GuildTaskEvaluateData {
	return dAtA.guildTaskEvaluateData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildTaskEvaluateDataArray() []*guild_data.GuildTaskEvaluateData {
	return dAtA.guildTaskEvaluateData.Array
}

func (dAtA *ConfigDatas) GuildTaskEvaluateData() *GuildTaskEvaluateDataConfig {
	return dAtA.guildTaskEvaluateData
}

type GuildTaskEvaluateDataConfig struct {
	Map   map[uint64]*guild_data.GuildTaskEvaluateData
	Array []*guild_data.GuildTaskEvaluateData

	MinKeyData *guild_data.GuildTaskEvaluateData
	MaxKeyData *guild_data.GuildTaskEvaluateData

	parserMap map[*guild_data.GuildTaskEvaluateData]*config.ObjectParser
}

func (d *GuildTaskEvaluateDataConfig) Get(key uint64) *guild_data.GuildTaskEvaluateData {
	return d.Map[key]
}

func (d *GuildTaskEvaluateDataConfig) Must(key uint64) *guild_data.GuildTaskEvaluateData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildTechnologyData(key uint64) *guild_data.GuildTechnologyData {
	return dAtA.guildTechnologyData.Map[key]
}

func (dAtA *ConfigDatas) GetGuildTechnologyDataArray() []*guild_data.GuildTechnologyData {
	return dAtA.guildTechnologyData.Array
}

func (dAtA *ConfigDatas) GuildTechnologyData() *GuildTechnologyDataConfig {
	return dAtA.guildTechnologyData
}

type GuildTechnologyDataConfig struct {
	Map   map[uint64]*guild_data.GuildTechnologyData
	Array []*guild_data.GuildTechnologyData

	MinKeyData *guild_data.GuildTechnologyData
	MaxKeyData *guild_data.GuildTechnologyData

	parserMap map[*guild_data.GuildTechnologyData]*config.ObjectParser
}

func (d *GuildTechnologyDataConfig) Get(key uint64) *guild_data.GuildTechnologyData {
	return d.Map[key]
}

func (d *GuildTechnologyDataConfig) Must(key uint64) *guild_data.GuildTechnologyData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetNpcGuildTemplate(key uint64) *guild_data.NpcGuildTemplate {
	return dAtA.npcGuildTemplate.Map[key]
}

func (dAtA *ConfigDatas) GetNpcGuildTemplateArray() []*guild_data.NpcGuildTemplate {
	return dAtA.npcGuildTemplate.Array
}

func (dAtA *ConfigDatas) NpcGuildTemplate() *NpcGuildTemplateConfig {
	return dAtA.npcGuildTemplate
}

type NpcGuildTemplateConfig struct {
	Map   map[uint64]*guild_data.NpcGuildTemplate
	Array []*guild_data.NpcGuildTemplate

	MinKeyData *guild_data.NpcGuildTemplate
	MaxKeyData *guild_data.NpcGuildTemplate

	parserMap map[*guild_data.NpcGuildTemplate]*config.ObjectParser
}

func (d *NpcGuildTemplateConfig) Get(key uint64) *guild_data.NpcGuildTemplate {
	return d.Map[key]
}

func (d *NpcGuildTemplateConfig) Must(key uint64) *guild_data.NpcGuildTemplate {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetNpcMemberData(key uint64) *guild_data.NpcMemberData {
	return dAtA.npcMemberData.Map[key]
}

func (dAtA *ConfigDatas) GetNpcMemberDataArray() []*guild_data.NpcMemberData {
	return dAtA.npcMemberData.Array
}

func (dAtA *ConfigDatas) NpcMemberData() *NpcMemberDataConfig {
	return dAtA.npcMemberData
}

type NpcMemberDataConfig struct {
	Map   map[uint64]*guild_data.NpcMemberData
	Array []*guild_data.NpcMemberData

	MinKeyData *guild_data.NpcMemberData
	MaxKeyData *guild_data.NpcMemberData

	parserMap map[*guild_data.NpcMemberData]*config.ObjectParser
}

func (d *NpcMemberDataConfig) Get(key uint64) *guild_data.NpcMemberData {
	return d.Map[key]
}

func (d *NpcMemberDataConfig) Must(key uint64) *guild_data.NpcMemberData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetHeadData(key string) *head.HeadData {
	return dAtA.headData.Map[key]
}

func (dAtA *ConfigDatas) GetHeadDataArray() []*head.HeadData {
	return dAtA.headData.Array
}

func (dAtA *ConfigDatas) HeadData() *HeadDataConfig {
	return dAtA.headData
}

type HeadDataConfig struct {
	Map   map[string]*head.HeadData
	Array []*head.HeadData

	MinKeyData *head.HeadData
	MaxKeyData *head.HeadData

	parserMap map[*head.HeadData]*config.ObjectParser
}

func (d *HeadDataConfig) Get(key string) *head.HeadData {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetHebiPrizeData(key uint64) *hebi.HebiPrizeData {
	return dAtA.hebiPrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetHebiPrizeDataArray() []*hebi.HebiPrizeData {
	return dAtA.hebiPrizeData.Array
}

func (dAtA *ConfigDatas) HebiPrizeData() *HebiPrizeDataConfig {
	return dAtA.hebiPrizeData
}

type HebiPrizeDataConfig struct {
	Map   map[uint64]*hebi.HebiPrizeData
	Array []*hebi.HebiPrizeData

	MinKeyData *hebi.HebiPrizeData
	MaxKeyData *hebi.HebiPrizeData

	parserMap map[*hebi.HebiPrizeData]*config.ObjectParser
}

func (d *HebiPrizeDataConfig) Get(key uint64) *hebi.HebiPrizeData {
	return d.Map[key]
}

func (d *HebiPrizeDataConfig) Must(key uint64) *hebi.HebiPrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetHeroLevelData(key uint64) *herodata.HeroLevelData {
	return dAtA.heroLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetHeroLevelDataArray() []*herodata.HeroLevelData {
	return dAtA.heroLevelData.Array
}

func (dAtA *ConfigDatas) HeroLevelData() *HeroLevelDataConfig {
	return dAtA.heroLevelData
}

type HeroLevelDataConfig struct {
	Map   map[uint64]*herodata.HeroLevelData
	Array []*herodata.HeroLevelData

	MinKeyData *herodata.HeroLevelData
	MaxKeyData *herodata.HeroLevelData

	parserMap map[*herodata.HeroLevelData]*config.ObjectParser
}

func (d *HeroLevelDataConfig) Get(key uint64) *herodata.HeroLevelData {
	return d.Map[key]
}

func (d *HeroLevelDataConfig) Must(key uint64) *herodata.HeroLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetI18nData(key string) *i18n.I18nData {
	return dAtA.i18nData.Map[key]
}

func (dAtA *ConfigDatas) GetI18nDataArray() []*i18n.I18nData {
	return dAtA.i18nData.Array
}

func (dAtA *ConfigDatas) I18nData() *I18nDataConfig {
	return dAtA.i18nData
}

type I18nDataConfig struct {
	Map   map[string]*i18n.I18nData
	Array []*i18n.I18nData

	MinKeyData *i18n.I18nData
	MaxKeyData *i18n.I18nData

	parserMap map[*i18n.I18nData]*config.ObjectParser
}

func (d *I18nDataConfig) Get(key string) *i18n.I18nData {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetIcon(key string) *icon.Icon {
	return dAtA.icon.Map[key]
}

func (dAtA *ConfigDatas) GetIconArray() []*icon.Icon {
	return dAtA.icon.Array
}

func (dAtA *ConfigDatas) Icon() *IconConfig {
	return dAtA.icon
}

type IconConfig struct {
	Map   map[string]*icon.Icon
	Array []*icon.Icon

	MinKeyData *icon.Icon
	MaxKeyData *icon.Icon

	parserMap map[*icon.Icon]*config.ObjectParser
}

func (d *IconConfig) Get(key string) *icon.Icon {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetLocationData(key uint64) *location.LocationData {
	return dAtA.locationData.Map[key]
}

func (dAtA *ConfigDatas) GetLocationDataArray() []*location.LocationData {
	return dAtA.locationData.Array
}

func (dAtA *ConfigDatas) LocationData() *LocationDataConfig {
	return dAtA.locationData
}

type LocationDataConfig struct {
	Map   map[uint64]*location.LocationData
	Array []*location.LocationData

	MinKeyData *location.LocationData
	MaxKeyData *location.LocationData

	parserMap map[*location.LocationData]*config.ObjectParser
}

func (d *LocationDataConfig) Get(key uint64) *location.LocationData {
	return d.Map[key]
}

func (d *LocationDataConfig) Must(key uint64) *location.LocationData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMailData(key string) *maildata.MailData {
	return dAtA.mailData.Map[key]
}

func (dAtA *ConfigDatas) GetMailDataArray() []*maildata.MailData {
	return dAtA.mailData.Array
}

func (dAtA *ConfigDatas) MailData() *MailDataConfig {
	return dAtA.mailData
}

type MailDataConfig struct {
	Map   map[string]*maildata.MailData
	Array []*maildata.MailData

	MinKeyData *maildata.MailData
	MaxKeyData *maildata.MailData

	parserMap map[*maildata.MailData]*config.ObjectParser
}

func (d *MailDataConfig) Get(key string) *maildata.MailData {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetJiuGuanData(key uint64) *military_data.JiuGuanData {
	return dAtA.jiuGuanData.Map[key]
}

func (dAtA *ConfigDatas) GetJiuGuanDataArray() []*military_data.JiuGuanData {
	return dAtA.jiuGuanData.Array
}

func (dAtA *ConfigDatas) JiuGuanData() *JiuGuanDataConfig {
	return dAtA.jiuGuanData
}

type JiuGuanDataConfig struct {
	Map   map[uint64]*military_data.JiuGuanData
	Array []*military_data.JiuGuanData

	MinKeyData *military_data.JiuGuanData
	MaxKeyData *military_data.JiuGuanData

	parserMap map[*military_data.JiuGuanData]*config.ObjectParser
}

func (d *JiuGuanDataConfig) Get(key uint64) *military_data.JiuGuanData {
	return d.Map[key]
}

func (d *JiuGuanDataConfig) Must(key uint64) *military_data.JiuGuanData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetJunYingLevelData(key uint64) *military_data.JunYingLevelData {
	return dAtA.junYingLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetJunYingLevelDataArray() []*military_data.JunYingLevelData {
	return dAtA.junYingLevelData.Array
}

func (dAtA *ConfigDatas) JunYingLevelData() *JunYingLevelDataConfig {
	return dAtA.junYingLevelData
}

type JunYingLevelDataConfig struct {
	Map   map[uint64]*military_data.JunYingLevelData
	Array []*military_data.JunYingLevelData

	MinKeyData *military_data.JunYingLevelData
	MaxKeyData *military_data.JunYingLevelData

	parserMap map[*military_data.JunYingLevelData]*config.ObjectParser
}

func (d *JunYingLevelDataConfig) Get(key uint64) *military_data.JunYingLevelData {
	return d.Map[key]
}

func (d *JunYingLevelDataConfig) Must(key uint64) *military_data.JunYingLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTrainingLevelData(key uint64) *military_data.TrainingLevelData {
	return dAtA.trainingLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetTrainingLevelDataArray() []*military_data.TrainingLevelData {
	return dAtA.trainingLevelData.Array
}

func (dAtA *ConfigDatas) TrainingLevelData() *TrainingLevelDataConfig {
	return dAtA.trainingLevelData
}

type TrainingLevelDataConfig struct {
	Map   map[uint64]*military_data.TrainingLevelData
	Array []*military_data.TrainingLevelData

	MinKeyData *military_data.TrainingLevelData
	MaxKeyData *military_data.TrainingLevelData

	parserMap map[*military_data.TrainingLevelData]*config.ObjectParser
}

func (d *TrainingLevelDataConfig) Get(key uint64) *military_data.TrainingLevelData {
	return d.Map[key]
}

func (d *TrainingLevelDataConfig) Must(key uint64) *military_data.TrainingLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTutorData(key uint64) *military_data.TutorData {
	return dAtA.tutorData.Map[key]
}

func (dAtA *ConfigDatas) GetTutorDataArray() []*military_data.TutorData {
	return dAtA.tutorData.Array
}

func (dAtA *ConfigDatas) TutorData() *TutorDataConfig {
	return dAtA.tutorData
}

type TutorDataConfig struct {
	Map   map[uint64]*military_data.TutorData
	Array []*military_data.TutorData

	MinKeyData *military_data.TutorData
	MaxKeyData *military_data.TutorData

	parserMap map[*military_data.TutorData]*config.ObjectParser
}

func (d *TutorDataConfig) Get(key uint64) *military_data.TutorData {
	return d.Map[key]
}

func (d *TutorDataConfig) Must(key uint64) *military_data.TutorData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMcBuildAddSupportData(key uint64) *mingcdata.McBuildAddSupportData {
	return dAtA.mcBuildAddSupportData.Map[key]
}

func (dAtA *ConfigDatas) GetMcBuildAddSupportDataArray() []*mingcdata.McBuildAddSupportData {
	return dAtA.mcBuildAddSupportData.Array
}

func (dAtA *ConfigDatas) McBuildAddSupportData() *McBuildAddSupportDataConfig {
	return dAtA.mcBuildAddSupportData
}

type McBuildAddSupportDataConfig struct {
	Map   map[uint64]*mingcdata.McBuildAddSupportData
	Array []*mingcdata.McBuildAddSupportData

	MinKeyData *mingcdata.McBuildAddSupportData
	MaxKeyData *mingcdata.McBuildAddSupportData

	parserMap map[*mingcdata.McBuildAddSupportData]*config.ObjectParser
}

func (d *McBuildAddSupportDataConfig) Get(key uint64) *mingcdata.McBuildAddSupportData {
	return d.Map[key]
}

func (d *McBuildAddSupportDataConfig) Must(key uint64) *mingcdata.McBuildAddSupportData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMcBuildGuildMemberPrizeData(key uint64) *mingcdata.McBuildGuildMemberPrizeData {
	return dAtA.mcBuildGuildMemberPrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetMcBuildGuildMemberPrizeDataArray() []*mingcdata.McBuildGuildMemberPrizeData {
	return dAtA.mcBuildGuildMemberPrizeData.Array
}

func (dAtA *ConfigDatas) McBuildGuildMemberPrizeData() *McBuildGuildMemberPrizeDataConfig {
	return dAtA.mcBuildGuildMemberPrizeData
}

type McBuildGuildMemberPrizeDataConfig struct {
	Map   map[uint64]*mingcdata.McBuildGuildMemberPrizeData
	Array []*mingcdata.McBuildGuildMemberPrizeData

	MinKeyData *mingcdata.McBuildGuildMemberPrizeData
	MaxKeyData *mingcdata.McBuildGuildMemberPrizeData

	parserMap map[*mingcdata.McBuildGuildMemberPrizeData]*config.ObjectParser
}

func (d *McBuildGuildMemberPrizeDataConfig) Get(key uint64) *mingcdata.McBuildGuildMemberPrizeData {
	return d.Map[key]
}

func (d *McBuildGuildMemberPrizeDataConfig) Must(key uint64) *mingcdata.McBuildGuildMemberPrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMcBuildMcSupportData(key uint64) *mingcdata.McBuildMcSupportData {
	return dAtA.mcBuildMcSupportData.Map[key]
}

func (dAtA *ConfigDatas) GetMcBuildMcSupportDataArray() []*mingcdata.McBuildMcSupportData {
	return dAtA.mcBuildMcSupportData.Array
}

func (dAtA *ConfigDatas) McBuildMcSupportData() *McBuildMcSupportDataConfig {
	return dAtA.mcBuildMcSupportData
}

type McBuildMcSupportDataConfig struct {
	Map   map[uint64]*mingcdata.McBuildMcSupportData
	Array []*mingcdata.McBuildMcSupportData

	MinKeyData *mingcdata.McBuildMcSupportData
	MaxKeyData *mingcdata.McBuildMcSupportData

	parserMap map[*mingcdata.McBuildMcSupportData]*config.ObjectParser
}

func (d *McBuildMcSupportDataConfig) Get(key uint64) *mingcdata.McBuildMcSupportData {
	return d.Map[key]
}

func (d *McBuildMcSupportDataConfig) Must(key uint64) *mingcdata.McBuildMcSupportData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcBaseData(key uint64) *mingcdata.MingcBaseData {
	return dAtA.mingcBaseData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcBaseDataArray() []*mingcdata.MingcBaseData {
	return dAtA.mingcBaseData.Array
}

func (dAtA *ConfigDatas) MingcBaseData() *MingcBaseDataConfig {
	return dAtA.mingcBaseData
}

type MingcBaseDataConfig struct {
	Map   map[uint64]*mingcdata.MingcBaseData
	Array []*mingcdata.MingcBaseData

	MinKeyData *mingcdata.MingcBaseData
	MaxKeyData *mingcdata.MingcBaseData

	parserMap map[*mingcdata.MingcBaseData]*config.ObjectParser
}

func (d *MingcBaseDataConfig) Get(key uint64) *mingcdata.MingcBaseData {
	return d.Map[key]
}

func (d *MingcBaseDataConfig) Must(key uint64) *mingcdata.MingcBaseData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcTimeData(key uint64) *mingcdata.MingcTimeData {
	return dAtA.mingcTimeData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcTimeDataArray() []*mingcdata.MingcTimeData {
	return dAtA.mingcTimeData.Array
}

func (dAtA *ConfigDatas) MingcTimeData() *MingcTimeDataConfig {
	return dAtA.mingcTimeData
}

type MingcTimeDataConfig struct {
	Map   map[uint64]*mingcdata.MingcTimeData
	Array []*mingcdata.MingcTimeData

	MinKeyData *mingcdata.MingcTimeData
	MaxKeyData *mingcdata.MingcTimeData

	parserMap map[*mingcdata.MingcTimeData]*config.ObjectParser
}

func (d *MingcTimeDataConfig) Get(key uint64) *mingcdata.MingcTimeData {
	return d.Map[key]
}

func (d *MingcTimeDataConfig) Must(key uint64) *mingcdata.MingcTimeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcWarBuildingData(key uint64) *mingcdata.MingcWarBuildingData {
	return dAtA.mingcWarBuildingData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcWarBuildingDataArray() []*mingcdata.MingcWarBuildingData {
	return dAtA.mingcWarBuildingData.Array
}

func (dAtA *ConfigDatas) MingcWarBuildingData() *MingcWarBuildingDataConfig {
	return dAtA.mingcWarBuildingData
}

type MingcWarBuildingDataConfig struct {
	Map   map[uint64]*mingcdata.MingcWarBuildingData
	Array []*mingcdata.MingcWarBuildingData

	MinKeyData *mingcdata.MingcWarBuildingData
	MaxKeyData *mingcdata.MingcWarBuildingData

	parserMap map[*mingcdata.MingcWarBuildingData]*config.ObjectParser
}

func (d *MingcWarBuildingDataConfig) Get(key uint64) *mingcdata.MingcWarBuildingData {
	return d.Map[key]
}

func (d *MingcWarBuildingDataConfig) Must(key uint64) *mingcdata.MingcWarBuildingData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcWarDrumStatData(key uint64) *mingcdata.MingcWarDrumStatData {
	return dAtA.mingcWarDrumStatData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcWarDrumStatDataArray() []*mingcdata.MingcWarDrumStatData {
	return dAtA.mingcWarDrumStatData.Array
}

func (dAtA *ConfigDatas) MingcWarDrumStatData() *MingcWarDrumStatDataConfig {
	return dAtA.mingcWarDrumStatData
}

type MingcWarDrumStatDataConfig struct {
	Map   map[uint64]*mingcdata.MingcWarDrumStatData
	Array []*mingcdata.MingcWarDrumStatData

	MinKeyData *mingcdata.MingcWarDrumStatData
	MaxKeyData *mingcdata.MingcWarDrumStatData

	parserMap map[*mingcdata.MingcWarDrumStatData]*config.ObjectParser
}

func (d *MingcWarDrumStatDataConfig) Get(key uint64) *mingcdata.MingcWarDrumStatData {
	return d.Map[key]
}

func (d *MingcWarDrumStatDataConfig) Must(key uint64) *mingcdata.MingcWarDrumStatData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcWarMapData(key uint64) *mingcdata.MingcWarMapData {
	return dAtA.mingcWarMapData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcWarMapDataArray() []*mingcdata.MingcWarMapData {
	return dAtA.mingcWarMapData.Array
}

func (dAtA *ConfigDatas) MingcWarMapData() *MingcWarMapDataConfig {
	return dAtA.mingcWarMapData
}

type MingcWarMapDataConfig struct {
	Map   map[uint64]*mingcdata.MingcWarMapData
	Array []*mingcdata.MingcWarMapData

	MinKeyData *mingcdata.MingcWarMapData
	MaxKeyData *mingcdata.MingcWarMapData

	parserMap map[*mingcdata.MingcWarMapData]*config.ObjectParser
}

func (d *MingcWarMapDataConfig) Get(key uint64) *mingcdata.MingcWarMapData {
	return d.Map[key]
}

func (d *MingcWarMapDataConfig) Must(key uint64) *mingcdata.MingcWarMapData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcWarMultiKillData(key uint64) *mingcdata.MingcWarMultiKillData {
	return dAtA.mingcWarMultiKillData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcWarMultiKillDataArray() []*mingcdata.MingcWarMultiKillData {
	return dAtA.mingcWarMultiKillData.Array
}

func (dAtA *ConfigDatas) MingcWarMultiKillData() *MingcWarMultiKillDataConfig {
	return dAtA.mingcWarMultiKillData
}

type MingcWarMultiKillDataConfig struct {
	Map   map[uint64]*mingcdata.MingcWarMultiKillData
	Array []*mingcdata.MingcWarMultiKillData

	MinKeyData *mingcdata.MingcWarMultiKillData
	MaxKeyData *mingcdata.MingcWarMultiKillData

	parserMap map[*mingcdata.MingcWarMultiKillData]*config.ObjectParser
}

func (d *MingcWarMultiKillDataConfig) Get(key uint64) *mingcdata.MingcWarMultiKillData {
	return d.Map[key]
}

func (d *MingcWarMultiKillDataConfig) Must(key uint64) *mingcdata.MingcWarMultiKillData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcWarNpcData(key uint64) *mingcdata.MingcWarNpcData {
	return dAtA.mingcWarNpcData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcWarNpcDataArray() []*mingcdata.MingcWarNpcData {
	return dAtA.mingcWarNpcData.Array
}

func (dAtA *ConfigDatas) MingcWarNpcData() *MingcWarNpcDataConfig {
	return dAtA.mingcWarNpcData
}

type MingcWarNpcDataConfig struct {
	Map   map[uint64]*mingcdata.MingcWarNpcData
	Array []*mingcdata.MingcWarNpcData

	MinKeyData *mingcdata.MingcWarNpcData
	MaxKeyData *mingcdata.MingcWarNpcData

	parserMap map[*mingcdata.MingcWarNpcData]*config.ObjectParser
}

func (d *MingcWarNpcDataConfig) Get(key uint64) *mingcdata.MingcWarNpcData {
	return d.Map[key]
}

func (d *MingcWarNpcDataConfig) Must(key uint64) *mingcdata.MingcWarNpcData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcWarNpcGuildData(key uint64) *mingcdata.MingcWarNpcGuildData {
	return dAtA.mingcWarNpcGuildData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcWarNpcGuildDataArray() []*mingcdata.MingcWarNpcGuildData {
	return dAtA.mingcWarNpcGuildData.Array
}

func (dAtA *ConfigDatas) MingcWarNpcGuildData() *MingcWarNpcGuildDataConfig {
	return dAtA.mingcWarNpcGuildData
}

type MingcWarNpcGuildDataConfig struct {
	Map   map[uint64]*mingcdata.MingcWarNpcGuildData
	Array []*mingcdata.MingcWarNpcGuildData

	MinKeyData *mingcdata.MingcWarNpcGuildData
	MaxKeyData *mingcdata.MingcWarNpcGuildData

	parserMap map[*mingcdata.MingcWarNpcGuildData]*config.ObjectParser
}

func (d *MingcWarNpcGuildDataConfig) Get(key uint64) *mingcdata.MingcWarNpcGuildData {
	return d.Map[key]
}

func (d *MingcWarNpcGuildDataConfig) Must(key uint64) *mingcdata.MingcWarNpcGuildData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcWarSceneData(key uint64) *mingcdata.MingcWarSceneData {
	return dAtA.mingcWarSceneData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcWarSceneDataArray() []*mingcdata.MingcWarSceneData {
	return dAtA.mingcWarSceneData.Array
}

func (dAtA *ConfigDatas) MingcWarSceneData() *MingcWarSceneDataConfig {
	return dAtA.mingcWarSceneData
}

type MingcWarSceneDataConfig struct {
	Map   map[uint64]*mingcdata.MingcWarSceneData
	Array []*mingcdata.MingcWarSceneData

	MinKeyData *mingcdata.MingcWarSceneData
	MaxKeyData *mingcdata.MingcWarSceneData

	parserMap map[*mingcdata.MingcWarSceneData]*config.ObjectParser
}

func (d *MingcWarSceneDataConfig) Get(key uint64) *mingcdata.MingcWarSceneData {
	return d.Map[key]
}

func (d *MingcWarSceneDataConfig) Must(key uint64) *mingcdata.MingcWarSceneData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcWarTouShiBuildingTargetData(key uint64) *mingcdata.MingcWarTouShiBuildingTargetData {
	return dAtA.mingcWarTouShiBuildingTargetData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcWarTouShiBuildingTargetDataArray() []*mingcdata.MingcWarTouShiBuildingTargetData {
	return dAtA.mingcWarTouShiBuildingTargetData.Array
}

func (dAtA *ConfigDatas) MingcWarTouShiBuildingTargetData() *MingcWarTouShiBuildingTargetDataConfig {
	return dAtA.mingcWarTouShiBuildingTargetData
}

type MingcWarTouShiBuildingTargetDataConfig struct {
	Map   map[uint64]*mingcdata.MingcWarTouShiBuildingTargetData
	Array []*mingcdata.MingcWarTouShiBuildingTargetData

	MinKeyData *mingcdata.MingcWarTouShiBuildingTargetData
	MaxKeyData *mingcdata.MingcWarTouShiBuildingTargetData

	parserMap map[*mingcdata.MingcWarTouShiBuildingTargetData]*config.ObjectParser
}

func (d *MingcWarTouShiBuildingTargetDataConfig) Get(key uint64) *mingcdata.MingcWarTouShiBuildingTargetData {
	return d.Map[key]
}

func (d *MingcWarTouShiBuildingTargetDataConfig) Must(key uint64) *mingcdata.MingcWarTouShiBuildingTargetData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMingcWarTroopLastBeatWhenFailData(key uint64) *mingcdata.MingcWarTroopLastBeatWhenFailData {
	return dAtA.mingcWarTroopLastBeatWhenFailData.Map[key]
}

func (dAtA *ConfigDatas) GetMingcWarTroopLastBeatWhenFailDataArray() []*mingcdata.MingcWarTroopLastBeatWhenFailData {
	return dAtA.mingcWarTroopLastBeatWhenFailData.Array
}

func (dAtA *ConfigDatas) MingcWarTroopLastBeatWhenFailData() *MingcWarTroopLastBeatWhenFailDataConfig {
	return dAtA.mingcWarTroopLastBeatWhenFailData
}

type MingcWarTroopLastBeatWhenFailDataConfig struct {
	Map   map[uint64]*mingcdata.MingcWarTroopLastBeatWhenFailData
	Array []*mingcdata.MingcWarTroopLastBeatWhenFailData

	MinKeyData *mingcdata.MingcWarTroopLastBeatWhenFailData
	MaxKeyData *mingcdata.MingcWarTroopLastBeatWhenFailData

	parserMap map[*mingcdata.MingcWarTroopLastBeatWhenFailData]*config.ObjectParser
}

func (d *MingcWarTroopLastBeatWhenFailDataConfig) Get(key uint64) *mingcdata.MingcWarTroopLastBeatWhenFailData {
	return d.Map[key]
}

func (d *MingcWarTroopLastBeatWhenFailDataConfig) Must(key uint64) *mingcdata.MingcWarTroopLastBeatWhenFailData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMonsterCaptainData(key uint64) *monsterdata.MonsterCaptainData {
	return dAtA.monsterCaptainData.Map[key]
}

func (dAtA *ConfigDatas) GetMonsterCaptainDataArray() []*monsterdata.MonsterCaptainData {
	return dAtA.monsterCaptainData.Array
}

func (dAtA *ConfigDatas) MonsterCaptainData() *MonsterCaptainDataConfig {
	return dAtA.monsterCaptainData
}

type MonsterCaptainDataConfig struct {
	Map   map[uint64]*monsterdata.MonsterCaptainData
	Array []*monsterdata.MonsterCaptainData

	MinKeyData *monsterdata.MonsterCaptainData
	MaxKeyData *monsterdata.MonsterCaptainData

	parserMap map[*monsterdata.MonsterCaptainData]*config.ObjectParser
}

func (d *MonsterCaptainDataConfig) Get(key uint64) *monsterdata.MonsterCaptainData {
	return d.Map[key]
}

func (d *MonsterCaptainDataConfig) Must(key uint64) *monsterdata.MonsterCaptainData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMonsterMasterData(key uint64) *monsterdata.MonsterMasterData {
	return dAtA.monsterMasterData.Map[key]
}

func (dAtA *ConfigDatas) GetMonsterMasterDataArray() []*monsterdata.MonsterMasterData {
	return dAtA.monsterMasterData.Array
}

func (dAtA *ConfigDatas) MonsterMasterData() *MonsterMasterDataConfig {
	return dAtA.monsterMasterData
}

type MonsterMasterDataConfig struct {
	Map   map[uint64]*monsterdata.MonsterMasterData
	Array []*monsterdata.MonsterMasterData

	MinKeyData *monsterdata.MonsterMasterData
	MaxKeyData *monsterdata.MonsterMasterData

	parserMap map[*monsterdata.MonsterMasterData]*config.ObjectParser
}

func (d *MonsterMasterDataConfig) Get(key uint64) *monsterdata.MonsterMasterData {
	return d.Map[key]
}

func (d *MonsterMasterDataConfig) Must(key uint64) *monsterdata.MonsterMasterData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetDailyBargainData(key uint64) *promdata.DailyBargainData {
	return dAtA.dailyBargainData.Map[key]
}

func (dAtA *ConfigDatas) GetDailyBargainDataArray() []*promdata.DailyBargainData {
	return dAtA.dailyBargainData.Array
}

func (dAtA *ConfigDatas) DailyBargainData() *DailyBargainDataConfig {
	return dAtA.dailyBargainData
}

type DailyBargainDataConfig struct {
	Map   map[uint64]*promdata.DailyBargainData
	Array []*promdata.DailyBargainData

	MinKeyData *promdata.DailyBargainData
	MaxKeyData *promdata.DailyBargainData

	parserMap map[*promdata.DailyBargainData]*config.ObjectParser
}

func (d *DailyBargainDataConfig) Get(key uint64) *promdata.DailyBargainData {
	return d.Map[key]
}

func (d *DailyBargainDataConfig) Must(key uint64) *promdata.DailyBargainData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetDurationCardData(key uint64) *promdata.DurationCardData {
	return dAtA.durationCardData.Map[key]
}

func (dAtA *ConfigDatas) GetDurationCardDataArray() []*promdata.DurationCardData {
	return dAtA.durationCardData.Array
}

func (dAtA *ConfigDatas) DurationCardData() *DurationCardDataConfig {
	return dAtA.durationCardData
}

type DurationCardDataConfig struct {
	Map   map[uint64]*promdata.DurationCardData
	Array []*promdata.DurationCardData

	MinKeyData *promdata.DurationCardData
	MaxKeyData *promdata.DurationCardData

	parserMap map[*promdata.DurationCardData]*config.ObjectParser
}

func (d *DurationCardDataConfig) Get(key uint64) *promdata.DurationCardData {
	return d.Map[key]
}

func (d *DurationCardDataConfig) Must(key uint64) *promdata.DurationCardData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetEventLimitGiftData(key uint64) *promdata.EventLimitGiftData {
	return dAtA.eventLimitGiftData.Map[key]
}

func (dAtA *ConfigDatas) GetEventLimitGiftDataArray() []*promdata.EventLimitGiftData {
	return dAtA.eventLimitGiftData.Array
}

func (dAtA *ConfigDatas) EventLimitGiftData() *EventLimitGiftDataConfig {
	return dAtA.eventLimitGiftData
}

type EventLimitGiftDataConfig struct {
	Map   map[uint64]*promdata.EventLimitGiftData
	Array []*promdata.EventLimitGiftData

	MinKeyData *promdata.EventLimitGiftData
	MaxKeyData *promdata.EventLimitGiftData

	parserMap map[*promdata.EventLimitGiftData]*config.ObjectParser
}

func (d *EventLimitGiftDataConfig) Get(key uint64) *promdata.EventLimitGiftData {
	return d.Map[key]
}

func (d *EventLimitGiftDataConfig) Must(key uint64) *promdata.EventLimitGiftData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetFreeGiftData(key uint64) *promdata.FreeGiftData {
	return dAtA.freeGiftData.Map[key]
}

func (dAtA *ConfigDatas) GetFreeGiftDataArray() []*promdata.FreeGiftData {
	return dAtA.freeGiftData.Array
}

func (dAtA *ConfigDatas) FreeGiftData() *FreeGiftDataConfig {
	return dAtA.freeGiftData
}

type FreeGiftDataConfig struct {
	Map   map[uint64]*promdata.FreeGiftData
	Array []*promdata.FreeGiftData

	MinKeyData *promdata.FreeGiftData
	MaxKeyData *promdata.FreeGiftData

	parserMap map[*promdata.FreeGiftData]*config.ObjectParser
}

func (d *FreeGiftDataConfig) Get(key uint64) *promdata.FreeGiftData {
	return d.Map[key]
}

func (d *FreeGiftDataConfig) Must(key uint64) *promdata.FreeGiftData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetHeroLevelFundData(key uint64) *promdata.HeroLevelFundData {
	return dAtA.heroLevelFundData.Map[key]
}

func (dAtA *ConfigDatas) GetHeroLevelFundDataArray() []*promdata.HeroLevelFundData {
	return dAtA.heroLevelFundData.Array
}

func (dAtA *ConfigDatas) HeroLevelFundData() *HeroLevelFundDataConfig {
	return dAtA.heroLevelFundData
}

type HeroLevelFundDataConfig struct {
	Map   map[uint64]*promdata.HeroLevelFundData
	Array []*promdata.HeroLevelFundData

	MinKeyData *promdata.HeroLevelFundData
	MaxKeyData *promdata.HeroLevelFundData

	parserMap map[*promdata.HeroLevelFundData]*config.ObjectParser
}

func (d *HeroLevelFundDataConfig) Get(key uint64) *promdata.HeroLevelFundData {
	return d.Map[key]
}

func (d *HeroLevelFundDataConfig) Must(key uint64) *promdata.HeroLevelFundData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetLoginDayData(key uint64) *promdata.LoginDayData {
	return dAtA.loginDayData.Map[key]
}

func (dAtA *ConfigDatas) GetLoginDayDataArray() []*promdata.LoginDayData {
	return dAtA.loginDayData.Array
}

func (dAtA *ConfigDatas) LoginDayData() *LoginDayDataConfig {
	return dAtA.loginDayData
}

type LoginDayDataConfig struct {
	Map   map[uint64]*promdata.LoginDayData
	Array []*promdata.LoginDayData

	MinKeyData *promdata.LoginDayData
	MaxKeyData *promdata.LoginDayData

	parserMap map[*promdata.LoginDayData]*config.ObjectParser
}

func (d *LoginDayDataConfig) Get(key uint64) *promdata.LoginDayData {
	return d.Map[key]
}

func (d *LoginDayDataConfig) Must(key uint64) *promdata.LoginDayData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetSpCollectionData(key uint64) *promdata.SpCollectionData {
	return dAtA.spCollectionData.Map[key]
}

func (dAtA *ConfigDatas) GetSpCollectionDataArray() []*promdata.SpCollectionData {
	return dAtA.spCollectionData.Array
}

func (dAtA *ConfigDatas) SpCollectionData() *SpCollectionDataConfig {
	return dAtA.spCollectionData
}

type SpCollectionDataConfig struct {
	Map   map[uint64]*promdata.SpCollectionData
	Array []*promdata.SpCollectionData

	MinKeyData *promdata.SpCollectionData
	MaxKeyData *promdata.SpCollectionData

	parserMap map[*promdata.SpCollectionData]*config.ObjectParser
}

func (d *SpCollectionDataConfig) Get(key uint64) *promdata.SpCollectionData {
	return d.Map[key]
}

func (d *SpCollectionDataConfig) Must(key uint64) *promdata.SpCollectionData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTimeLimitGiftData(key uint64) *promdata.TimeLimitGiftData {
	return dAtA.timeLimitGiftData.Map[key]
}

func (dAtA *ConfigDatas) GetTimeLimitGiftDataArray() []*promdata.TimeLimitGiftData {
	return dAtA.timeLimitGiftData.Array
}

func (dAtA *ConfigDatas) TimeLimitGiftData() *TimeLimitGiftDataConfig {
	return dAtA.timeLimitGiftData
}

type TimeLimitGiftDataConfig struct {
	Map   map[uint64]*promdata.TimeLimitGiftData
	Array []*promdata.TimeLimitGiftData

	MinKeyData *promdata.TimeLimitGiftData
	MaxKeyData *promdata.TimeLimitGiftData

	parserMap map[*promdata.TimeLimitGiftData]*config.ObjectParser
}

func (d *TimeLimitGiftDataConfig) Get(key uint64) *promdata.TimeLimitGiftData {
	return d.Map[key]
}

func (d *TimeLimitGiftDataConfig) Must(key uint64) *promdata.TimeLimitGiftData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTimeLimitGiftGroupData(key uint64) *promdata.TimeLimitGiftGroupData {
	return dAtA.timeLimitGiftGroupData.Map[key]
}

func (dAtA *ConfigDatas) GetTimeLimitGiftGroupDataArray() []*promdata.TimeLimitGiftGroupData {
	return dAtA.timeLimitGiftGroupData.Array
}

func (dAtA *ConfigDatas) TimeLimitGiftGroupData() *TimeLimitGiftGroupDataConfig {
	return dAtA.timeLimitGiftGroupData
}

type TimeLimitGiftGroupDataConfig struct {
	Map   map[uint64]*promdata.TimeLimitGiftGroupData
	Array []*promdata.TimeLimitGiftGroupData

	MinKeyData *promdata.TimeLimitGiftGroupData
	MaxKeyData *promdata.TimeLimitGiftGroupData

	parserMap map[*promdata.TimeLimitGiftGroupData]*config.ObjectParser
}

func (d *TimeLimitGiftGroupDataConfig) Get(key uint64) *promdata.TimeLimitGiftGroupData {
	return d.Map[key]
}

func (d *TimeLimitGiftGroupDataConfig) Must(key uint64) *promdata.TimeLimitGiftGroupData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetPushData(key uint64) *pushdata.PushData {
	return dAtA.pushData.Map[key]
}

func (dAtA *ConfigDatas) GetPushDataArray() []*pushdata.PushData {
	return dAtA.pushData.Array
}

func (dAtA *ConfigDatas) PushData() *PushDataConfig {
	return dAtA.pushData
}

type PushDataConfig struct {
	Map   map[uint64]*pushdata.PushData
	Array []*pushdata.PushData

	MinKeyData *pushdata.PushData
	MaxKeyData *pushdata.PushData

	parserMap map[*pushdata.PushData]*config.ObjectParser
}

func (d *PushDataConfig) Get(key uint64) *pushdata.PushData {
	return d.Map[key]
}

func (d *PushDataConfig) Must(key uint64) *pushdata.PushData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetPveTroopData(key uint64) *pvetroop.PveTroopData {
	return dAtA.pveTroopData.Map[key]
}

func (dAtA *ConfigDatas) GetPveTroopDataArray() []*pvetroop.PveTroopData {
	return dAtA.pveTroopData.Array
}

func (dAtA *ConfigDatas) PveTroopData() *PveTroopDataConfig {
	return dAtA.pveTroopData
}

type PveTroopDataConfig struct {
	Map   map[uint64]*pvetroop.PveTroopData
	Array []*pvetroop.PveTroopData

	MinKeyData *pvetroop.PveTroopData
	MaxKeyData *pvetroop.PveTroopData

	parserMap map[*pvetroop.PveTroopData]*config.ObjectParser
}

func (d *PveTroopDataConfig) Get(key uint64) *pvetroop.PveTroopData {
	return d.Map[key]
}

func (d *PveTroopDataConfig) Must(key uint64) *pvetroop.PveTroopData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetQuestionData(key uint64) *question.QuestionData {
	return dAtA.questionData.Map[key]
}

func (dAtA *ConfigDatas) GetQuestionDataArray() []*question.QuestionData {
	return dAtA.questionData.Array
}

func (dAtA *ConfigDatas) QuestionData() *QuestionDataConfig {
	return dAtA.questionData
}

type QuestionDataConfig struct {
	Map   map[uint64]*question.QuestionData
	Array []*question.QuestionData

	MinKeyData *question.QuestionData
	MaxKeyData *question.QuestionData

	parserMap map[*question.QuestionData]*config.ObjectParser
}

func (d *QuestionDataConfig) Get(key uint64) *question.QuestionData {
	return d.Map[key]
}

func (d *QuestionDataConfig) Must(key uint64) *question.QuestionData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetQuestionPrizeData(key uint64) *question.QuestionPrizeData {
	return dAtA.questionPrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetQuestionPrizeDataArray() []*question.QuestionPrizeData {
	return dAtA.questionPrizeData.Array
}

func (dAtA *ConfigDatas) QuestionPrizeData() *QuestionPrizeDataConfig {
	return dAtA.questionPrizeData
}

type QuestionPrizeDataConfig struct {
	Map   map[uint64]*question.QuestionPrizeData
	Array []*question.QuestionPrizeData

	MinKeyData *question.QuestionPrizeData
	MaxKeyData *question.QuestionPrizeData

	parserMap map[*question.QuestionPrizeData]*config.ObjectParser
}

func (d *QuestionPrizeDataConfig) Get(key uint64) *question.QuestionPrizeData {
	return d.Map[key]
}

func (d *QuestionPrizeDataConfig) Must(key uint64) *question.QuestionPrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetQuestionSayingData(key uint64) *question.QuestionSayingData {
	return dAtA.questionSayingData.Map[key]
}

func (dAtA *ConfigDatas) GetQuestionSayingDataArray() []*question.QuestionSayingData {
	return dAtA.questionSayingData.Array
}

func (dAtA *ConfigDatas) QuestionSayingData() *QuestionSayingDataConfig {
	return dAtA.questionSayingData
}

type QuestionSayingDataConfig struct {
	Map   map[uint64]*question.QuestionSayingData
	Array []*question.QuestionSayingData

	MinKeyData *question.QuestionSayingData
	MaxKeyData *question.QuestionSayingData

	parserMap map[*question.QuestionSayingData]*config.ObjectParser
}

func (d *QuestionSayingDataConfig) Get(key uint64) *question.QuestionSayingData {
	return d.Map[key]
}

func (d *QuestionSayingDataConfig) Must(key uint64) *question.QuestionSayingData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetRaceData(key int) *race.RaceData {
	return dAtA.raceData.Map[key]
}

func (dAtA *ConfigDatas) GetRaceDataArray() []*race.RaceData {
	return dAtA.raceData.Array
}

func (dAtA *ConfigDatas) RaceData() *RaceDataConfig {
	return dAtA.raceData
}

type RaceDataConfig struct {
	Map   map[int]*race.RaceData
	Array []*race.RaceData

	MinKeyData *race.RaceData
	MaxKeyData *race.RaceData

	parserMap map[*race.RaceData]*config.ObjectParser
}

func (d *RaceDataConfig) Get(key int) *race.RaceData {
	return d.Map[key]
}

func (d *RaceDataConfig) Must(key int) *race.RaceData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetEventOptionData(key uint64) *random_event.EventOptionData {
	return dAtA.eventOptionData.Map[key]
}

func (dAtA *ConfigDatas) GetEventOptionDataArray() []*random_event.EventOptionData {
	return dAtA.eventOptionData.Array
}

func (dAtA *ConfigDatas) EventOptionData() *EventOptionDataConfig {
	return dAtA.eventOptionData
}

type EventOptionDataConfig struct {
	Map   map[uint64]*random_event.EventOptionData
	Array []*random_event.EventOptionData

	MinKeyData *random_event.EventOptionData
	MaxKeyData *random_event.EventOptionData

	parserMap map[*random_event.EventOptionData]*config.ObjectParser
}

func (d *EventOptionDataConfig) Get(key uint64) *random_event.EventOptionData {
	return d.Map[key]
}

func (d *EventOptionDataConfig) Must(key uint64) *random_event.EventOptionData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetEventPosition(key uint64) *random_event.EventPosition {
	return dAtA.eventPosition.Map[key]
}

func (dAtA *ConfigDatas) GetEventPositionArray() []*random_event.EventPosition {
	return dAtA.eventPosition.Array
}

func (dAtA *ConfigDatas) EventPosition() *EventPositionConfig {
	return dAtA.eventPosition
}

type EventPositionConfig struct {
	Map   map[uint64]*random_event.EventPosition
	Array []*random_event.EventPosition

	MinKeyData *random_event.EventPosition
	MaxKeyData *random_event.EventPosition

	parserMap map[*random_event.EventPosition]*config.ObjectParser
}

func (d *EventPositionConfig) Get(key uint64) *random_event.EventPosition {
	return d.Map[key]
}

func (d *EventPositionConfig) Must(key uint64) *random_event.EventPosition {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetOptionPrize(key uint64) *random_event.OptionPrize {
	return dAtA.optionPrize.Map[key]
}

func (dAtA *ConfigDatas) GetOptionPrizeArray() []*random_event.OptionPrize {
	return dAtA.optionPrize.Array
}

func (dAtA *ConfigDatas) OptionPrize() *OptionPrizeConfig {
	return dAtA.optionPrize
}

type OptionPrizeConfig struct {
	Map   map[uint64]*random_event.OptionPrize
	Array []*random_event.OptionPrize

	MinKeyData *random_event.OptionPrize
	MaxKeyData *random_event.OptionPrize

	parserMap map[*random_event.OptionPrize]*config.ObjectParser
}

func (d *OptionPrizeConfig) Get(key uint64) *random_event.OptionPrize {
	return d.Map[key]
}

func (d *OptionPrizeConfig) Must(key uint64) *random_event.OptionPrize {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetRandomEventData(key uint64) *random_event.RandomEventData {
	return dAtA.randomEventData.Map[key]
}

func (dAtA *ConfigDatas) GetRandomEventDataArray() []*random_event.RandomEventData {
	return dAtA.randomEventData.Array
}

func (dAtA *ConfigDatas) RandomEventData() *RandomEventDataConfig {
	return dAtA.randomEventData
}

type RandomEventDataConfig struct {
	Map   map[uint64]*random_event.RandomEventData
	Array []*random_event.RandomEventData

	MinKeyData *random_event.RandomEventData
	MaxKeyData *random_event.RandomEventData

	parserMap map[*random_event.RandomEventData]*config.ObjectParser
}

func (d *RandomEventDataConfig) Get(key uint64) *random_event.RandomEventData {
	return d.Map[key]
}

func (d *RandomEventDataConfig) Must(key uint64) *random_event.RandomEventData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetRedPacketData(key uint64) *red_packet.RedPacketData {
	return dAtA.redPacketData.Map[key]
}

func (dAtA *ConfigDatas) GetRedPacketDataArray() []*red_packet.RedPacketData {
	return dAtA.redPacketData.Array
}

func (dAtA *ConfigDatas) RedPacketData() *RedPacketDataConfig {
	return dAtA.redPacketData
}

type RedPacketDataConfig struct {
	Map   map[uint64]*red_packet.RedPacketData
	Array []*red_packet.RedPacketData

	MinKeyData *red_packet.RedPacketData
	MaxKeyData *red_packet.RedPacketData

	parserMap map[*red_packet.RedPacketData]*config.ObjectParser
}

func (d *RedPacketDataConfig) Get(key uint64) *red_packet.RedPacketData {
	return d.Map[key]
}

func (d *RedPacketDataConfig) Must(key uint64) *red_packet.RedPacketData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetAreaData(key uint64) *regdata.AreaData {
	return dAtA.areaData.Map[key]
}

func (dAtA *ConfigDatas) GetAreaDataArray() []*regdata.AreaData {
	return dAtA.areaData.Array
}

func (dAtA *ConfigDatas) AreaData() *AreaDataConfig {
	return dAtA.areaData
}

type AreaDataConfig struct {
	Map   map[uint64]*regdata.AreaData
	Array []*regdata.AreaData

	MinKeyData *regdata.AreaData
	MaxKeyData *regdata.AreaData

	parserMap map[*regdata.AreaData]*config.ObjectParser
}

func (d *AreaDataConfig) Get(key uint64) *regdata.AreaData {
	return d.Map[key]
}

func (d *AreaDataConfig) Must(key uint64) *regdata.AreaData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetAssemblyData(key uint64) *regdata.AssemblyData {
	return dAtA.assemblyData.Map[key]
}

func (dAtA *ConfigDatas) GetAssemblyDataArray() []*regdata.AssemblyData {
	return dAtA.assemblyData.Array
}

func (dAtA *ConfigDatas) AssemblyData() *AssemblyDataConfig {
	return dAtA.assemblyData
}

type AssemblyDataConfig struct {
	Map   map[uint64]*regdata.AssemblyData
	Array []*regdata.AssemblyData

	MinKeyData *regdata.AssemblyData
	MaxKeyData *regdata.AssemblyData

	parserMap map[*regdata.AssemblyData]*config.ObjectParser
}

func (d *AssemblyDataConfig) Get(key uint64) *regdata.AssemblyData {
	return d.Map[key]
}

func (d *AssemblyDataConfig) Must(key uint64) *regdata.AssemblyData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBaozNpcData(key uint64) *regdata.BaozNpcData {
	return dAtA.baozNpcData.Map[key]
}

func (dAtA *ConfigDatas) GetBaozNpcDataArray() []*regdata.BaozNpcData {
	return dAtA.baozNpcData.Array
}

func (dAtA *ConfigDatas) BaozNpcData() *BaozNpcDataConfig {
	return dAtA.baozNpcData
}

type BaozNpcDataConfig struct {
	Map   map[uint64]*regdata.BaozNpcData
	Array []*regdata.BaozNpcData

	MinKeyData *regdata.BaozNpcData
	MaxKeyData *regdata.BaozNpcData

	parserMap map[*regdata.BaozNpcData]*config.ObjectParser
}

func (d *BaozNpcDataConfig) Get(key uint64) *regdata.BaozNpcData {
	return d.Map[key]
}

func (d *BaozNpcDataConfig) Must(key uint64) *regdata.BaozNpcData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetJunTuanNpcData(key uint64) *regdata.JunTuanNpcData {
	return dAtA.junTuanNpcData.Map[key]
}

func (dAtA *ConfigDatas) GetJunTuanNpcDataArray() []*regdata.JunTuanNpcData {
	return dAtA.junTuanNpcData.Array
}

func (dAtA *ConfigDatas) JunTuanNpcData() *JunTuanNpcDataConfig {
	return dAtA.junTuanNpcData
}

type JunTuanNpcDataConfig struct {
	Map   map[uint64]*regdata.JunTuanNpcData
	Array []*regdata.JunTuanNpcData

	MinKeyData *regdata.JunTuanNpcData
	MaxKeyData *regdata.JunTuanNpcData

	parserMap map[*regdata.JunTuanNpcData]*config.ObjectParser
}

func (d *JunTuanNpcDataConfig) Get(key uint64) *regdata.JunTuanNpcData {
	return d.Map[key]
}

func (d *JunTuanNpcDataConfig) Must(key uint64) *regdata.JunTuanNpcData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetJunTuanNpcPlaceData(key uint64) *regdata.JunTuanNpcPlaceData {
	return dAtA.junTuanNpcPlaceData.Map[key]
}

func (dAtA *ConfigDatas) GetJunTuanNpcPlaceDataArray() []*regdata.JunTuanNpcPlaceData {
	return dAtA.junTuanNpcPlaceData.Array
}

func (dAtA *ConfigDatas) JunTuanNpcPlaceData() *JunTuanNpcPlaceDataConfig {
	return dAtA.junTuanNpcPlaceData
}

type JunTuanNpcPlaceDataConfig struct {
	Map   map[uint64]*regdata.JunTuanNpcPlaceData
	Array []*regdata.JunTuanNpcPlaceData

	MinKeyData *regdata.JunTuanNpcPlaceData
	MaxKeyData *regdata.JunTuanNpcPlaceData

	parserMap map[*regdata.JunTuanNpcPlaceData]*config.ObjectParser
}

func (d *JunTuanNpcPlaceDataConfig) Get(key uint64) *regdata.JunTuanNpcPlaceData {
	return d.Map[key]
}

func (d *JunTuanNpcPlaceDataConfig) Must(key uint64) *regdata.JunTuanNpcPlaceData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetRegionAreaData(key uint64) *regdata.RegionAreaData {
	return dAtA.regionAreaData.Map[key]
}

func (dAtA *ConfigDatas) GetRegionAreaDataArray() []*regdata.RegionAreaData {
	return dAtA.regionAreaData.Array
}

func (dAtA *ConfigDatas) RegionAreaData() *RegionAreaDataConfig {
	return dAtA.regionAreaData
}

type RegionAreaDataConfig struct {
	Map   map[uint64]*regdata.RegionAreaData
	Array []*regdata.RegionAreaData

	MinKeyData *regdata.RegionAreaData
	MaxKeyData *regdata.RegionAreaData

	parserMap map[*regdata.RegionAreaData]*config.ObjectParser
}

func (d *RegionAreaDataConfig) Get(key uint64) *regdata.RegionAreaData {
	return d.Map[key]
}

func (d *RegionAreaDataConfig) Must(key uint64) *regdata.RegionAreaData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetRegionData(key uint64) *regdata.RegionData {
	return dAtA.regionData.Map[key]
}

func (dAtA *ConfigDatas) GetRegionDataArray() []*regdata.RegionData {
	return dAtA.regionData.Array
}

func (dAtA *ConfigDatas) RegionData() *RegionDataConfig {
	return dAtA.regionData
}

type RegionDataConfig struct {
	Map   map[uint64]*regdata.RegionData
	Array []*regdata.RegionData

	MinKeyData *regdata.RegionData
	MaxKeyData *regdata.RegionData

	parserMap map[*regdata.RegionData]*config.ObjectParser
}

func (d *RegionDataConfig) Get(key uint64) *regdata.RegionData {
	return d.Map[key]
}

func (d *RegionDataConfig) Must(key uint64) *regdata.RegionData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetRegionMonsterData(key uint64) *regdata.RegionMonsterData {
	return dAtA.regionMonsterData.Map[key]
}

func (dAtA *ConfigDatas) GetRegionMonsterDataArray() []*regdata.RegionMonsterData {
	return dAtA.regionMonsterData.Array
}

func (dAtA *ConfigDatas) RegionMonsterData() *RegionMonsterDataConfig {
	return dAtA.regionMonsterData
}

type RegionMonsterDataConfig struct {
	Map   map[uint64]*regdata.RegionMonsterData
	Array []*regdata.RegionMonsterData

	MinKeyData *regdata.RegionMonsterData
	MaxKeyData *regdata.RegionMonsterData

	parserMap map[*regdata.RegionMonsterData]*config.ObjectParser
}

func (d *RegionMonsterDataConfig) Get(key uint64) *regdata.RegionMonsterData {
	return d.Map[key]
}

func (d *RegionMonsterDataConfig) Must(key uint64) *regdata.RegionMonsterData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetRegionMultiLevelNpcData(key uint64) *regdata.RegionMultiLevelNpcData {
	return dAtA.regionMultiLevelNpcData.Map[key]
}

func (dAtA *ConfigDatas) GetRegionMultiLevelNpcDataArray() []*regdata.RegionMultiLevelNpcData {
	return dAtA.regionMultiLevelNpcData.Array
}

func (dAtA *ConfigDatas) RegionMultiLevelNpcData() *RegionMultiLevelNpcDataConfig {
	return dAtA.regionMultiLevelNpcData
}

type RegionMultiLevelNpcDataConfig struct {
	Map   map[uint64]*regdata.RegionMultiLevelNpcData
	Array []*regdata.RegionMultiLevelNpcData

	MinKeyData *regdata.RegionMultiLevelNpcData
	MaxKeyData *regdata.RegionMultiLevelNpcData

	parserMap map[*regdata.RegionMultiLevelNpcData]*config.ObjectParser
}

func (d *RegionMultiLevelNpcDataConfig) Get(key uint64) *regdata.RegionMultiLevelNpcData {
	return d.Map[key]
}

func (d *RegionMultiLevelNpcDataConfig) Must(key uint64) *regdata.RegionMultiLevelNpcData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetRegionMultiLevelNpcLevelData(key uint64) *regdata.RegionMultiLevelNpcLevelData {
	return dAtA.regionMultiLevelNpcLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetRegionMultiLevelNpcLevelDataArray() []*regdata.RegionMultiLevelNpcLevelData {
	return dAtA.regionMultiLevelNpcLevelData.Array
}

func (dAtA *ConfigDatas) RegionMultiLevelNpcLevelData() *RegionMultiLevelNpcLevelDataConfig {
	return dAtA.regionMultiLevelNpcLevelData
}

type RegionMultiLevelNpcLevelDataConfig struct {
	Map   map[uint64]*regdata.RegionMultiLevelNpcLevelData
	Array []*regdata.RegionMultiLevelNpcLevelData

	MinKeyData *regdata.RegionMultiLevelNpcLevelData
	MaxKeyData *regdata.RegionMultiLevelNpcLevelData

	parserMap map[*regdata.RegionMultiLevelNpcLevelData]*config.ObjectParser
}

func (d *RegionMultiLevelNpcLevelDataConfig) Get(key uint64) *regdata.RegionMultiLevelNpcLevelData {
	return d.Map[key]
}

func (d *RegionMultiLevelNpcLevelDataConfig) Must(key uint64) *regdata.RegionMultiLevelNpcLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetRegionMultiLevelNpcTypeData(key int) *regdata.RegionMultiLevelNpcTypeData {
	return dAtA.regionMultiLevelNpcTypeData.Map[key]
}

func (dAtA *ConfigDatas) GetRegionMultiLevelNpcTypeDataArray() []*regdata.RegionMultiLevelNpcTypeData {
	return dAtA.regionMultiLevelNpcTypeData.Array
}

func (dAtA *ConfigDatas) RegionMultiLevelNpcTypeData() *RegionMultiLevelNpcTypeDataConfig {
	return dAtA.regionMultiLevelNpcTypeData
}

type RegionMultiLevelNpcTypeDataConfig struct {
	Map   map[int]*regdata.RegionMultiLevelNpcTypeData
	Array []*regdata.RegionMultiLevelNpcTypeData

	MinKeyData *regdata.RegionMultiLevelNpcTypeData
	MaxKeyData *regdata.RegionMultiLevelNpcTypeData

	parserMap map[*regdata.RegionMultiLevelNpcTypeData]*config.ObjectParser
}

func (d *RegionMultiLevelNpcTypeDataConfig) Get(key int) *regdata.RegionMultiLevelNpcTypeData {
	return d.Map[key]
}

func (d *RegionMultiLevelNpcTypeDataConfig) Must(key int) *regdata.RegionMultiLevelNpcTypeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTroopDialogueData(key uint64) *regdata.TroopDialogueData {
	return dAtA.troopDialogueData.Map[key]
}

func (dAtA *ConfigDatas) GetTroopDialogueDataArray() []*regdata.TroopDialogueData {
	return dAtA.troopDialogueData.Array
}

func (dAtA *ConfigDatas) TroopDialogueData() *TroopDialogueDataConfig {
	return dAtA.troopDialogueData
}

type TroopDialogueDataConfig struct {
	Map   map[uint64]*regdata.TroopDialogueData
	Array []*regdata.TroopDialogueData

	MinKeyData *regdata.TroopDialogueData
	MaxKeyData *regdata.TroopDialogueData

	parserMap map[*regdata.TroopDialogueData]*config.ObjectParser
}

func (d *TroopDialogueDataConfig) Get(key uint64) *regdata.TroopDialogueData {
	return d.Map[key]
}

func (d *TroopDialogueDataConfig) Must(key uint64) *regdata.TroopDialogueData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTroopDialogueTextData(key uint64) *regdata.TroopDialogueTextData {
	return dAtA.troopDialogueTextData.Map[key]
}

func (dAtA *ConfigDatas) GetTroopDialogueTextDataArray() []*regdata.TroopDialogueTextData {
	return dAtA.troopDialogueTextData.Array
}

func (dAtA *ConfigDatas) TroopDialogueTextData() *TroopDialogueTextDataConfig {
	return dAtA.troopDialogueTextData
}

type TroopDialogueTextDataConfig struct {
	Map   map[uint64]*regdata.TroopDialogueTextData
	Array []*regdata.TroopDialogueTextData

	MinKeyData *regdata.TroopDialogueTextData
	MaxKeyData *regdata.TroopDialogueTextData

	parserMap map[*regdata.TroopDialogueTextData]*config.ObjectParser
}

func (d *TroopDialogueTextDataConfig) Get(key uint64) *regdata.TroopDialogueTextData {
	return d.Map[key]
}

func (d *TroopDialogueTextDataConfig) Must(key uint64) *regdata.TroopDialogueTextData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetAmountShowSortData(key uint64) *resdata.AmountShowSortData {
	return dAtA.amountShowSortData.Map[key]
}

func (dAtA *ConfigDatas) GetAmountShowSortDataArray() []*resdata.AmountShowSortData {
	return dAtA.amountShowSortData.Array
}

func (dAtA *ConfigDatas) AmountShowSortData() *AmountShowSortDataConfig {
	return dAtA.amountShowSortData
}

type AmountShowSortDataConfig struct {
	Map   map[uint64]*resdata.AmountShowSortData
	Array []*resdata.AmountShowSortData

	MinKeyData *resdata.AmountShowSortData
	MaxKeyData *resdata.AmountShowSortData

	parserMap map[*resdata.AmountShowSortData]*config.ObjectParser
}

func (d *AmountShowSortDataConfig) Get(key uint64) *resdata.AmountShowSortData {
	return d.Map[key]
}

func (d *AmountShowSortDataConfig) Must(key uint64) *resdata.AmountShowSortData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBaowuData(key uint64) *resdata.BaowuData {
	return dAtA.baowuData.Map[key]
}

func (dAtA *ConfigDatas) GetBaowuDataArray() []*resdata.BaowuData {
	return dAtA.baowuData.Array
}

func (dAtA *ConfigDatas) BaowuData() *BaowuDataConfig {
	return dAtA.baowuData
}

type BaowuDataConfig struct {
	Map   map[uint64]*resdata.BaowuData
	Array []*resdata.BaowuData

	MinKeyData *resdata.BaowuData
	MaxKeyData *resdata.BaowuData

	parserMap map[*resdata.BaowuData]*config.ObjectParser
}

func (d *BaowuDataConfig) Get(key uint64) *resdata.BaowuData {
	return d.Map[key]
}

func (d *BaowuDataConfig) Must(key uint64) *resdata.BaowuData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetConditionPlunder(key uint64) *resdata.ConditionPlunder {
	return dAtA.conditionPlunder.Map[key]
}

func (dAtA *ConfigDatas) GetConditionPlunderArray() []*resdata.ConditionPlunder {
	return dAtA.conditionPlunder.Array
}

func (dAtA *ConfigDatas) ConditionPlunder() *ConditionPlunderConfig {
	return dAtA.conditionPlunder
}

type ConditionPlunderConfig struct {
	Map   map[uint64]*resdata.ConditionPlunder
	Array []*resdata.ConditionPlunder

	MinKeyData *resdata.ConditionPlunder
	MaxKeyData *resdata.ConditionPlunder

	parserMap map[*resdata.ConditionPlunder]*config.ObjectParser
}

func (d *ConditionPlunderConfig) Get(key uint64) *resdata.ConditionPlunder {
	return d.Map[key]
}

func (d *ConditionPlunderConfig) Must(key uint64) *resdata.ConditionPlunder {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetConditionPlunderItem(key uint64) *resdata.ConditionPlunderItem {
	return dAtA.conditionPlunderItem.Map[key]
}

func (dAtA *ConfigDatas) GetConditionPlunderItemArray() []*resdata.ConditionPlunderItem {
	return dAtA.conditionPlunderItem.Array
}

func (dAtA *ConfigDatas) ConditionPlunderItem() *ConditionPlunderItemConfig {
	return dAtA.conditionPlunderItem
}

type ConditionPlunderItemConfig struct {
	Map   map[uint64]*resdata.ConditionPlunderItem
	Array []*resdata.ConditionPlunderItem

	MinKeyData *resdata.ConditionPlunderItem
	MaxKeyData *resdata.ConditionPlunderItem

	parserMap map[*resdata.ConditionPlunderItem]*config.ObjectParser
}

func (d *ConditionPlunderItemConfig) Get(key uint64) *resdata.ConditionPlunderItem {
	return d.Map[key]
}

func (d *ConditionPlunderItemConfig) Must(key uint64) *resdata.ConditionPlunderItem {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCost(key int) *resdata.Cost {
	return dAtA.cost.Map[key]
}

func (dAtA *ConfigDatas) GetCostArray() []*resdata.Cost {
	return dAtA.cost.Array
}

func (dAtA *ConfigDatas) Cost() *CostConfig {
	return dAtA.cost
}

type CostConfig struct {
	Map   map[int]*resdata.Cost
	Array []*resdata.Cost

	MinKeyData *resdata.Cost
	MaxKeyData *resdata.Cost

	parserMap map[*resdata.Cost]*config.ObjectParser
}

func (d *CostConfig) Get(key int) *resdata.Cost {
	return d.Map[key]
}

func (d *CostConfig) Must(key int) *resdata.Cost {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetGuildLevelPrize(key uint64) *resdata.GuildLevelPrize {
	return dAtA.guildLevelPrize.Map[key]
}

func (dAtA *ConfigDatas) GetGuildLevelPrizeArray() []*resdata.GuildLevelPrize {
	return dAtA.guildLevelPrize.Array
}

func (dAtA *ConfigDatas) GuildLevelPrize() *GuildLevelPrizeConfig {
	return dAtA.guildLevelPrize
}

type GuildLevelPrizeConfig struct {
	Map   map[uint64]*resdata.GuildLevelPrize
	Array []*resdata.GuildLevelPrize

	MinKeyData *resdata.GuildLevelPrize
	MaxKeyData *resdata.GuildLevelPrize

	parserMap map[*resdata.GuildLevelPrize]*config.ObjectParser
}

func (d *GuildLevelPrizeConfig) Get(key uint64) *resdata.GuildLevelPrize {
	return d.Map[key]
}

func (d *GuildLevelPrizeConfig) Must(key uint64) *resdata.GuildLevelPrize {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetPlunder(key uint64) *resdata.Plunder {
	return dAtA.plunder.Map[key]
}

func (dAtA *ConfigDatas) GetPlunderArray() []*resdata.Plunder {
	return dAtA.plunder.Array
}

func (dAtA *ConfigDatas) Plunder() *PlunderConfig {
	return dAtA.plunder
}

type PlunderConfig struct {
	Map   map[uint64]*resdata.Plunder
	Array []*resdata.Plunder

	MinKeyData *resdata.Plunder
	MaxKeyData *resdata.Plunder

	parserMap map[*resdata.Plunder]*config.ObjectParser
}

func (d *PlunderConfig) Get(key uint64) *resdata.Plunder {
	return d.Map[key]
}

func (d *PlunderConfig) Must(key uint64) *resdata.Plunder {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetPlunderGroup(key uint64) *resdata.PlunderGroup {
	return dAtA.plunderGroup.Map[key]
}

func (dAtA *ConfigDatas) GetPlunderGroupArray() []*resdata.PlunderGroup {
	return dAtA.plunderGroup.Array
}

func (dAtA *ConfigDatas) PlunderGroup() *PlunderGroupConfig {
	return dAtA.plunderGroup
}

type PlunderGroupConfig struct {
	Map   map[uint64]*resdata.PlunderGroup
	Array []*resdata.PlunderGroup

	MinKeyData *resdata.PlunderGroup
	MaxKeyData *resdata.PlunderGroup

	parserMap map[*resdata.PlunderGroup]*config.ObjectParser
}

func (d *PlunderGroupConfig) Get(key uint64) *resdata.PlunderGroup {
	return d.Map[key]
}

func (d *PlunderGroupConfig) Must(key uint64) *resdata.PlunderGroup {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetPlunderItem(key uint64) *resdata.PlunderItem {
	return dAtA.plunderItem.Map[key]
}

func (dAtA *ConfigDatas) GetPlunderItemArray() []*resdata.PlunderItem {
	return dAtA.plunderItem.Array
}

func (dAtA *ConfigDatas) PlunderItem() *PlunderItemConfig {
	return dAtA.plunderItem
}

type PlunderItemConfig struct {
	Map   map[uint64]*resdata.PlunderItem
	Array []*resdata.PlunderItem

	MinKeyData *resdata.PlunderItem
	MaxKeyData *resdata.PlunderItem

	parserMap map[*resdata.PlunderItem]*config.ObjectParser
}

func (d *PlunderItemConfig) Get(key uint64) *resdata.PlunderItem {
	return d.Map[key]
}

func (d *PlunderItemConfig) Must(key uint64) *resdata.PlunderItem {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetPlunderPrize(key uint64) *resdata.PlunderPrize {
	return dAtA.plunderPrize.Map[key]
}

func (dAtA *ConfigDatas) GetPlunderPrizeArray() []*resdata.PlunderPrize {
	return dAtA.plunderPrize.Array
}

func (dAtA *ConfigDatas) PlunderPrize() *PlunderPrizeConfig {
	return dAtA.plunderPrize
}

type PlunderPrizeConfig struct {
	Map   map[uint64]*resdata.PlunderPrize
	Array []*resdata.PlunderPrize

	MinKeyData *resdata.PlunderPrize
	MaxKeyData *resdata.PlunderPrize

	parserMap map[*resdata.PlunderPrize]*config.ObjectParser
}

func (d *PlunderPrizeConfig) Get(key uint64) *resdata.PlunderPrize {
	return d.Map[key]
}

func (d *PlunderPrizeConfig) Must(key uint64) *resdata.PlunderPrize {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetPrize(key int) *resdata.Prize {
	return dAtA.prize.Map[key]
}

func (dAtA *ConfigDatas) GetPrizeArray() []*resdata.Prize {
	return dAtA.prize.Array
}

func (dAtA *ConfigDatas) Prize() *PrizeConfig {
	return dAtA.prize
}

type PrizeConfig struct {
	Map   map[int]*resdata.Prize
	Array []*resdata.Prize

	MinKeyData *resdata.Prize
	MaxKeyData *resdata.Prize

	parserMap map[*resdata.Prize]*config.ObjectParser
}

func (d *PrizeConfig) Get(key int) *resdata.Prize {
	return d.Map[key]
}

func (d *PrizeConfig) Must(key int) *resdata.Prize {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetResCaptainData(key uint64) *resdata.ResCaptainData {
	return dAtA.resCaptainData.Map[key]
}

func (dAtA *ConfigDatas) GetResCaptainDataArray() []*resdata.ResCaptainData {
	return dAtA.resCaptainData.Array
}

func (dAtA *ConfigDatas) ResCaptainData() *ResCaptainDataConfig {
	return dAtA.resCaptainData
}

type ResCaptainDataConfig struct {
	Map   map[uint64]*resdata.ResCaptainData
	Array []*resdata.ResCaptainData

	MinKeyData *resdata.ResCaptainData
	MaxKeyData *resdata.ResCaptainData

	parserMap map[*resdata.ResCaptainData]*config.ObjectParser
}

func (d *ResCaptainDataConfig) Get(key uint64) *resdata.ResCaptainData {
	return d.Map[key]
}

func (d *ResCaptainDataConfig) Must(key uint64) *resdata.ResCaptainData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetCombatScene(key string) *scene.CombatScene {
	return dAtA.combatScene.Map[key]
}

func (dAtA *ConfigDatas) GetCombatSceneArray() []*scene.CombatScene {
	return dAtA.combatScene.Array
}

func (dAtA *ConfigDatas) CombatScene() *CombatSceneConfig {
	return dAtA.combatScene
}

type CombatSceneConfig struct {
	Map   map[string]*scene.CombatScene
	Array []*scene.CombatScene

	MinKeyData *scene.CombatScene
	MaxKeyData *scene.CombatScene

	parserMap map[*scene.CombatScene]*config.ObjectParser
}

func (d *CombatSceneConfig) Get(key string) *scene.CombatScene {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetSeasonData(key uint64) *season.SeasonData {
	return dAtA.seasonData.Map[key]
}

func (dAtA *ConfigDatas) GetSeasonDataArray() []*season.SeasonData {
	return dAtA.seasonData.Array
}

func (dAtA *ConfigDatas) SeasonData() *SeasonDataConfig {
	return dAtA.seasonData
}

type SeasonDataConfig struct {
	Map   map[uint64]*season.SeasonData
	Array []*season.SeasonData

	MinKeyData *season.SeasonData
	MaxKeyData *season.SeasonData

	parserMap map[*season.SeasonData]*config.ObjectParser
}

func (d *SeasonDataConfig) Get(key uint64) *season.SeasonData {
	return d.Map[key]
}

func (d *SeasonDataConfig) Must(key uint64) *season.SeasonData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetPrivacySettingData(key uint64) *settings.PrivacySettingData {
	return dAtA.privacySettingData.Map[key]
}

func (dAtA *ConfigDatas) GetPrivacySettingDataArray() []*settings.PrivacySettingData {
	return dAtA.privacySettingData.Array
}

func (dAtA *ConfigDatas) PrivacySettingData() *PrivacySettingDataConfig {
	return dAtA.privacySettingData
}

type PrivacySettingDataConfig struct {
	Map   map[uint64]*settings.PrivacySettingData
	Array []*settings.PrivacySettingData

	MinKeyData *settings.PrivacySettingData
	MaxKeyData *settings.PrivacySettingData

	parserMap map[*settings.PrivacySettingData]*config.ObjectParser
}

func (d *PrivacySettingDataConfig) Get(key uint64) *settings.PrivacySettingData {
	return d.Map[key]
}

func (d *PrivacySettingDataConfig) Must(key uint64) *settings.PrivacySettingData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBlackMarketData(key uint64) *shop.BlackMarketData {
	return dAtA.blackMarketData.Map[key]
}

func (dAtA *ConfigDatas) GetBlackMarketDataArray() []*shop.BlackMarketData {
	return dAtA.blackMarketData.Array
}

func (dAtA *ConfigDatas) BlackMarketData() *BlackMarketDataConfig {
	return dAtA.blackMarketData
}

type BlackMarketDataConfig struct {
	Map   map[uint64]*shop.BlackMarketData
	Array []*shop.BlackMarketData

	MinKeyData *shop.BlackMarketData
	MaxKeyData *shop.BlackMarketData

	parserMap map[*shop.BlackMarketData]*config.ObjectParser
}

func (d *BlackMarketDataConfig) Get(key uint64) *shop.BlackMarketData {
	return d.Map[key]
}

func (d *BlackMarketDataConfig) Must(key uint64) *shop.BlackMarketData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBlackMarketGoodsData(key uint64) *shop.BlackMarketGoodsData {
	return dAtA.blackMarketGoodsData.Map[key]
}

func (dAtA *ConfigDatas) GetBlackMarketGoodsDataArray() []*shop.BlackMarketGoodsData {
	return dAtA.blackMarketGoodsData.Array
}

func (dAtA *ConfigDatas) BlackMarketGoodsData() *BlackMarketGoodsDataConfig {
	return dAtA.blackMarketGoodsData
}

type BlackMarketGoodsDataConfig struct {
	Map   map[uint64]*shop.BlackMarketGoodsData
	Array []*shop.BlackMarketGoodsData

	MinKeyData *shop.BlackMarketGoodsData
	MaxKeyData *shop.BlackMarketGoodsData

	parserMap map[*shop.BlackMarketGoodsData]*config.ObjectParser
}

func (d *BlackMarketGoodsDataConfig) Get(key uint64) *shop.BlackMarketGoodsData {
	return d.Map[key]
}

func (d *BlackMarketGoodsDataConfig) Must(key uint64) *shop.BlackMarketGoodsData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBlackMarketGoodsGroupData(key uint64) *shop.BlackMarketGoodsGroupData {
	return dAtA.blackMarketGoodsGroupData.Map[key]
}

func (dAtA *ConfigDatas) GetBlackMarketGoodsGroupDataArray() []*shop.BlackMarketGoodsGroupData {
	return dAtA.blackMarketGoodsGroupData.Array
}

func (dAtA *ConfigDatas) BlackMarketGoodsGroupData() *BlackMarketGoodsGroupDataConfig {
	return dAtA.blackMarketGoodsGroupData
}

type BlackMarketGoodsGroupDataConfig struct {
	Map   map[uint64]*shop.BlackMarketGoodsGroupData
	Array []*shop.BlackMarketGoodsGroupData

	MinKeyData *shop.BlackMarketGoodsGroupData
	MaxKeyData *shop.BlackMarketGoodsGroupData

	parserMap map[*shop.BlackMarketGoodsGroupData]*config.ObjectParser
}

func (d *BlackMarketGoodsGroupDataConfig) Get(key uint64) *shop.BlackMarketGoodsGroupData {
	return d.Map[key]
}

func (d *BlackMarketGoodsGroupDataConfig) Must(key uint64) *shop.BlackMarketGoodsGroupData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetDiscountColorData(key uint64) *shop.DiscountColorData {
	return dAtA.discountColorData.Map[key]
}

func (dAtA *ConfigDatas) GetDiscountColorDataArray() []*shop.DiscountColorData {
	return dAtA.discountColorData.Array
}

func (dAtA *ConfigDatas) DiscountColorData() *DiscountColorDataConfig {
	return dAtA.discountColorData
}

type DiscountColorDataConfig struct {
	Map   map[uint64]*shop.DiscountColorData
	Array []*shop.DiscountColorData

	MinKeyData *shop.DiscountColorData
	MaxKeyData *shop.DiscountColorData

	parserMap map[*shop.DiscountColorData]*config.ObjectParser
}

func (d *DiscountColorDataConfig) Get(key uint64) *shop.DiscountColorData {
	return d.Map[key]
}

func (d *DiscountColorDataConfig) Must(key uint64) *shop.DiscountColorData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetNormalShopGoods(key uint64) *shop.NormalShopGoods {
	return dAtA.normalShopGoods.Map[key]
}

func (dAtA *ConfigDatas) GetNormalShopGoodsArray() []*shop.NormalShopGoods {
	return dAtA.normalShopGoods.Array
}

func (dAtA *ConfigDatas) NormalShopGoods() *NormalShopGoodsConfig {
	return dAtA.normalShopGoods
}

type NormalShopGoodsConfig struct {
	Map   map[uint64]*shop.NormalShopGoods
	Array []*shop.NormalShopGoods

	MinKeyData *shop.NormalShopGoods
	MaxKeyData *shop.NormalShopGoods

	parserMap map[*shop.NormalShopGoods]*config.ObjectParser
}

func (d *NormalShopGoodsConfig) Get(key uint64) *shop.NormalShopGoods {
	return d.Map[key]
}

func (d *NormalShopGoodsConfig) Must(key uint64) *shop.NormalShopGoods {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetShop(key uint64) *shop.Shop {
	return dAtA.shop.Map[key]
}

func (dAtA *ConfigDatas) GetShopArray() []*shop.Shop {
	return dAtA.shop.Array
}

func (dAtA *ConfigDatas) Shop() *ShopConfig {
	return dAtA.shop
}

type ShopConfig struct {
	Map   map[uint64]*shop.Shop
	Array []*shop.Shop

	MinKeyData *shop.Shop
	MaxKeyData *shop.Shop

	parserMap map[*shop.Shop]*config.ObjectParser
}

func (d *ShopConfig) Get(key uint64) *shop.Shop {
	return d.Map[key]
}

func (d *ShopConfig) Must(key uint64) *shop.Shop {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetZhenBaoGeShopGoods(key uint64) *shop.ZhenBaoGeShopGoods {
	return dAtA.zhenBaoGeShopGoods.Map[key]
}

func (dAtA *ConfigDatas) GetZhenBaoGeShopGoodsArray() []*shop.ZhenBaoGeShopGoods {
	return dAtA.zhenBaoGeShopGoods.Array
}

func (dAtA *ConfigDatas) ZhenBaoGeShopGoods() *ZhenBaoGeShopGoodsConfig {
	return dAtA.zhenBaoGeShopGoods
}

type ZhenBaoGeShopGoodsConfig struct {
	Map   map[uint64]*shop.ZhenBaoGeShopGoods
	Array []*shop.ZhenBaoGeShopGoods

	MinKeyData *shop.ZhenBaoGeShopGoods
	MaxKeyData *shop.ZhenBaoGeShopGoods

	parserMap map[*shop.ZhenBaoGeShopGoods]*config.ObjectParser
}

func (d *ZhenBaoGeShopGoodsConfig) Get(key uint64) *shop.ZhenBaoGeShopGoods {
	return d.Map[key]
}

func (d *ZhenBaoGeShopGoodsConfig) Must(key uint64) *shop.ZhenBaoGeShopGoods {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetPassiveSpellData(key uint64) *spell.PassiveSpellData {
	return dAtA.passiveSpellData.Map[key]
}

func (dAtA *ConfigDatas) GetPassiveSpellDataArray() []*spell.PassiveSpellData {
	return dAtA.passiveSpellData.Array
}

func (dAtA *ConfigDatas) PassiveSpellData() *PassiveSpellDataConfig {
	return dAtA.passiveSpellData
}

type PassiveSpellDataConfig struct {
	Map   map[uint64]*spell.PassiveSpellData
	Array []*spell.PassiveSpellData

	MinKeyData *spell.PassiveSpellData
	MaxKeyData *spell.PassiveSpellData

	parserMap map[*spell.PassiveSpellData]*config.ObjectParser
}

func (d *PassiveSpellDataConfig) Get(key uint64) *spell.PassiveSpellData {
	return d.Map[key]
}

func (d *PassiveSpellDataConfig) Must(key uint64) *spell.PassiveSpellData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetSpell(key uint64) *spell.Spell {
	return dAtA.spell.Map[key]
}

func (dAtA *ConfigDatas) GetSpellArray() []*spell.Spell {
	return dAtA.spell.Array
}

func (dAtA *ConfigDatas) Spell() *SpellConfig {
	return dAtA.spell
}

type SpellConfig struct {
	Map   map[uint64]*spell.Spell
	Array []*spell.Spell

	MinKeyData *spell.Spell
	MaxKeyData *spell.Spell

	parserMap map[*spell.Spell]*config.ObjectParser
}

func (d *SpellConfig) Get(key uint64) *spell.Spell {
	return d.Map[key]
}

func (d *SpellConfig) Must(key uint64) *spell.Spell {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetSpellData(key uint64) *spell.SpellData {
	return dAtA.spellData.Map[key]
}

func (dAtA *ConfigDatas) GetSpellDataArray() []*spell.SpellData {
	return dAtA.spellData.Array
}

func (dAtA *ConfigDatas) SpellData() *SpellDataConfig {
	return dAtA.spellData
}

type SpellDataConfig struct {
	Map   map[uint64]*spell.SpellData
	Array []*spell.SpellData

	MinKeyData *spell.SpellData
	MaxKeyData *spell.SpellData

	parserMap map[*spell.SpellData]*config.ObjectParser
}

func (d *SpellDataConfig) Get(key uint64) *spell.SpellData {
	return d.Map[key]
}

func (d *SpellDataConfig) Must(key uint64) *spell.SpellData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetSpellFacadeData(key uint64) *spell.SpellFacadeData {
	return dAtA.spellFacadeData.Map[key]
}

func (dAtA *ConfigDatas) GetSpellFacadeDataArray() []*spell.SpellFacadeData {
	return dAtA.spellFacadeData.Array
}

func (dAtA *ConfigDatas) SpellFacadeData() *SpellFacadeDataConfig {
	return dAtA.spellFacadeData
}

type SpellFacadeDataConfig struct {
	Map   map[uint64]*spell.SpellFacadeData
	Array []*spell.SpellFacadeData

	MinKeyData *spell.SpellFacadeData
	MaxKeyData *spell.SpellFacadeData

	parserMap map[*spell.SpellFacadeData]*config.ObjectParser
}

func (d *SpellFacadeDataConfig) Get(key uint64) *spell.SpellFacadeData {
	return d.Map[key]
}

func (d *SpellFacadeDataConfig) Must(key uint64) *spell.SpellFacadeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetStateData(key uint64) *spell.StateData {
	return dAtA.stateData.Map[key]
}

func (dAtA *ConfigDatas) GetStateDataArray() []*spell.StateData {
	return dAtA.stateData.Array
}

func (dAtA *ConfigDatas) StateData() *StateDataConfig {
	return dAtA.stateData
}

type StateDataConfig struct {
	Map   map[uint64]*spell.StateData
	Array []*spell.StateData

	MinKeyData *spell.StateData
	MaxKeyData *spell.StateData

	parserMap map[*spell.StateData]*config.ObjectParser
}

func (d *StateDataConfig) Get(key uint64) *spell.StateData {
	return d.Map[key]
}

func (d *StateDataConfig) Must(key uint64) *spell.StateData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetStrategyData(key uint64) *strategydata.StrategyData {
	return dAtA.strategyData.Map[key]
}

func (dAtA *ConfigDatas) GetStrategyDataArray() []*strategydata.StrategyData {
	return dAtA.strategyData.Array
}

func (dAtA *ConfigDatas) StrategyData() *StrategyDataConfig {
	return dAtA.strategyData
}

type StrategyDataConfig struct {
	Map   map[uint64]*strategydata.StrategyData
	Array []*strategydata.StrategyData

	MinKeyData *strategydata.StrategyData
	MaxKeyData *strategydata.StrategyData

	parserMap map[*strategydata.StrategyData]*config.ObjectParser
}

func (d *StrategyDataConfig) Get(key uint64) *strategydata.StrategyData {
	return d.Map[key]
}

func (d *StrategyDataConfig) Must(key uint64) *strategydata.StrategyData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetStrategyEffectData(key uint64) *strategydata.StrategyEffectData {
	return dAtA.strategyEffectData.Map[key]
}

func (dAtA *ConfigDatas) GetStrategyEffectDataArray() []*strategydata.StrategyEffectData {
	return dAtA.strategyEffectData.Array
}

func (dAtA *ConfigDatas) StrategyEffectData() *StrategyEffectDataConfig {
	return dAtA.strategyEffectData
}

type StrategyEffectDataConfig struct {
	Map   map[uint64]*strategydata.StrategyEffectData
	Array []*strategydata.StrategyEffectData

	MinKeyData *strategydata.StrategyEffectData
	MaxKeyData *strategydata.StrategyEffectData

	parserMap map[*strategydata.StrategyEffectData]*config.ObjectParser
}

func (d *StrategyEffectDataConfig) Get(key uint64) *strategydata.StrategyEffectData {
	return d.Map[key]
}

func (d *StrategyEffectDataConfig) Must(key uint64) *strategydata.StrategyEffectData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetStrongerData(key uint64) *strongerdata.StrongerData {
	return dAtA.strongerData.Map[key]
}

func (dAtA *ConfigDatas) GetStrongerDataArray() []*strongerdata.StrongerData {
	return dAtA.strongerData.Array
}

func (dAtA *ConfigDatas) StrongerData() *StrongerDataConfig {
	return dAtA.strongerData
}

type StrongerDataConfig struct {
	Map   map[uint64]*strongerdata.StrongerData
	Array []*strongerdata.StrongerData

	MinKeyData *strongerdata.StrongerData
	MaxKeyData *strongerdata.StrongerData

	parserMap map[*strongerdata.StrongerData]*config.ObjectParser
}

func (d *StrongerDataConfig) Get(key uint64) *strongerdata.StrongerData {
	return d.Map[key]
}

func (d *StrongerDataConfig) Must(key uint64) *strongerdata.StrongerData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBuildingEffectData(key int) *sub.BuildingEffectData {
	return dAtA.buildingEffectData.Map[key]
}

func (dAtA *ConfigDatas) GetBuildingEffectDataArray() []*sub.BuildingEffectData {
	return dAtA.buildingEffectData.Array
}

func (dAtA *ConfigDatas) BuildingEffectData() *BuildingEffectDataConfig {
	return dAtA.buildingEffectData
}

type BuildingEffectDataConfig struct {
	Map   map[int]*sub.BuildingEffectData
	Array []*sub.BuildingEffectData

	MinKeyData *sub.BuildingEffectData
	MaxKeyData *sub.BuildingEffectData

	parserMap map[*sub.BuildingEffectData]*config.ObjectParser
}

func (d *BuildingEffectDataConfig) Get(key int) *sub.BuildingEffectData {
	return d.Map[key]
}

func (d *BuildingEffectDataConfig) Must(key int) *sub.BuildingEffectData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetSurveyData(key string) *survey.SurveyData {
	return dAtA.surveyData.Map[key]
}

func (dAtA *ConfigDatas) GetSurveyDataArray() []*survey.SurveyData {
	return dAtA.surveyData.Array
}

func (dAtA *ConfigDatas) SurveyData() *SurveyDataConfig {
	return dAtA.surveyData
}

type SurveyDataConfig struct {
	Map   map[string]*survey.SurveyData
	Array []*survey.SurveyData

	MinKeyData *survey.SurveyData
	MaxKeyData *survey.SurveyData

	parserMap map[*survey.SurveyData]*config.ObjectParser
}

func (d *SurveyDataConfig) Get(key string) *survey.SurveyData {
	return d.Map[key]
}

func (dAtA *ConfigDatas) GetAchieveTaskData(key uint64) *taskdata.AchieveTaskData {
	return dAtA.achieveTaskData.Map[key]
}

func (dAtA *ConfigDatas) GetAchieveTaskDataArray() []*taskdata.AchieveTaskData {
	return dAtA.achieveTaskData.Array
}

func (dAtA *ConfigDatas) AchieveTaskData() *AchieveTaskDataConfig {
	return dAtA.achieveTaskData
}

type AchieveTaskDataConfig struct {
	Map   map[uint64]*taskdata.AchieveTaskData
	Array []*taskdata.AchieveTaskData

	MinKeyData *taskdata.AchieveTaskData
	MaxKeyData *taskdata.AchieveTaskData

	parserMap map[*taskdata.AchieveTaskData]*config.ObjectParser
}

func (d *AchieveTaskDataConfig) Get(key uint64) *taskdata.AchieveTaskData {
	return d.Map[key]
}

func (d *AchieveTaskDataConfig) Must(key uint64) *taskdata.AchieveTaskData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetAchieveTaskStarPrizeData(key uint64) *taskdata.AchieveTaskStarPrizeData {
	return dAtA.achieveTaskStarPrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetAchieveTaskStarPrizeDataArray() []*taskdata.AchieveTaskStarPrizeData {
	return dAtA.achieveTaskStarPrizeData.Array
}

func (dAtA *ConfigDatas) AchieveTaskStarPrizeData() *AchieveTaskStarPrizeDataConfig {
	return dAtA.achieveTaskStarPrizeData
}

type AchieveTaskStarPrizeDataConfig struct {
	Map   map[uint64]*taskdata.AchieveTaskStarPrizeData
	Array []*taskdata.AchieveTaskStarPrizeData

	MinKeyData *taskdata.AchieveTaskStarPrizeData
	MaxKeyData *taskdata.AchieveTaskStarPrizeData

	parserMap map[*taskdata.AchieveTaskStarPrizeData]*config.ObjectParser
}

func (d *AchieveTaskStarPrizeDataConfig) Get(key uint64) *taskdata.AchieveTaskStarPrizeData {
	return d.Map[key]
}

func (d *AchieveTaskStarPrizeDataConfig) Must(key uint64) *taskdata.AchieveTaskStarPrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetActiveDegreePrizeData(key uint64) *taskdata.ActiveDegreePrizeData {
	return dAtA.activeDegreePrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetActiveDegreePrizeDataArray() []*taskdata.ActiveDegreePrizeData {
	return dAtA.activeDegreePrizeData.Array
}

func (dAtA *ConfigDatas) ActiveDegreePrizeData() *ActiveDegreePrizeDataConfig {
	return dAtA.activeDegreePrizeData
}

type ActiveDegreePrizeDataConfig struct {
	Map   map[uint64]*taskdata.ActiveDegreePrizeData
	Array []*taskdata.ActiveDegreePrizeData

	MinKeyData *taskdata.ActiveDegreePrizeData
	MaxKeyData *taskdata.ActiveDegreePrizeData

	parserMap map[*taskdata.ActiveDegreePrizeData]*config.ObjectParser
}

func (d *ActiveDegreePrizeDataConfig) Get(key uint64) *taskdata.ActiveDegreePrizeData {
	return d.Map[key]
}

func (d *ActiveDegreePrizeDataConfig) Must(key uint64) *taskdata.ActiveDegreePrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetActiveDegreeTaskData(key uint64) *taskdata.ActiveDegreeTaskData {
	return dAtA.activeDegreeTaskData.Map[key]
}

func (dAtA *ConfigDatas) GetActiveDegreeTaskDataArray() []*taskdata.ActiveDegreeTaskData {
	return dAtA.activeDegreeTaskData.Array
}

func (dAtA *ConfigDatas) ActiveDegreeTaskData() *ActiveDegreeTaskDataConfig {
	return dAtA.activeDegreeTaskData
}

type ActiveDegreeTaskDataConfig struct {
	Map   map[uint64]*taskdata.ActiveDegreeTaskData
	Array []*taskdata.ActiveDegreeTaskData

	MinKeyData *taskdata.ActiveDegreeTaskData
	MaxKeyData *taskdata.ActiveDegreeTaskData

	parserMap map[*taskdata.ActiveDegreeTaskData]*config.ObjectParser
}

func (d *ActiveDegreeTaskDataConfig) Get(key uint64) *taskdata.ActiveDegreeTaskData {
	return d.Map[key]
}

func (d *ActiveDegreeTaskDataConfig) Must(key uint64) *taskdata.ActiveDegreeTaskData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetActivityTaskData(key uint64) *taskdata.ActivityTaskData {
	return dAtA.activityTaskData.Map[key]
}

func (dAtA *ConfigDatas) GetActivityTaskDataArray() []*taskdata.ActivityTaskData {
	return dAtA.activityTaskData.Array
}

func (dAtA *ConfigDatas) ActivityTaskData() *ActivityTaskDataConfig {
	return dAtA.activityTaskData
}

type ActivityTaskDataConfig struct {
	Map   map[uint64]*taskdata.ActivityTaskData
	Array []*taskdata.ActivityTaskData

	MinKeyData *taskdata.ActivityTaskData
	MaxKeyData *taskdata.ActivityTaskData

	parserMap map[*taskdata.ActivityTaskData]*config.ObjectParser
}

func (d *ActivityTaskDataConfig) Get(key uint64) *taskdata.ActivityTaskData {
	return d.Map[key]
}

func (d *ActivityTaskDataConfig) Must(key uint64) *taskdata.ActivityTaskData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBaYeStageData(key uint64) *taskdata.BaYeStageData {
	return dAtA.baYeStageData.Map[key]
}

func (dAtA *ConfigDatas) GetBaYeStageDataArray() []*taskdata.BaYeStageData {
	return dAtA.baYeStageData.Array
}

func (dAtA *ConfigDatas) BaYeStageData() *BaYeStageDataConfig {
	return dAtA.baYeStageData
}

type BaYeStageDataConfig struct {
	Map   map[uint64]*taskdata.BaYeStageData
	Array []*taskdata.BaYeStageData

	MinKeyData *taskdata.BaYeStageData
	MaxKeyData *taskdata.BaYeStageData

	parserMap map[*taskdata.BaYeStageData]*config.ObjectParser
}

func (d *BaYeStageDataConfig) Get(key uint64) *taskdata.BaYeStageData {
	return d.Map[key]
}

func (d *BaYeStageDataConfig) Must(key uint64) *taskdata.BaYeStageData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBaYeTaskData(key uint64) *taskdata.BaYeTaskData {
	return dAtA.baYeTaskData.Map[key]
}

func (dAtA *ConfigDatas) GetBaYeTaskDataArray() []*taskdata.BaYeTaskData {
	return dAtA.baYeTaskData.Array
}

func (dAtA *ConfigDatas) BaYeTaskData() *BaYeTaskDataConfig {
	return dAtA.baYeTaskData
}

type BaYeTaskDataConfig struct {
	Map   map[uint64]*taskdata.BaYeTaskData
	Array []*taskdata.BaYeTaskData

	MinKeyData *taskdata.BaYeTaskData
	MaxKeyData *taskdata.BaYeTaskData

	parserMap map[*taskdata.BaYeTaskData]*config.ObjectParser
}

func (d *BaYeTaskDataConfig) Get(key uint64) *taskdata.BaYeTaskData {
	return d.Map[key]
}

func (d *BaYeTaskDataConfig) Must(key uint64) *taskdata.BaYeTaskData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBranchTaskData(key uint64) *taskdata.BranchTaskData {
	return dAtA.branchTaskData.Map[key]
}

func (dAtA *ConfigDatas) GetBranchTaskDataArray() []*taskdata.BranchTaskData {
	return dAtA.branchTaskData.Array
}

func (dAtA *ConfigDatas) BranchTaskData() *BranchTaskDataConfig {
	return dAtA.branchTaskData
}

type BranchTaskDataConfig struct {
	Map   map[uint64]*taskdata.BranchTaskData
	Array []*taskdata.BranchTaskData

	MinKeyData *taskdata.BranchTaskData
	MaxKeyData *taskdata.BranchTaskData

	parserMap map[*taskdata.BranchTaskData]*config.ObjectParser
}

func (d *BranchTaskDataConfig) Get(key uint64) *taskdata.BranchTaskData {
	return d.Map[key]
}

func (d *BranchTaskDataConfig) Must(key uint64) *taskdata.BranchTaskData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBwzlPrizeData(key uint64) *taskdata.BwzlPrizeData {
	return dAtA.bwzlPrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetBwzlPrizeDataArray() []*taskdata.BwzlPrizeData {
	return dAtA.bwzlPrizeData.Array
}

func (dAtA *ConfigDatas) BwzlPrizeData() *BwzlPrizeDataConfig {
	return dAtA.bwzlPrizeData
}

type BwzlPrizeDataConfig struct {
	Map   map[uint64]*taskdata.BwzlPrizeData
	Array []*taskdata.BwzlPrizeData

	MinKeyData *taskdata.BwzlPrizeData
	MaxKeyData *taskdata.BwzlPrizeData

	parserMap map[*taskdata.BwzlPrizeData]*config.ObjectParser
}

func (d *BwzlPrizeDataConfig) Get(key uint64) *taskdata.BwzlPrizeData {
	return d.Map[key]
}

func (d *BwzlPrizeDataConfig) Must(key uint64) *taskdata.BwzlPrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetBwzlTaskData(key uint64) *taskdata.BwzlTaskData {
	return dAtA.bwzlTaskData.Map[key]
}

func (dAtA *ConfigDatas) GetBwzlTaskDataArray() []*taskdata.BwzlTaskData {
	return dAtA.bwzlTaskData.Array
}

func (dAtA *ConfigDatas) BwzlTaskData() *BwzlTaskDataConfig {
	return dAtA.bwzlTaskData
}

type BwzlTaskDataConfig struct {
	Map   map[uint64]*taskdata.BwzlTaskData
	Array []*taskdata.BwzlTaskData

	MinKeyData *taskdata.BwzlTaskData
	MaxKeyData *taskdata.BwzlTaskData

	parserMap map[*taskdata.BwzlTaskData]*config.ObjectParser
}

func (d *BwzlTaskDataConfig) Get(key uint64) *taskdata.BwzlTaskData {
	return d.Map[key]
}

func (d *BwzlTaskDataConfig) Must(key uint64) *taskdata.BwzlTaskData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetMainTaskData(key uint64) *taskdata.MainTaskData {
	return dAtA.mainTaskData.Map[key]
}

func (dAtA *ConfigDatas) GetMainTaskDataArray() []*taskdata.MainTaskData {
	return dAtA.mainTaskData.Array
}

func (dAtA *ConfigDatas) MainTaskData() *MainTaskDataConfig {
	return dAtA.mainTaskData
}

type MainTaskDataConfig struct {
	Map   map[uint64]*taskdata.MainTaskData
	Array []*taskdata.MainTaskData

	MinKeyData *taskdata.MainTaskData
	MaxKeyData *taskdata.MainTaskData

	parserMap map[*taskdata.MainTaskData]*config.ObjectParser
}

func (d *MainTaskDataConfig) Get(key uint64) *taskdata.MainTaskData {
	return d.Map[key]
}

func (d *MainTaskDataConfig) Must(key uint64) *taskdata.MainTaskData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTaskBoxData(key uint64) *taskdata.TaskBoxData {
	return dAtA.taskBoxData.Map[key]
}

func (dAtA *ConfigDatas) GetTaskBoxDataArray() []*taskdata.TaskBoxData {
	return dAtA.taskBoxData.Array
}

func (dAtA *ConfigDatas) TaskBoxData() *TaskBoxDataConfig {
	return dAtA.taskBoxData
}

type TaskBoxDataConfig struct {
	Map   map[uint64]*taskdata.TaskBoxData
	Array []*taskdata.TaskBoxData

	MinKeyData *taskdata.TaskBoxData
	MaxKeyData *taskdata.TaskBoxData

	parserMap map[*taskdata.TaskBoxData]*config.ObjectParser
}

func (d *TaskBoxDataConfig) Get(key uint64) *taskdata.TaskBoxData {
	return d.Map[key]
}

func (d *TaskBoxDataConfig) Must(key uint64) *taskdata.TaskBoxData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTaskTargetData(key uint64) *taskdata.TaskTargetData {
	return dAtA.taskTargetData.Map[key]
}

func (dAtA *ConfigDatas) GetTaskTargetDataArray() []*taskdata.TaskTargetData {
	return dAtA.taskTargetData.Array
}

func (dAtA *ConfigDatas) TaskTargetData() *TaskTargetDataConfig {
	return dAtA.taskTargetData
}

type TaskTargetDataConfig struct {
	Map   map[uint64]*taskdata.TaskTargetData
	Array []*taskdata.TaskTargetData

	MinKeyData *taskdata.TaskTargetData
	MaxKeyData *taskdata.TaskTargetData

	parserMap map[*taskdata.TaskTargetData]*config.ObjectParser
}

func (d *TaskTargetDataConfig) Get(key uint64) *taskdata.TaskTargetData {
	return d.Map[key]
}

func (d *TaskTargetDataConfig) Must(key uint64) *taskdata.TaskTargetData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTitleData(key uint64) *taskdata.TitleData {
	return dAtA.titleData.Map[key]
}

func (dAtA *ConfigDatas) GetTitleDataArray() []*taskdata.TitleData {
	return dAtA.titleData.Array
}

func (dAtA *ConfigDatas) TitleData() *TitleDataConfig {
	return dAtA.titleData
}

type TitleDataConfig struct {
	Map   map[uint64]*taskdata.TitleData
	Array []*taskdata.TitleData

	MinKeyData *taskdata.TitleData
	MaxKeyData *taskdata.TitleData

	parserMap map[*taskdata.TitleData]*config.ObjectParser
}

func (d *TitleDataConfig) Get(key uint64) *taskdata.TitleData {
	return d.Map[key]
}

func (d *TitleDataConfig) Must(key uint64) *taskdata.TitleData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTitleTaskData(key uint64) *taskdata.TitleTaskData {
	return dAtA.titleTaskData.Map[key]
}

func (dAtA *ConfigDatas) GetTitleTaskDataArray() []*taskdata.TitleTaskData {
	return dAtA.titleTaskData.Array
}

func (dAtA *ConfigDatas) TitleTaskData() *TitleTaskDataConfig {
	return dAtA.titleTaskData
}

type TitleTaskDataConfig struct {
	Map   map[uint64]*taskdata.TitleTaskData
	Array []*taskdata.TitleTaskData

	MinKeyData *taskdata.TitleTaskData
	MaxKeyData *taskdata.TitleTaskData

	parserMap map[*taskdata.TitleTaskData]*config.ObjectParser
}

func (d *TitleTaskDataConfig) Get(key uint64) *taskdata.TitleTaskData {
	return d.Map[key]
}

func (d *TitleTaskDataConfig) Must(key uint64) *taskdata.TitleTaskData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTeachChapterData(key uint64) *teach.TeachChapterData {
	return dAtA.teachChapterData.Map[key]
}

func (dAtA *ConfigDatas) GetTeachChapterDataArray() []*teach.TeachChapterData {
	return dAtA.teachChapterData.Array
}

func (dAtA *ConfigDatas) TeachChapterData() *TeachChapterDataConfig {
	return dAtA.teachChapterData
}

type TeachChapterDataConfig struct {
	Map   map[uint64]*teach.TeachChapterData
	Array []*teach.TeachChapterData

	MinKeyData *teach.TeachChapterData
	MaxKeyData *teach.TeachChapterData

	parserMap map[*teach.TeachChapterData]*config.ObjectParser
}

func (d *TeachChapterDataConfig) Get(key uint64) *teach.TeachChapterData {
	return d.Map[key]
}

func (d *TeachChapterDataConfig) Must(key uint64) *teach.TeachChapterData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetSecretTowerData(key uint64) *towerdata.SecretTowerData {
	return dAtA.secretTowerData.Map[key]
}

func (dAtA *ConfigDatas) GetSecretTowerDataArray() []*towerdata.SecretTowerData {
	return dAtA.secretTowerData.Array
}

func (dAtA *ConfigDatas) SecretTowerData() *SecretTowerDataConfig {
	return dAtA.secretTowerData
}

type SecretTowerDataConfig struct {
	Map   map[uint64]*towerdata.SecretTowerData
	Array []*towerdata.SecretTowerData

	MinKeyData *towerdata.SecretTowerData
	MaxKeyData *towerdata.SecretTowerData

	parserMap map[*towerdata.SecretTowerData]*config.ObjectParser
}

func (d *SecretTowerDataConfig) Get(key uint64) *towerdata.SecretTowerData {
	return d.Map[key]
}

func (d *SecretTowerDataConfig) Must(key uint64) *towerdata.SecretTowerData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetSecretTowerWordsData(key uint64) *towerdata.SecretTowerWordsData {
	return dAtA.secretTowerWordsData.Map[key]
}

func (dAtA *ConfigDatas) GetSecretTowerWordsDataArray() []*towerdata.SecretTowerWordsData {
	return dAtA.secretTowerWordsData.Array
}

func (dAtA *ConfigDatas) SecretTowerWordsData() *SecretTowerWordsDataConfig {
	return dAtA.secretTowerWordsData
}

type SecretTowerWordsDataConfig struct {
	Map   map[uint64]*towerdata.SecretTowerWordsData
	Array []*towerdata.SecretTowerWordsData

	MinKeyData *towerdata.SecretTowerWordsData
	MaxKeyData *towerdata.SecretTowerWordsData

	parserMap map[*towerdata.SecretTowerWordsData]*config.ObjectParser
}

func (d *SecretTowerWordsDataConfig) Get(key uint64) *towerdata.SecretTowerWordsData {
	return d.Map[key]
}

func (d *SecretTowerWordsDataConfig) Must(key uint64) *towerdata.SecretTowerWordsData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetTowerData(key uint64) *towerdata.TowerData {
	return dAtA.towerData.Map[key]
}

func (dAtA *ConfigDatas) GetTowerDataArray() []*towerdata.TowerData {
	return dAtA.towerData.Array
}

func (dAtA *ConfigDatas) TowerData() *TowerDataConfig {
	return dAtA.towerData
}

type TowerDataConfig struct {
	Map   map[uint64]*towerdata.TowerData
	Array []*towerdata.TowerData

	MinKeyData *towerdata.TowerData
	MaxKeyData *towerdata.TowerData

	parserMap map[*towerdata.TowerData]*config.ObjectParser
}

func (d *TowerDataConfig) Get(key uint64) *towerdata.TowerData {
	return d.Map[key]
}

func (d *TowerDataConfig) Must(key uint64) *towerdata.TowerData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetVipContinueDaysData(key uint64) *vip.VipContinueDaysData {
	return dAtA.vipContinueDaysData.Map[key]
}

func (dAtA *ConfigDatas) GetVipContinueDaysDataArray() []*vip.VipContinueDaysData {
	return dAtA.vipContinueDaysData.Array
}

func (dAtA *ConfigDatas) VipContinueDaysData() *VipContinueDaysDataConfig {
	return dAtA.vipContinueDaysData
}

type VipContinueDaysDataConfig struct {
	Map   map[uint64]*vip.VipContinueDaysData
	Array []*vip.VipContinueDaysData

	MinKeyData *vip.VipContinueDaysData
	MaxKeyData *vip.VipContinueDaysData

	parserMap map[*vip.VipContinueDaysData]*config.ObjectParser
}

func (d *VipContinueDaysDataConfig) Get(key uint64) *vip.VipContinueDaysData {
	return d.Map[key]
}

func (d *VipContinueDaysDataConfig) Must(key uint64) *vip.VipContinueDaysData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetVipLevelData(key uint64) *vip.VipLevelData {
	return dAtA.vipLevelData.Map[key]
}

func (dAtA *ConfigDatas) GetVipLevelDataArray() []*vip.VipLevelData {
	return dAtA.vipLevelData.Array
}

func (dAtA *ConfigDatas) VipLevelData() *VipLevelDataConfig {
	return dAtA.vipLevelData
}

type VipLevelDataConfig struct {
	Map   map[uint64]*vip.VipLevelData
	Array []*vip.VipLevelData

	MinKeyData *vip.VipLevelData
	MaxKeyData *vip.VipLevelData

	parserMap map[*vip.VipLevelData]*config.ObjectParser
}

func (d *VipLevelDataConfig) Get(key uint64) *vip.VipLevelData {
	return d.Map[key]
}

func (d *VipLevelDataConfig) Must(key uint64) *vip.VipLevelData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetResistXiongNuData(key uint64) *xiongnu.ResistXiongNuData {
	return dAtA.resistXiongNuData.Map[key]
}

func (dAtA *ConfigDatas) GetResistXiongNuDataArray() []*xiongnu.ResistXiongNuData {
	return dAtA.resistXiongNuData.Array
}

func (dAtA *ConfigDatas) ResistXiongNuData() *ResistXiongNuDataConfig {
	return dAtA.resistXiongNuData
}

type ResistXiongNuDataConfig struct {
	Map   map[uint64]*xiongnu.ResistXiongNuData
	Array []*xiongnu.ResistXiongNuData

	MinKeyData *xiongnu.ResistXiongNuData
	MaxKeyData *xiongnu.ResistXiongNuData

	parserMap map[*xiongnu.ResistXiongNuData]*config.ObjectParser
}

func (d *ResistXiongNuDataConfig) Get(key uint64) *xiongnu.ResistXiongNuData {
	return d.Map[key]
}

func (d *ResistXiongNuDataConfig) Must(key uint64) *xiongnu.ResistXiongNuData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetResistXiongNuScoreData(key uint64) *xiongnu.ResistXiongNuScoreData {
	return dAtA.resistXiongNuScoreData.Map[key]
}

func (dAtA *ConfigDatas) GetResistXiongNuScoreDataArray() []*xiongnu.ResistXiongNuScoreData {
	return dAtA.resistXiongNuScoreData.Array
}

func (dAtA *ConfigDatas) ResistXiongNuScoreData() *ResistXiongNuScoreDataConfig {
	return dAtA.resistXiongNuScoreData
}

type ResistXiongNuScoreDataConfig struct {
	Map   map[uint64]*xiongnu.ResistXiongNuScoreData
	Array []*xiongnu.ResistXiongNuScoreData

	MinKeyData *xiongnu.ResistXiongNuScoreData
	MaxKeyData *xiongnu.ResistXiongNuScoreData

	parserMap map[*xiongnu.ResistXiongNuScoreData]*config.ObjectParser
}

func (d *ResistXiongNuScoreDataConfig) Get(key uint64) *xiongnu.ResistXiongNuScoreData {
	return d.Map[key]
}

func (d *ResistXiongNuScoreDataConfig) Must(key uint64) *xiongnu.ResistXiongNuScoreData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetResistXiongNuWaveData(key uint64) *xiongnu.ResistXiongNuWaveData {
	return dAtA.resistXiongNuWaveData.Map[key]
}

func (dAtA *ConfigDatas) GetResistXiongNuWaveDataArray() []*xiongnu.ResistXiongNuWaveData {
	return dAtA.resistXiongNuWaveData.Array
}

func (dAtA *ConfigDatas) ResistXiongNuWaveData() *ResistXiongNuWaveDataConfig {
	return dAtA.resistXiongNuWaveData
}

type ResistXiongNuWaveDataConfig struct {
	Map   map[uint64]*xiongnu.ResistXiongNuWaveData
	Array []*xiongnu.ResistXiongNuWaveData

	MinKeyData *xiongnu.ResistXiongNuWaveData
	MaxKeyData *xiongnu.ResistXiongNuWaveData

	parserMap map[*xiongnu.ResistXiongNuWaveData]*config.ObjectParser
}

func (d *ResistXiongNuWaveDataConfig) Get(key uint64) *xiongnu.ResistXiongNuWaveData {
	return d.Map[key]
}

func (d *ResistXiongNuWaveDataConfig) Must(key uint64) *xiongnu.ResistXiongNuWaveData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetXuanyuanRangeData(key uint64) *xuanydata.XuanyuanRangeData {
	return dAtA.xuanyuanRangeData.Map[key]
}

func (dAtA *ConfigDatas) GetXuanyuanRangeDataArray() []*xuanydata.XuanyuanRangeData {
	return dAtA.xuanyuanRangeData.Array
}

func (dAtA *ConfigDatas) XuanyuanRangeData() *XuanyuanRangeDataConfig {
	return dAtA.xuanyuanRangeData
}

type XuanyuanRangeDataConfig struct {
	Map   map[uint64]*xuanydata.XuanyuanRangeData
	Array []*xuanydata.XuanyuanRangeData

	MinKeyData *xuanydata.XuanyuanRangeData
	MaxKeyData *xuanydata.XuanyuanRangeData

	parserMap map[*xuanydata.XuanyuanRangeData]*config.ObjectParser
}

func (d *XuanyuanRangeDataConfig) Get(key uint64) *xuanydata.XuanyuanRangeData {
	return d.Map[key]
}

func (d *XuanyuanRangeDataConfig) Must(key uint64) *xuanydata.XuanyuanRangeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetXuanyuanRankPrizeData(key uint64) *xuanydata.XuanyuanRankPrizeData {
	return dAtA.xuanyuanRankPrizeData.Map[key]
}

func (dAtA *ConfigDatas) GetXuanyuanRankPrizeDataArray() []*xuanydata.XuanyuanRankPrizeData {
	return dAtA.xuanyuanRankPrizeData.Array
}

func (dAtA *ConfigDatas) XuanyuanRankPrizeData() *XuanyuanRankPrizeDataConfig {
	return dAtA.xuanyuanRankPrizeData
}

type XuanyuanRankPrizeDataConfig struct {
	Map   map[uint64]*xuanydata.XuanyuanRankPrizeData
	Array []*xuanydata.XuanyuanRankPrizeData

	MinKeyData *xuanydata.XuanyuanRankPrizeData
	MaxKeyData *xuanydata.XuanyuanRankPrizeData

	parserMap map[*xuanydata.XuanyuanRankPrizeData]*config.ObjectParser
}

func (d *XuanyuanRankPrizeDataConfig) Get(key uint64) *xuanydata.XuanyuanRankPrizeData {
	return d.Map[key]
}

func (d *XuanyuanRankPrizeDataConfig) Must(key uint64) *xuanydata.XuanyuanRankPrizeData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetZhanJiangChapterData(key uint64) *zhanjiang.ZhanJiangChapterData {
	return dAtA.zhanJiangChapterData.Map[key]
}

func (dAtA *ConfigDatas) GetZhanJiangChapterDataArray() []*zhanjiang.ZhanJiangChapterData {
	return dAtA.zhanJiangChapterData.Array
}

func (dAtA *ConfigDatas) ZhanJiangChapterData() *ZhanJiangChapterDataConfig {
	return dAtA.zhanJiangChapterData
}

type ZhanJiangChapterDataConfig struct {
	Map   map[uint64]*zhanjiang.ZhanJiangChapterData
	Array []*zhanjiang.ZhanJiangChapterData

	MinKeyData *zhanjiang.ZhanJiangChapterData
	MaxKeyData *zhanjiang.ZhanJiangChapterData

	parserMap map[*zhanjiang.ZhanJiangChapterData]*config.ObjectParser
}

func (d *ZhanJiangChapterDataConfig) Get(key uint64) *zhanjiang.ZhanJiangChapterData {
	return d.Map[key]
}

func (d *ZhanJiangChapterDataConfig) Must(key uint64) *zhanjiang.ZhanJiangChapterData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetZhanJiangData(key uint64) *zhanjiang.ZhanJiangData {
	return dAtA.zhanJiangData.Map[key]
}

func (dAtA *ConfigDatas) GetZhanJiangDataArray() []*zhanjiang.ZhanJiangData {
	return dAtA.zhanJiangData.Array
}

func (dAtA *ConfigDatas) ZhanJiangData() *ZhanJiangDataConfig {
	return dAtA.zhanJiangData
}

type ZhanJiangDataConfig struct {
	Map   map[uint64]*zhanjiang.ZhanJiangData
	Array []*zhanjiang.ZhanJiangData

	MinKeyData *zhanjiang.ZhanJiangData
	MaxKeyData *zhanjiang.ZhanJiangData

	parserMap map[*zhanjiang.ZhanJiangData]*config.ObjectParser
}

func (d *ZhanJiangDataConfig) Get(key uint64) *zhanjiang.ZhanJiangData {
	return d.Map[key]
}

func (d *ZhanJiangDataConfig) Must(key uint64) *zhanjiang.ZhanJiangData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetZhanJiangGuanQiaData(key uint64) *zhanjiang.ZhanJiangGuanQiaData {
	return dAtA.zhanJiangGuanQiaData.Map[key]
}

func (dAtA *ConfigDatas) GetZhanJiangGuanQiaDataArray() []*zhanjiang.ZhanJiangGuanQiaData {
	return dAtA.zhanJiangGuanQiaData.Array
}

func (dAtA *ConfigDatas) ZhanJiangGuanQiaData() *ZhanJiangGuanQiaDataConfig {
	return dAtA.zhanJiangGuanQiaData
}

type ZhanJiangGuanQiaDataConfig struct {
	Map   map[uint64]*zhanjiang.ZhanJiangGuanQiaData
	Array []*zhanjiang.ZhanJiangGuanQiaData

	MinKeyData *zhanjiang.ZhanJiangGuanQiaData
	MaxKeyData *zhanjiang.ZhanJiangGuanQiaData

	parserMap map[*zhanjiang.ZhanJiangGuanQiaData]*config.ObjectParser
}

func (d *ZhanJiangGuanQiaDataConfig) Get(key uint64) *zhanjiang.ZhanJiangGuanQiaData {
	return d.Map[key]
}

func (d *ZhanJiangGuanQiaDataConfig) Must(key uint64) *zhanjiang.ZhanJiangGuanQiaData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetZhengWuCompleteData(key uint64) *zhengwu.ZhengWuCompleteData {
	return dAtA.zhengWuCompleteData.Map[key]
}

func (dAtA *ConfigDatas) GetZhengWuCompleteDataArray() []*zhengwu.ZhengWuCompleteData {
	return dAtA.zhengWuCompleteData.Array
}

func (dAtA *ConfigDatas) ZhengWuCompleteData() *ZhengWuCompleteDataConfig {
	return dAtA.zhengWuCompleteData
}

type ZhengWuCompleteDataConfig struct {
	Map   map[uint64]*zhengwu.ZhengWuCompleteData
	Array []*zhengwu.ZhengWuCompleteData

	MinKeyData *zhengwu.ZhengWuCompleteData
	MaxKeyData *zhengwu.ZhengWuCompleteData

	parserMap map[*zhengwu.ZhengWuCompleteData]*config.ObjectParser
}

func (d *ZhengWuCompleteDataConfig) Get(key uint64) *zhengwu.ZhengWuCompleteData {
	return d.Map[key]
}

func (d *ZhengWuCompleteDataConfig) Must(key uint64) *zhengwu.ZhengWuCompleteData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetZhengWuData(key uint64) *zhengwu.ZhengWuData {
	return dAtA.zhengWuData.Map[key]
}

func (dAtA *ConfigDatas) GetZhengWuDataArray() []*zhengwu.ZhengWuData {
	return dAtA.zhengWuData.Array
}

func (dAtA *ConfigDatas) ZhengWuData() *ZhengWuDataConfig {
	return dAtA.zhengWuData
}

type ZhengWuDataConfig struct {
	Map   map[uint64]*zhengwu.ZhengWuData
	Array []*zhengwu.ZhengWuData

	MinKeyData *zhengwu.ZhengWuData
	MaxKeyData *zhengwu.ZhengWuData

	parserMap map[*zhengwu.ZhengWuData]*config.ObjectParser
}

func (d *ZhengWuDataConfig) Get(key uint64) *zhengwu.ZhengWuData {
	return d.Map[key]
}

func (d *ZhengWuDataConfig) Must(key uint64) *zhengwu.ZhengWuData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) GetZhengWuRefreshData(key uint64) *zhengwu.ZhengWuRefreshData {
	return dAtA.zhengWuRefreshData.Map[key]
}

func (dAtA *ConfigDatas) GetZhengWuRefreshDataArray() []*zhengwu.ZhengWuRefreshData {
	return dAtA.zhengWuRefreshData.Array
}

func (dAtA *ConfigDatas) ZhengWuRefreshData() *ZhengWuRefreshDataConfig {
	return dAtA.zhengWuRefreshData
}

type ZhengWuRefreshDataConfig struct {
	Map   map[uint64]*zhengwu.ZhengWuRefreshData
	Array []*zhengwu.ZhengWuRefreshData

	MinKeyData *zhengwu.ZhengWuRefreshData
	MaxKeyData *zhengwu.ZhengWuRefreshData

	parserMap map[*zhengwu.ZhengWuRefreshData]*config.ObjectParser
}

func (d *ZhengWuRefreshDataConfig) Get(key uint64) *zhengwu.ZhengWuRefreshData {
	return d.Map[key]
}

func (d *ZhengWuRefreshDataConfig) Must(key uint64) *zhengwu.ZhengWuRefreshData {
	v := d.Map[key]
	if v != nil {
		return v
	}

	if key > 0 && int(key) <= len(d.Array) {
		return d.Array[key-1]
	}

	if key <= 0 {
		return d.MinKeyData
	}

	return d.MaxKeyData
}

func (dAtA *ConfigDatas) BaiZhanMiscData() *bai_zhan_data.BaiZhanMiscData {
	return dAtA.baiZhanMiscData
}

func (dAtA *ConfigDatas) CombatConfig() *combatdata.CombatConfig {
	return dAtA.combatConfig
}

func (dAtA *ConfigDatas) CombatMiscConfig() *combatdata.CombatMiscConfig {
	return dAtA.combatMiscConfig
}

func (dAtA *ConfigDatas) EquipCombineDatas() *combine.EquipCombineDatas {
	return dAtA.equipCombineDatas
}

func (dAtA *ConfigDatas) CountryMiscData() *country.CountryMiscData {
	return dAtA.countryMiscData
}

func (dAtA *ConfigDatas) BroadcastHelp() *data.BroadcastHelp {
	return dAtA.broadcastHelp
}

func (dAtA *ConfigDatas) TextHelp() *data.TextHelp {
	return dAtA.textHelp
}

func (dAtA *ConfigDatas) ExchangeMiscData() *dianquan.ExchangeMiscData {
	return dAtA.exchangeMiscData
}

func (dAtA *ConfigDatas) BuildingLayoutMiscData() *domestic_data.BuildingLayoutMiscData {
	return dAtA.buildingLayoutMiscData
}

func (dAtA *ConfigDatas) CityEventMiscData() *domestic_data.CityEventMiscData {
	return dAtA.cityEventMiscData
}

func (dAtA *ConfigDatas) MainCityMiscData() *domestic_data.MainCityMiscData {
	return dAtA.mainCityMiscData
}

func (dAtA *ConfigDatas) DungeonMiscData() *dungeon.DungeonMiscData {
	return dAtA.dungeonMiscData
}

func (dAtA *ConfigDatas) FarmMiscConfig() *farm.FarmMiscConfig {
	return dAtA.farmMiscConfig
}

func (dAtA *ConfigDatas) FishRandomer() *fishing_data.FishRandomer {
	return dAtA.fishRandomer
}

func (dAtA *ConfigDatas) GardenConfig() *gardendata.GardenConfig {
	return dAtA.gardenConfig
}

func (dAtA *ConfigDatas) EquipmentTaozConfig() *goods.EquipmentTaozConfig {
	return dAtA.equipmentTaozConfig
}

func (dAtA *ConfigDatas) GemDatas() *goods.GemDatas {
	return dAtA.gemDatas
}

func (dAtA *ConfigDatas) GoodsCheck() *goods.GoodsCheck {
	return dAtA.goodsCheck
}

func (dAtA *ConfigDatas) GuildLogHelp() *guild_data.GuildLogHelp {
	return dAtA.guildLogHelp
}

func (dAtA *ConfigDatas) NpcGuildSuffixName() *guild_data.NpcGuildSuffixName {
	return dAtA.npcGuildSuffixName
}

func (dAtA *ConfigDatas) HebiMiscData() *hebi.HebiMiscData {
	return dAtA.hebiMiscData
}

func (dAtA *ConfigDatas) HeroCreateData() *heroinit.HeroCreateData {
	return dAtA.heroCreateData
}

func (dAtA *ConfigDatas) HeroInitData() *heroinit.HeroInitData {
	return dAtA.heroInitData
}

func (dAtA *ConfigDatas) MailHelp() *maildata.MailHelp {
	return dAtA.mailHelp
}

func (dAtA *ConfigDatas) JiuGuanMiscData() *military_data.JiuGuanMiscData {
	return dAtA.jiuGuanMiscData
}

func (dAtA *ConfigDatas) JunYingMiscData() *military_data.JunYingMiscData {
	return dAtA.junYingMiscData
}

func (dAtA *ConfigDatas) McBuildMiscData() *mingcdata.McBuildMiscData {
	return dAtA.mcBuildMiscData
}

func (dAtA *ConfigDatas) MingcMiscData() *mingcdata.MingcMiscData {
	return dAtA.mingcMiscData
}

func (dAtA *ConfigDatas) EventLimitGiftConfig() *promdata.EventLimitGiftConfig {
	return dAtA.eventLimitGiftConfig
}

func (dAtA *ConfigDatas) PromotionMiscData() *promdata.PromotionMiscData {
	return dAtA.promotionMiscData
}

func (dAtA *ConfigDatas) QuestionMiscData() *question.QuestionMiscData {
	return dAtA.questionMiscData
}

func (dAtA *ConfigDatas) RaceConfig() *race.RaceConfig {
	return dAtA.raceConfig
}

func (dAtA *ConfigDatas) RandomEventDataDictionary() *random_event.RandomEventDataDictionary {
	return dAtA.randomEventDataDictionary
}

func (dAtA *ConfigDatas) RandomEventPositionDictionary() *random_event.RandomEventPositionDictionary {
	return dAtA.randomEventPositionDictionary
}

func (dAtA *ConfigDatas) RankMiscData() *rank_data.RankMiscData {
	return dAtA.rankMiscData
}

func (dAtA *ConfigDatas) JunTuanNpcPlaceConfig() *regdata.JunTuanNpcPlaceConfig {
	return dAtA.junTuanNpcPlaceConfig
}

func (dAtA *ConfigDatas) SeasonMiscData() *season.SeasonMiscData {
	return dAtA.seasonMiscData
}

func (dAtA *ConfigDatas) SettingMiscData() *settings.SettingMiscData {
	return dAtA.settingMiscData
}

func (dAtA *ConfigDatas) ShopMiscData() *shop.ShopMiscData {
	return dAtA.shopMiscData
}

func (dAtA *ConfigDatas) GoodsConfig() *singleton.GoodsConfig {
	return dAtA.goodsConfig
}

func (dAtA *ConfigDatas) GuildConfig() *singleton.GuildConfig {
	return dAtA.guildConfig
}

func (dAtA *ConfigDatas) GuildGenConfig() *singleton.GuildGenConfig {
	return dAtA.guildGenConfig
}

func (dAtA *ConfigDatas) MilitaryConfig() *singleton.MilitaryConfig {
	return dAtA.militaryConfig
}

func (dAtA *ConfigDatas) MiscConfig() *singleton.MiscConfig {
	return dAtA.miscConfig
}

func (dAtA *ConfigDatas) MiscGenConfig() *singleton.MiscGenConfig {
	return dAtA.miscGenConfig
}

func (dAtA *ConfigDatas) RegionConfig() *singleton.RegionConfig {
	return dAtA.regionConfig
}

func (dAtA *ConfigDatas) RegionGenConfig() *singleton.RegionGenConfig {
	return dAtA.regionGenConfig
}

func (dAtA *ConfigDatas) TagMiscData() *tag.TagMiscData {
	return dAtA.tagMiscData
}

func (dAtA *ConfigDatas) TaskMiscData() *taskdata.TaskMiscData {
	return dAtA.taskMiscData
}

func (dAtA *ConfigDatas) SecretTowerMiscData() *towerdata.SecretTowerMiscData {
	return dAtA.secretTowerMiscData
}

func (dAtA *ConfigDatas) VipMiscData() *vip.VipMiscData {
	return dAtA.vipMiscData
}

func (dAtA *ConfigDatas) ResistXiongNuMisc() *xiongnu.ResistXiongNuMisc {
	return dAtA.resistXiongNuMisc
}

func (dAtA *ConfigDatas) XuanyuanMiscData() *xuanydata.XuanyuanMiscData {
	return dAtA.xuanyuanMiscData
}

func (dAtA *ConfigDatas) ZhanJiangMiscData() *zhanjiang.ZhanJiangMiscData {
	return dAtA.zhanJiangMiscData
}

func (dAtA *ConfigDatas) ZhengWuMiscData() *zhengwu.ZhengWuMiscData {
	return dAtA.zhengWuMiscData
}

func (dAtA *ConfigDatas) ZhengWuRandomData() *zhengwu.ZhengWuRandomData {
	return dAtA.zhengWuRandomData
}

func ParseConfigDatas(gos *config.GameObjects) (*ConfigDatas, error) {

	var err error
	var intKeys []int
	var uint64Keys []uint64
	var stringKeys []string
	if len(intKeys)+len(uint64Keys)+len(stringKeys) > 0 {
	}

	dAtAs := &ConfigDatas{}

	// activitydata.ActivityCollectionData
	dAtAs.activityCollectionData = &ActivityCollectionDataConfig{}
	dAtAs.activityCollectionData.Map, dAtAs.activityCollectionData.parserMap, err = activitydata.LoadActivityCollectionData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.activityCollectionData.Map))
	for k := range dAtAs.activityCollectionData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.activityCollectionData.Array = make([]*activitydata.ActivityCollectionData, 0, len(dAtAs.activityCollectionData.Map))
	for _, k := range uint64Keys {
		dAtAs.activityCollectionData.Array = append(dAtAs.activityCollectionData.Array, dAtAs.activityCollectionData.Map[k])
	}
	dAtAs.activityCollectionData.MinKeyData = dAtAs.activityCollectionData.Array[0]
	dAtAs.activityCollectionData.MaxKeyData = dAtAs.activityCollectionData.Array[len(dAtAs.activityCollectionData.Array)-1]

	// activitydata.ActivityShowData
	dAtAs.activityShowData = &ActivityShowDataConfig{}
	dAtAs.activityShowData.Map, dAtAs.activityShowData.parserMap, err = activitydata.LoadActivityShowData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.activityShowData.Map))
	for k := range dAtAs.activityShowData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.activityShowData.Array = make([]*activitydata.ActivityShowData, 0, len(dAtAs.activityShowData.Map))
	for _, k := range uint64Keys {
		dAtAs.activityShowData.Array = append(dAtAs.activityShowData.Array, dAtAs.activityShowData.Map[k])
	}
	dAtAs.activityShowData.MinKeyData = dAtAs.activityShowData.Array[0]
	dAtAs.activityShowData.MaxKeyData = dAtAs.activityShowData.Array[len(dAtAs.activityShowData.Array)-1]

	// activitydata.ActivityTaskListModeData
	dAtAs.activityTaskListModeData = &ActivityTaskListModeDataConfig{}
	dAtAs.activityTaskListModeData.Map, dAtAs.activityTaskListModeData.parserMap, err = activitydata.LoadActivityTaskListModeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.activityTaskListModeData.Map))
	for k := range dAtAs.activityTaskListModeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.activityTaskListModeData.Array = make([]*activitydata.ActivityTaskListModeData, 0, len(dAtAs.activityTaskListModeData.Map))
	for _, k := range uint64Keys {
		dAtAs.activityTaskListModeData.Array = append(dAtAs.activityTaskListModeData.Array, dAtAs.activityTaskListModeData.Map[k])
	}
	dAtAs.activityTaskListModeData.MinKeyData = dAtAs.activityTaskListModeData.Array[0]
	dAtAs.activityTaskListModeData.MaxKeyData = dAtAs.activityTaskListModeData.Array[len(dAtAs.activityTaskListModeData.Array)-1]

	// activitydata.CollectionExchangeData
	dAtAs.collectionExchangeData = &CollectionExchangeDataConfig{}
	dAtAs.collectionExchangeData.Map, dAtAs.collectionExchangeData.parserMap, err = activitydata.LoadCollectionExchangeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.collectionExchangeData.Map))
	for k := range dAtAs.collectionExchangeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.collectionExchangeData.Array = make([]*activitydata.CollectionExchangeData, 0, len(dAtAs.collectionExchangeData.Map))
	for _, k := range uint64Keys {
		dAtAs.collectionExchangeData.Array = append(dAtAs.collectionExchangeData.Array, dAtAs.collectionExchangeData.Map[k])
	}
	dAtAs.collectionExchangeData.MinKeyData = dAtAs.collectionExchangeData.Array[0]
	dAtAs.collectionExchangeData.MaxKeyData = dAtAs.collectionExchangeData.Array[len(dAtAs.collectionExchangeData.Array)-1]

	// bai_zhan_data.JunXianLevelData
	dAtAs.junXianLevelData = &JunXianLevelDataConfig{}
	dAtAs.junXianLevelData.Map, dAtAs.junXianLevelData.parserMap, err = bai_zhan_data.LoadJunXianLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.junXianLevelData.Map))
	for k := range dAtAs.junXianLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.junXianLevelData.Array = make([]*bai_zhan_data.JunXianLevelData, 0, len(dAtAs.junXianLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.junXianLevelData.Array = append(dAtAs.junXianLevelData.Array, dAtAs.junXianLevelData.Map[k])
	}
	dAtAs.junXianLevelData.MinKeyData = dAtAs.junXianLevelData.Array[0]
	dAtAs.junXianLevelData.MaxKeyData = dAtAs.junXianLevelData.Array[len(dAtAs.junXianLevelData.Array)-1]

	// bai_zhan_data.JunXianPrizeData
	dAtAs.junXianPrizeData = &JunXianPrizeDataConfig{}
	dAtAs.junXianPrizeData.Map, dAtAs.junXianPrizeData.parserMap, err = bai_zhan_data.LoadJunXianPrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.junXianPrizeData.Map))
	for k := range dAtAs.junXianPrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.junXianPrizeData.Array = make([]*bai_zhan_data.JunXianPrizeData, 0, len(dAtAs.junXianPrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.junXianPrizeData.Array = append(dAtAs.junXianPrizeData.Array, dAtAs.junXianPrizeData.Map[k])
	}
	dAtAs.junXianPrizeData.MinKeyData = dAtAs.junXianPrizeData.Array[0]
	dAtAs.junXianPrizeData.MaxKeyData = dAtAs.junXianPrizeData.Array[len(dAtAs.junXianPrizeData.Array)-1]

	// basedata.HomeNpcBaseData
	dAtAs.homeNpcBaseData = &HomeNpcBaseDataConfig{}
	dAtAs.homeNpcBaseData.Map, dAtAs.homeNpcBaseData.parserMap, err = basedata.LoadHomeNpcBaseData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.homeNpcBaseData.Map))
	for k := range dAtAs.homeNpcBaseData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.homeNpcBaseData.Array = make([]*basedata.HomeNpcBaseData, 0, len(dAtAs.homeNpcBaseData.Map))
	for _, k := range uint64Keys {
		dAtAs.homeNpcBaseData.Array = append(dAtAs.homeNpcBaseData.Array, dAtAs.homeNpcBaseData.Map[k])
	}
	dAtAs.homeNpcBaseData.MinKeyData = dAtAs.homeNpcBaseData.Array[0]
	dAtAs.homeNpcBaseData.MaxKeyData = dAtAs.homeNpcBaseData.Array[len(dAtAs.homeNpcBaseData.Array)-1]

	// basedata.NpcBaseData
	dAtAs.npcBaseData = &NpcBaseDataConfig{}
	dAtAs.npcBaseData.Map, dAtAs.npcBaseData.parserMap, err = basedata.LoadNpcBaseData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.npcBaseData.Map))
	for k := range dAtAs.npcBaseData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.npcBaseData.Array = make([]*basedata.NpcBaseData, 0, len(dAtAs.npcBaseData.Map))
	for _, k := range uint64Keys {
		dAtAs.npcBaseData.Array = append(dAtAs.npcBaseData.Array, dAtAs.npcBaseData.Map[k])
	}
	dAtAs.npcBaseData.MinKeyData = dAtAs.npcBaseData.Array[0]
	dAtAs.npcBaseData.MaxKeyData = dAtAs.npcBaseData.Array[len(dAtAs.npcBaseData.Array)-1]

	// blockdata.BlockData
	dAtAs.blockData = &BlockDataConfig{}
	dAtAs.blockData.Map, dAtAs.blockData.parserMap, err = blockdata.LoadBlockData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.blockData.Map))
	for k := range dAtAs.blockData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.blockData.Array = make([]*blockdata.BlockData, 0, len(dAtAs.blockData.Map))
	for _, k := range uint64Keys {
		dAtAs.blockData.Array = append(dAtAs.blockData.Array, dAtAs.blockData.Map[k])
	}
	dAtAs.blockData.MinKeyData = dAtAs.blockData.Array[0]
	dAtAs.blockData.MaxKeyData = dAtAs.blockData.Array[len(dAtAs.blockData.Array)-1]

	// body.BodyData
	dAtAs.bodyData = &BodyDataConfig{}
	dAtAs.bodyData.Map, dAtAs.bodyData.parserMap, err = body.LoadBodyData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.bodyData.Map))
	for k := range dAtAs.bodyData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.bodyData.Array = make([]*body.BodyData, 0, len(dAtAs.bodyData.Map))
	for _, k := range uint64Keys {
		dAtAs.bodyData.Array = append(dAtAs.bodyData.Array, dAtAs.bodyData.Map[k])
	}
	dAtAs.bodyData.MinKeyData = dAtAs.bodyData.Array[0]
	dAtAs.bodyData.MaxKeyData = dAtAs.bodyData.Array[len(dAtAs.bodyData.Array)-1]

	// buffer.BufferData
	dAtAs.bufferData = &BufferDataConfig{}
	dAtAs.bufferData.Map, dAtAs.bufferData.parserMap, err = buffer.LoadBufferData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.bufferData.Map))
	for k := range dAtAs.bufferData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.bufferData.Array = make([]*buffer.BufferData, 0, len(dAtAs.bufferData.Map))
	for _, k := range uint64Keys {
		dAtAs.bufferData.Array = append(dAtAs.bufferData.Array, dAtAs.bufferData.Map[k])
	}
	dAtAs.bufferData.MinKeyData = dAtAs.bufferData.Array[0]
	dAtAs.bufferData.MaxKeyData = dAtAs.bufferData.Array[len(dAtAs.bufferData.Array)-1]

	// buffer.BufferTypeData
	dAtAs.bufferTypeData = &BufferTypeDataConfig{}
	dAtAs.bufferTypeData.Map, dAtAs.bufferTypeData.parserMap, err = buffer.LoadBufferTypeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.bufferTypeData.Map))
	for k := range dAtAs.bufferTypeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.bufferTypeData.Array = make([]*buffer.BufferTypeData, 0, len(dAtAs.bufferTypeData.Map))
	for _, k := range uint64Keys {
		dAtAs.bufferTypeData.Array = append(dAtAs.bufferTypeData.Array, dAtAs.bufferTypeData.Map[k])
	}
	dAtAs.bufferTypeData.MinKeyData = dAtAs.bufferTypeData.Array[0]
	dAtAs.bufferTypeData.MaxKeyData = dAtAs.bufferTypeData.Array[len(dAtAs.bufferTypeData.Array)-1]

	// captain.CaptainAbilityData
	dAtAs.captainAbilityData = &CaptainAbilityDataConfig{}
	dAtAs.captainAbilityData.Map, dAtAs.captainAbilityData.parserMap, err = captain.LoadCaptainAbilityData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.captainAbilityData.Map))
	for k := range dAtAs.captainAbilityData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.captainAbilityData.Array = make([]*captain.CaptainAbilityData, 0, len(dAtAs.captainAbilityData.Map))
	for _, k := range uint64Keys {
		dAtAs.captainAbilityData.Array = append(dAtAs.captainAbilityData.Array, dAtAs.captainAbilityData.Map[k])
	}
	dAtAs.captainAbilityData.MinKeyData = dAtAs.captainAbilityData.Array[0]
	dAtAs.captainAbilityData.MaxKeyData = dAtAs.captainAbilityData.Array[len(dAtAs.captainAbilityData.Array)-1]

	// captain.CaptainData
	dAtAs.captainData = &CaptainDataConfig{}
	dAtAs.captainData.Map, dAtAs.captainData.parserMap, err = captain.LoadCaptainData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.captainData.Map))
	for k := range dAtAs.captainData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.captainData.Array = make([]*captain.CaptainData, 0, len(dAtAs.captainData.Map))
	for _, k := range uint64Keys {
		dAtAs.captainData.Array = append(dAtAs.captainData.Array, dAtAs.captainData.Map[k])
	}
	dAtAs.captainData.MinKeyData = dAtAs.captainData.Array[0]
	dAtAs.captainData.MaxKeyData = dAtAs.captainData.Array[len(dAtAs.captainData.Array)-1]

	// captain.CaptainFriendshipData
	dAtAs.captainFriendshipData = &CaptainFriendshipDataConfig{}
	dAtAs.captainFriendshipData.Map, dAtAs.captainFriendshipData.parserMap, err = captain.LoadCaptainFriendshipData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.captainFriendshipData.Map))
	for k := range dAtAs.captainFriendshipData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.captainFriendshipData.Array = make([]*captain.CaptainFriendshipData, 0, len(dAtAs.captainFriendshipData.Map))
	for _, k := range uint64Keys {
		dAtAs.captainFriendshipData.Array = append(dAtAs.captainFriendshipData.Array, dAtAs.captainFriendshipData.Map[k])
	}
	dAtAs.captainFriendshipData.MinKeyData = dAtAs.captainFriendshipData.Array[0]
	dAtAs.captainFriendshipData.MaxKeyData = dAtAs.captainFriendshipData.Array[len(dAtAs.captainFriendshipData.Array)-1]

	// captain.CaptainLevelData
	dAtAs.captainLevelData = &CaptainLevelDataConfig{}
	dAtAs.captainLevelData.Map, dAtAs.captainLevelData.parserMap, err = captain.LoadCaptainLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.captainLevelData.Map))
	for k := range dAtAs.captainLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.captainLevelData.Array = make([]*captain.CaptainLevelData, 0, len(dAtAs.captainLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.captainLevelData.Array = append(dAtAs.captainLevelData.Array, dAtAs.captainLevelData.Map[k])
	}
	dAtAs.captainLevelData.MinKeyData = dAtAs.captainLevelData.Array[0]
	dAtAs.captainLevelData.MaxKeyData = dAtAs.captainLevelData.Array[len(dAtAs.captainLevelData.Array)-1]

	// captain.CaptainOfficialCountData
	dAtAs.captainOfficialCountData = &CaptainOfficialCountDataConfig{}
	dAtAs.captainOfficialCountData.Map, dAtAs.captainOfficialCountData.parserMap, err = captain.LoadCaptainOfficialCountData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.captainOfficialCountData.Map))
	for k := range dAtAs.captainOfficialCountData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.captainOfficialCountData.Array = make([]*captain.CaptainOfficialCountData, 0, len(dAtAs.captainOfficialCountData.Map))
	for _, k := range uint64Keys {
		dAtAs.captainOfficialCountData.Array = append(dAtAs.captainOfficialCountData.Array, dAtAs.captainOfficialCountData.Map[k])
	}
	dAtAs.captainOfficialCountData.MinKeyData = dAtAs.captainOfficialCountData.Array[0]
	dAtAs.captainOfficialCountData.MaxKeyData = dAtAs.captainOfficialCountData.Array[len(dAtAs.captainOfficialCountData.Array)-1]

	// captain.CaptainOfficialData
	dAtAs.captainOfficialData = &CaptainOfficialDataConfig{}
	dAtAs.captainOfficialData.Map, dAtAs.captainOfficialData.parserMap, err = captain.LoadCaptainOfficialData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.captainOfficialData.Map))
	for k := range dAtAs.captainOfficialData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.captainOfficialData.Array = make([]*captain.CaptainOfficialData, 0, len(dAtAs.captainOfficialData.Map))
	for _, k := range uint64Keys {
		dAtAs.captainOfficialData.Array = append(dAtAs.captainOfficialData.Array, dAtAs.captainOfficialData.Map[k])
	}
	dAtAs.captainOfficialData.MinKeyData = dAtAs.captainOfficialData.Array[0]
	dAtAs.captainOfficialData.MaxKeyData = dAtAs.captainOfficialData.Array[len(dAtAs.captainOfficialData.Array)-1]

	// captain.CaptainRarityData
	dAtAs.captainRarityData = &CaptainRarityDataConfig{}
	dAtAs.captainRarityData.Map, dAtAs.captainRarityData.parserMap, err = captain.LoadCaptainRarityData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.captainRarityData.Map))
	for k := range dAtAs.captainRarityData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.captainRarityData.Array = make([]*captain.CaptainRarityData, 0, len(dAtAs.captainRarityData.Map))
	for _, k := range uint64Keys {
		dAtAs.captainRarityData.Array = append(dAtAs.captainRarityData.Array, dAtAs.captainRarityData.Map[k])
	}
	dAtAs.captainRarityData.MinKeyData = dAtAs.captainRarityData.Array[0]
	dAtAs.captainRarityData.MaxKeyData = dAtAs.captainRarityData.Array[len(dAtAs.captainRarityData.Array)-1]

	// captain.CaptainRebirthLevelData
	dAtAs.captainRebirthLevelData = &CaptainRebirthLevelDataConfig{}
	dAtAs.captainRebirthLevelData.Map, dAtAs.captainRebirthLevelData.parserMap, err = captain.LoadCaptainRebirthLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.captainRebirthLevelData.Map))
	for k := range dAtAs.captainRebirthLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.captainRebirthLevelData.Array = make([]*captain.CaptainRebirthLevelData, 0, len(dAtAs.captainRebirthLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.captainRebirthLevelData.Array = append(dAtAs.captainRebirthLevelData.Array, dAtAs.captainRebirthLevelData.Map[k])
	}
	dAtAs.captainRebirthLevelData.MinKeyData = dAtAs.captainRebirthLevelData.Array[0]
	dAtAs.captainRebirthLevelData.MaxKeyData = dAtAs.captainRebirthLevelData.Array[len(dAtAs.captainRebirthLevelData.Array)-1]

	// captain.CaptainStarData
	dAtAs.captainStarData = &CaptainStarDataConfig{}
	dAtAs.captainStarData.Map, dAtAs.captainStarData.parserMap, err = captain.LoadCaptainStarData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.captainStarData.Map))
	for k := range dAtAs.captainStarData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.captainStarData.Array = make([]*captain.CaptainStarData, 0, len(dAtAs.captainStarData.Map))
	for _, k := range uint64Keys {
		dAtAs.captainStarData.Array = append(dAtAs.captainStarData.Array, dAtAs.captainStarData.Map[k])
	}
	dAtAs.captainStarData.MinKeyData = dAtAs.captainStarData.Array[0]
	dAtAs.captainStarData.MaxKeyData = dAtAs.captainStarData.Array[len(dAtAs.captainStarData.Array)-1]

	// captain.NamelessCaptainData
	dAtAs.namelessCaptainData = &NamelessCaptainDataConfig{}
	dAtAs.namelessCaptainData.Map, dAtAs.namelessCaptainData.parserMap, err = captain.LoadNamelessCaptainData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.namelessCaptainData.Map))
	for k := range dAtAs.namelessCaptainData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.namelessCaptainData.Array = make([]*captain.NamelessCaptainData, 0, len(dAtAs.namelessCaptainData.Map))
	for _, k := range uint64Keys {
		dAtAs.namelessCaptainData.Array = append(dAtAs.namelessCaptainData.Array, dAtAs.namelessCaptainData.Map[k])
	}
	dAtAs.namelessCaptainData.MinKeyData = dAtAs.namelessCaptainData.Array[0]
	dAtAs.namelessCaptainData.MaxKeyData = dAtAs.namelessCaptainData.Array[len(dAtAs.namelessCaptainData.Array)-1]

	// charge.ChargeObjData
	dAtAs.chargeObjData = &ChargeObjDataConfig{}
	dAtAs.chargeObjData.Map, dAtAs.chargeObjData.parserMap, err = charge.LoadChargeObjData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.chargeObjData.Map))
	for k := range dAtAs.chargeObjData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.chargeObjData.Array = make([]*charge.ChargeObjData, 0, len(dAtAs.chargeObjData.Map))
	for _, k := range uint64Keys {
		dAtAs.chargeObjData.Array = append(dAtAs.chargeObjData.Array, dAtAs.chargeObjData.Map[k])
	}
	dAtAs.chargeObjData.MinKeyData = dAtAs.chargeObjData.Array[0]
	dAtAs.chargeObjData.MaxKeyData = dAtAs.chargeObjData.Array[len(dAtAs.chargeObjData.Array)-1]

	// charge.ChargePrizeData
	dAtAs.chargePrizeData = &ChargePrizeDataConfig{}
	dAtAs.chargePrizeData.Map, dAtAs.chargePrizeData.parserMap, err = charge.LoadChargePrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.chargePrizeData.Map))
	for k := range dAtAs.chargePrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.chargePrizeData.Array = make([]*charge.ChargePrizeData, 0, len(dAtAs.chargePrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.chargePrizeData.Array = append(dAtAs.chargePrizeData.Array, dAtAs.chargePrizeData.Map[k])
	}
	dAtAs.chargePrizeData.MinKeyData = dAtAs.chargePrizeData.Array[0]
	dAtAs.chargePrizeData.MaxKeyData = dAtAs.chargePrizeData.Array[len(dAtAs.chargePrizeData.Array)-1]

	// charge.ProductData
	dAtAs.productData = &ProductDataConfig{}
	dAtAs.productData.Map, dAtAs.productData.parserMap, err = charge.LoadProductData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.productData.Map))
	for k := range dAtAs.productData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.productData.Array = make([]*charge.ProductData, 0, len(dAtAs.productData.Map))
	for _, k := range uint64Keys {
		dAtAs.productData.Array = append(dAtAs.productData.Array, dAtAs.productData.Map[k])
	}
	dAtAs.productData.MinKeyData = dAtAs.productData.Array[0]
	dAtAs.productData.MaxKeyData = dAtAs.productData.Array[len(dAtAs.productData.Array)-1]

	// combine.EquipCombineData
	dAtAs.equipCombineData = &EquipCombineDataConfig{}
	dAtAs.equipCombineData.Map, dAtAs.equipCombineData.parserMap, err = combine.LoadEquipCombineData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.equipCombineData.Map))
	for k := range dAtAs.equipCombineData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.equipCombineData.Array = make([]*combine.EquipCombineData, 0, len(dAtAs.equipCombineData.Map))
	for _, k := range uint64Keys {
		dAtAs.equipCombineData.Array = append(dAtAs.equipCombineData.Array, dAtAs.equipCombineData.Map[k])
	}
	dAtAs.equipCombineData.MinKeyData = dAtAs.equipCombineData.Array[0]
	dAtAs.equipCombineData.MaxKeyData = dAtAs.equipCombineData.Array[len(dAtAs.equipCombineData.Array)-1]

	// combine.GoodsCombineData
	dAtAs.goodsCombineData = &GoodsCombineDataConfig{}
	dAtAs.goodsCombineData.Map, dAtAs.goodsCombineData.parserMap, err = combine.LoadGoodsCombineData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.goodsCombineData.Map))
	for k := range dAtAs.goodsCombineData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.goodsCombineData.Array = make([]*combine.GoodsCombineData, 0, len(dAtAs.goodsCombineData.Map))
	for _, k := range uint64Keys {
		dAtAs.goodsCombineData.Array = append(dAtAs.goodsCombineData.Array, dAtAs.goodsCombineData.Map[k])
	}
	dAtAs.goodsCombineData.MinKeyData = dAtAs.goodsCombineData.Array[0]
	dAtAs.goodsCombineData.MaxKeyData = dAtAs.goodsCombineData.Array[len(dAtAs.goodsCombineData.Array)-1]

	// country.CountryData
	dAtAs.countryData = &CountryDataConfig{}
	dAtAs.countryData.Map, dAtAs.countryData.parserMap, err = country.LoadCountryData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.countryData.Map))
	for k := range dAtAs.countryData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.countryData.Array = make([]*country.CountryData, 0, len(dAtAs.countryData.Map))
	for _, k := range uint64Keys {
		dAtAs.countryData.Array = append(dAtAs.countryData.Array, dAtAs.countryData.Map[k])
	}
	dAtAs.countryData.MinKeyData = dAtAs.countryData.Array[0]
	dAtAs.countryData.MaxKeyData = dAtAs.countryData.Array[len(dAtAs.countryData.Array)-1]

	// country.CountryOfficialData
	dAtAs.countryOfficialData = &CountryOfficialDataConfig{}
	dAtAs.countryOfficialData.Map, dAtAs.countryOfficialData.parserMap, err = country.LoadCountryOfficialData(gos)
	if err != nil {
		return nil, err
	}

	intKeys = make([]int, 0, len(dAtAs.countryOfficialData.Map))
	for k := range dAtAs.countryOfficialData.Map {
		intKeys = append(intKeys, k)
	}
	sort.Sort(intSlice(intKeys))
	dAtAs.countryOfficialData.Array = make([]*country.CountryOfficialData, 0, len(dAtAs.countryOfficialData.Map))
	for _, k := range intKeys {
		dAtAs.countryOfficialData.Array = append(dAtAs.countryOfficialData.Array, dAtAs.countryOfficialData.Map[k])
	}
	dAtAs.countryOfficialData.MinKeyData = dAtAs.countryOfficialData.Array[0]
	dAtAs.countryOfficialData.MaxKeyData = dAtAs.countryOfficialData.Array[len(dAtAs.countryOfficialData.Array)-1]

	// country.CountryOfficialNpcData
	dAtAs.countryOfficialNpcData = &CountryOfficialNpcDataConfig{}
	dAtAs.countryOfficialNpcData.Map, dAtAs.countryOfficialNpcData.parserMap, err = country.LoadCountryOfficialNpcData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.countryOfficialNpcData.Map))
	for k := range dAtAs.countryOfficialNpcData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.countryOfficialNpcData.Array = make([]*country.CountryOfficialNpcData, 0, len(dAtAs.countryOfficialNpcData.Map))
	for _, k := range uint64Keys {
		dAtAs.countryOfficialNpcData.Array = append(dAtAs.countryOfficialNpcData.Array, dAtAs.countryOfficialNpcData.Map[k])
	}
	dAtAs.countryOfficialNpcData.MinKeyData = dAtAs.countryOfficialNpcData.Array[0]
	dAtAs.countryOfficialNpcData.MaxKeyData = dAtAs.countryOfficialNpcData.Array[len(dAtAs.countryOfficialNpcData.Array)-1]

	// country.FamilyNameData
	dAtAs.familyNameData = &FamilyNameDataConfig{}
	dAtAs.familyNameData.Map, dAtAs.familyNameData.parserMap, err = country.LoadFamilyNameData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.familyNameData.Map))
	for k := range dAtAs.familyNameData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.familyNameData.Array = make([]*country.FamilyNameData, 0, len(dAtAs.familyNameData.Map))
	for _, k := range uint64Keys {
		dAtAs.familyNameData.Array = append(dAtAs.familyNameData.Array, dAtAs.familyNameData.Map[k])
	}
	dAtAs.familyNameData.MinKeyData = dAtAs.familyNameData.Array[0]
	dAtAs.familyNameData.MaxKeyData = dAtAs.familyNameData.Array[len(dAtAs.familyNameData.Array)-1]

	// data.BroadcastData
	dAtAs.broadcastData = &BroadcastDataConfig{}
	dAtAs.broadcastData.Map, dAtAs.broadcastData.parserMap, err = data.LoadBroadcastData(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.broadcastData.Map))
	for k := range dAtAs.broadcastData.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.broadcastData.Array = make([]*data.BroadcastData, 0, len(dAtAs.broadcastData.Map))
	for _, k := range stringKeys {
		dAtAs.broadcastData.Array = append(dAtAs.broadcastData.Array, dAtAs.broadcastData.Map[k])
	}
	dAtAs.broadcastData.MinKeyData = dAtAs.broadcastData.Array[0]
	dAtAs.broadcastData.MaxKeyData = dAtAs.broadcastData.Array[len(dAtAs.broadcastData.Array)-1]

	// data.BuffEffectData
	dAtAs.buffEffectData = &BuffEffectDataConfig{}
	dAtAs.buffEffectData.Map, dAtAs.buffEffectData.parserMap, err = data.LoadBuffEffectData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.buffEffectData.Map))
	for k := range dAtAs.buffEffectData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.buffEffectData.Array = make([]*data.BuffEffectData, 0, len(dAtAs.buffEffectData.Map))
	for _, k := range uint64Keys {
		dAtAs.buffEffectData.Array = append(dAtAs.buffEffectData.Array, dAtAs.buffEffectData.Map[k])
	}
	dAtAs.buffEffectData.MinKeyData = dAtAs.buffEffectData.Array[0]
	dAtAs.buffEffectData.MaxKeyData = dAtAs.buffEffectData.Array[len(dAtAs.buffEffectData.Array)-1]

	// data.ColorData
	dAtAs.colorData = &ColorDataConfig{}
	dAtAs.colorData.Map, dAtAs.colorData.parserMap, err = data.LoadColorData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.colorData.Map))
	for k := range dAtAs.colorData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.colorData.Array = make([]*data.ColorData, 0, len(dAtAs.colorData.Map))
	for _, k := range uint64Keys {
		dAtAs.colorData.Array = append(dAtAs.colorData.Array, dAtAs.colorData.Map[k])
	}
	dAtAs.colorData.MinKeyData = dAtAs.colorData.Array[0]
	dAtAs.colorData.MaxKeyData = dAtAs.colorData.Array[len(dAtAs.colorData.Array)-1]

	// data.FamilyName
	dAtAs.familyName = &FamilyNameConfig{}
	dAtAs.familyName.Map, dAtAs.familyName.parserMap, err = data.LoadFamilyName(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.familyName.Map))
	for k := range dAtAs.familyName.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.familyName.Array = make([]*data.FamilyName, 0, len(dAtAs.familyName.Map))
	for _, k := range stringKeys {
		dAtAs.familyName.Array = append(dAtAs.familyName.Array, dAtAs.familyName.Map[k])
	}
	dAtAs.familyName.MinKeyData = dAtAs.familyName.Array[0]
	dAtAs.familyName.MaxKeyData = dAtAs.familyName.Array[len(dAtAs.familyName.Array)-1]

	// data.FemaleGivenName
	dAtAs.femaleGivenName = &FemaleGivenNameConfig{}
	dAtAs.femaleGivenName.Map, dAtAs.femaleGivenName.parserMap, err = data.LoadFemaleGivenName(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.femaleGivenName.Map))
	for k := range dAtAs.femaleGivenName.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.femaleGivenName.Array = make([]*data.FemaleGivenName, 0, len(dAtAs.femaleGivenName.Map))
	for _, k := range stringKeys {
		dAtAs.femaleGivenName.Array = append(dAtAs.femaleGivenName.Array, dAtAs.femaleGivenName.Map[k])
	}
	dAtAs.femaleGivenName.MinKeyData = dAtAs.femaleGivenName.Array[0]
	dAtAs.femaleGivenName.MaxKeyData = dAtAs.femaleGivenName.Array[len(dAtAs.femaleGivenName.Array)-1]

	// data.HeroLevelSubData
	dAtAs.heroLevelSubData = &HeroLevelSubDataConfig{}
	dAtAs.heroLevelSubData.Map, dAtAs.heroLevelSubData.parserMap, err = data.LoadHeroLevelSubData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.heroLevelSubData.Map))
	for k := range dAtAs.heroLevelSubData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.heroLevelSubData.Array = make([]*data.HeroLevelSubData, 0, len(dAtAs.heroLevelSubData.Map))
	for _, k := range uint64Keys {
		dAtAs.heroLevelSubData.Array = append(dAtAs.heroLevelSubData.Array, dAtAs.heroLevelSubData.Map[k])
	}
	dAtAs.heroLevelSubData.MinKeyData = dAtAs.heroLevelSubData.Array[0]
	dAtAs.heroLevelSubData.MaxKeyData = dAtAs.heroLevelSubData.Array[len(dAtAs.heroLevelSubData.Array)-1]

	// data.MaleGivenName
	dAtAs.maleGivenName = &MaleGivenNameConfig{}
	dAtAs.maleGivenName.Map, dAtAs.maleGivenName.parserMap, err = data.LoadMaleGivenName(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.maleGivenName.Map))
	for k := range dAtAs.maleGivenName.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.maleGivenName.Array = make([]*data.MaleGivenName, 0, len(dAtAs.maleGivenName.Map))
	for _, k := range stringKeys {
		dAtAs.maleGivenName.Array = append(dAtAs.maleGivenName.Array, dAtAs.maleGivenName.Map[k])
	}
	dAtAs.maleGivenName.MinKeyData = dAtAs.maleGivenName.Array[0]
	dAtAs.maleGivenName.MaxKeyData = dAtAs.maleGivenName.Array[len(dAtAs.maleGivenName.Array)-1]

	// data.SpriteStat
	dAtAs.spriteStat = &SpriteStatConfig{}
	dAtAs.spriteStat.Map, dAtAs.spriteStat.parserMap, err = data.LoadSpriteStat(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.spriteStat.Map))
	for k := range dAtAs.spriteStat.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.spriteStat.Array = make([]*data.SpriteStat, 0, len(dAtAs.spriteStat.Map))
	for _, k := range uint64Keys {
		dAtAs.spriteStat.Array = append(dAtAs.spriteStat.Array, dAtAs.spriteStat.Map[k])
	}
	dAtAs.spriteStat.MinKeyData = dAtAs.spriteStat.Array[0]
	dAtAs.spriteStat.MaxKeyData = dAtAs.spriteStat.Array[len(dAtAs.spriteStat.Array)-1]

	// data.Text
	dAtAs.text = &TextConfig{}
	dAtAs.text.Map, dAtAs.text.parserMap, err = data.LoadText(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.text.Map))
	for k := range dAtAs.text.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.text.Array = make([]*data.Text, 0, len(dAtAs.text.Map))
	for _, k := range stringKeys {
		dAtAs.text.Array = append(dAtAs.text.Array, dAtAs.text.Map[k])
	}
	dAtAs.text.MinKeyData = dAtAs.text.Array[0]
	dAtAs.text.MaxKeyData = dAtAs.text.Array[len(dAtAs.text.Array)-1]

	// data.TimeRuleData
	dAtAs.timeRuleData = &TimeRuleDataConfig{}
	dAtAs.timeRuleData.Map, dAtAs.timeRuleData.parserMap, err = data.LoadTimeRuleData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.timeRuleData.Map))
	for k := range dAtAs.timeRuleData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.timeRuleData.Array = make([]*data.TimeRuleData, 0, len(dAtAs.timeRuleData.Map))
	for _, k := range uint64Keys {
		dAtAs.timeRuleData.Array = append(dAtAs.timeRuleData.Array, dAtAs.timeRuleData.Map[k])
	}
	dAtAs.timeRuleData.MinKeyData = dAtAs.timeRuleData.Array[0]
	dAtAs.timeRuleData.MaxKeyData = dAtAs.timeRuleData.Array[len(dAtAs.timeRuleData.Array)-1]

	// domestic_data.BaseLevelData
	dAtAs.baseLevelData = &BaseLevelDataConfig{}
	dAtAs.baseLevelData.Map, dAtAs.baseLevelData.parserMap, err = domestic_data.LoadBaseLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.baseLevelData.Map))
	for k := range dAtAs.baseLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.baseLevelData.Array = make([]*domestic_data.BaseLevelData, 0, len(dAtAs.baseLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.baseLevelData.Array = append(dAtAs.baseLevelData.Array, dAtAs.baseLevelData.Map[k])
	}
	dAtAs.baseLevelData.MinKeyData = dAtAs.baseLevelData.Array[0]
	dAtAs.baseLevelData.MaxKeyData = dAtAs.baseLevelData.Array[len(dAtAs.baseLevelData.Array)-1]

	// domestic_data.BuildingData
	dAtAs.buildingData = &BuildingDataConfig{}
	dAtAs.buildingData.Map, dAtAs.buildingData.parserMap, err = domestic_data.LoadBuildingData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.buildingData.Map))
	for k := range dAtAs.buildingData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.buildingData.Array = make([]*domestic_data.BuildingData, 0, len(dAtAs.buildingData.Map))
	for _, k := range uint64Keys {
		dAtAs.buildingData.Array = append(dAtAs.buildingData.Array, dAtAs.buildingData.Map[k])
	}
	dAtAs.buildingData.MinKeyData = dAtAs.buildingData.Array[0]
	dAtAs.buildingData.MaxKeyData = dAtAs.buildingData.Array[len(dAtAs.buildingData.Array)-1]

	// domestic_data.BuildingLayoutData
	dAtAs.buildingLayoutData = &BuildingLayoutDataConfig{}
	dAtAs.buildingLayoutData.Map, dAtAs.buildingLayoutData.parserMap, err = domestic_data.LoadBuildingLayoutData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.buildingLayoutData.Map))
	for k := range dAtAs.buildingLayoutData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.buildingLayoutData.Array = make([]*domestic_data.BuildingLayoutData, 0, len(dAtAs.buildingLayoutData.Map))
	for _, k := range uint64Keys {
		dAtAs.buildingLayoutData.Array = append(dAtAs.buildingLayoutData.Array, dAtAs.buildingLayoutData.Map[k])
	}
	dAtAs.buildingLayoutData.MinKeyData = dAtAs.buildingLayoutData.Array[0]
	dAtAs.buildingLayoutData.MaxKeyData = dAtAs.buildingLayoutData.Array[len(dAtAs.buildingLayoutData.Array)-1]

	// domestic_data.BuildingUnlockData
	dAtAs.buildingUnlockData = &BuildingUnlockDataConfig{}
	dAtAs.buildingUnlockData.Map, dAtAs.buildingUnlockData.parserMap, err = domestic_data.LoadBuildingUnlockData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.buildingUnlockData.Map))
	for k := range dAtAs.buildingUnlockData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.buildingUnlockData.Array = make([]*domestic_data.BuildingUnlockData, 0, len(dAtAs.buildingUnlockData.Map))
	for _, k := range uint64Keys {
		dAtAs.buildingUnlockData.Array = append(dAtAs.buildingUnlockData.Array, dAtAs.buildingUnlockData.Map[k])
	}
	dAtAs.buildingUnlockData.MinKeyData = dAtAs.buildingUnlockData.Array[0]
	dAtAs.buildingUnlockData.MaxKeyData = dAtAs.buildingUnlockData.Array[len(dAtAs.buildingUnlockData.Array)-1]

	// domestic_data.CityEventData
	dAtAs.cityEventData = &CityEventDataConfig{}
	dAtAs.cityEventData.Map, dAtAs.cityEventData.parserMap, err = domestic_data.LoadCityEventData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.cityEventData.Map))
	for k := range dAtAs.cityEventData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.cityEventData.Array = make([]*domestic_data.CityEventData, 0, len(dAtAs.cityEventData.Map))
	for _, k := range uint64Keys {
		dAtAs.cityEventData.Array = append(dAtAs.cityEventData.Array, dAtAs.cityEventData.Map[k])
	}
	dAtAs.cityEventData.MinKeyData = dAtAs.cityEventData.Array[0]
	dAtAs.cityEventData.MaxKeyData = dAtAs.cityEventData.Array[len(dAtAs.cityEventData.Array)-1]

	// domestic_data.CityEventLevelData
	dAtAs.cityEventLevelData = &CityEventLevelDataConfig{}
	dAtAs.cityEventLevelData.Map, dAtAs.cityEventLevelData.parserMap, err = domestic_data.LoadCityEventLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.cityEventLevelData.Map))
	for k := range dAtAs.cityEventLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.cityEventLevelData.Array = make([]*domestic_data.CityEventLevelData, 0, len(dAtAs.cityEventLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.cityEventLevelData.Array = append(dAtAs.cityEventLevelData.Array, dAtAs.cityEventLevelData.Map[k])
	}
	dAtAs.cityEventLevelData.MinKeyData = dAtAs.cityEventLevelData.Array[0]
	dAtAs.cityEventLevelData.MaxKeyData = dAtAs.cityEventLevelData.Array[len(dAtAs.cityEventLevelData.Array)-1]

	// domestic_data.CombineCost
	dAtAs.combineCost = &CombineCostConfig{}
	dAtAs.combineCost.Map, dAtAs.combineCost.parserMap, err = domestic_data.LoadCombineCost(gos)
	if err != nil {
		return nil, err
	}

	intKeys = make([]int, 0, len(dAtAs.combineCost.Map))
	for k := range dAtAs.combineCost.Map {
		intKeys = append(intKeys, k)
	}
	sort.Sort(intSlice(intKeys))
	dAtAs.combineCost.Array = make([]*domestic_data.CombineCost, 0, len(dAtAs.combineCost.Map))
	for _, k := range intKeys {
		dAtAs.combineCost.Array = append(dAtAs.combineCost.Array, dAtAs.combineCost.Map[k])
	}
	dAtAs.combineCost.MinKeyData = dAtAs.combineCost.Array[0]
	dAtAs.combineCost.MaxKeyData = dAtAs.combineCost.Array[len(dAtAs.combineCost.Array)-1]

	// domestic_data.CountdownPrizeData
	dAtAs.countdownPrizeData = &CountdownPrizeDataConfig{}
	dAtAs.countdownPrizeData.Map, dAtAs.countdownPrizeData.parserMap, err = domestic_data.LoadCountdownPrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.countdownPrizeData.Map))
	for k := range dAtAs.countdownPrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.countdownPrizeData.Array = make([]*domestic_data.CountdownPrizeData, 0, len(dAtAs.countdownPrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.countdownPrizeData.Array = append(dAtAs.countdownPrizeData.Array, dAtAs.countdownPrizeData.Map[k])
	}
	dAtAs.countdownPrizeData.MinKeyData = dAtAs.countdownPrizeData.Array[0]
	dAtAs.countdownPrizeData.MaxKeyData = dAtAs.countdownPrizeData.Array[len(dAtAs.countdownPrizeData.Array)-1]

	// domestic_data.CountdownPrizeDescData
	dAtAs.countdownPrizeDescData = &CountdownPrizeDescDataConfig{}
	dAtAs.countdownPrizeDescData.Map, dAtAs.countdownPrizeDescData.parserMap, err = domestic_data.LoadCountdownPrizeDescData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.countdownPrizeDescData.Map))
	for k := range dAtAs.countdownPrizeDescData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.countdownPrizeDescData.Array = make([]*domestic_data.CountdownPrizeDescData, 0, len(dAtAs.countdownPrizeDescData.Map))
	for _, k := range uint64Keys {
		dAtAs.countdownPrizeDescData.Array = append(dAtAs.countdownPrizeDescData.Array, dAtAs.countdownPrizeDescData.Map[k])
	}
	dAtAs.countdownPrizeDescData.MinKeyData = dAtAs.countdownPrizeDescData.Array[0]
	dAtAs.countdownPrizeDescData.MaxKeyData = dAtAs.countdownPrizeDescData.Array[len(dAtAs.countdownPrizeDescData.Array)-1]

	// domestic_data.GuanFuLevelData
	dAtAs.guanFuLevelData = &GuanFuLevelDataConfig{}
	dAtAs.guanFuLevelData.Map, dAtAs.guanFuLevelData.parserMap, err = domestic_data.LoadGuanFuLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guanFuLevelData.Map))
	for k := range dAtAs.guanFuLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guanFuLevelData.Array = make([]*domestic_data.GuanFuLevelData, 0, len(dAtAs.guanFuLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.guanFuLevelData.Array = append(dAtAs.guanFuLevelData.Array, dAtAs.guanFuLevelData.Map[k])
	}
	dAtAs.guanFuLevelData.MinKeyData = dAtAs.guanFuLevelData.Array[0]
	dAtAs.guanFuLevelData.MaxKeyData = dAtAs.guanFuLevelData.Array[len(dAtAs.guanFuLevelData.Array)-1]

	// domestic_data.OuterCityBuildingData
	dAtAs.outerCityBuildingData = &OuterCityBuildingDataConfig{}
	dAtAs.outerCityBuildingData.Map, dAtAs.outerCityBuildingData.parserMap, err = domestic_data.LoadOuterCityBuildingData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.outerCityBuildingData.Map))
	for k := range dAtAs.outerCityBuildingData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.outerCityBuildingData.Array = make([]*domestic_data.OuterCityBuildingData, 0, len(dAtAs.outerCityBuildingData.Map))
	for _, k := range uint64Keys {
		dAtAs.outerCityBuildingData.Array = append(dAtAs.outerCityBuildingData.Array, dAtAs.outerCityBuildingData.Map[k])
	}
	dAtAs.outerCityBuildingData.MinKeyData = dAtAs.outerCityBuildingData.Array[0]
	dAtAs.outerCityBuildingData.MaxKeyData = dAtAs.outerCityBuildingData.Array[len(dAtAs.outerCityBuildingData.Array)-1]

	// domestic_data.OuterCityBuildingDescData
	dAtAs.outerCityBuildingDescData = &OuterCityBuildingDescDataConfig{}
	dAtAs.outerCityBuildingDescData.Map, dAtAs.outerCityBuildingDescData.parserMap, err = domestic_data.LoadOuterCityBuildingDescData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.outerCityBuildingDescData.Map))
	for k := range dAtAs.outerCityBuildingDescData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.outerCityBuildingDescData.Array = make([]*domestic_data.OuterCityBuildingDescData, 0, len(dAtAs.outerCityBuildingDescData.Map))
	for _, k := range uint64Keys {
		dAtAs.outerCityBuildingDescData.Array = append(dAtAs.outerCityBuildingDescData.Array, dAtAs.outerCityBuildingDescData.Map[k])
	}
	dAtAs.outerCityBuildingDescData.MinKeyData = dAtAs.outerCityBuildingDescData.Array[0]
	dAtAs.outerCityBuildingDescData.MaxKeyData = dAtAs.outerCityBuildingDescData.Array[len(dAtAs.outerCityBuildingDescData.Array)-1]

	// domestic_data.OuterCityData
	dAtAs.outerCityData = &OuterCityDataConfig{}
	dAtAs.outerCityData.Map, dAtAs.outerCityData.parserMap, err = domestic_data.LoadOuterCityData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.outerCityData.Map))
	for k := range dAtAs.outerCityData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.outerCityData.Array = make([]*domestic_data.OuterCityData, 0, len(dAtAs.outerCityData.Map))
	for _, k := range uint64Keys {
		dAtAs.outerCityData.Array = append(dAtAs.outerCityData.Array, dAtAs.outerCityData.Map[k])
	}
	dAtAs.outerCityData.MinKeyData = dAtAs.outerCityData.Array[0]
	dAtAs.outerCityData.MaxKeyData = dAtAs.outerCityData.Array[len(dAtAs.outerCityData.Array)-1]

	// domestic_data.OuterCityLayoutData
	dAtAs.outerCityLayoutData = &OuterCityLayoutDataConfig{}
	dAtAs.outerCityLayoutData.Map, dAtAs.outerCityLayoutData.parserMap, err = domestic_data.LoadOuterCityLayoutData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.outerCityLayoutData.Map))
	for k := range dAtAs.outerCityLayoutData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.outerCityLayoutData.Array = make([]*domestic_data.OuterCityLayoutData, 0, len(dAtAs.outerCityLayoutData.Map))
	for _, k := range uint64Keys {
		dAtAs.outerCityLayoutData.Array = append(dAtAs.outerCityLayoutData.Array, dAtAs.outerCityLayoutData.Map[k])
	}
	dAtAs.outerCityLayoutData.MinKeyData = dAtAs.outerCityLayoutData.Array[0]
	dAtAs.outerCityLayoutData.MaxKeyData = dAtAs.outerCityLayoutData.Array[len(dAtAs.outerCityLayoutData.Array)-1]

	// domestic_data.ProsperityDamageBuffData
	dAtAs.prosperityDamageBuffData = &ProsperityDamageBuffDataConfig{}
	dAtAs.prosperityDamageBuffData.Map, dAtAs.prosperityDamageBuffData.parserMap, err = domestic_data.LoadProsperityDamageBuffData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.prosperityDamageBuffData.Map))
	for k := range dAtAs.prosperityDamageBuffData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.prosperityDamageBuffData.Array = make([]*domestic_data.ProsperityDamageBuffData, 0, len(dAtAs.prosperityDamageBuffData.Map))
	for _, k := range uint64Keys {
		dAtAs.prosperityDamageBuffData.Array = append(dAtAs.prosperityDamageBuffData.Array, dAtAs.prosperityDamageBuffData.Map[k])
	}
	dAtAs.prosperityDamageBuffData.MinKeyData = dAtAs.prosperityDamageBuffData.Array[0]
	dAtAs.prosperityDamageBuffData.MaxKeyData = dAtAs.prosperityDamageBuffData.Array[len(dAtAs.prosperityDamageBuffData.Array)-1]

	// domestic_data.SoldierLevelData
	dAtAs.soldierLevelData = &SoldierLevelDataConfig{}
	dAtAs.soldierLevelData.Map, dAtAs.soldierLevelData.parserMap, err = domestic_data.LoadSoldierLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.soldierLevelData.Map))
	for k := range dAtAs.soldierLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.soldierLevelData.Array = make([]*domestic_data.SoldierLevelData, 0, len(dAtAs.soldierLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.soldierLevelData.Array = append(dAtAs.soldierLevelData.Array, dAtAs.soldierLevelData.Map[k])
	}
	dAtAs.soldierLevelData.MinKeyData = dAtAs.soldierLevelData.Array[0]
	dAtAs.soldierLevelData.MaxKeyData = dAtAs.soldierLevelData.Array[len(dAtAs.soldierLevelData.Array)-1]

	// domestic_data.TechnologyData
	dAtAs.technologyData = &TechnologyDataConfig{}
	dAtAs.technologyData.Map, dAtAs.technologyData.parserMap, err = domestic_data.LoadTechnologyData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.technologyData.Map))
	for k := range dAtAs.technologyData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.technologyData.Array = make([]*domestic_data.TechnologyData, 0, len(dAtAs.technologyData.Map))
	for _, k := range uint64Keys {
		dAtAs.technologyData.Array = append(dAtAs.technologyData.Array, dAtAs.technologyData.Map[k])
	}
	dAtAs.technologyData.MinKeyData = dAtAs.technologyData.Array[0]
	dAtAs.technologyData.MaxKeyData = dAtAs.technologyData.Array[len(dAtAs.technologyData.Array)-1]

	// domestic_data.TieJiangPuLevelData
	dAtAs.tieJiangPuLevelData = &TieJiangPuLevelDataConfig{}
	dAtAs.tieJiangPuLevelData.Map, dAtAs.tieJiangPuLevelData.parserMap, err = domestic_data.LoadTieJiangPuLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.tieJiangPuLevelData.Map))
	for k := range dAtAs.tieJiangPuLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.tieJiangPuLevelData.Array = make([]*domestic_data.TieJiangPuLevelData, 0, len(dAtAs.tieJiangPuLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.tieJiangPuLevelData.Array = append(dAtAs.tieJiangPuLevelData.Array, dAtAs.tieJiangPuLevelData.Map[k])
	}
	dAtAs.tieJiangPuLevelData.MinKeyData = dAtAs.tieJiangPuLevelData.Array[0]
	dAtAs.tieJiangPuLevelData.MaxKeyData = dAtAs.tieJiangPuLevelData.Array[len(dAtAs.tieJiangPuLevelData.Array)-1]

	// domestic_data.WorkshopDuration
	dAtAs.workshopDuration = &WorkshopDurationConfig{}
	dAtAs.workshopDuration.Map, dAtAs.workshopDuration.parserMap, err = domestic_data.LoadWorkshopDuration(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.workshopDuration.Map))
	for k := range dAtAs.workshopDuration.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.workshopDuration.Array = make([]*domestic_data.WorkshopDuration, 0, len(dAtAs.workshopDuration.Map))
	for _, k := range uint64Keys {
		dAtAs.workshopDuration.Array = append(dAtAs.workshopDuration.Array, dAtAs.workshopDuration.Map[k])
	}
	dAtAs.workshopDuration.MinKeyData = dAtAs.workshopDuration.Array[0]
	dAtAs.workshopDuration.MaxKeyData = dAtAs.workshopDuration.Array[len(dAtAs.workshopDuration.Array)-1]

	// domestic_data.WorkshopLevelData
	dAtAs.workshopLevelData = &WorkshopLevelDataConfig{}
	dAtAs.workshopLevelData.Map, dAtAs.workshopLevelData.parserMap, err = domestic_data.LoadWorkshopLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.workshopLevelData.Map))
	for k := range dAtAs.workshopLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.workshopLevelData.Array = make([]*domestic_data.WorkshopLevelData, 0, len(dAtAs.workshopLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.workshopLevelData.Array = append(dAtAs.workshopLevelData.Array, dAtAs.workshopLevelData.Map[k])
	}
	dAtAs.workshopLevelData.MinKeyData = dAtAs.workshopLevelData.Array[0]
	dAtAs.workshopLevelData.MaxKeyData = dAtAs.workshopLevelData.Array[len(dAtAs.workshopLevelData.Array)-1]

	// domestic_data.WorkshopRefreshCost
	dAtAs.workshopRefreshCost = &WorkshopRefreshCostConfig{}
	dAtAs.workshopRefreshCost.Map, dAtAs.workshopRefreshCost.parserMap, err = domestic_data.LoadWorkshopRefreshCost(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.workshopRefreshCost.Map))
	for k := range dAtAs.workshopRefreshCost.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.workshopRefreshCost.Array = make([]*domestic_data.WorkshopRefreshCost, 0, len(dAtAs.workshopRefreshCost.Map))
	for _, k := range uint64Keys {
		dAtAs.workshopRefreshCost.Array = append(dAtAs.workshopRefreshCost.Array, dAtAs.workshopRefreshCost.Map[k])
	}
	dAtAs.workshopRefreshCost.MinKeyData = dAtAs.workshopRefreshCost.Array[0]
	dAtAs.workshopRefreshCost.MaxKeyData = dAtAs.workshopRefreshCost.Array[len(dAtAs.workshopRefreshCost.Array)-1]

	// dungeon.DungeonChapterData
	dAtAs.dungeonChapterData = &DungeonChapterDataConfig{}
	dAtAs.dungeonChapterData.Map, dAtAs.dungeonChapterData.parserMap, err = dungeon.LoadDungeonChapterData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.dungeonChapterData.Map))
	for k := range dAtAs.dungeonChapterData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.dungeonChapterData.Array = make([]*dungeon.DungeonChapterData, 0, len(dAtAs.dungeonChapterData.Map))
	for _, k := range uint64Keys {
		dAtAs.dungeonChapterData.Array = append(dAtAs.dungeonChapterData.Array, dAtAs.dungeonChapterData.Map[k])
	}
	dAtAs.dungeonChapterData.MinKeyData = dAtAs.dungeonChapterData.Array[0]
	dAtAs.dungeonChapterData.MaxKeyData = dAtAs.dungeonChapterData.Array[len(dAtAs.dungeonChapterData.Array)-1]

	// dungeon.DungeonData
	dAtAs.dungeonData = &DungeonDataConfig{}
	dAtAs.dungeonData.Map, dAtAs.dungeonData.parserMap, err = dungeon.LoadDungeonData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.dungeonData.Map))
	for k := range dAtAs.dungeonData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.dungeonData.Array = make([]*dungeon.DungeonData, 0, len(dAtAs.dungeonData.Map))
	for _, k := range uint64Keys {
		dAtAs.dungeonData.Array = append(dAtAs.dungeonData.Array, dAtAs.dungeonData.Map[k])
	}
	dAtAs.dungeonData.MinKeyData = dAtAs.dungeonData.Array[0]
	dAtAs.dungeonData.MaxKeyData = dAtAs.dungeonData.Array[len(dAtAs.dungeonData.Array)-1]

	// dungeon.DungeonGuideTroopData
	dAtAs.dungeonGuideTroopData = &DungeonGuideTroopDataConfig{}
	dAtAs.dungeonGuideTroopData.Map, dAtAs.dungeonGuideTroopData.parserMap, err = dungeon.LoadDungeonGuideTroopData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.dungeonGuideTroopData.Map))
	for k := range dAtAs.dungeonGuideTroopData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.dungeonGuideTroopData.Array = make([]*dungeon.DungeonGuideTroopData, 0, len(dAtAs.dungeonGuideTroopData.Map))
	for _, k := range uint64Keys {
		dAtAs.dungeonGuideTroopData.Array = append(dAtAs.dungeonGuideTroopData.Array, dAtAs.dungeonGuideTroopData.Map[k])
	}
	dAtAs.dungeonGuideTroopData.MinKeyData = dAtAs.dungeonGuideTroopData.Array[0]
	dAtAs.dungeonGuideTroopData.MaxKeyData = dAtAs.dungeonGuideTroopData.Array[len(dAtAs.dungeonGuideTroopData.Array)-1]

	// farm.FarmMaxStealConfig
	dAtAs.farmMaxStealConfig = &FarmMaxStealConfigConfig{}
	dAtAs.farmMaxStealConfig.Map, dAtAs.farmMaxStealConfig.parserMap, err = farm.LoadFarmMaxStealConfig(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.farmMaxStealConfig.Map))
	for k := range dAtAs.farmMaxStealConfig.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.farmMaxStealConfig.Array = make([]*farm.FarmMaxStealConfig, 0, len(dAtAs.farmMaxStealConfig.Map))
	for _, k := range uint64Keys {
		dAtAs.farmMaxStealConfig.Array = append(dAtAs.farmMaxStealConfig.Array, dAtAs.farmMaxStealConfig.Map[k])
	}
	dAtAs.farmMaxStealConfig.MinKeyData = dAtAs.farmMaxStealConfig.Array[0]
	dAtAs.farmMaxStealConfig.MaxKeyData = dAtAs.farmMaxStealConfig.Array[len(dAtAs.farmMaxStealConfig.Array)-1]

	// farm.FarmOneKeyConfig
	dAtAs.farmOneKeyConfig = &FarmOneKeyConfigConfig{}
	dAtAs.farmOneKeyConfig.Map, dAtAs.farmOneKeyConfig.parserMap, err = farm.LoadFarmOneKeyConfig(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.farmOneKeyConfig.Map))
	for k := range dAtAs.farmOneKeyConfig.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.farmOneKeyConfig.Array = make([]*farm.FarmOneKeyConfig, 0, len(dAtAs.farmOneKeyConfig.Map))
	for _, k := range uint64Keys {
		dAtAs.farmOneKeyConfig.Array = append(dAtAs.farmOneKeyConfig.Array, dAtAs.farmOneKeyConfig.Map[k])
	}
	dAtAs.farmOneKeyConfig.MinKeyData = dAtAs.farmOneKeyConfig.Array[0]
	dAtAs.farmOneKeyConfig.MaxKeyData = dAtAs.farmOneKeyConfig.Array[len(dAtAs.farmOneKeyConfig.Array)-1]

	// farm.FarmResConfig
	dAtAs.farmResConfig = &FarmResConfigConfig{}
	dAtAs.farmResConfig.Map, dAtAs.farmResConfig.parserMap, err = farm.LoadFarmResConfig(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.farmResConfig.Map))
	for k := range dAtAs.farmResConfig.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.farmResConfig.Array = make([]*farm.FarmResConfig, 0, len(dAtAs.farmResConfig.Map))
	for _, k := range uint64Keys {
		dAtAs.farmResConfig.Array = append(dAtAs.farmResConfig.Array, dAtAs.farmResConfig.Map[k])
	}
	dAtAs.farmResConfig.MinKeyData = dAtAs.farmResConfig.Array[0]
	dAtAs.farmResConfig.MaxKeyData = dAtAs.farmResConfig.Array[len(dAtAs.farmResConfig.Array)-1]

	// fishing_data.FishData
	dAtAs.fishData = &FishDataConfig{}
	dAtAs.fishData.Map, dAtAs.fishData.parserMap, err = fishing_data.LoadFishData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.fishData.Map))
	for k := range dAtAs.fishData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.fishData.Array = make([]*fishing_data.FishData, 0, len(dAtAs.fishData.Map))
	for _, k := range uint64Keys {
		dAtAs.fishData.Array = append(dAtAs.fishData.Array, dAtAs.fishData.Map[k])
	}
	dAtAs.fishData.MinKeyData = dAtAs.fishData.Array[0]
	dAtAs.fishData.MaxKeyData = dAtAs.fishData.Array[len(dAtAs.fishData.Array)-1]

	// fishing_data.FishingCaptainProbabilityData
	dAtAs.fishingCaptainProbabilityData = &FishingCaptainProbabilityDataConfig{}
	dAtAs.fishingCaptainProbabilityData.Map, dAtAs.fishingCaptainProbabilityData.parserMap, err = fishing_data.LoadFishingCaptainProbabilityData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.fishingCaptainProbabilityData.Map))
	for k := range dAtAs.fishingCaptainProbabilityData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.fishingCaptainProbabilityData.Array = make([]*fishing_data.FishingCaptainProbabilityData, 0, len(dAtAs.fishingCaptainProbabilityData.Map))
	for _, k := range uint64Keys {
		dAtAs.fishingCaptainProbabilityData.Array = append(dAtAs.fishingCaptainProbabilityData.Array, dAtAs.fishingCaptainProbabilityData.Map[k])
	}
	dAtAs.fishingCaptainProbabilityData.MinKeyData = dAtAs.fishingCaptainProbabilityData.Array[0]
	dAtAs.fishingCaptainProbabilityData.MaxKeyData = dAtAs.fishingCaptainProbabilityData.Array[len(dAtAs.fishingCaptainProbabilityData.Array)-1]

	// fishing_data.FishingCostData
	dAtAs.fishingCostData = &FishingCostDataConfig{}
	dAtAs.fishingCostData.Map, dAtAs.fishingCostData.parserMap, err = fishing_data.LoadFishingCostData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.fishingCostData.Map))
	for k := range dAtAs.fishingCostData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.fishingCostData.Array = make([]*fishing_data.FishingCostData, 0, len(dAtAs.fishingCostData.Map))
	for _, k := range uint64Keys {
		dAtAs.fishingCostData.Array = append(dAtAs.fishingCostData.Array, dAtAs.fishingCostData.Map[k])
	}
	dAtAs.fishingCostData.MinKeyData = dAtAs.fishingCostData.Array[0]
	dAtAs.fishingCostData.MaxKeyData = dAtAs.fishingCostData.Array[len(dAtAs.fishingCostData.Array)-1]

	// fishing_data.FishingShowData
	dAtAs.fishingShowData = &FishingShowDataConfig{}
	dAtAs.fishingShowData.Map, dAtAs.fishingShowData.parserMap, err = fishing_data.LoadFishingShowData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.fishingShowData.Map))
	for k := range dAtAs.fishingShowData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.fishingShowData.Array = make([]*fishing_data.FishingShowData, 0, len(dAtAs.fishingShowData.Map))
	for _, k := range uint64Keys {
		dAtAs.fishingShowData.Array = append(dAtAs.fishingShowData.Array, dAtAs.fishingShowData.Map[k])
	}
	dAtAs.fishingShowData.MinKeyData = dAtAs.fishingShowData.Array[0]
	dAtAs.fishingShowData.MaxKeyData = dAtAs.fishingShowData.Array[len(dAtAs.fishingShowData.Array)-1]

	// function.FunctionOpenData
	dAtAs.functionOpenData = &FunctionOpenDataConfig{}
	dAtAs.functionOpenData.Map, dAtAs.functionOpenData.parserMap, err = function.LoadFunctionOpenData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.functionOpenData.Map))
	for k := range dAtAs.functionOpenData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.functionOpenData.Array = make([]*function.FunctionOpenData, 0, len(dAtAs.functionOpenData.Map))
	for _, k := range uint64Keys {
		dAtAs.functionOpenData.Array = append(dAtAs.functionOpenData.Array, dAtAs.functionOpenData.Map[k])
	}
	dAtAs.functionOpenData.MinKeyData = dAtAs.functionOpenData.Array[0]
	dAtAs.functionOpenData.MaxKeyData = dAtAs.functionOpenData.Array[len(dAtAs.functionOpenData.Array)-1]

	// gardendata.TreasuryTreeData
	dAtAs.treasuryTreeData = &TreasuryTreeDataConfig{}
	dAtAs.treasuryTreeData.Map, dAtAs.treasuryTreeData.parserMap, err = gardendata.LoadTreasuryTreeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.treasuryTreeData.Map))
	for k := range dAtAs.treasuryTreeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.treasuryTreeData.Array = make([]*gardendata.TreasuryTreeData, 0, len(dAtAs.treasuryTreeData.Map))
	for _, k := range uint64Keys {
		dAtAs.treasuryTreeData.Array = append(dAtAs.treasuryTreeData.Array, dAtAs.treasuryTreeData.Map[k])
	}
	dAtAs.treasuryTreeData.MinKeyData = dAtAs.treasuryTreeData.Array[0]
	dAtAs.treasuryTreeData.MaxKeyData = dAtAs.treasuryTreeData.Array[len(dAtAs.treasuryTreeData.Array)-1]

	// goods.EquipmentData
	dAtAs.equipmentData = &EquipmentDataConfig{}
	dAtAs.equipmentData.Map, dAtAs.equipmentData.parserMap, err = goods.LoadEquipmentData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.equipmentData.Map))
	for k := range dAtAs.equipmentData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.equipmentData.Array = make([]*goods.EquipmentData, 0, len(dAtAs.equipmentData.Map))
	for _, k := range uint64Keys {
		dAtAs.equipmentData.Array = append(dAtAs.equipmentData.Array, dAtAs.equipmentData.Map[k])
	}
	dAtAs.equipmentData.MinKeyData = dAtAs.equipmentData.Array[0]
	dAtAs.equipmentData.MaxKeyData = dAtAs.equipmentData.Array[len(dAtAs.equipmentData.Array)-1]

	// goods.EquipmentLevelData
	dAtAs.equipmentLevelData = &EquipmentLevelDataConfig{}
	dAtAs.equipmentLevelData.Map, dAtAs.equipmentLevelData.parserMap, err = goods.LoadEquipmentLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.equipmentLevelData.Map))
	for k := range dAtAs.equipmentLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.equipmentLevelData.Array = make([]*goods.EquipmentLevelData, 0, len(dAtAs.equipmentLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.equipmentLevelData.Array = append(dAtAs.equipmentLevelData.Array, dAtAs.equipmentLevelData.Map[k])
	}
	dAtAs.equipmentLevelData.MinKeyData = dAtAs.equipmentLevelData.Array[0]
	dAtAs.equipmentLevelData.MaxKeyData = dAtAs.equipmentLevelData.Array[len(dAtAs.equipmentLevelData.Array)-1]

	// goods.EquipmentQualityData
	dAtAs.equipmentQualityData = &EquipmentQualityDataConfig{}
	dAtAs.equipmentQualityData.Map, dAtAs.equipmentQualityData.parserMap, err = goods.LoadEquipmentQualityData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.equipmentQualityData.Map))
	for k := range dAtAs.equipmentQualityData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.equipmentQualityData.Array = make([]*goods.EquipmentQualityData, 0, len(dAtAs.equipmentQualityData.Map))
	for _, k := range uint64Keys {
		dAtAs.equipmentQualityData.Array = append(dAtAs.equipmentQualityData.Array, dAtAs.equipmentQualityData.Map[k])
	}
	dAtAs.equipmentQualityData.MinKeyData = dAtAs.equipmentQualityData.Array[0]
	dAtAs.equipmentQualityData.MaxKeyData = dAtAs.equipmentQualityData.Array[len(dAtAs.equipmentQualityData.Array)-1]

	// goods.EquipmentRefinedData
	dAtAs.equipmentRefinedData = &EquipmentRefinedDataConfig{}
	dAtAs.equipmentRefinedData.Map, dAtAs.equipmentRefinedData.parserMap, err = goods.LoadEquipmentRefinedData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.equipmentRefinedData.Map))
	for k := range dAtAs.equipmentRefinedData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.equipmentRefinedData.Array = make([]*goods.EquipmentRefinedData, 0, len(dAtAs.equipmentRefinedData.Map))
	for _, k := range uint64Keys {
		dAtAs.equipmentRefinedData.Array = append(dAtAs.equipmentRefinedData.Array, dAtAs.equipmentRefinedData.Map[k])
	}
	dAtAs.equipmentRefinedData.MinKeyData = dAtAs.equipmentRefinedData.Array[0]
	dAtAs.equipmentRefinedData.MaxKeyData = dAtAs.equipmentRefinedData.Array[len(dAtAs.equipmentRefinedData.Array)-1]

	// goods.EquipmentTaozData
	dAtAs.equipmentTaozData = &EquipmentTaozDataConfig{}
	dAtAs.equipmentTaozData.Map, dAtAs.equipmentTaozData.parserMap, err = goods.LoadEquipmentTaozData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.equipmentTaozData.Map))
	for k := range dAtAs.equipmentTaozData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.equipmentTaozData.Array = make([]*goods.EquipmentTaozData, 0, len(dAtAs.equipmentTaozData.Map))
	for _, k := range uint64Keys {
		dAtAs.equipmentTaozData.Array = append(dAtAs.equipmentTaozData.Array, dAtAs.equipmentTaozData.Map[k])
	}
	dAtAs.equipmentTaozData.MinKeyData = dAtAs.equipmentTaozData.Array[0]
	dAtAs.equipmentTaozData.MaxKeyData = dAtAs.equipmentTaozData.Array[len(dAtAs.equipmentTaozData.Array)-1]

	// goods.GemData
	dAtAs.gemData = &GemDataConfig{}
	dAtAs.gemData.Map, dAtAs.gemData.parserMap, err = goods.LoadGemData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.gemData.Map))
	for k := range dAtAs.gemData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.gemData.Array = make([]*goods.GemData, 0, len(dAtAs.gemData.Map))
	for _, k := range uint64Keys {
		dAtAs.gemData.Array = append(dAtAs.gemData.Array, dAtAs.gemData.Map[k])
	}
	dAtAs.gemData.MinKeyData = dAtAs.gemData.Array[0]
	dAtAs.gemData.MaxKeyData = dAtAs.gemData.Array[len(dAtAs.gemData.Array)-1]

	// goods.GoodsData
	dAtAs.goodsData = &GoodsDataConfig{}
	dAtAs.goodsData.Map, dAtAs.goodsData.parserMap, err = goods.LoadGoodsData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.goodsData.Map))
	for k := range dAtAs.goodsData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.goodsData.Array = make([]*goods.GoodsData, 0, len(dAtAs.goodsData.Map))
	for _, k := range uint64Keys {
		dAtAs.goodsData.Array = append(dAtAs.goodsData.Array, dAtAs.goodsData.Map[k])
	}
	dAtAs.goodsData.MinKeyData = dAtAs.goodsData.Array[0]
	dAtAs.goodsData.MaxKeyData = dAtAs.goodsData.Array[len(dAtAs.goodsData.Array)-1]

	// goods.GoodsQuality
	dAtAs.goodsQuality = &GoodsQualityConfig{}
	dAtAs.goodsQuality.Map, dAtAs.goodsQuality.parserMap, err = goods.LoadGoodsQuality(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.goodsQuality.Map))
	for k := range dAtAs.goodsQuality.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.goodsQuality.Array = make([]*goods.GoodsQuality, 0, len(dAtAs.goodsQuality.Map))
	for _, k := range uint64Keys {
		dAtAs.goodsQuality.Array = append(dAtAs.goodsQuality.Array, dAtAs.goodsQuality.Map[k])
	}
	dAtAs.goodsQuality.MinKeyData = dAtAs.goodsQuality.Array[0]
	dAtAs.goodsQuality.MaxKeyData = dAtAs.goodsQuality.Array[len(dAtAs.goodsQuality.Array)-1]

	// guild_data.GuildBigBoxData
	dAtAs.guildBigBoxData = &GuildBigBoxDataConfig{}
	dAtAs.guildBigBoxData.Map, dAtAs.guildBigBoxData.parserMap, err = guild_data.LoadGuildBigBoxData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildBigBoxData.Map))
	for k := range dAtAs.guildBigBoxData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildBigBoxData.Array = make([]*guild_data.GuildBigBoxData, 0, len(dAtAs.guildBigBoxData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildBigBoxData.Array = append(dAtAs.guildBigBoxData.Array, dAtAs.guildBigBoxData.Map[k])
	}
	dAtAs.guildBigBoxData.MinKeyData = dAtAs.guildBigBoxData.Array[0]
	dAtAs.guildBigBoxData.MaxKeyData = dAtAs.guildBigBoxData.Array[len(dAtAs.guildBigBoxData.Array)-1]

	// guild_data.GuildClassLevelData
	dAtAs.guildClassLevelData = &GuildClassLevelDataConfig{}
	dAtAs.guildClassLevelData.Map, dAtAs.guildClassLevelData.parserMap, err = guild_data.LoadGuildClassLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildClassLevelData.Map))
	for k := range dAtAs.guildClassLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildClassLevelData.Array = make([]*guild_data.GuildClassLevelData, 0, len(dAtAs.guildClassLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildClassLevelData.Array = append(dAtAs.guildClassLevelData.Array, dAtAs.guildClassLevelData.Map[k])
	}
	dAtAs.guildClassLevelData.MinKeyData = dAtAs.guildClassLevelData.Array[0]
	dAtAs.guildClassLevelData.MaxKeyData = dAtAs.guildClassLevelData.Array[len(dAtAs.guildClassLevelData.Array)-1]

	// guild_data.GuildClassTitleData
	dAtAs.guildClassTitleData = &GuildClassTitleDataConfig{}
	dAtAs.guildClassTitleData.Map, dAtAs.guildClassTitleData.parserMap, err = guild_data.LoadGuildClassTitleData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildClassTitleData.Map))
	for k := range dAtAs.guildClassTitleData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildClassTitleData.Array = make([]*guild_data.GuildClassTitleData, 0, len(dAtAs.guildClassTitleData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildClassTitleData.Array = append(dAtAs.guildClassTitleData.Array, dAtAs.guildClassTitleData.Map[k])
	}
	dAtAs.guildClassTitleData.MinKeyData = dAtAs.guildClassTitleData.Array[0]
	dAtAs.guildClassTitleData.MaxKeyData = dAtAs.guildClassTitleData.Array[len(dAtAs.guildClassTitleData.Array)-1]

	// guild_data.GuildDonateData
	dAtAs.guildDonateData = &GuildDonateDataConfig{}
	dAtAs.guildDonateData.Map, dAtAs.guildDonateData.parserMap, err = guild_data.LoadGuildDonateData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildDonateData.Map))
	for k := range dAtAs.guildDonateData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildDonateData.Array = make([]*guild_data.GuildDonateData, 0, len(dAtAs.guildDonateData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildDonateData.Array = append(dAtAs.guildDonateData.Array, dAtAs.guildDonateData.Map[k])
	}
	dAtAs.guildDonateData.MinKeyData = dAtAs.guildDonateData.Array[0]
	dAtAs.guildDonateData.MaxKeyData = dAtAs.guildDonateData.Array[len(dAtAs.guildDonateData.Array)-1]

	// guild_data.GuildEventPrizeData
	dAtAs.guildEventPrizeData = &GuildEventPrizeDataConfig{}
	dAtAs.guildEventPrizeData.Map, dAtAs.guildEventPrizeData.parserMap, err = guild_data.LoadGuildEventPrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildEventPrizeData.Map))
	for k := range dAtAs.guildEventPrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildEventPrizeData.Array = make([]*guild_data.GuildEventPrizeData, 0, len(dAtAs.guildEventPrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildEventPrizeData.Array = append(dAtAs.guildEventPrizeData.Array, dAtAs.guildEventPrizeData.Map[k])
	}
	dAtAs.guildEventPrizeData.MinKeyData = dAtAs.guildEventPrizeData.Array[0]
	dAtAs.guildEventPrizeData.MaxKeyData = dAtAs.guildEventPrizeData.Array[len(dAtAs.guildEventPrizeData.Array)-1]

	// guild_data.GuildLevelCdrData
	dAtAs.guildLevelCdrData = &GuildLevelCdrDataConfig{}
	dAtAs.guildLevelCdrData.Map, dAtAs.guildLevelCdrData.parserMap, err = guild_data.LoadGuildLevelCdrData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildLevelCdrData.Map))
	for k := range dAtAs.guildLevelCdrData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildLevelCdrData.Array = make([]*guild_data.GuildLevelCdrData, 0, len(dAtAs.guildLevelCdrData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildLevelCdrData.Array = append(dAtAs.guildLevelCdrData.Array, dAtAs.guildLevelCdrData.Map[k])
	}
	dAtAs.guildLevelCdrData.MinKeyData = dAtAs.guildLevelCdrData.Array[0]
	dAtAs.guildLevelCdrData.MaxKeyData = dAtAs.guildLevelCdrData.Array[len(dAtAs.guildLevelCdrData.Array)-1]

	// guild_data.GuildLevelData
	dAtAs.guildLevelData = &GuildLevelDataConfig{}
	dAtAs.guildLevelData.Map, dAtAs.guildLevelData.parserMap, err = guild_data.LoadGuildLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildLevelData.Map))
	for k := range dAtAs.guildLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildLevelData.Array = make([]*guild_data.GuildLevelData, 0, len(dAtAs.guildLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildLevelData.Array = append(dAtAs.guildLevelData.Array, dAtAs.guildLevelData.Map[k])
	}
	dAtAs.guildLevelData.MinKeyData = dAtAs.guildLevelData.Array[0]
	dAtAs.guildLevelData.MaxKeyData = dAtAs.guildLevelData.Array[len(dAtAs.guildLevelData.Array)-1]

	// guild_data.GuildLogData
	dAtAs.guildLogData = &GuildLogDataConfig{}
	dAtAs.guildLogData.Map, dAtAs.guildLogData.parserMap, err = guild_data.LoadGuildLogData(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.guildLogData.Map))
	for k := range dAtAs.guildLogData.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.guildLogData.Array = make([]*guild_data.GuildLogData, 0, len(dAtAs.guildLogData.Map))
	for _, k := range stringKeys {
		dAtAs.guildLogData.Array = append(dAtAs.guildLogData.Array, dAtAs.guildLogData.Map[k])
	}
	dAtAs.guildLogData.MinKeyData = dAtAs.guildLogData.Array[0]
	dAtAs.guildLogData.MaxKeyData = dAtAs.guildLogData.Array[len(dAtAs.guildLogData.Array)-1]

	// guild_data.GuildPermissionShowData
	dAtAs.guildPermissionShowData = &GuildPermissionShowDataConfig{}
	dAtAs.guildPermissionShowData.Map, dAtAs.guildPermissionShowData.parserMap, err = guild_data.LoadGuildPermissionShowData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildPermissionShowData.Map))
	for k := range dAtAs.guildPermissionShowData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildPermissionShowData.Array = make([]*guild_data.GuildPermissionShowData, 0, len(dAtAs.guildPermissionShowData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildPermissionShowData.Array = append(dAtAs.guildPermissionShowData.Array, dAtAs.guildPermissionShowData.Map[k])
	}
	dAtAs.guildPermissionShowData.MinKeyData = dAtAs.guildPermissionShowData.Array[0]
	dAtAs.guildPermissionShowData.MaxKeyData = dAtAs.guildPermissionShowData.Array[len(dAtAs.guildPermissionShowData.Array)-1]

	// guild_data.GuildPrestigeEventData
	dAtAs.guildPrestigeEventData = &GuildPrestigeEventDataConfig{}
	dAtAs.guildPrestigeEventData.Map, dAtAs.guildPrestigeEventData.parserMap, err = guild_data.LoadGuildPrestigeEventData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildPrestigeEventData.Map))
	for k := range dAtAs.guildPrestigeEventData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildPrestigeEventData.Array = make([]*guild_data.GuildPrestigeEventData, 0, len(dAtAs.guildPrestigeEventData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildPrestigeEventData.Array = append(dAtAs.guildPrestigeEventData.Array, dAtAs.guildPrestigeEventData.Map[k])
	}
	dAtAs.guildPrestigeEventData.MinKeyData = dAtAs.guildPrestigeEventData.Array[0]
	dAtAs.guildPrestigeEventData.MaxKeyData = dAtAs.guildPrestigeEventData.Array[len(dAtAs.guildPrestigeEventData.Array)-1]

	// guild_data.GuildPrestigePrizeData
	dAtAs.guildPrestigePrizeData = &GuildPrestigePrizeDataConfig{}
	dAtAs.guildPrestigePrizeData.Map, dAtAs.guildPrestigePrizeData.parserMap, err = guild_data.LoadGuildPrestigePrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildPrestigePrizeData.Map))
	for k := range dAtAs.guildPrestigePrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildPrestigePrizeData.Array = make([]*guild_data.GuildPrestigePrizeData, 0, len(dAtAs.guildPrestigePrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildPrestigePrizeData.Array = append(dAtAs.guildPrestigePrizeData.Array, dAtAs.guildPrestigePrizeData.Map[k])
	}
	dAtAs.guildPrestigePrizeData.MinKeyData = dAtAs.guildPrestigePrizeData.Array[0]
	dAtAs.guildPrestigePrizeData.MaxKeyData = dAtAs.guildPrestigePrizeData.Array[len(dAtAs.guildPrestigePrizeData.Array)-1]

	// guild_data.GuildRankPrizeData
	dAtAs.guildRankPrizeData = &GuildRankPrizeDataConfig{}
	dAtAs.guildRankPrizeData.Map, dAtAs.guildRankPrizeData.parserMap, err = guild_data.LoadGuildRankPrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildRankPrizeData.Map))
	for k := range dAtAs.guildRankPrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildRankPrizeData.Array = make([]*guild_data.GuildRankPrizeData, 0, len(dAtAs.guildRankPrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildRankPrizeData.Array = append(dAtAs.guildRankPrizeData.Array, dAtAs.guildRankPrizeData.Map[k])
	}
	dAtAs.guildRankPrizeData.MinKeyData = dAtAs.guildRankPrizeData.Array[0]
	dAtAs.guildRankPrizeData.MaxKeyData = dAtAs.guildRankPrizeData.Array[len(dAtAs.guildRankPrizeData.Array)-1]

	// guild_data.GuildTarget
	dAtAs.guildTarget = &GuildTargetConfig{}
	dAtAs.guildTarget.Map, dAtAs.guildTarget.parserMap, err = guild_data.LoadGuildTarget(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildTarget.Map))
	for k := range dAtAs.guildTarget.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildTarget.Array = make([]*guild_data.GuildTarget, 0, len(dAtAs.guildTarget.Map))
	for _, k := range uint64Keys {
		dAtAs.guildTarget.Array = append(dAtAs.guildTarget.Array, dAtAs.guildTarget.Map[k])
	}
	dAtAs.guildTarget.MinKeyData = dAtAs.guildTarget.Array[0]
	dAtAs.guildTarget.MaxKeyData = dAtAs.guildTarget.Array[len(dAtAs.guildTarget.Array)-1]

	// guild_data.GuildTaskData
	dAtAs.guildTaskData = &GuildTaskDataConfig{}
	dAtAs.guildTaskData.Map, dAtAs.guildTaskData.parserMap, err = guild_data.LoadGuildTaskData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildTaskData.Map))
	for k := range dAtAs.guildTaskData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildTaskData.Array = make([]*guild_data.GuildTaskData, 0, len(dAtAs.guildTaskData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildTaskData.Array = append(dAtAs.guildTaskData.Array, dAtAs.guildTaskData.Map[k])
	}
	dAtAs.guildTaskData.MinKeyData = dAtAs.guildTaskData.Array[0]
	dAtAs.guildTaskData.MaxKeyData = dAtAs.guildTaskData.Array[len(dAtAs.guildTaskData.Array)-1]

	// guild_data.GuildTaskEvaluateData
	dAtAs.guildTaskEvaluateData = &GuildTaskEvaluateDataConfig{}
	dAtAs.guildTaskEvaluateData.Map, dAtAs.guildTaskEvaluateData.parserMap, err = guild_data.LoadGuildTaskEvaluateData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildTaskEvaluateData.Map))
	for k := range dAtAs.guildTaskEvaluateData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildTaskEvaluateData.Array = make([]*guild_data.GuildTaskEvaluateData, 0, len(dAtAs.guildTaskEvaluateData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildTaskEvaluateData.Array = append(dAtAs.guildTaskEvaluateData.Array, dAtAs.guildTaskEvaluateData.Map[k])
	}
	dAtAs.guildTaskEvaluateData.MinKeyData = dAtAs.guildTaskEvaluateData.Array[0]
	dAtAs.guildTaskEvaluateData.MaxKeyData = dAtAs.guildTaskEvaluateData.Array[len(dAtAs.guildTaskEvaluateData.Array)-1]

	// guild_data.GuildTechnologyData
	dAtAs.guildTechnologyData = &GuildTechnologyDataConfig{}
	dAtAs.guildTechnologyData.Map, dAtAs.guildTechnologyData.parserMap, err = guild_data.LoadGuildTechnologyData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildTechnologyData.Map))
	for k := range dAtAs.guildTechnologyData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildTechnologyData.Array = make([]*guild_data.GuildTechnologyData, 0, len(dAtAs.guildTechnologyData.Map))
	for _, k := range uint64Keys {
		dAtAs.guildTechnologyData.Array = append(dAtAs.guildTechnologyData.Array, dAtAs.guildTechnologyData.Map[k])
	}
	dAtAs.guildTechnologyData.MinKeyData = dAtAs.guildTechnologyData.Array[0]
	dAtAs.guildTechnologyData.MaxKeyData = dAtAs.guildTechnologyData.Array[len(dAtAs.guildTechnologyData.Array)-1]

	// guild_data.NpcGuildTemplate
	dAtAs.npcGuildTemplate = &NpcGuildTemplateConfig{}
	dAtAs.npcGuildTemplate.Map, dAtAs.npcGuildTemplate.parserMap, err = guild_data.LoadNpcGuildTemplate(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.npcGuildTemplate.Map))
	for k := range dAtAs.npcGuildTemplate.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.npcGuildTemplate.Array = make([]*guild_data.NpcGuildTemplate, 0, len(dAtAs.npcGuildTemplate.Map))
	for _, k := range uint64Keys {
		dAtAs.npcGuildTemplate.Array = append(dAtAs.npcGuildTemplate.Array, dAtAs.npcGuildTemplate.Map[k])
	}
	dAtAs.npcGuildTemplate.MinKeyData = dAtAs.npcGuildTemplate.Array[0]
	dAtAs.npcGuildTemplate.MaxKeyData = dAtAs.npcGuildTemplate.Array[len(dAtAs.npcGuildTemplate.Array)-1]

	// guild_data.NpcMemberData
	dAtAs.npcMemberData = &NpcMemberDataConfig{}
	dAtAs.npcMemberData.Map, dAtAs.npcMemberData.parserMap, err = guild_data.LoadNpcMemberData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.npcMemberData.Map))
	for k := range dAtAs.npcMemberData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.npcMemberData.Array = make([]*guild_data.NpcMemberData, 0, len(dAtAs.npcMemberData.Map))
	for _, k := range uint64Keys {
		dAtAs.npcMemberData.Array = append(dAtAs.npcMemberData.Array, dAtAs.npcMemberData.Map[k])
	}
	dAtAs.npcMemberData.MinKeyData = dAtAs.npcMemberData.Array[0]
	dAtAs.npcMemberData.MaxKeyData = dAtAs.npcMemberData.Array[len(dAtAs.npcMemberData.Array)-1]

	// head.HeadData
	dAtAs.headData = &HeadDataConfig{}
	dAtAs.headData.Map, dAtAs.headData.parserMap, err = head.LoadHeadData(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.headData.Map))
	for k := range dAtAs.headData.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.headData.Array = make([]*head.HeadData, 0, len(dAtAs.headData.Map))
	for _, k := range stringKeys {
		dAtAs.headData.Array = append(dAtAs.headData.Array, dAtAs.headData.Map[k])
	}
	dAtAs.headData.MinKeyData = dAtAs.headData.Array[0]
	dAtAs.headData.MaxKeyData = dAtAs.headData.Array[len(dAtAs.headData.Array)-1]

	// hebi.HebiPrizeData
	dAtAs.hebiPrizeData = &HebiPrizeDataConfig{}
	dAtAs.hebiPrizeData.Map, dAtAs.hebiPrizeData.parserMap, err = hebi.LoadHebiPrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.hebiPrizeData.Map))
	for k := range dAtAs.hebiPrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.hebiPrizeData.Array = make([]*hebi.HebiPrizeData, 0, len(dAtAs.hebiPrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.hebiPrizeData.Array = append(dAtAs.hebiPrizeData.Array, dAtAs.hebiPrizeData.Map[k])
	}
	dAtAs.hebiPrizeData.MinKeyData = dAtAs.hebiPrizeData.Array[0]
	dAtAs.hebiPrizeData.MaxKeyData = dAtAs.hebiPrizeData.Array[len(dAtAs.hebiPrizeData.Array)-1]

	// herodata.HeroLevelData
	dAtAs.heroLevelData = &HeroLevelDataConfig{}
	dAtAs.heroLevelData.Map, dAtAs.heroLevelData.parserMap, err = herodata.LoadHeroLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.heroLevelData.Map))
	for k := range dAtAs.heroLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.heroLevelData.Array = make([]*herodata.HeroLevelData, 0, len(dAtAs.heroLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.heroLevelData.Array = append(dAtAs.heroLevelData.Array, dAtAs.heroLevelData.Map[k])
	}
	dAtAs.heroLevelData.MinKeyData = dAtAs.heroLevelData.Array[0]
	dAtAs.heroLevelData.MaxKeyData = dAtAs.heroLevelData.Array[len(dAtAs.heroLevelData.Array)-1]

	// i18n.I18nData
	dAtAs.i18nData = &I18nDataConfig{}
	dAtAs.i18nData.Map, dAtAs.i18nData.parserMap, err = i18n.LoadI18nData(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.i18nData.Map))
	for k := range dAtAs.i18nData.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.i18nData.Array = make([]*i18n.I18nData, 0, len(dAtAs.i18nData.Map))
	for _, k := range stringKeys {
		dAtAs.i18nData.Array = append(dAtAs.i18nData.Array, dAtAs.i18nData.Map[k])
	}
	dAtAs.i18nData.MinKeyData = dAtAs.i18nData.Array[0]
	dAtAs.i18nData.MaxKeyData = dAtAs.i18nData.Array[len(dAtAs.i18nData.Array)-1]

	// icon.Icon
	dAtAs.icon = &IconConfig{}
	dAtAs.icon.Map, dAtAs.icon.parserMap, err = icon.LoadIcon(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.icon.Map))
	for k := range dAtAs.icon.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.icon.Array = make([]*icon.Icon, 0, len(dAtAs.icon.Map))
	for _, k := range stringKeys {
		dAtAs.icon.Array = append(dAtAs.icon.Array, dAtAs.icon.Map[k])
	}
	dAtAs.icon.MinKeyData = dAtAs.icon.Array[0]
	dAtAs.icon.MaxKeyData = dAtAs.icon.Array[len(dAtAs.icon.Array)-1]

	// location.LocationData
	dAtAs.locationData = &LocationDataConfig{}
	dAtAs.locationData.Map, dAtAs.locationData.parserMap, err = location.LoadLocationData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.locationData.Map))
	for k := range dAtAs.locationData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.locationData.Array = make([]*location.LocationData, 0, len(dAtAs.locationData.Map))
	for _, k := range uint64Keys {
		dAtAs.locationData.Array = append(dAtAs.locationData.Array, dAtAs.locationData.Map[k])
	}
	dAtAs.locationData.MinKeyData = dAtAs.locationData.Array[0]
	dAtAs.locationData.MaxKeyData = dAtAs.locationData.Array[len(dAtAs.locationData.Array)-1]

	// maildata.MailData
	dAtAs.mailData = &MailDataConfig{}
	dAtAs.mailData.Map, dAtAs.mailData.parserMap, err = maildata.LoadMailData(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.mailData.Map))
	for k := range dAtAs.mailData.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.mailData.Array = make([]*maildata.MailData, 0, len(dAtAs.mailData.Map))
	for _, k := range stringKeys {
		dAtAs.mailData.Array = append(dAtAs.mailData.Array, dAtAs.mailData.Map[k])
	}
	dAtAs.mailData.MinKeyData = dAtAs.mailData.Array[0]
	dAtAs.mailData.MaxKeyData = dAtAs.mailData.Array[len(dAtAs.mailData.Array)-1]

	// military_data.JiuGuanData
	dAtAs.jiuGuanData = &JiuGuanDataConfig{}
	dAtAs.jiuGuanData.Map, dAtAs.jiuGuanData.parserMap, err = military_data.LoadJiuGuanData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.jiuGuanData.Map))
	for k := range dAtAs.jiuGuanData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.jiuGuanData.Array = make([]*military_data.JiuGuanData, 0, len(dAtAs.jiuGuanData.Map))
	for _, k := range uint64Keys {
		dAtAs.jiuGuanData.Array = append(dAtAs.jiuGuanData.Array, dAtAs.jiuGuanData.Map[k])
	}
	dAtAs.jiuGuanData.MinKeyData = dAtAs.jiuGuanData.Array[0]
	dAtAs.jiuGuanData.MaxKeyData = dAtAs.jiuGuanData.Array[len(dAtAs.jiuGuanData.Array)-1]

	// military_data.JunYingLevelData
	dAtAs.junYingLevelData = &JunYingLevelDataConfig{}
	dAtAs.junYingLevelData.Map, dAtAs.junYingLevelData.parserMap, err = military_data.LoadJunYingLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.junYingLevelData.Map))
	for k := range dAtAs.junYingLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.junYingLevelData.Array = make([]*military_data.JunYingLevelData, 0, len(dAtAs.junYingLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.junYingLevelData.Array = append(dAtAs.junYingLevelData.Array, dAtAs.junYingLevelData.Map[k])
	}
	dAtAs.junYingLevelData.MinKeyData = dAtAs.junYingLevelData.Array[0]
	dAtAs.junYingLevelData.MaxKeyData = dAtAs.junYingLevelData.Array[len(dAtAs.junYingLevelData.Array)-1]

	// military_data.TrainingLevelData
	dAtAs.trainingLevelData = &TrainingLevelDataConfig{}
	dAtAs.trainingLevelData.Map, dAtAs.trainingLevelData.parserMap, err = military_data.LoadTrainingLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.trainingLevelData.Map))
	for k := range dAtAs.trainingLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.trainingLevelData.Array = make([]*military_data.TrainingLevelData, 0, len(dAtAs.trainingLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.trainingLevelData.Array = append(dAtAs.trainingLevelData.Array, dAtAs.trainingLevelData.Map[k])
	}
	dAtAs.trainingLevelData.MinKeyData = dAtAs.trainingLevelData.Array[0]
	dAtAs.trainingLevelData.MaxKeyData = dAtAs.trainingLevelData.Array[len(dAtAs.trainingLevelData.Array)-1]

	// military_data.TutorData
	dAtAs.tutorData = &TutorDataConfig{}
	dAtAs.tutorData.Map, dAtAs.tutorData.parserMap, err = military_data.LoadTutorData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.tutorData.Map))
	for k := range dAtAs.tutorData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.tutorData.Array = make([]*military_data.TutorData, 0, len(dAtAs.tutorData.Map))
	for _, k := range uint64Keys {
		dAtAs.tutorData.Array = append(dAtAs.tutorData.Array, dAtAs.tutorData.Map[k])
	}
	dAtAs.tutorData.MinKeyData = dAtAs.tutorData.Array[0]
	dAtAs.tutorData.MaxKeyData = dAtAs.tutorData.Array[len(dAtAs.tutorData.Array)-1]

	// mingcdata.McBuildAddSupportData
	dAtAs.mcBuildAddSupportData = &McBuildAddSupportDataConfig{}
	dAtAs.mcBuildAddSupportData.Map, dAtAs.mcBuildAddSupportData.parserMap, err = mingcdata.LoadMcBuildAddSupportData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mcBuildAddSupportData.Map))
	for k := range dAtAs.mcBuildAddSupportData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mcBuildAddSupportData.Array = make([]*mingcdata.McBuildAddSupportData, 0, len(dAtAs.mcBuildAddSupportData.Map))
	for _, k := range uint64Keys {
		dAtAs.mcBuildAddSupportData.Array = append(dAtAs.mcBuildAddSupportData.Array, dAtAs.mcBuildAddSupportData.Map[k])
	}
	dAtAs.mcBuildAddSupportData.MinKeyData = dAtAs.mcBuildAddSupportData.Array[0]
	dAtAs.mcBuildAddSupportData.MaxKeyData = dAtAs.mcBuildAddSupportData.Array[len(dAtAs.mcBuildAddSupportData.Array)-1]

	// mingcdata.McBuildGuildMemberPrizeData
	dAtAs.mcBuildGuildMemberPrizeData = &McBuildGuildMemberPrizeDataConfig{}
	dAtAs.mcBuildGuildMemberPrizeData.Map, dAtAs.mcBuildGuildMemberPrizeData.parserMap, err = mingcdata.LoadMcBuildGuildMemberPrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mcBuildGuildMemberPrizeData.Map))
	for k := range dAtAs.mcBuildGuildMemberPrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mcBuildGuildMemberPrizeData.Array = make([]*mingcdata.McBuildGuildMemberPrizeData, 0, len(dAtAs.mcBuildGuildMemberPrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.mcBuildGuildMemberPrizeData.Array = append(dAtAs.mcBuildGuildMemberPrizeData.Array, dAtAs.mcBuildGuildMemberPrizeData.Map[k])
	}
	dAtAs.mcBuildGuildMemberPrizeData.MinKeyData = dAtAs.mcBuildGuildMemberPrizeData.Array[0]
	dAtAs.mcBuildGuildMemberPrizeData.MaxKeyData = dAtAs.mcBuildGuildMemberPrizeData.Array[len(dAtAs.mcBuildGuildMemberPrizeData.Array)-1]

	// mingcdata.McBuildMcSupportData
	dAtAs.mcBuildMcSupportData = &McBuildMcSupportDataConfig{}
	dAtAs.mcBuildMcSupportData.Map, dAtAs.mcBuildMcSupportData.parserMap, err = mingcdata.LoadMcBuildMcSupportData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mcBuildMcSupportData.Map))
	for k := range dAtAs.mcBuildMcSupportData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mcBuildMcSupportData.Array = make([]*mingcdata.McBuildMcSupportData, 0, len(dAtAs.mcBuildMcSupportData.Map))
	for _, k := range uint64Keys {
		dAtAs.mcBuildMcSupportData.Array = append(dAtAs.mcBuildMcSupportData.Array, dAtAs.mcBuildMcSupportData.Map[k])
	}
	dAtAs.mcBuildMcSupportData.MinKeyData = dAtAs.mcBuildMcSupportData.Array[0]
	dAtAs.mcBuildMcSupportData.MaxKeyData = dAtAs.mcBuildMcSupportData.Array[len(dAtAs.mcBuildMcSupportData.Array)-1]

	// mingcdata.MingcBaseData
	dAtAs.mingcBaseData = &MingcBaseDataConfig{}
	dAtAs.mingcBaseData.Map, dAtAs.mingcBaseData.parserMap, err = mingcdata.LoadMingcBaseData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcBaseData.Map))
	for k := range dAtAs.mingcBaseData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcBaseData.Array = make([]*mingcdata.MingcBaseData, 0, len(dAtAs.mingcBaseData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcBaseData.Array = append(dAtAs.mingcBaseData.Array, dAtAs.mingcBaseData.Map[k])
	}
	dAtAs.mingcBaseData.MinKeyData = dAtAs.mingcBaseData.Array[0]
	dAtAs.mingcBaseData.MaxKeyData = dAtAs.mingcBaseData.Array[len(dAtAs.mingcBaseData.Array)-1]

	// mingcdata.MingcTimeData
	dAtAs.mingcTimeData = &MingcTimeDataConfig{}
	dAtAs.mingcTimeData.Map, dAtAs.mingcTimeData.parserMap, err = mingcdata.LoadMingcTimeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcTimeData.Map))
	for k := range dAtAs.mingcTimeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcTimeData.Array = make([]*mingcdata.MingcTimeData, 0, len(dAtAs.mingcTimeData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcTimeData.Array = append(dAtAs.mingcTimeData.Array, dAtAs.mingcTimeData.Map[k])
	}
	dAtAs.mingcTimeData.MinKeyData = dAtAs.mingcTimeData.Array[0]
	dAtAs.mingcTimeData.MaxKeyData = dAtAs.mingcTimeData.Array[len(dAtAs.mingcTimeData.Array)-1]

	// mingcdata.MingcWarBuildingData
	dAtAs.mingcWarBuildingData = &MingcWarBuildingDataConfig{}
	dAtAs.mingcWarBuildingData.Map, dAtAs.mingcWarBuildingData.parserMap, err = mingcdata.LoadMingcWarBuildingData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcWarBuildingData.Map))
	for k := range dAtAs.mingcWarBuildingData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcWarBuildingData.Array = make([]*mingcdata.MingcWarBuildingData, 0, len(dAtAs.mingcWarBuildingData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcWarBuildingData.Array = append(dAtAs.mingcWarBuildingData.Array, dAtAs.mingcWarBuildingData.Map[k])
	}
	dAtAs.mingcWarBuildingData.MinKeyData = dAtAs.mingcWarBuildingData.Array[0]
	dAtAs.mingcWarBuildingData.MaxKeyData = dAtAs.mingcWarBuildingData.Array[len(dAtAs.mingcWarBuildingData.Array)-1]

	// mingcdata.MingcWarDrumStatData
	dAtAs.mingcWarDrumStatData = &MingcWarDrumStatDataConfig{}
	dAtAs.mingcWarDrumStatData.Map, dAtAs.mingcWarDrumStatData.parserMap, err = mingcdata.LoadMingcWarDrumStatData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcWarDrumStatData.Map))
	for k := range dAtAs.mingcWarDrumStatData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcWarDrumStatData.Array = make([]*mingcdata.MingcWarDrumStatData, 0, len(dAtAs.mingcWarDrumStatData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcWarDrumStatData.Array = append(dAtAs.mingcWarDrumStatData.Array, dAtAs.mingcWarDrumStatData.Map[k])
	}
	dAtAs.mingcWarDrumStatData.MinKeyData = dAtAs.mingcWarDrumStatData.Array[0]
	dAtAs.mingcWarDrumStatData.MaxKeyData = dAtAs.mingcWarDrumStatData.Array[len(dAtAs.mingcWarDrumStatData.Array)-1]

	// mingcdata.MingcWarMapData
	dAtAs.mingcWarMapData = &MingcWarMapDataConfig{}
	dAtAs.mingcWarMapData.Map, dAtAs.mingcWarMapData.parserMap, err = mingcdata.LoadMingcWarMapData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcWarMapData.Map))
	for k := range dAtAs.mingcWarMapData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcWarMapData.Array = make([]*mingcdata.MingcWarMapData, 0, len(dAtAs.mingcWarMapData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcWarMapData.Array = append(dAtAs.mingcWarMapData.Array, dAtAs.mingcWarMapData.Map[k])
	}
	dAtAs.mingcWarMapData.MinKeyData = dAtAs.mingcWarMapData.Array[0]
	dAtAs.mingcWarMapData.MaxKeyData = dAtAs.mingcWarMapData.Array[len(dAtAs.mingcWarMapData.Array)-1]

	// mingcdata.MingcWarMultiKillData
	dAtAs.mingcWarMultiKillData = &MingcWarMultiKillDataConfig{}
	dAtAs.mingcWarMultiKillData.Map, dAtAs.mingcWarMultiKillData.parserMap, err = mingcdata.LoadMingcWarMultiKillData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcWarMultiKillData.Map))
	for k := range dAtAs.mingcWarMultiKillData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcWarMultiKillData.Array = make([]*mingcdata.MingcWarMultiKillData, 0, len(dAtAs.mingcWarMultiKillData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcWarMultiKillData.Array = append(dAtAs.mingcWarMultiKillData.Array, dAtAs.mingcWarMultiKillData.Map[k])
	}
	dAtAs.mingcWarMultiKillData.MinKeyData = dAtAs.mingcWarMultiKillData.Array[0]
	dAtAs.mingcWarMultiKillData.MaxKeyData = dAtAs.mingcWarMultiKillData.Array[len(dAtAs.mingcWarMultiKillData.Array)-1]

	// mingcdata.MingcWarNpcData
	dAtAs.mingcWarNpcData = &MingcWarNpcDataConfig{}
	dAtAs.mingcWarNpcData.Map, dAtAs.mingcWarNpcData.parserMap, err = mingcdata.LoadMingcWarNpcData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcWarNpcData.Map))
	for k := range dAtAs.mingcWarNpcData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcWarNpcData.Array = make([]*mingcdata.MingcWarNpcData, 0, len(dAtAs.mingcWarNpcData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcWarNpcData.Array = append(dAtAs.mingcWarNpcData.Array, dAtAs.mingcWarNpcData.Map[k])
	}
	dAtAs.mingcWarNpcData.MinKeyData = dAtAs.mingcWarNpcData.Array[0]
	dAtAs.mingcWarNpcData.MaxKeyData = dAtAs.mingcWarNpcData.Array[len(dAtAs.mingcWarNpcData.Array)-1]

	// mingcdata.MingcWarNpcGuildData
	dAtAs.mingcWarNpcGuildData = &MingcWarNpcGuildDataConfig{}
	dAtAs.mingcWarNpcGuildData.Map, dAtAs.mingcWarNpcGuildData.parserMap, err = mingcdata.LoadMingcWarNpcGuildData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcWarNpcGuildData.Map))
	for k := range dAtAs.mingcWarNpcGuildData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcWarNpcGuildData.Array = make([]*mingcdata.MingcWarNpcGuildData, 0, len(dAtAs.mingcWarNpcGuildData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcWarNpcGuildData.Array = append(dAtAs.mingcWarNpcGuildData.Array, dAtAs.mingcWarNpcGuildData.Map[k])
	}
	dAtAs.mingcWarNpcGuildData.MinKeyData = dAtAs.mingcWarNpcGuildData.Array[0]
	dAtAs.mingcWarNpcGuildData.MaxKeyData = dAtAs.mingcWarNpcGuildData.Array[len(dAtAs.mingcWarNpcGuildData.Array)-1]

	// mingcdata.MingcWarSceneData
	dAtAs.mingcWarSceneData = &MingcWarSceneDataConfig{}
	dAtAs.mingcWarSceneData.Map, dAtAs.mingcWarSceneData.parserMap, err = mingcdata.LoadMingcWarSceneData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcWarSceneData.Map))
	for k := range dAtAs.mingcWarSceneData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcWarSceneData.Array = make([]*mingcdata.MingcWarSceneData, 0, len(dAtAs.mingcWarSceneData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcWarSceneData.Array = append(dAtAs.mingcWarSceneData.Array, dAtAs.mingcWarSceneData.Map[k])
	}
	dAtAs.mingcWarSceneData.MinKeyData = dAtAs.mingcWarSceneData.Array[0]
	dAtAs.mingcWarSceneData.MaxKeyData = dAtAs.mingcWarSceneData.Array[len(dAtAs.mingcWarSceneData.Array)-1]

	// mingcdata.MingcWarTouShiBuildingTargetData
	dAtAs.mingcWarTouShiBuildingTargetData = &MingcWarTouShiBuildingTargetDataConfig{}
	dAtAs.mingcWarTouShiBuildingTargetData.Map, dAtAs.mingcWarTouShiBuildingTargetData.parserMap, err = mingcdata.LoadMingcWarTouShiBuildingTargetData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcWarTouShiBuildingTargetData.Map))
	for k := range dAtAs.mingcWarTouShiBuildingTargetData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcWarTouShiBuildingTargetData.Array = make([]*mingcdata.MingcWarTouShiBuildingTargetData, 0, len(dAtAs.mingcWarTouShiBuildingTargetData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcWarTouShiBuildingTargetData.Array = append(dAtAs.mingcWarTouShiBuildingTargetData.Array, dAtAs.mingcWarTouShiBuildingTargetData.Map[k])
	}
	dAtAs.mingcWarTouShiBuildingTargetData.MinKeyData = dAtAs.mingcWarTouShiBuildingTargetData.Array[0]
	dAtAs.mingcWarTouShiBuildingTargetData.MaxKeyData = dAtAs.mingcWarTouShiBuildingTargetData.Array[len(dAtAs.mingcWarTouShiBuildingTargetData.Array)-1]

	// mingcdata.MingcWarTroopLastBeatWhenFailData
	dAtAs.mingcWarTroopLastBeatWhenFailData = &MingcWarTroopLastBeatWhenFailDataConfig{}
	dAtAs.mingcWarTroopLastBeatWhenFailData.Map, dAtAs.mingcWarTroopLastBeatWhenFailData.parserMap, err = mingcdata.LoadMingcWarTroopLastBeatWhenFailData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mingcWarTroopLastBeatWhenFailData.Map))
	for k := range dAtAs.mingcWarTroopLastBeatWhenFailData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mingcWarTroopLastBeatWhenFailData.Array = make([]*mingcdata.MingcWarTroopLastBeatWhenFailData, 0, len(dAtAs.mingcWarTroopLastBeatWhenFailData.Map))
	for _, k := range uint64Keys {
		dAtAs.mingcWarTroopLastBeatWhenFailData.Array = append(dAtAs.mingcWarTroopLastBeatWhenFailData.Array, dAtAs.mingcWarTroopLastBeatWhenFailData.Map[k])
	}
	dAtAs.mingcWarTroopLastBeatWhenFailData.MinKeyData = dAtAs.mingcWarTroopLastBeatWhenFailData.Array[0]
	dAtAs.mingcWarTroopLastBeatWhenFailData.MaxKeyData = dAtAs.mingcWarTroopLastBeatWhenFailData.Array[len(dAtAs.mingcWarTroopLastBeatWhenFailData.Array)-1]

	// monsterdata.MonsterCaptainData
	dAtAs.monsterCaptainData = &MonsterCaptainDataConfig{}
	dAtAs.monsterCaptainData.Map, dAtAs.monsterCaptainData.parserMap, err = monsterdata.LoadMonsterCaptainData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.monsterCaptainData.Map))
	for k := range dAtAs.monsterCaptainData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.monsterCaptainData.Array = make([]*monsterdata.MonsterCaptainData, 0, len(dAtAs.monsterCaptainData.Map))
	for _, k := range uint64Keys {
		dAtAs.monsterCaptainData.Array = append(dAtAs.monsterCaptainData.Array, dAtAs.monsterCaptainData.Map[k])
	}
	dAtAs.monsterCaptainData.MinKeyData = dAtAs.monsterCaptainData.Array[0]
	dAtAs.monsterCaptainData.MaxKeyData = dAtAs.monsterCaptainData.Array[len(dAtAs.monsterCaptainData.Array)-1]

	// monsterdata.MonsterMasterData
	dAtAs.monsterMasterData = &MonsterMasterDataConfig{}
	dAtAs.monsterMasterData.Map, dAtAs.monsterMasterData.parserMap, err = monsterdata.LoadMonsterMasterData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.monsterMasterData.Map))
	for k := range dAtAs.monsterMasterData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.monsterMasterData.Array = make([]*monsterdata.MonsterMasterData, 0, len(dAtAs.monsterMasterData.Map))
	for _, k := range uint64Keys {
		dAtAs.monsterMasterData.Array = append(dAtAs.monsterMasterData.Array, dAtAs.monsterMasterData.Map[k])
	}
	dAtAs.monsterMasterData.MinKeyData = dAtAs.monsterMasterData.Array[0]
	dAtAs.monsterMasterData.MaxKeyData = dAtAs.monsterMasterData.Array[len(dAtAs.monsterMasterData.Array)-1]

	// promdata.DailyBargainData
	dAtAs.dailyBargainData = &DailyBargainDataConfig{}
	dAtAs.dailyBargainData.Map, dAtAs.dailyBargainData.parserMap, err = promdata.LoadDailyBargainData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.dailyBargainData.Map))
	for k := range dAtAs.dailyBargainData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.dailyBargainData.Array = make([]*promdata.DailyBargainData, 0, len(dAtAs.dailyBargainData.Map))
	for _, k := range uint64Keys {
		dAtAs.dailyBargainData.Array = append(dAtAs.dailyBargainData.Array, dAtAs.dailyBargainData.Map[k])
	}
	dAtAs.dailyBargainData.MinKeyData = dAtAs.dailyBargainData.Array[0]
	dAtAs.dailyBargainData.MaxKeyData = dAtAs.dailyBargainData.Array[len(dAtAs.dailyBargainData.Array)-1]

	// promdata.DurationCardData
	dAtAs.durationCardData = &DurationCardDataConfig{}
	dAtAs.durationCardData.Map, dAtAs.durationCardData.parserMap, err = promdata.LoadDurationCardData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.durationCardData.Map))
	for k := range dAtAs.durationCardData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.durationCardData.Array = make([]*promdata.DurationCardData, 0, len(dAtAs.durationCardData.Map))
	for _, k := range uint64Keys {
		dAtAs.durationCardData.Array = append(dAtAs.durationCardData.Array, dAtAs.durationCardData.Map[k])
	}
	dAtAs.durationCardData.MinKeyData = dAtAs.durationCardData.Array[0]
	dAtAs.durationCardData.MaxKeyData = dAtAs.durationCardData.Array[len(dAtAs.durationCardData.Array)-1]

	// promdata.EventLimitGiftData
	dAtAs.eventLimitGiftData = &EventLimitGiftDataConfig{}
	dAtAs.eventLimitGiftData.Map, dAtAs.eventLimitGiftData.parserMap, err = promdata.LoadEventLimitGiftData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.eventLimitGiftData.Map))
	for k := range dAtAs.eventLimitGiftData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.eventLimitGiftData.Array = make([]*promdata.EventLimitGiftData, 0, len(dAtAs.eventLimitGiftData.Map))
	for _, k := range uint64Keys {
		dAtAs.eventLimitGiftData.Array = append(dAtAs.eventLimitGiftData.Array, dAtAs.eventLimitGiftData.Map[k])
	}
	dAtAs.eventLimitGiftData.MinKeyData = dAtAs.eventLimitGiftData.Array[0]
	dAtAs.eventLimitGiftData.MaxKeyData = dAtAs.eventLimitGiftData.Array[len(dAtAs.eventLimitGiftData.Array)-1]

	// promdata.FreeGiftData
	dAtAs.freeGiftData = &FreeGiftDataConfig{}
	dAtAs.freeGiftData.Map, dAtAs.freeGiftData.parserMap, err = promdata.LoadFreeGiftData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.freeGiftData.Map))
	for k := range dAtAs.freeGiftData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.freeGiftData.Array = make([]*promdata.FreeGiftData, 0, len(dAtAs.freeGiftData.Map))
	for _, k := range uint64Keys {
		dAtAs.freeGiftData.Array = append(dAtAs.freeGiftData.Array, dAtAs.freeGiftData.Map[k])
	}
	dAtAs.freeGiftData.MinKeyData = dAtAs.freeGiftData.Array[0]
	dAtAs.freeGiftData.MaxKeyData = dAtAs.freeGiftData.Array[len(dAtAs.freeGiftData.Array)-1]

	// promdata.HeroLevelFundData
	dAtAs.heroLevelFundData = &HeroLevelFundDataConfig{}
	dAtAs.heroLevelFundData.Map, dAtAs.heroLevelFundData.parserMap, err = promdata.LoadHeroLevelFundData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.heroLevelFundData.Map))
	for k := range dAtAs.heroLevelFundData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.heroLevelFundData.Array = make([]*promdata.HeroLevelFundData, 0, len(dAtAs.heroLevelFundData.Map))
	for _, k := range uint64Keys {
		dAtAs.heroLevelFundData.Array = append(dAtAs.heroLevelFundData.Array, dAtAs.heroLevelFundData.Map[k])
	}
	dAtAs.heroLevelFundData.MinKeyData = dAtAs.heroLevelFundData.Array[0]
	dAtAs.heroLevelFundData.MaxKeyData = dAtAs.heroLevelFundData.Array[len(dAtAs.heroLevelFundData.Array)-1]

	// promdata.LoginDayData
	dAtAs.loginDayData = &LoginDayDataConfig{}
	dAtAs.loginDayData.Map, dAtAs.loginDayData.parserMap, err = promdata.LoadLoginDayData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.loginDayData.Map))
	for k := range dAtAs.loginDayData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.loginDayData.Array = make([]*promdata.LoginDayData, 0, len(dAtAs.loginDayData.Map))
	for _, k := range uint64Keys {
		dAtAs.loginDayData.Array = append(dAtAs.loginDayData.Array, dAtAs.loginDayData.Map[k])
	}
	dAtAs.loginDayData.MinKeyData = dAtAs.loginDayData.Array[0]
	dAtAs.loginDayData.MaxKeyData = dAtAs.loginDayData.Array[len(dAtAs.loginDayData.Array)-1]

	// promdata.SpCollectionData
	dAtAs.spCollectionData = &SpCollectionDataConfig{}
	dAtAs.spCollectionData.Map, dAtAs.spCollectionData.parserMap, err = promdata.LoadSpCollectionData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.spCollectionData.Map))
	for k := range dAtAs.spCollectionData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.spCollectionData.Array = make([]*promdata.SpCollectionData, 0, len(dAtAs.spCollectionData.Map))
	for _, k := range uint64Keys {
		dAtAs.spCollectionData.Array = append(dAtAs.spCollectionData.Array, dAtAs.spCollectionData.Map[k])
	}
	dAtAs.spCollectionData.MinKeyData = dAtAs.spCollectionData.Array[0]
	dAtAs.spCollectionData.MaxKeyData = dAtAs.spCollectionData.Array[len(dAtAs.spCollectionData.Array)-1]

	// promdata.TimeLimitGiftData
	dAtAs.timeLimitGiftData = &TimeLimitGiftDataConfig{}
	dAtAs.timeLimitGiftData.Map, dAtAs.timeLimitGiftData.parserMap, err = promdata.LoadTimeLimitGiftData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.timeLimitGiftData.Map))
	for k := range dAtAs.timeLimitGiftData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.timeLimitGiftData.Array = make([]*promdata.TimeLimitGiftData, 0, len(dAtAs.timeLimitGiftData.Map))
	for _, k := range uint64Keys {
		dAtAs.timeLimitGiftData.Array = append(dAtAs.timeLimitGiftData.Array, dAtAs.timeLimitGiftData.Map[k])
	}
	dAtAs.timeLimitGiftData.MinKeyData = dAtAs.timeLimitGiftData.Array[0]
	dAtAs.timeLimitGiftData.MaxKeyData = dAtAs.timeLimitGiftData.Array[len(dAtAs.timeLimitGiftData.Array)-1]

	// promdata.TimeLimitGiftGroupData
	dAtAs.timeLimitGiftGroupData = &TimeLimitGiftGroupDataConfig{}
	dAtAs.timeLimitGiftGroupData.Map, dAtAs.timeLimitGiftGroupData.parserMap, err = promdata.LoadTimeLimitGiftGroupData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.timeLimitGiftGroupData.Map))
	for k := range dAtAs.timeLimitGiftGroupData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.timeLimitGiftGroupData.Array = make([]*promdata.TimeLimitGiftGroupData, 0, len(dAtAs.timeLimitGiftGroupData.Map))
	for _, k := range uint64Keys {
		dAtAs.timeLimitGiftGroupData.Array = append(dAtAs.timeLimitGiftGroupData.Array, dAtAs.timeLimitGiftGroupData.Map[k])
	}
	dAtAs.timeLimitGiftGroupData.MinKeyData = dAtAs.timeLimitGiftGroupData.Array[0]
	dAtAs.timeLimitGiftGroupData.MaxKeyData = dAtAs.timeLimitGiftGroupData.Array[len(dAtAs.timeLimitGiftGroupData.Array)-1]

	// pushdata.PushData
	dAtAs.pushData = &PushDataConfig{}
	dAtAs.pushData.Map, dAtAs.pushData.parserMap, err = pushdata.LoadPushData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.pushData.Map))
	for k := range dAtAs.pushData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.pushData.Array = make([]*pushdata.PushData, 0, len(dAtAs.pushData.Map))
	for _, k := range uint64Keys {
		dAtAs.pushData.Array = append(dAtAs.pushData.Array, dAtAs.pushData.Map[k])
	}
	dAtAs.pushData.MinKeyData = dAtAs.pushData.Array[0]
	dAtAs.pushData.MaxKeyData = dAtAs.pushData.Array[len(dAtAs.pushData.Array)-1]

	// pvetroop.PveTroopData
	dAtAs.pveTroopData = &PveTroopDataConfig{}
	dAtAs.pveTroopData.Map, dAtAs.pveTroopData.parserMap, err = pvetroop.LoadPveTroopData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.pveTroopData.Map))
	for k := range dAtAs.pveTroopData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.pveTroopData.Array = make([]*pvetroop.PveTroopData, 0, len(dAtAs.pveTroopData.Map))
	for _, k := range uint64Keys {
		dAtAs.pveTroopData.Array = append(dAtAs.pveTroopData.Array, dAtAs.pveTroopData.Map[k])
	}
	dAtAs.pveTroopData.MinKeyData = dAtAs.pveTroopData.Array[0]
	dAtAs.pveTroopData.MaxKeyData = dAtAs.pveTroopData.Array[len(dAtAs.pveTroopData.Array)-1]

	// question.QuestionData
	dAtAs.questionData = &QuestionDataConfig{}
	dAtAs.questionData.Map, dAtAs.questionData.parserMap, err = question.LoadQuestionData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.questionData.Map))
	for k := range dAtAs.questionData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.questionData.Array = make([]*question.QuestionData, 0, len(dAtAs.questionData.Map))
	for _, k := range uint64Keys {
		dAtAs.questionData.Array = append(dAtAs.questionData.Array, dAtAs.questionData.Map[k])
	}
	dAtAs.questionData.MinKeyData = dAtAs.questionData.Array[0]
	dAtAs.questionData.MaxKeyData = dAtAs.questionData.Array[len(dAtAs.questionData.Array)-1]

	// question.QuestionPrizeData
	dAtAs.questionPrizeData = &QuestionPrizeDataConfig{}
	dAtAs.questionPrizeData.Map, dAtAs.questionPrizeData.parserMap, err = question.LoadQuestionPrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.questionPrizeData.Map))
	for k := range dAtAs.questionPrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.questionPrizeData.Array = make([]*question.QuestionPrizeData, 0, len(dAtAs.questionPrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.questionPrizeData.Array = append(dAtAs.questionPrizeData.Array, dAtAs.questionPrizeData.Map[k])
	}
	dAtAs.questionPrizeData.MinKeyData = dAtAs.questionPrizeData.Array[0]
	dAtAs.questionPrizeData.MaxKeyData = dAtAs.questionPrizeData.Array[len(dAtAs.questionPrizeData.Array)-1]

	// question.QuestionSayingData
	dAtAs.questionSayingData = &QuestionSayingDataConfig{}
	dAtAs.questionSayingData.Map, dAtAs.questionSayingData.parserMap, err = question.LoadQuestionSayingData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.questionSayingData.Map))
	for k := range dAtAs.questionSayingData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.questionSayingData.Array = make([]*question.QuestionSayingData, 0, len(dAtAs.questionSayingData.Map))
	for _, k := range uint64Keys {
		dAtAs.questionSayingData.Array = append(dAtAs.questionSayingData.Array, dAtAs.questionSayingData.Map[k])
	}
	dAtAs.questionSayingData.MinKeyData = dAtAs.questionSayingData.Array[0]
	dAtAs.questionSayingData.MaxKeyData = dAtAs.questionSayingData.Array[len(dAtAs.questionSayingData.Array)-1]

	// race.RaceData
	dAtAs.raceData = &RaceDataConfig{}
	dAtAs.raceData.Map, dAtAs.raceData.parserMap, err = race.LoadRaceData(gos)
	if err != nil {
		return nil, err
	}

	intKeys = make([]int, 0, len(dAtAs.raceData.Map))
	for k := range dAtAs.raceData.Map {
		intKeys = append(intKeys, k)
	}
	sort.Sort(intSlice(intKeys))
	dAtAs.raceData.Array = make([]*race.RaceData, 0, len(dAtAs.raceData.Map))
	for _, k := range intKeys {
		dAtAs.raceData.Array = append(dAtAs.raceData.Array, dAtAs.raceData.Map[k])
	}
	dAtAs.raceData.MinKeyData = dAtAs.raceData.Array[0]
	dAtAs.raceData.MaxKeyData = dAtAs.raceData.Array[len(dAtAs.raceData.Array)-1]

	// random_event.EventOptionData
	dAtAs.eventOptionData = &EventOptionDataConfig{}
	dAtAs.eventOptionData.Map, dAtAs.eventOptionData.parserMap, err = random_event.LoadEventOptionData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.eventOptionData.Map))
	for k := range dAtAs.eventOptionData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.eventOptionData.Array = make([]*random_event.EventOptionData, 0, len(dAtAs.eventOptionData.Map))
	for _, k := range uint64Keys {
		dAtAs.eventOptionData.Array = append(dAtAs.eventOptionData.Array, dAtAs.eventOptionData.Map[k])
	}
	dAtAs.eventOptionData.MinKeyData = dAtAs.eventOptionData.Array[0]
	dAtAs.eventOptionData.MaxKeyData = dAtAs.eventOptionData.Array[len(dAtAs.eventOptionData.Array)-1]

	// random_event.EventPosition
	dAtAs.eventPosition = &EventPositionConfig{}
	dAtAs.eventPosition.Map, dAtAs.eventPosition.parserMap, err = random_event.LoadEventPosition(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.eventPosition.Map))
	for k := range dAtAs.eventPosition.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.eventPosition.Array = make([]*random_event.EventPosition, 0, len(dAtAs.eventPosition.Map))
	for _, k := range uint64Keys {
		dAtAs.eventPosition.Array = append(dAtAs.eventPosition.Array, dAtAs.eventPosition.Map[k])
	}
	dAtAs.eventPosition.MinKeyData = dAtAs.eventPosition.Array[0]
	dAtAs.eventPosition.MaxKeyData = dAtAs.eventPosition.Array[len(dAtAs.eventPosition.Array)-1]

	// random_event.OptionPrize
	dAtAs.optionPrize = &OptionPrizeConfig{}
	dAtAs.optionPrize.Map, dAtAs.optionPrize.parserMap, err = random_event.LoadOptionPrize(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.optionPrize.Map))
	for k := range dAtAs.optionPrize.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.optionPrize.Array = make([]*random_event.OptionPrize, 0, len(dAtAs.optionPrize.Map))
	for _, k := range uint64Keys {
		dAtAs.optionPrize.Array = append(dAtAs.optionPrize.Array, dAtAs.optionPrize.Map[k])
	}
	dAtAs.optionPrize.MinKeyData = dAtAs.optionPrize.Array[0]
	dAtAs.optionPrize.MaxKeyData = dAtAs.optionPrize.Array[len(dAtAs.optionPrize.Array)-1]

	// random_event.RandomEventData
	dAtAs.randomEventData = &RandomEventDataConfig{}
	dAtAs.randomEventData.Map, dAtAs.randomEventData.parserMap, err = random_event.LoadRandomEventData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.randomEventData.Map))
	for k := range dAtAs.randomEventData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.randomEventData.Array = make([]*random_event.RandomEventData, 0, len(dAtAs.randomEventData.Map))
	for _, k := range uint64Keys {
		dAtAs.randomEventData.Array = append(dAtAs.randomEventData.Array, dAtAs.randomEventData.Map[k])
	}
	dAtAs.randomEventData.MinKeyData = dAtAs.randomEventData.Array[0]
	dAtAs.randomEventData.MaxKeyData = dAtAs.randomEventData.Array[len(dAtAs.randomEventData.Array)-1]

	// red_packet.RedPacketData
	dAtAs.redPacketData = &RedPacketDataConfig{}
	dAtAs.redPacketData.Map, dAtAs.redPacketData.parserMap, err = red_packet.LoadRedPacketData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.redPacketData.Map))
	for k := range dAtAs.redPacketData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.redPacketData.Array = make([]*red_packet.RedPacketData, 0, len(dAtAs.redPacketData.Map))
	for _, k := range uint64Keys {
		dAtAs.redPacketData.Array = append(dAtAs.redPacketData.Array, dAtAs.redPacketData.Map[k])
	}
	dAtAs.redPacketData.MinKeyData = dAtAs.redPacketData.Array[0]
	dAtAs.redPacketData.MaxKeyData = dAtAs.redPacketData.Array[len(dAtAs.redPacketData.Array)-1]

	// regdata.AreaData
	dAtAs.areaData = &AreaDataConfig{}
	dAtAs.areaData.Map, dAtAs.areaData.parserMap, err = regdata.LoadAreaData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.areaData.Map))
	for k := range dAtAs.areaData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.areaData.Array = make([]*regdata.AreaData, 0, len(dAtAs.areaData.Map))
	for _, k := range uint64Keys {
		dAtAs.areaData.Array = append(dAtAs.areaData.Array, dAtAs.areaData.Map[k])
	}
	dAtAs.areaData.MinKeyData = dAtAs.areaData.Array[0]
	dAtAs.areaData.MaxKeyData = dAtAs.areaData.Array[len(dAtAs.areaData.Array)-1]

	// regdata.AssemblyData
	dAtAs.assemblyData = &AssemblyDataConfig{}
	dAtAs.assemblyData.Map, dAtAs.assemblyData.parserMap, err = regdata.LoadAssemblyData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.assemblyData.Map))
	for k := range dAtAs.assemblyData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.assemblyData.Array = make([]*regdata.AssemblyData, 0, len(dAtAs.assemblyData.Map))
	for _, k := range uint64Keys {
		dAtAs.assemblyData.Array = append(dAtAs.assemblyData.Array, dAtAs.assemblyData.Map[k])
	}
	dAtAs.assemblyData.MinKeyData = dAtAs.assemblyData.Array[0]
	dAtAs.assemblyData.MaxKeyData = dAtAs.assemblyData.Array[len(dAtAs.assemblyData.Array)-1]

	// regdata.BaozNpcData
	dAtAs.baozNpcData = &BaozNpcDataConfig{}
	dAtAs.baozNpcData.Map, dAtAs.baozNpcData.parserMap, err = regdata.LoadBaozNpcData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.baozNpcData.Map))
	for k := range dAtAs.baozNpcData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.baozNpcData.Array = make([]*regdata.BaozNpcData, 0, len(dAtAs.baozNpcData.Map))
	for _, k := range uint64Keys {
		dAtAs.baozNpcData.Array = append(dAtAs.baozNpcData.Array, dAtAs.baozNpcData.Map[k])
	}
	dAtAs.baozNpcData.MinKeyData = dAtAs.baozNpcData.Array[0]
	dAtAs.baozNpcData.MaxKeyData = dAtAs.baozNpcData.Array[len(dAtAs.baozNpcData.Array)-1]

	// regdata.JunTuanNpcData
	dAtAs.junTuanNpcData = &JunTuanNpcDataConfig{}
	dAtAs.junTuanNpcData.Map, dAtAs.junTuanNpcData.parserMap, err = regdata.LoadJunTuanNpcData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.junTuanNpcData.Map))
	for k := range dAtAs.junTuanNpcData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.junTuanNpcData.Array = make([]*regdata.JunTuanNpcData, 0, len(dAtAs.junTuanNpcData.Map))
	for _, k := range uint64Keys {
		dAtAs.junTuanNpcData.Array = append(dAtAs.junTuanNpcData.Array, dAtAs.junTuanNpcData.Map[k])
	}
	dAtAs.junTuanNpcData.MinKeyData = dAtAs.junTuanNpcData.Array[0]
	dAtAs.junTuanNpcData.MaxKeyData = dAtAs.junTuanNpcData.Array[len(dAtAs.junTuanNpcData.Array)-1]

	// regdata.JunTuanNpcPlaceData
	dAtAs.junTuanNpcPlaceData = &JunTuanNpcPlaceDataConfig{}
	dAtAs.junTuanNpcPlaceData.Map, dAtAs.junTuanNpcPlaceData.parserMap, err = regdata.LoadJunTuanNpcPlaceData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.junTuanNpcPlaceData.Map))
	for k := range dAtAs.junTuanNpcPlaceData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.junTuanNpcPlaceData.Array = make([]*regdata.JunTuanNpcPlaceData, 0, len(dAtAs.junTuanNpcPlaceData.Map))
	for _, k := range uint64Keys {
		dAtAs.junTuanNpcPlaceData.Array = append(dAtAs.junTuanNpcPlaceData.Array, dAtAs.junTuanNpcPlaceData.Map[k])
	}
	dAtAs.junTuanNpcPlaceData.MinKeyData = dAtAs.junTuanNpcPlaceData.Array[0]
	dAtAs.junTuanNpcPlaceData.MaxKeyData = dAtAs.junTuanNpcPlaceData.Array[len(dAtAs.junTuanNpcPlaceData.Array)-1]

	// regdata.RegionAreaData
	dAtAs.regionAreaData = &RegionAreaDataConfig{}
	dAtAs.regionAreaData.Map, dAtAs.regionAreaData.parserMap, err = regdata.LoadRegionAreaData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.regionAreaData.Map))
	for k := range dAtAs.regionAreaData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.regionAreaData.Array = make([]*regdata.RegionAreaData, 0, len(dAtAs.regionAreaData.Map))
	for _, k := range uint64Keys {
		dAtAs.regionAreaData.Array = append(dAtAs.regionAreaData.Array, dAtAs.regionAreaData.Map[k])
	}
	dAtAs.regionAreaData.MinKeyData = dAtAs.regionAreaData.Array[0]
	dAtAs.regionAreaData.MaxKeyData = dAtAs.regionAreaData.Array[len(dAtAs.regionAreaData.Array)-1]

	// regdata.RegionData
	dAtAs.regionData = &RegionDataConfig{}
	dAtAs.regionData.Map, dAtAs.regionData.parserMap, err = regdata.LoadRegionData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.regionData.Map))
	for k := range dAtAs.regionData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.regionData.Array = make([]*regdata.RegionData, 0, len(dAtAs.regionData.Map))
	for _, k := range uint64Keys {
		dAtAs.regionData.Array = append(dAtAs.regionData.Array, dAtAs.regionData.Map[k])
	}
	dAtAs.regionData.MinKeyData = dAtAs.regionData.Array[0]
	dAtAs.regionData.MaxKeyData = dAtAs.regionData.Array[len(dAtAs.regionData.Array)-1]

	// regdata.RegionMonsterData
	dAtAs.regionMonsterData = &RegionMonsterDataConfig{}
	dAtAs.regionMonsterData.Map, dAtAs.regionMonsterData.parserMap, err = regdata.LoadRegionMonsterData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.regionMonsterData.Map))
	for k := range dAtAs.regionMonsterData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.regionMonsterData.Array = make([]*regdata.RegionMonsterData, 0, len(dAtAs.regionMonsterData.Map))
	for _, k := range uint64Keys {
		dAtAs.regionMonsterData.Array = append(dAtAs.regionMonsterData.Array, dAtAs.regionMonsterData.Map[k])
	}
	dAtAs.regionMonsterData.MinKeyData = dAtAs.regionMonsterData.Array[0]
	dAtAs.regionMonsterData.MaxKeyData = dAtAs.regionMonsterData.Array[len(dAtAs.regionMonsterData.Array)-1]

	// regdata.RegionMultiLevelNpcData
	dAtAs.regionMultiLevelNpcData = &RegionMultiLevelNpcDataConfig{}
	dAtAs.regionMultiLevelNpcData.Map, dAtAs.regionMultiLevelNpcData.parserMap, err = regdata.LoadRegionMultiLevelNpcData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.regionMultiLevelNpcData.Map))
	for k := range dAtAs.regionMultiLevelNpcData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.regionMultiLevelNpcData.Array = make([]*regdata.RegionMultiLevelNpcData, 0, len(dAtAs.regionMultiLevelNpcData.Map))
	for _, k := range uint64Keys {
		dAtAs.regionMultiLevelNpcData.Array = append(dAtAs.regionMultiLevelNpcData.Array, dAtAs.regionMultiLevelNpcData.Map[k])
	}
	dAtAs.regionMultiLevelNpcData.MinKeyData = dAtAs.regionMultiLevelNpcData.Array[0]
	dAtAs.regionMultiLevelNpcData.MaxKeyData = dAtAs.regionMultiLevelNpcData.Array[len(dAtAs.regionMultiLevelNpcData.Array)-1]

	// regdata.RegionMultiLevelNpcLevelData
	dAtAs.regionMultiLevelNpcLevelData = &RegionMultiLevelNpcLevelDataConfig{}
	dAtAs.regionMultiLevelNpcLevelData.Map, dAtAs.regionMultiLevelNpcLevelData.parserMap, err = regdata.LoadRegionMultiLevelNpcLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.regionMultiLevelNpcLevelData.Map))
	for k := range dAtAs.regionMultiLevelNpcLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.regionMultiLevelNpcLevelData.Array = make([]*regdata.RegionMultiLevelNpcLevelData, 0, len(dAtAs.regionMultiLevelNpcLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.regionMultiLevelNpcLevelData.Array = append(dAtAs.regionMultiLevelNpcLevelData.Array, dAtAs.regionMultiLevelNpcLevelData.Map[k])
	}
	dAtAs.regionMultiLevelNpcLevelData.MinKeyData = dAtAs.regionMultiLevelNpcLevelData.Array[0]
	dAtAs.regionMultiLevelNpcLevelData.MaxKeyData = dAtAs.regionMultiLevelNpcLevelData.Array[len(dAtAs.regionMultiLevelNpcLevelData.Array)-1]

	// regdata.RegionMultiLevelNpcTypeData
	dAtAs.regionMultiLevelNpcTypeData = &RegionMultiLevelNpcTypeDataConfig{}
	dAtAs.regionMultiLevelNpcTypeData.Map, dAtAs.regionMultiLevelNpcTypeData.parserMap, err = regdata.LoadRegionMultiLevelNpcTypeData(gos)
	if err != nil {
		return nil, err
	}

	intKeys = make([]int, 0, len(dAtAs.regionMultiLevelNpcTypeData.Map))
	for k := range dAtAs.regionMultiLevelNpcTypeData.Map {
		intKeys = append(intKeys, k)
	}
	sort.Sort(intSlice(intKeys))
	dAtAs.regionMultiLevelNpcTypeData.Array = make([]*regdata.RegionMultiLevelNpcTypeData, 0, len(dAtAs.regionMultiLevelNpcTypeData.Map))
	for _, k := range intKeys {
		dAtAs.regionMultiLevelNpcTypeData.Array = append(dAtAs.regionMultiLevelNpcTypeData.Array, dAtAs.regionMultiLevelNpcTypeData.Map[k])
	}
	dAtAs.regionMultiLevelNpcTypeData.MinKeyData = dAtAs.regionMultiLevelNpcTypeData.Array[0]
	dAtAs.regionMultiLevelNpcTypeData.MaxKeyData = dAtAs.regionMultiLevelNpcTypeData.Array[len(dAtAs.regionMultiLevelNpcTypeData.Array)-1]

	// regdata.TroopDialogueData
	dAtAs.troopDialogueData = &TroopDialogueDataConfig{}
	dAtAs.troopDialogueData.Map, dAtAs.troopDialogueData.parserMap, err = regdata.LoadTroopDialogueData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.troopDialogueData.Map))
	for k := range dAtAs.troopDialogueData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.troopDialogueData.Array = make([]*regdata.TroopDialogueData, 0, len(dAtAs.troopDialogueData.Map))
	for _, k := range uint64Keys {
		dAtAs.troopDialogueData.Array = append(dAtAs.troopDialogueData.Array, dAtAs.troopDialogueData.Map[k])
	}
	dAtAs.troopDialogueData.MinKeyData = dAtAs.troopDialogueData.Array[0]
	dAtAs.troopDialogueData.MaxKeyData = dAtAs.troopDialogueData.Array[len(dAtAs.troopDialogueData.Array)-1]

	// regdata.TroopDialogueTextData
	dAtAs.troopDialogueTextData = &TroopDialogueTextDataConfig{}
	dAtAs.troopDialogueTextData.Map, dAtAs.troopDialogueTextData.parserMap, err = regdata.LoadTroopDialogueTextData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.troopDialogueTextData.Map))
	for k := range dAtAs.troopDialogueTextData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.troopDialogueTextData.Array = make([]*regdata.TroopDialogueTextData, 0, len(dAtAs.troopDialogueTextData.Map))
	for _, k := range uint64Keys {
		dAtAs.troopDialogueTextData.Array = append(dAtAs.troopDialogueTextData.Array, dAtAs.troopDialogueTextData.Map[k])
	}
	dAtAs.troopDialogueTextData.MinKeyData = dAtAs.troopDialogueTextData.Array[0]
	dAtAs.troopDialogueTextData.MaxKeyData = dAtAs.troopDialogueTextData.Array[len(dAtAs.troopDialogueTextData.Array)-1]

	// resdata.AmountShowSortData
	dAtAs.amountShowSortData = &AmountShowSortDataConfig{}
	dAtAs.amountShowSortData.Map, dAtAs.amountShowSortData.parserMap, err = resdata.LoadAmountShowSortData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.amountShowSortData.Map))
	for k := range dAtAs.amountShowSortData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.amountShowSortData.Array = make([]*resdata.AmountShowSortData, 0, len(dAtAs.amountShowSortData.Map))
	for _, k := range uint64Keys {
		dAtAs.amountShowSortData.Array = append(dAtAs.amountShowSortData.Array, dAtAs.amountShowSortData.Map[k])
	}
	dAtAs.amountShowSortData.MinKeyData = dAtAs.amountShowSortData.Array[0]
	dAtAs.amountShowSortData.MaxKeyData = dAtAs.amountShowSortData.Array[len(dAtAs.amountShowSortData.Array)-1]

	// resdata.BaowuData
	dAtAs.baowuData = &BaowuDataConfig{}
	dAtAs.baowuData.Map, dAtAs.baowuData.parserMap, err = resdata.LoadBaowuData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.baowuData.Map))
	for k := range dAtAs.baowuData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.baowuData.Array = make([]*resdata.BaowuData, 0, len(dAtAs.baowuData.Map))
	for _, k := range uint64Keys {
		dAtAs.baowuData.Array = append(dAtAs.baowuData.Array, dAtAs.baowuData.Map[k])
	}
	dAtAs.baowuData.MinKeyData = dAtAs.baowuData.Array[0]
	dAtAs.baowuData.MaxKeyData = dAtAs.baowuData.Array[len(dAtAs.baowuData.Array)-1]

	// resdata.ConditionPlunder
	dAtAs.conditionPlunder = &ConditionPlunderConfig{}
	dAtAs.conditionPlunder.Map, dAtAs.conditionPlunder.parserMap, err = resdata.LoadConditionPlunder(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.conditionPlunder.Map))
	for k := range dAtAs.conditionPlunder.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.conditionPlunder.Array = make([]*resdata.ConditionPlunder, 0, len(dAtAs.conditionPlunder.Map))
	for _, k := range uint64Keys {
		dAtAs.conditionPlunder.Array = append(dAtAs.conditionPlunder.Array, dAtAs.conditionPlunder.Map[k])
	}
	dAtAs.conditionPlunder.MinKeyData = dAtAs.conditionPlunder.Array[0]
	dAtAs.conditionPlunder.MaxKeyData = dAtAs.conditionPlunder.Array[len(dAtAs.conditionPlunder.Array)-1]

	// resdata.ConditionPlunderItem
	dAtAs.conditionPlunderItem = &ConditionPlunderItemConfig{}
	dAtAs.conditionPlunderItem.Map, dAtAs.conditionPlunderItem.parserMap, err = resdata.LoadConditionPlunderItem(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.conditionPlunderItem.Map))
	for k := range dAtAs.conditionPlunderItem.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.conditionPlunderItem.Array = make([]*resdata.ConditionPlunderItem, 0, len(dAtAs.conditionPlunderItem.Map))
	for _, k := range uint64Keys {
		dAtAs.conditionPlunderItem.Array = append(dAtAs.conditionPlunderItem.Array, dAtAs.conditionPlunderItem.Map[k])
	}
	dAtAs.conditionPlunderItem.MinKeyData = dAtAs.conditionPlunderItem.Array[0]
	dAtAs.conditionPlunderItem.MaxKeyData = dAtAs.conditionPlunderItem.Array[len(dAtAs.conditionPlunderItem.Array)-1]

	// resdata.Cost
	dAtAs.cost = &CostConfig{}
	dAtAs.cost.Map, dAtAs.cost.parserMap, err = resdata.LoadCost(gos)
	if err != nil {
		return nil, err
	}

	intKeys = make([]int, 0, len(dAtAs.cost.Map))
	for k := range dAtAs.cost.Map {
		intKeys = append(intKeys, k)
	}
	sort.Sort(intSlice(intKeys))
	dAtAs.cost.Array = make([]*resdata.Cost, 0, len(dAtAs.cost.Map))
	for _, k := range intKeys {
		dAtAs.cost.Array = append(dAtAs.cost.Array, dAtAs.cost.Map[k])
	}
	dAtAs.cost.MinKeyData = dAtAs.cost.Array[0]
	dAtAs.cost.MaxKeyData = dAtAs.cost.Array[len(dAtAs.cost.Array)-1]

	// resdata.GuildLevelPrize
	dAtAs.guildLevelPrize = &GuildLevelPrizeConfig{}
	dAtAs.guildLevelPrize.Map, dAtAs.guildLevelPrize.parserMap, err = resdata.LoadGuildLevelPrize(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.guildLevelPrize.Map))
	for k := range dAtAs.guildLevelPrize.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.guildLevelPrize.Array = make([]*resdata.GuildLevelPrize, 0, len(dAtAs.guildLevelPrize.Map))
	for _, k := range uint64Keys {
		dAtAs.guildLevelPrize.Array = append(dAtAs.guildLevelPrize.Array, dAtAs.guildLevelPrize.Map[k])
	}
	dAtAs.guildLevelPrize.MinKeyData = dAtAs.guildLevelPrize.Array[0]
	dAtAs.guildLevelPrize.MaxKeyData = dAtAs.guildLevelPrize.Array[len(dAtAs.guildLevelPrize.Array)-1]

	// resdata.Plunder
	dAtAs.plunder = &PlunderConfig{}
	dAtAs.plunder.Map, dAtAs.plunder.parserMap, err = resdata.LoadPlunder(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.plunder.Map))
	for k := range dAtAs.plunder.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.plunder.Array = make([]*resdata.Plunder, 0, len(dAtAs.plunder.Map))
	for _, k := range uint64Keys {
		dAtAs.plunder.Array = append(dAtAs.plunder.Array, dAtAs.plunder.Map[k])
	}
	dAtAs.plunder.MinKeyData = dAtAs.plunder.Array[0]
	dAtAs.plunder.MaxKeyData = dAtAs.plunder.Array[len(dAtAs.plunder.Array)-1]

	// resdata.PlunderGroup
	dAtAs.plunderGroup = &PlunderGroupConfig{}
	dAtAs.plunderGroup.Map, dAtAs.plunderGroup.parserMap, err = resdata.LoadPlunderGroup(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.plunderGroup.Map))
	for k := range dAtAs.plunderGroup.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.plunderGroup.Array = make([]*resdata.PlunderGroup, 0, len(dAtAs.plunderGroup.Map))
	for _, k := range uint64Keys {
		dAtAs.plunderGroup.Array = append(dAtAs.plunderGroup.Array, dAtAs.plunderGroup.Map[k])
	}
	dAtAs.plunderGroup.MinKeyData = dAtAs.plunderGroup.Array[0]
	dAtAs.plunderGroup.MaxKeyData = dAtAs.plunderGroup.Array[len(dAtAs.plunderGroup.Array)-1]

	// resdata.PlunderItem
	dAtAs.plunderItem = &PlunderItemConfig{}
	dAtAs.plunderItem.Map, dAtAs.plunderItem.parserMap, err = resdata.LoadPlunderItem(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.plunderItem.Map))
	for k := range dAtAs.plunderItem.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.plunderItem.Array = make([]*resdata.PlunderItem, 0, len(dAtAs.plunderItem.Map))
	for _, k := range uint64Keys {
		dAtAs.plunderItem.Array = append(dAtAs.plunderItem.Array, dAtAs.plunderItem.Map[k])
	}
	dAtAs.plunderItem.MinKeyData = dAtAs.plunderItem.Array[0]
	dAtAs.plunderItem.MaxKeyData = dAtAs.plunderItem.Array[len(dAtAs.plunderItem.Array)-1]

	// resdata.PlunderPrize
	dAtAs.plunderPrize = &PlunderPrizeConfig{}
	dAtAs.plunderPrize.Map, dAtAs.plunderPrize.parserMap, err = resdata.LoadPlunderPrize(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.plunderPrize.Map))
	for k := range dAtAs.plunderPrize.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.plunderPrize.Array = make([]*resdata.PlunderPrize, 0, len(dAtAs.plunderPrize.Map))
	for _, k := range uint64Keys {
		dAtAs.plunderPrize.Array = append(dAtAs.plunderPrize.Array, dAtAs.plunderPrize.Map[k])
	}
	dAtAs.plunderPrize.MinKeyData = dAtAs.plunderPrize.Array[0]
	dAtAs.plunderPrize.MaxKeyData = dAtAs.plunderPrize.Array[len(dAtAs.plunderPrize.Array)-1]

	// resdata.Prize
	dAtAs.prize = &PrizeConfig{}
	dAtAs.prize.Map, dAtAs.prize.parserMap, err = resdata.LoadPrize(gos)
	if err != nil {
		return nil, err
	}

	intKeys = make([]int, 0, len(dAtAs.prize.Map))
	for k := range dAtAs.prize.Map {
		intKeys = append(intKeys, k)
	}
	sort.Sort(intSlice(intKeys))
	dAtAs.prize.Array = make([]*resdata.Prize, 0, len(dAtAs.prize.Map))
	for _, k := range intKeys {
		dAtAs.prize.Array = append(dAtAs.prize.Array, dAtAs.prize.Map[k])
	}
	dAtAs.prize.MinKeyData = dAtAs.prize.Array[0]
	dAtAs.prize.MaxKeyData = dAtAs.prize.Array[len(dAtAs.prize.Array)-1]

	// resdata.ResCaptainData
	dAtAs.resCaptainData = &ResCaptainDataConfig{}
	dAtAs.resCaptainData.Map, dAtAs.resCaptainData.parserMap, err = resdata.LoadResCaptainData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.resCaptainData.Map))
	for k := range dAtAs.resCaptainData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.resCaptainData.Array = make([]*resdata.ResCaptainData, 0, len(dAtAs.resCaptainData.Map))
	for _, k := range uint64Keys {
		dAtAs.resCaptainData.Array = append(dAtAs.resCaptainData.Array, dAtAs.resCaptainData.Map[k])
	}
	dAtAs.resCaptainData.MinKeyData = dAtAs.resCaptainData.Array[0]
	dAtAs.resCaptainData.MaxKeyData = dAtAs.resCaptainData.Array[len(dAtAs.resCaptainData.Array)-1]

	// scene.CombatScene
	dAtAs.combatScene = &CombatSceneConfig{}
	dAtAs.combatScene.Map, dAtAs.combatScene.parserMap, err = scene.LoadCombatScene(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.combatScene.Map))
	for k := range dAtAs.combatScene.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.combatScene.Array = make([]*scene.CombatScene, 0, len(dAtAs.combatScene.Map))
	for _, k := range stringKeys {
		dAtAs.combatScene.Array = append(dAtAs.combatScene.Array, dAtAs.combatScene.Map[k])
	}
	dAtAs.combatScene.MinKeyData = dAtAs.combatScene.Array[0]
	dAtAs.combatScene.MaxKeyData = dAtAs.combatScene.Array[len(dAtAs.combatScene.Array)-1]

	// season.SeasonData
	dAtAs.seasonData = &SeasonDataConfig{}
	dAtAs.seasonData.Map, dAtAs.seasonData.parserMap, err = season.LoadSeasonData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.seasonData.Map))
	for k := range dAtAs.seasonData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.seasonData.Array = make([]*season.SeasonData, 0, len(dAtAs.seasonData.Map))
	for _, k := range uint64Keys {
		dAtAs.seasonData.Array = append(dAtAs.seasonData.Array, dAtAs.seasonData.Map[k])
	}
	dAtAs.seasonData.MinKeyData = dAtAs.seasonData.Array[0]
	dAtAs.seasonData.MaxKeyData = dAtAs.seasonData.Array[len(dAtAs.seasonData.Array)-1]

	// settings.PrivacySettingData
	dAtAs.privacySettingData = &PrivacySettingDataConfig{}
	dAtAs.privacySettingData.Map, dAtAs.privacySettingData.parserMap, err = settings.LoadPrivacySettingData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.privacySettingData.Map))
	for k := range dAtAs.privacySettingData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.privacySettingData.Array = make([]*settings.PrivacySettingData, 0, len(dAtAs.privacySettingData.Map))
	for _, k := range uint64Keys {
		dAtAs.privacySettingData.Array = append(dAtAs.privacySettingData.Array, dAtAs.privacySettingData.Map[k])
	}
	dAtAs.privacySettingData.MinKeyData = dAtAs.privacySettingData.Array[0]
	dAtAs.privacySettingData.MaxKeyData = dAtAs.privacySettingData.Array[len(dAtAs.privacySettingData.Array)-1]

	// shop.BlackMarketData
	dAtAs.blackMarketData = &BlackMarketDataConfig{}
	dAtAs.blackMarketData.Map, dAtAs.blackMarketData.parserMap, err = shop.LoadBlackMarketData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.blackMarketData.Map))
	for k := range dAtAs.blackMarketData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.blackMarketData.Array = make([]*shop.BlackMarketData, 0, len(dAtAs.blackMarketData.Map))
	for _, k := range uint64Keys {
		dAtAs.blackMarketData.Array = append(dAtAs.blackMarketData.Array, dAtAs.blackMarketData.Map[k])
	}
	dAtAs.blackMarketData.MinKeyData = dAtAs.blackMarketData.Array[0]
	dAtAs.blackMarketData.MaxKeyData = dAtAs.blackMarketData.Array[len(dAtAs.blackMarketData.Array)-1]

	// shop.BlackMarketGoodsData
	dAtAs.blackMarketGoodsData = &BlackMarketGoodsDataConfig{}
	dAtAs.blackMarketGoodsData.Map, dAtAs.blackMarketGoodsData.parserMap, err = shop.LoadBlackMarketGoodsData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.blackMarketGoodsData.Map))
	for k := range dAtAs.blackMarketGoodsData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.blackMarketGoodsData.Array = make([]*shop.BlackMarketGoodsData, 0, len(dAtAs.blackMarketGoodsData.Map))
	for _, k := range uint64Keys {
		dAtAs.blackMarketGoodsData.Array = append(dAtAs.blackMarketGoodsData.Array, dAtAs.blackMarketGoodsData.Map[k])
	}
	dAtAs.blackMarketGoodsData.MinKeyData = dAtAs.blackMarketGoodsData.Array[0]
	dAtAs.blackMarketGoodsData.MaxKeyData = dAtAs.blackMarketGoodsData.Array[len(dAtAs.blackMarketGoodsData.Array)-1]

	// shop.BlackMarketGoodsGroupData
	dAtAs.blackMarketGoodsGroupData = &BlackMarketGoodsGroupDataConfig{}
	dAtAs.blackMarketGoodsGroupData.Map, dAtAs.blackMarketGoodsGroupData.parserMap, err = shop.LoadBlackMarketGoodsGroupData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.blackMarketGoodsGroupData.Map))
	for k := range dAtAs.blackMarketGoodsGroupData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.blackMarketGoodsGroupData.Array = make([]*shop.BlackMarketGoodsGroupData, 0, len(dAtAs.blackMarketGoodsGroupData.Map))
	for _, k := range uint64Keys {
		dAtAs.blackMarketGoodsGroupData.Array = append(dAtAs.blackMarketGoodsGroupData.Array, dAtAs.blackMarketGoodsGroupData.Map[k])
	}
	dAtAs.blackMarketGoodsGroupData.MinKeyData = dAtAs.blackMarketGoodsGroupData.Array[0]
	dAtAs.blackMarketGoodsGroupData.MaxKeyData = dAtAs.blackMarketGoodsGroupData.Array[len(dAtAs.blackMarketGoodsGroupData.Array)-1]

	// shop.DiscountColorData
	dAtAs.discountColorData = &DiscountColorDataConfig{}
	dAtAs.discountColorData.Map, dAtAs.discountColorData.parserMap, err = shop.LoadDiscountColorData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.discountColorData.Map))
	for k := range dAtAs.discountColorData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.discountColorData.Array = make([]*shop.DiscountColorData, 0, len(dAtAs.discountColorData.Map))
	for _, k := range uint64Keys {
		dAtAs.discountColorData.Array = append(dAtAs.discountColorData.Array, dAtAs.discountColorData.Map[k])
	}
	dAtAs.discountColorData.MinKeyData = dAtAs.discountColorData.Array[0]
	dAtAs.discountColorData.MaxKeyData = dAtAs.discountColorData.Array[len(dAtAs.discountColorData.Array)-1]

	// shop.NormalShopGoods
	dAtAs.normalShopGoods = &NormalShopGoodsConfig{}
	dAtAs.normalShopGoods.Map, dAtAs.normalShopGoods.parserMap, err = shop.LoadNormalShopGoods(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.normalShopGoods.Map))
	for k := range dAtAs.normalShopGoods.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.normalShopGoods.Array = make([]*shop.NormalShopGoods, 0, len(dAtAs.normalShopGoods.Map))
	for _, k := range uint64Keys {
		dAtAs.normalShopGoods.Array = append(dAtAs.normalShopGoods.Array, dAtAs.normalShopGoods.Map[k])
	}
	dAtAs.normalShopGoods.MinKeyData = dAtAs.normalShopGoods.Array[0]
	dAtAs.normalShopGoods.MaxKeyData = dAtAs.normalShopGoods.Array[len(dAtAs.normalShopGoods.Array)-1]

	// shop.Shop
	dAtAs.shop = &ShopConfig{}
	dAtAs.shop.Map, dAtAs.shop.parserMap, err = shop.LoadShop(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.shop.Map))
	for k := range dAtAs.shop.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.shop.Array = make([]*shop.Shop, 0, len(dAtAs.shop.Map))
	for _, k := range uint64Keys {
		dAtAs.shop.Array = append(dAtAs.shop.Array, dAtAs.shop.Map[k])
	}
	dAtAs.shop.MinKeyData = dAtAs.shop.Array[0]
	dAtAs.shop.MaxKeyData = dAtAs.shop.Array[len(dAtAs.shop.Array)-1]

	// shop.ZhenBaoGeShopGoods
	dAtAs.zhenBaoGeShopGoods = &ZhenBaoGeShopGoodsConfig{}
	dAtAs.zhenBaoGeShopGoods.Map, dAtAs.zhenBaoGeShopGoods.parserMap, err = shop.LoadZhenBaoGeShopGoods(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.zhenBaoGeShopGoods.Map))
	for k := range dAtAs.zhenBaoGeShopGoods.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.zhenBaoGeShopGoods.Array = make([]*shop.ZhenBaoGeShopGoods, 0, len(dAtAs.zhenBaoGeShopGoods.Map))
	for _, k := range uint64Keys {
		dAtAs.zhenBaoGeShopGoods.Array = append(dAtAs.zhenBaoGeShopGoods.Array, dAtAs.zhenBaoGeShopGoods.Map[k])
	}
	dAtAs.zhenBaoGeShopGoods.MinKeyData = dAtAs.zhenBaoGeShopGoods.Array[0]
	dAtAs.zhenBaoGeShopGoods.MaxKeyData = dAtAs.zhenBaoGeShopGoods.Array[len(dAtAs.zhenBaoGeShopGoods.Array)-1]

	// spell.PassiveSpellData
	dAtAs.passiveSpellData = &PassiveSpellDataConfig{}
	dAtAs.passiveSpellData.Map, dAtAs.passiveSpellData.parserMap, err = spell.LoadPassiveSpellData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.passiveSpellData.Map))
	for k := range dAtAs.passiveSpellData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.passiveSpellData.Array = make([]*spell.PassiveSpellData, 0, len(dAtAs.passiveSpellData.Map))
	for _, k := range uint64Keys {
		dAtAs.passiveSpellData.Array = append(dAtAs.passiveSpellData.Array, dAtAs.passiveSpellData.Map[k])
	}
	dAtAs.passiveSpellData.MinKeyData = dAtAs.passiveSpellData.Array[0]
	dAtAs.passiveSpellData.MaxKeyData = dAtAs.passiveSpellData.Array[len(dAtAs.passiveSpellData.Array)-1]

	// spell.Spell
	dAtAs.spell = &SpellConfig{}
	dAtAs.spell.Map, dAtAs.spell.parserMap, err = spell.LoadSpell(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.spell.Map))
	for k := range dAtAs.spell.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.spell.Array = make([]*spell.Spell, 0, len(dAtAs.spell.Map))
	for _, k := range uint64Keys {
		dAtAs.spell.Array = append(dAtAs.spell.Array, dAtAs.spell.Map[k])
	}
	dAtAs.spell.MinKeyData = dAtAs.spell.Array[0]
	dAtAs.spell.MaxKeyData = dAtAs.spell.Array[len(dAtAs.spell.Array)-1]

	// spell.SpellData
	dAtAs.spellData = &SpellDataConfig{}
	dAtAs.spellData.Map, dAtAs.spellData.parserMap, err = spell.LoadSpellData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.spellData.Map))
	for k := range dAtAs.spellData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.spellData.Array = make([]*spell.SpellData, 0, len(dAtAs.spellData.Map))
	for _, k := range uint64Keys {
		dAtAs.spellData.Array = append(dAtAs.spellData.Array, dAtAs.spellData.Map[k])
	}
	dAtAs.spellData.MinKeyData = dAtAs.spellData.Array[0]
	dAtAs.spellData.MaxKeyData = dAtAs.spellData.Array[len(dAtAs.spellData.Array)-1]

	// spell.SpellFacadeData
	dAtAs.spellFacadeData = &SpellFacadeDataConfig{}
	dAtAs.spellFacadeData.Map, dAtAs.spellFacadeData.parserMap, err = spell.LoadSpellFacadeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.spellFacadeData.Map))
	for k := range dAtAs.spellFacadeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.spellFacadeData.Array = make([]*spell.SpellFacadeData, 0, len(dAtAs.spellFacadeData.Map))
	for _, k := range uint64Keys {
		dAtAs.spellFacadeData.Array = append(dAtAs.spellFacadeData.Array, dAtAs.spellFacadeData.Map[k])
	}
	dAtAs.spellFacadeData.MinKeyData = dAtAs.spellFacadeData.Array[0]
	dAtAs.spellFacadeData.MaxKeyData = dAtAs.spellFacadeData.Array[len(dAtAs.spellFacadeData.Array)-1]

	// spell.StateData
	dAtAs.stateData = &StateDataConfig{}
	dAtAs.stateData.Map, dAtAs.stateData.parserMap, err = spell.LoadStateData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.stateData.Map))
	for k := range dAtAs.stateData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.stateData.Array = make([]*spell.StateData, 0, len(dAtAs.stateData.Map))
	for _, k := range uint64Keys {
		dAtAs.stateData.Array = append(dAtAs.stateData.Array, dAtAs.stateData.Map[k])
	}
	dAtAs.stateData.MinKeyData = dAtAs.stateData.Array[0]
	dAtAs.stateData.MaxKeyData = dAtAs.stateData.Array[len(dAtAs.stateData.Array)-1]

	// strategydata.StrategyData
	dAtAs.strategyData = &StrategyDataConfig{}
	dAtAs.strategyData.Map, dAtAs.strategyData.parserMap, err = strategydata.LoadStrategyData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.strategyData.Map))
	for k := range dAtAs.strategyData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.strategyData.Array = make([]*strategydata.StrategyData, 0, len(dAtAs.strategyData.Map))
	for _, k := range uint64Keys {
		dAtAs.strategyData.Array = append(dAtAs.strategyData.Array, dAtAs.strategyData.Map[k])
	}
	dAtAs.strategyData.MinKeyData = dAtAs.strategyData.Array[0]
	dAtAs.strategyData.MaxKeyData = dAtAs.strategyData.Array[len(dAtAs.strategyData.Array)-1]

	// strategydata.StrategyEffectData
	dAtAs.strategyEffectData = &StrategyEffectDataConfig{}
	dAtAs.strategyEffectData.Map, dAtAs.strategyEffectData.parserMap, err = strategydata.LoadStrategyEffectData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.strategyEffectData.Map))
	for k := range dAtAs.strategyEffectData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.strategyEffectData.Array = make([]*strategydata.StrategyEffectData, 0, len(dAtAs.strategyEffectData.Map))
	for _, k := range uint64Keys {
		dAtAs.strategyEffectData.Array = append(dAtAs.strategyEffectData.Array, dAtAs.strategyEffectData.Map[k])
	}
	dAtAs.strategyEffectData.MinKeyData = dAtAs.strategyEffectData.Array[0]
	dAtAs.strategyEffectData.MaxKeyData = dAtAs.strategyEffectData.Array[len(dAtAs.strategyEffectData.Array)-1]

	// strongerdata.StrongerData
	dAtAs.strongerData = &StrongerDataConfig{}
	dAtAs.strongerData.Map, dAtAs.strongerData.parserMap, err = strongerdata.LoadStrongerData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.strongerData.Map))
	for k := range dAtAs.strongerData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.strongerData.Array = make([]*strongerdata.StrongerData, 0, len(dAtAs.strongerData.Map))
	for _, k := range uint64Keys {
		dAtAs.strongerData.Array = append(dAtAs.strongerData.Array, dAtAs.strongerData.Map[k])
	}
	dAtAs.strongerData.MinKeyData = dAtAs.strongerData.Array[0]
	dAtAs.strongerData.MaxKeyData = dAtAs.strongerData.Array[len(dAtAs.strongerData.Array)-1]

	// sub.BuildingEffectData
	dAtAs.buildingEffectData = &BuildingEffectDataConfig{}
	dAtAs.buildingEffectData.Map, dAtAs.buildingEffectData.parserMap, err = sub.LoadBuildingEffectData(gos)
	if err != nil {
		return nil, err
	}

	intKeys = make([]int, 0, len(dAtAs.buildingEffectData.Map))
	for k := range dAtAs.buildingEffectData.Map {
		intKeys = append(intKeys, k)
	}
	sort.Sort(intSlice(intKeys))
	dAtAs.buildingEffectData.Array = make([]*sub.BuildingEffectData, 0, len(dAtAs.buildingEffectData.Map))
	for _, k := range intKeys {
		dAtAs.buildingEffectData.Array = append(dAtAs.buildingEffectData.Array, dAtAs.buildingEffectData.Map[k])
	}
	dAtAs.buildingEffectData.MinKeyData = dAtAs.buildingEffectData.Array[0]
	dAtAs.buildingEffectData.MaxKeyData = dAtAs.buildingEffectData.Array[len(dAtAs.buildingEffectData.Array)-1]

	// survey.SurveyData
	dAtAs.surveyData = &SurveyDataConfig{}
	dAtAs.surveyData.Map, dAtAs.surveyData.parserMap, err = survey.LoadSurveyData(gos)
	if err != nil {
		return nil, err
	}

	stringKeys = make([]string, 0, len(dAtAs.surveyData.Map))
	for k := range dAtAs.surveyData.Map {
		stringKeys = append(stringKeys, k)
	}
	sort.Sort(stringSlice(stringKeys))
	dAtAs.surveyData.Array = make([]*survey.SurveyData, 0, len(dAtAs.surveyData.Map))
	for _, k := range stringKeys {
		dAtAs.surveyData.Array = append(dAtAs.surveyData.Array, dAtAs.surveyData.Map[k])
	}
	dAtAs.surveyData.MinKeyData = dAtAs.surveyData.Array[0]
	dAtAs.surveyData.MaxKeyData = dAtAs.surveyData.Array[len(dAtAs.surveyData.Array)-1]

	// taskdata.AchieveTaskData
	dAtAs.achieveTaskData = &AchieveTaskDataConfig{}
	dAtAs.achieveTaskData.Map, dAtAs.achieveTaskData.parserMap, err = taskdata.LoadAchieveTaskData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.achieveTaskData.Map))
	for k := range dAtAs.achieveTaskData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.achieveTaskData.Array = make([]*taskdata.AchieveTaskData, 0, len(dAtAs.achieveTaskData.Map))
	for _, k := range uint64Keys {
		dAtAs.achieveTaskData.Array = append(dAtAs.achieveTaskData.Array, dAtAs.achieveTaskData.Map[k])
	}
	dAtAs.achieveTaskData.MinKeyData = dAtAs.achieveTaskData.Array[0]
	dAtAs.achieveTaskData.MaxKeyData = dAtAs.achieveTaskData.Array[len(dAtAs.achieveTaskData.Array)-1]

	// taskdata.AchieveTaskStarPrizeData
	dAtAs.achieveTaskStarPrizeData = &AchieveTaskStarPrizeDataConfig{}
	dAtAs.achieveTaskStarPrizeData.Map, dAtAs.achieveTaskStarPrizeData.parserMap, err = taskdata.LoadAchieveTaskStarPrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.achieveTaskStarPrizeData.Map))
	for k := range dAtAs.achieveTaskStarPrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.achieveTaskStarPrizeData.Array = make([]*taskdata.AchieveTaskStarPrizeData, 0, len(dAtAs.achieveTaskStarPrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.achieveTaskStarPrizeData.Array = append(dAtAs.achieveTaskStarPrizeData.Array, dAtAs.achieveTaskStarPrizeData.Map[k])
	}
	dAtAs.achieveTaskStarPrizeData.MinKeyData = dAtAs.achieveTaskStarPrizeData.Array[0]
	dAtAs.achieveTaskStarPrizeData.MaxKeyData = dAtAs.achieveTaskStarPrizeData.Array[len(dAtAs.achieveTaskStarPrizeData.Array)-1]

	// taskdata.ActiveDegreePrizeData
	dAtAs.activeDegreePrizeData = &ActiveDegreePrizeDataConfig{}
	dAtAs.activeDegreePrizeData.Map, dAtAs.activeDegreePrizeData.parserMap, err = taskdata.LoadActiveDegreePrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.activeDegreePrizeData.Map))
	for k := range dAtAs.activeDegreePrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.activeDegreePrizeData.Array = make([]*taskdata.ActiveDegreePrizeData, 0, len(dAtAs.activeDegreePrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.activeDegreePrizeData.Array = append(dAtAs.activeDegreePrizeData.Array, dAtAs.activeDegreePrizeData.Map[k])
	}
	dAtAs.activeDegreePrizeData.MinKeyData = dAtAs.activeDegreePrizeData.Array[0]
	dAtAs.activeDegreePrizeData.MaxKeyData = dAtAs.activeDegreePrizeData.Array[len(dAtAs.activeDegreePrizeData.Array)-1]

	// taskdata.ActiveDegreeTaskData
	dAtAs.activeDegreeTaskData = &ActiveDegreeTaskDataConfig{}
	dAtAs.activeDegreeTaskData.Map, dAtAs.activeDegreeTaskData.parserMap, err = taskdata.LoadActiveDegreeTaskData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.activeDegreeTaskData.Map))
	for k := range dAtAs.activeDegreeTaskData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.activeDegreeTaskData.Array = make([]*taskdata.ActiveDegreeTaskData, 0, len(dAtAs.activeDegreeTaskData.Map))
	for _, k := range uint64Keys {
		dAtAs.activeDegreeTaskData.Array = append(dAtAs.activeDegreeTaskData.Array, dAtAs.activeDegreeTaskData.Map[k])
	}
	dAtAs.activeDegreeTaskData.MinKeyData = dAtAs.activeDegreeTaskData.Array[0]
	dAtAs.activeDegreeTaskData.MaxKeyData = dAtAs.activeDegreeTaskData.Array[len(dAtAs.activeDegreeTaskData.Array)-1]

	// taskdata.ActivityTaskData
	dAtAs.activityTaskData = &ActivityTaskDataConfig{}
	dAtAs.activityTaskData.Map, dAtAs.activityTaskData.parserMap, err = taskdata.LoadActivityTaskData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.activityTaskData.Map))
	for k := range dAtAs.activityTaskData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.activityTaskData.Array = make([]*taskdata.ActivityTaskData, 0, len(dAtAs.activityTaskData.Map))
	for _, k := range uint64Keys {
		dAtAs.activityTaskData.Array = append(dAtAs.activityTaskData.Array, dAtAs.activityTaskData.Map[k])
	}
	dAtAs.activityTaskData.MinKeyData = dAtAs.activityTaskData.Array[0]
	dAtAs.activityTaskData.MaxKeyData = dAtAs.activityTaskData.Array[len(dAtAs.activityTaskData.Array)-1]

	// taskdata.BaYeStageData
	dAtAs.baYeStageData = &BaYeStageDataConfig{}
	dAtAs.baYeStageData.Map, dAtAs.baYeStageData.parserMap, err = taskdata.LoadBaYeStageData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.baYeStageData.Map))
	for k := range dAtAs.baYeStageData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.baYeStageData.Array = make([]*taskdata.BaYeStageData, 0, len(dAtAs.baYeStageData.Map))
	for _, k := range uint64Keys {
		dAtAs.baYeStageData.Array = append(dAtAs.baYeStageData.Array, dAtAs.baYeStageData.Map[k])
	}
	dAtAs.baYeStageData.MinKeyData = dAtAs.baYeStageData.Array[0]
	dAtAs.baYeStageData.MaxKeyData = dAtAs.baYeStageData.Array[len(dAtAs.baYeStageData.Array)-1]

	// taskdata.BaYeTaskData
	dAtAs.baYeTaskData = &BaYeTaskDataConfig{}
	dAtAs.baYeTaskData.Map, dAtAs.baYeTaskData.parserMap, err = taskdata.LoadBaYeTaskData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.baYeTaskData.Map))
	for k := range dAtAs.baYeTaskData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.baYeTaskData.Array = make([]*taskdata.BaYeTaskData, 0, len(dAtAs.baYeTaskData.Map))
	for _, k := range uint64Keys {
		dAtAs.baYeTaskData.Array = append(dAtAs.baYeTaskData.Array, dAtAs.baYeTaskData.Map[k])
	}
	dAtAs.baYeTaskData.MinKeyData = dAtAs.baYeTaskData.Array[0]
	dAtAs.baYeTaskData.MaxKeyData = dAtAs.baYeTaskData.Array[len(dAtAs.baYeTaskData.Array)-1]

	// taskdata.BranchTaskData
	dAtAs.branchTaskData = &BranchTaskDataConfig{}
	dAtAs.branchTaskData.Map, dAtAs.branchTaskData.parserMap, err = taskdata.LoadBranchTaskData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.branchTaskData.Map))
	for k := range dAtAs.branchTaskData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.branchTaskData.Array = make([]*taskdata.BranchTaskData, 0, len(dAtAs.branchTaskData.Map))
	for _, k := range uint64Keys {
		dAtAs.branchTaskData.Array = append(dAtAs.branchTaskData.Array, dAtAs.branchTaskData.Map[k])
	}
	dAtAs.branchTaskData.MinKeyData = dAtAs.branchTaskData.Array[0]
	dAtAs.branchTaskData.MaxKeyData = dAtAs.branchTaskData.Array[len(dAtAs.branchTaskData.Array)-1]

	// taskdata.BwzlPrizeData
	dAtAs.bwzlPrizeData = &BwzlPrizeDataConfig{}
	dAtAs.bwzlPrizeData.Map, dAtAs.bwzlPrizeData.parserMap, err = taskdata.LoadBwzlPrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.bwzlPrizeData.Map))
	for k := range dAtAs.bwzlPrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.bwzlPrizeData.Array = make([]*taskdata.BwzlPrizeData, 0, len(dAtAs.bwzlPrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.bwzlPrizeData.Array = append(dAtAs.bwzlPrizeData.Array, dAtAs.bwzlPrizeData.Map[k])
	}
	dAtAs.bwzlPrizeData.MinKeyData = dAtAs.bwzlPrizeData.Array[0]
	dAtAs.bwzlPrizeData.MaxKeyData = dAtAs.bwzlPrizeData.Array[len(dAtAs.bwzlPrizeData.Array)-1]

	// taskdata.BwzlTaskData
	dAtAs.bwzlTaskData = &BwzlTaskDataConfig{}
	dAtAs.bwzlTaskData.Map, dAtAs.bwzlTaskData.parserMap, err = taskdata.LoadBwzlTaskData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.bwzlTaskData.Map))
	for k := range dAtAs.bwzlTaskData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.bwzlTaskData.Array = make([]*taskdata.BwzlTaskData, 0, len(dAtAs.bwzlTaskData.Map))
	for _, k := range uint64Keys {
		dAtAs.bwzlTaskData.Array = append(dAtAs.bwzlTaskData.Array, dAtAs.bwzlTaskData.Map[k])
	}
	dAtAs.bwzlTaskData.MinKeyData = dAtAs.bwzlTaskData.Array[0]
	dAtAs.bwzlTaskData.MaxKeyData = dAtAs.bwzlTaskData.Array[len(dAtAs.bwzlTaskData.Array)-1]

	// taskdata.MainTaskData
	dAtAs.mainTaskData = &MainTaskDataConfig{}
	dAtAs.mainTaskData.Map, dAtAs.mainTaskData.parserMap, err = taskdata.LoadMainTaskData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.mainTaskData.Map))
	for k := range dAtAs.mainTaskData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.mainTaskData.Array = make([]*taskdata.MainTaskData, 0, len(dAtAs.mainTaskData.Map))
	for _, k := range uint64Keys {
		dAtAs.mainTaskData.Array = append(dAtAs.mainTaskData.Array, dAtAs.mainTaskData.Map[k])
	}
	dAtAs.mainTaskData.MinKeyData = dAtAs.mainTaskData.Array[0]
	dAtAs.mainTaskData.MaxKeyData = dAtAs.mainTaskData.Array[len(dAtAs.mainTaskData.Array)-1]

	// taskdata.TaskBoxData
	dAtAs.taskBoxData = &TaskBoxDataConfig{}
	dAtAs.taskBoxData.Map, dAtAs.taskBoxData.parserMap, err = taskdata.LoadTaskBoxData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.taskBoxData.Map))
	for k := range dAtAs.taskBoxData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.taskBoxData.Array = make([]*taskdata.TaskBoxData, 0, len(dAtAs.taskBoxData.Map))
	for _, k := range uint64Keys {
		dAtAs.taskBoxData.Array = append(dAtAs.taskBoxData.Array, dAtAs.taskBoxData.Map[k])
	}
	dAtAs.taskBoxData.MinKeyData = dAtAs.taskBoxData.Array[0]
	dAtAs.taskBoxData.MaxKeyData = dAtAs.taskBoxData.Array[len(dAtAs.taskBoxData.Array)-1]

	// taskdata.TaskTargetData
	dAtAs.taskTargetData = &TaskTargetDataConfig{}
	dAtAs.taskTargetData.Map, dAtAs.taskTargetData.parserMap, err = taskdata.LoadTaskTargetData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.taskTargetData.Map))
	for k := range dAtAs.taskTargetData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.taskTargetData.Array = make([]*taskdata.TaskTargetData, 0, len(dAtAs.taskTargetData.Map))
	for _, k := range uint64Keys {
		dAtAs.taskTargetData.Array = append(dAtAs.taskTargetData.Array, dAtAs.taskTargetData.Map[k])
	}
	dAtAs.taskTargetData.MinKeyData = dAtAs.taskTargetData.Array[0]
	dAtAs.taskTargetData.MaxKeyData = dAtAs.taskTargetData.Array[len(dAtAs.taskTargetData.Array)-1]

	// taskdata.TitleData
	dAtAs.titleData = &TitleDataConfig{}
	dAtAs.titleData.Map, dAtAs.titleData.parserMap, err = taskdata.LoadTitleData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.titleData.Map))
	for k := range dAtAs.titleData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.titleData.Array = make([]*taskdata.TitleData, 0, len(dAtAs.titleData.Map))
	for _, k := range uint64Keys {
		dAtAs.titleData.Array = append(dAtAs.titleData.Array, dAtAs.titleData.Map[k])
	}
	dAtAs.titleData.MinKeyData = dAtAs.titleData.Array[0]
	dAtAs.titleData.MaxKeyData = dAtAs.titleData.Array[len(dAtAs.titleData.Array)-1]

	// taskdata.TitleTaskData
	dAtAs.titleTaskData = &TitleTaskDataConfig{}
	dAtAs.titleTaskData.Map, dAtAs.titleTaskData.parserMap, err = taskdata.LoadTitleTaskData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.titleTaskData.Map))
	for k := range dAtAs.titleTaskData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.titleTaskData.Array = make([]*taskdata.TitleTaskData, 0, len(dAtAs.titleTaskData.Map))
	for _, k := range uint64Keys {
		dAtAs.titleTaskData.Array = append(dAtAs.titleTaskData.Array, dAtAs.titleTaskData.Map[k])
	}
	dAtAs.titleTaskData.MinKeyData = dAtAs.titleTaskData.Array[0]
	dAtAs.titleTaskData.MaxKeyData = dAtAs.titleTaskData.Array[len(dAtAs.titleTaskData.Array)-1]

	// teach.TeachChapterData
	dAtAs.teachChapterData = &TeachChapterDataConfig{}
	dAtAs.teachChapterData.Map, dAtAs.teachChapterData.parserMap, err = teach.LoadTeachChapterData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.teachChapterData.Map))
	for k := range dAtAs.teachChapterData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.teachChapterData.Array = make([]*teach.TeachChapterData, 0, len(dAtAs.teachChapterData.Map))
	for _, k := range uint64Keys {
		dAtAs.teachChapterData.Array = append(dAtAs.teachChapterData.Array, dAtAs.teachChapterData.Map[k])
	}
	dAtAs.teachChapterData.MinKeyData = dAtAs.teachChapterData.Array[0]
	dAtAs.teachChapterData.MaxKeyData = dAtAs.teachChapterData.Array[len(dAtAs.teachChapterData.Array)-1]

	// towerdata.SecretTowerData
	dAtAs.secretTowerData = &SecretTowerDataConfig{}
	dAtAs.secretTowerData.Map, dAtAs.secretTowerData.parserMap, err = towerdata.LoadSecretTowerData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.secretTowerData.Map))
	for k := range dAtAs.secretTowerData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.secretTowerData.Array = make([]*towerdata.SecretTowerData, 0, len(dAtAs.secretTowerData.Map))
	for _, k := range uint64Keys {
		dAtAs.secretTowerData.Array = append(dAtAs.secretTowerData.Array, dAtAs.secretTowerData.Map[k])
	}
	dAtAs.secretTowerData.MinKeyData = dAtAs.secretTowerData.Array[0]
	dAtAs.secretTowerData.MaxKeyData = dAtAs.secretTowerData.Array[len(dAtAs.secretTowerData.Array)-1]

	// towerdata.SecretTowerWordsData
	dAtAs.secretTowerWordsData = &SecretTowerWordsDataConfig{}
	dAtAs.secretTowerWordsData.Map, dAtAs.secretTowerWordsData.parserMap, err = towerdata.LoadSecretTowerWordsData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.secretTowerWordsData.Map))
	for k := range dAtAs.secretTowerWordsData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.secretTowerWordsData.Array = make([]*towerdata.SecretTowerWordsData, 0, len(dAtAs.secretTowerWordsData.Map))
	for _, k := range uint64Keys {
		dAtAs.secretTowerWordsData.Array = append(dAtAs.secretTowerWordsData.Array, dAtAs.secretTowerWordsData.Map[k])
	}
	dAtAs.secretTowerWordsData.MinKeyData = dAtAs.secretTowerWordsData.Array[0]
	dAtAs.secretTowerWordsData.MaxKeyData = dAtAs.secretTowerWordsData.Array[len(dAtAs.secretTowerWordsData.Array)-1]

	// towerdata.TowerData
	dAtAs.towerData = &TowerDataConfig{}
	dAtAs.towerData.Map, dAtAs.towerData.parserMap, err = towerdata.LoadTowerData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.towerData.Map))
	for k := range dAtAs.towerData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.towerData.Array = make([]*towerdata.TowerData, 0, len(dAtAs.towerData.Map))
	for _, k := range uint64Keys {
		dAtAs.towerData.Array = append(dAtAs.towerData.Array, dAtAs.towerData.Map[k])
	}
	dAtAs.towerData.MinKeyData = dAtAs.towerData.Array[0]
	dAtAs.towerData.MaxKeyData = dAtAs.towerData.Array[len(dAtAs.towerData.Array)-1]

	// vip.VipContinueDaysData
	dAtAs.vipContinueDaysData = &VipContinueDaysDataConfig{}
	dAtAs.vipContinueDaysData.Map, dAtAs.vipContinueDaysData.parserMap, err = vip.LoadVipContinueDaysData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.vipContinueDaysData.Map))
	for k := range dAtAs.vipContinueDaysData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.vipContinueDaysData.Array = make([]*vip.VipContinueDaysData, 0, len(dAtAs.vipContinueDaysData.Map))
	for _, k := range uint64Keys {
		dAtAs.vipContinueDaysData.Array = append(dAtAs.vipContinueDaysData.Array, dAtAs.vipContinueDaysData.Map[k])
	}
	dAtAs.vipContinueDaysData.MinKeyData = dAtAs.vipContinueDaysData.Array[0]
	dAtAs.vipContinueDaysData.MaxKeyData = dAtAs.vipContinueDaysData.Array[len(dAtAs.vipContinueDaysData.Array)-1]

	// vip.VipLevelData
	dAtAs.vipLevelData = &VipLevelDataConfig{}
	dAtAs.vipLevelData.Map, dAtAs.vipLevelData.parserMap, err = vip.LoadVipLevelData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.vipLevelData.Map))
	for k := range dAtAs.vipLevelData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.vipLevelData.Array = make([]*vip.VipLevelData, 0, len(dAtAs.vipLevelData.Map))
	for _, k := range uint64Keys {
		dAtAs.vipLevelData.Array = append(dAtAs.vipLevelData.Array, dAtAs.vipLevelData.Map[k])
	}
	dAtAs.vipLevelData.MinKeyData = dAtAs.vipLevelData.Array[0]
	dAtAs.vipLevelData.MaxKeyData = dAtAs.vipLevelData.Array[len(dAtAs.vipLevelData.Array)-1]

	// xiongnu.ResistXiongNuData
	dAtAs.resistXiongNuData = &ResistXiongNuDataConfig{}
	dAtAs.resistXiongNuData.Map, dAtAs.resistXiongNuData.parserMap, err = xiongnu.LoadResistXiongNuData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.resistXiongNuData.Map))
	for k := range dAtAs.resistXiongNuData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.resistXiongNuData.Array = make([]*xiongnu.ResistXiongNuData, 0, len(dAtAs.resistXiongNuData.Map))
	for _, k := range uint64Keys {
		dAtAs.resistXiongNuData.Array = append(dAtAs.resistXiongNuData.Array, dAtAs.resistXiongNuData.Map[k])
	}
	dAtAs.resistXiongNuData.MinKeyData = dAtAs.resistXiongNuData.Array[0]
	dAtAs.resistXiongNuData.MaxKeyData = dAtAs.resistXiongNuData.Array[len(dAtAs.resistXiongNuData.Array)-1]

	// xiongnu.ResistXiongNuScoreData
	dAtAs.resistXiongNuScoreData = &ResistXiongNuScoreDataConfig{}
	dAtAs.resistXiongNuScoreData.Map, dAtAs.resistXiongNuScoreData.parserMap, err = xiongnu.LoadResistXiongNuScoreData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.resistXiongNuScoreData.Map))
	for k := range dAtAs.resistXiongNuScoreData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.resistXiongNuScoreData.Array = make([]*xiongnu.ResistXiongNuScoreData, 0, len(dAtAs.resistXiongNuScoreData.Map))
	for _, k := range uint64Keys {
		dAtAs.resistXiongNuScoreData.Array = append(dAtAs.resistXiongNuScoreData.Array, dAtAs.resistXiongNuScoreData.Map[k])
	}
	dAtAs.resistXiongNuScoreData.MinKeyData = dAtAs.resistXiongNuScoreData.Array[0]
	dAtAs.resistXiongNuScoreData.MaxKeyData = dAtAs.resistXiongNuScoreData.Array[len(dAtAs.resistXiongNuScoreData.Array)-1]

	// xiongnu.ResistXiongNuWaveData
	dAtAs.resistXiongNuWaveData = &ResistXiongNuWaveDataConfig{}
	dAtAs.resistXiongNuWaveData.Map, dAtAs.resistXiongNuWaveData.parserMap, err = xiongnu.LoadResistXiongNuWaveData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.resistXiongNuWaveData.Map))
	for k := range dAtAs.resistXiongNuWaveData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.resistXiongNuWaveData.Array = make([]*xiongnu.ResistXiongNuWaveData, 0, len(dAtAs.resistXiongNuWaveData.Map))
	for _, k := range uint64Keys {
		dAtAs.resistXiongNuWaveData.Array = append(dAtAs.resistXiongNuWaveData.Array, dAtAs.resistXiongNuWaveData.Map[k])
	}
	dAtAs.resistXiongNuWaveData.MinKeyData = dAtAs.resistXiongNuWaveData.Array[0]
	dAtAs.resistXiongNuWaveData.MaxKeyData = dAtAs.resistXiongNuWaveData.Array[len(dAtAs.resistXiongNuWaveData.Array)-1]

	// xuanydata.XuanyuanRangeData
	dAtAs.xuanyuanRangeData = &XuanyuanRangeDataConfig{}
	dAtAs.xuanyuanRangeData.Map, dAtAs.xuanyuanRangeData.parserMap, err = xuanydata.LoadXuanyuanRangeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.xuanyuanRangeData.Map))
	for k := range dAtAs.xuanyuanRangeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.xuanyuanRangeData.Array = make([]*xuanydata.XuanyuanRangeData, 0, len(dAtAs.xuanyuanRangeData.Map))
	for _, k := range uint64Keys {
		dAtAs.xuanyuanRangeData.Array = append(dAtAs.xuanyuanRangeData.Array, dAtAs.xuanyuanRangeData.Map[k])
	}
	dAtAs.xuanyuanRangeData.MinKeyData = dAtAs.xuanyuanRangeData.Array[0]
	dAtAs.xuanyuanRangeData.MaxKeyData = dAtAs.xuanyuanRangeData.Array[len(dAtAs.xuanyuanRangeData.Array)-1]

	// xuanydata.XuanyuanRankPrizeData
	dAtAs.xuanyuanRankPrizeData = &XuanyuanRankPrizeDataConfig{}
	dAtAs.xuanyuanRankPrizeData.Map, dAtAs.xuanyuanRankPrizeData.parserMap, err = xuanydata.LoadXuanyuanRankPrizeData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.xuanyuanRankPrizeData.Map))
	for k := range dAtAs.xuanyuanRankPrizeData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.xuanyuanRankPrizeData.Array = make([]*xuanydata.XuanyuanRankPrizeData, 0, len(dAtAs.xuanyuanRankPrizeData.Map))
	for _, k := range uint64Keys {
		dAtAs.xuanyuanRankPrizeData.Array = append(dAtAs.xuanyuanRankPrizeData.Array, dAtAs.xuanyuanRankPrizeData.Map[k])
	}
	dAtAs.xuanyuanRankPrizeData.MinKeyData = dAtAs.xuanyuanRankPrizeData.Array[0]
	dAtAs.xuanyuanRankPrizeData.MaxKeyData = dAtAs.xuanyuanRankPrizeData.Array[len(dAtAs.xuanyuanRankPrizeData.Array)-1]

	// zhanjiang.ZhanJiangChapterData
	dAtAs.zhanJiangChapterData = &ZhanJiangChapterDataConfig{}
	dAtAs.zhanJiangChapterData.Map, dAtAs.zhanJiangChapterData.parserMap, err = zhanjiang.LoadZhanJiangChapterData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.zhanJiangChapterData.Map))
	for k := range dAtAs.zhanJiangChapterData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.zhanJiangChapterData.Array = make([]*zhanjiang.ZhanJiangChapterData, 0, len(dAtAs.zhanJiangChapterData.Map))
	for _, k := range uint64Keys {
		dAtAs.zhanJiangChapterData.Array = append(dAtAs.zhanJiangChapterData.Array, dAtAs.zhanJiangChapterData.Map[k])
	}
	dAtAs.zhanJiangChapterData.MinKeyData = dAtAs.zhanJiangChapterData.Array[0]
	dAtAs.zhanJiangChapterData.MaxKeyData = dAtAs.zhanJiangChapterData.Array[len(dAtAs.zhanJiangChapterData.Array)-1]

	// zhanjiang.ZhanJiangData
	dAtAs.zhanJiangData = &ZhanJiangDataConfig{}
	dAtAs.zhanJiangData.Map, dAtAs.zhanJiangData.parserMap, err = zhanjiang.LoadZhanJiangData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.zhanJiangData.Map))
	for k := range dAtAs.zhanJiangData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.zhanJiangData.Array = make([]*zhanjiang.ZhanJiangData, 0, len(dAtAs.zhanJiangData.Map))
	for _, k := range uint64Keys {
		dAtAs.zhanJiangData.Array = append(dAtAs.zhanJiangData.Array, dAtAs.zhanJiangData.Map[k])
	}
	dAtAs.zhanJiangData.MinKeyData = dAtAs.zhanJiangData.Array[0]
	dAtAs.zhanJiangData.MaxKeyData = dAtAs.zhanJiangData.Array[len(dAtAs.zhanJiangData.Array)-1]

	// zhanjiang.ZhanJiangGuanQiaData
	dAtAs.zhanJiangGuanQiaData = &ZhanJiangGuanQiaDataConfig{}
	dAtAs.zhanJiangGuanQiaData.Map, dAtAs.zhanJiangGuanQiaData.parserMap, err = zhanjiang.LoadZhanJiangGuanQiaData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.zhanJiangGuanQiaData.Map))
	for k := range dAtAs.zhanJiangGuanQiaData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.zhanJiangGuanQiaData.Array = make([]*zhanjiang.ZhanJiangGuanQiaData, 0, len(dAtAs.zhanJiangGuanQiaData.Map))
	for _, k := range uint64Keys {
		dAtAs.zhanJiangGuanQiaData.Array = append(dAtAs.zhanJiangGuanQiaData.Array, dAtAs.zhanJiangGuanQiaData.Map[k])
	}
	dAtAs.zhanJiangGuanQiaData.MinKeyData = dAtAs.zhanJiangGuanQiaData.Array[0]
	dAtAs.zhanJiangGuanQiaData.MaxKeyData = dAtAs.zhanJiangGuanQiaData.Array[len(dAtAs.zhanJiangGuanQiaData.Array)-1]

	// zhengwu.ZhengWuCompleteData
	dAtAs.zhengWuCompleteData = &ZhengWuCompleteDataConfig{}
	dAtAs.zhengWuCompleteData.Map, dAtAs.zhengWuCompleteData.parserMap, err = zhengwu.LoadZhengWuCompleteData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.zhengWuCompleteData.Map))
	for k := range dAtAs.zhengWuCompleteData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.zhengWuCompleteData.Array = make([]*zhengwu.ZhengWuCompleteData, 0, len(dAtAs.zhengWuCompleteData.Map))
	for _, k := range uint64Keys {
		dAtAs.zhengWuCompleteData.Array = append(dAtAs.zhengWuCompleteData.Array, dAtAs.zhengWuCompleteData.Map[k])
	}
	dAtAs.zhengWuCompleteData.MinKeyData = dAtAs.zhengWuCompleteData.Array[0]
	dAtAs.zhengWuCompleteData.MaxKeyData = dAtAs.zhengWuCompleteData.Array[len(dAtAs.zhengWuCompleteData.Array)-1]

	// zhengwu.ZhengWuData
	dAtAs.zhengWuData = &ZhengWuDataConfig{}
	dAtAs.zhengWuData.Map, dAtAs.zhengWuData.parserMap, err = zhengwu.LoadZhengWuData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.zhengWuData.Map))
	for k := range dAtAs.zhengWuData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.zhengWuData.Array = make([]*zhengwu.ZhengWuData, 0, len(dAtAs.zhengWuData.Map))
	for _, k := range uint64Keys {
		dAtAs.zhengWuData.Array = append(dAtAs.zhengWuData.Array, dAtAs.zhengWuData.Map[k])
	}
	dAtAs.zhengWuData.MinKeyData = dAtAs.zhengWuData.Array[0]
	dAtAs.zhengWuData.MaxKeyData = dAtAs.zhengWuData.Array[len(dAtAs.zhengWuData.Array)-1]

	// zhengwu.ZhengWuRefreshData
	dAtAs.zhengWuRefreshData = &ZhengWuRefreshDataConfig{}
	dAtAs.zhengWuRefreshData.Map, dAtAs.zhengWuRefreshData.parserMap, err = zhengwu.LoadZhengWuRefreshData(gos)
	if err != nil {
		return nil, err
	}

	uint64Keys = make([]uint64, 0, len(dAtAs.zhengWuRefreshData.Map))
	for k := range dAtAs.zhengWuRefreshData.Map {
		uint64Keys = append(uint64Keys, k)
	}
	sort.Sort(uint64Slice(uint64Keys))
	dAtAs.zhengWuRefreshData.Array = make([]*zhengwu.ZhengWuRefreshData, 0, len(dAtAs.zhengWuRefreshData.Map))
	for _, k := range uint64Keys {
		dAtAs.zhengWuRefreshData.Array = append(dAtAs.zhengWuRefreshData.Array, dAtAs.zhengWuRefreshData.Map[k])
	}
	dAtAs.zhengWuRefreshData.MinKeyData = dAtAs.zhengWuRefreshData.Array[0]
	dAtAs.zhengWuRefreshData.MaxKeyData = dAtAs.zhengWuRefreshData.Array[len(dAtAs.zhengWuRefreshData.Array)-1]

	// bai_zhan_data.BaiZhanMiscData
	dAtAs.baiZhanMiscData, dAtAs.baiZhanMiscDataParser, err = bai_zhan_data.LoadBaiZhanMiscData(gos)
	if err != nil {
		return nil, err
	}

	// combatdata.CombatConfig
	dAtAs.combatConfig, dAtAs.combatConfigParser, err = combatdata.LoadCombatConfig(gos)
	if err != nil {
		return nil, err
	}

	// combatdata.CombatMiscConfig
	dAtAs.combatMiscConfig, dAtAs.combatMiscConfigParser, err = combatdata.LoadCombatMiscConfig(gos)
	if err != nil {
		return nil, err
	}

	// combine.EquipCombineDatas
	dAtAs.equipCombineDatas, dAtAs.equipCombineDatasParser, err = combine.LoadEquipCombineDatas(gos)
	if err != nil {
		return nil, err
	}

	// country.CountryMiscData
	dAtAs.countryMiscData, dAtAs.countryMiscDataParser, err = country.LoadCountryMiscData(gos)
	if err != nil {
		return nil, err
	}

	// data.BroadcastHelp
	dAtAs.broadcastHelp, dAtAs.broadcastHelpParser, err = data.LoadBroadcastHelp(gos)
	if err != nil {
		return nil, err
	}

	// data.TextHelp
	dAtAs.textHelp, dAtAs.textHelpParser, err = data.LoadTextHelp(gos)
	if err != nil {
		return nil, err
	}

	// dianquan.ExchangeMiscData
	dAtAs.exchangeMiscData, dAtAs.exchangeMiscDataParser, err = dianquan.LoadExchangeMiscData(gos)
	if err != nil {
		return nil, err
	}

	// domestic_data.BuildingLayoutMiscData
	dAtAs.buildingLayoutMiscData, dAtAs.buildingLayoutMiscDataParser, err = domestic_data.LoadBuildingLayoutMiscData(gos)
	if err != nil {
		return nil, err
	}

	// domestic_data.CityEventMiscData
	dAtAs.cityEventMiscData, dAtAs.cityEventMiscDataParser, err = domestic_data.LoadCityEventMiscData(gos)
	if err != nil {
		return nil, err
	}

	// domestic_data.MainCityMiscData
	dAtAs.mainCityMiscData, dAtAs.mainCityMiscDataParser, err = domestic_data.LoadMainCityMiscData(gos)
	if err != nil {
		return nil, err
	}

	// dungeon.DungeonMiscData
	dAtAs.dungeonMiscData, dAtAs.dungeonMiscDataParser, err = dungeon.LoadDungeonMiscData(gos)
	if err != nil {
		return nil, err
	}

	// farm.FarmMiscConfig
	dAtAs.farmMiscConfig, dAtAs.farmMiscConfigParser, err = farm.LoadFarmMiscConfig(gos)
	if err != nil {
		return nil, err
	}

	// fishing_data.FishRandomer
	dAtAs.fishRandomer, dAtAs.fishRandomerParser, err = fishing_data.LoadFishRandomer(gos)
	if err != nil {
		return nil, err
	}

	// gardendata.GardenConfig
	dAtAs.gardenConfig, dAtAs.gardenConfigParser, err = gardendata.LoadGardenConfig(gos)
	if err != nil {
		return nil, err
	}

	// goods.EquipmentTaozConfig
	dAtAs.equipmentTaozConfig, dAtAs.equipmentTaozConfigParser, err = goods.LoadEquipmentTaozConfig(gos)
	if err != nil {
		return nil, err
	}

	// goods.GemDatas
	dAtAs.gemDatas, dAtAs.gemDatasParser, err = goods.LoadGemDatas(gos)
	if err != nil {
		return nil, err
	}

	// goods.GoodsCheck
	dAtAs.goodsCheck, dAtAs.goodsCheckParser, err = goods.LoadGoodsCheck(gos)
	if err != nil {
		return nil, err
	}

	// guild_data.GuildLogHelp
	dAtAs.guildLogHelp, dAtAs.guildLogHelpParser, err = guild_data.LoadGuildLogHelp(gos)
	if err != nil {
		return nil, err
	}

	// guild_data.NpcGuildSuffixName
	dAtAs.npcGuildSuffixName, dAtAs.npcGuildSuffixNameParser, err = guild_data.LoadNpcGuildSuffixName(gos)
	if err != nil {
		return nil, err
	}

	// hebi.HebiMiscData
	dAtAs.hebiMiscData, dAtAs.hebiMiscDataParser, err = hebi.LoadHebiMiscData(gos)
	if err != nil {
		return nil, err
	}

	// heroinit.HeroCreateData
	dAtAs.heroCreateData, dAtAs.heroCreateDataParser, err = heroinit.LoadHeroCreateData(gos)
	if err != nil {
		return nil, err
	}

	// heroinit.HeroInitData
	dAtAs.heroInitData, dAtAs.heroInitDataParser, err = heroinit.LoadHeroInitData(gos)
	if err != nil {
		return nil, err
	}

	// maildata.MailHelp
	dAtAs.mailHelp, dAtAs.mailHelpParser, err = maildata.LoadMailHelp(gos)
	if err != nil {
		return nil, err
	}

	// military_data.JiuGuanMiscData
	dAtAs.jiuGuanMiscData, dAtAs.jiuGuanMiscDataParser, err = military_data.LoadJiuGuanMiscData(gos)
	if err != nil {
		return nil, err
	}

	// military_data.JunYingMiscData
	dAtAs.junYingMiscData, dAtAs.junYingMiscDataParser, err = military_data.LoadJunYingMiscData(gos)
	if err != nil {
		return nil, err
	}

	// mingcdata.McBuildMiscData
	dAtAs.mcBuildMiscData, dAtAs.mcBuildMiscDataParser, err = mingcdata.LoadMcBuildMiscData(gos)
	if err != nil {
		return nil, err
	}

	// mingcdata.MingcMiscData
	dAtAs.mingcMiscData, dAtAs.mingcMiscDataParser, err = mingcdata.LoadMingcMiscData(gos)
	if err != nil {
		return nil, err
	}

	// promdata.EventLimitGiftConfig
	dAtAs.eventLimitGiftConfig, dAtAs.eventLimitGiftConfigParser, err = promdata.LoadEventLimitGiftConfig(gos)
	if err != nil {
		return nil, err
	}

	// promdata.PromotionMiscData
	dAtAs.promotionMiscData, dAtAs.promotionMiscDataParser, err = promdata.LoadPromotionMiscData(gos)
	if err != nil {
		return nil, err
	}

	// question.QuestionMiscData
	dAtAs.questionMiscData, dAtAs.questionMiscDataParser, err = question.LoadQuestionMiscData(gos)
	if err != nil {
		return nil, err
	}

	// race.RaceConfig
	dAtAs.raceConfig, dAtAs.raceConfigParser, err = race.LoadRaceConfig(gos)
	if err != nil {
		return nil, err
	}

	// random_event.RandomEventDataDictionary
	dAtAs.randomEventDataDictionary, dAtAs.randomEventDataDictionaryParser, err = random_event.LoadRandomEventDataDictionary(gos)
	if err != nil {
		return nil, err
	}

	// random_event.RandomEventPositionDictionary
	dAtAs.randomEventPositionDictionary, dAtAs.randomEventPositionDictionaryParser, err = random_event.LoadRandomEventPositionDictionary(gos)
	if err != nil {
		return nil, err
	}

	// rank_data.RankMiscData
	dAtAs.rankMiscData, dAtAs.rankMiscDataParser, err = rank_data.LoadRankMiscData(gos)
	if err != nil {
		return nil, err
	}

	// regdata.JunTuanNpcPlaceConfig
	dAtAs.junTuanNpcPlaceConfig, dAtAs.junTuanNpcPlaceConfigParser, err = regdata.LoadJunTuanNpcPlaceConfig(gos)
	if err != nil {
		return nil, err
	}

	// season.SeasonMiscData
	dAtAs.seasonMiscData, dAtAs.seasonMiscDataParser, err = season.LoadSeasonMiscData(gos)
	if err != nil {
		return nil, err
	}

	// settings.SettingMiscData
	dAtAs.settingMiscData, dAtAs.settingMiscDataParser, err = settings.LoadSettingMiscData(gos)
	if err != nil {
		return nil, err
	}

	// shop.ShopMiscData
	dAtAs.shopMiscData, dAtAs.shopMiscDataParser, err = shop.LoadShopMiscData(gos)
	if err != nil {
		return nil, err
	}

	// singleton.GoodsConfig
	dAtAs.goodsConfig, dAtAs.goodsConfigParser, err = singleton.LoadGoodsConfig(gos)
	if err != nil {
		return nil, err
	}

	// singleton.GuildConfig
	dAtAs.guildConfig, dAtAs.guildConfigParser, err = singleton.LoadGuildConfig(gos)
	if err != nil {
		return nil, err
	}

	// singleton.GuildGenConfig
	dAtAs.guildGenConfig, dAtAs.guildGenConfigParser, err = singleton.LoadGuildGenConfig(gos)
	if err != nil {
		return nil, err
	}

	// singleton.MilitaryConfig
	dAtAs.militaryConfig, dAtAs.militaryConfigParser, err = singleton.LoadMilitaryConfig(gos)
	if err != nil {
		return nil, err
	}

	// singleton.MiscConfig
	dAtAs.miscConfig, dAtAs.miscConfigParser, err = singleton.LoadMiscConfig(gos)
	if err != nil {
		return nil, err
	}

	// singleton.MiscGenConfig
	dAtAs.miscGenConfig, dAtAs.miscGenConfigParser, err = singleton.LoadMiscGenConfig(gos)
	if err != nil {
		return nil, err
	}

	// singleton.RegionConfig
	dAtAs.regionConfig, dAtAs.regionConfigParser, err = singleton.LoadRegionConfig(gos)
	if err != nil {
		return nil, err
	}

	// singleton.RegionGenConfig
	dAtAs.regionGenConfig, dAtAs.regionGenConfigParser, err = singleton.LoadRegionGenConfig(gos)
	if err != nil {
		return nil, err
	}

	// tag.TagMiscData
	dAtAs.tagMiscData, dAtAs.tagMiscDataParser, err = tag.LoadTagMiscData(gos)
	if err != nil {
		return nil, err
	}

	// taskdata.TaskMiscData
	dAtAs.taskMiscData, dAtAs.taskMiscDataParser, err = taskdata.LoadTaskMiscData(gos)
	if err != nil {
		return nil, err
	}

	// towerdata.SecretTowerMiscData
	dAtAs.secretTowerMiscData, dAtAs.secretTowerMiscDataParser, err = towerdata.LoadSecretTowerMiscData(gos)
	if err != nil {
		return nil, err
	}

	// vip.VipMiscData
	dAtAs.vipMiscData, dAtAs.vipMiscDataParser, err = vip.LoadVipMiscData(gos)
	if err != nil {
		return nil, err
	}

	// xiongnu.ResistXiongNuMisc
	dAtAs.resistXiongNuMisc, dAtAs.resistXiongNuMiscParser, err = xiongnu.LoadResistXiongNuMisc(gos)
	if err != nil {
		return nil, err
	}

	// xuanydata.XuanyuanMiscData
	dAtAs.xuanyuanMiscData, dAtAs.xuanyuanMiscDataParser, err = xuanydata.LoadXuanyuanMiscData(gos)
	if err != nil {
		return nil, err
	}

	// zhanjiang.ZhanJiangMiscData
	dAtAs.zhanJiangMiscData, dAtAs.zhanJiangMiscDataParser, err = zhanjiang.LoadZhanJiangMiscData(gos)
	if err != nil {
		return nil, err
	}

	// zhengwu.ZhengWuMiscData
	dAtAs.zhengWuMiscData, dAtAs.zhengWuMiscDataParser, err = zhengwu.LoadZhengWuMiscData(gos)
	if err != nil {
		return nil, err
	}

	// zhengwu.ZhengWuRandomData
	dAtAs.zhengWuRandomData, dAtAs.zhengWuRandomDataParser, err = zhengwu.LoadZhengWuRandomData(gos)
	if err != nil {
		return nil, err
	}

	// SetRelatedData

	if err := activitydata.SetRelatedActivityCollectionData(dAtAs.activityCollectionData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := activitydata.SetRelatedActivityShowData(dAtAs.activityShowData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := activitydata.SetRelatedActivityTaskListModeData(dAtAs.activityTaskListModeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := activitydata.SetRelatedCollectionExchangeData(dAtAs.collectionExchangeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := bai_zhan_data.SetRelatedJunXianLevelData(dAtAs.junXianLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := bai_zhan_data.SetRelatedJunXianPrizeData(dAtAs.junXianPrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := basedata.SetRelatedHomeNpcBaseData(dAtAs.homeNpcBaseData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := basedata.SetRelatedNpcBaseData(dAtAs.npcBaseData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := blockdata.SetRelatedBlockData(dAtAs.blockData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := body.SetRelatedBodyData(dAtAs.bodyData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := buffer.SetRelatedBufferData(dAtAs.bufferData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := buffer.SetRelatedBufferTypeData(dAtAs.bufferTypeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := captain.SetRelatedCaptainAbilityData(dAtAs.captainAbilityData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := captain.SetRelatedCaptainData(dAtAs.captainData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := captain.SetRelatedCaptainFriendshipData(dAtAs.captainFriendshipData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := captain.SetRelatedCaptainLevelData(dAtAs.captainLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := captain.SetRelatedCaptainOfficialCountData(dAtAs.captainOfficialCountData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := captain.SetRelatedCaptainOfficialData(dAtAs.captainOfficialData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := captain.SetRelatedCaptainRarityData(dAtAs.captainRarityData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := captain.SetRelatedCaptainRebirthLevelData(dAtAs.captainRebirthLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := captain.SetRelatedCaptainStarData(dAtAs.captainStarData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := captain.SetRelatedNamelessCaptainData(dAtAs.namelessCaptainData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := charge.SetRelatedChargeObjData(dAtAs.chargeObjData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := charge.SetRelatedChargePrizeData(dAtAs.chargePrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := charge.SetRelatedProductData(dAtAs.productData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := combine.SetRelatedEquipCombineData(dAtAs.equipCombineData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := combine.SetRelatedGoodsCombineData(dAtAs.goodsCombineData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := country.SetRelatedCountryData(dAtAs.countryData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := country.SetRelatedCountryOfficialData(dAtAs.countryOfficialData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := country.SetRelatedCountryOfficialNpcData(dAtAs.countryOfficialNpcData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := country.SetRelatedFamilyNameData(dAtAs.familyNameData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedBroadcastData(dAtAs.broadcastData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedBuffEffectData(dAtAs.buffEffectData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedColorData(dAtAs.colorData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedFamilyName(dAtAs.familyName.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedFemaleGivenName(dAtAs.femaleGivenName.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedHeroLevelSubData(dAtAs.heroLevelSubData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedMaleGivenName(dAtAs.maleGivenName.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedSpriteStat(dAtAs.spriteStat.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedText(dAtAs.text.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedTimeRuleData(dAtAs.timeRuleData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedBaseLevelData(dAtAs.baseLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedBuildingData(dAtAs.buildingData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedBuildingLayoutData(dAtAs.buildingLayoutData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedBuildingUnlockData(dAtAs.buildingUnlockData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedCityEventData(dAtAs.cityEventData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedCityEventLevelData(dAtAs.cityEventLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedCombineCost(dAtAs.combineCost.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedCountdownPrizeData(dAtAs.countdownPrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedCountdownPrizeDescData(dAtAs.countdownPrizeDescData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedGuanFuLevelData(dAtAs.guanFuLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedOuterCityBuildingData(dAtAs.outerCityBuildingData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedOuterCityBuildingDescData(dAtAs.outerCityBuildingDescData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedOuterCityData(dAtAs.outerCityData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedOuterCityLayoutData(dAtAs.outerCityLayoutData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedProsperityDamageBuffData(dAtAs.prosperityDamageBuffData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedSoldierLevelData(dAtAs.soldierLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedTechnologyData(dAtAs.technologyData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedTieJiangPuLevelData(dAtAs.tieJiangPuLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedWorkshopDuration(dAtAs.workshopDuration.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedWorkshopLevelData(dAtAs.workshopLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedWorkshopRefreshCost(dAtAs.workshopRefreshCost.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := dungeon.SetRelatedDungeonChapterData(dAtAs.dungeonChapterData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := dungeon.SetRelatedDungeonData(dAtAs.dungeonData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := dungeon.SetRelatedDungeonGuideTroopData(dAtAs.dungeonGuideTroopData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := farm.SetRelatedFarmMaxStealConfig(dAtAs.farmMaxStealConfig.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := farm.SetRelatedFarmOneKeyConfig(dAtAs.farmOneKeyConfig.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := farm.SetRelatedFarmResConfig(dAtAs.farmResConfig.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := fishing_data.SetRelatedFishData(dAtAs.fishData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := fishing_data.SetRelatedFishingCaptainProbabilityData(dAtAs.fishingCaptainProbabilityData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := fishing_data.SetRelatedFishingCostData(dAtAs.fishingCostData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := fishing_data.SetRelatedFishingShowData(dAtAs.fishingShowData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := function.SetRelatedFunctionOpenData(dAtAs.functionOpenData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := gardendata.SetRelatedTreasuryTreeData(dAtAs.treasuryTreeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedEquipmentData(dAtAs.equipmentData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedEquipmentLevelData(dAtAs.equipmentLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedEquipmentQualityData(dAtAs.equipmentQualityData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedEquipmentRefinedData(dAtAs.equipmentRefinedData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedEquipmentTaozData(dAtAs.equipmentTaozData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedGemData(dAtAs.gemData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedGoodsData(dAtAs.goodsData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedGoodsQuality(dAtAs.goodsQuality.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildBigBoxData(dAtAs.guildBigBoxData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildClassLevelData(dAtAs.guildClassLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildClassTitleData(dAtAs.guildClassTitleData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildDonateData(dAtAs.guildDonateData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildEventPrizeData(dAtAs.guildEventPrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildLevelCdrData(dAtAs.guildLevelCdrData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildLevelData(dAtAs.guildLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildLogData(dAtAs.guildLogData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildPermissionShowData(dAtAs.guildPermissionShowData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildPrestigeEventData(dAtAs.guildPrestigeEventData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildPrestigePrizeData(dAtAs.guildPrestigePrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildRankPrizeData(dAtAs.guildRankPrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildTarget(dAtAs.guildTarget.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildTaskData(dAtAs.guildTaskData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildTaskEvaluateData(dAtAs.guildTaskEvaluateData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildTechnologyData(dAtAs.guildTechnologyData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedNpcGuildTemplate(dAtAs.npcGuildTemplate.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedNpcMemberData(dAtAs.npcMemberData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := head.SetRelatedHeadData(dAtAs.headData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := hebi.SetRelatedHebiPrizeData(dAtAs.hebiPrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := herodata.SetRelatedHeroLevelData(dAtAs.heroLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := i18n.SetRelatedI18nData(dAtAs.i18nData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := icon.SetRelatedIcon(dAtAs.icon.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := location.SetRelatedLocationData(dAtAs.locationData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := maildata.SetRelatedMailData(dAtAs.mailData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := military_data.SetRelatedJiuGuanData(dAtAs.jiuGuanData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := military_data.SetRelatedJunYingLevelData(dAtAs.junYingLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := military_data.SetRelatedTrainingLevelData(dAtAs.trainingLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := military_data.SetRelatedTutorData(dAtAs.tutorData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMcBuildAddSupportData(dAtAs.mcBuildAddSupportData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMcBuildGuildMemberPrizeData(dAtAs.mcBuildGuildMemberPrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMcBuildMcSupportData(dAtAs.mcBuildMcSupportData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcBaseData(dAtAs.mingcBaseData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcTimeData(dAtAs.mingcTimeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcWarBuildingData(dAtAs.mingcWarBuildingData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcWarDrumStatData(dAtAs.mingcWarDrumStatData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcWarMapData(dAtAs.mingcWarMapData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcWarMultiKillData(dAtAs.mingcWarMultiKillData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcWarNpcData(dAtAs.mingcWarNpcData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcWarNpcGuildData(dAtAs.mingcWarNpcGuildData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcWarSceneData(dAtAs.mingcWarSceneData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcWarTouShiBuildingTargetData(dAtAs.mingcWarTouShiBuildingTargetData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcWarTroopLastBeatWhenFailData(dAtAs.mingcWarTroopLastBeatWhenFailData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := monsterdata.SetRelatedMonsterCaptainData(dAtAs.monsterCaptainData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := monsterdata.SetRelatedMonsterMasterData(dAtAs.monsterMasterData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedDailyBargainData(dAtAs.dailyBargainData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedDurationCardData(dAtAs.durationCardData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedEventLimitGiftData(dAtAs.eventLimitGiftData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedFreeGiftData(dAtAs.freeGiftData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedHeroLevelFundData(dAtAs.heroLevelFundData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedLoginDayData(dAtAs.loginDayData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedSpCollectionData(dAtAs.spCollectionData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedTimeLimitGiftData(dAtAs.timeLimitGiftData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedTimeLimitGiftGroupData(dAtAs.timeLimitGiftGroupData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := pushdata.SetRelatedPushData(dAtAs.pushData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := pvetroop.SetRelatedPveTroopData(dAtAs.pveTroopData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := question.SetRelatedQuestionData(dAtAs.questionData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := question.SetRelatedQuestionPrizeData(dAtAs.questionPrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := question.SetRelatedQuestionSayingData(dAtAs.questionSayingData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := race.SetRelatedRaceData(dAtAs.raceData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := random_event.SetRelatedEventOptionData(dAtAs.eventOptionData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := random_event.SetRelatedEventPosition(dAtAs.eventPosition.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := random_event.SetRelatedOptionPrize(dAtAs.optionPrize.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := random_event.SetRelatedRandomEventData(dAtAs.randomEventData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := red_packet.SetRelatedRedPacketData(dAtAs.redPacketData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedAreaData(dAtAs.areaData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedAssemblyData(dAtAs.assemblyData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedBaozNpcData(dAtAs.baozNpcData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedJunTuanNpcData(dAtAs.junTuanNpcData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedJunTuanNpcPlaceData(dAtAs.junTuanNpcPlaceData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedRegionAreaData(dAtAs.regionAreaData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedRegionData(dAtAs.regionData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedRegionMonsterData(dAtAs.regionMonsterData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedRegionMultiLevelNpcData(dAtAs.regionMultiLevelNpcData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedRegionMultiLevelNpcLevelData(dAtAs.regionMultiLevelNpcLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedRegionMultiLevelNpcTypeData(dAtAs.regionMultiLevelNpcTypeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedTroopDialogueData(dAtAs.troopDialogueData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedTroopDialogueTextData(dAtAs.troopDialogueTextData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedAmountShowSortData(dAtAs.amountShowSortData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedBaowuData(dAtAs.baowuData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedConditionPlunder(dAtAs.conditionPlunder.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedConditionPlunderItem(dAtAs.conditionPlunderItem.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedCost(dAtAs.cost.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedGuildLevelPrize(dAtAs.guildLevelPrize.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedPlunder(dAtAs.plunder.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedPlunderGroup(dAtAs.plunderGroup.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedPlunderItem(dAtAs.plunderItem.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedPlunderPrize(dAtAs.plunderPrize.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedPrize(dAtAs.prize.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := resdata.SetRelatedResCaptainData(dAtAs.resCaptainData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := scene.SetRelatedCombatScene(dAtAs.combatScene.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := season.SetRelatedSeasonData(dAtAs.seasonData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := settings.SetRelatedPrivacySettingData(dAtAs.privacySettingData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := shop.SetRelatedBlackMarketData(dAtAs.blackMarketData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := shop.SetRelatedBlackMarketGoodsData(dAtAs.blackMarketGoodsData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := shop.SetRelatedBlackMarketGoodsGroupData(dAtAs.blackMarketGoodsGroupData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := shop.SetRelatedDiscountColorData(dAtAs.discountColorData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := shop.SetRelatedNormalShopGoods(dAtAs.normalShopGoods.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := shop.SetRelatedShop(dAtAs.shop.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := shop.SetRelatedZhenBaoGeShopGoods(dAtAs.zhenBaoGeShopGoods.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := spell.SetRelatedPassiveSpellData(dAtAs.passiveSpellData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := spell.SetRelatedSpell(dAtAs.spell.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := spell.SetRelatedSpellData(dAtAs.spellData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := spell.SetRelatedSpellFacadeData(dAtAs.spellFacadeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := spell.SetRelatedStateData(dAtAs.stateData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := strategydata.SetRelatedStrategyData(dAtAs.strategyData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := strategydata.SetRelatedStrategyEffectData(dAtAs.strategyEffectData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := strongerdata.SetRelatedStrongerData(dAtAs.strongerData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := sub.SetRelatedBuildingEffectData(dAtAs.buildingEffectData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := survey.SetRelatedSurveyData(dAtAs.surveyData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedAchieveTaskData(dAtAs.achieveTaskData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedAchieveTaskStarPrizeData(dAtAs.achieveTaskStarPrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedActiveDegreePrizeData(dAtAs.activeDegreePrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedActiveDegreeTaskData(dAtAs.activeDegreeTaskData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedActivityTaskData(dAtAs.activityTaskData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedBaYeStageData(dAtAs.baYeStageData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedBaYeTaskData(dAtAs.baYeTaskData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedBranchTaskData(dAtAs.branchTaskData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedBwzlPrizeData(dAtAs.bwzlPrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedBwzlTaskData(dAtAs.bwzlTaskData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedMainTaskData(dAtAs.mainTaskData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedTaskBoxData(dAtAs.taskBoxData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedTaskTargetData(dAtAs.taskTargetData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedTitleData(dAtAs.titleData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedTitleTaskData(dAtAs.titleTaskData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := teach.SetRelatedTeachChapterData(dAtAs.teachChapterData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := towerdata.SetRelatedSecretTowerData(dAtAs.secretTowerData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := towerdata.SetRelatedSecretTowerWordsData(dAtAs.secretTowerWordsData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := towerdata.SetRelatedTowerData(dAtAs.towerData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := vip.SetRelatedVipContinueDaysData(dAtAs.vipContinueDaysData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := vip.SetRelatedVipLevelData(dAtAs.vipLevelData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := xiongnu.SetRelatedResistXiongNuData(dAtAs.resistXiongNuData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := xiongnu.SetRelatedResistXiongNuScoreData(dAtAs.resistXiongNuScoreData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := xiongnu.SetRelatedResistXiongNuWaveData(dAtAs.resistXiongNuWaveData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := xuanydata.SetRelatedXuanyuanRangeData(dAtAs.xuanyuanRangeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := xuanydata.SetRelatedXuanyuanRankPrizeData(dAtAs.xuanyuanRankPrizeData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := zhanjiang.SetRelatedZhanJiangChapterData(dAtAs.zhanJiangChapterData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := zhanjiang.SetRelatedZhanJiangData(dAtAs.zhanJiangData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := zhanjiang.SetRelatedZhanJiangGuanQiaData(dAtAs.zhanJiangGuanQiaData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := zhengwu.SetRelatedZhengWuCompleteData(dAtAs.zhengWuCompleteData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := zhengwu.SetRelatedZhengWuData(dAtAs.zhengWuData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := zhengwu.SetRelatedZhengWuRefreshData(dAtAs.zhengWuRefreshData.parserMap, dAtAs); err != nil {
		return nil, err
	}

	if err := bai_zhan_data.SetRelatedBaiZhanMiscData(gos, dAtAs.baiZhanMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := combatdata.SetRelatedCombatConfig(gos, dAtAs.combatConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := combatdata.SetRelatedCombatMiscConfig(gos, dAtAs.combatMiscConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := combine.SetRelatedEquipCombineDatas(gos, dAtAs.equipCombineDatas, dAtAs); err != nil {
		return nil, err
	}

	if err := country.SetRelatedCountryMiscData(gos, dAtAs.countryMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedBroadcastHelp(gos, dAtAs.broadcastHelp, dAtAs); err != nil {
		return nil, err
	}

	if err := data.SetRelatedTextHelp(gos, dAtAs.textHelp, dAtAs); err != nil {
		return nil, err
	}

	if err := dianquan.SetRelatedExchangeMiscData(gos, dAtAs.exchangeMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedBuildingLayoutMiscData(gos, dAtAs.buildingLayoutMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedCityEventMiscData(gos, dAtAs.cityEventMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := domestic_data.SetRelatedMainCityMiscData(gos, dAtAs.mainCityMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := dungeon.SetRelatedDungeonMiscData(gos, dAtAs.dungeonMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := farm.SetRelatedFarmMiscConfig(gos, dAtAs.farmMiscConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := fishing_data.SetRelatedFishRandomer(gos, dAtAs.fishRandomer, dAtAs); err != nil {
		return nil, err
	}

	if err := gardendata.SetRelatedGardenConfig(gos, dAtAs.gardenConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedEquipmentTaozConfig(gos, dAtAs.equipmentTaozConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedGemDatas(gos, dAtAs.gemDatas, dAtAs); err != nil {
		return nil, err
	}

	if err := goods.SetRelatedGoodsCheck(gos, dAtAs.goodsCheck, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedGuildLogHelp(gos, dAtAs.guildLogHelp, dAtAs); err != nil {
		return nil, err
	}

	if err := guild_data.SetRelatedNpcGuildSuffixName(gos, dAtAs.npcGuildSuffixName, dAtAs); err != nil {
		return nil, err
	}

	if err := hebi.SetRelatedHebiMiscData(gos, dAtAs.hebiMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := heroinit.SetRelatedHeroCreateData(gos, dAtAs.heroCreateData, dAtAs); err != nil {
		return nil, err
	}

	if err := heroinit.SetRelatedHeroInitData(gos, dAtAs.heroInitData, dAtAs); err != nil {
		return nil, err
	}

	if err := maildata.SetRelatedMailHelp(gos, dAtAs.mailHelp, dAtAs); err != nil {
		return nil, err
	}

	if err := military_data.SetRelatedJiuGuanMiscData(gos, dAtAs.jiuGuanMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := military_data.SetRelatedJunYingMiscData(gos, dAtAs.junYingMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMcBuildMiscData(gos, dAtAs.mcBuildMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := mingcdata.SetRelatedMingcMiscData(gos, dAtAs.mingcMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedEventLimitGiftConfig(gos, dAtAs.eventLimitGiftConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := promdata.SetRelatedPromotionMiscData(gos, dAtAs.promotionMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := question.SetRelatedQuestionMiscData(gos, dAtAs.questionMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := race.SetRelatedRaceConfig(gos, dAtAs.raceConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := random_event.SetRelatedRandomEventDataDictionary(gos, dAtAs.randomEventDataDictionary, dAtAs); err != nil {
		return nil, err
	}

	if err := random_event.SetRelatedRandomEventPositionDictionary(gos, dAtAs.randomEventPositionDictionary, dAtAs); err != nil {
		return nil, err
	}

	if err := rank_data.SetRelatedRankMiscData(gos, dAtAs.rankMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := regdata.SetRelatedJunTuanNpcPlaceConfig(gos, dAtAs.junTuanNpcPlaceConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := season.SetRelatedSeasonMiscData(gos, dAtAs.seasonMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := settings.SetRelatedSettingMiscData(gos, dAtAs.settingMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := shop.SetRelatedShopMiscData(gos, dAtAs.shopMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := singleton.SetRelatedGoodsConfig(gos, dAtAs.goodsConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := singleton.SetRelatedGuildConfig(gos, dAtAs.guildConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := singleton.SetRelatedGuildGenConfig(gos, dAtAs.guildGenConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := singleton.SetRelatedMilitaryConfig(gos, dAtAs.militaryConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := singleton.SetRelatedMiscConfig(gos, dAtAs.miscConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := singleton.SetRelatedMiscGenConfig(gos, dAtAs.miscGenConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := singleton.SetRelatedRegionConfig(gos, dAtAs.regionConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := singleton.SetRelatedRegionGenConfig(gos, dAtAs.regionGenConfig, dAtAs); err != nil {
		return nil, err
	}

	if err := tag.SetRelatedTagMiscData(gos, dAtAs.tagMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := taskdata.SetRelatedTaskMiscData(gos, dAtAs.taskMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := towerdata.SetRelatedSecretTowerMiscData(gos, dAtAs.secretTowerMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := vip.SetRelatedVipMiscData(gos, dAtAs.vipMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := xiongnu.SetRelatedResistXiongNuMisc(gos, dAtAs.resistXiongNuMisc, dAtAs); err != nil {
		return nil, err
	}

	if err := xuanydata.SetRelatedXuanyuanMiscData(gos, dAtAs.xuanyuanMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := zhanjiang.SetRelatedZhanJiangMiscData(gos, dAtAs.zhanJiangMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := zhengwu.SetRelatedZhengWuMiscData(gos, dAtAs.zhengWuMiscData, dAtAs); err != nil {
		return nil, err
	}

	if err := zhengwu.SetRelatedZhengWuRandomData(gos, dAtAs.zhengWuRandomData, dAtAs); err != nil {
		return nil, err
	}

	// Init

	for _, v := range dAtAs.xuanyuanRangeData.Array {
		v.Init("轩辕会武/积分区间.txt", dAtAs.xuanyuanRangeData.Map)
	}

	dAtAs.xuanyuanMiscData.Init("轩辕会武/杂项.txt", dAtAs)

	for _, v := range dAtAs.buildingEffectData.Array {
		v.Init("内政/建筑效果.txt")
	}

	for _, v := range dAtAs.spellFacadeData.Array {
		v.Init("战斗/技能盒.txt")
	}

	for _, v := range dAtAs.blackMarketGoodsGroupData.Array {
		v.Init("商店/黑市商品分组.txt")
	}

	for _, v := range dAtAs.blackMarketData.Array {
		v.Init("商店/黑市.txt")
	}

	for _, v := range dAtAs.privacySettingData.Array {
		v.Init("设置/隐私设置.txt")
	}

	dAtAs.seasonMiscData.Init("季节/杂项.txt", dAtAs)

	for _, v := range dAtAs.amountShowSortData.Array {
		v.Init(dAtAs.amountShowSortData.parserMap[v], "杂项/展示排序.txt")
	}

	for _, v := range dAtAs.regionMultiLevelNpcTypeData.Array {
		v.Init("地图/多等级野怪类型.txt")
	}

	for _, v := range dAtAs.areaData.Array {
		v.Init("地图/区块链.txt")
	}

	for _, v := range dAtAs.randomEventData.Array {
		v.Init("随机事件/随机事件.txt")
	}

	for _, v := range dAtAs.raceData.Array {
		v.Init()
	}

	dAtAs.raceConfig.Init("杂项/职业杂项.txt", dAtAs)

	for _, v := range dAtAs.questionData.Array {
		v.Init("答题/答题问题.txt")
	}

	for _, v := range dAtAs.pveTroopData.Array {
		v.Init("杂项/pve部队.txt")
	}

	for _, v := range dAtAs.mingcWarTouShiBuildingTargetData.Array {
		v.Init("名城战/投石机目标.txt")
	}

	for _, v := range dAtAs.mingcWarNpcGuildData.Array {
		v.Init("名城战/初始城主联盟.txt")
	}

	for _, v := range dAtAs.mingcWarMapData.Array {
		v.Init("名城战/地图.txt")
	}

	for _, v := range dAtAs.mingcWarDrumStatData.Array {
		v.Init("名城战/鼓舞加成.txt")
	}

	for _, v := range dAtAs.mingcWarBuildingData.Array {
		v.Init("名城战/据点.txt")
	}

	dAtAs.mingcMiscData.Init("名城战/杂项.txt")

	for _, v := range dAtAs.mcBuildAddSupportData.Array {
		v.Init("名城营建/增加民心.txt", dAtAs)
	}

	for _, v := range dAtAs.i18nData.Array {
		v.Init99(gos, "i18n/语言.txt")
	}

	dAtAs.npcGuildSuffixName.Init("联盟/联盟Npc后缀名.txt")

	for _, v := range dAtAs.guildPrestigeEventData.Array {
		v.Init("联盟/联盟声望事件.txt")
	}

	for _, v := range dAtAs.guildPermissionShowData.Array {
		v.Init("联盟/联盟权限.txt", dAtAs)
	}

	for _, v := range dAtAs.guildLevelData.Array {
		v.Init("联盟/联盟等级.txt", dAtAs.guildLevelData.Map, dAtAs)
	}

	for _, v := range dAtAs.goodsData.Array {
		v.Init("物品/物品.txt", dAtAs)
	}

	for _, v := range dAtAs.gemData.Array {
		v.Init("物品/宝石.txt", dAtAs.gemData.Map)
	}

	for _, v := range dAtAs.equipmentRefinedData.Array {
		v.Init("物品/装备强化.txt", dAtAs.equipmentRefinedData.Map)
	}

	for _, v := range dAtAs.equipmentLevelData.Array {
		v.Init("物品/装备等级.txt", dAtAs.equipmentLevelData.Map)
	}

	dAtAs.farmMiscConfig.Init("农场/农场杂项.txt", dAtAs)

	for _, v := range dAtAs.workshopDuration.Array {
		v.Init("内政/装备作坊时间.txt")
	}

	dAtAs.mainCityMiscData.Init("内政/主城杂项.txt")

	for _, v := range dAtAs.buildingLayoutData.Array {
		v.Init("内政/建筑布局.txt")
	}

	for _, v := range dAtAs.baseLevelData.Array {
		v.Init("内政/主城等级.txt", dAtAs.baseLevelData.Map)
	}

	dAtAs.exchangeMiscData.Init("点券/点券杂项.txt")

	for _, v := range dAtAs.timeRuleData.Array {
		v.Init("杂项/时间规则.txt")
	}

	for _, v := range dAtAs.text.Array {
		v.Init("文字/文本.txt")
	}

	for _, v := range dAtAs.heroLevelSubData.Array {
		v.Init("内政/君主等级.txt", dAtAs.heroLevelSubData.Map)
	}

	for _, v := range dAtAs.colorData.Array {
		v.Init("文字/品质颜色.txt")
	}

	for _, v := range dAtAs.broadcastData.Array {
		v.Init("文字/广播.txt")
	}

	dAtAs.combatMiscConfig.Init("战斗/杂项.txt", dAtAs)

	for _, v := range dAtAs.chargeObjData.Array {
		v.Init("充值/充值项.txt")
	}

	for _, v := range dAtAs.namelessCaptainData.Array {
		v.Init("武将/无名武将.txt")
	}

	for _, v := range dAtAs.captainLevelData.Array {
		v.Init("武将/武将等级.txt", dAtAs.captainLevelData.Map)
	}

	for _, v := range dAtAs.bufferTypeData.Array {
		v.Init("策略/增益类型.txt", dAtAs)
	}

	for _, v := range dAtAs.bufferData.Array {
		v.Init("策略/增益.txt", dAtAs)
	}

	for _, v := range dAtAs.blockData.Array {
		v.Init(gos, "地图阻挡/阻挡块.txt")
	}

	dAtAs.baiZhanMiscData.Init("百战千军/杂项.txt")

	dAtAs.goodsConfig.Init("物品/物品杂项.txt", dAtAs)

	for _, v := range dAtAs.cost.Array {
		v.Init("杂项/消耗.txt")
	}

	dAtAs.randomEventDataDictionary.Init("随机事件/随机事件.txt", dAtAs)

	dAtAs.questionMiscData.Init("答题/答题杂项.txt", dAtAs)

	for _, v := range dAtAs.timeLimitGiftGroupData.Array {
		v.Init("福利/定时时限礼包组.txt", dAtAs)
	}

	for _, v := range dAtAs.mingcWarSceneData.Array {
		v.Init("名城战/场景.txt", dAtAs)
	}

	dAtAs.junYingMiscData.Init("军事/军营杂项.txt")

	dAtAs.jiuGuanMiscData.Init("军事/酒馆杂项.txt")

	dAtAs.hebiMiscData.Init("合璧/合璧杂项.txt", dAtAs)

	for _, v := range dAtAs.npcGuildTemplate.Array {
		v.Init("联盟/联盟Npc模板.txt", dAtAs)
	}

	for _, v := range dAtAs.guildTarget.Array {
		v.Init("联盟/联盟目标.txt", dAtAs)
	}

	dAtAs.gemDatas.Init("物品/宝石.txt", dAtAs)

	for _, v := range dAtAs.equipmentQualityData.Array {
		v.Init("物品/装备品质.txt", dAtAs.equipmentQualityData.Map, dAtAs)
	}

	for _, v := range dAtAs.equipmentData.Array {
		v.Init("物品/装备.txt")
	}

	for _, v := range dAtAs.fishingCostData.Array {
		v.Init("钓鱼/钓鱼消耗.txt")
	}

	for _, v := range dAtAs.tieJiangPuLevelData.Array {
		v.Init("内政/铁匠铺等级.txt", dAtAs)
	}

	for _, v := range dAtAs.soldierLevelData.Array {
		v.Init("军事/士兵等级.txt")
	}

	for _, v := range dAtAs.cityEventLevelData.Array {
		v.Init("内政/城内事件等级.txt")
	}

	for _, v := range dAtAs.buildingData.Array {
		v.Init("内政/建筑.txt", dAtAs.buildingData.Map, dAtAs)
	}

	for _, v := range dAtAs.equipCombineData.Array {
		v.Init("合成/装备合成.txt")
	}

	for _, v := range dAtAs.captainStarData.Array {
		v.Init("武将/武将星数.txt")
	}

	for _, v := range dAtAs.captainRebirthLevelData.Array {
		v.Init("武将/武将转生.txt", dAtAs.captainRebirthLevelData.Map, dAtAs)
	}

	for _, v := range dAtAs.titleData.Array {
		v.Init("任务/称号.txt", dAtAs.titleData.Map)
	}

	dAtAs.regionConfig.Init("地图/地区杂项.txt", dAtAs)

	dAtAs.shopMiscData.Init("商店/商店杂项.txt")

	for _, v := range dAtAs.heroLevelData.Array {
		v.Init("内政/君主等级.txt", dAtAs.heroLevelData.Map, dAtAs)
	}

	for _, v := range dAtAs.outerCityBuildingData.Array {
		v.Init("内政/外城建筑.txt", dAtAs)
	}

	dAtAs.cityEventMiscData.Init("内政/城内事件杂项.txt", dAtAs)

	for _, v := range dAtAs.buildingUnlockData.Array {
		v.Init("内政/建筑解锁.txt", dAtAs)
	}

	dAtAs.equipCombineDatas.Init("合成/装备合成.txt", dAtAs)

	for _, v := range dAtAs.activeDegreeTaskData.Array {
		v.Init("任务/活跃度任务.txt")
	}

	for _, v := range dAtAs.prize.Array {
		v.Init("杂项/奖励.txt")
	}

	for _, v := range dAtAs.plunderPrize.Array {
		v.Init("杂项/掉落_奖励.txt", dAtAs)
	}

	for _, v := range dAtAs.plunderItem.Array {
		v.Init("杂项/掉落项.txt")
	}

	for _, v := range dAtAs.plunderGroup.Array {
		v.Init("杂项/掉落组.txt")
	}

	for _, v := range dAtAs.plunder.Array {
		v.Init("杂项/掉落.txt", dAtAs)
	}

	for _, v := range dAtAs.baowuData.Array {
		v.Init("物品/宝物.txt", dAtAs.baowuData.Map)
	}

	for _, v := range dAtAs.redPacketData.Array {
		v.Init("杂项/红包.txt")
	}

	for _, v := range dAtAs.optionPrize.Array {
		v.Init("随机事件/选项奖励.txt")
	}

	for _, v := range dAtAs.questionPrizeData.Array {
		v.Init("答题/答题奖励.txt", dAtAs)
	}

	for _, v := range dAtAs.eventLimitGiftData.Array {
		v.Init("福利/事件时限礼包.txt", dAtAs)
	}

	dAtAs.eventLimitGiftConfig.Init("福利/事件时限礼包.txt", dAtAs)

	for _, v := range dAtAs.tutorData.Array {
		v.Init("军事/酒馆导师.txt")
	}

	for _, v := range dAtAs.jiuGuanData.Array {
		v.Init("军事/酒馆.txt")
	}

	for _, v := range dAtAs.mailData.Array {
		v.Init("文字/邮件.txt")
	}

	for _, v := range dAtAs.hebiPrizeData.Array {
		v.Init("合璧/合璧奖励.txt")
	}

	for _, v := range dAtAs.guildTechnologyData.Array {
		v.Init("联盟/联盟科技.txt", dAtAs)
	}

	for _, v := range dAtAs.guildTaskData.Array {
		v.Init("联盟/联盟任务.txt")
	}

	for _, v := range dAtAs.guildPrestigePrizeData.Array {
		v.Init("联盟/联盟声望礼包.txt")
	}

	for _, v := range dAtAs.guildEventPrizeData.Array {
		v.Init(dAtAs)
	}

	for _, v := range dAtAs.guildBigBoxData.Array {
		v.Init(dAtAs)
	}

	for _, v := range dAtAs.fishData.Array {
		v.Init("钓鱼/钓鱼数据.txt")
	}

	for _, v := range dAtAs.workshopLevelData.Array {
		v.Init("内政/装备作坊.txt", dAtAs)
	}

	for _, v := range dAtAs.outerCityLayoutData.Array {
		v.Init("内政/外城布局.txt")
	}

	for _, v := range dAtAs.outerCityData.Array {
		v.Init("内政/外城.txt", dAtAs)
	}

	for _, v := range dAtAs.countdownPrizeData.Array {
		v.Init("内政/倒计时礼包.txt", dAtAs.countdownPrizeData.Map)
	}

	for _, v := range dAtAs.captainFriendshipData.Array {
		v.Init("武将/武将羁绊.txt")
	}

	for _, v := range dAtAs.captainData.Array {
		v.Init("武将/武将.txt", dAtAs)
	}

	for _, v := range dAtAs.captainAbilityData.Array {
		v.Init(dAtAs.captainAbilityData.Map)
	}

	for _, v := range dAtAs.bodyData.Array {
		v.Init("杂项/形象.txt")
	}

	for _, v := range dAtAs.homeNpcBaseData.Array {
		v.Init("地图/玩家主城野怪.txt")
	}

	for _, v := range dAtAs.junXianPrizeData.Array {
		v.Init()
	}

	for _, v := range dAtAs.zhengWuData.Array {
		v.Init("政务/政务.txt", dAtAs)
	}

	for _, v := range dAtAs.zhanJiangGuanQiaData.Array {
		v.Init()
	}

	for _, v := range dAtAs.xuanyuanRankPrizeData.Array {
		v.Init("轩辕会武/排名奖励.txt", dAtAs.xuanyuanRankPrizeData.Map, dAtAs)
	}

	for _, v := range dAtAs.taskBoxData.Array {
		v.Init("任务/任务宝箱.txt", dAtAs.taskBoxData.Map)
	}

	for _, v := range dAtAs.bwzlPrizeData.Array {
		v.Init()
	}

	for _, v := range dAtAs.baYeStageData.Array {
		v.Init("任务/霸业任务阶段.txt", dAtAs.baYeStageData.Map)
	}

	for _, v := range dAtAs.surveyData.Array {
		v.Init()
	}

	dAtAs.miscGenConfig.Init("杂项/杂项.txt")

	dAtAs.miscConfig.Init("杂项/杂项.txt")

	dAtAs.militaryConfig.Init("军事/军事杂项.txt", dAtAs)

	dAtAs.guildConfig.Init("联盟/联盟杂项.txt", dAtAs)

	for _, v := range dAtAs.zhenBaoGeShopGoods.Array {
		v.Init2("商店/珍宝阁商品.txt")
	}

	for _, v := range dAtAs.monsterCaptainData.Array {
		v.Init("怪物/怪物武将.txt")
	}

	dAtAs.heroCreateData.Init("杂项/创建英雄基础数据.txt", dAtAs)

	for _, v := range dAtAs.headData.Array {
		v.Init("杂项/头像.txt")
	}

	for _, v := range dAtAs.fishingShowData.Array {
		v.Init("钓鱼/钓鱼展示.txt")
	}

	dAtAs.fishRandomer.Init("钓鱼/钓鱼数据.txt", dAtAs)

	dAtAs.combatConfig.Init("战斗/杂项.txt", dAtAs)

	dAtAs.zhengWuMiscData.Init("政务/其他.txt", dAtAs)

	for _, v := range dAtAs.monsterMasterData.Array {
		v.Init("怪物/怪物君主.txt", dAtAs)
	}

	for _, v := range dAtAs.npcBaseData.Array {
		v.Init("地图/野怪基础数据.txt")
	}

	for _, v := range dAtAs.junXianLevelData.Array {
		v.Init()
	}

	dAtAs.zhengWuRandomData.Init(dAtAs)

	for _, v := range dAtAs.resistXiongNuData.Array {
		v.Init("抗击匈奴/难度.txt", dAtAs)
	}

	dAtAs.guildGenConfig.Init("联盟/联盟杂项.txt")

	for _, v := range dAtAs.regionMultiLevelNpcLevelData.Array {
		v.Init("地图/多等级野怪等级.txt")
	}

	for _, v := range dAtAs.regionMultiLevelNpcData.Array {
		v.Init("地图/多等级野怪.txt", dAtAs)
	}

	for _, v := range dAtAs.regionData.Array {
		v.Init2("地图/地区.txt", dAtAs)
	}

	for _, v := range dAtAs.junTuanNpcData.Array {
		v.Init("地图/军团怪物.txt")
	}

	for _, v := range dAtAs.baozNpcData.Array {
		v.Init("地图/宝藏怪物.txt")
	}

	dAtAs.randomEventPositionDictionary.Init("随机事件/事件坐标.txt", dAtAs)

	for _, v := range dAtAs.junTuanNpcPlaceData.Array {
		v.Init("地图/军团怪物刷新.txt", dAtAs)
	}

	dAtAs.junTuanNpcPlaceConfig.Init("地图/军团怪物.txt", dAtAs)

	for _, v := range dAtAs.towerData.Array {
		v.Init("系统模块/千重楼.txt", dAtAs.towerData.Map)
	}

	for _, v := range dAtAs.secretTowerData.Array {
		v.Init("系统模块/重楼密室.txt")
	}

	for _, v := range dAtAs.mainTaskData.Array {
		v.Init("任务/主线任务.txt", dAtAs.mainTaskData.Map)
	}

	for _, v := range dAtAs.branchTaskData.Array {
		v.Init("任务/支线任务.txt", dAtAs.branchTaskData.Map)
	}

	for _, v := range dAtAs.achieveTaskData.Array {
		v.Init("任务/成就任务.txt", dAtAs.achieveTaskData.Map)
	}

	for _, v := range dAtAs.dungeonData.Array {
		v.Init("推图副本/推图副本.txt", dAtAs)
	}

	for _, v := range dAtAs.dungeonChapterData.Array {
		v.Init("推图副本/推图副本章节.txt", dAtAs)
	}

	for _, v := range dAtAs.taskTargetData.Array {
		v.Init("任务/任务目标.txt", dAtAs)
	}

	for _, v := range dAtAs.functionOpenData.Array {
		v.Init("功能开启/功能开启.txt")
	}

	for _, v := range dAtAs.dungeonGuideTroopData.Array {
		v.Init("推图副本/引导布阵.txt", dAtAs)
	}

	for _, v := range dAtAs.activityShowData.Array {
		v.Init("活动/活动展示.txt")
	}

	dAtAs.heroInitData.Init("singleton/hero_init_data.txt", dAtAs)

	for _, v := range dAtAs.vipLevelData.Array {
		v.Init("vip/vip等级.txt", dAtAs)
	}

	for _, v := range dAtAs.vipContinueDaysData.Array {
		v.Init("vip/连续登录奖励.txt", dAtAs)
	}

	// Init All

	(*zhanjiang.ZhanJiangChapterData)(nil).InitAll("过关斩将/章节.txt", dAtAs)

	(*xiongnu.ResistXiongNuScoreData)(nil).InitAll("抗击匈奴/评分.txt", dAtAs)

	(*towerdata.SecretTowerData)(nil).InitAll("系统模块/重楼密室.txt", dAtAs.secretTowerData.Array)

	(*taskdata.TitleData)(nil).InitAll("任务/称号.txt", dAtAs.titleData.Map)

	(*taskdata.TaskMiscData)(nil).InitAll(dAtAs)

	(*strongerdata.StrongerData)(nil).InitAll("杂项/变强.txt", dAtAs.strongerData.Array, dAtAs)

	(*strategydata.StrategyEffectData)(nil).InitAll("策略/策略效果.txt", dAtAs)

	(*shop.Shop)(nil).InitAll("商店/商店.txt", dAtAs)

	(*resdata.PlunderPrize)(nil).InitAll("杂项/掉落_奖励.txt", dAtAs.plunderPrize.Map, dAtAs)

	(*pvetroop.PveTroopData)(nil).InitAll("杂项/pve部队.txt", dAtAs)

	(*mingcdata.MingcWarTroopLastBeatWhenFailData)(nil).InitAll("名城战/舍命一击.txt", dAtAs)

	(*mingcdata.MingcTimeData)(nil).InitAll("名城战/时间.txt", dAtAs)

	(*mingcdata.McBuildMcSupportData)(nil).InitAll("名城营建/名城民心.txt", dAtAs)

	(*mingcdata.McBuildGuildMemberPrizeData)(nil).InitAll("名城营建/联盟成员奖励.txt", dAtAs)

	(*icon.Icon)(nil).InitAll("杂项/图标.txt", dAtAs)

	(*i18n.I18nData)(nil).InitAll("i18n/语言.txt", dAtAs.i18nData.Map)

	(*guild_data.GuildTechnologyData)(nil).InitAll("联盟/联盟科技.txt", dAtAs.guildTechnologyData.Array)

	(*guild_data.GuildRankPrizeData)(nil).InitAll("联盟/联盟排行奖励.txt", dAtAs)

	(*guild_data.GuildEventPrizeData)(nil).InitAll("联盟/联盟盟友礼包.txt", dAtAs)

	(*goods.GemData)(nil).InitAll("物品/宝石.txt", dAtAs)

	(*goods.EquipmentTaozData)(nil).InitAll("物品/装备套装.txt", dAtAs.equipmentTaozData.Array, dAtAs)

	(*farm.FarmOneKeyConfig)(nil).InitAll("农场/农场一键种植.txt", dAtAs.farmOneKeyConfig.Array, dAtAs)

	(*farm.FarmMaxStealConfig)(nil).InitAll("农场/农场偷菜上限.txt", dAtAs.farmMaxStealConfig.Array, dAtAs)

	(*dungeon.DungeonData)(nil).InitAll("推图副本/推图副本.txt", dAtAs)

	(*domestic_data.WorkshopRefreshCost)(nil).InitAll("内政/装备作坊刷新消耗.txt", dAtAs)

	(*domestic_data.TechnologyData)(nil).InitAll("内政/科技.txt", dAtAs.technologyData.Array)

	(*data.BuffEffectData)(nil).InitAll("杂项/buff.txt", dAtAs)

	(*data.BroadcastData)(nil).InitAll("文字/广播.txt", dAtAs)

	(*country.CountryOfficialData)(nil).InitAll("国家/官职.txt", dAtAs)

	(*country.CountryData)(nil).InitAll("国家/国家.txt", dAtAs)

	(*captain.CaptainOfficialData)(nil).InitAll("武将/武将官职.txt", dAtAs.captainOfficialData.Array)

	(*captain.CaptainOfficialCountData)(nil).InitAll(dAtAs.captainOfficialCountData.Array, dAtAs)

	(*buffer.BufferData)(nil).InitAll("策略/增益.txt", dAtAs)

	(*bai_zhan_data.JunXianLevelData)(nil).InitAll("百战千军/等级.txt", dAtAs)

	(*activitydata.ActivityTaskListModeData)(nil).InitAll(dAtAs)

	(*xiongnu.ResistXiongNuData)(nil).InitAll("抗击匈奴/难度.txt", dAtAs)

	(*vip.VipLevelData)(nil).InitAll("vip/vip等级.txt", dAtAs)

	(*teach.TeachChapterData)(nil).InitAll("教学/关卡.txt", dAtAs)

	(*goods.GoodsCheck)(nil).InitAll(dAtAs)

	(*domestic_data.ProsperityDamageBuffData)(nil).InitAll("内政/繁荣度buff.txt", dAtAs)

	(*bai_zhan_data.JunXianPrizeData)(nil).InitAll("百战千军/等级奖励.txt", dAtAs.junXianPrizeData.Array)

	// clear parser

	dAtAs.activityCollectionData.parserMap = nil

	dAtAs.activityShowData.parserMap = nil

	dAtAs.activityTaskListModeData.parserMap = nil

	dAtAs.collectionExchangeData.parserMap = nil

	dAtAs.junXianLevelData.parserMap = nil

	dAtAs.junXianPrizeData.parserMap = nil

	dAtAs.homeNpcBaseData.parserMap = nil

	dAtAs.npcBaseData.parserMap = nil

	dAtAs.blockData.parserMap = nil

	dAtAs.bodyData.parserMap = nil

	dAtAs.bufferData.parserMap = nil

	dAtAs.bufferTypeData.parserMap = nil

	dAtAs.captainAbilityData.parserMap = nil

	dAtAs.captainData.parserMap = nil

	dAtAs.captainFriendshipData.parserMap = nil

	dAtAs.captainLevelData.parserMap = nil

	dAtAs.captainOfficialCountData.parserMap = nil

	dAtAs.captainOfficialData.parserMap = nil

	dAtAs.captainRarityData.parserMap = nil

	dAtAs.captainRebirthLevelData.parserMap = nil

	dAtAs.captainStarData.parserMap = nil

	dAtAs.namelessCaptainData.parserMap = nil

	dAtAs.chargeObjData.parserMap = nil

	dAtAs.chargePrizeData.parserMap = nil

	dAtAs.productData.parserMap = nil

	dAtAs.equipCombineData.parserMap = nil

	dAtAs.goodsCombineData.parserMap = nil

	dAtAs.countryData.parserMap = nil

	dAtAs.countryOfficialData.parserMap = nil

	dAtAs.countryOfficialNpcData.parserMap = nil

	dAtAs.familyNameData.parserMap = nil

	dAtAs.broadcastData.parserMap = nil

	dAtAs.buffEffectData.parserMap = nil

	dAtAs.colorData.parserMap = nil

	dAtAs.familyName.parserMap = nil

	dAtAs.femaleGivenName.parserMap = nil

	dAtAs.heroLevelSubData.parserMap = nil

	dAtAs.maleGivenName.parserMap = nil

	dAtAs.spriteStat.parserMap = nil

	dAtAs.text.parserMap = nil

	dAtAs.timeRuleData.parserMap = nil

	dAtAs.baseLevelData.parserMap = nil

	dAtAs.buildingData.parserMap = nil

	dAtAs.buildingLayoutData.parserMap = nil

	dAtAs.buildingUnlockData.parserMap = nil

	dAtAs.cityEventData.parserMap = nil

	dAtAs.cityEventLevelData.parserMap = nil

	dAtAs.combineCost.parserMap = nil

	dAtAs.countdownPrizeData.parserMap = nil

	dAtAs.countdownPrizeDescData.parserMap = nil

	dAtAs.guanFuLevelData.parserMap = nil

	dAtAs.outerCityBuildingData.parserMap = nil

	dAtAs.outerCityBuildingDescData.parserMap = nil

	dAtAs.outerCityData.parserMap = nil

	dAtAs.outerCityLayoutData.parserMap = nil

	dAtAs.prosperityDamageBuffData.parserMap = nil

	dAtAs.soldierLevelData.parserMap = nil

	dAtAs.technologyData.parserMap = nil

	dAtAs.tieJiangPuLevelData.parserMap = nil

	dAtAs.workshopDuration.parserMap = nil

	dAtAs.workshopLevelData.parserMap = nil

	dAtAs.workshopRefreshCost.parserMap = nil

	dAtAs.dungeonChapterData.parserMap = nil

	dAtAs.dungeonData.parserMap = nil

	dAtAs.dungeonGuideTroopData.parserMap = nil

	dAtAs.farmMaxStealConfig.parserMap = nil

	dAtAs.farmOneKeyConfig.parserMap = nil

	dAtAs.farmResConfig.parserMap = nil

	dAtAs.fishData.parserMap = nil

	dAtAs.fishingCaptainProbabilityData.parserMap = nil

	dAtAs.fishingCostData.parserMap = nil

	dAtAs.fishingShowData.parserMap = nil

	dAtAs.functionOpenData.parserMap = nil

	dAtAs.treasuryTreeData.parserMap = nil

	dAtAs.equipmentData.parserMap = nil

	dAtAs.equipmentLevelData.parserMap = nil

	dAtAs.equipmentQualityData.parserMap = nil

	dAtAs.equipmentRefinedData.parserMap = nil

	dAtAs.equipmentTaozData.parserMap = nil

	dAtAs.gemData.parserMap = nil

	dAtAs.goodsData.parserMap = nil

	dAtAs.goodsQuality.parserMap = nil

	dAtAs.guildBigBoxData.parserMap = nil

	dAtAs.guildClassLevelData.parserMap = nil

	dAtAs.guildClassTitleData.parserMap = nil

	dAtAs.guildDonateData.parserMap = nil

	dAtAs.guildEventPrizeData.parserMap = nil

	dAtAs.guildLevelCdrData.parserMap = nil

	dAtAs.guildLevelData.parserMap = nil

	dAtAs.guildLogData.parserMap = nil

	dAtAs.guildPermissionShowData.parserMap = nil

	dAtAs.guildPrestigeEventData.parserMap = nil

	dAtAs.guildPrestigePrizeData.parserMap = nil

	dAtAs.guildRankPrizeData.parserMap = nil

	dAtAs.guildTarget.parserMap = nil

	dAtAs.guildTaskData.parserMap = nil

	dAtAs.guildTaskEvaluateData.parserMap = nil

	dAtAs.guildTechnologyData.parserMap = nil

	dAtAs.npcGuildTemplate.parserMap = nil

	dAtAs.npcMemberData.parserMap = nil

	dAtAs.headData.parserMap = nil

	dAtAs.hebiPrizeData.parserMap = nil

	dAtAs.heroLevelData.parserMap = nil

	dAtAs.i18nData.parserMap = nil

	dAtAs.icon.parserMap = nil

	dAtAs.locationData.parserMap = nil

	dAtAs.mailData.parserMap = nil

	dAtAs.jiuGuanData.parserMap = nil

	dAtAs.junYingLevelData.parserMap = nil

	dAtAs.trainingLevelData.parserMap = nil

	dAtAs.tutorData.parserMap = nil

	dAtAs.mcBuildAddSupportData.parserMap = nil

	dAtAs.mcBuildGuildMemberPrizeData.parserMap = nil

	dAtAs.mcBuildMcSupportData.parserMap = nil

	dAtAs.mingcBaseData.parserMap = nil

	dAtAs.mingcTimeData.parserMap = nil

	dAtAs.mingcWarBuildingData.parserMap = nil

	dAtAs.mingcWarDrumStatData.parserMap = nil

	dAtAs.mingcWarMapData.parserMap = nil

	dAtAs.mingcWarMultiKillData.parserMap = nil

	dAtAs.mingcWarNpcData.parserMap = nil

	dAtAs.mingcWarNpcGuildData.parserMap = nil

	dAtAs.mingcWarSceneData.parserMap = nil

	dAtAs.mingcWarTouShiBuildingTargetData.parserMap = nil

	dAtAs.mingcWarTroopLastBeatWhenFailData.parserMap = nil

	dAtAs.monsterCaptainData.parserMap = nil

	dAtAs.monsterMasterData.parserMap = nil

	dAtAs.dailyBargainData.parserMap = nil

	dAtAs.durationCardData.parserMap = nil

	dAtAs.eventLimitGiftData.parserMap = nil

	dAtAs.freeGiftData.parserMap = nil

	dAtAs.heroLevelFundData.parserMap = nil

	dAtAs.loginDayData.parserMap = nil

	dAtAs.spCollectionData.parserMap = nil

	dAtAs.timeLimitGiftData.parserMap = nil

	dAtAs.timeLimitGiftGroupData.parserMap = nil

	dAtAs.pushData.parserMap = nil

	dAtAs.pveTroopData.parserMap = nil

	dAtAs.questionData.parserMap = nil

	dAtAs.questionPrizeData.parserMap = nil

	dAtAs.questionSayingData.parserMap = nil

	dAtAs.raceData.parserMap = nil

	dAtAs.eventOptionData.parserMap = nil

	dAtAs.eventPosition.parserMap = nil

	dAtAs.optionPrize.parserMap = nil

	dAtAs.randomEventData.parserMap = nil

	dAtAs.redPacketData.parserMap = nil

	dAtAs.areaData.parserMap = nil

	dAtAs.assemblyData.parserMap = nil

	dAtAs.baozNpcData.parserMap = nil

	dAtAs.junTuanNpcData.parserMap = nil

	dAtAs.junTuanNpcPlaceData.parserMap = nil

	dAtAs.regionAreaData.parserMap = nil

	dAtAs.regionData.parserMap = nil

	dAtAs.regionMonsterData.parserMap = nil

	dAtAs.regionMultiLevelNpcData.parserMap = nil

	dAtAs.regionMultiLevelNpcLevelData.parserMap = nil

	dAtAs.regionMultiLevelNpcTypeData.parserMap = nil

	dAtAs.troopDialogueData.parserMap = nil

	dAtAs.troopDialogueTextData.parserMap = nil

	dAtAs.amountShowSortData.parserMap = nil

	dAtAs.baowuData.parserMap = nil

	dAtAs.conditionPlunder.parserMap = nil

	dAtAs.conditionPlunderItem.parserMap = nil

	dAtAs.cost.parserMap = nil

	dAtAs.guildLevelPrize.parserMap = nil

	dAtAs.plunder.parserMap = nil

	dAtAs.plunderGroup.parserMap = nil

	dAtAs.plunderItem.parserMap = nil

	dAtAs.plunderPrize.parserMap = nil

	dAtAs.prize.parserMap = nil

	dAtAs.resCaptainData.parserMap = nil

	dAtAs.combatScene.parserMap = nil

	dAtAs.seasonData.parserMap = nil

	dAtAs.privacySettingData.parserMap = nil

	dAtAs.blackMarketData.parserMap = nil

	dAtAs.blackMarketGoodsData.parserMap = nil

	dAtAs.blackMarketGoodsGroupData.parserMap = nil

	dAtAs.discountColorData.parserMap = nil

	dAtAs.normalShopGoods.parserMap = nil

	dAtAs.shop.parserMap = nil

	dAtAs.zhenBaoGeShopGoods.parserMap = nil

	dAtAs.passiveSpellData.parserMap = nil

	dAtAs.spell.parserMap = nil

	dAtAs.spellData.parserMap = nil

	dAtAs.spellFacadeData.parserMap = nil

	dAtAs.stateData.parserMap = nil

	dAtAs.strategyData.parserMap = nil

	dAtAs.strategyEffectData.parserMap = nil

	dAtAs.strongerData.parserMap = nil

	dAtAs.buildingEffectData.parserMap = nil

	dAtAs.surveyData.parserMap = nil

	dAtAs.achieveTaskData.parserMap = nil

	dAtAs.achieveTaskStarPrizeData.parserMap = nil

	dAtAs.activeDegreePrizeData.parserMap = nil

	dAtAs.activeDegreeTaskData.parserMap = nil

	dAtAs.activityTaskData.parserMap = nil

	dAtAs.baYeStageData.parserMap = nil

	dAtAs.baYeTaskData.parserMap = nil

	dAtAs.branchTaskData.parserMap = nil

	dAtAs.bwzlPrizeData.parserMap = nil

	dAtAs.bwzlTaskData.parserMap = nil

	dAtAs.mainTaskData.parserMap = nil

	dAtAs.taskBoxData.parserMap = nil

	dAtAs.taskTargetData.parserMap = nil

	dAtAs.titleData.parserMap = nil

	dAtAs.titleTaskData.parserMap = nil

	dAtAs.teachChapterData.parserMap = nil

	dAtAs.secretTowerData.parserMap = nil

	dAtAs.secretTowerWordsData.parserMap = nil

	dAtAs.towerData.parserMap = nil

	dAtAs.vipContinueDaysData.parserMap = nil

	dAtAs.vipLevelData.parserMap = nil

	dAtAs.resistXiongNuData.parserMap = nil

	dAtAs.resistXiongNuScoreData.parserMap = nil

	dAtAs.resistXiongNuWaveData.parserMap = nil

	dAtAs.xuanyuanRangeData.parserMap = nil

	dAtAs.xuanyuanRankPrizeData.parserMap = nil

	dAtAs.zhanJiangChapterData.parserMap = nil

	dAtAs.zhanJiangData.parserMap = nil

	dAtAs.zhanJiangGuanQiaData.parserMap = nil

	dAtAs.zhengWuCompleteData.parserMap = nil

	dAtAs.zhengWuData.parserMap = nil

	dAtAs.zhengWuRefreshData.parserMap = nil

	dAtAs.baiZhanMiscDataParser = nil

	dAtAs.combatConfigParser = nil

	dAtAs.combatMiscConfigParser = nil

	dAtAs.equipCombineDatasParser = nil

	dAtAs.countryMiscDataParser = nil

	dAtAs.broadcastHelpParser = nil

	dAtAs.textHelpParser = nil

	dAtAs.exchangeMiscDataParser = nil

	dAtAs.buildingLayoutMiscDataParser = nil

	dAtAs.cityEventMiscDataParser = nil

	dAtAs.mainCityMiscDataParser = nil

	dAtAs.dungeonMiscDataParser = nil

	dAtAs.farmMiscConfigParser = nil

	dAtAs.fishRandomerParser = nil

	dAtAs.gardenConfigParser = nil

	dAtAs.equipmentTaozConfigParser = nil

	dAtAs.gemDatasParser = nil

	dAtAs.goodsCheckParser = nil

	dAtAs.guildLogHelpParser = nil

	dAtAs.npcGuildSuffixNameParser = nil

	dAtAs.hebiMiscDataParser = nil

	dAtAs.heroCreateDataParser = nil

	dAtAs.heroInitDataParser = nil

	dAtAs.mailHelpParser = nil

	dAtAs.jiuGuanMiscDataParser = nil

	dAtAs.junYingMiscDataParser = nil

	dAtAs.mcBuildMiscDataParser = nil

	dAtAs.mingcMiscDataParser = nil

	dAtAs.eventLimitGiftConfigParser = nil

	dAtAs.promotionMiscDataParser = nil

	dAtAs.questionMiscDataParser = nil

	dAtAs.raceConfigParser = nil

	dAtAs.randomEventDataDictionaryParser = nil

	dAtAs.randomEventPositionDictionaryParser = nil

	dAtAs.rankMiscDataParser = nil

	dAtAs.junTuanNpcPlaceConfigParser = nil

	dAtAs.seasonMiscDataParser = nil

	dAtAs.settingMiscDataParser = nil

	dAtAs.shopMiscDataParser = nil

	dAtAs.goodsConfigParser = nil

	dAtAs.guildConfigParser = nil

	dAtAs.guildGenConfigParser = nil

	dAtAs.militaryConfigParser = nil

	dAtAs.miscConfigParser = nil

	dAtAs.miscGenConfigParser = nil

	dAtAs.regionConfigParser = nil

	dAtAs.regionGenConfigParser = nil

	dAtAs.tagMiscDataParser = nil

	dAtAs.taskMiscDataParser = nil

	dAtAs.secretTowerMiscDataParser = nil

	dAtAs.vipMiscDataParser = nil

	dAtAs.resistXiongNuMiscParser = nil

	dAtAs.xuanyuanMiscDataParser = nil

	dAtAs.zhanJiangMiscDataParser = nil

	dAtAs.zhengWuMiscDataParser = nil

	dAtAs.zhengWuRandomDataParser = nil

	return dAtAs, nil
}

func (dAtA *ConfigDatas) EncodeClient() *shared_proto.Config {

	c := &shared_proto.Config{}
	c.Gen = &shared_proto.ConfigGen{}

	for _, v := range dAtA.junXianLevelData.Array {
		c.JunXianLevel = append(c.JunXianLevel, v.Encode())
	}

	for _, v := range dAtA.junXianPrizeData.Array {
		c.JunXianLevelPrize = append(c.JunXianLevelPrize, v.Encode())
	}

	for _, v := range dAtA.npcBaseData.Array {
		c.NpcBaseData = append(c.NpcBaseData, v.Encode())
	}

	for _, v := range dAtA.bodyData.Array {
		c.Bodys = append(c.Bodys, v.Encode())
	}

	for _, v := range dAtA.bufferData.Array {
		c.Gen.BufferData = append(c.Gen.BufferData, v.Encode())
	}

	for _, v := range dAtA.bufferTypeData.Array {
		c.Gen.BufferTypeData = append(c.Gen.BufferTypeData, v.Encode())
	}

	for _, v := range dAtA.captainAbilityData.Array {
		c.CaptainAbility = append(c.CaptainAbility, v.Encode())
	}

	for _, v := range dAtA.captainData.Array {
		c.Gen.CaptainData = append(c.Gen.CaptainData, v.Encode())
	}

	for _, v := range dAtA.captainFriendshipData.Array {
		c.Gen.CaptainFriendshipData = append(c.Gen.CaptainFriendshipData, v.Encode())
	}

	for _, v := range dAtA.captainLevelData.Array {
		c.CaptainLevel = append(c.CaptainLevel, v.Encode())
	}

	for _, v := range dAtA.captainOfficialData.Array {
		c.CaptainOfficial = append(c.CaptainOfficial, v.Encode())
	}

	for _, v := range dAtA.captainRarityData.Array {
		c.Gen.CaptainRarityData = append(c.Gen.CaptainRarityData, v.Encode())
	}

	for _, v := range dAtA.captainRebirthLevelData.Array {
		c.CaptainRebirthLevel = append(c.CaptainRebirthLevel, v.Encode())
	}

	for _, v := range dAtA.namelessCaptainData.Array {
		c.Gen.NamelessCaptainData = append(c.Gen.NamelessCaptainData, v.Encode())
	}

	for _, v := range dAtA.chargeObjData.Array {
		c.Gen.ChargeObjData = append(c.Gen.ChargeObjData, v.Encode())
	}

	for _, v := range dAtA.chargePrizeData.Array {
		c.Gen.ChargePrizeData = append(c.Gen.ChargePrizeData, v.Encode())
	}

	for _, v := range dAtA.equipCombineData.Array {
		c.EquipCombine = append(c.EquipCombine, v.Encode())
	}

	for _, v := range dAtA.countryData.Array {
		c.Gen.CountryData = append(c.Gen.CountryData, v.Encode())
	}

	for _, v := range dAtA.countryOfficialData.Array {
		c.Gen.CountryOfficialData = append(c.Gen.CountryOfficialData, v.Encode())
	}

	for _, v := range dAtA.countryOfficialNpcData.Array {
		c.Gen.CountryOfficialNpcData = append(c.Gen.CountryOfficialNpcData, v.Encode())
	}

	for _, v := range dAtA.familyNameData.Array {
		c.Gen.FamilyNameData = append(c.Gen.FamilyNameData, v.Encode())
	}

	for _, v := range dAtA.buffEffectData.Array {
		c.Gen.BuffEffectData = append(c.Gen.BuffEffectData, v.Encode())
	}

	for _, v := range dAtA.heroLevelSubData.Array {
		c.Hero = append(c.Hero, v.Encode())
	}

	for _, v := range dAtA.baseLevelData.Array {
		c.BaseLevel = append(c.BaseLevel, v.Encode())
	}

	for _, v := range dAtA.buildingData.Array {
		c.BuildingData = append(c.BuildingData, v.Encode())
	}

	for _, v := range dAtA.buildingLayoutData.Array {
		c.BuildingLayout = append(c.BuildingLayout, v.Encode())
	}

	for _, v := range dAtA.buildingUnlockData.Array {
		c.BuildingUnlockData = append(c.BuildingUnlockData, v.Encode())
	}

	for _, v := range dAtA.cityEventData.Array {
		c.CityEventData = append(c.CityEventData, v.Encode())
	}

	for _, v := range dAtA.countdownPrizeDescData.Array {
		c.CountdownPrizeDesc = append(c.CountdownPrizeDesc, v.Encode())
	}

	for _, v := range dAtA.guanFuLevelData.Array {
		c.GuanFuLevelData = append(c.GuanFuLevelData, v.Encode())
	}

	for _, v := range dAtA.outerCityBuildingData.Array {
		c.Gen.OuterCityBuildingData = append(c.Gen.OuterCityBuildingData, v.Encode())
	}

	for _, v := range dAtA.outerCityData.Array {
		c.OuterCityDatas = append(c.OuterCityDatas, v.Encode())
	}

	for _, v := range dAtA.outerCityLayoutData.Array {
		c.OuterCityLayoutDatas = append(c.OuterCityLayoutDatas, v.Encode())
	}

	for _, v := range dAtA.prosperityDamageBuffData.Array {
		c.Gen.ProsperityDamageBuffData = append(c.Gen.ProsperityDamageBuffData, v.Encode())
	}

	for _, v := range dAtA.soldierLevelData.Array {
		c.Soldier = append(c.Soldier, v.Encode())
	}

	for _, v := range dAtA.technologyData.Array {
		c.TechnologyData = append(c.TechnologyData, v.Encode())
	}

	for _, v := range dAtA.tieJiangPuLevelData.Array {
		c.TieJiangPuLevel = append(c.TieJiangPuLevel, v.Encode())
	}

	for _, v := range dAtA.workshopRefreshCost.Array {
		c.WorkshopRefreshCosts = append(c.WorkshopRefreshCosts, v.Encode())
	}

	for _, v := range dAtA.dungeonChapterData.Array {
		c.DungeonChapter = append(c.DungeonChapter, v.Encode())
	}

	for _, v := range dAtA.farmResConfig.Array {
		c.FarmResConfig = append(c.FarmResConfig, v.Encode())
	}

	for _, v := range dAtA.fishingCaptainProbabilityData.Array {
		c.Gen.FishingCaptainProbabilityData = append(c.Gen.FishingCaptainProbabilityData, v.Encode())
	}

	for _, v := range dAtA.fishingCostData.Array {
		c.FishingCost = append(c.FishingCost, v.Encode())
	}

	for _, v := range dAtA.fishingShowData.Array {
		c.FishingShow = append(c.FishingShow, v.Encode())
	}

	for _, v := range dAtA.functionOpenData.Array {
		c.FunctionOpenDatas = append(c.FunctionOpenDatas, v.Encode())
	}

	for _, v := range dAtA.treasuryTreeData.Array {
		c.TreasuryTreeData = append(c.TreasuryTreeData, v.Encode())
	}

	for _, v := range dAtA.equipmentData.Array {
		c.Equipment = append(c.Equipment, v.Encode())
	}

	for _, v := range dAtA.equipmentQualityData.Array {
		c.EquipmentQuality = append(c.EquipmentQuality, v.Encode())
	}

	for _, v := range dAtA.equipmentRefinedData.Array {
		c.EquipmentRefined = append(c.EquipmentRefined, v.Encode())
	}

	for _, v := range dAtA.equipmentTaozData.Array {
		c.EquipmentTaoz = append(c.EquipmentTaoz, v.Encode())
	}

	for _, v := range dAtA.gemData.Array {
		c.Gem = append(c.Gem, v.Encode())
	}

	for _, v := range dAtA.goodsData.Array {
		c.Goods = append(c.Goods, v.Encode())
	}

	for _, v := range dAtA.goodsQuality.Array {
		c.GoodsQuality = append(c.GoodsQuality, v.Encode())
	}

	for _, v := range dAtA.guildBigBoxData.Array {
		c.GuildBigBox = append(c.GuildBigBox, v.Encode())
	}

	for _, v := range dAtA.guildClassLevelData.Array {
		c.GuildClassLevel = append(c.GuildClassLevel, v.Encode())
	}

	for _, v := range dAtA.guildClassTitleData.Array {
		c.GuildClassTitle = append(c.GuildClassTitle, v.Encode())
	}

	for _, v := range dAtA.guildDonateData.Array {
		c.GuildDonate = append(c.GuildDonate, v.Encode())
	}

	for _, v := range dAtA.guildEventPrizeData.Array {
		c.GuildEventPrize = append(c.GuildEventPrize, v.Encode())
	}

	for _, v := range dAtA.guildLevelData.Array {
		c.GuildLevel = append(c.GuildLevel, v.Encode())
	}

	for _, v := range dAtA.guildPermissionShowData.Array {
		c.GuildPermissionShow = append(c.GuildPermissionShow, v.Encode())
	}

	for _, v := range dAtA.guildPrestigePrizeData.Array {
		c.GuildPrestigePrize = append(c.GuildPrestigePrize, v.Encode())
	}

	for _, v := range dAtA.guildRankPrizeData.Array {
		c.Gen.GuildRankPrizeData = append(c.Gen.GuildRankPrizeData, v.Encode())
	}

	for _, v := range dAtA.guildTarget.Array {
		c.GuildTarget = append(c.GuildTarget, v.Encode())
	}

	for _, v := range dAtA.guildTaskData.Array {
		c.GuildTask = append(c.GuildTask, v.Encode())
	}

	for _, v := range dAtA.guildTaskEvaluateData.Array {
		c.GuildTaskEvaluate = append(c.GuildTaskEvaluate, v.Encode())
	}

	for _, v := range dAtA.guildTechnologyData.Array {
		c.GuildTechnology = append(c.GuildTechnology, v.Encode())
	}

	for _, v := range dAtA.headData.Array {
		c.Heads = append(c.Heads, v.Encode())
	}

	for _, v := range dAtA.hebiPrizeData.Array {
		c.HebiPrize = append(c.HebiPrize, v.Encode())
	}

	for _, v := range dAtA.i18nData.Array {
		c.I18N = append(c.I18N, v.Encode())
	}

	for _, v := range dAtA.icon.Array {
		c.Icons = append(c.Icons, v.Encode())
	}

	for _, v := range dAtA.locationData.Array {
		c.Gen.LocationData = append(c.Gen.LocationData, v.Encode())
	}

	for _, v := range dAtA.jiuGuanData.Array {
		c.JiuGuanData = append(c.JiuGuanData, v.Encode())
	}

	for _, v := range dAtA.junYingLevelData.Array {
		c.JunYingLevelData = append(c.JunYingLevelData, v.Encode())
	}

	for _, v := range dAtA.trainingLevelData.Array {
		c.TrainingLevel = append(c.TrainingLevel, v.Encode())
	}

	for _, v := range dAtA.mcBuildAddSupportData.Array {
		c.Gen.McBuildAddSupportData = append(c.Gen.McBuildAddSupportData, v.Encode())
	}

	for _, v := range dAtA.mcBuildGuildMemberPrizeData.Array {
		c.Gen.McBuildGuildMemberPrizeData = append(c.Gen.McBuildGuildMemberPrizeData, v.Encode())
	}

	for _, v := range dAtA.mcBuildMcSupportData.Array {
		c.Gen.McBuildMcSupportData = append(c.Gen.McBuildMcSupportData, v.Encode())
	}

	for _, v := range dAtA.mingcBaseData.Array {
		c.Gen.MingcBaseData = append(c.Gen.MingcBaseData, v.Encode())
	}

	for _, v := range dAtA.mingcTimeData.Array {
		c.Gen.MingcTimeData = append(c.Gen.MingcTimeData, v.Encode())
	}

	for _, v := range dAtA.mingcWarBuildingData.Array {
		c.Gen.MingcWarBuildingData = append(c.Gen.MingcWarBuildingData, v.Encode())
	}

	for _, v := range dAtA.mingcWarMapData.Array {
		c.Gen.MingcWarMapData = append(c.Gen.MingcWarMapData, v.Encode())
	}

	for _, v := range dAtA.mingcWarSceneData.Array {
		c.Gen.MingcWarSceneData = append(c.Gen.MingcWarSceneData, v.Encode())
	}

	for _, v := range dAtA.mingcWarTouShiBuildingTargetData.Array {
		c.Gen.MingcWarTouShiBuildingTargetData = append(c.Gen.MingcWarTouShiBuildingTargetData, v.Encode())
	}

	for _, v := range dAtA.mingcWarTroopLastBeatWhenFailData.Array {
		c.Gen.MingcWarTroopLastBeatWhenFailData = append(c.Gen.MingcWarTroopLastBeatWhenFailData, v.Encode())
	}

	for _, v := range dAtA.dailyBargainData.Array {
		c.Gen.DailyBargainData = append(c.Gen.DailyBargainData, v.Encode())
	}

	for _, v := range dAtA.durationCardData.Array {
		c.Gen.DurationCardData = append(c.Gen.DurationCardData, v.Encode())
	}

	for _, v := range dAtA.eventLimitGiftData.Array {
		c.Gen.EventLimitGiftData = append(c.Gen.EventLimitGiftData, v.Encode())
	}

	for _, v := range dAtA.freeGiftData.Array {
		c.Gen.FreeGiftData = append(c.Gen.FreeGiftData, v.Encode())
	}

	for _, v := range dAtA.heroLevelFundData.Array {
		c.Gen.HeroLevelFundData = append(c.Gen.HeroLevelFundData, v.Encode())
	}

	for _, v := range dAtA.loginDayData.Array {
		c.Gen.LoginDayData = append(c.Gen.LoginDayData, v.Encode())
	}

	for _, v := range dAtA.spCollectionData.Array {
		c.Gen.SpCollectionData = append(c.Gen.SpCollectionData, v.Encode())
	}

	for _, v := range dAtA.timeLimitGiftData.Array {
		c.Gen.TimeLimitGiftData = append(c.Gen.TimeLimitGiftData, v.Encode())
	}

	for _, v := range dAtA.timeLimitGiftGroupData.Array {
		c.TimeLimitGiftGroups = append(c.TimeLimitGiftGroups, v.Encode())
	}

	for _, v := range dAtA.pushData.Array {
		c.Gen.PushData = append(c.Gen.PushData, v.Encode())
	}

	for _, v := range dAtA.pveTroopData.Array {
		c.PveTroopDatas = append(c.PveTroopDatas, v.Encode())
	}

	for _, v := range dAtA.questionData.Array {
		c.Question = append(c.Question, v.Encode())
	}

	for _, v := range dAtA.questionPrizeData.Array {
		c.QuestionPrize = append(c.QuestionPrize, v.Encode())
	}

	for _, v := range dAtA.questionSayingData.Array {
		c.QuestionSaying = append(c.QuestionSaying, v.Encode())
	}

	for _, v := range dAtA.raceData.Array {
		c.RaceData = append(c.RaceData, v.Encode())
	}

	for _, v := range dAtA.randomEventData.Array {
		c.RandomEvent = append(c.RandomEvent, v.Encode())
	}

	for _, v := range dAtA.redPacketData.Array {
		c.Gen.RedPacketData = append(c.Gen.RedPacketData, v.Encode())
	}

	for _, v := range dAtA.assemblyData.Array {
		c.Gen.AssemblyData = append(c.Gen.AssemblyData, v.Encode())
	}

	for _, v := range dAtA.baozNpcData.Array {
		c.Gen.BaozNpcData = append(c.Gen.BaozNpcData, v.Encode())
	}

	for _, v := range dAtA.junTuanNpcData.Array {
		c.Gen.JunTuanNpcData = append(c.Gen.JunTuanNpcData, v.Encode())
	}

	for _, v := range dAtA.regionAreaData.Array {
		c.Gen.RegionAreaData = append(c.Gen.RegionAreaData, v.Encode())
	}

	for _, v := range dAtA.regionData.Array {
		c.RegionData = append(c.RegionData, v.Encode())
	}

	for _, v := range dAtA.regionMultiLevelNpcData.Array {
		c.MultiLevelNpcData = append(c.MultiLevelNpcData, v.Encode())
	}

	for _, v := range dAtA.regionMultiLevelNpcTypeData.Array {
		c.MultiLevelNpcType = append(c.MultiLevelNpcType, v.Encode())
	}

	for _, v := range dAtA.troopDialogueData.Array {
		c.Gen.TroopDialogueData = append(c.Gen.TroopDialogueData, v.Encode())
	}

	for _, v := range dAtA.troopDialogueTextData.Array {
		c.Gen.TroopDialogueTextData = append(c.Gen.TroopDialogueTextData, v.Encode())
	}

	for _, v := range dAtA.amountShowSortData.Array {
		c.AmountShowSortDatas = append(c.AmountShowSortDatas, v.Encode())
	}

	for _, v := range dAtA.baowuData.Array {
		c.Gen.BaowuData = append(c.Gen.BaowuData, v.Encode())
	}

	for _, v := range dAtA.seasonData.Array {
		c.SeasonDatas = append(c.SeasonDatas, v.Encode())
	}

	for _, v := range dAtA.privacySettingData.Array {
		c.Gen.PrivacySettingData = append(c.Gen.PrivacySettingData, v.Encode())
	}

	for _, v := range dAtA.blackMarketGoodsData.Array {
		c.Gen.BlackMarketGoodsData = append(c.Gen.BlackMarketGoodsData, v.Encode())
	}

	for _, v := range dAtA.discountColorData.Array {
		c.Gen.DiscountColorData = append(c.Gen.DiscountColorData, v.Encode())
	}

	for _, v := range dAtA.shop.Array {
		c.Shop = append(c.Shop, v.Encode())
	}

	for _, v := range dAtA.spell.Array {
		c.SpellConfig = append(c.SpellConfig, v.Encode())
	}

	for _, v := range dAtA.spellFacadeData.Array {
		c.Gen.SpellFacadeData = append(c.Gen.SpellFacadeData, v.Encode())
	}

	for _, v := range dAtA.strategyData.Array {
		c.Strategy = append(c.Strategy, v.Encode())
	}

	for _, v := range dAtA.strategyEffectData.Array {
		c.Gen.StrategyEffectData = append(c.Gen.StrategyEffectData, v.Encode())
	}

	for _, v := range dAtA.strongerData.Array {
		c.StrongerData = append(c.StrongerData, v.Encode())
	}

	for _, v := range dAtA.surveyData.Array {
		c.SurveyDatas = append(c.SurveyDatas, v.Encode())
	}

	for _, v := range dAtA.achieveTaskData.Array {
		c.AchieveTask = append(c.AchieveTask, v.Encode())
	}

	for _, v := range dAtA.achieveTaskStarPrizeData.Array {
		c.AchieveTaskStarPrize = append(c.AchieveTaskStarPrize, v.Encode())
	}

	for _, v := range dAtA.activeDegreePrizeData.Array {
		c.ActiveDegreePrize = append(c.ActiveDegreePrize, v.Encode())
	}

	for _, v := range dAtA.activeDegreeTaskData.Array {
		c.ActiveDegreeTask = append(c.ActiveDegreeTask, v.Encode())
	}

	for _, v := range dAtA.baYeStageData.Array {
		c.BaYeStage = append(c.BaYeStage, v.Encode())
	}

	for _, v := range dAtA.baYeTaskData.Array {
		c.BaYeTask = append(c.BaYeTask, v.Encode())
	}

	for _, v := range dAtA.branchTaskData.Array {
		c.BranchTask = append(c.BranchTask, v.Encode())
	}

	for _, v := range dAtA.bwzlPrizeData.Array {
		c.BwzlPrize = append(c.BwzlPrize, v.Encode())
	}

	for _, v := range dAtA.bwzlTaskData.Array {
		c.BwzlTask = append(c.BwzlTask, v.Encode())
	}

	for _, v := range dAtA.mainTaskData.Array {
		c.MainTask = append(c.MainTask, v.Encode())
	}

	for _, v := range dAtA.taskBoxData.Array {
		c.TaskBox = append(c.TaskBox, v.Encode())
	}

	for _, v := range dAtA.titleData.Array {
		c.Gen.TitleData = append(c.Gen.TitleData, v.Encode())
	}

	for _, v := range dAtA.titleTaskData.Array {
		c.Gen.TitleTaskData = append(c.Gen.TitleTaskData, v.Encode())
	}

	for _, v := range dAtA.teachChapterData.Array {
		c.TeachData = append(c.TeachData, v.Encode())
	}

	for _, v := range dAtA.secretTowerData.Array {
		c.SecretTower = append(c.SecretTower, v.Encode())
	}

	for _, v := range dAtA.secretTowerWordsData.Array {
		c.Gen.SecretTowerWordsData = append(c.Gen.SecretTowerWordsData, v.Encode())
	}

	for _, v := range dAtA.towerData.Array {
		c.Tower = append(c.Tower, v.Encode())
	}

	for _, v := range dAtA.vipContinueDaysData.Array {
		c.Gen.VipContinueDaysData = append(c.Gen.VipContinueDaysData, v.Encode())
	}

	for _, v := range dAtA.vipLevelData.Array {
		c.Gen.VipLevelData = append(c.Gen.VipLevelData, v.Encode())
	}

	for _, v := range dAtA.resistXiongNuData.Array {
		c.ResistXiongNuData = append(c.ResistXiongNuData, v.Encode())
	}

	for _, v := range dAtA.resistXiongNuScoreData.Array {
		c.ResistXiongNuScore = append(c.ResistXiongNuScore, v.Encode())
	}

	for _, v := range dAtA.xuanyuanRangeData.Array {
		c.Gen.XuanyuanRangeData = append(c.Gen.XuanyuanRangeData, v.Encode())
	}

	for _, v := range dAtA.xuanyuanRankPrizeData.Array {
		c.Gen.XuanyuanRankPrizeData = append(c.Gen.XuanyuanRankPrizeData, v.Encode())
	}

	for _, v := range dAtA.zhanJiangChapterData.Array {
		c.ZhanJiangChapter = append(c.ZhanJiangChapter, v.Encode())
	}

	for _, v := range dAtA.zhengWuRefreshData.Array {
		c.ZhengWuRefresh = append(c.ZhengWuRefresh, v.Encode())
	}

	c.BaiZhanMisc = dAtA.baiZhanMiscData.Encode()

	c.Gen.CombatMiscConfig = dAtA.combatMiscConfig.Encode()

	c.Gen.CountryMiscData = dAtA.countryMiscData.Encode()

	c.DianquanMisc = dAtA.exchangeMiscData.Encode()

	c.CityEventMisc = dAtA.cityEventMiscData.Encode()

	c.DungeonMisc = dAtA.dungeonMiscData.Encode()

	c.FarmMiscConfig = dAtA.farmMiscConfig.Encode()

	c.GradonConfig = dAtA.gardenConfig.Encode()

	c.HebiMisc = dAtA.hebiMiscData.Encode()

	c.JiuGuanMisc = dAtA.jiuGuanMiscData.Encode()

	c.JunYingMisc = dAtA.junYingMiscData.Encode()

	c.Gen.McBuildMiscData = dAtA.mcBuildMiscData.Encode()

	c.Gen.MingcMiscData = dAtA.mingcMiscData.Encode()

	c.Gen.PromotionMiscData = dAtA.promotionMiscData.Encode()

	c.QuestionMisc = dAtA.questionMiscData.Encode()

	c.RankMisc = dAtA.rankMiscData.Encode()

	c.SeasonMisc = dAtA.seasonMiscData.Encode()

	c.Gen.ShopMiscData = dAtA.shopMiscData.Encode()

	c.GoodsConfig = dAtA.goodsConfig.Encode()

	c.GuildConfig = dAtA.guildConfig.Encode()

	c.Gen.GuildGenConfig = dAtA.guildGenConfig.Encode()

	c.MilitaryConfig = dAtA.militaryConfig.Encode()

	c.MiscConfig = dAtA.miscConfig.Encode()

	c.Gen.MiscGenConfig = dAtA.miscGenConfig.Encode()

	c.RegionConfig = dAtA.regionConfig.Encode()

	c.Gen.RegionGenConfig = dAtA.regionGenConfig.Encode()

	c.TagMisc = dAtA.tagMiscData.Encode()

	c.TaskMiscData = dAtA.taskMiscData.Encode()

	c.SecretTowerMisc = dAtA.secretTowerMiscData.Encode()

	c.Gen.VipMiscData = dAtA.vipMiscData.Encode()

	c.ResistXiongNuMisc = dAtA.resistXiongNuMisc.Encode()

	c.Gen.XuanyuanMiscData = dAtA.xuanyuanMiscData.Encode()

	c.ZhanJiangMisc = dAtA.zhanJiangMiscData.Encode()

	c.ZhengWuMisc = dAtA.zhengWuMiscData.Encode()

	return c
}

type Configs interface {
	GetActivityCollectionData(key uint64) *activitydata.ActivityCollectionData
	GetActivityCollectionDataArray() []*activitydata.ActivityCollectionData
	ActivityCollectionData() *ActivityCollectionDataConfig

	GetActivityShowData(key uint64) *activitydata.ActivityShowData
	GetActivityShowDataArray() []*activitydata.ActivityShowData
	ActivityShowData() *ActivityShowDataConfig

	GetActivityTaskListModeData(key uint64) *activitydata.ActivityTaskListModeData
	GetActivityTaskListModeDataArray() []*activitydata.ActivityTaskListModeData
	ActivityTaskListModeData() *ActivityTaskListModeDataConfig

	GetCollectionExchangeData(key uint64) *activitydata.CollectionExchangeData
	GetCollectionExchangeDataArray() []*activitydata.CollectionExchangeData
	CollectionExchangeData() *CollectionExchangeDataConfig

	GetJunXianLevelData(key uint64) *bai_zhan_data.JunXianLevelData
	GetJunXianLevelDataArray() []*bai_zhan_data.JunXianLevelData
	JunXianLevelData() *JunXianLevelDataConfig

	GetJunXianPrizeData(key uint64) *bai_zhan_data.JunXianPrizeData
	GetJunXianPrizeDataArray() []*bai_zhan_data.JunXianPrizeData
	JunXianPrizeData() *JunXianPrizeDataConfig

	GetHomeNpcBaseData(key uint64) *basedata.HomeNpcBaseData
	GetHomeNpcBaseDataArray() []*basedata.HomeNpcBaseData
	HomeNpcBaseData() *HomeNpcBaseDataConfig

	GetNpcBaseData(key uint64) *basedata.NpcBaseData
	GetNpcBaseDataArray() []*basedata.NpcBaseData
	NpcBaseData() *NpcBaseDataConfig

	GetBlockData(key uint64) *blockdata.BlockData
	GetBlockDataArray() []*blockdata.BlockData
	BlockData() *BlockDataConfig

	GetBodyData(key uint64) *body.BodyData
	GetBodyDataArray() []*body.BodyData
	BodyData() *BodyDataConfig

	GetBufferData(key uint64) *buffer.BufferData
	GetBufferDataArray() []*buffer.BufferData
	BufferData() *BufferDataConfig

	GetBufferTypeData(key uint64) *buffer.BufferTypeData
	GetBufferTypeDataArray() []*buffer.BufferTypeData
	BufferTypeData() *BufferTypeDataConfig

	GetCaptainAbilityData(key uint64) *captain.CaptainAbilityData
	GetCaptainAbilityDataArray() []*captain.CaptainAbilityData
	CaptainAbilityData() *CaptainAbilityDataConfig

	GetCaptainData(key uint64) *captain.CaptainData
	GetCaptainDataArray() []*captain.CaptainData
	CaptainData() *CaptainDataConfig

	GetCaptainFriendshipData(key uint64) *captain.CaptainFriendshipData
	GetCaptainFriendshipDataArray() []*captain.CaptainFriendshipData
	CaptainFriendshipData() *CaptainFriendshipDataConfig

	GetCaptainLevelData(key uint64) *captain.CaptainLevelData
	GetCaptainLevelDataArray() []*captain.CaptainLevelData
	CaptainLevelData() *CaptainLevelDataConfig

	GetCaptainOfficialCountData(key uint64) *captain.CaptainOfficialCountData
	GetCaptainOfficialCountDataArray() []*captain.CaptainOfficialCountData
	CaptainOfficialCountData() *CaptainOfficialCountDataConfig

	GetCaptainOfficialData(key uint64) *captain.CaptainOfficialData
	GetCaptainOfficialDataArray() []*captain.CaptainOfficialData
	CaptainOfficialData() *CaptainOfficialDataConfig

	GetCaptainRarityData(key uint64) *captain.CaptainRarityData
	GetCaptainRarityDataArray() []*captain.CaptainRarityData
	CaptainRarityData() *CaptainRarityDataConfig

	GetCaptainRebirthLevelData(key uint64) *captain.CaptainRebirthLevelData
	GetCaptainRebirthLevelDataArray() []*captain.CaptainRebirthLevelData
	CaptainRebirthLevelData() *CaptainRebirthLevelDataConfig

	GetCaptainStarData(key uint64) *captain.CaptainStarData
	GetCaptainStarDataArray() []*captain.CaptainStarData
	CaptainStarData() *CaptainStarDataConfig

	GetNamelessCaptainData(key uint64) *captain.NamelessCaptainData
	GetNamelessCaptainDataArray() []*captain.NamelessCaptainData
	NamelessCaptainData() *NamelessCaptainDataConfig

	GetChargeObjData(key uint64) *charge.ChargeObjData
	GetChargeObjDataArray() []*charge.ChargeObjData
	ChargeObjData() *ChargeObjDataConfig

	GetChargePrizeData(key uint64) *charge.ChargePrizeData
	GetChargePrizeDataArray() []*charge.ChargePrizeData
	ChargePrizeData() *ChargePrizeDataConfig

	GetProductData(key uint64) *charge.ProductData
	GetProductDataArray() []*charge.ProductData
	ProductData() *ProductDataConfig

	GetEquipCombineData(key uint64) *combine.EquipCombineData
	GetEquipCombineDataArray() []*combine.EquipCombineData
	EquipCombineData() *EquipCombineDataConfig

	GetGoodsCombineData(key uint64) *combine.GoodsCombineData
	GetGoodsCombineDataArray() []*combine.GoodsCombineData
	GoodsCombineData() *GoodsCombineDataConfig

	GetCountryData(key uint64) *country.CountryData
	GetCountryDataArray() []*country.CountryData
	CountryData() *CountryDataConfig

	GetCountryOfficialData(key int) *country.CountryOfficialData
	GetCountryOfficialDataArray() []*country.CountryOfficialData
	CountryOfficialData() *CountryOfficialDataConfig

	GetCountryOfficialNpcData(key uint64) *country.CountryOfficialNpcData
	GetCountryOfficialNpcDataArray() []*country.CountryOfficialNpcData
	CountryOfficialNpcData() *CountryOfficialNpcDataConfig

	GetFamilyNameData(key uint64) *country.FamilyNameData
	GetFamilyNameDataArray() []*country.FamilyNameData
	FamilyNameData() *FamilyNameDataConfig

	GetBroadcastData(key string) *data.BroadcastData
	GetBroadcastDataArray() []*data.BroadcastData
	BroadcastData() *BroadcastDataConfig

	GetBuffEffectData(key uint64) *data.BuffEffectData
	GetBuffEffectDataArray() []*data.BuffEffectData
	BuffEffectData() *BuffEffectDataConfig

	GetColorData(key uint64) *data.ColorData
	GetColorDataArray() []*data.ColorData
	ColorData() *ColorDataConfig

	GetFamilyName(key string) *data.FamilyName
	GetFamilyNameArray() []*data.FamilyName
	FamilyName() *FamilyNameConfig

	GetFemaleGivenName(key string) *data.FemaleGivenName
	GetFemaleGivenNameArray() []*data.FemaleGivenName
	FemaleGivenName() *FemaleGivenNameConfig

	GetHeroLevelSubData(key uint64) *data.HeroLevelSubData
	GetHeroLevelSubDataArray() []*data.HeroLevelSubData
	HeroLevelSubData() *HeroLevelSubDataConfig

	GetMaleGivenName(key string) *data.MaleGivenName
	GetMaleGivenNameArray() []*data.MaleGivenName
	MaleGivenName() *MaleGivenNameConfig

	GetSpriteStat(key uint64) *data.SpriteStat
	GetSpriteStatArray() []*data.SpriteStat
	SpriteStat() *SpriteStatConfig

	GetText(key string) *data.Text
	GetTextArray() []*data.Text
	Text() *TextConfig

	GetTimeRuleData(key uint64) *data.TimeRuleData
	GetTimeRuleDataArray() []*data.TimeRuleData
	TimeRuleData() *TimeRuleDataConfig

	GetBaseLevelData(key uint64) *domestic_data.BaseLevelData
	GetBaseLevelDataArray() []*domestic_data.BaseLevelData
	BaseLevelData() *BaseLevelDataConfig

	GetBuildingData(key uint64) *domestic_data.BuildingData
	GetBuildingDataArray() []*domestic_data.BuildingData
	BuildingData() *BuildingDataConfig

	GetBuildingLayoutData(key uint64) *domestic_data.BuildingLayoutData
	GetBuildingLayoutDataArray() []*domestic_data.BuildingLayoutData
	BuildingLayoutData() *BuildingLayoutDataConfig

	GetBuildingUnlockData(key uint64) *domestic_data.BuildingUnlockData
	GetBuildingUnlockDataArray() []*domestic_data.BuildingUnlockData
	BuildingUnlockData() *BuildingUnlockDataConfig

	GetCityEventData(key uint64) *domestic_data.CityEventData
	GetCityEventDataArray() []*domestic_data.CityEventData
	CityEventData() *CityEventDataConfig

	GetCityEventLevelData(key uint64) *domestic_data.CityEventLevelData
	GetCityEventLevelDataArray() []*domestic_data.CityEventLevelData
	CityEventLevelData() *CityEventLevelDataConfig

	GetCombineCost(key int) *domestic_data.CombineCost
	GetCombineCostArray() []*domestic_data.CombineCost
	CombineCost() *CombineCostConfig

	GetCountdownPrizeData(key uint64) *domestic_data.CountdownPrizeData
	GetCountdownPrizeDataArray() []*domestic_data.CountdownPrizeData
	CountdownPrizeData() *CountdownPrizeDataConfig

	GetCountdownPrizeDescData(key uint64) *domestic_data.CountdownPrizeDescData
	GetCountdownPrizeDescDataArray() []*domestic_data.CountdownPrizeDescData
	CountdownPrizeDescData() *CountdownPrizeDescDataConfig

	GetGuanFuLevelData(key uint64) *domestic_data.GuanFuLevelData
	GetGuanFuLevelDataArray() []*domestic_data.GuanFuLevelData
	GuanFuLevelData() *GuanFuLevelDataConfig

	GetOuterCityBuildingData(key uint64) *domestic_data.OuterCityBuildingData
	GetOuterCityBuildingDataArray() []*domestic_data.OuterCityBuildingData
	OuterCityBuildingData() *OuterCityBuildingDataConfig

	GetOuterCityBuildingDescData(key uint64) *domestic_data.OuterCityBuildingDescData
	GetOuterCityBuildingDescDataArray() []*domestic_data.OuterCityBuildingDescData
	OuterCityBuildingDescData() *OuterCityBuildingDescDataConfig

	GetOuterCityData(key uint64) *domestic_data.OuterCityData
	GetOuterCityDataArray() []*domestic_data.OuterCityData
	OuterCityData() *OuterCityDataConfig

	GetOuterCityLayoutData(key uint64) *domestic_data.OuterCityLayoutData
	GetOuterCityLayoutDataArray() []*domestic_data.OuterCityLayoutData
	OuterCityLayoutData() *OuterCityLayoutDataConfig

	GetProsperityDamageBuffData(key uint64) *domestic_data.ProsperityDamageBuffData
	GetProsperityDamageBuffDataArray() []*domestic_data.ProsperityDamageBuffData
	ProsperityDamageBuffData() *ProsperityDamageBuffDataConfig

	GetSoldierLevelData(key uint64) *domestic_data.SoldierLevelData
	GetSoldierLevelDataArray() []*domestic_data.SoldierLevelData
	SoldierLevelData() *SoldierLevelDataConfig

	GetTechnologyData(key uint64) *domestic_data.TechnologyData
	GetTechnologyDataArray() []*domestic_data.TechnologyData
	TechnologyData() *TechnologyDataConfig

	GetTieJiangPuLevelData(key uint64) *domestic_data.TieJiangPuLevelData
	GetTieJiangPuLevelDataArray() []*domestic_data.TieJiangPuLevelData
	TieJiangPuLevelData() *TieJiangPuLevelDataConfig

	GetWorkshopDuration(key uint64) *domestic_data.WorkshopDuration
	GetWorkshopDurationArray() []*domestic_data.WorkshopDuration
	WorkshopDuration() *WorkshopDurationConfig

	GetWorkshopLevelData(key uint64) *domestic_data.WorkshopLevelData
	GetWorkshopLevelDataArray() []*domestic_data.WorkshopLevelData
	WorkshopLevelData() *WorkshopLevelDataConfig

	GetWorkshopRefreshCost(key uint64) *domestic_data.WorkshopRefreshCost
	GetWorkshopRefreshCostArray() []*domestic_data.WorkshopRefreshCost
	WorkshopRefreshCost() *WorkshopRefreshCostConfig

	GetDungeonChapterData(key uint64) *dungeon.DungeonChapterData
	GetDungeonChapterDataArray() []*dungeon.DungeonChapterData
	DungeonChapterData() *DungeonChapterDataConfig

	GetDungeonData(key uint64) *dungeon.DungeonData
	GetDungeonDataArray() []*dungeon.DungeonData
	DungeonData() *DungeonDataConfig

	GetDungeonGuideTroopData(key uint64) *dungeon.DungeonGuideTroopData
	GetDungeonGuideTroopDataArray() []*dungeon.DungeonGuideTroopData
	DungeonGuideTroopData() *DungeonGuideTroopDataConfig

	GetFarmMaxStealConfig(key uint64) *farm.FarmMaxStealConfig
	GetFarmMaxStealConfigArray() []*farm.FarmMaxStealConfig
	FarmMaxStealConfig() *FarmMaxStealConfigConfig

	GetFarmOneKeyConfig(key uint64) *farm.FarmOneKeyConfig
	GetFarmOneKeyConfigArray() []*farm.FarmOneKeyConfig
	FarmOneKeyConfig() *FarmOneKeyConfigConfig

	GetFarmResConfig(key uint64) *farm.FarmResConfig
	GetFarmResConfigArray() []*farm.FarmResConfig
	FarmResConfig() *FarmResConfigConfig

	GetFishData(key uint64) *fishing_data.FishData
	GetFishDataArray() []*fishing_data.FishData
	FishData() *FishDataConfig

	GetFishingCaptainProbabilityData(key uint64) *fishing_data.FishingCaptainProbabilityData
	GetFishingCaptainProbabilityDataArray() []*fishing_data.FishingCaptainProbabilityData
	FishingCaptainProbabilityData() *FishingCaptainProbabilityDataConfig

	GetFishingCostData(key uint64) *fishing_data.FishingCostData
	GetFishingCostDataArray() []*fishing_data.FishingCostData
	FishingCostData() *FishingCostDataConfig

	GetFishingShowData(key uint64) *fishing_data.FishingShowData
	GetFishingShowDataArray() []*fishing_data.FishingShowData
	FishingShowData() *FishingShowDataConfig

	GetFunctionOpenData(key uint64) *function.FunctionOpenData
	GetFunctionOpenDataArray() []*function.FunctionOpenData
	FunctionOpenData() *FunctionOpenDataConfig

	GetTreasuryTreeData(key uint64) *gardendata.TreasuryTreeData
	GetTreasuryTreeDataArray() []*gardendata.TreasuryTreeData
	TreasuryTreeData() *TreasuryTreeDataConfig

	GetEquipmentData(key uint64) *goods.EquipmentData
	GetEquipmentDataArray() []*goods.EquipmentData
	EquipmentData() *EquipmentDataConfig

	GetEquipmentLevelData(key uint64) *goods.EquipmentLevelData
	GetEquipmentLevelDataArray() []*goods.EquipmentLevelData
	EquipmentLevelData() *EquipmentLevelDataConfig

	GetEquipmentQualityData(key uint64) *goods.EquipmentQualityData
	GetEquipmentQualityDataArray() []*goods.EquipmentQualityData
	EquipmentQualityData() *EquipmentQualityDataConfig

	GetEquipmentRefinedData(key uint64) *goods.EquipmentRefinedData
	GetEquipmentRefinedDataArray() []*goods.EquipmentRefinedData
	EquipmentRefinedData() *EquipmentRefinedDataConfig

	GetEquipmentTaozData(key uint64) *goods.EquipmentTaozData
	GetEquipmentTaozDataArray() []*goods.EquipmentTaozData
	EquipmentTaozData() *EquipmentTaozDataConfig

	GetGemData(key uint64) *goods.GemData
	GetGemDataArray() []*goods.GemData
	GemData() *GemDataConfig

	GetGoodsData(key uint64) *goods.GoodsData
	GetGoodsDataArray() []*goods.GoodsData
	GoodsData() *GoodsDataConfig

	GetGoodsQuality(key uint64) *goods.GoodsQuality
	GetGoodsQualityArray() []*goods.GoodsQuality
	GoodsQuality() *GoodsQualityConfig

	GetGuildBigBoxData(key uint64) *guild_data.GuildBigBoxData
	GetGuildBigBoxDataArray() []*guild_data.GuildBigBoxData
	GuildBigBoxData() *GuildBigBoxDataConfig

	GetGuildClassLevelData(key uint64) *guild_data.GuildClassLevelData
	GetGuildClassLevelDataArray() []*guild_data.GuildClassLevelData
	GuildClassLevelData() *GuildClassLevelDataConfig

	GetGuildClassTitleData(key uint64) *guild_data.GuildClassTitleData
	GetGuildClassTitleDataArray() []*guild_data.GuildClassTitleData
	GuildClassTitleData() *GuildClassTitleDataConfig

	GetGuildDonateData(key uint64) *guild_data.GuildDonateData
	GetGuildDonateDataArray() []*guild_data.GuildDonateData
	GuildDonateData() *GuildDonateDataConfig

	GetGuildEventPrizeData(key uint64) *guild_data.GuildEventPrizeData
	GetGuildEventPrizeDataArray() []*guild_data.GuildEventPrizeData
	GuildEventPrizeData() *GuildEventPrizeDataConfig

	GetGuildLevelCdrData(key uint64) *guild_data.GuildLevelCdrData
	GetGuildLevelCdrDataArray() []*guild_data.GuildLevelCdrData
	GuildLevelCdrData() *GuildLevelCdrDataConfig

	GetGuildLevelData(key uint64) *guild_data.GuildLevelData
	GetGuildLevelDataArray() []*guild_data.GuildLevelData
	GuildLevelData() *GuildLevelDataConfig

	GetGuildLogData(key string) *guild_data.GuildLogData
	GetGuildLogDataArray() []*guild_data.GuildLogData
	GuildLogData() *GuildLogDataConfig

	GetGuildPermissionShowData(key uint64) *guild_data.GuildPermissionShowData
	GetGuildPermissionShowDataArray() []*guild_data.GuildPermissionShowData
	GuildPermissionShowData() *GuildPermissionShowDataConfig

	GetGuildPrestigeEventData(key uint64) *guild_data.GuildPrestigeEventData
	GetGuildPrestigeEventDataArray() []*guild_data.GuildPrestigeEventData
	GuildPrestigeEventData() *GuildPrestigeEventDataConfig

	GetGuildPrestigePrizeData(key uint64) *guild_data.GuildPrestigePrizeData
	GetGuildPrestigePrizeDataArray() []*guild_data.GuildPrestigePrizeData
	GuildPrestigePrizeData() *GuildPrestigePrizeDataConfig

	GetGuildRankPrizeData(key uint64) *guild_data.GuildRankPrizeData
	GetGuildRankPrizeDataArray() []*guild_data.GuildRankPrizeData
	GuildRankPrizeData() *GuildRankPrizeDataConfig

	GetGuildTarget(key uint64) *guild_data.GuildTarget
	GetGuildTargetArray() []*guild_data.GuildTarget
	GuildTarget() *GuildTargetConfig

	GetGuildTaskData(key uint64) *guild_data.GuildTaskData
	GetGuildTaskDataArray() []*guild_data.GuildTaskData
	GuildTaskData() *GuildTaskDataConfig

	GetGuildTaskEvaluateData(key uint64) *guild_data.GuildTaskEvaluateData
	GetGuildTaskEvaluateDataArray() []*guild_data.GuildTaskEvaluateData
	GuildTaskEvaluateData() *GuildTaskEvaluateDataConfig

	GetGuildTechnologyData(key uint64) *guild_data.GuildTechnologyData
	GetGuildTechnologyDataArray() []*guild_data.GuildTechnologyData
	GuildTechnologyData() *GuildTechnologyDataConfig

	GetNpcGuildTemplate(key uint64) *guild_data.NpcGuildTemplate
	GetNpcGuildTemplateArray() []*guild_data.NpcGuildTemplate
	NpcGuildTemplate() *NpcGuildTemplateConfig

	GetNpcMemberData(key uint64) *guild_data.NpcMemberData
	GetNpcMemberDataArray() []*guild_data.NpcMemberData
	NpcMemberData() *NpcMemberDataConfig

	GetHeadData(key string) *head.HeadData
	GetHeadDataArray() []*head.HeadData
	HeadData() *HeadDataConfig

	GetHebiPrizeData(key uint64) *hebi.HebiPrizeData
	GetHebiPrizeDataArray() []*hebi.HebiPrizeData
	HebiPrizeData() *HebiPrizeDataConfig

	GetHeroLevelData(key uint64) *herodata.HeroLevelData
	GetHeroLevelDataArray() []*herodata.HeroLevelData
	HeroLevelData() *HeroLevelDataConfig

	GetI18nData(key string) *i18n.I18nData
	GetI18nDataArray() []*i18n.I18nData
	I18nData() *I18nDataConfig

	GetIcon(key string) *icon.Icon
	GetIconArray() []*icon.Icon
	Icon() *IconConfig

	GetLocationData(key uint64) *location.LocationData
	GetLocationDataArray() []*location.LocationData
	LocationData() *LocationDataConfig

	GetMailData(key string) *maildata.MailData
	GetMailDataArray() []*maildata.MailData
	MailData() *MailDataConfig

	GetJiuGuanData(key uint64) *military_data.JiuGuanData
	GetJiuGuanDataArray() []*military_data.JiuGuanData
	JiuGuanData() *JiuGuanDataConfig

	GetJunYingLevelData(key uint64) *military_data.JunYingLevelData
	GetJunYingLevelDataArray() []*military_data.JunYingLevelData
	JunYingLevelData() *JunYingLevelDataConfig

	GetTrainingLevelData(key uint64) *military_data.TrainingLevelData
	GetTrainingLevelDataArray() []*military_data.TrainingLevelData
	TrainingLevelData() *TrainingLevelDataConfig

	GetTutorData(key uint64) *military_data.TutorData
	GetTutorDataArray() []*military_data.TutorData
	TutorData() *TutorDataConfig

	GetMcBuildAddSupportData(key uint64) *mingcdata.McBuildAddSupportData
	GetMcBuildAddSupportDataArray() []*mingcdata.McBuildAddSupportData
	McBuildAddSupportData() *McBuildAddSupportDataConfig

	GetMcBuildGuildMemberPrizeData(key uint64) *mingcdata.McBuildGuildMemberPrizeData
	GetMcBuildGuildMemberPrizeDataArray() []*mingcdata.McBuildGuildMemberPrizeData
	McBuildGuildMemberPrizeData() *McBuildGuildMemberPrizeDataConfig

	GetMcBuildMcSupportData(key uint64) *mingcdata.McBuildMcSupportData
	GetMcBuildMcSupportDataArray() []*mingcdata.McBuildMcSupportData
	McBuildMcSupportData() *McBuildMcSupportDataConfig

	GetMingcBaseData(key uint64) *mingcdata.MingcBaseData
	GetMingcBaseDataArray() []*mingcdata.MingcBaseData
	MingcBaseData() *MingcBaseDataConfig

	GetMingcTimeData(key uint64) *mingcdata.MingcTimeData
	GetMingcTimeDataArray() []*mingcdata.MingcTimeData
	MingcTimeData() *MingcTimeDataConfig

	GetMingcWarBuildingData(key uint64) *mingcdata.MingcWarBuildingData
	GetMingcWarBuildingDataArray() []*mingcdata.MingcWarBuildingData
	MingcWarBuildingData() *MingcWarBuildingDataConfig

	GetMingcWarDrumStatData(key uint64) *mingcdata.MingcWarDrumStatData
	GetMingcWarDrumStatDataArray() []*mingcdata.MingcWarDrumStatData
	MingcWarDrumStatData() *MingcWarDrumStatDataConfig

	GetMingcWarMapData(key uint64) *mingcdata.MingcWarMapData
	GetMingcWarMapDataArray() []*mingcdata.MingcWarMapData
	MingcWarMapData() *MingcWarMapDataConfig

	GetMingcWarMultiKillData(key uint64) *mingcdata.MingcWarMultiKillData
	GetMingcWarMultiKillDataArray() []*mingcdata.MingcWarMultiKillData
	MingcWarMultiKillData() *MingcWarMultiKillDataConfig

	GetMingcWarNpcData(key uint64) *mingcdata.MingcWarNpcData
	GetMingcWarNpcDataArray() []*mingcdata.MingcWarNpcData
	MingcWarNpcData() *MingcWarNpcDataConfig

	GetMingcWarNpcGuildData(key uint64) *mingcdata.MingcWarNpcGuildData
	GetMingcWarNpcGuildDataArray() []*mingcdata.MingcWarNpcGuildData
	MingcWarNpcGuildData() *MingcWarNpcGuildDataConfig

	GetMingcWarSceneData(key uint64) *mingcdata.MingcWarSceneData
	GetMingcWarSceneDataArray() []*mingcdata.MingcWarSceneData
	MingcWarSceneData() *MingcWarSceneDataConfig

	GetMingcWarTouShiBuildingTargetData(key uint64) *mingcdata.MingcWarTouShiBuildingTargetData
	GetMingcWarTouShiBuildingTargetDataArray() []*mingcdata.MingcWarTouShiBuildingTargetData
	MingcWarTouShiBuildingTargetData() *MingcWarTouShiBuildingTargetDataConfig

	GetMingcWarTroopLastBeatWhenFailData(key uint64) *mingcdata.MingcWarTroopLastBeatWhenFailData
	GetMingcWarTroopLastBeatWhenFailDataArray() []*mingcdata.MingcWarTroopLastBeatWhenFailData
	MingcWarTroopLastBeatWhenFailData() *MingcWarTroopLastBeatWhenFailDataConfig

	GetMonsterCaptainData(key uint64) *monsterdata.MonsterCaptainData
	GetMonsterCaptainDataArray() []*monsterdata.MonsterCaptainData
	MonsterCaptainData() *MonsterCaptainDataConfig

	GetMonsterMasterData(key uint64) *monsterdata.MonsterMasterData
	GetMonsterMasterDataArray() []*monsterdata.MonsterMasterData
	MonsterMasterData() *MonsterMasterDataConfig

	GetDailyBargainData(key uint64) *promdata.DailyBargainData
	GetDailyBargainDataArray() []*promdata.DailyBargainData
	DailyBargainData() *DailyBargainDataConfig

	GetDurationCardData(key uint64) *promdata.DurationCardData
	GetDurationCardDataArray() []*promdata.DurationCardData
	DurationCardData() *DurationCardDataConfig

	GetEventLimitGiftData(key uint64) *promdata.EventLimitGiftData
	GetEventLimitGiftDataArray() []*promdata.EventLimitGiftData
	EventLimitGiftData() *EventLimitGiftDataConfig

	GetFreeGiftData(key uint64) *promdata.FreeGiftData
	GetFreeGiftDataArray() []*promdata.FreeGiftData
	FreeGiftData() *FreeGiftDataConfig

	GetHeroLevelFundData(key uint64) *promdata.HeroLevelFundData
	GetHeroLevelFundDataArray() []*promdata.HeroLevelFundData
	HeroLevelFundData() *HeroLevelFundDataConfig

	GetLoginDayData(key uint64) *promdata.LoginDayData
	GetLoginDayDataArray() []*promdata.LoginDayData
	LoginDayData() *LoginDayDataConfig

	GetSpCollectionData(key uint64) *promdata.SpCollectionData
	GetSpCollectionDataArray() []*promdata.SpCollectionData
	SpCollectionData() *SpCollectionDataConfig

	GetTimeLimitGiftData(key uint64) *promdata.TimeLimitGiftData
	GetTimeLimitGiftDataArray() []*promdata.TimeLimitGiftData
	TimeLimitGiftData() *TimeLimitGiftDataConfig

	GetTimeLimitGiftGroupData(key uint64) *promdata.TimeLimitGiftGroupData
	GetTimeLimitGiftGroupDataArray() []*promdata.TimeLimitGiftGroupData
	TimeLimitGiftGroupData() *TimeLimitGiftGroupDataConfig

	GetPushData(key uint64) *pushdata.PushData
	GetPushDataArray() []*pushdata.PushData
	PushData() *PushDataConfig

	GetPveTroopData(key uint64) *pvetroop.PveTroopData
	GetPveTroopDataArray() []*pvetroop.PveTroopData
	PveTroopData() *PveTroopDataConfig

	GetQuestionData(key uint64) *question.QuestionData
	GetQuestionDataArray() []*question.QuestionData
	QuestionData() *QuestionDataConfig

	GetQuestionPrizeData(key uint64) *question.QuestionPrizeData
	GetQuestionPrizeDataArray() []*question.QuestionPrizeData
	QuestionPrizeData() *QuestionPrizeDataConfig

	GetQuestionSayingData(key uint64) *question.QuestionSayingData
	GetQuestionSayingDataArray() []*question.QuestionSayingData
	QuestionSayingData() *QuestionSayingDataConfig

	GetRaceData(key int) *race.RaceData
	GetRaceDataArray() []*race.RaceData
	RaceData() *RaceDataConfig

	GetEventOptionData(key uint64) *random_event.EventOptionData
	GetEventOptionDataArray() []*random_event.EventOptionData
	EventOptionData() *EventOptionDataConfig

	GetEventPosition(key uint64) *random_event.EventPosition
	GetEventPositionArray() []*random_event.EventPosition
	EventPosition() *EventPositionConfig

	GetOptionPrize(key uint64) *random_event.OptionPrize
	GetOptionPrizeArray() []*random_event.OptionPrize
	OptionPrize() *OptionPrizeConfig

	GetRandomEventData(key uint64) *random_event.RandomEventData
	GetRandomEventDataArray() []*random_event.RandomEventData
	RandomEventData() *RandomEventDataConfig

	GetRedPacketData(key uint64) *red_packet.RedPacketData
	GetRedPacketDataArray() []*red_packet.RedPacketData
	RedPacketData() *RedPacketDataConfig

	GetAreaData(key uint64) *regdata.AreaData
	GetAreaDataArray() []*regdata.AreaData
	AreaData() *AreaDataConfig

	GetAssemblyData(key uint64) *regdata.AssemblyData
	GetAssemblyDataArray() []*regdata.AssemblyData
	AssemblyData() *AssemblyDataConfig

	GetBaozNpcData(key uint64) *regdata.BaozNpcData
	GetBaozNpcDataArray() []*regdata.BaozNpcData
	BaozNpcData() *BaozNpcDataConfig

	GetJunTuanNpcData(key uint64) *regdata.JunTuanNpcData
	GetJunTuanNpcDataArray() []*regdata.JunTuanNpcData
	JunTuanNpcData() *JunTuanNpcDataConfig

	GetJunTuanNpcPlaceData(key uint64) *regdata.JunTuanNpcPlaceData
	GetJunTuanNpcPlaceDataArray() []*regdata.JunTuanNpcPlaceData
	JunTuanNpcPlaceData() *JunTuanNpcPlaceDataConfig

	GetRegionAreaData(key uint64) *regdata.RegionAreaData
	GetRegionAreaDataArray() []*regdata.RegionAreaData
	RegionAreaData() *RegionAreaDataConfig

	GetRegionData(key uint64) *regdata.RegionData
	GetRegionDataArray() []*regdata.RegionData
	RegionData() *RegionDataConfig

	GetRegionMonsterData(key uint64) *regdata.RegionMonsterData
	GetRegionMonsterDataArray() []*regdata.RegionMonsterData
	RegionMonsterData() *RegionMonsterDataConfig

	GetRegionMultiLevelNpcData(key uint64) *regdata.RegionMultiLevelNpcData
	GetRegionMultiLevelNpcDataArray() []*regdata.RegionMultiLevelNpcData
	RegionMultiLevelNpcData() *RegionMultiLevelNpcDataConfig

	GetRegionMultiLevelNpcLevelData(key uint64) *regdata.RegionMultiLevelNpcLevelData
	GetRegionMultiLevelNpcLevelDataArray() []*regdata.RegionMultiLevelNpcLevelData
	RegionMultiLevelNpcLevelData() *RegionMultiLevelNpcLevelDataConfig

	GetRegionMultiLevelNpcTypeData(key int) *regdata.RegionMultiLevelNpcTypeData
	GetRegionMultiLevelNpcTypeDataArray() []*regdata.RegionMultiLevelNpcTypeData
	RegionMultiLevelNpcTypeData() *RegionMultiLevelNpcTypeDataConfig

	GetTroopDialogueData(key uint64) *regdata.TroopDialogueData
	GetTroopDialogueDataArray() []*regdata.TroopDialogueData
	TroopDialogueData() *TroopDialogueDataConfig

	GetTroopDialogueTextData(key uint64) *regdata.TroopDialogueTextData
	GetTroopDialogueTextDataArray() []*regdata.TroopDialogueTextData
	TroopDialogueTextData() *TroopDialogueTextDataConfig

	GetAmountShowSortData(key uint64) *resdata.AmountShowSortData
	GetAmountShowSortDataArray() []*resdata.AmountShowSortData
	AmountShowSortData() *AmountShowSortDataConfig

	GetBaowuData(key uint64) *resdata.BaowuData
	GetBaowuDataArray() []*resdata.BaowuData
	BaowuData() *BaowuDataConfig

	GetConditionPlunder(key uint64) *resdata.ConditionPlunder
	GetConditionPlunderArray() []*resdata.ConditionPlunder
	ConditionPlunder() *ConditionPlunderConfig

	GetConditionPlunderItem(key uint64) *resdata.ConditionPlunderItem
	GetConditionPlunderItemArray() []*resdata.ConditionPlunderItem
	ConditionPlunderItem() *ConditionPlunderItemConfig

	GetCost(key int) *resdata.Cost
	GetCostArray() []*resdata.Cost
	Cost() *CostConfig

	GetGuildLevelPrize(key uint64) *resdata.GuildLevelPrize
	GetGuildLevelPrizeArray() []*resdata.GuildLevelPrize
	GuildLevelPrize() *GuildLevelPrizeConfig

	GetPlunder(key uint64) *resdata.Plunder
	GetPlunderArray() []*resdata.Plunder
	Plunder() *PlunderConfig

	GetPlunderGroup(key uint64) *resdata.PlunderGroup
	GetPlunderGroupArray() []*resdata.PlunderGroup
	PlunderGroup() *PlunderGroupConfig

	GetPlunderItem(key uint64) *resdata.PlunderItem
	GetPlunderItemArray() []*resdata.PlunderItem
	PlunderItem() *PlunderItemConfig

	GetPlunderPrize(key uint64) *resdata.PlunderPrize
	GetPlunderPrizeArray() []*resdata.PlunderPrize
	PlunderPrize() *PlunderPrizeConfig

	GetPrize(key int) *resdata.Prize
	GetPrizeArray() []*resdata.Prize
	Prize() *PrizeConfig

	GetResCaptainData(key uint64) *resdata.ResCaptainData
	GetResCaptainDataArray() []*resdata.ResCaptainData
	ResCaptainData() *ResCaptainDataConfig

	GetCombatScene(key string) *scene.CombatScene
	GetCombatSceneArray() []*scene.CombatScene
	CombatScene() *CombatSceneConfig

	GetSeasonData(key uint64) *season.SeasonData
	GetSeasonDataArray() []*season.SeasonData
	SeasonData() *SeasonDataConfig

	GetPrivacySettingData(key uint64) *settings.PrivacySettingData
	GetPrivacySettingDataArray() []*settings.PrivacySettingData
	PrivacySettingData() *PrivacySettingDataConfig

	GetBlackMarketData(key uint64) *shop.BlackMarketData
	GetBlackMarketDataArray() []*shop.BlackMarketData
	BlackMarketData() *BlackMarketDataConfig

	GetBlackMarketGoodsData(key uint64) *shop.BlackMarketGoodsData
	GetBlackMarketGoodsDataArray() []*shop.BlackMarketGoodsData
	BlackMarketGoodsData() *BlackMarketGoodsDataConfig

	GetBlackMarketGoodsGroupData(key uint64) *shop.BlackMarketGoodsGroupData
	GetBlackMarketGoodsGroupDataArray() []*shop.BlackMarketGoodsGroupData
	BlackMarketGoodsGroupData() *BlackMarketGoodsGroupDataConfig

	GetDiscountColorData(key uint64) *shop.DiscountColorData
	GetDiscountColorDataArray() []*shop.DiscountColorData
	DiscountColorData() *DiscountColorDataConfig

	GetNormalShopGoods(key uint64) *shop.NormalShopGoods
	GetNormalShopGoodsArray() []*shop.NormalShopGoods
	NormalShopGoods() *NormalShopGoodsConfig

	GetShop(key uint64) *shop.Shop
	GetShopArray() []*shop.Shop
	Shop() *ShopConfig

	GetZhenBaoGeShopGoods(key uint64) *shop.ZhenBaoGeShopGoods
	GetZhenBaoGeShopGoodsArray() []*shop.ZhenBaoGeShopGoods
	ZhenBaoGeShopGoods() *ZhenBaoGeShopGoodsConfig

	GetPassiveSpellData(key uint64) *spell.PassiveSpellData
	GetPassiveSpellDataArray() []*spell.PassiveSpellData
	PassiveSpellData() *PassiveSpellDataConfig

	GetSpell(key uint64) *spell.Spell
	GetSpellArray() []*spell.Spell
	Spell() *SpellConfig

	GetSpellData(key uint64) *spell.SpellData
	GetSpellDataArray() []*spell.SpellData
	SpellData() *SpellDataConfig

	GetSpellFacadeData(key uint64) *spell.SpellFacadeData
	GetSpellFacadeDataArray() []*spell.SpellFacadeData
	SpellFacadeData() *SpellFacadeDataConfig

	GetStateData(key uint64) *spell.StateData
	GetStateDataArray() []*spell.StateData
	StateData() *StateDataConfig

	GetStrategyData(key uint64) *strategydata.StrategyData
	GetStrategyDataArray() []*strategydata.StrategyData
	StrategyData() *StrategyDataConfig

	GetStrategyEffectData(key uint64) *strategydata.StrategyEffectData
	GetStrategyEffectDataArray() []*strategydata.StrategyEffectData
	StrategyEffectData() *StrategyEffectDataConfig

	GetStrongerData(key uint64) *strongerdata.StrongerData
	GetStrongerDataArray() []*strongerdata.StrongerData
	StrongerData() *StrongerDataConfig

	GetBuildingEffectData(key int) *sub.BuildingEffectData
	GetBuildingEffectDataArray() []*sub.BuildingEffectData
	BuildingEffectData() *BuildingEffectDataConfig

	GetSurveyData(key string) *survey.SurveyData
	GetSurveyDataArray() []*survey.SurveyData
	SurveyData() *SurveyDataConfig

	GetAchieveTaskData(key uint64) *taskdata.AchieveTaskData
	GetAchieveTaskDataArray() []*taskdata.AchieveTaskData
	AchieveTaskData() *AchieveTaskDataConfig

	GetAchieveTaskStarPrizeData(key uint64) *taskdata.AchieveTaskStarPrizeData
	GetAchieveTaskStarPrizeDataArray() []*taskdata.AchieveTaskStarPrizeData
	AchieveTaskStarPrizeData() *AchieveTaskStarPrizeDataConfig

	GetActiveDegreePrizeData(key uint64) *taskdata.ActiveDegreePrizeData
	GetActiveDegreePrizeDataArray() []*taskdata.ActiveDegreePrizeData
	ActiveDegreePrizeData() *ActiveDegreePrizeDataConfig

	GetActiveDegreeTaskData(key uint64) *taskdata.ActiveDegreeTaskData
	GetActiveDegreeTaskDataArray() []*taskdata.ActiveDegreeTaskData
	ActiveDegreeTaskData() *ActiveDegreeTaskDataConfig

	GetActivityTaskData(key uint64) *taskdata.ActivityTaskData
	GetActivityTaskDataArray() []*taskdata.ActivityTaskData
	ActivityTaskData() *ActivityTaskDataConfig

	GetBaYeStageData(key uint64) *taskdata.BaYeStageData
	GetBaYeStageDataArray() []*taskdata.BaYeStageData
	BaYeStageData() *BaYeStageDataConfig

	GetBaYeTaskData(key uint64) *taskdata.BaYeTaskData
	GetBaYeTaskDataArray() []*taskdata.BaYeTaskData
	BaYeTaskData() *BaYeTaskDataConfig

	GetBranchTaskData(key uint64) *taskdata.BranchTaskData
	GetBranchTaskDataArray() []*taskdata.BranchTaskData
	BranchTaskData() *BranchTaskDataConfig

	GetBwzlPrizeData(key uint64) *taskdata.BwzlPrizeData
	GetBwzlPrizeDataArray() []*taskdata.BwzlPrizeData
	BwzlPrizeData() *BwzlPrizeDataConfig

	GetBwzlTaskData(key uint64) *taskdata.BwzlTaskData
	GetBwzlTaskDataArray() []*taskdata.BwzlTaskData
	BwzlTaskData() *BwzlTaskDataConfig

	GetMainTaskData(key uint64) *taskdata.MainTaskData
	GetMainTaskDataArray() []*taskdata.MainTaskData
	MainTaskData() *MainTaskDataConfig

	GetTaskBoxData(key uint64) *taskdata.TaskBoxData
	GetTaskBoxDataArray() []*taskdata.TaskBoxData
	TaskBoxData() *TaskBoxDataConfig

	GetTaskTargetData(key uint64) *taskdata.TaskTargetData
	GetTaskTargetDataArray() []*taskdata.TaskTargetData
	TaskTargetData() *TaskTargetDataConfig

	GetTitleData(key uint64) *taskdata.TitleData
	GetTitleDataArray() []*taskdata.TitleData
	TitleData() *TitleDataConfig

	GetTitleTaskData(key uint64) *taskdata.TitleTaskData
	GetTitleTaskDataArray() []*taskdata.TitleTaskData
	TitleTaskData() *TitleTaskDataConfig

	GetTeachChapterData(key uint64) *teach.TeachChapterData
	GetTeachChapterDataArray() []*teach.TeachChapterData
	TeachChapterData() *TeachChapterDataConfig

	GetSecretTowerData(key uint64) *towerdata.SecretTowerData
	GetSecretTowerDataArray() []*towerdata.SecretTowerData
	SecretTowerData() *SecretTowerDataConfig

	GetSecretTowerWordsData(key uint64) *towerdata.SecretTowerWordsData
	GetSecretTowerWordsDataArray() []*towerdata.SecretTowerWordsData
	SecretTowerWordsData() *SecretTowerWordsDataConfig

	GetTowerData(key uint64) *towerdata.TowerData
	GetTowerDataArray() []*towerdata.TowerData
	TowerData() *TowerDataConfig

	GetVipContinueDaysData(key uint64) *vip.VipContinueDaysData
	GetVipContinueDaysDataArray() []*vip.VipContinueDaysData
	VipContinueDaysData() *VipContinueDaysDataConfig

	GetVipLevelData(key uint64) *vip.VipLevelData
	GetVipLevelDataArray() []*vip.VipLevelData
	VipLevelData() *VipLevelDataConfig

	GetResistXiongNuData(key uint64) *xiongnu.ResistXiongNuData
	GetResistXiongNuDataArray() []*xiongnu.ResistXiongNuData
	ResistXiongNuData() *ResistXiongNuDataConfig

	GetResistXiongNuScoreData(key uint64) *xiongnu.ResistXiongNuScoreData
	GetResistXiongNuScoreDataArray() []*xiongnu.ResistXiongNuScoreData
	ResistXiongNuScoreData() *ResistXiongNuScoreDataConfig

	GetResistXiongNuWaveData(key uint64) *xiongnu.ResistXiongNuWaveData
	GetResistXiongNuWaveDataArray() []*xiongnu.ResistXiongNuWaveData
	ResistXiongNuWaveData() *ResistXiongNuWaveDataConfig

	GetXuanyuanRangeData(key uint64) *xuanydata.XuanyuanRangeData
	GetXuanyuanRangeDataArray() []*xuanydata.XuanyuanRangeData
	XuanyuanRangeData() *XuanyuanRangeDataConfig

	GetXuanyuanRankPrizeData(key uint64) *xuanydata.XuanyuanRankPrizeData
	GetXuanyuanRankPrizeDataArray() []*xuanydata.XuanyuanRankPrizeData
	XuanyuanRankPrizeData() *XuanyuanRankPrizeDataConfig

	GetZhanJiangChapterData(key uint64) *zhanjiang.ZhanJiangChapterData
	GetZhanJiangChapterDataArray() []*zhanjiang.ZhanJiangChapterData
	ZhanJiangChapterData() *ZhanJiangChapterDataConfig

	GetZhanJiangData(key uint64) *zhanjiang.ZhanJiangData
	GetZhanJiangDataArray() []*zhanjiang.ZhanJiangData
	ZhanJiangData() *ZhanJiangDataConfig

	GetZhanJiangGuanQiaData(key uint64) *zhanjiang.ZhanJiangGuanQiaData
	GetZhanJiangGuanQiaDataArray() []*zhanjiang.ZhanJiangGuanQiaData
	ZhanJiangGuanQiaData() *ZhanJiangGuanQiaDataConfig

	GetZhengWuCompleteData(key uint64) *zhengwu.ZhengWuCompleteData
	GetZhengWuCompleteDataArray() []*zhengwu.ZhengWuCompleteData
	ZhengWuCompleteData() *ZhengWuCompleteDataConfig

	GetZhengWuData(key uint64) *zhengwu.ZhengWuData
	GetZhengWuDataArray() []*zhengwu.ZhengWuData
	ZhengWuData() *ZhengWuDataConfig

	GetZhengWuRefreshData(key uint64) *zhengwu.ZhengWuRefreshData
	GetZhengWuRefreshDataArray() []*zhengwu.ZhengWuRefreshData
	ZhengWuRefreshData() *ZhengWuRefreshDataConfig

	BaiZhanMiscData() *bai_zhan_data.BaiZhanMiscData

	CombatConfig() *combatdata.CombatConfig

	CombatMiscConfig() *combatdata.CombatMiscConfig

	EquipCombineDatas() *combine.EquipCombineDatas

	CountryMiscData() *country.CountryMiscData

	BroadcastHelp() *data.BroadcastHelp

	TextHelp() *data.TextHelp

	ExchangeMiscData() *dianquan.ExchangeMiscData

	BuildingLayoutMiscData() *domestic_data.BuildingLayoutMiscData

	CityEventMiscData() *domestic_data.CityEventMiscData

	MainCityMiscData() *domestic_data.MainCityMiscData

	DungeonMiscData() *dungeon.DungeonMiscData

	FarmMiscConfig() *farm.FarmMiscConfig

	FishRandomer() *fishing_data.FishRandomer

	GardenConfig() *gardendata.GardenConfig

	EquipmentTaozConfig() *goods.EquipmentTaozConfig

	GemDatas() *goods.GemDatas

	GoodsCheck() *goods.GoodsCheck

	GuildLogHelp() *guild_data.GuildLogHelp

	NpcGuildSuffixName() *guild_data.NpcGuildSuffixName

	HebiMiscData() *hebi.HebiMiscData

	HeroCreateData() *heroinit.HeroCreateData

	HeroInitData() *heroinit.HeroInitData

	MailHelp() *maildata.MailHelp

	JiuGuanMiscData() *military_data.JiuGuanMiscData

	JunYingMiscData() *military_data.JunYingMiscData

	McBuildMiscData() *mingcdata.McBuildMiscData

	MingcMiscData() *mingcdata.MingcMiscData

	EventLimitGiftConfig() *promdata.EventLimitGiftConfig

	PromotionMiscData() *promdata.PromotionMiscData

	QuestionMiscData() *question.QuestionMiscData

	RaceConfig() *race.RaceConfig

	RandomEventDataDictionary() *random_event.RandomEventDataDictionary

	RandomEventPositionDictionary() *random_event.RandomEventPositionDictionary

	RankMiscData() *rank_data.RankMiscData

	JunTuanNpcPlaceConfig() *regdata.JunTuanNpcPlaceConfig

	SeasonMiscData() *season.SeasonMiscData

	SettingMiscData() *settings.SettingMiscData

	ShopMiscData() *shop.ShopMiscData

	GoodsConfig() *singleton.GoodsConfig

	GuildConfig() *singleton.GuildConfig

	GuildGenConfig() *singleton.GuildGenConfig

	MilitaryConfig() *singleton.MilitaryConfig

	MiscConfig() *singleton.MiscConfig

	MiscGenConfig() *singleton.MiscGenConfig

	RegionConfig() *singleton.RegionConfig

	RegionGenConfig() *singleton.RegionGenConfig

	TagMiscData() *tag.TagMiscData

	TaskMiscData() *taskdata.TaskMiscData

	SecretTowerMiscData() *towerdata.SecretTowerMiscData

	VipMiscData() *vip.VipMiscData

	ResistXiongNuMisc() *xiongnu.ResistXiongNuMisc

	XuanyuanMiscData() *xuanydata.XuanyuanMiscData

	ZhanJiangMiscData() *zhanjiang.ZhanJiangMiscData

	ZhengWuMiscData() *zhengwu.ZhengWuMiscData

	ZhengWuRandomData() *zhengwu.ZhengWuRandomData
}

func intSlice(keys []int) sort.IntSlice {
	return sort.IntSlice(keys)
}

func stringSlice(keys []string) sort.StringSlice {
	return sort.StringSlice(keys)
}

type uint64Slice []uint64

func (a uint64Slice) Len() int           { return len(a) }
func (a uint64Slice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a uint64Slice) Less(i, j int) bool { return a[i] < a[j] }
