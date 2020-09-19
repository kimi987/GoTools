package survey

import (
	"github.com/lightpaw/male7/gen/pb/util"
	"github.com/lightpaw/male7/pb/shared_proto"
	"github.com/lightpaw/pbutil"
)

var (
	pool           = pbutil.Pool
	newProtoMsg    = util.NewProtoMsg
	newCompressMsg = util.NewCompressMsg
	safeMarshal    = util.SafeMarshal
	_              = shared_proto.ErrIntOverflowConfig
)

type marshaler util.Marshaler

const (
	MODULE_ID = 37

	C2S_COMPLETE = 2
)

func NewS2cCompleteMsg(id string) pbutil.Buffer {
	msg := &S2CCompleteProto{
		Id: id,
	}
	return NewS2cCompleteProtoMsg(msg)
}

var s2c_complete = [...]byte{37, 1} // 1
func NewS2cCompleteProtoMsg(object *S2CCompleteProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_complete[:], "s2c_complete")

}
