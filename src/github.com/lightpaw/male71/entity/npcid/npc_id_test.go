package npcid

import (
	. "github.com/onsi/gomega"
	"math"
	"testing"
	"github.com/lightpaw/male7/pb/shared_proto"
	"strings"
)

func TestValidType(t *testing.T) {
	RegisterTestingT(t)

	for name, v := range shared_proto.BaseTargetType_value {
		if strings.HasPrefix(strings.ToLower(name), "npc") {
			Ω(v < MaxNpcType).Should(BeTrue())
		}
	}
}

func TestIsNpcId(t *testing.T) {
	RegisterTestingT(t)

	Ω(IsNpcId(toNpcId(0))).Should(BeTrue())
	Ω(IsNpcId(toNpcId(100))).Should(BeTrue())
	Ω(IsNpcId(toNpcId(math.MaxInt64))).Should(BeTrue())

	Ω(toNpcId(0)).Should(Equal(int64(math.MinInt64)))

	Ω(toUid(toNpcId(0))).Should(Equal(uint64(0)))
	Ω(toUid(toNpcId(1))).Should(Equal(uint64(1)))
	Ω(toUid(toNpcId(100))).Should(Equal(uint64(100)))

	//Ω(toNpcId(math.MaxInt64)).Should(Equal(-math.MaxInt64))

	gen1 := NewNpcIdGen(0, NpcType_HomeNpc)
	Ω(gen1.Next(1)).Should(Equal(GetNpcId(1, 1, NpcType_HomeNpc)))
	Ω(gen1.Next(2)).Should(Equal(GetNpcId(2, 2, NpcType_HomeNpc)))

	next := gen1.Next(3)
	Ω(next).Should(Equal(GetNpcId(3, 3, NpcType_HomeNpc)))
	Ω(gen1.Sequence()).Should(Equal(uint64(3)))

	gen2 := NewNpcIdGen(3, NpcType_HomeNpc)

	Ω(gen1.Next(1)).Should(Equal(GetNpcId(4, 1, NpcType_HomeNpc)))
	Ω(gen2.Next(1)).Should(Equal(GetNpcId(4, 1, NpcType_HomeNpc)))

	Ω(gen1.Next(1)).Should(Equal(GetNpcId(5, 1, NpcType_HomeNpc)))
	Ω(gen2.Next(1)).Should(Equal(GetNpcId(5, 1, NpcType_HomeNpc)))
}

func TestLocal(t *testing.T) {
	RegisterTestingT(t)

	id := GetNpcId(0, 1, NpcType_HomeNpc)
	Ω(GetNpcIdSequence(id)).Should(Equal(uint64(0)))
	Ω(GetNpcIdType(id)).Should(Equal(NpcType(NpcType_HomeNpc)))
	Ω(GetNpcDataId(id)).Should(Equal(uint64(1)))

	id = GetNpcId(100, 2, NpcType_Monster)
	Ω(GetNpcIdSequence(id)).Should(Equal(uint64(100)))
	Ω(GetNpcIdType(id)).Should(Equal(NpcType(NpcType_Monster)))
	Ω(GetNpcDataId(id)).Should(Equal(uint64(2)))

	id = GetNpcId(0, 1, NpcType_BaoZang)
	Ω(GetNpcIdSequence(id)).Should(Equal(uint64(0)))
	Ω(GetNpcIdType(id)).Should(Equal(NpcType(NpcType_BaoZang)))
	Ω(GetNpcDataId(id)).Should(Equal(uint64(1)))

	id = GetNpcId(100, 2, NpcType_BaoZang)
	Ω(GetNpcIdSequence(id)).Should(Equal(uint64(100)))
	Ω(GetNpcIdType(id)).Should(Equal(NpcType(NpcType_BaoZang)))
	Ω(GetNpcDataId(id)).Should(Equal(uint64(2)))

}

func TestBaozId(t *testing.T) {
	RegisterTestingT(t)

	for b := uint64(0); b < 100; b++ {
		for i := uint64(0); i < 256; i++ {
			for d := uint64(0); d < 10; d++ {
				id := NewBaoZangNpcId(b, i, d)
				Ω(IsNpcId(id)).Should(BeTrue())
				Ω(GetNpcDataId(id)).Should(Equal(d))
				Ω(IsBaoZangNpcId(id)).Should(BeTrue())
				Ω(GetBaoZangBlock(id)).Should(Equal(b))
				Ω(GetBaoZangIndex(id)).Should(Equal(i))
			}
		}
	}

}
