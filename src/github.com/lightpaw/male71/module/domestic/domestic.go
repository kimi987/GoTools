package domestic

import (
	"context"
	"fmt"
	"github.com/lightpaw/logrus"
	"github.com/lightpaw/male7/config/buffer"
	"github.com/lightpaw/male7/config/data"
	"github.com/lightpaw/male7/config/domestic_data"
	"github.com/lightpaw/male7/config/domestic_data/sub"
	"github.com/lightpaw/male7/constants"
	"github.com/lightpaw/male7/entity"
	"github.com/lightpaw/male7/entity/cb"
	"github.com/lightpaw/male7/gen/iface"
	"github.com/lightpaw/male7/gen/pb/domestic"
	"github.com/lightpaw/male7/gen/pb/guild"
	"github.com/lightpaw/male7/gen/pb/military"
	"github.com/lightpaw/male7/gen/pb/misc"
	"github.com/lightpaw/male7/gen/pb/region"
	"github.com/lightpaw/male7/pb/server_proto"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/male7/service/conflict/heroservice/herolock"
	"github.com/lightpaw/male7/service/heromodule"
	"github.com/lightpaw/male7/service/operate_type"
	"github.com/lightpaw/male7/util"
	"github.com/lightpaw/male7/util/ctxfunc"
	"github.com/lightpaw/male7/util/i32"
	"github.com/lightpaw/male7/util/idbytes"
	"github.com/lightpaw/male7/util/imath"
	"github.com/lightpaw/male7/util/lock"
	"github.com/lightpaw/male7/util/must"
	"github.com/lightpaw/male7/util/timeutil"
	"github.com/lightpaw/male7/util/u64"
	"github.com/lightpaw/pbutil"
	"strings"
	"time"
	"github.com/lightpaw/male7/gamelogs"
)

func NewDomesticModule(dep iface.ServiceDep, seasonService iface.SeasonService, dbServuce iface.DbService,
	realmService iface.RealmService, baiZhanService iface.BaiZhanService, tssClient iface.TssClient,
	buffService iface.BuffService, region iface.RegionModule) *DomesticModule {

	techFirstLevelMap := make(map[uint64]*domestic_data.TechnologyData)
	for _, d := range dep.Datas().GetTechnologyDataArray() {
		if d.Level == 1 {
			techFirstLevelMap[d.Group] = d
		}
	}

	m := &DomesticModule{
		dep:                 dep,
		datas:               dep.Datas(),
		timeService:         dep.Time(),
		seasonService:       seasonService,
		dbServuce:           dbServuce,
		heroDataService:     dep.HeroData(),
		heroSnapshotService: dep.HeroSnapshot(),
		world:               dep.World(),
		realmService:        realmService,
		guildService:        dep.Guild(),
		baiZhanService:      baiZhanService,
		tssClient:           tssClient,
		techFirstLevelMap:   techFirstLevelMap,
		buffService:         buffService,
		region:              region,
	}
	return m
}

//gogen:iface
type DomesticModule struct {
	dep                 iface.ServiceDep
	datas               iface.ConfigDatas
	timeService         iface.TimeService
	seasonService       iface.SeasonService
	dbServuce           iface.DbService
	heroDataService     iface.HeroDataService
	heroSnapshotService iface.HeroSnapshotService
	world               iface.WorldService
	realmService        iface.RealmService
	guildService        iface.GuildService
	baiZhanService      iface.BaiZhanService
	tssClient           iface.TssClient
	buffService         iface.BuffService
	region              iface.RegionModule

	techFirstLevelMap map[uint64]*domestic_data.TechnologyData
}

//gogen:iface
func (m *DomesticModule) ProcessCreateResourceBuilding(proto *domestic.C2SCreateBuildingProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticCreateResourceBuilding)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		layoutId := u64.FromInt32(proto.GetId())

		// 位置无效
		layoutData := m.datas.GetBuildingLayoutData(layoutId)
		if layoutData == nil {
			logrus.Debugf("新建建筑，布局id无效, %s", layoutId)
			result.Add(domestic.ERR_CREATE_BUILDING_FAIL_INVALID_LAYOUT)
			return
		}

		buildingType := shared_proto.BuildingType(proto.GetType())
		heroDomestic := hero.Domestic()

		// 已经有建筑
		resourcePoint := heroDomestic.GetLayoutRes(layoutId)
		if resourcePoint != nil {
			logrus.Debugf("新建建筑，建筑已经存在, %s", layoutId)
			result.Add(domestic.ERR_CREATE_BUILDING_FAIL_NOT_EMPTY)
			return
		}

		ctime := m.timeService.CurrentTime()

		buildingId := domestic_data.BuildingId(buildingType, 1)
		building := m.datas.GetBuildingData(buildingId)
		if building == nil {
			// 为毛没有数据 ???
			logrus.Error("创建建筑没找到对应的建筑数据，%s-%s id: %s", proto.GetType(), 1, buildingId)
			result.Add(domestic.ERR_CREATE_BUILDING_FAIL_INVALID_TYPE)
			return
		}

		resEffect := building.GetResPointEffect()
		if resEffect == nil {
			logrus.Debugf("新建建筑，必须是资源点, %v", buildingType)
			result.Add(domestic.ERR_CREATE_BUILDING_FAIL_INVALID_TYPE)
			return
		}

		var workerPos int
		if building.WorkTime > 0 {
			// 建筑队休息
			workerPos = heroDomestic.GetFreeWorker(ctime)
			if workerPos < 0 {
				result.Add(domestic.ERR_CREATE_BUILDING_FAIL_WORKER_REST)
				return
			}
		}

		// 前提条件
		for _, requireBuilding := range building.RequireIds {
			heroBuilding := heroDomestic.GetBuilding(requireBuilding.Type)
			if heroBuilding == nil || heroBuilding.Level < requireBuilding.Level {
				logrus.Debugf("新建建筑，前置建筑还没有, %s", layoutId)
				result.Add(domestic.ERR_CREATE_BUILDING_FAIL_REQUIRE_NOT_REACH)
				return
			}
		}

		// 资源点不在你的势力范围，请先升级主城
		if hero.BaseLevel() < layoutData.RequireBaseLevel {
			logrus.Debugf("新建建筑，资源点不在你的势力范围")
			result.Add(domestic.ERR_CREATE_BUILDING_FAIL_RESOURCE_INVALID)
			return
		}

		// 资源点冲突
		if hero.IsConflictResourcePoint(layoutData) {
			logrus.Debugf("新建建筑，资源点冲突")
			result.Add(domestic.ERR_CREATE_BUILDING_FAIL_RESOURCE_CONFLICT)
			return
		}

		// 主城等级要求
		if building.BaseLevel != nil && hero.BaseLevel() < building.BaseLevel.Level {
			logrus.Debugf("新建建筑，主城等级不足")
			result.Add(domestic.ERR_CREATE_BUILDING_FAIL_REQUIRE_NOT_REACH)
			return
		}

		if !heromodule.TryReduceCostWithDianquanConvert(hctx, hero, result, hero.BuildingEffect().GetBuildingCost(building.Cost)) {
			logrus.Debugf("新建建筑，资源不足")
			result.Add(domestic.ERR_CREATE_BUILDING_FAIL_RESOURCE_NOT_ENOUGH)
			return
		}

		// 如果是资源点，添加领取时间
		hctx := heromodule.NewContext(m.dep, operate_type.DomesticCreateResourceBuilding)
		heromodule.AddExp(hctx, hero, result, building.HeroExp, ctime)
		result.Changed()

		var workerRestEndTime time.Time
		var seekHelp bool
		if building.WorkTime > 0 {
			workerRestEndTime, seekHelp = heroDomestic.AddWorkerRestEndTime(workerPos, ctime, building.CalculateWorkTime(heroDomestic.GetBuildingWorkerCoef(m.seasonService.Season().WorkerCdr)))
		} else {
			_, workerRestEndTime = heroDomestic.GetWorkerRestEndTime(workerPos)
		}
		resourcePoint = heroDomestic.SetResourcePoint(layoutData, building, ctime)

		// 计算产出
		result.Add(domestic.NewS2cCreateBuildingMsg(u64.Int32(layoutData.Id), u64.Int32(building.Id), imath.Int32(workerPos), timeutil.Marshal32(workerRestEndTime)))

		heromodule.SendResourcePointUpdateMsg(hero, result, resourcePoint)

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_RESOURCE_POINT_COUNT)

		if seekHelp {
			// 发消息通知，可以求助
			result.Add(guild.NewS2cUpdateSeekHelpMsg(constants.SeekTypeWorker, imath.Int32(workerPos), seekHelp))
		}

		result.Ok()
	})
}

//gogen:iface
func (m *DomesticModule) ProcessUpgradeResourceBuilding(proto *domestic.C2SUpgradeBuildingProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticUpgradeResourceBuilding)
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		layoutId := u64.FromInt32(proto.GetId())
		heroDomestic := hero.Domestic()

		// 还没建这个建筑
		resourcePoint := heroDomestic.GetLayoutRes(layoutId)
		if resourcePoint == nil {
			logrus.Debugf("升级建筑，这个建筑还未新建")
			result.Add(domestic.ERR_UPGRADE_BUILDING_FAIL_NOT_BUILDING)
			return
		}

		currentBuilding := resourcePoint.Building()

		// 建筑已经最高级
		nextLevelId := domestic_data.BuildingId(currentBuilding.Type, currentBuilding.Level+1)
		nextLevel := m.datas.GetBuildingData(nextLevelId)
		if nextLevel == nil {
			logrus.Debugf("升级建筑，已经是最高级了")
			result.Add(domestic.ERR_UPGRADE_BUILDING_FAIL_MAX_LEVEL)
			return
		}

		resEffect := nextLevel.GetResPointEffect()
		if resEffect == nil {
			logrus.Debugf("升级建筑，这个建筑不是资源点")
			result.Add(domestic.ERR_UPGRADE_BUILDING_FAIL_NOT_BUILDING)
			return
		}

		ctime := m.timeService.CurrentTime()

		var workerPos int
		if nextLevel.WorkTime > 0 {
			// 建筑队休息
			workerPos = heroDomestic.GetFreeWorker(ctime)
			if workerPos < 0 {
				logrus.Debugf("升级建筑，建筑队休息中")
				result.Add(domestic.ERR_UPGRADE_BUILDING_FAIL_WORKER_REST)
				return
			}
		}

		// 前提条件
		for _, requireBuilding := range nextLevel.RequireIds {
			heroBuilding := heroDomestic.GetBuilding(requireBuilding.Type)
			if heroBuilding == nil || heroBuilding.Level < requireBuilding.Level {
				logrus.Debugf("升级建筑，前置条件未达成")
				result.Add(domestic.ERR_UPGRADE_BUILDING_FAIL_REQUIRE_NOT_REACH)
				return
			}
		}

		// 资源点不在你的势力范围，请先升级主城
		if hero.BaseLevel() < resourcePoint.LayoutData().RequireBaseLevel {
			logrus.Debugf("升级建筑，资源点不足势力范围")
			result.Add(domestic.ERR_UPGRADE_BUILDING_FAIL_RESOURCE_INVALID)
			return
		}

		// 资源点冲突
		if hero.IsConflictResourcePoint(resourcePoint.LayoutData()) {
			logrus.Debugf("升级建筑，资源点冲突")
			result.Add(domestic.ERR_UPGRADE_BUILDING_FAIL_RESOURCE_CONFLICT)
			return
		}

		// 主城等级要求
		if nextLevel.BaseLevel != nil && hero.BaseLevel() < nextLevel.BaseLevel.Level {
			logrus.Debugf("升级建筑，主城等级不足")
			result.Add(domestic.ERR_UPGRADE_BUILDING_FAIL_REQUIRE_NOT_REACH)
			return
		}

		// 计算一下当前可以收获多少
		resourcePoint.CalculateCurrentOutput(hero, ctime)

		if !heromodule.TryReduceCostWithDianquanConvert(hctx, hero, result, hero.BuildingEffect().GetBuildingCost(nextLevel.Cost)) {
			logrus.Debugf("升级建筑，资源不足")
			result.Add(domestic.ERR_UPGRADE_BUILDING_FAIL_RESOURCE_NOT_ENOUGH)
			return
		}

		hctx := heromodule.NewContext(m.dep, operate_type.DomesticUpgradeResourceBuilding)
		heromodule.AddExp(hctx, hero, result, nextLevel.HeroExp, ctime)

		var workerRestEndTime time.Time
		var seekHelp bool
		if nextLevel.WorkTime > 0 {
			workerRestEndTime, seekHelp = heroDomestic.AddWorkerRestEndTime(workerPos, ctime, nextLevel.CalculateWorkTime(heroDomestic.GetBuildingWorkerCoef(m.seasonService.Season().WorkerCdr)))
		} else {
			_, workerRestEndTime = heroDomestic.GetWorkerRestEndTime(workerPos)
		}
		resourcePoint.SetBuilding(nextLevel)

		// 计算新的产出，发给客户端
		result.Add(domestic.NewS2cUpgradeBuildingMsg(u64.Int32(layoutId), u64.Int32(nextLevel.Id), imath.Int32(workerPos), timeutil.Marshal32(workerRestEndTime)))

		heromodule.SendResourcePointUpdateMsg(hero, result, resourcePoint)

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_RESOURCE_POINT_COUNT)

		if seekHelp {
			// 发消息通知，可以求助
			result.Add(guild.NewS2cUpdateSeekHelpMsg(constants.SeekTypeWorker, imath.Int32(workerPos), seekHelp))
		}

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DomesticModule) ProcessRebuildBuilding(proto *domestic.C2SRebuildResourceBuildingProto, hc iface.HeroController) {
	layoutId := u64.FromInt32(proto.GetId())

	// 位置无效
	layoutData := m.datas.GetBuildingLayoutData(layoutId)
	if layoutData == nil {
		logrus.Debugf("改建建筑，布局ID无效")
		hc.Send(domestic.ERR_REBUILD_RESOURCE_BUILDING_FAIL_INVALID_LAYOUT)
		return
	}

	// 只有一种类型，改建个P
	if len(layoutData.BuildingType) <= 1 {
		logrus.Debugf("改建建筑，这个建筑不能重建")
		hc.Send(domestic.ERR_REBUILD_RESOURCE_BUILDING_FAIL_INVALID_BUILDING)
		return
	}

	rebuildType := shared_proto.BuildingType(proto.GetType())

	invalidType := true
	for _, v := range layoutData.BuildingType {
		if v == rebuildType {
			invalidType = false
			break
		}
	}

	if invalidType {
		logrus.Debugf("改建建筑，重建类型无效")
		hc.Send(domestic.ERR_REBUILD_RESOURCE_BUILDING_FAIL_INVALID_BUILDING)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.DomesticRebuildBuilding)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroDomestic := hero.Domestic()

		// 空的
		resourcePoint := heroDomestic.GetLayoutRes(layoutId)
		if resourcePoint == nil {
			logrus.Debugf("改建建筑，这个建筑还未新建")
			result.Add(domestic.ERR_REBUILD_RESOURCE_BUILDING_FAIL_INVALID_BUILDING)
			return
		}
		originBuilding := resourcePoint.Building()

		ctime := m.timeService.CurrentTime()

		// 改建成同级的
		rebuildId := domestic_data.BuildingId(rebuildType, originBuilding.Level)
		newBuilding := m.datas.GetBuildingData(rebuildId)
		if newBuilding == nil {
			// 为毛没有数据 ???
			logrus.Error("改建建筑没找到，改建类型同级的建筑数据，%s-%s id: %s", proto.GetType(), originBuilding.Level, rebuildId)
			result.Add(domestic.ERR_REBUILD_RESOURCE_BUILDING_FAIL_INVALID_BUILDING)
			return
		}

		resEffect := newBuilding.GetResPointEffect()
		if resEffect == nil {
			logrus.Debugf("改建建筑，改建的类型不是资源点")
			hc.Send(domestic.ERR_REBUILD_RESOURCE_BUILDING_FAIL_INVALID_BUILDING)
			return
		}

		// 建筑队休息
		var workerPos int
		if newBuilding.WorkTime > 0 {
			workerPos = heroDomestic.GetFreeWorker(ctime)
			if workerPos < 0 {
				logrus.Debugf("改建建筑，建筑队休息中")
				result.Add(domestic.ERR_REBUILD_RESOURCE_BUILDING_FAIL_WORKER_REST)
				return
			}
		}

		// 前提条件 改建不涉及
		// 资源不足 改建不涉及

		// 资源点不在你的势力范围，请先升级主城
		if hero.BaseLevel() < layoutData.RequireBaseLevel {
			logrus.Debugf("改建建筑，资源点不在你的势力范围")
			result.Add(domestic.ERR_REBUILD_RESOURCE_BUILDING_FAIL_RESOURCE_INVALID)
			return
		}

		// 资源点冲突
		if hero.IsConflictResourcePoint(layoutData) {
			logrus.Debugf("改建建筑，资源点冲突")
			result.Add(domestic.ERR_REBUILD_RESOURCE_BUILDING_FAIL_RESOURCE_CONFLICT)
			return
		}

		// 先收取之前的资源
		originResType, ok := domestic_data.GetBuildingResType(originBuilding.Type)
		if ok {
			resourcePoint.CalculateCurrentOutput(hero, ctime)
			if resourcePoint.OutputAmount() > 0 {
				heromodule.AddUnsafeSingleResource(hctx, hero, result, originResType, resourcePoint.OutputAmount())
			}
		}

		var workerRestEndTime time.Time
		var seekHelp bool
		if newBuilding.WorkTime > 0 {
			workerRestEndTime, seekHelp = heroDomestic.AddWorkerRestEndTime(workerPos, ctime, newBuilding.CalculateWorkTime(heroDomestic.GetBuildingWorkerCoef(m.seasonService.Season().WorkerCdr)))
		} else {
			_, workerRestEndTime = heroDomestic.GetWorkerRestEndTime(workerPos)
		}
		resourcePoint.SetBuilding(newBuilding)

		// 计算产出
		result.Add(domestic.NewS2cRebuildResourceBuildingMsg(u64.Int32(layoutId), u64.Int32(newBuilding.Id), int32(workerPos), timeutil.Marshal32(workerRestEndTime)))

		heromodule.SendResourcePointUpdateMsg(hero, result, resourcePoint)

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_RESOURCE_POINT_COUNT)

		if seekHelp {
			// 发消息通知，可以求助
			result.Add(guild.NewS2cUpdateSeekHelpMsg(constants.SeekTypeWorker, imath.Int32(workerPos), seekHelp))
		}

		gamelogs.UpgradeBuildingLog(hero.Pid(), hero.Sid(), hero.Id(), uint64(newBuilding.Type), newBuilding.Level)

		result.Changed()
		result.Ok()
	})

}

