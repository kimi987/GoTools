package domestic

import (
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/pbutil"
)

var (
	pool           = pbutil.Pool
	newProtoMsg    = util.NewProtoMsg
	newCompressMsg = util.NewCompressMsg
	safeMarshal    = util.SafeMarshal
	_              = shared_proto.ErrIntOverflowConfig
)

type marshaler util.Marshaler

const (
	MODULE_ID = 2

	C2S_CREATE_BUILDING = 1

	C2S_UPGRADE_BUILDING = 4

	C2S_REBUILD_RESOURCE_BUILDING = 7

	C2S_UNLOCK_OUTER_CITY = 108

	C2S_UPDATE_OUTER_CITY_TYPE = 142

	C2S_UPGRADE_OUTER_CITY_BUILDING = 111

	C2S_COLLECT_RESOURCE = 15

	C2S_COLLECT_RESOURCE_V2 = 76

	C2S_REQUEST_RESOURCE_CONFLICT = 81

	C2S_LEARN_TECHNOLOGY = 18

	C2S_UNLOCK_STABLE_BUILDING = 87

	C2S_UPGRADE_STABLE_BUILDING = 24

	C2S_IS_HERO_NAME_EXIST = 90

	C2S_CHANGE_HERO_NAME = 30

	C2S_LIST_OLD_NAME = 33

	C2S_VIEW_OTHER_HERO = 35

	C2S_VIEW_FIGHT_INFO = 125

	C2S_MIAO_BUILDING_WORKER_CD = 41

	C2S_MIAO_TECH_WORKER_CD = 44

	C2S_FORGING_EQUIP = 51

	C2S_SIGN = 66

	C2S_VOICE = 69

	C2S_REQUEST_CITY_EXCHANGE_EVENT = 60

	C2S_CITY_EVENT_EXCHANGE = 63

	C2S_CHANGE_HEAD = 94

	C2S_CHANGE_BODY = 130

	C2S_COLLECT_COUNTDOWN_PRIZE = 114

	C2S_START_WORKSHOP = 119

	C2S_COLLECT_WORKSHOP = 122

	C2S_WORKSHOP_MIAO_CD = 127

	C2S_REFRESH_WORKSHOP = 133

	C2S_COLLECT_SEASON_PRIZE = 136

	C2S_BUY_SP = 147

	C2S_USE_BUF_EFFECT = 150

	C2S_OPEN_BUF_EFFECT_UI = 154

	C2S_USE_ADVANTAGE = 158

	C2S_WORKER_UNLOCK = 166
)

func NewS2cUpdateResourceBuildingMsg(id int32, amount int32, capcity int32, output int32, conflict bool, base_level_lock bool) pbutil.Buffer {
	msg := &S2CUpdateResourceBuildingProto{
		Id:            id,
		Amount:        amount,
		Capcity:       capcity,
		Output:        output,
		Conflict:      conflict,
		BaseLevelLock: base_level_lock,
	}
	return NewS2cUpdateResourceBuildingProtoMsg(msg)
}

var s2c_update_resource_building = [...]byte{2, 74} // 74
func NewS2cUpdateResourceBuildingProtoMsg(object *S2CUpdateResourceBuildingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_resource_building[:], "s2c_update_resource_building")

}

func NewS2cUpdateMultiResourceBuildingMsg(id []int32, amount []int32, capcity []int32, output []int32, conflict []bool, base_level_lock []bool) pbutil.Buffer {
	msg := &S2CUpdateMultiResourceBuildingProto{
		Id:            id,
		Amount:        amount,
		Capcity:       capcity,
		Output:        output,
		Conflict:      conflict,
		BaseLevelLock: base_level_lock,
	}
	return NewS2cUpdateMultiResourceBuildingProtoMsg(msg)
}

var s2c_update_multi_resource_building = [...]byte{2, 75} // 75
func NewS2cUpdateMultiResourceBuildingProtoMsg(object *S2CUpdateMultiResourceBuildingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_multi_resource_building[:], "s2c_update_multi_resource_building")

}

func NewS2cCreateBuildingMsg(id int32, building int32, worker_pos int32, worker_rest_end_time int32) pbutil.Buffer {
	msg := &S2CCreateBuildingProto{
		Id:                id,
		Building:          building,
		WorkerPos:         worker_pos,
		WorkerRestEndTime: worker_rest_end_time,
	}
	return NewS2cCreateBuildingProtoMsg(msg)
}

var s2c_create_building = [...]byte{2, 2} // 2
func NewS2cCreateBuildingProtoMsg(object *S2CCreateBuildingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_create_building[:], "s2c_create_building")

}

// 无效的布局id
var ERR_CREATE_BUILDING_FAIL_INVALID_LAYOUT = pbutil.StaticBuffer{3, 2, 3, 1} // 3-1

// 这个id上已经有建筑了
var ERR_CREATE_BUILDING_FAIL_NOT_EMPTY = pbutil.StaticBuffer{3, 2, 3, 2} // 3-2

// 建筑队在休息
var ERR_CREATE_BUILDING_FAIL_WORKER_REST = pbutil.StaticBuffer{3, 2, 3, 3} // 3-3

// 消耗资源不足
var ERR_CREATE_BUILDING_FAIL_RESOURCE_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 3, 4} // 3-4

// 前提条件未达成
var ERR_CREATE_BUILDING_FAIL_REQUIRE_NOT_REACH = pbutil.StaticBuffer{3, 2, 3, 5} // 3-5

// 资源点处于冲突状态
var ERR_CREATE_BUILDING_FAIL_RESOURCE_CONFLICT = pbutil.StaticBuffer{3, 2, 3, 6} // 3-6

// 资源点不在你的势力范围，请先升级主城
var ERR_CREATE_BUILDING_FAIL_RESOURCE_INVALID = pbutil.StaticBuffer{3, 2, 3, 7} // 3-7

// 无效的建筑类型
var ERR_CREATE_BUILDING_FAIL_INVALID_TYPE = pbutil.StaticBuffer{3, 2, 3, 8} // 3-8

// 服务器忙，请稍后再试
var ERR_CREATE_BUILDING_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 3, 9} // 3-9

func NewS2cUpgradeBuildingMsg(id int32, building int32, worker_pos int32, worker_rest_end_time int32) pbutil.Buffer {
	msg := &S2CUpgradeBuildingProto{
		Id:                id,
		Building:          building,
		WorkerPos:         worker_pos,
		WorkerRestEndTime: worker_rest_end_time,
	}
	return NewS2cUpgradeBuildingProtoMsg(msg)
}

var s2c_upgrade_building = [...]byte{2, 5} // 5
func NewS2cUpgradeBuildingProtoMsg(object *S2CUpgradeBuildingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_upgrade_building[:], "s2c_upgrade_building")

}

// 无效的布局id
var ERR_UPGRADE_BUILDING_FAIL_INVALID_LAYOUT = pbutil.StaticBuffer{3, 2, 6, 1} // 6-1

// 这个位置没有建筑
var ERR_UPGRADE_BUILDING_FAIL_NOT_BUILDING = pbutil.StaticBuffer{3, 2, 6, 2} // 6-2

// 这个建筑已经最高级了
var ERR_UPGRADE_BUILDING_FAIL_MAX_LEVEL = pbutil.StaticBuffer{3, 2, 6, 3} // 6-3

// 建筑队在休息
var ERR_UPGRADE_BUILDING_FAIL_WORKER_REST = pbutil.StaticBuffer{3, 2, 6, 4} // 6-4

// 消耗资源不足
var ERR_UPGRADE_BUILDING_FAIL_RESOURCE_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 6, 5} // 6-5

// 前提条件未达成
var ERR_UPGRADE_BUILDING_FAIL_REQUIRE_NOT_REACH = pbutil.StaticBuffer{3, 2, 6, 6} // 6-6

// 资源点处于冲突状态
var ERR_UPGRADE_BUILDING_FAIL_RESOURCE_CONFLICT = pbutil.StaticBuffer{3, 2, 6, 7} // 6-7

// 资源点不在你的势力范围，请先升级主城
var ERR_UPGRADE_BUILDING_FAIL_RESOURCE_INVALID = pbutil.StaticBuffer{3, 2, 6, 8} // 6-8

// 服务器忙，请稍后再试
var ERR_UPGRADE_BUILDING_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 6, 9} // 6-9

func NewS2cRebuildResourceBuildingMsg(id int32, building int32, worker_pos int32, worker_rest_end_time int32) pbutil.Buffer {
	msg := &S2CRebuildResourceBuildingProto{
		Id:                id,
		Building:          building,
		WorkerPos:         worker_pos,
		WorkerRestEndTime: worker_rest_end_time,
	}
	return NewS2cRebuildResourceBuildingProtoMsg(msg)
}

