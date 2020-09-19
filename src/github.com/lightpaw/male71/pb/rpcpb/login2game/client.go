
package login2game

import "golang.org/x/net/context"

type client interface {
	HandleBytes(ctx context.Context, handler, version string, key int64, proto []byte) ([]byte, error)
}
