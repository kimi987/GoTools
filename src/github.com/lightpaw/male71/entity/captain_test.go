package entity

import (
	"github.com/lightpaw/male7/config/captain"
	"testing"
)
import (
	. "github.com/onsi/gomega"
)

func TestUpgradeAbility(t *testing.T) {
	RegisterTestingT(t)

	n := 10
	dataMap := make(map[uint64]*captain.CaptainAbilityData, n)
	for i := 0; i < n; i++ {
		lv := uint64(i + 1)

		data := &captain.CaptainAbilityData{}
		data.Ability = lv
		data.UpgradeExp = 100
		data.MaxLevel = uint64(n)

		dataMap[lv] = data
	}

	for _, v := range dataMap {
		v.Init(dataMap)
	}

	// 每一级都是100经验
	d, exp := UpgradeAbility(dataMap[1], 99)
	Ω(d).Should(Equal(dataMap[1]))
	Ω(exp).Should(Equal(uint64(99)))

	// 升级了
	d, exp = UpgradeAbility(dataMap[1], 101)
	Ω(d).Should(Equal(dataMap[2]))
	Ω(exp).Should(Equal(uint64(1)))

	// 连升4级
	d, exp = UpgradeAbility(dataMap[3], 450)
	Ω(d).Should(Equal(dataMap[7]))
	Ω(exp).Should(Equal(uint64(50)))

	// 升满级
	d, exp = UpgradeAbility(dataMap[1], 1000)
	Ω(d).Should(Equal(dataMap[10]))
	Ω(exp).Should(Equal(uint64(0)))

	// 多的经验清0
	d, exp = UpgradeAbility(dataMap[1], 5000)
	Ω(d).Should(Equal(dataMap[10]))
	Ω(exp).Should(Equal(uint64(0)))
}
