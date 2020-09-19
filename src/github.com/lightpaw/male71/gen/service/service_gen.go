package service

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/kv"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/module"
	"github.com/lightpaw/male7/module/activity"
	"github.com/lightpaw/male7/module/bai_zhan"
	"github.com/lightpaw/male7/module/bai_zhan/bai_zhan_service"
	chat1 "github.com/lightpaw/male7/module/chat"
	"github.com/lightpaw/male7/module/client_config"
	country1 "github.com/lightpaw/male7/module/country"
	"github.com/lightpaw/male7/module/depot"
	"github.com/lightpaw/male7/module/dianquan"
	"github.com/lightpaw/male7/module/domestic"
	"github.com/lightpaw/male7/module/dungeon"
	"github.com/lightpaw/male7/module/equipment"
	"github.com/lightpaw/male7/module/farm"
	"github.com/lightpaw/male7/module/fishing"
	"github.com/lightpaw/male7/module/garden"
	"github.com/lightpaw/male7/module/gem"
	"github.com/lightpaw/male7/module/gm"
	"github.com/lightpaw/male7/module/guild"
	"github.com/lightpaw/male7/module/hebi"
	"github.com/lightpaw/male7/module/mail"
	"github.com/lightpaw/male7/module/military"
	mingc1 "github.com/lightpaw/male7/module/mingc"
	mingc_war1 "github.com/lightpaw/male7/module/mingc_war"
	"github.com/lightpaw/male7/module/misc"
	"github.com/lightpaw/male7/module/promotion"
	"github.com/lightpaw/male7/module/question"
	"github.com/lightpaw/male7/module/random_event"
	"github.com/lightpaw/male7/module/rank"
	"github.com/lightpaw/male7/module/realm"
	red_packet1 "github.com/lightpaw/male7/module/red_packet"
	"github.com/lightpaw/male7/module/region"
	"github.com/lightpaw/male7/module/relation"
	"github.com/lightpaw/male7/module/secret_tower"
	"github.com/lightpaw/male7/module/shop"
	"github.com/lightpaw/male7/module/strategy"
	"github.com/lightpaw/male7/module/stress"
	"github.com/lightpaw/male7/module/survey"
	"github.com/lightpaw/male7/module/tag"
	"github.com/lightpaw/male7/module/task"
	"github.com/lightpaw/male7/module/teach"
	"github.com/lightpaw/male7/module/tower"
	"github.com/lightpaw/male7/module/vip"
	"github.com/lightpaw/male7/module/xiongnu"
	"github.com/lightpaw/male7/module/xiongnu/xionnuservice"
	"github.com/lightpaw/male7/module/xuanyuan"
	"github.com/lightpaw/male7/module/zhanjiang"
	"github.com/lightpaw/male7/module/zhengwu"
	"github.com/lightpaw/male7/service"
	"github.com/lightpaw/male7/service/aws"
	"github.com/lightpaw/male7/service/broadcast"
	"github.com/lightpaw/male7/service/buff"
	"github.com/lightpaw/male7/service/chat"
	"github.com/lightpaw/male7/service/cluster"
	"github.com/lightpaw/male7/service/conflict/guilddataservice"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshot"
	"github.com/lightpaw/male7/service/conflict/heroservice"
	"github.com/lightpaw/male7/service/country"
	"github.com/lightpaw/male7/service/db"
	"github.com/lightpaw/male7/service/extratimesservice"
	"github.com/lightpaw/male7/service/fight"
	"github.com/lightpaw/male7/service/herosnapshot"
	"github.com/lightpaw/male7/service/httphandler"
	"github.com/lightpaw/male7/service/kafka"
	"github.com/lightpaw/male7/service/mingc"
	"github.com/lightpaw/male7/service/mingc_war"
	"github.com/lightpaw/male7/service/monitor"
	"github.com/lightpaw/male7/service/product"
	"github.com/lightpaw/male7/service/push"
	"github.com/lightpaw/male7/service/red_packet"
	"github.com/lightpaw/male7/service/reminder"
	"github.com/lightpaw/male7/service/seasonservice"
	"github.com/lightpaw/male7/service/status"
	"github.com/lightpaw/male7/service/ticker"
	"github.com/lightpaw/male7/service/timelimitgift"
	"github.com/lightpaw/male7/service/timeservice"
	"github.com/lightpaw/male7/service/tss/tssclient"
	"github.com/lightpaw/male7/service/world"
	"github.com/lightpaw/male7/tlog"
)

