package task

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
	MODULE_ID = 15

	C2S_COLLECT_TASK_PRIZE = 2

	C2S_COLLECT_TASK_BOX_PRIZE = 6

	C2S_COLLECT_BA_YE_STAGE_PRIZE = 9

	C2S_COLLECT_ACTIVE_DEGREE_PRIZE = 12

	C2S_COLLECT_ACHIEVE_STAR_PRIZE = 16

	C2S_CHANGE_SELECT_SHOW_ACHIEVE = 20

	C2S_COLLECT_BWZL_PRIZE = 23

	C2S_VIEW_OTHER_ACHIEVE_TASK_LIST = 26

	C2S_GET_TROOP_TITLE_FIGHT_AMOUNT = 29

	C2S_UPGRADE_TITLE = 31

	C2S_COMPLETE_BOOL_TASK = 35
)

func NewS2cUpdateTaskProgressMsg(id int32, progress int32, complete bool) pbutil.Buffer {
	msg := &S2CUpdateTaskProgressProto{
		Id:       id,
		Progress: progress,
		Complete: complete,
	}
	return NewS2cUpdateTaskProgressProtoMsg(msg)
}

var s2c_update_task_progress = [...]byte{15, 1} // 1
func NewS2cUpdateTaskProgressProtoMsg(object *S2CUpdateTaskProgressProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_task_progress[:], "s2c_update_task_progress")

}

func NewS2cCollectTaskPrizeMsg(id int32) pbutil.Buffer {
	msg := &S2CCollectTaskPrizeProto{
		Id: id,
	}
	return NewS2cCollectTaskPrizeProtoMsg(msg)
}

var s2c_collect_task_prize = [...]byte{15, 3} // 3
func NewS2cCollectTaskPrizeProtoMsg(object *S2CCollectTaskPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_task_prize[:], "s2c_collect_task_prize")

}

// 无效的任务id
var ERR_COLLECT_TASK_PRIZE_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 15, 4, 1} // 4-1

// 这个任务没有完成
var ERR_COLLECT_TASK_PRIZE_FAIL_NOT_COMPLETE = pbutil.StaticBuffer{3, 15, 4, 2} // 4-2

// 奖励已经领取
var ERR_COLLECT_TASK_PRIZE_FAIL_IS_COLLECTED = pbutil.StaticBuffer{3, 15, 4, 4} // 4-4

// 未到该天，无法领取(霸王之路任务)
var ERR_COLLECT_TASK_PRIZE_FAIL_CAN_NOT_COLLECT_THIS_DAY = pbutil.StaticBuffer{3, 15, 4, 5} // 4-5

// 服务器忙，请稍后再试
var ERR_COLLECT_TASK_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 15, 4, 3} // 4-3

func NewS2cNewTaskMsg(id int32, progress int32, complete bool, main bool) pbutil.Buffer {
	msg := &S2CNewTaskProto{
		Id:       id,
		Progress: progress,
		Complete: complete,
		Main:     main,
	}
	return NewS2cNewTaskProtoMsg(msg)
}

var s2c_new_task = [...]byte{15, 5} // 5
func NewS2cNewTaskProtoMsg(object *S2CNewTaskProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_new_task[:], "s2c_new_task")

}

func NewS2cCollectTaskBoxPrizeMsg(id int32) pbutil.Buffer {
	msg := &S2CCollectTaskBoxPrizeProto{
		Id: id,
	}
	return NewS2cCollectTaskBoxPrizeProtoMsg(msg)
}

var s2c_collect_task_box_prize = [...]byte{15, 7} // 7
func NewS2cCollectTaskBoxPrizeProtoMsg(object *S2CCollectTaskBoxPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_task_box_prize[:], "s2c_collect_task_box_prize")

}

// 无效的任务宝箱id
var ERR_COLLECT_TASK_BOX_PRIZE_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 15, 8, 1} // 8-1

// 这个任务完成个数没有达成
var ERR_COLLECT_TASK_BOX_PRIZE_FAIL_TASK_NOT_ENOUGH = pbutil.StaticBuffer{3, 15, 8, 2} // 8-2

// 这个宝箱已经领取过了
var ERR_COLLECT_TASK_BOX_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 15, 8, 3} // 8-3

// 前面还有宝箱未领取
var ERR_COLLECT_TASK_BOX_PRIZE_FAIL_PREV_NOT_COLLECTED = pbutil.StaticBuffer{3, 15, 8, 4} // 8-4

// 服务器忙，请稍后再试
var ERR_COLLECT_TASK_BOX_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 15, 8, 5} // 8-5

func NewS2cCollectBaYeStagePrizeMsg(stage []byte) pbutil.Buffer {
	msg := &S2CCollectBaYeStagePrizeProto{
		Stage: stage,
	}
	return NewS2cCollectBaYeStagePrizeProtoMsg(msg)
}

var s2c_collect_ba_ye_stage_prize = [...]byte{15, 10} // 10
func NewS2cCollectBaYeStagePrizeProtoMsg(object *S2CCollectBaYeStagePrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_ba_ye_stage_prize[:], "s2c_collect_ba_ye_stage_prize")

}