//gogen:iface
func (m *DomesticModule) ProcessCollectResource(proto *domestic.C2SCollectResourceProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticCollectResource)

	layoutId := u64.FromInt32(proto.GetId())

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroDomestic := hero.Domestic()

		// 还没建这个建筑
		resourcePoint := heroDomestic.GetLayoutRes(layoutId)
		if resourcePoint == nil {
			logrus.Debugf("收获资源点，这个建筑还未新建")
			result.Add(domestic.ERR_COLLECT_RESOURCE_FAIL_INVALID_LAYOUT)
			return
		}
		building := resourcePoint.Building()

		// 不是资源点
		resType, ok := domestic_data.GetBuildingResType(building.Type)
		if !ok {
			logrus.Debugf("收获资源点，这个建筑不是资源点")
			result.Add(domestic.ERR_COLLECT_RESOURCE_FAIL_CANT_COLLECTED)
			return
		}

		// 资源点不在你的势力范围，请先升级主城
		if hero.BaseLevel() < resourcePoint.LayoutData().RequireBaseLevel {
			logrus.Debugf("收获资源点，资源点不在你的势力范围")
			result.Add(domestic.ERR_COLLECT_RESOURCE_FAIL_RESOURCE_INVALID)
			return
		}

		// 资源点冲突
		if hero.IsConflictResourcePoint(resourcePoint.LayoutData()) {
			logrus.Debugf("收获资源点，资源点冲突")
			result.Add(domestic.ERR_COLLECT_RESOURCE_FAIL_RESOURCE_CONFLICT)
			return
		}

		ctime := m.timeService.CurrentTime()
		resourcePoint.CalculateCurrentOutput(hero, ctime)
		if resourcePoint.OutputAmount() <= 0 {
			// 没有可以收获的
			logrus.Debugf("收获资源点，资源点没有资源可以收")
			result.Add(domestic.ERR_COLLECT_RESOURCE_FAIL_FULL)
			return
		}

		//curRes := hero.GetUnsafeResource().GetRes(resType)
		//totalCapcity := heroDomestic.GetCapcity(resType)
		//if curRes >= totalCapcity {
		//	// 没有可以收获的
		//	logrus.Debugf("收获资源点，超出仓库上限")
		//	result.Add(domestic.ERR_COLLECT_RESOURCE_FAIL_FULL)
		//	return
		//}
		//toAddAmount := u64.Min(resourcePoint.OutputAmount(), u64.Sub(totalCapcity, curRes))
		toAddAmount := resourcePoint.OutputAmount()

		// 从资源点扣掉，加到英雄身上
		resourcePoint.ReduceOutputAmount(toAddAmount)
		heromodule.AddUnsafeSingleResource(hctx, hero, result, resType, toAddAmount)

		result.Add(domestic.NewS2cCollectResourceMsg(u64.Int32(layoutId), u64.Int32(toAddAmount)))

		// 历史累计数量
		hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CollectResource, toAddAmount)
		switch resType {
		case shared_proto.ResType_GOLD:
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CollectResourceGold, toAddAmount)
		case shared_proto.ResType_FOOD:
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CollectResourceFood, toAddAmount)
		case shared_proto.ResType_WOOD:
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CollectResourceWood, toAddAmount)
		case shared_proto.ResType_STONE:
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CollectResourceStone, toAddAmount)
		}

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_COLLECT_RESOURCE)

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DomesticModule) ProcessCollectResourceV2(proto *domestic.C2SCollectResourceV2Proto, hc iface.HeroController) {
	if _, ok := shared_proto.ResType_name[proto.ResType]; proto.ResType == 0 || !ok {
		logrus.WithField("res type", proto.ResType).Debugln("非法的资源类型")
		hc.Send(domestic.ERR_COLLECT_RESOURCE_V2_FAIL_INVALID_RESOURCE_TYPE)
		return
	}

	resType := shared_proto.ResType(proto.ResType)

	singleOutPut := m.datas.BuildingLayoutMiscData().SingleResOutPutAmount(resType)
	if singleOutPut == nil {
		logrus.Debugln("没有资源可以收集")
		hc.Send(domestic.ERR_COLLECT_RESOURCE_V2_FAIL_EMPTY)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.DomesticCollectResourceV2)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDomestic := hero.Domestic()

		//curRes := hero.GetRes(resType)
		//totalCapcity := heroDomestic.GetCapcity(resType)
		//if curRes >= totalCapcity {
		//	// 没有可以收获的
		//	logrus.Debugf("收获资源点，超出仓库上限")
		//	result.Add(domestic.ERR_COLLECT_RESOURCE_V2_FAIL_FULL)
		//	return
		//}

		resourceCollectTimes := heroDomestic.ResourceCollectTimes()
		if resourceCollectTimes >= m.datas.MiscConfig().MaxResourceCollectTimes {
			logrus.Debugf("收获资源点，没有次数了")
			hc.Send(domestic.ERR_COLLECT_RESOURCE_V2_FAIL_NO_TIMES)
			return
		}

		ctime := m.timeService.CurrentTime()
		if ctime.Before(heroDomestic.GetResourceNextCollectTime(resType)) {
			logrus.Debugf("收获资源点，倒计时未结束")
			hc.Send(domestic.ERR_COLLECT_RESOURCE_V2_FAIL_COUNTDOWN)
			return
		}

		var output uint64

		if extraOutput := heroDomestic.GetExtraOutput(resType); extraOutput != nil {
			output = data.TotalAmount(singleOutPut, extraOutput)
		} else {
			output = data.TotalAmount(singleOutPut)
		}

		if output <= 0 {
			logrus.Debugln("计算 output，发现采集量竟然小于0")
			hc.Send(domestic.ERR_COLLECT_RESOURCE_V2_FAIL_EMPTY)
			return
		}

		var toAddAmount uint64

		cubes := m.datas.RegionConfig().GetEvenOffsetCubesIncludeLowLevel(hero.BaseLevel())
		for _, cb := range cubes {
			layoutData := m.datas.RegionConfig().GetLayoutDataByEvenOffset(cb)
			if layoutData == nil {
				continue
			}

			if hero.IsConflictResourcePoint(layoutData) {
				continue
			}

			// 加资源量
			toAddAmount += output
		}

		if toAddAmount <= 0 {
			logrus.Debugln("没有资源可以收集")
			hc.Send(domestic.ERR_COLLECT_RESOURCE_V2_FAIL_EMPTY)
			return
		}

		//resourceCollectTimes.ReduceOneTimes(ctime)
		newTimes := heroDomestic.IncreseResourceCollectTimes()
		nextCollectTime := ctime.Add(m.datas.MiscConfig().ResourceRecoveryDuration)
		heroDomestic.SetResourceNextCollectTime(resType, nextCollectTime)

		heromodule.AddUnsafeSingleResource(hctx, hero, result, resType, toAddAmount)

		result.Add(domestic.NewS2cCollectResourceV2Msg(proto.ResType, u64.Int32(toAddAmount), u64.Int32(newTimes), timeutil.Marshal32(nextCollectTime)))

		// 历史累计数量
		hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CollectResource, toAddAmount)
		switch resType {
		case shared_proto.ResType_GOLD:
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CollectResourceGold, toAddAmount)
		case shared_proto.ResType_FOOD:
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CollectResourceFood, toAddAmount)
		case shared_proto.ResType_WOOD:
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CollectResourceWood, toAddAmount)
		case shared_proto.ResType_STONE:
			hero.HistoryAmount().Increase(server_proto.HistoryAmountType_CollectResourceStone, toAddAmount)
		}

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_COLLECT_RESOURCE)
		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_COLLECT_RESOURCE_TIMES)

		result.Changed()
		result.Ok()
	})
}

//gogen:iface c2s_request_resource_conflict
func (m *DomesticModule) ProcessRequestResourceConflict(hc iface.HeroController) {
	var homeRealm iface.Realm

	var flags []string
	var names []string

	if !hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if hero.BaseRegion() == 0 {
			return
		}

		homeRealm = m.realmService.GetBigMap()
		if homeRealm == nil {
			return
		}

		hero.WalkHomeNpcBase(func(base *entity.HomeNpcBase) bool {
			cube := cb.XYCube(base.GetData().EvenOffsetX, base.GetData().EvenOffsetY)
			layoutData := m.datas.RegionConfig().GetLayoutDataByEvenOffset(cube)
			if layoutData != nil && layoutData.RequireBaseLevel <= hero.BaseLevel() {
				flags = append(flags, "") // 没有帮旗
				names = append(names, base.GetData().Data.Npc.Name)
			}
			return false
		})
	}) {
		hc.Send(domestic.ERR_REQUEST_RESOURCE_CONFLICT_FAIL_SERVER_BUSY)
		return
	}

	if homeRealm == nil {
		hc.Send(domestic.ERR_REQUEST_RESOURCE_CONFLICT_FAIL_SERVER_BUSY)
		return
	}

	suc, conflictIds := homeRealm.GetConflictHeroIds(hc)
	if !suc {
		hc.Send(domestic.ERR_REQUEST_RESOURCE_CONFLICT_FAIL_SERVER_BUSY)
		return
	}

	for _, heroId := range conflictIds {
		if heroId == hc.Id() {
			continue
		}

		snapshot := m.heroSnapshotService.Get(heroId)
		if snapshot != nil {
			flags = append(flags, snapshot.GuildFlagName())
			names = append(names, snapshot.Name)
		}
	}

	hc.Send(domestic.NewS2cRequestResourceConflictMsg(flags, names))
}

