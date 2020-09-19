package chat

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
	MODULE_ID = 13

	C2S_WORLD_CHAT = 1

	C2S_GUILD_CHAT = 4

	C2S_SELF_CHAT_WINDOW = 8

	C2S_CREATE_SELF_CHAT_WINDOW = 21

	C2S_REMOVE_CHAT_WINDOW = 10

	C2S_LIST_HISTORY_CHAT = 12

	C2S_SEND_CHAT = 14

	C2S_READ_CHAT_MSG = 18

	C2S_GET_HERO_CHAT_INFO = 25
)

var WORLD_CHAT_S2C = pbutil.StaticBuffer{2, 13, 2} // 2

func NewS2cWorldOtherChatMsg(id []byte, name string, head string, guild_flag string, text string, white_flag_guild_flag_name string) pbutil.Buffer {
	msg := &S2CWorldOtherChatProto{
		Id:        id,
		Name:      name,
		Head:      head,
		GuildFlag: guild_flag,
		Text:      text,
		WhiteFlagGuildFlagName: white_flag_guild_flag_name,
	}
	return NewS2cWorldOtherChatProtoMsg(msg)
}

var s2c_world_other_chat = [...]byte{13, 3} // 3
func NewS2cWorldOtherChatProtoMsg(object *S2CWorldOtherChatProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_world_other_chat[:], "s2c_world_other_chat")

}

var GUILD_CHAT_S2C = pbutil.StaticBuffer{2, 13, 5} // 5

// 你当前没有联盟
var ERR_GUILD_CHAT_FAIL_NOT_GUILD = pbutil.StaticBuffer{3, 13, 7, 1} // 7-1

func NewS2cGuildOtherChatMsg(id []byte, name string, head string, text string, white_flag_guild_flag_name string) pbutil.Buffer {
	msg := &S2CGuildOtherChatProto{
		Id:   id,
		Name: name,
		Head: head,
		Text: text,
		WhiteFlagGuildFlagName: white_flag_guild_flag_name,
	}
	return NewS2cGuildOtherChatProtoMsg(msg)
}

var s2c_guild_other_chat = [...]byte{13, 6} // 6
func NewS2cGuildOtherChatProtoMsg(object *S2CGuildOtherChatProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_guild_other_chat[:], "s2c_guild_other_chat")

}

func NewS2cSelfChatWindowMsg(sender [][]byte, unread []int32) pbutil.Buffer {
	msg := &S2CSelfChatWindowProto{
		Sender: sender,
		Unread: unread,
	}
	return NewS2cSelfChatWindowProtoMsg(msg)
}

func NewS2cSelfChatWindowMarshalMsg(sender [][]byte, unread []int32) pbutil.Buffer {
	msg := &S2CSelfChatWindowProto{
		Sender: sender,
		Unread: unread,
	}
	return NewS2cSelfChatWindowProtoMsg(msg)
}

var s2c_self_chat_window = [...]byte{13, 9} // 9
func NewS2cSelfChatWindowProtoMsg(object *S2CSelfChatWindowProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_self_chat_window[:], "s2c_self_chat_window")

}

func NewS2cCreateSelfChatWindowMsg(target []byte, set_up bool) pbutil.Buffer {
	msg := &S2CCreateSelfChatWindowProto{
		Target: target,
		SetUp:  set_up,
	}
	return NewS2cCreateSelfChatWindowProtoMsg(msg)
}

var s2c_create_self_chat_window = [...]byte{13, 22} // 22
func NewS2cCreateSelfChatWindowProtoMsg(object *S2CCreateSelfChatWindowProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_create_self_chat_window[:], "s2c_create_self_chat_window")

}

func NewS2cRemoveChatWindowMsg(chat_type int32, chat_target []byte) pbutil.Buffer {
	msg := &S2CRemoveChatWindowProto{
		ChatType:   chat_type,
		ChatTarget: chat_target,
	}
	return NewS2cRemoveChatWindowProtoMsg(msg)
}

