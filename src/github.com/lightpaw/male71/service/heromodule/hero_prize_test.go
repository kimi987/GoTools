package heromodule

import (
	"testing"
	. "github.com/onsi/gomega"
	"github.com/lightpaw/male7/config/confpath"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/entity"
	"time"
)

func TestAddBaowu(t *testing.T) {
	RegisterTestingT(t)

	datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	Ω(err).Should(Succeed())
	firstLevel := datas.BaowuData().MinKeyData

	ctime := time.Now()
	hero := entity.NewHero(1, "hero", datas.HeroInitData(), ctime)

	if firstLevel.UpgradeNeedCount > 0 {
		ctx := &HeroContext{}
		for i := uint64(1); i < firstLevel.UpgradeNeedCount; i++ {
			AddBaowu(ctx, hero, result, firstLevel, 1, ctime)

			Ω(hero.Depot().GetBaowuCount(firstLevel.Id)).Should(Equal(i))
		}

		if firstLevel.GetNextLevel() != nil {
			// 再加1个触发自动合成
			AddBaowu(ctx, hero, result, firstLevel, 1, ctime)
			Ω(hero.Depot().GetBaowuCount(firstLevel.Id)).Should(Equal(uint64(0)))
			Ω(hero.Depot().GetBaowuCount(firstLevel.GetNextLevel().Id)).Should(Equal(uint64(1)))
		}
	}

}
