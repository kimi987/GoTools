package tlog

import (
	"bytes"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/pb/shared_proto"
)

// 服务器状态日志

func (s *TlogService) TlogGameSvrState() {
	s.WriteLog(s.buildGameSvrState())
}

func (s *TlogService) buildGameSvrState() string {
	return s.buildLog("GameSvrState", func() string {
		return s.BuildGameSvrState()
	})
}

func (s *TlogService) BuildGameSvrState() string {

	buf := &bytes.Buffer{}
	buf.WriteString("GameSvrState")
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetLocalAddStr())
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)玩家注册(注册完成时记录)

func (s *TlogService) TlogPlayerRegisterById(heroId int64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerRegisterById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogPlayerRegister(hero)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerRegisterById hero not found")
	}
}

func (s *TlogService) TlogPlayerRegister(heroInfo entity.TlogHero) {
	s.WriteLog(s.buildPlayerRegister(heroInfo))
}

func (s *TlogService) buildPlayerRegister(heroInfo entity.TlogHero) string {
	return s.buildLogHeroTx(heroInfo, "PlayerRegister", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildPlayerRegister(heroInfo, tencentInfo)
	})
}

func (s *TlogService) BuildPlayerRegister(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {

	buf := &bytes.Buffer{}
	buf.WriteString("PlayerRegister")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientVersion)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.RegChannel)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientSoftware)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientHardware)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientTelecom)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientNetwork)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.ScreenWidth)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.ScreenHight)
	buf.WriteString(sep)
	writeF32(buf, tencentInfo.Density)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.CpuHardware)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.Memory)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.GLRender)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.GLVersion)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.DeviceId)
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)玩家登陆(进入游戏时记录)

func (s *TlogService) TlogPlayerLoginById(heroId int64, iGuildID uint64, Citys uint64, FriendsNum uint64, QianChongLevel uint64, BaizhanLevel uint64, MiShiLevel uint64, TopPlayerLevel uint64, TopLevelPlayerId uint64, TopFightPower uint64, TopFightPowerPlayerId uint64, Vip uint64, PlayerCount uint64, TotalFightPower uint64, TopTeamFightPower uint64, Team1 []uint64, TeamFightPower1 uint64, Team2 []uint64, Team1FightPower2 uint64, Team3 []uint64, Team1FightPower3 uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerLoginById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogPlayerLogin(hero, iGuildID, Citys, FriendsNum, QianChongLevel, BaizhanLevel, MiShiLevel, TopPlayerLevel, TopLevelPlayerId, TopFightPower, TopFightPowerPlayerId, Vip, PlayerCount, TotalFightPower, TopTeamFightPower, Team1, TeamFightPower1, Team2, Team1FightPower2, Team3, Team1FightPower3)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerLoginById hero not found")
	}
}

func (s *TlogService) TlogPlayerLogin(heroInfo entity.TlogHero, iGuildID uint64, Citys uint64, FriendsNum uint64, QianChongLevel uint64, BaizhanLevel uint64, MiShiLevel uint64, TopPlayerLevel uint64, TopLevelPlayerId uint64, TopFightPower uint64, TopFightPowerPlayerId uint64, Vip uint64, PlayerCount uint64, TotalFightPower uint64, TopTeamFightPower uint64, Team1 []uint64, TeamFightPower1 uint64, Team2 []uint64, Team1FightPower2 uint64, Team3 []uint64, Team1FightPower3 uint64) {
	s.WriteLog(s.buildPlayerLogin(heroInfo, iGuildID, Citys, FriendsNum, QianChongLevel, BaizhanLevel, MiShiLevel, TopPlayerLevel, TopLevelPlayerId, TopFightPower, TopFightPowerPlayerId, Vip, PlayerCount, TotalFightPower, TopTeamFightPower, Team1, TeamFightPower1, Team2, Team1FightPower2, Team3, Team1FightPower3))
}

func (s *TlogService) buildPlayerLogin(heroInfo entity.TlogHero, iGuildID uint64, Citys uint64, FriendsNum uint64, QianChongLevel uint64, BaizhanLevel uint64, MiShiLevel uint64, TopPlayerLevel uint64, TopLevelPlayerId uint64, TopFightPower uint64, TopFightPowerPlayerId uint64, Vip uint64, PlayerCount uint64, TotalFightPower uint64, TopTeamFightPower uint64, Team1 []uint64, TeamFightPower1 uint64, Team2 []uint64, Team1FightPower2 uint64, Team3 []uint64, Team1FightPower3 uint64) string {
	return s.buildLogHeroTx(heroInfo, "PlayerLogin", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildPlayerLogin(heroInfo, tencentInfo, iGuildID, Citys, FriendsNum, QianChongLevel, BaizhanLevel, MiShiLevel, TopPlayerLevel, TopLevelPlayerId, TopFightPower, TopFightPowerPlayerId, Vip, PlayerCount, TotalFightPower, TopTeamFightPower, Team1, TeamFightPower1, Team2, Team1FightPower2, Team3, Team1FightPower3)
	})
}

func (s *TlogService) BuildPlayerLogin(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, iGuildID uint64, Citys uint64, FriendsNum uint64, QianChongLevel uint64, BaizhanLevel uint64, MiShiLevel uint64, TopPlayerLevel uint64, TopLevelPlayerId uint64, TopFightPower uint64, TopFightPowerPlayerId uint64, Vip uint64, PlayerCount uint64, TotalFightPower uint64, TopTeamFightPower uint64, Team1 []uint64, TeamFightPower1 uint64, Team2 []uint64, Team1FightPower2 uint64, Team3 []uint64, Team1FightPower3 uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("PlayerLogin")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.LoginChannel)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientVersion)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientSoftware)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientHardware)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientTelecom)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientNetwork)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.ScreenWidth)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.ScreenHight)
	buf.WriteString(sep)
	writeF32(buf, tencentInfo.Density)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.CpuHardware)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.Memory)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.GLRender)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.GLVersion)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.DeviceId)
	buf.WriteString(sep)
	writeU64(buf, iGuildID)
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, Citys)
	buf.WriteString(sep)
	writeU64(buf, FriendsNum)
	buf.WriteString(sep)
	writeU64(buf, QianChongLevel)
	buf.WriteString(sep)
	writeU64(buf, BaizhanLevel)
	buf.WriteString(sep)
	writeU64(buf, MiShiLevel)
	buf.WriteString(sep)
	writeU64(buf, TopPlayerLevel)
	buf.WriteString(sep)
	writeU64(buf, TopLevelPlayerId)
	buf.WriteString(sep)
	writeU64(buf, TopFightPower)
	buf.WriteString(sep)
	writeU64(buf, TopFightPowerPlayerId)
	buf.WriteString(sep)
	writeU64(buf, Vip)
	buf.WriteString(sep)
	writeU64(buf, PlayerCount)
	buf.WriteString(sep)
	writeU64(buf, TotalFightPower)
	buf.WriteString(sep)
	writeU64(buf, TopTeamFightPower)
	buf.WriteString(sep)
	writeU64Array(buf, Team1)
	buf.WriteString(sep)
	writeU64(buf, TeamFightPower1)
	buf.WriteString(sep)
	writeU64Array(buf, Team2)
	buf.WriteString(sep)
	writeU64(buf, Team1FightPower2)
	buf.WriteString(sep)
	writeU64Array(buf, Team3)
	buf.WriteString(sep)
	writeU64(buf, Team1FightPower3)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)玩家登出游戏时记录(包括所有登出类型)

func (s *TlogService) TlogPlayerLogoutById(heroId int64, OnlineTime uint64, LogoutType uint64, CurOnlineTime uint64, iGuildID uint64, Citys uint64, FriendsNum uint64, QianChongLevel uint64, BaizhanLevel uint64, MiShiLevel uint64, TopPlayerLevel uint64, TopLevelPlayerId uint64, TopFightPower uint64, TopFightPowerPlayerId uint64, Vip uint64, PlayerCount uint64, TotalFightPower uint64, TopTeamFightPower uint64, Team1 []uint64, TeamFightPower1 uint64, Team2 []uint64, Team1FightPower2 uint64, Team3 []uint64, Team1FightPower3 uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerLogoutById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogPlayerLogout(hero, OnlineTime, LogoutType, CurOnlineTime, iGuildID, Citys, FriendsNum, QianChongLevel, BaizhanLevel, MiShiLevel, TopPlayerLevel, TopLevelPlayerId, TopFightPower, TopFightPowerPlayerId, Vip, PlayerCount, TotalFightPower, TopTeamFightPower, Team1, TeamFightPower1, Team2, Team1FightPower2, Team3, Team1FightPower3)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerLogoutById hero not found")
	}
}

func (s *TlogService) TlogPlayerLogout(heroInfo entity.TlogHero, OnlineTime uint64, LogoutType uint64, CurOnlineTime uint64, iGuildID uint64, Citys uint64, FriendsNum uint64, QianChongLevel uint64, BaizhanLevel uint64, MiShiLevel uint64, TopPlayerLevel uint64, TopLevelPlayerId uint64, TopFightPower uint64, TopFightPowerPlayerId uint64, Vip uint64, PlayerCount uint64, TotalFightPower uint64, TopTeamFightPower uint64, Team1 []uint64, TeamFightPower1 uint64, Team2 []uint64, Team1FightPower2 uint64, Team3 []uint64, Team1FightPower3 uint64) {
	s.WriteLog(s.buildPlayerLogout(heroInfo, OnlineTime, LogoutType, CurOnlineTime, iGuildID, Citys, FriendsNum, QianChongLevel, BaizhanLevel, MiShiLevel, TopPlayerLevel, TopLevelPlayerId, TopFightPower, TopFightPowerPlayerId, Vip, PlayerCount, TotalFightPower, TopTeamFightPower, Team1, TeamFightPower1, Team2, Team1FightPower2, Team3, Team1FightPower3))
}

