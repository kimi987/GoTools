package entity

import (
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/config/resdata"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
)

func newBuildingEffect() *building_effect {
	b := &building_effect{}
	b.soldierRaceExtraStat = make([]*data.SpriteStat, len(shared_proto.Race_name)-1)

	b.goldEffect = &ResEffect{resType: shared_proto.ResType_GOLD}
	b.foodEffect = &ResEffect{resType: shared_proto.ResType_FOOD}
	b.woodEffect = &ResEffect{resType: shared_proto.ResType_WOOD}
	b.stoneEffect = &ResEffect{resType: shared_proto.ResType_STONE}

	return b
}

type building_effect struct {
	prosperityCapcity uint64 // 繁荣度最大上限

	protectedCapcity uint64

	// 兵营
	soldierCapcity         uint64
	soldierOutput          uint64
	woundedCapcity         uint64
	newSoldierCapcity      uint64
	newSoldierOutput       uint64
	extraLoad              uint64
	forceSoldier           uint64
	newSoldierRecruitCount uint64 // 每次能够征兵的数量

	// 兵种
	soldierRaceExtraStat  []*data.SpriteStat
	allSoldierExtraStat   *data.SpriteStat
	farSoldierExtraStat   *data.SpriteStat
	closeSoldierExtraStat *data.SpriteStat

	// 训练馆
	trainCoef float64 // 修炼馆系数

	// 官府
	buildingWorkerCdr float64

	// 书院
	techWorkerCdr float64

	homeWallStat      *data.SpriteStat // 主城城墙属性
	homeWallFixDamage uint64           // 主城城墙固定伤害

	goldEffect  *ResEffect
	foodEffect  *ResEffect
	woodEffect  *ResEffect
	stoneEffect *ResEffect

	// 建筑升级消耗折扣
	buildingCostReduceCoef float64

	// 科技升级消耗折扣
	techCostReduceCoef float64
}

func (d *building_effect) GetBuildingCostReduceCoef() float64 {
	return d.buildingCostReduceCoef
}

func (d *building_effect) GetBuildingCost(cost *resdata.Cost) *resdata.Cost {
	return getCostWithReduceCoef(cost, d.buildingCostReduceCoef)
}

func (d *building_effect) GetTechCostReduceCoef() float64 {
	return d.techCostReduceCoef
}

func (d *building_effect) GetTechCost(cost *resdata.Cost) *resdata.Cost {
	return getCostWithReduceCoef(cost, d.techCostReduceCoef)
}

func getCostWithReduceCoef(cost *resdata.Cost, costReduceCoef float64) *resdata.Cost {
	if costReduceCoef > 0 && costReduceCoef < 1 {
		return cost.MultipleF64(costReduceCoef)
	}
	return cost
}

type ResEffect struct {
	resType shared_proto.ResType

	// 仓库容量
	capcity uint64

	// 产出
	output        *data.Amount
	outputCapcity *data.Amount

	// 农场
	farmExtraOutput *data.Amount

	// 税收
	tax uint64
}

func (d *building_effect) getResEffect(resType shared_proto.ResType) *ResEffect {
	switch resType {
	case shared_proto.ResType_GOLD:
		return d.goldEffect
	case shared_proto.ResType_FOOD:
		return d.foodEffect
	case shared_proto.ResType_WOOD:
		return d.woodEffect
	case shared_proto.ResType_STONE:
		return d.stoneEffect
	default:
		return nil
	}
}

func (d *building_effect) GetTrainCoef() float64 {
	return d.trainCoef
}

func (d *building_effect) setAllSoldierExtraStat(toSet *data.SpriteStat) {
	d.allSoldierExtraStat = toSet
}

func (d *building_effect) GetAllSoldierExtraStat() *data.SpriteStat {
	return d.allSoldierExtraStat
}

func (d *building_effect) setSoldierExtraStat(race shared_proto.Race, toSet *data.SpriteStat) {

	index := int(race) - 1
	if index >= 0 && index < len(d.soldierRaceExtraStat) {
		d.soldierRaceExtraStat[index] = toSet
	}
}

func (d *building_effect) GetSoldierExtraStat(race shared_proto.Race) *data.SpriteStat {

	index := int(race) - 1
	if index >= 0 && index < len(d.soldierRaceExtraStat) {
		return d.soldierRaceExtraStat[index]
	}

	return data.EmptyStat()
}

