package rpc7

import (
	"google.golang.org/grpc"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "did not connect: %v", address)
	}

	return &Client{
		conn: conn,
		rsc:  NewRpcServiceClient(conn),
	}, nil
}

type Client struct {
	conn *grpc.ClientConn
	rsc  RpcServiceClient
}

func (c *Client) GetRaw() RpcServiceClient {
	return c.rsc
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Handle(ctx context.Context, handler, version string, key int64, proto marshaler) ([]byte, error) {
	if b, err := proto.Marshal(); err != nil {
		return nil, errors.Wrapf(err, "proto.Marshal fail")
	} else {
		return c.HandleBytes(ctx, handler, version, key, b)
	}
}

func (c *Client) HandleBytes(ctx context.Context, handler, version string, key int64, proto []byte) ([]byte, error) {

	resp, err := c.HandleRequest(ctx, &RpcRequest{
		Handler: handler,
		Version: version,
		Key:     key,
		Proto:   proto,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "rpc.Handle fail")
	}

	if resp.Code != RspCode_SUCCESS {
		return resp.Proto, &withCode{code: resp.Code, msg: resp.Msg}
	}
	return resp.Proto, nil
}

func (c *Client) HandleRequest(ctx context.Context, r *RpcRequest) (*RpcResponse, error) {
	return c.rsc.Handle(ctx, r)
}

func (c *Client) Check(ctx context.Context, r *HealthCheckRequest) (*HealthCheckResponse, error) {
	return c.rsc.Check(ctx, r)
}

var aliveRequest = &HealthCheckRequest{}

func (c *Client) CheckAlive(ctx context.Context) (error) {
	_, err := c.rsc.Check(ctx, aliveRequest)
	return err
}

type marshaler interface {
	Marshal() ([]byte, error)
}

type withCode struct {
	code RspCode
	msg  string
}

func (w *withCode) Code() RspCode { return w.code }
func (w *withCode) Error() string { return w.code.String() + ": " + w.msg }

func ErrorCode(err error) RspCode {
	type coder interface {
		Code() RspCode
	}

	if err != nil {
		if cause, ok := err.(coder); ok {
			return cause.Code()
		}
	}
	return RspCode_INVALID_CODE
}