func (s *TlogService) buildPlayerLogout(heroInfo entity.TlogHero, OnlineTime uint64, LogoutType uint64, CurOnlineTime uint64, iGuildID uint64, Citys uint64, FriendsNum uint64, QianChongLevel uint64, BaizhanLevel uint64, MiShiLevel uint64, TopPlayerLevel uint64, TopLevelPlayerId uint64, TopFightPower uint64, TopFightPowerPlayerId uint64, Vip uint64, PlayerCount uint64, TotalFightPower uint64, TopTeamFightPower uint64, Team1 []uint64, TeamFightPower1 uint64, Team2 []uint64, Team1FightPower2 uint64, Team3 []uint64, Team1FightPower3 uint64) string {
	return s.buildLogHeroTx(heroInfo, "PlayerLogout", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildPlayerLogout(heroInfo, tencentInfo, OnlineTime, LogoutType, CurOnlineTime, iGuildID, Citys, FriendsNum, QianChongLevel, BaizhanLevel, MiShiLevel, TopPlayerLevel, TopLevelPlayerId, TopFightPower, TopFightPowerPlayerId, Vip, PlayerCount, TotalFightPower, TopTeamFightPower, Team1, TeamFightPower1, Team2, Team1FightPower2, Team3, Team1FightPower3)
	})
}

func (s *TlogService) BuildPlayerLogout(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, OnlineTime uint64, LogoutType uint64, CurOnlineTime uint64, iGuildID uint64, Citys uint64, FriendsNum uint64, QianChongLevel uint64, BaizhanLevel uint64, MiShiLevel uint64, TopPlayerLevel uint64, TopLevelPlayerId uint64, TopFightPower uint64, TopFightPowerPlayerId uint64, Vip uint64, PlayerCount uint64, TotalFightPower uint64, TopTeamFightPower uint64, Team1 []uint64, TeamFightPower1 uint64, Team2 []uint64, Team1FightPower2 uint64, Team3 []uint64, Team1FightPower3 uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("PlayerLogout")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeU64(buf, OnlineTime)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientVersion)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientSoftware)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientHardware)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientTelecom)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientNetwork)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.ScreenWidth)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.ScreenHight)
	buf.WriteString(sep)
	writeF32(buf, tencentInfo.Density)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.LoginChannel)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.CpuHardware)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.Memory)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.GLRender)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.GLVersion)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.DeviceId)
	buf.WriteString(sep)
	writeU64(buf, LogoutType)
	buf.WriteString(sep)
	writeU64(buf, CurOnlineTime)
	buf.WriteString(sep)
	writeU64(buf, iGuildID)
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, Citys)
	buf.WriteString(sep)
	writeU64(buf, FriendsNum)
	buf.WriteString(sep)
	writeU64(buf, QianChongLevel)
	buf.WriteString(sep)
	writeU64(buf, BaizhanLevel)
	buf.WriteString(sep)
	writeU64(buf, MiShiLevel)
	buf.WriteString(sep)
	writeU64(buf, TopPlayerLevel)
	buf.WriteString(sep)
	writeU64(buf, TopLevelPlayerId)
	buf.WriteString(sep)
	writeU64(buf, TopFightPower)
	buf.WriteString(sep)
	writeU64(buf, TopFightPowerPlayerId)
	buf.WriteString(sep)
	writeU64(buf, Vip)
	buf.WriteString(sep)
	writeU64(buf, PlayerCount)
	buf.WriteString(sep)
	writeU64(buf, TotalFightPower)
	buf.WriteString(sep)
	writeU64(buf, TopTeamFightPower)
	buf.WriteString(sep)
	writeU64Array(buf, Team1)
	buf.WriteString(sep)
	writeU64(buf, TeamFightPower1)
	buf.WriteString(sep)
	writeU64Array(buf, Team2)
	buf.WriteString(sep)
	writeU64(buf, Team1FightPower2)
	buf.WriteString(sep)
	writeU64Array(buf, Team3)
	buf.WriteString(sep)
	writeU64(buf, Team1FightPower3)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)货币流水(货币获得/消耗时记录)

func (s *TlogService) TlogMoneyFlowById(heroId int64, iMoneyType uint64, AfterMoney uint64, iMoney uint64, Reason uint64, AddOrReduce uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.MoneyFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogMoneyFlow(hero, iMoneyType, AfterMoney, iMoney, Reason, AddOrReduce)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.MoneyFlowById hero not found")
	}
}

func (s *TlogService) TlogMoneyFlow(heroInfo entity.TlogHero, iMoneyType uint64, AfterMoney uint64, iMoney uint64, Reason uint64, AddOrReduce uint64) {
	s.WriteLog(s.buildMoneyFlow(heroInfo, iMoneyType, AfterMoney, iMoney, Reason, AddOrReduce))
}

func (s *TlogService) buildMoneyFlow(heroInfo entity.TlogHero, iMoneyType uint64, AfterMoney uint64, iMoney uint64, Reason uint64, AddOrReduce uint64) string {
	return s.buildLogHeroTx(heroInfo, "MoneyFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildMoneyFlow(heroInfo, tencentInfo, iMoneyType, AfterMoney, iMoney, Reason, AddOrReduce)
	})
}

func (s *TlogService) BuildMoneyFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, iMoneyType uint64, AfterMoney uint64, iMoney uint64, Reason uint64, AddOrReduce uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("MoneyFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, iMoneyType)
	buf.WriteString(sep)
	writeU64(buf, AfterMoney)
	buf.WriteString(sep)
	writeU64(buf, iMoney)
	buf.WriteString(sep)
	writeU64(buf, Reason)
	buf.WriteString(sep)
	writeU64(buf, AddOrReduce)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)道具流水表(道具获得/消耗时记录)

func (s *TlogService) TlogItemFlowById(heroId int64, ItemType uint64, ItemId uint64, Count uint64, AfterCount uint64, Reason uint64, AddOrReduce uint64, QualityLevel uint64, iItemId uint64, ItemDeltaCount int64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.ItemFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogItemFlow(hero, ItemType, ItemId, Count, AfterCount, Reason, AddOrReduce, QualityLevel, iItemId, ItemDeltaCount)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.ItemFlowById hero not found")
	}
}

func (s *TlogService) TlogItemFlow(heroInfo entity.TlogHero, ItemType uint64, ItemId uint64, Count uint64, AfterCount uint64, Reason uint64, AddOrReduce uint64, QualityLevel uint64, iItemId uint64, ItemDeltaCount int64) {
	s.WriteLog(s.buildItemFlow(heroInfo, ItemType, ItemId, Count, AfterCount, Reason, AddOrReduce, QualityLevel, iItemId, ItemDeltaCount))
}

func (s *TlogService) buildItemFlow(heroInfo entity.TlogHero, ItemType uint64, ItemId uint64, Count uint64, AfterCount uint64, Reason uint64, AddOrReduce uint64, QualityLevel uint64, iItemId uint64, ItemDeltaCount int64) string {
	return s.buildLogHeroTx(heroInfo, "ItemFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildItemFlow(heroInfo, tencentInfo, ItemType, ItemId, Count, AfterCount, Reason, AddOrReduce, QualityLevel, iItemId, ItemDeltaCount)
	})
}

func (s *TlogService) BuildItemFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, ItemType uint64, ItemId uint64, Count uint64, AfterCount uint64, Reason uint64, AddOrReduce uint64, QualityLevel uint64, iItemId uint64, ItemDeltaCount int64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("ItemFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, ItemType)
	buf.WriteString(sep)
	writeU64(buf, ItemId)
	buf.WriteString(sep)
	writeU64(buf, Count)
	buf.WriteString(sep)
	writeU64(buf, AfterCount)
	buf.WriteString(sep)
	writeU64(buf, Reason)
	buf.WriteString(sep)
	writeU64(buf, AddOrReduce)
	buf.WriteString(sep)
	writeU64(buf, QualityLevel)
	buf.WriteString(sep)
	writeU64(buf, iItemId)
	buf.WriteString(sep)
	writeI64(buf, ItemDeltaCount)
	buf.WriteString(line)

	str := buf.String()
	return str
}

// （必填)君主等级流水表（君主等级提升时记录）

func (s *TlogService) TlogKingExpFlowById(heroId int64, KingBeforeLevel uint64, KingAfterLevel uint64, KingReason uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.KingExpFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogKingExpFlow(hero, KingBeforeLevel, KingAfterLevel, KingReason)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.KingExpFlowById hero not found")
	}
}

func (s *TlogService) TlogKingExpFlow(heroInfo entity.TlogHero, KingBeforeLevel uint64, KingAfterLevel uint64, KingReason uint64) {
	s.WriteLog(s.buildKingExpFlow(heroInfo, KingBeforeLevel, KingAfterLevel, KingReason))
}

func (s *TlogService) buildKingExpFlow(heroInfo entity.TlogHero, KingBeforeLevel uint64, KingAfterLevel uint64, KingReason uint64) string {
	return s.buildLogHeroTx(heroInfo, "KingExpFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildKingExpFlow(heroInfo, tencentInfo, KingBeforeLevel, KingAfterLevel, KingReason)
	})
}

func (s *TlogService) BuildKingExpFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, KingBeforeLevel uint64, KingAfterLevel uint64, KingReason uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("KingExpFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, KingBeforeLevel)
	buf.WriteString(sep)
	writeU64(buf, KingAfterLevel)
	buf.WriteString(sep)
	writeU64(buf, KingReason)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)主城等级流水表(主城等级提升时记录)

func (s *TlogService) TlogCityExpFlowById(heroId int64, CityBeforeLevel uint64, CityAfterLevel uint64, CityReason uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.CityExpFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogCityExpFlow(hero, CityBeforeLevel, CityAfterLevel, CityReason)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.CityExpFlowById hero not found")
	}
}