var s2c_rebuild_resource_building = [...]byte{2, 8} // 8
func NewS2cRebuildResourceBuildingProtoMsg(object *S2CRebuildResourceBuildingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_rebuild_resource_building[:], "s2c_rebuild_resource_building")

}

// 无效的布局id
var ERR_REBUILD_RESOURCE_BUILDING_FAIL_INVALID_LAYOUT = pbutil.StaticBuffer{3, 2, 9, 1} // 9-1

// 此建筑不能改建
var ERR_REBUILD_RESOURCE_BUILDING_FAIL_INVALID_BUILDING = pbutil.StaticBuffer{3, 2, 9, 2} // 9-2

// 建筑队在休息
var ERR_REBUILD_RESOURCE_BUILDING_FAIL_WORKER_REST = pbutil.StaticBuffer{3, 2, 9, 3} // 9-3

// 消耗资源不足
var ERR_REBUILD_RESOURCE_BUILDING_FAIL_RESOURCE_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 9, 4} // 9-4

// 资源点处于冲突状态
var ERR_REBUILD_RESOURCE_BUILDING_FAIL_RESOURCE_CONFLICT = pbutil.StaticBuffer{3, 2, 9, 6} // 9-6

// 资源点不在你的势力范围，请先升级主城
var ERR_REBUILD_RESOURCE_BUILDING_FAIL_RESOURCE_INVALID = pbutil.StaticBuffer{3, 2, 9, 7} // 9-7

// 服务器忙，请稍后再试
var ERR_REBUILD_RESOURCE_BUILDING_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 9, 5} // 9-5

func NewS2cUnlockOuterCityMsg(outer_city []byte) pbutil.Buffer {
	msg := &S2CUnlockOuterCityProto{
		OuterCity: outer_city,
	}
	return NewS2cUnlockOuterCityProtoMsg(msg)
}

var s2c_unlock_outer_city = [...]byte{2, 109} // 109
func NewS2cUnlockOuterCityProtoMsg(object *S2CUnlockOuterCityProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_unlock_outer_city[:], "s2c_unlock_outer_city")

}

// 无效的外城id
var ERR_UNLOCK_OUTER_CITY_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 2, 110, 1} // 110-1

// 该外城已经解锁了
var ERR_UNLOCK_OUTER_CITY_FAIL_UNLOCKED = pbutil.StaticBuffer{3, 2, 110, 2} // 110-2

// 官府等级不足
var ERR_UNLOCK_OUTER_CITY_FAIL_REQUIRE_NOT_REACH = pbutil.StaticBuffer{3, 2, 110, 3} // 110-3

// 服务器忙，请稍后再试
var ERR_UNLOCK_OUTER_CITY_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 110, 4} // 110-4

func NewS2cUpdateOuterCityTypeMsg(id int32, t int32, ids []int32) pbutil.Buffer {
	msg := &S2CUpdateOuterCityTypeProto{
		Id:  id,
		T:   t,
		Ids: ids,
	}
	return NewS2cUpdateOuterCityTypeProtoMsg(msg)
}

var s2c_update_outer_city_type = [...]byte{2, 143, 1} // 143
func NewS2cUpdateOuterCityTypeProtoMsg(object *S2CUpdateOuterCityTypeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_outer_city_type[:], "s2c_update_outer_city_type")

}

// 无效的外城id
var ERR_UPDATE_OUTER_CITY_TYPE_FAIL_INVALID_ID = pbutil.StaticBuffer{4, 2, 144, 1, 1} // 144-1

// 该外城还未解锁
var ERR_UPDATE_OUTER_CITY_TYPE_FAIL_LOCKED = pbutil.StaticBuffer{4, 2, 144, 1, 2} // 144-2

// 已经是该类型
var ERR_UPDATE_OUTER_CITY_TYPE_FAIL_SAME_TYPE = pbutil.StaticBuffer{4, 2, 144, 1, 5} // 144-5

// 外城改建消耗不足
var ERR_UPDATE_OUTER_CITY_TYPE_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 2, 144, 1, 3} // 144-3

// 服务器忙，请稍后再试
var ERR_UPDATE_OUTER_CITY_TYPE_FAIL_SERVER_BUSY = pbutil.StaticBuffer{4, 2, 144, 1, 4} // 144-4

func NewS2cUpgradeOuterCityBuildingMsg(city_id int32, old_id int32, id int32) pbutil.Buffer {
	msg := &S2CUpgradeOuterCityBuildingProto{
		CityId: city_id,
		OldId:  old_id,
		Id:     id,
	}
	return NewS2cUpgradeOuterCityBuildingProtoMsg(msg)
}

var s2c_upgrade_outer_city_building = [...]byte{2, 112} // 112
func NewS2cUpgradeOuterCityBuildingProtoMsg(object *S2CUpgradeOuterCityBuildingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_upgrade_outer_city_building[:], "s2c_upgrade_outer_city_building")

}

// 无效的外城建筑id
var ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_INVALID_BUILDING_ID = pbutil.StaticBuffer{3, 2, 113, 1} // 113-1

// 外城未解锁
var ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_CITY_LOCKED = pbutil.StaticBuffer{3, 2, 113, 2} // 113-2

// 前置建筑未满足升级条件
var ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_REQUIRE_NOT_REACH = pbutil.StaticBuffer{3, 2, 113, 3} // 113-3

// 消耗不足
var ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 113, 4} // 113-4

// 服务器忙，请稍后再试
var ERR_UPGRADE_OUTER_CITY_BUILDING_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 113, 5} // 113-5

func NewS2cResourceUpdateMsg(gold int32, food int32, wood int32, stone int32, is_safe bool) pbutil.Buffer {
	msg := &S2CResourceUpdateProto{
		Gold:   gold,
		Food:   food,
		Wood:   wood,
		Stone:  stone,
		IsSafe: is_safe,
	}
	return NewS2cResourceUpdateProtoMsg(msg)
}

var s2c_resource_update = [...]byte{2, 13} // 13
func NewS2cResourceUpdateProtoMsg(object *S2CResourceUpdateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_resource_update[:], "s2c_resource_update")

}

func NewS2cResourceUpdateSingleMsg(res_type int32, amount int32, is_safe bool) pbutil.Buffer {
	msg := &S2CResourceUpdateSingleProto{
		ResType: res_type,
		Amount:  amount,
		IsSafe:  is_safe,
	}
	return NewS2cResourceUpdateSingleProtoMsg(msg)
}

var s2c_resource_update_single = [...]byte{2, 28} // 28
func NewS2cResourceUpdateSingleProtoMsg(object *S2CResourceUpdateSingleProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_resource_update_single[:], "s2c_resource_update_single")

}

func NewS2cResourceCapcityUpdateMsg(gold_capcity int32, food_capcity int32, wood_capcity int32, stone_capcity int32, protected_capcity int32) pbutil.Buffer {
	msg := &S2CResourceCapcityUpdateProto{
		GoldCapcity:      gold_capcity,
		FoodCapcity:      food_capcity,
		WoodCapcity:      wood_capcity,
		StoneCapcity:     stone_capcity,
		ProtectedCapcity: protected_capcity,
	}
	return NewS2cResourceCapcityUpdateProtoMsg(msg)
}

var s2c_resource_capcity_update = [...]byte{2, 14} // 14
func NewS2cResourceCapcityUpdateProtoMsg(object *S2CResourceCapcityUpdateProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_resource_capcity_update[:], "s2c_resource_capcity_update")

}

func NewS2cCollectResourceMsg(id int32, amount int32) pbutil.Buffer {
	msg := &S2CCollectResourceProto{
		Id:     id,
		Amount: amount,
	}
	return NewS2cCollectResourceProtoMsg(msg)
}

var s2c_collect_resource = [...]byte{2, 16} // 16
func NewS2cCollectResourceProtoMsg(object *S2CCollectResourceProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_resource[:], "s2c_collect_resource")

}

// 无效的布局id
var ERR_COLLECT_RESOURCE_FAIL_INVALID_LAYOUT = pbutil.StaticBuffer{3, 2, 17, 1} // 17-1

// 此建筑不是资源点
var ERR_COLLECT_RESOURCE_FAIL_CANT_COLLECTED = pbutil.StaticBuffer{3, 2, 17, 2} // 17-2

// 资源点处于冲突状态
var ERR_COLLECT_RESOURCE_FAIL_RESOURCE_CONFLICT = pbutil.StaticBuffer{3, 2, 17, 3} // 17-3

