package db

import (
	"context"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/guilddataservice/sharedguilddata"
	"github.com/lightpaw/male7/service/db/isql"
	"github.com/lightpaw/male7/util/atomic"
)

// autogen from SqlDbService

// adapter
type DbServiceAdapter interface {
	AddChatMsg(context.Context, int64, []byte, *shared_proto.ChatMsgProto) (int64, error)
	AddFarmSteal(context.Context, int64, int64, cb.Cube) error
	AddMcWarGuildRecord(context.Context, uint64, uint64, int64, *shared_proto.McWarTroopsInfoProto) error
	AddMcWarHeroRecord(context.Context, uint64, uint64, int64, *shared_proto.McWarTroopAllRecordProto) error
	AddMcWarRecord(context.Context, uint64, uint64, *shared_proto.McWarFightRecordProto) error
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

//gogen:iface
type DbService struct {
	adapter      DbServiceAdapter
	callingTimes atomic.Uint64
}

func (d *DbService) CallingTimes() uint64 {
	return d.callingTimes.Load()
}

func (d *DbService) AddChatMsg(a0 context.Context, a1 int64, a2 []byte, a3 *shared_proto.ChatMsgProto) (int64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.AddChatMsg(a0, a1, a2, a3)
}
func (d *DbService) AddFarmSteal(a0 context.Context, a1 int64, a2 int64, a3 cb.Cube) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.AddFarmSteal(a0, a1, a2, a3)
}
func (d *DbService) AddMcWarGuildRecord(a0 context.Context, a1 uint64, a2 uint64, a3 int64, a4 *shared_proto.McWarTroopsInfoProto) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.AddMcWarGuildRecord(a0, a1, a2, a3, a4)
}
func (d *DbService) AddMcWarHeroRecord(a0 context.Context, a1 uint64, a2 uint64, a3 int64, a4 *shared_proto.McWarTroopAllRecordProto) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.AddMcWarHeroRecord(a0, a1, a2, a3, a4)
}
func (d *DbService) AddMcWarRecord(a0 context.Context, a1 uint64, a2 uint64, a3 *shared_proto.McWarFightRecordProto) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.AddMcWarRecord(a0, a1, a2, a3)
}
func (d *DbService) Close() error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.Close()
}
func (d *DbService) CreateFarmCube(a0 context.Context, a1 *entity.FarmCube) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.CreateFarmCube(a0, a1)
}
func (d *DbService) CreateFarmLog(a0 context.Context, a1 *shared_proto.FarmStealLogProto) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.CreateFarmLog(a0, a1)
}
func (d *DbService) CreateGuild(a0 context.Context, a1 int64, a2 []byte) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.CreateGuild(a0, a1, a2)
}
func (d *DbService) CreateHero(a0 context.Context, a1 *entity.Hero) (bool, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.CreateHero(a0, a1)
}
func (d *DbService) CreateMail(a0 context.Context, a1 uint64, a2 int64, a3 []byte, a4 bool, a5 bool, a6 bool, a7 int32, a8 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.CreateMail(a0, a1, a2, a3, a4, a5, a6, a7, a8)
}
func (d *DbService) CreateOrder(a0 context.Context, a1 string, a2 uint64, a3 int64, a4 uint32, a5 uint32, a6 int64, a7 uint64, a8 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.CreateOrder(a0, a1, a2, a3, a4, a5, a6, a7, a8)
}
func (d *DbService) DelMcWarHeroRecord(a0 context.Context, a1 int32) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.DelMcWarHeroRecord(a0, a1)
}
func (d *DbService) DelMcWarHeroRecordWithHeroId(a0 context.Context, a1 int32, a2 uint64, a3 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.DelMcWarHeroRecordWithHeroId(a0, a1, a2, a3)
}
func (d *DbService) DeleteChatWindow(a0 context.Context, a1 int64, a2 []byte) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.DeleteChatWindow(a0, a1, a2)
}
func (d *DbService) DeleteGuild(a0 context.Context, a1 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.DeleteGuild(a0, a1)
}
func (d *DbService) DeleteMail(a0 context.Context, a1 uint64, a2 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.DeleteMail(a0, a1, a2)
}
func (d *DbService) DeleteMultiMail(a0 context.Context, a1 int64, a2 []uint64, a3 bool) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.DeleteMultiMail(a0, a1, a2, a3)
}
func (d *DbService) FindSettingsOpen(a0 context.Context, a1 shared_proto.SettingType, a2 []int64) ([]int64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.FindSettingsOpen(a0, a1, a2)
}
func (d *DbService) GMFarmRipe(a0 context.Context, a1 int64, a2 int64, a3 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.GMFarmRipe(a0, a1, a2, a3)
}
func (d *DbService) HeroId(a0 context.Context, a1 string) (int64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.HeroId(a0, a1)
}
func (d *DbService) HeroIdExist(a0 context.Context, a1 int64) (bool, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.HeroIdExist(a0, a1)
}
func (d *DbService) HeroIds(a0 context.Context) ([]int64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.HeroIds(a0)
}
func (d *DbService) HeroNameExist(a0 context.Context, a1 string) (bool, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.HeroNameExist(a0, a1)
}
func (d *DbService) InsertBaiZhanReplay(a0 context.Context, a1 int64, a2 int64, a3 *shared_proto.BaiZhanReplayProto, a4 bool, a5 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.InsertBaiZhanReplay(a0, a1, a2, a3, a4, a5)
}
func (d *DbService) InsertGuildLog(a0 context.Context, a1 int64, a2 *shared_proto.GuildLogProto) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.InsertGuildLog(a0, a1, a2)
}
func (d *DbService) InsertXuanyRecord(a0 context.Context, a1 int64, a2 []byte) (int64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.InsertXuanyRecord(a0, a1, a2)
}
func (d *DbService) IsCollectableMail(a0 context.Context, a1 uint64) (bool, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.IsCollectableMail(a0, a1)
}
func (d *DbService) ListHeroChatMsg(a0 context.Context, a1 []byte, a2 uint64) ([]*shared_proto.ChatMsgProto, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.ListHeroChatMsg(a0, a1, a2)
}
func (d *DbService) ListHeroChatWindow(a0 context.Context, a1 int64) ([]uint64, isql.BytesArray, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.ListHeroChatWindow(a0, a1)
}
func (d *DbService) LoadAllGuild(a0 context.Context) ([]*sharedguilddata.Guild, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadAllGuild(a0)
}
func (d *DbService) LoadAllHeroData(a0 context.Context) ([]*entity.Hero, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadAllHeroData(a0)
}
func (d *DbService) LoadAllRegionHero(a0 context.Context) ([]*entity.Hero, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadAllRegionHero(a0)
}
func (d *DbService) LoadBaiZhanRecord(a0 context.Context, a1 int64, a2 uint64) (isql.BytesArray, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadBaiZhanRecord(a0, a1, a2)
}
func (d *DbService) LoadCanStealCount(a0 context.Context, a1 int64, a2 int64, a3 int64, a4 uint64) (uint64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadCanStealCount(a0, a1, a2, a3, a4)
}
func (d *DbService) LoadCanStealCube(a0 context.Context, a1 int64, a2 int64, a3 int64, a4 uint64) ([]*entity.FarmCube, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadCanStealCube(a0, a1, a2, a3, a4)
}
func (d *DbService) LoadChatMsg(a0 context.Context, a1 int64) (*shared_proto.ChatMsgProto, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadChatMsg(a0, a1)
}
func (d *DbService) LoadCollectMailPrize(a0 context.Context, a1 uint64, a2 int64) (*resdata.Prize, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadCollectMailPrize(a0, a1, a2)
}
func (d *DbService) LoadFarmCube(a0 context.Context, a1 int64, a2 cb.Cube) (*entity.FarmCube, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadFarmCube(a0, a1, a2)
}
func (d *DbService) LoadFarmCubes(a0 context.Context, a1 int64) ([]*entity.FarmCube, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadFarmCubes(a0, a1)
}
func (d *DbService) LoadFarmHarvestCubes(a0 context.Context, a1 int64) ([]*entity.FarmCube, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadFarmHarvestCubes(a0, a1)
}
func (d *DbService) LoadFarmLog(a0 context.Context, a1 int64, a2 uint64) ([]*shared_proto.FarmStealLogProto, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadFarmLog(a0, a1, a2)
}
func (d *DbService) LoadFarmStealCount(a0 context.Context, a1 int64, a2 int64, a3 cb.Cube) (uint64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadFarmStealCount(a0, a1, a2, a3)
}
func (d *DbService) LoadFarmStealCubes(a0 context.Context, a1 int64, a2 int64, a3 uint64) ([]*entity.FarmCube, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadFarmStealCubes(a0, a1, a2, a3)
}
func (d *DbService) LoadGuild(a0 context.Context, a1 int64) (*sharedguilddata.Guild, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadGuild(a0, a1)
}
func (d *DbService) LoadGuildLogs(a0 context.Context, a1 int64, a2 shared_proto.GuildLogType, a3 int64, a4 uint64) ([]*shared_proto.GuildLogProto, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadGuildLogs(a0, a1, a2, a3, a4)
}
func (d *DbService) LoadHero(a0 context.Context, a1 int64) (*entity.Hero, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadHero(a0, a1)
}
func (d *DbService) LoadHeroCount(a0 context.Context) (uint64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadHeroCount(a0)
}
func (d *DbService) LoadHeroListByCountry(a0 context.Context, a1 uint64, a2 uint64, a3 uint64) ([]*entity.Hero, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadHeroListByCountry(a0, a1, a2, a3)
}
func (d *DbService) LoadHeroListByNameAndCountry(a0 context.Context, a1 string, a2 uint64, a3 uint64, a4 uint64) ([]*entity.Hero, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadHeroListByNameAndCountry(a0, a1, a2, a3, a4)
}
func (d *DbService) LoadHeroMailList(a0 context.Context, a1 int64, a2 uint64, a3 int32, a4 int32, a5 int32, a6 int32, a7 int32, a8 int32, a9 uint64) ([]*shared_proto.MailProto, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadHeroMailList(a0, a1, a2, a3, a4, a5, a6, a7, a8, a9)
}
func (d *DbService) LoadHerosByName(a0 context.Context, a1 string, a2 uint64, a3 uint64) ([]*entity.Hero, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadHerosByName(a0, a1, a2, a3)
}
func (d *DbService) LoadJoinedMcWarId(a0 context.Context, a1 int64) (*entity.JoinedMcWarIds, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadJoinedMcWarId(a0, a1)
}
func (d *DbService) LoadKey(a0 context.Context, a1 server_proto.Key) ([]byte, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadKey(a0, a1)
}
func (d *DbService) LoadMail(a0 context.Context, a1 uint64) ([]byte, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadMail(a0, a1)
}
func (d *DbService) LoadMailCountHasPrizeNotCollected(a0 context.Context, a1 int64) (int, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadMailCountHasPrizeNotCollected(a0, a1)
}
func (d *DbService) LoadMailCountHasReportNotReaded(a0 context.Context, a1 int64, a2 int32) (int, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadMailCountHasReportNotReaded(a0, a1, a2)
}
func (d *DbService) LoadMailCountNoReportNotReaded(a0 context.Context, a1 int64) (int, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadMailCountNoReportNotReaded(a0, a1)
}
func (d *DbService) LoadMcWarGuildRecord(a0 context.Context, a1 uint64, a2 uint64, a3 int64) (*shared_proto.McWarTroopsInfoProto, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadMcWarGuildRecord(a0, a1, a2, a3)
}
func (d *DbService) LoadMcWarHeroRecord(a0 context.Context, a1 uint64, a2 uint64, a3 int64) (*shared_proto.McWarTroopAllRecordProto, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadMcWarHeroRecord(a0, a1, a2, a3)
}
func (d *DbService) LoadMcWarRecord(a0 context.Context, a1 uint64, a2 uint64) (*shared_proto.McWarFightRecordProto, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadMcWarRecord(a0, a1, a2)
}
func (d *DbService) LoadNoGuildHeroListByName(a0 context.Context, a1 string, a2 uint64, a3 uint64) ([]*entity.Hero, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadNoGuildHeroListByName(a0, a1, a2, a3)
}
func (d *DbService) LoadRecommendHeros(a0 context.Context, a1 bool, a2 uint64, a3 uint64, a4 uint64, a5 uint64, a6 int64) ([]*entity.Hero, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadRecommendHeros(a0, a1, a2, a3, a4, a5, a6)
}
func (d *DbService) LoadUnreadChatCount(a0 context.Context, a1 int64) (uint64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadUnreadChatCount(a0, a1)
}
func (d *DbService) LoadUserMisc(a0 context.Context, a1 int64) (*server_proto.UserMiscProto, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadUserMisc(a0, a1)
}
func (d *DbService) LoadXuanyRecord(a0 context.Context, a1 int64, a2 int64, a3 bool) ([]int64, isql.BytesArray, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.LoadXuanyRecord(a0, a1, a2, a3)
}
func (d *DbService) MaxGuildId(a0 context.Context) (int64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.MaxGuildId(a0)
}
func (d *DbService) MaxMailId(a0 context.Context) (uint64, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.MaxMailId(a0)
}
func (d *DbService) OrderExist(a0 context.Context, a1 string) (bool, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.OrderExist(a0, a1)
}
func (d *DbService) PlantFarmCube(a0 context.Context, a1 int64, a2 cb.Cube, a3 int64, a4 int64, a5 uint64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.PlantFarmCube(a0, a1, a2, a3, a4, a5)
}
func (d *DbService) ReadChat(a0 context.Context, a1 int64, a2 []byte) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.ReadChat(a0, a1, a2)
}
func (d *DbService) ReadMultiMail(a0 context.Context, a1 int64, a2 []uint64, a3 bool) (*resdata.Prize, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.ReadMultiMail(a0, a1, a2, a3)
}
func (d *DbService) RemoveChatMsg(a0 context.Context, a1 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.RemoveChatMsg(a0, a1)
}
func (d *DbService) RemoveFarmCube(a0 context.Context, a1 int64, a2 cb.Cube) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.RemoveFarmCube(a0, a1, a2)
}
func (d *DbService) RemoveFarmLog(a0 context.Context, a1 int32) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.RemoveFarmLog(a0, a1)
}
func (d *DbService) RemoveFarmSteal(a0 context.Context, a1 int64, a2 []cb.Cube) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.RemoveFarmSteal(a0, a1, a2)
}
func (d *DbService) ResetConflictFarmCubes(a0 context.Context, a1 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.ResetConflictFarmCubes(a0, a1)
}
func (d *DbService) ResetFarmCubes(a0 context.Context, a1 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.ResetFarmCubes(a0, a1)
}
func (d *DbService) SaveFarmCube(a0 context.Context, a1 *entity.FarmCube) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.SaveFarmCube(a0, a1)
}
func (d *DbService) SaveGuild(a0 context.Context, a1 int64, a2 []byte) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.SaveGuild(a0, a1, a2)
}
func (d *DbService) SaveHero(a0 context.Context, a1 *entity.Hero) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.SaveHero(a0, a1)
}
func (d *DbService) SaveKey(a0 context.Context, a1 server_proto.Key, a2 []byte) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.SaveKey(a0, a1, a2)
}
func (d *DbService) SetFarmRipeTime(a0 context.Context, a1 int64, a2 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.SetFarmRipeTime(a0, a1, a2)
}
func (d *DbService) UpdateChatMsg(a0 context.Context, a1 int64, a2 *shared_proto.ChatMsgProto) bool {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateChatMsg(a0, a1, a2)
}
func (d *DbService) UpdateChatWindow(a0 context.Context, a1 int64, a2 []byte, a3 []byte, a4 bool, a5 int32, a6 bool) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateChatWindow(a0, a1, a2, a3, a4, a5, a6)
}
func (d *DbService) UpdateFarmCubeRipeTime(a0 context.Context, a1 int64, a2 cb.Cube, a3 int64, a4 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateFarmCubeRipeTime(a0, a1, a2, a3, a4)
}
func (d *DbService) UpdateFarmCubeState(a0 context.Context, a1 int64, a2 cb.Cube, a3 int64, a4 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateFarmCubeState(a0, a1, a2, a3, a4)
}
func (d *DbService) UpdateFarmStealTimes(a0 context.Context, a1 int64, a2 []cb.Cube) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateFarmStealTimes(a0, a1, a2)
}
func (d *DbService) UpdateHeroGuildId(a0 context.Context, a1 int64, a2 int64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateHeroGuildId(a0, a1, a2)
}
func (d *DbService) UpdateHeroName(a0 context.Context, a1 int64, a2 string, a3 string) bool {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateHeroName(a0, a1, a2, a3)
}
func (d *DbService) UpdateHeroOfflineBoolIfExpected(a0 context.Context, a1 int64, a2 isql.OfflineBool, a3 bool, a4 bool) (bool, error) {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateHeroOfflineBoolIfExpected(a0, a1, a2, a3, a4)
}
func (d *DbService) UpdateMailCollected(a0 context.Context, a1 uint64, a2 int64, a3 bool) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateMailCollected(a0, a1, a2, a3)
}
func (d *DbService) UpdateMailKeep(a0 context.Context, a1 uint64, a2 int64, a3 bool) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateMailKeep(a0, a1, a2, a3)
}
func (d *DbService) UpdateMailRead(a0 context.Context, a1 uint64, a2 int64, a3 bool) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateMailRead(a0, a1, a2, a3)
}
func (d *DbService) UpdateSettings(a0 context.Context, a1 int64, a2 uint64) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateSettings(a0, a1, a2)
}
func (d *DbService) UpdateUserMisc(a0 context.Context, a1 int64, a2 *server_proto.UserMiscProto) error {
	d.callingTimes.Inc()
	defer d.callingTimes.Dec()

	return d.adapter.UpdateUserMisc(a0, a1, a2)
}