//gogen:iface
func (m *DomesticModule) ProcessLearnTechnology(proto *domestic.C2SLearnTechnologyProto, hc iface.HeroController) {
	group := u64.FromInt32(proto.GetId())

	//technologyData := m.datas.GetTechnologyData(techId)
	//if technologyData == nil {
	//	logrus.Debugf("升级科技，科技id不存在")
	//	hc.Send(domestic.ERR_LEARN_TECHNOLOGY_FAIL_INVALID_ID)
	//	return
	//}
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticLearnTechnology)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroDomestic := hero.Domestic()

		learned := heroDomestic.GetTechnology(group)
		nextLevel := learned
		if learned != nil {
			nextLevel = learned.GetNextLevel()
			if nextLevel == nil {
				logrus.Debugf("升级科技，已经升到最高级")
				hc.Send(domestic.ERR_LEARN_TECHNOLOGY_FAIL_NOT_NEXT_LEVEL)
				return
			}
		} else {
			nextLevel = m.techFirstLevelMap[group]
			if nextLevel == nil {
				logrus.Debugf("升级科技，无效的group")
				hc.Send(domestic.ERR_LEARN_TECHNOLOGY_FAIL_INVALID_ID)
				return
			}
		}

		for _, requireBuilding := range nextLevel.RequireBuildingIds {
			building := heroDomestic.GetBuilding(requireBuilding.Type)
			if building == nil || building.Level < requireBuilding.Level {
				logrus.Debugf("升级科技，前置建筑等级不够")
				hc.Send(domestic.ERR_LEARN_TECHNOLOGY_FAIL_PRE_BUILDING_LEVEL_INVALID)
				return
			}
		}

		for _, requireTech := range nextLevel.RequireTechIds {
			tech := heroDomestic.GetTechnology(requireTech.Group)
			if tech == nil || tech.Level < requireTech.Level {
				logrus.Debugf("升级科技，前置科技等级不够")
				hc.Send(domestic.ERR_LEARN_TECHNOLOGY_FAIL_PRE_TECH_LEVEL_INVALID)
				return
			}
		}

		// 科研队休息
		ctime := m.timeService.CurrentTime()
		workerPos := heroDomestic.GetFreeTechWorker(ctime)
		if workerPos < 0 {
			logrus.Debugf("升级科技，科研队休息中")
			result.Add(domestic.ERR_LEARN_TECHNOLOGY_FAIL_WORKER_REST)
			return
		}

		if !heromodule.TryReduceCostWithDianquanConvert(hctx, hero, result, hero.BuildingEffect().GetTechCost(nextLevel.Cost)) {
			logrus.Debugf("升级科技，资源不足")
			result.Add(domestic.ERR_LEARN_TECHNOLOGY_FAIL_RESOURCE_NOT_ENOUGH)
			return
		}

		workerRestEndTime, seekHelp := heroDomestic.AddTechWorkerRestEndTime(workerPos, ctime, nextLevel.CalculateWorkTime(heroDomestic.GetTechWorkerCoef(m.seasonService.Season().WorkerCdr)))
		heroDomestic.SetTechnology(nextLevel.Group, nextLevel)
		result.Add(domestic.NewS2cLearnTechnologyMsg(u64.Int32(nextLevel.Id), int32(workerPos), timeutil.Marshal32(workerRestEndTime)))

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_TECH_LEVEL)
		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_TECH_LEVEL)

		heromodule.UpdateBuildingEffect(hero, result, m.datas, ctime, nextLevel.Effect)

		if seekHelp {
			result.Add(guild.NewS2cUpdateSeekHelpMsg(constants.SeekTypeTech, int32(workerPos), seekHelp))
		}

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_LEARN_TECH) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_LEARN_TECH)))
		}

		result.Changed()
		result.Ok()

		var oldLevel uint64
		if learned != nil {
			oldLevel = learned.Level
		}

		// tlog
		hctx.Tlog().TlogResearchFlow(hero, operate_type.ResearchTech, group, oldLevel, nextLevel.Level)

		gamelogs.UpgradeTechLog(hero.Pid(), hero.Sid(), hero.Id(), group, nextLevel.Level)
	})
}

//gogen:iface c2s_unlock_stable_building
func (m *DomesticModule) ProcessUnlockStableBuilding(proto *domestic.C2SUnlockStableBuildingProto, hc iface.HeroController) {

	if _, exist := shared_proto.BuildingType_name[proto.GetType()]; !exist {
		logrus.Debugf("解锁建筑，无效的类型, %v", proto.GetType())
		hc.Send(domestic.ERR_UNLOCK_STABLE_BUILDING_FAIL_INVALID_TYPE)
		return
	}

	buildingType := shared_proto.BuildingType(proto.GetType())
	if buildingType == shared_proto.BuildingType_InvalidBuildingType {
		logrus.Debugf("解锁建筑，无效的类型, %v", proto.GetType())
		hc.Send(domestic.ERR_UNLOCK_STABLE_BUILDING_FAIL_INVALID_TYPE)
		return
	}

	var realm iface.Realm
	var toAddProsperity uint64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroDomestic := hero.Domestic()

		currentBuilding := heroDomestic.GetBuilding(buildingType)
		if currentBuilding != nil {
			logrus.Debugln("解锁建筑，建筑已经解锁了")
			hc.Send(domestic.ERR_UNLOCK_STABLE_BUILDING_FAIL_UNLOCKED)
			return
		}

		unlockData := m.datas.GetBuildingUnlockData(uint64(buildingType))
		if unlockData == nil {
			logrus.Debugln("解锁建筑，解锁数据没找到，这种默认就解锁了")
			hc.Send(domestic.ERR_UNLOCK_STABLE_BUILDING_FAIL_UNLOCKED)
			return
		}

		if unlockData.GuanFuLevel != nil {
			if guanFu := heroDomestic.GetBuilding(shared_proto.BuildingType_GUAN_FU); guanFu == nil || unlockData.GuanFuLevel.Level > guanFu.Level {
				logrus.Debugln("解锁建筑，官府等级不够")
				hc.Send(domestic.ERR_UNLOCK_STABLE_BUILDING_FAIL_GUAN_FU_LEVEL_NOT_ENOUGH)
				return
			}
		}

		if unlockData.HeroLevel > hero.Level() {
			logrus.Debugln("解锁建筑，君主等级不够")
			hc.Send(domestic.ERR_UNLOCK_STABLE_BUILDING_FAIL_HERO_LEVEL_NOT_ENOUGH)
			return
		}

		if unlockData.MainTaskSequence > hero.TaskList().CompletedMainTaskSequence() {
			logrus.Debugln("解锁建筑，主线任务未达到")
			hc.Send(domestic.ERR_UNLOCK_STABLE_BUILDING_FAIL_MAIN_TASK_NOT_REACH)
			return
		}

		if unlockData.BaYeStage > hero.TaskList().GetCompletedBaYeStage() {
			logrus.Debugln("解锁建筑，霸业目标未达到")
			hc.Send(domestic.ERR_UNLOCK_STABLE_BUILDING_FAIL_BA_YE_STAGE_NOT_REACH)
			return
		}

		if hero.BaseLevel() <= 0 {
			logrus.Debugf("接受建筑，处于流亡状态")
			hc.Send(domestic.ERR_UNLOCK_STABLE_BUILDING_FAIL_BASE_DEAD)
			return
		}

		if hero.BaseRegion() != 0 {
			realm = m.realmService.GetBigMap()
		}

		heroDomestic.SetBuilding(unlockData.UnlockBuildingData)

		toAddProsperity = unlockData.UnlockBuildingData.Prosperity
		if currentBuilding != nil {
			toAddProsperity = u64.Sub(toAddProsperity, currentBuilding.Prosperity)
		}

		if toAddProsperity > 0 {
			hero.AddProsperityCapcity(toAddProsperity)
		}

		ctime := m.timeService.CurrentTime()
		hctx := heromodule.NewContext(m.dep, operate_type.DomesticUnlockStableBuilding)
		heromodule.AddExp(hctx, hero, result, unlockData.UnlockBuildingData.HeroExp, ctime)

		result.Add(domestic.NewS2cUnlockStableBuildingMsg(u64.Int32(unlockData.UnlockBuildingData.Id)))

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL)
		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL)

		heromodule.UpdateBuildingEffect(hero, result, m.datas, ctime, unlockData.UnlockBuildingData.Effect)

		m.updateBuildingLevel(hero, result, currentBuilding, unlockData.UnlockBuildingData, ctime)

		result.Changed()
		result.Ok()
	}) {
		return
	}

	if realm != nil && toAddProsperity > 0 {
		processed, err := realm.AddProsperity(hc.Id(), toAddProsperity)
		if err != nil {
			logrus.WithError(err).Error("英雄解锁建筑，加繁荣度出错")
		} else if !processed {
			logrus.Error("英雄解锁建筑，加繁荣度超时")
		}
	}
}

//gogen:iface
func (m *DomesticModule) ProcessUpgradeStableBuilding(proto *domestic.C2SUpgradeStableBuildingProto, hc iface.HeroController) {

	if _, exist := shared_proto.BuildingType_name[proto.GetType()]; !exist {
		logrus.Debugf("升级城内建筑，无效的类型, %v", proto.GetType())
		hc.Send(domestic.ERR_UPGRADE_STABLE_BUILDING_FAIL_INVALID_TYPE)
		return
	}

	buildingType := shared_proto.BuildingType(proto.GetType())
	if buildingType == shared_proto.BuildingType_InvalidBuildingType {
		logrus.Debugf("升级城内建筑，无效的类型, %v", proto.GetType())
		hc.Send(domestic.ERR_UPGRADE_STABLE_BUILDING_FAIL_INVALID_TYPE)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.DomesticUpgradeStableBuilding)

	var oldLevel, newLevel uint64
	var addProsperityFunc func()
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroDomestic := hero.Domestic()

		currentBuilding := heroDomestic.GetBuilding(buildingType)
		nextLevelAmount := uint64(1)
		if currentBuilding != nil {
			nextLevelAmount = currentBuilding.Level + 1
			oldLevel = currentBuilding.Level
		}

		if u64.FromInt32(proto.Level) != u64.Sub(nextLevelAmount, 1) {
			logrus.Debugf("升级城内建筑，等级不一致")
			result.Add(domestic.ERR_UPGRADE_STABLE_BUILDING_FAIL_DIFF_LEVEL)
			return
		}

		// 建筑已经最高级
		nextLevelId := domestic_data.BuildingId(buildingType, nextLevelAmount)
		nextLevel := m.datas.GetBuildingData(nextLevelId)
		if nextLevel == nil {
			logrus.Debugf("升级城内建筑，已经是最高级了")
			result.Add(domestic.ERR_UPGRADE_STABLE_BUILDING_FAIL_MAX_LEVEL)
			return
		}

		// 建筑队休息
		ctime := m.timeService.CurrentTime()
		workerPos := heroDomestic.GetFreeWorker(ctime)
		if workerPos < 0 {
			logrus.Debugf("升级城内建筑，建筑队休息中")
			result.Add(domestic.ERR_UPGRADE_STABLE_BUILDING_FAIL_WORKER_REST)
			return
		}

		// 前提条件
		for _, requireBuilding := range nextLevel.RequireIds {
			heroBuilding := heroDomestic.GetBuilding(requireBuilding.Type)
			if heroBuilding == nil || heroBuilding.Level < requireBuilding.Level {
				logrus.Debugf("升级城内建筑，前置建筑还没有")
				result.Add(domestic.ERR_UPGRADE_STABLE_BUILDING_FAIL_REQUIRE_NOT_REACH)
				return
			}
		}

		if hero.BaseLevel() <= 0 {
			logrus.Debugf("升级城内建筑，处于流亡状态")
			result.Add(domestic.ERR_UPGRADE_STABLE_BUILDING_FAIL_BASE_DEAD)
			return
		}

		// 主城等级要求
		if nextLevel.BaseLevel != nil && hero.BaseLevel() < nextLevel.BaseLevel.Level {
			logrus.Debugf("升级城内建筑，主城等级不足")
			result.Add(domestic.ERR_UPGRADE_STABLE_BUILDING_FAIL_REQUIRE_NOT_REACH)
			return
		}

		// 扣钱
		if !heromodule.TryReduceCostWithDianquanConvert(hctx, hero, result, hero.BuildingEffect().GetBuildingCost(nextLevel.Cost)) {
			logrus.Debugf("升级城内建筑，资源不足")
			result.Add(domestic.ERR_UPGRADE_STABLE_BUILDING_FAIL_RESOURCE_NOT_ENOUGH)
			return
		}

		heroDomestic.SetBuilding(nextLevel)

		toAddProsperity := nextLevel.Prosperity
		if currentBuilding != nil {
			toAddProsperity = u64.Sub(toAddProsperity, currentBuilding.Prosperity)
		}

		if toAddProsperity > 0 {
			hero.AddProsperityCapcity(toAddProsperity)
		}

		heromodule.AddExp(hctx, hero, result, nextLevel.HeroExp, ctime)

		workerRestEndTime, seekHelp := heroDomestic.AddWorkerRestEndTime(workerPos, ctime, nextLevel.CalculateWorkTime(heroDomestic.GetBuildingWorkerCoef(m.seasonService.Season().WorkerCdr)))

		result.Add(domestic.NewS2cUpgradeStableBuildingMsg(u64.Int32(nextLevel.Id), int32(workerPos), timeutil.Marshal32(workerRestEndTime)))

		heromodule.UpdateBuildingEffect(hero, result, m.datas, ctime, nextLevel.Effect)

		m.updateBuildingLevel(hero, result, currentBuilding, nextLevel, ctime)
		newLevel = nextLevel.Level

		if seekHelp {
			// 发消息通知，可以求助
			result.Add(guild.NewS2cUpdateSeekHelpMsg(constants.SeekTypeWorker, imath.Int32(workerPos), seekHelp))
		}

		// 解锁外城
		totalAddProsperity := toAddProsperity
		//if hero.Domestic().OuterCities().TryAutoUnlock() {
		//	toAdd := heromodule.TryUnlockOutsideCity(hero, result, m.datas, ctime)
		//	totalAddProsperity += toAdd
		//} else {
		//	// 更新任务进度
		//	heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL)
		//	heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL)
		//}

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL)
		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL)

		addProsperityFunc = m.realmService.AddProsperityFunc(hero.Id(), hero.BaseRegion(), totalAddProsperity, "解锁建筑")

		result.Changed()
		result.Ok()

		m.dep.Tlog().TlogStrenghBuildingFlow(hero, uint64(buildingType), oldLevel, newLevel)

		gamelogs.UpgradeBuildingLog(hero.Pid(), hero.Sid(), hero.Id(), uint64(buildingType), newLevel)

	}) {
		return
	}

	if addProsperityFunc != nil {
		addProsperityFunc()
	}

}