// 资源点不在你的势力范围，请先升级主城
var ERR_COLLECT_RESOURCE_FAIL_RESOURCE_INVALID = pbutil.StaticBuffer{3, 2, 17, 4} // 17-4

// 仓库已满
var ERR_COLLECT_RESOURCE_FAIL_FULL = pbutil.StaticBuffer{3, 2, 17, 5} // 17-5

// 没有资源可以采集
var ERR_COLLECT_RESOURCE_FAIL_EMPTY = pbutil.StaticBuffer{3, 2, 17, 6} // 17-6

// 服务器忙，请稍后再试
var ERR_COLLECT_RESOURCE_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 17, 7} // 17-7

func NewS2cCollectResourceV2Msg(res_type int32, amount int32, collect_times int32, next_collect_time int32) pbutil.Buffer {
	msg := &S2CCollectResourceV2Proto{
		ResType:         res_type,
		Amount:          amount,
		CollectTimes:    collect_times,
		NextCollectTime: next_collect_time,
	}
	return NewS2cCollectResourceV2ProtoMsg(msg)
}

var s2c_collect_resource_v2 = [...]byte{2, 77} // 77
func NewS2cCollectResourceV2ProtoMsg(object *S2CCollectResourceV2Proto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_resource_v2[:], "s2c_collect_resource_v2")

}

// 未知的资源类型
var ERR_COLLECT_RESOURCE_V2_FAIL_INVALID_RESOURCE_TYPE = pbutil.StaticBuffer{3, 2, 78, 1} // 78-1

// 没有资源可以采集
var ERR_COLLECT_RESOURCE_V2_FAIL_EMPTY = pbutil.StaticBuffer{3, 2, 78, 2} // 78-2

// 没有次数
var ERR_COLLECT_RESOURCE_V2_FAIL_NO_TIMES = pbutil.StaticBuffer{3, 2, 78, 3} // 78-3

// 仓库已满
var ERR_COLLECT_RESOURCE_V2_FAIL_FULL = pbutil.StaticBuffer{3, 2, 78, 4} // 78-4

// 倒计时未结束
var ERR_COLLECT_RESOURCE_V2_FAIL_COUNTDOWN = pbutil.StaticBuffer{3, 2, 78, 6} // 78-6

// 服务器忙，请稍后再试
var ERR_COLLECT_RESOURCE_V2_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 78, 5} // 78-5

func NewS2cCollectResourceTimesChangedMsg(start_recover_collect_time int32) pbutil.Buffer {
	msg := &S2CCollectResourceTimesChangedProto{
		StartRecoverCollectTime: start_recover_collect_time,
	}
	return NewS2cCollectResourceTimesChangedProtoMsg(msg)
}

var s2c_collect_resource_times_changed = [...]byte{2, 79} // 79
func NewS2cCollectResourceTimesChangedProtoMsg(object *S2CCollectResourceTimesChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_resource_times_changed[:], "s2c_collect_resource_times_changed")

}

func NewS2cResourcePointChangeV2Msg(data []byte) pbutil.Buffer {
	msg := &S2CResourcePointChangeV2Proto{
		Data: data,
	}
	return NewS2cResourcePointChangeV2ProtoMsg(msg)
}

var s2c_resource_point_change_v2 = [...]byte{2, 80} // 80
func NewS2cResourcePointChangeV2ProtoMsg(object *S2CResourcePointChangeV2Proto) pbutil.Buffer {

	return newProtoMsg(object, s2c_resource_point_change_v2[:], "s2c_resource_point_change_v2")

}

func NewS2cRequestResourceConflictMsg(flag []string, name []string) pbutil.Buffer {
	msg := &S2CRequestResourceConflictProto{
		Flag: flag,
		Name: name,
	}
	return NewS2cRequestResourceConflictProtoMsg(msg)
}

var s2c_request_resource_conflict = [...]byte{2, 82} // 82
func NewS2cRequestResourceConflictProtoMsg(object *S2CRequestResourceConflictProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_resource_conflict[:], "s2c_request_resource_conflict")

}

// 服务器忙，请稍后再试
var ERR_REQUEST_RESOURCE_CONFLICT_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 83, 1} // 83-1

var RESOURCE_CONFLICT_CHANGED_S2C = pbutil.StaticBuffer{2, 2, 84} // 84

func NewS2cLearnTechnologyMsg(id int32, worker_pos int32, worker_rest_end_time int32) pbutil.Buffer {
	msg := &S2CLearnTechnologyProto{
		Id:                id,
		WorkerPos:         worker_pos,
		WorkerRestEndTime: worker_rest_end_time,
	}
	return NewS2cLearnTechnologyProtoMsg(msg)
}

var s2c_learn_technology = [...]byte{2, 19} // 19
func NewS2cLearnTechnologyProtoMsg(object *S2CLearnTechnologyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_learn_technology[:], "s2c_learn_technology")

}

// 无效的科技id
var ERR_LEARN_TECHNOLOGY_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 2, 20, 1} // 20-1

// 这个科技已经学会了
var ERR_LEARN_TECHNOLOGY_FAIL_LEARNED = pbutil.StaticBuffer{3, 2, 20, 2} // 20-2

// 消耗资源不足
var ERR_LEARN_TECHNOLOGY_FAIL_RESOURCE_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 20, 3} // 20-3

// 科研队在休息
var ERR_LEARN_TECHNOLOGY_FAIL_WORKER_REST = pbutil.StaticBuffer{3, 2, 20, 4} // 20-4

// 选择的等级不是下一级数据，只能一级级的学
var ERR_LEARN_TECHNOLOGY_FAIL_NOT_NEXT_LEVEL = pbutil.StaticBuffer{3, 2, 20, 5} // 20-5

// 建筑等级不满足升级条件
var ERR_LEARN_TECHNOLOGY_FAIL_PRE_BUILDING_LEVEL_INVALID = pbutil.StaticBuffer{3, 2, 20, 7} // 20-7

// 科技等级不满足升级条件
var ERR_LEARN_TECHNOLOGY_FAIL_PRE_TECH_LEVEL_INVALID = pbutil.StaticBuffer{3, 2, 20, 8} // 20-8

// 服务器忙，请稍后再试
var ERR_LEARN_TECHNOLOGY_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 20, 6} // 20-6

func NewS2cUnlockStableBuildingMsg(building int32) pbutil.Buffer {
	msg := &S2CUnlockStableBuildingProto{
		Building: building,
	}
	return NewS2cUnlockStableBuildingProtoMsg(msg)
}

var s2c_unlock_stable_building = [...]byte{2, 88} // 88
func NewS2cUnlockStableBuildingProtoMsg(object *S2CUnlockStableBuildingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_unlock_stable_building[:], "s2c_unlock_stable_building")

}

// 无效的类型
var ERR_UNLOCK_STABLE_BUILDING_FAIL_INVALID_TYPE = pbutil.StaticBuffer{3, 2, 89, 1} // 89-1

// 主城已经流亡
var ERR_UNLOCK_STABLE_BUILDING_FAIL_BASE_DEAD = pbutil.StaticBuffer{3, 2, 89, 2} // 89-2

// 已经解锁
var ERR_UNLOCK_STABLE_BUILDING_FAIL_UNLOCKED = pbutil.StaticBuffer{3, 2, 89, 3} // 89-3

// 官府等级不够
var ERR_UNLOCK_STABLE_BUILDING_FAIL_GUAN_FU_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 89, 4} // 89-4

// 君主等级不够
var ERR_UNLOCK_STABLE_BUILDING_FAIL_HERO_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 89, 5} // 89-5

// 主线任务未达成
var ERR_UNLOCK_STABLE_BUILDING_FAIL_MAIN_TASK_NOT_REACH = pbutil.StaticBuffer{3, 2, 89, 6} // 89-6

// 霸业目标阶段未达成
var ERR_UNLOCK_STABLE_BUILDING_FAIL_BA_YE_STAGE_NOT_REACH = pbutil.StaticBuffer{3, 2, 89, 7} // 89-7

// 服务器忙，请稍后再试
var ERR_UNLOCK_STABLE_BUILDING_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 89, 8} // 89-8

func NewS2cUpgradeStableBuildingMsg(building int32, worker_pos int32, worker_rest_end_time int32) pbutil.Buffer {
	msg := &S2CUpgradeStableBuildingProto{
		Building:          building,
		WorkerPos:         worker_pos,
		WorkerRestEndTime: worker_rest_end_time,
	}
	return NewS2cUpgradeStableBuildingProtoMsg(msg)
}

var s2c_upgrade_stable_building = [...]byte{2, 25} // 25
func NewS2cUpgradeStableBuildingProtoMsg(object *S2CUpgradeStableBuildingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_upgrade_stable_building[:], "s2c_upgrade_stable_building")

}

