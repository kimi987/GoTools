package misc

import (
	"github.com/lightpaw/pbutil"
)

// disconect_reason
var (
	ErrDisconectReasonFailMustLogin = newMsgError("disconect_reason 没有登陆就发送了别的模块消息", ERR_DISCONECT_REASON_FAIL_MUST_LOGIN) // 5-1
	ErrDisconectReasonFailKick      = newMsgError("disconect_reason 你被顶号了", ERR_DISCONECT_REASON_FAIL_KICK)                // 5-2
	ErrDisconectReasonFailMoveBase  = newMsgError("disconect_reason 跨区域迁移失败", ERR_DISCONECT_REASON_FAIL_MOVE_BASE)         // 5-3
	ErrDisconectReasonFailLockFail  = newMsgError("disconect_reason Lock英雄失败", ERR_DISCONECT_REASON_FAIL_LOCK_FAIL)        // 5-4
	ErrDisconectReasonFailGm        = newMsgError("disconect_reason Gm模块踢人", ERR_DISCONECT_REASON_FAIL_GM)                 // 5-5
	ErrDisconectReasonFailClose     = newMsgError("disconect_reason 关服", ERR_DISCONECT_REASON_FAIL_CLOSE)                  // 5-6
	ErrDisconectReasonFailMsgRate   = newMsgError("disconect_reason 消息频率过快", ERR_DISCONECT_REASON_FAIL_MSG_RATE)           // 5-7
	ErrDisconectReasonFailHeartBeat = newMsgError("disconect_reason 心跳检查失败", ERR_DISCONECT_REASON_FAIL_HEART_BEAT)         // 5-8
)

// settings
var (
	ErrSettingsFailInvalidType = newMsgError("settings 非法的类型", ERR_SETTINGS_FAIL_INVALID_TYPE) // 23-1
)

// update_location
var (
	ErrUpdateLocationFailLocationError = newMsgError("update_location location 必须 >=0", ERR_UPDATE_LOCATION_FAIL_LOCATION_ERROR) // 33-1
)

// collect_charge_prize
var (
	ErrCollectChargePrizeFailInvalidId             = newMsgError("collect_charge_prize 无效的id", ERR_COLLECT_CHARGE_PRIZE_FAIL_INVALID_ID)              // 44-1
	ErrCollectChargePrizeFailCollected             = newMsgError("collect_charge_prize 奖励已领取", ERR_COLLECT_CHARGE_PRIZE_FAIL_COLLECTED)               // 44-2
	ErrCollectChargePrizeFailNotEnoughChargeAmount = newMsgError("collect_charge_prize 充值不够", ERR_COLLECT_CHARGE_PRIZE_FAIL_NOT_ENOUGH_CHARGE_AMOUNT) // 44-3
)

// collect_daily_bargain
var (
	ErrCollectDailyBargainFailInvalidId     = newMsgError("collect_daily_bargain 无效的id", ERR_COLLECT_DAILY_BARGAIN_FAIL_INVALID_ID)    // 48-1
	ErrCollectDailyBargainFailCannotCollect = newMsgError("collect_daily_bargain 无法领取", ERR_COLLECT_DAILY_BARGAIN_FAIL_CANNOT_COLLECT) // 48-2
)

// activate_duration_card
var (
	ErrActivateDurationCardFailInvalidId      = newMsgError("activate_duration_card 无效的id", ERR_ACTIVATE_DURATION_CARD_FAIL_INVALID_ID)     // 53-1
	ErrActivateDurationCardFailCannotActivate = newMsgError("activate_duration_card 无法购买", ERR_ACTIVATE_DURATION_CARD_FAIL_CANNOT_ACTIVATE) // 53-3
	ErrActivateDurationCardFailNoCharge       = newMsgError("activate_duration_card 充值未响应", ERR_ACTIVATE_DURATION_CARD_FAIL_NO_CHARGE)      // 53-2
)

// collect_duration_card_daily_prize
var (
	ErrCollectDurationCardDailyPrizeFailInvalidId = newMsgError("collect_duration_card_daily_prize 无效的id", ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_INVALID_ID) // 56-1
	ErrCollectDurationCardDailyPrizeFailNotActive = newMsgError("collect_duration_card_daily_prize 卡未激活", ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_NOT_ACTIVE)  // 56-4
	ErrCollectDurationCardDailyPrizeFailCollected = newMsgError("collect_duration_card_daily_prize 今日已领取", ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_COLLECTED)  // 56-2
	ErrCollectDurationCardDailyPrizeFailOverdue   = newMsgError("collect_duration_card_daily_prize 已过期，请续费", ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_OVERDUE)  // 56-3
)

// set_privacy_setting
var (
	ErrSetPrivacySettingFailInvalidType = newMsgError("set_privacy_setting 无效的类型", ERR_SET_PRIVACY_SETTING_FAIL_INVALID_TYPE) // 66-1
	ErrSetPrivacySettingFailDuplication = newMsgError("set_privacy_setting 重复设置", ERR_SET_PRIVACY_SETTING_FAIL_DUPLICATION)   // 66-2
)

// set_default_privacy_settings
var (
	ErrSetDefaultPrivacySettingsFailHasDefault = newMsgError("set_default_privacy_settings 已经恢复默认设置", ERR_SET_DEFAULT_PRIVACY_SETTINGS_FAIL_HAS_DEFAULT) // 68-1
)

// get_product_info
var (
	ErrGetProductInfoFailInvalidId = newMsgError("get_product_info 无效的商品id", ERR_GET_PRODUCT_INFO_FAIL_INVALID_ID) // 71-1
	ErrGetProductInfoFailCantBuy   = newMsgError("get_product_info 不能购买这个商品", ERR_GET_PRODUCT_INFO_FAIL_CANT_BUY)  // 71-2
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
