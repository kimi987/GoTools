package farm

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
	MODULE_ID = 38

	C2S_PLANT = 2

	C2S_HARVEST = 5

	C2S_CHANGE = 8

	C2S_ONE_KEY_PLANT = 12

	C2S_ONE_KEY_HARVEST = 28

	C2S_ONE_KEY_RESET = 52

	C2S_VIEW_FARM = 43

	C2S_STEAL = 18

	C2S_ONE_KEY_STEAL = 31

	C2S_STEAL_LOG_LIST = 39

	C2S_CAN_STEAL_LIST = 48
)

var FARM_IS_UPDATE_S2C = pbutil.StaticBuffer{2, 38, 50} // 50

func NewS2cPlantMsg(cube_x int32, cube_y int32, res_id int32) pbutil.Buffer {
	msg := &S2CPlantProto{
		CubeX: cube_x,
		CubeY: cube_y,
		ResId: res_id,
	}
	return NewS2cPlantProtoMsg(msg)
}

var s2c_plant = [...]byte{38, 3} // 3
func NewS2cPlantProtoMsg(object *S2CPlantProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_plant[:], "s2c_plant")

}

// 无效地块。不是自己的
var ERR_PLANT_FAIL_INVALID_CUBE = pbutil.StaticBuffer{3, 38, 4, 1} // 4-1

// 资源 id 错误
var ERR_PLANT_FAIL_INVALID_RES_ID = pbutil.StaticBuffer{3, 38, 4, 2} // 4-2

// 不是空闲地块
var ERR_PLANT_FAIL_NOT_IDLE_CUBE = pbutil.StaticBuffer{3, 38, 4, 3} // 4-3

// 服务器错误
var ERR_PLANT_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 4, 4} // 4-4

func NewS2cHarvestMsg(cube_x int32, cube_y int32, current_output int32) pbutil.Buffer {
	msg := &S2CHarvestProto{
		CubeX:         cube_x,
		CubeY:         cube_y,
		CurrentOutput: current_output,
	}
	return NewS2cHarvestProtoMsg(msg)
}

var s2c_harvest = [...]byte{38, 6} // 6
func NewS2cHarvestProtoMsg(object *S2CHarvestProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_harvest[:], "s2c_harvest")

}

// 无效地块。不是自己的、没有种植或者在冲突中
var ERR_HARVEST_FAIL_INVALID_CUBE = pbutil.StaticBuffer{3, 38, 7, 1} // 7-1

// 服务器错误
var ERR_HARVEST_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 7, 2} // 7-2

// 产量为0
var ERR_HARVEST_FAIL_NO_OUTPUT = pbutil.StaticBuffer{3, 38, 7, 3} // 7-3

func NewS2cChangeMsg(cube_x int32, cube_y int32, res_id int32, old_res_id int32, old_output int32) pbutil.Buffer {
	msg := &S2CChangeProto{
		CubeX:     cube_x,
		CubeY:     cube_y,
		ResId:     res_id,
		OldResId:  old_res_id,
		OldOutput: old_output,
	}
	return NewS2cChangeProtoMsg(msg)
}

var s2c_change = [...]byte{38, 9} // 9
func NewS2cChangeProtoMsg(object *S2CChangeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change[:], "s2c_change")

}

// 无效地块。不是自己的、没有种植或者在冲突中
var ERR_CHANGE_FAIL_INVALID_CUBE = pbutil.StaticBuffer{3, 38, 10, 1} // 10-1

// 服务器错误
var ERR_CHANGE_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 10, 2} // 10-2

func NewS2cOneKeyPlantMsg(cube_x []int32, cube_y []int32, res_id []int32) pbutil.Buffer {
	msg := &S2COneKeyPlantProto{
		CubeX: cube_x,
		CubeY: cube_y,
		ResId: res_id,
	}
	return NewS2cOneKeyPlantProtoMsg(msg)
}

var s2c_one_key_plant = [...]byte{38, 13} // 13
func NewS2cOneKeyPlantProtoMsg(object *S2COneKeyPlantProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_one_key_plant[:], "s2c_one_key_plant")

}

// 没有空白地块了
var ERR_ONE_KEY_PLANT_FAIL_NONE_IDLE_CUBE = pbutil.StaticBuffer{3, 38, 14, 2} // 14-2

// 服务器错误
var ERR_ONE_KEY_PLANT_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 14, 3} // 14-3