func (d *building_effect) setSoldierFightTypeExtraStat(isFar bool, toSet *data.SpriteStat) {
	if isFar {
		d.farSoldierExtraStat = toSet
	} else {
		d.closeSoldierExtraStat = toSet
	}
}

func (d *building_effect) GetSoldierFightTypeExtraStat(isFar bool) *data.SpriteStat {
	if isFar {
		return d.farSoldierExtraStat
	} else {
		return d.closeSoldierExtraStat
	}
}

func (d *building_effect) AddSoldierExtraStat(race shared_proto.Race, toAdd *data.SpriteStat) *data.SpriteStat {

	index := int(race) - 1
	if index >= 0 && index < len(d.soldierRaceExtraStat) {
		newStat := data.AppendSpriteStat(d.soldierRaceExtraStat[index], toAdd)
		d.soldierRaceExtraStat[index] = newStat
		return newStat
	}

	return data.EmptyStat()
}

func (d *building_effect) GoldCapcity() uint64 {
	return d.goldEffect.capcity
}

func (d *building_effect) FoodCapcity() uint64 {
	return d.foodEffect.capcity
}

func (d *building_effect) WoodCapcity() uint64 {
	return d.woodEffect.capcity
}

func (d *building_effect) StoneCapcity() uint64 {
	return d.stoneEffect.capcity
}

func (d *building_effect) GetCapcity(resType shared_proto.ResType) uint64 {
	if re := d.getResEffect(resType); re != nil {
		return re.capcity
	}
	return 0
}

func (d *building_effect) GetTax(resType shared_proto.ResType) uint64 {
	if re := d.getResEffect(resType); re != nil {
		return re.tax
	}
	return 0
}

func (d *building_effect) GetFarmExtraOutput(resType shared_proto.ResType) *data.Amount {
	if re := d.getResEffect(resType); re != nil {
		return re.farmExtraOutput
	}
	return nil
}

func (d *building_effect) GetExtraOutput(resType shared_proto.ResType) *data.Amount {
	if re := d.getResEffect(resType); re != nil {
		return re.output
	}
	return nil
}

func (d *building_effect) GetExtraOutputCapcity(resType shared_proto.ResType) *data.Amount {
	if re := d.getResEffect(resType); re != nil {
		return re.outputCapcity
	}
	return nil
}

func (d *building_effect) ProtectedCapcity() uint64 {
	return d.protectedCapcity
}

func (d *building_effect) ProsperityCapcity() uint64 {
	return d.prosperityCapcity
}

func (hero *Hero) UpdateProsperityCapcity() bool {
	newProsperityCapcity := hero.domestic.calculateProsperityCapcity()
	if newProsperityCapcity != hero.buildingEffect.prosperityCapcity {
		hero.buildingEffect.prosperityCapcity = newProsperityCapcity
		hero.home.prosperity = u64.Min(hero.home.prosperity, hero.buildingEffect.prosperityCapcity)
		return true
	}
	return false
}

func (d *building_effect) CalculateDomestic(hero *Hero) {
	d.prosperityCapcity = hero.domestic.calculateProsperityCapcity()

	d.calculateBuildingType(hero, shared_proto.BuildingType_CANG_KU)
	d.calculateBuildingType(hero, shared_proto.BuildingType_JUN_YING)
	d.calculateBuildingType(hero, shared_proto.BuildingType_SI_TU_FU)
	d.calculateBuildingType(hero, shared_proto.BuildingType_GUAN_FU)
	d.calculateBuildingType(hero, shared_proto.BuildingType_CHENG_QIANG)
	d.calculateBuildingType(hero, shared_proto.BuildingType_SHU_YUAN)
	d.calculateBuildingType(hero, shared_proto.BuildingType_XING_YING)

	d.calculateAllResExtraEffect(hero)

	d.CalculateWallStat(hero)
	d.CalculateWallFixDamage(hero)

	d.CalculateTax(hero)

	d.CalculateBuildingCostReduceCoef(hero)
	d.CalculateTechCostReduceCoef(hero)
}

func (d *building_effect) CalculateSoldierStat(hero *Hero) {
	// 以前没有更新士兵属性
	d.CalculateAllSoldierExtraStatEffect(hero)
	d.CalculateSoldierExtraStatEffect(hero, shared_proto.Race_BU)
	d.CalculateSoldierExtraStatEffect(hero, shared_proto.Race_QI)
	d.CalculateSoldierExtraStatEffect(hero, shared_proto.Race_GONG)
	d.CalculateSoldierExtraStatEffect(hero, shared_proto.Race_CHE)
	d.CalculateSoldierExtraStatEffect(hero, shared_proto.Race_XIE)
	d.CalculateSoldierFightTypeStatEffect(hero, true)
	d.CalculateSoldierFightTypeStatEffect(hero, false)
}

