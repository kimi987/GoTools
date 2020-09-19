package misc

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
	MODULE_ID = 5

	C2S_HEART_BEAT = 1

	C2S_BACKGROUND_HEART_BEAT = 35

	C2S_BACKGROUND_WEAKUP = 36

	C2S_CONFIG = 3

	C2S_CONFIGLUA = 76

	C2S_CLIENT_LOG = 7

	C2S_SYNC_TIME = 8

	C2S_BLOCK = 10

	C2S_PING = 15

	C2S_CLIENT_VERSION = 17

	C2S_UPDATE_PF_TOKEN = 19

	C2S_SETTINGS = 21

	C2S_SETTINGS_TO_DEFAULT = 24

	C2S_UPDATE_LOCATION = 31

	C2S_COLLECT_CHARGE_PRIZE = 42

	C2S_COLLECT_DAILY_BARGAIN = 46

	C2S_ACTIVATE_DURATION_CARD = 51

	C2S_COLLECT_DURATION_CARD_DAILY_PRIZE = 54

	C2S_SET_PRIVACY_SETTING = 64

	C2S_SET_DEFAULT_PRIVACY_SETTINGS = 67

	C2S_GET_PRODUCT_INFO = 69
)

func NewS2cConfigMsg(version string, config []byte) pbutil.Buffer {
	msg := &S2CConfigProto{
		Version: version,
		Config:  config,
	}
	return NewS2cConfigProtoMsg(msg)
}

func NewS2cConfigMarshalMsg(version string, config marshaler) pbutil.Buffer {
	msg := &S2CConfigProto{
		Version: version,
		Config:  safeMarshal(config),
	}
	return NewS2cConfigProtoMsg(msg)
}

var s2c_config = [...]byte{5, 4} // 4
func NewS2cConfigProtoMsg(object *S2CConfigProto) pbutil.Buffer {

	return util.NewGzipCompressMsg(object, s2c_config[:], "s2c_config")

}

func NewS2cConfigluaMsg(version string, config []byte) pbutil.Buffer {
	msg := &S2CConfigluaProto{
		Version: version,
		Config:  config,
	}
	return NewS2cConfigluaProtoMsg(msg)
}

var s2c_configlua = [...]byte{5, 77} // 77
func NewS2cConfigluaProtoMsg(object *S2CConfigluaProto) pbutil.Buffer {

	return util.NewGzipCompressMsg(object, s2c_configlua[:], "s2c_configlua")

}

// 没有登陆就发送了别的模块消息
var ERR_DISCONECT_REASON_FAIL_MUST_LOGIN = pbutil.StaticBuffer{3, 5, 5, 1} // 5-1

// 你被顶号了
var ERR_DISCONECT_REASON_FAIL_KICK = pbutil.StaticBuffer{3, 5, 5, 2} // 5-2

// 跨区域迁移失败
var ERR_DISCONECT_REASON_FAIL_MOVE_BASE = pbutil.StaticBuffer{3, 5, 5, 3} // 5-3

// Lock英雄失败
var ERR_DISCONECT_REASON_FAIL_LOCK_FAIL = pbutil.StaticBuffer{3, 5, 5, 4} // 5-4

// Gm模块踢人
var ERR_DISCONECT_REASON_FAIL_GM = pbutil.StaticBuffer{3, 5, 5, 5} // 5-5

// 关服
var ERR_DISCONECT_REASON_FAIL_CLOSE = pbutil.StaticBuffer{3, 5, 5, 6} // 5-6

// 消息频率过快
var ERR_DISCONECT_REASON_FAIL_MSG_RATE = pbutil.StaticBuffer{3, 5, 5, 7} // 5-7

// 心跳检查失败
var ERR_DISCONECT_REASON_FAIL_HEART_BEAT = pbutil.StaticBuffer{3, 5, 5, 8} // 5-8

var RESET_DAILY_S2C = pbutil.StaticBuffer{2, 5, 6} // 6

func NewS2cSyncTimeMsg(client_time int32, server_time int32) pbutil.Buffer {
	msg := &S2CSyncTimeProto{
		ClientTime: client_time,
		ServerTime: server_time,
	}
	return NewS2cSyncTimeProtoMsg(msg)
}

var s2c_sync_time = [...]byte{5, 9} // 9
func NewS2cSyncTimeProtoMsg(object *S2CSyncTimeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_sync_time[:], "s2c_sync_time")

}

func NewS2cBlockMsg(data []byte) pbutil.Buffer {
	msg := &S2CBlockProto{
		Data: data,
	}
	return NewS2cBlockProtoMsg(msg)
}

func NewS2cBlockMarshalMsg(data marshaler) pbutil.Buffer {
	msg := &S2CBlockProto{
		Data: safeMarshal(data),
	}
	return NewS2cBlockProtoMsg(msg)
}

