package goods

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/icon"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/check"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/util/weight"
)

var TypeArray = []shared_proto.EquipmentType{
	shared_proto.EquipmentType_WU_QI,
	shared_proto.EquipmentType_KAI_JIA,
	shared_proto.EquipmentType_HU_TUI,
	shared_proto.EquipmentType_SHI_PIN,
	shared_proto.EquipmentType_TOU_KUI,
}

//gogen:config
type EquipmentData struct {
	_  struct{} `file:"物品/装备.txt"`
	_  struct{} `proto:"shared_proto.EquipmentDataProto"`
	_  struct{} `protoconfig:"equipment"`
	Id uint64

	Name string

	Desc string

	Icon *icon.Icon `protofield:"IconId,%s.Id"` // 图标

	// 部件
	Type shared_proto.EquipmentType `type:"enum"`

	// 品质
	Quality *EquipmentQualityData `protofield:",config.U64ToI32(%s.Id)"`

	BaseStat      *data.SpriteStat
	BaseStatProto *shared_proto.SpriteStatProto `head:"-" protofield:"-"`
}

func (d *EquipmentData) Init(filename string) {
	check.PanicNotTrue(d.BaseStat.Sum4D() > 0, "装备配置%v 配置的属性至少需要有一项 > 0, %s, coef: %v", d.Id, filename)
	d.BaseStatProto = d.BaseStat.Encode4Init()
}

func (d *EquipmentData) DataId() uint64 {
	return d.Id
}

func (d *EquipmentData) GoodsType() GoodsType {
	return EQUIPMENT
}

func (e *EquipmentData) CalculateTotalStat(levelData *EquipmentQualityLevelData, refinedData *EquipmentRefinedData) (levelStat, refinedStat *shared_proto.SpriteStatProto, totalStat *data.SpriteStat) {

	act := e.BaseStat.Attack
	def := e.BaseStat.Defense
	str := e.BaseStat.Strength
	dex := e.BaseStat.Dexterity

	if levelData.LevelStat > 0 {
		levelStat = data.New4DStatProto(data.CalRate4DStat(levelData.LevelStat, act, def, str, dex))
	}

	if refinedData != nil {
		refinedStat = data.New4DStatProto(
			u64.MultiCoef(act, refinedData.StatCoef),
			u64.MultiCoef(def, refinedData.StatCoef),
			u64.MultiCoef(str, refinedData.StatCoef),
			u64.MultiCoef(dex, refinedData.StatCoef),
		)
	}

	b := data.NewSpriteStatBuilder()
	b.AddProto(e.BaseStatProto)
	b.AddProto(levelStat) // 里面有处理nil的情况
	b.AddProto(refinedStat)

	totalStat = b.Build()
	return
}

func (e *EquipmentData) CalculateUpgradeLevelStat(currentLevelData, nextLevelData *EquipmentQualityLevelData) *shared_proto.SpriteStatProto {
	if nextLevelData.LevelStat <= currentLevelData.LevelStat {
		return nil
	}

	cact := e.BaseStat.Attack
	cdef := e.BaseStat.Defense
	cstr := e.BaseStat.Strength
	cdex := e.BaseStat.Dexterity
	cact, cdef, cstr, cdex = data.CalRate4DStat(currentLevelData.LevelStat, cact, cdef, cstr, cdex)

	nact := e.BaseStat.Attack
	ndef := e.BaseStat.Defense
	nstr := e.BaseStat.Strength
	ndex := e.BaseStat.Dexterity
	nact, ndef, nstr, ndex = data.CalRate4DStat(nextLevelData.LevelStat, nact, ndef, nstr, ndex)

	return data.New4DStatProto(u64.Sub(nact, cact), u64.Sub(ndef, cdef), u64.Sub(nstr, cstr), u64.Sub(ndex, cdex))
}