// 没有霸业目标
var ERR_COLLECT_BA_YE_STAGE_PRIZE_FAIL_NO_BA_YE_STAGE = pbutil.StaticBuffer{3, 15, 11, 3} // 11-3

// 有霸业目标任务没有完成
var ERR_COLLECT_BA_YE_STAGE_PRIZE_FAIL_HAS_TASK_NOT_COMPLETE = pbutil.StaticBuffer{3, 15, 11, 4} // 11-4

// 有霸业目标任务没有领取奖励
var ERR_COLLECT_BA_YE_STAGE_PRIZE_FAIL_HAS_TASK_NOT_COLLECT_PRIZE = pbutil.StaticBuffer{3, 15, 11, 5} // 11-5

// 服务器忙，请稍后再试
var ERR_COLLECT_BA_YE_STAGE_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 15, 11, 2} // 11-2

func NewS2cCollectActiveDegreePrizeMsg(collect_index int32) pbutil.Buffer {
	msg := &S2CCollectActiveDegreePrizeProto{
		CollectIndex: collect_index,
	}
	return NewS2cCollectActiveDegreePrizeProtoMsg(msg)
}

var s2c_collect_active_degree_prize = [...]byte{15, 13} // 13
func NewS2cCollectActiveDegreePrizeProtoMsg(object *S2CCollectActiveDegreePrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_active_degree_prize[:], "s2c_collect_active_degree_prize")

}

// 没找到该活跃度奖励
var ERR_COLLECT_ACTIVE_DEGREE_PRIZE_FAIL_PRIZE_NOT_FOUND = pbutil.StaticBuffer{3, 15, 14, 1} // 14-1

// 活跃度不够
var ERR_COLLECT_ACTIVE_DEGREE_PRIZE_FAIL_DEGREE_NOT_ENOUGH = pbutil.StaticBuffer{3, 15, 14, 2} // 14-2

// 奖励已经领取了
var ERR_COLLECT_ACTIVE_DEGREE_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 15, 14, 3} // 14-3

// 服务器忙，请稍后再试
var ERR_COLLECT_ACTIVE_DEGREE_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 15, 14, 4} // 14-4

func NewS2cCollectAchieveStarPrizeMsg(star_count int32) pbutil.Buffer {
	msg := &S2CCollectAchieveStarPrizeProto{
		StarCount: star_count,
	}
	return NewS2cCollectAchieveStarPrizeProtoMsg(msg)
}

var s2c_collect_achieve_star_prize = [...]byte{15, 17} // 17
func NewS2cCollectAchieveStarPrizeProtoMsg(object *S2CCollectAchieveStarPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_achieve_star_prize[:], "s2c_collect_achieve_star_prize")

}

// 没找到该成就星数奖励
var ERR_COLLECT_ACHIEVE_STAR_PRIZE_FAIL_PRIZE_NOT_FOUND = pbutil.StaticBuffer{3, 15, 18, 1} // 18-1

// 星数不够
var ERR_COLLECT_ACHIEVE_STAR_PRIZE_FAIL_STAR_NOT_ENOUGH = pbutil.StaticBuffer{3, 15, 18, 2} // 18-2

// 奖励已经领取了
var ERR_COLLECT_ACHIEVE_STAR_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 15, 18, 3} // 18-3

// 服务器忙，请稍后再试
var ERR_COLLECT_ACHIEVE_STAR_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 15, 18, 4} // 18-4

func NewS2cAchieveReachMsg(id int32, reach_time int32) pbutil.Buffer {
	msg := &S2CAchieveReachProto{
		Id:        id,
		ReachTime: reach_time,
	}
	return NewS2cAchieveReachProtoMsg(msg)
}

var s2c_achieve_reach = [...]byte{15, 19} // 19
func NewS2cAchieveReachProtoMsg(object *S2CAchieveReachProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_achieve_reach[:], "s2c_achieve_reach")

}

func NewS2cChangeSelectShowAchieveMsg(achieve_type int32, add_or_remove bool) pbutil.Buffer {
	msg := &S2CChangeSelectShowAchieveProto{
		AchieveType: achieve_type,
		AddOrRemove: add_or_remove,
	}
	return NewS2cChangeSelectShowAchieveProtoMsg(msg)
}

var s2c_change_select_show_achieve = [...]byte{15, 21} // 21
func NewS2cChangeSelectShowAchieveProtoMsg(object *S2CChangeSelectShowAchieveProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_change_select_show_achieve[:], "s2c_change_select_show_achieve")

}

// 成就类型没找到
var ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 15, 22, 1} // 22-1

// 成就未解锁
var ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_LOCKED = pbutil.StaticBuffer{3, 15, 22, 2} // 22-2

// 已经在展示了
var ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_SHOWING = pbutil.StaticBuffer{3, 15, 22, 3} // 22-3

// 展示数量太多，请先移除其他展示成就
var ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_SHOW_TO_MANY = pbutil.StaticBuffer{3, 15, 22, 4} // 22-4

// 没在展示中，无法移除
var ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_NOT_SHOWING = pbutil.StaticBuffer{3, 15, 22, 5} // 22-5

