package entity

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

var (
	_ TlogHero = (*Hero)(nil)
)

type TlogHero interface {
	Id() int64
	Level() uint64
	VipLevel() uint64
	Name() string
	FriendsCount() uint64
	TotalOnlineTime() time.Duration
	BaseLevel() uint64
}

func NewSimpleTlogHeroInfo(heroId int64, name string) *TlogHeroInfo {
	return &TlogHeroInfo{
		HeroId:    heroId,
		name:      name,
		level:     1,
		baseLevel: 1,
	}
}

type TlogHeroInfo struct {
	HeroId          int64
	level           uint64
	baseLevel       uint64
	name            string
	vipLevel        uint64
	totalOnlineTime time.Duration

	OnlineTime                 uint64
	GuildId                    uint64
	OutCityCount               uint64
	FriendCount                uint64
	TowerMaxFloor              uint64
	JunXianLevel               uint64
	MaxSecretTowerId           uint64
	TopLevelCaptainId          uint64
	TopLevelCaptainLevel       uint64
	TopFightCaptainId          uint64
	TopFightCaptainFightAmount uint64
	CaptainCount               uint64
	AllFightAmount             uint64
	TopFightTroopFightAmount   uint64
	TroopFightAmount           [3]uint64
	TroopCaptainIds            [3][]uint64
}

func (hero *Hero) BuildFullTlogHeroInfo(ctime time.Time) *TlogHeroInfo {
	c := &TlogHeroInfo{}
	c.HeroId = hero.Id()
	c.name = hero.Name()
	c.level = hero.Level()
	c.baseLevel = hero.BaseLevel()
	c.vipLevel = hero.VipLevel()

	onlineTime := ctime.Sub(hero.loginTime)
	c.totalOnlineTime += onlineTime

	c.OnlineTime = uint64(onlineTime.Minutes())
	c.GuildId = u64.FromInt64(hero.GuildId())
	c.OutCityCount = u64.FromInt(len(hero.Domestic().OuterCities().CityIds()))
	c.FriendCount = hero.FriendsCount()
	c.TowerMaxFloor = hero.Tower().HistoryMaxFloor()
	c.JunXianLevel = 0
	c.MaxSecretTowerId = hero.SecretTower().MaxOpenSecretTowerId()

	var topFightCaptain *Captain
	var topLevelCaptain *Captain
	var allFightAmount uint64
	for _, c := range hero.Military().Captains() {
		if c == nil {
			continue
		}
		allFightAmount += c.FightAmount()
		if topFightCaptain == nil || c.FightAmount() > topFightCaptain.FightAmount() {
			topFightCaptain = c
		}
		if topLevelCaptain == nil || c.Level() > topLevelCaptain.Level() {
			topLevelCaptain = c
		}
	}

	if topLevelCaptain != nil {
		c.TopLevelCaptainId = topLevelCaptain.Level()
		c.TopLevelCaptainLevel = topLevelCaptain.Id()
	}

	if topFightCaptain != nil {
		c.TopFightCaptainId = topFightCaptain.FightAmount()
		c.TopFightCaptainFightAmount = topFightCaptain.Id()
	}

	c.CaptainCount = u64.FromInt(len(hero.Military().Captains()))
	c.AllFightAmount = allFightAmount

	var topFightTroopFightAmount uint64
	for i, t := range hero.Military().Troops() {
		if t == nil {
			continue
		}
		if i >= len(c.TroopFightAmount) {
			break
		}

		c.TroopFightAmount[i] = t.FullFightAmount()
		c.TroopCaptainIds[i] = t.CaptainIds()
		if t.FullFightAmount() > topFightTroopFightAmount {
			topFightTroopFightAmount = t.FullFightAmount()
		}
	}

	c.TopFightTroopFightAmount = topFightTroopFightAmount

	return c
}

