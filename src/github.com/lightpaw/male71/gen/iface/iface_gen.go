package iface

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/coreos/etcd/clientv3"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/basedata"
	"github.com/lightpaw/male7/config/blockdata"
	"github.com/lightpaw/male7/config/buffer"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/config/guild_data"
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/config/mingcdata"
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/pushdata"
	"github.com/lightpaw/male7/config/red_packet"
	"github.com/lightpaw/male7/config/regdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/config/season"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/face"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/module/bai_zhan/bai_zhan_objs"
	"github.com/lightpaw/male7/module/rank/rankface"
	"github.com/lightpaw/male7/module/realm/realmface"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuface"
	"github.com/lightpaw/male7/module/xiongnu/xiongnuinfo"
	"github.com/lightpaw/male7/pb/rpcpb/game2tss"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/guildsnapshotdata"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/db/isql"
	"github.com/lightpaw/male7/service/extratimesservice/extratimesface"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"github.com/lightpaw/male7/service/monitor/metrics"
	"github.com/lightpaw/male7/service/sender"
	"github.com/lightpaw/male7/service/ticker/tickdata"
	"github.com/lightpaw/male7/service/tss"
	"github.com/lightpaw/male7/util/atomic"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/i64/concurrent"
	"github.com/lightpaw/male7/util/msg"
	"github.com/lightpaw/pbutil"
	"github.com/lightpaw/rpc7"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

type ConfigDatas interface {
	config.Configs
	EncodeClient() *shared_proto.Config
}

type ActivityModule interface {
	Close()
	OnHeroOnline(HeroController)
}

type AwsService interface {
	InitFirehoseEventLog() bool
}

type BaiZhanModule interface {
	Challenge(HeroController) msg.ErrMsg
	ClearLastJunXian(HeroController)
	Close()
	CollectJunXianPrize(int32, HeroController) (bool, msg.ErrMsg)
	CollectSalary(HeroController) (bool, msg.ErrMsg)
	GmResetChallengeTimes(int64)
	GmResetDaily()
	GmSetJunXian(int64, HeroController)
	QueryBaiZhanInfo(HeroController)
	RequestRank(bool, uint64, HeroController)
	RequestSelfRank(HeroController)
	SelfRecord(int32, HeroController)
}

type BaiZhanService interface {
	Func(bai_zhan_objs.BaiZhanObjsFunc)
	GetBaiZhanObj(int64) bai_zhan_objs.RHeroBaiZhanObj
	GetHistoryMaxJunXianLevel(int64) uint64
	GetJunXianLevel(int64) uint64
	GetPoint(int64) uint64
	OnHeroOnline(HeroController)
	Stop(bai_zhan_objs.BaiZhanObjsFunc)
	TimeOutFunc(bai_zhan_objs.BaiZhanObjsFunc) bool
}

type BroadcastService interface {
	Broadcast(string, bool)
	Close()
	GetCaptainText(*entity.Captain) string
	GetEquipText(*goods.EquipmentData) string
}

type BuffService interface {

	// 新增或替换buff
	AddBuffToSelf(*data.BuffEffectData, int64) bool
	// 新增或替换buff
	AddBuffToTarget(*data.BuffEffectData, int64, int64) bool
	Cancel(int64, []*entity.BuffInfo) bool
	CancelGroup(int64, uint64) bool
	UpdatePerSecond(*entity.Hero, herolock.LockResult, *entity.BuffInfo)
}

type ChatModule interface {
}

type ChatService interface {
	AddMsg(*shared_proto.ChatMsgProto)
	BroadcastSystemChat(string)
	GetCacheMsg() pbutil.Buffer
	GetChatSender(int64) *shared_proto.ChatSenderProto
	GetSystemChatRecord(int64) pbutil.Buffer
	OnHeroOnline(HeroController)
	SaveDB(int64, int64, shared_proto.ChatType, []byte, []byte, *shared_proto.ChatMsgProto) int64
	// 系统自动聊天，有DB操作
	SysChat(int64, int64, shared_proto.ChatType, string, shared_proto.ChatMsgType, bool, bool, bool, bool)
	SysChatFunc(int64, int64, shared_proto.ChatType, string, shared_proto.ChatMsgType, bool, bool, bool, bool, constants.ChatFunc, constants.ChatFunc) int64
	SysChatProtoFunc(int64, int64, shared_proto.ChatType, string, shared_proto.ChatMsgType, bool, bool, bool, bool, constants.ChatFunc) int64
	SysChatSendFunc(int64, int64, shared_proto.ChatType, string, shared_proto.ChatMsgType, bool, bool, bool, bool, constants.ChatFunc) int64
	UpdateDBRedPacket(int64, bool) bool
}

type ClientConfigModule interface {
}

type ClusterService interface {
	Close()
	EtcdClient() *clientv3.Client
	GMUpdateClientVersion(string)
	GetClientVersionMsg(string, string) pbutil.Buffer
	GetConfig(string) ([]byte, error)
	LoginClient() *rpc7.Client
	RpcAddr() string
}

// 已在线的用户

type ConnectedUser interface {
	Disconnect(msg.ErrMsg)
	DisconnectAndWait(msg.ErrMsg)
	GetHeroController() HeroController
	Id() int64
	IsClosed() bool
	IsLoaded() bool
	LogoutType() uint64
	// 玩家杂项
	Misc() *server_proto.UserMiscProto
	// 玩家杂项
	MiscNeedOfflineSave() bool
	// 发送消息.
	Send(pbutil.Buffer)
	// 发送在线路繁忙时可以被丢掉的消息
	SendAll([]pbutil.Buffer)
	// 发送在线路繁忙时可以被丢掉的消息
	SendIfFree(pbutil.Buffer)
	SetHeroController(HeroController)
	SetLoaded()
	SetLogoutType(uint64)
	// 玩家杂项
	SetMisc(*server_proto.UserMiscProto)
	Sid() uint32
	TencentInfo() *shared_proto.TencentInfoProto
}

type CountryModule interface {
}

type CountryService interface {
	AddPrestige(uint64, uint64) (uint64, bool)
	AfterChangeCountry(int64, uint64, uint64, int32, int, bool)
	AfterChangeNameVoteStart(*entity.Country)
	AfterUpgradeTitle(int64, uint64, uint64)
	BroadcastCountry(pbutil.Buffer, uint64)
	CancelCountryDestroy(uint64)
	ChangeCountryHost(uint64, int64) bool
	ChangeKing(uint64, int64) bool
	Close()
	Countries() []*entity.Country
	CountriesMsg(uint64) pbutil.Buffer
	Country(uint64) *entity.Country
	CountryDestroy(uint64) bool
	CountryDetailMsg(uint64) pbutil.Buffer
	CountryFlagHeroName(int64) string
	CountryName(uint64) string
	CountryPrestigeMsg(uint64) pbutil.Buffer
	DetailMsgCacheDisable(uint64)
	ForceOfficialAppoint(uint64, int64, shared_proto.CountryOfficialType) bool
	ForceOfficialDepose(uint64, int64) bool
	GmAppointKing(uint64, int64)
	GmAppointOfficial(uint64, int64, shared_proto.CountryOfficialType)
	GmDeposeKing(uint64)
	GmDestroy(uint64)
	GmOfficialDeposeAll(uint64)
	HeroCountry(int64) uint64
	HeroOfficial(uint64, int64) shared_proto.CountryOfficialType
	IsCountryDestroyed(uint64) bool
	IsOnChangeNameVote(uint64) bool
	King(uint64) int64
	LockHeroCapital(int64) uint64
	LockHeroCountry(int64) uint64
	MsgCacheDisable(uint64)
	OfficialAppoint(int64, int64, uint64, shared_proto.CountryOfficialType, int32) (pbutil.Buffer, pbutil.Buffer)
	OfficialDepose(int64, int64, uint64) (shared_proto.CountryOfficialType, pbutil.Buffer, pbutil.Buffer)
	OfficialLeave(int64, uint64) (shared_proto.CountryOfficialType, pbutil.Buffer, pbutil.Buffer)
	OnHeroOnline(HeroController, uint64)
	ReducePrestige(uint64, uint64) (uint64, bool)
	TutorialCountriesProto() *shared_proto.CountriesProto
	UpdateMcWarMsg(*shared_proto.McWarProto)
	UpdateMingcsMsg(*shared_proto.MingcsProto)
	WalkCountryOnlineHero(uint64, CountryHeroWalker)
}

