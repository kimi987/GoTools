package game2login

import (
	"github.com/lightpaw/rpc7"
	"golang.org/x/net/context"
)

func NewVerifyLoginTokenHandler(f func(r *C2SVerifyLoginTokenProto) (*S2CVerifyLoginTokenProto, error)) (string, rpc7.HandlerFunc) {
	return "verify_login_token", func(ctx context.Context, r *rpc7.RpcRequest) (*rpc7.RpcResponse, error) {
		c2s := &C2SVerifyLoginTokenProto{}
		if err := c2s.Unmarshal(r.Proto); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_C2S_PROTO, Msg: "C2SVerifyLoginToken.Unmarshal() err: " + err.Error()}, nil
		}

		if s2c, err := f(c2s); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_HANDLE_ERROR, Msg: "VerifyLoginToken(r) err: " + err.Error()}, nil
		} else {
			if b, err := s2c.Marshal(); err != nil {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_S2C_PROTO, Msg: "S2CVerifyLoginToken.Marshal() err: " + err.Error()}, nil
			} else {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_SUCCESS, Proto: b}, nil
			}
		}
	}
}

func NewWriteTlogHandler(f func(r *C2SWriteTlogProto) (*S2CWriteTlogProto, error)) (string, rpc7.HandlerFunc) {
	return "write_tlog", func(ctx context.Context, r *rpc7.RpcRequest) (*rpc7.RpcResponse, error) {
		c2s := &C2SWriteTlogProto{}
		if err := c2s.Unmarshal(r.Proto); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_C2S_PROTO, Msg: "C2SWriteTlog.Unmarshal() err: " + err.Error()}, nil
		}

		if s2c, err := f(c2s); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_HANDLE_ERROR, Msg: "WriteTlog(r) err: " + err.Error()}, nil
		} else {
			if b, err := s2c.Marshal(); err != nil {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_S2C_PROTO, Msg: "S2CWriteTlog.Marshal() err: " + err.Error()}, nil
			} else {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_SUCCESS, Proto: b}, nil
			}
		}
	}
}

func NewPushHandler(f func(r *C2SPushProto) (*S2CPushProto, error)) (string, rpc7.HandlerFunc) {
	return "push", func(ctx context.Context, r *rpc7.RpcRequest) (*rpc7.RpcResponse, error) {
		c2s := &C2SPushProto{}
		if err := c2s.Unmarshal(r.Proto); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_C2S_PROTO, Msg: "C2SPush.Unmarshal() err: " + err.Error()}, nil
		}

		if s2c, err := f(c2s); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_HANDLE_ERROR, Msg: "Push(r) err: " + err.Error()}, nil
		} else {
			if b, err := s2c.Marshal(); err != nil {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_S2C_PROTO, Msg: "S2CPush.Marshal() err: " + err.Error()}, nil
			} else {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_SUCCESS, Proto: b}, nil
			}
		}
	}
}

func NewPushMultiHandler(f func(r *C2SPushMultiProto) (*S2CPushMultiProto, error)) (string, rpc7.HandlerFunc) {
	return "push_multi", func(ctx context.Context, r *rpc7.RpcRequest) (*rpc7.RpcResponse, error) {
		c2s := &C2SPushMultiProto{}
		if err := c2s.Unmarshal(r.Proto); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_C2S_PROTO, Msg: "C2SPushMulti.Unmarshal() err: " + err.Error()}, nil
		}

		if s2c, err := f(c2s); err != nil {
			return &rpc7.RpcResponse{Code: rpc7.RspCode_HANDLE_ERROR, Msg: "PushMulti(r) err: " + err.Error()}, nil
		} else {
			if b, err := s2c.Marshal(); err != nil {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_BAD_S2C_PROTO, Msg: "S2CPushMulti.Marshal() err: " + err.Error()}, nil
			} else {
				return &rpc7.RpcResponse{Code: rpc7.RspCode_SUCCESS, Proto: b}, nil
			}
		}
	}
}