// 无效的类型
var ERR_UPGRADE_STABLE_BUILDING_FAIL_INVALID_TYPE = pbutil.StaticBuffer{3, 2, 26, 6} // 26-6

// 主城已经流亡
var ERR_UPGRADE_STABLE_BUILDING_FAIL_BASE_DEAD = pbutil.StaticBuffer{3, 2, 26, 8} // 26-8

// 等级不一致
var ERR_UPGRADE_STABLE_BUILDING_FAIL_DIFF_LEVEL = pbutil.StaticBuffer{3, 2, 26, 7} // 26-7

// 这个建筑已经最高级了
var ERR_UPGRADE_STABLE_BUILDING_FAIL_MAX_LEVEL = pbutil.StaticBuffer{3, 2, 26, 1} // 26-1

// 建筑队在休息
var ERR_UPGRADE_STABLE_BUILDING_FAIL_WORKER_REST = pbutil.StaticBuffer{3, 2, 26, 2} // 26-2

// 消耗资源不足
var ERR_UPGRADE_STABLE_BUILDING_FAIL_RESOURCE_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 26, 3} // 26-3

// 前提条件未达成
var ERR_UPGRADE_STABLE_BUILDING_FAIL_REQUIRE_NOT_REACH = pbutil.StaticBuffer{3, 2, 26, 4} // 26-4

// 服务器忙，请稍后再试
var ERR_UPGRADE_STABLE_BUILDING_FAIL_SERVER_BUSY = pbutil.StaticBuffer{3, 2, 26, 5} // 26-5

func NewS2cHeroUpdateExpMsg(exp int32) pbutil.Buffer {
	msg := &S2CHeroUpdateExpProto{
		Exp: exp,
	}
	return NewS2cHeroUpdateExpProtoMsg(msg)
}

var s2c_hero_update_exp = [...]byte{2, 22} // 22
func NewS2cHeroUpdateExpProtoMsg(object *S2CHeroUpdateExpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_hero_update_exp[:], "s2c_hero_update_exp")

}

func NewS2cHeroUpgradeLevelMsg(exp int32, level int32) pbutil.Buffer {
	msg := &S2CHeroUpgradeLevelProto{
		Exp:   exp,
		Level: level,
	}
	return NewS2cHeroUpgradeLevelProtoMsg(msg)
}

var s2c_hero_upgrade_level = [...]byte{2, 21} // 21
func NewS2cHeroUpgradeLevelProtoMsg(object *S2CHeroUpgradeLevelProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_hero_upgrade_level[:], "s2c_hero_upgrade_level")

}

func NewS2cHeroUpdateProsperityMsg(prosperity int32, capcity int32) pbutil.Buffer {
	msg := &S2CHeroUpdateProsperityProto{
		Prosperity: prosperity,
		Capcity:    capcity,
	}
	return NewS2cHeroUpdateProsperityProtoMsg(msg)
}

var s2c_hero_update_prosperity = [...]byte{2, 27} // 27
func NewS2cHeroUpdateProsperityProtoMsg(object *S2CHeroUpdateProsperityProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_hero_update_prosperity[:], "s2c_hero_update_prosperity")

}

func NewS2cIsHeroNameExistMsg(name string, exist bool) pbutil.Buffer {
	msg := &S2CIsHeroNameExistProto{
		Name:  name,
		Exist: exist,
	}
	return NewS2cIsHeroNameExistProtoMsg(msg)
}

var s2c_is_hero_name_exist = [...]byte{2, 91} // 91
func NewS2cIsHeroNameExistProtoMsg(object *S2CIsHeroNameExistProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_is_hero_name_exist[:], "s2c_is_hero_name_exist")

}

// 无效的名字长度
var ERR_IS_HERO_NAME_EXIST_FAIL_INVALID_NAME = pbutil.StaticBuffer{3, 2, 92, 1} // 92-1

// 服务器忙，请稍后再试
var ERR_IS_HERO_NAME_EXIST_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 2, 92, 2} // 92-2

func NewS2cChangeHeroNameMsg(name string, next_change_name_time int32) pbutil.Buffer {
	msg := &S2CChangeHeroNameProto{
		Name:               name,
		NextChangeNameTime: next_change_name_time,
	}
	return NewS2cChangeHeroNameProtoMsg(msg)
}

var s2c_change_hero_name = [...]byte{2, 31} // 31
func NewS2cChangeHeroNameProtoMsg(object *S2CChangeHeroNameProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_hero_name[:], "s2c_change_hero_name")

}

// 无效的名字长度
var ERR_CHANGE_HERO_NAME_FAIL_INVALID_NAME = pbutil.StaticBuffer{3, 2, 32, 1} // 32-1

// 改的名字跟当前名字一样
var ERR_CHANGE_HERO_NAME_FAIL_SAME_NAME = pbutil.StaticBuffer{3, 2, 32, 2} // 32-2

// 消耗不足
var ERR_CHANGE_HERO_NAME_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 32, 3} // 32-3

// 改名cd中
var ERR_CHANGE_HERO_NAME_FAIL_CD = pbutil.StaticBuffer{3, 2, 32, 4} // 32-4

// 这个名字已经存在，不能使用
var ERR_CHANGE_HERO_NAME_FAIL_EXIST_NAME = pbutil.StaticBuffer{3, 2, 32, 5} // 32-5

// 名字包含敏感词
var ERR_CHANGE_HERO_NAME_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 2, 32, 7} // 32-7

// 服务器忙，请稍后再试
var ERR_CHANGE_HERO_NAME_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 2, 32, 6} // 32-6

var GIVE_FIRST_CHANGE_HERO_NAME_PRIZE_S2C = pbutil.StaticBuffer{2, 2, 93} // 93

func NewS2cHeroNameChangedBroadcastMsg(id []byte, name string) pbutil.Buffer {
	msg := &S2CHeroNameChangedBroadcastProto{
		Id:   id,
		Name: name,
	}
	return NewS2cHeroNameChangedBroadcastProtoMsg(msg)
}

var s2c_hero_name_changed_broadcast = [...]byte{2, 49} // 49
func NewS2cHeroNameChangedBroadcastProtoMsg(object *S2CHeroNameChangedBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_hero_name_changed_broadcast[:], "s2c_hero_name_changed_broadcast")

}

func NewS2cListOldNameMsg(name []string) pbutil.Buffer {
	msg := &S2CListOldNameProto{
		Name: name,
	}
	return NewS2cListOldNameProtoMsg(msg)
}

var s2c_list_old_name = [...]byte{2, 34} // 34
func NewS2cListOldNameProtoMsg(object *S2CListOldNameProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_old_name[:], "s2c_list_old_name")

}

// 无效的玩家id
var ERR_LIST_OLD_NAME_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 2, 37, 1} // 37-1

// 服务器忙，请稍后再试
var ERR_LIST_OLD_NAME_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 2, 37, 2} // 37-2

func NewS2cViewOtherHeroMsg(hero []byte) pbutil.Buffer {
	msg := &S2CViewOtherHeroProto{
		Hero: hero,
	}
	return NewS2cViewOtherHeroProtoMsg(msg)
}

var s2c_view_other_hero = [...]byte{2, 36} // 36
func NewS2cViewOtherHeroProtoMsg(object *S2CViewOtherHeroProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_other_hero[:], "s2c_view_other_hero")

}

// 无效的玩家id
var ERR_VIEW_OTHER_HERO_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 2, 38, 1} // 38-1

// 服务器忙，请稍后再试
var ERR_VIEW_OTHER_HERO_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 2, 38, 2} // 38-2

func NewS2cViewFightInfoMsg(id []byte, realm_fight_success int32, realm_fight_fail int32, realm_assist int32, realm_been_assist int32, inverstigation int32, been_inverstigation int32) pbutil.Buffer {
	msg := &S2CViewFightInfoProto{
		Id:                 id,
		RealmFightSuccess:  realm_fight_success,
		RealmFightFail:     realm_fight_fail,
		RealmAssist:        realm_assist,
		RealmBeenAssist:    realm_been_assist,
		Inverstigation:     inverstigation,
		BeenInverstigation: been_inverstigation,
	}
	return NewS2cViewFightInfoProtoMsg(msg)
}

var s2c_view_fight_info = [...]byte{2, 126} // 126
func NewS2cViewFightInfoProtoMsg(object *S2CViewFightInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_fight_info[:], "s2c_view_fight_info")

}

func NewS2cUpdateBuildingWorkerCoefMsg(coef int32) pbutil.Buffer {
	msg := &S2CUpdateBuildingWorkerCoefProto{
		Coef: coef,
	}
	return NewS2cUpdateBuildingWorkerCoefProtoMsg(msg)
}