func (d *building_effect) calculateBuildingType(hero *Hero, buildingType shared_proto.BuildingType) bool {
	switch buildingType {
	case shared_proto.BuildingType_CANG_KU:
		d.CalculateCangKuEffect(hero)
		return true
	case shared_proto.BuildingType_JUN_YING:
		d.CalculateJunYingEffect(hero)
		return true
	case shared_proto.BuildingType_SI_TU_FU:
		d.CalculateSiTuFuEffect(hero)
		return true
	case shared_proto.BuildingType_GUAN_FU:
		d.CalculateGuanFuEffect(hero)
		return true
	case shared_proto.BuildingType_SHU_YUAN:
		d.CalculateShuYuanEffect(hero)
		return true
	case shared_proto.BuildingType_CHENG_QIANG:
		d.CalculateChengQiangEffect(hero)
		return true
	default:
		if resType, ok := domestic_data.GetBuildingResType(buildingType); ok {
			d.CalculateResExtraEffect(hero, resType)
			return true
		}
	}

	return false
}

func (d *building_effect) CalculateTax(hero *Hero) {
	d.goldEffect.tax = hero.calculateTax(shared_proto.ResType_GOLD)
	d.foodEffect.tax = hero.calculateTax(shared_proto.ResType_FOOD)
	d.woodEffect.tax = hero.calculateTax(shared_proto.ResType_WOOD)
	d.stoneEffect.tax = hero.calculateTax(shared_proto.ResType_STONE)
}

func (d *building_effect) CalculateCangKuEffect(hero *Hero) {
	d.goldEffect.capcity = hero.calculateCapcity(shared_proto.ResType_GOLD)
	d.foodEffect.capcity = hero.calculateCapcity(shared_proto.ResType_FOOD)
	d.woodEffect.capcity = hero.calculateCapcity(shared_proto.ResType_WOOD)
	d.stoneEffect.capcity = hero.calculateCapcity(shared_proto.ResType_STONE)

	d.protectedCapcity = hero.calculateProtectedCapcity()
}

func (d *building_effect) calculateAllResExtraEffect(hero *Hero) {

	d.CalculateResExtraEffect(hero, shared_proto.ResType_GOLD)
	d.CalculateResExtraEffect(hero, shared_proto.ResType_FOOD)
	d.CalculateResExtraEffect(hero, shared_proto.ResType_WOOD)
	d.CalculateResExtraEffect(hero, shared_proto.ResType_STONE)

}

func (d *building_effect) CalculateResExtraEffect(hero *Hero, resType shared_proto.ResType) {

	switch resType {
	case shared_proto.ResType_GOLD:
		d.goldEffect.output = hero.calculateExtraOutput(resType)
		d.goldEffect.outputCapcity = hero.calculateExtraOutputCapcity(resType)
		d.goldEffect.farmExtraOutput = hero.calculateFarmExtraOutput(resType)
	case shared_proto.ResType_FOOD:
		d.foodEffect.output = hero.calculateExtraOutput(resType)
		d.foodEffect.outputCapcity = hero.calculateExtraOutputCapcity(resType)
		d.foodEffect.farmExtraOutput = hero.calculateFarmExtraOutput(resType)
	case shared_proto.ResType_WOOD:
		d.woodEffect.output = hero.calculateExtraOutput(resType)
		d.woodEffect.outputCapcity = hero.calculateExtraOutputCapcity(resType)
		d.woodEffect.farmExtraOutput = hero.calculateFarmExtraOutput(resType)
	case shared_proto.ResType_STONE:
		d.stoneEffect.output = hero.calculateExtraOutput(resType)
		d.stoneEffect.outputCapcity = hero.calculateExtraOutputCapcity(resType)
		d.stoneEffect.farmExtraOutput = hero.calculateFarmExtraOutput(resType)
	}
}

func (d *building_effect) calculateAllSoldierExtraStatEffect(hero *Hero) {

	d.CalculateResExtraEffect(hero, shared_proto.ResType_GOLD)
	d.CalculateResExtraEffect(hero, shared_proto.ResType_FOOD)
	d.CalculateResExtraEffect(hero, shared_proto.ResType_WOOD)
	d.CalculateResExtraEffect(hero, shared_proto.ResType_STONE)

}