func (s *TlogService) TlogCityExpFlow(heroInfo entity.TlogHero, CityBeforeLevel uint64, CityAfterLevel uint64, CityReason uint64) {
	s.WriteLog(s.buildCityExpFlow(heroInfo, CityBeforeLevel, CityAfterLevel, CityReason))
}

func (s *TlogService) buildCityExpFlow(heroInfo entity.TlogHero, CityBeforeLevel uint64, CityAfterLevel uint64, CityReason uint64) string {
	return s.buildLogHeroTx(heroInfo, "CityExpFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildCityExpFlow(heroInfo, tencentInfo, CityBeforeLevel, CityAfterLevel, CityReason)
	})
}

func (s *TlogService) BuildCityExpFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, CityBeforeLevel uint64, CityAfterLevel uint64, CityReason uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("CityExpFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, CityBeforeLevel)
	buf.WriteString(sep)
	writeU64(buf, CityAfterLevel)
	buf.WriteString(sep)
	writeU64(buf, CityReason)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)武将培养流水表(武将升级、成长、册封、转职时记录)

func (s *TlogService) TlogPlayerCultivateFlowById(heroId int64, PlayerId uint64, PlayerOpType uint64, BeforeLevel uint64, AfterLevel uint64, GrowBeforeLevel uint64, GrowAfterLevel uint64, OccupationBefore uint64, OccupationAfter uint64, PositionBefore uint64, PositionAfter uint64, PlayerReason uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerCultivateFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogPlayerCultivateFlow(hero, PlayerId, PlayerOpType, BeforeLevel, AfterLevel, GrowBeforeLevel, GrowAfterLevel, OccupationBefore, OccupationAfter, PositionBefore, PositionAfter, PlayerReason)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerCultivateFlowById hero not found")
	}
}

func (s *TlogService) TlogPlayerCultivateFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerOpType uint64, BeforeLevel uint64, AfterLevel uint64, GrowBeforeLevel uint64, GrowAfterLevel uint64, OccupationBefore uint64, OccupationAfter uint64, PositionBefore uint64, PositionAfter uint64, PlayerReason uint64) {
	s.WriteLog(s.buildPlayerCultivateFlow(heroInfo, PlayerId, PlayerOpType, BeforeLevel, AfterLevel, GrowBeforeLevel, GrowAfterLevel, OccupationBefore, OccupationAfter, PositionBefore, PositionAfter, PlayerReason))
}

func (s *TlogService) buildPlayerCultivateFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerOpType uint64, BeforeLevel uint64, AfterLevel uint64, GrowBeforeLevel uint64, GrowAfterLevel uint64, OccupationBefore uint64, OccupationAfter uint64, PositionBefore uint64, PositionAfter uint64, PlayerReason uint64) string {
	return s.buildLogHeroTx(heroInfo, "PlayerCultivateFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildPlayerCultivateFlow(heroInfo, tencentInfo, PlayerId, PlayerOpType, BeforeLevel, AfterLevel, GrowBeforeLevel, GrowAfterLevel, OccupationBefore, OccupationAfter, PositionBefore, PositionAfter, PlayerReason)
	})
}

func (s *TlogService) BuildPlayerCultivateFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, PlayerId uint64, PlayerOpType uint64, BeforeLevel uint64, AfterLevel uint64, GrowBeforeLevel uint64, GrowAfterLevel uint64, OccupationBefore uint64, OccupationAfter uint64, PositionBefore uint64, PositionAfter uint64, PlayerReason uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("PlayerCultivateFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, PlayerId)
	buf.WriteString(sep)
	writeU64(buf, PlayerOpType)
	buf.WriteString(sep)
	writeU64(buf, BeforeLevel)
	buf.WriteString(sep)
	writeU64(buf, AfterLevel)
	buf.WriteString(sep)
	writeU64(buf, GrowBeforeLevel)
	buf.WriteString(sep)
	writeU64(buf, GrowAfterLevel)
	buf.WriteString(sep)
	writeU64(buf, OccupationBefore)
	buf.WriteString(sep)
	writeU64(buf, OccupationAfter)
	buf.WriteString(sep)
	writeU64(buf, PositionBefore)
	buf.WriteString(sep)
	writeU64(buf, PositionAfter)
	buf.WriteString(sep)
	writeU64(buf, PlayerReason)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)武将穿装备(穿/脱装备成功时记录)

func (s *TlogService) TlogPlayerEquipFlowById(heroId int64, PlayerId uint64, PlayerLevel uint64, EquipId uint64, QualityLevel uint64, SlotIndex uint64, OpType uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerEquipFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogPlayerEquipFlow(hero, PlayerId, PlayerLevel, EquipId, QualityLevel, SlotIndex, OpType)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerEquipFlowById hero not found")
	}
}

func (s *TlogService) TlogPlayerEquipFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, EquipId uint64, QualityLevel uint64, SlotIndex uint64, OpType uint64) {
	s.WriteLog(s.buildPlayerEquipFlow(heroInfo, PlayerId, PlayerLevel, EquipId, QualityLevel, SlotIndex, OpType))
}

func (s *TlogService) buildPlayerEquipFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, EquipId uint64, QualityLevel uint64, SlotIndex uint64, OpType uint64) string {
	return s.buildLogHeroTx(heroInfo, "PlayerEquipFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildPlayerEquipFlow(heroInfo, tencentInfo, PlayerId, PlayerLevel, EquipId, QualityLevel, SlotIndex, OpType)
	})
}

func (s *TlogService) BuildPlayerEquipFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, PlayerId uint64, PlayerLevel uint64, EquipId uint64, QualityLevel uint64, SlotIndex uint64, OpType uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("PlayerEquipFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, PlayerId)
	buf.WriteString(sep)
	writeU64(buf, PlayerLevel)
	buf.WriteString(sep)
	writeU64(buf, EquipId)
	buf.WriteString(sep)
	writeU64(buf, QualityLevel)
	buf.WriteString(sep)
	writeU64(buf, SlotIndex)
	buf.WriteString(sep)
	writeU64(buf, OpType)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)强化装备流水(强化/升星装备操作时)

func (s *TlogService) TlogStrenghEquipmentFlowById(heroId int64, PlayerId uint64, PlayerLevel uint64, EquipId uint64, OpType uint64, IfInherit uint64, BeforeLevel uint64, AfterLevel uint64, BeforeStarLevel uint64, AfterStarLevel uint64, SlotIndex uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.StrenghEquipmentFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogStrenghEquipmentFlow(hero, PlayerId, PlayerLevel, EquipId, OpType, IfInherit, BeforeLevel, AfterLevel, BeforeStarLevel, AfterStarLevel, SlotIndex)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.StrenghEquipmentFlowById hero not found")
	}
}

func (s *TlogService) TlogStrenghEquipmentFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, EquipId uint64, OpType uint64, IfInherit uint64, BeforeLevel uint64, AfterLevel uint64, BeforeStarLevel uint64, AfterStarLevel uint64, SlotIndex uint64) {
	s.WriteLog(s.buildStrenghEquipmentFlow(heroInfo, PlayerId, PlayerLevel, EquipId, OpType, IfInherit, BeforeLevel, AfterLevel, BeforeStarLevel, AfterStarLevel, SlotIndex))
}

func (s *TlogService) buildStrenghEquipmentFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, EquipId uint64, OpType uint64, IfInherit uint64, BeforeLevel uint64, AfterLevel uint64, BeforeStarLevel uint64, AfterStarLevel uint64, SlotIndex uint64) string {
	return s.buildLogHeroTx(heroInfo, "StrenghEquipmentFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildStrenghEquipmentFlow(heroInfo, tencentInfo, PlayerId, PlayerLevel, EquipId, OpType, IfInherit, BeforeLevel, AfterLevel, BeforeStarLevel, AfterStarLevel, SlotIndex)
	})
}

func (s *TlogService) BuildStrenghEquipmentFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, PlayerId uint64, PlayerLevel uint64, EquipId uint64, OpType uint64, IfInherit uint64, BeforeLevel uint64, AfterLevel uint64, BeforeStarLevel uint64, AfterStarLevel uint64, SlotIndex uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("StrenghEquipmentFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, PlayerId)
	buf.WriteString(sep)
	writeU64(buf, PlayerLevel)
	buf.WriteString(sep)
	writeU64(buf, EquipId)
	buf.WriteString(sep)
	writeU64(buf, OpType)
	buf.WriteString(sep)
	writeU64(buf, IfInherit)
	buf.WriteString(sep)
	writeU64(buf, BeforeLevel)
	buf.WriteString(sep)
	writeU64(buf, AfterLevel)
	buf.WriteString(sep)
	writeU64(buf, BeforeStarLevel)
	buf.WriteString(sep)
	writeU64(buf, AfterStarLevel)
	buf.WriteString(sep)
	writeU64(buf, SlotIndex)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)装备宝石流水(宝石镶嵌成功时记录)

func (s *TlogService) TlogEquipmentAddStarFlowById(heroId int64, PlayerId uint64, PlayerLevel uint64, BeforeLevel uint64, AfterLevel uint64, GemSlotIndex uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.EquipmentAddStarFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogEquipmentAddStarFlow(hero, PlayerId, PlayerLevel, BeforeLevel, AfterLevel, GemSlotIndex)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.EquipmentAddStarFlowById hero not found")
	}
}

func (s *TlogService) TlogEquipmentAddStarFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, BeforeLevel uint64, AfterLevel uint64, GemSlotIndex uint64) {
	s.WriteLog(s.buildEquipmentAddStarFlow(heroInfo, PlayerId, PlayerLevel, BeforeLevel, AfterLevel, GemSlotIndex))
}

func (s *TlogService) buildEquipmentAddStarFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, BeforeLevel uint64, AfterLevel uint64, GemSlotIndex uint64) string {
	return s.buildLogHeroTx(heroInfo, "EquipmentAddStarFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildEquipmentAddStarFlow(heroInfo, tencentInfo, PlayerId, PlayerLevel, BeforeLevel, AfterLevel, GemSlotIndex)
	})
}

func (s *TlogService) BuildEquipmentAddStarFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, PlayerId uint64, PlayerLevel uint64, BeforeLevel uint64, AfterLevel uint64, GemSlotIndex uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("EquipmentAddStarFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, PlayerId)
	buf.WriteString(sep)
	writeU64(buf, PlayerLevel)
	buf.WriteString(sep)
	writeU64(buf, BeforeLevel)
	buf.WriteString(sep)
	writeU64(buf, AfterLevel)
	buf.WriteString(sep)
	writeU64(buf, GemSlotIndex)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)联盟流水(玩家有联盟行为时记录)

func (s *TlogService) TlogGuildFlowById(heroId int64, iActType uint64, iGuildID uint64, iGuildLevel uint64, iMemberNum uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.GuildFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogGuildFlow(hero, iActType, iGuildID, iGuildLevel, iMemberNum)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.GuildFlowById hero not found")
	}
}

func (s *TlogService) TlogGuildFlow(heroInfo entity.TlogHero, iActType uint64, iGuildID uint64, iGuildLevel uint64, iMemberNum uint64) {
	s.WriteLog(s.buildGuildFlow(heroInfo, iActType, iGuildID, iGuildLevel, iMemberNum))
}

func (s *TlogService) buildGuildFlow(heroInfo entity.TlogHero, iActType uint64, iGuildID uint64, iGuildLevel uint64, iMemberNum uint64) string {
	return s.buildLogHeroTx(heroInfo, "GuildFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildGuildFlow(heroInfo, tencentInfo, iActType, iGuildID, iGuildLevel, iMemberNum)
	})
}

func (s *TlogService) BuildGuildFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, iActType uint64, iGuildID uint64, iGuildLevel uint64, iMemberNum uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("GuildFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeU64(buf, iActType)
	buf.WriteString(sep)
	writeU64(buf, iGuildID)
	buf.WriteString(sep)
	writeU64(buf, iGuildLevel)
	buf.WriteString(sep)
	writeU64(buf, iMemberNum)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)邮件流水(玩家有邮件行为时记录)

func (s *TlogService) TlogMailFlowById(heroId int64, MailType uint64, MailId string) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.MailFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogMailFlow(hero, MailType, MailId)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.MailFlowById hero not found")
	}
}

func (s *TlogService) TlogMailFlow(heroInfo entity.TlogHero, MailType uint64, MailId string) {
	s.WriteLog(s.buildMailFlow(heroInfo, MailType, MailId))
}

func (s *TlogService) buildMailFlow(heroInfo entity.TlogHero, MailType uint64, MailId string) string {
	return s.buildLogHeroTx(heroInfo, "MailFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildMailFlow(heroInfo, tencentInfo, MailType, MailId)
	})
}

func (s *TlogService) BuildMailFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, MailType uint64, MailId string) string {

	buf := &bytes.Buffer{}
	buf.WriteString("MailFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeU64(buf, MailType)
	buf.WriteString(sep)
	writeString(buf, MailId)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)建筑升级流水(升级建筑成功时记录)

func (s *TlogService) TlogStrenghBuildingFlowById(heroId int64, BuildingID uint64, BeforeLevel uint64, AfterLevel uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.StrenghBuildingFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogStrenghBuildingFlow(hero, BuildingID, BeforeLevel, AfterLevel)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.StrenghBuildingFlowById hero not found")
	}
}

func (s *TlogService) TlogStrenghBuildingFlow(heroInfo entity.TlogHero, BuildingID uint64, BeforeLevel uint64, AfterLevel uint64) {
	s.WriteLog(s.buildStrenghBuildingFlow(heroInfo, BuildingID, BeforeLevel, AfterLevel))
}

func (s *TlogService) buildStrenghBuildingFlow(heroInfo entity.TlogHero, BuildingID uint64, BeforeLevel uint64, AfterLevel uint64) string {
	return s.buildLogHeroTx(heroInfo, "StrenghBuildingFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildStrenghBuildingFlow(heroInfo, tencentInfo, BuildingID, BeforeLevel, AfterLevel)
	})
}

func (s *TlogService) BuildStrenghBuildingFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, BuildingID uint64, BeforeLevel uint64, AfterLevel uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("StrenghBuildingFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, BuildingID)
	buf.WriteString(sep)
	writeU64(buf, BeforeLevel)
	buf.WriteString(sep)
	writeU64(buf, AfterLevel)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)摇钱树照顾流水(摇钱树照顾时记录)

func (s *TlogService) TlogCareFlowById(heroId int64, CareType uint64, pRoleID uint64, pRoleName string, BeforeLevel uint64, AfterLevel uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.CareFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogCareFlow(hero, CareType, pRoleID, pRoleName, BeforeLevel, AfterLevel)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.CareFlowById hero not found")
	}
}

func (s *TlogService) TlogCareFlow(heroInfo entity.TlogHero, CareType uint64, pRoleID uint64, pRoleName string, BeforeLevel uint64, AfterLevel uint64) {
	s.WriteLog(s.buildCareFlow(heroInfo, CareType, pRoleID, pRoleName, BeforeLevel, AfterLevel))
}

func (s *TlogService) buildCareFlow(heroInfo entity.TlogHero, CareType uint64, pRoleID uint64, pRoleName string, BeforeLevel uint64, AfterLevel uint64) string {
	return s.buildLogHeroTx(heroInfo, "CareFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildCareFlow(heroInfo, tencentInfo, CareType, pRoleID, pRoleName, BeforeLevel, AfterLevel)
	})
}

func (s *TlogService) BuildCareFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, CareType uint64, pRoleID uint64, pRoleName string, BeforeLevel uint64, AfterLevel uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("CareFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, CareType)
	buf.WriteString(sep)
	writeU64(buf, pRoleID)
	buf.WriteString(sep)
	writeString(buf, pRoleName)
	buf.WriteString(sep)
	writeU64(buf, BeforeLevel)
	buf.WriteString(sep)
	writeU64(buf, AfterLevel)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)刷新流水(刷新时记录)

func (s *TlogService) TlogRefreshFlowById(heroId int64, BuildingID uint64, RefreshType uint64, ReFreshID uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.RefreshFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogRefreshFlow(hero, BuildingID, RefreshType, ReFreshID)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.RefreshFlowById hero not found")
	}
}

func (s *TlogService) TlogRefreshFlow(heroInfo entity.TlogHero, BuildingID uint64, RefreshType uint64, ReFreshID uint64) {
	s.WriteLog(s.buildRefreshFlow(heroInfo, BuildingID, RefreshType, ReFreshID))
}

func (s *TlogService) buildRefreshFlow(heroInfo entity.TlogHero, BuildingID uint64, RefreshType uint64, ReFreshID uint64) string {
	return s.buildLogHeroTx(heroInfo, "RefreshFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildRefreshFlow(heroInfo, tencentInfo, BuildingID, RefreshType, ReFreshID)
	})
}

func (s *TlogService) BuildRefreshFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, BuildingID uint64, RefreshType uint64, ReFreshID uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("RefreshFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, BuildingID)
	buf.WriteString(sep)
	writeU64(buf, RefreshType)
	buf.WriteString(sep)
	writeU64(buf, ReFreshID)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)加速流水表

func (s *TlogService) TlogSpeedUpFlowById(heroId int64, SpeedUpType uint64, SpeedUpTime uint64, SpeedUpTimeBefore uint64, SpeedUpTimeAfter uint64, Reason uint64, HelperRoleID uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.SpeedUpFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogSpeedUpFlow(hero, SpeedUpType, SpeedUpTime, SpeedUpTimeBefore, SpeedUpTimeAfter, Reason, HelperRoleID)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.SpeedUpFlowById hero not found")
	}
}

func (s *TlogService) TlogSpeedUpFlow(heroInfo entity.TlogHero, SpeedUpType uint64, SpeedUpTime uint64, SpeedUpTimeBefore uint64, SpeedUpTimeAfter uint64, Reason uint64, HelperRoleID uint64) {
	s.WriteLog(s.buildSpeedUpFlow(heroInfo, SpeedUpType, SpeedUpTime, SpeedUpTimeBefore, SpeedUpTimeAfter, Reason, HelperRoleID))
}

func (s *TlogService) buildSpeedUpFlow(heroInfo entity.TlogHero, SpeedUpType uint64, SpeedUpTime uint64, SpeedUpTimeBefore uint64, SpeedUpTimeAfter uint64, Reason uint64, HelperRoleID uint64) string {
	return s.buildLogHeroTx(heroInfo, "SpeedUpFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildSpeedUpFlow(heroInfo, tencentInfo, SpeedUpType, SpeedUpTime, SpeedUpTimeBefore, SpeedUpTimeAfter, Reason, HelperRoleID)
	})
}

func (s *TlogService) BuildSpeedUpFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, SpeedUpType uint64, SpeedUpTime uint64, SpeedUpTimeBefore uint64, SpeedUpTimeAfter uint64, Reason uint64, HelperRoleID uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("SpeedUpFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, SpeedUpType)
	buf.WriteString(sep)
	writeU64(buf, SpeedUpTime)
	buf.WriteString(sep)
	writeU64(buf, SpeedUpTimeBefore)
	buf.WriteString(sep)
	writeU64(buf, SpeedUpTimeAfter)
	buf.WriteString(sep)
	writeU64(buf, Reason)
	buf.WriteString(sep)
	writeU64(buf, HelperRoleID)
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)SNS流水

