package strategy

import (
	"github.com/lightpaw/pbutil"
)

// use_stratagem
var (
	ErrUseStratagemFailInvalidStratagemId = newMsgError("use_stratagem 没有该计策", ERR_USE_STRATAGEM_FAIL_INVALID_STRATAGEM_ID)       // 3-1
	ErrUseStratagemFailLockedStratagemId  = newMsgError("use_stratagem 计策未解锁", ERR_USE_STRATAGEM_FAIL_LOCKED_STRATAGEM_ID)        // 3-2
	ErrUseStratagemFailSpNotEnough        = newMsgError("use_stratagem 体力值不足", ERR_USE_STRATAGEM_FAIL_SP_NOT_ENOUGH)              // 3-3
	ErrUseStratagemFailStratagemCd        = newMsgError("use_stratagem 计策冷却中", ERR_USE_STRATAGEM_FAIL_STRATAGEM_CD)               // 3-4
	ErrUseStratagemFailTimesLimit         = newMsgError("use_stratagem 该计策今日使用上限", ERR_USE_STRATAGEM_FAIL_TIMES_LIMIT)            // 3-5
	ErrUseStratagemFailInvalidTarget      = newMsgError("use_stratagem 施计对象错误", ERR_USE_STRATAGEM_FAIL_INVALID_TARGET)            // 3-6
	ErrUseStratagemFailTargetLimit        = newMsgError("use_stratagem 今日对该目标施计上限", ERR_USE_STRATAGEM_FAIL_TARGET_LIMIT)          // 3-7
	ErrUseStratagemFailTargetTrappedLimit = newMsgError("use_stratagem 今日该目标中计已达上限", ERR_USE_STRATAGEM_FAIL_TARGET_TRAPPED_LIMIT) // 3-8
	ErrUseStratagemFailSameTrapped        = newMsgError("use_stratagem 该目标正中此计，请选择其他目标", ERR_USE_STRATAGEM_FAIL_SAME_TRAPPED)     // 3-9
	ErrUseStratagemFailCostNotEnough      = newMsgError("use_stratagem 消耗不够", ERR_USE_STRATAGEM_FAIL_COST_NOT_ENOUGH)             // 3-10
	ErrUseStratagemFailServerErr          = newMsgError("use_stratagem 服务器错误", ERR_USE_STRATAGEM_FAIL_SERVER_ERR)                 // 3-11
	ErrUseStratagemFailTargetCannotEffect = newMsgError("use_stratagem 对方不符合计策条件", ERR_USE_STRATAGEM_FAIL_TARGET_CANNOT_EFFECT)   // 3-12
	ErrUseStratagemFailInvalidDataId      = newMsgError("use_stratagem 无效的配置id", ERR_USE_STRATAGEM_FAIL_INVALID_DATA_ID)          // 3-13
	ErrUseStratagemFailInvalidPos         = newMsgError("use_stratagem 无效的坐标", ERR_USE_STRATAGEM_FAIL_INVALID_POS)                // 3-14
	ErrUseStratagemFailLevelNotEnough     = newMsgError("use_stratagem 君主等级不足", ERR_USE_STRATAGEM_FAIL_LEVEL_NOT_ENOUGH)          // 3-15
	ErrUseStratagemFailExistBaoz          = newMsgError("use_stratagem 召唤的殷墟还未消失", ERR_USE_STRATAGEM_FAIL_EXIST_BAOZ)             // 3-16
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