type DbService interface {
	AddChatMsg(context.Context, int64, []byte, *shared_proto.ChatMsgProto) (int64, error)
	AddFarmSteal(context.Context, int64, int64, cb.Cube) error
	AddMcWarGuildRecord(context.Context, uint64, uint64, int64, *shared_proto.McWarTroopsInfoProto) error
	AddMcWarHeroRecord(context.Context, uint64, uint64, int64, *shared_proto.McWarTroopAllRecordProto) error
	AddMcWarRecord(context.Context, uint64, uint64, *shared_proto.McWarFightRecordProto) error
	CallingTimes() uint64
	Close() error
	CreateFarmCube(context.Context, *entity.FarmCube) error
	CreateFarmLog(context.Context, *shared_proto.FarmStealLogProto) error
	CreateGuild(context.Context, int64, []byte) error
	CreateHero(context.Context, *entity.Hero) (bool, error)
	CreateMail(context.Context, uint64, int64, []byte, bool, bool, bool, int32, int64) error
	CreateOrder(context.Context, string, uint64, int64, uint32, uint32, int64, uint64, int64) error
	DelMcWarHeroRecord(context.Context, int32) error
	DelMcWarHeroRecordWithHeroId(context.Context, int32, uint64, int64) error
	DeleteChatWindow(context.Context, int64, []byte) error
	DeleteGuild(context.Context, int64) error
	DeleteMail(context.Context, uint64, int64) error
	DeleteMultiMail(context.Context, int64, []uint64, bool) error
	FindSettingsOpen(context.Context, shared_proto.SettingType, []int64) ([]int64, error)
	GMFarmRipe(context.Context, int64, int64, int64) error
	HeroId(context.Context, string) (int64, error)
	HeroIdExist(context.Context, int64) (bool, error)
	HeroIds(context.Context) ([]int64, error)
	HeroNameExist(context.Context, string) (bool, error)
	InsertBaiZhanReplay(context.Context, int64, int64, *shared_proto.BaiZhanReplayProto, bool, int64) error
	InsertGuildLog(context.Context, int64, *shared_proto.GuildLogProto) error
	InsertXuanyRecord(context.Context, int64, []byte) (int64, error)
	IsCollectableMail(context.Context, uint64) (bool, error)
	ListHeroChatMsg(context.Context, []byte, uint64) ([]*shared_proto.ChatMsgProto, error)
	ListHeroChatWindow(context.Context, int64) ([]uint64, isql.BytesArray, error)
	LoadAllGuild(context.Context) ([]*sharedguilddata.Guild, error)
	LoadAllHeroData(context.Context) ([]*entity.Hero, error)
	LoadAllRegionHero(context.Context) ([]*entity.Hero, error)
	LoadBaiZhanRecord(context.Context, int64, uint64) (isql.BytesArray, error)
	LoadCanStealCount(context.Context, int64, int64, int64, uint64) (uint64, error)
	LoadCanStealCube(context.Context, int64, int64, int64, uint64) ([]*entity.FarmCube, error)
	LoadChatMsg(context.Context, int64) (*shared_proto.ChatMsgProto, error)
	LoadCollectMailPrize(context.Context, uint64, int64) (*resdata.Prize, error)
	LoadFarmCube(context.Context, int64, cb.Cube) (*entity.FarmCube, error)
	LoadFarmCubes(context.Context, int64) ([]*entity.FarmCube, error)
	LoadFarmHarvestCubes(context.Context, int64) ([]*entity.FarmCube, error)
	LoadFarmLog(context.Context, int64, uint64) ([]*shared_proto.FarmStealLogProto, error)
	LoadFarmStealCount(context.Context, int64, int64, cb.Cube) (uint64, error)
	LoadFarmStealCubes(context.Context, int64, int64, uint64) ([]*entity.FarmCube, error)
	LoadGuild(context.Context, int64) (*sharedguilddata.Guild, error)
	LoadGuildLogs(context.Context, int64, shared_proto.GuildLogType, int64, uint64) ([]*shared_proto.GuildLogProto, error)
	LoadHero(context.Context, int64) (*entity.Hero, error)
	LoadHeroCount(context.Context) (uint64, error)
	LoadHeroListByCountry(context.Context, uint64, uint64, uint64) ([]*entity.Hero, error)
	LoadHeroListByNameAndCountry(context.Context, string, uint64, uint64, uint64) ([]*entity.Hero, error)
	LoadHeroMailList(context.Context, int64, uint64, int32, int32, int32, int32, int32, int32, uint64) ([]*shared_proto.MailProto, error)
	LoadHerosByName(context.Context, string, uint64, uint64) ([]*entity.Hero, error)
	LoadJoinedMcWarId(context.Context, int64) (*entity.JoinedMcWarIds, error)
	LoadKey(context.Context, server_proto.Key) ([]byte, error)
	LoadMail(context.Context, uint64) ([]byte, error)
	LoadMailCountHasPrizeNotCollected(context.Context, int64) (int, error)
	LoadMailCountHasReportNotReaded(context.Context, int64, int32) (int, error)
	LoadMailCountNoReportNotReaded(context.Context, int64) (int, error)
	LoadMcWarGuildRecord(context.Context, uint64, uint64, int64) (*shared_proto.McWarTroopsInfoProto, error)
	LoadMcWarHeroRecord(context.Context, uint64, uint64, int64) (*shared_proto.McWarTroopAllRecordProto, error)
	LoadMcWarRecord(context.Context, uint64, uint64) (*shared_proto.McWarFightRecordProto, error)
	LoadNoGuildHeroListByName(context.Context, string, uint64, uint64) ([]*entity.Hero, error)
	LoadRecommendHeros(context.Context, bool, uint64, uint64, uint64, uint64, int64) ([]*entity.Hero, error)
	LoadUnreadChatCount(context.Context, int64) (uint64, error)
	LoadUserMisc(context.Context, int64) (*server_proto.UserMiscProto, error)
	LoadXuanyRecord(context.Context, int64, int64, bool) ([]int64, isql.BytesArray, error)
	MaxGuildId(context.Context) (int64, error)
	MaxMailId(context.Context) (uint64, error)
	OrderExist(context.Context, string) (bool, error)
	PlantFarmCube(context.Context, int64, cb.Cube, int64, int64, uint64) error
	ReadChat(context.Context, int64, []byte) error
	ReadMultiMail(context.Context, int64, []uint64, bool) (*resdata.Prize, error)
	RemoveChatMsg(context.Context, int64) error
	RemoveFarmCube(context.Context, int64, cb.Cube) error
	RemoveFarmLog(context.Context, int32) error
	RemoveFarmSteal(context.Context, int64, []cb.Cube) error
	ResetConflictFarmCubes(context.Context, int64) error
	ResetFarmCubes(context.Context, int64) error
	SaveFarmCube(context.Context, *entity.FarmCube) error
	SaveGuild(context.Context, int64, []byte) error
	SaveHero(context.Context, *entity.Hero) error
	SaveKey(context.Context, server_proto.Key, []byte) error
	SetFarmRipeTime(context.Context, int64, int64) error
	UpdateChatMsg(context.Context, int64, *shared_proto.ChatMsgProto) bool
	UpdateChatWindow(context.Context, int64, []byte, []byte, bool, int32, bool) error
	UpdateFarmCubeRipeTime(context.Context, int64, cb.Cube, int64, int64) error
	UpdateFarmCubeState(context.Context, int64, cb.Cube, int64, int64) error
	UpdateFarmStealTimes(context.Context, int64, []cb.Cube) error
	UpdateHeroGuildId(context.Context, int64, int64) error
	UpdateHeroName(context.Context, int64, string, string) bool
	UpdateHeroOfflineBoolIfExpected(context.Context, int64, isql.OfflineBool, bool, bool) (bool, error)
	UpdateMailCollected(context.Context, uint64, int64, bool) error
	UpdateMailKeep(context.Context, uint64, int64, bool) error
	UpdateMailRead(context.Context, uint64, int64, bool) error
	UpdateSettings(context.Context, int64, uint64) error
	UpdateUserMisc(context.Context, int64, *server_proto.UserMiscProto) error
}