// 服务器忙，请稍后再试
var ERR_CHANGE_SELECT_SHOW_ACHIEVE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 15, 22, 6} // 22-6

func NewS2cRemoveTaskMsg(id int32) pbutil.Buffer {
	msg := &S2CRemoveTaskProto{
		Id: id,
	}
	return NewS2cRemoveTaskProtoMsg(msg)
}

var s2c_remove_task = [...]byte{15, 15} // 15
func NewS2cRemoveTaskProtoMsg(object *S2CRemoveTaskProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_task[:], "s2c_remove_task")

}

func NewS2cCollectBwzlPrizeMsg(complete_count int32) pbutil.Buffer {
	msg := &S2CCollectBwzlPrizeProto{
		CompleteCount: complete_count,
	}
	return NewS2cCollectBwzlPrizeProtoMsg(msg)
}

var s2c_collect_bwzl_prize = [...]byte{15, 24} // 24
func NewS2cCollectBwzlPrizeProtoMsg(object *S2CCollectBwzlPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_bwzl_prize[:], "s2c_collect_bwzl_prize")

}

// 奖励没找到
var ERR_COLLECT_BWZL_PRIZE_FAIL_NOT_FOUND = pbutil.StaticBuffer{3, 15, 25, 1} // 25-1

// 已经领取了
var ERR_COLLECT_BWZL_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 15, 25, 2} // 25-2

// 条件未达成
var ERR_COLLECT_BWZL_PRIZE_FAIL_NOT_REACH = pbutil.StaticBuffer{3, 15, 25, 3} // 25-3

// 服务器忙，请稍后再试
var ERR_COLLECT_BWZL_PRIZE_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 15, 25, 4} // 25-4

func NewS2cViewOtherAchieveTaskListMsg(id []byte, list *shared_proto.OtherAchieveTaskListProto) pbutil.Buffer {
	msg := &S2CViewOtherAchieveTaskListProto{
		Id:   id,
		List: list,
	}
	return NewS2cViewOtherAchieveTaskListProtoMsg(msg)
}

var s2c_view_other_achieve_task_list = [...]byte{15, 27} // 27
func NewS2cViewOtherAchieveTaskListProtoMsg(object *S2CViewOtherAchieveTaskListProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_view_other_achieve_task_list[:], "s2c_view_other_achieve_task_list")

}

// 非法的id
var ERR_VIEW_OTHER_ACHIEVE_TASK_LIST_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 15, 28, 1} // 28-1

// 请不要请求自己
var ERR_VIEW_OTHER_ACHIEVE_TASK_LIST_FAIL_CANT_REQUEST_SELF = pbutil.StaticBuffer{3, 15, 28, 2} // 28-2

// 服务器忙，请稍后再试
var ERR_VIEW_OTHER_ACHIEVE_TASK_LIST_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 15, 28, 3} // 28-3

func NewS2cGetTroopTitleFightAmountMsg(troop_index int32, title_id int32, fight_amount int32) pbutil.Buffer {
	msg := &S2CGetTroopTitleFightAmountProto{
		TroopIndex:  troop_index,
		TitleId:     title_id,
		FightAmount: fight_amount,
	}
	return NewS2cGetTroopTitleFightAmountProtoMsg(msg)
}

var s2c_get_troop_title_fight_amount = [...]byte{15, 30} // 30
func NewS2cGetTroopTitleFightAmountProtoMsg(object *S2CGetTroopTitleFightAmountProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_troop_title_fight_amount[:], "s2c_get_troop_title_fight_amount")

}

func NewS2cUpgradeTitleMsg(complete_title_id int32) pbutil.Buffer {
	msg := &S2CUpgradeTitleProto{
		CompleteTitleId: complete_title_id,
	}
	return NewS2cUpgradeTitleProtoMsg(msg)
}

var s2c_upgrade_title = [...]byte{15, 32} // 32
func NewS2cUpgradeTitleProtoMsg(object *S2CUpgradeTitleProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_upgrade_title[:], "s2c_upgrade_title")

}

// 任务还未完成，不能升级
var ERR_UPGRADE_TITLE_FAIL_TASK_NOT_COMPLETE = pbutil.StaticBuffer{3, 15, 33, 1} // 33-1

// 已达到最高级称号
var ERR_UPGRADE_TITLE_FAIL_MAX_LEVEL = pbutil.StaticBuffer{3, 15, 33, 2} // 33-2

// 贡品不足
var ERR_UPGRADE_TITLE_FAIL_NOT_ENOUGH_COST = pbutil.StaticBuffer{3, 15, 33, 3} // 33-3

func NewS2cCompleteBoolTaskMsg(bool_type int32) pbutil.Buffer {
	msg := &S2CCompleteBoolTaskProto{
		BoolType: bool_type,
	}
	return NewS2cCompleteBoolTaskProtoMsg(msg)
}

var s2c_complete_bool_task = [...]byte{15, 36} // 36
func NewS2cCompleteBoolTaskProtoMsg(object *S2CCompleteBoolTaskProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_complete_bool_task[:], "s2c_complete_bool_task")

}
