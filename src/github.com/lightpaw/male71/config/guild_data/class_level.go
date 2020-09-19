package guild_data

import (
	"github.com/lightpaw/male7/config/i18n"
	"github.com/lightpaw/male7/pb/shared_proto"
)

//gogen:config
type GuildClassLevelData struct {
	_ struct{} `file:"联盟/联盟职位.txt"`
	_ struct{} `proto:"shared_proto.GuildClassLevelProto"`
	_ struct{} `protoconfig:"GuildClassLevel"`

	Level uint64

	Name string // 默认阶级名称

	CorePrestige bool `protofield:"-"` // 是否增加联盟核心声望

	VoteScore uint64 // 弹劾NPC盟主职位票数

	Permission *GuildPermissionData `type:"sub"`
}

//gogen:config
type GuildPermissionData struct {
	_ struct{} `proto:"shared_proto.GuildPermissionProto"`

	InvateOther bool // true表示允许邀请他人入盟

	AgreeJoin bool // true表示允许同意申请入盟

	UpdateText         bool // true表示允许修改联盟宣言
	UpdateInternalText bool // true表示允许修改联盟内部宣言
	UpdateLabel        bool // true表示允许修改联盟联盟标签

	UpdateFriendGuild bool // true表示允许修改联盟友盟
	UpdateEnemyGuild  bool // true表示允许修改联盟敌盟

	UpdateClassName bool // true表示允许修改阶级名称

	UpdateFlagType bool // true表示允许修改联盟旗帜

	UpdateLowerMemberClassLevel bool // true表示允许修改低等级盟友的阶级（变成副帮主之类的）

	UpdateClassTitle bool // true表示允许修改职称

	KickLowerMember bool // true表示允许踢人（低阶级盟友）

	UpdateJoinCondition bool // true表示允许修改入盟条件

	ImpeachNpcLeader bool // true表示允许修改弹劾NPC盟主

	UpgradeLevel      bool // true表示允许升级联盟等级
	UpgradeLevelCdr   bool // true表示允许加速联盟升级
	UpgradeBuilding   bool // true表示允许升级联盟建筑
	UpgradeTechnology bool // true表示允许升级联盟科技

	UpdatePrestigeTarget bool // true表示允许修改声望目标

	OpenResistXiongNu bool // true表示允许开启抗击匈奴

	SendToAllMembers bool // true表示允许发全体消息

	UpgradeTechnologyCdr bool //  true表示允许科技加速
	UpdateName           bool // 修改联盟名称
	UpdateFlagName       bool // 修改联盟简称
	LeaveGuild           bool // 退出联盟
	DismissGuild         bool // 解散联盟
	ChangeLeader         bool // 禅让盟主
	ChangeYinliang       bool // 使用银两
	UpdateMark           bool // 更新联盟标识

	ConveneMember bool // ture表示允许成员召集
	GetOnlineInfo bool // 查看成员在线情况

	Workshop bool // 联盟工坊

	RecommendMcBuild bool // 推荐营建名城
}

func (d *GuildPermissionData) HasPermission(t shared_proto.GuildPermissionType) bool {
	switch t {
	case shared_proto.GuildPermissionType_PermInvateOther:
		return d.InvateOther
	case shared_proto.GuildPermissionType_PermAgreeJoin:
		return d.AgreeJoin
	case shared_proto.GuildPermissionType_PermUpdateText:
		return d.UpdateText
	case shared_proto.GuildPermissionType_PermUpdateInternalText:
		return d.UpdateInternalText
	case shared_proto.GuildPermissionType_PermUpdateLabel:
		return d.UpdateLabel
	case shared_proto.GuildPermissionType_PermUpdateFriendGuild:
		return d.UpdateFriendGuild
	case shared_proto.GuildPermissionType_PermUpdateEnemyGuild:
		return d.UpdateEnemyGuild
	case shared_proto.GuildPermissionType_PermUpdateClassName:
		return d.UpdateClassName
	case shared_proto.GuildPermissionType_PermUpdateFlagType:
		return d.UpdateFlagType
	case shared_proto.GuildPermissionType_PermUpdateLowerMemberClassLevel:
		return d.UpdateLowerMemberClassLevel
	case shared_proto.GuildPermissionType_PermUpdateClassTitle:
		return d.UpdateClassTitle
	case shared_proto.GuildPermissionType_PermKickLowerMember:
		return d.KickLowerMember
	case shared_proto.GuildPermissionType_PermUpdateJoinCondition:
		return d.UpdateJoinCondition
	case shared_proto.GuildPermissionType_PermImpeachNpcLeader:
		return d.ImpeachNpcLeader
	case shared_proto.GuildPermissionType_PermUpgradeLevel:
		return d.UpgradeLevel
	case shared_proto.GuildPermissionType_PermUpgradeLevelCdr:
		return d.UpgradeLevelCdr
	case shared_proto.GuildPermissionType_PermUpgradeBuilding:
		return d.UpgradeBuilding
	case shared_proto.GuildPermissionType_PermUpgradeTechnology:
		return d.UpgradeTechnology
	case shared_proto.GuildPermissionType_PermUpdatePrestigeTarget:
		return d.UpdatePrestigeTarget
	case shared_proto.GuildPermissionType_PermOpenResistXiongNu:
		return d.OpenResistXiongNu
	case shared_proto.GuildPermissionType_PermSendToAllMembers:
		return d.SendToAllMembers
	case shared_proto.GuildPermissionType_PermUpgradeTechnologyCdr:
		return d.UpgradeTechnologyCdr
	case shared_proto.GuildPermissionType_PermUpdateName:
		return d.UpdateName
	case shared_proto.GuildPermissionType_PermUpdateFlagName:
		return d.UpdateFlagName
	case shared_proto.GuildPermissionType_PermLeaveGuild:
		return d.LeaveGuild
	case shared_proto.GuildPermissionType_PermDismissGuild:
		return d.DismissGuild
	case shared_proto.GuildPermissionType_PermChangeLeader:
		return d.ChangeLeader
	case shared_proto.GuildPermissionType_PermUpdateMark:
		return d.UpdateMark
	case shared_proto.GuildPermissionType_PermChangeYinliang:
		return d.ChangeYinliang
	case shared_proto.GuildPermissionType_PermConveneMember:
		return d.ConveneMember
	case shared_proto.GuildPermissionType_PermGetOnlineInfo:
		return d.GetOnlineInfo
	case shared_proto.GuildPermissionType_PermWorkshop:
		return d.Workshop
	case shared_proto.GuildPermissionType_PermRecommendMcBuild:
		return d.RecommendMcBuild
	}
	return false
}

//gogen:config
type GuildPermissionShowData struct {
	_ struct{} `file:"联盟/联盟权限.txt"`
	_ struct{} `proto:"shared_proto.GuildPermissionShowProto"`
	_ struct{} `protoconfig:"GuildPermissionShow"`

	Id uint64 `head:"-,uint64(%s.PermType)" protofield:"-"`

	PermType shared_proto.GuildPermissionType

	Name *i18n.I18nRef

	IsShow bool

	ClassLevel []uint64 `head:"-"`
}

func (d *GuildPermissionShowData) Init(filename string, configs interface {
	GetGuildClassLevelDataArray() []*GuildClassLevelData
}) {

	var cls []uint64
	for _, v := range configs.GetGuildClassLevelDataArray() {
		if v.Permission.HasPermission(d.PermType) {
			cls = append(cls, v.Level)
		}
	}
	d.ClassLevel = cls
}