var s2c_block = [...]byte{5, 11} // 11
func NewS2cBlockProtoMsg(object *S2CBlockProto) pbutil.Buffer {

	return util.NewGzipCompressMsg(object, s2c_block[:], "s2c_block")

}

func NewS2cOpenFunctionMsg(function_type int32) pbutil.Buffer {
	msg := &S2COpenFunctionProto{
		FunctionType: function_type,
	}
	return NewS2cOpenFunctionProtoMsg(msg)
}

var s2c_open_function = [...]byte{5, 12} // 12
func NewS2cOpenFunctionProtoMsg(object *S2COpenFunctionProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_open_function[:], "s2c_open_function")

}

func NewS2cOpenMultiFunctionMsg(function_type []int32) pbutil.Buffer {
	msg := &S2COpenMultiFunctionProto{
		FunctionType: function_type,
	}
	return NewS2cOpenMultiFunctionProtoMsg(msg)
}

var s2c_open_multi_function = [...]byte{5, 34} // 34
func NewS2cOpenMultiFunctionProtoMsg(object *S2COpenMultiFunctionProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_open_multi_function[:], "s2c_open_multi_function")

}

func NewS2cSetHeroBoolMsg(bool_type int32) pbutil.Buffer {
	msg := &S2CSetHeroBoolProto{
		BoolType: bool_type,
	}
	return NewS2cSetHeroBoolProtoMsg(msg)
}

var s2c_set_hero_bool = [...]byte{5, 14} // 14
func NewS2cSetHeroBoolProtoMsg(object *S2CSetHeroBoolProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_hero_bool[:], "s2c_set_hero_bool")

}

func NewS2cResetHeroBoolMsg(bool_type int32) pbutil.Buffer {
	msg := &S2CResetHeroBoolProto{
		BoolType: bool_type,
	}
	return NewS2cResetHeroBoolProtoMsg(msg)
}

var s2c_reset_hero_bool = [...]byte{5, 30} // 30
func NewS2cResetHeroBoolProtoMsg(object *S2CResetHeroBoolProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_reset_hero_bool[:], "s2c_reset_hero_bool")

}

func NewS2cScreenShowWordsMsg(json string) pbutil.Buffer {
	msg := &S2CScreenShowWordsProto{
		Json: json,
	}
	return NewS2cScreenShowWordsProtoMsg(msg)
}

var s2c_screen_show_words = [...]byte{5, 13} // 13
func NewS2cScreenShowWordsProtoMsg(object *S2CScreenShowWordsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_screen_show_words[:], "s2c_screen_show_words")

}

var PING_S2C = pbutil.StaticBuffer{2, 5, 16} // 16

func NewS2cClientVersionMsg(v string, os string, t string) pbutil.Buffer {
	msg := &S2CClientVersionProto{
		V:  v,
		Os: os,
		T:  t,
	}
	return NewS2cClientVersionProtoMsg(msg)
}

var s2c_client_version = [...]byte{5, 18} // 18
func NewS2cClientVersionProtoMsg(object *S2CClientVersionProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_client_version[:], "s2c_client_version")

}

var UPDATE_PF_TOKEN_S2C = pbutil.StaticBuffer{2, 5, 20} // 20

func NewS2cSettingsMsg(setting_type shared_proto.SettingType, open bool) pbutil.Buffer {
	msg := &S2CSettingsProto{
		SettingType: setting_type,
		Open:        open,
	}
	return NewS2cSettingsProtoMsg(msg)
}

var s2c_settings = [...]byte{5, 22} // 22
func NewS2cSettingsProtoMsg(object *S2CSettingsProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_settings[:], "s2c_settings")

}

// 非法的类型
var ERR_SETTINGS_FAIL_INVALID_TYPE = pbutil.StaticBuffer{3, 5, 23, 1} // 23-1

func NewS2cSettingsToDefaultMsg(setting_type *shared_proto.HeroSettingsProto) pbutil.Buffer {
	msg := &S2CSettingsToDefaultProto{
		SettingType: setting_type,
	}
	return NewS2cSettingsToDefaultProtoMsg(msg)
}

var s2c_settings_to_default = [...]byte{5, 25} // 25
func NewS2cSettingsToDefaultProtoMsg(object *S2CSettingsToDefaultProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_settings_to_default[:], "s2c_settings_to_default")

}

func NewS2cHeroBroadcastMsg(text string, name string, guild_flag string) pbutil.Buffer {
	msg := &S2CHeroBroadcastProto{
		Text:      text,
		Name:      name,
		GuildFlag: guild_flag,
	}
	return NewS2cHeroBroadcastProtoMsg(msg)
}

var s2c_hero_broadcast = [...]byte{5, 26} // 26
func NewS2cHeroBroadcastProtoMsg(object *S2CHeroBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_hero_broadcast[:], "s2c_hero_broadcast")

}

