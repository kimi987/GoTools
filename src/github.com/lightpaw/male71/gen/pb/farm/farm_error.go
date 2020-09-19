package farm

import (
	"github.com/lightpaw/pbutil"
)

// plant
var (
	ErrPlantFailInvalidCube  = newMsgError("plant 无效地块。不是自己的", ERR_PLANT_FAIL_INVALID_CUBE) // 4-1
	ErrPlantFailInvalidResId = newMsgError("plant 资源 id 错误", ERR_PLANT_FAIL_INVALID_RES_ID) // 4-2
	ErrPlantFailNotIdleCube  = newMsgError("plant 不是空闲地块", ERR_PLANT_FAIL_NOT_IDLE_CUBE)    // 4-3
	ErrPlantFailServerErr    = newMsgError("plant 服务器错误", ERR_PLANT_FAIL_SERVER_ERR)        // 4-4
)

// harvest
var (
	ErrHarvestFailInvalidCube = newMsgError("harvest 无效地块。不是自己的、没有种植或者在冲突中", ERR_HARVEST_FAIL_INVALID_CUBE) // 7-1
	ErrHarvestFailServerErr   = newMsgError("harvest 服务器错误", ERR_HARVEST_FAIL_SERVER_ERR)                   // 7-2
	ErrHarvestFailNoOutput    = newMsgError("harvest 产量为0", ERR_HARVEST_FAIL_NO_OUTPUT)                     // 7-3
)

// change
var (
	ErrChangeFailInvalidCube = newMsgError("change 无效地块。不是自己的、没有种植或者在冲突中", ERR_CHANGE_FAIL_INVALID_CUBE) // 10-1
	ErrChangeFailServerErr   = newMsgError("change 服务器错误", ERR_CHANGE_FAIL_SERVER_ERR)                   // 10-2
)

// one_key_plant
var (
	ErrOneKeyPlantFailNoneIdleCube = newMsgError("one_key_plant 没有空白地块了", ERR_ONE_KEY_PLANT_FAIL_NONE_IDLE_CUBE) // 14-2
	ErrOneKeyPlantFailServerErr    = newMsgError("one_key_plant 服务器错误", ERR_ONE_KEY_PLANT_FAIL_SERVER_ERR)       // 14-3
	ErrOneKeyPlantFailInvalidCount = newMsgError("one_key_plant 数量都是0", ERR_ONE_KEY_PLANT_FAIL_INVALID_COUNT)    // 14-4
)

// one_key_harvest
var (
	ErrOneKeyHarvestFailNoneIdleCube = newMsgError("one_key_harvest 还没有成熟的地块", ERR_ONE_KEY_HARVEST_FAIL_NONE_IDLE_CUBE) // 30-1
	ErrOneKeyHarvestFailServerErr    = newMsgError("one_key_harvest 服务器错误", ERR_ONE_KEY_HARVEST_FAIL_SERVER_ERR)        // 30-2
	ErrOneKeyHarvestFailResTypeErr   = newMsgError("one_key_harvest 资源类型错误", ERR_ONE_KEY_HARVEST_FAIL_RES_TYPE_ERR)     // 30-3
)

// one_key_reset
var (
	ErrOneKeyResetFailServerErr = newMsgError("one_key_reset 服务器错误", ERR_ONE_KEY_RESET_FAIL_SERVER_ERR) // 54-1
)

// view_farm
var (
	ErrViewFarmFailInvalidTarget = newMsgError("view_farm 无效的 target", ERR_VIEW_FARM_FAIL_INVALID_TARGET) // 45-1
	ErrViewFarmFailServerErr     = newMsgError("view_farm 服务器错误", ERR_VIEW_FARM_FAIL_SERVER_ERR)          // 45-2
)

// steal
var (
	ErrStealFailInvalidTarget       = newMsgError("steal 无效的 target", ERR_STEAL_FAIL_INVALID_TARGET)         // 20-1
	ErrStealFailInvalidCube         = newMsgError("steal 无效的 cube, 没有这块地或没种东西", ERR_STEAL_FAIL_INVALID_CUBE) // 20-2
	ErrStealFailInProtectedDuration = newMsgError("steal 在保护期", ERR_STEAL_FAIL_IN_PROTECTED_DURATION)        // 20-3
	ErrStealFailNoStealTime         = newMsgError("steal 这块地偷菜次数用完了", ERR_STEAL_FAIL_NO_STEAL_TIME)          // 20-4
	ErrStealFailServerErr           = newMsgError("steal 服务器错误", ERR_STEAL_FAIL_SERVER_ERR)                  // 20-5
	ErrStealFailCubeAlreadyStealed  = newMsgError("steal 已经偷过了", ERR_STEAL_FAIL_CUBE_ALREADY_STEALED)        // 20-6
	ErrStealFailDailyStealFull      = newMsgError("steal 今天偷满了", ERR_STEAL_FAIL_DAILY_STEAL_FULL)            // 20-7
)

// one_key_steal
var (
	ErrOneKeyStealFailInvalidTarget  = newMsgError("one_key_steal 无效的 target", ERR_ONE_KEY_STEAL_FAIL_INVALID_TARGET) // 37-1
	ErrOneKeyStealFailServerErr      = newMsgError("one_key_steal 服务器错误", ERR_ONE_KEY_STEAL_FAIL_SERVER_ERR)          // 37-2
	ErrOneKeyStealFailNoCanStealCube = newMsgError("one_key_steal 没有可偷的地块", ERR_ONE_KEY_STEAL_FAIL_NO_CAN_STEAL_CUBE) // 37-3
	ErrOneKeyStealFailDailyStealFull = newMsgError("one_key_steal 今天偷满了", ERR_ONE_KEY_STEAL_FAIL_DAILY_STEAL_FULL)    // 37-4
)

// steal_log_list
var (
	ErrStealLogListFailInvalidTarget    = newMsgError("steal_log_list target 非法", ERR_STEAL_LOG_LIST_FAIL_INVALID_TARGET)         // 42-2
	ErrStealLogListFailInvalidStartTime = newMsgError("steal_log_list start_time 非法", ERR_STEAL_LOG_LIST_FAIL_INVALID_START_TIME) // 42-3
	ErrStealLogListFailServerErr        = newMsgError("steal_log_list 服务器错误", ERR_STEAL_LOG_LIST_FAIL_SERVER_ERR)                 // 42-1
)

// can_steal_list
var (
	ErrCanStealListFailServerErr = newMsgError("can_steal_list 服务器错误", ERR_CAN_STEAL_LIST_FAIL_SERVER_ERR) // 49-1
)

func newMsgError(msg string, buffer pbutil.StaticBuffer) *error_msg {
	return &error_msg{
		msg:  msg,
		buff: buffer,
	}
}

type error_msg struct {
	msg  string
	buff pbutil.Buffer
}

func (f *error_msg) Error() string         { return f.msg }
func (f *error_msg) ErrMsg() pbutil.Buffer { return f.buff }