var s2c_remove_chat_window = [...]byte{13, 11} // 11
func NewS2cRemoveChatWindowProtoMsg(object *S2CRemoveChatWindowProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_remove_chat_window[:], "s2c_remove_chat_window")

}

func NewS2cListHistoryChatMsg(chat_msg [][]byte) pbutil.Buffer {
	msg := &S2CListHistoryChatProto{
		ChatMsg: chat_msg,
	}
	return NewS2cListHistoryChatProtoMsg(msg)
}

func NewS2cListHistoryChatMarshalMsg(chat_msg [][]byte) pbutil.Buffer {
	msg := &S2CListHistoryChatProto{
		ChatMsg: chat_msg,
	}
	return NewS2cListHistoryChatProtoMsg(msg)
}

var s2c_list_history_chat = [...]byte{13, 13} // 13
func NewS2cListHistoryChatProtoMsg(object *S2CListHistoryChatProto) pbutil.Buffer {

	return util.NewSnappyCompressMsg(object, s2c_list_history_chat[:], "s2c_list_history_chat")

}

func NewS2cSendChatMsg(chat_id []byte, receiver []byte, replace_text string) pbutil.Buffer {
	msg := &S2CSendChatProto{
		ChatId:      chat_id,
		Receiver:    receiver,
		ReplaceText: replace_text,
	}
	return NewS2cSendChatProtoMsg(msg)
}

func NewS2cSendChatMarshalMsg(chat_id []byte, receiver marshaler, replace_text string) pbutil.Buffer {
	msg := &S2CSendChatProto{
		ChatId:      chat_id,
		Receiver:    safeMarshal(receiver),
		ReplaceText: replace_text,
	}
	return NewS2cSendChatProtoMsg(msg)
}

var s2c_send_chat = [...]byte{13, 15} // 15
func NewS2cSendChatProtoMsg(object *S2CSendChatProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_send_chat[:], "s2c_send_chat")

}

// 君主等级不足，不能使用世界聊天频道
var ERR_SEND_CHAT_FAIL_WORLD_CHAT_LEVEL = pbutil.StaticBuffer{3, 13, 16, 1} // 16-1

// 世界聊天太频繁了
var ERR_SEND_CHAT_FAIL_WORLD_CHAT_TOO_FAST = pbutil.StaticBuffer{3, 13, 16, 2} // 16-2

// 你没有联盟，不能发送联盟聊天
var ERR_SEND_CHAT_FAIL_NOT_IN_GUILD = pbutil.StaticBuffer{3, 13, 16, 3} // 16-3

// 无效的目标id
var ERR_SEND_CHAT_FAIL_INVALID_TARGET = pbutil.StaticBuffer{3, 13, 16, 4} // 16-4

// 说话文字太长
var ERR_SEND_CHAT_FAIL_TEXT_TOO_LONG = pbutil.StaticBuffer{3, 13, 16, 5} // 16-5

// 无效的分享坐标
var ERR_SEND_CHAT_FAIL_INVALID_POS = pbutil.StaticBuffer{3, 13, 16, 6} // 16-6

// 无效的分享回放连接
var ERR_SEND_CHAT_FAIL_INVALID_REPLAY = pbutil.StaticBuffer{3, 13, 16, 7} // 16-7

// 无效的分享战报id
var ERR_SEND_CHAT_FAIL_INVALID_REPORT = pbutil.StaticBuffer{3, 13, 16, 8} // 16-8

// 文字中包含敏感词
var ERR_SEND_CHAT_FAIL_SENSITIVE_WORDS = pbutil.StaticBuffer{3, 13, 16, 9} // 16-9

// 服务器忙，请稍后再试
var ERR_SEND_CHAT_FAIL_SERVER_ERROR = pbutil.StaticBuffer{3, 13, 16, 10} // 16-10

// 小喇叭不够
var ERR_SEND_CHAT_FAIL_BROADCAST_GOODS_NOT_ENOUGH = pbutil.StaticBuffer{3, 13, 16, 11} // 16-11

