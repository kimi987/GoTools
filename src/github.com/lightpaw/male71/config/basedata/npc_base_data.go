package basedata

import (
	"github.com/lightpaw/male7/config/monsterdata"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/entity/npcid"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/i64"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

//gogen:config
type NpcBaseData struct {
	_ struct{} `file:"地图/野怪基础数据.txt"`
	_ struct{} `proto:"shared_proto.NpcBaseDataProto"`
	_ struct{} `protoconfig:"NpcBaseData"`

	Id uint64

	Name string `default:" " validator:"string"`

	// npc怪物
	Npc *monsterdata.MonsterMasterData

	// npc主城等级
	BaseLevel uint64

	// 外观配置
	Model string `default:" "`

	// 防守外观
	DefModel string `default:" " validator:"string"`

	// 繁荣度
	ProsperityCapcity uint64

	RobMaxDuration time.Duration `default:"0s"` // 持续掠夺最大时间

	// 繁荣度损失速度
	LostProsperityDuration    time.Duration `default:"0s"`
	LostProsperityPerDuration uint64        `validator:"uint" protofield:"-"`             // 每个duration减少的繁荣度
	FirstLoseProsperity       uint64        `default:"0" validator:"uint" protofield:"-"` // 第一次是否损失繁荣度

	TickDuration time.Duration `default:"0s"`
	TickIcon     []string      `default:""`

	TickPrize  *resdata.Prize `default:"nullable" protofield:"-"` // 持续掠夺获得的奖励
	FirstPrize *resdata.Prize `default:"nullable" protofield:"-"` // 抢第一下获得的奖励

	TickPlunder  *resdata.Plunder `default:"nullable" protofield:"-"` // 持续掠夺掉落的奖励
	FirstPlunder *resdata.Plunder `default:"nullable" protofield:"-"` // 抢第一下掉落的奖励

	TickConditionPlunder  *resdata.ConditionPlunder `default:"nullable" protofield:"-"`
	FirstConditionPlunder *resdata.ConditionPlunder `default:"nullable" protofield:"-"`

	ShowPrize    *resdata.Prize `default:"nullable"` // 展示奖励
	ShowSubPrize *resdata.Prize `default:"nullable"` // 展示奖励

	// 持续掠夺最大个数
	MaxRobbers uint64 `validator:"uint" protofield:"-"`

	DestroyWhenLose bool `default:"false" protofield:"-"`

	defenserBytes []byte
}

func (t *NpcBaseData) Init(filename string) {
	check.PanicNotTrue(t.Id <= npcid.NpcDataMask, "npc城池的配置数据的id最大不能超过 %d, id: %d", npcid.NpcDataMask, t.Id)

	t.defenserBytes = must.Marshal(t.encodeBaseDefenserProto())
}

func (t *NpcBaseData) GetDefenserBytes() []byte {
	return t.defenserBytes
}

func (t *NpcBaseData) encodeBaseDefenserProto() *shared_proto.BaseDefenserProto {
	proto := &shared_proto.BaseDefenserProto{}

	for i, c := range t.Npc.Captains {
		if c == nil {
			continue
		}

		proto.CaptainIndex = append(proto.CaptainIndex, imath.Int32(i+1))
		proto.Captains = append(proto.Captains, c.EncodeCaptainInfo())
	}

	return proto
}

func (m *NpcBaseData) EncodeSnapshot(id int64, realmId int64, baseX, baseY int) *shared_proto.HeroBasicSnapshotProto {
	proto := m.Npc.EncodeSnapshot(id)

	proto.BaseRegion = i64.Int32(realmId)
	proto.BaseLevel = u64.Int32(m.BaseLevel)
	proto.BaseX = imath.Int32(baseX)
	proto.BaseY = imath.Int32(baseY)

	return proto
}