// 数量都是0
var ERR_ONE_KEY_PLANT_FAIL_INVALID_COUNT = pbutil.StaticBuffer{3, 38, 14, 4} // 14-4

func NewS2cOneKeyHarvestMsg(cube_x []int32, cube_y []int32, gold_output []int32, stone_output []int32) pbutil.Buffer {
	msg := &S2COneKeyHarvestProto{
		CubeX:       cube_x,
		CubeY:       cube_y,
		GoldOutput:  gold_output,
		StoneOutput: stone_output,
	}
	return NewS2cOneKeyHarvestProtoMsg(msg)
}

var s2c_one_key_harvest = [...]byte{38, 29} // 29
func NewS2cOneKeyHarvestProtoMsg(object *S2COneKeyHarvestProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_one_key_harvest[:], "s2c_one_key_harvest")

}

// 还没有成熟的地块
var ERR_ONE_KEY_HARVEST_FAIL_NONE_IDLE_CUBE = pbutil.StaticBuffer{3, 38, 30, 1} // 30-1

// 服务器错误
var ERR_ONE_KEY_HARVEST_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 30, 2} // 30-2

// 资源类型错误
var ERR_ONE_KEY_HARVEST_FAIL_RES_TYPE_ERR = pbutil.StaticBuffer{3, 38, 30, 3} // 30-3

func NewS2cOneKeyResetMsg(cube_x []int32, cube_y []int32, gold_output []int32, stone_output []int32) pbutil.Buffer {
	msg := &S2COneKeyResetProto{
		CubeX:       cube_x,
		CubeY:       cube_y,
		GoldOutput:  gold_output,
		StoneOutput: stone_output,
	}
	return NewS2cOneKeyResetProtoMsg(msg)
}

var s2c_one_key_reset = [...]byte{38, 53} // 53
func NewS2cOneKeyResetProtoMsg(object *S2COneKeyResetProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_one_key_reset[:], "s2c_one_key_reset")

}

// 服务器错误
var ERR_ONE_KEY_RESET_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 54, 1} // 54-1

func NewS2cViewFarmMsg(target []byte, target_basic []byte, hero_farm []byte, next_level_cube_x []int32, next_level_cube_y []int32, can_steal bool) pbutil.Buffer {
	msg := &S2CViewFarmProto{
		Target:         target,
		TargetBasic:    target_basic,
		HeroFarm:       hero_farm,
		NextLevelCubeX: next_level_cube_x,
		NextLevelCubeY: next_level_cube_y,
		CanSteal:       can_steal,
	}
	return NewS2cViewFarmProtoMsg(msg)
}

var s2c_view_farm = [...]byte{38, 44} // 44
func NewS2cViewFarmProtoMsg(object *S2CViewFarmProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_view_farm[:], "s2c_view_farm")

}

// 无效的 target
var ERR_VIEW_FARM_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 38, 45, 1} // 45-1

// 服务器错误
var ERR_VIEW_FARM_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 45, 2} // 45-2

func NewS2cStealMsg(target []byte, cube_x int32, cube_y int32, steal_output int32) pbutil.Buffer {
	msg := &S2CStealProto{
		Target:      target,
		CubeX:       cube_x,
		CubeY:       cube_y,
		StealOutput: steal_output,
	}
	return NewS2cStealProtoMsg(msg)
}

var s2c_steal = [...]byte{38, 19} // 19
func NewS2cStealProtoMsg(object *S2CStealProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_steal[:], "s2c_steal")

}

// 无效的 target
var ERR_STEAL_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 38, 20, 1} // 20-1

// 无效的 cube, 没有这块地或没种东西
var ERR_STEAL_FAIL_INVALID_CUBE = pbutil.StaticBuffer{3, 38, 20, 2} // 20-2

// 在保护期
var ERR_STEAL_FAIL_IN_PROTECTED_DURATION = pbutil.StaticBuffer{3, 38, 20, 3} // 20-3

// 这块地偷菜次数用完了
var ERR_STEAL_FAIL_NO_STEAL_TIME = pbutil.StaticBuffer{3, 38, 20, 4} // 20-4

// 服务器错误
var ERR_STEAL_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 20, 5} // 20-5

// 已经偷过了
var ERR_STEAL_FAIL_CUBE_ALREADY_STEALED = pbutil.StaticBuffer{3, 38, 20, 6} // 20-6

// 今天偷满了
var ERR_STEAL_FAIL_DAILY_STEAL_FULL = pbutil.StaticBuffer{3, 38, 20, 7} // 20-7