var s2c_update_building_worker_coef = [...]byte{2, 39} // 39
func NewS2cUpdateBuildingWorkerCoefProtoMsg(object *S2CUpdateBuildingWorkerCoefProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_building_worker_coef[:], "s2c_update_building_worker_coef")

}

func NewS2cUpdateTechWorkerCoefMsg(coef int32) pbutil.Buffer {
	msg := &S2CUpdateTechWorkerCoefProto{
		Coef: coef,
	}
	return NewS2cUpdateTechWorkerCoefProtoMsg(msg)
}

var s2c_update_tech_worker_coef = [...]byte{2, 40} // 40
func NewS2cUpdateTechWorkerCoefProtoMsg(object *S2CUpdateTechWorkerCoefProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_tech_worker_coef[:], "s2c_update_tech_worker_coef")

}

func NewS2cUpdateBuildingWorkerFatigueDurationMsg(fatigue int32) pbutil.Buffer {
	msg := &S2CUpdateBuildingWorkerFatigueDurationProto{
		Fatigue: fatigue,
	}
	return NewS2cUpdateBuildingWorkerFatigueDurationProtoMsg(msg)
}

var s2c_update_building_worker_fatigue_duration = [...]byte{2, 55} // 55
func NewS2cUpdateBuildingWorkerFatigueDurationProtoMsg(object *S2CUpdateBuildingWorkerFatigueDurationProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_building_worker_fatigue_duration[:], "s2c_update_building_worker_fatigue_duration")

}

func NewS2cUpdateTechWorkerFatigueDurationMsg(fatigue int32) pbutil.Buffer {
	msg := &S2CUpdateTechWorkerFatigueDurationProto{
		Fatigue: fatigue,
	}
	return NewS2cUpdateTechWorkerFatigueDurationProtoMsg(msg)
}

var s2c_update_tech_worker_fatigue_duration = [...]byte{2, 56} // 56
func NewS2cUpdateTechWorkerFatigueDurationProtoMsg(object *S2CUpdateTechWorkerFatigueDurationProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_tech_worker_fatigue_duration[:], "s2c_update_tech_worker_fatigue_duration")

}

func NewS2cMiaoBuildingWorkerCdMsg(worker_pos int32) pbutil.Buffer {
	msg := &S2CMiaoBuildingWorkerCdProto{
		WorkerPos: worker_pos,
	}
	return NewS2cMiaoBuildingWorkerCdProtoMsg(msg)
}

var s2c_miao_building_worker_cd = [...]byte{2, 42} // 42
func NewS2cMiaoBuildingWorkerCdProtoMsg(object *S2CMiaoBuildingWorkerCdProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_miao_building_worker_cd[:], "s2c_miao_building_worker_cd")

}

// 无效的建筑队序号
var ERR_MIAO_BUILDING_WORKER_CD_FAIL_INVALID_POS = pbutil.StaticBuffer{3, 2, 43, 4} // 43-4

// 建筑队不是CD中
var ERR_MIAO_BUILDING_WORKER_CD_FAIL_NOT_WORKING = pbutil.StaticBuffer{3, 2, 43, 1} // 43-1

// 消耗不足
var ERR_MIAO_BUILDING_WORKER_CD_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 43, 2} // 43-2

// 秒CD功能没有开启
var ERR_MIAO_BUILDING_WORKER_CD_FAIL_ZERO_DURATION = pbutil.StaticBuffer{3, 2, 43, 3} // 43-3

func NewS2cMiaoTechWorkerCdMsg(worker_pos int32) pbutil.Buffer {
	msg := &S2CMiaoTechWorkerCdProto{
		WorkerPos: worker_pos,
	}
	return NewS2cMiaoTechWorkerCdProtoMsg(msg)
}

var s2c_miao_tech_worker_cd = [...]byte{2, 45} // 45
func NewS2cMiaoTechWorkerCdProtoMsg(object *S2CMiaoTechWorkerCdProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_miao_tech_worker_cd[:], "s2c_miao_tech_worker_cd")

}

// 无效的科研队序号
var ERR_MIAO_TECH_WORKER_CD_FAIL_INVALID_POS = pbutil.StaticBuffer{3, 2, 46, 4} // 46-4

// 科研队不是CD中
var ERR_MIAO_TECH_WORKER_CD_FAIL_NOT_WORKING = pbutil.StaticBuffer{3, 2, 46, 1} // 46-1

// 消耗不足
var ERR_MIAO_TECH_WORKER_CD_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 46, 2} // 46-2

// 秒CD功能没有开启
var ERR_MIAO_TECH_WORKER_CD_FAIL_ZERO_DURATION = pbutil.StaticBuffer{3, 2, 46, 3} // 46-3

func NewS2cUpdateYuanbaoMsg(yuanbao int32) pbutil.Buffer {
	msg := &S2CUpdateYuanbaoProto{
		Yuanbao: yuanbao,
	}
	return NewS2cUpdateYuanbaoProtoMsg(msg)
}

var s2c_update_yuanbao = [...]byte{2, 47} // 47
func NewS2cUpdateYuanbaoProtoMsg(object *S2CUpdateYuanbaoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_yuanbao[:], "s2c_update_yuanbao")

}

func NewS2cUpdateYuanbaoGiftLimitMsg(yuanbao_gift_limit int32) pbutil.Buffer {
	msg := &S2CUpdateYuanbaoGiftLimitProto{
		YuanbaoGiftLimit: yuanbao_gift_limit,
	}
	return NewS2cUpdateYuanbaoGiftLimitProtoMsg(msg)
}

var s2c_update_yuanbao_gift_limit = [...]byte{2, 169, 1} // 169
func NewS2cUpdateYuanbaoGiftLimitProtoMsg(object *S2CUpdateYuanbaoGiftLimitProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_yuanbao_gift_limit[:], "s2c_update_yuanbao_gift_limit")

}

func NewS2cUpdateDianquanMsg(dianquan int32) pbutil.Buffer {
	msg := &S2CUpdateDianquanProto{
		Dianquan: dianquan,
	}
	return NewS2cUpdateDianquanProtoMsg(msg)
}

var s2c_update_dianquan = [...]byte{2, 140, 1} // 140
func NewS2cUpdateDianquanProtoMsg(object *S2CUpdateDianquanProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_dianquan[:], "s2c_update_dianquan")

}

func NewS2cUpdateYinliangMsg(yinliang int32) pbutil.Buffer {
	msg := &S2CUpdateYinliangProto{
		Yinliang: yinliang,
	}
	return NewS2cUpdateYinliangProtoMsg(msg)
}

var s2c_update_yinliang = [...]byte{2, 141, 1} // 141
func NewS2cUpdateYinliangProtoMsg(object *S2CUpdateYinliangProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_yinliang[:], "s2c_update_yinliang")

}

func NewS2cUpdateHeroFightAmountMsg(fight_amount int32) pbutil.Buffer {
	msg := &S2CUpdateHeroFightAmountProto{
		FightAmount: fight_amount,
	}
	return NewS2cUpdateHeroFightAmountProtoMsg(msg)
}

var s2c_update_hero_fight_amount = [...]byte{2, 48} // 48
func NewS2cUpdateHeroFightAmountProtoMsg(object *S2CUpdateHeroFightAmountProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_hero_fight_amount[:], "s2c_update_hero_fight_amount")

}

func NewS2cRecoveryForgingTimeChangeMsg(times int32, next_time int32) pbutil.Buffer {
	msg := &S2CRecoveryForgingTimeChangeProto{
		Times:    times,
		NextTime: next_time,
	}
	return NewS2cRecoveryForgingTimeChangeProtoMsg(msg)
}

var s2c_recovery_forging_time_change = [...]byte{2, 54} // 54
func NewS2cRecoveryForgingTimeChangeProtoMsg(object *S2CRecoveryForgingTimeChangeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_recovery_forging_time_change[:], "s2c_recovery_forging_time_change")

}

func NewS2cForgingEquipMsg(slot int32) pbutil.Buffer {
	msg := &S2CForgingEquipProto{
		Slot: slot,
	}
	return NewS2cForgingEquipProtoMsg(msg)
}

var s2c_forging_equip = [...]byte{2, 52} // 52
func NewS2cForgingEquipProtoMsg(object *S2CForgingEquipProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_forging_equip[:], "s2c_forging_equip")

}

// 次数不够
var ERR_FORGING_EQUIP_FAIL_TIMES_NOT_ENOUGH = pbutil.StaticBuffer{3, 2, 53, 1} // 53-1