func (s *TlogService) TlogSnsFlowById(heroId int64, SNSType uint64, pRoleID uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.SnsFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogSnsFlow(hero, SNSType, pRoleID)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.SnsFlowById hero not found")
	}
}

func (s *TlogService) TlogSnsFlow(heroInfo entity.TlogHero, SNSType uint64, pRoleID uint64) {
	s.WriteLog(s.buildSnsFlow(heroInfo, SNSType, pRoleID))
}

func (s *TlogService) buildSnsFlow(heroInfo entity.TlogHero, SNSType uint64, pRoleID uint64) string {
	return s.buildLogHeroTx(heroInfo, "SnsFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildSnsFlow(heroInfo, tencentInfo, SNSType, pRoleID)
	})
}

func (s *TlogService) BuildSnsFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, SNSType uint64, pRoleID uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("SnsFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, SNSType)
	buf.WriteString(sep)
	writeU64(buf, pRoleID)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)新手引导信息表(触发/结束时记录)

func (s *TlogService) TlogGuideFlowById(heroId int64, iGuideID uint64, Status uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.GuideFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogGuideFlow(hero, iGuideID, Status)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.GuideFlowById hero not found")
	}
}

func (s *TlogService) TlogGuideFlow(heroInfo entity.TlogHero, iGuideID uint64, Status uint64) {
	s.WriteLog(s.buildGuideFlow(heroInfo, iGuideID, Status))
}

func (s *TlogService) buildGuideFlow(heroInfo entity.TlogHero, iGuideID uint64, Status uint64) string {
	return s.buildLogHeroTx(heroInfo, "GuideFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildGuideFlow(heroInfo, tencentInfo, iGuideID, Status)
	})
}

func (s *TlogService) BuildGuideFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, iGuideID uint64, Status uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("GuideFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, iGuideID)
	buf.WriteString(sep)
	writeU64(buf, Status)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// 单局结束数据流水(非任务类活动结束时触发日志)

func (s *TlogService) TlogRoundFlowById(heroId int64, BattleType uint64, BattleID uint64, BattleSource uint64, PlayerType uint64, Round uint64, Result uint64, Score uint64, PlayerNum uint64, Player1 uint64, Player2 uint64, Player3 uint64, Player4 uint64, Player5 uint64, Player1Occupation uint64, Player1Num uint64, Player2Occupation uint64, Player2Num uint64, Player3Occupation uint64, Player3Num uint64, Player4Occupation uint64, Player4Num uint64, Player5Occupation uint64, Player5Num uint64, Player1NumLeft uint64, Player2NumLeft uint64, Player3NumLeft uint64, Player4NumLeft uint64, Player5NumLeft uint64, UniqueId uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.RoundFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogRoundFlow(hero, BattleType, BattleID, BattleSource, PlayerType, Round, Result, Score, PlayerNum, Player1, Player2, Player3, Player4, Player5, Player1Occupation, Player1Num, Player2Occupation, Player2Num, Player3Occupation, Player3Num, Player4Occupation, Player4Num, Player5Occupation, Player5Num, Player1NumLeft, Player2NumLeft, Player3NumLeft, Player4NumLeft, Player5NumLeft, UniqueId)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.RoundFlowById hero not found")
	}
}

func (s *TlogService) TlogRoundFlow(heroInfo entity.TlogHero, BattleType uint64, BattleID uint64, BattleSource uint64, PlayerType uint64, Round uint64, Result uint64, Score uint64, PlayerNum uint64, Player1 uint64, Player2 uint64, Player3 uint64, Player4 uint64, Player5 uint64, Player1Occupation uint64, Player1Num uint64, Player2Occupation uint64, Player2Num uint64, Player3Occupation uint64, Player3Num uint64, Player4Occupation uint64, Player4Num uint64, Player5Occupation uint64, Player5Num uint64, Player1NumLeft uint64, Player2NumLeft uint64, Player3NumLeft uint64, Player4NumLeft uint64, Player5NumLeft uint64, UniqueId uint64) {
	s.WriteLog(s.buildRoundFlow(heroInfo, BattleType, BattleID, BattleSource, PlayerType, Round, Result, Score, PlayerNum, Player1, Player2, Player3, Player4, Player5, Player1Occupation, Player1Num, Player2Occupation, Player2Num, Player3Occupation, Player3Num, Player4Occupation, Player4Num, Player5Occupation, Player5Num, Player1NumLeft, Player2NumLeft, Player3NumLeft, Player4NumLeft, Player5NumLeft, UniqueId))
}

func (s *TlogService) buildRoundFlow(heroInfo entity.TlogHero, BattleType uint64, BattleID uint64, BattleSource uint64, PlayerType uint64, Round uint64, Result uint64, Score uint64, PlayerNum uint64, Player1 uint64, Player2 uint64, Player3 uint64, Player4 uint64, Player5 uint64, Player1Occupation uint64, Player1Num uint64, Player2Occupation uint64, Player2Num uint64, Player3Occupation uint64, Player3Num uint64, Player4Occupation uint64, Player4Num uint64, Player5Occupation uint64, Player5Num uint64, Player1NumLeft uint64, Player2NumLeft uint64, Player3NumLeft uint64, Player4NumLeft uint64, Player5NumLeft uint64, UniqueId uint64) string {
	return s.buildLogHeroTx(heroInfo, "RoundFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildRoundFlow(heroInfo, tencentInfo, BattleType, BattleID, BattleSource, PlayerType, Round, Result, Score, PlayerNum, Player1, Player2, Player3, Player4, Player5, Player1Occupation, Player1Num, Player2Occupation, Player2Num, Player3Occupation, Player3Num, Player4Occupation, Player4Num, Player5Occupation, Player5Num, Player1NumLeft, Player2NumLeft, Player3NumLeft, Player4NumLeft, Player5NumLeft, UniqueId)
	})
}

func (s *TlogService) BuildRoundFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, BattleType uint64, BattleID uint64, BattleSource uint64, PlayerType uint64, Round uint64, Result uint64, Score uint64, PlayerNum uint64, Player1 uint64, Player2 uint64, Player3 uint64, Player4 uint64, Player5 uint64, Player1Occupation uint64, Player1Num uint64, Player2Occupation uint64, Player2Num uint64, Player3Occupation uint64, Player3Num uint64, Player4Occupation uint64, Player4Num uint64, Player5Occupation uint64, Player5Num uint64, Player1NumLeft uint64, Player2NumLeft uint64, Player3NumLeft uint64, Player4NumLeft uint64, Player5NumLeft uint64, UniqueId uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("RoundFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, BattleType)
	buf.WriteString(sep)
	writeU64(buf, BattleID)
	buf.WriteString(sep)
	writeU64(buf, BattleSource)
	buf.WriteString(sep)
	writeU64(buf, PlayerType)
	buf.WriteString(sep)
	writeU64(buf, Round)
	buf.WriteString(sep)
	writeU64(buf, Result)
	buf.WriteString(sep)
	writeU64(buf, Score)
	buf.WriteString(sep)
	writeU64(buf, PlayerNum)
	buf.WriteString(sep)
	writeU64(buf, Player1)
	buf.WriteString(sep)
	writeU64(buf, Player2)
	buf.WriteString(sep)
	writeU64(buf, Player3)
	buf.WriteString(sep)
	writeU64(buf, Player4)
	buf.WriteString(sep)
	writeU64(buf, Player5)
	buf.WriteString(sep)
	writeU64(buf, Player1Occupation)
	buf.WriteString(sep)
	writeU64(buf, Player1Num)
	buf.WriteString(sep)
	writeU64(buf, Player2Occupation)
	buf.WriteString(sep)
	writeU64(buf, Player2Num)
	buf.WriteString(sep)
	writeU64(buf, Player3Occupation)
	buf.WriteString(sep)
	writeU64(buf, Player3Num)
	buf.WriteString(sep)
	writeU64(buf, Player4Occupation)
	buf.WriteString(sep)
	writeU64(buf, Player4Num)
	buf.WriteString(sep)
	writeU64(buf, Player5Occupation)
	buf.WriteString(sep)
	writeU64(buf, Player5Num)
	buf.WriteString(sep)
	writeU64(buf, Player1NumLeft)
	buf.WriteString(sep)
	writeU64(buf, Player2NumLeft)
	buf.WriteString(sep)
	writeU64(buf, Player3NumLeft)
	buf.WriteString(sep)
	writeU64(buf, Player4NumLeft)
	buf.WriteString(sep)
	writeU64(buf, Player5NumLeft)
	buf.WriteString(sep)
	writeU64(buf, UniqueId)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必选)百战千军军衔流水表

func (s *TlogService) TlogBaiZhanFlowById(heroId int64, BaiJiangLevelBefore uint64, BaiJiangLevelAfter uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.BaiZhanFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogBaiZhanFlow(hero, BaiJiangLevelBefore, BaiJiangLevelAfter)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.BaiZhanFlowById hero not found")
	}
}

func (s *TlogService) TlogBaiZhanFlow(heroInfo entity.TlogHero, BaiJiangLevelBefore uint64, BaiJiangLevelAfter uint64) {
	s.WriteLog(s.buildBaiZhanFlow(heroInfo, BaiJiangLevelBefore, BaiJiangLevelAfter))
}

func (s *TlogService) buildBaiZhanFlow(heroInfo entity.TlogHero, BaiJiangLevelBefore uint64, BaiJiangLevelAfter uint64) string {
	return s.buildLogHeroTx(heroInfo, "BaiZhanFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildBaiZhanFlow(heroInfo, tencentInfo, BaiJiangLevelBefore, BaiJiangLevelAfter)
	})
}

