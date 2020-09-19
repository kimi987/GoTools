package msg

import "github.com/lightpaw/pbutil"

type ErrMsg interface {
	error
	ErrMsg() pbutil.Buffer
}