type DepotModule interface {
}

type DianquanModule interface {
}

type DomesticModule interface {
	UseBuffGoods(*buffer.BufferData, uint64, HeroController)
	UseMianGoods(*buffer.BufferData, uint64, HeroController)
}

type DungeonModule interface {
}

type EquipmentModule interface {
	OnAddGoodsEvent(*entity.Hero, herolock.LockResult, uint64, uint64)
}

type ExtraTimesService interface {
	MultiLevelNpcMaxTimes() extratimesface.ExtraMaxTimes
}

type FarmModule interface {
}

type FarmService interface {
	Close()
	FuncNoWait(entity.FarmFuncType)
	FuncWait(string, pbutil.Buffer, HeroController, entity.FarmFuncType)
	GMCanSteal(int64)
	GMRipe(int64)
	ReduceRipeTime(int64, time.Duration)
	ReduceRipeTimePercent(int64, *entity.BuffInfo, *entity.BuffInfo)
	// 更新农场地块
	UpdateFarmCubeWithOffset(int64, uint64, int, int, []cb.Cube, bool, time.Time)
	/*
	 更新农场地块
	 absCubes 所有本次变化的地块
	 allConflictedBlocks 本次变化的地块中，有冲突的地块
	 npcConflictOffsets 所有 npc 造成的冲突地块
	*/
	UpdateFarmCubes(int64, uint64, []cb.Cube, []cb.Cube, []cb.Cube, int, int, time.Time)
}

type FightService interface {
	Close()
	SendFightRequest(*entity.TlogFightContext, *scene.CombatScene, int64, int64, *shared_proto.CombatPlayerProto, *shared_proto.CombatPlayerProto) *server_proto.CombatResponseServerProto
	SendFightRequestReturnResult(*entity.TlogFightContext, *scene.CombatScene, int64, int64, *shared_proto.CombatPlayerProto, *shared_proto.CombatPlayerProto, bool) *server_proto.CombatResponseServerProto
	SendMultiFightRequest(*entity.TlogFightContext, *scene.CombatScene, []int64, []int64, []*shared_proto.CombatPlayerProto, []*shared_proto.CombatPlayerProto, int32, int32, int32) *server_proto.MultiCombatResponseServerProto
	SendMultiFightRequestReturnResult(*entity.TlogFightContext, *scene.CombatScene, []int64, []int64, []*shared_proto.CombatPlayerProto, []*shared_proto.CombatPlayerProto, int32, int32, int32, bool) *server_proto.MultiCombatResponseServerProto
}

type FightXService interface {
	Close()
	// 新版战斗
	SendFightRequest(*entity.TlogFightContext, *scene.CombatScene, int64, int64, *shared_proto.CombatPlayerProto, *shared_proto.CombatPlayerProto) *server_proto.CombatXResponseServerProto
	SendFightRequestReturnResult(*entity.TlogFightContext, *scene.CombatScene, int64, int64, *shared_proto.CombatPlayerProto, *shared_proto.CombatPlayerProto, bool) *server_proto.CombatXResponseServerProto
}

type FishingModule interface {
	GmFishingRate(int64, uint64, uint64)
}

type GameExporter interface {
	GetMsgTimeCost() *metrics.MsgTimeCostSummary
	GetRegisterCounter() *atomic.Uint64
	Start() (http.Handler, error)
}

type GameServer interface {
	GetRpcPort() uint32
	GetTcpPort() uint32
	Serve(ServeListener, ConnHandler)
}

type GardenModule interface {
}

// 宝石模块

type GemModule interface {
}

type GmModule interface {
}

type GuildModule interface {
	GmAddBigBoxEnergy(int64, uint64)
	GmAddGuildBuildAmount(int64, uint64)
	GmAddGuildYinliang(int64, int64)
	GmGiveGuildEventPrize(*entity.Hero, herolock.LockResult, []*guild_data.GuildEventPrizeData)
	GmMiaoGuildTechCd(int64)
	GmOpenImpeachLeader(int64)
	GmRemoveNpcGuild()
	GmUpgradeGuildLevel(int64)
	HandleGiveGuildEventPrize(*entity.Hero, herolock.LockResult, int64, []*guild_data.GuildEventPrizeData, uint64)
	OnHeroOnline(HeroController, int64)
}

type GuildService interface {

	// 增加联盟任务进度
	AddGuildTaskProgress(int64, *guild_data.GuildTaskData, uint64)
	AddHufu(uint64, int64, int64, string, string) bool
	AddLog(int64, *shared_proto.GuildLogProto)
	AddLogWithMemberIds(int64, []int64, *shared_proto.GuildLogProto)
	AddRecommendInviteHeros(int64)
	Broadcast(int64, pbutil.Buffer)
	CheckAndAddRecommendInviteHeros(int64)
	ClearSelfGuildMsgCache(int64)
	Close()
	Func(sharedguilddata.Funcs) bool
	FuncGuild(int64, sharedguilddata.Func)
	GetGuildFlagName(int64) string
	GetGuildIdByFlagName(string) int64
	GetGuildIdByName(string) int64
	// 根据国家id获取声望排名列表的消息
	GetGuildPrestigeRankMsg(uint64, time.Time) pbutil.Buffer
	GetSnapshot(int64) *guildsnapshotdata.GuildSnapshot
	// 筛选推荐联盟
	RecommendGuildList(*snapshotdata.HeroSnapshot) []int64
	RecommendInviteHeroList(uint64, uint64, []int64) []*snapshotdata.HeroSnapshot
	RegisterCallback(guildsnapshotdata.Callback)
	RemoveRecommendInviteHero(int64)
	RemoveSnapshot(int64)
	SaveChangedGuild()
	SelfGuildMsgCache() concurrent.I64BufferMap
	SetGuildRankFunc(sharedguilddata.GetGuildRankFunc, sharedguilddata.GenerateRankMsgFunc)
	TimeoutFunc(sharedguilddata.Funcs) bool
	UpdateGuildHeroSnapshot(*sharedguilddata.Guild)
	UpdateSnapshot(*sharedguilddata.Guild) *guildsnapshotdata.GuildSnapshot
}