func (e *EquipmentData) CalculateRefinedStat(refinedData *EquipmentRefinedData) *shared_proto.SpriteStatProto {

	if refinedData != nil {
		act := e.BaseStat.Attack
		def := e.BaseStat.Defense
		str := e.BaseStat.Strength
		dex := e.BaseStat.Dexterity

		return data.New4DStatProto(
			u64.MultiCoef(act, refinedData.CurrentStatCoef),
			u64.MultiCoef(def, refinedData.CurrentStatCoef),
			u64.MultiCoef(str, refinedData.CurrentStatCoef),
			u64.MultiCoef(dex, refinedData.CurrentStatCoef),
		)
	}

	return nil
}

// 装备品质配置
//gogen:config
type EquipmentQualityData struct {
	_ struct{} `file:"物品/装备品质.txt"`
	_ struct{} `proto:"shared_proto.EquipmentQualityProto"`
	_ struct{} `protoconfig:"equipment_quality"`

	Id           uint64
	Level        uint64        `validator:"int>0" default:"1"` // 作废，用 GoodsQuality
	GoodsQuality *GoodsQuality `protofield:",config.U64ToI32(%s.Level)"`

	// 最高强化次数
	FirstLevelRefined *EquipmentRefinedData `default:"1" protofield:"-"`
	RefinedLevelLimit uint64

	// 装备熔炼返还物品
	SmeltBackCount uint64 `validator:"int>0"`

	LevelCostCoef float64 `protofield:"-"`
	LevelStatCoef float64 `protofield:"-"`

	LevelDatas []*EquipmentQualityLevelData `head:"-" protofield:"-"`
}

func (d *EquipmentQualityData) MustLevel(level uint64) *EquipmentQualityLevelData {
	if level > 0 && int(level) <= len(d.LevelDatas) {
		return d.LevelDatas[level-1]
	}

	if level <= 0 {
		return d.LevelDatas[0]
	}

	return d.LevelDatas[len(d.LevelDatas)-1]
}

func (d *EquipmentQualityData) calculateLevelStat(level uint64) uint64 {
	// 升级带来的总属性=装备品质*20*（装备等级-1）

	lv := u64.Sub(level, 1)
	if lv > 0 {
		return u64.MultiCoef(lv, d.LevelStatCoef)
	}

	return 0
}

func (d *EquipmentQualityData) Init(filename string, dataMap map[uint64]*EquipmentQualityData, config interface {
	GetEquipmentLevelDataArray() []*EquipmentLevelData
	GetHeroLevelSubDataArray() []*data.HeroLevelSubData
}) {

	check.PanicNotTrue(d.SmeltBackCount > 0, "装备品质配置%v 熔炼返还物品个数必须 > 0, %s, coef: %v", d.Id, filename, d.SmeltBackCount)

	levelMap := make(map[uint64]*EquipmentQualityLevelData, len(config.GetEquipmentLevelDataArray()))
	for _, v := range config.GetEquipmentLevelDataArray() {
		e := &EquipmentQualityLevelData{}
		e.Level = v.Level
		e.UpgradeLevelCost = u64.MultiCoef(1, d.LevelCostCoef*v.UpgradeCostCoef)
		e.LevelStat = d.calculateLevelStat(v.Level)

		levelMap[v.Level] = e
	}

	d.LevelDatas = make([]*EquipmentQualityLevelData, len(levelMap))
	for i := 0; i < len(levelMap); i++ {
		lv := u64.FromInt(i + 1)
		current := levelMap[lv]
		d.LevelDatas[i] = current
		check.PanicNotTrue(current != nil, "%s 装备品质配置%v 关联的%v级等级数据没找到", filename, d.Id, lv)

		if lv > 1 {
			prev := levelMap[lv-1]
			prev.nextLevel = current

			check.PanicNotTrue(current.UpgradeLevelCost > 0, "%s 装备品质配置[%v] 计算出来的%v级消耗个数必须>0", filename, d.Id, lv)
			check.PanicNotTrue(current.LevelStat > 0, "%s 装备品质配置%v 计算出来的%v级升级属性必须>0", filename, d.Id, lv)

			check.PanicNotTrue(prev.UpgradeLevelCost < current.UpgradeLevelCost, "%s 装备品质配置%v 计算出来的%v级升级消耗个数必须>上一等级", filename, d.Id, lv)
			check.PanicNotTrue(prev.LevelStat < current.LevelStat, "%s 装备品质配置%v 计算出来的%v级升级属性必须>上一等级", filename, d.Id, lv)

			current.CurrentUpgradeLevelCost = current.UpgradeLevelCost - prev.UpgradeLevelCost
			current.CurrentLevelStat = current.LevelStat - prev.LevelStat
		} else {
			current.CurrentUpgradeLevelCost = current.UpgradeLevelCost
			current.CurrentLevelStat = current.LevelStat
		}
	}

	heroMaxLevelData := config.GetHeroLevelSubDataArray()[len(config.GetHeroLevelSubDataArray())-1]
	check.PanicNotTrue(u64.FromInt(len(d.LevelDatas)) >= heroMaxLevelData.EquipmentLevelLimit, "%s 装备品质配置%v 关联的等级个数不足，当前版本装备等级上限是%v", filename, d.Id, heroMaxLevelData.EquipmentLevelLimit)

}

