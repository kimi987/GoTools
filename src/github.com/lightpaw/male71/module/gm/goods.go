package gm

import (
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/depot"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/util/u64"
	"strings"
	"strconv"
)

func (m *GmModule) newGoodsGmGroup() *gm_group {
	return &gm_group{
		tab: "物品",
		handler: []*gm_handler{
			newHeroStringHandler("加物品", "", m.addGoods),
			newHeroStringHandler("加装备", "", m.addEquip),
			newHeroStringHandler("加宝石", "", m.addGem),
			newHeroStringHandler("加宝石(没最高级)", "", m.addGemNoTopLevel),
			newHeroStringHandler("加宝物", "", m.addBaowu),
			newHeroStringHandler("清空物品", "", m.clearGoods),
			newHeroStringHandler("清空物品ID", "物品id 空格分隔", m.clearGoodsById),
			newHeroStringHandler("清空装备", "", m.clearEquips),
			newHeroStringHandler("清空宝石", "", m.clearGems),
			newHeroStringHandler("清空宝物", "", m.clearBaowu),
		},
	}
}

func (m *GmModule) addGoods(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	for _, data := range m.datas.GoodsData().Array {
		heromodule.AddGoods(m.hctx, hero, result, data, 10000)
	}
}

func (m *GmModule) addEquip(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	ctime := m.time.CurrentTime()
	for _, data := range m.datas.EquipmentData().Array {
		heromodule.AddEquipData(m.hctx, hero, result, data, 1, ctime)
	}
}

func (m *GmModule) addGem(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	heromodule.AddGemArrayGive1(m.hctx, hero, result, m.datas.GemData().Array, true)
}

func (m *GmModule) addGemNoTopLevel(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	for _, gemData := range m.datas.GemData().Array {
		if gemData.NextLevel != nil {
			heromodule.AddGem(m.hctx, hero, result, gemData, 1, true)
		}
	}
}

func (m *GmModule) addBaowu(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	ctime := m.time.CurrentTime()
	for _, data := range m.datas.BaowuData().Array {
		heromodule.AddBaowu(m.hctx, hero, result, data, 1, ctime)
	}
}

func (m *GmModule) clearGoods(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	for _, data := range m.datas.GoodsData().Array {
		if c := hero.Depot().GetGoodsCount(data.Id); c > 0 {
			heromodule.ReduceGoodsAnyway(m.hctx, hero, result, data, c)
		}
	}
}

func (m *GmModule) clearGoodsById(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {

	for _, v := range strings.Split(input, " ") {
		goodsId, _ := strconv.ParseUint(v, 10, 64)
		if data := m.datas.GetGoodsData(goodsId); data != nil {
			if c := hero.Depot().GetGoodsCount(goodsId); c > 0 {
				heromodule.ReduceGoodsAnyway(m.hctx, hero, result, data, c)
			}
		}
	}
}

func (m *GmModule) clearEquips(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	hero.Depot().WalkGenIdGoods(func(goods goods.GenIdGoods) {
		hero.Depot().RemoveGenIdGoods(goods.Id())
		result.Add(depot.NewS2cGoodsExpireTimeRemoveMsg([]int32{u64.Int32(goods.Id())}))
	})
}

func (m *GmModule) clearGems(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	for _, data := range m.datas.GemData().Array {
		if c := hero.Depot().GetGoodsCount(data.Id); c > 0 {
			heromodule.ReduceGemAnyway(m.hctx, hero, result, data, c)
		}
	}
}

func (m *GmModule) clearBaowu(input string, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
	for _, data := range m.datas.BaowuData().Array {
		if c := hero.Depot().GetGoodsCount(data.Id); c > 0 {
			heromodule.ReduceBaowuAnyway(hero, result, data, c)
		}
	}
}

func (m *GmModule) newSingleGoodsGmGroup() *gm_group {

	group := &gm_group{
		tab: "加物品",
	}

	for _, data := range m.datas.GoodsData().Array {
		toSet := data
		group.handler = append(group.handler, newHeroIntHandler(toSet.Name, "1", func(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
			m.addOneGoods(amount, hero, result, hc, toSet)
		}))
	}

	return group
}

func (m *GmModule) addOneGoods(count int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController, data *goods.GoodsData) {
	if count > 0 {
		heromodule.AddGoods(m.hctx, hero, result, data, uint64(count))
	}
}

func (m *GmModule) newSingleEquipmentGmGroup() *gm_group {

	group := &gm_group{
		tab: "加装备",
	}

	for _, data := range m.datas.EquipmentData().Array {
		toSet := data
		group.handler = append(group.handler, newHeroIntHandler(toSet.Name, "1", func(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
			m.addOneEquipment(amount, hero, result, hc, toSet)
		}))
	}

	return group
}

func (m *GmModule) addOneEquipment(count int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController, data *goods.EquipmentData) {
	if count > 0 {
		heromodule.AddEquipData(m.hctx, hero, result, data, uint64(count), m.time.CurrentTime())
	}
}

func (m *GmModule) newSingleGemGmGroup() *gm_group {

	group := &gm_group{
		tab: "加宝石",
	}

	for _, data := range m.datas.GemData().Array {
		toSet := data
		group.handler = append(group.handler, newHeroIntHandler(toSet.Name, "1", func(amount int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController) {
			m.addOneGem(amount, hero, result, hc, toSet)
		}))
	}

	return group
}

func (m *GmModule) addOneGem(count int64, hero *entity.Hero, result herolock.LockResult, hc iface.HeroController, data *goods.GemData) {
	if count > 0 {
		heromodule.AddGem(m.hctx, hero, result, data, uint64(count), true)
	}
}