// 这件装备不可以锻造
var ERR_FORGING_EQUIP_FAIL_CAN_NOT_FORGING_EQIUP = pbutil.StaticBuffer{3, 2, 53, 2} // 53-2

// 打造数量无效
var ERR_FORGING_EQUIP_FAIL_COUNT_INVALID = pbutil.StaticBuffer{3, 2, 53, 3} // 53-3

// 功能没开启
var ERR_FORGING_EQUIP_FAIL_FUNCTION_NOT_OPEN = pbutil.StaticBuffer{3, 2, 53, 4} // 53-4

// 一键没有开启
var ERR_FORGING_EQUIP_FAIL_ONE_KEY_NOT_OPEN = pbutil.StaticBuffer{3, 2, 53, 5} // 53-5

func NewS2cUpdateNewForgingPosMsg(new_forging_pos []int32) pbutil.Buffer {
	msg := &S2CUpdateNewForgingPosProto{
		NewForgingPos: new_forging_pos,
	}
	return NewS2cUpdateNewForgingPosProtoMsg(msg)
}

var s2c_update_new_forging_pos = [...]byte{2, 117} // 117
func NewS2cUpdateNewForgingPosProtoMsg(object *S2CUpdateNewForgingPosProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_new_forging_pos[:], "s2c_update_new_forging_pos")

}

var SIGN_S2C = pbutil.StaticBuffer{2, 2, 67} // 67

// 签名长度非法
var ERR_SIGN_FAIL_LEN_INVALID = pbutil.StaticBuffer{3, 2, 68, 1} // 68-1

// 输入包含敏感词
var ERR_SIGN_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 2, 68, 2} // 68-2

// 服务器忙，请稍后再试
var ERR_SIGN_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 2, 68, 3} // 68-3

var VOICE_S2C = pbutil.StaticBuffer{2, 2, 70} // 70

// 签名太长
var ERR_VOICE_FAIL_LEN_INVALID = pbutil.StaticBuffer{3, 2, 71, 1} // 71-1

func NewS2cBuildingWorkerTimeChangedMsg(worker_pos int32, worker_rest_end_time int32) pbutil.Buffer {
	msg := &S2CBuildingWorkerTimeChangedProto{
		WorkerPos:         worker_pos,
		WorkerRestEndTime: worker_rest_end_time,
	}
	return NewS2cBuildingWorkerTimeChangedProtoMsg(msg)
}

var s2c_building_worker_time_changed = [...]byte{2, 57} // 57
func NewS2cBuildingWorkerTimeChangedProtoMsg(object *S2CBuildingWorkerTimeChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_building_worker_time_changed[:], "s2c_building_worker_time_changed")

}

func NewS2cTechWorkerTimeChangedMsg(worker_pos int32, worker_rest_end_time int32) pbutil.Buffer {
	msg := &S2CTechWorkerTimeChangedProto{
		WorkerPos:         worker_pos,
		WorkerRestEndTime: worker_rest_end_time,
	}
	return NewS2cTechWorkerTimeChangedProtoMsg(msg)
}

var s2c_tech_worker_time_changed = [...]byte{2, 58} // 58
func NewS2cTechWorkerTimeChangedProtoMsg(object *S2CTechWorkerTimeChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_tech_worker_time_changed[:], "s2c_tech_worker_time_changed")

}

func NewS2cCityEventTimeChangedMsg(next_time int32) pbutil.Buffer {
	msg := &S2CCityEventTimeChangedProto{
		NextTime: next_time,
	}
	return NewS2cCityEventTimeChangedProtoMsg(msg)
}

var s2c_city_event_time_changed = [...]byte{2, 59} // 59
func NewS2cCityEventTimeChangedProtoMsg(object *S2CCityEventTimeChangedProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_city_event_time_changed[:], "s2c_city_event_time_changed")

}

func NewS2cRequestCityExchangeEventMsg(accept_times int32, event_id int32) pbutil.Buffer {
	msg := &S2CRequestCityExchangeEventProto{
		AcceptTimes: accept_times,
		EventId:     event_id,
	}
	return NewS2cRequestCityExchangeEventProtoMsg(msg)
}

var s2c_request_city_exchange_event = [...]byte{2, 61} // 61
func NewS2cRequestCityExchangeEventProtoMsg(object *S2CRequestCityExchangeEventProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_request_city_exchange_event[:], "s2c_request_city_exchange_event")

}

// 事件时间没到
var ERR_REQUEST_CITY_EXCHANGE_EVENT_FAIL_IN_CD = pbutil.StaticBuffer{3, 2, 62, 1} // 62-1

// 没有次数了
var ERR_REQUEST_CITY_EXCHANGE_EVENT_FAIL_NO_TIMES = pbutil.StaticBuffer{3, 2, 62, 2} // 62-2

// 功能没开启
var ERR_REQUEST_CITY_EXCHANGE_EVENT_FAIL_NOT_OPEN = pbutil.StaticBuffer{3, 2, 62, 3} // 62-3

func NewS2cCityEventExchangeMsg(give_up bool) pbutil.Buffer {
	msg := &S2CCityEventExchangeProto{
		GiveUp: give_up,
	}
	return NewS2cCityEventExchangeProtoMsg(msg)
}

var s2c_city_event_exchange = [...]byte{2, 64} // 64
func NewS2cCityEventExchangeProtoMsg(object *S2CCityEventExchangeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_city_event_exchange[:], "s2c_city_event_exchange")

}

// 事件时间没到
var ERR_CITY_EVENT_EXCHANGE_FAIL_IN_CD = pbutil.StaticBuffer{3, 2, 65, 1} // 65-1

// 条件不满足
var ERR_CITY_EVENT_EXCHANGE_FAIL_CONDITION_NOT_SATISFIED = pbutil.StaticBuffer{3, 2, 65, 2} // 65-2

// 没有次数了
var ERR_CITY_EVENT_EXCHANGE_FAIL_NO_TIMES = pbutil.StaticBuffer{3, 2, 65, 3} // 65-3

// 功能没开启
var ERR_CITY_EVENT_EXCHANGE_FAIL_NOT_OPEN = pbutil.StaticBuffer{3, 2, 65, 4} // 65-4

func NewS2cUpdateStrategyRestoreStartTimeMsg(time int32) pbutil.Buffer {
	msg := &S2CUpdateStrategyRestoreStartTimeProto{
		Time: time,
	}
	return NewS2cUpdateStrategyRestoreStartTimeProtoMsg(msg)
}

var s2c_update_strategy_restore_start_time = [...]byte{2, 72} // 72
func NewS2cUpdateStrategyRestoreStartTimeProtoMsg(object *S2CUpdateStrategyRestoreStartTimeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_strategy_restore_start_time[:], "s2c_update_strategy_restore_start_time")

}

func NewS2cUpdateStrategyNextUseTimeMsg(id int32, time int32) pbutil.Buffer {
	msg := &S2CUpdateStrategyNextUseTimeProto{
		Id:   id,
		Time: time,
	}
	return NewS2cUpdateStrategyNextUseTimeProtoMsg(msg)
}

var s2c_update_strategy_next_use_time = [...]byte{2, 73} // 73
func NewS2cUpdateStrategyNextUseTimeProtoMsg(object *S2CUpdateStrategyNextUseTimeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_strategy_next_use_time[:], "s2c_update_strategy_next_use_time")

}

func NewS2cUpdateJadeOreMsg(amount int32) pbutil.Buffer {
	msg := &S2CUpdateJadeOreProto{
		Amount: amount,
	}
	return NewS2cUpdateJadeOreProtoMsg(msg)
}

var s2c_update_jade_ore = [...]byte{2, 85} // 85
func NewS2cUpdateJadeOreProtoMsg(object *S2CUpdateJadeOreProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_jade_ore[:], "s2c_update_jade_ore")

}

func NewS2cUpdateJadeMsg(amount int32, history_jade int32, today_obtain_jade int32) pbutil.Buffer {
	msg := &S2CUpdateJadeProto{
		Amount:          amount,
		HistoryJade:     history_jade,
		TodayObtainJade: today_obtain_jade,
	}
	return NewS2cUpdateJadeProtoMsg(msg)
}

var s2c_update_jade = [...]byte{2, 86} // 86
func NewS2cUpdateJadeProtoMsg(object *S2CUpdateJadeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_jade[:], "s2c_update_jade")

}

func NewS2cChangeHeadMsg(head_id string) pbutil.Buffer {
	msg := &S2CChangeHeadProto{
		HeadId: head_id,
	}
	return NewS2cChangeHeadProtoMsg(msg)
}

