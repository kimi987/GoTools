package entity

import (
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"time"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/config/domestic_data/sub"
)

func NewResourcePoint(layoutData *domestic_data.BuildingLayoutData, building *domestic_data.BuildingData, outputStartTime time.Time) *ResourcePoint {
	r := &ResourcePoint{}
	r.layoutData = layoutData
	r.building = building
	r.outputStartTime = outputStartTime

	return r
}

type ResourcePoint struct {
	layoutData *domestic_data.BuildingLayoutData

	building *domestic_data.BuildingData

	outputStartTime time.Time

	// 累积的产出资源量
	outputAmount uint64
}

func (r *ResourcePoint) encodeClient(hero *Hero, ctime time.Time) *shared_proto.HeroResourcePointProto {
	proto := &shared_proto.HeroResourcePointProto{}

	proto.LayoutId = u64.Int32(r.layoutData.Id)
	proto.BuildingId = u64.Int32(r.building.Id)

	resEffect := r.building.GetResPointEffect()
	if resEffect != nil {
		outputPerHour, outputCapcity := hero.CalculateResourcePointOutput(resEffect)
		proto.OutputPerHour = u64.Int32(outputPerHour)
		proto.OutputCapcity = u64.Int32(outputCapcity)

		// 重新计算一次产出量
		r.calculateCurrentOutput(outputPerHour, outputCapcity, hero, ctime)
	}

	proto.OutputAmount = u64.Int32(r.outputAmount)

	proto.Conflicted = hero.IsConflictResourcePoint(r.layoutData)
	proto.BaseLevelLocked = hero.BaseLevel() < r.layoutData.RequireBaseLevel

	return proto
}

func (r *ResourcePoint) encodeServer() *server_proto.HeroResourcePointServerProto {
	proto := &server_proto.HeroResourcePointServerProto{}

	proto.LayoutId = r.layoutData.Id
	proto.BuildingId = r.building.Id
	proto.OutputStartTime = timeutil.Marshal64(r.outputStartTime)

	return proto
}

func (r *ResourcePoint) OutputAmount() uint64 {
	return r.outputAmount
}

func (r *ResourcePoint) AddOutputAmount(toAdd uint64) uint64 {
	r.outputAmount += toAdd
	return r.outputAmount
}

func (r *ResourcePoint) ReduceOutputAmount(toReduce uint64) uint64 {
	r.outputAmount = u64.Sub(r.outputAmount, toReduce)
	return r.outputAmount
}

func (hero *Hero) IsConflictResourcePoint(layoutData *domestic_data.BuildingLayoutData) bool {
	return hero.IsHeroConflictResourcePoint(layoutData) || hero.IsNpcConflictResourcePoint(layoutData)
}

func (hero *Hero) IsHeroConflictResourcePoint(layoutData *domestic_data.BuildingLayoutData) bool {
	return u64.Contains(hero.domestic.conflictResourcePoints, layoutData.Id)
}

func (hero *Hero) IsNpcConflictResourcePoint(layoutData *domestic_data.BuildingLayoutData) bool {
	for _, v := range hero.homeNpcBaseMap {
		if v.data.EvenOffsetX == layoutData.RegionOffsetX && v.data.EvenOffsetY == layoutData.RegionOffsetY {
			return true
		}
	}

	return false
}

func (hero *Hero) GetNpcConflictResourcePointOffset() []cb.Cube {
	result := make([]cb.Cube, 0)
	for _, v := range hero.homeNpcBaseMap {
		offset := cb.XYCube(v.data.EvenOffsetX, v.data.EvenOffsetY)
		result = append(result, offset)
	}

	return result
}

func (r *ResourcePoint) CalculateCurrentOutput(hero *Hero, ctime time.Time) {

	resEffect := r.building.GetResPointEffect()
	if resEffect == nil {
		logrus.Error("英雄计算资源点产出，这个建筑不是资源点，building.Effect == nil")
		return
	}

	outputPerHour, outputCapcity := hero.CalculateResourcePointOutput(resEffect)
	r.calculateCurrentOutput(outputPerHour, outputCapcity, hero, ctime)

}

