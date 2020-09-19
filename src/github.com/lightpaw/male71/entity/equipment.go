package entity

import (
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/goods"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/u64"
)

func NewEquipment(id uint64, data *goods.EquipmentData) *Equipment {
	e := &Equipment{}
	e.id = id
	e.data = data
	e.levelData = data.Quality.MustLevel(1)

	e.totalStat = data.BaseStatProto

	return e
}

type Equipment struct {
	id uint64

	// 装备数据
	data *goods.EquipmentData

	// 强化属性
	refinedData *goods.EquipmentRefinedData

	levelData *goods.EquipmentQualityLevelData

	totalStat   *shared_proto.SpriteStatProto
	levelStat   *shared_proto.SpriteStatProto
	refinedStat *shared_proto.SpriteStatProto

	upgradeCostCount uint64
}

func (e *Equipment) Id() uint64 {
	return e.id
}

func (e *Equipment) GoodsData() goods.GenIdGoodsData {
	return e.data
}

func (e *Equipment) Data() *goods.EquipmentData {
	return e.data
}

func (e *Equipment) Level() uint64 {
	return e.levelData.Level
}

func (e *Equipment) LevelData() *goods.EquipmentQualityLevelData {
	return e.levelData
}

func (e *Equipment) UpgradeLevel() {
	if e.levelData.NextLevel() != nil {
		e.levelData = e.levelData.NextLevel()
	}
}

func (e *Equipment) SetLevelData(toSet *goods.EquipmentQualityLevelData) {
	e.levelData = toSet
}

func (e *Equipment) RefinedLevel() uint64 {
	if e.refinedData != nil {
		return e.refinedData.Level
	}

	return 0
}

func (e *Equipment) RefinedStar() uint64 {
	if e.refinedData != nil {
		return e.refinedData.Star
	}

	return 0
}

func (e *Equipment) RefinedData() *goods.EquipmentRefinedData {
	return e.refinedData
}

func (e *Equipment) SetRefinedData(toSet *goods.EquipmentRefinedData) {
	e.refinedData = toSet
}

func (e *Equipment) Rebuild() (uint64, *goods.EquipmentRefinedData) {
	c := e.upgradeCostCount
	r := e.refinedData

	e.levelData = e.data.Quality.MustLevel(1)
	e.upgradeCostCount = 0
	e.refinedData = nil
	e.CalculateProperties()
	return c, r
}

func (e *Equipment) RebuildUpgrade() (uint64) {
	c := e.upgradeCostCount
	e.levelData = e.data.Quality.MustLevel(1)
	e.upgradeCostCount = 0
	e.CalculateProperties()
	return c
}

func (e *Equipment) RebuildRefine() (*goods.EquipmentRefinedData) {
	r := e.refinedData
	e.refinedData = nil
	e.CalculateProperties()
	return r
}

func (e *Equipment) AddUpgradeCostCount(toAdd uint64) {
	e.upgradeCostCount += toAdd
}

func (e *Equipment) CalculateProperties() {
	var totalStat *data.SpriteStat
	e.levelStat, e.refinedStat, totalStat = e.data.CalculateTotalStat(e.levelData, e.refinedData)
	e.totalStat = totalStat.Encode()
}

func (e *Equipment) TotalStat() *shared_proto.SpriteStatProto {
	return e.totalStat
}

func (e *Equipment) unmarshal(proto *server_proto.EquipmentServerProto, datas interface {
	EquipmentRefinedData() *config.EquipmentRefinedDataConfig
}) {

	e.levelData = e.data.Quality.MustLevel(proto.Level)

	if proto.RefinedLevel > 0 {
		e.refinedData = datas.EquipmentRefinedData().Must(proto.RefinedLevel)
	}

	e.upgradeCostCount = proto.UpgradeCostCount

	e.CalculateProperties()
}

func (e *Equipment) encodeServer() *server_proto.EquipmentServerProto {

	proto := &server_proto.EquipmentServerProto{}
	proto.Id = e.Id()
	proto.DataId = e.Data().Id
	proto.Level = e.Level()

	if e.refinedData != nil {
		proto.RefinedLevel = e.refinedData.Level
	}

	proto.UpgradeCostCount = e.upgradeCostCount

	return proto
}

func (e *Equipment) EncodeClient() *shared_proto.EquipmentProto {

	proto := &shared_proto.EquipmentProto{}
	proto.Id = u64.Int32(e.id)
	proto.DataId = u64.Int32(e.data.Id)
	proto.Level = u64.Int32(e.Level())

	if e.refinedData != nil {
		proto.RefinedLevel = u64.Int32(e.refinedData.Level)
	}

	proto.TotalStat = e.totalStat
	proto.LevelStat = e.levelStat
	proto.CurrentRefinedStat = e.refinedStat

	// 计算出来
	if e.levelData.NextLevel() != nil {
		proto.UpgradeLevelStat = e.data.CalculateUpgradeLevelStat(e.levelData, e.levelData.NextLevel())
		proto.UpgradeLevelCost = u64.Int32(e.levelData.NextLevel().CurrentUpgradeLevelCost)
		proto.UpgradeLevelTotalCost = u64.Int32(e.upgradeCostCount)
	}

	nextRefinedData := e.data.Quality.FirstLevelRefined
	if e.refinedData != nil {
		nextRefinedData = e.refinedData.NextLevel()
	}

	if nextRefinedData != nil {
		proto.RefinedStat = e.data.CalculateRefinedStat(nextRefinedData)
		proto.RefinedStatPercent = i32.MultiF64(100, nextRefinedData.CurrentStatCoef)
	}

	return proto
}
