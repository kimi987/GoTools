package data

import (
	"github.com/lightpaw/male7/util/u64"
	. "github.com/onsi/gomega"
	"testing"
)

func TestFightAmount(t *testing.T) {
	RegisterTestingT(t)

	stat1 := &SpriteStat{
		Attack: 1000,
	}

	Ω(stat1.Sum4D()).Should(Equal(uint64(1000)))

	Ω(stat1.FightAmount(0, 0)).Should(Equal(uint64(0)))
	Ω(stat1.FightAmount(1, 0)).Should(Equal(uint64(1000)))
	Ω(stat1.FightAmount(8, 0)).Should(Equal(uint64(2000)))

	stat2 := &SpriteStat{
		Attack:         250,
		Defense:        250,
		Strength:       250,
		Dexterity:      250,
		DamageIncrePer: 10000,
		DamageDecrePer: 10000,
	}
	Ω(stat2.Sum4D()).Should(Equal(uint64(1000)))

	Ω(stat2.FightAmount(0, 0)).Should(Equal(uint64(0)))
	Ω(stat2.FightAmount(4, 0)).Should(Equal(uint64(2000)))

	testProtoFightAmount(stat1)
	testProtoFightAmount(stat2)
}

func testProtoFightAmount(stat *SpriteStat) {
	for _, soldier := range []uint64{0, 10, 100, 1000, 100000} {
		Ω(u64.Int32(stat.FightAmount(soldier, 0))).Should(Equal(ProtoFightAmount(stat.Encode(), u64.Int32(soldier), 0)))
		Ω(u64.Int32(stat.FightAmount(soldier, 100))).Should(Equal(ProtoFightAmount(stat.Encode(), u64.Int32(soldier), 100)))
		Ω(u64.Int32(stat.FightAmount(soldier, 10000))).Should(Equal(ProtoFightAmount(stat.Encode(), u64.Int32(soldier), 10000)))
		Ω(u64.Int32(stat.FightAmount(soldier, 100000))).Should(Equal(ProtoFightAmount(stat.Encode(), u64.Int32(soldier), 100000)))
	}

	for i := 0; i < 100; i++ {
		soldier := uint64(i)
		Ω(u64.Int32(stat.FightAmount(soldier, 0))).Should(Equal(ProtoFightAmount(stat.Encode(), u64.Int32(soldier), 0)))
		Ω(u64.Int32(stat.FightAmount(soldier, 33333))).Should(Equal(ProtoFightAmount(stat.Encode(), u64.Int32(soldier), 33333)))
	}
}

func TestTroopFightAmount(t *testing.T) {
	RegisterTestingT(t)

	Ω(TroopFightAmount()).Should(Equal(uint64(0)))
	Ω(TroopFightAmount()).Should(Equal(uint64(0)))
	Ω(TroopFightAmount(100)).Should(Equal(uint64(100)))
	Ω(TroopFightAmount(1)).Should(Equal(uint64(1)))
	Ω(TroopFightAmount(1, 1)).Should(Equal(uint64(1)))
	Ω(TroopFightAmount(1, 1, 1)).Should(Equal(uint64(1)))
	Ω(TroopFightAmount(1, 1, 1, 1)).Should(Equal(uint64(1)))
	Ω(TroopFightAmount(1, 1, 1, 1, 1)).Should(Equal(uint64(1)))

	tfa := NewTroopFightAmount()
	Ω(tfa.ToU64()).Should(Equal(uint64(0)))

	tfa.Add(1)
	Ω(tfa.ToU64()).Should(Equal(uint64(1)))
	tfa.Add(1)
	Ω(tfa.ToU64()).Should(Equal(uint64(1)))
	tfa.Add(1)
	Ω(tfa.ToU64()).Should(Equal(uint64(1)))
	tfa.Add(1)
	Ω(tfa.ToU64()).Should(Equal(uint64(1)))
}

var benchStat = &SpriteStat{
	Attack:         250,
	Defense:        250,
	Strength:       250,
	Dexterity:      250,
	DamageIncrePer: 10000,
	DamageDecrePer: 10000,
}

func BenchmarkEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchStat.Encode()
	}
}

func BenchmarkEncodeInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchStat.Encode4Init()
	}
}

func TestCalRate4DStat(t *testing.T) {
	RegisterTestingT(t)

	newAct, newDef, newStr, newDex := CalRate4DStat(10, 0, 0, 0, 0)
	Ω(newAct).Should(Equal(uint64(0)))
	Ω(newDef).Should(Equal(uint64(0)))
	Ω(newStr).Should(Equal(uint64(0)))
	Ω(newDex).Should(Equal(uint64(0)))

	newAct, newDef, newStr, newDex = CalRate4DStat(10, 0, 0, 3, 0)
	Ω(newAct).Should(Equal(uint64(0)))
	Ω(newDef).Should(Equal(uint64(0)))
	Ω(newStr).Should(Equal(uint64(10)))
	Ω(newDex).Should(Equal(uint64(0)))

	newAct, newDef, newStr, newDex = CalRate4DStat(10, 0, 2, 3, 0)
	Ω(newAct).Should(Equal(uint64(0)))
	Ω(newDef).Should(Equal(uint64(4)))
	Ω(newStr).Should(Equal(uint64(6)))
	Ω(newDex).Should(Equal(uint64(0)))

	newAct, newDef, newStr, newDex = CalRate4DStat(10, 1, 2, 3, 0)
	Ω(newAct).Should(Equal(uint64(2)))
	Ω(newDef).Should(Equal(uint64(3)))
	Ω(newStr).Should(Equal(uint64(5)))
	Ω(newDex).Should(Equal(uint64(0)))

	newAct, newDef, newStr, newDex = CalRate4DStat(10, 1, 2, 3, 0)
	Ω(newAct).Should(Equal(uint64(2)))
	Ω(newDef).Should(Equal(uint64(3)))
	Ω(newStr).Should(Equal(uint64(5)))
	Ω(newDex).Should(Equal(uint64(0)))

}
