package stress

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
	MODULE_ID = 10

	C2S_ROBOT_PING = 1
)

func NewS2cRobotPingMsg(time int32) pbutil.Buffer {
	msg := &S2CRobotPingProto{
		Time: time,
	}
	return NewS2cRobotPingProtoMsg(msg)
}

var s2c_robot_ping = [...]byte{10, 2} // 2
func NewS2cRobotPingProtoMsg(object *S2CRobotPingProto) pbutil.Buffer {

	return newProtoMsg(object, s2c_robot_ping[:], "s2c_robot_ping")

}