func (d *building_effect) CalculateAllSoldierExtraStatEffect(hero *Hero) {

	b := data.NewSpriteStatBuilder()
	hero.walkEffect(func(effect *sub.BuildingEffectData) {
		s := effect.AllSoldierStat
		if s != nil {
			b.Add(s)
		}
	})
	d.setAllSoldierExtraStat(b.Build())
}

func (d *building_effect) CalculateSoldierExtraStatEffect(hero *Hero, race shared_proto.Race) {

	b := data.NewSpriteStatBuilder()
	hero.walkEffect(func(effect *sub.BuildingEffectData) {
		s := effect.GetSoldierStat(race)
		if s != nil {
			b.Add(s)
		}
	})
	d.setSoldierExtraStat(race, b.Build())
}

func (d *building_effect) CalculateSoldierFightTypeStatEffect(hero *Hero, isFar bool) {

	b := data.NewSpriteStatBuilder()
	hero.walkEffect(func(effect *sub.BuildingEffectData) {
		s := effect.GetSoldierFightTypeStat(isFar)
		if s != nil {
			b.Add(s)
		}
	})
	d.setSoldierFightTypeExtraStat(isFar, b.Build())
}

func (d *building_effect) CalculateJunYingEffect(hero *Hero) {

	// military
	var soldierCapcity uint64
	var soldierOutput uint64
	var woundedCapcity uint64
	var newSoldierCapcity uint64
	var newSoldierOutput uint64
	var newSoldierRecruitCount uint64
	var extraLoad uint64
	var forceSoldier uint64

	hero.walkEffect(func(effect *sub.BuildingEffectData) {
		soldierCapcity += effect.SoldierCapcity
		soldierOutput += effect.SoldierOutput
		woundedCapcity += effect.WoundedSoldierCapcity
		newSoldierCapcity += effect.NewSoldierCapcity
		newSoldierOutput += effect.NewSoldierOutput
		newSoldierRecruitCount += effect.RecruitSoldierCount
		extraLoad += effect.SoldierLoad
		forceSoldier += effect.ForceSoldier
	})

	d.soldierCapcity = soldierCapcity
	d.soldierOutput = soldierOutput
	d.woundedCapcity = woundedCapcity
	d.newSoldierCapcity = newSoldierCapcity
	d.newSoldierOutput = newSoldierOutput
	d.newSoldierRecruitCount = newSoldierRecruitCount
	d.extraLoad = extraLoad
	d.forceSoldier = forceSoldier
}

func (d *building_effect) GetForceSoldier() uint64 {
	return d.forceSoldier
}

func (d *building_effect) HomeWallStat() *data.SpriteStat {
	if d.homeWallStat != nil {
		return d.homeWallStat
	}
	return data.EmptyStat()
}

func (d *building_effect) HomeWallFixDamage() uint64 {
	return d.homeWallFixDamage
}

func (d *building_effect) CalculateWallStat(hero *Hero) {
	d.calculateHomeWallStat(hero)
}

func (d *building_effect) calculateHomeWallStat(hero *Hero) {
	homeStatBuilder := data.NewSpriteStatBuilder()
	hero.walkEffect(func(effect *sub.BuildingEffectData) {
		if effect.HomeWallStat != nil {
			homeStatBuilder.Add(effect.HomeWallStat)
		}
	})

	d.homeWallStat = homeStatBuilder.Build()
}

func (d *building_effect) CalculateWallFixDamage(hero *Hero) {
	var homeFixDamage uint64
	home := hero.domestic.GetBuilding(shared_proto.BuildingType_CHENG_QIANG)
	if home != nil && home.Effect != nil {
		homeFixDamage += home.Effect.HomeWallFixDamage
	}

	for _, data := range hero.domestic.technologys {
		homeFixDamage += data.Effect.HomeWallFixDamage
	}

	d.homeWallFixDamage = homeFixDamage
}

func (d *building_effect) CalculateBuildingCostReduceCoef(hero *Hero) float64 {

	var coef float64 = 1
	hero.walkEffect(func(effect *sub.BuildingEffectData) {
		if effect.BuildingCostReduceCoef > 0 && effect.BuildingCostReduceCoef < 1 {
			coef *= 1 - effect.BuildingCostReduceCoef
		}
	})

	d.buildingCostReduceCoef = coef
	return coef
}

