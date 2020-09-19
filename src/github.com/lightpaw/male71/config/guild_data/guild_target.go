package guild_data

import (
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"math"
)

// 联盟目标
//gogen:config
type GuildTarget struct {
	_ struct{} `file:"联盟/联盟目标.txt"`
	_ struct{} `proto:"shared_proto.GuildTargetProto"`
	_ struct{} `protoconfig:"GuildTarget"`

	Id   uint64 `validator:"int>0"`    // 目标id，任务顺序按照id来排序
	Name string `validator:"string>0"` // 目标名字
	Desc string `validator:"string>0"` // 目标描述
	Icon string                        // 目标图标

	Group uint64 `protofield:"-"`

	ButtonText string

	TargetType shared_proto.GuildTargetType // 任务目标

	// 联盟升级到x级, 通过 target 去找 GuildLevelProto
	// 声望
	Target uint64 `validator:"uint"` // 目标

	Order uint64 `protofield:"-"`

	OrderAmount uint64 `head:"-" protofield:"-"`
}

func (g *GuildTarget) Init(filename string, datas interface {
	GetGuildClassLevelData(level uint64) *GuildClassLevelData
	GetGuildLevelData(level uint64) *GuildLevelData
}) {
	switch g.TargetType {
	case shared_proto.GuildTargetType_GuildLevelUp:
		check.PanicNotTrue(datas.GetGuildLevelData(g.Target) != nil, "%s 配置联盟目标是 %v ，但是对应的联盟等级没找到 [%d]!", filename, g.TargetType, g.Target)
	case shared_proto.GuildTargetType_PrestigeUp:
		check.PanicNotTrue(g.Target > 0, "%s 配置联盟目标是 %v ，提升声望值必须>0", filename, g.TargetType, g.Target)
	default:
		check.PanicNotTrue(g.Target == 0, "%s 配置联盟目标是 %v ，target必须配置0!", filename, g.TargetType, g.Target)
	}

	g.OrderAmount = (g.Order << 16) | (g.Id & math.MaxUint16)
}

//func (g *GuildTarget) CheckIsComplete(targetType shared_proto.GuildTargetType) bool {
//	return g.CheckIsCompleteWithTarget(targetType, 0)
//}
//
//func (g *GuildTarget) CheckIsCompleteWithTarget(targetType shared_proto.GuildTargetType, target uint64) bool {
//	if g.TargetType != targetType {
//		return false
//	}
//
//	switch g.TargetType {
//	case shared_proto.GuildTargetType_UpgradeClassLevel:
//		return target == g.Target
//	default:
//		return target >= g.Target
//	}
//}
//
//func (g *GuildTarget) InitIsCompleted(guild interface {
//	LevelData() *GuildLevelData
//	MemberCount() int
//	Statue() (realmId int64)
//}) bool {
//	switch g.TargetType {
//	case shared_proto.GuildTargetType_GuildLevelUp:
//		return guild.LevelData().Level >= g.Target
//	case shared_proto.GuildTargetType_MemberCount:
//		return uint64(guild.MemberCount()) >= g.Target
//	case shared_proto.GuildTargetType_PlaceStatue:
//		return guild.Statue() > 0
//	case shared_proto.GuildTargetType_JoinCountry:
//		return false
//	default:
//		return false
//	}
//}