// 解锁外城
//gogen:iface
func (m *DomesticModule) ProcessUnlockOuterCity(proto *domestic.C2SUnlockOuterCityProto, hc iface.HeroController) {
	outerCityData := m.datas.GetOuterCityData(u64.FromInt32(proto.Id))
	if outerCityData == nil {
		logrus.Debugf("未知分城: %d", proto.Id)
		hc.Send(domestic.ERR_UNLOCK_OUTER_CITY_FAIL_INVALID_ID)
		return
	}

	var addProsperityFunc func()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDomestic := hero.Domestic()
		outerCities := heroDomestic.OuterCities()

		outerCity := outerCities.OuterCity(outerCityData)
		if outerCity != nil {
			logrus.Debugf("分城已经解锁: %d", proto.Id)
			result.Add(domestic.ERR_UNLOCK_OUTER_CITY_FAIL_UNLOCKED)
			return
		}

		guanFu := heroDomestic.GetBuilding(shared_proto.BuildingType_GUAN_FU)
		if guanFu == nil {
			logrus.Debugf("官府没有配置: %d", proto.Id)
			result.Add(domestic.ERR_UNLOCK_OUTER_CITY_FAIL_REQUIRE_NOT_REACH)
			return
		}

		if !hero.Function().IsFunctionOpened(outerCityData.GetFuncType()) {
			logrus.Debugf("分城解锁条件未达成")
			result.Add(domestic.ERR_UNLOCK_OUTER_CITY_FAIL_REQUIRE_NOT_REACH)
			return
		}

		ctime := m.timeService.CurrentTime()

		toAddProsperity := heromodule.UnlockOutsideCity(hero, result,
			[]*domestic_data.OuterCityData{outerCityData}, []uint64{u64.FromInt32(proto.T)},
			m.datas, ctime)

		addProsperityFunc = m.realmService.AddProsperityFunc(hero.Id(), hero.BaseRegion(), toAddProsperity, "解锁分城")

		result.Ok()
	})

	if addProsperityFunc != nil {
		addProsperityFunc()
	}
}

//// 改建外城
////gogen:iface
//func (m *DomesticModule) ProcessUpdateOuterCityType(proto *domestic.C2SUpdateOuterCityTypeProto, hc iface.HeroController) {
//	outerCityData := m.datas.GetOuterCityData(u64.FromInt32(proto.Id))
//	if outerCityData == nil {
//		logrus.Debugf("改建外城，未知分城: %d", proto.Id)
//		hc.Send(domestic.ERR_UPDATE_OUTER_CITY_TYPE_FAIL_INVALID_ID)
//		return
//	}
//
//	toUpdate := u64.FromInt32(proto.T)
//	if toUpdate > 1 {
//		logrus.Debugf("改建外城，无效的外城类型: %d", proto.T)
//		hc.Send(domestic.ERR_UPDATE_OUTER_CITY_TYPE_FAIL_INVALID_ID)
//		return
//	}
//
//	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
//		heroDomestic := hero.Domestic()
//		outerCities := heroDomestic.OuterCities()
//
//		outerCity := outerCities.OuterCity(outerCityData)
//		if outerCity == nil {
//			logrus.Debugf("改建外城，分城还未解锁: %d", proto.Id)
//			result.Add(domestic.ERR_UPDATE_OUTER_CITY_TYPE_FAIL_LOCKED)
//			return
//		}
//
//		if outerCity.Type() == toUpdate {
//			logrus.Debugf("改建外城，已经是这个类型的外城了")
//			result.Add(domestic.NewS2cUpdateOuterCityTypeMsg(proto.Id, proto.T))
//			return
//		}
//
//		hctx := heromodule.NewContext(m.dep, operate_type.DomesticUpdateOutcityType)
//		if !heromodule.TryReduceCost(hctx, hero, result, m.datas.MiscGenConfig().UpdateOuterCityTypeCost) {
//			logrus.Debugf("改建外城，消耗不足")
//			result.Add(domestic.ERR_UPDATE_OUTER_CITY_TYPE_FAIL_COST_NOT_ENOUGH)
//			return
//		}
//
//		// 更改前的建筑效果
//		var effects []*domestic_data.BuildingEffectData
//		outerCity.WalkLayouts(func(layout *domestic_data.OuterCityLayoutData, building *domestic_data.OuterCityBuildingData) {
//			effects = append(effects, building.BuildingData.Effect)
//		})
//
//		outerCity.SetType(toUpdate)
//
//		// 更改后的建筑效果
//		outerCity.WalkLayouts(func(layout *domestic_data.OuterCityLayoutData, building *domestic_data.OuterCityBuildingData) {
//			effects = append(effects, building.BuildingData.Effect)
//		})
//
//		ctime := m.timeService.CurrentTime()
//		heromodule.UpdateBuildingEffect(hero, result, m.datas, ctime, effects...)
//
//		result.Add(domestic.NewS2cUpdateOuterCityTypeMsg(proto.Id, proto.T))
//
//		result.Ok()
//	})
//
//}

// 新版改建外城（批量降级）
//gogen:iface
func (m *DomesticModule) ProcessUpdateOuterCityType(proto *domestic.C2SUpdateOuterCityTypeProto, hc iface.HeroController) {
	outerCityData := m.datas.GetOuterCityData(u64.FromInt32(proto.GetId()))
	if outerCityData == nil {
		logrus.Debugf("改建外城，未知分城: %d", proto.Id)
		hc.Send(domestic.ERR_UPDATE_OUTER_CITY_TYPE_FAIL_INVALID_ID)
		return
	}

	toUpdate := u64.FromInt32(proto.GetT())
	if toUpdate > 1 {
		logrus.Debugf("改建外城，无效的外城类型: %d", proto.T)
		hc.Send(domestic.ERR_UPDATE_OUTER_CITY_TYPE_FAIL_INVALID_ID)
		return
	}

	var updateProsperityFunc func()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDomestic := hero.Domestic()
		outerCities := heroDomestic.OuterCities()

		outerCity := outerCities.OuterCity(outerCityData)
		if outerCity == nil {
			logrus.Debugf("改建外城，分城还未解锁: %d", proto.Id)
			result.Add(domestic.ERR_UPDATE_OUTER_CITY_TYPE_FAIL_LOCKED)
			return
		}

		if outerCity.Type() == toUpdate {
			logrus.Debugf("改建外城，已经是这个类型的外城了")
			result.Add(domestic.ERR_UPDATE_OUTER_CITY_TYPE_FAIL_SAME_TYPE)
			return
		}

		cost := outerCity.GetChangeTypeCost()
		if cost == nil {
			logrus.Debugf("改建外城，改建消耗竟然为空")
			result.Add(domestic.ERR_UPDATE_OUTER_CITY_TYPE_FAIL_SERVER_BUSY)
			return
		}
		hctx := heromodule.NewContext(m.dep, operate_type.DomesticUpdateOutcityType)
		ctime := m.timeService.CurrentTime()
		if !heromodule.TryReduceCostWithDianquanConvert(hctx, hero, result, cost) {
			logrus.Debugf("改建外城，消耗不足")
			result.Add(domestic.ERR_UPDATE_OUTER_CITY_TYPE_FAIL_COST_NOT_ENOUGH)
			return
		}

		// 更改前的建筑效果
		var effects []*sub.BuildingEffectData
		outerCity.WalkLayouts(func(layout *domestic_data.OuterCityLayoutData, building *domestic_data.OuterCityBuildingData) {
			effects = append(effects, building.BuildingData.Effect)
		})

		// 开始改建
		outerCity.Change()
		outerCity.SetType(toUpdate)
		outerCities.UpdateUnlockBit()

		// 更改后的建筑效果，顺便获取改建后的建筑id列表
		var ids []int32
		outerCity.WalkLayouts(func(layout *domestic_data.OuterCityLayoutData, building *domestic_data.OuterCityBuildingData) {
			effects = append(effects, building.BuildingData.Effect)
			ids = append(ids, u64.Int32(layout.Id))
		})

		heromodule.UpdateBuildingEffect(hero, result, m.datas, ctime, effects...)

		result.Add(domestic.NewS2cUpdateOuterCityTypeMsg(proto.Id, u64.Int32(toUpdate), ids))

		// 更新更新外城类型，需要更新繁荣度，以及外城模型，繁荣度里面会有外城模型改变的处理，这里只处理繁荣度
		origionProsperity := hero.Prosperity()
		newProsperity := origionProsperity
		if hero.UpdateProsperityCapcity() {
			newProsperity = hero.Prosperity()
		}

		fmt.Println(origionProsperity, newProsperity)
		updateProsperityFunc = func() {
			if origionProsperity == newProsperity {
				m.realmService.GetBigMap().UpdateProsperity(hc.Id())
			} else if newProsperity < origionProsperity {
				m.realmService.GetBigMap().ReduceProsperity(hc.Id(), u64.Sub(origionProsperity, newProsperity))
			}
		}

		result.Ok()
	})

	if updateProsperityFunc != nil {
		updateProsperityFunc()
	}
}

// 升级外城建筑
//gogen:iface
func (m *DomesticModule) ProcessUpgradeOuterCityBuilding(proto *domestic.C2SUpgradeOuterCityBuildingProto, hc iface.HeroController) {
	layoutData := m.datas.GetOuterCityLayoutData(u64.FromInt32(proto.Id))
	if layoutData == nil {
		logrus.Debugf("未知分城布局: %d", proto.Id)
		hc.Send(domestic.ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_INVALID_BUILDING_ID)
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.DomesticUpgradeOuterCityBuilding)

	var buildingType shared_proto.BuildingType
	var oldLevel, newLevel uint64
	var realm iface.Realm
	var toAddProsperity uint64
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDomestic := hero.Domestic()
		outerCities := heroDomestic.OuterCities()

		outerCity := outerCities.OuterCity(layoutData.OuterCity)
		if outerCity == nil {
			logrus.Debugf("分城未解锁: %d", proto.Id)
			result.Add(domestic.ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_CITY_LOCKED)
			return
		}

		if hero.BaseLevel() <= 0 {
			logrus.Debugf("升级分城内建筑，处于流亡状态")
			result.Add(domestic.ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_REQUIRE_NOT_REACH)
			return
		}

		var nextLevel *domestic_data.OuterCityLayoutData
		curLayout := outerCity.Layout(layoutData)
		if curLayout != nil {
			oldLevel = curLayout.Level
			nextLevel = curLayout.NextLevel

			// 建筑已经最高级
			if nextLevel == nil {
				logrus.Debugf("升级分城内建筑，已经是最高级了")
				result.Add(domestic.ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_INVALID_BUILDING_ID)
				return
			}
		} else {
			nextLevel = layoutData
			if layoutData.Level != 1 {
				logrus.Debugf("升级分城内建筑，要升级的建筑非法")
				result.Add(domestic.ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_INVALID_BUILDING_ID)
				return
			}
		}

		nextLevelBuilding := nextLevel.GetBuilding(outerCity.Type())

		// 主城等级要求
		if nextLevelBuilding.BuildingData.BaseLevel != nil && hero.BaseLevel() < nextLevelBuilding.BuildingData.BaseLevel.Level {
			logrus.Debugf("升级分城内建筑，主城等级不足")
			result.Add(domestic.ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_REQUIRE_NOT_REACH)
			return
		}

		// 前提条件
		for _, requireBuilding := range nextLevelBuilding.BuildingData.RequireIds {
			heroBuilding := heroDomestic.GetBuilding(requireBuilding.Type)
			if heroBuilding == nil || heroBuilding.Level < requireBuilding.Level {
				logrus.Debugf("升级分城内建筑，前置主城建筑还没有")
				result.Add(domestic.ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_REQUIRE_NOT_REACH)
				return
			}
		}

		for _, requireBuilding := range nextLevel.UpgradeRequireIds {
			heroBuilding := heroDomestic.GetBuilding(requireBuilding.Type)
			if heroBuilding == nil || heroBuilding.Level < requireBuilding.Level {
				logrus.Debugf("升级分城内建筑，前置主城建筑还没有")
				result.Add(domestic.ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_REQUIRE_NOT_REACH)
				return
			}
		}

		if nextLevel.UpgradeRequireLayout != nil {
			// 升级需要其他的布局
			if layout := outerCity.Layout(nextLevel.UpgradeRequireLayout); layout == nil || layout.GetBuilding(outerCity.Type()).BuildingData.Level < nextLevel.UpgradeRequireLayout.Level {
				logrus.Debugf("升级分城内建筑，前置分城建筑等级没达到条件")
				result.Add(domestic.ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_REQUIRE_NOT_REACH)
				return
			}
		}

		// 建筑队休息
		ctime := m.timeService.CurrentTime()
		var workerPos int
		if nextLevelBuilding.BuildingData.WorkTime > 0 {
			workerPos = heroDomestic.GetFreeWorker(ctime)
			if workerPos < 0 {
				logrus.Debugf("升级分城内建筑，建筑队休息中")
				result.Add(domestic.ERR_UPGRADE_STABLE_BUILDING_FAIL_WORKER_REST)
				return
			}
		}

		// 扣钱
		if !heromodule.TryReduceCostWithDianquanConvert(hctx, hero, result, hero.BuildingEffect().GetBuildingCost(nextLevelBuilding.BuildingData.Cost)) {
			logrus.Debugf("升级分城内建筑，资源不足")
			result.Add(domestic.ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_COST_NOT_ENOUGH)
			return
		}

		if hero.BaseRegion() != 0 {
			realm = m.realmService.GetBigMap()
		}

		outerCity.SetLayout(nextLevel)

		toAddProsperity = nextLevelBuilding.BuildingData.Prosperity
		if curLayout != nil {
			toAddProsperity = u64.Sub(toAddProsperity, curLayout.GetBuilding(outerCity.Type()).BuildingData.Prosperity)
		}

		if toAddProsperity > 0 {
			hero.AddProsperityCapcity(toAddProsperity)
		}

		hctx := heromodule.NewContext(m.dep, operate_type.DomesticUpgradeOuterCityBuilding)
		heromodule.AddExp(hctx, hero, result, nextLevelBuilding.BuildingData.HeroExp, ctime)

		var seekHelp bool

		if nextLevelBuilding.BuildingData.WorkTime > 0 {
			var workerRestEndTime time.Time
			workerRestEndTime, seekHelp = heroDomestic.AddWorkerRestEndTime(workerPos, ctime, nextLevelBuilding.BuildingData.CalculateWorkTime(heroDomestic.GetBuildingWorkerCoef(m.seasonService.Season().WorkerCdr)))
			result.Add(domestic.NewS2cBuildingWorkerTimeChangedMsg(int32(workerPos), timeutil.Marshal32(workerRestEndTime)))
		}

		newLevel = nextLevel.Level

		// 更新任务进度
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILDING_LEVEL)
		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_UPGRADE_BUILDING_LEVEL)

		heromodule.UpdateBuildingEffect(hero, result, m.datas, ctime, nextLevelBuilding.BuildingData.Effect)

		if curLayout == nil {
			m.updateBuildingLevel(hero, result, nil, nextLevelBuilding.BuildingData, ctime)
		} else {
			m.updateBuildingLevel(hero, result, curLayout.GetBuilding(outerCity.Type()).BuildingData, nextLevelBuilding.BuildingData, ctime)
		}

		if seekHelp {
			// 发消息通知，可以求助
			result.Add(guild.NewS2cUpdateSeekHelpMsg(constants.SeekTypeWorker, imath.Int32(workerPos), seekHelp))
		}

		result.Add(domestic.NewS2cUpgradeOuterCityBuildingMsg(u64.Int32(layoutData.OuterCity.Id), u64.Int32(layoutData.Id), u64.Int32(nextLevel.Id)))

		result.Changed()
		result.Ok()

		m.dep.Tlog().TlogStrenghBuildingFlow(hero, uint64(buildingType), oldLevel, newLevel)

		gamelogs.UpgradeBuildingLog(hero.Pid(), hero.Sid(), hero.Id(), uint64(buildingType), newLevel)
	})

	if realm != nil {
		processed, err := realm.AddProsperity(hc.Id(), toAddProsperity)
		if err != nil {
			logrus.WithError(err).Error("英雄升级分城建筑，加繁荣度出错")
		} else if !processed {
			logrus.Error("英雄升级分城建筑，加繁荣度超时")
		}
	}
}

