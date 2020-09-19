package game2tss

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func UicJudgeUserInputNameV2(c client, ctx context.Context, msg string, if_replace bool) (*S2CUicJudgeUserInputNameV2Proto, error) {
	return UicJudgeUserInputNameV2Proto(c, ctx, &C2SUicJudgeUserInputNameV2Proto{

		Msg: msg,

		IfReplace: if_replace,
	})
}

func UicJudgeUserInputNameV2Proto(c client, ctx context.Context, proto *C2SUicJudgeUserInputNameV2Proto) (*S2CUicJudgeUserInputNameV2Proto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "uic_judge_user_input_name_v2 proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "uic_judge_user_input_name_v2", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "uic_judge_user_input_name_v2 fail")
		}

		s2c := &S2CUicJudgeUserInputNameV2Proto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "uic_judge_user_input_name_v2 s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}

func UicJudgeUserInputChatV2(c client, ctx context.Context, openid string, platid int32, world_id int64, msg_category int32, channel_id int64, client_ip int32, rold_id int64, role_name string, role_level int32, msg string, callback_data []byte, callback_addr string) (*S2CUicJudgeUserInputChatV2Proto, error) {
	return UicJudgeUserInputChatV2Proto(c, ctx, &C2SUicJudgeUserInputChatV2Proto{

		Openid: openid,

		Platid: platid,

		WorldId: world_id,

		MsgCategory: msg_category,

		ChannelId: channel_id,

		ClientIp: client_ip,

		RoldId: rold_id,

		RoleName: role_name,

		RoleLevel: role_level,

		Msg: msg,

		CallbackData: callback_data,

		CallbackAddr: callback_addr,
	})
}

func UicJudgeUserInputChatV2Proto(c client, ctx context.Context, proto *C2SUicJudgeUserInputChatV2Proto) (*S2CUicJudgeUserInputChatV2Proto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "uic_judge_user_input_chat_v2 proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "uic_judge_user_input_chat_v2", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "uic_judge_user_input_chat_v2 fail")
		}

		s2c := &S2CUicJudgeUserInputChatV2Proto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "uic_judge_user_input_chat_v2 s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}

func UicChatCallback(c client, ctx context.Context, openid string, platid int32, world_id int64, rold_id int64, msg_result_flag int32, replace_msg string, callback_data []byte) (*S2CUicChatCallbackProto, error) {
	return UicChatCallbackProto(c, ctx, &C2SUicChatCallbackProto{

		Openid: openid,

		Platid: platid,

		WorldId: world_id,

		RoldId: rold_id,

		MsgResultFlag: msg_result_flag,

		ReplaceMsg: replace_msg,

		CallbackData: callback_data,
	})
}

func UicChatCallbackProto(c client, ctx context.Context, proto *C2SUicChatCallbackProto) (*S2CUicChatCallbackProto, error) {
	if data, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "uic_chat_callback proto marshal fail")
	} else {
		result, err := c.HandleBytes(ctx, "uic_chat_callback", "", 0, data)
		if err != nil {
			return nil, errors.Wrapf(err, "uic_chat_callback fail")
		}

		s2c := &S2CUicChatCallbackProto{}
		if err := s2c.Unmarshal(result); err != nil {
			return nil, errors.Wrapf(err, "uic_chat_callback s2c.Unmarshal() fail")
		}
		return s2c, nil
	}
}