type GuildSnapshotService interface {
	GetGuildBasicProto(int64) *shared_proto.GuildBasicProto
	GetGuildLevel(int64) uint64
	GetSnapshot(int64) *guildsnapshotdata.GuildSnapshot
	// 注册监听snapshot变化callback
	RegisterCallback(guildsnapshotdata.Callback)
	RemoveSnapshot(int64)
	UpdateSnapshot(*guildsnapshotdata.GuildSnapshot)
}

type HebiModule interface {
	Close()
	UpdateGuildInfo(int64, int64)
	UpdateGuildInfoBatch([]int64, int64)
}

type HeroController interface {
	AddTickFunc(face.BFunc)
	Disconnect(msg.ErrMsg)
	Func(herolock.Func)
	FuncNotError(herolock.FuncNotError) bool
	FuncWithSend(herolock.SendFunc) bool
	//GetBlockIndex 获取地块索引
	GetBlockIndex() interface{}
	GetCareCondition() *server_proto.MilitaryConditionProto
	GetCareWaterTimesMap() map[int64]uint64
	GetClientIp() string
	GetClientIp32() uint32
	GetIsInBackgroud() bool
	GetPf() uint32
	GetViewArea() *realmface.ViewArea
	//GetWatchObjList 获取观察列表
	GetWatchObjList() map[interface{}]int
	Id() int64
	IdBytes() []byte
	IsClosed() bool
	LastClickTime() time.Time
	LockGetGuildId() (int64, bool)
	LockHeroCountry() uint64
	NextRefreshRecommendHeroTime() time.Time
	NextSearchHeroTime() time.Time
	NextSearchNoGuildHeros() time.Time
	Pid() uint32
	//RemoveBlockIndex 退出时 移除出场景
	RemoveBlockIndex()
	// 发送消息.
	Send(pbutil.Buffer)
	// 发送消息.
	SendAll([]pbutil.Buffer)
	// 发送在线路繁忙时可以被丢掉的消息
	SendIfFree(pbutil.Buffer)
	//SetBlockIndex 设置地块索引(AOI)
	SetBlockIndex(interface{}) interface{}
	SetCareCondition(*server_proto.MilitaryConditionProto)
	SetCareWaterTimesMap(map[int64]uint64)
	SetIsInBackgroud(time.Time, bool)
	SetLastClickTime(time.Time)
	SetNextSearchNoGuildHeros(time.Time)
	SetViewArea(*realmface.ViewArea)
	//AddWatchObjList 设置观察对象列表 如果设置nil,表示清空
	SetWatchObjList(map[interface{}]int)
	Sid() uint32
	TickFunc()
	TotalOnlineTime() time.Duration
	TryNextWriteOnlineLogTime(time.Time, time.Duration) bool
	UpdateIsInBackgroud(time.Time)
	UpdateNextRefreshRecommendHeroTime(time.Time)
	UpdateNextSearchHeroTime(time.Time)
}

type HeroDataService interface {
	Close()
	Create(*entity.Hero) error
	Exist(int64) (bool, error)
	Func(int64, herolock.Func)
	FuncNotError(int64, herolock.FuncNotError) bool
	FuncWithSend(int64, herolock.SendFunc) bool
	FuncWithSendError(int64, herolock.SendFunc) (bool, error)
	NewHeroLocker(int64) herolock.HeroLocker
	Put(*entity.Hero) error
}

type HeroSnapshotService interface {

	// 获得英雄的snapshot, 并不保证是最新的. 尽量
	// 就算英雄要从db中加载, 也不会触发callback.
	// 返回nil也可能是数据库报错, 英雄未必不存在
	Get(int64) *snapshotdata.HeroSnapshot
	GetBasicProto(int64) *shared_proto.HeroBasicProto
	GetBasicSnapshotProto(int64) *shared_proto.HeroBasicSnapshotProto
	GetFlagHeroName(int64) string
	// 只从Cache中获取，不读取DB
	GetFromCache(int64) *snapshotdata.HeroSnapshot
	GetHeroName(int64) string
	GetTlogHero(int64) entity.TlogHero
	// 创建个新的snapshot, 但是并没有保存. 等unlock后再调用Cache保存
	// 必须是lock住Hero的情况下才能调用, 确保此时没有其他人能访问hero对象
	NewSnapshot(*entity.Hero) *snapshotdata.HeroSnapshot
	// 英雄下线时调用, 把snapshot移动到lru中
	Offline(int64)
	// 英雄上线时调用, 保存snapshot. 这个snapshot必须是没有变化的数据的, 只是上个线而已. 不会触发callback
	// 如果英雄上线导致snapshot中缓存的数据有了变化, 必须再调用一次Update, 把这个变化告知其他系统
	Online(*snapshotdata.HeroSnapshot)
	// 注册监听snapshot变化callback
	RegisterCallback(snapshotdata.SnapshotCallback)
	// 改变了英雄数据后, 缓存英雄snapshot, 必须是在unlock后调用
	// 如果snapshot的版本号低于缓存中的版本号, 不会触发callback
	Update(*snapshotdata.HeroSnapshot)
}

type IndividualServerConfig interface {
	GetDisablePush() bool
	GetDontEncrypt() bool
	GetGameAppID() string
	GetHttpPort() int
	GetIgnoreHeartBeat() bool
	GetIsAllowRobot() bool
	GetIsDebug() bool
	GetKafkaBrokerAddr() []string
	GetKafkaStart() bool
	GetLocalAddStr() string
	GetPlatformID() int
	GetPort() int
	GetReplayPrefix() string
	GetServerID() int
	GetServerInfo() *shared_proto.HeroServerInfoProto
	GetServerStartTime() time.Time
	GetSkipHeader() bool
	GetTlogStart() bool
	GetTlogTopic() string
	GetZoneAreaID() int
}

type KafkaService interface {
	Close()
	NewProducerMsg(entity.Topic, []byte) *sarama.ProducerMessage
	SendAsync(*sarama.ProducerMessage)
	SendSync(*sarama.ProducerMessage) error
}

type LightpawHandler interface {
}

type LocationHeroCache interface {
	UpdateHero(*snapshotdata.HeroSnapshot, int64)
	UpdateLocation(int64, uint64)
}