func (c *TlogHeroInfo) Id() int64 {
	return c.HeroId
}
func (c *TlogHeroInfo) Level() uint64 {
	return c.level
}
func (c *TlogHeroInfo) VipLevel() uint64 {
	return c.vipLevel
}
func (c *TlogHeroInfo) Name() string {
	return c.name
}
func (c *TlogHeroInfo) FriendsCount() uint64 {
	return c.FriendCount
}
func (c *TlogHeroInfo) TotalOnlineTime() time.Duration {
	return c.totalOnlineTime
}
func (c *TlogHeroInfo) BaseLevel() uint64 {
	return c.baseLevel
}

type TlogFightContext struct {
	BattleType uint64 // (必填)战斗类型 （详见battletype）
	BattleID   uint64 // (必填)关卡id
	LeaderId   int64  // (必填)玩家组队类型，（队长，队员）
	Round      uint64 // (必填)回合数
}

func NewTlogFightContext(
	BattleType uint64, // (必填)战斗类型 （详见battletype）
	BattleID uint64,   // (必填)关卡id
	LeaderId int64,    // (必填)玩家组队类型，（队长，队员）
	Round uint64,      // (必填)回合数
) *TlogFightContext {
	c := &TlogFightContext{}
	c.BattleID = BattleID
	c.BattleType = BattleType
	c.LeaderId = LeaderId
	c.Round = Round
	return c
}

func NewTencentInfo(heroId int64, p *shared_proto.TencentInfoProto) *TencentInfo {
	c := &TencentInfo{heroID: heroId}
	c.openID = p.OpenID
	c.clientIP = p.ClientIP
	c.clientNetwork = p.ClientNetwork
	c.clientHardware = p.ClientHardware
	c.clientSoftware = p.ClientSoftware
	c.clientVersion = p.ClientVersion
	c.RegChannel = p.RegChannel
	c.ScreenWidth = u64.FromInt32(p.ScreenWidth)
	c.ScreenHight = u64.FromInt32(p.ScreenHight)
	c.Density = float64(p.Density)
	c.LoginChannel = p.LoginChannel
	c.CpuHardware = p.CpuHardware
	c.Memory = u64.FromInt32(p.Memory)
	c.GLRender = p.GLRender
	c.GLVersion = p.GLVersion
	c.DeviceId = p.DeviceId

	return c
}

type TencentInfo struct {
	heroID int64  // 玩家 ID
	openID string // (必填)用户OPENID号

	/* 客户端传来的 */
	platID         uint64 // (必填)ios 0/android 1
	clientIP       string // (必填)客户端IP(后台服务器记录与玩家通信时的IP地址)
	clientTelecom  string // (必填)运营商
	clientNetwork  string // (必填)3G/WIFI/2G
	clientHardware string // (必填)移动终端机型
	clientSoftware string // (可选)移动终端操作系统版本
	clientVersion  string // (必填)客户端版本

	RegChannel   string  // (必填)注册渠道
	ScreenWidth  uint64  // (可选)显示屏宽度
	ScreenHight  uint64  // (可选)显示屏高度
	Density      float64 // (可选)像素密度
	LoginChannel string  // (必填)登录渠道
	CpuHardware  string  // (可选)cpu类型-频率-核数
	Memory       uint64  // (可选)内存信息单位M
	GLRender     string  // (可选)opengl render信息
	GLVersion    string  // (可选)opengl版本信息
	DeviceId     string  // (可选)设备ID
}

func (t *TencentInfo) HeroID() int64 {
	return t.heroID
}

func (t *TencentInfo) PlatID() uint64 {
	return t.platID
}

func (t *TencentInfo) OpenID() string {
	return t.openID
}

func (t *TencentInfo) ClientIP() string {
	return t.clientIP
}

func (t *TencentInfo) ClientTelecom() string {
	return t.clientTelecom
}

func (t *TencentInfo) ClientNetwork() string {
	return t.clientNetwork
}

func (t *TencentInfo) ClientHardware() string {
	return t.clientHardware
}

func (t *TencentInfo) ClientSoftware() string {
	return t.clientSoftware
}

func (t *TencentInfo) ClientVersion() string {
	return t.clientVersion
}