// package private的实例
var (
	_                          = logrus.DebugLevel // prevent unused import
	individualServerConfig     *kv.IndividualServerConfig
	guildSnapshotService       *guildsnapshot.GuildSnapshotService
	configDatas                *config.ConfigDatas
	baiZhanService             *bai_zhan_service.BaiZhanService
	timeService                *timeservice.TimeService
	tickerService              *ticker.TickerService
	stressModule               *stress.StressModule
	seasonService              *seasonservice.SeasonService
	metricsRegister            *monitor.MetricsRegister
	locationHeroCache          *relation.LocationHeroCache
	gameServer                 *cluster.GameServer
	extraTimesService          *extratimesservice.ExtraTimesService
	dbService                  *db.DbService
	clientConfigModule         *client_config.ClientConfigModule
	awsService                 *aws.AwsService
	xiongNuService             *xionnuservice.XiongNuService
	worldService               *world.WorldService
	tssClient                  *tssclient.TssClient
	serverStartStopTimeService *status.ServerStartStopTimeService
	reminderService            *reminder.ReminderService
	kafkaService               *kafka.KafkaService
	heroSnapshotService        *herosnapshot.HeroSnapshotService
	heroDataService            *heroservice.HeroDataService
	gameExporter               *monitor.GameExporter
	farmService                *farm.FarmService
	clusterService             *cluster.ClusterService
	buffService                *buff.BuffService
	tlogBaseService            *tlog.TlogBaseService
	tagModule                  *tag.TagModule
	redPacketService           *red_packet.RedPacketService
	pushService                *push.PushService
	chatService                *chat.ChatService
	broadcastService           *broadcast.BroadcastService
	tlogService                *tlog.TlogService
	mailModule                 *mail.MailModule
	fightXService              *fight.FightXService
	fightService               *fight.FightService
	countryService             *country.CountryService
	guildService               *guilddataservice.GuildService
	surveyModule               *survey.SurveyModule
	rankModule                 *rank.RankModule
	mingcService               *mingc.MingcService
	serviceDep                 *service.ServiceDep
	secretTowerModule          *secret_tower.SecretTowerModule
	relationModule             *relation.RelationModule
	redPacketModule            *red_packet1.RedPacketModule
	realmService               *realm.RealmService
	randomEventModule          *random_event.RandomEventModule
	questionModule             *question.QuestionModule
	productService             *product.ProductService
	miscModule                 *misc.MiscModule
	mingcWarService            *mingc_war.MingcWarService
	mingcWarModule             *mingc_war1.MingcWarModule
	mingcModule                *mingc1.MingcModule
	militaryModule             *military.MilitaryModule
	lightpawHandler            *httphandler.LightpawHandler
	hebiModule                 *hebi.HebiModule
	gemModule                  *gem.GemModule
	gardenModule               *garden.GardenModule
	fishingModule              *fishing.FishingModule
	farmModule                 *farm.FarmModule
	equipmentModule            *equipment.EquipmentModule
	dungeonModule              *dungeon.DungeonModule
	dianquanModule             *dianquan.DianquanModule
	depotModule                *depot.DepotModule
	countryModule              *country1.CountryModule
	chatModule                 *chat1.ChatModule
	baiZhanModule              *bai_zhan.BaiZhanModule
	activityModule             *activity.ActivityModule
	zhengWuModule              *zhengwu.ZhengWuModule
	zhanJiangModule            *zhanjiang.ZhanJiangModule
	xuanyuanModule             *xuanyuan.XuanyuanModule
	xiongNuModule              *xiongnu.XiongNuModule
	vipModule                  *vip.VipModule
	towerModule                *tower.TowerModule
	timeLimitGiftService       *timelimitgift.TimeLimitGiftService
	teachModule                *teach.TeachModule
	taskModule                 *task.TaskModule
	strategyModule             *strategy.StrategyModule
	regionModule               *region.RegionModule
	guildModule                *guild.GuildModule
	domesticModule             *domestic.DomesticModule
	shopModule                 *shop.ShopModule
	promotionModule            *promotion.PromotionModule
	modules                    *module.Modules
	gmModule                   *gm.GmModule
)

