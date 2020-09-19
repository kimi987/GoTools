package red_packet

import (
	"github.com/lightpaw/pbutil"
)

// buy
var (
	ErrBuyFailInvalidDataId        = newMsgError("buy 红包 data.id 错误", ERR_BUY_FAIL_INVALID_DATA_ID)   // 3-1
	ErrBuyFailCostNotEnough        = newMsgError("buy 钱不够", ERR_BUY_FAIL_COST_NOT_ENOUGH)             // 3-2
	ErrBuyFailServerErr            = newMsgError("buy 服务器错误", ERR_BUY_FAIL_SERVER_ERR)                // 3-3
	ErrBuyFailRechangeYuanbaoLimit = newMsgError("buy 元宝赠送额度不够", ERR_BUY_FAIL_RECHANGE_YUANBAO_LIMIT) // 3-6
)

// create
var (
	ErrCreateFailInvalidDataId     = newMsgError("create 红包 data.id 错误", ERR_CREATE_FAIL_INVALID_DATA_ID) // 6-1
	ErrCreateFailInvalidChatType   = newMsgError("create 聊天类型错误", ERR_CREATE_FAIL_INVALID_CHAT_TYPE)      // 6-2
	ErrCreateFailNotBought         = newMsgError("create 没买", ERR_CREATE_FAIL_NOT_BOUGHT)                 // 6-4
	ErrCreateFailCountErr          = newMsgError("create 数量错误", ERR_CREATE_FAIL_COUNT_ERR)                // 6-5
	ErrCreateFailGuildLimit        = newMsgError("create 联盟人太少", ERR_CREATE_FAIL_GUILD_LIMIT)             // 6-6
	ErrCreateFailNoGuild           = newMsgError("create 发联盟红包，但没有联盟", ERR_CREATE_FAIL_NO_GUILD)          // 6-7
	ErrCreateFailServerErr         = newMsgError("create 服务器错误", ERR_CREATE_FAIL_SERVER_ERR)              // 6-3
	ErrCreateFailTextTooLong       = newMsgError("create 文字太长", ERR_CREATE_FAIL_TEXT_TOO_LONG)            // 6-8
	ErrCreateFailRechargeNotEnough = newMsgError("create 充钱不够", ERR_CREATE_FAIL_RECHARGE_NOT_ENOUGH)      // 6-9
)

// grab
var (
	ErrGrabFailInvalidId    = newMsgError("grab 红包id不存在或已经删掉了", ERR_GRAB_FAIL_INVALID_ID) // 9-1
	ErrGrabFailNotSameGuild = newMsgError("grab 不是同联盟", ERR_GRAB_FAIL_NOT_SAME_GUILD)     // 9-3
	ErrGrabFailExpired      = newMsgError("grab 过期了", ERR_GRAB_FAIL_EXPIRED)              // 9-4
	ErrGrabFailServerErr    = newMsgError("grab 服务器错误", ERR_GRAB_FAIL_SERVER_ERR)         // 9-2
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