func (s *TlogService) BuildBaiZhanFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, BaiJiangLevelBefore uint64, BaiJiangLevelAfter uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("BaiZhanFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, BaiJiangLevelBefore)
	buf.WriteString(sep)
	writeU64(buf, BaiJiangLevelAfter)
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (可选)兵种、科技流水表

func (s *TlogService) TlogResearchFlowById(heroId int64, ResearchType uint64, ResearchID uint64, iBeforeVipLevel uint64, iAfterVipLevel uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.ResearchFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogResearchFlow(hero, ResearchType, ResearchID, iBeforeVipLevel, iAfterVipLevel)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.ResearchFlowById hero not found")
	}
}

func (s *TlogService) TlogResearchFlow(heroInfo entity.TlogHero, ResearchType uint64, ResearchID uint64, iBeforeVipLevel uint64, iAfterVipLevel uint64) {
	s.WriteLog(s.buildResearchFlow(heroInfo, ResearchType, ResearchID, iBeforeVipLevel, iAfterVipLevel))
}

func (s *TlogService) buildResearchFlow(heroInfo entity.TlogHero, ResearchType uint64, ResearchID uint64, iBeforeVipLevel uint64, iAfterVipLevel uint64) string {
	return s.buildLogHeroTx(heroInfo, "ResearchFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildResearchFlow(heroInfo, tencentInfo, ResearchType, ResearchID, iBeforeVipLevel, iAfterVipLevel)
	})
}

func (s *TlogService) BuildResearchFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, ResearchType uint64, ResearchID uint64, iBeforeVipLevel uint64, iAfterVipLevel uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("ResearchFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, ResearchType)
	buf.WriteString(sep)
	writeU64(buf, ResearchID)
	buf.WriteString(sep)
	writeU64(buf, iBeforeVipLevel)
	buf.WriteString(sep)
	writeU64(buf, iAfterVipLevel)
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (可选)VIP等级流水表

func (s *TlogService) TlogVipLevelFlowById(heroId int64, iBeforeVipLevel uint64, iAfterVipLevel uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.VipLevelFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogVipLevelFlow(hero, iBeforeVipLevel, iAfterVipLevel)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.VipLevelFlowById hero not found")
	}
}

func (s *TlogService) TlogVipLevelFlow(heroInfo entity.TlogHero, iBeforeVipLevel uint64, iAfterVipLevel uint64) {
	s.WriteLog(s.buildVipLevelFlow(heroInfo, iBeforeVipLevel, iAfterVipLevel))
}

func (s *TlogService) buildVipLevelFlow(heroInfo entity.TlogHero, iBeforeVipLevel uint64, iAfterVipLevel uint64) string {
	return s.buildLogHeroTx(heroInfo, "VipLevelFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildVipLevelFlow(heroInfo, tencentInfo, iBeforeVipLevel, iAfterVipLevel)
	})
}

func (s *TlogService) BuildVipLevelFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, iBeforeVipLevel uint64, iAfterVipLevel uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("VipLevelFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, iBeforeVipLevel)
	buf.WriteString(sep)
	writeU64(buf, iAfterVipLevel)
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)任务流水(任务接取/完成/领取奖励时记录,不包括做了还没交,即还未开启下一个任务时的状态)

func (s *TlogService) TlogTaskFlowById(heroId int64, iTaskType uint64, iTaskID uint64, iState uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.TaskFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogTaskFlow(hero, iTaskType, iTaskID, iState)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.TaskFlowById hero not found")
	}
}

func (s *TlogService) TlogTaskFlow(heroInfo entity.TlogHero, iTaskType uint64, iTaskID uint64, iState uint64) {
	s.WriteLog(s.buildTaskFlow(heroInfo, iTaskType, iTaskID, iState))
}

func (s *TlogService) buildTaskFlow(heroInfo entity.TlogHero, iTaskType uint64, iTaskID uint64, iState uint64) string {
	return s.buildLogHeroTx(heroInfo, "TaskFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildTaskFlow(heroInfo, tencentInfo, iTaskType, iTaskID, iState)
	})
}

func (s *TlogService) BuildTaskFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, iTaskType uint64, iTaskID uint64, iState uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("TaskFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, iTaskType)
	buf.WriteString(sep)
	writeU64(buf, iTaskID)
	buf.WriteString(sep)
	writeU64(buf, iState)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)钓鱼表(点击钓鱼触发日志)

func (s *TlogService) TlogFishFlowById(heroId int64, FishType uint64, iMoney uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.FishFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogFishFlow(hero, FishType, iMoney)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.FishFlowById hero not found")
	}
}

func (s *TlogService) TlogFishFlow(heroInfo entity.TlogHero, FishType uint64, iMoney uint64) {
	s.WriteLog(s.buildFishFlow(heroInfo, FishType, iMoney))
}

func (s *TlogService) buildFishFlow(heroInfo entity.TlogHero, FishType uint64, iMoney uint64) string {
	return s.buildLogHeroTx(heroInfo, "FishFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildFishFlow(heroInfo, tencentInfo, FishType, iMoney)
	})
}

func (s *TlogService) BuildFishFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, FishType uint64, iMoney uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("FishFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, FishType)
	buf.WriteString(sep)
	writeU64(buf, iMoney)
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)武将修炼表(修炼涨经验经验时记录)

func (s *TlogService) TlogPlayerExpDrugFlowById(heroId int64, PlayerId uint64, BeforeLevel uint64, AfterLevel uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerExpDrugFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogPlayerExpDrugFlow(hero, PlayerId, BeforeLevel, AfterLevel)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerExpDrugFlowById hero not found")
	}
}

func (s *TlogService) TlogPlayerExpDrugFlow(heroInfo entity.TlogHero, PlayerId uint64, BeforeLevel uint64, AfterLevel uint64) {
	s.WriteLog(s.buildPlayerExpDrugFlow(heroInfo, PlayerId, BeforeLevel, AfterLevel))
}

func (s *TlogService) buildPlayerExpDrugFlow(heroInfo entity.TlogHero, PlayerId uint64, BeforeLevel uint64, AfterLevel uint64) string {
	return s.buildLogHeroTx(heroInfo, "PlayerExpDrugFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildPlayerExpDrugFlow(heroInfo, tencentInfo, PlayerId, BeforeLevel, AfterLevel)
	})
}

func (s *TlogService) BuildPlayerExpDrugFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, PlayerId uint64, BeforeLevel uint64, AfterLevel uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("PlayerExpDrugFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, PlayerId)
	buf.WriteString(sep)
	writeU64(buf, BeforeLevel)
	buf.WriteString(sep)
	writeU64(buf, AfterLevel)
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)更换上阵武将流水

func (s *TlogService) TlogChangeCaptainFlowById(heroId int64, Team uint64, Place uint64, OldPlayerOccupation uint64, OldPlayerFightPower uint64, OldPlayerID uint64, PlayerOccupation uint64, PlayerFightPower uint64, PlayerID uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.ChangeCaptainFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogChangeCaptainFlow(hero, Team, Place, OldPlayerOccupation, OldPlayerFightPower, OldPlayerID, PlayerOccupation, PlayerFightPower, PlayerID)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.ChangeCaptainFlowById hero not found")
	}
}

func (s *TlogService) TlogChangeCaptainFlow(heroInfo entity.TlogHero, Team uint64, Place uint64, OldPlayerOccupation uint64, OldPlayerFightPower uint64, OldPlayerID uint64, PlayerOccupation uint64, PlayerFightPower uint64, PlayerID uint64) {
	s.WriteLog(s.buildChangeCaptainFlow(heroInfo, Team, Place, OldPlayerOccupation, OldPlayerFightPower, OldPlayerID, PlayerOccupation, PlayerFightPower, PlayerID))
}

func (s *TlogService) buildChangeCaptainFlow(heroInfo entity.TlogHero, Team uint64, Place uint64, OldPlayerOccupation uint64, OldPlayerFightPower uint64, OldPlayerID uint64, PlayerOccupation uint64, PlayerFightPower uint64, PlayerID uint64) string {
	return s.buildLogHeroTx(heroInfo, "ChangeCaptainFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildChangeCaptainFlow(heroInfo, tencentInfo, Team, Place, OldPlayerOccupation, OldPlayerFightPower, OldPlayerID, PlayerOccupation, PlayerFightPower, PlayerID)
	})
}

func (s *TlogService) BuildChangeCaptainFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, Team uint64, Place uint64, OldPlayerOccupation uint64, OldPlayerFightPower uint64, OldPlayerID uint64, PlayerOccupation uint64, PlayerFightPower uint64, PlayerID uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("ChangeCaptainFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, Team)
	buf.WriteString(sep)
	writeU64(buf, Place)
	buf.WriteString(sep)
	writeU64(buf, OldPlayerOccupation)
	buf.WriteString(sep)
	writeU64(buf, OldPlayerFightPower)
	buf.WriteString(sep)
	writeU64(buf, OldPlayerID)
	buf.WriteString(sep)
	writeU64(buf, PlayerOccupation)
	buf.WriteString(sep)
	writeU64(buf, PlayerFightPower)
	buf.WriteString(sep)
	writeU64(buf, PlayerID)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// （必填)聊天流水

func (s *TlogService) TlogChatFlowById(heroId int64, Channel uint64, ChatTo uint64, ChatType uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.ChatFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogChatFlow(hero, Channel, ChatTo, ChatType)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.ChatFlowById hero not found")
	}
}

func (s *TlogService) TlogChatFlow(heroInfo entity.TlogHero, Channel uint64, ChatTo uint64, ChatType uint64) {
	s.WriteLog(s.buildChatFlow(heroInfo, Channel, ChatTo, ChatType))
}

func (s *TlogService) buildChatFlow(heroInfo entity.TlogHero, Channel uint64, ChatTo uint64, ChatType uint64) string {
	return s.buildLogHeroTx(heroInfo, "ChatFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildChatFlow(heroInfo, tencentInfo, Channel, ChatTo, ChatType)
	})
}