// 外部可见的每个实例的interface, 里面不包含所有的处理消息的方法
var (
	IndividualServerConfig     iface.IndividualServerConfig
	GuildSnapshotService       iface.GuildSnapshotService
	ConfigDatas                iface.ConfigDatas
	BaiZhanService             iface.BaiZhanService
	TimeService                iface.TimeService
	TickerService              iface.TickerService
	StressModule               iface.StressModule
	SeasonService              iface.SeasonService
	MetricsRegister            iface.MetricsRegister
	LocationHeroCache          iface.LocationHeroCache
	GameServer                 iface.GameServer
	ExtraTimesService          iface.ExtraTimesService
	DbService                  iface.DbService
	ClientConfigModule         iface.ClientConfigModule
	AwsService                 iface.AwsService
	XiongNuService             iface.XiongNuService
	WorldService               iface.WorldService
	TssClient                  iface.TssClient
	ServerStartStopTimeService iface.ServerStartStopTimeService
	ReminderService            iface.ReminderService
	KafkaService               iface.KafkaService
	HeroSnapshotService        iface.HeroSnapshotService
	HeroDataService            iface.HeroDataService
	GameExporter               iface.GameExporter
	FarmService                iface.FarmService
	ClusterService             iface.ClusterService
	BuffService                iface.BuffService
	TlogBaseService            iface.TlogBaseService
	TagModule                  iface.TagModule
	RedPacketService           iface.RedPacketService
	PushService                iface.PushService
	ChatService                iface.ChatService
	BroadcastService           iface.BroadcastService
	TlogService                iface.TlogService
	MailModule                 iface.MailModule
	FightXService              iface.FightXService
	FightService               iface.FightService
	CountryService             iface.CountryService
	GuildService               iface.GuildService
	SurveyModule               iface.SurveyModule
	RankModule                 iface.RankModule
	MingcService               iface.MingcService
	ServiceDep                 iface.ServiceDep
	SecretTowerModule          iface.SecretTowerModule
	RelationModule             iface.RelationModule
	RedPacketModule            iface.RedPacketModule
	RealmService               iface.RealmService
	RandomEventModule          iface.RandomEventModule
	QuestionModule             iface.QuestionModule
	ProductService             iface.ProductService
	MiscModule                 iface.MiscModule
	MingcWarService            iface.MingcWarService
	MingcWarModule             iface.MingcWarModule
	MingcModule                iface.MingcModule
	MilitaryModule             iface.MilitaryModule
	LightpawHandler            iface.LightpawHandler
	HebiModule                 iface.HebiModule
	GemModule                  iface.GemModule
	GardenModule               iface.GardenModule
	FishingModule              iface.FishingModule
	FarmModule                 iface.FarmModule
	EquipmentModule            iface.EquipmentModule
	DungeonModule              iface.DungeonModule
	DianquanModule             iface.DianquanModule
	DepotModule                iface.DepotModule
	CountryModule              iface.CountryModule
	ChatModule                 iface.ChatModule
	BaiZhanModule              iface.BaiZhanModule
	ActivityModule             iface.ActivityModule
	ZhengWuModule              iface.ZhengWuModule
	ZhanJiangModule            iface.ZhanJiangModule
	XuanyuanModule             iface.XuanyuanModule
	XiongNuModule              iface.XiongNuModule
	VipModule                  iface.VipModule
	TowerModule                iface.TowerModule
	TimeLimitGiftService       iface.TimeLimitGiftService
	TeachModule                iface.TeachModule
	TaskModule                 iface.TaskModule
	StrategyModule             iface.StrategyModule
	RegionModule               iface.RegionModule
	GuildModule                iface.GuildModule
	DomesticModule             iface.DomesticModule
	ShopModule                 iface.ShopModule
	PromotionModule            iface.PromotionModule
	Modules                    iface.Modules
	GmModule                   iface.GmModule
)