type MailModule interface {
	OnHeroOnline(HeroController)
	// 过期函数@AlbertFan
	// 发邮件
	SendMail(int64, uint64, string, string, bool, *shared_proto.FightReportProto, *shared_proto.PrizeProto, time.Time) bool
	SendProtoMail(int64, *shared_proto.MailProto, time.Time) bool
	SendReportMail(int64, *shared_proto.MailProto, time.Time) bool
}

type MetricsRegister interface {
	EnableDBMetrics() bool
	EnableMsgMetrics() bool
	EnableOnlineCountMetrics() bool
	EnablePanicMetrics() bool
	EnableRegisterCountMetrics() bool
	Register(prometheus.Collector)
	RegisterFunc(metrics.CollectFunc)
}

type MilitaryModule interface {
	GmRate(int64, uint64, uint64)
}

type MingcModule interface {
}

type MingcService interface {
	AllInOneGuild() int64
	Build(int64, int64, uint64, uint64) bool
	CaptainHostGuild(uint64) int64
	Close()
	Country(uint64) uint64
	// 国家当前占领的本国初始名城
	CountryHoldInitMcs(uint64) []*entity.Mingc
	DisableMcBuildLogCache(uint64)
	GetMcBuildGuildMemberPrize(uint64) *resdata.Prize
	GuildMingc(int64) *entity.Mingc
	IsHoldCountryCapital(int64, uint64) bool
	McBuildLogMsg(*entity.Mingc) pbutil.Buffer
	Mingc(uint64) *entity.Mingc
	MingcsMsg(uint64) pbutil.Buffer
	SetHostGuild(uint64, int64) bool
	UpdateMsg() pbutil.Buffer
	WalkMingcs(entity.MingcFunc)
}

type MingcWarModule interface {
}

type MingcWarService interface {
	ApplyAst(int64, bool, *mingcdata.MingcBaseData) (pbutil.Buffer, pbutil.Buffer)
	ApplyAstNotice(int64, int64) bool
	ApplyAtk(int64, *mingcdata.MingcBaseData, uint64) (pbutil.Buffer, pbutil.Buffer)
	ApplyAtkNotice(int64, int64) bool
	BuildFightStartMsg(time.Time) (pbutil.Buffer, bool)
	CancelApplyAst(int64, *mingcdata.MingcBaseData) (pbutil.Buffer, pbutil.Buffer)
	CatchHistoryRecord(int64, int64) pbutil.Buffer
	CatchTroopsRank(int64, uint64) (pbutil.Buffer, pbutil.Buffer)
	CleanOnGuildRemoved(int64)
	Close()
	CurrMcWarStage() (int32, time.Time, time.Time)
	GmApplyAtkGuild(uint64, int64)
	GmCampFail(uint64, bool)
	GmChangeStage(shared_proto.MingcWarState, HeroController, time.Time)
	GmNewMingcWar()
	GmSetAstAtkGuild(uint64, int64)
	GmSetAstDefGuild(uint64, int64)
	GmSetDefGuild(uint64, int64)
	GuildMcWarType(int64) (uint64, shared_proto.MingcWarGuildType)
	JoinFight(HeroController, uint64, []uint64, []int32) (pbutil.Buffer, pbutil.Buffer)
	JoiningFightMingc(int64) (uint64, bool)
	McWarStartEndTime() (time.Time, time.Time)
	QuitFight(int64) (pbutil.Buffer, pbutil.Buffer)
	QuitWatch(HeroController, uint64) (pbutil.Buffer, pbutil.Buffer)
	ReplyApplyAst(int64, int64, *mingcdata.MingcBaseData, bool) (pbutil.Buffer, pbutil.Buffer)
	SceneBack(int64) (pbutil.Buffer, pbutil.Buffer)
	SceneChangeMode(int64, shared_proto.MingcWarModeType) (pbutil.Buffer, pbutil.Buffer)
	SceneDrum(int64) (pbutil.Buffer, pbutil.Buffer)
	SceneMove(int64, cb.Cube) (pbutil.Buffer, pbutil.Buffer)
	SceneSpeedUp(int64, float64) (pbutil.Buffer, pbutil.Buffer)
	SceneTouShiBuildingFire(int64, cb.Cube) (pbutil.Buffer, pbutil.Buffer)
	SceneTouShiBuildingTurnTo(int64, cb.Cube, bool) (pbutil.Buffer, pbutil.Buffer)
	SceneTroopRelive(int64) (pbutil.Buffer, pbutil.Buffer)
	SendChat(int64, *shared_proto.ChatMsgProto) pbutil.Buffer
	UpdateMsg()
	ViewMcWarMcMsg(uint64) (pbutil.Buffer, pbutil.Buffer)
	ViewMcWarSceneMsg(uint64) (pbutil.Buffer, pbutil.Buffer)
	ViewMsg(uint64) pbutil.Buffer
	ViewSceneTroopRecord(int64) (pbutil.Buffer, pbutil.Buffer)
	ViewSelfGuildProto(int64) *shared_proto.McWarGuildProto
	Watch(HeroController, uint64) (pbutil.Buffer, pbutil.Buffer)
}

type MiscModule interface {
	GmSetClientVersion(string)
	OnHeroOnline(HeroController)
	PrintClientLog(int64, string, string, string)
	SendClientVersion(sender.Sender)
	SendConfig(*misc.C2SConfigProto, sender.Sender)
	SendLuaConfig(*misc.C2SConfigluaProto, sender.Sender)
	SyncTime(int32, sender.Sender)
}

type Modules interface {
	ActivityModule() ActivityModule
	BaiZhanModule() BaiZhanModule
	ChatModule() ChatModule
	ClientConfigModule() ClientConfigModule
	CountryModule() CountryModule
	DepotModule() DepotModule
	DianquanModule() DianquanModule
	DomesticModule() DomesticModule
	DungeonModule() DungeonModule
	EquipmentModule() EquipmentModule
	FarmModule() FarmModule
	FishingModule() FishingModule
	GardenModule() GardenModule
	GemModule() GemModule
	GuildModule() GuildModule
	HebiModule() HebiModule
	MailModule() MailModule
	MilitaryModule() MilitaryModule
	MingcModule() MingcModule
	MingcWarModule() MingcWarModule
	MiscModule() MiscModule
	PromotionModule() PromotionModule
	QuestionModule() QuestionModule
	RandomEventModule() RandomEventModule
	RankModule() RankModule
	RedPacketModule() RedPacketModule
	RegionModule() RegionModule
	RelationModule() RelationModule
	SecretTowerModule() SecretTowerModule
	ShopModule() ShopModule
	StrategyModule() StrategyModule
	StressModule() StressModule
	SurveyModule() SurveyModule
	TagModule() TagModule
	TaskModule() TaskModule
	TeachModule() TeachModule
	TowerModule() TowerModule
	VipModule() VipModule
	XiongNuModule() XiongNuModule
	XuanyuanModule() XuanyuanModule
	ZhanJiangModule() ZhanJiangModule
	ZhengWuModule() ZhengWuModule
}

type ProductService interface {
}

type PromotionModule interface {
}

type PushService interface {
	GmPush(int64)
	MultiPush(shared_proto.SettingType, []int64, int64)
	MultiPushFunc(shared_proto.SettingType, []int64, int64, pushdata.PushFunc)
	MultiPushTitleContent(shared_proto.SettingType, string, string, []int64, int64)
	Push(shared_proto.SettingType, int64)
	PushFunc(shared_proto.SettingType, int64, pushdata.PushFunc)
	PushTitleContent(shared_proto.SettingType, string, string, int64)
}