//gogen:config
type EquipmentLevelData struct {
	_     struct{} `file:"物品/装备等级.txt"`
	Level uint64

	// 等级消耗系数
	UpgradeCostCoef float64

	nextLevel *EquipmentLevelData `head:"-" protofield:"-"`
}

func (d *EquipmentLevelData) Init(filename string, dataMap map[uint64]*EquipmentLevelData) {

	if d.Level > 1 {
		prevLevel := dataMap[d.Level-1]
		check.PanicNotTrue(prevLevel != nil, "%s 装备等级，没有找到[%v]级的配置", filename, d.Level-1)

		check.PanicNotTrue(prevLevel.UpgradeCostCoef <= d.UpgradeCostCoef, "%s 装备等级配置%v级 配置的升级消耗比前一级的要少", filename, d.Level)

		prevLevel.nextLevel = d
	}
}

type EquipmentQualityLevelData struct {
	Level uint64

	UpgradeLevelCost        uint64
	CurrentUpgradeLevelCost uint64

	LevelStat        uint64
	CurrentLevelStat uint64

	nextLevel *EquipmentQualityLevelData
}

func (data *EquipmentQualityLevelData) NextLevel() *EquipmentQualityLevelData {
	return data.nextLevel
}

// 装备强化配置
//gogen:config
type EquipmentRefinedData struct {
	_     struct{} `file:"物品/装备强化.txt"`
	_     struct{} `proto:"shared_proto.EquipmentRefinedProto"`
	_     struct{} `protoconfig:"equipment_refined"`
	Level uint64   `validator:"int>0"`

	StatCoef        float64 `protofield:"-"`
	CurrentStatCoef float64 `head:"-" protofield:"-"`

	CostCount uint64 `validator:"int>0"`

	HeroLevelLimit uint64 // 君主等级限制

	TotalCostCount uint64 `head:"-"`

	nextLevel *EquipmentRefinedData `head:"-" protofield:"-"`

	Star uint64 `head:"-,%s.Level / 2" protofield:"-"`
}

func (data *EquipmentRefinedData) NextLevel() *EquipmentRefinedData {
	return data.nextLevel
}

func (d *EquipmentRefinedData) Init(filename string, dataMap map[uint64]*EquipmentRefinedData) {

	if d.Level > 1 {
		prevLevel := dataMap[d.Level-1]
		check.PanicNotTrue(prevLevel != nil, "装备强化等级配置%v 等级必须从1开始连续配置, %s", d.Level, filename)

		d.TotalCostCount = prevLevel.TotalCostCount + d.CostCount

		prevLevel.nextLevel = d

		check.PanicNotTrue(prevLevel.StatCoef < d.StatCoef, "装备强化等级配置%v 属性系数必须 > 上一等级, %s, coef: %v prev: %v", d.Level, filename, d.StatCoef, prevLevel.StatCoef)

		d.CurrentStatCoef = d.StatCoef - prevLevel.StatCoef
	} else {
		d.TotalCostCount = d.CostCount
		d.CurrentStatCoef = d.StatCoef
	}

	check.PanicNotTrue(d.StatCoef > 0, "装备强化等级配置%v 属性系数必须 > 0, %s, coef: %v", d.Level, filename, d.StatCoef)
}

