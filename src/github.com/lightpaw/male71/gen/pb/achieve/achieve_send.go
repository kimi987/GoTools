package achieve

import (
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/pbutil"
)

var (
	pool           = pbutil.Pool
	newProtoMsg    = util.NewProtoMsg
	newCompressMsg = util.NewCompressMsg
	safeMarshal    = util.SafeMarshal
)

type marshaler util.Marshaler

const (
	MODULE_ID = 28
)
