package equipment

import (
	"testing"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/mock"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/config/data"
	. "github.com/onsi/gomega"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/util/u64"
	"time"
)

func TestInhertRefine(t *testing.T) {
	//RegisterTestingT(t)
	//
	//m := newMockModule()
	//ctime := m.timeService.CurrentTime()
	//
	//result := mock.LockResult
	//
	//// 星级
	//inhertRefine(m, ctime, result)
	//// 等级
	//inhertUpgarde(m, ctime, result)
}

func inhertUpgarde(m *EquipmentModule, ctime time.Time, result herolock.LockResult) {
	oldEquip, newEquip, hero, captain := buildObjs(m, ctime, 3, 1)
	oldEquip.SetLevelData(oldEquip.Data().Quality.MustLevel(3))
	oldEquip.AddUpgradeCostCount(oldEquip.LevelData().UpgradeLevelCost)
	newEquip.SetLevelData(newEquip.Data().Quality.MustLevel(1))
	newEquip.AddUpgradeCostCount(newEquip.LevelData().UpgradeLevelCost)

	m.inheritUpgrade(oldEquip, newEquip, hero, result, captain)
	Ω(newEquip.Level()).Should(Equal(u64.FromInt32(4)))

	gid := m.datas.GoodsConfig().EquipmentRefinedGoods.Id
	Ω(hero.Depot().GetGoodsCount(gid)).Should(Equal(uint64(0)))
}

func inhertRefine(m *EquipmentModule, ctime time.Time, result herolock.LockResult) {
	oldEquip, newEquip, hero, captain := buildObjs(m, ctime, 1, 3)
	oldEquip.SetRefinedData(m.datas.EquipmentRefinedData().Get(3))
	newEquip.SetRefinedData(m.datas.EquipmentRefinedData().Get(1))

	m.inheritRefine(oldEquip, newEquip, hero, result, captain)
	Ω(newEquip.RefinedLevel()).Should(Equal(u64.FromInt32(3)))

	gid := m.datas.GoodsConfig().EquipmentRefinedGoods.Id
	Ω(hero.Depot().GetGoodsCount(gid)).Should(Equal(uint64(1)))
}

func buildObjs(m *EquipmentModule, ctime time.Time, equip1RefinedLimit uint64, equip2RefinedLimit uint64) (equip1 *entity.Equipment, equip2 *entity.Equipment, hero *entity.Hero, captain *entity.Captain) {
	hero = entity.NewHero(1, "hero1", m.datas.HeroInitData(), ctime)
	mock.DefaultHero(hero)

	mock.SetHeroTroopCaptain(hero, m.datas)
	captain = hero.GetTroopByIndex(0).GetCaptain(0)

	equip1Data := &goods.EquipmentData{Id: 1, Quality: m.datas.EquipmentQualityData().Get(10), BaseStat: data.New4DStat(1, 1, 1, 1), BaseStatProto: &shared_proto.SpriteStatProto{}}
	equip2Data := &goods.EquipmentData{Id: 2, Quality: m.datas.EquipmentQualityData().Get(10), BaseStat: data.New4DStat(1, 1, 1, 1), BaseStatProto: &shared_proto.SpriteStatProto{}}

	equip1 = entity.NewEquipment(hero.Depot().NewId(), equip1Data)
	equip2 = entity.NewEquipment(hero.Depot().NewId(), equip2Data)

	return
}

func newMockModule() *EquipmentModule {
	//datas, err := config.LoadConfigDatas(confpath.GetConfigPath())
	//Ω(err).Should(Succeed())
	return NewEquipmentModule(mock.MockDep())
}