//gogen:config
type EquipmentTaozData struct {
	_ struct{} `file:"物品/装备套装.txt"`
	_ struct{} `proto:"shared_proto.EquipmentTaozProto"`
	_ struct{} `protoconfig:"equipment_taoz"`

	Level uint64

	// 装备要求件数
	Count uint64 `validator:"int>0"`

	// 装备要求强化等级
	RefinedLevel uint64 `validator:"int>0"`

	// 士气
	Morale uint64 `validator:"int>0"`

	SpriteStat *data.SpriteStat

	// 星级
	Star uint64 `head:"-,%s.Level / 2" protofield:"-"`
}

func (*EquipmentTaozData) InitAll(filename string, array []*EquipmentTaozData, configs interface {
	EquipmentTaozConfig() *EquipmentTaozConfig
}) {
	for i := 0; i < len(array); i++ {
		tz := array[i]
		check.PanicNotTrue(tz.Level == u64.FromInt(i+1), "%s 装备套装配置没有找到[%v]级数据，等级必须从1开始连续配置", filename, i+1)

		if i > 0 {
			prev := array[i-1]
			check.PanicNotTrue(prev.Count < tz.Count || prev.RefinedLevel < tz.RefinedLevel, "%s 装备套装配置[%v]级数据，强化等级和件数不能小于上一级的数据", filename, tz.Level)
			check.PanicNotTrue(prev.Morale < tz.Morale, "%s 装备套装配置[%v]级数据，士气不能小于上一级的数据", filename, tz.Level)
		}
	}

	config := configs.EquipmentTaozConfig()
	config.countMap = make(map[uint64]*EquipmentTaozCountData)

	for _, tz := range array {
		countData := config.countMap[tz.Count]
		if countData == nil {
			countData = &EquipmentTaozCountData{}
			countData.Count = tz.Count
			countData.MinRefinedLevel = tz.RefinedLevel

			config.countMap[tz.Count] = countData
		}

		countData.levelDatas = append(countData.levelDatas, tz)
		countData.MinRefinedLevel = u64.Min(countData.MinRefinedLevel, tz.RefinedLevel)
	}

	// 构建rankRandomer
	for _, v := range config.countMap {
		levelArray := make([]uint64, len(v.levelDatas))
		for i, d := range v.levelDatas {
			levelArray[i] = d.RefinedLevel
		}

		r, err := weight.NewRankRandomer(levelArray)
		if err != nil {
			logrus.Panicf("%s 装备套装排序出错", filename)
		}
		v.levelRandomer = r
	}
}

//gogen:config
type EquipmentTaozConfig struct {
	_ struct{} `singleton:"true"`
	_ struct{} `file:"物品/taoz.txt"`

	countMap map[uint64]*EquipmentTaozCountData `head:"-" protofield:"-"`
}

func (c *EquipmentTaozConfig) GetTaoz(count, level uint64) *EquipmentTaozData {
	countData := c.countMap[count]
	if countData == nil {
		return nil
	}

	return countData.GetTaoz(level)
}

type EquipmentTaozCountData struct {
	Count uint64

	levelDatas []*EquipmentTaozData

	MinRefinedLevel uint64

	levelRandomer *weight.WeightRandomer
}

func (d *EquipmentTaozCountData) GetTaoz(level uint64) *EquipmentTaozData {

	if level < d.MinRefinedLevel {
		// 比最低要求还低
		return nil
	}

	return d.levelDatas[d.levelRandomer.Index(level)]
}