var s2c_change_head = [...]byte{2, 95} // 95
func NewS2cChangeHeadProtoMsg(object *S2CChangeHeadProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_head[:], "s2c_change_head")

}

// 无效的头像id
var ERR_CHANGE_HEAD_FAIL_INVALID_HEAD = pbutil.StaticBuffer{3, 2, 96, 4} // 96-4

// 将魂没有解锁
var ERR_CHANGE_HEAD_FAIL_CAPTAIN_SOUL_NOT_UNLOCK = pbutil.StaticBuffer{3, 2, 96, 2} // 96-2

// 君主等级不足
var ERR_CHANGE_HEAD_FAIL_HERO_LEVEL_TOO_LOW = pbutil.StaticBuffer{3, 2, 96, 5} // 96-5

// 服务器繁忙，请稍后再试
var ERR_CHANGE_HEAD_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 2, 96, 3} // 96-3

// 没有对应官职
var ERR_CHANGE_HEAD_FAIL_NOT_THIS_OFFICIAL = pbutil.StaticBuffer{3, 2, 96, 6} // 96-6

func NewS2cChangeBodyMsg(body_id int32) pbutil.Buffer {
	msg := &S2CChangeBodyProto{
		BodyId: body_id,
	}
	return NewS2cChangeBodyProtoMsg(msg)
}

var s2c_change_body = [...]byte{2, 131, 1} // 131
func NewS2cChangeBodyProtoMsg(object *S2CChangeBodyProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_body[:], "s2c_change_body")

}

// 无效的形象id
var ERR_CHANGE_BODY_FAIL_INVALID_BODY = pbutil.StaticBuffer{4, 2, 132, 1, 1} // 132-1

// 将魂没有解锁
var ERR_CHANGE_BODY_FAIL_CAPTAIN_SOUL_NOT_UNLOCK = pbutil.StaticBuffer{4, 2, 132, 1, 2} // 132-2

// 君主等级不足
var ERR_CHANGE_BODY_FAIL_HERO_LEVEL_TOO_LOW = pbutil.StaticBuffer{4, 2, 132, 1, 3} // 132-3

// 服务器繁忙，请稍后再试
var ERR_CHANGE_BODY_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 2, 132, 1, 4} // 132-4

// 没有对应官职
var ERR_CHANGE_BODY_FAIL_NOT_THIS_OFFICIAL = pbutil.StaticBuffer{4, 2, 132, 1, 5} // 132-5

func NewS2cCollectCountdownPrizeMsg(prize []byte, collect_time int32, desc_id int32, prosprity int32) pbutil.Buffer {
	msg := &S2CCollectCountdownPrizeProto{
		Prize:       prize,
		CollectTime: collect_time,
		DescId:      desc_id,
		Prosprity:   prosprity,
	}
	return NewS2cCollectCountdownPrizeProtoMsg(msg)
}

func NewS2cCollectCountdownPrizeMarshalMsg(prize marshaler, collect_time int32, desc_id int32, prosprity int32) pbutil.Buffer {
	msg := &S2CCollectCountdownPrizeProto{
		Prize:       safeMarshal(prize),
		CollectTime: collect_time,
		DescId:      desc_id,
		Prosprity:   prosprity,
	}
	return NewS2cCollectCountdownPrizeProtoMsg(msg)
}

var s2c_collect_countdown_prize = [...]byte{2, 115} // 115
func NewS2cCollectCountdownPrizeProtoMsg(object *S2CCollectCountdownPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_countdown_prize[:], "s2c_collect_countdown_prize")

}

// 倒计时未到
var ERR_COLLECT_COUNTDOWN_PRIZE_FAIL_TIME_NOT_REACHED = pbutil.StaticBuffer{3, 2, 116, 1} // 116-1

func NewS2cListWorkshopEquipmentMsg(refresh_time int32, equipment []int32, duration []int32, index int32, collect_time int32, workshop_refresh_times int32) pbutil.Buffer {
	msg := &S2CListWorkshopEquipmentProto{
		RefreshTime:          refresh_time,
		Equipment:            equipment,
		Duration:             duration,
		Index:                index,
		CollectTime:          collect_time,
		WorkshopRefreshTimes: workshop_refresh_times,
	}
	return NewS2cListWorkshopEquipmentProtoMsg(msg)
}

var s2c_list_workshop_equipment = [...]byte{2, 118} // 118
func NewS2cListWorkshopEquipmentProtoMsg(object *S2CListWorkshopEquipmentProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_list_workshop_equipment[:], "s2c_list_workshop_equipment")

}

func NewS2cStartWorkshopMsg(index int32, collect_time int32) pbutil.Buffer {
	msg := &S2CStartWorkshopProto{
		Index:       index,
		CollectTime: collect_time,
	}
	return NewS2cStartWorkshopProtoMsg(msg)
}

var s2c_start_workshop = [...]byte{2, 120} // 120
func NewS2cStartWorkshopProtoMsg(object *S2CStartWorkshopProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_start_workshop[:], "s2c_start_workshop")

}

// 无效的锻造装备索引
var ERR_START_WORKSHOP_FAIL_INVALID_INDEX = pbutil.StaticBuffer{3, 2, 121, 1} // 121-1

// 这个装备锻造中
var ERR_START_WORKSHOP_FAIL_FORGING = pbutil.StaticBuffer{3, 2, 121, 2} // 121-2

// 达到锻造个数上限
var ERR_START_WORKSHOP_FAIL_COUNT_LIMIT = pbutil.StaticBuffer{3, 2, 121, 3} // 121-3

func NewS2cCollectWorkshopMsg(index int32) pbutil.Buffer {
	msg := &S2CCollectWorkshopProto{
		Index: index,
	}
	return NewS2cCollectWorkshopProtoMsg(msg)
}

var s2c_collect_workshop = [...]byte{2, 123} // 123
func NewS2cCollectWorkshopProtoMsg(object *S2CCollectWorkshopProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_workshop[:], "s2c_collect_workshop")

}

// 无效的锻造装备索引
var ERR_COLLECT_WORKSHOP_FAIL_INVALID_INDEX = pbutil.StaticBuffer{3, 2, 124, 1} // 124-1

// 装备锻造没有完成，不能领取
var ERR_COLLECT_WORKSHOP_FAIL_CANT_COLLECT = pbutil.StaticBuffer{3, 2, 124, 2} // 124-2

func NewS2cWorkshopMiaoCdMsg(index int32) pbutil.Buffer {
	msg := &S2CWorkshopMiaoCdProto{
		Index: index,
	}
	return NewS2cWorkshopMiaoCdProtoMsg(msg)
}

var s2c_workshop_miao_cd = [...]byte{2, 128, 1} // 128
func NewS2cWorkshopMiaoCdProtoMsg(object *S2CWorkshopMiaoCdProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_workshop_miao_cd[:], "s2c_workshop_miao_cd")

}

// 无效的锻造装备索引
var ERR_WORKSHOP_MIAO_CD_FAIL_INVALID_INDEX = pbutil.StaticBuffer{4, 2, 129, 1, 1} // 129-1

// CD已结束
var ERR_WORKSHOP_MIAO_CD_FAIL_NOT_IN_CD = pbutil.StaticBuffer{4, 2, 129, 1, 4} // 129-4

// 消耗不够
var ERR_WORKSHOP_MIAO_CD_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 2, 129, 1, 3} // 129-3

var REFRESH_WORKSHOP_S2C = pbutil.StaticBuffer{3, 2, 134, 1} // 134

// 消耗不够
var ERR_REFRESH_WORKSHOP_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 2, 135, 1, 1} // 135-1

// 没有次数
var ERR_REFRESH_WORKSHOP_FAIL_TIMES_NOT_ENOUGH = pbutil.StaticBuffer{4, 2, 135, 1, 2} // 135-2

// 有装备没有锻造
var ERR_REFRESH_WORKSHOP_FAIL_HAS_EQUIP_NOT_FORG = pbutil.StaticBuffer{4, 2, 135, 1, 3} // 135-3

var COLLECT_SEASON_PRIZE_S2C = pbutil.StaticBuffer{3, 2, 137, 1} // 137

// 奖励已经领取
var ERR_COLLECT_SEASON_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{4, 2, 138, 1, 1} // 138-1

func NewS2cSeasonStartBroadcastMsg(season shared_proto.Season, start_time int32, is_reset bool) pbutil.Buffer {
	msg := &S2CSeasonStartBroadcastProto{
		Season:    season,
		StartTime: start_time,
		IsReset:   is_reset,
	}
	return NewS2cSeasonStartBroadcastProtoMsg(msg)
}