// oldLevelData可能为空
func (m *DomesticModule) updateBuildingLevel(hero *entity.Hero, result herolock.LockResult, oldLevelData *domestic_data.BuildingData, nextLevelData *domestic_data.BuildingData, ctime time.Time) {
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticUpdateBuildingLevel)

	// 铁匠铺
	switch nextLevelData.Type {
	case shared_proto.BuildingType_TIE_JIANG_PU:
		//// 这个老版本的不再用了 TODO 20号版本之后注释掉
		//newData := m.datas.TieJiangPuLevelData().Must(nextLevelData.Level)
		//hero.Domestic().AddNewForgingPos(newData.GetNewForgingEquipPos())
		//
		//if newData.GetNewForgingEquipPosMsg() != nil {
		//	result.Add(newData.GetNewForgingEquipPosMsg())
		//}

		// 新版锻造
		if oldLevelData == nil {
			// 1级铁匠铺，刷新一下锻造装备列表
			heromodule.RefreshWorkshopAnyway(hctx, hero, result, m.datas, ctime, operate_type.RefreshAuto, 0)
		}

		//case shared_proto.BuildingType_JUN_YING:
		//	recruitTimes := hero.Military().RecruitTimes()
		//	if oldLevelData == nil {
		//		// 给默认次数
		//		recruitTimes.DefaultGiveTimes(m.datas.JunYingMiscData().DefaultTimes, ctime)
		//	}
		//
		//	newJunYingLevelData := m.datas.JunYingLevelData().Must(nextLevelData.Level)
		//	recruitTimes.ChangeRecoveryDuration(newJunYingLevelData.RecoveryDuration, ctime)
		//	recruitTimes.ChangeMaxTimes(newJunYingLevelData.MaxTimes, ctime)
		//
		//	result.Add(military.NewS2cRecruitSoldierTimesChangedMsg(timeutil.Marshal32(recruitTimes.StartRecoveryTime())))
	case shared_proto.BuildingType_XIU_LIAN_GUAN:

		if oldLevelData != nil {
			// 修改武将的修炼进度，计算每个武将的累积时间
			globalStartTrainTime := hero.Military().GetGlobalTrainStartTime()

			endTime := globalStartTrainTime.Add(heromodule.GetTrainingMaxDuration(hero, m.datas))
			endTime = timeutil.Min(endTime, ctime)

			captainStartTrainTime := timeutil.Max(globalStartTrainTime, hero.Military().GetCaptainTrainStartTime())
			if d := endTime.Sub(captainStartTrainTime); d > 0 {
				trainingExp := u64.Multi(oldLevelData.Effect.TrainExpPerHour, d.Hours())

				hour := u64.Division2Float64(trainingExp, nextLevelData.Effect.TrainExpPerHour)

				newStartTime := endTime.Add(-timeutil.MultiDuration(hour, time.Hour))

				hero.Military().SetCaptainTrainStartTime(newStartTime)
			}

		} else {
			// 解锁修炼馆，修炼馆修炼开始计时
			hero.Military().SetGlobalTrainStartTime(ctime)
		}

		result.Add(military.NewS2cUpdateTrainingMsg(
			timeutil.Marshal32(hero.Military().GetGlobalTrainStartTime()),
			timeutil.Marshal32(hero.Military().GetCaptainTrainStartTime()),
			u64.Int32(nextLevelData.Effect.TrainExpPerHour),
			i32.MultiF64(1000, hero.BuildingEffect().GetTrainCoef())))
	}

	// 只要升级建筑，都有可能解锁功能
	heromodule.CheckFuncsOpened(hero, result)
}

//gogen:iface c2s_is_hero_name_exist
func (m *DomesticModule) ProcessIsHeroNameExist(proto *domestic.C2SIsHeroNameExistProto, hc iface.HeroController) {
	newName := strings.TrimSpace(proto.Name)
	if c := uint64(util.GetCharLen(newName)); c < m.datas.MiscConfig().MinNameCharLen || c > m.datas.MiscConfig().MaxNameCharLen {
		logrus.Debugf("检查君主改名，名字长度不符合规范，newName: %s， len: %d", proto.Name, c)
		hc.Send(domestic.ERR_IS_HERO_NAME_EXIST_FAIL_INVALID_NAME)
		return
	}

	if strings.HasPrefix(newName, constants.PlayerNamePrefix) {
		logrus.Debugf("检查君主改名，不可以以player开头")
		hc.Send(domestic.NewS2cIsHeroNameExistMsg(proto.Name, true))
		return
	}

	var originName string
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		originName = hero.Name()
		result.Ok()
	}) {
		// 已经被踢下线了
		return
	}

	if originName == newName {
		logrus.Debugf("君主改名，名字跟现在名字一样")
		hc.Send(domestic.NewS2cIsHeroNameExistMsg(proto.Name, true))
		return
	}

	heroNameExist := false
	if strings.ToLower(originName) != strings.ToLower(newName) {
		exist := false
		err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
			exist, err = m.dbServuce.HeroNameExist(ctx, newName)
			return
		})
		if err != nil {
			logrus.WithError(err).Errorf("君主改名，查询名字是否存在失败")
			hc.Send(domestic.ERR_IS_HERO_NAME_EXIST_FAIL_SERVER_ERROR)
			return
		}
		heroNameExist = exist
	}

	hc.Send(domestic.NewS2cIsHeroNameExistMsg(proto.Name, heroNameExist))
}

//gogen:iface
func (m *DomesticModule) ProcessChangeHeroName(proto *domestic.C2SChangeHeroNameProto, hc iface.HeroController) {

	newName := strings.TrimSpace(proto.Name)
	if c := uint64(util.GetCharLen(newName)); c < m.datas.MiscConfig().MinNameCharLen || c > m.datas.MiscConfig().MaxNameCharLen {
		logrus.Debugf("君主改名，名字长度不符合规范，newName: %s， len: %d", proto.Name, c)
		hc.Send(domestic.ERR_CHANGE_HERO_NAME_FAIL_INVALID_NAME)
		return
	}

	if strings.HasPrefix(newName, constants.PlayerNamePrefix) {
		logrus.Debugf("君主改名，不可以以player开头")
		hc.Send(domestic.ERR_CHANGE_HERO_NAME_FAIL_EXIST_NAME)
		return
	}

	if !util.IsValidName(newName) {
		logrus.Debugf("君主改名，名字包含非法字符，newName: %s", proto.Name)
		hc.Send(domestic.ERR_CHANGE_HERO_NAME_FAIL_INVALID_NAME)
		return
	}

	if !m.tssClient.TryCheckName("君主改名", hc, newName, domestic.ERR_CHANGE_HERO_NAME_FAIL_SENSITIVE_WORDS, domestic.ERR_CHANGE_HERO_NAME_FAIL_SERVER_ERROR) {
		return
	}

	hctx := heromodule.NewContext(m.dep, operate_type.DomesticChangeHeroName)

	var homeRealm iface.Realm
	var guildId int64
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		originName := hero.Name()
		if originName == newName {
			logrus.Debugf("君主改名，名字跟现在名字一样")
			result.Add(domestic.ERR_CHANGE_HERO_NAME_FAIL_SAME_NAME)
			return
		}

		ctime := m.timeService.CurrentTime()

		if ctime.Before(hero.GetNextChangeNameTime()) {
			logrus.Debugf("君主改名，CD中")
			result.Add(domestic.ERR_CHANGE_HERO_NAME_FAIL_CD)
			return
		}

		if strings.ToLower(originName) != strings.ToLower(newName) {
			exist := false
			err := ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
				exist, err = m.dbServuce.HeroNameExist(ctx, newName)
				return
			})
			if err != nil {
				logrus.WithError(err).Errorf("君主改名，查询名字是否存在失败")
				result.Add(domestic.ERR_CHANGE_HERO_NAME_FAIL_SERVER_ERROR)
				return
			}

			if exist {
				logrus.Debugf("君主改名，名字已经存在")
				result.Add(domestic.ERR_CHANGE_HERO_NAME_FAIL_EXIST_NAME)
				return
			}
		}

		var changeHeroNameTimes uint64 = hero.GetChangeHeroNameTimes()
		changeHeroNameCost := m.datas.MiscConfig().GetChangeHeroNameCost(changeHeroNameTimes)

		if !heromodule.HasEnoughCost(hero, changeHeroNameCost) {
			logrus.Debugf("君主改名，检查消耗不足")
			result.Add(domestic.ERR_CHANGE_HERO_NAME_FAIL_COST_NOT_ENOUGH)
			return
		}

		updateSuccess := false
		ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
			updateSuccess = m.dbServuce.UpdateHeroName(ctx, hc.Id(), originName, newName)
			return
		})
		if !updateSuccess {
			logrus.Debugf("君主改名，改的时候发现名字已经存在")
			result.Add(domestic.ERR_CHANGE_HERO_NAME_FAIL_EXIST_NAME)
			return
		}

		if !heromodule.TryReduceCost(hctx, hero, result, changeHeroNameCost) {
			logrus.Debugf("君主改名，扣消耗不足")

			// 没改成功，改回来
			updateSuccess := false
			ctxfunc.Timeout2s(func(ctx context.Context) (err error) {
				updateSuccess = m.dbServuce.UpdateHeroName(ctx, hc.Id(), newName, originName)
				return
			})
			if updateSuccess {
				// 名字改回来了，返回消耗不足
				result.Add(domestic.ERR_CHANGE_HERO_NAME_FAIL_COST_NOT_ENOUGH)
				return
			}

			logrus.Errorf("君主改名，扣元宝消耗不足，将名字改回来失败，算了这次改名当送给他了, heroId: %v newName: %s originName: %s", hc.Id(), newName, originName)
		}

		// hero
		hero.ChangeName(newName)

		nextChangeNameTime := time.Time{}
		if m.datas.MiscConfig().ChangeHeroNameDuration > 0 {
			nextChangeNameTime = ctime.Add(m.datas.MiscConfig().ChangeHeroNameDuration)
			hero.SetNextChangeNameTime(nextChangeNameTime)
		}

		hero.SetChangeHeroNameTimes(changeHeroNameTimes + 1)

		result.Add(domestic.NewS2cChangeHeroNameMsg(newName, timeutil.Marshal32(nextChangeNameTime)))

		if m.datas.MiscConfig().FirstChangeHeroNamePrize != nil && !hero.HasGiveFirstChangeHeroNamePrize() {
			hctx := heromodule.NewContext(m.dep, operate_type.DomesticChangeHeroName)
			hero.GiveFirstChangeHeroNamePrize()
			heromodule.AddPrize(hctx, hero, result, m.datas.MiscConfig().FirstChangeHeroNamePrize, ctime)
			result.Add(domestic.GIVE_FIRST_CHANGE_HERO_NAME_PRIZE_S2C)
		}

		if hero.BaseRegion() != 0 {
			homeRealm = m.realmService.GetBigMap()
		}

		guildId = hero.GuildId()

		result.Changed()
		result.Ok()
	}) {
		return
	}

	m.world.Broadcast(domestic.NewS2cHeroNameChangedBroadcastMsg(hc.IdBytes(), newName).Static())

	if homeRealm != nil {
		homeRealm.UpdateHeroBasicInfoNoBlock(hc.Id())
	}

	if guildId != 0 {
		if g := m.guildService.GetSnapshot(guildId); g != nil {
			if g.LeaderId == hc.Id() {
				// 盟主改名，更新联盟
				m.guildService.ClearSelfGuildMsgCache(g.Id)
				m.world.MultiSend(g.UserMemberIds, guild.SELF_GUILD_CHANGED_S2C)
			}
		}
	}
}