// 买小喇叭钱不够
var ERR_SEND_CHAT_FAIL_BROADCAST_MONEY_NOT_ENOUGH = pbutil.StaticBuffer{3, 13, 16, 12} // 16-12

// 名城战，没有参战
var ERR_SEND_CHAT_FAIL_MC_WAR_NOT_JOIN_FIGHT = pbutil.StaticBuffer{3, 13, 16, 14} // 16-14

// 不在名城战战斗阶段内
var ERR_SEND_CHAT_FAIL_MC_WAR_NOT_IN_FIGHT_STAGE = pbutil.StaticBuffer{3, 13, 16, 15} // 16-15

// 没有联盟权限
var ERR_SEND_CHAT_FAIL_GUILD_PERM_DENY = pbutil.StaticBuffer{3, 13, 16, 13} // 16-13

// 系统禁言
var ERR_SEND_CHAT_FAIL_BAN_CHAT = pbutil.StaticBuffer{3, 13, 16, 16} // 16-16

// 发送太频繁，请稍后再试
var ERR_SEND_CHAT_FAIL_TOO_FAST = pbutil.StaticBuffer{3, 13, 16, 17} // 16-17

// 通用聊天太频繁了
var ERR_SEND_CHAT_FAIL_CHAT_TOO_FAST = pbutil.StaticBuffer{3, 13, 16, 18} // 16-18

// 君主等级不足，不能使用私聊频道
var ERR_SEND_CHAT_FAIL_PRIVATE_CHAT_LEVEL = pbutil.StaticBuffer{3, 13, 16, 19} // 16-19

func NewS2cOtherSendChatMsg(chat_msg []byte) pbutil.Buffer {
	msg := &S2COtherSendChatProto{
		ChatMsg: chat_msg,
	}
	return NewS2cOtherSendChatProtoMsg(msg)
}

func NewS2cOtherSendChatMarshalMsg(chat_msg marshaler) pbutil.Buffer {
	msg := &S2COtherSendChatProto{
		ChatMsg: safeMarshal(chat_msg),
	}
	return NewS2cOtherSendChatProtoMsg(msg)
}

var s2c_other_send_chat = [...]byte{13, 17} // 17
func NewS2cOtherSendChatProtoMsg(object *S2COtherSendChatProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_other_send_chat[:], "s2c_other_send_chat")

}

func NewS2cReadChatMsgMsg(chat_type int32, chat_target []byte) pbutil.Buffer {
	msg := &S2CReadChatMsgProto{
		ChatType:   chat_type,
		ChatTarget: chat_target,
	}
	return NewS2cReadChatMsgProtoMsg(msg)
}

var s2c_read_chat_msg = [...]byte{13, 19} // 19
func NewS2cReadChatMsgProtoMsg(object *S2CReadChatMsgProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_read_chat_msg[:], "s2c_read_chat_msg")

}

var OFFLINE_CHAT_S2C = pbutil.StaticBuffer{2, 13, 24} // 24

func NewS2cGetHeroChatInfoMsg(id []byte, tower int32, jun_xian int32) pbutil.Buffer {
	msg := &S2CGetHeroChatInfoProto{
		Id:      id,
		Tower:   tower,
		JunXian: jun_xian,
	}
	return NewS2cGetHeroChatInfoProtoMsg(msg)
}

var s2c_get_hero_chat_info = [...]byte{13, 26} // 26
func NewS2cGetHeroChatInfoProtoMsg(object *S2CGetHeroChatInfoProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_get_hero_chat_info[:], "s2c_get_hero_chat_info")

}

func NewS2cBanChatMsg(end_time int32) pbutil.Buffer {
	msg := &S2CBanChatProto{
		EndTime: end_time,
	}
	return NewS2cBanChatProtoMsg(msg)
}

var s2c_ban_chat = [...]byte{13, 27} // 27
func NewS2cBanChatProtoMsg(object *S2CBanChatProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_ban_chat[:], "s2c_ban_chat")

}