type QuestionModule interface {
}

type RandomEventModule interface {
}

type RankModule interface {
	AddOrUpdateRankObj(rankface.RankObj)
	Close()
	CountryOfficial(int, string, uint64, shared_proto.CountryOfficialType) []*shared_proto.HeroBasicSnapshotProto
	RemoveRankObj(shared_proto.RankType, int64)
	// 百战千军类型的特殊排行榜不能从这里获取
	SingleRRankListFunc(shared_proto.RankType, rankface.RRankListFunc) bool
	SubTypeRRankListFunc(shared_proto.RankType, uint64, rankface.RRankListFunc) bool
	UpdateBaiZhanRankList([]rankface.RankObj)
	UpdateXuanyRankList([]rankface.RankObj)
}

type Realm interface {
	AddAstDefendLog(int64, time.Time, time.Time, string, uint64) bool
	// 把玩家的基地或行营加入进来, 可以是新英雄, 可以是随机迁城, 可以是高级快速迁城. 此时英雄的城必须不是流亡状态.
	// 英雄的主城或行营当前必须不能已经属于其他地区管理.
	// 加入前必须已预定坐标. processed为true的话, 就算err也不需要取消预定
	// isHome表示是不是主城.
	// return processed 是否已处理. err 表示是否处理有错. 根据err判断错误的类型
	AddBase(int64, int, int, realmface.AddBaseType) (bool, error)
	AddGuildWorkshop(int64, int, int, int32, int32) bool
	AddHeroBaoZangMonster(*regdata.BaozNpcData, int, int, int64, int32) (bool, bool)
	AddHomeNpc(HeroController, []*basedata.HomeNpcBaseData) (bool, error)
	AddInvasionMonster(int64, shared_proto.MultiLevelNpcType, uint64)
	// 英雄操作导致繁荣度增加. 传入增加的量, 由这里执行具体增加的操作
	// 升级建筑在英雄线程不要直接增加繁荣度, 要调这个方法来修改繁荣度.
	AddProsperity(int64, uint64) (bool, error)
	AddXiongNuBase(xiongnuface.RResistXiongNuInfo, int, int, int, int) (int64, int32, int32)
	AddXiongNuTroop(int64, []int64, []*monsterdata.MonsterMasterData) bool
	AroundBase(int, int, int, int) bool
	// 宝藏遣返
	BaozRepatriate(HeroController, int64, int64) (bool, error)
	CalcMoveSpeed(int64, float64) float64
	// 班师回朝
	CancelInvasion(HeroController, int64) (bool, error)
	// 取消预定的坐标. 由于某种原因, 哥来不了了.
	CancelReservedPos(int, int)
	// 取消缓慢迁城
	// 取消缓慢迁城，快速迁城，流亡，等等会自动取消缓慢迁移
	CancelSlowMoveBase(HeroController) (bool, error)
	// 变更个人签名
	ChangeSign(int64, string)
	// 变更个人签名
	ChangeTitle(int64, uint64)
	CheckCanMoveBase(int64, int, int, bool) error
	CheckIsFucked(int64) bool
	ClearAstDefendLog(int64)
	// 创建集结
	CreateAssembly(HeroController, int64, uint64, uint64, uint64, time.Duration) (bool, error)
	// 自己驱逐自己城里的坏人
	Expel(HeroController, int64, uint64) (bool, bool, string, error)
	GetAstDefendHeros(int64) []*shared_proto.HeroBasicProto
	GetAstDefendLogs() *server_proto.AllAstDefendLogProto
	GetAstDefendLogsByHero(int64) []*shared_proto.AstDefendLogProto
	GetAstDefendingTroopCount(int64) uint64
	GetBaseLevel(uint64) *domestic_data.BaseLevelData
	// 获得跟我土地有冲突的玩家id
	GetConflictHeroIds(HeroController) (bool, []int64)
	GetDefendingTroopCount(int64) uint64
	GetHeroBaozRoBase(int64) *server_proto.RoBaseProto
	GetMapData() *blockdata.StitchedBlocks
	GetMaxXiongNuTroopFightingAmount(int64, int64) (bool, uint64)
	GetRadius() uint64
	GetRoBase(int64) *server_proto.RoBaseProto
	GetRoBaseByPos(int, int) *server_proto.RoBaseProto
	GetRuinsBase(int, int) int64
	GetXiongNuInvateTargetCount(int64) i64.GetU64
	GetXiongNuTroopInfo(int64, int64) *shared_proto.XiongNuBaseTroopProto
	GmReduceProsperity(int64, uint64) bool
	GmRefreshBaoZangNpc()
	GmSpeedUpFightMe(int64)
	// 破坏联盟工坊
	HurtGuildWorkshop(HeroController, int64) (bool, bool)
	Id() int64
	// 出发攻打/帮忙驱逐
	// 没有err的话调用者还需要发送成功消息
	Invasion(HeroController, shared_proto.TroopOperate, int64, uint64, uint64, uint64) (bool, error)
	// 出发侦察
	InvasionInvestigate(HeroController, int64) (bool, error)
	IsEdgeNotHomePos(int, int) bool
	IsPosOpened(int, int) bool
	// 加入集结
	JoinAssembly(HeroController, int64, int64, uint64) (bool, error)
	Mian(int64, time.Time, bool) (bool, error)
	// 在同一个地图中移动基地
	// 移动前必须已预定坐标. processed为true的话, 就算err也不需要取消预定
	MoveBase(HeroController, int, int, int, int, bool) (bool, error)
	OnHeroLogin(HeroController)
	QueryTroopUnit(HeroController, int64, int64) (bool, error)
	RandomAroundBase(int, int) (int, int, bool)
	RandomBasePos() (int, int)
	ReduceProsperity(int64, uint64) bool
	// 把玩家在这个地图上的基地或者行营移除. 玩家的部队必须都已不在外面. 而且不能是流亡状态且归这个地图管. (流亡状态的话, 都不归这里管)
	RemoveBase(int64, bool, *i18n.I18nRef, *i18n.I18nRef) (bool, error, int, int)
	// 遣返
	Repatriate(HeroController, int64) (bool, error)
	// 预定一个随机主城坐标, 返回的是个可以建主城的位置
	ReserveNewHeroHomePos(uint64) (bool, int, int)
	// 预定一个坐标
	ReservePos(int, int) bool
	// 在同一个场景迁城，预定一个坐标
	ReservePosForMoveBase(int, int, int, int) bool
	// 预定一个随机主城坐标, 返回的是个可以建主城的位置
	ReserveRandomHomePos(realmface.RandomPointType) (bool, int, int)
	// 查看集结
	ShowAssembly(HeroController, int64, int64, int32)
	// 加速
	SpeedUp(HeroController, int64, int64, float64, uint64) (bool, error)
	StartCareMilitary(HeroController) bool
	// 开始关心这个地图, 获得地图中所有主城的信息
	StartCareRealm(HeroController, int, int, int, int) bool
	StopCareRealm(HeroController) bool
	TryRemoveBaseMian(int64) bool
	// 改变英雄基础信息（含帮派）
	UpdateHeroBasicInfoNoBlock(int64)
	UpdateHeroRealmInfo(int64, bool, bool, bool)
	UpdateProsperity(int64) bool
	UpdateProsperityBuff(int64, uint64, uint64, uint64)
	// 手动升级老家等级
	UpgradeBase(HeroController) (bool, error)
}