func (r *ResourcePoint) calculateCurrentOutput(outputPerHour, outputCapcity uint64, hero *Hero, ctime time.Time) {

	// 如果主城等级达不到，不加
	if hero.BaseLevel() < r.layoutData.RequireBaseLevel {
		return
	}

	// 当前点如果冲突，不加
	if hero.IsConflictResourcePoint(r.layoutData) {
		return
	}

	// 计算过了多少时间
	diff := ctime.Sub(r.outputStartTime)
	if diff <= 0 {
		return
	}

	// 看下这段时间可以加多少资源
	if r.outputAmount >= outputCapcity {
		return
	}

	toAdd := u64.Multi(outputPerHour, diff.Hours())
	if toAdd <= 0 {
		return
	}

	r.outputAmount += toAdd                                 // 加速本次获得的值
	r.outputAmount = u64.Min(r.outputAmount, outputCapcity) // 不能超出上限

	// 多扣了一些时间，不管
	r.outputStartTime = ctime

}

func (hero *Hero) CalculateResourcePointOutput(effect *sub.BuildingEffectData) (outputPerHour, outputCapcity uint64) {

	resType := effect.OutputType
	extraOutputPerHour := hero.Domestic().GetExtraOutput(resType)
	extraCapcity := hero.Domestic().GetExtraOutputCapcity(resType)

	return CalculateResourcePointOutputInfo(effect, extraOutputPerHour, extraCapcity)
}

func CalculateResourcePointOutputInfo(effect *sub.BuildingEffectData, extraOutputPerHour, extraCapcity *data.Amount) (outputPerHour, outputCapcity uint64) {

	outputPerHour = data.TotalAmount(effect.Output, extraOutputPerHour)
	outputCapcity = data.TotalAmount(effect.OutputCapcity, extraCapcity)
	return
}

func CalculateResourcePointFullTime(outputPerHour, outputCapcity uint64) (fullTime time.Duration) {
	if outputPerHour > 0 {
		fullTime = time.Duration(float64(time.Hour) * float64(outputCapcity) / float64(outputPerHour))
	}

	return
}

func (hero *Hero) TrySetHeroResourcePointConflicted(layoutData *domestic_data.BuildingLayoutData, conflicted bool, ctime time.Time) bool {
	if layoutData.IgnoreConflict || hero.IsHeroConflictResourcePoint(layoutData) == conflicted {
		return false
	}

	if conflicted {
		hero.SetHeroResourcePointConflicted(layoutData.Id, ctime)
	} else {
		hero.SetHeroResourcePointUnconflicted(layoutData.Id, ctime)
	}

	return true
}

func (hero *Hero) SetHeroResourcePointConflicted(layoutId uint64, ctime time.Time) {
	r := hero.domestic.resourcePoint[layoutId]

	if r != nil {
		r.CalculateCurrentOutput(hero, ctime)
	}

	// 设置layoutId 冲突了
	hero.domestic.conflictResourcePoints = u64.AddIfAbsent(hero.domestic.conflictResourcePoints, layoutId)
}

func (hero *Hero) SetHeroResourcePointUnconflicted(layoutId uint64, ctime time.Time) {
	r := hero.domestic.resourcePoint[layoutId]
	if r != nil {
		r.outputStartTime = ctime
	}

	// 删除layoutId 冲突

	hero.domestic.conflictResourcePoints = u64.RemoveIfPresent(hero.domestic.conflictResourcePoints, layoutId)

}

func (r *ResourcePoint) LayoutData() *domestic_data.BuildingLayoutData {
	return r.layoutData
}

func (r *ResourcePoint) Building() *domestic_data.BuildingData {
	return r.building
}

func (r *ResourcePoint) SetBuilding(toSet *domestic_data.BuildingData) {
	r.building = toSet
}

func (sr *ResourcePoint) OutputStartTime() time.Time {
	return sr.outputStartTime
}