//gogen:iface
func (m *DomesticModule) ProcessListOldName(proto *domestic.C2SListOldNameProto, hc iface.HeroController) {

	heroId := hc.Id()
	if len(proto.Id) > 0 {
		id, ok := idbytes.ToId(proto.Id)
		if !ok {
			logrus.Debugf("请求曾用名，id无效")
			hc.Send(domestic.ERR_LIST_OLD_NAME_FAIL_INVALID_ID)
			return
		}
		heroId = id
	}

	var toSend pbutil.Buffer
	m.heroDataService.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {

		if err != nil {
			if err == lock.ErrEmpty {
				toSend = domestic.ERR_LIST_OLD_NAME_FAIL_INVALID_ID
				return
			}

			logrus.WithError(err).Errorf("请求曾用名，lock hero fail", hc.Id())
			toSend = domestic.ERR_LIST_OLD_NAME_FAIL_SERVER_ERROR
			return
		}

		toSend = domestic.NewS2cListOldNameMsg(hero.OldName())
		return
	})

	hc.Send(toSend)
}

//gogen:iface
func (m *DomesticModule) ProcessViewOtherHero(proto *domestic.C2SViewOtherHeroProto, hc iface.HeroController) {

	heroId := hc.Id()
	if len(proto.Id) > 0 {
		id, ok := idbytes.ToId(proto.Id)
		if !ok {
			logrus.Debugf("查看其它玩家信息，id无效")
			hc.Send(domestic.ERR_VIEW_OTHER_HERO_FAIL_INVALID_ID)
			return
		}
		heroId = id
	}

	var toSend pbutil.Buffer
	m.heroDataService.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {
		if err != nil {
			if err == lock.ErrEmpty {
				toSend = domestic.ERR_VIEW_OTHER_HERO_FAIL_INVALID_ID
				return
			}

			logrus.WithError(err).Errorf("查看其它玩家信息，lock hero fail", hc.Id())
			toSend = domestic.ERR_VIEW_OTHER_HERO_FAIL_SERVER_ERROR
			return
		}

		toSend = domestic.NewS2cViewOtherHeroProtoMsg(&domestic.S2CViewOtherHeroProto{Hero: must.Marshal(hero.EncodeOther(m.timeService.CurrentTime(), m.guildService.GetSnapshot, m.baiZhanService))})

		return
	})

	hc.Send(toSend)
}

//gogen:iface
func (m *DomesticModule) ProcessViewFightInfo(proto *domestic.C2SViewFightInfoProto, hc iface.HeroController) {

	heroId := hc.Id()
	if len(proto.Id) > 0 {
		id, ok := idbytes.ToId(proto.Id)
		if !ok {
			logrus.Debugf("查看其它玩家战斗状态，id无效")
			hc.Send(domestic.NewS2cViewFightInfoMsg(proto.Id, 0, 0, 0, 0, 0, 0))
			return
		}
		heroId = id
	}

	var toSend pbutil.Buffer
	m.heroDataService.Func(heroId, func(hero *entity.Hero, err error) (heroChanged bool) {
		if err != nil {
			if err == lock.ErrEmpty {
				return
			}

			logrus.WithError(err).Errorf("查看其它玩家战斗状态，lock hero fail", hc.Id())
			return
		}

		toSend = domestic.NewS2cViewFightInfoMsg(
			hero.IdBytes(),
			u64.Int32(hero.HistoryAmount().Amount(server_proto.HistoryAmountType_RealmPvpSuccess)),
			u64.Int32(hero.HistoryAmount().Amount(server_proto.HistoryAmountType_RealmPvpFail)),
			u64.Int32(hero.HistoryAmount().Amount(server_proto.HistoryAmountType_RealmPvpAssist)),
			u64.Int32(hero.HistoryAmount().Amount(server_proto.HistoryAmountType_RealmPvpBeenAssist)),
			u64.Int32(hero.HistoryAmount().Amount(server_proto.HistoryAmountType_Inverstigation)),
			u64.Int32(hero.HistoryAmount().Amount(server_proto.HistoryAmountType_BeenInverstigation)),
		)

		return
	})

	if toSend == nil {
		toSend = domestic.NewS2cViewFightInfoMsg(proto.Id, 0, 0, 0, 0, 0, 0)
	}

	hc.Send(toSend)
}

//gogen:iface
func (m *DomesticModule) ProcessMiaoBuildingWorkerCd(proto *domestic.C2SMiaoBuildingWorkerCdProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticMiaoBuildingWorkerCd)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		idx := int(proto.WorkerPos)
		ok, workTime := hero.Domestic().GetWorkerRestEndTime(idx)
		if !ok {
			logrus.Debugf("秒建筑队cd，无效的WorkerPos, %v", proto.WorkerPos)
			result.Add(domestic.ERR_MIAO_BUILDING_WORKER_CD_FAIL_INVALID_POS)
			return
		}

		ctime := m.timeService.CurrentTime()
		if !workTime.After(ctime) {
			logrus.Debugf("秒建筑队cd，获得的workTime <= ctime")
			result.Add(domestic.ERR_MIAO_BUILDING_WORKER_CD_FAIL_NOT_WORKING)
			return
		}

		d := m.datas.MiscConfig().MiaoBuildingWorkerDuration
		if d <= 0 {
			logrus.Debugf("秒建筑队cd，功能未开放")
			result.Add(domestic.ERR_MIAO_BUILDING_WORKER_CD_FAIL_ZERO_DURATION)
			return
		}

		multi := int64((workTime.Sub(ctime) + d - 1) / d)
		if multi <= 0 {
			logrus.Errorf("秒建筑队cd，计算出来的multi <= 0")
			result.Add(domestic.ERR_MIAO_BUILDING_WORKER_CD_FAIL_NOT_WORKING)
			return
		}

		cost := m.datas.MiscConfig().MiaoBuildingWorkerCost.Multiple(uint64(multi))

		// 扣钱
		if !heromodule.TryReduceCost(hctx, hero, result, cost) {
			logrus.Debugf("秒建筑队cd，钱不够")
			result.Add(domestic.ERR_MIAO_BUILDING_WORKER_CD_FAIL_COST_NOT_ENOUGH)
			return
		}

		hero.Domestic().AddWorkerRestEndTime(idx, ctime, -d*time.Duration(multi))

		result.Add(domestic.NewS2cMiaoBuildingWorkerCdMsg(proto.WorkerPos))

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_WORKER_CDR) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_WORKER_CDR)))
		}

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DomesticModule) ProcessMiaoTechWorkerCd(proto *domestic.C2SMiaoTechWorkerCdProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticMiaoTechWorkerCd)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		idx := int(proto.WorkerPos)
		ok, workTime := hero.Domestic().GetTechWorkerRestEndTime(idx)
		if !ok {
			logrus.Debugf("秒科研队cd，无效的WorkerPos, %v", proto.WorkerPos)
			result.Add(domestic.ERR_MIAO_TECH_WORKER_CD_FAIL_INVALID_POS)
			return
		}

		ctime := m.timeService.CurrentTime()
		if !workTime.After(ctime) {
			logrus.Debugf("秒科研队cd，获得的workTime <= ctime")
			result.Add(domestic.ERR_MIAO_TECH_WORKER_CD_FAIL_NOT_WORKING)
			return
		}

		d := m.datas.MiscConfig().MiaoTechWorkerDuration
		if d <= 0 {
			logrus.Debugf("秒科研队cd，功能未开放")
			result.Add(domestic.ERR_MIAO_TECH_WORKER_CD_FAIL_ZERO_DURATION)
			return
		}

		multi := int64((workTime.Sub(ctime) + d - 1) / d)
		if multi <= 0 {
			logrus.Errorf("秒科研队cd，计算出来的multi <= 0")
			result.Add(domestic.ERR_MIAO_TECH_WORKER_CD_FAIL_NOT_WORKING)
			return
		}

		cost := m.datas.MiscConfig().MiaoTechWorkerCost.Multiple(uint64(multi))

		// 扣钱
		if !heromodule.TryReduceCost(hctx, hero, result, cost) {
			logrus.Debugf("秒科研队cd，钱不够")
			result.Add(domestic.ERR_MIAO_TECH_WORKER_CD_FAIL_COST_NOT_ENOUGH)
			return
		}

		hero.Domestic().AddTechWorkerRestEndTime(idx, ctime, -d*time.Duration(multi))

		result.Add(domestic.NewS2cMiaoTechWorkerCdMsg(proto.WorkerPos))

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_TECH_CDR) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_TECH_CDR)))
		}

		result.Changed()
		result.Ok()
	})
}

// 锻造装备
//gogen:iface
func (m *DomesticModule) ProcessForgingEquip(proto *domestic.C2SForgingEquipProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		heroDomestic := hero.Domestic()

		tieJiangPu := heroDomestic.GetBuilding(shared_proto.BuildingType_TIE_JIANG_PU)

		if tieJiangPu == nil {
			logrus.WithField("proto", proto).Debugln("玩家的铁匠铺没开启")
			result.Add(domestic.ERR_FORGING_EQUIP_FAIL_FUNCTION_NOT_OPEN)
			return
		}

		tieJiangPuLevelData := m.datas.TieJiangPuLevelData().Must(tieJiangPu.Level)

		equipPos := u64.FromInt32(proto.GetSlot())
		equip := tieJiangPuLevelData.GetForgingEquip(equipPos)
		if equip == nil {
			logrus.WithField("proto", proto).Debugln("玩家的铁匠铺里面要打造的装备找不到")
			result.Add(domestic.ERR_FORGING_EQUIP_FAIL_CAN_NOT_FORGING_EQIUP)
			return
		}

		forgingTimes := heroDomestic.GetForgingTimes()
		if forgingTimes.Times() >= tieJiangPuLevelData.MaxForgingTimes {
			logrus.WithField("proto", proto).Debugln("铁匠铺打造，打造次数不足")
			result.Add(domestic.ERR_FORGING_EQUIP_FAIL_TIMES_NOT_ENOUGH)
			return
		}

		ctime := m.timeService.CurrentTime()
		if ctime.Before(forgingTimes.NextTime()) {
			logrus.WithField("proto", proto).Debugln("铁匠铺打造，倒计时未结束")
			result.Add(domestic.ERR_FORGING_EQUIP_FAIL_TIMES_NOT_ENOUGH)
			return
		}

		// 移除新标签
		heroDomestic.RemoveNewForgingPos(equipPos)

		forgingTimes.IncreseTimes()
		forgingTimes.SetNextTime(ctime.Add(tieJiangPuLevelData.RecoveryForgingDuration))

		// 发消息
		result.Add(domestic.NewS2cRecoveryForgingTimeChangeMsg(forgingTimes.TimesAndNextTime()))
		result.Add(domestic.NewS2cForgingEquipMsg(proto.GetSlot()))

		hctx := heromodule.NewContext(m.dep, operate_type.EquipmentForge)
		heromodule.AddEquipData(hctx, hero, result, equip, 1, ctime)

		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_CombineEquip)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILD_EQUIP)
		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILD_EQUIP_DAILY)

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_FORGE_EQUIPMENT) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_FORGE_EQUIPMENT)))
		}

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DomesticModule) ProcessWorkshopStartForge(proto *domestic.C2SStartWorkshopProto, hc iface.HeroController) {

	if proto.Index < 1 {
		logrus.Debug("装备作坊开始生产，无效的索引")
		hc.Send(domestic.ERR_START_WORKSHOP_FAIL_INVALID_INDEX)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.Domestic().GetWorkshopIndex() > 0 {
			logrus.Debug("装备作坊开始生产，超出锻造个数上限")
			result.Add(domestic.ERR_START_WORKSHOP_FAIL_COUNT_LIMIT)
			return
		}

		idx := int(proto.Index - 1)
		if idx >= len(hero.Domestic().GetWorkshopEquipment()) {
			logrus.Debug("装备作坊开始生产，index超出上限")
			result.Add(domestic.ERR_START_WORKSHOP_FAIL_INVALID_INDEX)
			return
		}

		equipmentData := hero.Domestic().GetWorkshopEquipment()[idx]
		if equipmentData == nil {
			logrus.Error("装备作坊开始生产，equipmentData == nil")
			result.Add(domestic.ERR_START_WORKSHOP_FAIL_INVALID_INDEX)
			return
		}

		d := time.Hour
		if vipData := m.dep.Datas().GetVipLevelData(hero.VipLevel()); vipData != nil && vipData.WorkshopAutoCompleted {
			d = 0
		} else if durationData := m.datas.GetWorkshopDuration(equipmentData.Id); durationData != nil {
			d = durationData.Duration
		}

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_FORGE_EQUIPMENT) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_FORGE_EQUIPMENT)))
			d = 3 * time.Second
		}

		ctime := m.timeService.CurrentTime()
		collectTime := ctime.Add(d)
		hero.Domestic().SetCurrentWorkshop(uint64(proto.Index), collectTime)

		result.Add(domestic.NewS2cStartWorkshopMsg(proto.Index, timeutil.Marshal32(collectTime)))

		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_CombineEquip)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILD_EQUIP)
		heromodule.IncreTaskProgressOne(hero, result, shared_proto.TaskTargetType_TASK_TARGET_BUILD_EQUIP_DAILY)

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DomesticModule) ProcessWorkshopMiaoCd(proto *domestic.C2SWorkshopMiaoCdProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticWorkshopMiaoCd)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		d := m.datas.MiscConfig().MiaoWorkshopDuration
		if d <= 0 {
			logrus.Debugln("锻造装备秒 cd，功能未开放")
			result.Add(domestic.ERR_WORKSHOP_MIAO_CD_FAIL_NOT_IN_CD)
			return
		}

		idx := hero.Domestic().GetWorkshopIndex()
		if idx <= 0 {
			logrus.Debugln("锻造装备秒 cd，当前没有装备在锻造")
			result.Add(domestic.ERR_WORKSHOP_MIAO_CD_FAIL_INVALID_INDEX)
			return
		}

		if u64.FromInt32(proto.Index) != idx {
			logrus.Debugln("锻造装备秒 cd，index装备没有在锻造")
			result.Add(domestic.ERR_WORKSHOP_MIAO_CD_FAIL_INVALID_INDEX)
			return
		}

		ctime := m.timeService.CurrentTime()
		endTime := hero.Domestic().GetWorkshopCollectTime()
		if endTime.Before(ctime) {
			logrus.Debugln("锻造装备秒 cd，装备已经锻造完成")
			result.Add(domestic.ERR_WORKSHOP_MIAO_CD_FAIL_NOT_IN_CD)
			return
		}

		multi := int64((endTime.Sub(ctime) + d - 1) / d)
		if multi <= 0 {
			logrus.Errorln("装备锻造秒cd，计算出来的multi <= 0")
			result.Add(domestic.ERR_WORKSHOP_MIAO_CD_FAIL_COST_NOT_ENOUGH)
			return
		}

		cost := m.datas.MiscConfig().MiaoWorkshopCost.Multiple(u64.FromInt64(multi))
		if !heromodule.TryReduceCost(hctx, hero, result, cost) {
			logrus.Debugln("装备锻造秒cd，资源不足")
			result.Add(domestic.ERR_WORKSHOP_MIAO_CD_FAIL_COST_NOT_ENOUGH)
			return
		}

		hero.Domestic().SetWorkshopCollectTime(ctime)
		result.Add(domestic.NewS2cWorkshopMiaoCdMsg(u64.Int32(idx)))

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DomesticModule) ProcessWorkshopCollect(proto *domestic.C2SCollectWorkshopProto, hc iface.HeroController) {

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if hero.Domestic().GetWorkshopIndex() <= 0 {
			logrus.Debug("领取作坊装备，当前没有装备在锻造")
			result.Add(domestic.ERR_COLLECT_WORKSHOP_FAIL_INVALID_INDEX)
			return
		}

		idx := hero.Domestic().GetWorkshopIndex()
		if proto.Index >= 1 && u64.FromInt32(proto.Index) != idx {
			logrus.Debug("领取作坊装备，index装备没有在锻造")
			result.Add(domestic.ERR_COLLECT_WORKSHOP_FAIL_INVALID_INDEX)
			return
		}

		ctime := m.timeService.CurrentTime()
		if vipData := m.datas.GetVipLevelData(hero.VipLevel()); vipData == nil || !vipData.WorkshopAutoCompleted {
			if ctime.Before(hero.Domestic().GetWorkshopCollectTime()) {
				logrus.Debug("领取作坊装备，装备还未锻造完成")
				result.Add(domestic.ERR_COLLECT_WORKSHOP_FAIL_CANT_COLLECT)
				return
			}
		}

		hctx := heromodule.NewContext(m.dep, operate_type.DomesticWorkshopCollect)
		equipmentData := hero.Domestic().CollectWorkshop()
		if equipmentData != nil {
			heromodule.AddEquipData(hctx, hero, result, equipmentData, 1, ctime)
		}

		result.Add(domestic.NewS2cCollectWorkshopMsg(proto.Index))

		var equipmentIds []int32
		var durations []int32
		for _, e := range hero.Domestic().GetWorkshopEquipment() {
			equipmentIds = append(equipmentIds, u64.Int32(e.Id))

			d := time.Hour
			if durationData := m.datas.GetWorkshopDuration(e.Id); durationData != nil {
				d = durationData.Duration
			}
			durations = append(durations, timeutil.DurationMarshal32(d))
		}

		result.Add(domestic.NewS2cListWorkshopEquipmentMsg(
			timeutil.Marshal32(hero.Domestic().GetNextRefreshWorkshopTime()),
			equipmentIds, durations, 0, 0,
			u64.Int32(hero.Domestic().GetWorkshopRefreshTimes()),
		))

		if hero.Bools().TrySet(shared_proto.HeroBoolType_BOOL_WORKSHOP_COLLECT) {
			result.Add(misc.NewS2cSetHeroBoolMsg(int32(shared_proto.HeroBoolType_BOOL_WORKSHOP_COLLECT)))
		}

		result.Ok()
	})

}