func Init() {
	logrus.Debug("正在初始化 IndividualServerConfig")
	individualServerConfig = kv.NewIndividualServerConfig()
	logrus.Debug("初始化完成 IndividualServerConfig")
	logrus.Debug("正在初始化 GuildSnapshotService")
	guildSnapshotService = guildsnapshot.NewGuildSnapshotService()
	logrus.Debug("初始化完成 GuildSnapshotService")
	logrus.Debug("正在初始化 ConfigDatas")
	configDatas = config.NewConfigDatas()
	logrus.Debug("初始化完成 ConfigDatas")
	logrus.Debug("正在初始化 BaiZhanService")
	baiZhanService = bai_zhan_service.NewBaiZhanService()
	logrus.Debug("初始化完成 BaiZhanService")
	logrus.Debug("正在初始化 TimeService")
	timeService = timeservice.NewTimeService(individualServerConfig)
	logrus.Debug("初始化完成 TimeService")
	logrus.Debug("正在初始化 TickerService")
	tickerService = ticker.NewTickerService(timeService, configDatas)
	logrus.Debug("初始化完成 TickerService")
	logrus.Debug("正在初始化 StressModule")
	stressModule = stress.NewStressModule(individualServerConfig)
	logrus.Debug("初始化完成 StressModule")
	logrus.Debug("正在初始化 SeasonService")
	seasonService = seasonservice.NewSeasonService(timeService, individualServerConfig, tickerService, configDatas)
	logrus.Debug("初始化完成 SeasonService")
	logrus.Debug("正在初始化 MetricsRegister")
	metricsRegister = monitor.NewMetricsRegister(individualServerConfig)
	logrus.Debug("初始化完成 MetricsRegister")
	logrus.Debug("正在初始化 LocationHeroCache")
	locationHeroCache = relation.NewLocationHeroCache(timeService, configDatas)
	logrus.Debug("初始化完成 LocationHeroCache")
	logrus.Debug("正在初始化 GameServer")
	gameServer = cluster.NewGameServer(individualServerConfig)
	logrus.Debug("初始化完成 GameServer")
	logrus.Debug("正在初始化 ExtraTimesService")
	extraTimesService = extratimesservice.NewExtraTimesService(seasonService, timeService)
	logrus.Debug("初始化完成 ExtraTimesService")
	logrus.Debug("正在初始化 DbService")
	dbService = db.NewDbService(configDatas, individualServerConfig, timeService, metricsRegister)
	logrus.Debug("初始化完成 DbService")
	logrus.Debug("正在初始化 ClientConfigModule")
	clientConfigModule = client_config.NewClientConfigModule(individualServerConfig)
	logrus.Debug("初始化完成 ClientConfigModule")
	logrus.Debug("正在初始化 AwsService")
	awsService = aws.NewAwsService(individualServerConfig)
	logrus.Debug("初始化完成 AwsService")
	logrus.Debug("正在初始化 XiongNuService")
	xiongNuService = xionnuservice.NewXiongNuService(timeService, configDatas, tickerService, dbService)
	logrus.Debug("初始化完成 XiongNuService")
	logrus.Debug("正在初始化 WorldService")
	worldService = world.NewWorldService(metricsRegister, individualServerConfig, tickerService)
	logrus.Debug("初始化完成 WorldService")
	logrus.Debug("正在初始化 TssClient")
	tssClient = tssclient.NewTssClient(individualServerConfig, gameServer)
	logrus.Debug("初始化完成 TssClient")
	logrus.Debug("正在初始化 ServerStartStopTimeService")
	serverStartStopTimeService = status.NewServerStartStopTimeService(dbService, timeService)
	logrus.Debug("初始化完成 ServerStartStopTimeService")
	logrus.Debug("正在初始化 ReminderService")
	reminderService = reminder.NewReminderService(guildSnapshotService, worldService)
	logrus.Debug("初始化完成 ReminderService")
	logrus.Debug("正在初始化 KafkaService")
	kafkaService = kafka.NewKafkaService(individualServerConfig, dbService, timeService)
	logrus.Debug("初始化完成 KafkaService")
	logrus.Debug("正在初始化 HeroSnapshotService")
	heroSnapshotService = herosnapshot.NewHeroSnapshotService(dbService, configDatas, guildSnapshotService, baiZhanService, individualServerConfig)
	logrus.Debug("初始化完成 HeroSnapshotService")
	logrus.Debug("正在初始化 HeroDataService")
	heroDataService = heroservice.NewHeroDataService(individualServerConfig, dbService, heroSnapshotService, worldService)
	logrus.Debug("初始化完成 HeroDataService")
	logrus.Debug("正在初始化 GameExporter")
	gameExporter = monitor.NewGameExporter(metricsRegister, dbService)
	logrus.Debug("初始化完成 GameExporter")
	logrus.Debug("正在初始化 FarmService")
	farmService = farm.NewFarmService(dbService, worldService, configDatas, timeService, tickerService)
	logrus.Debug("初始化完成 FarmService")
	logrus.Debug("正在初始化 ClusterService")
	clusterService = cluster.NewClusterService(individualServerConfig, gameServer, worldService)
	logrus.Debug("初始化完成 ClusterService")
	logrus.Debug("正在初始化 BuffService")
	buffService = buff.NewBuffService(heroDataService, worldService, configDatas, timeService, farmService)
	logrus.Debug("初始化完成 BuffService")
	logrus.Debug("正在初始化 TlogBaseService")
	tlogBaseService = tlog.NewTlogBaseService(timeService, kafkaService, individualServerConfig)
	logrus.Debug("初始化完成 TlogBaseService")
	logrus.Debug("正在初始化 TagModule")
	tagModule = tag.NewTagModule(configDatas, timeService, heroDataService, guildSnapshotService, worldService)
	logrus.Debug("初始化完成 TagModule")
	logrus.Debug("正在初始化 RedPacketService")
	redPacketService = red_packet.NewRedPacketService(configDatas, timeService, heroSnapshotService, dbService)
	logrus.Debug("初始化完成 RedPacketService")
	logrus.Debug("正在初始化 PushService")
	pushService = push.NewPushService(configDatas, individualServerConfig, worldService, dbService, heroSnapshotService, clusterService)
	logrus.Debug("初始化完成 PushService")
	logrus.Debug("正在初始化 ChatService")
	chatService = chat.NewChatService(dbService, pushService, timeService, guildSnapshotService, heroSnapshotService, worldService, configDatas, redPacketService)
	logrus.Debug("初始化完成 ChatService")
	logrus.Debug("正在初始化 BroadcastService")
	broadcastService = broadcast.NewBroadcastService(guildSnapshotService, worldService, chatService, heroSnapshotService, configDatas, timeService)
	logrus.Debug("初始化完成 BroadcastService")
	logrus.Debug("正在初始化 TlogService")
	tlogService = tlog.NewTlogService(worldService, heroSnapshotService, tlogBaseService, timeService, individualServerConfig)
	logrus.Debug("初始化完成 TlogService")
	logrus.Debug("正在初始化 MailModule")
	mailModule = mail.NewMailModule(timeService, worldService, broadcastService, configDatas, tlogService, heroDataService, guildSnapshotService, dbService, chatService)
	logrus.Debug("初始化完成 MailModule")
	logrus.Debug("正在初始化 FightXService")
	fightXService = fight.NewFightXService(configDatas, heroDataService, individualServerConfig, clusterService, tlogService)
	logrus.Debug("初始化完成 FightXService")
	logrus.Debug("正在初始化 FightService")
	fightService = fight.NewFightService(configDatas, individualServerConfig, heroDataService, individualServerConfig, clusterService, tlogService)
	logrus.Debug("初始化完成 FightService")
	logrus.Debug("正在初始化 CountryService")
	countryService = country.NewCountryService(configDatas, dbService, tickerService, timeService, heroSnapshotService, heroDataService, buffService, worldService, mailModule, broadcastService)
	logrus.Debug("初始化完成 CountryService")
	logrus.Debug("正在初始化 GuildService")
	guildService = guilddataservice.NewGuildService(configDatas, timeService, dbService, countryService, heroDataService, heroSnapshotService, xiongNuService, guildSnapshotService, worldService, chatService)
	logrus.Debug("初始化完成 GuildService")
	logrus.Debug("正在初始化 SurveyModule")
	surveyModule = survey.NewSurveyModule(configDatas, heroDataService, timeService, guildService, mailModule)
	logrus.Debug("初始化完成 SurveyModule")
	logrus.Debug("正在初始化 RankModule")
	rankModule = rank.NewRankModule(configDatas, dbService, guildService, heroSnapshotService, timeService, serverStartStopTimeService, baiZhanService, tssClient)
	logrus.Debug("初始化完成 RankModule")
	logrus.Debug("正在初始化 MingcService")
	mingcService = mingc.NewMingcService(configDatas, dbService, tickerService, guildService, guildSnapshotService, countryService, timeService, individualServerConfig, mailModule, chatService, heroDataService, broadcastService, worldService, tlogService)
	logrus.Debug("初始化完成 MingcService")
	logrus.Debug("正在初始化 ServiceDep")
	serviceDep = service.NewServiceDep(configDatas, individualServerConfig, guildService, guildSnapshotService, worldService, broadcastService, heroSnapshotService, heroDataService, timeService, pushService, chatService, dbService, countryService, mingcService, fightService, fightXService, tlogService, mailModule)
	logrus.Debug("初始化完成 ServiceDep")
	logrus.Debug("正在初始化 SecretTowerModule")
	secretTowerModule = secret_tower.NewSecretTowerModule(serviceDep, seasonService, fightXService)
	logrus.Debug("初始化完成 SecretTowerModule")
	logrus.Debug("正在初始化 RelationModule")
	relationModule = relation.NewRelationModule(serviceDep, chatService, tssClient, locationHeroCache)
	logrus.Debug("初始化完成 RelationModule")
	logrus.Debug("正在初始化 RedPacketModule")
	redPacketModule = red_packet1.NewRedPacketModule(serviceDep, redPacketService)
	logrus.Debug("初始化完成 RedPacketModule")
	logrus.Debug("正在初始化 RealmService")
	realmService = realm.NewRealmService(serviceDep, seasonService, tlogService, fightXService, mailModule, extraTimesService, dbService, tickerService, pushService, reminderService, xiongNuService, farmService, buffService, baiZhanService)
	logrus.Debug("初始化完成 RealmService")
	logrus.Debug("正在初始化 RandomEventModule")
	randomEventModule = random_event.NewRandomEventModule(serviceDep, seasonService)
	logrus.Debug("初始化完成 RandomEventModule")
	logrus.Debug("正在初始化 QuestionModule")
	questionModule = question.NewQuestionModule(serviceDep)
	logrus.Debug("初始化完成 QuestionModule")
	logrus.Debug("正在初始化 ProductService")
	productService = product.NewProductService(serviceDep, configDatas, dbService, timeService, heroDataService)
	logrus.Debug("初始化完成 ProductService")
	logrus.Debug("正在初始化 MiscModule")
	miscModule = misc.NewMiscModule(serviceDep, configDatas, timeService, individualServerConfig, clusterService, locationHeroCache)
	logrus.Debug("初始化完成 MiscModule")
	logrus.Debug("正在初始化 MingcWarService")
	mingcWarService = mingc_war.NewMingcWarService(serviceDep)
	logrus.Debug("初始化完成 MingcWarService")
	logrus.Debug("正在初始化 MingcWarModule")
	mingcWarModule = mingc_war1.NewMingcWarModule(serviceDep, mingcWarService, realmService)
	logrus.Debug("初始化完成 MingcWarModule")
	logrus.Debug("正在初始化 MingcModule")
	mingcModule = mingc1.NewMingcModule(serviceDep, baiZhanService)
	logrus.Debug("初始化完成 MingcModule")
	logrus.Debug("正在初始化 MilitaryModule")
	militaryModule = military.NewMilitaryModule(serviceDep, configDatas, individualServerConfig, fightService, fightXService, realmService, tssClient)
	logrus.Debug("初始化完成 MilitaryModule")
	logrus.Debug("正在初始化 LightpawHandler")
	lightpawHandler = httphandler.NewLightpawHandler(serviceDep, timeService, heroDataService, individualServerConfig)
	logrus.Debug("初始化完成 LightpawHandler")
	logrus.Debug("正在初始化 HebiModule")
	hebiModule = hebi.NewHebiModule(serviceDep, mailModule, fightService, tickerService)
	logrus.Debug("初始化完成 HebiModule")
	logrus.Debug("正在初始化 GemModule")
	gemModule = gem.NewGemModule(serviceDep)
	logrus.Debug("初始化完成 GemModule")
	logrus.Debug("正在初始化 GardenModule")
	gardenModule = garden.NewGardenModule(serviceDep, individualServerConfig, seasonService)
	logrus.Debug("初始化完成 GardenModule")
	logrus.Debug("正在初始化 FishingModule")
	fishingModule = fishing.NewFishingModule(serviceDep, guildSnapshotService)
	logrus.Debug("初始化完成 FishingModule")
	logrus.Debug("正在初始化 FarmModule")
	farmModule = farm.NewFarmModule(serviceDep, farmService, dbService, worldService, configDatas, timeService, heroDataService, heroSnapshotService, guildSnapshotService, tickerService, seasonService)
	logrus.Debug("初始化完成 FarmModule")
	logrus.Debug("正在初始化 EquipmentModule")
	equipmentModule = equipment.NewEquipmentModule(serviceDep)
	logrus.Debug("初始化完成 EquipmentModule")
	logrus.Debug("正在初始化 DungeonModule")
	dungeonModule = dungeon.NewDungeonModule(serviceDep, fightXService, guildSnapshotService)
	logrus.Debug("初始化完成 DungeonModule")
	logrus.Debug("正在初始化 DianquanModule")
	dianquanModule = dianquan.NewDianquanModule(serviceDep, configDatas)
	logrus.Debug("初始化完成 DianquanModule")
	logrus.Debug("正在初始化 DepotModule")
	depotModule = depot.NewDepotModule(serviceDep, rankModule)
	logrus.Debug("初始化完成 DepotModule")
	logrus.Debug("正在初始化 CountryModule")
	countryModule = country1.NewCountryModule(serviceDep, buffService, rankModule)
	logrus.Debug("初始化完成 CountryModule")
	logrus.Debug("正在初始化 ChatModule")
	chatModule = chat1.NewChatModule(serviceDep, dbService, pushService, tssClient, individualServerConfig, chatService, baiZhanService, mingcWarService)
	logrus.Debug("初始化完成 ChatModule")
	logrus.Debug("正在初始化 BaiZhanModule")
	baiZhanModule = bai_zhan.NewBaiZhanModule(serviceDep, dbService, fightXService, tickerService, guildSnapshotService, baiZhanService, rankModule, mailModule)
	logrus.Debug("初始化完成 BaiZhanModule")
	logrus.Debug("正在初始化 ActivityModule")
	activityModule = activity.NewActivityModule(serviceDep, tickerService)
	logrus.Debug("初始化完成 ActivityModule")
	logrus.Debug("正在初始化 ZhengWuModule")
	zhengWuModule = zhengwu.NewZhengWuModule(serviceDep)
	logrus.Debug("初始化完成 ZhengWuModule")
	logrus.Debug("正在初始化 ZhanJiangModule")
	zhanJiangModule = zhanjiang.NewZhanJiangModule(serviceDep, fightService, militaryModule, guildSnapshotService)
	logrus.Debug("初始化完成 ZhanJiangModule")
	logrus.Debug("正在初始化 XuanyuanModule")
	xuanyuanModule = xuanyuan.NewXuanyuanModule(configDatas, timeService, heroDataService, heroSnapshotService, guildSnapshotService, fightXService, dbService, tickerService, worldService, serviceDep, rankModule)
	logrus.Debug("初始化完成 XuanyuanModule")
	logrus.Debug("正在初始化 XiongNuModule")
	xiongNuModule = xiongnu.NewXiongNuModule(serviceDep, configDatas, realmService, mailModule, tickerService, xiongNuService, pushService, rankModule)
	logrus.Debug("初始化完成 XiongNuModule")
	logrus.Debug("正在初始化 VipModule")
	vipModule = vip.NewVipModule(serviceDep)
	logrus.Debug("初始化完成 VipModule")
	logrus.Debug("正在初始化 TowerModule")
	towerModule = tower.NewTowerModule(serviceDep, fightXService, dbService, rankModule, xuanyuanModule, mailModule)
	logrus.Debug("初始化完成 TowerModule")
	logrus.Debug("正在初始化 TimeLimitGiftService")
	timeLimitGiftService = timelimitgift.NewTimeLimitGiftService(serviceDep)
	logrus.Debug("初始化完成 TimeLimitGiftService")
	logrus.Debug("正在初始化 TeachModule")
	teachModule = teach.NewTeachModule(serviceDep)
	logrus.Debug("初始化完成 TeachModule")
	logrus.Debug("正在初始化 TaskModule")
	taskModule = task.NewTaskModule(serviceDep, guildSnapshotService, rankModule, realmService)
	logrus.Debug("初始化完成 TaskModule")
	logrus.Debug("正在初始化 StrategyModule")
	strategyModule = strategy.NewStrategyModule(serviceDep, configDatas, realmService, farmService)
	logrus.Debug("初始化完成 StrategyModule")
	logrus.Debug("正在初始化 RegionModule")
	regionModule = region.NewRegionModule(serviceDep, configDatas, timeService, tlogService, worldService, heroDataService, extraTimesService, realmService, guildSnapshotService, mingcWarService, heroSnapshotService, baiZhanService, pushService, rankModule, mailModule)
	logrus.Debug("初始化完成 RegionModule")
	logrus.Debug("正在初始化 GuildModule")
	guildModule = guild.NewGuildModule(serviceDep, configDatas, dbService, hebiModule, mailModule, pushService, realmService, xiongNuModule, chatService, xiongNuService, rankModule, tickerService, countryService, tssClient, baiZhanService, mingcWarService)
	logrus.Debug("初始化完成 GuildModule")
	logrus.Debug("正在初始化 DomesticModule")
	domesticModule = domestic.NewDomesticModule(serviceDep, seasonService, dbService, realmService, baiZhanService, tssClient, buffService, regionModule)
	logrus.Debug("初始化完成 DomesticModule")
	logrus.Debug("正在初始化 ShopModule")
	shopModule = shop.NewShopModule(serviceDep, guildModule)
	logrus.Debug("初始化完成 ShopModule")
	logrus.Debug("正在初始化 PromotionModule")
	promotionModule = promotion.NewPromotionModule(configDatas, timeService, serviceDep, timeLimitGiftService, guildModule)
	logrus.Debug("初始化完成 PromotionModule")
	logrus.Debug("正在初始化 Modules")
	modules = module.NewModules(activityModule, baiZhanModule, chatModule, clientConfigModule, countryModule, depotModule, dianquanModule, domesticModule, dungeonModule, equipmentModule, farmModule, fishingModule, gardenModule, gemModule, guildModule, hebiModule, mailModule, militaryModule, mingcModule, mingcWarModule, miscModule, promotionModule, questionModule, randomEventModule, rankModule, redPacketModule, regionModule, relationModule, secretTowerModule, shopModule, strategyModule, stressModule, surveyModule, tagModule, taskModule, teachModule, towerModule, vipModule, xiongNuModule, xuanyuanModule, zhanJiangModule, zhengWuModule)
	logrus.Debug("初始化完成 Modules")
	logrus.Debug("正在初始化 GmModule")
	gmModule = gm.NewGmModule(serviceDep, dbService, individualServerConfig, configDatas, modules, realmService, reminderService, buffService, pushService, farmService, mingcWarService, mingcService, clusterService, seasonService, gameExporter, countryService, tickerService)
	logrus.Debug("初始化完成 GmModule")
	IndividualServerConfig = individualServerConfig
	GuildSnapshotService = guildSnapshotService
	ConfigDatas = configDatas
	BaiZhanService = baiZhanService
	TimeService = timeService
	TickerService = tickerService
	StressModule = stressModule
	SeasonService = seasonService
	MetricsRegister = metricsRegister
	LocationHeroCache = locationHeroCache
	GameServer = gameServer
	ExtraTimesService = extraTimesService
	DbService = dbService
	ClientConfigModule = clientConfigModule
	AwsService = awsService
	XiongNuService = xiongNuService
	WorldService = worldService
	TssClient = tssClient
	ServerStartStopTimeService = serverStartStopTimeService
	ReminderService = reminderService
	KafkaService = kafkaService
	HeroSnapshotService = heroSnapshotService
	HeroDataService = heroDataService
	GameExporter = gameExporter
	FarmService = farmService
	ClusterService = clusterService
	BuffService = buffService
	TlogBaseService = tlogBaseService
	TagModule = tagModule
	RedPacketService = redPacketService
	PushService = pushService
	ChatService = chatService
	BroadcastService = broadcastService
	TlogService = tlogService
	MailModule = mailModule
	FightXService = fightXService
	FightService = fightService
	CountryService = countryService
	GuildService = guildService
	SurveyModule = surveyModule
	RankModule = rankModule
	MingcService = mingcService
	ServiceDep = serviceDep
	SecretTowerModule = secretTowerModule
	RelationModule = relationModule
	RedPacketModule = redPacketModule
	RealmService = realmService
	RandomEventModule = randomEventModule
	QuestionModule = questionModule
	ProductService = productService
	MiscModule = miscModule
	MingcWarService = mingcWarService
	MingcWarModule = mingcWarModule
	MingcModule = mingcModule
	MilitaryModule = militaryModule
	LightpawHandler = lightpawHandler
	HebiModule = hebiModule
	GemModule = gemModule
	GardenModule = gardenModule
	FishingModule = fishingModule
	FarmModule = farmModule
	EquipmentModule = equipmentModule
	DungeonModule = dungeonModule
	DianquanModule = dianquanModule
	DepotModule = depotModule
	CountryModule = countryModule
	ChatModule = chatModule
	BaiZhanModule = baiZhanModule
	ActivityModule = activityModule
	ZhengWuModule = zhengWuModule
	ZhanJiangModule = zhanJiangModule
	XuanyuanModule = xuanyuanModule
	XiongNuModule = xiongNuModule
	VipModule = vipModule
	TowerModule = towerModule
	TimeLimitGiftService = timeLimitGiftService
	TeachModule = teachModule
	TaskModule = taskModule
	StrategyModule = strategyModule
	RegionModule = regionModule
	GuildModule = guildModule
	DomesticModule = domesticModule
	ShopModule = shopModule
	PromotionModule = promotionModule
	Modules = modules
	GmModule = gmModule
}