func NewS2cSysTimingBroadcastMsg(text string) pbutil.Buffer {
	msg := &S2CSysTimingBroadcastProto{
		Text: text,
	}
	return NewS2cSysTimingBroadcastProtoMsg(msg)
}

var s2c_sys_timing_broadcast = [...]byte{5, 28} // 28
func NewS2cSysTimingBroadcastProtoMsg(object *S2CSysTimingBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_sys_timing_broadcast[:], "s2c_sys_timing_broadcast")

}

func NewS2cSysBroadcastMsg(text string) pbutil.Buffer {
	msg := &S2CSysBroadcastProto{
		Text: text,
	}
	return NewS2cSysBroadcastProtoMsg(msg)
}

var s2c_sys_broadcast = [...]byte{5, 27} // 27
func NewS2cSysBroadcastProtoMsg(object *S2CSysBroadcastProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_sys_broadcast[:], "s2c_sys_broadcast")

}

func NewS2cUpdateLocationMsg(location int32) pbutil.Buffer {
	msg := &S2CUpdateLocationProto{
		Location: location,
	}
	return NewS2cUpdateLocationProtoMsg(msg)
}

var s2c_update_location = [...]byte{5, 32} // 32
func NewS2cUpdateLocationProtoMsg(object *S2CUpdateLocationProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_location[:], "s2c_update_location")

}

// location 必须 >=0
var ERR_UPDATE_LOCATION_FAIL_LOCATION_ERROR = pbutil.StaticBuffer{3, 5, 33, 1} // 33-1

func NewS2cCollectChargePrizeMsg(id int32, prize []byte) pbutil.Buffer {
	msg := &S2CCollectChargePrizeProto{
		Id:    id,
		Prize: prize,
	}
	return NewS2cCollectChargePrizeProtoMsg(msg)
}

var s2c_collect_charge_prize = [...]byte{5, 43} // 43
func NewS2cCollectChargePrizeProtoMsg(object *S2CCollectChargePrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_charge_prize[:], "s2c_collect_charge_prize")

}

// 无效的id
var ERR_COLLECT_CHARGE_PRIZE_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 5, 44, 1} // 44-1

// 奖励已领取
var ERR_COLLECT_CHARGE_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 5, 44, 2} // 44-2

// 充值不够
var ERR_COLLECT_CHARGE_PRIZE_FAIL_NOT_ENOUGH_CHARGE_AMOUNT = pbutil.StaticBuffer{3, 5, 44, 3} // 44-3

func NewS2cUpdateChargeAmountMsg(amount int32) pbutil.Buffer {
	msg := &S2CUpdateChargeAmountProto{
		Amount: amount,
	}
	return NewS2cUpdateChargeAmountProtoMsg(msg)
}

var s2c_update_charge_amount = [...]byte{5, 45} // 45
func NewS2cUpdateChargeAmountProtoMsg(object *S2CUpdateChargeAmountProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_charge_amount[:], "s2c_update_charge_amount")

}

func NewS2cCollectDailyBargainMsg(id int32, times int32, prize []byte) pbutil.Buffer {
	msg := &S2CCollectDailyBargainProto{
		Id:    id,
		Times: times,
		Prize: prize,
	}
	return NewS2cCollectDailyBargainProtoMsg(msg)
}

var s2c_collect_daily_bargain = [...]byte{5, 47} // 47
func NewS2cCollectDailyBargainProtoMsg(object *S2CCollectDailyBargainProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_daily_bargain[:], "s2c_collect_daily_bargain")

}

// 无效的id
var ERR_COLLECT_DAILY_BARGAIN_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 5, 48, 1} // 48-1

// 无法领取
var ERR_COLLECT_DAILY_BARGAIN_FAIL_CANNOT_COLLECT = pbutil.StaticBuffer{3, 5, 48, 2} // 48-2

func NewS2cActivateDurationCardMsg(id int32, end_time int32, prize []byte) pbutil.Buffer {
	msg := &S2CActivateDurationCardProto{
		Id:      id,
		EndTime: end_time,
		Prize:   prize,
	}
	return NewS2cActivateDurationCardProtoMsg(msg)
}

var s2c_activate_duration_card = [...]byte{5, 52} // 52
func NewS2cActivateDurationCardProtoMsg(object *S2CActivateDurationCardProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_activate_duration_card[:], "s2c_activate_duration_card")

}

// 无效的id
var ERR_ACTIVATE_DURATION_CARD_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 5, 53, 1} // 53-1

// 无法购买
var ERR_ACTIVATE_DURATION_CARD_FAIL_CANNOT_ACTIVATE = pbutil.StaticBuffer{3, 5, 53, 3} // 53-3

// 充值未响应
var ERR_ACTIVATE_DURATION_CARD_FAIL_NO_CHARGE = pbutil.StaticBuffer{3, 5, 53, 2} // 53-2

