package game2tss

import (
	"github.com/lightpaw/rpc7"
	"golang.org/x/net/context"
)

func NewUicJudgeUserInputNameV2Handler(f func(r *C2SUicJudgeUserInputNameV2Proto) (*S2CUicJudgeUserInputNameV2Proto, error)) (string, rpc7.HandlerFunc) {
	return "uic_judge_user_input_name_v2", func(ctx context.Context, r *rpc7.RpcRequest) (*rpc7.RpcResponse, error) {
		c2s := &C2SUicJudgeUserInputNameV2Proto{}
		if err := c2s.Unmarshal(r.Proto); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_C2S_PROTO, Msg: "C2SUicJudgeUserInputNameV2.Unmarshal() err: " + err.Error()}, nil
		}

		if s2c, err := f(c2s); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_HANDLE_ERROR, Msg: "UicJudgeUserInputNameV2(r) err: " + err.Error()}, nil
		} else {
			if b, err := s2c.Marshal(); err != nil {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_S2C_PROTO, Msg: "S2CUicJudgeUserInputNameV2.Marshal() err: " + err.Error()}, nil
			} else {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_SUCCESS, Proto: b}, nil
			}
		}
	}
}

func NewUicJudgeUserInputChatV2Handler(f func(r *C2SUicJudgeUserInputChatV2Proto) (*S2CUicJudgeUserInputChatV2Proto, error)) (string, rpc7.HandlerFunc) {
	return "uic_judge_user_input_chat_v2", func(ctx context.Context, r *rpc7.RpcRequest) (*rpc7.RpcResponse, error) {
		c2s := &C2SUicJudgeUserInputChatV2Proto{}
		if err := c2s.Unmarshal(r.Proto); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_C2S_PROTO, Msg: "C2SUicJudgeUserInputChatV2.Unmarshal() err: " + err.Error()}, nil
		}

		if s2c, err := f(c2s); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_HANDLE_ERROR, Msg: "UicJudgeUserInputChatV2(r) err: " + err.Error()}, nil
		} else {
			if b, err := s2c.Marshal(); err != nil {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_S2C_PROTO, Msg: "S2CUicJudgeUserInputChatV2.Marshal() err: " + err.Error()}, nil
			} else {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_SUCCESS, Proto: b}, nil
			}
		}
	}
}

func NewUicChatCallbackHandler(f func(r *C2SUicChatCallbackProto) (*S2CUicChatCallbackProto, error)) (string, rpc7.HandlerFunc) {
	return "uic_chat_callback", func(ctx context.Context, r *rpc7.RpcRequest) (*rpc7.RpcResponse, error) {
		c2s := &C2SUicChatCallbackProto{}
		if err := c2s.Unmarshal(r.Proto); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_C2S_PROTO, Msg: "C2SUicChatCallback.Unmarshal() err: " + err.Error()}, nil
		}

		if s2c, err := f(c2s); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_HANDLE_ERROR, Msg: "UicChatCallback(r) err: " + err.Error()}, nil
		} else {
			if b, err := s2c.Marshal(); err != nil {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_S2C_PROTO, Msg: "S2CUicChatCallback.Marshal() err: " + err.Error()}, nil
			} else {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_SUCCESS, Proto: b}, nil
			}
		}
	}
}
