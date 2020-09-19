package activity

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
	MODULE_ID = 51

	C2S_COLLECT_COLLECTION = 1
)

func NewS2cNoticeActivityShowMsg(show []*shared_proto.ActiviyShowProto) pbutil.Buffer {
	msg := &S2CNoticeActivityShowProto{
		Show: show,
	}
	return NewS2cNoticeActivityShowProtoMsg(msg)
}

var s2c_notice_activity_show = [...]byte{51, 9} // 9
func NewS2cNoticeActivityShowProtoMsg(object *S2CNoticeActivityShowProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_notice_activity_show[:], "s2c_notice_activity_show")

}

func NewS2cCollectCollectionMsg(id int32, exchange_id int32) pbutil.Buffer {
	msg := &S2CCollectCollectionProto{
		Id:         id,
		ExchangeId: exchange_id,
	}
	return NewS2cCollectCollectionProtoMsg(msg)
}

var s2c_collect_collection = [...]byte{51, 2} // 2
func NewS2cCollectCollectionProtoMsg(object *S2CCollectCollectionProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_collection[:], "s2c_collect_collection")

}

// 不在活动期间内
var ERR_COLLECT_COLLECTION_FAIL_OUT_TIME = pbutil.StaticBuffer{3, 51, 3, 1} // 3-1

// 无效的id
var ERR_COLLECT_COLLECTION_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 51, 3, 2} // 3-2

// 兑换上限
var ERR_COLLECT_COLLECTION_FAIL_LIMIT = pbutil.StaticBuffer{3, 51, 3, 3} // 3-3

// 收集不足
var ERR_COLLECT_COLLECTION_FAIL_NOT_ENOUGH_COST = pbutil.StaticBuffer{3, 51, 3, 4} // 3-4

func NewS2cNoticeCollectionMsg(activity []*shared_proto.ActivityCollectionProto) pbutil.Buffer {
	msg := &S2CNoticeCollectionProto{
		Activity: activity,
	}
	return NewS2cNoticeCollectionProtoMsg(msg)
}

var s2c_notice_collection = [...]byte{51, 5} // 5
func NewS2cNoticeCollectionProtoMsg(object *S2CNoticeCollectionProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_notice_collection[:], "s2c_notice_collection")

}

func NewS2cNoticeCollectionCountsMsg(counts *shared_proto.HeroAllCollectionProto) pbutil.Buffer {
	msg := &S2CNoticeCollectionCountsProto{
		Counts: counts,
	}
	return NewS2cNoticeCollectionCountsProtoMsg(msg)
}

var s2c_notice_collection_counts = [...]byte{51, 6} // 6
func NewS2cNoticeCollectionCountsProtoMsg(object *S2CNoticeCollectionCountsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_notice_collection_counts[:], "s2c_notice_collection_counts")

}

func NewS2cNoticeTaskListModeMsg(activity []*shared_proto.ActivityTaskListModeProto) pbutil.Buffer {
	msg := &S2CNoticeTaskListModeProto{
		Activity: activity,
	}
	return NewS2cNoticeTaskListModeProtoMsg(msg)
}

var s2c_notice_task_list_mode = [...]byte{51, 7} // 7
func NewS2cNoticeTaskListModeProtoMsg(object *S2CNoticeTaskListModeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_notice_task_list_mode[:], "s2c_notice_task_list_mode")

}

func NewS2cNoticeTaskListModeProgressMsg(activity []*shared_proto.HeroActivityTaskListModeProto) pbutil.Buffer {
	msg := &S2CNoticeTaskListModeProgressProto{
		Activity: activity,
	}
	return NewS2cNoticeTaskListModeProgressProtoMsg(msg)
}

var s2c_notice_task_list_mode_progress = [...]byte{51, 8} // 8
func NewS2cNoticeTaskListModeProgressProtoMsg(object *S2CNoticeTaskListModeProgressProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_notice_task_list_mode_progress[:], "s2c_notice_task_list_mode_progress")

}