func (s *TlogService) BuildChatFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, Channel uint64, ChatTo uint64, ChatType uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("ChatFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, Channel)
	buf.WriteString(sep)
	writeU64(buf, ChatTo)
	buf.WriteString(sep)
	writeU64(buf, ChatType)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)装备熔炼/重铸流水(装备熔炼/重铸成功时记录)

func (s *TlogService) TlogMountRefreshFlowById(heroId int64, OpType uint64, inherit uint64, ItemId uint64, QualityLevel uint64, ItemIndex uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.MountRefreshFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogMountRefreshFlow(hero, OpType, inherit, ItemId, QualityLevel, ItemIndex)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.MountRefreshFlowById hero not found")
	}
}

func (s *TlogService) TlogMountRefreshFlow(heroInfo entity.TlogHero, OpType uint64, inherit uint64, ItemId uint64, QualityLevel uint64, ItemIndex uint64) {
	s.WriteLog(s.buildMountRefreshFlow(heroInfo, OpType, inherit, ItemId, QualityLevel, ItemIndex))
}

func (s *TlogService) buildMountRefreshFlow(heroInfo entity.TlogHero, OpType uint64, inherit uint64, ItemId uint64, QualityLevel uint64, ItemIndex uint64) string {
	return s.buildLogHeroTx(heroInfo, "MountRefreshFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildMountRefreshFlow(heroInfo, tencentInfo, OpType, inherit, ItemId, QualityLevel, ItemIndex)
	})
}

func (s *TlogService) BuildMountRefreshFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, OpType uint64, inherit uint64, ItemId uint64, QualityLevel uint64, ItemIndex uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("MountRefreshFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, OpType)
	buf.WriteString(sep)
	writeU64(buf, inherit)
	buf.WriteString(sep)
	writeU64(buf, ItemId)
	buf.WriteString(sep)
	writeU64(buf, QualityLevel)
	buf.WriteString(sep)
	writeU64(buf, ItemIndex)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)活动玩法流水(参与非任务类活动结束时记录)

func (s *TlogService) TlogGameplayFlowById(heroId int64, GameplayType uint64, Difficulty uint64, DuaringTime uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.GameplayFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogGameplayFlow(hero, GameplayType, Difficulty, DuaringTime)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.GameplayFlowById hero not found")
	}
}

func (s *TlogService) TlogGameplayFlow(heroInfo entity.TlogHero, GameplayType uint64, Difficulty uint64, DuaringTime uint64) {
	s.WriteLog(s.buildGameplayFlow(heroInfo, GameplayType, Difficulty, DuaringTime))
}

func (s *TlogService) buildGameplayFlow(heroInfo entity.TlogHero, GameplayType uint64, Difficulty uint64, DuaringTime uint64) string {
	return s.buildLogHeroTx(heroInfo, "GameplayFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildGameplayFlow(heroInfo, tencentInfo, GameplayType, Difficulty, DuaringTime)
	})
}

func (s *TlogService) BuildGameplayFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, GameplayType uint64, Difficulty uint64, DuaringTime uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("GameplayFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, GameplayType)
	buf.WriteString(sep)
	writeU64(buf, Difficulty)
	buf.WriteString(sep)
	writeU64(buf, DuaringTime)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)武将将魂附身(穿/脱将魂成功时记录)

func (s *TlogService) TlogPlayerHaunterFlowById(heroId int64, PlayerId uint64, PlayerLevel uint64, Soul uint64, QualityLevel uint64, IndexID uint64, OpType uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerHaunterFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogPlayerHaunterFlow(hero, PlayerId, PlayerLevel, Soul, QualityLevel, IndexID, OpType)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.PlayerHaunterFlowById hero not found")
	}
}

func (s *TlogService) TlogPlayerHaunterFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, Soul uint64, QualityLevel uint64, IndexID uint64, OpType uint64) {
	s.WriteLog(s.buildPlayerHaunterFlow(heroInfo, PlayerId, PlayerLevel, Soul, QualityLevel, IndexID, OpType))
}

func (s *TlogService) buildPlayerHaunterFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, Soul uint64, QualityLevel uint64, IndexID uint64, OpType uint64) string {
	return s.buildLogHeroTx(heroInfo, "PlayerHaunterFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildPlayerHaunterFlow(heroInfo, tencentInfo, PlayerId, PlayerLevel, Soul, QualityLevel, IndexID, OpType)
	})
}

func (s *TlogService) BuildPlayerHaunterFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, PlayerId uint64, PlayerLevel uint64, Soul uint64, QualityLevel uint64, IndexID uint64, OpType uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("PlayerHaunterFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, PlayerId)
	buf.WriteString(sep)
	writeU64(buf, PlayerLevel)
	buf.WriteString(sep)
	writeU64(buf, Soul)
	buf.WriteString(sep)
	writeU64(buf, QualityLevel)
	buf.WriteString(sep)
	writeU64(buf, IndexID)
	buf.WriteString(sep)
	writeU64(buf, OpType)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)将魂进阶流水(将魂进阶时)

func (s *TlogService) TlogAdvanceSoulFlowById(heroId int64, PlayerId uint64, PlayerLevel uint64, SoulId uint64, OpType uint64, inherit uint64, BeforeStarLevel uint64, AfterStarLevel uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.AdvanceSoulFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogAdvanceSoulFlow(hero, PlayerId, PlayerLevel, SoulId, OpType, inherit, BeforeStarLevel, AfterStarLevel)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.AdvanceSoulFlowById hero not found")
	}
}

func (s *TlogService) TlogAdvanceSoulFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, SoulId uint64, OpType uint64, inherit uint64, BeforeStarLevel uint64, AfterStarLevel uint64) {
	s.WriteLog(s.buildAdvanceSoulFlow(heroInfo, PlayerId, PlayerLevel, SoulId, OpType, inherit, BeforeStarLevel, AfterStarLevel))
}

func (s *TlogService) buildAdvanceSoulFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, SoulId uint64, OpType uint64, inherit uint64, BeforeStarLevel uint64, AfterStarLevel uint64) string {
	return s.buildLogHeroTx(heroInfo, "AdvanceSoulFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildAdvanceSoulFlow(heroInfo, tencentInfo, PlayerId, PlayerLevel, SoulId, OpType, inherit, BeforeStarLevel, AfterStarLevel)
	})
}

func (s *TlogService) BuildAdvanceSoulFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, PlayerId uint64, PlayerLevel uint64, SoulId uint64, OpType uint64, inherit uint64, BeforeStarLevel uint64, AfterStarLevel uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("AdvanceSoulFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, PlayerId)
	buf.WriteString(sep)
	writeU64(buf, PlayerLevel)
	buf.WriteString(sep)
	writeU64(buf, SoulId)
	buf.WriteString(sep)
	writeU64(buf, OpType)
	buf.WriteString(sep)
	writeU64(buf, inherit)
	buf.WriteString(sep)
	writeU64(buf, BeforeStarLevel)
	buf.WriteString(sep)
	writeU64(buf, AfterStarLevel)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)资源存量(每日0点记录)

func (s *TlogService) TlogResourceStockFlowById(heroId int64, MT_Ingot uint64, MT_MONEY uint64, MT_SOLDIER uint64, MT_COPPER uint64, MT_STONE uint64, MT_CONTRIBUTION uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.ResourceStockFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogResourceStockFlow(hero, MT_Ingot, MT_MONEY, MT_SOLDIER, MT_COPPER, MT_STONE, MT_CONTRIBUTION)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.ResourceStockFlowById hero not found")
	}
}

func (s *TlogService) TlogResourceStockFlow(heroInfo entity.TlogHero, MT_Ingot uint64, MT_MONEY uint64, MT_SOLDIER uint64, MT_COPPER uint64, MT_STONE uint64, MT_CONTRIBUTION uint64) {
	s.WriteLog(s.buildResourceStockFlow(heroInfo, MT_Ingot, MT_MONEY, MT_SOLDIER, MT_COPPER, MT_STONE, MT_CONTRIBUTION))
}

func (s *TlogService) buildResourceStockFlow(heroInfo entity.TlogHero, MT_Ingot uint64, MT_MONEY uint64, MT_SOLDIER uint64, MT_COPPER uint64, MT_STONE uint64, MT_CONTRIBUTION uint64) string {
	return s.buildLogHeroTx(heroInfo, "ResourceStockFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildResourceStockFlow(heroInfo, tencentInfo, MT_Ingot, MT_MONEY, MT_SOLDIER, MT_COPPER, MT_STONE, MT_CONTRIBUTION)
	})
}

func (s *TlogService) BuildResourceStockFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, MT_Ingot uint64, MT_MONEY uint64, MT_SOLDIER uint64, MT_COPPER uint64, MT_STONE uint64, MT_CONTRIBUTION uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("ResourceStockFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, MT_Ingot)
	buf.WriteString(sep)
	writeU64(buf, MT_MONEY)
	buf.WriteString(sep)
	writeU64(buf, MT_SOLDIER)
	buf.WriteString(sep)
	writeU64(buf, MT_COPPER)
	buf.WriteString(sep)
	writeU64(buf, MT_STONE)
	buf.WriteString(sep)
	writeU64(buf, MT_CONTRIBUTION)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)过关斩将流水(挑战成功时记录)

func (s *TlogService) TlogGUOGUANFlowById(heroId int64, PlayerId uint64, PlayerLevel uint64, BeforeGrowth uint64, AfterGrowth uint64, Progress uint64, Honor uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.GUOGUANFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogGUOGUANFlow(hero, PlayerId, PlayerLevel, BeforeGrowth, AfterGrowth, Progress, Honor)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.GUOGUANFlowById hero not found")
	}
}

