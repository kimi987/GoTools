package singleton

import (
	"github.com/lightpaw/male7/config/captain"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/util/sortkeys"
	"time"
)

//gogen:config
type MilitaryConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"军事/军事杂项.txt"`
	_ struct{} `proto:"shared_proto.MilitaryConfigProto"`
	_ struct{} `protoconfig:"MilitaryConfig"`

	GenSeekCount   uint64        `validator:"int>0" default:"6"`
	SeekDuration   time.Duration `default:"1h" parser:"time.ParseDuration"`
	SeekMaxTimes   uint64        `validator:"int>0" default:"15"`
	DefenserCount  uint64        `validator:"int>0" default:"5"`
	FireLevelLimit uint64        `validator:"int>0" default:"15"`

	CombatRes           *scene.CombatScene `default:"Battle_Field_1" protofield:"-"`       // 地图默认资源
	CombatXLen          uint64             `validator:"int>0" default:"10" protofield:"-"` // 地图XY
	CombatYLen          uint64             `validator:"int>0" default:"5" protofield:"-"`
	CombatMaxRound      uint64             `validator:"int>0" default:"100" protofield:"-"`     // 最大回合数
	CombatCoef          float64            `validator:"float64>0" default:"10" protofield:"-"`  // 战斗系数
	CombatCritRate      float64            `validator:"float64>0" default:"0.3" protofield:"-"` // 暴击几率
	CombatRestraintRate float64            `validator:"float64>0" default:"0.2" protofield:"-"` // 克制技伤害系数

	CombatScorePercent []uint64 `validator:"float64>0" default:"33,66,90" protofield:"-"` // 评分

	TrainingHeroLevel  []uint64 `head:"-"`
	TrainingInitLevel  []uint64 `head:"-"`
	TrainingLevelLimit []uint64 `head:"-"`

	TrainingMaxDuration time.Duration `default:"8h"`
	CaptainInitTrainExp uint64        `default:"1000"`

	MinWallAttackRound          uint64  `validator:"int>0" default:"2" protofield:"-"`       // 最小城墙攻击轮次
	MaxWallAttachFixDamageRound uint64  `validator:"int>0" default:"4" protofield:"-"`       // 最大城墙攻击轮次
	MaxWallBeenHurtPercent      float64 `validator:"float64>0" default:"0.1" protofield:"-"` // 城墙受伤最大百分比

	// 队伍解锁等级(君主等级)，纯展示用（当那个队伍没人的时候展示出来，如果队伍有数据，显示队伍出来，主要防止配表错误）
	TroopsUnlockLevel []uint64 `head:"-"`

	CaptainSeekerCandidateCount uint64 `default:"5" protofield:"-"` // 招募武将--候选武将数

	allCaptainHeads []string
}

func (c *MilitaryConfig) Init(filename string, configs interface {
	GetHeroLevelSubDataArray() []*data.HeroLevelSubData
	GetCaptainDataArray() []*captain.CaptainData
	GetIconArray() []*icon.Icon
}) {

	prevCount := 0
	troopsCount := uint64(0)
	for _, data := range configs.GetHeroLevelSubDataArray() {
		if prevCount < len(data.CaptainTrainingLevel) {
			for i := prevCount; i < len(data.CaptainTrainingLevel); i++ {
				c.TrainingHeroLevel = append(c.TrainingHeroLevel, data.Level)
				c.TrainingInitLevel = append(c.TrainingInitLevel, data.CaptainTrainingLevel[i])
				c.TrainingLevelLimit = append(c.TrainingLevelLimit, data.CaptainTrainingLevelLimit[i])
			}
		}

		prevCount = len(data.CaptainTrainingLevel)

		if troopsCount < data.TroopsCount {
			troopsCount = data.TroopsCount

			c.TroopsUnlockLevel = append(c.TroopsUnlockLevel, data.Level)
		}
	}

	sortkeys.Uint64s(c.CombatScorePercent)

	c.allCaptainHeads = make([]string, 0)
	for _, ic := range configs.GetIconArray() {
		if ic.CaptainHead {
			c.allCaptainHeads = append(c.allCaptainHeads, ic.Id)
		}
	}

	for _, captain := range configs.GetCaptainDataArray() {
		captain.SetInitTrainExp(c.CaptainInitTrainExp)
	}
}

func (c *MilitaryConfig) AllCaptainHeads() []string {
	return c.allCaptainHeads
}