var s2c_season_start_broadcast = [...]byte{2, 139, 1} // 139
func NewS2cSeasonStartBroadcastProtoMsg(object *S2CSeasonStartBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_season_start_broadcast[:], "s2c_season_start_broadcast")

}

func NewS2cUpdateCostReduceCoefMsg(building int32, tech int32) pbutil.Buffer {
	msg := &S2CUpdateCostReduceCoefProto{
		Building: building,
		Tech:     tech,
	}
	return NewS2cUpdateCostReduceCoefProtoMsg(msg)
}

var s2c_update_cost_reduce_coef = [...]byte{2, 145, 1} // 145
func NewS2cUpdateCostReduceCoefProtoMsg(object *S2CUpdateCostReduceCoefProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_cost_reduce_coef[:], "s2c_update_cost_reduce_coef")

}

func NewS2cUpdateSpMsg(sp int32) pbutil.Buffer {
	msg := &S2CUpdateSpProto{
		Sp: sp,
	}
	return NewS2cUpdateSpProtoMsg(msg)
}

var s2c_update_sp = [...]byte{2, 146, 1} // 146
func NewS2cUpdateSpProtoMsg(object *S2CUpdateSpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_sp[:], "s2c_update_sp")

}

func NewS2cBuySpMsg(sp int32, buy_sp_times int32) pbutil.Buffer {
	msg := &S2CBuySpProto{
		Sp:         sp,
		BuySpTimes: buy_sp_times,
	}
	return NewS2cBuySpProtoMsg(msg)
}

var s2c_buy_sp = [...]byte{2, 148, 1} // 148
func NewS2cBuySpProtoMsg(object *S2CBuySpProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_buy_sp[:], "s2c_buy_sp")

}

// 无效的数据
var ERR_BUY_SP_FAIL_INVALID_DATA = pbutil.StaticBuffer{4, 2, 149, 1, 1} // 149-1

// 购买次数上限
var ERR_BUY_SP_FAIL_BUY_TIMES_LIMIT = pbutil.StaticBuffer{4, 2, 149, 1, 2} // 149-2

// 元宝不足
var ERR_BUY_SP_FAIL_NOT_ENOUGH_YUANBAO = pbutil.StaticBuffer{4, 2, 149, 1, 3} // 149-3

func NewS2cUseBufEffectMsg(buf_effect *shared_proto.BufferEffectProto) pbutil.Buffer {
	msg := &S2CUseBufEffectProto{
		BufEffect: buf_effect,
	}
	return NewS2cUseBufEffectProtoMsg(msg)
}

var s2c_use_buf_effect = [...]byte{2, 151, 1} // 151
func NewS2cUseBufEffectProtoMsg(object *S2CUseBufEffectProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_buf_effect[:], "s2c_use_buf_effect")

}

// 不存在的增益效果
var ERR_USE_BUF_EFFECT_FAIL_NO_BUF_EFFECT = pbutil.StaticBuffer{4, 2, 152, 1, 1} // 152-1

// 同种增益不可叠加
var ERR_USE_BUF_EFFECT_FAIL_SAME_BUF_EFFECT = pbutil.StaticBuffer{4, 2, 152, 1, 2} // 152-2

// 低等级增益不可覆盖高等级增益
var ERR_USE_BUF_EFFECT_FAIL_LEVEL_LIMIT = pbutil.StaticBuffer{4, 2, 152, 1, 3} // 152-3

// 低时长增益不可覆盖高时长增益
var ERR_USE_BUF_EFFECT_FAIL_KEEP_TIME_LIMIT = pbutil.StaticBuffer{4, 2, 152, 1, 4} // 152-4

// 道具不足
var ERR_USE_BUF_EFFECT_FAIL_ITEM_NOT_ENOUGH = pbutil.StaticBuffer{4, 2, 152, 1, 5} // 152-5

// buff更新失败
var ERR_USE_BUF_EFFECT_FAIL_BUFF_EFFECT_FAIL = pbutil.StaticBuffer{4, 2, 152, 1, 6} // 152-6

// 服务器繁忙，请稍后再试
var ERR_USE_BUF_EFFECT_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 2, 152, 1, 7} // 152-7

func NewS2cOpenBufEffectUiMsg(buffers *shared_proto.HeroBufferProto) pbutil.Buffer {
	msg := &S2COpenBufEffectUiProto{
		Buffers: buffers,
	}
	return NewS2cOpenBufEffectUiProtoMsg(msg)
}

var s2c_open_buf_effect_ui = [...]byte{2, 155, 1} // 155
func NewS2cOpenBufEffectUiProtoMsg(object *S2COpenBufEffectUiProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_open_buf_effect_ui[:], "s2c_open_buf_effect_ui")

}

// 服务器繁忙，请稍后再试
var ERR_OPEN_BUF_EFFECT_UI_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 2, 156, 1, 1} // 156-1

func NewS2cUseAdvantageMsg(id int32, start_time int32, end_time int32) pbutil.Buffer {
	msg := &S2CUseAdvantageProto{
		Id:        id,
		StartTime: start_time,
		EndTime:   end_time,
	}
	return NewS2cUseAdvantageProtoMsg(msg)
}

var s2c_use_advantage = [...]byte{2, 159, 1} // 159
func NewS2cUseAdvantageProtoMsg(object *S2CUseAdvantageProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_advantage[:], "s2c_use_advantage")

}

// 不存在的增益
var ERR_USE_ADVANTAGE_FAIL_INVALID_ID = pbutil.StaticBuffer{4, 2, 160, 1, 1} // 160-1

// 同种增益不可叠加
var ERR_USE_ADVANTAGE_FAIL_SAME_TYPE = pbutil.StaticBuffer{4, 2, 160, 1, 2} // 160-2

// 低等级增益不可覆盖高等级增益
var ERR_USE_ADVANTAGE_FAIL_LEVEL_LIMIT = pbutil.StaticBuffer{4, 2, 160, 1, 3} // 160-3

// 低时长增益不可覆盖高时长增益
var ERR_USE_ADVANTAGE_FAIL_KEEP_TIME_LIMIT = pbutil.StaticBuffer{4, 2, 160, 1, 4} // 160-4

// 道具不足
var ERR_USE_ADVANTAGE_FAIL_ITEM_NOT_ENOUGH = pbutil.StaticBuffer{4, 2, 160, 1, 5} // 160-5

// buff更新失败
var ERR_USE_ADVANTAGE_FAIL_BUFF_EFFECT_FAIL = pbutil.StaticBuffer{4, 2, 160, 1, 6} // 160-6

// 服务器繁忙，请稍后再试
var ERR_USE_ADVANTAGE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{4, 2, 160, 1, 7} // 160-7

func NewS2cUpdateAdvantageCountMsg(count int32) pbutil.Buffer {
	msg := &S2CUpdateAdvantageCountProto{
		Count: count,
	}
	return NewS2cUpdateAdvantageCountProtoMsg(msg)
}

var s2c_update_advantage_count = [...]byte{2, 161, 1} // 161
func NewS2cUpdateAdvantageCountProtoMsg(object *S2CUpdateAdvantageCountProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_advantage_count[:], "s2c_update_advantage_count")

}

func NewS2cWorkerUnlockMsg(pos int32, new_lock_start_time int32) pbutil.Buffer {
	msg := &S2CWorkerUnlockProto{
		Pos:              pos,
		NewLockStartTime: new_lock_start_time,
	}
	return NewS2cWorkerUnlockProtoMsg(msg)
}

var s2c_worker_unlock = [...]byte{2, 167, 1} // 167
func NewS2cWorkerUnlockProtoMsg(object *S2CWorkerUnlockProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_worker_unlock[:], "s2c_worker_unlock")

}

// pos错误
var ERR_WORKER_UNLOCK_FAIL_INVALID_POS = pbutil.StaticBuffer{4, 2, 168, 1, 1} // 168-1

// pos在解锁中
var ERR_WORKER_UNLOCK_FAIL_POS_IS_UNLOCKED = pbutil.StaticBuffer{4, 2, 168, 1, 2} // 168-2

// 钱不够
var ERR_WORKER_UNLOCK_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{4, 2, 168, 1, 3} // 168-3

func NewS2cWorkerAlwaysUnlockMsg(pos int32) pbutil.Buffer {
	msg := &S2CWorkerAlwaysUnlockProto{
		Pos: pos,
	}
	return NewS2cWorkerAlwaysUnlockProtoMsg(msg)
}

var s2c_worker_always_unlock = [...]byte{2, 165, 1} // 165
func NewS2cWorkerAlwaysUnlockProtoMsg(object *S2CWorkerAlwaysUnlockProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_worker_always_unlock[:], "s2c_worker_always_unlock")

}