func (s *TlogService) TlogGUOGUANFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, BeforeGrowth uint64, AfterGrowth uint64, Progress uint64, Honor uint64) {
	s.WriteLog(s.buildGUOGUANFlow(heroInfo, PlayerId, PlayerLevel, BeforeGrowth, AfterGrowth, Progress, Honor))
}

func (s *TlogService) buildGUOGUANFlow(heroInfo entity.TlogHero, PlayerId uint64, PlayerLevel uint64, BeforeGrowth uint64, AfterGrowth uint64, Progress uint64, Honor uint64) string {
	return s.buildLogHeroTx(heroInfo, "GUOGUANFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildGUOGUANFlow(heroInfo, tencentInfo, PlayerId, PlayerLevel, BeforeGrowth, AfterGrowth, Progress, Honor)
	})
}

func (s *TlogService) BuildGUOGUANFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, PlayerId uint64, PlayerLevel uint64, BeforeGrowth uint64, AfterGrowth uint64, Progress uint64, Honor uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("GUOGUANFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, PlayerId)
	buf.WriteString(sep)
	writeU64(buf, PlayerLevel)
	buf.WriteString(sep)
	writeU64(buf, BeforeGrowth)
	buf.WriteString(sep)
	writeU64(buf, AfterGrowth)
	buf.WriteString(sep)
	writeU64(buf, Progress)
	buf.WriteString(sep)
	writeU64(buf, Honor)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)农场流水(农场操作时记录)

func (s *TlogService) TlogFarmFlowById(heroId int64, LotId uint64, OpRpleID string, FarmOpType uint64, BeforeMoneyType uint64, AfterMoneyType uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.FarmFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogFarmFlow(hero, LotId, OpRpleID, FarmOpType, BeforeMoneyType, AfterMoneyType)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.FarmFlowById hero not found")
	}
}

func (s *TlogService) TlogFarmFlow(heroInfo entity.TlogHero, LotId uint64, OpRpleID string, FarmOpType uint64, BeforeMoneyType uint64, AfterMoneyType uint64) {
	s.WriteLog(s.buildFarmFlow(heroInfo, LotId, OpRpleID, FarmOpType, BeforeMoneyType, AfterMoneyType))
}

func (s *TlogService) buildFarmFlow(heroInfo entity.TlogHero, LotId uint64, OpRpleID string, FarmOpType uint64, BeforeMoneyType uint64, AfterMoneyType uint64) string {
	return s.buildLogHeroTx(heroInfo, "FarmFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildFarmFlow(heroInfo, tencentInfo, LotId, OpRpleID, FarmOpType, BeforeMoneyType, AfterMoneyType)
	})
}

func (s *TlogService) BuildFarmFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, LotId uint64, OpRpleID string, FarmOpType uint64, BeforeMoneyType uint64, AfterMoneyType uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("FarmFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, LotId)
	buf.WriteString(sep)
	writeString(buf, OpRpleID)
	buf.WriteString(sep)
	writeU64(buf, FarmOpType)
	buf.WriteString(sep)
	writeU64(buf, BeforeMoneyType)
	buf.WriteString(sep)
	writeU64(buf, AfterMoneyType)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)国家流水(国家操作时记录)

func (s *TlogService) TlogNationalFlowById(heroId int64, FarmOpType uint64, BeforeNationalID uint64, AfterNationalID uint64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.NationalFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogNationalFlow(hero, FarmOpType, BeforeNationalID, AfterNationalID)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.NationalFlowById hero not found")
	}
}

func (s *TlogService) TlogNationalFlow(heroInfo entity.TlogHero, FarmOpType uint64, BeforeNationalID uint64, AfterNationalID uint64) {
	s.WriteLog(s.buildNationalFlow(heroInfo, FarmOpType, BeforeNationalID, AfterNationalID))
}

func (s *TlogService) buildNationalFlow(heroInfo entity.TlogHero, FarmOpType uint64, BeforeNationalID uint64, AfterNationalID uint64) string {
	return s.buildLogHeroTx(heroInfo, "NationalFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildNationalFlow(heroInfo, tencentInfo, FarmOpType, BeforeNationalID, AfterNationalID)
	})
}

func (s *TlogService) BuildNationalFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, FarmOpType uint64, BeforeNationalID uint64, AfterNationalID uint64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("NationalFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, FarmOpType)
	buf.WriteString(sep)
	writeU64(buf, BeforeNationalID)
	buf.WriteString(sep)
	writeU64(buf, AfterNationalID)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)迁城流水(迁城操作时记录)

func (s *TlogService) TlogMoveCitylFlowById(heroId int64, MoveType uint64, BeforelocationX int64, BeforelocationY int64, AfterlocationX int64, AfterlocationY int64) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.MoveCitylFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogMoveCitylFlow(hero, MoveType, BeforelocationX, BeforelocationY, AfterlocationX, AfterlocationY)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.MoveCitylFlowById hero not found")
	}
}

func (s *TlogService) TlogMoveCitylFlow(heroInfo entity.TlogHero, MoveType uint64, BeforelocationX int64, BeforelocationY int64, AfterlocationX int64, AfterlocationY int64) {
	s.WriteLog(s.buildMoveCitylFlow(heroInfo, MoveType, BeforelocationX, BeforelocationY, AfterlocationX, AfterlocationY))
}

func (s *TlogService) buildMoveCitylFlow(heroInfo entity.TlogHero, MoveType uint64, BeforelocationX int64, BeforelocationY int64, AfterlocationX int64, AfterlocationY int64) string {
	return s.buildLogHeroTx(heroInfo, "MoveCitylFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildMoveCitylFlow(heroInfo, tencentInfo, MoveType, BeforelocationX, BeforelocationY, AfterlocationX, AfterlocationY)
	})
}

func (s *TlogService) BuildMoveCitylFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, MoveType uint64, BeforelocationX int64, BeforelocationY int64, AfterlocationX int64, AfterlocationY int64) string {

	buf := &bytes.Buffer{}
	buf.WriteString("MoveCitylFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, MoveType)
	buf.WriteString(sep)
	writeI64(buf, BeforelocationX)
	buf.WriteString(sep)
	writeI64(buf, BeforelocationY)
	buf.WriteString(sep)
	writeI64(buf, AfterlocationX)
	buf.WriteString(sep)
	writeI64(buf, AfterlocationY)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// (必填)答题流水(答题时记录)

func (s *TlogService) TlogAnswerFlowById(heroId int64, QuestionID uint64, Result bool) {
	if npcid.IsNpcId(heroId) {
		logrus.WithField("heroId", heroId).Debug("TlogService.AnswerFlowById npcid")
		return
	}
	if hero := s.heroSnapshotService.GetTlogHero(heroId); hero != nil {
		s.TlogAnswerFlow(hero, QuestionID, Result)
	} else {
		logrus.WithField("heroId", heroId).Debug("TlogService.AnswerFlowById hero not found")
	}
}

func (s *TlogService) TlogAnswerFlow(heroInfo entity.TlogHero, QuestionID uint64, Result bool) {
	s.WriteLog(s.buildAnswerFlow(heroInfo, QuestionID, Result))
}

func (s *TlogService) buildAnswerFlow(heroInfo entity.TlogHero, QuestionID uint64, Result bool) string {
	return s.buildLogHeroTx(heroInfo, "AnswerFlow", func(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildAnswerFlow(heroInfo, tencentInfo, QuestionID, Result)
	})
}

func (s *TlogService) BuildAnswerFlow(heroInfo entity.TlogHero, tencentInfo *shared_proto.TencentInfoProto, QuestionID uint64, Result bool) string {

	buf := &bytes.Buffer{}
	buf.WriteString("AnswerFlow")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeI64(buf, heroInfo.Id())
	buf.WriteString(sep)
	writeString(buf, heroInfo.Name())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.Level())
	buf.WriteString(sep)
	writeU64(buf, heroInfo.BaseLevel())
	buf.WriteString(sep)
	writeU64(buf, QuestionID)
	buf.WriteString(sep)
	writeBool(buf, Result)
	buf.WriteString(sep)
	writeU64(buf, uint64(heroInfo.TotalOnlineTime().Minutes()))
	buf.WriteString(line)

	str := buf.String()
	return str
}

// 账号注册表

func (s *TlogService) TlogAccountRegister(heroId int64) {
	s.WriteLog(s.buildAccountRegister(heroId))
}

func (s *TlogService) buildAccountRegister(heroId int64) string {
	return s.buildLogHeroIdTx(heroId, "AccountRegister", func(heroId int64, tencentInfo *shared_proto.TencentInfoProto) string {
		return s.BuildAccountRegister(heroId, tencentInfo)
	})
}

func (s *TlogService) BuildAccountRegister(heroId int64, tencentInfo *shared_proto.TencentInfoProto) string {

	buf := &bytes.Buffer{}
	buf.WriteString("AccountRegister")
	buf.WriteString(sep)
	writeInt(buf, s.config.GetServerID())
	buf.WriteString(sep)
	writeTime(buf, s.timeService.CurrentTime())
	buf.WriteString(sep)
	writeString(buf, s.config.GetGameAppID())
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.PlatID)
	buf.WriteString(sep)
	writeInt(buf, s.config.GetZoneAreaID())
	buf.WriteString(sep)
	writeString(buf, tencentInfo.OpenID)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientVersion)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.RegChannel)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientSoftware)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientHardware)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientTelecom)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.ClientNetwork)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.ScreenWidth)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.ScreenHight)
	buf.WriteString(sep)
	writeF32(buf, tencentInfo.Density)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.CpuHardware)
	buf.WriteString(sep)
	writeI32(buf, tencentInfo.Memory)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.GLRender)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.GLVersion)
	buf.WriteString(sep)
	writeString(buf, tencentInfo.DeviceId)
	buf.WriteString(line)

	str := buf.String()
	return str
}