//gogen:iface c2s_refresh_workshop
func (m *DomesticModule) ProcessRefreshWorkshop(hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticRefreshWorkshop)

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDomestic := hero.Domestic()
		workshopRefreshTimes := heroDomestic.GetWorkshopRefreshTimes()
		// 获得已经刷新了的次数
		refreshCostData := m.datas.WorkshopRefreshCost().Get(workshopRefreshTimes + 1)
		if refreshCostData == nil {
			logrus.Debugf("没有刷新次数了")
			result.Add(domestic.ERR_REFRESH_WORKSHOP_FAIL_TIMES_NOT_ENOUGH)
			return
		}

		if !heroDomestic.CanWorkshopRefresh() {
			logrus.Debugf("还有装备没有锻造")
			result.Add(domestic.ERR_REFRESH_WORKSHOP_FAIL_HAS_EQUIP_NOT_FORG)
			return
		}

		if !heromodule.TryReduceCost(hctx, hero, result, refreshCostData.Cost) {
			logrus.Debugf("消耗不够")
			result.Add(domestic.ERR_REFRESH_WORKSHOP_FAIL_COST_NOT_ENOUGH)
			return
		}

		// 刷新

		result.Changed()
		result.Ok()

		heroDomestic.IncWorkshopRefreshTimes()
		heromodule.RefreshWorkshopAnyway(hctx, hero, result, m.datas, heroDomestic.GetNextRefreshWorkshopTime(), operate_type.RefreshMoney, u64.FromInt(refreshCostData.Cost.Id))

		result.Add(domestic.REFRESH_WORKSHOP_S2C)
		heromodule.SendWorkshopMsg(hero, result, m.datas)
	})

}

// 请求城内事件
//gogen:iface c2s_request_city_exchange_event
func (m *DomesticModule) ProcessRequestCityExchangeEvent(hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		data := m.datas.GetCityEventLevelData(hero.Level())
		if data == nil {
			logrus.Debugln("请求城内事件，玩家没解锁")
			result.Add(domestic.ERR_REQUEST_CITY_EXCHANGE_EVENT_FAIL_NOT_OPEN)
			return
		}

		ctime := m.timeService.CurrentTime()
		cityEvent := hero.Domestic().CityEvent()

		// 放宽一定的时间
		if ctime.Add(time.Second * 5).Before(cityEvent.CanExchangeTime()) {
			logrus.Debugln("请求城内事件，时间没到")
			result.Add(domestic.ERR_REQUEST_CITY_EXCHANGE_EVENT_FAIL_IN_CD)
			return
		}

		if cityEvent.AcceptTimes() >= m.datas.CityEventMiscData().MaxTimes {
			logrus.Debugln("请求城内事件，没有次数了")
			result.Add(domestic.ERR_REQUEST_CITY_EXCHANGE_EVENT_FAIL_NO_TIMES)
			return
		}

		if cityEvent.EventData() == nil {
			cityEvent.Accept(data.Random())
		}

		result.Add(domestic.NewS2cRequestCityExchangeEventMsg(u64.Int32(cityEvent.AcceptTimes()), u64.Int32(cityEvent.EventData().Id)))

		result.Changed()
		result.Ok()
	})
}

var exchangeMsg = domestic.NewS2cCityEventExchangeMsg(true).Static()
var giveUpMsg = domestic.NewS2cCityEventExchangeMsg(false).Static()

// 请求城内事件兑换
//gogen:iface c2s_city_event_exchange
func (m *DomesticModule) ProcessCityEventExchange(proto *domestic.C2SCityEventExchangeProto, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticCityEventExchange)

	var toAddProsperityFunc func()

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		ctime := m.timeService.CurrentTime()
		cityEvent := hero.Domestic().CityEvent()

		if ctime.Before(cityEvent.CanExchangeTime()) {
			logrus.Debugln("请求城内事件兑换，时间没到")
			result.Add(domestic.ERR_CITY_EVENT_EXCHANGE_FAIL_IN_CD)
			return
		}

		// 相等的时候，可能下面还有事件
		if cityEvent.AcceptTimes() > m.datas.CityEventMiscData().MaxTimes {
			logrus.Debugln("请求城内事件兑换，没有次数了")
			result.Add(domestic.ERR_CITY_EVENT_EXCHANGE_FAIL_NO_TIMES)
			return
		}

		data := cityEvent.EventData()
		if data == nil {
			logrus.Debugln("请求城内事件兑换，兑换数据没有")
			result.Add(domestic.ERR_CITY_EVENT_EXCHANGE_FAIL_CONDITION_NOT_SATISFIED)
			return
		}

		if !proto.GetGiveUp() {
			if data.Prize.Prosperity > 0 && hero.BaseRegion() == 0 {
				logrus.Debugln("请求城内事件兑换，流亡不能加繁荣度")
				result.Add(domestic.ERR_CITY_EVENT_EXCHANGE_FAIL_CONDITION_NOT_SATISFIED)
				return
			}

			// 检查消耗够不够，兑换
			if !heromodule.TryReduceCombineCost(hctx, hero, result, data.Cost, ctime) {
				logrus.Debugln("请求城内事件兑换，消耗不够")
				result.Add(domestic.ERR_CITY_EVENT_EXCHANGE_FAIL_CONDITION_NOT_SATISFIED)
				return
			}

			cityEvent.ExchangeOrGiveUp()

			hctx := heromodule.NewContext(m.dep, operate_type.DomesticCityEventExchange)
			heromodule.AddPrize(hctx, hero, result, data.Prize, ctime)

			if data.Prize.Prosperity > 0 {
				baseRegion := hero.BaseRegion()
				toAddProsperityFunc = func() {
					heromodule.AddProsperity(m.realmService, hc.Id(), baseRegion, data.Prize.Prosperity)
				}
			}
			result.Add(exchangeMsg)
		} else {
			cityEvent.ExchangeOrGiveUp()
			result.Add(giveUpMsg)
		}

		nextCanExchangeTime := cityEvent.SetCanExchangeTime(ctime.Add(m.datas.CityEventMiscData().RecoverDuration))
		result.Add(domestic.NewS2cCityEventTimeChangedMsg(timeutil.Marshal32(nextCanExchangeTime)))

		result.Changed()
		result.Ok()
	})

	if toAddProsperityFunc != nil {
		toAddProsperityFunc()
	}
}

// 签名
//gogen:iface
func (m *DomesticModule) ProcessSign(proto *domestic.C2SSignProto, hc iface.HeroController) {
	len := util.GetCharLen(proto.Text)

	if uint64(len) > m.datas.MiscConfig().MaxSignLen {
		logrus.WithField("proto", proto).Debugln("签名长度超出了")
		hc.Send(domestic.ERR_SIGN_FAIL_LEN_INVALID)
		return
	}

	if len > 0 {
		if !m.tssClient.TryCheckName("更新签名", hc, proto.Text, domestic.ERR_SIGN_FAIL_SENSITIVE_WORDS, domestic.ERR_SIGN_FAIL_SERVER_ERROR) {
			return
		}
	}

	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.Domestic().SetSign(proto.GetText())
		result.Add(domestic.SIGN_S2C)
		result.Changed()
		result.Ok()
	}) {
		return
	}

	m.realmService.GetBigMap().ChangeSign(hc.Id(), proto.GetText())
}

// 语音
//gogen:iface
func (m *DomesticModule) ProcessVoice(proto *domestic.C2SVoiceProto, hc iface.HeroController) {
	len := len(proto.GetContent())

	if uint64(len) > m.datas.MiscConfig().MaxVoiceLen {
		logrus.WithField("proto", proto).Debugln("语音长度超出了")
		hc.Send(domestic.ERR_VOICE_FAIL_LEN_INVALID)
		return
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		hero.Domestic().SetVoice(proto.GetContent())
		result.Add(domestic.VOICE_S2C)
		result.Changed()
		result.Ok()
	})
}

// 头像
//gogen:iface
func (m *DomesticModule) ProcessChangeHead(proto *domestic.C2SChangeHeadProto, hc iface.HeroController) {

	if !strings.HasPrefix(proto.HeadId, "http") {
		head := m.datas.GetHeadData(proto.HeadId)
		if head == nil {
			logrus.WithField("id", proto.HeadId).Debugf("未知的头像")
			hc.Send(domestic.ERR_CHANGE_HEAD_FAIL_INVALID_HEAD)
			return
		}

		if head.CountryOfficialType != shared_proto.CountryOfficialType_COT_NO_OFFICIAL {
			countryId := hc.LockHeroCountry()
			if m.dep.Country().HeroOfficial(countryId, hc.Id()) != head.CountryOfficialType {
				hc.Send(domestic.ERR_CHANGE_HEAD_FAIL_NOT_THIS_OFFICIAL)
				return
			}
		}
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		if strings.HasPrefix(proto.HeadId, "http") {
			hero.SetHeadUrl(proto.HeadId)
			result.Add(domestic.NewS2cChangeHeadMsg(proto.HeadId))
		} else {
			head := m.datas.GetHeadData(proto.HeadId)
			if head == nil {
				logrus.WithField("id", proto.HeadId).Debugf("未知的头像")
				hc.Send(domestic.ERR_CHANGE_HEAD_FAIL_INVALID_HEAD)
				return
			}

			if head.UnlockNeedHeroLevel != nil && hero.Level() < head.UnlockNeedHeroLevel.Level {
				logrus.WithField("id", proto.HeadId).Debugf("君主等级不足")
				hc.Send(domestic.ERR_CHANGE_HEAD_FAIL_HERO_LEVEL_TOO_LOW)
				return
			}

			if head.UnlockNeedCaptain != nil && hero.Military().Captain(head.UnlockNeedCaptain.Id) == nil {
				logrus.WithField("id", proto.HeadId).Debugf("武将未解锁")
				hc.Send(domestic.ERR_CHANGE_HEAD_FAIL_CAPTAIN_SOUL_NOT_UNLOCK)
				return
			}

			hero.SetHead(head)
			result.Add(head.ChangeHeadMsg)
		}

		result.Changed()
		result.Ok()
	})

}

