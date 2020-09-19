package rpc7

import (
	"sync"
	"golang.org/x/net/context"
)

// NewHandlerMux allocates and returns a new HandlerMux.
func NewHandlerMux() *HandlerMux { return new(HandlerMux) }

// DefaultHandlerMux is the default HandlerMux used by Serve.
var DefaultHandlerMux = &defaultHandlerMux

var defaultHandlerMux HandlerMux

type Handler interface {
	Handle(ctx context.Context, r *RpcRequest) (*RpcResponse, error)
}

type HandlerFunc func(ctx context.Context, r *RpcRequest) (*RpcResponse, error)

func (f HandlerFunc) Handle(ctx context.Context, r *RpcRequest) (*RpcResponse, error) {
	return f(ctx, r)
}

type HandlerMux struct {
	mu              sync.RWMutex
	m               map[string]Handler
	beforeHandler   Handler // handler 之前执行
	notFoundHandler Handler // 没有匹配Handler的处理器
}

func HandleFunc(pattern string, handler func(ctx context.Context, r *RpcRequest) (*RpcResponse, error)) {
	DefaultHandlerMux.HandleFunc(pattern, handler)
}

func (mux *HandlerMux) HandleFunc(pattern string, handler func(ctx context.Context, r *RpcRequest) (*RpcResponse, error)) {
	mux.Handle(pattern, HandlerFunc(handler))
}

func Handle(pattern string, handler Handler) {
	DefaultHandlerMux.Handle(pattern, handler)
}

func (mux *HandlerMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	if pattern == "" {
		panic("rpc7: invalid pattern " + pattern)
	}
	if handler == nil {
		panic("rpc7: nil handler")
	}
	if mux.m == nil {
		mux.m = make(map[string]Handler)
	}
	if mux.m[pattern] != nil {
		panic("rpc7: multiple registrations for " + pattern)
	}

	mux.m[pattern] = handler
}

func BeforeHandleFunc(handler func(ctx context.Context, r *RpcRequest) (*RpcResponse, error)) {
	DefaultHandlerMux.BeforeHandleFunc(handler)
}

func (mux *HandlerMux) BeforeHandleFunc(handler func(ctx context.Context, r *RpcRequest) (*RpcResponse, error)) {
	mux.BeforeHandle(HandlerFunc(handler))
}

func BeforeHandle(handler Handler) {
	DefaultHandlerMux.BeforeHandle(handler)
}

func (mux *HandlerMux) BeforeHandle(handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if mux.beforeHandler != nil {
		panic("rpc7: multiple set before handler")
	}

	mux.beforeHandler = handler
}

func NotFoundHandleFunc(handler func(ctx context.Context, r *RpcRequest) (*RpcResponse, error)) {
	DefaultHandlerMux.NotFoundHandleFunc(handler)
}

func (mux *HandlerMux) NotFoundHandleFunc(handler func(ctx context.Context, r *RpcRequest) (*RpcResponse, error)) {
	mux.NotFoundHandle(HandlerFunc(handler))
}

func NotFoundHandle(handler Handler) {
	DefaultHandlerMux.NotFoundHandle(handler)
}

func (mux *HandlerMux) NotFoundHandle(handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if mux.notFoundHandler != nil {
		panic("rpc7: multiple set notfound handler")
	}

	mux.notFoundHandler = handler
}

func (m *HandlerMux) NewHandler() RpcServiceServer {
	m.mu.Lock()
	defer m.mu.Unlock()

	funcMap := make(map[string]Handler, len(m.m))
	for k, v := range m.m {
		funcMap[k] = v
	}

	return &RpcHandler{
		beforeHandler:   m.beforeHandler,
		m:               funcMap,
		notFoundHandler: m.notFoundHandler,
	}
}

type RpcHandler struct {
	m               map[string]Handler
	beforeHandler   Handler // handler 之前执行
	notFoundHandler Handler // 没有匹配Handler的处理器
}

var Continue = &RpcResponse{Code: RspCode_CONTINUE}

func (h *RpcHandler) Handle(ctx context.Context, r *RpcRequest) (*RpcResponse, error) {

	if h.beforeHandler != nil {
		if resp, err := h.beforeHandler.Handle(ctx, r); err != nil {
			return &RpcResponse{
				Code: RspCode_HANDLE_ERROR,
				Msg:  "before handler error, err: " + err.Error(),
			}, nil
		} else if resp != nil && resp.Code != RspCode_CONTINUE {
			return resp, nil
		}
	}

	handler := h.m[r.Handler]
	if handler == nil {
		if handler = h.notFoundHandler; handler == nil {
			return &RpcResponse{
				Code: RspCode_HANDLER_NOT_FOUND,
				Msg:  "Handler not found: " + r.Handler,
			}, nil
		}
	}

	// TODO Version 暂不处理
	return handler.Handle(ctx, r)
}

var (
	servingResp = &HealthCheckResponse{
		Status: HealthCheckResponse_SERVING,
	}
	notServingResp = &HealthCheckResponse{
		Status: HealthCheckResponse_NOT_SERVING,
	}
)

func (h *RpcHandler) Check(ctx context.Context, r *HealthCheckRequest) (*HealthCheckResponse, error) {
	if r.Service != "" && h.m[r.Service] == nil {
		return notServingResp, nil
	} else {
		return servingResp, nil
	}
}
