package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/male7/config/resdata"
)

// 所有分城
func newOuterCities() *OuterCities {
	return &OuterCities{cityMap: map[uint64]*OuterCity{}}
}

type OuterCities struct {
	cityMap    map[uint64]*OuterCity
	unlockBit  int32
	autoUnlock bool
}

func (c *OuterCities) SetAutoUnlock() {
	c.autoUnlock = true
}

func (c *OuterCities) TryAutoUnlock() bool {
	if c.autoUnlock {
		c.autoUnlock = false
		return true
	}
	return false
}

func (cities *OuterCities) UnlockBit() int32 {
	return cities.unlockBit
}

func (cities *OuterCities) CityIds() (ids []int32) {
	ids = make([]int32, 0, len(cities.cityMap))
	for id := range cities.cityMap {
		ids = append(ids, u64.Int32(id))
	}
	return
}

func (cities *OuterCities) WalkCities(walkFunc func(city *OuterCity)) {
	for _, city := range cities.cityMap {
		walkFunc(city)
	}
}

func (cities *OuterCities) WalkBuildings(walkFunc func(layout *domestic_data.OuterCityLayoutData, building *domestic_data.OuterCityBuildingData)) {
	if len(cities.cityMap) <= 0 {
		return
	}

	for _, city := range cities.cityMap {
		city.WalkLayouts(walkFunc)
	}
}

func (cities *OuterCities) Unlock(city *OuterCity) {
	cities.cityMap[city.Data().Id] = city
}

func (cities *OuterCities) UpdateUnlockBit() {

	var unlockBit int32
	for _, city := range cities.cityMap {
		oct := city.outerCityType + 1

		offset := uint8(city.data.Diraction)
		if offset > 0 {
			offset--
		}
		offset *= 2

		// 每2个bit表示一个数据，0-没有分城 1-军事分城 2-经济分城
		unlockBit |= int32(oct) << offset
	}

	cities.unlockBit = unlockBit
}

func (cities *OuterCities) OuterCity(data *domestic_data.OuterCityData) *OuterCity {
	return cities.cityMap[data.Id]
}

func (cities *OuterCities) Encode() *shared_proto.OuterCitiesProto {
	proto := &shared_proto.OuterCitiesProto{Cities: make([]*shared_proto.OuterCityProto, 0, len(cities.cityMap))}

	cities.WalkCities(func(city *OuterCity) {
		proto.Cities = append(proto.Cities, city.Encode())
	})

	return proto
}

func (cities *OuterCities) unmarshal(proto *shared_proto.OuterCitiesProto, datas *config.ConfigDatas) {
	if proto == nil {
		return
	}

	for _, cityProto := range proto.Cities {
		cityData := datas.GetOuterCityData(u64.FromInt32(cityProto.CityDataId))
		city := NewOuterCity(cityData, u64.FromInt32(cityProto.CityType))
		cities.Unlock(city)
		city.unmarshal(cityProto, datas)
	}
	cities.UpdateUnlockBit()
}

func NewOuterCity(data *domestic_data.OuterCityData, cityType uint64) *OuterCity {
	outerCity := &OuterCity{
		data:          data,
		outerCityType: cityType,
		layouts:       make(map[uint64]*domestic_data.OuterCityLayoutData, len(data.FirstLevelLayoutDatas)),
	}

	for _, layoutData := range data.FirstLevelLayoutDatas {
		if !layoutData.DefaultUnlocked {
			// 默认不解锁
			continue
		}

		// 解锁
		outerCity.layouts[layoutData.Layout] = layoutData
	}

	return outerCity
}

// 外城
type OuterCity struct {
	data          *domestic_data.OuterCityData                  // 配置
	layouts       map[uint64]*domestic_data.OuterCityLayoutData // 所有分城的布局
	outerCityType uint64
}

func (outerCity *OuterCity) Type() uint64 {
	return outerCity.outerCityType
}

func (outerCity *OuterCity) SetType(toSet uint64) {
	outerCity.outerCityType = toSet
}

func (outerCity *OuterCity) Data() *domestic_data.OuterCityData {
	return outerCity.data
}

func (outerCity *OuterCity) WalkLayouts(walkFunc func(layout *domestic_data.OuterCityLayoutData, building *domestic_data.OuterCityBuildingData)) {
	if len(outerCity.layouts) <= 0 {
		return
	}

	for _, layout := range outerCity.layouts {
		walkFunc(layout, layout.GetBuilding(outerCity.outerCityType))
	}
}

func (outerCity *OuterCity) Count() uint64 {
	return uint64(len(outerCity.layouts))
}

func (outerCity *OuterCity) Layout(data *domestic_data.OuterCityLayoutData) *domestic_data.OuterCityLayoutData {
	return outerCity.layouts[data.Layout]
}

func (outerCity *OuterCity) SetLayout(data *domestic_data.OuterCityLayoutData) {
	outerCity.layouts[data.Layout] = data
}

func (outerCity *OuterCity) Encode() *shared_proto.OuterCityProto {
	proto := &shared_proto.OuterCityProto{}

	proto.CityDataId = u64.Int32(outerCity.data.Id)
	for _, layout := range outerCity.layouts {
		proto.LayoutIds = append(proto.LayoutIds, u64.Int32(layout.Id))
	}
	proto.CityType = u64.Int32(outerCity.outerCityType)

	return proto
}

func (outerCity *OuterCity) unmarshal(proto *shared_proto.OuterCityProto, datas *config.ConfigDatas) {

	for _, layoutId := range proto.LayoutIds {
		layoutData := datas.GetOuterCityLayoutData(u64.FromInt32(layoutId))
		if layoutData == nil {
			logrus.WithField("layoutId", layoutId).Errorf("玩家分城建筑丢失!")
			continue
		}
		outerCity.SetLayout(layoutData)
	}
}

// 获取改建耗资
func (outerCity *OuterCity) GetChangeTypeCost() *resdata.Cost {
	for _, v := range outerCity.layouts {
		if v.IsMain {
			return v.ChangeTypeCost
		}
	}
	return nil
}

// 改建
func (outerCity *OuterCity) Change() {
	for k, v := range outerCity.layouts {
		outerCity.layouts[k] = v.ChangeTypeLevel
	}
}
