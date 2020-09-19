package towerdata

import (
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/config/scene"
	"github.com/lightpaw/male7/gen/pb/tower"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
)

//gogen:config
type TowerData struct {
	_ struct{} `file:"系统模块/千重楼.txt"`
	_ struct{} `proto:"shared_proto.TowerDataProto"`
	_ struct{} `protoconfig:"Tower"`

	// 楼层，从1开始
	Floor uint64 `key:"1"`

	// 首胜奖励
	FirstPassPrize      *resdata.Prize
	FirstPassPrizeProto *shared_proto.PrizeProto `head:"-" protofield:"-"`
	FirstPassPrizeBytes []byte                   `head:"-" protofield:"-"`

	// 展示奖励
	ShowPrize *resdata.Prize

	// 掉落奖励
	Plunder *resdata.Plunder `protofield:"-"`

	// 重楼奖励
	BoxPrize           *resdata.Prize `default:"nullable"`
	CollectBoxPrizeMsg pbutil.Buffer  `head:"-" protofield:"-"`

	// 是否是重点楼层
	CheckPoint bool `protofield:"-"`

	// 防守武将
	Monster *monsterdata.MonsterMasterData

	Desc string

	// 战斗场景
	CombatScene *scene.CombatScene `protofield:"-"`

	nextFloor *TowerData `head:"-" protofield:"-"`

	// 解锁的密室
	UnlockSecretTower *SecretTowerData `head:"-" protofield:"UnlockSecretTowerId,config.U64ToI32(%s.Id)"`
}

func (t *TowerData) Init(filename string, dataMap map[uint64]*TowerData) {

	t.FirstPassPrizeProto = t.FirstPassPrize.Encode4Init()
	t.FirstPassPrizeBytes = must.Marshal(t.FirstPassPrizeProto)

	if t.Floor > 1 {
		prevFloor := dataMap[t.Floor-1]
		check.PanicNotTrue(prevFloor != nil, "%s 千重楼配置数据，没找到%v层的数据，层数必须从1开始连续配置", filename, t.Floor-1)
		prevFloor.nextFloor = t
	}

	if t.BoxPrize != nil {
		t.CollectBoxPrizeMsg = tower.NewS2cCollectBoxMsg(u64.Int32(t.Floor)).Static()
	}

}

func (t *TowerData) NextFloor() *TowerData {
	return t.nextFloor
}

func (t *TowerData) GenChallengePrize() *resdata.Prize {
	return t.Plunder.Try()
}