func NewS2cCollectDurationCardDailyPrizeMsg(id int32, prize []byte) pbutil.Buffer {
	msg := &S2CCollectDurationCardDailyPrizeProto{
		Id:    id,
		Prize: prize,
	}
	return NewS2cCollectDurationCardDailyPrizeProtoMsg(msg)
}

var s2c_collect_duration_card_daily_prize = [...]byte{5, 55} // 55
func NewS2cCollectDurationCardDailyPrizeProtoMsg(object *S2CCollectDurationCardDailyPrizeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_collect_duration_card_daily_prize[:], "s2c_collect_duration_card_daily_prize")

}

// 无效的id
var ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 5, 56, 1} // 56-1

// 卡未激活
var ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_NOT_ACTIVE = pbutil.StaticBuffer{3, 5, 56, 4} // 56-4

// 今日已领取
var ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_COLLECTED = pbutil.StaticBuffer{3, 5, 56, 2} // 56-2

// 已过期，请续费
var ERR_COLLECT_DURATION_CARD_DAILY_PRIZE_FAIL_OVERDUE = pbutil.StaticBuffer{3, 5, 56, 3} // 56-3

var RESET_DAILY_ZERO_S2C = pbutil.StaticBuffer{2, 5, 62} // 62

var RESET_DAILY_MC_S2C = pbutil.StaticBuffer{2, 5, 72} // 72

func NewS2cSetPrivacySettingMsg(setting_type shared_proto.PrivacySettingType, open_or_close bool) pbutil.Buffer {
	msg := &S2CSetPrivacySettingProto{
		SettingType: setting_type,
		OpenOrClose: open_or_close,
	}
	return NewS2cSetPrivacySettingProtoMsg(msg)
}

var s2c_set_privacy_setting = [...]byte{5, 65} // 65
func NewS2cSetPrivacySettingProtoMsg(object *S2CSetPrivacySettingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_set_privacy_setting[:], "s2c_set_privacy_setting")

}

// 无效的类型
var ERR_SET_PRIVACY_SETTING_FAIL_INVALID_TYPE = pbutil.StaticBuffer{3, 5, 66, 1} // 66-1

// 重复设置
var ERR_SET_PRIVACY_SETTING_FAIL_DUPLICATION = pbutil.StaticBuffer{3, 5, 66, 2} // 66-2

// 已经恢复默认设置
var ERR_SET_DEFAULT_PRIVACY_SETTINGS_FAIL_HAS_DEFAULT = pbutil.StaticBuffer{3, 5, 68, 1} // 68-1

func NewS2cGetProductInfoMsg(id int32, product_id string, product_name string, cp_order_id string, money int32, gold int32, ext string, is_debug bool) pbutil.Buffer {
	msg := &S2CGetProductInfoProto{
		Id:          id,
		ProductId:   product_id,
		ProductName: product_name,
		CpOrderId:   cp_order_id,
		Money:       money,
		Gold:        gold,
		Ext:         ext,
		IsDebug:     is_debug,
	}
	return NewS2cGetProductInfoProtoMsg(msg)
}

var s2c_get_product_info = [...]byte{5, 70} // 70
func NewS2cGetProductInfoProtoMsg(object *S2CGetProductInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_product_info[:], "s2c_get_product_info")

}

// 无效的商品id
var ERR_GET_PRODUCT_INFO_FAIL_INVALID_ID = pbutil.StaticBuffer{3, 5, 71, 1} // 71-1

// 不能购买这个商品
var ERR_GET_PRODUCT_INFO_FAIL_CANT_BUY = pbutil.StaticBuffer{3, 5, 71, 2} // 71-2

var RESET_WEEKLY_S2C = pbutil.StaticBuffer{2, 5, 73} // 73

func NewS2cUpdateFirstRechargeMsg(id int32) pbutil.Buffer {
	msg := &S2CUpdateFirstRechargeProto{
		Id: id,
	}
	return NewS2cUpdateFirstRechargeProtoMsg(msg)
}

var s2c_update_first_recharge = [...]byte{5, 74} // 74
func NewS2cUpdateFirstRechargeProtoMsg(object *S2CUpdateFirstRechargeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_first_recharge[:], "s2c_update_first_recharge")

}

func NewS2cUpdateBuffNoticeMsg(group_id int32, buff *shared_proto.BuffInfoProto) pbutil.Buffer {
	msg := &S2CUpdateBuffNoticeProto{
		GroupId: group_id,
		Buff:    buff,
	}
	return NewS2cUpdateBuffNoticeProtoMsg(msg)
}

var s2c_update_buff_notice = [...]byte{5, 75} // 75
func NewS2cUpdateBuffNoticeProtoMsg(object *S2CUpdateBuffNoticeProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_update_buff_notice[:], "s2c_update_buff_notice")

}
