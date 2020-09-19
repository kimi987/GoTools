package npcid

import (
	"github.com/lightpaw/male7/util/atomic"
	"math"
	"github.com/lightpaw/male7/pb/shared_proto"
)

// 负数都属于Npc
func IsNpcId(id int64) bool {
	return id < 0
}

func IsHomeNpcId(id int64) bool {
	return IsNpcId(id) && GetNpcIdType(id) == NpcType_HomeNpc
}

func IsMonsterNpcId(id int64) bool {
	return IsNpcId(id) && GetNpcIdType(id) == NpcType_Monster
}

func IsMultiLevelMonsterNpcId(id int64) bool {
	return IsNpcId(id) && GetNpcIdType(id) == NpcType_MultiLevelMonster
}

func IsXiongNuNpcId(id int64) bool {
	return IsNpcId(id) && GetNpcIdType(id) == NpcType_XiongNu
}

func IsBaoZangNpcId(id int64) bool {
	return IsNpcId(id) && GetNpcIdType(id) == NpcType_BaoZang
}

func IsJunTuanNpcId(id int64) bool {
	return IsNpcId(id) && GetNpcIdType(id) == NpcType_JunTuan
}

type NpcIdGen interface {
	Next(dataId uint64) int64
	Sequence() uint64
}

func NewNpcIdGen(id uint64, t NpcType) NpcIdGen {
	return &npc_id_gen{
		t:   t,
		gen: atomic.NewUint64(id),
	}
}

const negativeBit = 1 << 63
const negativeMask = negativeBit - 1

func toNpcId(uid uint64) int64 {
	return int64(uid | negativeBit)
}

func toUid(uid int64) uint64 {
	return uint64(uid & negativeMask)
}

type npc_id_gen struct {
	t   NpcType
	gen *atomic.Uint64
}

func (g *npc_id_gen) Next(dataId uint64) int64 {
	return GetNpcId(g.gen.Inc(), dataId, g.t)
}

func (g *npc_id_gen) Sequence() uint64 {
	return g.gen.Load()
}

type NpcType = int32

const (
	NpcDataMask = 1<<25 - 1

	MaxNpcType = 1<<5 - 1 // 31 NpcType的最大id

	// Npc行营入侵
	NpcType_HomeNpc = NpcType(shared_proto.BaseTargetType_NpcHome)
	// 多等级野怪
	NpcType_MultiLevelMonster = NpcType(shared_proto.BaseTargetType_NpcMultiLevelMonster)
	// 野怪
	NpcType_Monster = NpcType(shared_proto.BaseTargetType_NpcMonster)
	// 联盟Npc
	NpcType_Guild = NpcType(shared_proto.BaseTargetType_NpcGuild)
	// 抗击匈奴Npc
	NpcType_XiongNu = NpcType(shared_proto.BaseTargetType_NpcXiongNu)
	// 宝藏Npc
	NpcType_BaoZang = NpcType(shared_proto.BaseTargetType_NpcBaoZang)
	// 名城Npc
	NpcType_MingCheng = NpcType(shared_proto.BaseTargetType_NpcMingCheng)
	// 军团怪Npc
	NpcType_JunTuan = NpcType(shared_proto.BaseTargetType_NpcJunTuan)
)

func NewGuildWorkshopId(guildId int64) int64 {
	return GetNpcId(uint64(guildId), 0, NpcType_Guild)
}

func GetWorkshopGuildId(id int64) int64 {
	return int64(GetNpcIdSequence(id))
}

func NewHomeNpcId(dataId uint64) int64 {
	return GetNpcId(0, dataId, NpcType_HomeNpc)
}

func NewNpcMemberId(guildId int64, sequence uint64) int64 {
	return GetNpcId(sequence, uint64(guildId), NpcType_Guild)
}

func NewMonsterId(monsterId uint64) int64 {
	return GetNpcId(monsterId, 0, NpcType_Monster)
}

func NewBaoZangNpcId(block, index, dataId uint64) int64 {
	sequence := block<<8 | index
	return GetNpcId(sequence, dataId, NpcType_BaoZang)
}

func GetBaoZangBlock(npcId int64) uint64 {
	sequence := GetNpcIdSequence(npcId)
	return sequence >> 8
}

func GetBaoZangIndex(npcId int64) uint64 {
	sequence := GetNpcIdSequence(npcId)
	return sequence & math.MaxUint8
}

func NewJunTuanNpcId(block, index, dataId uint64) int64 {
	sequence := block<<8 | index
	return GetNpcId(sequence, dataId, NpcType_JunTuan)
}

func GetJunTuanBlock(npcId int64) uint64 {
	sequence := GetNpcIdSequence(npcId)
	return sequence >> 8
}

func GetJunTuanIndex(npcId int64) uint64 {
	sequence := GetNpcIdSequence(npcId)
	return sequence & math.MaxUint8
}

func NewXiongNuNpcId(guildId int64, level uint64) int64 {
	return GetNpcId(uint64(guildId), level, NpcType_XiongNu)
}

func GetXiongNuGuildId(npcId int64) int64 {
	return int64(GetNpcIdSequence(npcId))
}

func GetNpcId(sequence uint64, dataId uint64, t NpcType) int64 {
	uid := sequence<<30 | dataId<<5 | uint64(t)
	return toNpcId(uid)
}

func GetNpcIdSequence(npcId int64) uint64 {
	uid := toUid(npcId)
	sequence := uid >> 30
	return sequence
}

func GetNpcDataId(npcId int64) uint64 {
	uid := toUid(npcId)
	sequence := (uid >> 5) & NpcDataMask
	return sequence
}

func GetNpcIdType(npcId int64) NpcType {
	uid := toUid(npcId)
	t := NpcType(uid & 31)
	return t
}

const TroopMaxSequence = math.MaxUint8

// troop id
func NewNpcTroopId(npcId int64, sequence uint64) int64 {
	uid := toUid(npcId)
	return toNpcId(uid<<8 | sequence)
}

func GetTroopNpcId(troopId int64) int64 {
	uid := toUid(troopId)
	return toNpcId(uid >> 8)
}

func GetTroopSequence(troopId int64) uint64 {
	uid := toUid(troopId)
	return uid & math.MaxUint8
}
