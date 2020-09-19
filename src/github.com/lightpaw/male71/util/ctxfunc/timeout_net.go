package ctxfunc

import (
	"time"
	"golang.org/x/net/context"
	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type NetFunc func(ctx context.Context) (err error)

func NetTimeout(d time.Duration, f Func) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	return f(ctx)
}

func NetTimeout1s(f Func) error {
	return NetTimeout(1*time.Second, f)
}

func NetTimeout2s(f Func) error {
	return NetTimeout(2*time.Second, f)
}

func NetTimeout3s(f Func) error {
	return NetTimeout(3*time.Second, f)
}

func NetTimeout10s(f Func) error {
	return NetTimeout(10*time.Second, f)
}

func NetTimeout30s(f Func) error {
	return NetTimeout(30*time.Second, f)
}

func NetTimeout1m(f Func) error {
	return NetTimeout(60*time.Second, f)
}

func IsRpcTimeout(err error) bool {
	err = errors.Cause(err)
	if code := status.Code(err); code == codes.DeadlineExceeded {
		return true
	}
	return err == context.DeadlineExceeded
}