// 有问题时, 注释上面的, 反注释下面的
/*
----------- 自动识别分割线 -----------
package service

import(
	"github.com/lightpaw/male7/gen/iface"
)

var(
        IndividualServerConfig iface.IndividualServerConfig
        GuildSnapshotService iface.GuildSnapshotService
        ConfigDatas iface.ConfigDatas
        BaiZhanService iface.BaiZhanService
        TimeService iface.TimeService
        TickerService iface.TickerService
        StressModule iface.StressModule
        SeasonService iface.SeasonService
        MetricsRegister iface.MetricsRegister
        LocationHeroCache iface.LocationHeroCache
        GameServer iface.GameServer
        ExtraTimesService iface.ExtraTimesService
        DbService iface.DbService
        ClientConfigModule iface.ClientConfigModule
        AwsService iface.AwsService
        XiongNuService iface.XiongNuService
        WorldService iface.WorldService
        TssClient iface.TssClient
        ServerStartStopTimeService iface.ServerStartStopTimeService
        ReminderService iface.ReminderService
        KafkaService iface.KafkaService
        HeroSnapshotService iface.HeroSnapshotService
        HeroDataService iface.HeroDataService
        GameExporter iface.GameExporter
        FarmService iface.FarmService
        ClusterService iface.ClusterService
        BuffService iface.BuffService
        TlogBaseService iface.TlogBaseService
        TagModule iface.TagModule
        RedPacketService iface.RedPacketService
        PushService iface.PushService
        ChatService iface.ChatService
        BroadcastService iface.BroadcastService
        TlogService iface.TlogService
        MailModule iface.MailModule
        FightXService iface.FightXService
        FightService iface.FightService
        CountryService iface.CountryService
        GuildService iface.GuildService
        SurveyModule iface.SurveyModule
        RankModule iface.RankModule
        MingcService iface.MingcService
        ServiceDep iface.ServiceDep
        SecretTowerModule iface.SecretTowerModule
        RelationModule iface.RelationModule
        RedPacketModule iface.RedPacketModule
        RealmService iface.RealmService
        RandomEventModule iface.RandomEventModule
        QuestionModule iface.QuestionModule
        ProductService iface.ProductService
        MiscModule iface.MiscModule
        MingcWarService iface.MingcWarService
        MingcWarModule iface.MingcWarModule
        MingcModule iface.MingcModule
        MilitaryModule iface.MilitaryModule
        LightpawHandler iface.LightpawHandler
        HebiModule iface.HebiModule
        GemModule iface.GemModule
        GardenModule iface.GardenModule
        FishingModule iface.FishingModule
        FarmModule iface.FarmModule
        EquipmentModule iface.EquipmentModule
        DungeonModule iface.DungeonModule
        DianquanModule iface.DianquanModule
        DepotModule iface.DepotModule
        CountryModule iface.CountryModule
        ChatModule iface.ChatModule
        BaiZhanModule iface.BaiZhanModule
        ActivityModule iface.ActivityModule
        ZhengWuModule iface.ZhengWuModule
        ZhanJiangModule iface.ZhanJiangModule
        XuanyuanModule iface.XuanyuanModule
        XiongNuModule iface.XiongNuModule
        VipModule iface.VipModule
        TowerModule iface.TowerModule
        TimeLimitGiftService iface.TimeLimitGiftService
        TeachModule iface.TeachModule
        TaskModule iface.TaskModule
        StrategyModule iface.StrategyModule
        RegionModule iface.RegionModule
        GuildModule iface.GuildModule
        DomesticModule iface.DomesticModule
        ShopModule iface.ShopModule
        PromotionModule iface.PromotionModule
        Modules iface.Modules
        GmModule iface.GmModule
)

func Init() {}
----------- 自动识别分割线 -----------
*/
