package face

import (
	"github.com/lightpaw/pbutil"
)

type Func func()

type BFunc func() bool

type MsgFunc func() pbutil.Buffer
