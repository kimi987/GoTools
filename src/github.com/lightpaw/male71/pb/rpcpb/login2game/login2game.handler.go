package login2game

import (
	"github.com/lightpaw/rpc7"
	"golang.org/x/net/context"
)

func NewCompleteQuestionnaireHandler(f func(r *C2SCompleteQuestionnaireProto) (*S2CCompleteQuestionnaireProto, error)) (string, rpc7.HandlerFunc) {
	return "complete_questionnaire", func(ctx context.Context, r *rpc7.RpcRequest) (*rpc7.RpcResponse, error) {
		c2s := &C2SCompleteQuestionnaireProto{}
		if err := c2s.Unmarshal(r.Proto); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_C2S_PROTO, Msg: "C2SCompleteQuestionnaire.Unmarshal() err: " + err.Error()}, nil
		}

		if s2c, err := f(c2s); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_HANDLE_ERROR, Msg: "CompleteQuestionnaire(r) err: " + err.Error()}, nil
		} else {
			if b, err := s2c.Marshal(); err != nil {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_S2C_PROTO, Msg: "S2CCompleteQuestionnaire.Marshal() err: " + err.Error()}, nil
			} else {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_SUCCESS, Proto: b}, nil
			}
		}
	}
}

func NewBuyProductHandler(f func(r *C2SBuyProductProto) (*S2CBuyProductProto, error)) (string, rpc7.HandlerFunc) {
	return "buy_product", func(ctx context.Context, r *rpc7.RpcRequest) (*rpc7.RpcResponse, error) {
		c2s := &C2SBuyProductProto{}
		if err := c2s.Unmarshal(r.Proto); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_C2S_PROTO, Msg: "C2SBuyProduct.Unmarshal() err: " + err.Error()}, nil
		}

		if s2c, err := f(c2s); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_HANDLE_ERROR, Msg: "BuyProduct(r) err: " + err.Error()}, nil
		} else {
			if b, err := s2c.Marshal(); err != nil {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_S2C_PROTO, Msg: "S2CBuyProduct.Marshal() err: " + err.Error()}, nil
			} else {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_SUCCESS, Proto: b}, nil
			}
		}
	}
}
