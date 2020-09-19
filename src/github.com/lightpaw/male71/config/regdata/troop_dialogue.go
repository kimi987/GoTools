package regdata

import (
	"time"
	"github.com/lightpaw/male7/pb/shared_proto"
)

func GetTroopDialogueId(targetType shared_proto.BaseTargetType, subType uint64) uint64 {
	return uint64(targetType)*10000 + subType + 1 // +1 是因为2个type可能都为0，导致id为0
}

// 对话数据
//gogen:config
type TroopDialogueData struct {
	_ struct{} `file:"地图/部队对话.txt"` // 起一个霸气的名字
	_ struct{} `protogen:"true"`
	_ struct{} `protoimport:"region.proto"`

	Id uint64 `head:"-,GetTroopDialogueId(%s.BaseTargetType%c %s.BaseTargetSubType)"`

	// 城池目标类型
	BaseTargetType shared_proto.BaseTargetType `validator:"string" white:"0" type:"enum"`

	// 目标子类型
	// 多等级怪 shared_proto.MultiLevelNpcType
	BaseTargetSubType uint64 `validator:"uint"`

	// 霸业任务章节
	BayeStage uint64 `validator:"uint" protofield:"-"`

	// 君主等级
	HeroLevel uint64 `validator:"uint" protofield:"-"`

	// 显示初始延时
	FirstDelay time.Duration

	// 延时多久显示下一个，0表示，不重复显示
	NextDelay time.Duration

	// true表示随机选择对话显示，否则表示顺序播放
	RandomText bool

	// 对话列表
	Texts []*TroopDialogueTextData `protofield:",config.U64a2I32a(GetTroopDialogueTextDataKeyArray(%s)),int32"`
}

func (d *TroopDialogueData) IsValid(heroCompleteBayeStage, heroLevel uint64) bool {
	if d.BayeStage > 0 && heroCompleteBayeStage >= d.BayeStage {
		return false
	}

	if d.HeroLevel > 0 && heroLevel >= d.HeroLevel {
		return false
	}

	return true
}

// 飘字
//gogen:config
type TroopDialogueTextData struct {
	_ struct{} `file:"地图/部队对话文字.txt"` // 起一个霸气的名字
	_ struct{} `protogen:"true"`

	Id uint64

	// 头像
	Head string

	// 说话的内容
	Text string

	// 持续的时间
	Duration time.Duration

	// 方向 0-左 1-右
	Direction uint64 `validator:"uint"`
}
