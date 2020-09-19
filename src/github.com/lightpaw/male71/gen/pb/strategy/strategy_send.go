package strategy

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
	MODULE_ID = 46

	C2S_USE_STRATAGEM = 1
)

func NewS2cUseStratagemMsg(id int32, name string, daily_used_times int32, next_useable_time int32, target_name string) pbutil.Buffer {
	msg := &S2CUseStratagemProto{
		Id:              id,
		Name:            name,
		DailyUsedTimes:  daily_used_times,
		NextUseableTime: next_useable_time,
		TargetName:      target_name,
	}
	return NewS2cUseStratagemProtoMsg(msg)
}

var s2c_use_stratagem = [...]byte{46, 2} // 2
func NewS2cUseStratagemProtoMsg(object *S2CUseStratagemProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_stratagem[:], "s2c_use_stratagem")

}

// 没有该计策
var ERR_USE_STRATAGEM_FAIL_INVALID_STRATAGEM_ID = pbutil.StaticBuffer{3, 46, 3, 1} // 3-1

// 计策未解锁
var ERR_USE_STRATAGEM_FAIL_LOCKED_STRATAGEM_ID = pbutil.StaticBuffer{3, 46, 3, 2} // 3-2

// 体力值不足
var ERR_USE_STRATAGEM_FAIL_SP_NOT_ENOUGH = pbutil.StaticBuffer{3, 46, 3, 3} // 3-3

// 计策冷却中
var ERR_USE_STRATAGEM_FAIL_STRATAGEM_CD = pbutil.StaticBuffer{3, 46, 3, 4} // 3-4

// 该计策今日使用上限
var ERR_USE_STRATAGEM_FAIL_TIMES_LIMIT = pbutil.StaticBuffer{3, 46, 3, 5} // 3-5

// 施计对象错误
var ERR_USE_STRATAGEM_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 46, 3, 6} // 3-6

// 今日对该目标施计上限
var ERR_USE_STRATAGEM_FAIL_TARGET_LIMIT = pbutil.StaticBuffer{3, 46, 3, 7} // 3-7

// 今日该目标中计已达上限
var ERR_USE_STRATAGEM_FAIL_TARGET_TRAPPED_LIMIT = pbutil.StaticBuffer{3, 46, 3, 8} // 3-8

// 该目标正中此计，请选择其他目标
var ERR_USE_STRATAGEM_FAIL_SAME_TRAPPED = pbutil.StaticBuffer{3, 46, 3, 9} // 3-9

// 消耗不够
var ERR_USE_STRATAGEM_FAIL_COST_NOT_ENOUGH = pbutil.StaticBuffer{3, 46, 3, 10} // 3-10

// 服务器错误
var ERR_USE_STRATAGEM_FAIL_SERVER_ERR = pbutil.StaticBuffer{3, 46, 3, 11} // 3-11

// 对方不符合计策条件
var ERR_USE_STRATAGEM_FAIL_TARGET_CANNOT_EFFECT = pbutil.StaticBuffer{3, 46, 3, 12} // 3-12

// 无效的配置id
var ERR_USE_STRATAGEM_FAIL_INVALID_DATA_ID = pbutil.StaticBuffer{3, 46, 3, 13} // 3-13

// 无效的坐标
var ERR_USE_STRATAGEM_FAIL_INVALID_POS = pbutil.StaticBuffer{3, 46, 3, 14} // 3-14

// 君主等级不足
var ERR_USE_STRATAGEM_FAIL_LEVEL_NOT_ENOUGH = pbutil.StaticBuffer{3, 46, 3, 15} // 3-15

// 召唤的殷墟还未消失
var ERR_USE_STRATAGEM_FAIL_EXIST_BAOZ = pbutil.StaticBuffer{3, 46, 3, 16} // 3-16

func NewS2cTrappedStratagemMsg(stratagem *shared_proto.TrappedStratagemProto) pbutil.Buffer {
	msg := &S2CTrappedStratagemProto{
		Stratagem: stratagem,
	}
	return NewS2cTrappedStratagemProtoMsg(msg)
}

var s2c_trapped_stratagem = [...]byte{46, 4} // 4
func NewS2cTrappedStratagemProtoMsg(object *S2CTrappedStratagemProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_trapped_stratagem[:], "s2c_trapped_stratagem")

}

func NewS2cUseStratagemFailMsg(id int32, name string, target_name string) pbutil.Buffer {
	msg := &S2CUseStratagemFailProto{
		Id:         id,
		Name:       name,
		TargetName: target_name,
	}
	return NewS2cUseStratagemFailProtoMsg(msg)
}

var s2c_use_stratagem_fail = [...]byte{46, 5} // 5
func NewS2cUseStratagemFailProtoMsg(object *S2CUseStratagemFailProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_use_stratagem_fail[:], "s2c_use_stratagem_fail")

}
