package random_event

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
	MODULE_ID = 45

	C2S_CHOOSE_OPTION = 1

	C2S_OPEN_EVENT = 4
)

func NewS2cChooseOptionMsg(pos_x int32, pos_y int32, option int32, success bool, prize []byte) pbutil.Buffer {
	msg := &S2CChooseOptionProto{
		PosX:    pos_x,
		PosY:    pos_y,
		Option:  option,
		Success: success,
		Prize:   prize,
	}
	return NewS2cChooseOptionProtoMsg(msg)
}

func NewS2cChooseOptionMarshalMsg(pos_x int32, pos_y int32, option int32, success bool, prize marshaler) pbutil.Buffer {
	msg := &S2CChooseOptionProto{
		PosX:    pos_x,
		PosY:    pos_y,
		Option:  option,
		Success: success,
		Prize:   safeMarshal(prize),
	}
	return NewS2cChooseOptionProtoMsg(msg)
}

var s2c_choose_option = [...]byte{45, 2} // 2
func NewS2cChooseOptionProtoMsg(object *S2CChooseOptionProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_choose_option[:], "s2c_choose_option")

}

// 失效的事件
var ERR_CHOOSE_OPTION_FAIL_INVALID_EVENT = pbutil.StaticBuffer{3, 45, 3, 1} // 3-1

// 无效的选项
var ERR_CHOOSE_OPTION_FAIL_INVALID_OPTION = pbutil.StaticBuffer{3, 45, 3, 2} // 3-2

// 选项消耗不足
var ERR_CHOOSE_OPTION_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 45, 3, 3} // 3-3

// 还未生成事件
var ERR_CHOOSE_OPTION_FAIL_NO_CATCH = pbutil.StaticBuffer{3, 45, 3, 4} // 3-4

func NewS2cOpenEventMsg(pos_x int32, pos_y int32, event_id int32, options []*shared_proto.EventOptionProto) pbutil.Buffer {
	msg := &S2COpenEventProto{
		PosX:    pos_x,
		PosY:    pos_y,
		EventId: event_id,
		Options: options,
	}
	return NewS2cOpenEventProtoMsg(msg)
}

var s2c_open_event = [...]byte{45, 5} // 5
func NewS2cOpenEventProtoMsg(object *S2COpenEventProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_open_event[:], "s2c_open_event")

}

// 失效的事件
var ERR_OPEN_EVENT_FAIL_INVALID_EVENT = pbutil.StaticBuffer{3, 45, 6, 1} // 6-1

func NewS2cNewEventMsg(arr_pos_x []int32, arr_pos_y []int32) pbutil.Buffer {
	msg := &S2CNewEventProto{
		ArrPosX: arr_pos_x,
		ArrPosY: arr_pos_y,
	}
	return NewS2cNewEventProtoMsg(msg)
}

var s2c_new_event = [...]byte{45, 8} // 8
func NewS2cNewEventProtoMsg(object *S2CNewEventProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_new_event[:], "s2c_new_event")

}

func NewS2cAddEventHandbookMsg(event_id int32) pbutil.Buffer {
	msg := &S2CAddEventHandbookProto{
		EventId: event_id,
	}
	return NewS2cAddEventHandbookProtoMsg(msg)
}

var s2c_add_event_handbook = [...]byte{45, 9} // 9
func NewS2cAddEventHandbookProtoMsg(object *S2CAddEventHandbookProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_add_event_handbook[:], "s2c_add_event_handbook")

}