type RealmService interface {
	AddProsperityFunc(int64, int64, uint64, string) Func
	CheckCanMoveBase(int64, int, int, bool) bool
	Close()
	DoMoveBase(shared_proto.GoodsMoveBaseType, Realm, HeroController, int, int, int, int, bool) bool
	GetBigMap() Realm
	GetRealm(int64) Realm
	OnGuildSnapshotRemoved(int64)
	// callback
	OnGuildSnapshotUpdated(*guildsnapshotdata.GuildSnapshot, *guildsnapshotdata.GuildSnapshot)
	// 坐标是已经占座占好了的, 直接调用realm.AddBase就可以了
	ReserveNewHeroHomePos(uint64) (Realm, int, int)
	// 坐标是已经占座占好了的, 直接调用realm.AddBase就可以了
	ReserveRandomHomePos(realmface.RandomPointType) (Realm, int, int)
	StartCareMilitary(HeroController)
}

type RedPacketModule interface {
}

type RedPacketService interface {
	AllGrabbed(int64) bool
	Close()
	Create(int64, *red_packet.RedPacketData, uint64, string, shared_proto.ChatType) (int64, string, msg.ErrMsg)
	Exist(int64) bool
	Expired(int64, time.Time) bool
	Grab(int64, int64, int64) (uint64, bool, *shared_proto.RedPacketProto, msg.ErrMsg)
	Grabbed(int64, int64) bool
	RedPacketChatId(int64) int64
	SetRedPacketChatId(int64, int64)
}

type RegionModule interface {
	InitHeroBase(HeroController, time.Time, uint64, realmface.AddBaseType) bool
	UseMianGoods(uint64, bool, HeroController) bool
}

type RelationModule interface {
}

type ReminderService interface {
	ChangeAttackOrRobCount(int64, int64, int64, int64, bool)
	OnHeroOnline(HeroController)
}

type SeasonService interface {
	GetSeasonTickTime() tickdata.TickTime
	OnHeroOnline(HeroController)
	Season() *season.SeasonData
	SeasonByTime(time.Time) *season.SeasonData
}

type SecretTowerModule interface {
	Close()
	OnHeroOffline(HeroController)
	OnHeroOnline(HeroController)
}

type ServerStartStopTimeService interface {
	IsNormalStop() bool
	SaveStartTime()
	SaveStopTime()
}

type ServiceDep interface {
	Broadcast() BroadcastService
	Chat() ChatService
	Country() CountryService
	Datas() ConfigDatas
	Db() DbService
	Fight() FightService
	FightX() FightXService
	Guild() GuildService
	GuildSnapshot() GuildSnapshotService
	HeroData() HeroDataService
	HeroSnapshot() HeroSnapshotService
	Mail() MailModule
	Mingc() MingcService
	Push() PushService
	SvrConf() IndividualServerConfig
	Time() TimeService
	Tlog() TlogService
	World() WorldService
}

type ShopModule interface {
}

type StrategyModule interface {
	GMStrategy(uint64, HeroController)
}

type StressModule interface {
}

type SurveyModule interface {
	GmGiveSurveyPrize(int64, string)
}

type TagModule interface {
}

type TaskModule interface {
	OnHeroOnline(HeroController)
}

type TeachModule interface {
}

type TickerService interface {
	Close()
	GetDailyMcTickTime() tickdata.TickTime
	GetDailyTickTime() tickdata.TickTime
	GetDailyZeroTickTime() tickdata.TickTime
	GetPer10MinuteTickTime() tickdata.TickTime
	GetPer30MinuteTickTime() tickdata.TickTime
	GetPerHourTickTime() tickdata.TickTime
	GetPerMinuteTickTime() tickdata.TickTime
	GetWeeklyTickTime() tickdata.TickTime
	TickPer10Minute(string, TickFunc) Func
	TickPer30Minute(string, TickFunc) Func
	TickPerDay(string, TickFunc) Func
	TickPerDayZero(string, TickFunc) Func
	TickPerHour(string, TickFunc) Func
	TickPerMinute(string, TickFunc) Func
	TickTickPerWeek(string, TickFunc) Func
}

type TimeLimitGiftService interface {
	Close()
	EncodeClient() []*shared_proto.TimeLimitGiftProto
	GetGiftEndTime(uint64) (time.Time, bool)
	OnHeroOnline(HeroController)
}

type TimeService interface {
	CurrentTime() time.Time
}

type TlogBaseService interface {
	Close()
	WriteTlog(string) bool
}