// 形象
//gogen:iface
func (m *DomesticModule) ProcessChangeBody(proto *domestic.C2SChangeBodyProto, hc iface.HeroController) {
	body := m.datas.GetBodyData(u64.FromInt32(proto.BodyId))
	if body == nil {
		logrus.WithField("id", proto.BodyId).Debugf("未知的头像")
		hc.Send(domestic.ERR_CHANGE_BODY_FAIL_INVALID_BODY)
		return
	}

	if body.CountryOfficialType != shared_proto.CountryOfficialType_COT_NO_OFFICIAL {
		countryId := hc.LockHeroCountry()
		if m.dep.Country().HeroOfficial(countryId, hc.Id()) != body.CountryOfficialType {
			hc.Send(domestic.ERR_CHANGE_BODY_FAIL_NOT_THIS_OFFICIAL)
			return
		}
	}

	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {

		if body.UnlockNeedHeroLevel != nil && hero.Level() < body.UnlockNeedHeroLevel.Level {
			logrus.WithField("id", proto.BodyId).Debugf("君主等级不足")
			hc.Send(domestic.ERR_CHANGE_BODY_FAIL_HERO_LEVEL_TOO_LOW)
			return
		}

		if body.UnlockNeedCaptain != nil && hero.Military().Captain(body.UnlockNeedCaptain.Id) == nil {
			logrus.WithField("id", proto.BodyId).Debugf("武将未解锁")
			hc.Send(domestic.ERR_CHANGE_BODY_FAIL_CAPTAIN_SOUL_NOT_UNLOCK)
			return
		}

		hero.SetBody(body)
		result.Add(body.ChangeBodyMsg)

		result.Changed()
		result.Ok()
	})

}

//gogen:iface c2s_collect_countdown_prize
func (m *DomesticModule) ProcessCollectCountdownPrize(hc iface.HeroController) {

	// 领取倒计时奖励
	var toAddProsperityFunc func()
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDomestic := hero.Domestic()
		cp := heroDomestic.GetCountdownPrize()

		ctime := m.timeService.CurrentTime()
		if ctime.Before(cp.CollectTime()) {
			logrus.Debug("领取倒计时奖励，时间未到")
			result.Add(domestic.ERR_COLLECT_COUNTDOWN_PRIZE_FAIL_TIME_NOT_REACHED)
			return
		}

		if hero.BaseLevel() != 0 {
			// cd到没到
			baseLevelData := m.datas.BaseLevelData().Must(hero.BaseLevel())
			if hero.Prosperity() < hero.ProsperityCapcity() &&
				ctime.Before(hero.GetMoveBaseRestoreProsperityBufEndTime()) &&
				heroDomestic.CollectProsperityDownCountDownPrizeTimes() < m.datas.MiscConfig().GiveAddCountDownPrizeTimes &&
				ctime.After(heroDomestic.NextCanCollectProsperityDownCountDownPrizeTime()) {

				heroDomestic.IncCollectProsperityDownCountDownPrizeTimes()
				heroDomestic.SetNextCanCollectProsperityDownCountDownPrizeTime(ctime.Add(m.datas.MiscConfig().GiveAddCountDownPrizeDuration))
				toAddProsperity := baseLevelData.AddCountdownPrizeProsperity
				if toAddProsperity > 0 {
					// 如果已经加过繁荣度，本次奖励顺延
					cp.NextCollectTime(ctime)

					result.Add(domestic.NewS2cCollectCountdownPrizeMsg(
						nil, timeutil.Marshal32(cp.CollectTime()), u64.Int32(baseLevelData.AddCountdownPrizeDesc.Id), u64.Int32(toAddProsperity),
					))

					baseRegion := hero.BaseRegion()
					toAddProsperityFunc = func() {
						heromodule.AddProsperity(m.realmService, hc.Id(), baseRegion, toAddProsperity)
					}
					return
				}

			}
		}

		// 没加过奖励
		toAddPrize := cp.Prize()
		descId := cp.Desc().Id

		newPrize := heroDomestic.CollectCountdownPrize(ctime)

		result.Add(domestic.NewS2cCollectCountdownPrizeMarshalMsg(
			toAddPrize.PrizeProto(),
			timeutil.Marshal32(newPrize.CollectTime()),
			u64.Int32(descId), 0,
		))

		hctx := heromodule.NewContext(m.dep, operate_type.DomesticCollectCountdownPrize)
		heromodule.AddPrize(hctx, hero, result, toAddPrize, ctime)

		hero.HistoryAmount().IncreaseOne(server_proto.HistoryAmountType_AccumCountDownPrizeTimes)
		heromodule.UpdateTaskProgress(hero, result, shared_proto.TaskTargetType_TASK_TARGET_ACCUM_COUNT_DOWN_PRIZE)
	})

	if toAddProsperityFunc != nil {
		toAddProsperityFunc()
	}

}

//gogen:iface c2s_collect_season_prize
func (m *DomesticModule) ProcessCollectSeasonPrize(hc iface.HeroController) {

	// 领取倒计时奖励
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		heroDomestic := hero.Domestic()

		if heroDomestic.IsCollectSeasonPrize() {
			logrus.Debug("领取季节奖励，奖励已经领取")
			result.Add(domestic.ERR_COLLECT_SEASON_PRIZE_FAIL_COLLECTED)
			return
		}

		heroDomestic.SetCollectSeasonPrize()

		seasonData := m.seasonService.Season()

		ctime := m.timeService.CurrentTime()
		hctx := heromodule.NewContext(m.dep, operate_type.DomesticCollectSeasonPrize)
		heromodule.AddPrize(hctx, hero, result, seasonData.Prize, ctime)

		result.Add(domestic.COLLECT_SEASON_PRIZE_S2C)

		result.Changed()
		result.Ok()
	})
}

//gogen:iface
func (m *DomesticModule) ProcessBuySp(proto *domestic.C2SBuySpProto, hc iface.HeroController) {
	buyTimes := proto.GetBuyTimes()
	if buyTimes <= 0 {
		logrus.Debug("购买体力值，错误数据")
		hc.Send(domestic.ERR_BUY_SP_FAIL_INVALID_DATA)
		return
	}
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		times := u64.FromInt32(buyTimes)
		miscConfig := m.datas.MiscGenConfig()
		var vipExtraTimes uint64
		if vipData := m.datas.GetVipLevelData(hero.VipLevel()); vipData != nil {
			vipExtraTimes = vipData.BuySpMaxTimes
		}
		if hero.GetBuySpTimes()+times > miscConfig.BuySpLimit+vipExtraTimes {
			logrus.Debug("购买体力值，购买次数上限")
			result.Add(domestic.ERR_BUY_SP_FAIL_BUY_TIMES_LIMIT)
			return
		}
		cost := miscConfig.BuySpCost * times
		if !hero.HasEnoughYuanbao(cost) {
			logrus.Debug("购买体力值，元宝不足")
			result.Add(domestic.ERR_BUY_SP_FAIL_NOT_ENOUGH_YUANBAO)
			return
		}
		hctx := heromodule.NewContext(m.dep, operate_type.DomesticBuySp)
		// 扣元宝
		heromodule.ReduceYuanbaoAnyway(hctx, hero, result, cost)
		// 累计次数
		hero.AddBuySpTimes(times)
		// 增加体力值
		hero.AddSp(miscConfig.BuySpValue * times)

		result.Add(domestic.NewS2cBuySpMsg(u64.Int32(hero.GetSp()), u64.Int32(hero.GetBuySpTimes())))

		result.Changed()
		result.Ok()
	})
}

// 使用主城增益
//gogen:iface
func (m *DomesticModule) ProcessUseAdvantage(proto *domestic.C2SUseAdvantageProto, hc iface.HeroController) {
	id := u64.FromInt32(proto.GetId())
	d := m.datas.GetBufferData(id)
	if d == nil {
		logrus.Debug("使用主城增益，无该增益")
		hc.Send(domestic.ERR_USE_ADVANTAGE_FAIL_INVALID_ID)
		return
	}

	if d.TypeData.IsMian {
		m.UseMianGoods(d, id, hc)
	} else {
		m.UseBuffGoods(d, id, hc)
	}
}

func (m *DomesticModule) UseMianGoods(d *buffer.BufferData, id uint64, hc iface.HeroController) {
	ctime := m.dep.Time().CurrentTime()
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		// 增益条件验证，这里用野外免战错误消息
		if hero.Depot().GetGoodsCount(d.BuffGoodsData.Id) <= 0 {
			logrus.Debug("主城增益，使用免战物品，物品个数不足")
			result.Add(region.ERR_USE_MIAN_GOODS_FAIL_COUNT_NOT_ENOUGH)
			return
		}

		// 免战中
		if ctime.Add(d.BuffGoodsData.GoodsEffect.MianDuration).Before(hero.GetMianDisappearTime()) {
			logrus.Debug("主城增益，使用免战物品，免战中")
			result.Add(region.ERR_USE_MIAN_GOODS_FAIL_MIAN)
			return
		}

		// CD
		if ctime.Before(hero.GetNextUseMianGoodsTime()) {
			logrus.Debug("主城增益，使用免战物品，免战物品CD中")
			result.Add(region.ERR_USE_MIAN_GOODS_FAIL_COOLDOWN)
			return
		}

		if hero.BaseRegion() == 0 || hero.BaseLevel() <= 0 || hero.Prosperity() <= 0 {
			logrus.Debug("主城增益，使用免战物品，主城流亡了")
			result.Add(region.ERR_USE_MIAN_GOODS_FAIL_HOME_NOT_ALIVE)
			return
		}

		result.Ok()
	}) {
		return
	}

	if m.region.UseMianGoods(d.BuffGoodsId, false, hc) {
		// 增益成功
		hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
			startTime := hero.GetMianStartTime()
			endTime := hero.GetMianDisappearTime()
			result.Add(domestic.NewS2cUseAdvantageMsg(u64.Int32(id), timeutil.Marshal32(startTime), timeutil.Marshal32(endTime)))
			result.Ok()
		})
	} else {
		hc.Send(domestic.ERR_USE_ADVANTAGE_FAIL_ITEM_NOT_ENOUGH)
	}
}

func (m *DomesticModule) UseBuffGoods(d *buffer.BufferData, id uint64, hc iface.HeroController) {
	hctx := heromodule.NewContext(m.dep, operate_type.DomesticUseBufEffect)
	ctime := m.dep.Time().CurrentTime()

	var reserveResult *entity.ReserveResult
	if hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		// 增益条件验证
		currBuff := hero.Buff().Buff(d.TypeData.BuffGroup)
		if currBuff != nil {
			newBuffData := m.datas.GetBuffEffectData(d.BuffGoodsData.GoodsEffect.BuffId)
			if newBuffData == nil {
				logrus.Errorf("找不到增益：%v 物品：%v 对应的 buff:%v", id, d.BuffGoodsData.Id, d.BuffGoodsData.GoodsEffect.BuffId)
				result.Add(domestic.ERR_USE_ADVANTAGE_FAIL_BUFF_EFFECT_FAIL)
				return
			}
			if newBuffData.Level < currBuff.EffectData.Level {
				result.Add(domestic.ERR_USE_ADVANTAGE_FAIL_LEVEL_LIMIT)
				return
			} else if newBuffData.Level == currBuff.EffectData.Level {
				if newBuffData.KeepDuration <= currBuff.EffectData.KeepDuration {
					result.Add(domestic.ERR_USE_ADVANTAGE_FAIL_KEEP_TIME_LIMIT)
					return
				}
			}
		}

		// 扣物品
		var ok bool
		if ok, reserveResult = heromodule.ReserveGoods(hctx, hero, result, d.BuffGoodsData, 1, ctime); !ok {
			result.Add(domestic.ERR_USE_ADVANTAGE_FAIL_ITEM_NOT_ENOUGH)
			return
		}
		result.Changed()
		result.Ok()
	}) {
		return
	}

	var buffUsedSucc bool
	heromodule.ConfirmReserveResult(hctx, hc, reserveResult, func() (success bool) {
		// 使用 buff
		buffData := m.dep.Datas().BuffEffectData().Get(d.BuffGoodsData.GoodsEffect.BuffId)
		if success = m.buffService.AddBuffToSelf(buffData, hc.Id()); !success {
			logrus.Warnf("主城增益使用 buff，没有成功。heroId:%v buffId:%v", hc.Id(), d.BuffGoodsData.GoodsEffect.BuffId)
			return
		}
		buffUsedSucc = true
		return
	})

	if !buffUsedSucc {
		hc.Send(domestic.ERR_USE_ADVANTAGE_FAIL_BUFF_EFFECT_FAIL)
		return
	}

	// 增益成功
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		var startTime, endTime time.Time
		if newBuff := hero.Buff().Buff(d.TypeData.BuffGroup); newBuff != nil {
			startTime = newBuff.StartTime
			endTime = newBuff.EndTime
		}
		result.Add(domestic.NewS2cUseAdvantageMsg(u64.Int32(id), timeutil.Marshal32(startTime), timeutil.Marshal32(endTime)))
		result.Ok()
	})
}

// 废弃
//gogen:iface
func (m *DomesticModule) ProcessUseBufEffect(proto *domestic.C2SUseBufEffectProto, hc iface.HeroController) {
	hc.Send(domestic.ERR_USE_BUF_EFFECT_FAIL_SERVER_ERROR)
}

// 废弃
//gogen:iface c2s_open_buf_effect_ui
func (m *DomesticModule) ProcessOpenBufEffectUi(hc iface.HeroController) {
	hc.Send(domestic.NewS2cOpenBufEffectUiMsg(&shared_proto.HeroBufferProto{}))
}

//gogen:iface
func (m *DomesticModule) ProcessWorkerUnlock(proto *domestic.C2SWorkerUnlockProto, hc iface.HeroController) {
	hc.FuncWithSend(func(hero *entity.Hero, result herolock.LockResult) {
		ctime := m.dep.Time().CurrentTime()
		dome := hero.Domestic()
		pos := int(proto.Pos)

		if !dome.WorkerPosLegal(pos) {
			result.Add(domestic.ERR_WORKER_UNLOCK_FAIL_INVALID_POS)
			return
		}

		if dome.WorkerIsUnlock(pos, ctime) {
			result.Add(domestic.ERR_WORKER_UNLOCK_FAIL_POS_IS_UNLOCKED)
			return
		}

		if !dome.WorkerNeverUnlocked(pos) {
			hctx := heromodule.NewContext(m.dep, operate_type.DomesticWorkerUnlock)
			if !heromodule.TryReduceCost(hctx, hero, result, m.dep.Datas().MiscConfig().SecondWorkerCost) {
				result.Add(domestic.ERR_WORKER_UNLOCK_FAIL_COST_NOT_ENOUGH)
				return
			}
		}

		newTime := dome.UnlockWorker(pos, ctime.Add(m.dep.Datas().MiscConfig().SecondWorkerUnlockDuration))
		result.Add(domestic.NewS2cWorkerUnlockMsg(int32(pos), timeutil.Marshal32(newTime)))
		result.Changed()
		result.Ok()
	})
}
