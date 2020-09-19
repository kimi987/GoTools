package xiongnuface

import (
	"github.com/lightpaw/male7/config/xiongnu"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/herosnapshot/snapshotdata"
	"time"
	"github.com/lightpaw/pbutil"
)

type RResistXiongNuInfo interface {
	Data() *xiongnu.ResistXiongNuData
	GuildId() int64
	WipeOutMonster(ctime time.Time, toReduceMorale uint64)
	Morale() uint64
	StartTime() time.Time
	InvadeEndTime() time.Time
	GivePrizeMembers() []int64
	AddHeroFightSoldier(heroId int64, killSoldier, beenKilledSoldier uint64)
}

type ResistXiongNuInfo interface {
	RResistXiongNuInfo
	BaseId() int64
	BaseX() int32
	BaseY() int32
	EncodeClient() *shared_proto.XiongNuProto
	Encode() *server_proto.XiongNuServerProto
	AddGivePrizeMember(id int64)
	SetBase(baseId int64, baseX, baseY int32)
	NeedSyncChange() bool
	Wave() uint64
	IncWave()
	NextWaveTime() time.Time
	Defenders() []int64
	WipeOutMonsterCount() uint64
	EndTime() time.Time
	ResistTime() time.Time
	AddMonsterCount() uint64
	IncAddMonsterCount(toAdd uint64)
	SetAddMonsterCount(toSet uint64)
	IsDefeated() bool
	SetDefeated()
	IsResist() bool
	SetResist()
	EncodeLast(scoreLevel uint64, heroSnapshotGetter func(id int64) *snapshotdata.HeroSnapshot) *shared_proto.LastResistXiongNuProto
	EncodeFight(heroSnapshotGetter func(id int64) *snapshotdata.HeroSnapshot) *shared_proto.ResistXiongNuFightProto
	NextRefreshMonsterTime() time.Time
	SetNextRefreshMonsterTime(toSet time.Time)
	UpdateDefenserFightAmount(defenserId []int64, newFightAmounts, enemyCounts []uint64,
		newExpireTime time.Time, old *DefenserFightAmount) *DefenserFightAmount
	GetDefenserFightAmount() *DefenserFightAmount
}

type WalkInfoFunc func(info ResistXiongNuInfo)

type DefenserFightAmount struct {
	Version int32

	FightAmounts []uint64
	EnemyCounts  []uint64

	DiffVersionMsg pbutil.Buffer
	SameVersionMsg pbutil.Buffer

	ExpireTime time.Time
}