type TlogService interface {
	BuildAccountRegister(int64, *shared_proto.TencentInfoProto) string
	BuildAdvanceSoulFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64) string
	BuildAnswerFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, bool) string
	BuildBaiZhanFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64) string
	BuildCareFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, string, uint64, uint64) string
	BuildChangeCaptainFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) string
	BuildChatFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string
	BuildCityExpFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string
	BuildEquipmentAddStarFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64) string
	BuildFarmFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, string, uint64, uint64, uint64) string
	BuildFishFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64) string
	BuildGUOGUANFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64) string
	BuildGameSvrState() string
	BuildGameplayFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string
	BuildGuideFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64) string
	BuildGuildFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64) string
	BuildItemFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, int64) string
	BuildKingExpFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string
	BuildMailFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, string) string
	BuildMoneyFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64) string
	BuildMountRefreshFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64) string
	BuildMoveCitylFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, int64, int64, int64, int64) string
	BuildNationalFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string
	BuildPlayerCultivateFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) string
	BuildPlayerEquipFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64) string
	BuildPlayerExpDrugFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string
	BuildPlayerHaunterFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64) string
	BuildPlayerLogin(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64) string
	BuildPlayerLogout(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64) string
	BuildPlayerRegister(entity.TlogHero, *shared_proto.TencentInfoProto) string
	BuildRefreshFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string
	BuildResearchFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64) string
	BuildResourceStockFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64) string
	BuildRoundFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) string
	BuildSnsFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64) string
	BuildSpeedUpFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64) string
	BuildStrenghBuildingFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string
	BuildStrenghEquipmentFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) string
	BuildTaskFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64, uint64) string
	BuildVipLevelFlow(entity.TlogHero, *shared_proto.TencentInfoProto, uint64, uint64) string
	Close()
	DontGenTlog() bool
	TlogAccountRegister(int64)
	TlogAdvanceSoulFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogAdvanceSoulFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogAnswerFlow(entity.TlogHero, uint64, bool)
	TlogAnswerFlowById(int64, uint64, bool)
	TlogBaiZhanFlow(entity.TlogHero, uint64, uint64)
	TlogBaiZhanFlowById(int64, uint64, uint64)
	TlogCareFlow(entity.TlogHero, uint64, uint64, string, uint64, uint64)
	TlogCareFlowById(int64, uint64, uint64, string, uint64, uint64)
	TlogChangeCaptainFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogChangeCaptainFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogChatFlow(entity.TlogHero, uint64, uint64, uint64)
	TlogChatFlowById(int64, uint64, uint64, uint64)
	TlogCityExpFlow(entity.TlogHero, uint64, uint64, uint64)
	TlogCityExpFlowById(int64, uint64, uint64, uint64)
	TlogEquipmentAddStarFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64)
	TlogEquipmentAddStarFlowById(int64, uint64, uint64, uint64, uint64, uint64)
	TlogFarmFlow(entity.TlogHero, uint64, string, uint64, uint64, uint64)
	TlogFarmFlowById(int64, uint64, string, uint64, uint64, uint64)
	TlogFishFlow(entity.TlogHero, uint64, uint64)
	TlogFishFlowById(int64, uint64, uint64)
	TlogGUOGUANFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogGUOGUANFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogGameSvrState()
	TlogGameplayFlow(entity.TlogHero, uint64, uint64, uint64)
	TlogGameplayFlowById(int64, uint64, uint64, uint64)
	TlogGuideFlow(entity.TlogHero, uint64, uint64)
	TlogGuideFlowById(int64, uint64, uint64)
	TlogGuildFlow(entity.TlogHero, uint64, uint64, uint64, uint64)
	TlogGuildFlowById(int64, uint64, uint64, uint64, uint64)
	TlogItemFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, int64)
	TlogItemFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, int64)
	TlogKingExpFlow(entity.TlogHero, uint64, uint64, uint64)
	TlogKingExpFlowById(int64, uint64, uint64, uint64)
	TlogMailFlow(entity.TlogHero, uint64, string)
	TlogMailFlowById(int64, uint64, string)
	TlogMoneyFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64)
	TlogMoneyFlowById(int64, uint64, uint64, uint64, uint64, uint64)
	TlogMountRefreshFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64)
	TlogMountRefreshFlowById(int64, uint64, uint64, uint64, uint64, uint64)
	TlogMoveCitylFlow(entity.TlogHero, uint64, int64, int64, int64, int64)
	TlogMoveCitylFlowById(int64, uint64, int64, int64, int64, int64)
	TlogNationalFlow(entity.TlogHero, uint64, uint64, uint64)
	TlogNationalFlowById(int64, uint64, uint64, uint64)
	TlogPlayerCultivateFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogPlayerCultivateFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogPlayerEquipFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogPlayerEquipFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogPlayerExpDrugFlow(entity.TlogHero, uint64, uint64, uint64)
	TlogPlayerExpDrugFlowById(int64, uint64, uint64, uint64)
	TlogPlayerHaunterFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogPlayerHaunterFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogPlayerLogin(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64)
	TlogPlayerLoginById(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64)
	TlogPlayerLogout(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64)
	TlogPlayerLogoutById(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, []uint64, uint64, []uint64, uint64, []uint64, uint64)
	TlogPlayerRegister(entity.TlogHero)
	TlogPlayerRegisterById(int64)
	TlogRefreshFlow(entity.TlogHero, uint64, uint64, uint64)
	TlogRefreshFlowById(int64, uint64, uint64, uint64)
	TlogResearchFlow(entity.TlogHero, uint64, uint64, uint64, uint64)
	TlogResearchFlowById(int64, uint64, uint64, uint64, uint64)
	TlogResourceStockFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogResourceStockFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogRoundFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogRoundFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogSnsFlow(entity.TlogHero, uint64, uint64)
	TlogSnsFlowById(int64, uint64, uint64)
	TlogSpeedUpFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogSpeedUpFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogStrenghBuildingFlow(entity.TlogHero, uint64, uint64, uint64)
	TlogStrenghBuildingFlowById(int64, uint64, uint64, uint64)
	TlogStrenghEquipmentFlow(entity.TlogHero, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogStrenghEquipmentFlowById(int64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64)
	TlogTaskFlow(entity.TlogHero, uint64, uint64, uint64)
	TlogTaskFlowById(int64, uint64, uint64, uint64)
	TlogVipLevelFlow(entity.TlogHero, uint64, uint64)
	TlogVipLevelFlowById(int64, uint64, uint64)
	WriteLog(string)
}

type TowerModule interface {
	Close()
}

type TssClient interface {
	CallbackAddr() string
	CheckName(string, bool) (*game2tss.S2CUicJudgeUserInputNameV2Proto, error)
	Client() *rpc7.Client
	Close()
	IsEnable() bool
	JudgeChat(*game2tss.C2SUicJudgeUserInputChatV2Proto) (*game2tss.S2CUicJudgeUserInputChatV2Proto, error)
	RegisterCallback(tss.MsgCategory, tss.Callback)
	TryCheckName(string, sender.Sender, string, pbutil.Buffer, pbutil.Buffer) bool
}

type VipModule interface {
}

type WorldService interface {

	// 广播消息. 同一个[]byte会发送给每一个人. 调用了之后不可以再修改[]byte中的内容
	// 消耗大, 会一个个shard锁住user map
	Broadcast(pbutil.Buffer)
	BroadcastIgnore(pbutil.Buffer, int64)
	Close()
	FuncHero(int64, HeroWalker) bool
	GetTencentInfo(int64) *shared_proto.TencentInfoProto
	GetUserCloseSender(int64) sender.ClosableSender
	GetUserSender(int64) sender.Sender
	IsDontPush(int64) bool
	IsOnline(int64) bool
	MultiSend([]int64, pbutil.Buffer)
	MultiSendIgnore([]int64, pbutil.Buffer, int64)
	MultiSendMsgs([]int64, []pbutil.Buffer)
	MultiSendMsgsIgnore([]int64, []pbutil.Buffer, int64)
	// 尝试放入用户, 如果用户已在里面, 则返回旧的用户和false, 如果用户放入成功, 则返回nil, ok
	PutConnectedUserIfAbsent(ConnectedUser) (ConnectedUser, bool)
	// 删除用户, 如果是同一个对象的话
	RemoveUserIfSame(ConnectedUser) bool
	Send(int64, pbutil.Buffer)
	SendFunc(int64, MsgFunc)
	SendMsgs(int64, []pbutil.Buffer)
	WalkHero(HeroWalker)
	WalkUser(UserWalker)
}

type XiongNuModule interface {

	// 启动一个定时任务，刷新
	Close()
	// gm 命令开启
	GmStart(HeroController, int64) bool
	// 加入联盟后的处理
	JoinGuild(int64, int64)
	OnHeroOnline(HeroController)
}

type XiongNuService interface {
	AddInfo(xiongnuface.ResistXiongNuInfo)
	GetInfo(int64) xiongnuface.ResistXiongNuInfo
	GetRInfo(int64) xiongnuface.RResistXiongNuInfo
	IsStarted(int64) bool
	IsTodayStarted(int64) bool
	RemoveInfo(xiongnuface.ResistXiongNuInfo)
	ResetDaily(time.Time)
	Save()
	SetTodayStarted(int64)
	TodayJoinMap() *xiongnuinfo.TodayJoinMap
	WalkInfo(xiongnuface.WalkInfoFunc)
	XiongNuInfoMsg(int64) pbutil.Buffer
}

type XuanyuanModule interface {
	AddChallenger(int64, *shared_proto.CombatPlayerProto, uint64)
	Close()
	GetResetTickTime() tickdata.TickTime
	GmReset()
	OnHeroOnline(HeroController)
}

type ZhanJiangModule interface {
}

type ZhengWuModule interface {
}