func NewS2cWhoStealFromMeMsg(target []byte, gold_output int32, stone_output int32) pbutil.Buffer {
	msg := &S2CWhoStealFromMeProto{
		Target:      target,
		GoldOutput:  gold_output,
		StoneOutput: stone_output,
	}
	return NewS2cWhoStealFromMeProtoMsg(msg)
}

var s2c_who_steal_from_me = [...]byte{38, 21} // 21
func NewS2cWhoStealFromMeProtoMsg(object *S2CWhoStealFromMeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_who_steal_from_me[:], "s2c_who_steal_from_me")

}

func NewS2cOneKeyStealMsg(gold_output int32, stone_output int32, cube_x []int32, cube_y []int32, cube_gold_output []int32, cube_stone_output []int32) pbutil.Buffer {
	msg := &S2COneKeyStealProto{
		GoldOutput:      gold_output,
		StoneOutput:     stone_output,
		CubeX:           cube_x,
		CubeY:           cube_y,
		CubeGoldOutput:  cube_gold_output,
		CubeStoneOutput: cube_stone_output,
	}
	return NewS2cOneKeyStealProtoMsg(msg)
}

var s2c_one_key_steal = [...]byte{38, 32} // 32
func NewS2cOneKeyStealProtoMsg(object *S2COneKeyStealProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_one_key_steal[:], "s2c_one_key_steal")

}

// 无效的 target
var ERR_ONE_KEY_STEAL_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 38, 37, 1} // 37-1

// 服务器错误
var ERR_ONE_KEY_STEAL_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 37, 2} // 37-2

// 没有可偷的地块
var ERR_ONE_KEY_STEAL_FAIL_NO_CAN_STEAL_CUBE = pbutil.StaticBuffer{3, 38, 37, 3} // 37-3

// 今天偷满了
var ERR_ONE_KEY_STEAL_FAIL_DAILY_STEAL_FULL = pbutil.StaticBuffer{3, 38, 37, 4} // 37-4

func NewS2cWhoOneKeyStealFromMeMsg(target_id []byte, cube_x []int32, cube_y []int32, steal_times []int32) pbutil.Buffer {
	msg := &S2CWhoOneKeyStealFromMeProto{
		TargetId:   target_id,
		CubeX:      cube_x,
		CubeY:      cube_y,
		StealTimes: steal_times,
	}
	return NewS2cWhoOneKeyStealFromMeProtoMsg(msg)
}

var s2c_who_one_key_steal_from_me = [...]byte{38, 33} // 33
func NewS2cWhoOneKeyStealFromMeProtoMsg(object *S2CWhoOneKeyStealFromMeProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_who_one_key_steal_from_me[:], "s2c_who_one_key_steal_from_me")

}

func NewS2cStealLogListMsg(newest bool, logs []byte) pbutil.Buffer {
	msg := &S2CStealLogListProto{
		Newest: newest,
		Logs:   logs,
	}
	return NewS2cStealLogListProtoMsg(msg)
}

func NewS2cStealLogListMarshalMsg(newest bool, logs marshaler) pbutil.Buffer {
	msg := &S2CStealLogListProto{
		Newest: newest,
		Logs:   safeMarshal(logs),
	}
	return NewS2cStealLogListProtoMsg(msg)
}

var s2c_steal_log_list = [...]byte{38, 40} // 40
func NewS2cStealLogListProtoMsg(object *S2CStealLogListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_steal_log_list[:], "s2c_steal_log_list")

}

// target 非法
var ERR_STEAL_LOG_LIST_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 38, 42, 2} // 42-2

// start_time 非法
var ERR_STEAL_LOG_LIST_FAIL_INVALID_START_TIME = pbutil.StaticBuffer{3, 38, 42, 3} // 42-3

// 服务器错误
var ERR_STEAL_LOG_LIST_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 42, 1} // 42-1

func NewS2cCanStealListMsg(can_steal_id [][]byte) pbutil.Buffer {
	msg := &S2CCanStealListProto{
		CanStealId: can_steal_id,
	}
	return NewS2cCanStealListProtoMsg(msg)
}

var s2c_can_steal_list = [...]byte{38, 46} // 46
func NewS2cCanStealListProtoMsg(object *S2CCanStealListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_can_steal_list[:], "s2c_can_steal_list")

}

// 服务器错误
var ERR_CAN_STEAL_LIST_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 38, 49, 1} // 49-1