func (d *building_effect) CalculateTechCostReduceCoef(hero *Hero) float64 {

	var coef float64 = 1
	hero.walkEffect(func(effect *sub.BuildingEffectData) {
		if effect.TechCostReduceCoef > 0 && effect.TechCostReduceCoef < 1 {
			coef *= 1 - effect.TechCostReduceCoef
		}
	})

	d.techCostReduceCoef = coef
	return coef
}

// 司徒府
func (d *building_effect) CalculateSiTuFuEffect(hero *Hero) {
	var trainCoef float64

	hero.walkEffect(func(effect *sub.BuildingEffectData) {
		trainCoef += effect.TrainCoef
	})

	d.trainCoef = trainCoef
}

func (d *building_effect) CalculateGuanFuEffect(hero *Hero) {

	cdr := float64(0)

	hero.walkEffect(func(effect *sub.BuildingEffectData) {
		cdr += effect.BuildingWorkerCdr
	})

	d.buildingWorkerCdr = cdr
}

func (d *building_effect) CalculateShuYuanEffect(hero *Hero) {

	cdr := float64(0)

	hero.walkEffect(func(effect *sub.BuildingEffectData) {
		cdr += effect.TechWorkerCdr
	})

	d.techWorkerCdr = cdr
}

func (d *building_effect) CalculateChengQiangEffect(hero *Hero) {
	d.calculateHomeWallStat(hero)
}

func (d *Hero) calculateExtraOutput(resType shared_proto.ResType) *data.Amount {
	b := data.NewAmountBuilder()
	d.walkEffect(func(effect *sub.BuildingEffectData) {
		b.AddAmount(effect.GetOutput(resType))
	})
	return b.Amount()
}

func (d *Hero) calculateFarmExtraOutput(resType shared_proto.ResType) *data.Amount {
	b := data.NewAmountBuilder()
	d.walkEffect(func(effect *sub.BuildingEffectData) {
		b.AddAmount(effect.GetFarmOutput(resType))
	})
	return b.Amount()
}

func (d *Hero) calculateExtraOutputCapcity(resType shared_proto.ResType) *data.Amount {
	bc := data.NewAmountBuilder()
	d.walkEffect(func(effect *sub.BuildingEffectData) {
		bc.AddAmount(effect.GetOutputCapcity(resType))
	})
	return bc.Amount()
}

func (d *Hero) calculateTax(resType shared_proto.ResType) uint64 {
	b := data.NewAmountBuilder()

	d.walkEffect(func(effect *sub.BuildingEffectData) {
		b.AddAmount(effect.GetTax(resType))
	})

	return b.TotalAmount()
}

func (d *Hero) calculateCapcity(resType shared_proto.ResType) uint64 {
	b := data.NewAmountBuilder()

	d.walkEffect(func(effect *sub.BuildingEffectData) {
		b.AddAmount(effect.GetCapcity(resType))
	})

	return b.TotalAmount()
}

func (d *Hero) calculateProtectedCapcity() uint64 {
	b := data.NewAmountBuilder()

	d.walkEffect(func(effect *sub.BuildingEffectData) {
		b.AddAmount(effect.ProtectedCapcity)
	})

	return b.TotalAmount()
}

func (hero *Hero) walkEffect(f func(effect *sub.BuildingEffectData)) {

	// 建筑，科技
	hero.domestic.WalkMainAndOuterCityBuildingsAndTechnologysHasEffect(f)

	// 武将内政技能
	for _, c := range hero.military.captains {
		if c.starData != nil {
			effects := c.starData.GetBuildingEffectSpell(c.AbilityData().UnlockSpellCount)
			if len(effects) > 0 {
				for _, v := range effects {
					f(v)
				}
			}
		}
	}

	// 国家官职属性
	if hero.countryMisc.effect != nil {
		f(hero.countryMisc.effect)
	}
}

func (d *HeroDomestic) calculateProsperityCapcity() uint64 {
	var max uint64 = 0
	d.WalkMainAndOuterCityBuildings(func(buildingData *domestic_data.BuildingData) {
		max += buildingData.Prosperity
	})
	return max
}

func (c *building_effect) ExtraLoad() uint64 {
	return c.extraLoad
}

func (c *building_effect) WoundedCapcity() uint64 {
	return c.woundedCapcity
}

func (c *building_effect) NewSoldierCapcity() uint64 {
	return c.newSoldierCapcity
}

func (c *building_effect) SoldierCapcity() uint64 {
	return c.soldierCapcity
}

func (c *building_effect) SoldierOutput() uint64 {
	return c.soldierOutput
}
