package regdata

import (
	"time"
	"github.com/lightpaw/male7/entity/npcid"
)

const (
	// 可被集结的类型
	AssemblyTypeHero    = 1 // 玩家
	AssemblyTypeXiongNu = 2 // 匈奴
	AssemblyTypeJunTuan = 3 // 军团怪
)

func GetAssemblyTypeByTarget(targetId int64) uint64 {
	if !npcid.IsNpcId(targetId) {
		return AssemblyTypeHero
	} else {
		targetNpcType := npcid.GetNpcIdType(targetId)
		switch targetNpcType {
		case npcid.NpcType_XiongNu:
			return AssemblyTypeXiongNu
		case npcid.NpcType_JunTuan:
			return AssemblyTypeJunTuan
		default:
			return 0
		}
	}
}

//gogen:config
type AssemblyData struct {
	_ struct{} `file:"地图/集结.txt"`
	_ struct{} `protogen:"true"`

	// 集结类型 1-玩家 2-匈奴 3-军团怪
	Id uint64

	// 集结成员个数（包含创建者）
	MemberCount uint64

	// 可选择的集结等待时间
	WaitDuration []time.Duration
}
